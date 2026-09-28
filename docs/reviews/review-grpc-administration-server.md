# Review-Report: gRPC-`Administration`-Service (ADR-0130) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan/ADR + Konventionen (Modul 10
§Drei Review-Arten). Die DoD-Frage bleibt beim Verifier.

**Gegenstand:** `git diff b317d529..HEAD` — vier Commits:
`9160ad32` (proto + generierter Code), `2c400b8a` (Handler + Interceptor +
Server), `23ad7424` (E2E-Test), `d888c01e` (§3.13-Suchlauf-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (geschärft,
vier repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext:**

- `docs/plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md` (Accepted) — vollständig gelesen
- `spec/pflichtenheft.md` §SPEC-031 (technische Ausgestaltung der ADR)
- `internal/adapters/driving/http/{consumer,registerconsumer,verwaltung,retention}.go` (Kontrakt-Referenz)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `harness/conventions.md` (MR-000/MR-001)
- `.harness/skills/reviewer.md` vollständig

---

## Vorgehen (zusammengefasst)

- Gesamten Diff gelesen, Zeile für Zeile gegen ADR-0130 (alle sechs
  Teilfragen) und gegen die vier HTTP-Handler-Dateien gehalten.
- `git diff --stat` gegen `proto/cdc/stream/v1/changestream.proto` und
  `gen/cdc/stream/v1/*` geprüft — leer, nicht Teil des Diffs (bestätigt
  „byte-identisch unverändert").
- `make generated-sync`, `make a-check`, `make gates`, `make fmt-check`,
  `make test`, `make kommentar-kennungen DIFF=b317d529` selbst ausgeführt,
  Exit-Codes direkt geprüft (nie durch Pipe, AGENTS.md §3.9).
- Eine der drei behaupteten Mutationen selbst nachgefahren (Edit direkt am
  Repo-File, `make test` gezielt, Rücknahme über `git checkout`, Baum
  danach clean bestätigt).
- `make test-integration` komplett und unabhängig frisch durchlaufen
  lassen (nicht nur den Bericht geglaubt) — bis zur letzten Zeile
  `Lauf abgeschlossen` unter `set -euo pipefail`, danach sauberer
  Compose-Teardown bestätigt.
- `spec/pflichtenheft.md`s SPEC-031-Historie-Lücke per `git log -S` gegen
  den Commit-Zeitpunkt geprüft (liegt außerhalb des Diffs).

---

## Findings

### F-1 — Neue gRPC-Betreiber-Oberfläche ohne Handbuch-Zug und ohne benannten Aufschub mit Adresse

- `kategorie`: HIGH
- `quelle`: `.harness/skills/reviewer.md` §HIGH „Neue Betreiber-Oberfläche
  ohne Handbuch-Zug"
- `pfad`: `docs/user/benutzerhandbuch.md` (im Diff nicht berührt);
  betroffen wären mindestens §4 „Aufgaben" (Consumer registrieren,
  Tabelle aktivieren/deaktivieren, Position bestätigen, Retention …) und
  §4 „Zugriff über den gRPC-Change-Stream" (Zeile 1238 ff.) sowie §5
  „Umgebungsvariablen des Feed-Containers" (Zeile 1592: `CDC_GRPC_ADDR`
  wird dort ausschließlich als „Horch-Adresse des gRPC-**Streaming**-
  Servers" beschrieben — diese Beschreibung ist seit diesem Diff
  unvollständig, dieselbe Adresse aktiviert jetzt zusätzlich den
  `Administration`-Service)
- `befund`: Der Diff führt neun neue gRPC-Endpunkte ein (ein neuer
  Service auf der bestehenden Horch-Adresse `CDC_GRPC_ADDR`) und lässt
  `docs/user/benutzerhandbuch.md` vollständig unberührt (`git diff --stat`
  bestätigt: keine Zeile). Es existiert **kein** committetes Artefakt, das
  einen Aufschub mit Adresse (Folge-Slice-ID) benennt: kein
  `docs/plan/planning/{open,next,in-progress}/`-Eintrag zu ADR-0130, kein
  Verweis in den vier Commit-Messages, kein Verweis in der ADR selbst
  (deren Folgepflichten-Liste Beispiel-Clients, SDK-Erweiterung und
  E2E-Beleg nennt, aber keine Handbuch-Zeile). Die Behauptung des
  Implementers, der Aufschub sei „mit dir/dem Auftraggeber abgestimmt",
  ist eine Aussage außerhalb des Repos — sie erfüllt nicht die vom Skill
  geforderte Form „benannter Aufschub **mit Adresse**": eine mündliche
  Abstimmung ist keine Übergabe (siehe Aufgabenstellung §Rollen-Trennung).
  Ohne eine solche Adresse ist nicht geprüfbar, ob und wann der
  Nachzug tatsächlich passiert — genau der Fall, den die Skill-Regel
  verhindern soll (Herkunft: `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  3× vor Aufnahme dieser Regel).
- `verifizierbar`: ja — `git diff --stat b317d529..HEAD -- docs/user/benutzerhandbuch.md`
  ist leer; `git grep -rn "ADR-0130\|grpc-administration\|Administration-Service"
  docs/plan/planning/{open,next,in-progress}` findet keinen Treffer (beide
  real ausgeführt).
- `klasse`: Neue Betreiber-Oberfläche ohne Handbuch-Zug

### F-2 — Kommentarblock mit zwei Kennungen in neuer `gen/`-Testdatei (vom Tool strukturell nicht erfasst)

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §MEDIUM „Herkunft als mehrere
  Felder, Kette, „ff.“ oder Spec-Wiederholung"; `AGENTS.md` §3.7 „Herkunft
  im Go-Kommentar ist ein Feld, kein Absatz"
- `pfad`: `gen/cdc/administration/v1/administration_test.go:1-10`
- `befund`: Der neue Datei-Kopf-Kommentar trägt zwei verschiedene
  Kennungen im selben Block: „…(`SPEC-031`, `ADR-0130`): die Feld-Getter,
  …". Das verletzt §3.7s „höchstens eine Kennung je Kommentarblock".
  `make kommentar-kennungen` meldet hierzu **0 Kandidaten** — nicht weil
  der Block konform ist, sondern weil `PATHS` den gesamten `gen/`-Baum
  strukturell ausnimmt (erzeugter Code; siehe
  `harness/sensors/kommentar-kennungen.md` §Aufruf). Diese Datei ist aber
  eine **von Hand geschriebene** Testdatei unter `gen/`, kein generierter
  Stub — der Ausschluss trifft sie als Kollateralschaden, nicht als
  Zweck. Derselbe Zwei-Kennungen-Musterblock existiert bereits im
  Vorbild `gen/cdc/stream/v1/changestream_test.go:2` (Bestandskandidat,
  außerhalb dieses Diffs) — die neue Datei kopiert das Muster wörtlich
  („dieselbe Testgrenze wie `gen/cdc/stream/v1/changestream_test.go`" laut
  eigenem Kommentar), statt die Gelegenheit zur Korrektur zu nutzen.
- `verifizierbar`: teilweise — die zwei Kennungen sind direkt im Text
  nachlesbar; dass das Werkzeug sie nicht meldet, ist über
  `make kommentar-kennungen DIFF=b317d529` nachvollziehbar (Exit 0, keine
  Ausgabe trotz des Blocks).
- `klasse`: Herkunft als Kette in Go-Kommentar (gen/-Blindstelle)

### F-3 (INFO) — Coverage-Zahl im Bericht weicht von der eigenen Messung ab

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (verwandter Grundsatz, hier nicht als
  Doku-Träger-Verstoß eingestuft, da die Zahl in keinem committeten
  Dokument steht, sondern nur im Handoff-Bericht des Implementers)
- `pfad`: `make gates` / `make coverage-gate` (Lauf)
- `befund`: Der Implementer berichtete eine Coverage von 80.80 % nach dem
  Nachzug der `gen`-Tests. Mein eigener, unabhängiger `make gates`-Lauf
  misst **80.70 %** (`coverage-gate: OK — Coverage 80.70% erfüllt
  Schwelle 80%"). Beide Werte erfüllen die Schwelle (80 %), die Diskrepanz
  hat keine Gate-Auswirkung, ist aber eine kleine Abweichung zwischen
  berichtetem und gemessenem Wert.
- `verifizierbar`: ja — `make gates` (durchgeführt, siehe oben)
- `klasse`: Bericht-vs-Messung-Abweichung (folgenlos)

---

## Bestätigte Implementer-Behauptungen (eigenständig nachgeprüft, nicht nur gelesen)

- **ADR-0130-Konformität (Teilfragen 1–6):** Service-/Paketname, RPC-Form,
  alle neun Nachrichtenschemata (Feldnamen 1:1), Rechtsklassen-Tabelle,
  Fehlercode-Tabelle und Server-Platzierung stimmen exakt mit ADR-0130
  überein — Zeile für Zeile gegen die ADR-Tabellen und gegen
  `spec/pflichtenheft.md` §SPEC-031 gehalten.
- **`changestream.proto`/`gen/cdc/stream/v1/*` unverändert:**
  `git diff --stat b317d529..HEAD` führt keine dieser Dateien;
  `make generated-sync` bestätigt Byte-Gleichheit für beide `.proto`-Quellen.
- **Kein zweiter Domänenpfad:** alle neun Handler in `administration.go`
  rufen exakt dieselben `inbound.*UseCase`-Interfaces wie ihr
  HTTP-Äquivalent auf; Zeile-für-Zeile-Vergleich mit
  `internal/adapters/driving/http/{consumer,registerconsumer,verwaltung,retention}.go`
  bestätigt Feld- und Fehlerpfad-Parität (`administrationError` ==
  `writeDomainError` in der Fehlerklassen-Zuordnung).
- **`authUnaryInterceptor` fail-closed:** Rechtsklassen-Tabelle
  (`administrationRPCRoles`) deckungsgleich mit ADR-0130 Teilfrage 4 (6×
  `roleAdmin`, 3× `roleReader`); `methodName()` ist robust gegen ein
  Format ohne `/` (fällt auf den vollen String zurück, trifft dann keinen
  Tabelleneintrag, landet im fail-closed-Zweig `roleAdmin`). Beide
  Interceptoren (`authStreamInterceptor`/`authUnaryInterceptor`)
  koexistieren konfliktfrei, weil `ChangeStream` ausschließlich
  Streaming- und `Administration` ausschließlich unäre RPCs trägt.
- **Eine der drei Mutationen selbst gefahren:** den `!ok`-Fail-closed-Zweig
  in `authUnaryInterceptor` direkt am Repo-File entfernt (Edit-Tool,
  nicht `sed -i`), `make test` gezielt laufen lassen —
  `TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin`
  färbte exakt wie im Kommentar vorhergesagt rot, alle anderen Pakete
  blieben grün. Rücknahme über `git checkout --`, `git status`/`git diff`
  danach leer.
- **`internal/bootstrap/wiring.go`:** keine Regression für den
  HTTP-only-Fall (unverändertes Verhalten); für `GRPCAddr`-only ohne
  `HTTPAddr` entstehen die vier Consumer-Use-Cases trotzdem (Bedingung
  `cfg.HTTPAddr != "" || cfg.GRPCAddr != ""`), kein Nullwert-Risiko im
  gRPC-Config.
- **Realer E2E-Beleg — selbst reproduziert, nicht nur gelesen:** einen
  komplett frischen `make test-integration`-Lauf gestartet und bis zum
  Ende beobachtet (Hintergrundprozess über ~23 Minuten begleitet). Die
  neue Phase lief exakt wie im Skript beschrieben: `ListTables` (reader),
  `RegisterConsumer` (admin, `cdc.consumer`-Beleg unabhängig bestätigt:
  `LISTED tables=6 retained=2` / `REGISTERED consumer_id=grpc-admin-e2e-consumer`),
  danach `REJECTED code=Unauthenticated` und `REJECTED code=PermissionDenied`.
  Der gesamte Lauf endete mit der letzten Skriptzeile „Lauf abgeschlossen"
  (unter `set -euo pipefail` nur erreichbar, wenn **keine** vorherige
  Prüfung im ~4300-Zeilen-Skript fehlgeschlagen ist), Compose-Teardown
  sauber, Arbeitsbaum danach unverändert (`git status` leer). Das deckt
  Implementer-Behauptung 6 vollständig und eigenständig ab, nicht nur die
  neue Phase, sondern den kompletten Bestands-E2E-Lauf inklusive Backfill,
  Transformationen, Leerlauf-Bestätigung, Upgrade-Tausch.
- **Coverage-Nachzug ist echt, kein Zahlentrick:**
  `gen/cdc/administration/v1/administration_test.go` prüft Nullwert- vs.
  gesetzte Getter mit echten Assertions (nicht nur Aufrufe ohne Prüfung),
  die Fehler-Weiterleitung des generierten unären Clients (Erfolg/Fehler
  je RPC), die neun `UnimplementedAdministrationServer`-Stubs
  (`codes.Unimplemented`) und die Handler-Verklebung aller neun Methoden
  (Dekodier-Fehler, mit/ohne Interceptor) — dieselbe, bereits akzeptierte
  Testgrenze wie `gen/cdc/stream/v1/changestream_test.go` (`ADR-0082`).
- **Gates/Sensoren — alle selbst gefahren, Exit-Codes direkt geprüft:**
  `make generated-sync` (OK, beide Quellen), `make a-check` (0 Befunde),
  `make gates` (Exit 0, komplette Kette inkl. baseline-verify,
  db-package-lists-check, coverage-gate, d-check ×2, commit-traceability,
  generated-sync, a-check), `make fmt-check` (277 Dateien, alle
  formatiert), `make test` (alle Pakete grün, inkl. der beiden neuen
  `gen/…/administration/v1` und `internal/adapters/driving/grpc`),
  `make kommentar-kennungen DIFF=b317d529` (Exit 0 — aber siehe F-2 zur
  Grenze dieses Befunds).
- **§3.13-Suchlauf (Commit `d888c01e`):** `harness/README.md`s
  `make proto-generate`-Zeile und `harness/sensors/generated-sync.md`
  korrekt auf „beide Quellen" nachgezogen; Diff geprüft, Formulierung
  stimmt mit dem tatsächlichen `Dockerfile`-/`proto-generate.sh`-Verhalten
  überein (ein `protoc`-Aufruf für beide `.proto`-Dateien).
- **SPEC-031-Beobachtung ist korrekt als vorbestehend eingeordnet:**
  `git log -S"SPEC-031" -- spec/pflichtenheft.md` zeigt den einzigen
  Treffer bei Commit `7892de6f`; `git merge-base --is-ancestor 7892de6f
  b317d529` bestätigt, dass dieser Commit **vor** dem Review-Startpunkt
  liegt — außerhalb des geprüften Diffs, `spec/pflichtenheft.md` ist in
  `git diff b317d529..HEAD --stat` nicht enthalten. Die Beobachtung
  („SPEC-031 hat keine §6/§7-Gegenzeile wie SPEC-020") ist sachlich
  richtig, aber kein Fund dieses Reviews.
- **Vier Commits, sauberer Arbeitsbaum:** `git status --short` nach `HEAD`
  ist leer; jeder der vier Commits trägt exakt die im Bericht genannten
  Pfade (`git show --stat` je Commit geprüft), keine unerwarteten
  Nebeneffekte.

---

## Negativbefunde

- geprüft, ohne Befund: `proto/cdc/administration/v1/administration.proto`
  (vollständig gegen ADR-0130 Teilfrage 2/3 gehalten)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/administration.go`
  (alle neun Handler, Zeile für Zeile gegen HTTP-Äquivalente)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/interceptor.go`
  (neue Teile: `administrationRPCRoles`, `methodName`,
  `authUnaryInterceptor`)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/server.go`
  (`Config`-Erweiterung, `New()`-Verdrahtung, Interceptor-Koexistenz)
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` (geteilte
  Use-Case-Instanzen, keine Regression für HTTP-only/gRPC-only/keins)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/administration_test.go`
  (Whitebox-Tests, alle an die Eingabeseite gebunden)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/interceptor_test.go`
  (Fitness-Function-Abdeckung vollständig: Unauthenticated,
  PermissionDenied, admin-erreicht-beides, Fail-closed-Fallback)
- geprüft, ohne Befund: `gen/cdc/administration/v1/administration.pb.go`,
  `administration_grpc.pb.go` (byte-identisch zur Generator-Ausgabe,
  `make generated-sync`)
- geprüft, mit Befund F-2: `gen/cdc/administration/v1/administration_test.go`
- geprüft, ohne Befund: `tools/harness/grpcadminclient/main.go`
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (neuer
  Abschnitt, real bis zum Ende durchgelaufen)
- geprüft, ohne Befund: `Dockerfile`, `tools/harness/proto-generate.sh`
  (Zwei-Quellen-Erweiterung, `make generated-sync` bestätigt)
- geprüft, ohne Befund: `harness/README.md`, `harness/sensors/generated-sync.md`
  (§3.13-Nachzug)
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` (auto-generiert,
  Zeilenverschiebung konsistent, script bestätigt „unverändert" bei
  erneutem Lauf)
- geprüft, mit Befund F-1: `docs/user/benutzerhandbuch.md` (im Diff nicht
  berührt)
- geprüft, ohne Befund (außerhalb des Diffs, zur Einordnung der
  Implementer-Beobachtung): `spec/pflichtenheft.md` §SPEC-031

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Neue Betreiber-Oberfläche ohne
Handbuch-Zug · Herkunft als Kette in Go-Kommentar (gen/-Blindstelle) ·
Bericht-vs-Messung-Abweichung (folgenlos)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) ist eine vom Skill explizit
kodifizierte Regel (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
3× vor dieser Regel-Aufnahme) und strukturell erfüllt: neue
Endpunkt-Oberfläche, keine Handbuch-Änderung, kein committeter Aufschub
mit Adresse. Bestreitet der Implementer dies unter Verweis auf eine
mündliche Abstimmung mit dem Auftraggeber, ist das ein Rollen-Widerspruch
auf HIGH-Ebene — nach Modul 8 läuft der Konflikt dann als Sequenz mit
Übergabe-Artefakten über den Architect, statt hier still herabgestuft zu
werden. F-2 (MEDIUM) ist unabhängig vom Ausgang von F-1 zu beheben
(Formfehler, kein Verhaltensfehler) — einfachster Fix: den Datei-Kopf-
Kommentar auf eine Kennung reduzieren (`ADR-0130`) und `SPEC-031` als
Prosa ohne Kennungs-Klammer erwähnen, oder umgekehrt.

Die Code-Substanz selbst (Handler, Interceptor, Server-Verdrahtung,
Fehlercode-Mapping, Wiring, Tests, generierter Code, E2E-Beleg) ist nach
eigenständiger, teils praktischer Nachprüfung (Mutation, voller
`make test-integration`-Lauf) **korrekt und ADR-konform** — keine
Befunde in dieser Substanz.

**Übergabe:** F-1 und F-2 gehen an den Implementer zurück. Die
Finding-Klassen gehen zusätzlich in die Slice-Closure §7 (sofern dieser
Zug über einen Slice-Plan geschlossen wird) und von dort in den
Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; DoD-/
Spec-Konformität prüft der Verifier separat.
