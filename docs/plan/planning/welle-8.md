# Welle 8: Black-Box-E2E und Integrationstest-Nachzug

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-8-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-12.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Keiner der drei bestehenden `*_endtoend_test.go`-Tests
(`test/integration/integration_test.go`, `internal/bootstrap/welle6_endtoend_test.go`,
`internal/bootstrap/walretention_endtoend_test.go`) behandelt das System als
Black Box — alle drei importieren interne Go-Pakete
(`postgresstorage`/`list`/`status`/`disable`/`bootstrap.Run`/
`bootstrap.RegisterConsumer`/`bootstrap.AcknowledgeConsumer`) und rufen sie
im selben Testprozess auf. Nach `ADR-0030`s Testpyramide ist das der
**Integrationstest**-Tier, nicht der getrennt geführte **E2E**-Tier — das
Repo hat aktuell **keinen einzigen** echten Black-Box-E2E-Test, obwohl die
externe Schnittstelle dafür bereits existiert:
[`cmd/pg-change-feed/main.go`](../../../cmd/pg-change-feed/main.go) hat
reale CLI-Unterbefehle (`register-consumer`, `acknowledge-consumer`), die im
Container laufen und bislang von keinem Test als Subprozess/`docker exec`
aufgerufen werden.

Zusätzlich ist `make test-integration`
([`tools/harness/run-integration-tests.sh`](../../../tools/harness/run-integration-tests.sh))
seit `welle-3` inhaltlich nicht über den ursprünglichen MVP-Zuschnitt
hinausgewachsen (`spec/lastenheft.md` §1 MVP-Schnitt) — mit einer Ausnahme
(der `cdc_capture_lag`-Lasttest-Beleg aus `welle-5`, selbst außerhalb des
MVP-Schnitts). Consumer-Zugriffsweg (`welle-6`), rollen-spezifische
DSN-Trennung (`slice-023`, `ADR-0047`) und WAL-Rückstand-Schwellen
(`welle-7`) werden nicht gegen den vollen, containerisierten Compose-Stack
geprüft — nur isoliert über `make test-store`/`make test-replication`.
`BEO-PGC/rollen-test-abdeckungsluecken` (1×) benennt dabei konkret: die
Replication-Stream-/ACK-Adapter sind gegen Rollen-Vertauschung nicht
testgesichert, weil `tools/harness/run-replication-tests.sh` keine
rollenbeschränkten Login-Test-Identitäten bereitstellt.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs:** Kein Slice-DoD allein
beweist, dass (a) ein echter externer Aufrufer — nicht ein Go-Test-Prozess —
das System über die tatsächlich ausgelieferte Schnittstelle bedienen kann,
und (b) der Compose-Integrationstest inzwischen alle seit `welle-3`
hinzugekommenen post-MVP-Fähigkeiten zusammen abdeckt, nicht nur den
ursprünglichen MVP-Rundlauf.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `welle-7` liegt in `done/`.
- Kein Slice liegt in `in-progress/`.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices (`slice-027`, `slice-028`) liegen in `done/`.
- `make gates` grün.
- Der neue Black-Box-E2E-Test läuft real gegen den Compose-Stack, ruft
  ausschließlich die CLI als externen Prozess auf (kein Go-Paket-Import
  interner Anwendungslogik) und deckt den vollen Consumer-Rundlauf
  (registrieren → lesen → bestätigen → Neustart simulieren → fortsetzen) ab.
- Der Compose-Integrationstest deckt nachweislich alle drei seit `welle-3`
  hinzugekommenen post-MVP-Fähigkeiten zusammen ab (Consumer-Zugriffsweg,
  rollen-spezifische DSN-Trennung, Nennung des WAL-Rückstand-Verhaltens als
  bewusst andernorts abgedeckt — siehe §6).
- Closure-Notiz in `welle-8-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-027 | Black-Box-E2E-Test über die bestehende CLI | [`LH-QA-POR-003`](../../../spec/lastenheft.md) |
| slice-028 | Compose-Integrationstest-Nachzug — Rollen-DSN-Trennung, MVP-Umbenennung | [`LH-QA-SEC-001`](../../../spec/lastenheft.md)…`003` |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keine andere Welle.
- Intern (Slice-Reihenfolge innerhalb dieser Welle): `slice-028` erweitert
  denselben Compose-Integrationstest, den `slice-027` um den ersten
  Black-Box-Testfall ergänzt — `slice-027` zuerst, damit `slice-028` auf
  demselben externen Aufrufmuster (CLI-Subprozess) aufbaut statt ein
  zweites zu erfinden.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **WAL-Rückstand-Schwellen-Verhalten wird nicht zusätzlich gegen den vollen
  Compose-Stack geprüft** — Bestand bleibt bewusst stehen:
  `internal/bootstrap/walretention_endtoend_test.go`
  (`TestWALRetentionThresholdEndToEnd`, `slice-026`) beweist dasselbe
  Verhalten bereits real gegen PostgreSQL über denselben Code-Pfad
  (`bootstrap.Run`), den auch der containerisierte Prozess durchläuft; eine
  zweite, langsamere Prüfung auf Compose-Ebene wäre Redundanz ohne
  zusätzliche Aussagekraft, kein Abdeckungsgewinn.
- **`LH-FA-SST-006`/`LH-FA-SST-007`** (HTTP/gRPC-API, NATS-Change-
  Notification) bleiben unberührt — Schicht-Abgrenzung: Diese Welle testet
  bestehende Schnittstellen, baut keine neuen.
- **`BEO-PGC/test-isolation-geteilter-zustand`** (geteilter Testzustand in
  `postgresstorage`-Paket-Tests) bleibt unberührt — anderer Vorgang: Das ist
  eine Ursachenbehebung an bestehenden Unit-/Adapter-Tests, keine neue
  Compose-/Black-Box-Testabdeckung.
- **`BEO-PGC/walsender-wirksamkeit`** (Publication-Entzug-Timing) bleibt
  unberührt — bereits als eigene, spätere Welle vorgemerkt (Roadmap
  *Nächste Wellen*), kein Teil der Testinfrastruktur-Arbeit hier.

## 7. Closure-Notiz

<!--
BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md §Verwendung,
Schritt 5) und darf deshalb nichts Tragendes halten.

- Erst nach Welle-Abschluss fuellen; nur die Nummer, nicht die volle Welle-ID.
- Ziel-Form der Ergebnis-Notiz: `welle-results.template.md` — Schwester-Vorlage
  im Template-Verzeichnis, kein Artefakt deines Repos. Sie ist von der
  Ruheort-Regel ausgenommen und faellt mit diesem Kommentar ohnehin weg.
-->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
