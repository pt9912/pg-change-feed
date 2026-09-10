# Slice slice-012: Health-Endpoint per Heartbeat

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-4.

**Bezug:** [`LH-FA-ADM-002`](../../../../spec/lastenheft.md), [`LH-QA-OPS-002`](../../../../spec/lastenheft.md), [`ADR-0024`](../../../../docs/plan/adr/README.md), [`ADR-0027`](../../../../docs/plan/adr/README.md), [`ADR-0046`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`ARC-002`](../../../../spec/architecture.md), [`ARC-006`](../../../../spec/architecture.md)
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-10.

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

**Ziel:** Health-Endpoint real: der Capture-Prozess schreibt periodisch
seinen Lebenszeichen-Zustand in eine `cdc.process_heartbeat`-Tabelle, eine
vierte SQL-Lese-View macht ihn automatisiert abfragbar (Architect-Verdikt
[`docs/plan/adr/architect-review-slice-011.md`](../../adr/architect-review-slice-011.md):
kein neuer Driving-Adapter, keine neue ADR nötig).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- HTTP-/gRPC-Health-Endpoint — [`ADR-0020`](../../adr/0020-http-grpc-optional.md)
  bleibt unberührt; der SQL-Kanal genügt für automatisierte Health Checks.
- Rollen-spezifische Verdrahtung (`BEO-PGC/rollen-verdrahtung`) — anderer
  Vorgang, Sicherheits-Thema; der Heartbeat-Schreiber nutzt die bestehende
  Instanz-DSN wie alle anderen Adapter dieses Bootstraps.
- Fehlerzustände über den Heartbeat hinaus (Erfassungs-Fehlerklassen,
  Stream-Abbrüche) — Folge-Slice slice-013 übernimmt das; der Heartbeat
  trägt nur Liveness, keine Fehlerdetails.

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

- [ ] Heartbeat-Schreiber im Capture-Prozess (periodischer Timer,
      Bootstrap-Verdrahtung) — Teil-Beleg zu
      [`LH-FA-ADM-002`](../../../../spec/lastenheft.md).
- [ ] `cdc.process_heartbeat`-View + `make test-store`-Beleg — Teil-Beleg
      zu [`LH-QA-OPS-002`](../../../../spec/lastenheft.md) (automatisierte
      Health Checks über SQL).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `compose.yaml` (Healthcheck liest den Heartbeat statt
      nur den Prozess-Start), falls berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/schema.yaml` | update | `process_heartbeat`-Tabelle in die d-migrate-Kette je [`ADR-0043`](../../../../docs/plan/adr/README.md) |
| `internal/application/port/outbound/heartbeat.go` | neu | Outbound-Port für den periodischen Schreib-Zug (Driven-Adapter-Fähigkeit je [`ADR-0024`](../../../../docs/plan/adr/README.md)) |
| `internal/adapters/driven/postgresstorage/heartbeat.go` | neu | PostgreSQL-Adapter, schreibt Zeitstempel + Prozess-Kennung |
| `internal/bootstrap/wiring.go` | update | Timer-Zug im Capture-Prozess (Application-/Bootstrap-Schicht je [`ADR-0027`](../../../../docs/plan/adr/README.md)) |
| `tools/schema/nacharbeit-heartbeat.sql` | neu | vierte SQL-Lese-View `cdc.heartbeat` (reine Projektion, [`ADR-0046`](../../../../docs/plan/adr/README.md)) — Ausweichform wie die bestehenden drei Views (`BEO-PGC/d-migrate-nacharbeit`) |
| `compose.yaml` | update | Healthcheck liest den Heartbeat-Zustand |

**Plan-Nachzug im selben Lauf (Baseline-Regelwerk `modul-09-implementierung.md`, seit slice-009):** über die obige Liste hinaus berührte dieser Lauf folgende Dateien — vor dem Sensor-Lauf nachgetragen, kein stilles Erweitern:

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | `UpsertHeartbeat`-SQL-Text — die SQL-Texte der postgresstorage-Adapter stehen ausschließlich hier (Paket-Kommentar), keine neue Komponente, direkte Abhängigkeit von `postgresstorage/heartbeat.go` |
| `internal/adapters/driven/postgresstorage/heartbeat_test.go` | neu | Adapter- und View-Test-Beleg (`make test-store`) für die zwei planmäßigen DoD-Punkte (Schreiber, View) |
| `internal/bootstrap/heartbeat_internal_test.go` | neu | Whitebox-Beleg (`package bootstrap`) für das §6-Risiko (Timer-Zug blockiert die Capture-Persist-ACK-Schleife nicht) und für die Healthcheck-Schwelle (`healthcheckVerdict`) |
| `Makefile` | update | fünfter `schema-rollout`-Schritt (`nacharbeit-heartbeat.sql`), direkte Folge der neuen Nacharbeit-Datei |
| `tools/schema/nacharbeit-observability.sql` | update | Kommentar korrigiert — die dort als „offene Architekturfrage" benannte Health-Endpoint-Lücke (Slice-Plan slice-011 §7) ist mit diesem Slice aufgelöst; ein stehengebliebener Kommentar wäre nach `AGENTS.md` §3.7 falsch |
| `spec/pflichtenheft.md` | update | [`SPEC-001`](../../../../spec/pflichtenheft.md)-Tabellenzeile: die dort bereits vorgesehene, nie detaillierte `cdc.capture_state`-Zeile („Betriebs-/Capture-Zustand") wird durch die jetzt realisierte `cdc.process_heartbeat` ersetzt — sonst zwei Namen für denselben Zweck (Fund beim Plan-vs-Bestand-Abgleich, Schritt 12 des Implementer-Workflows) |
| `cmd/pg-change-feed/main.go` | update | `--healthcheck`-Modus — der einzig ausführbare Compose-Healthcheck-Befehl im distroless Runtime-Image (kein Shell, kein `psql`, `Dockerfile`); ohne ihn bliebe der geplante `compose.yaml`-DoD-Punkt eine Doku-Behauptung ohne Wirkung |
| `tools/harness/run-integration-tests.sh` | update | Docker-Health-Status-Wartepunkt (`docker inspect .State.Health.Status`) — Beleg des Compose-Healthcheck-Vertrags am realen Container, nicht nur am Binary-Exit-Code |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-4 eröffnet, kein anderes Slice in
`in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Timer-Zug im
  Capture-Prozess plus die vierte SQL-View plus der Compose-Healthcheck-
  Umbau sprengen drei Liefer-Punkte — Rückzug mit Zerlegung (Schreiber
  getrennt von View/Healthcheck).
- `in-progress` → `open` (blockiert — Carveout?): der periodische Timer
  kollidiert mit der bestehenden Capture-Loop-Struktur und ist ohne
  größeren Umbau nicht sauber einzuhängen — Blocker, Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss ohne
  offenes HIGH-Finding; der Heartbeat wird am realen Compose-Healthcheck
  belegt (beobachtbar am Integrationstest-Diff).
- Lerneintrag §7: geschärfte Regel oder benannte Spec-Lücke — ohne ihn
  bleibt der Slice abgelegt, nicht fertig.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Heartbeat aus `BEO-PGC/health-endpoint-heartbeat` war in welle-3 als
  Umsetzungs-Frage benannt (Architect-Verdikt: kein Folge-ADR nötig) —
  **Ausgang:** eingetreten (Träger: dieser Slice).
- Der Timer-Zug im Capture-Prozess könnte die Persist-before-ACK-Ordnung
  ([`LH-QA-REL-001`](../../../../spec/lastenheft.md)`.a`) stören, wenn er in
  derselben Goroutine läuft — muss nachweislich unabhängig laufen. Wird bei
  Closure bewertet.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührten Sub-Areas
(Application/Outbound-Port, Store-Adapter, Bootstrap) sind Segmente des
GF-Baums der Modus-Deklaration (`PGC` Greenfield, Doc führt) — jede erfüllt
die Schwelle ≥ 2 von 3 Achsen: Strukturregeln im Spec-Stratum committet,
Phase-Reife 3–4, Evidenz-/Diskrepanz-Risiko niedrig. Keine Sub-Area zu grob.

**Vorgelagert — offene Beobachtungen sichten:** Register gelesen (neun
Einträge zum Zeitpunkt der Planung): `health-endpoint-heartbeat` 1× (dieser
Slice löst sie ein), `d-migrate-nacharbeit` 2× (berührt — vierte
Nacharbeit-Datei, Ausweichform trägt weiter, unter der Schwelle),
`rollen-verdrahtung` 1× (nicht berührt, siehe §1-Ausschluss). Übrige sechs
ohne Bezug zu diesem Slice. Kein Eintrag erreicht mit diesem Slice 3× —
keine Lücke.

**Modus-Begründungsblock — Umfang.** Reiner GF-Hinweis genügt (siehe oben);
kein Sub-Area-Block.
