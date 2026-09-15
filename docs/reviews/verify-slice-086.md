# Verifikationsbericht: slice-086 — 2026-09-15

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-086` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) — der
Endpunkt, der Inbound Port, die Filterachse, die Rechtsklasse —,
[`ADR-0057`](../plan/adr/0057-http-grpc-api.md) — die API und ihre
Token-Klassen (in zwei Klauseln teil-abgelöst) —,
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
— der View-Direktzugriff (**unberührt**) —,
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) — Transport-Typen
am Port —, [`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md)
Festlegung 1 — die opake `SourceTableID` ist kein Draht-Bezeichner —,
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) — der Digest-Beleg —;
`AGENTS.md` §3.1, §3.9, §3.10, §3.11). **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe; `review-slice-086.md` wurde als Kontext gelesen, **nicht**
als Beleg übernommen) und **nicht** gegen realen Bedarf (Validator, hier
nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf hat den Slice-Plan in der Fassung von
`HEAD` gelesen, dazu die sechs Entscheidungen, den Review-Report und die
berührten Artefakte. **Alle** Zahlen dieses Berichts stammen aus eigenen, in
dieser Sitzung gefahrenen Läufen; kein Beleg des Implementers, des Reviewers
oder des Planners wurde übernommen. Exit-Codes je **ungepiped** und in
eigenem Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung
getrennt beauftragt. Jeder Schema-/Testlauf verändert die committete
`tools/schema/plan.yaml`; sie wurde nach **jedem** Lauf real zurückgenommen —
der Arbeitsbaum ist am Ende dieses Laufs unverändert (`git status --porcelain`
leer).

**Gegenstand:** `slice-086`, geprüfter Stand `HEAD` = `e962048`
(„slice-086 DoD-Box Review nach Fixrunde gesetzt"), Zweig `main`. Der Slice
liegt weiterhin in `in-progress/`; der `git mv` nach `done/` ist **nicht**
erfolgt. Die Fixrunde umfasst `f20d2c9` (F-2), `42b3785` (F-4) und `e962048`
(F-3/F-5, DoD-Box); der Review-Stand war `88cb28b`.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (ungepiped, Ausgabe in Logdatei, Exit direkt gelesen) | **0** | baseline-verify `v6.5.0` OK (54 Dateien) · d-check **692** Dateien / **0** Befunde · commit-traceability OK (5 Commits, Betreffe ohne Struktur-ID) · a-check **0** Befunde · coverage-gate **OK — 72.00 %** erfüllt Schwelle 70 % |
| 2 | `make test` (netzlos, Race-Detector, sauberer Baum) | **0** | alle Pakete grün, kein `FAIL` |
| 3 | `make test-store` (reale PostgreSQL, Schwelle 70) | **0** | `DB-Adapter-Coverage: 73.38% (477/650; Profile gemergt: store,replication)` · `db-coverage: OK … Schwelle 70%` |
| 4 | `make commit-traceability RANGE=b2d1cc2..HEAD` | **0** | **12** Commits, je Betreff mit `LH-*`/`ADR-*`, Betreffe ohne `SPEC-*`/`ARC-*` |
| 5 | d-check Modul `vcs` (git-Diff-Immutabilität), `RANGE=b2d1cc2..HEAD` | **0** | 0 Befunde — kein Accepted-/Record-Artefakt inhaltlich überschrieben |
| 6 | d-check Modul `immutable`, `RANGE=b2d1cc2..HEAD` | **0** | 0 Befunde |
| 7 | **Gegenprobe Threshold** — `make coverage-gate THRESHOLD=99` (Ist 72.00 %) | **2** | `coverage-gate: FAIL — Coverage 72.00% unter Schwelle 99%` — die Schwelle **prüft** |
| 8 | **Gegenprobe Threshold** — `DB_COVERAGE_THRESHOLD=99 make test-store` (Ist 73.38 %) | **2** | Store-Tests **grün** (`ok … postgresstorage 4.087s`), dann `db-coverage: FAIL — DB-Adapter-Coverage 73.38% unter Schwelle 99%` — die Schwelle **prüft** |
| 9 | **Mutation A** — `queries.SelectChanges` Projektions-Ordnung `st.schema_name, st.table_name` → vertauscht | **2** | `make test-store` rot: `TestReadCarriesTableFilter` (`store_test.go:418`, „Schema:b Table:public, wollen public/b") — die **Projektions-/Klartext-Identität** ist real gebunden |
| 10 | **Mutation B** — `parseReadChangesQuery`: die `source`-Pflicht-Prüfung entfernt | **2** | `make test` rot: `TestReadChangesFehlendeQuelleEndetMit400` — die **Pflicht von `source`** ist real gebunden |
| 11 | `git diff b2d1cc2..HEAD -- tools/schema/` | — | **leer** — die View `cdc.changes` und der View-Direktzugriff sind unberührt |
| 12 | `git log b2d1cc2..HEAD -- .github/workflows/` | — | **leer** — der Slice berührt **keinen** Workflow (`AGENTS.md` §3.10 greift nicht) |
| 13 | Projektions-Vergleich `queries.SelectChanges` ↔ `tools/schema/schema.yaml` (`changes`-View) | — | **13** Spalten in **identischer** Ordnung, dieselben zwei Inner Joins, dieselbe `ORDER BY`; Unterschied allein `WHERE`/`LIMIT` |
| 14 | Zählprobe `docs/user/e2e-abdeckung.md` gegen `test/integration/integration_test.go` | — | alle Verweise um **+4** verschoben; die zwei Einfüge-Hunks vor `TestE2ECaptureFlow` ergeben `190 → 194`, danach netto 0 — arithmetisch konsistent |

**Nicht gefahren (bewusste Grenze):** `make test-integration`, `make image`
und der `e2e`-Workflow — nicht im Werkzeug-Satz dieses Laufs (siehe §7).

---

## 2. DoD-Konformität, Kriterium für Kriterium

Gelesen ist der Text in der Fassung `e962048`.

| # | DoD-Zeile (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| LP1-K1 | Die Filterachse des Leseports ist **eine** Form: `Schema`/`Table` als **Klartext** ersetzt `Table *model.SourceTableID` | **erfüllt** | `outbound.ChangeQuery` trägt `Schema string` und `Table string` (`changestore.go:43–44`); kein `Table *model.SourceTableID` mehr im Pfad (`grep` findet nur den erklärenden Doc-Kommentar). Mutation A bindet die neue Achse am Store. |
| LP1-K2 | `SelectChanges` zieht den `cdc.source_table`-Join nach; `mapper.ToChange` setzt `Schema`/`Table` | **erfüllt** | `queries.SelectChanges` liest `st.schema_name`/`st.table_name` über `JOIN cdc.source_table` (`queries.go:52–63`); `mapper.ToChange` setzt `change.Schema = row.Schema`/`change.Table = row.Table` (`mapper.go:134–135`). Projektions-Vergleich (Messung 13) deckungsgleich mit der View. |
| LP1-K3 | Port-/Mapper-Tests decken die neue Achse (einzeln, kombiniert, leer) | **erfüllt** | `TestChangeQueryValidateAcceptsOptionalTextFilter` (einzeln/kombiniert/leer, Port-Kontrakt), `TestToChangeCarriesSchemaAndTable` (Mapper), `TestReadCarriesTableFilter` (Store, einzeln/kombiniert/leer/ohne Treffer), `translate_test.go` (Projektions-Ordnung). |
| LP2-K1 | `ReadChangesUseCase` als Inbound Port, Transport-Typen **am Port** (`ADR-0042`); der Inbound Port importiert **nicht** outbound | **erfüllt** | `port/inbound/readchanges.go` definiert `ReadChangesQuery`/`ReadChange`/`ReadChangesResult` und importiert nur `context` + `domain/model`; `usecase/readchanges/service.go` führt sie als **Typ-Aliase** (`= inbound.…`). a-check grün (Messung 1). |
| LP2-K2 | Filter: `source` **Pflicht**, `schema`/`table` optional und **unabhängig**, `from` inklusiv / `to` exklusiv, `limit` optional — **kein** Default-Limit | **erfüllt** | `parseReadChangesQuery` verlangt `source` (Mutation B bindet es); SQL `>= $2` (Start) und `< $3` (End) (`queries.go:65–66`); `textArgument` gibt `nil` für leere Filter (unabhängig); `readChangesLimit` setzt **kein** Default, der Use Case normiert nichts (`TestReadChangesLeavesOptionalFiltersUnset`). |
| LP2-K3 | Netzlose Tests am Use Case: Filterwirkung, Bereichs-Grenzen, Leerfall | **erfüllt** | `service_test.go`: `TestReadChangesTranslatesQueryToPort`, `TestReadChangesLeavesOptionalFiltersUnset`, `TestReadChangesCarriesInvertedRange`, `TestReadChangesCarriesNonPositiveLimit`, `TestReadChangesEmptyResultIsSet`, `TestReadChangesRejectsMissingSource` — `make test` grün (Messung 2). |
| LP3-K1 | `GET /changes` in der `reader`-Rechtsklasse, kollisionsfrei neben `GET /changes/stream`; `ErrNonPositiveLimit`/`ErrRangeInverted` → **400**, **unbekannter Parameter → 400**, **kein Treffer → 200 mit `{"changes": []}`**, nie 404 | **erfüllt** | `server.go:102` registriert `mux.Handle("GET /changes", withToken(…, roleReader, …))`; `errors.go` bildet beide Sentinel auf 400 ab; `parseReadChangesQuery` verwirft jeden Namen außerhalb der geschlossenen Menge (`readChangesParams`); `readChangesHandler` initialisiert `make([]readChangeResponse, 0, …)` → `{"changes":[]}`. Fünf Adapter-Tests decken 401/400/200/500 und die Kollisionsfreiheit; `make test` grün. |
| LP3-K2 | Der Vertrag ist fortgeschrieben: `SPEC-018`s Satz nachgezogen, **`SPEC-022`** angelegt, Handbuch nennt den Endpunkt samt Querverweis | **erfüllt** | `SPEC-018` trägt jetzt den Zielwortlaut („Das Changes-Lesen … ist in `SPEC-022` ausgestaltet; Diagnose/Health bleibt außerhalb"); `SPEC-022` steht als eigener Abschnitt in der Form der bestehenden (Fähigkeits-Tabelle + Merkmal-Tabelle); Handbuch `Version: 1.15` mit Fähigkeits-Zeile `GET /changes` und Querverweis aus dem Lese-Abschnitt. **Kein `spec → adr`-Verweis**: `grep 'adr/'` über `spec/pflichtenheft.md` und `spec/lastenheft.md` ist **leer**; d-check Modul `matrix` grün. |
| LP3-K3 | `make gates` grün (Exit direkt, ungepiped) | **erfüllt** | Messung 1: **Exit 0**, alle fünf inneren Gates grün. |
| — | Review durchgeführt, Report unter `docs/reviews/` liegt vor; Fixrunde: 0 HIGH, F-2…F-5 geschlossen, F-1 mit `slice-087` adressiert | **erfüllt** | `review-slice-086.md` liegt vor; F-3 (`readchanges_test.go:378/382`, Name stimmt jetzt), F-4 (`SPEC-022` präzisiert, s. unten), F-5 (nur **eine** Leerzeile vor `---`, geprüft in `pflichtenheft.md:422–424`) sind am Artefakt bestätigt; `slice-087` liegt real in `open/` und **nimmt F-1 an** (§Ziel nennt die `ADR-0081`-Fitness-Function-Zeile). |
| — | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **offen (Planner)** | §7 trägt noch Platzhalter; der Slice ist `in-progress`. |
| — | Reconciliation-Register `../reconciliation.md` | **entfällt (trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht** — die §2-Zeile sieht genau diesen Entfall vor. |
| — | Beobachtungs-Register fortgeschrieben, **kein** Zähler gesetzt | **offen (Planner)** | Der Lauf steht aus; die §8-Sichtung nennt zwei Einträge. Eigene Zählprobe: `zahl-in-traeger-driftet-gegen-die-messung` = **1×** (slice-081), `dod-begruendung-unzutreffende-tatsachenbehauptung` = **3×** (slice-036/081/082) — die §8-Zahlen stimmen. |
| — | Jedes Risiko aus §6 trägt einen Ausgang | **offen (Planner)** | Die drei §6-Einträge tragen noch `<bei Closure>`. |
| — | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen, hier nicht zuständig** | Das Repo **hat** Wellen (`welle-20` offen); die Paarungen trägt die nächste Wellen-Closure (Modul 6 Schritt 3c) — die §2-Zeile sagt das. |

**Ergebnis §2:** Alle **drei Liefer-Punkte tragen vollständig** (LP1-K1…K3,
LP2-K1…K3, LP3-K1…K3 bestätigt). Die nicht abgehakten Posten sind die
Closure-Pflichten des **Planners** (§7, Register, Risiko-Ausgänge), kein
Code-/Vertrags-Defekt; der Slice ist bewusst noch `in-progress`.

---

## 3. Kein zweiter Lesepfad — die tragende Zusage

**Bestätigt, am Code.** Es gibt **genau eine** Kette:

`http` → `inbound.ReadChangesUseCase` → `readchanges.ReadChangesService` →
`outbound.ChangeStorePort.ReadChanges` → `PostgresChangeStoreAdapter.
ReadChanges` → `queries.SelectChanges`.

Kein zweites SQL, keine eigene Sortierung, keine Filterentscheidung im
Adapter, kein `make(…, 0, …)`-Default. Die Sortierung
(`ORDER BY t.commit_position, c.transaction_id, c.sequence`) ist unverändert
und mit der View zeichengleich in Ordnung und Joins (Messung 13). Die
Projektions-/Klartext-Identität ist durch **Mutation A** real gebunden
(Projektions-Ordnung vertauscht → Store-Test rot).

**Die View `cdc.changes` ist unberührt:** `git diff b2d1cc2..HEAD --
tools/schema/` ist **leer** (Messung 11). Der Join lebt ausschließlich im
Lesepfad (`SelectChanges`), nicht in der View; kein DDL, kein
Schema-Eingriff. `ADR-0046` bleibt unberührt.

---

## 4. Entscheidungs-Konformität

| ADR | Prüfung | Verdikt |
|---|---|---|
| [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) Teilfrage 1 | Neuer Inbound Port, Transport-Typen am Port, Use Case als **Übersetzung** (bildet auf `ChangeQuery` ab und delegiert), kein neuer Outbound-Port, keine zweite Sortierung | **konform** |
| `ADR-0081` Teilfrage 2 | Endpunkt `GET /changes` (nicht-streamende zweite Form), Rechtsklasse `reader`, geschlossene Parameter-Menge, **kein** Default-Limit, Antwortfeldnamen im API-Vokabular (`schema`/`table`, `old_image`/`new_image`, `source_table_id` daneben), leere Menge `{"changes": []}` nie `null` | **konform** (Feldnamen und Form am `readChangeResponse` bestätigt) |
| `ADR-0081` Teilfrage 3 | Klartext-Achse **eine** Form, `cdc.source_table`-Join, Mapper-Feld, `ADR-0046` unberührt | **konform** |
| `ADR-0081` Teilfrage 4 | Pflichtfeld fehlt → 400; unbekannter Parameter → 400 **mit Namen**; unlesbare Zahl / `< 1` / `from > to` / `limit < 1` → 400; kein Treffer → 200; 401; Store-Fehler → 500 — **kein 404-Pfad** | **konform** |
| `ADR-0081` **§Fitness Function**, Zeile `make test-integration` | verlangt einen **realen** `GET /changes` gegen den laufenden Feed-Container über den HTTP-Wegwerf-Client | **offen — mit benanntem Ausgang** | Kein Aufrufer (`grep '"/changes' tools/ test/` findet nur `/changes/stream`); der Slice-Plan §1 nimmt die Wegwerf-Clients aus und §5 setzt „(Adapter-Test)" an die Stelle des Netzwerk-Belegs. Der Ausgang ist real benannt: `slice-087` (`open/`) trägt genau diese Zeile. Dies ist **Review-F-1**; ich habe ihn als Kontext gelesen, **nicht** nachgeprüft (ausgelagert laut Auftrag). |
| [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) | Die zwei Token-Klassen und ihre Hierarchie; `GET /changes` in der niedrigeren, lesenden Klasse (`admin` deckt sie implizit ab); kein 403-Pfad auf dem Endpunkt | **konform** — die zwei in `ADR-0081` abgelösten Klauseln (Changes-Lesen-Hälfte) sind sauber ersetzt, alles Übrige (Protokoll, Adapter-Platzierung, Diagnose-Hälfte) unberührt |
| [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) | View-Direktzugriff des SQL-Kanals bleibt; keine Entscheidungslogik in SQL, keine zweite Leselogik in Go — nur eine zusätzliche Spalte im bestehenden Join | **konform** (View unverändert, Messung 11) |
| [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) | Transport-Typen am Port, Use Case aliasiert statt neu zu definieren; keine neue Schichten-Kante | **konform** — a-check 0 Befunde; `.a-check.yml` unverändert (kein Diff) |
| [`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md) Festlegung 1 | Die opake `SourceTableID` ist **kein** Draht-Bezeichner; die Filterachse ist Klartext | **konform** — `SourceTableID` bleibt nur als zusätzliches Antwortelement (kein Adressierungsmittel) |
| [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) | Build-Kontext-Dateien geändert → `make image` vor der Closure, eigener Digest-Commit; Digest selbst ist lauf-gebunden, nicht inhalts-vergleichbar | **form-konform** — `c051b7f` ist ein eigener `chore(image)`-Commit (`LH-FA-SST-006`/`ADR-0081`), der `harness/image-hash.txt` neu stempelt; der Digest-**Wert** ist nicht nachrechenbar (§(2) der ADR) |

**Einordnung zu F-1:** Der Slice liefert genau den Schnitt, den `ADR-0081`
§Folgepflicht (Slice-Schnitt-Empfehlung) als **drei** Liefer-Punkte führt.
Die `make test-integration`-Zeile steht nur in §Fitness Function/§Testabdeckung
der ADR — dort als „Erwartung, keine abschließende Festlegung" geführt — und
wird von einem **eigenen** Vorgang (`slice-087`) getragen. Das ist eine
bewusst getrennte Lieferung mit benanntem Ausgang, **keine** stille Lücke;
der Slice-Text (§1, §5) benennt die Trennung ausdrücklich.

---

## 5. Plan-vs-Code-Diff

**Zuschnitt (§3):** Der Slice nennt als Träger den Leseport samt Adapter und
Mapper, den Inbound Port/Use Case, den HTTP-Adapter und die Vertrags-Dateien
(`spec/pflichtenheft.md`, `docs/user/benutzerhandbuch.md`). Der Diff über
`b2d1cc2..HEAD` berührt **genau** diese Bereiche.

**Der §3-Nachtrag des Implementer-Laufs (vier Erweiterungen) ist gerechtfertigt:**

1. `internal/bootstrap/wiring.go` — **gerechtfertigt.** Der lesende Endpunkt
   braucht einen `ChangeStorePort`; nur `cdc_admin` trägt
   `SELECT, DELETE ON cdc.transaction, cdc.change`
   (`tools/schema/nacharbeit-roles.sql:67`), `cdc_reader` ausschließlich die
   vier Views (Z. 105). Die Verdrahtung wählt `cfg.AdminDSN` und folgt dem
   Muster von `retentionStore`. Die Rolle bleibt von der API-Rechtsklasse
   unabhängig; `reader` erhält keine Fähigkeit, die nicht ohnehin lesend im
   Port liegt.
2. `test/integration/integration_test.go` — **erzwungen.** Die Typänderung
   `Table *model.SourceTableID` → `Table string` erzwingt den Nachzug
   (Kompilier-Bindung); das neue `table`-Feld trägt den Klartext-Namen.
3. `docs/user/e2e-abdeckung.md` — **gerechtfertigt (Erzeugnis).** Die
   Zeilennummern folgen dem verschobenen Ort ihres Quelltexts; Zählprobe
   (Messung 14) bestätigt die **+4**-Verschiebung arithmetisch (Generator
   nicht neu gefahren, siehe §7).
4. `harness/image-hash.txt` — **gerechtfertigt** ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)). Eigener
   Digest-Commit, Build-Kontext geändert.

**Die „Nicht in dieser Liste"-Zeile trägt:** `git diff b2d1cc2..HEAD` enthält
**keine** Datei aus `tools/schema/**`, **keine** `.a-check.yml`-Änderung,
**keine** `examples/**`- und **keine** `Makefile`/`harness/mk/**`-Änderung
(Messung 11 und `git diff --name-only`).

**Kein stiller Umfangszuwachs:** Die berührten Doku-Dateien tragen den
geänderten Zustand (keine neue Zusage); kein Produktionscode-Kommentar nennt
`slice-`/`welle-` (`AGENTS.md` §3.7); kein host-lokaler Pfad (§3.11, d-check
`hostpaths` grün); keine `//nolint`-Suppression; der Diff enthält nur `A`/`M`
(§3.3).

---

## 6. `AGENTS.md` §3.10

**§3.10 greift nicht.** `git log b2d1cc2..HEAD -- .github/workflows/` ist
**leer** (Messung 12); der Diff enthält keine Workflow-Datei. Der Slice führt
**keinen** neuen Workflow und **keine** strukturelle Änderung an einem
bestehenden ein — die Hard Rule ist damit nicht berührt. (Die
strukturellen Grenzen des Werkzeugs für die bestehenden Workflows bleiben
gleichwohl bestehen, siehe §7.)

---

## 7. Was ich nicht prüfen konnte

- **`make test-integration`**, **`make image`** und der **`e2e`-Workflow**
  (Compose-Stack bzw. gehosteter Runner): nicht im Werkzeug-Satz dieses Laufs.
  Damit ist der **Netzwerk-Beleg** zu `LH-FA-SST-006` aus `ADR-0081`
  §Fitness Function **nicht** durch einen eigenen Lauf bestätigt — genau die
  Lücke, die Review-F-1 als ausgelagert führt (§4).
- **Der Post-Push-Lauf** auf GitHub: kein eigener Lauf, kein `gh`-Zugriff
  möglich; `AGENTS.md` §3.10 bleibt eine Regel, deren Beleg ich nur
  **einsehen**, nicht **erzeugen** kann — hier ohnehin nicht berührt (§6).
- **Der Digest-Wert** in `harness/image-hash.txt`: nach [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) §(2)
  builder- und lauf-gebunden und **nicht** inhaltlich nachrechenbar — geprüft
  ist allein die **Form** des Commits (eigener `chore(image)`-Commit).
- **Der E2E-Abdeckungs-Generator**: nicht neu gefahren; die Tabelle ist
  über Messung 14 **arithmetisch** konsistent, nicht über einen Generator-Neulauf.
- **Die Vollständigkeit des Belegs „kein zweiter Lesepfad"** über die
  gesamte Testhistorie: geprüft ist der Diff dieses Slice und die
  Bindungs-Kraft von zwei eigenen Mutationen, nicht jeder frühere Stand.

---

## 8. Befunde dieses Laufs

Es gibt **keinen** neuen Verifier-Fund, der eine Änderung verlangt. Der eine
offene Posten ist bereits benannt und adressiert:

### V-1 — `ADR-0081` §Fitness Function (`make test-integration`) ist mit diesem Slice noch nicht eingelöst (Adressat: Planner/Architect — **Status, kein neuer Defekt**)

- `pfad`: `docs/plan/adr/0081-changes-lesen-ueber-die-http-api.md`
  §Fitness Function (Zeile `make test-integration`) und §Testabdeckung
  (letzter Aufzählungspunkt) gegen `docs/plan/planning/in-progress/slice-086-changes-lesen-api.md` §1/§5
- `befund`: Die ADR verlangt einen **realen** `GET /changes` gegen den
  laufenden Feed-Container; kein Aufrufer im Baum ruft den nicht-streamenden
  Endpunkt an (`grep '"/changes' tools/ test/` → nur `/changes/stream`). Der
  Slice-Plan ersetzt das durch einen **Adapter-Test** und nimmt die
  Wegwerf-Clients aus. Der Ausgang ist **benannt und angenommen**: `slice-087`
  liegt in `open/` und nennt in §Ziel genau diese Fitness-Function-Zeile.
- `kategorie`: Fitness-Function-Zusage ohne Träger im aktuellen Stand
- `verifizierbar`: ja — `grep -rn '"/changes' tools/ test/`
- **Rollen-Grenze:** Dieser Punkt ist Review-F-1 und laut Auftrag
  **ausgelagert**; ich habe ihn **nicht** nachgeprüft, sondern nur seinen
  Entscheidungs-Konformitäts-**Stand** und seinen **Ausgang** festgehalten.
  Der Verifier repariert nicht.

**Negativbefunde (geprüft, ohne Befund):**

- Der `Source == ""`-Guard steht im Use Case **und** im Port-Kontrakt
  (`TestReadChangesRejectsMissingSource`) — als Eingangs-Guard begründbar,
  kein Drift (Review-F-7).
- `ADR-0081` §Kontext („zehn Routen") und `SPEC-018` („neun Port-gedeckte
  Fähigkeiten") sind **snapshot-korrekt**: die ADR beschreibt den
  gemessenen Vorzustand, `SPEC-018` scoped weiterhin die neun, `SPEC-022`
  trägt die zehnte — kein Widerspruch.
- Die Antwortform trennt sauber: `old_image`/`new_image` als eingebettetes
  JSON, fehlendes Bild `null` (`rowImage`), `committed_at` als
  `RFC3339Nano` in UTC — deckungsgleich mit `SPEC-022` (das den
  Bruchteil-Nullen-Fall ausdrücklich benennt, Review-F-4).

---

## 9. Verdikt

**Der Slice ist closure-fähig** — die drei Liefer-Punkte und ihr Vertrag
tragen vollständig; es fehlen allein die **Planner**-Closure-Pflichten
(§7-Notiz, Register, Risiko-Ausgänge), die naturgemäß erst beim Übergang nach
`done/` entstehen.

Getragen ist alles, was der Slice zusagt:

- **Kein zweiter Lesepfad:** eine Kette, derselbe Port, dieselbe Delegation,
  dieselbe Sortierung, dieselbe Projektion wie die View; `tools/schema/**`
  **unberührt** (Messungen 11/13) — der Kern von `ADR-0046`/`ADR-0081`.
- **Der Endpunkt verhält sich wie zugesagt:** `reader`, kollisionsfrei,
  `source` Pflicht (eigene Mutation B rot), unbekannter Parameter → 400,
  Bereichs-/Limit-Grenzen → 400, kein Treffer → 200 `{"changes": []}`, nie
  404, **kein Default-Limit**; Feld- und Bildform wie `SPEC-022`.
- **Der Vertrag ist fortgeschrieben**, ohne `spec → adr`-Verweis; Handbuch
  nachgezogen.
- **`make gates` grün** (Exit 0); die Schwellen **prüfen** noch (eigene
  Gegenproben: `coverage-gate THRESHOLD=99` rot, `DB_COVERAGE_THRESHOLD=99`
  rot); `make test`/`make test-store` grün.
- **Die vier §3-Erweiterungen** sind je gerechtfertigt; der Digest-Commit ist
  form-korrekt nach [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md).
- **§3.10 greift nicht** — kein Workflow berührt.
- Die **eigenen Mutationen** bestätigen zwei Zusagen, die kein anderer Lauf
  mutiert hat: die Projektions-/Klartext-Ordnung (Store) und die Pflicht von
  `source` (HTTP).

**Was fehlt / offen bleibt:** (1) die `ADR-0081`-Fitness-Function-Zeile
`make test-integration` — getragen von `slice-087` (V-1, Review-F-1, nicht
neu); (2) die Planner-Closure-Posten. **Nicht closure-blockierend** ist der
Adapter-Negativtest zu `limit < 1` (F-2): die gelieferte Semantik ist
korrekt, und die Regressions-Verteidigung ist in `f20d2c9` real gebunden
(der Fake-Store prüft jetzt die übergebene Abfrage; Mutation der
Implementer-Hälfte `raw=="0" → nil` würde `Limit == nil` verletzen).

---

**Übergabe:** Dieser Bericht geht an den **Planner** (Closure-Entscheidung;
V-1 ist sein bereits benannter, an `slice-087` adressierter Posten) und als
Entlastung an den **Implementer** (der Diff ist nachweislich tragend); an den
**Validator** nur, falls dieser Slice als MVP-Slice validiert wird — der reale
Bedarf liegt hier repo-extern nicht vor. Der Bericht ist ein **Lauf-Beleg**
(dieser Stand, diese Läufe, dieses Verdikt) und ersetzt weder das Review noch
die Planner-Closure.
