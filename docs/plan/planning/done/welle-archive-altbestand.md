# Welle welle-archive-altbestand: Erste Archivierung dieses Repos — wellenloser Altbestand und `welle-d-check`

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-archive-altbestand-results.md`). Der Zustand ist
die Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**Dieses Repo hat noch nie archiviert** (`find . -iname "*archiv*"` findet nur
die zwei vendored Templates unter `.harness/baseline/`). Zwei Träger stehen
dafür bereit: das externe Werkzeug `ai-harness-init archive-welle`
(Hausform-Zitation, `AGENTS.md` §3.11 — vendored unter
`.harness/state/bin/ai-harness-init`, **nicht** ins Repo kopiert, läuft
gegen den Repo-Baum wie ein Fremd-Programm) und der bereits geschlossene
`welle-d-check` (2 Mitglieder-Slices, Plan und Ergebnisnotiz liegen bereits
in `done/`).

**Real gemessen (Dry-Run, `--vorschau welle-d-check`, 2026-09-18, sauberer
Baum):** Ein Lauf gegen `welle-d-check` sammelt **nicht nur** dessen 2
Mitglieder ein, sondern **jeden** wellenlosen `done/`-Slice dieses Repos
zugleich — 41 an der Zahl — plus 64 fremde (bleiben liegen) und 99
Review-Reports. Der Grund ist strukturell, nicht ein Fehler der Kennung: Das
Werkzeug sammelt einen wellenlosen Slice unabhängig von der übergebenen
Welle-ID ein, solange **keine** Untergrenze existiert
(`internal/archive/scan.go`s `untergrenze()` — ein vorhandenes
`done/*/archiv.zip`, irgendeine Kennung). Da dieses Repo noch nie archiviert
hat, gibt es diese Untergrenze nicht — jeder erste Lauf, gleich welcher
Kennung, würde denselben Altbestand mitreißen. Zwei Sperren stehen dem
schreibenden Lauf heute im Weg: `[untergrenze]` (kein
`done/*/archiv.zip` existiert) und `[haenger]` (mehrere `Accepted`-ADRs
zitieren Review-Reports, die der Lauf ersatzlos löschen würde).

**Das Mehr gegenüber einem Einzel-Slice:** Die Zuordnungsfrage — welcher
Schlüssel den Altbestand einsammelt, und in welchem Verhältnis das zu
`welle-d-check`s eigener Archivierung steht — ist eine architektonische
Entscheidung mit dauerhafter Konsequenz (der gewählte Schlüssel ist
ortsfest, siehe Präzedenzfall unten), **und** ihr Vollzug hängt von einer
externen, bereits separat laufenden Bereinigung ab (`[haenger]`). Das ist
mehr als eine einzelne DoD trägt — Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht.

**Präzedenzfall (Hausform, kein host-lokaler Pfad, `AGENTS.md` §3.11):**
`ai-harness-init`s eigenes Repo stand vor genau dieser Frage und hat sie mit
einer eigenen ADR entschieden (Sammel-Archiv unter einem eigenen, dedizierten
Schlüssel statt Vermischung mit einer echten Welle) — dieselbe Klasse
Entscheidung, hier für dieses Repo neu zu treffen, nicht zu kopieren.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln.

- WIP-Limit frei (`in-progress/` trägt aktuell keinen Slice).
- Keine externe Vorbedingung für die **Eröffnung** — die ADR-Entscheidung
  (slice-archive-altbestand-adr) ist der erste Slice dieser Welle, nicht ihr
  Start-Trigger.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

- Beide Slices in `done/`.
- Ein `done/*/archiv.zip` existiert (Untergrenze gesetzt) — Beleg:
  `ls -d docs/plan/planning/done/*/`.
- `welle-d-check` ist tatsächlich archiviert (nicht nur geschlossen) —
  Beleg: sein Plan und seine Slices tragen einen Stub, kein Volltext mehr.
- `make gates` grün auf dem Endstand.
- Closure-Notiz `welle-archive-altbestand-results.md` geschrieben.

## 4. Slices in dieser Welle

| Slice | Titel | Bezug |
|---|---|---|
| slice-archive-altbestand-adr | Eigene ADR: welcher Schlüssel sammelt den wellenlosen Altbestand ein | Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 4 |
| slice-archive-altbestand-vollzug | Realer `archive-welle`-Lauf (Altbestand und `welle-d-check`) nach ADR-Ausgang | slice-archive-altbestand-adr |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Wird blockiert von:** einem separaten, bereits laufenden Vorgang außerhalb
  dieser Welle — einer ADR-Zitat-Korrektur-Runde (`AGENTS.md` §3.5, Ausnahme
  nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md))
  an den `Accepted`-ADRs, die heute Review-Reports referenzieren, welche ein
  Altbestands-Archivlauf ersatzlos löschen würde (`[haenger]`-Sperre, real
  gemessen: u. a. `ADR-0062`, `ADR-0069`, `ADR-0070`, `ADR-0075`, `ADR-0083`,
  `ADR-0084`, `ADR-0086`, `ADR-0087`, `ADR-0089`, `ADR-0091`, `ADR-0092`,
  `ADR-0093` — nicht abschließend, der reale Lauf zum Vollzugs-Zeitpunkt
  entscheidet). **Kennung dieses externen Vorgangs:** wird hier nachgetragen,
  sobald sie vorliegt — bis dahin blockiert `slice-archive-altbestand-vollzug`
  ohne benannte Kennung, aber mit benannter Bedingung
  (`[haenger]`-Sperre entfällt am realen Baum).
- Blockiert nichts Bestehendes — kein anderer offener Slice oder Welle
  hängt an dieser.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1.

- **Die `[haenger]`-Bereinigung selbst** (Zitat-Korrektur an den betroffenen
  `Accepted`-ADRs) — läuft als eigener, bereits beauftragter Vorgang außerhalb
  dieser Welle (siehe §5). Diese Welle **konsumiert** sein Ergebnis
  (Vorbedingung), liefert es nicht.
- **Eine Änderung an `ai-harness-init` selbst** (z. B. Nachrüsten der
  `altbestand`-Sperren-Aufhebung, falls die hier vendorte Werkzeug-Version sie
  noch nicht trägt — real geprüft: sie hebt für den Schlüssel `altbestand`
  heute `ergebnisnotiz`/`kein-plan` **nicht** auf, anders als die neuere
  Fassung im Schwester-Repo). Bestand bleibt bewusst stehen — es ist ein
  externes Werkzeug, keine Baustelle dieses Repos; die ADR (slice 1) wägt
  einen Umgang **ohne** Werkzeugänderung ab.
- **Der `structure`-Modul-Blindfleck nach dem Move** (`done/slice-*.md` ist
  ein flacher Glob, kein `**` — ein archivierter Slice unter
  `done/<schlüssel>/` verlässt den Prüfbereich der Closure-Notiz-Regel). Wird
  in der ADR **benannt** (Folgepflicht, analog zum Präzedenzfall), aber nicht
  zwingend in dieser Welle geheilt — die ADR entscheidet, ob das ein
  Blocker oder ein akzeptiertes Negativ ist.
- **Eine zweite Archivierungs-Runde für künftig schließende Wellen.** Sobald
  die Untergrenze aus dieser Welle steht, trägt die laufende Baseline-Regel
  ohne weitere Zuordnung (kein Wiederholungsfall).

## 7. Closure-Notiz

<!--
BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md §Verwendung,
Schritt 5) und darf deshalb nichts Tragendes halten.
-->

Ergebnis: `welle-archive-altbestand-results.md`, Geschwister im Ruheort `done/`.
Zähler: `../observations/README.md`, eine Ebene über dem Ruheort.
