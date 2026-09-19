# Slice release-doku-releasing: `docs/user/releasing.md` für Betreiber/Maintainer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(derivativ — kein eigener Entscheidungspunkt, dokumentiert die Summe der
vier vorangehenden Slices).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum;
`docs/user/*` selbst ist Rang 6 der Source Precedence, kein Spec-Stratum).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "dann mach weiter bis die welle geschlossen ist").

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** `docs/user/releasing.md` neu anlegen — beschreibt für Betreiber
und Maintainer den nach den vier vorangehenden Slices real existierenden
Release-Prozess: Versionierung (`docs/user/version.md`), wie ein Release
ausgelöst wird (Tag-Form `v<SemVer>`, wer taggen darf), was automatisch
passiert (Build über `make image`, Push nach GHCR und Docker Hub,
GitHub-Release mit Digest-Pin, Docker-Hub-Beschreibungs-Sync), und die
drei advisory-Begleit-Workflows (`image-scan.yml`, `upstream-drift.yml`)
samt ihrer Bedeutung für den Betrieb.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Behauptung, dass bereits ein echter Release stattgefunden hat —
  `docs/user/releasing.md` beschreibt den **Mechanismus**, nennt aber
  explizit, dass zum Zeitpunkt dieses Slice noch kein erster Tag gesetzt
  wurde (`AGENTS.md` §3.12: keine ungeprüfte Behauptung als Feststellung).
- Änderungen an `docs/user/benutzerhandbuch.md` — das Handbuch adressiert
  Betreiber im laufenden Betrieb (Umgebungsvariablen, SQL-Zugriffe), nicht
  den Release-Prozess des Projekts selbst; beide Dokumente bleiben
  getrennt (unterschiedliche Zielgruppen-Fragen).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `docs/user/releasing.md` existiert, trägt Version/Stand-Kopf analog
      `benutzerhandbuch.md`, und beschreibt real existierende Artefakte
      (kein Verweis auf einen noch nicht implementierten Workflow) — jede
      genannte Datei/jeder genannte Workflow real gegen den Baum geprüft
      (`tools/harness/release-tag-info.sh`, `docs/user/version.md`,
      `.github/workflows/release.yml`/`image-scan.yml`/
      `upstream-drift.yml`, `README.md` — alle real existent). Alle sechs
      im Dokument genannten Tag-Beispiele real gegen
      `release-tag-info.sh` verifiziert (drei gültig, drei ungültig,
      jeweils mit dem im Dokument behaupteten Ergebnis). Enthält
      ausdrücklich den Hinweis, dass noch kein realer Release-Tag
      gesetzt wurde (`AGENTS.md` §3.12: keine ungeprüfte Behauptung als
      Feststellung).
- [x] `harness/README.md` Source-Precedence-Zeile 6 (`docs/user/*`) —
      real geprüft: der `<!-- d-check:ignore -->`-Kommentar traf nicht
      mehr zu (das Verzeichnis trägt bereits sechs Dateien, jetzt sieben
      mit `releasing.md` — die im Kommentar genannte Bedingung „im
      frischen Repo selten vorhanden" ist damit widerlegt). Entfernt und
      auf einen echten Link umgestellt (`[docs/user/](../docs/user/)`,
      analog Zeile 4 derselben Tabelle). Dieselbe, wortgleiche Zeile
      existierte zusätzlich (im Slice-Plan nicht vorab benannt) in
      `AGENTS.md` §2 Zeile 6 — im selben Zug konsistent mitgezogen
      (Plan-Nachzug, §3 unten).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      1 HIGH (F-1: „Zitat nennt die falsche Stelle" — `§3` statt `§4` für
      den `version.md`-Abgleich) und 1 INFO (F-2: `releasing.md` fehlt in
      der „Siehe:"-Liste von `README.md`, vorbestehende Lücke, kein
      neues Muster) in derselben Fixrunde behoben bzw. zur Kenntnis
      genommen, kein offenes HIGH.
- [ ] Doku-Update für `docs/user/releasing.md` selbst — entfällt als
      eigener Punkt, da bereits §2 Zeile 1 dieser DoD dieselbe Datei
      trägt; kein zusätzlicher öffentlicher Vertrag jenseits ihrer.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neuer Beleg `evidence/slice-release-doku-releasing.md` unter der bestehenden `BEO-PGC/zitat-nennt-die-falsche-stelle/` (bereits verkörpert, jetzt 7. Beleg, dritter Fundort in dieser Welle).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice ist der **letzte** von `welle-release-pipeline-adr-0051` — die Welle schließt unmittelbar im Anschluss, die Prüfung läuft dort.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/releasing.md` | neu | Betreiber-/Maintainer-Doku des realen Release-Prozesses. |
| `AGENTS.md` §2 Zeile 6 | update (Plan-Nachzug) | dieselbe, wortgleiche `docs/user/*`-Zeile mit demselben veralteten `<!-- d-check:ignore -->`-Kommentar existierte hier zusätzlich zu `harness/README.md` — im selben Zug konsistent mitgezogen, sonst bliebe eine der beiden Kopien stehen. |
| `harness/README.md` | ggf. update | `docs/user/*`-Zeile, `<!-- d-check:ignore -->`-Kommentar auf Zutreffen geprüft. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `release-version-und-workflow`,
`release-image-scan`, `release-upstream-drift` und
`release-hub-description` liegen alle in `done/` (siehe Welle-Datei §5
Abhängigkeiten — dokumentiert den vollständigen Endzustand).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — reiner Doku-Slice.
- `in-progress` → `open` (blockiert — Carveout?): einer der vier
  Vorgänger-Slices liegt entgegen der Annahme noch nicht in `done/` —
  dann zurück nach `open`, bis die Voraussetzung real erfüllt ist.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Eine Doku, die einen Prozess beschreibt, der real noch nie durchlaufen
  wurde (kein Tag je gesetzt), kann Details enthalten, die erst beim
  ersten echten Release als falsch auffallen (z. B. ein Feld im
  GitHub-Release-Formular, das anders aussieht als angenommen).
  **Ausgang:** weiter offen, löst sich mit dem ersten echten Release —
  dann Nachtrag/Korrektur als eigener kleiner Folge-Vorgang, kein
  Blocker für diesen Slice.

## 7. Closure-Notiz

- **Was hat funktioniert:** Jede in `docs/user/releasing.md` genannte
  Datei, jedes genannte Make-Target und alle sechs Tag-Beispiele wurden
  vor der Übergabe real gegen den Baum bzw. `tools/harness/
  release-tag-info.sh` geprüft (`AGENTS.md` §4/§3.12) — kein
  Halluzinieren nicht existierender Targets, real vom Reviewer und
  Verifier unabhängig reproduziert. Der Nutzer-Hinweis auf das
  Schwester-Repo d-check (`release-hub-description`) ließ sich direkt
  als Betriebswissen in dieses Dokument übernehmen (Docker-Hub-Token-
  Scope `read/write/delete`), bevor der Nutzer das eigene Token überhaupt
  angelegt hat.
- **Was ging anders als geplant:** Der Reviewer fand 1 HIGH (F-1:
  „Zitat nennt die falsche Stelle" — `§3` statt `§4` für den
  `docs/user/version.md`-Abgleich) und 1 INFO (F-2: `releasing.md` fehlt
  in der „Siehe:"-Liste von `README.md`, vorbestehende Lücke). F-1
  behoben, F-2 als bereits bestehende, nicht-blockierende Lücke zur
  Kenntnis genommen statt in diesem Slice mitgelöst (kein Teil dieses
  Slice-Umfangs). Zusätzlich, außerhalb des ursprünglichen §3-Plans:
  `AGENTS.md` §2 Zeile 6 trug denselben veralteten
  `<!-- d-check:ignore -->`-Kommentar wie `harness/README.md` — im selben
  Zug konsistent mitgezogen, sonst wäre eine der beiden Kopien stehen
  geblieben.
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine geschärfte Regel —
  F-1 ist der 7. Beleg der bereits verkörperten Beobachtung
  `BEO-PGC/zitat-nennt-die-falsche-stelle`, und der **dritte** Fundort
  dieser Klasse in drei aufeinanderfolgenden Slices derselben Welle
  (`release-hub-description` F-1/F-2, jetzt hier). Bemerkenswert: alle
  drei Fundstellen dieser Welle sind interne Selbst-Verweise innerhalb
  eines neu geschriebenen Dokuments (Abschnitts-/ADR-Zitat auf dieselbe
  bzw. eine benachbarte Datei), nicht Verweise auf ein fremdes Dokument
  wie bei den vier früheren Belegen — eine mögliche Verfeinerung der
  Beobachtung für den Lese-Schritt der Welle-Closure, nicht hier
  entschieden.
- **Beobachtungs-Register (`../observations/`):** neuer Beleg
  `evidence/slice-release-doku-releasing.md` unter der bestehenden
  `BEO-PGC/zitat-nennt-die-falsche-stelle/` (bereits verkörpert, jetzt
  7. Beleg; Zähler in `state.md` nachgezogen).
- **Folge-Slices:** keine — letzter Slice der Welle
  `welle-release-pipeline-adr-0051`.
- **Risiken aus §6:** ein Risiko, ein Ausgang — Doku beschreibt einen
  real noch nie durchlaufenen Prozess (kein Tag je gesetzt) →
  **weiter offen**, löst sich mit dem ersten echten Release (dann
  Nachtrag/Korrektur als eigener kleiner Folge-Vorgang, kein Blocker für
  diesen Slice).
- **Drei Paarungen:** dieser Slice ist der letzte von
  [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md)
  — die Prüfung läuft regelkonform bei deren unmittelbar folgender
  Closure (§2 DoD-Zeile „im Repo mit Wellen von der nächsten
  Welle-Closure"). Vorab-Hinweis für diese Prüfung: kein
  `liegt in`-Feld in diesem Slice; kein Folge-Slice; der neue
  Beobachtungs-Beleg liegt unter
  `BEO-PGC/zitat-nennt-die-falsche-stelle/evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „Betreiber-Dokumentation"
(`docs/user/*`) — bereits mehrfach berührt (`benutzerhandbuch.md`,
`ADR-0051` selbst nennt `docs/user/*` als Rang-6-Ort), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (3×, verkörpert)
betrifft strukturell dieselbe Klasse (Meta-Pflicht an einer
Doku-Datei übersehen) — hier auf `releasing.md` statt
`benutzerhandbuch.md` angewandt: auch `releasing.md` braucht künftig eine
eigene Versionshistorie-Disziplin, sobald es geändert wird. Kein weiterer
Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
