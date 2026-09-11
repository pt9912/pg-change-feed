# Beleg: slice-016

Vorgang: slice-016 (d-migrate-1.3.1-Views-Retirement).

Fund: Implementer, Reviewer und Verifier haben unabhängig voneinander
real reproduziert, dass ein zweiter `make schema-rollout`-Lauf gegen eine
bereits migrierte DB mit Exit 8 (`DropView` auf `cdc.heartbeat`/
`cdc.metrics`) blockiert — pin-unabhängig (Reviewer: reproduziert auch
gegen den alten Pin 1.3.0 in einem separaten Git-Worktree). Kein
bestehender Sensor übt diesen Pfad aus (alle Testläufe fahren die
Umgebung vorher frisch hoch). Bewusst nicht in slice-016 behoben —
anderer Vorgang (§1-Klasse 3): beträfe `nacharbeit-observability.sql`/
`nacharbeit-heartbeat.sql`, nicht die Views-Ausweichform dieses Slices.

Quelle: `docs/reviews/review-slice-016.md`, `docs/reviews/verify-slice-016.md`.
