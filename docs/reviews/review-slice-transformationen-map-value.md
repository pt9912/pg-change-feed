# Review-Report: slice-transformationen-map-value — 2026-09-26

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec-Stellen und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-transformationen-map-value` (Welle `welle-transformationen`), Diff-Range
`e5a11979..072d65b2` (9 Commits, 14 Dateien, +1455/−97 einschließlich der drei Lifecycle-Commits `ce7acc8c`,
`fa5c2e29`, `4c5000e9`; Baum sauber, nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen
ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

**Ablage:** Der Reviewer-Lauf durfte den Report nicht selbst schreiben; der Text ist sein Ergebnis, die Sitzung des
Auftraggebers hat ihn unverändert in diese Datei übertragen (Kennungen in Inline-Code gesetzt).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-transformationen-map-value` (§1, §2 DoD-Wortlaut als Bezug, §3 Plan mit Auflösung der Übergabe aus
  `slice-transformationen-backfill-pfad` und `slice-transformationen-antragsweg-usecase`, Suchlauf-Feld, §6)
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 2, 3, 5, 6,
  Folgepflicht 4 und 7), [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md),
  [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
  [`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md),
  [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) (unverändert gelesen),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`SPEC-030`](../../spec/pflichtenheft.md), [`SPEC-019`](../../spec/pflichtenheft.md) (K1–K4)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../spec/lastenheft.md),
  [`LH-FA-DAT-005`](../../spec/lastenheft.md)
- `AGENTS.md` (§3.1, §3.2, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md` (`MR-000` bis `MR-003`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-transformationen-backfill-pfad.md`

**Eigene Messungen:**

- **Diff-Stat nach DoD Punkt 2:**
  `git diff --stat e5a11979 -- internal/application/usecase internal/bootstrap tools/schema internal/adapters spec ':!*_test.go'`
  ist **leer**. Ohne den Ausschluss nennt der Stat genau sieben Testdateien: `administrationrequest_test.go`,
  `transformation_mapvalue_test.go` und `transformation_test.go` im `mapper`-Paket, `backfill/transformation_test.go`,
  `settransformation/service_test.go`, `administration_internal_test.go`, `backfill_image_parity_test.go`.
- **Kodierung der Zuordnung:** Ein Wegwerf-Test in der Kopie, nicht committet, zog 300 000 zufällige Zuordnungen mit
  bis zu vier Paaren, Feldern bis drei Zeichen aus einem Alphabet mit U+0000/1/2, `é` als ein Zeichen und als
  `e`+U+0301 und dem Byte 0xFF. Ergebnis: keine Kollision zweier verschiedener Zuordnungen auf dieselbe
  Zeichenkette; `Values()` liefert jede Zuordnung unverändert zurück; `lookupMappedValue` findet jeden Schlüssel und
  lehnt jede Sonde ab, die keiner ist. Die Längenpräfix-Kodierung ist damit für Felder beliebigen Inhalts eindeutig.
  Die Längen-Grenze steht in F-3.
- **Läufe** (Exit-Codes ungefiltert gesichert und einzeln gelesen): `make gates` Exit 0; `make test` Exit 0, 45 Pakete
  ok; `make test-store` Exit 0, der neue Store-Test läuft im Paket `postgresstorage` mit; `make fmt-check` 257 Dateien
  formatiert; `make kommentar-kennungen DIFF=e5a11979` kein Kandidat; `make suchlauf-nachmessen` auf dem Plan „22 Zeilen
  stimmen“.
- **Mutationen** liefen an einer Kopie im Scratchpad (`git archive HEAD`, je Mutation `sed … Datei > Kopie`, danach `cp`
  zurück; kein `sed -i`, kein Host-python). Dangling Volumes vorher 34, nachher 34, kein prune.

## Findings

### F-1 — Der Paritätstest sagt für `map_value` zu, was er nicht bindet

- `kategorie`: MEDIUM
- `quelle`: `ADR-0083` (Kommentar-Zusage; „Zusage ohne Bindung an ihre Eingabeseite“ des Skills); `LH-QA-SEC-004`
- `pfad`: `internal/bootstrap/backfill_image_parity_test.go:132-140` (Doc-Kommentar
  `TestBackfillAndWALImagesAreByteEqualWithRules`), `:76-83`, `:96-99`, `:231-238` (`parityRuled.silent`)
- `befund`: Der Kommentar sagt zu, die sichtbare Wirkung der Regel — bei `map_value` „der abgebildete Wert statt des
  Quellwerts“ — stehe im Bild, wo die Spalte einen Wert und keinen Ausschluss trägt, „und fehlt sonst“. Für `map_value`
  trägt `parityRuled` kein `silent` (`nil`), der Zweig „NULL oder ausgeschlossen“ prüft nichts und endet mit
  `continue`. Gemessen: Die Mutation `containsName(excluded, column) && !mapsColumn(rules, column)` in `BuildRowImage`
  hebt den Ausschluss nur für eine Spalte mit `map_value`-Regel auf. Sie lässt
  `TestBackfillAndWALImagesAreByteEqualWithRules` grün, färbt aber `TestExcludedColumnIsUnreachableForEveryRuleKind`
  (`mapper`), `TestExecuteRulesNeverLeakExcludedColumns` (Run) und `TestBuildRowImageMapValue` (Domäne) rot. Die
  Sicherheitsaussage ist an drei anderen Stellen gebunden, die Zusage dieses Kommentars für `map_value` nicht.
- `verifizierbar`: ja — Mutation wie beschrieben, nur `internal/bootstrap` bleibt grün
- `klasse`: Kommentar-Zusage ohne Bindung (Test-Kommentar)
- Einstufung: Der Skill führt „Zusage ohne Bindung an ihre Eingabeseite“ unter HIGH. Hier gilt die Klasse nur für den
  Satz im Kommentar dieses einen Tests. Die Eigenschaft selbst ist gemessen an drei Stellen rot färbbar, und der
  Paritätstest trägt sie für `rename_column`. MEDIUM statt HIGH ist ein Urteil des Reviewers und hier begründet.

### F-2 — DoD-Zeilen 1 und 2 tragen ihren Wortlaut unverändert, §3 legt sie aus

- `kategorie`: LOW
- `quelle`: Maintainability; Skill-Klasse „Nachzug widerspricht dem Nachbarn im selben Träger“
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-map-value.md:66-87` (DoD 1, 2) gegen `:124-145` (§3)
- `befund`: DoD 1 verlangt „bestehende Tests ohne geänderte Erwartung“, §3 nennt zwei geänderte Erwartungen
  (`TestTransformationKindsIsAClosedSet`, `TestTransformationSpecBuild`). DoD 2 nimmt nur Testdateien in
  `internal/adapters/**` aus, §3 legt die Ausnahme auf jedes Verzeichnis der Aufzählung aus. Die Übergabe aus
  `slice-transformationen-backfill-pfad` sagt, DoD 2 müsse die Ausnahme selbst nennen. Beide DoD-Zeilen stehen mit
  `[x]` und unverändertem Wortlaut. Der Store-Test `administrationrequest_test.go` (`adapters/**`) zählt keine
  Regeltypen auf und liegt damit auch außerhalb der Auslegung; §3 nennt ihn selbst „über den Plan hinaus, nur Test“.
- `verifizierbar`: nein (Lese-Handlung; die Messung selbst stimmt, der Stat ohne Testdateien ist leer)
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-3 — Der Kodierungstest sagt „jede Stufe“/„jeder Länge“, die höchste Stufe fehlt

- `kategorie`: LOW
- `quelle`: `ADR-0083` (Aussage breiter als Messung)
- `pfad`: `internal/domain/model/transformation_mapvalue_test.go:314-323` (Kommentar und Name
  `TestMapValueEncodingCarriesFieldsOfAnyLength`), `internal/domain/model/transformation.go:196-215`
- `befund`: Der Kommentar nennt Feldlängen bis 70 000 Byte „jede Stufe, an der ein weiteres Längen-Byte gebraucht
  wird“. Das vierte Längen-Byte (ab 16 777 216 Byte) wird nie geübt. Gemessen: Ein Präfix von drei statt vier Byte,
  Schreiben und Lesen zugleich geändert, lässt `model`, `mapper`, `backfill` und `bootstrap` grün. Die Längenangabe
  ist ein `uint32`. Eine Feldlängen-Grenze steht nirgends im Kommentar. Über die `jsonb`-Spalte (255 MiB je Wert) ist
  sie oberhalb 16 MiB erreichbar.
- `verifizierbar`: ja — Mutation wie beschrieben
- `klasse`: Test-Aussage breiter als die Messung

### F-4 — Doc-Kommentar der Nichtanwendbarkeit beschreibt nur den Fall mit Zielname

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger außerhalb des Diffs, Suchmuster ohne die bewegte Beschreibung)
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:58-64` (Doc-Kommentar `ErrTransformationNotApplicable`)
- `befund`: Der Kommentar nennt als Nichtanwendbarkeit „ihre Spalte fehlt in der Relation, oder ihr Zielname gleicht
  einer Spalte der Relation“. `map_value` trägt keinen Zielnamen, für ihn gilt nur der erste Fall. `CheckApplicable`
  im Diff trägt die Beschreibung bereits als „Zielname von `rename_column`“. Das committete Suchlauf-Feld sucht
  Regeltypnamen, `TransformationKinds`, Zählwörter und die Handbuch-Träger, nicht das Wort „Zielname“ außerhalb der
  Domäne. `git grep -n -i Zielname -- internal cmd ':!*_test.go'` außerhalb `internal/domain` findet neun Zeilen.
  Davon ist nur diese eine definierend; die übrigen nennen K3 und die Kollision zweier Regeln und bleiben für
  `rename_column` wahr.
- `verifizierbar`: ja — `git grep` wie beschrieben
- `klasse`: Träger-Nachzug (Suchmuster ohne die bewegte Beschreibung)

### F-5 — Streichen von `map_value` aus der Domänen-Menge färbt genau einen Test rot

- `kategorie`: INFO
- `quelle`: `ADR-0112` (Fitness Function); `ADR-0083`
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-map-value.md:88-94` (DoD 3)
- `befund`: DoD 3 nennt als Beleg die Mutation, die `map_value` aus der Domänen-Menge entfernt („die Tests färben sich
  rot“). Gemessen färbt `TransformationKinds()` ohne `map_value` allein `TestTransformationKindsIsAClosedSet`
  (`model`) rot. Die Eigenschaftstests in `mapper`, Run und `bootstrap` laufen dann über einen Typ und bleiben grün.
  Die Gegenrichtung trägt: Ein dritter Typ in der Menge färbt fünf Tests rot:
  `TestTransformationKindsIsAClosedSet`, `TestExcludedColumnIsUnreachableForEveryRuleKind`,
  `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`,
  `TestBackfillAndWALImagesAreByteEqualWithRules` (Abbruch „Regeltyp … ohne Fall“). Die Aussage „erfasst `map_value`
  ohne manuelle Ergänzung“ trägt also die Schleife über die Menge; die Sicherung der Menge trägt der eine Test mit
  fester Liste. Der Implementer meldet genau diesen Befund. Ein Befund für den Verifier (DoD-Formulierung), keine
  Aktion am Code.
- `verifizierbar`: ja
- `klasse`: Beleg trägt seinen Satz nicht (Zählwort)

### F-6 — Die lineare Suche der Zuordnung ist benannt, ihre Größenordnung nicht

- `kategorie`: INFO
- `quelle`: `ADR-0083`; `SPEC-030`
- `pfad`: `internal/domain/model/transformation.go:209-219` (`lookupMappedValue`), Plan §3
- `befund`: Die Grenze steht im Indikativ im Doc-Kommentar („linear in der Zahl der Zuordnungen“) und im Plan. Sie
  ist eine Herleitung aus dem Quelltext und nicht als Messung ausgegeben, das ist konform. `SPEC-030` bindet die Zahl
  der Paare von `values` nicht nach oben. `slice-transformationen-betriebsdoku` nennt „linear“ oder „Paare“ nicht
  (gesucht in `docs/plan/planning`). Die Kosten fallen je Wert, je Spalte mit Regel und je Zeile eines Backfill-Blocks
  an. Es gibt keine Messung im Repo.
- `verifizierbar`: nein
- `klasse`: benannte Grenze ohne Adresse für den Betreiber

### F-7 — Anhängen von Testtext per `cat >>` und die Grenze von `AGENTS.md` §3.1

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.1 (Text-Umschreiben), `MR-003`
- `pfad`: keine Spur im Diff (`make fmt-check`, `make kommentar-kennungen` ohne Befund); die Angabe stammt aus dem
  Auftrag des Reviews
- `befund`: Ein Anhängen per Shell-Umleitung ist weder `sed -i`, `perl -pi` noch ein Host-Interpreter und schreibt
  nichts um. Der HIGH-Punkt „Docker-only-Verstoß“ des Skills und der Verbotssatz in §3.1 nennen in-place-Umschreiben.
  Der Weg fällt unter die „flaglosen Schreibwege“, die der Guard nach `MR-003` nicht liest und für die §3.1 kein
  ausdrückliches Urteil trägt. Die Überschrift des Absatzes („Text-Umschreiben im Repo ist Sache der Datei-Werkzeuge
  des Laufs“) lässt eine strengere Lesung zu. Kein HIGH: Die Regeln nennen ihn nicht, und das Ergebnis steht im
  geprüften Diff. Hinweis an den Architect, nicht an den Implementer.
- `verifizierbar`: nein
- `klasse`: Regelgrenze Umleitung/Anhängen

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model` Produktivcode (`transformation.go`, `transformationspec.go`):
  Kodierung eindeutig (300 000 Zufallszuordnungen); Konstruktor-Invarianten gegen `SPEC-030` (leeres `values`,
  Abbildung auf sich selbst, mehrere Quellwerte auf einen Zielwert, leere Schlüssel/Werte, U+0000 in `column`);
  Parser strikt (Nicht-Zeichenketten-Werte, `null`, unbekannte Schlüssel, Prüfreihenfolge a bis d); Doppelschlüssel
  liest den letzten Wert (`jsonb`-Vertrag `ADR-0126`); exakte Gleichheit ohne Normalisierung (`é` und `e`+U+0301
  bleiben getrennte Schlüssel).
- geprüft, ohne Befund: Ausschluss vor Regel in beiden Pfaden (`BuildRowImage` ist die eine Konstruktionsstelle für
  WAL-Mapper und Backfill): Mutation „Ausschluss entfällt“ färbt Tests in vier Paketen rot; mit vertauschter
  Typ-Reihenfolge in der Menge bricht der Eigenschaftstest bei `map_value` ab, nicht erst bei `rename_column`; der
  Run-Test färbt beide Typen rot; `map_value` lässt Schlüssel und Position der Quellspalte; ein abgebildeter Wert,
  der einer Spalte gleicht, ist kein Schlüsselkonflikt.
- geprüft, ohne Befund: Kombinationen `rename_column` + `map_value` (K1–K4). Auf verschiedenen Spalten unabhängig und
  ordnungsfrei (Mutation `rules[:1]` rot), auf derselben Spalte von K2 abgelehnt (Domänen- und Verdrahtungstest), K3
  trifft `map_value` nicht (`s.to != ""`-Wache, Mutation rot).
- geprüft, ohne Befund: DoD 2 „kein Diff im Use Case“ (Stat mit `':!*_test.go'` leer):
  `TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain` und das Verdrahtungs-Pendant werden rot, wenn der
  Use Case einen Regeltyp beim Namen nennt. Die Auslegung in §3 ist redlich für die fünf Dateien, die Regeltypen
  aufzählen oder den geforderten Use-Case-Test tragen (Wortlaut und Store-Test: F-2).
- geprüft, ohne Befund: Die zwei geänderten bestehenden Erwartungen sind zwangsläufig:
  `TestTransformationKindsIsAClosedSet` kann bei zwei Typen nicht „genau `rename_column`“ tragen;
  `TestTransformationSpecBuild` führte `map_value` als Platzhalter eines unbekannten Typs. Die übrigen geänderten
  Erwartungen in `mapper`, Run und `bootstrap` sind für `rename_column` gleichwertig umgeformt, teils strenger
  (exakter Wert statt Präfix).
- geprüft, ohne Befund: Mutationsangaben der Test-Kommentare, stichprobenartig, jeweils rot wie behauptet: Zweig
  `map_value` in `applyTransformations` entfernt; Schlüssel aus `rule.to`; Suche per `HasPrefix` und per `EqualFold`;
  Sortierung in `encodeValueMap` entfernt (`model`, `backfill`); `values` beim Bau nicht gesetzt; Typprüfung in
  `jsonStringMap` entfernt; Wache `s.to != ""`; Zielnamen-Prüfung für `map_value` in `CheckApplicable`; NUL-Prüfung
  von `column`; `rules[:1]`.
- geprüft, ohne Befund: `internal/application/usecase`, `internal/bootstrap`, `tools/schema`, `spec`: kein
  Produktivcode im Diff. `sameSet[T comparable]` bleibt gültig, weil alle fünf Felder von `Transformation`
  Zeichenketten sind; der Typ-Kommentar nennt die koppelnde Stelle.
- geprüft, ohne Befund: Kommentare nach `AGENTS.md` §3.7 in den hinzugefügten Zeilen: keine Kennungskette, keine
  Kompaktform, kein „ff.“; Konjunktiv-/Chronik-Muster ohne Treffer außer zwei Zustandsbeschreibungen („nicht mehr
  lesbare Zeile“, „bildet den Namen nicht mehr ab“); der Block von `ErrInvalidTransformation` trägt eine statt zwei
  Kennungen.
- geprüft, ohne Befund: Suchlauf-Feld (§3.13): alle 22 Zeilen an beiden Ständen gleich; Trefferaufschlüsselung im Plan
  stimmt (11 Zeilen `rename_column` im Code: sechs Domäne, fünf in zwei Runner-Skripten); Handbuch, SDK-Pfade,
  `examples/`, `proto/` und `README` nennen keinen Regeltyp; das Handbuch trägt einen benannten Aufschub mit Adresse
  (`slice-transformationen-betriebsdoku`, dessen §2 nennt `map_value` und die Rückfall-Grenze). F-4 ist die einzige
  Lücke des Musters.
- geprüft, ohne Befund: Traceability/Lifecycle: alle neun Betreffe tragen eine `ADR-*`- oder `LH-*`-Kennung, keine
  `SPEC-*`/`ARC-*` im Betreff; die beiden Moves (`ce7acc8c`, `4c5000e9`) sind reine Renames (0 Zeilen), die
  Inhaltsänderung `fa5c2e29` ein eigener Commit (§3.3).
- geprüft, ohne Befund: Docker-only (§3.1) am Diff: kein Skript und kein Host-Werkzeug im Diff; kein `//nolint`,
  keine Coverage-Ausnahme (§3.2).
- geprüft, ohne Befund: Handbuch-Regeln des Skills: `docs/user` im Diff unberührt; `map_value` erweitert die Werte
  von `rule_spec`, keine neue Funktion, keine neue Umgebungsvariable.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kommentar-Zusage ohne Bindung (Test-Kommentar) · Nachzug widerspricht dem Nachbarn
im selben Träger · Test-Aussage breiter als die Messung · Träger-Nachzug (Suchmuster ohne die bewegte Beschreibung) ·
Beleg trägt seinen Satz nicht (Zählwort) · benannte Grenze ohne Adresse für den Betreiber · Regelgrenze
Umleitung/Anhängen

## Verdikt

**Merge-blockierend:** nein für Domäne und Verdrahtung — 0 HIGH, kein Produktivcode-Befund. F-1 (MEDIUM) ist ein
Test-Kommentar samt fehlender Assertion; der Fix liegt im Test. Kleine Fixrunde am Implementer (F-1, F-2 bis F-4 nach
dessen Ermessen). Die DoD-Zeile „Review durchgeführt“ bleibt offen und wird bei Schritt 21 des Implementer-Ablaufs
nachgezogen.

**Übergabe:** Findings an den Implementer. F-5 und F-2 (DoD-Wortlaut) zusätzlich an den Verifier, F-6 an den Planner
(Übergabe an `slice-transformationen-betriebsdoku`), F-7 an den Architect. Die Finding-Klassen gehen in die
Slice-Closure §7 und von dort in den Zähler. Der Report ist ein Lauf-Beleg und ersetzt keine Verifikation.
