# Review-Report: slice-060 — 2026-09-14

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-060` — drei Commits, Elter `a6f45e2` (reiner
`next→in-progress`-Move). Zwischen Elter und erstem Implementierungs-Commit
liegt `eee3909` (unabhängiger Lastenheft-Commit des Planners,
`LH-FA-SST-008`) — nicht Teil dieses Slice-Diffs, per Diff-Range
`eee3909..60afc38` ausgeklammert und real per `git diff --stat` bestätigt:

- `f5ff80e` — acht restliche Port-gedeckte Fähigkeiten über den
  HTTP/JSON-Adapter (`consumer.go`, `verwaltung.go`, `retention.go`,
  `errors.go`), Routing-Erweiterung in `server.go`
- `073c597` — `SPEC-018`-Erweiterung in `spec/pflichtenheft.md`
- `60afc38` — vier Implementierungs-DoD-Punkte abgehakt, inkl. dokumentierter
  Sechs-vs-Acht-Abweichung

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft
2026-09-13: vier repo-spezifische HIGH-Regeln plus Slice-/Wellen-Chronik-
und Handbuch-Versionshistorie-Regel)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-060-http-api-restliche-faehigkeiten.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8 Sub-Area)
  — vollständig gelesen
- `docs/plan/adr/0057-http-grpc-api.md` — bindend: Umfang (Teilfrage 2),
  Token-Authn (Teilfrage 3), Adapter-Platzierung (Teilfrage 4),
  Slice-Schnitt-Empfehlung (Konsequenzen/Folgepflicht)
- `docs/plan/planning/done/slice-059-http-api-adapter-grundgeruest-registerconsumer.md`
  als Vorbild/Baseline (§1 Out-of-Scope benennt bereits alle acht
  Fähigkeiten als `slice-060`-Scope)
- `docs/reviews/review-slice-059.md` — etablierter Qualitätsstandard dieses
  Adapters
- `docs/reviews/review-slice-052.md` — Präzedenzfall für die
  Chronik-in-Produktionscode-Regel (F-1 dort)
- `LH-FA-SST-006`, `LH-FA-CON-003/004/006`, `LH-FA-CFG-001…004`,
  `LH-FA-RET-002…004` (`spec/lastenheft.md`), `SPEC-018`
  (`spec/pflichtenheft.md`, erweitert)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.9), §5
  Traceability-Regeln
- `harness/conventions.md` (MR-000 ID-Schema)
- Reale Gate-Läufe dieses Reviews: `make gates` (baseline-verify,
  docs-check ×2, commit-traceability, coverage-gate, a-check) — alle grün,
  Exit-Code direkt und ungefiltert geprüft (siehe Verdikt)

---

## Bewertung der Sechs-vs-Acht-Diskrepanz

`slice-060`s Titel und die erste Zeile von §1 sprechen von „restlichen
**sechs** Port-gedeckten Fähigkeiten"; die Aufzählung in derselben Klammer
zählt jedoch bereits acht Elemente (Acknowledge-/Position-/Remove-Consumer =
3, Enable-/Disable-/Status-/List-Table = 4, Retention-Lauf = 1). Drei
unabhängige Quellen bestätigen, dass **acht** der korrekte Umfang ist, nicht
sechs:

- `ADR-0057` Teilfrage 2 (Option B, gewählt) und §Konsequenzen/Folgepflicht
  (Slice-Schnitt-Empfehlung „Slice B") zählen für den gesamten
  API-Erstumfang neun Fähigkeiten (`RegisterConsumer` + die acht übrigen);
  `server.go`s eigener `Config`-Kommentar zieht dieselbe Zahl („die neun
  Port-gedeckten Use Cases").
- `slice-059` §1 Out-of-Scope benannte, **bevor** `slice-060` geschrieben
  wurde, bereits exakt dieselben acht Fähigkeiten als `slice-060`s Scope.
- `slice-060` §2 DoD und §3 Plan-Tabelle (drei neue Dateien für genau diese
  drei Fähigkeitsgruppen) sind mit acht Fähigkeiten intern konsistent; die
  Implementierung deckt alle acht real ab (`consumer.go`: 3,
  `verwaltung.go`: 4, `retention.go`: 1).

Es gibt **keine** plausible Lesart, in der „sechs" die beabsichtigte Grenze
war und zwei der acht Fähigkeiten hätten zurückgestellt werden sollen — im
Gegenteil, sowohl `ADR-0057` als auch `slice-059`s eigener Out-of-Scope-
Abschnitt sahen exakt diese acht bereits vorab vor. Die Diskrepanz ist damit
ein reiner **Zähl-/Prosafehler im Slice-Titel und in §1s erster Zeile**,
keine Scope-Frage. Der Implementer hat sie selbst korrekt erkannt, in §2 DoD
transparent dokumentiert und den *tatsächlichen* (korrekten) Umfang
umgesetzt — bewertet als **LOW** (F-3 unten), keine Fixrunde nötig; der
Slice-Titel bleibt zur Closure hin unkorrigiert stehen, was für die
Nachvollziehbarkeit unschädlich ist, da §2 die Abweichung bereits erklärt.

## Findings

### F-1 — Slice-Chronik in Produktionscode-Kommentar (`writeDomainError`, `removeConsumerResponse`)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Hard Rule „Ein Kommentar beschreibt, was da
  ist") · `.harness/skills/reviewer.md` §Klassifikation (HIGH-Bullet
  „Slice-/Wellen-Chronik in Produktionscode-Kommentar") · Präzedenzfall
  `docs/reviews/review-slice-052.md` F-1 und
  `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar/`
  (4× gezählt, Architect-Verdikt vom 2026-09-13)
- `pfad`: `internal/adapters/driving/http/errors.go:13-21`
  (`writeDomainError`-Godoc); `internal/adapters/driving/http/consumer.go:103-108`
  (`removeConsumerResponse`-Godoc)
- `befund`: Beide Godoc-Kommentare stehen über **Produktionscode**-Typen/
  -Funktionen (nicht `Test*`) und begründen eine Design-Aussage mit einer
  Slice-Kennung statt mit `ADR-*`/`LH-*` oder dem Herkunfts-Anker
  `· seit slice-<NNN>`: `errors.go:15` schreibt „(`slice-060` §1 Ziel)" als
  Begründung für das einheitliche Fehler-Mapping und behauptet dabei
  zusätzlich „alle sechs Fähigkeiten dieses Slice" — faktisch acht, siehe
  Bewertung oben; `consumer.go:105-106` schreibt „(`slice-060` §2 DoD: der
  Use Case behandelt die Entfernung als Idempotenz …)" als Begründung für
  das `Removed`-Feld. Satzsubjekt ist in beiden Fällen der
  Produktionscode-Pfad, nicht ein Testfall — exakt die in
  `review-slice-052.md` F-1 etablierte Unterscheidung. Kein Gate fängt das.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Klassen
- `klasse`: „Slice-Chronik in Produktionscode-Kommentar"

### F-2 — Fehlender Negativtest: `GetStatus`-404-Pfad (`ErrSourceTableMissing`)

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation (MEDIUM-Bullet
  „fehlende Negativtests bei neuem öffentlichem Vertrag") · `SPEC-018`
  (dokumentiert `404` für `GetStatus` explizit: „physisch fehlende Tabelle
  an der Quelle → `404`")
- `pfad`: `internal/adapters/driving/http/verwaltung_test.go:165-227`
  (`fakeGetStatusUseCase`/`TestGetStatus*`)
- `befund`: `fakeGetStatusUseCase` implementiert bereits ein
  `tableExists`-Feld und den `ErrSourceTableMissing`-Fehlerpfad (Zeilen
  172-178), aber keiner der beiden `GetStatus`-Tests setzt `tableExists:
  false`, um den dokumentierten `404`-Pfad real zu belegen — anders als bei
  `EnableTable`/`DisableTable`, die je einen eigenen
  `TestXFehlendeTabelleEndetMit404`-Test tragen. Der zugrunde liegende
  Use Case (`internal/application/usecase/status/service.go:58`) gibt
  `inbound.ErrSourceTableMissing` real zurück; der Adapter-Pfad dafür ist
  ungetestet.
- `verifizierbar`: ja — ein `make test`-Lauf mit ergänztem Testfall würde es
  bestätigen (aktuell keine Aussage möglich, da der Pfad fehlt)
- `klasse`: „Fehlender Negativtest trotz vorbereiteter Test-Fake-Unterstützung"

### F-3 — Sechs-vs-Acht-Zählfehler in Slice-Titel/§1 und mehreren Testkommentaren

- `kategorie`: LOW
- `quelle`: Maintainability (Prosa-Genauigkeit, kein Konformitäts-Defekt)
- `pfad`: `docs/plan/planning/in-progress/slice-060-http-api-restliche-faehigkeiten.md:1,28`
  (Titel, §1 erste Zeile); zusätzlich dieselbe Falschzahl in
  `internal/adapters/driving/http/errors_test.go:18`,
  `internal/adapters/driving/http/consumer_test.go:278`,
  `internal/adapters/driving/http/verwaltung_test.go:74` (Test-Godocs,
  daher kein Chronik-Verstoß, siehe F-1-Abgrenzung — reine Zahlungenauigkeit)
- `befund`: Titel und §1 sprechen von „sechs" restlichen Fähigkeiten, die
  eigene Klammer-Aufzählung im selben Satz zählt acht; dieselbe Falschzahl
  wiederholt sich in drei Test-Godoc-Kommentaren. Siehe Bewertung oben:
  reiner Zähl-/Prosafehler, kein Scope-Defekt — der Implementer hat die
  Diskrepanz bereits selbst im DoD (§2) benannt und korrekt bewertet.
- `verifizierbar`: nein — kein Gate zählt Fähigkeiten in Prosa
- `klasse`: „Zählfehler in Slice-Titel/Prosa, korrekt selbst geflaggt"

### F-4 — Admin-Token-Erreichbarkeit von `GetStatus`/`ListTables` ohne eigenen Positive-Test

- `kategorie`: INFO
- `quelle`: Maintainability · `ADR-0057` Teilfrage 3 (admin deckt implizit
  die lesende Klasse ab)
- `pfad`: `internal/adapters/driving/http/verwaltung_test.go` (`GetStatus`-
  und `ListTables`-Testblöcke)
- `befund`: Für `GetConsumerPosition` existieren sowohl ein
  `ReaderToken`- als auch ein `AdminToken`-Positive-Test; für `GetStatus`
  und `ListTables` nur je ein `ReaderToken`-Test. Die
  Admin-deckt-Reader-Hierarchie ist bereits generisch durch
  `TestGetConsumerPositionAdminTokenLiefertPosition` und die
  Middleware-Tests aus `slice-059` belegt — keine erwartete Aktion an
  diesem Slice.
- `verifizierbar`: ja — `make test` (aktueller Zustand: Admin-Zugriff auf
  diese beiden Endpunkte ist ungetestet, aber durch generische
  Middleware-Logik gedeckt)
- `klasse`: „Asymmetrische Positive-Test-Abdeckung zwischen gleichrangigen Reader-Endpunkten"

## Negativbefunde

- geprüft, ohne Befund: Scope-Treue gegen `slice-060` §1 — `git diff
  eee3909 60afc38 --name-only` bestätigt: keine Änderung an
  `compose.yaml`, `harness/README.md`,
  `tools/harness/run-integration-tests.sh`, `internal/bootstrap/wiring.go`;
  keine gRPC-Artefakte; nur die acht neuen Fähigkeiten plus Routing plus
  `SPEC-018`-Erweiterung
- geprüft, ohne Befund: Fehler-Mapping (`errors.go`, `writeDomainError`) —
  `inbound.ErrSourceTableMissing` → `404`, die sechs benannten
  Domänen-Invarianten (`ADR-0029`) → `400`, jeder übrige Fehler → `500` und
  geloggt; konsistent mit `ADR-0057`s Fehler-Mapping-Vorgabe und dem in
  `slice-059` etablierten Muster (401/403 über dieselbe, unveränderte
  Middleware); `404` ist gegenüber `slice-059` neu und durch
  `inbound.ErrSourceTableMissing` sauber begründet (physisch fehlende
  Tabelle an der Quelle, `EnableTable`/`DisableTable` real getestet)
- geprüft, ohne Befund: Token-Klassen-Zuordnung je Endpunkt (`server.go:75-92`)
  — `GetConsumerPosition`, `GetStatus`, `ListTables` mit `roleReader`
  verdrahtet (erreichbar für `reader`- **und** `admin`-Token über die
  bestehende Hierarchie), alle übrigen fünf mit `roleAdmin`; deckungsgleich
  mit `ADR-0057` Teilfrage 3 und der `SPEC-018`-Tabelle; real getestet über
  `TestAcknowledgeConsumerReaderTokenEndetMit403`,
  `TestRemoveConsumerReaderTokenEndetMit403`,
  `TestEnableTableReaderTokenEndetMit403`,
  `TestDisableTableReaderTokenEndetMit403`,
  `TestRunRetentionReaderTokenEndetMit403` sowie die Reader-Positive-Tests
  für alle drei lesenden Endpunkte
- geprüft, ohne Befund: `RemoveConsumer`-Idempotenz — eine unbekannte
  Kennung liefert `removed=false`/`200`, nicht `404`
  (`TestRemoveConsumerUnbekannteKennungBleibtIdempotent`), fachlich
  gleichwertig zu CLI/SQL (`LH-FA-SST-006` Boundary); Risiko aus
  `slice-060` §6 („Idempotenz vs. `404`-Konsistenz") damit real belegt
- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-018` — per
  `awk`/`grep` real bestätigt: **keine** `ADR-*`-Referenz innerhalb des
  erweiterten Eintrags; die acht neuen Tabellenzeilen (Endpunkt, Methode,
  Rechtsklasse, Request-/Response-Schema, Statuscodes) stimmen mit der
  tatsächlichen Implementierung überein (Pfade, JSON-Feldnamen,
  Rechtsklassen je Endpunkt); `404` nur bei `EnableTable`/`DisableTable`/
  `GetStatus` dokumentiert — `ListTables`/`RunRetention`-Use-Cases geben
  `ErrSourceTableMissing` tatsächlich nicht zurück (real im
  `list`-Use-Case-Quellcode nachvollzogen), kein Doku-Code-Auseinanderlaufen
- geprüft, ohne Befund: Layering-Konformität — `consumer.go`,
  `verwaltung.go`, `retention.go`, `errors.go` importieren ausschließlich
  `application/port/inbound`, `application/port/outbound` (`LogPort`),
  `domain/model`, `domain/errors` — keine Driven-Adapter-Interna, keine
  Application-Interna, identisches Muster zu `slice-059`
- geprüft, ohne Befund: Commit-Traceability — alle drei Commit-Betreffs
  (`f5ff80e`, `073c597`, `60afc38`) tragen `LH-FA-SST-006`, `f5ff80e`
  zusätzlich `ADR-0057`; keine `SPEC-*`/`ARC-*`-Kennung in einem Betreff;
  kein `Co-Authored-By:`- oder Session-Trailer
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) im übrigen
  Diff — Zusage-/Kopplungs-/Grenz-Kommentare mit `ADR-*`/`LH-*`/`SPEC-*`-
  Bezug; keine konjunktivische Aussage über eine verworfene Alternative,
  kein abwesender Text, kein mitten im Satz abgebrochener Kommentar,
  abgesehen von den beiden unter F-1 benannten Stellen
  (`middleware.go:26`, eine dritte Slice-Referenz mit demselben Muster,
  liegt außerhalb dieses Diffs — `slice-059`-Bestand, unverändert, kein
  Finding-Kandidat dieses Reviews)
- geprüft, ohne Befund: reale Sensor-Läufe dieses Reviews — `make gates`
  ungefiltert ausgeführt: `baseline-verify` OK, `docs-check` (475 Dateien,
  0 Befunde, Vollset und Commit-Modul), `commit-traceability` (OK, 5
  Commits im Fenster, Betreffs ohne Struktur-ID), `a-check` (0 Befunde,
  unveränderter Abdeckungs-Hinweis für `tools/harness/natssub/main.go`),
  `coverage-gate` (44,40 % ≥ 35 %) — Exit-Code unmittelbar nach dem
  ungepipten `make gates`-Aufruf geprüft: `0`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Slice-Chronik in Produktionscode-Kommentar
· Fehlender Negativtest trotz vorbereiteter Test-Fake-Unterstützung ·
Zählfehler in Slice-Titel/Prosa, korrekt selbst geflaggt · Asymmetrische
Positive-Test-Abdeckung zwischen gleichrangigen Reader-Endpunkten

## Verdikt

**Merge-blockierend:** ja — ein HIGH-Finding (F-1) und ein MEDIUM-Finding
(F-2) blockieren die Übergabe an den Verifier.

**Übergabe:** F-1 und F-2 gehen als Findings an den Implementer zurück
(Rückkante Review → Implementer, Modul 8) für eine Fixrunde — beide sind
lokal begrenzte Korrekturen (Kommentartext ohne Verhaltensänderung bei F-1;
ein zusätzlicher Testfall ohne Produktionscode-Änderung bei F-2) ohne
Rollen-Widerspruch und ohne drittes Auftreten eines neuen Konflikttyps im
Sinne von Modul 8 §Konflikt-Pfad — keine Architect-Beteiligung nötig. F-3
und F-4 gehen informativ mit, lösen aber keine Fixrunde aus. Da eine
Fixrunde folgt, bleibt die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan **offen** (kein Nachzug in diesem
Commit, siehe `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde, Grenzfall „bleibt die Checkbox offen"). Die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler. Dieser Report selbst ist ein **Lauf-Beleg**
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und wird
über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).

---

## Fixrunden-Nachtrag — 2026-09-14

Drei Fix-Commits real geprüft (`git show`, nicht nur Commit-Messages):
`bf38260` (F-1), `5b2921c` (F-2 + F-4), `2364860` (F-3).

### F-1 — behoben

`git show bf38260` bestätigt: `errors.go:13-19` (`writeDomainError`-Godoc)
und `consumer.go:104-107` (`removeConsumerResponse`-Godoc) tragen keine
`slice-060`-Referenz mehr; beide beschreiben den Ist-Zustand ausschließlich
über `ADR-*`/`LH-*`/`SPEC-*`-Bezüge. `errors.go` verliert zugleich die
falsche „sechs"-Zählung aus demselben Satz. `grep -n "slice-0"
internal/adapters/driving/http/errors.go internal/adapters/driving/http/consumer.go`
liefert real keinen Treffer mehr.

### F-2 — behoben

`git show 5b2921c` bestätigt: `TestGetStatusFehlendeTabelleEndetMit404`
setzt `fakeGetStatusUseCase{tableExists: false}` und erwartet real
`http.StatusNotFound` — übt damit den zuvor ungetesteten
`ErrSourceTableMissing`→`404`-Pfad tatsächlich aus, nicht nur namentlich.
`make gates` (Coverage-Stufe, s. u.) lief mit diesem neuen Test grün durch.

### F-3 — behoben, mit Selbstkorrektur meines eigenen Erstlauf-Befunds

`git show 2364860` und `grep -rn "sechs" internal/adapters/driving/http/
docs/plan/planning/in-progress/slice-060*.md` (selbst ausgeführt): **kein
Treffer mehr** — Titel-Zeile in §1, DoD, Plan-Tabelle (§3), §4 und §6 des
Slice-Plans sowie der Test-Godoc in `errors_test.go` sind durchgängig auf
„acht" korrigiert; die DoD-Abweichungsnotiz ist konsequent entfernt, da sie
gegenstandslos wurde.

Bei der Verifikation fällt eine Ungenauigkeit meines **eigenen**
Erstlauf-Befunds auf, die hier festgehalten wird statt stillschweigend
übernommen zu werden: F-3 hatte `consumer_test.go:278` und
`verwaltung_test.go:74` als weitere Fundstellen „derselben Falschzahl"
gelistet. Ein realer Blick auf den dortigen Text (bereits im Erstlauf-Stand
`60afc38`) zeigt: Diese beiden Stellen enthielten nie eine Zahl „sechs" —
sie zitieren `slice-060 §2 DoD` als zulässige Testfall-Provenienz (Subjekt
ist der Testfall, kein Zählfehler). Der Implementer hat sie im Fix-Commit
deshalb zu Recht **nicht** angefasst; mein Erstlauf-`pfad`-Feld war an
diesen zwei Stellen überzogen. Der tatsächliche Kern von F-3 (Titel-Zeile
§1, DoD, §3, §4, §6, `errors_test.go`) ist vollständig und korrekt behoben.

### F-4 — behoben

`git show 5b2921c` bestätigt zwei neue Tests:
`TestGetStatusAdminTokenLiefertStatus` und
`TestListTablesAdminTokenLiefertListe`, beide mit `testAdminToken` gegen
den jeweiligen lesenden Endpunkt und realer Prüfung auf
`http.StatusOK`. Schließt die zuvor benannte Asymmetrie zwischen den drei
lesenden Endpunkten sinnvoll.

### Gate-Lauf nach der Fixrunde

`make gates` erneut selbst ausgeführt (ungepiped, Exit-Code direkt
geprüft): **Exit-Code `0`** — `baseline-verify`, `docs-check` (476 Dateien,
0 Befunde), `commit-traceability` (5 Commits im Fenster, keine
Struktur-ID im Betreff), `coverage-gate` (44,40 % ≥ 35 %, neue Tests liefen
grün durch), `a-check` (0 Befunde) — alle grün.

### Aktualisiertes Verdikt

Alle vier Findings sind real geprüft und behoben. Kein weiterer
Fixrunden-Bedarf. Die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" ist im Slice-Plan auf `[x]` nachgezogen (dieser
Commit), mit Verweis auf diesen Report.
