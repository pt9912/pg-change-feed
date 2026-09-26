# Slice transformationen-antragsweg-schema: Antragsweg, Schema — Antragsarten `set_transformation`/`remove_transformation`, Spalten `rule_name`/`rule_spec`, zwei SQL-Funktionen, Grants, Idempotenz-Guard

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Konfiguration der
Transformation), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md)
(Administration über SQL), [`LH-QA-SEC-002`](../../../../spec/lastenheft.md)
(Least-Privilege der Rollen),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 1 und Folgepflicht 3,
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue und Live-Reload),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
(Schemamigrationen — Idempotenz-Guard),
[`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) (Rollen),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(keine Domänenlogik in SQL).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md)
(Antrags-Datensatz, durch `slice-transformationen-spec-nachzug`),
[`ARC-005`](../../../../spec/architecture.md) (SQL-Funktionen als Driving
Adapter) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Antrags-Queue nimmt die zwei neuen Antragsarten entgegen: das
Schema trägt die Spalten `rule_name` und `rule_spec`, die geschlossene
`request_kind`-Menge trägt zwei Werte mehr, und die SQL-Funktionen
`cdc.set_transformation(source_id, schema_name, table_name, rule_name,
rule_spec jsonb)` und `cdc.remove_transformation(source_id, schema_name,
table_name, rule_name)` — die Parameterlisten legt dieser Slice fest, die Spec
([`SPEC-019`](../../../../spec/pflichtenheft.md)) nennt die Funktionen ohne
Parameter —
schreiben **nur** den Antrag und senden `pg_notify` —
beide ausschließlich der Rolle `cdc_admin` ausführbar (`SECURITY DEFINER`,
gepinnter `search_path`, `REVOKE … FROM PUBLIC`). Die Domäne kennt die
Antragsarten, der Store liest die zwei Spalten. Die Verarbeitung eines solchen
Antrags trägt der Folge-Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Use Cases, K1–K4, Regelstand-Ableitung, Verdrahtung in
  `applyAdministrationRequest`** — `antragsweg-usecase`. Bis dahin endet ein
  Antrag der neuen Arten im `default`-Zweig von `applyAdministrationRequest`
  als `failed` mit Text: sichtbar, nicht still (§6).
- **Regel-Parsing in SQL** — die Funktionen prüfen keine Regelform; die
  Validierung liegt vollständig in Go
  ([`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 1).
- **Wirkung im Assembler** — `kern-rename`.
- **Betreiber-Handbuch** — der Slice liefert Betreiber-Oberfläche (zwei
  SQL-Funktionen), die noch nichts bewirkt; ihre Beschreibung ist aufgeschoben
  mit Adresse `slice-transformationen-betriebsdoku`, dessen §2 den
  aufgeschobenen Gegenstand vollständig nennt (SQL-Funktionen, Rollen, Wirkung,
  Abhilfe, Glossar, Änderungshistorie;
  `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, offen, 2×). Aussagen über
  eine Wirkung stehen erst, wenn sie belegt ist
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- **Eine Regel-Sicht (View)** —
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6: der Regelstand ist über `cdc.administration_request` lesbar,
  eine eigene Sicht ist nicht Teil.

## 2. Definition of Done

- [ ] Antragsweg im Schema: `cdc.set_transformation` und
      `cdc.remove_transformation` schreiben ausschließlich einen Antrag (Art,
      `rule_name`, `rule_spec` nur bei `set_transformation`, `column_name`
      leer) und senden `pg_notify`; `PUBLIC` hat kein `EXECUTE`, ein Login ohne
      `cdc_admin`-Mitgliedschaft scheitert mit „permission denied for
      function“; die `request_kind`-Menge trägt genau die Werte des
      Parent-Stands plus die zwei neuen; `rule_name` und `rule_spec` sind
      nullable; der zweite Schema-Träger
      `internal/adapters/driven/postgresstorage/schema.sql` trägt die zwei
      Spalten. *Zu belegen durch:* `make test-store`, `make schema-rollout`
      zweimal hintereinander (Exit 0) und
      `tools/harness/run-schema-rollout-guard-test.sh` (alle Läufe);
      `plan.yaml` und `down.sql` regeneriert, falls der Rollout sie verändert.
      Der **Alt-Tag-Lauf** desselben Skripts
      ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
      Entscheidung 7: das Schema des jüngsten `v*`-Tags per `git archive`
      ausrollen, danach den Arbeitsbaum — Exit 0 zweimal, der zuvor eingefügte
      Datenstand über `cdc.changes` lesbar; der Bericht nennt den Tag und die
      gedruckten Exit-Codes) trägt den Upgrade-Beleg für die zwei nullable
      Spalten (`jsonb` statt `text` ist über einen Alt-Bestand ungemessen) und
      die zwei Funktionen.
- [ ] Idempotenz-Guard und Rollen-Test: `knownForeignObjects` in
      `tools/schema/rolloutguard/guard.go` trägt beide Funktionen
      (Signatur-Schreibweise am realen `--plan-only`-Report gemessen, nicht
      angenommen — der Parameter `jsonb` ist neu gegenüber den bestehenden
      Funktionen), `guard_test.go` deckt sie;
      `internal/bootstrap/roles_rollout_file_internal_test.go` liest den
      Grant-Text beider Funktionen und färbt sich bei entferntem
      `GRANT`/fehlendem `REVOKE` rot. *Zu belegen durch:* `make test` (beide
      Tests) und je eine Mutation, die den Test rot färbt.
- [ ] Domäne und Store: `model.AdministrationRequestKind` trägt die zwei Arten;
      der Antrags-Konstruktor erzwingt `rule_name` für beide Arten und
      `rule_spec` für `set_transformation` (Invariante wie bei `column_name`);
      der `AdministrationRequestAdapter` liest die zwei Spalten (`ListPending`
      liefert sie); ein Antrag der neuen Arten wird in
      `applyAdministrationRequest` als `failed` mit Text vermerkt, und der
      Fehlertext nennt die tatsächliche Menge. *Zu belegen durch:* `make test`
      (Whitebox in `internal/bootstrap`, Domänen-Test) und `make test-store`
      (Round-Trip der zwei Spalten, skopierte Bereinigung —
      `BEO-PGC/test-isolation-geteilter-zustand`, offen, 1×).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: aufgeschoben mit Adresse — die Betreiber-Oberfläche entsteht
      hier, ihre Beschreibung trägt `slice-transformationen-betriebsdoku`
      (siehe §1); die Änderungshistorie des Handbuchs trägt dieser Slice
      **nicht**, weil er das Handbuch nicht berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/schema.yaml` | update | Spalten `rule_name text`, `rule_spec jsonb` (nullable) im neutralen Modell von `administration_request`; die `description` nennt sie. |
| `tools/schema/nacharbeit-administration.sql` | update | CHECK-Menge um zwei Werte (dort steht die Menge außerhalb des neutralen Modells), die zwei Funktionen, `REVOKE … FROM PUBLIC`, `GRANT EXECUTE … TO cdc_admin`, Kopfkommentar. |
| `internal/adapters/driven/postgresstorage/schema.sql` | update | zweiter Schema-Träger (Store-Tier) trägt die zwei Spalten. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | falls der Rollout sie verändert ([`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)). |
| `tools/schema/rolloutguard/guard.go` (+ `guard_test.go`) | update | Einträge beider Funktionen; der Kommentar, der die Objekte zählt, zählt neu; der Kommentar am Feld `allowDestructive` („nur bekannte Fremdobjekte blockieren“) nennt neben den bekannten Fremdobjekten die Klasse „View-Signatur“ (Adresse der Meldung aus `slice-backfill-change-origin`, `BEO-PGC/gemeldete-ungenauigkeit-ohne-traeger`). |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | Grant-Text der zwei Funktionen (Rolle `cdc_admin`, `PUBLIC` ohne Recht). Die Prüfung (4b) aus `slice-backfill-run-store` hält die Tabellen-Grants von `cdc_admin` auf `cdc.administration_request` (`SELECT`, `UPDATE`; kein `INSERT`, kein `DELETE`; kein Recht der beiden anderen Rollen); die zwei neuen Spalten fallen unter diese Tabellen-Grants und brauchen keinen eigenen Grant. |
| `internal/domain/model/administrationrequest.go` (+ Test) | update | zwei Antragsarten, Felder `RuleName`/`RuleSpec`, Konstruktor-Invarianten; der Doc-Kommentar zählt die Menge auf. |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` (+ Test) | update | Lesen der zwei Spalten, Abbildung der Antragsarten. |
| `internal/bootstrap/wiring.go` (+ Test) | update | nur der Fehlertext im `default`-Zweig von `applyAdministrationRequest` (nennt die geschlossene Menge). |
| `tools/harness/run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md` | prüfen / update | Zahl der Fremdobjekte und Beschreibung der Läufe; der Alt-Tag-Lauf ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7) wird ausgeführt **und geändert** (Übergabe aus `slice-backfill-sql-administration`): Lauf 5 vergleicht die `request_kind`-Menge des Alt-Bestands nach dem Upgrade exakt mit den **fünf** Werten (`backfill,disable,enable,exclude_column,include_column`, samt Fehlermeldung „fünf Antragsarten") und prüft `EXECUTE` allein für die Funktion `cdc.backfill_table`; der Slice zieht Menge und Meldung auf sieben Werte nach und ergänzt `EXECUTE` (nur `cdc_admin`, nicht `PUBLIC`) je neue Funktion — sonst färbt Lauf 5 mit der Erweiterung rot. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die geschlossene
`request_kind`-Menge (Parent-Stand plus zwei)“, „die Menge der Fremdobjekte
außerhalb des neutralen Modells“, „die Spaltenform von
`administration_request`“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufzählungen der Antragsarten | `grep -rn 'exclude_column' internal tools docs spec harness` | *(Implementer trägt ein)* | Fehlertext in `applyAdministrationRequest`, Doc-Kommentare, `schema.yaml`-Beschreibung nachziehen; Handbuch-Stellen an `betriebsdoku` melden; Spec-Stellen trägt `spec-nachzug` |
| Zahl der Fremdobjekte (Wortform des Parent-Stands, am Start gemessen) | `grep -rn 'sechs\|sieben\|acht' harness Makefile tools docs/user` | *(Implementer trägt ein)* | Guard-Kommentar, `harness/targets/schema-rollout.md`, die `harness/README.md`-Zeile `make example-demo-up`, Sensor-Dateien nachziehen; `Accepted` ADRs nicht ändern |
| Läufe des Guard-Tests | Lesen von `tools/harness/run-schema-rollout-guard-test.sh` und `harness/targets/schema-rollout.md` §Belege | *(Implementer trägt ein)* | Zahl und Beschreibung nachziehen, falls sich die Läufe ändern |
| Spaltenform von `administration_request` in weiteren Trägern | `grep -rn 'administration_request' internal tools --include=*.go --include=*.sql --include=*.yaml` | *(Implementer trägt ein)* | beide Schema-Träger und Test-Fixtures tragen dieselben Spalten |
| Fehlertext der geschlossenen Menge im `default`-Zweig | `grep -n 'geschlossene Menge' internal/bootstrap/wiring.go` | *(Implementer trägt ein)* | nennt die tatsächliche Menge; ein Test liest ihn |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-backfill-sql-administration` in
`done/` liegt (Kopplung K3 der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md) §5: beide Umsetzungen
berühren die geschlossene `request_kind`-Menge in
`tools/schema/nacharbeit-administration.sql`,
[`SPEC-019`](../../../../spec/pflichtenheft.md), `applyAdministrationRequest`
und den Idempotenz-Guard — dieser Slice erweitert die dort stehende Menge
additiv), `slice-transformationen-kern-rename` in `done/` liegt und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Schema, Guard,
  Rollen-Test, Domäne und Store nicht in einem Review tragen — der abtrennbare
  Teil ist der dritte Liefer-Punkt (Domäne und Store-Lesen) als eigener Slice
  `antragsweg-store` mit Start nach diesem; die Welle nimmt ihn dann in §4 auf.
- `in-progress` → `open` (blockiert): falls das neutrale Modell die zwei
  nullable Spalten samt CHECK-Änderung nicht konvergiert (dann Nacharbeit-SQL
  und ein zweiter Guard-Eintrag — Architect-Frage nach
  [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)) oder der
  Idempotenz-Guard die Funktion mit `jsonb`-Parameter nicht erkennt (Exit 8).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test`, `make test-store` und `make
schema-rollout` (zweimal) real grün + der Alt-Tag-Lauf von
`tools/harness/run-schema-rollout-guard-test.sh` real grün + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Idempotenz-Guard erkennt die Funktion mit `jsonb`-Parameter nicht** —
  ein zweiter Rollout bräche mit Exit 8. Die Signatur-Schreibweise der
  bestehenden Funktionen (`enable_table(in:text,in:text,in:text)` in
  `guard.go`) ist gelesen, die der neuen mit `jsonb` ist ungemessen. *Erwartet,
  zu belegen durch:* zweiter `make schema-rollout` und
  `run-schema-rollout-guard-test.sh`. **Ausgang:** *(bei Closure)*
- **Die zwei Spalten, besonders `rule_spec jsonb`, konvergieren nicht über
  einen Alt-Bestand.** Eine nullable `text`-Spalte an einer bestehenden Tabelle
  rollt über einen Alt-Bestand (gemessen im Architect-Verdikt
  `architect-verdict-schema-rollout-view-signatur`, Szenario 5); `jsonb` statt
  `text` ist nicht gemessen. *Erwartet, zu belegen durch:* der Alt-Tag-Lauf von
  `run-schema-rollout-guard-test.sh`. **Ausgang:** *(bei Closure)*
- **Zwei Schema-Träger driften** (`tools/schema/schema.yaml` und
  `postgresstorage/schema.sql` tragen dieselbe Spaltenform, jeder für einen
  anderen Tier). *Erwartet, zu belegen durch:* der Suchlauf-Eintrag und `make
  test-store`, das die Spalten real liest. **Ausgang:** *(bei Closure)*
- **Der Rollen-Test deckt den neuen Grant nicht.** Der Parser liest literale
  `GRANT … ;`-Anweisungen (Kopf von `roles_rollout_file_internal_test.go`);
  eine Funktion mit `jsonb`-Parameter ist gegen ihn ungemessen
  (`BEO-PGC/rollen-test-abdeckungsluecken`, offen, 2×). *Erwartet, zu belegen
  durch:* die Mutation aus DoD Punkt 2. **Ausgang:** *(bei Closure)*
- **Zwischenzustand: Funktion vorhanden, Wirkung fehlt.** Ein Betreiber, der
  `cdc.set_transformation` zwischen diesem Slice und `antragsweg-usecase`
  aufruft, erhält einen `failed`-Antrag statt einer Regel; der Fehlertext nennt
  die Ursache. Der Slice ist unveröffentlicht und die Folge-Slices folgen
  unmittelbar. *Erwartet, zu belegen durch:* der Test des `default`-Zweigs.
  **Ausgang:** *(bei Closure: entfallen mit der Closure von
  `slice-transformationen-antragsweg-usecase`)*
- **Das Handbuch nennt die Funktionen nicht**
  (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  verkörpert, 3×): die Adresse `slice-transformationen-betriebsdoku` muss den
  aufgeschobenen Gegenstand tragen. *Erwartet, zu belegen durch:* Review liest
  §2 von `betriebsdoku` gegen §1 dieses Slice. **Ausgang:** *(bei Closure)*
- **Idempotenz gegen die Fremdobjekte.** Jede neue `nacharbeit-*`-Funktion
  vergrößert die Menge, die ein zweiter Rollout als Blocker sieht
  (`BEO-PGC/schema-rollout-fremdobjekte`, verkörpert, 3×;
  `BEO-PGC/d-migrate-nacharbeit`, verkörpert, 6×). *Erwartet, zu belegen
  durch:* zweiter Rollout, `run-schema-rollout-guard-test.sh`. **Ausgang:**
  *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Schema-Skripte, Store-Adapter und Composition Root sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×) und
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×) — Guard-Eintrag und Rückführung
§4, `BEO-PGC/rollen-verdrahtung` (eingetreten, Grants der Rolle `cdc_admin`),
`BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×, einschlägig — Risiko §6),
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×, einschlägig — der
Handbuch-Aufschub),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, je 3×, Aufschub
mit Adresse), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×,
einschlägig — Store-Tests skopiert bereinigen),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
