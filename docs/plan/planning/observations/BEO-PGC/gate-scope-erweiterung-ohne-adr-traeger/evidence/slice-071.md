**Vorgang:** slice-071
**Fund:** Der Implementer erweiterte `.a-check.yml`s `composition_root` um
`tools/**`, damit der neue Wegwerf-Client den generierten gRPC-Stub
importieren darf — eine Gate-Scope-Erweiterung ohne ADR, begründet mit dem
Präzedenzfall `test/integration/**`. Der Reviewer hat sie als HIGH bestätigt:
`ADR-0041` nennt `composition_root` in der Definition des deklarativen
Standes, aber **nicht** in der Ausnahme-Klausel (die nur „Schichten-Globs und
Edges" führt); zudem gab die Aufnahme dem ganzen `tools/**`-Baum
unbeschränktes Importrecht (eine Probe mit Produktions-Use-Case-Import blieb
grün). Aufgelöst über `ADR-0068` — begrenzte Kante `tooling → adapters`
statt Glob-Erweiterung; die Klausel ist seither an beiden Enden geschärft.
Erstauftreten.
