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
  Fehlen einer Race nicht beweisen (bestanden ohne Ordnung wäre Glück), der
  Beleg der Ordnung selbst gehört in diesen Slice, die Wirkung am System in
  `e2e-abhilfe`.
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
- [ ] Der Dauerbetrieb bleibt unverändert: `make coverage-gate` grün (Teil von
      `make gates`, nächster Haken) und ein realer, grüner
      `make test-integration`-Lauf mit den bestehenden Belegen der
      SQL-Administration einschließlich des Live-Reload-Belegs ohne
      `docker restart`. Der Lauf zeigt den Dauerbetrieb und den Vorlauf mit
      **leerer** Queue bei jedem Start; er zeigt **nicht** den Vorlauf mit einem
      beim Start wartenden Antrag: `git grep -n "'pending'" --
      tools/harness/run-integration-tests.sh test/integration` trifft keine Zeile
      (gemessen in der Fixrunde), also legt kein Runner-Schritt einen `pending`
      stehenden Antrag über einen Start. Diese Wirkung am System trägt
      `slice-transformationen-e2e-abhilfe`. Ein committeter Beleg-Anker für den
      Integrationslauf liegt nicht vor; *zu belegen durch:* der Lauf des
      Verifiers.
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
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` | update | **Geliefert.** Neue Funktion `runStreamAfterAdministrationPass(ctx, deps, startLoop, runStream)`: erst `processAdministrationRequests` (Vorlauf), dann `startLoop` (die Goroutine), dann `runStream`. Die `administrationDeps` werden vor dem Start gebildet (Variable `administration`), der Goroutinen-Start ist die Funktion `startAdministration`; `Run` ruft die Sequenz mit `stream.Run` als viertem Argument. Der Vorlauf steht **vor** dem Goroutinen-Start, damit zu keinem Zeitpunkt zwei Durchläufe dieselbe Queue lesen. Gelesen am Start (Schritt 12): `reconcileBackfillRuns` und der Start des Backfill-Workers stehen in `Run` vor dem Aufruf der Sequenz (Risiko §6, zweiter Punkt), der Kontext des Vorlaufs ist `streamCtx` (abgeleitet aus `ctx`, ein Abbruch durch die WAL-Schwelle beendet den Vorlauf mit). **Fixrunde (F-2, F-3):** der Godoc der Sequenz sagt die Ordnung mit ihrer Bedingung zu — sie gilt für einen Antrag, den der Vorlauf liest und verarbeitet; bei einem Lesefehler von `ListPending` bleibt jeder Antrag `pending` und der Stream startet mit dem bisherigen Regelstand, bei einem gescheiterten `MarkApplied` steht die Wirkung in der Bindung und der Antrag bleibt `pending`, bis die Goroutine ihn erneut verarbeitet (aus `processAdministrationRequests` und `applyAdministrationRequest` gelesen; der Lesefehler-Test bindet „hält den Start nicht an“ und den Rückgabewert, nicht die Folge für den Antrag). Der Godoc nennt außerdem, was den Vorlauf begrenzt (`ctx`, in `Run` `streamCtx`) und dass er keine eigene Frist trägt. |
| `internal/bootstrap/administration_startorder_internal_test.go` | neu (Plan-Nachzug) | Sechs Tests ohne Datenbank: Ordnung an einem `pending` stehenden `remove_transformation`-Antrag (Stream-Fake sieht Vermerk `applied` und Rohform), Ordnung für `enable`, Lesefehler hält den Start nicht an und der Stream-Ausgang bleibt der Rückgabewert, hängende Abfrage hält den Stream-Start an und der Abbruch des Kontexts beendet den Vorlauf, und zwei Tests über die Quelltext-Gestalt von `Run` (`TestRunSourceTextPassesStreamRunOnlyAsArgumentOfTheSequence`: `stream.Run` nur als Argument der Sequenz; `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall`: Aufrufe von Abgleich und Worker-Start vor dem Aufruf der Sequenz). Grund der Quelltext-Tests: `Run` braucht die Datenbank, kein netzloser Lauf erreicht die Aufrufstelle. **Fixrunde (F-4):** die beiden Namen tragen die Messung (Quelltext), nicht mehr eine Lauf-Reihenfolge. **Benannte Grenze (F-8, F-9):** die Aufrufstelle in `Run` ist nur über Quelltext-Positionen gebunden; der Rumpf von `startAdministration` (Start der Goroutine) ist an keinen netzlosen Test gebunden — eine Mutation, die ihn leert, bleibt grün (Mutation M10 des Reviews, **übernommen**), und ein Aufruf in einer toten Verzweigung (`if false { … }`) färbt den Test der Reihenfolge nicht rot (in einer Kopie gesehen, `go test -race` grün). Den Goroutinen-Start deckt am laufenden Prozess `make test-integration` (Live-Reload-Belege, Haken 3). Ein verhaltensbasierter Weg (Abgleich, Worker-Start, Vorlauf, Goroutine, Stream als Funktionswerte einer Start-Funktion) ersetzt die Quelltext-Tests nur, wenn `Run` diese Funktionswerte aus seinen Pools und Kontexten selbst baut — ein Umbau der Verdrahtung über die kleine Extraktion hinaus (§4, Rückführung), und die Mutation M10 bliebe dabei ungebunden, weil der Rumpf von `startAdministration` in `Run` entsteht. Nicht geliefert. |
| `internal/adapters/driving/replication/mapper/mapper.go`, `internal/adapters/driving/replication/receive/receive.go` | update (Kommentare, Plan-Nachzug, F-7) | Vier Kommentare nannten die Administrations-Goroutine als Aufrufer der Bindungs-Nachträge (`AddBinding`, `ExcludeColumn`, `RemoveBinding` in `mapper.go`, `Stream.Assembler` in `receive.go`); seit dem Slice ruft sie auch der Vorlauf vor dem Stream-Lauf. Sie nennen jetzt die Administrations-Verarbeitung (Goroutine und Vorlauf), je Kommentar höchstens eine Kennung (`AddBinding` und `ExcludeColumn` trugen drei bzw. zwei; `make kommentar-kennungen DIFF=ce229b8e` meldet jeden berührten Block mit mehr als einer Kennung). Kein Verhaltens-Diff, kein Eingriff in den Capture-Pfad (§1, dritter Punkt: `stream.Run`, `CaptureService`, ACK bleiben unberührt). Der Struktur-Kommentar am `Assembler` (`mapper.go`, Synchronisation der Bindungstabelle) beschreibt die Nebenläufigkeit im Dauerbetrieb und bleibt wahr: der Vorlauf läuft vor dem Stream-Lauf. |
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
| Andere Pläne und ADRs, die die Ordnung als Träger nennen | Zeilen 7–8 und 17–18 im Block unten | **Gefunden.** Zeilen 7–8 („Startreihenfolge“): Parent 10 (vier in `welle-transformationen`, zwei in `slice-transformationen-e2e-abhilfe`, drei in `ADR-0112`/`ADR-0113`, einer der Code-Kommentar zum Backfill-Start), Diff 11 (dazu der Kopfkommentar der neuen Testdatei). `slice-transformationen-e2e-abhilfe` (§1, §2, §4, §6) nennt „den Ordnungs-Test aus `start-reihenfolge`“ und die Startreihenfolge als Voraussetzung; `welle-transformationen` §3/§4/§5 beschreibt den Slice. Beide Aussagen bleiben wahr. **Nichtgefunden:** kein Plan behauptet eine Beschaffenheit der Ordnung, die der Slice ändert. Ein Träger außerhalb der Änderung: [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) führt die Abhilfe im gescheiterten Prozess als „erwartet, nicht am Code belegt“ (Zeilen 17–18: 1 Treffer in `ADR-0112`, Parent und Diff) — die Aussage ist eine `Accepted` ADR und bleibt unberührt (`AGENTS.md` §3.5); ihren Beleg trägt `slice-transformationen-e2e-abhilfe`. | Meldung an `e2e-abhilfe` (Fremddatei): die Tests tragen die Namen `TestRunStreamAfterAdministrationPass…` und `TestRunSourceText…` in `internal/bootstrap/administration_startorder_internal_test.go` (Frist: Closure dieses Slice, danach zieht der Planner nach oder benennt den Träger mit Adresse, `AGENTS.md` §3.13) |

```suchlauf
ce229b8e 134 -n -E 'runAdministration|stream\.Run|processAdministrationRequests' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 144 -n -E 'runAdministration|stream\.Run|processAdministrationRequests' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 141 -n -E 'Administrations-Goroutine|Antrags-Queue' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 139 -n -E 'Administrations-Goroutine|Antrags-Queue' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 7 -n -E 'vor .?stream\.Run|nach .?stream\.Run' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 8 -n -E 'vor .?stream\.Run|nach .?stream\.Run' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 10 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 11 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 4 -n -F 'assembliert wird' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 4 -n -F 'assembliert wird' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 11 -n -E 'erste Transaktion|bevor der Stream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 15 -n -E 'erste Transaktion|bevor der Stream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
ce229b8e 0 -n -E 'Startreihenfolge|Start-Reihenfolge|assembliert wird|erste Transaktion|bevor der Stream|vor .?stream\.Run|nach .?stream\.Run' -- docs/user harness
diff 0 -n -E 'Startreihenfolge|Start-Reihenfolge|assembliert wird|erste Transaktion|bevor der Stream|vor .?stream\.Run|nach .?stream\.Run' -- docs/user harness
ce229b8e 22 -n -F 'Administrations-Goroutine' -- internal ':!*_test.go'
diff 20 -n -F 'Administrations-Goroutine' -- internal ':!*_test.go'
ce229b8e 1 -n -F 'nicht am Code belegt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 1 -n -F 'nicht am Code belegt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-e2e-wirkung`
in `done/` liegt (Reihenfolge der Welle) und
`slice-backfill-sql-administration` in `done/` liegt (beide Slices ändern `Run`
im Bereich der Administrations-Goroutine; der Backfill-Zweig und der
Start-Abgleich stehen dann fest) und `slice-antragsqueue-lesefehler-failed` in
`done/` liegt (Kante: der Vorlauf vor `stream.Run` liest dieselbe Queue, und
eine Zeile, die die Lesung anhält, hielte auch ihn an — hergeleitet aus
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
  dieselbe. Eine eigene Frist je Antrag oder je Durchlauf trägt der Vorlauf
  nicht. Belegt: ein Test mit blockierendem Fake zeigt, dass ein Abbruch des
  Kontexts den Vorlauf beendet
  (`TestRunStreamAfterAdministrationPassHoldsTheStreamUntilThePassEndsAndContextCancelEndsIt`).
  **Hergeleitet aus dem Code, nicht am laufenden Prozess gemessen:** ein Antrag,
  der lange läuft, ohne zu hängen (etwa `ALTER PUBLICATION … ADD TABLE`, die auf
  eine Tabellen-Sperre wartet), oder eine lange Queue hält die Erfassung der
  Quelle an; der Heartbeat startet vor der Sequenz und der `--healthcheck` liest
  nur dessen Alter — ein wartender Vorlauf erscheint als gesund. Die
  Betreiber-Aussage dazu trägt `slice-transformationen-betriebsdoku` (Aufschub
  mit Adresse). Ob der Vorlauf eine Zeitgrenze bekommt (etwa je Antrag), ist eine
  Entscheidung, die dieser Plan nicht trifft: Frage an den Planner/Architect im
  Bericht der Fixrunde. **Ausgang:** *(bei Closure)*
- **Ein `backfill`-Antrag im Vorlauf trifft einen Zustand ohne Abgleich oder
  Worker.** Der Backfill-Zweig nimmt den Antrag an (Run-Zeile `queued` und Antrag
  `applied` in einer Transaktion, dauerhaft, `slice-backfill-sql-administration`)
  und weckt den Worker; der Worker liest die `queued`-Zeilen beim Start, nach dem
  Abgleich `running` → `interrupted`. Der Vorlauf darf den Zweig nicht in einen
  Zustand bringen, in dem der Abgleich noch aussteht — ein `running`-Run einer
  früheren Prozessinstanz ließe die Prüfung „kein aktiver Run" den Antrag als
  `failed` enden (erwartet, hergeleitet aus der Annahme-Regel von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1). *Erwartet, zu belegen durch:* Lesen von `Run` am
  Start (Suchlauf-Tabelle, vierte Zeile: Startpfad des Backfill-Workers) und
  `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall`
  (Quelltext-Positionen, kein Fake-Worker).
  **Ausgang:** *(bei Closure)*
- **Der Vorlauf ändert das Verhalten für `enable`/`disable`**: eine Bindung
  entsteht vor dem Stream statt danach. Das ist die gewollte Eigenschaft; ein
  Test belegt die Reihenfolge für eine `enable`-Anfrage. **Ausgang:** *(bei
  Closure)*
- **Ein Kommentar behauptet mehr, als der Code trägt** — etwa, dass jede Race
  beseitigt sei (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`,
  offen, 2×): der Kommentar sagt nur zu, dass der Vorlauf vor dem Stream läuft.
  Der Review fand die Zusage ohne ihre Bedingung (F-2, Lesefehler und
  gescheiterter Vermerk); der Godoc nennt die Bedingung seit der Fixrunde.
  *Erwartet, zu belegen durch:* Review der Fixrunde. **Ausgang:** *(bei
  Closure)*
- **Abweichung vom ADR-Wortlaut.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 5 legt „offene Anträge vor `stream.Run` verarbeiten“ in den
  E2E-Slice, **wenn** Kriterium (c) nicht getragen wird; dieser Plan
  entscheidet die Bedingung am Code (keine Synchronisation, siehe §1) und macht
  den Zug zum eigenen Slice vor dem E2E-Beleg
  (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, offen, 1×). Als Frage
  an den Auftraggeber geführt, nicht als Entscheidung dargestellt (Welle §4,
  Abweichung 2). **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die Composition Root ist keine eigene Sub-Area — kein
Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
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
