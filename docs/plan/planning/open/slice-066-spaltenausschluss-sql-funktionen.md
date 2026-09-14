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

**Verantwortlich:** —.

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

- [ ] `cdc.administration_request` erweitert (neue Spalte für den
      Spaltennamen, erweiterte `request_kind`-CHECK-Klausel um
      `exclude_column`/`include_column`) — über `tools/schema/schema.yaml`
      **oder**, falls d-migrate für eine reine Spalten-/CHECK-Änderung an
      einer bestehenden Tabelle real scheitert (Präzedenzfall
      `BEO-PGC/d-migrate-nacharbeit`, dort nur für Funktionen belegt, nicht
      für Tabellenänderungen), über eine geeignete Migrationsform —
      Implementer-Entscheidung, Plan-Nachzug. Beleg: real ausgerollt über
      `make schema-rollout`.
- [ ] `cdc.exclude_column(...)`/`cdc.include_column(...)` real als SQL-
      Funktionen angelegt (`tools/schema/nacharbeit-administration.sql`,
      analog `cdc.enable_table`/`cdc.disable_table` — dieselbe
      d-migrate-Ausweichform, `BEO-PGC/d-migrate-nacharbeit`); schreiben
      ausschließlich einen Antrags-Datensatz und senden `pg_notify`. Neuer
      Outbound Port mit `ColumnExists`, neue Inbound Ports
      (`ExcludeColumnUseCase`/`IncludeColumnUseCase`), `ErrSourceColumnMissing`.
      `applyAdministrationRequest` um `exclude_column`/`include_column`
      erweitert. Beleg: `internal/adapters/driven/postgresstorage/administrationrequest_test.go`
      (neue Testfälle, real gegen PostgreSQL: Happy Path `applied`,
      Negative-Fall nicht existierende Spalte `failed` mit Fehlertext).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `spec/architecture.md`s Sequenzdiagramm zu `ARC-005`
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
| `tools/schema/schema.yaml` | update | neue Spalte + `request_kind`-CHECK-Erweiterung an `cdc.administration_request` |
| `tools/schema/nacharbeit-administration.sql` | update | zwei neue SQL-Funktionen `cdc.exclude_column`/`cdc.include_column` |
| `internal/application/port/outbound/` | neu | neuer Outbound Port mit `ColumnExists` |
| `internal/application/port/inbound/verwaltung.go` | update | `ExcludeColumnUseCase`/`IncludeColumnUseCase`, `ErrSourceColumnMissing` |
| `internal/application/service/` (o. ä.) | neu | `ExcludeColumnService`/`IncludeColumnService` |
| `internal/bootstrap/wiring.go` | update | `applyAdministrationRequest` um zwei `case`-Zweige |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update | Tests real gegen PostgreSQL |
| `spec/architecture.md` | update | `ARC-005`-Sequenzdiagramm-Korrektur (`ADR-0059` Folgepflicht) |
| `spec/pflichtenheft.md` | update | neuer `SPEC-*`-Eintrag (`ADR-0059` Folgepflicht) |

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
