# Carveouts — PG Change Feed

**Derivativ:** Quelle der Wahrheit sind die Carveout-Dateien
`CO-<NNN>-<titel>.md` in diesem Verzeichnis (Kennungsform: `harness/conventions.md`
MR-000); bei jedem neuen oder aufgelösten Carveout zieht dieser Index mit.

Aktive Carveouts mit Auflösungs-Trigger. Aufgelöste Carveouts wandern nach
`done/` (reiner `git mv`).

## Aktive Carveouts

(keine)

Ein neuer Carveout bekommt hier eine Zeile mit den Spalten `ID` (Link auf die
Datei), `Titel`, `Gate` (das betroffene Make-Ziel), `Trigger` und
`Folge-Slice` (`slice-<Kennung>`); die Bindung-Spalte des Gates in
`harness/README.md` §Sensors verweist auf dieselbe `CO-<NNN>`-Kennung.

## Aufgelöste Carveouts

(noch keine)

## Konventionen

- Jeder aktive Carveout braucht: Trigger, Folge-Slice, letzten Prüf-Termin.
- Bei Welle-Closure: Carveout-Audit zwingend — welche gültig, welche aufgelöst?
  Ohne Welle trägt die Slice-Closure diesen Trigger-Audit (Baseline-Regelwerk
  `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).
- Siehe Baseline-Regelwerk `modul-07-carveouts.md`.
