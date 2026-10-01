# Slice sdk-sse-client-schema-table-filter: SSE-Client der drei SDK-Packages und SSE-Beispiele — die optionalen Parameter `schema` und `table`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming),
[`LH-FA-SST-009`](../../../../spec/lastenheft.md) (SDK-Packages),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Teilfrage 4:
SSE nutzt dieselben Query-Parameter `schema`/`table` wie der gRPC-Stream und
`GET /changes`),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`),
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) (Beispiel-Clients).

**Berührte Spec-Stellen:** — (die Wire-Form des SSE-Endpunkts steht in
[`SPEC-021`](../../../../spec/pflichtenheft.md); die SDKs setzen sie um und ändern
sie nicht).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Closure von
[slice-routing-sdk-beispiel-target](../in-progress/slice-routing-sdk-beispiel-target.md)
(dort §6, Aufschub „`schema`/`table` am SSE-Client der Packages“). **Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der SSE-Client jedes der drei SDK-Packages (C# `PgChangeFeed.Client`, Kotlin
`pgchangefeed-kotlin`, Python `pgchangefeed`) und die SSE-Beispiele in Go, C# und
Kotlin setzen die Query-Parameter `schema` und `table` des Endpunkts
`GET /changes/stream`, die der Server seit
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) auswertet. Heute
ist `target` der einzige Filter der SSE-Clients der Packages und des SSE-Beispiels;
ein Anwender, der nach Tabelle filtern will, baut die Anfrage selbst. Der Parameter
ist je Sprache additiv und optional (dieselbe Form wie `target`: letzter Parameter,
nur gesetzt auf dem Draht). Drei Liefer-Punkte:

- (A) **SDK-Packages:** der SSE-Client der drei Packages nimmt `schema` und `table`
  optional entgegen; Unit-Tests mit Fake-Transport je Sprache (gesetzt → Parameter auf
  dem Draht, leer → Aufruf byte-gleich zum Bestand, Konjunktion mit `target`; dieselbe
  Eingabetabelle je Sprache);
- (B) **SSE-Beispiele:** die Flags `-schema`/`-table` (Go) bzw. `--schema`/`--table`
  (C#, Kotlin) am SSE-Beispiel, mit den Tests des jeweiligen Beispiels;
- (C) **Öffentlicher Text und Nachzug:** READMEs der drei Packages, `examples/README.md`
  und die SSE-Passagen des Benutzerhandbuchs; die Aussage „`target` ist dort der
  einzige Filter“ entfällt, das Zählwort „drei optionale“ steht erst nach (A) am
  SSE-Absatz.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release oder eine Versionsänderung der Packages** — jedes SDK-Release ist eine
  Freigabe des Auftraggebers (neue Version, Tag-Workflow, externes Konto).
- **Der NATS-Vollinhalts-Client** — der Weg filtert über das Subjekt
  (`cdc.stream.<source_id>.<schema>.<table>`), nicht über Query-Parameter; sein
  Subjekt-Bau ist unverändert.
- **Realserver-Belege der SDK-Tiers** — dieser Slice belegt auf Unit-Ebene mit
  Fake-Transport; ein Beleg am laufenden Server ist ein eigener Zug im Muster der
  SDK-Realserver-Tiers (`make test-sdk-*-integration`).
- **Server-Änderungen** — der Server filtert seit
  [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md); ein gefundener
  Server-Fehler wird gemeldet, nicht hier repariert.

## 2. Definition of Done

- [ ] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (A): der SSE-Client jedes
      Packages trägt `schema` und `table` als optionale Parameter, die nur gesetzt auf
      dem Draht erscheinen (Prozent-Kodierung wie bei `target`); jede Sprache hat Tests
      für gesetzt, leer und die Konjunktion mit `target`, mit derselben Eingabetabelle
      (`null`/`""`/`eu`/`a&b=c`;
      `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). *Zu belegen durch:*
      `make sdk-pack-csharp`, `make sdk-pack-python`, `make sdk-pack-kotlin` (die Tests
      laufen im Bau; ein Lauf aus dem Docker-Schicht-Cache druckt keine Testzeile und
      belegt keine Ausführung — der Belegbefehl ist der Bau mit `--no-cache` der
      Test-Stufe oder eine Mutation, die den Test rot färbt).
- [ ] [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (B): die SSE-Beispiele in Go
      (`examples/sse-client/`), C# (`examples/csharp/sse-client/`) und Kotlin
      (`examples/kotlin/sse-client/`) nehmen `schema` und `table` über Flags entgegen;
      ihre Tests belegen Flag → Anfrage. *Zu belegen durch:* `make test` (Go-Beispiel),
      `make examples-csharp`, `make examples-kotlin` (gleiche Cache-Bedingung wie in (A)).
- [ ] [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) und Nachzug
      (C): `make sdk-public-doc-check` Exit 0; die READMEs, `examples/README.md` und das
      Benutzerhandbuch (Version und Änderungshistorie) nennen die Parameter; die Sätze
      „`target` ist der einzige Filter“ sind ersetzt, das Zählwort an den SSE-Absätzen
      steht erst nach (A). *Zu belegen durch:* `make sdk-public-doc-check`,
      `make docs-check`, Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-sse-client-schema-table-filter.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: siehe (C); kein SPEC-/ARC-Eintrag. Die Package-Versionen bleiben.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Welle für diesen Slice hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Sse/PgChangeFeedSseClient.cs`, `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/PgChangeFeedSseClient.kt`, `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py` | update | `schema` und `table` als optionale Parameter des SSE-Clients, nach `target` angefügt (kein Verschieben bestehender Parameter), Query-Aufbau mit der vorhandenen Prozent-Kodierung. |
| Testquellen der drei Packages (SSE-Tests je Sprache), die drei READMEs | update | Fake-Transport-Tests je Sprache mit derselben Eingabetabelle, Regression ohne Parameter; README (Englisch, ohne interne Kennung). |
| `examples/sse-client/` (Go), `examples/csharp/sse-client/`, `examples/kotlin/sse-client/` | update | Flags nach dem Muster von `examples/grpc-client/`; Tests Flag → Anfrage. |
| `examples/README.md`, `docs/user/benutzerhandbuch.md` (SSE-Absätze, Version, Änderungshistorie) | update | Nachzug (Liefer-Punkt C). |

**Ansatz:** Das Formvorbild ist der `schema`/`table`-Filter des gRPC-Clients der drei
Packages und der Parameter `target` am SSE-Client (letzter Parameter, nur gesetzt auf
dem Draht). Der Implementer liest beide und setzt `schema`/`table` an dieselbe Stelle des
SSE-Query-Aufbaus.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der SSE-Client der Packages
und das SSE-Beispiel setzen `schema`/`table` nicht, `target` ist ihr einziger Filter“;
Parent ist `dd5377bc`; der Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und
Nichtgefundenes ein):**

```suchlauf
dd5377bc 5 -n -i -E 'only filter|does not set|einzige Filter' -- sdks docs/user examples
dd5377bc 67 -n -E 'StreamChangesAsync|stream_changes|streamChanges' -- sdks ':!*Test*' ':!*test*' ':!dist'
dd5377bc 14 -n -E 'schema|table' -- sdks/csharp/PgChangeFeed.Client/Sse sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py
dd5377bc 6 -n -E 'schema|table' -- examples/sse-client examples/csharp/sse-client examples/kotlin/sse-client
dd5377bc 0 -n -E 'ADR-|LH-FA|LH-QA|SPEC-|ARC-' -- sdks ':!dist'
```

| Träger | Messung am Parent (`dd5377bc`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aussagen „`target` ist der einzige Filter“ (READMEs, Handbuch) | Zeile 1: 5 Zeilen (drei READMEs, Handbuch-Absatz am SSE-Stream, Änderungshistorie 1.86) | jede Zeile lesen: die Aussage entfällt oder wird auf den neuen Stand berichtigt; die Historienzeile 1.86 bleibt als Historie stehen. **Befund am Diff:** vom Implementer einzutragen. |
| Stellen der Stream-Methoden | Zeile 2: 67 Nicht-Test-Zeilen | jede SSE-Stelle trägt `schema`/`table` oder die Auslassung ist begründet. **Befund am Diff:** vom Implementer einzutragen. |
| `schema`/`table` im SSE-Teil der Packages | Zeile 3: 14 Zeilen (Nachrichtenmodell `Change` der Antwort, Docstring), kein Treffer am Aufrufparameter | die Treffer sind die Felder der empfangenen Change, nicht der Anfrage; der Anfrageparameter kommt hinzu. **Befund am Diff:** vom Implementer einzutragen. |
| Flags der SSE-Beispiele | Zeile 4: 6 Zeilen (Testdaten der Antwort), kein Flag | das Paar kommt je Sprache hinzu. **Befund am Diff:** vom Implementer einzutragen. |
| Interne Kennungen unter `sdks/` | Zeile 5: 0 Zeilen | muss 0 bleiben; `make sdk-public-doc-check` ist der Wächter. **Befund am Diff:** vom Implementer einzutragen. |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-sdk-beispiel-target` liegt in
`done/` und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Voraussetzung am
Start: Docker mit Netzzugang für die drei `make sdk-pack-*`-Ziele und `make examples-*`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): braucht eine Sprache mehr als
  drei Fixrunden, teilt sich der Slice an der Sprach-Wurzel (`sdks/<sprache>/`).
- `in-progress` → `open` (blockiert): der SSE-Endpunkt wertet `schema`/`table` am
  laufenden Server nicht wie in [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md)
  beschrieben aus — der Fund geht an den Architect, kein Umgehen im SDK.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), die drei `make sdk-pack-*`-Ziele und `make examples-csharp`/
`make examples-kotlin` mit ausgeführten Tests belegt (Bau ohne Cache der Test-Stufe oder
rote Mutation), Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Drei Sprachen, ein Randfall, drei Lesarten.** Leerer Wert gegen nicht gesetzt,
  Zeichen im Namen, Kodierung im Query-Parameter
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). — **Ausgang:** bei der Closure
  einzutragen (dieselbe Eingabetabelle je Sprache).
- **Test läuft nicht, Bau grün.** Ein Lauf aus dem Docker-Schicht-Cache druckt keine
  Testzeile. — **Ausgang:** bei der Closure einzutragen (Beleg: gedruckte Testzahl im Bau
  ohne Cache oder rote Mutation).
- **Interne Kennung in öffentlichem Text.** Docstrings, KDoc, XML-Doku und Fehlertexte
  erreichen die Anwender über die Packages. — **Ausgang:** bei der Closure einzutragen
  (`make sdk-public-doc-check` Exit 0).

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
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `sdks/`, `examples/` und `docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); Treffer:
`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 3× nach der Closure von
[slice-routing-sdk-beispiel-target](../in-progress/slice-routing-sdk-beispiel-target.md)),
`BEO-PGC/docker-cache-ueberspringt-tests-still` (offen, 1×),
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert, 3×),
`BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen, 1×),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
