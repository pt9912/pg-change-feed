# Architect-Verdikt — Intermittenz der Phase „Leerlauf-Bestätigung“ in `make test-integration`

**Datum:** 2026-09-27 · **Stand:** `4a4da1b8` (Baum sauber) · **Rolle:** Architect (frischer
Kontext) · **Anlass:** offene Frage im Register
`BEO-PGC/test-integration-retention-timing-flake` — die bestehende Phase „Leerlauf-Bestätigung“
war im `e2e.yml`-Lauf 36287009221 (Leg PostgreSQL 18, erster Versuch) einmal rot, der
Wiederholungsversuch grün. **Vollmacht:** Empfehlung statt Optionsliste, keine Rückfragen an
den Auftraggeber.

**Bezug:** [`LH-QA-REL-001`](../../spec/lastenheft.md),
[`LH-QA-REL-003`](../../spec/lastenheft.md),
[`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (Haupt-Bezug),
[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md),
[`SPEC-013`](../../spec/pflichtenheft.md),
[`AGENTS.md`](../../AGENTS.md) §3.5, §3.6, §3.12, §3.13.

---

## 0. Ergebnis in einem Blick

- **(1) Ursache: der Testaufbau, nicht das Produkt.** Die Anweisung des Runners — ein
  einzelnes `INSERT` mit 16,0 MB WAL, dem 1,9-fachen der Fehlerschwelle der Phase — ist ein
  **Stoß**: der Rückstand (WAL-Ende minus `confirmed_flush_lsn`) steht danach für einige
  zehn Millisekunden (lokal) bei der **ganzen Last**, bevor die Bestätigung folgt (gemessen
  am echten Stream, §1 M5). Ein Prüf-Takt (5 s), der in dieses Fenster fällt, beendet den
  Container. Das ist ein Rennen zwischen Stoß und Takt, kein Mangel der
  Leerlauf-Bestätigung. Die Entscheidung `ADR-0120` gilt; die Phase hat breiter behauptet,
  als der Mechanismus trägt (Konflikt-Pfad, erstes Verdikt).
- **(2) Ein kleiner Umsetzungs-Slice** `slice-leerlauf-phase-last-in-stuecken` (wellenlos,
  Umfang S — Schätzung): die Last der Phase in Stücken unter der Warnschwelle, nach jedem
  Stück die Wartebedingung „`confirmed_flush_lsn` hat das Stück erreicht“; die Summe bleibt
  über der Fehlerschwelle. **Erprobt** am echten Stream (§1 M6, lokal): höchster Rückstand
  32 % der Fehlerschwelle bei gleicher Gesamtlast. Direkt hinter
  `slice-capture-leerlauf-quellbelege`, **vor** `slice-wal-fehlerschwelle-ausgangsklasse`.
- **(3) Ein Rot dieser Phase ist bis zum Slice weder Beleg noch Widerlegung** der Phasen
  dahinter; nach dem Slice ist ein Rot mit derselben Signatur ein **Befund** und ein neuer
  Architect-Zug.
- **(4) Sofortmaßnahme:** zwei `gh`-Befehle im Register (§5), beide gefahren, kein Sensor,
  kein Gate.

Keine `Accepted`-ADR wird berührt, keine neue ADR geschrieben (Begründung §2.4). Zwei
Angaben im Register waren falsch und sind berichtigt (§7).

---

## 1. Was gelesen und was gemessen ist

**Gelesen** (Stand `4a4da1b8`):

| Nr | Fundstelle | Lesung |
|---|---|---|
| L1 | `compose.yaml`, Zeile 30 | die Compose-PostgreSQL läuft mit `wal_sender_timeout=2000` (2 s) |
| L2 | `internal/bootstrap/wiring.go`, `heartbeatInterval` (Zeile 150) und der Aufruf in Zeile 1038 | der Prüf-Takt des WAL-Rückstands ist 5 s |
| L3 | `wiring.go`, `runWALRetentionCheck` (Zeilen 1205–1233) | je Takt **eine** Messung; `bytes > errorBytes` setzt den Fehler, ruft `stopStream` und kehrt zurück — keine Wiederholung, keine Karenz |
| L4 | `receive/walretention.go`, `Measure` | Rückstand = `IDENTIFY_SYSTEM` (WAL-Ende) minus `confirmed_flush_lsn` |
| L5 | `receive/receive.go`, `handleCopyData`/`confirmIdle` | der Stream bestätigt nur auf eine Keepalive-Nachricht, nicht inmitten einer Quelltransaktion |
| L6 | `tools/harness/run-integration-tests.sh`, Phase „Leerlauf-Bestätigung“ (Zeilen 3379–3464) | Schwellen 4 MiB / 8 MiB; der Schreiber auf die nicht aktivierte Tabelle ist **eine** Anweisung `INSERT … generate_series(1, 60000)`; danach hält der Runner 12 s und prüft jede Sekunde, dass der Container läuft |

**Gemessen** (Ursprung je Zeile; „lokal“ = Entwicklungshost, nicht der CI-Runner):

| Nr | Messung | Instanz und Ursprung |
|---|---|---|
| M1 | Lauf 36287009221, Versuch 1, Leg PostgreSQL 18: Container gestartet 02:02:20,78 Z; `INSERT 0 60000` gedruckt 02:02:45,76 Z; „Lauf beendet mit Fehler“ 02:02:45,88 Z; Rückstand 15.238.216 B über 8.388.608 B; die Schritte 5 bis 7 des Jobs `skipped` | `gh api repos/pt9912/pg-change-feed/actions/jobs/108529548391/logs` und `…/jobs/108529548391`, 2026-09-27, selbst gelesen |
| M2 | Last des Schreibers 16.019.264 B (PostgreSQL 18, Job 108531740887) und 16.010.992 B (PostgreSQL 17, Job 108529548457); Backfill-Run 22.190.184 B und 22.378.528 B; höchste Probe des Containers im 5-s-Takt 64 B und 96 B | gedruckte Abschlusszeile der Phase „Leerlauf-Bestätigung … belegt“ beider Job-Logs, selbst gelesen. 16.019.264 / 8.388.608 = 1,91 (abgeleitet) |
| M3 | Häufigkeit: 40 abgeschlossene `e2e.yml`-Läufe auf Commits, die `00caec48` (Einführung der Phase) enthalten, je zwei Legs = **80** Ausführungen der Phase im ersten Versuch: 79 grün, 1 rot; ein weiterer Lauf (36290713488) lief noch. Von allen 444 Läufen von `e2e.yml` hat genau einer `run_attempt` > 1 | `git merge-base --is-ancestor 00caec48 <sha>` je Lauf und `gh api …/runs/<id>/attempts/1/jobs` je Lauf, 2026-09-27, selbst gefahren |
| M4 | Standard-Empfänger (`pg_recvlogical`) an einer Wegwerf-Instanz PostgreSQL 18.6 (Pin `PG_TEST_IMAGE`) mit `wal_sender_timeout=2000`, dieselbe Last 60.000 Zeilen, 6 Läufe, Proben alle 20 ms: der Rückstand steigt in Stufen (0,33 MB, 8,7 MB, 16,0 MB) auf die **ganze Last** (16.019.264 B in 5 Läufen, 16.019.288 B im sechsten) und fällt erst 0,19 bis 0,81 s nach dem Ende der Anweisung auf 0 | lokal, Wegwerf-Skript im Scratchpad, Container mit `docker rm -f -v` entfernt |
| M5 | **Echter `receive.Stream`** mit `WALRetentionChecker.Measure` daneben, dieselbe Instanzform, dieselbe Last, 5 Läufe, Proben alle 20 ms: Anweisung 144–153 ms; Spitze des Rückstands 16.019.288 B in 4 Läufen, 12.957.136 B im fünften; Fenster von der ersten Probe über 8 MiB bis zur letzten über 1 MiB **42, 42, 44, 63, 63 ms**; die letzte Probe über 1 MiB liegt **23, 47, 55, 65, 80 ms** nach dem Ende der Anweisung | lokal, Wegwerf-Go-Test im Paket `receive` (nicht committet, danach entfernt; `git status` sauber), gepinnte Images `TOOLCHAIN_IMAGE`/`PG_TEST_IMAGE` |
| M6 | Dieselbe Gesamtlast in **6 Stücken zu je 10.000 Zeilen**, nach jedem Stück Warten auf `confirmed_flush_lsn` ≥ WAL-Position hinter dem Stück, 3 Läufe: Gesamtlast 16.019.528 B, **höchster Rückstand 2.676.816 B** (32 % der Fehlerschwelle, 64 % der Warnschwelle), Wartezeit je Stück 11 bis 33 ms | lokal, derselbe Wegwerf-Go-Test, Proben alle 10 ms, echter Stream |

Nicht gemessen: die Dauer der Anweisung auf dem CI-Runner (der Runner druckt sie nicht), die
Latenz der Bestätigung dort, und das Verhalten anderer Nebenversionen als 18.6 (das
CI-Leg PostgreSQL 17 lief grün, mehr ist dazu nicht belegt).

---

## 2. Verdikt (1) — Ursache

### 2.1 Der Mechanismus

**Gemessen** (M4/M5): bei einem Stoß folgt der Rückstand dem WAL-Ende sofort — das WAL-Ende
springt in etwa 0,1 bis 0,2 s auf 16 MB, `confirmed_flush_lsn` steht noch am Anfang der Last —
und fällt erst nach dem Ende der Anweisung (M5: 23 bis 80 ms danach, lokal). Im Fenster
dazwischen steht er bei fast der ganzen Last, über der Fehlerschwelle, weil die Last sie um
das 1,9-fache übersteigt (M2).

**Hergeleitet** (Ablauf in der Quelle, nicht gelesen): die Bestätigung kann erst folgen, wenn
die Quelle den Commit der Anweisung gelesen, die 60.000 Änderungen gegen die Publication
verworfen und die Keepalive-Nachricht gesendet hat; der Stream beantwortet sie sofort. Er
antwortet also im Fenster nicht zu spät — ihm fehlt bis dahin, was zu bestätigen wäre.

Der gemessene Fehlerfall passt dazu: 15.238.216 B, 0,12 s nach der Rückkehr der Anweisung
(M1), das sind 95 % einer Last von etwa 16,0 MB (M2; die Last des roten Legs selbst ist nicht
gedruckt — *hergeleitet* aus den zwei grünen Legs mit 16,02 und 16,01 MB).

### 2.2 Das Rennen — Größenordnung (hergeleitet)

Der Takt läuft 5 s (L2); die Wahrscheinlichkeit, dass ein Takt in das Fenster W fällt, ist
W / 5 s. Mit dem lokal gemessenen W = 42 bis 63 ms (M5) sind das 0,8 bis 1,3 % je Ausführung
der Phase; bei 80 Ausführungen (M3) erwartet das Modell 0,7 bis 1,0 Treffer, bei einem
dreimal langsameren Runner 2 bis 3. Beobachtet ist **einer**. Das Modell ist mit der
Beobachtung **vereinbar**, mehr sagt es nicht: die Takt-Phase ist nicht gleichverteilt (der
Ablauf des Runners ist bis auf Jitter deterministisch), die Fensterbreite auf dem CI-Runner ist
nicht gemessen. Die Häufigkeit der Beobachtung ist 1 von 80 Ausführungen; das exakte
95-%-Intervall der Binomialverteilung ist etwa 0,03 % bis 6,8 % (*hergeleitet*, von Hand
gerechnet, gerundet); je Lauf mit zwei Legs 1 von 40, Intervall etwa 0,06 % bis 13,2 %. Die
Punktschätzung 1,25 % je Ausführung trägt keine Aussage über einen Trend; sie sagt: selten,
aber nicht null, und jedes Rot verzögert die Belege der Phasen dahinter.

### 2.3 Warum Testaufbau und nicht Produkt

1. **`ADR-0120` verspricht das Fenster nicht weg.** Die Zusage ist: der Rückstand von WAL
   ohne Inhalt für die Publication **sinkt** im Leerlauf (Form X1 im Tier `test-replication`
   misst **nach** der Last, nicht während ihr). Dass ein Container einen Stoß über der Schwelle
   ohne Rot übersteht, ist eine Aussage der **Runner-Phase**, nicht der ADR.
2. **Der Stoß ist das Konstrukt des Runners.** Ein `INSERT … SELECT generate_series` schreibt
   16 MB in etwa 0,15 s (M5, lokal), also rund 100 MB/s. Ein Schreiber des Betriebs liegt
   darunter, und die Bestätigung folgt ihm mit: im Handbuch stehen neun Backfill-Runs (je drei
   mit 10.000, 50.000 und 200.000 Zeilen), in denen die Spitze des Rückstands bei 0 MiB lag
   (Lauf `20260925T032925Z` **übernommen**, `20260925T043056Z` im Review-Report von
   `slice-backfill-slot-leerlauf-bestaetigung` gedruckt und mit gleichem Ergebnis
   nachgemessen; im 1-bis-2-s-Takt gelesen und auf ganze MiB gerundet, also grob). Der
   Rückstand ist im Bestätigungs-Betrieb höchstens `max(Stoß, Rate · Latenz)` (*hergeleitet*
   aus M5/M6). Die Phase wählt einen Stoß über der Schwelle und macht den Ausgang damit zu
   einer Frage der Takt-Phase.
3. **Produkt-Schwelle und Test-Schwelle sind verschieden groß gegen dieselbe Rate.** Die
   Fehlerschwelle des Betriebs ist 1 GiB ([`SPEC-013`](../../spec/pflichtenheft.md)); die
   Phase senkt sie auf 8 MiB, damit die Last klein bleibt. `Rate · Latenz` bleibt gleich;
   bei 8 MiB liegt es dicht an der Schwelle (M5: 100 MB/s · 0,05 s ≈ 5 MB), bei 1 GiB müsste
   ein Schreiber mehr als 1 GiB innerhalb einer Bestätigungs-Runde schreiben (Größenordnung
   10 GiB/s bei einer Runde von 0,1 s — *hergeleitet*, nicht gemessen). Für die Schwelle des
   Betriebs ist das Fenster **kein** praktisches Risiko; es ist ein Artefakt der gesenkten
   Test-Schwelle.
4. **Ein Produkt-Eingriff wäre die teurere Lösung** (§2.5) und behöbe keinen Betreiber-Fall.

### 2.4 Reichweite von `ADR-0120` — Präzisierung, keine Folge-ADR

`ADR-0120` §Festlegung 2 beschreibt die Metrik als „das WAL, das die Quelle dem Feed
geliefert, der Feed aber noch nicht bestätigt hat“; §Konsequenzen sagt, dasselbe gelte für
jeden Schreiber auf Tabellen ohne Publication-Bezug. Beide Sätze sind für **anhaltende** Last
richtig (M4/M5 nach dem Stoß; die neun Runs oben). Sie decken den **Stoß** nicht: die Metrik
enthält auch WAL, das die Quelle noch **nicht** geliefert hat (M5: Spitze = ganze Last, bevor
der Commit gelesen ist). Die Präzisierung steht hier und wird vom Slice an den Trägern geführt
(§3, Punkt 2). Eine Folge-ADR mit `Supersedes` ist nicht nötig: keine Entscheidung ändert sich
(weder Schwelle noch Code noch Spec-Wortlaut), die zwei Sätze sind beschreibend, und der einzige
Träger, der sie als Zusage liest, ist das Handbuch (ein Satz, Slice Punkt 2). Die Abkürzung
trägt, weil eine Folge-ADR nur Text gegen Text tauschte, ohne dass ein Leser anders handelte;
die Kette Register → Verdikt führt zu ihr.

### 2.5 Verworfen: Karenz der Fehlerschwelle im Produkt

Eine Karenz (die Fehlerschwelle gilt erst bei zwei aufeinanderfolgenden Messungen) beseitigte
das Rot in jedem Fall, in dem ein Stoß kürzer als ein Takt nachwirkt. Sie ist eine
**Lockerung des Gates** ([`AGENTS.md`](../../AGENTS.md) §3.6): eine ADR mit teilweisem
`Supersedes` der Abbruch-Aussage von `ADR-0049`, ein Nachzug in `SPEC-013`/`SPEC-008` und im
Handbuch (die Fehlerschwelle beendet dann erst nach einem zweiten Takt) und Änderungen an den
Tests der Schwellen-Prüfung — für ein Fenster, das bei der Schwelle des Betriebs nicht
vorkommt (2.3, Punkt 3). Der Aufwand steht in keinem Verhältnis zum Nutzen. Neubewertung nur
über den Trigger in 2.6.

### 2.6 Akzeptiertes Negativ, mit Trigger

**Bleibt:** ein Betreiber mit einer **gesenkten** `wal_retention_error_bytes` (Override der
Konfigurationsdatei) unter dem WAL, das seine Quelle in einer Bestätigungs-Runde schreibt,
kann den Container an einem Stoß verlieren. Der Weg ist der, den das Handbuch schon nennt: die
Fehlerschwelle über den erwarteten Rückstand heben. Ein Satz im Handbuch macht die Regel
„Schwelle deutlich über dem, was die Quelle in wenigen Sekunden schreibt“ lesbar (Slice,
Punkt 2).

**Trigger** (Neubewertung durch einen Architect-Zug): ein Rot mit der Signatur „Fehlerklasse
`replication` … WAL-Rückstand … über Fehlerschwelle“ in einer Phase des Runners, **in der die
Persistierung nicht gehalten wird** (also nicht in „Fehlerschwelle beendet den Container“) —
entweder in „Leerlauf-Bestätigung“ **nach** dem Slice (dann ist §2 falsch, die Karenz aus 2.5
wird neu entschieden) oder in einem Betriebsbericht bei nicht gesenkter Schwelle.

---

## 3. Verdikt (2) — Abhilfe und Slice-Vorschlag an den Planner

**Kennung:** `slice-leerlauf-phase-last-in-stuecken` (wellenlos, Umfang S — Schätzung: ein
Schleifen-Block im Runner, vier Träger-Texte).

**Start-Trigger:** `slice-capture-leerlauf-quellbelege` liegt in `done/` — erfüllt. Der Slice
ist der erste der Kette; `slice-wal-fehlerschwelle-ausgangsklasse` erhält ihn als Vorgänger.

**Reihenfolge:** `slice-capture-leerlauf-quellbelege` (done) →
**`slice-leerlauf-phase-last-in-stuecken`** → `slice-wal-fehlerschwelle-ausgangsklasse` →
`slice-start-vorlauf-grenze` → `slice-transformationen-e2e-abhilfe` →
`slice-transformationen-betriebsdoku`. **Grund für die Position vor
`slice-wal-fehlerschwelle-ausgangsklasse`:** beide ändern denselben Runner
(`tools/harness/run-integration-tests.sh`) an **verschiedenen** Stellen (dieser die Last der
Phase „Leerlauf-Bestätigung“, Zeilen 3436–3444, und ihre Texte; jener die Phase „Fehlerschwelle
beendet den Container“ ab Zeile 3478); ein Rot in der Leerlauf-Phase lässt die Phase dahinter
ungelaufen (M1: Schritte 5 bis 7 `skipped`), und der kleine Slice zuerst nimmt der Folge-Kette
das Masken-Rot.

**DoD-Kern:**

1. **Last in Stücken.** Der Schreiber auf `$WAL_FOREIGN` läuft in Stücken (Vorschlag:
   6 × 10.000 Zeilen, je etwa 2,68 MB WAL — M6, dort **erprobt**: Stelle „echter Stream,
   Wegwerf-Go-Test, lokal“, gesehen: höchster Rückstand 2.676.816 B in 3 von 3 Läufen). Nach
   jedem Stück wartet der Runner mit Frist (Vorschlag 30 s) darauf, dass `confirmed_flush_lsn`
   des Slots die WAL-Position hinter dem Stück erreicht (`pg_replication_slots`, dieselbe
   Abfrage-Form wie Zeile 3332 des Runners); ein Ablauf der Frist färbt die Phase rot mit
   benanntem Text. Zwei Wächter im Runner: das WAL je Stück liegt **unter** `WAL_WARN_BYTES`,
   und die **Summe** liegt über `WAL_ERROR_BYTES` (die bestehende Prüfung bleibt, Zeile 3443).
2. **Träger nachgezogen** (Suchlauf nach `AGENTS.md` §3.13 im Plan, beide Stände gemessen):
   der Kommentar der Phase (Zeilen 3371–3378), die Zeile `abdeckung_declare` (3369), die
   Ausgabezeile (3464), die `make test-integration`-Zeile in `harness/README.md` §Sensors
   („ein Schreiber auf eine nicht aktivierte Tabelle mit ebenfalls größerem WAL …“) und das
   Handbuch (`docs/user/benutzerhandbuch.md` bei „dasselbe gilt für jeden Schreiber auf eine
   nicht aktivierte Tabelle“ und bei der Beschreibung von `cdc_wal_retention_bytes`): die
   Zusage heißt „Last in Stücken unter der Warnschwelle, Summe über der Fehlerschwelle,
   bestätigt zwischen den Stücken“; ein Handbuch-Satz nennt die Grenze aus 2.6.
   `docs/user/e2e-abdeckung.md` ist Erzeugnis des Runners.
3. **Backfill-Hälfte der Phase.** Der Backfill-Run (22 MB in einer Transaktion, geschrieben
   vom Worker) trug in den 80 Ausführungen von M3 kein Rot; dass er kein Stoß ist, ist
   **hergeleitet** (die Kopierdauer des Runs druckt der Runner nicht). Der Implementer misst
   die Spitze des Rückstands im Run einmal mit Proben ≤ 50 ms an einer Wegwerf-Instanz und
   berichtet Stelle, Instanz und Zahl; liegt sie über der Hälfte der Fehlerschwelle, meldet er
   es an den Architect, statt die Run-Größe still zu ändern.
4. **Beleg.** Ein lokaler, grüner `make test-integration` und der erste grüne
   `e2e.yml`-Lauf mit beiden Legs (Lauf, Versuchsnummer und Job-Kennungen im Bericht, wie es
   der Abschluss von `slice-capture-leerlauf-quellbelege` bei der Behandlung von F-1 vorgibt).
   Die Mutation „die Leerlauf-Bestätigung des Streams bestätigt nichts“ (z. B. `ConfirmIdle`
   liefert kein `Acknowledged`) färbt die Phase an der Wartebedingung rot — *erwartet, zu
   erproben*: eine Stelle, ein Lauf, Farbe gedruckt; ohne Lauf steht sie als *hergeleitet*.
   Eine Aussage über die Stabilität in CI ist keine DoD (ein grüner Lauf beweist bei 1,25 %
   nichts); die Wirkung trägt M6 und der Trigger aus 2.6.
5. **Register:** der Eintrag `BEO-PGC/test-integration-retention-timing-flake` führt den
   Ausgang (Nachzug des Planners bei der Closure; der Zwischenstand steht im `state.md`).

**Abgrenzung:** keine Änderung an Schwellen, an `runWALRetentionCheck`, an
`ADR-0049`/`ADR-0120`, am Workflow `e2e.yml` (daher greift `AGENTS.md` §3.10 nicht) und an der
Phase „Fehlerschwelle beendet den Container“ — sie **soll** einen Stoß über der Schwelle nicht
überleben, dort ist das Fenster kein Fehler.

---

## 4. Verdikt (3) — Wirkung auf die Belege der Folge-Slices

**Was ein Rot hier auslöst** (gemessen, M1): der Runner endet mit `bf_fail`; jede Phase dahinter
läuft nicht — „Fehlerschwelle beendet den Container“, die Transformations-Rundläufe, der
Upgrade-Rundlauf — und im Job die Schritte dahinter (zwei Coverage-Schritte,
„Replication-Tier“). Ihr Beleg ist dann **nicht widerlegt, sondern ungelaufen**.

**Regel bis zum Slice:** ein Rot des ersten Versuchs mit der Signatur aus M1 in der Phase
„Leerlauf-Bestätigung“ ist kein Befund gegen den Slice, der gerade schließt, und kein Beleg
für ihn; der Verifier führt den Wiederholungsversuch (`gh run rerun <Lauf> --failed`) und
nennt in seinem Bericht Lauf, Versuchsnummer und Job-Kennungen beider Versuche. Ein Rot anderer
Signatur oder in einer anderen Phase gilt nicht als dieser Fall. Die Erwartung, dass es in einem
Slice-Zyklus zu einem solchen Rot kommt, ist klein (bei 1 von 40 je Lauf und zwei Läufen je
Zyklus etwa 5 %, *hergeleitet*), trifft aber `slice-wal-fehlerschwelle-ausgangsklasse` an der
Stelle mit dem höchsten Preis: seine Runner-Phase folgt der Leerlauf-Phase unmittelbar.

**Regel nach dem Slice:** dasselbe Rot ist ein Befund und ein Architect-Zug (Trigger aus 2.6);
die Wiederholung macht es nicht grün im Sinne des Belegs.

---

## 5. Verdikt (4) — Sofortmaßnahme

Zwei Befehle, beide am 2026-09-27 gefahren; kein Sensor, kein Gate, kein neuer Vorgang. Sie
stehen im `state.md` des Registers:

```text
gh api --paginate 'repos/pt9912/pg-change-feed/actions/workflows/e2e.yml/runs?per_page=100' --jq '.workflow_runs[] | select(.run_attempt > 1) | [.id, .head_sha[0:8], .run_attempt, .created_at] | @tsv'
gh run view <Lauf> --attempt 1 --log-failed | grep -c 'Leerlauf-Bestätigung — nach dem Schreiber'
```

Der erste listet jeden Lauf von `e2e.yml`, der wiederholt wurde (Ausgabe heute: genau
`36287009221`, `e4b77a05`, Versuch 2, von 444 Läufen); der zweite prüft die Signatur im ersten
Versuch (Ausgabe für `36287009221`: `1`). Ein wiederholter Lauf mit Ausgabe ≥ 1 ist ein
weiteres Rot der Phase — der Trigger aus 2.6, wenn es **nach** dem Slice liegt.

---

## 6. Register

`BEO-PGC/test-integration-retention-timing-flake`: die **Resthälfte** („ein
`make test-integration`-Lauf endet rot bei unverändertem Stand, die Wiederholung ist grün“)
war mit 2× geführt (`slice-057`, `slice-capture-leerlauf-quellbelege`). Die zwei Fälle haben
**verschiedene Ursachen**: `slice-057` (Retention-Lebenszyklus-Timing) bleibt 1× ohne Ursache
und ohne Ausgang; der Leerlauf-Fall hat eine Ursache (Stoß gegen Takt) und den Ausgang „Slice
`slice-leerlauf-phase-last-in-stuecken`“. Der Eintrag ist entsprechend nachgezogen (`state.md`
und `evidence/`); die Zahl-Hälfte (Coverage) bleibt verkörpert. Kein neuer Eintrag.

## 7. Berichtigte Aussagen im Register

Zwei Angaben im Register hielten der Nachmessung nicht stand:

1. **„Der `INSERT` über 60.000 Zeilen lief etwa 19,8 s.“** Die zwei Stempel 02:02:26,00 Z
   (die Ausgabe `INSERT 0 30000`/`CREATE TABLE` der Einrichtung) und 02:02:45,76 Z (`INSERT 0
   60000`) begrenzen den **Abschnitt** der Phase — Aktivierung, Backfill-Run über 30.000
   Zeilen, 12 s Halten, dann der Schreiber; `psql` druckt seine Zeile erst bei der Rückkehr.
   Die Dauer der Anweisung auf dem Runner ist nicht gemessen; lokal sind es 0,14 bis 0,25 s
   (M4/M5). Die Aussage „Last rund das Doppelte der Schwelle **in etwa 20 s**“ (`state.md`)
   trifft nicht zu: es ist ein Stoß, kein Zulauf.
2. **„1 von 40 Läufen“** zählt Läufe; die Phase läuft je Lauf in zwei Legs. Die Zahl der
   Ausführungen ist 80 (M3), das rote Leg ist 1 davon.

Records (`done/`, `docs/reviews/**`) bleiben stehen. Der Review-Report von
`slice-capture-leerlauf-quellbelege` nennt die Last „rund das Doppelte der Schwelle“ (stimmt,
M2) und nicht die 20 s.

## 8. Was nicht getan wurde — und die Nachzüge

Kein Produktionscode, keine neue oder geänderte ADR, keine Änderung an `spec/**`,
`docs/user/**`, `harness/README.md` und an den Plänen. Geschrieben: dieses Verdikt und der
Register-Nachzug (`state.md`, `evidence/architect-verdict-leerlauf-bestaetigung-intermittenz.md`,
Berichtigung von `evidence/slice-capture-leerlauf-quellbelege.md`). Die Wegwerf-Messungen
(Scratchpad, ein nicht committeter Go-Test) sind entfernt; die Wegwerf-Container und das
Wegwerf-Netz sind abgeräumt (`docker rm -f -v`, `docker network rm`; dangling Volumes vorher und
nachher 36).

Offen bei anderen Rollen:

1. **Planner:** den Plan `slice-leerlauf-phase-last-in-stuecken` anlegen (Adresse: §3;
   Frist: der Start von `slice-wal-fehlerschwelle-ausgangsklasse`); die Kante
   „Leerlauf-Stabilisierung → `slice-wal-fehlerschwelle-ausgangsklasse`“ als Start-Trigger dort
   und in `slice-start-vorlauf-grenze` §4 nachziehen; `welle-transformationen` §5 („Bedingte
   Kante: Stabilisierung …“) auf **beauftragt** setzen und die Kette dort und in der Roadmap
   um den neuen Slice ergänzen; die Regel aus §4 in den Plänen der drei Folge-Slices als
   Verifier-Hinweis führen.
2. **Implementer** des Slice: Punkte 1 bis 5 aus §3.
