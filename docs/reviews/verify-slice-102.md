# Verifikationsbericht: slice-102 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-102` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) §Entscheidung
Festlegung 2 (der benannte Zusatzkontext), Festlegung 3 (Stub im Bau, nicht
committet), Festlegung 4 (Kandidaten-Tabelle
`Grpc.Net.Client`/`Grpc.Tools`/`Google.Protobuf`), §Fitness Function (der Bau
bricht ohne den benannten Zusatzkontext ab) sowie
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (Server-Stream,
TLS-Terminierung Betreiber-Pflicht) — und die Hard Rules `AGENTS.md` §3.1
(Docker-only), §3.9 (Exit-Code-Disziplin), §3.12/§3.13 (Beleg-Herkunft,
Träger-Nachzug). **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe,
[`review-slice-102.md`](review-slice-102.md)) und **nicht** gegen realen
Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), `ADR-0090` vollständig (§Kontext, §Entscheidung Festlegung 1–8,
§Verglichene Alternativen, §Fitness Function, §Slice-Schnitt-Empfehlung,
§Re-Evaluierungs-Trigger), den Review-Report (`dad65c1`), die drei
Implementierungs-Commits (`4f27c0f`, `1749438`, `a71b425`) samt Diff-Stat,
`examples/csharp/grpc-client/{Program,Cli,Config,Format}.cs` und die
`GrpcClient.Tests`, `examples/csharp/Directory.Packages.props`,
`examples/csharp/Dockerfile`, `harness/mk/examples.mk`, `harness/README.md`
(Diff), `docs/user/benutzerhandbuch.md` §4 „Zugriff über den
gRPC-Change-Stream" (Diff), das Go-Form-Vorbild `examples/grpc-client/main.go`
gegen `ADR-0060`, sowie die im Slice-Kopf/§8 zitierten
Beobachtungs-Register-Pfade. Zahlen und Befunde aus dem Review waren
**Kontext**, nicht übernommen — jede Aussage dieses Berichts stammt aus einem
hier selbst gefahrenen Lauf oder einer hier selbst gelesenen/abgefragten
Quelle, einschließlich zweier eigener `docker build`-Läufe (einer ohne, einer
mit dem benannten Zusatzkontext, letzterer `--no-cache`) und einer eigenen
NuGet-Registry-Gegenprobe.

**Gegenstand.** Der Vorgang ist `4f27c0f` (LP1, benannter Zusatzkontext) →
`1749438` (LP2, C#-gRPC-Client) → `a71b425` (LP3, Handbuch-/Träger-Nachzug) →
`dad65c1` (Review-Report, 0 HIGH/MEDIUM, 1 LOW, 2 INFO). Der Slice liegt in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt — erwartungsgemäß,
das ist Planner-Arbeit bei der Closure. `git status` ist sauber, kein
paralleler, gegenstandsfremder Commit landete während dieser Sitzung auf
`main`.

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9); kein Lauf im selben Batch wie eine
Folgehandlung.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `docker build --target runtime-grpc examples/csharp` (**ohne** `--build-context proto=proto`) | **1** | reale, unmittelbare Ablehnung: „pull access denied … `docker.io/library/proto:latest`" — Docker interpretiert den unbenannten Kontextnamen `proto` als Image-Referenz; kein Fallback, kein übersprungener Schritt (`Dockerfile:69`, `COPY --from=proto …`) |
| 2 | `make examples-csharp` (Docker-Layer-Cache) | **0** | vier Images gebaut (`:csharp`, `:csharp-sse`, `:csharp-nats`, `:csharp-grpc`), alle Stufen `CACHED` — inkl. `COPY --from=proto cdc/stream/v1/changestream.proto grpc-client/proto/changestream.proto` |
| 3 | `docker build --no-cache --build-context proto=proto --target runtime-grpc examples/csharp` (unabhängige Gegenprobe ohne Cache) | **0** | alle vier Testsuiten real neu gelaufen: „`HttpClient.Tests.dll`: 8/8", „`SseClient.Tests.dll`: 8/8", „`NatsClient.Tests.dll`: 9/9", „`GrpcClient.Tests.dll`: 7/7" — nicht nur Cache-Treffer; der Stub wird real aus der über den Zusatzkontext gelesenen `.proto` erzeugt (kein Fehler in der `dotnet build grpc-client/…`-Stufe) |
| 4 | `make gates` (auf `HEAD` = `dad65c1`, unkontaminiert) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 836 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` |
| 5 | `curl api.nuget.org/v3-flatcontainer/{grpc.net.client,grpc.tools,google.protobuf}/index.json` (eigene Gegenprobe) | — | letzte **stabile** (nicht `-pre`/`-rc`) Versionen: `2.83.0`, `2.84.0`, `3.36.1` — exakt die im Manifest gepinnten und in `ADR-0090` §Entscheidung Festlegung 4 gemessenen Versionen; `4.0.0-rc1`/`-rc2` korrekt ausgeschlossen |
| 6 | `git ls-files examples/csharp \| grep -i "\.pb\.cs\|Grpc\.cs"` | — | kein Treffer — kein committeter C#-Stub, wie `ADR-0090` Festlegung 3 verlangt; `.gitignore` schließt `bin/`/`obj/` bereits aus |
| 7 | `git diff a990143..HEAD --stat -- examples/grpc-client .a-check.yml Dockerfile .dockerignore examples/kotlin` | — | leer — keine der fünf Out-of-Scope-Zusagen aus §1 des Plans verletzt |
| 8 | Lektüre `examples/csharp/grpc-client/Program.cs` gegen `examples/grpc-client/main.go` und `ADR-0060` | — | identisches Credential-Muster (`Http2UnencryptedSupport` ↔ `insecure.NewCredentials()`), identischer Metadata-Key `authorization`/`Bearer `-Präfix, identischer Fehlerpfad (fehlende Adresse/Token → Exit 2, `RpcException`/Stream-Ende → Exit 1) |
| 9 | `git diff a990143..a71b425 -- docs/user/benutzerhandbuch.md harness/README.md` | — | `**Beispiele:**`-Block (Go, C#) in §4 „Zugriff über den gRPC-Change-Stream", Änderungshistorie-Zeile 1.24; `harness/README.md`-Zeile „drei"→„vier" Images korrigiert, benannter Zusatzkontext benannt |
| 10 | `ls docs/plan/planning/observations/BEO-PGC/{handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche,handbuch-versionshistorie-uebersprungen,nicht-blockierender-workflow-alarmmuedigkeit,github-actions-unverifizierbar-lokal,arbeit-ueberholt-stehenden-traeger}/evidence/` | — | Zähler: `handbuch-nicht-nachgezogen…` 3×, `handbuch-versionshistorie…` 3×, `nicht-blockierender-workflow-alarmmuedigkeit` 1×, `github-actions-unverifizierbar-lokal` **7×** (nicht 5×, siehe §5), `arbeit-ueberholt-stehenden-traeger` 8× — die vom Slice erwartete `evidence/slice-102.md` in keinem der fünf Verzeichnisse **noch nicht** angelegt (erwartungsgemäß, Closure-Arbeit) |
| 11 | `grep -n "erreichen" docs/plan/planning/done/slice-095-beispiel-clients-drei.md docs/plan/planning/done/slice-097-umzug-vertragsflaeche.md` | — | Treffer **nur** in `slice-097:107` („Bau-Kontexte erreichen sie heute **nicht**"), **kein** Treffer in `slice-095` — bestätigt Review-F-3 unabhängig |
| 12 | `git log -1 --format=%ad -- docs/plan/planning/done/slice-095-beispiel-clients-drei.md` gegen die Erstellungszeit des Slice-Plans (`git log -1 -- docs/plan/planning/*/slice-102*`) | — | Plan-Erstellung `2026-09-17 10:05:42`, slice-095-Closure `2026-09-17 10:19:41` — die Aussage „`slice-095` (aktuell `in-progress/`)" in §1 des Plans war zum Schreibzeitpunkt **korrekt**; die spätere Closure von `slice-095` macht sie **nachträglich** unzutreffend, ohne dass eine Arbeit *dieses* Slice sie bewegt hätte (kein §3.13-Fall) |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — der `.proto`-Weg im Bau

| Kriterium (§2) | Befund |
|---|---|
| `examples/csharp/Dockerfile` liest die `.proto` über einen zusätzlichen, benannten Bau-Kontext | **erfüllt** — `Dockerfile:69`, `COPY --from=proto …`; `harness/mk/examples.mk` ruft alle vier Bauten mit `--build-context proto=proto` auf |
| Bau bricht **ohne** Kontext real ab (kein Fallback, kein übersprungener Schritt, `ADR-0090` §Fitness Function) | **erfüllt, selbst reproduziert** — #1: reale Ablehnung mit Exit 1, kein stiller Erfolg |
| Bau erzeugt den Stub real **mit** Kontext | **erfüllt, selbst reproduziert** — #2 (Cache) **und** #3 (unabhängige `--no-cache`-Gegenprobe: `GrpcClient.Tests` real 7/7 grün, Stub-Kompilierung im selben Lauf) |

**LP1: erfüllt, gemessen — inklusive der geforderten negativen Probe.**

### LP2 — der Client

| Kriterium (§2) | Befund |
|---|---|
| `examples/csharp/grpc-client/` öffnet real `ChangeStream/StreamChanges`, gibt jede Nachricht aus | **erfüllt** — #8, Form-Vorbild-treu |
| Liest Adresse/Token aus `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER`, Flag-Override | **erfüllt** — `Cli.cs`/`Program.cs`, Fehlerpfad Exit 2 bei fehlender Adresse/Token |
| Netzlos prüfbare Teile getestet, laufen über `examples-csharp` | **erfüllt, selbst gebaut** — #2 (Cache) **und** #3 (`--no-cache`-Gegenprobe: 7/7) |

**LP2: erfüllt, gemessen.**

### LP3 — die Handbuch-Zeile

| Kriterium (§2) | Befund |
|---|---|
| `**Beispiele:**`-Block trägt Go+C# im gRPC-Abschnitt | **erfüllt** — #9 |
| Änderungshistorie-Zeile | **erfüllt** — #9, Zeile 1.24, inhaltlich korrekt |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt, unkontaminiert** — #4, Exit 0, sechs Checks |
| Review durchgeführt, Report vorliegend | **erfüllt** — `review-slice-102.md` (`dad65c1`), 0 HIGH/MEDIUM, 1 LOW, 2 INFO, keine Fixrunde nötig |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-Arbeit bei Closure (siehe §6) |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, inklusive der von
diesem Auftrag verlangten eigenen negativen Bauprobe (Exit 1 ohne Kontext) und
einer eigenen, cache-unabhängigen positiven Bauprobe (Exit 0, 7/7 Tests real
neu gelaufen). `make gates` ist am Gegenstand grün.

---

## 3. ADR-Konformität

| Zusage | Befund |
|---|---|
| `ADR-0090` Festlegung 2 (der benannte Zusatzkontext trägt real, nicht nur als isolierte Probe) | **eingehalten, jetzt real belegt** — die ADR selbst nennt den offiziellen Sprach-Bau als noch ausstehend („der Zusatzkontext selbst ist noch nicht im offiziellen Bau gemessen, nur in einer Probe", §B4-Contra); dieser Slice liefert genau diesen Beleg (#1–#3) |
| `ADR-0090` Festlegung 3 (Stub im Bau erzeugt, nicht committet; `gen/**` bleibt Go-Bindung) | **eingehalten** — #3, #6: kein committeter C#-Stub, Testlauf beweist reale Erzeugung im Bau |
| `ADR-0090` Festlegung 4 (Kandidaten-Tabelle `Grpc.Net.Client`/`Grpc.Tools`/`Google.Protobuf`) | **eingehalten, exakter Treffer** — #5, alle drei Versionen sind die zum Zeitpunkt der ADR gemessenen und aktuell neuesten stabilen Versionen |
| `ADR-0090` §Fitness Function (Bau bricht ohne Zusatzkontext ab) | **eingehalten, selbst reproduziert** — #1 |
| `ADR-0090` Festlegung 6 (`.a-check.yml` unverändert, `examples`-Gruppe bleibt Go-Aussage) | **eingehalten** — #7, kein Diff an `.a-check.yml` |
| `ADR-0060` (Server-Stream-Semantik, TLS-Terminierung Betreiber-Pflicht) | **eingehalten** — #8: `Http2UnencryptedSupport` ist die dokumentierte, zum Go-Vorbild (`insecure.NewCredentials()`) analoge Wahl; Review hatte das bereits bestätigt, hier unabhängig gegengeprüft und bestätigt |
| `AGENTS.md` §3.1 (Docker-only) | **eingehalten** — Bau ausschließlich über `docker build`/`make` |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten** — jeder Lauf dieses Berichts einzeln, ungepiped, Exit-Code direkt gelesen |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten** — `a71b425` zieht `harness/README.md` (Image-Zahl) nach; kein weiterer im Diff bewegter Träger gefunden |
| Out-of-Scope-Disziplin (§1 des Plans, sechs Ausschlüsse) | **eingehalten** — #7: kein Kotlin-Code, kein Umbau des Go-Vorbilds, keine `.proto`-Inhaltsänderung, kein committeter Stub, `.a-check.yml`/Wurzel-`Dockerfile`/`.dockerignore` unberührt |

**Abweichung bei `AGENTS.md` §3.12 — eigener Fund, nicht aus dem Review
übernommen:** Der Slice-Kopf/§6/§8 zitiert
`BEO-PGC/github-actions-unverifizierbar-lokal` mit „**5×**". Die eigene
Zählung (#10) ergibt **7×** — die Zahl war zum Zeitpunkt der Plan-Niederschrift
(`2026-09-17 10:05:42`, praktisch zeitgleich mit `slice-101`) korrekt (die
Evidence-Dateien für `slice-098`/`-099` entstanden erst danach, um 10:58 bzw.
11:47), ist aber jetzt — zum Zeitpunkt dieser Verifikation — überholt, ohne
dass eine Arbeit dieses Slice sie bewegt hat. Dieselbe Klasse wie
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-101.md`
(dort dieselbe Register-Zeile, derselbe Effekt). Ändert an der
Schlussfolgerung nichts (der Eintrag ist bereits `verkörpert`, unabhängig ob
5× oder 7×), gehört aber als Fund in die Closure-Notiz und als weitere
Evidence-Datei dieses bereits bestehenden Beobachtungs-Eintrags.

**Zweiter, nicht-blockierender Fund — Plan-Prosa, kein §3.12/§3.13-Fall:**
Die Aussage in §1 des Plans, `slice-095` sei „aktuell `in-progress/`", war zum
Schreibzeitpunkt korrekt (#12) und wurde durch die *spätere* Closure von
`slice-095` — eine fremde Arbeit, nicht diese — nachträglich unzutreffend.
Das ist kein Verstoß gegen §3.13 (die bewegende Arbeit war nicht die dieses
Slice) und kein Verstoß gegen §3.12 (die Aussage trug ihren Ursprung korrekt
zum Zeitpunkt der Niederschrift); es ist eine benannte Randnotiz für die
Closure-Notiz, keine eigene Beobachtungsklasse.

---

## 4. Risiken aus §6 — Status und Materialisierung (sieben, nicht fünf)

Alle sieben Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler.

| # | Risiko | Materialisiert? |
|---|---|---|
| 1 | Der benannte Zusatzkontext trägt im realen Sprach-Bau nicht | **nicht eingetreten** — #1/#2/#3: der Zusatzkontext trägt exakt wie in `ADR-0090` §Kontext isoliert gemessen, jetzt im offiziellen `examples-csharp`-Bau bestätigt |
| 2 | Eine der drei gRPC-/Protobuf-Bibliotheken nicht mehr auflösbar | **nicht eingetreten** — #3/#5: Bau lief real ohne Cache gegen NuGet, alle drei Pakete real auflösbar |
| 3 | Die C#-gRPC-Werkzeugkette verlangt eine nicht-öffentliche Quelle | **nicht eingetreten** — #3: Bau lief ausschließlich gegen die öffentliche NuGet-Registry, ohne Zugangsdaten |
| 4 | Der fremdsprachige gRPC-Bau kann seinen Stub nicht mehr aus der `.proto` erzeugen | **nicht eingetreten** — #3: Stub-Erzeugung und Kompilierung liefen real und fehlerfrei |
| 5 | Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr | **nicht entscheidbar vor einem realen Runner-Lauf** — der Workflow selbst ist in diesem Diff unverändert (kein `.github/workflows/`-Diff im Slice-Umfang); `AGENTS.md` §3.10 wird durch diesen Slice nicht neu ausgelöst, das bestehende Risiko bleibt am bereits verkörperten Register-Eintrag hängen |
| 6 | Der Handbuch-Nachzug wird vergessen | **nicht eingetreten** — #9: Block (zweite Zeile, C#) und Änderungshistorie-Zeile vorhanden |
| 7 | Ein Träger wird überholt, den dieser Slice nicht anfasst | **kein schädlicher Fund, aber real gefunden und korrigiert** — #9: der §3.13-Suchlauf des Implementers fand und korrigierte `harness/README.md` (drei→vier Images) im selben Slice; Review-F-3 fand zusätzlich eine **falsche** Zieladresse für diesen Risiko-Typ (`slice-095` statt `slice-097` als Beleg für „Bau-Kontexte erreichen sie heute nicht") — bestätigt in #11 |

**Keines der sieben Risiken stellt die DoD-Konformität dieses Vorgangs
infrage.** Risiko 7 braucht bei der Closure die korrekte Zieladresse
(`slice-097`, nicht `slice-095`) für den Beleg des
`arbeit-ueberholt-stehenden-traeger`-Eintrags — das ist Review-F-3, hier
unabhängig reproduziert (#11).

---

## 5. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `dad65c1`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/csharp/Dockerfile` — update | Generator-Stufe (`Grpc.Tools`), benannter Zusatzkontext (`COPY --from=proto …`) | **Plan eingehalten** |
| `examples/csharp/<Paket-Manifest>` — update | `Directory.Packages.props`, drei Pins mit Herkunfts-Kommentar | **Plan eingehalten** |
| `examples/csharp/grpc-client/**` — neu | vorhanden, Form-Vorbild-treu, sieben Tests | **Plan eingehalten** |
| `Makefile`/`harness/mk/*.mk` — update | `examples-csharp` ruft alle vier Bauten mit `--build-context proto=proto` auf | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update | `**Beispiele:**`-Block + Änderungshistorie 1.24 | **Plan eingehalten** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht als eigene Zeile
nennt:** `harness/README.md` (Bild-Zahl-Korrektur, transparent als
§3.13-Suchlauf im Commit-Text ausgewiesen) und `docs/reviews/review-slice-102.md`
(Modul-8-Übergabe, kein stiller Umfangs-Zuwachs).

---

## 6. Verdikt

**DoD-konform, mit zwei benannten, nicht-blockierenden Einschränkungen.**
Alle drei Liefer-Punkte (LP1–LP3) sind gemessen erfüllt — inklusive der von
diesem Slice zum ersten Mal geforderten **negativen** Bauprobe (Bau ohne den
benannten Zusatzkontext bricht real mit Exit 1 ab, kein Fallback) und einer
eigenen, cache-unabhängigen `docker build --no-cache`-Gegenprobe, die belegt,
dass der C#-Stub real im Bau aus der über den Zusatzkontext gelesenen `.proto`
entsteht und alle vier Testsuiten (inkl. `GrpcClient.Tests` 7/7) real grün
laufen. `make gates` läuft unkontaminiert grün (Exit 0, sechs Checks). Die
ADR-0090-Zusagen (Festlegung 2/3/4, §Fitness Function) sind eingehalten, die
`ADR-0060`-TLS-Wahl (`Http2UnencryptedSupport`, analog zu
`insecure.NewCredentials()`) ist unabhängig gegengeprüft und bestätigt. Damit
liefert dieser Slice genau den Beleg, den `ADR-0090` selbst als offen
ausweist: dass der benannte Zusatzkontext nicht nur in der isolierten Probe,
sondern im offiziellen Sprach-Bau trägt.

**Zwei Einschränkungen, keine davon blockiert die DoD-Konformität:**

1. Der im Slice-Kopf/§6/§8 genannte Zähler für
   `BEO-PGC/github-actions-unverifizierbar-lokal` (5×) ist zum
   Verifikationszeitpunkt überholt (real 7×) — zum Schreibzeitpunkt des Plans
   war er korrekt, wurde aber durch zwei zwischenzeitlich geschlossene Slices
   (`slice-098`, `slice-099`) überholt. Dieselbe Klasse wie bei `slice-101`.
   Ändert nichts am Ausgang (Eintrag ist bereits verkörpert), gehört aber in
   die Closure-Notiz und als weitere Evidence-Datei in
   `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`.
2. Review-F-3 (Herkunfts-Zitat zeigt auf die falsche `done/`-Datei —
   `slice-095` statt `slice-097` für die Aussage „ihre Bau-Kontexte erreichen
   sie heute nicht") ist unabhängig reproduziert (#11) und muss bei der
   Closure korrekt adressiert werden: der Beleg für
   `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Risiko 7) gehört gegen
   `slice-097`, nicht `slice-095`.

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die Closure-Checkboxen (LP1–LP3, `make gates`, Verifikation jetzt mit
   diesem Bericht erledigt; Closure-Notiz, Beobachtungs-Register,
   Risiko-Ausgänge, drei Paarungen stehen noch aus).
2. §6 — alle sieben Risiken brauchen ihren formalen Ausgang; keines ist
   *eingetreten* im schädlichen Sinn (§4 dieses Berichts): 1/2/3/4/6 →
   voraussichtlich *entfallen* mit Begründung, 5 → *weiter offen* (Bezug auf
   `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`, bis ein realer
   Post-Push-Lauf sichtbar wird), 7 → *entfallen* mit Begründung (real
   gefunden und im selben Slice korrigiert — Beleg gegen `slice-097`, nicht
   `slice-095`).
3. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag; die naheliegende
   Kandidaten-Regel ist bereits im Plan benannt (erste reale Bauprobe des
   benannten Zusatzkontexts im offiziellen Sprach-Ziel, übertragbar auf
   `slice-103`/Kotlin).
4. Die beiden in §3/§5 dieses Berichts benannten Funde (Zähler-Drift 5×→7×,
   Herkunfts-Zitat `slice-095`→`slice-097`) gehören namentlich in die
   Closure-Notiz, mit je einer weiteren Evidence-Datei in den bereits
   bestehenden Beobachtungs-Einträgen (`zahl-in-traeger-driftet-gegen-die-messung`,
   `arbeit-ueberholt-stehenden-traeger`).
5. Der `git mv` nach `done/` folgt erst nach 1–4 (Modul 5 — Inhalt vor Move
   bei Closure-Übergängen).

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst) und die Übertragbarkeit der hier gemessenen
Form auf `slice-103` (Kotlin) — das ist Gegenstand des Folge-Slice selbst.
