# Slice routing-betriebsdoku: Benutzerhandbuch — Routing-Regeln konfigurieren, Ziele lesen, Konsequenzen und Abhilfe

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (SQL-Administration),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Zugriffswege),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 8 (Handbuch-Hälfte — Träger dieses Slice) und Entscheidung 4
(Handbuch-Hinweis zu DELETE und Inhaltsregeln).

**Berührte Spec-Stellen:** — (das Handbuch beschreibt Betreiber-Oberfläche und
verweist für Zusagen auf die Spec; es ist selbst keine Spec-Stelle).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `docs/user/benutzerhandbuch.md` beschreibt die Routing-Konfiguration für
Betreiber und Integratoren — nur das, was `slice-routing-e2e` am laufenden System
belegt hat. Drei Liefer-Punkte:

- (A) **Abschnitt „Routing-Regel konfigurieren"** (Formvorbild: der Abschnitt
  „Transformationsregel konfigurieren"): Voraussetzung (`cdc_admin`), Vorgehen mit
  `cdc.set_route`/`cdc.remove_route`, die Form der `rule_spec`, ein lauffähiges
  Beispiel (unter der genannten Rolle ausgeführt), R1–R6 mit ihrer Wirkung (führende
  Stelle für die Fehlertexte bleibt die Spec), die **Konsequenzen** (das Label ist zum
  Erfassungszeitpunkt fest, keine Rückwirkung; ein Change ohne Treffer hat
  `route_target IS NULL` und erscheint nur ungefiltert und über diese Abfrage; der
  Altbestand wird über einen neuen Backfill-Run mit dem aktuellen Regelstand
  gekennzeichnet; eine Inhaltsregel auf eine Nicht-Schlüsselspalte wirkt bei DELETE ohne
  volle Replica-Identität nicht — mit der Messung des E2E-Slice an PostgreSQL 17 und 18),
  die Abhilfe (Regel entfernen, Neustart — soweit `slice-routing-e2e` sie belegt hat)
  und der Hinweis, dass das Ziel Auswahl und kein Zugriffsschutz ist;
- (B) **Lesewege:** der Parameter `target` je Weg (`GET /changes`, gRPC-Stream, SSE,
  RPC `ReadChanges` nach Verdikt V1), die Konjunktion mit `schema`/`table`, das
  Subjekt `cdc.route.<source_id>.<ziel>` mit fire-and-forget und der Lücke-schließen-
  Hinweis über `cdc.changes`, die Spalte `route_target` in „Änderungen lesen"; die
  Kosten-Aussage der zweiten Veröffentlichung trägt ihren Ursprung aus
  `slice-routing-nats-subjekt`;
- (C) **Fehlerklassen, Rollen, Grenzwerte, Historie:** die Zeile `schema` der
  Fehlerklassen-Tabelle und der Abschnitt „Neustart nach einem Fehler" nennen die
  nicht anwendbare Routing-Regel; die Rollen-Tabelle nennt die zwei Funktionen; die
  Grenzwerte nennen Alphabet und Länge des Zielnamens; Version und
  Änderungshistorie tragen die Zeile (Version 1.83 am Parent, gemessen).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Aussagen, die kein Beleg trägt** — das Handbuch nennt nur gemessene Wirkung; eine
  Aussage, die `slice-routing-e2e` nicht belegt hat, steht als *erwartet* oder gar
  nicht ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- **SDK- und Beispiel-Clients** — `slice-routing-sdk-beispiel-target`: das Handbuch
  nennt den Parameter am Roh-Weg; die Abschnitte der drei SDKs und die Flags der
  Beispiele ziehen dort nach (§3, Suchlauf, benannte Zwischenzeit).
- **Spezifikation und Code** — keine Änderung; eine im Schreiben gefundene Lücke ist
  ein Befund, kein stiller Nachzug.
- **Version und Release** — `docs/user/version.md` bleibt; die Handbuch-Version
  (Kopfzeile) und die Änderungshistorie ändern sich.

## 2. Definition of Done

**Übergabe-Block: aufgeschobene Gegenstände aus den Slices der Welle**
([welle-routing](welle-routing.md)) — diese Gegenstände haben die Slices mit
Adresse hierher aufgeschoben; sie gehören zu (A), (B) oder (C):

- aus `slice-routing-kern-label` — die Spalte `route_target` in der View `cdc.changes`
  und ihre Lesesemantik im Abschnitt „Änderungen lesen" und in der Beschreibung der
  Zugriffswege; `NULL` heißt „nicht geroutet", nicht „wal";
- aus `slice-routing-antragsweg` — `cdc.set_route`/`cdc.remove_route` mit Parametern
  und Rolle `cdc_admin` (Rollen-Tabelle, Abschnitt „Zugriff und Rollen"), die
  `rule_spec`-Form, R1–R6 mit Fehlertexten (führende Stelle: die Spec), Abhilfe, und der
  Fall R4: eine Abschlussregel ohne `when` muss die höchste `order` tragen — wer eine
  Regel mit höherer `order` ergänzen will, entfernt die Abschlussregel und setzt sie neu;
  die `order`-Obergrenze 2147483647 (`MaxRouteOrder`, Setzung des Slice
  `routing-antragsweg`, `SPEC-032` nennt keine) und die Annahme jeder JSON-Zahl mit dem
  Wert einer positiven ganzen Zahl (`10`, `10.0`, `1e1`), während ein Bruchteil, 0, ein
  negativer Wert und ein Wert über der Obergrenze `rule_spec ist ungültig` enden;
- aus `slice-routing-backfill-pfad` — Backfill-Bestand trägt das Label des Regelstands
  zum Run; Neuerzeugung des Altbestands über einen neuen Run (Abschnitt „Bestand als
  Backfill überführen"). Übergabe-Block (Umsetzungsstand des Slice, Tests
  `internal/application/usecase/backfill/routing_test.go`): ein `set_route` oder
  `remove_route`, das während des Runs `applied` wird (zwischen zwei Blöcken oder vor
  dem Commit), endet den Run `failed` mit Klasse `configuration` und dem Text
  „Routing-Regelstand während des Backfills geändert“, vor jeder Sichtbarkeit; eine
  Routing-Regel, deren `when.column` der Snapshot nicht trägt, endet ihn `failed` mit
  Klasse `schema` vor der ersten Zeile (Text nennt Regelname und Spalte), Abhilfe
  `cdc.remove_route` und neuer Antrag. Handbuch-Stellen, die bisher nur Ausschluss und
  Transformationsregeln nennen: `benutzerhandbuch.md` Aufzählung „Ausgeschlossene
  Spalten“ im Abschnitt „Bestand als Backfill überführen“ (Wechsel des Standes während
  des Runs) und die Zeile `schema` der Fehlerklassen-Tabelle in §6 „Fehlerbehebung“ (Run mit nicht
  anwendbarer Routing-Regel); Messung der Stellen:
  `git grep -n -E "Ausgeschlossene Spalten|Transformationsregel ist auf die Relation" -- docs/user/benutzerhandbuch.md`;
- aus `slice-routing-lesewege` — Parameter `target` je Weg, Konjunktion mit
  `schema`/`table`, das Verhalten eines Servers ohne diesen Parameter (Alt-Server
  ignoriert das Feld am gRPC-Weg, lehnt den Parameter an HTTP und SSE ab; die Aussage
  trägt ihren Ursprung aus dem Slice), und: das Ziel ist Auswahl, kein Zugriffsschutz;
  Zählwort „zwei optionale Felder/Parameter“ (`schema`/`table`) wird „drei“ (mit
  `target`); Handbuch-Zeilen (Stand `docs/user/benutzerhandbuch.md` am Ende von
  `slice-routing-lesewege`): `:1301` (gRPC-Stream, „zwei optionale“ Felder), `:1364`–`:1365`
  (C#-Stream, „zwei optionale Parameter“), `:1378` (Kotlin-Stream), `:1403` (Python-Stream),
  `:1456`–`:1457` (`ReadChanges` ohne `target` in der Filter-Aufzählung), `:1557`
  (SSE, „Zwei optionale“ Query-Parameter) und `:1149` (Abschnitt „Zugriff über die
  HTTP-/JSON-API“, „optional gefiltert über `schema` und `table`“); die Aussage „ein
  `target` außerhalb des Alphabets liefert leer, keinen Fehler“ steht im Handbuch mit dem
  Qualifier „bei sonst gültiger Anfrage“ (Fehler des Lese-Kontrakts — leere Quelle,
  `limit < 1`, invertierter Bereich, Start-/End-Position einer anderen Quelle — haben
  Vorrang; Norm: `SPEC-022`, `SPEC-031`);
- aus `slice-routing-nats-subjekt` — das Zusatz-Subjekt `cdc.route.<source_id>.<ziel>`,
  fire-and-forget und die Kosten-Aussage mit Ursprung (Abschnitt „Zugriff über den
  NATS-Vollinhalts-Stream"). Übergabe-Block (Umsetzungsstand des Slice, Tests
  `internal/adapters/driven/natsstream/publisher_nats_test.go`): ein Zielname mit `-`
  und `_` (`eu-west_1`) ist am realen NATS-Server ein einzelnes Subjekt-Token (erprobt,
  `make test-notify`), `cdc.route.<source_id>.*` und `.>` empfangen jedes Ziel; die
  Wecksignal-Wurzel `cdc.changes` bekommt von der zweiten Veröffentlichung nichts (nur
  negativ belegt: der Test läuft ohne Wecksignal-Sender, er zeigt null Nachrichten; das
  Wecksignal selbst trägt der unberührte `natsnotify`-Adapter mit seinem Test im selben
  Lauf). Die zweite Veröffentlichung ist von der ersten in der Prüfung unabhängig
  (Tabellen-Subjekt zuerst, danach das Ziel-Subjekt): ein Fehlschlag oder Überspringen der
  Ziel-Veröffentlichung verändert die erste nicht, ein übersprungenes Tabellen-Subjekt
  (leerer oder reservierter Name) lässt das Ziel-Subjekt bestehen (Adapter-Test).
  Kosten-Aussage, **gemessen** in einem Lauf von `make test-notify` (Zeit für
  `publish` samt Flush zum Server, je 10 000 Changes, Median von fünf Läufen, ohne
  Abonnent, Testcontainer-NATS), zwei Läufe des Implementers: „ohne Ziel 18.242683ms
  (548165 Changes/s), mit Ziel 21.239865ms (470813 Changes/s)“ und „ohne Ziel
  17.479304ms (572105 Changes/s), mit Ziel 22.739539ms (439763 Changes/s)“; Verhältnis
  mit/ohne 1,16 und 1,30 (**abgeleitet**). Drei weitere Läufe (Review: 1,05; Fixrunde:
  1,14, gedruckte Zeile „ohne Ziel 18.450212ms (541999 Changes/s), mit Ziel 21.111558ms
  (473674 Changes/s)“; Verifier: 0,92, gedruckte Zeile „ohne Ziel 20,856 ms (479477
  Changes/s), mit Ziel 19,272 ms (518898 Changes/s)“) ergeben über fünf Läufe eine
  Streuung von 0,92 bis 1,30; das Review-Verhältnis ist **übernommen** aus dem
  Review-Report, die Verifier-Zeile stammt aus dem Verifikations-Report (ebenfalls
  übernommen), das Verhältnis 0,92 ist **abgeleitet**. Ein Aufschlag der zweiten
  Veröffentlichung ist an diesem Messaufbau nicht auflösbar (das Verhältnis liegt im
  Rauschen). Die
  Messung hat keine Schwelle und deckt die Verteilung an Abonnenten nicht ab; der
  Re-Evaluierungs-Trigger der ADR bleibt „Messung der zweiten NATS-Veröffentlichung
  zeigt Druck am Publisher“.
- aus `slice-routing-e2e` — Übergabe-Block, **gemessen** in `make test-integration`
  (Läufe vom 2026-10-01 am Arbeitsbaum von `slice-routing-e2e` nach der Fixrunde, je
  einmal mit dem Image des 18-Legs aus `.github/workflows/e2e.yml`
  (`postgres:18-alpine@sha256:63bdc97d67b5133bf0e5ebd500bec6d046fa851dc81340d838f0347e616107e8`
  = PostgreSQL 18.6) und des 17-Legs
  (`postgres:17-alpine@sha256:7456ef82e5f5bc43d997f4781bbd7c0d6389bff397564649a356e206ba473aee`
  = PostgreSQL 17.11),
  die Versionen stehen in der gedruckten Zeile `ROUTING-DELETE-MESSUNG` der Testfunktion
  `TestE2ERoutingDeleteWithoutFullReplicaIdentity`):
  DELETE einer Zeile ohne volle Replica-Identität trägt im Alt-Bild nur den Schlüssel
  (`{"id": "1"}`); eine Inhaltsregel auf eine Nicht-Schlüsselspalte trifft dort nicht —
  mit einer Abschlussregel ohne `when` bestimmt diese das Ziel (gedruckt `DELETE="sonstige"`),
  ohne sie bleibt `route_target` NULL (`DELETE=""`, nicht gesetzt); unter
  `REPLICA IDENTITY FULL` trifft die Inhaltsregel (`DELETE="eu"`, Alt-Bild mit allen
  Spalten); INSERT und UPDATE derselben Tabellen tragen das Ziel der Inhaltsregel (`eu`).
  Beide Versionen ergaben dieselben Ziele; die Erwartung von `ADR-0137` Entscheidung 4 ist
  bestätigt, es gibt keinen Befund. Dieselbe Phase `Routing-Nichtanwendbarkeit und
  Abhilfe` (gleicher Lauf, beide Versionen) misst die Erreichbarkeit nach
  [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md):
  die nicht anwendbare Regel entsteht am System in (b) Erstaktivierung ohne Spaltenform
  (null Zeilen in `cdc.table_schema` nach der Aktivierung; Regel auf `region`, Spalte
  entfernt, erste Change der Tabelle) und in (c) Publication mit Spaltenliste `id,name`;
  der Erfassungspfad endet mit der Klasse `schema` und dem Log-Text „Routing-Regel auf
  die Änderung nicht anwendbar“, keine Change der Tabelle ist persistiert; die Abhilfe
  (`cdc.remove_route` bei stehendem Prozess, Antrag `pending`; Neustart; Antrag `applied`)
  legt die zuvor nicht bestätigte Zeile ohne Ziel und ohne zweiten `schema`-Fehler vor.
  In (a) — Spalte entfernt, nachdem eine Change der Tabelle die Spaltenform angelegt hat —
  endet der Pfad als inkompatible Schemaänderung (Log-Text „Relation-Änderung nicht sicher
  als Obermenge interpretierbar“), die Change nach der Entfernung hat keine Zeile in
  `cdc.changes`; dass das Entfernen der Regel dort nicht genügt, ist **hergeleitet**
  (`ADR-0140` Entscheidung 2), der Neustart nach (a) ist nicht gefahren. Replay und
  Backfill (gemessen): dieselbe `change_id` liefert nach einer Regeländerung dasselbe Ziel
  über `cdc.changes` und `GET /changes?target=`, eine vor der Regel erfasste Change bleibt
  ohne Ziel, ein Backfill-Run trägt das Ziel des Regelstands zum Run (`origin = 'backfill'`).
  Messbare Stelle im Test: `git grep -n "ROUTING-" -- test/integration/routing_e2e_test.go`.

- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (A): der Abschnitt „Routing-Regel
      konfigurieren" steht im Abschnitt „Aufgaben"; jedes SQL-Beispiel ist unter der
      genannten Rolle (`cdc_admin`) gegen die Compose-Umgebung ausgeführt worden, die
      gedruckte Antwort steht im Bericht
      (`BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig`); R1–R6, die
      Konsequenzen, die Abhilfe und der DELETE-Hinweis stehen mit ihrem Ursprung (gemessen
      von `slice-routing-e2e`, übernommen aus der Spec, oder als erwartet gekennzeichnet);
      der DELETE-Hinweis nennt die Menge der Messung (PostgreSQL-Versionen und Fälle).
      *Zu belegen durch:* Lesen des Abschnitts gegen den Bericht von `slice-routing-e2e`,
      Ausführen der Beispiele, `make docs-check`. *Beleg:* Beispiel, `GET /changes?target=`,
      R-Tabelle (neun von zehn Zeilen), Formprüfungen, `order`-Grenze, DELETE und Backfill
      vom Verifier am System gefahren, gedruckte Ausgabe im
      [Verifikations-Report](../../../reviews/verifikation-slice-routing-betriebsdoku.md) §2
      (Implementer und Reviewer fuhren das Beispiel unter `cdc_admin`/`cdc_reader` ebenfalls);
      die R4-Zeile „höchste `order`“ ist nicht gefahren und im Handbuch als Text der Spec
      gekennzeichnet; Ursprung je Aussage im Handbuch (übernommen aus dem E2E-Lauf mit
      Verweis auf dessen Verifikations-Report; die Phase „Routing-Nichtanwendbarkeit“ ist
      dort nur an PostgreSQL 17.11 vom Verifier gedruckt, 18.6 übernommen).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (B): die Abschnitte zu den
      Zugriffswegen nennen `target` und das NATS-Subjekt, „Änderungen lesen" die Spalte
      `route_target`; die Aussage zu Altservern und die Kosten-Aussage stehen mit ihrem
      Ursprung; was für die SDK-Packages gilt, steht als „folgt mit den SDK-Packages"
      oder gar nicht, bis `slice-routing-sdk-beispiel-target` es einlöst.
      *Zu belegen durch:* Lesen, `make docs-check`, Suchlauf in §3. *Beleg:* Verifikations-Report
      §3 Zeilen 2 bis 4 und §4; die SDK-Abschnitte tragen „folgt mit dem Package“; die Kosten-Spanne
      nennt nur Einzelwerte aus verlinkten Berichten (V-1 nachgezogen).
- [x] (C): die Zeile `schema` der Fehlerklassen-Tabelle und der Abschnitt „Neustart nach
      einem Fehler" nennen die nicht anwendbare Routing-Regel (Pfad nach Messung von
      `slice-routing-e2e`); der Routing-Fehlerblock trennt nach
      [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
      die zwei Ursachen der Klasse `schema` — nicht anwendbare Regel (Abhilfe: Regel
      entfernen, Neustart) und entfernte Spalte (inkompatible Schemaänderung; die Regel
      auf diese Spalte ist zusätzlich zu entfernen) — und sagt für die entfernte Spalte
      nicht „Regel entfernen“ als Abhilfe; die Rollen-Tabelle nennt die zwei Funktionen; „Grenzwerte"
      nennt Alphabet und Länge des Zielnamens (Ursprung: die Spec) und die Zahl der Regeln,
      soweit die Spec eine Grenze sagt; die Handbuch-Version (Kopfzeile) und die
      Änderungshistorie tragen die nächste Zeile (Parent: Version 1.83, gemessen).
      *Zu belegen durch:* `make docs-check`; die Historie-Zeile zeigt die Version der
      Kopfzeile (`BEO-PGC/handbuch-versionshistorie-uebersprungen`). *Beleg:* Verifikations-Report
      §3 Zeilen 5 und 6 (Kopfzeile `Version: 1.84`, Historienzeile 1.84 als letzte Zeile hinter 1.83).
      Eine Anleitung zur Abhilfe der inkompatiblen Schemaänderung (Ursache 2) steht nicht im Handbuch
      und wird nicht versprochen; benannte Hinnahme, Adresse
      `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (§6).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-betriebsdoku`](../../../reviews/review-slice-routing-betriebsdoku.md)
      (3 LOW, 2 INFO, kein HIGH/MEDIUM) und Verifier-Gegenprüfung der Fixrunde
      [`verifikation-slice-routing-betriebsdoku`](../../../reviews/verifikation-slice-routing-betriebsdoku.md)
      (Wahrheitsprobe am System ausgeführt, einschließlich „beide Regeln entfernen, Neustart: Fehler
      bleibt“). Es gab **kein separates Re-Review** der Fixrunde: die Fixrunde änderte nur
      Handbuchtext, und ein anderer Kontext hat ihre neue Tatsachenaussage ausgeführt — engere Fassung
      von `BEO-PGC/fixrunde-ohne-reviewer-lesung`.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-betriebsdoku.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: der Slice **ist** das Doku-Update des Handbuchs; `README.md` im
      Repo-Wurzelverzeichnis nur, soweit der Suchlauf eine bewegte Beschreibung findet
      (der Suchlauf fand keine, siehe Feld in §3).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (§7).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](welle-routing.md) (die Roadmap führt sie
      unter *Abgeschlossene Wellen*, die Closure erfolgte am 2026-10-02). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` §4 „Aufgaben" (neuer Abschnitt „Routing-Regel konfigurieren", Formvorbild „Transformationsregel konfigurieren") | neu | Liefer-Punkt A. |
| `docs/user/benutzerhandbuch.md` §4 („Änderungen lesen", „Bestand als Backfill überführen", „Zugriff über …": HTTP, gRPC-Stream, gRPC-Verwaltungs-API, SSE, NATS-Vollinhalt) | update | Liefer-Punkt B. |
| `docs/user/benutzerhandbuch.md` §2 „Zugriff und Rollen", §6 „Fehlerklassen" und „Neustart nach einem Fehler", §9 „Grenzwerte" und „Änderungshistorie", Kopfzeile (Version, Stand) | update | Liefer-Punkt C. |
| `docs/user/benutzerhandbuch.md` Transformationen, „Fehler: Erfassung endet mit der Fehlerklasse `schema`“ (Bestandsbefund, Parent `2629d544` etwa Z. 387; *hergeleitet* nach `ADR-0112` Teilfrage 4) | update | Handbuch-Nachzug zu `ADR-0140`: Ursache „`column` fehlt in der Relation“ mit Lösung „Regel entfernen, Neustart“ ist für die entfernte Spalte falsch (der Pfad endet dort schon an der inkompatiblen Schemaänderung); die Ursachen trennen wie im Routing-Block. Das Handbuch wird in diesem Plan nicht geändert, der Zug gehört diesem Slice. |
| `docs/user/benutzerhandbuch.md` „Änderungen lesen“ und Feldliste (Bestandsbefund Review F-4, Parent `2629d544`: `:516` „`origin` ist die letzte Spalte der View“, `:1462` Feldliste endet bei `origin`) | update | `route_target` ist seit der Spalte `route_target` die letzte Spalte der View `cdc.changes`; beide Stellen nachziehen (Liefer-Punkt B). |
| `README.md` | lesen | nur, wenn der Suchlauf eine bewegte Beschreibung findet. |
| `docs/user/benutzerhandbuch.md` §2 „Zugriff und Rollen" (Rollen-Tabelle, Zeile `cdc_admin`, Z. 76 am Parent) | update | `cdc.set_route`/`cdc.remove_route` neben den Transformations-Funktionen; Liefer-Punkt C. |
| `docs/user/benutzerhandbuch.md` §4 „Bestand als Backfill überführen" (Aufzählung „Ausgeschlossene Spalten", Z. 756–759, und „Zustellung", Z. 760–763) | update | Routing-Regelstand im Run: Label des Regelstands zum Run, Wechsel des Standes endet den Run `failed` (`configuration`), nicht anwendbare Regel endet ihn `failed` (`schema`) vor der ersten Zeile; Lesefehler = Klasse der Ursache (`ADR-0141`). Liefer-Punkte A/C. |
| `docs/user/benutzerhandbuch.md` §4 „Zugriff über die HTTP-/JSON-API" (Z. 1148–1172), „… den gRPC-Change-Stream" (Z. 1301–1308), „… die gRPC-Verwaltungs-API" (Z. 1454–1464), „… Server-Sent-Events" (Z. 1557–1566), „… das NATS-Wecksignal" (Z. 1659–1662), „… den NATS-Vollinhalts-Stream" (Z. 1697–1712) | update | Liefer-Punkt B: `target` je Weg, Konjunktion, Qualifier „bei sonst gültiger Anfrage", Zählwort „zwei" → „drei" (Roh-Wege), NATS-Zusatz-Subjekt, Wecksignal trägt kein Ziel, Altserver-Aussage als *hergeleitet*. |
| `docs/user/benutzerhandbuch.md` §4 „Schema aktualisieren" (Z. 1033, 1059), §8 „Glossar" (Z. 1975), §9 „Grenzwerte" | update | zweite zusätzliche View-Spalte `route_target` neben `origin`; Glossar-Einträge Routing-Regel/Zustellziel; Zielname-Alphabet, `order`-Obergrenze. |

**Konkretisierung (Implementer, vor der ersten Änderung, Parent `30fd6cb5`):**

- **Neuer Abschnitt** „Routing-Regel konfigurieren" steht hinter „Transformationsregel konfigurieren" und vor „Aktivierte Tabellen auflisten": Voraussetzung, Vorgehen (`rule_spec`-Tabelle, Aufruf, Status-Poll), ein Beispiel, Regel entfernen, R1–R6 als Klartext-Tabelle (Fehlertexte: führende Stelle `SPEC-019`), Fehlerblock der Klasse `schema` (zwei Ursachen nach `ADR-0140`, dazu der Run), „Ziel lesen" (Tabelle je Weg), Hinweise (Label fest, kein Treffer, abwesender Wert, DELETE, Auswahl statt Zugriffsschutz, Verhältnis zu Transformation und Ausschluss).
- **Beispiel und Fehlertexte** werden in einer Wegwerf-Compose-Umgebung (`compose.yaml` mit Überschreibung im Scratchpad, PostgreSQL 18.6, Image `:dev`) unter Login-Identitäten mit `cdc_admin`- bzw. `cdc_reader`-Mitgliedschaft ausgeführt; die gedruckte Ausgabe steht im Bericht. Was nur der E2E-Slice gemessen hat (Neustart, API-Aktivierung, Mengen der Nichtanwendbarkeit, DELETE an PostgreSQL 17.11/18.6), steht mit Ursprung „übernommen aus dem E2E-Lauf von `make test-integration`" samt Verweis auf den Verifikations-Report des Quell-Slice (Fixrunde nach dem Review; im Handbuch nie als dort selbst gemessen); was nur Unit-Tests tragen (Backfill-Run mit nicht anwendbarer Regel, Wechsel des Standes), steht als „durch Unit-Tests belegt, am System nicht gefahren" ohne Datei- oder Testnamen (die Namen stehen in den Berichten der Quell-Slices).
- **Zählwort an den SDK-Stellen (Abweichung vom Wortlaut des Übergabe-Blocks, begründet):** die drei SDK-Packages tragen am Stream genau zwei optionale Parameter (`git grep -n -i target -- sdks/csharp/PgChangeFeed.Client sdks/python/pgchangefeed/pgchangefeed sdks/kotlin/pgchangefeed-kotlin/src/main` findet nur Treffer in Build-Dateien und einem Kommentar, keinen Parameter `target`; gemessen am Parent). „zwei optionale Parameter `schema`/`table`" bleibt an `:1364`–`:1365`, `:1378` und `:1403` deshalb **wahr** und steht weiter; der Satz „tragen denselben Filter wie oben beschrieben" wird auf `schema`/`table` eingegrenzt und um „`target` folgt mit den SDK-Packages" ergänzt. Das Zählwort wird an den Roh-Wegen (`:1149`, `:1301`/`:1305`, `:1456`–`:1457`, `:1557`) auf drei gezogen.
- **Transformations-Fehlerblock** (Z. 387–399): die Ursache „`column` fehlt in der Relation" wird nach `ADR-0140` getrennt (entfernte Spalte bei bekannter Spaltenform endet als inkompatible Schemaänderung; Regel entfernen genügt dort nicht — *hergeleitet*). Fixrunde nach dem Review: Der Satz „Spaltenform noch nicht bekannt“ ist gestrichen (`ADR-0112` Teilfrage 4 nennt ihn nicht; die Transformations-Phase des E2E-Laufs belegt nur die Namenskollision) — die Ursache nennt die fehlende `column` in der Relation der Change. Im Routing-Block gehört eine Publication mit Spaltenliste an einer Tabelle **mit** bekannter Spaltenform zur Ursache 2 (am System gefahren: Exit 1, Klasse `schema`, „Relation-Änderung nicht sicher als Obermenge interpretierbar“; `ADR-0140` Option 3 deckt den Fall, keine neue ADR), zur Ursache 1 nur die Tabelle ohne bekannte Spaltenform.
- **Nicht in diesem Zug:** die Abschnitte der drei SDK-Packages (Methoden, Parameter) und die Flags der Beispiel-Clients (`slice-routing-sdk-beispiel-target`); `examples/README.md` und SDK-READMEs bleiben unberührt.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „das Handbuch kennt zwei
Transformations-Funktionen und keine Routing-Funktion"; „Filter über `schema`/`table`
sind die einzigen Auswahl-Parameter"; „Version 1.83"; Parent ist `30fd6cb5`; der
Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und Nichtgefundenes ein):**

```suchlauf
30fd6cb5 10 -n -E 'cdc\.set_transformation|cdc\.remove_transformation' -- docs/user
30fd6cb5 36 -n -E 'cdc\.changes' -- docs/user
30fd6cb5 1 -n -E '^Version: ' -- docs/user/benutzerhandbuch.md
30fd6cb5 4 -n -E '^\| 1\.8[0-9] ' -- docs/user/benutzerhandbuch.md
30fd6cb5 18 -n -E 'nehmen den Filter|über `-schema`/`-table`|-schema|`schema` und `table`' -- docs/user examples/README.md
30fd6cb5 0 -n -E 'cdc\.set_route|cdc\.remove_route' -- docs/user
diff 18 -n -E 'cdc\.set_route|cdc\.remove_route' -- docs/user
diff 10 -n -E 'cdc\.set_transformation|cdc\.remove_transformation' -- docs/user
30fd6cb5 0 -n -E 'route_target' -- docs/user
diff 29 -n -E 'route_target' -- docs/user
diff 1 -n -E '^Version: 1\.84$' -- docs/user/benutzerhandbuch.md
diff 5 -n -E '^\| 1\.8[0-9] ' -- docs/user/benutzerhandbuch.md
30fd6cb5 4 -n -i -E 'zwei optionale' -- docs/user examples/README.md
diff 2 -n -i -E 'zwei optionale' -- docs/user examples/README.md
diff 2 -n -i -E 'drei optionale' -- docs/user
diff 17 -n -E 'nehmen den Filter|über `-schema`/`-table`|-schema|`schema` und `table`' -- docs/user examples/README.md
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, an denen die Transformations-Funktionen stehen (Muster für die Routing-Funktionen) | Zeile 1: 10 Zeilen in `docs/user` | jede Stelle lesen: Rollen-Tabelle und Aufzählungen, die um die zwei Routing-Funktionen wachsen; **Befund am Diff (gemessen):** Transformations-Funktionen weiter 10 Zeilen (`diff 10`); Routing-Funktionen 0 am Parent, 18 am Diff — Fundstellen der Rollen-Tabelle (Zeile `cdc_admin`), des neuen Abschnitts und des Glossars; Gefunden und nachgezogen: Rollen-Tabelle, Glossar. Nicht gefunden: eine Aufzählung der Antragsfunktionen außerhalb der Rollen-Tabelle, die die Transformations-Funktionen nennt und die Routing-Funktionen nicht (die zehn Parent-Treffer wurden einzeln gelesen). |
| Beschreibungen von `cdc.changes` | Zeile 2: 36 Zeilen | jede Stelle, die die Spalten der View aufzählt, trägt `route_target`; **Befund am Diff (gemessen):** `route_target` 0 Zeilen am Parent, 29 am Diff in `docs/user`. Gefunden und nachgezogen: „Änderungen lesen“ (SELECT-Liste und Satz zur letzten Spalte, Parent `:503`/`:516`), HTTP „Changes lesen“ (Antwort ohne `route_target`, Parent `:1155`–`:1158`), „Schema aktualisieren“ (Parent `:1033`, `:1059`: zwei zusätzliche letzte Spalten), Feldliste der gRPC-Verwaltungs-API (Parent `:1462`: das Ziel gehört nicht zu den dreizehn Feldern). Nicht gefunden: eine weitere vollständige Spaltenliste von `cdc.changes` (die übrigen Parent-Treffer sind Projektionen wie die Fortsetzungs-Abfrage `:533` und die Retention-Abfrage `:813`; sie zählen nicht auf, was die View trägt). |
| Version und Historie | Zeilen 3 und 4: 1 Kopfzeile, 4 Historie-Zeilen der Reihe 1.8x | die nächste Version ist 1.84 (*hergeleitet*; der Implementer misst am Start, falls ein anderer Slice die Version bewegt hat); Kopfzeile und Historie-Zeile tragen dieselbe Version; **Befund am Diff (gemessen):** am Start war 1.83 die Kopfzeile und die Reihe 1.8x hatte 4 Zeilen; am Diff `Version: 1.84` (1 Treffer), 5 Zeilen der Reihe 1.8x, die neue Zeile 1.84 steht hinter der letzten Zeile 1.83 (Tabellenende, letzte Zeile der Datei). |
| Filter-Beschreibung der Clients (`-schema`/`-table`) und der Stream-/Lese-Wege (`schema` und `table`) | Zeile 5: 18 Zeilen (das Muster trifft auch die Schreibweise mit Backticks um beide Wörter, `benutzerhandbuch.md:1149`), darunter `benutzerhandbuch.md:1308` („Go, C# und Kotlin nehmen den Filter über `-schema`/`-table`") | **benannte Zwischenzeit:** diese Aussage zum `target`-Flag der Clients zieht `slice-routing-sdk-beispiel-target` nach; dieser Slice schreibt dazu nichts Behauptendes (§2 B). **Befund am Diff (gemessen):** Zeile 5 jetzt 17 Treffer (die HTTP-Stelle `:1149` nennt `schema`, `table` und `target` und trifft das Muster nicht mehr); „zwei optionale“: 4 am Parent (`:1301` gRPC-Stream, `:1378` Kotlin, `:1403` Python, `:1557` SSE; die C#-Stelle `:1364`–`:1365` bricht zwischen „zwei“ und „optionale“ um und trifft das Zeilenmuster nie), 2 am Diff (Kotlin, Python; plus C#) — diese drei SDK-Stellen bleiben „zwei“, weil die Packages zwei Parameter tragen (§3 Konkretisierung); „drei optionale“: 2 Treffer am Diff (gRPC-Stream „drei optionale, unabhängig setzbare Felder“ und SSE „Drei optionale“). Gefunden und nachgezogen: `:1149` (HTTP), `:1301`/`:1305` (gRPC-Stream, „beide leer“ → „alle drei Felder leer“), `:1456`–`:1457` (`ReadChanges`), `:1557` (SSE). Zwischenzeit: `:1308` („Go, C# und Kotlin nehmen den Filter über `-schema`/`-table`“), die SDK-Abschnitte und die Go-Zeile `:1339` („beide leer“ für die zwei Flags des Beispiels) bleiben und tragen den Hinweis „folgt mit den SDK-Packages“. Nicht gefunden: `README.md` im Wurzelverzeichnis, `examples/README.md` und die SDK-READMEs nennen weder `route_target` noch `cdc.set_route` (`git grep -n -E 'route_target\|cdc\.set_route' -- README.md examples sdks`: 0 Treffer). |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-e2e` liegt in `done/` (das
Handbuch beschreibt nur Gemessenes, Welle §4 Abweichung 1) und kein anderer Slice
liegt in `in-progress/` (WIP-Limit 1). Der Übergabe-Block in §2 ist vollständig, wenn
alle fünf Quell-Slices in `done/` liegen (sie liegen vor `slice-routing-e2e`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Übergabe-Block und die
  drei Liefer-Punkte sprengen den Umfang eines Doku-Zugs nicht erwartet; sprengten
  sie ihn, trennt sich (B) ab (`slice-routing-betriebsdoku-lesewege`).
- `in-progress` → `open` (blockiert): der Bericht von `slice-routing-e2e` trägt eine
  Aussage nicht, die das Handbuch braucht (z. B. die Abhilfe der Negative, V3) — dann
  steht die Aussage als erwartet oder gar nicht, und die Lücke geht als Befund an den
  Architect; das Handbuch wartet nicht auf einen Beleg, den es nicht gibt.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), SQL-Beispiele unter der genannten Rolle ausgeführt, Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Zwei Quellen für dieselben Fakten.** Fehlertexte, `rule_spec`-Form und Alphabet
  stehen in der Spec und im Handbuch; die Spec ist führend
  (`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`, offen, 3×; Ausgang beim
  Lese-Schritt der Welle-Closure). — **Ausgang:** bei der Closure einzutragen (das
  Handbuch verweist je Sachverhalt auf die führende Stelle, statt Texte zu
  wiederholen). — **Ausgang: weiter offen.** Das Handbuch wiederholt die zehn Fehlertexte der
  R-Tabelle mit ihrer Adresse und verweist für Wortlaut und Reihenfolge auf `SPEC-019`; neun der
  zehn Zeilen sind vom Verifier am System gegen den Text gehalten (Verifikations-Report §2), die
  zehnte (R4 „höchste `order`“) ist als nicht gefahren gekennzeichnet. Die Drift-Klasse selbst
  bleibt beim Register (`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`, offen, 3×);
  Adresse: Lese-Schritt der Closure von [welle-routing](welle-routing.md).
- **Beispiele unter der falschen Rolle.** Ein SQL-Beispiel, das unter `cdc_admin`
  läuft, aber ohne Rollenangabe dasteht, scheitert beim Leser
  (`BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig`). — **Ausgang:**
  bei der Closure einzutragen (Beispiele gelaufen, gedruckte Zeile). — **Ausgang: entfallen.**
  Das Beispiel lief unter einer Login-Identität mit `cdc_admin` (Anträge) und einer mit
  `cdc_reader` (Lesen); der Verifier wiederholte es in einer eigenen Wegwerf-Umgebung mit
  gedruckter Tabelle der Zeilen 1 bis 6, identisch mit dem Handbuch (Verifikations-Report §2).
- **Aussage ohne Messung.** Die DELETE-Aussage, die Altserver-Aussage und die
  Kosten-Aussage sind Erwartungen, bis ihr Slice sie belegt hat
  (`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`, verkörpert, 11×). —
  **Ausgang:** bei der Closure einzutragen (jede Aussage mit Ursprung: gemessen,
  übernommen, abgeleitet, erwartet). — **Ausgang: teils entfallen, teils weiter offen, je
  Aussage mit Beleg-Anker.** *Entfallen (gemessen von Reviewer und Verifier selbst):* Beispiel,
  `GET /changes?target=`, neun R-Zeilen, Formprüfungen, `order`-Grenze, DELETE ohne volle
  Replica-Identität und Backfill-Label, Ursache 2 samt „beide Regeln entfernen, Neustart: Fehler
  bleibt“ (Verifikations-Report §2). *Übernommen, im Handbuch so gekennzeichnet:* die
  E2E-Phase „Routing-Nichtanwendbarkeit und Abhilfe“ ist im Verifikations-Report von
  `slice-routing-e2e` nur an PostgreSQL 17.11 vom Verifier gedruckt; PostgreSQL 18.6 ist Lauf des
  Implementers und der `e2e.yml`-Legs (V-3, Handbuch auf diese Aussage berichtigt); gRPC-/SSE-/
  NATS-Zielauswahl aus dem E2E-Verifikations-Report. *Durch Unit-Tests belegt, am System nicht
  gefahren:* die Backfill-Run-Fälle (Wechsel des Standes, nicht anwendbare Regel). *Weiter offen
  (nicht gefahren):* R4 „höchste `order`“ (V-4, Text der Spec; Adresse: ein Lauf, der diese
  Zeile am System erzeugt, wenn ein Folge-Zug das Handbuch berührt). *Hergeleitet, weiter
  hergeleitet:* Altserver-Verhalten der Lesewege, Schlüsselspalten-Bedingung bei `DELETE`,
  „Regel zusätzlich entfernen“ bei Ursache 2, die Abhilfe (a) nach `ADR-0140` Entscheidung 2.
  *Weiter offen (Lücke, benannte Hinnahme, V-2):* die Abhilfe der inkompatiblen Schemaänderung
  ist im Handbuch nicht beschrieben (das Handbuch sagt es, ohne Versprechen); Adresse:
  `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, jetzt 2×).
- **Ursprung nicht auflösbar (V-1, Review F-1).** Die Kosten-Spanne trug ihren Ursprung als
  „gemessen und abgeleitet in diesem Handbuch-Zug“ und „Berichte der Umsetzung“. — **Ausgang:
  eingetreten und behoben:** das Handbuch nennt die Spanne als Spanne mehrerer Läufe ohne
  Schwelle und führt nur die Einzelwerte, die in den verlinkten, committeten Berichten des
  Quell-Slice stehen; alle anderen Einzelzahlen sind entfernt. Register:
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (neue Form, siehe §7).
- **Handbuch nicht nachgezogen / Historie übersprungen**
  (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`, 3×;
  `BEO-PGC/handbuch-versionshistorie-uebersprungen`, 3×; beide verkörpert). —
  **Ausgang: entfallen.** Version 1.84 in Kopfzeile und als letzte Historienzeile hinter
  1.83; der Suchlauf-Block ist nachgemessen (`make suchlauf-nachmessen`, Exit 0) und trägt
  Gefundenes wie Nichtgefundenes je Träger.
- **Zwischenzeit mit den SDK-Packages.** Zwischen diesem Slice und
  `slice-routing-sdk-beispiel-target` beschreibt das Handbuch `target` am Roh-Weg, die
  SDK-Abschnitte nicht. — **Ausgang: weiter offen.** Die drei SDK-Stellen (C#, Kotlin, Python)
  tragen „zwei optionale Parameter `schema`/`table`“ und „das Feld `target` folgt mit dem
  Package“; Nachzug: [slice-routing-sdk-beispiel-target](slice-routing-sdk-beispiel-target.md)
  (Suchlauf-Zeile 5 in §3).

## 7. Closure-Notiz

- **Was hat funktioniert:** Das Handbuch 1.84 beschreibt die Routing-Konfiguration für
  Betreiber: Abschnitt „Routing-Regel konfigurieren“, `target` an den Roh-Wegen, das
  Zusatz-Subjekt, die Spalte `route_target`, die Fehlerklasse `schema` in zwei Ursachen. Die
  Wahrheitsproben sind von Reviewer und Verifier **selbst** gefahren (Wegwerf-Umgebung, Beispiel
  unter `cdc_admin`/`cdc_reader`, neun R-Zeilen, `GET /changes?target=`, Ursache 2 einschließlich
  „beide Regeln entfernen, Neustart: Fehler bleibt“); die Implementer-Beispiele sind in dieser
  Form **übernommen**, der Verifier hat sie wiederholt (Verifikations-Report §2). Der
  Verifier-Lauf ergab: DoD getragen, kein HIGH/MEDIUM.
- **Was ging anders als geplant:** Eine Fixrunde (`e465ca1c`) zu drei LOW und zwei INFO des
  Reviews (Ursache-Zuordnung bei der Publication mit Spaltenliste, Ursprungs-Anker, Adressliste,
  Test- und Dateinamen in Betreiber-Prosa). Ein separates Re-Review gab es nicht: die Fixrunde
  änderte nur Handbuchtext, ein anderer Kontext (der Verifier) hat ihre neue Tatsachenaussage am
  System ausgeführt (engere Fassung von `BEO-PGC/fixrunde-ohne-reviewer-lesung`). Nach der
  Verifikation blieben zwei Formpunkte, in dieser Closure gezogen: **V-1** (Ursprung der
  Kosten-Spanne nicht auflösbar: nur noch Einzelwerte mit verlinktem Bericht, „Handbuch-Zug“
  und „Berichte der Umsetzung“ entfernt), **V-3** (Ursache 1 nannte „17.11 und 18.6“ als
  gemessen; der E2E-Verifikations-Report druckt die Phase nur für 17.11, 18.6 steht als
  Lauf des Implementers und der `e2e.yml`-Legs, übernommen). **V-2** (Abhilfe der inkompatiblen
  Schemaänderung fehlt im Handbuch) ist eine Vorbestands-Lücke und als benannte Hinnahme geführt
  (§6); das Handbuch sagt sie ausdrücklich. **V-4** (R4-Zeile nicht gefahren) steht im Handbuch
  als „Text der Spec, am System nicht gefahren“. V-5 war der Stand der Platzhalter dieser Notiz.
  Die Handbuch-Version bleibt 1.84; die Historienzeile trägt die Nachzüge.
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor (`AGENTS.md` §3.12, Ursprung
  als Anker): **In einem Betreiber-Handbuch ist ein Ursprung nur so viel wert, wie der Leser ihn
  auflösen kann.** „Gemessen in diesem Handbuch-Zug“ oder „aus den Berichten der Umsetzung“
  ist kein Anker, weil der Leser den Zug nicht kennt und die gedruckte Zeile in keinem
  committeten Report steht. Die Handlung des Planners beim Schreiben: pro Zahl den
  **committeten, verlinkten** Report nennen, in dem die gedruckte Zeile steht, oder keine
  Einzelzahl nennen und die Spanne als Spanne führen. Zweitens: Wahrheitsproben (Reviewer,
  Verifier) fuhren die Aussagen selbst, die Implementer-Beispiele sind **übernommen**; ein
  Handbuch darf eine Messung des Quell-Slice nur so weit zusammenfassen („beide Versionen“),
  wie dessen Report sie trägt (V-3). §3.7: keine Chronik-Sprache im Handbuch; „Handbuch-Zug“ und
  „Review, Fixrunde, Verifikation“ waren Prozess-Vokabular und sind entfernt. Träger: die
  Lese-Handlung des Reviewers (HIGH-Punkt „Zahl im Träger ohne Ursprung“) und der Verifier;
  keine neue Regel im Text von `AGENTS.md`, nichts neu verkörpert.
- **Beobachtungs-Register (`../observations/`):** zwei Dateien, ein Vermerk, zwei
  Deckel-Entscheidungen:
  - **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** (verkörpert, Deckel bei 23×):
    `evidence/slice-routing-betriebsdoku.md`, Zähler **28×** → **29×** (`ls evidence | wc -l`).
    Datei trotz Deckel und trotz Schwere LOW, weil es eine **neue Form** ist: der Ursprung ist
    genannt, aber für den Leser nicht auflösbar (V-1, Review F-1); die bisherigen Formen
    sind Zahl ohne Ursprung und Zahl gegen die Messung. Gefunden von Reviewer und Verifier
    vor dem Merge.
  - **`BEO-PGC/kein-admin-weg-schema-fehler-recovery`** (offen, Ausgang weiter offen):
    `evidence/slice-routing-betriebsdoku.md`, Zähler **1×** → **2×** (V-2: die Abhilfe der
    inkompatiblen Schemaänderung fehlt im Handbuch, am System nur durch Verwerfen des Slots
    erholt); die benannte Hinnahme dieses Slice hat hier ihre Adresse.
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×): **sechster Gegenbeleg, keine
    Datei.** Der Verifier hat die Fixrunde ausgeführt (Wahrheitsprobe am System); die engere
    Fassung trägt auch diesen Fall. Im `state.md` vermerkt, der Ausgang bleibt beim Lese-Schritt
    der Closure von [welle-routing](welle-routing.md).
  - **Deckel, keine Datei:** Review F-2 (LOW; das Handbuch ordnete den Fall „Publication mit
    Spaltenliste bei bekannter Spaltenform“ der falschen Ursache zu, am System vom Reviewer
    gefahren) trifft `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert,
    Deckel bei 10×): vor dem Merge von einem Leser gefunden, ≤ LOW, bekannter Träger-Typ
    (Aussage über den Gegenstand breiter als ihre Messung); steht hier mit Finding-Kennung.
    Review F-3 bis F-5 (LOW, INFO: Transformations-Fehlerblock, Adresse des unbekannten
    Schlüssels, Test- und Dateinamen in Betreiber-Prosa) sind Träger-Nachzüge im Slice und
    haben keinen Eintrag.
- **Folge-Slices:** keine neuen. Übergabe ohne neuen Slice:
  [slice-routing-sdk-beispiel-target](slice-routing-sdk-beispiel-target.md) (SDK-Abschnitte,
  Beispiel-Flags) und [slice-routing-sdk-realserver-e2e](slice-routing-sdk-realserver-e2e.md).
- **Risiken aus §6:** Zwei Quellen **weiter offen** (Register, Lese-Schritt der Welle-Closure);
  Beispiele unter falscher Rolle **entfallen**; Aussage ohne Messung **teils entfallen, teils
  weiter offen** (R4 „höchste `order`“ nicht gefahren; Altserver, Schlüsselspalte bei `DELETE` und
  Abhilfe (a) hergeleitet; Abhilfe der inkompatiblen Schemaänderung nicht beschrieben, V-2);
  Ursprung nicht auflösbar **eingetreten und behoben**; Handbuch nicht nachgezogen/Historie
  **entfallen**; Zwischenzeit mit den SDK-Packages **weiter offen** (Adresse
  `slice-routing-sdk-beispiel-target`).
- **Drei Paarungen:** der Slice gehörte zu [welle-routing](welle-routing.md) — die
  Prüfung lief bei deren Closure (2026-10-02), die DoD-Zeile ist abgehakt. (a) Anker: der
  Lerneintrag verkörpert nichts neu; (b) Folge-Slice: die genannten Pläne
  (`sdk-beispiel-target`, `sdk-realserver-e2e`) liegen als Dateien in `done/`; (c) Register: die genannten Kennungen existieren als Verzeichnis, jede trägt ein
  nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — die SDK-Abschnitte stehen noch aus, das Routing ist für
  Betreiber erst mit ihnen als Ganzes nutzbar; der Nutzer-Bedarf
  ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch den Wellen-Beleg validierbar.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem Pfad
`docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](welle-routing.md) §6:
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (offen, 3×),
`BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig`,
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (verkörpert, 3×),
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, 3×),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×),
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 11×). Der
offene Eintrag `zwei-quellen-drift-handbuch-gegen-pflichtenheft` steht bereits bei 3×;
sein Ausgang gehört dem Lese-Schritt der Welle-Closure.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

