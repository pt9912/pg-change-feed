# Review-Report: slice-meldungscodes-kennungsfreie-ausgaben — 2026-10-03

**Review-Art:** Code (gegen Plan, Entscheidung und Hard Rules)

**Gegenstand:** `git diff 47a204b4 HEAD` (Implementer-Commits `e2b206c6`, `df0a440f` plus Plan)

**Skill:** `.harness/skills/reviewer.md` @ HEAD `df0a440f`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- Slice-Plan `slice-meldungscodes-kennungsfreie-ausgaben` (in-progress)
- [ADR-0144](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), insbesondere Festlegung 10
- Architect-Verdikt `architect-verdict-meldungscodes-statt-interner-kennungen`
- [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen)
- `AGENTS.md` §3.1, 3.2, 3.6, 3.7, 3.9, 3.12, 3.13, 3.15, §4

---

## Findings

### F-1 — Shell-Kommentar-Ausnahme des Wächters ist im Tabellentest nicht gebunden

- `kategorie`: HIGH
- `quelle`: Reviewer-Skill, HIGH-Klasse „Zusage ohne Bindung an ihre Eingabeseite"; Plan-DoD (Mutationsprobe der Kommentar-Ausnahme)
- `pfad`: `tools/harness/run-ausgabe-kennungen-check-tests.sh:101`, `tools/harness/ausgabe-kennungen-check.sh:43`
- `befund`: Mutation an einer Kopie im Scratchpad: `skip_sh` auf ein nie passendes Muster gesetzt, Tabellentest bleibt bei 43 von 43 grün. Der einzige Fall zur Shell-Kommentar-Ausnahme (`# Ausgabe (…) im Kommentar`) trägt kein `echo`/`printf` und wird schon vom Shell-Muster nicht getroffen, die Ausnahme wird nie ausgeübt. Direkter Lauf gegen eine Kopie mit der Zeile `# echo "<Kennung>"` unter `examples/` endet korrekt mit Exit 0, aber nur die Produktions-Logik, nicht der Test, bindet das. Die Go-Seite ist gebunden (gleiche Mutation an `skip_go` färbt den Test rot).
- `verifizierbar`: ja — Mutation von `skip_sh` an einer Kopie, `make test-ausgabe-kennungen-check`
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-2 — Go-Muster übersieht Kennung in Raw-String mit Anführungszeichen im Literal und nach Rune-Literal `'"'`; nicht unter den Grenzen

- `kategorie`: LOW
- `quelle`: [ADR-0144](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10 (Grenzen „im Sensor-Vertrag zu nennen")
- `pfad`: `harness/sensors/ausgabe-kennungen-check.md` §Grenze, `tools/harness/ausgabe-kennungen-check.sh:40`
- `befund`: Gemessen an einer Kopie: `` var s = `{"a": "<Kennung>"}` `` (einzeiliges Raw-String-Literal mit inneren `"`) und `y := "<Kennung>"` hinter einem Rune-Literal `'"'` enden mit Exit 0. Beide sind einzeilig, also nicht die in Grenze 1 genannte Klasse; im Bestand gibt es keinen Treffer (loses Gegenmuster über den echten Baum: 0), die Lücke ist nur nicht benannt.
- `verifizierbar`: ja — Kopie mit den beiden Zeilen
- `klasse`: Grenze nicht benannt

### F-3 — Datenbank-Beschreibungen (16 `description:`): Messung nicht nachgefahren

- `kategorie`: INFO
- `quelle`: [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) / `AGENTS.md` §3.12
- `pfad`: Plan §3, Zeile „Datenbank-Metadaten"
- `befund`: Die Messung „pg_description 0" stammt aus dem Implementer-Lauf und ist hier **übernommen**, nicht nachgemessen. Schwache Stützen gefunden: die d-migrate-Erzeugnisse unter `.tmp/schema-rollout` (`down.sql`, `plan.yaml`, `rollout-precheck.yaml`) tragen keine `COMMENT ON`-Anweisung, sie sind aber keine Aufwärts-DDL und beweisen nichts über die Datenbank.
- `verifizierbar`: ja — Wegwerf-PostgreSQL, `bash tools/schema/apply-rollout.sh`, Zählung in `pg_description`
- `klasse`: Übernommener Messwert

## Geprüft ohne Befund (Schwerpunkte)

- **(a) Ausgaben, Zeile für Zeile (`git diff --word-diff`):** alle `+`-Zeilen in `wiring.go` (12), `config_file.go`, `rolloutguard/guard.go`, `rollout.sh` (3) gelesen. Sachaussage erhalten, Format-Verben (`%s`, `%.3fs`, `%.0f`, `%d`) unverändert, keine Klammer-/Interpunktionsreste, kein verlorenes Leerzeichen (insbesondere `rollout.sh`: `Blocker - --execute`, `Vorlauf - View-…`). Der Fehlertext in `config_file.go` endet sauber auf `env-var-exklusiv`.
- **Erwartungen:** `diagnose_test.go` (7), `run-integration-tests.sh` (12, davon drei `grep -qE` ohne maskierte Klammer, korrekt gegen `Speicherverbrauch cdc_storage_bytes:` und `Blockierender Consumer: .*\(<Name>\)`), `run-schema-rollout-guard-test.sh` (2) stimmen mit den neuen Ausgaben überein; `make test` Exit 0. `make test-integration` und `run-schema-rollout-guard-test.sh` nicht gefahren: **übernommen** (Implementer: grün).
- **Handbuch 1.89:** Version und Stand hochgezählt, Historienzeile in Betreibersicht ohne Kennung, Beispiel `Orders-Reader (cli-e2e-consumer), …` entspricht `%s (%s)`, Hinweissatz und Platzhalter entfernt, `docs/user/version.md` unverändert, `make handbuch-public-doc-check` Exit 0.
- **(c) Wächter:** `tools/harness/ausgabe-kennungen-check.sh` vollständig gelesen: ERE ohne `-P`, Muster wörtlich nach Festlegung 10, Gegenstand als aufgezählte Menge, fail-closed (fehlende Wurzel, `find`-Fehler, `grep`-Exit ≥ 2, leerer Gegenstand → Exit 2 mit gebundenem Meldungstext), keine `cmd | grep -q`-Pipes hinter Erzeugern, Trap/Scratch mit `:?`-Schutz. Eigene Gegenbeispiele: `Errorf`-Literal, `slog`-Attribut, Slice-Literal über mehrere Zeilen, Konstante außerhalb eines Aufrufs, mehrzeilige `+`-Verkettung, Kommentar mit Zitat (laut, benannt) → Exit 1 bzw. wie dokumentiert; Struct-Tag, Test-Datei, erzeugte Datei → grün. Falsch-Positiv-Quote am echten Baum 0.
- **Mutationen einzeln nachgefahren (Kopie per `git archive`, Echtrepo unberührt):** `skip_go` wirkungslos → rot; Backtick aus dem Muster → rot (Fall „Raw-String in Backticks"); Test-Ausnahme aus `find` → rot; `*.pb.go`-Ausnahme → rot; Lesefehler-Zweig `grep` → rot; `find`-Fehlerzweig → rot; leerer Gegenstand erlaubt → rot; `-a` entfernt → rot; Gegenstandsliste `tools/schema` aus `go_roots` → rot; `examples` aus `sh_roots` → rot. Einzige überlebende Mutation: `skip_sh` (F-1).
- **Verdrahtung:** `harness/mk/doc-gate.mk` (`GATE_CHECKS +=`), `make -n gates` nennt das Gate, `make gates` Exit 0 (ungepiped, Ausgabezeile des Gates im Log), `harness/README.md` Gate-Index und Werkzeug-Zeile, Sensor-Vertrag mit Grenzen vollständig. Gate-Aufnahme ab Ist-Stand 0 durch die ADR gedeckt (`AGENTS.md` §3.6), `LH-QA-OPS-001` im Lastenheft vorhanden.
- **(d) Läufe selbst:** `make test-ausgabe-kennungen-check` (43/43), `make ausgabe-kennungen-check`, `make gates`, `make test`, `make fmt-check`, `make kommentar-kennungen DIFF=47a204b4` (0), `make suchlauf-nachmessen` (26 Zeilen stimmen), `make docs-check`, `make sdk-public-doc-check`, `make handbuch-public-doc-check`, `make commit-traceability`: alle Exit 0. `make test-store`, `make test-integration` und der reale `diagnose`-Lauf: nicht gefahren, **übernommen** aus dem Implementer-Bericht im Plan.
- **(e) Reichweite:** der Diff berührt weder `AGENTS.md`, `.claude`, `.harness`, noch Paket-Versionen oder Fehlerklassen-Semantik, noch Quellcode-Kommentare (`kommentar-kennungen DIFF` leer).
- **(f) Plan-Klärungen:** die Zahlen 97 gegen 142 und 9 gegen 12 sind nachgemessen: die Suchlauf-Zeilen (97 am Parent, 9 plus 3 maskierte) stimmen per `make suchlauf-nachmessen`; die 142 des Vorläufer-Plans bleiben im Plan als „nicht reproduziert" gekennzeichnet, die Begründung ist plausibel.

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/`, `tools/schema/` (Ausgabe-Literale)
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh`, `tools/harness/run-schema-rollout-guard-test.sh`
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md`
- geprüft, ohne Befund: `harness/mk/doc-gate.mk`, `harness/README.md`, `harness/targets/schema-rollout.md`
- geprüft, ohne Befund: `docs/plan/planning/in-progress/slice-meldungscodes-kennungsfreie-ausgaben.md` (Suchlauf-Feld)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite · Grenze nicht benannt · Übernommener Messwert

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (HIGH nach Klassifikation des Skills). Der Fix ist klein: ein Tabellenfall, der an der Shell-Kommentar-Ausnahme rot werden kann. Die Produktionslogik selbst ist korrekt. Die DoD-Zeile „Review durchgeführt" bleibt offen (Fixrunde nötig) und wird beim Nachlauf gesetzt.

**Übergabe:** F-1 und F-2 an den Implementer; F-3 an Planner/Verifier (Messung nachfahren oder als übernommen führen). Dieser Report ersetzt keine Verifikation.
