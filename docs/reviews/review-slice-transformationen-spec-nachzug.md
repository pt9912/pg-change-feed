# Review-Report: slice-transformationen-spec-nachzug — 2026-09-26

**Review-Art:** Plan/Design — der Diff ist reine Spec-Doku (Pflichtenheft und
Architektur-Sicht) plus Slice-Plan; geprüft gegen Plan, ADRs und `AGENTS.md`
Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `1718546b..6ad71a49`, sechs Commits `e758fc32` …
`6ad71a49`; drei Dateien: `spec/pflichtenheft.md`, `spec/architecture.md`,
Slice-Plan `slice-transformationen-spec-nachzug` (Lifecycle `open` → `next` →
`in-progress`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Träger-Nachzug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-transformationen-spec-nachzug` (§1 Ziel, §2 DoD, §3 Plan
  und Suchlauf, §6 Risiken, §8) und Welle `welle-transformationen`
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  (Accepted) — Teilfragen 1–8, K1–K4, Folgepflichten 1, 5 und 7
- [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) (Accepted)
  — Festlegungen 1–5 (Klasse `schema` im Run, Regelstand-Wechsel im Lauf);
  [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md)
  (Schema-Version-Referenz, gelesen)
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  und [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  (Antrags-Queue, Statusvokabular, dauerhafter Stand)
- Lastenheft [`LH-FA-CFG-007`](../../spec/lastenheft.md) (Version 0.13.0),
  [`LH-FA-CFG-008`](../../spec/lastenheft.md),
  [`LH-FA-CFG-005`](../../spec/lastenheft.md),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md) und die Pflichtenheft-Stellen
  [`LH-FA-CFG-007.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md),
  [`SPEC-002`](../../spec/pflichtenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md),
  [`SPEC-019`](../../spec/pflichtenheft.md),
  [`SPEC-030`](../../spec/pflichtenheft.md),
  [`ARC-001`](../../spec/architecture.md)
- `AGENTS.md` (Hard Rules §3.3, §3.4, §3.5, §3.7, §3.9, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md` der vendored Baseline
  (Repo-Form nach den bestehenden `review-slice-*`-Reports)

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen):

- Gate-Läufe, Exit-Code jeweils ungefiltert in eine Datei gesichert und
  gesondert gelesen (`AGENTS.md` §3.9): `make docs-check` → Exit 0,
  `d-check: 1192 Datei(en) geprüft, 0 Befund(e)`; `make gates` → Exit 0 mit
  `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%`,
  `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`,
  `generated-sync: OK`, `gesamt: 0 Befund(e)`; `make doc-trace` → Exit 0,
  `80 Anforderung(en), 2 Waise(n)`.
- Mutationen auf Spec-Ebene (Eingabeseite der Zusage „die Spec-Straten tragen
  keinen ADR-Bezug und keine nackte Kennung"), je danach `git checkout`,
  `git status` sauber: (1) ein Link auf `ADR-0112` ans Ende von
  `spec/pflichtenheft.md` → `make docs-check` Exit 2, Befund `matrix-forbidden`
  („Referenz spec → adr ist nicht erlaubt"); (2) ein nacktes `LH-FA-CFG-007`
  am Ende von `spec/architecture.md` → Exit 2, Befund `id-unlinked`. Beide
  Zusagen sind an ihrer Eingabeseite rot färbbar.
- Suchlauf über den Diff: `git diff 1718546b..HEAD -- spec | grep '^+' | grep -c
  'ADR-0\|slice-\|welle-'` → `0` (`AGENTS.md` §3.4); `git diff
  1718546b..HEAD -- docs/user | wc -l` → `0` (Handbuch unberührt).
- Kennung: `git grep -h -o 'SPEC-0[0-9][0-9]' 1718546b -- spec/pflichtenheft.md
  | sort -u | tail -1` → `SPEC-029`, am Kopf `SPEC-030`; die neue Kennung ist
  die nächste freie, kein anderer Träger außerhalb der Pläne vergibt sie.
- Suchlauf-Feld (Slice-Plan §3) an **beiden** Ständen nachgefahren (`git grep`
  an `1718546b` und `HEAD`): (1) `transformationsform` — Parent: die
  Überschrift `spec/pflichtenheft.md:268`, Anker-Suche `lh-fa-cfg-007a|
  transformationsform-offen` 0 Treffer; Kopf: Überschrift Zeile 275, ein
  Treffer der Anker-Suche in der Plan-Zelle selbst; (2) `exclude_column` in
  `spec docs/user harness README.md` — 10 Zeilen in 5 Dateien am Parent
  (`architecture` 1, `pflichtenheft` 5, Handbuch 2, `harness/README.md` 1,
  `schema-rollout.md` 1), 11 am Kopf; (3) `Antragsart` in `spec` — 15 → 26
  Zeilen (Architektur 5 → 7, Pflichtenheft 10 → 19); Zählwörter
  „fünf Arten"/„drei übrigen"/„(fünf Werte)" am Parent an den genannten
  Stellen, „zehn Feld" nur `SPEC-021`/`SPEC-024` an beiden Ständen;
  (4) `nicht sicher interpretierbar` — 2 Zeilen an beiden Ständen
  (`pflichtenheft` :783 → :963, Handbuch :1524 unverändert); (5) `CFG-007` in
  `docs harness .d-check.yml` ohne Records — 1 Zeile (`harness/README.md:134`)
  an beiden Ständen. Alle Zahlen des Feldes stimmen; die Adressen der
  gemeldeten Träger sind unten unter den Negativbefunden bewertet.
- Zählwörter im Diff-Stand gelesen: `request_kind` sieben Werte
  (`spec/pflichtenheft.md:616`); `column_name` „fünf übrigen" (7 − 2 = 5,
  Aufzählung der fünf Namen stimmt), `rule_name` „fünf übrigen", `rule_spec`
  „sechs übrigen" (7 − 1 = 6); Architektur „sieben Arten" bei sieben
  Tabellenzeilen; `grep -n -i 'fünf\|sieben\|sechs' spec/*.md` findet kein
  weiteres Zählwort zu Antragsarten.
- Statusvokabular gegen den Bestand: `tools/schema/nacharbeit-administration.sql`
  (Zeilen 77–153: `'pending'`), `internal/adapters/driven/postgresstorage/queries/queries.go`
  (Zeilen 329, 358, 366: `status = 'pending'`) und
  [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  (Status `pending`/`applied`/`failed`) → der Wert lautet `pending`.
- Fehlertext-Form gegen den Bestand: `internal/application/port/inbound/verwaltung.go:26`
  (Sentinel „Spalte existiert nicht an der Quelle") und
  `internal/application/usecase/includecolumn/service.go:44`
  (`"%w: %s.%s.%s"`) → die Regel „Klartext, Doppelpunkt, Leerzeichen,
  Adresse `schema.table.column`" der neuen Tabelle entspricht dem Bestand.
- Text-Abgleich `LH-FA-CFG-007.a`, `SPEC-019`, `SPEC-030`, `SPEC-008`,
  `LH-FA-CAP-009.a` und Architektur mit `ADR-0112` (Teilfragen 1–8, K1–K4,
  Folgepflicht 1, 5, 7) und `ADR-0117` (Festlegungen 1–5) Satz für Satz; die
  Abweichungen und Ergänzungen stehen in F-7 und F-8.
- Register-Zahlen des Plans §6/§8 gegen `docs/plan/planning/observations/BEO-PGC/`
  (`ls <eintrag>/evidence | wc -l`, `state.md`) gemessen, siehe F-10.
- Lifecycle und Commit-Struktur: `git show --stat -M` je Commit —
  `e758fc32` und `ac8f757c` sind reine Renames (0 Zeilen), Inhalt liegt in
  getrennten Commits (`AGENTS.md` §3.3); Betreffs ohne `SPEC-`/`ARC-`-Kennung,
  jede Message mit `LH-FA-CFG-007` und `ADR-0112`.
- Dangling-Docker-Volumes vor dem Lauf 34, nach dem Lauf 34 (`docker volume ls
  -qf dangling=true | wc -l`); kein `prune`.

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Bezeichner-Vergleich und zulässige Zeichenmenge von `column`, `to` und `rule_name` sind nicht festgelegt

- `kategorie`: MEDIUM
- `quelle`: Maintainability (unklare Fehlerbehandlung am Rand des Spec-Bereichs);
  Bindung an den Folge-Slice `slice-transformationen-kern-rename` (Domäne/
  `Assembler`)
- `pfad`: `spec/pflichtenheft.md:892-899` (Tabelle in
  [`SPEC-030`](../../spec/pflichtenheft.md): „nicht leere Zeichenketten"),
  `spec/pflichtenheft.md:670-672` (K3 „verschieden von jedem Spaltennamen"),
  `spec/pflichtenheft.md:907-909` („Vergleich" — nur für Werte)
- `befund`: Nur der Wertvergleich von `map_value` ist als zeichengenau bestimmt.
  Ob `column` und `to` gegen die Spaltennamen der Relation exakt, ohne
  Groß-/Kleinschreibungs-Faltung und ohne Quoting-Interpretation verglichen
  werden, ob `to` beliebige Zeichen (Leerzeichen, Punkt, Steuerzeichen) tragen
  darf und ob K3 „verschieden" exakt oder gefaltet meint, steht nirgends; die
  Adresse `schema.table.rule_name` ist bei einem Punkt im Namen nicht mehr
  eindeutig zerlegbar. Der Folge-Slice implementiert Parser und K3-Prüfung
  gegen diesen Wortlaut.
- `verifizierbar`: nein — kein Gate liest die Lücke; Falsifikation ist das
  Lesen von `SPEC-030` mit der Frage „welcher Test entscheidet `Name` gegen
  `name`?"
- `klasse`: Randfall der Regelform nicht festgelegt

### F-2 — „Position" und „byte-gleich" sind über die gespeicherten und gelesenen Wege nicht beobachtbar

- `kategorie`: LOW
- `quelle`: Maintainability (Zusage ohne beobachtbaren Träger); `AGENTS.md` §3.12
  Instanz B (was die Umsetzung liefern muss, muss belegbar sein)
- `pfad`: `spec/pflichtenheft.md:297-301` (Auswertungsreihenfolge),
  `spec/pflichtenheft.md:910-911` (`SPEC-030` „Position"),
  `spec/pflichtenheft.md:437-442` (`SPEC-002`, `old_data`/`new_data` als
  `jsonb`)
- `befund`: Die Zusagen „der umbenannte Schlüssel behält die Position seiner
  Quellspalte" und „byte-gleiches Image" gelten nur für die Ausgabe des
  Zusammensetzens; `jsonb` ordnet Schlüssel nach Länge und Bytes, sodass
  `cdc.changes` und `GET /changes` die Position nicht tragen — der Wirkort-Satz
  „live wie beim erneuten Lesen dieselbe Form" ist nur für die Schlüsselmenge
  und die Werte wahr. Die Zusage ist an der Stelle prüfbar, an der sie
  steht, aber nicht über den E2E-Rundlauf.
- `verifizierbar`: nein — Falsifikation ist das Lesen gegen `SPEC-002`
  (`jsonb`); ein Unit-Test des `Assembler` deckt die Position, kein Lesezugriff
- `klasse`: Zusage ohne beobachtbaren Träger am Lesepfad

### F-3 — Fehlertext-Tabelle lässt drei Eingaben ohne Zeile

- `kategorie`: LOW
- `quelle`: Maintainability (unklare Fehlerbehandlung am Rand des Spec-Bereichs);
  der Slice-Plan §6 nennt den Wortlaut „bindend für
  `slice-transformationen-antragsweg-usecase`"
- `pfad`: `spec/pflichtenheft.md:655-657` (`rule_name` Pflicht — kein
  Fehlertext), `spec/pflichtenheft.md:685` (Formzeile), `spec/pflichtenheft.md:680-681`
  (Prüfreihenfolge)
- `befund`: (a) Ein `set_transformation` mit SQL-`NULL` oder JSON-`null` als
  `rule_spec` und ein Antrag mit leerem `rule_name` erscheinen in keiner Zeile
  (dieselbe Lage wie die Pflicht von `column_name`, dort ebenfalls ohne
  Fehlertext); (b) bei einem unbekannten `kind` sind die „Pflichtschlüssel des
  Typs" nicht bestimmbar, die Formzeile „ein Pflichtschlüssel fehlt" steht aber
  vor der Zeile „unbekannter Regeltyp" in der Prüfreihenfolge.
- `verifizierbar`: nein
- `klasse`: Fehlertext-Tabelle mit Lücke am Rand

### F-4 — `SPEC-008` Zeile `schema`: „Abhilfe: Regelstand ändern" steht unbedingt neben Ursachen ohne Regelbezug

- `kategorie`: LOW
- `quelle`: Maintainability; Prüfmuster „Nachzug widerspricht dem Nachbarn im
  selben Träger"
- `pfad`: `spec/pflichtenheft.md:963` (Zeile `schema` der Fehlerklassen-Tabelle)
- `befund`: Die Bedingungs-Zelle führt zwei Ursachen (nicht sicher
  interpretierbare Schemaänderung/Dekodierfehler; nicht anwendbare Regel); die
  Aktions-Zelle nennt „Abhilfe: Regelstand ändern" und die
  Erfassungspfad-/Run-Aussagen ohne Zuordnung zur Ursache. Für einen
  Dekodierfehler oder eine inkompatible Relation ist der Regelstand keine
  Abhilfe; der Absatz darunter (Zeilen 968-978) trennt die Ursache, die Zeile
  selbst nicht.
- `verifizierbar`: nein
- `klasse`: Zellinhalt ohne Ursachen-Zuordnung

### F-5 — `SPEC-030` trägt kein Beispiel

- `kategorie`: LOW
- `quelle`: Maintainability (Lesbarkeit des Vertrags für den Folge-Slice)
- `pfad`: `spec/pflichtenheft.md:885-930` (`SPEC-030`)
- `befund`: Die Regelform ist als Tabelle mit Schlüsseln und Formen beschrieben;
  weder ein vollständiges `rule_spec`-Objekt je Regeltyp noch ein Vorher-/
  Nachher-Image noch ein abgelehnter Antrag ist gezeigt. Randfälle (leeres
  `values`, Abbildung auf sich selbst, `NULL`-Wert) stehen nur als Prosa.
- `verifizierbar`: nein
- `klasse`: Datenform ohne Beispiel

### F-6 — zwei nach einer Teilersetzung unbrochene Zeilen

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `spec/architecture.md:402` (134 Zeichen: der neue Satz „… gelesen
  wird; der Erfassungspfad bleibt davon unberührt." läuft in „Der
  Snapshot-Leser trägt den Snapshot-Export, den Import, die Lesesperre und"
  weiter), `spec/pflichtenheft.md:208-209` (Zeile 209 länger als der Umbruch der
  Nachbarzeilen)
- `befund`: Beide Absätze folgen sonst dem Umbruch von etwa 80 Zeichen; die
  Ersetzung hat je einen Satz ohne Umbruch stehen lassen.
- `verifizierbar`: nein — `structure` prüft Zellenlängen, keine Zeilenlängen
- `klasse`: Umbruch nach Teilersetzung nicht nachgezogen

### F-7 — `ADR-0112` Folgepflicht 5 (b) nennt `requested`, die Spec `pending`

- `kategorie`: INFO
- `quelle`: Zwei-Quellen-Prüfung mit deklariertem Gewinner (Source Precedence:
  Pflichtenheft Rang 2 vor ADR Rang 4)
- `pfad`: `spec/pflichtenheft.md:338` und `spec/pflichtenheft.md:618` gegen
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 5 (b)
- `befund`: Der Bestand kennt nur `pending`/`applied`/`failed`
  (`nacharbeit-administration.sql`, `queries.go`, `ADR-0050`); die Spec folgt
  ihm, das Wort `requested` der ADR ist ein Wortlaut-Versehen ohne
  Entsprechung im Schema. Die Spec ist richtig; der Wortlaut der `Accepted` ADR
  bleibt (`AGENTS.md` §3.5). Der Folge-Slice `slice-transformationen-e2e-abhilfe`
  übersetzt das Wort bereits selbst (Plan Zeile 80: „Zustand `requested` des
  ADR-Wortlauts ist der Status `pending`"). Eine Berichtigungs-ADR ist dafür
  nicht nötig.
- `verifizierbar`: ja — `grep -n "'pending'" tools/schema/nacharbeit-administration.sql`
- `klasse`: Wortlaut-Drift ADR gegen Bestand (Gewinner deklariert)

### F-8 — Festlegungen der Spec über den Text von `ADR-0112` hinaus (Bewertung)

- `kategorie`: INFO
- `quelle`: Spec-Stratum-Regel „präzisieren ja, erweitern nie" (Kopf des
  Pflichtenhefts); [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfragen 2 und 3
- `pfad`: `spec/pflichtenheft.md:912-918` (Randfälle), `spec/pflichtenheft.md:613`
  (`column_name`), `spec/pflichtenheft.md:677-693` (Fehlertexte)
- `befund`: Von der Spec ergänzt und in der ADR nicht bestimmt: leeres `values`
  endet `failed` (strenger als die ADR, die nur unbekannte Schlüssel/`kind` und
  Nicht-Zeichenketten nennt); Zielname gleich Quellname endet `failed`
  (folgt wörtlich aus K3: die Quellspalte ist ein Spaltenname der
  Quelltabelle); Abbildung eines Werts auf sich selbst zulässig; mehrere
  Quellwerte auf denselben Zielwert zulässig (die ADR erwartet den Fall
  ausdrücklich in ihren Konsequenzen); `column_name` bleibt für beide
  Transformations-Antragsarten `NULL` (die Spalte der Regel steht in
  `rule_spec`); Fehlertext-Wortlaute und Prüfreihenfolge. Keine dieser
  Festlegungen widerspricht einer Aussage der ADR oder von `ADR-0117`; keine
  führt eine neue bindende Anforderung ein, sie präzisieren `LH-FA-CFG-007.a`.
  Eine Berichtigungs-ADR ist nicht nötig; die Wortlaute sind bindend für den
  Folge-Slice und werden bei Abweichung dort als Plan-Nachzug hierher
  zurückgeführt (Plan §6).
- `verifizierbar`: nein
- `klasse`: Festlegung jenseits des ADR-Textes (verträglich)

### F-9 — Anwendbarkeit, Reihenfolge und Abhilfe stehen je an zwei Stellen ohne benannten Gewinner

- `kategorie`: INFO
- `quelle`: Prüfmuster „Zwei-Quellen-Drift" (latent, heute gleichlautend)
- `pfad`: `spec/pflichtenheft.md:297-301` und `spec/pflichtenheft.md:909-911`
  (Reihenfolge/Position); `spec/pflichtenheft.md:328-336` und
  `spec/pflichtenheft.md:919-926` (Anwendbarkeit); `spec/pflichtenheft.md:337-343`
  und `spec/pflichtenheft.md:968-978` (Abhilfe)
- `befund`: Dieselbe Definition steht in `LH-FA-CFG-007.a` und in `SPEC-030`
  bzw. `SPEC-008`. Die Texte sind gleichwertig; die Abhilfe in
  `LH-FA-CFG-007.a` nennt nur `cdc.remove_transformation`, `SPEC-008` zusätzlich
  das Ersetzen der Regel, und das Kriterium „keine zweite Neustart-Schleife"
  aus `ADR-0112` Folgepflicht 5 (d) steht in keiner der beiden Stellen. Ein
  Gewinner ist nicht deklariert; die Backfill-Vorlage
  (`LH-FA-CAP-009.a`/`SPEC-029`) führt es ebenso.
- `verifizierbar`: nein
- `klasse`: Mehrfach geführte Definition ohne Gewinner-Zeiger

### F-10 — Register-Zahlen in Plan §6/§8 driften gegen das Register

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahl mit Ursprung; der Träger liegt
  außerhalb der geänderten Zeilen des Diffs — Planner-Text vom Eröffnungstag)
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-spec-nachzug.md`
  §8 (Zeilen 307-317) und §6 (Zeile 261)
- `befund`: Gemessen an `docs/plan/planning/observations/BEO-PGC/`:
  `arbeit-ueberholt-stehenden-traeger` 32 Evidence-Dateien (Plan: 26×),
  `zahl-in-traeger-driftet-gegen-die-messung` 23 (Plan: 13×),
  `dod-begruendung-unzutreffende-tatsachenbehauptung` 8 (Plan: 6×),
  `zitat-nennt-die-falsche-stelle` 8 (Plan: 7×),
  `nachzug-laesst-ueberholten-text-stehen` 12 Evidence-Dateien, Zustand
  „verkörpert" (Plan §6 und §8: „offen, 2×"); `adr-folgepflicht-ohne-traeger-slice`
  1, offen — stimmt. Die Sichtung nennt keinen Lauf und keinen Messzeitpunkt.
- `verifizierbar`: ja — `ls docs/plan/planning/observations/BEO-PGC/<eintrag>/evidence | wc -l`
- `klasse`: Zahl im Träger driftet gegen die Messung

### F-11 — Träger außerhalb des Diffs: gemeldete Adressen stimmen, ein Zahlwort ist bereits überholt

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug: Meldung statt stiller Mitänderung)
- `pfad`: `harness/README.md:134` („3 Waisen"), `harness/targets/schema-rollout.md:74`,
  `docs/user/benutzerhandbuch.md:265` und `:1524`, `harness/README.md:141`
- `befund`: Die Zuordnung des Slice stimmt: Handbuch §4/§6 → `slice-transformationen-betriebsdoku`
  (Plan Zeilen 135, 148), `schema-rollout.md:74` (fünf SQL-Funktionen, CHECK mit
  fünf Antragsarten) → `slice-transformationen-antragsweg-schema` (Suchlauf-Zeile
  „Aufzählungen der Antragsarten", `grep … harness`), Waisen-Aussage und
  E2E-Aufzählung → `slice-transformationen-e2e-wirkung`/`-e2e-abhilfe` (Plan
  Zeilen 147/157). Die Zahl „3 Waisen" in `harness/README.md:134` trägt die
  Datierung 2026-09-23 und ist gegen die heutige Messung (`make doc-trace`:
  2 Waisen, `LH-FA-CFG-007`/`LH-FA-CFG-008`) überholt, unabhängig von diesem
  Diff — `LH-FA-CAP-009` ist zwischenzeitlich belegt.
- `verifizierbar`: ja — `make doc-trace`
- `klasse`: Träger außerhalb des Diffs driftet (gemeldet, adressiert)

---

## Negativbefunde

- geprüft, ohne Befund: `spec/architecture.md` — kein ADR-/Slice-/Wellen-Bezug
  im Diff (`AGENTS.md` §3.4, Suchlauf 0); Antragsart-Tabelle mit sieben Zeilen,
  „sieben Arten"; `ARC-001`-Zeile trägt „Transformationsregel"; Aussage „Form der
  Row Images" (erst Ausschluss, dann Regeln, vor Serialisierung und
  Persistierung) stimmt mit `ADR-0112` Teilfrage 6 überein; Fail-closed-Zeile,
  Absatz zur Ausführung und §5-Zeile („die Regel wird nie übersprungen, die
  Change nie roh ausgeliefert") stimmen mit `ADR-0112` Teilfrage 4 und `ADR-0117`
  Festlegung 2/3/5 überein; „Letzte Änderung" gezogen.
- geprüft, ohne Befund: `LH-FA-CFG-007.a` — Überschrift ohne „offen"; die
  Einleitung nennt die Sätze „Zusagen an die Umsetzung, keine Messergebnisse"
  (dieselbe Form wie `LH-FA-CAP-009.a`); die Abhilfe steht ausdrücklich als
  „Zusage" und „erst mit dem Beleg am laufenden System eine Tatsache"
  (`ADR-0112`: „erwartet, nicht am Code belegt"); kein zweiter Konfigurationsweg
  (Optionen B/C der ADR verworfen); Reihenfolge, K1–K4, Ausschluss-zuerst,
  Wirkort, Rohform-Konsequenz (nicht rückwirkend, nicht rekonstruierbar,
  `map_value` nicht umkehrbar bei Mehrfachabbildung), Nachrichtenschemata
  unverändert — alle mit der ADR gleichlautend. `LH-FA-CFG-008.a` bleibt offen
  und ist abgegrenzt.
- geprüft, ohne Befund: Lastenheft-Bindung — Happy Path, Boundary („Auflösung
  bei Mehrdeutigkeit definiert": statischer Ausschluss, expliziter Fehler) und
  Negative („sichtbarer Fehlerzustand, nicht still unveränderte Auslieferung":
  Klasse `schema`, keine Persistierung) von
  [`LH-FA-CFG-007`](../../spec/lastenheft.md) haben je eine Entsprechung; das
  Pflichtenheft erweitert nicht (Spec-Stratum-Regel): `SPEC-030` ist eine
  Struktur-Kennung, keine neue `LH-*`-Anforderung.
- geprüft, ohne Befund: [`SPEC-019`](../../spec/pflichtenheft.md) — Spalten
  `rule_name`/`rule_spec` mit Pflicht je Antragsart, `request_kind` mit sieben
  Werten, K1–K4 gleichlautend mit `ADR-0112` Teilfrage 3 (K3 ergänzt „Quellspalte
  selbst" als wörtliche Folge), zweite Bedeutung von `applied` und Ordnung
  `requested_at`/`administration_request_id` gleichlautend mit `ADR-0065`, Antrag
  gegen eine Tabelle ohne Bindung `applied`, Grants unverändert; die
  Fehlertext-Form entspricht dem Bestand (`ErrSourceColumnMissing`).
- geprüft, ohne Befund: [`SPEC-002`](../../spec/pflichtenheft.md) und
  `LH-FA-CAP-009.a` — Regelwirkung im Backfill-Bild, Regelstand in der
  Fail-closed-Prüfung mit Klasse `configuration` (`ADR-0117` Festlegung 5),
  Run-Fehler der Klasse `schema` vor der ersten Zeile und run-lokal
  (`ADR-0117` Festlegung 2/3), neuer Antrag nach Abhilfe (Festlegung 4).
- geprüft, ohne Befund: `LH-QA-SEC-004` — Ausschluss zuerst, weder Quell- noch
  Zielname noch Wert im Image, K3 hält Zielnamen von ausgeschlossenen Spalten
  fern, `cdc.exclude_column`/`cdc.include_column` gegen eine Spalte mit Regel
  zulässig: alle drei Sätze stehen in `LH-FA-CFG-007.a` und tragen die
  `ADR-0112`-Aussage.
- geprüft, ohne Befund: Zukunftsform (`AGENTS.md` §3.7/§3.12) — jede neue
  Aussage einzeln gelesen: die Sätze in `LH-FA-CFG-007.a` stehen unter dem
  Vorbehalt „Zusagen an die Umsetzung"; `SPEC-019`/`SPEC-030`/`SPEC-002` sind
  Festlegungen der Datenform (dieselbe Form wie `SPEC-029`), keine Aussagen über
  einen gemessenen Ist-Zustand; `SPEC-008` verweist für die Abhilfe im
  Erfassungspfad auf die Zusage in `LH-FA-CFG-007.a` und wiederholt die
  ordnungsabhängige Abfolge nicht als Tatsache; die Abhilfe im Run nennt keinen
  Prozessneustart als überflüssig (`ADR-0117` Festlegung 4 „erwartet").
- geprüft, ohne Befund: Chronik/Forensik-Sprache in Spec-Prosa und
  Historienzeilen — die neuen Absätze und Zeilen nennen den Zustand
  (Fundstellen von „zuvor"/„neu": Sachbezug „die zuvor nicht bestätigte
  Transaktion", „erst entfernen, dann neu setzen", kein Vorher/Nachher der
  Spec); die vier Historienzeilen vom 2026-09-26 tragen keinen ADR-/Slice-Bezug
  im Text.
- geprüft, ohne Befund: Zählwörter und Aufzählungen im Diff-Stand —
  `request_kind` sieben, „sieben Arten", „fünf/sechs übrigen", „zehn Felder"
  unverändert (`SPEC-021`/`SPEC-024`), keine Nennung von „dreizehn" im
  Pflichtenheft; kein überholtes Zählwort zu Antragsarten im selben Dokument
  (`grep` über beide Spec-Dateien).
- geprüft, ohne Befund: `make doc-trace` — `LH-FA-CFG-007` bleibt Waise (ADR-Spalte
  `ADR-0112, ADR-0117`, Slices und Coverage `—`), `LH-FA-CFG-008` bleibt Waise;
  beide erwartungsgemäß bis zu den Folge-Slices. Die Zahl im Slice-Plan (2
  Waisen) stimmt mit der eigenen Messung.
- geprüft, ohne Befund: Slice-Plan §3–§8 — Umfang deckt sich mit dem Diff
  (zwei Spec-Dateien, keine Lastenheft-, Handbuch-, Code- oder Schema-Änderung);
  die „Plan-Nachzug"-Zeilen in §3 sind als solche gekennzeichnet; §6-Ausgänge
  stehen unbesetzt („bei Closure"), §7 unbesetzt; die Übergänge `open` → `next`
  → `in-progress` sind reine Renames; das Suchlauf-Feld trägt je Zeile Fund und
  Nichtfund an beiden Ständen, die genannten Zahlen stimmen (siehe
  Prüfungsliste); WIP-Limit 1 erfüllt (`in-progress/` trägt Roadmap und diesen
  Slice).
- geprüft, ohne Befund: Handbuch-Regeln des Skills — der Diff berührt
  `docs/user/benutzerhandbuch.md` nicht (keine Versionshistorie-Pflicht); die
  neuen SQL-Funktionen erscheinen nur in der Spec, der Handbuch-Zug ist mit
  Adresse `slice-transformationen-betriebsdoku` aufgeschoben (kein Code, keine
  Betreiber-Oberfläche im Diff).
- geprüft, ohne Befund: `AGENTS.md` §3.5 (keine `Accepted` ADR im Diff), §3.6
  (keine Gate-Lockerung), §3.11 (keine host-lokalen Pfade in den Diff-Zeilen),
  §3.1 (kein Skript, kein Make-Target im Diff), Traceability/ID-Schema (jede
  Commit-Message mit `LH-FA-CFG-007`/`ADR-0112`, Betreffs ohne `SPEC-`/`ARC-`).
- Schwerpunkte laut Skill ohne Anwendungsfall im Diff (kein Code, kein
  Kommentar in Produktionscode, keine Workflow-Datei, kein Handbuch-Inhalt):
  Chronik in Code-Kommentar, Suppression, Sicherheits-/Korrektheitspfad im Code,
  Docker-only-Skripte — nicht berührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Randfall der Regelform nicht festgelegt ·
Zusage ohne beobachtbaren Träger am Lesepfad · Fehlertext-Tabelle mit Lücke am
Rand · Zellinhalt ohne Ursachen-Zuordnung · Datenform ohne Beispiel · Umbruch
nach Teilersetzung nicht nachgezogen · Wortlaut-Drift ADR gegen Bestand
(Gewinner deklariert) · Festlegung jenseits des ADR-Textes (verträglich) ·
Mehrfach geführte Definition ohne Gewinner-Zeiger · Zahl im Träger driftet gegen
die Messung · Träger außerhalb des Diffs driftet (gemeldet, adressiert)

## Verdikt

**Merge-blockierend:** nein — kein HIGH; die Konsistenz von Spec und
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)/
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) ist gegeben
(keine Berichtigungs-ADR nötig, F-7/F-8), die Doku-Gates sind grün, beide
Spec-Mutationen färben das Gate rot. F-1 (MEDIUM) ist ein Ein-Absatz-Vorbehalt
in `SPEC-030` und wird vor dem Start von `slice-transformationen-kern-rename`
geschlossen, weil dessen Parser- und K3-Tests gegen diesen Wortlaut laufen; die
Entscheidung über die Fixrunde liegt beim Planner. Die DoD-Zeile „Review
durchgeführt" bleibt deshalb offen und wird nach der Fixrunde regulär
nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug: Grenze
„braucht der Slice eine Fixrunde, bleibt die Checkbox offen").

**Übergabe:** Findings F-1 bis F-6 gehen an den Implementer (Fixrunde bzw. als
Plan-Nachzug an die Folge-Slices, F-2 an `slice-transformationen-kern-rename`
und `slice-transformationen-e2e-wirkung` als Hinweis auf den Beobachtungsort der
Position); F-7 bis F-9 gehen als Kenntnisnahme an Verifier/Architect; F-10 und
F-11 gehen an den Planner (Register-Sichtung bei Closure, Träger-Adressen). Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-
Konformität prüft der Verifier separat (Modul 11).
