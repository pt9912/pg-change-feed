# ADR-0148: Kotlin-SDK nennt `io.grpc:grpc-api` im README; Upgrade-Trigger von ADR-0064 erfüllt

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (SDK-Packages),
[`LH-QA-OPS-005`](../../../spec/lastenheft.md) (Upgrade-Sicherheit),
[ADR-0145](0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md),
[ADR-0109](0109-kotlin-github-packages-drittes-sdk-package.md),
[ADR-0064](0064-lh-qa-ops-005-testansatz-korrektur.md),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.12,
Slice `sdk-0-6-kompatibilitaet-messen`,
Beobachtung
[`kein-echter-versionswechsel-upgrade-test`](../planning/observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md).

**Schärft:** [ADR-0064](0064-lh-qa-ops-005-testansatz-korrektur.md) (Re-Evaluierungs-Trigger 1;
ADR-0064 bleibt unverändert). Teil 1 (Kotlin) schärft keine Spec-Stelle und widerspricht keiner
Entscheidung (gemessen: `git grep -n -E 'implementation|\bapi\(' -- docs/plan/adr/0109-*.md docs/plan/adr/0110-*.md`
trifft nichts; der einzige Treffer in `0123` betrifft die Konsumentenseite, nicht den Scope des SDK).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

### Teil 1 — Kotlin: gRPC-Typen in der öffentlichen Signatur

**Gemessen** (`git grep -n 'io.grpc' -- sdks/kotlin/pgchangefeed-kotlin/src/main`, 2026-10-04):
`PgChangeFeedGrpcException` trägt `Status.Code` als öffentliches Feld und `StatusException` als
`cause`-Parameter; `PgChangeFeedGrpcClient` und `PgChangeFeedAdministrationClient` nehmen einen
`io.grpc.Channel` im öffentlichen Konstruktor. `build.gradle.kts` deklariert alle Bibliotheken
(Gson, gRPC, protobuf, coroutines, jnats) als `implementation`.

**Gemessen am veröffentlichten Artefakt** (Cloudsmith `…/pgchangefeed-kotlin/0.6.0/`, `.pom` und
`.module`, im Container gelesen, 2026-10-04): die POM führt nur `kotlin-stdlib` mit Scope
`compile`, alle übrigen Bibliotheken mit Scope `runtime`; die Gradle-Metadatei führt unter
`apiElements` nur `kotlin-stdlib`, die übrigen unter `runtimeElements`. Ein Konsument sieht die
gRPC-Typen deshalb nicht auf seinem Compile-Classpath (deckt sich mit der Messung des Slice:
der Kotlin-Gast brauchte `io.grpc:grpc-api` ausdrücklich; gilt für 0.5.0 wie 0.6.0).

Das README nennt dieses Verhalten bereits als Prinzip (§„Dependencies of the library“: die
Bibliotheken werden nicht durchgereicht, der Konsument ergänzt die, deren Typen er benutzt) und
führt coroutines, protobuf und Gson auf, **nicht** `io.grpc:grpc-api`.

## Entscheidung

**Teil 1.** Wir wählen **Option B: ein README-Eintrag**, kein `api`-Umbau. Der Block
„Dependencies of the library“ des Kotlin-README erhält `io.grpc:grpc-api` (Version der
gRPC-BOM des Builds, `1.84.0`, abgelesen aus `build.gradle.kts`) mit dem Anlass „für `Status`,
`StatusException` und `Channel` der gRPC-Fehlertypen und -Konstruktoren“. Keine Änderung an
`build.gradle.kts`, kein neues Release: die Zeile wirkt mit dem nächsten ohnehin fälligen
Release (ob Cloudsmith/GitHub Packages das README tragen, ist nicht geprüft). Slice-Zuschnitt:
**kein eigener Slice**; die Zeile geht als Ein-Datei-Änderung des Implementers oder mit der
nächsten Kotlin-SDK-Änderung. Release: keiner ohne Nutzerfreigabe.

**Teil 2.** Der Re-Evaluierungs-Trigger 1 von ADR-0064 ist **erfüllt**: eine Release-Historie
existiert (`git tag` führt `v0.1.0` bis `v0.6.0`, gemessen 2026-10-04). Der Slice
`sdk-0-6-kompatibilitaet-messen` erfüllt ihn **nicht**: er misst SDK-Stände gegeneinander und
SDKs gegen einen Alt-Server (Fehlertypen, Diagnose), nicht den Datenstand über einen
Server-Wechsel. Die Folge-Entscheidung ist die **kleinste tragfähige Form** des
Alt-gegen-Neu-Vergleichs.

Der bestehende Upgrade-Tausch-Rundlauf (`make test-integration`, `--force-recreate` desselben
`:dev`-Images) misst: Datenstand vor/nach dem Tausch über `cdc.changes` identisch lesbar, danach
eingefügte Zeile erfasst (ADR-0064). Er misst **nicht**, dass ein *anderer Server-Build* den
Datenstand liest und fortsetzt. Neu wird **eine Phase** ergänzt: der Feed-Container startet aus
dem veröffentlichten Image 0.5.0 (Pin-Form `<Tag>@<Index-Digest>` wie `SDK_ALTSERVER_IMAGE` in
`tools/harness/run-sdk-altserver-tests.sh`; dieser Server läuft gegen das aktuelle Schema,
Befund B0 des Slice), erfasst Zeilen und wird per `--force-recreate` durch das `:dev`-Image
ersetzt; erwartet wird derselbe Prüfumfang wie im bestehenden Rundlauf (identisch lesbar,
Erfassung setzt fort). Das Schema bleibt **konstant** (das aktuelle); ein Schemawechsel über
Versionen bleibt ungemessen — **akzeptiertes Negativ**: der Schema-Stand des 0.5.0-Baums wäre
nur über Tag-Checkout mit eigenem d-migrate-Rollout herstellbar (Aufwand ohne bisher
beobachteten Fehler, Option C unten). Slice-Zuschnitt (nur benannt):
`upgrade-versionswechsel-alt-image`, Träger `tools/harness/` (Erweiterung des Altserver- oder
des Integrations-Runners; die Wahl trifft der Planner). Die Aussage „der Alt-Build läuft mit
konstantem Schema durch den Tausch“ ist im Slice **erwartet**, bis die Phase gelaufen ist.

## Verglichene Alternativen

### Teil 1

| Option | Pro | Contra |
|---|---|---|
| A — `api` für die gRPC-Typen der Signatur | Compile-Classpath des Konsumenten trägt sie ohne eigene Zeile | ändert die veröffentlichten Scopes (hergeleitet: `apiElements` und POM-`compile`; ob `kotlin("jvm")` ohne `java-library` die Konfiguration `api` bereitstellt, ist nicht gemessen — *erwartet*, ein Lauf von `make sdk-pack-kotlin` entscheidet); erzwingt ein Release; inkonsistent, solange protobuf, coroutines und Gson (ebenfalls in öffentlichen Signaturen: `ByteString`, `Flow`, `JsonElement`) `implementation` blieben |
| **B — README-Eintrag (gewählt)** | folgt dem bereits dokumentierten Prinzip des README; kein Release, keine Scope-Änderung | der Konsument setzt die Zeile selbst (wie bei den drei genannten) |
| C — A und B | doppelte Absicherung | Aufwand von A ohne Nutzen gegenüber B |
| D — nichts tun | kein Aufwand | der Übersetzungsfehler beim ersten Gebrauch bleibt unerklärt (gemessen im Slice) |

### Teil 2

| Option | Pro | Contra |
|---|---|---|
| A — „Trigger nicht erfüllt“ | keine Arbeit | falsch: die Bedingung des Triggers (Tags, Release-Historie) ist erfüllt (gemessen); die Lücke bliebe ohne Grund offen |
| **B — eine Alt-Image-Phase bei konstantem Schema (gewählt)** | schließt den Kern der Lücke (anderer Server-Build liest den Datenstand); nutzt Pin-Mechanismus und Runner, die es gibt | Schemawechsel bleibt ungemessen (benannt) |
| C — voller Versionswechsel mit Schema des 0.5.0-Tags | misst auch den Schemawechsel | Tag-Checkout und eigener Rollout im Runner; kein beobachteter Fehler, der den Aufwand trüge |

## Konsequenzen

- Positiv: der Kotlin-Konsument bekommt die fehlende Zeile; die Beobachtung „kein echter
  Versionswechsel“ bekommt einen konkreten Träger statt weiteren Wartens.
- Negativ: Teil 1 lässt die Compile-Scopes unverändert (bewusst); Teil 2 misst den
  Schemawechsel nicht (akzeptiertes Negativ, oben begründet).
- Folgepflicht: README-Zeile (Teil 1, Implementer, kein Release). Slice
  `upgrade-versionswechsel-alt-image` (Teil 2, der Planner legt ihn an). Die Beobachtung
  `kein-echter-versionswechsel-upgrade-test` nennt diese ADR als Träger der Auflösung, ihr
  Zähler bleibt 2×; sie wird aufgelöst, wenn der Slice gelaufen ist. `harness/README.md`
  §Sensors (Zeilen `make test-integration` und `make test-sdk-altserver`) zieht der Slice nach.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `git grep -n 'grpc-api' -- sdks/kotlin/pgchangefeed-kotlin/README.md` | Teil 1: mindestens ein Treffer im Block „Dependencies of the library“ — *erwartet* nach dem Implementer-Zug, nicht gefahren | — |
| Alt-Image-Phase (neu) | Teil 2: Datenstand des 0.5.0-Builds nach dem Tausch auf `:dev` identisch lesbar, Erfassung setzt fort — *erwartet*, nicht erprobt | `make test-integration` bzw. `make test-sdk-altserver` (kein Gate) |

## Re-Evaluierungs-Trigger

Teil 1: ein Konsumenten-Befund, dass die README-Zeile nicht genügt, oder eine Entscheidung, die
Scopes der veröffentlichten Metadaten für alle Bibliotheken neu zu fassen. Teil 2: ein
beobachteter Fehler eines Betriebs-Upgrades, der am konstanten Schema vorbeiläuft (dann
Option C).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted — Anlass: offene Fragen V-5/F-4 und Re-Evaluierung ADR-0064 aus der Closure des Slice `sdk-0-6-kompatibilitaet-messen` | Slice `sdk-0-6-kompatibilitaet-messen` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0148` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
