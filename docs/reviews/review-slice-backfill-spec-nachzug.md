# Review-Report: slice-backfill-spec-nachzug — 2026-09-24

**Review-Art:** Plan/Design — der Diff ist reine Spec-Doku (Pflichtenheft und
Architektur-Sicht) plus Slice-Plan; geprüft gegen Plan, ADRs und `AGENTS.md`
Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `77c60bd1..HEAD`, sieben Commits `8a47eb3a` …
`46d731c5`; drei Dateien: `spec/pflichtenheft.md`, `spec/architecture.md`,
Slice-Plan `slice-backfill-spec-nachzug` (Lifecycle `open` → `next` →
`in-progress`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Träger-Nachzug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-spec-nachzug` (§1 Ziel, §2 DoD, §3 Plan und
  Suchlauf, §6 Risiken) und Welle `welle-backfill-bestand`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  (Accepted) — Teilfragen 1–8, Festlegungen 1–3 des Auftraggebers,
  Folgepflicht 1
- [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  (Accepted, Supersedes `ADR-0111` Teilfrage 5 teilweise) — Festlegung 1
  (Rollenschnitt/Grants), 2 (Aufnahme), 3 (Warn-Kriterium), Folgepflicht 1
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  und [`ADR-0012`](../plan/adr/0012-at-least-once.md) (Abgrenzung
  Coverage-Fläche; Consumer-Idempotenz für die Replay-Invariante)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-FA-CON-005`](../../spec/lastenheft.md),
  [`LH-FA-REA-004`](../../spec/lastenheft.md) und die Pflichtenheft-Stellen
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md),
  [`SPEC-001`](../../spec/pflichtenheft.md),
  [`SPEC-002`](../../spec/pflichtenheft.md),
  [`SPEC-019`](../../spec/pflichtenheft.md),
  [`SPEC-020`](../../spec/pflichtenheft.md),
  [`SPEC-021`](../../spec/pflichtenheft.md),
  [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-029`](../../spec/pflichtenheft.md),
  [`ARC-006`](../../spec/architecture.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.4, §3.7, §3.9, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md` (repo-geführte
  Form, spiegelt das Output-Schema des Skills; die vendored Baseline führt
  eine Tabellenvariante desselben Gerüsts — die Repo-Form folgt den
  bestehenden `review-slice-*`-Reports)

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen):

- `git log --stat`/`--follow` über die sieben Commits: `8a47eb3a` (`open` →
  `next`) und `9e865160` (`next` → `in-progress`) sind reine Renames ohne
  Zeilenänderung; `f411b352` (Feld `Verantwortlich`) und die drei
  Inhalts-Commits liegen getrennt davon — `AGENTS.md` §3.3 erfüllt.
- `make commit-traceability RANGE=77c60bd1..HEAD` → Exit 0, „7 Commit(s) …
  Betreffs ohne Struktur-ID"; jede Message trägt `LH-FA-CAP-009` und/oder
  `ADR-0111`.
- Suchlauf über den Diff: `git diff 77c60bd1..HEAD -- spec | grep '^+' |
  grep -Ei 'ADR-|slice|welle'` → 0 Treffer; `grep -c 'ADR-0'` auf beide
  Spec-Dateien → 0/0 (`AGENTS.md` §3.4). Der einzige Treffer von
  `welle|slice` in `spec/architecture.md` ist der bestehende Kopf-Satz der
  Hard Rule (Zeile 8), nicht Teil des Diffs.
- Host-Pfade: Suchlauf nach Wurzel-Segmenten von Entwicklerrechnern und
  Laufwerks-Mustern über die Diff-Zeilen → 0 Treffer (`AGENTS.md` §3.11).
- Kennung `SPEC-029`: höchste Kennung am Parent-Stand
  (`git show 77c60bd1:spec/pflichtenheft.md`) ist `SPEC-028`; am Kopf ist
  `SPEC-029` die einzige neue; keine andere Datei unter `docs/plan` vergibt
  `SPEC-029`/`SPEC-030` (nur Verweise der Folge-Slices auf die neue Kennung).
- Zahlen des Suchlauf-Feldes (Slice-Plan §3) nachgemessen mit `git grep` an
  `77c60bd1` und `HEAD`: `exclude_column` 9/9 Treffer mit den genannten
  Zeilennummern; `Felder` 11/11; `ARC-005|ARC-006` in `spec/architecture.md`
  20 → 24; `backfill-mechanismus` 1/1 (nur die eigene Zelle);
  `CAP-009` außerhalb der Backfill-Plan-/ADR-Dateien 13 Zeilen am Parent;
  `CAP-009\.a` außerhalb derselben Dateien: `ADR-0112` Z. 393/411 und
  Pflichtenheft. Alle genannten Werte stimmen. Die Handbuch-Zeilen
  (302, 375, 377, 639, 641, 721) tragen den beschriebenen Inhalt;
  `docs/user/` ist im Diff unberührt.
- `SPEC-022`-Zeile „Position und `limit`" gegen
  `internal/adapters/driven/postgresstorage/queries/queries.go`
  (`SelectChanges`, Zeilen 44–70) gelesen: Filter `commit_position >= $2`
  (from inklusiv), `< $3` (to exklusiv), `ORDER BY t.commit_position,
  c.transaction_id, c.sequence`, `LIMIT $6`. `from = <letzte Position> + 1`
  überspringt den Rest einer angeschnittenen Position; ein Lesen ab
  derselben Position liefert dieselben ersten `limit` Zeilen. Die Aussage
  stimmt. Der Schlüsselvergleich auf der View ist tragfähig: `cdc.changes`
  führt `commit_position`, `transaction_id`, `sequence`
  (`tools/schema/schema.yaml` Z. 309–333).
- Text-Abgleich `SPEC-029`, `SPEC-019` und `LH-FA-CAP-009.a` mit
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 1 (Grants: `cdc_admin` `SELECT`/`INSERT`, `cdc_capture`
  `SELECT`/`UPDATE`, `cdc_reader` nur View, niemand `DELETE`), Festlegung 2
  (Aufnahme beim Start und Wecksignal, `(requested_at, run_id)`),
  Festlegung 3 (Punkt 1: keine Toleranz-Zahl; Punkt 2: keine Richtgröße;
  Punkt 3: zwei Warnungen ohne Statusänderung; Punkt 5: `estimated_rows`
  `NULL` = „unbekannt", nie `0`; Punkt 6: „geschätzt" an jeder Nennung der
  Zeilenzahl): stimmt überein. `grep -nE '10 Minuten|Minuten'` über die
  hinzugefügten Spec-Zeilen → 0 Treffer.
- Host-Python-Einsatz für Markdown-Ersetzungen (Bericht des Implementers):
  im Diff ist kein Skript, kein Make-Target und keine Konfiguration
  enthalten; `git status` ist sauber. Bewertung siehe Negativbefunde.

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — `SPEC-022`: Fortsetzungsregel der Zeile „Reihenfolge" ohne Vorbehalt neben der neuen Zeile „Position und `limit`"

- `kategorie`: MEDIUM
- `quelle`: `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (Wiederholung
  der Klasse — dritter Beleg); Maintainability
- `pfad`: `spec/pflichtenheft.md:614` (Zeile „Reihenfolge") gegen
  `spec/pflichtenheft.md:618` (neue Zeile „Position und `limit`")
- `befund`: Die Zeile „Reihenfolge" trägt unverändert „die Fortsetzung ist
  `from = <letzte gelieferte commit_position> + 1`", ohne Vorbehalt; vier
  Zeilen darunter sagt die neue Zeile, dass genau diese Fortsetzung bei einer
  Position mit mehr Changes als `limit` den Rest nicht liefert. Ein Leser der
  Zeile „Reihenfolge" trifft die Regel ohne die Ausnahme (zwei Aussagen im
  selben Abschnitt, keine verweist auf die andere).
- `verifizierbar`: nein — kein Gate liest Prosa-Kohärenz innerhalb einer
  Tabelle; Falsifikation ist das Lesen des Abschnitts
- `klasse`: Nachzug lässt überholten Text im selben Dokument stehen

### F-2 — §3.13-Suchlauf ohne die Plan-Träger der festgelegten Warn-Spalten und der Antragsarten-Menge

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Suchlauf über die Träger nach der bewegten
  Eigenschaft, Meldung statt stiller Mitänderung)
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-spec-nachzug.md`
  §3 Suchlauf (Zeilen 151–157) gegen Träger außerhalb des Diffs:
  `slice-backfill-run-store` §1 (Z. 40–41: „leer, solange keine Auswertung
  sie setzt"), `slice-backfill-sql-administration` (Z. 70, 77, 120–121:
  Warn-Spalte(n) „leer"), `slice-transformationen-spec-nachzug` §6
  (Z. 237–240: „nennt die vier SQL-Funktionen", Satz „vier Arten")
- `befund`: Die Suchbefehle des Feldes lesen `spec`, `docs/user`, `harness`
  und `README.md`, nicht `docs/plan/planning`; `SPEC-029` legt die
  Warn-Spalten als `boolean NOT NULL DEFAULT false` fest (nie „leer"), und
  `SPEC-019` nennt jetzt fünf Funktionen/Arten — die drei offenen Pläne
  beschreiben den Vorzustand, und das Feld enthält weder einen Fund noch eine
  Meldung dazu.
- `verifizierbar`: ja — `grep -rn 'Warn-Spalte\|vier SQL-Funktionen\|vier
  Arten' docs/plan/planning`
- `klasse`: Suchlauf deckt die Plan-Träger der bewegten Eigenschaft nicht

### F-3 — Konjunktiv über die verworfene Alternative in Spec-Prosa

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Ist-Zustand statt verworfener Alternative;
  die Abwägung gehört in die ADR) — Geltung für Doku-Prosa nach der im
  Repo geübten Auslegung, Wortlaut der Regel nennt Code, Konfiguration,
  Skripte und Zustandsfelder
- `pfad`: `spec/pflichtenheft.md:694-696` („eine gemeinsame Spalte hätte
  zwei Schreiber für einen Wert"); `spec/architecture.md:301-304` („weil ein
  Run … sonst entweder ohne angenommenen Antrag oder mit einem Antrag
  stünde, dessen Wiederholung … scheitert")
- `befund`: Beide Sätze begründen eine Festlegung durch das, was bei der
  verworfenen Alternative geschähe („hätte", „sonst … stünde"), statt den
  geltenden Zustand zu nennen.
- `verifizierbar`: nein — kein Gate liest Konjunktiv-Sprache
- `klasse`: Konjunktiv über verworfene Alternative

### F-4 — Sicht: „nichts für Leser sichtbar" gegen den sichtbaren Fortschritt in `SPEC-029`

- `kategorie`: LOW
- `quelle`: Maintainability (Zwei-Quellen-Aussage innerhalb der Spec)
- `pfad`: `spec/architecture.md:342` gegen `spec/pflichtenheft.md:685`
  (`rows_copied`: „außerhalb der Daten-Transaktion fortgeschrieben und damit
  für Leser sichtbar")
- `befund`: Die Sicht sagt „Bis zum einen Commit ist nichts für Leser
  sichtbar", während der Run-Fortschritt (`rows_copied`) laut `SPEC-029`
  bereits während des Laufs für Leser der View sichtbar ist; gemeint sind
  offenbar nur die Backfill-Changes, der Satz sagt es nicht. „bis zum einen
  Commit" ist zudem grammatisch unrund.
- `verifizierbar`: nein
- `klasse`: Unpräzise Sichtbarkeitsaussage neben einer Gegenstelle

### F-5 — Warnungen für große Tabellen ohne Lastenheft-Anker

- `kategorie`: INFO
- `quelle`: Spec-Stratum-Regel „präzisieren ja, erweitern nie" (Kopf von
  `spec/pflichtenheft.md`); Zuständigkeit Verifier/Architect
- `pfad`: `spec/pflichtenheft.md:238-245` (Punkt „Große Tabellen"),
  `spec/pflichtenheft.md:687-688` (Warn-Spalten)
- `befund`: `LH-FA-CAP-009` nennt große Tabellen nur im Out-of-Scope
  (Durchsatz als Ausbaustufe); die zwei Warnungen sind in
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Festlegung 3 des Auftraggebers und
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3 begründet, tragen im Pflichtenheft aber keinen Bezug auf
  einen Lastenheft-Text. Als Verfeinerung unter `LH-FA-CAP-009.a` gelesen
  kein Verstoß; ob sie dort als Präzisierung oder als neue Zusage zählt, ist
  Deutung, kein Reviewer-Urteil.
- `verifizierbar`: nein
- `klasse`: Technik-Zusage ohne Lastenheft-Anker (zur Deutung)

### F-6 — Architektur-Sicht ohne Zukunfts-Kennzeichnung für noch nicht gebaute Rollen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz B (Zukunfts-Aussage als Zusage
  formuliert oder mit Anker)
- `pfad`: `spec/architecture.md:268-352` (Abschnitt „Bestand als Backfill
  überführen")
- `befund`: Das Pflichtenheft kennzeichnet `LH-FA-CAP-009.a` ausdrücklich
  als Zusagen an die Umsetzung (Zeilen 150–152); die Sicht beschreibt Annahme,
  Worker, Snapshot-Leser und Commit im Indikativ ohne diese Kennzeichnung —
  das entspricht dem Stil der übrigen Sichten (Greenfield, Doku führt), die
  Umsetzung liegt aber in den Folge-Slices. Kein Widerspruch zum Plan, der für
  die Sicht keine Markierung fordert.
- `verifizierbar`: nein
- `klasse`: Zukunfts-Aussage im Indikativ ohne Marker

---

## Negativbefunde

- geprüft, ohne Befund: `spec/architecture.md` — kein ADR-/Slice-/Wellen-Bezug
  im Diff (`AGENTS.md` §3.4); `ARC-006`-Zeile und Knoten „Driven Adapters"
  tragen dieselben Rollen (Snapshot-Leser, Annahme); Antragsarten-Aufzählung
  auf fünf Werte gezogen, „vier Arten", „drei übrigen Antragsarten" und
  „beiden Tabellen-Antragsarten" nicht mehr vorhanden (`grep`); Frische-Marker
  „Letzte Änderung" gezogen.
- geprüft, ohne Befund: `spec/pflichtenheft.md` `LH-FA-CAP-009.a` — Überschrift
  ohne „offen"; Einleitung nennt die Sätze ausdrücklich „Zusagen an die
  Umsetzung, keine Messergebnisse"; Mechanismus, Markierung, Position und
  Ordnung, Überlappungs-Verhalten (keine Lücke, begrenzte idempotente
  Dopplung), Sichtbarkeits-Grenze, Neubeginn, Fail-closed, Fehlerklassen,
  Zustellung/Retention stimmen mit
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  überein; die als „hergeleitet"/„erwartet" geführten Aussagen
  (Lückenfreiheit, Replay-Invariante) stehen als Zusage, nicht als geprüfte
  Tatsache (`AGENTS.md` §3.12 Instanz B, soweit im Diff).
- geprüft, ohne Befund: `SPEC-029` — zwei Warn-Spalten
  `warn_estimated_size`/`warn_duration` (`boolean`, `NOT NULL`, Default
  `false`) mit je einem Schreiber; keine Toleranz-Zahl, keine Richtgröße, kein
  „10 Minuten"; `estimated_rows` `NULL` = „unbekannt", nie `0`; Grants je
  Rolle nach
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 1; jede Nennung der Zeilenzahl trägt „geschätzt"; Richtgröße
  heißt nie „Grenze"/„Limit"/„maximal".
- geprüft, ohne Befund: `SPEC-001`, `SPEC-002` (Feld `origin`, fehlender Wert
  liest als `wal`, letzte Spalte der View), `SPEC-019` (fünf Werte, `applied`
  = „angenommen", eine Transaktion), `SPEC-020`/`SPEC-021` (Ausnahme
  `origin`; `SPEC-024` verweist auf `SPEC-021` und bleibt mit „zehn Feldern"
  wahr), `SPEC-022` (Feld `origin`, Zeile „Herkunft", Zeile „Position und
  `limit`" gegen `SelectChanges` bestätigt).
- geprüft, ohne Befund: Kennungs-Vergabe — `SPEC-029` ist die nächste freie
  (Parent-Maximum `SPEC-028`), keine Kollision; alle neuen Kennungen der Spec
  sind auflösbar verlinkt (Gate `ids`/`anchors`, siehe Gate-Lauf).
- geprüft, ohne Befund: Historie (`spec/pflichtenheft.md` §7) — drei neue
  Zeilen ohne ADR-/Slice-Bezug im Text; die Zeile vom 2026-09-19 („offene,
  ADR-pflichtige Frage") ist ein Zeitpunkt-Record und wahr zu ihrem Datum,
  kein überholter Text.
- geprüft, ohne Befund: Lifecycle und Commit-Struktur — reine Renames
  `open` → `next` und `next` → `in-progress`, Inhalt getrennt; Betreffs ohne
  `SPEC-`/`ARC-`-Kennung, jede Message mit `LH-FA-CAP-009`/`ADR-0111`
  (`make commit-traceability` Exit 0).
- geprüft, ohne Befund: `docs/user/` — Benutzerhandbuch im Diff unberührt
  (kein neuer Betreiber-Pfad ohne Handbuch-Zug, da der Slice keine
  Betreiber-Oberfläche einführt, sondern Spec nachzieht); die Meldungen an
  `slice-backfill-change-origin` (Handbuch Z. 639, 377) und
  `slice-backfill-sql-administration` (Handbuch „Lesen ohne `Limit`",
  Fortsetzungs-Idiome Z. 375 ff./641 ff.) stehen im Suchlauf-Feld §3, und beide
  Empfänger-Pläne führen den Handbuch-Zug bereits im eigenen Umfang
  (`docs/user/benutzerhandbuch.md` in `slice-backfill-change-origin` §3;
  „Lesen ohne `Limit`" in `slice-backfill-sql-administration`, DoD-Punkt zum Handbuch).
- geprüft, ohne Befund: `AGENTS.md` §3.1 Docker-only — der berichtete
  Host-Python-Einsatz für Markdown-Ersetzungen ist ein einmaliges
  Text-Editieren ohne Installation, ohne Artefakt im Repo und ohne Bezug zu
  einem Build-/Test-/Betriebspfad; die Regel verbietet lokale
  Toolchain-Installation für Build/Test/Betrieb, das Werkzeug ist kein
  Bestandteil des Repos. Kein Befund.
- geprüft, ohne Befund: `AGENTS.md` §3.11 (keine host-lokalen Pfade), §3.5
  (keine `Accepted`-ADR im Diff), §3.6 (keine Gate-Lockerung).
- Schwerpunkte laut Skill ohne Anwendungsfall im Diff (kein Code, kein
  Kommentar in Produktionscode, keine Workflow-Datei, kein Handbuch-Inhalt):
  Chronik in Code-Kommentar, Handbuch-Versionshistorie, Betreiber-Oberfläche,
  Docker-only-Skripte, Suppression, Sicherheits-/Korrektheitspfad im Code —
  nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Nachzug lässt überholten Text im selben
Dokument stehen · Suchlauf deckt die Plan-Träger der bewegten Eigenschaft
nicht · Konjunktiv über verworfene Alternative · Unpräzise
Sichtbarkeitsaussage neben einer Gegenstelle · Technik-Zusage ohne
Lastenheft-Anker (zur Deutung) · Zukunfts-Aussage im Indikativ ohne Marker

## Verdikt

**Merge-blockierend:** nein für HIGH (keine); F-1 (MEDIUM) ist ein
Ein-Zeilen-Vorbehalt und wird vor der Closure des Slice behoben — die
Entscheidung über die Fixrunde liegt beim Planner. Die DoD-Zeile „Review
durchgeführt" bleibt deshalb offen und wird nach der Fixrunde regulär
nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug: Grenze
„braucht der Slice eine Fixrunde, bleibt die Checkbox offen").

**Übergabe:** Findings gehen an den Implementer (F-1 bis F-4); F-5 und F-6
gehen als Deutungsfragen an Verifier/Architect. Die **Finding-Klassen**
gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler; die
Klasse „Nachzug lässt überholten Text im selben Dokument stehen" ist mit
diesem Lauf der dritte Beleg von
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` und erreicht die
Schärfungs-Schwelle des Steering-Loops (Modul 10 §Pflege). Dieser Report
selbst ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft
der Verifier separat (Modul 11).
