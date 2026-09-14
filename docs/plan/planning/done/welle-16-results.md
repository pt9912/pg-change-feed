# Welle 16 — HTTP/JSON-API mit Token-Authn — Closure-Notiz

**Welle:** welle-16
**Abschluss:** 2026-09-14
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP/JSON-API mit
  Token-Authn) ist vollständig erfüllt: alle neun bereits Port-gedeckten
  Fähigkeiten (`RegisterConsumer`, `AcknowledgeConsumer`,
  `GetConsumerPosition`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `GetStatus`, `ListTables`, Retention-Lauf) sind über einen neuen
  Driving-Adapter `internal/adapters/driving/http/` erreichbar, geschützt
  durch eine Token-Middleware mit zwei Rechtsklassen
  (`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`) und einheitlichem
  Fehler-Mapping (`400`/`401`/`403`/`404`/`500`).
- [`ADR-0057`](../../adr/0057-http-grpc-api.md) (Accepted, Supersedes
  `ADR-0020`) trifft die Architektur-Entscheidung vollständig — Protokoll
  (HTTP/JSON), Umfang (Verwaltungs-API, kein Changes-Lesen, keine
  Diagnose/Health), Token-Authn (statische Bearer-Tokens, zwei
  Rechtsklassen), Adapter-Platzierung — und wird von dieser Welle in drei
  Slices umgesetzt, ohne dass eine der drei Vorentscheidungen revidiert
  werden musste.
- `slice-059`: Adapter-Grundgerüst, Token-Middleware, `RegisterConsumer`
  als erste exponierte Fähigkeit — real getestet: kein/unbekanntes Token
  → `401`, `reader`-Token gegen schreibenden Endpunkt → `403`,
  `admin`-Token erreicht den Endpunkt. `SPEC-018` neu angelegt.
- `slice-060`: die restlichen acht Fähigkeiten, zentrales Fehler-Mapping
  (`errors.go`), Reader-Endpunkte (`GetConsumerPosition`, `GetStatus`,
  `ListTables`) akzeptieren beide Token-Klassen, die übrigen nur
  `admin`-Token. `RemoveConsumer` behandelt eine unbekannte Kennung als
  Idempotenz (`200`/`removed=false`), `EnableTable`/`DisableTable` melden
  eine an der Quelle fehlende Tabelle real als `404`. `SPEC-018` um die
  acht Endpunkte erweitert.
- `slice-061`: Beispiel-Client (`tools/harness/httpclient/`),
  `compose.yaml`-Port-/Token-Verdrahtung (`CDC_HTTP_ADDR: ":8090"`, alle
  Interfaces im Compose-Netz), realer End-to-End-Rundlauf in
  `tools/harness/run-integration-tests.sh` — mindestens eine Fähigkeit je
  Token-Klasse (`RegisterConsumer` mit `admin`-Token, unabhängig per
  `docker exec … psql` gegen `cdc.consumer` bestätigt; `ListTables` mit
  `reader`-Token, realer JSON-Antwortkörper geprüft) gegen den laufenden
  Feed-Container. Das ist der Beleg für das *Mehr*, das keiner der drei
  Slice-DoDs allein trägt (§Verifikation unten).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- `ADR-0057` traf alle drei Architektur-Fragen (Protokoll, Umfang,
  Token-Authn) bereits vollständig vorab — kein Slice dieser Welle musste
  eine offene Design-Frage klären; Reviewer und Verifier fanden über alle
  drei Slices hinweg zusammen nur 1 HIGH/1 MEDIUM/1 LOW (alle in
  `slice-060`, per Fixrunde behoben) und 2 INFO.
- Der bestehende Driving-Adapter-Zuschnitt
  (`internal/adapters/driving/replication` als Strukturvorbild) und ein an
  einer Stelle gebündeltes Fehler-Mapping (`errors.go`) ließen sich
  verlustfrei von einer auf acht weitere Fähigkeiten skalieren, ohne die
  Token-Middleware selbst anzufassen.
- Der reale End-to-End-Rundlauf in `slice-061` erbrachte genau den
  Nachweiswert, für den `ADR-0057` ihn bewusst als eigenen, dritten Slice
  vorsah: Er deckte eine echte Bootstrap-Verdrahtungslücke aus
  `slice-059`/`slice-060` auf (acht Handler mit `nil`-Use-Case im
  Bootstrap), die deren eigene Whitebox-Unit-Tests strukturell nicht sehen
  konnten.
- Reviewer und Verifier reproduzierten in jedem der drei Slices die
  zentralen Behauptungen real selbst (Token-Middleware-Statuscodes,
  DSN-Rollen-Zuordnung gegen `ADR-0047`, den End-to-End-Rundlauf-Log) statt
  Implementer-Aussagen zu übernehmen.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- `internal/bootstrap/wiring.go` musste in `slice-061` außerplanmäßig
  nachgezogen werden: `apihttp.Config` war bis dahin nur mit
  `RegisterConsumer` befüllt, die acht in `slice-060` gebauten Handler
  liefen am realen Feed-Container mangels Verdrahtung ins Leere
  (`nil`-Use-Case). Reviewer und Verifier prüften die Korrektur
  unabhängig gegen `ADR-0047` (Least-Privilege) und bestätigten sie als
  legitim — reine Wiederverwendung bereits bestehender Instanzen, kein
  neuer Scope, kein Carveout-Fall. Neu registriert:
  `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (1×, unter der
  Schwelle).
- Ein Commit des Implementers in `slice-059` trug `SPEC-018` im Betreff
  (verbotene Struktur-ID nach `AGENTS.md` §5) — vom Planner-Koordinator
  erst beim eigenen `make gates`-Lauf gefangen (kein Vorab-Hook
  existiert), non-interaktiv korrigiert. Drittes Auftreten von
  `BEO-PGC/commit-traceability-kein-vorab-hook` — Schwelle mit dieser
  Welle erreicht, siehe Steering-Loop-Einträge.
- `slice-060`s Review fand 1 HIGH (Slice-Chronik in
  Produktionscode-Kommentar), 1 MEDIUM (fehlender `404`-Test für
  `GetStatus`), 1 LOW (Sechs-vs-Acht-Zählfehler im Slice-Kopf) und 1 INFO
  — Fixrunde behob alle vier, vom selben Reviewer real gegengeprüft. Der
  Chronik-Fund ist das fünfte Auftreten der bereits seit `slice-041`/
  `slice-044`/`slice-052` (Welle 15) 3×-/4×-verkörperten Klasse
  `BEO-PGC/slice-chronik-in-code-kommentar` — reguläre Fortschreibung ohne
  neue Eskalation, die Verteidigungslinie (Reviewer als tragende Instanz)
  trug erneut vor jedem Merge.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Folge-Slice geplant** (kein Ein-Zug-Sensor/keine Ein-Zug-Prosa-Regel):
  ein lokaler `commit-msg`-Git-Hook, der die beiden bestehenden
  Commit-Traceability-Regeln (`ADR-0045`) vor dem `git commit`-Abschluss
  statt erst über `make gates` meldet — Design-Vorgaben (bash-only, kein
  Docker-Aufruf im Hook, Opt-in-Aktivierung) im Architect-Verdikt bereits
  festgelegt: `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`.
  Kein Ein-Zug-Fix, weil ein neues, lauffähiges Skript mit eigenen
  Fehlerfällen und Duplikations-Risiko entsteht (anders als die beiden
  bisherigen Prosa-Präzedenzfälle dieser Beobachtungsklasse) — Umsetzung
  in `slice-073` (`open/`).
  Auslöser: `BEO-PGC/commit-traceability-kein-vorab-hook` (`slice-038`,
  `review-slice-041`, `slice-059` — 3×, Lese-Schritt dieser
  Welle-Closure).
- **Bereits verkörpert, keine neue Aktion dieser Welle** (Feststellung,
  kein neuer Eintrag): `BEO-PGC/slice-chronik-in-code-kommentar` erreichte
  mit `slice-060`s Chronik-Fund den 5. Beleg — die Klasse ist seit
  `slice-052`/`welle-15` bereits als eigener HIGH-Punkt in
  `.harness/skills/reviewer.md` und als Grenz-Klarstellung in
  `.claude/commands/implement-slice.md` Schritt 20 verkörpert
  (`seit slice-052`, siehe `docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`);
  dieses Auftreten bestätigt nur, dass die tragende Verteidigungslinie
  (unabhängiger Reviewer) weiterhin trägt — kein weiterer Architect-Zug
  nötig.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. In
dieser Welle **neu angelegt**: `adapter-unittest-verdeckt-bootstrap-luecke`
(1×, `slice-061`, Ausgang *weiter offen*, unter der Schwelle). In dieser
Welle **auf 3× gehoben und mit Ausgang versehen** (siehe
Steering-Loop-Einträge): `commit-traceability-kein-vorab-hook` (3×, Ausgang
*geplant* → `slice-073`). Durchgesehen ohne Zähler-Änderung in allen drei
Slices (§8-Sichtungen, keine Treffer oder unter der Schwelle unverändert):
`rollen-test-abdeckungsluecken` (2×), `adapter-fehler-ausgang` (2×),
`coverage-stage-dockerignore-blockiert-tooling` (1×),
`test-runner-stiller-ausschluss` (1×), `a-check-null-abdeckung` (bereits
verkörpert seit `welle-1`), `rollen-verdrahtung` (bereits eingetreten seit
`slice-023`). Unverändert, nicht von dieser Welle berührt: alle übrigen
Registereinträge aus `welle-15` und früher.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine, die diese Welle selbst fortsetzen — `welle-16` schließt vollständig
mit ihren drei Slices (`slice-059`, `slice-060`, `slice-061`);
`LH-FA-SST-006` ist vollständig erfüllt. `welle-19` (Live-Change-Streaming,
`LH-FA-SST-008`, `ADR-0060`/`ADR-0061`) ist eine eigenständige, bereits vor
dieser Closure auf dem Hauptzweig eröffnete Folge-Welle für eine andere
Fähigkeit, kein Folge-Slice dieser Welle — siehe §Trigger-Audit unten.
`slice-073` (lokaler `commit-msg`-Git-Hook) ist ein wellenloser Folge-Slice
aus dem Lese-Schritt, nicht aus einem offenen Risiko dieser Welle.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices (`slice-059`, `slice-060`, `slice-061`) liegen in
  `done/`.
- `make gates` grün — Planner-Lauf zur Closure, Exit-Code explizit
  geprüft, nicht durch eine Pipe maskiert (`AGENTS.md` §3.9): Coverage-Gate
  44,40 % ≥ 35 % Schwelle, `d-check` 0 Befunde, `commit-traceability` OK,
  `a-check` 0 Befunde (`tools/harness/httpclient/main.go` und
  `tools/harness/natssub/main.go` erwartungsgemäß schichtlos, konsistent
  mit `ADR-0057`s explizitem `a-check`-Ausschluss für Wegwerf-Testwerkzeug).
- **Welle-16-Closure-Trigger (c)** — der reale End-to-End-Rundlauf über
  beide Token-Klassen: Dieser Beleg liegt bereits **dreifach unabhängig
  bestätigt** vor, aus `slice-061`s `make test-integration`-Lauf, und wird
  hier zusammengefasst statt erneut erzeugt.
  - **Implementer-Behauptung:** `tools/harness/run-integration-tests.sh`s
    neuer HTTP-Block ruft `RegisterConsumer` real per HTTP mit
    `admin`-Token und `ListTables` real per HTTP mit `reader`-Token auf.
  - **Reviewer** (`docs/reviews/review-slice-061.md`): eigene Lektüre von
    `tools/harness/httpclient/main.go` und des Skript-Diffs, eigener
    `make gates`/`make test-integration`-Lauf (Exit 0, ungepipt), realer
    Log-Auszug geprüft (`REGISTERED …`, `LISTED … "feed_e2e_full" …`).
  - **Verifier** (`docs/reviews/verify-slice-061.md` §2/§4/§6): eigener,
    unabhängiger `make gates`/`make test-integration`-Lauf (Exit 0,
    ungepipt), eigene Code-Lektüre von `httpclient/main.go`,
    `compose.yaml`-Diff und dem `run-integration-tests.sh`-Diff, explizite
    Prüfung, dass der Rundlauf **vor** dieser Welle mit keinem der beiden
    Vorgänger-Slices möglich war (beide testeten ausschließlich per
    `httptest.NewServer`, nie über `bootstrap.Run`) — Urteil: „das *Mehr*
    der Welle gegenüber den einzelnen Slice-DoDs ist erbracht".
  - Kein neuer Testlauf in dieser Closure nötig — der vorliegende,
    dreifach unabhängig reproduzierte Beleg trägt den Closure-Trigger.
- **Trigger-Audit der Welle** (Modul 6 §Wellen-Closure-Prozedur, Schritt 2
  — drei Artefaktklassen):
  - **Carveouts:** `docs/plan/carveouts/` enthält ausschließlich
    `.gitkeep` — kein offener Carveout im Repo.
  - **Bootstrap-aware Gates:** keines Teil dieser Welle (HTTP/JSON-API ist
    kein Gate-Gegenstand).
  - **ADR mit Re-Evaluierungs-Trigger:** `ADR-0057`s Re-Evaluierungs-
    Trigger („sobald … der beobachtete Consumer-Bedarf sich um typed
    Contracts oder Streaming erweitert, z. B. ein konkreter
    gRPC-Consumer benennt sich") **ist real eingetreten** — aber bereits
    korrekt über eine **disjunkte** Folge-ADR behandelt: `ADR-0060`
    (gRPC-Server-Streaming) und `ADR-0061` (HTTP/SSE), beide explizit
    „keine Korrektur und keine Supersedes-ADR zu `ADR-0057`" — sie
    entscheiden eine **andere** Fähigkeit (`LH-FA-SST-008`,
    Live-Change-Streaming), nicht `ADR-0057`s eigene Verwaltungs-API.
    `ADR-0057` bleibt inhaltlich unverändert `Accepted`, kein
    `Supersedes ADR-0057` existiert. **Bewusste Feststellung: Trigger
    eingetreten, korrekt separat behandelt, keine Nacharbeit an
    `ADR-0057` selbst nötig.**
- Drei Paarungen (Anker · Folge-Slice · Register): **Anker** — der eine
  Steering-Loop-Eintrag mit `liegt in`-artigem Verweis
  (Architect-Verdikt-Datei) existiert real
  (`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`).
  **Folge-Slice** — `slice-073` existiert real als Datei in `open/`
  (`docs/plan/planning/open/slice-073-commit-msg-git-hook.md`). **Register**
  — beide in dieser Welle berührten Verzeichnisse
  (`adapter-unittest-verdeckt-bootstrap-luecke`,
  `commit-traceability-kein-vorab-hook`) existieren mit nicht leerem
  `evidence/` — grün.

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; alle drei Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
