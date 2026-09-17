// examples/kotlin/settings.gradle.kts — Sprach-Wurzel-Wurzelprojekt der
// Kotlin-Beispiel-Clients (ADR-0087 Festlegung 3, ADR-0090 Festlegung 1/5).
// Ein Modul je Client (mirror der C#-Sprach-Wurzel, `examples/csharp/`):
// `http-client`, `sse-client`, `nats-client` (`ADR-0090`
// §Slice-Schnitt-Empfehlung, Zeile 4, slice-101), `grpc-client` (Zeile 6,
// slice-103). `dependencyResolutionManagement` zentralisiert die
// Paket-Quelle für alle Module — Netz ist ohnehin Voraussetzung
// (NuGet-/Maven-Paketbezug, ADR-0090 Festlegung 5), `make gates` bleibt
// davon unberührt. `gradlePluginPortal()` trägt zusätzlich das
// Gradle-Plugin `com.google.protobuf` (`grpc-client/build.gradle.kts`) —
// die drei anderen Module brauchen kein Fremd-Plugin.
rootProject.name = "cdc-examples-kotlin"

pluginManagement {
    repositories {
        gradlePluginPortal()
        mavenCentral()
    }
}

dependencyResolutionManagement {
    repositories {
        mavenCentral()
    }
}

include("http-client")
include("sse-client")
include("nats-client")
include("grpc-client")
