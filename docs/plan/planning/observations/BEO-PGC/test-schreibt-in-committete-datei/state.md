Zustand: offen — Ausgang: **weiter offen** (unter der 3×-Schärfungsschwelle).
Der Guard-Test stellt `plan.yaml` und `down.sql` selbst wieder her; für die Läufe
von `make test-store` und `make test-replication` (`tools/schema/apply-rollout.sh`
ruft `make schema-rollout` ohne eigenes Ziel für den Report) besteht keine
Wiederherstellung und kein Slice, der sie liefert. Gelesen wird der Eintrag im
Sichtungs-Schritt der Slice-Planung. Zähler (abgeleitet): 2×
(evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-snapshot-reader.md); der zweite Beleg ist derselbe Pfad
an den zwei DB-Tier-Läufen (Verifikation V-6).
