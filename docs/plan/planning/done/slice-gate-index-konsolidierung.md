# Slice slice-gate-index-konsolidierung: `AGENTS.md` §4 auf Regel + Zeiger kürzen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die über die eigene
DoD dieses Slice hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht). Dieses Repo läuft durchgehend wellenlos.

**Bezug:** kein einschlägiges `LH-*`, `ADR-*` oder `CO-*` — reine
Harness-Struktur-Konsolidierung ohne Vertrags- oder Entscheidungsgegenstand
(derselbe Befund wie bei `slice-105` §Bezug für den Baseline-Hub selbst).
Dies ist **Migrations-Folge-Slice 2 von 3** aus der Recherche
`/tmp/regelwerk-migration-v6.5.0-zu-v6.9.0.md` (Kurs-Welle 129,
`pt9912/ai-harness-course`), angekündigt in `slice-105` §1 Abgrenzung und
§7 Closure-Notiz (`done/slice-105-baseline-v6.9.0-materialisieren.md`).
Namensform nach `harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md`
— erster Slice-Plan dieses Repos mit Namens- statt Nummern-Kennung.

**Berührte Spec-Stellen:** — (geprüft: `grep -n "AGENTS.md\|harness/README" spec/lastenheft.md spec/pflichtenheft.md spec/architecture.md` liefert keine Treffer; reiner Harness-Doku-Vorgang ohne Vertrags- oder Draht-Bezug).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** `AGENTS.md` §4 (Quality Gates) wird von einer zweiten, den
Vertrag duplizierenden Gate-Tabelle auf **Regel + Zeiger** gekürzt — „kein
behauptetes Gate ohne Deckung in `harness/README.md` §Sensors" plus ein
Verweis dorthin. `harness/README.md` §Sensors bleibt die **einzige**
vollständige Gate-Tabelle (Target, Vertrag, Bindung inkl. ADR-Links,
Schwellen, Carveout-Verweisen). Anlass: Kurs-Welle 129
(`grundlagen-harness-dateien.md`, vendored unter
`.harness/baseline/v6.9.0/regelwerk/`) — eine reale Drift zwischen zwei
Fassungen derselben Gate-Zählung in einem anderen Repo („vier Module" vs.
„fünf Module"), weil kein Sensor zwei Tabellen auf Gleichheit hielt. Dieses
Repo hat dieselbe Struktur-Verwundbarkeit (zwei Tabellen, ein Zustand) auch
ohne dass bislang eine Zahlen-Drift real aufgetreten ist — verifiziert per
Zeile-für-Zeile-Vergleich am 2026-09-17: `AGENTS.md` §4 dupliziert in der
`Zweck`-Spalte real ADR-Links und Vertragsdetails (Schwellen-Rampe
„Einstiegsstufe 70 % → Endstufe 80 %", `SPEC-*`/`ARC-*`-Ausschluss-Regel,
Docker-Stufen-Namen), die inhaltsgleich bereits in `harness/README.md`
§Sensors stehen — trotz des bereits vorhandenen Hinweissatzes am
Tabellenende („Diese Tabelle listet auf; definiert wird hier nichts. Die
Bindung eines Targets … steht in `harness/README.md` §Sensors").

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Inhaltliche Änderung an einem Gate-Vertrag** (Schwellen, ADR-Bezug,
  Deckungsfläche). *Es wäre ein anderer Vorgang* — dieser Slice ändert nur,
  **wo** ein Vertrag steht, nie **was** er besagt; eine Schwellen-Änderung
  braucht ohnehin eine eigene ADR (`AGENTS.md` §3.6).
- **Struktur-Umbau von `harness/README.md` über §Sensors hinaus** (z. B.
  neue Spalten, neue Abschnitte). *Es wäre ein anderer Vorgang* — die
  Zieldatei bleibt formstabil, nur die Konkurrenz-Tabelle in `AGENTS.md`
  verschwindet.
- **Ein Sensor, der künftig zwei Gate-Tabellen gegeneinander hält.** *Bestand
  bleibt bewusst stehen* — nach diesem Slice gibt es nur noch **eine**
  Volltabelle; ein Vergleichs-Sensor für eine nicht mehr existierende zweite
  Tabelle wäre Zeremonie ohne Objekt. Genau das ist die Pointe der
  Kurs-Welle: das Problem strukturell auflösen (eine Quelle) statt es
  mechanisch bewachen (zwei Quellen + Wächter).
- **Nachträgliche Änderung an bereits abgeschlossenen `docs/reviews/**`- und
  `done/**`-Artefakten**, die die alte `AGENTS.md` §4-Tabellenform zitieren
  oder ihren Wortlaut referenzieren. *Bestand bleibt bewusst stehen* — diese
  Dateien sind Records eines abgeschlossenen Laufs (`AGENTS.md` §3.5-Analogie
  für Records: Zitat bleibt Beleg des damaligen Standes), keine lebenden
  Dokumente; keiner der gefundenen Verweise ist ein Hyperlink auf eine
  konkrete Tabellenzeile, sondern Prosa-Erwähnung von „`AGENTS.md` §4" als
  Abschnitt — der Abschnitt bleibt unter dieser Nummer bestehen.
- **Der dritte Folge-Slice dieser Migrations-Reihe** (DoD-Rot-vor-Grün-Regel
  + E2E-Gate-Typ-Abgrenzung, Kurs-Welle 135). *Ein Folge-Slice übernimmt
  es* — eigener, noch zu vergebender Slice-Name, inhaltlich unabhängig von
  der reinen Darstellungs-Konsolidierung hier.

## 2. Definition of Done

- [x] **LP1:** `AGENTS.md` §4 gekürzt auf Regel („kein behauptetes Gate ohne
      Deckung in `harness/README.md` §Sensors") + Zeiger dorthin — keine
      Vertragsdetails (ADR-Links, Schwellen, Docker-Stufen) mehr dupliziert;
      die Warnung „nur real existierende Targets, keine Halluzination"
      bleibt sinngemäß erhalten (gilt jetzt für die Zieltabelle).
- [x] **LP2:** Vollständigkeits-Grep über `.github/workflows/*.yml` erneut
      selbst ausgeführt (nicht nur aus dieser Planung übernommen,
      `AGENTS.md` §3.12/§3.13) — geprüft, ob eine **dritte** Fassung
      derselben Gate-Zählung existiert (Kurs-Welle 129 fand eine solche im
      Ursprungs-Repo in `.github/workflows/checks.yml`). Planungs-Vorabbefund
      (2026-09-17): kein Treffer — `ci.yml`/`e2e.yml` nennen Gate-**Namen**
      in einem Kommentar/Schritt-Titel und verweisen explizit auf
      `harness/README.md` §Sensors als Quelle, duplizieren aber keine
      Vertragsdetails (keine ADR-Links, keine Schwellenwerte). Falls der
      Implementer-Lauf denselben Befund bestätigt: keine Änderung an
      `.github/workflows/*.yml` nötig, aber die Bestätigung selbst gehört in
      den Bericht.
      **Implementer-Lauf bestätigt (2026-09-17):** eigenständiger Grep über
      alle drei Workflow-Dateien (`ci.yml`, `e2e.yml`, **und** `examples.yml`,
      das über die reine Planungs-Vorprüfung hinausgeht) — kein Treffer einer
      dritten Vertragsdetail-Fassung; nur Namens-/ADR-Kommentare mit Verweis
      auf `harness/README.md` §Sensors. Zusätzlich `Makefile`,
      `harness/mk/*.mk` und `docs/user/*.md` gegrept — kein Treffer einer
      weiteren Gate-Tabelle.
- [x] **LP3:** `make docs-check` grün — vor der Kürzung `grep -rn "AGENTS.md.*§4\|AGENTS\.md#" .` (ohne `.harness/baseline/`) laufen lassen: alle gefundenen Verweise sind Prosa-Erwähnungen des Abschnitts, keine Anker-Hyperlinks auf einzelne Tabellenzeilen — die Abschnittsnummer `§4` bleibt nach der Kürzung erhalten, daher keine Broken-Anchor-Erwartung; trotzdem nach der Änderung erneut `make docs-check` laufen lassen statt es anzunehmen.
      **Implementer-Lauf:** Vorab-Grep fand ausschließlich Prosa-Erwähnungen
      in `done/**`- und `docs/reviews/**`-Records (Out-of-Scope, §1) sowie in
      bereits `Accepted` ADRs (§3.5-immutabel, historische Folgepflicht-
      Beschreibungen zum jeweiligen ADR-Zeitpunkt) und eine Zeile in
      `.claude/agents/implementer.md`, die weiterhin trägt (§6 Risiko 4). Kein
      Anker-Hyperlink auf eine Tabellenzeile gefunden. `make docs-check` nach
      der Kürzung erneut ausgeführt: `d-check: 867 Datei(en) geprüft, 0
      Befund(e)`, `EXIT=0`.
- [x] `make gates` grün — Exit-Code direkt und ungepiped geprüft (`AGENTS.md`
      §3.9): `EXIT=0`.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Report:** `docs/reviews/review-slice-gate-index-konsolidierung.md`
      (0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO; kein Informationsverlust bestätigt,
      keine Fixrunde nötig).
- [x] Doku-Update: `AGENTS.md` §4 (siehe LP1); `harness/README.md` bleibt
      inhaltlich unverändert (bereits die vollständige Fassung).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-PGC/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; kein Zähler wird gesetzt, er folgt aus den Dateien. Keine
      Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert. **Beleg:** neues Verzeichnis
      `BEO-PGC/bindung-spalte-uneinheitlich-tief/` (Reviewer F-1, INFO, 1×,
      offen) — siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen). Siehe §6 unten — alle vier Risiken entfallen.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      geprüft in dieser Closure (Repo ohne Wellen-Betrieb).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §4 | update | Tabelle durch Regel + Zeiger auf `harness/README.md` §Sensors ersetzen — Kurs-Welle 129, siehe §1 Ziel |
| `.github/workflows/ci.yml`, `.github/workflows/e2e.yml` | prüfen, **kein Update erwartet** | Vorab-Grep fand nur Namens-Kommentare mit Verweis auf `harness/README.md` §Sensors, keine Vertragsdetail-Duplikate — Implementer bestätigt eigenständig (`AGENTS.md` §3.13) |
| `docs/plan/planning/observations/BEO-PGC/` | ggf. neues Verzeichnis oder `evidence/`-Ergänzung | siehe §8 Sichtungs-Schritt und §7 Closure-Notiz |

## 4. Trigger

**Start** (`next` → `in-progress`): sofort verfügbar — kein Abhängigkeits-Slice,
`in-progress/` ist zum Zeitpunkt dieser Planung leer (WIP-Limit frei,
geprüft 2026-09-17). Nachfolge-Slice von `slice-105` (`done/`), der diesen
Gegenstand ausdrücklich ausgeklammert und als Folge-Slice angekündigt hat.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Zeile-für-Zeile-Vergleich herausstellt, dass `AGENTS.md` §4 an mehreren,
  strukturell unabhängigen Stellen echte, nicht ableitbare Zusatzinformation
  trägt (z. B. eine eigene Ausnahmeregel, die `harness/README.md` nicht
  kennt) und die Kürzung dadurch zu einer inhaltlichen Entscheidung würde,
  die über reine Darstellungs-Konsolidierung hinausgeht — dann zurück zur
  Zerlegung, weil der Slice dann zwei Lieferwerte trägt (Kürzung *und*
  Entscheidung, was aus der verlorenen Information wird).
- `in-progress` → `open` (blockiert — Carveout?): falls `make docs-check`
  nach der Kürzung einen Anker- oder Struktur-Befund meldet, der sich nicht
  im Rahmen dieses Slice beheben lässt (z. B. ein d-check-Modul, das die
  §4-Tabelle strukturell voraussetzt) — dann Blocker, Carveout prüfen.

## 5. Closure-Trigger

DoD vollständig (§2) **und** `make gates` grün **und** Closure-Notiz mit
Lerneintrag geschrieben — kein Datum, keine externe Bedingung.

## 6. Risiken und offene Punkte

- **Die Kürzung verliert echte, nirgendwo sonst stehende Information aus
  `AGENTS.md` §4** (z. B. eine Nuance, die der Implementer beim genauen
  Vergleich übersieht). — **Ausgang: entfallen.** Zwei unabhängige
  Zeile-für-Zeile-Vergleiche (Reviewer: alle zehn gestrichenen Zeilen gegen
  `harness/README.md` §Sensors; Verifier: eigener Stichproben-Vergleich an
  drei weiteren, selbst gewählten Zeilen) finden keinen Informationsverlust
  — jedes Gegenstück ist inhaltsgleich oder detaillierter. Nicht eingetreten.
- **Eine dritte Fassung derselben Gate-Zählung wird übersehen** (Kurs-Welle
  129 fand eine in `.github/workflows/checks.yml` des Ursprungs-Repos;
  dieses Repo könnte eine analoge Stelle außerhalb der bereits geprüften
  `.github/workflows/*.yml` haben, z. B. `Makefile`-Kommentare,
  `docs/user/*`). — **Ausgang: entfallen.** Implementer, Reviewer und
  Verifier haben je eigenständig per `grep` über
  `.github/workflows/*.yml` (inkl. `examples.yml`), `Makefile`,
  `harness/mk/*.mk` und `docs/user/*.md` geprüft — kein Treffer einer
  dritten Vertragsdetail-Fassung, nur Namens-/Kommentar-Nennungen mit
  Verweis auf `harness/README.md` §Sensors. Nicht eingetreten.
- **Der Reviewer-Skill hat keine eigene HIGH-Regel für „zwei
  Doku-Tabellen, die denselben Vertrag beschreiben" außerhalb der bereits
  bestehenden Zwei-Quellen-Drift-Klasse** — geprüft, ob diese Klasse
  (`.harness/skills/reviewer.md` §HIGH) den Fall bereits deckt oder ob eine
  Schärfung nötig ist. — **Ausgang: entfallen.** Die bestehende HIGH-Klasse
  „Zwei-Quellen-Drift" deckt den Fall bereits; dieser Slice löst genau
  diese Klasse strukturell auf (eine Volltabelle statt zwei), statt eine
  neue Regel zu brauchen — von Reviewer (§8 Sichtungs-Schritt) und Verifier
  (§7) übereinstimmend bestätigt. Nichts zu schärfen, nicht eingetreten.
- **Referenz-Bruch in `.claude/agents/implementer.md`** (nennt
  `AGENTS.md` §4 als Beleg für „halluzinierte Targets sind verboten") — die
  Regel bleibt nach der Kürzung sinngemäß erhalten (Pointer statt Tabelle),
  eine Prüfung, ob die dortige Formulierung noch trägt, ist trotzdem nötig.
  — **Ausgang: entfallen.** Verifier hat Zeile 52 eigenständig geprüft: die
  Regel trägt sinngemäß weiter, die neue §4-Fassung führt sie sogar
  prominenter. Kein Referenzbruch, nicht eingetreten.

## 7. Closure-Notiz

**Besonderes Gewicht:** Dies ist der **erste** Slice-Plan mit
namensbasierter Kennung nach `harness/conventions/MR-002-slice-welle-
kennungen-sind-namen.md` (`slice-105` war der letzte nummerierte).

- **Was hat funktioniert:** Die neue Namenskonvention aus `MR-002` trug
  in ihrem ersten Anwendungsfall ohne jede Reibung. Der Slug
  `slice-gate-index-konsolidierung` war eindeutig, kollisionsfrei und ohne
  Nachschlagen einer nächsten freien Nummer vergeben — genau der Vorteil,
  den `MR-002` als Begründung nennt (Doppelvergabe wird an der Namensform
  selbst unwahrscheinlicher, nicht nur an durchgehaltener Disziplin). Der
  Name blieb über den gesamten Rollen-Durchlauf (Implementer, Reviewer,
  Verifier) stabil referenzierbar, ohne dass irgendeine Rolle ihn
  verwechselte oder eine numerische Nebenform bildete. Die inhaltliche
  Konsolidierung selbst (Regel + Zeiger statt Doppel-Tabelle) lief ebenso
  reibungslos: Implementer, Reviewer und Verifier haben unabhängig
  voneinander per Zeile-für-Zeile-Vergleich bzw. Stichprobe bestätigt, dass
  keine Information verloren ging.
- **Was ging anders als geplant:** Der Reviewer fand eine vorbestehende,
  von diesem Diff nicht verursachte Asymmetrie in `harness/README.md`
  §Sensors — die `Bindung`-Spalte von `generated-sync` verlinkt nur die
  Sensor-Datei statt die beiden ADRs inline zu tragen wie die übrigen
  Gates (F-1, INFO, kein Befund gegen diesen Diff). Nicht im Plan
  antizipiert, weil sie außerhalb des eigentlichen Änderungsgegenstands
  liegt; als eigenständige, kleine Beobachtung festgehalten (siehe unten),
  statt sie stillschweigend zu übergehen oder überzubewerten.
- **Steering-Loop-Eintrag:** Die neue Namenskonvention (`MR-002`) hat sich
  im ersten Anwendungsfall bewährt — kein Kollisions-, Lesbarkeits- oder
  Verwechslungsproblem trat auf. Das ist keine Regelschärfung und kein
  neuer Sensor (die Regel selbst ist unverändert richtig), sondern die
  erste reale Bestätigung einer bereits getroffenen Adaptions-Entscheidung.
  Kein neuer Zielort verkörpert — kein `liegt in`-Feld fällig
  (Baseline-Regelwerk `modul-06-roadmap.md`: das Feld steht nur, wenn mit
  diesem Slice wirklich etwas verkörpert wurde).
- **Beobachtungs-Register (`../observations/`):** Register auf eine
  passende bestehende Klasse geprüft
  (`grep -rli "generated-sync\|bindung-spalte\|inline.*adr-link"
  docs/plan/planning/observations/`) — keine passende Beobachtung
  gefunden. Neues Verzeichnis angelegt:
  `BEO-PGC/bindung-spalte-uneinheitlich-tief/` — Erstbeleg
  (`evidence/slice-gate-index-konsolidierung.md`), Stand `offen`, 1×. Kein
  Informationsverlust (beide ADRs über den Sensor-Link erreichbar), daher
  bewusst niedrigschwellig registriert statt als Risiko geführt.
- **Folge-Slices:** **Keine automatisch angelegt.** Zwei mögliche, noch
  nicht angelegte Folge-Themen aus dem laufenden Nutzer-Dialog, ohne eigene
  Kennung — Vergabe ist Sache der nächsten Planungs-Runde:
  (a) eine „Bewusstes Brechen"-Regel für DoD-Testbehauptungen (Kurs-Welle
  135, dritter Migrations-Folge-Slice aus §1 Abgrenzung), (b) möglicherweise
  ein Slice zur Aktivierung eines `tracked`-Moduls — gerade erst im
  Nutzer-Dialog besprochen, noch nicht bestätigt.
- **Risiken aus §6:** alle vier Risiken **entfallen** — keines ist real
  eingetreten (Einzelbegründungen siehe §6 oben, jeweils von mindestens
  zwei Rollen unabhängig geprüft).
- **Drei Paarungen:** *(nach dem `git mv` nach `done/`, siehe Commit 3)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Dieser Slice berührt ausschließlich
die Default-Sub-Area `*` (Kürzel `PGC`, `harness/conventions.md` §Modus-
Deklaration pro Sub-Area) — reine Harness-Meta-Doku (`AGENTS.md`,
`harness/README.md`), kein Code, kein separater Sub-Area-Kandidat. Die
Schwelle-≥2-Prüfung entfällt damit auf die bereits deklarierte
Default-Sub-Area; keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`grep -rli "gate.index\|doppelte.*tabelle\|zwei.*quellen\|zwei-quellen" docs/plan/planning/observations/`,
2026-09-17) — **keine Treffer** für eine bereits bestehende Beobachtung zu
genau diesem Gate-Index-Duplikat. Verwandt, aber nicht deckungsgleich:
`BEO-PGC/lese-doppelquelle` (bereits **verkörpert** seit `slice-029` — betrifft
SQL-View- vs. Go-Use-Case-Semantik, ein anderer Gegenstand). Der Reviewer-Skill
führt bereits die HIGH-Klasse **Zwei-Quellen-Drift** (`.harness/skills/reviewer.md`
§HIGH: „derselbe Zustand wird in zwei Dateien geführt, ohne dass der Gewinner
deklariert ist") — dieser Slice ist eine Instanz davon und braucht nach
aktueller Einschätzung keine neue HIGH-Regel, nur die Anwendung der
bestehenden; wird im Review gegengeprüft (§6 Risiko 3).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Default,
`harness/conventions.md`) — reiner Text-/Struktur-Slice ohne Code-Berührung,
kein eigener Begründungsblock pro Sub-Area nötig.
