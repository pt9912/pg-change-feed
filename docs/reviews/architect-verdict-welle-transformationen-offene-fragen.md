# Architect-Verdikt — offene Fragen von `welle-transformationen` (Zeitgrenze des Vorlaufs, Lesekosten, Restfläche, Regelgrenze, Paare-Obergrenze, Wortliste, Register-Lesung)

**Datum:** 2026-09-27 · **Stand:** `ca9aaea7` (Baum sauber, nichts gepusht) · **Rolle:**
Architect (frischer Kontext; jede Zahl dieses Zugs ist am genannten Stand selbst gemessen oder
als übernommen bzw. abgeleitet gekennzeichnet, Ursprung je Zeile in §1) · **Anlass:** Zug vor
dem Start von `slice-capture-leerlauf-quellbelege` und `slice-transformationen-e2e-abhilfe`;
drei Closure-Kriterien von `welle-transformationen` (Zeitgrenze des Vorlaufs, Lesekosten des
Backfill-Runs, Restfläche der Zustellwege), die Fragen (a) bis (e) in §5 der Welle und vier
Register-Einträge des Lese-Schritts · **Vollmacht:** „Entscheidungen nach Evidenz, Empfehlung
statt Optionsliste“ (Auftraggeber).

**Bezug:** [`LH-FA-CFG-007`](../../spec/lastenheft.md),
[`LH-FA-ADM-003`](../../spec/lastenheft.md),
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
[`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
(mit diesem Zug geschrieben), [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
[`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md),
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md),
[`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.12, §3.13,
[`MR-003`](../../harness/conventions/MR-003-guard-inplace-textwerkzeug.md). Vorbild für Form
und Tiefe: [`architect-verdict-welle-backfill-bestand-lese-schritt`](architect-verdict-welle-backfill-bestand-lese-schritt.md).

---

## 0. Ergebnis in einem Blick

- **(A) Zeitgrenze des Vorlaufs — Umsetzung, ADR und Slice.** Der Vorlauf wartet nicht nur ohne
  Grenze und „gesund“ (am laufenden Prozess gemessen), er **beendet den Prozess** mit der Klasse
  `replication`, sobald er `wal_sender_timeout` überschreitet, weil `NewStream` den Strom vor dem
  Vorlauf startet. [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md):
  `START_REPLICATION` wandert in `Stream.Run`, der Vorlauf trägt 30 s Frist, bei Ablauf startet
  der Stream und der Antrag bleibt `pending`; Umsetzung `slice-start-vorlauf-grenze` (wellenlos,
  vor `slice-transformationen-e2e-abhilfe`, §9).
- **(B) Lesekosten des Backfill-Runs — akzeptiertes Negativ mit Messung und Trigger.** Eine
  Lesung kostet bei 10 000 Zeilen der Queue rund 6 ms (gemessen und abgeleitet, §3), ein Run
  liest zweimal je Block; das sind etwa 8 % einer Blockdauer, bei 100 000 Zeilen rund 70 %.
  Trigger: mehr als 10 000 Zeilen in der Queue einer Quelle; die Behebung wäre ein Index samt
  tabellenbezogenem Lesezugriff (Adresse: §3).
- **(C) Restfläche der Zustellwege — als „hergeleitet“ geführt, kein Folge-Slice.** Die drei
  Adapter kopieren `OldImage`/`NewImage` einer Change unverändert; die gemessene Menge (INSERT,
  Neu-Bild, fünf Wege) steht schon in der Abdeckungstabelle. `slice-sdk-regel-realserver-e2e`
  trägt die Fläche nicht und soll sie nicht tragen (§4).
- **(D) `AGENTS.md` §3.1 und Umleitungen — geschärft.** Der Absatz nennt jetzt das Schreiben und
  Anhängen von Text an eine Repo-Datei per Umleitung als verboten (Ausnahme: Scratchpad und
  Temp-Verzeichnis, ganze Dateien kopieren und verschieben); der Guard liest sie weiter nicht,
  das Review ist der Wächter (§5, umgesetzt in `AGENTS.md`, Reviewer-Skill, `implement-slice`).
- **(E) Obergrenze der Paare von `map_value` — benannter Verzicht.** Nachgemessen: 6,5 bis
  6,7 µs je Suche bei 1000 Paaren, 0,66 bis 0,70 ms bei 100 000; erst ab rund 100 000 Paaren
  reichte die Suche allein nicht mehr für die Stufe „groß“ (abgeleitet). Keine Änderung an
  `SPEC-030`; Trigger und Handbuch-Zahlen in §6.
- **(F) Wortliste des Kandidatenlaufs — erweitert.** Der Lauf trifft den Fund F-1 von
  `slice-transformationen-e2e-wirkung` („trüge“) jetzt (1 statt 0 Zeile am Diff des Fundes); am
  Bestand steigt er von 322 auf 359 Zeilen. Umgesetzt in `.claude/commands/implement-slice.md`
  Schritt 20 (§7).
- **(G) Vier Register-Einträge gelesen:** `test-name-behauptet-mehr-als-der-test-treibt` (4×) und
  `plan-zusage-erfuellung-ohne-committeten-anker` (3×) **gestrichen** als akzeptierte Negative,
  `gemeldete-ungenauigkeit-ohne-traeger` (3×) **verkörpert** in einer Zeile des
  Closure-Note-Reviewer-Skills, `adapter-unittest-verdeckt-bootstrap-luecke` (4× mit der
  Messung dieses Zugs) **geplant** über `slice-start-vorlauf-grenze` (§8).
- **Ein neuer Slice-Vorschlag** (`slice-start-vorlauf-grenze`), **eine neue ADR**
  ([`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)),
  **keine Frage an den Nutzer** (§10).

---

## 1. Messungen dieses Zugs

Host Linux 6.8.0-139-generic, Intel Core i9-13900H (20 Threads), Docker, Stand `ca9aaea7`,
2026-09-27. Alle Skripte und Benchmark-Dateien liegen im Scratchpad des Zugs und sind nicht
committet; der Arbeitsbaum ist unberührt. Docker-Ausgangsstand 34 dangling Volumes (kein
`prune`); die Last anderer Läufe auf dem Host ist nicht gemessen.

| Nr | Was, womit | Ergebnis (gedruckt) |
|---|---|---|
| M1 | Go-Benchmark in einer Kopie von `go.mod`, `go.sum`, `internal/domain` (`git archive HEAD`), `golang:1.27-alpine@sha256:cf6fca66…` mit `--network none`, ohne `-race`, `-benchtime 300x -count 3`; Schlüssel der Suche fehlt (durchläuft alle Paare), Regeln `map_value` und `rename_column` | `lookupMappedValue` bei 10 Paaren 54 bis 123 ns, bei 1000: 6,5 bis 6,7 µs, bei 10 000: 65 bis 80 µs, bei 100 000: 0,66 bis 0,70 ms; `ParseTransformationSpec` plus `Build` bei 10 Paaren (309 B) 8,8 bis 12,5 µs, bei 1000 (26 KB) 0,71 bis 1,07 ms, bei 10 000 (260 KB) 8,0 bis 8,2 ms, bei 100 000 (2,6 MB) 82 bis 83 ms; `FoldTransformations` über 10, 1000 und 100 000 kleine Anträge (fünf Regelnamen im Wechsel `set`/`remove`) 8,8 bis 15,6 µs, 0,91 bis 0,97 ms, 94 ms |
| M2 | PostgreSQL 18 (`postgres:18-alpine@sha256:63bdc97d…`), `psql` im Container über den Unix-Socket, `\timing`, Ergebnis nach `/dev/null`; eine Tabelle mit den Spalten von `cdc.administration_request` und nur dem Primärschlüssel, gefüllt mit `n` Zeilen (fünf Antragsarten im Wechsel, davon zwei die Transformationsarten, 10 % einer fremden Quelle, alle `applied`); der Text von `SelectAppliedTransformationRequests`, fünf Läufe je Größe (der erste kalt) | `n = 100` (30 Zeilen gelesen, 64 kB): 4,0 (kalt), 0,29, 0,31, 0,30, 0,26 ms; `n = 10 000` (3000 gelesen, 1,7 MB): 8,6 (kalt), 3,4, 2,9, 3,8, 3,1 ms; `n = 100 000` (30 000 gelesen, 17 MB): 32,4 (kalt), 28,3, 29,8, 28,2, 28,0 ms; `n = 1 000 000` (300 000 gelesen, 167 MB): 203 (kalt), 192, 186, 201, 201 ms |
| M3 | drei Läufe „Vorlauf wartet an einer Sperre“ am Feed-Container, Aufbau und gedruckte Zeilen in [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md) §Gemessen | Healthcheck-Exit 0 über die ganze Wartezeit (22 und 23 Abfragen), Heartbeat-Alter höchstens 7,1 s, eingefügte Zeile nie sichtbar; der Server beendet den Walsender nach 2,0 s (Lauf B) bzw. 60,0 s (Lauf C, PostgreSQL-Log „terminating walsender process due to replication timeout“), der Prozess endet, sobald der Vorlauf endet, mit Klasse `replication` |
| M4 | `psql "dbname=postgres user=postgres replication=database"` gegen PostgreSQL 18 mit `-c wal_sender_timeout=2000`: `IDENTIFY_SYSTEM`, 7 s, `IDENTIFY_SYSTEM`, `CREATE_REPLICATION_SLOT … LOGICAL pgoutput`, 7 s, `IDENTIFY_SYSTEM` | alle Kommandos antworten, das Log trägt keine Zeile „terminating“ |
| M5 | `git grep -nE "(//\|#).*(<Liste>)" -- '*.go' '*.sh' '*.awk'` am Baum, Liste alt und neu (§7) | alt 322 Zeilen (davon 62 ohne `sonst`/`statt`), neu 359 (103 ohne `sonst`/`statt`); 42 Zeilen tragen eine der neuen Formen |
| M6 | `git diff -U0 c246ba4f 66f60c8b -- '*.go' '*.sh' '*.awk' \| grep -cE '^\+.*(//\|#).*(<Liste>)'` (Diff des Fundes F-1 von `slice-transformationen-e2e-wirkung`) | alt 0, neu 1 (die Zeile mit „trüge“) |
| M7 | `git tag --contains 06719d4f` (Erst-Commit der Regeltypen), `git tag -l` | keine Ausgabe; die Server-Tags tragen `v0.1.0` bis `v0.2.0`, keiner enthält Transformationen |
| M8 | `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence \| wc -l` je Register-Eintrag aus §8 | `wartegrenze-…` 1, `inplace-textwerkzeug-…` 5, `vorher-nachher-…` 7, `test-name-…` 4, `gemeldete-ungenauigkeit-…` 3, `plan-zusage-…` 3, `adapter-unittest-…` 3, `lesesperre-…` 1 |

**Übernommen** (nicht nachgemessen): die Blockdauer des Backfill-Runs von 0,12 bis 0,13 s bei
Blockgröße 1000 und die Kopierrate 4088 bis 9425 Zeilen/s aus `docs/user/benutzerhandbuch.md`
§Backfill (dort als übernommen aus Lauf-Berichten des Messwerkzeugs `tools/bench-backfill.sh`
geführt); die Stufe „groß“ von `SPEC-014` (1000 Changes/s).

---

## 2. (A) Zeitgrenze des Vorlaufs vor `stream.Run`

**Verdikt: die Entscheidung ist eine Umsetzung, kein akzeptiertes Negativ.**
[`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
trägt Begründung, Alternativen, Fitness Function und Trigger; hier steht die Beantwortung der
vier Teilfragen (1) bis (4) der Welle und die Einordnung.

### 2.1 Was die Messung gegen die Herleitung der Welle geändert hat

Die Welle (§5 Frage (e), Register `wartegrenze-ohne-zeitgrenze-im-startpfad`) führte zwei
Behauptungen als hergeleitet: ein wartender Vorlauf halte die Erfassung an, und er erscheine
gesund. Beide sind jetzt **gemessen** (M3): der Healthcheck meldet `0` über die ganze Wartezeit,
das Heartbeat-Alter bleibt unter 7,1 s (Schwelle 15 s), eine eingefügte Zeile ist nie über
`cdc.changes` sichtbar.

Die Messung fand außerdem etwas, das kein Plan, kein Review und keine Verifikation hatte: der
Prozess **stirbt** an dem Warten. `receive.NewStream` sendet `START_REPLICATION`
(`receive.go`, `session.StartReplication`) und liegt in `Run` vor dem gesamten Rest der
Verdrahtung und dem Vorlauf (`wiring.go:641` gegen `wiring.go:1041`); der Server beendet einen
Strom, der `wal_sender_timeout` lang nicht liest (PostgreSQL-Log, M3), und der Prozess bemerkt es
am ersten Schreiben nach dem Vorlauf („broken pipe“, Klasse `replication`, Heartbeat-Fehlerzustand
`replication`). Standard von `wal_sender_timeout` sind 60 s (Lauf C), die Compose-Umgebung setzt
2 s (Lauf A und B). Das ist der Grund, warum Unit-Tests der Sequenz (Fakes, Quelltext-Lesung) und
der E2E-Beleg mit einem Antrag in Millisekunden es nicht sahen: die Aufrufstelle in `Run` ist nur
über zwei Quelltext-Tests gebunden (Review F-4 und F-9 zu `slice-transformationen-start-reihenfolge`),
und kein Lauf hielt den Vorlauf länger als die Server-Einstellung an.

### 2.2 Die vier Teilfragen

1. **Frist.** Ja, **30 s je Prozessstart** (ein Durchlauf, nicht je Antrag). Wert hergeleitet aus
   der Fehlergrenze von 60 s des Capture-Abstands (`SPEC-013`), ohne Messung des Werts selbst;
   Konstante des Codes wie `heartbeatInterval`, keine Konfigurationsachse (Trigger für einen
   Konfigurationswert: Betreiber-Bericht, `ADR-0128` §Re-Evaluierungs-Trigger).
2. **Bei Ablauf.** Der Stream startet, der unterbrochene Antrag und jeder dahinter bleiben
   `pending`, die Administrations-Goroutine verarbeitet sie ohne Frist; kein Vermerk `failed`
   (Alternative E in der ADR: ein gesunder, langsamer Antrag würde verworfen, der Abhilfe-Antrag
   könnte verloren gehen). Wirkung auf die Abhilfe (c) aus `ADR-0112` Folgepflicht 5: sie hält für
   jeden Antrag, den der Vorlauf innerhalb der Frist erreicht (im Normalfall Millisekunden); geht
   ihm ein Antrag voraus, der länger als 30 s läuft, startet der Stream mit dem bisherigen
   Regelstand, und ist die Regel nicht anwendbar, endet der Prozess wieder mit `schema` — sichtbar,
   ohne Datenverlust, der Betreiber startet neu, sobald der Antrag `applied` ist.
3. **Sichtbarkeit.** Keine Anzeige in `diagnose` und `--healthcheck`: eine Wartezeit von höchstens
   30 s ist kein Fehlerzustand im Sinn von `LH-FA-ADM-003`; der Träger ist ein Warn-Eintrag im Log
   beim Ablauf. Die Aussage „gesund während des Wartens“ steht damit **gemessen** in der ADR
   (M3), der Wert der Frist selbst ist erwartet, bis der Rundlauf des Slice läuft.
4. **Handbuch.** `slice-transformationen-betriebsdoku` §2 führt „ein hängender Antrag hält den
   Stream-Start höchstens 30 s an“ mit der Messung des Rundlaufs, den der Slice liefert (§9).

### 2.3 Warum nicht kleiner

Eine Frist allein (Alternative B) beseitigt den Absturz nicht: sie müsste unter
`wal_sender_timeout` liegen, dessen Wert der Quelle gehört und im Compose-Pin 2 s beträgt
(Prozessstart gegen eine Instanz mit 2 s: jeder Vorlauf über 2 s bräche). Ein `SHOW`-abgeleiteter Wert
(Alternative C) koppelt den Start an eine Server-Einstellung. Die Verschiebung von
`START_REPLICATION` macht beides überflüssig, ist an M4 erprobt (Kommando-Zustand) und im Diff
klein (§9).

**Nicht erprobt:** der Adapter mit verschobenem `START_REPLICATION` (der Kommando-Zustand ist mit
`psql` erprobt, PostgreSQL 17 nicht gefahren), der Ablauf der Frist am Prozess (kein Lauf hat
eine Frist), eine andere Antragsart als `enable`.

---

## 3. (B) Lesekosten des Backfill-Runs

**Verdikt: akzeptiertes Negativ mit Messung und benanntem Trigger; kein Folge-Slice.**

Der Run liest den Regelstand einmal zu Beginn, je Block und vor dem Commit (`Blockzahl + 2`
Lesungen) und den Ausschlussstand je Block und vor dem Commit (`Blockzahl + 1`); jede Lesung liest
die `applied`-Zeilen der Antragsarten aller Tabellen der Quelle (`copyBlocks`, Statements
`SelectAppliedTransformationRequests` und Ausschluss, ohne Tabellenfilter; beide am Stand gelesen).
Die Kosten je Lesung sind die SQL-Lesung (M2) plus die Faltung in Go (M1):

| Zeilen der Queue (gesamt / gelesen) | SQL je Lesung (M2, gemessen) | Faltung der Regeln (abgeleitet aus M1: 0,91 bis 0,97 µs je Antrag) | Lesungen je Block (Regelstand, Ausschluss) | Anteil einer Blockdauer von 0,12 bis 0,13 s (abgeleitet) |
|---|---|---|---|---|
| 100 / 30 | 0,26 bis 0,31 ms | ≈ 0,03 ms | 2 | ≈ 0,5 % |
| 10 000 / 3000 | 2,9 bis 3,8 ms | ≈ 2,9 ms | 2 | ≈ 8 % |
| 100 000 / 30 000 | 28,0 bis 29,8 ms | ≈ 28 ms | 2 | ≈ 70 % |
| 1 000 000 / 300 000 | 186 bis 201 ms | ≈ 282 ms | 2 | ≈ 530 % |

*Rechnung (abgeleitet):* je Block = SQL des Regelstands + Faltung + SQL des Ausschlusses (dessen
Faltung ist kleiner und nicht gerechnet), also bei 10 000 Zeilen ≈ 3,4 + 2,9 + 3,4 = 9,7 ms
gegen 125 ms. Die Faltung bei 30 000 und 300 000 gelesenen Anträgen ist aus dem gemessenen
Wert bei 100 000 Anträgen (94 ms) **linear hochgerechnet**, nicht gemessen; die SQL-Werte sind an
PostgreSQL 18 im Container ohne Netz gemessen, ohne Index (die Tabelle trägt nur den
Primärschlüssel, wie im Schema).

**Aussage, die trägt:** die Kosten je Lesung sind **klein bis rund 10 000 Zeilen** der Queue
einer Quelle (≈ 8 % einer Blockdauer bei zwei Lesungen je Block). Eine Queue **realer** Größe gibt
es nicht: kein Server-Tag trägt die Transformations-Antragsarten (M7), die größte Queue der
Testumgebungen liegt unter 100 Zeilen (Größenordnung der Läufe von `make test-integration`,
hergeleitet, nicht gezählt). Die Messung ist deshalb synthetisch, und die Aussage endet dort, wo
sie gemessen ist.

**Trigger** (beobachtbar): eine Quelle trägt mehr als 10 000 Zeilen in
`cdc.administration_request` (`SELECT count(*) FROM cdc.administration_request WHERE source_id =
…`), oder ein Backfill-Run bleibt bei größerer Queue erkennbar unter der Richtgröße von
`docs/user/benutzerhandbuch.md`, oder `DefaultBlockSize` sinkt unter 100 (dann steigt die Zahl
der Lesungen um das Zehnfache). **Behebung dann:** ein Index auf `(source_id, status,
request_kind)` **und** ein tabellenbezogener Lesezugriff — ein Filter im Port allein
(`TransformationRules(ctx, source, table)`) spart nur den Transfer und die Faltung der vier
anderen Tabellen, nicht den Scan der Tabelle (M2 misst den Scan). Die Zeilen der Queue lassen sich
nicht archivieren: die `applied`-Zeilen sind die Herkunft des Regelstands und des
Ausschlussstands (`ADR-0065`, `ADR-0112` Teilfrage 6). Die Nebenwirkung derselben Größe auf den
Fallback-Poll (`ListPending` alle 5 s liest dieselbe Tabelle) ist nicht Gegenstand dieses Zugs;
bei einer Queue von 1 000 000 Zeilen läse eine Lesung 186 bis 203 ms (M2) alle 5 s, hergeleitet
und Teil derselben Behebung.

**Träger der Aussage:** dieses Verdikt (§3) und die Zeile in `welle-transformationen` §3;
`slice-transformationen-betriebsdoku` §2 nennt den Betreiber-Hinweis („die Zahl der Antrags-Zeilen
ist nicht begrenzt; ab rund 10 000 Zeilen einer Quelle kostet die Lesung je Block messbar“, mit
den Zahlen dieses Abschnitts und ihrem Ursprung).

---

## 4. (C) Restfläche der Zustellwege

**Verdikt: kein Folge-Slice; die Aussage bleibt „hergeleitet“ mit der gemessenen Menge.**

**Gemessene Menge** (Abdeckungstabelle `docs/user/e2e-abdeckung.md`, Runner-Phase
„Transformationen-Happy-Path (fünf Zustellwege)“): eine per `cdc.set_transformation` beantragte
`rename_column`- und `map_value`-Regel prägt eine danach **eingefügte Zeile** (INSERT, Neu-Bild)
auf den fünf Wegen; UPDATE, DELETE und das Alt-Bild sind über `cdc.changes` belegt
(`TestE2ETransformationRulesShapeBothImages`).

**Herleitung für den Rest** (zwei Glieder, beide gelesen): (1) die Regel wirkt im `Assembler`
auf beide Bilder jeder Operation (belegt über `cdc.changes`); (2) die drei Live-Adapter und der
HTTP-Lesezugriff bilden die Bilder einer Change unverändert ab —
`toProtoChange` (`internal/adapters/driving/grpc/server.go`, `OldImage: change.OldImage`),
`toStreamChange` im SSE-Adapter (`http/sse.go`, `rowImage(change.OldImage)`) und im
NATS-Publisher (`natsstream/publisher.go`, `rowImage(change.OldImage)`); der Live-Weg reicht
dasselbe `model.Change` weiter, das der `CaptureService` persistiert (`capture/service.go`,
`s.stream.Publish(ctx, &changes[i])`). Ein Unit-Test mit einem Alt-Bild, das nicht `null` ist,
steht für gRPC (`server_test.go`), HTTP-Lesen (`readchanges_test.go`) und SSE (`sse_test.go`);
für den NATS-Publisher steht er nicht (`publisher_test.go` trägt nur den Fall
`OldImage = nil`): dort ist das Glied (2) nur **gelesen**.

**Warum kein Slice.** Ein Folge-Slice trüge einen realen Lauf pro Weg und Operation; das Risiko,
das er abdeckt, ist ein Adapter, der ein Bild verändert — ein Fehler, den ein Unit-Test der
Abbildung fängt, wo er entsteht, und der in einer Abbildung ohne Zustand keinen Ort hat, an dem
die Regel eine Rolle spielt. Kein Betreiber-Bericht und kein Beleg im Register sprechen dagegen.

**Trigger** für einen realen Lauf: eine Änderung an einer der vier Abbildungsstellen, ein neuer
Zustellweg, oder ein Bericht über ein abweichendes Alt-Bild an einem Weg. Ein Unit-Test mit
nicht-leerem Alt-Bild am NATS-Publisher ist die kleinste Erweiterung und gehört in den ersten
Slice, der `natsstream/` berührt (Adresse: Trigger, kein eigener Slice).

**`slice-sdk-regel-realserver-e2e`** trägt die Fläche nicht und soll sie nicht tragen: er prüft die
Modelle der drei SDK-Clients (Schlüssel opak) mit einer Change, die der Server erzeugt; das
Alt-Bild einer UPDATE-Change verlangte eine Operation mehr je Tier-Phase, ohne etwas über die
Modelle auszusagen, das die INSERT-Change nicht schon aussagt.

**Träger, die etwas sagen:** die Abdeckungstabelle (Erzeugnis, sagt die gemessene Menge exakt);
`slice-transformationen-e2e-wirkung` §3 (Record, benennt die Restfläche als hergeleitet);
`welle-transformationen` §3 (dieses Verdikt); `slice-transformationen-betriebsdoku` §2 — das
Handbuch nennt „alle Wege dieselbe Form“ **nur** mit der gemessenen Menge und kennzeichnet
UPDATE/DELETE/Alt-Bild auf den vier weiteren Wegen als hergeleitet.

---

## 5. (D) `AGENTS.md` §3.1 und flaglose Schreibwege

**Verdikt: der Wortlaut wird geschärft; die Entscheidung ist eindeutig.**

Der Absatz „Text-Umschreiben im Repo ist Sache der Datei-Werkzeuge des Laufs“ trägt seine Regel
schon in der Überschrift und im zweiten Satz („Die Änderung läuft über Edit/Write des Laufs oder
über ein Repo-Werkzeug hinter `make`“); der Verbotssatz zählte nur die in-place Formen auf. Die
Lücke war die Aufzählung, nicht die Regel: `cat >> Datei` erfüllt keinen der beiden erlaubten Wege
(Beleg: `evidence/slice-transformationen-map-value.md` in
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`). Kein Weg außer Edit/Write und dem
Repo-Werkzeug verlangt eine benannte Ausnahme: der Bestand kennt nur zwei legitime Umleitungen,
beide ohne Text einer Repo-Datei zu schreiben — die Umleitung in den Scratchpad oder das
Temp-Verzeichnis (bereits ausdrücklich erlaubt: `sed … Datei > Kopie`) und das Kopieren und
Verschieben ganzer Dateien (`cp` der Vorlage nach `AGENTS.md` §1, `git mv`, `git checkout`).

**Diff-Text** (umgesetzt in `AGENTS.md` §3.1):

```text
- Datei —, ist verboten, auch für „nur einen Schnelltest“. Die Änderung läuft
- über Edit/Write des Laufs oder über ein Repo-Werkzeug hinter `make`. Eine
+ Datei —, ist verboten, auch für „nur einen Schnelltest“. **Dasselbe gilt für
+ das Schreiben und Anhängen von Text an eine Repo-Datei durch eine Umleitung**
+ (`> Datei`, `>> Datei`, `cat <<EOF > Datei`, `tee`, `dd of=`): der Text einer
+ Repo-Datei entsteht über Edit/Write des Laufs oder über ein Repo-Werkzeug hinter
+ `make`. Ganze Dateien verschieben oder kopieren (`git mv`, `cp` einer Vorlage,
+ `git checkout`) schreibt keinen Text und bleibt erlaubt, ebenso eine Umleitung in
+ eine Datei im Scratchpad oder Temp-Verzeichnis. Eine
```

Dazu drei Folgeänderungen im selben Zug: das Beispiel „Falsch“ trägt `cat >> docs/x.md`; die
„Begründung“ und die „Durchsetzung“ nennen das Umleitungs-Verbot; der Reviewer-Skill-Punkt
„Docker-only-Verstoß“ und die Docker-only-Zeile in `implement-slice` nennen den Weg (der Reviewer
ist der Wächter, weil der Guard Umleitungen nicht liest).

**Kein Guard-Ausbau.** `MR-003` (Obergrenze der Quote-Lesung) hält Heredocs und Umleitungen als
Sache einer Sandbox-Ausführung fest; eine Zeichenketten-Regel gegen `> Repo-Pfad` blockte zugleich
`make gates > /tmp/log` und jeden Scratchpad-Lauf und ließe `tee` und jeden Umweg über Variablen
offen. Der Auflösungs-Trigger von `MR-003` (eine weitere Beleg-Datei mit Wirkung auf eine
Repo-Datei) bleibt der Weg dorthin. `tools/harness/blocked/go` bleibt eine Frage des Nutzers und
ist nicht berührt.

---

## 6. (E) Obergrenze der Paare von `map_value` (`SPEC-030`)

**Verdikt: benannter Verzicht auf eine Obergrenze; keine Änderung am Pflichtenheft.**

**Nachgemessen** (M1; die Verifikation von `slice-transformationen-map-value` maß 80 ns, 6,8 µs und
0,72 ms bei 10, 1000 und 100 000 Paaren, übernommen; die Nachmessung bestätigt die Größenordnung):
die Suche einer Zuordnung, wenn der Schlüssel fehlt (ungünstigster Fall), kostet je Wert und Regel 54
bis 123 ns bei 10 Paaren, 6,5 bis 6,7 µs bei 1000, 65 bis 80 µs bei 10 000 und 0,66 bis 0,70 ms bei
100 000. Das Lesen der Regel (Faltung, je Prozessstart und je Lesung des Backfills) kostet bei 1000
Paaren 0,71 bis 1,07 ms, bei 100 000 Paaren 82 bis 83 ms.

**Abgeleitete Wirkung** (Stufe „groß“, `SPEC-014`: 1000 Changes/s, im Stream einer Quelle in einer
Goroutine; jede Change trägt im ungünstigen Fall zwei Bilder mit dem Wert der Regelspalte, also zwei
Suchen): bei 1000 Paaren 13 ms je Sekunde (1,3 % eines Kerns), bei 10 000 Paaren 130 bis 160 ms
(13 bis 16 %), bei 100 000 Paaren 1,3 bis 1,4 s je Sekunde — die Erfassung dieser Quelle hielte die
Stufe „groß“ nicht. Für die Stufen „klein“ und „mittel“ (höchstens 100 Changes/s) bleibt es bei
höchstens 14 % auch bei 100 000 Paaren.

**Warum kein Verbot.** Der Fehler ist selbstgemacht (`cdc_admin` legt die Regel an), im Betrieb
sichtbar (`cdc_capture_lag` wächst über die Warn-Grenze von 5 s, `SPEC-013`), und mit
`cdc.remove_transformation` in Sekunden behoben; eine Zuordnung von 100 000 Paaren ist ein
Verweis-Datensatz und keine Transformation. Eine Obergrenze im Pflichtenheft verlangte einen neuen
Fehlertext in `SPEC-019`, eine Prüfung in `ParseTransformationSpec`, Tests und Handbuch-Zeilen — ein
Slice, dessen Nutzen ein Bericht aus dem Betrieb noch nie angefordert hat.

**Trigger:** ein Betreiber meldet ein Wachstum von `cdc_capture_lag` mit einer aktiven
`map_value`-Regel, oder eine Regel mit mehr als 10 000 Paaren steht in der Queue. **Dann:**
Folge-ADR mit einer Obergrenze (Vorschlag 10 000 Paare: die abgeleitete Last der Stufe „groß“ liegt
dort bei höchstens 16 % eines Kerns) samt Fehlertext in `SPEC-019` und `SPEC-030`.

**Träger:** dieses Verdikt und `slice-transformationen-betriebsdoku` §2 („Größenordnung der
Suche“): das Handbuch nennt **diese** Zahlen mit Messbedingungen (Instanz und Ursprung wie in M1) und
führt keine Obergrenze. Die Übernahme-Kennzeichnung im Plan von `betriebsdoku` („übernommen, nicht
nachgemessen“) fällt damit weg.

---

## 7. (F) Konjunktiv-Wortliste in `implement-slice` Schritt 20

**Verdikt: erweitert, weil messbar besser.**

**Befund.** Die Wortliste `wäre|waere|würde|wuerde|hielte|hätte|haette` (mit `sonst|statt`) traf
den Fund F-1 von `slice-transformationen-e2e-wirkung` („trüge“) nicht: 0 Zeilen am Diff `c246ba4f`
bis `66f60c8b` (M6), obwohl die Zeile darin steht.

**Messung.** Ergänzt um `trüge`, `bliebe`, `ließe`, `könnte`, `müsste`, `bräuchte`, `läge`,
`stünde`, `käme`, `wären`, `gäbe`, `ginge`, `fände`, `dürfte`, `hieße`, `brächte` (jeweils auch in
der transliterierten Form) **mit Wortgrenze `\b`** — ohne sie träfen `bliebe` „geblieben“, `ließe`
„schließen“ und `stünde` „entstünde“ (gemessen: ohne `\b` 28 statt 20 Zeilen für `bliebe`, 34 statt
6 für `ließe`, 5 statt 1 für `stünde`): am Diff des Fundes 1 statt 0 Zeile, am Bestand 359 statt 322
(M5; ohne `sonst`/`statt` 103 statt 62). Der Mehrbedarf des Kandidatenlaufs ist **+37 Zeilen am
Bestand (+11,5 %)**; er ist klein gegen die 271 Zeilen der beiden Wörter `sonst` und `statt`, die
im Lauf bleiben.

**Falsch-Positive** (gelesen: die 42 Bestandszeilen mit einer der neuen Formen, fast alle
Mutationsbeschreibungen in Test-Godocs: „ohne X bliebe Y“, „sonst trüge die Messung …“, eine
Lese-Einschätzung, keine Zählung): sie sind nach dem Wortlaut von Schritt 20 **zulässig**
(„Mutationsbeschreibungen in Test-Godocs“). Der Nutzen ist der Recall einer bekannten, HIGH
eingestuften Klasse (7×); die Kosten sind einige Urteile je Slice. Die gemessene Abwägung geht
zugunsten der Erweiterung aus, aber knapp; der Zusatz in Schritt 20 nennt die Mutationsbeschreibung
ausdrücklich als zulässigen Treffer.

**Umgesetzt** in `.claude/commands/implement-slice.md` Schritt 20: die Befehls-Zeile, ein Satz mit
den zwei Zahlen und ein Zusatz zu den zulässigen Treffern. Der Reviewer-Skill trägt keine
Wortliste (nur den Beispielsatz „würde diese verzögern“); keine Änderung dort. Eine Grenze bleibt: eine
Konjunktiv-Form außerhalb der Liste trifft der Lauf nicht; der Reviewer ist die tragende Linie
(Schritt 20, Grenze der Selbstprüfung).

---

## 8. (G) Register-Einträge bei 3× oder darüber, dem Lese-Schritt der Welle zugeordnet

Zähler aus M8. Vier Einträge, je mit Ausgang.

| Eintrag | Zähler | Verdikt | Träger / Adresse |
|---|---|---|---|
| `test-name-behauptet-mehr-als-der-test-treibt` | 4× | **gestrichen** — akzeptiertes Negativ | Alle vier Belege hat der Reviewer (oder der Verifier) vor dem Merge gefunden, Schwere LOW; die Eskalation „Wiederholung eines Musters, das schon zweimal LOW war“ steht als MEDIUM-Punkt im Reviewer-Skill und wirkt beim nächsten Auftreten. Ein eigener Skill-Punkt kostete jedem Review-Lauf Kontext für eine Klasse, die nachweislich ohne ihn gefunden wird. **Trigger:** ein Auftreten nach dem Merge oder mit Schwere ≥ MEDIUM |
| `gemeldete-ungenauigkeit-ohne-traeger` | 3× | **verkörpert** — eine Zeile im Closure-Note-Reviewer-Skill | `.harness/skills/closure-note-reviewer.md`, MEDIUM-Zeile „Ungenauigkeit oder Grenze als „gemeldet“ geführt, ohne eine Adresse, die eintreten kann“ (umgesetzt in diesem Zug). Die drei Funde des Eintrags sind beantwortet: der Kommentar an der gRPC-Nachricht in `proto/cdc/stream/v1/changestream.proto` nennt heute „without the commit position, commit time and origin“ (am Baum gelesen, Datei zuletzt geändert mit `b522980f`), der Kommentar am Feld `allowDestructive` in `tools/schema/rolloutguard/guard.go` nennt die Klasse „View-Signatur“ (Zeilen 59 bis 62, am Baum gelesen), die lineare Suche hat ihre Antwort in §6, die Restfläche in §4 |
| `plan-zusage-erfuellung-ohne-committeten-anker` | 3× | **gestrichen** — akzeptiertes Negativ | Drei Belege, alle LOW, alle vom Reviewer vor dem Merge gefunden; die Behebung ist jedes Mal derselbe Zug (der Plan trägt den Verifikations-Report als Anker samt Reichweite) und liegt in der Rolle, die den Haken setzt (Planner-Closure, Nachbar `dod-checkbox-nachzug`, verkörpert in `implement-slice` Schritt 18 und 21). Eine weitere Zeile im Implementer-Ablauf verlangte von der Rolle einen Anker, den erst der Verifier erzeugt. **Trigger:** ein Auftreten, dessen Haken bis in `done/` ohne Anker steht |
| `adapter-unittest-verdeckt-bootstrap-luecke` | 3× (4× mit diesem Zug) | **geplant** → `slice-start-vorlauf-grenze` | Das vierte Auftreten ist die Messung dieses Zuges: die Sequenz aus `Run` hat Unit-Tests mit Fakes und zwei Quelltext-Tests, und die Eigenschaft, die sie nicht zeigen (der Strom läuft schon, während der Vorlauf wartet), erschien erst am komponierten Prozess. Der vorgeschlagene Träger des Eintrags (ein Smoke-Test der Aufrufstelle) hätte sie nicht gefunden — er hätte die Datenbank gebraucht wie `Run` selbst und keinen Zeitverlauf gefahren (hergeleitet). Der Träger der Klasse ist die Phase in `make test-integration` **in dem Slice, der die Eigenschaft einführt** (Testpyramide, `ADR-0030`); für diesen Fall liefert sie `slice-start-vorlauf-grenze` |

Die Ausgangs-Sätze stehen in den `state.md` der vier Einträge; für zwei Einträge kommt eine
Beleg-Datei hinzu (das vierte Auftreten von `adapter-unittest-…`, das zweite von
`wartegrenze-…`, §9).

**Nicht gelesen** (nicht zugeordnet): die übrigen Einträge bei 3× oder darüber. Der Lese-Schritt der
Closure von `welle-transformationen` liest sie; er bleibt ein Kriterium in §3 der Welle.

---

## 9. Übergaben an den Planner (Träger-Nachzüge und Slice-Vorschlag)

### 9.1 Slice-Vorschlag `slice-start-vorlauf-grenze`

- **Kennung** (MR-002: ein Name): `slice-start-vorlauf-grenze`; **wellenlos** (keine Closure-Bedingung,
  die von seiner DoD verschieden wäre; das Closure-Kriterium der Welle zur Zeitgrenze ist mit der
  Beauftragung dieses Slice erfüllt, nicht mit seiner Lieferung). Bezug:
  [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  Folgepflichten 1 bis 5, [`LH-FA-CFG-007`](../../spec/lastenheft.md),
  [`LH-QA-REL-001`](../../spec/lastenheft.md).
- **Umfang** (Vorschlag für die DoD, jeder Punkt mit „zu belegen durch“):
  1. `receive.NewStream` sendet kein `START_REPLICATION` mehr; `Stream.Run` sendet es als erste
     Handlung; der Test im Replication-Tier (Instanz mit `wal_sender_timeout` 2 s: `NewStream`,
     länger warten, `Run`, eine danach committete Änderung wird geliefert) läuft an PostgreSQL 18
     **und** 17 (`PG_TEST_IMAGE` auf den Digest aus `e2e.yml`, wie in
     `slice-capture-leerlauf-quellbelege`); Mutation: `StartReplication` zurück in `NewStream`.
  2. `runStreamAfterAdministrationPass` trägt die Frist von 30 s; bei Ablauf: Warn-Eintrag,
     kein `MarkFailed` für den unterbrochenen und die folgenden Anträge (auch nicht mit dem beendeten
     Kontext), Goroutine, Stream; Whitebox-Test mit verkürzter Frist und einer Queue, die bis zum
     Kontext-Ende blockiert; Mutation: Vorlauf ohne Frist, `MarkFailed` bei Ablauf.
  3. Eine Phase im Runner von `make test-integration` (Sperre an einer eigenen Tabelle länger als
     die Frist, `pending`-Antrag `enable`, Neustart des Feed-Containers, Healthcheck `0`, eine
     währenddessen committete Änderung ist binnen Frist plus Toleranz sichtbar, Antrag nach
     Freigabe `applied`, Prozess ohne Fehlerklasse; Laufzeit der Phase rund 45 s, hergeleitet aus
     der Frist); Deklarations-Anker mit `LH-FA-CFG-007` und `-run`-Abgleich.
  4. Pflichtenheft: der Absatz „Abhilfe (Zusage)“ von `LH-FA-CFG-007.a` nennt die Ordnung innerhalb
     der Frist (`ADR-0128` Folgepflicht 2).
  5. §3.13-Suchlauf nach `START_REPLICATION`/`NewStream`/„Vorlauf“ (Träger nach `ADR-0128`
     Folgepflicht 5: Godoc von `NewStream`, Kommentar in `wiring.go` „Der Slot besteht an dieser
     Stelle bereits“, Godoc von `runStreamAfterAdministrationPass`, dessen Satz „Der Vorlauf trägt
     keine eigene Frist“ dann falsch ist).
  6. `harness/README.md` §Sensors (`make test-integration`, `make test-replication`): der Beleg in
     der Aufzählung.
- **Start-Trigger** (Vorab-Bedingung): [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  `Accepted` (erfüllt, mit diesem Zug), `slice-capture-leerlauf-quellbelege` in `done/` (beide
  Slices erweitern `tools/harness/run-integration-tests.sh` und den Tier `make test-replication`;
  WIP-Limit 1), kein Slice in `in-progress/`.
- **Einordnung und Reihenfolge:** `slice-capture-leerlauf-quellbelege` → `slice-start-vorlauf-grenze` →
  `slice-transformationen-e2e-abhilfe` (dessen Start-Trigger um „`slice-start-vorlauf-grenze` in
  `done/`“ zu erweitern ist: der Slice belegt die Abhilfe am Startpfad, den diese Lieferung ändert;
  sein Risiko „Der Vorlauf trägt keine Zeitgrenze“ (§6) und die Zeile im Bereich „Berührte Stellen“
  ändern sich mit ihr). `slice-transformationen-betriebsdoku` bleibt der letzte Slice der Welle.
- **Risiken des Slice** (für den Plan): `Stream.Run` sendet `START_REPLICATION` erstmals selbst —
  Fakes der Tests in `receive` und der ACK-Adapter, der dieselbe Verbindung teilt, sind zu prüfen
  (`ADR-0128` §Konsequenzen, „Erwartet, nicht am Code belegt“); die Phase verlängert
  `make test-integration` um rund 45 s je Leg der CI-Matrix (`e2e.yml`, kein Workflow-Zug,
  `AGENTS.md` §3.10 greift nicht).

### 9.2 Nachzüge, die dieser Zug selbst trägt

Umgesetzt in diesem Zug: [`ADR-0128`](../plan/adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
samt ADR-Index; `AGENTS.md` §3.1; `.harness/skills/reviewer.md` (eine Klausel „Docker-only-Verstoß“);
`.harness/skills/closure-note-reviewer.md` (eine MEDIUM-Zeile); `.claude/commands/implement-slice.md`
(zwei Stellen); `welle-transformationen` (§3 drei Kriterien, §5 Fragen und Kante); die `state.md`
der betroffenen Register-Einträge und zwei Beleg-Dateien.

### 9.3 Nachzüge, die dem Planner gehören

| Träger | Zug | Frist |
|---|---|---|
| `slice-transformationen-e2e-abhilfe` §4 (Start-Trigger), §6 (Risiko „Der Vorlauf trägt keine Zeitgrenze“), §1 („NICHT“-Punkt zur Ordnung) | die Kante zu `slice-start-vorlauf-grenze` und den Stand der Entscheidung (Frist 30 s, Ablauf, Wartesicht) nachtragen; das Risiko wird zu „die Frist ändert den Startpfad, den dieser Slice belegt“ | vor dem Start von `e2e-abhilfe` |
| `slice-transformationen-betriebsdoku` §2 | fünf Aussagen mit ihrem Ursprung: (1) hängender Antrag hält den Stream-Start höchstens 30 s an (Messung des neuen Rundlaufs, nicht Herleitung), (2) Größenordnung der Suche einer Zuordnung mit den Zahlen aus §6 (nachgemessen, kein „übernommen“), (3) Kosten der Lesung im Backfill mit den Zahlen und der Grenze von 10 000 Zeilen aus §3, (4) „alle Wege dieselbe Form“ nur mit der gemessenen Menge, der Rest als hergeleitet (§4), (5) Bezugspunkt des Regelstands im Run (schon dort geführt) | Start von `betriebsdoku` |
| `docs/plan/planning/in-progress/roadmap.md` | `slice-start-vorlauf-grenze` bei *Offene Slices* führen (wellenlos, Kante vor `e2e-abhilfe`), Eintrag im Drift-Log der Umplanung | mit der Anlage des Slice |
| Slice-Plan `slice-start-vorlauf-grenze` | aus §9.1 anlegen (`cp` aus der Vorlage, `open/`) | vor dem Start von `slice-transformationen-e2e-abhilfe` |

Die Beleg-Dateien zu `wartegrenze-ohne-zeitgrenze-im-startpfad` (Messung dieses Zugs, 2×) und
`adapter-unittest-verdeckt-bootstrap-luecke` (4×) sind in diesem Zug angelegt; die Ausgänge stehen in
den `state.md`.

---

## 10. Fragen an den Nutzer

**Keine entscheidungspflichtige.** Die Frist von 30 s ist eine Konstante des Codes (Trigger für
einen Konfigurationswert: `ADR-0128`), die Zusage der Abhilfe wird um „innerhalb der Frist“ ergänzt
(Pflichtenheft ist fortschreibbar, Bezug ist die ADR). `tools/harness/blocked/go` (Host-Toolchain-Sperre)
ist unberührt; der Auflösungs-Trigger von `MR-003` ist unverändert.

Grenzen dieses Zugs, nicht verschwiegen: kein Lauf fuhr eine Antragsart außer `enable` in den
Warte-Läufen; PostgreSQL 17 ist an keiner der Messungen dieses Zugs gefahren (M4 gilt für 18); die
SQL-Kosten in §3 sind an einer synthetischen Queue ohne Netz gemessen; die Fitness-Function-Zeilen
von `ADR-0128` sind Erwartungen.
