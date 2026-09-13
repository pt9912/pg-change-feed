# Review-Report: slice-058 — 2026-09-13

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-058` — sechs Commits auf `main`, bereits gepusht,
Elter `6d259f2` (reiner `next→in-progress`-Move):

- `ea9da9d` — `ChangeNotificationPort.Notify` auf vier-Token-Subjekt +
  defensive Validierung
- `f35ee72` — `model.Change` trägt `Schema`/`Table` für den Notify-Pfad
- `345d61b` — `CaptureService.Capture()` dedupliziert je `(schema, table)`
- `afa55bd` — Happy-Path-Testbeleg (`run-integration-tests.sh`) auf
  vier-Token-Subjekt nachgezogen
- `29ac256` — Benutzerhandbuch-Korrektur + Änderungshistorie 1.11
- `ac874b7` — Image-Digest-Nachzug (`harness/image-hash.txt`)

**Skill:** `.harness/skills/reviewer.md` @ `fbb6042` (Stand 2026-09-13, vier
repo-spezifische HIGH-Regeln, Slice-/Wellen-Chronik-Regel ergänzt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-058-nats-tabellen-granulares-subjekt.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §6 Risiken, §8 Sub-Area)
- `ADR-0056` (`docs/plan/adr/0056-nats-tabellen-granulares-subjekt.md`) —
  bindend: Subjekt-Schema, Notify-Kardinalität, Port-/Modell-Erweiterung
- `ADR-0055` (Punkte 1, 3, 4, 5 unverändert in Kraft)
- `LH-FA-SST-007`, `SPEC-017` (`spec/pflichtenheft.md`, bereits durch
  `ADR-0056` auf das vier-Token-Schema aktualisiert), `ARC-013`
  (`spec/architecture.md`, ebenfalls bereits aktualisiert)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.8), §6 Workflow
- `harness/conventions.md` (MR-000 ID-Schema)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-053.md`
  (F-1 dort verweist explizit auf diesen Slice als Folge-Slice)
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/`
  (Anlass für F-1 unten)
- Reale Gate-Läufe dieses Reviews: `make a-check`, `make docs-check`,
  `make commit-traceability`, `make test`, `make coverage-gate` — alle
  grün (siehe Negativbefunde)

---

## Findings

### F-1 — Benutzerhandbuch-Versionshistorie in `slice-053` übersprungen: drittes Auftreten, Schwelle erreicht

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Das
  Beobachtungs-Register (Ausgang bei 3×: „wird zur verkörperten Regel");
  Beobachtung `BEO-PGC/handbuch-versionshistorie-uebersprungen`
- `pfad`: `docs/user/benutzerhandbuch.md` — Commit `6ddb7ae` (`slice-053`,
  bereits in `done/`), nicht Teil des `slice-058`-Diffs selbst
- `befund`: Unabhängig verifiziert per `git log -p --follow -- docs/user/benutzerhandbuch.md`
  und `git show 6ddb7ae^:docs/user/benutzerhandbuch.md` /
  `git show 6ddb7ae:docs/user/benutzerhandbuch.md`: `6ddb7ae` fügte die
  `CDC_NATS_URL`-Tabellenzeile ein (öffentlicher Vertrag, Doku-Update im
  Sinne der DoD von `slice-053`), ohne den `Version:`-Kopf (blieb bei
  1.10) oder die `### Änderungshistorie`-Tabelle fortzuschreiben. Das
  Register führt zu diesem Muster bereits zwei Belege
  (`evidence/slice-045.md`, `evidence/slice-046.md`, Stand `state.md`:
  „2× — unter der Schwelle"). `6ddb7ae` ist ein **drittes**, bisher nicht
  im Register erfasstes Auftreten und hebt den Zähler auf 3× — die
  Schwelle, ab der der Eintrag laut Regelwerk einen Ausgang
  (*verkörpert* oder *geplant*) statt `offen` trägt.
  `slice-058`s eigener Commit `29ac256` macht genau diesen Fehler
  **nicht** — er korrigiert die Zeile *und* bumpt Version/Changelog auf
  1.11 korrekt. Der Fund betrifft also nicht den geprüften Diff selbst,
  sondern eine noch nicht registrierte dritte Instanz aus einem bereits
  geschlossenen Vorgänger-Slice, die durch die Arbeit an diesem Slice erst
  sichtbar wurde (derselbe Datei-Diff lag im Blick).
- `verifizierbar`: ja — `git log -p --follow -- docs/user/benutzerhandbuch.md`
  zeigt `6ddb7ae` ohne Version-/Changelog-Änderung, `git diff 6ddb7ae^..6ddb7ae`
  bestätigt `Version: 1.10` unverändert.
- `klasse`: „Handbuch-Versionshistorie übersprungen (3. Auftreten, Schwelle
  erreicht)"

**Einordnung zur Kernfrage:** Kein Fixrunden-Fall für `slice-058` — der
Fund betrifft einen bereits in `done/` liegenden Vorgänger-Commit, den
dieser Diff nicht ändert und nicht ändern muss. Die Registrierung als
drittes `evidence/`-Beleg-Dokument und die Zuweisung eines Ausgangs
(*verkörpert*: Sensor/Regel-Schärfung · *geplant*: Kennung eines
Folge-Slice) ist Planner-Arbeit bei der nächsten Slice- oder
Welle-Closure (Baseline-Regelwerk `modul-08-agentenrollen.md` §Rollen-Sequenz
für eine Welle, Schritt 3a/3b) — hier nur festgehalten, damit der Fund
nicht verloren geht.

### F-2 — `model.Change.Schema`/`.Table` ohne Konstruktor-Invariante: begründeter, aber bewusst ungewächteter Tradeoff

- `kategorie`: INFO
- `quelle`: Maintainability · `ADR-0056` §Konsequenzen ("`model.Change`
  wächst um Felder, die nur der Notify-Pfad braucht … vertretbar")
- `pfad`: `internal/domain/model/change.go:39-50` (Felder ohne
  `NewChange`-Prüfung), `internal/adapters/driving/replication/mapper/mapper.go:245-246`
  (einzige produktive Setzstelle)
- `befund`: `Schema`/`Table` werden nach `NewChange(...)` per
  Direktzuweisung gesetzt, nicht als Konstruktor-Parameter — damit bleibt
  der zweite Konstruktionsort (`postgresstorage/mapper.ToChange`,
  SQL-Lesepfad) unverändert, der die Felder nie braucht. Das ist mit
  `ADR-0056` explizit vereinbar ("Name und exakte Platzierung … legt der
  umsetzende Slice fest"). Die Kehrseite: Ein hypothetischer *dritter*
  Konstruktionsort, der die beiden Zuweisungszeilen vergisst, würde vom
  Compiler nicht erkannt — die einzige Absicherung ist die defensive
  Prüfung in `natsnotify.Notify` (leere Schema-/Tabellennamen werden dort
  abgelehnt, Fehler landet als `Warn`-Log, nicht als Testfehler) und nur
  dann wirksam, wenn `CDC_NATS_URL` überhaupt aktiv ist. Aktuelle Inventur
  (Grep gegen `model.NewChange(`) bestätigt: genau zwei
  Konstruktionsstellen existieren, exakt wie `ADR-0056` §Kontext
  behauptet — das in §6 des Slice-Plans notierte Risiko 1 ist damit durch
  reale Inventur gedeckt, nicht nur Annahme.
- `verifizierbar`: ja — `grep -rn "model.NewChange("` bestätigt genau zwei
  Fundstellen; `make test` (`mapper_test.go` prüft `Schema`/`Table` an der
  einzigen produktiven Setzstelle).
- `klasse`: „Domänenfeld ohne Konstruktor-Invariante, Absicherung nur am
  Adapter-Rand"

## Negativbefunde

- geprüft, ohne Befund: `internal/application/port/outbound/changenotification.go` —
  Signaturänderung exakt nach `ADR-0056` (`sourceID, schema, table`
  getrennt, nicht vorkombiniert), Godoc beschreibt aktuellen Zustand +
  ADR-Bezug, keine Slice-Chronik
- geprüft, ohne Befund: `internal/adapters/driven/natsnotify/notify.go` —
  Subjekt-Bildung real vier Tokens (`cdc.changes.<source_id>.<schema>.<table>`),
  Validierung (`.`, `*`, `>`, Whitespace) läuft vor dem `conn.Publish`-Aufruf,
  nicht danach; Reihenfolge sourceID → schema/table leer → reservierte
  Zeichen ist konsistent mit `ADR-0056` Folgepflicht
- geprüft, ohne Befund: `internal/adapters/driven/natsnotify/notify_test.go` —
  fünf Negativ-Fälle für reservierte Zeichen/Whitespace (Punkt in Schema,
  Punkt in Tabelle, Stern, Größer-als, Leerzeichen), leere Schema-/
  Tabellenfälle separat getestet
- geprüft, ohne Befund: `internal/application/usecase/capture/service.go` —
  `distinctTables` dedupliziert real über eine Map, Aufruf-Reihenfolge in
  erster Auftrittsreihenfolge; Fehlerpfad bleibt best-effort (`Warn`-Log,
  kein Rückgabefehler), unverändert aus `ADR-0055` Punkt 4
- geprüft, ohne Befund: `internal/application/usecase/capture/service_test.go` —
  zwei neue Regressionstests (ein Notify für drei Changes derselben
  Tabelle; zwei distinkte Notifys für zwei Tabellen), beide grün
  (`make test`)
- geprüft, ohne Befund: `internal/adapters/driving/replication/mapper/mapper.go`/`mapper_test.go` —
  `event.Relation.Schema`/`.Name` fließen ohne neuen Lookup in
  `Change.Schema`/`.Table`, Regressionstest prüft konkrete Werte
- geprüft, ohne Befund: Inventur der `model.NewChange`-Aufrufstellen (Grep)
  — genau zwei, wie in `ADR-0056` §Kontext behauptet; kein dritter,
  unentdeckter Aufrufer
- geprüft, ohne Befund: `tools/harness/natssub/main.go` — unverändert und
  zu Recht unverändert: das Subjekt kommt als generischer Aufrufparameter,
  keine Subjekt-Bildung im Tool selbst
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` — Subjekt
  real auf `cdc.changes.src-mvp.public.feed_mvp_full` nachgezogen,
  Kommentar beschreibt aktuellen Zustand, keine Slice-Chronik
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (dieser Diff,
  `29ac256`) — Env-Var-Zeile korrekt auf das vier-Token-Schema korrigiert,
  `Version:`-Kopf UND Änderungshistorie-Zeile 1.11 konsistent mit dem
  etablierten Tabellenmuster (Datum, LH-/ADR-Bezug, Slice-Kennung in der
  Chronik-Spalte — dort zulässig, weil die Änderungshistorie-Tabelle
  selbst der deklarierte Chronik-Träger ist, kein Produktionscode-Kommentar)
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) in allen
  geänderten Dateien — Zusage-/Kopplungs-/Grenz-Kommentare, `ADR-*`/`LH-*`-
  Bezüge statt Slice-Chronik; keine Aussage über eine verworfene
  Alternative im Konjunktiv, kein abwesender Text
- geprüft, ohne Befund: Commit-Trailer aller sechs Commits — keine
  `Co-Authored-By:`/`Claude-Session:`-Zeilen; jede Message referenziert
  `ADR-0056`/`LH-FA-SST-007`, keine `SPEC-*`/`ARC-*`-ID im Betreff
- geprüft, ohne Befund: Commit-Reihenfolge Test-Commit (`afa55bd`) vor
  Image-Digest-Commit (`ac874b7`) — die Behauptung „real gegen einen
  frisch gebauten Feed-Container geprüft" in `afa55bd`s Message ist
  plausibel nur, wenn `make image` vor diesem Testlauf bereits lokal
  ausgeführt wurde; der Digest-Commit selbst folgt als eigener, später
  gestellter `chore`-Commit. Das ist zulässig (`AGENTS.md` §3.3 verlangt
  nur getrennte Commits, keine bestimmte Reihenfolge zwischen Test- und
  Digest-Commit), aber die Commit-Abfolge liest sich beim ersten
  Durchgang irreführend — kein Regelverstoß, daher kein eigenes Finding
- geprüft, ohne Befund: reale Sensor-Läufe dieses Reviews — `make a-check`
  (0 Befunde, unveränderter Abdeckungs-Hinweis für
  `tools/harness/natssub/main.go`), `make docs-check` (420 Dateien/0
  Befunde), `make commit-traceability` (OK, 5 Commits im Fenster),
  `make test` (alle Pakete grün, inkl. `-race`), `make coverage-gate`
  (40,90 % ≥ 35 %) — alle real in diesem Lauf ausgeführt, Exit 0
- geprüft, ohne Befund: `spec/pflichtenheft.md` (`SPEC-017`),
  `spec/architecture.md` (`ARC-013`) — bereits durch `ADR-0056` auf das
  vier-Token-Schema aktualisiert, dieser Slice liest nur, ändert nicht;
  keine Drift zum tatsächlich implementierten Subjekt gefunden

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Handbuch-Versionshistorie übersprungen
(3. Auftreten, Schwelle erreicht) · Domänenfeld ohne
Konstruktor-Invariante, Absicherung nur am Adapter-Rand

## Verdikt

**Merge-blockierend:** nein — F-1 betrifft nicht den geprüften
`slice-058`-Diff, sondern einen bereits in `done/` liegenden
Vorgänger-Commit (`slice-053`); ihre Auflösung (Register-Eintrag als
drittes `evidence/`-Beleg-Dokument, Ausgang *verkörpert* oder *geplant*)
ist ein Planner-Schritt bei der nächsten Closure, kein
Reviewer→Implementer-Rückgabepfeil. F-2 ist INFO — ein bereits von
`ADR-0056` explizit abgewogener und akzeptierter Tradeoff, durch reale
Inventur bestätigt, keine erwartete Aktion an diesem Slice.

**Übergabe:** Kein Fixrunden-Pfad zum Implementer nötig. F-1 geht als
Fund an die Planungsebene (Register-Eintrag `BEO-PGC/handbuch-versionshistorie-uebersprungen`,
drittes `evidence/`-Dokument für `6ddb7ae`/`slice-053`, Ausgangs-Zuweisung
bei der nächsten Closure). Nach `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde wird die DoD-Zeile „Review
durchgeführt" im Slice-Plan in diesem Commit mitgezogen. Die
Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in
den Beobachtungs-Register-Zähler. Dieser Report selbst ist ein
**Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und wird über Läufe hinweg nicht wieder gelesen. Der Report
ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier
separat (Modul 11).
