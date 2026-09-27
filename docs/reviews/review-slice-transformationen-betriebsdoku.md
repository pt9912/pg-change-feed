# Review-Report: slice-transformationen-betriebsdoku — 2026-09-27

**Review-Art:** Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `git diff b80f45cf..97de226c` (Feature-Commit `97de226c`, letzter Slice von
`welle-transformationen`).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-27 (geschärft 2026-09-09, seit
slice-harness-mutationsbild-und-verweigerte-aktion erweitert)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-transformationen-betriebsdoku.md` (Plan-Datei, inkl. der
  beiden „Übergabe aus …“-Blöcke §2)
- `spec/lastenheft.md` — `LH-FA-CFG-007`, `LH-FA-ADM-001`, `LH-FA-SST-009`
- `spec/pflichtenheft.md` — `SPEC-019` (Antrags-Queue, K1–K4, Zeilen ohne Antrag),
  `SPEC-030` (Regelform-Beispiele)
- `docs/plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md`,
  `0125-transformationen-parametertyp-regelform-json.md`,
  `0126-transformationen-annahmemenge-rule-spec.md`,
  `0117-backfill-run-fehlerklasse-schema.md`,
  `0127-antrags-queue-requested-at-aufrufzeitpunkt.md`,
  `0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md`
- `docs/reviews/review-slice-transformationen-antragsweg-schema.md` (zitierte Messung des
  `::jsonb`-Fehlertexts)
- `docs/reviews/architect-verdict-welle-transformationen-offene-fragen.md` (Quelle der
  übernommenen/gemessenen Zahlen im neuen Abschnitt)
- `AGENTS.md` §3.7, §3.9, §3.11, §3.12, §3.13

---

## Findings

### F-1 — Zitierter PostgreSQL-Fehlertext stimmt nicht mit dem genannten Beleg überein

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz B / Reviewer-Skill HIGH „Beleg trägt seinen Satz nicht“
- `pfad`: `docs/user/benutzerhandbuch.md:337`
- `befund`: Der neue Abschnitt zitiert als „gemessen im Review-Report“ (Link auf
  `docs/reviews/review-slice-transformationen-antragsweg-schema.md`) den Wortlaut
  `` `function cdc.set_transformation(text, text, text, text, jsonb) does not exist` ``. Der
  zitierte Review-Report misst an dieser Stelle jedoch einen anderen Wortlaut:
  `„function cdc.set_transformation(unknown, unknown, unknown, unknown, jsonb) does not exist“`
  (Zeile 102 des Reports — die vier Textliteral-Parameter erscheinen dort als `unknown`, nicht als
  `text`, weil PostgreSQL nicht gecastete String-Literale beim Signatur-Fehler als `unknown`
  typisiert). Auch `ADR-0125` Festlegung 1 — die zweite genannte Quelle — führt an dieser Stelle nur
  die abgekürzte Form „function … does not exist“, nie den ausgeschriebenen Parametersatz. Eine
  repoweite Suche (`git grep -n "text, text, text, text, jsonb"`) findet die im Handbuch verwendete
  Zeichenkette **nirgends sonst** im Repository — auch nicht als Variante mit `json` statt `jsonb`,
  die an einer völlig anderen Stelle (`ADR-0125` §Gemessen, d-migrate-Rollout-Fehler beim
  `DROP FUNCTION`, nicht der Laufzeit-Aufruf) tatsächlich vorkommt. Ein Betreiber, der diesen realen
  Fehler erhält und den im Handbuch zitierten Wortlaut zum Vergleich oder zur Suche heranzieht, sieht
  einen anderen Text als die zitierte Handbuch-Zeile.
- `verifizierbar`: ja — `SELECT cdc.set_transformation('s','public','t','r','{"kind":"x"}'::jsonb)`
  gegen eine Wegwerf-Instanz mit dem aktuellen Schema-Rollout liefert den tatsächlichen Wortlaut;
  `git grep -n "text, text, text, text, jsonb" -- '*.md'` zeigt, dass die Zeichenkette nur in der
  neuen Handbuch-Zeile auftaucht.
- `klasse`: Beleg trägt seinen Satz nicht (Herkunft: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`)

### F-2 — Übernahme der Übergabe „Zeilen, die kein Antrag sind“ ist unvollständig

- `kategorie`: MEDIUM
- `quelle`: Plan-Datei §2 DoD-Punkt 2 (Übernahme aus `antragsqueue-lesefehler-failed`), `SPEC-019`
  Absatz „Zeilen, die kein Antrag sind“
- `pfad`: `docs/user/benutzerhandbuch.md:118-124`
- `befund`: Der Plan verlangt für diesen Absatz ausdrücklich drei Beispiele einer vom
  Antrags-Konstruktor verworfenen `pending`-Zeile: „leeres Schema, leerer Tabellenname,
  `exclude_column`/`include_column` mit leerer Spalte“ (Plan §2, Übergabe aus
  `antragsqueue-lesefehler-failed`, belegt durch den Store-Test
  `TestAdministrationRequestListPendingPassesRejectedRowsThrough`, der genau diese drei Fälle
  prüft). Der geschriebene Handbuch-Absatz nennt nur zwei: „leeres Schema, leerer Tabellenname“ —
  das dritte Beispiel (leere Spalte bei `exclude_column`/`include_column`) fehlt vollständig, obwohl
  der einleitende Satz „Quelle, Schema, Tabelle **und Spalte**“ nennt und damit selbst ein
  Spalten-Beispiel ankündigt. `SPEC-019`s eigener Absatz „Zeilen, die kein Antrag sind“ führt
  zusätzlich „leere Quelle“ und „unbekannte Antragsart“ als Gründe — die der Plan bewusst vom
  Handbuch ausschließt (Fremdschlüssel/CHECK, nur am Fake-Test belegt); dieser Ausschluss ist im
  Plan begründet und kein Fund. Das fehlende Spalten-Beispiel ist dagegen nicht begründet
  weggelassen, sondern schlicht nicht übernommen.
- `verifizierbar`: ja — Textvergleich Plan §2 gegen den geschriebenen Absatz;
  `TestAdministrationRequestListPendingPassesRejectedRowsThrough` (Zeile 1439-1442) zeigt die drei
  geprüften Fälle.
- `klasse`: Nachzug lässt überholten/unvollständigen Text stehen (verwandt zu
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, hier: unvollständige statt widersprechende
  Übernahme)

## Negativbefunde

- geprüft, ohne Befund: SDK-Suchlauf-Nichtbefund (§3 der Plan-Datei) — alle 137 Treffer in
  `sdks/**/*.{py,cs,kt}` selbst gelesen (nicht nur stichprobenartig); jede Produktivcode-Fundstelle
  (`sdks/csharp/PgChangeFeed.Client/{Http,Sse,Nats}/Models/*.cs`,
  `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../{http,sse,nats}/model/*.kt` +
  `grpc/PgChangeFeedGrpcClient.kt`, `sdks/python/pgchangefeed/src/pgchangefeed/{models,grpc_client,sse_client}.py`)
  behandelt `old_image`/`new_image` tatsächlich als opakes oberstes Feld (`JsonElement?`,
  `Any | None`, `bytes`); kein Zugriff auf einen benannten Schlüssel innerhalb eines Row Image. Alle
  übrigen Treffer liegen in Tests/Integrationstests (Fixture-Daten) oder Docstrings. `make
  suchlauf-nachmessen` bestätigt 137/137 an beiden Ständen.
- geprüft, ohne Befund: `sdks/csharp/README.md`, `sdks/python/README.md`,
  `sdks/kotlin/pgchangefeed-kotlin/README.md` tragen bereits je einen `## The change object`-Abschnitt
  — die Behauptung „keine Lücke“ ist zutreffend.
- geprüft, ohne Befund: numerische Werte des neuen Abschnitts (ns/µs/ms-Bereiche der `map_value`-Suche,
  Prozentangaben der CPU-Last, Backfill-Lesekosten, 30-s-Frist, 31-s-Messung, 7,1-s-Heartbeat-Alter)
  gegen `architect-verdict-welle-transformationen-offene-fragen.md` §1/§3/§6 und
  `docs/plan/planning/done/slice-start-vorlauf-grenze.md` abgeglichen — jede Zahl trägt eine korrekte
  gemessen/übernommen/abgeleitet-Kennzeichnung und stimmt mit der genannten Quelle überein.
- geprüft, ohne Befund: K1–K4-Fehlertexte der Konfliktfreiheits-Tabelle gegen
  `internal/domain/errors/errors.go`, `internal/application/port/inbound/verwaltung.go` und
  `spec/pflichtenheft.md` §K1–K4 — Wortlaut exakt identisch.
- geprüft, ohne Befund: Aufrufform-Absatz (`json` vs. `jsonb`-Parameter, `\u0000`, Zahlbereich
  `numeric`, Stapeltiefe, PostgreSQL 17/18) gegen `ADR-0125`/`ADR-0126` — inhaltlich korrekt
  wiedergegeben.
- geprüft, ohne Befund: Backfill-Bezug/Nichtanwendbarkeit im Run (Klasse `schema` vor der ersten
  Kopie, `configuration` bei Regelstand-Wechsel während des Runs) gegen `ADR-0117` — korrekt.
- geprüft, ohne Befund: „Reihenfolge der Aufrufe“ gegen `ADR-0127` Folgepflicht 4 — Wortlaut fast
  identisch übernommen; die zugehörige Fitness Function
  (`TestAdministrationRequestListPendingOrdersTiesByRequestID`,
  `internal/adapters/driven/postgresstorage/administrationrequest_order_test.go`) existiert bereits
  im Bestand und deckt die Aussage, wird im Handbuch-Text aber nicht namentlich zitiert (nur der
  ADR-Anker) — keine eigene Kategorie, da der ADR-Anker als Beleg-Anker im Sinn von `AGENTS.md`
  §3.12 ausreicht.
- geprüft, ohne Befund: Abgrenzung „kein allgemeiner Schema-Fehler-Recovery-Weg“ (Plan §1) — der
  Nichtanwendbarkeits-/Abhilfe-Absatz bleibt eng auf die Transformationsregel bezogen; „Container
  startet nicht“ und „Neustart nach einem Fehler“ bleiben unverändert und verweisen nicht auf die
  neue Abhilfe.
- geprüft, ohne Befund: §3.7 (keine Chronik-/Vorher-Nachher-Sprache außerhalb der
  Änderungshistorie-Tabelle) im gesamten neuen Abschnitt.
- geprüft, ohne Befund: §3.11 (keine host-lokalen absoluten Pfade) im gesamten neuen Abschnitt.
- geprüft, ohne Befund: Handbuch-Versionshistorie — `Version:`-Kopf 1.66→1.67 und neue Zeile in
  `### Änderungshistorie` im selben Diff vorhanden.
- geprüft, ohne Befund: `make gates` (real ausgeführt, Exit-Code direkt geprüft, EXIT:0),
  `make docs-check` (EXIT:0), `make suchlauf-nachmessen PLAN=…` (8/8 „stimmen“),
  `make sdk-public-doc-check` (EXIT:0, keine interne Kennung unter `sdks/`) — alle vier real
  nachgefahren, nicht nur dem Bericht geglaubt.
- geprüft, ohne Befund: `make fmt-check`/`make kommentar-kennungen` zu Recht nicht gelaufen — der
  Diff berührt keine `.go`- oder `tools/schema/*.sql`-Datei (`git diff --stat` bestätigt).
- geprüft, ohne Befund: DoD-Checkbox-Zustände der Plan-Datei — nur Implementer-eigene Punkte
  (SDK-Beleg, Betriebsdokumentation, Zahlen-Ursprung, `make gates`, §3.13-Suchlauf, Doku-Update,
  Reconciliation-Register entfällt) sind `[x]`; Review/Closure/Verifikations-Punkte bleiben `[ ]`.
- geprüft, ohne Befund: Traceability — Plan-Datei und Handbuch-Änderungshistorie nennen `LH-*`- und
  `ADR-*`-Kennungen durchgehend.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Nachzug lässt
überholten/unvollständigen Text stehen

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) zitiert einen konkreten, für Betreiber nachschlagbaren
Fehlertext falsch; das gehört korrigiert, bevor das Handbuch als Betriebsdokumentation gilt. F-2
(MEDIUM) ist eine benannte Lücke in einer explizit im Plan geforderten Übernahme.

**Übergabe:** Beide Findings gehen an den Implementer zur Fixrunde (Reviewer → Implementer). Die
DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt vor“ bleibt bewusst `[ ]` —
es kommt eine Fixrunde, der reguläre Nachzug läuft über Schritt 21 des Implementer-Workflows nach
Abschluss der Fixrunde (Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde: die dortige Ausnahme
gilt nur ohne Fixrunde).
