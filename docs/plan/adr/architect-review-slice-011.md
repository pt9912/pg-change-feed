# Architect-Review slice-011 — Verdikt zu F-1/F-2 (Heartbeat-Pattern)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-10.
**Eingang:** [`docs/reviews/review-slice-011.md`](../../reviews/review-slice-011.md)
F-1 (MEDIUM, Architektur-Frage) und F-2 (MEDIUM, Closure-Trigger nicht mehr
erfüllbar) · Slice-Plan
[`docs/plan/planning/done/slice-011-sicherheit-observability.md`](../planning/done/slice-011-sicherheit-observability.md)
§1/§3/§6 · [`ADR-0020`](0020-http-grpc-optional.md) (Accepted, permanent) ·
[`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md) (Accepted,
Supersedes ADR-0018) · [`ADR-0024`](0024-observability-ausserhalb-der-domain.md)
(Accepted, permanent) · [`ADR-0027`](0027-capture-application-service.md)
(Accepted, permanent).
**Ausgang:** Übergabe-Artefakt an Planner — **kein Folge-ADR nötig**; Verdikt
zu F-1 und Anpassungs-Empfehlung zu F-2.
**Harte Regel eingehalten:** keine Accepted-ADR in-place geändert; keine
Folge-ADR-Datei angelegt (die Prüfung unten verneint den Bedarf); der
Slice-Plan selbst wird von diesem Lauf **nicht** editiert — §5-Anpassung ist
Planner-Sache im nächsten Zug (Modul 8 §Konflikt-Pfad, Sequenz).

---

## Verdikt zu F-1 — Heartbeat-Pattern trägt ohne neue ADR

**Kurzform:** Der Reviewer hat recht. Der Implementer-Rückzug war zu breit
begründet. Health-Endpoint per Heartbeat-Tabelle + SQL-View ist eine offene
**Umsetzungs**-Frage (Folge-Slice), keine offene **Architektur**-Frage
(kein Folge-ADR).

### Prüfung: Schreib-Seite (Heartbeat-Tabelle)

Der laufende CDC-Prozess (Capture Application Service, ADR-0027) schreibt
periodisch seinen Lebenszeichen-Zustand in `cdc.process_heartbeat`. Das ist
architektonisch **kein** SQL-Driving-Adapter im Sinne von ADR-0046 — dort
geht es um SQL-Objekte, die ein **externer** Client aufruft (Views/
Funktionen). Der Heartbeat-Schreiber ist der umgekehrte Fall: der Prozess
selbst persistiert Betriebsinformation über die eigene DB-Verbindung — exakt
das Muster, das [`ADR-0024`](0024-observability-ausserhalb-der-domain.md)
bereits entschieden hat: „Betriebsinformationen … müssen aus allen Schichten
ankommen"; Träger sind Outbound Ports + Driven Adapters (dort für
`MetricsPort`/`EventSinkPort`, dieselbe Klasse trägt eine
Heartbeat-Persistenz). `spec/architecture.md` §`ARC-006` führt „Telemetry"
bereits als eigene Driven-Adapter-Fähigkeit — ein Heartbeat-Schreiber ist
eine Erweiterung dieser bereits deklarierten Fläche, kein neuer
Adapter-**Typ**.

Der periodische Schreib-**Zeitpunkt** (Timer statt „ein Aufruf pro
eingehender Replication-Message") ist ein Ausführungsdetail der
Composition-Root-Verdrahtung ([`ADR-0026`](0026-composition-root.md)), nicht
eine neue Schicht- oder Abhängigkeitsentscheidung: Er fügt dem bereits
laufenden, langlebigen Capture-Prozess (ADR-0027) einen weiteren
Aufruf-Auslöser für denselben Port-Adapter-Mechanismus hinzu, den ADR-0024
schon autorisiert. Keine Domain-Berührung, keine neue Abhängigkeitsrichtung,
kein neuer Driving-Adapter. Ein Timer/Goroutine-Zuschnitt ist Implementierung
unter bereits Entschiedenem, keine Architekturentscheidung.

### Prüfung: Lese-Seite (SQL-View)

Eine vierte Lese-View (nach `active_tables`, `consumer_status`, `changes`,
`metrics`) liest `cdc.process_heartbeat` und projiziert Alter/Zeitstempel des
letzten Lebenszeichens. Das ist deckungsgleich mit
[`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md) Kategorie C:
reine Projektion über bereits persistierte, bereits validierte Daten, keine
Domänenregel. Die `SPEC-007`-Klassifikation (`HEALTH_STATES`) bleibt beim
lesenden System — dasselbe bereits gelebte Muster wie bei `cdc.metrics`
(Review-Negativbefund review-slice-011, „`cdc.metrics` als reine
Projektion"). Kein neuer Adapter-Typ, keine neue Entscheidung.

### Prüfung: Berührt das `ADR-0020`-Sperrfeld?

Nein. [`ADR-0020`](0020-http-grpc-optional.md) stellt **HTTP/gRPC als
Driving Adapter** zurück — neue Netz-Listener, die einen konkreten
API-Consumer-Bedarf voraussetzen (Re-Evaluierungs-Trigger: „Beobachtbarer
Bedarf eines API-Consumers"). Weder die Tabellen-Schreibseite (kein Driving
Adapter, kein Netz-Listener) noch die View-Leseseite (bereits etablierter
SQL-Driving-Adapter-Kanal, `ARC-005`) sind HTTP oder gRPC. `LH-FA-ADM-002`/
`LH-QA-OPS-002` verlangen zudem nicht HTTP — „automatisierte Health Checks"
sind über einen SQL-Aufruf (`psql`/Compose-Healthcheck gegen die neue View)
ebenso automatisierbar wie über HTTP; der SQL-Kanal ist für Status-/
Diagnoseabfragen bereits der etablierte Weg (`spec/lastenheft.md`
`LH-FA-ADM-001`-Akzeptanzkriterium: „`SELECT * FROM cdc.tables`,
`SELECT * FROM cdc.consumers`"). Die Prämisse im Plan-Nachzug („eine echte
Prozess-Liveness-Antwort braucht einen neuen Driving-Adapter-Zuschnitt
(HTTP-Listener, Heartbeat-Datei o. ä.)") trifft für die HTTP-Variante zu, ist
für die Heartbeat-Tabellen-Variante aber falsch — genau die Lücke, die der
Reviewer benennt.

### Einordnung nach Modul 8 §Konflikt-Pfad (sinngemäß, kein Rollenwiderspruch)

Kein Fall „ADR wird per Folge-ADR abgelöst" und kein Fall „Lockerung legitim,
aber undokumentiert". Es ist der erste der drei Pfade: **ADR-0020 gilt
unverändert** (bleibt `Accepted`, permanent, kein neuer Entscheidungsanlass
eingetreten) — **der Plan-Nachzug hat ihre Sperrwirkung zu breit auf einen
Fall ausgedehnt, den sie nicht erfasst.** Der tatsächlich tragende Grund für
den Rückzug ist eine **Schicht-Abgrenzung**: Ein Heartbeat-Schreiber
berührt den laufenden Prozess (Application/Bootstrap-Verdrahtung), diesen
Slice hält sich aber „auf die DB-Schicht (DDL/Compose)" (§1) — exakt dieselbe
Grenze, die der Implementer im selben §1 für den `wiring.go`-Ausschluss
bereits selbst zieht. Die Eskalation an den Architect (Folge-ADR-Pflicht)
war damit unnötig; ein einfacher Folge-Slice-Verweis hätte genügt.

**Ergebnis: Kein Folge-ADR.** Health-Endpoint per Heartbeat-Muster ist
Folge-Slice-Arbeit, keine offene Architektur-Frage.

---

## Disposition zu F-2 — §5 Closure-Trigger

**Trigger-Anpassung: ja.** Da der Health-Teil (Heartbeat schreiben + View
lesen) mehrere Schichten gleichzeitig berührt — Application/Capture-Service
(Timer-Auslöser), Bootstrap/Composition-Root (Verdrahtung des neuen
Outbound-Port-Adapters), DB-Schicht (neue Tabelle + View) —, ist er nach
denselben Größen-Kriterien, die dieser Slice bereits selbst für
`wiring.go` anwendet (Modul 5: „mehrere Schichten betroffen" →
`in-progress→next`-Kriterium), **nicht** mehr innerhalb von slice-011s
eigenem, bereits auf die DB-Schicht reduzierten Rahmen lieferbar. Der Slice
in diesem Lauf zusätzlich um den Heartbeat-Mechanismus zu erweitern hieße,
exakt die Schicht-Grenze wieder zu überschreiten, die der Implementer für
`wiring.go` bewusst gezogen hat, um nicht in die vorab benannte
Rückführung „zu groß" (§4) zu laufen.

**Empfehlung an den Planner** (nicht bindend — Architect greift nicht in den
Plan ein, Modul 8):

- §5 zweiteilen: Metriken-Closure („Metriken-Endpoint … am verdrahteten
  System belegt") bleibt Closure-Bedingung *dieses* Slice. Der Health-Teil
  verlässt §5 und wandert nach §1 (Ausschluss-Klasse **Folge-Slice**, mit
  Kennung, sobald die Welle den Folge-Slice eröffnet) bzw. §6 (Ausgang
  „weiter offen" nur solange keine Kennung existiert — sobald der
  Folge-Slice angelegt ist, Ausgang „eingetreten" mit dessen Kennung).
- Der neue Folge-Slice trägt Bezug `LH-FA-ADM-002`, `LH-QA-OPS-002` und
  Kopf-Referenz auf `ADR-0024`/`ADR-0027`/`ADR-0046` (kein `ADR-0020`-Bezug
  nötig, da nicht berührt) sowie auf dieses Verdikt-Dokument als Beleg, warum
  kein Folge-ADR vorausgesetzt ist.
- Kein Carveout nötig: Das ist kein rotes Gate, sondern eine
  Scope-Reduktion des Closure-Triggers auf Basis einer bereits im Plan
  vorab benannten Rückführungs-Grenze (§4) — die Trigger-Anpassung *ist*
  der saubere Ausgang, kein Ausnahme-Mechanismus.

---

## Beleg-Anker (Kurzfassung)

| Aussage | Tragende ADR/Spec |
|---|---|
| Heartbeat-Schreiben ist Driven-Adapter-Erweiterung, kein neuer Adapter-Typ | `ADR-0024` · `ARC-006` (Telemetry) |
| Heartbeat-Schreib-Auslöser (Timer) ist Bootstrap-/Verdrahtungsdetail | `ADR-0026` · `ADR-0027` |
| Heartbeat-Lese-View ist reine Projektion, Kategorie C | `ADR-0046` |
| Weder Schreib- noch Lese-Seite berühren HTTP/gRPC | `ADR-0020` (Sperrfeld bleibt unberührt) |
| Automatisierte Health Checks erfordern kein HTTP | `LH-FA-ADM-002`, `LH-QA-OPS-002`, `LH-FA-ADM-001`-Akzeptanzkriterium (SQL-Kanal) |
| Health-Heartbeat berührt mehrere Schichten → nicht mehr in slice-011s reduziertem Rahmen | Modul 5 §Ziel-Form: Slice (Größen-Kriterium „mehrere Schichten"), slice-011 §1/§4 (Schicht-Abgrenzung `wiring.go`) |
