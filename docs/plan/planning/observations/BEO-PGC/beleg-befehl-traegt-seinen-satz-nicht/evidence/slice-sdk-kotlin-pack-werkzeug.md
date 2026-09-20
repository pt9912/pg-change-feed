**Vorgang:** slice-sdk-kotlin-pack-werkzeug (Review-Fund F-2, Planner-Closure)

**Fund:** `sdks/kotlin/Dockerfile`, `harness/mk/sdk.mk`, `harness/README.md`
und die Commit-Message behaupten wortgleich, GitHub Packages „verlangt laut
offizieller Dokumentation … keines" (Sources-/Javadoc-Jar) und zitieren dafür
`docs.github.com/en/actions/publishing-packages/publishing-java-packages-with-gradle`.
Ein eigener `curl`-Abruf dieser Seite (Reviewer, danach unabhängig ein
zweites Mal vom Verifier wiederholt) zeigt: Die Seite erwähnt
Sources-/Javadoc-Jars an **keiner** Stelle — weder fordernd noch
freistellend. Der zitierte Beleg **schweigt** zur behaupteten Aussage,
er widerspricht ihr nicht. Der praktische Schluss bleibt trotzdem haltbar
(gestützt auf ein zweites, tatsächlich tragendes, aber nicht als solches
benanntes Faktum: Gradles `from(components["java"])` ohne
`withSourcesJar()`/`withJavadocJar()` packt nur den Haupt-Jar,
`docs.gradle.org` bestätigt das real) — die Einstufung bleibt deshalb
MEDIUM statt HIGH. Erster Beleg dieser Klasse an einer **Web-Recherche**
statt einem Befehl/einer Codestelle/einer Adresse/einer Assertion (die
vier bislang belegten Formen).

**Ausgang:** kein Fix in diesem Slice (0 HIGH, nicht merge-blockierend,
Risiko strukturell durch `ADR-0109` §Re-Evaluierungs-Trigger 4 aufgefangen)
— für `slice-sdk-kotlin-publish-workflow` vorgemerkt: die Attribution
präzisieren oder auf das tragende Gradle-Kern-Faktum umstellen, falls der
Zug die Formulierung ohnehin berührt.
