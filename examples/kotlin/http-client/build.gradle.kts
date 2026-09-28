// examples/kotlin/http-client/build.gradle.kts — Programm-Manifest des
// Kotlin-HTTP-Clients (ADR-0090 Festlegung 1/4, LP2 des Slice-Plans). Der
// Client selbst benutzt `java.net.http` aus der JDK-Standardbibliothek für
// den Transport; für die neun Fähigkeiten mit typisiertem JSON-Vertrag
// (`ConsumerClient`, `TablesAdminClient`, `ChangesClient`,
// `RetentionClient`) braucht er zusätzlich einen JSON-Decoder — weder
// `jnats` noch die JDK-Standardbibliothek tragen einen öffentlichen. Die
// Tabellen-Auflistung (`TablesClient`) liest die rohe Server-Antwort und
// bleibt unverändert ohne JSON-Decoder. `com.google.code.gson:gson` ist
// dieselbe, bereits gepinnte, öffentliche Client-Bibliothek wie in
// `examples/kotlin/nats-stream-client/build.gradle.kts` (Version/Lizenz/
// transitive Abhängigkeiten dort bereits bewertet, kein neuer Pin).
// Ausschließlich der Testlauf zieht zusätzlich JUnit 5 über Kotlins
// `kotlin-test`-Fassade. Gepinnte, exakte Versionen — dasselbe Prinzip wie
// die Digest-Pinnung der Basis-Images im Dockerfile.
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
    implementation("com.google.code.gson:gson:2.14.0")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

tasks.test {
    useJUnitPlatform()
}
