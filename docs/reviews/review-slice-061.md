# Review-Report: slice-061 — 2026-09-14

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-061` — drei Commits, Elter `44618bb` (reiner
`next→in-progress`-Move). Zwischen Elter und erstem Implementierungs-Commit
liegen `8224b4d`/`c40543c` (unabhängige Lastenheft-/ADR-Commits des Planners
— `LH-FA-SST-008`-Umformulierung, `ADR-0061`) — nicht Teil dieses
Slice-Diffs, real per `git show --stat` je Commit bestätigt (ihre Dateien
tauchen in keinem der drei Slice-Commits auf):

- `82e4898` — Bootstrap-Verdrahtung vervollständigt (`internal/bootstrap/wiring.go`,
  `harness/image-hash.txt`)
- `1088cd6` — Beispiel-Client (`tools/harness/httpclient/main.go`),
  `compose.yaml`-Port-/Token-Verdrahtung, `run-integration-tests.sh`-
  HTTP-Rundlauf, `harness/README.md`-Sensors-Update
- `2e80970` — fünf Implementierungs-DoD-Punkte abgehakt

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft
2026-09-13: vier repo-spezifische HIGH-Regeln plus Slice-/Wellen-Chronik-
und Handbuch-Versionshistorie-Regel)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-061-http-api-beispielclient-e2e.md`
  (§1–§8 vollständig gelesen)
- `docs/plan/adr/0057-http-grpc-api.md` — bindend: Testabdeckungs-Erwartung
  (Beispiel-Client, E2E-Beleg), Token-Klassen-Zuordnung (Teilfrage 3), die
  ausdrückliche Orthogonalitäts-Klausel („welche DSN der darunterliegende
  Adapter tatsächlich benutzt, bleibt … die bei der Verdrahtung (`ADR-0047`)
  fixierte")
- `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` — bindend für
  die Bewertung der Wiring-Nachbesserung (DSN-je-Aufrufer-Tabelle)
- `docs/plan/planning/done/slice-059-http-api-adapter-grundgeruest-registerconsumer.md`,
  `docs/plan/planning/done/slice-060-http-api-restliche-faehigkeiten.md`
  (§1, §2, §3, §7, §8 je vollständig gelesen)
- `docs/reviews/review-slice-059.md`, `docs/reviews/review-slice-060.md`,
  `docs/reviews/verify-slice-059.md`, `docs/reviews/verify-slice-060.md`
  — Präzedenz für Test-Zuschnitt (Whitebox `httptest` mit direkt
  konstruiertem `Config{Fake…}`, nie über `bootstrap.Run`/`ConfigFromEnv`)
- `LH-FA-SST-006` (`spec/lastenheft.md`), `SPEC-018` (`spec/pflichtenheft.md`,
  nur gelesen)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.9), §5 Traceability-Regeln
- `harness/conventions.md` (MR-000 ID-Schema)

---

## Der zentrale Prüfpunkt: Bootstrap-Verdrahtungs-Korrektur in `82e4898`

**Behauptung des Implementers:** `internal/bootstrap/wiring.go` hatte
`apihttp.Config` bislang nur mit `RegisterConsumer` befüllt; die acht in
`slice-060` gebauten Handler liefen gegen einen `nil`-Use-Case. Der Fix
nutze ausschließlich bereits bestehende Service-Instanzen, keine neue
Fähigkeit, kein Fehler-Mapping geändert.

**Eigene Prüfung, nicht die Selbstauskunft übernommen:**

1. **`git show 82e4898` vollständig gelesen.** Der Diff ist exakt zwei
   Dateien: `internal/bootstrap/wiring.go` (Import-Zeilen für vier bereits
   existierende Use-Case-Pakete + acht neue `Config`-Feldzuweisungen) und
   `harness/image-hash.txt` (Digest-Nachzug, weil sich der Build-Kontext
   änderte — konform zu `harness/README.md` §Werkzeuge, `make image`-Zeile).
   Keine neue `pgxpool`/DSN-Konstruktion: `apiConsumerState` (Zeile 623)
   existierte bereits für `RegisterConsumer` und wird jetzt zusätzlich für
   `AcknowledgeConsumer`/`GetConsumerPosition`/`RemoveConsumer`
   wiederverwendet; `activation`/`enableTables`/`disableTables`/
   `retentionUseCase` (Zeilen 415–468) sind alle **vor** dem HTTP-Block
   konstruierte, seit langem bestehende Instanzen der übrigen Verdrahtung
   (Aktivierungspfad, Retention-Zyklus). Die Behauptung „keine neue
   Fähigkeit" ist damit real bestätigt.
2. **DSN-Rollen-Prüfung gegen `ADR-0047` (Least-Privilege, `LH-QA-SEC-001`
   …`003`).** `activation` und `apiConsumerState` sind beide an
   `cfg.AdminDSN`/`cdc_admin` gebunden. `ADR-0047`s Aufrufer-Tabelle weist
   „Tabellen-Aktivierung (`postgresstorage.NewTableActivation`,
   `EnableTableUseCase`)" **und** „`RegisterConsumer`/`AcknowledgeConsumer`
   (CLI-Sondermodi, `postgresstorage.NewConsumerState`)" ausdrücklich
   `CDC_ADMIN_DSN`/`cdc_admin` zu — die ADR bindet den **Aufrufer/Adapter**
   an eine Rolle, nicht die einzelne Operation. `GetStatus`/`ListTables`
   nutzen denselben `TableActivationPort` wie `EnableTable`/`DisableTable`
   (`internal/application/usecase/status/service.go`,
   `internal/application/usecase/list/service.go` — beide nehmen exakt
   diesen Port als einzige Abhängigkeit), und dessen einzige
   Produktionsimplementierung
   (`internal/adapters/driven/postgresstorage/tableactivation.go`) liest
   `Registered`/`List` direkt gegen `cdc.source_table` — eine Tabelle, auf
   die laut `tools/schema/nacharbeit-roles.sql` **nur** `cdc_admin` ein
   `SELECT`-Grant hat (`cdc_reader` hat nur `SELECT` auf die vier Views,
   `cdc.source_table` ist keine davon). Eine reader-DSN-gebundene Instanz
   wäre an dieser Stelle real nicht lauffähig gewesen — es gibt keine
   Alternative, die dem Fix zur Verfügung stand. `ADR-0057` selbst
   benennt diese Orthogonalität explizit: „welche DSN der darunterliegende
   Adapter tatsächlich benutzt, bleibt unverändert die bei der Verdrahtung
   (`ADR-0047`) fixierte" — der HTTP-Reader-Token entscheidet nur, *wer*
   den Aufruf erreichen darf, nicht welche DB-Rolle die dahinterliegende
   Query fährt. **Befund: kein Least-Privilege-Verstoß, kein
   ADR-Widerspruch — der Fix reiht sich korrekt in eine vorab getroffene,
   `Accepted`-ADR-Zuordnung ein.**
3. **Scope-Frage — lag der Fund/Fix wirklich außerhalb des ursprünglichen
   `slice-059`/`060`-Scopes?** Bestätigt über `review-slice-059.md`
   (Negativbefund: „`internal/bootstrap/wiring.go` — die drei neuen
   Umgebungsvariablen … additiv gelesen") und `review-slice-060.md`
   (Negativbefund: „Scope-Treue … keine Änderung an … `internal/bootstrap/
   wiring.go`") — `slice-060` fügte acht Handler und `Config`-Felder in
   `internal/adapters/driving/http/` hinzu, fasste `wiring.go` aber
   nachweislich **nicht** an. Beide Slices testeten ausschließlich über
   `httptest.NewServer` mit direkt konstruiertem `Config{Fake…UseCase}`
   (real gegengeprüft: `retention_test.go`, `verwaltung_test.go`,
   `consumer_test.go` — jeder Testaufruf baut `Config` lokal im
   `internal/adapters/driving/http`-Package, keiner läuft über
   `bootstrap.Run`/`ConfigFromEnv`). Diese Testform kann eine fehlende
   Zuweisung in `wiring.go` strukturell nicht sehen — sie liegt außerhalb
   ihres Objektbereichs. `ADR-0057` selbst verschiebt den einzigen
   Nachweis, der `wiring.go` real durchläuft (`make test-integration`),
   ausdrücklich auf „Slice C" (= `slice-061`, Konsequenzen/Folgepflicht,
   Slice-Schnitt-Empfehlung). **Befund: Der Fund entstand strukturell
   zwangsläufig erst in `slice-061` — weder `slice-059` noch `slice-060`
   hätten ihn mit ihrer eigenen Testform fangen können, und beide
   Reviews/Verifikationen haben ihren jeweiligen Scope korrekt geprüft.**
4. **Legitimität des Nachholens innerhalb `slice-061` (Modul 5/7).** Kein
   rotes Gate wurde bei `slice-060`s Closure toleriert — der Defekt war zu
   diesem Zeitpunkt unbekannt, kein dokumentierter Carveout nötig oder
   möglich. Der Fix behebt eine bereits akzeptierte, aber fehlerhafte
   Verdrahtung eines Vorgänger-Slice, gefunden beim reellen Testaufbau des
   direkten Folge-Slice, transparent im Slice-Plan §3 als „Abweichung vom
   Plan" mit Begründung dokumentiert (nicht stillschweigend verschwiegen)
   und in der Commit-Message ausführlich begründet. Das entspricht der in
   `AGENTS.md` §6 Schritt 4 vorgesehenen „kleinste sinnvolle Änderung" —
   kein eigener Carveout-Fall, keine eigene Slice-Abspaltung nötig, weil
   der Fix ohne ihn DoD-Punkt 3 dieses Slices strukturell unerfüllbar
   gemacht hätte (Blocker-Charakter, nicht Erweiterung).

**Verdikt zu diesem Punkt: legitim.** Der Implementer hat den Fund korrekt
eingeordnet, im Scope belassen (reine Wiederverwendung bestehender
Instanzen) und transparent dokumentiert; die eigene Prüfung bestätigt sowohl
die technische Behauptung als auch die ADR-Konformität der DSN-Zuordnung.

## Weitere Findings

Keine HIGH/MEDIUM/LOW-Findings. Ein INFO unten.

### F-1 — Whitebox-Testform des HTTP-Adapters kann Bootstrap-Wiring-Lücken strukturell nicht fangen

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/adapters/driving/http/*_test.go` (Package-weit,
  Muster über `slice-059`/`slice-060`/`slice-061` hinweg unverändert)
- `befund`: Alle Whitebox-Tests des HTTP-Adapters konstruieren `Config`
  direkt mit Fake-Use-Cases im Adapter-Package und laufen nie über
  `bootstrap.Run`/`ConfigFromEnv`. Das ist genau der Grund, warum die in
  diesem Slice behobene Bootstrap-Verdrahtungslücke zwei Slices lang
  (`slice-059`/`slice-060`) unentdeckt blieb, bevor der reale E2E-Rundlauf
  in `slice-061` sie fand — kein Fehlverhalten der beiden Vorgänger-Reviews
  (siehe oben), sondern eine strukturelle Grenze dieser Testform. Ein
  künftiger neuer Handler in diesem Adapter trägt dasselbe Risiko, bis der
  nächste E2E-Slice ihn real aufruft.
- `verifizierbar`: nein — kein Gate-Lauf bestätigt eine Abwesenheit
  zukünftiger Wiring-Lücken; die Beobachtung selbst ist der Beleg.
- `klasse`: „Whitebox-Adapter-Test ohne Bootstrap-Durchlauf verdeckt
  Verdrahtungslücken"

Kein erwarteter Reviewer-Schritt daraus — Hinweis an die Planner-Rolle zur
Erwägung, ob dies als neuer Eintrag ins Beobachtungs-Register
(`docs/plan/planning/observations/`) gehört, da die bestehenden Einträge
(`BEO-PGC/rollen-test-abdeckungsluecken`, `BEO-PGC/adapter-fehler-ausgang`)
eine andere Deckungslücke beschreiben.

## Negativbefunde

- geprüft, ohne Befund: Scope-Treue — `git show --stat` je Commit
  (`82e4898`, `1088cd6`, `2e80970`) bestätigt: ausschließlich
  `internal/bootstrap/wiring.go`, `harness/image-hash.txt`,
  `tools/harness/httpclient/main.go`, `compose.yaml`,
  `tools/harness/run-integration-tests.sh`, `harness/README.md`, die
  Slice-Plan-Datei selbst — deckungsgleich mit Plan §3 (inkl. der dort
  dokumentierten Abweichung); keine Änderung an
  `internal/adapters/driving/http/*`, `spec/pflichtenheft.md`,
  `.a-check.yml`, keine gRPC-Artefakte.
- geprüft, ohne Befund: `tools/harness/httpclient/main.go` — Godoc- und
  Inline-Kommentare beschreiben ausschließlich den Ist-Zustand
  (Zweck, Träger-Skript, Verhalten der `call`-Hilfsfunktion); keine
  Slice-/Wellen-Chronik, keine Vorher/Nachher-Sprache, kein `<!-- -->`-Nur-
  Norm-Fall (`AGENTS.md` §3.7).
- geprüft, ohne Befund: der Kommentar-Umbau in `internal/bootstrap/wiring.go`
  Zeilen 615–621 — beschreibt den geltenden Zustand („die Tabellen-
  Verwaltung und der Retention-Lauf nutzen dieselben Use-Case-Instanzen wie
  die übrige Verdrahtung"), keine abwesende Alternative, kein Abbruch
  mitten im Satz.
- geprüft, ohne Befund: `compose.yaml` — `CDC_HTTP_ADDR: ":8090"` (führende
  Doppelpunkt-Form bindet auf allen Interfaces, löst das in Plan §6
  benannte Risiko); kein Service in `compose.yaml` published Host-Ports
  (repo-weit konsistent, `httpclient` erreicht den Feed-Container über den
  Netzwerk-Alias `pg-change-feed`, analog `nats://nats:4222`) — „exponiert
  als Port" bedeutet hier konsistent mit dem Bestand: TCP-Listen-Port im
  Compose-Netz, kein Host-Publish.
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` —
  `http_status=$?` direkt nach der Kommandosubstitution ohne
  zwischengeschaltete Pipe (`AGENTS.md` §3.9-konform); `RegisterConsumer`
  wird zusätzlich unabhängig per `docker exec … psql` gegen `cdc.consumer`
  bestätigt (kein reines Client-Output-Vertrauen); `ListTables`-Beleg prüft
  den realen JSON-Antwortkörper auf die dauerhaft aktivierte Tabelle.
- geprüft, ohne Befund: `harness/README.md` §Sensors — die
  `make test-integration`-Zeile ist um genau einen neuen Satz erweitert,
  trägt `· seit slice-061` im etablierten Muster, keine Chronik-Sprache
  darüber hinaus.
- geprüft, ohne Befund: DoD-Checkbox-Commit `2e80970` — ausschließlich
  Checkbox-Toggles und die Plan-Tabellen-Ergänzung der dokumentierten
  Abweichung, kein Inhalts-Nachzug an anderer Stelle.
- geprüft, ohne Befund: Traceability — alle drei Commit-Betreffs tragen
  `LH-FA-SST-006`/`ADR-0057`, keiner trägt eine `SPEC-*`/`ARC-*`-Kennung im
  Betreff.
- geprüft, ohne Befund: `make gates` — Exit-Code direkt geprüft (nicht
  gepiped), Exit 0: `baseline-verify`, `docs-check` (2×, inkl.
  Commit-Range), `commit-traceability`, `coverage-gate` (44.40 % ≥ 35 %
  Schwelle), `a-check` (0 Befunde; `tools/harness/httpclient/main.go` und
  `tools/harness/natssub/main.go` als erwartungsgemäß schichtlos gemeldet,
  konsistent mit `ADR-0057`s expliziter `a-check`-Ausschlussbegründung).
- geprüft, ohne Befund: `make test-integration` — Exit-Code direkt geprüft
  (nicht gepiped), Exit 0; HTTP-Rundlauf-Zeile real im Log:
  `REGISTERED body={"consumer_id":"http-e2e-consumer",…}`,
  `LISTED body={"tables":[…,"table":"feed_e2e_full"…]}`, unabhängig
  gegen `cdc.consumer` per `psql` bestätigt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Whitebox-Adapter-Test ohne
Bootstrap-Durchlauf verdeckt Verdrahtungslücken

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW; das eine INFO
erwartet keine Implementer-Aktion.

**Übergabe:** Keine Fixrunde nötig. Gemäß `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde zieht dieser Report die DoD-Zeile
„Review durchgeführt …" im selben Commit, der ihn anlegt, auf `[x]` nach.
Das INFO-Finding geht als Hinweis an die Planner-Rolle für die
Slice-Closure (§7/Beobachtungs-Register-Erwägung), keine Rückkante an den
Implementer. Dieser Report ist ein Lauf-Beleg (Audit: dieser Diff, dieser
Skill, dieses Modell, dieses Verdikt) und wird über Läufe hinweg nicht
wieder gelesen. Verifikation gegen DoD/Spec bleibt Aufgabe des Verifiers
(Modul 11).
