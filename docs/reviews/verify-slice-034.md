# Verifier-Report: slice-034 — 2026-09-13

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §1 (Ziel/Abgrenzung),
§2 (Definition of Done, aktueller Stand nach `6ef4aec`), §6 (Risiken,
noch ohne Ausgang), §8 (Sub-Area-Sichtung), sowie DoD-/Spec-Konformität
gegen [`LH-FA-CFG-003`](../../spec/lastenheft.md), [`LH-FA-CFG-004`](../../spec/lastenheft.md)
(`spec/lastenheft.md`) und den Review-Report
[`review-slice-034.md`](review-slice-034.md) (0 Findings, Verdikt
merge-blockierend: keins).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen oder ausgeführt: vollständiger Slice-Plan
(§1/§2/§3/§6/§7/§8), Review-Report, Testdatei-Volltext (Zeilen 478–620),
`tools/harness/run-integration-tests.sh` (`-run`-Muster), `spec/lastenheft.md`
(`LH-FA-CFG-003`/`004`-Akzeptanzkriterien), die pfadseitig zitierten
Beobachtungsverzeichnisse, `make gates` einmal und `make test-integration`
**dreimal in Folge** selbst gestartet, `git status`/`git diff` am Ende
sauber.

**Gegenstand:** `f3a45f6` (neuer Testfall
`TestMVPActiveTablesViewMatchesActivationState`, `-run`-Muster-Nachzug),
`6ef4aec` (DoD-Häkchen, Plan-Nachzug §3), `c0e0e5d` (Review-Report).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-034-black-box-status-liste.md`
  (vollständig, §1–§8, aktueller Stand)
- `docs/reviews/review-slice-034.md` (0 Findings)
- `spec/lastenheft.md` (`LH-FA-CFG-003`, `LH-FA-CFG-004`,
  vollständige Akzeptanzkriterien-Blöcke)
- `test/integration/integration_test.go` (`newMVPEnv`,
  `TestMVPActivationState` Zeile 478–536, `TestMVPActiveTablesViewMatchesActivationState`
  Zeile 543–606, `TestMVPDisableRetainedState` Zeile 615ff.)
- `internal/application/usecase/list/service.go` und
  `service_test.go` (`TestListTablesEmpty`, Go-Use-Case-Ebene)
- `tools/harness/run-integration-tests.sh` (`-run`-Muster Zeile 241,
  `feed_mvp_idle`-Anlage Zeile 163, `compose.yaml`-`CDC_TABLES`-Bindung)
- `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/`,
  `.../verwaltung-keine-sql-administration/` (Existenzprüfung)
- `Makefile` (Ziel-Existenz-Prüfung für `doc-commits`/`doc-immutable`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify`: v6.5.0, 54 Dateien OK · `d-check` Standardlauf: 288 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 288 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-integration` (Lauf 1/3) | alle 7 Compose-Testfälle PASS, `TestMVPActiveTablesViewMatchesActivationState` PASS (0,02s), CLI-Rundlauf belegt | **0** |
| `make test-integration` (Lauf 2/3) | alle 7 PASS, `TestMVPActiveTablesViewMatchesActivationState` PASS (0,03s) | **0** |
| `make test-integration` (Lauf 3/3) | alle 7 PASS, `TestMVPActiveTablesViewMatchesActivationState` PASS (0,03s) | **0** |
| `git status`/`git diff` (am Ende dieses Laufs) | einzige Restspur `tools/schema/plan.yaml` (erwartetes d-migrate-Rollout-Artefakt, bereits im Review-Report vermerkt) — per `git checkout --` zurückgesetzt, danach sauber | — |

**Zu `make doc-commits`/`make doc-immutable`:** Beide Targets existieren in
diesem Repos `Makefile` **nicht** (eigene Prüfung: `grep -n
"^[a-zA-Z_-]*:" Makefile`) — nach `AGENTS.md` §4 werden nur real
existierende Targets genannt. Die inhaltliche Traceability-Deckung läuft
über `commit-traceability` (Teil von `make gates`, oben grün); eine
dedizierte ADR-Immutabilitäts-Prüfung ist hier ohne Berührungspunkt, da
keiner der drei slice-034-Commits eine ADR-Datei ändert (eigene `git show
--stat` über alle drei Commits).

`make test-integration` wurde entgegen der bloßen Implementer-/Reviewer-
Behauptung in diesem Lauf **selbst dreimal in Folge** ausgeführt (nicht
übernommen) — das ist die häufigste Verifier-Lücke, und genau sie wird
hier geschlossen.

## Prüfpunkt 1 — Black-Box-Reinheit des neuen Testfalls (eigenständig am Code nachvollzogen)

**Bestanden.** `TestMVPActiveTablesViewMatchesActivationState`
(`integration_test.go:551–606`) liest beide neuen Werte über
`env.pool.QueryRow` gegen `cdc.active_tables` — rohes SQL, kein
`list.NewListTablesService`-Import, kein `status`-Paket-Import am
Lesepfad selbst. Der `status.NewGetStatusService`-Aufruf (Zeile 578–581,
596–599) folgt *nach* der jeweiligen SQL-Lesung als unabhängiger
Kontrollwert (`enabled`, `idle`), in eigenen `Fatalf`-Bedingungen — keine
Vermischung der beiden Lesewege in einer Assertion. Deckt sich mit dem
Reviewer-Negativbefund.

## Prüfpunkt 2 — DoD-Punkt 1/2 gegen `LH-FA-CFG-003`/`004` real erfüllt?

**Ja, mit einer benennenswerten Präzisierung (kein Blocker).**

- **Happy Path** (`LH-FA-CFG-003`: „Given aktiviert, dann `aktiviert`
  gemeldet"; `LH-FA-CFG-004`: „Given `t1`/`t2` aktiviert, dann Liste
  enthält mindestens beide"): Der Testfall prüft `feed_mvp_full` real
  über `cdc.active_tables` (Zeile 566–576) — `activeCount == 1`, gefiltert
  auf `source_table_id`. Erfüllt für den geprüften Einzelfall.
- **Boundary `LH-FA-CFG-003`** (`Given t wurde nie aktiviert, dann „nicht
  aktiviert"`): `feed_mvp_idle` wird nie per `EnableTable` aktiviert
  (eigene Prüfung: `compose.yaml` `CDC_TABLES` enthält nur `feed_mvp_flow`,
  `feed_mvp_full`, `feed_mvp_schema`) — `idleCount == 0` gegen
  `cdc.active_tables`, exakt diese Kriterium-Form. **Voll erfüllt.**
- **Boundary `LH-FA-CFG-004`** (wörtlich: „Given **keine** Tabelle ist
  aktiviert, dann **leere Liste**"): Das ist eine andere Situation als die
  oben getestete — die Situation „gar keine Tabelle aktiviert" ist im
  gemeinsamen Compose-Lauf strukturell nicht herstellbar, weil
  `feed_mvp_flow`/`feed_mvp_full`/`feed_mvp_schema` bereits durch die
  Container-Verdrahtung aktiv sind (`ADR-0028`). Der neue Testfall prüft
  stattdessen, dass eine **einzelne** nie aktivierte Tabelle nicht
  erscheint — das ist eine Instanz von `LH-FA-CFG-003`s Boundary, keine
  wörtliche Instanz von `LH-FA-CFG-004`s Boundary. Der Slice-Kopf (§1)
  formuliert das leicht verkürzend („`LH-FA-CFG-004` Boundary — bei
  keiner aktivierten Tabelle liefert die Sicht eine leere
  Ergebnismenge"), was bei wörtlicher Lesung so klingt, als würde genau
  dieses Szenario hier neu bewiesen.
  **Eigenständig geprüft, ob die wörtliche `LH-FA-CFG-004`-Boundary
  anderswo real abgedeckt ist:** Ja — `TestListTablesEmpty`
  (`internal/application/usecase/list/service_test.go:117`) belegt sie
  bereits auf Go-Use-Case-Ebene mit einem Fake-Adapter (0 aktivierte
  Tabellen → leere Liste), referenziert in `service.go:43` explizit als
  `LH-FA-CFG-004` Boundary. Der SQL-View-Beweis dieses konkreten
  Szenarios (0 Zeilen in `cdc.active_tables` bei komplett leerem
  Aktivierungszustand) bleibt am Compose-Stack ungeführt — strukturell
  bedingt, nicht durch diesen Slice verschuldet, und außerhalb dessen
  explizitem §1-Scope (der nur „eine nie aktivierte Tabelle" verspricht).
  Kein DoD-Verstoß, aber eine Formulierungs-Unschärfe im Slice-Kopf — siehe
  Eigener Befund V-1.

## Prüfpunkt 3 — Testisolation und Plan-Nachzug (eigenständig nachvollzogen)

**Bestanden.** `feed_mvp_full` wird an keiner Stelle deaktiviert (eigene
`grep -n "Disable\|feed_mvp_full"` über die gesamte Testdatei); die
Happy-Path-Abfrage filtert zusätzlich hart auf `source_table_id`. Der
`-run`-Muster-Nachzug in `tools/harness/run-integration-tests.sh:241`
enthält `TestMVPActiveTablesViewMatchesActivationState` — ohne diese
Zeile liefe der Testfall nie in `make test-integration` (Risiko aus
`BEO-PGC/test-runner-stiller-ausschluss`, im Slice-Kopf §3 korrekt als
Plan-Nachzug benannt statt stillschweigend nachgetragen). Der dreifache
eigene `make test-integration`-Lauf bestätigt real: alle drei Läufe
PASS, keine Reihenfolge-Abhängigkeit sichtbar geworden.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt, Stand nach `6ef4aec`)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-FA-CFG-003`/`004` erfüllt: neuer Testfall liest `cdc.active_tables` real über SQL | **bestätigt, mit Präzisierung** | Prüfpunkt 2 — Happy Path und `CFG-003`-Boundary voll erfüllt; die wörtliche `CFG-004`-Boundary (leere Gesamtliste) bleibt strukturell außerhalb des Compose-Stacks, ist aber bereits auf Use-Case-Ebene (`TestListTablesEmpty`) belegt — kein neuer Defekt, siehe V-1 |
| 2 | Vertragstest hält SQL-Sicht gegen `status.NewGetStatusService` | **bestätigt** | Prüfpunkt 1 — derselbe Testfall, beide Lesewege in unabhängigen Assertions |
| 3 | `make gates` grün, `make test-integration` dreimal in Folge grün | **eigenständig real reproduziert** | Sensor-Belege oben — nicht aus Implementer-/Reviewer-Bericht übernommen, selbst dreimal gefahren |
| 4 | Review durchgeführt, Report liegt vor | **bestätigt** | `docs/reviews/review-slice-034.md` (`c0e0e5d`), 0 Findings, Verdikt „keins" merge-blockierend |
| 5 | Doku-Update, falls öffentlicher Vertrag berührt | **bestätigt, entfällt korrekt** | reine Testdatei-/Runner-Skript-Änderung, keine API-/CLI-/Schema-Änderung — eigene `git show --stat` bestätigt |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **noch offen — korrekt, Planner-Closure-Schritt** | §7 ist Platzhalter; das ist der Übergabepunkt an den Planner nach dieser Verifikation, kein Implementer-Versäumnis |
| 7 | Reconciliation-Register, falls einschlägig | **entfällt korrekt** | `../reconciliation.md` existiert nicht (Greenfield, `harness/conventions.md` MR-000) |
| 8 | Beobachtungs-Register fortgeschrieben | **noch offen — korrekt, Planner-Closure-Schritt** | §8 zitiert zwei existierende Verzeichnisse korrekt (eigene Prüfung: beide vorhanden); die eigentliche Fortschreibung (neues Verzeichnis oder `evidence/`-Datei) ist noch nicht erfolgt |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **noch offen — korrekt, Planner-Closure-Schritt** | §6-Risiko (Testisolation) trägt noch `<bei Closure einzutragen>`; Prüfpunkt 3 zeigt, dass das Risiko real nicht eingetreten ist (kein Kollisionsfall in drei Läufen) — Ausgang „entfallen" wäre plausibel, bleibt aber Planner-Entscheidung |
| 10 | Drei Paarungen getragen | **noch offen — korrekt, Welle-11-Closure vorbehalten** | Slice-Kopf nennt `welle-11` als Träger; Paarungen sind laut Modul 6 Sache der nächsten Welle-Closure |

**Zwischenstand:** Die drei Kern-DoD-Punkte, die der Implementer als
erledigt beansprucht (1–3), sind eigenständig real nachgeprüft und
tragen — Punkt 1 mit einer benannten, nicht blockierenden Präzisierung
(V-1). Die fünf offenen Punkte (6, 8, 9, 10 sowie das noch fehlende
Checkbox-Häkchen für Punkt 4 im Plan-Text selbst) sind korrekt der
Planner-Closure vorbehalten (Modul 5/6/8) und kein Implementer- oder
Reviewer-Versäumnis.

## Negativbefunde

- geprüft, ohne Befund: **Kein ADR-Inhalt geändert** — `git show --stat`
  über alle drei slice-034-Commits (`f3a45f6`, `6ef4aec`, `c0e0e5d`):
  keiner berührt eine Datei unter `docs/plan/adr/` (Hard Rule 3.5 gewahrt).
- geprüft, ohne Befund: **Kein öffentlicher Vertrag geändert** — Diff
  beschränkt sich auf `test/integration/integration_test.go` und
  `tools/harness/run-integration-tests.sh`; kein Pfad unter `internal/`
  oder `cmd/`, also kein Image-Rebuild-Trigger (deckt sich mit
  Reviewer-Negativbefund, eigenständig anhand `git show --stat`
  nachvollzogen).
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs
  tragen `LH-FA-CFG-003`/`004`, kein `SPEC-*`/`ARC-*` im Betreff;
  `commit-traceability` (Teil von `make gates`) lief in diesem Lauf über
  `HEAD~5..HEAD` grün.
- geprüft, ohne Befund: **§8-Register-Zitate existieren** —
  `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/`
  und `.../verwaltung-keine-sql-administration/` beide vorhanden (eigene
  `find`).
- geprüft, ohne Befund: **`make gates`** — ein eigener Lauf, 0 Befunde
  (288 Dateien Standardlauf, 288 Dateien Range-Lauf, `commit-traceability`
  OK, `a-check` 0 Befunde).
- geprüft, ohne Befund: **`make test-integration` dreimal in Folge** —
  drei eigene, unabhängige Läufe, alle 7 Testfälle PASS in jedem Lauf,
  `TestMVPActiveTablesViewMatchesActivationState` PASS in allen drei
  (0,02s/0,03s/0,03s).
- geprüft, ohne Befund: **Arbeitsverzeichnis am Ende dieses Laufs** —
  einzige Restspur war das erwartete `tools/schema/plan.yaml`-Artefakt
  (Rollout-Nebenwirkung von `make test-integration`, bereits im
  Review-Report vermerkt), per `git checkout --` zurückgesetzt;
  `git status --porcelain` danach leer.

## Eigene Befunde

### V-1 — Slice-Kopf §1 formuliert die `LH-FA-CFG-004`-Boundary-Zuordnung enger als tatsächlich getestet

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-034-black-box-status-liste.md`
  §1, Klammerzusatz „(`LH-FA-CFG-003` Boundary, `LH-FA-CFG-004` Boundary —
  bei keiner aktivierten Tabelle liefert die Sicht eine leere
  Ergebnismenge)"
- `befund`: Siehe Prüfpunkt 2. Der neue Testfall beweist, dass eine
  **einzelne**, nie aktivierte Tabelle nicht in `cdc.active_tables`
  erscheint — das ist wörtlich `LH-FA-CFG-003`s Boundary-Kriterium
  (Status einer bestimmten Tabelle `t`). `LH-FA-CFG-004`s Boundary ist
  wörtlich eine andere Situation: die **gesamte** Liste ist leer, weil
  **keine** Tabelle aktiviert ist. Diese Situation ist im Compose-Stack
  strukturell nicht herstellbar (andere Tabellen sind durch die
  Container-Verdrahtung bereits aktiv) und wird vom neuen Testfall nicht
  geprüft. Die Formulierung im Slice-Kopf erweckt durch die Nennung
  beider IDs nebeneinander den Eindruck, beide Boundary-Kriterien seien
  hier gleichermaßen real bewiesen.
- `verifizierbar`: ja — Wortlaut des Slice-Kopfs, Wortlaut von
  `LH-FA-CFG-004` in `spec/lastenheft.md` Zeile 236f., Testcode-Zeilen
  592–606.
- **Kein DoD-Defekt:** Die wörtliche `LH-FA-CFG-004`-Boundary ist bereits
  real belegt — nur nicht von diesem Slice und nicht am Compose-Stack,
  sondern von `TestListTablesEmpty` auf Go-Use-Case-Ebene (siehe
  Prüfpunkt 2). Es gibt keine unbelegte Lücke in der Spec-Konformität,
  nur eine unscharfe Zuordnung im Slice-Kopf.
- **Für die Closure:** kein Blocker. Empfehlung an den Planner: bei der
  Closure-Notiz (§7 „Was ging anders als geplant") oder als kleine
  Präzisierung in §1 festhalten, dass die SQL-View-Instanz von
  `LH-FA-CFG-004`s Boundary (leere Ergebnismenge bei **komplett**
  fehlender Aktivierung) am Compose-Stack strukturell offen bleibt —
  damit ein späterer Leser des Slice-Kopfs nicht annimmt, dieses
  spezifische Szenario sei am realen Stack bewiesen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** Die drei implementierten Kern-Punkte (1–3) sind
eigenständig real nachgeprüft; Punkt 1 trägt eine benannte, nicht
blockierende Präzisierung (V-1). Die fünf noch offenen Punkte (4-Häkchen
im Plan-Text, 6, 8, 9, 10) sind korrekt der anstehenden Planner-Closure
vorbehalten und kein Implementer-/Reviewer-Versäumnis. Kein DoD-Defekt im
Sinn eines unbelegten „bestätigt"-Punkts.

## Verdikt

**DoD-Konformität (soweit vom Implementer beansprucht): bestätigt.** Der
neue Testfall liest `cdc.active_tables` real über rohes SQL, ohne
Vermischung mit dem internen Lesepfad, und hält beide Werte im selben
Testfall gegeneinander. `LH-FA-CFG-003`s Happy Path und Boundary sind
damit jetzt zusätzlich am realen Compose-Stack über die SQL-Sicht
bewiesen. `LH-FA-CFG-004`s Happy Path ist ebenfalls real über die Sicht
belegt; `LH-FA-CFG-004`s wörtliche Boundary (leere Gesamtliste) bleibt am
Compose-Stack strukturell ungeführt, ist aber bereits real auf
Go-Use-Case-Ebene bewiesen (`TestListTablesEmpty`) — kein neuer Defekt,
nur eine im Slice-Kopf zu eng formulierte Zuordnung (V-1, LOW).

**Eigene Reproduktion:** `make gates` einmal, `make test-integration`
**dreimal in Folge**, alle vier Läufe grün, real in diesem Kontext
ausgeführt — nicht aus dem Implementer- oder Reviewer-Bericht übernommen.

**Plan-vs-Code-Diff:** Der Diff (`f3a45f6`) hält sich exakt an §3 des
Slice-Plans (Testdatei + Runner-Skript-Nachzug), keine unangekündigte
Schicht- oder Umfangs-Erweiterung. §1s Abgrenzung (Negative-Fall,
aktiviert-vs-retained, CLI-Exposition) wird nicht berührt.

**Closure-Bereitschaft:** Die Verifikation der real gelieferten Punkte
ist abgeschlossen und positiv. Vor dem `git mv` nach `done/` fehlen noch
die Planner-Closure-Schritte: DoD-Häkchen 4 im Plan-Text nachziehen,
§7-Closure-Notiz schreiben (inkl. Behandlung von V-1), Beobachtungs-
Register fortschreiben, §6-Risiko-Ausgang eintragen (Prüfpunkt 3 legt
„entfallen" nahe, Urteil bleibt beim Planner) — die drei Paarungen bleiben
regulär der `welle-11`-Closure vorbehalten.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Slice-Plan,
Register und Roadmap wurden von diesem Lauf nicht verändert
(`git status`/`git diff` am Ende sauber, bis auf den erwarteten,
zurückgesetzten `tools/schema/plan.yaml`-Rollout-Rest).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: v6.5.0, 54 Dateien OK; `d-check` Standardlauf:
288 Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`:
288 Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits;
`a-check`: 0 Befunde). `make test-integration` dreimal in Folge
ausgeführt, Exit 0 in jedem Lauf, `TestMVPActiveTablesViewMatchesActivationState`
PASS in allen drei. `make doc-commits`/`make doc-immutable` existieren in
diesem Repo nicht (eigene Makefile-Prüfung) — Traceability-Deckung läuft
über `commit-traceability` (Teil von `make gates`), ADR-Immutabilität ist
hier ohne Berührungspunkt (kein Commit ändert eine ADR-Datei). `git
status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
