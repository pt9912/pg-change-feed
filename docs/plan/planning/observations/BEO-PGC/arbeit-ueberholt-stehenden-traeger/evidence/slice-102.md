**Vorgang:** slice-102

**Fund (a) — direkt behoben:** Die Einführung der vierten Programmspur je
Sprache (`examples/csharp/grpc-client`, eine neue `runtime-grpc`-Stufe neben
den unveränderten `runtime`-/`runtime-sse`-/`runtime-nats`-Stufen) machte
einen Satz in `harness/README.md` §Sensors falsch, der nicht im Diff stand:
die `make examples-csharp`-Zeile beschrieb seit `slice-101` „drei Images" je
Sprache — bei Niederschrift (`slice-101`) wahr, seit diesem Slice falsch,
weil jetzt vier Images gebaut werden.

Der Implementer-eigene §3.13-Suchlauf fand den Satz und korrigierte ihn im
selben Commit (`a71b425`) auf „vier Images" mit dem benannten Zusatzkontext
erwähnt. Reviewer und Verifier bestätigen die Korrektur unabhängig als
vollständig (Review zu `slice-102`, Negativbefund
„`harness/README.md`" · Verifikationsbericht zu `slice-102`, #9). Wie bei
`slice-093`/`slice-094`/`slice-100`/`slice-101` ist dies ein Fall, in dem der
vorgeschriebene Suchlauf selbst den Fund lieferte.

**Fund (b) — gemeldet, nicht geändert:** `slice-097` §1 (`done/`, immutabel)
beschreibt den C#-/Kotlin-Weg als „ihre Bau-Kontexte erreichen sie heute
nicht". Diese Aussage war bei Niederschrift von `slice-097` wahr und wird
durch `slice-102` **teilweise** überholt: Für C# stimmt sie seit diesem
Slice nicht mehr, für Kotlin (bis `slice-103` geliefert ist) weiterhin. Das
ist kein Zitat-Korrektur-Fall nach `ADR-0073` — die Aussage trug ihren
Ursprung korrekt und war zum Zeitpunkt ihrer Niederschrift wahr — und keine
Editier-Gelegenheit: `AGENTS.md` §3.13 verlangt hier ausdrücklich **Melden,
nicht Ändern**, weil `slice-097` als geschlossener Vorgang in `done/` liegt.

**Ein Vorgang, eine Gelegenheit:** Beide Funde entstammen demselben Vorgang
(`slice-102`) und zählen deshalb als **eine** Evidenzdatei, auch wenn sie
zwei unterschiedliche Träger betreffen (Code-/Harness-Doku vs. eine
immutable Planning-Datei) und zwei unterschiedliche Behandlungsformen
verlangen (direkt behoben vs. gemeldet). Das ist dieselbe Zählregel wie bei
`slice-093` (vier Stellen in zwei Dateien, eine Gelegenheit) und `slice-096`
(zwei ADRs, eine Runde).

Quelle: Review zu `slice-102` (Negativbefund
„`harness/README.md`") · Verifikationsbericht zu `slice-102` (#9) ·
`docs/plan/planning/done/slice-097-umzug-vertragsflaeche.md:107` ·
`AGENTS.md` §3.13.
