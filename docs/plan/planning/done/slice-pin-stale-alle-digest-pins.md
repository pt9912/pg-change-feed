# Slice pin-stale-alle-digest-pins: Sensor P10 — jede Digest-Pin-Referenz des Baums nächtlich gegen die Registry

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei,
`ls docs/plan/planning/*.md` nennt nur `README.md`), siehe Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht (Modul 6).

**Abhängigkeit:** startet erst, wenn
[`slice-pin-digests-aktualisieren-2026-10-b`](../done/slice-pin-digests-aktualisieren-2026-10-b.md)
in `done/` liegt. Der Sensor färbt den Nachtlauf am Tag seiner Einführung rot,
solange die dort gehobenen Pins driften
([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
Frage (d): erst Hebung, dann Sensor).

**Bezug:** [`LH-QA-POR-001`](../../../../spec/lastenheft.md),
[`LH-QA-POR-003`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Festlegung 1 bis 6, Folgepflicht Slice 2),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7
(Pin-Inventar P1 bis P9, `upstream-drift.yml`; bleibt unverändert),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md)
(Muster: Tabellentest bleibt Werkzeug).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum; zu belegen durch
den Suchlauf in §3, der in `spec/` keinen Treffer der Pin-Inventar-Begriffe
findet — Erwartung, nicht gemessen).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, direkt beauftragt (kein Architect: die Entscheidung
liegt mit [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
vor). **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Anlass:** das Pin-Inventar von
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 ist
eine Liste (P1 bis P9) und veraltet mit jedem neuen Pin; von 15 verschiedenen
Referenzen am Stand `28a6242a` standen neun außerhalb der Liste ohne Leser, unter
ihnen der PostgreSQL-17-Pin
([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
Kontext 2 und 3, dort am Stand `28a6242a` gemessen — **übernommen**, hier nicht
nachgemessen; der Implementer misst die Menge am Arbeitsstand nach, §3).

**Ziel:** Ein Sensor P10 prüft jede Referenz der Form
`<image>[:<tag>]@sha256:<digest>` in einer getrackten Datei außerhalb von
`docs/` und `.harness/` nächtlich gegen den Registry-Digest, ohne dass ein
neuer Pin in eine Liste eingetragen wird — Skript, Tabellentest, Make-Target,
Workflow-Schritt, Zeile in `harness/README.md` und Sensor-Vertrag.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Hebung gedrifteter Pins** — Vorgänger-Slice
  [`slice-pin-digests-aktualisieren-2026-10-b`](../done/slice-pin-digests-aktualisieren-2026-10-b.md)
  (liegt vor diesem in `done/`); mischte dieser Slice beides, wäre der erste
  Sensorlauf rot und der Nachweis „P10 findet Drift“ nicht von „die Hebung ist
  unvollständig“ zu trennen. Zeigt die Startmessung einen frischen Drift, wird
  er als eigener Hebe-Commit **vor** dem Sensor-Commit gehoben (derselbe
  Vorgang wie im Vorgänger, Befund in §7), nicht im Sensor-Commit.
- **Eine Pflicht zur Drift-Auflösung, ein Alarm- oder Eigentümer-Mechanismus** —
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Frage (d); die Beobachtung
  [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
  bleibt ohne Änderung.
- **Der Tag-Wechsel** (neuerer Major eines Tags) — Grenze von
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Festlegung 6; der Sensor meldet Digest-Drift am gepinnten Tag.
- **Eine Änderung an P1 bis P9** — sie bleiben unverändert bestehen, Überschneidungen
  melden denselben Drift doppelt (Festlegung 1 der ADR); die Sensoren
  `make image-stale` und `make pin-stale-*` und ihr Verhalten bleiben
  byte-gleich in Ausgabe und Exit-Codes (das Teilen des Vergleichs ist
  Refactor, kein Verhaltenswechsel; Beleg §2 Liefer-Punkt 2).
- **Ein Gate** — P10 ist wie P1 bis P9 advisory und steht nicht in
  `GATE_CHECKS`/`make gates` (braucht Netz,
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7);
  ein blockierendes Gate wäre eine eigene ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Sensor (Liefer-Punkt 1).** Ein Skript unter `tools/harness/`
      (Arbeitsname `pin-stale-all.sh`) zählt die Referenzen auf
      ([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
      Festlegung 1: `git grep` mit dem Muster der ADR, Wertemenge dedupliziert,
      `docs/` und `.harness/` ausgenommen, kein Fundort-Verzeichnis), vergleicht
      jede gegen den Index-Digest (Festlegung 2), nimmt für Pins ohne Tag
      `:latest` (Festlegung 3) und endet fail-open je Referenz (Festlegung 4:
      `DRIFT` → Exit 1, `UNBESTIMMT` ohne `DRIFT` → Exit 2, sonst 0; leerer
      Gegenstand Exit 2). Der Digest-Vergleich ist **wiederverwendet, nicht
      dupliziert**: die Vergleichslogik von `tools/harness/pin-stale.sh` liegt
      in einer gemeinsamen, von beiden Skripten eingelesenen Datei (Arbeitsname
      `tools/harness/lib-pin-compare.sh`); Ausgabe und Exit-Codes von
      `make image-stale` und `make pin-stale-race|-pgtest|-dmigrate|-acheck`
      sind vor und nach dem Refactor an derselben Registry gemessen **gleich**
      (gedruckte Zeilen beider Stände in §7). Die Meldungsform trägt je Referenz
      `OK`/`DRIFT`/`UNBESTIMMT`, den Fundort (`datei:zeile`) und beide Digests.
      Weder das Skript noch sein Tabellentest noch der Sensor-Vertrag trägt eine
      **vollständige** Referenz `<image>@sha256:<64 Hex>` als Literal (das
      Muster fände sie und meldete sie für immer als `DRIFT`); Fixture-Digests
      entstehen zur Laufzeit.
- [x] **Tabellentest und Mutationsbeleg (Liefer-Punkt 2).** Ein Tabellentest
      (Arbeitsname `tools/harness/run-pin-stale-all-tests.sh`,
      `make test-pin-stale-all`, netzlos über ein Stub-`docker` im
      Wegwerf-Repo im Temp-Verzeichnis) deckt je Zweig der Festlegungen 1 bis 4
      einen Fall mit Meldungstext: Treffer in `.yml`, im Shell-Default und in
      `compose.yaml`; kein Treffer (Digest gleich); Pin mit Tag und Pin ohne Tag
      (gegen `:latest`); `docs/` und `.harness/` ausgenommen; Registry nicht
      erreichbar → `UNBESTIMMT` ohne Abbruch der übrigen Referenzen; `DRIFT`
      neben `UNBESTIMMT` → Exit 1; ein Einzelplattform-Digest wird als `DRIFT`
      gemeldet; leerer Gegenstand → Exit 2. **Mutationsbeleg (Zusage, zu belegen
      durch Läufe an einer Kopie im Scratchpad, nicht an der Datei im
      Arbeitsbaum — [`AGENTS.md`](../../../../AGENTS.md) §3.1):** der
      Implementer mutiert das Skript an **drei benannten Stellen** — das
      `docs/`-Ausschluss-Pathspec (Festlegung 1), die Abbildung „leerer
      Gegenstand → Exit 2“ auf Exit 0 (Festlegung 4) und die `:latest`-Ableitung
      für Pins ohne Tag (Festlegung 3) — je **eine** Mutation, am **Tabellentest**
      als Instanz (Prüfling per Umgebungsvariable auf die Kopie gerichtet, im
      Muster von `GUARD=` in `make test-command-guard`); der Beleg nennt je
      Mutation Stelle, Instanz und die gesehene Farbe (rot erwartet). Jede
      Verallgemeinerung darüber hinaus („alle Zweige tragen“) steht als
      *hergeleitet* ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).
- [x] **Verdrahtung und realer Lauf (Liefer-Punkt 3).** Make-Target
      (Arbeitsname `pin-stale-all`, neben den `pin-stale-*`-Zielen im
      `Makefile`, `.PHONY`, Hilfetext) und `make test-pin-stale-all` laufen; ein
      zehnter Schritt mit `if: always()` in
      `.github/workflows/upstream-drift.yml` ruft dasselbe Target (Name in
      `harness/README.md` §Sensors und im Schritt **identisch**, siehe
      [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
      Festlegung 5; jede neu eingeführte `uses:`-Zeile bliebe SHA-gepinnt mit
      Tag-Kommentar, [`AGENTS.md`](../../../../AGENTS.md) §3.8 — der Schritt nutzt
      `run:`, keine neue Action); Kopf- und Jobname-Kommentare des Workflows
      nennen den neuen Umfang. Eine Zeile in
      [`harness/README.md`](../../../../harness/README.md) §Sensors und die
      Vertragsdatei `harness/sensors/pin-stale-all.md` (Vertrag, Ausgänge,
      Grenzen aus der ADR, Host-Werkzeuge) liegen vor. `make pin-stale-all` am
      Arbeitsstand nach dem Vorgänger-Slice endet mit Exit 0 (gedruckte Zeilen
      und Parent-Kennung in §7, **gemessen**, nicht aus der ADR übernommen). Nach
      dem Push ist `upstream-drift.yml` per `workflow_dispatch` real gelaufen
      (`gh workflow run`, `gh run view`): zehn Schritte, Gesamtlauf `success`
      ([`AGENTS.md`](../../../../AGENTS.md) §3.10; Push gehört nicht zu diesem
      Plan, die Lesung ist Closure-Pflicht).

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [x] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0
      (`make docs-check` liest `harness/README.md`, den Sensor-Vertrag und diese
      Pläne: Kennungen verlinkt).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      ([`review-slice-pin-stale-alle-digest-pins`](../../../reviews/review-slice-pin-stale-alle-digest-pins.md),
      F-1 MEDIUM in der Fixrunde behoben;
      `.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (neue Zeile, Zeile zu `upstream-drift.yml` auf den neuen Umfang),
      [`docs/maintainer/releasing.md`](../../../maintainer/releasing.md) (Satz zum
      Umfang von `upstream-drift.yml`, samt Versionshistorie-Zeile, falls das
      Dokument eine führt) — gemeldete Träger fremder Dateien werden mit der
      Closure nachgezogen (§3 Suchlauf, [`AGENTS.md`](../../../../AGENTS.md)
      §3.13).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Zu tragen: der Zustand von
      [`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
      (Träger `ADR-0146` steht bereits in dessen `state.md`) wird mit dem
      ersten grünen Nachtlauf aktualisiert (Ausgang nur nach gelesenem Lauf).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/pin-stale-all.sh` (Arbeitsname) | neu | Sensor P10: Aufzählung per `git grep`, Vergleich je Referenz, Ausgänge nach Festlegung 4 |
| `tools/harness/lib-pin-compare.sh` (Arbeitsname) | neu | der Index-Digest-Vergleich aus `tools/harness/pin-stale.sh`, von beiden Skripten eingelesen — Wiederverwendung statt Duplikat (Festlegung 5) |
| `tools/harness/pin-stale.sh` | refactor | liest die gemeinsame Datei; Ausgabe und Exit-Codes unverändert (Beleg: gedruckte Zeilen vor und nach dem Refactor, §7); `tools/harness/pin-stale-dcheck.sh` ruft `pin-stale.sh` und ist damit mitgemessen |
| `tools/harness/run-pin-stale-all-tests.sh` (Arbeitsname) | neu | Tabellentest je Zweig der Festlegungen 1 bis 4, Stub-`docker`; Prüfling per Umgebungsvariable übersteuerbar (Mutationsläufe an Kopien) |
| `Makefile` (neben `pin-stale-actions`, Zeilen 68 bis 88), `harness/mk/doc-gate.mk` oder ein bestehender `*.mk` | update | `pin-stale-all` und `test-pin-stale-all`, je Hilfetext, `.PHONY`; **nicht** in `GATE_CHECKS` |
| `.github/workflows/upstream-drift.yml` | update | zehnter Schritt `if: always()`; Kopf-Kommentar (Zeile 5) und `name:` (Zeile 42) auf den neuen Umfang; strukturelle Workflow-Änderung → §3.10 |
| `harness/README.md` (§Sensors, Zeile zu `upstream-drift.yml` und neue Zeile) | update | Zeile für `make pin-stale-all` mit Vertrag und Bindung, `make test-pin-stale-all` im Werkzeug-Block |
| `harness/sensors/pin-stale-all.md` | neu | Sensor-Vertrag: Gegenstand, Ausgänge, Grenzen ([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) §Fitness Function, Grenzen), Host-Werkzeuge (`bash`, `git`, `docker`) |
| `docs/maintainer/releasing.md` (Zeile 371) | update | Satz zum Umfang von `upstream-drift.yml`, Versionshistorie-Zeile (letzte Zeile der Tabelle ist `1.13`) |
| `docs/plan/planning/observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/` | update | Zustand mit dem ersten grünen Nachtlauf (DoD Register-Zeile) |

- **Reihenfolge:** (1) Startmessung am Arbeitsstand: `git grep` mit dem Muster
  der ADR und `make image-stale` samt `make pin-stale-*` — zeigt sie frischen
  Drift, wird er zuerst als eigener Hebe-Commit gehoben (§1); die gedruckten
  Zeilen der P3–P6-Läufe für den Refactor-Vergleich werden **vor** dem
  Refactor gesichert. (2) Refactor `pin-stale.sh` mit Vergleich der Ausgabe;
  (3) Sensor-Skript; (4) Tabellentest, Mutationsläufe an Kopien im Scratchpad;
  (5) Make-Ziele, Workflow-Schritt, README, Vertrag, `releasing.md`;
  (6) `make pin-stale-all` real (Exit 0 erwartet nach dem Vorgänger),
  `make test-pin-stale-all`, `make gates`; (7) nach dem Push
  `workflow_dispatch`.
- **Menge am Arbeitsstand nachmessen:** `git grep -ohE
  '[a-z0-9][a-z0-9./_-]*(:[A-Za-z0-9._-]+)?@sha256:[0-9a-f]{64}' -- . ':!docs'
  ':!.harness'` am Stand der Messung, Zahl der verschiedenen Referenzen und
  Dateien in §7; die 15 und 20 der ADR sind Messwerte am Stand `28a6242a`, kein
  Soll ([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Frage (c)). Findet die Messung eine Referenz, die **kein** echter Pin ist
  (Fixture mit erfundenem Digest in einer getrackten Datei), ist das ein Befund
  für den Planner, kein Anlass für eine Ausnahmeliste im Sensor
  (Festlegung 1: keine Liste von Fundorten).
- **Vergleich je Referenz:** der Maßstab ist `docker buildx imagetools inspect
  <ref ohne @digest> --format '{{.Manifest.Digest}}'`; Pins ohne Tag gegen
  `:latest` (Festlegung 3). Ein Einzelplattform-Digest meldet der Sensor als
  `DRIFT` (Festlegung 2, gewollt).
- **Suchlauf (§3.13, bewegte Eigenschaft: der Umfang des Pin-Inventars — „P1
  bis P9“, „Neun-Achsen“ — und die Aufrufform von `pin-stale.sh`).** Der Parent
  ist `39e27242eb58ed1212b381a15da55c6ee1f5033c` (`git rev-parse HEAD` am
  Planungsstand, vor dem Plan-Commit; nie `HEAD`). Suchraum: ganzer Baum ohne
  `docs/reviews`, `done/`, `observations/`, `.harness/baseline`, die offenen
  Pläne und `docs/plan/adr` (Accepted ADRs halten ihren Stand,
  [`AGENTS.md`](../../../../AGENTS.md) §3.5). Parent-Zeilen gemessen; die
  `diff`-Zeilen tragen die **Erwartung** nach der Arbeit (hergeleitet, nicht
  gemessen):

```suchlauf
39e27242eb58ed1212b381a15da55c6ee1f5033c 4 -n -E 'Neun-Achsen|neun Achsen|P1[–-]P9|P1 bis P9' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
39e27242eb58ed1212b381a15da55c6ee1f5033c 0 -n -E 'Pin-Inventar|pin-stale|upstream-drift' -- spec
39e27242eb58ed1212b381a15da55c6ee1f5033c 5 -n -E 'pin-stale\.sh' -- Makefile harness/mk
diff 0 -n -E 'Neun-Achsen|neun Achsen|P1[–-]P9|P1 bis P9' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline' ':!docs/plan/adr'
diff 0 -n -E 'Pin-Inventar|pin-stale|upstream-drift' -- spec
diff 5 -n -E 'pin-stale\.sh' -- Makefile harness/mk
```

  **Erwartung am `diff`-Stand** (hergeleitet): die vier Parent-Treffer
  (`.github/workflows/upstream-drift.yml` Zeilen 5 und 42, `harness/README.md`
  Zeile 147, `docs/maintainer/releasing.md` Zeile 371) tragen nach dem Nachzug
  den neuen Umfang („P1 bis P10“, „zehn Achsen“ ist die Wortwahl des
  Implementers, solange keines der vier Muster in der ersten Zeile trifft);
  die Aufrufer von `pin-stale.sh` im `Makefile` bleiben fünf, weil das neue Ziel
  `pin-stale-all.sh` aufruft (nicht `pin-stale.sh`). Das Muster liest
  Symbolnamen und Zählwörter; ein Zeilen-Lokator oder eine Zählung in Prosa
  („neun“ in anderer Wortform) trifft es nicht — diese Hälfte bleibt Lese-Handlung
  des Reviewers ([`AGENTS.md`](../../../../AGENTS.md) §3.13 Grenze).

### Implementer-Beleg (Stand der Lieferung, Parent `3eb60de6`)

- **Plan-Nachzug (Abweichung vom Arbeitsnamen):** `pin-stale-all` und
  `test-pin-stale-all` liegen im `Makefile` neben `pin-stale-actions` (kein
  `*.mk`); der Vergleich eines Registry-Aufrufs ist über
  `PIN_COMPARE_TIMEOUT` (Default 60 s im Sensor, ohne Wert unbegrenzt in
  `pin-stale.sh`) begrenzt — Antwort auf das Laufzeit-Risiko aus §6. Die
  Referenz mit Registry-Host samt Port ist eine benannte Grenze im
  Sensor-Vertrag (das Muster der ADR beginnt mit `[a-z0-9]`, ein `:` vor dem
  Namen trennt es).
- **Menge (gemessen, Stand Arbeitsbaum nach Slice 1):** `git grep -ohE '…' -- .
  ':!docs' ':!.harness' | sort -u | wc -l` ergibt 15 verschiedene Referenzen,
  `git grep -lE` ergibt 20 Dateien — gleich den Messwerten der ADR am Stand
  `28a6242a`. Jede der 15 ist ein echter Pin (die Zeilen des realen Laufs nennen
  Fundort und Image); kein Fixture-Fund, keine Ausnahmeliste.
- **Refactor-Vergleich (gemessen):** `make image-stale`,
  `make pin-stale-race|-pgtest|-dmigrate|-acheck|-dcheck` vor und nach dem
  Refactor an derselben Registry gefahren, Ausgabe je Ziel mit `cmp` byte-gleich
  und Exit je Ziel 0 (gedruckte Zeile u. a. `OK          PG_TEST_IMAGE
  (postgres:18-alpine) == sha256:77f58511…`). Die Fehlerpfade (`DRIFT`,
  `UNBESTIMMT` bei Registry-Ausfall, Wert ohne Digest, Variable fehlt) sind mit
  einem Stub-`docker` an der Fassung von `3eb60de6` und der neuen Fassung
  verglichen: Ausgabe und Exit je Fall gleich (1, 0, 2, 2, 2).
- **Realer Lauf (gemessen):** `make pin-stale-all` endet mit Exit 0, Laufzeit
  `real 0m28,845s`, gedruckte Schlusszeile `pin-stale-all: 15 Referenzen — 15 OK,
  0 DRIFT, 0 UNBESTIMMT`; darunter die PostgreSQL-17-Referenz am Fundort
  `.github/workflows/e2e.yml:80` mit Index-Digest.
- **Tabellentest:** `make test-pin-stale-all` endet mit Exit 0, gedruckte Zeile
  `run-pin-stale-all-tests: alle 34 Prüfungen bestanden`.
- **Mutationsbeleg (Instanz: der Tabellentest, Prüfling per `PROG=<Kopie im
  Scratchpad>`, je **eine** Mutation, Eingabeseite des Sensors):**

| Zusage | Stelle der Mutation | Instanz | gesehene Farbe |
|---|---|---|---|
| `docs/` ist ausgenommen (Festlegung 1) | Pathspec `':!docs'` aus dem `git grep`-Aufruf entfernt | `run-pin-stale-all-tests.sh` | rot, Exit 1: vier Prüfungen (Fall „docs/ und .harness/ ausgenommen“ Exit 2 statt 0, Referenz aus `docs/` gemeldet, Zusammenfassung, Meldung „einzige Referenz unter docs/“) |
| leerer Gegenstand ist Exit 2 (Festlegung 4) | `exit 2` hinter der Meldung „leerer Gegenstand“ auf `exit 0` | `run-pin-stale-all-tests.sh` | rot, Exit 1: zwei Prüfungen (Fälle „keine Referenz im Baum“ und „einzige Referenz liegt unter docs/“ Exit 0 statt 2) |
| Pin ohne Tag wird gegen `:latest` verglichen (Festlegung 3) | `target="$image:latest"` auf `target="$image"` | `run-pin-stale-all-tests.sh` | rot, Exit 1: vier Prüfungen (Exit 2 statt 0, Meldung ohne `:latest`, Registry-Aufruf mit `:latest` 0 statt 1, Aufruf ohne Tag 1 statt 0) |

  *Hergeleitet, nicht gefahren* ([`AGENTS.md`](../../../../AGENTS.md) §3.12
  Instanz B): dass die übrigen Zweige (Deduplizierung, `.harness/`-Ausschluss,
  `DRIFT` neben `UNBESTIMMT`, Einzelplattform-Digest) ebenso rot
  färben; für sie liegt eine Prüfung im Tabellentest, aber keine gesehene Farbe.
  Der Fall „`.harness/` ausgenommen“ teilt sich die Prüfung mit `docs/`
  (Mutation 1 entfernt nur `docs`). Das Zeitlimit (`PIN_COMPARE_TIMEOUT`) hat
  seit der Fixrunde einen Tabellenfall (siehe dort).

### Implementer-Beleg Fixrunde (Review `review-slice-pin-stale-alle-digest-pins`, F-1 bis F-4)

- **F-1 Zeitlimit:** `run-pin-stale-all-tests.sh` trägt einen Fall mit
  `PIN_COMPARE_TIMEOUT=1` und einem Stub, der 6 s schläft (`exec sleep`):
  erwartet Exit 2, `UNBESTIMMT … Registry nicht erreichbar`, weder `OK` noch
  `DRIFT` für diese Referenz, die schnelle Nachbarin bleibt `OK`, Gesamtdauer
  unter 5 s; dazu ein Gegenfall mit demselben Limit ohne Verzögerung (Exit 0).
  `make test-pin-stale-all` am Original: `run-pin-stale-all-tests: alle 47
  Prüfungen bestanden` (Dauer 1,6 s).
- **Mutation (Instanz: Tabellentest, Prüfling per `PROG=<Kopie im Scratchpad>`,
  eine Stelle):** in `lib-pin-compare.sh` `limit=(timeout "$PIN_COMPARE_TIMEOUT")`
  auf `limit=()`. Gesehene Farbe: rot, Exit 1, fünf Prüfungen (Exit 1 statt 2,
  Zeile `DRIFT` statt `UNBESTIMMT`, „kein DRIFT“, Zusammenfassung, Lauf dauerte
  6 s). Der Stub antwortet ohne Limit mit leerem Digest, deshalb `DRIFT`.
- **F-2:** Grenze 3 im Vertrag benennt, dass der Treffer am Port beginnt; ein
  Tabellenfall (Stub, `localhost:5000/…`) belegt Treffer `5000/…` und Aufruf
  ohne Host. Die Antwort einer realen Registry darauf ist *hergeleitet*.
- **F-3:** [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 4 nennt „Registry nicht erreichbar“, nicht das
  Zeitlimit; dass ein Zeitablauf darunter fällt (`UNBESTIMMT`, nicht `DRIFT`),
  ist *hergeleitet* aus Festlegung 4 und durch den Tabellenfall aus F-1 am
  Sensor belegt, nicht in der ADR festgeschrieben.
- **F-4 / Grenze und Laufzeit:** Sensor-Vertrag Grenze 5 (Abruflimit 429 ist
  `UNBESTIMMT`, Exit 2, färbt den Nachtschritt) und Grenze 6 (Obergrenze je
  Referenz Limit mal Zahl, 15 × 60 s = 15 min *abgeleitet*, nicht unter dem
  Job-Limit). Gemessen am 2026-10-03: `make pin-stale-all` Exit 2, `real
  0m17,491s`, Zeile `pin-stale-all: 15 Referenzen — 6 OK, 0 DRIFT, 9
  UNBESTIMMT`, Ursache `429 Too Many Requests` (Docker Hub, an diesem Host).
  Der Lauf zum Plan-Beleg (28,845 s, 15 OK) war vorher; das Limit ist
  zeitabhängig.
- **Suchlauf (§3.13) — Gefundenes:** `make suchlauf-nachmessen PLAN=…` am Stand
  der Lieferung: `suchlauf-nachmessen: 6 Zeilen stimmen`, die `diff`-Zeilen
  (0, 0, 5) stimmen mit dem Soll. Nachgezogen, weil sie den alten Umfang
  nannten: `.github/workflows/upstream-drift.yml` (Kopf-Kommentar, Job-`name:`),
  `harness/README.md` (Zeile zu `upstream-drift.yml`),
  `docs/maintainer/releasing.md` (§5, Version 1.14). **Nichtgefunden:** kein
  Treffer in `spec/`, in `docs/user/` und im Quelltext außerhalb der genannten
  Träger; kein weiterer Aufrufer von `pin-stale.sh` außer den fünf im
  `Makefile` und `pin-stale-dcheck.sh`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), **und**
[`slice-pin-digests-aktualisieren-2026-10-b`](../done/slice-pin-digests-aktualisieren-2026-10-b.md)
liegt in `done/` **und** der Implementer hat die Startmessung aus §3 am
aktuellen Arbeitsstand gefahren: `make image-stale` und alle `make
pin-stale-*` enden mit Exit 0. (Ein frischer Drift seit der Hebung ist kein
Hindernis des Starts, sondern der erste Commit, §3 Reihenfolge.)

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Wiederverwendung
  verlangt eine Verhaltensänderung an `pin-stale.sh`, `image-stale.sh` oder
  den Aufrufern (Ausgabe oder Exit-Code einer Achse ändert sich) — der Refactor
  wird dann ein eigener Slice vor dem Sensor.
- `in-progress` → `open` (blockiert — Carveout?): der Vorgänger-Slice liegt
  noch nicht in `done/` (Abhängigkeit), oder ein Registry-/Netzausfall hindert
  die reale Messung am Arbeitsstand; kein Carveout, weil kein Gate rot ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make test-pin-stale-all` endet mit Exit 0, `make pin-stale-all` endet
am Stand der Closure mit Exit 0, der reale `workflow_dispatch`-Lauf von
`upstream-drift.yml` ist gelesen (zehn Schritte, `success`), und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; Kandidat: der Sensor P10 selbst und die Grenze „Form, nicht
Sinn“). Ein Gate, das am Stand der Closure rot ist, geht nur mit dokumentiertem
Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen** (Platzhalter `Ausgang: offen bis Closure`).

- **Der erste Lauf mit P10 ist rot** (ein Pin driftete zwischen der Hebung im
  Vorgänger-Slice und dem Start; oder eine Kopie blieb auf dem alten Digest):
  der Nachtlauf bekäme einen weiteren roten Lauf
  ([`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md),
  zwei Evidenz-Dateien am Planungsstand gezählt). Gegenmaßnahme ist die
  Startmessung (§4) und der Hebe-Commit vor dem Sensor (§3). —
  **Ausgang:** entfallen für den ersten Runner-Lauf — er ist gelesen (Nachtrag in §7:
  15 OK, 0 UNBESTIMMT, kein 429); die 429-Möglichkeit bleibt benannte Grenze 5
  des Sensor-Vertrags. Am Arbeitsstand ist `make pin-stale-all` Exit 0
  (15 OK, gemessen im Verifikations-Report §5). Dazu die Verifier-Feststellung V-1: neun der 15 Referenzen liegen
  bei Docker Hub, das anonyme Abruflimit (HTTP 429) färbte am Messhost den
  Sensor einmal auf `9 UNBESTIMMT`, Exit 2, ohne dass ein Pin driftete; ein
  roter erster Runner-Lauf kann also am Limit liegen, nicht an Drift. Der
  Re-Evaluierungs-Trigger „vier rote Nachtläufe“ von
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  zählt Läufe, nicht Ursachen — ein 429-Rot ist beim Zählen auszunehmen und
  zu benennen. Anker:
  [`harness/sensors/pin-stale-all.md`](../../../../harness/sensors/pin-stale-all.md)
  Grenze 5 (429 ist `UNBESTIMMT`) und
  [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md);
  der Hauptlauf liest den ersten Lauf und trägt Ursache und Ausgang nach.
- **Die Muster-Form fängt eine Referenz nicht oder fängt zu viel**: ein Pin, den
  das Muster nicht trifft (Digest in einer zur Laufzeit zusammengesetzten
  Variablen), wird nicht gesehen (akzeptiertes Negativ der ADR); eine
  Referenz-ähnliche Zeichenkette in einer getrackten Datei, die kein Pin ist
  (Fixture, Beispiel in einem Skript-Kommentar), wird als Pin gelesen und endet
  für immer als `DRIFT` oder `UNBESTIMMT`. Der eigene Tabellentest und der
  Sensor-Vertrag tragen deshalb keine vollständige Referenz als Literal (§2
  Liefer-Punkt 1). — **Ausgang:** entfallen am Stand: die Mengen-Messung
  (15 Referenzen, 20 Dateien) nennt jede Referenz als echten Pin, und die
  Gegenprobe des Verifiers fand keine Referenz der Form `<image>@sha256:`
  außerhalb des Musters (einzige Digest-Quelle außerhalb: `DCHECK_DIGEST`,
  eigene Achse P7). Das akzeptierte Negativ der ADR (zur Laufzeit
  zusammengesetzte Digests) bleibt als Grenze 3 im Sensor-Vertrag.
- **Die Wiederverwendung ändert das Verhalten von `pin-stale.sh`** (eine
  Zeichenfolge der Meldung, ein Exit-Code), und `make pin-stale-dcheck` oder der
  Workflow-Parser liest sie anders: `tools/harness/pin-stale-dcheck.sh` ruft
  `pin-stale.sh`. — **Ausgang:** entfallen: Ausgabe und Exit je Fall
  byte-gleich (Implementer: fünf Fehlerpfade; Verifier: neun Stub-Fälle, alle
  `SAME`; reale Läufe der sechs Ziele Exit 0). Ein Vorher/Nachher an der echten
  Registry fuhr niemand — die Gleichheit trägt der Stub-Vergleich.
- **Mutationsbeleg zu schmal**: drei Mutationen an drei Stellen am Tabellentest
  belegen nicht, dass alle Zweige tragen; die Verallgemeinerung steht als
  *hergeleitet*
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12). — **Ausgang:** entfallen: der
  Beleg nennt Stellen, Instanz und Farbe; sechs Mutationen an sechs Stellen
  (Implementer drei plus Zeitlimit, Verifier alle fünf nachgefahren plus
  Deduplizierung und Exit bei `DRIFT`), alle rot. Der Rest bleibt *hergeleitet*
  (`.harness/`-Ausschluss einzeln, Lesefehler von `git grep`, `UNBESTIMMT` ohne
  `DRIFT`); die Instanz ist ein Stub-`docker`.
- **Workflow-Änderung ist lokal unbeweisbar** ([`AGENTS.md`](../../../../AGENTS.md)
  §3.10): der zehnte Schritt, der Schritt-Name und die Make-Ziel-Namen müssen auf
  dem Runner stimmen; ein lokal grünes `make pin-stale-all` zeigt das nicht. —
  **Ausgang:** entfallen — der Push ist erfolgt und der
  `workflow_dispatch`-Lauf von `upstream-drift.yml` (zehn Schritte, `success`)
  ist gelesen (Nachtrag in §7); Liefer-Punkt 3 der DoD ist abgehakt.
  Träger:
  [`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md).
- **Laufzeit des Sensors**: je Referenz ein Registry-Aufruf; der Job hat
  `timeout-minutes: 15` ([`upstream-drift.yml`](../../../../.github/workflows/upstream-drift.yml)),
  die übrigen neun Schritte laufen davor. Die Laufzeit des Sensors ist nicht
  gemessen; ein Netzhänger je Referenz ohne Zeitlimit des Aufrufs wäre ein
  Dauerhänger bis zum Job-Limit. — **Ausgang:** offen bis Closure (zu belegen: die
  gemessene Laufzeit des realen Laufs, gedruckt in §7). —
  **Ausgang:** entfallen: gemessen `real 0m28,845s` (Implementer) und
  `real 0m34,461s` (Verifier) bei 15 Referenzen, `real 0m17,161s` bis
  `0m17,491s` bei 9 Referenzen mit 429; das Zeitlimit je Aufruf
  (`PIN_COMPARE_TIMEOUT`, Tabellenfall) ist belegt. Die Obergrenze 15 × 60 s
  liegt nicht unter dem Job-Limit (Vertrag Grenze 6, *abgeleitet*) — benannt,
  nicht behoben.
- **Plattformabhängigkeit des Vergleichs**: der Index-Digest ist plattformunabhängig
  (Festlegung 2), ein Einzelplattform-Pin wird `DRIFT` — das ist gewollt, aber ein
  Pin, den eine Kopie bewusst als Einzelplattform trägt, wäre ein Dauer-Rot ohne
  Hebe-Weg. — **Ausgang:** entfallen am Stand: der PostgreSQL-17-Pin
  (`.github/workflows/e2e.yml`) ist ein Index-Digest, alle 15 Referenzen `OK`
  (Verifikations-Report §5); tritt ein bewusster Einzelplattform-Pin später auf,
  ist das eine Entscheidung des Eigentümers, nicht dieses Slice.

## 7. Closure-Notiz

Wird mit der Closure gefüllt (Inhalt, dann `git mv`, dann Häkchen der
Paarungs-Zeile — [`AGENTS.md`](../../../../AGENTS.md) §3.3). Hier stehen: die
Mengen-Messung der Referenzen (Zahl der verschiedenen Referenzen und Dateien,
Parent-Kennung), die gedruckten Zeilen von `make pin-stale-race|-pgtest|-dmigrate|-acheck`
und `make image-stale` vor und nach dem Refactor, der Mutationsbeleg (je
Mutation Stelle, Instanz, Farbe), die Zeilen des realen `make pin-stale-all`,
die gemessene Laufzeit und der gelesene `workflow_dispatch`-Lauf.

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Belege (Mengen-Messung, Refactor-Vergleich, Mutationsbeleg, Zeilen des realen
Laufs, Laufzeit) stehen in §3 „Implementer-Beleg“ und in der Fixrunde; die
unabhängige Nachmessung steht im
[Verifikations-Report](../../../reviews/verifikation-slice-pin-stale-alle-digest-pins.md)
(bestanden mit einer Bedingung: dem Runner-Lauf, DoD Liefer-Punkt 3).

- **Was hat funktioniert:** Die Wiederverwendung statt Duplikat
  (`lib-pin-compare.sh`) ließ `pin-stale.sh` byte-gleich in Ausgabe und Exit;
  der Stub-`docker` machte den Tabellentest netzlos und die Mutationsläufe an
  Kopien reproduzierbar (der Verifier fuhr sechs Mutationen, alle rot). Der
  Review fand das fehlende Zeitlimit als MEDIUM, bevor es ein Lauf tat.
- **Was ging anders als geplant:** Die Make-Ziele liegen im `Makefile`, nicht in
  einem `*.mk`; `PIN_COMPARE_TIMEOUT` kam als Antwort auf das Laufzeit-Risiko
  hinzu; das Abruflimit der Registry (429) färbte am Messhost einen Lauf auf
  `9 UNBESTIMMT`, ein Fall, den der Plan als Risiko nicht trug (Vertrag
  Grenze 5, §6 erstes Risiko). Der Runner-Lauf steht aus.
- **Steering-Loop-Eintrag (neuer Sensor, Spec-Lücke benannt):** (1) Neuer Sensor
  P10 `make pin-stale-all`: das Pin-Inventar ist eine quantifizierte Regel („jede
  Digest-Referenz des Baums“), keine Liste; ein neuer Pin braucht keinen Eintrag
  mehr, um gelesen zu werden. liegt in: `harness/sensors/pin-stale-all.md §Vertrag`. (2) Benannte Grenze, nicht Regel: ein Nachtlauf-Sensor, der viele
  Aufrufe an ein anonymes Abruflimit richtet, kann ohne Drift rot werden, und
  seine Ausgabe nennt die Ursache nicht (`Registry nicht erreichbar`, kein
  Statuscode). Ein Zähler „rote Läufe in Folge“ muss deshalb Ursachen lesen, nicht
  nur Farben; das trägt der Vertrag (Grenze 5) und §6, kein Sensor.
- **Validator:** entfällt: der Slice liefert ein advisory Betreiber- und
  Wartungswerkzeug ohne End-Nutzer-Wert (kein Produktverhalten, kein
  SDK-/Image-Pin berührt).
- **Release-Folge:** keine.
- **Beobachtungs-Register (`../observations/`):**
  [`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md):
  Zustand bleibt **weiter offen**, Träger Sensor gebaut; er wird erst nach dem
  gelesenen grünen Nachtlauf fortgeschrieben (Nachtrag im Hauptlauf).
  [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md):
  nicht gezählt — V-1 (429-Rot ohne Drift) zählt erst nach dem Runner-Lauf und
  nur, wenn er eintritt. [`BEO-PGC/gate-prueft-existenz-nicht-passung`](../observations/BEO-PGC/gate-prueft-existenz-nicht-passung/observation.md):
  **kein** Querbezug, kein Zähler — die Klasse ist die Zuordnung einer
  Kennung zur Ursache der Stelle (Sinn einer vorhandenen Kennung); P10 prüft
  die Aktualität eines Digests gegen die Registry, und die Aussage „Form, nicht
  Sinn“ steht bereits als Grenze 1 im Vertrag. Ein Zähler aus einer
  Verwandtschaft des Wortlauts verwässerte die Zählung (der Zähler folgt aus
  `evidence/`-Dateien, und es gibt keine passende). V-2 (Ausgabe unterscheidet
  429 nicht von „nicht erreichbar“, LOW): **keine** neue Beobachtung — Vertrag
  Grenze 5 trägt es bereits; ein Eintrag wäre Doppelung, und eine
  Beobachtung braucht ein reales Vorkommnis auf dem Runner. Tritt ein
  429-Rot im Nachtlauf auf und kostet das Lesen der Ursache Zeit, legt der
  Hauptlauf dann ein Register-Verzeichnis mit diesem Lauf als Evidenz an; kein
  Folge-Slice. V-3 (INFO, Stub-Instanz) steht als *hergeleitet* im Beleg.
- **Folge-Slices:** keiner.
- **Risiken aus §6:** je ein Ausgang, an den Zeilen von §6 ablesbar: entfallen —
  Muster zu eng/zu weit, Verhaltensänderung `pin-stale.sh`, Mutationsbeleg zu
  schmal, Laufzeit, Einzelplattform-Pin; weiter offen — erster Lauf rot (mit
  V-1, Anker: Vertrag Grenze 5 und die Alarmmüdigkeits-Beobachtung), Workflow
  lokal unbeweisbar (Anker: `BEO-PGC/github-actions-unverifizierbar-lokal`).
- **Drei Paarungen:** Anker — `liegt in: harness/sensors/pin-stale-all.md` ist
  eine existierende Datei mit §Vertrag und §Grenze; Folge-Slice — keiner
  geplant, nichts offen ohne Adresse (die beiden offenen Risiken tragen
  Register-Anker); Register — alle genannten Kennungen existieren als
  Verzeichnis mit `evidence/` (gemessen: `ls` am Endstand).
- **Nachtrag (erledigt, 2026-10-03):** Der Stand `eb4155e5` ist gepusht
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Gelesen wurde der reale Lauf auf
  dem Runner:
  - `ci`, `e2e` (beide Legs) und `examples` desselben Stands: `success`.
  - `upstream-drift.yml`, Lauf `37144584033` (`workflow_dispatch`, Stand
    `eb4155e5`): Gesamtlauf `success`, alle zehn Schritte P1/P2 bis P10
    `success`. Das Log von P10 druckt
    `pin-stale-all: 15 Referenzen — 15 OK, 0 DRIFT, 0 UNBESTIMMT`; die
    Zeitstempel des Logs (`Run make pin-stale-all` 18:33:04, Schlusszeile
    18:33:25) ergeben rund 21 s. Kein HTTP 429 im Log.

  Liefer-Punkt 3 und die Register-Zeile sind damit abgehakt, beide offenen
  Risiken aus §6 für den ersten Runner-Lauf entfallen. Der Zustand von
  [`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
  ist im selben Zug fortgeschrieben. Ein Lauf per `workflow_dispatch` ist ein
  Lauf; der erste planmäßige Nachtlauf (`schedule`) mit P10 ist nicht gelesen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist allein die Default-Sub-Area
`*`/`PGC` (Harness-Werkzeuge unter `tools/harness/`, Make-Ziele, Workflow,
Sensor-Vertrag; kein Domänen- oder Adapter-Code). Eine feinere Aufteilung trägt
der Slice nicht.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`, Namen gelesen, Evidenz-Dateien am
Planungsstand gezählt). Treffer:
[`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
(Zähler 1) — der Gegenstand dieses Slice, Träger
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md).
[`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
(Zähler 2): §6. [`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md)
(Zähler 8, eigener Träger in `AGENTS.md` §3.10): §6. Kein Treffer für die
übrigen Einträge.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt `*` als
Greenfield).
