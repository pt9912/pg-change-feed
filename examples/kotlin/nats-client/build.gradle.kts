// examples/kotlin/nats-client/build.gradle.kts — Programm-Manifest des
// Kotlin-NATS-Clients (ADR-0090 Festlegung 1/4, LP2 des Slice-Plans,
// slice-101). Der Client selbst braucht **eine** gepinnte, öffentliche
// Client-Bibliothek (ADR-0090 Festlegung 4): `io.nats:jnats`, das
// tabellen-granulare Wecksignal-Subjekt abonnieren (`ADR-0055`/`ADR-0056`);
// die Änderung selbst holt er wie die anderen Kotlin-Beispiele über
// `java.net.http.HttpClient` aus der JDK-Standardbibliothek. Ausschließlich
// der Testlauf zieht zusätzlich JUnit 5 über Kotlins `kotlin-test`-Fassade,
// dasselbe Muster wie `examples/kotlin/http-client/build.gradle.kts`.
//
// Gemessen am 2026-09-17 (Maven Central maven-metadata.xml,
// repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml): neueste stabile
// Version 2.26.3 (`<latest>`/`<release>`, kein Snapshot).
//
// Transitive Abhängigkeit, real gemessen (jnats-2.26.3.pom /
// jnats-2.26.3.module, beide identisch): `org.bouncycastle:bcprov-lts8on:2.73.12.1`
// (compile-scope, Apache-artige Bouncy-Castle-Lizenz — Krypto-Bibliothek für
// NKey/Ed25519-Authentifizierung) und `org.jspecify:jspecify:1.0.0`
// (Apache-2.0, reine Nullability-Annotationen, keine Laufzeit-Logik). Anders
// als beim C#-Pendant (`NATS.Net`, siehe
// `examples/csharp/Directory.Packages.props` — dort ausschließlich
// gleichprojektige `NATS.*`- und `Microsoft.Extensions.*`-Pakete für
// net10.0) zieht `io.nats:jnats` damit **eine** echte Fremdbibliothek
// außerhalb des nats-io-Projekts nach. Bewertet: Bouncy Castle ist eine
// seit Jahrzehnten breit eingesetzte, quelloffene Kryptobibliothek
// (u. a. Android-, Spring-Ökosystem) mit permissiver Lizenz; sie wird
// unconditional (compile-scope) mitgezogen, auch wenn dieses Beispiel keine
// NKey-Credentials nutzt. Diese Bewertung ist der Lerneintrag aus dem
// Closure-Trigger dieses Slice (§5/§7 des Slice-Plans) — kein Grund für die
// Rückführung `in-progress → next` (§4), weil `io.nats:jnats` in
// `ADR-0090` Festlegung 4 bereits als Accepted benannt ist und Bouncy
// Castle keine undurchsichtige, unlizenzierte oder verwaiste Abhängigkeit
// ist.
plugins {
    kotlin("jvm")
    application
}

kotlin {
    jvmToolchain(21)
}

application {
    mainClass.set("cdcexamples.nats.MainKt")
}

dependencies {
    implementation("io.nats:jnats:2.26.3")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

tasks.test {
    useJUnitPlatform()
}
