# BEO-PGC/alt-tag-lauf-vorbedingung-am-juengsten-tag

**Sub-Area:** Schemamigration (d-migrate; Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: Lauf 5 von `tools/harness/run-schema-rollout-guard-test.sh` (Alt-Tag-Lauf, `ADR-0114` Entscheidung 7) rollt das Schema des jüngsten `v*`-Tags aus und prüft vor dem Upgrade Vorbedingungen, die das **Delta** des laufenden Slice beschreiben (Objekte, die der Tag noch nicht trägt). Die Vorbedingungen hängen damit an dem Tag, der zum Zeitpunkt des Schreibens der jüngste ist: ein späterer Release-Tag, der das Delta trägt, macht sie falsch. Belegt: die Vorbedingungen aus dem Stand vor der Antragsart `backfill` (`cdc_admin` ohne `UPDATE` auf die Antrags-Tabelle, `backfill_table` fehlt, die Menge trägt kein `backfill`) waren gegen `v0.2.0` falsch und wurden am Start von `slice-transformationen-antragsweg-schema` gemessen und auf das neue Delta gelegt. Ein Tag, der das jetzige Delta trägt, lässt den Lauf laut scheitern („Vorbedingung fehlgeschlagen“), kein stilles Grün.

**Abgrenzung.** `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` betrifft die Release-Dokumentation; hier ist die Kopplung eines Testlaufs an den Zustand des Tag-Bestands die Ursache.

**Warum das zählt:** Der Lauf ist der einzige Upgrade-Beleg gegen einen echten Alt-Bestand; nach jedem Release verschiebt sich sein Bezugspunkt, ohne dass ein Gate den Lauf fährt.

Deklaration: `slice-transformationen-antragsweg-schema`.
