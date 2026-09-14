**Vorgang:** slice-061
**Fund:** Acht in `slice-060` gebaute HTTP-Handler
(Acknowledge/Position/Remove-Consumer, Enable/Disable/Status/List-Table,
Retention-Lauf) waren im Bootstrap (`internal/bootstrap/wiring.go`) mit
`nil`-Use-Cases verdrahtet — ihre Whitebox-Unit-Tests aus `slice-060`
konstruierten `Config` direkt mit Fakes und liefen nie über die reale
Verdrahtung, daher blieb die Lücke dort unsichtbar. Erst `slice-061`s
realer `make test-integration`-Rundlauf gegen den vollständig
bootstrap-verdrahteten Feed-Container hätte den Fehler (Panic bei realem
Aufruf) offengelegt; der Implementer fand und behob ihn vorher beim Bauen
des Rundlaufs. Erstauftreten.
