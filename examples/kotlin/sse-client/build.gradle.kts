// examples/kotlin/sse-client/build.gradle.kts — Programm-Manifest des
// Kotlin-SSE-Clients (ADR-0090 Festlegung 1/4, LP2 des Slice-Plans). Der
// Client selbst braucht kein Fremdmodul — er liest den Stream über
// `java.net.http.HttpClient`/`HttpResponse.BodyHandlers.ofLines()` aus der
// JDK-Standardbibliothek (**keine neue Abhängigkeit**, ADR-0090
// Festlegung 1: kein `okhttp-sse`); ausschließlich der Testlauf zieht
// Fremdmodule (JUnit 5 über Kotlins `kotlin-test`-Fassade), dasselbe Muster
// wie `examples/kotlin/http-client/build.gradle.kts`.
//
// Gemessen am 2026-09-17 (Maven Central maven-metadata.xml, deckungsgleich
// mit dem http-client-Manifest):
//   org.jetbrains.kotlin:kotlin-test-junit5   -> 2.4.20
//   org.junit.platform:junit-platform-launcher -> 6.1.3
plugins {
    kotlin("jvm")
    application
}

kotlin {
    jvmToolchain(21)
}

application {
    mainClass.set("cdcexamples.sse.MainKt")
}

dependencies {
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

tasks.test {
    useJUnitPlatform()
}
