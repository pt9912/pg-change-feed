# Verifier-Report: slice-008 — 2026-09-10

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
**10 Items**), §3 (Plan-vs-Code, Range `296faee..HEAD` inkl. Plan-Nachzug
`81fc8f2`), §6 (Risiko-Ausgänge) und ADR-Konformität
([`ADR-0028`](../plan/adr/0028-inbound-use-cases.md) ·
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) ·
[`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) ·
[`ADR-0026`](../plan/adr/0026-composition-root.md) ·
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) ·
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail (Reviewer,
`review-slice-008.md`, Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`../plan/planning/in-progress/slice-008-cdc-verwaltung-use-cases.md` (Welle
welle-3) · Range `296faee..HEAD` (**HEAD am Prüfzeitpunkt `81fc8f2`** — der
Plan-Nachzug liegt über dem Review-Range-Head `4d072c9`; alle Sensor-Läufe
unten fahren am HEAD `81fc8f2`) · Implementer-Commits `814a166`, `4d06164`,
`c7a725b`, `54ab209`, `82d3250`, `4d072c9` · Review-Report `bd0efde` ·
Fix-Commits `bbc9bc2` (F-2), `4a935dd` (F-4), `9e7979d` (F-6/F-7/F-9) ·
Image-Beleg `05fe8f9` · Plan-Nachzug `81fc8f2` (F-1/F-8 + F-3-Risiko).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über stdout
(Tabelle unten). Die Plan-Datei und der Code blieben unberührt; die einzige
Schreibaktion dieses Laufs ist dieser Report (`docs/reviews/verify-slice-008.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am HEAD `81fc8f2` (inkl. Plan-Nachzug) ·
  `review-slice-008.md` (F-1…F-10, committet in `bd0efde`) · Fix-Commits im
  Volltext
- `../../spec/lastenheft.md` ([`LH-FA-CFG-001`](../../spec/lastenheft.md)…004
    samt Akzeptanzkriterien), [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md)/[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)/[`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md)/[`ADR-0026`](../plan/adr/0026-composition-root.md) im Volltext
- Beobachtungs-Register (`BEO-PGC/*`, drei Einträge) ·
  `../plan/planning/welle-3.md` (flach — Repo **mit** Wellen-Betrieb) ·
  `harness/image-hash.txt`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `81fc8f2`) | `baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` · `d-check: 114 Datei(en) geprüft, 0 Befund(e)` (voll) · `d-check … --range HEAD~5..HEAD`: `0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | **0** |
| `make test` (netzlos, gepinnter Container `golang:1.27-alpine@sha256:cf6fca…`) | **15 Pakete `ok`** — inkl. `usecase/{enable,disable,status,list}`, `postgresstorage` (mit `identifier_test.go`), `bootstrap`, `test/integration` (Skip ohne DSN) | **0** |
| `make test-store` (Stichprobe, reale PostgreSQL, gepinnte Digests) | **grün** — `postgresstorage 1.654s` (reale DB-Läufe gegen `0.005s` treiberfrei im `make test`-Lauf): Register/Unregister/Publish/Unpublish/List, `TestRetainedStateView`, `TestPublishedMissingPublication`, `TestIdentifierVerweigerung`; PostgreSQL-Digest `ab07f3725c1f…` | **0** |
| `docker images ghcr.io/pt9912/pg-change-feed` | lokales `:dev`-Image `de204223c27c…` == committeter Lauf-Beleg `harness/image-hash.txt` (`05fe8f9`, nach allen drei Fix-Commits) — [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Semantik (Lauf-Beleg, kein Inhalts-Fingerabdruck) | — |
| `git status --porcelain` (vor diesem Report) | leer — der Status-Snapshot des Auftrags (Dockerfile, `0044`) war veraltet; Arbeitsbaum clean am Prüf-HEAD | — |

Nicht selbst gefahren: `make test-integration` (kein DoD-Item dieses Slice;
`TestMVPDisableRetainedState`/`TestMVPActivationState` skippen ohne
`CDC_INTEGRATION_DSN` — deren Belege tragen die Adapter-Tests aus
`make test-store` am selben Bestand; die Capture-Zeit-Hälfte der
Deaktivierung ist als offenes Risiko benannt, siehe §6-(b) unten) ·
`make doc-commits`/`doc-immutable` (das Standing-Gate-Fenster
`HEAD~5..HEAD` deckt den Range `296faee..HEAD` nicht vollständig; die
Traceability aller Commits ist über den vollständigen d-check-Lauf oben und
den Review-Negativbefund „Traceability aller sechs Commits" getragen).

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt — 10 Items)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `EnableTableUseCase`/`DisableTableUseCase` am Inbound-Port ([`LH-FA-CFG-001`](../../spec/lastenheft.md)/002) | **bestätigt** | `verwaltung.go:109-118` trägt beide Port-Interfaces; Services `usecase/enable`/`disable` je mit `var _ inbound.…UseCase`-Verankerung; LH-Akzeptanzkriterien testgetragen: Happy (`TestEnableHappyPath`, `TestDisableHappyPath`), Boundary idempotent (`TestEnableIdempotent` über `AlreadyEnabled`, `TestDisableIdempotent`), Negative sichtbar (`TestEnableMissingTable`/`TestDisableMissingTable` über `ErrSourceTableMissing`) — alle 15 Pakete grün in meinem `make test`-Lauf. Aktivierung am Produktionspfad: `wiring.go:167-185` ruft EnableTable vor dem Stream-Start ([`ADR-0026`](../plan/adr/0026-composition-root.md)) |
| 2 | `GetStatusUseCase`/`ListTablesUseCase` (Kanon-Bezeichner je [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md), [`LH-FA-CFG-003`](../../spec/lastenheft.md)/004) | **bestätigt** | `verwaltung.go:120-130` trägt beide unter den Kanon-Namen (`GetStatusUseCase` — der Plan-Text trägt den Namen nach dem Nachzug `81fc8f2` korrekt); `GetStatusService`/`ListTablesService` lesen den Zustand über Bindungs-Zeile **samt** `Published`-Mitgliedschaft — LH-003 Happy/Boundary/Negative (`TestStatusEnabled`/`TestStatusNotActivated`/`TestStatusMissingTable`), LH-004 Happy/Boundary (`TestListTables`, `TestListTablesEmpty`); Negative für LH-004 ist im Lastenheft ausdrücklich „—" |
| 3 | `make gates` grün | **bestätigt** | Vier Gates grün am HEAD `81fc8f2` (Tabelle oben, Exit 0), inkl. [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)-Standing-Gate (5 Commits OK); §2-Häkchen noch offen — Closure-Pflicht |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **bestätigt** | `review-slice-008.md` committet in `bd0efde` (F-1…F-10, Verdikt „merge-blockierend: nein"); Rollenwechsel nach Schritt 8 eingehalten, kein Self-Review |
| 5 | Doku-Update ENV-Container-Vertrag (`compose.yaml`), falls berührt — geprüft: ENV-Vertrag unverändert, Item entfällt | **entfällt — korrekt geprüft** | Meine Stichprobe: `compose.yaml` trägt unverändert `CDC_SOURCE_DSN`/`CDC_SOURCE_ID`/`CDC_PUBLICATION`/`CDC_SLOT`/`CDC_TABLES`; die compose-Änderung im Range ist der Seed-SQL-Entfall (§3-Zeile), keine Vertrags-Änderung. Die `harness/README.md`-Werkzeuge-Zeile trägt den Container-Vertrag seit dem slice-007-Nachzug (V-2 von verify-slice-007 erledigt) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure (normal)** | §7 trägt Platzhalter; Slice korrekt in `in-progress/`. Closure-Lern- und Vermerk-Pflichten unten (Summary) |
| 7 | Reconciliation-Register fortgeschrieben, falls Inventur-Fund | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen; kein Inventur-Fund im Range |
| 8 | Beobachtungs-Register fortgeschrieben | **offene Closure-Pflicht — Ausgang benannt** | Risiko §6 (b) trägt den Ausgang „weiter offen → `BEO-PGC/walsender-wirksamkeit` im Register (bei Closure)" — das Verzeichnis existiert **noch nicht** (Register führt drei Einträge: `d-migrate-nacharbeit`, `adapter-fehler-ausgang`, `a-check-null-abdeckung`); korrekter Vor-Closure-Stand, Beleg `evidence/slice-008.md` fällig beim Übergang |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **Voreinstellung getragen — Endbelege fällig bei Closure** | (a) Aktivierung über Use Case vs. Seed-SQL — **eingetreten**, Träger dieser Slice: belegt — der Runner trägt keine Seed-Aktivierung mehr (`grep seed` im Runner: nur Kommentar „EnableTable Use Case, nicht über Seed-SQL“), die Aktivierung läuft an genau einer Stelle (`wiring.go`), der Runner-Wächter prüft Slot- und Feed-Bestand. (b) Wirksamkeits-Grenze am laufenden Walsender — **weiter offen**, Träger: Register-Eintrag bei Closure (Item 8). Beide Ausgänge stammen aus der geschlossenen Dreier-Menge; der Register-Beleg für (b) fällt bei der Closure an |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (`welle-3.md` flach vorhanden; der Slice trägt `**Welle:** welle-3`): der DoD-Wortlaut weist die Prüfung der nächsten Welle-Closure zu |

**Zwischenstand: 4/10 Punkte jetzt erfüllt** (Items 1–4, Belege selbst
gefahren) · 2 entfallen (Items 5, 7 — beide korrekt geprüft) · 3 erfüllen
sich bei Closure (Items 6, 8, 9-Endbelege) · 1 delegiert (Item 10). Der
Zustand ist der **normale eines Slice in `in-progress/`** — alle §2-Häkchen
stehen offen; die offenen Items sind Closure-Pflichten. Ein neuer
blockierender DoD-Defekt ist **nicht** aufgetreten.

## Plan-vs-Code-Diff (Range `296faee..HEAD`, Plan-Nachzug `81fc8f2`)

**Plan-Nachzug (`81fc8f2`) gegen den Review-Blocker F-1:** bestätigt — §3
trägt jetzt die vollständige Datei-Menge (acht Ergänzungen inkl.
Outbound-Port, `status`/`list`, Adapter mit `queries.go`, `wiring.go`,
`compose.yaml`, `mvp_test.go`); §2 trägt die vier F-8-Defekt-Korrekturen
(Dup-DoD entfernt, Platzhalter aufgelöst, Kanon-Name, Ausschluss Klasse 2);
§6 führt Risiko (b) mit Ausgang (F-3). Die Message trägt
[`LH-FA-CFG-001`](../../spec/lastenheft.md)/[`ADR-0028`](../plan/adr/0028-inbound-use-cases.md).

**Deckung §3 (am HEAD `81fc8f2`) gegen den Code-Stand:** die acht §3-Zeilen
decken alle Produkt-Dateien des Ranges — `verwaltung.go`,
`port/outbound/tableactivation.go`, `usecase/{enable,disable,status,list}/*.go`
(deckt auch die vier `*_test.go`), Adapter `tableactivation.go` +
`queries/queries.go` (deckt `tableactivation_test.go`), `wiring.go`,
Runner, `compose.yaml`, `mvp_test.go`.

**Rest außerhalb §3 (V-2 unten):** `internal/adapters/driving/replication/receive/receive.go`
(M — die Alphabet-Quell-Deklaration des F-6-Fixes `9e7979d`). Eine
Produkt-Datei im Range ohne §3-Zeile; getragen ist die Änderung inhaltlich
(Review F-6-Disposition), aber sie ist Plan-arbeitslos. Die reinen
Test-Dateien (`identifier_test.go` — F-4-Fix) decken sich über die
Komponenten-Zeilen des Adapters.

**Nicht in §3, zulässig:** `harness/image-hash.txt` (Lauf-Beleg,
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md), `05fe8f9` — nach
allen Fix-Commits geordnet, Beleg stimmt mit dem lokalen Image überein,
Tabelle oben) · Review-Report `bd0efde` · Plan-Nachzug `81fc8f2` (Übergabe-
Artefakte).

## Review-Dispositions-Check (F-1…F-10)

| Finding | Disposition | Verdikt |
|---|---|---|
| F-1 (Plan-Nachzug fehlt, 8. Auftreten) | Plan-Commit `81fc8f2` — §3 vollständig | **getragen** (Rest: V-2, `receive.go`) |
| F-2 (Zustands-Doppeldeutigkeit) | Code-Fix `bbc9bc2` + Tests — `Published`-Mitgliedschaft trennt Erfassungs-Zustand von Herkunft am Status- **und** Listen-Pfad (`status/service.go`, `list/service.go`, Adapter `Published`-View); Test-Träger: `TestStatusRetained`, `TestListTablesRetainedPartition`, `TestRetainedStateView` (`make test-store` grün), `TestMVPDisableRetainedState` (Integration); Plan §3-Zeile nennt den Zustands-View | **getragen** |
| F-3 (Walsender-Wirksamkeits-Grenze nur Kommentar) | Plan §6 Risiko (b) mit Ausgang „weiter offen → Register" in `81fc8f2`; der Port-Kommentar bleibt der Code-Träger der Grenze | **getragen als Risiko-Ausgang** — Register-Beleg fällig bei Closure |
| F-4 (Bezeichner-Vertrag ohne Negativtest, 3. Auftreten) | `4a935dd` — `identifier_test.go`: Alphabet-Verweigerung, Längen-Grenze, Publication-DDL vor der Interpolation; läuft grün in meinem `make test`- und `test-store`-Lauf | **getragen** |
| F-5 (`ListTablesUseCase`/`TableActivationPort` außerhalb der ADR-Listen) | **nicht getragen** — korrekt an die Closure verwiesen (Closure-Vermerk §7 oder Folge-ADR); kein ADR im Range berührt ([`ADR-0028`](../plan/adr/0028-inbound-use-cases.md) unverändert, Accepted-immutable) | **offen — Closure-Vermerk fällig (V-3)** |
| F-6 (Alphabet doppelt geführt) | `9e7979d` — `receive.go` deklariert sich als Quelle des Alphabet-Vertrags, der Aktivierungs-Adapter als Halter desselben Ausdrucks (Kopplung benannt, Gewinner deklariert) — im Code gelesen | **getragen** |
| F-7 (Konjunktiv-Klausel) | `9e7979d` — Enable-Kommentar trägt jetzt die geltende Zusage indikativ („die Bindungs-Zeile läuft vor der Publication …") — im Code gelesen | **getragen** |
| F-9 (Drei Verbindungen ohne Grenze) | `9e7979d` — `wiring.go:156-161` trägt die Abgrenzung indikativ — im Code gelesen | **getragen** |
| F-8 (Plan-Form-Defekte) | `81fc8f2` — alle vier Punkte im Plan-Diff sichtbar | **getragen** |
| F-10 (Anfangs-Version hart 1, INFO) | kein Träger nötig — geht in die Closure §7 | **notiert** |

Alle vier Review-Blockierpunkte (F-1, F-2, F-3, F-4/F-5-Komplex) sind
beantwortet; F-5 ist bewusst als Closure-Vermerk stehen geblieben und wird
unten als Closure-Pflicht weitergeführt.

## Befunde

### V-1 — §8-vorgelagerte Prüfungen stehen als Vorlagen-Platzhalter

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form:
  Sub-Area-Modus-Begründung (zwei vorgelagerte Prüfungen „stehen in **jedem**
  Slice-Plan") · F-9-Klasse aus review-slice-007 (Platzhalter im aktiven Plan,
  dort als V-3 belegt)
- `pfad`: `../plan/planning/in-progress/slice-008-cdc-verwaltung-use-cases.md`
  §8 („Vorgelagert — Sub-Area-Wahl prüfen: <je berührter Sub-Area …>" und
  „Vorgelagert — offene Beobachtungen sichten: <Register durchgegangen …>")
- `befund`: Die zwei vorgelagerten Prüfungen tragen Vorlagen-Platzhalter statt
  einer aufgezeichneten Antwort — der Sichtungs-Schritt (offene
  Beobachtungen gegen `observations/`) ist nicht belegt, obwohl er der einzige
  Register-Leser auf der Planungsseite ist. Die Modus-Begründung selbst trägt
  den GF-Hinweis („Reiner GF-Hinweis genügt; kein Sub-Area-Block") — nur die
  vorgelagerten Hälfte fehlt. Planner-Arbeit, keine Code-Arbeit; review-slice-008
  F-8 listete die §1/§2-Defekte, nicht §8.
- `verifizierbar`: ja — Lese der Plan-Datei am HEAD `81fc8f2`
- `klasse`: Vorlagen-Platzhalter im aktiven Plan (F-9-Klasse — Sichtungs-Schritt)

### V-2 — Eine Produkt-Datei des Fix-Zugs liegt außerhalb der §3-Tabelle

- `kategorie`: LOW
- `quelle`: Slice-Plan §3 („Wer später mitnimmt …, hat den Plan geändert") ·
  Modul 5
- `pfad`: `internal/adapters/driving/replication/receive/receive.go` (M in
  `9e7979d`, F-6-Fix: Alphabet-Quell-Deklaration)
- `befund`: Der Plan-Nachzug `81fc8f2` schloss die F-1-Lücke für den
  Implementer-Range vollständig, aber der nachfolgende Fix-Zug (`bbc9bc2`,
  `4a935dd`, `9e7979d`) berührt eine Produkt-Datei, die keine §3-Zeile trägt.
  Die Änderung ist Kommentar-Half der F-6-Disposition (Review-getragen), aber
  der Plan-vs-Code-Vergleich weist sie als Deckungs-Rest aus. Bei der Closure
  in §7 („Was ging anders als geplant") zu benennen oder §3 zu ergänzen.
- `verifizierbar`: ja — `git diff 296faee..HEAD --name-only` gegen §3
- `klasse`: Plan-Erweiterung-Klasse (9. Auftreten, minimale Form — einzelne
  Fix-Datei, Review-getragen)

### V-3 — F-5 (`ListTablesUseCase`/`TableActivationPort` außerhalb der ADR-Listen) steht als Closure-Vermerk an

- `kategorie`: LOW
- `quelle`: review-slice-008 F-5 (MEDIUM) · Modul 8 (kein stiller
  ADR-Kanon-Ausbau)
- `pfad`: `verwaltung.go:126-130` (`ListTablesUseCase`) ·
  `port/outbound/tableactivation.go` (`TableActivationPort`) ·
  [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md)-Liste (neun Einträge,
  unverändert) · [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)
- `befund`: Kein Träger existiert bisher: kein Folge-ADR, kein Vermerk, keine
  Closure-Zeile. Der Review-Verdikt-Punkt 4 nennt beide Ausgangs-Wege
  (Folge-ADR oder Closure-Vermerk §7 mit ADR-Bezug). **Notiert, nicht neu
  bewertet** — die Pflicht liegt bei der Closure; ohne einen der beiden
  Träger geht der Slice nicht nach `done/`.
- `verifizierbar`: ja — ADR-Liste gegen den Port-/Use-Case-Baum
- `klasse`: ADR-Liste still erweitert (Ausgang offen bis Closure)

### V-4 — Register-Beleg für §6 (b) und der Capture-Zeit-Beleg von [`LH-FA-CFG-002`](../../spec/lastenheft.md) Happy Path

- `kategorie`: INFO
- `quelle`: Modul 5 („weiter offen" → Beobachtungs-Register) ·
  [`LH-FA-CFG-002`](../../spec/lastenheft.md) Happy Path („fortan keine
  Änderungen mehr erfasst")
- `pfad`: Slice-Plan §6 (b) · `TestMVPDisableRetainedState`
  (`test/integration/mvp_test.go:393 ff.`)
- `befund`: Der Happy Path von [`LH-FA-CFG-002`](../../spec/lastenheft.md) ist
  am Bindungs-/Publication-Zustand belegt (`Unpublish` trägt den Stopp; der
  Integrationstest belegt Retained/Status/Liste nach der Deaktivierung), aber
  **kein Test zählt die Changes nach dem Deaktivierungs-Zeitpunkt** — das
  Capture-Zeitverhalten am laufenden Walsender bleibt die offene Hälfte, und
  genau das benannt der Plan §6 (b): „der Zustands-View (Kataloge) belegt die
  Bindung, nicht das Capture-Zeitverhalten". Die Ehrlichkeit des Risikos ist
  korrekt; der Register-Eintrag (`BEO-PGC/walsender-wirksamkeit` mit
  `evidence/slice-008.md`) ist der fällige Endbeleg bei der Closure. Nicht
  blockierend über die ohnehin offene Closure hinaus.
- `verifizierbar`: ja — Plan §6 gegen `docs/plan/planning/observations/`
- `klasse`: Wirksamkeits-Grenze am laufenden Walsender (Register-Ausgang
  vorbereitet, Beleg fällig)

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS §3.1)** — alle Go-Belege im
  gepinnten Toolchain-Container (`--network none`, Modul-Cache im Volume);
  `make gates`/`test-store` über die gepinnten d-check-, a-check- und
  PostgreSQL-Digests
- geprüft, ohne Befund: **[`ADR-0026`](../plan/adr/0026-composition-root.md)**
  — `wiring.go` verdrahtet den Aktivierungs-Adapter und ruft EnableTable vor
  dem Stream-Start; `main` im Range unberührt; `a-check` 0 Befunde im
  `make gates`-Lauf
- geprüft, ohne Befund: **[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)**
  — Transport-Typen am Inbound-Port definiert, vier Alias-Blöcke je
  Use-Case-Paket (`= inbound.X`), keine Neuedefinition in den Services
- geprüft, ohne Befund: **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)**
  — der Bezeichner-Vertrag endet sichtbar über `ErrActivationConfiguration`
  (Klasse `configuration`) vor dem ersten SQL-Aufruf, jetzt mit
  Negativtest-Träger (`4a935dd`); `ErrSourceTableMissing` als Port-Sentinel
  über alle drei Negative-Pfade
- geprüft, ohne Befund: **ENV-Vertrag (`compose.yaml`)** — die fünf
  `CDC_*`-Variablen unverändert gegen den slice-007-Stand; die
  compose-Änderung im Range trägt nur den Seed-SQL-Entfall
- geprüft, ohne Befund: **WIP-Limit 1** — nur slice-008 in `in-progress/`
- geprüft, ohne Befund: **Arbeitsbaum** — `git status --porcelain` leer vor
  diesem Report; die einzige Schreibaktion dieses Laufs ist der Report-Commit
- geprüft, ohne Befund: **Traceability der Fix-Zug-Commits** — `bbc9bc2`,
  `4a935dd`, `9e7979d`, `05fe8f9`, `81fc8f2` tragen je [`LH-*`](../../spec/lastenheft.md)-/[`ADR-*`](../plan/adr/README.md)-Kennungen,
  keine Struktur-ID im Betreff (`commit-traceability` grün im `make gates`-Lauf)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Vorlagen-Platzhalter im aktiven Plan (§8
Sichtungs-Schritt — F-9-Klasse) · Plan-Erweiterung-Klasse, minimale Form
(`receive.go` außerhalb §3, Review-getragen) · ADR-Liste still erweitert
(F-5, Ausgang offen bis Closure) · Capture-Zeit-Hälfte der Deaktivierung
(§6 (b), Register-Beleg fällig).

**Zusammenfassung DoD:** **4/10 Punkte jetzt erfüllt** (Items 1–4, Belege
selbst gefahren: `make gates` grün am HEAD `81fc8f2` mit Exit 0, 15/15
Test-Pakete netzlos, `test-store` grün am realen PostgreSQL, Image-Beleg
`de204223…` lokal verifiziert) · Items 5 und 7 **entfallen korrekt** ·
Items 6/8/9 erfüllen sich bei Closure · Item 10 delegiert an die
Welle-3-Closure. Kein neuer DoD-Defekt.

**Alle Review-Blockierpunkte sind beantwortet:** F-1 (`81fc8f2`,
vollständiger §3-Nachzug), F-2 (`bbc9bc2` mit Status-/Listen-Trennung über
`Published`, testgetragen), F-3 (§6-Risiko (b) mit Ausgang), F-4
(`4a935dd`); F-6/F-7/F-9 (`9e7979d`) und F-8 (`81fc8f2`) getragen, F-5
steht als Closure-Vermerk an (V-3).

## Verdikt

**Merge-blockierend:** nein — kein DoD-Defekt, kein HIGH-Finding; die
Verwaltungs-Use-Cases tragen [`LH-FA-CFG-001`](../../spec/lastenheft.md)…004
mit LH-Akzeptanzkriterien je Test-Ebene, `make gates` läuft grün am HEAD, und
die vier Review-Blockierpunkte sind getragen.

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`, plus zwei echte Punkte):**

1. **V-3 (F-5):** Closure-Vermerk §7 mit ADR-Bezug **oder** Folge-ADR für
   `ListTablesUseCase`/`TableActivationPort` — ohne Träger kein Übergang nach
   `done/`.
2. **V-1:** §8-vorgelagerte Prüfungen füllen (Sub-Area-Wahl, offene
   Beobachtungen sichten) — Planner-Arbeit vor der Closure.
3. **Closure-Pflichten:** §2-Häkchen, Closure-Notiz §7 mit
   Steering-Loop-Lerneintrag (Kandidaten: Plan-Erweiterung-Klasse — 8./9.
   Auftreten, Sequenz läuft), Register-Beleg für §6 (b)
   (`BEO-PGC/walsender-wirksamkeit`, `evidence/slice-008.md` — V-4),
   `receive.go`-Rest in §7 benennen (V-2), Risiko-Ausgänge je §6, Paarungen
   an die Welle-3-Closure delegiert (Item 10).

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert.