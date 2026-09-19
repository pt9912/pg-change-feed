# Review-Report: rtm-letzte-zwoelf-tag-only — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/rtm-letzte-zwoelf-tag-only.md`) und
`AGENTS.md` Hard Rules (nicht gegen die DoD — das ist Verifier-Aufgabe,
Modul 11).

**Gegenstand:** Commit `748ac7bc` (Parent `a88ebbf4`) — reiner
Tag-Nachtrag ohne neue Testfälle/Prüfzeilen an sechs bestehenden Stellen
(`test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh`,
Erzeugnis `docs/user/e2e-abdeckung.md`, `harness/README.md`).

**Skill:** `.harness/skills/reviewer.md` @ `70098d3` · **Modell:**
claude-sonnet-5 · **Datum:** 2026-09-19

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/rtm-letzte-zwoelf-tag-only.md` (voll gelesen)
- `spec/lastenheft.md` §`LH-FA-CON-001`…`005`, §`LH-FA-DAT-002`/`003`/`005`,
  §`LH-FA-REA-001`, §`LH-FA-RET-001`, §`LH-QA-REL-003`/`004` (Volltext
  gelesen, inkl. Akzeptanzkriterien/Messmethoden)
- `AGENTS.md` §3 (Hard Rules), insbes. §3.7, §3.12, §3.13
- `.harness/skills/reviewer.md` (Klassifikation, HIGH-Liste)
- `git show 748ac7bc` — voller Diff aller fünf geänderten Dateien einzeln gelesen
- Tatsächlicher Testcode an jeder der sechs getaggten Stellen gelesen
  (nicht nur die neuen Kommentare/`abdeckung_declare`-Strings):
  `TestE2ECaptureFlow` (inkl. Wiederlese-Block),
  `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer` (vollständig),
  „Black-Box-CLI-Rundlauf" (Zeilen ~811–1150), „HTTP-API-Rundlauf" (Zeilen
  ~1987–2100, inkl. `tools/harness/httpclient/main.go`), „CLI-Diagnose-Beleg
  (Fehlerzustand)" (Zeilen ~1281–1320), „Upgrade-Sicherheits-Rundlauf"
  (Zeilen ~2687–2790)
- `make doc-trace` real ausgeführt (netzlos) zur Verifikation der
  RTM-Zahlen
- `make docs-check` real ausgeführt (netzlos, 751 Dateien)
- `make commit-traceability` real ausgeführt (netzlos, `RANGE=HEAD~5..HEAD`)

---

## Zentrale Prüfung: inhaltliche Rechtfertigung der zwölf Tags

Je Tag wurde der tatsächliche Testcode gegen Beschreibung/Akzeptanzkriterien
des Lastenhefts gehalten (nicht nur der neue Kommentartext). Ergebnis:

| Tag | Fundstelle | Befund |
|---|---|---|
| `LH-FA-CON-001` (Registrierung) | Black-Box-CLI-Rundlauf | Belegt real: `register-consumer` per `docker exec`, danach `SELECT count(*) FROM cdc.consumer WHERE consumer_id = …` = 1. Happy Path exakt getroffen. |
| `LH-FA-CON-003` (Positions-Persistierung) | Black-Box-CLI-Rundlauf | Belegt real: nach `acknowledge-consumer` wird `cdc.consumer_position.acknowledged_position` zurückgelesen und exakt gegen den bestätigten Wert geprüft (`acked_first != first_position` bricht ab). Happy Path exakt getroffen. |
| `LH-FA-CON-004` (Bestätigung) | Black-Box-CLI-Rundlauf | Belegt real: `acknowledge-consumer` erfolgreich, Rücklesung bestätigt. Happy Path exakt getroffen. |
| `LH-FA-CON-005` (Fortsetzung nach Neustart) | Black-Box-CLI-Rundlauf | Belegt real: echter `docker restart`, danach `restored_position` unabhängig aus `cdc.consumer_position` gelesen (`== first_position`), und die ID-Folge ab dieser Position (`resumed_ids`) enthält exakt die neue Zeile (id=96), nicht die bereits bestätigte (id=95). Happy Path exakt getroffen — stärker als nötig (echter Prozess-Neustart, kein simulierter Zustand). |
| `LH-FA-DAT-002` (Quelltabelle identifizierbar) | `TestE2ECaptureFlow` | Belegt real: `record.Change.SourceTableID != env.tableID` wird für alle drei Changes geprüft. Happy Path getroffen. |
| `LH-FA-DAT-003` (Operationstyp) | `TestE2ECaptureFlow` | Belegt real: `operations := []model.Operation{Insert, Update, Delete}` gegen `record.Change.Operation` je Change geprüft. Happy Path getroffen. |
| `LH-FA-DAT-005` (Relevante Datenwerte) | `TestE2ECaptureFlow` | Belegt real: `insertImage`/`updateImage`/`deleteImage` mit tatsächlichen Spaltenwerten (`id`, `name`) je Operationstyp geprüft, inkl. korrekter Ab-/Anwesenheit von Alt-/Neu-Image. Happy Path getroffen. |
| `LH-QA-REL-004` (Idempotente Verarbeitung) | `TestE2ECaptureFlow`, Wiederlese-Block | Belegt real und **exakt nach Messmethode**: zweites `ReadChanges` desselben Bereichs, Länge und Inhalt (`Change.ID`, `Position.Offset`, `NewImage`) gegen den Erstsatz geprüft. Trifft die Lastenheft-Messmethode wörtlich („doppeltes Lesen desselben Bereichs liefert dieselben Changes"). |
| `LH-FA-REA-001` (Lesebereich zwischen zwei Positionen) | HTTP-API-Rundlauf | Belegt real: `http_read_to = http_read_position + 1` wird als eigener `to`-Query-Parameter neben `from` an `GET /changes` übergeben (`tools/harness/httpclient/main.go`, `readChanges()`); die Antwort trägt genau die eigens eingefügte Zeile. Ein echter Zwei-Positionen-Bereich, kein unbegrenzter Lesezugriff — die im Kontext gestellte Prüffrage ist mit Nein zu beantworten (es ist **kein** unbegrenzter Read). |
| `LH-FA-RET-001` (Persistente Speicherung über Neustart) | Upgrade-Sicherheits-Rundlauf | Belegt real und **stärker als die Wortwahl im Kommentar nahelegt**: `upgrade_before_still` prüft explizit, dass die **vor** dem Container-Tausch erfasste Zeile (id=250) **danach** noch mit `count(*) = 1` über `cdc.changes` lesbar ist — nicht nur, dass die Erfassung danach weiterläuft (das wird zusätzlich, separat, über id=251 geprüft). Die im Kontext gestellte Prüffrage ist mit Ja zu beantworten. |
| `LH-QA-REL-003` (Sichtbarer Unzuverlässigkeitszustand) | CLI-Diagnose-Beleg (Fehlerzustand) | Belegt real und passt zur Messmethode wörtlich: Fehlerzustand wird direkt in `cdc.process_heartbeat.error_class` injiziert, `diagnose`-Ausgabe muss „Fehlerzustand (`LH-FA-ADM-003`): schema" zeigen **und** darf nicht „keiner (Normalbetrieb)" zeigen — exakt „Fehlerinjektionstest: Zustand herbeiführen und Erkennbarkeit über `LH-FA-ADM-003` prüfen". |

Elf der zwölf Tags sind damit **inhaltlich vollständig gerechtfertigt** —
kein kosmetischer Nachtrag, sondern reale, bereits vor diesem Slice grün
laufende Prüfzeilen, die exakt das behaupten, was der Kommentar/die
`abdeckung_declare`-Beschreibung sagt. Zum zwölften Tag (`LH-FA-CON-002`)
siehe F-1 unten — dort ist die Substanz real vorhanden, aber schmaler als
die im Kommentar gewählte Formulierung.

---

## Findings

### F-1 — `LH-FA-CON-002`-Tag: reale Teilevidenz, aber weder Happy-Path- noch Boundary-GWT wörtlich getroffen

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Reviewer-Skill „Beleg trägt seinen Satz
  nicht" (HIGH-Klasse), hier mit begründeter Abstufung — der Beleg trägt
  die Aussage **teilweise**, nicht **gar nicht**
- `pfad`: `test/integration/integration_test.go:503-505` (neuer
  Kommentarsatz), Testkörper `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer`
  (`test/integration/integration_test.go:511-619`)
- `befund`: Der neue Satz „beide Consumer bestätigen unabhängig
  voneinander, die Position des einen bleibt von der Bestätigung des
  anderen unberührt" wird nur in **einer** Richtung real geprüft: nachdem
  `behindConsumer` seine Position bestätigt hat, bestätigt
  `aheadConsumer` eine **andere** Position (nicht denselben Bereich, wie
  es `LH-FA-CON-002`s Happy-Path-GWT verlangt: „beide denselben
  Änderungsbereich lesen und bestätigen"), und danach wird geprüft, dass
  `behindConsumer`s bereits gespeicherte Position exakt erhalten blieb
  (`blockers[0].position == behindRows[0].commitPosition`). Das ist reale,
  nicht-kosmetische Evidenz gegen ein Überschreiben durch eine spätere
  fremde Bestätigung — aber weder der Happy-Path (dieselbe Bestätigung
  desselben Bereichs durch beide) noch der Boundary-Pfad („Given `c1`
  bestätigt Position `p`, when `c2` anschließend **liest** [nicht:
  bestätigt], then bleibt die gelesene Position von `c2` unverändert")
  werden wörtlich exerciert: `aheadConsumer`s eigene gespeicherte Position
  wird nie unabhängig vor ihrer eigenen Bestätigung gelesen und geprüft,
  und beide Consumer wirken nie auf demselben Bereich. Die Prüffrage aus
  dem Auftrag „nicht nur zwei Consumer nacheinander" ist damit nicht
  eindeutig mit Nein zu beantworten: es sind real zwei Consumer, die real
  nacheinander (nicht gleichzeitig, nicht auf demselben Bereich)
  bestätigen — die Nichtbeeinflussung ist an genau einer Stelle bewiesen
  (later-write-doesn't-clobber-earlier), nicht umfassend.
- `verifizierbar`: ja — Quelltext-Lektüre der zitierten Zeilen; keine
  Ausführung nötig, die Lücke ist strukturell im Testablauf sichtbar.
- `klasse`: „GWT-Szenario des Lastenhefts nicht wörtlich exerciert, obwohl
  Kommentar es so formuliert"

### F-2 — Nachbar-Inline-Kommentar zum Wiederlese-Block trägt die neu ergänzte Kennung nicht mit

- `kategorie`: INFO
- `quelle`: Maintainability (Nachbar-Frage zu `AGENTS.md` §3.13) — bereits
  einmal als eigenes Finding im Vorgänger-Review aufgetreten (F-5,
  `review-slice-rtm-reste-sst-cfg-por.md`), hier zweites Auftreten,
  weiterhin unter der 3×-Schärfungsschwelle
- `pfad`: `test/integration/integration_test.go:256-257` (bestehender
  Inline-Kommentar direkt über dem Wiederlese-Block: „Das Wiederlesen
  desselben Bereichs trägt dieselben Changes in derselben Reihenfolge
  (`LH-FA-REA-004.a`, `LH-FA-REA-005`)."), verglichen mit dem neuen
  Funktions-Kopfkommentar (Zeile 197-198), der `LH-QA-REL-004` dort
  bereits nennt
- `befund`: Die neu ergänzte Kennung `LH-QA-REL-004` steht ausschließlich
  im Funktions-Kopfkommentar, nicht im direkt benachbarten,
  block-lokalen Inline-Kommentar, der von diesem Diff nicht angefasst
  wurde. Kein Widerspruch, nur eine unvollständige Querverweisung an der
  Stelle, die ein Lesender zuerst sieht.
- `verifizierbar`: ja — Diff-/Quelltext-Lektüre.
- `klasse`: —

---

## Positiv-Befunde (was real geprüft wurde und trägt)

- **RTM-Zahlen real nachgemessen, keine Drift:** `make doc-trace` netzlos
  ausgeführt liefert exakt `76 Anforderung(en), 0 Waise(n)` —
  deckungsgleich mit der Commit-Message („make doc-trace: 76
  Anforderungen, 0 Waisen") und mit `harness/README.md`s neuer
  `make doc-trace`-Zeile („76 Anforderungen, 55 Waisen ohne
  `trace.coverage`, **0 Waisen mit**"). Kein Fall der Klasse „Zahl im
  Träger driftet gegen die Messung" (`AGENTS.md` §3.12).
- **Elf von zwölf Tags inhaltlich vollständig gerechtfertigt** (siehe
  Tabelle oben) — jeweils Happy-Path-GWT oder Messmethode wörtlich
  getroffen, an bereits bestehendem, unverändertem Testcode. Kein
  einziger der elf ist eine Behauptung ohne Substanz.
- **`LH-FA-RET-001`-Beleg tatsächlich stärker als nötig:** Der reale
  Container-**Tausch** (`--force-recreate`, neue Container-ID bestätigt)
  ist eine strengere Form von „Neustart" als ein bloßer
  Prozess-Restart — die Feed-Anwendung verliert dabei jeden In-Memory-
  Zustand vollständig, während PostgreSQL/NATS unberührt bleiben
  (`--no-deps`, Container-IDs beider unverändert geprüft).
- **`docs/user/e2e-abdeckung.md` korrekt regeneriert:** Alle
  verschobenen Zeilen-Lokatoren im Diff (z. B. `integration_test.go:197`
  → `:202`, `:305` → `:310`) wurden gegen den aktuellen Quellcode
  nachgeprüft und stimmen exakt (`grep -n "^func TestE2ECaptureFlow"` →
  Zeile 202; `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer` →
  Zeile 511) — kein Fall von „Beleg trägt seinen Satz nicht" durch
  veraltete Zeilenangaben.
- **`AGENTS.md` §3.7 eingehalten:** Alle neuen/geänderten Kommentarstellen
  (Godoc über zwei `Test*`-Funktionen, vier `abdeckung_declare`-Strings)
  beschreiben den gegenwärtigen Zustand im Indikativ, keine
  Vorher/Nachher-Sprache, kein Abbruch mitten im Satz, kein Bezug auf eine
  verworfene Alternative.
- **Keine Slice-/Wellen-Chronik in Produktionscode:** Die betroffenen
  Kommentare stehen ausschließlich über `Test*`-Funktionen (zulässige
  Testfall-Provenienz-Form) bzw. in einem Shell-Testharness-String, nicht
  über Produktionscode-Pfaden.
- **`make docs-check` real ausgeführt:** 751 Dateien, 0 Befunde — keine
  kaputten Links/Anker/hostpaths-Treffer durch die geänderten Dateien.
- **`make commit-traceability` real ausgeführt:** `HEAD~5..HEAD`, 0
  Befunde — Commit `748ac7bc` trägt alle zwölf `LH-*`-Kennungen im Body,
  kein `SPEC-*`/`ARC-*` im Betreff.
- **Spec-Stratum unangetastet:** `spec/lastenheft.md`,
  `spec/pflichtenheft.md`, `spec/architecture.md` in diesem Diff nicht
  berührt (laut `git show --stat 748ac7bc`) — reiner Tag-Nachtrag ohne
  neue/geänderte Anforderung.
- **Keine Betreiber-Oberfläche berührt:** keine neue `CDC_*`-Variable,
  keine `cdc.*`-SQL-Funktion, kein Endpunkt; `docs/user/benutzerhandbuch.md`
  zu Recht nicht angefasst.
- **Ehrliche Selbstbegrenzung im Slice-Plan:** §6 benennt bereits vor
  diesem Review offen, dass Boundary-/Negative-Akzeptanzkriterien (z. B.
  `LH-FA-CON-001`s „erneute Registrierung idempotent", `LH-FA-DAT-005`s
  „nicht lieferbarer Wert", `LH-FA-DAT-002`s „zwei gleichnamige Tabellen
  in verschiedenen Schemata") ohne eigenen Testbeleg bleiben — das deckt
  sich mit der eigenen Nachprüfung dieses Reviews, keine versteckte
  Lücke.

## Negativbefunde

- geprüft, ohne Befund: `spec/lastenheft.md`, `spec/pflichtenheft.md`,
  `spec/architecture.md` — unangetastet, Spec-Stratum-Grenze intakt
- geprüft, ohne Befund: `.github/workflows/*.yml` — unangetastet,
  `AGENTS.md` §3.8 (Action-Pinning) nicht berührt
- geprüft, ohne Befund: Docker-only-Disziplin (§3.1) — keine lokale
  Toolchain-Installation im Diff
- geprüft, ohne Befund: Suppression-Verbot (§3.2) — kein `//nolint` o. ä.
- geprüft, ohne Befund: `git mv`-Disziplin (§3.3) — keine
  Datei-Bewegungen in diesem Diff
- geprüft, ohne Befund: ADR-Immutabilität (§3.5) — dieser Diff berührt
  keine ADR-Datei
- geprüft, ohne Befund: host-lokale absolute Pfade (§3.11) —
  `docs-check`-Modul `hostpaths` bestätigt 0 Befunde repo-weit
- geprüft, ohne Befund: `docs/plan/adr/README.md`, `.d-check.yml` —
  unangetastet, keine neue ADR/Coverage-Dimension in diesem Slice nötig
- geprüft, ohne Befund: Referenzen auf die alte Waisenzahl (`12 Waisen`)
  oder auf „ohne eigenen Kennungs-Tag" bei den zwölf Anforderungen in
  anderen Markdown-Dateien — repo-weiter `grep` findet keinen stehen
  gebliebenen Verweis (`AGENTS.md` §3.13)
- nicht ausgeführt: `make test-integration` (voller Docker-/DB-Stack,
  außerhalb des Zeitbudgets dieses Reviews) — die Regenerierung von
  `docs/user/e2e-abdeckung.md` wurde stattdessen strukturell gegen den
  aktuellen Quellcode geprüft (Zeilen-Lokatoren, s. o.), nicht durch einen
  eigenen Lauf reproduziert
- nicht ausgeführt: `make gates` vollständig (u. a. `coverage-gate`,
  `a-check`, `baseline-verify`, `generated-sync`) — dieser Diff ändert
  ausschließlich Kommentare/String-Literale in einer Testdatei und einem
  Shell-Testharness sowie zwei generierte/dokumentarische Markdown-Dateien,
  keinen Go-Produktionscode-Pfad und keinen generierten Code;
  `docs-check`/`commit-traceability`/`doc-trace` gezielt separat
  ausgeführt (s. o.)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** GWT-Szenario des Lastenhefts nicht
wörtlich exerciert, obwohl Kommentar es so formuliert.

## Verdikt

**Merge-blockierend:** ja — ein MEDIUM-Finding (F-1). Die Abweichung von
einer automatischen HIGH-Einstufung wird hier begründet: der genannte
Beleg trägt die Aussage teilweise (reale Evidenz gegen Überschreiben
durch eine spätere fremde Bestätigung), nicht gar nicht — anders als ein
Beleg, der etwas völlig anderes misst.

**Übergabe:** F-1 (MEDIUM) geht an den Implementer zurück. Fixrunde
empfohlen, aber mechanisch günstig zu beheben **ohne neue Testlogik**
(deckt sich mit dem im Slice-Plan §"Ausdrücklich NICHT in diesem Slice"
gesetzten Rahmen): entweder die Kommentarformulierung auf das real
Geprüfte präzisieren (z. B. „eine spätere Bestätigung des einen
überschreibt nicht die bereits gespeicherte Position des anderen" statt
der umfassenderen „beide Consumer bestätigen unabhängig voneinander … die
Position des einen bleibt von der Bestätigung des anderen unberührt"),
oder die Lücke zum Happy-Path-/Boundary-GWT explizit als akzeptierten
Teil-Beleg im Slice-Plan §6 benennen (analog zu den dort bereits
gelisteten Boundary-Lücken). F-2 (INFO) ist ein Hinweis ohne zwingende
Aktion. Da dieses Verdikt eine Fixrunde vorsieht, wird die DoD-Checkbox
„Review durchgeführt" in
`docs/plan/planning/in-progress/rtm-letzte-zwoelf-tag-only.md` **nicht**
in diesem Report nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne
Fixrunde — Grenzfall trifft hier nicht zu) — der reguläre Nachzug läuft
nach der Fixrunde am Implementer-Workflow-Schritt 21.

Dieser Report ist ein Lauf-Beleg (Audit: dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/
Spec-Konformität prüft der Verifier separat (Modul 11).
