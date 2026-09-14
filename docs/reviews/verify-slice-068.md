# Verifikationsbericht: slice-068 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-068`
§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5 Closure-Trigger, §6
Risiken, §8 Sub-Area) und die bindenden Entscheidungen
([`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) — umgesetzt,
unverändert `Accepted`; [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
— **Abgrenzung**: die Dauerhaftigkeit des Ausschlussstandes ist ausdrücklich
**nicht** Gegenstand dieses Slice, ihr Träger ist der wellenlose Folge-Slice).
Nicht gegen den Diff als solchen (Reviewer-Aufgabe;
[`review-slice-068`](review-slice-068.md) vollständig gelesen, aber nur als
Kontext) und nicht gegen realen Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8,
Stand `HEAD = 8b8fa6f`), die vollständige
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) und
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md), den
vollständigen Review-Report, die beiden Vorgänger
`docs/plan/planning/done/slice-066-spaltenausschluss-sql-funktionen.md` und
`docs/plan/planning/done/slice-067-assembler-filterung-live-reload.md`, die
Wellen-Datei [`welle-18`](../plan/planning/welle-18.md) und die betroffenen
Abschnitte des Lastenhefts. Die Sensoren wurden in dieser Sitzung
**eigenständig real ausgeführt** (Exit-Code je in einem eigenen, ungepipten
Schritt, `AGENTS.md` §3.9) — kein Implementer- oder Reviewer-Beleg ungeprüft
übernommen. Beide Mutationsläufe sind **eigene**, am Skript-Code gesetzt.

**Gegenstand:** `docs/plan/planning/in-progress/slice-068-e2e-spaltenausschluss.md`
zum Stand `HEAD = 8b8fa6f`. Drei Commits seit `d088bb0` (reiner
`next→in-progress`-Move, 0 Inhaltszeilen): `7fa5784` (Rundlauf + `harness/README.md`
+ `plan.yaml` + Plan-Nachzüge + DoD-Häkchen), `bd8b161` (Review-Report +
Review-Häkchen), `8b8fa6f` (Plan-Korrektur F-1/F-2 und `LH-QA-SEC-004`-Bindung
im `harness/README.md`).

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — die implementierungs-/
reviewbezogenen Zeilen sind Prüfgegenstand, die Closure-Zeilen **müssen** offen
bleiben, bis der Planner sie schließt.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `LH-FA-CFG-005` Happy Path real belegt (Antrag → Poll `applied` → neue Change ohne Schlüssel, historische Change geprüft) | **erfüllt, selbst reproduziert** | Eigener `make test-integration`-Lauf Exit **0**; die Beleg-Zeile steht real in der Ausgabe. Am Skript-Code: eigene Tabelle `feed_e2e_column_exclusion` (nicht in `CDC_TABLES`, `compose.yaml:89`), Aktivierung über `cdc.enable_table`, `SELECT cdc.exclude_column(...)`, Poll auf `status = 'applied'`, danach drei Zusagen zur Change `id=2` (Schlüssel fehlt in `new_data`, Wert in `new_data`/`old_data` nirgends, `name` bleibt), die Change `id=1` bleibt unverändert lesbar. Mutation A (§3) macht genau diese Zusage rot. |
| 2 | `LH-FA-CFG-005` Negative real belegt (`failed` samt `ErrSourceColumnMissing`-Text) | **erfüllt, selbst reproduziert** | Skript `:443` ruft `cdc.exclude_column(… 'nicht_vorhandene_spalte')` gegen die laufende Bindung ab, pollt auf `status = 'failed'`, prüft `error_message` gegen `"Spalte existiert nicht an der Quelle"` — wörtlich `ErrSourceColumnMissing.Error()` (`internal/application/port/inbound/verwaltung.go:26`). Eigener Lauf grün; Mutation B (§3) macht genau diese Zusage rot. |
| 3 | `make gates` grün, `make test-integration` real grün (voller Compose-Stack) | **erfüllt, selbst reproduziert** | Zwei eigene, ungepipste Läufe, beide Exit **0** (§2). `make test-integration` ist **kein** Gate (`harness/README.md` §Sensors, Bindung `kein Gate, ADR-0030, LH-QA-POR-003`) — der Beleg trägt über den Exit und die zwei Ausgabezeilen, nicht über die Gate-Summe. |
| 4 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | [`review-slice-068`](review-slice-068.md) vorhanden, im eigenen Commit `bd8b161`; Verdikt 0 HIGH / 1 MEDIUM / 1 LOW / 3 INFO deckungsgleich mit §2 der DoD-Zeile. Das Review-Häkchen wurde in **genau diesem** Commit gezogen (§6), nicht im Implementer-Commit. |
| 5 | Doku-Update: `harness/README.md` §Sensors, Zeile `make test-integration` | **erfüllt** | Der Satz nennt real, was das Skript tut (aktivierte, nicht in `CDC_TABLES` gelistete Tabelle, Poll auf `applied`, drei Aussagen zur danach erfassten Change, unveränderte Lesbarkeit der davor erfassten, Negative mit Fehlertext); nach F-5 zusätzlich die tragende `LH-QA-SEC-004`-Kennung (`8b8fa6f`). Herkunfts-Anker `· seit slice-068`, Bindungs-Spalte unverändert. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<bei Closure>`-Platzhalter (Volltext gelesen). Planner-Arbeit nach diesem Bericht. |
| 7 | Reconciliation-Register — entfällt (Greenfield) | **korrekt offen, Entfall-Vermerk trägt** | Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`); `docs/plan/planning/reconciliation.md` existiert real nicht. |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/` real gelistet: kein `slice-068`-Beleg; die §8-Sichtung nennt zwei Treffer (`test-integration-retention-timing-flake`, `test-isolation-geteilter-zustand`) — beide real **1×** (je eine `evidence/`-Datei), keiner erreicht 3×. Die Eintragung gehört in die Closure. |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Einträge tragen wörtlich `<bei Closure zuzuweisen>`; kein Risiko vorzeitig geschlossen. |
| 10 | Die drei Paarungen | **korrekt offen** | Slice liegt real in `in-progress/`; `Welle: welle-18`-Feld vorhanden; die DoD-Zeile verweist korrekt auf die `welle-18`-Closure. |

**Ergebnis §1:** Alle fünf implementierungs-/reviewbezogenen DoD-Punkte (1–5)
sind real erfüllt; 1–3 wurden **selbst reproduziert**, die zwei Kernzusagen
zusätzlich über zwei eigene Mutationen als einzeln tragend belegt. Die fünf
Closure-Punkte (6–10) sind korrekt noch offen und wurden nicht vorweggenommen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` (Integrität + Vollständigkeit, netzlos); d-check 531 Dateien/0 Befunde; d-check `--enable commits --range HEAD~5..HEAD` 0 Befunde; `commit-traceability: OK — 5 Commit(s)`, Betreffe ohne Struktur-ID; a-check gesamt 0 Befunde (2 unverschichtete Werkzeug-Dateien, bekannte Hinweise); `coverage-gate: OK — Coverage 45.80% erfüllt Schwelle 35%` |
| `make test-integration` | **0** | voller Compose-Stack-Lauf; beide neuen Beleg-Zeilen real in der Ausgabe (`Spaltenausschluss-Rundlauf (LH-FA-CFG-005 Happy Path, ADR-0059) belegt …` und `Spaltenausschluss-Negative-Beleg (LH-FA-CFG-005) …`) |
| `make doc-commits RANGE=d088bb0..HEAD` | **0** | 531 Dateien, 0 Befunde — jeder der drei Slice-Commits ist kennungstragend |
| `make doc-immutable RANGE=d088bb0..HEAD` | **0** | 0 Befunde — keine `Accepted`-ADR im Range überschrieben (`ADR-0059` bleibt `Accepted`, `ADR-0065` nur ergänzt) |

`git status --porcelain` nach allen Sensor- und Mutationsläufen: **leer** —
insbesondere reproduziert der Lauf `tools/schema/plan.yaml` lokal
(F-4-Beobachtung des Reviews bestätigt: die committete Ziel-DSN-Zeile
`cdc-test-postgres:5432/cdc` stellt sich im selben Lauf wieder her).

## 3. Mutations-Stichprobe — DoD-Zusage real rot gesehen (Verifier-only-Nachweis)

Der Implementer und der Reviewer nennen je zwei rot gesehene Mutationen. Ich
habe **beide** selbst gestellt, jede in einem eigenen `make test-integration`-
Aufruf, danach per `git checkout` zurückgenommen. Der Blob-Hash der
zurückgenommenen Datei stimmt mit `HEAD` überein
(`git hash-object` = `git rev-parse HEAD:tools/harness/run-integration-tests.sh`
= `cc3f399…`), `git status --porcelain` ist nach jeder Rücknahme leer.

| # | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| A | Antragsart in der Happy-Path-Anforderung verfälscht: `cdc.exclude_column(...)` → `cdc.include_column(...)` (`run-integration-tests.sh:349`) | rot an der Row-Image-Zusage | **rot** (Exit 2; genau `…die nach dem Ausschluss erfasste Change (id=2, feed_e2e_column_exclusion) trägt den ausgeschlossenen Spaltenschlüssel secret weiterhin im Row Image`) |
| B | Katalog-Blick ersetzt die nicht existierende Spalte: `'nicht_vorhandene_spalte'` → `'name'` (`run-integration-tests.sh:443`) | rot an der Negative-Zusage | **rot** (Exit 2; genau `…cdc.exclude_column gegen die nicht existierende Spalte endete mit status=applied, wollen failed (LH-FA-CFG-005 Negative: expliziter Fehlerpfad)`) |

**Ergebnis §3:** Beide Zusagen des Abschnitts sind **einzeln** tragend, keine
ist Dekoration. Der Happy-Path-Beleg hängt am Ausschluss (Mutation A zeigt: ohne
den `exclude_column`-Antrag bleibt der Schlüssel im Row Image) und nicht am
Dekodier-/Publication-Pfad; der Negative-Beleg hängt an der realen
Nicht-Existenz der Spalte (Mutation B: eine existierende Spalte endet
`applied` und bricht die Zusage). Im Mutationslauf B blieb der Happy-Path-Beleg
**grün** und nur die Negative-Zeile rot — die zwei Hälften sind unabhängig
voneinander sensibel, keine trägt die andere mit.

## 4. Lage vor der Container-Ende-Grenze und Weiterlaufen des Feed-Containers (am Skript-Code, nicht am Log)

- Der neue Abschnitt beginnt bei `run-integration-tests.sh:280` und endet bei
  `:481`. Die Container-Ende-Grenze der beiden Schema-Negative-Testfunktionen
  liegt bei `:1765` (`TestE2ESchemaChangeDropColumn`, `:1868`) und
  `:1946` (`TestE2ESchemaChangeIncompatibleTypeChange`, `:1963`). Der Rundlauf
  liegt also **vollständig davor**.
- Zwischen `:280` und `:481` steht **kein** `docker restart`/`docker pause`/
  `up`-Aufruf für den Feed-Container; die ersten Neustart-/Pause-Aufrufe danach
  sind `docker pause` (`:540`, Lasttest-Beleg) und `docker restart` (`:838`,
  simulierter Container-Neustart des Black-Box-Rundlaufs) — beides **nach**
  unserem Abschnitt. Die beiden `docker inspect --format '{{.State.Running}}'`-
  Wächter (`:431`, `:475`) sind reine Beobachtungen und bestätigen nach jeder
  Beleg-Hälfte real `true`.
- Der Abschnitt ist eine Shell-Strecke und berührt das Go-`-run`-Muster
  (`:277-279`) nicht; er legt keine Tabelle an, die ein anderer Abschnitt liest
  (die Tabelle kommt ausschließlich hier vor, geprüft über den ganzen Bestand),
  und liest keine, die ein anderer anlegt. `BEO-PGC/test-isolation-geteilter-zustand`
  und `BEO-PGC/test-runner-stiller-ausschluss` werden dadurch weder berührt noch
  verschärft.
- Das Skript läuft unter `set -euo pipefail` (`:44`); jede der fünfzehn
  `if …; then echo … >&2; exit 1; fi`-Prüfungen des Abschnitts bricht den Lauf
  an der zugesagten Stelle ab — kein stiller Durchlauf.

## 5. Baseline und Kontrolle — die zwei Alternativerklärungen, real ausgeschlossen

Die Frage war, ob der Happy-Path-Beleg auch anders erklärbar wäre (eine Spalte,
die nie im Row Image war; ein leeres Row Image). Beide Erklärungen sind durch
zwei **eigene, im selben Lauf liegende** Beobachtungen ausgeschlossen:

- **„Die Spalte war nie Teil des Row Image".** Vor dem Ausschluss legt der
  Abschnitt eine Change (`id=1`) mit einem eindeutigen Sentinel
  (`ColumnBeforeExclusionSentinel`) an und verlangt real
  `new_data->>'secret' = 'ColumnBeforeExclusionSentinel'`. Damit ist belegt,
  dass `secret` **ohne** Ausschluss real erfasst wird — eine
  publication-/decoder-seitige Spaltenauswahl als Alternativerklärung fällt.
  Die einzige Variable zwischen `id=1` und `id=2` ist der `exclude_column`-Antrag;
  Mutation A bindet die Zusage zusätzlich an genau diesen Antrag.
- **„Das Row Image ist leer".** In derselben Change, deren Schlüssel-Abwesenheit
  geprüft wird, verlangt der Abschnitt `new_data->>'name' = 'ColumnAfter'`
  (nicht ausgeschlossene Spalte bleibt). Das unterscheidet „Filter wirkt gezielt"
  von „Filter wirft alles weg" — ein leeres Row Image wäre hier rot.

Der am Skript-Code gelesene Wirkort (`rowImage`,
`internal/adapters/driving/replication/mapper/mapper.go:524`) überspringt einen
Spaltenschlüssel über **dieselbe** `continue`-Verzweigung wie einen `nil`-Wert
(`values[i] == nil || containsColumn(excluded, column.Name)`) und filtert beide
Bilder (`:233`/`:237`). Die Baseline-Konstruktion prüft damit real das, wogegen
`LH-FA-CFG-005` Happy Path formuliert.

**Grenze (kein Verstoß, benannt):** Der Rundlauf übt nur den INSERT-Pfad
(`new_data`); die `old_data`-Hälfte eines UPDATE ist ausschließlich
Unit-Test-belegt (Review-F-3, INFO). Die §1/§2-Zusage dieses Slice („der Wert
steht weder in `new_data` noch in `old_data`", „historische Change bleibt über
`cdc.changes` lesbar") ist davon nicht berührt — der Wert-Check läuft real über
beide Felder, und für einen INSERT ist `old_data` null.

## 6. Urteil zur §1-Korrektur (F-1), zu F-2 und zur DoD-Häkchen-Trennung

**F-1 (§1 „künftige **und historische** Changes … den ausgeschlossenen Wert
nicht mehr tragen") — die Korrektur trifft den Befund.** Der alte Satz verlangte
für die vor dem Ausschluss erfasste Change einen Wertverlust, den derselbe Satz
über den Wirkort (Row-Image-Konstruktion) strukturell ausschließt. Die
korrigierte Fassung (`8b8fa6f`) sagt jetzt: „dann realer Beleg, dass **künftige**
Changes … den ausgeschlossenen Wert nicht mehr tragen … Eine **vor** dem
Ausschluss erfasste Change bleibt über `cdc.changes` unverändert lesbar; ein
rückwirkendes Entfernen … ist mit dem … Wirkort … nicht verbunden und deshalb
auch nicht zugesagt." Das ist **deckungsgleich mit dem Code**: die Filterung
liegt in `rowImage` zur Bauzeit, und `cdc.changes` ist eine reine Projektion
(`tools/schema/schema.yaml:303-323`, `SELECT … c.old_data, c.new_data FROM
cdc.change c`) — eine persistierte Change wird nicht umgeschrieben. Der §3-
Plan-Nachzug „Auslegung des historischen Changes" trägt dieselbe Lesart im
Detail. F-1 ist damit wirklich adressiert, nicht nur umformuliert.

**Rest-Beobachtung (Textpräzision, kein DoD-Verstoß):** Die §2-DoD-Zeile 1 formuliert
weiterhin, die historische Change sei „gemäß `LH-FA-CFG-005`s Akzeptanzkriterien
geprüft" — die Akzeptanzkriterien der Anforderung adressieren wörtlich nur
„künftige Changes von `t`". Die Zeile entschärft das selbst mit dem Klammer-Zusatz
„(Auslegung des historischen Changes in §3 Plan-Nachzug)" und delegiert an die
dort ausformulierte Lesart; die Prüfung selbst ist korrekt. Wer die DoD-Zeile
noch schärfen will, kann sie beim Closure-Nachzug an §1 angleichen — erzwungen
ist das nicht.

**F-2 (Kopf „alle drei Akzeptanzkriterien") — adressiert.** Der `Bezug:`-Kopf
nennt jetzt „Happy Path und Negative; die Boundary-Klausel ist ausdrücklich
nicht Gegenstand dieses Slice, sie verweist auf `LH-FA-SCH-003` und ist dort
belegt". Deckungsgleich mit §1 (Boundary-Ausschluss) und mit dem realen
Konvergenz-Test in `slice-067`.

**F-5 (`LH-QA-SEC-004` nicht genannt) — adressiert.** Der Sensors-Satz nennt
nun `LH-QA-SEC-004` mit dem Hinweis „dessen Messmethode ausschließlich über
`LH-FA-CFG-005` läuft" — passend zur Skript-Zusage „der Wert erreicht die
Persistenzschicht nie". F-3/F-4 bleiben INFO ohne erwartete Aktion.

**DoD-Häkchen-Trennung — korrekt.** Per `git show <commit> -- <slice-plan>`
geprüft: `7fa5784` (Implementer) hakt genau die drei Liefer-Punkte plus
`make gates`/`make test-integration` und `Doku-Update` ab; `bd8b161` (Reviewer)
hakt **genau** die Review-Zeile ab; `8b8fa6f` (Planner-Korrektur) ändert **kein**
Häkchen. Die fünf Closure-Zeilen bleiben über alle drei Commits offen. Keine
überschrittene Checkbox.

## 7. `welle-18`-Closure-Reife

Der Closure-Trigger der Welle ([`welle-18`](../plan/planning/welle-18.md) §3)
hat vier Bedingungen; ich prüfe jede einzeln gegen den gemergten Stand:

| Bedingung | Stand meiner Prüfung | Erfüllt? |
|---|---|---|
| `slice-066`, `slice-067`, `slice-068` liegen in `done/` | `slice-066`/`slice-067` real in `docs/plan/planning/done/`; **`slice-068` liegt in `in-progress/`** | **nein** |
| `make gates` grün | eigener Lauf Exit **0** | **ja** |
| Der reale E2E-Beleg aus `slice-068` läuft grün (das *Mehr* gegenüber den Einzel-DoDs) | eigener `make test-integration`-Lauf Exit **0**, beide Beleg-Zeilen real; SQL-Antrag → Live-Reload → gefiltertes Row Image am **laufenden** Feed-Container ohne Neustart | **ja** |
| Closure-Notiz in `welle-18-results.md` | `welle-18-results.md` existiert real nicht; `welle-18` trägt §7 noch als Platzhalter | **nein** |

**Ergebnis §7:** Das *Mehr* der Welle — der reale, grüne E2E-Rundlauf, den keine
Einzel-DoD für sich beansprucht — ist durch meinen eigenen Lauf belegt, und
`make gates` ist grün. Die Welle ist trotzdem **noch nicht schließbar**: es
fehlt die `slice-068`-eigene Closure (Closure-Notiz mit Steering-Loop-Eintrag,
Risiko-Ausgänge aus §6, Beobachtungs-Register-Entscheidung), der `git mv` nach
`done/`, und danach die sechsstufige Wellen-Closure (Trigger-Audit, Closure-Notiz,
Archivierung, Self-Close-Commit, Roadmap-Fortschreibung, die drei Paarungen).
Der in `welle-18` §6 geführte Dauerhaftigkeits-Vorbehalt (`ADR-0065`,
`slice-075`) ist ausdrücklich **nicht** Teil des Closure-Triggers und blockiert
ihn nicht.

## Negativbefunde

- **geprüft, ohne Befund: DoD-Punkt-Trennung und -Vollständigkeit.** Zehn Punkte,
  fünf abgehakt (drei Liefer-Punkte, Gate-Läufe, Review, Doku), fünf offen
  (Closure-Notiz, Reconciliation-Entfall, Register, Risiko-Ausgänge, Paarungen) —
  keine vorweggenommene Closure-Zeile.
- **geprüft, ohne Befund: §5-Closure-Trigger des Slice.** „DoD vollständig und
  `make test-integration` real grün und `make gates` grün und Closure-Notiz
  geschrieben" — die drei Sensor-Hälften sind grün, die Closure-Notiz ist die
  noch fehlende Planner-Leistung.
- **geprüft, ohne Befund: §3-Plan-Nachzüge (kein Scope-Creep).** Platzierung,
  Auslegung des historischen Changes, Baseline/Kontrolle und `plan.yaml` sind
  echte Nachzüge desselben Umfangs; kein vierter Liefer-Punkt, keine stille
  §3-Tabellen-Erweiterung.
- **geprüft, ohne Befund: `ADR-0059` unverändert `Accepted`.** `git diff
  --name-only d088bb0..HEAD` berührt keine Datei unter `docs/plan/adr/`;
  `make doc-immutable` (Range) 0 Befunde.
- **geprüft, ohne Befund: Abgrenzung zu `ADR-0065`.** Der Diff fügt keinen
  dauerhaften Träger hinzu; §1/§8 und `welle-18` §6 führen die Grenze als dem
  Folge-Slice zugeordnet, `ADR-0065` §Folgepflicht verweist auf `slice-075`.
- **geprüft, ohne Befund: keine ungewollte Zustands-Überschneidung.** Die
  Aktivierung von `feed_e2e_column_exclusion` wirkt nur auf diese Tabelle; kein
  anderer Abschnitt liest sie (Bestands-Grep über `tools/`, `test/`,
  `internal/`, `compose.yaml`).
- **geprüft, ohne Befund: Traceability.** `make doc-commits RANGE=d088bb0..HEAD`
  Exit 0; die drei Commits tragen `LH-FA-CFG-005`/`ADR-0059`, kein `SPEC-*`/
  `ARC-*` im Betreff.

## Summary

| Prüfung | Ergebnis |
|---|---|
| DoD-Punkte 1–5 (implementierungs-/reviewbezogen) | **erfüllt**, 1–3 selbst reproduziert |
| DoD-Punkte 6–10 (Closure) | korrekt offen, keine vorweggenommen |
| Vier Sensoren (`gates`, `test-integration`, `doc-commits`, `doc-immutable`) | **je Exit 0** |
| Zwei eigene Mutationen (Happy Path, Negative) | **zwei Mal erwarteter roter Lauf**, Rücknahme verifiziert, Blatt-Identität wiederhergestellt |
| Lage vor der Container-Ende-Grenze / kein Neustart | **bestätigt** — am Skript-Code, beide `State.Running`-Wächter grün im Lauf |
| Baseline + Kontrolle | **beide Alternativerklärungen ausgeschlossen** |
| §1-Korrektur zu F-1 | **deckungsgleich mit dem Code**; eine Textpräzision in §2 als Rest-Beobachtung |
| F-2 / F-5 | adressiert |
| DoD-Häkchen-Trennung | korrekt (drei Commits, kein Häkchen überschritten) |
| `welle-18`-Closure-Reife | **noch nicht schließbar** — `slice-068` nicht in `done/`, Closure-Notiz fehlt; das E2E-*Mehr* und `make gates` sind belegt |

**Verdikt:** `slice-068` erfüllt seine DoD. Die beiden Kernzusagen sind am
laufenden Feed-Container real reproduziert (eigener `make test-integration`-Lauf
Exit 0) und durch zwei eigene Mutationen als **einzeln** tragend belegt; `make gates`
läuft Exit 0, die Skript-Belege liegen vollständig vor der Container-Ende-Grenze
und ohne Neustart, die Baseline trennt „Wert abwesend" von „Spalte war nie im Row
Image", die Kontrolle trennt gezielten von totalem Filter, der Negative-Beleg
landet real `failed` samt Sentinel-Text. Die Plan-Korrektur `8b8fa6f` trifft F-1
und F-2 wirklich: die nachkorrigierte §1-Fassung ist deckungsgleich mit dem, was
der Code belegt.

**Übergabe an den Planner (Verifier → Planner):** Der Slice ist doD-konform;
über diese Verifikation hinaus ist **kein** Implementer-Rückweg nötig. Für den
`slice-068`-Closure-Schritt bleiben: Closure-Notiz mit Steering-Loop-Eintrag,
die zwei §6-Risiko-Ausgänge (beide bislang `<bei Closure zuzuweisen>`),
die Beobachtungs-Register-Entscheidung (Sichtung nennt zwei 1×-Treffer, keine
Neueintragung erwartet) und der `git mv` nach `done/`. Für die danach folgende
`welle-18`-Closure: mein Lauf dieses Berichts ist der repo-weite
Verifikations-Beleg aus Closure-Schritt 1 (Modul 6/8); die Welle ist erst
schließbar, wenn der `slice-068`-Closure-Lauf vorliegt. Rest-Beobachtung ohne
Zwang: die §2-DoD-Zeile 1 könnte bei der Closure an §1s Wortlaut angeglichen
werden (`LH-FA-CFG-005`s Akzeptanzkriterien adressieren wörtlich nur „künftige
Changes").
