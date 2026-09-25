Zustand: **geplant** — Ausgang: **geplant** →
`slice-backfill-speicher-untersuchung`
(Messreihe mit `tools/bench-backfill.sh` samt Grundlinie, Ursache benennen; wächst
der Speicher mit der Tabellengröße, folgen ein Befund und ein eigener
Änderungs-Slice, sonst trägt das Handbuch die gemessene Grenze) · seit
welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3 und §5 (h)).

Gegenstand: der Speicher des Feed-Containers im Backfill liegt über dem Bedarf eines
Blocks (Handbuch, Abschnitt „Grenzwerte“: Spitzen von einigen hundert MiB gegen etwa
74 KB je Block, übernommen); die Ursache ist nicht untersucht. Die Grenze der
Blockgröße („zählt Zeilen, nicht Bytes“) steht in der Port-Doku
(`internal/application/port/outbound/tablesnapshot.go`, `NextBlock`) und im Kommentar
an `DefaultBlockSize`.

Zähler: 1× (Datei unter `evidence/`).
