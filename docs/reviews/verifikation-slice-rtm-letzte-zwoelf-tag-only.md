# Verifikationsbericht: slice-rtm-letzte-zwoelf-tag-only — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/rtm-letzte-zwoelf-tag-only.md` §2) und
die §6-Risiko-Ausgänge. **Nicht** gegen den Diff als solchen (Reviewer-
Aufgabe, mit
[`review-slice-rtm-letzte-zwoelf-tag-only.md`](review-slice-rtm-letzte-zwoelf-tag-only.md)
abgeschlossen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Gegenstand:** zwei Commits auf `main`:

- `748ac7bc` — ursprünglicher Implementer-Commit (zwölf Tag-Nachträge,
  kein neuer Testcode).
- `d1136723` — Fixrunde nach einem MEDIUM-Finding (F-1) des unabhängigen
  Reviewers: Kommentarpräzisierung an
  `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer`
  (`LH-FA-CON-002`), Slice-Plan §6 um den Review-Fund ergänzt, DoD-Checkbox
  „Review durchgeführt" nachgezogen.

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan, den Review-Report,
beide Commit-Diffs, `spec/lastenheft.md` (Volltext der zwölf betroffenen
Anforderungen) und den tatsächlichen Testcode an jeder der sechs
getaggten Stellen gelesen. Nichts aus Slice-Plan, Commit-Message oder
Review-Report wurde übernommen, ohne es selbst nachzuprüfen: eigener
`make doc-trace`-Lauf, eigene Codelektüre an allen sechs Fundstellen
gegen den jeweiligen Lastenheft-Wortlaut, eigener `grep` gegen stehen
gebliebene Referenzen auf die alte Waisenzahl, eigener Zeilen-Abgleich
der generierten `docs/user/e2e-abdeckung.md` gegen den aktuellen
Quellcode, eigener, ungepipter `make gates`-Lauf mit direkter
Exit-Code-Prüfung (`AGENTS.md` §3.9).

**Working-Tree-Hinweis (methodisch relevant für §4):** Im
Haupt-Arbeitsbaum lief zum Zeitpunkt dieser Verifikation parallel ein
anderer Verifier-Agent an einem anderen, bereits committeten Slice
(`rtm-reste-sst-cfg-por`, Commits `a88ebbf4`/`b3e1c100`). Ein erster
`make gates`-Lauf im Haupt-Arbeitsbaum scheiterte an `docs-check`
(`hostpath-forbidden` in dessen eigener, damals unfertiger
Report-Datei) — nicht an einer Datei dieses Slice-Diffs. Diese Datei
wurde **nicht** angefasst oder verschoben (Auftrag: ignorieren). Um
`make gates` unabhängig von dieser Fremdinterferenz auf dem tatsächlichen
`HEAD` (`d1136723`) zu prüfen, wurde ein vollständiger, eigenständiger
`git clone` (kein `git worktree` — dessen `.git`-Datei-Indirektion bricht
den `docker run -v`-Mount des Sensors) in einem temporären
Geschwister-Verzeichnis außerhalb dieses Repos angelegt, auf `d1136723`
ausgecheckt, dort `make gates` gefahren und danach wieder entfernt. Der
Haupt-Arbeitsbaum wurde durch diese Verifikation nicht verändert (siehe
§4).

---

## 1. `make doc-trace` real ausgeführt — 76/0 gegengeprüft

Eigener Lauf:

```
make doc-trace
```

Letzte Zeile der Ausgabe: **`76 Anforderung(en), 0 Waise(n).`** Zusätzlich
gezielt gegen alle zwölf Ziel-Kennungen gefiltert — jede zeigt `ok`:

| Kennung | Träger (`make doc-trace`) | Ergebnis |
|---|---|---|
| `LH-FA-CON-001` | E2E | `ok` |
| `LH-FA-CON-002` | E2E | `ok` |
| `LH-FA-CON-003` | E2E | `ok` |
| `LH-FA-CON-004` | E2E | `ok` |
| `LH-FA-CON-005` | E2E | `ok` |
| `LH-FA-DAT-002` | E2E | `ok` |
| `LH-FA-DAT-003` | E2E | `ok` |
| `LH-FA-DAT-005` | E2E | `ok` |
| `LH-FA-REA-001` | E2E | `ok` |
| `LH-FA-RET-001` | E2E | `ok` |
| `LH-QA-REL-003` | E2E | `ok` |
| `LH-QA-REL-004` | E2E | `ok` |

**Ergebnis: deckungsgleich** mit der Commit-Message („make doc-trace: 76
Anforderungen, 0 Waisen") und mit der DoD-Behauptung. Kein Fall von „Zahl
im Träger driftet gegen die Messung" (`AGENTS.md` §3.12).

## 2. Alle zwölf Tags unabhängig gegen `spec/lastenheft.md` UND den tatsächlichen Testcode gehalten

Eigene Codelektüre an jeder der sechs Fundstellen (nicht nur der neue
Kommentar-/`abdeckung_declare`-Text), gegen den vollständigen
Akzeptanzkriterien-Wortlaut aus `spec/lastenheft.md`:

| Kennung | Fundstelle | Eigener Befund |
|---|---|---|
| `LH-FA-CON-001` (Registrierung) | `tools/harness/run-integration-tests.sh:1024-1034` | `register-consumer` per `docker exec`, danach `SELECT count(*) FROM cdc.consumer` = 1. Happy Path exakt getroffen. |
| `LH-FA-CON-003` (Positions-Persistierung) | `run-integration-tests.sh:1054-1064` | `acknowledge-consumer`, danach `cdc.consumer_position.acknowledged_position` exakt gegen den bestätigten Wert zurückgelesen. Happy Path exakt getroffen. |
| `LH-FA-CON-004` (Bestätigung) | `run-integration-tests.sh:1054-1064` | Bestätigung erfolgreich, Rücklesung bestätigt `p` als verarbeitete Position. Happy Path exakt getroffen. |
| `LH-FA-CON-005` (Fortsetzung nach Neustart) | `run-integration-tests.sh:1071-1123` | echter `docker restart` (Zeile 1071); `restored_position` unabhängig aus `cdc.consumer_position` zurückgelesen (Zeile 1111-1116, `== first_position`); `resumed_ids` ab dieser Position enthält exakt die neue Zeile (id=96), nicht die bereits bestätigte (id=95, Zeile 1118-1123). Happy Path exakt getroffen, stärker als nötig (echter Prozess-Neustart). |
| `LH-FA-DAT-002` (Quelltabelle identifizierbar) | `test/integration/integration_test.go:228-230` | `record.Change.SourceTableID != env.tableID` für alle drei Changes geprüft. Happy Path getroffen; Boundary („zwei gleichnamige Tabellen in verschiedenen Schemata") bleibt unbelegt — bereits in Plan §6 benannt. |
| `LH-FA-DAT-003` (Operationstyp) | `integration_test.go:221-227` | INSERT/UPDATE/DELETE je Change gegen `record.Change.Operation` geprüft. Happy Path exakt getroffen. |
| `LH-FA-DAT-005` (Relevante Datenwerte) | `integration_test.go:247-267` | `insertImage`/`updateImage`/`deleteImage` mit realen Spaltenwerten je Operationstyp, inkl. korrekter Ab-/Anwesenheit von Alt-/Neu-Image. Happy Path getroffen; Boundary (Spaltenausschluss) und Negative (nicht lieferbarer Wert) bleiben unbelegt — bereits in Plan §6 benannt. |
| `LH-QA-REL-004` (Idempotente Verarbeitung) | `integration_test.go:269-286` | Zweites `ReadChanges` desselben Bereichs, Länge und Inhalt (`Change.ID`, `Position.Offset`, `NewImage`) exakt gegen den Erstsatz geprüft — trifft die Messmethode wörtlich („doppeltes Lesen desselben Bereichs liefert dieselben Changes"). |
| `LH-FA-REA-001` (Lesebereich zwischen zwei Positionen) | `run-integration-tests.sh:2019-2063`, `httpclient/main.go:readChanges` | `http_read_to = http_read_position + 1` als eigener `to`-Parameter neben `from` an `GET /changes`; Antwort trägt genau die eigens eingefügte Zeile. Ein echter, begrenzter Zwei-Positionen-Bereich, kein unbegrenzter Read. Happy Path getroffen. |
| `LH-FA-RET-001` (Persistente Speicherung über Neustart) | `run-integration-tests.sh:2726-2762` | `upgrade_before_still` prüft `count(*) = 1` für die vor dem realen Container-Tausch (`--force-recreate`) erfasste Zeile (id=250), danach über `cdc.changes` weiterhin identisch lesbar. Happy Path getroffen, stärker als der Wortlaut „Neustart" nahelegt (echter Container-Tausch). |
| `LH-QA-REL-003` (Sichtbarer Unzuverlässigkeitszustand) | `run-integration-tests.sh:1298-1320` | Fehlerzustand direkt in `cdc.process_heartbeat.error_class` injiziert; `diagnose`-Ausgabe zeigt „Fehlerzustand (`LH-FA-ADM-003`): schema" und **nicht** „keiner (Normalbetrieb)". Trifft die Messmethode wörtlich („Fehlerinjektionstest ... über `LH-FA-ADM-003`"). |
| `LH-FA-CON-002` (Unabhängige Consumer) | `integration_test.go:502-508` (Kommentar, nach Fixrunde), Testkörper `:514-622` | Siehe eigener Abschnitt unten — **Teilbeleg**, Kommentar nach der Fixrunde akkurat. |

**Elf der zwölf Tags sind eigenständig als inhaltlich vollständig
gerechtfertigt bestätigt** — jeweils Happy-Path- oder
Messmethode-Wortlaut real getroffen, an unverändertem, bereits vor
diesem Slice grün laufendem Testcode.

### `LH-FA-CON-002` — Präzisierung nach F-1 eigenständig gegen den Code gehalten

Vorher (Commit `748ac7bc`, vom Review als F-1/MEDIUM beanstandet):

> „… trägt zugleich `LH-FA-CON-002`: beide Consumer bestätigen unabhängig
> voneinander, die Position des einen bleibt von der Bestätigung des
> anderen unberührt"

Das behauptet wörtlich `LH-FA-CON-002`s volle Happy-Path-Formulierung
(„beide denselben Änderungsbereich lesen und bestätigen") — der Testcode
lässt aber `behindConsumer` und `aheadConsumer` **unterschiedliche**
Positionen bestätigen, nie denselben Bereich gleichzeitig.

Nachher (Commit `d1136723`, aktueller Stand,
`integration_test.go:502-508`):

> „… trägt zugleich einen Teilbeleg für `LH-FA-CON-002`:
> `behindConsumer`s gespeicherte Position bleibt exakt sein eigener
> bestätigter Wert, unverändert durch `aheadConsumer`s spätere,
> unabhängige Bestätigung — die Happy-Path-/Boundary-Formulierung aus
> `spec/lastenheft.md` selbst, denselben Bereich gleichzeitig lesend,
> prüft dieser Testfall nicht"

Eigene Prüfung gegen den Testkörper (`integration_test.go:514-622`):
`behindConsumer` bestätigt `behindPosition` (Zeile 570), danach bestätigt
`aheadConsumer` `aheadPosition` (Zeile 573) — eine andere Position, kein
gemeinsamer Bereich. Die anschließende Prüfung
(`blockers[0].position == behindRows[0].commitPosition`, Zeile 609) zeigt
exakt und ausschließlich, dass `behindConsumer`s bereits gespeicherte
Position durch `aheadConsumer`s spätere, unabhängige Bestätigung
**nicht** überschrieben wird. Die neue Formulierung behauptet genau das
— nicht mehr, nicht weniger — und benennt explizit und korrekt, was sie
**nicht** zeigt (gleichzeitiges Lesen desselben Bereichs durch beide
Consumer, der wörtliche Happy-Path aus dem Lastenheft).

**Eigenes Ergebnis: Die präzisierte Formulierung ist jetzt akkurat.**
Sie überzeichnet den Testumfang nicht mehr — F-1 ist real behoben, nicht
nur behauptet. Der verbleibende Restumfang (kein Beleg für gleichzeitiges
Lesen desselben Bereichs durch zwei Consumer) ist im Slice-Plan §6 als
Review-Fund-Risiko explizit benannt (siehe §6 unten).

## 3. `AGENTS.md` §3.13 — `harness/README.md`s `make doc-trace`-Zeile gegen eigene Messung

`harness/README.md:129` behauptet wörtlich: „… nach dem Tag-Nachtrag für
die letzten zwölf RTM-Waisen (…) an bereits bestehenden E2E-Belegen: 76
Anforderungen, 55 Waisen ohne `trace.coverage`, **0 Waisen mit** — jede
Anforderung trägt mindestens einen Beleg."

**Ergebnis: deckungsgleich** mit der eigenen Messung aus §1 (76/0 mit
`trace.coverage`). Kein Fall der Klasse „Regel weiter als ihr Sensor"
oder „Zahl im Träger driftet gegen die Messung". Ein repo-weiter `grep`
nach der alten Waisenzahl („zwölf RTM-Waisen" als offene Zahl, „ohne
eigenen Kennungs-Tag" bei den zwölf Anforderungen) findet außerhalb
dieses Slices und des vorangegangenen `rtm-reste-sst-cfg-por`-Plans (der
die Zahl korrekt als *historischen* Ausgangspunkt referenziert, nicht als
aktuellen Stand) keinen stehen gebliebenen Verweis.

## 4. `make gates` real ausgeführt — auf `HEAD` (`d1136723`), isoliert von Fremdinterferenz

Aus den oben genannten Working-Tree-Gründen lief die maßgebliche Messung
in einem eigenständigen `git clone` von `d1136723` (nicht im
Haupt-Arbeitsbaum, dort durch eine fremde, unfertige Report-Datei eines
parallelen Verifier-Agents kontaminiert — nicht Teil dieses Diffs).
Ungepipter Lauf, Exit-Code direkt geprüft (`AGENTS.md` §3.9):

```
cd <eigenständiger Klon von d1136723>
make gates > gates.log 2>&1
ec=$?
```

Ergebnis: **`EXIT=0`** (zweifach reproduziert, erster und zweiter Lauf
identisch grün). Einzelbelege aus demselben Lauf:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 752 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check: 752 Datei(en) geprüft, 0 Befund(e)` (commits-Modul) + `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Zusätzlich bestätigt: `.harness/state/gates-passed.diffsha` wurde im
isolierten Klon frisch geschrieben — der Nachweis-Stempel lief tatsächlich
nach allen sechs grünen Checks. Der Haupt-Arbeitsbaum dieses Repos wurde
durch diese Messung **nicht** verändert (der isolierte Klon wurde nach
Abschluss der Prüfung wieder entfernt); ein `git status` im
Haupt-Arbeitsbaum vor und nach dieser Verifikation zeigt ausschließlich
Änderungen des parallelen, fremden Vorgangs (`rtm-reste-sst-cfg-por`),
keine dieses Slices.

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf `[x]`
setzbar** — real, ungepiped, mit `EXIT=0` bestätigt.

## 5. DoD-Checkbox „Review durchgeführt" — berechtigt gesetzt?

Ausgangslage: Der Review-Report fand ein MEDIUM (F-1) und ein INFO (F-2)
und verlangte für F-1 ausdrücklich eine Fixrunde; er zog die
DoD-Checkbox bewusst **nicht** selbst nach — die Reviewer-Skill-Regel
„DoD-Checkbox-Nachzug ohne Fixrunde" greift laut Report explizit nicht
(hier läuft der Nachzug regulär am Implementer-Workflow-Schritt 21, nach
der Fixrunde). Das deckt sich mit der eigenen Lektüre von
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde: die
Ausnahme gilt nur, wenn *keine* Fixrunde nötig ist; hier war sie nötig.

Eigene, unabhängige Prüfung des Fixrunden-Commits `d1136723` gegen das
Finding:

- **F-1** (MEDIUM, Kommentar überzeichnet Testumfang bei
  `LH-FA-CON-002`): real behoben, siehe eigene Prüfung in §2 oben — die
  neue Formulierung ist akkurat, nicht nur behauptet akkurat.
- **F-2** (INFO, Nachbar-Inline-Kommentar zieht `LH-QA-REL-004` nicht
  mit): laut Review-Verdikt kein Rückgabe-Zwang (unter der
  3×-Schärfungsschwelle) — korrekt nicht Gegenstand der Fixrunde.

Der Fixrunden-Commit ergänzt zusätzlich Slice-Plan §6 um den Review-Fund
als eigenen Risiko-Eintrag (siehe §6 unten) — mehr als vom Review
gefordert (das Review nannte Kommentarpräzisierung *oder*
Plan-Dokumentation als hinreichend; der Implementer tat beides).

**Ergebnis: Die Checkbox ist berechtigt auf `[x]` gesetzt.** Das eine
MEDIUM ist real behoben, kein offenes HIGH/MEDIUM verbleibt, das INFO ist
korrekt ohne Aktion dokumentiert.

## 6. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout oder Folge-Slice mit ID · *entfallen*
→ gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | Reiner Tag-Nachtrag trägt keinen zusätzlichen Regressionsschutz | „entfallen als Risiko, das ist die Art dieses Slice" | ✓ entfallen | Begründung trägt — das ist tatsächlich die deklarierte Natur dieses Slice-Typs (wie schon bei `rtm-reste-sst-cfg-por`). |
| 2 | Boundary-/Negative-Akzeptanzkriterien bleiben ohne Testbeleg | „weiter offen, benannt als unsanierte Bestandslücke (wie bereits in `harness/README.md` §`make doc-trace` … beschrieben)" | ✓ weiter offen | Klasse formal korrekt. Kein *neuer* `BEO-PGC/…`-Eintrag angelegt — der Anker ist stattdessen die bereits bestehende Beschreibung in `harness/README.md`. Das deckt sich mit dem bereits im Vorgänger-Slice (`rtm-reste-sst-cfg-por` §6, dortiges Risiko 2) etablierten Muster (Anker in `docs/user/ci-matrix-abdeckung.md` statt neuem BEO-Eintrag) — kein Einzelfall, aber ein wiederkehrendes Muster, das eine strengere Lesart von Modul 5 als eigenen `BEO-PGC`-Eintrag verlangen könnte. **Nicht blockierend**, aber der Planner sollte bei der nächsten Sichtung entscheiden, ob dieses Muster verallgemeinert werden soll. |
| 3 (F-1) | `LH-FA-CON-002`-Tag belegt nur Teilbereich | „eingetreten, in derselben Fixrunde behoben" | ✓ eingetreten | Klasse formal korrekt und inhaltlich bestätigt (§2/§5 oben: real behoben, kein Rest-Gap außer dem bereits in Risiko 2 gebündelten Boundary-Punkt). Anzumerken: Modul 5 nennt für *eingetreten* wörtlich „Carveout oder Folge-Slice mit ID" — hier wurde stattdessen direkt und vollständig innerhalb derselben Fixrunde behoben, ohne dass ein Carveout oder Folge-Slice nötig wurde. Das ist ein stärkeres Ergebnis als beide genannten Instrumente (nichts bleibt offen, das ein Carveout absichern oder ein Folge-Slice aufgreifen müsste) und daher inhaltlich tragend, auch wenn es wörtlich keines der zwei benannten Instrumente ist. **Nicht blockierend.** |

**Ergebnis: Alle drei Risiken tragen einen der drei zulässigen
Ausgangsklassen.** Zwei kleinere, nicht blockierende Beobachtungen (Risiko
2 und 3, oben ausgeführt) sind für die Planner-Sichtung notiert, ändern
aber nichts an der DoD-Konformität dieses Diffs.

## 7. Zusätzlicher eigener Fund: Zeilen-Lokator-Drift in `docs/user/e2e-abdeckung.md`

Nicht Teil der sechs beauftragten Prüfpunkte, aber beim Quellcode-Abgleich
in §2 aufgefallen: Die Fixrunde (`d1136723`) verlängerte den Doc-Kommentar
über `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer` um netto 3
Zeilen (3 ersetzt durch 6). `docs/user/e2e-abdeckung.md` wurde danach
**nicht** erneut über `make test-integration` regeneriert — die Zeile für
`LH-FA-RET-005`/`LH-FA-CON-002` zitiert weiterhin
`test/integration/integration_test.go:511`, während die Funktion jetzt
real bei Zeile 514 beginnt (eigene Messung: `grep -n "^func
TestE2ERetentionBlockersViewShowsFurthestBehindConsumer"` → 514).

Die übrigen fünf Zeilen-Lokatoren in derselben Datei (Zeilen zu
`TestE2ECaptureFlow`, „Black-Box-CLI-Rundlauf", „CLI-Diagnose-Beleg
(Fehlerzustand)", „HTTP-API-Rundlauf", „Upgrade-Sicherheits-Rundlauf")
wurden einzeln gegen den aktuellen Quellcode nachgeprüft und stimmen
exakt — nur `tools/harness/run-integration-tests.sh` war von der
Fixrunde nicht betroffen, `integration_test.go` schon.

**Einordnung:** Dies ist exakt die in `AGENTS.md` §3.13 §Grenze benannte
Blindstelle des `grep`-Nachzugs (Zeilen-Lokatoren verschieben sich ohne
wiederholbare Textspur). Inhaltlich harmlos — die Tag-Präsenz selbst ist
korrekt, `make doc-trace` zeigt weiterhin `ok` (die RTM-Prüfung hängt
nicht an der exakten Zeile) —, aber ein Leser, der dem Lokator folgt,
landet drei Zeilen daneben. Kein DoD-Blocker, da die DoD-Checkbox
„`make test-integration` grün, regeneriert `docs/user/e2e-abdeckung.md`
mit allen zwölf neuen Tags" für den Stand nach `748ac7bc` zutraf; die
Drift entstand erst durch die *nachfolgende* Fixrunde, die selbst keinen
erneuten `make test-integration`-Lauf verlangt oder durchgeführt hat.
**Empfehlung an den Planner:** vor `git mv` nach `done/` einmal `make
test-integration` erneut laufen lassen (kostenlos im Sinne des
Tag-Nachtrag-Rahmens — kein neuer Testcode, nur ein Regenerierungslauf)
oder den Lokator manuell auf `:514` korrigieren.

---

## Verdikt

**DoD-Konformität bestätigt**, mit einem benannten, nicht blockierenden
Nachbesserungspunkt (§7). Alle sechs beauftragten Prüfpunkte wurden real
und unabhängig nachgemessen, nicht aus Bericht oder Commit-Message
übernommen:

1. `make doc-trace` real ausgeführt: `76/0`, deckungsgleich mit der
   Commit-Message und der DoD-Behauptung; alle zwölf Ziel-Kennungen
   einzeln als `ok` bestätigt.
2. Alle zwölf Tags eigenständig gegen `spec/lastenheft.md` und den
   tatsächlichen Testcode gehalten — elf vollständig gerechtfertigt, der
   zwölfte (`LH-FA-CON-002`) nach der Fixrunde jetzt akkurat als
   Teilbeleg formuliert (F-1 real, nicht nur behauptet, behoben).
3. `harness/README.md`s `make doc-trace`-Zeile („0 Waisen") ist
   deckungsgleich mit der eigenen Messung — `AGENTS.md` §3.13 erfüllt.
4. `make gates` lief eigenständig, ungepiped, mit `EXIT=0` — auf einem
   isolierten Klon von `HEAD` (`d1136723`), notwendig wegen einer
   fremden, parallel laufenden Kontamination des Haupt-Arbeitsbaums, die
   nicht zu diesem Slice gehört.
5. Die DoD-Checkbox „Review durchgeführt" ist berechtigt gesetzt: das
   eine MEDIUM real behoben, kein offenes HIGH/MEDIUM, das INFO korrekt
   ohne Aktion dokumentiert.
6. Alle drei §6-Risiken (inklusive des neu ergänzten Risikos zu F-1)
   tragen eine der drei zulässigen Ausgangsklassen; zwei kleinere,
   nicht blockierende Beobachtungen für die Planner-Sichtung notiert.

**Offener, nicht blockierender Punkt:** `docs/user/e2e-abdeckung.md`
trägt für die `LH-FA-CON-002`-Zeile einen um 3 Zeilen veralteten
Datei-Lokator (`:511` statt `:514`), entstanden durch die Fixrunde ohne
begleitenden `make test-integration`-Lauf — vor Closure zu korrigieren
(§7).

**Freigabe an den Planner:** Der Slice ist bereit für Closure nach
`done/`, sobald der Lokator in `docs/user/e2e-abdeckung.md` aufgefrischt
(§7) sowie Closure-Notiz, Beobachtungs-Register-Eintrag, der formale
Risiko-Ausgangs-Nachzug und die drei Paarungen ergänzt sind.
