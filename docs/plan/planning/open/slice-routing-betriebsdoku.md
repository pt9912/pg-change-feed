# Slice routing-betriebsdoku: Benutzerhandbuch — Routing-Regeln konfigurieren, Ziele lesen, Konsequenzen und Abhilfe

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (SQL-Administration),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Zugriffswege),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 8 (Handbuch-Hälfte — Träger dieses Slice) und Entscheidung 4
(Handbuch-Hinweis zu DELETE und Inhaltsregeln).

**Berührte Spec-Stellen:** — (das Handbuch beschreibt Betreiber-Oberfläche und
verweist für Zusagen auf die Spec; es ist selbst keine Spec-Stelle).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
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
([welle-routing](../welle-routing.md)) — diese Gegenstände haben die Slices mit
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
  (Lauf vom 2026-10-01 am Arbeitsbaum von `slice-routing-e2e`, je einmal mit
  `postgres:18-alpine` = PostgreSQL 18.6 und `postgres:17-alpine` = PostgreSQL 17.11,
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

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (A): der Abschnitt „Routing-Regel
      konfigurieren" steht im Abschnitt „Aufgaben"; jedes SQL-Beispiel ist unter der
      genannten Rolle (`cdc_admin`) gegen die Compose-Umgebung ausgeführt worden, die
      gedruckte Antwort steht im Bericht
      (`BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig`); R1–R6, die
      Konsequenzen, die Abhilfe und der DELETE-Hinweis stehen mit ihrem Ursprung (gemessen
      von `slice-routing-e2e`, übernommen aus der Spec, oder als erwartet gekennzeichnet);
      der DELETE-Hinweis nennt die Menge der Messung (PostgreSQL-Versionen und Fälle).
      *Zu belegen durch:* Lesen des Abschnitts gegen den Bericht von `slice-routing-e2e`,
      Ausführen der Beispiele, `make docs-check`.
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (B): die Abschnitte zu den
      Zugriffswegen nennen `target` und das NATS-Subjekt, „Änderungen lesen" die Spalte
      `route_target`; die Aussage zu Altservern und die Kosten-Aussage stehen mit ihrem
      Ursprung; was für die SDK-Packages gilt, steht als „folgt mit den SDK-Packages"
      oder gar nicht, bis `slice-routing-sdk-beispiel-target` es einlöst.
      *Zu belegen durch:* Lesen, `make docs-check`, Suchlauf in §3.
- [ ] (C): die Zeile `schema` der Fehlerklassen-Tabelle und der Abschnitt „Neustart nach
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
      Kopfzeile (`BEO-PGC/handbuch-versionshistorie-uebersprungen`).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-betriebsdoku.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: der Slice **ist** das Doku-Update des Handbuchs; `README.md` im
      Repo-Wurzelverzeichnis nur, soweit der Suchlauf eine bewegte Beschreibung findet.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` §4 „Aufgaben" (neuer Abschnitt „Routing-Regel konfigurieren", Formvorbild „Transformationsregel konfigurieren") | neu | Liefer-Punkt A. |
| `docs/user/benutzerhandbuch.md` §4 („Änderungen lesen", „Bestand als Backfill überführen", „Zugriff über …": HTTP, gRPC-Stream, gRPC-Verwaltungs-API, SSE, NATS-Vollinhalt) | update | Liefer-Punkt B. |
| `docs/user/benutzerhandbuch.md` §2 „Zugriff und Rollen", §6 „Fehlerklassen" und „Neustart nach einem Fehler", §9 „Grenzwerte" und „Änderungshistorie", Kopfzeile (Version, Stand) | update | Liefer-Punkt C. |
| `docs/user/benutzerhandbuch.md` Transformationen, „Fehler: Erfassung endet mit der Fehlerklasse `schema`“ (Bestandsbefund, Parent `2629d544` etwa Z. 387; *hergeleitet* nach `ADR-0112` Teilfrage 4) | update | Handbuch-Nachzug zu `ADR-0140`: Ursache „`column` fehlt in der Relation“ mit Lösung „Regel entfernen, Neustart“ ist für die entfernte Spalte falsch (der Pfad endet dort schon an der inkompatiblen Schemaänderung); die Ursachen trennen wie im Routing-Block. Das Handbuch wird in diesem Plan nicht geändert, der Zug gehört diesem Slice. |
| `docs/user/benutzerhandbuch.md` „Änderungen lesen“ und Feldliste (Bestandsbefund Review F-4, Parent `2629d544`: `:516` „`origin` ist die letzte Spalte der View“, `:1462` Feldliste endet bei `origin`) | update | `route_target` ist seit der Spalte `route_target` die letzte Spalte der View `cdc.changes`; beide Stellen nachziehen (Liefer-Punkt B). |
| `README.md` | lesen | nur, wenn der Suchlauf eine bewegte Beschreibung findet. |

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
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, an denen die Transformations-Funktionen stehen (Muster für die Routing-Funktionen) | Zeile 1: 10 Zeilen in `docs/user` | jede Stelle lesen: Rollen-Tabelle und Aufzählungen, die um die zwei Routing-Funktionen wachsen; Befund am Diff: einzutragen |
| Beschreibungen von `cdc.changes` | Zeile 2: 36 Zeilen | jede Stelle, die die Spalten der View aufzählt, trägt `route_target`; Befund: einzutragen |
| Version und Historie | Zeilen 3 und 4: 1 Kopfzeile, 4 Historie-Zeilen der Reihe 1.8x | die nächste Version ist 1.84 (*hergeleitet*; der Implementer misst am Start, falls ein anderer Slice die Version bewegt hat); Kopfzeile und Historie-Zeile tragen dieselbe Version; Befund: einzutragen |
| Filter-Beschreibung der Clients (`-schema`/`-table`) und der Stream-/Lese-Wege (`schema` und `table`) | Zeile 5: 18 Zeilen (das Muster trifft auch die Schreibweise mit Backticks um beide Wörter, `benutzerhandbuch.md:1149`), darunter `benutzerhandbuch.md:1308` („Go, C# und Kotlin nehmen den Filter über `-schema`/`-table`") | **benannte Zwischenzeit:** diese Aussage zum `target`-Flag der Clients zieht `slice-routing-sdk-beispiel-target` nach; dieser Slice schreibt dazu nichts Behauptendes (§2 B). Befund: einzutragen |

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
  wiederholen).
- **Beispiele unter der falschen Rolle.** Ein SQL-Beispiel, das unter `cdc_admin`
  läuft, aber ohne Rollenangabe dasteht, scheitert beim Leser
  (`BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig`). — **Ausgang:**
  bei der Closure einzutragen (Beispiele gelaufen, gedruckte Zeile).
- **Aussage ohne Messung.** Die DELETE-Aussage, die Altserver-Aussage und die
  Kosten-Aussage sind Erwartungen, bis ihr Slice sie belegt hat
  (`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`, verkörpert, 11×). —
  **Ausgang:** bei der Closure einzutragen (jede Aussage mit Ursprung: gemessen,
  übernommen, abgeleitet, erwartet).
- **Handbuch nicht nachgezogen / Historie übersprungen**
  (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`, 3×;
  `BEO-PGC/handbuch-versionshistorie-uebersprungen`, 3×; beide verkörpert). —
  **Ausgang:** bei der Closure einzutragen.
- **Zwischenzeit mit den SDK-Packages.** Zwischen diesem Slice und
  `slice-routing-sdk-beispiel-target` beschreibt das Handbuch `target` am Roh-Weg, die
  SDK-Abschnitte nicht. — **Ausgang:** bei der Closure einzutragen (Suchlauf-Zeile 5;
  Nachzug dort).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem Pfad
`docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (offen, 3×),
`BEO-PGC/handbuch-beispiel-nicht-unter-genannter-rolle-lauffaehig`,
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (verkörpert, 3×),
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, 3×),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×),
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 11×). Der
offene Eintrag `zwei-quellen-drift-handbuch-gegen-pflichtenheft` steht bereits bei 3×;
sein Ausgang gehört dem Lese-Schritt der Welle-Closure.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

