# Review-Report: slice-089 — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier, Modul 11),
§6-Risiko-Ausgänge, Beobachtungs-Register und die drei Paarungen
(Planner-Closure) sind **nicht** Gegenstand dieses Reports.

**Gegenstand:** `slice-089`, Diff `c9c3963..960fee4` (zwei Commits, HEAD
`960fee4`, Baum sauber) — **sieben Dateien, alle Doku**: `AGENTS.md`,
`.harness/skills/reviewer.md`, `.claude/commands/implement-slice.md`,
`harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`,
`docs/plan/planning/welle-20.md` und der Slice-Plan selbst. Keine hinzugefügte
Datei, kein Produktionscode, `harness/README.md` und `AGENTS.md` §4 unberührt.

**Skill:** `.harness/skills/reviewer.md` @ `960fee4` — die zwei neuen
HIGH-Unterpunkte gehören zum Gegenstand und werden hier zum ersten Mal
angewandt · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-089` vollständig (§1–§8); offene Welle `welle-20`
- `ADR-0083` (Herkunfts-Regel: §Entscheidung 1/2 Wortlaut, 4 durchsetzende
  Hälfte, 6 Geltungsbereich, §Konsequenzen Folgepflichten, §Die benannte
  Grenze), `ADR-0084` (Sync-Gate — ausdrücklich **nicht** Gegenstand),
  `ADR-0078` §Fitness Function Zeile 3, `ADR-0082` (Nenner/Cluster),
  `ADR-0071`/`ADR-0077` (Messgegenstand, Rampe; `Accepted` = immutabel)
- `docs/reviews/architect-verdict-negativtest-eingabeseite-4x.md` (der
  „Wortlaut — eine Regel, beide Träger"), `review-slice-088.md` (F-1),
  `review-slice-088-delta.md`, `verify-slice-085.md` (V-1/V-2/V-3),
  `verify-slice-081.md`
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/` (Zähler
  `ls evidence/` nachgezählt, nicht aus Prosa übernommen)
- `AGENTS.md` §3.5, §3.7, §3.9, §3.11, §3.12 (der neue Träger), §4, §5, §6 ·
  `harness/conventions.md` (MR-000 ID-Schema)

---

## Findings

### F-1 — §Zählbasis nennt `71.94%` als Ausgabe des Gate-Skripts; das Skript gibt `71.90%`

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz A („Ist sie eine Messung, trägt sie
  zusätzlich den **Lauf**, aus der sie stammt (die **gedruckte Zeile**)") ·
  `ADR-0083` §Entscheidung 1 · Klasse
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `harness/sensors/coverage-gate.md:73` (gegen `:62`, gegen
  `tools/coverage-gate.sh:35`, `:47`)
- `befund`: Der Absatz begründet, die gedruckte Zeile sei „dieselbe Messung in
  anderer Ausgabepräzision", und schließt: „`go tool cover` druckt eine
  Nachkommastelle, `tools/coverage-gate.sh` formatiert `%.2f` — daraus werden
  `71.9%` und `71.94%`." Das Skript formatiert den **geparsten gedruckten**
  Wert (`total_pct` aus der `total:`-Zeile): `71.9` → `71.90%`; gemessen am
  synthetischen Profil **und** live am eigenen Gate-Lauf (gedruckt `74.7%` →
  `coverage-gate: OK — Coverage 74.70%`). `71,94 %` ist die **deduplizierte
  Rechnung** (1369/1903) und steht vier Zeilen darüber als „**71,94 %**,
  gedruckt `71.90%`" — dieselbe Sektion führt damit **zwei verschiedene
  gedruckte Zeilen** für denselben Stand. Die Zeile liegt in einem Träger
  dieses Diffs und ist in der §3-Nachzug-Liste des Slice-Plans **nicht** als
  liegen gelassene Fundstelle geführt (dort stehen nur V-1 §4.3/§4.4 und V-3).
- `verifizierbar`: ja — `LC_ALL=C bash tools/coverage-gate.sh <file> 70` über
  ein `total:`-File mit `71.9%` druckt `Coverage 71.90% erfüllt Schwelle 70%`
  (Exit 0); die Gegenprobe mit `71.94%` druckt `71.94%` — der Wert stammt aus
  der Eingabezeile, nicht aus der Formatierung
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`

### F-2 — Die Herkunfts-Erklärung nennt vier Werte, die der Abschnitt nicht führt

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 (eine Aussage über den Gegenstand trägt ihren
  Beleg-Anker) · Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `docs/plan/planning/welle-20.md:117-118` (gegen `:84-113`)
- `befund`: Die Erklärung sagt „Die **Summen, Differenzen und Prozente**
  **dieser Sektion** (`229–231`, `175–177`, `165–167`, `1523`, `154`, `13`,
  `63`, `3,3 pp`, `81,2 %`, `84,1 %`, `83,6 %`)". `175–177`, `165–167` und
  `83,6 %` kommen in `welle-20.md` **nirgends** vor (`grep` über die ganze
  Datei: Treffer nur innerhalb der Aufzählung selbst); `1679` steht in §3
  (`:25`), nicht in §4. Der Satz beansprucht damit Zahlen als Zahlen des
  Abschnitts, die nur in `ADR-0082` stehen (`175–177` §Kontext (4b)/§Schnittmaß,
  `165–167` §Schnittmaß, `83,6 %` §Kontext-Tabelle „Realistische Decke").
- `verifizierbar`: ja — `grep -nE "175|177|165|167|83,6" docs/plan/planning/welle-20.md`
  liefert genau die zwei Zeilen der Aufzählung
- `klasse`: „Herkunfts-Erklärung zählt Werte, die der Träger nicht führt"

### F-3 — §8 zählt drei Einträge über der Schwelle; die Liste darüber markiert vier

- `kategorie`: MEDIUM
- `quelle`: Maintainability · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung (Schritt 2 *Offene Beobachtungen sichten*) · Klasse
  `zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `docs/plan/planning/in-progress/slice-089-regeln-verkoerpern.md:295`
  (gegen `:282-289`)
- `befund`: Die Ergebnis-Zeile sagt „drei Einträge stehen **über** der
  Schwelle, und dieser Slice trägt ihre Ausgänge". Die fünf Zeilen darüber
  markieren **vier** Einträge mit „Schwelle erreicht" (`zahl-in-traeger-…`,
  `negativtest-…`, `dod-begruendung-…`, `generierte-artefakte-ohne-sync-sensor`
  — alle vier tragen vier Beleg-Dateien) und einen mit „unter der Schwelle".
  Gelesen als Zählung der gesichteten Einträge ist „drei" um eins zu niedrig;
  gelesen als Aussage über das Register ist sie weiter daneben (eigene
  Zählung: **21 von 59** Einträgen tragen ≥ 3 Beleg-Dateien). Dieselbe Form
  wie `review-slice-088` F-3 (§8-Zählung gegen die Ablage), dort ebenfalls
  MEDIUM.
- `verifizierbar`: ja — `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence/`
  gegen die Markierungen der Liste
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`

### F-4 — Schritt 19 begründet mit einer Mutation, die nirgends dokumentiert ist

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz B (Beleg-Anker oder als **erwartet**
  formuliert) — durchsetzende Hälfte laut `ADR-0083` §Entscheidung 4 ist der
  **Verifier**, nicht der Reviewer
- `pfad`: `.claude/commands/implement-slice.md:186` (gegen
  `docs/reviews/review-slice-088.md:235`)
- `befund`: Der Satz sagt, in `slice-088` habe die Pflicht „zwei Aussagen
  trotzdem ungebunden" gelassen, „weil mutiert wurde, was der Test
  **zurückgibt**". `review-slice-088.md` hält ausdrücklich fest: „Die fünf
  Mutationen des Implementers sind im Diff nicht dokumentiert" — die
  Richtungs-Aussage über sie hat damit keinen auflösbaren Anker; belegt sind
  nur die Mutationen des Reviewers (M1–M5) und deren Ergebnis. Die tragende
  Hälfte des Satzes („fünf gefahren, zwei Aussagen ungebunden,
  `review-slice-088.md` F-1") trägt.
- `verifizierbar`: ja — `review-slice-088.md` §Eigene Messungen gegen die
  Behauptung über den Implementer-Lauf
- `klasse`: „Tatsachenbehauptung ohne Beleg-Anker (Instanz B)"

### F-5 — Die liegen gelassenen Fundstellen sind benannt, aber ohne Adresse

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Offene Risiken werden bei Closure aufgelöst · Slice-Plan §6 Risiko 3
  („weitere Funde gehören als eigene Adresse geführt")
- `pfad`: `docs/plan/planning/in-progress/slice-089-regeln-verkoerpern.md:168-170`
  (und `:256`, `:158-170`)
- `befund`: V-1 §4.3 (die „132 Positionen × 2"-Lesart in
  `db-adapter-coverage.md`) und V-3 (der `go list -f '{{len .TestGoFiles}}'`-Befehl
  in `coverage-gate.md` §Grenze Punkt 1, der 25 statt fünf liefert) sind als
  „nicht nachgezogen" benannt und begründet; eine **Adresse** (Kennung eines
  Folge-Slice) trägt keine von ihnen, und §7 der Closure-Notiz ist noch
  Platzhalter. Die **dritte** liegen gelassene Fundstelle — die F-1-Zeile —
  ist in §3 gar nicht genannt. Für V-3 ist die Fundstelle die von LP3
  namentlich genannte Sektion — die Zahlenseite dort ist nachgezogen, die
  Beleg-Befehl-Seite nicht.
- `verifizierbar`: ja — die §3-Nachzug-Liste und §7 des Plans gegen die
  genannten Fundstellen
- `klasse`: „benannter Fund ohne Adresse"

## Negativbefunde

- geprüft, ohne Befund: `AGENTS.md` §3.12 (Z. 366–422) — **wörtlich**: die
  beiden `> `-Blöcke sind byte-identisch zu `ADR-0083:143-167` (eigene Probe,
  leerer `diff`); die Träger-Aussage („die vier belegten Fälle stehen in
  `…/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`, er hat sie durch
  eigenes Nachmessen gefunden") löst auf: Verzeichnis existiert, vier
  Beleg-Dateien (`slice-081`/`-084`/`-085`/`-088`), und alle vier Berichte
  führen eigene Messungen des Reviewers. Die weggelassene ADR-Einleitungszeile
  („Der Unterabschnitt trägt:") und die aus zwei ADR-Stellen gezogene
  Grenz-Formulierung sind **Form**, keine Inhaltsänderung: der gestrichene
  Halbsatz „und wird nicht als künftiger bestellt" widerspräche sogar
  `ADR-0083` §Re-Evaluierungs-Trigger (a); die Beleg-Klausel „sie hat in
  `slice-081`, `-084` und `-085` sechs Werte im Baum widerlegt" trägt keine
  Regel-Aussage. `· seit slice-089` steht in der Form von §3.11.
- geprüft, ohne Befund: `.harness/skills/reviewer.md` (Z. 100–132) — beide
  Unterpunkte stehen in der HIGH-Liste, in der Hausform der übrigen
  repo-spezifischen Punkte; die Herkunfts-Angaben halten (`4×` und die
  Beleg-Listen stimmen mit `ls evidence/` überein, „in drei der vier Fälle …
  durch Mutieren der Eingabeseite" stimmt mit dem Architect-Verdikt überein);
  die Abgrenzung zur MEDIUM-Klasse und „kein Gate fängt das" sind gesetzt; der
  Wortlaut der Mutations-Regel ist identisch mit dem Verdikt-Zitat (eigener
  `diff`, leer). Die Erweiterung auf **zwei** Unterpunkte ist Folgepflicht aus
  `ADR-0083` §Entscheidung 4 (dort ist der Skill-Unterpunkt für Instanz A
  namentlich als Zielort geführt) und in §3 des Plans als Nachzug dokumentiert
  → **im Zuschnitt**. Die Zeile „vier repo-spezifische HIGH-Regeln" (Z. 180)
  trägt ihren Stand („seit 2026-09-09") und ist damit regelkonform, als
  Zählung heute überholt.
- geprüft, ohne Befund: `.claude/commands/implement-slice.md` Schritt 19
  (Z. 181–196) — „fünf gefahren, zwei Aussagen ungebunden, F-1" hält gegen
  `review-slice-088.md` F-1 und `evidence/slice-088.md`; die Regel steht
  wörtlich wie im Verdikt; „Enumerations-Pflicht statt Erinnerung" ist die vom
  Verdikt verlangte Form. Nur die kausale Hälfte ist F-4.
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` §Kalibrierungs-Bindung,
  §Grenze Punkt 1, §Ausgabe (Z. 32, 93, 106-111, 162-165) — alle dort
  nachgezogenen Werte habe ich **nachgemessen** (§Eigene Messungen): `20/20`
  (`mapper`), `86/61` und die abgeleitete Zeile `70,9 %`, `49/0` (`cmd`), der
  Nenner `1903`; die Lauf-Anker (`slice-089`, `slice-081`) sind auflösbar
  (`verify-slice-081.md` trägt beide realen Kommandos mit `71.30%`). Nur
  `:73` ist F-1.
- geprüft, ohne Befund: `harness/sensors/db-adapter-coverage.md` (Z. 88, 169)
  — `477 von 650 = 73,38 %` trägt denselben Lauf in der Kalibrierungs-Zelle und
  im Rot/Grün-Beleg; `slice-081` ist der richtige Lauf (die Naht entsteht in
  `slice-081`: `ADR-0077` §Geschichte). Die Aufnahme dieses Dokuments ist die
  Lesart, die V-1 §4.2 erfüllt (V-1 nennt **namentlich** zwei
  §Ausgabe-Abschnitte, einen je Dokument) → LP3-Zeile **erfüllt**, kein
  Überlauf über den Zuschnitt.
- geprüft, ohne Befund: `docs/plan/planning/welle-20.md` §4 Herkunfts-Absatz
  (Z. 110–122) in der Sache — `ein eigener Messstand: Lauf slice-089` hält
  (mein Dedup misst `1903`, gedruckt `74.7 %`); die Zuordnung „übernommen aus
  `ADR-0082`" trifft für `336`, `308`, `26–28`, `198`, `0 von 198`, `591`,
  `327` zu. Nur die Aufzählung ist F-2.
- geprüft, ohne Befund: `harness/README.md` §Sensors und `AGENTS.md` §4 —
  **unberührt** (`git diff --name-status` führt `harness/README.md` nicht;
  der `AGENTS.md`-Hunk endet vor `## 4. Quality Gates`). Kein Target, das es
  nicht gibt, wird genannt; `ADR-0084` bleibt in beiden ausgespart.
- geprüft, ohne Befund: Hygiene des Diff-Ranges — keine hinzugefügte Datei,
  kein Lauf-Artefakt (kein `tools/schema/plan.yaml`/`down.sql`), keine
  Vorlagen-Reste, keine Platzhalter (`<NNN>`, `<…>`, `TODO`), keine
  host-lokalen Pfade (§3.11; `hostpaths`-Modul 0 Befunde), beide Betreffe
  nennen `ADR-0083` und keine Struktur-ID, kein Trailer, Historie linear auf
  `main`, Baum sauber.
- geprüft, ohne Befund: Register-Zähler des Plans §8 — `zahl-in-traeger-…` 4×,
  `negativtest-…` 4×, `dod-begruendung-…` 4×, `generierte-artefakte-…` 4×,
  `db-gegenstand-…` 2× — alle fünf stimmen mit `ls evidence/` überein (nicht
  aus Prosa gelesen). Nur die Ergebnis-Zeile ist F-3.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-status c9c3963..960fee4` | **0** | 7 `M`, keine `A`/`D` |
| Wörtlich-Probe: `sed -n '143,167p' ADR-0083 \| grep '^>'` → `21` Zeilen; §3.12-Zitate (`awk`-Bereich) → `21` Zeilen; `diff` | **0** | **leerer `diff`** — die „20 vs. 20" der Übergabe sind um eins daneben (die zwei leeren `>`-Zeilen), die Substanz trägt |
| `make gates` (Log in Datei, Exit **danach** in eigener Datei gelesen, §3.9) | **0** | `baseline-verify v6.5.0 OK` (54 Dateien) · d-check **725** Dateien / **0** Befunde · `commit-traceability: OK — 5 Commit(s), Betreffe ohne Struktur-ID` · `coverage-gate: OK — Coverage 74.70% erfüllt Schwelle 70%` · a-check **0** Befunde |
| `make docs-check` (eigener Aufruf, Exit direkt) | **0** | `d-check: 725 Datei(en) geprüft, 0 Befund(e)` |
| `docker run --rm --entrypoint sh pg-change-feed:coverage -c 'cat /out/coverage.out'` + `awk`-Dedup über die Block-Position | **0** | **1422 von 1903 = 74,72 %**; `postgresstorage/mapper` **20 von 20**; `grpc/streamv1` **61 von 86** (61/86 = 70,9 %); `cmd/pg-change-feed` **0 von 49** — **alle drei Zahlen der §Grenze Punkt 1 halten** |
| `LC_ALL=C bash tools/coverage-gate.sh <file 71.9%> 70` | **0** | `coverage-gate: OK — Coverage 71.90% erfüllt Schwelle 70%` (F-1) |
| `LC_ALL=C bash tools/coverage-gate.sh <file 71.94%> 70` | **0** | `coverage-gate: OK — Coverage 71.94% erfüllt Schwelle 70%` — der Wert kommt aus der Eingabezeile |
| `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence/` (fünf Slugs) | **0** | 4 · 4 · 4 · 4 · 2 — deckt §8 des Plans |
| `for d in …/BEO-PGC/*/; do ls evidence/ \| wc -l; done` (Register gesamt) | **0** | **59** Einträge, davon **21** mit ≥ 3 Beleg-Dateien (Zählung für F-3) |
| `grep -nE "175\|177\|165\|167\|83,6" docs/plan/planning/welle-20.md` | **0** | Treffer **nur** in der Aufzählung (Z. 117/118) — F-2 |

Der Gate-Lauf und seine Auswertung sind **zwei Schritte**: `make gates` wurde
ungepiped in eine Log-Datei gefahren, der Exit-Code in
`/tmp/gates-089.ec` separat geschrieben und erst danach gelesen; kein
Folgekommando hing an einer Pipe.

## Antwort auf die Schwerpunkte

**(1) „Wörtlich" — trägt, mechanisch nachgeprüft.** Meine Probe ergibt 21 zu
21 Zitat-Zeilen und einen **leeren** `diff` (nicht „20 zu 20"; die Abweichung
ist die Zählung der zwei leeren `>`-Zeilen und liegt nur in der Übergabe, in
keinem Artefakt). Die beiden Arrangements sind **Form**: die gestrichene
Zeile „Der Unterabschnitt trägt:" ist ADR-Regie, und die zusammengezogene
Grenze lässt nur „und wird nicht als künftiger bestellt" sowie die
Beleg-Klausel weg — ersteres stünde in Spannung zu `ADR-0083`
§Re-Evaluierungs-Trigger (a), letzteres trägt keine Regel-Aussage. Der Zusatz
in §3.12 (durchsetzende Leser = Reviewer/Verifier, Träger außerhalb des Diffs
= Messung) ist ADR-Inhalt (§Entscheidung 4 Tabelle, §Die benannte Grenze),
keine Erweiterung.

**(2) Die zwei Erweiterungen — beide tragen.**
`.harness/skills/reviewer.md` **muss** zwei Unterpunkte tragen: `ADR-0083`
§Entscheidung 4 führt „`AGENTS.md` §3.12 **+ ein benannter HIGH-Unterpunkt in
`.harness/skills/reviewer.md`**" als Zielort von Instanz A, §Konsequenzen
wiederholt es als Folgepflicht, und die ADR bestellt ausdrücklich **keinen**
Sensor — ohne den Unterpunkt hätte die Zahlen-Hälfte keinen Leser. Der Plan
hatte in LP2 nur die Mutations-Hälfte genannt; der Nachzug ist in §3
dokumentiert und damit plan-konform. `harness/sensors/db-adapter-coverage.md`
ist ebenfalls **im Zuschnitt**: V-1 §4.2 nennt die zwei §Ausgabe-Abschnitte
**namentlich mit ihren Werten** (`73,38 %` in `db-adapter-coverage.md`,
`71,30 %` in `coverage-gate.md`) — die Lesart, die beide erfüllt, ist die
einzige, die V-1 §4.2 ganz abdeckt; `coverage-gate.md` allein hätte eines der
zwei namentlich genannten Belege stehen lassen. Die LP3-Zeile ist damit
**erfüllt**, nicht übererfüllt.

**(3) Die Regel gegen sich selbst — vier Stellen gefunden.** Gegen §3.12
geprüft wurden §3.12 selbst, beide Skill-Unterpunkte, Schritt 19, die drei
Sensor-Dokumente und `welle-20.md`: **F-1** (falscher Ursprung für `71.94 %`
in `coverage-gate.md` — die Zahl, die als „gedruckt" ausgegeben wird, ist
nicht die gedruckte, und die Sektion widerspricht sich selbst), **F-2** (die
neue Herkunfts-Erklärung in `welle-20.md` nennt vier Werte, die der Abschnitt
nicht führt), **F-3** (§8-Zählung des Plans gegen die eigene Liste und die
Ablage), **F-4** (Schritt 19 begründet mit einer nicht dokumentierten
Mutation). Alle Zahlen der **nachgezogenen** Stellen (§Kalibrierungs-Bindung,
§Grenze Punkt 1, §Ausgabe, `db-adapter-coverage.md`) habe ich nachgemessen —
sie halten. Die zwei vom Implementer gemeldeten Selbst-Korrekturen sind damit
nicht die einzigen; die Regel ist in ihrer eigenen Einführung **nicht**
vollständig durchgesetzt, und der schwerste Fall (F-1) stand gar nicht in der
Nachzug-Liste.

**(4) Die liegen gelassenen Fundstellen — V-1 §4.3/§4.4 und V-3 tragen als
Benennung; der neue Fund trägt.** V-1 §4.3 („132 Positionen × 2") und V-3
(`go list`-Befehl) sind in §3 benannt, begründet und über §6 Risiko 3 als
eigene Klasse geführt — das ist die vom Plan verlangte Form (eine Adresse
fehlt noch, F-5, und sie gehört in §7 der Closure). Der **neue** Fund trägt:
`tools/coverage-gate.sh:35` parst den gedruckten Wert, `:47` formatiert ihn mit
`%.2f`; `71.9` wird zu `71.90%` — live am eigenen Gate-Lauf bestätigt
(gedruckt `74.7%` → `Coverage 74.70%`). Damit ist `71.94 %` nicht die Ausgabe
des Skripts, sondern die deduplizierte Rechnung, die dieselbe §Zählbasis vier
Zeilen darüber bereits als „gedruckt `71.90%`" ausweist.

**(5) Kein Gate, das es nicht gibt — hält.** `AGENTS.md` §4 und
`harness/README.md` §Sensors sind diff-frei (`git diff --name-status` führt
`harness/README.md` nicht; der `AGENTS.md`-Hunk endet vor `## 4. Quality
Gates`). `ADR-0084`s Sync-Gate wird in keiner der beiden Stellen genannt; der
Plan schließt es in §3 ausdrücklich aus. Kein halluziniertes Target.

**(6) Hygiene — sauber.** Keine hinzugefügte Datei, kein Lauf-Artefakt, keine
Vorlagen-Reste, keine Platzhalter, keine host-lokalen Pfade (`hostpaths`-Modul
im Gate: 0 Befunde), Betreffe ohne `SPEC-*`/`ARC-*` und mit `ADR-0083`, kein
Trailer, lineare Historie auf `main`, Arbeitsbaum sauber. Der Betreff des
zweiten Commits nennt einen Typo-Zug im Slice-Plan — der Hunk ist die
DoD-/§3-Fortschreibung, kein Fremdinhalt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:**
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (F-1, F-3) ·
„Herkunfts-Erklärung zählt Werte, die der Träger nicht führt" (F-2) ·
„Tatsachenbehauptung ohne Beleg-Anker (Instanz B)" (F-4) ·
„benannter Fund ohne Adresse" (F-5)

## Verdikt

**Merge-blockierend:** ja — 1 HIGH, 2 MEDIUM.

**Auf der Implementer-Seite:** F-1 und F-2 liegen in Trägern dieses Diffs
(`coverage-gate.md` ist eine der zwei Sensor-Dokumente, die LP3 regelkonform
machen soll; `welle-20.md` §4 hat dieser Zug geschrieben). Beide sind
Ein-Stellen-Korrekturen im bereits berührten Träger, kein wachsender
Berichtigungs-Zug über fremde Dokumente — §6 Risiko 3 bleibt unberührt, weil
nichts Neues aufgenommen, sondern eine falsche Herkunfts-Aussage berichtigt
wird. Sie erfordern eine **Rückgabe an den Implementer**.

**Review → Plan (Rückkante):** F-3 ist ein Defekt des Plans (§8), nicht der
Ausführung — er geht an den Planner, wie `review-slice-088` F-3.

**INFO:** F-4 verweist auf die **Verifier**-Hälfte von `ADR-0083`
§Entscheidung 4 (Instanz B prüft der Verifier, nicht der Reviewer); F-5 ist
der Adress-Nachtrag, den §7 der Closure tragen muss.

**Übergabe:** Rückgabe-Pfeil Reviewer → Implementer ist **nötig** (F-1, F-2);
die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine
Verifikation (Modul 11).

**DoD-Häkchen „Review durchgeführt":** bleibt **offen** — der Slice braucht
eine Fixrunde; der Nachzug läuft regulär über Schritt 21 des
Implementer-Workflows (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug
ohne Fixrunde, *Grenze*). Der Slice-Plan ist in diesem Commit **nicht**
angefasst; die Findings sind die Eingabe für die Fixrunde.
