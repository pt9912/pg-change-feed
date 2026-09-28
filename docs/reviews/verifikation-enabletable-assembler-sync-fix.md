# Verifikations-Report: EnableTable/DisableTable-Assembler-Sync-Fix (`LH-FA-CFG-001`/`LH-FA-CFG-002`) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Konformitätsprüfung
gegen einen kritischen Datenverlust-Fix, in frischem Kontext, nach
Implementierung, Review und Fixrunde.

**Gegenstand:** `git log --oneline 17b5efd4..61e7dcf3` — drei Commits:
`a40b4809` (Implementierung), `eeaa2895` (Review-Report), `61e7dcf3`
(Fixrunde). `58a32bab` (`plan(welle): welle-sdk-grpc-administration-flaeche
eroeffnet`) liegt in diesem Bereich, ist aber ein unabhängiger, wellenloser
Planungs-Zug ohne Berührung der geprüften Dateien und **nicht** Gegenstand
dieser Prüfung (Auftrag).

**Eingangs-Kontext:**

- [`review-enabletable-assembler-sync-fix.md`](review-enabletable-assembler-sync-fix.md)
  — vollständig gelesen (1 HIGH F-1, 1 LOW F-2, Verdikt merge-blockierend
  wegen F-1, fachliche Korrektheit bereits vollständig bestätigt)
- `git show a40b4809` — vollständig gelesen
- `git show 61e7dcf3` — vollständig gelesen (Fixrunde)
- [`LH-FA-CFG-001`](../../spec/lastenheft.md), [`LH-FA-CFG-002`](../../spec/lastenheft.md)
  im Lastenheft nachgeschlagen (Happy-Path-Wortlaut: „dann werden fortan
  Änderungen an `t` erfasst")
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13, §3.15

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien`; `d-check: 1398 Datei(en) geprüft, 0 Befund(e)` (voller Modul-Bündel-Lauf, `docs-check`); `commit-traceability` positive Hälfte (`--enable commits --range HEAD~5..HEAD`): `1398 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK` (beide `.proto`-Quellen, alle vier generierten Dateien geprüft); `a-check: gesamt: 0 Befund(e)`; `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%` |
| `make a-check` (isoliert) | **Exit 0** | `gesamt: 0 Befund(e)` |
| `git grep -c "assembler\.AddBinding(\|assembler\.RemoveBinding(" -- '*.go'` ohne `_test.go` | **genau 2**, beide in `internal/bootstrap/assemblersync.go` | bestätigt die Kernbehauptung: ein einziger geteilter Mechanismus, keine zweite Aufrufstelle |
| `git grep -n "replication/mapper" -- 'internal/adapters/driving/http/*.go' 'internal/adapters/driving/grpc/*.go'` | **0 Treffer** | keine neue Adapter-zu-Adapter-Kante — deckungsgleich mit `a-check: 0 Befunde` |
| `make kommentar-kennungen DIFF=17b5efd4` | **Exit 0** | keine Ausgabe, 0 Kandidaten über den gesamten Drei-Commit-Diff (inkl. der Fixrunde) |
| `make test-integration` (real, Hintergrundprozess, gegen ein frisch aus `61e7dcf3` gebautes Image, `make gates` hat dafür bereits `make image` intern gebaut) | **komplett durchgelaufen** | siehe §4 — alle drei neuen Phasen (gRPC-Enable, gRPC-Disable, HTTP-Enable) real bestätigt, Lauf endet mit „Lauf abgeschlossen — E2E-Abdeckungstabelle aus 17 Go-Zeilen und 49 Bash-Zeilen"; kein `FAIL`/`Error` im gesamten Log; `docker ps -a --filter name=cdc-test` danach leer (sauberer Teardown) |
| `git status --short` (vor meinem Commit, nur meine Zieldateien betrachtet) | sauber für alle geprüften Dateien | ein **unabhängiger, gleichzeitig laufender Agentenlauf** (slice-sdk-csharp-grpc-administration-flaeche) hinterließ währenddessen uncommittete Änderungen unter `sdks/csharp/**` sowie neue Commits auf `HEAD` (`4794a4ee`, `eaa30cc6`, `4034e7c6`) — `git diff --stat 61e7dcf3..HEAD -- <meine sieben Zieldateien>` liefert **keine Ausgabe**: keine dieser fremden Commits berührt eine der hier geprüften Dateien |

Alle Sensoren dieser Tabelle wurden in diesem Lauf **selbst ausgeführt** —
keine Behauptung des Implementers oder Reviewers wurde ohne eigenen Beleg
übernommen (`AGENTS.md` §3.12 Instanz B, Modul 11 §Falle).

**Umgebungshinweis (kein Befund an diesem Fix):** Während dieser Verifikation
lief parallel ein unabhängiger Agentenlauf an einem anderen Slice
(`slice-sdk-csharp-grpc-administration-flaeche`, C#-SDK) im selben
Arbeitsbaum. Er erzeugte drei zusätzliche Commits auf `HEAD` sowie
uncommittete Änderungen unter `sdks/csharp/**`. Keiner dieser Commits oder
Änderungen berührt eine der sieben hier geprüften Dateien (per `git diff
--stat` bestätigt); mein `make gates`-Lauf (§1) und der reale
`make test-integration`-Lauf (§4) liefen beide erfolgreich trotz dieser
gleichzeitigen fremden Aktivität im Arbeitsbaum.

## 2. F-1 (HIGH) — Vorher/Nachher-Chronik in Produktionscode-Kommentar: geprüft, vollständig behoben

**Behauptung der Fixrunde (`61e7dcf3`):** die drei benannten Stellen
(`assemblersync.go:18-23`, `assemblersync.go:68`,
`run-integration-tests.sh` `abdeckung_declare`-Beschreibungstext) sind auf
Ist-Zustand umformuliert; der Anker-Mechanismus (4. Argument von
`abdeckung_declare`) bleibt unverändert, nur das 3. Argument (Beschreibungstext)
ist betroffen.

**Eigen geprüft (nicht nur gelesen):**

- `git show 61e7dcf3` direkt am Diff gelesen: der Godoc über
  `syncAssemblerAddBinding` trägt jetzt „… aufrufen — beide Pfade tragen
  dieselbe Aktivierung unabhängig vom Auslöser in den laufenden Prozess nach
  (`LH-FA-CFG-001`)" — kein „bislang", kein „blieb … unerfasst", reiner
  Ist-Zustand, Subjekt bleibt der Produktionscode-Pfad.
- Der Godoc über `enableTableWithAssemblerSync` trägt jetzt „… nach, damit
  eine über diesen Weg aktivierte Tabelle ab diesem Aufruf erfasst wird
  (`LH-FA-CFG-001`)" — keine „Ohne-X-Konsequenz"-Rhetorik mehr.
- Die `run-integration-tests.sh`-Zeile trägt jetzt „… derselbe
  Live-Reload-Vertrag wie der SQL-Antragsqueue-Pfad, wegunabhängig" — kein
  „bis dahin nicht nachgetragenen" mehr.
- **Anker-Mechanismus-Probe (kritisch für die Wiederholbarkeit des
  E2E-Laufs):** im Diff der `abdeckung_declare`-Zeile ist ausschließlich das
  3. Argument (Freitext) verändert; das 4. Argument
  (`"Direkter gRPC-Zugriffsweg — EnableTable"`, der Anker-Suchtext, den das
  Skript später beim Druck der tatsächlichen Zeile sucht) ist **byte-gleich**
  geblieben — per `git show 61e7dcf3` direkt am Einzeiler-Diff gegengelesen.
  Bestätigt durch den realen `make test-integration`-Lauf selbst (§4): kein
  Abbruch am Anker, die Phase erschien wie erwartet.
- Gezielter `grep -n -i` gegen alle vier betroffenen Dateien nach den
  bekannten Chronik-Markern (`bislang`, `bis dahin`, `blieb bis zum
  nächsten`, `riefen bislang`, `ohne diesen nachtrag`, `nicht
  nachgetragen`): **kein Treffer** in den vier vom Fix berührten
  Stellen. Ein `bislang`-Treffer bleibt in `mapper.go:440`
  („für eine bislang nicht aktivierte Tabelle") — das ist eine
  Zustandsbeschreibung der Vorbedingung (nicht-aktivierte Tabelle), kein
  Vorher/Nachher-Marker über den Fix selbst, und lag außerhalb der von F-1
  benannten drei Stellen. Die übrigen `bislang`-Treffer in
  `run-integration-tests.sh` liegen an unveränderten, unrelated Stellen
  (NATS-Verbindung, Metriken-Dimensionen, Prozessstart-Vorlauf-Test,
  Lasttest) — nicht Teil dieses Fixes.

**Verdikt F-1:** vollständig und korrekt behoben, Anker-Mechanismus intakt,
real durch den kompletten `make test-integration`-Lauf bestätigt (die drei
Phasen erschienen unverändert korrekt, siehe §4).

## 3. F-2 (LOW) — Träger-Nachzug außerhalb des Diffs: drei von vier Stellen korrigiert, vierte bewusst optional offen

**Behauptung der Fixrunde:** `mapper.go` (`AddBinding`, `RemoveBinding`) und
`receive.go` (`Assembler()`) sind von einem SQL-Antragsqueue-Alleinstellungs-
anspruch auf „SQL-Antragsqueue oder direkter HTTP-/gRPC-Zugriffsweg"
korrigiert.

**Eigen geprüft (nicht nur gelesen):**

- `git show 61e7dcf3 -- internal/adapters/driving/replication/mapper/mapper.go`:
  `AddBinding`-Godoc trägt jetzt „… im Dauerbetrieb die
  Administrations-Goroutine (SQL-Antragsqueue) oder der direkte
  HTTP-/gRPC-Zugriffsweg …"; `RemoveBinding`-Godoc trägt „… über die
  SQL-Antragsqueue (Goroutine und Vorlauf vor dem Stream-Lauf) oder den
  direkten HTTP-/gRPC-Zugriffsweg" — beide korrekt auf Wegunabhängigkeit
  umformuliert, kein Alleinstellungsanspruch mehr.
- `git show 61e7dcf3 -- internal/adapters/driving/replication/receive/receive.go`:
  `Assembler()`-Godoc trägt jetzt „(SQL-Antragsqueue-Goroutine, ihr Vorlauf
  vor `Run`, oder der direkte HTTP-/gRPC-Zugriffsweg) trägt über ihn …" —
  ebenfalls korrigiert.
- Die im Review als „vierte, schwächere Instanz" genannte Stelle
  (`internal/bootstrap/administration_internal_test.go:220-224`,
  `fakeTableActivationPort`-Kommentar) ist **nicht** angefasst — bestätigt
  per `git show 61e7dcf3 --stat` (Datei nicht im Diff) und per Lesen der
  aktuellen Datei. Das ist konform zum Review-Verdikt: dort stand diese
  vierte Instanz ausdrücklich als „optional" neben den drei Pflicht-Stellen
  (Übergabe-Absatz: „optional die drei Fremd-Doc-Stellen in
  `mapper.go`/`receive.go`/`administration_internal_test.go`"), und F-2
  selbst war als LOW klassifiziert, nicht merge-blockierend. Die verbliebene
  Test-Kommentar-Stelle beschreibt eine Testfixtur in einem Whitebox-Test
  für den SQL-Antragsqueue-Zweig — sie behauptet keinen Alleinstellungs-
  anspruch für die Produktionsmethoden selbst, sondern beschreibt, welchen
  Zweig dieser eine Test abdeckt.

**Verdikt F-2:** die drei nicht-optionalen Stellen sind vollständig und
korrekt behoben; die vierte, explizit als optional geführte Stelle bleibt
unangetastet — kein DoD-Verstoß, da LOW und ausdrücklich optional.

## 4. Realer `make test-integration`-Lauf — eigenständig reproduziert

Vollständiger, frischer Lauf gegen den committeten Stand `61e7dcf3`
(Hintergrundprozess, Laufzeit ca. 60 Minuten inkl. Backfill-, Leerlauf- und
Transformationsphasen). Kein `FAIL`, kein unerwarteter `Error`; die beiden
einzigen Treffer für das Wort „failed"/„ERROR" im Log sind erwartete,
im Testtext selbst benannte Negative-Belege (Backfill-DDL-Fenster,
inkompatible Typänderung durch PostgreSQL selbst abgelehnt) — keine
Fehlschläge des Fixes.

Die drei neuen, vom Fix eingeführten Phasen erschienen real und exakt wie
vom Reviewer zitiert:

```
run-integration-tests: Direkter gRPC-Zugriffsweg — EnableTable(feed_e2e_api_enable_grpc) ohne SQL-Antragsqueue und ohne Neustart verarbeitet, Änderung id=1 real erfasst: ENABLED table_id=tbl-e2e-api-enable-grpc already_enabled=false
run-integration-tests: Direkter gRPC-Zugriffsweg — DisableTable(feed_e2e_api_enable_grpc) ohne SQL-Antragsqueue verarbeitet, Änderung id=2 nicht erfasst, Feed-Container läuft unverändert weiter: DISABLED removed=false retained=true
run-integration-tests: Direkter HTTP-Zugriffsweg — EnableTable(feed_e2e_api_enable_http) ohne SQL-Antragsqueue und ohne Neustart verarbeitet, Änderung id=1 real erfasst: ENABLED body={"table_id":"tbl-e2e-api-enable-http","source":"src-e2e","schema":"public","table":"feed_e2e_api_enable_http","already_enabled":false}
```

Der Lauf endete mit:

```
run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 17 Go-Zeilen und 49 Bash-Zeilen
```

— Zeile 4883 von `tools/harness/run-integration-tests.sh`, die letzte Zeile
des Skripts (`set -euo pipefail` an Zeile 62: das Skript hätte bei jedem
Fehlschlag einer nicht abgesicherten Anweisung vorzeitig abgebrochen, keine
weitere Zeile mehr gedruckt). Nach dem Lauf ist `docker ps -a --filter
name=cdc-test` leer — sauberer Teardown, konsistent mit einem vollständig
und erfolgreich durchlaufenen Skript.

**Einschränkung, offen benannt:** Der Lauf wurde über `nohup make
test-integration > log 2>&1 &` gestartet und lief vollständig detached; der
numerische Exit-Code von `make` selbst wurde dadurch nicht in einer Datei
gesichert (`$?` eines detachten Hintergrundprozesses ist über separate
Bash-Aufrufe nicht mehr abrufbar). Die Ist-Zustand-Belege dieses Abschnitts
(letzte Skriptzeile erreicht, kein `FAIL`, sauberer Teardown, `set -euo
pipefail` im Skript) sind eine starke, aber keine literal über `$?`
bestätigte Aussage — dieselbe Klasse Beleg, die der Reviewer bereits in
seinem eigenen `make test-integration`-Lauf verwendet hat.

## 5. Handbuch-Konformität

`docs/user/benutzerhandbuch.md` (Version 1.77→1.78, `a40b4809`):

- Neue Zeile in `### Änderungshistorie` (1.78), sachlich, keine
  Chronik-Sprache, kein ADR-/Review-Verweis im Fließtext.
- Neuer Hinweisabsatz in §„Zugriff über die HTTP-/JSON-API": „Tabelle
  aktivieren/deaktivieren wirkt sofort" — nennt den Vertrag präzise
  (aktualisiert den laufenden Erfassungsprozess unmittelbar, ohne Neustart
  und ohne die Antrags-Queue) und verweist auf „Tabelle live aktivieren"
  als Vergleichspunkt.
- Neuer Hinweisabsatz in §„Zugriff über die gRPC-Verwaltungs-API": „wirken
  wie ihr HTTP-Äquivalent sofort" — Präsens, Ist-Zustand, kein
  ADR-Verweis im Fließtext.
- `61e7dcf3` ändert das Handbuch nicht weiter — die Fixrunde betrifft
  ausschließlich Code-Kommentare und den Runner-Beschreibungstext, konsistent
  mit ihrer Commit-Message.

## 6. Verdikt

**DoD erfüllt: ja.**

- Der fachliche Kern des Fixes ist vollständig, korrekt und real bestätigt:
  genau 2 `AddBinding`/`RemoveBinding`-Aufrufstellen (beide in
  `assemblersync.go`), 0 neue Adapter-zu-Adapter-Kanten (`a-check`: 0
  Befunde, `git grep`: 0 Treffer), `make gates` mit Exit 0 (Coverage 80,50 %
  ≥ 80 %), und ein eigenständig gefahrener, vollständiger
  `make test-integration`-Lauf bestätigt alle drei neuen E2E-Phasen
  (gRPC-Enable, gRPC-Disable, HTTP-Enable) mit exakt den vom Reviewer
  zitierten Log-Zeilen.
- F-1 (HIGH, merge-blockierend) ist an allen drei benannten Stellen
  vollständig behoben — keine Vorher/Nachher-Chronik mehr, der
  Anker-Mechanismus für den E2E-Lauf blieb unangetastet (real durch den
  kompletten Testlauf bestätigt).
- F-2 (LOW) ist an den drei nicht-optionalen Stellen vollständig behoben;
  die vierte, im Review ausdrücklich als optional geführte Stelle bleibt
  offen — kein DoD-Bruch.
- `make kommentar-kennungen DIFF=17b5efd4` bestätigt 0 Kandidaten über den
  gesamten Drei-Commit-Diff (inkl. Fixrunde).
- Handbuch ist vollständig und korrekt nachgezogen, keine Chronik-Sprache,
  kein ADR-/Review-Verweis im Fließtext.

**Verbleibendes Restrisiko (kein DoD-Bruch, aus dem Review übernommen und
eigenständig für plausibel befunden):** Ein transientes Fehlerfenster
zwischen dem Commit des inneren `EnableTableUseCase.Enable` und den vier
Folge-Lesungen in `syncAssemblerAddBinding` bleibt bestehen — bei einem
Fehlschlag genau in diesem schmalen Fenster erhält der HTTP-/gRPC-Aufrufer
einen expliziten `500`/`Internal`-Fehler, ohne dass ein Betreiber-sichtbares
Äquivalent zum `failed`-Vermerk der SQL-Antragsqueue existiert. Dies
reproduziert einen bereits seit `ADR-0050` bestehenden Fehlermodus des
Antragsqueue-Pfads (derselbe extrahierte Code), führt keinen neuen
Fehlerpfad ein, und der Implementer hat es selbst bereits gemeldet; der
Reviewer hat es geprüft und als nicht blockierend (MEDIUM, benannte
Beobachtung statt eigenes Finding) eingestuft. Diese Einschätzung wird hier
geteilt.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md`
§3.9). `make a-check` (isoliert) — Exit 0, 0 Befunde. `make
kommentar-kennungen DIFF=17b5efd4` — Exit 0, 0 Kandidaten. `make
test-integration` — real vollständig durchgelaufen bis zur letzten
Skriptzeile („Lauf abgeschlossen"), sauberer Teardown, kein Exit-Code
literal gesichert (Einschränkung, siehe §4).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand — `61e7dcf3` — dieser
Lauf) und ersetzt keine künftige Verifikation an einem späteren Stand.
