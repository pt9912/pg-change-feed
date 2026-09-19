# Slice rtm-reste-sst-cfg-por: die restlichen sechs RTM-Waisen aus der Neun-Klassifikation

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — kein Closure-Kriterium jenseits der eigenen DoD.

**Bezug:** [`LH-FA-SST-001`](../../../../spec/lastenheft.md),
[`LH-FA-SST-005`](../../../../spec/lastenheft.md),
[`LH-FA-CFG-006`](../../../../spec/lastenheft.md),
[`LH-QA-PER-004`](../../../../spec/lastenheft.md),
[`LH-QA-POR-001`](../../../../spec/lastenheft.md),
[`LH-QA-POR-002`](../../../../spec/lastenheft.md);
[`ADR-0105`](../../adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md) (neu).

**Berührte Spec-Stellen:** keine — alle sechs Anforderungen sind bereits
real erfüllt (siehe je Punkt in §1); dieser Slice macht sie nur für
`make doc-trace` sichtbar bzw. ergänzt einen bislang fehlenden
Struktur-Beleg (`LH-FA-CFG-006`).

**Verantwortlich:** — (wellenlos, direkt umgesetzt, siehe Autor).

**Autor:** Implementer-Agent, direkt beauftragt durch den Auftraggeber im
Anschluss an `bench-schwellen-per-001-002-003` ("Jetzt noch die
restlichen 6 Anforderungen abdecken"). Kein separater
Priorisierungs-Schritt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Nach `bench-schwellen-per-001-002-003` (löste `LH-QA-PER-001`/
`002`/`003`) verblieben sechs der ursprünglich neun "andere Ebene als
`test-integration`"-RTM-Waisen. Dieser Slice schließt alle sechs:

1. `LH-FA-SST-001` (Schnittstelle zur PostgreSQL-Quelle) — Tag-only:
   `TestE2ECaptureFlow`s Doc-Kommentar (`test/integration/integration_test.go`)
   trägt die Kennung; der Test selbst läuft bereits Ende-zu-Ende über
   Logical Replication, nur ohne eigenen Kennungs-Tag.
2. `LH-FA-SST-005` (Vorbereitung späterer HTTP-/gRPC-API) — Tag-only:
   die bestehende "HTTP-API-Rundlauf"-Phase
   (`tools/harness/run-integration-tests.sh`) belegt bereits real, dass
   die spätere API ohne Änderung am internen CDC-Modell entstand.
3. `LH-QA-PER-004` (Latenz bis zur CDC-Verfügbarkeit) — Tag-only: die
   bestehende "Lasttest-Beleg cdc_capture_lag"-Phase misst bereits
   dieselbe Latenz, die `LH-QA-PER-004` über `LH-FA-ADM-004` verlangt
   (`SPEC-013`-Schwelle, siehe `ADR-0054`/`ADR-0104`).
4. `LH-FA-CFG-006` (keine Anwendungscode-Anpassung) — neuer struktureller
   Beleg: die bestehende "SQL-Administration Live-Reload
   (enable)"-Phase nimmt jetzt einen DDL-Fingerabdruck (Spaltenliste +
   Trigger-Anzahl der Quelltabelle) vor und nach `cdc.enable_table`
   und vergleicht — `cdc.enable_table` liest/schreibt ausschließlich
   `cdc.active_table`, nie die Quelltabelle selbst.
5. `LH-QA-POR-001`/`LH-QA-POR-002` (PostgreSQL-Major-Versionen /
   primäre Zielplattform Linux) — neue generierte Coverage-Dimension
   `docs/user/ci-matrix-abdeckung.md` (`ADR-0105`): fragt reale,
   bereits abgeschlossene GitHub-Actions-Läufe von `e2e.yml`
   (PostgreSQL-17/18-Matrix) und `ci.yml` (Linux-Plattform-Assertion)
   über die GitHub-REST-API ab, zitiert Lauf-ID/Commit/URL als Beleg.

**Ausdrücklich NICHT in diesem Slice:**

- Ein neuer Sensor/Gate für `LH-QA-POR-001`/`002` in `make gates` —
  bleibt strukturell ausgeschlossen (`AGENTS.md` §3.10: GitHub-Actions-
  Runner-Verhalten ist von einem Docker-only-Sensor nicht lokal
  nachbildbar).
- Die zwölf real bereits über bestehende E2E-Phasen mitbelegten, aber
  tag-losen RTM-Waisen aus der ursprünglichen Klassifikation
  (`LH-FA-CON-001`…`005`, `LH-FA-DAT-002`/`003`/`005`, `LH-FA-REA-001`,
  `LH-FA-RET-001`, `LH-QA-REL-003`/`004`) — eigener, unsanierter
  Bestandslücken-Vorgang (siehe `harness/README.md` §`make doc-trace`),
  keine Benchmark- oder CI-Matrix-Sichtbarkeitsfrage wie die sechs hier.

## 2. Definition of Done

- [x] `LH-FA-SST-001` zeigt in `make doc-trace` `ok` (E2E-Coverage,
      Tag in `TestE2ECaptureFlow`).
- [x] `LH-FA-SST-005` zeigt in `make doc-trace` `ok` (E2E-Coverage,
      Tag in der HTTP-API-Rundlauf-Phase).
- [x] `LH-QA-PER-004` zeigt in `make doc-trace` `ok` (E2E-Coverage,
      Tag in der Lasttest-Beleg-Phase).
- [x] `LH-FA-CFG-006` zeigt in `make doc-trace` `ok` — neuer,
      real ausgeführter DDL-Fingerabdruck-Vergleich in der
      SQL-Administration-Live-Reload-Phase (nicht nur ein Tag ohne
      neue Prüfung).
- [x] `LH-QA-POR-001`/`LH-QA-POR-002` zeigen in `make doc-trace` `ok`
      (CI-Matrix-Coverage, `ADR-0105`, real gegen die GitHub-REST-API
      geprüft — Lauf-ID/Commit/URL im generierten Beleg).
- [x] `make test-integration` grün, regeneriert
      `docs/user/e2e-abdeckung.md` mit den vier neuen/erweiterten
      Kennungs-Tags.
- [x] `make doc-ci-matrix` grün, erzeugt
      `docs/user/ci-matrix-abdeckung.md`.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] Doku-Update: `harness/README.md` (`make doc-trace`-Zeile mit der
      real gemessenen finalen Waisenzahl nachgezogen), `docs/plan/adr/README.md`
      (`ADR-0105`-Index-Zeile).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben (falls Review-Funde anfallen).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | `TestE2ECaptureFlow`s Doc-Kommentar trägt zusätzlich `LH-FA-SST-001`. |
| `tools/harness/run-integration-tests.sh` | update | "HTTP-API-Rundlauf"-`abdeckung_declare` trägt zusätzlich `LH-FA-SST-005`; "Lasttest-Beleg cdc_capture_lag"-`abdeckung_declare` trägt zusätzlich `LH-QA-PER-004`; "SQL-Administration Live-Reload (enable)"-Phase bekommt einen DDL-Fingerabdruck-Vergleich (Spalten + Trigger, vor/nach `cdc.enable_table`) und trägt zusätzlich `LH-FA-CFG-006`. |
| `docs/plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md` | neu | Entscheidung: Docker-gekapselte GitHub-REST-API-Abfrage + generierte Coverage-Datei für `LH-QA-POR-001`/`002`. |
| `docs/plan/adr/README.md` | update | Index-Zeile für `ADR-0105`. |
| `tools/harness/ci-matrix-abdeckung.sh` | neu | Fragt `e2e.yml`/`ci.yml` über die GitHub-REST-API ab (Container-gekapselt, `TOOLCHAIN_IMAGE` wiederverwendet), schreibt `docs/user/ci-matrix-abdeckung.md`. |
| `Makefile` | update | Neues Ziel `doc-ci-matrix` (kein Gate, braucht Netz). |
| `.d-check.yml` | update | `trace.coverage` um einen dritten Eintrag (`docs/user/ci-matrix-abdeckung.md`, Label „CI-Matrix") ergänzt. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | von `make test-integration` neu geschrieben. |
| `docs/user/ci-matrix-abdeckung.md` | Erzeugnis | von `make doc-ci-matrix` neu geschrieben. |
| `harness/README.md` | update | `make doc-trace`-Zeile mit der finalen real gemessenen Waisenzahl. |

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — direkter Auftrag, kein
externer Trigger.

**Rückführungen:**

- `in-progress` → `next` (zu groß): entfällt — Umfang bereits real
  umgesetzt (alle sechs Mechanismen laufen).
- `in-progress` → `open` (blockiert): entfällt aus demselben Grund.

## 5. Closure-Trigger

DoD vollständig + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der CI-Matrix-Beleg (`ADR-0105`) zitiert den letzten **erfolgreichen**
  GitHub-Actions-Lauf, nicht zwingend einen, der zum aktuellen HEAD
  gehört — bei frisch gepushtem, noch nicht durchgelaufenem Commit zeigt
  die generierte Datei auf den Vorgänger-Lauf. **Ausgang:** weiter offen,
  bewusst akzeptiert (`ADR-0105` §Konsequenzen), dasselbe Muster wie
  `docs/user/bench-abdeckung.md` ("stabile Abdeckungs-Deklaration, kein
  Lauf-Beleg für den aktuellen Stand").
- `make doc-ci-matrix` braucht Netz und die öffentliche Erreichbarkeit
  der GitHub-REST-API — schlägt bei Netzausfall mit `exit 1` fehl (kein
  stiller veralteter Erfolg), aber blockiert dann auch keinen Gate-Lauf,
  da außerhalb von `make gates`. **Ausgang:** entfallen als Risiko für
  `make gates` (strukturell nicht Teil davon), weiter offen als
  Bedienungshinweis (dokumentiert in `docs/user/ci-matrix-abdeckung.md`
  selbst).
- Der DDL-Fingerabdruck für `LH-FA-CFG-006` (Spaltenliste +
  Trigger-Anzahl, sortiert nach `ordinal_position`) deckt keine
  Änderung an Constraints/Indizes der Quelltabelle ab — ein `enable_table`,
  das (hypothetisch) einen Index anlegen würde, bliebe unentdeckt.
  **Ausgang:** weiter offen, praktisch irrelevant: `cdc.enable_table`s
  Implementierung schreibt nachweislich ausschließlich in
  `cdc.active_table` (Code-Lektüre, keine DDL-Anweisung im Pfad), der
  Fingerabdruck ist ein zusätzlicher Laufzeit-Beleg, kein einziger
  Wahrheitsanker.
- Drei der sechs Punkte (`SST-001`, `SST-005`, `PER-004`) sind reine
  Tag-Ergänzungen ohne neue Prüfung — sie tragen keinen zusätzlichen
  Regressionsschutz über das hinaus, was die jeweilige Phase bereits vor
  diesem Slice belegte. **Ausgang:** entfallen als Risiko, das ist die
  bewusste Einordnung dieser drei Anforderungen (Beleg liegt bereits vor,
  nur der Kennungs-Tag fehlte).

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Areas „E2E-Test-Harness"
(`tools/harness/run-integration-tests.sh`, `test/integration/`) und neu
„CI-Matrix-Sichtbarkeit" (`tools/harness/ci-matrix-abdeckung.sh`) — beide
GF (Modus-Deklaration `harness/conventions.md`, Default `PGC`/Greenfield).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer zu diesem konkreten Gegenstand.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
