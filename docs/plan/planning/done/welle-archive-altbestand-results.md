# Welle welle-archive-altbestand — Erste Archivierung dieses Repos — Closure-Notiz

**Welle:** welle-archive-altbestand
**Abschluss:** 2026-09-18
**Verantwortlich:** pt9912

## Was wurde geliefert?

- Erstmalige Untergrenze für „wellenlos seit der letzten Closure" gesetzt:
  `docs/plan/planning/done/altbestand/archiv.zip` (41 wellenlose Slices +
  Plan-Äquivalent, 99 Review-Reports entfernt) und
  `docs/plan/planning/done/welle-d-check/archiv.zip` (2 Mitglieder-Slices +
  Plan, 0 Review-Reports entfernt — alle sechs Kandidaten anderswo weiter
  zitiert).
- `ADR-0096` entscheidet den dedizierten Schlüssel `altbestand` für den
  Altbestand, getrennt von `welle-d-check`s eigener Archivierung —
  bestätigt durch den Präzedenzfall des Schwester-Repos, der exakt
  denselben Schlüsselnamen für dieselbe Funktion trägt.
- `make docs-check`: 0 Befunde auf dem Endstand (819 Dateien), vollständiger
  Verweis-Nachzug bestätigt.

## Was hat funktioniert?

Die vorgelagerte ADR-Zitat-Korrektur (`ADR-0073`/`ADR-0094`/`ADR-0095`/
`ADR-0097`, außerhalb dieser Welle beauftragt) hat die `[haenger]`-Sperre
vollständig gelöst, bevor der reale Lauf versucht wurde — kein einziger
Hänger-Fund während des schreibenden Laufs selbst. Ein frischer Build des
vendorten Werkzeugs aus der aktuellen Schwester-Repo-Quelle trug bereits
die Sonderbehandlung für den Schlüssel `altbestand`, was `ADR-0096`s
Entscheidung nachträglich bestätigte, ohne sie revidieren zu müssen.

## Was ging anders als geplant?

Der lokale, optionale `commit-msg`-Hook lehnte die fest einprogrammierten
Commit-Messages des externen Werkzeugs ab (keine `LH-*`/`ADR-*`-Kennung).
Nutzer-Entscheidung: Bypass über eine prozess-scoped
`GIT_CONFIG_*`-Umgebungsvariable statt einer persistenten Config-Änderung.
Das reale Standing-Gate (`make commit-traceability`) markierte die vier
resultierenden Commits vorübergehend im gleitenden Fenster — durch die
nachfolgenden, Kennung-tragenden Closure-Commits von selbst aufgelöst, ohne
Sonderbehandlung.

## Steering-Loop-Einträge

Keine — beide während dieser Welle entstandenen Beobachtungen
(`BEO-PGC/externes-werkzeug-committet-ohne-kennung`) liegen bei 1×, unter
der 3×-Schwelle. Kein bestehender ≥3×-Eintrag stand beim Lese-Schritt
dieser Closure noch offen (alle bereits in früheren Wellen verkörpert).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`docs/plan/planning/observations/`](../observations/).
Neu während dieser Welle: `BEO-PGC/externes-werkzeug-committet-ohne-kennung`
(1×).

## Folge-Slices

Keine.

## Verifikation

- `docs/reviews/review-slice-archive-altbestand-adr.md` (0 HIGH).
- `docs/reviews/review-slice-archive-altbestand-vollzug.md` (0 HIGH, 2
  MEDIUM — behoben, siehe Folge-Commit).
- `make gates`: grün auf dem Endstand (alle sechs Gates, `commit-traceability`
  eingeschlossen).
- Archiv-Integrität unabhängig geprüft (Reviewer): beide `archiv.zip`
  CRC-intakt, Stichproben-Volltext korrekt.
