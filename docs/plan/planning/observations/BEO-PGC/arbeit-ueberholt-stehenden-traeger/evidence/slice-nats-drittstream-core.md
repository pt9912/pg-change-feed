# Beleg: slice-nats-drittstream-core

Vorgang: `slice-nats-drittstream-core` — NATS als dritter,
vollinhaltstragender Zustellweg (`ADR-0100`); der neue Schlüssel
`nats_stream_token` wächst in die Klasse der zugangsdaten-tragenden
Konfigurationsdatei-Felder und hebt ihre Zahl von sechs auf sieben.

Fund A (Zahl der Klasse): Der Implementer-Suchlauf lief über `spec/`,
`docs/user/` und `docs/plan/adr/` und fand dort zwei lebende Träger; die
dritte lebende Stelle — `internal/bootstrap/config_file.go` — lag
außerhalb dieser Wurzeln und blieb stehen: der `fileConfig`-Kommentar sagte
weiter „die zwei Token-Schlüssel", 15 Zeilen unter einem Kommentar, den
derselbe Commit bereits auf „drei" gezogen hatte. Gefunden hat das nicht der
Suchlauf des Slice, sondern der unabhängige Reviewer (F-2, HIGH). Der
nachgeholte Suchlauf über `internal/` fand an derselben bewegten Eigenschaft
zwei weitere Stellen (`mergeConfig`-Kommentar „die zwei Token-Klassen";
`config_file_internal_test.go` „jeder ihrer sechs Schlüssel" bei einer Liste
mit sieben Einträgen) — beide nachgezogen.

Fund B (Gate-Index-Zeile): `harness/README.md`s `make test-integration`-Zeile
enumeriert jede Rundlauf-Phase und wurde bei jedem vergleichbaren Rundlauf im
selben Commit mitgezogen (SSE `89d31d1`, gRPC `b835dde`, `GET /changes`
`93ac92e`); dieser Commit fügte eine vollständige neue Phase hinzu und ließ
die Zeile unverändert. Ebenfalls vom Reviewer gefunden (F-3, HIGH), nicht vom
eigenen Suchlauf — die Zeile stand nicht im Diff.

Fund C (Zahl in Accepted-ADRs): Die Zahl-Aussagen der `Accepted`-ADRs
`ADR-0088`/`0089`/`0091`/`0092` über dieselbe Klasse trugen die alte Sechs
weiter; die erste Meldung des Slice nannte als Wächter eine von `ADR-0092`
bereits supersedierte ADR (`ADR-0089`) und schlug eine zu enge Remediation
vor. Über zwei Architect-Züge (`ADR-0101`, `ADR-0102`) klauselgenau
supersediert; `ADR-0102` ersetzt die Aufzählung durch eine **Regel über
einen definierten Umfang** und fand dabei eine weitere überfahrene Stelle
(`ADR-0089:73`).

Quelle: Review zu `slice-nats-drittstream-core` F-2/F-3/F-4 ·
`docs/plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md` ·
`docs/plan/adr/0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md` ·
`docs/plan/planning/in-progress/slice-nats-drittstream-core.md` §6.
