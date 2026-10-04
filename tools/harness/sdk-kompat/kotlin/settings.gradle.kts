rootProject.name = "kompat-gast"

// Die veröffentlichte Bibliothek kommt aus dem öffentlichen Cloudsmith-Repository
// (Pfad der Installationsanleitung des Kotlin-SDK), ihre eigenen Abhängigkeiten
// aus Maven Central.
dependencyResolutionManagement {
    repositories {
        mavenCentral()
        maven { url = uri("https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/") }
    }
}
