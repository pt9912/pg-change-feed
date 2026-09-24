Zustand: offen — Ausgang: **weiter offen**, adressiert. Die Grenze steht im Handbuch
(`docs/user/benutzerhandbuch.md`, §4 „Bestand als Backfill überführen“, Absatz „Sperre der
Tabelle“, Gegenrichtung); der Kontext-Abbruch an der wartenden Anweisung ist im Store-Tier
gebunden (`TestImportLockWaitEndsWithTheContext`, `make test-replication`). Eine
Zeitgrenze der Sperranweisung (`lock_timeout`) ist eine Designänderung an `ADR-0118`
Festlegung 1 und nicht entschieden; sie bräuchte eine Folge-ADR. Zähler (abgeleitet): 1×
(evidence/slice-backfill-e2e.md).
