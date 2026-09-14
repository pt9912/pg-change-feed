# Verifikationsbericht: slice-061 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-061` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8
Sub-Area) und die bindende [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) —
nicht gegen Diff (Reviewer-Aufgabe, bereits abgeschlossen ohne Fixrunde:
`docs/reviews/review-slice-061.md`, vollständig gelesen, aber als Kontext,
nicht als Ersatz für eigene Prüfung übernommen) und nicht gegen realen Bedarf
(Validator — hier nicht ausgelöst, siehe Begründung unten).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan, die
vollständige `ADR-0057`, den vollständigen Review-Report, den tatsächlichen
Diff seit `44618bb` (`git log`, `git diff --stat`), `internal/bootstrap/wiring.go`
(vollständiger Diff `82e4898`), `tools/harness/httpclient/main.go`
(vollständig), den neuen Abschnitt in `tools/harness/run-integration-tests.sh`
(vollständiger Diff `1088cd6`), `compose.yaml`-Diff, `harness/README.md`-Diff,
`tools/schema/nacharbeit-roles.sql` (vollständig) sowie
`internal/adapters/driven/postgresstorage/tableactivation.go` (Ausschnitt
`Registered`/`List`). `make gates` und `make test-integration` wurden in
dieser Sitzung **eigenständig real ausgeführt**, keine Implementer- oder
Reviewer-Behauptung ungeprüft übernommen.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-061-http-api-beispielclient-e2e.md`
zum Stand `HEAD = 4064ce9`. Sieben Commits seit `44618bb` (reiner
`next→in-progress`-Move):

- `8224b4d`/`c40543c` — **nicht Teil dieses Slice-Diffs** (unabhängige
  Lastenheft-/ADR-Commits `LH-FA-SST-008`/`ADR-0061`, Planner-/
  Architect-Läufe) — per Instruktion ausgeklammert, real per
  `git show --stat` je Commit bestätigt: ihre Dateien
  (`spec/lastenheft.md`, `docs/plan/adr/0061-…md`) tauchen in keinem der
  drei Slice-Commits auf.
- `82e4898` — Bootstrap-Verdrahtung vervollständigt
  (`internal/bootstrap/wiring.go`, `harness/image-hash.txt`)
- `1088cd6` — Beispiel-Client, `compose.yaml`-Verdrahtung,
  `run-integration-tests.sh`-HTTP-Rundlauf, `harness/README.md`-Update
- `2e80970` — fünf Implementierungs-DoD-Punkte abgehakt
- `2621d89` — **ebenfalls nicht Teil dieses Slice-Diffs**, real per
  `git show --stat 2621d89` bestätigt: ausschließlich
  `docs/plan/adr/0061-…md` (README-Index-Eintrag), `docs/plan/planning/welle-19.md`,
  vier neue `open/slice-069…072`-Dateien, `roadmap.md` — kein Bezug zu
  `slice-061`s Plan-Tabelle (§3) oder Out-of-Scope (§1). Kein DoD-Mangel,
  gleiches Muster wie in `verify-slice-060.md` bereits dokumentiert (der
  Hauptzweig ist zwischen `next→in-progress` und Closure nicht slice-exklusiv).
- `4064ce9` — Review-Report-Commit (0 HIGH/0 MEDIUM/0 LOW/1 INFO, keine
  Fixrunde), zieht die DoD-Zeile „Review durchgeführt …" nach

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `tools/harness/httpclient/` liefert ein lauffähiges Werkzeug, ruft mind. einen `admin`- und einen `reader`-Endpunkt real per HTTP auf | **erfüllt, selbst reproduziert** | `tools/harness/httpclient/main.go` vollständig gelesen: `RegisterConsumer` (`POST /consumers`, `Authorization: Bearer <adminToken>`, erwartet `201`) und `ListTables` (`GET /tables?...`, `Bearer <readerToken>`, erwartet `200`) — beide über echtes `net/http`, kein Mock, kein `docker exec`. |
| 2 | `compose.yaml` exponiert `CDC_HTTP_ADDR` als Port am Feed-Container | **erfüllt, selbst reproduziert** | `git show 1088cd6 -- compose.yaml` real gelesen: `CDC_HTTP_ADDR: ":8090"` (führende Doppelpunkt-Form, alle Interfaces im Compose-Netz), zusätzlich `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` — kein Host-Port-Publish, konsistent zum repo-weiten Muster (`nats://nats:4222`). |
| 3 | `run-integration-tests.sh` fährt einen echten HTTP-Rundlauf über mindestens eine Fähigkeit je Token-Klasse, realer Erfolgsbeleg | **erfüllt, selbst reproduziert** | `git show 1088cd6` vollständig gelesen: `http_output=$(docker run … go run ./tools/harness/httpclient …)`, `http_status=$?` **direkt** nach der Kommandosubstitution (kein Pipe dazwischen, `AGENTS.md` §3.9-konform); Skript prüft `REGISTERED`-Zeile und `"table":"feed_e2e_full"` in der `ListTables`-Antwort, bestätigt `RegisterConsumer` zusätzlich unabhängig per `docker exec … psql` gegen `cdc.consumer`, und prüft, dass der Feed-Container danach weiterläuft (kein Neustart). Eigener `make test-integration`-Lauf (§2 unten) bestätigt den Text real im Log. |
| 4 | `make test-integration` grün mit dem neuen HTTP-Rundlauf | **erfüllt, selbst reproduziert** | Eigener, vollständiger, ungefiltert ausgeführter Lauf, Exit-Code unmittelbar danach geprüft: `0` (§2 unten). Log-Zeile real vorhanden: `HTTP-API-Rundlauf (LH-FA-SST-006, ADR-0057) belegt — RegisterConsumer real per HTTP mit admin-Token (http-e2e-consumer, cdc.consumer bestätigt), ListTables real per HTTP mit reader-Token (feed_e2e_full in der Antwort)`. |
| 5 | `make gates` grün | **erfüllt, selbst reproduziert** | Eigener, vollständiger, ungefiltert ausgeführter Lauf, Exit-Code unmittelbar danach geprüft: `0` (§2 unten). |
| 6 | Review durchgeführt, Report liegt vor | **erfüllt** | `docs/reviews/review-slice-061.md` vollständig gelesen: 0 HIGH/0 MEDIUM/0 LOW, 1 INFO, keine Fixrunde nötig; die zentrale Prüfung (Bootstrap-Verdrahtungs-Korrektur gegen `ADR-0047`) real gegengeprüft (§3 unten, nicht nur Reviewer-Aussage übernommen). |
| 7 | Doku-Update `harness/README.md` §Sensors | **erfüllt, selbst reproduziert** | `git show 1088cd6 -- harness/README.md` real gelesen: genau ein Satz an die bestehende `make test-integration`-Zeile angehängt, `· seit slice-061` im etablierten Muster, keine Chronik-Sprache darüber hinaus, keine sonstige Änderung an der Datei. |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich Platzhalter (`<…>`), real per Volltext-Lektüre bestätigt — Planner-Arbeit, beginnt laut Rollen-Sequenz (Modul 8) erst nach diesem Bericht. |
| 9 | Reconciliation-Register — entfällt | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` real geprüft: existiert nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `grep -rl "slice-061" docs/plan/planning/observations/` real ausgeführt: kein Treffer — kein neues Verzeichnis, keine neue `evidence/`-Datei. Die beiden in §8 als Treffer benannten Einträge (`coverage-stage-dockerignore-blockiert-tooling`, `test-runner-stiller-ausschluss`) real existent und weiterhin `offen`/1× (unverändert durch diesen Slice, da der Client nie in eine Docker-Stage eingebaut wurde und `run-integration-tests.sh`s bestehendes `-run`-Muster den neuen Block ohne Anpassung mit erfasst). |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Zeilen tragen noch wörtlich `<bei Closure zu füllen>` — real per Volltext-Lektüre bestätigt. Für die Planner-Closure vorgeprüft: Risiko 1 (Bind-Adresse) ist durch `CDC_HTTP_ADDR: ":8090"` real entschärft (Reviewer-Negativbefund bestätigt dies, eigene Lektüre von `compose.yaml` bestätigt es ebenfalls); Risiko 2 (`.dockerignore`/Alpine-`bash`-Fallstrick) ist nicht eingetreten, weil der Client ausschließlich per `go run` im Toolchain-Container läuft, nicht in eine eigene Docker-Build-Stage eingebaut wurde (`git diff 44618bb..HEAD -- Dockerfile` real leer). Beide sind damit reif für den Ausgang „entfallen" bzw. „nicht eingetreten" — die endgültige Formulierung bleibt Planner-Urteil. |
| 12 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/` (real per `ls` bestätigt); die Paarungen suchen in `done/` und sind vor dem `git mv` nicht sinnvoll prüfbar — Slice trägt `Welle: welle-16`, DoD verweist korrekt auf die Welle-16-Closure. |

**Ergebnis §1:** Alle sieben implementierungs-/reviewbezogenen DoD-Punkte
(1–7) sind real erfüllt und selbst reproduziert, nicht nur behauptet. Die
fünf Closure-Punkte (8–12) sind korrekt noch offen und wurden **nicht**
vom Implementer oder Reviewer vorweggenommen — die DoD-Checkbox-Trennung
(7 abgehakt: 5 Implementierung + Review + Doku-Update, 5 offen: sämtliche
Closure-Pflichten) ist sauber.

## 2. Sensor-Läufe (selbst ausgeführt)

**`make gates`** — vollständiger, ungefiltert ausgeführter Lauf:

```
coverage-gate: OK — Coverage 44.40% erfüllt Schwelle 35%
d-check: 485 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 485 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
  Hinweis: tools/harness/httpclient/main.go, tools/harness/natssub/main.go
  liegen in keiner Schicht (erwartungsgemäß, ADR-0057 §Teilfrage 4 /
  §Testabdeckung — Wegwerf-Testwerkzeug, explizit von a-check ausgeschlossen)
```

Exit-Code direkt nach dem ungepipten Aufruf geprüft (`AGENTS.md` §3.9):
**0**.

**`make test-integration`** (voller Compose-Lauf, gegen den echten
Feed-Container): alle `test/integration`-Testfunktionen `PASS`, u. a. der
neue HTTP-Rundlauf-Block real im Log:

```
run-integration-tests: HTTP-API-Rundlauf (LH-FA-SST-006, ADR-0057) belegt —
RegisterConsumer real per HTTP mit admin-Token (http-e2e-consumer,
cdc.consumer bestätigt), ListTables real per HTTP mit reader-Token
(feed_e2e_full in der Antwort): REGISTERED body={"consumer_id":
"http-e2e-consumer","name":"HTTP E2E Consumer","already_registered":false}

LISTED body={"tables":[{"table_id":"tbl-e2e-flow", … "feed_e2e_full" …}
```

Exit-Code direkt nach dem ungepipten Aufruf geprüft (`AGENTS.md` §3.9):
**0**. Der Rundlauf lief **vor** `run-integration-tests.sh`s abschließender
Schema-Negative-Testfunktion (die den Feed-Container beendet) — bestätigt
durch die Log-Reihenfolge und den Skript-Kommentar dazu.

## 3. Bootstrap-Verdrahtungs-Korrektur (`82e4898`) — eigenständig gegen `ADR-0047` geprüft

Nicht die Reviewer-Aussage übernommen, sondern eigenständig nachvollzogen:

1. **`git show 82e4898` vollständig gelesen.** Diff ist exakt zwei Dateien:
   `internal/bootstrap/wiring.go` (vier neue Import-Zeilen für bereits
   bestehende Use-Case-Pakete `list`/`position`/`remove`/`status`, acht
   `Config`-Feldzuweisungen) und `harness/image-hash.txt`
   (Digest-Nachzug). `acknowledge` war bereits vor diesem Commit
   importiert (`grep -n "acknowledge" internal/bootstrap/wiring.go` zeigt
   den Import an unveränderter Stelle, Zeile 46 im aktuellen Stand) —
   keine neue Fremd-Abhängigkeit. Keine neue `pgxpool`/DSN-Konstruktion:
   `apiConsumerState`, `activation`, `enableTables`, `disableTables`,
   `retentionUseCase` sind alle bereits **vor** dem HTTP-Block bestehende
   Instanzen der übrigen Verdrahtung. **Bestätigt: keine neue Fähigkeit.**
2. **DSN-Rollen-Prüfung gegen `ADR-0047` real anhand `tools/schema/nacharbeit-roles.sql`
   nachvollzogen** (nicht nur den Reviewer-Text übernommen): Die Datei
   grantet `cdc_reader` ausschließlich `SELECT` auf vier Views
   (`cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`,
   `cdc.retention_blockers`) — **kein** Grant auf `cdc.source_table`. Die
   Produktionsimplementierung des `TableActivationPort`
   (`internal/adapters/driven/postgresstorage/tableactivation.go`,
   Funktionen `Registered`/`List`) liest über `queries.SelectSourceTable`
   direkt gegen `cdc.source_table` — real im Quellcode bestätigt (Zeile 92:
   `a.pool.QueryRow(ctx, queries.SelectSourceTable, …)`). Eine
   reader-DSN-gebundene Instanz dieses Ports wäre an dieser Stelle real
   nicht lauffähig gewesen (fehlendes Grant, SQLSTATE 42501 zu erwarten).
   `ADR-0057` selbst benennt die Orthogonalität explizit: Der HTTP-Token
   entscheidet, *wer* den Aufruf erreichen darf; welche DB-Rolle die
   dahinterliegende Query fährt, bleibt die bei der Verdrahtung
   (`ADR-0047`) fixierte Zuordnung — `activation` ist seit jeher
   `CDC_ADMIN_DSN`-gebunden, unabhängig davon, ob der Aufrufer über einen
   `reader`- oder `admin`-Token kam. **Eigenes Verdikt: kein
   Least-Privilege-Verstoß, kein `ADR-0047`-Widerspruch — bestätigt.**
3. **Legitimität des Nachholens (Modul 5/7).** Kein rotes Gate wurde bei
   `slice-060`s Closure toleriert (`slice-060` liegt in `done/`, sein
   Review benennt den Defekt nicht — die Whitebox-Testform konnte ihn
   strukturell nicht sehen, siehe §5). Der Fix ist im Slice-Plan §3
   transparent als „Abweichung vom Plan" mit Begründung dokumentiert, in
   der Commit-Message ausführlich benannt, reine Wiederverwendung
   bestehender Instanzen (kein neuer Scope, keine Carveout-Pflicht).

**Eigenes Urteil, unabhängig vom Reviewer-Text:** Die Korrektur ist
technisch korrekt, ADR-konform und legitim im Scope dieses Slice
nachgeholt.

## 4. `httpclient`/`run-integration-tests.sh` — Gegenkontrolle des realen Rundlaufs

- `tools/harness/httpclient/main.go` real gelesen (§1 Punkt 1 oben):
  `POST /consumers` mit `adminToken` (erwartet `201`), `GET /tables?…` mit
  `readerToken` (erwartet `200`) — zwei unterschiedliche Tokens für zwei
  unterschiedliche Rechtsklassen, kein gemeinsamer Token für beide Aufrufe.
- `run-integration-tests.sh`s neuer Block real gelesen: Die
  `RegisterConsumer`-Wirkung wird **zusätzlich unabhängig** per
  `docker exec … psql -tAc "SELECT consumer_id FROM cdc.consumer WHERE
  consumer_id = '$HTTP_CONSUMER'"` bestätigt — kein reines
  Client-Output-Vertrauen; die eigene `make test-integration`-Ausführung
  (§2) bestätigt, dass diese Prüfung real durchlief (Exit 0, kein Abbruch
  bei der `psql`-Prüfung).
- `ListTables`-Beleg prüft den realen JSON-Antwortkörper auf
  `"table":"feed_e2e_full"` — real im selbst erzeugten Log vorhanden (§2).

**Bestätigt: beide Behauptungen (echter HTTP-Aufruf, unabhängige
`psql`-Gegenprüfung) sind real, nicht nur Skript-Text.**

## 5. F-1 (Whitebox-Testform kann Bootstrap-Wiring-Lücken strukturell nicht fangen) — eigene Einschätzung

**Geteilt, mit derselben Begründung, eigenständig nachvollzogen statt
übernommen:** Alle Whitebox-Tests unter
`internal/adapters/driving/http/*_test.go` konstruieren `Config` lokal mit
Fake-Use-Cases (real gegengeprüft — die Testdateien importieren
ausschließlich das eigene Package und `application/port/inbound`, keiner
importiert `internal/bootstrap`). Diese Testform kann per Konstruktion nie
sehen, ob `bootstrap.Run` die `Config`-Felder tatsächlich befüllt — das
liegt außerhalb ihres Objektbereichs, unabhängig davon, wie vollständig die
Handler-Logik selbst getestet ist. Der reale Beleg ist genau das, was in
diesem Slice geliefert wurde: ein Durchlauf über `bootstrap.Run` per echtem
HTTP-Request. Die Lücke blieb zwei Slices (`slice-059`/`slice-060`) lang
unentdeckt, weil `ADR-0057` selbst den einzigen `wiring.go`-durchlaufenden
Beleg (`make test-integration`) bewusst erst „Slice C" zuordnete — das ist
eine architektonische Vorentscheidung, kein Testversäumnis der
Vorgänger-Slices. **Eigenes Urteil: F-1 trägt, die Klassifikation als
INFO (kein Merge-Blocker, aber ein wiederkehrendes strukturelles Risiko
für jeden künftigen Handler in diesem Adapter) ist angemessen.** Ob daraus
ein neuer Beobachtungs-Registereintrag wird, ist Planner-Urteil bei der
Closure — die bestehenden Einträge (`rollen-test-abdeckungsluecken`,
`adapter-fehler-ausgang`) beschreiben real andere Deckungslücken (Rollen-
Grant-Regression bzw. Fehler-Ausgang-Konsistenz), kein Duplikat.

## 6. Welle-16-Closure-Trigger — Erfüllbarkeit geprüft

`docs/plan/planning/welle-16.md` §3 verlangt: (a) `slice-059`, `slice-060`,
`slice-061` liegen in `done/`; (b) `make gates` grün; (c) ein real belegter
E2E-Rundlauf über beide Token-Klassen gegen den laufenden Feed-Container;
(d) Closure-Notiz.

- `slice-059`/`slice-060` liegen real in `done/` (`ls` bestätigt).
  `slice-061` liegt noch in `in-progress/` — korrekt, der `git mv` nach
  `done/` ist Planner-Arbeit nach diesem Bericht, nicht Verifier-Arbeit.
- `make gates`: eigenständig grün, Exit 0 (§2).
- Punkt (c) ist mit **diesem** Slice erstmals real erbracht: Vor
  `slice-061` gab es keinen Beleg, der beide Token-Klassen (`reader` und
  `admin`) über einen echten Netzwerk-Request gegen den laufenden
  Feed-Container zusammenführt — `slice-059`/`slice-060` testeten
  ausschließlich per `httptest.NewServer` mit lokal konstruierter `Config`
  (siehe §5). Der neue Block in `run-integration-tests.sh` deckt exakt
  dieses *Mehr* ab, das keiner der drei Slice-DoDs allein trägt.
- Punkt (d) ist noch offen (§7 des Slice-Plans, Platzhalter) — Planner-Arbeit.

**Urteil:** Der Welle-16-Closure-Trigger ist mit dieser DoD-Bestätigung
inhaltlich **erfüllbar** — die einzigen verbleibenden Schritte sind
Planner-Closure-Arbeit (Lerneintrag, Register, Risiko-Ausgänge, `git mv`
nach `done/`, die drei Paarungen, Welle-Closure-Prozedur Modul 6), keine
weitere Implementierungs- oder Review-Arbeit.

## 7. Hard Rules

- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `44618bb` real per
  `git show --stat` bestätigt: reiner `next→in-progress`-Move (nicht
  Bestandteil des geprüften Bereichs, aber als Elter-Commit bestätigt).
- **3.7 (Kommentar-/Chronik-Disziplin):** `tools/harness/httpclient/main.go`
  und der Kommentar-Umbau in `wiring.go` real gelesen — beide beschreiben
  ausschließlich den Ist-Zustand, keine Chronik-Sprache über den Diff hinaus.
- **3.9 (Exit-Code nie gepiped):** in dieser Sitzung durchgehend beachtet
  (§2) — `make gates`/`make test-integration` jeweils ungefiltert
  ausgeführt, Exit-Code unmittelbar danach in einer eigenen Zeile geprüft.
  Ebenso im geprüften Skript-Diff selbst bestätigt (`http_status=$?` direkt
  nach der Kommandosubstitution, keine Pipe dazwischen).
- **3.1 (Docker-only):** beide Sensor-Läufe liefen über die gepinnten
  Toolchain-/PostgreSQL-/NATS-Container, keine lokale Go-Installation genutzt.

## 8. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 12) — Slice liegt noch in `in-progress/`.
Closure-Notiz, Beobachtungs-Register-Neueintrag und §6-Risiko-Ausgänge
(Planner-Closure-Arbeit, beginnt laut Rollen-Sequenz Modul 8 erst nach
diesem Bericht). Validierung gegen realen Bedarf: **kein Validator-Zug
ausgelöst** — `slice-061` ist kein MVP-Meilenstein-Slice, sondern der
abschließende E2E-Beleg-Slice einer bereits laufenden, additiven
API-Erweiterung; die beiden Validator-Kanten (Modul 8 §Die neun Übergaben)
greifen hier nicht.

## Verdikt

**DoD-Konformität: bestätigt** für alle sieben implementierungs-/
reviewbezogenen Punkte (1–7), jeweils selbst reproduziert (`make gates`,
`make test-integration`, direkte Code-Lektüre gegen `httpclient/main.go`,
`wiring.go`-Diff, `nacharbeit-roles.sql`, `tableactivation.go`). Die fünf
verbleibenden Closure-Punkte (8–12) sind korrekt noch offen und wurden
nicht vorweggenommen.

**Bootstrap-Verdrahtungs-Korrektur (`82e4898`): eigenständig bestätigt.**
Kein Least-Privilege-Verstoß gegen `ADR-0047` — `GetStatus`/`ListTables`
teilen sich den `cdc_admin`-gebundenen `TableActivationPort` mit
`EnableTable`/`DisableTable`, weil `cdc_reader` real kein SELECT-Grant auf
`cdc.source_table` trägt (nur die vier Views) und `ADR-0057` Token-Klasse
und DB-Rolle explizit orthogonal hält. Legitim im Scope dieses Slice
nachgeholt (kein rotes Gate toleriert, transparent im Plan §3 dokumentiert).

**F-1: geteilt.** Die strukturelle Grenze der Whitebox-Testform ist real
nachvollzogen und die INFO-Klassifikation angemessen.

**Welle-16-Closure-Trigger: erfüllbar.** Alle drei Slices liegen inhaltlich
fertig vor (zwei bereits in `done/`, dieser hier DoD-bestätigt), `make
gates` grün, und der reale E2E-Beleg über beide Token-Klassen liegt jetzt
vor — das *Mehr* der Welle gegenüber den einzelnen Slice-DoDs ist erbracht.
Verbleibende Schritte sind reine Planner-Closure-Arbeit.

**Sensor-Läufe:** `make gates` — Exit-Code **0** (ungepipt, unmittelbar
geprüft). `make test-integration` — Exit-Code **0** (ungepipt, unmittelbar
geprüft; HTTP-Rundlauf-Zeile real im Log, unabhängig gegen `cdc.consumer`
bestätigt).

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden — und mit ihm der Welle-16-Closure-Trigger. Für die Closure-Notiz
vorzumerken: beide §6-Risiko-Ausgänge (Bind-Adresse real entschärft durch
`CDC_HTTP_ADDR: ":8090"`; `.dockerignore`/Alpine-Fallstrick nicht
eingetreten, Client läuft ausschließlich per `go run`), der Beobachtungs-
Registereintrag (keine neue Beobachtung durch diesen Slice selbst
angefallen; F-1 als mögliche neue Beobachtung bleibt Planner-Erwägung), der
Fremd-Commit `2621d89` im Diff-Fenster (kein Mangel, Welle-19-Eröffnung
parallel auf dem Hauptzweig), und die anschließende Welle-16-Closure-
Prozedur (Modul 6: Trigger-Audit, Lese-Schritt Beobachtungs-Register,
`welle-16-results.md`, `git mv` nach `done/`, drei Paarungen,
Wave-Self-Close-Commit, Roadmap-Fortschreibung). Kein Validator-Zug
ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
