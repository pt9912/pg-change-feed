**Vorgang:** slice-routing-lesewege (Review F-4, LOW; neunte Datei des Eintrags)

**Fund:** Der Godoc eines Tests in `internal/application/usecase/capture/service_test.go` sprach von der „Strecke des Zustellziels vom Assembler zum Stream“. Der Test baut die Transaktion von Hand (das Ziel ist dort per `WithRouteTarget` gesetzt) und endet am Fake des Stream-Ports; der Assembler-Schritt und der `Broadcaster` liegen nicht auf der gefahrenen Strecke. Der Reviewer fand es durch Nachfahren der zugesagten Strecke. Die Fixrunde begrenzte den Godoc auf „von der Transaktion … zum Stream-Port-Fake“ und nannte den Assembler-Schritt ausdrücklich als ausgenommen.

**Form (Ausprägung):** **Allaussage über einen Fall** — der Kommentar sagt eine Strecke zu, der Test fährt ein Teilstück. Vor dem Merge gefunden, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-routing-lesewege.md` (F-4) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-lesewege.md` (§5, Zeile F-4). <!-- d-check:status-provenance -->
