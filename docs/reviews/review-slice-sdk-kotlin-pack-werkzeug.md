# Review-Report: slice-sdk-kotlin-pack-werkzeug — 2026-09-20

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-pack-werkzeug.md`),
`ADR-0109` (Accepted, Festlegung 5, §Konsequenzen Folgepflicht 2/3) sowie
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich —
das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Commit `84293e14` (`feat(sdk): make sdk-pack-kotlin +
Pflichtenheft-Traeger-Nachzug`) gegen Elternstand `9c4efd42`
(`git diff 9c4efd42..84293e14`), Slice `slice-sdk-kotlin-pack-werkzeug`,
Welle `welle-sdk-kotlin-lh-fa-sst-009`. Ein Commit, sieben geänderte
Dateien: `harness/mk/sdk.mk` (neues `sdk-pack-kotlin`-Target),
`sdks/kotlin/Dockerfile` (neue Stufen `pack`/`pack-export`),
`sdks/kotlin/.gitignore` (`dist/`-Eintrag), `tools/harness/sdk-pack-kotlin.sh`
(neu), `spec/pflichtenheft.md` (§1 `LH-FA-SST-009.a`-Nachzug, §6
`SPEC-028`, §7 Historie), `harness/README.md` (§Werkzeuge-Zeile),
Slice-Plan (DoD-Häkchen).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(HEAD zum Review-Zeitpunkt) · **Modell:** claude-sonnet-5 · **Datum:** 2026-09-20

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-kotlin-pack-werkzeug.md` (§1, §2 DoD, §3 Plan, §6 Risiken)
- `ADR-0109` (Accepted) — Festlegung 1/3/5, §Kontext Recherche, §Konsequenzen Folgepflicht
- `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a`, §6 Externe Verträge, §7 Historie
- `harness/README.md` §Sensors/§Werkzeuge, `harness/conventions.md` (MR-000/MR-001)
- Formvorbilder: `harness/mk/sdk.mk` (`sdk-pack-csharp`/`sdk-pack-python`), `sdks/csharp/Dockerfile`,
  `sdks/python/Dockerfile`, `tools/harness/sdk-pack-csharp.sh`, `tools/harness/sdk-pack-python.sh`,
  `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md`, `docs/reviews/review-slice-sdk-python-pack-werkzeug.md`
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13, §4

---

## Findings

### F-1 — C#-Absatz-Schlusszeile bleibt stale, obwohl der Nachbar-Absatz jetzt zum dritten Mal wächst

- `kategorie`: LOW
- `quelle`: Maintainability (`spec/pflichtenheft.md` §1 `LH-FA-SST-009.a`)
- `pfad`: `spec/pflichtenheft.md:182-184` (C#-Absatz: „Eine zweite Sprache
  oder ein zweiter Vertriebsweg bleibt offen — diese Kennung bleibt ihre
  Adresse.")
- `befund`: Der C#-Absatz behauptet weiterhin, eine zweite Sprache sei
  offen — real ist inzwischen die dritte (Kotlin, dieser Diff) beantwortet.
  Die Staleness selbst entstand bereits mit dem Python-Nachzug (vor diesem
  Diff, nicht durch diesen Zug verursacht) und wurde in keinem der beiden
  vorangegangenen Reviews (`review-slice-sdk-python-pack-werkzeug.md`)
  benannt. Der Slice-Plan dieses Zuges (§3, Zeile zu
  `spec/pflichtenheft.md` §1) benennt das Problem explizit und entscheidet
  bewusst, es nicht zu beheben („außerhalb des DoD-Umfangs") — eine
  begründete, aber nicht risikofreie Abgrenzung: mit jeder weiteren Sprache
  wird der stehen gebliebene Satz für einen Abschnittsweise-Leser
  irreführender.
- `verifizierbar`: nein — kein Sensor prüft Prosa-Kohärenz zwischen zwei
  Absätzen derselben Pflichtenheft-Sektion.
- `klasse`: „Nachzug lässt Nachbarabsatz stale" (Geschwister-Fall zu
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, dort bislang nur für
  Slice-Plan-Dokumente belegt, nicht für `spec/pflichtenheft.md`
  Fließtext — dies wäre der erste Beleg an dieser Trägerklasse).

### F-2 — Sources-/Javadoc-Jar-Freistellung: zitierter Beleg bestätigt die Aussage nicht explizit, sondern schweigt dazu

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz B (Herkunft von Aussagen — Beleg-Anker)
- `pfad`: `sdks/kotlin/Dockerfile:74-80`, `harness/mk/sdk.mk:47-49`,
  `harness/README.md` (neue `sdk-pack-kotlin`-Zeile), Commit-Message
- `befund`: Alle vier Stellen behaupten wortgleich, GitHub Packages
  „verlangt laut offizieller Dokumentation … keines" (Sources-/Javadoc-Jar).
  Eigener `curl`-Abruf der genannten Seite
  (`docs.github.com/en/actions/publishing-packages/publishing-java-packages-with-gradle`)
  sowie zweier verwandter Registry-Seiten
  (`.../working-with-the-gradle-registry`) zeigt: Der Text erwähnt
  Sources-/Javadoc-Jars an **keiner** Stelle — weder als Pflicht noch als
  Freistellung. Die zitierte Seite belegt die *Mechanik* (kein
  Dritt-Plugin, `GITHUB_TOKEN`, `GitHubPackages`-Repository-Block —
  diese Teile sind wörtlich korrekt und real bestätigt) und **nicht** die
  Freistellungs-Aussage selbst. Der Schluss ist plausibel und stützt sich
  zusätzlich auf ein real bestätigtes, aber getrenntes Gradle-Kern-Faktum
  (`docs.gradle.org/current/userguide/publishing_maven.html`:
  `from(components["java"])` ohne `withSourcesJar()`/`withJavadocJar()`
  packt nur den Haupt-Jar) — dieses zweite Faktum wird in Dockerfile/
  Commit-Message aber nicht als die tragende Stütze benannt, sondern der
  GitHub-Seite zugeschrieben, die dazu schweigt. Das Risiko ist durch
  `ADR-0109` §Re-Evaluierungs-Trigger 4 bereits strukturell aufgefangen
  (ein realer Publish-Fehlschlag löst eine Folge-ADR aus) — deshalb MEDIUM
  statt HIGH trotz Nähe zur Skill-Klasse „Beleg trägt seinen Satz nicht".
- `verifizierbar`: nein — nur durch eigenen Web-Abruf/Lesen der zitierten
  Stelle, kein Gate.
- `klasse`: „Beleg unterstützt Aussage nur durch Schweigen, nicht durch
  Bestätigung" (Nachbar-Form zu `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`,
  hier zum ersten Mal an einer Web-Recherche statt einem Befehl/einer
  Codestelle).

---

## Negativbefunde

- geprüft, ohne Befund: `make sdk-pack-kotlin` (`harness/mk/sdk.mk`) —
  Docker-only (`@bash tools/harness/sdk-pack-kotlin.sh`), kein
  `GATE_CHECKS`-Eintrag (`grep -rn GATE_CHECKS Makefile harness/mk/*.mk`
  zeigt `sdk.mk` nicht in der Liste), Target alphabetisch zwischen
  `sdk-pack-csharp`/`sdk-pack-python` eingefügt, identisches
  `.PHONY`/Doku-Kommentar-Muster.
- geprüft, ohne Befund: roter Test bricht den Bau ab — eigene
  Mutations-Probe (`assertEquals("test-token", …)` →
  `assertEquals("wrong-token", …)` in `PgChangeFeedClientOptionsTest.kt`,
  danach zurückgesetzt, `git status` nach Rücksetzung sauber):
  `bash tools/harness/sdk-pack-kotlin.sh` bricht mit echtem Exit-Code `1`
  ab (ungepiped gemessen, `AGENTS.md` §3.9), `./gradlew test` schlägt am
  mutierten Assert fehl, `docker build` erreicht die `pack`-Stufe nicht.
- geprüft, ohne Befund: reales `.jar` — eigener Lauf von
  `bash tools/harness/sdk-pack-kotlin.sh` (nach `rm -rf sdks/kotlin/dist`)
  erzeugt `sdks/kotlin/dist/pgchangefeed-kotlin-0.1.0.jar`, 112591 Bytes
  (`ls -la`), `unzip -l` zeigt ein valides Jar-Archiv mit
  `META-INF/MANIFEST.MF` und den erwarteten `cdc/stream/v1/*.class`-Einträgen
  — deckungsgleich mit der in Commit-Message/Dockerfile-Kommentar genannten
  Byte-Größe, kein Drift (`AGENTS.md` §3.12 Instanz A).
- geprüft, ohne Befund: `--build-context proto=proto` — sowohl
  `sdks/kotlin/Dockerfile` (`COPY --from=proto …`) als auch
  `tools/harness/sdk-pack-kotlin.sh` (`docker build --build-context
  proto=proto …`) tragen den zusätzlichen Bau-Kontext; eigener Bau-Lauf
  bestätigt, dass die Stufe real durchläuft.
- geprüft, ohne Befund: Export-Mechanik `pack-export` — identisches
  `tar -cf - -C /out .`-Entrypoint- und
  `docker run --rm --network none <image> | tar -x -C …`-Extraktionsmuster
  wie `sdks/csharp/Dockerfile`/`sdks/python/Dockerfile` und
  `sdk-pack-csharp.sh`/`sdk-pack-python.sh`; `tools/harness/sdk-pack-kotlin.sh`
  ist strukturell nahezu Zeile-für-Zeile identisch mit
  `sdk-pack-csharp.sh` (Namen ausgetauscht), inklusive
  `set -euo pipefail`, `repo_root`-Ermittlung über `git rev-parse
  --show-toplevel` und der Pipe-Disziplin-Begründung im Kommentarkopf
  (`AGENTS.md` §3.9).
- geprüft, ohne Befund: `spec/pflichtenheft.md` §1
  `LH-FA-SST-009.a`-Nachzugsatz — der neue Kotlin-Absatz schließt korrekt
  an den (jetzt bereinigten) Python-Absatz an, „Eine vierte Sprache …
  bleibt offen" ist zum Zeitpunkt dieses Diffs sachlich richtig; die vom
  Implementer behauptete Entfernung der stale gewordenen
  Python-Schlusszeile („Eine dritte Sprache … bleibt offen") ist real
  vollzogen (`git diff` zeigt die Zeile entfernt, kein Wortlaut-Widerspruch
  zurückgelassen) — notwendig und korrekt ausgeführt
  (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`). Siehe F-1 für den
  abgegrenzten, nicht behobenen C#-Nachbarfall.
- geprüft, ohne Befund: kein `ADR-0109`-Verweis im Pflichtenheft-Fließtext
  (`grep -n "ADR-" spec/pflichtenheft.md` zeigt an den neuen Zeilen keinen
  Treffer) — `matrix`-Modul-Konformität (`spec → adr` verboten) gewahrt,
  wie im Slice-Plan §2 gefordert.
- geprüft, ohne Befund: `SPEC-028`-Kollision — `grep -rn "SPEC-028"
  spec/*.md docs/` zeigt vor diesem Diff keinen Treffer außerhalb des
  Slice-Plans selbst; Nummer real neu vergeben, kein Parallel-Zug hat sie
  bereits belegt.
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — neue Zeile
  korrekt zwischen `sdk-pack-csharp`/`sdk-pack-python` positioniert
  (alphabetisch/thematisch konsistent), Bindung auf `ADR-0109` Festlegung 5,
  Format identisch mit den beiden Nachbarzeilen (Bau-Kontext, Stufen,
  Export-Mechanik, `.gitignore`-Hinweis, „kein Gate"-Begründung).
- geprüft, ohne Befund: Out-of-Scope-Einhaltung (§1 des Slice-Plans) —
  kein `.github/workflows/`-Diff, kein `docs/user/version.md`-Diff, kein
  `publish`-Aufruf in Dockerfile/Skript/Makefile, keine Secret-Anlage, kein
  Eingriff in `sdks/kotlin/pgchangefeed-kotlin/src/**` (nur Bau-Infrastruktur
  geändert) — `git diff --stat` bestätigt exakt die sieben im Plan
  vorgesehenen Dateien.
- geprüft, ohne Befund: root `README.md`/`README.de.md`/
  `docs/user/releasing.md` bewusst nicht angefasst — Präzedens real gegen
  `git log` verifiziert: `1b3b4009`/`3be590f9` (README-Erwähnung bzw.
  Releasing-Nachzug) liefen jeweils **nach** der realen Veröffentlichung
  auf NuGet.org/PyPI, nicht im jeweiligen `sdk-pack-*-werkzeug`-Slice
  selbst (`0e42b801`, der Python-Analogcommit, ließ dieselben Dateien
  unangetastet). Die Begründung „Kotlin noch nicht real veröffentlicht"
  trägt.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — kein
  Konjunktiv über eine verworfene Alternative, kein abwesender Text.
  Slice-ID-Nennung in Klammern (`sdks/kotlin/Dockerfile:19-24,68-70`,
  `harness/mk/sdk.mk:31`) folgt exakt dem bereits etablierten,
  vorab-reviewten Muster in `sdks/csharp/Dockerfile`/`harness/mk/sdk.mk`
  (dort in `review-slice-sdk-csharp-pack-werkzeug.md` bereits als
  „etablierte Herkunfts-Anker-Form" akzeptiert, nicht als
  Slice-/Wellen-Chronik-Verstoß) — kein neuer Verstoß, keine Verschärfung.
- geprüft, ohne Befund: JDK-Basis-Digest — `eclipse-temurin:21-jdk@sha256:085eb9…`
  in `sdks/kotlin/Dockerfile` ist byte-identisch mit dem bereits vorher
  (Projektgerüst-Slice) gepinnten Digest in `sdks/kotlin/Dockerfile` selbst
  und mit `examples/kotlin/Dockerfile` — kein neuer, in diesem Diff
  eingeführter Pin.
- geprüft, ohne Befund: Backtick-Parität — alle drei von diesem Diff
  berührten Markdown-Dateien (`docs/plan/planning/in-progress/slice-sdk-kotlin-pack-werkzeug.md`,
  `harness/README.md`, `spec/pflichtenheft.md`) haben eine gerade Anzahl
  Backticks im committeten Endstand (362/1690/1454 — je `% 2 == 0`).
- geprüft, ohne Befund: `make gates` — eigener Lauf, ungepiped, Exit-Code
  direkt geprüft: `0` (grün), inklusive `generated-sync`/`a-check`
  (`gesamt: 0 Befund(e)`).
- geprüft, ohne Befund: DoD-Checkboxen dieses Diffs — die neun mit `[x]`
  markierten Zeilen entsprechen real vollzogenen Schritten (Target
  existiert, Recherche dokumentiert, `.jar` real erzeugt,
  Pflichtenheft/`harness/README.md`-Nachzüge vorhanden, `make gates`
  grün); „Reconciliation-Register entfällt" stimmt (`find … -iname
  reconciliation*` liefert keinen Treffer im Repo).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Nachzug lässt Nachbarabsatz stale ·
Beleg unterstützt Aussage nur durch Schweigen, nicht durch Bestätigung

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. Die eine MEDIUM-Finding (F-2)
betrifft eine Doku-/Belegpräzisions-Frage ohne Auswirkung auf Bau, Test
oder das reale Erzeugnis dieses Slice (kein Publish in diesem Umfang,
Risiko bereits über `ADR-0109` §Re-Evaluierungs-Trigger 4 abgefangen); die
LOW-Finding (F-1) ist vorbestehende, nicht durch diesen Diff verursachte
Textschuld. Beide sind für den Implementer dieses Slice keine
Merge-Blocker; F-2 ist am ehesten beim nächsten Zug relevant, der die
Formulierung berührt oder real veröffentlicht (`slice-sdk-kotlin-publish-workflow`),
F-1 beim nächsten Nachzug, der C#s Absatz ohnehin anfasst.

**Übergabe:** Keine Fixrunde am Implementer erforderlich (0 HIGH, beide
Findings sind für spätere Züge vorgemerkt, kein Reviewer→Implementer-
Rückgabe-Pfeil für diesen Slice). Gemäß Skill-Abschnitt „DoD-Checkbox-
Nachzug ohne Fixrunde" zieht dieser Report die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan im
selben Commit nach, der diesen Report anlegt. Die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Steering-Loop-
Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt keine Verifikation
(Modul 11, separater Kontext).
