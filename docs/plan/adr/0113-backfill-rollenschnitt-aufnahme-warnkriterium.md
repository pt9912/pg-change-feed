# ADR-0113: Backfill — Rollenschnitt der `queued`-Zeile, Aufnahme beim Start, Warn-Kriterium für große Tabellen (Supersedes ADR-0111, Teilfrage 5 teilweise)

**Status:** Proposed — Supersedes [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
in **Teilfrage 5, teilweise** (genau zwei Absätze, siehe §Entscheidung);
alles Übrige von `ADR-0111` bleibt in Kraft

**Datum:** 2026-09-23

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch zwei Start-
Bedingungen der geplanten Backfill-Slices — die Schreib-Rolle der
`queued`-Zeile und das Warn-Kriterium für große Tabellen — sowie durch die
Auftraggeber-Vorentscheidungen vom 2026-09-23, die diese ADR begründet und
am Code prüft)

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill des
Bestands; Negative: Neubeginn ohne Verlust),
[`LH-FA-ADM-001`](../../../spec/lastenheft.md) (Administration über SQL),
[`LH-FA-SST-003`](../../../spec/lastenheft.md) (CLI-Diagnose),
[`LH-QA-SEC-001`](../../../spec/lastenheft.md) (Least-Privilege),
[`LH-QA-SEC-002`](../../../spec/lastenheft.md) (getrennte Berechtigbarkeit),
[`LH-QA-SEC-003`](../../../spec/lastenheft.md) (Zugriff ohne Berechtigung
scheitert), [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
(teilweise superseded — Haupt-Bezug),
[ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md) (Rollen je DSN),
[ADR-0050](0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue), [ADR-0034](0034-ports-nach-faehigkeiten.md) (Ports nach
Fähigkeiten), [ADR-0028](0028-inbound-use-cases.md) (Inbound Use Cases),
[ADR-0040](0040-clockport.md) (`ClockPort`),
[ADR-0023](0023-fehlerklassifikation.md) (Fehlerklassen),
[ADR-0054](0054-coverage-gate-und-benchmark-infrastruktur.md) (Bench-
Infrastruktur), [ADR-0071](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(nur zur Abgrenzung)

**Schärft:** [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (die
Antwort auf Auslösung und Sichtbarkeit wird um Rollenschnitt, Aufnahme und
Warnungen ergänzt), [`SPEC-019`](../../../spec/pflichtenheft.md) (`applied`
heißt für die Antragsart `backfill` „angenommen", die Run-Zeile entsteht in
derselben Transaktion), die von `slice-backfill-spec-nachzug` vergebene
`SPEC-*`-Kennung für `cdc.backfill_run`/`cdc.backfill_status` (Grants,
Warn-Spalten, Bedeutung von `estimated_rows` = `NULL`),
[`ARC-005`](../../../spec/architecture.md) (View `cdc.backfill_status`),
[`ARC-006`](../../../spec/architecture.md) (Adapter-Rolle der Annahme)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Drei Stellen seiner Teilfrage 5 und ihrer
Kanten tragen eine Lücke, die vor dem Start der Umsetzung entschieden sein
muss; keine ist eine Zitat-Korrektur, weil jede den Referenten ändert.

### Befunde am Bestand (jeweils mit Anker)

- **Zwei Rollen, ein Text.** `ADR-0111` Teilfrage 5 lässt die
  Administrations-Goroutine die Run-Zeile `queued` anlegen und den Antrag
  `applied` vermerken; im Absatz „Verbindungen und Rollen" nennt er
  `SELECT/INSERT/UPDATE` auf `cdc.backfill_run` ausschließlich für den Pool
  aus `CDC_CAPTURE_DSN`. Die Administrations-Goroutine liest ihre Anträge
  über `postgresstorage.NewAdministrationRequest(ctx, cfg.AdminDSN, …)`
  (`internal/bootstrap/wiring.go`) — Rolle `cdc_admin` (Zuordnung
  `ADR-0047`, Tabelle „Verdrahtung"). `tools/schema/nacharbeit-roles.sql`
  vergibt an `cdc_admin` DML auf vier Verwaltungstabellen und
  `SELECT`/`DELETE` auf `cdc.transaction`/`cdc.change`; `cdc.backfill_run`
  existiert noch nicht (Suche `grep -rn backfill_run tools internal cmd`
  am 2026-09-23: kein Treffer). Der Text von `ADR-0111` wörtlich umgesetzt
  ließe die Goroutine mit einem Recht scheitern, das nur der Worker-Pool
  trägt — **erwartet**, im Store-Tier zu belegen (§Fitness Function).
- **Annahme und Zeile sind heute zwei getrennte Aufrufe.**
  `processAdministrationRequests` (`internal/bootstrap/wiring.go`) ruft
  `applyAdministrationRequest` und danach `deps.requests.MarkApplied` — zwei
  Statements über den Pool, keine gemeinsame Transaktion. Bei den
  bestehenden Antragsarten trägt die Idempotenz der Use Cases die Lücke
  dazwischen (Kommentar an `applyAdministrationRequest`: derselbe Antrag
  trifft nach einem Neustart dieselbe Zeile); ein Run hat diese
  Idempotenz nicht, weil ein zweiter Antrag bei aktivem Run als `failed`
  endet (`ADR-0111` Teilfrage 4).
- **Die Transaktion ist am Adapter machbar.** `AdministrationRequestAdapter`
  hält seinen Pool als `sqlexec.DB`, dessen Schnittstelle `Begin` trägt
  (`internal/adapters/driven/postgresstorage/sqlexec/seam.go`);
  `MarkApplied` läuft heute als `UPDATE … WHERE … AND status = 'pending'`
  (`queries.UpdateAdministrationRequestApplied`), ein zweiter Aufruf trifft
  keine Zeile und ist kein Fehler. Der Port `AdministrationRequestPort`
  (`internal/application/port/outbound/administrationrequest.go`) kennt nur
  `ListPending`, `MarkApplied`, `MarkFailed` — er trägt weder eine
  Run-Zeile noch eine Transaktion über beide Tabellen.
- **Lesen läuft über `cdc_reader`.** `bootstrap.Diagnose` erhält die DSN
  der Lese-Rolle (`cmd/pg-change-feed/main.go`: `cfg.ReaderDSN`); eine
  Warnung, die `diagnose` zeigen soll, muss deshalb in einer View oder
  Tabelle stehen, die `cdc_reader` lesen darf.
- **Der Prozessstart kennt nur `running → interrupted`.** `ADR-0111`
  Teilfrage 4 sagt „`queued`-Runs laufen weiter" und beschreibt keine
  Regel, wer sie wann aufnimmt. Das Muster der Schleife steht im
  Bestand: `runAdministration` verarbeitet zuerst und wartet dann auf das
  nächste Signal (`internal/bootstrap/wiring.go`).
- **Das Warn-Kriterium fehlt.** `ADR-0111` Festlegung 3 fordert „eine
  Warnung, keine Ablehnung" und eine **zu messende** Richtgröße; welche
  Größe wann warnt, steht nirgends. Die Beobachtung
  `BEO-PGC/geschaetzter-wert-als-grenze` beschreibt, wie ein geschätzter
  Wert über mehrere Träger zur Grenze wird — genau die Form, die
  `pg_class.reltuples` (eine Schätzung) hier annehmen kann.

### Konstraints

- `ADR-0111` Teilfrage 3/4 (Position `X`, Ein-Transaktions-Form,
  Fail-closed vor dem Commit), Teilfrage 5 (Antragsart, SQL-Funktion
  schreibt nur den Antrag, Diagnose, Fehlerklassen, Wecksignal) und
  Festlegungen 1–3 bleiben unberührt.
- `ADR-0050`: SQL-Funktionen schreiben ausschließlich einen Antrag und
  senden `pg_notify`.
- Kein Gate wird gelockert, keine Schwelle gesenkt (`AGENTS.md` §3.6).

## Entscheidung

Wir wählen **drei Festlegungen**: (1) `cdc_admin` legt die `queued`-Zeile in
**derselben Transaktion** an, die den Antrag auf `applied` setzt; die
Grants folgen daraus. (2) Der Backfill-Worker nimmt `queued`-Zeilen beim
Prozessstart und bei jedem Wecksignal in Antragsreihenfolge auf. (3) Das
Warn-Kriterium für große Tabellen ist ein **Startpunkt ohne Messung**:
Toleranz 10 Minuten Kopierdauer als Konstante im Code, zwei Warnungen,
Richtgröße in Zeilen erst aus der Messung.

**Was diese ADR von `ADR-0111` ersetzt** — genau zwei Absätze seiner
Teilfrage 5: (a) im Absatz „Ausführung" der Satzteil, der die Run-Zeile
`queued` anlegen und den Antrag `applied` vermerken lässt, ohne die
Zusammengehörigkeit beider Schritte zu binden (gilt jetzt: Festlegung 1);
(b) im Absatz „Verbindungen und Rollen" die Grant-Angabe
`SELECT/INSERT/UPDATE` des Worker-Pools (gilt jetzt: `SELECT/UPDATE`, das
`INSERT` trägt `cdc_admin`). Teilfrage 4 (Satz „`queued`-Runs laufen
weiter") und Festlegung 3 des Auftraggebers werden **geschärft**, nicht
ersetzt (Festlegung 2 und 3 dieser ADR). Alles Übrige — auch die Aussage,
dass `cdc_reader` `SELECT` auf die View `cdc.backfill_status` bekommt —
bleibt gültig und wird hier nicht wiederholt.

Nicht Gegenstand: die Gate-Zuordnung des Snapshot-Adapters (ein neues
DB-Paket) — `ADR-0071` Re-Evaluierungs-Trigger (a) regelt sie, die Regel
greift ohne Textänderung, nachzuziehen sind die namentlichen Listen; und die
Ortsfrage der Replay-Invariante.

### Festlegung 1 — Rollenschnitt der `queued`-Zeile

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: `ADR-0111` wörtlich (Administrations-Goroutine legt an, Grants nur an `cdc_capture`) | keine Änderung | nicht lauffähig: die Goroutine läuft als `cdc_admin` ohne Recht auf die Tabelle (Befund oben, erwartet: `permission denied`); der Start-Trigger des Run-Use-Case-Slice bleibt offen |
| B — der Worker legt die Zeile an (`cdc_capture` mit `INSERT`) | passt zum Wortlaut der Grant-Angabe | zwei Rollen und zwei Pools ohne gemeinsame Transaktion mit dem Antragsvermerk: Vermerk zuerst → `applied` ohne Run nach einem Absturz; Run zuerst → nach einem Absturz `pending` plus `queued`-Zeile, die Wiederholung scheitert an der Regel „kein zweiter aktiver Run" (`ADR-0111` Teilfrage 4) und der Antrag endet `failed`, obwohl der Run läuft; zusätzlich schreibt die Datenpfad-Rolle Verwaltungszustand (`LH-QA-SEC-001`/`-002`) |
| C — beide Rollen mit `INSERT` | löst den Wortlaut-Konflikt | Rechte-Aufweichung ohne Nutzen: jede Rolle könnte Runs erzeugen, die Atomarität aus B bleibt ungelöst |
| D — die SQL-Funktion `cdc.backfill_table` legt Antrag **und** Run atomar an (`SECURITY DEFINER`) | strukturell atomar, keine Port-Änderung | die Vorbedingungen (Tabelle aktiviert und gebunden, Schätzung) prüft Go; die Run-Zeile entstünde vor der Prüfung, ein Vorbedingungs-Fehler hinterließe einen `queued`-Run ohne Ausführbarkeit; widerspricht `ADR-0111` („schreibt nur den Antrag") und der `ADR-0050`-Fitness-Function; zweiter Träger derselben Aussage in SQL |
| **E — `cdc_admin` legt an, in einer Transaktion mit dem Antragsvermerk (gewählt)** | Annahme ≙ Run existiert: ein Absturz hinterlässt entweder nichts (Antrag bleibt `pending`, dieselbe Verarbeitung wiederholt sich) oder Antrag `applied` samt `queued`-Zeile; getrennte Berechtigbarkeit bleibt (die Datenpfad-Rolle legt nichts an); Vorbedingungs-Fehler hinterlassen keine Zeile | ein neuer Port-Schnitt und eine Transaktion im Administrations-Adapter |

**Festlegung.**

1. **`cdc_admin` legt die Run-Zeile mit Status `queued` an**, und zwar in
   **derselben Transaktion**, die den Antrag von `pending` auf `applied`
   setzt. „`applied`" heißt für die Antragsart `backfill` **angenommen**;
   die Ausführung steht in `cdc.backfill_run` (`ADR-0111`).
2. **In derselben Transaktion** steht die Prüfung „kein aktiver Run
   (`queued`/`running`) derselben Tabelle" (Lesen vor Einfügen; ein aktiver
   Run endet den Antrag als `failed` mit Grund, `ADR-0111` Teilfrage 4) — der
   Antragsvermerk muss genau **eine** `pending`-Zeile treffen, sonst
   Rollback und Fehler. Eine Sperre zwischen mehreren Annehmenden ist nicht
   nötig: `processAdministrationRequests` verarbeitet Anträge sequenziell in
   **einer** Goroutine, und je Quelle nimmt eine Instanz Anträge an
   (`ADR-0111` Teilfrage 4, `ADR-0050`-Muster; Annahme, zu belegen im
   Store-Tier-Test der Antragsverarbeitung).
3. **Grants auf `cdc.backfill_run`** (Datei `tools/schema/nacharbeit-roles.sql`,
   Zuordnung der Anschlüsse nach `ADR-0047`):

   | Rolle | Recht | Träger |
   |---|---|---|
   | `cdc_admin` | `SELECT`, `INSERT` | Annahme-Transaktion (Prüfung „kein aktiver Run", Anlage `queued`) |
   | `cdc_capture` | `SELECT`, `UPDATE` | Worker: `queued` lesen, Statuswechsel, `started_at`/`finished_at`, `snapshot_position`, `rows_copied`, `error_message`, Warnung; Start-Abgleich `running → interrupted` |
   | `cdc_reader` | **kein** Recht auf die Basistabelle; `SELECT` auf die View `cdc.backfill_status` | Lesezugriff und `diagnose` (Definer-Semantik der Views, wie bei den übrigen Views dieser Datei) |

   Niemand trägt `DELETE`; `cdc_admin` trägt kein `UPDATE`, `cdc_capture`
   kein `INSERT`.
4. **Port-Schnitt (Arbeitsnamen; nicht Teil dieser Entscheidung, zu
   implementieren im Run-Use-Case- und Run-Store-Slice).**
   - Ein **neuer Outbound-Port** für die Fähigkeit „Antrag annehmen"
     (`BackfillAdmissionPort`, eine Methode `Admit(ctx, requestID, run)`),
     nach `ADR-0034` als eigener Fähigkeits-Port und **nicht** als vierte
     Methode des `AdministrationRequestPort` — jener trägt die
     Antrags-Queue aller Antragsarten. Sein Adapter liegt im Paket
     `postgresstorage` über den **Pool der Administrations-Goroutine**
     (`CDC_ADMIN_DSN`) und nutzt `Begin` der Naht `sqlexec.DB`: Prüfung →
     `INSERT` → `UPDATE … applied` → Commit, jede Abweichung Rollback.
   - Der **Run-Zustands-Port** des Worker-Pools (Pool aus `CDC_CAPTURE_DSN`)
     hat **keine** Operation „anlegen"; er trägt Übergänge, Fortschritt,
     Abschluss, `queued` der eigenen Quelle in Antragsreihenfolge lesen und
     den Abgleich `running → interrupted`.
   - `BackfillTableUseCase.Request` ruft `Admit` als **letzten** Schritt,
     nach den Vorbedingungen und der Schätzung. Die bestehende Schleife
     `processAdministrationRequests` bleibt unverändert lauffähig: ihr
     anschließendes `MarkApplied` trifft wegen der `WHERE`-Klausel keine
     `pending`-Zeile mehr und ist kein Fehler (Anker: die Idempotenz-Klausel
     oben).

**Akzeptiertes Negativ.** Ein `INSERT`-Recht auf Tabellenebene lässt
`cdc_admin` einen beliebigen `status` einfügen; nur der Code beschränkt die
Anlage auf `queued`. Die Rolle ist die Verwaltungsrolle des Prozesses; der
Schutz, den `LH-QA-SEC-001`/`-002` verlangen (Trennung nach Rolle), bleibt
erhalten, ein Zeilenfilter innerhalb einer Rolle wäre Aufwand ohne
benannten Angreifer.

### Festlegung 2 — Aufnahme beim Start und bei Wecksignal

| Option | Pro | Contra |
|---|---|---|
| A — nur der Start-Abgleich von `ADR-0111` (`queued` nur durch Übergabe) | keine Änderung | eine `queued`-Zeile überlebt einen Neustart, wird aber nie ausgeführt: stiller Stillstand, sichtbar nur als ewig `queued` |
| B — periodischer Poll auf `queued` | einfach, kein Signal | ein weiterer Ticker mit Latenz für ein Ereignis, das der Prozess ohnehin kennt |
| C — eigene `LISTEN`-Verbindung des Workers auf einem DB-Kanal | entkoppelt vom Administrations-Pfad | zweite dedizierte Verbindung samt Wiederverbindungslogik (wie `AdministrationListener`) für ein Signal, das nur diese Instanz erzeugt |
| **D — Aufnahme beim Prozessstart plus in-Prozess-Wecksignal, Abarbeiten vor Warten (gewählt)** | deckt Neustart und Neuannahme mit einem Muster, das der Bestand kennt (`runAdministration`); kein zweiter Kanal | setzt voraus, dass nur die Administrations-Goroutine dieser Instanz `queued`-Zeilen anlegt (Festlegung 1) |

**Festlegung.**

1. **Reihenfolge und Ein-Worker-Semantik von `ADR-0111` bleiben**: eine
   Goroutine, ein Run zugleich, Aufnahme in Antragsreihenfolge —
   `ORDER BY requested_at, run_id` (der Zweitschlüssel ist eine Setzung, weil
   `ListPending` heute nur nach `requested_at` ordnet).
2. **Aufnahme-Ereignisse:** der Prozessstart und jedes Wecksignal. Beim
   Prozessstart läuft zuerst der Bindungsaufbau der Composition Root, dann
   der Abgleich `running → interrupted`, dann startet der Worker und liest
   die `queued`-Zeilen **seiner Quelle**. Eine `queued`-Zeile überlebt einen
   Neustart und wird ausgeführt; eine `interrupted`-Zeile wird **nicht**
   aufgenommen (`ADR-0111`: kein automatischer Neustart).
3. **Das Wecksignal** ist ein nicht blockierendes Signal (Kapazität 1,
   Signale verschmelzen), das die Administrations-Goroutine nach einem
   erfolgreichen `Admit` sendet; die Kette beginnt beim `pg_notify` von
   `cdc.backfill_table`, das die Administrations-Goroutine weckt. Die
   Worker-Schleife arbeitet **zuerst** alle `queued`-Zeilen ab und wartet
   **dann**: ein Signal, das während eines Runs eintrifft, geht nicht
   verloren, weil die Schleife nach jedem Run erneut liest. Ein zweiter
   `LISTEN`-Kanal und ein Fallback-Poll sind nicht nötig, weil nur die
   Administrations-Goroutine dieser Instanz `queued`-Zeilen anlegt und der
   Start-Lesegang jedes Absturz-Fenster deckt.
4. **Vorbedingungen vor der Ausführung erneut prüfen:** vor dem Öffnen des
   Slots prüft der Worker Bindung und Publication-Mitgliedschaft
   (`TableActivationPort.Registered`/`Published`); Abweichung endet den Run
   als `failed` (Fehlerklasse `configuration`, kein Slot, keine Kopie). Die
   Zeile wartet jetzt hinter einem anderen Run oder über einen Neustart;
   die Prüfung spart die Kopie. Die Sicherheits-Zusage trägt weiterhin die
   Fail-closed-Prüfung vor dem Commit (`ADR-0111` Teilfrage 4).

### Festlegung 3 — Warn-Kriterium für große Tabellen (Startpunkt)

**Status dieser Festlegung: Startpunkt, Setzung ohne Messung.** Der
Startwert unten ist **nicht gemessen**; Nachschärfen ist erwünscht und kein
Bruch dieser ADR.

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (nur Handbuch-Hinweis) | kein Aufwand | `ADR-0111` Festlegung 3 verlangt die Warnung; der Bench-Slice hat kein Kriterium, aus dem er eine Richtgröße ableiten könnte |
| B — Ablehnung ab einer Zeilenzahl | einfach durchsetzbar | `ADR-0111` Festlegung 3: „keine Ablehnung"; die Grundlage ist eine Schätzung (`pg_class.reltuples`) — eine Schätzung würde zur Grenze (`BEO-PGC/geschaetzter-wert-als-grenze`) |
| C — Zeilenzahl-Schwelle sofort im Code | Warnung ab dem ersten Tag | eine Setzung ohne Messung, die an jedem Träger zur Grenze wird (dieselbe Beobachtung) |
| D — Toleranz und Richtgröße als Konfigurationsschlüssel | ohne Neubau anpassbar | neue Betreiber-Oberfläche (Konfigurationsdatei, Handbuch, Tests) für einen Wert ohne Messung; eine Konstante ist billiger nachzuschärfen |
| E — Laufzeit-Warnung berechnen in der View | keine Ablage nötig | die Toleranz stünde ein zweites Mal in SQL (View) neben dem Code |
| **F — Toleranz als Konstante, zwei Warnungen, Richtgröße aus der Messung (gewählt)** | eine Konstante an einer Stelle; die Laufzeit-Warnung braucht keine Schätzung und fängt, was die Schätzung nicht sieht (Zeilenbreite, anderer Host) | die Richtgröße in Zeilen ist eine Orientierung für den gemessenen Host, keine Aussage über die Tragfähigkeit |

**Festlegung.**

1. **Toleranz:** Startwert **10 Minuten Kopierdauer**, eine benannte
   Konstante im Code an genau **einer** Stelle (Use-Case-Paket des Runs);
   die Uhr ist der `ClockPort` (`ADR-0040`), also mit Fake-Uhr prüfbar.
   **Kein neuer Konfigurationsschlüssel im Erstumfang.** Die Kopierdauer
   beginnt mit dem Übergang nach `running` (`started_at`); die Wartezeit in
   `queued` zählt nicht. Der **geltende** Wert steht an der Konstante und im
   Handbuch, nicht in dieser ADR: sie nennt nur den Startwert. Die Toleranz
   ist kein Eintrag in `spec/pflichtenheft.md` §3: sie steuert ausschließlich
   eine Warnung, die kein Gate und keine Anforderung liest; ein §3-Eintrag
   („über ADR schärfbar") würde jedes Nachschärfen ADR-pflichtig machen und
   dem Startpunkt-Charakter widersprechen.
2. **Richtgröße in Zeilen:** entsteht **erst aus der Messung** des
   Bench-Slice und steht vorher **nirgends** — nicht im Code, nicht im
   Handbuch, nicht in einem Test. Sie ist die Zeilenzahl, die der Bench-Lauf
   in der Toleranz kopiert, abgerundet auf eine runde Zahl (die
   Rundungsregel nennt der Bench-Bericht). Kopiert der Bench nicht die volle
   Toleranzdauer, ist die Zahl aus der gemessenen Rate hochgerechnet und als
   **abgeleitet** gekennzeichnet (`AGENTS.md` §3.12); sie nennt den
   gemessenen Host und den Lauf. Die Kopierdauer, die der Bench misst, ist
   dieselbe wie die der Toleranz: `finished_at − started_at` des Runs.
3. **Zwei Warnungen — keine Ablehnung, kein Abbruch, keine
   Statusänderung:**
   - **(1) beim Antrag:** die **geschätzte** Zeilenzahl liegt über der
     Richtgröße. Wirksam erst, wenn die Richtgröße aus der Messung vorliegt.
   - **(2) zur Laufzeit:** ein Run läuft länger als die Toleranz. Wirksam ab
     dem Erstumfang, unabhängig von einer Schätzung; geprüft bei jedem
     Fortschritts-Update (je Block) und beim Abschluss.
4. **Ablage und Sichtbarkeit:** die Auswertung liegt an **einer** Stelle
   (im Use Case des Runs); ihr Ergebnis wird als Spalte(n) an der Run-Zeile
   festgehalten — Warnung (1) schreibt `Admit`, Warnung (2) der Worker per
   `UPDATE` —, `cdc.backfill_status` reicht sie durch und `diagnose` liest
   sie aus der View. Toleranz und Richtgröße stehen dadurch nirgends ein
   zweites Mal (nicht in der View, nicht in `diagnose`). Ob eine oder zwei
   Spalten und ihre Bezeichner legen `slice-backfill-run-store` und der
   Spec-Nachzug fest.
5. **Unbekannte Schätzung:** `estimated_rows` ist dann `NULL`; Status und
   `diagnose` zeigen „unbekannt", **nie** `0`; Warnung (1) entfällt, Warnung
   (2) greift. Für nie analysierte Tabellen liefert `pg_class.reltuples` in
   neueren PostgreSQL-Versionen `−1` — Wissen aus der
   PostgreSQL-Dokumentation, **in diesem Repo nicht gemessen**; **erwartet**,
   für PostgreSQL 17 und 18 zu belegen im Snapshot-Leser-Slice an der realen
   Datenbank.
6. **Benennung (gegen `BEO-PGC/geschaetzter-wert-als-grenze`):** jede Stelle,
   die die Zeilenzahl nennt, trägt das Wort „geschätzt"; die Richtgröße heißt
   „Richtgröße" oder „Orientierung", nie „Grenze", „Limit" oder „maximal";
   die Toleranz heißt „Startwert" mit dem Zusatz „Setzung ohne Messung", bis
   eine Messung sie ersetzt.

## Konsequenzen

- Positiv: die Annahme ist atomar — ein angenommener Antrag hat einen Run,
  ein Absturz hinterlässt nichts Halbes; Vorbedingungs-Fehler hinterlassen
  keine Zeile.
- Positiv: die Datenpfad-Rolle `cdc_capture` schreibt keine Verwaltungszeile
  an; die Rechte pro Rolle sind je eine Aussage (`LH-QA-SEC-001`/`-002`).
- Positiv: die Warnungen haben einen Wert und einen Ort; keine Zahl steht
  vor ihrer Messung.
- Negativ: ein neuer Port und eine Transaktion im Administrations-Adapter
  (Aufwand im Run-Use-Case- und Run-Store-Slice).
- Negativ: die Richtgröße in Zeilen ignoriert Zeilenbreite und Host; die
  Laufzeit-Warnung (2) ist die Korrektur, kein Ersatz.
- **Akzeptierte Negative** (kurz begründet, keine Folgepflicht):
  - *Fortschritts-Takt der Warnung (2):* sie wird je Block und beim Abschluss
    geprüft; ein einzelner Block, der länger als die Toleranz braucht, warnt
    erst danach — die Warnung ist eine Orientierung, kein Alarm.
  - *Run-Zeilen wachsen:* niemand trägt `DELETE`; je Antrag entsteht eine
    kleine Zeile, `cdc.backfill_status` liest nur den letzten Run je Tabelle.
  - *`INSERT` auf Tabellenebene:* siehe Festlegung 1.

### Folgepflichten

Jede Pflicht hat einen Träger; **die genannten Plan-Dateien ändert diese ADR
nicht** — der Planner zieht sie nach (`AGENTS.md` §3.13).

1. **`slice-backfill-spec-nachzug`:** die neue `SPEC-*`-Kennung trägt die
   Grants dieser ADR, die Warn-Spalte(n), `estimated_rows` = `NULL` als
   „unbekannt"; [`SPEC-019`](../../../spec/pflichtenheft.md) sagt, dass
   `applied` bei `backfill` „angenommen" heißt; die Architektur-Sicht
   beschreibt die Annahme-Sequenz **ohne** ADR- und Slice-Bezug
   (`AGENTS.md` §3.4).
2. **`slice-backfill-run-usecase`:** der Start-Trigger „Schreib-Rolle der
   `queued`-Zeile festgelegt" ist mit dem Status `Accepted` dieser ADR
   erfüllt; §3 (Port-Schnitt) und §6 (erstes Risiko) folgen dem Port-Schnitt
   aus Festlegung 1 — der Run-Zustands-Port verliert „anlegen", der neue
   Annahme-Port kommt hinzu; `Request` endet mit `Admit`; `Execute` prüft die
   Vorbedingungen erneut (Festlegung 2, Punkt 4).
3. **`slice-backfill-run-store`:** die Grants (`cdc_admin` `SELECT`/`INSERT`,
   `cdc_capture` `SELECT`/`UPDATE`) statt „`SELECT`, `INSERT`, `UPDATE` für
   `cdc_capture`"; `estimated_rows` nullable; die Warn-Spalte(n); der
   Annahme-Adapter samt Rollback-Test; die Rollen-Tests in
   `internal/bootstrap/roles_rollout_file_internal_test.go` und im
   Store-Tier (Fitness Function).
4. **`slice-backfill-sql-administration`:** der Verarbeitungs-Zweig ruft
   `Request`; die Worker-Schleife mit Start-Aufnahme und Signal
   (Festlegung 2) samt Start-Reihenfolge; View und `diagnose` zeigen die
   Warn-Spalte(n) (leer, solange die Auswertung nicht liefert) — die dortige
   Abgrenzung „ohne Schwellenwert" bleibt wahr, weil die Konstanten erst
   im Bench-Slice hinzukommen.
5. **`slice-backfill-bench-richtgroesse`:** der Start-Trigger „Warn-Kriterium
   festgelegt" ist mit dem Status `Accepted` dieser ADR erfüllt; die Messung
   ist `finished_at − started_at` (statt „Antrag bis `completed`"); die
   Auswertung beider Warnungen liegt an **einer** Stelle im Use Case des
   Runs, nicht in `bootstrap.Diagnose`, und die Toleranz-Konstante kommt mit
   diesem Slice; die dortige Abgrenzung „keine Warnung am Status-View"
   entfällt, weil die View das Ergebnis trägt (nicht die Konstante); das
   Handbuch nennt die Toleranz als „Startwert, Setzung ohne Messung" und die
   Richtgröße mit ihrem Ursprung (gemessen oder abgeleitet), Host und Lauf.
6. **`slice-backfill-snapshot-reader`:** die Erwartung „`reltuples` = `−1`
   für nie analysierte Tabellen" ist seine Messung an PostgreSQL 17 und 18;
   sein DoD trägt sie bereits.
7. **`slice-backfill-e2e`:** Empfehlung an den Planner — der Negativ-Beleg
   um „eine `queued`-Zeile überlebt einen Neustart und wird ausgeführt"
   erweitern (Festlegung 2).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (`roles_test.go`-Muster) | `cdc_capture`-Login: `INSERT` auf `cdc.backfill_run` scheitert (SQLSTATE `42501`), `UPDATE` gelingt; `cdc_admin`-Login: `INSERT` gelingt, `UPDATE` und `DELETE` scheitern; `cdc_reader`: `SELECT` auf die Basistabelle scheitert, auf `cdc.backfill_status` gelingt | `make test-store` |
| Go-Test, reale PostgreSQL | Annahme atomar: trifft der Antragsvermerk keine `pending`-Zeile, bleibt weder Run-Zeile noch Vermerk; ein zweiter Antrag bei aktivem Run endet ohne zweite Zeile | `make test-store` |
| Go-Test (Fakes), Whitebox `internal/bootstrap` | Worker: `queued` beim Start wird in `(requested_at, run_id)`-Reihenfolge ausgeführt, `interrupted` nicht; ein Signal während eines Runs geht nicht verloren; fehlende Bindung → `failed` (`configuration`) ohne Slot | `make test` |
| Go-Test (Fake-Uhr) | Warnung (2) genau an der Toleranz (darunter keine, darüber gesetzt), je mit einer Mutation; unbekannte Schätzung → keine Warnung (1) und nie `0`; eine Warnung ändert weder Status noch Ablauf | `make test` |
| `a-check` | der Annahme-Port liegt in `ports`, sein Adapter in `adapters`, kein Adapter importiert einen anderen | `make a-check` |
| Review-Prüfpflicht | die Toleranz steht an genau einer Stelle; keine Zeilenzahl-Konstante vor der Messung; „geschätzt" an jeder Nennung der Zeilenzahl | — (kein Gate) |

## Re-Evaluierungs-Trigger

- **Eine Messung liegt vor** (Bench-Slice) **oder ein Betrieb zeigt eine
  andere Toleranz:** die Konstante wird nachgeschärft, das Handbuch zieht
  mit — kein Bruch, keine Folge-ADR.
- **Toleranz oder Richtgröße sollen konfigurierbar sein**, oder die Warnung
  wird Vertrag (Anforderung, Abnahme, `spec/pflichtenheft.md` §3): Folge-ADR
  (Konfigurationsschlüssel bzw. `SPEC-*`-Konstante).
- **Die Warnung soll ablehnen oder abbrechen:** Folge-ADR (`ADR-0111`
  Festlegung 3: „keine Ablehnung").
- **Mehr als eine Instanz je Quelle nimmt Anträge an, oder ein zweiter Worker
  kommt** (Parallelisierung, `ADR-0111`-Ausbaustufe): Aufnahme und
  Annahme-Prüfung neu entscheiden (DB-seitige Beanspruchung, Kanal).
- Sonst permanent — Rollenschnitt und Aufnahme-Regel gelten, solange
  `cdc.backfill_run` von den beiden Rollen geschrieben wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-23 | Proposed — Architect-Vorschlag; die Auftraggeber-Vorentscheidungen zu Rollenschnitt, Aufnahme und Warn-Kriterium sind an den Beleg-Ankern im Code geprüft | [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0113` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
