# Review-Report: slice-harness-integration-runner-vollstaendigkeit — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan
([`slice-harness-integration-runner-vollstaendigkeit`](../plan/planning/in-progress/slice-harness-integration-runner-vollstaendigkeit.md)),
Architect-Verdikt
[`architect-verdict-welle-routing-lese-schritt`](architect-verdict-welle-routing-lese-schritt.md) §3.2/§5,
[`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0030`](../plan/adr/0030-testpyramide.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und die Hard Rules in
[`AGENTS.md`](../../AGENTS.md) §3.1/§3.7/§3.9/§3.12/§3.13. Kein DoD-Abgleich (Verifier).

**Gegenstand:** Commit `eeefa622` (Parent `af53b900`): neue Datei
`test/integration/runner_vollstaendigkeit_test.go` (211 Zeilen) und Plan-Nachzug (Zeile §3 der
Dateitabelle, Suchlauf-Feld).

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md). **Modell:**
claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage / Arbeitsweise:** Der Reviewer-Lauf hat den Report selbst geschrieben (Write). Alle
Mutationen liefen einzeln an je einer eigenen `git archive HEAD`-Kopie im Scratchpad (Mutation
über `sed … > Kopie`, nie `sed -i`, nie am Echtrepo); Ausführung wie `make test` (Toolchain-Image
`TOOLCHAIN_IMAGE`, `--network none`, `go test -v -count=1 -run TestRunner ./test/integration/`).
`git status --short` im Echtrepo nach den Mutationen: leer. Keine verweigerte Aktion im Lauf
(`AGENTS.md` §3.15); ein `rm -rf` mit Shell-Variable wurde von der Sicherheitsprüfung verweigert,
das Ziel (frische Kopien) war ohne Löschen erreichbar und wurde mit neuen Verzeichnisnamen
erreicht.

---

## Eigene Messungen

**Lauf am echten Skript, `-v` (nach der letzten Edit-Änderung der Fehlermeldung):**
`runner_vollstaendigkeit_test.go:61: Runner-Vollständigkeit: 21 von 21 func TestE2E* in einem -run-Wert von tools/harness/run-integration-tests.sh erfasst`
— `--- PASS: TestRunnerFuehrtJedeE2EFunktionAus`; `TestRunnerLeserDreiZustaende` sieben
Unterfälle, alle PASS.

**Mutationen, jede einzeln (Stelle → gesehene Farbe):**

| # | Mutation | Instanz | Farbe |
|---|---|---|---|
| M1 | `TestE2ERoutingDeleteWithoutFullReplicaIdentity` aus dem Sammelmuster (Zeile 466) entfernt | Go-Test am mutierten Skript | **rot**: `20 von 21 … erfasst`, Meldung nennt genau den Namen |
| M2 | wie M1, der Name steht zusätzlich nur in einer `#`-Zeile | dito | **rot**, derselbe Name |
| M3 | Überspringen der `#`-Zeilen im Leser (`strings.HasPrefix(…, "#")` → `false`), Skript wie M1 plus `# go test -run '^<Name>$' ./x` | Tabellentest | Wächter-Lauf grün (`21 von 21`), **`TestRunnerLeserDreiZustaende/Name_nur_im_Zeilenkommentar` rot** (`fehlend = [], erwartet [TestE2EBeta]`); die übrigen sechs Fälle grün |
| M4 | wie M1, dazu eine Zeile `true # -run '^<Name>$'` (Zeilenende-Kommentar) | Go-Test am mutierten Skript | **grün**, `21 von 21` — Lücke, siehe F-1 |
| M5 | Sammelmuster durch `-run 'TestE2E'` ersetzt, dazu M1 | dito | **grün**, `21 von 21` — siehe F-2 |

Eingabe-Bindung (b): Die Negativfälle des Tabellentests vergleichen mit `reflect.DeepEqual` die
Liste der fehlenden Namen (`[]string{"TestE2EBeta"}`), nicht nur „rot“; M3 belegt, dass der
Kommentarfall an der Eingabeseite (Skript-Text) rot wird. Gebunden.

**Weitere Sonden am Leser (gelesen, nicht gefahren):** `-run` mit Wert auf der Folgezeile
(Zeilenfortsetzung `\` nach `-run`) fällt in die Gruppe 3 und endet laut mit „ohne
Anführungszeichen“ (fail-closed). `-run "$VAR"` kompiliert als Muster `$VAR`, trifft nichts und
färbt rot (fail-closed). Ein `-run` in `echo "…-run '^X$'…"` oder in einem Heredoc-Text wird wie
ein Argument gelesen (fail-open, wie M4).

**Träger:** `make test` (Race-Detector, `./...`) Exit 0, `test/integration` `ok`;
`make fmt-check` Exit 0; `make kommentar-kennungen DIFF=af53b900` Exit 0, kein Kandidat;
`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-integration-runner-vollstaendigkeit.md`
Exit 0, „9 Zeilen stimmen“; `make docs-check` siehe Commit-Abschnitt unten.

**(g) Lauf unter `make test-integration`:** Der Runner `tools/harness/run-integration-tests.sh`
fährt über `-run`-Werte nur `TestE2E*` und `TestAbdeckungstabelleZeilen`; die beiden
`TestRunner*`-Funktionen werden dort **nicht** mitgefahren, laufen aber unter `make test`
(`./...`, netzlos, ohne Stack — gemessen). Kein Selbstbezug: der Namen-Leser zählt nur
`TestE2E*` (`strings.HasPrefix(…, "TestE2E")`), und `TestAbdeckungstabelleZeilen`
(`integration_test.go:1465`) filtert auf dasselbe Präfix, die Abdeckungstabelle bleibt unberührt.

---

## Findings

### F-1 — Ein Zeilenende-Kommentar (und ein `-run` in `echo`/Heredoc-Text) zählt als erfasst; der Tabellentest bindet diese Grenze nicht

- `kategorie`: MEDIUM
- `quelle`: Plan §2 DoD-Zeile 1 („Kommentare und Namen außerhalb eines `-run`-Arguments zählen
  nicht als erfasst“); `LH-QA-POR-003`
- `pfad`: `test/integration/runner_vollstaendigkeit_test.go:99–118` (`runMuster`, Skip nur für
  Zeilen mit führendem `#`), Godoc `:36–37`
- `befund`: M4 zeigt: `true # -run '^<Name>$'` färbt den Wächter nicht rot, obwohl die
  Funktion keinen realen `go test`-Aufruf trifft. Der Godoc benennt die Grenze ehrlich („Ein
  Kommentar am Zeilenende … wird mitgelesen“), die DoD-Zeile sagt aber ohne Einschränkung, dass
  Kommentare nicht zählen; der Tabellentest hat weder einen Fall, der die Grenze festhält, noch
  einen, der sie schließt. Dieselbe Lücke gilt für ein `-run` in einem `echo`-String oder
  Heredoc-Text.
- `verifizierbar`: ja — M4, reproduzierbar an einer Kopie
- `klasse`: Zusage weiter als ihre Messung (Wächter fail-open an einer benannten, nicht
  gebundenen Stelle)

### F-2 — Ein unverankerter `-run`-Wert (`TestE2E`) gilt als Erfassung aller Namen

- `kategorie`: INFO
- `quelle`: Plan §1 (Zuschnitt, Trefferregel wie `go test -run`)
- `pfad`: `test/integration/runner_vollstaendigkeit_test.go:120–123` (Godoc von
  `fehlendeImRunner`), `:132–139`
- `befund`: M5: ein einziges `-run 'TestE2E'` macht den Wächter grün für alle 21 Namen. Das
  entspricht der Semantik von `go test -run` (ein solcher Wert fährt tatsächlich alle), ist also
  keine Falschaussage über „läuft nie“; der Wächter sagt aber nichts über Reihenfolge (die
  Container-beendenden Funktionen müssen laut Runner-Kommentar zuletzt laufen) und der Godoc
  nennt diese Grenze nicht ausdrücklich. Auch ein `-run`-Wert mit `/` (Untertest-Trennung bei
  `go test`) wird als ganzes Muster gelesen und färbt konservativ rot. Keine Aktion erwartet;
  der Architect kann die Grenze im Godoc benennen lassen.
- `verifizierbar`: ja — M5
- `klasse`: Grenze des Sensors (Anwesenheit im Muster, nicht Ausführungsreihenfolge)

### F-3 — Suchlauf-Feld schränkt den Suchraum der Zahl-Zeile auf zwei Pfade ein, ohne den Grund zu nennen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 „Suchform“ (Suchraum ist der ganze Baum; jede Einschränkung steht
  mit Grund im Feld)
- `pfad`: Plan §3, Suchlauf-Block, Zeilen `af53b900 0 …` und `diff 0 …`
  (`-- harness/README.md docs/user`)
- `befund`: Das Muster `(21|einundzwanzig) (func|Test|E2E)` wird nur in zwei Pfaden gesucht;
  der Text nennt die Pfade als Träger-Kandidaten, gibt aber keinen Grund für den Ausschluss des
  übrigen Baums (z. B. `harness/sensors/`, `.d-check.yml`). Die Zahlen selbst stimmen
  (`make suchlauf-nachmessen`: 9 Zeilen).
- `verifizierbar`: ja — `git grep` ohne Pathspec an beiden Ständen
- `klasse`: Suchform §3.13 unvollständig (Einschränkung ohne Grund)

### F-4 — Beschreibende Träger kennen nur die Deklarations-Hälfte

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger in fremder Datei wird gemeldet, nicht still mitgeändert)
- `pfad`: `harness/sensors/docs-check.md:99` („Dieselbe Klasse führt
  `BEO-PGC/test-runner-stiller-ausschluss` für die `-run`-Muster“), `docs/user/e2e-abdeckung.md:5`,
  `.d-check.yml:251` (nennen nur `TestAbdeckungstabelleZeilen` als tragende Garantie)
- `befund`: Mit dem neuen Wächter ist die dort als offen beschriebene `-run`-Hälfte getragen
  (im Rahmen von F-1/F-2). Der Plan führt den Register-Zielort „verkörpert“; die drei
  Fundstellen sind Meldungen an den Planner mit Frist Closure dieses Slice, nicht Diff-Fehler.
- `verifizierbar`: ja — `git grep -n 'test-runner-stiller-ausschluss' -- harness .d-check.yml docs/user`
- `klasse`: Träger-Nachzug (Meldung)

### F-5 — Planabweichung: Schwesterdatei statt `integration_test.go`

- `kategorie`: INFO
- `quelle`: Plan §3 Dateitabelle
- `pfad`: `test/integration/runner_vollstaendigkeit_test.go`
- `befund`: Der Plan ließ „oder eine Schwesterdatei“ ausdrücklich zu; die §3-Zeile ist
  nachgezogen („geliefert: Schwesterdatei … neu (statt update …)“). Akzeptabel; die zweite
  Plan-Zeile („Testdatei im selben Paket, neu / update“) deckt nun dieselbe Datei und bleibt als
  Doppelung stehen. Der Test-Name `TestRunnerLeserDreiZustaende` trägt sieben Fälle (drei
  Zustände plus vier Eingabe-/Grenzfälle) — nur Benennung, keine Wirkung.
- `verifizierbar`: ja — Diff
- `klasse`: Planabweichung innerhalb der Plan-Erlaubnis

---

## Architect-Fragen

1. F-1: Soll die Lücke (Zeilenende-Kommentar, `-run` im String/Heredoc) **geschlossen** werden
   (Leser erkennt Kommentar-/String-Kontext) oder als Grenze **benannt und per Tabellenfall
   gebunden** werden, mit entsprechender Umformulierung der DoD-Zeile „Kommentare … zählen nicht
   als erfasst“? Beides ist ein Zug am Plan/DoD, kein Reviewer-Vorschlag im Code.
2. F-2: Reicht die Benennung der Trefferregel „Anwesenheit im Muster“ im Godoc, oder soll der
   Wächter einen unverankerten Wert wie `TestE2E` ablehnen?

## Verdikt

Kein HIGH, ein offenes MEDIUM (F-1). **Merge-blockierend: nein** (der Wächter trägt seine
Hauptzusage: M1/M2 rot, M3 am Tabellentest rot, 21 von 21 am echten Skript, Godoc benennt die
Grenze); die Entscheidung zu F-1 braucht eine Antwort des Architects. Wegen des offenen MEDIUM
bleibt die DoD-Zeile „Review durchgeführt“ im Plan **offen** (Skill §DoD-Checkbox-Nachzug: nur
ohne offenes HIGH/MEDIUM).

## Geprüft, ohne Befund

- `test/integration/` (`runner_vollstaendigkeit_test.go`: Godoc/Kommentare, höchstens eine
  Kennung je Block, `make kommentar-kennungen` leer; Fehlermeldungen nennen Datei und Namen;
  `make fmt-check` sauber): ohne weiteren Befund.
- `tools/harness/run-integration-tests.sh` (nur gelesen, nicht im Diff): `-run`-Werte Zeilen
  466, 3842, 5106, 5203, 5244 gelesen; kein Befund.
- Docker-only (`AGENTS.md` §3.1): der Test liest nur Repo-Dateien, keine Host-Werkzeuge; ohne
  Befund.
- `docs/plan/planning/in-progress/` (Plan-Nachzug): Suchlauf-Zahlen stimmen, Befund nur F-3/F-5.
