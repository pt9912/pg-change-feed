Zustand: offen — **Punkt 1 ist geschlossen** (`slice-093`, s. unten), **Punkt 2
bleibt offen**; Ausgang: **weiter offen** → eine
Rollen-Vertauschungsprüfung, die real die Replication-Stream-/ACK-Adapter
(`receive.NewStream`/`postgresack.New`) mit rollenbeschränkten
Verbindungen aufruft (Punkt 2 — `slice-028` schloss nur die
PostgreSQL-Server-Ebene, nicht die Adapter-Ebene), sind eigenständiger
Aufwand; kein Slice dafür existiert.
Zähler (abgeleitet): **3×** (evidence/slice-023.md, evidence/slice-028.md,
evidence/slice-backfill-run-store.md) — **Schwelle erreicht mit
`slice-backfill-run-store`**; der Ausgang gehört dem Lese-Schritt der Closure von
`welle-backfill-bestand` (Modul 6), nicht entschieden.

**Der dritte Beleg trifft Punkt 2 in einer weiteren Ausprägung** (Review F-1, F-2 zu
`slice-backfill-run-store`): die Zusage der drei Backfill-Adapter und des
Administrationswegs hängt an der **Rolle der Verbindung**, und jeder Test lief als
Superuser — `cdc_admin` trug auf `cdc.administration_request` kein Recht, und kein Test sah
es. Geschlossen ist damit die Adapter-Ebene **für diese Adapter und diesen Weg**:
`backfillroles_test.go` (Annahme unter `cdc_admin`-, Run-Zustand und Schreiber unter
`cdc_capture`-Login, je unter der Rolle des anderen mit SQLSTATE `42501`) und
`internal/bootstrap/administration_roles_internal_test.go` (`processAdministrationRequests` mit
allen vier Antragsarten unter einem `cdc_admin`- und einem `cdc_capture`-Login); der Verifier
färbte sie mit der Mutation „Grant streichen" rot. Für Replication-Stream und ACK-Adapter
bleibt Punkt 2 **offen**. **Ausgangs-Kandidat** (Architect-Entscheidung im Lese-Schritt, nicht
getroffen): eine Rollen-Zusage (Grant, Verbot, Zuordnung Adapter → Rolle) braucht einen Test
unter dem **Login der Rolle** (`CREATE ROLE … LOGIN … IN ROLE`, kein Superuser, kein Eigentum
an `cdc`-Objekten), nicht unter dem Superuser; die Eingabeseite eines Rechte-Tests ist die
Rolle der Verbindung (derselbe Mechanismus wie `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
dort nicht doppelt gezählt). Träger-Kandidaten: Schritt 19 von `implement-slice` und der
HIGH-Punkt „Zusage ohne Bindung an ihre Eingabeseite" im Reviewer-Skill; das Muster der Logins
liefern `newTestLoginRole` (`internal/bootstrap`) und `backfillroles_test.go` (`postgresstorage`).

**Punkt 1 ist mit `slice-093` geschlossen — der Vorgang legt keinen Beleg an.**
Der Slice hat die Lücke **behoben**, statt sie erneut zu treffen: der neue Test
`TestRolloutDateiTraegtDieRechteDerVerdrahtung` liest die **echte**
`tools/schema/nacharbeit-roles.sql` und bindet sieben Zusagen an ihren Inhalt;
eine Mutation in der Datei färbt **Test und Gate** (Bau-EC 1).
**Präzise, mit der Verifikation (F-4):** geschlossen ist die Lücke **im
Gate-Bündel** — nicht „nirgends"; der Integrationslauf liest den Heartbeat-Grant
real (er rollt die Datei aus und verlangt den Compose-Status `healthy`).
**Punkt 2 bleibt offen** (Rollen-Vertauschung für Replication-Stream/ACK auf
Adapter-Ebene; sie braucht rollenbeschränkte Logins in
`tools/harness/run-replication-tests.sh`, unberührt).
Ein **Beleg** wird dafür nicht angelegt: der Eintrag wird nicht erneut getroffen,
sondern aufgelöst — und ein Vorgang, der eine Lücke schließt, ist kein Vorkommen.
Sein Ausgang gehört dem **Lese-Schritt der `welle-20`-Closure** (Modul 6).
