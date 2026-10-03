# Slice meldungscodes-registry-fehlerkopf: Code-Tabelle, Fehlerkopf `Fehlerklasse <klasse> [<code>]`, Handbuch-Katalog und Gate `meldungscodes-check` (Teil 2 von 4)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (Betriebsfähigkeit),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (Fehlerklassen),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegungen 1–5, 8, 9 T2),
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) und
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Klassen, unverändert),
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(Handbuch-Gate; Abgrenzung in Festlegung 7),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md) (neuer Absatz
„Meldungscode“; der Satz „der Fehlertext beginnt mit der Klasse“ wird zum Kopf nach
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3),
[`SPEC-029`](../../../../spec/pflichtenheft.md) (`error_message` des Backfill-Runs),
[`SPEC-019`](../../../../spec/pflichtenheft.md) (`error_message` des Antrags) — Liefer-Punkt
dieses Slice, nicht des Auftraggebers.

**Reihenfolge:** nach [T1 `meldungscodes-kennungsfreie-ausgaben`](../done/slice-meldungscodes-kennungsfreie-ausgaben.md)
(sein Gate hält die Ausgaben bei 0, bevor dieser Slice die Fehlertext-Literale anfasst); vor
[T3](../open/slice-meldungscodes-warnungen-heartbeat-diagnose.md) und
[T4](../open/slice-meldungscodes-http-grpc-fehlerkoerper.md), die an der Tabelle dieses Slice hängen.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Die Quelle der Fehlertexte trägt ein einheitliches Präfix `Fehlerklasse <klasse>: …`:
**28** Literale in Produktions-Go. `classifyRunError` (`internal/bootstrap/wiring.go`)
bildet Sentinels über `errors.Is` auf die sieben Klassen ab; der Backfill-Dienst entfernt den
Kopf per Textersetzung (`internal/application/usecase/backfill/service.go`,
`strings.Replace(cause.Error(), "Fehlerklasse "+string(class)+": ", "", 1)`) — diese
Textchirurgie ersetzt der Slice durch den Code des Fehlerwerts. Das Handbuch nennt das Wort
`Fehlerklasse` in 36 Zeilen. Kein `PCF-` im Baum (0 Treffer).

**Ziel:** Es gibt eine Code-Tabelle als Quelle der Wahrheit (Paket `internal/domain`),
jeder klassifizierte Fehler trägt den Code seiner Einzelursache oder den Rückfall seiner
Klasse, die Wege dieses Slice zeigen den Kopf `Fehlerklasse <klasse> [<code>]: …`, das
Handbuch trägt den Katalog, und das Gate `meldungscodes-check` hält Quelltext, Tabelle und
Katalog gleich ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegungen 1–5).

**Ausdrücklich NICHT in diesem Slice:**

- **Warnungs-Codes (`W`), Heartbeat-Spalte `error_code`, Diagnose** — T3.
- **`code` im HTTP-Fehlerkörper, gRPC `ErrorInfo`** — T4. Die Ablehnungs-Codes `E8…`
  entstehen in der Tabelle hier (der Antrag trägt sie in `error_message`), ihre Netz-Wege in T4.
- **Änderung der Fehlerklassen-Semantik, Exit-Codes** — Klassen und Ausgang 1 bleiben
  ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2, 3).
- **Code als Metrik-Label** — ausgeschlossen (Festlegung 2).
- **SDK-Änderungen, Release, Tag** — Freigaben des Auftraggebers; der Slice ändert `sdks/`
  nicht (Server-Release-Frage: Empfehlung 3 des Verdikts, Auftraggeber-Entscheidung).
- **Lastenheft** — keine Änderung nötig; eine Anforderung wäre Auftraggeber-Entscheidung.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [x] **(A) Tabelle, Kopf und Wege.** Die Tabelle in `internal/domain` (Code, Schwere, Klasse,
      Status) enthält die sieben Rückfall-Codes `PCF-E1000`…`PCF-E7000` und die Einzelursachen
      der Sentinels (Granularität: unterschiedliche Maßnahme des Betreibers,
      [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2) sowie
      die Ablehnungs-Codes `E8…` des Antrags; die 28 `Fehlerklasse`-Literale verwenden die
      Konstanten der Tabelle (kein zweites Literal); der Kopf `Fehlerklasse <klasse> [<code>]: …`
      steht in Fehlerwert, Prozess-Ende-Zeile und `error`-Attribut; `error_message` von Run
      (`<klasse> [<code>]: …`) und Antrag (`abgelehnt [<code>]: …`); die Textchirurgie im
      Backfill-Dienst entfällt; Rollout-Fehlerzeilen als `FEHLER [<code>]: …`; Spec-Nachzug
      (`SPEC-008` Absatz, Satz zum Kopf, `SPEC-029`, `SPEC-019`). *Zu belegen durch:* Go-Test in
      `make test` (Format gegen ERE `PCF-[EWI][0-9]{4}`, erste Ziffer gegen Klasse, keine Dopplung,
      Rückfall je Klasse; ein Code mit falscher Klasse färbt rot — in der ADR nur *hergeleitet*,
      der Implementer fährt die Probe und berichtet Stelle und Farbe); Suchlauf §3;
      `make test-integration` und `make test-store` grün nach `make image`.
- [x] **(B) Katalog und Gate.** Handbuch-Katalog (kennungsfrei: Code, Klasse bzw. Bereich,
      Bedeutung, Maßnahme; zurückgezogene Codes mit Vermerk; der Satz „Der Text ist nicht
      Vertrag, der Code ist es“; Änderungshistorie in Betreibersicht); Gate `meldungscodes-check`
      in `GATE_CHECKS` (Mengengleichheit Quelltext/Tabelle/Katalog, nur `bash`/`git`/`grep`),
      Heimat `harness/mk/doc-gate.mk`, Sensor-Vertrag `harness/sensors/meldungscodes-check.md`,
      Tabellentest `make test-meldungscodes-check` (Werkzeug ohne Gate-Aufnahme), Zeilen in
      `harness/README.md` §Sensors. *Zu belegen durch:* Tabellentest (Code im Quelltext ohne
      Tabelle → Exit 1, Tabellen-Code ohne Katalog → Exit 1, zurückgezogener Code bleibt in
      beiden Mengen, saubere Menge → Exit 0) und `make handbuch-public-doc-check` grün **mit**
      Katalog (Gate-Lauf; in der ADR nur aus den Mustern *hergeleitet*, Festlegung 7);
      `make gates` Exit 0.
- [x] **(C) Textzusage und Bestand geprüft.** Der Implementer liest die Handbuch-Abschnitte
      Diagnose, Fehlerklassen und Exit-Codes auf eine Zusage des Fehlertexts
      ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Konsequenzen: im Handbuch
      nicht nachgemessen); findet er eine, geht der Slice mit dieser Frage zurück an den
      Auftraggeber. Er misst außerdem, ob `sdks/` Servertexte der Form `Fehlerklasse …` in Tests
      erwartet (Messung am Parent: 0 Zeilen mit `Fehlerklasse` unter `sdks/`). *Zu belegen
      durch:* Befehle und Ergebnis im Bericht (§3.12: gemessen, nicht übernommen).
- [x] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review; der Reviewer liest den Katalog
      gegen den Code (akzeptiertes Negativ der ADR: das Gate vergleicht Mengen, nicht Sinn).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-registry-fehlerkopf.md`
      endet mit Exit 0.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten, nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/` (neue Tabelle), Go-Test dazu | neu | Quelle der Wahrheit, Registry-Test (A). |
| 28 `Fehlerklasse`-Literale (u. a. `mapper.go`, `receive.go`, `decode.go`, `tableactivation.go`, `internal/application/port/outbound/*.go`, `internal/bootstrap/wiring.go`) | update | Kopf mit Code aus der Tabelle. |
| `internal/bootstrap/wiring.go` (`classifyRunError`), `internal/application/usecase/backfill/service.go` | update | Code des Fehlerwerts statt Textersetzung; `error_message` des Runs. |
| Antragsverarbeitung (`error_message` einer abgelehnten Antrags-Zeile) | update | `abgelehnt [<code>]: …`, Codes aus dem Bereich `E8`. |
| `tools/schema/rollout.sh`, `tools/schema/rolloutguard/` | update | `FEHLER [<code>]: …`. |
| `spec/pflichtenheft.md` (`SPEC-008`, `SPEC-029`, `SPEC-019`) | update | Spec-Nachzug (Liefer-Punkt dieses Slice). |
| `docs/user/benutzerhandbuch.md` | update | Katalog, Satz zur Stabilität, Änderungshistorie. Rebase auf `main`, nur eigene Abschnitte committen. |
| Gate-Skript unter `tools/harness/`, Tabellentest-Läufer, `harness/sensors/meldungscodes-check.md`, `harness/README.md` §Sensors, `harness/mk/doc-gate.mk`, `GATE_CHECKS` | neu / update | (B). |
| `*_test.go` mit Erwartung `Fehlerklasse <klasse>:` (u. a. `internal/application/usecase/backfill/service_test.go`), Läufer unter `tools/harness/` | update | Erwartungen an den Kopf. |
| `internal/domain/messagecode/{codes,messagecode,messagecode_test}.go` | neu (Umsetzung von Zeile 1) | `codes.go` ist die Tabelle (66 Codes: 7 Rückfälle `…000` der Fehlerklassen, 26 Einzelursachen-Codes der Klassen 1 bis 6, 10 Codes `E7001` bis `E7010` der Klasse `internal` — die zehn Sentinels, die `classifyRunError` vor diesem Slice nicht kannte und auf `internal` fielen, tragen sie, damit ihre Klasse unverändert `internal` bleibt, und der Text eines Runs trägt für dieselben Ursachen weiter die Codes der Klassen 1 bis 6 (`failureCode`) —, 23 Codes der Ablehnungen `E8…` mit dem Rückfall `E8000` und 22 Gründen); `messagecode.go` trägt Typ `Code`, `Class`, `Entry`, `Error` (Sentinel-Fehlerwert mit Kopf), `From`, `Codes` (alle Codes einer Kette, Grundlage der Vorrangfolge von `classifyRunError`), `WithoutHead`, `RunMessage`, `RejectionMessage`; der Test prüft Form, erste Ziffer gegen Klasse, Dopplung, Rückfall je Klasse. Granularität (Maßnahme des Betreibers): ein Code je Speicher-Port (sechs Codes für `storage`: Change, Backfill, Consumer-Stand, Heartbeat, Schema, Snapshot; Antrags-Queue und Diagnose tragen Codes der Klasse `internal`, `E7002`, `E7003`), weil der Betreiber an der Fehlerzeile erkennt, welcher Speicher gestört ist; die drei Stream-Ordnungs-Sentinels des Mappers teilen einen Code. |
| `internal/domain/model/backfillrun.go` | update | `Run.Fail(at, code, cause)` — die Klasse folgt aus dem Code, `error_message` = `<Klasse> [<Code>]: <Ursache>`; ein Code ohne Fehlerklasse endet als `ErrInvalidErrorClass`. |
| `internal/bootstrap/rejection.go`, `internal/bootstrap/messagecodes_internal_test.go` | neu | `administrationFailureText`/`rejectionCode`: Ablehnungs-Code je Grund (E8…), unerwartete Fehler tragen den Kopf `internal`; der Test bindet jeden Sentinel an Code und Klasse und jeden Ablehnungsgrund an seinen Code. |
| `internal/adapters/driven/postgresstorage/sqlexec/translate.go` (`rejectionMessage`) | update | `abgelehnt [<Code>]: …` für die vom Antrags-Konstruktor verworfenen Zeilen (E8001 bis E8006, Rückfall E8000). |
| `test/integration/{routing,transformation}_e2e_test.go`, `tools/harness/run-integration-tests.sh` | update | Erwartungen an `error_message` (Kopf `abgelehnt [<Code>]`, Run `schema [PCF-E4004]`, DDL-Fenster `storage [PCF-E5008]`/`transient [PCF-E1003]`) und an die Log-Zeile (Kopf mit Code im Container-Log, `docker logs`). |
| `harness/targets/schema-rollout.md` | update | die `FEHLER`-Zeile trägt den Code. |
| `tools/schema/rolloutguard/` | **nicht geändert** | seine Zeilen (`rolloutguard: …`) sind Diagnose der Wache, keine `FEHLER`-Zeilen des Rollout-Skripts; [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3 nennt dort nur `FEHLER [<code>]` des Skripts. Eine Reduktion gegenüber der Zeile oben, kein stilles Weglassen. |
| `harness/sensors/meldungscodes-check.md`, `tools/harness/run-meldungscodes-check-tests.sh` | neu | Sensor-Vertrag und Tabellentest (71 Prüfungen) zu Liefer-Punkt (B). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Fehlertext beginnt mit
`Fehlerklasse <klasse>: `“; Parent ist `ba60c7bc`; die `diff`-Zeilen und die Befunde trägt der
Implementer ein):**

```suchlauf
ba60c7bc 28 -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'
ba60c7bc 1 -n -F 'strings.Replace(cause.Error(), "Fehlerklasse "' -- internal/application/usecase/backfill/service.go
ba60c7bc 1 -n -F 'func classifyRunError' -- internal/bootstrap/wiring.go
ba60c7bc 1 -n -F 'der Fehlertext beginnt mit der Klasse' -- spec/pflichtenheft.md
ba60c7bc 36 -n -F 'Fehlerklasse' -- docs/user/benutzerhandbuch.md
ba60c7bc 2 -n -E 'Fehlerklasse [a-z]+:' -- internal/application/usecase/backfill/service_test.go
ba60c7bc 0 -n -F 'Fehlerklasse' -- sdks
ba60c7bc 0 -n -F 'PCF-' -- internal tools sdks docs/user spec cmd examples
diff 0 -n -E '"Fehlerklasse [a-z]+: ' -- '*.go' ':!*_test.go'
diff 0 -n -F 'strings.Replace(cause.Error(), "Fehlerklasse "' -- internal/application/usecase/backfill/service.go
diff 1 -n -F 'func classifyRunError' -- internal/bootstrap/wiring.go
diff 0 -n -F 'der Fehlertext beginnt mit der Klasse' -- spec/pflichtenheft.md
diff 43 -n -F 'Fehlerklasse' -- docs/user/benutzerhandbuch.md
diff 0 -n -E 'Fehlerklasse [a-z]+:' -- internal/application/usecase/backfill/service_test.go
diff 0 -n -F 'Fehlerklasse' -- sdks
diff 0 -n -F 'PCF-' -- sdks
diff 0 -n -F 'PCF-' -- examples
```

| Träger | Messung am Parent (`ba60c7bc`, 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| `Fehlerklasse`-Literale in Produktions-Go | Zeile 1: 28 (die Tabellenzeile des Vorläufer-Plans gehört jetzt hierher, nicht zu T1) | Diff Zeile 9: 0 — alle 28 Sentinels sind `messagecode.New(<Konstante>, "<Ursache>")`, kein zweites Literal |
| Textchirurgie, Klassifikation | Zeilen 2, 3: je 1 | Diff Zeile 10: 0 (die Textersetzung entfällt, `messagecode.WithoutHead` liest den Kopf des Fehlerwerts); Diff Zeile 11: 1 (`classifyRunError` bleibt, leitet die Klasse aus den Codes der Kette ab, mit der Vorrangfolge der Klassen des Parents) |
| Spec-Satz zum Kopf | Zeile 4: 1 | Diff Zeile 12: 0 — auf Klasse und Meldungscode umgeschrieben (`SPEC-008` Absatz „Meldungscode“, `SPEC-029`, `SPEC-019`) |
| Handbuch | Zeile 5: 36 | Diff Zeile 13: 43 = 36 + 6 Zeilen des Abschnitts „Meldungscodes“ + 1 Zeile der Änderungshistorie. Jede der 36 Zeilen gelesen: Textzusage des Fehlertexts (z. B. „beginnt mit der Klasse“) gefunden: **keine**; drei Zeilen zitieren einen Fehlertext und sind mitgezogen (Beispiele des Log-Texts `schema [PCF-E4005]`/`schema [PCF-E4003]`, Beschreibung von `error_message` in der Run-Statustabelle), zwei Tabellen mit Fehlertexten von Anträgen (Konflikte der Transformationsregeln, Prüfungen der Routing-Regeln) tragen den Kopf `abgelehnt [<Code>]: `, der Backfill-Hinweis `transient: …` trägt den Code |
| Tests | Zeile 6: 2 | Diff Zeile 14: 0 — Erwartungen an Kopf und Code mitgezogen in `internal/application/usecase/backfill/{service,routing,transformation}_test.go`, `internal/bootstrap/{administration,routing,backfill}_internal_test.go`, `internal/bootstrap/{backfill_endtoend,backfill_routing_store_internal}_test.go`, `internal/adapters/driven/postgresstorage/{sqlexec/translate,sqlexec/rejection_internal,backfill*}_test.go`, `internal/domain/model/backfillrun_test.go`, `test/integration/{routing,transformation}_e2e_test.go` und `tools/harness/run-integration-tests.sh` |
| SDKs | Zeile 7: 0 | Diff Zeile 15: 0 — Nichtgefunden belegt; auch kein SDK-Test oder Beispiel erwartet einen Servertext (`git grep` nach den Fehlertexten der Ablehnungen und Sentinels unter `sdks/` und `examples/`: 0 Treffer) |
| Code-Präfix im Baum | Zeile 8: 0 | Diff: Treffer in Tabelle, Katalog, Tests und Läufern (Zahl nicht fortgeschrieben, sie bewegt sich mit jedem Test); Diff Zeilen 16 und 17: 0 unter `sdks` und `examples` |

**Läufe des Implementers (Handoff-Stand, gemessen am Arbeitsbaum, Exit-Codes ungefiltert):**

| Befehl | Exit | gedruckte Zeile |
|---|---|---|
| `make test` | 0 | 50 Pakete `ok`, darunter `internal/domain/messagecode`, `internal/bootstrap`, `internal/application/usecase/backfill` |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 82.99% erfuellt Schwelle 80%` |
| `make test-replication` | 0 | `ok  github.com/pt9912/pg-change-feed/internal/adapters/driving/replication/receive 40.325s` (PostgreSQL 18.6) |
| `make image`, `make test-integration` | 0 | `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen`; die Phase DDL-Fenster druckt `failed (storage [PCF-E5008]: Lesefehler im Tabellen-Snapshot …` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | `run-schema-rollout-guard-test: OK — alle Belege real erbracht` |
| `make schema-validate SCHEMA_SOURCE=tools/schema/does-not-exist.yaml` | 2 | `FEHLER [PCF-E2007]: tools/schema/does-not-exist.yaml fehlt — …` |
| `docker run --rm --network none ghcr.io/pt9912/pg-change-feed:dev` | 2 | `pg-change-feed: Fehlerklasse configuration [PCF-E2001]: Verdrahtung ohne vollständige Vorbedingung: CDC_CAPTURE_DSN fehlt` (der Startfehler der Klasse `configuration` endet mit Ausgang 2, nicht 1) |
| `make meldungscodes-check` | 0 | `meldungscodes-check: 66 Codes in Tabelle und Katalog gleich, Quelltext (124 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` (Fixrunde nach dem Review) |
| `make test-meldungscodes-check` | 0 | `run-meldungscodes-check-tests: 71 Prüfungen bestanden` (Fixrunde nach dem Review) |
| `make handbuch-public-doc-check` (mit Katalog) | 0 | `handbuch-public-doc-check: keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make ausgabe-kennungen-check`, `make sdk-public-doc-check`, `make a-check`, `make fmt-check` | je 0 | `… 126 Go-Dateien und 4 Skripten`, `keine interne Kennung unter sdks`, `gesamt: 0 Befund(e)`, `328 Go-Dateien geprüft, alle formatiert` |

**Läufe der Fixrunde nach dem Review (Arbeitsbaum, Exit-Codes ungefiltert):**

| Befehl | Exit | gedruckte Zeile |
|---|---|---|
| `make test` | 0 | 50 Pakete `ok` |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 82.99% erfuellt Schwelle 80%` |
| `make image`, `make test-integration` | 0 | `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen`; das Container-Log trägt die Köpfe `Fehlerklasse schema [PCF-E4004]`, `[PCF-E4005]` und `[PCF-E4003]` (Prüfung des Läufers) |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | `run-schema-rollout-guard-test: OK — alle Belege real erbracht` |
| `make fmt-check`, `make a-check` | je 0 | `328 Go-Dateien geprüft, alle formatiert`, `gesamt: 0 Befund(e)` |

**Mutationen (Zusage · mutierte Eingabe · gesehenes Rot; alle an Kopien im Scratchpad, der Arbeitsbaum bleibt unberührt):**

| Zusage | mutierte Eingabe | Rot |
|---|---|---|
| Quelltext nur mit Codes der Tabelle | `tools/schema/rollout.sh`: `PCF-E2007` → `PCF-E2999` | `meldungscodes-check` Exit 1: `Code ohne Eintrag in der Tabelle internal/domain/messagecode/codes.go: tools/schema/rollout.sh:69:PCF-E2999` |
| Tabelle gleich Katalog | Handbuch: Katalog-Zeile `PCF-E5004` entfernt (Fixrunde; `PCF-E5005` gibt es nicht mehr) | Exit 1: `Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch.md: PCF-E5004` |
| erste Ziffer ist die Klasse | `codes.go`: Zeile `StorageFallback` mit Klasse `schema` | `TestTableClassMatchesFirstDigit` und `TestEveryClassHasItsFallback` rot |
| Rückfall je Klasse | `codes.go`: Tabellenzeile `PermissionFallback` entfernt | `TestEveryClassHasItsFallback`: `Rückfall "PCF-E3000" der Klasse "permission" fehlt in der Tabelle` |
| jeder Sentinel trägt Kopf und Code | `decode.ErrSchema` mit dem früheren Wortlaut ohne Code | `TestSentinelsCarryTheirCodeAndClass`, `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`, `TestReportFaultWritesClassifiedFaultOnNonNilError` rot |
| Run-Code folgt der Einzelursache | `failureCode`: Zweig `ErrSnapshotPermission` liefert `SnapshotConfiguration` | `TestExecuteClassifiesFailures` rot |
| Ablehnungs-Code je Grund (Antrag) | `rejectionCode`: `ErrColumnHasRule` liefert `RejectedRuleNameTaken` | `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts`, `…MapValueTakesEffectLive`, `TestAdministrationFailureTextHasACodeForEveryFailure` rot |
| Ablehnungs-Code je Grund (verworfene Zeile) | `rejectionMessage`: leere Quelle liefert `RejectedSchemaEmpty` | `TestReadPendingRequestsPassesRejectedRowsThrough` rot |
| keine Textchirurgie, Kopf einmal | `WithoutHead` entfernt nichts | `TestExecuteFailureTextCarriesClassOnce`, `TestNewBindsCodeToClass` rot |
| `error_message` des Runs trägt den Code | `RunMessage` ohne Code | `TestMessageForms`, `TestBackfillRunFields`, `TestExecuteRoutingStateChangeEndsRunAsConfiguration` rot |
| Klasse folgt aus dem Code | `classifyRunError` liefert immer `internal` | `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`, `TestSentinelsCarryTheirCodeAndClass`, `TestRunStreamWithRetryGesamtfensterFuenfMinuten` rot |
| Wächter-Logik: Tabelle gegen Katalog | Prüfling `comm -23 table.set table.set` (Kopie des Skripts) | `make test-meldungscodes-check` rot (u. a. `Tabellen-Code ohne Katalog-Zeile`, `Codes nur in Prosa`, `zurückgezogener Code fehlt im Katalog`) |
| Wächter-Logik: Form | Prüfling mit `[0-9]{4,5}` statt `[0-9]{4}` | Tabellentest rot: Fall „fünf Ziffern“ (Meldung `a.go:3:PCF-E40010 ` fehlt) |
| Wächter-Logik: Quelltext in Tabelle | Prüfling prüft nur das Handbuch, nicht den Quelltext | Tabellentest rot: Go-, cmd- und Skript-Fälle ohne Tabelle |
| Fixrunde: Klasse der zehn Sentinels bleibt `internal` | `codes.go` (Kopie): `SnapshotReadFault` auf `PCF-E5009` mit Klasse `storage` | `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes/Snapshot_lesen` und `TestInternalSentinelCodesStayInternal` rot |
| Fixrunde: Vorrangfolge der Klassen in einer Kette | `classifyRunError` (Kopie) liest nur den ersten Code der Kette | fünf Ketten-Fälle von `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` rot (`storage` vor `replication`, `replication` vor `configuration`, `internal` vor `storage`, `schema` vor `transient`, `errors.Join`) |
| Fixrunde: Binärzeichen in einer gelesenen Datei | `meldungscodes-check.sh` (Kopie): `grep -aon` → `grep -on` | Tabellentest rot: Fall „Go-Datei mit NUL-Byte und Code ohne Tabelle“ (Exit 0 statt 1) und seine Meldungsbindung |
| Fixrunde: Mengenvergleich ohne Auswertung | Kopien: `comm`-Zweig ohne `die2`; `sort -u` ohne `die2`; `sed`-Zweig ohne `die2` | je Tabellentest rot (Stub-`comm`/`sort`/`sed` mit Exit 2: Meldung `Lesefehler: Mengenvergleich …` bzw. `Mengenbildung der Tabelle …` fehlt) |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, [T1](../done/slice-meldungscodes-kennungsfreie-ausgaben.md) liegt in
`done/`, **und das Code-Präfix `PCF-` ist vom Auftraggeber bestätigt oder geändert**
([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Empfehlung 1 des
Verdikts: bisher nicht bestätigt). Wird es geändert, entsteht vor dem Start eine Folge-ADR
mit `Supersedes ADR-0144`, die Form und Muster nachzieht (Festlegung 1, 5, 7); nach diesem
Slice ist das Präfix nur noch per Folge-ADR änderbar. Der Hauptlauf klärt die Frage vor dem Start.

**Rückführung:** `in-progress` → `open` (blockiert), wenn der Implementer eine Zusage des
Fehlertexts im Handbuch findet (Liefer-Punkt C) oder ein SDK-Test den Servertext erwartet
(Auftraggeber-Entscheidung; die Änderung geht nie als Anpassung einer Erwartung durch).
`in-progress` → `next` (zu groß), wenn Tabelle und Wege zusammen nicht in einen Diff passen:
Schnitt entlang `application`/`adapters` (Einzelursachen-Codes) gegen Katalog/Gate, die ADR
bleibt unverändert.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration`, `make test-store` grün nach `make image`, Registry-
und Gate-Probe in den Bericht, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag.
Der sichtbar geänderte Wortlaut des Fehlerkopfs ist eine Betreiber-Änderung: ob sie ein
Server-Release mit eigener Freigabe braucht, entscheidet der Auftraggeber (Verdikt, Empfehlung 3).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **Präfix unbestätigt.** `PCF-` ist Empfehlung der ADR, nicht Auftraggeber-Entscheidung;
  nach diesem Slice ist es nur per Folge-ADR änderbar. Gegenmittel: §4-Startbedingung. —
  **Ausgang:** entfallen — `PCF-` ist vom Auftraggeber am 2026-10-03 bestätigt (*Ursprung:*
  Hauptlauf, Auftraggeber-Antwort in der Sitzung auf die Frage „Welches Präfix sollen die
  nutzerseitigen Meldungscodes tragen?“, Option „PCF-E4003 (Recommended)“; **übernommen**, kein
  committetes Artefakt trägt die Antwort außer dieser Zeile).
- **Ausgaben-Stabilität für Betreiber.** Der Kopf ändert sich zu
  `Fehlerklasse <klasse> [<code>]: …`; ein `grep` auf `Fehlerklasse schema:` bricht. Ausgang der
  ADR: keine Übergangsregel, Text nicht Vertrag; das Handbuch wird hier gemessen (Liefer-Punkt C). —
  **Ausgang:** eingetreten und aufgelöst — der Kopf ändert sich wie angekündigt; das Handbuch trägt
  den Satz „Der Text ist nicht Vertrag, der Code ist es“ und eine Änderungshistorie in
  Betreibersicht (Verifikation §2 (B), §9). Ob die sichtbare Änderung ein Server-Release braucht:
  siehe unten, Auftraggeber.
- **SDKs und Clients reichen Servertexte durch.** Ein SDK zeigt künftig den Code im Text mit;
  SDK-Tests, die den Wortlaut prüfen, wären betroffen. Gemessen am Parent: 0 Zeilen
  `Fehlerklasse` unter `sdks/`; Test-Erwartungen anderer Fehlertexte misst der Implementer. —
  **Ausgang:** entfallen — 0 Treffer für `Fehlerklasse` unter `sdks` und für `PCF-` unter `sdks`
  und `examples`, `sdks/` und `examples/` ohne Diff (Verifier **gemessen**, Verifikation §2 (C)).
- **Granularität der Einzelursachen.** Zu feine Codes erzeugen Pflege ohne Nutzen, zu grobe
  verwischen die Maßnahme; Regel der ADR: Maßnahme des Betreibers. Gegenmittel: Rückfall
  `…000` garantiert, dass kein klassifizierter Fehler ohne Code bleibt; der Reviewer liest die
  Zuordnung. — **Ausgang:** eingetreten, bewusst — für dieselbe Ursache existieren zwei Codes
  (Run-Text `E5008`, Log und Prozessende `E7010`; analog sechs weitere Paare, Verifikation V-3),
  das Handbuch erklärt es zeilenweise; Kehrseite von „Klassen unverändert“. Die Auflösung ist die
  Auftraggeber-Frage zur Klasse der zehn Sentinels (unten).
- **Handbuch-Gate mit Katalog.** Dass `make handbuch-public-doc-check` mit Codes grün bleibt,
  ist aus den Mustern hergeleitet, nicht gelaufen. — **Ausgang:** entfallen — gelaufen, Exit 0 mit
  Katalog (`handbuch-public-doc-check: keine interne Kennung in 3 Nutzerdokumenten unter
  docs/user/`, Implementer und Verifier **gemessen**).
- **Offene Frage an den Auftraggeber: Klasse der zehn Sentinels, die `classifyRunError` vor diesem
  Slice nicht kannte.** `outbound.ErrSchemaStoreStorage`, `ErrNotify`, `ErrAdministrationStorage`,
  `ErrBackfillStorage`, `ErrDiagnosticsStorage` und die fünf `ErrSnapshotPermission`,
  `ErrSnapshotConfiguration`, `ErrSnapshotTransient`, `ErrSnapshotReplication`, `ErrSnapshotStorage`
  fielen am Parent (`cb8afd07:internal/bootstrap/wiring.go`) auf `internal`; sie bleiben `internal`
  (Codes `E7001` bis `E7010`), weil [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
  Festlegung 2 „Klassen unverändert“ sagt. Eine Neuzuordnung (etwa `ErrSchemaStoreStorage` zu
  `storage`, `ErrNotify` zu `transient`) ist eine eigene Entscheidung des Auftraggebers und ändert
  `error_class` im Heartbeat und das Metrik-Label dieser Fälle. Sichtbare Folge der heutigen
  Zuordnung: der Kopf eines `failed` Antrags an der Antrags-Queue (`ErrAdministrationStorage`) lautet
  `Fehlerklasse internal [PCF-E7002]`, am Parent stand im Text `Fehlerklasse storage:`. Die Text-Folge
  gilt für alle zehn Sentinels `E7001` bis `E7010` (etwa `ErrNotify`: im Fehlertext vorher
  `transient`, jetzt `internal`), also in der Log- und der Prozessende-Zeile; `error_class` und
  Metrik-Label bleiben unverändert (*Ursprung:* Verifikation, gelesen am Diff gegen den Parent, nicht
  gefahren). —
  **Ausgang:** weiter offen → Auftraggeber-Entscheidung, Anker dieser Abschnitt (Neuzuordnung
  `E7001` bis `E7010` ändert `error_class` im Heartbeat und das Metrik-Label dieser Fälle; sie
  wäre eine Betreiber-sichtbare Klassenänderung und braucht eine eigene Freigabe). Bis dahin
  bleiben die zehn Sentinels `internal` (Verifier **gemessen**: 28 von 28 Sentinels, Klasse
  vorher und nachher unverändert, Verifikation §5).
- **Reserve ohne Verwender.** Die Rückfälle `E1000`, `E3000`, `E4000`, `E5000`, `E6000` und
  `messagecode.Fallback` haben im Produktionscode keinen Aufrufer (nur `E2000` und `E7000` werden
  emittiert), und `StatusWithdrawn` kennt weder Verhalten noch Test, solange kein Code
  zurückgezogen ist: bewusste Reserve nach [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
  Festlegung 1 (jede Klasse führt ihren Rückfall) und Festlegung 4 (zurückgezogen, nie neu
  belegt); der Test für den ersten zurückgezogenen Code entsteht mit diesem Code. —
  **Ausgang:** weiter offen, bewusst — der erste Verwender in T3 oder T4 und der erste
  zurückgezogene Code lösen es; Anker: Review F-4 und F-6.
- **Kollision mit parallelen Arbeiten am Handbuch.** Rebase auf `main`, nur eigene Abschnitte
  committen. — **Ausgang:** entfallen — kein Konflikt im Diff (Verifikation §9).
- **Server-Release für die sichtbar geänderte Fehlerkopf-Form (§5).** — **Ausgang:** weiter
  offen → Auftraggeber-Entscheidung; dieser Slice setzt weder Tag noch Version.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** (1) Die Tabelle als einzige Quelle (66 Codes) mit dem Gate
  `meldungscodes-check` als Mengenvergleich Quelltext / Tabelle / Katalog: Tabelle und Katalog
  sind 66 zu 66 gleich (das Gate **gemessen**, `make meldungscodes-check` und `make gates`
  Exit 0). Die 66 sind 7 Rückfälle + 26 Einzelursachen + 10 Codes der Klasse `internal`
  (`E7001` bis `E7010`) + 23 Ablehnungen (`E8…`). (2) Der Vergleich gegen den Parent
  (`git show cb8afd07:…`) in Review und Verifikation hat die Klassifikationsänderung gefunden
  und danach die Gleichheit belegt: 28 von 28 Sentinels, `error_class` und Metrik-Label vorher
  und nachher unverändert (Verifier **gemessen**, Verifikation §5); vier Go-Einzelmutationen
  der Fixrunde (Verifier) färbten rot. (3) Reale Läufe (Verifier **gemessen**): der Start ohne
  Umgebung druckt `pg-change-feed: Fehlerklasse configuration [PCF-E2001]: …` mit Ausgang 2,
  `make schema-validate` mit fehlender Quelle `FEHLER [PCF-E2007]: …`; `make test-integration`
  und `make test-store` grün. (4) Tabellentest des Gates: 73 Prüfungen (Implementer **gemessen**
  nach der Kleinstkorrektur `f3cc6964`; 71 waren es bei der Verifikation vor ihr, Verifier
  **gemessen**). `make gates` Exit 0 (Verifier **gemessen**, ungepiped).
- **Was ging anders als geplant:** (1) **Review und Fixrunde:** der Review fand F-1 (MEDIUM): der
  Umbau auf code-getriebene Klassen ordnete zehn Sentinels, die `classifyRunError` am Parent auf
  `internal` fallen ließ, einer anderen Klasse zu. Behoben durch die Fixrunde `c080136a` (die
  zehn tragen Codes der Klasse `internal`, `E7001` bis `E7010`; `classifyRunError` liest alle
  Codes der Kette und wendet die Vorrangfolge des Parents an). **Ein Re-Review nach der
  Fixrunde fand nicht statt** — ehrliche Abweichung vom Ablauf: der Verifier hat den
  Fixrunden-Diff (16 Dateien) selbst gelesen, vier Go-Einzelmutationen gefahren (rot) und einen
  weiteren Reviewer-Durchgang nicht für nötig befunden (Verifikation §4); die Review-Reports
  bleiben unverändert und tragen den Stand vor der Fixrunde. F-2 (LOW, Binärzeichen und
  Prozess-Substitution im Gate) in `fd155e1c` behoben, F-3 (LOW, Ursprungsangaben im Handbuch)
  in `3399cdd5`; F-4 bis F-6 (INFO) sind benannt (Reserve ohne Verwender, „Ausgang 1“ der ADR,
  `withdrawn` ohne Verhalten). Die Kleinstkorrektur `f3cc6964` nach der Verifikation: V-1 (Gate
  meldet den `sort`-Fehler der Dateiliste), V-5 (zwei Ursprungsangaben der Routing-Tabelle im
  Handbuch, `E8032` und `E8010`) und V-2 (Plan §6 nennt alle zehn Sentinels). Die Angaben zu
  `E8032` und `E8010` im Handbuch stehen als **übernommen** aus der Verifikation, nicht neu
  gefahren. (2) **Abweichung der ADR vom Code:** [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
  sagt „Ausgang 1“; gemessen endet ein Startfehler der Klasse `configuration` mit Ausgang 2
  (Bestand, `cmd/` ohne Diff, Verifikation §1); Handbuch und Spec sagen „Ausgang stabil“ und sind
  wahr. Eine ADR-Änderung entfällt ([`AGENTS.md`](../../../../AGENTS.md) §3.5). (3) **Offene
  Grenze (V-4):** der Spec-Satz „Die Klasse folgt aus dem Code“ ist enger als das Verhalten bei
  mehreren Codes in einer Kette (Vorrangfolge); ein halber Satz bei der nächsten Berührung von
  `SPEC-008` (T3). (4) **Pfad-Zitate in `ADR-0144` (V-6):** gemessen mit `git grep -n
  "slice-meldungscodes" -- docs/plan/adr`: zwei Zeilen, beide auf
  `done/slice-meldungscodes-kennungsfreie-ausgaben.md` (T1, bereits `done/`); zu T2 steht nur der
  Name im Fließtext, kein Pfad — **keine Zitat-Korrektur** nach
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) nötig. (5) **CI:**
  das Gate `meldungscodes-check` ist neu und lief auf dem Runner noch nicht; kein Workflow
  berührt, [`AGENTS.md`](../../../../AGENTS.md) §3.10 löst keine Pflicht aus; der erste Lauf nach
  dem Push ist zu beobachten (Erwartung, nicht gemessen).
- **Steering-Loop-Eintrag:** *Geschärfte Regel:* keine neue Regel im Wortlaut — Kandidat, 1×,
  unter der Schwelle: **ein Umbau mit der Zusage „Verhalten unverändert“ braucht einen Vergleich
  gegen den Parent über alle betroffenen Pfade, und der Plan nennt den Vergleichsbefehl.** Die
  Zusage „Klassen unverändert“ ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
  Festlegung 2) war verletzt, ohne dass ein Test, ein Plan-Satz oder die ADR die zehn Fälle
  nannte; der Fänger war der Parent-Vergleich des Reviewers, nicht ein Sensor. *Neuer Sensor:*
  keiner — die Fixrunde band die zehn Sentinels und die Vorrangfolge an einen Test
  (`TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`, `TestInternalSentinelCodesStayInternal`);
  ein allgemeiner Parent-Vergleich ist ohne Wissen, welche Pfade betroffen sind, nicht
  mechanisierbar. *Benannte Spec-Lücke:* V-4 (Vorrangfolge bei mehreren Codes steht in keiner
  Spec-Aussage) und die Klassenfrage der zehn Sentinels (§6).
- **Beobachtungs-Register (`../observations/`):** NEU
  `BEO-PGC/umbau-aendert-still-beobachtbares-verhalten` (1×, offen — unter der Schwelle; Beleg
  `evidence/slice-meldungscodes-registry-fehlerkopf.md`, Review F-1 und Fix `c080136a`).
  `BEO-PGC/intern-kennungen-in-ausgelieferten-texten`: Zustandsfeld fortgeschrieben — das Fangnetz
  für Programm-Ausgaben ist zweistufig (`ausgabe-kennungen-check` prüft Kennungsfreiheit,
  `meldungscodes-check` gleicht Codes ab); Zähler **1×** unverändert (kein neuer Vorgang der
  Klasse). Kein Beleg in `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`: F-1 ist eine
  Verhaltensabweichung eines Umbaus, keine Bindungslücke eines Negativtests; F-2 (LOW, Gate
  ohne Bindung an Binärzeichen) trifft einen bekannten Träger-Typ (Wächter-Skript mit
  Tabellentest), vor dem Merge gefunden, Schwere ≤ LOW — nach der Deckel-Regel ohne eigene Datei,
  hier genannt.
- **Folge-Slices:** T3 [`meldungscodes-warnungen-heartbeat-diagnose`](../open/slice-meldungscodes-warnungen-heartbeat-diagnose.md)
  und T4 [`meldungscodes-http-grpc-fehlerkoerper`](../open/slice-meldungscodes-http-grpc-fehlerkoerper.md)
  (beide in `open/`, untereinander unabhängig; Startbedingung: dieser Slice in `done/`). Kein
  Release, keine Versionsänderung.
- **Risiken aus §6:** Präfix: entfallen (vom Auftraggeber am 2026-10-03 bestätigt, **übernommen**);
  Ausgaben-Stabilität: eingetreten und aufgelöst; SDKs/Clients: entfallen (0 Treffer);
  Granularität: eingetreten, bewusst (zwei Codes für dieselbe Ursache, Handbuch erklärt es);
  Handbuch-Gate mit Katalog: entfallen (grün); **Klasse der zehn Sentinels: weiter offen →
  Auftraggeber-Frage** (Neuzuordnung ändert `error_class` und Metrik-Label dieser Fälle); Reserve
  ohne Verwender: weiter offen, bewusst; Kollision mit parallelen Handbuch-Arbeiten: entfallen;
  **Server-Release für die geänderte Fehlerkopf-Form: weiter offen → Auftraggeber**, kein Tag in
  diesem Slice.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); eine feinere Aufteilung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start auf Treffer zu
„Fehlertext“ und „Kennung in Ausgabe“ (nicht gemessen, kein Treffer behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
