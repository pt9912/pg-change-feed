# Welle 15 — NATS-Change-Notification — Closure-Notiz

**Welle:** welle-15
**Abschluss:** 2026-09-14
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-SST-007`](../../../../spec/lastenheft.md) (NATS-Change-Notification)
  ist jetzt vollständig erfüllt — alle drei Akzeptanzkriterien real belegt:
  **Happy Path** (`slice-053`/`058`): ein Consumer, der
  `cdc.changes.<source_id>.<schema>.<table>` abonniert, empfängt real ein
  leeres Wecksignal nach jeder Change dieser Tabelle. **Boundary**
  (`slice-054`): eine Change, die entsteht, während kein Consumer
  verbunden ist, bleibt vollständig über `cdc.changes` lesbar, der
  Feed-Container läuft unbeeinflusst weiter. **Negative** (`slice-055`):
  nach einer realen Verbindungstrennung (`docker network disconnect`)
  und Wiederverbindung holt ein Consumer verpasste Changes ausschließlich
  über den bestehenden SQL-Zugriffsweg nach — Core NATS liefert nichts
  nach, real durch Log-Abwesenheitsbeleg bestätigt.
- Grundgerüst (`slice-052`, [ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md)):
  neuer Outbound-Port `ChangeNotificationPort`, Driven-Adapter
  `natsnotify` (Core NATS, kein JetStream, leerer Payload,
  `transient`-Fehlerklasse, best-effort nach `ACK Source`), `CDC_NATS_URL`
  vollständig optional.
- Compose-Verdrahtung (`slice-053`): digest-gepinnter NATS-Service
  (`nats:2-alpine@sha256:065e8355…`), Healthcheck über den
  HTTP-Monitor-Port 8222, `CDC_NATS_URL`-Wiring für den Feed-Container.
- Korrektur-Slice (`slice-058`, [ADR-0056](../../adr/0056-nats-tabellen-granulares-subjekt.md),
  `Supersedes ADR-0055` Punkt 2): das Subjekt-Schema wurde von
  `cdc.changes.<source_id>` (ein Subjekt je Quelle) auf
  `cdc.changes.<source_id>.<schema>.<table>` (tabellen-granular)
  korrigiert — reale Nutzeranforderung während der laufenden Welle,
  Konflikt-Pfad-Verdikt 2 (Folge-ADR mit `Supersedes`, keine stille
  Lockerung). Wildcard-Erhalt (`cdc.changes.<source_id>.>`) bleibt für
  „alle Tabellen einer Quelle"-Consumer erhalten; die bisherige exakte
  Drei-Token-Subscription funktioniert danach bewusst nicht mehr
  (benannter Breaking Change).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der Referenz-Adapter `postgresack` trug als Kopiervorlage für
  `natsnotify` die gesamte Konstruktions-/Options-/Fehler-Wrapping-Form
  fast unverändert — beide Slices, die einen neuen NATS-Adapter oder
  seine Erweiterung bauten (`slice-052`, `slice-058`), profitierten
  davon.
- Die Functional-Option-Erweiterung von `CaptureService` (`WithChangeNotification`)
  hielt alle fünf bestehenden Aufrufstellen unverändert kompilierbar —
  kein einziger Rückführungsfall (`in-progress → next`) trat in der
  gesamten Welle ein, obwohl mehrfach als Risiko benannt.
- Der eigenständige, netzlose NATS-Testcontainer
  (`tools/harness/run-notify-tests.sh`, `make test-notify`) und die
  drei aufeinander aufbauenden E2E-Testabschnitte in
  `run-integration-tests.sh` (Happy-Path/Boundary/Negative) trugen ein
  konsistentes, wiederverwendbares Muster (`natssub`-Wegwerf-Client,
  `docker logs`-Polling auf `READY`/`RECEIVED`).
- Reviewer und Verifier haben in jedem der fünf Slices dieser Welle
  mindestens einen der berichteten Mechanismen (NATS-Monitor-Endpunkte,
  `docker inspect`-Netzwerkstatus, Wildcard-Subscriptions) **selbst,
  unabhängig voneinander** gegen frische Wegwerf-Container reproduziert
  — nie nur die Implementer-Behauptung übernommen.
- Der vorgezogene Architect-Zug für `ADR-0056` (während `slice-053` noch
  lief) verhinderte, dass zwei weitere, noch ungeplante Slices
  (`slice-054`/`055`) mit dem inzwischen überholten Subjekt-Schema
  geplant worden wären — der Plan-Nachzug geschah, bevor sie aktiviert
  wurden, nicht danach.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Der Auftraggeber (pt9912) bemängelte während `slice-053`s laufender
  Implementierung real das ursprüngliche, quellen-weite Subjekt-Schema
  aus `ADR-0055` — ein Consumer mit Tabellen-Interesse hätte jedes
  Signal der ganzen Quelle erhalten. Das führte zu einer echten
  Architektur-Korrektur mitten in der Welle: `ADR-0056` (`Supersedes
  ADR-0055` Punkt 2) und dem zusätzlichen Korrektur-Slice `slice-058`,
  der ursprünglich nicht Teil der Welle-Planung war.
- `slice-055`s ursprünglich geplanter `/subsz`-Bestätigungsweg für die
  reale Verbindungstrennung erwies sich als zu träge (NATS' eigene
  Dead-Connection-Erkennung ist ping-basiert, nicht sofort) — ein echter
  roter Zwischenstand während der Implementierung, behoben durch den
  Wechsel auf `docker inspect`s synchronen Netzwerk-Status.
- Der Godoc-Kommentar über `natsnotify.New` (`slice-052`) trug
  versehentlich eine Slice-Chronik — das vierte Auftreten der bereits
  3×-verkörperten Klasse `BEO-PGC/slice-chronik-in-code-kommentar`,
  behoben per Fixrunde und mit einem vorgezogenen Architect-Zug
  bewertet (Rollenteilung Implementer-Selbstprüfung/Reviewer-
  Sicherheitsnetz explizit verkörpert, siehe Steering-Loop-Einträge).
- `slice-053`s Commit führte die `CDC_NATS_URL`-Doku-Zeile ohne
  Versionshistorie-Nachzug ein — das dritte Auftreten von
  `BEO-PGC/handbuch-versionshistorie-uebersprungen`, erst bei `slice-058`s
  Implementierung entdeckt und per Architect-Zug verkörpert.
- Über die gesamte Welle hinweg trat dreimal ein Pipe-/Wrapper-Exit-Code-
  Fallstrick auf (`make gates`/`make test-integration`, deren Fehlschlag
  durch `| tail`, `| grep` oder einen Hintergrund-Task-Wrapper maskiert
  wurde) — einmal traf es den Planner selbst real bis zum Push. Per
  Architect-Zug als neue Hard Rule `AGENTS.md` §3.9 verkörpert.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **`.harness/skills/reviewer.md`** geschärft: ein eigener, benannter
  HIGH-Unterpunkt „Slice-/Wellen-Chronik in Produktionscode-Kommentar"
  ersetzt die bisher implizite Subsumtion unter „Kommentar trägt keine
  der Kommentar-Klassen"; zusätzlich trägt `.claude/commands/implement-slice.md`
  Schritt 20 seither eine Grenz-Klarstellung (Implementer-Selbstprüfung
  ist erste, nicht tragende Verteidigungslinie)
  — liegt in `.harness/skills/reviewer.md` (HIGH-Liste) und
  `.claude/commands/implement-slice.md` (Schritt 20).
  Auslöser: `BEO-PGC/slice-chronik-in-code-kommentar` (`slice-041`,
  `slice-041`-Fixrunde, `slice-044`, `slice-052` — 4×, vorgezogener
  Architect-Zug während dieser Welle).
- **`.harness/skills/reviewer.md`** und **`.claude/commands/implement-slice.md`**
  geschärft: ein neuer HIGH-Punkt „Handbuch-Versionshistorie nicht
  fortgeschrieben" (tragende Linie, Reviewer) sowie eine Implementer-
  Selbstprüf-Instruktion (Schritt 17, erste, nicht tragende Linie)
  stellen sicher, dass eine inhaltliche Änderung an
  `docs/user/benutzerhandbuch.md` künftig immer mit `Version:`-Kopf und
  Änderungshistorie-Zeile im selben Diff einhergeht
  — liegt in `.harness/skills/reviewer.md` (HIGH-Liste) und
  `.claude/commands/implement-slice.md` (Schritt 17).
  Auslöser: `BEO-PGC/handbuch-versionshistorie-uebersprungen`
  (`slice-045`, `slice-046`, `slice-053` — 3×, vorgezogener Architect-Zug
  während dieser Welle).
- **`AGENTS.md`** geschärft: neue Hard Rule §3.9 „Exit-Code eines
  Gate-Laufs wird direkt geprüft, nie durch eine Pipe/einen Wrapper
  hindurch" (mit Falsch/Richtig-Beispiel), plus Kurzverweis in §6
  Schritt 6 und in `harness/README.md` §Minimal agent workflow Schritt 6
  — liegt in `AGENTS.md` §3.9.
  Auslöser: `BEO-PGC/pipe-maskiert-make-exit-code` (`slice-058`,
  `slice-054`, `slice-055` — 3×, Lese-Schritt dieser Welle-Closure).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. In
dieser Welle **neu angelegt**: `nats-notify-validierung-koennte-bestand-ablehnen`
(1×, `slice-058`, Ausgang *weiter offen*), `pipe-maskiert-make-exit-code`
(3×, verkörpert `seit welle-15`, siehe Steering-Loop-Einträge). In dieser
Welle **auf 3× bzw. 4× gehoben** (Belege aus vorherigen Wellen plus einem
Beleg aus dieser Welle, beide per vorgezogenem Architect-Zug verkörpert,
siehe Steering-Loop-Einträge): `slice-chronik-in-code-kommentar` (4×,
Beleg `slice-052`), `handbuch-versionshistorie-uebersprungen` (3×, Beleg
`slice-053`, entdeckt bei `slice-058`). Unverändert, nicht von dieser
Welle berührt: alle übrigen Registereinträge aus `welle-14` und früher.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — `welle-15` schließt vollständig mit ihren fünf Slices (`slice-052`,
`053`, `054`, `055`, `058`); keine Fortsetzung wurde als Folge-Slice
angelegt. `LH-FA-SST-007` ist vollständig erfüllt, `ADR-0056`s
Re-Evaluierungs-Trigger (feinere Filterung als Tabellenebene) ist nicht
eingetreten.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle fünf Slices (`slice-052`, `053`, `054`, `055`, `058`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure, Commit `0a5cfbb` und danach
  unverändert, Exit-Code explizit geprüft, nicht durch eine Pipe
  maskiert — Coverage-Gate 40,90 % ≥ 35 %, `d-check` 437 Dateien/0
  Befunde, `commit-traceability` OK, `a-check` 0 Befunde).
- `make test-integration` real grün über alle drei NATS-Testabschnitte
  (Happy-Path, Boundary, Negative) — real von Implementer, Reviewer und
  Verifier je unabhängig reproduziert, mit dem korrigierten
  vier-Token-Subjekt (`slice-058`).
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): kein
  offener Carveout im Repo (`docs/plan/carveouts/` nur `.gitkeep`). Kein
  bootstrap-aware Gate Teil dieser Welle (NATS ist kein Gate). `ADR-0055`s
  Re-Evaluierungs-Trigger (Anforderung, dass NATS selbst ohne
  SQL-Rückgriff Nachvollziehbarkeit trägt) und `ADR-0056`s
  Re-Evaluierungs-Trigger (Filterung feiner als Tabellenebene) sind
  beide nicht eingetreten — beide ADRs bleiben unverändert `Accepted`.
- Drei Paarungen (Anker · Folge-Slice · Register): **Anker** — drei
  Steering-Loop-Einträge mit `liegt in`, alle drei Zielorte
  (`.harness/skills/reviewer.md`, `.claude/commands/implement-slice.md`,
  `AGENTS.md` §3.9) existieren real. **Folge-Slice** — keiner genannt,
  nichts zu prüfen. **Register** — alle in dieser Welle berührten
  Verzeichnisse (`nats-notify-validierung-koennte-bestand-ablehnen`,
  `pipe-maskiert-make-exit-code`, `slice-chronik-in-code-kommentar`,
  `handbuch-versionshistorie-uebersprungen`) existieren mit nicht leerem
  `evidence/` — grün.

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; alle fünf Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
