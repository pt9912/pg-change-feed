# ADR-0048: Heartbeat-Grant-Korrektur — `SELECT` ergänzt den `ON CONFLICT`-Zweig

**Status:** Accepted — Supersedes [`ADR-0047`](0047-rollenspezifische-dsn-verdrahtung.md) (nur Konsequenzen-Bullet zum Heartbeat-Grant; die Rollen-Zuordnung und der Drei-DSN-Vertrag aus `ADR-0047` bleiben unverändert)

**Datum:** 2026-09-12

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-QA-SEC-001`](../../../spec/lastenheft.md) (Least-Privilege),
[`LH-QA-SEC-002`](../../../spec/lastenheft.md) (Getrennte
Berechtigbarkeit), [`LH-QA-SEC-003`](../../../spec/lastenheft.md)
(Beschränkbarkeit von CDC-Datenzugriffen), [ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md)
(Rollen-spezifische DSN-Verdrahtung — die hier korrigierte Folgepflicht
steht dort im Abschnitt „Konsequenzen"), [ADR-0043](0043-schemamigrationen-mit-d-migrate.md)
(Nacharbeit-Ausweichform für Rollen-DDL), Commit `a32a2c9`
(`internal/bootstrap/roles_wiring_test.go`, real gegen PostgreSQL-Testcontainer,
dreimal `make test-store` grün)

**Schärft:** [`ARC-007`](../../../spec/architecture.md) — Bootstrap/
Composition Root (dieselbe Sicht-Stelle, die bereits [ADR-0026](0026-composition-root.md)
und [ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md) schärfen; diese
ADR ändert an der Sicht-Aussage selbst nichts — sie korrigiert nur den
exakten Grant-Text einer dort schon vorgesehenen DDL-Nacharbeit).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0047`](0047-rollenspezifische-dsn-verdrahtung.md) benannte in ihrem
Kontext-Befund 3 eine vorbestehende Lücke: `cdc.process_heartbeat` trug
keinen Grant an irgendeine der drei CDC-Rollen, obwohl der
Heartbeat-Adapter (`postgresstorage.NewHeartbeat`, `Beat`/`Fault`) per
`INSERT … ON CONFLICT … DO UPDATE` gegen diese Tabelle schreibt. Als
Folgepflicht nannte ihr Abschnitt „Konsequenzen" den exakten SQL-Text:

```sql
GRANT INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin;
```

Beim realen Ausrollen dieses Grants und dem zugehörigen Fitness-Function-Test
(`internal/bootstrap/roles_wiring_test.go`, Commit `a32a2c9`, dreimal in
Folge `make test-store` grün gegen den PostgreSQL-Testcontainer) zeigte
sich: Dieser exakte Grant **reicht nicht**. PostgreSQL benötigt für den
`INSERT … ON CONFLICT … DO UPDATE`-Zweig zusätzlich ein `SELECT`-Recht auf
der Zieltabelle — der `DO UPDATE`-Teil liest implizit die bestehende Zeile,
um den Konflikt aufzulösen. Ohne `SELECT` scheitert der Schreibpfad real
mit `SQLSTATE 42501` (insufficient privilege). Das ist kein
Implementierungsfehler im Adapter, sondern PostgreSQL-`ON CONFLICT`-Semantik
und war in `ADR-0047` nicht antizipiert.

Der Implementer hat `ADR-0047` selbst **nicht** verändert (Hard Rule 3.5,
`AGENTS.md` §3.5 — Accepted-ADRs sind immutable), sondern nur den
tatsächlichen Grant in `tools/schema/nacharbeit-roles.sql` korrigiert und
real bestätigt. Diese ADR trägt die dadurch fällige Korrektur nach —
Verdikt 2 aus Modul 8 §Konflikt-Pfad: die ADR wird per Folge-ADR
`supersedes`d, nicht der Implementer-Fund als „Lockerung" durchgewinkt.

**Was von diesem Fund *nicht* betroffen ist:** Die eigentliche Entscheidung
von `ADR-0047` — drei getrennte Verbindungs-DSNs, welcher Aufrufer welche
Rolle bekommt (`cdc_capture`/`cdc_admin`/`cdc_reader`) — bleibt unverändert
stehen. Betroffen ist ausschließlich der exakte SQL-Text einer in ihren
Konsequenzen genannten Folgepflicht.

## Entscheidung

Wir korrigieren den in `ADR-0047` §Konsequenzen genannten Grant-Text auf:

```sql
GRANT SELECT, INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin;
```

Alle übrigen Entscheidungen aus `ADR-0047` (Drei-DSN-Vertrag, Rollen-Zuordnung
je Aufrufer, `NOLOGIN`-Gruppenrollen-Modell, `REPLICATION`-Attribut-Vererbung)
bleiben unverändert bestehen und werden hier nicht wiederholt — verbindlich
bleibt `ADR-0047` selbst für alles außer diesem einen Grant-Text.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — **`SELECT` zum bestehenden Grant ergänzen (gewählt)** | minimale Änderung — ein zusätzliches Recht auf derselben Rolle, derselben Tabelle; `cdc_admin` bleibt exakt die Rolle, die `ADR-0047` bereits für den Heartbeat-Schreibpfad vorsah; kein Eingriff in Adapter-Code oder Rollen-Zuordnung nötig; real getestet (dreimal `make test-store` grün) | `cdc_admin` bekommt damit auch ein reines Lese-Recht auf `cdc.process_heartbeat`, das über den strikten Schreib-Bedarf hinausgeht — minimale Privilegien-Ausweitung ggü. dem ursprünglich geplanten Zuschnitt (bleibt aber innerhalb derselben Tabelle, kein neuer Objekt-Zugriff) |
| B — `UPSERT` über eine SQL-Funktion mit `SECURITY DEFINER` umgehen (die Funktion liefe mit den Rechten ihres Eigentümers, `cdc_admin` bräuchte nur `EXECUTE`) | `cdc_admin` behielte exakt `INSERT`/`UPDATE`, kein `SELECT`-Grant nötig; Least-Privilege im engsten Sinn gewahrt | verlagert die Lese-Semantik nicht weg, sondern nur den Ort, an dem sie erteilt wird — die `SECURITY DEFINER`-Funktion selbst braucht dieselbe implizite Leseberechtigung, nur unter einer anderen Identität (Funktionseigentümer statt Rolle); zusätzliche Angriffsfläche durch `SECURITY DEFINER` (privilegien-eskalierender Code-Pfad, den PostgreSQL explizit als Vorsichtsfall dokumentiert); ADR-0046 hat gerade erst SQL-Funktionen als reine Driving-Adapter ohne eigene Businesslogik etabliert — eine neue `SECURITY DEFINER`-Funktion für einen reinen Grant-Umweg widerspräche dieser gerade getroffenen Trennung ohne neuen fachlichen Bedarf |
| C — den `ON CONFLICT`-Zweig im Heartbeat-Adapter vermeiden (z. B. `SELECT … FOR UPDATE` + bedingtes `INSERT`/`UPDATE` in Anwendungslogik) | vermeidet das `SELECT`-Erfordernis von PostgreSQLs `ON CONFLICT`-Mechanik vollständig | ersetzt eine atomare Datenbank-Operation durch zwei Round-Trips mit Race-Window (Read-Modify-Write ohne Datenbank-seitige Atomarität) oder erzwingt explizite Sperren-Disziplin im Adapter — Mehraufwand und neues Fehlerpotenzial für ein Problem, das ein einzeiliger Grant löst; kein Codepfad existiert dafür bereits, reine Neu-Implementierung ohne fachlichen Mehrwert |
| D — Status quo beibehalten (`GRANT INSERT, UPDATE` wie in `ADR-0047` benannt) | keine Änderung nötig | real widerlegt: Der Heartbeat-Schreibpfad scheitert nachweislich mit `SQLSTATE 42501` (Commit `a32a2c9`, dreimal reproduziert) — verfehlt `LH-QA-SEC-001`/`002`/`003` nicht durch zu weite, sondern durch **zu enge**, nicht funktionsfähige Rechte; keine gangbare Option |

## Konsequenzen

- Positiv: Der Heartbeat-Schreibpfad funktioniert real unter `cdc_admin`
  (belegt durch `internal/bootstrap/roles_wiring_test.go`, Commit `a32a2c9`).
- Positiv: Die Korrektur bleibt lokal — kein Eingriff in die Rollen-Zuordnung
  oder den Drei-DSN-Vertrag aus `ADR-0047`.
- Negativ: `cdc_admin` trägt nun ein `SELECT`-Recht auf
  `cdc.process_heartbeat`, das über den reinen Schreib-Bedarf hinausgeht —
  eine minimale, dokumentierte Abweichung vom engsten denkbaren
  Least-Privilege-Zuschnitt, innerhalb derselben Tabelle.
- Folgepflicht: keine weitere — `tools/schema/nacharbeit-roles.sql` trägt
  den korrigierten Grant bereits (Commit `a32a2c9`); diese ADR dokumentiert
  die Entscheidung nachträglich, ändert am Code nichts mehr.
- Folgepflicht (Doku): `docs/plan/planning/in-progress/slice-023-rollen-spezifische-dsn-verdrahtung.md` <!-- d-check:status-provenance -->
  Kopf-`Bezug:`-Feld verweist zusätzlich auf diese ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `internal/bootstrap/roles_wiring_test.go` gegen PostgreSQL-Testcontainer | Der Heartbeat-Schreibpfad (`INSERT … ON CONFLICT … DO UPDATE` gegen `cdc.process_heartbeat`) gelingt unter `cdc_admin` nur, wenn `SELECT` **und** `INSERT`/`UPDATE` erteilt sind — mit nur `INSERT`/`UPDATE` schlägt er mit `SQLSTATE 42501` fehl (Vorher/Nachher-Beleg, bereits in Commit `a32a2c9` real getestet) | `make test-store` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Permanent — der Grant folgt aus PostgreSQLs `ON CONFLICT`-Implementierung,
kein produktspezifischer oder versions-gebundener Zustand, der sich
absehbar ändert. Re-Evaluierung nur fällig, falls der Heartbeat-Adapter den
`ON CONFLICT`-Zweig künftig durch einen anderen Schreibmechanismus ersetzt
(dann entfiele das `SELECT`-Erfordernis möglicherweise, wäre aber ein neuer
fachlicher Anlass, keine Bestätigung dieser ADR).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-12 | Accepted — Anlass: real widerlegter Grant-Text aus `ADR-0047` §Konsequenzen (Commit `a32a2c9`, `internal/bootstrap/roles_wiring_test.go`, dreimal `make test-store` grün); Verdikt 2 aus Modul 8 §Konflikt-Pfad (Folge-ADR statt stiller Implementer-Lockerung) | `slice-023` (in `in-progress/`) <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0048` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
