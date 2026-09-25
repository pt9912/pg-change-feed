Zustand: **gestrichen** — Ausgang: **gestrichen** (akzeptiertes Negativ) · seit
welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3 und §5 (f)).

Begründung: jede Schicht des Setzpfads ist an ihrer Eingabeseite gebunden
(Mutationen rot, View-Test gegen reale PostgreSQL); ein Ende-zu-Ende-Lauf bräuchte
eine Tabelle über der Richtgröße (`estimatedRowsGuideline = 4_000_000` in
`internal/application/usecase/backfill/warn.go`) oder einen Konfigurations-Hebel, den
`ADR-0113` ausschließt („Toleranz oder Richtgröße sollen konfigurierbar sein:
Folge-ADR“). Kein Betrieb existiert (kein Server-Tag trägt Backfill-Änderungen). Die
Grenze „Setz-Pfade nur durch Unit- und View-Test belegt“ steht im Verifikations-Report
und in der Closure-Notiz von `slice-backfill-bench-richtgroesse`.

Zähler: 1× (Datei unter `evidence/`).
