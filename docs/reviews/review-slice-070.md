# Review-Report: slice-070 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-070`, Diff `6f9e9d8..f5e00d0` (= `HEAD`) — vier Commits:
`78d697f` (Publish-Kette + Tests, vier Dateien), `efd8d87` (Digest-Stempel),
`48733d7` (Kommentar-Fixrunde), `f5e00d0` (Digest-Neustempel). Netto fünf
Dateien (+362/−17). Der Elter-Commit `6f9e9d8` ist der reine
`next → in-progress`-Move und trägt **keine** Fremd-Commits;
`git diff --name-only 6f9e9d8..HEAD` nennt genau die fünf erwarteten Pfade.

**Skill:** `.harness/skills/reviewer.md` @ `f07ba3d` (letzte Schärfung
2026-09-13, unverändert seitdem).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-14.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-070-grpc-capture-integration.md`
  vollständig (§1–§8, inkl. Plan-Korrektur und Implementer-Nachzug in §3)
- [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  vollständig — §Kontext, §Entscheidung (Festlegungen 1–4), §Verglichene
  Alternativen (Option B und C), §Konsequenzen, §Fitness Function,
  §Re-Evaluierungs-Trigger
- [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) vollständig —
  insbesondere Teilfrage 2 (Port-/Adapter-Design, Aufrufkette), Teilfrage 3,
  Teilfrage 6, §Fitness Function, §Slice-Schnitt-Empfehlung
- [`ADR-0011`](../plan/adr/0011-persist-before-ack.md),
  [`ADR-0027`](../plan/adr/0027-capture-application-service.md),
  [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)
  Punkt 1/4, [`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md),
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
- [`LH-FA-SST-008`](../../spec/lastenheft.md) (Happy Path/Boundary),
  [`SPEC-020`](../../spec/pflichtenheft.md) (Zeilen *Zustellgarantie*,
  *Erzeuger-Blockade*)
- `AGENTS.md` §3.1/§3.2/§3.5/§3.7/§3.9/§4/§5/§6; `harness/conventions.md`
  (`MR-000`); `.a-check.yml`
- [`review-slice-069`](review-slice-069.md) (beide Läufe),
  [`verify-slice-069`](verify-slice-069.md), der Vorgänger-Plan
  `docs/plan/planning/done/slice-069-grpc-streaming-adapter-grundgeruest.md`,
  `internal/adapters/driven/grpcstream/broadcaster.go`,
  `internal/application/port/outbound/changestream.go`
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/`:
  `slice-chronik-in-code-kommentar` (verkörpert), `adapter-fehler-ausgang`
  (2×, weiter offen), `report-nackte-id-ohne-link` (verkörpert),
  `commit-traceability-kein-vorab-hook`

---

## Findings

### F-1 — `ADR-0066`s dritte Fitness-Function-Zeile fordert eine Zeit-Isolation, die dieselbe ADR an der Aufrufstelle verbietet

- `kategorie`: HIGH
- `quelle`: [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  §Fitness Function (3. Zeile) gegen §Entscheidung und §Verglichene
  Alternativen Option B · `AGENTS.md` §3.5 (Accepted-ADR wird nicht inhaltlich
  überschrieben)
- `pfad`: `docs/plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md:276`
  gegen `:117-123` und `:210` ·
  `internal/application/usecase/capture/service.go:148-155` ·
  `internal/application/usecase/capture/service_test.go:648-689` ·
  `docs/plan/planning/in-progress/slice-070-grpc-capture-integration.md:109-111`,
  `:143-161`
- `befund`: Die Zeile verlangt einen Testfall, der deckt, dass ein
  `ChangeStreamPort` mit **nicht zurückkehrendem** `Publish` `Capture()` nicht
  anhält; dieselbe ADR entscheidet, dass der Aufrufer `Publish` **synchron**
  ohne caller-seitige Goroutine/Deadline ruft, und verwirft genau diese
  Maßnahme als Option B. Bei einem synchronen Aufruf hält ein nicht
  zurückkehrender `Publish` `Capture()` zwangsläufig an — der geforderte Test
  wäre nur mit der in derselben ADR verworfenen Maßnahme grün; die
  Implementierung trägt die Reduktion in §3 und als Fußnote der DoD-Zeile, die
  `Accepted`-ADR-Zeile selbst bleibt unverändert.
- `verifizierbar`: nein — Textwiderspruch innerhalb einer Datei, kein
  Gate-Gegenstand; die **Unmöglichkeit** der Klammer-Hälfte an dieser Schicht
  ist dagegen maschinell belegbar (§Mutations- und Probenläufe: `make a-check`
  mit einer Probe-Datei, Exit 2)
- `klasse`: „Fitness-Function-Zeile gegen die eigene Entscheidung"

### F-2 — Der neue `wiring.go`-Kommentar behauptet einen Fehlerpfad, den der Code nicht trägt

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Hard Rule: „Ein Kommentar beschreibt, was da
  ist")
- `pfad`: `internal/bootstrap/wiring.go:557-560`, gegen `:685`, `:713` und
  `:788`
- `befund`: Der Kommentar sagt zu, „der gRPC-Server-Startfehler erreicht das
  Ergebnis von `Run` weiter unten auf demselben Pfad wie jeder andere
  Adapter-Startfehler" — jeder Adapter-Startfehler wird jedoch in seiner
  Goroutine per `log.Error` verworfen (`:685` HTTP, `:713` gRPC, kein
  Fehlerkanal, kein `errgroup`), und `Run` gibt ausschließlich
  `mergeStreamAndWALFaultOutcome(streamErr, &walFault)` zurück (`:788`). Eine
  fehlkonfigurierte `CDC_GRPC_ADDR` beendet den Lauf danach nicht; der
  Kommentar legt das Gegenteil nahe.
- `verifizierbar`: nein — Kommentar-Wahrheit, kein Gate-Gegenstand (der
  `make gates`-Lauf ist davon unberührt grün)
- `klasse`: „Kommentar behauptet einen Fehlerpfad, den der Code nicht trägt"

### F-3 — Der NATS-Kommentar steht nach dem Diff nicht mehr bei seinem Gegenstand

- `kategorie`: LOW
- `quelle`: Maintainability (`AGENTS.md` §3.7)
- `pfad`: `internal/bootstrap/wiring.go:541-550`, gegen `:567`
- `befund`: Der Block, der die NATS-Verbindung als Start-Vorbedingung erklärt
  („ein Betreiber, der das Feature einschaltet … soll das beim Start
  bemerken"), steht jetzt 17 Zeilen über `if cfg.NatsURL != ""`; dazwischen
  liegen der eingefügte Broadcaster-Block (`:551-560`) und der Aufbau von
  `captureOpts`. Die Zuordnung Kommentar → Anweisung ist an dieser Stelle
  verloren.
- `verifizierbar`: nein
- `klasse`: „Kommentar nicht mehr bei seinem Gegenstand"

### F-4 — `erwarteteZustellung` zählt zurückgekehrte `Publish`-Aufrufe, nicht zugestellte Changes

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/application/usecase/capture/service_test.go:585-593`,
  `:625-627`, gegen `:444-464`
- `befund`: Im Fall „nicht lesender Empfänger" erwartet der Test den Wert `2`,
  obwohl der Doppel die zweite Change wegen voller Warteschlange (Kapazität 1)
  verwirft — gezählt wird die Zahl der ohne Fehler **zurückgekehrten** Aufrufe
  (`fakeStream.published`). Die Fehlermeldung benennt das korrekt
  („erfolgreiche Aufrufe"), der Feldname legt dagegen Zustellung nahe.
- `verifizierbar`: nein
- `klasse`: „Test-Feldname verwechselt Aufruf-Rückkehr mit Zustellung"

### F-5 — Zwei Digest-Stempel für einen Zug, der erste nach drei Minuten überholt

- `kategorie`: INFO
- `quelle`: [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
  (Digest-Beleg ist lauf-/builder-gebunden)
- `pfad`: `harness/image-hash.txt`; Commits `efd8d87`, `48733d7`, `f5e00d0`
- `befund`: `efd8d87` stempelte den Beleg eines echten `make image`-Laufs,
  danach änderte `48733d7` eine Build-Kontext-Datei (`service.go`, doppelter
  Kommentarblock), `f5e00d0` stempelte den zweiten Lauf — regelkonform
  ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) verlangt genau das).
  Der finale Digest ist als echter Lauf-Beleg verifizierbar: er stimmt mit der
  ID des lokalen `ghcr.io/pt9912/pg-change-feed:dev` überein
  (`created=2026-09-14T20:38:38+02:00`, also zwischen `48733d7` 20:38:07 und
  `f5e00d0` 20:39:33).
- `verifizierbar`: ja — `docker image inspect ghcr.io/pt9912/pg-change-feed:dev`
  gegen `harness/image-hash.txt`
- `klasse`: „Digest-Stempel durch nachgereichte Kommentar-Korrektur verdoppelt"

### F-6 — Die Begründung des verworfenen `tx.Changes()`-Fehlers steht wortgleich zweimal in derselben Datei

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/application/usecase/capture/service.go:145-147`, gegen
  `:169-171`
- `befund`: Derselbe dreizeilige Satz („Der Fehler von `tx.Changes()` … deshalb
  wird er hier verworfen") steht über dem neuen Publish-Block und —
  vorbestehend — über `distinctTables`. Beide Stellen sind sachlich gedeckt
  (der `CommitPosition`-Wächter `:110-113` liegt davor, `Changes()` scheitert
  nur an einer offenen Transaktion); die zweite Fassung ist eine wörtliche
  Kopie.
- `verifizierbar`: nein
- `klasse`: „wörtliche Begründungs-Doppelung im selben Produktionspfad"

### F-7 — Die Klammer-Hälfte der Fitness-Function-Zeile ist an dieser Schicht nicht gegen den realen Port prüfbar

- `kategorie`: INFO
- `quelle`: `.a-check.yml:12-21` (keine `app → adapters`-Kante) ·
  [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  §Fitness Function (3. Zeile)
- `pfad`: `internal/application/usecase/capture/service_test.go:444-464`, gegen
  `.a-check.yml:12-21`
- `befund`: Der Capture-Test trägt die Zustellsemantik in einem Doppel nach
  (`fakeStream` sendet nicht-blockierend in eine Warteschlange mit Kapazität 1
  und verwirft den Überlauf), weil
  `internal/application/usecase/capture` kein Adapterpaket importieren darf —
  eine Probe-Datei mit `grpcstream`-Import färbt `make a-check` rot
  (`app-impurity`, Exit 2). Der Test belegt damit, dass **die Aufrufstelle**
  kein eigenes Warten hinzufügt, nicht dass der reale `Broadcaster` nicht
  blockiert; diesen Beleg trägt `TestNichtLesenderAbonnentHaeltPublishNichtAn`
  im Paket `internal/adapters/driven/grpcstream` (aus der `slice-069`-Fixrunde,
  im Review dort rot gesehen).
- `verifizierbar`: ja — `make a-check` mit einer Probe-Datei (hier real
  ausgeführt, Exit 2); die Regel selbst bleibt grün (`make gates` Exit 0)
- `klasse`: „FF-Hälfte an der Schichtgrenze nur gegen Port-Doppel prüfbar"

### F-8 — Während dieses Laufs erscheinen zwei fremde, ungetrackte Dateien im Arbeitsbaum

- `kategorie`: INFO
- `quelle`: Maintainability (Rollen-Trennung, Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Rollen-Sequenz für einen Slice)
- `pfad`: `docs/reviews/…`-Commit dieses Laufs:
  `docs/plan/planning/observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/`
  (`observation.md`, `state.md`, `evidence/slice-059.md`, `-066.md`, `-069.md`)
  und `docs/plan/planning/open/slice-077-handbuch-betreiber-stand.md`
- `befund`: Beide Pfade sind während dieses Reviews entstanden (der
  Arbeitsbaum war beim Start und bis nach den Mutationsläufen leer) und sind
  Planer-Arbeit — ein neuer Registereintrag und ein `slice-077`-Stub, der noch
  die `<Platzhalter>` der Vorlage trägt. Sie sind nicht Gegenstand dieses
  Reports und werden von ihm **nicht** mitcommittet (der Commit ist per
  Pfad-Argument auf `docs/reviews/review-slice-070.md` beschränkt).
- `verifizierbar`: ja — `git status --porcelain` (beide Pfade `??`)
- `klasse`: „Fremdänderung im Arbeitsbaum während des Reviews"

---

## Zum ADR-Widerspruch — Kernfrage dieses Auftrags

**Verdikt: legitime Abweichung — die Fitness-Function-Zeile ist der Defekt,
nicht der Code. Die Implementierung ist vollständig gegenüber dem, was
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
**entschieden** hat.** Der Widerspruch liegt **innerhalb von `ADR-0066`**; die
Code-Seite folgt der einen Hälfte und muss deshalb der anderen nicht folgen.

Die drei betroffenen Stellen im Wortlaut:

1. §Entscheidung (`:119-123`): „Der Erzeuger (`CaptureService`) hält nie auf
   einen Abonnenten an. Der **Aufrufer bleibt unverändert** — `slice-070` ruft
   `Publish` weiter synchron in der Best-Effort-Kette und isoliert den Fehler;
   eine caller-seitige Goroutine entfällt."
2. §Verglichene Alternativen, Option B (`:210`, **verworfen**): „**B** — nur der
   Aufrufer: `slice-070` ruft `Publish` nicht-blockierend (eigene Goroutine je
   Change oder `Publish` mit kurz-Deadline-`ctx`)" · Contra dort: „… eine
   Deadline-`ctx` verlagert den Halt nur in ein Zeitfenster und legt **jedem**
   Publish eine Wartezeit auf …".
3. §Fitness Function, 3. Zeile (`:276`): „Ein `ChangeStreamPort`, dessen
   `Publish` nicht zurückkehrt (oder dessen Subscriber nicht liest), darf
   `Capture()` nicht anhalten; der Fehler-/Zeitisolations-Regressionstest aus
   `slice-070` deckt **beide Hälften**."

**Der Widerspruch, in einer Zeile:** Satz 1 schreibt den synchronen Aufruf ohne
Zeit-Isolation fest, Satz 3 verlangt für den Fall „`Publish` kehrt nicht
zurück" genau die Zeit-Isolation, die Satz 2 als verworfene Option B benennt.
Mit einem synchronen Aufruf im Capture-Goroutine-Kontext gibt es keinen
Mechanismus, der `Capture()` zurückkehren ließe, während `Publish` hängt — die
einzigen Kandidaten (eigene Goroutine je Change, Deadline-`ctx`) sind exakt die
verworfene Option B. Die drei Sätze sind nicht gleichzeitig erfüllbar; einer
muss fallen, und gefallen ist der maschinell prüfbare — die Fitness Function.

**Der Beleg liegt im ausgelieferten Test selbst.**
`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`
(`internal/application/usecase/capture/service_test.go:656-689`) parkt den
`Capture()`-Aufruf über `fakeBlockingStream` (`:470-482`) **in** `Publish`,
liest in dieser Parkposition Store- und ACK-Stand ab und gibt ihn erst mit
`close(stream.freigabe)` (`:685`) wieder frei. Der Test zeigt damit exakt die
Gegenrichtung der FF-Zeile: `Capture()` hält an, solange `Publish` nicht
zurückkehrt — und die Capture-kritische Kette `Persist → ACK` ist trotzdem, nach
`ADR-0060` Teilfrage 2 und
[`ADR-0011`](../plan/adr/0011-persist-before-ack.md), vollständig. Die
Implementer-Begründung trägt; sie ist hier unabhängig nachvollzogen, nicht
übernommen.

**Zwei Gründe, warum die FF-Zeile auch in ihrer zweiten Lesart nicht zu retten
ist:** Die Klammer „oder dessen Subscriber nicht liest" könnte man so verstehen,
dass der Capture-Test den realen `Broadcaster` anbinden soll. Das verbietet die
Schichtregel: `.a-check.yml` führt keine `app → adapters`-Kante, und `a-check`
scannt `**/*.go` inklusive `_test.go` — eine hier real gestellte Probe-Datei
(`grpcstream`-Import im Paket `capture_test`) endet mit `app-impurity`, Exit 2.
Ein Test gegen den realen Port ist an dieser Stelle strukturell unmöglich; wo
die Eigenschaft real belegt wird, ist das Paket
`internal/adapters/driven/grpcstream`. `ADR-0066`s FF-Zeile weist der
Capture-Schicht damit eine Prüfung zu, die diese Schicht nicht tragen kann.

**Was der Architect-Zug ohne Rückfrage vorfindet** — drei Verdikte, alle nach
`AGENTS.md` §3.5 (in-place ist keine Option, `ADR-0066` ist `Accepted`):

1. *FF-Zeile ist falsch, Entscheidung ist richtig* (der Befund dieses Reports):
   Folge-ADR mit `Supersedes`-Vermerk auf `ADR-0066`, deren FF-Zeile den Anteil der Schichten
   trennt — Port-Ebene: „ein registrierter, nicht lesender Abonnent hält
   `Publish` nicht an" (Paket `internal/adapters/driven/grpcstream`, real
   getragen); Aufrufer-Ebene: „die Capture-kritische Kette ist vollständig,
   bevor der erste Change das Haus verlässt, und die Aufrufstelle fügt kein
   eigenes Warten hinzu" (Paket `internal/application/usecase/capture`, gegen
   einen vertragstreuen Doppel). **Keine Code-Änderung, kein Implementer-Rückweg.**
2. *Zeit-Isolation gehört doch an den Aufrufer*: Folge-ADR, die §Entscheidung
   und Option B dieser ADR `supersedes`; Folge: Implementer-Runde
   (Goroutine/Deadline **plus** Frist-Test) — und die in Option B notierten
   Contras (unbeschränkte Goroutine je Change, Wartezeit auf jedem Publish)
   gelten dann unverändert.
3. *Zeile bleibt, Reduktion bleibt* ist kein Verdikt: eine Accepted-ADR-Aussage
   ohne Ausgang stehen zu lassen ist die stille Korrektur, die `AGENTS.md` §3.5
   und Modul 8 §Konflikt-Pfad ausschließen.

Der minimale Träger ist damit ein Architect-Artefakt (Folge-ADR oder Verdikt,
wie im Repo bei `ADR-0059`/`ADR-0065`/`ADR-0062` und beim Verdikt
`architect-verdict-publish-blockiert-capture-pfad.md` bereits praktiziert).
Nichts davon ist Implementer-Arbeit — deshalb steht in der Summary auch bei zwei
HIGH nur **ein** Rückweg-Kandidat am Code (F-2).

---

## Negativbefunde

- **geprüft, ohne Befund: der Publish-Aufruf liegt real nach `ACK Source` und
  nach dem Notify-Schritt (`ADR-0060` Teilfrage 2).** `service.go:118-120` ACK,
  `:129-135` Notify, `:148-155` Stream-Publish; der Happy-Path-Test pinnt die
  Ordnung als Ereignisliste (`service_test.go:504`:
  `persist → ack → notify → stream:t-1-1 → stream:t-1-2`). Selbst gegengeprüft:
  Mutation (ii) (Block vor `ack.Acknowledge`) färbt genau die drei
  ordnungstragenden Tests rot (Exit 2, Ereignisse sichtbar verschoben).
- **geprüft, ohne Befund: ein `Publish` je `Change` aus `tx.Changes()`, in
  Reihenfolge, ohne Dedup nach Tabelle.** `service.go:149-154` iteriert die
  Kopie aus `Changes()` (`transaction.go:98-103` liefert
  `append([]Change(nil), t.changes...)` in Anhang-Reihenfolge); zwei Changes
  derselben Tabelle ergeben zwei Aufrufe (`service_test.go:491-514`), anders als
  beim tabellen-granularen Notify. Mutation (iii) (`break` nach dem ersten
  Change) färbt die Tests rot (Exit 2) — die Schleife ist lasttragend, keine
  Deko.
- **geprüft, ohne Befund: Fehlerisolation und `s.stream == nil`.**
  `service.go:151-153` fängt jeden `Publish`-Fehler in eine `LogPort.Warnung`,
  der Rückgabewert von `Capture()` bleibt `CaptureResult{Acknowledged: position},
  nil` (`:156`); der Test `TestCaptureSucceedsDespiteFailingStream`
  (`:549-575`) belegt es einschließlich „der Aufruf läuft über alle Changes
  weiter". Mutation (i) (`Warn` → `return CaptureResult{}, err`) färbt genau
  diesen Test rot (Exit 2). Ohne Port (`nil`) unterbleibt der Versuch
  vollständig (`:148`), belegt durch
  `TestCaptureWithoutStreamPortLeavesBehaviourUnchanged` (`:520-539`) —
  bit-identisches Bestandsverhalten für Aufrufer ohne diese Option.
- **geprüft, ohne Befund: `Persist`/`ACK` durch den neuen Schritt unberührt
  (`ADR-0011`/`ADR-0027`).** Der Diff ändert im Rumpf von `Capture()` keine
  bestehende Zeile: die `store.PersistTransaction`- und
  `ack.Acknowledge`-Aufrufe sowie der Notify-Block sind zeichengleich
  (`git show 78d697f -- internal/application/usecase/capture/service.go`), der
  neue Block steht danach. Das §6-Risiko „Publish versehentlich vor `ACK`" ist
  damit am Code **nicht** eingetreten — und wäre von drei Tests gefangen worden
  (Mutation (ii)).
- **geprüft, ohne Befund: Scope-Fidelity (Auftrags-Punkt 7).**
  `git diff --name-only 6f9e9d8..HEAD` nennt ausschließlich
  `docs/plan/planning/in-progress/slice-070-grpc-capture-integration.md`,
  `harness/image-hash.txt`,
  `internal/application/usecase/capture/service.go`,
  `internal/application/usecase/capture/service_test.go`,
  `internal/bootstrap/wiring.go`. **Kein** Eingriff unter
  `internal/adapters/driven/grpcstream/**` (Port-Seite bleibt `slice-069`s
  Gegenstand), **keine** Änderung an `.a-check.yml` (`ADR-0060` Teilfrage 5
  Punkt 2), keine Änderung an `ChangeNotificationPort`/`natsnotify`, kein
  `internal/adapters/driving/**`.
- **geprüft, ohne Befund: Kommentar-Hygiene der geänderten `.go`-Dateien
  (Auftrags-Punkt 5).** `grep -nE "slice-[0-9]+|welle-[0-9]+"` über `service.go`,
  `service_test.go`, `wiring.go`: **kein** Treffer (Klasse
  `BEO-PGC/slice-chronik-in-code-kommentar`). Der in `48733d7` entfernte
  doppelte Kommentarblock ist real weg; eine Suche nach doppelten
  Kommentarzeilen in `service.go`/`service_test.go` liefert keinen Treffer mehr
  (`wiring.go` hat drei vorbestehende Doppelungen, alle bei `:1236`+ — außerhalb
  des Diff-Bereichs `:541-560`/`:691-706`). F-6 betrifft die *inhaltliche*
  Wiederholung einer Begründung, nicht einen Block-Doppelgänger.
- **geprüft, ohne Befund: Traceability und ID-Schema.** Alle vier Betreffs
  tragen `LH-FA-SST-008` und `ADR-0060`, keiner ein `SPEC-*`/`ARC-*`;
  `make commit-traceability` meldet im Gate-Lauf 5 Commits, Betreffe ohne
  Struktur-ID (Exit 0). Die neuen Kennungen im Diff folgen `MR-000`; neu
  vergeben wird keine ID.
- **geprüft, ohne Befund: kein Suppression-Pfad, kein Docker-only-Verstoß.**
  `grep -rn "nolint\|noqa\|SuppressMessage"` über die geänderten Produktions-
  und Testdateien: kein Treffer (`AGENTS.md` §3.2). Kein Build-/Test-Betrieb
  außerhalb von `make` (`AGENTS.md` §3.1); der Digest-Beleg stammt aus einem
  realen `make image` (F-5).
- **geprüft, ohne Befund: der `CDC_GRPC_ADDR`-Zweig in `Run()` ohne Unit-Test
  ist konsistent mit dem Bestand (Auftrags-Punkt 6).** `wiring_test.go` trägt
  sieben Tests, alle gegen `ConfigFromEnv` (u. a.
  `TestConfigFromEnvGRPCAddrBleibtOptional`, `:155-171`); `Run()` wird in
  `make test` von **keinem** Unit-Test durchlaufen — auch nicht für den
  HTTP- oder NATS-Zweig. Der reale `bootstrap.Run`-Pfad liegt in
  `walretention_endtoend_test.go` und braucht PostgreSQL (Store-Tier, nicht
  `make test`). Der neue Zweig ist damit nicht schlechter gestellt als der
  Bestand; die Verdrahtungs-Identität (derselbe `*grpcstream.Broadcaster` für
  Server und CaptureService) erzwingt der Compiler über die Variable
  (`wiring.go:561-566` gegen `:706`). Kein Finding.
- **geprüft, ohne Befund: die Verdrahtung ist additiv und trägt beide Hälften
  korrekt.** `grpcBroadcaster` entsteht nur bei `cfg.GRPCAddr != ""`
  (`:563-566`), der `WithChangeStream`-Aufruf hängt an derselben Bedingung; ohne
  Adresse kein Broadcaster, kein Port, kein Publish-Schritt (`ADR-0060`
  Teilfrage 6). Der gRPC-Server liest aus **derselben** Instanz (`:706`), nicht
  aus einer zweiten — die Zustellkette Capture → Broadcaster → Server ist
  geschlossen.
- **geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` zu Recht unberührt.**
  Der Diff führt keine neue Umgebungsvariable ein (nur `internal/bootstrap/wiring.go`
  und die Capture-Schicht); `CDC_GRPC_ADDR` bleibt wie `CDC_HTTP_ADDR` außerhalb
  des Handbuchs. Kein HIGH der Klasse
  `BEO-PGC/handbuch-versionshistorie-uebersprungen` (die Datei ist nicht im
  Diff).
- **geprüft, ohne Befund: DoD-Häkchen-Trennung im Slice-Plan.** Der Plan-Diff
  (`78d697f`) zieht genau vier Zeilen auf `[x]` — Happy Path, Teilfrage 2 (mit
  der Fußnote zur FF-Hälfte), `make gates`, Doku-Update `wiring.go`-Kommentar —
  und lässt „Review durchgeführt", Closure-Notiz, Regel-Einträge, Risiko-Ausgänge
  und die drei Paarungen offen. Keine vorweggenommene Closure-Zeile, kein
  Übergreifen in eine andere Rolle; die Fußnote trägt die Reduktion sichtbar,
  nicht stillschweigend.
- **geprüft, ohne Befund: kein Falschzitat im Implementer-Nachzug (§3).** Die
  drei Einträge benennen die tatsächlich gebauten Tests (`:143-161`) und den
  Pfad, der die Schichtregel trägt (`internal/application/usecase/capture`
  importiert `grpcstream` nicht) — beides am Baum bestätigt.
- **geprüft, ohne Befund: `internal/application/usecase/capture/service.go`,
  `.../service_test.go`, `internal/bootstrap/wiring.go`,
  `harness/image-hash.txt`, der Slice-Plan** — über die oben benannten Punkte
  hinaus kein Befund.

---

## Sensor-Läufe und Mutationen (alle selbst ausgeführt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code danach in
einem **eigenen, ungeketteten** Schritt ermittelt (`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify OK; d-check 547 Dateien/0 Befunde (links/anchors/ids/matrix/versions/structure) + `commits`-Modul 0 Befunde; commit-traceability 5 Commits; a-check 0 Befunde (2 unverschichtete Werkzeug-Dateien, bekannte Hinweise); `coverage-gate: OK — 48.40 % ≥ 35 %` |
| `make test` | **0** | 28 Pakete `ok`, kein `FAIL`; `internal/application/usecase/capture` grün |
| `make a-check` + Probe-Datei (`grpcstream`-Import in `capture_test`) | **2** | `app-impurity` — Beleg für F-7; Probe zurückgenommen, `git status --porcelain` leer |
| `make test` + Mutation (i) (`Warn` → `return CaptureResult{}, err`) | **2** | rot genau an `TestCaptureSucceedsDespiteFailingStream` (`service_test.go:559`) |
| `make test` + Mutation (ii) (Publish-Block vor `ack.Acknowledge`) | **2** | rot an `TestCapturePublishesEachChangeOnceAfterAckAndNotify`, `TestCaptureSucceedsDespiteFailingStream`, `TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`; Log zeigt die verschobene Ordnung `persist, stream…, ack, notify` |
| `make test` + Mutation (iii) (`break` nach dem ersten Change) | **2** | rot an `TestCapturePublishesEachChangeOnceAfterAckAndNotify`, `TestCaptureSucceedsDespiteFailingStream` und beiden Unterfällen von `TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck` |

**Rücknahme und Blatt-Identität nach jeder Mutation** (`git checkout --`, dann
eigener Schritt): `git hash-object` von
`internal/application/usecase/capture/service.go` = `34917edb2ab0e9ec3e8207a19466296f74a72153`
= `git rev-parse HEAD:internal/application/usecase/capture/service.go`;
`git status --porcelain` leer. Die grünen Läufe oben gelten damit für die
committeten Bytes.

**Zum Nachweis-Wert der Mutationen.** Die drei vom Implementer genannten
Mutationen sind alle drei real sensibel — keine ist eine Dekoration, keine
greift ins Leere. Bemerkenswert: Mutation (ii) färbt auch
`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn` rot, den Test also,
der die Parkposition in `Publish` ausnutzt — die „kritische Kette steht vor dem
Publish"-Aussage ist damit ordnungs-, nicht nur anwesenheitsgetragen.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Fitness-Function-Zeile gegen die eigene
Entscheidung · Kommentar behauptet einen Fehlerpfad, den der Code nicht trägt ·
Kommentar nicht mehr bei seinem Gegenstand · Test-Feldname verwechselt
Aufruf-Rückkehr mit Zustellung · Digest-Stempel durch nachgereichte
Kommentar-Korrektur verdoppelt · wörtliche Begründungs-Doppelung im selben
Produktionspfad · FF-Hälfte an der Schichtgrenze nur gegen Port-Doppel prüfbar · Fremdänderung im Arbeitsbaum während des Reviews

## Verdikt

**Merge-blockierend:** ja — 2 HIGH, davon **genau einer am Code**: F-2 (eine
Kommentarzeile in `internal/bootstrap/wiring.go:557-560`, die einen
Fehlerpfad behauptet, den der Code nicht hat) braucht eine Rückkante
Reviewer → Implementer. **F-1 blockiert den Merge nicht, sondern die Closure:**
der Defekt liegt im `Accepted`-Text von
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
seine Behebung ist eine Architect-Entscheidung (Folge-ADR/Verdikt, Modul 8
§Konflikt-Pfad), und der Slice darf nicht still abschließen, solange die
FF-Zeile unrevidiert dasteht. F-3 bis F-7 sind ohne erwartete Rückkante; F-5/F-6
sind Hinweise, F-7 ist die für die Closure nützliche Schicht-Aussage.

**Auftrags-Kernfrage — Antwort in einem Satz:** legitime Abweichung; die
Implementierung ist gegenüber `ADR-0066`s Entscheidung vollständig, und die
dritte Fitness-Function-Zeile derselben ADR ist der Defekt, weil sie die
Zeit-Isolation verlangt, die die ADR an der Aufrufstelle ausdrücklich
ausschließt. Die ausführliche Fassung samt drei entscheidungsfähigen Verdikten
steht in §„Zum ADR-Widerspruch".

**DoD-Checkbox-Nachzug:** **kein Nachzug** — weil F-2 eine echte Fixrunde am
Implementer nach sich zieht, bleibt die Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan **offen** und wird regulär bei
Schritt 21 des Implementer-Workflows nachgezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde greift genau dann, wenn **keine** Fixrunde
kommt).

**Übergabe:** F-2 an den Implementer (eine Zeile, kein Verhalten). F-1 an die
**Architect-Rolle** (Folge-ADR mit `Supersedes`-Vermerk auf `ADR-0066` oder Verdikt; keine
Code-Arbeit) — die Formulierung ist rückfragefrei entscheidungsfähig gehalten.
Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff,
dieser Skill, dieses Modell, dieses Verdikt) und wird über Läufe hinweg nicht
wieder gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-11-*.md`;
anderes Prüf-Artefakt, anderer Eingabe-Kontext).


---

## Fixrunde (2. Lauf) — 2026-09-14

**Gegenstand:** Commit `37717c6` (`wiring.go`, `service.go`, `service_test.go`,
Slice-Plan-§3) — real gelesen, nicht aus der Commit-Message übernommen.
Zusätzliche Prüfgrundlage dieses Laufs: [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
vollständig (`20d643a`, samt Index-Zeile `docs/plan/adr/README.md:79-80`),
`f714a88` (Plan-Nachzug), `docs/plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md`
(gegen die Korrektur gelesen), `spec/pflichtenheft.md:339-361` (`SPEC-020`,
auf eine Rest-Hälfte geprüft).

**Eigene Sensor-Läufe** (Exit-Code je eigener, ungepipter, mechanisch
konditionierter Schritt, `AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify `v6.5.0` 54 Dateien OK; d-check 555 Dateien/0 Befunde (links/anchors/ids/matrix/versions/structure), `commits`-Modul 0 Befunde; commit-traceability 5 Commits; a-check 0 Befunde; `coverage-gate: OK — 48.40 % ≥ 35 %` |
| `make test` (Stand `37717c6`) | **0** | kein `FAIL`; `capture`, `grpcstream`, `bootstrap` grün |
| `make test` + Mutation A (`erwarteteErfolgreicheAufrufe: 2 → 1`, Fall „nicht lesender Empfänger") | **2** | rot genau an `TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck/nicht_lesender_Empfänger` (`service_test.go:626`, „Stream trägt 2 erfolgreiche Aufrufe, wollen 1") |
| `make test` + Mutation B (`log.Warn` → `return CaptureResult{}, err` im Publish-Block) | **2** | rot genau an `TestCaptureSucceedsDespiteFailingStream` (`service_test.go:559`) |

**Rücknahme und Blatt-Identität:** nach `git checkout --` stimmt
`git hash-object` beider Dateien mit `git rev-parse HEAD:<pfad>` überein
(`service_test.go` `e1d86977…`, `service.go` `b6a912c0…`),
`git status --porcelain` ist leer.

### Verdikt je Finding

| Finding | Verdikt | Beleg (eigene Prüfung am Code/Text) |
|---|---|---|
| F-1 (HIGH) | **behoben — durch `ADR-0067`, nicht am Code** | `0067-…md:114-176` (Entscheidung + Ersatztext), `:252-257` (neue Fitness Function), `README.md:79-80`; kein Code-Eingriff, keiner nötig |
| F-2 (HIGH) | **behoben** | `wiring.go:547-550` gegen `:563-566`, `:685`, `:713`, `:788` |
| F-3 (LOW) | **behoben** | `wiring.go:553-562` steht jetzt unmittelbar über `if cfg.NatsURL != ""` (`:563`) |
| F-4 (LOW) | **behoben** | `service_test.go:586-593`, `:625`; Name trifft die gezählte Größe |
| F-5 (INFO) | **Einordnung geteilt** — kein „Repair" nötig | Finaler Digest bleibt der echte Lauf-Beleg (`ADR-0044`), der Zwischenstempel ist Historie in `git`; unverändert gegenüber dem ersten Lauf |
| F-6 (INFO) | **behoben** | `service.go:157-164` (Helfer), `:175` (`distinctTables`), `:146-149`; die Begründung steht genau einmal |
| F-7 (INFO) | **beantwortet durch `ADR-0067`** | `0067-…md:83-92` (Klammer-Hälfte mit dem `app-impurity`-Beleg), `:167-168`, `:234-237` (Doppel statt realer Port als Schichtgrenze, kein Mangel) |
| F-8 (INFO) | **erledigt, kein Handlungsbedarf** | Die zwei Pfade sind inzwischen von der Planer-Rolle committet (`8423b28`, `68d2ebd`); `git status --porcelain` ist jetzt leer |

**F-2 — behoben, am Code nachgeprüft.** Der neue Satz
(`wiring.go:547-550`) lautet: „Die Bindung ist keine Start-Vorbedingung: der
gRPC-Server startet weiter unten in eigener Goroutine, ein Startfehler wird
dort über `log.Error` gemeldet und geht nicht in das Ergebnis von `Run` ein —
wie beim HTTP-Adapter." Beide Hälften stimmen: der Server startet in einer
eigenen Goroutine mit `defer grpcDone.Done()` und `log.Error` (`:708-714`), der
HTTP-Adapter ebenso (`:680-687`), und `Run` gibt ausschließlich
`mergeStreamAndWALFaultOutcome(streamErr, &walFault)` zurück (`:788`). Der Satz
beschreibt jetzt den Ist-Zustand; die alte Fehlaussage ist weg, nicht
umformuliert.

**F-3 — behoben.** Der NATS-Block steht zeichengleich direkt über seinem `if`
(`:553-562`/`:563`); die Broadcaster-Erklärung und der Aufbau von `captureOpts`
liegen jetzt davor. Kommentar und Gegenstand sind wieder zusammen.

**F-4 — behoben, und der Name trägt den Sachverhalt.**
`erwarteteErfolgreicheAufrufe` beschreibt genau die gezählte Größe
(`len(stream.published)` — Aufrufe, die ohne Fehler zurückkehren); im Fall
„nicht lesender Empfänger" ist der Wert `2` bei einer real verworfenen zweiten
Change weiterhin korrekt, weil der Doppel die Rückkehr, nicht die Zustellung
zählt. Mutation A färbt genau diesen Fall rot — der Wert ist lasttragend, nicht
Dekoration.

**F-6 — behoben, und die Aussage „Notify-Block bleibt unverändert" ist nicht
falsch geworden.** Die Begründung des verworfenen `tx.Changes()`-Fehlers steht
jetzt einmal im Helfer `changesOfCommittedTransaction`
(`service.go:157-164`); `Capture()` (`:149`) und `distinctTables` (`:175`)
nutzen ihn. Der Helfer tut zeichengleich das, was vorher an beiden Stellen
stand (`changes, _ := tx.Changes(); return changes`) — dieselbe Kopie
(`transaction.go:102`), dieselbe Reihenfolge, derselbe verworfene Fehler. Der
**Notify-Block** in `Capture()` (`:129-135`) ist byte-identisch unverändert; die
*eine* berührte Zeile liegt im Helfer der Notify-Deduplizierung und ist
verhaltensneutral — die Notify-Tests (`TestCaptureNotifiesAfterAckOnSuccess`,
`…NotifiesOnceForSameTableMultipleChanges`, `…NotifiesDistinctlyForTwoTables`)
sind unverändert und grün. Die Extraktion ist im selben §3, einen Bullet
unterhalb der Abgrenzungs-Aussage, ausdrücklich benannt. Kein Befund — siehe
aber F-9 zum Wortlaut der Abgrenzung.

**Zur Vorgeschichte F-1 (entschieden durch `ADR-0067`).** Der Architect-Zug
bestätigt den von F-1 belegten Widerspruch und wählt genau das, was der Report
als Verdikt 1 skizziert hat: die Fitness-Function-Zeile ist der Defekt, die
Entscheidung gilt. Alle drei im Report skizzierten Verdikte sind in §Verglichene
Alternativen abgebildet (dort A = Verdikt 3, B = Verdikt 1, C = Verdikt 2; dazu
zwei weitere: D „Zeile streichen", E „Zeile in Prosa"), und die Wahl ist mit
demselben Grund begründet, den der Report trägt (die Entscheidung ist am Port
strukturell gedeckt, eine caller-seitige Maßnahme wäre eine zweite Stelle).

**Trägt `ADR-0067` den Befund vollständig?** Ja — beide Hälften sind
geschlossen:

1. *Der Selbstwiderspruch:* Die dritte Zeile ist ersetzt
   (`0067-…md:158-172`), der Ersatztext grenzt die unerreichbare Behauptung
   ausdrücklich aus („**Nicht** Gegenstand dieser Zeile: ein
   `ChangeStreamPort`, dessen `Publish` nicht zurückkehrt"), und §Entscheidung
   sowie Option B der `ADR-0066` sind als bestätigt erklärt (`:3-12`, `:122-134`).
   Die neue Fitness Function (`:252-257`) nennt genau die zwei real gebauten
   Tests als Regel — beide existieren, beide sind grün, und Mutation A belegt,
   dass der zweite nicht blind ist.
2. *Die unerreichbare Klammer-Hälfte (F-7):* `ADR-0067` übernimmt den
   `app-impurity`-Beleg (`:83-92`) und macht die Schichtgrenze zum explizit
   benannten Bestandteil der Ersatz-Zeile („gegen einen Port-Doppel, weil diese
   Schicht keinen Adapter importieren darf, `.a-check.yml`", `:160-168`;
   Konsequenz „Negativ", `:234-237`).

Keine offene Hälfte gefunden. Zusätzlich geprüft: `SPEC-020` ist zu Recht
unberührt (`spec/pflichtenheft.md:356-358`: die Zeile *Erzeuger-Blockade* ist
eine **Port**-Aussage — „`Publish` blockiert nie auf einen Abonnenten" — und
verlangt keine Aufrufstellen-Zeit-Isolation; der Widerspruch lag allein in der
Fitness-Function-Zeile); `ADR-0060`s Fitness-Function-Zeilen sind unberührt;
`ADR-0066`s §Kontext („isoliert nur den Fehler, nicht die Zeit") und
§Konsequenzen („best-effort in Zeit und Fehler") bleiben unter der bestätigten
Entscheidung kohärent; der ADR-Index trägt die Zeile und den
`→ ADR-0067`-Vermerk (`README.md:79-80`, `AGENTS.md` §5).

### Neues Finding dieses Laufs

### F-9 — Der Wortlaut der §1/§3-Abgrenzung deckt die verhaltensneutrale Helfer-Extraktion im Notify-Pfad nicht mehr

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.7 — ein Kommentar/Text beschreibt,
  was da ist)
- `pfad`: `docs/plan/planning/in-progress/slice-070-grpc-capture-integration.md:76-77`
  (§1) und `:170-172` (§3), gegen
  `internal/application/usecase/capture/service.go:175`
- `befund`: §1 sagt „dieser Slice fasst keine Zeile des Kernpfads vor dem
  Stream-Publish-Aufruf an", §3 wiederholt das als „Keine Zeile des Kernpfads
  vor dem Stream-Publish-Aufruf angefasst"; die Fixrunde hat mit
  `distinctTables`s erster Zeile (Helfer der Notify-Deduplizierung) genau so
  eine Zeile berührt — verhaltensneutral, und im selben §3 einen Bullet tiefer
  ausdrücklich benannt. Die semantische Aussage hält (der Notify-Block
  `service.go:129-135` ist byte-identisch, der Helfer tut dasselbe wie zuvor);
  die *wörtliche* „keine Zeile"-Form ist seit der Fixrunde enger als der
  Vorgang.
- `verifizierbar`: ja — `git show 37717c6 -- internal/application/usecase/capture/service.go`
  gegen `git diff 6f9e9d8..37717c6 -- .../service.go`
- `klasse`: „Abgrenzungs-Wortlaut und getroffene Zeile fallen auseinander"

## Summary (Fixrunde)

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Abgrenzungs-Wortlaut und getroffene Zeile
fallen auseinander

## Verdikt (Fixrunde)

**Merge-blockierend:** nein — **keine Fixrunde mehr nötig.** F-2, F-3, F-4 und
F-6 sind am Code behoben und einzeln nachgeprüft (nicht an der Commit-Message),
F-1 ist durch [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
entschieden und trägt den Befund vollständig, F-5/F-7/F-8 sind INFO ohne
erwartete Aktion, F-9 ist ein Wortlaut-Hinweis ohne Fix-Zwang. Die Kette läuft
damit Reviewer → Architect (Folge-ADR `ADR-0067`) → Implementer (Fixrunde) →
Reviewer und schließt hier; der nächste Rollenwechsel ist Reviewer → Verifier
(DoD-/Spec-Konformität, Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-11-*.md`).

**DoD-Nachzug:** Mit diesem Verdikt ist die Rückkante geschlossen — die Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des
Slice-Plans ist in demselben Commit wie dieser Vermerk auf `[x]` gezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Nur diese
eine Zeile; Closure-Notiz, Beobachtungs-Register, die Risiko-Ausgänge und die
drei Paarungen bleiben offen (Planner-Arbeit).

**Übergabe:** F-9 geht als Wortlaut-Hinweis an die Planner-Seite (kein
Implementer-Rückweg); die **Finding-Klassen beider Läufe** gehen in die
Slice-Closure §7 und von dort in den Zähler.
