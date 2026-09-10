# Verifier-Report: slice-010 — 2026-09-10

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
**11 Zeilen** in §2, davon 10 materiell zählbar — siehe DoD-Prüfung), §3
(Plan-vs-Code, Range `79c0c12..HEAD`, inkl. beider Nachzüge `9a7b3c1` und
`5365ee2`), §6 (Risiko-Ausgang) und Entscheidungs-Konformität
([`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md) ·
[`ADR-0018`](../plan/adr/0018-sql-driving-adapter.md) (superseded) ·
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) ·
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail über die
F-1-Disposition hinaus (Reviewer, `review-slice-010.md`, Verdikt dort),
realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`../plan/planning/in-progress/slice-010-lese-vollabdeckung-sql.md` (Welle
welle-3) · Range `79c0c12..HEAD` (**HEAD am Prüfzeitpunkt `5365ee2`**) ·
Implementer-Commits `9a7b3c1` (Plan-Nachzug), `fe1ff63` (SQL-Views),
`c534fec` (SQL-View-Tests) · Planner-Zug `c57215a` (Register-Beleg) ·
Review-Report `28e48de`/`c292bca` (F-1 HIGH, F-2 INFO) · Architect-Sequenz
`ab00601` (`ADR-0046` Accepted, `Supersedes ADR-0018`), `672a2d7`
(`ADR-0046`-Fix für `make gates`), `5365ee2` (Architektur-Korrektur +
Bezug-/Plan-Nachzug für `ADR-0046`).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren, inklusive einer eigenen
Mutationsprobe gegen einen frischen Testcontainer (die drei vom
Implementer berichteten Mutationsproben hatte der Reviewer nur als
plausibel bewertet, nicht selbst ausgeführt — Aufgabe explizit an den
Verifier delegiert). Die Plan-Datei und der Code blieben unberührt; die
einzige Schreibaktion dieses Laufs ist dieser Report
(`docs/reviews/verify-slice-010.md`).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am HEAD `5365ee2` · `review-slice-010.md` (F-1 HIGH,
  F-2 INFO, committet in `28e48de`/`c292bca`)
- `ADR-0046` (Accepted, Supersedes `ADR-0018`) und `ADR-0018` (Status
  jetzt `Superseded by ADR-0046`) im Volltext
- `spec/lastenheft.md` (`LH-FA-REA-001`…006, `LH-FA-SST-002`),
  `spec/architecture.md` §1/§2/§4 (korrigierter SQL-Kanal-Abschnitt zu
  `LH-FA-REA-002`, Commit `5365ee2`)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
  (`observation.md`, `state.md`, `evidence/slice-006.md`,
  `evidence/slice-010.md`) — 2× wie im Plan zitiert
- `tools/schema/nacharbeit-views.sql`, `tools/schema/schema.yaml`
  (Kommentar-Block), `internal/adapters/driven/postgresstorage/sqlviews_test.go`,
  `Makefile` (`schema-rollout`-Target) im Volltext
- `.harness/baseline/v6.5.0/regelwerk/modul-05-planning-harness.md`,
  `modul-06-roadmap.md`, `modul-08-agentenrollen.md` — DoD-, Register- und
  Rollen-Regeln

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `5365ee2`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 127 Datei(en) geprüft, 0 Befund(e)` (voll) · `d-check … --range HEAD~5..HEAD`: `0 Befund(e)` · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | **0** |
| `make test` (netzlos, gepinnter Container `golang:1.27-alpine@sha256:cf6fca…`) | **19 Pakete `ok`/`[no test files]`**, keine Fehlschläge | **0** |
| `make schema-validate` | `d-migrate schema validate` — 7 Tabellen, 26 Spalten, 1 Index, 6 Constraints, `Validation passed: 0 warning(s)` | **0** |
| `make test-store` (reale PostgreSQL, gepinnte Digests, Rollout inkl. beider Nacharbeit-Schritte vor dem Testlauf) | **grün** — `postgresstorage 2.758s`; Rollout-Log zeigt beide psql-Nacharbeit-Schritte (`nacharbeit-operation-check.sql`, `nacharbeit-views.sql`, drei `CREATE VIEW`); Rollout-Report (`plan.yaml`) und Rollback-Artefakt (`down.sql`) erzeugt | **0** |
| Eigene Mutationsprobe (isolierter Testcontainer `cdc-verify-test-pg`, Netz `cdc-verify-test`, vollständiger Schema-Rollout inkl. der drei Views) — Baseline-Lauf `go test -run TestChangesViewCarriesRangeLimitAndFilter` | `--- PASS (0.06s)` gegen die reale, unveränderte `cdc.changes`-View | **0** |
| Dieselbe Probe nach Mutation `JOIN cdc.source_table st ON false` in `cdc.changes` (dritte der drei vom Implementer berichteten Mutationen) | `--- FAIL`: `sqlviews_test.go:246: Bereich [1100,1300) Limit 1 = [], wollen genau [vt-changes-c-100]` — der Inner Join liefert bei `ON false` keine Zeile, der Test schlägt real fehl | **1** |
| Dieselbe Probe nach Rückbau der Mutation (View auf Originalform zurückgesetzt) | `--- PASS (0.06s)` — Wächter-Charakter bestätigt, kein Seiteneffekt zurückgelassen | **0** |
| Mutmaßlicher Testcontainer-Nebeneffekt: `git status --porcelain` nach allen Läufen | `M tools/schema/plan.yaml` — Report-Artefakt der `schema-rollout`-Sensor-Läufe, bereits vor diesem Verifikationslauf modifiziert (siehe Git-Status zu Sitzungsbeginn); kein Plan-/Code-Inhalt geändert | — |

**Mutationsprobe im Detail (Nachvollzug der Implementer-Behauptung, siehe
`review-slice-010.md`, Negativbefund „geprüft, plausibel, nicht im
Review-Lauf reproduziert"):** eigener Testcontainer `cdc-verify-test-pg`
(gepinntes `postgres:18-alpine`-Image, wie `run-store-tests.sh`) auf
eigenem Docker-Netz `cdc-verify-test`, vollständiger `make schema-rollout`
gegen diese Instanz (beide Nacharbeit-Schritte inklusive der drei Views).
Danach: (1) Baseline-Lauf des dritten berichteten Mutations-Kandidaten
(„`source_table`-Join in `changes` per `ON false` gekappt") — Test grün
gegen die unveränderte View. (2) Mutation direkt per `psql`
(`CREATE OR REPLACE VIEW cdc.changes … JOIN cdc.source_table st ON
false`) angewendet, Zeilen zwischen den Läufen per `TRUNCATE` bereinigt.
(3) Derselbe Test erneut — **schlägt real fehl**, nicht weil eine
Fixture-Kollision vorlag (die erste Fehlermeldung war ein
Duplicate-Key-Artefakt aus ungereinigten Zeilen des Vorlaufs, korrigiert
per `TRUNCATE` vor dem zweiten Versuch), sondern weil die View bei
`ON false` keine Zeile mehr liefert — die genaue Fehlform, die der
Implementer für diese Mutation berichtet hatte. (4) View zurückgesetzt,
Test erneut grün. Container und Netz danach abgeräumt
(`docker rm -f` / `docker network rm`). **Ergebnis: Die View-Tests sind
echte Wächter, keine Tautologien — mindestens für die dritte der drei
berichteten Mutationen in diesem Lauf am realen Adapter bestätigt.**

Nicht selbst gefahren: `make test-integration` (kein DoD-Item dieses
Slice, die SQL-Views sind nicht Teil des Bootstrap-Verdrahtungspfads) ·
`make doc-commits`/`make doc-immutable` über den vollen Slice-Range
separat — Traceability aller zwölf Range-Commits stattdessen einzeln
gegen den vollen `git log`-Betreff geprüft (Tabelle unten), da das
5er-Gate-Fenster (`HEAD~5..HEAD`) den vollen Range `79c0c12..HEAD` (12
Commits) nicht deckt; alle zwölf Betreffs tragen mindestens eine
`LH-*`-/`ADR-*`-Kennung und keine `SPEC-*`/`ARC-*`-Kennung.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

Vorbemerkung — **Formsache, siehe V-1 unten:** §2 zählt in der Vorlage
eigentlich 10 unterscheidbare Kriterien, trägt aber **11 Zeilen**, weil
„`make gates` grün" doppelt steht (Zeilen 69/70) — ein Planner-Defekt aus
der Welle-3-Eröffnung (`7fc9d23`), nicht aus einem Implementer-Commit. Die
Tabelle unten zählt die 10 unterscheidbaren Kriterien; die Dopplung wird
unter V-1 gesondert behandelt.

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Lesen-Vollabdeckung am `PostgresChangeStoreAdapter` — Teil-Beleg `LH-FA-REA-001`…006 | **bestätigt** | Die fünf zitierten Tests (`TestReadCarriesPositionsAndRanges`, `TestReadCarriesLimit`, `TestReadCarriesTableFilter`, `TestReadIsDeterministicallySorted`, `TestReadLeavesPersistedStateUnchanged`) existieren in `store_test.go`, liefen in diesem Lauf real über `make test-store` (grün, reale PostgreSQL) und decken Bereich/Startposition/Limit/Ordnung/Wiederlesen/Tabellenfilter — Implementer-Behauptung „seit slice-004" durch eigenen Sensor-Lauf bestätigt, nicht nur übernommen |
| 2 | SQL-Views im `cdc`-Schema für SST-002 | **bestätigt** | Drei Views real angelegt (`tools/schema/nacharbeit-views.sql`, per `make schema-rollout` in `make test-store` ausgerollt), drei reale Tests (`sqlviews_test.go`) grün gegen die reale PostgreSQL; eigene Mutationsprobe bestätigt Wächter-Charakter (siehe oben) |
| 3 | `make gates` grün | **bestätigt** | Exit 0 am HEAD `5365ee2` (Sensor-Tabelle oben) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **bestätigt** | `review-slice-010.md` committet (`28e48de`, Fix `c292bca`); F-1 (HIGH) durch Architect-Sequenz `ab00601`/`672a2d7`/`5365ee2` disponiert (siehe F-1-Dispositions-Prüfung unten), F-2 (INFO) an Planner delegiert, kein Closure-Blocker laut Review-Verdikt nach Disposition |
| 5 | Doku-Update für `<Schnittstelle X>` falls öffentlicher Vertrag berührt | **bestätigt, mit Formhinweis (V-1)** | `LH-FA-SST-002`s eigenes Boundary-Kriterium verlangt „die beteiligten Objekte … sind dokumentiert": `tools/schema/schema.yaml:16-27` dokumentiert alle drei Views indikativ (Zweck, Grenzen — u. a. die bewusste `active_tables`-Lücke); `spec/architecture.md` (`5365ee2`) korrigiert zusätzlich den SQL-Kanal der Sicht. Der Plan-Platzhalter `<Schnittstelle X>` selbst ist nie durch einen konkreten Namen ersetzt worden — Formsache, materiell erfüllt |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **noch offen** | §7 trägt am Prüf-HEAD weiterhin die Vorlagen-Platzhalter (`<…>`) — korrekt, da die Verifikation vor der Closure läuft; kein Defekt, sondern der reguläre Vor-Closure-Zustand |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) |
| 8 | Beobachtungs-Register fortgeschrieben | **bestätigt** | `BEO-PGC/d-migrate-nacharbeit/evidence/slice-010.md` neu (`c57215a`); Zähler abgeleitet 2× (`evidence/slice-006.md`, `evidence/slice-010.md`) — unter der 3×-Schwelle, `state.md` korrekt `offen`/`weiter offen` mit Zähler-Vermerk |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **noch offen** | Risiko (a) (SQL-Views konkurrieren mit Use-Case-Lese-Pfaden) trägt am Prüf-HEAD noch „wird bei Closure bewertet" — kein Ausgang aus der geschlossenen Drei-Menge (eingetreten/entfallen/weiter offen) ist bereits eingetragen; `BEO-PGC/lese-doppelquelle` existiert korrekt noch nicht. Regulärer Vor-Closure-Zustand, kein Defekt |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt delegiert** | §2-Zeile und §7-Platzhalter verweisen richtig auf die nächste Welle-3-Closure (Repo **mit** Wellen-Betrieb) — kein Slice-lokaler Prüfpunkt |

**Zwischenstand: 6/10 Kriterien materiell erfüllt und in diesem Lauf
selbst geprüft, 1 korrekt entfallen (Item 7), 1 korrekt an die
Welle-Closure delegiert (Item 10), 2 korrekt noch offen, weil
Verifikation vor Closure läuft (Items 6 und 9).** Kein DoD-*Defekt* im
Sinn eines unbelegten „bestätigt"-Punkts; die beiden offenen Punkte sind
die eigentliche Closure-Arbeit des Planners, kein Implementer-Fehlen.

## F-1-Dispositions-Prüfung (HIGH, review-slice-010.md)

**Frage:** Hält `ADR-0046` tatsächlich, was der Architect-Verdikt
verspricht — enthalten die drei Views wirklich **keine**
Entscheidungslogik?

**Befund je View, gelesen gegen `tools/schema/nacharbeit-views.sql`:**

- `cdc.active_tables` — reiner Join/Scalar-Subquery über
  `cdc.source_table`/`cdc.schema_version`, liefert die höchste
  Schema-Version je Tabelle. Keine `WHERE`-Klausel, kein `CASE`, kein
  Ausschluss nach Zustand — reine Projektion. Die im Schema-Kommentar
  benannte Grenze (Publication-Mitgliedschaft nicht nachgebildet) ist eine
  bewusst **fehlende** Aussage, keine implizite Entscheidung.
- `cdc.consumer_status` — `LEFT JOIN` plus Scalar-Subquery (höchste
  `commit_position`), keine Filterung, kein Schwellenwert-Vergleich; der
  Rückstand (`latest - acknowledged`) wird **nicht** in der View berechnet,
  sondern erst vom lesenden Client (bestätigt: die View liefert nur die
  beiden Roh-Positionsspalten, `TestConsumerStatusViewCarriesBacklog`
  berechnet die Differenz selbst in der Test-Assertion, nicht die View).
- `cdc.changes` — zwei `JOIN`s über Fremdschlüssel, `ORDER BY`. Kein
  `WHERE` in der View selbst; Bereich/Limit/Filter kommen laut Plan und
  Testkommentar vom aufrufenden `WHERE`/`LIMIT` des SQL-Clients — durch
  die eigene Mutationsprobe oben bestätigt: Der Join ist strukturell
  reine Projektion (eine Kappung per `ON false` bricht den Test, was nur
  Sinn ergibt, wenn der Join tatsächlich Daten durchreicht und keine
  Entscheidungslogik enthält, die das Ergebnis unabhängig vom Join
  bestimmen würde).

**Kein `INSERT`/`UPDATE`/`DELETE`, keine Funktionsdefinition mit
Seiteneffekt, kein `CASE`/`WHERE`, das eine Autorisierungs- oder
Domänenentscheidung kodiert** — alle drei Views erfüllen die
Fitness-Function-Beschreibung aus `ADR-0046` (`SQL-Objekte … enthalten
kein INSERT/UPDATE/DELETE, keine Funktionsdefinition mit Seiteneffekt und
keine WHERE-Klausel, die eine Autorisierungs- oder Domänenentscheidung
kodiert`) wörtlich. **F-1 ist damit ein geschlossener Fund, kein neuer.**

**Architektur-Korrektur geprüft (`5365ee2`):** `spec/architecture.md`
trennt den SQL-Kanal jetzt vom CLI-Kanal (eigenes Sequenzdiagramm, SQL
liest direkt gegen die gespeicherten Tabellen). Die Sicht selbst trägt
**keinen ADR-Bezug** — Hard Rule 3.4 gehalten (der Diff `5365ee2` zeigt
keine `ADR-*`-Zeichenkette innerhalb `spec/architecture.md`; die Herkunft
trägt `ADR-0046` allein in seinem `Schärft:`-Feld). Die verwendeten
ARC-Kennungen (`ARC-005` für die SQL-View/den SQL-Kanal, `ARC-006` für
die gespeicherten Tabellen) sind korrekt aus §1 der Architecture-Datei
abgeleitet (`ARC-005` = Driving Adapters, `ARC-006` = Driven Adapters,
Zeilen 58/59 in `spec/architecture.md`) und stimmen mit der bereits
bestehenden Zeile 151 überein (die den SQL-Kanal textlich schon vorher als
„`ARC-005`" führte).

**`ADR-0046` selbst geprüft gegen MADR-Form:** drei Optionen mit
Pro/Contra (A wörtliche Durchsetzung / B Streichen / C gewählt),
Re-Evaluierungs-Trigger benannt (technische Brücke FDW/`dblink`/Extension
— sonst permanent), Folgepflicht (Architektur-Korrektur) benannt und in
`5365ee2` eingelöst, Geschichte-Tabelle trägt den Anlass. `ADR-0018` ist
korrekt `Superseded` mit `Supersedes`-Verweis und unverändertem
historischen Inhalt (Geschichte-Zeile ergänzt, kein Overwrite der
`Entscheidung`-Sektion) — Hard Rule 3.5 gehalten.

## Plan-vs-Code-Diff (Range `79c0c12..HEAD`, beide Nachzüge)

**Deckung §3 gegen den vollen Datei-Diff:**

```
git diff --stat 79c0c12..HEAD -- . ':!docs/plan/planning' ':!docs/reviews' ':!docs/plan/adr'
```

liefert `tools/schema/schema.yaml`, `tools/schema/nacharbeit-views.sql`,
`Makefile`, `internal/adapters/driven/postgresstorage/sqlviews_test.go`,
`spec/architecture.md`, `harness/image-hash.txt`,
`tools/schema/{plan.yaml,down.sql}`.

Die vier §3-Zeilen (nach dem Plan-Nachzug `9a7b3c1`) decken:
`tools/schema/schema.yaml` (update), `tools/schema/nacharbeit-views.sql`
(neu), `Makefile` (update, zweiter psql-Schritt), `sqlviews_test.go`
(neu). Die fünfte, entfallene Zeile (`queries/*.go`) korrekt ohne Diff.

**`spec/architecture.md` steht nicht in §3** — Folgepflicht aus
`ADR-0046`s Konsequenzen-Abschnitt, nachgezogen im selben Commit
(`5365ee2`), der auch den Bezug-Nachzug (`ADR-0046` im Kopf) und den
zweiten Plan-Nachzug (§3-Zeile für `schema.yaml` um den
`ADR-0046`-Verweis ergänzt) trägt. Das ist zulässig und korrekt benannt
(Commit-Message nennt „slice-010 Bezug- und Plan-Nachzug"), aber §3
selbst wurde **nicht** um eine eigene Zeile für `spec/architecture.md`
erweitert — ein Grenzfall: Die Architektur-Korrektur ist Folgepflicht
einer ADR, kein eigener Liefer-Punkt dieses Slice (vergleichbar mit
`harness/image-hash.txt`/`tools/schema/{plan.yaml,down.sql}`, die
verify-slice-009 bereits als „an der §3-Zeile hängendes Pflicht-Nebenprodukt,
kein eigener Liefer-Punkt" eingeordnet hat). Zählt nicht gegen die
≤3-Liefer-Punkte-Grenze (weiterhin 2 Liefer-Punkte: Lesen-Vollabdeckung
war bereits erledigt, SQL-Views ist der einzige tatsächliche
Liefer-Punkt dieses Slice).

**Kein Produkt-Code außerhalb §3/ADR-Folgepflicht.** `harness/image-hash.txt`
war laut Sitzungsstart-Git-Status bereits vor diesem Slice modifiziert
(Lauf-Beleg, `ADR-0044`, nicht Gegenstand dieses Diffs) —
`tools/schema/{plan.yaml,down.sql}` sind Pflicht-Nebenprodukt des
`schema-rollout`-Sensors.

## Review-Dispositions-Check (F-1, F-2)

| Finding | Disposition | Verdikt |
|---|---|---|
| F-1 (HIGH — SQL-Views widersprechen `ADR-0018`) | Architect-Sequenz: `ADR-0046` Accepted (`ab00601`, Fix `672a2d7`), `Supersedes ADR-0018`; Architektur-Korrektur `5365ee2` | **getragen** — siehe F-1-Dispositions-Prüfung oben: die drei Views tragen tatsächlich keine Entscheidungslogik, `ADR-0046` trägt Form (drei Optionen, Trigger, Folgepflicht eingelöst), Hard Rule 3.4/3.5 gehalten |
| F-2 (INFO — §8 vorgelagerte Prüfungen Platzhalter, 2. Auftreten) | keine Fix — Review-Verdikt weist F-2 an den Planner, außerhalb des Implementer-Diffs, kein Closure-Blocker | **korrekt unadressiert im Implementer-Range** — §8 trägt am Prüf-HEAD weiterhin beide Platzhalter; bleibt für die Closure offen (siehe V-2 unten) |

## Befunde

### V-1 — Dup-DoD und Vorlagen-Platzhalter im aktiven Plan (4. bzw. 3. Auftreten derselben Klassen, review-übergreifend nicht erkannt)

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form:
  Slice · review-slice-005 F-11 (Dup-DoD, 1. Auftreten) ·
  review-slice-008 F-8 (Dup-DoD 2., Platzhalter 1.) · review-slice-009
  F-6 (Dup-DoD 3., Platzhalter 2.)
- `pfad`: `docs/plan/planning/in-progress/slice-010-lese-vollabdeckung-sql.md`
  §2 Zeilen 69/70 („`make gates` grün" doppelt) · §2 Zeile 74
  („Doku-Update für `<Schnittstelle X>`" — Vorlagen-Platzhalter, nie durch
  einen konkreten Namen ersetzt)
- `befund`: Dieselben zwei Plan-Form-Defekte wie in slice-005/008/009
  stehen im aktiven Plan von slice-010 — Dup-DoD erreicht damit sein
  **viertes**, der Platzhalter sein **drittes** Auftreten. Beide Defekte
  stammen aus der Welle-3-Eröffnung (`7fc9d23`, Planner-Commit,
  gleichzeitig mit slice-009 angelegt) und wurden vom Plan-Nachzug
  `9a7b3c1` — der ausschließlich §3 berührte — nicht mitgenommen.
  `review-slice-010.md` hat diesen Defekt **nicht** gefunden: Sein
  Gegenstand war explizit auf die drei Implementer-Commits begrenzt und
  schloss den Planner-Commit `7fc9d23` aus (so wie bei F-2 für §8), aber
  anders als bei F-2 wurde für §2 gar kein Hinweis-Finding gesetzt — das
  ist eine echte Lücke im Review-Lauf, keine bewusste Abgrenzung wie bei
  F-2. **Der Platzhalter erreicht mit diesem Auftreten die 3×-Schwelle**
  (review-slice-009 F-6 selbst benannte das Muster „Wiederholung, das
  schon zweimal LOW war" für Dup-DoD bei dessen 3. Auftreten) — nach der
  Register-Logik (Modul 6, §Das Beobachtungs-Register) wäre das der Punkt,
  an dem aus „Notiz" eine „Lücke" wird, die einen eigenen Ausgang braucht.
  Es existiert aber **keine** `BEO-PGC`-Beobachtung für diese Klasse — die
  vier Fundstellen (slice-005/008/009/010) wurden bislang ausschließlich
  über Review-Report-Prosa gezählt, nie ins Register eingetragen. Das ist
  selbst ein Symptom derselben Lücke, die das Register schließen soll
  (Zähler wird "geführt", nicht "abgeleitet").
- `verifizierbar`: ja — Lese der Plan-Datei gegen die Vorlage
  (`.harness/baseline/v6.5.0/templates/docs/plan/planning/slice.template.md`)
  und gegen die vier genannten Review-Reports
- `klasse`: Dup-DoD (4. Auftreten) · Vorlagen-Platzhalter im aktiven Plan
  (3. Auftreten — Schwelle erreicht, kein Register-Eintrag vorhanden)
- **Für die Closure:** (1) beide Form-Defekte in slice-010 §2 vor dem
  `git mv` nach `done/` bereinigen (Dopplung streichen, Platzhalter durch
  einen konkreten Bezug ersetzen — z. B. auf `tools/schema/schema.yaml`s
  Kommentar-Block als erfüllenden Beleg verweisen); (2) für die Klasse
  selbst gehört jetzt ein `BEO-PGC/<slug>`-Eintrag angelegt (mit den vier
  Vorgängen als Belege, sofern sie noch als abgeschlossene Vorgänge
  auflösen) statt einer fünften Review-Prosa-Wiederholung beim nächsten
  Slice — das ist Planner-Entscheidung, kein Verifier-Fix.

### V-2 — §8 vorgelagerte Prüfungen weiterhin Vorlagen-Platzhalter (F-2 bestätigt offen)

- `kategorie`: INFO
- `quelle`: review-slice-010 F-2 (2. Auftreten) · Baseline-Regelwerk
  `modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung
- `pfad`: Slice-Plan §8, beide `<…>`-Zeilen
- `befund`: Bestätigt unverändert seit F-2 — kein Fix-Commit im
  geprüften Range. Kein Closure-Blocker laut Review-Verdikt, aber vor
  dem `git mv` nach `done/` zu füllen (die Sichtung wurde laut Review
  faktisch getan — `BEO-PGC/d-migrate-nacharbeit` wird in §4 zitiert —,
  ist aber in §8 nirgends eingetragen)
- `verifizierbar`: ja — Lese von §8 gegen die Vorlage
- `klasse`: Vorgelagerte §8-Prüfungen unausgefüllt (bestätigt 2. Auftreten,
  keine neue Zählung durch diesen Report)

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS §3.1)** — alle Go-/
  d-migrate-/psql-Belege liefen im gepinnten Container bzw. isolierten
  Testnetz, kein lokales Toolchain-Install (auch die eigene Mutationsprobe
  lief vollständig über `docker run`/`docker exec`)
- geprüft, ohne Befund: **Suppression-Verbot** — keine
  `nolint`/`noqa`/`SuppressMessage`-Marker in den drei Implementer-Commits
- geprüft, ohne Befund: **Kommentar-Klassen (`AGENTS.md` §3.7)** — die
  Kommentare in `nacharbeit-views.sql`, `schema.yaml`, `Makefile`,
  `sqlviews_test.go` sind indikativ über den geltenden Zustand, keine
  Konjunktiv-Klausel über eine verworfene Alternative
- geprüft, ohne Befund: **`ADR-0043`-Ausweichform, 2. Beleg** —
  `BEO-PGC/d-migrate-nacharbeit/evidence/slice-010.md` beschreibt exakt
  denselben `raw-sql-text-drift`-Befund, in diesem Lauf am Rollout-Log
  von `make test-store` bestätigt (`Post-execute compare detected drift`
  träte ohne die Nacharbeit-Datei auf; die Nacharbeit selbst lief sichtbar
  als zweiter psql-Schritt, drei `CREATE VIEW`)
- geprüft, ohne Befund: **Traceability aller zwölf Range-Commits**
  (`38c3c99`, `679b94b`, `4b355e4`, `9a7b3c1`, `fe1ff63`, `c534fec`,
  `c57215a`, `28e48de`, `c292bca`, `ab00601`, `672a2d7`, `5365ee2`) —
  jeder Betreff trägt mindestens eine `LH-*`-/`ADR-*`-Kennung, keine
  Struktur-ID im Betreff (`commit-traceability` grün im
  `make gates`-Lauf für das 5er-Fenster; die übrigen sieben per Lese
  bestätigt)
- geprüft, ohne Befund: **WIP-Limit und Lifecycle** — genau ein Slice
  (`slice-010`) in `in-progress/` am Prüfzeitpunkt
- geprüft, ohne Befund: **a-check-Sichtbarkeitsgrenze** — `.a-check.yml`
  kennt nur `**/*.go`-Globs; die drei SQL-Views liegen außerhalb jedes
  Layer-Globs und sind für a-check strukturell unsichtbar, wie
  `ADR-0046`s Fitness-Function-Zeile selbst feststellt — bestätigt durch
  eigene Lese von `.a-check.yml`
- geprüft, ohne Befund: **Register-Sichtung außerhalb dieses Slices** —
  `BEO-PGC/walsender-wirksamkeit`, `BEO-PGC/adapter-fehler-ausgang`
  betreffen Sub-Areas außerhalb dieses Slice-Umfangs;
  `BEO-PGC/a-check-null-abdeckung` bereits verkörpert; `BEO-PGC/plan-nachzug`
  (2×, unter Schwelle) von diesem Slice nicht berührt — der Plan-Nachzug
  `9a7b3c1` lag vor den Code-Commits im selben Lauf
- geprüft, ohne Befund: **Arbeitsbaum** — außer dem generierten
  Rollout-Report `tools/schema/plan.yaml` (Byprodukt der
  Sensor-Läufe dieses und vorheriger Läufe, kein Plan-/Code-Inhalt) keine
  unerwarteten Änderungen; die einzige Schreibaktion dieses Laufs ist
  dieser Report

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Dup-DoD (V-1, 4. Auftreten) ·
Vorlagen-Platzhalter im aktiven Plan (V-1, 3. Auftreten — Schwelle
erreicht, kein Register-Eintrag) · Vorgelagerte §8-Prüfungen unausgefüllt
(V-2, bestätigt 2. Auftreten, kein neuer Zähler-Schritt).

**Zusammenfassung DoD:** **6/10 Kriterien materiell erfüllt und in
diesem Lauf selbst geprüft** (`make gates` grün Exit 0, `make test` grün
netzlos, `make schema-validate` grün, `make test-store` grün am realen
PostgreSQL inkl. eigener Mutationsprobe), 1 Item korrekt entfallen
(Reconciliation-Register), 1 Item korrekt an die Welle-3-Closure
delegiert (drei Paarungen), 2 Items regulär noch offen, weil
Verifikation vor Closure läuft (Closure-Notiz/Lerneintrag,
Risiko-Ausgang). **Kein DoD-Defekt.**

**F-1 (HIGH) ist geschlossen:** Die drei Views enthalten nachweislich
keine Entscheidungslogik (reine Projektion/Join, kein Schreibzugriff,
keine Business-Regel), `ADR-0046` trägt die MADR-Form vollständig, die
Architektur-Korrektur hält Hard Rule 3.4, `ADR-0018` ist korrekt
`Superseded` unter Hard Rule 3.5. F-2 (INFO) bleibt korrekt offen für den
Planner.

**Neuer Verifier-only-Fund (V-1):** Zwei Plan-Form-Defekte, review-über-
greifend das vierte bzw. dritte Auftreten derselben Klassen, in keinem
der vier betroffenen Review-Läufe für slice-010 selbst erkannt (der
Platzhalter-Defekt erreicht hier die 3×-Schwelle, ohne dass ein
Register-Eintrag existiert) — genau die Art Befund, für die die
Verifier-Rolle vom Reviewer getrennt ist.

## Verdikt

**Merge-/Closure-Konformität:** DoD materiell erfüllt, kein HIGH-Finding
offen; F-1 sauber über die Architect-Sequenz disponiert und in diesem
Lauf eigenständig nachgeprüft (Views ohne Entscheidungslogik, ADR-Form
korrekt, Hard Rules 3.4/3.5 gehalten). Die Lesen-Vollabdeckung und die
SQL-Views sind real und sensor-belegt, nicht nur behauptet — inklusive
einer eigenen, real am Adapter ausgeführten Mutationsprobe.

**Blockierend für Closure (`git mv` nach `done/`):**

1. **V-1:** Dup-DoD- und Platzhalter-Zeilen in §2 bereinigen, vor dem
   `git mv` — dieselbe Reihenfolge wie bei slice-009s V-1 (Inhalt vor
   Move).
2. **§6-Risiko-Ausgang** setzen (eingetreten / entfallen / weiter offen
   → `BEO-PGC/lese-doppelquelle`), regulärer Closure-Schritt.
3. **§7 Closure-Notiz** mit Steering-Loop-Lerneintrag füllen — dabei
   V-1 als Kandidat für einen neuen `BEO-PGC`-Eintrag erwägen (die
   3×-Schwelle des Platzhalter-Defekts ist mit diesem Slice erreicht).
4. **Paarungen (Item 10 der Vorlage, delegiert):** Anker-, Folge-Slice-
   und Register-Paarung prüft die Welle-3-Closure — auch für diesen
   Slice.
5. §8-Platzhalter (V-2) füllen — kein Closure-Blocker laut Review, aber
   offene Formsache.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; alle temporären
Testcontainer/-netze dieses Laufs wurden wieder entfernt.
