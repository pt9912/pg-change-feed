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
- [x] Der Dauerbetrieb bleibt unverändert: `make test-integration` bleibt real
      grün (bestehende Belege der SQL-Administration einschließlich des
      Live-Reload-Belegs ohne `docker restart`); `make coverage-gate` grün. *Zu
      belegen durch:* ein realer, grüner `make test-integration`-Lauf.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
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
| `internal/bootstrap/wiring.go` | update | **Geliefert.** Neue Funktion `runStreamAfterAdministrationPass(ctx, deps, startLoop, runStream)`: erst `processAdministrationRequests` (Vorlauf), dann `startLoop` (die Goroutine), dann `runStream`. Die `administrationDeps` werden vor dem Start gebildet (Variable `administration`), der Goroutinen-Start ist die Funktion `startAdministration`; `Run` ruft die Sequenz mit `stream.Run` als viertem Argument. Der Vorlauf steht **vor** dem Goroutinen-Start, damit zu keinem Zeitpunkt zwei Durchläufe dieselbe Queue lesen. Gelesen am Start (Schritt 12): `reconcileBackfillRuns` und der Start des Backfill-Workers stehen in `Run` vor dem Aufruf der Sequenz (Risiko §6, zweiter Punkt), der Kontext des Vorlaufs ist `streamCtx` (abgeleitet aus `ctx`, ein Abbruch durch die WAL-Schwelle beendet den Vorlauf mit). |
| `internal/bootstrap/administration_startorder_internal_test.go` | neu (Plan-Nachzug) | Sechs Tests ohne Datenbank: Ordnung an einem `pending` stehenden `remove_transformation`-Antrag (Stream-Fake sieht Vermerk `applied` und Rohform), Ordnung für `enable`, Lesefehler hält den Start nicht an und der Stream-Ausgang bleibt der Rückgabewert, hängende Abfrage hält den Stream-Start an und der Abbruch des Kontexts beendet den Vorlauf, und zwei Tests über die Quelltext-Gestalt von `Run` (`stream.Run` nur als Argument der Sequenz; Abgleich und Worker-Start vor ihr). Grund der Quelltext-Tests: `Run` braucht die Datenbank, kein netzloser Lauf erreicht die Aufrufstelle. |
| `harness/README.md`, `docs/user/benutzerhandbuch.md` | prüfen | **Ergebnis: keine Änderung.** Keine der beiden Dateien beschreibt die Reihenfolge von Antrags-Verarbeitung und Stream-Start (Suchlauf unten, Zeilen 3 bis 4 und 7 bis 8); die Handbuch-Stellen zur Administrations-Goroutine (§4) sagen „verarbeitet offene Anträge … ohne Neustart“ und bleiben wahr. Der Abhilfe-Ablauf gehört zu `slice-transformationen-betriebsdoku` (Adresse). Keine Betreiber-Oberfläche berührt: `git diff --name-only ce229b8e -- internal/bootstrap/ tools/schema/ internal/adapters/driving/` trifft nur `wiring.go` und die Testdatei, keine Umgebungsvariable, keine SQL-Funktion, keinen Endpunkt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Reihenfolge von
Antrags-Verarbeitung und Stream-Start“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibungen der Reihenfolge im Code | Zeilen 1–2 im Block unten (Parent `ce229b8e`) | **Gefunden.** Symbole `runAdministration`/`stream.Run`/`processAdministrationRequests` in `wiring.go`: Parent 15, Diff 18. Die drei Doc-Kommentare (`Run`-Kopf: beschreibt keine Reihenfolge und bleibt; Kommentar an der Goroutine, Kommentar an `processAdministrationRequests`) sind gelesen: der Kommentar an der Goroutine nennt jetzt den Vorlauf und den Ort der Sequenz; `runStreamAfterAdministrationPass` trägt die Zusage der Sequenz. **Nichtgefunden:** kein Kommentar des Parent nennt die Reihenfolge von Goroutinen-Start und `stream.Run` (Zeile 5 im Block unten trifft am Parent 0). | Kommentar an der Goroutine im Ist-Ton, neuer Doc-Kommentar |
| Beschreibungen in Doku | Zeilen 3–4 im Block unten | **Gefunden.** `Administrations-Goroutine`/`Antrags-Queue` in `docs/user`, `harness`, `spec`: Parent 10, Diff 10 (unverändert). Gelesen: Handbuch §4 (zwei Stellen) sagt „verarbeitet offene Anträge (`LISTEN`/`NOTIFY`-Weckung mit periodischem Fallback-Poll)“, die Zeilen zu Rollen und Umgebungsvariable nennen die Queue, `SPEC-019` beschreibt die Queue. **Nichtgefunden:** keine Stelle in `docs/user`, `harness`, `spec` beschreibt die Reihenfolge von Antrags-Verarbeitung und Stream-Start (Zeilen 7–8 suchen beide Schreibweisen des Worts Startreihenfolge und treffen in Träger-Dokumenten ausschließlich Pläne der Welle und Code). | keiner; die Prozedur „Abhilfe nach einer nicht anwendbaren Regel“ (mit der Aussage zur Reihenfolge) trägt `slice-transformationen-betriebsdoku` (Adresse, Frist: dessen Closure) |
| Kommentare, die eine Reihenfolge behaupten | Zeilen 5–6 im Block unten | **Gefunden.** „vor“ oder „nach“ vor `stream.Run` in `internal` (Zeilen 5–6): Parent 0, Diff 1 (`wiring.go`, Kommentar an der Goroutine: „steht mit ihm vor `stream.Run`“). **Nichtgefunden:** kein weiterer Kommentar behauptet eine Reihenfolge zu `stream.Run`. Die eine Behauptung trägt `TestRunCallsTheStreamOnlyThroughTheOrderedSequence` (Aufrufstelle) und die Tests der Sequenz. | keiner |
| Startpfad des Backfill-Workers und des Start-Abgleichs (aus `slice-backfill-sql-administration`) | Lesen von `Run` in `internal/bootstrap/wiring.go` am Start | **Gefunden.** Die Folge in `Run` ist Bindungsaufbau (`activatedTableBindings`, `stream.BindCapture`) → `reconcileBackfillRuns` → Worker-Start → Sequenz (Vorlauf, Goroutine, `stream.Run`); der Abgleich `running` → `interrupted` und der Worker-Start stehen vor dem Vorlauf. **Nichtgefunden:** kein Pfad, auf dem der Vorlauf vor dem Abgleich läuft; das Wecksignal `backfillWake` besteht seit dem Anlegen vor dem Worker-Start (Kapazität 1, verschmelzend). Gebunden durch `TestRunReconcilesAndStartsTheBackfillWorkerBeforeTheAdministrationPass`. | keiner |
| Andere Pläne, die die Ordnung als Träger nennen | Zeilen 7–8 im Block unten | **Gefunden.** Zeilen 7–8: Parent 7 (vier in `welle-transformationen`, zwei in `slice-transformationen-e2e-abhilfe`, einer der Code-Kommentar zum Backfill-Start), Diff 8 (dazu der Kopfkommentar der neuen Testdatei). `slice-transformationen-e2e-abhilfe` (§1, §2, §4, §6) nennt „den Ordnungs-Test aus `start-reihenfolge`“ und die Startreihenfolge als Voraussetzung; `welle-transformationen` §3/§4/§5 beschreibt den Slice. Beide Aussagen bleiben wahr. **Nichtgefunden:** kein Plan behauptet eine Beschaffenheit der Ordnung, die der Slice ändert. Ein Träger außerhalb der Suche (ausgenommen `docs/plan/adr`): [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) führt die Abhilfe im gescheiterten Prozess als „erwartet, nicht am Code belegt“ (`git grep -n 'nicht am Code belegt' -- docs/plan/adr`: 1 Treffer in `ADR-0112`, Parent und Diff) — die Aussage ist eine `Accepted` ADR und bleibt unberührt (`AGENTS.md` §3.5); ihren Beleg trägt `slice-transformationen-e2e-abhilfe`. | Meldung an `e2e-abhilfe` (Fremddatei): die Tests tragen die Namen `TestRunStreamAfterAdministrationPass…` und `TestRun…` in `internal/bootstrap/administration_startorder_internal_test.go` (Frist: Start von `e2e-abhilfe`) |

```suchlauf
ce229b8e 15 -n -E 'runAdministration|stream\.Run|processAdministrationRequests' -- internal/bootstrap/wiring.go
diff 18 -n -E 'runAdministration|stream\.Run|processAdministrationRequests' -- internal/bootstrap/wiring.go
ce229b8e 10 -n -E 'Administrations-Goroutine|Antrags-Queue' -- docs/user harness spec
diff 10 -n -E 'Administrations-Goroutine|Antrags-Queue' -- docs/user harness spec
ce229b8e 0 -n -E 'vor .?stream\.Run|nach .?stream\.Run' -- internal
diff 1 -n -E 'vor .?stream\.Run|nach .?stream\.Run' -- internal
ce229b8e 7 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/adr'
diff 8 -n -E 'Startreihenfolge|Start-Reihenfolge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/adr'
ce229b8e 1 -n -F 'nicht am Code belegt' -- docs/plan/adr
diff 1 -n -F 'nicht am Code belegt' -- docs/plan/adr
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
  im Vorlauf hält den Start an). *Erwartet, zu belegen durch:* der Kontext des
  Vorlaufs trägt dieselbe Abbruch-Semantik wie die Goroutine; ein Test mit
  blockierendem Fake belegt, dass ein Abbruch des Kontexts den Vorlauf beendet.
  **Ausgang:** *(bei Closure)*
- **Ein `backfill`-Antrag im Vorlauf trifft einen Zustand ohne Abgleich oder
  Worker.** Der Backfill-Zweig nimmt den Antrag an (Run-Zeile `queued` und Antrag
  `applied` in einer Transaktion, dauerhaft, `slice-backfill-sql-administration`)
  und weckt den Worker; der Worker liest die `queued`-Zeilen beim Start, nach dem
  Abgleich `running` → `interrupted`. Der Vorlauf darf den Zweig nicht in einen
  Zustand bringen, in dem der Abgleich noch aussteht — ein `running`-Run einer
  früheren Prozessinstanz ließe die Prüfung „kein aktiver Run" den Antrag als
  `failed` enden (erwartet, hergeleitet aus der Annahme-Regel von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1). *Erwartet, zu belegen durch:* Lesen von `Run` am
  Start (Suchlauf, Zeile 4) und ein Test mit Fake-Worker.
  **Ausgang:** *(bei Closure)*
- **Der Vorlauf ändert das Verhalten für `enable`/`disable`**: eine Bindung
  entsteht vor dem Stream statt danach. Das ist die gewollte Eigenschaft; ein
  Test belegt die Reihenfolge für eine `enable`-Anfrage. **Ausgang:** *(bei
  Closure)*
- **Ein Kommentar behauptet mehr, als der Code trägt** — etwa, dass jede Race
  beseitigt sei (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`,
  offen, 2×): der Kommentar sagt nur zu, dass der Vorlauf vor dem Stream läuft.
  *Erwartet, zu belegen durch:* Review. **Ausgang:** *(bei Closure)*
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
