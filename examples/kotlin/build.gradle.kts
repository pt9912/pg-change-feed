// examples/kotlin/build.gradle.kts — gemeinsame Plugin-Version-Verwaltung für
// alle Module der Kotlin-Sprach-Wurzel (ADR-0087 Festlegung 3, ADR-0090
// Festlegung 1/5). Die Kotlin-Gradle-Plugin-Version wird hier **einmal**
// gepinnt (`apply false`) und von jedem Modul ohne eigene Versionsnummer
// referenziert — Pin-Hebung ist ein bewusster Commit, keine stille Auflösung
// (`SPEC-023`).
//
// Gemessen am 2026-09-17 (Maven Central,
// repo1.maven.org/maven2/org/jetbrains/kotlin/kotlin-gradle-plugin/maven-metadata.xml):
// aktuellste Version 2.4.20.
//
// `com.google.protobuf` (Gradle Plugin Portal) trägt ausschließlich
// `grpc-client` (ADR-0090 Festlegung 2/4, slice-103) — hier trotzdem zentral
// gepinnt, dasselbe Muster wie die Kotlin-Plugin-Version. Gemessen am
// 2026-09-17 (plugins.gradle.org/m2/com/google/protobuf/com.google.protobuf.gradle.plugin/maven-metadata.xml):
// aktuellste Version 0.10.0.
plugins {
    kotlin("jvm") version "2.4.20" apply false
    id("com.google.protobuf") version "0.10.0" apply false
}
