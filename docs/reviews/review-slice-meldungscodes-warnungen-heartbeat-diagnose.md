# Review-Report: slice-meldungscodes-warnungen-heartbeat-diagnose — 2026-10-03

**Review-Art:** Code (gegen Plan, Entscheidung und Hard Rules; kein DoD-Abgleich, der gehört dem Verifier)

**Gegenstand:** `git diff a9767e87 HEAD` — Implementer-Commits `4e2a8f9b`, `82f6c931`, `5df69b02`; 53 Dateien, +1034/−195.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `5df69b02`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Ablage und Arbeitsweise:** Alle Mutationen liefen an `git archive HEAD`-Kopien im Scratchpad (Mutation per
`sed … > Kopie` und `mv` innerhalb der Kopie, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei);
`git status --short` im Echtrepo war nach den Läufen leer. Die Läufe unten liefen gegen den Arbeitsbaum,
der Exit-Code je Ziel wurde in einer eigenen Datei gesichert, nicht durch eine Pipe gelesen
([`AGENTS.md`](../../AGENTS.md) §3.9). Keine verweigerte Aktion im Lauf ([`AGENTS.md`](../../AGENTS.md) §3.15).

**Eingangs-Kontext:**

- Slice-Plan `slice-meldungscodes-warnungen-heartbeat-diagnose` (in-progress), §1 bis §3, §6
- [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegungen 1 bis 10, Verdikt
  `architect-verdict-meldungscodes-statt-interner-kennungen`
- [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md), [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
  [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- Vorgänger `slice-meldungscodes-registry-fehlerkopf` und `slice-meldungscodes-kennungsfreie-ausgaben` samt
  ihren Reviews (Lerneinträge)
- [`SPEC-008`](../../spec/pflichtenheft.md), [`SPEC-018`](../../spec/pflichtenheft.md), [`SPEC-031`](../../spec/pflichtenheft.md);
  [`LH-FA-ADM-003`](../../spec/lastenheft.md), [`LH-QA-REL-003`](../../spec/lastenheft.md), [`LH-QA-OPS-001`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.5, §3.6, §3.7, §3.9, §3.12, §3.13, §3.15, §4

---

## Läufe (selbst gefahren, Exit je Ziel gesichert)

| Ziel | Exit | gedruckte Zeile (Auszug) |
|---|---|---|
| `make test` | 0 | alle Pakete `ok` |
| `make test-store` | 0 | Paket `postgresstorage` `ok` |
| `make fmt-check` | 0 | — |
| `make a-check` | 0 | `gesamt: 0 Befund(e)` |
| `make generated-sync` | 0 | alle `gen/`-Dateien geprüft |
| `make meldungscodes-check` | 0 | `88 Codes in Tabelle und Katalog gleich, Quelltext (124 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make test-meldungscodes-check` | 0 | `73 Prüfungen bestanden` |
| `make ausgabe-kennungen-check` | 0 | — |
| `make handbuch-public-doc-check` | 0 | — |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make docs-check` | 0 | — |
| `make commit-traceability` | 0 | — |
| `make kommentar-kennungen DIFF=a9767e87` | 0 | kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | 0 | alle zwölf Zeilen stimmen |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | `Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf) … error_code text:YES … source_id,heartbeat_at,error_class,age_seconds,error_code … schema|NULL`; Schluss `OK — alle Belege real erbracht` |
| `make gates` (ungepiped, Exit in eigener Datei) | 0 | — |

**Übernommen, nicht selbst gefahren:** `make test-integration`, `make test-replication`, `make image`,
`make sdk-pack-*`, `make examples-*` (der Implementer meldet sie grün; die Zeile des Fehlerzustand-Belegs im
Plan §3 ist *übernommen*).

**Nachgemessene Zählung (Schwerpunkt a):** `git grep -n '\.Warn(' -- internal cmd` ohne Pathspec-Ausnahme: 39 Zeilen;
davon 4 in Test-Dateien (`slog_levels_internal_test.go`, `log_test.go` zweimal, ein Kommentar in
`wiring_rest_internal_test.go`); Produktion **35** (wie im Plan). Aufteilung gelesen: 31 mit `W`-Code, 1 mit
`E`-Code (`heartbeat: Fehlerzustand gemeldet`), 3 ohne Code (`grpc/administration.go:63`, `http/errors.go:36`,
`http/registerconsumer.go:51`); `git grep -n 'LogKey' -- internal ':!*_test.go'` zählt 34 Zeilen (32 Aufrufe plus
Definition und `Area`-Kommentar, konsistent mit der Suchlauf-Zeile `diff 32`).

**Mutationsproben (selbst, je einzeln, an Kopien):**

| Mutation | Ergebnis |
|---|---|
| `Fault` schreibt `""` statt des Codes | `make test-store` rot: `TestFaultWritesErrorClass`, `TestHeartbeatViewProjectsErrorClass` |
| `CASE` der View `cdc.heartbeat` entfernt | `make test-store` rot: `TestHeartbeatViewHidesCodeWithoutClass` |
| `classifyRunFault` liefert `codes[0]` | rot: `TestClassifyRunFaultCodeBelongsToTheClass` |
| HTTP-Diagnose ohne `ErrorCode` | rot: `TestDiagnoseReaderTokenLiestBericht` (erster Versuch mit falschem Muster nicht mutiert, wiederholt) |
| gRPC-Diagnose ohne `ErrorCode` | rot: `TestDiagnoseRuftUseCaseMitDerQuelleAufUndUebersetztDenBericht` |
| Warn-Code in `http/sse.go` vertauscht (`W1005` → `W1002`) | rot: `TestStreamNichtKodierbareChangeBeendetDenStream` |
| Vertauschung je Datei: `W2003`→`W2004`, `W2002`→`W2001` (`bootstrap/backfill.go`); `W2001`→`W1001` (`backfill/service.go`); `W1002`→`W1001` (`capture/service.go`); `W1004`→`W1003` (`publisher.go`); `W3003`→`W3001`, `W4007`→`W4005`, `W4003`→`W4004` (`wiring.go`) | je rot, je mit einem benannten Test (u. a. `TestBackfillWorkerDoesNotSpinOnARunThatStaysQueued`, `TestRunRetentionCleanupContinuesAfterError`, `TestRunStreamAfterAdministrationPassWithTimeoutLeavesRequestsPendingAfterTheDeadline`) |

Die Bindung der Zusagen an ihre Eingabeseite trägt: jede der 13 Mutationen färbt einen Test rot.

---

## Findings

### F-1 — Katalog-Bedeutung von `PCF-W1002` trifft die Stelle in der Erfassung nicht

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2 (Stellen mit gleicher Maßnahme teilen einen Code), Skill „Beleg trägt seinen Satz nicht"
- `pfad`: `internal/application/usecase/capture/service.go:151`, `docs/user/benutzerhandbuch.md` (Zeile `PCF-W1002`), `internal/adapters/driven/grpcstream/broadcaster.go:96`
- `befund`: Die Stelle `capture: Stream-Publish fehlgeschlagen` ruft `ChangeStreamPort.Publish` des prozessinternen `grpcstream.Broadcaster` auf (`wiring.go:699` verdrahtet `changeBroadcaster`); dessen `Publish` liefert nur bei `nil`-Change oder beendetem Kontext einen Fehler und berührt NATS nicht. Der Katalog nennt für `PCF-W1002` „Veröffentlichung … im NATS-Vollinhalts-Stream" mit der Maßnahme „NATS-Server und Verbindung prüfen"; für die Erfassungs-Stelle ist diese Maßnahme nicht die des Betreibers, die Gruppierung „gleiche Maßnahme = gleicher Code" (Plan §3, Auswahl der Warn-Stellen) trägt dort nicht.
- `verifizierbar`: ja — Lesen von `grpcstream/broadcaster.go:96` bis `:123` und `wiring.go:699`; kein Gate bildet es ab (`make meldungscodes-check` vergleicht Mengen, nicht Bedeutungen).
- `klasse`: Zuordnung gegen die Bedeutung nicht getragen

### F-2 — Warn-Zeile trägt einen `E`-Code im Attribut `code`

- `kategorie`: LOW
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1 und 3 (Log-Zeile mit Warnung → `code=<code>`, Warnungen tragen keine Klasse), [`SPEC-008`](../../spec/pflichtenheft.md) Absatz „Meldungscode"
- `pfad`: `internal/adapters/driven/postgresstorage/heartbeat.go:111`
- `befund`: `heartbeat: Fehlerzustand gemeldet` ist eine `Warn`-Zeile und trägt den `E`-Code des Fehlerzustands unter demselben Schlüssel `code` wie die 31 Warnungen mit `W`-Code; die Spec („bei einer Warnung ist die erste Ziffer der Bereich") und der Katalog („Log-Zeile einer Warnung: `code=PCF-W…`") nennen diese Form nicht, ein Betreiber, der auf `code=PCF-W` filtert, sieht die Zeile nicht. Die Lesart ist eine Auslegung von Festlegung 3, kein Verstoß gegen ihren Wortlaut.
- `verifizierbar`: nein — Auslegungsfrage (Architect-Frage 1).
- `klasse`: Auslegung einer Festlegung ohne Träger-Zeile

### F-3 — Begründung „Anhängen genügt" für das Fehlen des Vorlaufs benennt nicht den wirkenden Mechanismus

- `kategorie`: LOW
- `quelle`: [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 1 (Klasse „View-Signatur" ist `ReplaceView` einer im neutralen Modell deklarierten View), `tools/schema/rolloutguard/guard.go` `knownForeignObjects`
- `pfad`: `tools/schema/nacharbeit-heartbeat.sql:22`, Plan §3 Zeile „Nicht realisiert (Abweichung vom Plan)"
- `befund`: Der Kommentar und der Plan führen das Fehlen eines Vorlaufs darauf zurück, dass `error_code` als letzte Spalte angehängt wird und `CREATE OR REPLACE VIEW` deshalb genügt. `cdc.heartbeat` steht aber in `knownForeignObjects` als `DropView`, d-migrate meldet es als destruktiven Blocker und der Rollout löscht es unter `--allow-destructive` vor der Nacharbeit; die Klasse des Vorlaufs trifft die View nicht, weil sie nicht im neutralen Modell steht. Das Ergebnis (Exit 0, Spalten in Soll-Reihenfolge, Alt-Zeile `schema|NULL`) ist in Lauf 5 belegt und stimmt; die genannte Reihenfolge-Begründung trägt es nicht, und Lauf 5 trennt beide Erklärungen nicht (die View wird dort ohnehin neu angelegt).
- `verifizierbar`: ja — `guard_test.go:315` (Report mit `DropView … heartbeat`, `rendered: true`) und Lauf 5.
- `klasse`: Begründung nennt nicht den wirkenden Mechanismus

### F-4 — Gemischte Versionen: ein älterer Server kann einen Code neben fremder Klasse hinterlassen

- `kategorie`: LOW
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2 (Klasse folgt aus dem Code), Plan §6 (Ausgang offen)
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go` (`UpsertHeartbeatFault`), `tools/schema/nacharbeit-heartbeat.sql:35`
- `befund`: Der Upsert eines Servers ohne `error_code` (`error_class = EXCLUDED.error_class`) lässt einen früher von einer neuen Version geschriebenen `error_code` stehen; die View blendet ihn nur bei `error_class IS NULL` aus. Nach Rückkehr auf die ältere Version und einem Fehlerzustand anderer Klasse zeigt `cdc.heartbeat` Klasse und Code verschiedener Klassen. Der Handbuch-Abschnitt „Reihenfolge beim Upgrade" behandelt nur die Richtung neuer Server gegen altes Schema; der Rückweg ist nicht genannt.
- `verifizierbar`: nein — kein Lauf gegen einen alten Server (Folge von Schreibzügen abgeleitet, nicht gefahren).
- `klasse`: Rückkompatibilität nur in einer Richtung benannt

### F-5 — Plan widerspricht seinem Nachbarn: Risiko und Liefer-Punkt sprechen von View-Signatur-Änderung und Vorlauf

- `kategorie`: LOW
- `quelle`: Skill „Nachzug widerspricht dem Nachbarn im selben Träger", [`AGENTS.md`](../../AGENTS.md) §3.13
- `pfad`: `docs/plan/planning/in-progress/slice-meldungscodes-warnungen-heartbeat-diagnose.md` §1 Ziel, §2 DoD (B), §6 erstes Risiko
- `befund`: §6 nennt weiter „Die Spalte ändert die Signatur von `cdc.heartbeat`; ohne Vorlauf scheitert der Rollout", §2 (B) verlangt „Vorlauf bei View-Signatur-Änderung", während §3 in der Zeile „Nicht realisiert" Vorlauf und `rolloutguard`-Änderung als entfallen ausweist; keine der beiden Stellen verweist auf die andere. Der Ausgang des Risikos steht „(bei Closure)" noch aus.
- `verifizierbar`: ja — Lesen von §2, §3, §6 des Plans.
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-6 — SPEC-031: `error_code` steht in der Aufzählung hinter dem Gedankenstrich von `error_class`

- `kategorie`: LOW
- `quelle`: [`SPEC-031`](../../spec/pflichtenheft.md) Hilfsnachrichten von `Diagnose`, Maintainability
- `pfad`: `spec/pflichtenheft.md` (Absatz „Hilfsnachrichten von `Diagnose`", `HeartbeatStatus`)
- `befund`: Der Satz lautet „`error_class` — leer bedeutet Normalbetrieb, nur gültig wenn `known`, `error_code` (Feldnummer 4) — der Meldungscode …"; `error_code` ist grammatisch Teil der Erläuterung von `error_class`, nicht ein eigener Aufzählungspunkt wie `known`, `age_seconds`, `error_class`.
- `verifizierbar`: nein — Lese-Handlung.
- `klasse`: Aufzählung nicht lesbar abgesetzt

### F-7 — `PCF-W1003` und `PCF-W1004` im Katalog

- `kategorie`: INFO
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2
- `pfad`: `internal/adapters/driven/natsstream/publisher.go:260` bis `:279`, Handbuch-Katalog
- `befund`: `PCF-W1003` fasst „Schema oder Tabelle leer" und „reserviertes Zeichen" zusammen (gemeinsame Maßnahme „Namen prüfen", im Rahmen der Festlegung). Für `PCF-W1004` verweist die Maßnahme auf die Korrektur des Zielnamens einer Regel, die Stelle benennt jedoch selbst, dass das Alphabet des Zielnamens das reservierte Zeichen ausschließt; die Warnung ist über den Regelweg nicht erreichbar.
- `verifizierbar`: nein.
- `klasse`: Katalogzeile ohne erreichbaren Auslöser

### F-8 — Maßnahme von `PCF-W2001` nennt eine Seite

- `kategorie`: INFO
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 5 (Katalog: Maßnahme)
- `pfad`: `internal/application/usecase/backfill/service.go:571` und `:582`, Handbuch-Katalog
- `befund`: Die Zeile fasst das Schließen des Snapshots (Quelle) und den Rollback der Schreibtransaktion (CDC-Speicher) zusammen; die Maßnahme nennt nur „Erreichbarkeit der Quelle prüfen".
- `verifizierbar`: nein.
- `klasse`: Maßnahme deckt eine von zwei Stellen

### F-9 — Erinnerung für den Architect: Bereich 5

- `kategorie`: INFO
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1
- `pfad`: `internal/domain/messagecode/codes.go`, `messagecode.go` (`Area`)
- `befund`: `Area` akzeptiert die Bereiche 1 bis 5, vergeben sind Bereiche 1 bis 4; der Handbuch-Katalog listet nur 1 bis 4, die Spec alle fünf. Kein Widerspruch, der Bereich 5 bleibt ohne Warnung im Bestand.
- `verifizierbar`: nein.
- `klasse`: Hinweis

---

## Architect-Fragen

1. **Fehlerzustand als Warn-Zeile (F-2).** Soll `heartbeat: Fehlerzustand gemeldet` den `E`-Code unter `code` tragen
   (Auslegung von [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3, so umgesetzt),
   einen `W`-Code des Bereichs 1 erhalten oder ohne Code bleiben? Spec und Katalog brauchen in jedem Fall eine Zeile dazu.
2. **Prozessinterner Broadcaster (F-1).** Ist ein Fehlschlag von `ChangeStreamPort.Publish` im Erfassungspfad eine Warnung
   mit Betreiber-Maßnahme (eigener Code) oder ohne (kein Code wie die drei Stellen für T4)? Die Entscheidung gehört in den
   Plan, nicht in diesen Report.
3. **Mechanismus des Vorlaufs (F-3).** Genügt dem Architect die Feststellung „`cdc.heartbeat` ist Fremdobjekt und wird
   beim Rollout ohnehin ersetzt", oder soll die Begründung in `harness/targets/schema-rollout.md` ausdrücklich stehen?

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/messagecode/` (Tabelle: 22 `W`-Codes in vier Bereichen je einmal, `Area` und `LogKey`; Registry-Test Ziffer gegen Bereich grün)
- geprüft, ohne Befund: `classifyRunFault` gegen `a9767e87:internal/bootstrap/wiring.go` — Klassenwahl unverändert (die Vorrangfolge ist um `internal` am Ende erweitert, das Ergebnis ohne Code der sechs Klassen bleibt `internal`), `classifyRunError` delegiert; `reportFault` ignoriert wie zuvor den Rückgabewert
- geprüft, ohne Befund: `postgresstorage/heartbeat.go` (Port-Grenze: Code einer anderen Klasse oder ohne Tabelleneintrag erreicht keinen SQL-Aufruf), `queries.go` (`Beat` löscht Klasse und Code)
- geprüft, ohne Befund: Alt-Zeilen (`error_code` NULL, in der View `NULL`, CLI-Zeile `Fehlerzustand: schema`), Alt-Server gegen neues Schema (Spalte nullable, Upsert ohne die Spalte gültig), Neu-Server gegen altes Schema (Schreibzug scheitert, im Handbuch als *abgeleitet* gekennzeichnet)
- geprüft, ohne Befund: Proto `administration.proto` (additive Feldnummer 4 in `HeartbeatStatus`, keine Reservierung berührt), `gen/` byte-gleich (`make generated-sync`), Go-Beispiel `examples/grpc-client/diagnose.go`, `tools/harness/grpcadminclient`
- geprüft, ohne Befund: Reichweite — `git diff --name-only a9767e87 HEAD -- sdks AGENTS.md .claude .harness spec/architecture.md spec/lastenheft.md` leer; kein Fehlerkörper-Code in HTTP/gRPC (T4), keine Paketversion, kein Release
- geprüft, ohne Befund: Handbuch 1.91 (Kopf hochgezählt, Änderungshistorie in Betreibersicht, Katalog 22 Warncodes gleich Tabelle, `make handbuch-public-doc-check` Exit 0); zehn Maßnahme-Zeilen gegen den Code gelesen (`W1001`, `W1003`, `W1005`, `W1006`, `W2003`, `W3002`, `W4001`, `W4003`, `W4004`, `W4006`)
- geprüft, ohne Befund: Spec — [`SPEC-008`](../../spec/pflichtenheft.md) (Warncodes, Vorrangfolge bei mehreren Codes, erster Code der gewählten Klasse, Rückfall `PCF-E7000`), Heartbeat-Tabelle, [`SPEC-018`](../../spec/pflichtenheft.md), Änderungstabelle; `spec/architecture.md` unberührt ([`AGENTS.md`](../../AGENTS.md) §3.4)
- geprüft, ohne Befund: Kommentare im Diff ([`AGENTS.md`](../../AGENTS.md) §3.7): `make kommentar-kennungen DIFF=a9767e87` ohne Kandidat; Lesen der neuen und geänderten Kommentare ohne Vorher/Nachher-Sprache (Ausnahme F-3 für die Begründung)
- geprüft, ohne Befund: [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) unberührt (Accepted, [`AGENTS.md`](../../AGENTS.md) §3.5); keine Suppression ([`AGENTS.md`](../../AGENTS.md) §3.2), keine Gate-Lockerung ([`AGENTS.md`](../../AGENTS.md) §3.6), kein Host-Werkzeug am Repo ([`AGENTS.md`](../../AGENTS.md) §3.1)
- geprüft, ohne Befund: Commit-Messages tragen [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), keine `SPEC-*`/`ARC-*` im Betreff (`make commit-traceability` Exit 0)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zuordnung gegen die Bedeutung nicht getragen · Auslegung einer Festlegung ohne Träger-Zeile · Begründung nennt nicht den wirkenden Mechanismus · Rückkompatibilität nur in einer Richtung benannt · Nachzug widerspricht dem Nachbarn im selben Träger · Aufzählung nicht lesbar abgesetzt

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM): der Katalog beschreibt für `PCF-W1002` eine Maßnahme, die für eine der beiden Stellen nicht gilt; der Katalog ist der Kern dieses Slice. Keine HIGH-Findings. Die DoD-Zeile „Review durchgeführt" im Plan bleibt offen, bis F-1 behoben oder vom Architect anders entschieden ist (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift nur ohne Fixrunde).

**Übergabe:** F-1 bis F-6 an den Implementer (F-2 und F-3 nach der Antwort des Architect); die Architect-Fragen gehen an den Architect. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
