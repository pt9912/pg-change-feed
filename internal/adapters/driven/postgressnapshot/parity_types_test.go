package postgressnapshot

import (
	"fmt"
	"strings"
)

// typeColumn ist eine Zeile der Typ-Tabelle des Paritätstests (`ADR-0115`
// Festlegung 4): eine Spalte mit Typ-DDL und je einem SQL-Ausdruck für die
// Wert-Zeile und die Rand-Zeile. Ein leerer `edge` wiederholt den Wert.
type typeColumn struct {
	name  string
	ddl   string
	value string
	edge  string
}

// typeSet trägt die Typ-Tabelle samt den benutzerdefinierten Typen, die sie
// braucht (Enum, Range, Domains, Composites — je Test mit eigenem Suffix).
type typeSet struct {
	create  []string
	drop    []string
	columns []typeColumn
}

// requiredTypes nennt die Typ-DDLs, die der Typ-Satz tragen muss
// (`ADR-0115` Festlegung 4); `TestImageParityWalAndBackfill` prüft sie
// gegen die Tabelle, damit das Streichen einer Zeile rot wird.
var requiredTypes = []string{
	"boolean", "char(5)", "inet", "cidr", "macaddr", "money", "bit(5)", "varbit", "uuid", "xml",
	"oid", "regclass", "int4range", "int4multirange", "tsvector", "tsquery", "point", "polygon",
	"timetz", "hstore", "citext", "ltree", "hstore[]",
}

// typeTable liefert den Typ-Satz der Bild-Parität, einschließlich der
// Extension-Typen `hstore`, `citext` und `ltree`: sie liegen in einem
// eigenen Schema `snap_ext_<sfx>`, das `set.drop` samt den Extensions
// entfernt.
func typeTable(sfx string) typeSet {
	mood, rng := "snap_mood_"+sfx, "snap_rng_"+sfx
	c1, c2 := "snap_c1_"+sfx, "snap_c2_"+sfx
	dInt, dText, dArr, dBool := "snap_dint_"+sfx, "snap_dtext_"+sfx, "snap_darr_"+sfx, "snap_dbool_"+sfx
	ext := "snap_ext_" + sfx
	set := typeSet{
		create: []string{
			"CREATE SCHEMA " + ext,
			"CREATE EXTENSION hstore SCHEMA " + ext,
			"CREATE EXTENSION citext SCHEMA " + ext,
			"CREATE EXTENSION ltree SCHEMA " + ext,
			"CREATE TYPE " + mood + " AS ENUM ('sad', 'ok', 'hä ppy')",
			"CREATE TYPE " + rng + " AS RANGE (subtype = integer)",
			"CREATE TYPE " + c1 + " AS (a int, b text)",
			"CREATE TYPE " + c2 + " AS (f bool, i inet, c char(4))",
			"CREATE DOMAIN " + dInt + " AS integer CHECK (VALUE > 0)",
			"CREATE DOMAIN " + dText + " AS text",
			"CREATE DOMAIN " + dArr + " AS integer[]",
			"CREATE DOMAIN " + dBool + " AS boolean",
		},
		drop: []string{
			"DROP DOMAIN IF EXISTS " + dInt + ", " + dText + ", " + dArr + ", " + dBool + " CASCADE",
			"DROP TYPE IF EXISTS " + mood + ", " + rng + ", " + c1 + ", " + c2 + " CASCADE",
			"DROP EXTENSION IF EXISTS hstore, citext, ltree CASCADE",
			"DROP SCHEMA IF EXISTS " + ext + " CASCADE",
		},
	}
	add := func(name, ddl, value, edge string) {
		if edge == "" {
			edge = value
		}
		set.columns = append(set.columns, typeColumn{name: name, ddl: ddl, value: value, edge: edge})
	}

	// Skalare Grundtypen.
	add("bo", "boolean", "true", "false")
	add("i2", "smallint", "12", "-32768")
	add("i4", "integer", "2147483647", "-2147483648")
	add("i8", "bigint", "9223372036854775807", "-9223372036854775808")
	add("rl", "real", "0.1", "'NaN'")
	add("f8", "double precision", "0.1::float8 + 0.2::float8", "'-Infinity'")
	add("nu", "numeric", "12345.6700", "'NaN'")
	add("mo", "money", "1234567.89", "-0.5")
	add("ch5", "char(5)", "'ab'", "''")
	add("ch1", "char(1)", "' '", "''")
	add("vc", "varchar(10)", "'héllo'", "''")
	add("tx", "text", `E'a\tb\n"c"\\d,{}'`, "''")
	add("tu", "text", "'日本語 ✓'", "'  padded  '")
	add("qc", `"char"`, "'A'", "''")
	add("nm", "name", "'sample_name'", "''")
	add("bya", "bytea", `'\xdeadbeef'`, `'\x'`)
	add("bt", "bit(5)", "B'10101'", "B'00000'")
	add("vb", "varbit", "B'1011'", "B''")
	add("ui", "uuid", "'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'", "'00000000-0000-0000-0000-000000000000'")

	// Zeit.
	add("dt", "date", "'2026-09-24'", "'-infinity'")
	add("tm", "time", "'08:11:12.5'", "'24:00:00'")
	add("tz", "timetz", "'08:11:12.5+02'", "'00:00:00-05:30'")
	add("ts", "timestamp", "'2026-09-24 08:11:12.5'", "'infinity'")
	add("tt", "timestamptz", "'2026-09-24 08:11:12.5+00'", "'-infinity'")
	add("iv", "interval", "'1 year 2 mons 3 days 04:05:06.7'", "'-1 days -02:03:04.5'")

	// Netz- und Adress-Typen.
	add("ia", "inet", "'1.2.3.4'", "'2001:db8::1'")
	add("ib", "inet", "'10.0.0.1/8'", "'::1/64'")
	add("cd", "cidr", "'10.1.2.0/24'", "'2001:db8::/32'")
	add("ma", "macaddr", "'08:00:2b:01:02:03'", "'ff:ff:ff:ff:ff:ff'")
	add("m8", "macaddr8", "'08:00:2b:01:02:03:04:05'", "'ff:ff:ff:ff:ff:ff:ff:ff'")

	// Dokument-Typen.
	add("js", "json", `'{"a": [1,  2], "b": null}'`, `'""'`)
	add("jb", "jsonb", `'{"a": [1,  2], "b": null}'`, "'{}'")
	add("jp", "jsonpath", "'$.a[*] ? (@ > 1)'", "'$'")
	add("xm", "xml", `'<?xml version="1.0"?><b/>'`, "'<a>&amp;</a>'")

	// Objekt- und Systemtypen.
	add("oi", "oid", "12345", "'4294967295'")
	add("rc", "regclass", "'pg_catalog.pg_class'", "'pg_catalog.pg_type'")
	add("rt", "regtype", "'integer'", "'text'")
	add("rp", "regproc", "'now'", "'now'")
	add("xi", "xid", "'12345'", "'0'")
	add("td", "tid", "'(0,1)'", "'(4294967295,65535)'")
	add("ls", "pg_lsn", "'16/B374D848'", "'0/0'")
	add("sn", "pg_snapshot", "'10:20:10,14'", "'10:20:'")

	// Volltext und Geometrie.
	add("tv", "tsvector", "'a:1 fat:2 cat:3'", "''")
	add("tq", "tsquery", "'fat & !cat'", "'a'")
	add("pt", "point", "'(1.5,2)'", "'(-0,1e300)'")
	add("ln", "line", "'{1,2,3}'", "'{0,1,0}'")
	add("ls2", "lseg", "'[(0,0),(1,1)]'", "'[(-1,-1),(0,0)]'")
	add("bx", "box", "'(1,1),(0,0)'", "'(0,0),(0,0)'")
	add("pa", "path", "'[(0,0),(1,1),(2,0)]'", "'((0,0),(1,1))'")
	add("pg", "polygon", "'((0,0),(1,1),(2,0))'", "'((0,0))'")
	add("ci", "circle", "'<(0,0),1.5>'", "'<(1,1),0>'")

	// Ranges und Multiranges.
	add("r4", "int4range", "'[1,10)'", "'empty'")
	add("rt2", "tstzrange", "'[2026-09-24 00:00+00,infinity)'", "'empty'")
	add("rd", "daterange", "'[2026-09-01,2026-10-01)'", "'empty'")
	add("mr", "int4multirange", "'{[1,3),[5,9)}'", "'{}'")
	add("rx", rng, "'[1,5)'", "'empty'")

	// Extension-Typen.
	add("hs", ext+".hstore", `'a=>1, "b c"=>NULL, d=>"x\"y"'`, "''")
	add("ci2", ext+".citext", "'MiXed Case'", "''")
	add("lt", ext+".ltree", "'Top.Science.Astronomy'", "''")

	// Enum, Domains, Composites.
	add("en", mood, "'ok'", "'hä ppy'")
	add("di", dInt, "3", "1")
	add("dx", dText, "'x'", "''")
	add("da", dArr, "'{1,2}'", "'{}'")
	add("db", dBool, "true", "false")
	add("c1", c1, "ROW(1, 'x y')", "ROW(NULL, '')")
	add("c2", c2, "ROW(true, '1.2.3.4', 'ab')", "ROW(false, '::1/64', '')")

	// Arrays.
	add("a4", "integer[]", "'{1,2,NULL}'", "'{}'")
	add("at", "text[]", `ARRAY['x', 'y z', NULL, '', 'a"b', 'c\d']`, "ARRAY[]::text[]")
	add("ab", "boolean[]", "'{t,f,NULL}'", "'{}'")
	add("ai", "inet[]", "'{1.2.3.4,10.0.0.1/8}'", "'{::1}'")
	add("ac", "char(3)[]", "ARRAY['a', 'bc']::char(3)[]", "ARRAY['']::char(3)[]")
	add("ay", "bytea[]", `ARRAY['\xdead'::bytea, '\x'::bytea]`, "ARRAY[]::bytea[]")
	add("ag", "timestamptz[]", "ARRAY['2026-09-24 08:11:12.5+00'::timestamptz]", "ARRAY['infinity'::timestamptz]")
	add("af", "float8[]", "'{{1.5,2},{3,4}}'", "'[0:1]={5,6}'")
	add("am", "money[]", "ARRAY[1.5::money, 2::money]", "ARRAY[]::money[]")
	add("ae", mood+"[]", "ARRAY['ok', 'sad']::"+mood+"[]", "ARRAY[]::"+mood+"[]")
	add("ao", c1+"[]", "ARRAY[ROW(1, 'a')::"+c1+"]", "ARRAY[]::"+c1+"[]")
	add("aj", "interval[]", "ARRAY['1 day'::interval]", "ARRAY[]::interval[]")
	add("an", "numeric[]", "ARRAY[1.10, 2]", "ARRAY[]::numeric[]")
	add("au", "uuid[]", "ARRAY['a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'::uuid]", "ARRAY[]::uuid[]")
	add("ad", "cidr[]", "'{10.1.2.0/24}'", "'{}'")
	add("abt", "bit(3)[]", "ARRAY[B'101']", "ARRAY[]::bit(3)[]")
	add("ar", "int4range[]", "ARRAY['[1,3)'::int4range]", "ARRAY[]::int4range[]")
	add("ax", dBool+"[]", "ARRAY[true]::"+dBool+"[]", "ARRAY[]::"+dBool+"[]")
	add("ah", ext+".hstore[]", "ARRAY['a=>1'::"+ext+".hstore, 'b=>NULL'::"+ext+".hstore]", "ARRAY[]::"+ext+".hstore[]")
	return set
}

// definition ergibt die Spaltendefinitionen der Test-Tabelle: `id`, die
// Typ-Spalten, eine gelöschte und eine generierte Spalte (`STORED`, unter
// PostgreSQL 18 zusätzlich `VIRTUAL`), die im WAL und im Backfill fehlen.
func (s typeSet) definition(virtualGenerated bool) string {
	defs := []string{"id int PRIMARY KEY", "dropme int"}
	for _, column := range s.columns {
		defs = append(defs, column.name+" "+column.ddl)
	}
	defs = append(defs, "gen int GENERATED ALWAYS AS (id * 2) STORED")
	if virtualGenerated {
		defs = append(defs, "genv int GENERATED ALWAYS AS (id * 3) VIRTUAL")
	}
	return strings.Join(defs, ", ")
}

// insert ergibt die Zeilen der Tabelle: Zeile 1 mit den Werten, Zeile 2
// mit den Rand-Werten, Zeile 3 mit NULL in jeder Typ-Spalte.
func (s typeSet) insert(qualified string) string {
	names := []string{"id"}
	values := [3][]string{{"1"}, {"2"}, {"3"}}
	for _, column := range s.columns {
		names = append(names, column.name)
		values[0] = append(values[0], fmt.Sprintf("CAST((%s) AS %s)", column.value, column.ddl))
		values[1] = append(values[1], fmt.Sprintf("CAST((%s) AS %s)", column.edge, column.ddl))
		values[2] = append(values[2], fmt.Sprintf("CAST(NULL AS %s)", column.ddl))
	}
	rows := make([]string, len(values))
	for i, row := range values {
		rows[i] = "(" + strings.Join(row, ", ") + ")"
	}
	return "INSERT INTO " + qualified + " (" + strings.Join(names, ", ") + ") VALUES " + strings.Join(rows, ", ")
}

// missingRequired nennt die Typen aus `requiredTypes`, die keine Spalte der
// Tabelle trägt; ein Typ eines Schemas zählt über seinen Namen nach dem Punkt.
func (s typeSet) missingRequired() []string {
	var missing []string
	for _, want := range requiredTypes {
		found := false
		for _, column := range s.columns {
			if column.ddl == want || strings.HasSuffix(column.ddl, "."+want) {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, want)
		}
	}
	return missing
}

// ddlOf liefert den Typ einer Spalte für die Fehlermeldung.
func (s typeSet) ddlOf(name string) string {
	for _, column := range s.columns {
		if column.name == name {
			return column.ddl
		}
	}
	return "?"
}
