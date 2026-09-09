# Verifier-Report: slice-003 — 2026-09-09

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
11 Items), §3 (Plan-vs-Code, Range `a227033..9bfcd7a` inkl. Plan-Nachzug
`4a41522`), §6 (Risiko-Ausgänge) und ADR-Konformität
([`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Supersedes
[`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)) ·
[`ADR-0011`](../plan/adr/0011-persist-before-ack.md) ·
[`ADR-0012`](../plan/adr/0012-at-least-once.md) ·
[`ADR-0029`](../plan/adr/0029-domain-invarianten.md) Regel 1 ·
[`ADR-0027`](../plan/adr/0027-capture-application-service.md) ·
[`ADR-0028`](../plan/adr/0028-inbound-use-cases.md) ·
[`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) ·
[`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md)). Nicht
geprüft: Diff gegen Plan/Hard Rules (Reviewer, `review-slice-003.md`,
Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-003-capture-persist.md` · Range
`a227033..9bfcd7a` (4 Commits: 58f0192, cd1a6de, 1983e84, 9bfcd7a) ·
Fortsetzung nach Review-Schluss: `4a41522` (Review-Report + Plan-Nachzug +
Register-Korrektur) · `55ce39a` ([`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)), Architect-Verdikt der
Konflikt-Sequenz) · `c263797` (Klasse-B-Verkörperung: Commit-Message-Regel
im Implementer-Briefing) · Arbeitsbaum-HEAD `c263797` · **Fix-Commit
`9bfcd7a` trägt F-4/F-5/F-6 aus dem Review.**

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über
stdout, der Arbeitsbaum wurde read-only gehalten (go-Belege über einen
`git worktree` unter `/tmp`, nach Abschluss entfernt; der Hauptbaum steht
clean auf dem Handoff-Stand, Gate-Stempel deckungsgleich — unten).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Slice-Plan §1–§8 · `harness/README.md` (Sensors-Tabelle) · `Makefile`,
  `.a-check.yml`, `.d-check.yml`
- [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Accepted,
  Supersedes [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)) ·
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)
  (Superseded, Struktur-Regeln als Rest fort) ·
  [`ADR-0011`](../plan/adr/0011-persist-before-ack.md) ·
  [`ADR-0012`](../plan/adr/0012-at-least-once.md) ·
  [`ADR-0029`](../plan/adr/0029-domain-invarianten.md) ·
  [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) ·
  [`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md)
- `docs/reviews/review-slice-003.md` (F-1…F-9, Übergabe-Artefakt der
  Vorgängerrolle) · Fix-Commit `9bfcd7a` im Volltext · Plan-Nachzug
  `4a41522` · `55ce39a` · `c263797` im Volltext
- `spec/lastenheft.md` ([`LH-QA-REL-001`](../../spec/lastenheft.md),
  [`LH-QA-REL-002`](../../spec/lastenheft.md)) · `spec/pflichtenheft.md`
  ([`LH-QA-REL-001.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), `LH-FA-CAP-005`,
  [`SPEC-001`](../../spec/pflichtenheft.md), `SPEC-008`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `c263797`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 82 Datei(en) geprüft, 0 Befund(e)` · `a-check … gesamt: 0 Befund(e)` — a-check im Bündel; mit `usecase/**` matcht jetzt auch der `app`-Layer-Glob echten Content | 0 |
| `go test -count=1 ./internal/...` im gepinnten Toolchain-Container (`golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125`, go1.27.1, `/src:ro`, `--network none`, GOCACHE im Container) | `ok … internal/application/usecase/capture` · `ok … internal/application/port/outbound` · `ok … internal/domain/model` (3 Pakete mit Tests) | 0 |
| `go vet ./...` (derselbe Container, dieselbe Referenz) | `VET_CLEAN` — keine Ausgabe | 0 |
| `gofmt -l .` (derselbe Container-Lauf) | `GOFMT_CLEAN` | 0 |
| Gate-Stempel | `.harness/state/gates-passed.diffsha` = `b44b3ed3cb00620726fdd730b7b0c56009061bf45525d378416e587476a76f8c` = `tools/harness/working-tree-hash.sh` am Prüfpunkt — deckungsgleich | — |
| Mutations-Probe A (leere committed Transaktion wird abgelehnt statt durchgelassen — Worktree, revertiert) | `--- FAIL: TestCapturePersistsEmptyCommittedTransaction` · `Capture: Transaktion ist nicht committed` | erwartet rot |
| Mutations-Probe B3 (committed-Flag-Guard außer Kraft: `if !committed && false`, Worktree, revertiert) | `--- FAIL: TestCaptureRejectsOpenTransaction` · `Capture-Fehler = <nil>, wollen Transaktion ist nicht committed` | erwartet rot |
| `make doc-commits RANGE=a227033..c263797` | `d-check: 82 Datei(en) geprüft, 0 Befund(e)` (Traceability je Commit, volle Range inkl. Fix, Plan-Nachzug, [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md), Verkörperung) | 0 |
| `make doc-immutable RANGE=a227033..c263797` | `d-check: 82 Datei(en) geprüft, 0 Befund(e)` (MR-/Doku-Immutabilität) | 0 |
| `git ls-remote origin main` | `929e449ef5d32c5162c3db080c5ccf13fdc8a145` — **nicht** `c263797` (siehe V-1) | — |

Hinweis zu den Mutations-Proben: Die drei vom Implementer/Review gemeldeten
Proben galten dem **Vor-Fix**-Stand (Probe 3 traf den `Changes()`-Guard, den
F-5 inzwischen ersetzt hat). Für den Endstand habe ich beide Wächter des
Fix-Commits selbst durch eigene Proben rot gesehen (A → F-4-Grenze, B3 →
F-5-Guard); die Worktrees sind entfernt, der Hauptbaum blieb read-only
(`git status` clean, `git worktree list` trägt nur den Vorbestand von
slice-002 — V-3).

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `CaptureInboundPort` + `CaptureService` tragen die Persist-before-ACK-Ordnung; Test referenziert zu [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) | **bestätigt** | `service.go:57-78` — die einzige Orchestrierungs-Stelle (persist → ack, kein zweiter ACK-Pfad, Guard vor Persist); Ordnung in drei konsistenten Kontrakt-Restates (Inbound-Port `capture.go:34-37`, Store-Port `changestore.go:15-19`, Service-Kommentar); [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) als Port-Kontrakt-Verankerung zitiert in Port, Service und Tests; `TestCapturePersistsBeforeAck` grün, Ereignisbeleg `persist:t-1 → ack:100` an einer geteilten Liste |
| 2 | Application-Tests mit Fakes: kein ACK bei Persistenzfehler; Crash zwischen Persistenz und ACK → höchstens erneute Verarbeitung | **bestätigt** | `TestCaptureDoesNotAckOnPersistenceError` (Persist-Fehler → kein ACK-Ereignis, kein `ack.acked`-Eintrag) · `TestCrashBetweenPersistAndAckLeadsToReprocessing` (Persistenz besteht, ACK fehlt, Wiederholung führt zu Ende; Doppel-Persist im Fake ehrlich dokumentiert mit [`SPEC-002`](../../spec/pflichtenheft.md)-Hinweis an den realen Adapter) — beide grün im Container-Lauf; beide Wächter mutations-geprüft (Probe B3 und Review-Proben) |
| 3 | Application-Testsuite grün (`go test ./internal/application/...`) | **bestätigt** | Volle `./internal/...` im gepinnten 1.27-Container grün (Tabelle oben, 3 Test-Pakete, `-count=1`) |
| 4 | `make gates` grün | **bestätigt** | Drei Gates grün am HEAD (Tabelle oben); Stempel deckungsgleich mit dem Arbeitsbaum-Hash |
| 5 | Review-Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/review-slice-003.md` committet (`4a41522`); Rollenwechsel nach Schritt 8 eingehalten — Implementer-Diff → Review (F-1…F-9) → Fix `9bfcd7a` → Sequenz-Artefakte → Verifier; kein Self-Review |
| 6 | Doku-Update bzw. begründete Aussage „kein öffentlicher Vertrag berührt" | **bestätigt** | Code-Range `a227033..9bfcd7a` berührt ausschließlich `internal/**` (5 neue Dateien, `git diff --stat`); Spec-Straten und `harness/README.md` unverändert — **kein öffentlicher Vertrag berührt, die Aussage trägt**. Die Folge-Commits sind Rollen-Artefakte, keine Vertragsberührung: `55ce39a` ([`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) + ADR-Index, Traceability-Regel „neue ADRs in den Index" erfüllt), `c263797` (Implementer-Briefing) |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure** | §7 trägt nur Platzhalter; Slice liegt korrekt noch in `in-progress/` — Planner in Arbeit, fällig vor dem `git mv` nach `done/` (inkl. Finding-Klassen des Reviews, V-1 und V-2 dieses Berichts) |
| 8 | Reconciliation-Register | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen |
| 9 | Beobachtungs-Register fortgeschrieben | **offen bis Closure — läuft** | `state.md` ist auf die Abgeleitet-Regel zurückgebaut (F-7-Korrektur in `4a41522`: geführte Zähler-Zeile entfernt, „der Zähler ist die Zahl der evidence-Dateien (2× nach slice-002) und wird nie gespeichert") — der gespeicherte Zähler ist weg, genau die geforderte Berichtigung. `evidence/` trägt `slice-001.md`, `slice-002.md` = **2×**; mit diesem Slice erreicht die Beobachtung **3×** (F-9 des Reviews: der `app`-Glob matcht seit diesem Slice echten Content — `usecase/capture` in der Range). Fällig vor dem `git mv`: `evidence/slice-003.md` und der Ausgang im Lese-Schritt (Regel-Material steht bereit: [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) trägt die Alias-Pflicht, a-check bleibt die Maschinenform) |
| 10 | Jedes §6-Risiko mit Ausgang | **teilweise erfüllt — Endwert fällig bei Closure** | Beide §6-Risiken stehen mit Voreinstellung „weiter offen": (a) Fake-Ordnung stärker als realer Treiber — unverändert tragend, Welle 2; (b) reale Idempotenz des Stores — durch F-6 **vorwärts getragen**: die Pflicht steht jetzt am Port-Kontrakt (`changestore.go:21-25`), ihr voller Beleg bleibt beim realen Adapter (Welle 2). Bewertung der Endausgänge bei Closure |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (flache `docs/plan/planning/welle-1.md` vorhanden): der DoD-Wortlaut weist die Prüfung der Welle-1-Closure zu; hier nicht fällig, notiert |

## Plan-vs-Code-Diff (Range `a227033..9bfcd7a`, Plan-Nachzug `4a41522`)

**Deckung §3:** Die Range enthält genau die fünf §3-Dateien —
`port/inbound/capture.go`, `port/outbound/changestore.go` (inkl. der
plan-vermerkten Grenz-Träger-Zeile zu [`ADR-0029`](../plan/adr/0029-domain-invarianten.md)
Regel 1),
`port/outbound/replicationack.go`, `usecase/capture/service.go`,
`usecase/capture/service_test.go` (429 Zeilen, alle neu). Keine
unbudgetierte Code-Datei.

**Plan-Nachzug (`4a41522`):** Die drei Nachzugs-Zeilen stehen in §3 und
tragen die Design-Entscheidungen auf dem Endstand:
Transport-Typen am Inbound-Port mit Aliasen (Zeile „Plan-Nachzug (Review
F-2, Konflikt-Sequenz Klasse A 3×)" — benennt die Schärfung als
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)), `context.Context`-Signaturen +
`ErrMissingTransaction` als Port-Fehler, und die F-4-Grenzentscheidung
(leere committed Transaktion: ablehnen oder Grenze am Port-Kontrakt + Test).
Plan und Diff sind vor der Closure deckungsgleich; die drei
Review-Blockierpunkte sind erfüllt (Blocker 1 → `4a41522` + `55ce39a`,
Blocker 2 → `c263797`, Blocker 3 → `9bfcd7a`).

**Nicht in §3, aber zulässig:** `docs/reviews/review-slice-003.md` und der
Plan-Nachzug selbst (`4a41522`) — Review-Report ist Lauf-Beleg pro Slice und
zählt nicht zum Umfang (Modul 5); `docs/plan/adr/0042-*.md`,
ADR-Index-Update und `.d-check.yml`-Exempt (`55ce39a`) — Übergabe-Artefakt
der Konflikt-Sequenz, nicht Slice-Umfang; `state.md`-Korrektur (`4a41522`) —
Register-Sache; `.claude/commands/implement-slice.md` (`c263797`) —
Verkörperung der F-3-Sequenz, trägt den Herkunfts-Anker `seit slice-003`.

## ADR-Konformität

- **[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)** (Transport-Typen am Port, Accepted — Supersedes [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)): **konform.** `CaptureCommand`/`CaptureResult` sind am Inbound-Port definiert (`capture.go:18-27`), der Use Case führt beide als Typ-Aliase (`service.go:21-24`) und implementiert den Port (`var _ inbound.CaptureInboundPort = …`); der Re-Evaluierungs-Trigger der ADR („ein Port trägt einen Transport-Typ, den sein Use Case nicht aliasiert") ist nicht erfüllt — kein Alias fehlt. Die Fitness-Zeile (a-check: `app → ports`-Kante, keine `adapters → app`-Kante) ist durch `make a-check` (0 Befunde) belegt. *Abweichung am Rand: V-2 — der Service-Kommentar zitiert noch die superseded [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md) an genau der geschärften Stelle.*
- **[`ADR-0011`](../plan/adr/0011-persist-before-ack.md)** (Persist-before-ACK): **konform.** Die Invariante `ACK(position) => durable(…)` ist an **einer** Stelle orchestriert (`service.Capture` — persist, dann ack, kein zweiter ACK-Pfad); die Grenz-Träger-Zeile steht am Store-Port-Kontrakt (`changestore.go:15-19`, plan-vermerkt aus dem slice-002-Lerneintrag); die Idempotenz-Pflicht der ADR („Pflicht, kein Optimismus") trägt nach F-6 die Kontrakt-Zeile `changestore.go:21-25` (deduplizierbar über die interne Transaktions-ID, [`SPEC-001`](../../spec/pflichtenheft.md)).
- **[`ADR-0012`](../plan/adr/0012-at-least-once.md)** (At-Least-Once): **konform.** `TestCrashBetweenPersistAndAckLeadsToReprocessing` belegt das Crash-Fenster (Persistenz bleibt, ACK fehlt, Wiederholung führt zu Ende); die F-4-Grenzentscheidung (leere Transaktion durchlassen) ist aus der ADR-Hälfte „Wiederholung wird gegenüber Datenverlust bevorzugt" sauber abgeleitet — das Ablehnen würde dieselbe leere Transaktion endlos neu liefern.
- **[`ADR-0029`](../plan/adr/0029-domain-invarianten.md) Regel 1** (Grenz-Träger): **konform — in der plan-getrennten Form.** Regel 1 hat im [`ARC-001`](../../spec/architecture.md)-Modell-Satz kein Subjekt (slice-002, V-1); der Träger ist die plan-vermerkte Port-Kontrakt-Zeile (`changestore.go:16-18` nennt Regel 1 ausdrücklich) plus die Orchestrierung im Service — drei konsistente Restates (Inbound-Port, Store-Port, Service), kein Drift zwischen ihnen. Die Spannung „erzwungen im Domain Core" gegen den Träger am Port bleibt als bekannte Grenze bestehen (kein neuer Befund dieses Laufs; die Schicht-Sicht trägt die Ordnung im Application Layer).
- **[`ADR-0027`](../plan/adr/0027-capture-application-service.md) / [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md) / [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) / [`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md)** (Port-/Service-Formen): **konform.** `CaptureService` orchestriert über zwei Outbound-Ports (`ChangeStorePort`, `ReplicationAckPort`) ohne Adapter-Logik; `CaptureInboundPort` trägt seinen Namen und seine Transport-Typen; der ACK ist Core-gesteuert (die Application ruft ihn erst nach dem Store-Commit); der Store ist die einzige Persistenz-Grenze (kein Repository je Entity).

## Fix-Commit `9bfcd7a` gegen die Review-Findings

| Finding | getragen im Fix? | Beleg |
|---|---|---|
| F-4 (Grenze „leere committed Transaktion" unbenannt, kein Test) | **ja** | Kontrakt-Zeile am [`CaptureInboundPort`](../../spec/pflichtenheft.md) (`capture.go:39-44`: Grenze benannt, Begründung aus [`ADR-0012`](../../spec/pflichtenheft.md) — Ablehnen würde endlose Wiederholung erzeugen) **und** `TestCapturePersistsEmptyCommittedTransaction` (leere committed Transaktion → persistiert + geackt, Ereignisbeleg `persist:t-1 → ack:100`). **Rote Probe gedeckt:** meine Mutation A (leere Transaktion abgelehnt) lässt genau diesen Test FAIL (`Capture: Transaktion ist nicht committed`) — der Test trägt die Grenze in beide Richtungen |
| F-5 (committed-Flag ignoriert, Kopplung nur per Kommentar) | **ja** | Der Guard prüft das zweite Ergebnis selbst: `position, committed := tx.CommitPosition(); if !committed { … ErrTransactionNotCommitted }` — das ignorierte Ergebnis und die Kommentar-Kopplung sind entfernt. **Rote Probe gedeckt:** Mutation B3 (Guard außer Kraft) lässt `TestCaptureRejectsOpenTransaction` FAIL (`Capture-Fehler = <nil>`) |
| F-6 ([`ADR-0011`](../plan/adr/0011-persist-before-ack.md)-Idempotenz ohne Kontrakt-Träger) | **ja** | Kontrakt-Zeile am `ChangeStorePort` (`changestore.go:21-25`): PersistTransaction ist deduplizierbar, Wiederholung nach Crash zwischen Persistenz und ACK zugelassen, Deduplizierungsbasis interne Transaktions-ID ([`SPEC-001`](../../spec/pflichtenheft.md), `cdc.transaction`) |
| F-1/F-2 (Typ-Platzierung + Plan-Nachzug — Blockierpunkte 1) | **ja, über die Sequenz** | Plan-Nachzug `4a41522` (drei §3-Zeilen auf Endstand) + Architect-Verdikt [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (`55ce39a`, Supersedes [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md), drei Alternativen, Fitness-Zeile, Trigger) |
| F-3 (Struktur-ID in Commit-Message — Blockierpunkt 2) | **ja, verkörpert** | `c263797`: Regel „Struktur-IDs (`SPEC-*`, `ARC-*`) gehören NICHT in die Commit-Message — nur `LH-*`/`ADR-*` · seit slice-003" im Implementer-Briefing (`.claude/commands/implement-slice.md`); Messagen der Fix- und Folge-Commits tragen keine `SPEC-*`/`ARC-*` (grep-Beleg über die Range; `SPEC-008` steckt nur in `1983e84` und bleibt in der Historie unveränderlich) |
| F-7 (stale Zähler-Zeile im Register) | **ja** | `4a41522`: `state.md` ohne geführten Zähler — „der Zähler ist die Zahl der `evidence/`-Dateien … und wird nie gespeichert"; Stand-Anker `evidence/slice-002.md` |

Alle IDs der Fix-Message existieren ([`LH-QA-REL-001`](../../spec/lastenheft.md)/002,
`LH-FA-CAP-005/006/007`, [`ADR-0011`](../plan/adr/0011-persist-before-ack.md),
[`ADR-0012`](../plan/adr/0012-at-least-once.md) — grep-Beleg je Treffer in Lastenheft/Pflichtenheft/ADR-Index).

## Befunde

### V-1 — Handoff-Behauptung „HEAD (c263797) gleich origin/main" trifft nicht zu: die Range ist ungepusht

- `kategorie`: MEDIUM
- `quelle`: Implementer-Handoff („HEAD ist auf `c263797` (origin/main
  gleich)") · Modul 5 §Lifecycle („Aussagen über die Verzeichnis-Position
  gelten für den gemergten Stand") — Behauptung ohne Beleg
- `pfad`: `git ls-remote origin main` → `929e449ef5d32c5162c3db080c5ccf13fdc8a145`
  (= slice-002-Closure, **vor** der Range); der lokale `main` steht auf
  `c263797`, das Remote-Tracking ist aktuell (beide Quellen stimmen überein)
- `befund`: Zehn Commits liegen nur lokal (`7f5af2c … c263797` — inkl. der
  gesamten Code-Range, des Review-Reports, [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
  und der Verkörperung). Kein DoD-Punkt verlangt den Push, und der Slice
  liegt korrekt in `in-progress/` — aber die Handoff-Aussage ist falsch,
  und die Modul-5-Lesart (Verzeichnis-Position team-weit = Hauptzweig) gilt
  für die Range noch nicht. Vor der Welle-1-Closure (Paarungen, Register
  team-weit lesen) ist der Push Pflicht, nicht Option.
- `verifizierbar`: ja — `git ls-remote origin main` gegen `git rev-parse HEAD`
- `klasse`: Handoff-Behauptung ohne Beleg (Push-Stand)

### V-2 — Superseded-ADR-Anker an der geschärften Stelle: `service.go` zitiert `ADR-0039` für die Transport-Typen

- `kategorie`: LOW
- `quelle`: [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Supersedes
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md), seit
  `55ce39a`) · Traceability-Regel „keine superseded Referenz"
  (Review-Prüfpflicht, von d-check auf Code-Kommentare nicht gedeckt)
- `pfad`: `internal/application/usecase/capture/service.go:17` („die
  Transport-Typen des Capture Use Cases (`ADR-0039`)") · Plan §8 Zeile 234
  („Ports/Use-Cases in `ADR-0028`/`ADR-0039` verankert") · Plan §3 Zeile 128
  (Original-Zeile)
- `befund`: Genau die Regel, die [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
  geschärft hat, wird im Code-Kommentar noch mit ihrem superseded Anker
  zitiert. Der Kommentar-Inhalt ist mit [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
  konsistent (Definition am Port, Aliase), und [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
  trägt die Struktur-Regeln von [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)
  als Rest fort — deshalb kein Verstoß gegen §3.5 (kein Widerspruch, kein
  in-place-Edit einer Accepted-ADR). Aber der Anker an der
  Definitions-Stelle verweist auf die abgelöste Fassung; die Nachfolge-
  Kette trägt der ADR-Index, nicht der Code. Korrektur wirkt nur vorwärts
  (nächster Lauf oder §7-Vermerk) — der Plan-Nachzug Zeile 130 nennt
  [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) bereits korrekt. `source.go:3`
  zitiert [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md) für die
  *Paketstruktur* (Fortgeltungs-Rest) — kein Befund.
- `verifizierbar`: ja — grep über `internal/` und Plan gegen den
  Supersession-Status
- `klasse`: Superseded-ADR-Zitat am geschärften Punkt

### V-3 — Vorbestand aus slice-002: Worktree `/tmp/verify-slice-002-wt` noch registriert

- `kategorie`: INFO
- `quelle`: `git worktree list` · `verify-slice-002.md` behauptete
  „Worktree … nach Abschluss entfernt"
- `pfad`: `/tmp/verify-slice-002-wt` (Stand `4d6ecfa`, detached)
- `befund`: Der Worktree des Vorgänger-Verifier-Laufs ist nicht entfernt —
  der Bericht behauptete es. Kein Einfluss auf diesen Slice (Hauptbaum
  clean, alle Läufe dieses Laufs im eigenen, entfernten Worktree);
  Hauskeeping für den Planner: `git worktree remove /tmp/verify-slice-002-wt`.
- `verifizierbar`: ja — `git worktree list`
- `klasse`: Verifier-Seiteneffekt aus Vorgänger-Lauf, Bericht-Diskrepanz

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS.md §3.1)** — alle go-Belege
  dieses Laufs liefen im gepinnten Toolchain-Container
  (`golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125`,
  `--network none`, `/src:ro`, GOCACHE/GOMODCACHE im Container); kein
  Host-Toolchain-Aufruf, kein Schreibzugriff auf den Hauptbaum (Mutations-
  Proben im `/tmp`-Worktree, entfernt; `git status` clean)
- geprüft, ohne Befund: **§1-Abgrenzung** — kein reales PostgreSQL-Adapter-
  Artefakt, keine Transaktions-/Spooling-Logik, keine Consumer-/Retention-
  Behandlung in der Range; die drei Ausschlüsse sind im Diff unberührt
- geprüft, ohne Befund: **Traceability der Range** — jeder Commit trägt
  mindestens eine `LH-*`-/`ADR-*`-Kennung (doc-commits 0 Befunde über die
  volle Range `a227033..c263797`); alle genannten IDs existieren; keine
  `SPEC-*`/`ARC-*` in Fix- und Folge-Messagen (`1983e84` ist F-3, vorwärts
  verkörpert); doc-immutable 0 Befunde
- geprüft, ohne Befund: **Spec-Stratum** — die Code-Range berührt keine
  Spec-Datei; die Kontrakt-Zeilen zitieren [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), [`SPEC-001`](../../spec/pflichtenheft.md)/008 als
  Rang-Zeiger, ohne das Technik-Stratum zu erweitern
- geprüft, ohne Befund: **Kommentar-Klassen im Fix (AGENTS.md §3.7)** — die
  F-4-Grenz-Zeile am Inbound-Port (Grenze/Zusage im Indikativ, Grund über
  [`ADR-0012`](../../spec/pflichtenheft.md)), die F-6-Idempotenz-Zeile am Store-Port
  (Zusage), der F-5-Kommentar (Zustand, kein Chronik-Text) — kein
  abwesender Text, keine verworfene Alternative im Konjunktiv
- geprüft, ohne Befund: **a-check-Abdeckung** — `make gates` (a-check im
  Bündel) 0 Befunde am HEAD; der `app`-Glob matcht seit diesem Slice
  echten Content (`usecase/capture`), keine undeclarierte Kante entstanden
  (Fitness-Hälfte von [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) maschinell
  belegt)
- geprüft, ohne Befund: **Verkörperungs-Anker** — `c263797` trägt die Regel
  mit Herkunfts-Anker `seit slice-003` am Zielort
  (`.claude/commands/implement-slice.md:51-52`, grep-Beleg)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Handoff-Behauptung ohne Beleg
(Push-Stand) · Superseded-ADR-Zitat am geschärften Punkt ·
Worktree-Vorbestand aus Vorgänger-Lauf

**Zusammenfassung DoD:** **6/11 Punkte jetzt erfüllt** (Items 1–6, Belege
selbst gefahren) · 1 entfällt (Reconciliation-Register) · 3 erfüllen sich
erst bei Closure (Closure-Notiz §7, Beobachtungs-Register inkl.
`evidence/slice-003.md` + 3×-Ausgang, Endwert der beiden §6-Risiko-Ausgänge)
· 1 an die Welle-1-Closure delegiert (drei Paarungen, DoD-Wortlaut).
Abweichungen: V-1 (ungepusht, Handoff-Behauptung falsch) und V-2
(superseded-ADR-Zitat) — kein DoD-Bruch.

**Review-Blockierpunkte:** alle drei erfüllt — Plan-Nachzug + [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
(`4a41522`, `55ce39a`), F-3-Verkörperung (`c263797`, Herkunfts-Anker
`seit slice-003`), F-4-Nacharbeit (`9bfcd7a`, Grenze + Test, rote Probe
durch eigene Replikation gedeckt). F-7 (Register) im Plan-Nachzug
berichtigt; F-9 (3×) ist Closure-Sache, Zähler-Kontext oben notiert.

## Verdikt

**Merge-blockierend:** nein — alle jetzt prüfbaren DoD-Punkte sind durch
eigene Sensor-Läufe belegt; `make gates` (inkl. a-check mit `app`-Glob auf
echtem Content) ist grün am HEAD, `go test`/`go vet`/`gofmt` im gepinnten
1.27-Container sind clean, die Persist-before-ACK-Ordnung ist am Service an
einer Stelle getragen und durch eigene Mutations-Proben (A, B3) gedeckt,
und der Fix trägt exakt F-4/F-5/F-6. Die Konflikt-Sequenz ist mit
übergabe-Artefakten abgeschlossen (Plan-Nachzug, [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
als Architect-Verdikt, Klasse-B-Verkörperung mit Herkunfts-Anker).

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`, kein Befund gegen den Handoff):** DoD-Punkte 7, 9 und der
Endwert von 10 sind vor dem `git mv` nach `done/` zu erbringen —
Closure-Notiz mit Lerneintrag (§7, inkl. der Review-Finding-Klassen, V-1
und V-2), Register-Ausgang für `BEO-PGC/a-check-null-abdeckung` bei 3×
(`evidence/slice-003.md` + Ausgang im Lese-Schritt), Endausgang für beide
§6-Risiken. **Zusätzlich zu klären vor der Welle-1-Closure:** V-1 — die
zehn Commits sind zu pushen, bevor team-weite Aussagen (Paarungen,
Register-Lese-Schritt) auf dem gemergten Stand gelten.

**Übergabe:** Bericht an den Planner. Keine Reparaturen. V-2 als
vorwärts-Korrektur beim nächsten Code-Kontakt (Zitat auf
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) umstellen) oder als §7-Zeile;
V-3 als `git worktree remove` in der Closure-Handarbeit. F-8
(a-check-Grenze Inbound/Outbound, review-slice-003) bleibt
Architect-Übergabe aus dem Review.