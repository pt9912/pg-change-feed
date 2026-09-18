# Harness

---

## Purpose

Dieser Harness verbindet bestehende Spezifikationen, ADRs,
Planning-Dokumente und Gates. Er ist **kein Ersatz** für `spec/` oder
`docs/`, sondern ein **Einstiegspunkt** für Menschen und AI-Code-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, **gewinnt die
kanonische Quelle**, und diese Datei wird angepasst.

Strukturregeln (Verzeichniskonvention, ID-Schemata, Modus-Deklarationen
pro Sub-Area, Zusatzklassen für Sensors-Bindung) sowie Adaptionen ggü.
der adoptierten Baseline leben in [`conventions.md`](conventions.md).
Diese Datei dupliziert sie nicht.

## Source precedence

| Rang | Datei | Charakter |
|---|---|---|
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | vertraglich abnahmebindend |
| 2 | [`spec/pflichtenheft.md`](../spec/pflichtenheft.md) | technisch fortschreibbar |
| 3 | [`spec/architecture.md`](../spec/architecture.md) | Komponenten/Sequenzen, meilensteinfrei |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Wellen-Sequenz |
| 6 | `docs/user/*` *(falls vorhanden)* | Operations, Quality, Releasing | <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
| 7 | [`README.md`](../README.md) | Projekt-Überblick |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 9 | diese Datei | Harness-Einstieg |

> Die Ränge 1–3 sind die **drei Spec-Straten** — Vertrag, Technik, Sicht —,
> und alle drei sind obligatorisch (Baseline-Regelwerk
> `grundlagen-referenz-richtung.md` §Spec-Straten). **Adaption ist die
> Zwei-Straten-Form, nicht die Drei-Straten-Form**: Wer Rang 2 streicht,
> deklariert das als `MR-<NNN>` in [`conventions.md`](conventions.md) und
> nummeriert neu (dann acht Ränge).

## Guides (Feedforward-Quellen)

<!--
Was lenkt den Agenten *vor* der Handlung? Pointer, kein Inhalt.
-->

| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen, IDs, Akzeptanzkriterien |
| [`spec/pflichtenheft.md`](../spec/pflichtenheft.md) | technische Details, Defaults |
| [`spec/architecture.md`](../spec/architecture.md) | Komponenten, Schichten, Constraints |
| [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |
| `.harness/skills/reviewer.md` | Reviewer-Skill: HIGH-Liste, Kategorien-Regeln, Negativbefund-Pflicht, Output-Schema (Modul 10) — nächste Rolle nach Schritt 8 des Minimal Agent Workflow, nicht Teil der Implementer-Eingabe |
| `.harness/baseline/<tag>/regelwerk/` (vendored; `README.md` = Index) | adoptiertes Betriebsregelwerk in Agenten-Kurzform — **präsente nachschlagbare Vertiefung**, pro Entscheidung abschnittsweise (siehe [`AGENTS.md`](../AGENTS.md) §1); derivativ, Stand/Tag siehe [`conventions.md`](conventions.md) §Baseline |
| `.harness/baseline/<tag>/templates/` (vendored, parallel) | Referenz-Form der Skelette, auf die das Regelwerk mit `../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu `regelwerk/`); Vorlagen zum Kopieren-und-Ausfüllen |

## Sensors (Feedback-Gates)

<!--
WICHTIG: Nur Befehle aufzählen, die im Makefile *existieren*.
Halluzinierte Gates sind die häufigste Form von Harness-Lüge (Modul 13).

Drei Spalten — kein Lauf-Status:
- Target:  der Make-Befehl.
- Vertrag: was prüft das Gate (was wäre verletzt, wenn es rot wird).
- Bindung: strukturelle Referenzen — Carveout-ID (`CO-<NNN>`),
  Slice-ID, Schwelle, Image-Hash, ADR-ID. NICHT der Lauf-Status,
  sondern was das Gate *strukturell trägt*.

Lauf-Wahrheit pro Commit liegt in CI (Badge/Dashboard), nicht hier
(`harness/README.md` ist Rang 9 in der Source Precedence).
Strukturell rote Gates (dauerhaft rot) bekommen einen Carveout in
`docs/plan/carveouts/CO-<NNN>-…` mit Auflösungs-Trigger und Folge-Slice
(Modul 7); die Bindung-Spalte verweist auf die `CO-<NNN>`-ID, die
Begründung lebt im Carveout, nicht hier.

Bei d-check-Einsatz (≥ v0.73.0) deckt Modul `reviews` (Ziel `doc-reviews`,
Review-Report-Deckung für `done/`-Slices mit Review-DoD-Haken) die
Code→Review-Kante ab, Modul `planning` (Ziel `doc-planning`,
Planning-Lifecycle-Konsistenz) die Verify→Closure-Kante — beide aus dem
Lebenszyklus-Diagramm in Modul 1. Nur eintragen, wenn das Ziel im
Makefile existiert (siehe oben).

NICHT-GATES: Ein Target, das der Agent braucht, aber das nichts über den
Zustand des Repos urteilt — es *bewegt* (Slice-Move), *misst* (Latenz) oder
*sagt*, was ein schreibender Lauf täte — steht in der zweiten Tabelle und
trägt `kein Gate` IN DER ZEILE SELBST, in der Spalte, die hier die Bindung
führt — nicht in Prosa daneben. Weglassen ist nur für Targets richtig, die
niemand braucht (Modul 13 §Vorhanden ≠ behauptet).

WÄCHST DIE SEKTION: Die Tabelle bleibt klein, die Prosa darunter nicht.
Braucht ein Gate mehr als EINEN SATZ — Deckungsgrenze, Ausgabe-Bedeutung,
Exit-Codes, Abbruch-Bedingungen —, wandert das nach. Ob der Überhang schon
unter der Tabelle steht oder in die Zelle gedrängt wurde, ist dieselbe Sache:
eine Zelle, die zum Absatz geworden ist, ist der Fund, nicht die Ausnahme.
`harness/sensors/<target>.md`, und die **Target-Zelle wird zum Link darauf**
— wie die `MR`-Zelle im Adaptions-Block. Der Link ist kein Komfort: Er ist die
einzige Fassung dieser Zuordnung, die der Link-Sensor prüft. Eine bloße
Namenskonvention (`make X` -> `sensors/X.md`) bleibt still grün, wenn die
Datei verschwindet und die Zeile stehen bleibt. Seine Grenze: Er prüft EINE
Richtung — ob das Ziel existiert; eine Datei ohne Index-Zeile und eine Zeile
auf die falsche Datei bleiben still grün, und geprüft wird nur, wo ein
Link-Sensor über `harness/` läuft. Kein `sensors/done/`:
ein retiriertes Gate verschwindet, `git` hält seine Geschichte. Was das
Werkzeug selbst deckt (welcher Test welche Hälfte trägt), gehört NICHT
dorthin, sondern in seine ADR/Spec-Zeile/seinen Skriptkopf.
-->

| Target | Vertrag | Bindung |
|---|---|---|
| `make baseline-verify` | verifiziert die vendored Baseline netzlos (Integrität + Vollständigkeit gegen `SHA256SUMS`) | [`harness/sensors/baseline-verify.md`](sensors/baseline-verify.md) |
| `make docs-check` | kaputte Referenzen in der Markdown-Doku (links, anchors, ids, matrix, versions, structure, hostpaths, tracked) | [`harness/sensors/docs-check.md`](sensors/docs-check.md) |
| `make a-check` | prüft die Hexagon-Schichten-Edges aus `.a-check.yml` gegen den Go-Baum — netzlos, read-only, digest-gepinntes Release-Image | [`ADR-0041`](../docs/plan/adr/0041-a-check-maschinenform-architekturpruefung.md) · [`harness/sensors/a-check.md`](sensors/a-check.md) |
| `make commit-traceability` | Commit-Message-Traceability: je Message der Range ≥ 1 `LH-*`-/`ADR-*`-Kennung (d-check Modul `commits`, Befund `commit-untraceable`) und keine `SPEC-*`/`ARC-*`-Kennung im Betreff (`tools/harness/commit-traceability.sh`); Standing-Gate über die letzten 5 Commits (`RANGE=base..head` überschreibt) | [`ADR-0045`](../docs/plan/adr/0045-commit-traceability-standing-gate.md) · seit slice-006 |
| `make coverage-gate` | Go-Test-Coverage über die netzlos prüfbare Fläche (`./internal/...`+`./cmd/...`+`./gen/...` ohne die Pakete, deren Testlauf einen externen Dienst voraussetzt) gegen `THRESHOLD` — bootstrap-aware Gate, Rampe Einstieg 70 % → Endstufe 80 %; geltende Stufe in [`harness/mk/coverage.mk`](mk/coverage.mk) | [`ADR-0071`](../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) · [`ADR-0054`](../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) · [`ADR-0077`](../docs/plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md) · [`harness/sensors/coverage-gate.md`](sensors/coverage-gate.md) · seit slice-049 |
| `make generated-sync` | prüft, dass der committete Protobuf-/gRPC-Code byte-gleich der Ausgabe des **gepinnten** Generators (Dockerfile-Stufe `proto`) aus den committeten `.proto`-Quellen ist; erzeugt wird in ein Temp-Verzeichnis, der Arbeitsbaum bleibt unberührt | [`harness/sensors/generated-sync.md`](sensors/generated-sync.md) · seit slice-090 |
| `make gates` | alle inneren Gates (baseline-verify, docs-check, a-check, commit-traceability, coverage-gate, generated-sync), Nachweis-Stempel zuletzt | — |
| `<make-target>` | volle Closure | Image-Hash `sha256:…` (Modul 14) |

**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate:**

| Target | Tut was | Bindung |
|---|---|---|
| `make image` | baut das OCI-Image (Multi-Stage, digest-gepinnt); **Image-Digest** nach `harness/image-hash.txt` — der Digest ist der **Lauf-Beleg** des letzten offiziellen `make image`-Laufs (builder- und lauf-gebunden, kein Inhalts-Fingerabdruck); ein Digest-Vergleich über Umgebungen/Läufe ist kein Staleness-Beweis, Inhalts-Streits entscheidet der sha256 des extrahierten Binaries (Container-Export). Ein Zug, der Build-Kontext-Dateien ändert, läuft `make image` vor seiner Closure; der Digest-Commit entfällt bei unverändertem Digest | kein Gate, [`ADR-0044`](../docs/plan/adr/0044-image-beleg-semantik.md), [`ADR-0042`](../docs/plan/adr/README.md) |
| `make image-stale` | meldet Base-Image-Drift: FROM-Digests des Dockerfile gegen die aktuellen Registry-Digests derselben Tags, plus Existenz des nächsten Major-Tags; braucht Netz | kein Gate, [`ADR-0039`](../docs/plan/adr/README.md) (Update = bewusster Digest-Commit, Modul 14) |
| `make doc-tracked` | isolierter Einzel-Lauf des `tracked`-Moduls (Getrackt-Status auflösbarer, existierender Link-/Bild-Ziele gegen den git-Index, `--enable tracked` mit allen anderen Modulen `--disable`) — das Modul läuft bereits im `modules:`-Bündel mit und damit in `make docs-check`/`make gates`; dieses Ziel bleibt daneben als isoliertes Diagnose-Werkzeug bestehen, netzlos, braucht `.git` im Mount | kein Gate, [`harness/sensors/docs-check.md`](sensors/docs-check.md) §Bindung · seit slice-d-check-tracked-modul |
| `make doc-trace` | gibt die Requirements Traceability Matrix (RTM) auf stdout aus (`--trace`, **kein** eigenes Modul, advisory — Exit 0 auch mit Requirements-Waisen); `trace.requirements.id-pattern` in `.d-check.yml` verdrahtet die RTM-Quelle auf die `LH-(FA|QA)-[A-Z]{3}-\d{3}`-Kennungskonvention dieses Repos statt d-checks generischen Default, `trace.coverage` liest zusätzlich `docs/user/e2e-abdeckung.md` als kuratierte Coverage-Dimension, damit eine Anforderung mit ausschließlichem E2E-Nachweis nicht fälschlich als Waise erscheint (real gemessen: 76 Anforderungen, 9 Waisen ohne `trace.coverage`, 7 Waisen mit). `make doc-complete`/`--require-complete` bleibt bewusst **nicht** in `GATE_CHECKS`/`make gates` — die verbleibenden 7 Waisen (`LH-FA-CFG-006`, `LH-FA-CON-002`, `LH-FA-DAT-002`, `LH-FA-DAT-003`, `LH-FA-SST-001`, `LH-FA-SST-005`, `LH-QA-REL-004`) sind unsanierte Bestandslücken, deren Sanierung ein eigener inhaltlicher Vorgang wäre, kein Konfigurations-Vorgang | kein Gate, wie `make image-stale`, [`harness/sensors/docs-check.md`](sensors/docs-check.md) §Grenze · seit slice-d-check-trace-rtm |
| `make proto-generate` | erzeugt den Go-Code des Live-Change-Streams aus `proto/cdc/stream/v1/changestream.proto` — Docker-only: die Dockerfile-Stufe `proto` trägt `protoc` und die gepinnten Plugins `protoc-gen-go`/`protoc-gen-go-grpc`; die Folgestufe `proto-export` (seit slice-104) kopiert die `.proto`-Quelle per `COPY` hinein und ruft `protoc` **zur Build-Zeit** auf (kein Bind-Mount, kein `--user`-Workaround), ihr `ENTRYPOINT` gibt das Erzeugnis als `tar`-Stream über stdout aus; `tools/harness/proto-generate.sh` baut die Stufe und extrahiert **host-seitig** über `docker run --rm --network none <image> | tar -x -C .` — die Pipe läuft explizit unter `bash` mit `set -o pipefail`, damit ein `docker run`-Fehlschlag den Lauf sichtbar rot macht statt hinter `tar`s Exit-Code zu verschwinden (`AGENTS.md` §3.9; das Makefile-Rezept selbst liefe unter `/bin/sh`, dort `dash` ohne `pipefail`) — die extrahierten Dateien gehören dadurch automatisch dem aufrufenden Nutzer. Die erzeugten Dateien liegen committet im Baum und werden von `make test`/`make image` mitkompiliert; der Generator läuft nur, wenn sich die `.proto`-Quelle ändert | kein Gate, [`ADR-0060`](../docs/plan/adr/0060-grpc-streaming-mechanismus.md) |
| `make test` | Unit-Tests mit Race-Detector im gepinnten Toolchain-Container (netzlos; Modul-Cache im Docker-Volume, Vorbereitung: `make mod-download`); `TOOLCHAIN_RACE_IMAGE` ist Debian- statt Alpine-basiert, weil `go test -race` einen C-Compiler zum Linken braucht (`gcc`), den das Alpine-Image nicht trägt · seit slice-037 (`ADR-0050` Fitness Function: `mapper.Assembler` gegen gleichzeitigen Zugriff aus zwei Goroutinen) | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) |
| `make test-store` | Adapter-Tests gegen reale PostgreSQL im Testcontainer (gepinnte Digests: Toolchain + PostgreSQL; Schema-Rollout über d-migrate vor dem Testlauf — dieselbe Kette wie der E2E-Lauf; DB-Daten im Container); legt das `-coverprofile` des Store-Teils der DB-Adapter-Coverage ab | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md), [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md); DB-Adapter-Coverage ([`sensors/db-adapter-coverage.md`](sensors/db-adapter-coverage.md), [`ADR-0071`](../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) |
| `make test-replication` | Replication-Stream-Tests gegen reale PostgreSQL mit Publication und Logical Replication Slot (`wal_level=logical`, gepinnte Digests: Toolchain + PostgreSQL; Schema-Stand des Tier-Laufs über denselben d-migrate-Rollout wie der Betrieb, `tools/schema/apply-rollout.sh`; DB-Daten im Container); legt das `-coverprofile` des Replication-Teils der DB-Adapter-Coverage ab, mergt beide Profile zur Zahl und prüft sie gegen die Schwelle | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md); DB-Adapter-Coverage ([`sensors/db-adapter-coverage.md`](sensors/db-adapter-coverage.md), [`ADR-0071`](../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) |
| `make test-notify` | `natsnotify`-Adapter-Tests gegen einen echten NATS-Server im Testcontainer (gepinnter Digest, `tools/harness/run-notify-tests.sh`): realer Publish-Erfolgsbeleg über eine Subscription auf `cdc.changes.>` — Subjekt-Form und leerer Payload (`SPEC-017`); der Konstruktions-, Empty-Source- und Fehlerklassen-Wrapping-Teil läuft bereits ohne Server in `make test` (Verbindungsfehler bleiben auf Loopback, auch unter `--network none`) | kein Gate, [`ADR-0055`](../docs/plan/adr/0055-nats-change-notification-wecksignal.md), [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) · seit slice-052 |
| `make test-integration` | Compose-Integrationstest gegen die Compose-Umgebung (`compose.yaml`), inhaltlich über den ursprünglichen MVP-Zuschnitt hinausgewachsen: Umgebung frisch hochfahren → Schema-Rollout über d-migrate ([`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)) → Rollen-DSN-Verifikation direkt gegen die Compose-PostgreSQL-Instanz (`LH-QA-SEC-001`…`003`, [`ADR-0047`](../docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md): ein cdc_reader-Login scheitert an einem schreibenden Aufruf, eine Login-Identität ohne REPLICATION-Attribut scheitert an einer Replication-Verbindung, eine cdc_capture-Login-Identität mit dem Attribut gelingt) → Feed-Container fährt die CDC-Runtime (ENV-Vertrag `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`/`CDC_SOURCE_ID`/`CDC_PUBLICATION`/`CDC_SLOT`/`CDC_TABLES` ist der Container-Vertrag in `compose.yaml`) → Toolchain-Container gegen das Compose-Netz; gepinnte Digests, DB-Daten im Container. Am Ende ein Black-Box-Rundlauf ([`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) E2E-Tier): `register-consumer`/`acknowledge-consumer` laufen ausschließlich per `docker exec` gegen den laufenden Feed-Container (kein Go-Paket-Import), über einen simulierten Container-Neustart hinweg — Fortsetzen ab der bestätigten Position, gelesen über den bestehenden SQL-Lesezugriffsweg `cdc.changes`. Zusätzlich ein SQL-Administration-Live-Reload-Beleg ([`ADR-0050`](../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)): `SELECT cdc.enable_table(...)`/`cdc.disable_table(...)` gegen eine nicht in `CDC_TABLES` gelistete Tabelle, Poll auf `status = 'applied'`, dann realer Erfassungs- bzw. Nicht-Erfassungs-Beleg am **laufenden** Feed-Container — ohne `docker restart` · seit slice-037. Zusätzlich ein CLI-Diagnose-Beleg ([`LH-FA-SST-003`](../spec/lastenheft.md), deckt `LH-FA-ADM-002`…`005`): derselbe externe `docker exec`-Zugriffsweg gegen den neuen `diagnose`-Sondermodus, einmal im Normalbetrieb (Betriebsstatus, `cdc_capture_lag`, Consumer-Rückstand) und einmal mit einem direkt in `cdc.process_heartbeat` geschriebenen Fehlerzustand (`LH-FA-ADM-003`-Boundary: erkennbar von Normalbetrieb unterscheidbar) — ohne den laufenden Feed-Container zu beenden · seit slice-038. Zusätzlich ein isolierter Publication-Entzug-Wirksamkeitsbeleg (`LH-FA-CFG-002`): `ALTER PUBLICATION ... DROP TABLE` direkt per `psql`, ohne den regulären `cdc.disable_table`-Antragsweg, damit die App-seitige `Assembler`-Filterung als Alternativerklärung ausgeschlossen ist — belegt real, dass PostgreSQLs bereits laufende Decoding-Session die entzogene Tabelle sofort ausfiltert · seit slice-051. Zusätzlich ein HTTP-API-Rundlauf ([`LH-FA-SST-006`](../spec/lastenheft.md), [`ADR-0057`](../docs/plan/adr/0057-http-grpc-api.md)): ein Wegwerf-Client (`tools/harness/httpclient`) ruft `RegisterConsumer` mit dem `admin`-Token und `ListTables` mit dem `reader`-Token real per HTTP gegen den laufenden Feed-Container auf (`CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` im Container-Vertrag) — ein echter Netzwerk-Request über den Compose-Netz-Alias, kein `docker exec` und kein Mock; die reale Registrierung wird zusätzlich gegen `cdc.consumer` bestätigt · seit slice-061. Zusätzlich liest derselbe Wegwerf-Client (`tools/harness/httpclient`) über `GET /changes` mit dem `reader`-Token Changes im Bereich `[from, to)` um eine unmittelbar zuvor erfasste Zeile — geprüft am Inhalt (Tabelle, Operation, Row Image) und über ihre `change_id` gegen den Lesezugriffsweg `cdc.changes` gehalten ([`LH-FA-SST-006`](../spec/lastenheft.md), [`ADR-0081`](../docs/plan/adr/0081-changes-lesen-ueber-die-http-api.md)) · seit slice-087. Zusätzlich ein gRPC-Stream-Rundlauf ([`LH-FA-SST-008`](../spec/lastenheft.md), [`ADR-0060`](../docs/plan/adr/0060-grpc-streaming-mechanismus.md)): ein Wegwerf-Client (`tools/harness/grpcclient`) öffnet real über gRPC den Server-Stream gegen den laufenden Feed-Container (`CDC_GRPC_ADDR` im Container-Vertrag) und empfängt eine danach committete Änderung — belegt am Stream-Image über Tabelle, Operation und den Wert der Spalte `name`, und über ihre `change_id` gegen den Lesezugriffsweg `cdc.changes` gehalten (die Feldvollständigkeit des Nachrichtenschemas trägt `internal/adapters/driving/grpc/server_test.go` auf Unit-Ebene) —, während ein Stream-Öffnungsversuch ohne gültiges Token mit gRPC-Status `Unauthenticated` abgelehnt wird — ein echter Netzwerk-Request über den Compose-Netz-Alias, kein `docker exec` und kein Mock · seit slice-071. Zusätzlich zwei Struktur-/Verhaltens-Belege im Go-Testpaket (`test/integration/integration_test.go`): `TestE2ESchemaChangeDropColumn` (`LH-FA-SCH-003`, [`ADR-0063`](../docs/plan/adr/0063-lh-fa-sch-003-testform-korrektur.md) Supersedes [`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md) Entscheidung 1) — eine eigene, wegwerfbare Spalte auf `feed_e2e_schema` real per `ALTER TABLE … ADD/DROP COLUMN` entfernt; die reale Entfernung löst denselben `relationOther`/`ErrIncompatibleSchemaChange`-Pfad aus wie eine inkompatible Typänderung (`LH-FA-SCH-004`, Konvergenz bereits von `ADR-0059` Teilfrage 4 akzeptiert) und beendet den Erfassungspfad des Feed-Containers dauerhaft sichtbar über `error_class=schema` — die zuvor erfasste Änderung bleibt über `cdc.changes` inklusive historischem Wert unverändert lesbar (Boundary); läuft deshalb als eigener Aufruf nach der ursprünglichen Container-Ende-Grenze, mit explizitem Replication-Slot-Neuanlage/Schema-Version-Nachtrag/Neustart/Health-Poll davor — und `TestE2EChangeTableMetadataExtensibility` (`LH-FA-DAT-006`, [`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md) Entscheidung 2) — ein reales `ALTER TABLE cdc.change ADD COLUMN` auf dem internen Change-Schema; vor und nach der Erweiterung erfasste Zeilen bleiben über `cdc.changes` lesbar, die neue Spalte ist über die View nicht sichtbar (explizite `columns:`-Signatur, `tools/schema/schema.yaml`), `t.Cleanup` nimmt sie real zurück · seit slice-062. Zusätzlich ein Upgrade-Sicherheits-Rundlauf (`LH-QA-OPS-005`, [`ADR-0064`](../docs/plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md) Supersedes [`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md) Entscheidung 3): ein realer Container-Tausch über `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` ersetzt den Feed-Container durch eine neue Instanz desselben `:dev`-Images (`container_name` bleibt `cdc-test-feed`) — `postgres`/`nats` bleiben durch `--no-deps` unberührt und healthy; Health-Poll wie beim bestehenden simulierten Neustart-Rundlauf, der vor dem Tausch erfasste Datenstand bleibt über `cdc.changes` identisch lesbar, eine danach eingefügte Zeile wird weiterhin erfasst · seit slice-063. Zusätzlich ein Spaltenausschluss-Rundlauf ([`LH-FA-CFG-005`](../spec/lastenheft.md), [`ADR-0059`](../docs/plan/adr/0059-spaltenauswahl-mechanismus.md); trägt zugleich `LH-QA-SEC-004`, dessen Messmethode ausschließlich über `LH-FA-CFG-005` läuft): eine eigene, nicht in `CDC_TABLES` gelistete Tabelle wird über `cdc.enable_table` aktiviert, `SELECT cdc.exclude_column(...)` gegen die laufende Aktivierung wird real verarbeitet (Poll auf `status = 'applied'`), die danach erfasste Change trägt den ausgeschlossenen Spaltenschlüssel nicht mehr im Row Image und den Wert nirgends, während die nicht ausgeschlossene Spalte darin bleibt und die vor dem Ausschluss erfasste Change unverändert lesbar bleibt; ein `cdc.exclude_column` gegen eine nicht existierende Spalte landet real `failed` samt Fehlertext — alles am **laufenden** Feed-Container, ohne Neustart · seit slice-068. Zusätzlich ein SSE-Stream-Rundlauf ([`LH-FA-SST-008`](../spec/lastenheft.md), [`ADR-0061`](../docs/plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)): ein Wegwerf-Client (`tools/harness/sseclient`) öffnet real per HTTP den Endpunkt `GET /changes/stream` gegen den laufenden Feed-Container (`CDC_HTTP_ADDR` im Container-Vertrag) und empfängt eine danach committete Änderung — belegt am Event über Tabelle, Operation und den Wert der Spalte `name`, und über ihre `change_id` gegen den Lesezugriffsweg `cdc.changes` gehalten (die Feldvollständigkeit des Nachrichtenschemas trägt `internal/adapters/driving/http/sse_test.go` auf Unit-Ebene) —, während ein Aufruf ohne gültiges Token mit HTTP-Status `401` abgelehnt wird; ein echter Netzwerk-Request über den Compose-Netz-Alias, kein `docker exec` und kein Mock · seit slice-072. Zusätzlich ein Spaltenausschluss-Neustart-Beleg ([`ADR-0065`](../docs/plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md), trägt `LH-QA-SEC-004`): nach einem **realen** `docker restart` des Feed-Containers leitet der Prozessstart den dauerhaften Ausschlussstand aus den `applied`-Zeilen der beiden Spalten-Antragsarten ab ([`LH-FA-CFG-005`](../spec/lastenheft.md) über dieselbe Tabelle wie im Spaltenausschluss-Rundlauf) — die danach erfasste Change trägt den ausgeschlossenen Spaltenschlüssel nicht mehr im Row Image und den Wert nirgends, während die nicht ausgeschlossene Spalte darin bleibt; der Beleg läuft vor der Container-Ende-Grenze und vor dem Upgrade-Tausch · seit slice-075. Zusätzlich ein NATS-Vollinhalts-Stream-Rundlauf ([`LH-FA-SST-008`](../spec/lastenheft.md), [`ADR-0100`](../docs/plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)): ein Wegwerf-Client (`tools/harness/natsstreamsub`) verbindet sich real per NATS mit dem gültigen `CDC_NATS_STREAM_TOKEN` gegen den laufenden Feed-Container (`CDC_NATS_URL` im Container-Vertrag, Server-Auth über `--auth`) und empfängt eine danach committete Änderung als vollständiges JSON-Event über `cdc.stream.<source_id>.<schema>.<table>` — belegt am Event über Tabelle, Operation und den Wert der Spalte `name`, und über ihre `change_id` gegen den Lesezugriffsweg `cdc.changes` gehalten (die Feldvollständigkeit des Nachrichtenschemas trägt `internal/adapters/driven/natsstream/publisher_test.go` auf Unit-Ebene) —, während ein Verbindungsversuch ohne Token und einer mit falschem Token vom NATS-Server abgelehnt werden; das bestehende Wecksignal (`tools/harness/natssub`) funktioniert mit demselben Test-Token unverändert weiter — ein echter NATS-Request über den Compose-Netz-Alias, kein `docker exec` und kein Mock · seit slice-nats-drittstream-core. Zusätzlich schreibt der Lauf `docs/user/e2e-abdeckung.md` — die E2E-Abdeckungstabelle je Spec-Kennung, abgeleitet aus dem Quelltext der `func TestE2E*` und den Phasen-Deklarationen des Runners (Erzeugnis, **kein** Lauf-Beleg) · seit slice-074 | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md), [`LH-QA-POR-003`](../spec/lastenheft.md) |
| `.github/workflows/e2e.yml` | automatisiert `make image` gefolgt von `make test-integration` auf jeden Pull Request und Push (derselbe Trigger wie `ci.yml`) — läuft aber als **eigener, nicht-blockierender** Workflow statt als Job in `ci.yml`: ein roter Lauf ist im PR sichtbar, blockiert aber keinen Merge (kein Required-Status-Check — das ist eine Repository-Einstellung, keine Workflow-Eigenschaft), damit der deutlich längere Compose-Stack-Aufbau `ci.yml`s kurzes Signal nicht verzögert. `make image` läuft zuerst, weil `compose.yaml` das lokal gebaute Image ohne eigenen `build:`-Block referenziert ([`ADR-0044`](../docs/plan/adr/0044-image-beleg-semantik.md)). Seit slice-064 trägt der Job eine `strategy: matrix:` über die beiden in `SPEC-012` festgelegten PostgreSQL-Digests (`LH-QA-POR-001`, [`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md) Entscheidung 4) — `fail-fast: false`, beide Legs laufen parallel und unabhängig sichtbar; jedes Leg setzt `PG_TEST_IMAGE` auf Job-Ebene, `compose.yaml`s `${PG_TEST_IMAGE}`-Interpolation liest die Variable, `make image`/`make test-integration` laufen unverändert je Leg | kein Gate, [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md), [`LH-QA-POR-001`](../spec/lastenheft.md), [`LH-QA-POR-003`](../spec/lastenheft.md) |
| `make schema-validate` | prüft das neutrale Schema-YAML mit d-migrate (netzlos; Vorlauf vor jedem `generate`/`migrate`) | kein Gate, [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md) |
| `make bench` | Performance-Benchmark-Infrastruktur (`ADR-0054` §(b)): drei eigenständige Bench-Skripte hinter einem Ziel, je eine eigene, von `tools/harness/run-integration-tests.sh` unabhängige PostgreSQL-/Feed-Umgebung (`docker network`/`docker run`, kein `compose.yaml`-Bezug) samt Schema-Rollout über d-migrate. `tools/bench-source-impact.sh` (`LH-QA-PER-001`): dieselbe Insert-Last einmal ohne jeden Replication-Slot, einmal mit aktivem Slot und laufendem Feed-Container, Laufzeit-Vergleich. `tools/bench-scaling.sh` (`LH-QA-PER-002`): durchläuft real die drei `SPEC-014`-Lastenstufen (klein/mittel/groß) gegen einen aktiven Feed-Container, dokumentiert erreichten Durchsatz und `cdc_capture_lag` je Stufe — Default-Modus fährt stark verkürzte Dauern bei unveränderter Ziel-Rate, `--full` fährt die tatsächlichen `SPEC-014`-Dauern. `tools/bench-batch-vs-single.sh` (`LH-QA-PER-003`): Batch- vs. Einzelabruf beim Lesen über `cdc.changes`. Alle drei dokumentieren Aufwand/Ergebnis auf stdout ohne Pass/Fail-Schwellenwert — anders als das Coverage-Gate | kein Gate, [ADR-0054](../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b) |
| `make schema-rollout` | rollt das Schema mit d-migrate aus — `schema migrate --execute` mit Pflicht-Report (`tools/schema/plan.yaml`) und Rollback-Artefakt (`tools/schema/down.sql`); braucht DB-Zugang. Die drei Views (`active_tables`, `consumer_status`, `changes`) sind seit slice-016 deklarativ im `views:`-Knoten von `tools/schema/schema.yaml` überführt (`source_dialect: postgresql`, sichtbare `columns:`-Signatur) — d-migrate 1.3.1 behebt den Post-Compare-Drift auf frisch angelegten Sichten (fehlendes `ViewDefinition.sourceDialect` im Fingerabdruck-Vergleich); real gegen eine leere DB (Erstanlage) UND einen Folgelauf gegen eine bereits migrierte DB (`ReplaceView`) getestet, beide Exit 0, kein `VIEW_SIGNATURE_UNKNOWN`-Blocker. Die Ausweichform `nacharbeit-views.sql` ist zurückgebaut · seit slice-016. Ein zweiter Lauf gegen ein bereits migriertes Ziel ist seit `schema-rollout-zentrale-idempotenz-wache` insgesamt idempotent: ein vorgelagerter `--plan-only`-Lauf plus `tools/schema/rolloutguard` erkennt, wenn ausschließlich die sechs bekannten, außerhalb des neutralen Modells liegenden Fremdobjekte blockieren (vier SQL-Funktionen, zwei Views aus den `nacharbeit-*.sql`-Dateien), und lässt in diesem Fall den nachfolgenden `--execute`-Schritt zusätzlich mit `--allow-destructive` laufen statt mit Exit 8 abzubrechen; `--execute` selbst läuft dabei immer, damit eine echte, gleichzeitig anstehende, nicht-destruktive Änderung im selben Plan nicht verlustig geht (ein reines Überspringen von `--execute` verwarf genau das real, siehe Slice-Plan §1) — jeder andere Blocker bricht weiterhin real mit Exit 8 ab, ohne `--allow-destructive` — real geprüft mit `tools/harness/run-schema-rollout-guard-test.sh` (vier Läufe: frisch, `--allow-destructive`-Pfad, echte Änderung bleibt wirksam, Negativ-Abbruch) · seit schema-rollout-zentrale-idempotenz-wache | kein Gate, [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md) |
| `make examples-csharp` | baut die vier Werkzeugketten-Images der C#-Sprach-Wurzel (Bau-Kontext `examples/csharp/`, eigenes digest-gepinntes Dockerfile, gemeinsame `build`-Stufe) — `pg-change-feed-examples:csharp` (http-client), `pg-change-feed-examples:csharp-sse` (sse-client), `pg-change-feed-examples:csharp-nats` (nats-client) und `pg-change-feed-examples:csharp-grpc` (grpc-client, seit slice-102) — und fährt darin `dotnet restore`/`build`/`test` für alle vier Programme; ein roter Test bricht den `docker build` ab; die Beispiele bleiben Doku mit Bau-Bindung, kein Lauf-Beleg. Für grpc-client liest der Bau die `.proto` zusätzlich über einen benannten Bau-Kontext (`--build-context proto=proto`, `COPY --from=proto …`) — jeder der vier Aufrufe trägt ihn zwingend, ohne ihn bricht der Bau an dieser Stelle ab. Braucht Netz (NuGet-Paketbezug), deshalb Werkzeug statt Gate | kein Gate, [`ADR-0087`](../docs/plan/adr/0087-beispiel-clients-csharp-kotlin.md), [`ADR-0090`](../docs/plan/adr/0090-beispiel-clients-volle-matrix.md) · seit slice-098, erweitert seit slice-100, seit slice-101, seit slice-102 |
| `make examples-kotlin` | baut die vier Werkzeugketten-Images der Kotlin-Sprach-Wurzel (Bau-Kontext `examples/kotlin/`, eigenes digest-gepinntes Dockerfile, gemeinsame `build`-Stufe) — `pg-change-feed-examples:kotlin` (http-client), `pg-change-feed-examples:kotlin-sse` (sse-client), `pg-change-feed-examples:kotlin-nats` (nats-client) und `pg-change-feed-examples:kotlin-grpc` (grpc-client, seit slice-103) — und fährt darin den Gradle-Wrapper (`test`/`installDist`) für alle vier Module; ein roter Test bricht den `docker build` ab; die Beispiele bleiben Doku mit Bau-Bindung, kein Lauf-Beleg. Für grpc-client liest der Bau die `.proto` zusätzlich über einen benannten Bau-Kontext (`--build-context proto=proto`, `COPY --from=proto …`) — jeder der vier Aufrufe trägt ihn zwingend, ohne ihn bricht der Bau an dieser Stelle ab. Braucht Netz (Maven-Central-/Gradle-Plugin-Portal-Paketbezug), deshalb Werkzeug statt Gate | kein Gate, [`ADR-0087`](../docs/plan/adr/0087-beispiel-clients-csharp-kotlin.md), [`ADR-0090`](../docs/plan/adr/0090-beispiel-clients-volle-matrix.md) · seit slice-099, erweitert seit slice-100, seit slice-101, seit slice-103 |
| `make example-run-go` | baut bei Bedarf (Docker-Layer-Cache) und startet real eines der vier Go-Beispiele — Pflicht-Argument `SURFACE=http\|sse\|grpc\|nats`, optional `ARGS="…"` für die Flag-Übersteuerung je Beispiel; ein fehlendes/unbekanntes `SURFACE` bricht mit `$(error …)` ab, bevor `docker build`/`docker run` läuft. Bau-Quelle ist `examples/Dockerfile` (Bau-Kontext Repo-Wurzel, isoliert über `examples/Dockerfile.dockerignore`, digest-gepinnte Basis wiederverwendet aus der Wurzel-`Dockerfile`); der Start ist `docker run --rm --network cdc-examples --env-file examples/.env <image-tag> $(ARGS)` — das Docker-Netzwerk `cdc-examples` und `examples/.env` legt `slice-beispiele-compose-bootstrap` an, dieses Ziel referenziert beide nur. Analog dazu, aber ohne eigenen Bau-Schritt: `example-run-csharp`/`example-run-kotlin` (siehe Zeile darunter). Braucht das benannte Docker-Netzwerk, deshalb Werkzeug statt Gate | kein Gate, [`ADR-0098`](../docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) |
| `make example-run-csharp`/`make example-run-kotlin` | startet real eines der vier bereits gebauten C#- bzw. Kotlin-Beispiele — Pflicht-Argument `SURFACE=http\|sse\|grpc\|nats`, optional `ARGS="…"` für die Flag-Übersteuerung je Beispiel; ein fehlendes/unbekanntes `SURFACE` bricht mit `$(error …)` ab, bevor `docker run` läuft. Anders als `example-run-go` bauen diese Ziele **nichts** — der Bau liegt bei `make examples-csharp`/`make examples-kotlin` (`ADR-0087`/`ADR-0090`, unverändert); der Start ist `docker run --rm --network cdc-examples --env-file examples/.env pg-change-feed-examples:csharp[-<surface>]`/`:kotlin[-<surface>] $(ARGS)` gegen den bereits vorhandenen Image-Tag — läuft der Bau-Schritt noch nicht, scheitert `docker run` real und sichtbar mit einem Docker-eigenen „image not found", kein stiller Vorab-Bau. Das Docker-Netzwerk `cdc-examples` und `examples/.env` legt `slice-beispiele-compose-bootstrap` an, dieses Ziel referenziert beide nur. Braucht das benannte Docker-Netzwerk, deshalb Werkzeug statt Gate | kein Gate, [`ADR-0098`](../docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) · seit slice-beispiele-csharp-kotlin-start-target |
| `make example-demo-up`/`make example-demo-down` | fährt die Demo-/Quickstart-Umgebung unter `examples/` real hoch bzw. räumt sie ab (`examples/compose.yaml`, eigenes Docker-Netzwerk `cdc-examples`, eigenständig von der Wurzel-`compose.yaml` — beide können gleichzeitig laufen, real geprüft): PostgreSQL+NATS starten, Schema-Rollout über d-migrate treiben (`tools/schema/apply-rollout.sh`-Muster, `examples/bootstrap.sh`), eine Beispiel-Quelle (`cdc.source`) und -Tabelle (`public.orders`, eine Demo-Zeile) registrieren, dann den Feed-Container starten — der aktiviert `public.orders` beim eigenen Start automatisch über `CDC_TABLES` (`ADR-0028`). Die Demo-Zeile wird erst **nach** Anlage des Replication-Slots eingefügt (logische Replikation trägt keinen initialen Snapshot-Export dieses Wegs — eine frühere Zeile bliebe über `GET /changes` unsichtbar, real geprüft); danach liefert `GET /changes?source=demo-source` real die Zeile, ohne dass ein Nutzer manuell SQL oder `make schema-rollout` aufruft. Idempotent: ein zweiter `example-demo-up`-Lauf überspringt den Rollout-Schritt gegen ein bereits migriertes Ziel über den eigenen, aufruferseitigen Existenz-Check von `examples/bootstrap.sh`; `make schema-rollout` selbst ist seit `schema-rollout-zentrale-idempotenz-wache` ebenfalls idempotent gegen die sechs bekannten Fremdobjekte (siehe `make schema-rollout`-Zeile oben) — der Aufrufer-Check bleibt trotzdem bestehen. Fügt die Demo-Zeile nur ein, wenn `public.orders` noch leer ist. `examples/.env` (committet, keine Vorlage, `ADR-0098` Festlegung 3) ist die gemeinsame Umgebungsdatei-Quelle für den Feed-Container (`env_file:`) und jedes `make example-run-*`-Ziel (`--env-file`). Braucht DB-Zugang und das benannte Docker-Netzwerk, deshalb Werkzeug statt Gate | kein Gate, [`ADR-0098`](../docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) · seit slice-beispiele-compose-bootstrap |
| `make <mover>` | bewegt <…>, prüft nichts | kein Gate |
| `make <messung>` | misst <…> gegen <Schwelle> | kein Gate, ADR-<NNNN> |
| `make <vorschau>` | sagt, was <schreibender Lauf> täte; Ausgänge und Sperren in der verlinkten Datei | kein Gate |

**Aktueller Lauf-Status:** [![ci](https://github.com/pt9912/pg-change-feed/actions/workflows/ci.yml/badge.svg)](https://github.com/pt9912/pg-change-feed/actions/workflows/ci.yml)
bzw. lokal `make help` / `make gates`. Der Workflow
[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) automatisiert
diesen Gate-Lauf (plus `make test`) auf jeden Pull Request und Push — kein
neues Gate, nur die Automatisierung des bestehenden
([`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md)).
Unmittelbar nach dem Checkout, vor dem `Gates`-Schritt, prüft ein eigener
Workflow-Schritt `uname -s`/`go env GOOS` gegen `Linux`/`linux` mit
sichtbarem Fehlschlag bei Abweichung — der sichtbare Log-Beleg für
[`LH-QA-POR-002`](../spec/lastenheft.md)
([`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md)
Entscheidung 5, kein Gate).
**Rote Gates:** Begründung im verlinkten `CO-<NNN>` (siehe Bindung-Spalte), Modul 7.
**Nicht behauptet** (geplant): `make image-cve` (CVE-Scan des gebauten
Images, advisory) — die Aktivierungsbedingung ist seit slice-001 eingetreten
(der erste `make image`-Lauf läuft grün); das Target ist noch nicht
implementiert, bis dahin bleibt der Scan unausgesprochen.

<!-- Domänenspezifische Gates ergänzen, je nach Repo-Klasse: -->

## Traceability rules

- PRs/Commits **müssen** mindestens eine `<LH-*>` oder `ADR-*`-ID nennen.
- Mechanisch getragen (Standing-Gate, [`ADR-0045`](../docs/plan/adr/0045-commit-traceability-standing-gate.md)): `make commit-traceability` prüft je Message der letzten 5 Commits (`RANGE=base..head` überschreibt) beide Grenzen — positive Hälfte via d-check Modul `commits` (Befund `commit-untraceable`), Betreff-Grenze via `tools/harness/commit-traceability.sh` · seit slice-006 (ausgelöst durch das dritte Auftreten der Verstoß-Klasse, review-slice-006 F-2).
- Optional, **nicht-durchsetzend** ([`ADR-0062`](../docs/plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md), Zusage einseitig nach [`ADR-0069`](../docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md)): der lokale `commit-msg`-Hook unter `.githooks/commit-msg` meldet vorab in den Klassen, in denen er liest — er fängt **nicht** jeden Verstoß, den `make commit-traceability` fängt, und weist keinen Commit zurück, den das Gate zulässt. Aktivierung ist ein einmaliger, lokaler Schritt — `git config core.hooksPath .githooks` (`core.hooksPath` ist lokale Konfiguration, sie reist nicht mit dem Klon; `git commit --no-verify` umgeht den Hook). Er ersetzt `make commit-traceability` nicht: er läuft ausschließlich lokal vor `git commit` und ist kein Gate-Ziel · seit slice-073.
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Neue ADRs müssen im ADR-Index ergänzt werden.
- Änderungen an Planning-Dokumenten müssen die Lifecycle-Regeln beachten (open → next → in-progress → done; reine `git mv`-Commits siehe AGENTS.md §3.3).

## Safety and scope boundaries

<!--
Repo-spezifisch formulieren. Beispiele:

Für ein Referenz-Repo:
- Dies ist kein produktiver Service.
- Externer Cloud-Zugriff darf nicht für lokale Demo-Abnahme vorausgesetzt werden.
- Determinismus und Replayability sind Kernverträge.

Für ein Safety/Control-Repo:
- Markt-/Optimierungs-Output muss durch Statemachine, Constraint-Limiter, Ramp-Limiter fließen.
- Software-Stop ersetzt keine Hardware-Sicherheitsfunktionen.
- Produktion-Profile müssen fail-closed sein.

Für ein Policy/Compliance-Repo:
- Dieses Werkzeug ist keine Rechts-/Steuer-/Fachberatung.
- KI-Funktionen liefern Vorschläge, keine verbindlichen Entscheidungen.
-->

- <…>
- <…>

## Minimal agent workflow

1. Diese Datei lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`) — Exit-Code direkt prüfen, nie durch eine Pipe/einen Wrapper hindurch (`AGENTS.md` §3.9).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe `harness/README.md` §Guides) →
Verifier. Kein Self-Review — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).

## Leseordnung

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt — die Menschen-Hälfte des Einstiegs:
drei bis fünf **geordnete** Zeiger, was ein neuer Mensch zuerst liest und was
bei Bedarf; eine Leseordnung, die alles nennt, ist keine.

1. <zuerst — z. B. `AGENTS.md` §Hard Rules>
2. <dann — z. B. `spec/lastenheft.md`>
3. <bei Bedarf — z. B. `harness/conventions.md`>
