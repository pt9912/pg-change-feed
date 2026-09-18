# Beleg: slice-091

Vorgang: `slice-091` — Coverage Cluster C (Zustell- und Betriebs-Rand).

Fund: Der Slice gab dem Paket `internal/adapters/driving/grpc/streamv1` eine
eigene Testdatei (`changestream_test.go`, `package streamv1_test`) — und machte
damit **drei** Sätze in `harness/sensors/coverage-gate.md` falsch, einer Datei,
die **nicht im Diff** lag und die niemand angefasst hat:

1. die Überschrift des Aufzählungspunkts — „**Fünf** Pakete des Gegenstands
   [ohne] Testdatei": `streamv1` hat seit diesem Vorgang eine;
2. der Schlusssatz — der Lauf weise „die beiden Pakete mit Statements" als
   `coverage: 0.0% of statements` aus: er weist für `streamv1`
   `ok … coverage: 2.4%` aus, die `0.0%`-Zeile trägt heute nur
   `cmd/pg-change-feed`;
3. die **Gruppierung selbst**.

**Beide Nachbarsätze waren vorher wahr.** Gemessen am Parent-Stand
(`git archive f90c3f4^`, `go list -f '{{.ImportPath}} Test={{len .TestGoFiles}}
XTest={{len .XTestGoFiles}}'`): dort waren es **genau fünf** Pakete mit
`Test=0 XTest=0` — `postgresstorage/queries`, `application/port/inbound`,
`domain/errors`, `streamv1`, `cmd/pg-change-feed`. Die Liste **war** mechanisch.
Die Arbeit hat sie es nicht mehr sein lassen.

**Die Reparatur hat den Fall zunächst verschärft.** Der erste Versuch begründete
`streamv1`s Zugehörigkeit mit `go list -f '{{len .TestGoFiles}}'`. Gemessen
trifft diese Zählung **23** der 31 Pakete des Gegenstands (darunter `natsnotify`
0/1, `port/outbound` 0/3, **13** `usecase/*` 0/1) und ist damit keine
Gruppierungsregel; sie ist für **jedes** Paket mit externem Testpaket wahr.
Gefunden hat das der **Delta-Review** zu `slice-091` (D-1) — nicht das
Review und nicht die Verifikation: beide hatten keinen Anlass, eine Datei zu
öffnen, die nicht im Diff lag. Die vierte Runde hat die Gruppe dann auf ihren
mechanischen Bestand gebracht (vier Pakete mit `Test=0 XTest=0`) und `streamv1`
eine **eigene**, substanziierte Rolle gegeben (45 von 86 Statements aus dem
eigenen Testlauf gegen 76 von 86 im Gegenstand).

Quelle: Delta-Review zu `slice-091` (D-1) ·
`harness/sensors/coverage-gate.md` §Grenze Punkt 1 (berichtigt in `af3ea9f`) ·
`git archive f90c3f4^` + `go list` (Parent-Stand, gemessen) ·
Verifikationsbericht zu `slice-091` (V-6).
