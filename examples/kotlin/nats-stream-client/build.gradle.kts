// examples/kotlin/nats-stream-client/build.gradle.kts — Programm-Manifest
// des Kotlin-Clients für den dritten, vollinhaltstragenden NATS-Zustellweg
// (ADR-0100, ADR-0090 Festlegung 1/4). Denselben gepinnten
// Verbindungsweg wie `examples/kotlin/nats-client`: `io.nats:jnats:2.26.3`
// (Version/Lizenz/transitive Abhängigkeiten dort bereits bewertet, kein
// neuer Pin).
//
// Anders als das Wecksignal trägt dieser Zustellweg vollständigen
// JSON-Inhalt (`SPEC-024`), den dieses Programm in benannte Felder
// dekodiert (`change_id`, `table`, `operation`, `new_image` — dieselbe
// Ausgabeform wie `examples/kotlin/grpc-client`). Weder `jnats` noch die
// JDK-Standardbibliothek tragen einen öffentlichen JSON-Decoder; keines der
// bestehenden Kotlin-Beispiele braucht bislang einen (sie reichen HTTP-/
// SSE-Antworten unverändert als Text durch). `com.google.code.gson:gson`
// ist die gepinnte, öffentliche Client-Bibliothek dieses Programms
// (ADR-0090 Festlegung 4) für genau diesen Zweck.
//
// Gemessen am 2026-09-18 (Maven Central maven-metadata.xml,
// repo1.maven.org/maven2/com/google/code/gson/gson/maven-metadata.xml):
// neueste stabile Version 2.14.0 (`<latest>`/`<release>`, das einzige
// Prerelease in der Liste ist `2.13.2-rc1`, kein Kandidat). Transitive
// Abhängigkeit, real gemessen (gson-2.14.0.pom): ausschließlich
// `com.google.errorprone:error_prone_annotations:2.48.0` im
// `compile`-Scope (reine Annotationen, keine Laufzeit-Logik — dieselbe
// Klasse wie `jspecify` beim jnats-Pin oben); alle übrigen Abhängigkeiten
// (`junit`, `com.google.truth`, `guava`/`guava-testlib`) sind
// `test`-Scope und werden nicht mitgezogen. Apache-2.0-Lizenz, keine
// zweite, eigenständig zu bewertende Fremdbibliothek wie bei `jnats`
// (Bouncy Castle).
plugins {
    kotlin("jvm")
    application
}

kotlin {
    jvmToolchain(21)
}

application {
    mainClass.set("cdcexamples.natsstream.MainKt")
}

dependencies {
    implementation("io.nats:jnats:2.26.3")
    implementation("com.google.code.gson:gson:2.14.0")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

tasks.test {
    useJUnitPlatform()
}
