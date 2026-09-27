**Vorgang:** slice-harness-mutationsbild-und-verweigerte-aktion (Review H-1; Verifikation §1/§5)

**Fund:** Die Zusage „`SRC` wird vor dem Aufruf zu einem absoluten Pfad aufgelöst“
(`realpath -- "$src"` in `tools/harness/image-mutation.sh`, der einzige Schutz gegen
einen relativen `SRC`, der zufällig mit der Repo-Wurzel identisch ist) war ohne
Testfall an ihrer Eingabeseite gebunden: jeder Fall des Tabellentests übergab `SRC`
bereits als absoluten Pfad. Der Reviewer erzeugte eine Kopie ohne die `realpath`-Zeile
(`src_abs="$src"`) — der ganze Tabellentest blieb grün — und fuhr den Exploit real: mit
`SRC=.` aus einem Wegwerf-Git-Repo baute die mutierte Fassung tatsächlich ein Image aus
dem Arbeitsverzeichnis, während die unveränderte Fassung mit „`SRC '.' ist die
Repo-Wurzel oder liegt unter ihr“ ablehnte. Die Fixrunde ergänzte drei Testfälle
(relativer Pfad `.`, relativer Pfad `sub`, Symlink auf die Repo-Wurzel); der Verifier
fuhr denselben Exploit gegen den erweiterten Tabellentest erneut — Exit 1, exakt die
drei neuen Fälle rot, alle 27 übrigen grün. Schwere HIGH (echter Arbeitsbaum-Bau
möglich, `AGENTS.md` §3.1): eine Datei trotz Deckel des verkörperten Eintrags.

Quelle: `docs/reviews/review-slice-harness-mutationsbild-und-verweigerte-aktion.md` (H-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-harness-mutationsbild-und-verweigerte-aktion.md` (§1, §5). <!-- d-check:status-provenance -->
