# Review-Report: slice-sdk-tls-optionen — 2026-10-05

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules

**Gegenstand:** Diff 267a58b6..c1fcd9a8, beschränkt auf die fünf Commits a44f432d (C#), ecff358f (Kotlin), 6bd41f37 (Python), d0e781b2 (Harness), c1fcd9a8 (Plan-Nachzug); der fremde Commit ad151ab7 gehört nicht zum Gegenstand

**Skill:** `.harness/skills/reviewer.md` @ c1fcd9a8 (Stand Baseline v6.14.0)
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-05

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-sdk-tls-optionen`
- [`LH-FA-SST-013`](../../spec/lastenheft.md) (Lastenheft 0.16.0), [`SPEC-037`](../../spec/pflichtenheft.md)
- [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md), [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.12, §3.13

**Eigene Proben dieses Laufs** (alle an Kopien im Scratchpad, nichts am Repo):
`make suchlauf-nachmessen PLAN=<Plan>` (6 Zeilen stimmen), `make sdk-public-doc-check`
und `make handbuch-public-doc-check` (je Exit 0), `make kommentar-kennungen DIFF=267a58b6 PATHS=sdks`
(kein Kandidat; Form, nicht Wahrheit), ein C#-Chain-Probeprogramm im gepinnten
`dotnet/sdk`-Image (`--network none`, siehe F-2). Die Unit-Tests, die Mutationsläufe und die
Realserver-Läufe des Implementers sind **nicht** wiederholt, ihre Zahlen sind **übernommen**.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Implementer meldet, eine Zeile an `run-sdk-csharp-integration-tests.sh` per `cat >>` angehängt zu haben (Host-Umleitung auf eine Repo-Datei). Die Spur im Diff ist nur die behobene fehlende End-Newline der letzten `echo`-Zeile; der Wortlaut der Umleitung steht in keinem Artefakt. Der Guard liest Umleitungen nicht. Der resultierende Dateiinhalt ist gelesen und ohne Befund. | `AGENTS.md` §3.1 (Docker-only, Umleitungs-Verbot); Skill HIGH „Docker-only-Verstoß“ | `tools/harness/run-sdk-csharp-integration-tests.sh` · `echo "run-sdk-csharp-integration-tests: TLS-Belege (LH-FA-SST-013, ADR-0150) grün` | nein — **übernommen** aus der Angabe des Implementers; kein Gate-Lauf liest eine Umleitung | Umleitung schreibt Repo-Datei |
| F-2 | MEDIUM | Ein Server-Zertifikat, das **nicht selbstsigniert** ist (von einer privaten CA ausgestellt), als alleiniger Anker scheitert in C# an `PartialChain`; gemessen im Probeprogramm: `leaf-only-anchor: False PartialChain`, `ca-anchor: True`. `SPEC-037`, README und Handbuch sagen „Zertifikat des Servers oder seines Ausstellers“; Kotlin (JDK-Anker) und Python (OpenSSL-Teilkette) sind dafür nicht gemessen. Szenario: Betreiber legt das ausgestellte Server-Zertifikat als Anker in die Datei, der C#-Client verweigert, der Kotlin-Client verbindet — der Inhalt der Option ist nicht in allen drei SDKs derselbe. Die Tests decken nur selbstsigniertes Server-Zertifikat und Aussteller. | `SPEC-037` (Pflichtenheft, „Inhalt in allen drei SDKs derselbe“, „Aussteller oder Server“); `LH-FA-SST-013` | `sdks/csharp/PgChangeFeed.Client/TlsTransport.cs` · `chain.ChainPolicy.TrustMode = X509ChainTrustMode.CustomRootTrust;` | ja — ein Unit-Test mit `CreateIssued`-Server-Zertifikat als alleinigem Anker (heute nicht vorhanden) wird rot | Anker-Variante nicht getragen / ungetestet |
| F-3 | MEDIUM | Die Zusage „mit Anker gelten die Vertrauensanker des Betriebssystems nicht zusätzlich“ (`SPEC-037`) hat in keinem der drei SDKs einen Test; der Plan benennt das selbst unter „Nicht durch eine Mutation gebunden“ (a). Szenario: eine Änderung stellt in C# `TrustMode` auf `System` oder in Kotlin auf den Standard-Speicher um, alle Tests bleiben grün, ein von einer öffentlichen CA ausgestelltes fremdes Zertifikat wird trotz Anker akzeptiert. | `SPEC-037`; Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `docs/plan/planning/in-progress/slice-sdk-tls-optionen.md` · `Nicht durch eine Mutation gebunden:` | ja — ein Test mit einem im Laufzeit-Speicher (z. B. per `SSL_CERT_FILE` im Test-Container) anerkannten Zertifikat und fremdem Anker wird bei der Mutation rot; für Python gilt dies *hergeleitet* (`cadata` ohne Standard-Speicher) | Systemanker-Ausschluss ungebunden |
| F-4 | INFO | `MapException` meldet einen `Internal`-Status mit `HttpRequestException` als `DebugException` jetzt als `PgChangeFeedGrpcUnexpectedStatusException` statt `PgChangeFeedGrpcInternalException`; in Kotlin spricht ein `https`-Adressraum, der bisher (Schema ignoriert) Klartext sprach, jetzt TLS. Beides ist eine sichtbare Verhaltensänderung für Aufrufer, die auf den bisherigen Typ fangen. `make test-sdk-kompat` deckt Konstruktor- und Leseflächen der Fehlertypen, nicht dieses Mapping, und ist nicht gefahren; die Option selbst ist additiv (Zweit-Konstruktor, Python-Default, Kotlin-Zweit-Konstruktor), ein Minor-Sprung ist damit vertretbar. | `SPEC-037` („Versionsgrenze … SemVer-Minor“); `ADR-0145` | `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedAdministrationClient.cs` · `StatusCode.Internal when ex.Status.DebugException is HttpRequestException` | ja — `make test-sdk-kompat` vor dem Release-Zug | Verhaltensänderung ohne Kompat-Messung |
| F-5 | INFO | Die Versionen bleiben `0.6.1`, und die Tags `sdk-csharp-v0.6.1`, `sdk-kotlin-v0.6.1`, `sdk-python-v0.6.1` existieren: `make sdk-pack-*` baut Artefakte gleichen Namens mit anderem Inhalt als die veröffentlichten. Der Plan benennt die Entscheidung (Minor-Sprung gehört in den Release-Zug mit neuer Freigabe); das Anheben ist vor jeder Veröffentlichung Pflicht. | `SPEC-037` (SemVer-Minor); Plan-Zeile „Package-Versionen … nicht geändert“ | `docs/plan/planning/in-progress/slice-sdk-tls-optionen.md` · `die Versionen bleiben (\`0.6.1\`)` | nein — Release-Zug außerhalb des Diffs | Versionsstand vs. Inhalt |
| F-6 | INFO | Weitere Mutationsbelege fehlen und stehen im Plan offen: Python-Ablauf-Test ohne Mutation, die Ausstellerfälle in allen drei Sprachen unmutiert, Kotlin-Realserverlauf mit `SDK_TLS_CA_FILE=/tls/other.pem` nicht gefahren. Die Aussagen sind im Plan als nicht gebunden bzw. *hergeleitet* gekennzeichnet; kein Verstoß gegen die Herkunftsregel. | `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-sdk-tls-optionen.md` · `Nicht gefahren: Kotlin mit` | nein | Mutationsbeleg offen, benannt |
| F-7 | INFO | Die Fixture legt die Wegwerf-Schlüssel (`good-key.pem`) mit Modus 0644 in ein `mktemp -d` mit 0755 ab, damit der Container als `nonroot` sie liest; auf einem Mehrbenutzer-Host sind sie für die Laufzeit lesbar. Im Kommentar benannt (Wegwerf-Material), der Cleanup-Trap entfernt das Verzeichnis; im Arbeitsbaum liegt kein Zertifikat (Prüfung in der Fixture, `git status` leer). | Maintainability | `tools/harness/lib-sdk-tls-fixture.sh` · `chmod 0644 "$SDK_TLS_TMP"/*.pem` | nein | Schlüsselmaterial lesbar im Temp |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Sicherheit Kotlin (`TlsSupport.kt`: `ValidityCheckingTrustManager`, `httpClient`, `ownedChannel`) | geprüft, ohne Befund — Delegat prüft Kette und (über die Engine-/Socket-Varianten) Namen, Gültigkeit zusätzlich; Anker mit nicht-`https`-Adresse wirft; kein Klartext-Rückfall |
| Sicherheit Python (`tls.py`, `options.py`) | geprüft, ohne Befund — `create_default_context(cadata=…)` lädt keinen Standard-Speicher, `check_hostname` bleibt an, Anker plus `http://` wirft, gRPC ohne Schema und ohne Anker bleibt Klartext wie bisher; übergebene eigene Clients/Kanäle erhalten den Anker nicht und das steht in Docstring, README und Handbuch ausdrücklich (kein stilles Versprechen) |
| Sicherheit C#, Namen/Ablauf/Rückfall (`TlsTransport.cs`) | geprüft, ohne Befund außer F-2 — `NameMismatch`/`NotAvailable` der Plattform werden übernommen, Kette wird gegen die Anker neu gebaut (Gültigkeit durch `Build`), kein Überspringen, kein Override des Namens |
| `sdks/` Public-API und Kommentare | geprüft, ohne Befund — Konstruktoren additiv, `make sdk-public-doc-check` Exit 0, `make kommentar-kennungen DIFF=…` ohne Kandidat; READMEs englisch, keine interne Kennung |
| `docs/user/benutzerhandbuch.md`, `docs/user/sdk-e2e-abdeckung.md` | geprüft, ohne Befund — Version 1.100 mit Historienzeile in Betreibersicht ohne Kennung, `make handbuch-public-doc-check` Exit 0, kein Chronik-Satz; der „keine TLS-Einstellung“-Satz entfällt, die eine Resttreffer-Zeile ist die Historie 1.96 |
| Suchlauf und §3.13-Träger | geprüft, ohne Befund — alle sechs Zeilen stimmen (Soll gleich Ist an beiden Ständen); „plaintext gRPC“-Aussagen in `sdks/` sind auf 0 gefallen, die 37 verbleibenden `plaintext`-Treffer sind Klartext als Zweig neben TLS |
| `harness/README.md`, `harness/mk/sdk.mk`, Runner (`run-sdk-*-integration-tests.sh`, `lib-sdk-tls-fixture.sh`) | geprüft, ohne Befund außer F-1/F-7 — `compose.yaml` unverändert (Override im Temp), Wiederherstellung ohne Override, Cleanup-Trap entfernt das Temp-Verzeichnis, `:?`-Guards für `TOOLCHAIN_IMAGE`/`GO_MODCACHE_VOLUME`, kein Zertifikat im Repo, `cryptography>=43` nur im Test-Extra |
| Plan `slice-sdk-tls-optionen` (Zahlen) | geprüft, ohne Befund — Nenner der Mutationsproben (C# 238 vor, 246 nach den zwei Ausstellerfällen) sind als Stand der Probe gekennzeichnet; Lauf-Zahlen sind als gemessen/übernommen erkennbar; nachgemessen nur der Suchlauf |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Umleitung schreibt Repo-Datei · Anker-Variante nicht getragen / ungetestet · Systemanker-Ausschluss ungebunden · Verhaltensänderung ohne Kompat-Messung · Versionsstand vs. Inhalt · Mutationsbeleg offen, benannt · Schlüsselmaterial lesbar im Temp

## Verdikt

**Merge-blockierend:** ja — F-2 und F-3 verlangen eine Fixrunde am Implementer (C#-Verhalten bzw. Tests für die Anker-Variante und den Systemanker-Ausschluss). F-1 betrifft einen Prozessverstoß ohne Korrektur im Diff (Inhalt der Datei ist gelesen und in Ordnung); er ist **übernommen**, nicht gemessen, und geht als Beobachtung an den Planner (Register `BEO-PGC`, Guard liest Umleitungen nicht), ohne Rückgabe-Pfeil an den Implementer. Die DoD-Zeile „Review durchgeführt“ bleibt offen und wird bei der Fixrunde regulär nachgezogen.

**Übergabe:** F-2 und F-3 gehen an den Implementer; F-1 an den Planner; F-4 bis F-7 sind Hinweise (F-4: `make test-sdk-kompat` vor dem Release-Zug an den Verifier/Planner). Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier.
