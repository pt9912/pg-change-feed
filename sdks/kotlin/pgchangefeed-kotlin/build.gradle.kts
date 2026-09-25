// sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts — Maven metadata and build
// configuration of the Kotlin SDK package. A standalone project, not a
// submodule of examples/kotlin/; it imports nothing from a private tree of
// this repository.
//
// Kotlin Gradle plugin version measured against Maven Central
// (`repo1.maven.org/maven2/org/jetbrains/kotlin/kotlin-gradle-plugin/maven-metadata.xml`,
// 2026-09-20): the latest version is 2.4.20 (`lastUpdated` 2026-09-07), the
// same as in `examples/kotlin/build.gradle.kts` (measured 2026-09-17).
//
// `maven-publish` is the built-in Gradle core plugin, not a third-party plugin
// such as `com.vanniktech.maven-publish`. The `publishing` block carries the
// Maven coordinate and two `repositories{}` entries, one publish task each:
//   GitHubPackages -> https://maven.pkg.github.com/pt9912/pg-change-feed,
//     credentials from `GITHUB_ACTOR`/`GITHUB_TOKEN`
//     (`publishMavenPublicationToGitHubPackagesRepository`)
//   Cloudsmith -> https://maven.cloudsmith.io/pt9912/pg-change-feed/,
//     credentials from `CLOUDSMITH_USERNAME`/`CLOUDSMITH_API_KEY`
//     (`publishMavenPublicationToCloudsmithRepository`)
// `GITHUB_ACTOR` is a default value GitHub Actions provides; the other three
// variables are set explicitly by the calling workflow
// (`.github/workflows/sdk-kotlin-release.yml`, one job per target). No
// credential value appears in this file. The aggregate `publish` task runs
// both targets and fails without the Cloudsmith values, so the workflow calls
// the per-target tasks. Outside that workflow (locally, in `make
// sdk-pack-kotlin`) all four variables stay empty; `make sdk-pack-kotlin`
// builds only the `pack-export` Docker stage (`test`/`build`, no publish
// upload). `sdks/kotlin/Dockerfile` also carries a `publish` stage (built on
// `build`) that only the workflow builds and starts, with the credentials
// passed by `docker run -e` at run time.
//
// `com.google.code.gson:gson` is the JSON library of the HTTP client, the SSE
// client and the NATS client, at the same pinned version as
// `examples/kotlin/nats-stream-client/build.gradle.kts` (license and transitive
// dependencies checked there: only
// `com.google.errorprone:error_prone_annotations` in the `compile` scope,
// Apache-2.0). Neither `java.net.http` nor the JDK standard library carries a
// public JSON decoder. Measured 2026-09-20 (Maven Central maven-metadata.xml,
// repo1.maven.org/maven2/com/google/code/gson/gson/maven-metadata.xml): the
// latest version is `2.14.0` (`<latest>`/`<release>`).
//
// gRPC stream client: the same coordinates as
// `examples/kotlin/grpc-client/build.gradle.kts`, measured again on 2026-09-20
// (Maven Central maven-metadata.xml per artifact, Gradle Plugin Portal for the
// `com.google.protobuf` plugin):
//   io.grpc:grpc-kotlin-stub        -> 1.5.0 (`<latest>`/`<release>` of the
//     Maven metadata show a commit-hash entry — a CI snapshot artifact, not a
//     release, lastUpdated 2025-09-16; 1.5.0 is the latest regular release)
//   io.grpc:protoc-gen-grpc-kotlin  -> 1.5.0 (same anomaly, same resolution)
//   io.grpc:grpc-netty-shaded       -> 1.84.0
//   io.grpc:grpc-bom                -> 1.84.0
//   io.grpc:protoc-gen-grpc-java    -> 1.84.0
//   org.jetbrains.kotlinx:kotlinx-coroutines-core -> 1.11.0
//   com.google.protobuf (Gradle plugin, Gradle Plugin Portal) -> 0.10.0
//   com.google.protobuf:protoc -> 4.36.2 (`examples/kotlin/grpc-client`, measured
//     2026-09-17, has 4.36.1. The Maven metadata also carry a `21.0-rc-1` entry,
//     lastUpdated 2026-09-17, a release candidate of a new version numbering;
//     4.36.2 is the latest stable version, obtained from
//     repo1.maven.org/maven2/com/google/protobuf/protoc/4.36.2/)
//   com.google.protobuf:protobuf-java -> 4.36.2 (the same version as `protoc`;
//     both must be on the same major line. `io.grpc:grpc-protobuf:1.84.0` still
//     pulls in `protobuf-java:3.25.9` transitively, so the explicit raise stays
//     necessary)
//
// NATS stream client: `io.nats:jnats` measured 2026-09-22 (Maven Central
// maven-metadata.xml, repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml):
// `2.26.3` (`<latest>`/`<release>`), the same version as in
// `examples/kotlin/nats-stream-client/build.gradle.kts`. No second JSON decoder
// is needed: the pinned `com.google.code.gson:gson:2.14.0` also deserializes
// the NATS messages, which share their schema with the SSE events.
//
// Version `0.2.2` (SemVer): raised deliberately at each release, not derived
// automatically; the top-level `version` (jar name, checked against the release
// tag) and the `version` of the publication (POM) carry the same value. The
// POM fields `name`/`description`/`url`/`licenses`/`scm` in the `publishing`
// block are the package metadata published with the artifacts; the jar
// carries no README. They contain no internal identifiers. The `java` block
// adds a sources jar to the publication.
plugins {
    kotlin("jvm") version "2.4.20"
    `maven-publish`
    id("com.google.protobuf") version "0.10.0"
}

group = "io.github.pt9912"
version = "0.2.2"

kotlin {
    jvmToolchain(21)
}

java {
    withSourcesJar()
}

dependencies {
    implementation("com.google.code.gson:gson:2.14.0")
    implementation(platform("io.grpc:grpc-bom:1.84.0"))
    implementation("io.grpc:grpc-kotlin-stub:1.5.0")
    implementation("io.grpc:grpc-protobuf")
    implementation("io.grpc:grpc-stub")
    implementation("com.google.protobuf:protobuf-java:4.36.2")
    implementation("org.jetbrains.kotlinx:kotlinx-coroutines-core:1.11.0")
    implementation("io.nats:jnats:2.26.3")
    runtimeOnly("io.grpc:grpc-netty-shaded")
    testImplementation("org.jetbrains.kotlin:kotlin-test-junit5:2.4.20")
    testRuntimeOnly("org.junit.platform:junit-platform-launcher:6.1.3")
}

// The `.proto` source is not in the committed tree — it arrives only in the
// Docker build, at `src/main/proto/` (`sdks/kotlin/Dockerfile`), copied from
// the additional named build context `proto`
// (`docker build --build-context proto=proto …`, the same pattern as
// `examples/kotlin/Dockerfile`/`harness/mk/examples.mk`). Without that context
// the build stops at the `COPY` line of the Dockerfile — there is no silent
// fallback.
protobuf {
    protoc {
        artifact = "com.google.protobuf:protoc:4.36.2"
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

// Integration source set: the real-server tests need a running server
// instance (compose network). They do not run in the network-free unit run
// (`./gradlew test`/`build`, `make sdk-pack-kotlin`; the `integrationTest`
// task is deliberately not attached to `check`) but only through
// `make test-sdk-kotlin-integration`
// (tools/harness/run-sdk-kotlin-integration-tests.sh), which sets a `--tests`
// filter on a test class name per phase (explicit phase selection, no silent
// exclusion); phase addresses and auth arrive as environment variables that
// the test JVM inherits.
val integrationTestImplementation by configurations.creating {
    extendsFrom(configurations.testImplementation.get())
}

// Resolvable runtime configuration of the integration test class path: the
// runtimeOnly dependencies of the main module (grpc-netty-shaded) are not
// resolvable themselves — this configuration inherits them (and the test
// runtime) and is resolved in the task below.
val integrationTestRuntimeClasspath by configurations.creating {
    extendsFrom(
        configurations.getByName("integrationTestImplementation"),
        configurations.testRuntimeOnly.get(),
        configurations.runtimeOnly.get(),
    )
    isCanBeResolved = true
    isVisible = false
}

sourceSets {
    create("integrationTest") {
        kotlin.srcDir("src/integrationTest/kotlin")
    }
}

dependencies {
    // The main classes (including the gRPC stubs generated in the build) are
    // the subject of the integration tests: the compiled client assembly is
    // what is tested.
    "integrationTestImplementation"(sourceSets["main"].output)
}

val integrationTest by tasks.registering(Test::class) {
    description =
        "Real-server integration test of the client surfaces — needs a running server instance; one call per phase with a --tests filter."
    group = "verification"
    testClassesDirs = sourceSets["integrationTest"].output.classesDirs
    // The run class path carries the runtime pieces explicitly: the
    // implementation dependencies (via integrationTestImplementation), the
    // test runtime (junit-launcher, testRuntimeOnly) and the runtimeOnly
    // dependencies of the main module (grpc-netty-shaded — without it the gRPC
    // channel ends in a ProviderNotFoundException).
    classpath = sourceSets["integrationTest"].output +
        configurations["integrationTestRuntimeClasspath"]
    useJUnitPlatform()
    // The runner markers (READY/RECEIVED/REJECTED) are read by the runner via
    // `docker logs` while the test is running — the standard streams must be
    // passed through live (JVM stdout buffering).
    testLogging {
        events("passed", "failed")
        showStandardStreams = true
    }
}

publishing {
    publications {
        create<MavenPublication>("maven") {
            groupId = "io.github.pt9912"
            artifactId = "pgchangefeed-kotlin"
            version = "0.2.2"
            from(components["java"])
            pom {
                name.set("PG Change Feed Kotlin client")
                description.set(
                    "Kotlin/JVM client library for PG Change Feed: read PostgreSQL changes over HTTP and receive " +
                        "them live over gRPC, SSE and NATS.",
                )
                url.set("https://github.com/pt9912/pg-change-feed")
                licenses {
                    license {
                        name.set("MIT")
                        url.set("https://github.com/pt9912/pg-change-feed/blob/main/LICENSE")
                    }
                }
                scm {
                    url.set("https://github.com/pt9912/pg-change-feed")
                }
            }
        }
    }
    repositories {
        maven {
            name = "GitHubPackages"
            url = uri("https://maven.pkg.github.com/pt9912/pg-change-feed")
            credentials {
                username = System.getenv("GITHUB_ACTOR")
                password = System.getenv("GITHUB_TOKEN")
            }
        }
        maven {
            name = "Cloudsmith"
            url = uri("https://maven.cloudsmith.io/pt9912/pg-change-feed/")
            credentials {
                username = System.getenv("CLOUDSMITH_USERNAME")
                password = System.getenv("CLOUDSMITH_API_KEY")
            }
        }
    }
}
