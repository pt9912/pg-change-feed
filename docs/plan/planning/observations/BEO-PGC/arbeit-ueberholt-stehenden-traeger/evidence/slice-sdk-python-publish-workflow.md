# Beleg: slice-sdk-python-publish-workflow (Coordinator-Fix `38a137f7`)

Vorgang: Während des Implementer-Zugs von `slice-sdk-python-publish-workflow`
fand der Coordinator beim Gegenlesen von `docs/user/releasing.md`
zwei Absätze, die noch „kein realer `sdk-csharp-v*`-Tag gesetzt" behaupteten
— obwohl `sdk-csharp-v0.1.0` bereits real gepusht und `PgChangeFeed.Client`
0.1.0 bereits real auf NuGet.org veröffentlicht war (`dotnet nuget push`
bestätigte `201 Created`).

**Die neue Form innerhalb dieser Klasse:** In allen bisherigen achtzehn
Belegen dieses Eintrags lag das überholende Ereignis selbst in einem Commit
dieses Repos (ein neuer Test, ein neues Runtime-Image, eine gelieferte
HTTP-Fläche). Hier liegt das überholende Ereignis **außerhalb jedes Diffs
und außerhalb jeder Versionskontrolle**: ein realer Tag-Push
(`sdk-csharp-v0.1.0`) und ein realer `dotnet nuget push` gegen die externe
NuGet-Registry — beides Handlungen des Repository-Betreibers nach der
`slice-sdk-csharp-publish-workflow`-Closure (`AGENTS.md` §3.10-Ereignis, kein
Commit). `docs/user/releasing.md` stand seit jener Closure mit einer
korrekten, aber inzwischen überholten Aussage da, und **kein** Sensor und
**kein** Diff dieses Repos konnte das je zeigen — die Wahrheit über den
Post-Push-Stand lebt strukturell außerhalb des Repos (PyPI-/NuGet-API,
`gh`-CLI).

Gefunden hat den Drift der Coordinator, nicht der Implementer-eigene
§3.13-Suchlauf dieses Slice — der Suchlauf dieses Slice war korrekt auf
`sdk-python-release\|PYPI_API_TOKEN` beschränkt (siehe Closure-Notiz §7 des
Slice-Plans) und hatte keinen Anlass, eine Aussage über den **C#**-Release-
Stand zu prüfen; das lag außerhalb des Gegenstands, den dieser Slice bewegt.
Der Fund entstand aus der allgemeinen Übung, beim Schreiben eines
Release-Workflow-Slices den ganzen Nachbarabschnitt in `docs/user/releasing.md`
mitzulesen, nicht aus einem gezielten §3.13-Suchlauf-Treffer.

Behoben im eigenständigen Commit `38a137f7` (getrennt vom Slice-Scope
committet, vom Reviewer dieses Slice ausschließlich auf Korrektheit
geprüft — nicht als Teil der DoD dieses Slice) und vom Reviewer sowie
Verifier dieses Slices real gegen NuGet.org
(`curl https://api.nuget.org/v3-flatcontainer/pgchangefeed.client/index.json`
→ `{"versions": ["0.1.0"]}`) und gegen `git tag -l "sdk-csharp-v*"`
nachgeprüft — dreifach unabhängig bestätigt korrekt.

**Warum dieser Beleg trotz der externen Auslösung zählt:** Die Klasse ist
nicht „ein Commit überholt einen Träger", sondern „eine bewegte Eigenschaft
überholt einen sie beschreibenden Träger" — und eine Eigenschaft kann sich
auch durch ein Ereignis außerhalb der Versionskontrolle bewegen. Kein
Sensor kann das fangen (`docs-check` liest keine externen Registries); der
einzige Wächter bleibt ein Mensch oder Agent, der beim Schreiben eines
benachbarten Abschnitts den bestehenden Nachbartext mitliest, statt sich auf
den eigenen Diff-Rand zu beschränken.

Quelle: `docs/reviews/review-slice-sdk-python-publish-workflow.md` <!-- d-check:status-provenance -->
§Eingangs-Kontext/„Eigenständig durchgeführte Prüfungen" (Coordinator-Fix-Prüfung) ·
`docs/reviews/verifikation-slice-sdk-python-publish-workflow.md` <!-- d-check:status-provenance -->
§4 (Live-Prüfung gegen NuGet, dritte unabhängige Bestätigung) · Commit `38a137f7`.
