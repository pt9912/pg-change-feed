**Vorgang:** slice-sdk-csharp-publish-workflow

**Fund:** Der Implementer-Commit (`b0fbfa91`) legte einen neuen,
eigenständigen Release-Trigger an — `.github/workflows/sdk-csharp-release.yml`,
Tag-Namensraum `sdk-csharp-v*`, eigenes Secret `NUGET_API_KEY` — ohne
`docs/user/releasing.md` mitzuziehen. Der Reviewer fand das real (Review zu
`slice-sdk-csharp-publish-workflow`, F-1, MEDIUM): kein Treffer für
`sdk-csharp`, `NUGET_API_KEY` oder `PgChangeFeed.Client` im
gesamten Dokument, obwohl `releasing.md` laut eigenem §1-Zweck genau der
Träger ist, der beschreibt, wie ein Release ausgelöst wird und welche
Secrets dafür nötig sind. Der Slice-Plan (§1 „Ausdrücklich NICHT in diesem
Slice") listete diesen Doku-Träger nicht als bewusst zurückgestellt — die
Auslassung war unbedacht, nicht begründet.

Der Implementer hat den Fund in der Fixrunde (`c582eaaf`) real behoben:
`docs/user/releasing.md` §4 trägt jetzt einen eigenen Unterabschnitt
„SDK-Release: NuGet.org-Publish für `PgChangeFeed.Client`" samt
Secret-Tabellen-Zeile. Reviewer-Fixrunden-Nachprüfung und Verifier haben den
Nachzug unabhängig gegen die reale Workflow-Datei gehalten — kein Drift.

Der Reviewer selbst markierte den Fund als „erstes Auftreten" einer neuen
Klasse und ausdrücklich „kein Steering-Loop-Zähler-Eintrag nötig unter 3x" —
dieser Eintrag legt den Zähler jetzt an (Entscheidung: eigene, vom
`benutzerhandbuch.md`-Eintrag getrennte Klasse, siehe `observation.md`
Abgrenzung), statt den Fund unregistriert zu lassen.

Quelle: Review zu `slice-sdk-csharp-publish-workflow` F-1 ·
Fixrunden-Nachprüfung desselben Reports · Planner-Closure-Entscheidung.
