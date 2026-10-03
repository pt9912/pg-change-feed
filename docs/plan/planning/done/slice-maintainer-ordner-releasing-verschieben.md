# Slice maintainer-ordner-releasing-verschieben: Ordner `docs/maintainer/` anlegen und `releasing.md` dorthin verschieben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(`Accepted`; Festlegung 2: `releasing.md` dauerhaft außerhalb des Gates; das Gate
und seine Listen sind betroffen),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (`Accepted`; Rang 6
und `docs/user/version.md` als Quelle der Wahrheit der Version),
[`ADR-0123`](../../adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) (`Accepted`;
drei Links auf `releasing.md`), [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Zitat-Korrektur für die drei ADR-Fundstellen). Anforderungs-Bezug: keiner — reine
Ablage-Ordnung der Dokumentation; die Betreiber-Dokumentation
([`LH-QA-OPS-001`](../../../../spec/lastenheft.md)) bleibt unter `docs/user/`.

**Berührte Spec-Stellen:** — (keine Spec-Stelle genannt; `git grep -n "releasing" spec`
am Stand `81d0fa96` zu messen, erwartet 0 Treffer — vom Implementer zu belegen).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, auf Auftrag des Auftraggebers (wörtlich): „können wir in
docs einen maintainer Ordner anlegen und das releasing.md Dokument dorthin
verschieben“. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

**Ziel:** `docs/user/releasing.md` liegt per reinem `git mv` unter
`docs/maintainer/releasing.md`, jeder Verweis darauf im aktiven Baum und jeder
dadurch gebrochene Link in Records ist nachgezogen, und das Handbuch-Gate führt keine
Ausnahme mehr für eine Datei, die nicht mehr in `docs/user/` liegt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`docs/user/version.md` bleibt.** Sie ist die Quelle der Wahrheit der Version, der
  Release-Workflow (`.github/workflows/release.yml`) und zwei SDK-Workflow-Kommentare lesen
  den Pfad (`git grep -n -F "docs/user/version.md"` am Stand `81d0fa96`: 39 Trefferzeilen
  im aktiven Baum ohne Records); ein Pfadwechsel wäre ein Workflow-Eingriff (`AGENTS.md`
  §3.10: realer Post-Push-Lauf nötig) ohne Nutzen für den Auftrag.
  `releasing.md` §2 nennt die Beziehung weiterhin, nur der Linkpfad ändert sich.
- **Keine Inhaltsänderung an `releasing.md`** außer Pfaden, Verweisen und dem
  Zweck-/Zielgruppen-Absatz in §1 — der Auftrag ist eine Verschiebung; Inhalt
  und Versionsstand des Dokuments (Zeile `Version:`) bleiben bis auf die Änderungshistorie-Zeile
  für diese Pfadänderung stehen.
- **Keine Aufteilung von `releasing.md`.** Ob Betreiber-nahe Teile (§2 Versionierung,
  Image-Tags auf GHCR/Docker Hub) im Nutzerbereich bleiben sollten, ist eine
  Auftraggeber-Frage (§6, Frage 2), keine Slice-Arbeit.
- **`AGENTS.md` wird nicht geändert.** Rang 6 nennt `docs/user/` mit „Operations,
  Quality, Releasing“ (`AGENTS.md` Zeile 66) — das ist Regeltext; die Änderung
  braucht die Freigabe des Auftraggebers (§6, Frage 1, mit Textvorschlag).
- **Kein Release, keine Versionsänderung.** Anderer Vorgang.
- **Keine Inhaltsänderung an Accepted-ADRs.** Nur Zitat-Korrekturen nach
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md); §Entscheidung,
  §Konsequenzen und §Verglichene Alternativen bleiben unberührt (`AGENTS.md` §3.5).
- **Keine Änderung an `docs/user/benutzerhandbuch.md`.** Es verweist nicht auf
  `releasing.md` (`git grep -n releasing -- docs/user/benutzerhandbuch.md
  docs/user/benutzerhandbuch-standard.md` am Stand `81d0fa96`: 0 Treffer); das
  Handbuch richtet sich an Betreiber des laufenden Feed-Containers, `releasing.md` an
  Maintainer, ein Link vom Handbuch nach `docs/maintainer/` entsteht nicht.

**Messung des Suchlaufs (`AGENTS.md` §3.13, Stand `81d0fa96`).** Die bewegte Eigenschaft ist
der Pfad `docs/user/releasing.md`. Messung am Plan-Stand: 45 Trefferzeilen in 25 Dateien
für `releasing.md` im aktiven Baum ohne Records (`docs/reviews`, `done/`,
vendored Baseline); darin 24 Zeilen mit dem Klartext-Pfad `docs/user/releasing.md`;
im ganzen Baum ohne Baseline 8 Zeilen in **Link-Form** (`](…releasing.md`), nur diese
brechen `make docs-check` (`links`-Modul): `README.md` 1, `README.de.md` 1,
`docs/plan/adr/0123-…` 3, `docs/reviews/` 3.

```suchlauf
81d0fa96 45 releasing\.md -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
81d0fa96 24 -F docs/user/releasing.md -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
81d0fa96 8 -E \]\([^)]*releasing\.md -- . :!.harness/baseline
81d0fa96 3 -E \]\([^)]*releasing\.md -- docs/reviews
81d0fa96 3 -F 'Operations, Quality, Releasing' -- AGENTS.md harness/README.md docs/plan/adr/0051-cicd-pipeline-github-actions.md
81d0fa96 1 -F excluded=(releasing.md -- tools/harness
diff 0 -F excluded=(releasing.md -- tools/harness
diff 3 -E \]\([^)]*releasing\.md -- docs/reviews
diff 0 -E \]\([^)]*user/releasing\.md -- docs/reviews README.md README.de.md
diff 0 -F 'Operations, Quality, Releasing' -- AGENTS.md harness/README.md
diff 41 releasing\.md -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 17 -F docs/user/releasing.md -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
```

**Messung am Stand `diff` (Implementer, nach den Commits; Parent `de24629d`).**
Von den 22 Klartext-Treffern `docs/user/releasing.md` im aktiven Baum liegen 1 in
`ADR-0051` (Klartext, bleibt), 4 in `ADR-0123` (3 Linktexte mit korrigiertem Ziel,
1 Klartext, bleiben), 1 in `ADR-0143` (Klartext, bleibt), 14 im
Beobachtungs-Register (Klartext, bleiben; die lebende Adresse in
`plattform-verhalten-nur-vom-betreiber-pruefbar/state.md` ist nachgezogen) und 2 in
`docs/plan/planning/open/slice-sdk-meldungscodes-in-fehlertypen.md` (fremde Datei,
nach dem Start dieses Slice angelegt: **gemeldet**, Frist die Closure dieses Slice;
der Planner zieht nach). Die zwei Review-Links, die der Plan als dritten zählte
(`verifikation-slice-release-doku-releasing.md`), sind Dateinamen-Links auf
`review-slice-release-doku-releasing.md`, kein Pfad nach `releasing.md`; die Link-Form
mit Ziel `…/user/releasing.md` hat 0 Treffer im Baum.

**Messung nach dem Nachzug der Closure (Stand `diff`, Parent `5c4474ae`).** Der
Klartext-Treffer `docs/user/releasing.md` im aktiven Baum sinkt von 22 auf 17: die
drei Linktexte in `ADR-0123` (Commit `1c90030f`, Zitat-Korrektur nach
[`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)) und die zwei
Nennungen im SDK-Plan sind nachgezogen; die Suchlauf-Zeile `diff 17` trägt den neuen Wert
(gemessen mit `make suchlauf-nachmessen`).

(Die Zeile mit `-F Operations, …` zählt Rang-6-Zeilen von `AGENTS.md`,
`harness/README.md` und [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) — gemessen am Stand: je 1. Die Plan-Datei ist vom
Werkzeug ausgenommen; gemessen wird mit `make suchlauf-nachmessen PLAN=<Plan-Datei>`.)

**Gruppierung der 25 Dateien, Entscheidung je Gruppe:**

| Gruppe | Dateien | Entscheidung |
|---|---|---|
| Nutzer-Einstieg | `README.md`, `README.de.md` | Link-Ziel auf `docs/maintainer/releasing.md` (Link bricht sonst); Satz bleibt kennungsfrei |
| Harness | `harness/README.md` (Z. 28 Rang 6, Z. 121 Gate-Zeile), `.harness/skills/nutzerdoku-schreiben.md` Z. 13, `harness/sensors/handbuch-public-doc-check.md` (Z. 30, Z. 63), `tools/harness/handbuch-public-doc-check.sh`, `tools/harness/run-handbuch-public-doc-check-tests.sh` | anpassen (Nicht-Regel-Dateien, Pfade/Verweise); Rang-6-Zeile von `harness/README.md` nur zusammen mit `AGENTS.md` (§6 Frage 1) |
| Accepted-ADRs | `0051` (1 Zeile, Klartext), `0123` (4 Zeilen: 3 Links, 1 Klartext), `0143` (11 Zeilen, alle Klartext) | Zitat-Korrektur nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) nur, wo die Fundstelle ein **Link** ist: `0123` Z. 48, 139, 304 (Linkziel `../../user/releasing.md` → `../../maintainer/releasing.md`), eine neue §Geschichte-Zeile. `0051` Z. 212 (Klartext im Entscheidungs-Vergleich: beschreibt das damalige Vorbild) und `0123` Z. 275 (Klartext im Alternativen-Vergleich) und alle `0143`-Zeilen (Klartext, beschreiben die damalige Lage; Festlegung 2 „dauerhaft außerhalb“ bleibt der Wortlaut) bleiben stehen — die Entscheidung der ADR ändert sich nicht, ihre Lage räumlich bestätigt sich (Klarstellung: `docs/maintainer/` liegt außerhalb von `docs/user/`, das Gate prüft es dort nie). Für `0051` und `0143` entsteht keine §Geschichte-Zeile |
| Register | 12 Dateien unter `docs/plan/planning/observations/BEO-PGC/**` (`evidence/`, `state.md`, `observation.md`), alle Klartext, 0 Links (Messung: 0 Link-Form-Treffer unter `observations/`) | unverändert belassen (Records bzw. Klartext-Pfade, kein Linkbruch; `observation.md`-Zitate beschreiben die damalige Lage). Der Implementer prüft, dass `state.md`-Adressen (`plattform-verhalten-nur-vom-betreiber-pruefbar`: „Adresse: `docs/user/releasing.md`“) eine **Adresse** tragen: diese eine Zeile wird auf den neuen Pfad nachgezogen, weil sie ein lebender Zeiger ist, kein Record — die übrigen bleiben |
| Records | `docs/reviews/**` (3 Link-Fundstellen), `done/**` (0 Link-Fundstellen) | nur das Linkziel der 3 Review-Links korrigieren (wie bei früheren Verschiebungen, Commit nennt [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)); Klartext-Pfade bleiben |

## 2. Definition of Done

- [x] **(1) Verschiebung.** `docs/maintainer/` existiert; `docs/user/releasing.md` ist per
      reinem `git mv` (eigener Commit, `AGENTS.md` §3.3) nach `docs/maintainer/releasing.md`
      verschoben. Zu belegen durch `git diff --stat -M <Parent>..<Move-Commit>` (Similarity
      100 %, kein Inhalt) und `git log --follow`.
- [x] **(2) Verweise nachgezogen** (eigene Commits, getrennt vom Move): die 8 Link-Form-Stellen
      (README.md, README.de.md, [`ADR-0123`](../../adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) ×3 per Zitat-Korrektur, `docs/reviews/` ×3 nur Linkziel),
      `releasing.md` selbst (Link auf `version.md` → `../user/version.md`, `docs/user/version.md`-Klartext
      bleibt wahr, Zweck/Zielgruppe §1 geschärft: Maintainer als Zielgruppe, Betreiber nur
      als Leser des Mechanismus; Version/Stand und Änderungshistorie-Zeile fortgeschrieben), die
      Harness-Träger der Gruppe „Harness“ (ohne Rang 6, sofern die Freigabe aussteht). Zu
      belegen durch `make docs-check` Exit 0 und den Suchlauf aus §1 am Stand `diff`: `docs/user/releasing.md`
      im aktiven Baum ohne Records erwartet 0 Trefferzeilen außer den zitierenden
      Klartext-Stellen in `docs/plan/adr/0051-…` (1) und `docs/plan/adr/0143-…`/`0123-…`
      (Klartext, nach Entscheidung je Datei oben) — die genaue Restzahl misst der Implementer
      und trägt sie als gemessen ein.
- [x] **(3) Handbuch-Gate bereinigt.** `excluded=(releasing.md …)` in
      `tools/harness/handbuch-public-doc-check.sh` führt `releasing.md` nicht mehr (die Datei
      liegt nicht mehr in `docs/user/`; eine genannte, fehlende Datei endet mit Exit 2);
      Tabellentest `tools/harness/run-handbuch-public-doc-check-tests.sh` (`all_files`,
      Fall „Kennung in releasing.md ausgenommen“, Löschfall) und Sensor-Vertrag
      `harness/sensors/handbuch-public-doc-check.md` (Reichweite-Tabelle, Absatz 3) sowie die
      Gate-Zeile in `harness/README.md` (Ausnahme-Aufzählung) sind angepasst. Zu belegen
      durch `make handbuch-public-doc-check` Exit 0 und `make test-handbuch-public-doc-check`
      Exit 0. Gate-Verhalten bleibt unverändert (kein Lockern, `AGENTS.md` §3.6): geprüft
      sind weiter dieselben drei Dateien, ausgenommen die vier Erzeugnisse.
- [x] `make gates` grün (Exit direkt geprüft, `AGENTS.md` §3.9); `make docs-check` Exit 0;
      `make doc-immutable RANGE=<Parent>..HEAD` 0 Befunde;
      `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0.
- [x] Zitat-Korrektur nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md):
      `docs/plan/adr/0123-…` trägt genau **eine** neue §Geschichte-Zeile (Datum, Ereignis,
      Commit-Kennungen; in der Closure um die Linktexte erweitert, Commit `1c90030f`); die Commit-Message nennt `ADR-0073`; `0051` und `0143` tragen keine
      neue Zeile (keine Fundstelle geändert).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      Report: `docs/reviews/review-slice-maintainer-ordner-releasing-verschieben.md` (0 HIGH, 0 MEDIUM, keine Fixrunde).
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `harness/README.md` (Gate-Zeile `make handbuch-public-doc-check`; Rang 6 nach
      Auftraggeber-Antwort), `README.md`/`README.de.md` (Link).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Kandidat: der Suchlauf trennt Link-Form (bricht
      `docs-check`) von Klartext (bleibt Record); der Lerneintrag nennt, was die Messung trug.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben, oder „keine Beobachtung angefallen“
      in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — die Roadmap
      führt keine offene Welle, deshalb trägt die Slice-Closure selbst die drei Paarungen
      (Modul 6 §Was der wellenlose Betrieb selbst auslöst).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/maintainer/releasing.md` (von `docs/user/releasing.md`) | `git mv` (Commit 1, rein) | Auftrag; Ordner entsteht durch den Move |
| `docs/maintainer/releasing.md` | update (Commit 2) | Link `version.md` → `../user/version.md`; §1 Zielgruppe; Versionszeile und Änderungshistorie (Tiefe der übrigen relativen Links unverändert: beide Ordner liegen eine Ebene unter `docs/`; Auflösung der relativen Link-Ziele ist durch `make docs-check` zu belegen) |
| `docs/maintainer/README.md` | neu (klein) | Index: Zweck („Dokumente für Maintainer des Projekts, nicht für Betreiber des Feed-Containers“), Zielgruppe, eine Zeile je Dokument (derzeit `releasing.md`); kein Regeltext, kein Kennungs-Zwang. Entscheidung: ja, weil ein einzelnes Dokument in einem neuen Ordner sonst keine Zielgruppenaussage trägt |
| `README.md`, `README.de.md` | update | Link-Ziel Zeile 38 |
| `docs/plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md` | Zitat-Korrektur | 3 Linkziele + eine §Geschichte-Zeile |
| `docs/reviews/architect-verdict-sdk-kotlin-cloudsmith.md`, `…verifikation-slice-meldungscodes-http-grpc-fehlerkoerper.md`, `…verifikation-slice-release-doku-releasing.md` | Zitat-Korrektur (nur Linkziel) | `docs/reviews/**` Records; Klartext bleibt |
| `docs/plan/planning/observations/BEO-PGC/plattform-verhalten-nur-vom-betreiber-pruefbar/state.md` | update (1 Zeile) | lebende Adresse |
| `tools/harness/handbuch-public-doc-check.sh`, `tools/harness/run-handbuch-public-doc-check-tests.sh`, `harness/sensors/handbuch-public-doc-check.md`, `harness/README.md` Z. 121 | update | Gruppe (3) der DoD |
| `.harness/skills/nutzerdoku-schreiben.md` Z. 13 | update | „Nicht Gegenstand: `releasing.md` (Maintainer-Doku)“ → nennt den neuen Ort `docs/maintainer/releasing.md` bzw. entfällt, weil die Datei nicht mehr unter `docs/user/` liegt (Skill gilt für `docs/user/`); der Implementer wählt die kürzere wahre Form |
| `harness/README.md` Z. 28, `AGENTS.md` Z. 66 | **nach Freigabe** | Rang 6, §6 Frage 1 |

**`.d-check.yml` (geprüft):** `scan.roots: ["."]` — `docs/maintainer/` wird ohne Änderung gescannt
(`links`, `ids`, `hostpaths`, `tracked`); `matrix`-Klassen führen `docs/user/` nicht
(`git grep -n "docs/user" .d-check.yml` am Stand `81d0fa96`: nur die vier `*-abdeckung.md`-Pfade
in `trace.coverage` und einer Regel Zeile 256), also gilt keine Regel, die für
`releasing.md` in `docs/user/` galt, anders in `docs/maintainer/`. Keine `.d-check.yml`-Änderung
erwartet; ein neuer Root ist nicht nötig. Zu belegen durch `make docs-check` am Ende.

**Release-Werkzeuge (geprüft):** `git grep -n "releasing" -- .github tools Makefile harness/mk`
am Stand `81d0fa96`: Treffer in `tools/harness/handbuch-public-doc-check.sh` und im Tabellentest
(4 Zeilen, siehe §1-Suchlauf) — kein Workflow und kein Make-Ziel liest `releasing.md`
(Messung der Befehlsform ohne `.md`: der Implementer wiederholt sie am Stand des Starts); `release.yml` liest `docs/user/version.md`
(bleibt).

**Handbuch-Frage (Antwort):** `releasing.md` ist für Handbuch-Leser (Betreiber) nicht
notwendig; das Handbuch verlinkt es nicht, und das Handbuch-Gate verbietet Links nach
`docs/plan/` und `docs/reviews/`, nicht nach `docs/maintainer/` — es entsteht kein Link und
damit kein Fall. Die Frage, ob das Gate künftig auch Links nach `docs/maintainer/` verbietet, ist
nicht Gegenstand (Gate-Verhalten bleibt, §3.6).

## 4. Trigger

**Start** (`next` → `in-progress`): Auftraggeber hat den Slice zur Umsetzung freigegeben;
die Antwort auf §6 Frage 1 (Rang-6-Wortlaut) liegt vor oder ist bewusst auf „später“ gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Fundstellen-Menge wächst über die
  25 Dateien hinaus, weil der Suchlauf am Stand des Starts mehr Link-Form-Treffer findet als
  die 8 des Plans.
- `in-progress` → `open` (blockiert): der Auftraggeber beantwortet §6 Frage 2 mit „aufteilen“;
  dann wird die Aufteilung ein eigener Slice, dieser wartet.

## 5. Closure-Trigger

DoD vollständig, `make gates` Exit 0 (direkt geprüft), `make docs-check` Exit 0,
`make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, Review- und Verifikations-Report
unter `docs/reviews/`, Closure-Notiz geschrieben; der Lerneintrag nennt die geschärfte
Regel, den neuen Sensor oder die benannte Spec-Lücke.

## 6. Risiken und offene Punkte

**Auftraggeber-Fragen** (Regeltext-Berührung; der Planner ändert `AGENTS.md` nicht):

1. **Rang-6-Wortlaut.** `AGENTS.md` Zeile 66: „6. `docs/user/` — Operations, Quality, Releasing.“ Nach dem Move
   stimmt „Releasing“ nicht mehr für `docs/user/`. Textvorschlag:
   „6. `docs/user/` — Operations, Quality, Version; `docs/maintainer/` — Releasing.“ (in
   `AGENTS.md` und `harness/README.md` in der dort üblichen Link-Form der Rang-Zeile)
   Dieselbe Zeile in `harness/README.md` Zeile 28 (Nicht-Regel-Datei, im Slice erlaubt) folgt
   dem Wortlaut der Freigabe, damit beide Dateien gleich lauten.
2. **Aufteilung.** Sollen Betreiber-nahe Teile von `releasing.md` (§2 Versionierung,
   Image-Tags und Herkunft eines Images auf GHCR/Docker Hub) im Nutzerbereich bleiben oder ins
   Handbuch wandern? Der Slice verschiebt das Dokument ungeteilt.

Risiken:

- **Externe Rückwärts-Links.** GitHub kennt keine Weiterleitung für verschobene Dateien;
  externe Links auf `docs/user/releasing.md` brechen. — **Ausgang:** weiter offen (bewusst
  in Kauf genommen: kein Beleg für externe Links im Repo; Auftraggeber kennt den Preis des Moves).
- **Docker-Hub-Beschreibung.** `hub-description.yml` spiegelt `README.md`; der geänderte Link
  erscheint dort erst beim nächsten Release oder `workflow_dispatch`. — **Ausgang:** weiter
  offen bis zum nächsten Lauf (kein Fehler, nur Verzögerung; der Workflow ist nicht geändert,
  `AGENTS.md` §3.10 greift nicht).
- **Klartext-Pfade in Records und Accepted-ADRs** zeigen auf einen nicht mehr existierenden Pfad. —
  **Ausgang:** weiter offen als Nebenwirkung der Immutabilität (Records beschreiben die damalige Lage;
  `docs-check` prüft Klartext nicht, nur Links).
- **Gate-Ausnahme räumlich bestätigt, nicht entschieden.** [`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Festlegung 2 („dauerhaft
  außerhalb“) wird durch den Move bestätigt, nicht geändert; ein Reviewer könnte die Entfernung
  des Listen-Eintrags als Lockerung lesen. — **Ausgang:** entfallen: das Gate prüft dieselben Dateien
  wie zuvor; ein Eintrag für eine nicht mehr vorhandene Datei würde Exit 2 erzeugen
  (`excluded`-Eintrag ohne Datei), der Eintrag muss weg.
- **Reihenfolge-Fehler bei gleichzeitigem Move und Gate.** Zwischen Move-Commit und
  Gate-Anpassung ist `make handbuch-public-doc-check` rot (Exit 2, genannte Datei fehlt). —
  **Ausgang:** entfallen durch Disziplin: Gate-Anpassung im Commit direkt nach dem Move; nur der
  Endstand wird gegated. Die `commit-traceability`-Range prüft Messages, nicht Gate-Farbe je Commit.

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
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** Der reine Move (Similarity 100 %, eigener Commit), das
  Handbuch-Gate ohne Lockerung (Verifier-Mutationen 1 bis 4: eine fehlende genannte Datei
  endet Exit 2, eine Kennung in `docs/maintainer/` bleibt außerhalb) und der Suchlauf mit
  Nachmessen an beiden Ständen; er fand den Link-Bruch (`docs-check`) und trennte ihn vom
  Klartext, der als Record stehen bleibt.
- **Was ging anders als geplant:** Der Plan zählte 3 Review-Links, tatsächlich waren es
  2 Linkziel-Änderungen: die dritte Fundstelle ist ein Dateiname-Link auf
  `review-slice-release-doku-releasing.md`, kein Pfad nach `releasing.md`. Die
  Zitat-Korrektur in `ADR-0123` brauchte einen zweiten Schritt: der Review (F-1) und der
  Verifier (V-1, V-4) fanden, dass der sichtbare Linktext das Linkziel nicht begleitet hatte
  und die §Geschichte-Zeile den Move statt der korrigierenden Commits nannte (Commits
  `1c90030f`, `df503ad5`). Der Kopfkommentar des Gate-Skripts (V-2) war nach der
  Listenänderung veraltet und ist im Commit `5c4474ae` bereinigt. Die zwei Klartext-Nennungen
  im fremden SDK-Plan (gemeldet, Frist die Closure) sind in `1c90030f` nachgezogen.
- **Steering-Loop-Eintrag:** Lehre, nicht verkörpert: ein Suchlauf über einen bewegten Pfad
  zählt **Link-Form** (bricht `docs-check`) und **Klartext** (bleibt Record oder wird als
  lebende Adresse nachgezogen) getrennt; wer die Fundstellen im Plan zählt, misst beide
  Formen und die Linktexte neben den Linkzielen, sonst stimmt die Zahl der Korrekturen
  nicht (hier 3 geplant, 2 tatsächlich) und der sichtbare Text bleibt stehen. Die Messung
  trug das Ergebnis; ein neuer Sensor entsteht daraus nicht (Grenze der Suchform:
  `AGENTS.md` §3.13). Nicht verkörpert, nur gezählt.
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-maintainer-ordner-releasing-verschieben.md` in
  `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform/` ergänzt — Zähler steht damit bei 2×.
  `BEO-PGC/zitat-nennt-die-falsche-stelle`: kein Vorkommen (das Linkziel nannte die richtige
  Datei, nur der sichtbare Text blieb stehen). „Sichtbarer Linktext gegen Linkziel" ist als
  eigene Klasse ein erstes Vorkommen und hier vermerkt, kein eigener Eintrag. Die lebende
  Adresse in `BEO-PGC/plattform-verhalten-nur-vom-betreiber-pruefbar/state.md` ist
  nachgezogen, Zähler unverändert.
- **Folge-Slices:** keiner angelegt. Benannt, nicht angelegt: eine Folge-ADR des Architect,
  die die Abschnitte-Liste von
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
  an die Kurzform angleicht (Auslöser `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`,
  2×). Der SDK-Plan
  [`slice-sdk-meldungscodes-in-fehlertypen`](../in-progress/slice-sdk-meldungscodes-in-fehlertypen.md)
  bleibt in `open/`.
- **Validator:** nicht nötig. Der Slice ist Pflegearbeit an der Ablage: er ändert keine
  Spec-Stelle (`git grep -n releasing spec` zählt 0 Treffer), keine Anforderung und kein
  Verhalten des Produkts; das Gate behält seinen Gegenstand (keine Lockerung,
  `AGENTS.md` §3.6). Die Abnahme gegen Anforderungen läuft über keine `LH-*`-Kennung, die
  Review und Verifikation tragen die Prüfung.
- **Risiken aus §6:** Externe Rückwärts-Links: weiter offen. Docker-Hub-Beschreibung: weiter
  offen bis zum nächsten Release oder `workflow_dispatch`. Klartext-Pfade in Records und
  `Accepted` ADRs: weiter offen (Nebenwirkung der Immutabilität). Gate-Ausnahme räumlich
  bestätigt: entfallen. Reihenfolge Move/Gate: entfallen. Die Auftraggeber-Fragen 1 und 2
  sind beantwortet (Rang-6-Wortlaut in `AGENTS.md` Z. 66 und `harness/README.md`
  übernommen, `releasing.md` ungeteilt verschoben).
- **Drei Paarungen:** Anker `docs/maintainer/releasing.md` (Datei, aufgelöst durch
  `make docs-check`) · Folge-Slice: keiner, die Folge-ADR ist beim Architect benannt ·
  Register: `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` bei 2×.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist allein die Default-Sub-Area `*`
(Dokumentation, Harness-Skripte des Handbuch-Gates); `harness/conventions.md`
§Modus-Deklaration führt nur diese eine Zeile, keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer:
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (18 Dateien unter `evidence/`) und
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (34 Dateien) — beide beschreiben die Klasse
„Arbeit bewegt eine beschriebene Eigenschaft, Träger bleiben stehen“; dieser Slice trägt sie
als Suchlauf (§1) und als Träger-Gruppierung. `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(1 Datei) betrifft `releasing.md` selbst, nicht den Pfad. Keine Beobachtung erreicht
mit diesem Slice eine neue Stufe, die einen eigenen Slice verlangte.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (`*`, Modus Greenfield nach `harness/conventions.md`
§Modus-Deklaration) — kein Modus-Begründungsblock nötig.
