// examples/kotlin/grpc-client/build.gradle.kts — Programm-Manifest des
// Kotlin-gRPC-Clients (ADR-0090 Festlegung 1/2/4, LP1/LP2 des Slice-Plans,
// slice-103). Anders als http-/sse-/nats-client braucht dieses Programm
// einen Stub-Generator: das Gradle-Plugin `com.google.protobuf` orchestriert
// `protoc` und zwei Codegen-Plugins — `protoc-gen-grpc-java` (der
// Java-Service-Deskriptor, den der Kotlin-Stub referenziert) und
// `protoc-gen-grpc-kotlin` (der eigentliche Coroutine-Stub, Plugin-Id
// `grpckt`). Die `.proto` selbst liegt NICHT in diesem Verzeichnis — sie
// kommt erst im Docker-Bau nach `src/main/proto/` (examples/kotlin/Dockerfile),
// kopiert aus dem zusätzlichen, benannten Bau-Kontext `proto` (ADR-0090
// Festlegung 2). Dieser Pfad existiert deshalb nur innerhalb des Baus, nicht
// im committeten Baum (ADR-0090 Festlegung 3 — kein committeter Stub, keine
// committete Kopie der `.proto`).
//
// Gepinnte, exakte Versionen (Maven Central, gemessen 2026-09-17,
// maven-metadata.xml je Artefakt — deckungsgleich mit den in ADR-0090
// §Entscheidung Festlegung 4 gemessenen Kandidaten):
//   io.grpc:grpc-kotlin-stub        -> 1.5.0 (Coroutine-Stub-Laufzeit)
//   io.grpc:protoc-gen-grpc-kotlin  -> 1.5.0 (Codegen-Plugin `grpckt`, jdk8-Jar)
//   io.grpc:grpc-netty-shaded       -> 1.84.0 (gRPC-Transport)
// Die übrigen `io.grpc:*`-Artefakte werden über die BOM `io.grpc:grpc-bom`
// auf dieselbe 1.84.0-Zeile gehalten (Standard-Muster des grpc-java-Projekts
// für binär konsistente `io.grpc:*`-Versionen) — ohne die BOM zöge
// `grpc-kotlin-stub` transitiv `grpc-stub:1.62.2`, eine ältere Zeile als der
// hier gepinnte Transport. `grpc-kotlin-stub` bleibt davon unberührt: seine
// Versionierung läuft unabhängig von `io.grpc:grpc-bom` (siehe Kommentar in
// examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Main.kt zur
// übertragenen Form aus slice-102).
//
// Zusätzlich zur Werkzeugkette benötigt (nicht Teil der drei in ADR-0090
// namentlich benannten Bibliotheken, sondern deren Infrastruktur —
// vergleichbar mit dem C#-Pendant `Grpc.Tools`, das `protoc` und das
// C#-Codegen-Plugin intern bündelt; die Kotlin-/Gradle-Werkzeugkette braucht
// dieselben Bausteine als eigenständige, hier einzeln gepinnte Artefakte):
//   com.google.protobuf:protoc     -> 4.36.1 (protoc-Binary, Maven Central;
//                                     deckungsgleich mit protobuf-java/-kotlin)
//   io.grpc:protoc-gen-grpc-java   -> 1.84.0 (Codegen-Plugin `grpc`, erzeugt
//                                     den Java-Service-Deskriptor)
//   com.google.protobuf (Gradle-Plugin, Gradle Plugin Portal) -> 0.10.0
//   org.jetbrains.kotlinx:kotlinx-coroutines-core -> 1.11.0. Real beobachtet
//     (slice-103): `grpc-kotlin-stub`s eigenes POM erklärt
//     `kotlinx-coroutines-core-jvm` nur als `runtime`-Scope-Abhängigkeit,
//     nicht als `compile`/`api` — der generierte Coroutine-Stub referenziert
//     `kotlinx.coroutines.flow.Flow` aber bereits zur Kompilierzeit dieses
//     Moduls. Ohne diese explizite `implementation`-Abhängigkeit bricht
//     `compileKotlin` mit „Unresolved reference 'kotlinx'" ab — derselbe
//     Coroutine-Import, den auch `Main.kt` direkt benutzt (`runBlocking`,
//     `Flow.collect`).
//   com.google.protobuf:protobuf-java -> 4.36.1. Real beobachtet
//     (slice-103): `io.grpc:grpc-protobuf:1.84.0` zieht transitiv
//     `protobuf-java:3.25.9` — eine ältere Laufzeit-Zeile, als der oben
//     gepinnte `protoc` (4.36.1) Java-Quelltext dafür erzeugt (u. a.
//     `com.google.protobuf.GeneratedFile`, `RuntimeVersion`,
//     `GeneratedMessage.isStringEmpty`, alles erst ab der 4.x-Zeile
//     vorhanden). Ohne diese explizite Anhebung bricht `compileJava` des
//     erzeugten Stubs mit `cannot find symbol` ab — `protoc` und
//     `protobuf-java` müssen dieselbe Major-Zeile tragen, dasselbe Prinzip
//     wie die `io.grpc:grpc-bom`-Angleichung oben.
plugins {
    kotlin("jvm")
    application
    id("com.google.protobuf")
}

kotlin {
    jvmToolchain(21)
}

application {
    mainClass.set("cdcexamples.grpc.MainKt")
}

dependencies {
    implementation(platform("io.grpc:grpc-bom:1.84.0"))
    implementation("io.grpc:grpc-kotlin-stub:1.5.0")
    implementation("io.grpc:grpc-protobuf")
    implementation("io.grpc:grpc-stub")
    implementation("com.google.protobuf:protobuf-java:4.36.1")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.11.0")
    runtimeOnly("io.grpc:grpc-netty-shaded")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:4.36.1"
    }
    plugins {
        create("grpc") {
            artifact = "io.grpc:protoc-gen-grpc-java:1.84.0"
        }
        create("grpckt") {
            artifact = "io.grpc:protoc-gen-grpc-kotlin:1.5.0:jdk8@jar"
        }
    }
    generateProtoTasks {
        all().forEach { task ->
            task.plugins {
                create("grpc")
                create("grpckt")
            }
        }
    }
}

tasks.test {
    useJUnitPlatform()
}
