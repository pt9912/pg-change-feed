# Slice slice-023: Rollen-spezifische DSN-Verdrahtung (Least-Privilege-Adoption)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die DoD dieses Slice (rollen-spezifische
Verdrahtung real gebunden, `make gates` grün) ist bereits die volle
Closure-Bedingung; es gibt kein *Mehr*, das eine Welle beobachten müsste
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).
Zeigt sich beim Zerlegen, dass mehr als drei Liefer-Punkte oder mehr als
zwei Schichten nötig sind, ist eine Welle die richtige Antwort — das
entscheidet die Planung beim Ausplanen, nicht dieser Entwurf.

**Bezug:** [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)
(Least-Privilege), [`LH-QA-SEC-002`](../../../../spec/lastenheft.md)
(Getrennte Berechtigbarkeit), [`LH-QA-SEC-003`](../../../../spec/lastenheft.md)
(Beschränkbarkeit von CDC-Datenzugriffen). Architect-Verdikt:
[`architect-review-welle-6.md`](../../adr/architect-review-welle-6.md)
Zug 2 — Ausgang `geplant` für `BEO-PGC/rollen-verdrahtung` (3×), dieser
Slice ist die zugewiesene Kennung. **[`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)**
(Accepted) — der Konfigurationsvertrag (drei Verbindungs-DSNs,
Rollen-Zuordnung je Aufrufer) ist damit entschieden; dieser Slice setzt
sie um, ohne selbst noch zu entscheiden. **[`ADR-0048`](../../adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md)**
(Accepted, Supersedes `ADR-0047` teilweise) — korrigiert den in
`ADR-0047` §Konsequenzen genannten Heartbeat-Grant-Text auf
`GRANT SELECT, INSERT, UPDATE`; die Rollen-Zuordnung selbst bleibt
unverändert.

**Berührte Spec-Stellen:** — (kein `SPEC-*`/`ARC-*`-Eintrag zur
DSN-Rollenbindung; die drei Rollen `cdc_capture`/`cdc_admin`/`cdc_reader`
sind DDL-Artefakt, kein Sicht-Gegenstand).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `internal/bootstrap/wiring.go` verdrahtet Store-, Aktivierungs-,
Stream- und die beiden CLI-Zugriffswege (`register-consumer`,
`acknowledge-consumer`) jeweils über die zur Aufgabe passende PostgreSQL-Rolle
(`cdc_capture` für den Stream-Adapter, `cdc_admin` für administrative/
schreibende Zugriffe, `cdc_reader` für reine Lesezugriffe) statt — wie heute —
durchgehend über eine gemeinsame Instanz-DSN. Die in `LH-QA-SEC-001`…003
geforderte Trennung existiert bereits als DDL-Artefakt und ist einzeln per
`SET ROLE` testbar (`BEO-PGC/rollen-verdrahtung/observation.md`); dieser Slice
macht sie in der tatsächlichen Laufzeit-Verdrahtung wirksam.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Rollen oder DDL-Änderungen an `cdc_capture`/`cdc_admin`/
  `cdc_reader`** — Bestand bleibt bewusst stehen: Die drei Rollen existieren
  bereits (`slice-011`) und werden von diesem Slice nur *verwendet*, nicht neu
  geschnitten.
- **Authentifizierung/Autorisierung des externen CLI-/SQL-/späteren
  API-Zugriffswegs** (wer darf welchen CLI-Unterbefehl aufrufen) — anderer
  Vorgang: Diese Frage betrifft den *externen* Aufrufer, nicht die *interne*
  DB-Rollenbindung, die dieser Slice herstellt.
- **Verdrahtung künftiger, noch nicht existierender Zugriffswege** (z. B.
  eine spätere HTTP-/gRPC-API nach `LH-FA-SST-006`) — Folge-Slice, sobald ein
  solcher Zugriffsweg tatsächlich gebaut wird; dieser Slice schließt nur die
  Lücke am *heutigen* Bestand (Capture, Store, beide CLI-Unterbefehle).

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

- [x] `LH-QA-SEC-001`/`002`/`003` erfüllt: Jeder Verdrahtungspfad in
      `internal/bootstrap/wiring.go` verwendet die zur Aufgabe passende Rolle,
      nicht mehr die gemeinsame Instanz-DSN — real getestet (z. B.
      Verbindungsaufbau mit `cdc_reader`-Rolle scheitert an einem
      schreibenden Aufruf). Beleg: `internal/bootstrap/roles_wiring_test.go`
      (`TestCdcReaderLoginConnectionRejectsWrite`,
      `TestCdcCaptureLoginConnectionRejectsAdminWrite`,
      `TestCdcAdminHeartbeatWriteRequiresGrant`), real gegen PostgreSQL
      dreimal in Folge grün (`make test-store`).
- [x] Konfigurationsvertrag erweitert gemäß [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md):
      `CDC_CAPTURE_DSN`, `CDC_ADMIN_DSN`, `CDC_READER_DSN` ersetzen
      `CDC_SOURCE_DSN` ersatzlos (Breaking Change, kein Fallback —
      Greenfield, kein veröffentlichtes Image) — `compose.yaml` und
      `docs/user/benutzerhandbuch.md` entsprechend nachgezogen.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für den erweiterten Konfigurationsvertrag
      (`docs/user/benutzerhandbuch.md`, `compose.yaml`) — bereits Teil des
      ersten DoD-Punkts, hier kein eigener Liefer-Punkt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      `BEO-PGC/rollen-verdrahtung` (3×, Ausgang bereits `geplant` →
      dieser Slice, siehe Kopf-Feld) bekommt bei Abschluss dieses Slice
      seinen finalen Ausgang *eingetreten*.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` | refactor | zentrale Verdrahtungsstelle (`ADR-0026`) — jede Adapter-Konstruktion bekommt ihre Rollen-DSN statt `cfg.DSN`, Zuordnung je Aufrufer nach `ADR-0047` |
| `internal/bootstrap/config.go` (oder Äquivalent — Implementer prüft den tatsächlichen Dateinamen) | update | `Config` trägt `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN` statt eines einzelnen `DSN`-Felds; `CDC_SOURCE_DSN` entfällt |
| `tools/schema/nacharbeit-roles.sql` (oder neue gleichartige Nacharbeit-Datei) | update | fehlenden Grant nachtragen: `GRANT INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin` (`ADR-0047` Kontext-Befund 3) — Lückenschließung, kein Neu-Zuschnitt der drei Rollen |
| `compose.yaml` | update | Env-Vertrag für die drei neuen DSN-Variablen |
| `docs/user/benutzerhandbuch.md` | update | Betreiber-Doku für den neuen Konfigurationsvertrag, inkl. Hinweis zum `REPLICATION`-Attribut auf der `CDC_CAPTURE_DSN`-Login-Identität (`ADR-0047` Kontext-Befund 2) |
| `internal/bootstrap/*_test.go` | update/neu | Rollenbindung real gegen PostgreSQL-Testcontainer geprüft (nicht nur behauptet) |

Der Implementer erweitert diese Liste im ersten Lauf um alle
Driven-/Driving-Adapter-Konstruktoren, die heute `cfg.DSN` nutzen (u. a.
`RegisterConsumer`, `AcknowledgeConsumer`, der Stream-Adapter, der
`EnableTableUseCase`-Pfad) — vollständig zu identifizieren, nicht zu raten.

**Plan-Nachzug (Implementer-Lauf, vor dem Sensor-Lauf eingetragen —
`modul-09-implementierung.md` §Plan-Nachzug im selben Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (kein eigenes `config.go`) | Klarstellung | `Config`/`ConfigFromEnv` liegen bereits in `wiring.go` (Composition Root, `ADR-0026`) — kein separates `config.go` im Bestand, deshalb kein neuer Dateiname |
| `cmd/pg-change-feed/main.go` | update | `--healthcheck`-Aufruf verwendet `cfg.ReaderDSN` statt der entfallenen `cfg.DSN` |
| `internal/bootstrap/wiring_test.go`, `register_test.go`, `acknowledge_test.go`, `welle6_endtoend_test.go` | update | `bootstrap.Config{DSN: …}` auf `AdminDSN`/drei-DSN-Form nachgezogen (Feldumbenennung) |
| `internal/bootstrap/roles_wiring_test.go` | neu | reale Verbindungsaufbau-Tests mit anmeldefähigen Test-Login-Identitäten je Rolle (`ADR-0047` §Fitness Function, beide Zeilen) — Login-Rolle ist Test-Fixture, kein Rollen-DDL |
| `harness/README.md` | update | ENV-Vertrag-Zeile für `make test-integration` auf die drei neuen DSN-Variablen umgestellt (öffentlicher Vertrag berührt, Workflow-Schritt 17) |
| `harness/image-hash.txt` | update | `make image` neu gelaufen (Build-Kontext geändert — Go-Quellen), neuer Lauf-Beleg-Digest |

**Korrektur ggü. Plan-Text (`tools/schema/nacharbeit-roles.sql`):** Der
Slice-Plan-Text (und `ADR-0047` Kontext-Befund 3/Konsequenzen) nennt
`GRANT INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin`. Real gegen
PostgreSQL getestet (`internal/bootstrap/roles_wiring_test.go`,
`TestCdcAdminHeartbeatWriteRequiresGrant`) zeigte: `INSERT, UPDATE` allein
lässt den Heartbeat-Schreibpfad weiterhin mit SQLSTATE 42501 scheitern —
PostgreSQL verlangt für den `ON CONFLICT (…) DO UPDATE`-Zweig zusätzlich
`SELECT` auf der Zieltabelle, auch wenn die `SET`-Klausel selbst keinen
bestehenden Spaltenwert liest. Der tatsächliche Grant lautet deshalb
`GRANT SELECT, INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin` — eine
Ergänzung um ein drittes Recht auf derselben einen Tabelle, keine Änderung
der Rollen-Zuordnung selbst (`ADR-0047`s Entscheidung bleibt unberührt;
diese Korrektur betrifft nur die SQL-Mechanik der bereits beschlossenen
Lückenschließung). Kandidat für einen Steering-Loop-Eintrag bei der
Planner-Closure (§7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einer
anderen Welle oder einem anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Ergeben sich beim
  Aufzählen der betroffenen Adapter-Konstruktoren mehr als drei
  Liefer-Punkte oder mehr als zwei Schichten (§3-Liste wächst über Bootstrap
  + mehrere Driven-Adapter hinaus).
- `in-progress` → `open` (blockiert — Carveout?): Der Konfigurationsvertrag
  (welche Rollen-DSNs, wie benannt) braucht eine eigene Architect-
  Entscheidung, bevor die Verdrahtung beginnen kann.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Review ohne offenes HIGH
**und** Closure-Notiz geschrieben (inkl. `BEO-PGC/rollen-verdrahtung` →
*eingetreten*).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Konfigurationsvertrag-Bruch: bestehende `CDC_SOURCE_DSN`-Nutzer (Betreiber,
  `compose.yaml`, Betreiberdoku) müssten auf mehrere Variablen umstellen —
  ein Migrationspfad/Kompatibilitäts-Fallback könnte nötig werden.
  **Architect-Entscheidung ([`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)):**
  kein Fallback — reiner Breaking Change. Begründung: Repo ist Greenfield,
  es gibt kein veröffentlichtes Image (`docs/user/benutzerhandbuch.md`
  §„Es gibt aktuell kein veröffentlichtes Container-Image"), also keine
  bestehenden externen Nutzer, deren Konfiguration bräche. **Ausgang:**
  formal bei Closure einzutragen, voraussichtlich *entfallen* (mit dieser
  Begründung) — kein Migrationspfad ist nötig, weil die Voraussetzung für
  das Risiko (ein bestehender externer Nutzer) nicht gegeben ist.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührte Sub-Area ist `*` (Default,
`PGC`, Greenfield laut `harness/conventions.md`) — die einzige deklarierte
Sub-Area dieses Repos, Schwelle also trivial erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** `BEO-PGC/rollen-verdrahtung`
(3×, Ausgang `geplant` → dieser Slice) ist der Anlass dieses Slice selbst.
Keine weiteren Treffer bei erster Durchsicht.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Refactor
auf bestehender Greenfield-Codebasis, keine Inventur-Diskrepanz zu
erwarten).
