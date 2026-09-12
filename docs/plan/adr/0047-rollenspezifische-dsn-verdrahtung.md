# ADR-0047: Rollen-spezifische DSN-Verdrahtung (eine Verbindungs-DSN je PostgreSQL-Rolle)

**Status:** Accepted

**Datum:** 2026-09-12

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-QA-SEC-001`](../../../spec/lastenheft.md) (Least-Privilege),
[`LH-QA-SEC-002`](../../../spec/lastenheft.md) (Getrennte
Berechtigbarkeit), [`LH-QA-SEC-003`](../../../spec/lastenheft.md)
(Beschränkbarkeit von CDC-Datenzugriffen),
[ADR-0026](0026-composition-root.md) (Composition Root),
[ADR-0043](0043-schemamigrationen-mit-d-migrate.md) (Nacharbeit-Ausweichform
für Rollen-DDL), Architect-Verdikt
[`architect-review-welle-6.md`](architect-review-welle-6.md) Zug 2
(`BEO-PGC/rollen-verdrahtung`, 3×, Ausgang `geplant` →
`slice-023`) <!-- d-check:status-provenance -->

**Schärft:** [`ARC-007`](../../../spec/architecture.md) — Bootstrap/
Composition Root (bereits von [ADR-0026](0026-composition-root.md)
geschärft; diese ADR schärft dieselbe Sicht-Stelle um die konkrete
Verbindungs-/Rollen-Zuordnung der Verdrahtung).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`internal/bootstrap/wiring.go` verdrahtet heute **jeden** Aufrufer — den
Store-Adapter (`postgresstorage.New`), die Tabellen-Aktivierung
(`postgresstorage.NewTableActivation`), den Heartbeat-Adapter
(`postgresstorage.NewHeartbeat`), den Replication-Stream
(`receive.NewStream`), den ACK-Adapter (teilt die Stream-Verbindung) sowie
die beiden CLI-Sondermodi `RegisterConsumer`/`AcknowledgeConsumer`
(`postgresstorage.NewConsumerState`) — über dieselbe gemeinsame `cfg.DSN`
(`CDC_SOURCE_DSN`). Drei PostgreSQL-Rollen mit unterschiedlichem
Least-Privilege-Zuschnitt existieren bereits als DDL seit `slice-011` <!-- d-check:status-provenance -->
(`tools/schema/nacharbeit-roles.sql`):

| Rolle | Grants (Ist-Stand) |
|---|---|
| `cdc_capture` | `SELECT` auf `cdc.source_table`/`cdc.schema_version`, `INSERT` auf `cdc.transaction`/`cdc.change`, Attribut `REPLICATION` |
| `cdc_admin` | `SELECT`/`INSERT`/`UPDATE`/`DELETE` auf `cdc.source_table`/`cdc.schema_version`/`cdc.consumer`/`cdc.consumer_position`, `SELECT` auf `cdc.transaction`/`cdc.change`, `CREATE` auf der Datenbank |
| `cdc_reader` | `SELECT` auf `cdc.active_tables`, `cdc.consumer_status`, `cdc.changes`, `cdc.metrics` (`nacharbeit-observability.sql`), `cdc.heartbeat` (`nacharbeit-heartbeat.sql`) |

Diese Trennung wird von der Verdrahtung nicht genutzt — jeder Aufrufer
bekäme dieselben (impliziten) Rechte, weil es nur eine Konfigurationsquelle
gibt. `BEO-PGC/rollen-verdrahtung` hat mit `slice-011`, `slice-021`, <!-- d-check:status-provenance -->
`slice-022` dreimal denselben Befund reproduziert (physische <!-- d-check:status-provenance -->
Verdrahtungslücke, kein Workflow-Defizit — Architect-Verdikt
`architect-review-welle-6.md` Zug 2) und ist damit `geplant → slice-023`. <!-- d-check:status-provenance -->
`slice-023` benennt die Entscheidung selbst als Blocker für <!-- d-check:status-provenance -->
`in-progress → open` (§4): „Der Konfigurationsvertrag … braucht eine
eigene Architect-Entscheidung, bevor die Verdrahtung beginnen kann." Diese
ADR ist diese Entscheidung.

**Drei zusätzliche technische Befunde, beim Nachvollziehen der
Grant-Zuordnung entdeckt — keiner davon ändert eine bestehende
`Accepted`-Entscheidung, alle drei begrenzen die Optionen unten:**

1. **Alle drei Rollen sind `NOLOGIN`** (`nacharbeit-roles.sql`, bewusste
   Grenze: „Gruppenrollen, keine Anmelde-Identität mit Passwort"). Eine DSN
   kann nur gegen eine anmeldefähige Identität aufgebaut werden — jede
   DSN-basierte Option setzt voraus, dass der Betreiber pro Rolle (oder für
   alle) eine Login-Identität bereitstellt, die Mitglied der jeweiligen
   Gruppenrolle ist (Muster bereits dokumentiert:
   `docs/user/benutzerhandbuch.md` §2 „`GRANT cdc_reader TO ihr_login;`").
   Diese ADR entscheidet **nicht**, wie der Betreiber diese Login-Identität
   anlegt — nur, wie viele Verbindungs-Endpunkte die Verdrahtung braucht.
2. **PostgreSQL vererbt Rollen-*Attribute*** (`LOGIN`, `REPLICATION`,
   `SUPERUSER`, `CREATEDB`, `CREATEROLE`, `BYPASSRLS`) **nicht über
   Mitgliedschaft** — nur Objekt-Privilegien (`GRANT SELECT/INSERT/…`) tun
   das. Das `REPLICATION`-Attribut liegt direkt auf `cdc_capture`
   (`CREATE ROLE cdc_capture NOLOGIN REPLICATION`); eine Login-Rolle, die
   nur `GRANT cdc_capture TO <login>` erhält, bekommt dadurch **nicht**
   automatisch die Fähigkeit, eine Replication-Verbindung aufzubauen — das
   Attribut muss direkt auf der Login-Rolle gesetzt werden
   (`ALTER ROLE <login> REPLICATION;`). Das ist dieselbe Klasse offener
   Betriebs-Vorbedingung, die `nacharbeit-roles.sql` für die
   Tabellen-Eigentümerschaft von `cdc_admin` bereits dokumentiert (Zeile
   72ff., real getestet in `roles_test.go`).
3. **`cdc.process_heartbeat` trägt heute keinen Grant an irgendeine der
   drei Rollen.** Der Heartbeat-Adapter (`postgresstorage.NewHeartbeat`,
   `Beat`/`Fault`) schreibt per `INSERT … ON CONFLICT` gegen diese Tabelle
   (`internal/adapters/driven/postgresstorage/queries/queries.go`); weder
   `cdc_capture` noch `cdc_admin` noch `cdc_reader` hat dafür ein
   `INSERT`/`UPDATE`-Recht. Ohne einen ergänzenden Grant kann **keine** der
   drei Rollen den Heartbeat-Schreibpfad übernehmen — das ist eine
   vorbestehende Lücke in `nacharbeit-roles.sql`/`nacharbeit-heartbeat.sql`
   (`slice-011`), keine neue Anforderung dieser ADR, aber eine <!-- d-check:status-provenance -->
   Voraussetzung, um die hier getroffene Zuordnung real umzusetzen.

## Entscheidung

Wir wählen **A — eine eigene Verbindungs-DSN je Rolle**, drei insgesamt:
`CDC_CAPTURE_DSN`, `CDC_ADMIN_DSN`, `CDC_READER_DSN`. `CDC_SOURCE_DSN`
entfällt ersatzlos (Breaking Change ohne Fallback — Begründung siehe
Konsequenzen). Jeder bestehende Aufrufer in `internal/bootstrap/wiring.go`
bindet an die zur Aufgabe passende Rolle:

| Aufrufer | Env-Var | Rolle | Begründung |
|---|---|---|---|
| Store-Adapter (`postgresstorage.New`, Persist Transaction/Change) | `CDC_CAPTURE_DSN` | `cdc_capture` | deckt sich exakt mit den Ist-Grants: `SELECT` auf `source_table`/`schema_version`, `INSERT` auf `transaction`/`change` |
| Replication-Stream (`receive.NewStream`) | `CDC_CAPTURE_DSN` | `cdc_capture` | `REPLICATION`-Attribut liegt direkt auf `cdc_capture` (siehe Kontext, Befund 2) |
| ACK-Adapter (`postgresack.New(stream.Conn())`) | — (teilt die Stream-Verbindung) | `cdc_capture` | dieselbe physische Verbindung wie der Stream (`ADR-0007` Option C) — kein eigenes DSN |
| Tabellen-Aktivierung (`postgresstorage.NewTableActivation`, `EnableTableUseCase`) | `CDC_ADMIN_DSN` | `cdc_admin` | DML auf `source_table`/`schema_version` + `CREATE`-Recht auf der Datenbank für `CREATE`/`ALTER PUBLICATION` |
| Heartbeat-Adapter (`postgresstorage.NewHeartbeat`, `Beat`/`Fault`) | `CDC_ADMIN_DSN` | `cdc_admin` | operatives Verwaltungs-/Statusschreiben (Betriebszustand der Quelle), kein Erfassungs-Nutzlast-Pfad — bewusst getrennt von `cdc_capture`, das strikt auf den Datenpfad `transaction`/`change` beschränkt bleibt; setzt den ergänzenden Grant aus Kontext-Befund 3 voraus |
| `RegisterConsumer`/`AcknowledgeConsumer` (CLI-Sondermodi, `postgresstorage.NewConsumerState`) | `CDC_ADMIN_DSN` | `cdc_admin` | DML auf `consumer`/`consumer_position` — löst zugleich `docs/user/benutzerhandbuch.md`s bisherigen Hinweis „eine gesonderte Rolle für diesen Zugriffsweg ist nicht verdrahtet" auf |
| `Healthcheck` (Lese-Zugriff auf `cdc.heartbeat`) | `CDC_READER_DSN` | `cdc_reader` | `SELECT`-Grant auf `cdc.heartbeat` liegt bereits bei `cdc_reader` |
| künftige, heute noch nicht existierende Lesezugriffe | `CDC_READER_DSN` | `cdc_reader` | Default für neue reine Lesepfade, solange kein Schreibbedarf entsteht |

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — **drei getrennte DSNs, eine je Rolle (gewählt)** | echte, unabhängig widerrufbare Zugangsdaten je Rolle — ein kompromittierter Capture-Prozess kennt kein Passwort mit Admin-Rechten, kein Code-Pfad kann versehentlich ohne `SET ROLE`-Disziplin eskalieren; deckt `LH-QA-SEC-002`/`003` **strukturell** statt nur prozedural ab | Betreiber muss drei Zugangsdaten/Login-Identitäten bereitstellen statt einer; `compose.yaml`/Secrets-Handling wächst um zwei Variablen |
| B — eine Basis-Login-DSN + `SET ROLE <rolle>` je Adapter zur Laufzeit | ein einziges Zugangsdatum zu verwalten; PostgreSQL unterstützt `SET ROLE` nach dem Verbindungsaufbau nativ | das Basis-Login muss Mitglied **aller drei** Gruppenrollen sein (oder `SUPERUSER`), um `SET ROLE` in jede von ihnen zu erlauben — das eine Passwort trägt damit faktisch das Vereinigungs-Privileg aller drei Rollen; vergisst ein Adapter-Konstruktor den `SET ROLE`-Aufruf, läuft er unbemerkt mit den vollen Rechten des Basis-Logins statt der beabsichtigten Einschränkung (Fehler durch Unterlassung, nicht durch Handlung) — das unterläuft `LH-QA-SEC-003` („ohne die vergebene Berechtigung schlägt der Zugriff fehl") strukturell, nicht nur im Testfall; zusätzlich löst `SET ROLE` das `REPLICATION`-Attribut-Problem aus Kontext-Befund 2 nicht: auch nach `SET ROLE cdc_capture` bleibt die tatsächliche Session-Berechtigung für `START_REPLICATION` an das Attribut der ursprünglich angemeldeten Rolle gebunden, nicht an die per `SET ROLE` aktive Rolle |
| C — Status quo beibehalten (eine gemeinsame `CDC_SOURCE_DSN` für alle Aufrufer) | keine Änderung, kein Migrationsaufwand | verfehlt `LH-QA-SEC-001`/`002` vollständig — die drei bereits vorhandenen Rollen blieben ungenutzte DDL, „Berechtigungsprüfung: Lese-Rolle ohne Admin-Rolle kann lesen, nicht administrieren" (`LH-QA-SEC-002`-Messmethode) ist mit einer einzigen DSN nicht herstellbar; das ist genau der 3×-Befund, den `BEO-PGC/rollen-verdrahtung` markiert |

## Konsequenzen

- Positiv: `LH-QA-SEC-001`/`002`/`003` werden **strukturell**, nicht nur
  prozedural erfüllt — die Trennung besteht unabhängig davon, ob ein
  künftiger Adapter-Konstruktor diszipliniert bleibt.
- Positiv: Der Blast-Radius eines kompromittierten Capture-Prozesses
  beschränkt sich auf `INSERT` gegen `transaction`/`change` (plus Lesen der
  Bindungstabellen) — kein Zugriff auf `consumer`/`consumer_position`, kein
  `CREATE`-Recht.
- Negativ: Betreiber verwaltet drei Zugangsdaten/Login-Identitäten statt
  einer — Mehraufwand in Secrets-Handling und `compose.yaml`.
- Negativ / Folgepflicht (DDL-Ergänzung, **kein** Neu-Zuschnitt der drei
  Rollen — Bestand aus `slice-011` bleibt unverändert stehen, es wird nur <!-- d-check:status-provenance -->
  der eine fehlende Grant nachgetragen): `tools/schema/nacharbeit-roles.sql`
  (oder eine neue, gleichartige Nacharbeit-Datei) ergänzt
  `GRANT INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin` — ohne
  diesen Grant kann der Heartbeat-Adapter unter keiner der drei Rollen
  schreiben (Kontext-Befund 3). Dieser Grant ist eine reine Lückenschließung
  am bestehenden Rollen-DDL, keine Änderung der bereits erteilten
  Privilegien der drei Rollen — `slice-023`s Ausschluss „Neue Rollen oder <!-- d-check:status-provenance -->
  DDL-Änderungen … Bestand bleibt bewusst stehen" bezieht sich auf das
  **Zuschneiden** der Privilegien, nicht auf das Schließen einer Lücke, ohne
  die diese Entscheidung nicht real umsetzbar wäre.
- Folgepflicht: `compose.yaml` (Env-Vertrag), `docs/user/benutzerhandbuch.md`
  §2 „Zugriff und Rollen", §3 „Schneller Einstieg", §4 „Consumer
  registrieren"/„Position bestätigen"/„Betriebsstatus prüfen", §5
  „Konfiguration" werden auf die drei neuen Variablen umgestellt;
  `docs/user/benutzerhandbuch.md` ergänzt den Betriebs-Hinweis aus
  Kontext-Befund 2 (`REPLICATION`-Attribut direkt auf die
  `CDC_CAPTURE_DSN`-Login-Identität setzen).
- Folgepflicht: `internal/bootstrap/config.go`/`Config` trägt künftig drei
  DSN-Felder statt eines (Implementer-Arbeit, `slice-023`). <!-- d-check:status-provenance -->

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `internal/bootstrap/*_test.go` gegen PostgreSQL-Testcontainer | Ein Verbindungsaufbau mit der `cdc_reader`-Login-Identität scheitert an einem schreibenden Aufruf (z. B. `INSERT` auf `cdc.change`); ein Verbindungsaufbau mit der `cdc_capture`-Login-Identität scheitert an einem administrativen Aufruf (z. B. `INSERT` auf `cdc.consumer`) | `make test-store` |
| `internal/bootstrap/*_test.go` gegen PostgreSQL-Testcontainer (Replikation) | Der Heartbeat-Schreibpfad gelingt nur, nachdem der ergänzende `GRANT` auf `cdc.process_heartbeat` ausgerollt ist — real getestet, nicht angenommen | `make test-replication`/`make test-store` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Re-Evaluierung fällig, wenn eines eintritt: (a) eine vierte CDC-Rolle mit
eigenem Aufgabenschnitt wird eingeführt (Frage: bleibt „eine DSN je Rolle"
tragfähig, oder wird die Variablen-Zahl unhandlich?); (b) ein
Secrets-Management-System (z. B. dynamische Vault-Datenbank-Leases) verlangt
eine andere Bereitstellungsform als statische DSN-Strings in
Umgebungsvariablen — dann ist zu prüfen, ob Option A in dieser Form noch
trägt oder durch eine Lease-basierte Variante ersetzt wird. Sonst
**permanent**.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-12 | Accepted — Anlass: `BEO-PGC/rollen-verdrahtung` 3× (`architect-review-welle-6.md` Zug 2), Ausgang `geplant` → `slice-023`; diese ADR ist die dort angeforderte Architect-Entscheidung vor Implementierungsbeginn | `slice-023` (in `in-progress/`) <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0047` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
