# Verifier-Report: slice-020 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code), §6 (Risiko-Ausgang, bleibt
Planner-Entscheidung), §8 (Sub-Area-Prüfung), sowie Entscheidungs-Konformität
gegen `ADR-0023` (Fehlerklassifikation, Accepted) und `SPEC-008` (Referenz-
Tabelle der sieben Fehlerklassen). Nicht geprüft: Diff gegen Plan/Hard Rules
im Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe, bereits erledigt,
siehe [`review-slice-020.md`](review-slice-020.md), 0 HIGH / 1 MEDIUM (F-1,
pre-existing, kein Merge-Blocker) / 1 LOW (F-2) / 1 INFO (F-3)), realer
Bedarf (Validator — hier nicht einschlägig, reiner Doku-Slice, kein
MVP-Grenz-Slice).

**Gegenstand:** drei Commits im Slice-Fenster: `bbdd5f4` (Implementer:
Handbuch-Tabelle auf sieben Klassen erweitert, DoD-Checkboxen teilweise
nachgezogen), `051096e` (Planner-Korrektur: „· seit slice-020" aus der
Änderungshistorie entfernt), `c679509` (Review-Report). Reiner
Dokumentations-Slice — kein Produktionscode geändert; entsprechend kein
Mutationstest-Apparat in diesem Lauf (nichts Ausführbares wurde geändert).

**Grundsatz:** Keine Behauptung wurde übernommen — jede Tabellenzeile, jeder
Code-Pfad und jeder Sensor-Lauf unten wurde in diesem Lauf selbst geprüft.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-020-handbuch-fehlerklassen-vollstaendig.md`)
- `spec/pflichtenheft.md` `SPEC-008` (Referenz-Tabelle, wörtlich gelesen)
- `docs/plan/adr/0023-fehlerklassifikation.md` (Accepted)
- `docs/user/benutzerhandbuch.md` §6 (Fehlerklassen), §9 (Änderungshistorie)
- Code im Volltext: `internal/domain/model/errorstate.go`,
  `internal/bootstrap/wiring.go` (`Run`, `classifyRunError`, `reportFault`,
  `runHeartbeat`), `cmd/pg-change-feed/main.go`
- `docs/reviews/review-slice-020.md` (0 HIGH, 1 MEDIUM F-1, 1 LOW F-2, 1 INFO F-3)
- `git show`/`git log` über `bbdd5f4`, `051096e`, `c679509`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 222 Datei(en), 0 Befund(e)` (Standardlauf und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s)` · `a-check: 0 Befund(e)` | **0** |
| `grep -rn "ErrorClassTransient\|ErrorClassPermission" --include="*.go" .` | Treffer nur in `errorstate.go` (Konstanten-Deklaration + Enum-Test), `errorstate_test.go`, `heartbeat_test.go` (Test ruft `Fault` direkt mit vorgegebener Klasse auf — kein Adapter klassifiziert selbst) | — |
| `grep -n "func classifyRunError" -A 22 internal/bootstrap/wiring.go` | genau ein `default:`-Zweig → `model.ErrorClassInternal` | — |
| `grep -n "slice-[0-9]" docs/user/benutzerhandbuch.md` | keine Treffer (leer) | — |
| `git show --stat bbdd5f4/051096e/c679509` | genau die in Plan-§3 genannten Dateien, keine unangekündigte Datei | — |
| `git status` | sauber, keine Restspuren dieses Laufs | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Handbuch-Tabelle führt alle sieben Klassen mit Bedingung/Aktion; `internal` als real erreichbarer Fallback benannt; `transient`/`permission` als deklariert, aber nicht konstruiert | **bestätigt, mit einer Einschränkung — siehe F-1 (bereits im Review erfasst, hier unabhängig nachvollzogen)** | Alle sieben Zeilen (`transient`, `configuration`, `permission`, `schema`, `storage`, `replication`, `internal`) stehen in `benutzerhandbuch.md` §6, in `SPEC-008`-Reihenfolge, Bedingung und Aktion je Zeile inhaltlich deckungsgleich mit `SPEC-008` — **außer** der `replication`-Zeile: Handbuch sagt „Sichtbarer Fehler" (und ordnet sie über den Sammelsatz dem Fatal-Pfad zu), `SPEC-008` sagt „Überwachung über Schwellen (§5); kontrollierte Fortsetzung". Eigene Code-Prüfung (`Run` in `wiring.go`, `main.go`) bestätigt: Das Handbuch beschreibt den **heutigen** Code korrekt (siehe Abschnitt unten); `SPEC-008` beschreibt an dieser einen Zeile einen Ziel-Zustand, der noch nicht implementiert ist. `internal`-Fallback: `classifyRunError` hat exakt einen `default:`-Zweig → `ErrorClassInternal`, Handbuch-Aussage korrekt. `transient`/`permission`: `grep` bestätigt, kein Adapter/keine Verdrahtung konstruiert sie real — nur Konstanten-Deklaration und Tests, Handbuch-Aussage korrekt |
| 2 | Änderungshistorie (§9) trägt einen Eintrag für diese Korrektur | **bestätigt** | Version 1.2, 2026-09-12, Zeile nennt `ADR-0023`/`SPEC-008` und die drei ergänzten Klassen; keine `slice-`-Referenz mehr (`051096e` hat „· seit slice-020" entfernt — durch `git show 051096e` und eigenen `grep` über das gesamte Dokument bestätigt: 0 Treffer) |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, 0 Befunde in allen vier Gates |
| 4 | Review durchgeführt, kein offenes HIGH | **inhaltlich bestätigt — Checkbox-Formfehler, siehe V-1** | `review-slice-020.md` liegt vor (`c679509`), 0 HIGH, 1 MEDIUM (F-1), 1 LOW (F-2), 1 INFO (F-3), kein Merge-Blocker. Die DoD-Checkbox selbst steht im aktuellen Plan-Stand aber noch auf `[ ]`, obwohl der Review-Commit zeitlich vor der letzten Prüfung dieses Laufs liegt — siehe V-1 |
| 5 | Doku-Update für `error_class`-Sichtbarkeit (Duplikat von Punkt 1) | **bestätigt, korrekt annotiert** | `[x]` mit Vermerk „entfällt hier als Duplikat" — konsistent, Punkt 1 trägt den Inhalt |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin vollständig den Platzhalter-Vorlagentext — Planner-Arbeit, wie erwartet |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `find . -iname reconciliation.md` unter `docs/plan/planning/` liefert nichts; Repo durchgehend GF (`harness/conventions.md`) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/` unverändert seit vor diesem Slice (nur `BEO-PGC/` und `README.md`, keine neue Datei durch diesen Slice) — Planner-Arbeit bei Closure |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | das einzige Risiko (Drift von `classifyRunError` zwischen Planung und Umsetzung) trägt noch keinen Ausgang — Text markiert „wird bei Closure eingetragen"; Planner-Arbeit |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt vermerkt** | Slice ist wellenlos (Kopf-Feld „ohne Welle") — Paarungen laufen bei dieser Closure selbst, noch nicht fällig, da Closure aussteht |

**Zwischenstand: 6/10 Kriterien materiell erfüllt und in diesem Lauf selbst
nachgeprüft (real, inkl. eigener Code-Lektüre von `Run`/`classifyRunError`
und eigenem `grep`), 1 Item korrekt entfallen (7), 3 Items regulär offen als
Planner-Closure-Arbeit (6, 8, 9). Punkt 1 trägt einen bereits vom Reviewer
erfassten, hier unabhängig bestätigten Befund (F-1); Punkt 4 trägt einen
neuen, eigenständigen Formfehler (V-1, unten) — beide non-blocking.**

## F-1 — eigenständig nachvollzogen: Handbuch oder `SPEC-008` — was stimmt für den heutigen Code?

**Frage aus dem Auftrag:** Beendet ein `replication`-klassifizierter Fehler
den Prozess sofort (wie das Handbuch sagt), oder gibt es bereits eine Form
von „Überwachung über Schwellen; kontrollierte Fortsetzung" (wie `SPEC-008`
verlangt)?

**Eigene Prüfung von `internal/bootstrap/wiring.go`s `Run`-Funktion
(vollständig gelesen, Zeile 230–349):**

- `stream.Run(ctx)` liefert `streamErr` direkt an die Rückgabe der Funktion
  zurück: `streamErr := stream.Run(ctx); stopHeartbeat(); heartbeatDone.Wait(); return streamErr`.
  Es gibt keinen Zwischen-Schritt, der einen `replication`-klassifizierten
  Fehler gegen eine Schwelle zählt oder den Lauf fortsetzt — jeder von
  `stream.Run` zurückgegebene Fehler propagiert unverändert nach oben.
- `reportFault` (per `defer`, LIFO vor `heartbeat.Close()`) schreibt den
  klassifizierten Fehlerzustand über `Fault(...)`, **beendet den Lauf aber
  nicht selbst** — er läuft, *weil* `Run` bereits mit einem Fehler
  zurückkehrt, nicht als eigener Kontrollfluss.
- Der eigene Code-Kommentar direkt über `reportFault` bestätigt die Absicht
  explizit: „die Sichtbarkeit gilt für „Erfassung kann nicht fortsetzen"
  (`Run` endet auf jeden Adapter-Fehler, Dateikommentar oben)".
- `cmd/pg-change-feed/main.go` Zeile 100–102: `if err := bootstrap.Run(ctx, cfg); err != nil { fmt.Fprintf(...); os.Exit(1) }`
  — jeder nicht-`nil`-Fehler aus `Run` (inklusive `replication`) beendet den
  Prozess mit Ausgang 1. Kein Schwellen-Zähler, kein Wiederanlauf-Pfad
  innerhalb desselben Prozesses.
- Es existiert **keine** Stelle im Repo, die WAL-Rückstand oder eine
  Fehler-Häufigkeit gegen eine Schwelle zählt und daraufhin *fortsetzt*
  statt abzubrechen — die einzige Schwellen-Überwachung im Repo betrifft
  Metrik-Warnschwellen (`SPEC-009`/§5 des Pflichtenhefts), nicht einen
  Kontrollfluss-Zweig in `Run`.

**Eigenes Verdikt (unabhängig vom Reviewer):** Das **Handbuch beschreibt den
heutigen Code korrekt** — ein `replication`-Fehler beendet den Prozess sofort
mit Ausgang 1, wie alle sechs anderen Klassen auch. `SPEC-008`s Aktionsspalte
für `replication` („Überwachung über Schwellen; kontrollierte Fortsetzung")
ist an dieser Stelle ein **Ziel-Zustand, der im Code noch nicht existiert** —
die Divergenz liegt in `SPEC-008` selbst (Soll-Beschreibung vor der
Umsetzung), nicht im Handbuch (Ist-Beschreibung) und nicht im Code. Das
deckt sich mit dem Reviewer-Befund F-1, ist hier aber durch eigene
Code-Lektüre der vollständigen `Run`-Funktion und `main.go`s Fehlerpfad
bestätigt, nicht übernommen.

**Einordnung:** F-1 ist damit **kein DoD-Verstoß dieses Slices** — §1 dieses
Slice grenzt „Code-Änderung an `classifyRunError`/`ErrorClass`" bewusst aus
(Bestand bleibt stehen), und die Divergenz betrifft eine Zeile, die von
diesem Diff nicht neu geschrieben, sondern nur umsortiert wurde (durch
`git show bbdd5f4` bestätigt: `replication`-Zeileninhalt unverändert
gegenüber dem Vorgänger-Stand `673d667`). Non-blocking, wie im Review
bewertet; gehört als Risiko-Ausgang oder Beobachtungs-Register-Eintrag in
die Slice-Closure, nicht in eine Fixrunde dieses Slices.

## Eigener Befund

### V-1 — DoD-Checkbox „Review durchgeführt" steht noch auf `[ ]`, obwohl der Review-Report bereits committet ist

- `kategorie`: LOW (Formfehler, kein Substanz-Defekt — die zugrundeliegende
  Tatsachenbehauptung wird durch diesen Lauf unabhängig bestätigt: Review
  liegt vor, kein offenes HIGH)
- `pfad`: `docs/plan/planning/in-progress/slice-020-handbuch-fehlerklassen-vollstaendig.md:82`
- `befund`: Der Review-Report (`c679509`, Zeitstempel 13:34:26) wurde nach
  dem Implementer-Commit (`bbdd5f4`, 13:24:30) und der Planner-Korrektur
  (`051096e`, 13:26:17) erstellt und liegt vollständig unter
  `docs/reviews/review-slice-020.md` vor — 0 HIGH, kein Merge-Blocker. Die
  DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt vor"
  im Slice-Plan trägt aber weiterhin `[ ]` statt `[x]`, weder in `c679509`
  selbst (der ausschließlich die Report-Datei hinzufügt, `git show --stat`
  bestätigt) noch in einem Folge-Commit nachgezogen. Verglichen mit dem
  Muster aus `slice-023` (dort wurde derselbe Checkbox-Nachzug explizit als
  eigener kleiner Fixrunden-Schritt vollzogen) ist dieser Nachzug hier noch
  offen.
- `verifizierbar`: ja — `git show --stat c679509` zeigt nur die
  Report-Datei; `grep -n "Review durchgeführt" -A0` im aktuellen
  Plan-Stand zeigt `[ ]`.
- **Für die Closure:** Kein Blocker — die Checkbox untertreibt den
  tatsächlichen Stand (konservative Richtung, keine Überclaiming), sollte
  aber vor `git mv` nach `done/` auf `[x]` nachgezogen werden, damit die
  Closure-Notiz nicht auf einer Lücke zwischen Plan-Text und Repo-Zustand
  aufsetzt.

## Plan-vs-Code-Diff (gegen Plan-§3)

`git show --stat bbdd5f4` zeigt genau zwei Dateien: die Slice-Plan-Datei
(DoD-Nachzug) und `docs/user/benutzerhandbuch.md` — exakt die in §3
genannte Datei/Komponente (`benutzerhandbuch.md` §6 und §9). `051096e`
ändert ausschließlich `docs/user/benutzerhandbuch.md` (1 Zeile, Entfernung
der Slice-Chronik-Referenz) — eine Planner-Korrektur, die den Plan nicht
erweitert, sondern eine Hard-Rule-3.7-Konformität nachträgt (keine
Prozess-Chronik in Betreiberdoku). `c679509` fügt ausschließlich den
Review-Report hinzu. **Keine unangekündigte Datei, kein Code-Diff, keine
Abweichung von §3.**

§1-Abgrenzung gewahrt: kein Code-Diff an `classifyRunError`/`ErrorClass`
(eigener `grep`/Lektüre bestätigt), kein neuer Sensor, kein
Adapter-Konstruktions-Fund für `transient`/`permission`.

## Entscheidungs-Konformität gegen `ADR-0023`/`SPEC-008`

- Die sieben Klassen der Handbuch-Tabelle sind exakt die sieben
  `ErrorClass`-Konstanten aus `internal/domain/model/errorstate.go` und
  exakt die sieben `SPEC-008`-Zeilen — geschlossene Menge, keine Abweichung.
- `ADR-0023` selbst unverändert seit `Accepted` (kein Commit in diesem
  Slice-Fenster berührt die ADR-Datei).
- Die einzige inhaltliche Diskrepanz (F-1, `replication`-Aktion) liegt
  zwischen `SPEC-008` und dem **Code**, nicht zwischen Handbuch und Code —
  das Handbuch bildet den Code korrekt ab und ist damit für diesen Slice
  (der laut §1 die Handbuch-Tabelle gegen `ADR-0023`/`SPEC-008`s
  *Klassen-Menge* vervollständigt, nicht jede Aktionsspalte neu
  auditiert) konform.

**Ergebnis: vollständig konform**, mit der bereits im Review als F-1
(MEDIUM, pre-existing) erfassten und hier unabhängig bestätigten
Spec-Code-Divergenz bei der `replication`-Aktionsspalte.

## Negativbefunde

- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde in allen
  vier Gates.
- geprüft, ohne Befund: **Vollständigkeit der sieben Klassen** — Wortlaut
  jeder der sieben Handbuch-Zeilen gegen `SPEC-008` verglichen (nicht nur
  Zeilenzahl gezählt); sechs von sieben Zeilen inhaltlich deckungsgleich,
  eine Abweichung bereits als F-1 erfasst.
- geprüft, ohne Befund: **`internal`-Fallback-Behauptung** — eigener
  `classifyRunError`-Lesevorgang, genau ein `default:`-Zweig.
- geprüft, ohne Befund: **`transient`/`permission` unkonstruiert** — eigener
  `grep`, keine Adapter-Konstruktion, nur Konstanten und Tests.
- geprüft, ohne Befund: **Chronik-Freiheit (Hard Rule 3.7)** — `grep -n
  "slice-[0-9]" docs/user/benutzerhandbuch.md` liefert keinen Treffer.
- geprüft, ohne Befund: **Plan-vs-Code-Diff** — alle drei Commits decken
  sich exakt mit §3, keine unangekündigte Datei.
- geprüft, ohne Befund: **AGENTS.md §3.3** (`git mv` + Inhaltsänderung) —
  keiner der drei Commits verschiebt eine Datei; nicht einschlägig.
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs nennen
  `LH-FA-ADM-003` (und `bbdd5f4` zusätzlich `ADR-0023`), kein `SPEC-*`/
  `ARC-*` im Betreff.
- geprüft, ohne Befund: **§6/§7/Register unangetastet** — alle drei bleiben
  korrekt im Platzhalter- bzw. Vorlagen-Zustand, konsistent mit
  Lifecycle-Zustand `in-progress`.
- geprüft, ohne Befund: **`git status`** — sauber nach diesem Lauf.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 (F-1 bereits vom Reviewer erfasst, hier unabhängig bestätigt statt neu gezählt) |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 6/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (real, inkl. eigener Code-Lektüre von `Run`/
`classifyRunError`/`main.go` und eigenem `grep`), 1 Item korrekt entfallen
(Reconciliation-Register), 3 Items regulär offen als Planner-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, Risiko-Ausgang). **Kein DoD-Defekt
im Sinn eines unbelegten „bestätigt"-Punkts** — V-1 ist ein Formfehler
(Checkbox untertreibt den Stand), kein Substanz-Defekt.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit einem offenen
Klein-Befund (1 LOW, V-1; non-blocking) und einer unabhängig
nachvollzogenen, bereits im Review korrekt bewerteten pre-existing
MEDIUM-Divergenz (F-1).** Die Handbuch-Tabelle ist vollständig (alle sieben
Klassen, `SPEC-008`-Reihenfolge) und für sechs von sieben Zeilen inhaltlich
exakt; die eine Abweichung (`replication`-Aktion) liegt zwischen `SPEC-008`
und dem Code, nicht zwischen Handbuch und Code — **das Handbuch ist für den
heutigen Code die zutreffende Aussage**, `SPEC-008` beschreibt hier einen
noch nicht umgesetzten Ziel-Zustand (eigene, vollständige Lektüre von
`Run` und `main.go`, nicht nur Reviewer-Befund übernommen). Kein
Code-Diff, kein neuer Sensor, keine unangekündigte Datei — reiner
Doku-Slice, Plan-vs-Code-Diff deckungsgleich.

**Vor `git mv` nach `done/` zu klären (Planner):**

1. DoD-Checkbox „Review durchgeführt" auf `[x]` nachziehen (V-1) —
   Beleg-Zeile auf `review-slice-020.md`.
2. Closure-Notiz §7 schreiben (Was hat funktioniert / was ging anders /
   Steering-Loop-Eintrag, falls einer entsteht).
3. §6-Risiko (Drift von `classifyRunError`) disponieren — Ausgang
   voraussichtlich *entfallen* (Code hat sich zwischen Planung und
   Umsetzung nicht geändert, durch diesen Lauf bestätigt) oder *eingetreten*
   je nach Planner-Urteil.
4. F-1 (`replication`-Aktionsspalte, `SPEC-008` vs. Code) als eigenständigen
   Beobachtungs-Register-Eintrag oder Folge-Slice disponieren — betrifft die
   Sub-Area `*`/`PGC`; Register aktuell ohne einschlägigen Treffer (§8
   dieses Slice-Plans bereits vermerkt).
5. F-2/F-3 aus dem Review-Report (DoD-Annotation-Praxis,
   Doku-Code-Drift-Risiko ohne Sensor) in Closure-Notiz bzw. Register
   aufnehmen, wie im Review-Report selbst vorgesehen.
6. Drei Paarungen (Anker · Folge-Slice · Register) laufen bei dieser
   Closure selbst (Slice ist wellenlos), nicht bei einer Welle-Closure.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert.

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `docs-check`/d-check: 222
Dateien, 0 Befunde, Standardlauf und `--range HEAD~5..HEAD`;
`commit-traceability`: OK über 5 Commits; `a-check`: 0 Befunde). `git
status` am Ende dieses Laufs sauber.
