**Vorgang:** slice-meldungscodes-registry-fehlerkopf (Review F-1; Verifikation §4, §5)

**Fund:** Der Umbau von `classifyRunError` (`internal/bootstrap/wiring.go`) von einer
`errors.Is`-Liste auf eine Klasse, die aus dem Code des Fehlerwerts folgt, ordnete zehn Sentinels
(`ErrSchemaStoreStorage`, `ErrNotify`, `ErrAdministrationStorage`, `ErrBackfillStorage`,
`ErrDiagnosticsStorage`, `ErrSnapshotPermission`, `ErrSnapshotConfiguration`,
`ErrSnapshotTransient`, `ErrSnapshotReplication`, `ErrSnapshotStorage`), die der Parent
(`cb8afd07`) auf `internal` fallen ließ, einer anderen Klasse zu, und wählte bei einer Kette mit
mehreren Sentinels den ersten Fund statt der Vorrangfolge des Parents. Betroffen wären
`error_class` im Heartbeat und das Metrik-Label dieser Fälle gewesen. Die Zusage „Klassen
unverändert“ ([`ADR-0144`](../../../../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2) stand in der ADR,
im Plan und im Handoff-Bericht; der Test `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`
führte die zehn Fälle nicht. Der Reviewer fand die Abweichung durch den Vergleich von
`git show cb8afd07:internal/bootstrap/wiring.go` mit dem Kopf. Schwere MEDIUM, Träger-Typ
Klassifikationsfunktion im Verdrahtungspaket. Behoben in `c080136a`, gebunden an einen Test über
die zehn Sentinels und fünf Ketten-Fälle; der Verifier fuhr vier Einzelmutationen (rot) und
maß 28 von 28 Sentinels vorher und nachher unverändert. Ein Re-Review nach der Fixrunde fand
nicht statt.

Quelle: `docs/reviews/review-slice-meldungscodes-registry-fehlerkopf.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-meldungscodes-registry-fehlerkopf.md` (§4, §5). <!-- d-check:status-provenance -->
