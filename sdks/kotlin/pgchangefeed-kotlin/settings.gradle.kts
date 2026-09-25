// sdks/kotlin/pgchangefeed-kotlin/settings.gradle.kts — root project of the
// Kotlin SDK package: a standalone project under sdks/kotlin/, not a
// submodule of examples/kotlin/. Unlike examples/kotlin/settings.gradle.kts
// (five modules, one language root), this project is a single flat Gradle
// module — one package, one Maven coordinate.
// `dependencyResolutionManagement` centralizes the package source, like
// examples/kotlin/settings.gradle.kts; network access is a prerequisite anyway
// (Maven Central package download), `make gates` is unaffected.
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
