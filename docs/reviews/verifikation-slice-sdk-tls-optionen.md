# Verifikations-Report: slice-sdk-tls-optionen — 2026-10-05

**Review-Art:** Verifikation (Modul 11) — DoD- und Entscheidungs-Konformität plus Plan-vs-Code-Diff; Frage „Bauen wir es richtig?“, nicht Review des Diffs.

**Gegenstand:** `slice-sdk-tls-optionen`, Commits a44f432d (C#), ecff358f (Kotlin), 6bd41f37 (Python), d0e781b2 (Harness-Runner), c1fcd9a8 (Plan-Nachzug, Handbuch 1.100), Fixrunde fb4e8d70 / 4e4a2a7e; Stand der Messungen: Arbeitsbaum auf ec0d043d (clean).

**Rolle:** `verifier` · **Modell:** Sonnet 5.5 · **Datum:** 2026-10-05

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad
> in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-sdk-tls-optionen` (DoD §2, Belege §3, Risiken §6)
- `LH-FA-SST-013` (Lastenheft 0.16.0), `SPEC-037` (Pflichtenheft), `ADR-0150`, `ADR-0110`, `ADR-0145`
- Reviews `review-slice-sdk-tls-optionen` (F-1 bis F-7) und `review-slice-sdk-tls-optionen-fixrunde` (N-1 bis N-4)
- `AGENTS.md` §3.1, §3.9, §3.12, §3.13

---

## 1. Eigene Messungen (gefahren in diesem Lauf, Ausgabe gelesen)

Alle Exit-Codes direkt aus dem `make`-Aufruf gelesen, nicht durch eine Pipe. Im Scratchpad lagen Ausgangsdateien früherer Läufe; der `make gates`-Lauf, auf den dieser Report sich stützt, schrieb in ein frisches Verzeichnis.

| Messung | Ergebnis (gedruckte Zeile) | Stand |
|---|---|---|
| `make gates` | Exit 0; `baseline-verify: v6.14.0 OK — 54 Dateien`; `d-check: 1695 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s)`; `coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`; `generated-sync: OK`; `sdk-public-doc-check: keine interne Kennung unter sdks`; `handbuch-public-doc-check`, `ausgabe-kennungen-check`, `meldungscodes-check` grün; `a-check: gesamt: 0 Befund(e)` | **gemessen** |
| `make suchlauf-nachmessen PLAN=…` | `suchlauf-nachmessen: 6 Zeilen stimmen` (Soll = Ist in allen sechs; `diff 38` und `diff 0` am Arbeitsbaum) | **gemessen** |
| `make sdk-public-doc-check`, `make handbuch-public-doc-check` | Exit 0 je | **gemessen** |
| `make doc-commits RANGE=267a58b6..HEAD` | Exit 0, `0 Befund(e)` | **gemessen** |
| `make doc-immutable RANGE=267a58b6..HEAD` | Exit 0, `0 Befund(e)`; ohne `RANGE` endet das Ziel mit Exit 2 (`flag needs an argument: --range`) — Aufrufform, kein Befund am Slice | **gemessen** |
| C# Unit (Docker-Bau der Stufe `build`, `--no-cache`) | `Passed!  - Failed:     0, Passed:   254, Skipped:     0, Total:   254` | **gemessen**, gleich der Zusage (254) |
| Python Unit (`--no-cache`) | `252 passed, 6 warnings` | **gemessen**, gleich der Zusage (252) |
| Kotlin Unit (`--no-cache`) | `BUILD SUCCESSFUL`; Ergebnisdatei `TlsClientTest`: `tests="13" … failures="0" errors="0"`; Summe aller Testdateien 135 (grobe Summe über die XML-Köpfe, nicht als Zusage gelesen) | **gemessen**, `TlsClientTest` gleich der Zusage (13) |
| `make sdk-pack-csharp`, `-python`, `-kotlin` | je Exit 0; Artefakte `PgChangeFeed.Client.0.6.1.nupkg`, `pgchangefeed-kotlin-0.6.1.jar` und `…-sources.jar`, `pgchangefeed-0.6.1-py3-none-any.whl` und `.tar.gz` | **gemessen** |
| `make image`, dann `make test-sdk-python-integration` | Exit 0; Zeile: `run-sdk-python-integration-tests: TLS-Belege (LH-FA-SST-013, ADR-0150) grün — … HTTP (consumer_id=python-sdk-tls-fcac5647d8ea), SSE (change_id=945-1), gRPC-Stream (change_id=948-1) und gRPC-Verwaltung (consumer_id=python-sdk-tls-admin-41ec64c4ff4f) arbeiteten über TLS mit trust_anchor_file …; ohne Anker, mit fremdem Anker und mit einem Servernamen außerhalb des Zertifikats scheiterte dieselbe Verbindung an der TLS-Prüfung` | **gemessen** |
| `make test-sdk-csharp-integration` | Exit 0; gleiche Zeile mit `TrustAnchorFile` (consumer_id=csharp-sdk-tls-20261005071143, change_id 954-1 und 957-1) | **gemessen** |
| `make test-sdk-kotlin-integration` | Exit 0; gleiche Zeile mit `trustAnchorFile` (consumer_id=kotlin-sdk-tls-20261005071616, change_id 980-1 und 985-1) | **gemessen** |
| Nach den Läufen `git status` | clean (der Runner schreibt `docs/user/sdk-e2e-abdeckung.md` idempotent, kein Diff) | **gemessen** |
| Kein Zertifikat/Schlüssel im Repo | `git ls-files` auf `*.pem/*.key/*.crt/*.cer/*.p12/*.pfx/*.jks`: 0 Treffer; `git grep 'BEGIN CERTIFICATE|PRIVATE KEY'` außerhalb von `.md`: 0 Treffer | **gemessen** |
| Versionen | `.csproj` `0.6.1`, `build.gradle.kts` `0.6.1`, `pyproject.toml` `0.6.1`; der Diff 267a58b6..HEAD an diesen drei Dateien ändert nur das Test-Extra `cryptography>=43` | **gemessen** |
| Server-Code unverändert | `git diff 267a58b6..HEAD --stat` über `internal`, `cmd`, `gen`, `proto`, `spec`, `compose.yaml`, `Makefile`: leer | **gemessen** |

**Nicht wiederholt, aus dem Plan und den Reviews übernommen:** alle Mutationsproben (C# 3 + 3, Kotlin 4, Python 4, Realserver-Gegenläufe mit vertauschten Ankern), der rote Lauf von `make test-sdk-kompat` und der Registry-Modus, der Kotlin-Gegenlauf `SDK_TLS_CA_FILE=/tls/other.pem`. Ihre Zahlen stehen im Plan als gemessen/übernommen/hergeleitet gekennzeichnet; ich bestätige hier nur, dass die Tests, an die sie binden, in dieser Zahl und grün existieren.

## 2. DoD-Abgleich (Plan §2)

| DoD-Punkt | Befund | Beleg |
|---|---|---|
| C#: Option, vier Flächen, Happy/Boundary/Negative | erfüllt | `TlsClientTests` (13 Methoden × 4 Flächen: Happy, Aussteller, Blatt eines Ausstellers als alleiniger Anker, fremder Anker, mehrere Zertifikate, Namensfehler, abgelaufen, ohne Anker, `http://` unverändert, Anker plus `http://`), 254/0 gemessen |
| Kotlin und Python: dieselben Fälle, derselbe Inhalt | erfüllt | Fallnamen in Kotlin (13) und Python (siehe `test_tls.py`) spiegeln die C#-Fälle; Python 252, Kotlin 13 gemessen |
| Realserver-TLS-Phase je SDK, Abdeckungs-Träger | erfüllt | alle drei Runner Exit 0 mit gedruckter Abschlusszeile (§1); `docs/user/sdk-e2e-abdeckung.md` trägt je Sprache die Zeile mit `LH-FA-SST-013`; Server und `compose.yaml` unverändert |
| `make gates` grün inkl. `sdk-public-doc-check` | erfüllt | §1 |
| Review durchgeführt, Report liegt vor, kein Self-Review | erfüllt | zwei Reports vorhanden; Fixrunde ohne offenes HIGH/MEDIUM (nur der Prozessfund F-1, siehe §4) |
| Doku-Update (READMEs, Handbuch 1.100) | erfüllt | Handbuch `Version: 1.100`, Historienzeile 1.100, Abschnitt „Schnittstellen mit TLS verschlüsseln“ mit Absatz „Client-Pakete“; drei SDK-READMEs mit Abschnitt „Connect over TLS“; Suchlauf stimmt |
| Closure-Notiz, Register, Risiken-Ausgang, drei Paarungen | **offen** | Plan §7 steht auf „—“, §6 trägt bei den Risiken noch „weiter offen“; das sind Closure-Pflichten des Planners, kein Mangel der Lieferung |

Liefer-Punkte (drei: je SDK Option, Tests, Realserver-Phase): alle drei **erfüllt**.

## 3. Entscheidungs-Konformität und Aussagen-Abgleich (§3.12)

| Aussage (Träger) | Abgleich gegen Code/Messung | Stand |
|---|---|---|
| „Ein Vertrauensanker = Pfad einer PEM-Datei, ein oder mehrere Zertifikate, Server oder Aussteller“ (`SPEC-037`, Handbuch, READMEs) | Unit-Fälle `IssuerAsTrustAnchor`, `IssuedServerCertificateAsOnlyTrustAnchor`, `AnchorFileWithSeveralCertificates` in allen drei Sprachen grün (gemessen); am Realserver nur das selbstsignierte Server-Zertifikat (Runner-Zeile), der Ausstellerfall ist dort nicht gefahren, weil `certgen` selbstsigniert (Plan, übernommen) | gemessen (Unit), Realserver für Aussteller nicht belegt |
| „Mit Anker vertraut die Verbindung genau diesen Zertifikaten; die Anker des Betriebssystems gelten nicht zusätzlich“ | Python HTTP und SSE an der Eingabeseite gebunden (Plan: Mutation rot, übernommen); C#, Kotlin, Python gRPC nicht gebunden, Plan nennt Grund (prozessweiter Wurzelspeicher) und kennzeichnet die Aussage dort als *hergeleitet* | teils gebunden, Rest *hergeleitet*, benannt |
| „Kette, Gültigkeitszeitraum, Servername werden immer geprüft; keine Abschaltung, kein Überschreiben“ | Fälle `ForeignAnchor`, `ServerNameNotInCertificate`, `ExpiredCertificate` in allen drei Sprachen grün; Realserver-Zeile nennt „ohne Anker, mit fremdem Anker und mit einem Servernamen außerhalb des Zertifikats scheiterte dieselbe Verbindung“ in allen drei Läufen; öffentliche Option-Fläche trägt keinen Schalter (`TlsTransport.cs`, `tls.py`, `TlsSupport.kt` gelesen: nur Anker-Pfad) | gemessen |
| „Anker plus `http://` wirft“ (READMEs) | Fall `TrustAnchorWithPlaintextAddress…` / `a trust anchor with a plaintext address is refused…` / `test_trust_anchor_with_plaintext_address_is_refused` grün in allen drei SDKs. **Grenze:** `SPEC-037` und `ADR-0150` sagen dazu nichts; es ist eine Festlegung der Umsetzung (C#: `http`-Adresse; Kotlin: jedes Schema außer `https`; Python gRPC: `host:port` ohne Schema geht mit Anker über TLS). In den READMEs so dokumentiert. | gemessen, **über die Spec hinaus** (V-2) |
| „Eigene Clients/Kanäle erhalten den Anker nicht“ (READMEs, Handbuch; `SPEC-037` Abgrenzung: übergebene Konfiguration bleibt zulässig) | aus Code und Konstruktorform **hergeleitet** (Konstruktoren `…(httpClient, options)`, `CallInvoker`/`Channel`, Python liest nur `create_http_client`/`create_grpc_channel`); **kein Test** bindet, dass ein übergebener Client unverändert bleibt (V-1) | *hergeleitet*, nicht gemessen |
| „Verbindungsfehler der Fläche, ohne Meldungscode des Servers“ (`SPEC-037`) | C#-`MapException` bildet `Internal` mit `HttpRequestException` auf `PgChangeFeedGrpcUnexpectedStatusException` ab; Unit-Fälle `ExpectConnectionFailureAsync` grün; die Verhaltensänderung steht als Upgrading-Eintrag im README | gemessen (Unit), Kompat-Ziel misst sie nicht (N-3 des Fixrunden-Reviews) |
| „NATS unberührt“ | `git diff` an NATS-Flächen nicht im Slice-Umfang; READMEs nennen die Grenze | geprüft |
| „Versionen bleiben 0.6.1“ | gemessen (§1); README-Überschrift „0.7.0“ in C# und Kotlin ist der Inhalt der nächsten Minor-Version, das gebaute Paket trägt `0.6.1` | gemessen; Spannung benannt (V-4) |

**Plan-vs-Code-Diff:** die Plan-Tabelle (Datei/Komponente) und der Plan-Nachzug sind im Baum wiedergefunden: `TlsTransport.cs`, `TlsSupport.kt`, `tls.py`/`options.py`, Tests je SDK, `lib-sdk-tls-fixture.sh`, die drei Runner, `harness/mk/sdk.mk`, Abdeckungs-Träger, Handbuch, `harness/README.md` (drei Sensors-Zeilen). Der Diff trägt keine Datei außerhalb dieses Umfangs: kein Server-Code, kein `compose.yaml`, keine Versionsänderung.

## 4. Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | LOW | Die Zusage „ein selbst übergebener `HttpClient`/Kanal behält seine eigene TLS-Einstellung, der Anker wirkt dort nicht“ steht in den drei READMEs und im Handbuch, hat aber in keinem SDK einen Test; sie ist aus der Konstruktorform *hergeleitet*. Failure-Szenario: eine künftige Änderung wendet den Anker auch auf einen übergebenen Client an, kein Test wird rot. | `SPEC-037` Abgrenzung; `AGENTS.md` §3.12 | `sdks/csharp/README.md` · `are not touched by the trust anchor option` | ja — ein Fall je SDK mit übergebenem Client und fremdem Anker, der die Verbindung nicht ändert | Zusage ohne Test, *hergeleitet* |
| V-2 | INFO | „Anker plus nicht-TLS-Adresse wirft“ ist in keiner Norm festgelegt (weder `SPEC-037` noch `ADR-0150`), die drei SDKs setzen es gleich um und dokumentieren es, mit sprachlichen Unterschieden (Kotlin: jedes Nicht-`https`; Python gRPC: schemaloses `host:port` mit Anker wird TLS). Konsistent, aber eine Festlegung ohne Norm-Anker. | `SPEC-037` | `sdks/kotlin/pgchangefeed-kotlin/README.md` · `throws an \`IllegalArgumentException\` when the client is created` | nein — Normfrage | Umsetzungsentscheidung ohne Norm |
| V-3 | INFO | Der Realserver belegt nur das selbstsignierte Server-Zertifikat als Anker; die Zusage „oder des Ausstellers“ ist in den Unit-Fällen aller drei SDKs gemessen, am Realserver nicht (kein CA-ausgestelltes Zertifikat in `certgen`). Die Handbuch-Messzeile behauptet dazu nichts. | `LH-FA-SST-013` | `tools/harness/lib-sdk-tls-fixture.sh` | ja — `certgen`-Erweiterung | Ausstellerfall nur auf Unit-Ebene |

Offene Punkte aus den Reviews, hier nur **gemessen oder bestätigt**, nicht neu bewertet:

- **N-1 (LOW):** Systemanker-Ausschluss außerhalb Python HTTP/SSE ungebunden; Plan benennt Grund und kennzeichnet *hergeleitet*. Bestätigt.
- **N-2 (LOW):** README-Überschrift „0.7.0“ (C# und Kotlin) gegenüber Package `0.6.1`: **gemessen** (§1). Auflösung ist der Release-Zug (Versionen anheben vor jeder Veröffentlichung).
- **N-3 (LOW):** `make test-sdk-kompat` Werkzeug-Drift: `tools/harness/sdk-kompat/csharp/Dockerfile` ist auf Package `0.6.0` festgelegt (gemessen: Kopf- und Stufenkommentare nennen 0.6.0), `sdks/csharp/dist` trägt `0.6.1`. Das Ziel selbst **nicht gefahren** (braucht Netz); der Rotbefund ist aus dem Plan **übernommen**. Es fehlt weiter eine Adresse mit Frist: das ist eine Planner-Übergabe.
- **N-4 (INFO):** Untergrenze `NotBefore` im Gleichheitszweig von `TlsTransport.cs` ungebunden. Übernommen.
- **Prozessfund F-1 des Erst-Reviews** (`cat >>` auf eine Repo-Datei, §3.1): **übernommen**, nicht messbar; im Beobachtungs-Register unter `BEO-PGC` finde ich keinen Eintrag zur Umleitung auf Repo-Dateien. Eintrag und Zähler liegen beim Planner.

## 5. Übergaben (offen)

| Was | An | Adresse / Frist |
|---|---|---|
| Prozessfund `cat >>` (Host-Umleitung auf Repo-Datei, Guard liest sie nicht) | Planner | Beobachtungs-Register `BEO-PGC`, bei der Closure von `slice-sdk-tls-optionen` |
| `make test-sdk-kompat`: Werkzeug an Package-Version gekoppelt (`0.6.0` fest, dist `0.6.1`; bricht beim Versionssprung erneut) | Planner | Frist: Closure dieses Slice oder Folge-Slice vor dem Release-Zug; benannte Adresse fehlt noch |
| README-Überschrift „0.7.0“ gegen Package `0.6.1` | Release-Zug | Versionen anheben vor jeder Veröffentlichung (braucht neue Freigabe des Auftraggebers) |
| Closure-Notiz, Register, Risiken-Ausgang (§6: drei Risiken stehen auf „weiter offen“), drei Paarungen | Planner | Closure |
| V-1 (Test für übergebene Clients), V-3 (Ausstellerfall am Realserver) | Planner | Folge-Slice oder Verzicht mit Begründung |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Gates (`make gates`, zehn Ziele) | gemessen, grün |
| Unit-Zahlen (C# 254, Python 252, Kotlin `TlsClientTest` 13) | gemessen, gleich der Zusage |
| Realserver-TLS-Phasen aller drei SDKs | gemessen, Exit 0 und Abschlusszeile gedruckt |
| Schlüsselmaterial im Repo | gemessen, keines |
| Package-Versionen | gemessen, unverändert `0.6.1` |
| Server und `compose.yaml` | gemessen, unverändert |
| Suchlauf (§3.13) | gemessen, sechs Zeilen stimmen |
| Traceability und Immutabilität (`267a58b6..HEAD`) | gemessen, 0 Befunde |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

(Dazu die übernommenen Review-Befunde N-1 bis N-4 und der Prozessfund, hier nicht neu gezählt.)

**Finding-Klassen dieses Laufs:** Zusage ohne Test, *hergeleitet* · Umsetzungsentscheidung ohne Norm · Ausstellerfall nur auf Unit-Ebene

## Verdikt

**DoD-Lieferung: bestanden.** Alle drei Liefer-Punkte und alle Punkte der DoD bis einschließlich Doku-Update sind durch eigene Messung belegt; die vier Closure-Punkte sind regulär offen und liegen beim Planner. Kein HIGH, kein MEDIUM, keine Merge-Blockade durch die Verifikation. Die Aussagen „Server oder Aussteller als Anker“ (Unit gemessen, Realserver nur Server), „Anker plus `http://` wirft“ (gemessen, aber über die Spec hinaus) und „eigene Clients/Kanäle erhalten den Anker nicht“ (*hergeleitet*, ungetestet) sind oben mit ihrem Stand gekennzeichnet. Der Slice kann nach Closure-Notiz nach `done/` gehen; die Übergaben in §5 sind vorher mit Adresse zu versehen.
