Zustand: **verkörpert** — Ausgang: **verkörpert** →
`slice-capture-leerlauf-quellbelege`:
(1) die einmalige Messung eines Keepalive inmitten einer Transaktion ist ein committeter
Test im Tier `make test-replication` (`TestSourceKeepaliveInsideTransactionDeliversWholeTransaction`,
`internal/adapters/driving/replication/receive/sourcekeepalive_test.go`, eigener Lauf in
`tools/harness/run-replication-tests.sh` mit `--- PASS`-Wächter), gelaufen an PostgreSQL 17.11
und 18.6 (Verifikations-Report `verifikation-slice-capture-leerlauf-quellbelege` §1, gedruckte
Zeilen; Mutation „Position `+ 1 GiB`“ rot an beiden, §4 K1) und in CI je Leg
(`e2e.yml`-Lauf 36287009221: Schritt „Replication-Tier“ Leg PostgreSQL 17 im ersten, Leg
PostgreSQL 18 im zweiten Versuch `success`); die Ergänzung der Grenze trägt `ADR-0129`
(`Supersedes ADR-0121` teilweise; die Lese-Prüfung der ADR gegen `AGENTS.md` §3.12 steht
aus, Risiko §6 im Plan);
(2) die Kette „Fehlerschwelle erreicht → Container endet“ ist die Runner-Phase „Fehlerschwelle
beendet den Container“ in `make test-integration` (Ende, Ausgang 1, Abbruch-Zeile, Fehlerzustand;
die Klasse des Ausgangs als benannte Grenze mit Träger `slice-wal-fehlerschwelle-ausgangsklasse`);
die Mutation „`stopStream` in der Schwellen-Prüfung entfernt“ färbt sie rot (vom Verifier
nachgefahren, Verifikations-Report §4 P2) · seit slice-capture-leerlauf-quellbelege
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3, §5 (b) und (e)).

Nicht getragen, weiter benannt: `proto_version` 2 mit Streaming großer Transaktionen, eine
gleichzeitige zweite Quelltransaktion, mehr als ein Keepalive je Transaktion, die Startposition
des Neustarts und die Wartezeit auf die Inaktivität des Slots im Test (`ADR-0129`
Festlegung 2 Punkte 1, 2 und 4).

Zähler: 1× (Datei unter `evidence/`; zwei Fälle in diesem einen Vorgang zählen
einmal).
