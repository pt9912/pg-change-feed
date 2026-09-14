# Slice slice-066: Spaltenausschluss — SQL-Funktionen und Antrags-Verarbeitung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-18`](../welle-18.md) — erster Slice; `slice-067` (Filterung)
und `slice-068` (E2E-Beleg) bauen auf ihm auf, ohne die Anträge blieben ohne
Wirkung.

**Bezug:** [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Haupt-Bezug),
[`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Messmethode verweist
ausschließlich auf `LH-FA-CFG-005`),
[`ADR-0059`](../../../../docs/plan/adr/0059-spaltenauswahl-mechanismus.md)
(nur umgesetzt — keine aktive ADR wird geändert, `ADR-0059` bleibt
`Accepted`), [`ADR-0050`](../../../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(die erweiterte Antrags-Queue), [`ADR-0028`](../../../../docs/plan/adr/0028-inbound-use-cases.md)/[`ADR-0034`](../../../../docs/plan/adr/0034-ports-nach-faehigkeiten.md)
(Inbound-/Outbound-Port-Muster für die zwei neuen Fähigkeiten).

**Berührte Spec-Stellen:** [`ARC-005`](../../../../spec/architecture.md)
(Sequenzsicht SQL-Administration → Antragsqueue → Administrations-
Hintergrundzug — bekommt zwei weitere Antragsarten auf demselben Pfad,
`ADR-0059` Folgepflicht), [`ARC-002`](../../../../spec/architecture.md)
(zwei neue Use Cases), [`ARC-004`](../../../../spec/architecture.md) (ein
neuer Outbound Port). Dieser Slice trägt zusätzlich die `ADR-0059`
Folgepflicht, einen neuen `SPEC-*`-Eintrag in `spec/pflichtenheft.md` für
die konkrete Feldform des erweiterten Antrags-Datensatzes anzulegen — hier
zugeordnet, weil die konkrete Feldform (Spaltenname-Spalte,
`request_kind`-Erweiterung) in diesem Slice erstmals entsteht.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `ADR-0059`s Antrags-Seite für Spaltenausschluss/-einschluss real
bauen: `cdc.administration_request` (`tools/schema/schema.yaml`) um eine
Spalte für den Spaltennamen und eine erweiterte `request_kind`-CHECK-
Klausel (`exclude_column`/`include_column` neben `enable`/`disable`)
erweitern; zwei neue SQL-Funktionen `cdc.exclude_column(schema, table,
column)`/`cdc.include_column(schema, table, column)`, die **ausschließlich**
einen Antrags-Datensatz schreiben und `pg_notify` senden (`ADR-0018`/
`ADR-0046`-Disziplin, kein Prüf-/Domänenlogik-Anteil); ein neuer Outbound
Port mit `ColumnExists`-Fähigkeit (analog
`TableActivationPort.TableExists`); neue Inbound Ports
(`ExcludeColumnUseCase`/`IncludeColumnUseCase`); ein neuer Sentinel
`ErrSourceColumnMissing` (Muster `ErrSourceTableMissing`,
`internal/application/port/inbound/verwaltung.go`); Erweiterung von
`applyAdministrationRequest` (`internal/bootstrap/wiring.go`) um die zwei
neuen `case`-Zweige, die bei Erfolg **noch keine** `Assembler`-Bindung
nachtragen (das übernimmt `slice-067`s Live-Reload-Methode — hier reicht der
Antrag als `applied`/`failed`). Real gegen PostgreSQL getestet: der
Funktionsaufruf legt real eine Zeile mit Status `pending` an, die
Administrations-Goroutine verarbeitet sie zu `applied` bzw. `failed`
(Negative-Fall: nicht existierende Spalte).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`Assembler`-Filterung der Row-Image-Konstruktion und die neue
  synchronisierte Live-Reload-Methode** — Folge-Slice `slice-067` übernimmt
  das explizit; dieser Slice liefert nur die Antrags-Seite, ohne dass
  irgendetwas den verarbeiteten Antrag bislang in eine Filterwirkung
  übersetzt.
- **E2E-Beleg am laufenden Feed-Container** — eigener Slice `slice-068`
  (braucht `slice-066` **und** `slice-067` zusammen, siehe `welle-18` §5).
- **Quellen-/musterweiter Ausschluss über mehrere Tabellen hinweg** —
  Bestand bleibt bewusst außen vor: `ADR-0059` Teilfrage 2 schließt Option
  C bewusst aus (eigener Re-Evaluierungs-Trigger), diese Welle liefert
  ausschließlich Granularität pro Tabelle.
- **Direktes Schreiben von `cdc.source_table`/Publication durch die neuen
  SQL-Funktionen** — Bestand bleibt bewusst außen vor, dieselbe
  `ADR-0050`-Disziplin wie bei `enable_table`/`disable_table`: eine
  SQL-Funktion, die das umginge, widerspräche `ADR-0018`/`ADR-0046` direkt.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `cdc.administration_request` erweitert (neue Spalte für den
      Spaltennamen, erweiterte `request_kind`-CHECK-Klausel um
      `exclude_column`/`include_column`) — über `tools/schema/schema.yaml`
      **oder**, falls d-migrate für eine reine Spalten-/CHECK-Änderung an
      einer bestehenden Tabelle real scheitert (Präzedenzfall
      `BEO-PGC/d-migrate-nacharbeit`, dort nur für Funktionen belegt, nicht
      für Tabellenänderungen), über eine geeignete Migrationsform —
      Implementer-Entscheidung, Plan-Nachzug. Beleg: real ausgerollt über
      `make schema-rollout`.
- [x] `cdc.exclude_column(...)`/`cdc.include_column(...)` real als SQL-
      Funktionen angelegt (`tools/schema/nacharbeit-administration.sql`,
      analog `cdc.enable_table`/`cdc.disable_table` — dieselbe
      d-migrate-Ausweichform, `BEO-PGC/d-migrate-nacharbeit`); schreiben
      ausschließlich einen Antrags-Datensatz und senden `pg_notify`. Neuer
      Outbound Port mit `ColumnExists`, neue Inbound Ports
      (`ExcludeColumnUseCase`/`IncludeColumnUseCase`), `ErrSourceColumnMissing`.
      `applyAdministrationRequest` um `exclude_column`/`include_column`
      erweitert. Beleg: `internal/bootstrap/administration_endtoend_test.go`
      (neu, real gegen PostgreSQL: Happy Path `applied` über
      `cdc.exclude_column` mit vorhandener Spalte, Negative-Fall `failed`
      samt Fehlertext über `cdc.include_column` mit nicht existierender
      Spalte); `internal/adapters/driven/postgresstorage/administrationrequest_test.go`
      (neue Testfälle, real gegen PostgreSQL: Anlage über beide
      SQL-Funktionen, Rücklesen von `column_name`/`request_kind` und
      `cdc_reader`-Ablehnung).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-066.md`](../../../../docs/reviews/review-slice-066.md)
      (Fixrunde zu F-1…F-8 geprüft; F-9/F-10 gehen als Closure-Nachzug an
      den Planner, keine weitere Implementer-Runde).
- [x] Doku-Update: `spec/architecture.md`s Sequenzdiagramm zu `ARC-005`
      um die beiden neuen Antragsarten ergänzt (`ADR-0059` Folgepflicht);
      neuer `SPEC-*`-Eintrag in `spec/pflichtenheft.md` für die konkrete
      Feldform des erweiterten Antrags-Datensatzes (`ADR-0059`
      Folgepflicht).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb (`welle-18` offen) — Prüfung läuft bei der `welle-18`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/schema.yaml` | update | neue Spalte `column_name` an `cdc.administration_request`; die `request_kind`-CHECK-Klausel verlässt das deklarative Modell (Plan-Nachzug unten) |
| `tools/schema/nacharbeit-administration.sql` | update | vier SQL-Funktionen `cdc.exclude_column`/`cdc.include_column` zusätzlich zu `cdc.enable_table`/`cdc.disable_table`; `request_kind`-CHECK-Klausel mit den vier Antragsarten |
| `tools/schema/{down.sql,plan.yaml}` | update | generierte Rollout-Artefakte zu obigen Schema-Änderungen |
| `internal/application/port/outbound/columnexclusion.go` | neu | neuer Outbound Port mit `ColumnExists` |
| `internal/application/port/outbound/administrationrequest.go` | update | Port-Doku auf die vier Antragsarten |
| `internal/application/port/inbound/verwaltung.go` | update | `ExcludeColumnUseCase`/`IncludeColumnUseCase`, `ErrSourceColumnMissing` |
| `internal/application/usecase/excludecolumn/service.go` | neu | `ExcludeColumnService` |
| `internal/application/usecase/includecolumn/service.go` | neu | `IncludeColumnService` |
| `internal/domain/model/administrationrequest.go` | update | `Column`-Feld, zwei Werte der geschlossenen Menge, Invarianten-Prüfung im Konstruktor |
| `internal/domain/errors/errors.go` | update | Kommentar der geschlossenen Antragsarten-Menge (`ErrInvalidAdministrationRequestKind`) auf die vier Werte nachgezogen — der Sentinel `ErrSourceColumnMissing` liegt in `internal/application/port/inbound/verwaltung.go` |
| `internal/adapters/driven/postgresstorage/{administrationrequest,queries}.go` | update | `column_name` in Lese-/Schreibweg und Query-Text |
| `internal/adapters/driven/postgresstorage/tableactivation.go` | update | `ColumnExists`-Implementierung des neuen Ports |
| `internal/bootstrap/wiring.go` | update | `applyAdministrationRequest` um zwei `case`-Zweige |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update | Tests real gegen PostgreSQL |
| `internal/application/usecase/{exclude,include}column/service_test.go` | neu | Use-Case-Tests |
| `internal/bootstrap/administration_endtoend_test.go` | neu | Ende-zu-Ende-Beleg des Spaltenwegs real gegen PostgreSQL |
| `internal/bootstrap/administration_internal_test.go` | update | `MarkFailed`-Zweig für die beiden Spalten-Antragsarten |
| `internal/domain/model/administrationrequest_test.go` | neu | Invarianten-Tests des Konstruktors |
| `spec/architecture.md` | update | `ARC-005`-Sequenzdiagramm-Korrektur (`ADR-0059` Folgepflicht) |
| `spec/pflichtenheft.md` | update | neuer `SPEC-*`-Eintrag (`ADR-0059` Folgepflicht) |

**Plan-Nachzug (Implementer, 2026-09-14) — Code-Umfang der Tabelle.** Die
Tabelle trug vor der Umsetzung zwei Sammelpfade als Platzhalter
(`internal/application/port/outbound/`, `internal/application/service/`
(o. ä.)); implementiert wurden sie an den oben eingetragenen konkreten
Pfaden, dazu die produktiv nötigen Folgeänderungen an
`internal/domain/model/administrationrequest.go`,
`internal/domain/errors/errors.go`, den beiden
`postgresstorage`-Dateien und dem `tableactivation`-Adapter
(`ColumnExists`). Die Tabelle ist auf diesen Stand nachgezogen.

**Plan-Nachzug (Implementer, 2026-09-14) — Migrationsform der
`request_kind`-CHECK-Klausel.** Die DoD-Zeile zu
`cdc.administration_request` hat den Fall vorgesehen („oder, falls d-migrate
für eine reine Spalten-/CHECK-Änderung an einer bestehenden Tabelle real
scheitert … über eine geeignete Migrationsform — Implementer-Entscheidung,
Plan-Nachzug"); er ist eingetreten. Real gemessen gegen eine bestehende
Instanz (PostgreSQL 18, d-migrate 1.3.1): die Spalten-Erweiterung
(`column_name`) konvergiert (Exit 0), die CHECK-Erweiterung **nicht** —
`POST_EXECUTE_DRIFT` (Exit 5), und die bestehende Klausel entfällt dabei;
dasselbe Bild mit umbenanntem Constraint, also unabhängig vom Namen, und
ohne dass die neue Klausel entsteht. Die Anlage einer neuen CHECK-Klausel an
einer bestehenden Tabelle ist von d-migrate 1.3.1 also nicht getragen. Der
Implementer hat deshalb die Klausel **aus dem deklarativen Modell
genommen** (`tools/schema/schema.yaml` trägt nur noch
`chk_administration_request_status`) und in die etablierte Ausweichform
`tools/schema/nacharbeit-administration.sql` gelegt — idempotent
(`DROP CONSTRAINT IF EXISTS` vor `ADD CONSTRAINT`), dieselbe Form wie die
Antrags-Funktionen. Beleg (real gemessen von der Reviewer-Rolle gegen
PostgreSQL 18 / d-migrate 1.3.1, dokumentiert in
[`docs/reviews/review-slice-066.md`](../../../reviews/review-slice-066.md)
§Eigene Nachmessung; vom umsetzenden Lauf übernommen, nicht dort erneut
gefahren): ein **frischer Rollout** mit dem aktuellen Baum läuft Exit 0; danach trägt
`cdc.administration_request` die Spalte `column_name`, und
`chk_administration_request_kind` entsteht anschließend mit den vier Arten
(`nacharbeit-administration.sql:58` meldet zuvor
`NOTICE: constraint … does not exist, skipping`). Die Ausweichform ist
idempotent — dreimaliges `psql -f nacharbeit-administration.sql` endet
dreimal Exit 0 mit unverändertem Endzustand. Ein Rollout mit dem aktuellen
Baum gegen eine Instanz auf dem Stand *vor* dieser Änderung endet dagegen
**nicht** mit Exit 0, sondern mit **Exit 8**
(`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`); `make` bricht vor den
psql-Nacharbeitsschritten ab, der Endzustand (Spalte + vierwertige Klausel)
entsteht so nicht. Der Blocker ist nicht neu — derselbe Lauf mit dem
*Eltern*-Modell gegen dieselbe Instanz endet ebenfalls Exit 8; die seit
`slice-011`/`slice-012`/`slice-036` undeklarierten Nacharbeitsobjekte
(Funktionen/Views) tragen ihn. **Neu** durch diesen Diff ist eine
zusätzliche destruktive Plan-Operation auf
`administration_request.request_kind` (`AlterColumnType`, erste Anweisung
`ALTER TABLE … DROP CONSTRAINT IF EXISTS "chk_administration_request_kind"`):
sie entsteht erst dadurch, dass die Klausel aus dem deklarativen Modell
genommen wurde, während sie im Katalog liegt; manuelles Entfernen der
Klausel lässt den Plan von 11 auf 10 Operationen fallen und genau diese
Operation verschwinden. Die Ausweichform macht ein Objekt, das d-migrate im
Katalog vorfindet und das Modell nicht mehr deklariert, für den
Bestands-Rollout damit zum **destruktiven** Plan-Eintrag; der reale Ausgang
dieser Upgrade-Konstellation ist Exit 8, kein Exit 0. Die
`column_name`-Spalte bleibt bewusst deklarativ — sie konvergiert. Das ist
eine **neue Objektklasse** derselben Beobachtung
(`BEO-PGC/d-migrate-nacharbeit`): „CHECK-Änderung/-Anlage an einer
bestehenden Tabelle", verschieden von der seit `slice-015` deklarativ
gelösten Erstanlage-Konvergenz; §8 dieser Planung hatte „keinen neuen
Zähler-Beitrag" erwartet, solange keine neue Objektklasse betroffen ist —
sie ist betroffen, der Beleg gehört in die Closure.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Schema-Erweiterung, beide SQL-Funktionen, neuer Port und beide Use Cases
  zusammen mehr als drei Liefer-Punkte oder mehr als zwei Schichten in
  einer Review-Sitzung nicht mehr prüfbar machen, gehört das zurück zur
  Zerlegung (z. B. `exclude`/`include` als getrennte Slices).
- `in-progress` → `open` (blockiert — Carveout?): d-migrate scheitert real
  auch an der reinen Schema-Erweiterung von `cdc.administration_request`
  (nicht nur an neuen Funktionen) und es gibt keine saubere
  Ausweichform ohne Datenverlust am Bestand — dann Carveout statt
  stillem Workaround.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Adapter-Test gegen reale
PostgreSQL grün (Happy Path `applied`, Negative-Fall `failed`) **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- d-migrate könnte auch an einer reinen Spalten-/CHECK-Erweiterung einer
  bestehenden Tabelle scheitern (der bislang belegte
  `POST_EXECUTE_DRIFT`-Fall aus `BEO-PGC/d-migrate-nacharbeit` betrifft nur
  Funktionen/Prozeduren, nicht Tabellenänderungen — unklar, ob dieselbe
  Grenze auch hier greift). — **Ausgang:** <bei Closure zuzuweisen>
- Der Fehlerpfad für `ErrSourceColumnMissing` folgt dem
  `applyAdministrationRequest`-Muster mit Prozess-Fortsetzung bei
  Adapter-Fehlern (`BEO-PGC/adapter-fehler-ausgang`, 2×, weiter offen) —
  ein dritter, ähnlich gelagerter Fund würde die 3×-Schwelle erreichen. —
  **Ausgang:** <bei Closure zuzuweisen>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <bei Closure>
- **Was ging anders als geplant:** <bei Closure>
- **Steering-Loop-Eintrag:** <bei Closure>
- **Beobachtungs-Register (`../observations/`):** <bei Closure>
- **Folge-Slices:** keine — `slice-067`/`slice-068` stehen bereits in
  `welle-18` §4 als vorgesehene nächste Slices.
- **Risiken aus §6:** <bei Closure>
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-18` offen) —
  Prüfung läuft bei der `welle-18`-Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC` mit Bezug zu Antrags-Queue/SQL-Administration/Schema-
Erweiterung: `BEO-PGC/d-migrate-nacharbeit` (5×, verkörpert — dritte Klasse
„SQL-Funktionen/Prozeduren" bleibt real offen; dieser Slice fügt zwei
weitere Funktionen genau dieser Klasse hinzu und nutzt bewusst dieselbe
etablierte Ausweichform `nacharbeit-administration.sql`, kein neuer
Zähler-Beitrag zu erwarten, solange keine neue Objektklasse betroffen ist),
`BEO-PGC/adapter-fehler-ausgang` (2×, weiter offen — der neue
`ErrSourceColumnMissing`-Pfad folgt demselben `applyAdministrationRequest`-
Verarbeitungsmuster wie die bestehenden Fälle; kein dritter Fund a priori
zu erwarten, aber im Blick zu behalten, siehe §6),
`BEO-PGC/rollen-verdrahtung` (eingetreten, aufgelöst seit `slice-023` —
keine neue Rollenfrage, `cdc_admin`-Grant-Muster aus `slice-036` wird
unverändert übernommen), `BEO-PGC/verwaltung-keine-sql-administration`
(verkörpert seit `welle-12` — dieser Slice erweitert den bereits
verkörperten Mechanismus, keine neue Lücke). Keiner der Treffer erreicht
mit diesem Slice 3× erstmals oder verlangt eine Sonderbehandlung über die
in §6 benannten Risiken hinaus.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
