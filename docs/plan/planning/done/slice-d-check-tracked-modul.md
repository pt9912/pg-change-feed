# Slice slice-d-check-tracked-modul: `d-check`-Modul `tracked` aktivieren

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-d-check` — Closure-Bedingung ist ein **kombinierter**
`make gates`-Lauf mit diesem Slice **und** `slice-d-check-trace-rtm`
gleichzeitig aktiv in derselben `.d-check.yml` (siehe
[`welle-d-check.md`](welle-d-check.md) §1/§3); das geht über die eigene
DoD dieses Slice hinaus.

**Bezug:** kein einschlägiges `LH-*`, `ADR-*` oder `CO-*` — reine
`d-check`-Konfigurationserweiterung ohne Vertrags- oder Entscheidungsgegenstand.
Ob die Aktivierung eines vormals opt-in-Moduls in das verpflichtende
`modules:`-Bündel eine eigene ADR braucht (Präzedenzfall:
[`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) für
das `hostpaths`-Modul, `BEO-PGC/gate-modul-abgeschaltet-trotz-regel`), klärt
der Architect-Zug bei Priorisierung (§4 Trigger, §6 Risiko 1) — dieser Punkt
ist bewusst offen gelassen, keine Vorwegnahme.

**Berührte Spec-Stellen:** — (geprüft:
`grep -n "d-check\|tracked" spec/lastenheft.md spec/pflichtenheft.md spec/architecture.md`
liefert keine Treffer; reiner Harness-Werkzeug-Vorgang ohne Vertrags- oder
Draht-Bezug).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das `d-check`-Modul `tracked` (Getrackt-Status auflösbarer,
existierender Link-/Bild-Datei-Ziele gegen den git-Index) wird in
`.d-check.yml` scharf geschaltet — Config-Block `tracked:` mit
`exempt-targets:` sowie Aufnahme von `tracked` in die `modules:`-Liste, womit
das Modul Teil von `make docs-check` (und damit `make gates`) wird, nicht nur
über sein bereits vorhandenes `doc-tracked`-Einzel-Target (`d-check.mk`,
`--enable tracked`) erreichbar. Anlass: `docs/user/e2e-abdeckung.md` und
`gen/cdc/stream/v1/*.pb.go`/`*_test.go` sind generierte, aber committete
Dateien mit vielfachen Verweisen (`harness/README.md`, diverse ADRs/Slices) —
eine künftige, zu breite `.gitignore`-Regel könnte sie unbemerkt aus dem
git-Index nehmen, während `links`/`anchors` auf dem Erzeuger-Rechner weiterhin
grün blieben (der Sensor liest das Dateisystem, nicht den Index) und erst ein
frischer Klon mit `target-missing` bricht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`--trace`/Requirements Traceability Matrix und `trace.coverage`.** *Ein
  Folge-Slice übernimmt es* — `slice-d-check-trace-rtm`, dieselbe Welle,
  unabhängiger Config-Block (`trace:` statt `tracked:`) derselben Datei. Die
  Adresse nimmt die Sendung an: Beide Slices sind Teil derselben
  Welle-Eröffnung, `slice-d-check-trace-rtm` schließt nicht vor diesem.
- **Eine neue ADR für die Bündel-Aufnahme.** *Es wäre ein anderer Vorgang* —
  ob eine ADR nötig ist, ist selbst eine architektonische Entscheidung
  (Modul 8: Planner→Architect-Zug); dieser Slice liefert die Konfiguration
  und die reale Grün-Messung, die ADR-Frage wird bei Priorisierung geklärt
  (§4/§6).
- **Sanierung tatsächlich gefundener `target-untracked`-Befunde.** *Bestand
  bleibt bewusst stehen, solange keiner auftritt* — der Planungslauf
  (2026-09-17, gepinnter Digest) fand real **0 Befunde** über 872 Dateien;
  taucht beim Implementer-Lauf dennoch einer auf, ist er zu beheben (kein
  neuer Vorgang), aber es ist keine vorab bekannte Arbeit.
- **Verzeichnis-Ziele und Symlink-Referenzen.** *Schicht-Abgrenzung* — das
  Modul `tracked` prüft sie laut Benutzerhandbuch (§4.20) explizit nicht;
  dieser Slice erweitert den Modul-Funktionsumfang nicht.
- **Versionsanhebung des gepinnten `d-check`-Digests.** *Bestand bleibt
  bewusst stehen* — der bereits gepinnte Stand (`v0.75.0`,
  `sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`)
  trägt `tracked` bereits seit `v0.37.0` (Benutzerhandbuch §11); kein
  funktionaler Grund für einen Versions-Sprung in diesem Slice.

## 2. Definition of Done

- [x] **LP1:** `.d-check.yml` trägt einen `tracked:`-Block
      (`exempt-targets: []`, oder konkrete Globs, falls der Implementer-Lauf
      reale `target-untracked`-Befunde findet, die bewusst ausgenommen werden
      sollen — mit Begründung je Glob) **und** `tracked` steht in der
      `modules:`-Liste. `make docs-check` (ungefiltert, Exit-Code direkt
      geprüft, `AGENTS.md` §3.9) bleibt grün — real erneut laufen lassen,
      nicht die Planungs-Vorabmessung übernehmen (`AGENTS.md` §3.12: eine
      Messung, kein „wird schon").
- [x] **LP2:** `harness/sensors/docs-check.md` §Grenze Punkt 3
      („Opt-in-Module nicht im Bündel") nachgezogen — `tracked` verlässt
      diese Aufzählung (es läuft jetzt im Bündel), die verbleibenden
      (`planning`, `vcs`, `commits`, `reviews`) bleiben stehen; die
      `## Bindung`-Sektion derselben Datei bekommt einen Eintrag für
      `tracked` (Modul, Vertrag, Konfigurationsort) analog zum bestehenden
      `hostpaths`-Eintrag. `.d-check.yml`s eigener Kommentar zur
      `tracked`-Sektion wird ergänzt (Aktivierungsstand, nicht mehr nur
      als Beispiel-Kommentar).
- [x] **LP3:** `make gates` grün — Exit-Code direkt und ungepiped geprüft
      (`AGENTS.md` §3.9). Zusätzlich (Closure-Trigger der Welle, nicht Teil
      dieser Einzel-DoD): sobald `slice-d-check-trace-rtm` ebenfalls in
      `.d-check.yml` geschrieben hat, ein gemeinsamer `make gates`-Lauf ohne
      Config-Kollision.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Ausgangs-Report [`review-slice-d-check-tracked-modul.md`](../../../reviews/review-slice-d-check-tracked-modul.md)
      (2 HIGH, 1 MEDIUM, 1 INFO), Delta-Review der Fixrunde
      [`review-slice-d-check-tracked-modul-fixrunde.md`](../../../reviews/review-slice-d-check-tracked-modul-fixrunde.md)
      (0 HIGH, nicht merge-blockierend).
- [x] Doku-Update: `harness/sensors/docs-check.md` (siehe LP2); `.d-check.yml`
      eigener Kommentarblock.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — siehe §7.
      Insbesondere geklärt: dieser Vorgang trägt **keine** weitere
      `evidence/`-Datei zu `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` —
      andere Beobachtungsklasse (dort: Modul fehlte trotz gewollter Regel;
      hier: Modul war bewusst opt-in dokumentiert, keine stille Abschaltung).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      geprüft von der **Welle-Closure** (`welle-d-check`), nicht hier
      einzeln (Repo mit Wellen-Betrieb für diese Welle).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` | update | `tracked:`-Block + Aufnahme in `modules:`-Liste (§1 Ziel) |
| `harness/sensors/docs-check.md` | update | §Grenze Punkt 3 nachziehen, `tracked`-Bindung ergänzen (LP2) |
| `harness/README.md` | update | **Plan-Nachzug:** `make doc-tracked` als Werkzeug (kein Gate) in §Sensors ergänzt, analog `make image-stale` — nicht im ursprünglichen §3, aber Träger-Nachzug derselben Aktivierung |
| `docs/plan/planning/observations/BEO-PGC/` | ggf. neues Verzeichnis oder `evidence/`-Ergänzung | siehe §8 Sichtungs-Schritt und §7 Closure-Notiz |

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung entscheidet auch die offene
ADR-Frage aus dem `Bezug`-Feld — der Architect-Zug (Planner→Architect,
Modul 8) bestätigt vor Implementierungsbeginn, ob eine eigene ADR nötig ist
oder der Verweis auf `ADR-0072`/`ADR-0075` als Präzedenzmuster in der
Closure-Notiz genügt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Implementer-Lauf reale `target-untracked`-Befunde findet, deren Auflösung
  selbst mehrschichtige Arbeit ist (z. B. ein generiertes Artefakt, das aus
  einem strukturellen Grund nicht getrackt werden *soll* und stattdessen sein
  Erzeugungs-Gate braucht) — dann ist „Modul aktivieren" und „Sanierung"
  zwei Lieferwerte, zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): falls `tracked` im
  `modules:`-Bündel eine Interaktion mit einem der sieben bestehenden Module
  auslöst, die sich nicht in diesem Slice beheben lässt (z. B. ein
  Konfigurationsfehler, den `.d-check.yml`s Schema-Validierung nicht
  vorhersagt) — dann Blocker, Carveout prüfen.

## 5. Closure-Trigger

DoD vollständig (§2) **und** `make gates` grün **und** Closure-Notiz mit
Lerneintrag geschrieben — kein Datum, keine externe Bedingung. Die
welle-weite Bedingung (kombinierter Lauf mit `slice-d-check-trace-rtm`) steht
zusätzlich in `welle-d-check.md` §3 und wird von der Welle-Closure geprüft,
nicht von diesem Slice allein.

## 6. Risiken und offene Punkte

- **Die Aufnahme in `modules:` verlangt eine eigene ADR, die dieser Slice
  nicht mitliefert** (Präzedenz `ADR-0072`/`ADR-0075` für `hostpaths`). —
  **Ausgang: entfallen** — der nachträgliche Architect-Zug
  ([`docs/reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md`](../../../reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md))
  hat die Frage geklärt: keine eigene ADR nötig, tragender Präzedenzfall ist
  Commit `f9e5a3c` (`structure`-Modul-Aktivierung), nicht `ADR-0072`/`ADR-0075`
  (die belegen den gegenteiligen Fall — eine Aktivierung mit realen Befunden,
  Korrektur an `Accepted`-ADRs und neuer Hard Rule).
- **`exempt-targets` bleibt bei Planung unbekannt** — der reale Umfang
  auflösbarer, existierender Link-/Bild-Ziele könnte beim vollen `modules:`-Lauf
  (statt der isolierten `--enable tracked`-Planungsmessung) anders ausfallen,
  weil andere Module denselben Scan-Baum anders beschneiden. —
  **Ausgang: entfallen** — sowohl der isolierte `make doc-tracked`-Lauf als
  auch der volle `make docs-check`/`make gates`-Bündel-Lauf melden über alle
  Stände dieses Slice-Laufs (Implementer, Review, Delta-Review, Verifikation:
  875→876→877→878 geprüfte Dateien, je nach neu hinzugekommenen Dateien im
  selben Lauf) durchgehend **0 Befunde** — identisches Ergebnis, kein
  `exempt-targets`-Eintrag nötig.
- **`harness/sensors/docs-check.md`s Grenze-Aufzählung (Punkt 3) wird nicht
  vollständig nachgezogen**, weil sie an mehreren Stellen der Datei referenziert
  wird (z. B. auch implizit im Bindung-Abschnitt). —
  **Ausgang: entfallen** — der Reviewer hat die Vollständigkeit gegen die
  ganze Datei geprüft (Ausgangs-Review, Negativbefund): §Grenze Punkt 3/4
  konsistent nachgezogen, einziger echter Fehler war der Präzedenz-Beleg in
  §Bindung (F-2, separat als eigenes Risiko behandelt und über den
  Architect-Zug behoben).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln.

- **Was hat funktioniert:** Die reale Grün-Messung vor der Priorisierung
  (872/875 Dateien, 0 Befunde am Planungsstand) hat sich über den ganzen
  Slice-Lauf hinweg bestätigt — `tracked` lief an jedem gemessenen Stand
  (Implementer, Review, Delta-Review, Verifikation: 875→878 Dateien) mit 0
  Befunden. Der Reviewer hat eigenständig recherchiert und den korrekten
  Präzedenzfall (`structure`-Modul, Commit `f9e5a3c`) gefunden, statt nur den
  falschen Implementer-Beleg zurückzuweisen (F-4) — das hat die
  Architect-Fixrunde erheblich beschleunigt.
- **Was ging anders als geplant:** Der Architect-Zug, den §4 Trigger
  explizit zur Start-Bedingung machte, wurde beim `next→in-progress`-Übergang
  übersprungen (Commit `262bcda`); der Implementer hat die ADR-Frage
  selbst beantwortet und dabei den falschen Präzedenzfall
  (`ADR-0072`/`ADR-0075`, tragen die gegenteilige Aussage) zitiert. Das löste
  eine volle Fixrunde aus: nachträglicher Architect-Zug → Implementer-Fix
  (Träger-Nachzug + Zitat-Korrektur) → Delta-Review (0 Findings) → Verifikation.
  Ohne den übersprungenen Architect-Zug wäre die Zitat-Korrektur vermutlich
  vor dem ersten Review-Lauf bereits richtig gestanden.
- **Steering-Loop-Eintrag:** Keine neue Regel geschärft — alle drei
  betroffenen Beobachtungsklassen liegen noch unter der 3×-Schwelle
  (`start-trigger-ohne-uebergabe-artefakt`, 1×) oder ihr Lese-Schritt gehört
  der Welle-Closure (`zitat-nennt-die-falsche-stelle`, 3×, Schwelle erreicht,
  Ausgang bei `welle-d-check`-Closure). `arbeit-ueberholt-stehenden-traeger`
  ist bereits seit `welle-20` in `AGENTS.md` §3.13 verkörpert — dieser Vorgang
  ist ein weiterer Beleg für eine bereits geschriebene Regel, kein neuer
  Steering-Loop-Auslöser.
- **Beobachtungs-Register (`../observations/`):** Drei Einträge fortgeschrieben:
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (weiterer Beleg, bereits
  verkörpert) · `BEO-PGC/zitat-nennt-die-falsche-stelle` (dritter Beleg,
  Schwelle erreicht, Ausgang bei Welle-Closure) · neu angelegt:
  `BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (1×, Erstauftreten: der
  Architect-Zug als Plan-Start-Bedingung fand nicht statt). Keine Beobachtung
  zu `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` — andere Klasse
  (dokumentiertes Opt-in statt stiller Abschaltung, siehe §2 DoD).
- **Folge-Slices:** keine — die ADR-Frage ist über den Architect-Zug
  abschließend geklärt (keine ADR nötig), keine offene Arbeit verbleibt.
- **Risiken aus §6:** alle drei mit Ausgang „entfallen" — siehe §6.
- **Drei Paarungen:** von der Welle-Closure `welle-d-check` geprüft (Repo mit
  Wellen-Betrieb für diese Welle, Modul 6 §Wann Arbeit eine Welle braucht).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Dieser Slice berührt ausschließlich
die Default-Sub-Area `*` (Kürzel `PGC`, `harness/conventions.md` §Modus-
Deklaration pro Sub-Area) — reine Harness-Werkzeug-Konfiguration
(`.d-check.yml`, `harness/sensors/docs-check.md`), kein Code, kein eigener
Sub-Area-Kandidat. Die Schwelle-≥2-Prüfung entfällt damit auf die bereits
deklarierte Default-Sub-Area; keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`grep -rli "tracked\|target-untracked\|gitignore" docs/plan/planning/observations/`,
2026-09-17) — **keine Treffer** für eine bestehende Beobachtung zu genau
diesem Gegenstand. Direkt verwandt, aber nicht deckungsgleich:

- `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` (**verkörpert**, 1×,
  `slice-078`, `hostpaths`) — dieselbe *Klasse* („eine gewollte Regel steht,
  aber ihr Modul fehlt im Bündel"), anderes Modul; ob dieser Slice eine
  zweite `evidence/`-Datei zu diesem Eintrag beiträgt oder eine eigene
  Beobachtung braucht, entscheidet die Closure (siehe DoD-Punkt oben) — hier
  bewusst **nicht** vorweggenommen, weil `tracked` anders als `hostpaths`
  bereits *dokumentiert* opt-in war (`harness/sensors/docs-check.md` §Grenze
  Punkt 3), nicht stillschweigend abgeschaltet.
- `BEO-PGC/regel-weiter-als-ihr-sensor` (**offen**, 3×, Schwelle erreicht,
  Ausgang seit `welle-20`s Schließung noch nicht zugewiesen) — anderer
  Gegenstand (Reichweite einer Regel vs. Sensor-Deckung bei `hostpaths`/
  `coverage-gate`), der fällige Lese-Schritt gehört der Closure von
  `welle-d-check` (Modul 6 Closure-Schritt 3a), nicht diesem Slice.
- `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (**verkörpert**, 4×,
  geschlossen mit `make generated-sync`) — verwandtes Thema (generierte,
  committete Artefakte), aber bereits aufgelöst über einen anderen
  Mechanismus (Byte-Sync-Gate statt Tracked-Status); `tracked` prüft eine
  andere Eigenschaft (Index-Mitgliedschaft, nicht Inhaltsgleichheit) und ist
  komplementär, kein Duplikat.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Default,
`harness/conventions.md`) — reiner Konfigurations-/Doku-Slice ohne
Produktionscode-Berührung, kein eigener Begründungsblock pro Sub-Area nötig.
