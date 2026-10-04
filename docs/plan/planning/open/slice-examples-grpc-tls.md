# Slice examples-grpc-tls: die gRPC-Beispiele in Go, C# und Kotlin verbinden wahlweise über TLS

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (gemessen: `ls docs/plan/planning/*.md` am Planungsstand nennt nur
`README.md`; die Roadmap führt unter *Offene Wellen* keine Welle), siehe
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
(Modul 6).

**Abhängigkeit:** keine harte. Der Server-TLS-Weg liegt mit
[`slice-tls-http-grpc-server`](../done/slice-tls-http-grpc-server.md) vor. Die
Beispiele sind eigenständige Programme mit eigenem gRPC-Aufbau und hängen
**nicht** an den SDK-Packages; die TLS-Optionen der Packages sind ein eigener,
noch nicht angelegter Folge-Slice (`sdk-tls-optionen`, wartet auf eine
Folge-Anforderung, §1) und keine Voraussetzung.

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (gRPC-Stream, den die
Beispiele zeigen),
[`LH-FA-SST-011`](../../../../spec/lastenheft.md) (Scope: der Server spricht TLS,
wenn konfiguriert),
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) und
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Beispiele als Dokumentation mit Bau-Bindung, Start über `make`),
[`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) (Festlegung 1),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).

**Berührte Spec-Stellen:** [`SPEC-034`](../../../../spec/pflichtenheft.md)
(TLS-Konfiguration des Servers; der Slice ändert sie nicht, er liest sie als
Gegenseite). Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, bei der Closure von
[`slice-tls-http-grpc-server`](../done/slice-tls-http-grpc-server.md) angelegt
(Review F-5 und Verifikation nennen den Beispiel-Nachzug als Adresse; kein
Architect: keine neue Entscheidung, die Beispiele folgen
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)). **Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass:** Ein Betreiber, der den Feed-Container mit TLS betreibt
([`LH-FA-SST-011`](../../../../spec/lastenheft.md)), findet unter `examples/` drei
gRPC-Beispiele (Go, C#, Kotlin), die nur im Klartext verbinden. Gemessen am
Parent `53c65763`: Go `grpc.WithTransportCredentials(insecure.NewCredentials())`
(`examples/grpc-client/main.go`), C# der Switch
`Http2UnencryptedSupport` (`examples/csharp/grpc-client/Program.cs`), Kotlin
`usePlaintext()` (`examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Main.kt`
und im Test `DispatcherTest.kt`); §3 Suchlauf Zeile 1 zählt 7 Trefferzeilen. Das
Benutzerhandbuch sagt dazu nur, dass die gRPC-Beispiele nicht mit einem
TLS-Server verbinden (§3 Suchlauf Zeile 2: 2 Trefferzeilen mit „gRPC-Beispiele“,
eine davon die Stelle des Satzes). Der Schluss „verbinden deshalb nicht“ ist
*hergeleitet*; kein Beispiel ist gegen einen TLS-Server gefahren.

**Ziel:** Jedes der drei gRPC-Beispiele nimmt eine Vertrauensanker-Datei
(Zertifikat als PEM) entgegen, verbindet dann über TLS und holt denselben
Change-Stream wie im Klartext; ohne die Angabe bleibt das Verhalten unverändert
(Klartext, Demo-Umgebung); beides ist am laufenden Container belegt, und das
Handbuch nennt die Beispiele nach Messung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **TLS-Optionen der SDK-Packages** (C#, Kotlin, Python) — `sdk-tls-optionen`,
  **noch keine Datei**: er wartet auf die in
  [`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md) §Konsequenzen verlangte
  Folge-Anforderung im Lastenheft (Zug des Auftraggebers oder des Architects, nicht
  der Planung eines Slice). Die Beispiele bauen ihren gRPC-Kanal selbst; sie
  brauchen die Packages dafür nicht, deshalb ist dies keine Voraussetzung.
- **Die HTTP- und SSE-Beispiele** (Go, C#, Kotlin) — sie tragen keinen
  TLS-spezifischen Code (am Slice `tls-http-grpc-server` gemessen: `git grep -i
  -E 'tls|https://|ssl'` in `examples/http-client`, `examples/sse-client`: 0
  Treffer) und nehmen eine Basis-URL entgegen; ob `https://` dort mit dem
  Vertrauensanker des Hosts funktioniert, ist **nicht gemessen** und kein
  Gegenstand dieses Slice. Der Implementer nennt im Bericht, ob ein Aufruf gegen
  den TLS-Container gefahren ist; ein Ergebnis wird im Handbuch nur als Messung
  genannt, nicht als Erwartung.
- **Client-Zertifikate (mTLS), eine eigene CA-Verwaltung, Ablauf- und
  Namensprüfung durch die Beispiele** — der Server fordert kein
  Client-Zertifikat ([`LH-FA-SST-011`](../../../../spec/lastenheft.md)); die
  Standard-Prüfung der Bibliothek (Kette und Name) bleibt, kein
  Überspringen der Prüfung (`InsecureSkipVerify`-Äquivalente sind ausgeschlossen).
- **Die Demo-Umgebung unter `examples/compose.yaml`** — sie läuft im Klartext
  und bleibt so ([`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
  Festlegung 3: die Werte gelten nur im Netzwerk `cdc-examples`). Der
  Realserver-Beleg dieses Slice nutzt eine Compose-Override-Datei im
  Temp-Verzeichnis, kein committetes Zertifikat.
- **Eine Änderung des Servers oder der Wegwerf-Clients unter `tools/harness/`** —
  der Server-Vertrag liegt vor; die Wegwerf-Clients tragen den TLS-Weg
  (`HARNESS_TLS_CA_FILE`) bereits.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Alle Beleg-Angaben dieser Liste sind **Zusagen** („zu belegen durch …“): der
Planungsstand hat keinen davon gefahren
([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
Instanz B). Die gedruckte Zeile bzw. der Exit-Code steht im Bericht des
Implementers und in §7.

- [ ] **Go-Beispiel mit TLS-Wahl (Liefer-Punkt 1,
      [`LH-FA-SST-011`](../../../../spec/lastenheft.md) Happy Path, Boundary,
      Negative).** `examples/grpc-client` bekommt eine Option für die
      Vertrauensanker-Datei (Flag und Umgebungsvariable nach dem Muster der
      bestehenden Optionen des Programms); gesetzt → `credentials.NewTLS` mit dem
      Anker, ungesetzt → Klartext wie bisher; eine nicht lesbare oder
      PEM-lose Datei endet vor dem Verbindungsaufbau mit einer Fehlermeldung ohne
      Kennung. Unit-Tests (`make test`, netzlos, Loopback) mit einem im Test
      erzeugten Zertifikat gegen einen gRPC-Server im Test: TLS-Weg erreicht den
      Server, Klartext-Weg gegen den TLS-Server scheitert, ungesetzt gegen einen
      Klartext-Server gelingt. **Mutationsprobe (Zusage, an einer Kopie im
      Scratchpad, [`AGENTS.md`](../../../../AGENTS.md) §3.1):** die Wahl
      „gesetzt → TLS“ (Bedingung negiert) und das Lesen des Ankers (falsche
      Datei) färben je einen Test rot; die Eingabeseite ist der Optionswert, nicht
      der Fake
      ([`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)).
- [ ] **C#- und Kotlin-Beispiel mit TLS-Wahl (Liefer-Punkt 2,
      [`LH-FA-SST-011`](../../../../spec/lastenheft.md) Happy Path, Boundary,
      Negative).** Dieselbe Option je Programm: C# `HttpClientHandler` mit dem
      Anker als Vertrauensquelle und ohne den h2c-Switch, wenn gesetzt; Kotlin
      `NettyChannelBuilder` mit `SslContext` aus dem Anker statt `usePlaintext()`,
      wenn gesetzt (die genaue Form der Bibliothek ist **nach Kenntnisstand**, der
      Implementer misst sie am Arbeitsstand und nennt sie im Bericht).
      `make examples-csharp` und `make examples-kotlin` Exit 0 (die Tests laufen im
      Docker-Bau; ein roter Test bricht ihn ab); die Tests je Sprache binden die
      Wahl an den Optionswert wie in Liefer-Punkt 1, mit je einer benannten
      Mutation und gesehener Farbe im Bericht. Der Bau-Beleg ist als Cache-Treffer
      zu kennzeichnen, wo er einer ist
      ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)).
- [ ] **Realserver-Beleg und Handbuch (Liefer-Punkt 3,
      [`LH-FA-SST-011`](../../../../spec/lastenheft.md) Happy Path, Negative).**
      Die drei Beispiele laufen gegen einen Feed-Container mit TLS (Compose-Override
      im Temp-Verzeichnis, Zertifikat mit dem vorhandenen Hilfsprogramm
      `tools/harness/certgen`; kein Schlüssel im Repo, `git status --porcelain` am
      Ende ohne `.pem`) und empfangen eine danach committete Änderung (Tabelle,
      Operation, Wert; `change_id` gegen `cdc.changes` gehalten); ein Aufruf ohne
      Anker gegen denselben Container scheitert sichtbar (Fehlermeldung des
      Programms, kein Absturz ohne Text). Der Beleg ist je Sprache **ein Lauf**,
      keine Runner-Phase von `make test-integration`, sofern der Implementer keine
      anlegt (Entscheidung im Bericht, mit Begründung; ein Runner-Eingriff wäre
      ein vierter Liefer-Punkt und hieße: zurück zur Zerlegung). Das
      Benutzerhandbuch ersetzt den Satz, dass die gRPC-Beispiele nicht mit einem
      TLS-Server verbinden, durch die gemessene Aussage samt Aufrufform;
      `examples/README.md` nennt die Option je Sprache; Handbuch-Historienzeile mit
      der nächsten freien Nummer
      ([`BEO-PGC/handbuch-versionshistorie-uebersprungen`](../observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/observation.md));
      `make handbuch-public-doc-check` Exit 0.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make test`, `make examples-csharp` und `make examples-kotlin` Exit 0,
      `make a-check` Exit 0, `make fmt-check` Exit 0.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Zeilen `make examples-csharp`, `make examples-kotlin`,
      `make example-run-go`: je ein Satz zur TLS-Option); gemeldete Träger
      fremder Dateien mit der Closure nachgezogen (§3 Suchlauf,
      [`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/grpc-client/main.go` und, wo nötig, `format.go`/`main_test.go` | update | Option, Wahl der Transport-Credentials, Test mit im Test erzeugtem Zertifikat (Liefer-Punkt 1) |
| `examples/csharp/grpc-client/Program.cs`, `Config.cs`, `Cli.cs`, `GrpcClient.Tests/` | update | Option, Handler mit Vertrauensanker, h2c-Switch nur im Klartext-Weg (Liefer-Punkt 2) |
| `examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/Main.kt`, `Cli.kt` (Name nach Bestand), `src/test/…` | update | Option, `SslContext` aus dem Anker statt `usePlaintext()`, wenn gesetzt (Liefer-Punkt 2) |
| `docs/user/benutzerhandbuch.md` | update | Satz zu den Beispielen nach Messung, Historienzeile (Liefer-Punkt 3) |
| `examples/README.md`, `harness/README.md` | update | Option je Sprache; Sensors-Zeilen |

- **Zertifikat im Test.** Der Go-Test erzeugt es mit `crypto/x509` zur Laufzeit;
  C# und Kotlin erzeugen es mit der Bibliothek der Sprache im Test oder lesen
  eine zur Laufzeit geschriebene Datei im Temp-Verzeichnis. Kein Zertifikat und
  kein Schlüssel liegen im Repo (ein committeter privater Schlüssel — auch ein
  Test-Schlüssel — wäre ein Befund für den Reviewer).
- **Reihenfolge:** (1) Startmessung (Suchlauf unten); (2) Go; (3) C# und Kotlin;
  (4) Realserver-Läufe gegen einen TLS-Container (`make image` vorher, der Zug
  ändert keine Build-Kontext-Dateien des Servers, aber `make examples-*` baut die
  Beispiele); (5) Handbuch, READMEs, `make gates`.
- **Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
  Eigenschaft: **wie die gRPC-Beispiele ihren Kanal aufbauen** (Symbole
  `insecure.NewCredentials`, `usePlaintext`, `Http2UnencryptedSupport`) und die
  Beschreibung der Beispiele als **Klartext-Clients** (Hedge „Klartext“,
  „gRPC-Beispiele“). Parent ist der Stand `53c65763…` (`git rev-parse HEAD` am
  Planungsstand; nie `HEAD`). Suchraum: Zeile 1 und 3 `examples`, Zeile 2
  `docs/user`, jede Einschränkung mit Grund (die Träger der Aussage liegen dort;
  Reviews, Records und ADRs halten ihren Stand). **Parent-Zeilen gemessen; der
  Implementer setzt `diff`-Zeilen nach der Arbeit und liest jede Trefferzeile.**
  **Erwartung am Endstand** (hergeleitet, nicht gemessen): Zeile 1 wächst nicht
  (die Klartext-Wege bleiben, Tests tragen sie weiter); Zeile 2 sinkt oder bleibt,
  je nach Wortlaut der Handbuch-Fassung; Zeile 3 wird um die Kommentare der drei
  Programme neu gelesen, die sie als Klartext-Clients beschreiben — jeder dieser
  Träger wird auf „Klartext ohne Anker, TLS mit Anker“ nachgezogen.

```suchlauf
53c65763ee4bd7a45840ef74e2b87a6e04d79d5f 7 -n -E 'insecure\.NewCredentials|usePlaintext|Http2UnencryptedSupport' -- examples
53c65763ee4bd7a45840ef74e2b87a6e04d79d5f 2 -n -E 'gRPC-Beispiele' -- docs/user
53c65763ee4bd7a45840ef74e2b87a6e04d79d5f 4 -n -E 'Klartext' -- examples
```

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), **und** der Implementer hat die Startmessung gefahren: die
Suchlauf-Zeilen aus §3 am Arbeitsstand neu gemessen (neue Commit-Kennung als
Parent, `make suchlauf-nachmessen PLAN=` mit diesem Plan Exit 0 nach Anpassung
von Parent und Soll; jede Abweichung mit Ursache im Bericht). Eine harte
Abhängigkeit zu einem anderen offenen Slice besteht nicht; die Reihenfolge zu
[`slice-otlp-metrik-export`](../in-progress/slice-otlp-metrik-export.md) ist frei (andere
Dateien).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Realserver-Beleg
  verlangt eine Runner-Phase in `tools/harness/run-integration-tests.sh` oder einen
  Umbau der Wegwerf-Clients; oder der Slice wächst auf die HTTP- und
  SSE-Beispiele (je Sprache eine Teilung, etwa „Go“ und „C#/Kotlin“).
- `in-progress` → `open` (blockiert — Carveout?): die gRPC-Bibliothek einer Sprache
  bietet keinen Weg, einen Vertrauensanker aus einer Datei zu setzen, ohne die
  Prüfung abzuschalten; kein Carveout, weil dann kein Gate rot ist — das wäre eine
  Frage an den Architect.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit
Exit 0 (die `diff`-Zeilen auf den Endstand gesetzt, jede Trefferzeile gelesen), und
die Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor
oder benannte Spec-Lücke; Kandidat: der Befund, wie jede der drei Bibliotheken
einen Vertrauensanker aus einer Datei nimmt, und ob ein Test den Optionswert
bindet). Ein Gate, das am Stand der Closure rot ist, geht nur mit dokumentiertem
Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen**.

- **Die Wahl „Anker gesetzt → TLS“ ist an den Optionswert nicht gebunden:** ein
  Test mit einem Fake, der unabhängig von der Option antwortet, bleibt grün
  ([`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/observation.md)).
  Gegenmaßnahme: je Sprache eine Mutation der Wahl, gesehene Farbe im Bericht. —
  **Ausgang:** offen bis Closure.
- **Der Bau-Beleg von `make examples-csharp`/`make examples-kotlin` ist ein
  Cache-Treffer** (am Verifikationslauf des Server-Slice beobachtet) und belegt
  den Compile nicht frisch. Gegenmaßnahme: die Tests laufen im Bau; der Bericht
  nennt, ob die Schichten gebaut oder aus dem Cache kamen. — **Ausgang:** offen
  bis Closure.
- **Ein Zertifikat oder Schlüssel landet im Repo oder im Image.** Gegenmaßnahme:
  Zertifikat im Test zur Laufzeit, im Realserver-Lauf in einem Temp-Verzeichnis;
  der Reviewer prüft `git status`. — **Ausgang:** offen bis Closure.
- **Das Handbuch sagt über die Beispiele mehr, als gemessen ist** (HTTP- und
  SSE-Beispiele gegen `https://` sind nicht gefahren). Gegenmaßnahme: nur das
  Gemessene, je Beispiel mit Aufruf
  ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) Instanz B). —
  **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

- **Was hat funktioniert:** (bei der Closure zu füllen)
- **Was ging anders als geplant:** (bei der Closure zu füllen)
- **Steering-Loop-Eintrag:** (bei der Closure zu füllen)
- **Beobachtungs-Register (`../observations/`):** (bei der Closure zu füllen)
- **Folge-Slices:** (bei der Closure zu füllen)
- **Risiken aus §6:** (bei der Closure zu füllen, je genau ein Ausgang)
- **Drei Paarungen:** (bei der Closure zu füllen: Anker · Folge-Slice · Register)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt **eine** Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`examples/`, `docs/user/`, `harness/`) liegen alle in ihr, es
entsteht keine neue Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) für die berührte Sub-Area
durchgegangen; Treffer, Zähler = Zahl der `evidence/`-Dateien am Planungsstand
(Dateien der Eintragsverzeichnisse gezählt, die Zahl der Beispiel-Einträge ist
nicht nachgemessen — **übernommen** aus dem Plan des Server-Slice):
`negativtest-ohne-bindung-an-seine-eingabe` (Schwelle erreicht; Gegenmaßnahme: je
Sprache eine Mutation der Wahl, §2),
`handbuch-versionshistorie-uebersprungen` (3, Historienzeile in Liefer-Punkt 3),
`kommentar-herkunft-als-kette` (höchstens eine Kennung je Kommentar im neuen
Code, [`AGENTS.md`](../../../../AGENTS.md) §3.7). Die Einträge mit mindestens drei
Dateien sind **mit** diesem Slice berührt und tragen hier ihre Gegenmaßnahme.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
