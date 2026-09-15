package main

import (
	"net/url"
)

// Subject leitet das tabellen-granulare Wecksignal-Subjekt aus Quelle, Schema
// und Tabelle ab (`SPEC-017`): `cdc.changes.<source_id>.<schema>.<table>`.
// Das Subjekt wird abgeleitet, nicht handgetippt — dieselbe Form, die der
// Feed-Container beim Publizieren bildet (`ADR-0056`).
func Subject(sourceID, schema, table string) string {
	return "cdc.changes." + sourceID + "." + schema + "." + table
}

// ChangesURL baut die Lese-Adresse für die Änderungen einer Quelle/Tabelle
// über die HTTP-/JSON-API (`LH-FA-SST-006`): ein `GET` auf den Changes-
// Lesezugriff mit den drei Filtern aus dem Wecksignal-Subjekt. `base` ist die
// Horch-Adresse des Feed-Containers (`CDC_HTTP_ADDR`, Form `host:port`).
func ChangesURL(base, sourceID, schema, table string) string {
	u := url.URL{Scheme: "http", Host: base, Path: "/changes"}
	q := u.Query()
	q.Set("source", sourceID)
	q.Set("schema", schema)
	q.Set("table", table)
	u.RawQuery = q.Encode()
	return u.String()
}
