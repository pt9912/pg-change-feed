# BEO-PGC/cdc-capture-lag-real

**Sub-Area:** Observability/Replication-Adapter (reales `cdc_capture_lag`;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: [`LH-FA-ADM-004`](../../../../../../spec/lastenheft.md)
(messbarer CDC-Abstand) ist nur als Näherung geliefert
(`now() − max(committed_at)`, Persistenz-Zeit-Proxy, slice-013). Ein
reales `cdc_capture_lag` (Quell-Commit-Zeit → CDC-Verfügbarkeit) bräuchte
den bereits vorhandenen `pglogrepl.CommitMessage.CommitTime`-Zeitstempel
bis in `cdc.transaction.committed_at` gespiegelt — keine neue
Wire-Verbindung, aber eine Änderung an vier Schichten zugleich
(Replication-Decoder, Domain, Application/Ports, Store-Adapter), zu groß
für einen einzelnen Slice neben anderer Arbeit.

Deklaration: slice-013-Plan §1 (Ausschluss, Klasse Folge-Slice), `tools/
schema/nacharbeit-observability.sql` (Näherungs-Kommentar, Klasse Grenze).
