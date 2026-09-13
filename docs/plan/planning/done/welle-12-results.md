# Welle 12 — Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose — Closure-Notiz

**Welle:** welle-12
**Abschluss:** 2026-09-13
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  real umgesetzt: eine Antrags-Queue (`cdc.administration_request`) mit
  `LISTEN`/`NOTIFY`, verarbeitet von einer neuen Administrations-Goroutine
  im laufenden Capture-Prozess (`slice-037`).
- `LH-FA-ADM-001` (SQL-Administration, Aktivierung/Deaktivierung):
  `cdc.enable_table`/`cdc.disable_table` als `SECURITY DEFINER`-SQL-Funktionen,
  `SET search_path = cdc, pg_temp`, ausschließlich über `cdc_admin`
  ausführbar (`slice-036`).
- `LH-FA-CFG-002` (Live-Deaktivierung): ein echter Live-Zugriffsweg
  existiert jetzt — `SELECT cdc.disable_table(...)` beendet die Erfassung
  einer Tabelle am **laufenden** Prozess, ohne Neustart, real gegen den
  Compose-Stack belegt (`slice-037`).
- Boot-Wechsel: `Assembler.tables` wird beim Start jetzt aus
  `cdc.source_table` über `TableActivationPort.List` gebaut, nicht mehr
  ausschließlich aus `CDC_TABLES` — `CDC_TABLES` bleibt Seed für die
  Erstaktivierung (`slice-037`).
- `LH-FA-SST-003` (CLI-Diagnose): neuer `diagnose`-Subcommand
  (`cmd/pg-change-feed/main.go`) deckt `LH-FA-ADM-002`…`005` real ab
  (Betriebsstatus, sichtbare Fehlerzustände, messbarer CDC-Abstand,
  sichtbarer Verarbeitungsrückstand) über die bestehenden
  `cdc.heartbeat`/`cdc.metrics`-Views (`slice-038`).
- `spec/architecture.md`s Sequenzdiagramm zu `LH-FA-CFG-001.a` in
  CLI-synchronen und SQL-asynchronen Pfad getrennt (`slice-037`,
  `ADR-0050`s Folgepflicht).
- `BEO-PGC/verwaltung-keine-sql-administration` (real fehlende
  Fähigkeiten, gefunden bei `welle-11`s Eröffnungs-Recherche) ist damit
  aufgelöst — siehe unten.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die zweite Fork-Recherche vor dem Schneiden der Slices (Eröffnung dieser
  Welle) fand real, dass `ADR-0046` die SQL→Go-Port-Brücke bewusst offen
  gelassen hatte und `Assembler.tables` keinen Reload-Mechanismus besaß —
  beide Fragen wurden vor der Implementierung durch `ADR-0050` (Architect-
  Entscheidung, Option D: Antrags-Queue + Live-Reload) entschieden, statt
  während der Implementierung ad hoc gelöst zu werden.
- `ADR-0050`s eigene Fitness Function (`go test -race` gegen
  `Assembler.AddBinding`/`RemoveBinding`) erwies sich als diskriminierend
  genau geplant: Implementer, Reviewer und Verifier reproduzierten
  unabhängig voneinander denselben roten Kontrollfall (Sperren entfernt →
  reale `DATA RACE`) und denselben grünen Zustand (Sperren vorhanden).
- `slice-036`s d-migrate-`POST_EXECUTE_DRIFT`-Grenze für SQL-Funktionen
  wurde sauber über die etablierte `nacharbeit-*.sql`-Ausweichform
  aufgefangen, statt den Rollout zu erzwingen — dieselbe Disziplin wie bei
  den bereits gelösten Fällen (CHECK-Constraint, Views).
- `slice-038` fand von sich aus zwei reale Probleme (NULL-sicherer
  `cdc_consumer_lag`-Scan, `latest_commit_position`s tatsächliche
  Quellenweite statt Consumer-/Tabellenspezifität) vor dem Review, nicht
  erst danach.
- Der Reviewer traf bei `slice-038`s §6-Risiko 1 (Ersatzbeleg für einen
  strukturell unbeobachtbaren realen Fehlerzustand) eine klare, begründete
  Bewertungsentscheidung statt die Frage offen zu lassen — das machte die
  Closure-Entscheidung eindeutig.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Ein Implementer-Dispatch für `slice-037` scheiterte einmal an einem
  echten Netzwerkfehler (DNS/API-Erreichbarkeit) mitten im Lauf — `git
  status`/`git log` bestätigten einen sauberen, unveränderten Stand, der
  identische Auftrag wurde erfolgreich neu gestartet. Kein Code- oder
  Prozessfehler, reine Infrastruktur-Störung.
- `slice-037`s Reviewer wurde durch einen parallel laufenden `EnterPlanMode`-
  Aufruf (orthogonales CI/CD-Vorhaben) kurzzeitig pausiert und schrieb
  seinen Zwischenstand in eine Plan-Datei statt fertigzustellen — nach
  `ExitPlanMode` sauber mit demselben Kontext fortgesetzt, kein
  Arbeitsverlust.
- Zwei kleine Nachzieh-Lücken traten bei den Closures auf: `slice-037`s
  eigene „Review durchgeführt"-Checkbox blieb nach dem Review-Abschluss
  zunächst `[ ]` (vom Verifier als Finding V-1 markiert, bei der Closure
  nachgezogen); `slice-038`s Beobachtungs-Register-`state.md` zeigte nach
  der Fixrunde weiterhin den alten 1×-Zähler statt 2× (Reviewer-
  Bestätigungslauf F-4, direkt behoben). Beide sind dieselbe Klasse wie
  bereits verkörpert in `BEO-PGC/dod-checkbox-nachzug` — keine neue
  Beobachtung, unter der bestehenden Schwelle unverändert.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle 3× — der Normalfall.
`BEO-PGC/verwaltung-keine-sql-administration` wird unten als eigener,
direkter Auflösungsfall (nicht über die 3×-Schwelle) behandelt.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. In
dieser Welle **verkörpert** (direkt, nicht über die 3×-Schwelle, analog zu
`schema-evolution-nicht-dynamisch`): `verwaltung-keine-sql-administration`
(0×, seit welle-12). In dieser Welle **fortgeschrieben**:
`adapter-fehler-ausgang` (1× → 2×, `slice-038` bestätigt dieselbe
strukturelle Eigenschaft erneut). In dieser Welle **neu angelegt, direkt
gestrichen** (wellenlos, aus einer Nutzerfrage zum ADR-Index, nicht Teil
der Welle-Lieferung selbst): `architect-verdikt-ablageort-uneinheitlich`
(0×, real behoben). Neu angelegt (wellenlos, orthogonal zur CI/CD-Arbeit,
nicht Teil dieser Welle): `github-actions-unverifizierbar-lokal` (1×).
Unter der Schwelle, unverändert von dieser Welle:
`dod-checkbox-nachzug-architect-pfad` (1×), `retention-keine-loeschausfuehrung`
(0×, benannt nicht gezählt), `rollen-test-abdeckungsluecken` (2×),
`schema-rollout-fremdobjekte` (1×), `test-isolation-geteilter-zustand` (1×),
`test-runner-stiller-ausschluss` (1×), `walsender-wirksamkeit` (1×).
Bereits verkörpert, unverändert: `a-check-null-abdeckung` (3×),
`d-migrate-nacharbeit` (5×, `slice-036` fügte den fünften Beleg für die
weiterhin offene dritte Objektklasse — SQL-Funktionen — hinzu, ohne den
bereits verkörperten Ausgang zu ändern), `dod-checkbox-nachzug` (3×),
`lese-doppelquelle` (3×, `seit slice-029`), `plan-nachzug` (2×),
`plan-vorlagen-defekt` (3×), `schema-evolution-nicht-dynamisch` (2×,
`seit slice-033`). Bereits eingetreten, unverändert: `cdc-capture-lag-real`
(2×), `health-endpoint-heartbeat` (2×), `rollen-verdrahtung` (4×),
`spec008-replication-luecke` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — `welle-12` schließt vollständig mit ihren drei Slices; keine
Fortsetzung wurde als Folge-Slice angelegt.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices (`slice-036`, `slice-037`, `slice-038`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure, Commit `c70a236` und
  danach unverändert).
- Alle in `welle-12` §3 genannten Closure-Kriterien real erfüllt:
  `SELECT cdc.enable_table(...)` real ausgeführt → Administrations-
  Goroutine verarbeitet den Antrag → eine danach ausgeführte Änderung wird
  vom laufenden Prozess real erfasst, ohne Neustart
  (`tools/harness/run-integration-tests.sh`, Abschnitt „SQL-Administration
  Live-Reload-Beleg (enable)"); derselbe Nachweis für
  `cdc.disable_table(...)` (Abschnitt „… (disable)" — Erfassung endet,
  Feed-Container läuft unverändert weiter). Beide dreifach in Folge grün
  reproduziert (`slice-037`s Verifikation). Der `diagnose`-CLI-Befehl
  liefert real alle vier `LH-FA-ADM-002`…`005`-Signale
  (`slice-038`, eigenständig vom Verifier reproduziert).
  `BEO-PGC/verwaltung-keine-sql-administration` erreicht Ausgang
  *verkörpert* (siehe oben).
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig". Kein offener Carveout im Repo
  (`docs/plan/carveouts/` leer). Kein bootstrap-aware Gate berührt.
  `ADR-0050`s und `ADR-0046`s Re-Evaluierungs-Trigger (technische
  SQL→Port-Brücke wie FDW/`dblink`) sind nicht eingetreten — beide bleiben
  unverändert `Accepted`.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — kein
  Steering-Loop-Eintrag mit `liegt in` in dieser Welle, nichts zu prüfen.
  Folge-Slice — keiner genannt, nichts zu prüfen. Register — die in dieser
  Welle verkörperte `BEO-PGC/verwaltung-keine-sql-administration` und die
  fortgeschriebene `BEO-PGC/adapter-fehler-ausgang` existieren beide als
  Verzeichnis mit nicht leerem `evidence/` (Verwaltung: 0 Belege, bereits
  im eigenen `observation.md` unter „Benannt, nicht gezählt" transparent
  ausgewiesen — direkte Auflösung ohne 3×-Beleg, keine stille Lücke;
  Adapter-Fehler-Ausgang: 2 Belege real vorhanden).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; alle drei Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
