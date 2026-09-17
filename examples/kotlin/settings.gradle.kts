// examples/kotlin/settings.gradle.kts — Sprach-Wurzel-Wurzelprojekt der
// Kotlin-Beispiel-Clients (ADR-0087 Festlegung 3, ADR-0090 Festlegung 1/5).
// Ein Modul je Client (mirror der C#-Sprach-Wurzel, `examples/csharp/`), das
// erste war `http-client`, `sse-client` folgt (`ADR-0090`
// §Slice-Schnitt-Empfehlung, Zeile 3). `dependencyResolutionManagement`
// zentralisiert die Paket-Quelle für alle Module — Netz ist ohnehin
// Voraussetzung (NuGet-/Maven-Paketbezug, ADR-0090 Festlegung 5), `make
// gates` bleibt davon unberührt.
rootProject.name = "cdc-examples-kotlin"

dependencyResolutionManagement {
    repositories {
        mavenCentral()
    }
}

include("http-client")
include("sse-client")
