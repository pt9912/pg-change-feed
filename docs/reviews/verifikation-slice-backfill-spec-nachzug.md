# Verifikations-Report: slice-backfill-spec-nachzug — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Spec-Konformität gegen die
Entscheidungen + Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-backfill-spec-nachzug.md`](review-slice-backfill-spec-nachzug.md);
Hausform dieses Reports:
[`verifikation-slice-sdk-python-http-reale2e.md`](verifikation-slice-sdk-python-http-reale2e.md).

**Gegenstand:** `Slice-Plan` `slice-backfill-spec-nachzug` (Welle
`welle-backfill-bestand`), Diff-Range `77c60bd1..HEAD` (`6d6ee7ea`): Implementer-
Commits `8a47eb3a` … `46d731c5` (drei Lifecycle-/Meta-Commits, zwei Spec-Commits,
zwei Plan-Commits), Review-Report `88bee004`, Fixrunde `1034840f`/`6d6ee7ea`.
9 Dateien, 637 Insertions / 55 Deletions (`git diff --stat 77c60bd1..HEAD`): zwei
Spec-Dateien, der Slice-Plan, fünf Folge-Pläne, der Review-Report. Dieser Lauf
ändert keine Spec, keinen Plan und keinen Code; er schreibt nur diesen Report.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make gates` (vor dem Report-Commit) | **EXIT=0** | Exit-Code separat gesichert (`make gates > <log> 2>&1; echo EXIT=$? > <datei>`), Log danach gelesen. baseline-verify v6.9.0 OK (54 Dateien) · `docs-check` — d-check `967 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability` OK (5 Commits in `HEAD~5..HEAD`, Betreffs ohne Struktur-ID; das `commits`-Modul ebenfalls 0 Befunde) · coverage-gate über die netzlose Fläche `total: (statements) 82.7%`, Exit 0 · generated-sync OK (byte-gleich, Stufe `proto-export`) · a-check `gesamt: 0 Befund(e)` |
| Stempel | **GLEICH** | `bash tools/harness/working-tree-hash.sh` = `.harness/state/gates-passed.diffsha` = `f3e4feb8d107462b8a67d8b7cc800af4fd4f87ecaac65f57ae48184426d2dd6e` (Arbeitsbaum vor dem Report sauber) |
| `make commit-traceability RANGE=77c60bd1..HEAD` | **EXIT=0** | „OK — 10 Commit(s) in `77c60bd1..HEAD`, Betreffs ohne Struktur-ID" |
| `make doc-commits RANGE=77c60bd1..HEAD` | **EXIT=0** | 967 Dateien, 0 Befunde |
| `make doc-immutable RANGE=77c60bd1..HEAD` | **EXIT=0** | 967 Dateien, 0 Befunde (kein `Accepted`-ADR berührt; `git diff --stat` nennt keine Datei unter `docs/plan/adr/`). Ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab (EXIT=2) — Aufruf-Form, kein Befund |
| `make doc-trace` | **EXIT=0** (advisory) | **80 Anforderung(en), 3 Waise(n)** — `LH-FA-CAP-009` (Zeile: ADRs `ADR-0111`, `ADR-0112`, `ADR-0113`; Slices `—`; Coverage `—`; `WAISE`), `LH-FA-CFG-007`, `LH-FA-CFG-008` |

**Zu `make doc-trace`:** [`LH-FA-CAP-009`](../../spec/lastenheft.md) bleibt
`WAISE` — erwartet, denn der Slice ist ein reiner Spec-Zug; die Zeile trägt
weder Slice noch Coverage, bis ein umsetzender Slice die Kennung führt und
`e2e` die Abdeckungs-Zeile schreibt. Die Zählung 80/3 ist **nicht** vom Diff
verursacht: `LH-FA-CFG-008` kam mit Lastenheft 0.13.0 (`712dc26d`), das vor dem
Parent (`77c60bd1`) liegt (`git merge-base --is-ancestor` bestätigt); der Diff
berührt kein Lastenheft und keine der 13 Träger-Zeilen zu
`CAP-009` (siehe §3, Suchlauf-Zeile 5). Die Momentaufnahme „79 Anforderungen,
2 Waisen" in `harness/README.md` (Zeile `make doc-trace`) ist damit veraltet —
Beobachtung V-2.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) beantwortet: Überschrift ohne „offen"; Mechanismus, Markierung, Überlappung, Sichtbarkeits-Grenze, Neubeginn als Zusagen | **erfüllt** | `spec/pflichtenheft.md` Z. 144: `### LH-FA-CAP-009.a — Backfill-Mechanismus` (Parent-Z. 144 trug „… offen"). Der Abschnitt (Z. 144–246) trägt je einen Punkt „Mechanismus" (Bulk-Copy im Snapshot eines je Run angelegten temporären Slots, eine Store-Transaktion), „Markierung" (`origin = backfill`, `INSERT`, vollständiges `new_data`), „Position und Ordnung", „Überlappungs-Verhalten" (Position `X`, keine Lücke, begrenzte idempotente Dopplung), „Sichtbarkeits-Grenze", „Unterbrechung und Neubeginn"; Zukunfts-Form durch den Vorsatz „Die folgenden Sätze sind Zusagen an die Umsetzung, keine Messergebnisse." Kein Satz im Abschnitt behauptet eine Messung; der Diff enthält in den Spec-Dateien keinen Konjunktiv über verworfene Alternativen (Regex-Lauf über die `+`-Zeilen: nur „neu registrierte" und „neu starten/neuem Snapshot" — Sachwörter, keine Chronik). `make gates` grün |
| 2 | Datenstrukturen: [`SPEC-002`](../../spec/pflichtenheft.md) (`origin`), [`SPEC-019`](../../spec/pflichtenheft.md) (`backfill`, fünf Werte, `applied` = „angenommen"), [`SPEC-022`](../../spec/pflichtenheft.md) (`origin`, Positions-Anmerkung), [`SPEC-001`](../../spec/pflichtenheft.md), [`SPEC-029`](../../spec/pflichtenheft.md) neu; nächste freie Kennung; §7 Historie ohne ADR-/Slice-Bezug | **erfüllt** | `SPEC-002` Z. 338–350 (Feld `origin`, geschlossene Menge, fehlender Wert liest `wal`, letzte Spalte der View); `SPEC-019` Z. 511/512 (`request_kind` = fünf Werte, `column_name`-Zelle „die drei übrigen") und Z. 526–535 („`applied` … **angenommen**", eine Transaktion, genau eine `pending`-Zeile); `SPEC-022` Z. 603 (Antwortfeld `origin`), Z. 617 („Herkunft"), Z. 618 („Position und `limit`"); `SPEC-001` Z. 324 (Tabellenliste); `SPEC-029` Z. 668–712 (**genau ein** `### SPEC-029`-Kopf). Kennungs-Vergabe nachgemessen: `git show 77c60bd1:spec/pflichtenheft.md \| grep -o 'SPEC-0[0-9][0-9]' \| sort -u` endet auf `SPEC-028`; HEAD endet auf `SPEC-029`. §7 Historie (Z. 839–841): drei neue Zeilen, `grep` über die `+`-Zeilen des Spec-Diffs nach `ADR-[0-9]`/`slice-`/`welle-` liefert **0 Treffer** |
| 3 | `spec/architecture.md`: Backfill-Sequenz (Annahme in einer Transaktion, Aufnahme beim Start), Rollen in [`ARC-006`](../../spec/architecture.md), ohne ADR-/Slice-/Wellen-Bezug; Antragsarten-Aufzählungen gezogen | **erfüllt** | Abschnitt `### Use-Case: LH-FA-CAP-009.a — Bestand als Backfill überführen` (Z. 268 ff.) mit zwei Sequenzdiagrammen (Annahme + Aufnahme; Ausführung); `ARC-006`-Zeile (Z. 59) und der Knoten „Driven Adapters" (Z. 41) tragen Snapshot-Leser und Annahme; „fünf Arten" + Tabellenzeile `backfill` (Z. 221–230). `grep -nE 'ADR-0[0-9]\|slice-[0-9a-z]\|welle-' spec/architecture.md` → **0 Treffer** (ganze Datei), der Diff der Datei enthält kein `ADR`/`Slice`/`Welle`. `make docs-check` (`matrix`-Modul) grün |
| 4 | `make gates` grün, Exit-Code ungefiltert gesichert | **erfüllt** | eigener Lauf EXIT=0 + Stempel-Gleichheit (§1) |
| 5 | Review durchgeführt, Report unter `docs/reviews/`, kein offenes HIGH/MEDIUM | **erfüllt** | `docs/reviews/review-slice-backfill-spec-nachzug.md` liegt vor (`88bee004`); Summary 0 HIGH / 1 MEDIUM / 3 LOW / 2 INFO — die Zählung im Plan stimmt mit den sechs Findings F-1…F-6 überein. Behebung von F-1…F-4 selbst geprüft (§5) |
| 6 | §3.13-Suchlauf: committetes Feld trägt Gefundenes **und** Nichtgefundenes je Träger, beide Stände gemessen | **erfüllt** | Plan §3 (sieben Zeilen) trägt je Zeile Parent-/Diff-Befund und „Nicht gefunden"; Stichprobe nachgemessen in §3 dieses Reports — alle geprüften Zahlen stimmen |
| 7 | Doku-Update entfällt (Slice **ist** das Doku-Update) | **erfüllt** | `docs/user/` unberührt (`git diff --stat`); die Handbuch-Träger sind als Meldung an `slice-backfill-change-origin`/`-sql-administration` im Suchlauf-Feld geführt |
| 8 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" — Rollen-Sequenz (Closure nach Verifier) |
| 9 | Reconciliation-Register | **korrekt offen** | Zeile trägt „entfällt", Haken folgt bei Closure |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; die Review-Klassen (F-1: `nachzug-laesst-ueberholten-text-stehen`, dritter Beleg; F-2: Suchlauf-Deckung; F-3: Konjunktiv) stehen für den Zähler an — nicht mein Gegenstand |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | alle sechs Risiko-Zeilen tragen „**Ausgang:** *(bei Closure)*" |
| 12 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | hängt an der Closure der Welle (`welle-backfill-bestand` steht offen) |

Kein `[x]` ohne Beleg, kein `[ ]`, das bereits belegt wäre.

## 3. Plan-vs-Code-Diff (§3 + Plan-Nachzug + §3.13-Suchlauf-Feld)

**Geplant und geliefert:** Alle zwölf Plan-Zeilen (§3, Zeilen 144–155) sind im Diff
vertreten: `spec/pflichtenheft.md` §1 (`LH-FA-CAP-009.a`), §2 (`SPEC-001`, `-002`,
`-019`, `-022`, `-029` neu), §7; `spec/architecture.md` §1/§4, Kopf
(`Letzte Änderung: 2026-09-24`) und Diagramm-Knoten; `SPEC-020`/`-021` (je ein
Halbsatz: das Feld `origin` gehört nicht zur Live-Nachricht — Z. 566, 585); die
Warn-Spalten-Festlegung (`warn_estimated_size`, `warn_duration`, beide
`boolean NOT NULL DEFAULT false`); die Fixrunden-Zeilen F-1 (Z. 614), F-3
(`SPEC-029`-Absatz und `architecture.md`), F-4 (`architecture.md` Z. 342), F-2 (fünf
Folge-Pläne). **Ungeplant:** keine Datei außerhalb dieser Liste (die neun Dateien
des `--stat` decken sich mit den Plan-Zeilen plus Review-Report). **Nicht erfolgt:**
nichts aus §1 „Ausdrücklich NICHT" berührt — Lastenheft, Benutzerhandbuch, Code,
Schema und die Transformations-Regeltypen sind unverändert (`git diff --stat`).

**Lifecycle (`AGENTS.md` §3.3):** die beiden Moves sind reine Renames —
`git show --stat -M 8a47eb3a`: `{open => next}/slice-backfill-spec-nachzug.md | 0`,
`git show --stat -M 9e865160`: `{next => in-progress}/… | 0`, je „1 file changed, 0
insertions(+), 0 deletions(-)". Die Inhaltsänderung „Verantwortlich gesetzt"
(`f411b352`) liegt in einem eigenen Commit **zwischen** den Moves. Regelkonform.

**§3.13-Suchlauf-Feld — gegen beide Stände (`77c60bd1` und `HEAD`) selbst
nachgemessen:**

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| 1 (Anker „… offen") | `git grep -n 'backfill-mechanismus' 77c60bd1 -- '*.md'` → 1 Treffer (die Suchbefehl-Zelle des Plans); `HEAD` → 2 Treffer (Plan-Zelle + eine Zeile im Review-Report, der den Suchlauf zitiert); Überschrift Parent Z. 144 „… offen", HEAD ohne; kein Link auf den Anker; `CAP-009.a`-Fundstellen außerhalb der Backfill-Dateien: `ADR-0112` Z. 393/411 (nackte Kennung) und Pflichtenheft-Historie | **bestätigt** (der Zusatztreffer im Report entstand nach dem Plan-Stand und ist Zitat, kein Link) |
| 2 (Antragsarten-Aufzählungen) | `git grep -n 'exclude_column' <Stand> -- spec docs/user harness README.md`: Parent **9**, HEAD **9**; Zählformulierungen („vier/fünf Antragsarten", „drei übrigen", „beiden Tabellen-Antragsarten") in `spec`/`docs/user`/`harness`/READMEs/`AGENTS.md` auf HEAD: nur wahre Restformen („beiden Spalten-Antragsarten", `column_name` „die drei übrigen … (`enable`, `disable`, `backfill`)"); `harness/README.md` Z. 151 „vier SQL-Funktionen" beschreibt die Fremdobjekte des Idempotenz-Guards und bleibt wahr | **bestätigt** |
| 3 (Feldlisten) | `git grep -c 'Felder' <Stand> -- spec/pflichtenheft.md spec/architecture.md`: Parent 11/0, HEAD 11/0; „zehn Felder" in `SPEC-021`/`SPEC-024` und `SPEC-020` tragen nach dem Nachzug die Ausnahme `origin` bzw. verweisen auf `SPEC-021` | **bestätigt** |
| 4 (`ARC-005`/`ARC-006`) | `git grep -n 'ARC-005\|ARC-006' <Stand> -- spec/architecture.md \| wc -l`: Parent **20**, HEAD **24** | **bestätigt** |
| 5 (`CAP-009` in Trägern) | `git grep -n 'CAP-009' <Stand> -- docs harness .d-check.yml`, Pfade ohne `slice-backfill`/`welle-backfill`/`adr/0111-`/`adr/0113-`: Parent **13**, HEAD **13**; Träger: `ADR-0112` (2), ADR-Index, Roadmap, `slice-transformationen-backfill-pfad`, `welle-transformationen`, drei Verifikationsberichte (je 2), `harness/README.md` — kein Träger in `docs/user/` und `.d-check.yml` | **bestätigt** |
| 6 (Fortsetzungs-Aussage, F-1) | `git grep -nE 'letzte gelieferte\|letzte.*commit_position\|\+ 1' <Stand> -- spec`: Parent 1 Treffer (Z. 506), HEAD 2 (Z. 614 mit Vorbehalt, Z. 618 „Position und `limit`"); `spec/architecture.md` 0 Treffer | **bestätigt** |
| 7 (Plan-Träger der Warn-Spalten/Antragsarten, F-2) | Muster aus der Plan-Zelle über `docs/plan/planning/open` und `docs/plan/planning/welle-*.md`: Parent **9** Zeilen in 5 Dateien (`run-store` 1, `run-usecase` 2, `sql-administration` 3, `transformationen-antragsweg-schema` 1, `transformationen-spec-nachzug` 2), `88bee004` identisch, HEAD **0** | **bestätigt** |

Die fünf Folge-Pläne wurden dabei **direkt geändert** statt nur gemeldet (`AGENTS.md`
§3.13: „gemeldet statt still mitgeändert"). Die Änderung ist nicht still: die Plan-
Zeilen 154/155 und die Suchlauf-Zeile 7 nennen sie, und der Diff ist minimal („leer"
→ „`false`", „vier" → ohne Zahl; die Absicht der Pläne bleibt unberührt — von mir am
Diff `88bee004..HEAD` der fünf Dateien gelesen). Form-Anmerkung, kein Befund.

## 4. Spec-Konformität gegen die Entscheidungen

Gegenstand: `LH-FA-CAP-009.a`, [`SPEC-001`](../../spec/pflichtenheft.md),
[`SPEC-002`](../../spec/pflichtenheft.md), [`SPEC-019`](../../spec/pflichtenheft.md),
[`SPEC-020`](../../spec/pflichtenheft.md), [`SPEC-021`](../../spec/pflichtenheft.md),
[`SPEC-022`](../../spec/pflichtenheft.md), [`SPEC-029`](../../spec/pflichtenheft.md)
und die Architektur-Sicht gegen [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
(`Accepted`), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
(`Accepted`, Supersedes `ADR-0111` Teilfrage 5 teilweise) und
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(nur Abgrenzung).

| Festlegung der Entscheidung | Spec-Stelle | Verdikt |
|---|---|---|
| Warn-Spalten `warn_estimated_size`/`warn_duration`, `boolean NOT NULL DEFAULT false`; Ein-/Zwei-Spalten-Wahl liegt beim Spec-Nachzug (`ADR-0113` Festlegung 3 Punkt 4) | `SPEC-029` Z. 687/688: beide Spalten, Typ und Default wörtlich; (1) schreibt die Annahme, (2) der Worker, geprüft „je Block oder beim Abschluss" (deckt `ADR-0113` Festlegung 3 Punkt 3) | konform |
| Grants: `cdc_admin` `SELECT`/`INSERT`; `cdc_capture` `SELECT`/`UPDATE`; `cdc_reader` nur View; niemand `DELETE` (`ADR-0113` Festlegung 1 Punkt 3) | `SPEC-029` Z. 697–706: Tabelle und „Niemand trägt `DELETE`; `cdc_admin` trägt kein `UPDATE`, `cdc_capture` kein `INSERT`"; `cdc_reader` „**keines**" auf die Basistabelle, `SELECT` auf `cdc.backfill_status` | konform (wörtlich) |
| `estimated_rows` `NULL` = „unbekannt", nie `0`; „geschätzt" an jeder Nennung (`ADR-0113` Festlegung 3 Punkt 5/6) | `SPEC-029` Z. 686 und Absatz zur View (Z. 708–712); `CAP-009.a` „**geschätzte** Zeilenzahl"; Sicht „Zeilenzahl schätzen" | konform |
| **keine** Toleranz, **keine** Richtgröße, **keine** „10 Minuten" in der Spec (`ADR-0113` Festlegung 3 Punkt 1/2; nicht in `spec/pflichtenheft.md` §3) | `grep -nE '10 Minuten\|zehn Minuten\|Toleranz\|Richtgröße' spec/*.md`: nur Verweis-Sätze (Z. 242–244, 687–692: „sind keine Konstanten dieses Dokuments"), **kein Zahlenwert**; Richtgröße heißt nie „Grenze/Limit/maximal" | konform |
| `applied` bei `backfill` = „angenommen", Run-Zeile in derselben Transaktion (`ADR-0113` Festlegung 1 Punkt 1/2) | `SPEC-019` Z. 526–535; `CAP-009.a` Punkt „Annahme und Aufnahme"; Sicht Z. 268 ff. | konform |
| Fünf `request_kind`-Werte (`ADR-0111` Teilfrage 5) | `SPEC-019` Z. 512 und Sicht Z. 221–230 | konform |
| Aufnahme beim Prozessstart + Wecksignal, `ORDER BY requested_at, run_id`, Start-Reihenfolge, Neu-Prüfung von Bindung/Publication (`ADR-0113` Festlegung 2) | `CAP-009.a` Punkt „Annahme und Aufnahme"; `SPEC-029` Z. 681 („Aufnahme-Ordnung ist (`requested_at`, `run_id`)"); Sicht Z. 307–311 | konform |
| Mechanismus, Markierung, Position `X`, Block = Transaktion `0bf-<run-id>-<Blocknummer>`, Überlappung, Sichtbarkeits-Grenze, Ein-Transaktions-Neubeginn, Fail-closed, run-lokale Fehler, Retention (`ADR-0111` Teilfragen 1–7) | `CAP-009.a` Z. 144–246: Punkt für Punkt gegen die ADR gelesen, keine Abweichung, kein zusätzlicher Wert | konform |
| Live-Wege unverändert (`ADR-0111` Teilfrage 8) | `SPEC-020` Z. 566, `SPEC-021` Z. 585: `origin` „gehört nicht zur Nachricht"; `SPEC-002` „Die Live-Wege … tragen es nicht"; kein Proto-/`gen/`-Diff | konform |
| Sicht ohne ADR-/Slice-/Wellen-Bezug (`AGENTS.md` §3.4) | 0 Treffer, siehe §2 Zeile 3 | konform |
| Lastenheft unberührt; „präzisieren ja, erweitern nie" | `git diff --stat`: `spec/lastenheft.md` nicht im Diff. Inhaltlich: Happy Path (Bestand „zum Startzeitpunkt", Herkunft erkennbar), Boundary (Zeilen mit WAL-Änderung zwischen Start und Ende: Commit ≤ `X` gedoppelt und begrenzt, Commit > `X` sortiert dahinter), Negative (Neubeginn ohne Verlust) sind in `CAP-009.a` abgedeckt | konform (Deutung F-5 siehe §6) |

Ein Widerspruch zwischen `ADR-0111` und `ADR-0113` wird in der Spec nicht
übernommen: die Spec folgt an beiden Superseded-Stellen (`SPEC-029` Grants:
`cdc_capture` **nicht** `INSERT`; `applied` = „angenommen") der jüngeren ADR.
`ADR-0071` wird nicht berührt (kein Coverage-/Gate-Bezug in der Spec-Änderung).

## 5. Fixrunde — F-1…F-4 behoben?

Die Stellen selbst gelesen (Diff `88bee004..HEAD` und Stand `HEAD`), nicht dem Bericht
geglaubt:

- **F-1 (MEDIUM, `SPEC-022` „Reihenfolge"):** behoben. Z. 614 trägt „die Fortsetzung
  ist `from = <letzte gelieferte commit_position> + 1`, **wenn** das Lesen die letzte
  Position vollständig erfasst hat; enthält sie mehr Changes als `limit`, gilt die Zeile
  „Position und `limit`"". Z. 618 sagt: bei mehr Changes als `limit` liefert genau diese
  Fortsetzung den Rest der Position nicht. Kein Widerspruch mehr; die Aussage ist
  am Code gegengeprüft — `queries.SelectChanges` filtert
  `t.commit_position >= $2` / `< $3` und schneidet mit `LIMIT $6`
  (`internal/adapters/driven/postgresstorage/queries/queries.go`, Z. 65–70), also Zeilen,
  nicht Positionen. Ebenso konsistent: `CAP-009.a` „Sichtbarkeits-Grenze" (Z. 215–218)
  und Sicht (kein `limit`-Satz). **Kein Restwiderspruch.**
- **F-2 (LOW, Suchlauf deckt Plan-Träger nicht):** behoben. Muster-Lauf über die Plan-
  Träger: HEAD 0 Treffer (Parent 9). Neue Suchlauf-Zeile 7 im Plan (bestätigt, §3).
  Restbeobachtung V-1.
- **F-3 (LOW, Konjunktiv):** behoben. `SPEC-029` (Absatz nach der Tabelle): „Es sind zwei
  Spalten mit je einem Schreiber … jede Spalte hat deshalb genau eine schreibende
  Rolle." Sicht (Z. 301–304): „Annahme und Antragsvermerk sind eine Transaktion: nach
  einem Absturz … bleibt entweder beides bestehen oder beides aus". Regex-Lauf über
  alle `+`-Zeilen des Spec-Diffs nach `hätte|wäre|würde|bliebe|stünde|sonst` → keine
  Konjunktiv-Treffer.
- **F-4 (LOW, Sichtbarkeit):** behoben. Sicht Z. 342: „Sichtbar für Leser ist bis zum
  Commit nur der Fortschritt der Run-Zeile, die kopierten Changes werden mit dem einen
  Commit sichtbar." Das deckt sich mit `SPEC-029` `rows_copied` („außerhalb der
  Daten-Transaktion fortgeschrieben und damit für Leser sichtbar") und mit `CAP-009.a`
  „Alle Blöcke werden in **einer** Store-Transaktion geschrieben, die einmal am Ende …
  committet". **Kein Restwiderspruch** zwischen Sicht und `SPEC-029`; `grep -nE
  'nichts (für\|sichtbar)\|bis zum (einen )?Commit' spec/*.md` findet nur noch die
  präzisierte Stelle.

## 6. Deutungsfragen aus dem Review (F-5, F-6) — Urteil

### F-5 — Warnungen für große Tabellen ohne Lastenheft-Anker: **kein Befund**

Die zwei Warnungen sind eine **zulässige Präzisierung** von
[`LH-FA-CAP-009`](../../spec/lastenheft.md), keine neue Zusage, die einen
Lastenheft-Anker bräuchte. Begründung am Lastenheft-Wortlaut:

1. **Kein Kriterium wird erweitert.** `LH-FA-CAP-009` verlangt im Happy Path
   „jede zum Startzeitpunkt vorhandene Zeile … erkennbar als Backfill-Herkunft", in der
   Boundary ein definiertes Verhalten im Überlappungsfenster, in der Negative den
   Neubeginn ohne stillen Verlust. Die Warnung ändert an keinem dieser Ausgänge etwas
   („**keine Ablehnung, kein Abbruch, keine Statusänderung**", `CAP-009.a` Punkt „Große
   Tabellen"; `SPEC-029`: eine gesetzte Warnung ändert weder `status` noch den Ablauf).
   Ein Consumer, der nur die Lastenheft-Kriterien prüft, sieht sie nie. „Erweitern"
   im Sinn des Pflichtenheft-Kopfs („präzisieren ja, erweitern nie") hieße, ein
   abnahmerelevantes Verhalten hinzuzufügen; das geschieht hier nicht.
2. **Das Lastenheft schließt die Frage ausdrücklich aus dem Muss.** Out-of-Scope:
   „Parallelisierung/Durchsatz eines Backfills über sehr große Tabellen hinweg ist eine
   Ausbaustufe, keine Voraussetzung dieser Anforderung." Große Tabellen sind also weder
   gefordert noch verboten; die Warnung ist die Betriebs-Begleitung eines Nachteils der
   gewählten Ein-Transaktions-Form (`ADR-0111` Konsequenzen), nicht eine Fähigkeit, die
   das Lastenheft nennt. Ein Lastenheft-Anker würde einen Abnahme-Gegenstand
   behaupten, der keiner ist.
3. **Ein indirekter Anker existiert.** Die Warn-Spalten sind Zustand, der über die View
   `cdc.backfill_status` und die CLI-Diagnose sichtbar wird — dieselbe Klasse wie
   [`LH-FA-ADM-002`](../../spec/lastenheft.md)/[`003`](../../spec/lastenheft.md)
   (Betriebsstatus, sichtbare Zustände) und [`LH-FA-SST-003`](../../spec/lastenheft.md)
   (Diagnose-Abdeckung); `SPEC-029` zitiert `LH-FA-SST-003` bereits.
4. **Die Entscheidung liegt beim Auftraggeber.** `ADR-0111` „Festlegungen des
   Auftraggebers" Nr. 3 („Es gibt eine Warnung, keine Ablehnung") und `ADR-0113`
   Festlegung 3 sind `Accepted`; das Pflichtenheft ist laut Kopf „fortschreibbar ohne
   Change Request; eine ADR darf sie schärfen".

**Maßnahme:** keine. Würde der Auftraggeber die Warnung später zum Vertrag machen
(Abnahme, Anforderung), ist das der in `ADR-0113` benannte Re-Evaluierungs-Trigger
(Folge-ADR bzw. Lastenheft-Änderung), nicht ein Nachtrag in diesem Slice.

### F-6 — Sicht beschreibt noch nicht gebaute Rollen im Indikativ: **kein Befund**

Das passt zum Greenfield-Stil der Datei; ein Zukunftsmarker wäre sogar regelwidrig.

1. **Belegter Stil.** Die Sicht wurde mit `f734a6b2` („Architektur-Sicht 0.1", 2026-09-09)
   **vor** dem ersten Go-Bootstrap (`88e1517c`, derselbe Tag) angelegt; sie beschreibt
   Application-Schicht, Consumer-Verwaltung, Retention, Konfiguration und alle
   Driven Adapters im Indikativ, obwohl nichts davon gebaut war. Das ist die
   Modus-Deklaration der `harness/conventions.md` („Greenfield: Doc führt").
2. **Die Regel verbietet den Marker.** `spec/architecture.md` Kopf und `AGENTS.md` §3.4:
   keine Wellen, Slices, Commit-Hashes, Closure-Daten, keine Historie; `Letzte
   Änderung` ist ein Frische-Marker, kein Protokoll. „Noch nicht gebaut / folgt in
   Slice x" ist ein Zeit-/Meilensteinbezug. Die Sicht ist derivativ; die
   Zusage-Kennzeichnung („Zusagen an die Umsetzung, keine Messergebnisse") trägt das
   Pflichtenheft, auf dem die Sicht aufsetzt.
3. **§3.12 Instanz B trifft nicht.** Die Sicht behauptet keine Tatsache über den
   Gegenstand mit Beleg-Anspruch; sie beschreibt eine Ziel-Zerlegung. Eine erfundene
   Messung steht nirgends.

**Maßnahme:** keine.

## 7. Weitere Beobachtungen (nicht blockierend)

- **V-1 (INFO) — Rest-Hedge „Warn-Spalte(n)" in den Folge-Plänen.** Die Zahl steht
  mit `SPEC-029` fest (zwei Spalten, benannt). Die Folge-Pläne tragen weiter die
  Form „Warn-Spalte(n)"/„Spalte(n)": `slice-backfill-run-store` (Z. 21, 40, 84),
  `slice-backfill-sql-administration` (Z. 21, 68, 70, 77, 120, 121),
  `slice-backfill-bench-richtgroesse` (Z. 62, 146, 147),
  `welle-backfill-bestand` (Z. 251, 253) — gemessen mit
  `git grep -nE 'Warn-Spalte\(n\)|Spalte\(n\)' HEAD -- docs/plan/planning`. Das ist
  mit den zwei festgelegten Spalten vereinbar und nirgends falsch; die
  Suchlauf-Zeile 7 des Slice suchte nach „leer"-Formen, nicht nach der Zahl. **An den
  Planner:** bei der Erst-Sichtung der Folge-Slices auf „die zwei Warn-Spalten
  (`SPEC-029`)" ziehen (Meldung, keine Änderung durch mich).
- **V-2 (INFO) — `harness/README.md` Momentaufnahme `make doc-trace`.** Die Zeile
  nennt „79 Anforderungen, **2 Waisen** — `LH-FA-CAP-009`/`LH-FA-CFG-007`"; real
  ausgegeben wird heute **80 Anforderung(en), 3 Waise(n)** (zusätzlich
  `LH-FA-CFG-008`, Lastenheft 0.13.0). Die Drift ist **älter als dieser Slice** (0.13.0
  liegt vor dem Parent) und die Datei liegt außerhalb seines Diffs; der Plan (§3,
  Suchlauf-Zeile 5) weist sie ausdrücklich dem Slice zu, der sie ändert. Der
  Waisen-Stand ändert dieser Slice nicht (`CAP-009` bleibt `WAISE`). **An den Planner:**
  Nachzug in einem passenden Zug (zuständig ist, wer die Zahl bewegt; hier der Zug
  zu `LH-FA-CFG-008`).
- **V-3 (INFO) — „neu registrierte" Consumer ohne „erwartet"-Kennzeichnung.** `CAP-009.a`
  „Sichtbarkeits-Grenze" nennt neu registrierte Consumer als „vor `X`"; `ADR-0111`
  Festlegung des Auftraggebers Nr. 1 führt die Startposition eines frisch
  registrierten Consumers als „erwartet, nicht geprüft". Der Code stützt die Aussage
  (`internal/application/usecase/position/service.go` Z. 26: der Nullwert liest die
  „definierte Anfangsposition eines Consumers ohne Bestätigung"), und die Messung ist
  bereits Träger des Slice `slice-backfill-e2e` (Titel und DoD: „Startposition
  gemessen und im Handbuch mit ihrem Lauf genannt"). Kein Handlungsbedarf in diesem
  Slice; die Zusagen-Form des Abschnitts trägt den Satz.

## 8. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (7) | **7 von 7 erfüllt**, je mit eigenem Beleg |
| DoD §2 — `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz) |
| Spec-Konformität ggü. `ADR-0111`/`ADR-0113` | **konform**, keine Abweichung, kein zusätzlicher Wert |
| Plan-vs-Code-Diff | **deckungsgleich**, keine ungeplante Datei, keine still gestrichene Plan-Zeile |
| Harte Regeln (`AGENTS.md` §3.3, §3.4, §3.7, §3.11, §3.12, §3.13) | **erfüllt** (§3.11 durch `docs-check` 0 Befunde; §3.12/§3.13 nachgemessen) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=77c60bd1..HEAD`) |
| Fixrunde F-1…F-4 | **behoben**, kein Restwiderspruch |
| Gates | **`make gates` EXIT=0**, Stempel gleich; `make doc-trace`: `LH-FA-CAP-009` `WAISE` (erwartet) |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers und die Spec-Zusage sind durch
eigene Messung belegt; die Spec trägt genau die Festlegungen von `ADR-0111` und
`ADR-0113`, ohne Toleranz, Richtgröße oder Zahlenwert, ohne ADR-/Slice-/Wellen-Bezug in
Sicht und Pflichtenheft. Kein V-Befund mit Blockade; V-1…V-3 sind INFO-Beobachtungen
an den Planner. F-5 und F-6: **kein Befund, keine Maßnahme**.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure (Notiz mit Lerneintrag,
Beobachtungs-Register, §6-Ausgänge, Welle-Paarungen); danach der reine `git mv` nach
`done/` (`AGENTS.md` §3.3, Fall 2).
