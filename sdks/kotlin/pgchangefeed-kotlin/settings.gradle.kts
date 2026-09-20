// sdks/kotlin/pgchangefeed-kotlin/settings.gradle.kts — Wurzelprojekt des
// Kotlin-SDK-Pakets (ADR-0109 Festlegung 3: eigenständiges Projekt unter
// sdks/kotlin/, kein Untermodul von examples/kotlin/). Anders als
// examples/kotlin/settings.gradle.kts (fünf Module, eine Sprach-Wurzel) ist
// dieses Projekt ein einzelnes, flaches Gradle-Modul — ein Package, eine
// Maven-Koordinate (ADR-0109 Festlegung 1: "ein Package, nicht vier/fünf").
// `dependencyResolutionManagement` zentralisiert die Paket-Quelle, analog
// examples/kotlin/settings.gradle.kts; Netz ist ohnehin Voraussetzung
// (Maven-Central-Paketbezug), `make gates` bleibt davon unberührt.
rootProject.name = "pgchangefeed-kotlin"

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
