// sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts — Maven-Metadaten und
// Bau-Konfiguration des Kotlin-SDK-Pakets (ADR-0109 Festlegung 1/2/3/4/5).
// Eigenständiges Projekt, kein Untermodul von examples/kotlin/ (Festlegung 3)
// — kein Import auf einen privaten Baum dieses Repos.
//
// Kotlin-Gradle-Plugin-Version real gegen Maven Central nachgemessen
// (`repo1.maven.org/maven2/org/jetbrains/kotlin/kotlin-gradle-plugin/maven-metadata.xml`,
// heute, 2026-09-20): aktuellste Version weiterhin 2.4.20 (`lastUpdated`
// 2026-09-07, keine Drift ggü. der in `examples/kotlin/build.gradle.kts` am
// 2026-09-17 gemessenen Version — real re-verifiziert, nicht blind
// übernommen, AGENTS.md §3.12).
//
// `maven-publish` ist das eingebaute Gradle-Kern-Plugin (ADR-0109
// Festlegung 5, F1) — kein Dritt-Plugin wie `com.vanniktech.maven-publish`.
// Der `publishing`-Block trägt die Maven-Koordinate UND den
// GitHub-Packages-`repositories{}`-Eintrag (real recherchiertes
// Minimalrezept, ADR-0109 §Kontext Recherche:
// `docs.github.com/…/publishing-java-packages-with-gradle`) — Registry-URL
// `https://maven.pkg.github.com/pt9912/pg-change-feed`, Zugangsdaten aus den
// Umgebungsvariablen `GITHUB_ACTOR`/`GITHUB_TOKEN`. `GITHUB_ACTOR` ist ein
// von GitHub Actions automatisch bereitgestellter Default-Umgebungswert;
// `GITHUB_TOKEN` wird vom aufrufenden Workflow
// (`.github/workflows/sdk-kotlin-release.yml`, ADR-0109 §Konsequenzen
// Folgepflicht 1) explizit als Umgebungsvariable gesetzt (`secrets.GITHUB_TOKEN`)
// — kein neues Repository-Secret, kein `secrets.<NAME>`-Verweis hier.
// Außerhalb dieses Workflows (lokal, in `make sdk-pack-kotlin`) bleiben
// beide Umgebungsvariablen leer — `make sdk-pack-kotlin` baut ausschließlich
// die `pack-export`-Docker-Stufe (`test`/`build`, kein `publish`-Task).
// `sdks/kotlin/Dockerfile` trägt daneben eine eigene `publish`-Stufe (baut
// auf `build` auf), die ausschließlich `.github/workflows/sdk-kotlin-release.yml`
// baut und per `docker run -e GITHUB_ACTOR=... -e GITHUB_TOKEN=...` mit
// echten Zugangsdaten zur Laufzeit startet (ADR-0109 Festlegung 5: Docker-only
// bis einschließlich `publish` — Korrektur von Review-Finding F-1,
// docs/reviews/review-slice-sdk-kotlin-publish-workflow.md).
//
// `com.google.code.gson:gson` ist die JSON-Bibliothek der HTTP-Client-Fläche
// (slice-sdk-kotlin-http-client-flaeche, `SPEC-018`/`SPEC-022`) — dieselbe
// bereits im selben Repo real bewertete, gepinnte Version wie
// `examples/kotlin/nats-stream-client/build.gradle.kts` (dort Lizenz/
// transitive Abhängigkeiten bereits geprüft: nur
// `com.google.errorprone:error_prone_annotations` im `compile`-Scope,
// Apache-2.0). Weder `java.net.http` noch die JDK-Standardbibliothek tragen
// einen öffentlichen JSON-Decoder. Gemessen am 2026-09-20 (Maven Central
// maven-metadata.xml, repo1.maven.org/maven2/com/google/code/gson/gson/
// maven-metadata.xml): weiterhin `2.14.0` (`<latest>`/`<release>`, keine
// Drift ggü. der Messung vom 2026-09-18 in `nats-stream-client` —
// `AGENTS.md` §3.12).
//
// gRPC-Stream-Client-Fläche (`slice-sdk-kotlin-grpc-client-flaeche`,
// `ADR-0109` Festlegung 1/3, `SPEC-020`): dieselben Koordinaten wie
// `examples/kotlin/grpc-client/build.gradle.kts` (`slice-103`) — real am
// heutigen Bau-Zeitpunkt dieses Slice (2026-09-20, Maven Central
// maven-metadata.xml je Artefakt, Gradle Plugin Portal für das
// `com.google.protobuf`-Plugin) neu gemessen, nicht aus `examples/kotlin/`
// oder `ADR-0109`s Kontext-Messung (2026-09-17) unbesehen übernommen
// (`AGENTS.md` §3.12):
//   io.grpc:grpc-kotlin-stub        -> 1.5.0 (unverändert; `<latest>`/
//     `<release>` der Maven-Metadaten zeigen einen Commit-Hash-Eintrag
//     — ein CI-Snapshot-Artefakt, kein echtes Release, lastUpdated
//     2025-09-16, also bereits vor `examples/kotlin/grpc-client`s eigener
//     2026-09-17-Messung vorhanden und dort korrekt ignoriert; 1.5.0 bleibt
//     die tatsächlich zuletzt veröffentlichte, reguläre Version)
//   io.grpc:protoc-gen-grpc-kotlin  -> 1.5.0 (dieselbe Anomalie, dieselbe
//     Auflösung wie oben)
//   io.grpc:grpc-netty-shaded       -> 1.84.0 (unverändert)
//   io.grpc:grpc-bom                -> 1.84.0 (unverändert)
//   io.grpc:protoc-gen-grpc-java    -> 1.84.0 (unverändert)
//   org.jetbrains.kotlinx:kotlinx-coroutines-core -> 1.11.0 (unverändert)
//   com.google.protobuf (Gradle-Plugin, Gradle Plugin Portal) -> 0.10.0
//     (unverändert)
//   com.google.protobuf:protoc -> 4.36.2 (REALE DRIFT ggü. der
//     2026-09-17-Messung in `examples/kotlin/grpc-client`: dort 4.36.1.
//     Maven-Metadaten tragen zusätzlich einen `21.0-rc-1`-Eintrag —
//     lastUpdated 2026-09-17, ein Release-Candidate einer neuen
//     Versionszählung, kein stabiles Release; 4.36.2 bleibt die aktuellste
//     stabile Version, real bezogen unter
//     repo1.maven.org/maven2/com/google/protobuf/protoc/4.36.2/)
//   com.google.protobuf:protobuf-java -> 4.36.2 (dieselbe reale Drift wie
//     `protoc` oben — beide müssen dieselbe Major-Zeile tragen, siehe
//     Kommentar unten zur `io.grpc:grpc-protobuf`-Transitiv-Version; real
//     erneut geprüft: `io.grpc:grpc-protobuf:1.84.0` zieht weiterhin
//     transitiv `protobuf-java:3.25.9`, unverändert ggü. der
//     2026-09-17-Messung — die explizite Anhebung bleibt aus demselben
//     Grund nötig)
//
// NATS-Vollinhalts-Stream-Client-Fläche (`slice-sdk-kotlin-nats-stream-client-flaeche`,
// `ADR-0109` Festlegung 1, `SPEC-024`): `io.nats:jnats` real am heutigen
// Bau-Zeitpunkt dieses Slice (2026-09-22, Maven Central maven-metadata.xml,
// repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml) neu gemessen,
// nicht aus `examples/kotlin/nats-stream-client/build.gradle.kts`s
// 2026-09-18-Messung unbesehen übernommen (`AGENTS.md` §3.12): weiterhin
// `2.26.3` (`<latest>`/`<release>`, keine Drift, dieselbe Version wie bei
// jenem Beispiel). Kein zweiter JSON-Decoder nötig — dieselbe bereits
// gepinnte `com.google.code.gson:gson:2.14.0` deserialisiert auch die
// NATS-Nachrichten (`SPEC-024` teilt das Schema mit `SPEC-021`).
//
// Diese Version-Hebung (`0.1.0` -> `0.2.0`) trägt außerdem den
// SSE-Client-Fläche (`slice-sdk-kotlin-sse-client-flaeche`, `SPEC-021`) —
// beide Flächen bündeln ihren Version-Bump gemeinsam in diesem letzten
// Flächen-Slice der Welle (`welle-sdk-kotlin-vollabdeckung` §1). Additiv,
// rückwärtskompatible Erweiterung — SemVer-Minor, keine ADR-pflichtige
// Ausnahme (Slice-Plan §6).
plugins {
    kotlin("jvm") version "2.4.20"
    `maven-publish`
    id("com.google.protobuf") version "0.10.0"
}

group = "io.github.pt9912"
version = "0.2.0"

kotlin {
    jvmToolchain(21)
}

dependencies {
    implementation("com.google.code.gson:gson:2.14.0")
    implementation(platform("io.grpc:grpc-bom:1.84.0"))
    implementation("io.grpc:grpc-kotlin-stub:1.5.0")
    implementation("io.grpc:grpc-protobuf")
    implementation("io.grpc:grpc-stub")
    implementation("com.google.protobuf:protobuf-java:4.36.2")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.11.0")
    implementation("io.nats:jnats:2.26.3")
    runtimeOnly("io.grpc:grpc-netty-shaded")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

// Die `.proto`-Quelle liegt NICHT im committeten Baum (`ADR-0109`
// Festlegung 3, kein committeter Stub) — sie kommt erst im Docker-Bau nach
// `src/main/proto/` (`sdks/kotlin/Dockerfile`), kopiert aus dem
// zusätzlichen, benannten Bau-Kontext `proto`
// (`docker build --build-context proto=proto …`, Muster
// `examples/kotlin/Dockerfile`/`harness/mk/examples.mk`). Ohne diesen
// Kontext bricht der Bau an der `COPY`-Zeile im Dockerfile ab — kein
// stiller Fallback.
protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:4.36.2"
    }
    plugins {
        create("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:1.84.0"
        }
        create("grpckt") {
            artifact = "io.grpc:protoc-gen-grpc-kotlin:1.5.0:jdk8@jar"
        }
    }
    generateProtoTasks {
        all().forEach { task ->
            task.plugins {
                create("grpc")
                create("grpckt")
            }
        }
    }
}

tasks.test {
    useJUnitPlatform()
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            groupId = "io.github.pt9912"
            artifactId = "pgchangefeed-kotlin"
            version = "0.2.0"
            from(components["java"])
        }
    }
    repositories {
        maven {
            name = "GitHubPackages"
            url = uri("https://maven.pkg.github.com/pt9912/pg-change-feed")
            credentials {
                username = System.getenv("GITHUB_ACTOR")
                password = System.getenv("GITHUB_TOKEN")
            }
        }
    }
}
