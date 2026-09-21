**Vorgang:** slice-sdk-kotlin-publish-workflow

**Fund:** Review-Finding F-1 (HIGH,
`docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` <!-- d-check:status-provenance -->): Im
ursprünglichen Implementer-Commit (`546a7503`) lief `./gradlew publish`
direkt auf dem GitHub-hosted Runner statt, wie `ADR-0109` §Entscheidung
Festlegung 5 wörtlich verlangt, im gepinnten `eclipse-temurin:21-jdk`-Image
— eine echte Rebuild-/Codegenerierungs-Duplikation außerhalb von Docker
(dieselbe Bau-Kette, die `make sdk-pack-kotlin` bereits vollständig im
Docker-Image durchläuft, lief ein zweites Mal auf dem Runner, mit frischem,
ungecachtem Dependency-Bezug), nicht bloß eine CLI-Stil-Frage. Der
Slice-Plan (§3 „Beim Schreiben getroffene Entscheidungen") begründete die
Abweichung mit einer eigenen Lesart von `AGENTS.md` §3.1 („bindet die
Host-Toolchain-Sperre an lokale Entwicklung, nicht an den ephemeren
GitHub-hosted Runner"), ohne sie als offene Frage zu kennzeichnen oder eine
Supersede-ADR zu erwägen — der Reviewer bewertete das ausdrücklich als
„plausible Position", die aber „im Widerspruch zu dem [steht], was
`ADR-0109` als Accepted-Entscheidung bereits wörtlich festgelegt hat, ohne
dass eine Folge-ADR mit `Supersedes ADR-0109` diese Änderung trägt".

Aufgelöst über eine Fixrunde (`21872c3d`): eine neue Docker-Stufe
`publish` in `sdks/kotlin/Dockerfile` (baut auf `build` auf) trägt den
Publish-Schritt jetzt Docker-only; der Fixrunden-Review
(`docs/reviews/review-slice-sdk-kotlin-publish-workflow-fixrunde.md` <!-- d-check:status-provenance -->)
bestätigte real über einen eigenen `docker build --target publish` +
`docker inspect`-Lauf (`Cmd=["./gradlew","--no-daemon","publish"]`), dass
`./gradlew publish` nicht mehr auf dem Runner läuft — 0 HIGH/MEDIUM
verbleibend, 1 INFO (Detailschärfung des ohnehin bereits offenen
`AGENTS.md` §3.10-Post-Push-Risikos). Kein `Supersedes ADR-0109` nötig —
der Fund betraf den Ausführungsort, keine inhaltliche ADR-Korrektur, genau
die Unterscheidung, die `ADR-0109` §Re-Evaluierungs-Trigger 4 bereits
selbst zieht.

Quelle: Review zu `slice-sdk-kotlin-publish-workflow` (F-1) · Fixrunden-
Commit `21872c3d` · Anlage durch den Planner bei der Slice-Closure.
