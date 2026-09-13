# Verifikationsbericht: slice-042 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen
Plan (`slice-042` §1/§2 DoD/§6/§8) und die reale Implementierung aus
`slice-037`/`slice-036`, nicht gegen Diff (Reviewer-Aufgabe, bereits
abgeschlossen) und nicht gegen realen Bedarf (Validator, hier nicht
ausgelöst — reine Doku-Nacharbeit, kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen, aktuellen
Slice-Plan, `slice-037` (`done/`) vollständig, `tools/schema/nacharbeit-administration.sql`,
beide Review-Reports und den tatsächlichen Diff/Code selbst — keine
Behauptung aus einem Bericht wird ungeprüft übernommen; jeder unten
genannte Sensor-/Werkzeuglauf wurde in dieser Sitzung **selbst**
ausgeführt, nicht aus den Reports zitiert.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-042-handbuch-sql-administration.md`
zum Stand `HEAD = 5e85c4c`. Commits: `3d5029b` (`next→in-progress`, vor
der Arbeit auf `main`), `294d165` (Implementierung: zwei neue
§4-Abschnitte, Changelog 1.7), `7daf5e9` (Review, 0 HIGH/1 MEDIUM/2 LOW),
`a98e099` (themenfremder Docs-Check-Fix am Review-Report — nicht
Gegenstand dieser Prüfung, geprüft nur auf Nicht-Berührung des
Prüfgegenstands), `9cae8f9` (Fixrunde F-1/F-2/F-3), `5e85c4c`
(Review-Bestätigung, fügt ausschließlich
`docs/reviews/review-slice-042-fixrunde.md` hinzu). Lifecycle-Reihenfolge
verifiziert (`git log --oneline 3d5029b^..5e85c4c`): Move → Implementierung
→ Review → Fix → Fixrunde-Bestätigung, keine Rolle springt rückwärts ohne
Übergabe-Artefakt. `git show a98e099 --stat` eigenständig geprüft: berührt
nur einen Link in einem Review-Report, nicht `benutzerhandbuch.md` und
nicht die Plan-Datei — kein Einfluss auf diese DoD-Prüfung.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Neuer Abschnitt „Tabelle live aktivieren": Voraussetzung, Vorgehen (`SELECT cdc.enable_table(...)`), Ergebnis (Antrag in `cdc.administration_request`, Goroutine verarbeitet ohne Neustart, `status` `applied`/`failed`) | **erfüllt** | `docs/user/benutzerhandbuch.md:180-213` (aktueller Stand) selbst gelesen. Signatur exakt gegen `tools/schema/nacharbeit-administration.sql:32` (`cdc.enable_table(p_source_id text, p_schema_name text, p_table_name text)`) geprüft — Parameterzahl und -reihenfolge identisch mit `SELECT cdc.enable_table('<source_id>', '<schema>', '<tabelle>')`. Asynchroner Fluss sachlich korrekt gegen `slice-037` (`done/`) geprüft: `cdc.enable_table` schreibt nur einen `pending`-Datensatz und `pg_notify` (siehe SQL-Datei Kommentar Z. 14-16), die reale Wirkung (Publication-Mitgliedschaft, Erfassungs-Bindung) trägt ausschließlich die Administrations-Goroutine (`slice-037` §1) — der Doku-Text „die Tabelle ist damit noch nicht aktiv … Erst nach `applied` erfasst der laufende Prozess" spiegelt das exakt, kein „sofort aktiv"-Framing |
| 2 | Neuer Abschnitt „Tabelle deaktivieren" (analoges Format, `SELECT cdc.disable_table(...)`, Ergebnis: Erfassung endet, Prozess läuft weiter) | **erfüllt** | `docs/user/benutzerhandbuch.md:215-250`. Signatur exakt gegen `cdc.disable_table(p_source_id text, p_schema_name text, p_table_name text)` geprüft. Nach der Fixrunde (`9cae8f9`, F-2) ist der Poll-Codeblock strukturell deckungsgleich mit „Tabelle live aktivieren" (Vorgehen → Ergebnis-Absatz → Poll-Codeblock → Status-Erläuterung) — selbst Zeile für Zeile verglichen, nicht nur aus dem Fixrunden-Report übernommen. „Feed-Container läuft für alle übrigen aktivierten Tabellen unverändert weiter" deckt sich mit `disable/service.go`: `Unpublish` wirkt gezielt auf `command.Schema`/`command.Table`, keine globale Wirkung |
| 3 | Änderungshistorie: Versionsfeld → 1.7, neuer Eintrag mit Bezug auf `LH-FA-ADM-001`/`LH-FA-CFG-002`, `ADR-0050`, `slice-036`/`037` | **erfüllt** | `docs/user/benutzerhandbuch.md:3` (`Version: 1.7`), Zeile 631 (Changelog-Tabelle): „SQL-Administration nachdokumentiert (`LH-FA-ADM-001`, `LH-FA-CFG-002`, `ADR-0050`, slice-036, slice-037, slice-042)" — nach der Fixrunde (F-3) im reinen, kommagetrennten Format der Zeilen 1.1–1.6, kein Prosa-Einschub mehr (`git show 9cae8f9` selbst verglichen) |
| 4 | `make gates` grün | **erfüllt, eigenständig reproduziert** | Selbst ausgeführt gegen `HEAD = 5e85c4c`: `baseline-verify` (v6.5.0, 54 Dateien) OK, `d-check` (341 Dateien, 0 Befunde, volle Modul-Menge inkl. `links`/`anchors`), `commit-traceability` (`HEAD~5..HEAD`, 5 Commits, „Betreffs ohne Struktur-ID" OK), `a-check` (0 Befunde) — alle vier Gates grün, kein Carveout nötig. Der Slice-Plan schreibt in §5 explizit nur `make gates` als Closure-Trigger vor (kein `make test`/`make test-integration` — anders als `slice-037`, dessen §5 diese zusätzlich verlangt); reine Doku-Änderung ohne Code-Berührung bestätigt (`git diff 3d5029b..5e85c4c --stat`: nur `docs/`-Pfade) |
| 5 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **inhaltlich erfüllt, Formular-Diskrepanz** | Beide Reports vollständig gelesen: `docs/reviews/review-slice-042.md` (0 HIGH, 1 MEDIUM F-1, 2 LOW F-2/F-3) und `docs/reviews/review-slice-042-fixrunde.md` (alle drei real verifiziert behoben, `make gates` eigenständig gegengeprüft). Die Bedingung ist damit tatsächlich erfüllt. **Aber:** Die Checkbox in §2 des Slice-Plans steht weiterhin auf `- [ ]` (Zeile 91) — weder `7daf5e9` noch `9cae8f9` noch `5e85c4c` hat sie auf `[x]` gesetzt (`git show <commit> --stat` für alle drei geprüft: keiner berührt den Checkbox-Bereich der Plan-Datei außer der ursprüngliche Implementierungs-Commit `294d165`, der die Punkte 1–4/6 setzte, Punkt 5 aber noch nicht setzen konnte, da das Review zu diesem Zeitpunkt noch nicht gelaufen war). Siehe Finding V-1 unten — dieselbe Klasse wie `verify-slice-039.md` V-1 |
| 6 | Doku-Update — entfällt zusätzlich (dieser Slice **ist** das Doku-Update) | **erfüllt** | Korrekt als `[x]` markiert; kein weiterer öffentlicher Vertrag außerhalb `docs/user/benutzerhandbuch.md` berührt (`git diff 3d5029b..5e85c4c --stat` bestätigt) |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); §7 des Plans ist noch die Bedienhinweis-Vorlage (`<…>`-Platzhalter). Kein Verifikations-Gegenstand dieser Prüfung |
| 8 | Reconciliation-Register fortgeschrieben | **entfällt strukturell** | `find docs/plan/planning -maxdepth 1 -name reconciliation.md` ohne Treffer — Repo ist Greenfield (`harness/conventions.md` §Modus-Deklaration, ausschließlich Sub-Area `*`/`PGC`), kein Brownfield-Bootstrap. Das Item nennt diese Bedingung selbst; strukturell nie zutreffend für dieses Repo, kein offener Planner-Punkt |
| 9 | Beobachtungs-Register fortgeschrieben | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit. Zur Einordnung selbst geprüft: `BEO-PGC/verwaltung-keine-sql-administration/state.md` steht bereits auf `verkörpert` (`seit welle-12`, ausgelöst durch `slice-036`/`037`/`038`) — dieser Slice trägt keine neue Fähigkeit nach, sondern nur ihre Dokumentation; §8 des Plans sichtet das Register korrekt und findet keinen spezifisch dokumentationsbezogenen Treffer |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit; das einzige Risiko in §6 trägt noch `<bei Closure einzutragen>` |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit, wellenlos hier statt bei einer Welle-Closure fällig, aber erst **nach** dem `git mv` nach `done/` sinnvoll prüfbar |

## 2. Finding V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen

- **Klasse:** Verifier-only — für Tests und Review unsichtbar, weil beide
  Rollen ihre eigene Arbeit erledigt haben; nur ein Blick auf den
  *Formular-Zustand nach* der Review-Sequenz deckt die Lücke auf.
- **Befund:** DoD-Punkt 5 in
  `slice-042-handbuch-sql-administration.md:91` steht auf `- [ ]`, obwohl
  die Bedingung — Review durchgeführt, Report liegt vor — seit Commit
  `7daf5e9` faktisch erfüllt ist und seit `5e85c4c` zusätzlich als
  „behoben" bestätigt vorliegt. Keiner der Review-/Fix-Commits hat die
  Checkbox im Plan-Dokument nachgezogen.
- **Einordnung:** kein inhaltlicher Mangel an den beiden neuen
  Handbuch-Abschnitten oder der Review-Substanz — beide sind, wie oben
  belegt, real und reproduzierbar erfüllt. Es ist eine Diskrepanz
  zwischen dem DoD-Formular und der tatsächlichen Sachlage, exakt die
  Klasse, für die der Verifier existiert. Dieselbe Klasse trat bereits in
  `verify-slice-039.md` V-1 auf — zweites Auftreten dieser Finding-Klasse
  (Steering-Loop-Beobachtung, aber kein eigenständiges Register-Item
  dieser Prüfung; das gehört in die Closure-Notiz des Planners).
- **Erwartete Korrektur:** Checkbox auf `[x]` setzen, in einem eigenen,
  kleinen Commit vor dem `git mv` nach `done/` (Inhalt vor Move, Modul 5
  §git mv + Inhaltsänderung). Kein Rollback, keine Rückführung
  (`in-progress→next`/`open`) — die Korrektur ist ein Ein-Zeilen-Nachzug.

## 3. Sachliche Korrektheit gegen `slice-037`/`slice-036` (eigenständig, zentrale Prüffrage)

- **Asynchroner Antrags-Fluss:** Beide Abschnitte beschreiben den Aufruf
  korrekt als asynchron — `pending` → `applied`/`failed`, kein
  „sofort aktiv/inaktiv"-Framing. Selbst gegen `tools/schema/nacharbeit-administration.sql`
  (Z. 32-68, `enable_table`/`disable_table` schreiben ausschließlich
  `INSERT … VALUES (…, 'pending')` + `pg_notify`) und gegen `slice-037`
  §1 (die Goroutine trägt die reale Wirkung nach) gegengelesen.
- **Funktionssignaturen exakt:** `cdc.enable_table(p_source_id text,
  p_schema_name text, p_table_name text)` /
  `cdc.disable_table(p_source_id text, p_schema_name text, p_table_name
  text)` — Reihenfolge und Anzahl der Parameter stimmen exakt mit den
  Doku-Aufrufen überein (selbst Zeichen für Zeichen verglichen).
- **Voraussetzungen nach der Fixrunde vollständig:** Beide Abschnitte
  nennen jetzt „die physische Tabelle existiert; `REPLICA IDENTITY` ist
  wie benötigt gesetzt" mit Verweis auf `[Tabelle
  aktivieren](#tabelle-aktivieren)`, zusätzlich `cdc_admin`-Mitgliedschaft
  über `CDC_ADMIN_DSN`. Gegen `internal/application/usecase/enable/service.go`
  (Kommentar Z. 26-31: „die Existenz der physischen Tabelle ist die
  Vorbedingung … die `REPLICA IDENTITY` bleibt vom Aktivieren unberührt")
  und `disable/service.go` (`TableExists`-Prüfung mit explizitem
  Fehlerpfad `ErrSourceTableMissing`) selbst gelesen: Die Existenzprüfung
  ist code-real (`activation.TableExists`), `REPLICA IDENTITY` ist keine
  code-geprüfte, sondern eine operative Vorbedingung für korrektes
  Replikationsverhalten — identisch zur bereits bestehenden „Tabelle
  aktivieren"-Sektion (CLI-Pfad), auf die beide neuen Abschnitte korrekt
  verweisen statt sie zu duplizieren.
- **Keine Chronik-Sprache außerhalb der Änderungshistorie:** `grep -n
  "slice-\|welle-\|Welle" docs/user/benutzerhandbuch.md` liefert
  ausschließlich Treffer in der Changelog-Tabelle (Zeilen 615-631), keinen
  in den operativen §4-Abschnitten — `AGENTS.md` §3.7 auf Markdown-Prosa
  angewandt, selbst geprüft.
- **Anker lösen auf:** `#tabelle-aktivieren`, `#zugriff-und-rollen` — durch
  den eigenständigen vollen `make gates`-Lauf (d-check mit aktiviertem
  `anchors`/`links`-Modul, 0 Befunde) mitbestätigt, nicht nur behauptet.

## 4. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

Eigenständig geprüft, nicht aus dem Review übernommen:

- `git diff 3d5029b..5e85c4c --stat`: ausschließlich
  `docs/user/benutzerhandbuch.md`, die Plan-Datei selbst und
  `docs/reviews/review-slice-042*.md` — kein Treffer in `internal/`,
  `tools/schema/`, `spec/`. Ausschluss 1 (Code-/Verhaltensänderung) und
  Ausschluss 3 (`pflichtenheft`/`architecture`-Änderungen) beide
  eingehalten.
- `grep -n "register-consumer\|acknowledge-consumer\|Consumer" docs/user/benutzerhandbuch.md`
  in den neuen Abschnittszeilen (181-250): kein Treffer — Ausschluss 2
  (Consumer-Verwaltung über SQL) eingehalten, `welle-12` (`done/`) bleibt
  die korrekte Referenz für diesen bereits bestehenden Ausschluss.

## 5. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Wie in Auftrag benannt und durch §Träger im Repo ohne Wellen-Betrieb
(Modul 6/8) gedeckt: Closure-Notiz, Reconciliation-/Beobachtungs-Register-
Fortschreibung, Risiko-Ausgang (§6), die drei Paarungen. Alle zugehörigen
§2-Häkchen sind **korrekt unbeansprucht** — kein Mangel, sondern der
vorgesehene Zustand vor dem nächsten Rollenwechsel an den Planner. Auch
nicht Gegenstand: Validierung gegen realen Bedarf (kein
MVP-Meilenstein-Slice, kein Validator-Zug ausgelöst; reine Doku-Nacharbeit
zu bereits validierter Fähigkeit aus `welle-12`).

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-blockierenden Formular-Diskrepanz (V-1: Checkbox 5 nicht
nachgezogen — inhaltlich längst erfüllt). Alle drei substanziellen
Implementer-Liefer-Punkte (1–3) sind durch eigene, unabhängige Prüfung
gegen `slice-037`/`slice-036` (reale Implementierung) und
`tools/schema/nacharbeit-administration.sql` (exakte Signaturen) gedeckt:
beide neuen Handbuch-Abschnitte beschreiben den asynchronen
Antrags-Queue-Fluss sachlich korrekt, kein „sofort wirksam"-Framing, die
Voraussetzungen (Tabellenexistenz, `REPLICA IDENTITY`, `cdc_admin`) stehen
seit der Fixrunde vollständig in beiden Abschnitten. `make gates` selbst
grün gelaufen (§4 des Plans verlangt hier korrekt kein `make test`/`make
test-integration` — reine Doku-Änderung ohne Code-Berührung, verifiziert
über `git diff --stat`).

**Zur zentralen Prüffrage dieser Verifikation** (sachliche Korrektheit
gegen den realen asynchronen Antrags-Fluss aus `slice-037`, exakte
Funktionssignaturen, vollständige Voraussetzungen): **bestätigt**, eigenständig
gegen den Code und die reale Implementierungsbeschreibung nachvollzogen,
nicht aus den Review-Reports übernommen.

**Die vier Planner-Closure-Punkte** (Closure-Notiz, Reconciliation-/
Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) sind **korrekt
unbeansprucht** — bewusst nicht Gegenstand dieser Prüfung (Modul 8, Träger
Planner). Das Reconciliation-Register-Item entfällt zusätzlich strukturell
(Greenfield-Repo, Datei existiert nicht).

**Keine Rückführung nötig.** Weder `in-progress→next` (der Slice ist
nicht zu groß — ein einziger Liefer-Punkt bleibt innerhalb einer Schicht,
Doku, keine vierte Berührung entstanden) noch `in-progress→open` (kein
Blocker). Vor dem `git mv` nach `done/` ist lediglich die
Checkbox-Korrektur aus V-1 fällig — ein Ein-Zeilen-Commit, kein
Zerlegungs- oder Blocker-Fall.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für die
Closure-Entscheidung, mit dem Hinweis V-1 zur Nachbesserung vor dem
`git mv`. Kein Validator-Zug ausgelöst — `slice-042` ist reine
Doku-Nacharbeit zu einer bereits geschlossenen Welle
([`LH-FA-ADM-001`](../../spec/lastenheft.md)), kein MVP-Meilenstein-Slice
im Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
