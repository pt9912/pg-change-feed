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

**Verantwortlich:** — (noch nicht priorisiert).

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
Synchronisation (gelesen am 2026-09-23 an `internal/bootstrap/wiring.go`,
Funktionen `Run` und `runAdministration`) — die Reihenfolge ist Zufall;
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

- [ ] Der Vorlauf steht: `Run` verarbeitet die offenen Anträge (ein synchroner
      Durchlauf von `processAdministrationRequests` über dieselben
      `administrationDeps` wie die Goroutine) **vor** `stream.Run`; die
      Goroutine läuft danach unverändert weiter; ein Lesefehler im Vorlauf
      blockiert den Start nicht (best effort, unverändert zum Dauerbetrieb,
      [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)),
      sondern wird protokolliert; die Verarbeitung eines Antrags im Vorlauf
      trägt die Bindung im `Assembler` nach, bevor der Stream startet. *Zu
      belegen durch:* `make test` (Whitebox in `internal/bootstrap`, Fakes für
      Antrags-Port und Assembler).
- [ ] Die Ordnung ist ohne Datenbank prüfbar und an ihre Eingabe gebunden: die
      Sequenz „Vorlauf, dann Stream“ steht an einer Stelle, die ein Test mit
      Fakes aufruft (kleine Extraktion aus `Run`, kein Umbau der Verdrahtung);
      ein Test belegt, dass ein `pending` stehender
      `remove_transformation`-Antrag verarbeitet ist, wenn der
      Stream-Start-Fake aufgerufen wird, und färbt sich rot, wenn die
      Reihenfolge vertauscht wird (Mutation). *Zu belegen durch:* `make test`
      mit Race-Detector und die Mutation im Bericht.
- [ ] Der Dauerbetrieb bleibt unverändert: `make test-integration` bleibt real
      grün (bestehende Belege der SQL-Administration einschließlich des
      Live-Reload-Belegs ohne `docker restart`); `make coverage-gate` grün. *Zu
      belegen durch:* ein realer, grüner `make test-integration`-Lauf.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt — keine Betreiber-Oberfläche; die Aussage zur
      Startreihenfolge steht, sofern das Handbuch sie trägt (Suchlauf), im
      Handbuch-Abschnitt von `slice-transformationen-betriebsdoku`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
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
| `internal/bootstrap/wiring.go` (+ Whitebox-Test) | update | Vorlauf-Durchlauf vor `stream.Run`; die Sequenz an einer testbaren Stelle; `administrationDeps` wird vor dem Goroutinen-Start gebildet. |
| `harness/README.md`, `docs/user/benutzerhandbuch.md` | prüfen | Aussagen über die Startreihenfolge und über „Anträge werden verarbeitet“ (Suchlauf) nachziehen bzw. an `betriebsdoku` melden. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Reihenfolge von
Antrags-Verarbeitung und Stream-Start“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibungen der Reihenfolge im Code | `grep -n 'runAdministration\|stream.Run\|processAdministrationRequests' internal/bootstrap/wiring.go` | *(Implementer trägt ein)* | Doc-Kommentare an `Run`, `runAdministration` und `processAdministrationRequests` beschreiben den Vorlauf |
| Beschreibungen in Doku | `grep -rn 'Administrations-Goroutine\|Antrags-Queue' docs/user harness spec` | *(Implementer trägt ein)* | Aussagen über den Start nachziehen; Handbuch-Stellen an `betriebsdoku` melden |
| Kommentare, die eine Reihenfolge behaupten | `grep -rn 'vor stream.Run\|nach stream.Run' internal --include=*.go` | *(Implementer trägt ein)* | jede Behauptung ist durch den Test dieses Slice getragen oder wird gestrichen (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, offen, 2×) |
| Startpfad des Backfill-Workers und des Start-Abgleichs (aus `slice-backfill-sql-administration`) | Lesen von `Run` in `internal/bootstrap/wiring.go` | *(Implementer trägt ein)* | Worker und Abgleich sind vor dem Vorlauf gestartet, sonst verarbeitet der Vorlauf einen `backfill`-Antrag gegen einen fehlenden Worker (Risiko §6) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-e2e-wirkung`
in `done/` liegt (Reihenfolge der Welle) und
`slice-backfill-sql-administration` in `done/` liegt (beide Slices ändern `Run`
im Bereich der Administrations-Goroutine; der Backfill-Zweig und der
Start-Abgleich stehen dann fest) und kein anderer Slice in `in-progress/` liegt
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  kleiner Zug im Startpfad; sprengt die Extraktion der testbaren Stelle den
  Umfang, ist der abtrennbare Teil der Test der Ordnung (zweiter Liefer-Punkt)
  als eigener Slice.
- `in-progress` → `open` (blockiert): falls `Run` ohne Umbau, der über die
  kleine Extraktion hinausgeht, keine Stelle für einen ordnungsprüfenden Test
  bietet (dann Architect-Frage zum Zuschnitt der Composition Root) oder falls
  der Vorlauf mit einem `backfill`-Antrag kollidiert, den der Worker noch nicht
  annimmt.

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
- **Ein `backfill`-Antrag im Vorlauf trifft keinen Worker.** Der Backfill-Zweig
  übergibt an einen Worker und legt eine `queued`-Run-Zeile an (dauerhaft,
  `slice-backfill-sql-administration`); der Vorlauf darf den Zweig nicht in
  einen Zustand bringen, in dem der Worker fehlt. *Erwartet, zu belegen durch:*
  Lesen von `Run` am Start (Suchlauf, Zeile 4) und ein Test mit Fake-Worker.
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
