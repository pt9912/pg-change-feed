**Vorgang:** slice-sdk-readme-nutzerdoku (Review F-11, INFO; Verifikation §3 Zeile F-11)

**Fund:** Zwei Quellen sagten dasselbe verschieden — die Ausprägung liegt zwischen einem **Quelltext-Kommentar** und der README desselben Packages, nicht zwischen Handbuch und Pflichtenheft; die Klasse ist dieselbe. Der XML-Kommentar von `PgChangeFeedGrpcClient` forderte, der Aufrufer schalte den `AppContext`-Schalter `Http2UnencryptedSupport` einmal selbst ein, bevor er sich mit einem Klartext-Server verbindet; die C#-README und der Realserver-Test (`GrpcRealserverTests.cs`) verbinden sich mit `http://` ohne den Schalter. Der Kommentar lag außerhalb des README-Diffs. Behoben in der Fixrunde (der Kommentar nennt: eine `http://`-Adresse verbindet auf .NET 10 ohne weitere Einrichtung, der Konstruktor setzt keinen prozessweiten Schalter); der Verifier las den Kommentar gegen README und Test, den Realserver-Test nicht gefahren.

Quelle: `docs/reviews/review-slice-sdk-readme-nutzerdoku.md` (F-11) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-readme-nutzerdoku.md` (§3 Zeile F-11). <!-- d-check:status-provenance -->
