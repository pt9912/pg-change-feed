# Slice pin-digests-aktualisieren-2026-10: Gedriftete Upstream-Pins heben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht (Modul 6).

**Bezug:** [`LH-QA-POR-001`](../../../../spec/lastenheft.md),
[`LH-QA-POR-003`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7
(Upstream-Pin-Freshness; „Digest-Hebung bleibt bewusster Commit“),
[`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md),
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md).

**Berührte Spec-Stellen:** — (gemessen: `git grep` nach `sha256:` in `spec/`,
`docs/user/`, `README.md`, `AGENTS.md` findet keinen der Digests; `SPEC-012`
nennt nur die Major-Versionen 17 und 18, keinen Digest). Eine Spec-Berührung
entsteht nur, wenn die Messung unter §3 Schritt 1 eine neue PostgreSQL-Major-
oder eine Tag-Änderung zeigt; dann wird sie benannt, nicht in diesem Slice
entschieden.

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Anlass (übernommen, nicht gemessen):** der nächtliche Lauf
[upstream-drift 37112891025](https://github.com/pt9912/pg-change-feed/actions/runs/37112891025)
(advisory, Messzeitpunkt 2026-10-03, Repo-Stand `7aa31562`) meldete Drift an
`golang:1.27-alpine`, `TOOLCHAIN_RACE_IMAGE`, `PG_TEST_IMAGE`, `D_MIGRATE_IMAGE`,
`A_CHECK_IMAGE` und die Tag-Frische von `DCHECK_IMAGE`. Der Lauf nannte `P9` in
seiner Fehlerliste. **Nachgemessen am Arbeitsstand `e5820a03`** (2026-10-03,
`make image-stale` und `make pin-stale-race`/`-pgtest`/`-dmigrate`/`-acheck`/
`-dcheck`/`-baseline`/`-actions`): dieselben sechs Befunde als `DRIFT`
(`image-stale`, die vier Digest-Achsen und `pin-stale-dcheck` enden über
`make` mit Exit 2), dazu `OK` für den Digest von `DCHECK_DIGEST`
(`v0.77.0`; nur die Tag-Frische liegt hinter `v0.79.0`), `P2` `OK`, `P8` `OK`
(Baseline `v6.13.0`) und
`P9` `OK` (alle vier Actions: Tag-Mutation SHA unverändert, Tag-Frische gleich
dem neuesten Release; `astral-sh/setup-uv` steht seit dem Bump `06ffb834` auf
`v10.2.0`). Dass `P9` im Lauf rot war, am Arbeitsstand aber grün ist, ist eine
**Erwartung** (der Bump `06ffb834` liegt nach dem Repo-Stand des Laufs), zu
belegen durch den Blick in das Log des Laufs im Schritt 1 des Plans — die
vollständigen Digests stehen in den Messzeilen der Befunde und werden dort
neu gemessen, nicht aus der Anfrage übernommen.

**Ziel:** Jeder gemessene Drift der Upstream-Pins P1 und P3 bis P7 ist als
bewusster Digest-Commit gehoben, jeder Träger des alten Pins nachgezogen, und
die Test- und Gate-Läufe belegen, dass die neuen Bau- und Prüfwerkzeuge das
Repo gleich bewerten wie vorher.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Befunde, die ein gehobenes Prüfwerkzeug erstmals meldet** (`a-check`,
  `d-check`, `d-migrate`) — eine Lockerung wäre eine Schwellen-Senkung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6) und braucht eine ADR; die
  Bereinigung ist eigene Arbeit. Ein solcher Befund hält den betroffenen Pin
  auf dem Parent (siehe §6), der Slice hebt dann die übrigen Pins.
- **Eine Änderung der PostgreSQL-Hauptversionen oder des Image-Tags**
  (`17`/`18`, `-alpine`) — das wäre eine Änderung von
  [`SPEC-012`](../../../../spec/pflichtenheft.md), kein Digest-Commit.
- **Ein Sensor für den PostgreSQL-17-Pin in `.github/workflows/e2e.yml`** —
  das Pin-Inventar von [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
  führt ihn nicht, und `ADR-0051` ist `Accepted`; ein Sensor für ihn ist eine
  Folge-ADR mit eigenem Slice (§6, offene Frage). Der Slice **misst** den Pin
  einmal von Hand und hebt ihn mit.
- **Pins ohne Eintrag im Inventar** (die gepinnten Basis-Images der SDK-,
  Beispiel- und Schema-Dockerfiles außer `examples/Dockerfile`) — kein Sensor
  liest sie, die Drift ist unbekannt; die Messung gehört in den Folge-Befund
  des §7, nicht in diese DoD.
- **Die in `Accepted` ADRs und in `docs/reviews/`, `done/` und `observations/`
  stehenden alten Digests** — sie halten den Stand ihrer Messung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); eine Zitat-Korrektur nach
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  trifft sie nicht, weil der Referent (der Pin zu seiner Zeit) unverändert ist.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Messung (Liefer-Punkt 1).** Alle Achsen P1 bis P9 sind am Arbeitsstand
      des Implementers neu gemessen (`make image-stale`, `make pin-stale-race`,
      `make pin-stale-pgtest`, `make pin-stale-dmigrate`,
      `make pin-stale-acheck`, `make pin-stale-dcheck`,
      `make pin-stale-baseline`, `make pin-stale-actions`; der Pin
      `postgres:17-alpine` in `.github/workflows/e2e.yml` von Hand gegen die
      Registry); die gedruckten Zeilen mit den vollen Digests und der Parent-
      Kennung stehen in §7. Zeigt `P9` oder `P8` Drift, steht der Befund samt
      Ausgang dort; eine Achse ohne Drift wird nicht angefasst.
- [x] **Hebung und Nachzug (Liefer-Punkt 2).** Je gedriftete Achse ein eigener
      Digest-Commit (Commit-Betreff nennt
      [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md), keine
      `SPEC-`/`ARC-`-Kennung): `golang:1.27-alpine` (`Dockerfile`, `Makefile`
      `TOOLCHAIN_IMAGE`, `examples/Dockerfile` und die in §3 genannten
      Skript-Vorgaben), `TOOLCHAIN_RACE_IMAGE`, `PG_TEST_IMAGE` (`Makefile`,
      `compose.yaml`, `examples/compose.yaml`, `.github/workflows/e2e.yml` und
      Skript-Vorgaben), `D_MIGRATE_IMAGE` (`Makefile`, `tools/bench-lib.sh`),
      `A_CHECK_IMAGE` (`a-check.mk`), `DCHECK_IMAGE`/`DCHECK_DIGEST`
      (`d-check.mk`, beide Zeilen in einem Commit). Der Suchlauf aus §3
      (`diff`-Stand) zeigt keinen alten Digest mehr in einem lebenden Träger;
      die verbleibenden Treffer sind die in §1 ausgeschlossenen Records.
- [x] **Belege (Liefer-Punkt 3).** Je gehobene Achse läuft der engste
      betroffene Lauf aus §3 und die Ausgabe wird gegen den Parent verglichen
      (Befunde, Meldungsform); `make image` endet erfolgreich, `make
      test-store`, `make test-replication` an PostgreSQL 18 und mit
      `PG_TEST_IMAGE` auf den gehobenen PostgreSQL-17-Digest, und
      `make test-integration` enden grün. Ein Lauf, der am Parent grün und am
      Diff rot ist, ist ein Befund (§6), kein Anlass, den Lauf zu lockern.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [x] `make gates` grün (Exit direkt ausgewertet, Exit 0 am Endstand, Zeilen
      `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`,
      `d-check: 1620 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)`;
      `make test` Exit 0 nach P3, `make docs-check` Exit 0;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make test` grün, `make
      docs-check` grün (Exit direkt); `make mod-download` nur, falls `go.mod`
      berührt wird — erwartet nicht.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      ([`review-slice-pin-digests-aktualisieren-2026-10`](../../../reviews/review-slice-pin-digests-aktualisieren-2026-10.md),
      `.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: kein öffentlicher Vertrag berührt (gemessen am Diff: kein
      Digest-Treffer in `spec/`, `docs/user/`, `README.md`, `AGENTS.md`; erwartet — gemessen
      unter §Berührte Spec-Stellen: kein Digest in `spec/`, `docs/user/`);
      die Träger-Nachzüge stehen in der Liefer-Punkt-2-Zeile.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` (deps-Stage), `Makefile` (`TOOLCHAIN_IMAGE`), `examples/Dockerfile`, `tools/harness/ci-matrix-abdeckung.sh`, `tools/harness/lib-github-api.sh`, `tools/harness/run-integration-tests.sh`, `tools/harness/run-notify-tests.sh`, `tools/harness/run-replication-tests.sh`, `tools/harness/run-store-tests.sh` | update | P1: neuer `golang:1.27-alpine`-Digest; die Skript-Vorgaben tragen denselben Digest als Default (`git grep`, siehe Suchlauf) — Liefer-Punkt 2 |
| `Makefile` (`TOOLCHAIN_RACE_IMAGE`) | update | P3: neuer `golang:1.27`-Digest (Debian-Variante für `-race`) |
| `Makefile` (`PG_TEST_IMAGE`), `compose.yaml`, `examples/compose.yaml`, `.github/workflows/e2e.yml` (PG-18-Leg), `tools/bench-lib.sh`, `tools/harness/run-replication-tests.sh`, `tools/harness/run-store-tests.sh`, `tools/harness/run-schema-rollout-guard-test.sh` | update | P4: neuer `postgres:18-alpine`-Digest; `e2e.yml` hält den Pin als Matrix-Eintrag, der Kommentar dort nennt „derselbe Digest-Pin wie Makefile“ |
| `.github/workflows/e2e.yml` (PG-17-Leg, Kommentar „ermittelt 2026-09-14“) | update, falls die Handmessung Drift zeigt | PG-17-Digest ist nicht im Inventar; Messung von Hand, Kommentar-Datum und Digest gemeinsam |
| `Makefile` (`D_MIGRATE_IMAGE`), `tools/bench-lib.sh` | update | P5: neuer `d-migrate:latest`-Digest |
| `a-check.mk` (`A_CHECK_IMAGE`) | update | P6: neuer `a-check:latest`-Digest |
| `d-check.mk` (`DCHECK_IMAGE`, `DCHECK_DIGEST`) | update | P7: Tag `v0.79.0` und der dazu gemessene Digest, beide Zeilen in einem Commit (der Digest sticht den Tag, ein Tag ohne neuen Digest wäre wirkungslos) |
| `.d-check.yml` | update, nur falls der Lauf mit der neuen Version eine Schlüssel-/Modul-Änderung verlangt | eine Konfigurationsänderung ist ein Befund-Anlass (§6), keine stille Mitänderung |
| `docs/plan/planning/in-progress/slice-pin-digests-aktualisieren-2026-10.md` | update | §7 Messzeilen, Vergleich gegen den Parent |

- **Reihenfolge:** (1) messen und die Zeilen in §7 festhalten; (2) je Achse
  einen Digest hochziehen, den engsten Lauf fahren (P1/P3: `make test`, `make
  image`; P4: `make test-store`, `make test-replication`; P5: `make
  schema-validate`, `make test-store`; P6: `make a-check`; P7: `make
  docs-check`), dann committen; (3) am Ende `make gates`, `make image`,
  `make test-integration`. Eine Achse nach der anderen, damit ein roter Lauf
  einer Achse zuzuordnen ist.
- **Digest-Messung:** der neue Digest wird mit dem Messwerkzeug des Pins
  gelesen (`tools/harness/pin-stale.sh`, `tools/harness/pin-stale-dcheck.sh`
  drucken ihn in der `aktuell`-Spalte); der Implementer trägt den Wert aus
  der gedruckten Zeile ein, nicht aus diesem Plan.
- **PG-17-Messung:** der Implementer nennt Befehl und Zeile (der Kommentar in
  `e2e.yml` nennt `docker manifest inspect postgres:17-alpine`, amd64) in §7.
- **d-migrate:** die vier Nacharbeits-Schritte `tools/schema/nacharbeit-*.sql`
  sind eine Antwort auf das Verhalten eines früheren d-migrate-Stands
  ([`BEO-PGC/d-migrate-nacharbeit`](../observations/BEO-PGC/d-migrate-nacharbeit/observation.md));
  ein neuer Stand kann sie überflüssig oder falsch machen. Der Lauf von
  `make schema-rollout` gegen eine Wegwerf-Datenbank
  (`tools/harness/run-schema-rollout-guard-test.sh`, sechs Läufe) und
  `make test-store` zeigen es; ein Abweichen ist ein Befund (§6).
- **Suchlauf (§3.13, bewegte Eigenschaft: der Digest-Wert eines Pins).** Der
  Parent ist `e5820a030a2dea7e2deafc474fadbf973673cd6a`
  (`git rev-parse HEAD` am Planungsstand; nie `HEAD`). Suchraum: ganzer Baum
  ohne `docs/reviews`, `done/`, `observations/` und `.harness/baseline`
  (Records und vendored Baseline). Gemessen am Parent (Soll = Trefferzeilen);
  der Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und
  Nichtgefundenes ein:

```suchlauf
e5820a030a2dea7e2deafc474fadbf973673cd6a 31 -n -E 'cf6fca66|b475798f|63bdc97d|862dfb04|34d3dfb5' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 14 -n -E 'cf6fca66' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 2 -n -E 'b475798f' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 10 -n -E '63bdc97d' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 3 -n -E '862dfb04' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 2 -n -E '34d3dfb5' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 1 -n -E 'd-check:v0\.77\.0' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
e5820a030a2dea7e2deafc474fadbf973673cd6a 1 -n -E '7456ef82' -- . ':!docs/reviews' ':!docs/plan/planning/observations' ':!.harness/baseline' ':!docs/plan/planning/done'
diff 10 -n -E 'cf6fca66|b475798f|63bdc97d|862dfb04|34d3dfb5' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 5 -n -E 'cf6fca66' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 1 -n -E 'b475798f' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 2 -n -E '63bdc97d' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 1 -n -E '862dfb04' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 1 -n -E '34d3dfb5' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 0 -n -E 'd-check:v0\.77\.0|3f84502b' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness/baseline'
diff 0 -n -E '7456ef82' -- . ':!docs/reviews' ':!docs/plan/planning/observations' ':!.harness/baseline' ':!docs/plan/planning/done'
```

  **Erwartung am `diff`-Stand nach der Hebung** (hergeleitet aus der
  Aufteilung der 31 Parent-Treffer in lebende Träger und Records, noch nicht
  gemessen): die Sammelzeile 10 (je Muster: `cf6fca66` 5, `b475798f` 1,
  `63bdc97d` 2, `862dfb04` 1, `34d3dfb5` 1) — die Treffer sind
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (fünf Zeilen des
  Pin-Inventars),
  [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md),
  [`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md),
  [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  und `harness/sensors/db-adapter-coverage.md` (Zeile 300, Messkontext einer
  früheren Messung — der Implementer liest sie und entscheidet, ob sie einen
  Stand beschreibt, der weiter gilt; sonst Nachzug). `d-check:v0.77.0` 0.
  `7456ef82` 0 nur, wenn die Handmessung den PG-17-Pin als gedriftet zeigt,
  sonst bleibt 1.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), und der Implementer hat den Messlauf aus Liefer-Punkt 1 am
aktuellen Arbeitsstand gefahren — zeigt er keine Drift mehr, entfällt der
Slice (siehe unten).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zeigt
  mehr als die sieben genannten Achsen oder ein gehobener Pin verlangt
  Quelltext-Änderungen über den Pin hinaus (neue Befunde, geänderte
  Konfigurationsschlüssel) — der Slice hebt dann nur die unbeteiligten Achsen,
  der Rest wird ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): ein Registry- oder Netzausfall
  hindert die Messung oder das Ziehen eines Digests; kein Carveout, weil kein
  Gate rot ist.
- Entfällt der Gegenstand (der nächtliche Lauf zeigt vor dem Start nichts mehr),
  geht die Datei nach `done/` mit der Zeile `Gegenstand:` in §7.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, ein erneuter Lauf der Drift-Messung (`make image-stale`, die
`make pin-stale-*`-Ziele) druckt für jede gehobene Achse `OK`, und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor
oder benannte Spec-Lücke). Ein Gate, das am Stand der Closure rot ist, geht nur
mit dokumentiertem Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register).

- **Neue Befunde eines gehobenen Prüfwerkzeugs** (`a-check` P6, `d-check` P7
  über zwei Releases `v0.77.0` bis `v0.79.0`, `d-migrate` P5): `make a-check`
  und `make docs-check` können am Diff Befunde melden, die der Parent nicht
  hat; eine Meldungsform- oder Schlüsseländerung kann `.d-check.yml` brechen.
  Ein Befund ist Folge-Arbeit; die Gate-Schwelle bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6). Der Pin der Achse bleibt dann
  auf dem Parent. — **Ausgang:** entfallen: am Diff meldet weder `make a-check`
  (`gesamt: 0 Befund(e)`) noch `make docs-check` (0 Befund(e), v0.79.0) einen
  neuen Befund (Parent-Vergleich im Klon, Verifikationsbericht §5).
- **Tag-Frische von `DCHECK_IMAGE` ohne Digest-Messung:** `pin-stale-dcheck`
  meldet den Tag; der zugehörige Digest ist eine eigene Messung. Ein Tag-Wechsel
  ohne neuen `DCHECK_DIGEST` bliebe wirkungslos (der Digest sticht den Tag,
  `d-check.mk`). — **Ausgang:** entfallen: `make pin-stale-dcheck` druckt
  `OK … DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.79.0) == sha256:b4b8756b…` und
  `OK … Tag-Frische (v0.79.0)` (Verifikationsbericht §1, eigener Lauf).
- **Neue Patch-Version von PostgreSQL oder Go ändert ein Testverhalten**
  (`make test-store`, `make test-replication`, `make test-integration`, beide
  PostgreSQL-Versionen): ein roter Lauf am Diff, grün am Parent, ist ein
  Befund; das Zeitverhalten der Replikations-Tests (`wal_sender_timeout`-Fenster)
  ist am wahrscheinlichsten betroffen. — **Ausgang:** entfallen: `make test`,
  `make test-store`, `make test-replication` (PostgreSQL 18.6 und 17.11) und
  `make test-integration` endeten am Diff mit Exit 0, beim Implementer und im
  eigenen Lauf des Verifiers.
- **CI-Matrix `e2e.yml` PG 17/18:** die Digests dort sind Pins des Workflows;
  ein lokal grüner Lauf zeigt nicht, dass der Runner grün läuft. Der Slice
  berührt `.github/workflows/e2e.yml` (Pin-Wert, keine Struktur) — nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 ist die strukturelle Änderung
  Auslöser der Pflicht; eine reine Digest-Hebung ist es nicht, der reale
  Post-Push-Lauf von `e2e.yml` (beide Matrix-Legs) wird dennoch gelesen
  (`gh run list --workflow e2e.yml`), weil sich die Pins ändern, gegen die die
  Matrix läuft. Der Push ist nicht Teil dieses Auftrags; der Lauf steht aus.
  Die Closure berührt `.github/workflows/e2e.yml` ein zweites Mal (nur
  Kommentarzeilen zu `SPEC-012`, keine Struktur, kein Pin); der Post-Push-Lauf
  deckt beides.
  — **Ausgang:** weiter offen, bis der Post-Push-Lauf von `e2e.yml` (beide
  Legs) und `upstream-drift.yml` gelesen ist (`gh run list --workflow
  e2e.yml`, `gh run list --workflow upstream-drift.yml`;
  [`AGENTS.md`](../../../../AGENTS.md) §3.10); Träger der Nachverfolgung ist
  [`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md).
- **`image-stale`-Befund `P1` und `make image`:** der neue `golang:1.27-alpine`-Digest
  ändert das Toolchain-Image jedes Docker-only-Laufs
  ([`AGENTS.md`](../../../../AGENTS.md) §3.1); ein Rückfall auf den Parent-Digest
  in einem der 9 Träger ließe zwei Toolchains koexistieren. — **Ausgang:**
  entfallen: der `diff`-Suchlauf aus §3 trifft die
  Erwartung (`make suchlauf-nachmessen`: 16 Zeilen stimmen), und die
  Träger-Auszählung des Verifiers zeigt je Achse genau einen Wert.
- **Offene Frage — Sensor-Lücke PG 17:** der Pin `postgres:17-alpine` in
  `e2e.yml` fehlt im Pin-Inventar von
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) und damit in
  `upstream-drift.yml`; die Drift dort bliebe unsichtbar.
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) ist `Accepted`;
  die Aufnahme ist eine Folge-ADR (Architect). Der Slice benennt sie, entscheidet
  sie nicht. — **Ausgang:** weiter offen:
  [`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
  (die Handmessung zeigte Drift, die der Sensor nicht gemeldet hätte); die
  Folge-ADR ist eine Frage an den Architect.
- **Offene Frage — stale Zitat:** [`harness/README.md`](../../../../harness/README.md)
  nennt für `make image-stale` die Bindung `ADR-0039` („Update =
  bewusster Digest-Commit“); der [ADR-Index](../../adr/README.md) führt diese
  Kennung als `Superseded` und mit anderem Gegenstand (Paketstruktur). Die tragende ADR für
  die Digest-Hebung ist [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
  Entscheidung 7. Der Slice zieht das nicht nach (fremde, nicht berührte
  Zeile); gemeldet. — **Ausgang:** eingetreten und mit der Closure nachgezogen
  (Bindung `ADR-0051` Entscheidung 7 in `harness/README.md`, eigener Commit).
- **Keine neue ADR nötig** — Begründung: Entscheidung 7 von
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) führt die
  Hebung ausdrücklich als „bewusster Commit“; der Slice ändert weder Tag noch
  Schwelle. Das gilt, solange die Messung keine Tag- oder Hauptversions-
  Änderung zeigt. — **Ausgang:** entfallen: §7 nennt keine solche Änderung
  (`go1.27.1` am alten wie am neuen Digest, PostgreSQL 18.6 und 17.11).

## 7. Closure-Notiz

Wird mit der Closure gefüllt (Inhalt, dann `git mv`, dann Häkchen der
Paarungs-Zeile — [`AGENTS.md`](../../../../AGENTS.md) §3.3). Messzeilen der
Planung (gemessen am Arbeitsstand `e5820a03`, 2026-10-03, volle Digests aus den
gedruckten Zeilen in `harness/`-Ausgaben des Implementers einzutragen, nicht
hier abgekürzt zu übernehmen):

- **Drift-Messung P1–P9 (gemessen, 2026-10-03, Arbeitsstand `8aeec6f2`
  (Parent des ersten Hebungs-Commits), Befehle `make image-stale` und
  `make pin-stale-race|-pgtest|-dmigrate|-acheck|-dcheck|-baseline|-actions`,
  Exit je `make` direkt gelesen):**

  ```text
  DRIFT golang:1.27-alpine: gepinnt sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125, aktuell sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414   (P1, Exit 2)
  OK    gcr.io/distroless/static-debian12:nonroot == sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab   (P2)
  DRIFT TOOLCHAIN_RACE_IMAGE (golang:1.27): gepinnt sha256:b475798fb16158e6c38e8b5ca2d870fbeaa8b7fec0fc8ec64b3dc20966040635, aktuell sha256:e0174e51e81218523251d85d248a90d24c3d5e81543b4f07a5d66229397db190   (P3, Exit 2)
  DRIFT PG_TEST_IMAGE (postgres:18-alpine): gepinnt sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8, aktuell sha256:77f585114c32fbca283dc835b0596f4e52b51b4c6662d7810b2f4084f60a1873   (P4, Exit 2)
  DRIFT D_MIGRATE_IMAGE (ghcr.io/pt9912/d-migrate:latest): gepinnt sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89, aktuell sha256:af9d3eb323a6cfd13eb63f012823788e510f764509918933d2ff98c8e68e4212   (P5, Exit 2)
  DRIFT A_CHECK_IMAGE (ghcr.io/pt9912/a-check:latest): gepinnt sha256:34d3dfb50e44d99ea735186a35e1040589c4681dcfa2a51ed0f2aaea718cdd2d, aktuell sha256:e8208764b119c606c92f82722813386277a65b12812d23b6107ea7a14dc25da1   (P6, Exit 2)
  OK    DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.77.0) == sha256:3f84502b09af65246fff38b1c3893130050e50581943a0434da95bf68091e337
  DRIFT DCHECK_IMAGE Tag-Frische: gepinnt v0.77.0, neuester Release v0.79.0   (P7, Exit 2)
  OK    Kurs-Baseline v6.13.0 == neuester Release   (P8, Exit 0)
  OK    actions/checkout@v7.0.1, docker/setup-buildx-action@v4.4.1, docker/login-action@v4.6.0, astral-sh/setup-uv@v10.2.0: Tag-Mutation und Tag-Frische je OK   (P9, Exit 0)
  ```

  Keine Tag- und keine Hauptversionsänderung: `golang:1.27` meldet
  `go version go1.27.1` am alten wie am neuen Digest (gemessen mit
  `docker run --entrypoint go … version` an beiden neuen Digests), PostgreSQL
  bleibt `18.6` (gedruckte Zeile des Replikations-Tests) und `17.11`.
  P8/P9 ohne Drift, nicht angefasst. **P9 im Lauf `37112891025` (gemessen
  in der Closure, Log des Jobs über `gh api repos/pt9912/pg-change-feed/actions/jobs/111174082159/logs`,
  Repo-Stand des Laufs `7aa31562`):** der Schritt war `failure`, gedruckt
  `DRIFT       astral-sh/setup-uv@v10.1.0 Tag-Frische: gepinnt v10.1.0, neuester
  Release v10.2.0` und `make: *** [Makefile:88: pin-stale-actions] Error 1`; die
  drei anderen Actions druckten je zwei `OK`. Die Plan-Erwartung trifft: der
  Bump `06ffb834` auf `v10.2.0` liegt nach dem Stand des Laufs, am Arbeitsstand
  ist `P9` `OK`.
- **PG-17-Handmessung (gemessen, 2026-10-03):** `docker buildx imagetools
  inspect postgres:17-alpine --format '{{.Manifest.Digest}}'` druckte den
  Index-Digest `sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24`;
  der amd64-Eintrag (Form des bisherigen Pins, `docker manifest inspect
  postgres:17-alpine`) trägt `sha256:aa90e97ee862e558111d34cfb8b2c4bec768c2b039fb791341686928560263b3`
  (Annotation `17.11-alpine3.24`); der alte Pin
  `sha256:7456ef82e5f5bc43d997f4781bbd7c0d6389bff397564649a356e206ba473aee`
  ist ein Einzel-Manifest (`application/vnd.oci.image.manifest.v1+json`), also
  ebenfalls amd64-Form. Die Drift ist real und vom Sensor von `upstream-drift`
  nicht gemeldet worden (Inventar führt PG 17 nicht, siehe §6).
- **Hebung (sieben Commits, je Achse einer):** P1 `cc93ca3b`, P3 `2d871258`,
  P4 `f69ce37e`, PG-17-Pin `f3b904ec`, P5 `79e40cbf`, P6 `7b5865f0`, P7
  `168db825`. DCHECK_DIGEST `sha256:b4b8756b40d3dcd2670a3f83526cb5e5d727d1a850571f73be31edba248abb40`
  gemessen mit `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.79.0
  --format '{{.Manifest.Digest}}'`; das Image trägt das Label
  `org.opencontainers.image.version` `0.79.0`; `make pin-stale-dcheck` druckt
  danach `OK … DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.79.0) == sha256:b4b8756b…`
  und `OK … Tag-Frische (v0.79.0)`, Exit 0.
- **Läufe am Diff (Exit direkt):** `make test` nach P1 und nach P3 Exit 0
  (zweiter Lauf 51 `ok`-Zeilen, kein `FAIL`); `make image` Exit 0 (nach P1 und
  am Endstand); `make test-store` Exit 0 nach P4 und nach P5, gedruckt
  `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%`;
  `make test-replication` Exit 0 an PostgreSQL 18 (`PostgreSQL 18.6: Keepalive
  inmitten der Transaktion …`) und mit `PG_TEST_IMAGE=postgres:17-alpine@sha256:aa90e97e…`
  Exit 0 (`PostgreSQL 17.11: Keepalive inmitten der Transaktion …`);
  `make schema-validate` Exit 0 (`Validation passed: 0 warning(s)`) und
  `bash tools/harness/run-schema-rollout-guard-test.sh` Exit 0 (`alle Belege
  real erbracht … make-Exit 2/2 mit d-migrate-Exit 8`) mit dem neuen
  d-migrate-Digest, die Nacharbeits-Schritte also weiter tragend;
  `make a-check` Exit 0 mit `gesamt: 0 Befund(e)` am alten **und** am neuen
  Digest; `make docs-check` Exit 0 mit `d-check: 1620 Datei(en) geprüft, 0
  Befund(e)` am alten (v0.77.0) **und** am neuen (v0.79.0) Digest;
  `make test-integration` Exit 0 am Endstand (`E2E-Abdeckungstabelle
  unverändert`, `Lauf abgeschlossen`). Ein erster Lauf von
  `make test-integration` wurde vom 30-Minuten-Limit des
  Hintergrund-Wrappers abgebrochen (kein Testfehler, Abbruch in der
  Upgrade-Phase, Container abgeräumt) und mit längerer Frist wiederholt. Der
  Parent wurde nur dort gefahren, wo ein Diff-Rot zu vergleichen gewesen wäre
  (`a-check`, `docs-check`); die übrigen Läufe sind am Diff grün, ein
  Parent-Lauf entfällt mangels Rot.
- **`diff`-Suchlauf (§3):** die acht `diff`-Zeilen des Blocks treffen die
  Erwartung (10 = 5+1+2+1+1; `d-check:v0.77.0`/`3f84502b` 0; `7456ef82` 0).
  Gefunden (alle Records bzw. Messkontext, bewusst nicht geändert):
  `docs/plan/adr/0051-…` (fünf Zeilen Pin-Inventar), `0071-…` (zwei),
  `0098-…` (eine), `0128-…` (eine), `harness/sensors/db-adapter-coverage.md`
  Zeile 300 (Rot-/Grün-Beleg eines Laufs von `slice-081`; beschreibt den Stand
  dieses Laufs, kein geltender Pin — nicht nachgezogen). Nichtgefunden: kein
  alter Digest in `Makefile`, `*.mk`, `Dockerfile`, `examples/`, `compose*.yaml`,
  `.github/workflows/`, `tools/`. Grenze: der Suchlauf liest Digest-Präfixe, keine
  Tag-Texte; das Makefile-Kommentar `go1.27.1` stimmt mit der Messung überein.
- **Gemeldete Träger fremder Dateien:** `harness/README.md` nannte für
  `make image-stale` die Bindung `ADR-0039` (stale, siehe §6) — mit der
  Closure auf `ADR-0051` Entscheidung 7 nachgezogen (der ADR-Index führt
  `ADR-0039` als `Superseded`, Gegenstand Paketstruktur).
- **Folge-Befund Out-of-Scope-Pins (gelesen in der Verifikation, `git grep -E
  '^FROM .*@sha256|image: .*@sha256'` über `sdks/`, `examples/` ohne
  `examples/Dockerfile` und `examples/compose.yaml`, und `tools/schema/`;
  Drift nicht gemessen):** fünf gepinnte Basis-Images ohne Sensor —
  `mcr.microsoft.com/dotnet/sdk:10.0` (`sdks/csharp`, `examples/csharp`),
  `mcr.microsoft.com/dotnet/runtime:10.0` (`examples/csharp`),
  `eclipse-temurin:21-jdk` (`sdks/kotlin`, `examples/kotlin`),
  `eclipse-temurin:21-jre` (`examples/kotlin`) und `python:3.14-slim`
  (`sdks/python`). `tools/schema/Dockerfile` bezieht `D_MIGRATE_IMAGE` und
  `TOOLCHAIN_IMAGE` als Build-Argument; sie wandern mit P1 und P5 mit. Adresse:
  `BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`; eine Messung oder ein
  Sensor ist Sache einer Folge-ADR (Frage an den Architect).
- **Post-Push-Lauf `e2e.yml` (beide Legs) und `upstream-drift.yml`:** steht aus
  (kein Push). Der Slice berührt `.github/workflows/e2e.yml` (nur Pin-Werte und
  das Datum im Kommentar, keine Struktur); `AGENTS.md` §3.10 meldet: der
  Post-Push-Lauf ist offen. Die Closure ergänzt vier Kommentarzeilen in
  `e2e.yml` (Z. 41ff. und die beiden `(Major nach SPEC-012)`-Zeilen), die
  Digests und die Struktur bleiben.
- **Beobachtung:** Der PG-17-Pin driftete ohne Sensor (Messung oben) — angelegt
  als `BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`.

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

- **Was hat funktioniert:** je Achse ein Commit mit engstem Lauf davor
  machte jeden Lauf einer Achse zuordenbar; kein gehobenes Prüfwerkzeug
  meldete einen neuen Befund (`a-check` und `d-check` am alten und neuen
  Digest gleich, Parent-Vergleich im Klon des Verifiers); der unabhängige
  Nachlauf des Verifiers bestätigte die Handmessung des PG-17-Pins.
- **Was ging anders als geplant:** ein erster Lauf von `make test-integration`
  brach am 30-Minuten-Limit des Hintergrund-Wrappers ab (kein Testfehler) und
  lief mit längerer Frist grün. Der Plan-Satz zu `P9` blieb bis zur Closure
  Erwartung; das Log des Laufs `37112891025` bestätigt ihn (§7 oben). Die
  Nacharbeit-Schritte von `d-migrate` tragen auch am neuen Pin
  ([`BEO-PGC/d-migrate-nacharbeit`](../observations/BEO-PGC/d-migrate-nacharbeit/observation.md):
  kein neuer Beleg, keine neue Evidenz-Datei).
- **Steering-Loop-Eintrag:** geschärfte Beobachtung, noch keine verkörperte
  Regel: **Ein Pin ohne Zeile im Inventar hat keinen Leser und driftet
  unsichtbar; eine Messung von Hand ist eine Momentaufnahme, kein Sensor.**
  Der PG-17-Pin war gedriftet, während alle Sensoren grün oder nur an anderen
  Achsen rot waren. Benannte ADR-Lücke: die Aufnahme des Pins und die Form
  des Sensors (`pin-stale.sh` liest Makefile-Variablen, der Pin ist ein
  YAML-Matrix-Wert); dazu die Inventar-Form, die Digest-Präfixe als Werte nennt
  (sie veralten mit jeder Hebung; eine Folge-ADR verweist auf die Variable).
  Zu zitierende Zahl: die Dateizahl von `d-check` nicht als Kennzahl führen —
  sie ändert sich am HEAD mit jedem Report (1620 beim Implementer, 1621 beim
  Verifier). Frage an den Architect (nicht angelegt): Folge-ADR zu
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7.
- **Beobachtungs-Register (`../observations/`):** neu
  [`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
  (Zähler 1). Zweite Evidenz-Datei zu
  [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
  (Zähler 2): die Klasse „advisory, nicht blockierend, rot ohne Reaktion“ trifft
  den Nachtlauf `upstream-drift.yml`, dessen 14 Läufe seit 2026-09-20 alle
  `failure` endeten (gemessen mit `gh run list`). Kein Eintrag erreicht 3×
  über diesen Slice.
- **Folge-Slices:** keiner angelegt. Die Hebung der fünf Basis-Images und ein
  Sensor für den PG-17-Pin hängen an der Folge-ADR des Architects.
- **Risiken aus §6:** je Risiko ein Ausgang in §6 (fünf entfallen, eines
  eingetreten und nachgezogen, zwei weiter offen: CI-Matrix `e2e.yml` samt
  `upstream-drift.yml` bis zum Post-Push-Lauf, Sensor-Lücke PG 17 mit
  Register-Eintrag).
- **Validator (Modul 8):** entfällt, weil der Slice Pflegearbeit an
  Upstream-Pins ist und keinen End-Nutzer-Wert liefert (kein
  Nutzerdokument, kein Vertrag berührt, gemessen am Diff).
- **Drei Paarungen:** (a) Anker: kein Steering-Loop-Eintrag trägt das Feld
  `liegt in`, nichts zu prüfen; (b) Folge-Slice: keiner genannt; (c) Register:
  beide genannten `BEO-PGC/`-Verzeichnisse existieren, jedes mit nicht leerem
  `evidence/`.

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
`*`/`PGC` (Build-, Test- und Gate-Konfiguration; kein Domänen- oder
Adapter-Code). Eine feinere Aufteilung trägt der Slice nicht: er ändert
Pin-Werte, keine Verhaltens-Pfade.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`, Namen gelesen). Treffer:
[`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
(Zähler 1: eine Datei unter `evidence/`) — benennt die Alarmmüdigkeit advisory-
Läufe und trifft hier: der nächtliche `upstream-drift`-Lauf meldete Drift seit
einem unbekannten Zeitpunkt; der Slice benennt ihn als Anlass, hält die
Beobachtung aber nicht für erreicht (kein dritter Beleg).
[`BEO-PGC/d-migrate-nacharbeit`](../observations/BEO-PGC/d-migrate-nacharbeit/observation.md)
(Zähler 8) betrifft die Hebung von `D_MIGRATE_IMAGE` unmittelbar — siehe §3
(d-migrate) und §6; dass sie bereits ≥ 3 Belege trägt, ist nicht Auslöser
dieses Slice (sie hat ihren eigenen Träger). Kein Treffer für die übrigen
Achsen.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt `*` als
Greenfield).

## Nachtrag — Post-Push-Läufe (2026-10-03)

Der Stand `38c9ff93` ist gepusht. Gelesen wurde der reale Lauf auf dem Runner
([`AGENTS.md`](../../../../AGENTS.md) §3.10):

- `e2e.yml`, Lauf `37125336646`: `gh run view` nennt beide Matrix-Legs,
  `image + test-integration (PostgreSQL 18) success` und
  `image + test-integration (PostgreSQL 17) success`. `ci` und `examples`
  desselben Stands: `success`.
- `upstream-drift.yml`, Lauf `37126796088` (`workflow_dispatch`, Stand
  `38c9ff93`): `success`, also ohne Drift auf P1 bis P9.

Das CI-Matrix-Risiko aus §6 ist damit **entfallen**. Weiter offen bleibt allein
die Sensor-Lücke des PG-17-Pins
([`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)).
