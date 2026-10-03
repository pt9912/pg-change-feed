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

- [ ] **Messung (Liefer-Punkt 1).** Alle Achsen P1 bis P9 sind am Arbeitsstand
      des Implementers neu gemessen (`make image-stale`, `make pin-stale-race`,
      `make pin-stale-pgtest`, `make pin-stale-dmigrate`,
      `make pin-stale-acheck`, `make pin-stale-dcheck`,
      `make pin-stale-baseline`, `make pin-stale-actions`; der Pin
      `postgres:17-alpine` in `.github/workflows/e2e.yml` von Hand gegen die
      Registry); die gedruckten Zeilen mit den vollen Digests und der Parent-
      Kennung stehen in §7. Zeigt `P9` oder `P8` Drift, steht der Befund samt
      Ausgang dort; eine Achse ohne Drift wird nicht angefasst.
- [ ] **Hebung und Nachzug (Liefer-Punkt 2).** Je gedriftete Achse ein eigener
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
- [ ] **Belege (Liefer-Punkt 3).** Je gehobene Achse läuft der engste
      betroffene Lauf aus §3 und die Ausgabe wird gegen den Parent verglichen
      (Befunde, Meldungsform); `make image` endet erfolgreich, `make
      test-store`, `make test-replication` an PostgreSQL 18 und mit
      `PG_TEST_IMAGE` auf den gehobenen PostgreSQL-17-Digest, und
      `make test-integration` enden grün. Ein Lauf, der am Parent grün und am
      Diff rot ist, ist ein Befund (§6), kein Anlass, den Lauf zu lockern.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet,
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make test` grün, `make
      docs-check` grün (Exit direkt); `make mod-download` nur, falls `go.mod`
      berührt wird — erwartet nicht.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: kein öffentlicher Vertrag berührt (erwartet — gemessen
      unter §Berührte Spec-Stellen: kein Digest in `spec/`, `docs/user/`);
      die Träger-Nachzüge stehen in der Liefer-Punkt-2-Zeile.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist.

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
| `docs/plan/planning/open/slice-pin-digests-aktualisieren-2026-10.md` | update | §7 Messzeilen, Vergleich gegen den Parent |

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
  auf dem Parent. — **Ausgang:** offen bis Closure.
- **Tag-Frische von `DCHECK_IMAGE` ohne Digest-Messung:** `pin-stale-dcheck`
  meldet den Tag; der zugehörige Digest ist eine eigene Messung. Ein Tag-Wechsel
  ohne neuen `DCHECK_DIGEST` bliebe wirkungslos (der Digest sticht den Tag,
  `d-check.mk`). — **Ausgang:** offen bis Closure; Gegenprobe: `make docs-check`
  druckt am Diff die Version, deren Digest gezogen wurde (der Implementer
  nennt die gedruckte Zeile).
- **Neue Patch-Version von PostgreSQL oder Go ändert ein Testverhalten**
  (`make test-store`, `make test-replication`, `make test-integration`, beide
  PostgreSQL-Versionen): ein roter Lauf am Diff, grün am Parent, ist ein
  Befund; das Zeitverhalten der Replikations-Tests (`wal_sender_timeout`-Fenster)
  ist am wahrscheinlichsten betroffen. — **Ausgang:** offen bis Closure.
- **CI-Matrix `e2e.yml` PG 17/18:** die Digests dort sind Pins des Workflows;
  ein lokal grüner Lauf zeigt nicht, dass der Runner grün läuft. Der Slice
  berührt `.github/workflows/e2e.yml` (Pin-Wert, keine Struktur) — nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 ist die strukturelle Änderung
  Auslöser der Pflicht; eine reine Digest-Hebung ist es nicht, der reale
  Post-Push-Lauf von `e2e.yml` (beide Matrix-Legs) wird dennoch gelesen
  (`gh run list --workflow e2e.yml`), weil sich die Pins ändern, gegen die die
  Matrix läuft. Der Push ist nicht Teil dieses Auftrags; der Lauf steht aus.
  — **Ausgang:** weiter offen bis zum Post-Push-Lauf, dann Eintrag in §7.
- **`image-stale`-Befund `P1` und `make image`:** der neue `golang:1.27-alpine`-Digest
  ändert das Toolchain-Image jedes Docker-only-Laufs
  ([`AGENTS.md`](../../../../AGENTS.md) §3.1); ein Rückfall auf den Parent-Digest
  in einem der 9 Träger ließe zwei Toolchains koexistieren. — **Ausgang:** durch
  den `diff`-Suchlauf aus §3 entfallen, wenn die Zahl der Treffer dort die
  Erwartung trifft.
- **Offene Frage — Sensor-Lücke PG 17:** der Pin `postgres:17-alpine` in
  `e2e.yml` fehlt im Pin-Inventar von
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) und damit in
  `upstream-drift.yml`; die Drift dort bliebe unsichtbar.
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) ist `Accepted`;
  die Aufnahme ist eine Folge-ADR (Architect). Der Slice benennt sie, entscheidet
  sie nicht. — **Ausgang:** weiter offen: Beobachtung im Register
  (`BEO-PGC/...`, der Implementer legt den Eintrag an, falls die Handmessung
  Drift zeigt, die der Sensor nicht gemeldet hätte).
- **Offene Frage — stale Zitat:** [`harness/README.md`](../../../../harness/README.md)
  nennt für `make image-stale` die Bindung `ADR-0039` („Update =
  bewusster Digest-Commit“); der [ADR-Index](../../adr/README.md) führt diese
  Kennung als `Superseded` und mit anderem Gegenstand (Paketstruktur). Die tragende ADR für
  die Digest-Hebung ist [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
  Entscheidung 7. Der Slice zieht das nicht nach (fremde, nicht berührte
  Zeile); gemeldet. — **Ausgang:** Folge-Slice oder Nachzug durch den Planner
  mit der Closure.
- **Keine neue ADR nötig** — Begründung: Entscheidung 7 von
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) führt die
  Hebung ausdrücklich als „bewusster Commit“; der Slice ändert weder Tag noch
  Schwelle. Das gilt, solange die Messung keine Tag- oder Hauptversions-
  Änderung zeigt. — **Ausgang:** entfallen, wenn §7 keine solche Änderung nennt.

## 7. Closure-Notiz

Wird mit der Closure gefüllt (Inhalt, dann `git mv`, dann Häkchen der
Paarungs-Zeile — [`AGENTS.md`](../../../../AGENTS.md) §3.3). Messzeilen der
Planung (gemessen am Arbeitsstand `e5820a03`, 2026-10-03, volle Digests aus den
gedruckten Zeilen in `harness/`-Ausgaben des Implementers einzutragen, nicht
hier abgekürzt zu übernehmen):

- Drift-Messung P1–P9 (Implementer): Befehl, gedruckte Zeile, Parent-Kennung.
- PG-17-Handmessung (Implementer): Befehl, gedruckte Zeile.
- Vergleich der Gate-/Lauf-Ausgaben Parent gegen Diff je gehobene Achse.
- Ergebnis des `diff`-Suchlaufs aus §3, Gefundenes und Nichtgefundenes.
- Post-Push-Lauf `e2e.yml` (beide Legs) und `upstream-drift.yml`, sobald gepusht.

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

- **Was hat funktioniert:** (mit der Closure)
- **Was ging anders als geplant:** (mit der Closure)
- **Steering-Loop-Eintrag:** (mit der Closure; Kandidaten für einen
  Lerneintrag: die Sensor-Lücke PG 17 als benannte Spec-/ADR-Lücke, das stale
  Zitat `ADR-0039` in `harness/README.md`, ein Befund eines gehobenen
  Prüfwerkzeugs)
- **Beobachtungs-Register (`../observations/`):** (mit der Closure; Sichtung
  der Planung unter §8)
- **Folge-Slices:** (mit der Closure)
- **Risiken aus §6:** (mit der Closure, je Risiko genau ein Ausgang)
- **Drei Paarungen:** (mit der Closure, von der Slice-Closure selbst getragen)

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
