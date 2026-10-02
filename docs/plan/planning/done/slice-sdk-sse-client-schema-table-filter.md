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
[slice-routing-sdk-beispiel-target](../done/slice-routing-sdk-beispiel-target.md)
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
  dem Draht, `null` → Aufruf byte-gleich zum Bestand, `""` → `schema=`/`table=` auf dem
  Draht (der Server liest einen leeren Wert als „kein Filter“, wie beim vorhandenen
  `target=""`), Konjunktion mit `target`; dieselbe Eingabetabelle je Sprache);
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

- [x] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (A): der SSE-Client jedes
      Packages trägt `schema` und `table` als optionale Parameter, die nur gesetzt auf
      dem Draht erscheinen (Prozent-Kodierung wie bei `target`); jede Sprache hat Tests
      für gesetzt, `null` (Parameter fehlt, Aufruf byte-gleich zum Bestand), `""`
      (erscheint als `schema=`/`table=`, vom Server als „kein Filter“ gelesen:
      `parseStreamChangesFilter`, `MatchesFilter`) und die Konjunktion mit `target`,
      mit derselben Eingabetabelle
      (`null`/`""`/`eu`/`a&b=c`/`a+b`/`100%`/`ü`/`a b`; ein Leerzeichen sendet Python
      als `+`, C# und Kotlin als `%20`, der Server dekodiert beides gleich;
      `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). *Zu belegen durch:*
      `make sdk-pack-csharp`, `make sdk-pack-python`, `make sdk-pack-kotlin` (die Tests
      laufen im Bau; ein Lauf aus dem Docker-Schicht-Cache druckt keine Testzeile und
      belegt keine Ausführung — der Belegbefehl ist der Bau mit `--no-cache` der
      Test-Stufe oder eine Mutation, die den Test rot färbt).
- [x] [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (B): die SSE-Beispiele in Go
      (`examples/sse-client/`), C# (`examples/csharp/sse-client/`) und Kotlin
      (`examples/kotlin/sse-client/`) nehmen `schema` und `table` über Flags entgegen;
      ihre Tests belegen Flag → Anfrage. *Zu belegen durch:* `make test` (Go-Beispiel),
      `make examples-csharp`, `make examples-kotlin` (gleiche Cache-Bedingung wie in (A)).
- [x] [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) und Nachzug
      (C): `make sdk-public-doc-check` Exit 0; die READMEs, `examples/README.md` und das
      Benutzerhandbuch (Version und Änderungshistorie) nennen die Parameter; die Sätze
      „`target` ist der einzige Filter“ sind ersetzt, das Zählwort an den SSE-Absätzen
      steht erst nach (A). *Zu belegen durch:* `make sdk-public-doc-check`,
      `make docs-check`, Suchlauf in §3.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). F-1 (MEDIUM) bis F-3 sind durch die Fixrunde geschlossen,
      ein Re-Review nach der Fixrunde fand nicht statt (§7).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-sse-client-schema-table-filter.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: siehe (C); kein SPEC-/ARC-Eintrag. Die Package-Versionen bleiben.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Welle für diesen Slice hier geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Sse/PgChangeFeedSseClient.cs`, `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/PgChangeFeedSseClient.kt`, `sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py` | update | `schema` und `table` als optionale Parameter des SSE-Clients, nach `target` angefügt (kein Verschieben bestehender Parameter), Query-Aufbau mit der vorhandenen Prozent-Kodierung. |
| Testquellen der drei Packages (SSE-Tests je Sprache), die drei READMEs | update | Fake-Transport-Tests je Sprache mit derselben Eingabetabelle, Regression ohne Parameter; README (Englisch, ohne interne Kennung). |
| `examples/sse-client/` (Go), `examples/csharp/sse-client/`, `examples/kotlin/sse-client/` | update | Flags nach dem Muster von `examples/grpc-client/`; Tests Flag → Anfrage. |
| `examples/sse-client/main.go` (`parseConfig`, `config.streamURL`), `examples/sse-client/main_test.go` (neu), `examples/csharp/sse-client/SseStream.cs` (Überladung `StreamUrl(Config)`), `examples/kotlin/sse-client/.../SseStream.kt` (Überladung `streamUrl(Config)`), die Tests `CliTests.cs`/`CliTest.kt` | update | Fixrunde nach dem Review: Die Verdrahtung Flag → Anfrage liegt in einer testbaren Funktion; die Tests fahren die echten Flag-Argumente bis zur URL (Bindung auch für das vorhandene `-target`). |
| Docstring in `sse_client.py` und die SSE-Testtabellen der drei Packages | update | Fixrunde: Eingabetabelle um `a+b`, `100%`, `ü`, `a b` erweitert; die Leerzeichen-Divergenz (Python `+`, C#/Kotlin `%20`) ist im Test benannt und gebunden. |
| `examples/README.md`, `docs/user/benutzerhandbuch.md` (SSE-Absätze, Version, Änderungshistorie) | update | Nachzug (Liefer-Punkt C). |

**Ansatz:** Das Formvorbild ist der `schema`/`table`-Filter des gRPC-Clients der drei
Packages und der Parameter `target` am SSE-Client (letzter Parameter, nur gesetzt auf
dem Draht). Der Implementer liest beide und setzt `schema`/`table` an dieselbe Stelle des
SSE-Query-Aufbaus.

**Umsetzung (Implementer):** Reihenfolge auf dem Draht `schema`, `table`, `target` (wie
`GET /changes`); in den Packages gilt wie bei `target`: `null` fehlt auf dem Draht, ein
leerer Wert erscheint als `schema=` (der Server liest ihn als „kein Filter“,
`MatchesFilter`); in den Beispielen ist leer „nicht gesetzt“ (wie beim vorhandenen
`-target`). Der Server hat keinen Randfall für `table` ohne `schema` (Treffer in jedem
Schema); das Python-Test `test_readme_examples` bindet die API-Tabelle des README an die
Signatur.

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
diff 1 -n -i -E 'only filter|does not set|einzige Filter' -- sdks docs/user examples
diff 67 -n -E 'StreamChangesAsync|stream_changes|streamChanges' -- sdks ':!*Test*' ':!*test*' ':!dist'
diff 43 -n -E 'schema|table' -- sdks/csharp/PgChangeFeed.Client/Sse sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse sdks/python/pgchangefeed/src/pgchangefeed/sse_client.py
diff 88 -n -E 'schema|table' -- examples/sse-client examples/csharp/sse-client examples/kotlin/sse-client
diff 0 -n -E 'ADR-|LH-FA|LH-QA|SPEC-|ARC-' -- sdks ':!dist'
```

| Träger | Messung am Parent (`dd5377bc`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aussagen „`target` ist der einzige Filter“ (READMEs, Handbuch) | Zeile 1: 5 Zeilen (drei READMEs, Handbuch-Absatz am SSE-Stream, Änderungshistorie 1.86) | jede Zeile lesen: die Aussage entfällt oder wird auf den neuen Stand berichtigt; die Historienzeile 1.86 bleibt als Historie stehen. **Befund am Diff:** Zeile 1 am Diff: 1 (die Historienzeile 1.86, bewusst stehen gelassen). Gefunden und berichtigt: die drei READMEs (je der SSE-Absatz, dazu die API-Tabellenzeile und der Abschnitt „Known limits“ „… filtered by delivery target only“, den das Suchmuster nicht trifft), im Handbuch der SSE-Absatz am Ende von „Zugriff über Server-Sent-Events“ (Beispiel-Clients und Packages nehmen alle drei Parameter entgegen), die drei SDK-Abschnitte und die Beispielliste dort; der Docstring-Hinweis „target only“ in `sse_client.py`. Nicht gefunden: weitere Träger in `harness/README.md`, `docs/user/` außerhalb des Handbuchs oder `spec/` (Suche nach `SSE`-Zeilen mit `target`/`schema`/`filter`: nur die E2E-Abdeckungstabellen, die den Server-Pfad beschreiben, nicht die Clients). |
| Stellen der Stream-Methoden | Zeile 2: 67 Nicht-Test-Zeilen | jede SSE-Stelle trägt `schema`/`table` oder die Auslassung ist begründet. **Befund am Diff:** 67 am Diff (Zahl unverändert, weil nur Parameter und Zeilen innerhalb bestehender Zeilen hinzukommen). Die SSE-Stelle je Package ist `StreamChangesAsync`/`streamChanges`/`stream_changes` des SSE-Clients (jetzt mit `schema`/`table`), dazu die README-Zeilen (nachgezogen). Die gRPC- und NATS-Stellen sind nicht Gegenstand (gRPC trägt die Parameter bereits, NATS filtert über das Subjekt). |
| `schema`/`table` im SSE-Teil der Packages | Zeile 3: 14 Zeilen (Nachrichtenmodell `Change` der Antwort, Docstring), kein Treffer am Aufrufparameter | die Treffer sind die Felder der empfangenen Change, nicht der Anfrage; der Anfrageparameter kommt hinzu. **Befund am Diff:** 43 am Diff; der Anfrageparameter steht jetzt in allen drei Clients (Signatur, Query-Aufbau, Docstring). Die 14 Treffer des Parents (Felder der Antwort) bleiben unverändert. |
| Flags der SSE-Beispiele | Zeile 4: 6 Zeilen (Testdaten der Antwort), kein Flag | das Paar kommt je Sprache hinzu. **Befund am Diff:** 88 am Diff (nach der Fixrunde: 76 + die Verdrahtungs-Tests und -Funktionen); das Flag-Paar steht in Go (`main.go`), C# (`Cli.cs`) und Kotlin (`Cli.kt`), jeweils mit Tests Flag → Anfrage, die die echten Flag-Argumente bis zur URL fahren. |
| Interne Kennungen unter `sdks/` | Zeile 5: 0 Zeilen | muss 0 bleiben; `make sdk-public-doc-check` ist der Wächter. **Befund am Diff:** 0 am Diff, `make sdk-public-doc-check` Exit 0 (siehe Bericht). |

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
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). — **Ausgang: entfallen.**
  Dieselbe Eingabetabelle (`null`/`""`/`eu`/`a&b=c`/`a+b`/`100%`/`ü`/`a b`) steht in allen
  drei Sprachen; die erwarteten Draht-Bytes sind gleich bis auf das Leerzeichen (Python `+`,
  C# und Kotlin `%20`), das im Test benannt und dekodiert-äquivalent gebunden ist (Server
  dekodiert beide Formen zum selben Wert; gemessen vom Verifier, Verifikation §4).
- **Test läuft nicht, Bau grün.** Ein Lauf aus dem Docker-Schicht-Cache druckt keine
  Testzeile. — **Ausgang: eingetreten, gelöst.** Alle fünf `make sdk-pack-*`-/`make
  examples-*`-Ziele liefen bei Reviewer und Verifier mit Exit 0 aus dem Cache ohne Testzeile;
  der Verifier holte die Zahlen mit `docker build --no-cache --target build` nach (Python
  `150 passed`, C# `Passed: 159`, C#-SSE-Beispiel `Passed: 25`) und sah zehn Einzelmutationen
  rot (Verifikation §1 bis §3).
- **Binärkompatibilität bei künftigem Release.** Die neuen Parameter (C# optionale
  Parameter, Kotlin ohne `@JvmOverloads`) ändern die binäre Signatur; gegen 0.2.x
  kompilierte Aufrufer brauchen Neukompilierung (quellkompatibel). — **Ausgang: weiter
  offen** bis zum nächsten SDK-Release (Freigabe des Auftraggebers, nicht Teil dieses
  Slice); Hinweis für die Release-Notiz: C# und Kotlin neu kompilieren.
- **Interne Kennung in öffentlichem Text.** Docstrings, KDoc, XML-Doku und Fehlertexte
  erreichen die Anwender über die Packages. — **Ausgang: entfallen**
  (`make sdk-public-doc-check` Exit 0, Suchlauf-Zeile `diff 0`).

## 7. Closure-Notiz

Ursprung der Angaben: **gemessen** = vom Reviewer
([`review-slice-sdk-sse-client-schema-table-filter`](../../../reviews/review-slice-sdk-sse-client-schema-table-filter.md))
oder Verifier
([`verifikation-slice-sdk-sse-client-schema-table-filter`](../../../reviews/verifikation-slice-sdk-sse-client-schema-table-filter.md))
im eigenen Lauf; **übernommen** = aus einem dieser Berichte ohne eigene Nachmessung. Der
Planner hat diese Zahlen nicht nachgemessen (übernommen), ausgenommen der Zählwort-Suchlauf
unten (gemessen, Closure).

- **Was hat funktioniert:** Die Eingabetabelle je Sprache (Plan-Pflicht) hielt die Draht-Bytes
  der drei Packages gleich (gemessen, Verifikation §4). Der Reviewer fand die Verdrahtungs-Lücke
  F-1 mit Mutationen vor dem Merge. Der Verifier sah zehn Einzelmutationen rot (G1 bis G3, C1,
  C2, K1, K2, P-C#, P-Py, P-Kt; gemessen) und `make gates` Exit 0 (gemessen).
- **Was ging anders als geplant:** F-1 (MEDIUM: Verdrahtung Flag → Anfrage der Beispiele
  ungebunden) bis F-3 (LOW) wurden durch die Fixrunde `81ad7643`/`410e22ac`/`9b1023c7`
  geschlossen; Beleg ist die Verifikation (zehn Einzelmutationen rot). **Ein Re-Review nach der
  Fixrunde fand nicht statt (ehrliche Abweichung);** die Review-Reports tragen deshalb weiter
  „merge-blockierend: ja" bei F-1 (Records, unverändert). Der Verifier hielt dies als
  Bedingung B1 fest und stufte einen Re-Review als nicht erforderlich ein, weil die Fixrunde
  nur Beispiel-Verdrahtung, Tests und Plantext änderte. Verbleibende LOW-Lücke: der Aufruf
  der URL-Funktion in `main`/`Program.cs`/`Main.kt` ist von keinem Test gebunden
  (hergeleitet, nicht gemutet; das vorhandene `-target` trug sie schon vorher). INFO: das
  Go-Beispiel beendet `-h` jetzt mit Exit 2 statt 0.
- **Steering-Loop-Eintrag (Lerneintrag):** Die Verdrahtung Flag → Anfrage ist eine eigene
  Bindung; Parser und URL-Bau einzeln zu testen genügt nicht, der Test fährt die echten
  Flag-Argumente bis zur URL. Das ist die Anwendung der verkörperten Regel
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` („mutiere den Eingabewert"), kein neuer
  Eintrag. Neu und real aufgetreten: die Docker-Cache-Falle. Alle fünf `make sdk-pack-*`-/
  `make examples-*`-Ziele sind aus dem Schicht-Cache grün ohne Testzeile; der Belegbefehl ist
  `docker build --no-cache` an der Test-Stufe oder eine rote Mutation, nie der Exit des Ziels
  allein.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/docker-cache-ueberspringt-tests-still`
  um `evidence/slice-sdk-sse-client-schema-table-filter.md` fortgeschrieben, Zähler **2×**
  (Schwelle 3× nicht erreicht, kein Ausgang zu entscheiden).
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (gestrichen, 3×): Vermerk „Gegenmaßnahme
  beobachtet, kein neuer Anfall" im state, Zähler bleibt 3×, keine neue Datei.
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 23×, Schwelle bereits
  erreicht und beschieden): F-1 ist ein weiterer Träger-Typ-Treffer, vom Reviewer vor dem Merge
  gefunden, Schwere MEDIUM; die Eintragsform verlangt für Schwere über LOW eine `evidence/`-Datei,
  deshalb weder Deckel-Vermerk noch Zählung hier entschieden. Die Notiz meldet den Fall: ob eine
  `evidence/`-Datei (24×) oder ein Verweis in dieser Notiz genügt, ist Sache des nächsten
  Lese-Schritts; dieser Slice hat sie nicht angelegt.
- **Folge-Slices:** keiner angelegt. Empfehlung (optional, nicht angelegt): ein Realserver-Beleg
  für `schema`/`table` am SSE-Client der drei Packages und der Beispiele als eigener Slice in
  `open/` im Muster `make test-sdk-*-integration`. Grund: Das Handbuch sagt, Packages und
  Beispiele nähmen alle drei Parameter entgegen; belegt ist das nur auf Unit-Ebene mit
  Fake-Transport, ein Realserver-Beleg war laut Plan §1 ausgeschlossen.
- **Zählwort „drei optionale":** Suchlauf am Arbeitsbaum
  (`git grep -n -i -E 'drei optionale' -- sdks docs/user examples`, gemessen bei der Closure):
  vier Treffer, alle im Handbuch (Zeilen am gRPC-Request, am SSE-Absatz und an der
  `GET /changes`-Passage); keine in den READMEs. Das Zählwort steht an den SSE-Absätzen nach (A).
- **Risiken aus §6:** Randfall-Kopie entfallen · Test läuft nicht (Cache) eingetreten, gelöst ·
  Binärkompatibilität (C# und Kotlin) weiter offen bis zum SDK-Release, Hinweis für die
  Release-Notiz · interne Kennung entfallen (Ausgänge in §6).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `sdks/`, `examples/` und `docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); Treffer:
`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 3× nach der Closure von
[slice-routing-sdk-beispiel-target](../done/slice-routing-sdk-beispiel-target.md)),
`BEO-PGC/docker-cache-ueberspringt-tests-still` (offen, 1×),
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert, 3×),
`BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen, 1×),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
