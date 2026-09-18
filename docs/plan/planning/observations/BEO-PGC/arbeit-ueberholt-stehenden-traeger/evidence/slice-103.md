**Vorgang:** slice-103

**Fund (a) — direkt behoben:** Die Einführung der vierten Programmspur je
Sprache auf der Kotlin-Seite (`examples/kotlin/grpc-client`, eine neue
`runtime-grpc`-Stufe neben den unveränderten `runtime`-/`runtime-sse`-/
`runtime-nats`-Stufen) machte einen Satz in `harness/README.md` §Sensors
falsch, der nicht im Diff stand: die `make examples-kotlin`-Zeile beschrieb
seit `slice-101` „drei Images" je Sprache — bei Niederschrift wahr, seit
diesem Slice falsch, weil jetzt vier Images gebaut werden.

Der Implementer-eigene §3.13-Suchlauf fand den Satz und korrigierte ihn im
selben Commit (`111f1cb`) auf „vier Images" mit dem benannten Zusatzkontext
erwähnt. Reviewer und Verifier bestätigen die Korrektur unabhängig als
vollständig (Verifikationsbericht zu `slice-103`, #10). Wie bei
`slice-093`/`slice-094`/`slice-100`/`slice-101`/`slice-102` ist dies ein
Fall, in dem der vorgeschriebene Suchlauf selbst den Fund lieferte.

**Fund (b) — gemeldet, nicht geändert:** `slice-097` §1 (`done/`, immutabel)
beschreibt den C#-/Kotlin-Weg als „ihre Bau-Kontexte erreichen sie heute
nicht". Diese Aussage war bei Niederschrift von `slice-097` wahr. Mit
`slice-102` wurde sie **teilweise** überholt (für C# seither falsch, für
Kotlin weiterhin richtig, siehe `evidence/slice-102.md`). Mit **diesem**
Slice ist sie jetzt auch für Kotlin überholt — die Aussage ist damit für
**beide** in `slice-097` genannten Sprachen falsch geworden, weil die volle
Matrix mit `slice-103` vollständig geliefert ist. Kein Zitat-Korrektur-Fall
nach `ADR-0073` (die Aussage trug ihren Ursprung korrekt und war bei
Niederschrift wahr) und keine Editier-Gelegenheit — `AGENTS.md` §3.13
verlangt hier ausdrücklich **Melden, nicht Ändern**, weil `slice-097` als
geschlossener Vorgang in `done/` liegt.

**Ein Vorgang, eine Gelegenheit:** Beide Funde entstammen demselben Vorgang
(`slice-103`) und zählen deshalb als **eine** Evidenzdatei, auch wenn sie
zwei unterschiedliche Träger betreffen und zwei unterschiedliche
Behandlungsformen verlangen (direkt behoben vs. gemeldet) — dieselbe
Zählregel wie bei `slice-102`, `slice-093` und `slice-096`.

Quelle: Implementer-Commit `111f1cb` (§3.13-Suchlauf, Fund a) ·
Review zu `slice-103` (Negativbefund „`harness/README.md`") ·
Verifikationsbericht zu `slice-103` (#10) · Planner-Closure-Sichtung (Fund
b, Vervollständigung von `evidence/slice-102.md` Fund b) ·
`docs/plan/planning/done/slice-097-umzug-vertragsflaeche.md:107` ·
`AGENTS.md` §3.13.
