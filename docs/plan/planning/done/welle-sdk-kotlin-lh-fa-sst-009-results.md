# Welle welle-sdk-kotlin-lh-fa-sst-009 — Closure-Notiz

**Welle:** welle-sdk-kotlin-lh-fa-sst-009
**Abschluss:** 2026-09-21
**Verantwortlich:** pt9912

## Was wurde geliefert?

[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(`Accepted`) vollständig umgesetzt — das dritte offizielle
Client-Bibliothek-Package für
[`LH-FA-SST-009`](../../../../spec/lastenheft.md): Kotlin, GitHub Packages, in
fünf Slices, mit derselben Zwei-Flächen-Parallelisierungs-Form wie die
C#-Welle (HTTP + gRPC in einem Erst-Release):

- **`slice-sdk-kotlin-projektgeruest`**: `sdks/kotlin/pgchangefeed-kotlin/`
  neu angelegt — Gradle-Projekt (`build.gradle.kts`, Koordinate
  `io.github.pt9912:pgchangefeed-kotlin`, `version = "0.1.0"`),
  eigenständiges, digest-gepinntes Docker-Bau-Setup
  (`eclipse-temurin:21-jdk`, JDK 21/Gradle 8.14 wie `examples/kotlin/`),
  englischsprachiges `README.md`, ein minimales Konfigurations-Skelett.
- **`slice-sdk-kotlin-http-client-flaeche`** und
  **`slice-sdk-kotlin-grpc-client-flaeche`**: öffentliche Kotlin-API-Fläche
  für alle neun Port-gedeckten Fähigkeiten von `SPEC-018` plus das
  Changes-Lesen (`SPEC-022`), und `StreamChanges` über gRPC (`SPEC-020`,
  `.proto`-Bezug über den benannten Zusatzkontext `proto=proto`,
  Metadata-Auth) — parallelisierbar geschnitten, beide hängen nur vom
  Projektgerüst ab.
- **`slice-sdk-kotlin-pack-werkzeug`**: `make sdk-pack-kotlin` (Docker-only,
  kein Gate, `./gradlew build`/`test`) baut, testet und paketiert beide
  Flächen zu einem realen `pgchangefeed-kotlin-0.1.0.jar` — Träger-Nachzug
  in `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 (neue
  `SPEC-028`-Zeile) sowie `harness/README.md` §Werkzeuge.
- **`slice-sdk-kotlin-publish-workflow`**: `.github/workflows/sdk-kotlin-release.yml`
  (Trigger `push: tags: ['sdk-kotlin-v*']`, eigener Tag-Namensraum),
  Tag-vs-`build.gradle.kts`-Versionsabgleich (dritter Konsument von
  `tools/harness/semver-regex.sh`), `./gradlew publish` gegen
  `https://maven.pkg.github.com/pt9912/pg-change-feed` — authentifiziert
  ausschließlich über das eingebaute `GITHUB_TOKEN`
  (`permissions: contents: read` / `packages: write`), **kein** externes
  Repository-Secret; `docs/user/releasing.md`-Nachzug im selben Slice
  (vorab geplant). Eine reale Reviewer-Fixrunde korrigierte den
  Publish-Schritt von einem Runner-Aufruf zu einer eigenen, gepinnten
  Docker-Stufe `publish` — siehe „Was ging anders als geplant".

`make gates` grün auf dem Endstand (874 Dateien, 0 `docs-check`-Befunde,
Coverage 82,70 % ≥ 80 %, `a-check`/`generated-sync`/`commit-traceability`/
`baseline-verify` je ohne Befund). Ein real gebautes
`sdks/kotlin/dist/pgchangefeed-kotlin-0.1.0.jar` existiert, frisch für
diese Closure mit einem eigenständigen `make sdk-pack-kotlin`-Lauf
(Exit 0) erneut erzeugt — nicht nur ein grüner `docker build`, das
Artefakt selbst liegt vor und trägt die erwartete Version `0.1.0`. Kein
realer `sdk-kotlin-v*`-Tag wurde gepusht (`git tag -l "sdk-kotlin-v*"`
leer) — diese Welle liefert bewusst nur den **Mechanismus**; anders als
bei den beiden Vorgänger-SDKs entfällt hier das Secret-Anlage-Risiko
strukturell (`GITHUB_TOKEN` ist immer verfügbar, kein
`gh secret list`-Bezug nötig).

## Was hat funktioniert?

Über den gesamten Wellen-Zyklus hinweg trug eine durchgehende, mehrfach
unabhängige Verifikation:

- `slice-sdk-kotlin-projektgeruest`: 0 HIGH/MEDIUM, 3 LOW — keine
  Fixrunde, ein Genus-Tippfehler (F-3) direkt vom Planner bei der Closure
  korrigiert.
- `slice-sdk-kotlin-http-client-flaeche`: 1 HIGH (F-1, KDoc-Aussage über
  Kotlin-`internal` als vermeintliche JVM-Zugriffsbeschränkung) — real
  behoben (`9181ae21`), Verifikation bestätigt DoD-Konformität.
- `slice-sdk-kotlin-grpc-client-flaeche`: 0 HIGH/MEDIUM/LOW/INFO auf
  Anhieb — der sauberste Zyklus der Welle, keine Fixrunde.
- `slice-sdk-kotlin-pack-werkzeug`: 0 HIGH, 1 MEDIUM (F-1, Träger-Nachzug-
  Lücke), 1 LOW (F-2, Beleg-Schweigen statt Beleg-Aussage) — beide
  Findings direkt dem Beobachtungs-Register zugeführt statt neu
  formuliert, kein neuer Fixrunden-Zyklus nötig.
- `slice-sdk-kotlin-publish-workflow`: 1 HIGH (F-1) — die einzige echte
  Fixrunde dieser Welle, mit einem **zweiten, unabhängigen** Reviewer-Lauf
  (kein Self-Review, keine Übernahme der Implementer-Einschätzung): Der
  erste Reviewer-Lauf empfahl explizit **keine** einseitige
  Implementer-Rückgabe, sondern eine Fixrunde mit frischem Zweit-Review —
  genau das lief so, der Fix (`21872c3d`) löste das Finding real auf, der
  Fixrunden-Reviewer bestätigte über einen eigenen `docker build --target
  publish` + `docker inspect`-Lauf (0 HIGH/MEDIUM, 1 INFO), der Verifier
  wiederholte das ein drittes Mal unabhängig (vierter Docker-Build
  insgesamt). Dieses 3-Schritt-Muster (Review → Fix → Fixrunden-Review)
  trug real, ohne dass eine Rolle die Einschätzung der vorherigen
  unbesehen übernahm.

Das etablierte 3-Commit-Move-Muster (`git mv` · Inhalt · `git mv`) hielt
`AGENTS.md` §3.3 über alle fünf Slices sauber. Die parallelisierbare
Zwei-Flächen-Form (HTTP + gRPC) funktionierte wie bei der C#-Welle ohne
gegenseitige Blockade.

## Was ging anders als geplant?

- **`slice-sdk-kotlin-publish-workflow` trug eine Implementierungs-
  Abweichung von einer bereits `Accepted`-ADR, nicht nur einen
  Textdefekt.** Der Slice-Plan sah beim Schreiben (§3) den Publish-Schritt
  bewusst auf dem GitHub-hosted Runner vor, mit einer eigenen Lesart von
  `AGENTS.md` §3.1 als Begründung, ohne das als offene Frage zu
  kennzeichnen — entgegen `ADR-0109` §Entscheidung Festlegung 5s
  wörtlichem Docker-only-Wortlaut. Anders als bei den Fixrunden der beiden
  Vorgänger-SDK-Wellen (C#: Slice-Chronik im Docstring; Python:
  Träger-Nachzug-Lücke) ist dies eine **neue Fehlerklasse** für das
  Beobachtungs-Register — festgehalten als neuer Eintrag
  `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (1×), nicht als
  `AGENTS.md`-Regelschärfung, weil ein Erstauftreten unter der
  3×-Schwelle liegt. Die Auflösung blieb ein reiner Code-Fix, keine
  Supersede-ADR — `ADR-0109` §Re-Evaluierungs-Trigger 4 hatte diesen
  Unterschied (Workflow-Bugfix ohne Entscheidungsänderung vs. inhaltliche
  Korrektur) bereits selbst vorgesehen, bevor der Fall real eintrat; siehe
  `slice-sdk-kotlin-publish-workflow.md` §7.
- **`slice-sdk-kotlin-projektgeruest` fand real drei Instanzen derselben
  wortgleichen README-Formulierungslücke über alle drei SDK-Sprachpakete
  hinweg** (`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`,
  rückwirkend bei 3× angelegt — siehe „Steering-Loop-Einträge" unten).
- **`slice-sdk-kotlin-pack-werkzeug` fand eine zweite Instanz einer bereits
  bekannten Nachzug-Klasse** (`BEO-PGC/nachzug-laesst-ueberholten-text-
  stehen`, 1× → 2×, weiter unter der Schwelle) sowie eine weitere Instanz
  der bereits verkörperten `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
  (jetzt 7×, unverändert verkörpert).
- Die beiden aus der Eröffnungs-Sichtung vorab benannten Vermeidungs-
  Strategien griffen real: kein drittes Auftreten von `BEO-PGC/workaround-
  uebersieht-etabliertes-muster-im-bestand` (weiterhin 2×) und kein
  zweites Auftreten von `BEO-PGC/release-mechanismus-nicht-in-releasing-
  doku-nachgezogen` (weiterhin 1×) — beide Vermeidungen waren bereits
  vorab in den jeweiligen Slice-Plänen als eigener DoD-Punkt eingeplant.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — hier stehen die Beobachtungen, die im
Register 3× erreicht haben (Lese-Schritt dieser Closure).

- **`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`** — erreichte
  während dieser Welle real die 3×-Schwelle
  (`slice-sdk-kotlin-projektgeruest`, Beleg über `csharp`/`python`/`kotlin`
  gleichermaßen, rückwirkend erfasst: der Fehler war in allen drei
  README-Dateien real vorhanden, aber nur beim Kotlin-Review zum ersten
  Mal formal benannt). Sein eigenes `state.md` benennt die Closure dieser
  Welle ausdrücklich als „nächste reguläre Lesegelegenheit" — **diese
  Lesung findet hiermit statt.** Die eigentliche Regelschärfung (eine
  Reviewer-Skill-Ergänzung, „prüfe eine von einem Formvorbild übernommene
  README-/Doku-Formulierung auf verbliebene deutsche Fachwörter", oder
  eine Fitness Function) bleibt bewusst **außerhalb** dieser
  Planner-geführten Wellen-Closure: Es ist eine Architect-Entscheidung
  (Modul 4/8), kein Planner-Zug. Dieser Lese-Schritt trägt den Fund
  deshalb als offene Übergabe weiter, statt ihn einseitig zu embodyen —
  `state.md` bleibt entsprechend `offen`, mit dem Vermerk, dass die
  Lesung stattgefunden hat und der nächste Schritt bei einer künftigen
  Architect-Sichtung liegt.
- **`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`** — neu angelegt
  in dieser Welle (`slice-sdk-kotlin-publish-workflow`, 1×), unter der
  3×-Schwelle. Kein Verkörperungsbedarf jetzt — siehe „Was ging anders
  als geplant" oben.
- **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** (bereits verkörpert
  seit `welle-20` in `.harness/skills/reviewer.md`): ein neuer, siebter
  Beleg in dieser Welle (`slice-sdk-kotlin-pack-werkzeug`). Bestätigt:
  bereits verkörpert, keine erneute Verkörperung nötig.
- **`BEO-PGC/report-nackte-id-ohne-link`** (verkörpert seit `slice-063` in
  `AGENTS.md` §3.9): geprüft, kein neuer Beleg in dieser Welle — Zähler
  bleibt real bei 8×.
- **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (verkörpert seit
  `AGENTS.md` §3.13): geprüft, kein neuer Beleg in dieser Welle — Zähler
  bleibt real bei 19×.
- **`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`**
  (bereits verkörpert, kombinierte Verkörperung): geprüft, kein neuer
  Beleg — Zähler bleibt real bei 3×, der Handbuch-Hinweis lag bereits als
  DoD-Punkt in den beiden Flächen-Slices.

**Unter der Schwelle, unverändert oder neu, kein Verkörperungsbedarf:**

- `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` — weiterhin
  2× (`slice-104`, `slice-sdk-csharp-projektgeruest`), kein drittes
  Auftreten in dieser Welle.
- `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` —
  weiterhin 1× (`slice-sdk-csharp-publish-workflow`), kein zweites
  Auftreten — `slice-sdk-kotlin-publish-workflow` nahm den
  `docs/user/releasing.md`-Nachzug bewusst vorab als eigenen DoD-Punkt auf.
- `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` — 1× →
  2× in dieser Welle (`slice-sdk-kotlin-projektgeruest`, KDoc-Wortfragment
  wie bei C#), weiterhin unter der Schwelle.
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` — 1× → 2× in dieser
  Welle (`slice-sdk-kotlin-pack-werkzeug`), weiterhin unter der Schwelle.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert als
  `AGENTS.md` §3.10, Zähler weiterhin 7×) — das Post-Push-Risiko dieser
  Welle bleibt strukturell offen (siehe Trigger-Audit unten), kein neuer
  Beleg, da kein realer Tag-Push stattfand.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/). Ein Eintrag
erreichte in dieser Welle **erstmals** die 3×-Schwelle
(`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`) — der Lese-Schritt
dazu steht oben unter *Steering-Loop-Einträge*; die eigentliche
Regelschärfung bleibt eine offene Übergabe an eine künftige
Architect-Sichtung, kein Teil dieser Wellen-Closure. Ein neuer Eintrag
liegt unter der Schwelle (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`,
1×).

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine belegte
Feststellung.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert `ADR-0109` oder
  eine dieser fünf Slices.
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate`
  steht unverändert bei 82,70 % ≥ 80 %-Schwelle (Endstufe, ausgeschöpft)
  und wurde von keinem der fünf Slices berührt (reiner Kotlin-Baum,
  außerhalb der netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+
  `./gen/...`); `.a-check.yml` `languages: go` liest `sdks/kotlin/**`
  strukturell nicht (real bestätigt: `a-check` meldet 0 Befunde auf dem
  Endstand).
- **Entscheidung/ADR —
  [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  §Re-Evaluierungs-Trigger, durch diese Welle **nicht** ausgelöst:**
  - Trigger 1 (vierte Sprache/vierter Vertriebsweg verlangt) — nicht
    eingetreten, out-of-scope §6 der Welle-Datei.
  - Trigger 2 (real auf GitHub Packages veröffentlicht + Nutzungsdaten
    legen eine andere Priorisierung nahe) — **nicht** eingetreten. Das
    Package ist real paketierbar (`pgchangefeed-kotlin-0.1.0.jar`
    existiert, frisch für diese Closure erneut gebaut), aber **nicht**
    real veröffentlicht: `git tag -l "sdk-kotlin-v*"` ist leer, kein
    Tag-Push erfolgte, `./gradlew publish` lief nie real gegen GitHub
    Packages. Ohne reale Veröffentlichung gibt es keine Download-Zahlen
    und keine Issue-Nachfrage. **Beobachtbare Auslösebedingung**: der
    erste reale `git tag sdk-kotlin-v<SemVer> && git push --tags`, der
    einen grünen `sdk-kotlin-release.yml`-Lauf erzeugt, und danach
    messbare Nutzungssignale (GitHub-Packages-Downloadzahlen,
    Issue-Nachfrage) — prüfbar über `git tag -l "sdk-kotlin-v*"` (nicht
    mehr leer) und einen grünen Actions-Lauf. Bis dahin bleibt `ADR-0109`
    Festlegung 1 (Umfang: HTTP + gRPC, kein SSE/NATS) unverändert in
    Kraft.
  - Trigger 3 (Gradle real auf ≥ 9.1.0 angehoben) — nicht eingetreten,
    `examples/kotlin/`/`sdks/kotlin/` bleiben bei Gradle 8.14/JDK 21.
  - Trigger 4 (der reale GitHub-Packages-Publish-Workflow scheitert im
    ersten realen Lauf an einer Berechtigungs-/Registry-Eigenheit) —
    **nicht** eingetreten im Sinn dieses Triggers: Es gab **keinen**
    realen Post-Push-Lauf in dieser Welle (kein Tag-Push). Das reale
    HIGH-Finding der Fixrunde (`slice-sdk-kotlin-publish-workflow` F-1)
    betraf den Ausführungsort **vor** jedem realen Push, nicht einen
    gescheiterten realen Lauf — genau die vom Trigger-Text selbst
    vorgesehene Unterscheidung („ein reiner Workflow-Bugfix ohne
    Entscheidungsänderung braucht keine neue ADR") traf zu, ohne dass der
    Trigger selbst ausgelöst wurde. Der Post-Push-Nachweis bleibt nach
    `AGENTS.md` §3.10 weiter offen, unabhängig von diesem Trigger.
  - Trigger 5 (Authentifizierungspflicht beim Lesen erweist sich als
    echtes Consumer-Problem) — nicht eingetreten, kein Consumer hat das
    Package real bezogen.

## Nebenbefunde (außerhalb des Scopes dieser Welle, hier nur gemeldet)

- **Kein `sdk-kotlin-v*`-Secret-Bedarf, real bestätigt:** `GITHUB_TOKEN`
  ist der von GitHub Actions automatisch bereitgestellte Workflow-Token —
  anders als bei den beiden Vorgänger-SDK-Wellen entfällt hier jede
  Prüfung eines externen Repository-Secrets strukturell.
  `gh secret list` zeigt zum Zeitpunkt dieser Closure `DOCKERHUB_TOKEN`,
  `DOCKERHUB_USERNAME`, `NUGET_API_KEY`, `PYPI_API_TOKEN` — keines davon
  ist für den Kotlin-Publish-Weg relevant.
- **`PYPI_API_TOKEN` existiert jetzt real** (`gh secret list`,
  2026-09-21) — seit der Python-Welle-Closure extern angelegt; außerhalb
  des Scopes dieser Welle, hier nur zur Vollständigkeit vermerkt (ändert
  nichts an `ADR-0107`s eigenem Re-Evaluierungs-Trigger, der separat von
  dieser Welle geführt wird).

## Folge-Slices

Keine. `ADR-0109` ist mit dieser Welle für Kotlin/GitHub Packages
vollständig umgesetzt. Was ansteht, sind keine weiteren Slices dieser
Welle, sondern eigenständige, künftige Entscheidungen:

- Ein realer `sdk-kotlin-v0.1.0`-Tag-Push — eine irreversible, extern
  sichtbare Aktion (öffentlicher GitHub-Packages-Push), die nur nach
  expliziter, gesonderter Rückfrage beim Auftraggeber läuft
  (Welle-Datei §3 vierter Punkt, `AGENTS.md` §3.10).
- Eine vierte SDK-Sprache oder ein vierter Vertriebsweg — bleibt laut
  `ADR-0109` Re-Evaluierungs-Trigger 1 offen für eine künftige, separate
  ADR; kein Slice dieser Welle nimmt das vorweg.
- Ein Folge-Package für SSE/NATS-Vollinhalt in Kotlin — bleibt laut
  `ADR-0109` §Entscheidung Festlegung 1/§Re-Evaluierungs-Trigger 2 offen
  für ein künftiges Release, nicht Teil dieser Welle.
- Eine künftige Architect-Sichtung von
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (3×, siehe
  Steering-Loop-Einträge oben) — kein Slice dieser Welle, eine
  eigenständige Regelschärfungs-Entscheidung.

## Verifikation

- `docs/reviews/review-slice-sdk-kotlin-projektgeruest.md` (0 HIGH/MEDIUM,
  3 LOW, keine Fixrunde) +
  `docs/reviews/verifikation-slice-sdk-kotlin-projektgeruest.md`.
- `docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` (1 HIGH,
  real behoben `9181ae21`) +
  `docs/reviews/verifikation-slice-sdk-kotlin-http-client-flaeche.md`
  (DoD erfüllt).
- `docs/reviews/review-slice-sdk-kotlin-grpc-client-flaeche.md` (0
  HIGH/MEDIUM/LOW/INFO, keine Fixrunde) +
  `docs/reviews/verifikation-slice-sdk-kotlin-grpc-client-flaeche.md`.
- `docs/reviews/review-slice-sdk-kotlin-pack-werkzeug.md` (0 HIGH, 1
  MEDIUM, 1 LOW, keine Fixrunde — beide Findings dem Beobachtungs-Register
  zugeführt) +
  `docs/reviews/verifikation-slice-sdk-kotlin-pack-werkzeug.md`.
- `docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` (1 HIGH,
  Fixrunde `21872c3d` real behoben) +
  `docs/reviews/review-slice-sdk-kotlin-publish-workflow-fixrunde.md` (0
  HIGH/MEDIUM, 1 INFO) +
  `docs/reviews/verifikation-slice-sdk-kotlin-publish-workflow.md` (DoD
  konform).
- `make gates`: grün auf dem Endstand (874 Dateien, 0 Befunde; Coverage
  82,70 % ≥ 80 %; `a-check`/`generated-sync`/`commit-traceability`/
  `baseline-verify` je ohne Befund), ungepiped geprüft (`AGENTS.md` §3.9)
  nach jedem Commit dieser Closure sowie erneut vor dieser Welle-Closure.
- Reales `pgchangefeed-kotlin-0.1.0.jar`-Artefakt — mehrfach unabhängig
  real erzeugt (Implementer, Reviewer, Fixrunden-Reviewer, Verifier, diese
  Planner-Closure — fünfter unabhängiger Bau, eigener `make
  sdk-pack-kotlin`-Lauf, Exit 0).
- `git tag -l "sdk-kotlin-v*"`: leer — bestätigt, dass kein realer
  GitHub-Packages-Publish-Versuch stattfand (Welle-Datei §3 vierter Punkt
  eingehalten, Post-Push-Risiko bleibt strukturell offen nach `AGENTS.md`
  §3.10).
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente (kein
`archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per `grep`,
dieselbe Feststellung wie bei den beiden vorigen SDK-Wellen) — die
Bedingung für Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine
Handarbeit als Ersatz.
