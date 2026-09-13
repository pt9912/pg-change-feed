# Slice slice-038: CLI-Diagnose-Befehl

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-12 — dritter und letzter Slice, schließt
`BEO-PGC/verwaltung-keine-sql-administration` ab (siehe `welle-12.md` §1/§3).

**Bezug:**
[`LH-FA-SST-003`](../../../../spec/lastenheft.md) (CLI),
[`LH-FA-ADM-002`](../../../../spec/lastenheft.md) (Betriebsstatus),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbare Fehlerzustände),
[`LH-FA-ADM-004`](../../../../spec/lastenheft.md) (messbarer CDC-Abstand),
[`LH-FA-ADM-005`](../../../../spec/lastenheft.md) (sichtbarer
Verarbeitungsrückstand).

**Berührte Spec-Stellen:** — (reine CLI-Ergänzung auf bereits bestehenden
SQL-Lese-Views, keine Verhaltensänderung an einer Architektur-Sicht-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein neuer CLI-Diagnose-Befehl (`cmd/pg-change-feed/main.go`, Muster
identisch zum bestehenden `--healthcheck`) liest über `cfg.ReaderDSN`
(`ADR-0047`) die bereits bestehenden, `cdc_reader`-gegrantenen SQL-Lese-Views
`cdc.heartbeat` (`tools/schema/nacharbeit-heartbeat.sql`) und `cdc.metrics`
(`tools/schema/nacharbeit-observability.sql`) und gibt eine
menschenlesbare Zusammenfassung aus, die `LH-FA-SST-003`s Boundary
("deckt mindestens die Status-/Diagnoseabfragen ab, die `LH-FA-ADM-002`
… `005` nennen") real erfüllt:
- **`LH-FA-ADM-002`** (Betriebsstatus) und **`LH-FA-ADM-003`** (sichtbare
  Fehlerzustände) — aus `cdc.heartbeat` (`age_seconds`, `error_class`).
- **`LH-FA-ADM-004`** (messbarer CDC-Abstand) — aus `cdc.metrics`,
  Zeile `cdc_capture_lag`.
- **`LH-FA-ADM-005`** (sichtbarer Verarbeitungsrückstand) — aus
  `cdc.metrics`, Zeilen `cdc_consumer_lag{consumer}` (je Consumer).

Kein neuer Domänentyp, kein neuer Port: alle vier Werte stehen bereits als
Lese-Views bereit (`LH-FA-SST-002`, seit `welle-11`/früher); dieser Slice
liefert ausschließlich den fehlenden CLI-Zugriffsweg darauf.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue SQL-Views oder -Spalten** — `cdc.heartbeat`/`cdc.metrics` decken
  alle vier `LH-FA-ADM-00N`-Signale bereits vollständig ab (siehe oben);
  dieser Slice ist reiner Lesezugriffsweg, keine Datenmodell-Änderung
  (Schicht-Abgrenzung: CLI-Schicht, nicht Schema-Schicht).
- **Maschinenlesbares Ausgabeformat (JSON o. ä.)** — `LH-FA-SST-003`
  verlangt nur, dass die Abfragen "über eine CLI erfolgen" können, kein
  bestimmtes Format; `cdc.metrics` selbst ist bereits das
  maschinenlesbare Long-Format (`LH-FA-SST-004`) und über SQL direkt
  erreichbar (`LH-FA-SST-002`) — ein zweites, CLI-eigenes
  Maschinenformat wäre eine Doppelung ohne neuen Anforderungs-Träger.
  Bleibt Bestand: der bestehende SQL-Zugriffsweg auf dieselben Views.
- **Schwellenwert-Interpretation (healthy/degraded/unhealthy,
  `SPEC-007`)** — dieselbe bewusste Nicht-Entscheidung wie bei
  `cdc.heartbeat`/`cdc.metrics` selbst (siehe deren Kommentare): die
  Views liefern Rohwerte, keine Klassifikation; der neue Befehl gibt sie
  unverändert weiter, trifft keine neue Schwellenwert-Entscheidung, die
  ein ADR bräuchte (`AGENTS.md` §3.6).
- **Ein Datenbankzugriff außerhalb `cdc_reader`s Grant-Fläche**
  (z. B. `pg_stat_replication`, Relationsgrößen) — dieselbe bewusste
  Nicht-Abdeckung wie in `nacharbeit-observability.sql` dokumentiert
  (`cdc_wal_retention_bytes`, `cdc_storage_bytes` u. a.); ein anderer
  Vorgang, der die Least-Privilege-Fläche erweitern müsste.

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

- [x] Neuer CLI-Befehl (`cmd/pg-change-feed/main.go`, z. B. `diagnose` —
      Implementer entscheidet den exakten Namen und begründet ihn im
      Plan-Nachzug) liest `cdc.heartbeat` + `cdc.metrics` über
      `cfg.ReaderDSN` und gibt `age_seconds`/`error_class`
      (`LH-FA-ADM-002`/`003`), `cdc_capture_lag`
      (`LH-FA-ADM-004`) und `cdc_consumer_lag{consumer}` je Consumer
      (`LH-FA-ADM-005`) menschenlesbar aus.
- [x] `LH-FA-SST-003` real erfüllt: ein Integrationstest (Ergänzung in
      `tools/harness/run-integration-tests.sh`, analog zum bestehenden
      `--healthcheck`-Testabschnitt) ruft den neuen Befehl gegen den
      laufenden Compose-Feed-Container per `docker exec` auf und
      bestätigt, dass alle vier Signale in der Ausgabe erscheinen —
      sowohl im Normalbetrieb als auch (mindestens für ADM-003) in einem
      erkennbar von Normalbetrieb unterscheidbaren Fehlerzustand
      (Boundary-Kriterium von `LH-FA-ADM-002`/`003`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-038.md`](../../../reviews/review-slice-038.md)
      (0 HIGH, 2 MEDIUM, 1 LOW), Fixrunde behoben in Commit `95db5a5`,
      bestätigt in
      [`docs/reviews/review-slice-038-fixrunde.md`](../../../reviews/review-slice-038-fixrunde.md)
      (dabei neuer Nebenbefund F-4 LOW, in Commit `af00c0c` behoben).
      Verifikation in
      [`docs/reviews/verify-slice-038.md`](../../../reviews/verify-slice-038.md)
      (DoD eigenständig nachgeprüft, alle vier `LH-FA-ADM-002`…`005`-Signale
      real bestätigt, keine Rückführung nötig).
- [x] Doku-Update: `README.md`/`harness/README.md`, falls dort die
      vorhandenen CLI-Befehle aufgezählt sind (Implementer prüft und
      begründet im Plan-Nachzug, ob eine Stelle existiert, die den neuen
      Befehl nennen muss).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht (siehe §3 Plan-Nachzug).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — `BEO-PGC/adapter-fehler-ausgang` um `evidence/slice-038.md` ergänzt (2×, Zähler in Commit `af00c0c` nachgezogen).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — beide *entfallen*.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo mit Wellen-Betrieb (`welle-12`) — Prüfung läuft bei der `welle-12`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `cmd/pg-change-feed/main.go` | update | neuer Subcommand, Muster identisch zu `--healthcheck`/`register-consumer` |
| `internal/bootstrap/wiring.go` | update | neue Funktion (Muster `Healthcheck`) — öffnet kurzlebige `ReaderDSN`-Verbindung, liest `cdc.heartbeat` + `cdc.metrics`, formatiert die Ausgabe |
| `tools/harness/run-integration-tests.sh` | update | realer E2E-Beleg des neuen Befehls gegen den laufenden Feed-Container |
| `README.md` oder `harness/README.md` | update, falls zutreffend | Implementer prüft, ob eine Stelle die vorhandenen CLI-Befehle aufzählt |

### Plan-Nachzug (nach Code)

**Exakter Befehlsname:** `diagnose` (ohne führende `--`, wie `register-consumer`/
`acknowledge-consumer` — kein Infra-Flag wie `--healthcheck`/`--version`,
sondern ein admin-facing Sondermodus ohne Argumente:
`pg-change-feed diagnose` bzw. `docker exec <container> /pg-change-feed diagnose`).

**Implementierung** (`internal/bootstrap/wiring.go`, Funktion `Diagnose`,
Muster identisch zu `Healthcheck`): kurzlebige `pgxpool`-Verbindung über
`cfg.ReaderDSN`, `context.WithTimeout(3s)`, liest `cdc.heartbeat`
(`age_seconds`, `error_class`) und zwei `cdc.metrics`-Abfragen
(`cdc_capture_lag`, alle `cdc_consumer_lag`-Zeilen). Anders als
`Healthcheck` trifft der Befehl keine binäre Verdikt-Entscheidung — Ausgang
0 heißt „Lesezugriff erfolgreich", unabhängig vom Berichtsinhalt (ein
gemeldeter Fehlerzustand oder Rückstand ist kein Befehlsfehler); Ausgang 1
nur bei Verbindungs-/Query-Fehler (DSN ungültig, Instanz nicht erreichbar,
View nicht lesbar).

**Abweichung vom Plan — NULL-sicherer Scan für `cdc_consumer_lag`:** Der
`test-store`-Lauf deckte real auf, dass `cdc.metrics.value` für eine
`cdc_consumer_lag`-Zeile `NULL` sein kann, obwohl `WHERE
cs.acknowledged_position IS NOT NULL` filtert — die
`latest_commit_position`-Unterabfrage in `cdc.consumer_status` liefert
`NULL`, wenn die an die Position gebundene Quelle noch nie eine Transaktion
trug; die Differenz `NULL - x` bleibt `NULL`. Ursprünglich als `float64`
geplant, scannt der Consumer-Zeilen-Loop jetzt über `*float64` und gibt in
diesem Fall „unbekannt (Quelle trug noch nie eine Transaktion)" statt eines
irreführenden Rückstands aus. Dieselbe defensive Formulierung („nur
Consumer mit mindestens einer bestätigten Position") schließt zugleich
Risiko 2 aus §6: ein nie bestätigender Consumer erscheint gar nicht erst,
und der Text sagt das ausdrücklich, statt sein Fehlen als „kein Rückstand"
lesbar zu lassen.

**Risiko 1 aus §6 (realer Fehlerzustand im laufenden Container):**
`reportFault` (`internal/bootstrap/wiring.go`) schreibt `error_class` nur
unmittelbar vor `os.Exit` — jeder von `Run()` klassifizierte Fehler beendet
den Prozess, `restart: "no"` hält den Container danach beendet stehen; ein
`docker exec` gegen einen bereits beendeten Container ist nicht mehr
möglich. Ein genuin vom Erfassungspfad ausgelöster Fehlerzustand lässt sich
an einem **laufenden** Container deshalb strukturell nicht per CLI
beobachten, ohne den Container zu beenden — das bestätigt die im Risiko
vorab benannte Schwierigkeit. Sowohl der Unit-Test
(`internal/bootstrap/diagnose_test.go`,
`TestDiagnoseReportsErrorState`) als auch der E2E-Beleg
(`tools/harness/run-integration-tests.sh`) schreiben denselben Spaltenwert
(`cdc.process_heartbeat.error_class`), den `reportFault` im realen Fehlerfall
schriebe, direkt über SQL — derselbe Lesepfad (View → CLI-Ausgabe), ohne den
laufenden Prozess zu beenden. Der E2E-Beleg schreibt den Fehlerzustand in
einer bis zu 20 Versuche kurzen Schleife, weil der periodische
Heartbeat-Takt (5s) ihn beim nächsten erfolgreichen Beat unabhängig davon
wieder auf `NULL` zurücksetzt — 20 Versuche liegen weit innerhalb eines
einzigen 5s-Takts und sind kein Flakiness-Kompromiss, sondern eine
Sicherheitsmarge gegen die Taktphase. **Ausgang siehe §6, Risiko 1** (in §6
oben ist es das erste der beiden Risiken).

**Zweite reale Erkenntnis beim E2E-Beleg:** `cdc.consumer_status.
latest_commit_position` ist **quellenweit**, nicht consumer- oder
tabellenspezifisch (`SELECT max(commit_position) FROM cdc.transaction WHERE
source_id = cp.source_id`, unabhängig von der Tabelle). Die ursprünglich im
Testskript angenommene Erwartung „`CLI_CONSUMER` zeigt Rückstand 0" war
deshalb falsch, sobald nach seiner letzten Bestätigung eine weitere
Transaktion auf `src-mvp` läuft (hier: die `BACKLOG_CONSUMER`-Belege auf
derselben Tabelle) — der E2E-Abschnitt prüft für `CLI_CONSUMER` jetzt nur
noch generisch eine numerische Zeile, für `BACKLOG_CONSUMER` (dessen zweite
Bestätigung unmittelbar davor lief) weiterhin exakt 0.

**Doku-Update — Ort geprüft, nicht README.md/harness/README.md selbst:**
Weder `README.md` noch `harness/README.md` zählen die vorhandenen
CLI-Befehle auf (geprüft: kein Treffer für `register-consumer`,
`--healthcheck` o. ä. in `README.md`; `harness/README.md` nennt sie nur
innerhalb der `make test-integration`-Sensor-Beschreibung). Die tatsächliche
„Bedienung"-Stelle ist `docs/user/benutzerhandbuch.md` (von `README.md` §Was
kann ich heute tun? verlinkt) — sie führt `register-consumer`,
`acknowledge-consumer` und `--healthcheck` mit Beispielaufrufen. Aktualisiert:
neue Untersektion „Diagnose ausführen" (§4, zwischen „Metriken lesen" und
„WAL-Rückstand prüfen"), die `cdc_reader`-Zeile und die
`CDC_READER_DSN`-Zeile, Versionsfeld 1.1→1.5 und Änderungshistorie-Eintrag
1.5 (die Versions-/Datums-Kopfzeile hatte bereits vor diesem Slice hinter der
Änderungshistorie zurückgelegen — mit diesem Eintrag korrigiert).
`harness/README.md`s `make test-integration`-Zeile zusätzlich um den neuen
CLI-Diagnose-Beleg ergänzt (`seit slice-038`).

**Image neu gebaut:** `main.go`/`wiring.go` sind Build-Kontext-Dateien
(`harness/README.md` §Sensors, `make image`-Zeile) — `make image` gelaufen,
`harness/image-hash.txt` trägt den neuen Digest, committet im selben Zug.

**Reconciliation-Register (§2-DoD-Punkt):** entfällt — dieses Repo führt
kein `docs/plan/planning/reconciliation.md` (Greenfield, kein
Brownfield-Bootstrap, `harness/conventions.md` §Modus-Deklaration).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-036` und `slice-037` liegen in
`done/` (unabhängig davon laut `welle-12` §5 auch parallel startbar,
tatsächlich aber sequenziell nach diesem Slice-Zuschnitt), `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  der neue Befehl plus sein E2E-Beleg zusammen mehr als drei Liefer-Punkte
  oder mehr als zwei Schichten in einer Review-Sitzung nicht mehr prüfbar
  machen, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün (inkl. des neuen CLI-Diagnose-Belegs) **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein realer Fehlerzustand (`LH-FA-ADM-003`s Boundary: „unterscheidbar von
  normalem Betrieb") im laufenden Compose-Stack künstlich zu erzeugen,
  könnte schwieriger sein als angenommen (bisherige Fehlerklassen-Tests,
  z. B. `TestMVPSchemaChangeIncompatibleTypeChange`, brauchen einen
  eigenen, isolierten Feed-Container-Lauf) — **Ausgang: entfallen.** Die
  Befürchtung bestätigte sich in einer schärferen Form als angenommen: ein
  genuin vom Erfassungspfad ausgelöster Fehlerzustand ist an einem
  **laufenden** Container strukturell nicht per `docker exec` beobachtbar,
  weil `reportFault` nur unmittelbar vor `os.Exit` schreibt (§3
  Plan-Nachzug). Die eigentliche Sorge hinter dem Risiko — dass das
  DoD-Item „Fehlerzustand real erkennbar" damit unerfüllbar bliebe —
  entfällt aber, weil der wörtliche DoD-Wortlaut nur einen „erkennbar von
  Normalbetrieb unterscheidbaren Fehlerzustand" verlangt, nicht dessen
  reale Auslösung durch den Erfassungspfad: derselbe Spaltenwert
  (`cdc.process_heartbeat.error_class`) und derselbe Lesepfad (View → CLI)
  ohne Prozessende erfüllen das. Der Reviewer hat diesen Ersatzbeleg
  eigenständig geprüft und als „akzeptabel, keine unzulässige Verwässerung"
  bestätigt (`review-slice-038.md`). Kein Carveout, kein Folge-Slice nötig.
- `cdc.metrics`s `cdc_consumer_lag`-Zeilen existieren nur für Consumer mit
  mindestens einer bestätigten Position (`WHERE cs.acknowledged_position
  IS NOT NULL`); ein frisch registrierter, noch nie bestätigender
  Consumer erscheint dort nicht — die CLI-Ausgabe könnte das
  fälschlich als „kein Rückstand" statt „noch nie gemessen" lesen lassen.
  **Ausgang: entfallen.** Behoben durch explizite Ausgabe-Formulierung
  („nur Consumer mit mindestens einer bestätigten Position") statt eines
  Silent-Fallbacks, real durch `TestDiagnoseReportsNoConfirmedConsumer`
  (Fixrunde-Ergänzung, Commit `95db5a5`, vom Reviewer eigenständig
  reproduziert bestätigt) abgesichert.

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

- **Was hat funktioniert:** Der Implementer fand von sich aus den
  NULL-sicheren Scan-Bugfix für `cdc_consumer_lag` und die reale
  `latest_commit_position`-Quellenweit-statt-Consumer-spezifisch-Erkenntnis
  während der E2E-Belegarbeit — beide real vor dem Review, nicht erst
  danach gefunden. Die Wiederverwendung des `Healthcheck`-Musters
  (kurzlebige `ReaderDSN`-Verbindung, Timeout, klare Fehlerklassen) trug
  ohne Anpassung. Der Reviewer traf eine klare, begründete
  Bewertungsentscheidung zur Ersatzbeleg-Frage (§6 Risiko 1), statt sie
  offen zu lassen — das machte die Closure-Entscheidung eindeutig.
- **Was ging anders als geplant:** Der Reviewer fand 2 MEDIUM (fehlende
  Regressionstests für drei reale `Diagnose`-Zweige; ein übersehener
  Beobachtungs-Register-Treffer `BEO-PGC/adapter-fehler-ausgang`) und
  1 LOW (Risiko-Nummerierung im Plan-Nachzug widersprach §6), alle drei in
  der Fixrunde behoben. Die Fixrunden-Bestätigung selbst fand einen
  weiteren, kleinen Nebenbefund F-4 (Register-`state.md` zeigte nach der
  Fixrunde weiterhin den alten 1×-Zähler statt 2×) — bei der Closure
  direkt nachgezogen, kein weiterer Rollenwechsel nötig für eine reine
  Zähler-Korrektur im Register selbst.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× —
  `BEO-PGC/adapter-fehler-ausgang` steht jetzt bei 2× (`slice-007`,
  `slice-038`), `BEO-PGC/verwaltung-keine-sql-administration` bleibt bei
  0× (siehe unten).
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-038.md` in `BEO-PGC/adapter-fehler-ausgang/` ergänzt —
  Zähler steht bei 2×, dieselbe strukturelle Eigenschaft („jeder
  klassifizierte Fehler ist terminal") erneut bestätigt.
- **Folge-Slices:** keine — dies ist der letzte Slice von `welle-12`;
  ihre Closure ist der nächste Schritt.
- **Risiken aus §6:** beide *entfallen* — siehe §6 für Begründung.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-12`) — Prüfung
  läuft bei der `welle-12`-Closure.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/verwaltung-keine-sql-administration` (0×,
benannt nicht gezählt — dieser Slice liefert den dritten und letzten
Baustein der Auflösung, siehe `welle-12` §3 Closure-Trigger),
`BEO-PGC/github-actions-unverifizierbar-lokal` (1×, thematisch nicht
berührt von diesem Slice — reine CDC-Fähigkeit, kein CI/CD-Bezug) und
`BEO-PGC/adapter-fehler-ausgang` (1× seit `slice-007` — thematisch
einschlägig, nachträglich ergänzt nach `review-slice-038` F-3: der
Plan-Nachzug oben bestätigt real dieselbe strukturelle Eigenschaft, dass
jeder von `Run()` klassifizierte Fehler terminal ist, als Grund, warum §6
Risiko 1 nur per SQL-Ersatzbeleg statt eines realen, prozess-ausgelösten
Fehlerzustands testbar ist; Beleg `evidence/slice-038.md` nachgetragen,
Zähler damit bei 2× — Schwelle 3× noch nicht erreicht, kein
Steering-Loop-Eintrag fällig). Keiner der Treffer erreicht mit diesem
Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
