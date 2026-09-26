# ADR-0128: Prozessstart — der Vorlauf der Antrags-Queue trägt eine Frist, und der Replikationsstrom beginnt erst im Stream-Lauf (schärft ADR-0112)

**Status:** Accepted — Supersedes: keine. Schärft
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5
(Bedingung (c), die Startreihenfolge) um Frist und Ausgang des Vorlaufs. Der Befund am Bestand
in [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) („danach startet `NewStream` sofort
`START_REPLICATION`“) beschreibt den Stand vor dieser ADR; er bleibt als Befund stehen und wird
mit der Umsetzung ungenau (§Konsequenzen, Folgepflicht 5).

**Datum:** 2026-09-27

**Autor:** Architect-Agent (Modul 8), Architect-Zug zu Frage (e) von `welle-transformationen`
(Vollmacht des Auftraggebers); jede Tatsachenaussage trägt ihren Beleg-Anker (gedruckte
Messzeile, §Gemessen) oder ist als hergeleitet gekennzeichnet.

**Bezug:** [`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Transformationen, Abhilfe),
[`LH-FA-ADM-003`](../../../spec/lastenheft.md) (erkennbarer Fehlerzustand),
[`LH-QA-PER-004`](../../../spec/lastenheft.md) (Commit bis CDC-Verfügbarkeit),
[`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0049`](0049-replication-fehlerklassen-schwellen.md),
[`ADR-0007`](0007-source-ack-outbound-port.md)

**Schärft:** [`LH-FA-CFG-007.a`](../../../spec/pflichtenheft.md) (Zusage der Abhilfe: offene
Anträge werden beim Start verarbeitet, bevor die erste Transaktion der Tabelle assembliert
wird) und [`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md) Schritt 1 (Receive: Aufbau der
Verbindung und Beginn des Empfangs).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`Run` (`internal/bootstrap/wiring.go`) verarbeitet die offenen Anträge der Antrags-Queue in einem
synchronen Vorlauf, bevor der Stream startet (`runStreamAfterAdministrationPass`); die Ordnung
trägt Bedingung (c) der Abhilfe von `ADR-0112` Folgepflicht 5: die Abhilfe-Anträge stehen
`applied` in der `Assembler`-Bindung, bevor die erste Transaktion der Tabelle assembliert wird.
Der Vorlauf trägt keine eigene Frist; ihn beendet allein der Kontext des Streams (Prozess-`ctx`
und WAL-Fehlerschwelle). Vor dem Vorlauf lief die Administrations-Goroutine nebenläufig zum
Stream (Stand des Schnitts der Welle, gelesen an `wiring.go`), ein langer Antrag hielt damals nur
die Queue an. Frage (e) fragt nach einer Frist, dem Verhalten bei Ablauf und der Sichtbarkeit des
Wartens.

Der Aufbau des Stream-Adapters ist die zweite Hälfte der Frage: `receive.NewStream` baut die
Verbindung auf, legt Slot und Publication an **und sendet `START_REPLICATION`**
(`internal/adapters/driving/replication/receive/receive.go`, Aufruf `session.StartReplication`);
`Stream.Run` beginnt mit dem Lesen. Zwischen beiden liegt in `Run` der gesamte Rest der
Verdrahtung und der Vorlauf (`NewStream` an `wiring.go:641`, die Sequenz an `wiring.go:1041`,
gelesen am Stand `ca9aaea7`).

### Gemessen

Alle Läufe am Feed-Container `ghcr.io/pt9912/pg-change-feed:dev` (Bau 2026-09-26 23:41 CEST,
Stand des Vorlaufs in `wiring.go` enthalten), Compose-Umgebung aus `compose.yaml` mit
PostgreSQL 18 (`postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8`),
Stand des Arbeitsbaums `ca9aaea7`, Host Linux 6.8.0-139-generic. **Ablauf je Lauf** (ein
Wegwerf-Skript im Scratchpad des Zugs, nicht committet): Schema-Rollout, Tabellen und Quelle
anlegen, Feed starten und auf `healthy` warten, Feed stoppen, eine zweite Sitzung sperrt
`public.feed_wait_probe` mit `LOCK TABLE … IN ACCESS EXCLUSIVE MODE` und hält die Sperre mit
`pg_sleep`, `SELECT cdc.enable_table('src-e2e','public','feed_wait_probe')` legt einen
`pending`-Antrag an, Feed starten, dann im Abstand von rund 3 s: Heartbeat-Alter
(`SELECT age_seconds FROM cdc.heartbeat`), `docker exec cdc-test-feed /pg-change-feed
--healthcheck` (Exit), Status des Antrags, und ab etwa 12 s eine in die aktivierte Tabelle
`public.feed_e2e_full` eingefügte Zeile, deren Sichtbarkeit über `cdc.changes` gezählt wird.

| Lauf | `wal_sender_timeout` | Sperre | Abfragen mit Healthcheck-Exit 0 | größtes Heartbeat-Alter | eingefügte Zeile in `cdc.changes` | Prozess |
|---|---|---|---|---|---|---|
| A | 2000 ms (Compose-Einstellung) | 50 s | 22 Abfragen von t = 1 s bis 48 s (danach der Ausgang 1) | 7,1 s | 0 bis zum Ende des Vorlaufs | endet bei Freigabe der Sperre (t ≈ 49 s), Antrag danach `applied` |
| B | 2000 ms (Compose-Einstellung; ein `ALTER SYSTEM` blieb wirkungslos, die Kommandozeile hat Vorrang) | 75 s | 23 Abfragen von t = 0 s bis 71 s | 6,8 s | 0 bis zum Ende des Vorlaufs | endet bei Freigabe (t ≈ 74 s) |
| C | 60 s (PostgreSQL-Standard; Compose-Override im Scratchpad ohne die Einstellung) | 75 s | 23 Abfragen von t = 0 s bis 72 s | 6,8 s | 0 bis zum Ende des Vorlaufs | endet bei Freigabe (t ≈ 75 s) |

Der Healthcheck-Exit ist `0` und das Heartbeat-Alter bleibt unter der Schwelle von 15 s
(`heartbeatStaleAfter`) über die **ganze** Dauer des Vorlaufs: der Vorlauf ist von einem
gesunden Lauf nicht zu unterscheiden, und die Erfassung der ganzen Quelle steht (die Zeile ist
während des Vorlaufs nie sichtbar; nach dem Prozessende in den Läufen ebenfalls nicht, weil
kein Stream mehr läuft).

**Der Prozess endet mit der Fehlerklasse `replication`, nicht mit einem Warten.** Gedruckte Zeilen:
Lauf C — Feed-Log `2026-09-26T22:42:33.862Z … "replication: Stream gestartet"`, PostgreSQL-Log
`2026-09-26 22:43:33.863 UTC … LOG:  terminating walsender process due to replication timeout`
(60,0 s später), Feed-Log `2026-09-26T22:43:45.6297Z … "administrationrequest: Antrag erledigt"`
und im selben Zeitpunkt `… "replicationack: Bestätigungsfehler","error":"write failed: … write:
broken pipe"`, `… "replication: Stream beendet mit Fehler"`, danach
`"heartbeat: Fehlerzustand gemeldet","class":"replication"`. Lauf B — Feed-Log `Stream gestartet`
`22:40:16.2807Z`, PostgreSQL-Log `terminating walsender process due to replication timeout`
`22:40:18.281` (2,0 s später). Ursache: `NewStream` hat den Strom vor dem Vorlauf gestartet, der
Server beendet einen Strom, der `wal_sender_timeout` lang nicht liest und nicht antwortet, der
Prozess bemerkt es erst am ersten Schreiben nach dem Vorlauf.

**Eine Replikationsverbindung im Kommando-Zustand läuft über `wal_sender_timeout` hinaus.**
`psql "dbname=postgres user=postgres replication=database"` gegen PostgreSQL 18 (dasselbe
Image) mit `-c wal_sender_timeout=2000`: `IDENTIFY_SYSTEM`, 7 s Pause, `IDENTIFY_SYSTEM`,
`CREATE_REPLICATION_SLOT idle_probe LOGICAL pgoutput`, 7 s Pause, `IDENTIFY_SYSTEM` — alle drei
Kommandos antworten, das Log trägt keine Zeile „terminating“. PostgreSQL 17 nicht gefahren
(hergeleitet: dieselbe Server-Einstellung, die Verbindung wartet im selben Protokollzustand).

Kein Lauf fuhr eine andere Antragsart als `enable`; dass der Vorlauf für **jede** Antragsart
dieselbe Wartestelle ist, ist hergeleitet aus `processAdministrationRequests` (eine Schleife über
`ListPending` mit `applyAdministrationRequest`), nicht je Art erprobt.

## Entscheidung

Vier Festlegungen.

1. **Der Replikationsstrom beginnt im Stream-Lauf.** `NewStream` baut Verbindung, Slot,
   Publication und `Assembler` auf und lässt die Verbindung im Kommando-Zustand; `Stream.Run`
   sendet `START_REPLICATION` als erste Handlung. Solange kein Strom läuft, gilt
   `wal_sender_timeout` nicht für diese Verbindung (Gemessen, dritter Absatz). Damit ist die
   Länge des Vorlaufs für den Prozess folgenlos, unabhängig davon, welchen Wert die Quelle für
   `wal_sender_timeout` trägt (Compose-Umgebung: 2000 ms, PostgreSQL-Standard: 60 s).
2. **Der Vorlauf trägt eine Frist von 30 s je Prozessstart** (ein Durchlauf von
   `processAdministrationRequests` vor dem Stream-Start, nicht je Antrag). Der Wert ist
   **hergeleitet**: die Hälfte der Fehlergrenze von 60 s des Capture-Abstands
   (`SPEC-013`, `cdc_capture_lag`), damit die Verzögerung, die der Vorlauf einer zu diesem
   Zeitpunkt committeten Änderung hinzufügt, unter der Fehlergrenze bleibt und die
   Verfügbarkeit des Streams selbst noch Raum hat. Er ist eine Konstante des Codes wie
   `heartbeatInterval`, keine Konfigurationsachse.
3. **Bei Ablauf startet der Stream, der Antrag bleibt `pending`.** Der Vorlauf endet mit einem
   Warn-Eintrag im Log (Frist und der Antrag, an dem sie ablief); der Antrag, den die Frist
   unterbricht, und jeder Antrag dahinter werden **nicht** als `failed` vermerkt, sie bleiben
   `pending`; die Administrations-Goroutine startet und verarbeitet sie ohne Frist; der Stream
   startet mit dem Regelstand, den der Vorlauf bis dahin nachgetragen hat. Die Ordnung „Antrag
   `applied` vor der ersten Transaktion der Tabelle“ gilt damit für jeden Antrag, den der Vorlauf
   **innerhalb der Frist** erreicht; für einen späteren gilt sie nicht.
4. **Das Warten wird nicht in `diagnose` und `--healthcheck` angezeigt.** Begründung: ein
   Warten hat mit Festlegung 1 und 2 eine Obergrenze von 30 s und heilt sich selbst — es ist
   kein Fehlerzustand im Sinn von `LH-FA-ADM-003`; der Log-Eintrag der Festlegung 3 ist sein
   Träger. Während dieser Zeit ist der Prozess für den Healthcheck gesund (Gemessen: die
   Läufe A bis C; das Verhalten ist an den drei Läufen **gemessen**, für die Frist selbst
   **erwartet**, bis der Rundlauf der Fitness Function läuft).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Vorlauf ohne Frist, Strom beginnt in `NewStream`) | keine Änderung | gemessen (Läufe A bis C): die Erfassung steht unbegrenzt, der Prozess erscheint gesund, und der Prozess endet mit der Klasse `replication`, sobald der Vorlauf `wal_sender_timeout` überschreitet (Standard 60 s); der Ausgang nennt die Ursache nicht |
| B — nur eine Frist, Strom beginnt weiter in `NewStream` | kleiner Diff | die Frist müsste unter `wal_sender_timeout` liegen, der Wert gehört der Quelle (Compose: 2 s, Standard: 60 s, `0` schaltet ihn ab) und ist bei Bau des Prozesses unbekannt; jeder kleinere Betreiberwert bringt den Ausgang der Läufe B und C zurück |
| C — Frist aus `SHOW wal_sender_timeout` ableiten (Hälfte des Werts) | löst B ohne Änderung an `NewStream` | koppelt den Prozessstart an eine Server-Einstellung und ihre Sonderfälle (`0`, Einheiten), jede weitere Wartestelle zwischen `NewStream` und `Run` trägt dieselbe Falle; der Kommando-Zustand macht die Kopplung überflüssig |
| D — Vorlauf wieder nebenläufig (Goroutine vor `stream.Run`, keine Ordnung) | kein Warten, kein Ablauf | verwirft Bedingung (c) der Abhilfe; das ist der Re-Evaluierungs-Trigger 4 von `ADR-0112` |
| E — Frist mit Vermerk `failed` bei Ablauf | einfache Zuordnung, ein Antrag hat nach dem Start einen Endzustand | ein gesunder, nur langsamer Antrag (Sperre des Betreibers) würde als fehlgeschlagen vermerkt und müsste neu beantragt werden; der Abhilfe-Antrag selbst kann so verloren gehen (Neustart ohne Regelstand, erneutes `schema`) |
| **F — Strom beginnt in `Stream.Run`, Frist 30 s, Rest bleibt `pending` (gewählt)** | Folge des Wartens ist begrenzt und für jeden Serverwert von `wal_sender_timeout` harmlos; kein Antrag geht verloren; die Ordnung gilt im normalen Fall (Antrag in Millisekunden) unverändert | die Erfassung startet bei einem hängenden Antrag bis zu 30 s später; die Ordnungszusage gilt nur innerhalb der Frist |

## Konsequenzen

- Positiv: der Prozess stirbt nicht mehr an einem langen Vorlauf; der Ausgang „Klasse
  `replication`, Ursache im Log nicht zu finden“ entfällt (Läufe B und C).
- Positiv: das Warten ist auf 30 s begrenzt; die Regression gegenüber dem nebenläufigen Stand
  (ein langer Antrag hält die ganze Erfassung an) ist bis auf diese Grenze zurückgenommen.
- Negativ: bei einem Antrag, der beim Start länger als 30 s wartet, startet die Erfassung mit dem
  bisherigen Regelstand. Ist die Regel dabei nicht anwendbar und der wartende Antrag die
  Abhilfe, endet der Prozess erneut mit der Klasse `schema`; der Betreiber startet neu, sobald der
  Antrag `applied` ist (Antrag über `cdc.administration_request` lesbar). Gilt für den
  Abhilfe-Antrag nur, wenn ihm ein Antrag vorausgeht, der länger als die Frist läuft
  (hergeleitet aus der Ordnung `requested_at`, `administration_request_id`).
- Negativ: der Slot hält WAL, solange die Verbindung im Kommando-Zustand wartet; die Messung
  des WAL-Rückstands und ihre Fehlerschwelle laufen vor dem Vorlauf (`walRetention`-Goroutine
  vor der Sequenz, gelesen am Stand `ca9aaea7`), sie beenden den Vorlauf über `streamCtx`
  (hergeleitet, nicht gefahren).
- Erwartet, nicht am Code belegt: `Stream.Run` kann `START_REPLICATION` als erste Handlung
  senden, ohne dass der ACK-Adapter, der dieselbe Verbindung teilt (`ADR-0007`), vor dem ersten
  Lesen schreibt (hergeleitet aus `receive.go`; der Implementer prüft es).
- Folgepflicht 1: **Umsetzung** — Slice `slice-start-vorlauf-grenze` (`open/`, wellenlos; Kante
  vor `slice-transformationen-e2e-abhilfe`): `START_REPLICATION` nach `Stream.Run`, Frist im
  Vorlauf, kein `failed` bei Ablauf, Warn-Eintrag, die drei Tests der Fitness Function.
- Folgepflicht 2: **Pflichtenheft** — der Absatz „Abhilfe (Zusage)“ von `LH-FA-CFG-007.a` nennt
  die Ordnung „innerhalb der Frist des Vorlaufs“ (Träger: derselbe Slice, sein Plan trägt die
  Zeile).
- Folgepflicht 3: **Handbuch** — `slice-transformationen-betriebsdoku` §2 führt die
  Betreiber-Aussage „ein hängender Antrag hält den Stream-Start höchstens 30 s an“ mit der
  Messung des Rundlaufs der Fitness Function, nicht als Herleitung.
- Folgepflicht 4: **Beobachtungs-Register** — `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad`
  bekommt den Ausgang *geplant* mit dem Slice als Anker.
- Folgepflicht 5: **Träger, die den Stand vor der Umsetzung beschreiben** — der Suchlauf des
  Umsetzungs-Slice nach `START_REPLICATION` und `NewStream` (§3.13): der Befund in `ADR-0111`
  bleibt stehen (`Accepted`); die Kommentare in `receive.go` (Godoc von `NewStream`) und
  `wiring.go` (Zeile „Der Slot besteht an dieser Stelle bereits“) beschreiben den Aufbau und
  werden vom Slice nachgezogen.

## Fitness Function (falls maschinell prüfbar)

Alle drei Zeilen sind **Erwartungen**: sie sind nicht erprobt; der Implementer des Slice fährt sie.
Die Mutation jeder Zeile ist **hergeleitet**, an keiner Stelle gefahren.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Whitebox, `internal/bootstrap`, Fake-Queue) | die Queue-Lesung blockiert bis zum Ende des übergebenen Kontexts; nach der (im Test verkürzten) Frist startet die Goroutine und danach der Stream; der Antrag bleibt `pending` und wird nicht als `failed` vermerkt; ein Warn-Eintrag steht im Log. Mutation (hergeleitet): der Vorlauf ohne Frist (`ctx` unverändert) färbt den Test rot; ein `MarkFailed` bei Ablauf färbt ihn rot | `make test` |
| Go-Test (Replication-Tier, reale PostgreSQL mit `wal_sender_timeout` 2 s) | `NewStream`, danach länger als `wal_sender_timeout` warten, danach `Run`: eine danach committete Änderung wird geliefert und bestätigt. Mutation (hergeleitet): `StartReplication` zurück in `NewStream` färbt den Test rot; erprobt ist an der Quelle nur der Kommando-Zustand mit `psql` (Gemessen, dritter Absatz), nicht der Adapter | `make test-replication` (und mit `PG_TEST_IMAGE` auf dem PostgreSQL-17-Digest) |
| Realer Rundlauf | ein `pending`-Antrag wartet an einer Sperre des Betreibers länger als die Frist: der Prozess läuft, der Healthcheck meldet `0`, eine währenddessen committete Änderung ist über `cdc.changes` binnen Frist plus Toleranz sichtbar, der Antrag ist nach Freigabe der Sperre `applied`, der Prozess endet nicht mit einer Fehlerklasse | `make test-integration` |

## Re-Evaluierungs-Trigger

- **Ein Betreiber meldet, dass 30 s zu lang oder zu kurz sind**: die Frist wird ein
  Konfigurationswert (Folge-ADR).
- **Ein Betreiber meldet einen Ablauf der Frist, den er nicht aus dem Log erkannt hat**: eine
  Anzeige in `diagnose` wird geprüft (Festlegung 4).
- **Eine weitere synchrone Wartestelle zwischen `NewStream` und dem Stream-Lauf entsteht**:
  Festlegung 1 gilt für sie; sie braucht keine eigene Frist gegen `wal_sender_timeout`.
- **`Stream.Run` sendet `START_REPLICATION` nicht mehr als erste Handlung**: diese ADR ist zu
  prüfen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-27 | Accepted — Architect-Zug (Vollmacht des Auftraggebers) zu Frage (e) von `welle-transformationen`; das Warten an einer Sperre an drei Läufen gemessen (§Gemessen) | [`LH-FA-CFG-007`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0128` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
