# Slice transformationen-start-reihenfolge: Startreihenfolge — offene Anträge sind verarbeitet, bevor der Stream die erste Transaktion assembliert

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Negative:
Abhilfe), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 (Bedingung (c) des Abhilfe-Akzeptanzkriteriums) und
§Re-Evaluierungs-Trigger (die Abhilfe wirkt nicht ohne Eingriff in die
Startreihenfolge),
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue, Live-Reload),
[`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) (Persist-before-ACK —
nicht berührt).

**Berührte Spec-Stellen:** [`ARC-007`](../../../../spec/architecture.md)
(Composition Root) — gelesen, nicht geändert;
[`SPEC-019`](../../../../spec/pflichtenheft.md) (Verarbeitung offener Anträge)
— gelesen.

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Antrag, der beim Prozessstart `pending` steht, ist verarbeitet,
**bevor** `stream.Run` die erste Transaktion assembliert — deterministisch,
nicht durch die Planung der Goroutinen. Die Composition Root führt vor
`stream.Run` einen synchronen Durchlauf von `processAdministrationRequests`
aus; die Administrations-Goroutine läuft danach unverändert weiter. Grund:
`runAdministration` startet in `Run` per `go func()` **vor** `stream.Run`,
jeder ihrer Durchläufe beginnt mit `processAdministrationRequests`, zwischen
dem Start der Goroutine und dem Aufruf von `stream.Run` liegt keine
Synchronisation (Stand `c82d3333`, gemessen am 2026-09-25 mit `git grep -n -E
'runAdministration|stream\.Run|go func' c82d3333 -- internal/bootstrap/wiring.go`:
`go func()` mit `runAdministration` in den Zeilen 834 und 836, `stream.Run` in
Zeile 1009; zwischen beiden Zeilen stehen fünf weitere `go func()`-Starts, in den
Zeilen 865, 922, 950, 983 und 1004 — ob sie die Reihenfolge berühren, ist nicht
gelesen; der Implementer liest `Run` am Start erneut) — die Reihenfolge ist Zufall;
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
führt die Abhilfe im gescheiterten Prozess ausdrücklich als „erwartet, nicht am
Code belegt“.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Abhilfe-Beleg am realen System** — `e2e-abhilfe`; ein E2E-Lauf kann das
  Fehlen einer Race nicht beweisen (er besteht auch ohne Ordnung, wenn die
  Planung der Goroutinen sie zufällig liefert), der Beleg der Ordnung selbst
  gehört in diesen Slice, die Wirkung am System in `e2e-abhilfe`.
- **Eine Änderung der Verarbeitungs-Semantik einzelner Antragsarten** — jede
  Art wird wie im Dauerbetrieb verarbeitet; der Slice verschiebt nur den ersten
  Durchlauf vor den Stream.
- **Eine Änderung am Capture-Pfad** (`stream.Run`, `CaptureService`, ACK) — der
  Slice berührt die kritische Sektion von
  [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) nicht; er ändert nur,
  was **vor** dem Stream-Start geschieht.
- **Ein Recovery-Weg für andere Schema-Fehler** —
  `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×) bleibt offen;
  dieser Slice und `e2e-abhilfe` decken die Nichtanwendbarkeit einer
  Transformationsregel, keine allgemeine Recovery.

## 2. Definition of Done

- [x] Der Vorlauf steht: `Run` verarbeitet die offenen Anträge (ein synchroner
      Durchlauf von `processAdministrationRequests` über dieselben
      `administrationDeps` wie die Goroutine) **vor** `stream.Run`; die
      Goroutine läuft danach unverändert weiter; ein Lesefehler im Vorlauf
      blockiert den Start nicht (best effort, unverändert zum Dauerbetrieb,
      [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)),
      sondern wird protokolliert; die Verarbeitung eines Antrags im Vorlauf
      trägt die Bindung im `Assembler` nach, bevor der Stream startet. *Zu
      belegen durch:* `make test` (Whitebox in `internal/bootstrap`, Fakes für
      Antrags-Port und Assembler).
- [x] Die Ordnung ist ohne Datenbank prüfbar und an ihre Eingabe gebunden: die
      Sequenz „Vorlauf, dann Stream“ steht an einer Stelle, die ein Test mit
      Fakes aufruft (kleine Extraktion aus `Run`, kein Umbau der Verdrahtung);
      ein Test belegt, dass ein `pending` stehender
      `remove_transformation`-Antrag verarbeitet ist, wenn der
      Stream-Start-Fake aufgerufen wird, und färbt sich rot, wenn die
      Reihenfolge vertauscht wird (Mutation). *Zu belegen durch:* `make test`
      mit Race-Detector und die Mutation im Bericht.
- [x] Der Dauerbetrieb bleibt unverändert: `make coverage-gate` grün (Teil von
      `make gates`, nächster Haken) und ein realer, grüner
      `make test-integration`-Lauf mit den bestehenden Belegen der
      SQL-Administration einschließlich des Live-Reload-Belegs ohne
      `docker restart`. **Reichweite:** der Lauf belegt den Dauerbetrieb der
      Administrations-Goroutine und den Vorlauf mit **leerer** Queue bei jedem
      Prozessstart des Laufs; er belegt **nicht** den Vorlauf mit einem beim
      Start wartenden Antrag: `git grep -n "'pending'" --
      tools/harness/run-integration-tests.sh test/integration` trifft keine Zeile
      (Exit 1, gemessen in der Fixrunde, von der Verifikation und in der
      Closure nachgemessen), also legt kein Runner-Schritt einen `pending`
      stehenden Antrag über einen Start. Diese Wirkung am System trägt
      `slice-transformationen-e2e-abhilfe`. **Beleg-Anker:** der Lauf des
      Verifiers am `HEAD` `cd5ff4f2`, Report
      `docs/reviews/verifikation-slice-transformationen-start-reihenfolge.md` §1
      — `make test-integration` Exit 0, 347 s, die Live-Reload-Zeilen
      (`enable` und `disable` ohne Neustart) im Log gedruckt (**übernommen**
      aus dem Report); der Code- und Skript-Stand ist seither unverändert
      (`git diff --stat cd5ff4f2 HEAD -- internal tools test harness Makefile`
      druckt nichts, gemessen in der Closure). Die Deckung trägt der eigene
      `make gates`-Lauf der Closure (§7).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-transformationen-start-reihenfolge.md`,
      `.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8); die Fixrunde löst das HIGH (F-1) und die beiden
      MEDIUM (F-2, F-3) auf, kein HIGH und kein MEDIUM bleibt offen.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — keine Betreiber-Oberfläche; die Aussage zur
      Startreihenfolge steht, sofern das Handbuch sie trägt (Suchlauf), im
      Handbuch-Abschnitt von `slice-transformationen-betriebsdoku`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` | update | **Geliefert.** Neue Funktion `runStreamAfterAdministrationPass(ctx, deps, startLoop, runStream)`: erst `processAdministrationRequests` (Vorlauf), dann `startLoop` (die Goroutine), dann `runStream`. Die `administrationDeps` werden vor dem Start gebildet (Variable `administration`), der Goroutinen-Start ist die Funktion `startAdministration`; `Run` ruft die Sequenz mit `stream.Run` als viertem Argument. Der Vorlauf steht **vor** dem Goroutinen-Start, damit zu keinem Zeitpunkt zwei Durchläufe dieselbe Queue lesen. Gelesen am Start (Schritt 12): `reconcileBackfillRuns` und der Start des Backfill-Workers stehen in `Run` vor dem Aufruf der Sequenz (Risiko §6, zweiter Punkt), der Kontext des Vorlaufs ist `streamCtx` (abgeleitet aus `ctx`, ein Abbruch durch die WAL-Schwelle beendet den Vorlauf mit). **Fixrunde (F-2, F-3):** der Godoc der Sequenz sagt die Ordnung mit ihrer Bedingung zu — sie gilt für einen Antrag, den der Vorlauf liest und verarbeitet; bei einem Lesefehler von `ListPending` bleibt jeder Antrag `pending` und der Stream startet mit dem bisherigen Regelstand, bei einem gescheiterten `MarkApplied` steht die Wirkung in der Bindung und der Antrag bleibt `pending`, bis die Goroutine ihn erneut verarbeitet (aus `processAdministrationRequests` und `applyAdministrationRequest` gelesen; der Lesefehler-Test bindet „hält den Start nicht an“ und den Rückgabewert, nicht die Folge für den Antrag). Der Godoc nennt außerdem, was den Vorlauf begrenzt (`ctx`, in `Run` `streamCtx`) und dass er keine eigene Frist trägt. **Reichweite der Godoc-Sätze (Verifikation V-2, V-3):** die Herkunft des Kontexts im Aufruf (`runStreamAfterAdministrationPass(streamCtx, …)`) ist aus `Run` gelesen und an keinen Test gebunden — die Mutation „Vorlauf unter `ctx` statt `streamCtx`“ bleibt grün (Verifikation V9b, **übernommen**); der Halte-Test bindet den Kontext innerhalb der Sequenz. Der Satz „die Ordnung gilt für jede Antragsart“ gilt für die Arten mit Bindungswirkung (`enable`, `disable`, `exclude_column`, `include_column`, `set_transformation`, `remove_transformation`; die sieben Arten der Domäne, `model.AdministrationRequestKinds`, gelesen an `applyAdministrationRequest`) wörtlich; für `backfill` heißt „verarbeitet“ angenommen (Run-Zeile `queued`, Antrag `applied`), eine Bindung wird nicht berührt, und ein `backfill`-Antrag einer anderen Quelle bleibt im Vorlauf `pending` (`processAdministrationRequests`, gelesen). |
| `internal/bootstrap/administration_startorder_internal_test.go` | neu (Plan-Nachzug) | Sechs Tests ohne Datenbank: Ordnung an einem `pending` stehenden `remove_transformation`-Antrag (Stream-Fake sieht Vermerk `applied` und Rohform), Ordnung für `enable`, Lesefehler hält den Start nicht an und der Stream-Ausgang bleibt der Rückgabewert, hängende Abfrage hält den Stream-Start an und der Abbruch des Kontexts beendet den Vorlauf, und zwei Tests über die Quelltext-Gestalt von `Run` (`TestRunSourceTextPassesStreamRunOnlyAsArgumentOfTheSequence`: `stream.Run` nur als Argument der Sequenz; `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall`: Aufrufe von Abgleich und Worker-Start vor dem Aufruf der Sequenz). Grund der Quelltext-Tests: `Run` braucht die Datenbank, kein netzloser Lauf erreicht die Aufrufstelle. **Fixrunde (F-4):** die beiden Namen benennen die Messung (Quelltext-Lesung), keine Lauf-Reihenfolge. **Benannte Grenze (F-8, F-9, Verifikation V-2):** die Aufrufstelle in `Run` ist nur über Quelltext-Positionen gebunden; vier Mutationen bleiben im Unit-Lauf grün (Verifikation V9b, V10b, V12, V14, **übernommen**): die Herkunft des Kontexts des Vorlaufs (`streamCtx`), der Rumpf von `startAdministration` (Start der Goroutine, Mutation M10 des Reviews), der Inhalt der Wertegruppe `administration` und ein Aufruf in einer toten Verzweigung (`if false { … }`, in einer Kopie gesehen, `go test -race` grün). Den Goroutinen-Start deckt am laufenden Prozess `make test-integration`: die Mutation „die Goroutine startet nie“ färbt den Integrationslauf rot (Verifikation V10b, Exit 2, `TestE2ETransformationRulesShapeBothImages` und `TestE2ETransformationConflictsFailWithSpecText` FAIL, **übernommen**); für „die Goroutine startet zu früh oder zweimal“ ist das nicht gefahren. Ein verhaltensbasierter Weg (Abgleich, Worker-Start, Vorlauf, Goroutine, Stream als Funktionswerte einer Start-Funktion) ersetzt die Quelltext-Tests nur, wenn `Run` diese Funktionswerte aus seinen Pools und Kontexten selbst baut — ein Umbau der Verdrahtung über die kleine Extraktion hinaus (§4, Rückführung), und die Mutation M10 bliebe dabei ungebunden, weil der Rumpf von `startAdministration` in `Run` entsteht. Nicht geliefert. |
| `internal/adapters/driving/replication/mapper/mapper.go`, `internal/adapters/driving/replication/receive/receive.go` | update (Kommentare, Plan-Nachzug, F-7) | Vier Kommentare (`AddBinding`, `ExcludeColumn`, `RemoveBinding` in `mapper.go`, `Stream.Assembler` in `receive.go`) nennen die Administrations-Verarbeitung (Goroutine und Vorlauf) als Aufrufer der Bindungs-Nachträge; der Vorlauf ruft sie vor dem Stream-Lauf, die Goroutine danach. Je Kommentar höchstens eine Kennung (`make kommentar-kennungen DIFF=ce229b8e` meldet jeden berührten Block mit mehr als einer Kennung; Exit 0, kein Kandidat, im Review und in der Verifikation gemessen). Kosmetik ohne Änderung (Verifikation V-5, INFO): der Kommentar von `RemoveBinding` bricht mit einer kurzen Zeile um, der von `ExcludeColumn` trägt eine überlange Zeile; `gofmt` und `make kommentar-kennungen` melden nichts, der Slice ändert sie nicht. Kein Verhaltens-Diff, kein Eingriff in den Capture-Pfad (§1, dritter Punkt: `stream.Run`, `CaptureService`, ACK bleiben unberührt). Der Struktur-Kommentar am `Assembler` (`mapper.go`, Synchronisation der Bindungstabelle) beschreibt die Nebenläufigkeit im Dauerbetrieb und bleibt wahr: der Vorlauf läuft vor dem Stream-Lauf. |
| `harness/README.md`, `docs/user/benutzerhandbuch.md` | prüfen | **Ergebnis: keine Änderung.** Keine der beiden Dateien beschreibt die Reihenfolge von Antrags-Verarbeitung und Stream-Start (Suchlauf unten, Zeilen 3–4 und 13–14); die Handbuch-Stellen zur Administrations-Goroutine (§4) sagen „verarbeitet offene Anträge … ohne Neustart“ und bleiben wahr. Der Abhilfe-Ablauf gehört zu `slice-transformationen-betriebsdoku` (Adresse). Keine Betreiber-Oberfläche berührt: `git diff --name-only ce229b8e -- internal/bootstrap/ tools/schema/ internal/adapters/driving/` trifft `wiring.go`, die Testdatei und (Fixrunde) die beiden Kommentar-Dateien, keine Umgebungsvariable, keine SQL-Funktion, keinen Endpunkt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Reihenfolge von
Antrags-Verarbeitung und Stream-Start“; beide Stände gemessen; Suchraum der
Zeilen 1–12 und 17–18 der ganze Baum, ausgenommen `docs/reviews/**`, die
Records unter `done/` und `.harness/baseline/**`; die Zeilen 13–16 sind
eingeschränkt, Grund je Zeile in der Tabelle):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Code: Symbole des Startpfads und Beschreibungen der Reihenfolge | Zeilen 1–2, 5–6 und 11–12 im Block unten (Parent `ce229b8e`) | **Gefunden.** Zeilen 1–2 (Symbole `runAdministration`/`stream.Run`/`processAdministrationRequests`): Parent 134, Diff 144; Diff nach Wurzel `internal/bootstrap` 79, `internal/adapters` 23, `docs/plan/planning` 20, `docs/plan/adr` 18, `harness` 4 (je Datei mit `git grep -c` gezählt, summiert). Gelesen: die Kommentare in `wiring.go` an der Goroutine, an `runStreamAfterAdministrationPass`, `runAdministration` und `processAdministrationRequests` (der Kommentar an der Goroutine nennt jetzt den Vorlauf und den Ort der Sequenz, `runStreamAfterAdministrationPass` trägt die Zusage samt Bedingung), die vier Treffer in `harness/sensors/coverage-gate.md` (`runAdministration` als Coverage-Zweig, keine Reihenfolge) und alle Treffer der Zeilen 5–6 und 11–12. Die übrigen Symbol-Treffer (Aufrufe und Tests in `internal`, Pläne, ADRs) sind gezählt, nicht je Zeile gelesen — **Grenze dieser Zeilen**. Zeilen 5–6 („vor“/„nach“ vor `stream.Run`): Parent 7, Diff 8; die sieben des Parent stehen in Plänen, `ADR-0112` und Roadmap, keiner im Code, der Zuwachs ist der Kommentar an der Goroutine (`wiring.go`). Zeilen 11–12 („erste Transaktion“/„bevor der Stream“): Parent 11, Diff 15; neu sind der Godoc der Sequenz und drei Zeilen der Testdatei. **Nichtgefunden:** kein Code-Kommentar des Parent nennt die Reihenfolge von Goroutinen-Start und `stream.Run` (die Treffer der Zeile 5 am Parent liegen alle außerhalb `internal`). | Kommentar an der Goroutine im Ist-Ton, Godoc der Sequenz an den Code gebunden |
| Doku: Beschreibungen der Reihenfolge | Zeilen 3–4, 9–10 und 13–14 im Block unten | **Gefunden.** `spec/pflichtenheft.md:341`, Abschnitt `LH-FA-CFG-007.a` (Zeilen 336–343): „offene Anträge werden beim Start verarbeitet, **bevor** die erste Transaktion der Tabelle assembliert wird … Diese Abfolge muss die Umsetzung liefern; sie ist erst mit dem Beleg am laufenden System eine Tatsache“ — die Abhilfe-Zusage, die dieser Slice umsetzt (Zeilen 9–10, Parent 4, Diff 4: `spec/pflichtenheft.md`, `ADR-0112`, `welle-transformationen`, `slice-transformationen-e2e-abhilfe`). Die Aussage bleibt wahr: der Beleg am System steht bei `e2e-abhilfe` aus; ein Nachzug ist nicht fällig. Handbuch §4 (Zeilen 238 und 271) sagt „Administrations-Goroutine des laufenden Feed-Containers verarbeitet offene Anträge“ — wahr für den Dauerbetrieb, der Vorlauf beim Start steht dort nicht. `SPEC-019` und das Handbuch (Zeilen 76, 1436) nennen die Queue (Zeilen 3–4: Parent 141, Diff 139, ganzer Baum). **Nichtgefunden:** in `docs/user` und `harness` trifft keine Zeile eine Beschreibung der Reihenfolge (Zeilen 13–14, Parent 0, Diff 0; Suchraum dieser Zeilen: nur `docs/user` und `harness`, Grund: sie tragen die Aussage „Nichtgefunden“ für diese beiden Bäume, die Zeilen 1–12 decken den übrigen Baum). | keiner; die Prozedur „Abhilfe nach einer nicht anwendbaren Regel“ und die Beschreibung des Vorlaufs beim Start (samt der Aussage, dass ein langer Antrag den Stream-Start anhält, §6 erster Punkt) trägt `slice-transformationen-betriebsdoku` (Adresse, Frist: dessen Closure) |
| Kommentare außerhalb des Diffs, die die Administrations-Goroutine nennen (F-7) | Zeilen 15–16 im Block unten (Suchraum: `internal` ohne Testdateien, Grund: Kommentare und Godocs der Produktivdateien) | **Gefunden.** Zeilen mit „Administrations-Goroutine“: Parent 22, Diff 20. Gelesen im Kontext: `mapper.go` (`AddBinding`, `ExcludeColumn`, `RemoveBinding`, Struktur-Kommentar am `Assembler`) und `receive.go` (`Stream.Assembler`); die Zeilen in Ports, Adaptern und Modell als Listing. **Nachgezogen:** die vier Kommentare der Bindungs-Nachträge (§3, dritte Zeile). **Nicht nachgezogen:** die Kommentare in `postgresstorage/administrationrequest.go`, `postgresstorage/backfilladmission.go`, `postgresstorage/queries/queries.go`, `port/outbound/administrationrequest.go` und `domain/model/administrationrequest.go` nennen die Goroutine als Leserin und Verarbeiterin der Queue; der Vorlauf ruft dieselbe Funktion über dieselben `administrationDeps` (Pool, Rolle) — „Administrations-Goroutine“ steht dort für die Verarbeitung, kein Widerspruch (Lesung der Aufrufstelle, nicht gemessen). **Nichtgefunden:** kein weiterer Kommentar behauptet, die Goroutine sei der einzige Aufrufer der Bindungs-Nachträge. | vier Kommentare im Ist-Ton |
| Startpfad des Backfill-Workers und des Start-Abgleichs (aus `slice-backfill-sql-administration`) | Lesen von `Run` in `internal/bootstrap/wiring.go` am Start | **Gefunden.** Die Folge in `Run` ist Bindungsaufbau (`activatedTableBindings`, `stream.BindCapture`) → `reconcileBackfillRuns` → Worker-Start → Sequenz (Vorlauf, Goroutine, `stream.Run`); der Abgleich `running` → `interrupted` und der Worker-Start stehen vor dem Vorlauf. **Nichtgefunden:** kein Pfad, auf dem der Vorlauf vor dem Abgleich läuft; das Wecksignal `backfillWake` besteht seit dem Anlegen vor dem Worker-Start (Kapazität 1, verschmelzend). Gebunden durch `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall` (Quelltext-Positionen der Aufrufe, Grenze in §3 zweite Zeile). | keiner |
| Andere Pläne und ADRs, die die Ordnung als Träger nennen | Zeilen 7–8 und 17–18 im Block unten | **Gefunden.** Zeilen 7–8 („Startreihenfolge“): Parent 10 (vier in `welle-transformationen`, zwei in `slice-transformationen-e2e-abhilfe`, drei in `ADR-0112`/`ADR-0113`, einer der Code-Kommentar zum Backfill-Start), Diff 11 (dazu der Kopfkommentar der neuen Testdatei). `slice-transformationen-e2e-abhilfe` (§1, §2, §4, §6) nennt „den Ordnungs-Test aus `start-reihenfolge`“ und die Startreihenfolge als Voraussetzung; `welle-transformationen` §3/§4/§5 beschreibt den Slice. Beide Aussagen bleiben wahr. **Nichtgefunden:** kein Plan behauptet eine Beschaffenheit der Ordnung, die der Slice ändert. Ein Träger außerhalb der Änderung: [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) führt die Abhilfe im gescheiterten Prozess als „erwartet, nicht am Code belegt“ (Zeilen 17–18: 1 Treffer in `ADR-0112`, Parent und Diff) — die Aussage ist eine `Accepted` ADR und bleibt unberührt (`AGENTS.md` §3.5); ihren Beleg trägt `slice-transformationen-e2e-abhilfe`. | Meldung an `e2e-abhilfe` (Fremddatei): die Tests tragen die Namen `TestRunStreamAfterAdministrationPass…` und `TestRunSourceText…` in `internal/bootstrap/administration_startorder_internal_test.go` (Frist: Closure dieses Slice, danach zieht der Planner nach oder benennt den Träger mit Adresse, `AGENTS.md` §3.13). **In der Closure gezogen:** `slice-transformationen-e2e-abhilfe` §1, §2 und §6 nennen die Namen der Tests, ihre Grenze (Quelltext-Lesung der Aufrufstelle, Fakes für die Sequenz, kein Beleg am System) und den Runner-Schritt, der einen beim Start `pending` stehenden Antrag über einen Prozessstart legt |

Die „Diff“-Stände der Zeilen dieses Blocks stehen seit der Closure auf dem Commit
`cd5ff4f2` (der Stand, den Implementer und Verifikation nachmaßen): der Arbeitsbaum
bewegt sich mit jeder Änderung an Registern, Wellen-Plan und Nachbar-Plänen, die die
Closure vornimmt (`make suchlauf-nachmessen` mit `diff` fand dort vier Abweichungen bei den Zeilen des
Slice-Diffs, etwa 146 statt 144 Treffer der Symbole, weil die Evidence-Dateien anderer
Einträge und der Wellen-Plan sie nennen).

```suchlauf
ce229b8e 134 -n -E 'runAdministration|stream\.Run|processAdministrationRequests' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 144 -n -E 'runAdministration|stream\.Run|processAdministrationRequests' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 141 -n -E 'Administrations-Goroutine|Antrags-Queue' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 139 -n -E 'Administrations-Goroutine|Antrags-Queue' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 7 -n -E 'vor .?stream\.Run|nach .?stream\.Run' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 8 -n -E 'vor .?stream\.Run|nach .?stream\.Run' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 10 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 11 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 4 -n -F 'assembliert wird' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 4 -n -F 'assembliert wird' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 11 -n -E 'erste Transaktion|bevor der Stream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 15 -n -E 'erste Transaktion|bevor der Stream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 0 -n -E 'Startreihenfolge|Start-Reihenfolge|assembliert wird|erste Transaktion|bevor der Stream|vor .?stream\.Run|nach .?stream\.Run' -- docs/user harness
cd5ff4f2 0 -n -E 'Startreihenfolge|Start-Reihenfolge|assembliert wird|erste Transaktion|bevor der Stream|vor .?stream\.Run|nach .?stream\.Run' -- docs/user harness
ce229b8e 22 -n -F 'Administrations-Goroutine' -- internal ':!*_test.go'
cd5ff4f2 20 -n -F 'Administrations-Goroutine' -- internal ':!*_test.go'
ce229b8e 1 -n -F 'nicht am Code belegt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cd5ff4f2 1 -n -F 'nicht am Code belegt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-e2e-wirkung`
in `done/` liegt (Reihenfolge der Welle) und
`slice-backfill-sql-administration` in `done/` liegt (beide Slices ändern `Run`
im Bereich der Administrations-Goroutine; der Backfill-Zweig und der
Start-Abgleich stehen dann fest) und `slice-antragsqueue-lesefehler-failed` in
`done/` liegt (Kante: der Vorlauf vor `stream.Run` liest dieselbe Queue, und
eine Zeile, die die Lesung anhält, hält auch ihn an — hergeleitet aus
`processAdministrationRequests`, nicht erprobt; beide Slices ändern die
Verarbeitung der Queue in `internal/bootstrap/wiring.go`) und kein anderer Slice
in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  kleiner Zug im Startpfad; sprengt die Extraktion der testbaren Stelle den
  Umfang, ist der abtrennbare Teil der Test der Ordnung (zweiter Liefer-Punkt)
  als eigener Slice.
- `in-progress` → `open` (blockiert): falls `Run` ohne Umbau, der über die
  kleine Extraktion hinausgeht, keine Stelle für einen ordnungsprüfenden Test
  bietet (dann Architect-Frage zum Zuschnitt der Composition Root) oder falls
  der Vorlauf mit einem `backfill`-Antrag kollidiert, den Abgleich und Worker
  noch nicht aufgenommen haben.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) und ein
realer, grüner `make test-integration`-Lauf + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Der Vorlauf verzögert den Stream-Start** (eine hängende Datenbank-Abfrage
  im Vorlauf hält den Start an). Der Stream-Start wartet auf den vollständigen
  Vorlauf. Der Vorlauf läuft unter `streamCtx` (den der Prozess-`ctx` und die
  WAL-Fehlerschwelle beenden), die Goroutine danach unter `administrationCtx`
  (den nur der Prozess-`ctx` beendet): die Abbruch-Semantik beider ist nicht
  dieselbe. Dass der Vorlauf unter `streamCtx` läuft, ist aus `Run` gelesen und
  an keinen Test gebunden (Verifikation V-2; Grenze in §3). Eine eigene Frist je
  Antrag oder je Durchlauf trägt der Vorlauf nicht. Belegt: ein Test mit blockierendem Fake zeigt, dass ein Abbruch des
  Kontexts den Vorlauf beendet
  (`TestRunStreamAfterAdministrationPassHoldsTheStreamUntilThePassEndsAndContextCancelEndsIt`).
  **Hergeleitet aus dem Code, nicht am laufenden Prozess gemessen:** ein Antrag,
  der lange läuft, ohne zu hängen (etwa `ALTER PUBLICATION … ADD TABLE`, die auf
  eine Tabellen-Sperre wartet), oder eine lange Queue hält die Erfassung der
  Quelle an; der Heartbeat startet vor der Sequenz und der `--healthcheck` liest
  nur dessen Alter — ein wartender Vorlauf erscheint als gesund. Die
  Betreiber-Aussage dazu trägt `slice-transformationen-betriebsdoku` (Aufschub
  mit Adresse). Ob der Vorlauf eine Zeitgrenze trägt (Frist je Antrag oder je
  Durchlauf), was bei Ablauf geschieht und ob `diagnose`/`--healthcheck` das
  Warten zeigen, ist eine Entscheidung, die dieser Plan nicht trifft.
  **Ausgang: weiter offen.** Träger: der Register-Eintrag
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` (neu in der Closure, 1×)
  und die Frage (e) in [welle-transformationen](../welle-transformationen.md) §5
  samt Closure-Kriterium in §3; Adresse: ein Architect-Verdikt vor der Closure
  der Welle (das Ereignis kann eintreten: die Welle trägt den Closure-Trigger
  in §3), im Umfang von `slice-transformationen-e2e-abhilfe` (der Beleg fährt
  den Vorlauf mit wartendem Antrag am System) und
  `slice-transformationen-betriebsdoku` (Betreiber-Aussage).
- **Ein `backfill`-Antrag im Vorlauf trifft einen Zustand ohne Abgleich oder
  Worker.** Der Backfill-Zweig nimmt den Antrag an (Run-Zeile `queued` und Antrag
  `applied` in einer Transaktion, dauerhaft, `slice-backfill-sql-administration`)
  und weckt den Worker; der Worker liest die `queued`-Zeilen beim Start, nach dem
  Abgleich `running` → `interrupted`. Der Vorlauf darf den Zweig nicht in einen
  Zustand bringen, in dem der Abgleich noch aussteht — bei einem `running`-Run
  einer früheren Prozessinstanz endet der Antrag an der Prüfung „kein aktiver
  Run" als `failed` (hergeleitet aus der Annahme-Regel von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1). *Zu belegen durch:* Lesen von `Run` am
  Start (Suchlauf-Tabelle, vierte Zeile: Startpfad des Backfill-Workers) und
  `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall`
  (Quelltext-Positionen, kein Fake-Worker).
  **Ausgang: entfallen.** Die Folge in `Run` ist Bindungsaufbau →
  `reconcileBackfillRuns` → Worker-Start (das Wecksignal `backfillWake` besteht
  seit dem Anlegen vor dem Worker-Start) → Sequenz; gelesen im Implementer-Lauf,
  vom Reviewer und vom Verifier (`wiring.go`, Aufrufe in Zeilen 823, 830 und
  1041, in der Closure mit `grep -n` nachgemessen), und die Mutation „Abgleich
  hinter die Sequenz“ färbt den Test rot (Review M11, Verifikation V11,
  **übernommen**). Grenze: ein Aufruf des Abgleichs in einer toten Verzweigung
  bleibt grün (Verifikation V14, in §3 benannt); ein beim Start wartender
  `backfill`-Antrag am System ist nicht gefahren (DoD 3: kein Runner-Schritt
  legt einen `pending` stehenden Antrag über einen Start).
- **Der Vorlauf ändert das Verhalten für `enable`/`disable`**: eine Bindung
  entsteht vor dem Stream statt danach. Das ist die gewollte Eigenschaft; ein
  Test belegt die Reihenfolge für eine `enable`-Anfrage. **Ausgang:
  eingetreten** als gewollte Eigenschaft, kein Carveout und kein Folge-Slice:
  `TestRunStreamAfterAdministrationPassBindsAnEnabledTableBeforeTheStreamStarts`
  belegt die Bindung vor dem Stream-Start und färbt sich bei einer Antragsart
  `disable` rot (Mutation E2, Verifikation, **übernommen**); die Wirkung am
  System (ein beim Start wartender `enable`-Antrag) trägt kein Runner-Schritt
  (DoD 3).
- **Ein Kommentar behauptet mehr, als der Code trägt** — etwa, dass jede Race
  beseitigt sei (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`,
  verkörpert, 6× mit dieser Closure): der Kommentar sagt nur zu, dass der
  Vorlauf vor dem Stream läuft. Der Review fand die Zusage ohne ihre Bedingung
  (F-2, Lesefehler und gescheiterter Vermerk); der Godoc nennt die Bedingung.
  *Zu belegen durch:* Nachlesen des Godocs gegen `processAdministrationRequests`
  in der Verifikation (§5 Zeile F-2). **Ausgang: entfallen** für den Godoc der
  Sequenz: er nennt Bedingung, Lesefehler-Pfad, Vermerk-Pfad und die fehlende
  Frist (Verifikation, gelesen gegen den Code). **Rest:** „die Ordnung gilt für
  jede Antragsart“ meint für `backfill` „angenommen“ (Verifikation V-3, INFO);
  der Rest ist in §3 benannt, im Godoc nicht geändert und mit V-3 in der
  Evidence-Datei dieses Slice unter
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` geführt.
- **Abweichung vom ADR-Wortlaut.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 5 legt „offene Anträge vor `stream.Run` verarbeiten“ in den
  E2E-Slice, **wenn** Kriterium (c) nicht getragen wird; dieser Plan
  entscheidet die Bedingung am Code (keine Synchronisation, siehe §1) und macht
  den Zug zum eigenen Slice vor dem E2E-Beleg
  (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, offen, 2×). Als Frage
  an den Auftraggeber geführt, nicht als Entscheidung dargestellt (Welle §4,
  Abweichung 2). **Ausgang: entfallen** für das Ergebnis: die Ordnung steht im
  Code und ist an die Eingabeseite gebunden (Verifikation §4, V1 bis V3 und E1
  bis E4 rot), die Aussage von Folgepflicht 5 ist getragen, und die
  ADR-Präambel überlässt den Schnitt in Slices dem Planner („der Planner formt
  daraus Welle und Slices“, gelesen an
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Folgepflichten); die Verifikation liest die Lieferung als konform. Der
  Vorgang zählt im Register nicht als drittes Auftreten (Begründung in der
  `state.md` des Eintrags); die Zuschnitts-Frage bleibt in
  [welle-transformationen](../welle-transformationen.md) §4 (Abweichung 2)
  geführt und trägt keine Verantwortung dieses Slice mehr.

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Sequenz „Vorlauf, Goroutinen-Start, Stream“
  steht als Funktion mit drei Zeilen Rumpf (`runStreamAfterAdministrationPass`)
  an einer Stelle, die vier Tests ohne Datenbank aufrufen; die Leser-Kette lief
  mit einer Fixrunde ohne Rückführung: der Review nennt 1 HIGH, 2 MEDIUM, 3 LOW
  und 3 INFO, die Verifikation 0 HIGH, 0 MEDIUM, 2 LOW und 4 INFO (beide
  **übernommen** aus den Summary-Tabellen der Reports), F-1 bis F-7 sind durch
  die Fixrunde (`eb4c2b1c`, `990608d3`, `cd5ff4f2`) behoben, F-8 und F-9 stehen
  als benannte Grenze in §3. (2) Die Belege binden an die Eingabeseite: der
  Review fuhr 18 Mutationen (17 aussagekräftig, 16 rot, M10 grün), die
  Verifikation 18 (14 rot; grün V9b, V10b, V12, V14 — alle vier an der
  Aufrufstelle in `Run`; beide Zahlen **übernommen** aus den Reports); V10b fuhr
  die Verifikation am System und sah `make test-integration` rot (Exit 2), damit
  ist der Goroutinen-Start am laufenden Prozess gedeckt (**übernommen**). (3) Die
  Leser fanden, was der Implementer-Lauf nicht fand: den falschen Satz
  „Nichtgefunden“ im Suchlauf-Feld (F-1, HIGH), die Zusage ohne Bedingung im
  Godoc (F-2) und die fehlende Benennung der Wartegrenze (F-3); die Verifikation
  bestätigte die Fixrunde je Finding gegen den Code, ihre Gegensuche fand keinen
  weiteren Träger der Reihenfolge. (4) Die Closure maß nach, was Report und Plan
  nur nannten: `git grep -n "'pending'" -- tools/harness/run-integration-tests.sh
  test/integration` Exit 1 (kein Runner-Schritt legt einen `pending` stehenden
  Antrag über einen Start), `git diff --stat cd5ff4f2 HEAD -- internal tools test
  harness Makefile` ohne Ausgabe (Code- und Skript-Stand seit dem Lauf des
  Verifiers unverändert), `make suchlauf-nachmessen` („18 Zeilen stimmen“ vor den
  Closure-Zeilen dieser Notiz), die Zeilen 823, 830 und 1041 in `wiring.go` (Abgleich,
  Worker-Start, Sequenz), die sieben Antragsarten von `model.AdministrationRequestKinds`
  und die Zähler der Register-Einträge mit `ls evidence | wc -l`.
- **Was ging anders als geplant:** (1) **Das Suchlauf-Feld trug ein falsches
  „Nichtgefunden“** (F-1, HIGH): die Muster stammten aus dem Planungs-Feld, das
  vor der Suchform von [`AGENTS.md`](../../../../AGENTS.md) §3.13 geschnitten ist;
  `make suchlauf-nachmessen` war grün, weil es Zahlen und Stände misst, nicht
  Suchraum und Muster. Die Spec trug die Zusage („beim Start … bevor … assembliert“),
  das Muster suchte das Wort „Startreihenfolge“. (2) **Der Godoc der Sequenz sagte
  die Ordnung ohne die Bedingung des Lesens zu** (F-2, MEDIUM) und kannte die
  Grenze der Wartezeit nicht (F-3, MEDIUM). (3) **Der Vorlauf ist ein neuer
  Warte-Zustand vor der Erfassung**, den weder Spec noch ADR benennen (V-4): die
  Entscheidung steht als Frage (e) bei der Welle (Steering-Loop (c)). (4) **Zwei
  Testnamen sagten Lauf-Reihenfolge, gelesen wird Quelltext** (F-4, LOW). (5)
  **DoD-Haken 3 stand auf `[x]` ohne committeten Anker** (F-5, LOW): die Fixrunde
  setzte ihn zurück, die Closure trägt den Report der Verifikation als Anker samt
  Reichweite ein. (6) **Reichweite des Integrationslaufs:** der Lauf des Verifiers
  belegt den Dauerbetrieb und den Vorlauf mit leerer Queue, nicht den Vorlauf mit
  wartendem Antrag; diese Wirkung trägt `slice-transformationen-e2e-abhilfe`. (7)
  **Umfang:** 11 Commits `ce229b8e..f212b904` (`git rev-list --count`), 8 Dateien mit
  1311 Einfügungen und 284 Löschungen einschließlich der drei Lifecycle-Moves und der
  zwei Berichte; ohne `docs/reviews` 6 Dateien mit 691 und 284, die Go-Dateien 4 mit
  391 und 40 (gemessen in der Closure mit `git diff --shortstat`; die Plan-Datei
  erscheint als Neuanlage und Löschung, weil der Diff ab `ce229b8e` den Move nicht als
  Rename führt). (8) **Coverage:** 85.30 % im ersten `make gates`-Lauf der Closure (gedruckt:
  „coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%“, Exit 0), dieselbe Zahl im
  Lauf der Verifikation (**übernommen**); Schwelle 80 % — der Beleg je eines Laufs, kein
  Ist-Stand. (9) **Nicht gefahren:** `make
  test-store`, `make test-replication`, `make bench` (kein Adapter im Diff; Aussage
  der Verifikation, **übernommen**); nicht gepusht, kein Workflow im Diff — kein
  Lauf von `e2e.yml` nach dem Push, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift
  nicht.
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor:* keiner gebaut. F-1
  fiel durch `make suchlauf-nachmessen` (Exit 0), weil das Werkzeug Zahlen und
  Stände prüft; ein Sensor, der prüft, ob das Muster die Aussage trägt, müsste
  wissen, mit welchen Wörtern die Spec die bewegte Eigenschaft beschreibt (eine
  Semantik-Entscheidung, benannt in `harness/sensors/suchlauf-nachmessen.md`
  §Grenze). Die verfügbare Falsifikation ist die Gegensuche; sie fuhren Reviewer und
  Verifier. *(b) Geschärfte Regel:* keine neue Regel; zwei Anwendungen mit Adresse.
  Erstens die Suchform selbst: ein Planungs-Feld, das vor ihr geschnitten ist,
  führt Symbolnamen und enge Suchräume und stützt darauf ein „Nichtgefunden“
  (Register `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, 15×); nachgezogen in
  der Closure ist `slice-transformationen-e2e-abhilfe` §3, die übrigen Pläne unter
  `open/` tragen ihre Adresse im Register-Eintrag (Start des jeweiligen Slice).
  Zweitens die Klausel „Zusage“ des Reviewer-Skills: eine Zusage über eine Ordnung
  nennt die Bedingung, von der sie abhängt (Register
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, 6×); die Probe „den
  zugesagten Pfad nachfahren“ hat vor dem Merge gegriffen. *(c) Benannte
  Spec-Lücke:* das Verhalten des Erfassungsstarts **während** des Vorlaufs — Frist,
  Verhalten bei Ablauf, Sichtbarkeit in `diagnose` und `--healthcheck`. Das
  Pflichtenheft (`LH-FA-CFG-007.a`) verlangt die Abfolge „Antrag vor erster
  Transaktion“ ohne Aussage zur Dauer des Wartens, `LH-FA-ADM-003` verlangt einen
  erkennbaren Fehlerzustand, kein erkennbares Warten (`git grep -n -i -E
  'Vorlauf|Zeitgrenze|eigene Frist|Wartegrenze' -- spec` trifft keine Zeile,
  gemessen in der Closure). Adresse: Frage (e) in
  [welle-transformationen](../welle-transformationen.md) §5, Closure-Kriterium in
  §3 und der Register-Eintrag `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad`.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-transformationen-start-reihenfolge.md`; Zähler = Zahl der Dateien
  (gemessen mit `ls evidence | wc -l` am Stand dieser Closure). *Neue Belege:*
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` **15×** (F-1 HIGH, daher Datei trotz
  Deckel; verkörpert); `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`
  **6×** (F-2 MEDIUM, V-3 INFO; verkörpert);
  `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt` **4×** (F-4 LOW; offen,
  Schwelle seit dem dritten Beleg erreicht);
  `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` **3×** (F-5 LOW; offen,
  **Schwelle erreicht**); `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` **3×**
  (V-2 LOW, F-8 INFO; offen, **Schwelle erreicht**). *Neuer Eintrag:*
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` **1×** (F-3 MEDIUM, V-4 LOW;
  offen, adressiert: Frage (e) und Closure-Kriterium der Welle). *Ohne Zählung
  fortgeschrieben:* `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` bleibt
  **2×** (`state.md` trägt, warum der Zuschnitt von `start-reihenfolge` nicht zählt).
  *Deckel-Fälle ohne Datei, Finding-Kennung hier* (vor dem Merge von Reviewer bzw.
  Verifier gefunden, Schwere ≤ LOW, bekannter Träger-Typ): F-6 (LOW,
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, Deckel bei 10×: die Meldung an
  `e2e-abhilfe` nannte im committeten Feld als Frist den Start des fremden Slice statt
  der Closure des meldenden); F-7 (INFO, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
  Deckel bei 32×: vier Kommentare außerhalb des Diffs nannten die
  Administrations-Goroutine als einzigen Aufrufer der Bindungs-Nachträge). *Kein
  Register-Anfall:* F-9 (INFO, eine Alternative, im Plan §3 als nicht geliefert
  begründet), V-1 (INFO, Reichweite, im Plan benannt), V-5 (INFO, Kosmetik ohne
  Änderung, §3), V-6 (INFO, Übernommenes im Report benannt). *Lese-Schritt der
  Closure von `welle-transformationen`:* aus diesem Slice erreichen
  `plan-zusage-erfuellung-ohne-committeten-anker` und
  `adapter-unittest-verdeckt-bootstrap-luecke` die Schwelle ohne Ausgang;
  `test-name-behauptet-mehr-als-der-test-treibt` steht mit 4× dort. Kein Slice ist
  wegen einer Beobachtung dieses Slice fällig: die Träger sind Zeilen in Skill oder
  Regelwerk (die beiden ersten) und die Frage eines Bootstrap-Smoke-Tests (der
  dritte), also Architect-Züge des Lese-Schritts der Welle-Closure, deren Ereignis
  eintreten kann (die Welle trägt den Closure-Trigger „Lese-Schritt“ in §3).
- **Folge-Slices:** keine angelegt. Nächster startbarer Slice der Kette:
  `slice-capture-leerlauf-quellbelege` (`open/`, wellenlos) — die Tabelle der Welle
  (§4) nennt als nächsten `slice-transformationen-e2e-abhilfe`, dessen Start-Trigger
  (§4 dort) `slice-capture-leerlauf-quellbelege` in `done/` verlangt (Belegaufbau
  der Container-Ende-Grenze im selben Runner); der Start-Trigger von
  `slice-capture-leerlauf-quellbelege` ist mit dem Move dieses Slice erfüllt (kein
  anderer Slice in `in-progress/`, gelesen mit `ls` in der Closure). Danach folgt
  `slice-transformationen-e2e-abhilfe`, zuletzt `slice-transformationen-betriebsdoku`
  (Start nach `e2e-abhilfe` und `slice-sdk-regel-realserver-e2e`, dessen Start-Trigger
  seit `e2e-wirkung` erfüllt ist). Übergaben mit Adresse (in der Closure gezogen):
  `slice-transformationen-e2e-abhilfe` §1, §2, §3, §6, §8 (Namen und Grenze der
  Ordnungs-Tests, Runner-Schritt mit wartendem Antrag, Suchform, Risiko der
  Zeitgrenze); `welle-transformationen` §3 (Closure-Kriterium), §4 (Stand des
  Schnitts in Abweichung 2) und §5 (Frage (e)). Ein Umsetzungs-Slice im Startpfad
  entsteht nur, wenn das Architect-Verdikt zu Frage (e) eine Frist oder eine Anzeige
  beschließt.
- **Risiken aus §6:** je ein Ausgang, mit Begründung in §6. *Weiter offen:* der
  Vorlauf verzögert den Stream-Start (Register
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad`, Frage (e), Architect-Verdikt vor
  der Closure der Welle). *Eingetreten:* der Vorlauf ändert das Verhalten für
  `enable`/`disable` — als gewollte Eigenschaft, kein Carveout, kein Folge-Slice.
  *Entfallen:* `backfill`-Antrag im Vorlauf ohne Abgleich oder Worker · Kommentar
  behauptet mehr, als der Code trägt (für den Godoc der Sequenz; Rest V-3 in der
  Evidence-Datei) · Abweichung vom ADR-Wortlaut (für das Ergebnis; die
  Zuschnitts-Frage bleibt bei der Welle, §4 Abweichung 2).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure. Vorab gelesen: *Anker* — dieser Slice trägt
  kein Feld „liegt in“; *Folge-Slice* — `slice-capture-leerlauf-quellbelege`,
  `slice-transformationen-e2e-abhilfe` und `slice-transformationen-betriebsdoku` liegen
  als Dateien in `open/` (gelesen mit `ls` in der Closure); *Register* — die sieben
  genannten Kennungen `BEO-PGC/…` sind Verzeichnisse mit nicht leerem `evidence/`
  (gemessen mit `ls evidence | wc -l` in der Closure: 15, 6, 4, 3, 3, 1 und 2).

**§3.13-Suchlauf der Closure** (bewegte Eigenschaften: „die Reihenfolge von
Antrags-Verarbeitung und Stream-Start“, die Beschreibungen des Wartens vor der
Erfassung; Parent `f212b904` = `HEAD` vor den Closure-Änderungen, Diff = Commit `9320715c`
mit den Änderungen an Registern, Wellen-Plan und `e2e-abhilfe`; Suchraum der ganze Baum,
ausgenommen `docs/reviews/**`, die Records unter `done/` und `.harness/baseline/**`;
das Werkzeug schließt jede Datei mit dem Basisnamen der Plan-Datei aus, also auch die
sieben Evidence-Dateien dieses Slice). *Gefunden:* die
Beschreibungen der Reihenfolge „keine Synchronisation“ und „ohne Reihenfolge-Garantie“
stehen in `ADR-0112` (`Accepted`, unberührt, `AGENTS.md` §3.5) und in `welle-transformationen`
§4 (Abweichung 2; nachgezogen: „Stand des Schnitts“, die Ordnung liefert dieser Slice);
der Treffer in `mapper.go` betrifft die Bindungstabelle, nicht die Startreihenfolge. Das
Wort „Startreihenfolge“ steht unverändert in 11 Zeilen (Zeilen 5–6 des Blocks). Die
Zeilen zu „Zeitgrenze“, „eigene Frist“ und „Wartegrenze“ wachsen von 22 auf 31 um die
Zeilen dieser Closure (Frage (e) und Closure-Kriterium im Wellen-Plan, Risiko in
`e2e-abhilfe`, Zeilen des Registers, im Diff gelesen); die 22 Zeilen des Parent
betreffen die Lesesperre des Backfills und Test-Zeitgrenzen, keine beschreibt eine
Frist des Vorlaufs. *Nichtgefunden:* kein Träger
außer `wiring.go` (Godoc) beschreibt die Frist des Vorlaufs; `docs/user` und `harness`
tragen keine Beschreibung der Reihenfolge (Zeilen 13–14 des Feldes in §3). Grenze:
Träger, die die Reihenfolge in anderen Worten beschreiben, trifft das Muster nicht
(Lese-Handlung des Reviewers).

```suchlauf
f212b904 4 -n -E 'Reihenfolge-Garantie|keine Synchronisation|ohne Synchronisation|ohne Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
9320715c 4 -n -E 'Reihenfolge-Garantie|keine Synchronisation|ohne Synchronisation|ohne Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
f212b904 22 -n -E 'Zeitgrenze|eigene Frist|Wartegrenze' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
9320715c 31 -n -E 'Zeitgrenze|eigene Frist|Wartegrenze' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
f212b904 11 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
9320715c 11 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die Composition Root ist keine eigene Sub-Area — kein
Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
zum Stand der Planung; die Zähler der Closure stehen in §7) —
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×,
einschlägig — Risiko §6), `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`
(offen, 1×, einschlägig — Risiko §6),
`BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×, einschlägig —
Abgrenzung §1), `BEO-PGC/adapter-fehler-ausgang` (offen, 2×, gesichtet — der
Vorlauf bleibt best effort, ein Retry ist nicht Teil),
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (offen, 2×, einschlägig —
die Ordnung ist durch einen Bootstrap-Test getragen, nicht durch
Adapter-Tests), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×, DoD Punkt 2), `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
