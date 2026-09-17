# Review-Report: slice-103 — 2026-09-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten): geprüft wird der Diff gegen den Slice-Plan, `ADR-0090`
und Modul 5/8 des Baseline-Regelwerks, nicht gegen die DoD (das ist
Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 5b10620..HEAD` — Commits `911050c`, `6fe60f7`,
`111f1cb` (Slice `slice-103`, Plan
`docs/plan/planning/in-progress/slice-103-grpc-client-kotlin.md`), letzter
der sechs C#/Kotlin-Matrix-Slices.

**Skill:** `.harness/skills/reviewer.md` @ `0a18679` · **Modell:**
claude-sonnet-5 · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-103-grpc-client-kotlin.md` (§1
  Abgrenzung, §2 LP1–LP3, §4 Trigger, §6 Risiken, §8 Sub-Area-Prüfungen)
- `docs/plan/planning/done/slice-102-grpc-client-csharp.md` und
  `docs/reviews/review-slice-102.md`/`architect-verdict-slice-101-jnats-
  bouncycastle.md` als Form- und Präzedenz-Referenz
- `ADR-0090` (§Kontext, §Entscheidung Festlegung 1–4, §Fitness Function,
  §Re-Evaluierungs-Trigger)
- `ADR-0060` (Server-Stream-Semantik)
- Baseline-Regelwerk `modul-05-planning-harness.md` §Trigger je
  Lifecycle-Übergang; `modul-08-agentenrollen.md` §Rollen-Regeln,
  §Konflikt-Pfad als Rollen-Sequenz
- `AGENTS.md` §3.1/§3.5/§3.9/§3.12/§3.13
- Eigene Nachprüfung mit Netzzugriff: zwei reale `docker build`-Läufe gegen
  `examples/kotlin/Dockerfile` (ohne und mit `--build-context
  proto=proto`), eine reale Mutation von
  `examples/kotlin/grpc-client/build.gradle.kts` (Entfernen des
  `protobuf-java:4.36.1`-Pins, Gegenbau), Abruf von
  `grpc-kotlin-stub-1.5.0.pom` (Maven Central), `grep` gegen
  `spec/pflichtenheft.md` und `docs/plan/planning/observations/`

---

## Findings

### F-1 — Zusatzkontext-Kopplung bleibt strukturell für alle vier Sprach-Programme zwingend (Fortsetzung von `review-slice-102.md` F-1)

- `kategorie`: LOW
- `quelle`: Maintainability (Plan-Genauigkeit gegen `ADR-0090` §Fitness
  Function)
- `pfad`: `examples/kotlin/Dockerfile:67` (COPY-Zeile in der gemeinsamen
  `build`-Stufe), `harness/mk/examples.mk:8-44`
- `befund`: Eigener Nachbau bestätigt: `docker build examples/kotlin` (ohne
  `--build-context proto=proto`) bricht bereits beim Auflösen der
  `docker.io/library/proto:latest`-Referenz ab — unabhängig vom
  angeforderten `--target`, weil `COPY --from=proto …` unbedingt in der von
  allen vier Modulen geteilten `build`-Stufe steht. Dieselbe, bereits in
  `slice-102` als LOW (F-1) akzeptierte Nebenwirkung derselben
  Dockerfile-Konstruktion (eine gemeinsame `build`-Stufe für alle Programme
  einer Sprache), hier korrekt übertragen und in `harness/mk/examples.mk`
  transparent kommentiert („unabhängig vom angeforderten `--target`"). Kein
  neuer Defekt, aber die zweite Instanz desselben Musters.
- `verifizierbar`: ja — `docker build examples/kotlin` ohne
  `--build-context` reproduziert den Abbruch
- `klasse`: Zusatzkontext-Kopplung breiter als DoD-Wortlaut, transparent
  dokumentiert (2. Instanz — noch unter der 3×-MEDIUM-Schwelle der
  Reviewer-Skill-Klassifikation)

## Negativbefunde

- geprüft, ohne Befund: **Die Trigger-Frage selbst (§4 Rückführung 1)** —
  siehe Verdikt unten. Kein Fund; die Klassifikation des Implementers
  („normale Anpassung", kein Trigger) ist nach eigener Prüfung korrekt.
- geprüft, ohne Befund: `examples/kotlin/build.gradle.kts`,
  `examples/kotlin/grpc-client/build.gradle.kts`,
  `examples/kotlin/settings.gradle.kts` — alle vom Implementer-Bericht
  genannten Zusatz-Artefakte sind real vorhanden und mit spezifischer,
  überprüfbarer technischer Begründung kommentiert (nicht pauschal):
  `io.grpc:grpc-bom:1.84.0` (Version-Alignment für `io.grpc:*`-Transitive),
  `io.grpc:protoc-gen-grpc-java:1.84.0` (Java-Service-Deskriptor, den
  `grpc-kotlin` als Vorstufe braucht), `com.google.protobuf:protoc:4.36.1`/
  `protobuf-java:4.36.1` (explizit gegen die transitiv gezogene
  `3.25.9`-Zeile angehoben), `kotlinx-coroutines-core:1.11.0`
  (`implementation`, weil `grpc-kotlin-stub` es nur `runtime` deklariert)
- geprüft, mit eigenem Mutationsbeleg bestätigt: **Die
  `protobuf-java`-Versionsanhebung ist real notwendig, keine
  Vorsichtsmaßnahme.** Eigener Gegenbau (Kopie des Verzeichnisses,
  `implementation("com.google.protobuf:protobuf-java:4.36.1")` entfernt,
  `docker build --build-context proto=… --target runtime-grpc`) reproduziert
  exakt die im Kommentar benannten Fehler: 27 `cannot find symbol`-Fehler in
  `Changestream.java`, darunter `GeneratedMessage.isStringEmpty`,
  `FileDescriptor.getMessageType`, `resolveAllFeaturesImmutable` und
  `package com.google.protobuf.RuntimeVersion does not exist` — dieselben
  Symbole, die der Kommentar in `grpc-client/build.gradle.kts` nennt. Kein
  neues transitives Risiko identifiziert: `protoc` (4.36.1) und
  `protobuf-java` (4.36.1) tragen dieselbe Major-Zeile, der reale Bau
  (positiv, mit Pin) kompiliert und testet fehlerfrei; ein Editions-Sprung
  3.x→4.x ist damit intern konsistent gehalten, nicht offen gelassen.
- geprüft, mit eigenem Beleg bestätigt: `kotlinx-coroutines-core`-Scope —
  `grpc-kotlin-stub-1.5.0.pom` (Maven Central, eigener Abruf) deklariert
  `kotlinx-coroutines-core-jvm:1.10.2` wörtlich mit
  `<scope>runtime</scope>`, keinem `compile`/`api`. Die Kommentar-Aussage
  ist exakt zutreffend, keine Übertreibung.
- geprüft, ohne Befund: **Zusatzkontext-Mechanismus identisch zu
  `slice-102`.** `examples/kotlin/Dockerfile` und `harness/mk/
  examples.mk` sind strukturell deckungsgleich mit `examples/csharp/
  Dockerfile`/derselben `.mk`-Sektion aus `slice-102`: derselbe
  Kontextname `proto`, dieselbe `COPY --from=proto …`-Form, derselbe
  Aufrufer-seitige `--build-context proto=proto`. Der mildernde Kommentar
  über der `COPY`-Zeile in `examples/csharp/Dockerfile` ist tatsächlich —
  nicht nur behauptet — mitgeführt (`examples/kotlin/Dockerfile:59-66`,
  wortnah identische Erklärung „kein Fallback, kein übersprungener
  Schritt").
- geprüft, mit eigenem Mutationsbeleg bestätigt: **Negativ-Beleg.** Eigener
  `docker build -t review-negative-test examples/kotlin` (ohne
  `--build-context`) scheitert real mit „pull access denied, repository
  does not exist or may require authorization … `docker.io/library/
  proto:latest`" — identisch zur in `slice-102` gemessenen, irreführenden
  Meldung.
- geprüft, mit eigenem Mutationsbeleg bestätigt: **Positiv-Beleg.** Eigener
  `docker build --build-context proto=proto --target runtime-grpc
  examples/kotlin` läuft bis zum Image-Export durch (vollständiger
  Cache-Hit über alle 28 Layer inkl. `RUN ./gradlew … test`/
  `installDist`) — die reale Testausführung, die dieser Cache-Treffer
  voraussetzt, hat bereits real stattgefunden und nicht nur behauptet.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` §4 „Zugriff über
  den gRPC-Change-Stream" — dritte und letzte Zeile (Kotlin) im
  `**Beispiele:**`-Block korrekt ergänzt. Die Vollständigkeits-Aussage der
  vollen Matrix steht **explizit** in der Änderungshistorie-Zeile 1.25
  („Mit dieser Zeile ist die volle Matrix … vollständig") — im
  Abschnitts-Fließtext selbst steht sie nicht, was konsistent mit
  `ADR-0090` §7 ist („die Prosa des Abschnitts bleibt unverändert und
  einfach"; keine Matrix-Übersicht als eigenes Artefakt, §1 Ausschlussliste
  dieses Slice-Plans). Kein Fund — der Plan-Wortlaut in §1 selbst steht in
  der Closure-Notiz, nicht zwingend im Handbuchtext.
- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-023` — bereits vor
  diesem Diff (`git log -S"volle Matrix"` → Commit `e892a67`, Teil der
  `slice-099`-Closure, außerhalb des hier geprüften Bereichs `5b10620..
  HEAD`) auf die volle Matrix gezogen; die Implementer-Aussage „bereits von
  `slice-099` aktualisiert" ist zutreffend
- geprüft, ohne Befund: Out-of-Scope-Disziplin — `git diff --stat` gegen
  `.a-check.yml`, `Makefile`, `docs/plan/adr/`, `spec/` in
  `5b10620..HEAD` ist leer; kein Go-gRPC-Client, keine `.proto`-
  Inhaltsänderung, kein committeter Stub (`git ls-files examples/kotlin`
  ohne generierte `.java`/`.kt`-Dateien, `build/` bleibt gitignored)
- geprüft, ohne Befund: Commit-Traceability — `make commit-traceability
  RANGE=5b10620..HEAD` real ausgeführt, 3 Commits, 0 Befunde; kein
  `SPEC-*`/`ARC-*` im Betreff
- geprüft, ohne Befund: `examples/kotlin/grpc-client/src/{main,test}/**` —
  Cli/Config/Format/Main folgen demselben Fehlerpfad-Muster wie
  `examples/csharp/grpc-client` (Exit 2 bei fehlender Adresse/Token, Exit 1
  bei Stream-Ende/`StatusException`); Tests (`CliTest`, `FormatTest`) sind
  echte, netzlose Assertions ohne Bindungs-Lücke
- geprüft, ohne Befund: Beobachtungs-Register-Zitate in §8 des Slice-Plans
  (`handbuch-nicht-nachgezogen-…`, `handbuch-versionshistorie-…`,
  `nicht-blockierender-workflow-alarmmuedigkeit`,
  `github-actions-unverifizierbar-lokal`,
  `arbeit-ueberholt-stehenden-traeger`) — alle fünf Verzeichnisse
  existieren mit nicht leerem `evidence/`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Zusatzkontext-Kopplung breiter als
DoD-Wortlaut, transparent dokumentiert (2. Instanz)

## Verdikt zur Trigger-Frage (§4 Rückführung 1)

**Kein Trigger-Ereignis — Granularitäts-Unterschied, kein
Mechanismus-Unterschied.** Der Plan-Trigger benennt als Beispiel
ausdrücklich „einen anderen **Kontext-Zugriff**" (`ADR-0090` Festlegung 2:
der eine, benannte Zusatzkontext, `docker buildx build --build-context …`,
`COPY --from=proto …`). Dieser Mechanismus ist real gegenbaut identisch zu
`slice-102`: gleicher Kontextname, gleiche COPY-Form, gleicher
Abbruch-Fehler ohne Kontext, gleicher Erfolgspfad mit Kontext. Was
gewachsen ist, ist die Zahl der Bibliotheks-Koordinaten, die die
Gradle-/Maven-Werkzeugkette braucht, um denselben Generator-Schritt
auszuführen, den `Grpc.Tools` in einem einzigen NuGet-Paket bündelt — eine
Paketierungs-Konvention der jeweiligen Ökosysteme (artefakt-je-Konzern bei
Gradle/Maven vs. monolithisch bei NuGet), keine andere Zugriffsform auf die
`.proto`. `ADR-0090` selbst delegiert das ausdrücklich: „Der Pin gehört dem
umsetzenden Zug" (Festlegung 4) — die dort genannten drei Bibliotheken sind
explizit als „Kandidaten", nicht als abschließende Liste markiert.

Abgrenzung zum Präzedenzfall (`slice-101`,
`architect-verdict-slice-101-jnats-bouncycastle.md`): Dort hatte der
Plan-Trigger eine **deterministische** Form („X tritt ein → Annahme ist
falsch", ohne Ermessens-Spielraum), und die Implementer-Bewertung ersetzte
eine vorab verlangte *Rolle* (Architect/Planner entscheidet über die
Konsequenz), nicht nur eine Tatsachenfeststellung. Hier verlangt der
Trigger selbst eine **zusammengesetzte** Bedingung („trägt nicht
unverändert" **und** „eigenständige, abweichende Lösung nötig") — die
Einordnung, ob eine Beobachtung darunterfällt, ist ein normaler Teil der
Umsetzungsarbeit, keine Ersetzung einer vorbehaltenen Entscheidung. Beide
Kriterien des Triggers sind nach eigener, unabhängiger Prüfung (zwei reale
Baus, eine reale Mutation) nicht erfüllt: Der Zugriffsmechanismus trägt
unverändert; es gibt keine „eigenständige, abweichende Lösung" — die
zusätzlichen Artefakte sind Zubehör desselben, unveränderten Wegs.

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. F-1 ist eine bereits in
`slice-102` akzeptierte, transparent dokumentierte Nebenwirkung derselben
Konstruktion.

Da keine Fixrunde am Implementer nötig ist, wird die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
`docs/plan/planning/in-progress/slice-103-grpc-client-kotlin.md` in diesem
Commit selbst auf `[x]` nachgezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** F-1 geht als Klasse in die Slice-Closure §7. Geprüft
(`grep -rli "zusatzkontext" docs/plan/planning/observations/` → kein
Fund): Es existiert **noch kein** Beobachtungs-Register-Verzeichnis für
diese Klasse — `review-slice-102.md` F-1 wurde als akzeptierte LOW-Notiz in
der Closure-Prosa behandelt, aber nicht als Registereintrag angelegt. Mit
diesem Fund liegt die Klasse **2×** vor (`slice-102`, `slice-103`), beide
Male dieselbe strukturelle Ursache (eine gemeinsame `build`-Stufe für alle
Programme einer Sprache). Empfehlung an die Planner-Closure: neues
Verzeichnis `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/`
mit `evidence/slice-102.md` und `evidence/slice-103.md` anlegen, Stand
`offen` (2×, unter der 3×-Schwelle) — keine Entscheidung des Reviewers,
nur ein benannter Kandidat. Dieser Report selbst ist ein Lauf-Beleg und
wird über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität (inkl. `make gates` grün) prüft der
Verifier separat (Modul 11).
