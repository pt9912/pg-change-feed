# Review — `slice-nats-drittstream-example-csharp-kotlin`

**Review-Art:** Code und Doku gegen Plan, Entscheidung und Hard Rules
(`.harness/skills/reviewer.md`, Modul 10).

**Gegenstand:** `git diff c8778bac~1..374a6baa` (3 Commits: 1× reiner `git
mv` `open` → `in-progress`, 1× reiner Ruhe-Marker-Entfernungs-Commit
(`roadmap.md`, 2 Zeilen), 1× feat) gegen den Slice-Plan zu
`slice-nats-drittstream-example-csharp-kotlin` (bei diesem Lauf unter
`in-progress/`), [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(`Accepted`) und [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
(`Accepted`, volle Matrix).

**Datum:** 2026-09-18.

**Selbst ausgeführt (Exit-Code jeweils direkt und ungepiped, `AGENTS.md`
§3.9):**

- `git show --stat`/`--find-renames` auf allen drei Commits einzeln —
  Beleg zu Hard Rule 3.3 (siehe Negativbefunde).
- `make examples-csharp` — Exit 0, eigener unabhängiger Lauf.
- `make examples-kotlin` — Exit 0, eigener unabhängiger Lauf.
- `make gates` — Exit 0 (d-check 727 Dateien, 0 Befunde · coverage-gate
  82,80 % ≥ 80 % · commit-traceability OK · a-check 0 Befunde ·
  generated-sync OK).
- Realer End-to-End-Smoke-Test gegen die laufende Demo-Umgebung
  (`cdc-examples-feed`/`-postgres`/`-nats`, bereits aus einem früheren Lauf
  hochgefahren), unabhängig von den im DoD notierten Werten nachvollzogen:
  `make example-run-csharp SURFACE=nats-stream` im Hintergrund gestartet,
  danach `INSERT INTO public.orders (customer, amount) VALUES
  ('review-csharp-smoke-test', 77.77)` — Client druckte `change_id=1829-1
  table=public.orders operation=INSERT
  new_image={"id":"6","customer":"review-csharp-smoke-test","amount":"77.77"}`.
  Analog für Kotlin: `INSERT ... ('review-kotlin-smoke-test', 88.88)` —
  Client druckte `change_id=1834-1 ...
  new_image={"id":"7","customer":"review-kotlin-smoke-test","amount":"88.88"}`.
- `curl` gegen `repo1.maven.org/maven2/com/google/code/gson/gson/maven-metadata.xml`
  und das reale `gson-2.14.0.pom` — bestätigt Versions- und
  Transitiv-Abhängigkeits-Claim im `build.gradle.kts`-Kommentar wörtlich
  (siehe Negativbefunde).
- `grep -n "languages" .a-check.yml` — bestätigt `go: ["**/*.go"]` als
  einzigen Sprachschlüssel.

**Verdikt:** 0 HIGH, 0 MEDIUM, 1 LOW (real gefundener, nach
Skill-Definition aber folgenloser Duplikat-Absatz) → **keine Fixrunde am
Implementer nötig.** DoD-Checkbox „Review durchgeführt" wird in diesem
Lauf selbst nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug
ohne Fixrunde).

---

## HIGH

*(keine)*

## MEDIUM

*(keine)*

## LOW

### F-1 — Doppelter Absatz in `docs/user/benutzerhandbuch.md` nach dem C#/Kotlin-Nachzug

- **klasse:** Editier-Artefakt — Absatz beim Ersetzen des Platzhalters
  versehentlich verdoppelt
- **pfad:** `docs/user/benutzerhandbuch.md:883-889`
- **befund:** Der Diff ersetzt den Platzhalter-Bullet „C#/Kotlin: folgen in
  einem eigenen Folge-Vorgang" durch die zwei neuen C#-/Kotlin-Bullets und
  fügt danach den Absatz „Die Beispiele sind zum Lesen und Nachbauen
  gedacht; die E2E-Testclients des Harness liegen unter `tools/harness/`
  (`natsstreamsub`) und sind kein Vorbild." neu ein — derselbe Absatz stand
  im Elternstand bereits unverändert direkt im Anschluss an den
  Platzhalter-Bullet. Ergebnis: der Absatz steht jetzt zweimal identisch
  hintereinander (Zeilen 883–885 und 887–889). Kein semantischer Fehler
  (beide Kopien sagen dasselbe Richtige), aber ein sichtbarer Kopier-Fehler
  im offiziellen Benutzerhandbuch.
- **verifizierbar:** ja — `grep -n "Beispiele sind zum Lesen und Nachbauen
  gedacht" docs/user/benutzerhandbuch.md` liefert drei Treffer; eine dritte,
  unabhängige Stelle an anderer Position im Dokument (Zeile 825, gRPC-
  Abschnitt) ist real vorhanden und korrekt einmalig — die hier relevante
  Dopplung sind die beiden Treffer bei Zeile 883 und 887. `make
  docs-check`/`d-check` fängt das nicht (prüft Referenzen, nicht
  Inhalts-Duplikation): real bestätigt, `make gates` lief mit 0 Befunden
  trotz der Dopplung.

## INFO

- **I-1** — Die „Ansatz — Plan-Korrektur"-Selbstauskunft des Slice-Plans
  (§3) ist akkurat: real geprüft, dass C# keinen neuen NuGet-Pin braucht
  (`System.Text.Json` liegt in der .NET-BCL, `Directory.Packages.props`
  bleibt in diesem Diff unverändert) und dass Kotlins
  `com.google.code.gson:gson:2.14.0`-Pin-Kommentar wörtlich mit einer
  eigenen Live-Messung übereinstimmt (Maven-Central-`maven-metadata.xml`:
  `<latest>`/`<release>` = `2.14.0`; `gson-2.14.0.pom`: genau eine
  `compile`-scope-Laufzeitabhängigkeit,
  `com.google.errorprone:error_prone_annotations:2.48.0`, alle übrigen
  Deps `test`-scope) — dieselbe Rigorosität wie die Vorbild-Kommentare in
  `examples/kotlin/nats-client/build.gradle.kts` und
  `examples/csharp/Directory.Packages.props`.
- **I-2** — Weder [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
  (Festlegung 4, Abhängigkeits-Kandidaten-Tabelle) noch
  [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)
  (Teilfrage 6) nennt oder erwartet eine JSON-Bibliothek für einen
  Kotlin-Beispiel-Client — beide wurden vor der (in
  `slice-nats-drittstream-example-go` getroffenen und dort bereits
  unbeanstandet durchgewinkten) Design-Entscheidung geschrieben, den
  Vollinhalts-Stream-Client an `grpc-client`s Endlosschleife/
  Feldausgabe statt an `nats-client`s Rohdurchreiche-Stil anzulehnen. Die
  dadurch nötig gewordene erste JSON-Fremdbibliothek für ein
  Kotlin-Beispiel ist damit Folge einer bereits akzeptierten Vorentscheidung,
  keine neue architektonische Entscheidung dieses Slice — der
  `SPEC-023`-Nachzug (Klarstellungssatz zur optionalen JSON-Bibliothek)
  beschreibt nur, was geschah, erweitert aber keine bindende Anforderung und
  bleibt innerhalb der bereits geltenden Zulässigkeitsregel („öffentliche
  Fremdmodule", Plural, unbegrenzt). Kein Befund — reine Kontext-Notiz für
  die Architect-Rolle, falls das Muster ein drittes Mal auftritt
  (Steering-Loop-Schwelle).
- **I-3** — `.a-check.yml`s `languages:`-Schlüssel trägt ausschließlich
  `go: ["**/*.go"]` (real geprüft) — C#-/Kotlin-Quellen liegen vollständig
  außerhalb der Prüffläche von `a-check`, und der bestehende
  `examples`-Gruppen-Kommentar zählt ohnehin nur die Go-Beispiele auf. Kein
  Update nötig, keine Lücke.

---

## Negativbefunde (geprüft, ohne Befund)

- Hard Rule 3.3 (git mv + Inhalt = zwei Commits) — alle drei Commits
  einzeln per `git show --stat --find-renames` geprüft: `c8778bac` ist
  reiner `git mv` (0 Zeilen geändert), `fd9cdc49` ändert ausschließlich
  `docs/plan/planning/in-progress/roadmap.md` (Ruhe-Marker-Entfernung, 2
  Zeilen, kein Rename), `374a6baa` trägt jede Inhaltsänderung ohne
  begleitenden Rename. Kein Verstoß.
- Kotlin-Gson-Pin (`build.gradle.kts`) — Version und
  Transitiv-Abhängigkeits-Claim gegen eine eigene Live-Abfrage von Maven
  Central nachgemessen, deckungsgleich mit dem Kommentar.
- `new_image=null`-Übersetzung für ein fehlendes Row Image —
  `FormatTests.cs`/`FormatTest.kt` prüfen beide explizit den `DELETE`-Fall
  mit `new_image=null` (JSON-Literal, nicht leerer String), dieselbe
  Übersetzung wie `internal/adapters/driven/natsstream/publisher.go`s
  `rowImage()` und wie die bereits gefixte Go-Fassung
  (`examples/nats-stream-client/format_test.go`) — beide neuen Sprachports
  haben die Lektion von Anfang an richtig übernommen, kein
  Wiederholungsfall der ursprünglichen Go-Verwechslung.
- Startform-Docstring (`Program.cs`, `Main.kt`) — beide zitieren
  „Startform ist ein Container-Aufruf, kein Host-Aufruf
  ([`ADR-0087`](../plan/adr/0087-beispiel-clients-csharp-kotlin.md)
  Festlegung 3)", wortgleich mit allen bestehenden C#-/Kotlin-Geschwistern
  (`grpc-client`, `http-client`, `nats-client`, `sse-client`); C#/Kotlin
  trugen nie die vom Go-Slice gefundene `go run`-Drift (andere
  Startform-Konvention von Anfang an), also keine vierte Instanz derselben
  Fehlerklasse.
- `.a-check.yml` — `languages:`-Schlüssel trägt ausschließlich `go`,
  C#-/Kotlin-Dateien liegen außerhalb der Prüffläche; kein Update nötig
  (siehe INFO I-3).
- `SPEC-023`-Nachzug (`spec/pflichtenheft.md`) — Präzisierung einer
  bestehenden Zeile (vier → fünf Zugriffs-Oberflächen, `SPEC-024`-Verweis,
  JSON-Bibliothek-Klarstellung), keine neue bindende Anforderung;
  Historie-Zeile korrekt datiert angehängt, keine ADR-würdige Erweiterung
  (siehe INFO I-2).
- `examples/csharp/Dockerfile`/`examples/kotlin/Dockerfile` — neue Stufe
  `runtime-nats-stream` je Sprache, real gebaut (eigener Lauf, s.o.);
  COPY-/Test-Reihenfolge konsistent mit den vier bestehenden Programmen.
- `harness/mk/examples.mk` — fünfter `docker build --target
  runtime-nats-stream`-Aufruf je Sprache, `nats-stream` in
  `example-run-csharp`/`example-run-kotlin`s `$(filter …)`/`$(error …)`
  ergänzt; `example-run-go`s Filter war bereits im Go-Slice erweitert, hier
  korrekt nicht erneut berührt.
- `examples/kotlin/settings.gradle.kts` — `include("nats-stream-client")`
  ergänzt, Plan-Nachzug korrekt benannt und umgesetzt.
- `examples/README.md`/`docs/user/benutzerhandbuch.md` Aggregat-Zahlen
  „fünf × drei, fünfzehn Programme" — real gezählt (`ls examples/*-client
  examples/csharp/*-client examples/kotlin/*-client`): exakt 15
  Programmverzeichnisse existieren; die Aussage ist jetzt wahr.
- `docs/user/benutzerhandbuch.md` Versionshistorie — 1.29 → 1.30 und neue
  Zeile im selben Diff wie die inhaltliche Änderung (abgesehen von F-1
  oben).
- Realer End-to-End-Beleg gegen die Demo-Umgebung — unabhängig von den im
  DoD notierten Werten nachgestellt (beide Sprachen, s.o.), beide Male ein
  reales, unmittelbar zuvor eingefügtes Event empfangen.
- `make gates` — selbst ausgeführt, Exit 0 (Coverage 82,80 %, d-check 0
  Befunde über 727 Dateien, a-check 0 Befunde, generated-sync OK,
  commit-traceability OK).
- Traceability — alle drei Commits nennen mindestens eine `LH-*`/`ADR-*`-
  Kennung im Betreff.
- C#-Paket-Pins (`Directory.Packages.props`) — im Diff unverändert; kein
  neuer Pin nötig, da `System.Text.Json` in der .NET-BCL liegt.
- Import-Grenze (`SPEC-023`) — kein Import eines privaten Paketbaums dieses
  Repositories in `examples/csharp/nats-stream-client/**`/
  `examples/kotlin/nats-stream-client/**`.

---

## Nachzug — DoD-Checkbox ohne Fixrunde

0 HIGH, 0 MEDIUM, ein einziges LOW ohne Reviewer→Implementer-Rückgabe-Pfeil
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" im
Slice-Plan wird deshalb in diesem Lauf selbst auf `[x]` nachgezogen, mit
Verweis auf diesen Report-Pfad. Der Duplikat-Absatz aus F-1 bleibt ein
benannter, nicht blockierender Rest — empfehlenswert vor Closure in einem
trivialen Ein-Zeilen-Fix zu entfernen, aber kein Grund für eine Fixrunde am
Implementer.
