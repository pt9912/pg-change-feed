# Verifikations-Report: slice-transformationen-e2e-wirkung — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-e2e-wirkung.md`](review-slice-transformationen-e2e-wirkung.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-backfill-pfad.md`](verifikation-slice-transformationen-backfill-pfad.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`).

**Gegenstand:** Slice-Plan `slice-transformationen-e2e-wirkung` (Welle `welle-transformationen`), `HEAD` =
`c3348898`, Diff-Range `c246ba4f..HEAD`, 10 Commits, 6 Dateien (+1048/−58, davon 321 Zeilen der Review-Report).
Slice-Inhalt: Lifecycle (`631748d2`, `c73d4292`, `38c7b3bc`), Implementer-Lauf (`fd1cc90d`, `87281d08`, `21a34696`,
`66f60c8b`), Review-Report (`7b786343`; 1 HIGH, 1 MEDIUM, 4 LOW, 4 INFO), Fixrunde (`4800d75b`, `c3348898`). Der
Stand ist nicht gepusht (`git rev-list --count origin/main..HEAD` = 10, gleich der Range). Produktivcode liegt nicht im
Diff (`git diff --stat c246ba4f..HEAD` für `internal/`, `cmd/`, `Dockerfile`, `go.mod`: leer). Dieser Lauf ändert weder
Code noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an Kopien im Scratchpad (Runner-Kopien
aus `sed`/`awk` mit Ausgabe nach stdout, Go-Kopie und Produktivcode-Kopie in einem eigenen `git archive`-Baum mit
eigenem `git init`); das Repo blieb unberührt (`git status --short` nach jedem schweren Lauf leer).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei, der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen.
Stand aller Läufe ohne Mutation: `HEAD` = `c3348898`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich
(`free -m` vor dem ersten schweren Lauf: 15,6 GB verfügbar).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` (einmal) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1292 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 85.50% erfüllt Schwelle 80%` · `generated-sync: OK` · `gesamt: 0 Befund(e)` (a-check) |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | `OK — 10 Commit(s) in "origin/main..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1292 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1292 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 20 Zeilen stimmen` |
| `make fmt-check` | **Exit 0** | `fmt-check: 258 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=c246ba4f` | **Exit 0** | kein Kandidat |
| `make doc-trace` (Arbeitsbaum) | **Exit 0** | `80 Anforderung(en), 1 Waise(n).`; `LH-FA-CFG-007` Coverage `E2E`, Status `ok`; `LH-FA-CFG-008` `WAISE` |
| `make doc-trace` (Parent-Stand `38c7b3bc`, `git archive` im Scratchpad) | **Exit 0** | `80 Anforderung(en), 1 Waise(n).`; `LH-FA-CFG-007` Coverage `E2E`, Status `ok` (V-3) |
| `make image` | **Exit 0** | Digest `sha256:f9eee81b…` (Layer-Cache; gleich dem Digest vor dem Lauf, das `:dev`-Image entspricht `HEAD`: `internal/`, `cmd/`, `Dockerfile` und `go.mod` liegen nicht im Diff) |
| `make test-integration` (ein Lauf, Wiederholung des im Plan als **übernommen** gekennzeichneten Laufs) | **Exit 0** | Dauer 21:43:47 bis 21:49:33 = 5 min 46 s = 346 s (gemessen mit Uhrzeit-Stempeln um den Aufruf; der Implementer nannte 349 s, **übernommen**, Abstand 3 s abgeleitet); Zeilen unten |

**Integrationslauf (gedruckt).**

- Go, erster `go test`-Aufruf: `--- PASS: TestE2ETransformationRulesShapeBothImages (1.52s)`,
  `--- PASS: TestE2ETransformationConflictsFailWithSpecText (3.24s)`,
  `ok  github.com/pt9912/pg-change-feed/test/integration  6.157s`.
- Phase Happy Path: „Transformationen-Happy-Path (`LH-FA-CFG-007`) belegt — auf feed_e2e_transform prägten
  rename_column (name zu customer_name) und map_value (status) 1 Stream-Zeile(n) auf allen fünf Wegen in derselben
  Form und zwei weitere Zeilen (nicht abgebildeter und fehlender Wert) über cdc.changes und GET /changes; gRPC:
  RECEIVED change_id=2156-1 table=feed_e2e_transform operation=INSERT
  new_image={"id":"11","customer_name":"TfNeu","status":"open"}; SSE: (dieselbe Zeile); NATS: (dieselbe Zeile); die
  Zeile vor den Regeln blieb in Rohform, nach dem Entfernen der Regeln trug die nächste Change wieder die Rohform“.
- Phase Neustart und Ausschluss: „… trug die Change nach einem realen docker restart (Startzeit
  2026-09-26T19:48:31.247358798Z, danach 2026-09-26T19:48:54.972507728Z) die Form beider Regeln; … nach exclude_column
  auf name trug das Bild vor und nach dem zweiten docker restart (Startzeit 2026-09-26T19:48:54.972507728Z, danach
  2026-09-26T19:49:02.12769693Z) weder name noch customer_name noch den Wert, die Spalte status blieb (mit Regel offen,
  ohne Regel roh)“.
- Danach: „E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand“ und „Lauf
  abgeschlossen — E2E-Abdeckungstabelle aus 16 Go-Zeilen und 38 Bash-Zeilen“. `git status --short` nach dem Lauf
  leer: das Erzeugnis `docs/user/e2e-abdeckung.md` liegt am Endstand unverändert, der Lauf bewegt keine andere
  committete Datei.

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach dem
Integrationslauf **34**, nach den Mutationsläufen und der Wiederherstellung des `:dev`-Images **34**; kein `prune`,
kein `system prune`. Die eigene Wegwerf-Umgebung der Mutationen (Compose-Stack im Scratchpad-Baum) ist mit
`docker compose down -v --remove-orphans` abgeräumt (kein Container `cdc-test-*`, kein Netz `cdc-feed-test` danach).
Nach den zwei Produktivcode-Mutationen (§4) ist `make image` am unmutierten Repo erneut gelaufen: Digest wieder
`sha256:f9eee81b…`. Nicht gefahren: `make test`, `make test-store`, `make test-replication`, `make bench` (kein
Produktivcode im Diff).

**Coverage-Zahl als Lauf-Beleg mit Streuung ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A).** Mein
`coverage-gate` im Lauf von `make gates` druckt 85.50 %; der Implementer nannte 85,30 % (**übernommen** aus seinem
Bericht, nicht im Plan). Der Plan trägt die Zahl nicht (Diff gelesen); der Abstand zur Schwelle (80 %) ist mehr als
fünf Punkte.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 7 `[x]`-Zeilen, 5 `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 7, `grep -c '^- \[ \]'` 5).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Happy Path und Boundary am laufenden Container: `rename_column`/`map_value` per `cdc.set_transformation`, `applied`; dieselbe Form über `cdc.changes`, `GET /changes`, gRPC, SSE, NATS; frühere Change in Rohform; K1–K4 je `failed` mit Text der Spec; danach weiter die gültige Form; Gegenprobe `applied` | **erfüllt** | eigener Lauf Exit 0 (§1). Die fünf Wege tragen dasselbe Bild `{"id":"11","customer_name":"TfNeu","status":"open"}` (gedruckt: gRPC, SSE, NATS; `cdc.changes` und `GET /changes` über `tf_expect_form` und die Bild-Gleichheit je `change_id`). Boundary: sechs Fälle (K1, K2, K3 in beiden Formen, K4 in beiden Formen) an den Klartext von [`SPEC-019`](../../spec/pflichtenheft.md) (`spec/pflichtenheft.md:753-758`) wörtlich gebunden, Gegenprobe `applied`. Grenze der Aussage „alle Wege“: INSERT-Neu-Bild (V-2). Bindung an die Eingabeseite: 21 Mutationen der Eingabeseite rot (§4) |
| 2 | Neustart-Festigkeit und Ausschluss+Regel: realer `docker restart`, Regelstand aus den `applied`-Zeilen; `cdc.remove_transformation` stellt die Rohform her; `cdc.exclude_column` auf Spalte mit Regel: weder Quell- noch Zielname noch Wert vor und nach dem Neustart, nicht ausgeschlossene Spalte bleibt; Phase vor Container-Ende-Grenze und Upgrade-Tausch | **erfüllt** | eigener Lauf Exit 0 (§1), zwei Startzeiten je Neustart gedruckt (verschieden). Lage im Runner: Phasen bei Zeile 3464–3758 vor `abdeckung_declare "Upgrade-Sicherheits-Rundlauf"` (Zeile 3760) und vor den beiden Schema-Tests der Container-Ende-Grenze (Lauf-Log: `TestE2ESchemaChangeDropColumn` erst danach). Der zweite Neustart ist an seinen eigenen Startzeit-Vergleich gebunden (Mutationen R1, R2 einzeln rot); der Prozessstart-Beleg an Produktivcode-Mutationen (P1: Regelstand beim Prozessstart verworfen, P2: Ausschlussstand beim Prozessstart verworfen — beide rot, §4). [`LH-QA-SEC-004`](../../spec/lastenheft.md): Schlüsselmenge exakt `{"status":"open"}` (Quell- **und** Zielname fehlen) plus Wert-Prüfung `tf_no_value` |
| 3 | Abdeckung getragen: jede neue `func TestE2E*` und Runner-Phase trägt [`LH-FA-CFG-007`](../../spec/lastenheft.md) im Anker, wird von einem `-run`-Muster erfasst (Befehl im Bericht), Zeile im Erzeugnis; `make doc-trace` führt es nicht mehr als Waise; [`LH-FA-CFG-008`](../../spec/lastenheft.md) bleibt Waise | **erfüllt** (mit V-3) | Abgleich selbst gefahren: `git grep -hoE '^func TestE2E[A-Za-z0-9_]+' -- test/integration` = 16 Namen; jeder steht in mindestens einer Zeile `-run '` des Runners (`git grep -h -F -e "-run '" -- tools/harness/run-integration-tests.sh`), kein Name ohne Treffer; die 2 neuen im ersten Aufruf. Erzeugnis: 5 Zeilen mit `LH-FA-CFG-007` (`git grep -c`), Lauf schreibt „unverändert“. `make doc-trace` Arbeitsbaum: 1 Waise, `LH-FA-CFG-008`. `make docs-check` im Gate-Lauf grün |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | Report liegt vor (`7b786343`); alle zehn Findings gegen die Fixrunde nachgemessen (§5), kein offenes HIGH/MEDIUM |
| 6 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes, beide Stände | **erfüllt** | 20 Zeilen mit dem Werkzeug Exit 0 (§1); Lese-Prüfung von Suchraum und Muster (§6) |
| 7 | Doku-Update: `harness/README.md` §Sensors (Zeile `make test-integration`: vier Belege; Zeile `make doc-trace`: nachgemessene Waisen-Aussage); Handbuch unberührt bis `betriebsdoku` | **erfüllt** | Diff der beiden Zeilen gelesen: „Zusätzlich vier Transformations-Belege“ (Happy Path, Bilder aller Operationen, Boundary, Neustart mit Ausschluss = vier), Herkunftsfeld `· seit slice-transformationen-e2e-wirkung`; Zeile `make doc-trace` „gemessen 2026-09-26 … 80 Anforderungen, **1 Waise** — `LH-FA-CFG-008`“ stimmt mit meinen beiden `make doc-trace`-Läufen. `docs/user/benutzerhandbuch.md` liegt nicht im Diff; der Aufschub mit Adresse `slice-transformationen-betriebsdoku` steht im Plan (§6 letzte Zeile; die DoD-Zeile nennt „§1“, V-1) |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt | **korrekt offen** | Planner |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt keine Register-Datei |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: alle sieben Zeilen tragen „**Ausgang:** *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Die Zeilen 8 bis 12 sind
Planner-Arbeit und nicht Teil dieser Prüfung; ich habe keinen DoD-Haken gesetzt.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** Wirkung am laufenden Feed-Container, fünf Wege, Neustart, Ausschluss, Abdeckungs-Träger. Die
„Ausdrücklich NICHT“-Punkte sind eingehalten: `git diff --stat c246ba4f..HEAD` nennt weder `internal/` noch `cmd/` noch
`spec/` noch `docs/plan/adr/` noch `.github/` noch `compose.yaml`; die vier Wegwerf-Clients (`tools/harness/httpclient`,
`grpcclient`, `sseclient`, `natsstreamsub`) und der Erzeuger der Abdeckungstabelle sind unverändert; Nichtanwendbarkeit
und Abhilfe (`e2e-abhilfe`), der Backfill-Beleg (`backfill-pfad`) und das SDK-Realserver-E2E stehen nicht im Diff.

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `test/integration/integration_test.go` „nicht realisiert“ | unverändert (nicht im Stat); Grund im Plan: die Funktionen stehen in der neuen Datei nach dem Muster von `backfill_e2e_test.go`, der Erzeuger liest jede Go-Datei des Verzeichnisses (Lauf: 16 Go-Zeilen) |
| `test/integration/transformation_e2e_test.go` neu | 324 Zeilen, zwei Funktionen mit `LH-FA-CFG-007` im Doc-Kommentar; importiert nur die Standardbibliothek (`context`, `encoding/json`, `fmt`, `reflect`, `testing`, `time`), spricht das System über `pgx`-Hilfen (`newBackfillEnv`) an — kein Import interner Anwendungslogik, damit E2E-Tier im Sinn von [`ADR-0030`](../plan/adr/0030-testpyramide.md). Der Grund der Abweichung „Boundary in Go statt im Runner“ steht jetzt im Plan (Review F-10) |
| `tools/harness/run-integration-tests.sh` update | 301 Einfügungen, 2 Löschungen (`git diff --numstat 38c7b3bc HEAD`): Kopfkommentar, `-run`-Muster, Block der zwei Phasen mit den Helfern `tf_*`; je Phase ein `abdeckung_declare`-Anker mit `LH-FA-CFG-007`; `tf_restart_feed` trägt beide Neustarts; die Kopplung `TF_RULE_*` ist im Kommentar benannt |
| vier Wegwerf-Clients „geprüft, unverändert“ | `git grep -n -F 'new_image'` über die vier Verzeichnisse: 7 Treffer an beiden Ständen, kein Diff dort |
| `docs/user/e2e-abdeckung.md` Erzeugnis | 40 Einfügungen, 36 Löschungen (Zeilen-Lokatoren plus 4 neue Zeilen); die `Ort`-Angaben der neuen Zeilen zeigen an `HEAD` auf die Zeilen 253 (Go-Boundary), 3699 und 3758 (Runner-Abschlussmeldungen der Phasen, wie bei den Bestandsphasen) — nachgeprüft mit `sed -n`; der Lauf meldet „unverändert“ |
| `harness/README.md` §Sensors | 2 Einfügungen, 2 Löschungen (Zeile `make test-integration`, Zeile `make doc-trace`) |
| `compose.yaml` „geprüft, unverändert“ | nicht im Diff |

Nicht im Plan-Feld, im Diff: der Review-Report (`docs/reviews/review-slice-transformationen-e2e-wirkung.md`, Rolle
Reviewer, erwartet). Im Plan, nicht im Diff: nichts. Keine `Accepted` ADR inhaltlich geändert (`make doc-immutable`
Exit 0). Keine Schwelle, keine Gate-Konfiguration verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6). Lifecycle:
`631748d2` und `38c7b3bc` sind reine Renames (`git show -M --stat`: je 0 Zeilen), die Inhaltsänderung `c73d4292` liegt
in einem eigenen Commit dazwischen ([`AGENTS.md`](../../AGENTS.md) §3.3). Docker-only ([`AGENTS.md`](../../AGENTS.md)
§3.1): die hinzugefügten Runner-Zeilen tragen kein `sed -i`, kein `perl`, kein `python`; `//nolint` steht nirgends
([`AGENTS.md`](../../AGENTS.md) §3.2, Suche über die hinzugefügten Zeilen: 0 Treffer).

## 4. Mutationen der Eingabeseite (dieser Lauf)

Eigene Wegwerf-Umgebung (PostgreSQL, NATS, Feed-Container aus den Zeilen 1 bis 402 des committeten Runners; eigener
Baum aus `git archive HEAD`), je Lauf eigene Tabellennamen (Suffix je Mutation). Unmutierte Wiederholung zuerst: beide Go-Tests
PASS und die beiden Runner-Phasen (Zeilen 3464–3759) Exit 0. Danach je Mutation genau eine Änderung an der Kopie; jedes
`rot` bedeutet: der Lauf endete mit Exit 1 an der genannten Stelle (Log gelesen, nicht nur der Exit-Code).

Runner-Phasen (12):

| # | Zusage | Mutation (Eingabeseite) | Gesehene Farbe |
|---|---|---|---|
| R1 | der **erste** Neustart ist ein realer `docker restart` | in `tf_restart_feed` wird der Aufruf `docker restart` beim ersten Durchlauf übersprungen | **rot**: „die Startzeit des Feed-Containers blieb nach docker restart gleich“ |
| R2 | der **zweite** Neustart ist ein realer `docker restart` (Review F-2) | derselbe Aufruf wird beim zweiten Durchlauf übersprungen | **rot**: dieselbe Meldung, Phase endet vor der Zeile mit dem Wert nach dem Neustart |
| R3 | die Regel `kundenname` wirkt auf den Stream-Wegen | `tf_set … kundenname` der Happy-Path-Phase entfällt | **rot**: gRPC-Bild `{"id":"11","name":"TfNeu","status":"open"}`, „erwartet 't', gelesen 'f'“ |
| R4 | nach dem Entfernen der Regeln trägt die nächste Change die Rohform | `tf_remove … status_lesbar` entfällt | **rot**: „Rohform nach dem Entfernen der Regeln“ |
| R5 | `cdc.exclude_column` auf der Spalte mit Regel entfernt beide Namen | `tf_exclude … name` entfällt | **rot**: „Bild mit Ausschluss und Regeln vor dem Neustart“ |
| R6 | ein nicht abgebildeter Wert bleibt unverändert | der eingefügte Wert `x` wird `o` | **rot**: „cdc.changes, nicht abgebildeter Wert“ |
| R7 | der NATS-Client empfängt am Subjekt der Tabelle | Subjekt auf `…public.nope` | **rot**: „nicht jeder Stream-Client empfing eine der 6 Zeilen“ (gRPC und SSE hatten `RECEIVED`, NATS nicht) |
| R8 | `GET /changes` liest den Bereich der Tabelle | Bereich `[from, from)` statt `[from, to)` | **rot**: „leere Changes-Liste für den erwarteten Bestand“ |
| R9 | das erneute Setzen der Regeln vor dem Ausschluss wirkt | zweites `tf_set … status_lesbar` der Neustart-Phase entfällt | **rot**: „Form nach dem erneuten Setzen“ |
| R10 | der gRPC-Client authentifiziert sich mit dem Reader-Token | Token `wrong-token` | **rot**: „nicht jeder Stream-Client empfing …“ (gRPC leer) |
| R11 | das Bild nach dem Neustart trägt die Form mit `open` | eingefügter Wert `c` statt `o` | **rot**: „Bild mit Ausschluss und Regeln nach dem Neustart“ |
| R12 | die Zeile vor den Regeln bleibt in Rohform | der Einfüge-Schritt der Zeile 1 wandert hinter die zwei `tf_set` | **rot**: „cdc.changes, Change vor den Regeln“ |

Go-Tests (9, je unmutierte Gegenprobe grün):

| # | Zusage | Mutation | Gesehene Farbe |
|---|---|---|---|
| G1 | K2: Spalte trägt bereits eine Regel | Spalte `name` → `note` in der K2-Regelform | **rot**: „Antrag endete applied …, erwartet failed“ |
| G2 | K3 (Zielname gleicht einer anderen Regel) | `"to":"customer_name"` → `"to":"cust"` | **rot**: „endete applied“ |
| G3 | K4: Regelname nicht geführt | Regelname `gibt_es_nicht` → `kundenname` (geführt) | **rot**: „endete applied“ |
| G4 | K4: Spalte fehlt an der Quelle | `nicht_vorhanden` → `note` | **rot**: „endete applied“ |
| G5 | `map_value` prägt beide Bilder | `"c":"closed"` → `"c":"shut"` | **rot** in beiden Tests (UPDATE-Neu-Bild, Change nach der Gegenprobe) |
| G6 | `rename_column` prägt das Alt-Bild | `setRule … kundenname` des Bilder-Tests entfällt | **rot**: „UPDATE, Alt-Bild“ trägt `name` statt `customer_name` |
| G7 | DELETE trägt das Alt-Bild der gelöschten Zeile | `DELETE … WHERE id = 2` → `id = 3` | **rot**: „DELETE, Alt-Bild“ |
| G8 | ein nicht abgebildeter Wert bleibt | Wert `x` → `o` in der INSERT-Zeile | **rot**: „INSERT mit nicht abgebildetem Wert“ |
| G9 | der Regelstand bleibt nach abgelehnten Anträgen | `setRule … status_lesbar` des Boundary-Tests entfällt | **rot**: „Change nach den abgelehnten Anträgen“ trägt `status:o` |

Produktivcode (2, Eingabeseite der Verdrahtung des Prozessstarts; Bau des Images aus der Kopie, Feed-Container neu
angelegt, danach `make image` am unmutierten Repo):

| # | Zusage | Mutation | Gesehene Farbe |
|---|---|---|---|
| P1 | der Prozessstart leitet den Regelstand aus den `applied`-Zeilen ab | `internal/bootstrap/wiring.go`, Bindung des Prozessstarts: `Transformations: rules[…][:0]` | **rot**: Phase 1 grün (kein Neustart), Phase 2 an „Form nach dem Neustart — erwartet 't', gelesen 'f'“ |
| P2 | der Prozessstart leitet den Ausschlussstand ab, auch mit Regel | dieselbe Bindung: `ExcludedColumns: excluded[…][:0]` | **rot**: „Bild mit Ausschluss und Regeln nach dem Neustart“ |

Zusammen 23 Mutationen, alle rot; kein Fall grün. Die 15 Mutationen des Reviews (Kopien, Wegwerf-Umgebung) sind
**übernommen**, nicht wiederholt; meine Auswahl deckt dieselben Zusagen an `HEAD` ab und ergänzt die zwei Neustarts
einzeln (R1, R2), die Bindung der Zusage „nach dem Neustart“ an den Produktivcode (P1, P2) und die K3-/K4-Formen und
die Alt-Bilder (G2–G4, G6, G7). Die Zahl in DoD Nr. 1 („21 Mutationen der Eingabeseite“) ist die Summe der zwölf
Runner- und neun Go-Mutationen (abgeleitet).

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Konjunktiv-Kommentar im Runner | Kommentar vor `tf_restart_feed "$TF_PHASE"` selbst gelesen: „Der Neustart ist ein echter Prozess-Neustart desselben Containers: der Prozessstart leitet den Regelstand aus den applied-Zeilen ab, und die Form der nächsten Change ist der Beleg dieser Ableitung. Der Beleg des Neustarts selbst ist die Startzeit des Containers (tf_restart_feed).“ — Indikativ, Klassen Zusage/Kopplung. Suche über alle hinzugefügten Kommentarzeilen von Runner und Go-Datei nach `würde/wäre/trüge/hätte/früher/vorher/nicht mehr/statt/sonst/ursprünglich/bisher`: keine Konjunktiv- oder Vorher-Nachher-Zeile; `make kommentar-kennungen DIFF=c246ba4f` Exit 0 | **behoben** |
| F-2 (MEDIUM) zweiter Neustart nicht gebunden | `tf_restart_feed` liest die Startzeit vor und nach `docker restart` und beendet die Phase bei Gleichheit; beide Neustarts der Phase rufen ihn auf; die Abschlussmeldung druckt beide Startzeit-Paare (gedruckt, §1). Mutationen R1 und R2 einzeln rot; die Fixrunde meldete den ersten Aufruf als nicht einzeln mutiert — ich habe beide einzeln gefahren | **behoben** |
| F-3 (LOW) „in genau einem Feld“ | Doc-Kommentar jetzt: fünf `set_transformation`-Negativfälle weichen in genau einem Feld von der Gegenprobe ab; der `remove_transformation`-Fall nennt einen nicht geführten Regelnamen; das Entfernen einer geführten Regel endet in `removeRulesAtEnd` `applied`. Am Test nachgemessen: K1 (Regelname), K2 (Spalte), K3 zweimal (Zielname), K4 Spalte (Spalte) — je ein Feld; G1–G4 rot | **behoben** |
| F-4 (LOW) Fenster READY → erste Zeile | Plan §6 trägt das Risiko im Ist-Ton mit der Messung am Host des Reviews, dem Grenzsatz „auf dem Runner von `e2e.yml` nicht gemessen“, dem Beleg-Anker (erster realer `e2e.yml`-Lauf nach dem Push, [`AGENTS.md`](../../AGENTS.md) §3.10) und `Ausgang: (bei Closure)`. Nachgemessen: `nproc` = 20; `go build` kalt im Toolchain-Container 5,32 s (`natsstreamsub`), 4,82 s (`sseclient`), 6,18 s (`grpcclient`) — gleich der Plan-Angabe 4,85 s bis 6,15 s im Rahmen der Messstreuung; Fristen im Code gelesen (NATS `NextMsg(30 * time.Second)`, gRPC und SSE `WithTimeout(…, 60*time.Second)`); im Lauf von `make test-integration` und in der unmutierten Wiederholung der Phasen genügte der erste Einfüge-Versuch (beide gedruckt `1 Stream-Zeile(n)`). Der Workflow `e2e.yml` ist unverändert, [`AGENTS.md`](../../AGENTS.md) §3.10 greift nicht; der reale Lauf nach dem Push ist offen | **im Plan geführt; Ausgang offen** (Übergabe, V-4) |
| F-5 (LOW) Beleg-Anker nicht committet | Plan-Zeile „Jede neue `TestE2E*`-Funktion …“ trägt jetzt: der Lauf ist **übernommen** aus dem Bericht des Implementers, „die Wiederholung trägt die Verifikation“. Die Wiederholung ist dieser Report (§1: Exit 0, beide Tests `PASS`). Suchlauf-Anker: 20 Zeilen mit dem Werkzeug Exit 0, alle an beiden Ständen | **behoben** (mit diesem Report) |
| F-6 (LOW) fremder Träger `welle-transformationen.md` | Plan §3 nennt Adresse (Planner) und Frist (Closure dieses Slice) und die Behauptung „das Ziel ist bereits vor diesem Slice erreicht“. Nachgemessen: `git show fc0b8d38^:docs/user/e2e-abdeckung.md`, durch `grep -c CFG-007` gezählt, ergibt 0, `git show fc0b8d38:…` ergibt 1, `git show 38c7b3bc:…` ergibt 1; `make doc-trace` am Parent-Stand `38c7b3bc`: 1 Waise (`LH-FA-CFG-008`) — die Behauptung stimmt. `git grep -n -E 'verlässt die Waisen\|nicht mehr (unter den )?Waise' -- docs/plan/planning/welle-transformationen.md`: Zeile 53 und Zeile 168 (die Zeile 168 schreibt diesem Slice die Wirkung „`LH-FA-CFG-007` verlässt die Waisen“ zu). Der Träger ist gemeldet, nicht nachgezogen | **Meldung konform; Nachzug offen beim Planner** (V-3) |
| F-7 (INFO) Alt-Bild und UPDATE/DELETE nur über `cdc.changes` | Plan §3 (Zeile der Clients) und der Deklarations-Anker der Phase im Runner nennen die Grenze; die Zeile im Erzeugnis trägt sie (gelesen, `Transformationen-Happy-Path`) | **benannte Grenze** (V-2) |
| F-8 (INFO) Kopplung `TF_RULE_*` zwischen den Phasen | Kommentar am Block der ersten Phase nennt die Kopplung (gelesen) | **behoben** |
| F-9 (INFO) Erzeugnis an Zwischenständen | Plan §3 (Zeile `docs/user/e2e-abdeckung.md`) nennt die Verschiebung und die Kommentar-Kürzung in `66f60c8b`; am Endstand ist die Datei byte-gleich dem Erzeuger-Ergebnis (mein Lauf: „unverändert“, `git status` leer; die Ort-Angaben 253/3699/3758 stimmen mit `sed -n`) | **behoben**; Bisect-Zwischenstände tragen abweichende Nummern (im Plan benannt) |
| F-10 (INFO) Abweichung Boundary in Go ohne Grund | Plan §3 nennt den Grund (nur SQL-Funktionen, `cdc.administration_request` und `cdc.changes`, kein Container-Zugriff, Klartext wörtlich, erster `go test`-Aufruf, kein Import interner Logik); Importliste der Datei gelesen | **behoben** |

Kein offenes HIGH, kein offenes MEDIUM.

## 6. Suchlauf-Feld (§3 des Plans), Lese-Prüfung

Das Werkzeug prüft Zahlen und Stände, nicht Suchraum und Muster (Exit 0, 20 Zeilen, §1). Gelesen:

- Suchraum: der ganze Baum ohne `docs/reviews`, `docs/plan/planning/done` und `.harness/baseline`; die Plan-Datei
  schließt das Werkzeug selbst aus. Muster: Symbolname einer bestehenden Phase (`Spaltenausschluss-Rundlauf`), Zählwort
  `2 Waisen`, Beschreibung `verlässt die Waisen` bzw. `nicht mehr … Waise`, `func TestE2E`, die zwei neuen Testnamen im
  Runner, `new_image` in den Clients, die `^func TestE2E`- und `^abdeckung_declare `-Zeilen — drei Arten (Symbol,
  Zählwort, Beschreibung) vorhanden.
- Eigene Gegensuche nach Trägern der bewegten Eigenschaften (`git grep -l -E 'Leerlauf-Bestätigung|Backfill-Regelstand|Spaltenausschluss-Rundlauf'`
  und `git grep -n -i -E 'Transformations-(Rundl|Phase|Beleg)|slice-transformationen-e2e-wirkung'`, Baum ohne Records):
  weitere Nennungen der Phasenliste stehen in `harness/sensors/db-adapter-coverage.md` und `harness/targets/bench-backfill.md`
  (Backfill-/Leerlauf-Gegenstände, nicht die Aufzählung der Belege von `make test-integration`) und in `open/`-Plänen der
  Folge-Slices, die den Slice nur als Start-Kante nennen; kein weiterer Träger der Belegliste oder einer Waisen-Zahl
  außerhalb des Feldes. Der einzige Träger außerhalb des Feldes bleibt `welle-transformationen.md:168` (V-3).
- Die Zahlen des Plans stimmen an beiden Ständen (Werkzeug): 12/12 `Spaltenausschluss-Rundlauf`, 2 → 1 `2 Waisen`,
  19 → 21 und 14 → 16 `func TestE2E`, 36 → 38 `abdeckung_declare`, 1 → 5 `LH-FA-CFG-007` im Erzeugnis, 7/7
  `new_image`, Stat 324/301/40/2 (gleich am Endstand, `git diff --numstat 38c7b3bc HEAD`).

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5 (E2E-Belege)**
  und §Fitness Function „Realer Rundlauf“: Happy Path am laufenden Feed-Container (Regel per SQL beantragt, `applied`,
  danach erfasste Change über `cdc.changes` und die Zustellwege mit umbenanntem Schlüssel), Neustart-Festigkeit,
  Ausschluss+Regel ([`LH-QA-SEC-004`](../../spec/lastenheft.md): weder Quell- noch Zielname noch Wert) sind belegt
  (§2, §4). Die **Nichtanwendbarkeit** und das Abhilfe-Akzeptanzkriterium (a)–(d) gehören nach dem Schnitt der Welle
  (`welle-transformationen.md` §4 Abweichung 2) zu `e2e-abhilfe`; das ist im Plan §1 ausdrücklich ausgenommen und kein
  Verstoß. Teilfrage 6 Option D (alle Wege dieselbe Form): für das Neu-Bild einer eingefügten Zeile auf allen fünf
  Wegen erprobt, für UPDATE/DELETE und das Alt-Bild über `cdc.changes` (V-2). Konform mit benannter Grenze.
- **[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md):** der Ausschlussstand überlebt den
  Neustart auch neben Regeln (P2 rot, R5 rot). Konform.
- **[`ADR-0030`](../plan/adr/0030-testpyramide.md) (E2E-Tier):** die Belege laufen ausschließlich über SQL-Funktionen
  und Views, HTTP, die Wegwerf-Clients und `docker restart`; die neue Go-Datei importiert keine interne
  Anwendungslogik. Konform.
- **[`SPEC-019`](../../spec/pflichtenheft.md):** die sechs Fehlertexte (K1, K2, K3 zweimal, K4 zweimal) stehen wörtlich
  in der Zeilenfolge der Spec-Tabelle (`spec/pflichtenheft.md:753-758`), die Adresse trägt die Form
  `schema.table.name`; der Test vergleicht exakt (G1–G4 rot bei einer abweichenden Eingabe). Konform.
- **[`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-FA-CFG-005`](../../spec/lastenheft.md),
  [`LH-FA-REA-005`](../../spec/lastenheft.md), [`LH-FA-SST-002`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md):** Ausschluss gilt zuerst
  (Schlüsselmenge `{"status":"open"}`), dieselbe Change beim erneuten Lesen (zweite Lesung im Go-Test gleich der ersten),
  die drei Zustellwege plus `GET /changes` und `cdc.changes` (§1). Die Feldzahlen der Nachrichten im Plan (zehn Felder
  der Live-Nachricht, dreizehn der HTTP-Antwort) am Code nachgezählt: `message Change` in
  `proto/cdc/stream/v1/changestream.proto` 10 Felder, `readChangeResponse` in
  `internal/adapters/driving/http/readchanges.go` 13 Felder.
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (Moves rein), §3.5/§3.6 (keine `Accepted` ADR überschrieben, keine Schwelle
  gesenkt), §3.7 (Kommentare im Ist-Ton, ein Herkunftsfeld je Go-Kommentar, `make kommentar-kennungen` ohne Kandidat),
  §3.9 (meine Läufe ungepiped, Exit gesondert), §3.10 (kein Workflow im Diff; der reale Lauf von `e2e.yml` nach dem Push
  bleibt für F-4 offen), §3.12 (Zahlen des Plans tragen Herkunft oder Nachmessung; der Lauf ist als **übernommen**
  gekennzeichnet und durch diese Wiederholung getragen), §3.13 (Suchlauf beide Stände, Meldung des fremden Trägers mit
  Adresse und Frist). Docker-only: keine Host-Toolchain im Diff.
- **Commits:** `make doc-commits` und `make commit-traceability` über die Range Exit 0; jeder der zehn Betreffs nennt
  `LH-FA-CFG-007` und/oder `ADR-0112`, keiner trägt `SPEC-`/`ARC-` im Betreff.
- **Handbuch:** `docs/user/benutzerhandbuch.md` liegt nicht im Diff; keine neue Betreiber-Oberfläche
  (`internal/`, `cmd/`, `tools/schema` unberührt, keine `CDC_*`-Variable).

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | INFO | **Adresse einer DoD-Zeile.** Die DoD-Zeile „Doku-Update“ nennt als Ort des Handbuch-Aufschubs „(§1, Welle §4 Abweichung 3)“; §1 des Plans trägt keinen Satz zum Handbuch, der Aufschub steht in §6 (letzte Zeile, Adresse `slice-transformationen-betriebsdoku`). Die Welle §4 Abweichung 3 ist vorhanden (`welle-transformationen.md` Zeile 410 nennt den Handbuch-Anteil an `betriebsdoku`). Kein Befund am Code; der Verweis „§1“ ist ein Planner-Text (Vorschlag: „§6“). Zusatz: die Plan-Zeile zur Go-Datei nennt „K1, K2, K3 in beiden Formen, K4 in beiden Formen“; die Spec führt K1 und K2 in je einer Form, K3 und K4 in je zwei — getestet sind alle sechs Zeilen, nur die Formulierung ist mehrdeutig | Plan §2 (DoD 7), §3 (Zeile Go-Datei); `spec/pflichtenheft.md:753-758` | ja — `git grep -n 'Abweichung 3' -- docs/plan/planning/welle-transformationen.md` |
| V-2 | LOW | **Aussagegrenze „alle Wege dieselbe Form“ ohne Träger für die Restfläche.** Die Wegwerf-Clients drucken das Neu-Bild (`new_image`); die Phase belegt die Form auf den fünf Wegen für eine eingefügte Zeile. UPDATE, DELETE und das Alt-Bild sind nur über `cdc.changes` belegt (`TestE2ETransformationRulesShapeBothImages`). Der Plan, der Runner-Anker und das Erzeugnis benennen die Grenze ehrlich; eine Adresse für ihre Schließung (Slice oder Register-Eintrag) trägt keiner. Da die Transformation vor der Persistenz wirkt und alle Wege dieselbe `model.Change` lesen (Architektur von [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 6), ist ein Unterschied auf den Wegen nicht zu erwarten; belegt ist er für UPDATE/DELETE auf den Stream-Wegen nicht. **Bewertung:** keine DoD-Verletzung (die DoD verlangt die Form der „danach erfassten Change“); der Ausgang des §6-Risikos „Ein Zustellweg trägt eine abweichende Form“ ist damit „entfallen für das INSERT-Neu-Bild, Grenze für die Restfläche benannt“ und gehört so in die Closure-Notiz. Ob die Restfläche einen Träger braucht (etwa `slice-sdk-regel-realserver-e2e`, dessen Umfang SDK-Clients trägt, nicht die Wegwerf-Clients dieses Repos), entscheidet der Planner | Plan §3 (Zeile der Clients), §6 (erste Zeile) | ja — `git grep -n 'new_image' -- tools/harness/grpcclient tools/harness/sseclient tools/harness/natsstreamsub tools/harness/httpclient` |
| V-3 | LOW | **Träger-Nachzug `welle-transformationen.md:168` offen; das DoD-Kriterium „nicht mehr Waise“ war vor dem Slice erfüllt.** `make doc-trace` am Parent-Stand `38c7b3bc` (eigener Lauf): 1 Waise, `LH-FA-CFG-007` Coverage `E2E` (Zeile der Phase „Backfill-Regelstand“ aus `fc0b8d38`). Die Aussage der Zeile `make doc-trace` in `harness/README.md` („2 Waisen“) war seit `fc0b8d38` veraltet und ist mit diesem Slice richtiggestellt (2 → 1); die Ursache der Zahl ist nicht dieser Slice, sondern der Backfill-Regelstand. Die Welle-Tabelle (Zeile 168) schreibt `e2e-wirkung` „`LH-FA-CFG-007` verlässt die Waisen“ zu; der Satz in Zeile 53 („ist im RTM-Lauf nicht mehr Waise“) beschreibt den Zustand der Welle und bleibt wahr. **Meldung an den Planner** (Träger in fremder Datei): Zeile 168 bei der Closure an den Ist-Stand ziehen (dieser Slice fügt vier Träger hinzu, das Ziel war erreicht); Adresse und Frist stehen im Plan (Planner, Closure dieses Slice) | `docs/plan/planning/welle-transformationen.md:168`; `harness/README.md` Zeile `make doc-trace` | ja — `git grep -n 'verlässt die Waisen' -- docs/plan/planning/welle-transformationen.md`; `make doc-trace` am Parent-Stand |
| V-4 | INFO | **Reale Bestätigung von `e2e.yml` nach dem Push ist Closure-Pflicht.** Der Workflow `e2e.yml` fährt das Testpaket unverändert (§3.10 verlangt keinen Workflow-Beleg), aber die zwei neuen Phasen wachsen die Laufzeit und tragen das Fenster READY → erste Zeile (Plan §6, F-4); auf dem Runner ist beides ungemessen. Meine Laufzeit: 346 s für den ganzen `make test-integration` (gemessen, lokal, 20 Kerne); die Vorgänger-Verifikation nannte 314 s für `make image` und `make test-integration` (**übernommen** aus `verifikation-slice-transformationen-backfill-pfad.md`, andere Bedingungen) — ein Unterschied von 32 s (abgeleitet, nicht auf zwei Phasen zurückführbar). `e2e.yml` trägt `timeout-minutes: 60` (gelesen, `.github/workflows/e2e.yml` Zeile 69). Nach dem Push: `gh run list --workflow e2e.yml` bzw. `gh run view`, beide PostgreSQL-Legs, und die beiden Risiko-Ausgänge (Laufzeit, Fenster) aus dem Lauf nachtragen | Plan §6 (Zeilen „Laufzeit“ und „Fenster“); [`AGENTS.md`](../../AGENTS.md) §3.10 | ja — `gh run list --workflow e2e.yml` nach dem Push |
| V-5 | INFO | **Coverage-Zahl ist ein Lauf-Beleg mit Streuung.** 85.50 % (mein `make gates`); Implementer 85,30 % (**übernommen**); frühere Verifikation 85.00–85.10 % bei anderem Stand. Der Plan trägt die Zahl nicht | §1 | ja — `make coverage-gate` |
| V-6 | INFO | **Übernommen, nicht nachgemessen:** die 15 Mutationen des Reviews (mein §4 wiederholt sie nicht, sondern ergänzt sie), die 349 s und die 85,30 % des Implementers; nicht gefahren: `make test`, `make test-store`, `make test-replication`, `make bench` (kein Produktivcode im Diff); nicht nachgemessen: das Verhalten des Fensters READY → erste Zeile auf dem Runner von `e2e.yml` (V-4) | §1, §4 | nein |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der sieben `[x]`-Zeilen (Nr. 1–7) ist am Ist-Zustand belegt; die fünf `[ ]`-Zeilen
(Nr. 8–12) sind korrekt offen (Planner-Closure). **Plan-vs-Code:** keine unbenannte Abweichung; die einzige nicht im
Plan-Feld stehende Datei ist der Review-Report (Rolle Reviewer). Produktivcode liegt nicht im Diff.
**Entscheidungs-Konformität:** [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 ist für Happy Path, Neustart und Ausschluss+Regel erfüllt; Nichtanwendbarkeit und Abhilfe liegen nach
dem Schnitt der Welle bei `e2e-abhilfe`. [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
[`ADR-0030`](../plan/adr/0030-testpyramide.md), [`SPEC-019`](../../spec/pflichtenheft.md) konform.
**Review-Findings:** F-1 bis F-10 gemessen — F-1, F-2, F-3, F-5, F-8, F-9, F-10 behoben; F-4 im Plan geführt, Ausgang
offen (V-4); F-6 Meldung konform, Nachzug offen (V-3); F-7 benannte Grenze (V-2). **Gates:** `make gates` (Exit 0),
`make commit-traceability`/`doc-commits`/`doc-immutable` über `origin/main..HEAD` (Exit 0), `make suchlauf-nachmessen`
(20 Zeilen), `make fmt-check`, `make kommentar-kennungen`, `make doc-trace` (1 Waise) und der reale `make image` +
`make test-integration`-Lauf (346 s, Exit 0, beide neuen Phasen und beide Go-Tests im gedruckten Log, Erzeugnis am
Endstand unverändert) im eigenen Lauf grün; 23 von 23 Eingabeseiten-Mutationen rot (12 Runner, 9 Go, 2 Produktivcode).

**Übergabe:** an den Planner — (1) Closure-Notiz mit Lerneintrag; Klassen aus Review und diesem Report: „Kommentar mit
Konjunktiv über die verworfene Alternative“, „Zusage ohne Bindung an die Eingabeseite (zweiter Neustart)“, „Träger-Nachzug
in fremder Datei mit Frist“, „Aussagegrenze der Zustellwege ohne Träger“. (2) Ausgänge der §6-Risiken (Vorschläge):
*Ein Zustellweg trägt eine abweichende Form* — entfallen für das INSERT-Neu-Bild auf allen fünf Wegen (Lauf 2026-09-26),
Restfläche benannt (V-2); *Laufzeit* — 346 s, gemessen 2026-09-26 (dieser Lauf), Runner-Zahl aus dem ersten
`e2e.yml`-Lauf nach dem Push; *Eine neue Testfunktion fällt still aus dem Runner* — entfallen (Abgleich 16/16, §2 Nr. 3);
*Geteilter Zustand zwischen Phasen* — entfallen (eigene Tabelle je Phase, Rücknahme der Regeln in beiden Phasen und in
beiden Go-Tests, Lauf grün); *Neustart-Beleg misst Timing statt Zustand* — entfallen (Poll auf Zustand mit Frist, Startzeit
als Beleg, R1/R2/P1/P2 rot); *Fenster READY → erste Zeile* — weiter offen bis zum ersten `e2e.yml`-Lauf nach dem Push
(§3.10, V-4); *Das Handbuch entsteht erst danach* — Aufschub mit Adresse `slice-transformationen-betriebsdoku`
(eingetreten als geplanter Aufschub). (3) Träger `welle-transformationen.md:168` (V-3, Frist: Closure dieses Slice).
(4) V-1 (DoD-Verweis „§1“ statt „§6“) als Textpflege. (5) Beobachtungs-Register: keine neue Datei im Diff — Anfall oder
„keine Beobachtung angefallen“ notiert der Planner; die Klassen des Reviews (Konjunktiv-Kommentar, ungebundener Zweit-Neustart)
sind Kandidaten für den Zähler der bestehenden Einträge. (6) Paarungen bei der Closure der Welle (Zusagen an
`slice-transformationen-e2e-abhilfe`, `slice-transformationen-start-reihenfolge`, `slice-sdk-regel-realserver-e2e`
bleiben unberührt). Nicht gepusht; der Push und der `e2e.yml`-Beleg folgen der Freigabe. Dieser Report ist ein
**Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
