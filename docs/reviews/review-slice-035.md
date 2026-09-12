# Review-Report: slice-035 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-035`, §1/§2/§3/§6) und
`AGENTS.md` §3 Hard Rules (Modul 10 §Drei Review-Arten). Kein aktives ADR
berührt (Slice-Kopf: „Kein aktives ADR wird geändert").

**Gegenstand:** Commits `aa88ed5` (neue Testfälle `TestMVPHeartbeatHealthy`
in `test/integration/integration_test.go` sowie ein neuer
Verarbeitungsrückstand-Beleg-Abschnitt in
`tools/harness/run-integration-tests.sh`), `9d6dfd8` (DoD-Häkchen,
Plan-Nachzug §3).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-035-black-box-observability-konsolidierung.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken, §8 Register-Sichtung), inklusive Diff `9d6dfd8` gegen den ursprünglichen Plan-Stand (`78d5ae3`…`4f564e9`)
- `spec/lastenheft.md` (`LH-FA-ADM-002`, `LH-FA-ADM-003`, `LH-FA-ADM-004`, `LH-FA-ADM-005`, `LH-FA-SST-002`)
- `test/integration/integration_test.go` (vollständig, insbesondere `TestMVPHeartbeatHealthy`, `awaitHeartbeatErrorClass`, `newMVPEnv`, `mvpSource`, `TestMVPUpdateOldImageWithFullReplicaIdentity`)
- `tools/harness/run-integration-tests.sh` (vollständig: `-run`-Muster, Lasttest-Beleg-Block, Black-Box-CLI-Rundlauf, neuer Rückstands-Beleg-Block, `TestMVPSchemaChangeIncompatibleTypeChange`-Container-Ende-Kommentar)
- `tools/schema/schema.yaml` (`views.consumer_status`-Query, LEFT JOIN gegen `cdc.consumer_position`)
- `internal/bootstrap/wiring.go` (`heartbeatInterval`, `heartbeatStaleAfter`)
- `internal/application/usecase/acknowledge/service.go`, `internal/adapters/driven/postgresstorage/queries/queries.go` (`UpsertConsumerPosition` — bestätigt: `cdc.consumer_position`-Zeile entsteht ausschließlich bei `Acknowledge`, nicht bei `Register`)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7
- `harness/conventions.md` (MR-000 ID-Schema, Modus-Deklaration `PGC`/GF)
- `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/` (Zähler-Stand, `evidence/`), `.../verwaltung-keine-sql-administration/` (im Slice-Plan §8 zitiert)
- `docs/reviews/review-slice-034.md` (Format-Vorlage)

---

## Findings

Keine HIGH/MEDIUM-Findings mit Merge-Relevanz.

- **[LOW]** — `quelle`: Maintainability (`AGENTS.md` §3.7)
  `pfad`: `tools/harness/run-integration-tests.sh:574-575`
  `befund`: Der unveränderte Kommentar vor
  `TestMVPSchemaChangeIncompatibleTypeChange` zählt auf, was den laufenden
  Feed-Container noch braucht und deshalb davor läuft
  („Lasttest-Beleg, Black-Box-CLI-Rundlauf oben") — der neue
  Verarbeitungsrückstand-Beleg-Block (Zeilen 485–567), der ebenfalls den
  laufenden Container braucht und jetzt ebenfalls davor liegt, fehlt in
  dieser Aufzählung. Kein Laufzeitfehler (die Platzierung selbst ist
  korrekt), aber die Aufzählung beschreibt den Ist-Zustand nach diesem
  Diff nicht mehr vollständig.
  `verifizierbar`: nein (kein Sensor auf Kommentar-Vollständigkeit)
  `klasse`: „Aufzählungs-Kommentar nicht mit Diff nachgezogen" (erstes
  Auftreten in diesem Skill-Lauf)

## Negativbefunde

- geprüft, ohne Befund: **Plan-Abweichung bei `LH-FA-ADM-005` ist sachlich
  korrekt begründet und keine Verwässerung — eher eine höhere Treue zum
  Lastenheft-Wortlaut als die ursprüngliche DoD-Formulierung.**
  `tools/schema/schema.yaml:226-247` (`views.consumer_status`) liest
  `latest_commit_position` über eine auf `cp.source_id` korrelierte
  Unterabfrage (`SELECT max(t.commit_position) FROM cdc.transaction t
  WHERE t.source_id = cp.source_id`) hinter einem `LEFT JOIN
  cdc.consumer_position cp ON cp.consumer_id = c.consumer_id`. Eigenständig
  im Code nachvollzogen (nicht aus dem Plan-Nachzug übernommen):
  `internal/adapters/driven/postgresstorage/queries/queries.go`s
  `UpsertConsumerPosition` wird ausschließlich vom
  `Acknowledge`-Anwendungsfall aufgerufen (`internal/application/usecase/acknowledge/service.go`);
  `RegisterConsumer` (`internal/bootstrap/wiring.go:659ff.`) legt nur eine
  `cdc.consumer`-Zeile an, keine `cdc.consumer_position`-Zeile. Für einen
  Consumer, der *nie* bestätigt hat, ist `cp` beim `LEFT JOIN` also `NULL`
  — `cp.source_id` ist `NULL`, die Unterabfrage vergleicht `t.source_id =
  NULL` (nie wahr, SQL-Dreiwertlogik), liefert daher `NULL` statt eines
  Maximalwerts, und `latest_commit_position - acknowledged_position` ist
  `NULL - NULL = NULL` — keine Zahl, erst recht keine Zahl `> 0`. Die
  Behauptung im Plan-Nachzug trägt damit. Wichtiger: `LH-FA-ADM-005`s
  Happy-Path-Wortlaut selbst („Given unbestätigte Changes **jenseits der
  bestätigten Position eines Consumers**") setzt voraus, dass der Consumer
  bereits eine bestätigte Position hat — die ursprüngliche
  DoD-Formulierung („noch nicht bestätigender Consumer") war die losere
  Paraphrase, nicht das Lastenheft-Kriterium selbst. Der real umgesetzte
  Ablauf (erste Bestätigung bindet `source_id` → weitere real erfasste
  Änderung hebt `latest_commit_position` über die bestätigte Position →
  Rückstand sichtbar `> 0` → zweite Bestätigung senkt ihn auf `0`) trifft
  den Wortlaut *exakter* als die ursprüngliche Formulierung. Auch das
  Boundary-Kriterium („Given kein Rückstand, when die Messung gelesen
  wird, then zeigt sie das auch") ist damit real belegt (`backlog_after =
  0` nach der zweiten Bestätigung, real gegen `cdc.consumer_status`
  gelesen). Kein Akzeptanzkriterium von `LH-FA-ADM-005` bleibt durch diese
  Präzisierung unbelegt.
- geprüft, ohne Befund: **Schwelle `age_seconds < 15` für
  `TestMVPHeartbeatHealthy` ist korrekt aus dem Produktionswert
  hergeleitet, nicht neu erfunden.** Eigenständig gegen
  `internal/bootstrap/wiring.go` nachvollzogen: `heartbeatInterval = 5 *
  time.Second` (Zeile 76), `heartbeatStaleAfter = 3 * heartbeatInterval`
  (Zeile 84) ergibt 15s — exakt der im Testfall verwendete
  `freshnessThresholdSeconds`-Wert. Dieselbe Schwelle prüft bereits der
  Compose-Healthcheck (`--healthcheck`-Lauf, Zeile 723/769 in
  `wiring.go`) gegen dieselbe `cdc.heartbeat`-Zeile. Der Testfall
  übernimmt damit einen bereits produktiv bindenden Wert, statt eine
  eigene Zahl zu setzen.
- geprüft, ohne Befund: **Testisolation/Kollisionsfreiheit des neuen
  Bash-Abschnitts.** `tools/harness/run-integration-tests.sh:485-567`
  liegt zwischen dem bestehenden Black-Box-CLI-Rundlauf (`slice-027`,
  endet Zeile 483) und dem Kommentar/Aufruf von
  `TestMVPSchemaChangeIncompatibleTypeChange` (`slice-033`, ab Zeile 569)
  — also *vor* dem Container-Ende, wie es der Slice-Plan §3 verlangt. Die
  neuen IDs `120`/`121` auf `feed_mvp_full` kollidieren mit keiner
  bestehenden ID auf derselben Tabelle: `grep -n "feed_mvp_full\|VALUES
  ("` zeigt bereits vergebene IDs `1` (Go-Testfall
  `TestMVPUpdateOldImageWithFullReplicaIdentity`, läuft im vorgelagerten
  `go test`-Aufruf, blockierend, also abgeschlossen bevor das Bash-Skript
  fortfährt), `90`/`91` (Lasttest-Beleg), `95`/`96` (CLI-Rundlauf) — der
  neue Consumer `cli-e2e-backlog-consumer` ist zudem ein eigener,
  ungenutzter `consumer_id`-Wert. Kein gemeinsamer Zustand wird zwischen
  den Blöcken implizit geteilt außer der bereits vorher aktivierten
  Tabelle selbst.
- geprüft, ohne Befund: **`BEO-PGC/test-runner-stiller-ausschluss` bleibt
  bei 1× — dieser Slice erzeugt keinen neuen Beleg, weil er die Lücke
  vermeidet, nicht wiederholt.** `TestMVPHeartbeatHealthy` ist im
  `-run`-Muster des ersten `go test`-Aufrufs aufgenommen
  (`tools/harness/run-integration-tests.sh:241`,
  `...|TestMVPSchemaChangeAddColumn|TestMVPHeartbeatHealthy)$`) — real im
  Testlauf bestätigt (`--- PASS: TestMVPHeartbeatHealthy` in allen drei
  Review-Läufen, siehe unten). Der neue Bash-Abschnitt für
  `LH-FA-ADM-005` läuft ohnehin unbedingt: Bash-Codeblöcke im
  Runner-Skript unterliegen keinem `go test -run`-Filter — dieser läuft
  ausschließlich auf die beiden `docker run … go test`-Aufrufe, nicht auf
  die dazwischenliegenden reinen Shell-Abschnitte (`register-consumer`,
  `psql`-Inserts, SQL-Lesungen). Ein „stiller Ausschluss" im Sinn dieser
  Beobachtung ist für Bash-Abschnitte strukturell nicht möglich, nur für
  neue Go-Testfunktionen. Register-Verzeichnis
  `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/evidence/`
  enthält weiterhin nur `slice-033.md` — konsistent mit Plan §8
  („1×, weiter offen").
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7), bis auf
  das oben gemeldete LOW-Finding.** Der neue Doc-Kommentar auf
  `TestMVPHeartbeatHealthy` (`integration_test.go:747-760`) beschreibt den
  geltenden Testzweck im Indikativ, referenziert auflösbare Anker
  (`LH-FA-ADM-002`, `tools/schema/nacharbeit-heartbeat.sql`,
  `LH-QA-OPS-002`, `internal/bootstrap/wiring.go`) und die
  Kopplungs-Begründung zur Testreihenfolge (läuft vor
  `TestMVPSchemaChangeIncompatibleTypeChange`, weil jener `error_class`
  dauerhaft setzt) — zulässige Kopplungs-Klasse, keine Chronik einer
  verworfenen Alternative. Der neue Bash-Kommentar
  (`run-integration-tests.sh:485-498`) beschreibt ebenfalls den geltenden
  Mechanismus (LEFT-JOIN-/NULL-Verhalten) im Indikativ, keine
  Vorwärtsverweise auf andere Slices, kein Abbruch mitten im Satz. Kein
  Kommentar in diesem Diff behauptet einen abwesenden Text oder eine
  verworfene Alternative.
- geprüft, ohne Befund: **DoD-Aktualisierung ohne Overclaiming.** Die
  sechs in `9d6dfd8` auf `[x]` gesetzten Punkte sind durch den Diff und
  durch eigene Testläufe in dieser Sitzung gedeckt: beide Testfälle
  existieren und liefen in drei unabhängigen `make test-integration`-Läufen
  grün; `make gates` lief real und grün; kein öffentlicher Vertrag berührt
  (reine Testdatei- und Runner-Skript-Änderung); das
  Reconciliation-Register-Item ist korrekt als „entfällt" markiert (Repo
  laut `harness/conventions.md` durchgehend GF, keine
  `reconciliation.md` vorhanden). Die vier offen gelassenen Punkte
  (Review, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei
  Paarungen) sind korrekt unbeansprucht — Planner-Closure-Arbeit nach
  Modul 8, kein Self-Review. Die im DoD-Text erwähnten „vier grünen
  Läufe plus zwei gezielt rot geführte Mutationsläufe" sind eine
  Implementer-Eigenaussage („siehe Bericht an den Reviewer") und wurden in
  dieser Review-Sitzung nicht erneut reproduziert (Mutationsläufe sind
  kein DoD-Pflichtnachweis); die drei geforderten grünen
  `make test-integration`-Läufe wurden unabhängig, real und vollständig
  nachvollzogen (siehe unten).
- geprüft, ohne Befund: **Plan-Nachzug §3 ehrlich benannt.** Der Diff
  `9d6dfd8` beschreibt zutreffend zwei Implementer-Entscheidungen (Bash
  statt Go-Test wegen fehlendem Docker-Socket im Toolchain-Container;
  Präzisierung des Rückstandsszenarios) als das, was sie sind —
  Design-Entscheidungen ggü. der ursprünglichen Formulierung, nicht
  rückwirkende Umdeutung des Slice-Ziels. Der Verweis auf die ursprüngliche
  wörtliche Formulierung bleibt im DoD-Text sichtbar stehen
  („Präzisierung … nicht die wörtliche „niemals bestätigt"-Situation"),
  statt sie stillschweigend zu ersetzen.
- geprüft, ohne Befund: **Kein Image-Rebuild-Trigger.** Der Diff berührt
  ausschließlich `test/integration/integration_test.go` und
  `tools/harness/run-integration-tests.sh` — keine Datei unter `internal/`
  oder `cmd/`. `harness/image-hash.txt` bleibt unangetastet, konsistent
  mit `ADR-0044`/`harness/README.md` §Werkzeuge.
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `LH-FA-ADM-002`/`LH-FA-ADM-005`, kein `SPEC-*`/`ARC-*` im Betreff;
  `make commit-traceability` lief in dieser Sitzung über `HEAD~5..HEAD`
  grün (0 Befunde, „Betreffs ohne Struktur-ID").
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — der neue
  Bash-Abschnitt läuft ausschließlich über `docker exec`/`docker run`
  gegen bereits laufende Compose-Container, kein lokales Toolchain-Install.
  **3.2** (Suppression-Verbot) — kein `#noqa`/`//nolint`/`[SuppressMessage]`
  im Diff. **3.3** (git-mv/Inhalt-Trennung) — nicht einschlägig, keine
  Umbenennung in diesem Diff. **3.5**/**3.6** — kein ADR im Diff berührt,
  keine Gate-Schwelle gelockert.
- geprüft, ohne Befund: **§8 Sub-Area-Sichtung konsistent mit dem
  Register-Bestand.** Beide im Plan §8 zitierten Beobachtungsverzeichnisse
  (`BEO-PGC/test-runner-stiller-ausschluss`,
  `BEO-PGC/verwaltung-keine-sql-administration`) existieren im
  Beobachtungs-Register mit dem im Plan genannten Zähler-Stand (1× bzw.
  0×/benannt). Modus-Begründung „alle berührten Sub-Areas GF (nur
  `*`/`PGC`)" ist korrekt — kein neuer Code-Pfad, reine Testinfrastruktur
  gegen bestehende, bereits dokumentierte SQL-Sichten.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify` (54 Dateien), `docs-check` (290 Dateien, 0 Befunde),
  `commit-traceability` (5 Commits, `HEAD~5..HEAD`, OK), `a-check` (0
  Befunde) — alle grün.
- geprüft, ohne Befund: **`make test-integration`, real dreimal in Folge
  ausgeführt** (diese Review-Sitzung, nicht aus dem Implementer-Bericht
  übernommen) — `TestMVPHeartbeatHealthy` PASS in allen drei Läufen
  (0.01s), alle acht Go-Testfälle des ersten Compose-`go test`-Aufrufs
  PASS, `TestMVPSchemaChangeIncompatibleTypeChange` PASS im jeweils
  letzten Aufruf; der neue Verarbeitungsrückstand-Beleg meldete in allen
  drei Läufen identisch „Rückstand vor der zweiten Bestätigung 1464,
  danach 0"; `git status --porcelain` am Ende dieser Sitzung leer (bis auf
  das erwartete, zurückgesetzte `tools/schema/plan.yaml`-Rollout-Artefakt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Aufzählungs-Kommentar nicht mit Diff
nachgezogen" (1×, erstes Auftreten — kein Steering-Loop-Eintrag fällig,
Schwelle ist 3×).

## Verdikt

**Merge-blockierend:** keins. Das einzige Finding ist LOW und
semantik-frei (die Platzierung des neuen Bash-Abschnitts ist korrekt,
nur eine bestehende Kommentar-Aufzählung wurde nicht um den neuen Block
ergänzt).

**Zur Plan-Abweichung (zentrale Prüffrage dieses Laufs):** Die
Abweichung von der ursprünglichen, wörtlichen DoD-Formulierung („niemals
bestätigt" → Rückstand > 0) ist sachlich korrekt begründet und
eigenständig gegen `tools/schema/schema.yaml` sowie den Go-Code-Pfad von
`register-consumer`/`acknowledge-consumer` nachvollzogen worden, nicht
aus dem Plan-Nachzug übernommen. Sie ist keine Verwässerung von
`LH-FA-ADM-005`: Der real umgesetzte Ablauf trifft den
Lastenheft-Wortlaut des Happy-Path-Kriteriums exakter als die
ursprüngliche DoD-Paraphrase, und das Boundary-Kriterium bleibt
unverändert real belegt.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** die SQL-/Go-Code-Herleitung der `NULL`-Behauptung
(`schema.yaml`, `queries.go`, `wiring.go`, `service.go`), die
Herleitung der 15s-Schwelle aus `wiring.go`, die ID-Kollisionsfreiheit
über den vollständigen Bestand von `run-integration-tests.sh` und
`integration_test.go`, die strukturelle Begründung, warum Bash-Abschnitte
keinem `-run`-Filter unterliegen, und `make gates`/`make test-integration`
(real dreimal) als eigenständige Läufe in dieser Sitzung.

**Übergabe:** Ein LOW-Finding, kein Rollen-Konflikt, keine
Architect-Sequenz nötig (Modul 8; die Sequenz greift erst ab HIGH mit
Rollen-Widerspruch oder ab dem dritten gleichen Konflikttyp). Dieser
Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat.
