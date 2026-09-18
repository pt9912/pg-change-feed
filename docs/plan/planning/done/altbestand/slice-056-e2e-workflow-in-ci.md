# Slice slice-056: E2E-Workflow in CI (`make test-integration` auf GitHub Actions)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist ausschließlich die eigene
DoD (ein neuer, nicht-blockierender Workflow); kein repo-weites *Mehr* wie
bei welle-13/14/15 (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht). Orthogonal zu `welle-15` (NATS) —
reine CI-Infrastruktur, keine CDC-Fähigkeit, blockiert und wird nicht
blockiert von `welle-15`.

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (vollständige
Testumgebung automatisiert und reproduzierbar), [ADR-0051](../../adr/0051-cicd-pipeline-github-actions.md)
(Entscheidung 2 — Folgepflicht: „die DB-gebundenen Testziele … laufen nach
Bedarf, nicht zwingend bei jedem PR (Umfang: Folge-Slice)" — dieser Slice
löst genau diese Folgepflicht ein, mit dem hier getroffenen konkreten
Design).

**Berührte Spec-Stellen:** — (reine CI-Infrastruktur, keine Spec-Stelle
berührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein neuer Workflow `.github/workflows/e2e.yml` ruft
`make test-integration` bei jedem Pull Request und Push auf (derselbe
Trigger wie `ci.yml`, `tags-ignore: ['**']`), ist aber **nicht** als
Required-Status-Check im Branch-Schutz hinterlegt — ein roter Lauf ist
sichtbar (rotes Check-Symbol im PR), blockiert aber keinen Merge
(Nutzerentscheidung 2026-09-13: „Bei jedem PR, aber als eigener,
nicht-blockierender Workflow"). Das löst `ADR-0051`s Folgepflicht ein, die
DB-gebundenen Testziele „nach Bedarf" statt zwingend an `ci.yml` zu binden
— hier: eigener Workflow statt eigener Job in `ci.yml`, damit sein
(erwartbar längerer) Laufzeit- und Ressourcenbedarf `ci.yml`s Signal nicht
verzögert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`test-store`/`test-replication` in CI** — beide sind eigene,
  DB-gebundene Testziele mit eigenem Testcontainer-Aufbau
  (`tools/harness/run-store-tests.sh`/`run-replication-tests.sh`); `ADR-0051`
  spricht sie zusammen mit `test-integration` als Gruppe an, aber dieser
  Slice liefert bewusst nur den Black-Box-E2E-Rundlauf (`LH-QA-POR-003`s
  konkreteste Ausprägung) — ein Folge-Slice kann `test-store`/
  `test-replication` nach demselben Muster ergänzen, falls gewünscht.
- **Required-Status-Check / Branch-Protection-Konfiguration** — die
  GitHub-Repository-Einstellung selbst (welche Checks required sind) ist
  eine Repository-Admin-Aktion außerhalb des Repo-Inhalts; dieser Slice
  liefert nur den Workflow so, dass er *nicht* als required eingetragen
  werden muss, um nicht-blockierend zu sein (kein `needs:`/`required`-
  Constraint im Workflow selbst kann das erzwingen — das ist eine
  GitHub-Einstellung, keine YAML-Eigenschaft).
- **`release.yml`/`image-scan.yml`/`upstream-drift.yml`/`hub-description.yml`/
  `version.md`** — vollständig separate Folgepflicht aus `ADR-0051`
  (dort als „Folge-Slice `slice-040`" benannt, bisher nicht angelegt);
  dieser Slice berührt ausschließlich den PR-/Push-Testlauf, nicht den
  Release-Pfad.
- **Neue Make-Targets oder Änderungen an `tools/harness/run-integration-tests.sh`**
  — der Workflow ruft das bestehende `make test-integration` unverändert
  auf (`AGENTS.md` §3.1, Docker-only: kein Workflow-Schritt dupliziert
  oder umgeht ein bestehendes Ziel).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `.github/workflows/e2e.yml` neu: Trigger `pull_request` + `push`
      (`branches: ['**']`, `tags-ignore: ['**']`, wie `ci.yml`),
      `permissions: {}` auf Workflow-Ebene, `contents: read` auf Job-Ebene,
      Action-Pinning nach `AGENTS.md` §3.8 (SHA + `# vX.Y.Z`-Kommentar).
      Ein Job ruft ausschließlich `make test-integration` auf (kein Inline-
      Shell, das das Target umgeht — `AGENTS.md` §3.1).
- [x] Der Workflow ist real über einen Push auf `main` ausgelöst worden und
      zeigt ein sichtbares Check-Ergebnis — vier unabhängige reale Läufe,
      alle `completed`/`success`, von Reviewer UND Verifier unabhängig per
      `gh run view`/`gh run list` gegengeprüft (siehe
      `docs/reviews/verify-slice-056.md`).
- [x] Rot-Beleg für Nicht-Blockierung — **substanziell erfüllt, mit
      dokumentierter Abweichung vom geplanten Beleg-Weg:** Kein
      dedizierter Test-PR mit absichtlich fehlschlagendem Schritt wurde
      erzeugt. Der Verifier ersetzte ihn durch einen stärkeren
      strukturellen Beleg (`gh api .../branches/main/protection` → `404
      "Branch not protected"`) — siehe §7 *Was ging anders als geplant*.
- [x] `make gates` grün (unverändert — dieser Slice ändert kein bestehendes
      Gate-Target).
- [x] Review durchgeführt, Report unter `docs/reviews/review-slice-056.md` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `harness/README.md` §Werkzeuge neuer Eintrag für den
      Workflow (kein Gate — analog zum bestehenden `ci.yml`-Eintrag im
      Abschnitt „Aktueller Lauf-Status").
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo führt kein Brownfield-Bootstrap, die Datei existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/e2e.yml` | neu | eigener, nicht-blockierender `make test-integration`-Workflow |
| `harness/README.md` | update | §Werkzeuge-Eintrag für den neuen Workflow |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-039` liegt in `done/` (`ci.yml`
existiert bereits als Vorbild für Trigger-/Permissions-/Pinning-Form),
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten — ein einzelner neuer Workflow nach bereits etabliertem
  `ci.yml`-Muster.
- `in-progress` → `open` (blockiert — Carveout?): Falls sich zeigt, dass
  `make test-integration` auf dem GitHub-hosted Runner strukturell nicht
  läuft (z. B. Ressourcen-/Zeitlimit des Compose-Stacks auf dem
  Standard-Runner überschritten) — dann Carveout mit Folge-Slice
  (größerer Runner, Selbst-gehosteter Runner, oder Reduktion des
  Testumfangs).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** ein realer,
sichtbarer Workflow-Lauf auf GitHub nach einem Push **und** Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- GitHub-Actions-Workflows lassen sich nicht in Docker simulieren
  (`BEO-PGC/github-actions-unverifizierbar-lokal`, bereits 1× aus
  `slice-039`) — ob `make test-integration` auf dem echten Runner
  tatsächlich innerhalb des Zeitlimits durchläuft, bleibt bis zum ersten
  realen Push unbewiesen. **Ausgang: entfallen** — vier unabhängige,
  reale Läufe auf `main` (u. a. Run `34788082848`), von Reviewer UND
  Verifier unabhängig per `gh run view`/`gh run list` gegengeprüft, alle
  `completed`/`success` in 3m29s–3m46s. Zugleich neuer Beleg
  `evidence/slice-056.md` für `BEO-PGC/github-actions-unverifizierbar-lokal`
  (Zähler damit 2×, weiter unter der 3×-Schwelle) — das strukturelle
  Verifikationsgrenze-Muster (kein Docker-only-Beleg vor dem ersten realen
  Push) trat erneut auf, auch wenn der Lauf selbst grün war.
- `make test-integration` ist deutlich länger als `make gates`/`make test`
  (Compose-Stack hochfahren, Schema-Rollout, mehrere Testphasen) — auf
  dem Standard-GitHub-hosted-Runner könnte die Gesamtlaufzeit das
  Standard-Timeout (hier: wie `ci.yml`, 30 Minuten) unerwartet knapp
  werden lassen. **Ausgang: entfallen** — der reale Lauf blieb
  durchgehend unter 7 % des 60-Minuten-Budgets (das großzügiger als
  `ci.yml`s 30 Minuten gewählte Timeout, siehe §3/Implementierung), über
  vier unabhängige Läufe hinweg reproduziert.
- Ohne Required-Status-Check könnte ein dauerhaft roter `e2e.yml`-Lauf
  unbemerkt bleiben (niemand ist gezwungen hinzusehen) — dasselbe
  Alarmmüdigkeits-Risiko, das `ADR-0051` für die advisory-Läufe
  (CVE-Scan, Pin-Freshness) bereits benennt, hier auf einen
  PR-Trigger-Workflow übertragen statt auf einen Nachtlauf. **Ausgang:
  weiter offen** — neuer Registereintrag
  `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×,
  `evidence/slice-056.md`), vom Verifier empfohlen: Das Risiko ist real
  und durch diesen Slice allein nicht reduzierbar (fehlende
  Required-Status-Check-Konfiguration ist bewusst außerhalb des Umfangs,
  siehe §1).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Das etablierte `ci.yml`-Muster (Trigger,
  `permissions: {}`, Action-Pinning) ließ sich verlustfrei auf einen
  zweiten, eigenständigen Workflow übertragen — kein Duplikat-Gate, nur
  Automatisierung des bestehenden `make test-integration`-Targets
  (`AGENTS.md` §3.1). Die zwingende Reihenfolge `make image` →
  `make test-integration` (kein `build:`-Block in `compose.yaml`,
  `ADR-0044`) wurde von Implementer, Reviewer und Verifier unabhängig
  korrekt hergeleitet. Vier reale, unabhängige GitHub-Actions-Läufe
  bestätigten alle drei ursprünglich offenen Risiken in einem Zug.
- **Was ging anders als geplant:** Der geplante Rot-Beleg (ein absichtlich
  fehlschlagender Testlauf, der einen PR trotzdem mergebar zeigt) wurde
  **nicht** als Verhaltens-Experiment durchgeführt — weder von
  Implementer/Reviewer noch vom Verifier. Der Verifier ersetzte ihn durch
  einen stärkeren strukturellen Beleg: `gh api
  repos/pt9912/pg-change-feed/branches/main/protection` liefert `404
  "Branch not protected"` — ohne jede Branch-Protection-Regel kann
  kategorisch kein Workflow einen Merge blockieren, nicht nur im
  getesteten Einzelfall. Diskrepanz zwischen DoD-Wortlaut
  (Verhaltensbeleg) und tatsächlichem Beleg-Weg (Konfigurationsbeleg)
  vom Verifier benannt — dieselbe Finding-Klasse wie Review-F-1
  (DoD-Wortlaut deckt den tatsächlichen Umfang nicht präzise ab); kein
  eigener Fixrunden-Anlass, da der gelieferte Beleg stärker ist als der
  geplante.
- **Steering-Loop-Eintrag:** keiner — beide Findings (Review-F-1,
  DoD-Wortlaut-Diskrepanz beim Rot-Beleg) sind isolierte
  LOW-Formulierungs-Ungenauigkeiten in genau diesem Slice-Plan, kein
  wiederkehrendes Muster über mehrere Vorgänge (Schwelle 3× nicht
  erreicht, jeweils Erstauftreten).
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-056.md` in `BEO-PGC/github-actions-unverifizierbar-lokal/`
  ergänzt — Zähler steht damit bei 2×, weiter unter der 3×-Schwelle.
  Zusätzlich neu angelegt: `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/`,
  Beleg `evidence/slice-056.md` (1×, vom Verifier empfohlen).
- **Folge-Slices:** keine — `test-store`/`test-replication`-CI-Integration
  und der Release-Pfad (`slice-040`-Familie) bleiben wie in §1 benannter,
  unbeanspruchter Bestand ohne eigene Kennung in diesem Slice.
- **Risiken aus §6:** zwei entfallen (Runner-Unsicherheit,
  Timeout-Knappheit — beide durch vier reale grüne Läufe widerlegt), eines
  weiter offen (Alarmmüdigkeit bei stillem Dauer-Rot — neuer
  Registereintrag).
- **Drei Paarungen:** Anker — kein `liegt in`-Feld in diesem Eintrag (kein
  Steering-Loop-Eintrag verkörpert), Paarung entfällt für diesen Slice.
  Folge-Slice — keiner genannt, Paarung entfällt. Register — beide
  zitierten Kennungen (`BEO-PGC/github-actions-unverifizierbar-lokal`,
  `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`) existieren als
  Verzeichnis mit nicht leerem `evidence/` — geprüft nach dem `git mv`
  nach `done/`.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Ein Treffer: `BEO-PGC/github-actions-unverifizierbar-lokal`, Stand 1×
(aus `slice-039`) — dieser Slice wird ihn voraussichtlich auf 2× heben
(derselbe strukturelle Grund: kein Docker-only-Beleg vor dem echten Push
möglich), noch unter der 3×-Schwelle.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
