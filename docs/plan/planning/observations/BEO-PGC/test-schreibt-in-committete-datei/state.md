Zustand: offen — Ausgang: **weiter offen**, **Schwelle 3× erreicht mit
`slice-backfill-run-store`**; der Ausgang gehört dem Lese-Schritt der Closure von
`welle-backfill-bestand` (Modul 6), nicht entschieden.
Der Guard-Test stellt `plan.yaml` und `down.sql` selbst wieder her; für die Läufe
von `make test-store` und `make test-replication` (`tools/schema/apply-rollout.sh`
ruft `make schema-rollout` ohne eigenes Ziel für den Report) besteht keine
Wiederherstellung und kein Slice, der sie liefert. Gelesen wird der Eintrag im
Sichtungs-Schritt der Slice-Planung. Zähler (abgeleitet): **4×**
(evidence/slice-backfill-change-origin.md,
evidence/slice-backfill-snapshot-reader.md,
evidence/slice-backfill-run-store.md,
evidence/slice-backfill-sql-administration.md); der vierte Beleg
(`slice-backfill-sql-administration`, Verifikation V-5) ist derselbe Pfad in den
Läufen zweier weiterer Rollen; der zweite Beleg ist derselbe Pfad
an den zwei DB-Tier-Läufen (Verifikation V-6), der dritte derselbe Pfad in den Läufen
dreier Rollen desselben Vorgangs. **Ausgangs-Kandidat** (Architect-Entscheidung, nicht
getroffen): `tools/schema/apply-rollout.sh` nimmt `plan.yaml` und `down.sql` nach dem
Lauf zurück (wie der Guard-Test es tut) — ein Träger, der die Rücknahme mechanisch
trägt statt der Disziplin jedes Lesers.
