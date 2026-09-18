# Sammelpunkt `altbestand` — Ergebnisnotiz

**Schlüssel:** `altbestand`
**Abschluss:** 2026-09-18
**Verantwortlich:** pt9912

## Was wurde gesammelt?

Der erste `ai-harness-init archive-welle`-Lauf dieses Repos — vor ihm
existierte kein einziges `done/*/archiv.zip`, die Untergrenze für
„wellenlos seit der letzten Closure" war deshalb unbeobachtbar
([`ADR-0096`](../../adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md)).
Dieser Lauf sammelt bewusst den **gesamten** wellenlosen Altbestand seit
Repo-Beginn ein (41 Slices, real gezählt 2026-09-18 — `grep -l
'^\*\*Welle:\*\* ohne Welle' docs/plan/planning/done/*.md`), damit jede
künftige Wellen-Closure eine beobachtbare Untergrenze vorfindet.

## Was hat funktioniert?

Die vorbereitende ADR-Zitat-Korrektur (`ADR-0073`/`ADR-0094`/`ADR-0095`/`ADR-0097`)
hat die `[haenger]`-Sperre vollständig aufgelöst, bevor dieser Lauf
versucht wurde — kein Fund während des schreibenden Laufs selbst.

## Was ging anders als geplant?

Der ursprünglich vendorte `ai-harness-init`-Stand kannte die
Sonderbehandlung für den Schlüssel `altbestand` (die eigene ADR des
Schwester-Repos zur Betriebsart „Schlüssel ohne Welle" in `archive-welle`
— hebt `[ergebnisnotiz]`/`[kein-plan]`/`[untergrenze]` für genau diesen
Schlüssel auf) noch nicht; ein frischer
Build aus der aktuellen Quelle des Schwester-Repos (`.harness/state/bin/`,
lokal, nicht versioniert) hat das behoben, ohne dass `slice-archive-altbestand-adr`s
Entscheidung (dedizierter Schlüssel) revidiert werden musste — im
Gegenteil, die Sonderbehandlung bestätigt exakt diese Wahl.

## Folge-Slices

Keine.

## Verifikation

- `ai-harness-init archive-welle --vorschau altbestand`: 0 Sperren vor dem
  schreibenden Lauf.
- `make gates` grün vor dem schreibenden Lauf (Voraussetzung für einen
  sauberen Arbeitsbaum).
