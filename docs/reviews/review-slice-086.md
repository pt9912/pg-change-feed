# Review-Report: slice-086 — 2026-09-15

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `slice-086`, Diff `b2d1cc2..b198b9c` — sieben Commits:
`1690ce8` (Leseport in Klartext), `91795c3` (Use Case + Inbound Port),
`2a51cc8` (Endpunkt + Verdrahtung), `c8d24c4` (Spec/Handbuch), `ba36d28`
(Projektions-Ordnung), `c051b7f` (Digest-Beleg), `b198b9c` (E2E-Abdeckungs-
tabelle). 23 Dateien, 1214 Einfügungen / 90 Löschungen. Kein
`tools/schema/**`-Eingriff, keine `.a-check.yml`-Änderung, kein `Makefile`/
`harness/mk/**`.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-086-changes-lesen-api.md`
  vollständig (§1–§8)
- [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  (der Gegenstand), [`ADR-0057`](../plan/adr/0057-http-grpc-api.md)
  (Token-Klassen, teil-abgelöst), [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (View-Direktzugriff — unberührt), [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)
  (Transport-Typen am Port), [`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md)
  Festlegung 1 (opake Kennung als Draht-Bezeichner), [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
  (Digest-Beleg), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md)
  (Rollen-DSNs), [`ADR-0030`](../plan/adr/0030-testpyramide.md) (Test-Tiers)
- Berührte `LH-*`: [`LH-FA-SST-006`](../../spec/lastenheft.md) (Hauptbezug),
  [`LH-FA-REA-001`](../../spec/lastenheft.md)…`006`,
  [`LH-FA-CAP-008`](../../spec/lastenheft.md) (Boundary der Row Images),
  [`LH-QA-POR-003`](../../spec/lastenheft.md) (E2E-Tabelle)
- `AGENTS.md` §3.1, §3.3, §3.7, §3.9, §3.11, §5 · `harness/conventions.md`
  (MR-000 ID-Schema), `spec/architecture.md` §3/§4, `spec/lastenheft.md`,
  `spec/pflichtenheft.md` §2/§4/§7
- Vorherige Läufe am gleichen Modul: `docs/reviews/review-slice-081.md`,
  `review-slice-082.md` (zuletzt `slice-chronik-in-code-kommentar` und
  `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`)

---

## Eigene Messungen dieses Laufs

Docker-only; jeder Exit-Code **ungepiped** und in einem eigenen Schritt
gelesen (`AGENTS.md` §3.9). Das Lauf-Artefakt `tools/schema/plan.yaml`
(entsteht bei `make test-store`) wurde nach jedem Lauf per `git checkout --`
zurückgenommen; `git status --porcelain` ist vor diesem Commit leer.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` | **0** | alle fünf inneren Gates grün; `coverage-gate: OK — 71.90% erfüllt Schwelle 70%`; `a-check: gesamt 0 Befund(e)`; `d-check: 690 Dateien, 0 Befunde`; `commit-traceability: OK` |
| 2 | `make test` | **0** | Race-Detector-Lauf durchweg grün, inkl. `usecase/readchanges` und `adapters/driving/http` |
| 3 | `make test-store` | **0** | reale PostgreSQL; `DB-Adapter-Coverage: 73.38% (477/650) … erfuellt Schwelle 70%` |
| 4 | **Mutation A** — `ORDER BY t.commit_position` → `DESC` (`queries.go`) | **2** | `TestReadIsDeterministicallySorted`, `TestReadCarriesLimit`, `TestDeleteChangesRemovesOnlyGivenChanges` rot — die Ordnungszusage ist real gebunden |
| 5 | **Mutation B** — Bereichsende `commit_position < $3` → `<= $3` | **2** | `TestReadCarriesPositionsAndRanges` rot (`to` exklusiv ist real gebunden) |
| 6 | **Mutation C** — Tabellenfilter verlangt Schema (`$5`-Prädikat an `$4 IS NOT NULL` gekoppelt) | **2** | `TestReadCarriesTableFilter` rot (`store_test.go:383`, „Filter table=a = [], wollen nur seine 2 Changes") — die Unabhängigkeit der beiden Filter ist real gebunden |
| 7 | **Mutation D** — `readChangesLimit`: `raw == "0"` → `(nil, nil)` (statt `&0`) | **0** | `make test` bleibt **grün**, obwohl `GET /changes?limit=0` danach still unbegrenzt liest → **F-2** |
| 8 | `grep -rn '"/changes' tools/ test/` | — | nur `GET /changes/stream` (`tools/harness/sseclient/main.go:44,121`); **kein** Aufruf des nicht-streamenden Endpunkts → **F-1** |
| 9 | Projektions-/Ordnungs-Gegenprobe `queries.SelectChanges` ↔ `tools/schema/schema.yaml` (`changes`-View) | — | 13 Spalten in **identischer** Ordnung (`source_id, commit_position, change_id, transaction_id, source_table_id, schema_name, table_name, sequence, operation, old_data, new_data, schema_version, committed_at`), dieselben **zwei Inner Joins**, dieselbe `ORDER BY`; Unterschied allein `WHERE`/`LIMIT` — `ba36d28`s Behauptung bestätigt |
| 10 | Zählprobe `docs/user/e2e-abdeckung.md` gegen `test/integration/integration_test.go` | — | **alle 13** Go-Zeilen-Verweise um genau **+4** verschoben; die drei Einfüge-Hunks (`+3`, `+1`, `0`) ergeben am Bezugspunkt `TestE2ECaptureFlow` genau `190 → 194` — die Tabelle ist mit dem Lauf arithmetisch konsistent, nicht von Hand verstellt |
| 11 | Rollen-Nachweis der Verdrahtung gegen `tools/schema/nacharbeit-roles.sql` | — | `cdc_admin` trägt `GRANT SELECT, DELETE ON cdc.transaction, cdc.change` (Z. 67); `cdc_reader` trägt ausschließlich die vier Views (Z. 105) — die Wahl von `cfg.AdminDSN` ist die einzige, die `ChangeStorePort` überhaupt tragen kann |

**Nicht gemessen (bewusste Grenze):** `make test-integration` wurde **nicht**
ausgeführt (Compose-Stack, nicht im Werkzeug-Satz dieses Laufs). Die
E2E-Abdeckungstabelle ist daher über Messung 10 arithmetisch, nicht über
einen Generator-Neulauf geprüft. `harness/image-hash.txt` ist ein
builder-/lauf-gebundener Digest (`ADR-0044`) und ist **nicht**
inhaltlich nachrechenbar — geprüft ist allein die Form (Negativbefund).

---

## Findings

### F-1 — Die ADR-0081-Fitness-Function-Zeile `make test-integration` hat keinen Träger

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  §Fitness Function (Zeile `make test-integration`) und §Testabdeckung
  (letzter Aufzählungspunkt) — beide verlangen einen **realen**
  `GET /changes` gegen den laufenden Feed-Container über den bestehenden
  HTTP-Wegwerf-Client als Netzwerk-Beleg zu `LH-FA-SST-006`
- `pfad`: tools/harness/run-integration-tests.sh:1951 („HTTP-API-Rundlauf"),
  tools/harness/httpclient/main.go:20
- `befund`: Im Baum ruft kein Aufrufer den nicht-streamenden Endpunkt an —
  `grep -rn '"/changes' tools/ test/` findet ausschließlich
  `GET /changes/stream`. Der Slice-Plan §1 nimmt die Wegwerf-Clients aus und
  §5 setzt „real über die API erreichbar **(Adapter-Test)**" an die Stelle des
  Netzwerk-Belegs; der Verweis `slice-083` trägt diesen Punkt nicht (jener
  Slice ist `examples/nats-client`, nicht der E2E-Runner). Die ADR widerspricht
  sich dabei selbst: ihre §Folgepflicht „Slice-Schnitt-Empfehlung" führt genau
  die drei Liefer-Punkte, die der Slice liefert — die E2E-Erweiterung steht
  nur in §Testabdeckung/§Fitness Function.
- `verifizierbar`: ja — `grep -rn '"/changes' tools/ test/`; ein grüner
  `make test-integration` ohne diesen Aufruf bestätigt die Lücke
- `klasse`: Fitness-Function-Zusage ohne Träger

### F-2 — Der Adapter-Negativtest für `limit < 1` bindet den 400-Status nicht an die Eingabe

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  §Testabdeckung (Whitebox-Unit-Tests des Adapters:
  „`400` für … `limit < 1`") · `LH-FA-REA-003` Negative
- `pfad`: internal/adapters/driving/http/readchanges_test.go:258
  (`TestReadChangesNichtpositivesLimitEndetMit400`), :25/:31
  (`fakeReadChangesUseCase`), internal/adapters/driving/http/readchanges.go:190
  (`readChangesLimit`)
- `befund`: Der Test ruft `?limit=0` gegen einen Fake, der
  `outbound.ErrNonPositiveLimit` **unabhängig von der übergebenen Abfrage**
  zurückgibt — der 400-Status hängt damit an keinem Parameterwert und ist für
  `?limit=10` derselbe. Eigene Mutation D: `readChangesLimit` liefert für `"0"`
  `nil` statt `&0` (der Aufruf liest dann still unbegrenzt, statt `400` zu
  enden) — `make test` bleibt **Exit 0**. Die gelieferte Implementierung ist
  korrekt; die Verteidigungslinie gegen ihre Regression trägt nicht.
- `verifizierbar`: ja — `make test` (Mutation D dieses Laufs, Exit 0)
- `klasse`: Negativtest ohne Bindung an die Eingabe

### F-3 — Doc-Kommentar über einem Test nennt eine nicht existierende Funktion

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (ein Kommentar beschreibt, was
  da ist)
- `pfad`: internal/adapters/driving/http/readchanges_test.go:355
- `befund`: Der Kommentar beginnt mit
  `TestReadChangesHandlerOhneUseCaseTraegtDerFehlerpfad`; die darunter
  deklarierte Funktion heißt
  `TestReadChangesOhneVerdrahtungEndetVorDemHandler`. Die Beschreibung des
  geprüften Verhaltens selbst ist zutreffend — es ist ein Namensrelikt einer
  Umbenennung, kein inhaltlich falscher Satz (deshalb LOW und nicht die
  HIGH-Klasse „Kommentar trägt keine der Kommentar-Klassen").
- `verifizierbar`: nein — kein Gate liest Doc-Kommentare
- `klasse`: Kommentar nennt abwesende Kennung

### F-4 — `SPEC-022` sagt „mit Nanosekunden", der Code schreibt `time.RFC3339Nano`

- `kategorie`: LOW
- `quelle`: [`SPEC-022`](../../spec/pflichtenheft.md) (Antwort-Zelle) ·
  `LH-FA-REA-004`
- `pfad`: spec/pflichtenheft.md:407 ·
  internal/adapters/driving/http/readchanges.go:91
- `befund`: `time.RFC3339Nano` entfernt Nullen am Bruchteil-Ende; ein
  Commit-Zeitpunkt ohne Bruchteil-Sekunden erscheint als
  `…T10:30:00Z`, nicht als `…T10:30:00.000000000Z`. Der Spec-Wortlaut
  „`committed_at` (`<RFC 3339 mit Nanosekunden, UTC>`)" trifft damit nicht für
  jeden Wert wörtlich.
- `verifizierbar`: nein — kein Gate; die Aussage folgt aus der Semantik der
  `time`-Konstante `RFC3339Nano`
- `klasse`: Spec-Wortlaut überzeichnet die Implementierung

### F-5 — Zwei Leerzeilen-Reste vor dem Abschnittstrenner nach `SPEC-022`

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: spec/pflichtenheft.md:423–426
- `befund`: Nach der letzten Tabellenzeile von `SPEC-022` stehen drei
  Leerzeilen vor `---` (dort, wo eine genügte). Ohne Wirkung auf den Vertrag.
- `verifizierbar`: nein
- `klasse`: Formatierungsrest

### F-6 — `ARC-005` der Sicht trägt weiter „später HTTP-/gRPC" (außerhalb des Diff)

- `kategorie`: INFO
- `quelle`: [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  §Schärft (`ARC-005`)
- `pfad`: spec/architecture.md:58
- `befund`: Die Sicht-Zeile `ARC-005 Driving Adapters` endet auf „… später
  HTTP-/gRPC", während HTTP (jetzt **zehn** Routen) und gRPC live sind. Der
  Fund liegt **außerhalb** dieses Diff und ist vorgefunden, nicht von diesem
  Zug erzeugt; er steht hier, weil `ADR-0081` ihre Bindung an `ARC-005`
  erklärt.
- `verifizierbar`: nein
- `klasse`: vorgefundene Sicht-Drift

### F-7 — Die `Source == ""`-Prüfung steht zweimal (Use Case und Port-Kontrakt)

- `kategorie`: INFO
- `quelle`: Maintainability · [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  §Festlegungen Teilfrage 1
  („Der Use Case ist eine Übersetzung, keine Logik")
- `pfad`: internal/application/port/inbound/readchanges.go (Port-Doku) ·
  internal/application/usecase/readchanges/service.go:51
- `befund`: Der Use Case weist eine leere Quellen-Kennung selbst ab, obwohl
  `outbound.ChangeQuery.Validate()` dieselbe Bedingung über
  `domainerrors.ErrEmptyIdentifier` trägt. Dieselbe Entscheidung an zwei
  Stellen, gegen denselben Sentinel — als Eingangs-Guard vor dem Port-Aufruf
  begründbar und durch `TestReadChangesRejectsMissingSource` abgedeckt; kein
  Handlungsbedarf, nur benannt.
- `verifizierbar`: ja — `make test` (`TestReadChangesRejectsMissingSource`)
- `klasse`: benannte Doppelung ohne Drift

---

## Negativbefunde

- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/**` — der
  Lesepfad trägt **dieselbe** Delegation, **dieselbe** Sortierung
  (`ORDER BY t.commit_position, c.transaction_id, c.sequence`, unverändert)
  und **dieselbe** Projektion wie `cdc.changes` (Messung 9). Der
  Tabellenfilter ist **eine** Form (`st.schema_name`/`st.table_name`,
  parametergebunden, kein String-Bau); die beiden Inner Joins sind wegen der
  FK `cdc.change.source_table_id → cdc.source_table.source_table_id`
  (`tools/schema/schema.yaml:145`) erhaltend, der Retention-Lesepfad verliert
  dadurch keine Zeile. **Kein zweiter Lesepfad.**
- geprüft, ohne Befund: `tools/schema/**` — **unberührt**. Der Diff enthält
  keine Datei dieses Verzeichnisses; die View `cdc.changes` behält ihren
  direkten Projektionszugriff (`ADR-0046` unberührt), der Join lebt
  ausschließlich im Lesepfad (`SELECT_CHANGES`). Kein DDL, kein Schema-Eingriff.
- geprüft, ohne Befund: `internal/adapters/driving/http/**` — `GET /changes`
  in der `reader`-Rechtsklasse (`roleReader`, admin deckt sie implizit ab,
  kein 403-Pfad); geschlossene Parameter-Menge mit `400` samt Parameternamen;
  `from`/`to` `< 1` als `400` **direkt und eingabegebunden** (im Unterschied
  zu `limit < 1`, siehe F-2); kein Treffer → `200` mit `{"changes": []}`
  (nicht `null`, nie `404`); Store-Fehler → `500` geloggt. Der Handler berührt
  kein Speichermedium und spricht ausschließlich über den Inbound Port.
  Kollisionsfrei neben `GET /changes/stream` (exakte ServeMux-Muster, durch
  `TestReadChangesUndStreamKollidierenNicht` belegt).
- geprüft, ohne Befund: `internal/application/**` — `port/inbound/readchanges.go`
  importiert nur `context` + `domain/model`, **kein** outbound; die
  Transport-Typen liegen am Port und der Use Case führt sie als Aliase
  (`ADR-0042`); der Use Case bildet ab und delegiert, setzt **kein**
  Default-Limit und keine Obergrenze.
- geprüft, ohne Befund: `internal/domain/model/change.go` — der Kommentar
  beschreibt jetzt den Ist-Zustand (zwei füllende Pfade), keine
  Vorher/Nachher-Sprache, kein Slice-/Wellen-Bezug. Kein Produktionscode-
  Kommentar des Diff nennt `slice-`/`welle-` (§3.7 geprüft).
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` — die Verdrahtung ist
  gerechtfertigt: `postgresstorage.New(ctx, cfg.AdminDSN, …)` ist die
  **einzige** Rolle mit `SELECT` auf `cdc.change`/`cdc.transaction`
  (Messung 11), und sie folgt dem bestehenden Muster von `retentionStore`
  (wiring.go:491). Die Rolle bleibt von der API-Rechtsklasse unabhängig; der
  `reader`-Token erhält keine Fähigkeit, die `ChangeStorePort` nicht ohnehin
  lesend trägt (kein `DeleteChanges`/`PersistTransaction` im Pfad).
- geprüft, ohne Befund: `test/integration/integration_test.go` und
  `docs/user/e2e-abdeckung.md` — die Typänderung des Filters
  (`*model.SourceTableID` → `string`) **erzwingt** den Nachzug (Kompilier-
  Bindung); die Abdeckungstabelle ist arithmetisch konsistent regeneriert
  (Messung 10). `harness/image-hash.txt` ist form-korrekt nach `ADR-0044`
  (eigener `chore(image)`-Commit, Build-Kontext geändert, neuer Digest).
- geprüft, ohne Befund: `spec/pflichtenheft.md`, `docs/user/benutzerhandbuch.md`
  — `SPEC-018`s Satz trägt den Zielwortlaut aus `ADR-0081` („ist in `SPEC-022`
  ausgestaltet; Diagnose/Health bleibt außerhalb"); `SPEC-022` steht in der
  Form der neun bestehenden Abschnitte (Kennung, Titel, Fließtext, Merkmals-
  Tabellen) und zitiert **keine** ADR und **keinen** Slice (Referenzrichtung
  `spec → adr` gewahrt, `matrix` grün); die Historie hat ihre Zeile. Das
  Handbuch ist auf `Version: 1.15` gehoben, mit neuer Änderungszeile, neuer
  Fähigkeits-Zeile `GET /changes` und Querverweis aus dem Lese-Abschnitt
  (die HIGH-Klassen „Handbuch-Versionshistorie nicht fortgeschrieben" und
  „neue Betreiber-Oberfläche ohne Handbuch-Zug" sind **nicht** erfüllt).
- geprüft, ohne Befund: Hygiene — kein Lauf-Artefakt im Diff (`tools/schema/
  plan.yaml` bleibt unberührt), keine Vorlagen-Reste, keine unaufgelösten
  Kennungen, kein host-lokaler Pfad (§3.11), keine `//nolint`-Suppression,
  keine Rename-/Move-Vermischung (§3.3 — der Diff enthält ausschließlich `A`
  und `M`), alle sieben Betreffs tragen `LH-*`/`ADR-*` und **keine**
  `SPEC-*`/`ARC-*`-Kennung, kein Trailer.

---

## Antworten auf die Schwerpunkte

**1 — Kein zweiter Lesepfad (tragende Zusage): bestätigt, am Code.**
`http` → `inbound.ReadChangesUseCase` → `readchanges.ReadChangesService` →
`outbound.ChangeStorePort.ReadChanges` → derselbe `SelectChanges`. Kein
zweites SQL, keine eigene Sortierung, keine Filterentscheidung im Adapter,
kein `make(…,0,…)`-Default. Die Sortierung ist unverändert und durch Mutation A
real gebunden; die Projektion ist mit der View zeichengleich in Ordnung und
Joins (Messung 9). **`cdc.changes` selbst ist unberührt** — `tools/schema/**`
kommt im Diff nicht vor; der Join lebt im Lesepfad.

**2 — Verhalten des Endpunkts: erfüllt.** `source` Pflicht; `schema`/`table`
optional und unabhängig (Mutation C gebunden); `from` inklusiv / `to` exklusiv
(Mutation B gebunden); `limit` optional, **kein** Default, **keine** harte
Obergrenze; unbekannter Parameter → `400` mit Namen; kein Treffer → `200`
mit `{"changes": []}`, nie `404`; Rechtsklasse `reader`; kollisionsfrei neben
`GET /changes/stream`. Ausnahme: die **Verteidigungslinie** für `limit < 1`
trägt nicht (F-2), die gelieferte Semantik ist korrekt.

**3 — Der Vertrag: erfüllt.** `SPEC-022` in der Form der bestehenden acht
Abschnitte; `SPEC-018`s Satz trägt den Zielwortlaut aus `ADR-0081`;
Handbuch nennt den Endpunkt samt Querverweis; **die Herkunft bleibt
auffindbar, ohne dass der Spec-Text eine ADR zitiert** — die ADR trägt ihr
`Schärft:`, der Spec-Text ist selbsttragend (kein `spec → adr`-Link,
`matrix` grün).

**4 — Die §3-Erweiterungen des ersten Laufs: je Datei gerechtfertigt.**
`internal/bootstrap/wiring.go` — ja, und zwar nachweisbar: ohne die
`cdc_admin`-gebundene Verbindung trägt `ChangeStorePort` nicht (Messung 11,
`nacharbeit-roles.sql`); sie folgt dem Muster von `retentionStore`.
`test/integration/integration_test.go` — ja, die Typänderung erzwingt sie.
`docs/user/e2e-abdeckung.md` — ja, das Erzeugnis folgt dem verschobenen Ort
seiner Quellen (Messung 10). `harness/image-hash.txt` — ja, form-korrekt nach
`ADR-0044` (Build-Kontext geändert → Digest-Commit; der Digest selbst ist
lauf-gebunden und nicht gegenprüfbar).

**5 — `tools/schema/**` unberührt: bestätigt.** Kein DDL, kein Schema-Eingriff;
`git diff --name-only` enthält keine Datei dieses Verzeichnisses.

**6 — Hygiene: sauber** (siehe Negativbefund „Hygiene"). Die zwei Abzüge
liegen nicht in der Hygiene, sondern in der Verteidigungslinie (F-2) und in
der ADR-Zusage (F-1).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Fitness-Function-Zusage ohne Träger ·
Negativtest ohne Bindung an die Eingabe · Kommentar nennt abwesende Kennung ·
Spec-Wortlaut überzeichnet die Implementierung · Formatierungsrest ·
vorgefundene Sicht-Drift · benannte Doppelung ohne Drift

## Verdikt

**Merge-blockierend:** nein — 0 HIGH; die gelieferte Semantik entspricht
`ADR-0081` in allen vier Teilfragen, das zweite Medium betrifft eine
Verteidigungs-, keine Verhaltenslücke, und F-1 verlangt eine
Planer-/Architect-Entscheidung (ADR-Zusage gegen Slice-Zuschnitt), keine
Codeänderung. Nicht still entschieden: **beide MEDIUMs brauchen einen
Ausgang**, bevor der Slice nach `done/` geht.

**Übergabe:** F-2 und F-3 gehen als **Rückkante Reviewer → Implementer**
zurück (Test-Bindung, Kommentar-Name); F-4 kann mit ihnen reisen. F-1 geht an
**Planner/Architect**: entweder die Fitness-Function-Zeile der `ADR-0081` wird
auf den tatsächlichen Träger gezogen (Folge-ADR), oder der E2E-Netzwerk-Beleg
bekommt einen benannten Folge-Slice. Die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Der Report selbst ist
Lauf-Beleg und wird über Läufe hinweg nicht gelesen; er ersetzt keine
Verifikation (Modul 11).

**DoD-Häkchen „Review durchgeführt":** wird **nicht** in diesem Commit
nachgezogen — F-2 verlangt eine Fixrunde am Implementer, die Checkbox gehört
damit in Schritt 21 des zweiten Implementer-Laufs
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde, Grenze).
