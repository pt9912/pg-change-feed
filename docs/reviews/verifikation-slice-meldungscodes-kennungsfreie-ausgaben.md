# Verifikations-Report: slice-meldungscodes-kennungsfreie-ausgaben — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 9 T1 und 10,
Architect-Verdikt
[`architect-verdict-meldungscodes-statt-interner-kennungen`](architect-verdict-meldungscodes-statt-interner-kennungen.md))
+ Plan-vs-Code-Diff. Plan:
[`slice-meldungscodes-kennungsfreie-ausgaben`](../plan/planning/done/slice-meldungscodes-kennungsfreie-ausgaben.md)
(wellenlos, Teil 1 von 4). Review:
[`review-slice-meldungscodes-kennungsfreie-ausgaben`](review-slice-meldungscodes-kennungsfreie-ausgaben.md)
(1 HIGH F-1, 1 LOW F-2, 1 INFO F-3). Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md).

**Gegenstand:** `git diff 47a204b4 HEAD` — vier Commits: Implementierung `e2b206c6`, Gate `df0a440f`,
Review `e1037bd9`, Fixrunde `1cead04f` (alle mit `ADR-0144`, kein `SPEC-*`/`ARC-*` im Betreff).
Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Mutationen liefen
auf einer `git archive`-Kopie im Scratchpad (Änderungen per `sed … Datei > Kopie` und `cp`, kein
`-i`); `git status --short` im Echtrepo war nach jedem Schritt leer. Das `:dev`-Image habe ich vor den
Läufen mit `make image` (Exit 0) frisch gebaut. Nicht gepusht.

## 1. Eigene Sensor-Belege (ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Gedruckter Beleg |
|---|---|---|
| `make gates` | 0 | `d-check: 1593 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s)`; `ausgabe-kennungen-check: keine interne Kennung in Ausgabe-Literalen von 123 Go-Dateien und 4 Skripten`; `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `generated-sync: OK`; a-check `gesamt: 0 Befund(e)`; `make -n gates` nennt das Gate (1 Treffer) |
| `make test` | 0 | alle Pakete `ok`, auch `tools/schema/rolloutguard` |
| `make test-store` | 0 | `DB-Adapter-Coverage: 82.99%`, `db-coverage: OK` |
| `make test-integration` | 0 | letzte Zeilen: `E2E-Abdeckungstabelle unverändert`, `Lauf abgeschlossen`; `git status` danach leer (`docs/user/e2e-abdeckung.md` byte-gleich) |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | `run-schema-rollout-guard-test: OK — alle Belege real erbracht (…)` |
| `make test-ausgabe-kennungen-check` | 0 | `alle 53 Fälle bestanden` (erwartet 53) |
| `make ausgabe-kennungen-check` | 0 | wie oben |
| `make handbuch-public-doc-check` / `make sdk-public-doc-check` | 0 / 0 | `3 Nutzerdokumenten` / `unter sdks` ohne Kennung |
| `make fmt-check` | 0 | 323 Go-Dateien formatiert |
| `make kommentar-kennungen DIFF=47a204b4` | 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | 0 | `26 Zeilen stimmen` |
| `make docs-check` / `make commit-traceability` | 0 / 0 | 0 Befunde / `OK — 5 Commit(s)` |
| `make doc-commits RANGE=47a204b4..HEAD`, `make doc-immutable RANGE=47a204b4..HEAD` | 0 / 0 | je 0 Befunde |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |

**Realer `diagnose`-Lauf (gemessen, nicht übernommen).** Wegwerf-PostgreSQL 18 (Netz und Container
im Lauf angelegt und entfernt), `bash tools/schema/apply-rollout.sh`, Quelle `demo` mit Heartbeat-Zeile,
dann `docker run … ghcr.io/pt9912/pg-change-feed:dev diagnose`. Gedruckt (Auszug): `Betriebsstatus:
Lebenszeichen vor 0.331s`, `Fehlerzustand: keiner (Normalbetrieb)`, `CDC-Abstand cdc_capture_lag:
0.000s`, `Blockierender Consumer: kein Blocker (kein Consumer hat je gegen diese Quelle bestätigt)`,
`Speicherverbrauch cdc_storage_bytes: 0 Bytes`; mit `error_class = 'schema'`: `Fehlerzustand: schema`.
Keine Zeile trägt eine Kennung; die Zeilenform stimmt mit dem Handbuch-Beispiel überein. Die
Blocker-Zeile mit Namen und Kennung (`Orders-Reader (cli-e2e-consumer), …`) habe ich in diesem
Wegwerf-Lauf **nicht** erzeugt (kein bestätigter Consumer); sie ist durch Code
(`%s (%s)` in `wiring.go`) und die grüne Erwartung des `make test-integration`-Laufs
(`Blockierender Consumer: .*\($CLI_CONSUMER\), bestätigte Position …`) getragen. Der Fehlertext der
Konfigurationsdatei lautet gedruckt: `… trägt den Schlüssel "capture_dsn" — Zugangsdaten bleiben env-var-exklusiv`.

**Datenbank-Metadaten (Review F-3, jetzt nachgemessen).** Am selben Wegwerf-Rollout:
`pg_description` für Klassen im Schema `cdc` → 0, für Funktionen im Schema `cdc` → 0, gesamt mit
Kennungsfilter → 0; `\d+ cdc.source` zeigt leere `Description`-Spalten. Die 16 `description:`-Felder von
`tools/schema/schema.yaml` kommen beim Betreiber nicht an; „bleibt als Quelltext-Dokumentation" folgt der
Regel der ADR. F-3 ist damit **gemessen** statt übernommen.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Alle DoD-Zeilen stehen im Plan auf `[ ]`; ich setze keinen Haken (Planner-Vorrecht).

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| A | Ausgaben kennungsfrei (14 Go-Zeilen, 3 `echo`-Zeilen, Erwartungen, Handbuch, `description:` entschieden) | **bestätigt** | Suchlauf-Block 26/26 am Parent und am Diff; Vollläufe `make test`, `make test-store`, `make test-integration`, Wache-Test grün; realer `diagnose`-Lauf und `pg_description`-Messung (oben); Handbuch 1.89 gelesen (Hinweissatz und Platzhalter weg, Beispiel `Orders-Reader (cli-e2e-consumer)` passt zu `%s (%s)`, Historienzeile in Betreibersicht ohne Kennung, `make handbuch-public-doc-check` Exit 0) |
| B | Gate `ausgabe-kennungen-check` in `GATE_CHECKS`, Exit 0 am Ist-Stand, Tabellentest, Vertrag, README-Zeilen, Mutationsprobe | **bestätigt** | `harness/mk/doc-gate.mk` (`GATE_CHECKS +=`), `make gates` Exit 0, 53 Fälle grün; Vertrag `harness/sensors/ausgabe-kennungen-check.md` mit sechs benannten Grenzen; Mutationsprobe von mir erprobt (§3); `harness/README.md` Gate-Zeile, Werkzeug-Zeile und Gates-Zeile gelesen |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf Exit 0 |
| 4 | Review durchgeführt | **bestätigt mit Einschränkung** | Report liegt vor; F-1 durch Fixrunde behoben und von mir mutiert (§3); ein Re-Review fand nicht statt, zur Notwendigkeit siehe §4 |
| 5 | §3.13-Suchlauf | **bestätigt** | `make suchlauf-nachmessen` Exit 0, 26 Zeilen; Plan-Feld trägt Gefundenes und Nichtgefundenes je Träger |
| 6 | Closure-Notiz (§7) | **korrekt offen** | Plan §7 trägt noch Platzhalter `—` |
| 7 | Beobachtungs-Register | **korrekt offen** | Fortschreibung bei Closure (§6) |
| 8 | Risiken §6 mit Ausgang | **korrekt offen** | alle fünf `(bei Closure)`; Vorschläge §7 |

## 3. Mutationen (einzeln, Kopie im Scratchpad, Echtrepo unberührt)

| # | Mutation | Erwartet | Gesehen |
|---|---|---|---|
| M0 | Ausgangslage Kopie | Exit 0 | Tabellentest 53 Fälle grün, Gate Exit 0 |
| M1 | `skip_sh` wirkungslos (`'^NEVER$'`) | Tabellentest rot (4 Fälle) | **rot, genau 4 `FEHLER`-Zeilen**: Shell-Kommentar mit `echo` und Kennung, je eingerückt und nicht eingerückt, je in den Wurzeln `examples` und `tools/schema` (2 × 2 Fälle) |
| M2 | `skip_go` wirkungslos | rot | **rot** (1 Fall: eingerückte Kommentarzeile mit Literal und Kennung) |
| M3 | Kennung ins Go-Literal von `wiring.go:2232` (`Betriebsstatus (LH-FA-ADM-002): …`) | Gate rot mit Datei:Zeile | **Exit 1**, `internal/bootstrap/wiring.go:2232:` gedruckt |
| M4 | dieselbe Kennung als Kommentarzeile davor, Literal sauber | grün | **Exit 0** |
| M5 | Kennung in `echo` von `tools/schema/rollout.sh:139` (`Vorlauf (ADR-0114)`) | rot | **Exit 1**, `tools/schema/rollout.sh:139:` gedruckt; Rücknahme per `cp` → Exit 0 |

Hinweis zu M1: die Fixrunde bindet die Shell-Kommentar-Ausnahme je Skript-Wurzel; die vier roten Fälle
liegen an beiden Wurzeln (`examples`, `tools/schema`). Die Variable `skip_sh` gilt für beide Wurzeln
(`ausgabe-kennungen-check.sh:106`); eine wurzelweise getrennte Mutation habe ich nicht gefahren.

## 4. Fixrunde `1cead04f` — eigene Code-Lesung ohne Re-Review

`git diff e1037bd9 1cead04f`: drei Dateien, 20 Einfügungen, 2 Löschungen — Tabellentest
(`run-ausgabe-kennungen-check-tests.sh`: sechs Fälle zur Shell-Kommentar-Ausnahme je Wurzel, zwei
„benannte Grenze"-Fälle), Sensor-Vertrag (Grenze 5, die zwei einzeiligen Lücken aus F-2) und
`harness/README.md` (Gate-Zeile nennt dieselben Grenzen). **Das Gate-Skript
`tools/harness/ausgabe-kennungen-check.sh` ist in der Fixrunde unverändert.** Die Grenz-Fälle dokumentieren
F-2 als Exit 0 und färben sich bei einer künftigen Schließung rot (so im Vertrag zugesagt). Die
Fixrunde ändert keine Produktionslogik, keinen Ausgabetext, keine Verdrahtung.

**Ist ein weiterer Reviewer-Durchgang Bedingung des Verdikts? Nein.** Begründung: die Fixrunde ist
reine Test-/Vertragsergänzung, ihr Zweck (Bindung von `skip_sh`) ist von mir mit M1 unabhängig erprobt
und färbt rot, das Gate-Skript selbst steht byte-gleich zum bereits reviewten Stand. Der Planner kann
die DoD-Zeile „Review durchgeführt" setzen, sobald er die Fixrunde im Closure-Bericht nennt.

## 5. Entscheidungs-Konformität ([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md))

- **Festlegung 9 T1:** Ausgaben kennungsfrei, keine Codes (kein Code-Muster, keine Tabelle im Diff) — bestätigt.
  Die 28 `Fehlerklasse <klasse>: `-Literale unverändert (Suchlauf diff 28), `description:`-Felder
  nach der Regel der ADR entschieden (nicht sichtbar → bleibt, gemessen).
- **Festlegung 10 Wächter:** ERE ohne `-P`, zweistufig, fail-closed (Exit 2 bei Lesefehler, fehlender
  Wurzel, leerem Gegenstand), Gegenstand wie festgelegt, Läufer unter `tools/harness/` ausgenommen.
  Die in der ADR nur *hergeleitete* Mutationsprobe ist jetzt erprobt (M3 bis M5: Stelle `wiring.go:2232`
  und `rollout.sh:139`, Instanz echter Gate-Lauf, Farben rot/grün/rot).
- **Reichweite:** Diff berührt weder `AGENTS.md` noch `.claude`/`.harness`, keine Paketversion, kein Release, keine
  Workflow-Datei, keine Fehlerklassen-Semantik, keine Quellcode-Kommentare (`kommentar-kennungen DIFF` leer),
  `sdks/` unberührt (Suchlauf 0). Keine `Accepted`-ADR im Diff (`make doc-immutable` Exit 0).

## 6. Plan-vs-Code-Diff

`git diff 47a204b4 HEAD --stat`: 16 Dateien, deckungsgleich mit Plan §3: `wiring.go` (24), `config_file.go`,
`guard.go`, `rollout.sh`, vier Erwartungs-Träger (`diagnose_test.go`, `run-integration-tests.sh`,
`run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md`), Handbuch, Gate (Skript,
Tabellentest-Läufer, `doc-gate.mk`, Sensor-Vertrag, `harness/README.md`), Plan, Review. Kein unbenannter
Nebeneffekt. Das Plan-Feld nennt `guard_test.go` als Träger: im Diff unverändert (der Plan führt die
verbleibende Kennung als Test-Kommentar `guard_test.go:123` — Läufer/Test-Kommentar, außerhalb der Reichweite).

## 7. CI ([`AGENTS.md`](../../AGENTS.md) §3.10) und Vorschläge für die Closure

**CI.** Der T1-Stand ist nicht gepusht; kein Workflow berührt, §3.10 löst keine Pflicht aus. Der erste
Lauf nach dem Push trägt das neue Gate erstmals auf dem Runner: `ci.yml` ruft `make gates` (Zeile 67), das
Gate braucht nur `bash`, `git`, `grep` (mit `-a`), `find` (`-print0`), `mktemp`, `sed` — auf dem
Ubuntu-Runner vorhanden, nichts Weiteres fehlt in `ci.yml`. Ob der Lauf dort grün ist, ist **nicht gemessen**
(Erwartung, kein Beleg).

**Register (melden, nicht ändern).**
- `BEO-PGC/intern-kennungen-in-ausgelieferten-texten`: Fortschreibung — Gate `ausgabe-kennungen-check` ist
  Fangnetz für Programm-Ausgaben mit den Grenzen aus dem Sensor-Vertrag (mehrzeiliges Raw-String, Heredoc,
  Raw-String mit inneren Anführungszeichen, Rune-Literal `'"'`, nachgestelltes Kommentar-Zitat,
  Block-Kommentar); Codes bleiben T2 bis T4.
- **Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) fällig (Planner):**
  [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) trägt den Pfad
  `planning/open/slice-meldungscodes-kennungsfreie-ausgaben.md` im Code-Span an zwei Stellen (Kopf Zeile 19,
  §Geschichte Zeile 323); der Plan liegt in `in-progress/` und geht nach `done/`. Commit mit `ADR-0073` in der
  Message und einer Zeile in der §Geschichte-Tabelle.
- Handbuch 1.89 (Stichprobe): Zeile in Betreibersicht ohne Kennung, `docs/user/version.md` unverändert,
  Gate Exit 0.

**§6-Ausgang-Vorschläge.**
1. *Ausgaben-Stabilität für Betreiber:* **nicht eingetreten** — Handbuch-Diagnose und -Rollout sagen keinen
   Anker mit der entfallenden Klammer zu (Stichprobe der geänderten Stellen; kein Voll-Lesen des Handbuchs).
2. *Tests erwarten Texte:* **nicht eingetreten** — `make test-integration`, `make test-store`, Wache-Test und
   `make test` grün am Stand `1cead04f`.
3. *Datenbank-Metadaten:* **entfallen** — `pg_description` 0, gemessen (§1).
4. *Wächter mit Falsch-Positiven/-Negativen:* **eingetreten und benannt, nicht offen** — Grenzen im Vertrag und
   als Tabellenfälle gebunden; Mutationsprobe erprobt. Optional ein Register-Vermerk (siehe oben).
5. *Kollision mit parallelen Handbuch-Arbeiten:* **entfallen** — kein Konflikt im Diff.

## 8. Verdikt

**Bestanden.** Alle DoD-Liefer-Punkte (A, B) und `make gates` sind mit eigenen, gedruckten Belegen
bestätigt; die vom Review übernommenen Messwerte (Integration, Store, Wache-Test, Diagnose, `pg_description`)
sind nachgefahren. Keine Nachbesserung nötig.

**Bedingungen:** keine zwingende. Vor dem `git mv` nach `done/`: Closure-Notiz und Register (offene
DoD-Zeilen), Zitat-Korrektur der [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) (§7).

**Offene Punkte für den Planner:** (1) Zitat-Korrektur [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) `open/`→`done/`; (2) Register
`intern-kennungen-in-ausgelieferten-texten` fortschreiben; (3) Risiko-Ausgänge §6 (Vorschläge oben); (4) CI-Lauf nach dem Push beobachten.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
