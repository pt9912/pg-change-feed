# ADR-0055: NATS-Change-Notification als reines Wecksignal (Core NATS, kein JetStream)

**Status:** Accepted

**Datum:** 2026-09-13

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-13)

**Bezug:** [`LH-FA-SST-007`](../../../spec/lastenheft.md), [`LH-FA-REA-001`](../../../spec/lastenheft.md), [`LH-QA-REL-001`](../../../spec/lastenheft.md), [ADR-0007](0007-source-ack-outbound-port.md), [ADR-0011](0011-persist-before-ack.md), [ADR-0023](0023-fehlerklassifikation.md), [ADR-0027](0027-capture-application-service.md)

**Schärft:** [`ARC-004`](../../../spec/architecture.md), [`ARC-006`](../../../spec/architecture.md), [`ARC-013`](../../../spec/architecture.md) (neu, diese ADR ergänzt die Zeile), [`SPEC-017`](../../../spec/pflichtenheft.md) (neu, diese ADR legt den Inhalt fest)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-SST-007`](../../../spec/lastenheft.md) verlangt, dass ein über NATS
verbundener Consumer über neue Changes erfährt, ohne den bestehenden
Lesezugriffsweg (SQL/API, [`LH-FA-REA-001`](../../../spec/lastenheft.md) ff.)
pollen zu müssen. Das Lastenheft grenzt dabei ausdrücklich ab, was **nicht**
seine Sache ist: Ob NATS nur ein Wecksignal oder die vollständigen
Change-Inhalte trägt, das konkrete Subjekt-/Stream-Schema und die
Zustellgarantie (JetStream vs. Core NATS) sind explizit als
Architektur-/Spezifikationsfragen benannt. Sein Boundary- und
Negative-Akzeptanzkriterium sind dabei bindend: Ein nicht verbundener
Consumer verliert keine Änderung — sie bleibt über den bestehenden
Zugriffsweg abrufbar; nach einer unterbrochenen NATS-Verbindung muss der
Consumer über denselben bestehenden Zugriffsweg nachholen. Wörtlich: „NATS
ersetzt nicht die Nachvollziehbarkeit des Zugriffswegs."

Diese Formulierung ist der entscheidende Architektur-Hebel: Die
Nachvollziehbarkeit (Persistenz, geordnete Historie, Nachholbarkeit) liegt
bereits vollständig und bewährt beim `cdc.changes`-Lesepfad
(`ChangeStorePort`/SQL-Views, `ARC-004`/`ARC-005`/`ARC-006`). NATS muss
diese Eigenschaft nicht selbst tragen — ein verpasstes oder verlorenes
NATS-Signal ist folgenlos, weil jeder Consumer ohnehin nur über den
bestehenden Pfad autoritativ liest. Das reduziert den Lösungsraum auf eine
reine Weck-/Trigger-Semantik ohne eigene Zustellgarantie.

Grüne Wiese: Kein NATS-Treffer in `spec/pflichtenheft.md`,
`spec/architecture.md`, `docs/plan/adr/`, `go.mod`, `compose.yaml` oder
`Dockerfile` — diese Entscheidung eröffnet das Feld vollständig neu.

Die bestehende Persist-before-ACK-Kette (`ADR-0011`, `ADR-0027`,
`CaptureService.Capture()`, `internal/application/usecase/capture/service.go`)
trägt die Ordnung `Receive → Decode → Persist → COMMIT Store → ACK Source`
als Kritischen Pfad: Ein Fehler an jedem dieser Schritte beendet den
Replication-Stream (`internal/adapters/driving/replication/receive/receive.go`,
Funktion `process`: „Ein Fehler endet den Stream"). Ein Notify-Schritt, der
sich naiv in dieselbe Fehlerpropagation einreiht, würde eine bewusst
best-effort gedachte Zusatzfunktion zu einem Verfügbarkeits-Risiko für den
bereits erfolgreich abgeschlossenen kritischen Pfad machen — genau das
verhindert diese Entscheidung explizit (siehe Punkt 4 unten).

## Entscheidung

Wir wählen **C — Core NATS als reines, verlustbehaftetes Wecksignal ohne
Payload-Inhalt**, mit fünf Festlegungen:

1. **Technologie und Zustellsemantik: Core NATS, kein JetStream.** Core NATS
   liefert Publish/Subscribe ohne Persistenz und ohne Zustellgarantie — ein
   Abonnent, der nicht verbunden ist, erhält die Nachricht nie (Fire-and-
   Forget). Das ist für dieses Merkmal **kein Mangel, sondern die richtige
   Passung**: Da `LH-FA-SST-007` die Nachvollziehbarkeit explizit dem
   bestehenden Zugriffsweg zuweist, braucht kein Notify-Signal selbst zu
   überleben. JetStream (Streams, Konsum-Bestätigung, Replay, At-Least-Once)
   würde eine zweite, parallele Nachvollziehbarkeits-Infrastruktur neben dem
   bereits vorhandenen `cdc.changes`-Pfad aufbauen — Zustandshaltung
   (Stream-Retention, Consumer-Gruppen, Acknowledgment-Fenster), die kein
   Akzeptanzkriterium verlangt und die Frage aufwirft, welche der beiden
   Wahrheiten (JetStream-Stream vs. SQL-Store) im Konfliktfall gilt. Core
   NATS vermeidet diese Frage strukturell: Es gibt nur eine Quelle der
   Wahrheit, den bestehenden Store.

2. **Subjekt-Schema: `cdc.changes.<source_id>`.** Ein Subjekt pro Quelle,
   abgeleitet aus der bereits bestehenden `CDC_SOURCE_ID` — konsistent mit
   der übrigen Skopierung (Publication-, Slot-Namen tragen dieselbe
   Quelle-Bindung, `compose.yaml`). Ein Consumer, der eine bestimmte Quelle
   verfolgen will, abonniert ihr Subjekt direkt; ein Consumer, der über
   mehrere Quellen hinweg lauschen will, abonniert `cdc.changes.>`
   (NATS-Wildcard) — keine gesonderte Aggregations-Infrastruktur nötig, weil
   Core NATS Wildcard-Subscriptions nativ trägt.

3. **Nachrichteninhalt: leerer Payload, reines Trigger-Signal — kein
   Change-Inhalt, keine Positionsangabe.** Verworfene Variante: ein
   JSON-Payload mit `source_id` und neuester Position. Begründung für die
   leere Form: Eine im Payload mitgeführte Position wäre eine zweite,
   potenziell veraltete Repräsentation desselben Zustands, den der
   Store bereits autoritativ trägt (dieselbe Zweitpfad-Sorge wie in
   [`LH-FA-SST-006`](../../../spec/lastenheft.md), Boundary: „dieselbe
   Domänenlogik, kein Zweitpfad" — hier auf Datenrepräsentation statt
   Zugriffsweg übertragen). Ein leeres Signal lässt keinen Konsumenten je
   in Versuchung geraten, eine Position aus NATS statt aus dem Store zu
   übernehmen; jede Benachrichtigung bedeutet ausschließlich „lies erneut
   über den bestehenden Zugriffsweg", ohne dass ihr Fehlen oder ihre
   Verdopplung (At-Least-Once wird hier nicht einmal versprochen) eine
   semantische Bedeutung trüge.

4. **Fehlerklasse und kritischer Pfad — die wichtigste Einzelentscheidung:**
   Ein Notify-Fehlschlag klassifiziert als `transient` (`SPEC-008`,
   `ADR-0023`): eine unterbrochene NATS-Verbindung ist strukturell dieselbe
   Kategorie wie eine „vorübergehend nicht verfügbare Quelle/Speicher" —
   der nats.go-Client trägt bereits eingebautes automatisches Reconnect mit
   Backoff, sodass die in `SPEC-008` für `transient` vorgesehene Aktion
   („Erneut versuchen mit begrenztem Backoff") von der Client-Bibliothek auf
   Verbindungsebene getragen wird, nicht von einer eigenen Retry-Schleife
   auf Nachrichtenebene. **Ein Notify-Fehler darf niemals die bereits
   erfolgte Persistierung oder das bereits erfolgte Source-ACK rückgängig
   machen, verzögern oder blockieren.** Der Notify-Aufruf reiht sich als
   dritter, **optionaler** Outbound-Port-Aufruf **nach** `ACK Source` in
   `CaptureService.Capture()` ein
   (`Receive → Decode → Persist → COMMIT Store → ACK Source → Notify (best
   effort)`); sein Fehler wird an der Aufrufstelle in `CaptureService`
   abgefangen und **nicht** in den Rückgabewert von `Capture()` übernommen
   — anders als bei `store`/`ack`, deren Fehler den Replication-Stream
   beenden (`internal/adapters/driving/replication/receive/receive.go`,
   `process`). Sichtbarkeit bleibt gewahrt über denselben strukturierten
   Fehler-Log-Mechanismus wie beim bestehenden `PostgresReplicationAckAdapter`
   (injizierter `LogPort`, `ADR-0024`) und, sofern verdrahtet, die
   Fehler-Metrik `cdc_errors_total{class="transient"}` (`SPEC-009`). Jede
   nachfolgende Change löst einen neuen Notify-Versuch aus — das ist die
   natürliche Wiederholung auf Nachrichtenebene; ein verpasstes einzelnes
   Signal hat wegen Punkt 3 keine Konsequenz, die über „der Consumer liest
   ohnehin wieder über den bestehenden Weg" hinausginge.

5. **Betrieb: `CDC_NATS_URL`, vollständig optional.** Neue Umgebungsvariable
   nach demselben Namensschema wie die bestehenden `CDC_*`-Variablen. Ist
   sie nicht gesetzt, bleibt das Notify-Feature vollständig deaktiviert —
   `CaptureService` erhält keinen (bzw. einen `nil`/No-Op-)
   `ChangeNotificationPort`, keine NATS-Verbindung wird aufgebaut, das
   bestehende Verhalten bleibt für jedes heutige Deployment (inklusive
   `compose.yaml`) unverändert. Das ist konsistent mit `LH-FA-SST-007`s
   Boundary-Kriterium: Der bestehende Zugriffsweg ist die Grundlage, NATS
   ist eine additive Fähigkeit, keine Voraussetzung. Der NATS-Server selbst
   wird als digest-gepinntes Image geführt, Docker-only-Disziplin
   (`AGENTS.md` §3.1): `nats:2-alpine@sha256:065e8355c20a5575b3c77224be1855e8103fd148b68fba05130b9b8ddfa40ccc`
   (linux/amd64-Manifest desselben Multi-Arch-Tags, ermittelt am
   2026-09-13 über `docker buildx imagetools inspect nats:2-alpine`,
   entspricht `nats-server` 2.14.6 auf Alpine 3.22 — dieselbe Pin-Form wie
   `PG_TEST_IMAGE`/`TOOLCHAIN_IMAGE` im Makefile: die plattformspezifische
   Manifest-Digest, nicht die Multi-Arch-Index-Digest). Der Go-Client
   `github.com/nats-io/nats.go` wird die erste direkte Nicht-PostgreSQL-
   Abhängigkeit dieses Repos (`go.mod` trägt heute nur `pgx`/`pglogrepl`).

**Port-/Adapter-Design** (Folgepflicht des umsetzenden Slices, hier
architektonisch festgelegt): neuer Outbound-Port
`internal/application/port/outbound/changenotification.go`
(`ChangeNotificationPort`, Methode `Notify(ctx context.Context, sourceID
string) error`, analog zum bestehenden Port-Zuschnitt nach Fähigkeiten,
`ADR-0034`) und ein neuer Driven-Adapter
`internal/adapters/driven/natsnotify/` (Konstruktions- und Options-Muster
analog zu `postgresack`: `New(conn *nats.Conn, opts ...Option)`,
`WithLog`, Fehler-Wrapping über einen neuen Sentinel `ErrNotify` in die
Klasse `transient`). `.a-check.yml` braucht **keine** Änderung — die
Glob-Layer `ports: ["internal/application/port/**"]` und `adapters:
["internal/adapters/**"]` fassen die neuen Dateien bereits.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun (kein NATS, Consumer pollt weiterhin ausschließlich über SQL/API) | kein neuer Fremdkörper, keine neue Abhängigkeit, keine neue Betriebs-Komponente | `LH-FA-SST-007` bleibt unerfüllt; Consumer, die auf niedrige Latenz angewiesen sind, müssen aggressiv pollen (Last auf dem Store) statt ereignisgetrieben zu reagieren |
| B — JetStream (Streams mit Persistenz, Consumer-Gruppen, Replay, At-Least-Once) | eigene Zustellgarantie, Nachholbarkeit direkt aus NATS ohne Rückgriff auf den SQL-Pfad, Konsum-Fortschritt serverseitig verwaltet | baut eine zweite, parallele Nachvollziehbarkeits-Infrastruktur neben dem bereits vorhandenen und bewährten `cdc.changes`-Pfad auf; Zustandshaltung (Retention, Ack-Fenster, Consumer-Gruppen) verlangt Betriebsaufwand und Kapazitätsplanung, den kein Akzeptanzkriterium fordert; zwei Wahrheiten (JetStream-Stream-Position vs. Store) mit ungeklärter Konfliktauflösung; widerspricht der expliziten Lastenheft-Aussage, dass NATS die Nachvollziehbarkeit nicht trägt |
| **C — Core NATS als reines Wecksignal, leerer Payload, `transient`-Fehlerklasse, best-effort nach ACK (gewählt)** | minimale Implementierung (ein dünner Publish-Aufruf); keine neue Zustandshaltung; keine zweite Wahrheit; Notify-Fehlschlag kann den kritischen Pfad strukturell nicht gefährden, weil er nach ACK liegt und sein Fehler nicht propagiert wird; deckt alle drei Akzeptanzkriterien von `LH-FA-SST-007` exakt ab (Happy Path: Signal bei Connected; Boundary: Store bleibt Quelle der Wahrheit; Negative: Reconnect-Nachholen läuft ausschließlich über den bestehenden Weg, den NATS nie ersetzt) | kein Zustellversprechen — ein Consumer, der zwischen zwei Changes kurz getrennt war, verpasst das Signal ersatzlos (ist aber laut Lastenheft explizit zulässig); Consumer muss selbst einen Resync/Poll-Fallback implementieren, den NATS ihm nicht abnimmt |

## Konsequenzen

- Positiv: `LH-FA-SST-007` ist mit minimaler neuer Infrastruktur erfüllbar
  — ein Container, eine Umgebungsvariable, ein dünner Port/Adapter — ohne
  die bestehende Persist-before-ACK-Garantie (`ADR-0011`) anzutasten.
- Positiv: Die einzige Quelle der Wahrheit bleibt der bestehende
  `cdc.changes`-Store; NATS kann nicht zu einer zweiten, divergierenden
  Nachvollziehbarkeits-Quelle werden, weil es strukturell keinen Zustand
  trägt, der divergieren könnte.
- Positiv: Bestehende Deployments (`compose.yaml`, jede heutige
  Env-only-Installation) bleiben ohne jede Änderung lauffähig — additiv,
  kein Breaking Change (analog zu `ADR-0052`s Precedence-Disziplin).
- Negativ: Ohne Zustellgarantie muss jeder NATS-Consumer eigenständig einen
  Resync-Pfad über die SQL-/API-Seite implementieren; NATS nimmt ihm diese
  Verantwortung nicht ab (das ist eine bewusste Entscheidung, keine
  übersehene Lücke).
- Negativ: `github.com/nats-io/nats.go` ist die erste direkte
  Nicht-PostgreSQL-Abhängigkeit — ein neuer Vendor-/Supply-Chain-Faktor,
  der in `go.sum` sichtbar wird.
- Folgepflicht: Ein Slice implementiert den Outbound-Port
  `ChangeNotificationPort`, den Driven-Adapter `natsnotify`, die
  Verdrahtung in `CaptureService` (dritter, optionaler Port-Parameter) und
  `internal/bootstrap` (`CDC_NATS_URL`-Verdrahtung, No-Op bei ungesetzter
  Variable).
- Folgepflicht: `compose.yaml` bekommt einen optionalen NATS-Service samt
  `CDC_NATS_URL`-Wiring für den Integrationstest — Gegenstand des
  umsetzenden Slices, nicht dieser ADR (dieselbe Aufgabenteilung wie in
  `ADR-0052`).
- Folgepflicht: `spec/pflichtenheft.md` §2 erhält `SPEC-017` (Subjekt- und
  Nachrichtenform) und §6 die zugehörige Zeile unter *Externe Verträge* —
  **bereits Teil dieser ADR** (siehe unten), weil Subjekt-Schema und
  Payload-Form hier vollständig entschieden sind, anders als bei
  `ADR-0052`, wo die konkrete Feldform erst im Slice entstand.
- Folgepflicht: `spec/architecture.md` §3 erhält `ARC-013` (NATS als neue
  externe Abhängigkeit) — **bereits Teil dieser ADR**.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/application/usecase/capture`, Whitebox) | Ein fehlschlagender `ChangeNotificationPort` darf den Rückgabewert von `CaptureService.Capture()` nicht zu einem Fehler machen, wenn `store`/`ack` erfolgreich waren — Regressionstest gegen versehentliche Fehlerpropagation | `make test` |
| Go-Unit-Test (`internal/adapters/driven/natsnotify`, Whitebox) | Verbindungsfehler/Publish-Fehler übersetzen in die Klasse `transient` (`errors.Is`-Prüfung gegen einen neuen Sentinel), analog zu `postgresack`s `ErrReplication`-Test | `make test` |
| `.a-check` | unverändert — neue Port-/Adapter-Dateien liegen vollständig innerhalb der bestehenden Glob-Layer `ports`/`adapters`, kein neuer Hexagon-Schichten-Edge nötig | `make a-check` |
| `make test-integration` | optionaler NATS-Service in `compose.yaml`: ein reales Publish/Subscribe belegt Happy Path (`LH-FA-SST-007`); ein Lauf ohne `CDC_NATS_URL` belegt das Boundary-Kriterium (Feature deaktiviert, bestehender Zugriffsweg unverändert) — Umsetzung Gegenstand des Slices | `make test-integration` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Wird jemals eine Anforderung formuliert, dass NATS selbst — ohne Rückgriff
auf den SQL-Lesezugriffsweg — Nachvollziehbarkeit tragen soll (z. B.
Replay über NATS für Consumer, die keinen SQL-Zugriff haben dürfen): dann
widerspricht diese Entscheidung der neuen Anforderung fundamental (Core
NATS kann das nicht) und braucht eine Folge-ADR, die JetStream (Option B)
neu bewertet — nicht als Korrektur, sondern als eigene Entscheidung, weil
sich die Prämisse aus `LH-FA-SST-007` geändert hätte. Sonst permanent — die
Weck-Signal-Semantik dieser ADR gilt unabhängig vom NATS-Server-Versions-
stand (Digest-Hebung ist Modul-14-Routine, keine Re-Evaluierung dieser
Entscheidung).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Architect-Eröffnung der Welle „NATS-Change-Notification" (Roadmap-Trigger `slice-051` <!-- d-check:status-provenance --> liegt in `done/`, erfüllt); `LH-FA-SST-007` bereits im Lastenheft committet, kein bestehender NATS-Bezug im Repo (grüne Wiese) | [`spec/lastenheft.md`](../../../spec/lastenheft.md), [roadmap.md](../planning/in-progress/roadmap.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0055` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
