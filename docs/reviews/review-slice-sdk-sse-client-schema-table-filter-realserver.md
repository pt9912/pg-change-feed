# Review-Report: slice-sdk-sse-client-schema-table-filter-realserver — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice [sdk-sse-client-schema-table-filter-realserver](../plan/planning/in-progress/slice-sdk-sse-client-schema-table-filter-realserver.md),
Implementer-Commit `8b0e1e52` (Parent der Plan-Anlage: `448ee4a5`; die Commits `40371908`, `de905535`, `d79086d1` sind Planner-Commits ohne Code).
Diff `git diff 448ee4a5 HEAD`: 13 Dateien, +1191/−14 (drei Testdateien, zwei `PhaseEnvironment`-Hilfen, die neue Hilfsdatei
`tools/harness/lib-sdk-filter-fixture.sh`, drei Runner, `harness/README.md`, `harness/mk/sdk.mk`, drei Abdeckungs-Zeilen, Plan).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage und Arbeitsweise:** Der Reviewer-Lauf hat diesen Report mit dem Write-Werkzeug geschrieben (kein Edit-Werkzeug im Lauf). Mutationen liefen
in einer Kopie des Arbeitsbaums im Scratchpad (`git archive HEAD` nach `rvf/base`, je Mutation eine Kopie; Mutation per `sed … > Datei im Scratchpad`
und `cp` in die Kopie, kein `sed -i`); die Läufe der Mutationen fuhren `make test-sdk-<sprache>-integration` in der Kopie.
`git status --short` im echten Repo war nach den Läufen leer. `make image` wurde nicht verweigert und lief einmal vor den Tier-Läufen
(Exit 0, `:dev` war als `ghcr.io/pt9912/pg-change-feed:dev` geladen; der Server-Code ist durch den Slice unverändert, `git diff 448ee4a5 HEAD --stat -- internal cmd proto gen` leer);
die Mutationsläufe nutzten dasselbe `:dev`. Der Lauf-Beleg des Images ist der Digest `sha256:48e5699d9b2fbf9761fa1f853d6d2faa6b0ba3e43007e81b320e8a57d466a878`
(`harness/image-hash.txt`, lokal, nicht committet, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)).

**Eingangs-Kontext:**

- Slice-Plan `sdk-sse-client-schema-table-filter-realserver` (§1 Abgrenzung, §2 DoD, §3 Plan und Suchlauf-Feld, §6 Risiken)
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Teilfrage 1 und 4),
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Festlegung 2),
  [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
- [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.7, §3.9, §3.12, §3.13, §3.15), [`harness/conventions.md`](../../harness/conventions.md)
- Mechanik-Vorbild und vorheriger Report am selben Modul: [`review-slice-routing-sdk-realserver-e2e`](review-slice-routing-sdk-realserver-e2e.md),
  [`review-slice-sdk-sse-client-schema-table-filter`](review-slice-sdk-sse-client-schema-table-filter.md)

**Eigene Messungen** (Exit-Codes je als eigener Schritt ausgewertet; „gefahren“ nur für Selbstgefahrenes):

- `make fmt-check` Exit 0, gedruckt „fmt-check: 323 Go-Dateien geprüft, alle formatiert“ (prüft nur Go; die Kotlin-Datei ist davon nicht gedeckt, siehe F-1).
- `make sdk-public-doc-check` Exit 0, gedruckt „sdk-public-doc-check: keine interne Kennung unter sdks“.
- `make docs-check` Exit 0 vor Anlage dieses Reports, gedruckt „d-check: 1548 Datei(en) geprüft, 0 Befund(e)“ (nach Anlage siehe Verdikt).
- `make test` Exit 0 (kein Go-Code im Diff; Lauf als Gegenprobe).
- `make kommentar-kennungen DIFF=448ee4a5` Exit 0 (kein Kandidat). Der richtige Parent ist `448ee4a5`: die Plan-Commits berühren keinen Code, der Lauf mit
  `DIFF=d79086d1` liefert dasselbe Ergebnis (Exit 0); der Unterschied ist für diesen Diff ohne Wirkung.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-sse-client-schema-table-filter-realserver.md` Exit 0, gedruckt „suchlauf-nachmessen: 12 Zeilen stimmen“.
- `git diff 448ee4a5 HEAD --name-only -- sdks`: fünf Pfade, alle unter `PgChangeFeed.Client.Integration/`, `src/integrationTest/` bzw. `integration/`; keine `.csproj`/`pyproject.toml`/`build.gradle.kts`.
- `git grep -nE "zwölf|12 Phasen"` über den Baum (ohne `docs/reviews`, Baseline, `done/`): keine Phasenzahl-Nennung der Tiers außerhalb des Plans (Roadmap „zwölf Zustellweg-Flächen“ zählt 3 × 4 Flächen, nicht Phasen).
  „dreizehn“ steht in `harness/mk/sdk.mk` (3 Hilfetexte) und in den Köpfen der drei Runner; `harness/README.md` nennt keine Phasenzahl, nur die Filter-Phase.
- `docs/user/sdk-e2e-abdeckung.md`: `git diff --stat` zeigt 3 Einfügungen (je Runner eine Zeile); die drei unmutierten Läufe haben die Datei nicht erneut verändert (`git status --short` danach leer).

**Unmutierte Tier-Läufe** (jeder mit eigenem Exit-Code; gedruckte Zeile ist der Beleg, nicht der Bau):

| Tier | Exit | gedruckte Zeile |
|---|---|---|
| `make test-sdk-csharp-integration` | 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15`, „SEEN nach 376ms ab dem letzten Commit (Versuch 1)“ |
| `make test-sdk-kotlin-integration` | 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |
| `make test-sdk-python-integration` | 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |

Die Alt-Phasen der drei Läufe liefen mit unveränderter Erwartung grün (Exit 0 der Gesamtläufe); der in der Übergabe gemeldete einmalige Flake der NATS-Phase trat in meinem einen C#-Lauf nicht auf (siehe F-4).

**Mutationsproben** (Eingabeseite des Tests bzw. des Packages, Instanz: realer Tier-Lauf in einer Kopie; jede Farbe mit der gedruckten Zeile):

| Nr. | Sprache | Stelle der Mutation | Farbe | gedruckte Zeile |
|---|---|---|---|---|
| M1 | C# | F1 ohne `table:` (`SseFilterRealserverTests.cs:87`) | rot (Exit 2, Runner: „fremde Changes am Filter-Client“) | `FILTER_RESULT f1=2 f1_foreign=1 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |
| M2 | C# | F2 ohne `schema:` (`SseFilterRealserverTests.cs:90`) | rot (Exit 2) | `FILTER_RESULT f1=1 f1_foreign=0 f2=3 f2_foreign=2 unfiltered=3 quiet_seconds=15` |
| M2 | Kotlin | F2 ohne `schema =` (`SseFilterRealserverTest.kt:80`) | rot (Exit 2) | `FILTER_RESULT f1=1 f1_foreign=0 f2=3 f2_foreign=2 unfiltered=3 quiet_seconds=15` |
| M1 | Kotlin | F1 ohne `table =` (`SseFilterRealserverTest.kt:76`) | rot (Exit 2) | `FILTER_RESULT f1=2 f1_foreign=1 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |
| M1 | Python | F1 mit `table=None` (`test_sse_filter_realserver.py:141`) | rot (Exit 2) | `FILTER_RESULT f1=2 f1_foreign=1 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15` |
| P1 | Python (Package) | `sse_client.py:90`: `("table", table)` aus der Query-Menge entfernt | rot (Exit 2), aber **im Bau-Schritt `pytest`** der Stufe `build` (`tests/test_sse_client.py:158` u. a.), nicht in der Realserver-Phase; keine `FILTER_RESULT`-Zeile | — |

Nicht gefahren: M3 (je Sprache) und die Package-Mutationen für C# und Kotlin; die Aussage zu diesen ist **hergeleitet**, nicht erprobt (siehe F-5).
Die Fremdwerte der Mutationen bestätigen die Zählung am Empfang: M1 → B (ein Fremdwert; die Tabelle des zweiten Schemas gehört nicht zu `public`), M2 → A und B (zwei Fremdwerte).

---

## Findings

### F-1 — Fehlendes Leerzeichen in der Kotlin-Warteschleife

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/SseFilterRealserverTest.kt:87`
- `befund`: Die Bedingung trägt `&&hasAllThree(unfiltered.rows)` ohne Leerzeichen hinter dem Operator; die Zeile ist zudem die längste der Datei. `make fmt-check` liest nur Go und deckt die Datei nicht.
- `verifizierbar`: nein — kein Sensor liest Kotlin-Formatierung.
- `klasse`: Formatierungs-Drift ohne Gate (Nicht-Go)

### F-2 — Negativ-Beleg setzt die Verbindung der gefilterten Clients vor dem ersten Commit voraus, ohne darauf zu warten

- `kategorie`: LOW
- `quelle`: [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4; Plan §6 „Negativ-Beleg über Abwesenheit“
- `pfad`: `sdks/csharp/PgChangeFeed.Client.Integration/SseFilterRealserverTests.cs:95`, `tools/harness/lib-sdk-filter-fixture.sh:114`
- `befund`: `READY` wird gedruckt, sobald die drei Konsumenten gestartet sind, nicht sobald ihre Verbindungen stehen (gleich in Kotlin und Python); der Runner committet danach die Gruppe (B, zweites Schema, A). Ein gefilterter Client, der erst zwischen den Commits verbindet, bleibt in dieser Gruppe ohne Fremdwert, weil ihn die fremde Change nie erreichte; `SEEN` verlangt von F1 nur eine Change der Tabelle A. In allen sechs gefahrenen Läufen fiel `SEEN` im Versuch 1 und die Mutationen färbten rot, die Lücke ist also eine Robustheitsfrage, kein gesehener Ausfall.
- `verifizierbar`: nein — ein Gate-Lauf würde die Reihenfolge nicht belegen; ein Zwischenlauf mit verzögerter Verbindung wäre die Probe, nicht gefahren.
- `klasse`: Negativ-Beleg ohne Nachweis der Vorbedingung

### F-3 — Die Zusage „genau seine Tabelle“ prüft Reinheit, nicht Vollständigkeit

- `kategorie`: INFO
- `quelle`: Plan §1 Ziel („empfängt … genau die Changes seiner Tabelle“)
- `pfad`: `tools/harness/lib-sdk-filter-fixture.sh:202`, `:207`, `:208`
- `befund`: Der Runner verlangt `RECEIVED_F1`/`RECEIVED_F2`-Zeilen in der Zahl des Fremdwert-Zählers und mindestens eine Zeile, vergleicht sie aber nicht mit der Zahl der committeten Changes der Auswahl; eine ausfallende Change der eigenen Tabelle bliebe bei Zahl ≥ 1 unbemerkt. Die Gegenlesung ist unabhängig (Schema und Tabelle kommen aus `cdc.changes`, nicht aus dem Client-Ergebnis; Treffer verlangt zusätzlich den Sentinel der Phase), der Fremdwert-Zähler selbst stammt vom Client.
- `verifizierbar`: nein
- `klasse`: Zusage breiter als ihre Messung (Reinheit statt Vollständigkeit)

### F-4 — Gemeldeter Flake der alten NATS-Phase: Ursache hergeleitet, kein Zusammenhang mit der neuen Phase erkennbar

- `kategorie`: INFO
- `quelle`: Maintainability; [`AGENTS.md`](../../AGENTS.md) §3.12 (Ursprung: übernommen)
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh:242`, `sdks/csharp/PgChangeFeed.Client.Integration/NatsRealserverTests.cs:55`
- `befund`: Die Filter-Phase läuft im Runner nach allen Alt-Phasen in eigenem Container; sie kann die NATS-Phase kausal nicht beeinflussen. Lesbar ist ein enges Zeitfenster im unveränderten Code: der Runner sucht den Marker `REJECTED token-rejected` bis 20 s (40 × 0,5 s) nach dem Empfangsmarker, der Test bricht den Reject-Versuch erst mit `RejectCts` nach 15 s ab. Die Aussage zur Ursache ist **hergeleitet**; der Flake trat in meinem einen C#-Lauf nicht auf, ein Wiederholungslauf zur Messung fand nicht statt.
- `verifizierbar`: nein — ohne Wiederholungsläufe.
- `klasse`: Flake einer Alt-Phase ohne Bezug zum Diff

### F-5 — Package-Mutation am Realserver-Beleg ist für Python nicht von den Unit-Tests zu trennen

- `kategorie`: INFO
- `quelle`: Plan §2 „Mutationsproben je Sprache“ (Package-Mutation), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py:90` (Mutationsstelle), `sdks/python/Dockerfile` (Stufe `build` fährt `pytest` vor `integration`)
- `befund`: Die Package-Mutation P1 (`table` nicht auf den Draht) wird in der Stufe `build` von den Unit-Tests abgefangen, bevor die Realserver-Phase läuft; die Farbe stammt damit nicht vom Realserver-Beleg. Ob die Filter-Phase die Bindung des Packages für sich trägt, zeigt erst eine Package-Mutation, die zugleich die Unit-Erwartung mitnimmt; die Probe fuhr ich nicht, die Aussage ist hergeleitet. Der Plan nimmt die Package-Mutation je Sprache in die DoD; deren Beleg liegt beim Implementer-Bericht.
- `verifizierbar`: ja — durch eine Package-Mutation samt gespiegelter Unit-Erwartung und gedruckter `FILTER_RESULT`-Zeile.
- `klasse`: Beleg trägt seinen Satz nicht (Farbe aus anderer Stufe)

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/lib-sdk-filter-fixture.sh` — Wiederholbarkeit (jeder Lauf baut die Compose-Umgebung neu auf, `trap cleanup` mit `down -v`; `CREATE SCHEMA` ohne `IF NOT EXISTS` ist deshalb unkritisch), Namenskollision (`feed_e2e_sdkfilt_*` kommt im Baum nur in der neuen Datei vor; ID-Basis 1000 mit Schritt 3, je Tabelle eigener Schlüsselraum), Rollen (kein Rollenwechsel im Diff).
- geprüft, ohne Befund: Zeitfenster-Konstanten — 15 s Ruhefenster (derselbe Wert wie das Routing-Fenster), 5 Dreiergruppen à 12 s, Test-Frist 90 s; die Zustelldauer lag in meinem C#-Lauf bei 376 ms, das Ruhefenster beginnt erst nach `SEEN` des Clients ohne Filter.
- geprüft, ohne Befund: die drei Testdateien — F1/F2/U aufgebaut wie im Plan, `FILTER_RESULT` wird vor den Assertions gedruckt, die Fremdwerte zählen über alle empfangenen Zeilen (nicht nur die mit Sentinel).
- geprüft, ohne Befund: drei Runner — Phasenzahl „dreizehn“, `source`-Zeile, Sentinel je Sprache, Schlusszeile mit `SDK_FILTER_REPORT`, Abdeckungs-Zeile im marker-gegrenzten Abschnitt (idempotent: zweiter Lauf ohne Änderung).
- geprüft, ohne Befund: `harness/mk/sdk.mk`, `harness/README.md` §Sensors (drei Zeilen), Köpfe der Runner — Ist-Umfang stimmt mit den Runnern.
- geprüft, ohne Befund: `sdks/` — nur Test-Code und `PhaseEnvironment`-Hilfen, keine Version, kein Produktivcode, `make sdk-public-doc-check` Exit 0; keine Kennung in den neuen Testdateien.
- geprüft, ohne Befund: Kommentare (§3.7) — Indikativ, je Block höchstens eine Kennung; `make kommentar-kennungen DIFF=448ee4a5` ohne Kandidat.
- geprüft, ohne Befund: §3.1/§3.2/§3.15 — kein Host-Toolchain-Aufruf, keine Suppression, `make image` lief in meinem Lauf ohne Verweigerung; eine Verweigerung im Lauf des Implementers ist aus dem Diff nicht ablesbar.
- geprüft, ohne Befund: Träger der Phasenzahl im ganzen Baum (§3.13) — keine weitere Nennung gefunden, Suchlauf-Block stimmt an beiden Ständen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Formatierungs-Drift ohne Gate (Nicht-Go) · Negativ-Beleg ohne Nachweis der Vorbedingung · Zusage breiter als ihre Messung · Flake einer Alt-Phase ohne Bezug zum Diff · Beleg trägt seinen Satz nicht

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die Zusagen der drei Phasen sind an der Eingabeseite gebunden (M1/M2 in C#, Kotlin- und Python-Mutationen rot mit gedrucktem Fremdwert); die Gegenlesung über `cdc.changes` ist unabhängig vom Client-Ergebnis. Die beiden LOW-Funde (F-1, F-2) gehen an den Implementer ohne Fixrunde-Pflicht; F-3 bis F-5 sind Hinweise.

**Architect-Fragen:** keine zu klären; F-4 (Flake der Alt-Phase) bekommt einen eigenen Slice oder Beobachtungs-Eintrag, wenn er ein zweites Mal auftritt.

**Übergabe:** Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD-/Spec-Konformität prüft der Verifier). Die DoD-Zeile „Review durchgeführt“ ist im Plan im selben Commit nachgezogen ([`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md) §DoD-Checkbox-Nachzug), weil keine Fixrunde am Implementer nötig ist.
