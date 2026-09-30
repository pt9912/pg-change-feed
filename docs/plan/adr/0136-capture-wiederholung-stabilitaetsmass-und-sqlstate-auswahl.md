# ADR-0136: Capture-Wiederholung — Stabilitätsmaß der Episode und SQLSTATE-Auswahl (Supersedes ADR-0135, teilweise)

**Status:** Accepted — Supersedes [`ADR-0135`](0135-capture-transient-wiederholung-stream-zyklus.md)
in genau **zwei Klauseln**: dem Satz „Ein erfolgreicher Stream-Zyklus setzt die
Episode zurück" in Festlegung 2 (die Messgröße „erfolgreich" war offen und ist
im Code als Zyklus-Dauer entstanden) und der Formulierung „genau die
Transport-/Verbindungsstörung … Kette `receive.ErrReplication` oder
`outbound.ErrReplication`" in Festlegung 3 (die Fehlermenge hat eine
SQLSTATE-Auswahl, die dort nicht steht). Alles Übrige aus `ADR-0135` bleibt
in Kraft: Ort und Wiederholeinheit (Festlegung 1), die Backoff-Werte
2 s / Faktor 2 / 30 s / 5 min, die Ausschlüsse `configuration`/`schema`/
Ordnungs-Verletzung/`storage`, Sichtbarkeit, Ausgang bei Erschöpfung,
`restart: "no"` (Festlegungen 4–6).

**Datum:** 2026-09-30

**Autor:** Architect-Agent (Modul 8), Folge-Entscheidung zu den Findings N-1 und
N-2 der Fixrunde 1 im Review-Report zu `slice-capture-transient-wiederholung`

**Bezug:** [`LH-QA-REL-002`](../../../spec/lastenheft.md),
[`LH-QA-REL-001`](../../../spec/lastenheft.md),
[`LH-FA-ADM-003`](../../../spec/lastenheft.md),
[ADR-0023](0023-fehlerklassifikation.md),
[ADR-0049](0049-replication-fehlerklassen-schwellen.md),
[ADR-0128](0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)

**Schärft:** [`SPEC-008`](../../../spec/pflichtenheft.md) (Zeile `transient` —
Bedingung der wiederholten Fehlermenge)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0135` ist umgesetzt (`runStreamWithRetry`, `receive.serverFault`). Der
Re-Review der Fixrunde 1 fand zwei Setzungen, die in keiner Entscheidung stehen:

- **N-1.** Die Episode wird zurückgesetzt, wenn ein Zyklus mindestens
  `streamRetryStableAfter` (30 s) gedauert hat. Gemessen wird ab **vor**
  `NewStream` — die Dauer des Verbindungsaufbaus zählt mit. Ein Zyklus, der
  31 s im Verbindungsaufbau hängt und mit `receive.ErrReplication` endet,
  gilt als stabil; die Schleife kehrt dann nie mit `ErrTransientExhausted`
  zurück (Probe des Reviewers: 201 Zyklen, 1 h 50 min Test-Uhr, keine
  Erschöpfung — **übernommen** aus dem Review-Report, hier nicht nachgefahren).
  Gerade „Quelle unerreichbar" behielte damit keine Grenze, und der
  Heartbeat-Fehlerzustand der Klasse `transient` bliebe aus. Ob ein 31 s
  hängender Verbindungsaufbau real vorkommt, ist **hergeleitet**: `pgconn`
  trägt ohne `connect_timeout` in der DSN keine eigene Frist; die Frist des
  Betriebssystems für einen verworfenen TCP-Verbindungsversuch liegt bei
  Minuten. Nicht am Netz gemessen.
- **N-2.** `receive.serverFault` wiederholt nur SQLSTATE-Klassen 08, 40, 53,
  55, 57, 58 sowie Fehler ohne SQLSTATE (Verbindungsabbruch); 42501 und
  Klasse 28 werden `permission`; jeder andere SQLSTATE trägt zusätzlich
  `ErrRejected` und endet sofort. `ADR-0135` Festlegung 3 sagt dagegen, wiederholt
  werde „genau" die `ErrReplication`-Kette — `ErrRejected` trägt diese Kette,
  wird aber nicht wiederholt.

Der Auftrag ist, beides zu entscheiden, ohne die Wiederholung zu verkomplizieren:
ein Maß und eine Auswahlregel, kein neuer Konfigurationsschlüssel.

## Entscheidung

**1. „Erfolgreicher Zyklus" = der Server hat `START_REPLICATION` bestätigt und
der Stream lief danach mindestens 30 s.** Die Stabilitätsdauer wird ab dem
Zeitpunkt gemessen, an dem `Stream.Run` die Bestätigung des Servers für
`START_REPLICATION` erhalten hat (Rückkehr von `StartReplication`, dieselbe
Stelle, an der der Stream-Adapter „Stream gestartet" loggt) bis zum Fehler des
Zyklus. Verbindungsaufbau, Katalog-/Slot-Abfragen und ein abgewiesenes
`START_REPLICATION` (zum Beispiel 55006) zählen nie zur Stabilität: ein Zyklus,
der den Streaming-Zustand nicht erreicht hat, setzt die Episode nicht zurück,
gleich wie lange er dauerte. Es wird bewusst **nicht** verlangt, dass eine
Transaktion geliefert oder eine Position bestätigt wurde: ein leerer Quell-Stream
liefert nichts, und ein ruhiger, gesunder Stream, der nach einer Stunde stirbt,
soll seine Episode zurücksetzen. Die Schwelle bleibt **30 s** (`streamRetryMaxDelay`); sie
ist **abgeleitet**, keine Messung: ein Stream, der länger streamt als der
größte Warteschritt, hat länger gearbeitet, als eine weitere Wartezeit ihn
gekostet hätte.

**2. Der Zyklus-Aufbau bekommt eine Frist von 30 s.** `NewStream` (Verbindungsaufbau,
Katalog- und Slot-Abfragen) läuft unter einer Kontext-Frist von
`streamRetryMaxDelay`; ihr Ablauf ist ein Fehler ohne SQLSTATE und damit
wiederholbar (`ErrReplication`). Damit ist die Zeit bis zum nächsten Fehlerpunkt
im Aufbau begrenzt, und das Gesamtfenster (5 min) bleibt auch bei „Quelle
unerreichbar" wirksam. Die Frist gilt nur dem Aufbau; sie darf die
Lebensdauer des aufgebauten Streams nicht binden (`Run` erhält den
Prozess-Kontext). Die Frist ist an der Stelle gesetzt, wo der Zyklus aufgebaut
wird (Composition Root), nicht in der DSN — ein Betreiber-`connect_timeout`
bleibt wirksam und kürzer möglich. Dass `pgconn` bei Fristablauf des
übergebenen Kontexts den Verbindungsversuch abbricht und mit einem Fehler ohne
SQLSTATE zurückkehrt, ist **hergeleitet**, nicht erprobt (siehe
Fitness Function).

**3. Die wiederholte Fehlermenge ist eine Auswahl nach SQLSTATE — Ursache im
Zustand von Server oder Verbindung, nicht in der Anfrage.** Wiederholt werden
(Menge der geprüften Klassen: die in `serverFault` gelisteten und die unten
genannten Grenzfälle; **hergeleitet** aus der SQLSTATE-Tabelle der PostgreSQL-
Dokumentation, am Server nicht einzeln erprobt):

- Fehler **ohne SQLSTATE** (Verbindungsabbruch, Fristablauf, EOF);
- Klasse **08** (Verbindungsausnahme), **40** (Rollback wegen Deadlock oder
  Serialisierung), **53** (unzureichende Ressourcen, zum Beispiel 53300
  `too_many_connections`), **55** (Objekt nicht im erforderlichen Zustand,
  einschließlich 55006 „Slot noch aktiv"), **57** (Betreibereingriff,
  einschließlich 57P01 `admin_shutdown`, 57P02 `crash_shutdown`, 57P03
  `cannot_connect_now`), **58** (Systemfehler);
- zusätzlich der einzelne Code **25006** (`read_only_sql_transaction`): der
  Knoten ist während eines Failovers noch oder wieder schreibgeschützt; der
  Zustand geht vorbei.

**Nicht** wiederholt werden `permission` (42501, Klasse 28) und jede
Server-Abweisung außerhalb der obigen Menge (`ErrRejected`: zum Beispiel 42704
Slot/Objekt existiert nicht, 3D000 Datenbank existiert nicht, 0A000, Klasse XX
`internal_error`/`data_corrupted`, Klasse 25 außer 25006). Eine unbekannte
Klasse endet **sofort sichtbar** statt nach fünf Minuten — die Kosten eines
fälschlich nicht wiederholten transienten Fehlers sind ein Neustart des
Aufrufers (`LH-QA-REL-002`), die eines fälschlich wiederholten Dauerfehlers
fünf Minuten Verzögerung der Meldung. Die Regel ist damit **Positivliste**:
wer einen weiteren SQLSTATE aufnehmen will, tut es in einer Folge-ADR mit dem
Fall, der ihn trägt. Die Klassen 57 und 58 werden als Ganzes geführt, obwohl
einzelne Codes darin (57P04 `database_dropped`, 58P01 `undefined_file`)
nicht heilen: sie enden nach dem Fenster mit `transient` — der Preis der
Klassen-Granularität, bewusst gewählt gegenüber einer Code-Liste, die niemand
pflegt.

**4. Sichtbarkeit: das INFO „fortgesetzt" folgt dem Signal aus Festlegung 1.**
Der Eintrag „Stream-Zyklus nach Wiederholung fortgesetzt" wird geschrieben, wenn
der Zyklus den Streaming-Zustand erreicht (dasselbe Signal wie die
Stabilitäts-Messung), nicht erst am Zyklus-Ende. Ein weiterlaufender Zyklus meldet
damit seine Fortsetzung. Die Rücksetzung der Episode bleibt an die 30 s aus
Festlegung 1 gebunden.

**5. `SPEC-008` wird nachgezogen.** Die Zeile `transient` nennt als Bedingung
weiter die Transport-/Verbindungsstörung (`receive.ErrReplication`/
`outbound.ErrReplication`), zusätzlich: nicht wiederholt enden Server-
Abweisungen ohne transiente Ursache (`receive.ErrRejected`, Auswahl nach
SQLSTATE, Festlegung 3 dieser ADR). Die Zeile ist eine Bedingung, die
`ADR-0135` als Folgepflicht schon der Festlegung 3 zugeordnet hat; sie
unterscheidet bisher nicht zwischen `ErrReplication` mit und ohne
`ErrRejected`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Status quo: Zyklus-Dauer ab vor `NewStream`, 30 s | kein Code | Hängender Verbindungsaufbau setzt die Episode zurück, das Gesamtfenster ist umgehbar (N-1) |
| B — Dauer ab Streaming-Beginn (`START_REPLICATION` bestätigt), 30 s (**gewählt**) | schließt N-1; trennt „Aufbau" von „Streamen" genau an der Stelle, die der Adapter schon loggt; kein Verkehr nötig, ruhige Streams zählen | braucht ein Signal aus `receive.Stream` in die Composition Root |
| C — Rücksetzung nur bei gelieferter Transaktion oder bestätigter Position | strengster Beleg für „hat gearbeitet" | ein ruhiger Stream ohne Änderungen setzt nie zurück; nach zwei Störungen im Abstand von Stunden wäre das 5-min-Fenster verbraucht — falsche Erschöpfung |
| D — höhere Schwelle (zum Beispiel 5 min) statt neuem Maß | ein Wert | verlagert N-1 (ein Aufbau kann auch dann hängen), verzögert die Rücksetzung gesunder Streams |
| E — nur eine Frist für den Aufbau, ohne neues Stabilitäts-Maß | begrenzt den Aufbau | ein Aufbau, der knapp unter der Frist scheitert, dauert weiter fast 30 s und setzt bei Maß A die Episode zurück; daher **ergänzend** zu B, nicht statt B |
| F — SQLSTATE: Negativliste (alles wiederholen außer bekannter Dauerfehler) | robust gegen unbekannte transiente Codes | unbekannte Dauerfehler verzögern die Meldung um das Fenster; die Liste der Dauerfehler ist offen |
| **G — SQLSTATE: Positivliste 08/40/53/55/57/58 + 25006 (gewählt)** | fail-fast bei Unbekanntem; entspricht dem Code | ein nicht gelisteter transienter Code kostet einen Aufrufer-Neustart |

## Konsequenzen

- Positiv: das Gesamtfenster von `ADR-0135` ist auch bei „Quelle unerreichbar"
  (Aufbau) wirksam; die
  Fehlermenge steht in einer Entscheidung, `SPEC-008` folgt.
- Negativ: ein flappender Stream, der jedes Mal mindestens 30 s streamt und dann
  fällt, erschöpft das Fenster nie — die Wiederholung kaschiert dann eine
  Dauerstörung. Sichtbar bleibt er über das WARN je Wiederholung; eine
  Obergrenze für Flapping-Episoden ist **nicht** Teil dieser Entscheidung
  (akzeptiertes Negativ: ein Stream, der 30 s arbeitet, liefert Änderungen und
  ist nicht ausgefallen).
- Neu: die Frist von 30 s für den Aufbau ist eine Setzung ohne Messung, wie die
  Werte in `ADR-0135`.
- Folgepflichten (Träger: `slice-capture-transient-wiederholung`, Fixrunde 2):
  1. **Code.** `receive.Stream` meldet der Composition Root den Beginn des
     Streamings (Form dem Implementer überlassen: Callback oder Kanal, aufgerufen
     nach Rückkehr von `StartReplication`); der Zyklus gibt `streamStart` (Null,
     wenn nie erreicht) an `runStreamWithRetry` zurück; die Rücksetzung
     prüft `now - streamStart >= streamRetryStableAfter`, nie `cycleStart`.
     `NewStream` läuft unter `context.WithTimeout(ctx, streamRetryMaxDelay)`; `Run`
     unter dem Prozess-Kontext. `serverFault` nimmt 25006 in die wiederholte
     Menge auf; Kommentare an `ErrRejected`, `serverFault` und
     `retryableStreamError` nennen Positivliste und `ADR-0136`. INFO
     „fortgesetzt" aus dem Signal.
  2. **Tests.** Je Punkt eine Mutation, die rot färben muss (Instanz und Stelle
     hier **erwartet**, nicht erprobt): (a) Zyklus 31 s ohne Streaming-Signal,
     Fehler `ErrReplication`, Erwartung `ErrTransientExhausted` — Mutation: Messung
     zurück auf `cycleStart`; (b) Zyklus mit Signal und ≥ 30 s Streaming setzt zurück,
     mit 29 s nicht — Mutation: Schwelle 60 s / 10 s; (c) `serverFault`-Tabelle je
     Klasse 08, 40, 53, 55, 57, 58 (bisher 40 und 58 ungebunden, N-3), 25006
     wiederholt, 25P02, XX000, 42704, 3D000 `ErrRejected`, 42501/28000 `permission`
     — Mutation je Präfix entfernen, 25006 entfernen; (d) Fristtest mit einem
     Loopback-Listener, der TCP annimmt und nie antwortet: `NewStream` kehrt
     nach der (im Test verkürzten) Frist mit `ErrReplication` zurück und eine
     spätere `Run` desselben Streams ist von der Frist unberührt (belegt die
     hergeleitete `pgconn`-Aussage, netzlos); (e) ein Test für `Stream.Close`
     und die Schließ-Pfade der Zyklus-Closure (N-4a).
  3. **Träger.** `SPEC-008` Zeile `transient` (`spec/pflichtenheft.md`,
     Festlegung 5 dieser ADR); Handbuch „Neustart nach einem Fehler" („gestreamt"
     bedeutet dann: nach Bestätigung von `START_REPLICATION`; Frist des Aufbaus;
     25006 in der Fehlermenge); die Träger-Tabelle des Slices; Suchlauf
     nach der bewegten Eigenschaft („30 s", „stabil", `streamRetryStableAfter`) an
     beiden Ständen (`AGENTS.md` §3.13). `spec/architecture.md` bleibt frei von
     diesen Bezügen (`AGENTS.md` §3.4).

## Fitness Function (falls maschinell prüfbar)

Die Zeilen sind **Zusagen** des Slices, keine Erprobungen — vor der Umsetzung
gibt es keinen Test. Erprobt ist nur die Ausgangslage: N-1 und die
Mutationen `n1`–`n5` im Review-Report (durch den Reviewer, Kopie im Scratchpad,
Race-Image ohne Netz; **übernommen**). Nicht erprobt: das Verhalten von
`pgconn` bei Fristablauf, alle SQLSTATE-Zuordnungen am Server, jede hier
genannte künftige Mutation.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Fake-Uhr) | Zyklus ohne Streaming-Signal setzt nie zurück, auch bei 31 s; Zyklus mit Signal und ≥ 30 s Streaming setzt zurück, mit 29 s nicht | `make test` |
| Go-Test | `serverFault`-Tabelle je SQLSTATE-Klasse und -Grenzfall (Folgepflicht 2c) | `make test` |
| Go-Test (Loopback-Listener ohne Antwort) | `NewStream` endet nach Fristablauf mit `ErrReplication`; `Run` nach `NewStream` nicht von der Frist gebunden | `make test` |

## Re-Evaluierungs-Trigger

- Für die Schwelle und die Aufbau-Frist (je 30 s): **nicht permanent.** Zeigt
  der Betrieb, dass ein Aufbau regelmäßig länger dauert oder Flapping-Episoden
  das Fenster umgehen, schärft eine Folge-ADR mit `Supersedes ADR-0136`.
- Für die SQLSTATE-Positivliste: **permanent, solange die Menge klein bleibt**;
  ein realer, nicht gelisteter transienter SQLSTATE (zum Beispiel aus einem
  Failover-Lauf) ist der Anlass, ihn per Folge-ADR aufzunehmen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-30 | Accepted | Review-Report zu `slice-capture-transient-wiederholung`, Fixrunde 1, N-1/N-2 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0136` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
