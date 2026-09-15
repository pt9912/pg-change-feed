# ADR-0080: Nähte der pgconn-Adapter — Treiber-Hülle, kein Subjekt-Transfer

**Status:** Accepted — **kein** Supersedes. Diese ADR entscheidet das Design der
zwei aus `slice-081` entstandenen Vorgänge — [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) <!-- d-check:status-provenance -->
Punkt 5 hat es ausdrücklich offengelassen, [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
hat den zweiten Vorgang als „Transfer" **erwartet**. Sie **wendet** dessen
Regel an, statt sie zu ändern, und **bestätigt** beide ADRs im Übrigen.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Planner-Zug, der die
zwei Slices geschnitten hat, und als der Implementer-Lauf von `slice-081`, <!-- d-check:status-provenance -->
dessen Messung hier nachgeprüft wird — Modul 8 §Rollen-Regeln)

**Bezug:** [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Punkt 1/3/4/5 — Gegenstands-Partition, Naht-Zulässigkeit) ·
[`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(Neu-Bemessung bei Transfer) ·
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(der dreiteilige Transfer-Nachweis, hier als **Null-Befund** geführt) ·
[`ADR-0023`](0023-fehlerklassifikation.md) (Fehlerklassen an der Adapter-Grenze) ·
[`ADR-0032`](0032-postgresql-adapterdetail.md) (der Treiber bleibt Adapterdetail) ·
[`ADR-0007`](0007-source-ack-outbound-port.md) (Stream und ACK als getrennte
Rollen an **einer** technischen Verbindung — Option C) ·
[`ADR-0006`](0006-replication-stream-driving-adapter.md) · [`ADR-0049`](0049-replication-fehlerklassen-schwellen.md) ·
`AGENTS.md` §3.5 (Accepted-ADRs sind immutable) · §3.6 (Schwellen nur per ADR) ·
§3.7 (Ist-Zustand) ·
`Dockerfile` (Stufe `coverage`, Paket-Filter) · `tools/harness/db-coverage.sh`
(`DB_COVERAGE_PKGS`, `DB_COVERAGE_THRESHOLD`) · `harness/mk/coverage.mk`
(`THRESHOLD`) · `internal/adapters/driven/postgresstorage/sqlexec/seam.go`
(das Muster der gezogenen Naht — die Zusicherung, die hier **nicht** trägt) ·
`docs/plan/planning/open/slice-084-postgresack-naht.md` <!-- d-check:status-provenance -->
und `docs/plan/planning/open/slice-085-receive-naht.md` <!-- d-check:status-provenance -->
(deren §1 und §4-Start-Trigger dieser Zug erfüllt) ·
`.a-check.yml` (unverändert)

**Schärft:** — (Design-ADR ohne Spec-Stratum; beide Nähte liegen **innerhalb**
der Adapter und berühren weder Vertrag noch Technik noch Sicht)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Zwei geschnittene Vorgänge warten auf genau eine Aussage.** `slice-084` <!-- d-check:status-provenance -->
(`internal/adapters/driven/postgresack`) und `slice-085` <!-- d-check:status-provenance -->
(`internal/adapters/driving/replication/receive`) sagen in §1 je ausdrücklich:
„Das Design entscheidet der Architect beim Start dieses Slice — dieser Plan
zeichnet es **nicht** vor." Ihr §4 führt es als **Start-Trigger**
(`next` → `in-progress`). Dieser Zug entscheidet es; der Trigger ist damit
erfüllt, sobald der Planner den Plan nachgezogen hat (§Konsequenzen).

**(2) Die Fläche ist gemessen, nicht vermutet.** Aufruf-Symbole aus `pgconn`/
`pglogrepl` in den beiden Paketen (am Quelltext gezählt):

| Paket | Symbole | darunter |
|---|---|---|
| `postgresack` (`ack.go`, `ack_test.go`) | **6** | `pgconn.ParseConfig`, `pgconn.ConnectConfig`, `pgconn.PgConn`, `pglogrepl.LSN`, `pglogrepl.SendStandbyStatusUpdate`, `pglogrepl.StandbyStatusUpdate` |
| `receive` (`receive.go`, `walretention.go`, `stream_test.go`) | **18** | `pglogrepl.StartReplication`, `pglogrepl.CreateReplicationSlot`, `pglogrepl.IdentifySystem`, `pglogrepl.ParseXLogData`, `pglogrepl.ParsePrimaryKeepaliveMessage`, `pglogrepl.SendStandbyStatusUpdate`, `pgconn.ErrorResponseToPgError` und weitere |

Die zwei Zahlen sind eine **andere Größenordnung** — sie sind der Grund, warum
es zwei Vorgänge sind und nicht einer (`slice-081` §4). <!-- d-check:status-provenance -->

**(3) Die Lehre aus `slice-081`.** Dort hat sich **compilerseitig** gezeigt, <!-- d-check:status-provenance -->
dass zwei Zusagen sich ausschließen können: `*pgxpool.Pool` erfüllt `Executor`
**genau dann**, wenn `Query` `pgx.Rows` liefert; eine auf vier Methoden verengte
Zeilen-Schnittstelle verträgt sich damit nicht (`wrong type for method Query`).
Der Slice hat **die Zusicherung behalten** (`var _ DB = (*pgxpool.Pool)(nil)`,
`sqlexec/seam.go`) und die Verengung fallen gelassen (Review `review-slice-081` <!-- d-check:status-provenance -->
F-2, Plan §1). Die tragende Frage für jedes weitere Naht-Design ist damit
gestellt: **welcher Teil der Fläche lässt sich überhaupt so schneiden, dass die
Zusicherung hält — und welcher verlangt den konkreten Typ?**

**(4) Dieser Zug hat die Signatur-Frage gemessen — an der gepinnten Bibliothek
und am Compiler.** Ergebnis in zwei Spalten, die sich auf dieser Fläche
**ausschließen**:

| Operation der Naht | erfüllt `*pgconn.PgConn` sie **direkt**? | von einem **Fake** konstruierbar? |
|---|---|---|
| `Exec(ctx, sql) *pgconn.MultiResultReader` | **ja** | **nein** — `*pgconn.MultiResultReader` hat **keinen** exportierten Konstruktor; allein `PgConn.Exec` erzeugt einen |
| `Exec(ctx, sql) ([]*pgconn.Result, error)` | **nein** — `wrong type for method Exec` | **ja** — `pgconn.Result` trägt exportierte Felder (`Rows [][][]byte`, …) |
| `IdentifySystem`, `CreateReplicationSlot`, `StartReplication`, `SendStandbyStatusUpdate` | **nein** — `does not implement … (missing method …)`: die vier sind **Paketfunktionen**, die `conn *pgconn.PgConn` **in der Signatur** verlangen (`pglogrepl.go`) | **ja** — `IdentifySystemResult`, `CreateReplicationSlotResult`, `StandbyStatusUpdate`, `LSN` sind exportierte Werttypen |
| `ReceiveMessage(ctx) (pgproto3.BackendMessage, error)`, `Close(ctx) error` | **ja** | **ja** — `pgproto3.CopyData`/`CopyDone`/`ErrorResponse` sind exportiert |

Gemessen mit `go vet` gegen eine **Wegwerf-Kopie** des Baums im gepinnten
Toolchain-Container (der Arbeitsbaum bleibt unberührt): die **Hülle** (sieben
delegierende Methoden über `*pgconn.PgConn`) und ein **Fake** (nur exportierte
Typen) erfüllen **dieselbe** Schnittstelle — `SEAM-CHECK: OK`; die
Gegenprobe `var _ seam = (*pgconn.PgConn)(nil)` bricht mit
`does not implement … (missing method CreateReplicationSlot)` bzw.
`wrong type for method Exec`.

**Das ist der Befund:** Der einzige Teil der Fläche, den der Treiber **direkt**
erfüllt und der eine Ergebnis-Menge trägt, ist der **unfake-bare**; jede
Operation, die Logik trägt und einen Fake braucht, ist eine **Paketfunktion**,
die der Treiber nicht erfüllt. **Die Zusicherungsform aus `slice-081` ist auf <!-- d-check:status-provenance -->
dieser Fläche nicht zu haben** — die Naht muss eine **Hülle** tragen.

**(5) Was ein Naht-Zug überhaupt bewegen kann.** Der Messgegenstand ist eine
**Paket**-Partition, keine Datei- oder Funktions-Partition:

- Die **netzlos prüfbare Fläche** (Unit) ist `go list ./internal/... ./cmd/...`
  **ohne** die namentlich gefilterten Pakete — `Dockerfile`, Stufe `coverage`:
  `grep -vE '(^|/)(postgresstorage|postgresack|replication/receive)$'`. Der
  `$`-Anker trifft **genau** diese Pakete: ein **Unterpaket** eines
  ausgenommenen Pakets zählt im Unit-Gegenstand (Präzedenz:
  `postgresstorage/sqlexec`, [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  §Kontext (1)).
- Die **DB-Adapter-Coverage** ist in `tools/harness/db-coverage.sh`
  (`DB_COVERAGE_PKGS`) namentlich verankert; die Nenner der zwei hier
  betroffenen Pakete sind **23**
  (`postgresack`) und **155** (`receive`) Statements
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3, [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  §Kontext (1) — zitiert, nicht in diesem Zug gemessen).

Ein „Transfer" im Sinne von [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
setzt deshalb einen **Paketwechsel des Gegenstands** voraus. Eine Naht, die
**im Paket** liegt, bewegt **nichts** — die Logik wird netzlos **testbar**, ohne
den Nenner zu wechseln.

**(6) Die Accepted-ADRs erwarten einen Transfer — die Messung korrigiert die
Erwartung, nicht die Regel.** [`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
§Kontext (5) und §Konsequenzen sowie [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
§Konsequenzen führen als Folge „jeder weiteren Naht" einen schrumpfenden
DB-Gegenstand und damit eine Neu-Bemessung. Das ist eine **Erwartung**, keine
Festlegung — und sie trifft nur, solange die Naht die Logik **aus dem Paket**
zieht. [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
Trigger (a) deckt den anderen Fall bereits ausdrücklich ab: „**Kein** Transfer:
der Einstieg bleibt stehen." Diese ADR muss dort daher nichts berichtigen und
nichts supersedet (anders als `ADR-0078` §Kontext (4) für eine falsche *Zahl*).

## Entscheidung

Wir wählen: **Die Naht ist eine adapter-eigene Treiber-Hülle plus eine schmale,
für einen Fake erreichbare Fläche — und sie bleibt im Paket. Ein
Subjekt-Transfer tritt nicht ein.** Sieben Festlegungen:

1. **Die Nahtform ist Hülle + Schnittstelle, nicht Zusicherung.** Beide Pakete
   bekommen eine **adapter-eigene Schnittstelle**, deren Methoden genau die
   Operationen sind, die der Adapter wirklich ausführt, und deren Typen ein
   **Test konstruieren kann** (exportierte Werttypen, `[]*pgconn.Result`,
   `pgproto3.BackendMessage`). Eine **Hülle** im selben Paket hält die konkrete
   `*pgconn.PgConn` und implementiert jede Methode als **1:1-Delegation**
   (`pglogrepl.X(ctx, c.conn, …)`, `c.conn.Exec(ctx, sql).ReadAll()`). Eine
   Zusicherung der Form `var _ = (*pgconn.PgConn)(nil)` **entfällt** — sie trägt
   hier nicht (gemessen, §Kontext (4)).

2. **Der Schnitt liegt am Treiber, nicht am Adapter-Port.** Die Naht spricht die
   Operationen des Treibers in ihrer **vollen Wertform** aus; die **Übersetzung**
   bleibt darüber und wird dadurch netzlos prüfbar:
   `Acknowledge` (Null-Positions-Grenze, LSN-Form, Fehlerklassen-Wrapping) bei
   `postgresack`; Schleifen-Dispatch, Slot- und Publication-Auflösung,
   Katalog-Zeilen-Übersetzung und Rückstands-Messung bei `receive`.
   **Ausdrücklich nicht** gewählt:
   - eine Naht-Methode in der Sprache des **Anwendungsfalls** (etwa
     `Acknowledge(pos)`) — sie schöbe die Logik in die Hülle und machte sie
     damit wieder unprüfbar;
   - ein **adapter-eigener Draht-Typ** statt `pglogrepl.StandbyStatusUpdate` —
     er verlangte eine Übersetzung, die genau die Information verliert, die
     verlustfrei durchgereicht werden kann. Das wäre der Fehler aus `slice-081` <!-- d-check:status-provenance -->
     in der Gegenrichtung: dort wurde ein **Ergebnis**-Typ verengt, hier würde
     ein **Eingabe**-Typ verengt.

3. **Die konkrete Verbindung bleibt — in Hülle und Dial.** Der konkrete Typ
   bleibt im Paket, und er bleibt **klein**: die Hülle (eine delegierende
   Methode je Naht-Operation) und der Verbindungsaufbau
   (`connectReplication`, Replication-Modus am Treiber). Was ihn nutzt, ist
   ausschließlich der Treiber-Rand; was **darüber** liegt, kennt ihn nicht mehr.

4. **Kein Subjekt-Transfer — der Nachweis wird als Null-Befund geführt.**
   Die Naht liegt in `postgresack` bzw. `receive`; beide Pakete bleiben in der
   DB-Adapter-Coverage, der Unit-Gegenstand bleibt unverändert. Der
   dreiteilige Nachweis aus [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
   hat damit ein **Subjekt ohne Bewegung**:
   **(a) Arithmetik:** `k_ab = 0` — kein Gegenstand gibt Statements ab; der
   Regressions-Riegel (`k_auf < k_ab`) greift nicht, weil keine Neu-Bemessung
   beansprucht wird. **(b) Paket-Diff:** er zeigt, was er zeigen soll — **kein**
   Paketwechsel; der DB-Gegenstand wächst um die Hüllen-Statements (Größe: der
   Lauf misst sie). **(c) kein Verhalten verloren:** `make test-store`,
   `make test-replication` und der Tier-Lauf bleiben grün, **kein Testfall wird
   entfernt** — die Fake-Seite **ergänzt** die realen Tests, sie ersetzt sie
   nicht (`slice-084`/`slice-085` §1). <!-- d-check:status-provenance -->
   Angewandt wird damit [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
   Trigger (a) zweiter Fall: **der Einstieg bleibt stehen.** Die Rampen bleiben
   `DB_COVERAGE_THRESHOLD=70` und `THRESHOLD=70`, die Endstufen je **80 %**;
   eine Schwellen-ADR wird **nicht** fällig.

5. **Je Paket.** Die Naht konkret:

   - **`postgresack` — eine Schnittstelle, eine Methode.**
     `SendStandbyStatusUpdate(ctx, pglogrepl.StandbyStatusUpdate) error`; Hülle
     `connSender{conn *pgconn.PgConn}`. Der **Konstruktions-Vertrag bleibt
     unverändert**: `New(conn *pgconn.PgConn, opts ...Option)` samt
     `nil`-Grenze und die Verdrahtung in der Composition Root
     (`internal/bootstrap/wiring.go`) bleiben, wie sie sind — das ist die
     kleinste Änderung am Bestand und hält die reale Verdrahtung unverändert.
     Für die netzlosen Tests tritt ein **paket-interner** Einstieg hinzu (die
     Schnittstelle ist adapter-eigen und unexportiert; kein neuer öffentlicher
     Rand). Der netzlos prüfbare Teil ist die **Standby-Status-Form** (drei
     LSN-Positionen aus dem Offset), die **Null-Positions-Grenze** und das
     **Fehlerklassen-Wrapping** — je als reine Funktion mit eigenem Test.

   - **`receive` — eine Schnittstelle, sieben Operationen.**
     `IdentifySystem`, `CreateReplicationSlot`, `StartReplication`,
     `SendStandbyStatusUpdate`, `Exec(ctx, sql) ([]*pgconn.Result, error)`,
     `ReceiveMessage(ctx) (pgproto3.BackendMessage, error)`, `Close(ctx) error`;
     Hülle `connSession{conn *pgconn.PgConn}`. Der **öffentliche Rand bleibt
     unverändert** — `Config`, `NewStream`, `Stream.Run/Conn/Assembler/
     BindCapture`, `WALRetentionChecker`, die zwei Fehlerklassen: die
     Composition Root wird **nicht** berührt. Der netzlos prüfbare Teil ist die
     **Empfangs-Schleife** (Byte-ID-Dispatch `XLogData`/`Keepalive`/`CopyDone`/
     `ErrorResponse`, Keepalive-Antwort mit der letzten bestätigten Position,
     Kontext-Ende, leeres `CopyData`), die **Slot-Auflösung** (bestehender Slot
     → `confirmed_flush_lsn`; fehlender → Anlage über `CreateReplicationSlot`),
     die **Publication-Grenze**, die **Katalog-Zeilen-Übersetzung** und die
     **Rückstands-Messung** — fahrbar über einen Fake, der
     `pgproto3.CopyData`/`CopyDone`/`ErrorResponse` liefert.
     `replication/decode` und `replication/mapper` **bleiben außerhalb** des
     Gegenstands: sie liegen bereits im Unit-Gegenstand (Geschwister-Pakete, vom
     Filter nicht getroffen) und die Naht ändert sie nicht.

6. **Zuschnitt und Reihenfolge: `slice-084` zuerst, `slice-085` danach.** <!-- d-check:status-provenance -->
   Keiner der beiden sprengt die Größenregel **in dieser Fassung**: je ein
   Paket, je ein Liefer-Fokus, drei Liefer-Punkte. Die Reihenfolge ist keine
   Code-Abhängigkeit (die Pakete teilen keinen Code — Adapter-→-Adapter-Import
   bleibt ausgeschlossen) — sie ist die **kleinste Probe zuerst**: `084` führt
   die Form an **einer** Hüllen-Methode und **einem** Fake vor, bevor die
   **18-Symbol-Fläche** von `085` auf dieselbe Form baut. Die Rückführung
   `in-progress` → `next`, die `slice-085` §4 vorab benennt, wird auf **diesen** <!-- d-check:status-provenance -->
   Fall geschärft: sie greift, wenn die Naht **einen Paket- oder API-Umzug**
   verlangt (Fassade, Composition Root, Gegenstands-Liste) — **nicht**, wenn
   `decode`/`mapper` „mit hineinmüssen" (das ist hiermit entschieden: sie bleiben
   draußen).

7. **`.a-check.yml` bleibt unverändert.** Beide Pakete liegen unter
   `internal/adapters/**` und damit in der Schicht `adapters`; die Naht liegt
   **innerhalb** dieser Schicht, die Kanten (`adapters → ports`, `adapters →
   domain`) werden weder neu gezogen noch verletzt. Auch ein **Unterpaket**
   bliebe in `adapters` — die Wahl dieser ADR braucht aber keines.

### Warum die Hülle kein Rückfall auf `slice-081`s verworfene Option D ist <!-- d-check:status-provenance -->

`slice-081` hat einen „Vermittler" verworfen, der `pgx.Rows` auf vier Methoden <!-- d-check:status-provenance -->
eindampft. Diese ADR wählt eine Hülle — und das ist **nicht** derselbe Fall:

| | `slice-081`s verworfener Vermittler | die Hülle dieser ADR | <!-- d-check:status-provenance -->
|---|---|---|
| was sie tut | **verengt** einen gelieferten Typ (`pgx.Rows` → 4 Methoden) und verliert den Rest | **delegiert 1:1** — jede Methode ist genau der Bibliotheksaufruf, den der Adapter sonst direkt machte |
| was sie kostet | eine zusätzliche Schale **und** den Verlust der Pool-Zusicherung | die Hülle **selbst** (sieben delegierende Methoden) |
| was sie erhält | nichts, was die Zusicherung nicht schon erhält | alles: die Ergebnis-Werttypen werden unverändert durchgereicht, jede Prüfung der Bibliothek (etwa die LSN-Analyse) läuft weiter |

Die Hülle ist hier also keine Verengung, sondern die **einzige** Form, in der die
zwei Eigenschaften „der Fake kann Treiber spielen" und „der konkrete Typ bleibt
draußen" gleichzeitig zu haben sind — sie ist in §Kontext (4) gemessen, nicht
behauptet.

### Bestätigt im Buchstaben, widerlegt in der Folge — der Verdacht aus `slice-085` §6 <!-- d-check:status-provenance -->

`slice-085` §6 führt als Risiko: „`pglogrepl.StartReplication` und <!-- d-check:status-provenance -->
`CreateReplicationSlot` verlangen den **konkreten** `*pgconn.PgConn` — eine
schmale Schnittstelle drückt das womöglich nicht aus."

- **Bestätigt:** es sind **vier** Paketfunktionen (nicht zwei), die den
  konkreten Typ in der Signatur tragen, und `*pgconn.PgConn` kann für sie
  **keine** Schnittstelle erfüllen (gemessen, §Kontext (4)). `slice-081`s <!-- d-check:status-provenance -->
  Zusicherungsform ist auf dieser Fläche nicht zu haben.
- **Widerlegt:** daraus folgt **nicht**, dass die Naht nicht trägt. Der konkrete
  Typ wird nicht *ersetzt*, sondern **eingeschlossen** — die Hülle erfüllt die
  Naht (gemessen: `SEAM-CHECK: OK`), und alles oberhalb ist netzlos fahrbar.
  Der Blockerfall aus `slice-085` §4 („ohne Verhaltensänderung nicht zu haben") <!-- d-check:status-provenance -->
  tritt damit **nicht** ein.

### Die benannte Grenze: der Unit-Nenner sieht die neue Prüfbarkeit nicht

Diese Naht macht Logik netzlos **testbar**, aber sie macht sie nicht zum Teil
des **Unit-Gegenstands**: die Statements liegen weiter in den zwei
ausgenommenen Paketen. [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 beschreibt als Wirkung der Naht, sie ziehe die Logik „in die prüfbare
Fläche und hebt die Decke" — für diese zwei Pakete trifft das **nicht** zu, weil
ihre Naht im Paket bleibt. Das ist die Grenze, nicht ein Versehen:

- **Der Unit-Nenner bekommt nichts.** Wer den Nutzen der Naht in der
  Unit-Zahl sucht, sucht ihn an der falschen Stelle; ihr Nutzen ist das
  **Design** ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 5: „Ihre Rechtfertigung ist Design … nicht die Zahl").
- **Der DB-Nenner wird leicht verdünnt.** Die neuen netzlosen Tests laufen im
  Messlauf der DB-Adapter-Coverage mit (dasselbe `go test`) und decken
  Statements, die keine echte PostgreSQL-Instanz berührt haben — die Zahl
  heißt weiterhin richtig „Coverage des DB-Gegenstands", aber ein wachsender
  Teil davon stammt aus Fakes. Diese Verdünnung existiert **heute schon** in
  kleinem Maß (`stream_test.go` trägt einen netzlosen Fall); die Naht
  vervielfacht sie. **Trigger:** wird sie materiell, ist die design-honeste
  Antwort die **Verlagerung** des netzlos prüfbaren Stücks in ein gezähltes
  (Unter-)Paket — als eigener, design-getragener Vorgang (das Muster
  `postgresstorage/sqlexec`), **nicht** als Beigabe zu dieser Naht. Eine
  Verlagerung, deren einziger Zweck die Zahl ist, wäre genau das Goodhart, das
  [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Option D verwirft.

### Was diese ADR nicht entscheidet

- **Nicht den Lifecycle der zwei Slices** — `Verantwortlich:`, `open` → `next`
  → `in-progress` und die WIP-Frage sind Planner-Arbeit (Modul 5). Dieser Zug
  liefert die Design-Hälfte des Start-Triggers.
- **Nicht die Bezeichner** — Namen von Schnittstelle, Hülle und Einstieg sind
  Implementer-Arbeit; die *Form* ist entschieden (Festlegung 1/5).
- **Nicht die Rampen-Zahlen** — kein Transfer, keine Neu-Bemessung
  (Festlegung 4); sollte ein Lauf wider Erwarten einen **Paketwechsel**
  messen, gilt [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  unverändert und die drei Belege werden Bringschuld jenes Laufs.
- **Nicht die Verlagerung in ein gezähltes Paket** (die Verdünnungs-Antwort) —
  sie hat hier einen **Trigger**, keinen Beschluss.
- **Nicht die Frage, ob `postgresack` ganz aus dem DB-Gegenstand fällt** (der
  reale Test zöge um, die Liste in `Dockerfile`/`db-coverage.sh` wäre
  nachzuziehen): das ist eine Gegenstands-Frage
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Trigger (a)) und für 23 Statements der falsche Preis in diesem Zug
  (§Verglichene Alternativen, Option E).
- **Nicht eine Berichtigung an [`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)/[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)** —
  ihre Festlegungen bleiben unverändert gültig; ihre **Erwartung** eines
  weiteren Transfers trifft für diese zwei Pakete nicht (kein Supersedes, kein
  Zitat-Eingriff; §Kontext (6)).
- **Nicht die Wellen-Zuordnung** — beide Slices sind ohne Welle geschnitten
  (ihre Closure-Bedingung ist die eigene DoD); der Sichtungs-Schritt ihres §8
  bleibt Planner-Sache.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Zusicherungsform wie in `slice-081`: `*pgconn.PgConn` erfüllt die Naht direkt | der schönste Beleg: der reale Treiber trägt die Schnittstelle, kein Vermittler | **gemessen unmöglich** (§Kontext (4)): für die vier `pglogrepl`-Operationen existiert keine Methode am Typ; und die einzige direkt erfüllbare `Exec`-Form (`*pgconn.MultiResultReader`) ist **nicht fake-bar** | <!-- d-check:status-provenance -->
| B — Naht in der Sprache des Anwendungsfalls (etwa `Acknowledge(pos)` / `Receive()`) | die kleinste denkbare Schnittstelle am Adapter | schiebt die Logik **in** die Hülle: die reine Form, die Null-Positions-Grenze und die Fehlerklassifikation wären wieder nur hinter einer lebenden Verbindung erreichbar — die Naht hätte ihr Ziel verfehlt |
| **C — Hülle + schmale, fake-fähige Fläche, im Paket, kein Subjekt-Transfer (gewählt)** | erreicht das Design-Ziel („schmale Abhängigkeit statt konkreter Typ") vollständig und **verlustfrei**; lässt Composition Root, öffentlichen Rand, Gegenstands-Listen, Schwellen und Rampen unberührt; beide Slices bleiben klein und einzeln lieferbar | der Unit-Nenner sieht die neue Prüfbarkeit nicht, und der DB-Nenner verdünnt leicht (§„Die benannte Grenze") — beides benannt, beides mit Trigger |
| D — C **plus** Verlagerung des netzlos prüfbaren Stücks in ein gezähltes (Unter-)Paket | realer, zählbarer Transfer; die zwei Gegenstände bleiben subjektiv ehrlich getrennt (das `sqlexec`-Muster) | für `receive` verlangt sie einen **Fassaden- und/oder API-Umzug** (Composition Root, Gegenstands-Liste) — der Slice-Schnitt bricht; für `postgresack` wäre es ein Paket für wenige Statements, und die Rechtfertigung wäre die **Zahl**, also der Goodhart-Fall aus [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Option D |
| E — `postgresack` **ganz** aus dem DB-Gegenstand nehmen (realen Test umziehen, Paketliste nachziehen) | 23 Statements wandern real in den Unit-Nenner; die Paket-Zuordnung folgt dann exakt der Eigenschaft („Testlauf braucht keinen Dienst") | bewegt für 23 Statements die Gegenstands-Liste **und** die Kalibrierungsbasis der Unit-Rampe; der Test-Umzug ist der tragende Teil — und der reale ACK-Beleg liegt bereits stärker in `internal/bootstrap`; der Preis steht in keinem Verhältnis zum Ertrag. Bleibt als benannter Folgeweg mit Trigger (Verdünnung), nicht als Teil dieses Zugs |
| F — nichts tun (die zwei Slices bleiben im `open/`, ihr Start-Trigger unerfüllt) | kein Aufwand | die zwei Slices sind aus einer **Abweichung** von `slice-081` entstanden und ausdrücklich als eigener, design-begründeter Vorgang geschnitten ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 5); sie ohne Entscheidung liegen zu lassen, stellt den Start-Trigger nie her | <!-- d-check:status-provenance -->

**Fazit:** C. A ist gemessen unmöglich, B verfehlt das Ziel, D und E bezahlen mit
dem Slice-Schnitt bzw. mit der Gegenstands-Liste für einen Ertrag, dessen
Rechtfertigung die Zahl wäre, F lässt den Auftrag offen. C erreicht die
Design-Hälfte vollständig, hält die Zähl-Hälfte unberührt und macht die
Nicht-Bewegung zu einem **geführten** Befund statt zu einem stillen.

## Konsequenzen

- Positiv: Beide Adapter hängen an einer **adapter-eigenen** Fläche statt an der
  vollen `*pgconn.PgConn`-API; ihre Übersetzung (Fehlerklassen, LSN-Form,
  Slot-Auflösung, Meldungs-Dispatch, Zeilen-Übersetzung) wird **netzlos**
  fahrbar. Der Fake prüft die **Verklebung** — die realen Tests bleiben der
  Wächter des Protokolls (`slice-084`/`slice-085` §1). <!-- d-check:status-provenance -->
- Positiv: Der Start-Trigger beider Slices ist mit dieser ADR erfüllt; die
  Composition Root, der öffentliche Rand, `Dockerfile`-Filter, `DB_COVERAGE_PKGS`,
  beide Schwellen und `.a-check.yml` bleiben **unverändert**.
- Positiv: Der Nachweis wird trotzdem geführt — als **Null-Befund**: der
  Paket-Diff zeigt „kein Paketwechsel", statt dass jemand einen Transfer
  behauptet, den die Arithmetik nicht deckt.
- Negativ mit Grenze: **Der Unit-Gegenstand wächst nicht** (§„Die benannte
  Grenze"). Wer die Naht über die Zahl rechtfertigt, findet sie dort nicht —
  und die DB-Zahl zählt künftig einen kleinen Anteil fake-gedeckter Statements
  mit. Beides ist benannt.
- Negativ mit Grenze: **Die Hülle ist Code, den es ohne die Naht nicht gäbe**
  (eine delegierende Methode je Operation). Sie ist der Preis dafür, dass der
  konkrete Typ nicht durch die Logik reist; der Paket-Diff macht sie sichtbar.
- Folgepflicht (Planner-Zug): **die zwei Slice-Pläne nachziehen** — §1 (die
  Nahtform benennen und den `Acknowledge`/„Minimal-Typ"-Vorgriff ersetzen),
  §2 Liefer-Punkt 1 (nicht „statt der konkreten `*pgconn.PgConn`", sondern „die
  **Logik** hängt an der Schnittstelle; der konkrete Typ bleibt in Hülle und
  Dial") und Liefer-Punkt 3 (der Nachweis als **Null-Befund** statt als
  erwarteter Transfer), §3 (Datei-Zuschnitt: Schnittstelle, Hülle, reine
  Funktionen, Fakes — **kein** Unterpaket, **kein** Composition-Root-Zugriff),
  §4 (`slice-085`: Rückführung auf „Paket-/API-Umzug" schärfen, nicht auf <!-- d-check:status-provenance -->
  `decode`/`mapper`), §6 (`slice-085`: das Naht-Risiko ist mit dieser ADR <!-- d-check:status-provenance -->
  beantwortet — bestätigt im Buchstaben, widerlegt in der Folge), §8
  (`slice-085`: das „Evidenz-/Diskrepanz-Risiko hoch" hing an genau diesem <!-- d-check:status-provenance -->
  Verdacht und ist nachzuziehen).
- Folgepflicht (Implementer-Zug, je Slice, **kein** Produkt-Code außerhalb der
  zwei Pakete): Schnittstelle + Hülle + paket-interner Einstieg für die
  netzlosen Tests; reine Funktionen mit eigenen Tests; Fakes, die die
  **Verklebung** fahren (Aufruf mit der richtigen Meldung, Fehlerpfad,
  Dispatch, Leerfall); die **realen** Tests bleiben unverändert; der
  Paket-Diff wird als Null-Befund belegt.
- Folgepflicht (Closure): die Nicht-Bewegung wird in der Closure-Notiz des
  jeweiligen Slice als Nachweis-Ausgang geführt („kein Transfer"), nicht als
  ausgefallenes DoD-Item.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Compiler (Zusicherung **im Paket**, nicht im Test) | Die Hülle erfüllt die Naht: `var _ <seam> = connSender{}` bzw. `connSession{}`. Ändert `pglogrepl` eine der vier Paketfunktions-Signaturen, bricht dieser Bau — und nicht erst ein Adapter-Aufruf | `make test` |
| `make test-store` / `make test-replication` (Tier-Phase) | **Kein Verhalten verloren:** die realen ACK-/Stream-Läufe bleiben grün, kein Testfall entfernt; die Fakes sind Zusatz, kein Ersatz | `make test-replication` |
| (Disziplin, kein Sensor) | Der **Null-Befund** (Paket-Diff: kein Subjektwechsel; `k_ab = 0`) ist Bringschuld des Laufs und Review-Gegenstand — die aggregierte Summe bleibt als Beleg entwertet ([`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md) §Konsequenzen) | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Fünf benannte Trigger, sonst permanent:

**(a)** Die **Verdünnung** des DB-Nenners wird materiell (ein wachsender Anteil
fake-gedeckter Statements bei unverändertem Gegenstand) — dann ist die
Verlagerung in ein gezähltes (Unter-)Paket als eigener, design-getragener
Vorgang zu entscheiden (Option D), mit dem Transfer-Nachweis aus
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md).

**(b)** Ein ausgenommenes Paket hat **keinen** Testlauf mehr, der einen externen
Dienst voraussetzt (etwa weil der reale Test umgezogen ist) — dann greift
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Trigger (a): die namentliche Liste in `Dockerfile` und `DB_COVERAGE_PKGS` ist
nachzuziehen, und die **Unit-Rampe** ist gegen ihren neuen Nenner zu prüfen.

**(c)** `pgx`/`pglogrepl` bietet eine **Schnittstelle** für die vier
Treiber-Operationen oder einen exportierten Konstruktor für
`*pgconn.MultiResultReader` — dann wird die Zusicherungsform aus `slice-081` <!-- d-check:status-provenance -->
auf dieser Fläche wieder möglich und die Hülle kann auf die Teile schrumpfen,
die der Treiber dann selbst trägt.

**(d)** Ein **dritter** Adapter braucht dieselbe Naht — dann ist die Form zu
verallgemeinern (gemeinsame Schnittstelle), **ohne** Adapter-→-Adapter-Import
(`ADR-0002`/`ADR-0031`, `.a-check.yml`).

**(e)** Ein Lauf misst wider Erwarten einen **Paketwechsel** (der Implementer
verlagert doch) — dann gilt [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
unverändert, und die drei Belege (Ankunft, Paket-Diff, kein Verhalten
verloren) werden Bringschuld dieses Laufs.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: `slice-084`/`slice-085` §1 („Das Design entscheidet der Architect beim Start dieses Slice") und ihr §4-Start-Trigger. Entscheidet die Nahtform (**Hülle + fake-fähige Fläche, keine Zusicherung** — die Zusicherungsform ist auf dieser Signatur-Fläche gemessen unmöglich), den **Null-Subjekt-Transfer** (`ADR-0078` Trigger (a) zweiter Fall) und die Reihenfolge `084` → `085`; bestätigt `ADR-0071` Punkt 5 im Buchstaben und korrigiert die **Erwartung** eines weiteren Transfers, ohne zu supersedet | eigene Signatur- und Compiler-Gegenprobe (Wegwerf-Kopie, gepinntes Toolchain-Image), `docs/plan/planning/open/slice-084-postgresack-naht.md`, `docs/plan/planning/open/slice-085-receive-naht.md` | <!-- d-check:status-provenance -->

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0080` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
