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

- [ ] [`LH-FA-SST-009`](../../../../spec/lastenheft.md) und
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
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (B): die Beispiel-Clients der
      vier Flächen in Go, C# und Kotlin (Parent: Go unter `examples/<fläche>-client/`,
      C# unter `examples/csharp/<fläche>-client/`, Kotlin unter `examples/kotlin/`)
      nehmen das Ziel über ein Flag nach dem Muster von `-schema`/`-table`; die Tests
      des jeweiligen Beispiels belegen Flag → Anfrage; `examples/README.md` nennt es. Die
      Beispiele bleiben Doku mit Bau-Bindung, kein Lauf-Beleg. *Zu belegen durch:*
      `make test` (Go-Beispiele), `make examples-csharp`, `make examples-kotlin`
      (Docker-only, braucht Netz).
- [ ] [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) und Nachzug
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
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9;
      `sdk-public-doc-check` ist Teil von `make gates`).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-sdk-beispiel-target.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: SDK-READMEs, `examples/README.md`, Handbuch (siehe C); kein
      SPEC-/ARC-Eintrag. Die Package-Versionen bleiben (Abgrenzung §1).
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
| `sdks/csharp/PgChangeFeed.Client/` (`Http/PgChangeFeedHttpClient.cs`, `Grpc/PgChangeFeedGrpcClient.cs`, `Grpc/PgChangeFeedAdministrationClient.cs` (`ReadChanges`, `ADR-0138`), `Sse/PgChangeFeedSseClient.cs`, `Nats/PgChangeFeedNatsStreamClient.cs`) | update | optionaler Parameter je Fläche, additiv; die Signatur ohne Parameter bleibt aufrufbar. Gelandet: `target` ist in C# der **letzte** Parameter (hinter dem `CancellationToken`), damit ein Aufruf mit positionalem Token unverändert übersetzt; der `ReadChanges`-Aufruf des Administration-Clients reicht den generierten Request unverändert durch (das Feld `Target` kommt aus dem Stub, nur Docstring geändert). NATS: `BuildTargetSubject`/`BuildSourceTargetsSubject` bauen `cdc.route.<source>.<ziel>` bzw. `cdc.route.<source>.>` mit derselben Token-Prüfung wie `BuildSubject`. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/` (neu: `Http/PgChangeFeedHttpClientTargetTests.cs`, `Grpc/PgChangeFeedGrpcClientTargetTests.cs`, `Sse/PgChangeFeedSseClientTargetTests.cs`; erweitert: `Nats/SubjectTests.cs`), `sdks/csharp/README.md` | update | Tests je Fläche (Draht, Regression ohne Parameter, derselbe Eingabesatz `null`/`""`/`eu`/`a&b=c`); README (Englisch). |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/` (`http/`, `grpc/`, `sse/`, `nats/`) und `README.md` | update | dasselbe für Kotlin (`target` als letzter Parameter mit Default `null`); Testquellen unter `src/test` (neu: `http/…TargetTest.kt`, `grpc/…TargetTest.kt`, `sse/…TargetTest.kt`; erweitert: `nats/…SubjectTest.kt`). Die Prozent-Kodierung des HTTP-Clients ist als `internal fun percentEncode` aus der Klasse herausgelöst und vom SSE-Client mitbenutzt. |
| `sdks/python/pgchangefeed/src/pgchangefeed/` (`http_client.py`, `grpc_client.py`, `administration_client.py` (`ReadChanges`, `ADR-0138`), `sse_client.py`, `nats_stream_client.py`) und `sdks/python/README.md` | update | dasselbe für Python (`target` als letzter Parameter); Tests unter `sdks/python/pgchangefeed/tests/` (erweitert in fünf Dateien). Python-NATS: `stream_changes(timeout, target)` abonniert mit gesetztem Ziel `cdc.route.<source_id>.<ziel>`; leeres Ziel behält den Namensraum der Quelle, ein ungültiges Zeichen im Ziel endet mit `ValueError` vor dem Verbindungsaufbau. |
| `examples/<fläche>-client/` (Go), `examples/csharp/<fläche>-client/`, `examples/kotlin/<fläche>-client/` | update | Flag `-target`/`--target` an `http-client` (Verb `changes`), `grpc-client` (Verben `stream` und `read-changes`), `sse-client`; `nats-stream-client` nimmt `-source` zusammen mit `-target` und abonniert `cdc.route.<source_id>.<ziel>`. Ein Beispiel-Client des SSE- und des NATS-Wegs trug bisher weder `schema` noch `table` (anders als die Beschreibung der Welle annimmt); sie bekommen nur `target`. Neue Hilfsdateien für die Bindung Flag → Anfrage: `examples/nats-stream-client/subject.go` (+ Test), `examples/csharp/nats-stream-client/SubscribeSubject.cs` (+ Test), `examples/csharp/grpc-client/StreamRequest.cs`, `examples/kotlin/nats-stream-client/…/SubscribeSubject.kt` (+ Test), `examples/kotlin/grpc-client/…/StreamRequest.kt`. |
| `examples/README.md`, `docs/user/benutzerhandbuch.md` (SDK-Passagen, Version 1.85, Änderungshistorie) | update | Nachzug (Liefer-Punkt C). |

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
diff 19 -n -i -E 'filter' -- sdks/csharp/README.md sdks/python/README.md sdks/kotlin/pgchangefeed-kotlin/README.md examples/README.md
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
| Aussagen über den Filter in README und Beispiel-Dokumentation | Zeile 1: 9 Zeilen (darunter `sdks/csharp/README.md:251` „The gRPC stream can be filtered by schema/table …; the SSE stream c…" — die Aussage ist vollständig zu lesen: ob sie nach dem Stand des Servers noch stimmt, ist am Start zu prüfen) | jede Zeile lesen: ergänzt um `target` oder als veraltet berichtigt (Berichtigung desselben Satzes, kein neuer Vorgang). **Befund am Diff (19 Zeilen, `diff`-Zeile 1):** alle neun Parent-Zeilen sind berührt — die drei Sätze in `examples/README.md` (gRPC-Beispiele je Sprache) nennen `-target`/`--target`; die beiden Sätze „takes two optional parameters" (C#, Kotlin) heißen jetzt „three"; die Python-Tabellenzeile nennt `target=None`; die drei Schlusssätze „the SSE stream cannot yet" (C#, Kotlin) bzw. „gRPC and SSE streams cannot be filtered by table" (Python) sagen jetzt, dass SSE (und in Python das NATS-Abonnement) nach dem Ziel, nicht nach `schema`/`table` filtert. Die Aussage „cannot yet" stimmte für den SDK-Client weiter (der SDK-SSE-Client hatte weder `schema` noch `table`), der Server kann beides seit `ADR-0133`; der Satz sagt es jetzt getrennt. Nichtgefunden: kein weiterer Satz mit „Filter" in den vier Dateien, den dieser Lauf nicht ändert |
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
  **Ausgang:** bei der Closure einzutragen (dieselbe Eingabetabelle je Sprache).
- **Interne Kennung in öffentlichem Text.** Docstrings, KDoc, XML-Doku und Fehlertexte
  erreichen die Anwender über die Packages; ein Kommentar mit `ADR-0137` ist ein Befund
  des Gates (`BEO-PGC/intern-kennungen-in-ausgelieferten-texten`, offen, 1×). —
  **Ausgang:** bei der Closure einzutragen (`make sdk-public-doc-check` Exit 0).
- **Neues Anfrage-Feld, altes Server-Verhalten.** Ein Package mit `target` gegen einen
  Server ohne den Parameter bekommt ungefilterte Changes (gRPC) oder `400` (HTTP, SSE);
  die README sagt es mit dem Ursprung aus `slice-routing-lesewege` §6
  (`BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen`, gestrichen, 2×). —
  **Ausgang:** bei der Closure einzutragen.
- **Kein Realserver-Beleg.** Die Unit-Ebene beweist den Draht gegen Fakes, nicht gegen
  den laufenden Server (Abgrenzung §1; der Realserver-Beleg liegt bei
  `slice-routing-sdk-realserver-e2e`). — **Ausgang:** bei der Closure einzutragen
  (Aussage im Bericht auf die Menge „Unit-Ebene, Fake-Transport" beschränkt).
- **Netzbezug der Bauten.** `make sdk-pack-*` und `make examples-*` brauchen Netz
  (Paketquellen) und Zeit; ein Netzausfall ist Umgebung, kein Befund. — **Ausgang:** bei
  der Closure einzutragen.
- **Versionsstand der Packages.** Die Packages tragen die Version der letzten
  Veröffentlichung; eine Quelländerung ohne Versionsänderung ist kein Release — der
  Bericht sagt, dass weder Tag noch Version bewegt wurden. — **Ausgang:** bei der
  Closure einzutragen.

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

