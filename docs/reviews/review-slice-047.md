# Review-Report: slice-047 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-047`, §1/§2/§3/§4/§6/§7/§8),
`slice-045`/`slice-046` (Vorgänger, Referenz für View-Semantik und
Testmuster) sowie `AGENTS.md` §3 Hard Rules (§3.1, §3.3, §3.5, §3.7) —
Rollentrennung Modul 8: diese Prüfung läuft gegen Plan/ADR/Hard Rules
(Maintainability), nicht gegen DoD (Verifier-Aufgabe).

**Gegenstand:** Commit `8a0ccee` (`feat(diagnose): Retention-Sichtbarkeit in
CLI-Diagnose (LH-FA-SST-003)`) gegen Elter `860e7c0` (isolierter Tree-Diff
über `git show 8a0ccee --stat`, nicht der volle Bereich ab `c6ba18f`, der
den unabhängigen Commit `860e7c0` mitträgt); Diff:
`internal/bootstrap/wiring.go`, `tools/harness/run-integration-tests.sh`,
`docs/user/benutzerhandbuch.md`, `harness/image-hash.txt`,
`docs/plan/planning/in-progress/slice-047-diagnose-cli-retention-sichtbarkeit.md`
(Plan-Nachzug, DoD-Häkchen, §6/§7).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft
2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-047-diagnose-cli-retention-sichtbarkeit.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug vier
  Punkte, §4 Trigger, §6 Risiken mit Ausgang, §7 Closure-Notiz, §8)
- `docs/plan/planning/done/slice-045-blockierende-consumer-sichtbarkeit.md`,
  `.../slice-046-cdc-storage-bytes-metrik.md`, `.../welle-13.md`
  (Vorgänger-Slices; Referenz für View-Verträge und Testmuster)
- `docs/plan/planning/open/slice-048-retention-lebenszyklus-kombinierter-e2e-rundlauf.md`
  (Isolations-Prüfung gegen künftige Überschneidung)
- `tools/schema/schema.yaml` (§`retention_blockers`, §`consumer_status`),
  `tools/schema/nacharbeit-observability.sql` (§`cdc.metrics`-Definition),
  `tools/schema/nacharbeit-roles.sql` (§`cdc_reader`-Grants)
- `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/`,
  `.../BEO-PGC/test-isolation-geteilter-zustand/`,
  `.../BEO-PGC/slice-chronik-in-code-kommentar/` (Register-Einträge, real
  gegen `observation.md`/`state.md`/`evidence/` geprüft)
- `spec/lastenheft.md` `LH-FA-SST-003`, `LH-FA-RET-005`, `LH-FA-RET-006`
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git-mv + Inhalt = zwei Commits —
  hier nicht einschlägig, kein `git mv` in diesem Commit), §3.5
  (ADR-Immutabilität), §3.7 (Kommentar-Disziplin —
  `BEO-PGC/slice-chronik-in-code-kommentar`, 3×/verkörpert vor diesem
  Lauf, besonders geprüft)
- `git log -1 --format='%s' 8a0ccee` (Commit-Traceability)
- `make docs-check` (lokal ausgeführt: 368 Dateien, 0 Befunde)
- `make commit-traceability` (lokal ausgeführt gegen `c6ba18f..8a0ccee`
  bzw. `HEAD~5..HEAD`: OK, keine Struktur-ID im Betreff)
- `make a-check`, `make gates` (lokal ausgeführt: 0 Befunde, grüner
  Nachweis-Stempel)
- `grep -nE '^\+.*\b(slice-[0-9]{3}|welle-[0-9]+)\b'` über den
  vollständigen `.go`/`.sh`-Diff dieses Commits (kein Treffer)
- `compose.yaml` (§Werkzeuge-Prüfung: kein `build:`-Block, bestätigt die
  in §7 geschilderte Ursache des ersten roten `make test-integration`-Laufs)

---

## Findings

### F-1 — Beispiel-Ausgabe im Benutzerhandbuch zeigt für denselben Consumer zwei widersprüchliche Rückstandswerte

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Dokumentationskorrektheit; kein Hard-Rule-Verstoß)
- `pfad`: `docs/user/benutzerhandbuch.md:451-454`
- `befund`: Die neue „Ausgabe (Beispiel)“ zeigt für denselben Consumer
  `cli-e2e-consumer` in derselben Momentaufnahme `cdc_consumer_lag: 0`
  (Zeile 452) und zugleich `Blockierender Consumer … Rückstand 3` (Zeile
  453). Laut Schema berechnen beide Werte für denselben Consumer/dieselbe
  Quelle dieselbe Formel — `cdc_consumer_lag` aus `cdc.consumer_status`
  (`tools/schema/schema.yaml:277-281`: `latest_commit_position -
  acknowledged_position`) und `retention_blockers.backlog`
  (`tools/schema/schema.yaml:337-341`: `max(commit_position) -
  acknowledged_position`) sind für denselben Consumer gegen dieselbe
  Quelle identisch. Die erfundenen Beispielzahlen (0 vs. 3) widersprechen
  sich damit intern, ohne dass der Fließtext eine Erklärung für die
  Differenz liefert — ein Leser könnte daraus fälschlich schließen, beide
  Metriken maßen etwas Verschiedenes.
- `verifizierbar`: nein — kein Gate prüft den semantischen Inhalt einer
  Beispiel-Ausgabe in Prosa-Dokumentation.
- `klasse`: „Beispiel-Ausgabe widerspricht der eigenen Schema-Formel“

### F-2 — Beispiel-Ausgabe widerspricht dem im selben Dokument festgehaltenen Namens-/Kennungs-Invariant

- `kategorie`: LOW
- `quelle`: Maintainability (Dokumentationskorrektheit)
- `pfad`: `docs/user/benutzerhandbuch.md:453` vs. `docs/user/benutzerhandbuch.md:278`
- `befund`: Zeile 278 desselben Dokuments hält fest: „`<name>` trägt
  zugleich Kennung und Namen des Consumers“ (`register-consumer`ist der
  einzige dokumentierte Registrierungsweg). Die neue Beispielzeile 453
  zeigt trotzdem einen abweichenden Namen `CLI E2E Consumer` neben der
  Kennung `cli-e2e-consumer` — für einen über die CLI registrierten
  Consumer strukturell nicht erreichbar. Geringe Auswirkung (rein
  illustrativ, keine Prosa-Aussage hängt daran), aber ein Widerspruch
  innerhalb desselben Dokuments.
- `verifizierbar`: nein — kein Gate prüft Beispielinhalte gegen
  Fließtext-Invarianten.
- `klasse`: „Beispiel-Ausgabe widerspricht dokumentiertem Invariant“

## Negativbefunde

- geprüft, ohne Befund: Kommentar-Disziplin (§3.7,
  `BEO-PGC/slice-chronik-in-code-kommentar`, vor diesem Lauf
  3×/verkörpert) — vollständiger `grep`-Lauf gegen
  `slice-[0-9]{3}`/`welle-[0-9]+` über alle neuen/geänderten Zeilen in
  `internal/bootstrap/wiring.go` und `tools/harness/run-integration-tests.sh`:
  kein Treffer. Alle neuen Kommentare (Go-Godoc, Shell-Blockkommentare)
  beschreiben den Ist-Zustand mit auflösbaren `LH-*`-Herkunfts-Ankern,
  keine Chronik, kein Verweis auf verworfene Alternativen. Die einzige
  im Diff sichtbare `slice-038`-Erwähnung in
  `run-integration-tests.sh:20-22` ist unveränderter Kontext (Teil eines
  bereits vor diesem Slice bestehenden Kommentarblocks), keine
  Neueinführung durch diesen Commit.
- geprüft, ohne Befund: docs-check-ID-Link-Falle — `make docs-check` real
  ausgeführt (368 Dateien, 0 Befunde); die neue `Bezug:`-Zeile in
  `slice-047-…md` (`` [`LH-FA-SST-003`](…) ``, backtick-innerhalb-Link)
  folgt exakt demselben, bereits in `slice-045`/`046` etablierten und
  gate-grünen Muster — keine neue Verletzung. Spätere Fließtext-Nennungen
  derselben IDs ohne Link (Benutzerhandbuch-Changelog, Prosa) sind laut
  realem Gate-Lauf zulässig.
- geprüft, ohne Befund: „Kein Blocker“-Randfall (`pgx.ErrNoRows`) —
  `internal/bootstrap/wiring.go:1246-1250` behandelt die Abwesenheit
  einer Zeile in `cdc.retention_blockers` explizit als eigenen
  `switch`-Zweig mit Text „kein Blocker …“, kein Fehlerausgang; dieselbe
  Lesart wie der bestehende `cdc.heartbeat`-Zweig unmittelbar darüber in
  derselben Funktion. Kein struktureller Unterschied zu `slice-038`s §6
  Risiko 2 bzw. `slice-045`s §6 Risiko 1.
- geprüft, ohne Befund: reales rotes Gegenbeispiel (§7 Closure-Notiz) —
  `compose.yaml` trägt real keinen `build:`-Block
  (`compose.yaml:5,44-45`, mit erklärendem Kommentar), das beschriebene
  Szenario (Compose-Stack führt noch das alte
  `ghcr.io/pt9912/pg-change-feed:dev`-Image, bis `make image` läuft) ist
  damit strukturell plausibel und deckt sich mit
  `harness/README.md` §Werkzeuge; kein Hinweis auf ein
  Persistenz-/Cache-Problem stattdessen.
- geprüft, ohne Befund: `harness/image-hash.txt` — real geändert
  (`sha256:68b3859f…` → `sha256:c303e5a4…`, gültiges 64-Zeichen-Format),
  Build-Kontext (`wiring.go`) ist in demselben Commit geändert, Commit
  enthält beide Änderungen zusammen (kein separater, vergessener
  Digest-Commit).
- geprüft, ohne Befund: bestehendes `diagnose`-Ausgabeformat — Diff zeigt
  ausschließlich angehängte neue Zeilen nach dem bestehenden
  `cdc_consumer_lag`-Block; keine bestehende `LH-FA-ADM-002`…`005`-Zeile
  wurde umformuliert oder verschoben, sowohl im Go-Code als auch in den
  zugehörigen Shell-Assertions von `run-integration-tests.sh` (Zeilen
  644-654 bleiben unverändert, neue Prüfungen sind zusätzliche
  `if`-Blöcke danach).
- geprüft, ohne Befund: Commit-Traceability — `git log -1 --format=%s
  8a0ccee` trägt `LH-FA-SST-003`, keine `SPEC-*`/`ARC-*`-Struktur-ID im
  Betreff; `make commit-traceability` real ausgeführt (`c6ba18f..8a0ccee`
  und `HEAD~5..HEAD`): beide OK.
- geprüft, ohne Befund: DoD-Checkboxen (§2) — der Implementer-Bericht
  „bereits selbst nachgezogen“ trifft real zu: alle Punkte außer „Review
  durchgeführt“ und „drei Paarungen“ sind in diesem Diff auf `[x]`
  gesetzt und mit konkreten Beleg-Verweisen versehen; kein Auftreten der
  `BEO-PGC/dod-checkbox-nachzug`-Klasse.
- geprüft, ohne Befund: Beobachtungs-Register-Sichtung (§8) — beide
  zitierten Einträge (`BEO-PGC/test-runner-stiller-ausschluss`,
  `BEO-PGC/test-isolation-geteilter-zustand`) existieren real mit
  jeweils genau einer `evidence/`-Datei (1×, unter der Schwelle); die
  Einschätzung „kein direkter Bezug“ für Ersteren ist zutreffend (kein
  neuer `-run`-gefilterter Go-Testfall entsteht, nur zusätzliche
  Shell-Assertions gegen einen bereits erfassten `docker exec
  diagnose`-Aufruf) und für Letzteren ebenfalls (die Wiederverwendung des
  bestehenden `CLI_CONSUMER`-Zustands ist sequenzielle Positionierung in
  einem einzelnen Compose-Lauf, kein paralleler `postgresstorage`-
  Paket-Testzustand, den die Beobachtung tatsächlich beschreibt).
- geprüft, ohne Befund: §6 Risiken — beide Risiken tragen den Ausgang
  „entfallen“ mit nachvollziehbarer, real durch den beschriebenen
  Testlauf gedeckter Begründung; keine der drei Ausgangsklassen falsch
  zugeordnet.
- geprüft, ohne Befund: `cdc_reader`-Rollenfläche — `cdc.retention_blockers`
  trägt das `SELECT`-Grant an `cdc_reader` bereits seit `slice-045`
  (`tools/schema/nacharbeit-roles.sql:105`, unverändert in diesem Diff);
  keine neue Rollenfläche nötig, Kommentarbehauptung in `wiring.go`
  („alle drei Views“) trifft zu.
- geprüft, ohne Befund: Überschneidung mit `slice-048` — dessen Plan
  (`open/slice-048-…md` §1) schließt CLI-Diagnose-Erweiterungen explizit
  als bereits von `slice-047` abgedeckt aus und beschränkt sich auf die
  reine SQL/`psql`-Ebene; die neuen Abschnitte in
  `run-integration-tests.sh` sind mit eigenen, benannten
  Kommentarüberschriften („Retention-Sichtbarkeits-Beleg (CLI), Zustand
  1/2“) klar abgegrenzt und stehen an einer Stelle (vor der ersten
  `register-consumer`-Zeile bzw. im bestehenden Normalbetrieb-Block), die
  einem künftigen kombinierten SQL-Abschnitt nicht im Weg steht.
- geprüft, ohne Befund: Schicht-Grenzen (`AGENTS.md` §3.4) — keine
  Berührung, `spec/architecture.md` wird in diesem Diff nicht geändert;
  `make a-check` real ausgeführt: 0 Befunde.
- geprüft, ohne Befund: ADR-Immutabilität (§3.5) — keine ADR-Datei in
  diesem Diff geändert.
- geprüft, ohne Befund: Docker-only (§3.1) — keine lokale
  Toolchain-Installation in diesem Diff.
- geprüft, ohne Befund: §1 Ziel und Abgrenzung — drei Ausschlüsse mit
  Begründung, konsistent mit den vier Klassen (Bestand/anderer
  Vorgang/Schicht-Abgrenzung), keine erfundene Zusatzzahl.
- geprüft, ohne Befund: SQL-Form der neuen Abfrage — `WHERE source_id =
  $1` gegen `cdc.retention_blockers` ist konsistent mit deren
  `DISTINCT ON (source_id)`-Vertrag (höchstens eine Zeile je Quelle,
  `tools/schema/schema.yaml:332`); `backlog *int64`-Zeiger-Lesart ist
  defensiv gegen einen laut Plan-Nachzug strukturell nicht erreichbaren
  Fall, analog zum bestehenden `cdc_consumer_lag`-Zweig in derselben
  Funktion — kein Shadowing-/Wiederverwendungsfehler bei der `err`-Variable.
- geprüft, ohne Befund: `make gates` — real ausgeführt, alle vier inneren
  Gates (`baseline-verify`, `docs-check`, `commit-traceability`,
  `a-check`) grün, 0 Befunde.

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Beispiel-Ausgabe widerspricht der
eigenen Schema-Formel“ · „Beispiel-Ausgabe widerspricht dokumentiertem
Invariant“

## Verdikt

**Merge-blockierend:** nein. Beide Findings betreffen ausschließlich die
illustrative „Ausgabe (Beispiel)“ in `docs/user/benutzerhandbuch.md`, nicht
Code, nicht die reale `docker exec`-Testassertion (die prüft nur Textmuster
und numerische Form, keine konkreten Zahlenwerte, und ist davon nicht
betroffen). Kein Rollen-Widerspruch, keine Konflikt-Sequenz nach Modul 8
erforderlich — beide Findings sind isolierte Dokumentationskorrekturen ohne
Implementer-Widerspruch bislang.

Alle repo-spezifisch besonders geprüften Punkte aus dem Auftrag
(Kommentar-Disziplin, docs-check-Link-Falle, „Kein Blocker“-Randfall, reales
rotes Gegenbeispiel, Image-Digest-Aktualität, bestehendes Ausgabeformat,
Commit-Traceability, DoD-Checkboxen, Isolation gegenüber `slice-048`) sind
ohne Befund.

**Übergabe:** F-1 und F-2 gehen an den Implementer zur Korrektur der
Beispiel-Ausgabe (konsistente Zahlen wählen, bei denen `cdc_consumer_lag`
und `Rückstand` für denselben Consumer übereinstimmen; Name gleich Kennung
setzen oder einen Consumer mit tatsächlich abweichendem Namen — z. B. über
einen nicht-CLI-Registrierungsweg — als Beispiel wählen und das im Text
kennzeichnen) — keine Fixrunde am Code nötig. Dieser Report ist ein
Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen; die
Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5).
