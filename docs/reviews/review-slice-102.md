# Review-Report: slice-102 — 2026-09-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten): geprüft wird der Diff gegen den Slice-Plan, `ADR-0090`
und Modul 5/8 des Baseline-Regelwerks, nicht gegen die DoD (das ist
Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff a990143..HEAD` — Commits `4f27c0f`, `1749438`,
`a71b425` (Slice `slice-102`, Plan
`docs/plan/planning/in-progress/slice-102-grpc-client-csharp.md`).

**Skill:** `.harness/skills/reviewer.md` @ `97ea8fb` · **Modell:**
claude-sonnet-5 · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-102-grpc-client-csharp.md` (§1
  Abgrenzung, §2 LP1–LP3, §4 Trigger, §6 Risiken, §8 Sub-Area-Prüfungen)
- `ADR-0090` (§Kontext — die isolierte Zusatzkontext-Probe, `buildx`
  v0.37.1; §Entscheidung Festlegung 2/3/4; §Fitness Function)
- `ADR-0060` (Server-Stream-Semantik, TLS-Terminierung Betreiber-Pflicht)
- `ADR-0087` (Sprach-Form-Bestand, soweit von `ADR-0090` bestätigt)
- Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice,
  §Offene Risiken werden bei Closure aufgelöst; `modul-08-agentenrollen.md`
  §Rollen-Regeln
- `AGENTS.md` §3.1/§3.5/§3.9/§3.12/§3.13
- Form-Vorbild `examples/grpc-client` (Go, bereits in `main`),
  `tools/harness/grpcclient/main.go` (Wegwerf-Client, gleiches
  Credential-Muster)
- Eigene Nachprüfung mit Netzzugriff: NuGet-Registry
  (`api.nuget.org/v3-flatcontainer/{grpc.net.client,grpc.tools,
  google.protobuf}/index.json`); eigener `docker build`-Lauf gegen
  `examples/csharp/Dockerfile` mit und ohne `--build-context proto=proto`
  (`--target runtime-grpc` und `--target runtime-sse`)

---

## Findings

### F-1 — Zusatzkontext ist strukturell für alle vier Sprach-Programme zwingend, nicht nur für `grpc-client`

- `kategorie`: LOW
- `quelle`: Maintainability (Plan-Genauigkeit gegen `ADR-0090` §Fitness
  Function)
- `pfad`: `examples/csharp/Dockerfile:65-69` (COPY-Zeile in der
  gemeinsamen `build`-Stufe), `harness/mk/examples.mk:8-34`
- `befund`: Eigener Nachbau bestätigt: `docker build --target runtime-sse
  examples/csharp` (ohne `--build-context proto=proto`) bricht mit
  demselben Fehler ab wie `--target runtime-grpc` — „pull access denied,
  repository does not exist" beim Versuch, `docker.io/library/proto:latest`
  zu ziehen. Der Zusatzkontext ist also nicht auf den `grpc`-Aufrufpfad
  begrenzt, sondern zwingend für **alle vier** `docker build`-Aufrufe von
  `examples-csharp` (http-client, sse-client, nats-client eingeschlossen),
  weil die `COPY --from=proto …`-Zeile unbedingt in der von allen vier
  Programmen geteilten `build`-Stufe steht. Das ist real transparent
  dokumentiert (`harness/mk/examples.mk`-Kommentar „unabhängig vom
  angeforderten `--target`", `harness/README.md`-Zeile „jeder der vier
  Aufrufe trägt ihn zwingend") — keine verschwiegene Kopplung. Die DoD-Zeile
  LP1 selbst spricht nur von „der Bau" allgemein und ist damit nicht
  verletzt; die schärfer gefasste Formulierung des Slice-Plans, „für grpc
  zusätzlich" (`ADR-0090` §Fitness Function), trifft die tatsächliche
  Reichweite nicht exakt — sie ist aber keine Zusage des Plans selbst,
  sondern eine ADR-Paraphrase in der Fitness-Function-Zeile. Strukturell
  ist die Kopplung eine Fortsetzung des bereits vor diesem Slice
  bestehenden Musters (eine gemeinsame `build`-Stufe für alle Programme
  einer Sprache) und keine neu eingeführte Vereinfachung dieses Zuges — der
  Implementer hatte innerhalb des vorgegebenen Dockerfile-Musters keine
  andere Option, ohne die `build`-Stufe je Programm aufzuspalten (eine
  Änderung, die `ADR-0090` nicht verlangt und die den Slice-Umfang
  gesprengt hätte).
- `verifizierbar`: ja — `docker build --target runtime-sse examples/csharp`
  ohne `--build-context` reproduziert den Abbruch
- `klasse`: Zusatzkontext-Kopplung breiter als DoD-Wortlaut, transparent dokumentiert

### F-2 — Fehlermeldung bei fehlendem Zusatzkontext liest sich wie ein Registry-/Netzwerkproblem

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `examples/csharp/Dockerfile:69`
- `befund`: Eigener Nachbau bestätigt den vom Implementer berichteten
  Mutationsbeleg exakt: `docker build` ohne `--build-context proto=proto`
  scheitert mit „pull access denied, repository does not exist or may
  require authorization … docker.io/library/proto:latest" — Docker
  interpretiert den unbenannten Kontextnamen `proto` als Image-Referenz.
  Der Fehlermechanismus ist der korrekte, real gemessene (kein Fallback,
  kein stiller Erfolg), aber die Meldung selbst deutet ohne Vorwissen eher
  auf ein Registry-/Auth-Problem als auf einen fehlenden Bau-Kontext hin.
  Der erklärende Kommentar direkt über der COPY-Zeile und in
  `harness/mk/examples.mk`/`harness/README.md` mildert das für jeden, der
  den Code vor dem Bau liest.
- `verifizierbar`: ja — reproduziert wie oben
- `klasse`: Docker-Fehlermeldung bei fehlendem benannten Bau-Kontext selbsterklärend nicht

### F-3 — Herkunfts-Zitat im Slice-Plan verweist auf die falsche `done/`-Datei

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.13, betrifft die Genauigkeit
  der Closure-Vorbereitung — nicht Teil des geprüften Diffs selbst)
- `pfad`: `docs/plan/planning/in-progress/slice-102-grpc-client-csharp.md`
  (§1 „Warum dieser Slice vor `slice-103`" nahe Zeile 60 sowie §8
  „Vorgelagert — offene Beobachtungen sichten") gegen
  `docs/plan/planning/done/slice-097-umzug-vertragsflaeche.md:106-108`
- `befund`: Der Plan zitiert den Satz „ihre Bau-Kontexte erreichen sie
  heute nicht" als „`slice-095`'s §1". Gemessen (`grep -n "erreichen"
  docs/plan/planning/done/slice-095-beispiel-clients-drei.md` → kein
  Treffer): Der Satz steht wörtlich in `slice-097` §1
  („Ausdrücklich NICHT in diesem Slice", zweiter Punkt), nicht in
  `slice-095`. `slice-095` (done) enthält weder das Wort „erreichen" noch
  „Kotlin". Dies ist kein Fund im Diff dieses Zuges (die Plan-Datei ist
  Bestandteil des Kontexts, nicht des Diffs), sondern ein Melde-Fund für
  die Closure: Wird bei der Slice-Closure ein Beleg für
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` gegen `slice-095` angelegt,
  zeigt er auf die falsche Datei — der tatsächlich teilweise überholte
  Träger ist `slice-097` §1.
- `verifizierbar`: ja — `grep` gegen beide `done/`-Dateien
- `klasse`: Herkunfts-Zitat zeigt auf falsche done/-Datei

## Negativbefunde

- geprüft, ohne Befund: `examples/csharp/grpc-client/{Program,Cli,Config,
  Format}.cs` — öffnet real den Server-Stream, Flag-/Env-Parität mit
  `examples/csharp/sse-client`/`nats-client`, Fehlerpfade (fehlende Adresse/
  Token: Exit 2; `RpcException`/Stream-Ende: Exit 1) analog zum Go-Vorbild
- geprüft, ohne Befund: TLS/Plaintext-Wahl — `Http2UnencryptedSupport` ist
  real analog zu `insecure.NewCredentials()` in `examples/grpc-client`
  (Go) und `tools/harness/grpcclient`; deckt sich mit `ADR-0060`
  („TLS-Terminierung vor dem gRPC-Endpunkt bleibt Betreiber-Pflicht") —
  keine Abweichung, korrekt kommentiert
- geprüft, ohne Befund: `examples/csharp/grpc-client/GrpcClient.Tests/*` —
  sieben xUnit-Tests, Flag-Override/Env-Default/Fehlerfälle (Cli),
  Payload/Leerfall (Format); netzlos, Form-Vorbild `sse-client`-Tests
  eingehalten
- geprüft, ohne Befund: `examples/csharp/Directory.Packages.props` — exakte
  Pins (keine Ranges); `Grpc.Net.Client` 2.83.0, `Grpc.Tools` 2.84.0,
  `Google.Protobuf` 3.36.1 gegen NuGet-Registry eigenständig nachgemessen:
  jeweils die neueste **stabile** Version zum Zeitpunkt der Prüfung
  (`Google.Protobuf` 4.0.0 ist `-rc1`/`-rc2` und korrekt ausgeschlossen)
- geprüft, ohne Befund: kein committeter C#-Stub — `git ls-files
  examples/csharp` ohne `.pb.cs`/generierte Protobuf-Dateien,
  `examples/csharp/.gitignore` schließt `bin/`/`obj/` bereits aus
  (`ADR-0090` Festlegung 3)
- geprüft, ohne Befund: `examples/csharp/Dockerfile` — vierte Programmspur
  konsistent zur bestehenden `runtime-sse`/`runtime-nats`-Form, `runtime`
  bleibt letzte (Default-)Stufe; Basis-Digests unverändert gegenüber
  `ADR-0087`; eigener realer `docker build --target runtime-grpc` Lauf
  erfolgreich (Cache-Hit, `dotnet test` bereits real gelaufen)
- geprüft, ohne Befund: `harness/README.md` — „drei"→„vier" Images-Korrektur
  präzise, inklusive der ehrlichen Offenlegung der Zusatzkontext-Kopplung
  (siehe F-1)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` §4 „Zugriff über
  den gRPC-Change-Stream" — `**Beispiele:**`-Block (Go, C#) und
  Versionshistorie-Zeile 1.24 fortlaufend und inhaltlich zutreffend
- geprüft, ohne Befund: Out-of-Scope-Disziplin — kein Kotlin-Code, kein
  Umbau von `examples/grpc-client` (Go), keine `.proto`-Inhaltsänderung,
  `.a-check.yml`/Wurzel-`Dockerfile`/`.dockerignore` unangetastet
  (`git diff --stat` gegen diese Pfade leer)
- geprüft, ohne Befund: Commit-Traceability aller drei Commits (`LH-FA-
  SST-008`, `ADR-0090`/`ADR-0060` je Betreff, kein `SPEC-*`/`ARC-*` im
  Betreff)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusatzkontext-Kopplung breiter als
DoD-Wortlaut, transparent dokumentiert · Docker-Fehlermeldung bei fehlendem
benannten Bau-Kontext selbsterklärend nicht · Herkunfts-Zitat zeigt auf
falsche done/-Datei

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. F-1 und F-2 sind
dokumentierte, real geprüfte Nebenwirkungen einer vom Plan gedeckten
Konstruktion, kein Defekt. F-3 ist ein Melde-Fund für die Closure (kein
Fund im Diff selbst) und blockiert diesen Zug nicht.

Da keine Fixrunde am Implementer nötig ist, wird die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
`docs/plan/planning/in-progress/slice-102-grpc-client-csharp.md` in diesem
Commit selbst auf `[x]` nachgezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** Alle drei Findings gehen als Klassen in die Slice-Closure §7
und von dort in den Zähler (Beobachtungs-Register); F-3 zusätzlich
namentlich, damit die Closure den `BEO-PGC/arbeit-ueberholt-stehenden-
traeger`-Beleg gegen `slice-097` statt `slice-095` anlegt. Dieser Report
selbst ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen.
Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität (inkl. `make
gates` grün) prüft der Verifier separat (Modul 11).
