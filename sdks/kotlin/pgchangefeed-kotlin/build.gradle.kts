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
// Der `publishing`-Block hier trägt nur die Maven-Koordinate; der reale
// GitHub-Packages-`repositories{}`-Eintrag (`GITHUB_TOKEN`, Registry-URL)
// ist Sache des Publish-Workflow-Zuges (`slice-sdk-kotlin-publish-workflow`,
// ADR-0109 §Konsequenzen Folgepflicht 1) — kein Publish-Aufruf in diesem
// Slice.
plugins {
    kotlin("jvm") version "2.4.20"
    `maven-publish`
}

group = "io.github.pt9912"
version = "0.1.0"

kotlin {
    jvmToolchain(21)
}

dependencies {
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

tasks.test {
    useJUnitPlatform()
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            groupId = "io.github.pt9912"
            artifactId = "pgchangefeed-kotlin"
            version = "0.1.0"
            from(components["java"])
        }
    }
}
