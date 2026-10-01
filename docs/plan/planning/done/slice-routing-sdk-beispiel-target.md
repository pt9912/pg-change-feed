# Slice routing-sdk-beispiel-target: SDK-Packages und Beispiel-Clients — der Parameter `target` an allen Zustellweg-Flächen

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
Haupt-Bezug), [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (SDK-Packages),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Gleichwertigkeit der Zugriffswege),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 8 (SDK-/Beispiel-Hälfte — Träger dieses Slice) und Teilfrage 5,
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Formvorbild:
`schema`/`table` in den SDKs und Beispielen),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`),
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) und
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) (Beispiel-Clients).

**Berührte Spec-Stellen:** — (die Wire-Formen stehen in
[`SPEC-020`](../../../../spec/pflichtenheft.md) bis
[`SPEC-024`](../../../../spec/pflichtenheft.md); die SDKs setzen sie um und ändern
sie nicht).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die drei SDK-Packages (C# `PgChangeFeed.Client`, Kotlin
`pgchangefeed-kotlin`, Python `pgchangefeed`) und die Beispiel-Clients (Go, C#,
Kotlin) wählen das Ziel einer Change an den Flächen, die der Server mit `target`
bedient — der HTTP-Lesezugriff, der gRPC-Stream (und der RPC `ReadChanges` nach
`ADR-0138` Festlegung 1), der SSE-Stream — und abonnieren das Zusatz-Subjekt des NATS-Vollinhalts-
Wegs; ohne Angabe bleibt jeder Aufruf unverändert (additiv und optional). Drei
Liefer-Punkte:

- (A) **SDK-Packages:** der optionale Parameter je Fläche in allen drei Sprachen
  (Formvorbild: der `schema`/`table`-Filter, `ADR-0133`), mit Unit-Tests je Fläche
  und einem Regressionstest für den Aufruf ohne Parameter; die README jedes Packages
  beschreibt ihn (Englisch, ohne interne Kennung);
- (B) **Beispiel-Clients:** das Flag für das Ziel an den Beispielen der vier Flächen in
  Go, C# und Kotlin (Formvorbild: `-schema`/`-table`, Parent: 54 Zeilen unter
  `examples/`), mit den Tests des jeweiligen Beispiels und `examples/README.md`;
- (C) **Öffentlicher Text und Nachzug:** `make sdk-public-doc-check` bleibt grün (keine
  Kennung der Spezifikations-, Entscheidungs- und Anforderungsdokumente und kein Slice-/
  Welle-Name unter `sdks/`); die SDK-Passagen des Benutzerhandbuchs nennen den Parameter
  (der Nachzug zu `slice-routing-betriebsdoku`: Zwischenzeit endet hier).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release oder eine Versionsänderung der Packages** — jedes SDK-Release ist eine
  Freigabe des Auftraggebers (neue Version, Tag-Workflow, externes Konto); diese Welle
  ändert Quellen, nicht Versionen (Welle §6).
- **Realserver-Belege der SDK-Tiers** (`make test-sdk-*-integration`) —
  `slice-routing-sdk-realserver-e2e` (Slice 10 der Welle, Entscheidung des
  Auftraggebers); dieser Slice belegt auf Unit-Ebene mit Fake-Transport und
  Fake-Invoker (Formvorbild der Slices der Welle
  `welle-sdk-grpc-administration-flaeche`) und dem Bau der Packages.
- **Änderung des Nachrichtenmodells der SDKs** — das Label ist nicht Teil der
  Nachrichten (zehn bzw. dreizehn Felder bleiben); die Decoder bekommen kein neues Feld
  zu lesen.
- **Ein Admin-Client für `set_route`/`remove_route`** — die Antragsarten haben keinen
  gRPC-Weg (`slice-routing-antragsweg` §1); der Betreiber konfiguriert über SQL.
- **Server-Änderungen** — alle Wire-Formen stehen aus den Slices davor; ein gefundener
  Server-Fehler wird gemeldet, nicht hier repariert (anderer Vorgang).

## 2. Definition of Done

- [x] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) und
      [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (A): in jedem der drei Packages
      tragen der HTTP-Lesezugriff, der gRPC-Stream, der SSE-Client (und, nach
      `ADR-0138` Festlegung 1 mit `ReadChangesRequest.target = 7`, der
      `ReadChanges`-Aufruf des Administration-Clients) einen optionalen
      Parameter für das Ziel, der nur gesetzt auf dem Draht erscheint; der NATS-Stream-
      Client bietet das Abonnement des Zusatz-Subjekts je Ziel (Subjekt-Bau mit der
      Prüfung der reservierten Zeichen wie beim Bestand); jede Fläche hat einen Test
      mit der Eingabe (Ziel gesetzt → Parameter auf dem Draht; leer → Aufruf byte-gleich
      zum Bestand) und das Verhalten bei einem ungültigen Ziel folgt der Spec
      (clientseitige Prüfung nur, soweit die Spec sie zusagt); die drei Sprachen lesen
      dieselben Randfälle gleich (Test mit derselben Eingabetabelle je Sprache,
      `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). *Zu belegen durch:*
      `make sdk-pack-csharp`, `make sdk-pack-python`, `make sdk-pack-kotlin` (die Tests
      laufen im Bau; ein roter Test bricht ihn ab; Docker-only, braucht Netz).
      *Beleg:* [Verifikations-Report](../../../reviews/verifikation-slice-routing-sdk-beispiel-target.md)
      §1 und §3 — die drei `make sdk-pack-*` kamen dort aus dem Docker-Schicht-Cache ohne
      Testzeile und belegen keine Ausführung; die Ausführung belegen die Test-Stufen mit
      `docker build --no-cache` (gedruckt: Python `144 passed`, C# `Passed: 137`; Kotlin druckt
      im grünen Lauf keine Testzahl, `:test` lief) und zehn selbst gefahrene Mutationen, alle rot.
      Wire-Wirkung gegen einen laufenden Server ist nicht belegt (Menge: Unit-Ebene,
      Fake-Transport, Fake-Invoker).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (B): die Beispiel-Clients der
      vier Flächen in Go, C# und Kotlin (Parent: Go unter `examples/<fläche>-client/`,
      C# unter `examples/csharp/<fläche>-client/`, Kotlin unter `examples/kotlin/`)
      nehmen das Ziel über ein Flag nach dem Muster von `-schema`/`-table`; die Tests
      des jeweiligen Beispiels belegen Flag → Anfrage; `examples/README.md` nennt es. Die
      Beispiele bleiben Doku mit Bau-Bindung, kein Lauf-Beleg. *Zu belegen durch:*
      `make test` (Go-Beispiele), `make examples-csharp`, `make examples-kotlin`
      (Docker-only, braucht Netz). *Beleg:* Verifikations-Report §1 und §3 — `make test` Exit 0 (fünf
      Beispiel-Pakete `ok`, nicht gecacht); C#-Beispiele im Bau ohne Cache (`Passed: 54`, `12`, `9`,
      `52`, `20`, je `Failed: 0`), Kotlin-Beispiele im Bau ohne Cache (`:test` aller fünf Module
      ausgeführt, keine Testzahl gedruckt); vier Go-Mutationen rot. Die C#- und Kotlin-Beispiele sind
      vom Verifier nicht mutiert (nur ausgeführt gelesen).
- [x] [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) und Nachzug
      (C): `make sdk-public-doc-check` Exit 0 (keine Kennung `SPEC-`, `ADR-`, `ARC-`,
      `LH-FA-`, `LH-QA-` und kein Slice-/Welle-Name unter `sdks/`, auch nicht in
      Docstrings, KDoc, XML-Doku, Fehlertexten); die generierten Python-Stubs prüft
      `sdks/python/pgchangefeed/tests/test_public_text.py` im Bau; die README-Texte sind
      Englisch ohne deutsche Fachwörter
      (`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`); die SDK-Passagen in
      `docs/user/benutzerhandbuch.md` und `examples/README.md` nennen den Parameter, die
      Aussage „Go, C# und Kotlin nehmen den Filter über `-schema`/`-table`" ist
      vollständig; Handbuch-Version und Änderungshistorie tragen die Zeile.
      *Zu belegen durch:* `make sdk-public-doc-check`, `make docs-check`, Suchlauf in §3.
      *Beleg:* `make sdk-public-doc-check` Exit 0 („keine interne Kennung unter sdks“), Handbuch
      `Version: 1.86` mit Historienzeile 1.86 als letzter Zeile (Verifikations-Report §4 Zeile 3).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9;
      `sdk-public-doc-check` ist Teil von `make gates`). *Beleg:* Verifikations-Report §1
      (Exit 0) und der Lauf der Closure (§7).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-sdk-beispiel-target`](../../../reviews/review-slice-routing-sdk-beispiel-target.md)
      (0 HIGH, 1 MEDIUM F-1 durch die Fixrunde `9e8bc1ce` geschlossen, 2 LOW, 2 INFO) und
      Verifier-Gegenprüfung der Fixrunde
      [`verifikation-slice-routing-sdk-beispiel-target`](../../../reviews/verifikation-slice-routing-sdk-beispiel-target.md)
      (F-1 bis F-5 am Stand `HEAD` gelesen, Test-Stufen ohne Cache ausgeführt, zehn Mutationen). Es gab
      **kein separates Re-Review** der Fixrunde: engere Fassung von
      `BEO-PGC/fixrunde-ohne-reviewer-lesung` (Produktionslogik oder Norm geändert, oder kein anderer
      Kontext hat sie ausgeführt) — die Fixrunde änderte nur READMEs, einen Docstring, das Handbuch und den
      Plan, und der Verifier hat sie ausgeführt (Verifikations-Report §7).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-sdk-beispiel-target.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg:* Exit 0, „18 Zeilen
      stimmen“ (Verifikations-Report §1 und der Lauf der Closure, §7).
- [x] Doku-Update: SDK-READMEs, `examples/README.md`, Handbuch (siehe C); kein
      SPEC-/ARC-Eintrag. Die Package-Versionen bleiben (Abgrenzung §1). *Beleg:* Verifikations-Report §2
      (`git diff --stat`: keine Datei unter `spec/`, `proto/`, `internal/`, `cmd/`, `gen/`; keine
      Versionsdatei, kein Tag).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7) — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/` (`Http/PgChangeFeedHttpClient.cs`, `Grpc/PgChangeFeedGrpcClient.cs`, `Grpc/PgChangeFeedAdministrationClient.cs` (`ReadChanges`, `ADR-0138`), `Sse/PgChangeFeedSseClient.cs`, `Nats/PgChangeFeedNatsStreamClient.cs`) | update | optionaler Parameter je Fläche, additiv; die Signatur ohne Parameter bleibt aufrufbar. Gelandet: `target` ist in C# der **letzte** Parameter (hinter dem `CancellationToken`), damit ein Aufruf mit positionalem Token unverändert übersetzt; der `ReadChanges`-Aufruf des Administration-Clients reicht den generierten Request unverändert durch (das Feld `Target` kommt aus dem Stub, nur Docstring geändert). NATS: `BuildTargetSubject`/`BuildSourceTargetsSubject` bauen `cdc.route.<source>.<ziel>` bzw. `cdc.route.<source>.>` mit derselben Token-Prüfung wie `BuildSubject`. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/` (neu: `Http/PgChangeFeedHttpClientTargetTests.cs`, `Grpc/PgChangeFeedGrpcClientTargetTests.cs`, `Sse/PgChangeFeedSseClientTargetTests.cs`; erweitert: `Nats/SubjectTests.cs`), `sdks/csharp/README.md` | update | Tests je Fläche (Draht, Regression ohne Parameter, derselbe Eingabesatz `null`/`""`/`eu`/`a&b=c`); README (Englisch). |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/` (`http/`, `grpc/`, `sse/`, `nats/`) und `README.md` | update | dasselbe für Kotlin (`target` als letzter Parameter mit Default `null`); Testquellen unter `src/test` (neu: `http/…TargetTest.kt`, `grpc/…TargetTest.kt`, `sse/…TargetTest.kt`; erweitert: `nats/…SubjectTest.kt`). Die Prozent-Kodierung des HTTP-Clients ist als `internal fun percentEncode` aus der Klasse herausgelöst und vom SSE-Client mitbenutzt. |
| `sdks/python/pgchangefeed/src/pgchangefeed/` (`http_client.py`, `grpc_client.py`, `administration_client.py` (`ReadChanges`, `ADR-0138`), `sse_client.py`, `nats_stream_client.py`) und `sdks/python/README.md` | update | dasselbe für Python (`target` als letzter Parameter); Tests unter `sdks/python/pgchangefeed/tests/` (erweitert in fünf Dateien). Python-NATS: `stream_changes(timeout, target)` abonniert mit gesetztem Ziel `cdc.route.<source_id>.<ziel>`; leeres Ziel behält den Namensraum der Quelle, ein ungültiges Zeichen im Ziel endet mit `ValueError` beim ersten `next()` des Iterators (`stream_changes` ist ein Generator), vor dem Verbindungsaufbau. C# und Kotlin kennen kein leeres Ziel am Subjekt-Bau (`BuildTargetSubject("")` wirft); die Differenz steht in den drei READMEs. |
| `examples/<fläche>-client/` (Go), `examples/csharp/<fläche>-client/`, `examples/kotlin/<fläche>-client/` | update | Flag `-target`/`--target` an `http-client` (Verb `changes`), `grpc-client` (Verben `stream` und `read-changes`), `sse-client`; `nats-stream-client` nimmt `-source` zusammen mit `-target` und abonniert `cdc.route.<source_id>.<ziel>`. Ein Beispiel-Client des SSE- und des NATS-Wegs trug bisher weder `schema` noch `table` (anders als die Beschreibung der Welle annimmt); sie bekommen nur `target`. Neue Hilfsdateien für die Bindung Flag → Anfrage: `examples/nats-stream-client/subject.go` (+ Test), `examples/csharp/nats-stream-client/SubscribeSubject.cs` (+ Test), `examples/csharp/grpc-client/StreamRequest.cs`, `examples/kotlin/nats-stream-client/…/SubscribeSubject.kt` (+ Test), `examples/kotlin/grpc-client/…/StreamRequest.kt`. |
| `examples/README.md`, `docs/user/benutzerhandbuch.md` (SDK-Passagen, Version 1.86, Änderungshistorie) | update | Nachzug (Liefer-Punkt C). |

**Ansatz:** Das Formvorbild je Fläche ist der `schema`/`table`-Filter
([`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md)); der
Implementer liest ihn je Sprache und setzt `target` an dieselben Stellen. Die drei
Sprachen bauen die gRPC-Stubs im eigenen Bau aus der `.proto` (`--build-context
proto=proto`), das neue Feld stammt aus `slice-routing-lesewege`. Die Stubs tragen
das Feld ohne `.proto`-Änderung (`StreamChangesRequest.target = 3`,
`ReadChangesRequest.target = 7`); die gebauten Tests setzen und lesen es.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die SDK- und
Beispiel-Flächen filtern nach `schema`/`table`"; Parent ist `30fd6cb5`; der
Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und Nichtgefundenes ein):**

```suchlauf
30fd6cb5 9 -n -i -E 'filter' -- sdks/csharp/README.md sdks/python/README.md sdks/kotlin/pgchangefeed-kotlin/README.md examples/README.md
30fd6cb5 59 -n -E 'StreamChangesAsync|stream_changes|streamChanges' -- sdks ':!*Test*' ':!*test*' ':!dist'
30fd6cb5 23 -n -E 'cdc\.stream\.' -- sdks ':!*Test*' ':!*test*' ':!dist'
30fd6cb5 54 -n -E '\-schema' -- examples
30fd6cb5 0 -n -E 'ADR-|LH-FA|LH-QA|SPEC-|ARC-' -- sdks ':!dist'
30fd6cb5 3 -n -E '\-target' -- examples
30fd6cb5 55 -n -E 'ReadChangesAsync|read_changes|readChanges' -- sdks ':!*Test*' ':!*test*' ':!dist'
30fd6cb5 0 -n -E 'cdc\.route\.' -- sdks ':!*Test*' ':!*test*' ':!dist'
30fd6cb5 4 -n -E 'folgt mit|noch nicht als eigenen|zwei optionale|nehmen .target. nicht' -- docs/user/benutzerhandbuch.md examples/README.md sdks
diff 21 -n -i -E 'filter' -- sdks/csharp/README.md sdks/python/README.md sdks/kotlin/pgchangefeed-kotlin/README.md examples/README.md
diff 67 -n -E 'StreamChangesAsync|stream_changes|streamChanges' -- sdks ':!*Test*' ':!*test*' ':!dist'
diff 23 -n -E 'cdc\.stream\.' -- sdks ':!*Test*' ':!*test*' ':!dist'
diff 58 -n -E '\-schema' -- examples
diff 0 -n -E 'ADR-|LH-FA|LH-QA|SPEC-|ARC-' -- sdks ':!dist'
diff 59 -n -E '\-target' -- examples
diff 55 -n -E 'ReadChangesAsync|read_changes|readChanges' -- sdks ':!*Test*' ':!*test*' ':!dist'
diff 14 -n -E 'cdc\.route\.' -- sdks ':!*Test*' ':!*test*' ':!dist'
diff 0 -n -E 'folgt mit|noch nicht als eigenen|zwei optionale|nehmen .target. nicht' -- docs/user/benutzerhandbuch.md examples/README.md sdks
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aussagen über den Filter in README und Beispiel-Dokumentation | Zeile 1: 9 Zeilen (darunter `sdks/csharp/README.md:251` „The gRPC stream can be filtered by schema/table …; the SSE stream c…" — die Aussage ist vollständig zu lesen: ob sie nach dem Stand des Servers noch stimmt, ist am Start zu prüfen) | jede Zeile lesen: ergänzt um `target` oder als veraltet berichtigt (Berichtigung desselben Satzes, kein neuer Vorgang). **Befund am Diff (21 Zeilen, `diff`-Zeile 1; die letzten zwei sind die Sätze der Fixrunde zum leeren Ziel im C#- und Kotlin-NATS-Absatz, „no filter“):** alle neun Parent-Zeilen sind berührt — die drei Sätze in `examples/README.md` (gRPC-Beispiele je Sprache) nennen `-target`/`--target`; die beiden Sätze „takes two optional parameters" (C#, Kotlin) heißen jetzt „three"; die Python-Tabellenzeile nennt `target=None`; die drei Schlusssätze „the SSE stream cannot yet" (C#, Kotlin) bzw. „gRPC and SSE streams cannot be filtered by table" (Python) sagen jetzt, dass SSE (und in Python das NATS-Abonnement) nach dem Ziel, nicht nach `schema`/`table` filtert. Die Aussage „cannot yet" stimmte für den SDK-Client weiter (der SDK-SSE-Client hatte weder `schema` noch `table`), der Server kann beides seit `ADR-0133`; der Satz sagt es jetzt getrennt. Nichtgefunden: kein weiterer Satz mit „Filter" in den vier Dateien, den dieser Lauf nicht ändert |
| Stellen, die Stream-Filter an die Server-Anfrage geben | Zeile 2: 59 Nicht-Test-Zeilen in den drei Packages (Methoden `StreamChangesAsync`, `stream_changes`, `streamChanges` samt Docstrings) | jede Stelle trägt `target` oder die Auslassung ist begründet. **Befund am Diff (67, +8):** gRPC (`StreamChangesAsync`/`streamChanges`/`stream_changes` der gRPC-Clients) und SSE tragen `target`; die +8 Zeilen sind README-Zeilen (C# +1, Kotlin +1, Python +4) und je eine Docstring-Zeile in `nats_stream_client.py` und `sse_client.py`. Die Methoden der drei NATS-Clients tragen **kein** `target` in `StreamChangesAsync(subject)`/`streamChanges(subject)` — der Aufrufer baut das Subjekt über `BuildTargetSubject`/`buildTargetSubject` (Auslassung begründet: der Subjekt-Parameter ist der Wähler, ein zweiter Parameter wäre eine zweite Wahl desselben Gegenstands); in Python hat `stream_changes` kein Subjekt-Argument, dort trägt es `target`. **Lücke des Suchmusters:** die Muster trifft `ReadChanges*`/`read_changes`/`readChanges` nicht; die `diff`-Zeile `ReadChanges…` (55, gleich dem Parent) zählt diese Stellen — `PgChangeFeedHttpClient.ReadChangesAsync`/`readChanges`/`read_changes` tragen `target`, die Administration-Clients reichen das Request-Feld durch (nur Docstring geändert), `models.py`/`__init__.py` tragen den Namen nur als Typ- bzw. Exportnamen. Die Zahl bleibt gleich, weil die Änderungen am Parameter, nicht am Namen liegen |
| Stellen mit dem NATS-Subjekt-Schema im SDK | Zeile 3: 23 Zeilen | das Zusatz-Subjekt kommt je Sprache an die Stelle des Subjekt-Baus. **Befund am Diff (23, gleich dem Parent):** die 23 `cdc.stream.`-Stellen bleiben unverändert (Tabellen-Subjekt-Bau und Docstrings); das Zusatz-Subjekt steht an 14 neuen Stellen (`diff`-Zeile `cdc\.route\.`: Subjekt-Bau und Docstrings der drei NATS-Clients, drei READMEs). Nichtgefunden: kein weiterer Subjekt-Bau in den Packages |
| Flags der Beispiele | Zeile 4: 54 Zeilen `-schema` unter `examples/` | jede Stelle mit dem Paar `-schema`/`-table` bekommt das Ziel-Flag, die Hilfetexte eingeschlossen. **Befund am Diff (58 `-schema`, 59 `-target`; Parent 54 bzw. 3):** von den 54 Parent-Zeilen tragen nur die `grpc-client`-Stellen den Stream-Filter (Flag, Hilfetext, Go-Test-Kommentar, `examples/README.md`, Handbuch) — dort steht jetzt `-target`; die übrigen `-schema`-Zeilen sind die Pflichtfelder der Tabellen-Verben (`enable-table` u. a.) und das Wecksignal-Beispiel `nats-client` (kein Ziel: Wecksignal ohne Ziel, `SPEC-017`), sie bleiben. Die `-target`-Zahl 3 am Parent sind die `--target`-Kommentarzeilen der drei Dockerfiles unter `examples/`. Nichtgefunden: ein Filter-Paar `-schema`/`-table` am SSE- und am NATS-Vollinhalts-Beispiel — die gab es nicht (die Annahme der Welle-Eröffnung, alle vier Flächen trügen das Paar, traf nur auf `grpc-client` und das `changes`-Verb von `http-client` zu); sie bekommen `-target`, nicht das Paar |
| Interne Kennungen unter `sdks/` | Zeile 5: 0 Zeilen | muss 0 bleiben; `make sdk-public-doc-check` ist der Wächter, der Suchlauf zählt gegen. **Befund am Diff:** 0 (`diff`-Zeile 5); `make sdk-public-doc-check` Exit 0 |
| Handbuch-Passage zu den Clients | `benutzerhandbuch.md:1308` („Go, C# und Kotlin nehmen den Filter über `-schema`/`-table`") liegt im Suchraum von `slice-routing-betriebsdoku` (Zeile 5 dort) | diese Passage zieht dieser Slice. **Befund am Diff:** die Passage (jetzt im Abschnitt „Zugriff über den gRPC-Change-Stream") nennt `-target`/`--target` und die drei Package-Signaturen; die vier Sätze „folgt mit dem Package" / „nehmen `target` nicht als Aufrufparameter entgegen" / „noch nicht als eigenen Aufrufparameter" (Parent 4 Zeilen, `diff`-Zeile 9: 0) sind ersetzt; Beispiel-Absätze der HTTP-, gRPC-, SSE- und NATS-Vollinhalts-Abschnitte und die SDK-Absätze von SSE und NATS nennen den Parameter bzw. die Subjekt-Bauer. Nichtgefunden: kein weiterer Satz zu „zwei optionalen" Parametern im Handbuch |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-betriebsdoku` liegt in `done/`
(Welle §4 Abweichung 1 und Folgepflicht 8: nach dem Slice 7; die Handbuch-Passage der
Clients ist die Zwischenzeit, die dieser Slice schließt) und kein anderer Slice liegt
in `in-progress/` (WIP-Limit 1). Voraussetzung am Start: Docker mit Netzzugang für die
drei `make sdk-pack-*`-Ziele und `make examples-*` (NuGet-, PyPI-, Maven-Bezug).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): mit drei Sprachen, vier
  Flächen und drei Beispiel-Sätzen ist der Slice der umfangreichste Doku-/Client-Zug der
  Welle. Zeigt sich nach der ersten Fixrunde, dass eine Sprache mehr als drei Fixrunden
  braucht, teilt er sich nach dem Muster der Welle `welle-sdk-grpc-administration-flaeche`
  in `slice-routing-sdk-csharp-target`, `…-kotlin-…`, `…-python-…` und einen
  Beispiel-Slice; die Naht ist die Sprach-Wurzel (`sdks/<sprache>/`,
  `examples/<sprache>/`).
- `in-progress` → `open` (blockiert): `make generated-sync` oder ein SDK-Bau scheitert
  am Proto-Feld aus `slice-routing-lesewege` (z. B. Namenskonflikt im generierten Code
  einer Sprache) — der Fund geht an `slice-routing-lesewege` zurück, kein Umgehen im SDK.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert, `sdk-public-doc-check` eingeschlossen), die drei `make sdk-pack-*`-Ziele
und `make examples-csharp`/`make examples-kotlin` real grün, Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Drei Sprachen, ein Randfall, drei Lesarten.** Leerer Wert gegen nicht gesetzt,
  Zeichen im Zielnamen, Kodierung im Query-Parameter
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`, offen, 2×;
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`, verkörpert, 4×). —
  **Ausgang: teils entfallen, teils eingetreten.** *Entfallen:* HTTP, SSE und gRPC lesen die
  Eingabetabelle `null`/`""`/`eu`/`a&b=c` in allen drei Sprachen gleich (Tests je Sprache; Reviewer
  mutierte 13 Stellen in 10 Läufen, der Verifier 10 Stellen, je selbst gefahren, alle rot —
  Review „Mutationen“, Verifikations-Report §3). *Eingetreten:* der Satz der Python-README zum
  NATS-Stream war enger als der Code (Review F-3, vor dem Merge gefunden, in der Fixrunde
  berichtigt) — drittes Auftreten von `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`, Datei
  `evidence/slice-routing-sdk-beispiel-target.md`, siehe den Ausgang des Python-NATS-Randfalls unten.
- **Interne Kennung in öffentlichem Text.** Docstrings, KDoc, XML-Doku und Fehlertexte
  erreichen die Anwender über die Packages; ein Kommentar mit `ADR-0137` ist ein Befund
  des Gates (`BEO-PGC/intern-kennungen-in-ausgelieferten-texten`, offen, 1×). —
  **Ausgang: entfallen.** `make sdk-public-doc-check` Exit 0 („keine interne Kennung unter sdks“,
  Verifikations-Report §1), Suchlauf `diff 0` für `ADR-|LH-FA|LH-QA|SPEC-|ARC-` unter `sdks/`
  (§3); `make kommentar-kennungen DIFF=101e24cd~1` ohne Kandidat (Probe der Form, kein Beleg).
- **Neues Anfrage-Feld, altes Server-Verhalten.** Ein Package mit `target` gegen einen
  Server ohne den Parameter bekommt ungefilterte Changes (gRPC) oder `400` (HTTP, SSE);
  die README sagt es mit dem Ursprung aus `slice-routing-lesewege` §6
  (`BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen`, gestrichen, 2×). —
  **Ausgang: weiter offen, hergeleitet.** Die READMEs sagen die Aussage mit dem Ursprung („has not
  been run against such a release“); ein Lauf gegen ein Release ohne den Parameter ist weder von
  Reviewer noch von Verifier gefahren (Review „Prüfung nach den Schwerpunkten“ (b), Verifikations-Report
  §6). Der Beleg ist ein Lauf, der einen Server ohne den Parameter braucht; er hat keinen eigenen Träger
  (der Realserver-Slice läuft gegen den Stand des Repos).
- **Kein Realserver-Beleg.** Die Unit-Ebene beweist den Draht gegen Fakes, nicht gegen
  den laufenden Server (Abgrenzung §1; der Realserver-Beleg liegt bei
  `slice-routing-sdk-realserver-e2e`). — **Ausgang: weiter offen**, Adresse
  [slice-routing-sdk-realserver-e2e](../in-progress/slice-routing-sdk-realserver-e2e.md); alle Aussagen dieses
  Slice gelten für die Menge „Unit-Ebene, Fake-Transport, Fake-Invoker“ (Verifikations-Report §6).
- **Netzbezug der Bauten.** `make sdk-pack-*` und `make examples-*` brauchen Netz
  (Paketquellen) und Zeit; ein Netzausfall ist Umgebung, kein Befund. — **Ausgang: entfallen.**
  Alle Bauten von Implementer, Reviewer und Verifier endeten mit Exit 0, kein Netzausfall trat auf
  (Verifikations-Report §1).
- **Versionsstand der Packages.** Die Packages tragen die Version der letzten
  Veröffentlichung; eine Quelländerung ohne Versionsänderung ist kein Release — der
  Bericht sagt, dass weder Tag noch Version bewegt wurden. **Hinweis für den
  Release-Entscheid des Auftraggebers** (Review F-2, beabsichtigt): `target` ist in
  C# der neue letzte Parameter (hinter dem `CancellationToken`), in Kotlin ein neuer
  letzter Parameter mit Default `null`. Quellkompatibel: ja; binärkompatibel: nein —
  ein gegen die veröffentlichte Version kompilierter Aufrufer bindet die alte
  Signatur und ist neu zu übersetzen. Ein Java-Aufrufer von
  `PgChangeFeedSseClient.streamChanges` übersetzt nicht mehr ohne Argument (kein
  `@JvmOverloads`; keine Code-Änderung dafür). Der Release trägt einen
  Versionssprung und einen Hinweis im Release. — **Ausgang: weiter offen bis zum
  Release-Entscheid des Auftraggebers** (kein Release in der Welle). Beleg: `git diff --stat` über
  `*.csproj`, `*pyproject.toml`, `*build.gradle.kts` leer und kein Tag nach `101e24cd`
  (Verifikations-Report §2, F-2).
- **`schema`/`table` am SSE-Client der Packages.** Der Server filtert den SSE-Stream
  seit `ADR-0133` nach `schema`/`table`; die SSE-Clients der drei Packages und das
  SSE-Beispiel setzen beides nicht, `target` ist ihr einziger Filter. Die Nachlieferung
  ist ein eigener Slice (die öffentliche Doku nennt keinen Slice-Namen). —
  **Ausgang: weiter offen**, Adresse
  [slice-sdk-sse-client-schema-table-filter](../open/slice-sdk-sse-client-schema-table-filter.md); der
  Plan dort trägt den Gegenstand (`git grep -n -E 'SSE-Client|schema|table' -- docs/plan/planning/open/slice-sdk-sse-client-schema-table-filter.md`
  trifft §1 bis §3; seine Zusage „drei optionale“ steht erst nach der Lieferung an den SSE-Absätzen).
- **Python-NATS-Randfall `target=""`** (Review F-3):
  `drei-sprachen-kopie-divergiert-am-randfall` — Python liest ein leeres Ziel als
  „kein Ziel" (`stream_changes(target="")` abonniert den Namensraum der Quelle), C# und
  Kotlin weisen ein leeres Ziel am Subjekt-Bau ab. Die Differenz folgt dem
  Parameter-Schnitt (Python hat kein Subjekt-Argument) und steht in den READMEs und im
  Docstring. — **Ausgang: eingetreten und eingeordnet.** Das Verhalten ist eine gewollte
  API-Form-Differenz, kein Divergenz-Fund; gezählt wird der Satz der Python-README, der dem Wortlaut der
  Schwester-READMEs folgte und enger als der Code war (gleiche Form wie der zweite Beleg der Klasse):
  drittes Auftreten, Zähler **3×**, Ausgang beim Lese-Schritt der Closure von
  [welle-routing](../welle-routing.md) (Begründung im `state.md` der Beobachtung).
- **Iterator-Test für ein Ziel aus Leerraum (Verifier V-2, LOW).** Die README-Aussage zu Python-NATS
  („nur Leerraum“ wirft beim ersten `next()`) ist für diese Zeichenklasse am Subjekt-Bau getestet, am
  Iterator nur mit `"a.b"`. — **Ausgang: hingenommen.** Beide Eingaben laufen im Code durch denselben
  Pfad (`_target_subject`, gelesen vom Verifier, Verifikations-Report §5 F-3 und §8); der Fehlerzeitpunkt ist
  durch den Iterator-Test mit `"a.b"` gebunden. Kein Folgezug.
- **Kotlin-Testzahl im grünen Lauf (Verifier V-1, INFO).** Der grüne Gradle-Lauf druckt keine Testzahl; die
  Zahl 110 ist aus den roten Mutationsläufen **abgeleitet** (`110 tests completed, 2 failed` und
  `110 tests completed, 1 failed`), nicht aus einem grünen Lauf gemessen. — **Ausgang: hingenommen**, Eigenschaft
  des Sensors; die Zahlen 137 (C#) und 144 (Python) sind im Bau ohne Cache gedruckt.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei SDK-Packages (C#, Kotlin, Python) und die Beispiele in Go, C#
  und Kotlin tragen `target` an HTTP-Lesezugriff, gRPC-Stream, `ReadChanges` des Administration-Clients
  und SSE sowie das NATS-Zusatz-Subjekt; ohne Angabe bleibt der Aufruf unverändert (Tests je Sprache mit
  derselben Eingabetabelle `null`/`""`/`eu`/`a&b=c`). Belegt auf der Menge „Unit-Ebene, Fake-Transport,
  Fake-Invoker“, nicht am laufenden Server. **Ursprung der Zahlen (§3.12):** Testzahlen C# 137 und Python
  144 hat der Verifier im Bau ohne Cache gedruckt gesehen; Kotlin 110 ist **abgeleitet** aus den
  roten Mutationsläufen des Verifiers (der grüne Lauf druckt keine Zahl); die C#-Beispiele (54, 12, 9, 52, 20)
  sind vom Verifier im Bau ohne Cache gedruckt. Mutationen: der Reviewer fuhr 13 Stellen in 10 Läufen selbst, der
  Verifier 10 Mutationen selbst (Go 4, Python 2, C# 2, Kotlin 2), alle rot; die Mutationen des Implementers sind
  **übernommen**, nicht nachgefahren. Die Wahrheitsprobe kam von zwei unabhängigen Kontexten, die den Bau
  selbst ausführten.
- **Was ging anders als geplant:** (1) Die Annahme der Welle-Eröffnung, alle vier Beispiel-Flächen trügen das Paar
  `-schema`/`-table`, traf nur `grpc-client` und das `changes`-Verb von `http-client` zu; SSE- und NATS-Beispiel
  bekamen nur `target` (Plan §3, Suchlauf-Tabelle). (2) Eine Fixrunde (`9e8bc1ce`) zu F-1 (MEDIUM: der Aufschub
  „offener Folge-Schritt“ für `schema`/`table` am SSE-Client trug keine Adresse; die Aussage wurde auf den
  Ist-Zustand berichtigt, die Adresse steht in §6), F-3 (Python-README zum leeren Ziel) und F-4
  (Fehlerzeitpunkt des Python-Generators); sie änderte nur READMEs, einen Docstring, das Handbuch (1.86) und den Plan.
  Ein Re-Review gab es nicht (engere Fassung von `BEO-PGC/fixrunde-ohne-reviewer-lesung`, Verifier hat
  ausgeführt). (3) Die drei `make sdk-pack-*` und beide `make examples-*` kamen bei Reviewer und Verifier aus dem
  Docker-Schicht-Cache und druckten keine Testzeile.
- **Steering-Loop-Eintrag:** geschärfte Handlung, kein neuer Sensor, nicht verkörpert (gezählt, 1×): **Ein
  Docker-Cache kann Testläufe still überspringen — ein `make sdk-pack-*`- oder `make examples-*`-Lauf, der
  aus dem Cache kommt und keine Testzeile druckt, belegt keine Testausführung.** Die Lese-Handlung der
  Rollen (Reviewer, Verifier, Planner beim Abhaken von „Tests grün“): die gedruckte Testzeile eines Baus
  ohne Cache an der Test-Stufe (`docker build --no-cache`, Argumente wie das Skript des Makefile-Wegs) oder eine
  Mutation, die die Stufe rot färbt, ist der Beleg; „Exit 0“ allein ist es nicht. Benannte Lücke: kein Sensor
  erzwingt den Bau ohne Cache, und der grüne Gradle-Lauf druckt keine Testzahl (V-1). Zweitens, als Beobachtung
  (nicht als Regel): eine Mutationsreihe nennt absolute Pfade und prüft `git status` nach jedem Lauf — der
  Reviewer mutierte nach einem fehlgeschlagenen `cd` drei Repo-Dateien (zurückgenommen, selbst berichtet).
  Beide ohne `liegt in`, weil in diesem Slice nichts verkörpert wurde.
- **Beobachtungs-Register (`../observations/`):** drei Dateien, zwei Vermerke, ein Deckel-Vermerk.
  - **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`** (offen): `evidence/slice-routing-sdk-beispiel-target.md`,
    Zähler **2×** → **3×** (`ls evidence | wc -l`), **Schwelle erreicht**. Einordnung nach der Zählregel (gezählt
    wird der Satz, nicht das Verhalten): das Verhalten ist eine gewollte API-Form-Differenz (Plan §3 benennt sie),
    der Satz der Python-README folgte dem Wortlaut der Schwester-READMEs und war enger als der Code — dieselbe Form
    wie der zweite Beleg; Begründung im `state.md`. Der Ausgang gehört zum Lese-Schritt der Closure von
    [welle-routing](../welle-routing.md).
  - **`BEO-PGC/docker-cache-ueberspringt-tests-still`** (neu, offen, **1×**): Beobachtung, Zustand und
    `evidence/slice-routing-sdk-beispiel-target.md` angelegt. Eigenes Verzeichnis statt Beleg bei
    `BEO-PGC/test-runner-stiller-ausschluss` (Skript-Filter) und `BEO-PGC/test-methode-lauft-still-nicht`
    (Signatur-Form): hier ist der Test vorhanden und erfasst, die Stufe läuft nicht; Abgrenzung in der
    `observation.md`.
  - **`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`** (verkörpert, 8× → 9×, unter dem Deckel von 10×):
    `evidence/slice-routing-sdk-beispiel-target.md` — Mutationsprobe des Reviewers an drei Repo-Dateien nach einem
    fehlgeschlagenen `cd`, selbst bemerkt und zurückgenommen; im `state.md` als neunter Beleg vermerkt, kein
    Beleg für das erste Neubewertungs-Kriterium.
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×): **siebter Gegenbeleg, keine Datei** (Verifier hat die
    Fixrunde ausgeführt); im `state.md` vermerkt, der Ausgang bleibt beim Lese-Schritt der Closure von
    [welle-routing](../welle-routing.md).
  - **Kein Eintrag, nach Zählregel entschieden:** Review F-1 (MEDIUM, Aufschub ohne Adresse im Handbuch) — die
    Regel „Ein Aufschub ohne Adresse ist keiner“ steht in `.claude/commands/implement-slice.md` Schritt 17, der
    Reviewer fand die Verletzung vor dem Merge mit der dortigen Probe, die Fixrunde schloss sie; es entsteht keine
    neue Datei, weil keine bestehende Beobachtung den Fall „Adresse fehlt ganz“ als Klasse führt und er mit F-1
    einmal vorlag (unter der Schwelle, wird bei einem weiteren Auftreten neu angelegt). Review F-2 (Binärkompatibilität),
    F-4, F-5 sind Träger-Nachzüge im Slice. `BEO-PGC/intern-kennungen-in-ausgelieferten-texten`: kein Auftreten
    (Gate Exit 0); `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`: kein Fund in den READMEs.
- **Folge-Slices:** [slice-sdk-sse-client-schema-table-filter](../open/slice-sdk-sse-client-schema-table-filter.md)
  (SSE-Client der drei Packages und SSE-Beispiele: `schema`/`table`; ist eine Datei in `open/`, ohne Welle) und
  [slice-routing-sdk-realserver-e2e](../in-progress/slice-routing-sdk-realserver-e2e.md) (Realserver-Beleg, Welle §4
  Abweichung 4).
- **Risiken aus §6:** Randfall in drei Sprachen **teils entfallen, teils eingetreten** (Satz der Python-README, 3×);
  Interne Kennung **entfallen**; Altserver **weiter offen, hergeleitet**; Realserver **weiter offen**
  (`slice-routing-sdk-realserver-e2e`); Netzbezug **entfallen**; Versionsstand **weiter offen** bis zum
  Release-Entscheid des Auftraggebers (Binärinkompatibilität C#/Kotlin benannt); `schema`/`table` am SSE-Client
  **weiter offen** (`slice-sdk-sse-client-schema-table-filter`); Python-NATS `target=""` **eingetreten und
  eingeordnet**; V-2 und V-1 **hingenommen**.
- **Drei Paarungen:** der Slice gehört zu [welle-routing](../welle-routing.md) (offen) — die Prüfung läuft bei deren
  Closure; die DoD-Zeile bleibt deshalb `[ ]`. (a) Anker: nichts verkörpert; (b) Folge-Slice: beide genannten Pläne liegen
  als Datei in `open/`; (c) Register: die genannten Kennungen existieren als Verzeichnis, jede trägt ein nicht leeres
  `evidence/`.
- **Validator (Modul 8):** entfällt — der Nutzer-Bedarf ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst
  durch den Wellen-Beleg am laufenden Server validierbar (`slice-routing-sdk-realserver-e2e`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `sdks/` (drei Sprach-Wurzeln), `examples/` und `docs/user/` — eine Sub-Area;
die Sprach-Wurzeln sind die Nahtstellen einer möglichen Teilung (§4).

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 2×),
`BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen, 1×),
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert, 3×),
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (verkörpert, 4×),
`BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen` (gestrichen, 2×),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×). Ein weiteres
Auftreten von `drei-sprachen-kopie-divergiert-am-randfall` erreicht 3× und braucht
dann einen Ausgang.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

