# Slice slice-011: Sicherheit und Observability-Basis

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-3.

**Bezug:** [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…003, [`LH-FA-ADM-002`](../../../../spec/lastenheft.md)/003, [`LH-FA-SST-004`](../../../../spec/lastenheft.md), [`ADR-0043`](../../../../docs/plan/adr/README.md), [`ADR-0046`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`SPEC-009`](../../../../spec/pflichtenheft.md), [`ARC-005`](../../../../spec/architecture.md)
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Plan-Nachzug (Implementer-Lauf, vor dem Sensor-Lauf — §3/§1/§8 unten):**
`ARC-006`/`ARC-011` (Driven-Adapter, Telemetrie-Backend) ersetzt durch
`ARC-005` (Driving Adapters, SQL-Funktionen/Views) — die tatsächliche
Lieferung erweitert die bestehende SQL-Driving-Adapter-Fläche
([`ADR-0046`](../../../../docs/plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md))
um zwei Nacharbeit-Dateien; kein Driven-Adapter/Telemetrie-Backend-Code
entstand in diesem Lauf. Begründung und Abgrenzung: §1, §3, §8.

**Verantwortlich:** pt9912.
**Autor:** pt9912. **Datum:** 2026-09-09.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Sicherheit: Least-Privilege-Rollen (Capture/Reader/Admin getrennt je [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…003), Observability-Basis: Health-Endpoint ([`LH-FA-ADM-002`](../../../../spec/lastenheft.md)) und Metriken-Minimum ([`LH-FA-SST-004`](../../../../spec/lastenheft.md), [`SPEC-009`](../../../../spec/pflichtenheft.md)-Kern).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produktions-Deployment (HA/Kubernetes) — Out-of-Scope des MVP; der
  Slice liefert die Rollen und Endpoints, nicht die Betriebsumgebung.
- Vollständige Observability-Abdeckung — Basis hier (Health +
  Metriken-Minimum); die volle Metriken-Abdeckung folgt.

**Plan-Nachzug (Implementer-Lauf, vor dem Sensor-Lauf, Baseline-Regelwerk
`modul-09-implementierung.md` §Plan-Nachzug):** Zwei Punkte gehen über die
ursprüngliche Abgrenzung hinaus und werden hier nachgetragen — Klasse
**anderer Vorgang** bzw. **Schicht-Abgrenzung**, keine Kennung vorhanden,
weil keine der beiden Erweiterungen bereits als eigener Slice existiert:

- **Health-Endpoint ([`LH-FA-ADM-002`](../../../../spec/lastenheft.md), [`LH-QA-OPS-002`](../../../../spec/lastenheft.md)) — nicht realisiert.**
  **Planner-Korrektur nach Architect-Verdikt**
  ([`docs/plan/adr/architect-review-slice-011.md`](../../adr/architect-review-slice-011.md)):
  Die ursprüngliche Implementer-Begründung (neue ADR nötig) trug nicht —
  der Architect bestätigte den Reviewer (review-slice-011 F-1): ein
  Heartbeat-Pattern (Prozess schreibt periodisch in eine
  `cdc.process_heartbeat`-Tabelle, eine vierte SQL-View liest sie) deckt
  [`LH-FA-ADM-002`](../../../../spec/lastenheft.md)/[`LH-QA-OPS-002`](../../../../spec/lastenheft.md)
  vollständig innerhalb der bestehenden Entscheidungslage
  ([`ADR-0024`](../../../../docs/plan/adr/0024-observability-ausserhalb-der-domain.md),
  [`ADR-0027`](../../../../docs/plan/adr/0027-capture-application-service.md),
  [`ADR-0046`](../../../../docs/plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md))
  — **kein Folge-ADR nötig**, [`ADR-0020`](../../../../docs/plan/adr/0020-http-grpc-optional.md)s
  HTTP/gRPC-Sperre ist nicht berührt. Der tragende Ausschlussgrund ist
  **Schicht-Abgrenzung**, nicht **anderer Vorgang**: Der periodische
  Schreib-Zug berührt Application- und Bootstrap-Schicht zugleich
  (Timer im Capture-Prozess) — genau dieselbe Grenze, die dieser Slice
  bereits für `wiring.go` unten zieht; ihn hier mitzuliefern hätte den
  bereits vorab benannten „zu groß"-Rückführungs-Trigger (§4) ausgelöst.
  **Klasse: Folge-Slice** (Kennung folgt bei der nächsten Slice-Eröffnung
  — `Bezug: LH-FA-ADM-002, LH-QA-OPS-002`, Kopf-Referenz
  `ADR-0024`/`ADR-0027`/`ADR-0046`, Beleg-Zeiger auf das Architect-
  Verdikt oben). Siehe §5 (Closure-Trigger geteilt) und §6 (Risiko).
- **`internal/bootstrap/wiring.go` bleibt unverändert** — die drei Rollen
  (`cdc_capture`/`cdc_admin`/`cdc_reader`) entstehen als DDL/Compose-Artefakt
  (Rollen existieren, sind über `SET ROLE` einzeln testbar,
  `internal/adapters/driven/postgresstorage/roles_test.go`), aber die
  Verdrahtung trägt weiterhin eine gemeinsame Instanz-DSN (`Config.DSN`)
  für Store-, Aktivierungs- und Stream-Verbindung. Rollen-spezifische DSNs
  je Adapter zu verdrahten berührt Bootstrap **und** alle drei
  Driven-Adapter-Konstruktoren zugleich — **Schicht-Abgrenzung**: dieser
  Slice hält sich auf die DB-Schicht (DDL/Compose), die Anwendungs-/
  Verdrahtungsschicht bleibt unberührt; sonst wäre der bereits im Plan
  vorab benannte Rückführungs-Trigger „zu groß" (§4) eingetreten. Siehe §6.

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Least-Privilege-Rollen in Compose/DDL getrennt — Teil-Beleg zu
      [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…003.
- [x] Metriken-Minimum (`cdc.metrics`-View) — Teil-Beleg zu
      [`LH-FA-SST-004`](../../../../spec/lastenheft.md). Health-Endpoint
      ist Folge-Slice-Arbeit (§1, §5 — Architect-Verdikt).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `harness/README.md` (Sensors-Tabelle), falls
      berührt — geprüft: `make test-store`-Bindung trägt bereits den
      Rollout-Schritt (seit slice-009); kein neuer Vertrag, Item
      entfällt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| ~~`internal/adapters/driving/` (Health/Metrics)~~ — **Plan-Nachzug: nicht realisiert** | — | ersetzt durch die beiden Zeilen unten; Health-Endpoint fehlt eine ADR (§1 Plan-Nachzug), kein neuer Go-Driving-Adapter entstand |
| `tools/schema/nacharbeit-roles.sql` — **Plan-Nachzug: neue Datei statt „Rollen-DDL"-Sammelzeile** | neu | drei Least-Privilege-Rollen (`cdc_capture`/`cdc_admin`/`cdc_reader`) je [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…003, angewandt über `make schema-rollout` (Nacharbeit-Schritt wie `nacharbeit-views.sql`, [`ADR-0043`](../../../../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)) |
| `tools/schema/nacharbeit-observability.sql` — **Plan-Nachzug: neue Datei** | neu | Metriken-Minimum-View `cdc.metrics` je [`LH-FA-SST-004`](../../../../spec/lastenheft.md)/[`SPEC-009`](../../../../spec/pflichtenheft.md)-Teilmenge, als vierte Lese-View auf der bestehenden SQL-Driving-Adapter-Fläche ([`ARC-005`](../../../../spec/architecture.md), [`ADR-0046`](../../../../docs/plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)) |
| `Makefile` (`schema-rollout`) — **Plan-Nachzug: neue Zeile** | update | zwei zusätzliche Nacharbeit-Schritte für die beiden Dateien oben |
| `internal/adapters/driven/postgresstorage/roles_test.go` — **Plan-Nachzug: neue Datei** | neu | Berechtigungsprüfung (`SET ROLE`) gegen reale PostgreSQL — Teil-Beleg [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…003, `make test-store` |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-010 liegt in `done/`; kein anderes Slice in `in-progress/`
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Sicherheit (Rollen) plus Observability (Health + Metriken) sprengen
  drei Liefer-Punkte — Rückzug mit Zerlegung (Rollen-DDL als eigener
  Slice).
- `in-progress` → `open` (blockiert — Carveout?): Compose/DDL-Umgebung weicht vom Least-Privilege-Vertrag ab und ist
  vor dem Fix-Zug nicht herstellbar — Blocker, Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss
      ohne offenes HIGH-Finding; Least-Privilege-Rollen und
      Metriken-Minimum-View am realen Adapter belegt (`make test-store`).
      **Planner-Korrektur nach Architect-Verdikt**
      ([`docs/plan/adr/architect-review-slice-011.md`](../../adr/architect-review-slice-011.md),
      review-slice-011 F-2): Der ursprüngliche Trigger verlangte „Health-
      und Metriken-Endpoint am verdrahteten System" — Health-Endpoint ist
      als Folge-Slice-Arbeit ausgegliedert (§1), dieser Slice schließt auf
      Sicherheit + Metriken-Minimum allein; kein Carveout, da dies eine
      Scope-Reduktion des Triggers ist, kein rotes Gate.
- Lerneintrag §7: geschärfte Regel oder benannte Spec-Lücke —
  ohne ihn bleibt der Slice abgelegt, nicht fertig.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Produktions-Deployment (HA/Kubernetes) — **Ausgang:** entfallen; das
  ist eine Out-of-Scope-Deklaration aus §1, keine eigene Risiko-Klasse
  (redundante Übernahme aus der welle-3-Fill-Routine, ohne Fehlwirkung
  — nichts zu tun).
- Vollständige Observability-Abdeckung — **Ausgang:** entfallen;
  ebenfalls redundante §1-Übernahme, kein eigenständiges Risiko. Der
  konkrete Rest ([`LH-QA-OPS-003`](../../../../spec/lastenheft.md))
  ist kein akutes Risiko dieses Slice, sondern erwarteter MVP-Rand.
- **Health-Endpoint** ([`LH-FA-ADM-002`](../../../../spec/lastenheft.md),
  [`LH-QA-OPS-002`](../../../../spec/lastenheft.md)) — **Ausgang:**
  weiter offen → `BEO-PGC/health-endpoint-heartbeat` im Register
  (Eintrag angelegt, Beleg `evidence/slice-011.md`). Architect-Verdikt
  klärte: kein Folge-ADR nötig, Folge-Slice-Arbeit (§1, §5).
- **Rollen-Verdrahtung** (`internal/bootstrap/wiring.go` trägt weiter
  eine gemeinsame Instanz-DSN) — **Ausgang:** weiter offen →
  `BEO-PGC/rollen-verdrahtung` im Register (Eintrag angelegt, Beleg
  `evidence/slice-011.md`).
- **Beobachtungs-Register gesichtet (§8):** `BEO-PGC/lese-doppelquelle`
  trifft `cdc.metrics` strukturell mit (`evidence/slice-011.md`
  ergänzt, Zähler 2×, unter der Schwelle).

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

- **Was hat funktioniert:** Der Implementer eskalierte die
  Health-Endpoint-Frage sauber statt selbst zu entscheiden (Modul 8
  §Konflikt-Pfad) — kein stiller ADR-Widerspruch, sondern ein
  benannter Plan-Nachzug mit Reviewer-/Architect-Hinweis. Der
  Review deckte trotzdem auf, dass der Rückzug zu breit begründet
  war (F-1) — die Rollen-Trennung fing zwei Findings echt ab: F-1
  (Architektur-Rückzug zu breit) und F-3 (ungetesteter
  Least-Privilege-Kernpfad, real als SQLSTATE 42501 nachgewiesen).
- **Was ging anders als geplant:** Health-Endpoint und Metriken-Minimum
  wurden getrennt statt gemeinsam geliefert — Architect-Verdikt
  ([`docs/plan/adr/architect-review-slice-011.md`](../../adr/architect-review-slice-011.md))
  klärte, dass Health-per-Heartbeat keine neue ADR braucht, aber eine
  eigene Schicht-Abgrenzung ist (Application-/Bootstrap-Zug), die
  dieser Slice bewusst nicht mitliefert. `cdc_admin`s
  `GRANT CREATE ON DATABASE` reichte entgegen der ersten Annahme
  nicht für `ALTER PUBLICATION … ADD TABLE` — PostgreSQL verlangt
  Tabellen-Ownership; als operative Vorbedingung dokumentiert statt
  über-großzügig nachgegeben.
- **Steering-Loop-Eintrag:** nichts verkörpert — der Normalfall in
  diesem Slice (die Verkörperung lief bereits über slice-009/010s
  Steering-Loop-Einträge). Die drei neuen Register-Einträge (unten)
  bleiben unter der 3×-Schwelle.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/health-endpoint-heartbeat/`
  neu angelegt (Beleg `evidence/slice-011.md`, 1×); `BEO-PGC/rollen-verdrahtung/`
  neu angelegt (Beleg `evidence/slice-011.md`, 1×); `BEO-PGC/lese-doppelquelle/`
  ergänzt (`evidence/slice-011.md`, Zähler jetzt 2×).
- **Folge-Slices:** keine mit Kennung — Health-per-Heartbeat und
  Rollen-Verdrahtung bleiben Register-Beobachtungen bis zum nächsten
  Schneiden.
- **Risiken aus §6:** zwei entfallen (redundante §1-Übernahme aus der
  welle-3-Fill-Routine), zwei weiter offen → Register (siehe oben).
- **Drei Paarungen:** entfällt hier — das Repo arbeitet mit Wellen;
  die Welle-3-Closure prüft Anker · Folge-Slice · Register auch für
  diesen Slice.

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

**Vorgelagert — Sub-Area-Wahl prüfen (nachgetragen im Implementer-Lauf):**
`harness/conventions.md` deklariert genau eine Sub-Area (`*`, Kürzel `PGC`,
Modus Greenfield) für das gesamte Repo — keine feinere Deklaration
existiert, gegen die dieser Slice zu grob wäre. Die berührte Fläche
(DDL/Compose-Rollen, SQL-Lese-View) fällt vollständig unter dieses eine
`PGC`; keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten (nachgetragen im
Implementer-Lauf):** Register `docs/plan/planning/observations/BEO-PGC/`
durchgegangen (sieben Einträge). Treffer mit Zähler-Stand:
`lese-doppelquelle` (1×, `evidence/slice-010.md`) — die neue
`cdc.metrics`-View liest über `cdc.consumer_status` dieselbe
Rückstands-Semantik wie die bestehenden Views, verstärkt also dieselbe
Beobachtungsklasse (kein Sensor hält SQL-View-Semantik gegen
Go-Use-Case-Semantik), siehe §6. `d-migrate-nacharbeit` (2×) — nicht
unmittelbar getroffen: die beiden neuen Nacharbeit-Dateien folgen
demselben etablierten Muster (`nacharbeit-views.sql`), lösen aber keinen
neuen `raw-sql-text-drift`-Fall aus (getestet, `make test-store` grün).
Übrige fünf Einträge (`a-check-null-abdeckung`, `adapter-fehler-ausgang`,
`plan-nachzug`, `plan-vorlagen-defekt`, `walsender-wirksamkeit`) ohne
Bezug zu diesem Slice.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

*Reiner GF-Hinweis genügt (siehe oben); kein Sub-Area-Block.*
