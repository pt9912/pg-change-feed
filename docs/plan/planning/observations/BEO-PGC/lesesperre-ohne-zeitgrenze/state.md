Zustand: **gestrichen** — Ausgang: **gestrichen** (akzeptiertes Negativ) · seit
welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3 und §5 (d)).

Begründung: die Wartezeit hängt an einer **fremden** Transaktion, die
`ACCESS EXCLUSIVE` hält, also an DDL des Betreibers. Das Handbuch nennt die
Gegenrichtung (`docs/user/benutzerhandbuch.md` §4, Absatz „Sperre der Tabelle“), und
der Kontext-Abbruch beendet die wartende Anweisung
(`TestImportLockWaitEndsWithTheContext`, `make test-replication`). Ein `lock_timeout`
bräche einen gesunden Run wegen einer vorübergehenden Sperre ab und ist eine
Designänderung an `ADR-0118` Festlegung 1; `ADR-0119` trägt den Trigger „Eine
Zeitgrenze der Sperranweisung wird eingeführt“. Ein Betreiber-Bericht liegt nicht vor
(kein Server-Tag trägt Backfill-Änderungen).

Zähler: 1× (Datei unter `evidence/`).
