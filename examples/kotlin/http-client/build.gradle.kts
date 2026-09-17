// examples/kotlin/http-client/build.gradle.kts — Programm-Manifest des
// Kotlin-HTTP-Clients (ADR-0090 Festlegung 1/4, LP2 des Slice-Plans). Der
// Client selbst braucht kein Fremdmodul — er liest JSON/HTTP nicht (die
// Antwort ist Klartext) und benutzt ausschließlich `java.net.http` aus der
// JDK-Standardbibliothek; ausschließlich der Testlauf zieht Fremdmodule
// (JUnit 5 über Kotlins `kotlin-test`-Fassade). Gepinnte, exakte Versionen —
// dasselbe Prinzip wie die Digest-Pinnung der Basis-Images im Dockerfile.
//
// Gemessen am 2026-09-17 (Maven Central maven-metadata.xml):
//   org.jetbrains.kotlin:kotlin-test-junit5   -> 2.4.20 (deckungsgleich mit
//                                                der Kotlin-Gradle-Plugin-Version)
//   org.junit.platform:junit-platform-launcher -> 6.1.3
plugins {
    kotlin("jvm")
    application
}

kotlin {
    jvmToolchain(21)
}

application {
    mainClass.set("cdcexamples.http.MainKt")
}

dependencies {
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

tasks.test {
    useJUnitPlatform()
}
