// Gast-Projekt der Kompatibilitätsmessung (make test-sdk-kompat): bindet die
// veröffentlichte Bibliothek in der Version -PpgcfVersion und übersetzt den
// Quelltext src/<-PgastQuelle>/kotlin — alt: nur Formen der 0.5.x-Fehlertypen,
// neu: Formen mit Meldungscode. `kopiereLaufzeit` legt die aufgelöste
// Laufzeit-Klassenpfad-Menge (Bibliothek und Abhängigkeiten) als Jars ab.
plugins {
    kotlin("jvm") version "2.4.20"
}

val pgcfVersion = providers.gradleProperty("pgcfVersion").getOrElse("0.5.0")
val gastQuelle = providers.gradleProperty("gastQuelle").getOrElse("alt")

kotlin {
    jvmToolchain(21)
}

sourceSets {
    main {
        kotlin.srcDir("src/$gastQuelle/kotlin")
    }
}

dependencies {
    implementation("io.github.pt9912:pgchangefeed-kotlin:$pgcfVersion")
    // Die Fehlertypen der Bibliothek tragen io.grpc.Status in ihrer Signatur, die
    // Bibliothek veröffentlicht grpc-api aber nur für die Laufzeit: der
    // Übersetzer des Aufrufers braucht es ausdrücklich (gleicher Stand wie die
    // grpc-bom der Bibliothek).
    implementation("io.grpc:grpc-api:1.84.0")
}

tasks.register<Copy>("kopiereLaufzeit") {
    from(configurations.runtimeClasspath)
    into(layout.buildDirectory.dir("laufzeit"))
}
