**Vorgang:** slice-sdk-tls-optionen (Erst-Review F-1, Planner-Nachtrag bei der Closure).

**Fund:** Die Umleitungs-Hälfte von `AGENTS.md` §3.1 in der Rolle Implementer: Text wurde per
`cat >>` an `tools/harness/run-sdk-csharp-integration-tests.sh` angehängt, statt per Edit/Write.
Der Reviewer (Review-Report zu `slice-sdk-tls-optionen`, F-1) und der Verifier (Verifikation,
§4 „Prozessfund F-1“) halten fest, dass der Guard Umleitungen nicht liest und der Diff die
Schreibweise nicht zeigt. Verifizierbar am Repo: nein; Ursprung der Angabe **übernommen**, nicht
am Guard gemessen.

**Form (Ausprägung):** Umleitungs-Anhang trotz geltender Regel; kein neuer Mechanismus, die
Regel in `AGENTS.md` §3.1 galt und der Guard-Vertrag (`MR-003`) nennt die Grenze (Umleitungen
werden nicht gelesen). Zwölfter Beleg dieses Eintrags.
