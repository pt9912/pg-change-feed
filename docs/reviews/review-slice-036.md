# Review-Report: slice-036 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-036`, §1/§2/§3/§6/§8) und
`ADR-0050` (Accepted, zentraler Prüfmaßstab dieses Laufs) sowie `AGENTS.md`
§3 Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `b6a1c9f` (neue Tabelle `cdc.administration_request`
in `tools/schema/schema.yaml`; neue SQL-Funktionen `cdc.enable_table`/
`cdc.disable_table` in `tools/schema/nacharbeit-administration.sql`;
neuer Testfall `internal/adapters/driven/postgresstorage/administrationrequest_test.go`;
`Makefile`-Anpassung Rollout-Reihenfolge; neuer Beleg
`docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/evidence/slice-036.md`),
`3cd1f64` (DoD-Häkchen und Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-036-antragsqueue-sql-funktionen.md` (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4–6, §8, Plan-Nachzug)
- `docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md` (vollständig — zentraler Prüfmaßstab)
- `docs/plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md`, `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` (referenziert von `ADR-0050`)
- `spec/lastenheft.md` (`LH-FA-ADM-001`)
- `tools/schema/nacharbeit-administration.sql`, `tools/schema/schema.yaml` (Diff vollständig gelesen, nicht nur Implementer-Zusammenfassung)
- `internal/adapters/driven/postgresstorage/administrationrequest_test.go` (vollständig)
- `Makefile` (Rollout-Reihenfolge, `schema-rollout`-Target)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/` (`observation.md`, `state.md`, alle fünf `evidence/*.md`-Dateien)
- `harness/README.md` §Sensors/Werkzeuge (Zeile `make schema-rollout`)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.3, §3.7
- d-migrate-Quellcode (lokal verfügbar unter `/Development/d-migrate`,
  Digest `sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89`,
  `hexagon/core/.../MigrationFingerprint.kt`, `RawSqlTextProjection.kt`) —
  zur unabhängigen Prüfung der d-migrate-Behauptung, nicht Teil des Diffs
- `docs/reviews/review-slice-035.md` (Format-Vorlage)

---

## Findings

- **[MEDIUM]** — `quelle`: Modul 6 (Beobachtungs-Register) / Maintainability
  `pfad`: `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/state.md`
  `befund`: Mit `evidence/slice-036.md` liegt die **fünfte** Beleg-Datei in
  diesem Verzeichnis vor, `state.md`s Notiz nennt aber weiterhin „Zähler
  (abgeleitet): 4× (evidence/slice-006.md, evidence/slice-010.md,
  evidence/slice-015.md, evidence/slice-016.md)" — ohne `slice-036.md`.
  Zusätzlich deckt der bisherige Text nur die zwei historisch betroffenen
  Fälle (CHECK-Constraint, Views) als „technisch aufgelöst" ab; der mit
  diesem Slice neu belegte **dritte** Fall (SQL-Funktionen, weiterhin
  aktiv über die Ausweichform, nicht aufgelöst) hat keine Entsprechung in
  der `state.md`-Notiz. Der Zähler selbst ist laut Regel abgeleitet und
  nie falsch (Modul 6 §Das Beobachtungs-Register), aber diese freitextliche
  Momentaufnahme ist es faktisch geworden, sobald eine neue Beleg-Datei
  hinzukam, ohne sie nachzuziehen — ein künftiger Sichtungs-Schritt (§8
  eines Folge-Slice), der nur `state.md` liest, sieht weder den aktuellen
  Belegstand noch den offenen Funktions-Fall.
  `verifizierbar`: nein (kein Sensor prüft Freitext-Konsistenz in `state.md`
  gegen die Dateizahl in `evidence/`)
  `klasse`: „Beobachtungs-Register-Notiz nicht mit neuem Beleg nachgezogen"
  (erstes Auftreten in diesem Skill-Lauf)

- **[MEDIUM]** — `quelle`: Maintainability (Reviewer-Skill §Klassifikation,
  „fehlende Negativtests bei neuem öffentlichem Vertrag")
  `pfad`: `internal/adapters/driven/postgresstorage/administrationrequest_test.go`
  `befund`: Der Plan-Nachzug behauptet „Real geprüft: ein Login ohne
  `cdc_admin`-Mitgliedschaft scheitert an `SELECT cdc.enable_table(...)`
  mit „permission denied for function", ein Login mit Mitgliedschaft
  gelingt" — im Diff existiert dafür kein automatisierter Testfall. Die
  beiden vorhandenen Tests rufen die Funktionen ausschließlich über den
  Pool auf, dessen Rolle nicht geprüft/variiert wird. Die
  Least-Privilege-Durchsetzung (`REVOKE …/GRANT … TO cdc_admin`,
  SQL-seitig korrekt implementiert, siehe Negativbefunde) ist damit ein
  neuer sicherheitsrelevanter öffentlicher Vertrag ohne Regressionsschutz
  — eine künftige Änderung an `nacharbeit-administration.sql` oder
  `nacharbeit-roles.sql`, die das Grant versehentlich lockert, würde von
  keinem Gate/Test erkannt.
  `verifizierbar`: ja — ein Testfall mit einer `cdc_reader`-DSN gegen
  `SELECT cdc.enable_table(...)` würde das reproduzieren
  `klasse`: „fehlende Negativtests bei neuem öffentlichem Vertrag"

- **[MEDIUM]** — `quelle`: Maintainability (DoD-Begründung)
  `pfad`: `docs/plan/planning/in-progress/slice-036-antragsqueue-sql-funktionen.md`
  (DoD-Punkt 4, Plan-Nachzug)
  `befund`: Die Begründung, warum `harness/README.md` §Sensors nicht
  aktualisiert wurde, lautet: „die Zeile nennt schon heute keine
  einzelnen `nacharbeit-*.sql`-Dateien namentlich, das bleibt
  Makefile-Kommentar-Ebene". Das ist unzutreffend — die `make
  schema-rollout`-Zeile in `harness/README.md` (§Werkzeuge) nennt bereits
  explizit eine einzelne Datei beim Namen: „Die Ausweichform
  `nacharbeit-views.sql` ist zurückgebaut · seit slice-016". Es gibt damit
  einen Präzedenzfall dafür, dass Status-Änderungen einzelner
  `nacharbeit-*.sql`-Dateien in dieser Zeile auftauchen. Die DoD-Aussage
  selbst könnte trotzdem im Ergebnis vertretbar sein (die Zeile
  dokumentiert bislang nur *aufgelöste* Fälle, nicht neu entstandene
  Ausweichformen) — aber die angegebene *Begründung* trägt nicht, und eine
  nicht geprüfte Übernahme dieser Begründung hätte die Unstimmigkeit nicht
  aufgedeckt.
  `verifizierbar`: ja (Datei-Diff von `harness/README.md` gegen den
  zitierten Text)
  `klasse`: „DoD-Begründung mit unzutreffender Tatsachenbehauptung"

## Negativbefunde

- geprüft, ohne Befund: **`cdc.enable_table`/`cdc.disable_table` schreiben
  ausschließlich `cdc.administration_request` und rufen `pg_notify`** —
  eigenständig im SQL-Quelltext nachvollzogen (`tools/schema/nacharbeit-administration.sql`),
  nicht aus dem Implementer-Bericht übernommen. Kein `INSERT`/`UPDATE`/
  `DELETE` auf `cdc.source_table`/`cdc.schema_version`, kein `ALTER
  PUBLICATION`, kein Aufruf einer anderen schreibenden Funktion. Erfüllt
  `ADR-0050`s zentrale Entscheidung und die Slice-Abgrenzung (§1) sowie
  die Fitness-Function-Zeile 1 aus `ADR-0050`.
- geprüft, ohne Befund: **Rollenwahl/Least-Privilege (`ADR-0047`)
  wirksam.** `REVOKE EXECUTE ON FUNCTION cdc.enable_table(...),
  cdc.disable_table(...) FROM PUBLIC;` läuft **nach** beiden
  `CREATE OR REPLACE FUNCTION`-Anweisungen (die implizit `EXECUTE TO
  PUBLIC` vergeben) und **vor** dem gezielten `GRANT … TO cdc_admin` —
  kein impliziter PUBLIC-Zugriff bleibt zwischen den Anweisungen bestehen,
  da alles in einem einzigen, sequenziell abgearbeiteten `psql -f`-Lauf
  steht. `CREATE OR REPLACE FUNCTION` erhält laut PostgreSQL-Semantik
  bestehende Grants bei erneutem Rollout (Idempotenz gewahrt). Real
  bestätigt: eigener `make test-store`-Lauf in dieser Sitzung zeigt
  `REVOKE`/`GRANT` ohne Fehler (`-v ON_ERROR_STOP=1`).
- geprüft, ohne Befund: **Testfall ist real, keine Attrappe.**
  `administrationrequest_test.go` überspringt ohne `CDC_STORE_TEST_DSN`,
  baut für `LISTEN` eine dedizierte `pgx.Connect`-Verbindung (korrekt —
  Pool-Verbindungen sind für `WaitForNotification` ungeeignet), ruft die
  SQL-Funktion über den Pool auf und prüft sowohl die Notify-Payload als
  auch die resultierende Zeile (`request_kind`, `status = pending`) via
  Datenbank-Lesezugriff. Eigenständig gegen eine reale PostgreSQL-18-Instanz
  nachvollzogen (siehe unten, `make test-store`).
- geprüft, ohne Befund: **d-migrate-`POST_EXECUTE_DRIFT`-Behauptung ist
  zutreffend — verifiziert auf zwei unabhängigen Wegen, nicht aus dem
  Implementer-Bericht übernommen.**
  1. **Quellcode-Analyse** (d-migrate liegt lokal unter `/Development/d-migrate`
     vor, exakt der gepinnte Digest): `MigrationFingerprint.appendFunctions`
     (`hexagon/core/.../MigrationFingerprint.kt:605-627`) nimmt
     `fn.body`/`fn.sourceDialect` **unverändert** in den Post-Compare-
     Fingerabdruck auf. `RawSqlTextProjection.blank()` — die Stelle, die
     laut Changelog 1.3.1 genau dieses Problem für Sichten behoben hat
     (`view.query`/`view.sourceDialect` werden vor dem Fingerabdruck
     ausgeblendet) — deckt laut eigenem Quelltext-Kommentar „fünf
     Stellen" ab: `ViewDefinition.query`, `ConstraintDefinition.expression`,
     `IndexDefinition.where`, `IndexColumn.expression`,
     `ColumnGeneration.Computed.expression`. **Funktionen/Prozeduren/
     Trigger sind nicht darunter** — ihr Rumpf bleibt im Fingerabdruck,
     obwohl PostgreSQL ihn beim Zurücklesen umformatiert (`pg_get_functiondef`),
     genau die Bedingung, die für Sichten schon einmal zu genau diesem
     Fehlerbild führte.
  2. **Eigene Reproduktion** (unabhängig vom Implementer-Lauf, eigener
     Docker-Container, dieselbe gepinnte PostgreSQL-18- und d-migrate-
     1.3.1-Version): ein Schema mit nur einer Tabelle rollt mit Exit 0;
     dasselbe Schema mit zusätzlich einer trivialen No-Arg-`plpgsql`-Funktion
     bricht mit `POST_EXECUTE_DRIFT` (Exit 5) ab — reale
     `docker run`-Läufe, kein Konstrukt. Damit ist die Implementer-Aussage
     nicht nur plausibel, sondern positiv verifiziert.
- geprüft, ohne Befund: **Beleg `evidence/slice-036.md` liegt im selben
  Beobachtungs-Rahmen wie die vier vorherigen Belege dieses
  Verzeichnisses.** Alle fünf Belege (`slice-006`, `slice-010`,
  `slice-015`, `slice-016`, `slice-036`) beschreiben dieselbe
  Fehlerklasse: der Post-Compare von `schema migrate --execute` vergleicht
  rohen, vom Server umformatierten SQL-Text gegen den Autorentext und
  meldet Drift, wo keine ist — zuerst bei CHECK-Ausdrücken, dann bei
  Sichten, jetzt bei Funktionen. Dieselbe Ausweichform
  (berichtete manuelle Nacharbeit über `nacharbeit-*.sql`) trägt alle drei
  Fälle. Kein neu erfundener Bucket — siehe auch das oben unabhängig
  bestätigte Quellcode-Muster (`RawSqlTextProjection` deckt Tabellen/
  Sichten ab, nicht Routinen — dieselbe Lücken-Familie).
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7).** Kein
  Kommentar in `nacharbeit-administration.sql`, `schema.yaml`, `Makefile`
  oder `administrationrequest_test.go` referenziert `slice-037`/`slice-038`
  namentlich; Verweise auf künftige Arbeit bleiben bei „Folge-Slice"
  (generisch, kein Vorwärtsverweis auf eine konkrete Kennung). Alle
  Kommentare beschreiben den geltenden Zustand im Indikativ mit
  auflösbaren Ankern (`ADR-0050`, `ADR-0043`, `LH-FA-ADM-001`,
  `BEO-PGC/d-migrate-nacharbeit`), keine Chronik einer verworfenen
  Alternative, kein Abbruch mitten im Satz.
- geprüft, ohne Befund: **DoD-Aktualisierung, bis auf das oben gemeldete
  MEDIUM-Finding zu DoD-Punkt 4.** Die übrigen `[x]`-Punkte sind durch den
  Diff und eigene Testläufe in dieser Sitzung gedeckt (Tabelle real
  ausgerollt, Funktionen real getestet, `make gates` real grün). Die vier
  offen gelassenen Punkte (Review, Closure-Notiz, Risiko-Ausgänge, drei
  Paarungen) sind korrekt unbeansprucht — Planner-Closure-Arbeit nach
  Modul 8, kein Self-Review. Das Reconciliation-Item ist korrekt als
  „entfällt" markiert (Repo durchgehend GF laut `harness/conventions.md`).
- geprüft, ohne Befund: **§8 Sub-Area-Sichtung.** Einzige berührte
  Sub-Area `*`/`PGC`, korrekt als GF eingestuft. Der zitierte Treffer
  `BEO-PGC/verwaltung-keine-sql-administration` existiert im Register mit
  Stand „0×, benannt". `BEO-PGC/d-migrate-nacharbeit` wird in §8 nicht
  separat aufgeführt, obwohl dieser Slice ihn mit einer neuen Beleg-Datei
  bedient — kein Fehler in der Sichtung selbst (die Sichtung lief vor der
  Implementierung, der d-migrate-Fund entstand erst während der
  Implementierung), aber siehe das erste MEDIUM-Finding zur
  Register-Nachpflege.
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — jeder Rollout-/
  Testschritt läuft über `docker run`/`make`-Targets, kein lokales
  Toolchain-Install. **3.3** (git-mv/Inhalt-Trennung) — nicht einschlägig,
  kein Lifecycle-Übergang in diesem Diff. **3.5/3.6** — kein Accepted-ADR
  im Diff verändert, keine Gate-Schwelle gelockert.
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `ADR-0050`, kein `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability`
  lief in dieser Sitzung über `HEAD~5..HEAD` grün (0 Befunde).
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf,
  eigenständig ausgeführt) — `baseline-verify` (54 Dateien),
  `d-check` (297 Dateien, 0 Befunde — 1 Datei mehr als im Implementer-Bericht,
  konsistent mit den seit dessen Lauf hinzugekommenen Review-Vorlage-Läufen,
  kein Befund dadurch berührt), `commit-traceability` (5 Commits, OK),
  `a-check` (0 Befunde) — alle grün.
- geprüft, ohne Befund: **`make test-store`, real ausgeführt** (diese
  Review-Sitzung, nicht aus dem Implementer-Bericht übernommen) —
  `schema-validate`/`schema-rollout` inkl. `nacharbeit-administration.sql`
  liefen fehlerfrei (`CREATE FUNCTION` ×2, `REVOKE`, `GRANT` ohne Fehler),
  das komplette Go-Testpaket lief grün, einschließlich
  `internal/adapters/driven/postgresstorage` (3.007s, beide neuen
  Testfälle eingeschlossen).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Beobachtungs-Register-Notiz nicht mit
neuem Beleg nachgezogen" (1×, erstes Auftreten), „fehlende Negativtests bei
neuem öffentlichem Vertrag" (1×, bereits benannte Skill-Klasse), „DoD-
Begründung mit unzutreffender Tatsachenbehauptung" (1×, erstes Auftreten) —
kein Steering-Loop-Eintrag fällig, Schwelle ist 3×.

## Verdikt

**Merge-blockierend:** keins. Alle drei Findings sind MEDIUM ohne
Rollen-Widerspruch — keine Architect-Sequenz nötig (Modul 8; die Sequenz
greift erst ab HIGH mit Rollen-Widerspruch oder ab dem dritten gleichen
Konflikttyp).

**Zur zentralen Prüffrage dieses Laufs (`ADR-0050`-Konformität):** Die
beiden SQL-Funktionen sind `ADR-0050`-konform — eigenständig im SQL-Quelltext
verifiziert, nicht aus dem Implementer-Bericht übernommen: kein direkter
Zugriff auf `cdc.source_table`/`cdc.schema_version`/die Publication, nur
Antrags-Zeile plus `pg_notify`. Die Rollenwahl (`cdc_admin`, `ADR-0047`)
ist korrekt und wirksam umgesetzt (REVOKE-vor-GRANT-Reihenfolge schließt
den impliziten PUBLIC-Zugriff). Die d-migrate-`POST_EXECUTE_DRIFT`-Behauptung
ist nicht nur plausibel, sondern über eine unabhängige Quellcode-Analyse
und eine eigene Reproduktion positiv bestätigt — echter, bislang
ungepatchter Deckungslücken-Bug in d-migrate 1.3.1 (Funktionen/Prozeduren/
Trigger fehlen dort, wo Sichten seit 1.3.1 bereits behoben sind), die
Ausweichform ist gerechtfertigt.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** der SQL-Quelltext der beiden Funktionen, die Wirksamkeit der
REVOKE/GRANT-Reihenfolge, die d-migrate-Fingerabdruck-Logik im
Quellcode sowie eine eigene, isolierte Reproduktion des
`POST_EXECUTE_DRIFT`-Fehlers, und `make gates`/`make test-store` als
eigenständige Läufe in dieser Sitzung.

**Übergabe:** Drei MEDIUM-Findings, kein Rollen-Konflikt. Der Implementer
entscheidet über Annahme oder Begründung (Modul 8: die Sequenz mit
Architect-Übergabe ist bei isoliertem MEDIUM Overkill). Dieser Report ist
ein Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen — die
Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5). Er ersetzt
keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
