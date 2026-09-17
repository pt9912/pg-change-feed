package main

import (
	"net/url"
)

// TablesURL baut die Lese-Adresse der Tabellen-Auflistung über die
// HTTP-/JSON-API (`LH-FA-SST-006`): ein `GET` auf den `reader`-Endpunkt
// `/tables` mit seinen zwei Pflichtfeldern `source` und `publication`. `addr`
// ist die Horch-Adresse des Feed-Containers (`CDC_HTTP_ADDR`, Form
// `host:port`).
func TablesURL(addr, source, publication string) string {
	u := url.URL{Scheme: "http", Host: addr, Path: "/tables"}
	q := u.Query()
	q.Set("source", source)
	q.Set("publication", publication)
	u.RawQuery = q.Encode()
	return u.String()
}
