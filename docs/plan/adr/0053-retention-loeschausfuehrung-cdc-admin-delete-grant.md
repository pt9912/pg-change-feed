# ADR-0053: Retention-Löschausführung bindet an `cdc_admin` — `DELETE`-Grant-Ergänzung

**Status:** Accepted — Supersedes [`ADR-0047`](0047-rollenspezifische-dsn-verdrahtung.md)
(nur die Rollen-Zuordnungstabelle wird um die Retention-Löschausführung als
neuen Aufrufer ergänzt; der Drei-DSN-Vertrag und alle übrigen
Rollen-Zuordnungen aus `ADR-0047` bleiben unverändert — dieselbe begrenzte
Supersession-Form wie [`ADR-0048`](0048-heartbeat-grant-korrektur-select-ergaenzung.md))

**Datum:** 2026-09-13

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-RET-002`](../../../spec/lastenheft.md)…[`004`](../../../spec/lastenheft.md),
[`LH-QA-SEC-001`](../../../spec/lastenheft.md) (Least-Privilege),
[`LH-QA-SEC-002`](../../../spec/lastenheft.md) (Getrennte Berechtigbarkeit),
[`LH-QA-SEC-003`](../../../spec/lastenheft.md) (Beschränkbarkeit von
CDC-Datenzugriffen), [ADR-0014](0014-retention-domain-policy.md)
(Retention als Domain Policy), [ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md)
(Rollen-spezifische DSN-Verdrahtung — die hier ergänzte Zuordnung fehlte dort
vollständig, weil die Retention-Löschausführung zum Zeitpunkt von `ADR-0047`
noch kein realer Aufrufer war), [ADR-0048](0048-heartbeat-grant-korrektur-select-ergaenzung.md)
(Präzedenzfall: Grant-Korrektur an `ADR-0047` per Folge-ADR, nicht per
stiller Implementer-Lockerung), der vorausgehende Architect-Verdikt zur
Retention-Löschausführung
(Frage 1 — prüfte Domain-/Port-/ADR-Ebene, nicht die Rollen-/Grant-Konsequenz),
der Architect-Verdikt zu `slice-044` und dem Rollen-Grant <!-- d-check:status-provenance -->
(dieser Konflikt-Pfad, Modul 8), `docs/plan/planning/welle-13.md` §6,
`docs/plan/planning/in-progress/slice-044-retention-hintergrundjob.md` §3 <!-- d-check:status-provenance -->
Plan-Nachzug Punkt 6, `docs/plan/planning/observations/BEO-PGC/architect-verdikt-rollen-scope-luecke`,
Commit `824e001` (`tools/schema/nacharbeit-roles.sql`)

**Schärft:** [`ARC-007`](../../../spec/architecture.md) — Bootstrap/
Composition Root (dieselbe Sicht-Stelle, die bereits [ADR-0026](0026-composition-root.md),
[ADR-0047](0047-rollenspezifische-dsn-verdrahtung.md) und
[ADR-0048](0048-heartbeat-grant-korrektur-select-ergaenzung.md) schärfen;
diese ADR ändert an der Sicht-Aussage selbst nichts — sie ergänzt die
Rollen-Zuordnung um einen zur Bootstrap-Zeit von `ADR-0047` noch nicht
existierenden Aufrufer).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0014` (Accepted, `permanent`) entscheidet, dass der Storage-Adapter die
physische Löschung freigegebener `cdc.change`-Zeilen ausführt, während der
Domain Core (`RetentionPolicy.AllowsDeletion`) allein entscheidet, *was*
gelöscht werden darf. `ADR-0014` trifft **keine** Aussage darüber, unter
welcher PostgreSQL-Rolle diese physische Löschung läuft — das ist der
Gegenstand von `ADR-0047`. `ADR-0047`s Rollen-Zuordnungstabelle (datiert
2026-09-12) listet jeden zum damaligen Zeitpunkt existierenden Aufrufer
(Store-Adapter für `Persist`, Tabellen-Aktivierung, Heartbeat,
Consumer-Registrierung, Healthcheck) — die Retention-Löschausführung fehlt
darin vollständig, weil `RunRetentionUseCase` zu diesem Zeitpunkt noch nicht
existierte (er entstand erst mit `slice-043`, einen Tag später). <!-- d-check:status-provenance -->

Der Architect-Verdikt zur Retention-Löschausführung, vor Eröffnung von
`welle-13`, prüfte in
Frage 1 ausführlich, ob die neue `ChangeStorePort`-Delete-Methode und
`RunRetentionUseCase` gegen `ADR-0009`/`0011`/`0012`/`0014`/`0029` bestehen
— auf **Domain-/Port-/Anwendungsfall-Ebene**. Er prüfte **nicht**, unter
welcher PostgreSQL-Rolle die physische Löschung in der Verdrahtung real
läuft und ob diese Rolle das dafür nötige `DELETE`-Grant trägt. `welle-13`
§6 schloss eine „Erweiterung der Least-Privilege-Rollen" deshalb explizit
als Out-of-Scope aus, mit der ausdrücklichen Bedingung: „träfe das während
der Umsetzung doch zu, wäre das ein eigener Architect-Zug, kein stiller
Fortschritt dieser Welle."

Während der Umsetzung von `slice-044` (Hintergrund-Job, der <!-- d-check:status-provenance -->
`RunRetentionUseCase` real verdrahtet) zeigte ein realer Grant-Abgleich: Die
Retention-Pools binden an `cfg.AdminDSN`/`cdc_admin` — dieselbe
Kategorisierung wie Heartbeat, Tabellen-Aktivierung und
Consumer-Registrierung (alle nicht Teil des Erfassungs-INSERT-Pfads) —, aber
weder `cdc_capture` noch `cdc_admin` trugen vor `slice-044` ein <!-- d-check:status-provenance -->
`DELETE`-Grant auf `cdc.transaction`/`cdc.change`. Ohne dieses Grant
scheitert die Löschausführung real, sobald `CDC_ADMIN_DSN` an eine
rollenbeschränkte Login-Identität bindet (der in
`docs/user/benutzerhandbuch.md` dokumentierte Betriebsmodus) — im
Compose-E2E-Lauf unsichtbar, weil dort alle drei DSNs mit dem Superuser
verbunden sind.

Die Implementer-Session lief ohne separaten Architect-Kontext. Statt den
Slice ungelöst zurückzuführen (`in-progress` → `next`/`open`) oder still
fortzuschreiten, wurde die minimale Grant-Ergänzung vorgenommen (`GRANT
SELECT, DELETE ON cdc.transaction, cdc.change TO cdc_admin;`, Commit
`824e001`), real regressionsgetestet
(`TestCdcAdminRetentionDeleteChangesRequiresGrant`,
`TestCdcWiringCallerRejectsWrongRoleAssignment` neuer Fall) und transparent
als Beobachtung geführt
(`BEO-PGC/architect-verdikt-rollen-scope-luecke`). Das ist real dieselbe
Konstellation wie bei `ADR-0048`: Der Implementer hat `ADR-0047` selbst
**nicht** verändert (Hard Rule 3.5) und keine der drei Rollen neu
zugeschnitten, aber `welle-13` §6s ausdrücklicher Bedingung „ein eigener
Architect-Zug" nicht entsprochen. Diese ADR trägt den fälligen
Architect-Zug nach — Verdikt 2 aus Modul 8 §Konflikt-Pfad: Folge-ADR,
`supersedes`, nicht Herabstufung des Befunds, weil der Implementer aus
guter Absicht und transparent gehandelt hat.

**Was von diesem Fund *nicht* betroffen ist:** Der Drei-DSN-Vertrag aus
`ADR-0047` (`CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`), die
`NOLOGIN`-Gruppenrollen-Form und alle übrigen Zuordnungen der Tabelle
bleiben unverändert. `cdc_capture` bleibt strikt auf `INSERT` gegen
`transaction`/`change` beschränkt (`LH-QA-SEC-002`) — unangetastet von
dieser Entscheidung.

## Entscheidung

Wir wählen **A — `cdc_admin` trägt zusätzlich `SELECT`, `DELETE` auf
`cdc.transaction`/`cdc.change`** (bereits umgesetzt in
`tools/schema/nacharbeit-roles.sql`, Commit `824e001`; hiermit nachträglich
architektonisch bestätigt statt zurückgebaut). `ADR-0047`s Rollen-Tabelle
wird um folgende Zeile ergänzt (Fortschreibung, kein Widerspruch zur
bestehenden Tabelle):

| Aufrufer | Env-Var | Rolle | Begründung |
|---|---|---|---|
| Retention-Löschausführung (`RunRetentionUseCase`, `ChangeStorePort.DeleteChanges`/`DeleteOrphanedTransactions`) | `CDC_ADMIN_DSN` | `cdc_admin` | periodischer Hintergrundzug, kein Erfassungs-Nutzlast-Pfad — dieselbe Kategorisierung wie Heartbeat/Tabellen-Aktivierung/Consumer-Registrierung; physisch löscht nur, was `RetentionPolicy.AllowsDeletion` im Domain Core bereits freigegeben hat (`ADR-0014`) |

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — **`cdc_admin` bekommt `SELECT`/`DELETE` (gewählt, bereits umgesetzt)** | keine neue Infrastruktur (kein viertes DSN/Login/Compose-Var/Doku-Abschnitt); folgt exakt dem in `ADR-0047` bereits etablierten Kategorisierungs-Prinzip „alles außer dem Erfassungs-`INSERT`-Pfad läuft über `cdc_admin`" (Heartbeat, Tabellen-Aktivierung, Consumer-Registrierung sind bereits so eingeordnet); `cdc_admin` trägt bereits volles DML auf vier Steuerungstabellen und `CREATE` auf der Datenbank — die zusätzliche Fähigkeit erweitert keine neue Vertrauensstufe, nur die Objektmenge; real regressionsgetestet (`make test-store`, zwei neue Fälle) | vergrößert den Blast-Radius eines kompromittierten `cdc_admin`-Logins um die Fähigkeit, `cdc.transaction`/`cdc.change`-Zeilen unabhängig von `RetentionPolicy.AllowsDeletion` zu löschen — das SQL-Grant selbst ist nicht auf policy-freigegebene Zeilen beschränkt, nur die Anwendungslogik ruft `DELETE` ausschließlich für freigegebene Zeilen auf; bislang war der schlimmste Fall eines kompromittierten `cdc_admin`-Logins die Korruption von Steuerungsdaten (Consumer-Positionen, Registrierung, Heartbeat-Status) — im Prinzip rekonstruierbar; jetzt zusätzlich der unwiderrufliche Verlust erfasster Change-Historie |
| B — eigene, vierte Rolle `cdc_retention` (nur `SELECT`/`DELETE` auf `transaction`/`change`, eigenes `CDC_RETENTION_DSN`) | stärkste Isolation: ein kompromittiertes `cdc_admin`-Login (Registrierung/Aktivierung/Heartbeat) kann keine Kernstore-Daten löschen, und umgekehrt; entspricht demselben Isolations-Prinzip, das `ADR-0047` bereits für das `REPLICATION`-Attribut auf `cdc_capture` anwendet (genuin unterschiedliche Risikoklassen auf unterschiedliche Prinzipale) | ist genau das Ereignis, das `ADR-0047`s eigener Re-Evaluierungs-Trigger (a) benennt („eine vierte CDC-Rolle mit eigenem Aufgabenschnitt wird eingeführt") — der Trigger fragt dort ausdrücklich, ob „eine DSN je Rolle" bei vier Rollen noch trägt oder unhandlich wird; viertes DSN/Login/Compose-Var/Doku-Abschnitt für einen einzelnen Hintergrundzug, während alle bisherigen Hintergrundzüge (Heartbeat, Tabellen-Aktivierung, Administration, WAL-Retention-Check) ohne eigene Rolle auskommen; kein Aufgabenschnitt-Unterschied zu den übrigen `cdc_admin`-Hintergrundzügen außer der Ziel-Tabelle — die Kategorie „periodischer Verwaltungs-Hintergrundzug, kein Erfassungspfad" bleibt dieselbe |
| C — `cdc_capture` bekommt das `DELETE`-Grant statt `cdc_admin` | keine neue Rolle, kein neues DSN | widerspricht `ADR-0047`s expliziter, unmissverständlicher Kategorisierung „der Erfassungspfad bleibt `cdc_capture` auf `INSERT` beschränkt" (`LH-QA-SEC-002`) — `cdc_capture` ist die Rolle, die über den Replication-Stream läuft und strukturell am nächsten an nicht selbst kontrollierten Quelldaten sitzt; ihr eine destruktive Fähigkeit auf denselben Tabellen zu geben, in die sie schreibt, wäre die am wenigsten geeignete Stelle für diese Erweiterung; keine gangbare Option |
| D — Status quo beibehalten (kein Grant) | keine Änderung | real widerlegt: Die Löschausführung scheitert nachweislich, sobald `CDC_ADMIN_DSN` an eine rollenbeschränkte Login-Identität bindet — im Compose-E2E-Lauf unsichtbar (Superuser für alle drei DSNs), aber der von `docs/user/benutzerhandbuch.md` dokumentierte reale Betriebsmodus; verfehlt `LH-FA-RET-002`…`004` vollständig; keine gangbare Option |

**Warum A statt B:** Der Unterschied zwischen `cdc_admin`s bisherigen
Fähigkeiten (destruktive DML auf Steuerungstabellen) und der neuen Fähigkeit
(destruktive DML auf Kernstore-Tabellen) ist real und wird in den
Konsequenzen unten benannt, nicht kleingeredet — er trägt aber nicht, eine
vierte Rolle **jetzt** einzuführen: Der Aufgabenschnitt der
Retention-Löschausführung (periodischer Hintergrundzug, kein
Erfassungspfad, physisch ausführend für eine bereits im Domain Core
getroffene Freigabe-Entscheidung) ist derselbe wie bei Heartbeat und
Tabellen-Aktivierung, nur die betroffene Tabelle unterscheidet sich. Eine
vierte Rolle wäre eine Entscheidung, die allein aus der higher-stakes-Natur
der Zieltabelle folgt, nicht aus einem neuen Aufgabenschnitt — genau die
Unterscheidung, die der Re-Evaluierungs-Trigger unten offen hält, statt sie
hier vorwegzunehmen.

## Konsequenzen

- Positiv: `RunRetentionUseCase` funktioniert real unter einer
  rollenbeschränkten `cdc_admin`-Login-Identität (Fitness Function unten) —
  `LH-FA-RET-002`…`004` sind damit auch im dokumentierten Least-Privilege-
  Betriebsmodus erfüllbar, nicht nur im Compose-Superuser-Lauf.
- Positiv: keine neue Infrastruktur (DSN/Login/Compose-Variable/
  Benutzerhandbuch-Abschnitt) — die Erweiterung bleibt lokal auf eine
  Grant-Zeile in `nacharbeit-roles.sql` beschränkt.
- Negativ: `cdc_admin` trägt nun eine destruktive Fähigkeit auf den beiden
  Kernstore-Tabellen (`cdc.transaction`, `cdc.change`) — eine kompromittierte
  `cdc_admin`-Login-Identität kann diese Zeilen unabhängig von
  `RetentionPolicy.AllowsDeletion` löschen, weil das SQL-Grant die
  Domain-Policy nicht durchsetzt, nur die Anwendungslogik tut das. Dies ist
  eine bewusst in Kauf genommene, dokumentierte Erweiterung des
  `cdc_admin`-Blast-Radius, keine übersehene.
- Folgepflicht: keine weitere Code-Änderung — `tools/schema/nacharbeit-roles.sql`
  trägt den Grant bereits (Commit `824e001`); diese ADR dokumentiert die
  Entscheidung nachträglich und ergänzt `ADR-0047`s Rollen-Tabelle formal.
- Folgepflicht (Doku): `docs/plan/planning/welle-13.md` §6 erhält eine
  Klarstellung, dass die dort formulierte Bedingung („ein eigener
  Architect-Zug") real eingetreten und mit dieser ADR erfüllt ist.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `internal/adapters/driven/postgresstorage/*_test.go` gegen PostgreSQL-Testcontainer | `TestCdcAdminRetentionDeleteChangesRequiresGrant`: `DeleteChanges`/`DeleteOrphanedTransactions` scheitert unter `cdc_admin` ohne das `DELETE`-Grant und gelingt real, sobald es erteilt ist | `make test-store` |
| `internal/bootstrap/roles_wiring_test.go` gegen PostgreSQL-Testcontainer | `TestCdcWiringCallerRejectsWrongRoleAssignment`, Fall „Retention Store-Adapter … mit `cdc_capture`-Login": ein `DELETE`-Aufruf gegen `cdc.change` scheitert real unter `cdc_capture` (`INSERT`-only bleibt durchgesetzt) | `make test-store` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Re-Evaluierung fällig, wenn eines eintritt: (a) ein realer Sicherheitsvorfall
oder eine Compliance-Anforderung (Audit, Zertifizierung) verlangt eine
stärkere Isolation zwischen dem Registrierungs-/Betriebs-Schreibpfad und dem
Löschpfad auf Kernstore-Daten — dann Umsetzung von Option B (eigene
`cdc_retention`-Rolle) prüfen; (b) eine weitere destruktive Fähigkeit auf
Kernstore-Tabellen entsteht, die denselben Aufrufer-Typ (periodischer
Hintergrundzug, kein Erfassungspfad) mit der hier entschiedenen teilt — dann
Bündel-Neubewertung, ob „`cdc_admin` trägt alles außer dem
Erfassungs-`INSERT`-Pfad" als Kategorie noch trägt oder in Sub-Kategorien
zerlegt werden sollte (dann zugleich `ADR-0047`s Re-Evaluierungs-Trigger (a)
fällig — „eine vierte CDC-Rolle mit eigenem Aufgabenschnitt"). Sonst
**permanent**.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Anlass: `BEO-PGC/architect-verdikt-rollen-scope-luecke` (1×, unter der 3×-Schwelle, aber realer Rollen-Konflikt nach Modul 8 §Konflikt-Pfad), Commit `824e001` (Implementer-Fund + Grant-Ergänzung ohne separaten Architect-Zug, transparent dokumentiert); Verdikt 2 aus Modul 8 §Konflikt-Pfad (Folge-ADR statt stille Implementer-Lockerung, analog `ADR-0048`) | `slice-044` (in `in-progress/`) <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0053` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
