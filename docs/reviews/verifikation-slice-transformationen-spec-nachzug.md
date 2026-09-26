# Verifikations-Report: slice-transformationen-spec-nachzug — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich +
Spec-Konformität gegen die Entscheidungen + Plan-vs-Diff + Gates.
Review-Artefakt des Reviewers:
[`review-slice-transformationen-spec-nachzug.md`](review-slice-transformationen-spec-nachzug.md);
Hausform dieses Reports:
[`verifikation-slice-backfill-spec-nachzug.md`](verifikation-slice-backfill-spec-nachzug.md).

**Gegenstand:** Slice-Plan `slice-transformationen-spec-nachzug` (Welle
`welle-transformationen`), Stand `HEAD` = `26d6fde7`, Diff-Range
`1718546b..HEAD`. Enthalten: drei Lifecycle-/Meta-Commits, zwei Spec-Commits
(`71627d67`, `8e977328`), drei Plan-Commits, der Review-Report `a32a1931`
(kein Slice-Inhalt). Fünf Dateien (`git diff --stat 1718546b..HEAD`): die
Pflichtenheft-Datei (+291), die Architektur-Sicht (+40), der Slice-Plan (Move
`open/` → `in-progress/` samt Inhalt), der Review-Report. Dieser Lauf ändert
keine Spec, keinen Plan und keinen Code; er schreibt nur diesen Report.

## 1. Eigene Sensor-Belege (dieser Lauf, Exit-Code direkt in Log-Dateien — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Beleg aus meinem Lauf |
|---|---|---|
| `make gates` | **0** | baseline-verify v6.9.0 OK (54 Dateien) · d-check `1193 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability` OK (5 Commits in `HEAD~5..HEAD`, Betreffs ohne Struktur-ID) · coverage-gate `Coverage 83.40% erfüllt Schwelle 80%` · generated-sync OK · a-check `gesamt: 0 Befund(e)` |
| `make docs-check` | **0** | `d-check: 1193 Datei(en) geprüft, 0 Befund(e)` (acht Module) |
| `make commit-traceability RANGE=1718546b..HEAD` | **0** | „OK — 9 Commit(s) in 1718546b..HEAD, Betreffs ohne Struktur-ID" |
| `make doc-commits RANGE=1718546b..HEAD` | **0** | 1193 Dateien, 0 Befunde |
| `make doc-immutable RANGE=1718546b..HEAD` | **0** | 1193 Dateien, 0 Befunde (kein `Accepted`-ADR im Diff berührt) |
| `make doc-trace` | **0** (advisory) | `80 Anforderung(en), 2 Waise(n)` — [`LH-FA-CFG-007`](../../spec/lastenheft.md) (ADRs `ADR-0112`, `ADR-0117`; Slices `—`; Coverage `—`) und [`LH-FA-CFG-008`](../../spec/lastenheft.md) (`ADR-0112`) — erwartet: der Slice ist ein reiner Spec-Zug, die Waisen bleiben bis zu den umsetzenden Folge-Slices |
| Spec-Ebene-Mutation | **2** | angehängte Zeile mit nackter Kennung `LH-FA-CFG-007` am Ende der Pflichtenheft-Datei → `make docs-check`: `1 Befund(e)`, `spec/pflichtenheft.md:1138 LH-FA-CFG-007 id-unlinked`; per `git checkout -- spec/pflichtenheft.md` zurückgenommen, `git status` sauber |
| Docker-Volumes | — | dangling-Volumes (`docker volume ls -qf dangling=true`): **34 vor, 34 nach** dem Lauf; kein Prune |

Lifecycle ([`AGENTS.md`](../../AGENTS.md) §3.3): `git show --stat -M e758fc32` und
`ac8f757c` zeigen je `{open => next}` bzw. `{next => in-progress}` mit
`0 insertions, 0 deletions` — reine Moves. Die Inhaltsänderung „Verantwortlich"
(`69addf43`) liegt in einem eigenen Commit dazwischen.

## 2. DoD — Verdikt je Zeile (Plan §2)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | [`LH-FA-CFG-007.a`](../../spec/pflichtenheft.md) beantwortet, Zukunftsform | **erfüllt** | `spec/pflichtenheft.md:275` `### LH-FA-CFG-007.a — Transformationsform` (Parent `:268` trug „offen"). Der Abschnitt trägt den Vorsatz „Die folgenden Sätze sind Zusagen an die Umsetzung, keine Messergebnisse" und je einen Punkt für Konfigurationsmechanismus, Ausdrucksform, Auswertungsreihenfolge, Mehrdeutigkeit, Verhältnis zum Ausschluss, Wirkort, nicht anwendbare Regel, Abhilfe (Zusage), Abgrenzung. Die Abhilfe schließt mit „Diese Abfolge muss die Umsetzung liefern; sie ist erst mit dem Beleg am laufenden System eine Tatsache" (deckt „erwartet, nicht am Code belegt" der ADR). Chronik-Suche über die `+`-Zeilen des Spec-Diffs (`früher|vorher|bisher|nicht mehr|inzwischen|…`): ein Treffer „die zuvor nicht bestätigte Transaktion" — Sachwort, keine Chronik |
| 2 | Datenstrukturen: `SPEC-019`, Regelform, `SPEC-002`, Zeile `schema`, Run-Satz; nächste freie Kennung; §7 ohne ADR-/Slice-Bezug | **erfüllt** | `SPEC-019`: `request_kind` mit sieben Werten (`:616`), Spalten `rule_name`/`rule_spec` (`:614`/`:615`), zweite Bedeutung von `applied` und Ordnung `requested_at`, dann `administration_request_id`, K1–K4, Fehlertext-Tabelle. `SPEC-030` neu (`:891` ff.). `SPEC-002`-Absatz zur Schlüsselmenge. Zeile `schema` in der Fehlerklassen-Tabelle (`:1020`) plus Absatz „Nicht anwendbare Regel (Klasse `schema`)". Run-Satz in `LH-FA-CAP-009.a` („Sichtbarkeit und Fehler des Runs"). Kennung nachgemessen: `git grep -h -o 'SPEC-0[0-9][0-9]' 1718546b -- spec/pflichtenheft.md \| sort -u \| tail -1` endet auf `SPEC-029`; `SPEC-030` ist die nächste freie und genau einmal als Kopf vergeben. §7: fünf neue Zeilen (`:1132`–`:1136`); `git diff 1718546b..HEAD -- spec \| grep '^+' \| grep -c 'ADR-0\|slice-\|welle-'` druckt **0** |
| 3 | `spec/architecture.md`: Antragsart-Tabelle mit den zwei weiteren Arten, Aussage zur Regelform vor der Persistierung, ohne ADR-/Slice-/Wellen-Bezug | **erfüllt** | Tabelle `:258`–`:267` („sieben Arten", `set_transformation`/`remove_transformation` mit Inbound Port); Absatz „Form der Row Images" (`:133`) mit der Reihenfolge Ausschluss, dann Regeln, vor Serialisierung und Persistierung; `ARC-001`-Zeile nennt „Transformationsregel"; Backfill-Sequenz und §5-Zeile gezogen. Dieselbe Null-Messung wie in Zeile 2 (der Diff der Sicht trägt keinen `ADR`/`Slice`/`Welle`-Bezug); `make docs-check` (`matrix`-Modul) 0 Befunde |
| 4 | `make gates` grün, Exit-Code ungefiltert und gesondert ausgewertet | **erfüllt** | eigener Lauf Exit 0 (§1). Der Plan-Beleg „Exit 0 am Stand nach dem Plan-Nachzug-Commit" stimmt mit meiner Messung am Endstand überein |
| 5 | Review durchgeführt, kein offenes HIGH/MEDIUM | **erfüllt** | `docs/reviews/review-slice-transformationen-spec-nachzug.md` (`a32a1931`): 0 HIGH / 1 MEDIUM / 5 LOW / 5 INFO; das MEDIUM F-1 ist in der Fixrunde behoben (Beleg in §5) |
| 6 | §3.13-Suchlauf: Gefundenes **und** Nichtgefundenes je Träger, beide Stände gemessen | **erfüllt, mit zwei Zahl-/Lokator-Drifts und einer falschen Träger-Aussage** (V-1, V-2) | Die Felder stehen (zwei Tabellen: Erst-Lauf und Fixrunde). Nachgefahren: §4 dieses Reports. Sieben von zwölf Zeilen stimmen exakt; drei Angaben sind überholt oder nicht zutreffend |
| 7 | Doku-Update entfällt (der Slice ist das Doku-Update) | **erfüllt** | `git diff 1718546b..HEAD -- docs/user harness \| wc -l` druckt **0**: Handbuch und Harness-Dateien unberührt, die Träger sind als Meldung geführt (§4) |
| 8 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(zu tragen bei Closure)*" — Rollen-Sequenz (Closure nach Verifier) |
| 9 | Reconciliation-Register | **erfüllt** | `entfällt` (Greenfield), Haken gesetzt und begründet |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; V-1 ist ein Kandidat für einen weiteren Beleg (§6) |
| 11 | Risiken aus §6 mit Ausgang | **korrekt offen** | alle sechs Risiko-Punkte tragen „**Ausgang:** *(bei Closure)*" |
| 12 | Drei Paarungen | **korrekt offen** | hängt an der Closure der Welle `welle-transformationen` (offen) |

Kein `[x]` ohne Beleg; kein `[ ]`, das bereits belegt wäre.

## 3. Konformität gegen die Entscheidungen

Gegenstand: die geänderten Stellen von `spec/pflichtenheft.md` und
`spec/architecture.md` gegen
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md),
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (nur Muster/Bezeichner)
und [`LH-FA-CFG-007`](../../spec/lastenheft.md)/[`LH-FA-CFG-008`](../../spec/lastenheft.md).

| Festlegung der Entscheidung | Spec-Stelle | Verdikt |
|---|---|---|
| Regeltypen `rename_column` (`column`, `to`) und `map_value` (`column`, `values`), Pflichtschlüssel `kind`, unbekannte Schlüssel/`kind` enden `failed` (`ADR-0112` Teilfrage 2) | `SPEC-030` Tabelle; Fehlertext-Zeilen „unbekannter Regeltyp", „unbekannter Schlüssel in rule_spec" | konform (Wortlaut der Wirkungs-Spalte gleich) |
| Abwesenheit bleibt Abwesenheit; Wirkung nur auf Schlüssel und Werte; Tabellen-Identität bleibt Quell-Identität | `SPEC-030` „Abwesenheit", „Wirkung nur auf das Bild"; `LH-FA-CFG-007.a` Wirkort (Feldliste `change_id` … `schema_version`) | konform |
| K1–K4 mit Wirkung „`failed`, Regelstand unverändert" (Teilfrage 3) | `SPEC-019` K1–K4 (K3 einschließlich ausgeschlossener Spalten und der Quellspalte selbst; K4 einschließlich nicht geführtem Regelnamen bei `remove_transformation`) | konform; Fehlertext-Wortlaute und Prüfreihenfolge sind Festlegungen der Spec über die ADR hinaus (siehe unten) |
| Ausschluss vor Regelstand, `cdc.exclude_column`/`include_column` bleiben zulässig (Teilfrage 5) | `LH-FA-CFG-007.a` „Verhältnis zum Spaltenausschluss" | konform (wörtlich) |
| Wirkort vor der Persistierung, eine Stelle für alle Erzeugungspfade; Rohform nicht gespeichert, Regeländerung wirkt nur nach vorn, `map_value` nicht umkehrbar (Teilfrage 6, Konsequenzen) | `LH-FA-CFG-007.a` „Wirkort"; Sicht „Form der Row Images" | konform |
| Nichtanwendbarkeit: Erfassungspfad endet sichtbar mit `schema`, keine Persistierung, keine Bestätigung, Fortsetzung ab bestätigter Position (Teilfrage 4) | `LH-FA-CFG-007.a` „Nicht anwendbare Regel"; `SPEC-008` Zeile `schema` | konform |
| Abhilfe-Akzeptanzkriterium (Folgepflicht 5 (a)–(d)): Regel entfernen (Antrag bleibt offen, solange der Prozess steht), Neustart, Antrag vor der ersten Assemblierung `applied`, zuvor unbestätigte Transaktion erscheint, keine zweite Neustart-Schleife | `LH-FA-CFG-007.a` „Abhilfe (Zusage)": `pending` statt `requested` (Bestand, `nacharbeit-administration.sql`), offene Anträge „**bevor** die erste Transaktion der Tabelle assembliert wird", „ohne dass der Prozess erneut an derselben Regel endet"; als Zusage, nicht Tatsache, formuliert | konform |
| `ADR-0117` Festlegung 1/2/3: Klasse `schema` im Run, Prüfung gegen die Snapshot-Spalten vor Schreibtransaktion und erster Zeile, run-lokal, sichtbar in `cdc.backfill_status` und Diagnose | `LH-FA-CAP-009.a` („bevor die erste Zeile gelesen wird … ohne Change und ohne den Erfassungspfad zu berühren"); `SPEC-008` Zeile und Absatz; Sicht Backfill-Sequenz | konform |
| `ADR-0117` Festlegung 4: Abhilfe = Regelstand ändern, dann **neuer** Antrag `cdc.backfill_table`; ein `failed`-Run wird nicht fortgesetzt | `SPEC-008` Absatz „Nach der Abhilfe im Run beginnt ein **neuer** Antrag …" | konform (siehe V-5: „Neustart" nicht ausdrücklich ausgeschlossen) |
| `ADR-0117` Festlegung 5: Regelstand-Wechsel im Lauf endet `configuration` | `LH-FA-CAP-009.a` „Fail-closed vor dem Commit" (`configuration`); im Bestand bildet `classifyError` die Ausschluss-Abweichung auf `configuration` ab (`internal/application/usecase/backfill/service.go`) | konform |
| `ADR-0059`: der Bestand vergleicht Spaltennamen zeichengenau | `SPEC-030` „Bezeichner", `SPEC-019` K3 | konform, Bestand nachgeprüft (unten) |
| Lastenheft unberührt; „präzisieren ja, erweitern nie" | `git diff --stat`: `spec/lastenheft.md` nicht im Diff. Happy Path (Regel wirkt auf die ausgelieferte Form), Boundary (Mehrdeutigkeit definiert statisch), Negative (sichtbarer Fehlerzustand statt unveränderter/verworfener Auslieferung) sind in `LH-FA-CFG-007.a` abgedeckt; das Routing bleibt bei `LH-FA-CFG-008.a` (Abgrenzung) | konform |

**Festlegungen des Slice gegen den Bestand nachgeprüft** (Plan §3 Tabelle):

| Festlegung | Bestands-Beleg aus meinem Lauf | Verdikt |
|---|---|---|
| `column` zeichengenau, keine Faltung, kein Quoting | `internal/domain/model/rowimage.go:60` `containsName`: `candidate == name`; `queries.SelectTableColumnExists`: `column_name = $3` gegen `information_schema.columns` (`internal/adapters/driven/postgresstorage/queries/queries.go`); `ColumnExists` reicht den Namen als Parameter durch | bestätigt |
| `column`: nichtleer; U+0000 abgelehnt | `NewAdministrationRequest` (`internal/domain/model/administrationrequest.go`) verlangt nur nichtleere Spalte für die Spalten-Arten; die U+0000-Begründung (Katalog-Parameter scheitert) ist keine Bestandsaussage, sondern eine Festlegung — als solche im Plan geführt | bestätigt (Festlegung, nicht Bestand) |
| `rule_name` `[a-z0-9_]{1,63}` | `identifierShape = ^[a-z0-9_]{1,63}$` (`internal/adapters/driven/postgresstorage/tableactivation.go:28`), angewandt auf Schema und Tabelle (`:350`) | bestätigt |
| `to` ≤ 63 Byte, nicht leer, ohne U+0000 | Bestand prüft keinen Zielnamen (`grep` über den Baum: kein Treffer), 63 = Bezeichner-Länge der Quelle | bestätigt (Festlegung) |
| `column_name` bleibt NULL | `NewAdministrationRequest` kennt nur die zwei Spalten-Arten mit Spalte; die Tabelle trägt `column_name` bereits nullable | bestätigt |
| Status `pending` statt `requested` | `tools/schema/nacharbeit-administration.sql` (`'pending'` in den Funktionen); `ADR-0112` Folgepflicht 5 (b) nennt `requested` — Wortlaut-Versehen der `Accepted` ADR, die Spec folgt dem Bestand | bestätigt |
| `request_kind`: Menge des Parent-Stands plus zwei Werte | `nacharbeit-administration.sql:62`: CHECK mit fünf Werten (`enable`, `disable`, `exclude_column`, `include_column`, `backfill`); Spec `:616`: sieben | bestätigt |
| Leeres `values` `failed`; Zielname gleich Quellname `failed` (folgt aus K3); Abbildung auf sich selbst zulässig; Mehrfachabbildung zulässig | Kein Bestand; `ADR-0112` Teilfrage 3 (K3: „jeder Spaltenname der Quelltabelle") und Konsequenzen (Mehrfachabbildung erwartet) tragen sie; keine widerspricht einer ADR-Aussage | bestätigt |
| Schlüsselreihenfolge im Image nicht zugesagt | `SPEC-002` (`jsonb`), Bestand `BuildRowImage` baut den Text in Relation-Spaltenreihenfolge — der Assembler-Ausgang hat eine Reihenfolge, der Lesepfad nicht. **Weicht in der Stärke von `ADR-0112` Teilfrage 3 ab** („der umbenannte oder abgebildete Schlüssel behält die Position seiner Quellspalte") | konform als Abschwächung auf Spec-Ebene, aber siehe V-3: die Plan-Zeile „jede berührt keine ADR-Aussage" trifft für diese Zeile nicht zu |

## 4. Plan-vs-Diff und §3.13-Suchlauf — gegen beide Stände selbst nachgefahren

**Plan-Zeilen (§3) im Diff vertreten:** Pflichtenheft §1 (`LH-FA-CFG-007.a`, `LH-FA-CAP-009.a`), §2 (`SPEC-019`, `SPEC-002`, `SPEC-030` neu), §4 (Zeile `schema` und Absatz), §7; Architektur §1 (`ARC-001`-Zeile), §2/§4 (Form der Row Images, Antragsart-Tabelle, Backfill-Sequenz), §5. **Ungeplant:** nichts (die fünf Dateien des `--stat` decken sich mit den Plan-Zeilen plus Report). **Nicht erfolgt:** nichts aus „Ausdrücklich NICHT" berührt (Lastenheft, Routing, Handbuch, Code, Schema — `git diff --stat`).

| Plan-Zeile (Suchlauf) | Meine Messung | Ergebnis |
|---|---|---|
| Anker „Transformationsform offen" | `git grep -n 'Transformationsform' 1718546b -- spec/pflichtenheft.md` → `:268` „… Transformationsform offen"; HEAD → `:275` „… Transformationsform"; `git grep -n 'lh-fa-cfg-007a\|transformationsform-offen' 1718546b` → 0 | bestätigt |
| Antragsarten-Aufzählungen (`exclude_column`) | `git grep -c 'exclude_column' 1718546b -- spec docs/user harness README.md`: `benutzerhandbuch.md` 2, `harness/README.md` 1, `schema-rollout.md` 1, `architecture.md` 1, `pflichtenheft.md` 5 = **10 Zeilen in 5 Dateien** ✓. **HEAD: 12**, nicht 11 (Plan) | Parent bestätigt; Diff-Zahl überholt (V-2) |
| Träger außerhalb des Diffs (`benutzerhandbuch.md:265`, `:1524`; `schema-rollout.md:74`; `harness/README.md:141`) | `:265` ist ein SQL-Beispiel des Ausschlusses, `:1524` die `schema`-Zeile („nicht sicher interpretierbar"), `schema-rollout.md:74` nennt „die fünf SQL-Funktionen … CHECK … (fünf Antragsarten)", `harness/README.md:141` (`make test-integration`-Zeile) trägt `exclude_column` | Adressen bestätigt; Zuordnung (Handbuch → `slice-transformationen-betriebsdoku`, `schema-rollout.md` → `slice-transformationen-antragsweg-schema`, `harness/README.md:141` → E2E-Slices) stimmig |
| Zählwörter zu Antragsarten und Werten | `git grep -n -i 'fünf Arten\|sieben Arten\|fünf übrigen\|drei übrigen\|sechs übrigen\|fünf Werte\|sieben Werte' 1718546b -- spec`: 3 Treffer (`architecture.md:248` „fünf Arten", `pflichtenheft.md:535` „die drei übrigen", `:876` „(fünf Werte)"); Diff: `architecture.md:258` „sieben Arten", `:613` „die fünf übrigen", `:614` „die fünf übrigen", `:615` „die sechs übrigen", Historie `:1122` „fünf Werte" (Protokoll, bleibt), `:1133` „sieben Werte" | bestätigt; `git grep -n 'Antragsart' 1718546b -- spec` 15 Zeilen, HEAD 26 ✓ |
| Zeile `schema` in weiteren Trägern | Parent: `pflichtenheft.md:783` und `benutzerhandbuch.md:1524`; HEAD: `pflichtenheft.md:1020` und `benutzerhandbuch.md:1524` (2 Zeilen) | Zahlen bestätigt; der Plan nennt in dieser Zeile den Diff-Lokator **`:963`**, gemessen ist **`:1020`** (V-2) |
| `LH-FA-CFG-007` in Trägern der Abdeckung | `git grep -n 'CFG-007' 1718546b -- docs harness .d-check.yml` (ohne Records): 1 Zeile, `harness/README.md:134`. Diese Zeile lautet am **Parent** bereits „gemessen 2026-09-24 … 80 Anforderungen, **2 Waisen** — `LH-FA-CFG-007`/`LH-FA-CFG-008`"; `harness/README.md` liegt nicht im Diff (`git diff --stat`) | **widerlegt** — der Plan behauptet „3 Waisen — `LH-FA-CAP-009`/… vom 2026-09-23" und meldet die Drift an den Planner; die Zahl stimmt mit der heutigen Messung (2 Waisen) überein (V-1) |
| Architektur-Sicht `ARC-002`/`ARC-004`/`ARC-005` | `git grep -n 'ARC-002\|ARC-004\|ARC-005' 1718546b -- spec/architecture.md`; die drei Zeilen in §1 nennen „Konfiguration", „Fähigkeitsschnittstellen", „SQL-Funktionen/Views" — keine Aufzählung von Antragsarten oder Regeltypen; die Sicht trägt die Transformationen jetzt im Text (`ARC-004` als „ein Outbound Port (`ARC-004`)") | bestätigt |
| Fixrunden-Suchlauf `zeichengenau|Groß-/Kleinschreibung` | Parent `a32a1931` 1 Zeile, HEAD 10 Zeilen; `architecture.md` und `docs/user`: kein Treffer | bestätigt |
| Fixrunden-Suchlauf `byte-gleich|Position seiner|Position der Quellspalte` | `a32a1931`: `pflichtenheft.md:205`, `:299`, `:301`; HEAD: nur `:205` | bestätigt (Rest siehe V-4) |
| Fixrunden-Suchlauf `drei/fünf Formzeilen`, `Regelname ist ungültig` | `a32a1931`: `:681` „drei Formzeilen"; HEAD: `:683` „fünf Formzeilen", `:689` die neue Zeile | bestätigt |
| Fixrunden-Suchlauf `Abhilfe` | `a32a1931` `spec` 8 Zeilen, HEAD 11 (`:332`, `:335`, `:985`, `:1020`, `:1027`, `:1028`, `:1030`, `:1031`, drei Historie); `docs/user harness`: an beiden Ständen 3 | bestätigt (Führung: `LH-FA-CFG-007.a` `:335`) |
| Register-Zahlen §6/§8 (`ls …/evidence \| wc -l`) | `nachzug-laesst-ueberholten-text-stehen` **12**, `arbeit-ueberholt-stehenden-traeger` **32**, `dod-begruendung-unzutreffende-tatsachenbehauptung` **8**, `zitat-nennt-die-falsche-stelle` **8**, `zahl-in-traeger-driftet-gegen-die-messung` **23**, `adr-folgepflicht-ohne-traeger-slice` **1**; `state.md`: die ersten fünf „verkörpert", der letzte „offen" | bestätigt (F-10 behoben, Ursprung und Messlauf genannt) |

**Suchlauf-Lücke (gegen die eigene Zusage):** Der Suchlauf der Fixrunde suchte
die bewegten Formulierungen („Position seiner Quellspalte", „byte-gleich")
nur in `spec` und `docs/user`. Die Folge-Pläne unter
`docs/plan/planning/open/` tragen dieselben Formulierungen
(`grep -rn -E 'Position seiner|byte-gleich' docs/plan/planning/open/slice-transformationen-*.md`):
`slice-transformationen-kern-rename` (DoD-Zeile „an der Position seiner
Quellspalte", Determinismus-Test „byte-gleiches Image", Risiko
„Schlüsselposition nach der Umbenennung"), `slice-transformationen-map-value`
und `slice-transformationen-backfill-pfad` (je „byte-gleiches Image/Bild"). Die
Backfill-Welle hat ihre Folge-Pläne im Nachzug gezogen; hier steht die Meldung
aus (V-3).

## 5. Fixrunde — F-1 bis F-11 behoben? (Stellen selbst gelesen)

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 MEDIUM — Bezeichner-Vergleich | **behoben** | `SPEC-030` Absatz „Bezeichner" (`:917` ff.): `column`/`to`/Regelname zeichengenau, keine Faltung, kein Quoting; `column` nichtleer ohne U+0000 und gegen den Katalog; `to` nichtleer, ohne U+0000, ≤ 63 Byte UTF-8, sonst jedes Zeichen; `rule_name` 1–63 Zeichen `a`–`z`/`0`–`9`/`_`, Ungültiges endet `failed`, nicht gefaltet; die Adresse `schema.table.name` ist eindeutig zerlegbar (Schema und Tabelle tragen keinen Punkt); `SPEC-019` K3 „zeichengenau verschieden". Die Folge-Frage „welcher Test entscheidet `Name` gegen `name`?" ist beantwortet: sie sind verschieden |
| F-2 LOW — „Position"/„byte-gleich" | **behoben in `SPEC-030`/`LH-FA-CFG-007.a`; ein Rest in `LH-FA-CAP-009.a`** | „Reihenfolge und Inhalt" statt „Position"; Wirkort-Satz „denselben Inhalt der Row Images (Schlüsselmenge und Werte)". Rest: `:205` (V-4) |
| F-3 LOW — Fehlertext-Tabelle | **behoben** | fünf Formzeilen in der Reihenfolge Regelname → Form der `rule_spec` (NULL/kein Objekt/`kind` fehlt oder keine Zeichenkette) → unbekannter Regeltyp → unbekannter Schlüssel → Pflichtschlüssel/Werte; die Reihenfolge trennt „unbekannter `kind`" von „Pflichtschlüssel fehlt" (der `kind` wird vor den Pflichtschlüsseln des Typs bestimmt) |
| F-4 LOW — Zeile `schema` ohne Ursachen-Zuordnung | **behoben** | Aktions-Zelle: „Abhilfe nur bei nicht anwendbarer Regel: Regelstand ändern (Absatz unten)"; der Absatz verweist für Ursache, Erfassungspfad und Abhilfe und trägt nur den Run |
| F-5 LOW — kein Beispiel | **behoben** | vier Beispiele in `SPEC-030`: `rename_column` (Vorher-/Nachher-Image), `map_value` (`o` → `open`, `x` bleibt, `NULL` fehlt), abgelehnter Antrag (`to` = `id` → `Zielname kollidiert mit einer Spalte der Tabelle: public.orders.id`), nicht anwendbare Regel (Spalte `customer_name` nachträglich angelegt) |
| F-6 LOW — Umbrüche | **behoben** | `spec/architecture.md` und `spec/pflichtenheft.md`: die zwei benannten Absätze folgen dem Umbruch der Nachbarzeilen (Diff `8e977328`, nur Form) |
| F-7 INFO — `requested` vs. `pending` | **Gewinner deklariert** | Plan §3 „Bekanntes Wortlaut-Versehen der ADR"; Bestand nachgeprüft (§3) |
| F-8 INFO — Festlegungen jenseits des ADR-Textes | **im Plan §3 als Tabelle geführt** (mit Bestands-Beleg je Zeile) | siehe V-3 zur einen Zeile, die eine ADR-Aussage abschwächt |
| F-9 INFO — Führung je Sachverhalt | **behoben** | Plan §3 „Führende Stelle je Sachverhalt"; in der Spec: `SPEC-030` „Anwendbarkeit (die führende Stelle; … verweisen hierher)", `LH-FA-CFG-007.a` „Abhilfe (Zusage). Diese Stelle führt die Abhilfe …", `SPEC-008` verweist. Die Abhilfe trägt das Ersetzen der Regel und „ohne dass der Prozess erneut an derselben Regel endet" (das Kriterium „keine zweite Neustart-Schleife") |
| F-10 INFO — Register-Zahlen | **behoben** | §4 (letzte Zeile) |
| F-11 INFO — Träger-Adressen | **Adressen bestätigt, eine Aussage widerlegt** | siehe §4 und V-1: die Zahl „3 Waisen" steht in `harness/README.md:134` an keinem der beiden Stände |

## 6. Beobachtungen und Nachzüge an den Planner

**V-1 (LOW) — Plan-Aussage über `harness/README.md:134` ist falsch (§3.12 Instanz A).**
Der Suchlauf-Eintrag „`LH-FA-CFG-007` in Trägern der Abdeckung" meldet „3 Waisen
— `LH-FA-CAP-009`/… vom 2026-09-23" als überholte Zahl an den Planner. Gemessen:
die Zeile lautet am Parent `1718546b` und am `HEAD` „gemessen 2026-09-24 … 80
Anforderungen, **2 Waisen** — `LH-FA-CFG-007`/`LH-FA-CFG-008`" und deckt sich mit
`make doc-trace` (`80 Anforderung(en), 2 Waise(n)`). Es gibt keine Drift und
keine Meldung. **Nachzug:** die Plan-Zeile berichtigen (Parent-Befund „bereits
aktuell, keine Meldung"), den zweiten Teil des Review-Findings F-11 ebenso
lesen; die Fehlerart (eine Zahl aus einer Nachbardatei übernommen, ohne an
beiden Ständen nachzumessen) ist ein weiterer Beleg für
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, Kandidat für das
Beobachtungs-Register bei der Closure.

**V-2 (INFO) — zwei Plan-Zahlen sind nach der Fixrunde überholt.** Zeile
„Aufzählungen der Antragsarten": „Diff: 11 Zeilen" — gemessen am `HEAD` **12**
(die Fixrunde hat mit `spec/pflichtenheft.md:917` einen Satz zur Spaltenauswahl
im Absatz „Bezeichner" ergänzt; am Stand `71627d67` waren es 11). Zeile „Zeile
`schema` in weiteren Trägern": Diff-Lokator `:963`, gemessen **`:1020`** (die
Fixrunden-Tabelle nennt `:1020` bereits). **Nachzug:** beide Angaben auf den
Endstand ziehen oder als Stand `71627d67` kennzeichnen.

**V-3 (LOW) — Folge-Pläne tragen „Position" und „byte-gleich"; eine Festlegungs-Zeile schwächt eine ADR-Aussage ab.**
Der Suchlauf der Fixrunde deckt die Folge-Pläne nicht (§4). Die Zeile
„Schlüsselreihenfolge im Image nicht zugesagt" in Plan §3 schwächt die Aussage
„behält die Position seiner Quellspalte" aus `ADR-0112` Teilfrage 3 auf der
Spec-Ebene ab — der Vorspann der Tabelle behauptet „jede berührt keine
ADR-Aussage". **Nachzug:** (a) den Vorspann berichtigen („eine Zeile schwächt
auf Spec-Ebene ab; die Position bleibt Festlegung der ADR und darf der
Assembler-Test prüfen"); (b) `slice-transformationen-kern-rename` DoD-Zeile 1,
Determinismus-Test und Risiko „Schlüsselposition" mit einem Satz versehen, dass
Position und Byte-Gleichheit ADR-Festlegungen am Assembler-Ausgang sind, keine
Spec-Zusage am Lesepfad; (c) `slice-transformationen-map-value` und
`slice-transformationen-backfill-pfad` gleich. Die Backfill-Welle hat ihre
Folge-Pläne im Nachzug gezogen.

**V-4 (LOW) — verbleibende Zeile `spec/pflichtenheft.md:205`.**
`LH-FA-CAP-009.a` „Markierung": „Das Row Image ist byte-gleich dem WAL-Image
derselben Zeile: dieselbe Bild-Konstruktion, ausgeschlossene Spalten … fehlen,
die Transformationsregeln der Tabelle … wirken". Sie gilt dem Ausgang der einen
gemeinsamen Bild-Konstruktion (Backfill-Stand, unverändert) und ist über den
`jsonb`-Lesepfad nicht als Byte-Gleichheit beobachtbar — dort ist der Inhalt
(Schlüsselmenge und Werte) gleich, nicht die Reihenfolge, die `SPEC-030`
ausdrücklich nicht zusagt. Zwei Zusagen mit verschiedener Stärke stehen in einer
Datei. **Nachzug (Empfehlung):** auf „ist inhaltsgleich (Schlüsselmenge und
Werte) dem WAL-Image derselben Zeile: dieselbe Bild-Konstruktion" umformulieren;
`slice-transformationen-backfill-pfad` (dort „byte-gleiches Bild wie die
WAL-Change") im selben Zug lesen. Die Umformulierung ist kein Blocker dieses
Slice: die Zeile stammt nicht aus dem Diff, und ihre Stärke ist am Assembler
belegbar, nicht am Lesepfad.

**V-5 (INFO) — Abhilfe im Run und „Neustart".** `SPEC-008` (Absatz „Nicht
anwendbare Regel") sagt „Die Abhilfe ist die Regelstand-Änderung der
Abhilfe-Zusage"; die Abhilfe-Zusage in `LH-FA-CFG-007.a` enthält den
Prozess-Neustart des Erfassungspfads. `ADR-0117` Festlegung 4 nennt für den Run
keinen Neustart („erwartet, zu belegen im E2E-Beleg des Backfill-Pfads"). Die
Spec schließt den Neustart im Run nicht aus und fordert ihn nicht; ein Leser
kann ihn als Pflicht lesen. **Nachzug (optional):** ein Halbsatz „ein
Prozessneustart ist für den neuen Run nicht Teil der Zusage".

**V-6 (INFO) — `LH-FA-CFG-007` und `LH-FA-CFG-008` bleiben Waisen.**
Erwartet; `docs/user/e2e-abdeckung.md` trägt die Zeile erst mit
`slice-transformationen-e2e-wirkung`, ein umsetzender Slice führt die Kennung im
Träger. Kein Handlungsbedarf in diesem Slice.

**Konjunktiv-/Chronik-Prüfung:** Regex-Lauf über die `+`-Zeilen des Spec-Diffs
nach `hätte|wäre|würde|bliebe|stünde|sonst` und nach Vorher-/Nachher-Wörtern:
keine Konjunktiv-Chronik; Zusagen stehen im Indikativ der Norm mit der
Kennzeichnung „Zusagen an die Umsetzung, keine Messergebnisse"; die Sicht trägt
keinen Zukunftsmarker (Greenfield-Stil, [`AGENTS.md`](../../AGENTS.md) §3.4).

## 7. Anschlussfähigkeit für `slice-transformationen-kern-rename`

**Ja, mit vier benannten Lücken (Nachzug an den Planner, keine Blockade).**
Ein Implementer kann aus `SPEC-030`, `SPEC-019` und `LH-FA-CFG-007.a` ohne
Rückfrage umsetzen:

- **Parser der `rule_spec`:** JSON-Objekt, Pflichtschlüssel `kind`
  (Zeichenkette), je Regeltyp die genauen Schlüssel und Typen, fünf Formzeilen
  mit Klartext, Adresse und Prüfreihenfolge (`SPEC-019`); Bezeichner-Form je
  Feld (`SPEC-030`).
- **K1–K4:** je mit Bedingung, Fehlertext und Adresse; K3 zeichengenau und
  einschließlich ausgeschlossener Spalten und der Quellspalte selbst; K4 über
  den bestehenden Katalog-Abfrageweg.
- **Wirkung:** Ausschluss zuerst, dann Spaltenregeln; `rename_column` und
  `map_value` mit Abwesenheits-Vereinbarung (`NULL`, TOAST, ausgeschlossen,
  fehlendes Bild); Identität und übrige Felder unverändert; Beispiele mit
  Vorher-/Nachher-Image.
- **Anwendbarkeit:** Spalte in der Relation der Change, Zielname ohne Kollision
  mit einer Spalte der Relation; hängt an Regel- und Spaltenmenge, nie am Wert;
  Klasse `schema`, keine Persistierung, keine Bestätigung.

Lücken:

1. **Zuordnung der Form-Prüfungen zu Schichten (LOW).** Die Spec legt fest,
   *was* geprüft wird (Alphabet, 63 Byte, U+0000, nichtleer), nicht, ob die
   Domänen-Konstruktor-Invarianten (`kern-rename`) oder der Use Case
   (`antragsweg-usecase`) sie tragen. Der Plan von `kern-rename` nennt als
   Domänen-Invarianten nur „Spalten- und Zielname nicht leer, verschieden".
   **Nachzug:** benennen, welche der Form-Prüfungen `kern-rename` trägt (U+0000,
   63 Byte, Alphabet des Regelnamens).
2. **Sentinel-Abbildung für „Zielname gleich Quellname" (LOW).** Die Spec
   ordnet den Fall K3 zu („`Zielname kollidiert mit einer Spalte der Tabelle`",
   nach den Formzeilen). Verletzte die Domäne die Invariante „Zielname
   verschieden von Spaltenname" bereits im Konstruktor und der Use Case bildete
   sie auf „`rule_spec ist ungültig`" ab, träfe der Fehlertext nicht die Spec.
   **Nachzug:** im Plan festhalten, dass ein solcher Sentinel im Use Case auf
   den K3-Text abgebildet wird (Prüfreihenfolge: nach den Formzeilen).
3. **Gestaffelte Lieferung der Regeltypen (INFO).** `kern-rename` liefert nur
   `rename_column`; die Spec führt den geschlossenen Satz aus zwei Typen.
   Ein Zwischenstand, der `map_value` als „unbekannter Regeltyp" ablehnt,
   verletzt `SPEC-019` Zeile 3, ist aber erst über den Antragsweg erreichbar
   (`kern-rename` hat keinen). **Nachzug:** im Plan von `antragsweg-usecase`
   festhalten, dass die Regeltyp-Menge dort bereits der Spec entspricht oder der
   Antragsweg erst nach `map-value` freigeschaltet wird.
4. **Signatur der SQL-Funktionen (INFO).** `cdc.set_transformation(...)`/
   `cdc.remove_transformation(...)` stehen ohne Parameterliste (wie die übrigen
   Funktionen). Kein Problem für `kern-rename`; `antragsweg-schema` legt sie fest.

## 8. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, Zeile 6 mit den Nachzügen V-1/V-2 |
| DoD §2 — `[ ]`-Zeilen (4) | **korrekt offen** (Closure-/Rollen-Sequenz) |
| Konformität gegen `ADR-0112`, `ADR-0117`, `ADR-0059`, `LH-FA-CFG-007` | **konform**; eine Festlegung (Schlüsselreihenfolge) schwächt `ADR-0112` Teilfrage 3 auf Spec-Ebene ab (V-3) |
| Plan-vs-Diff | **deckungsgleich**, keine ungeplante Datei |
| Harte Regeln (§3.3, §3.4, §3.7, §3.11, §3.12, §3.13) | **erfüllt**; §3.13-Suchlauf deckt die Folge-Pläne nicht (V-3) |
| Review-Findings F-1…F-11 | **behoben** (F-2 mit Rest V-4) |
| Sensoren | `make gates` 0, `make docs-check` 0, `commit-traceability`/`doc-commits`/`doc-immutable` 0, `make doc-trace` 2 Waisen (erwartet), Mutation Exit 2 `id-unlinked` |

## Verdikt

**Bestätigt.** Die DoD-Behauptungen des Implementers und die Spec-Zusagen sind
durch eigene Messung belegt; die Spec trägt die Festlegungen von `ADR-0112` und
`ADR-0117` ohne ADR-/Slice-/Wellen-Bezug, der Bezeichner-Vergleich entspricht dem
Bestand. Kein Befund blockiert. V-1 bis V-5 und die Lücken 1–2 aus §7 gehen als
Nachzug an den Planner; sie ändern keinen Spec-Satz dieses Slice, außer der
optionalen Umformulierung von `spec/pflichtenheft.md:205` (V-4).

**Übergabe:** Verifier → Planner. Offen bleibt die Closure (Notiz mit
Lerneintrag, Beobachtungs-Register, §6-Ausgänge, Welle-Paarungen); danach der
reine `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
