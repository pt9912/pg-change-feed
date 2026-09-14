# Verifikationsbericht: slice-070 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-070` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5
Closure-Trigger, §6 Risiken, §8 Sub-Area) und die bindenden Entscheidungen
([`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
— `Accepted`, superseded genau die **dritte Fitness-Function-Zeile** von
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
Prüfgrundlage dieses Laufs; [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
Teilfrage 2; [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
§Entscheidung/Festlegungen 1–4). Nicht gegen den Diff als solchen
(Reviewer-Aufgabe; [`review-slice-070`](review-slice-070.md) in beiden Läufen
vollständig gelesen, aber nur als Kontext) und nicht gegen realen Bedarf
(Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8,
Stand `HEAD = aacc6b2`), die vollständige
[`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md),
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) Teilfrage 2,
`SPEC-020` (`spec/pflichtenheft.md:339-361`), den vollständigen Review-Report
(beide Läufe), `service.go`, `service_test.go`, `wiring.go`,
`internal/domain/model/transaction.go`, `broadcaster_test.go`, `welle-19.md`
und den ADR-Index. Die Sensoren wurden in dieser Sitzung **eigenständig real
ausgeführt** (Exit-Code je in einem eigenen, ungepipten, mechanisch
konditionierten Schritt, `AGENTS.md` §3.9) — kein Implementer- oder
Reviewer-Beleg ungeprüft übernommen. Beide Mutationsstichproben sind
**eigene**, selbst gesetzt.

**Gegenstand:** `docs/plan/planning/in-progress/slice-070-grpc-capture-integration.md`
zum Stand `HEAD = aacc6b2`. Der Elter-`6f9e9d8` ist der reine
`next → in-progress`-Move (`git show 6f9e9d8 --stat`: `{next => in-progress}`,
0 Zeilen Änderung). Der Slice-eigene Code-Umfang liegt in `78d697f`
(Publish-Kette `+` Tests), `efd8d87`/`f5e00d0` (Digest-Stempel), `48733d7`
(Kommentar-Fix), `37717c6` (Review-Fixrunde) — genau
`service.go`, `service_test.go`, `wiring.go`, der Slice-Plan und
`harness/image-hash.txt`. Die Commits dazwischen (`d0790b3` Report,
`8423b28` Beobachtung, `20d643a` [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md),
`f714a88` Plan-Nachzug, `68d2ebd` Handbuch/Skill, `aacc6b2` Fixrunden-Verdikt)
sind Planner-/Architect-/Reviewer-Arbeit und **nicht** slice-070-Code-Umfang.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — die implementierungs-/
reviewbezogenen Zeilen sind Prüfgegenstand, die Closure-Zeilen **müssen** offen
bleiben, bis der Planner sie schließt.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | [`LH-FA-SST-008`](../../spec/lastenheft.md) Happy Path: `Publish` nach `ACK Source` und nach Notify, genau einmal je `Change` der committed Transaktion | **erfüllt, selbst reproduziert** | `service.go:118-120` (`Acknowledge`), `:129-135` (Notify), `:146-153` (Publish-Block); `service_test.go:504` pinnt die Ordnung als Ereignisliste `persist → ack → notify → stream:t-1-1 → stream:t-1-2`, `:508` die Change-Reihenfolge. Eigener `make test`-Lauf grün (§2). |
| 2 | [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) Teilfrage 2 vollständig: Fehlerisolations-Regression **und** Fire-and-Forget bei getrenntem **und** registriertem, nicht lesendem Empfänger | **erfüllt, selbst reproduziert** | `TestCaptureSucceedsDespiteFailingStream` (`service_test.go:549-575`) — selbst **rot gesehen** (§3, Mutation 1); `TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck` (`:585-646`) — selbst **rot gesehen** (§3, Mutation 2); `TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn` (`:656-689`) parkt `Capture()` in `Publish` und liest Store-/ACK-Stand in der Parkposition. Die [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)-Fußnote (§2) trägt die **zutreffende** Zeile, keine Reduktion (§4). |
| 3 | `make gates` grün | **erfüllt, selbst ausgeführt** | Eigener Lauf Exit **0** (§2). |
| 4 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | [`review-slice-070`](review-slice-070.md) vorhanden (2 Läufe, Fixrunden-Vermerk); Verdikt 0 HIGH / 0 MEDIUM im zweiten Lauf, deckungsgleich mit §2 der DoD-Zeile. |
| 5 | Doku-Update `wiring.go`-Kommentar, falls Aktivierungspfad sich ändert | **erfüllt** | `wiring.go:547-550` trägt den Ist-Zustand des Fehlerpfads (§5); die beiden `CDC_GRPC_ADDR`-Blöcke sind nachgezogen. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<…>`-Platzhalter (Volltext gelesen). Planner-Arbeit nach diesem Bericht. |
| 7 | Reconciliation-Register — entfällt (GF-Repo) | **korrekt offen, Entfall-Vermerk trägt** | `docs/plan/planning/reconciliation.md` existiert real nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `grep -rl "slice-070" docs/plan/planning/observations/` → **kein** Treffer. Die §8-Sichtung nennt die Treffer; die Eintragung gehört in die Closure. |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Einträge (`:232-236`, `:237-240`) tragen wörtlich „Ausgang: wird bei Closure zugewiesen"; keiner ist vorzeitig geschlossen. |
| 10 | Die drei Paarungen | **korrekt offen** | Slice liegt real in `in-progress/`; `Welle: welle-19`-Feld vorhanden; die DoD-Zeile verweist korrekt auf die `welle-19`-Closure. |

**Ergebnis §1:** Alle fünf implementierungs-/reviewbezogenen DoD-Punkte (1–5)
sind real erfüllt; 1 und 2 wurden **selbst reproduziert** — die Kernzusage
zusätzlich über zwei eigene Mutationen als tragend belegt. Die fünf
Closure-Punkte (6–10) sind korrekt offen und nicht vorweggenommen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | `baseline-verify v6.5.0` 54 Dateien OK; d-check 555 Dateien/0 Befunde (links/anchors/ids/matrix/versions/structure), `commits`-Modul `HEAD~5..HEAD` 0 Befunde; `commit-traceability: OK — 5 Commit(s)`, Betreffe ohne Struktur-ID; a-check gesamt 0 Befunde (2 unverschichtete Werkzeug-Dateien, bekannte Hinweise); `coverage-gate: OK — Coverage 48.40% erfüllt Schwelle 35%` (bootstrap-aware Einstiegsstufe [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)) |
| `make test` | **0** | 28 Pakete `ok`, **kein `FAIL`**; `internal/application/usecase/capture` und `internal/adapters/driven/grpcstream` grün |
| `make doc-commits RANGE=6f9e9d8..HEAD` | **0** | 555 Dateien, 0 Befunde — jeder Commit des Slice-Fensters ist kennungstragend |
| `make doc-immutable RANGE=6f9e9d8..HEAD` | **0** | 0 Befunde — keine `Accepted`-ADR überschrieben; `ADR-0066`s Textkörper ist im Fenster **unberührt** (§5) |

`git status --porcelain` nach allen Sensor- und Mutationsläufen: **leer** (§3).

## 3. Mutations-Stichprobe — Kernzusage selbst rot/grün gesehen (Verifier-only)

Der Implementer/Reviewer nennt Mutationen mit Exit-Codes. Ich habe die zwei
tragenden **selbst gestellt**, jede in einem eigenen `make test`-Aufruf, danach
per `git checkout --` zurückgenommen.

| # | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| 1 | `s.log.Warn(…)` im Publish-Block durch `return CaptureResult{}, err` ersetzt (`service.go:149-151`) | rot genau an der Fehlerisolation | **rot** (Exit 2; genau `TestCaptureSucceedsDespiteFailingStream`, `service_test.go:560`: „Capture: Stream nicht verfügbar, wollen keinen Fehler trotz fehlschlagendem Stream-Publish") |
| 2 | `erwarteteErfolgreicheAufrufe: 2 → 1` im Fall „nicht lesender Empfänger" (`service_test.go:592`) | rot genau an diesem Unterfall | **rot** (Exit 2; genau `TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck/nicht_lesender_Empfänger`, `:626`: „Stream trägt 2 erfolgreiche Aufrufe, wollen 1 (der Port kehrt je Change zurück)") |

**Rücknahme und Blatt-Identität:** nach `git checkout --` stimmt
`git hash-object` beider Dateien mit `git rev-parse HEAD:<pfad>` überein
(`service.go` `b6a912c0…`, `service_test.go` `e1d86977…`),
`git status --porcelain` ist leer.

**Ergebnis §3:** Die Fehlerisolation (Mutation 1) und der Fire-and-Forget-Beleg
(Mutation 2) sind **einzeln tragend** und **rot gesehen** — genau an der
zugesagten Stelle. Beide Erwartungen sind lasttragend, keine Dekoration; keine
greift ins Leere. Die Ordnung der Kette trägt zusätzlich
`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`, der bei verletzter
Vor-`Publish`-Ordnung rot würde (Review-Mutation (ii), hier nicht wiederholt;
der Test liest den Store-/ACK-Stand in der Parkposition, §4).

## 4. `ADR-0067` gegen den ausgelieferten Code

Die neue Fitness-Function-Zeile der
[`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
(`:158-172`, `:252-257`) trägt zwei Hälften, und beide sind an genau je einem
Test verankert. Ich habe geprüft, dass die Tests **genau das** ausüben — und
**nicht** heimlich die von der ADR ausgeschlossene Behauptung.

- **Hälfte (a) — kritische Kette abgeschlossen, bevor `Publish` betreten wird.**
  `TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`
  (`service_test.go:656-689`) startet `Capture()` in einer Goroutine, wartet auf
  `stream.eingetreten` (der Doppel sendet in `Publish`, `:475-482`) und liest
  **in dieser Parkposition** `store.persisted == 1`, `ack.acked == [100]` und
  `events == [persist:t-1, ack:100]` ab (`:675-683`) — erst danach `close(freigabe)`.
  Der Test zeigt damit real: `Persist`/`ACK` sind geschrieben, **bevor** der
  erste Change das Haus verlässt. **Deckungsgleich mit (a).**
- **Hälfte (b) — `Capture()` kehrt gegen einen vertragstreuen Port ohne eigene
  Zeit-Isolation zurück.** `TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck`
  (`:585-646`) nutzt `fakeStream` (`:444-464`) — einen **vertragstreuen** Doppel
  (nicht-blockierender Send mit Drop-Newest, kehrt immer zurück) — und prüft
  beide Unterfälle: getrennter Client (`empfaenger == nil`) und registrierter,
  nicht lesender Empfänger (`make(chan *model.Change, 1)`). Beide kehren
  innerhalb der 2-Sekunden-Frist zurück (`:613-623`), mit 2 erfolgreichen
  Aufrufen. **Deckungsgleich mit (b)** — inklusive der von der ADR ausdrücklich
  genannten Klammer-Hälfte („registrierter, nicht lesender Empfänger“) gegen
  einen Doppel.

**Nichts ausgeschlossen wird heimlich getestet.** Die von der ADR
ausgeschlossene Behauptung („ein `ChangeStreamPort`, dessen `Publish` seinen
Vertrag verletzt und nicht zurückkehrt, hält `Capture()` nicht an“) wird von
**keinem** der beiden Tests geprüft: `PublishOhneRueckkehr…` *verlangt* im
Gegenteil, dass `Capture()` in `Publish` parkt, und gibt erst dann frei;
`Kehrt…` verwendet ausschließlich den vertragstreuen Doppel. Kein Test greift
zur verworfenen Option B (caller-seitige Goroutine oder Deadline-`ctx`).
Bestätigend: `internal/application/usecase/capture` importiert `grpcstream`
nicht (`grep` leer), kein Adapter-Paket im Test importiert.

**Port-Hälfte (Zeile 1)** bleibt in Kraft und real belegt:
`TestNichtLesenderAbonnentHaeltPublishNichtAn` existiert in
`internal/adapters/driven/grpcstream/broadcaster_test.go:117` und ist im eigenen
`make test` grün. Die zweite FF-Zeile der
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(Fire-and-Forget ohne Subscriber) trägt derselbe Lauf.

**Zum Supersede-Umfang.** [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
erklärt in Status und §Entscheidung ausdrücklich, dass **nur** die dritte Zeile
ersetzt und §Entscheidung (synchron, keine caller-seitige Goroutine) sowie die
beiden ersten Zeilen **bestätigt** sind. Der ADR-Index trägt die Zeile samt
`supersedes`-Vermerk und den Rückzeiger bei
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(`docs/plan/adr/README.md:79-80`). `make doc-immutable` 0 Befunde bestätigt, dass
`ADR-0066`s Datei nicht in-place überschrieben wurde (`AGENTS.md` §3.5).

**Ergebnis §4:** Die neue Fitness-Function-Zeile ist am Code **vollständig
getragen** und grenzt die unerreichbare Behauptung korrekt aus. Die von der ADR
geforderte Doppel-Form ist die Schichtgrenze (`.a-check.yml`), kein Mangel —
und der Beleg für die reale Nicht-Blockade liegt am Port.

## 5. Publish-Kette am Code, Helfer-Neutralität, `wiring.go`-Kommentar

- **Kette und Ordnung.** `Capture()` ruft `store.PersistTransaction` (`:114`),
  `ack.Acknowledge` (`:118`), den Notify-Block (`:129-135`) und **danach** den
  Stream-Publish-Block (`:146-153`). Der Publish liegt real **nach** `ACK Source`
  und **nach** dem Notify-Schritt ([`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  Teilfrage 2). **Deckungsgleich.**
- **Ein `Publish` je `Change`, in Reihenfolge, ohne Dedup.** `changes := changesOfCommittedTransaction(tx)`,
  dann `for i := range changes { s.stream.Publish(ctx, &changes[i]) }` (`:147-152`);
  `tx.Changes()` liefert die Changes in Anhang-Reihenfolge als Kopie
  (`transaction.go:98-103`). Keine Tabellen-Deduplizierung (anders als das
  Notify). `TestCapturePublishesEachChangeOnceAfterAckAndNotify` belegt zwei
  Aufrufe für zwei Changes derselben Tabelle (`:498`, `:508`). **Deckungsgleich.**
- **`store`/`ack` durch den neuen Schritt unberührt** ([`ADR-0011`](../plan/adr/0011-persist-before-ack.md)/[`ADR-0027`](../plan/adr/0027-capture-application-service.md)).
  Der neue Block ruft ausschließlich `s.stream.Publish` und im Fehlerfall
  `s.log.Warn`; kein `store`-/`ack`-Aufruf darin. **Bestätigt.**
- **Helfer `changesOfCommittedTransaction` (F-6) — verhaltensneutral.**
  `git show 37717c6 -- service.go` zeigt die Extraktion: an beiden Stellen
  (`Capture` und `distinctTables`) stand zuvor zeichengleich
  `changes, _ := tx.Changes()`; jetzt ruft jede den Helfer
  (`:157-164`, aufgerufen `:147` und `:178`), der genau diesen Rumpf trägt.
  Dieselbe Kopie, dieselbe Reihenfolge, derselbe verworfene Fehler (der
  `CommitPosition`-Wächter `:110-113` liegt davor, `Changes()` scheitert nur an
  einer offenen Transaktion). **Deckungsgleich — verhaltensneutral.** Die drei
  Notify-Tests (`TestCaptureNotifiesAfterAckOnSuccess`,
  `…OnceForSameTableMultipleChanges`, `…DistinctlyForTwoTables`) sind grün.
- **`wiring.go`-Kommentar (F-2) — beschreibt, was da ist.** `wiring.go:547-550`
  sagt zu: „Die Bindung ist keine Start-Vorbedingung: der gRPC-Server startet
  weiter unten in eigener Goroutine, ein Startfehler wird dort über `log.Error`
  gemeldet und geht nicht in das Ergebnis von `Run` ein — wie beim HTTP-Adapter."
  Alle drei Hälften stimmen: der gRPC-Server startet in eigener Goroutine mit
  `defer grpcDone.Done()` und `log.Error` (`:708-714`), der HTTP-Adapter ebenso
  (`:680-687`), und `Run` gibt ausschließlich
  `mergeStreamAndWALFaultOutcome(streamErr, &walFault)` zurück (`:788`). Kein
  Fehlerkanal, kein `errgroup`. **Deckungsgleich.**
- **Unverändert-Zusagen.** `grep "grpcstream"` über
  `internal/application/usecase/capture/` → leer (kein Adapter-Import,
  `.a-check.yml` unverändert); kein `nolint`/`noqa`/`SuppressMessage` in den
  geänderten `.go`-Dateien (`AGENTS.md` §3.2); kein `slice-`/`welle-`-Treffer in
  den Kommentaren der drei geänderten `.go`-Dateien
  (`BEO-PGC/slice-chronik-in-code-kommentar`).

## 6. DoD-Häkchen-Trennung

Der DoD-Block ist in zwei Zustände getrennt: Zeilen 1–5 (Happy Path, Teilfrage 2,
`make gates`, Review, Doku) `[x]`; Zeilen 6–10 (Closure-Notiz, Reconciliation,
Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) `[ ]`. Der Reviewer hat
in der Fixrunde (Commit `aacc6b2`) **genau** die Review-Zeile nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde) — keine
Closure-Zeile vorweggenommen. **Korrekt.**

## 7. Urteil zu F-9 (Abgrenzungs-Wortlaut)

**F-9 — kein Blocker, aber ein realer Text-Defekt; Planner-Korrektur am Wortlaut
empfohlen.** §1 (`:76-77`) sagt „dieser Slice fasst keine Zeile des Kernpfads vor
dem Stream-Publish-Aufruf an", §3 (`:178-180`) wiederholt das als „Keine Zeile
des Kernpfads vor dem Stream-Publish-Aufruf angefasst". Die Fixrunde hat mit der
Helfer-Extraktion genau **eine** Zeile in `distinctTables` (`service.go:178`)
berührt — den Helfer der Notify-Deduplizierung, also im Kernpfad **vor** dem
Stream-Publish. Die Aussage ist damit seit der Fixrunde **wörtlich falsch**, aber
**semantisch nicht gebrochen**: der Notify-Block (`:129-135`) ist byte-identisch,
der Helfer tut zeichengleich dasselbe wie zuvor (§5), und die getragene Invariante
(`Persist → ACK` unantastbar) hält. Die Extraktion ist im selben §3 einen Bullet
tiefer (F-6-Bullet, `:187-189`) ausdrücklich benannt.

Kein DoD-Bezug: die „keine Zeile"-Formulierung steht **nicht** in §2 (DoD),
sondern in der Out-of-Scope-Deklaration §1/§3. Sie ist damit kein Gegenstand der
DoD-Konformität — die DoD ist erfüllt (§1). Nach `AGENTS.md` §3.7 beschreibt der
Satz aber nicht mehr, was da ist. Der minimale korrekte Träger ist eine
**Planner-Wortlaut-Korrektur** an §1/§3 (wie die [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)-Folgepflicht
sie für §2/§3 bereits verlangt hat) — z. B. den Satz auf die **bestehende**
Kette (`Persist`/`ACK`/Notify-Block) einschränken. Ein Implementer-Rückweg ist
nicht nötig; wenn der Planner die Drift als immateriell beurteilt, genügt eine
benannte Closure-Notiz. **Mein Verdikt: Planner-Wortlaut-Korrektur, kein
Verifier-Blocker.**

## 8. `welle-19`-Stand — was nach dieser Closure noch fehlt

Der Closure-Trigger (`welle-19.md` §3) hat fünf Bedingungen; ich prüfe jede
gegen den Stand:

| Bedingung | Stand meiner Prüfung | Erfüllt? |
|---|---|---|
| Alle vier Slices (`slice-069`…`slice-072`) in `done/` | `slice-069` in `done/`; `slice-070` in `in-progress/`; `slice-071`/`slice-072` in `open/` | **nein** |
| `make gates` grün | eigener Lauf Exit **0** | **ja** |
| Realer gRPC-E2E-Rundlauf (`slice-071`, `make test-integration`) | `slice-071` in `open/`, nicht begonnen; kein Beleg | **nein** |
| Realer SSE-E2E-Rundlauf (`slice-072`, `make test-integration`) | `slice-072` in `open/`, nicht begonnen; kein Beleg | **nein** |
| Closure-Notiz in `welle-19-results.md` | Datei existiert real nicht; `welle-19.md` §7 noch Platzhalter | **nein** |

**Ergebnis §8:** Nach der `slice-070`-Closure fehlen zur Welle noch **zwei
Slices** — `slice-071` (gRPC-Beispiel-Client/E2E) und `slice-072` (HTTP/SSE) —,
die beide noch in `open/` liegen und erst `next → in-progress` und dann `done`
durchlaufen müssen, sowie die beiden realen E2E-Belege und die Wellen-Closure
(`welle-19-results.md`, `git mv` der Welle-Datei nach `done/`). `make gates` ist
bereits grün. `slice-072` kann laut `welle-19.md` §4 parallel zu `slice-071`
laufen, sobald `slice-070` in `done/` liegt.

## Negativbefunde

- **geprüft, ohne Befund: DoD-Punkt-Trennung und -Vollständigkeit.** Zehn
  Punkte, fünf abgehakt (Happy Path, Teilfrage 2, Gate, Review, Doku), fünf
  offen (Closure-Notiz, Reconciliation, Register, Risiko-Ausgänge, Paarungen) —
  keine vorweggenommene Closure-Zeile (§6).
- **geprüft, ohne Befund: die zwei FF-Tests testen nicht die ausgeschlossene
  Behauptung.** `PublishOhneRueckkehr…` parkt bewusst und gibt frei;
  `Kehrt…` nutzt den vertragstreuen Doppel. Kein Test verlangt die verworfene
  Option B (§4).
- **geprüft, ohne Befund: Scope-Fidelity der Slice-eigenen Commits.** `78d697f`,
  `48733d7`, `37717c6`, `efd8d87`, `f5e00d0` berühren ausschließlich
  `service.go`, `service_test.go`, `wiring.go`, den Slice-Plan und
  `harness/image-hash.txt` — kein Eingriff unter
  `internal/adapters/driven/grpcstream/**` (Port-Seite bleibt `slice-069`,
  **unverändert**), kein `internal/adapters/driving/**`, keine Änderung an
  `ChangeNotificationPort`/`natsnotify` oder `.a-check.yml`.
- **geprüft, ohne Befund: `ADR-0066`s Textkörper unverändert.**
  `make doc-immutable RANGE=6f9e9d8..HEAD` Exit 0; der `Supersedes`-Vermerk steht
  am ADR-Index, die Klausel-Supersession ist eng auf die dritte Zeile gefasst.
  `AGENTS.md` §3.5 eingehalten.
- **geprüft, ohne Befund: Helfer-Extraktion verhaltensneutral** (§5) — die
  Notify-Tests sind unverändert und grün.
- **geprüft, ohne Befund: Kommentar-Hygiene der geänderten `.go`-Dateien.**
  Kein `slice-`/`welle-`-Chronik-Treffer in `service.go`, `service_test.go`,
  `wiring.go`; kein Suppression-Marker (`AGENTS.md` §3.2).
- **geprüft, ohne Befund: Traceability.** `make doc-commits RANGE=6f9e9d8..HEAD`
  Exit 0; die Slice-Commits tragen `LH-FA-SST-008`/`ADR-0060`/`ADR-0067`, kein
  `SPEC-*`/`ARC-*` im Betreff.
- **geprüft, ohne Befund: `internal/application/usecase/capture` importiert kein
  Adapter-Paket.** `grep "grpcstream"` leer; die Schichtregel (`.a-check.yml`)
  ist nicht verletzt, `make a-check` 0 Befunde.

## Summary

| Prüfung | Ergebnis |
|---|---|
| DoD-Punkte 1–5 (implementierungs-/reviewbezogen) | **erfüllt**, 1/2 selbst reproduziert |
| DoD-Punkte 6–10 (Closure) | korrekt offen, keine vorweggenommen |
| Vier Sensoren (`gates`, `test`, `doc-commits`, `doc-immutable`) | **je Exit 0** |
| Mutation 1 (`log.Warn` → `return CaptureResult{}, err`) | **rot** an `TestCaptureSucceedsDespiteFailingStream`; Rücknahme `+` Blatt-Identität verifiziert |
| Mutation 2 (`erwarteteErfolgreicheAufrufe: 2 → 1`) | **rot** am Unterfall „nicht lesender Empfänger"; Rücknahme `+` Blatt-Identität verifiziert |
| Publish-Kette (nach `ACK`/Notify, je `Change`, ohne Dedup, `store`/`ack` unberührt) | **deckungsgleich** |
| [`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md) neue FF-Zeile gegen den Code | **vollständig getragen**; kein Test heimlich gegen die Ausschluss-Klausel |
| Helfer `changesOfCommittedTransaction` (F-6) | **verhaltensneutral**, Notify-Tests grün |
| `wiring.go`-Kommentar (F-2) | **deckungsgleich** mit dem realen Fehlerpfad |
| DoD-Häkchen-Trennung | korrekt |
| F-9 | **Planner-Wortlaut-Korrektur, kein Blocker** |
| `welle-19`-Closure-Reife | **noch nicht schließbar** — zwei Slices offen, zwei E2E-Belege fehlen |

**Verdikt:** `slice-070` erfüllt seine DoD. Die Kernzusage
([`LH-FA-SST-008`](../../spec/lastenheft.md) Happy Path: ein `Publish` je `Change`
nach `ACK Source` und Notify, fehlerisoliert, nicht-blockierend) ist am Code
deckungsgleich umgesetzt und durch **zwei eigene** Mutationen als **rot** belegt
— genau an den zugesagten Stellen, mit Rücknahme und verifizierter
Blatt-Identität. Die über
[`ADR-0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
korrigierte dritte Fitness-Function-Zeile der
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md) wird
von genau den zwei genannten Tests getragen, die die kritische Kette vor
`Publish` beziehungsweise die Rückkehr gegen einen vertragstreuen Doppel
ausüben — und **nicht** die von der ADR ausgeschlossene Behauptung. Die
Helfer-Extraktion (F-6) ist verhaltensneutral, der `wiring.go`-Kommentar (F-2)
beschreibt den Ist-Zustand, die DoD-Häkchen sind korrekt getrennt. Alle vier
selbst ausgeführten Sensoren laufen Exit 0.

**Einzige Rest-Beobachtung ohne Zwang:** der §1/§3-Wortlaut „keine Zeile des
Kernpfads … angefasst" ist durch die verhaltensneutrale Helfer-Extraktion der
Fixrunde wörtlich überholt (F-9); die semantische Invariante hält. Empfehlung:
eine Planner-Wortlaut-Korrektur; kein Implementer-Rückweg.

**Übergabe an den Planner (Verifier → Planner):** Der Slice ist DoD-konform;
über diese Verifikation hinaus ist **kein** Implementer-Rückweg nötig. Für den
`slice-070`-Closure-Schritt bleiben: Closure-Notiz mit Steering-Loop-Eintrag
(inkl. der beiden Finding-Klassen aus dem Review und der F-9-Klasse), die zwei
§6-Risiko-Ausgänge (heute „wird bei Closure zugewiesen"), die Beobachtungs-
Register-Entscheidung, die F-9-Wortlaut-Korrektur und der `git mv` nach `done/`.
Für die danach folgende `welle-19`-Closure: mein Lauf dieses Berichts ist der
repo-weite Verifikations-Beleg aus Closure-Schritt 1 (Modul 6/8); die Welle ist
erst schließbar, wenn `slice-071` und `slice-072` in `done/` liegen und beide
realen E2E-Belege (gRPC über `slice-071`, SSE über `slice-072`) grün vorliegen.
