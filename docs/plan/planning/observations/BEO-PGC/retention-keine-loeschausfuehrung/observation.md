# BEO-PGC/retention-keine-loeschausfuehrung

**Sub-Area:** Retention/ChangeStore (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: `spec/lastenheft.md`s `LH-FA-RET-002`…`006` (konfigurierbare,
zeit- und consumer-basierte Aufbewahrung, Erkennbarkeit blockierender
Consumer, Kontrolle des Datenwachstums) haben in `internal/domain/model/retention.go`
(`RetentionPolicy`, `AllowsDeletion`, [`ADR-0014`](../../../../adr/0014-retention-domain-policy.md))
eine reine Entscheidungslogik — sie bewertet, *ob* ein Change gelöscht werden
dürfte. Es gibt aber **keinen Aufrufer**: kein Use-Case, kein CLI-Unterbefehl
und kein periodischer Job ruft `AllowsDeletion` auf oder löscht je eine Zeile
aus `cdc.change`. `RetentionPolicy` wird reponweit nur in ihrer eigenen
Testdatei verwendet. Die zugehörige Metrik `cdc_storage_bytes`
(`LH-FA-RET-006`) fehlt ebenfalls (`tools/schema/nacharbeit-observability.sql`
nennt sie explizit als „nicht abgedeckt"). `LH-FA-RET-001`
(persistente Speicherung, kein Datenverlust bei Neustart) ist davon
unberührt und real erfüllt — das ist keine Regression, sondern eine noch nie
gebaute Hälfte der Fähigkeit.

## Benannt, nicht gezählt

Gefunden bei einer repo-weiten Fähigkeiten-Bestandsaufnahme (Fork-Recherche,
2026-09-12), nicht bei einer Slice-/Welle-Closure — kein abgeschlossener
Vorgang trägt bisher einen Beleg. Der Zähler bleibt bei 0×, bis ein
Slice/Review/Verify diese Beobachtung tatsächlich berührt.
