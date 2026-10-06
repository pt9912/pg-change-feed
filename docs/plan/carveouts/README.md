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
  Dieses Repo führt Wellen: Die Welle-Closure liest dabei auch die seit der
  letzten Welle geschlossenen Slices ohne Wellen-Zugehörigkeit
  (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht,
  „Achse zuerst“). Führt die Roadmap unter *Offene Wellen* keine Welle, läuft
  der Trigger-Audit zusätzlich bei der Closure eines wellenlosen Slice — so
  geschehen bei der Slice-Closure, aus der
  [`ADR-0070`](../adr/0070-supersede-reichweite-und-klassengrenze.md) hervorging.
- Siehe Baseline-Regelwerk `modul-07-carveouts.md`.
