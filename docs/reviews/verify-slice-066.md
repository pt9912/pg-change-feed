# Verifikationsbericht: slice-066 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-066` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5
Closure-Trigger, §6 Risiken, §8 Sub-Area) und die bindende `ADR-0059`
(Teilfragen 1–5, §Konsequenzen, §Re-Evaluierungs-Trigger) — **nicht** gegen
den Diff als solchen (Reviewer-Aufgabe; `docs/reviews/review-slice-066.md`
vollständig gelesen, inkl. Fixrunden-Vermerk, aber nur als Kontext, nicht als
Ersatz für eigene Prüfung) und **nicht** gegen realen Bedarf (Validator — hier
nicht ausgelöst, `slice-066` ist kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8),
die vollständige `ADR-0059`, den vollständigen Review-Report inkl.
Fixrunden-Vermerk, den tatsächlichen Diff seit `529f021` (reiner
`next→in-progress`-Move, 0 Zeilen) bis `HEAD = f658728`, und — als
Gegenstände — `internal/domain/model/administrationrequest.go` +
`administrationrequest_test.go`, `internal/bootstrap/{wiring.go,
administration_endtoend_test.go, administration_internal_test.go}`,
`internal/adapters/driven/postgresstorage/{administrationrequest.go,
administrationrequest_test.go, tableactivation.go, queries/queries.go}`,
die beiden neuen Use-Case-Pakete samt Tests, `internal/application/port/
{inbound/verwaltung.go, outbound/{administrationrequest,columnexclusion}.go}`,
`internal/domain/errors/errors.go`, `tools/schema/{schema.yaml,
nacharbeit-administration.sql}`, `spec/architecture.md` und
`spec/pflichtenheft.md`, `spec/lastenheft.md` (`LH-FA-CFG-005`). Die drei
Sensoren wurden in dieser Sitzung **eigenständig real ausgeführt** (Exit-Code
je in einem eigenen, ungepipten Schritt, `AGENTS.md` §3.9) — kein
Implementer- oder Reviewer-Beleg ungeprüft übernommen.

**Gegenstand:** `docs/plan/planning/in-progress/
slice-066-spaltenausschluss-sql-funktionen.md` zum Stand `HEAD = f658728`.
Neun Commits seit `529f021`, die den Slice tragen (`82ce83e` Schema,
`0d2030f` Verarbeitung, `131fd98` Spec, `2246f63` DoD-Häkchen, `0b0ae8e`
Review-Report, `58cddac` Fixrunde, `e11e544` Fixrunden-Vermerk, `f658728`
Planner-Nachzug) — dazu **im Range, aber nicht Teil dieses Slice:**
`8e346e2` (`slice-074` neu angelegt, fremder Slice).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `cdc.administration_request` erweitert (Spalte für den Spaltennamen + `request_kind`-CHECK um `exclude_column`/`include_column`), real ausgerollt über `make schema-rollout` | **erfüllt, selbst reproduziert** | Eigener Wegwerf-Container (gepinnter `postgres:18-alpine`-Digest wie im Makefile), `cdc`-Schema + `search_path`, dann `make schema-rollout` gegen ihn: **Exit 0**. Danach real gelesen: `column_name` liegt in `information_schema.columns` an Position 6; `pg_get_constraintdef` für `chk_administration_request_kind` = `CHECK ((request_kind = ANY (ARRAY['enable','disable','exclude_column','include_column'])))`; ein Einfügeversuch mit `request_kind='truncate'` scheitert real am Constraint (§4 unten). Die Ausweichform-Entscheidung (`schema.yaml` deklarativ für die Spalte, CHECK in `nacharbeit-administration.sql`) deckt sich mit §3 Plan-Nachzug. |
| 2 | `cdc.exclude_column`/`cdc.include_column` als SQL-Funktionen + neuer Outbound Port `ColumnExists` + Inbound Ports `ExcludeColumnUseCase`/`IncludeColumnUseCase` + `ErrSourceColumnMissing` + `applyAdministrationRequest` um zwei `case`-Zweige, belegt über `administration_endtoend_test.go` und `administrationrequest_test.go` | **erfüllt, selbst reproduziert** | Alle Artefakte im Diff real gelesen (`nacharbeit-administration.sql:99–135` vier Funktionen, `:137–138` REVOKE/GRANT; `outbound/columnexclusion.go`; `inbound/verwaltung.go` Sentinel + zwei Use-Case-Interfaces; `tableactivation.go:71–93` `ColumnExists` auf demselben `a.pool`; `wiring.go:1069–1082` die zwei neuen `case`-Zweige ohne `Assembler`-Bindung, wie §1 es fordert). `make test-store` grün (**Exit 0**), u. a. `internal/bootstrap` und `.../postgresstorage` — die beiden genannten Belegdateien laufen real gegen PostgreSQL (§2 unten). |
| 3 | `make gates` grün | **erfüllt, selbst reproduziert** | Eigener, ungefilterter Lauf, Ausgabe in Log-Datei umgeleitet, Exit-Code danach in eigenem Schritt geprüft: **0** (§2 unten). |
| 4 | Review durchgeführt, Report liegt vor | **erfüllt** | `docs/reviews/review-slice-066.md` vollständig gelesen: Ausgang 1 HIGH/2 MEDIUM/2 LOW/3 INFO, Fixrunde `58cddac` real geprüft, Stand danach 0 HIGH/0 MEDIUM/2 LOW (F-9/F-10 als Planner-Nachzug) — Fixrunden-Vermerk `e11e544`. Die DoD-Review-Zeile wurde in `e11e544` real von `[ ]` auf `[x]` gezogen (per `git show e11e544` bestätigt), deckungsgleich mit dem Vermerk. |
| 5 | Doku-Update: `spec/architecture.md` (`ARC-005`-Sequenzsicht um die zwei neuen Antragsarten) + neuer `SPEC-*`-Eintrag in `spec/pflichtenheft.md` | **erfüllt, selbst reproduziert** | `git diff 529f021..HEAD -- spec/architecture.md spec/pflichtenheft.md` real gelesen: die Sicht trägt eine Vierzeilen-Tabelle (Antragsart · SQL-Funktion · Inbound Port) ohne `ADR-*`-/Slice-Bezug (`AGENTS.md` §3.4), der Satz über das Diagramm ist auf „am Beispiel `enable`" präzisiert (die frühere Über-Behauptung F-6 ist weg); `SPEC-019` beschreibt die Feldform deckungsgleich mit Modell und Nacharbeit und nennt keinen `ADR-*`-Rückverweis. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<bei Closure>`-Platzhalter (Volltext gelesen) — Planner-Arbeit nach diesem Bericht. |
| 7 | Reconciliation-Register — entfällt (Greenfield) | **korrekt offen, Entfall-Vermerk trägt** | Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`); `docs/plan/planning/reconciliation.md` real geprüft: existiert nicht. Der Entfall-Vermerk in der DoD-Zeile ist zutreffend; das Item bleibt bis zur Closure unangetastet (zulässig). |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/evidence/` real gelistet: `slice-006/010/015/036` = 5 Belege, **kein** `slice-066`-Beleg. §3 Plan-Nachzug kündigt einen neuen Beleg an (neue Objektklasse „CHECK-Änderung an bestehender Tabelle") — dieser gehört in die Closure, nicht in den Diff. Kein widersprüchlicher Zustand. |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Einträge real gelesen: beide tragen noch wörtlich `<bei Closure zuzuweisen>`. |
| 10 | Die drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/` (real bestätigt), `Welle: welle-18`-Feld vorhanden; DoD-Zeile verweist korrekt auf die `welle-18`-Closure. |

**Ergebnis §1:** Alle fünf implementierungs-/reviewbezogenen DoD-Punkte (1–5)
sind real erfüllt; Punkt 1, 2, 3 und 5 wurden **selbst reproduziert**, nicht
nur behauptet. Die fünf verbleibenden Closure-Punkte (6–10) sind korrekt noch
offen und wurden nicht vorweggenommen.

## 2. Sensor-Läufe (alle drei selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | d-check 519 Dateien / 0 Befunde; d-check `--enable commits --range HEAD~5..HEAD` 519/0; `commit-traceability` OK (5 Commits, Betreffe ohne Struktur-ID); a-check 0 Befunde (2 unverschichtete Werkzeug-Dateien, bekannte Hinweise); coverage-gate OK — 44.90 % ≥ Schwelle 35 % |
| `make test` | **0** | vollständige Suite im Race-Container, alle Pakete `ok`, inkl. `internal/application/usecase/{exclude,include}column` und `internal/domain/model` |
| `make test-store` | **0** | Rollout + alle realen PostgreSQL-Tests grün, u. a. `internal/bootstrap` und `internal/adapters/driven/postgresstorage` |

`git status --porcelain` nach allen Läufen: leer — kein unbeabsichtigter
Seiteneffekt auf den Arbeitsbaum.

**Nicht als Gate ausgeführt (existiert im Repo nicht):** `make doc-commits`,
`make doc-immutable`. Die vom Rollen-Briefing genannten Ziele existieren in
diesem Repo nicht (`AGENTS.md` §4 listet nur real existierende Targets). Die
Traceability je Commit trägt `make commit-traceability` (Bestandteil von
`make gates`, oben: 5 Commits, 0 Befunde); MR-/ADR-Immutabilität ist über
`ADR-0059` als Hard Rule real geprüft (§7 unten).

## 3. Mutations-Stichprobe — DoD-Zusage real rot gesehen (Verifier-only-Nachweis)

Der Implementer hat je Zusage rot gesehene Gegenbeispiele dokumentiert. Zwei
davon wurden **eigenständig nachgestellt** — Mutation gesetzt, Rotlauf gesehen,
Mutation zurückgenommen, Arbeitsbaum danach gegen den Commit geprüft.

**Mutation A — Leer-Spalten-Prüfung entfernt** (`administrationrequest.go`,
der `if column == ""`-Rumpf der beiden Spalten-Antragsarten):

```
MUTATION_TEST_EXIT=1
--- FAIL: TestNewAdministrationRequestRejectsInvariantViolations (0.00s)
    --- FAIL: .../Spalten-Antragsart_ohne_Spalte (0.00s)
    administrationrequest_test.go:73: Art "exclude_column" ohne Spalte: Fehler = <nil>, wollen ErrEmptyIdentifier
FAIL  github.com/pt9912/pg-change-feed/internal/domain/model
```

**Mutation B — `default`-Zweig der geschlossenen Menge geleert** (unbekannte
Antragsarten würden akzeptiert):

```
MUTATION2_TEST_EXIT=1
    --- FAIL: .../Antragsart_außerhalb_der_geschlossenen_Menge (0.00s)
FAIL  github.com/pt9912/pg-change-feed/internal/domain/model
```

Nach jeder Mutation `cp`-Restore aus der Vorlage, danach:
`git status --porcelain` **leer** und `git diff HEAD` **leer** — der
Arbeitsbaum ist identisch zum Commit. Beide Zusagen (Spalten-Antragsart ohne
Spalte → `ErrEmptyIdentifier`; unbekannte Art → `ErrInvalidAdministration-
RequestKind`) sind damit nicht nur behauptet, sondern **rot-fähig** belegt.

## 4. Ausweichform `tools/schema/nacharbeit-administration.sql` — Idempotenz und Nicht-Deklaration, selbst geprüft

Auf einem eigenen Wegwerf-Container (gepinntes `postgres:18-alpine`), nach
`make schema-rollout` (Exit 0) als Ausgangszustand:

- **Lauf 2** (`psql -f nacharbeit-administration.sql`): **Exit 0**.
- **Lauf 3**: **Exit 0**.
- `pg_get_constraintdef` vor und nach den beiden Läufen **identisch**:
  `CHECK ((request_kind = ANY (ARRAY['enable','disable','exclude_column','include_column'])))`.
- Die vier Funktionen `enable_table`/`disable_table`/`exclude_column`/
  `include_column` liegen unverändert in `pg_proc`.
- **Constraint ist wirksam, nicht nur vorhanden:** ein Einfügeversuch mit
  `request_kind='truncate'` endet real mit
  `ERROR: new row for relation "administration_request" violates check
  constraint "chk_administration_request_kind"` (Exit 1).

**`tools/schema/schema.yaml` deklariert die CHECK-Klausel nicht mehr** — real
gelesen: der `constraints:`-Block von `administration_request` trägt nur noch
`chk_administration_request_status`; ein Kommentar am Ort nennt die
Ausweichform und die d-migrate-Grenze. Bestätigt.

Container und Docker-Netz danach abgeräumt; die beim Rollout regenerierten
`tools/schema/{plan.yaml,down.sql}` wurden aus der Sicherung wiederhergestellt
(der **einzige** Unterschied der Regeneration war das umgebungsabhängige
`target`-Feld der DSN im Report — kein Inhaltsdefekt; `down.sql` byte-identisch).

## 5. F-9 / F-10 — Planner-Nachträge (`f658728`), selbst geprüft

- **F-9 (Herkunft der Messung auflösbar) — trifft.** Der Beleg-Absatz in §3
  nennt jetzt die Herkunft: „real gemessen von der Reviewer-Rolle gegen
  PostgreSQL 18 / d-migrate 1.3.1, dokumentiert in
  [`docs/reviews/review-slice-066.md`](../reviews/review-slice-066.md)
  §Eigene Nachmessung; vom umsetzenden Lauf übernommen, nicht dort erneut
  gefahren". Die Quelle ist damit auflösbar (Rolle + Artefakt + Abschnitt).
- **F-10 (`errors.go`-Zeile beschreibt den realen Änderungsinhalt) — trifft.**
  Die §3-Zeile führt `internal/domain/errors/errors.go` jetzt mit „Kommentar
  der geschlossenen Antragsarten-Menge (`ErrInvalidAdministrationRequestKind`)
  auf die vier Werte nachgezogen — der Sentinel `ErrSourceColumnMissing` liegt
  in `internal/application/port/inbound/verwaltung.go`". Der reale Diff
  bestätigt das exakt: in `errors.go` ändert sich **ausschließlich** der
  Kommentar von `ErrInvalidAdministrationRequestKind` (auf die vierwertige
  Menge); der Sentinel steht in `verwaltung.go`.

## 6. Zwei offen gebliebene Punkte — eigenes Urteil

**F-7 (INFO, bewusst nicht behoben).** Die geschlossene `request_kind`-Menge
liegt nach dieser Änderung an drei Orten (Go-Konstanten in
`administrationrequest.go`, Prosa-Festlegung in `SPEC-019`, handgeschriebene
DDL-Klausel in `nacharbeit-administration.sql`); ein Sensor auf die Paarung
gibt es nicht. **Ich teile die Verwerfung.** Ein Sensor müsste drei
Sprachräume zugleich auswerten (Go, Prosa, DDL-Text) oder die lebende DB
befragen — kein kleines, dependency-freies Skript, dieselbe Verwerfung wie
beim Chronik-Sensor (`docs/reviews/architect-verdict-slice-chronik-in-
code-kommentar.md`). Die mittelbare Deckung trägt real: die beiden
Adapter-/E2E-Tests schreiben Zeilen **aller vier** Arten gegen die lebende
Klausel (`make test-store`, Exit 0, §2). Die benannte Grenze — eine Klausel,
die *weiter* ist als die Go-Menge, bliebe unentdeckt — ist heute harmlos,
weil kein Aufrufer einen Wert außerhalb der Go-Menge erzeugt. INFO ohne
erwartete Aktion: **zutreffend.**

**`make image` nicht gefahren, Digest nicht neu gestempelt.** Der Slice ändert
Go-Build-Kontext; `harness/README.md` §Werkzeuge verlangt für einen solchen Zug
`make image` **vor seiner Closure** (Digest-Commit nur bei geändertem Digest).
`make image` ist ausdrücklich **kein Gate** und steht **nicht** in der
`slice-066`-DoD; der §5-Closure-Trigger fordert DoD + `make gates` +
Adapter-Test + Closure-Notiz — kein `image`. **Für diese Verifikation also kein
Blocker: ein sauber übergebener Rest an die Planner-Closure.** Einschränkung,
benannt statt verschwiegen: der Rest ist bislang in **keinem committeten
Artefakt** festgehalten (§7 trägt nur `<bei Closure>`) — er lebt in der
Implementer-Übergabe. Empfehlung an den Planner: den Schritt in §7 aufnehmen,
damit er bei der Closure nicht aus dem Kontext fällt. Da der Go-Code sich
ändert, ist mit einem geänderten Digest und damit mit einem Digest-Commit zu
rechnen.

## 7. Hard Rules und Plan-vs-Code-Diff

- **3.3 (`git mv` + Inhalt = zwei Commits):** `529f021` ist real der Elter —
  `git show --stat 529f021` zeigt 0 Zeilen Änderung (reiner Move); nicht
  Bestandteil des geprüften Bereichs.
- **3.5 (Accepted-ADR immutable):** `ADR-0059` in dieser Sitzung nicht
  verändert — der Diff berührt `docs/plan/adr/` nicht (`git diff
  529f021..HEAD --stat` bestätigt); der Slice *umsetzt* die ADR nur.
- **3.7 (Kommentar-/Chronik-Disziplin):** der Review-HIGH F-1 (Fund-Referenz
  im Produktionscode-Godoc) ist im Fix-Commit behoben — `administrationrequest.go:51`
  trägt jetzt den Rang-Zeiger `ARC-001` statt `Review-Finding F-4`; die zwei
  weiteren Chronik-Stellen in `administrationrequest.go`/`wiring.go` sind
  ebenfalls in die Indikativ-Form überführt (im Fix-Diff real gelesen). Kein
  Slice-/Wellen-Bezug in den neuen Produktionskommentaren dieses Diff.
- **3.8 (Action-Pinning):** keine `uses:`-Zeile im Diff — nicht berührt
  (`.github/` nicht im Diff).
- **3.9 (Exit-Code nie gepiped):** in dieser Sitzung durchgehend beachtet —
  jede der drei Sensoren in eigene Log-Datei umgeleitet, Exit-Code in eigenem,
  ungekettetem Schritt geprüft; Mutationsläufe ebenso (§2/§3).
- **Plan-vs-Code-Diff / Scope:** `git diff --stat 529f021..HEAD` zeigt 27
  Dateien. Der Slice-Anteil deckt sich mit §3 (nach F-4-Nachzug konkret
  benannt). `internal/adapters/driving/http`, `compose.yaml`,
  `.github/workflows/`, `harness/` und fremde Planungsdokumente sind **nicht**
  berührt. **Einziger Fremdanteil im Range:** `8e346e2` (`slice-074` neu
  angelegt, 443 Zeilen) — ein anderer Slice, nicht Teil von `slice-066`; die
  §3-Tabelle muss ihn nicht führen. Kein Scope-Creep des Slice.
- **Architektur-Kanten:** `make a-check` grün (0 Befunde, §2) — die neuen
  Ports liegen auf den deklarierten Schichten; `ColumnExclusionPort` wird nur
  im Adapter implementiert und in der Application nur als Interface gesehen.

## 8. Beobachtung zur Reproduzierbarkeit des Review-Reports (nicht DoD-relevant)

F-10 des Review-Reports gibt als Verifikationshinweis
`git show 82ce83e -- internal/domain/errors/errors.go` an. Dieser Befehl
liefert **keinen** Treffer: `82ce83e` (Schema/Antrags-Queue) berührt
`errors.go` nicht; die Datei wurde real in `0d2030f` geändert
(`git log 529f021..HEAD -- internal/domain/errors/errors.go`). Der **Nachtrag
selbst ist korrekt** (§5 oben) — nur der Hinweis zeigt auf den falschen
Commit. Der Review-Report ist Lauf-Beleg und wird über Läufe hinweg nicht
gelesen; diese Zeile ist deshalb eine Beobachtung, kein Blocker.

## Negativbefunde

- **geprüft, ohne Befund: Punkt 2 (SQL-Funktionen/Ports/Verarbeitung).** Die
  Funktionen schreiben nachweislich nur einen Antrags-Datensatz (`INSERT` in
  `cdc.administration_request`) und senden `pg_notify` — kein Zugriff auf
  `cdc.source_table` oder die Publication (`ADR-0018`/`ADR-0046`-Disziplin).
  `applyAdministrationRequest` trägt für die beiden Spalten-Arten **keine**
  `Assembler`-Bindung nach und fällt nicht in den `default`-Fehlerzweig; der
  `default`-Text nennt jetzt die vierwertige Menge. Der Zwischenzustand
  „`applied`, aber noch keine Filterwirkung" ist in §1 ausdrücklich als
  Zustand dieses Slice benannt (`slice-067`/`068` schließen ihn).
- **geprüft, ohne Befund: neuer Outbound Port / Pool / Layering.** Ein
  `ColumnExclusionPort` mit genau `ColumnExists`; die Implementierung sitzt am
  `TableActivationAdapter` auf demselben `a.pool` (ein `cdc_admin`-Pool, kein
  zweiter, kein neuer Grant); `validateIdentifier` für Schema/Tabelle, der
  Spaltenname geht als **Wert** in die parameterisierte Katalog-Abfrage
  (`queries.SelectTableColumnExists`), nicht in DDL-Text. Die `REVOKE`/`GRANT`-Zeile
  nimmt beide neuen Signaturen mit.
- **geprüft, ohne Befund: Testform der zwei Belegdateien.** Der neue
  `administrationrequest_test.go`-Fall ruft beide SQL-Funktionen real auf,
  liest `column_name`/`request_kind` zurück und prüft die `cdc_reader`-Ablehnung
  für `cdc.exclude_column` (SQLSTATE 42501); der `applied`/`failed`-Übergang
  wird dort korrekt mit konstruierter Meldung gesetzt (die reale
  Spaltenprüfung liegt in `administration_endtoend_test.go`, das genau den
  `applied`- und `failed`-Pfad gegen den echten `ColumnExclusionPort` fährt).
  Beide unter `make test-store` real grün (§2).
- **geprüft, ohne Befund: Traceability/ID-Schema der Slice-Commits.** Alle
  acht `slice-066`-Commits tragen `LH-FA-CFG-005` und `ADR-0059` im Betreff,
  kein `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability` grün (§2).
- **geprüft, ohne Befund: §7/§8-Konsistenz.** §8 listet den Register-Stand
  korrekt (u. a. `BEO-PGC/d-migrate-nacharbeit` 5×, verkörpert;
  `BEO-PGC/adapter-fehler-ausgang` 2×, weiter offen) und kündigt den neuen
  Beleg an; die Sub-Area-Wahl ist konsistent mit `harness/conventions.md`
  (`*`/`PGC`, GF — Begründungsblock entfällt, Abschnitt bleibt).

## Summary

| Kategorie | Anzahl |
|---|---|
| DoD-Punkte erfüllt (1–5) | 5 |
| Closure-Punkte korrekt offen (6–10) | 5 |
| Verifier-Findings (Blocker) | 0 |
| Beobachtungen ohne Blocker | 2 (F-7-Urteil teilt die Verwerfung; `make image`-Rest sauber übergeben) |
| Fremdanteil im Diff-Range | 1 Commit (`8e346e2`, fremder Slice) |

## Verdikt

**DoD-Konformität: bestätigt.** Alle fünf implementierungs-/reviewbezogenen
DoD-Punkte (1–5) sind real erfüllt; Punkt 1, 2, 3 und 5 wurden **selbst
reproduziert** (eigener Rollout, drei eigene Sensor-Läufe, eigene
Mutationsläufe), nicht übernommen. Die fünf Closure-Punkte (6–10) sind korrekt
noch offen und wurden nicht vorweggenommen.

**Mutations-Stichprobe: bestanden.** Zwei dokumentierte DoD-Gegenbeispiele
eigenständig nachgestellt — beide liefen real rot (Exit 1), beide Restores
hinterließen einen byte-identischen Arbeitsbaum (`git status`/`git diff HEAD`
leer).

**Ausweichform: bestätigt.** `nacharbeit-administration.sql` ist idempotent
(Lauf 2 und 3 Exit 0, unveränderter Endzustand), die Klausel ist wirksam
(fünfter Wert abgelehnt), und `schema.yaml` deklariert sie nicht mehr.

**`ADR-0059`-Konformität: bestätigt.** Die Antrags-Seite ist so gebaut, wie
die ADR sie in Teilfrage 1/5 und §Konsequenzen festlegt; die vollständige
Erfüllung von `LH-FA-CFG-005` (Happy-Path-Wirkung) liegt planmäßig erst bei
`slice-067`/`068` — die Slice-Abgrenzung deckt das ausdrücklich.

**F-9/F-10-Nachträge: beide treffen.**

**Sensor-Exit-Codes:** `make gates` **0** · `make test` **0** ·
`make test-store` **0** · (Mutationsläufe je **1**, erwartetes Rot).

**Übergabe an Planner:** Der Slice kann an die Closure übergeben werden.
Für die Closure-Notiz vorzumerken: (a) §6-Risiko 1 — Ausgang zuweisen; die
Bestands-Instanz-Konstellation endet real mit **Exit 8**
(`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`, in §3 Plan-Nachzug dokumentiert),
der frische Rollout dagegen **Exit 0** (selbst bestätigt, §1/§4); (b)
§6-Risiko 2 — Ausgang zuweisen (`BEO-PGC/adapter-fehler-ausgang` bekommt in
diesem Diff real keinen dritten Fund); (c) Beobachtungs-Register:
`evidence/slice-066.md` zu `BEO-PGC/d-migrate-nacharbeit` (neue Objektklasse
„CHECK-Änderung an bestehender Tabelle", §3 Plan-Nachzug); (d) `make image`
vor der Closure fahren und den Digest-Commit nur bei geändertem Digest setzen
(§6 dieses Berichts); (e) der F-7-INFO-Eintrag bleibt begründet stehen, keine
Aktion.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
