# Review-Report: slice-harness-integration-runner-vollstaendigkeit — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan
([`slice-harness-integration-runner-vollstaendigkeit`](../plan/planning/done/slice-harness-integration-runner-vollstaendigkeit.md)),
Architect-Verdikt
[`architect-verdict-welle-routing-lese-schritt`](architect-verdict-welle-routing-lese-schritt.md) §3.2/§5,
[`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0030`](../plan/adr/0030-testpyramide.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und die Hard Rules in
[`AGENTS.md`](../../AGENTS.md) §3.1/§3.7/§3.9/§3.12/§3.13. Kein DoD-Abgleich (Verifier).

**Gegenstand:** Lauf 1: Commit `eeefa622` (Parent `af53b900`): neue Datei
`test/integration/runner_vollstaendigkeit_test.go` und Plan-Nachzug. Lauf 2 (Re-Review der
Fixrunde): Commit `631a95cd` (Parent `fa7e07a9`).

**Skill:** [`.harness/skills/reviewer.md`](../../.harness/skills/reviewer.md). **Modell:**
claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Arbeitsweise:** Der Reviewer-Lauf hat den Report selbst geschrieben (Write). Alle Mutationen
liefen einzeln an je einer eigenen `git archive HEAD`-Kopie im Scratchpad (Mutation über
`sed … > Kopie`, nie `sed -i`, nie am Echtrepo); Ausführung wie `make test` (Toolchain-Image,
`--network none`, `go test -v -count=1 -run TestRunner ./test/integration/`).
`git status --short` im Echtrepo nach den Mutationen: leer. Keine verweigerte Aktion im Lauf
(`AGENTS.md` §3.15).

---

# Lauf 1 — Review von `eeefa622`

## Eigene Messungen

**Lauf am echten Skript, `-v`:** `Runner-Vollständigkeit: 21 von 21 func TestE2E* in einem -run-Wert
von tools/harness/run-integration-tests.sh erfasst`; Tabellentest sieben Unterfälle, alle PASS.

**Mutationen, jede einzeln:**

| # | Mutation | Instanz | Farbe |
|---|---|---|---|
| M1 | `TestE2ERoutingDeleteWithoutFullReplicaIdentity` aus dem Sammelmuster (Zeile 466) entfernt | Go-Test am mutierten Skript | rot: `20 von 21`, Meldung nennt den Namen |
| M2 | wie M1, der Name steht zusätzlich nur in einer `#`-Zeile | dito | rot |
| M3 | Überspringen der `#`-Zeilen im Leser entfernt | Tabellentest | Wächter grün, `Name_nur_im_Zeilenkommentar` rot |
| M4 | wie M1, dazu `true # -run '^<Name>$'` | Go-Test am mutierten Skript | grün — Lücke, F-1 |
| M5 | Sammelmuster durch `-run 'TestE2E'` ersetzt, dazu M1 | dito | grün — F-2 |

## Findings Lauf 1 und ihr Stand nach der Fixrunde

| Finding | Kategorie | Stand nach Lauf 2 |
|---|---|---|
| F-1 Zeilenende-Kommentar (und `-run` in `echo`/Heredoc) zählt als erfasst | MEDIUM | **geschlossen** für den Zeilenende-Kommentar (`ohneKommentar`, Tabellenfall); der `echo`-/Heredoc-Teil ist als Grenze benannt und gebunden (siehe Lauf 2, Prüfpunkt 2) |
| F-2 unverankerter `-run`-Wert (`TestE2E`) erfasst alle Namen | INFO | **geschlossen**: Godoc von `fehlendeImRunner` nennt Trefferregel und Reihenfolge-Grenze, Tabellenfall bindet sie |
| F-3 Suchlauf-Feld schränkt den Suchraum ohne Grund ein | LOW | **geschlossen** (Lauf 2, Prüfpunkt 3) |
| F-4 Beschreibende Träger kennen nur die Deklarations-Hälfte (`harness/sensors/docs-check.md:99`, `docs/user/e2e-abdeckung.md:5`, `.d-check.yml:251`) | INFO | **offen** — Meldung an den Planner, Frist Closure dieses Slice; nicht Gegenstand der Fixrunde |
| F-5 Planabweichung: Schwesterdatei statt `integration_test.go` | INFO | unverändert, ohne Aktion |

---

# Lauf 2 — Re-Review der Fixrunde `631a95cd`

## Prüfpunkt 1 — F-1 geschlossen? Shell-Regel von `ohneKommentar`

Gelesen (`test/integration/runner_vollstaendigkeit_test.go:99–127`) und mit eigenen Gegenbeispielen
durchgegangen (Zeichen für Zeichen nach dem Zustandsautomaten; Schaltreihenfolge: Einfach-Zustand,
Backslash, Doppel-Zustand, Öffnen, `#`):

| Eingabe | Erwartung nach Shell | Leser |
|---|---|---|
| `'it'\''s' # c` | Kommentar ab ` # c` | korrekt (Einfach schließt, `\'` übersprungen, Einfach öffnet/schließt, dann Kommentar) |
| `"a#b" #c` | Kommentar ab ` #c` | korrekt |
| `\#` / `echo \\# c` | kein Kommentar | korrekt (Backslash überspringt; `\\#`: `#` steht nicht nach Leerraum) |
| `"a\"b # c"` | `# c` liegt im String | korrekt (Backslash-Zweig steht vor dem Doppel-Zweig, `\"` schließt nicht) |
| `$#`, `${#x}`, `a#b` | kein Kommentar | korrekt (Vorzeichen `$`/`{`/Buchstabe, kein Leerraum) |
| `'a\'b' # c` (Einfach: Backslash ist wörtlich) | Einfach schließt am `\'` | korrekt (Einfach-Zweig steht vor dem Backslash-Zweig) |

Fail-open bleiben Formen, die der Leser nicht kennt (R-3, INFO). Kein Gegenbeispiel, das einen
echten Kommentar-Rest fälschlich zählt **und** ein reales `-run` fälschlich entfernt; am echten
Skript `21 von 21`.

**Mutationen dieses Laufs, jede einzeln an eigener Kopie (Stelle → gesehene Farbe):**

| # | Mutation | Farbe |
|---|---|---|
| mA | Strip entfernt (`zeile = ohneKommentar(zeile)` → `zeile = zeile + ""`) | Wächter grün (`21 von 21`), Tabelle **rot**: `Name_nur_im_Zeilenkommentar`, `Name_nur_im_Zeilenende-Kommentar` (`fehlend = [], erwartet [TestE2EBeta]`) |
| mB | `#`-Anker entfernt (`case c == '#':`) | Tabelle **rot**: `$#_vor_dem_-run_beginnt_keinen_Kommentar` (`kein -run-Argument gefunden`) |
| mC1 | Einfach-Quote-Verfolgung entfernt (`einfach = true` → `false`) | Tabelle **rot**: `#_in_Anführungszeichen_im_Muster_bleibt_gelesen` |
| mC2 | Doppel-Quote-Verfolgung entfernt (`doppelt = true` → `false`) | **grün** — nicht gebunden, R-1 |
| mD | Backslash-Überspringen entfernt (`i++` → `_ = i`) | **grün** — nicht gebunden, R-1 |

## Prüfpunkt 2 — Godoc: benannte Grenze und Tabellenbindung

Der Godoc von `TestRunnerFuehrtJedeE2EFunktionAus` ist im Indikativ („Benannte Grenze: der Leser
parst keine Shell-Grammatik und liest ein `-run` in einem String … oder Heredoc-Text wie ein
Argument mit“) und trägt eine Kennung (`LH-QA-POR-003`); `make kommentar-kennungen DIFF=fa7e07a9`
Exit 0. Die Grenze ist durch den Tabellenfall `-run nur in echo-String zählt als erfasst (benannte
Grenze)` gebunden (Erwartung `[]`: kippt das Verhalten, wird der Fall rot). Der Heredoc-Teil der
Aussage hat keinen eigenen Fall; er läuft über denselben Pfad (Zeile ohne führendes `#` →
Argumentmuster), daher kein eigener Befund. Der Fall `unverankerter Wert erfasst alle Namen`
bindet die F-2-Aussage von `fehlendeImRunner`.

## Prüfpunkt 3 — F-3 (Suchlauf)

Die Zeilen `fa7e07a9 1 …` und `diff 1 …` suchen `(21|einundzwanzig) (func|Test|E2E)` über den
ganzen Baum (`-- .`) mit den Ausschlüssen `docs/reviews`, `docs/plan/planning/done`,
`.harness/baseline`; der Plan nennt den Grund je Ausschluss (Berichte, Records, vendored
Fremdbestand). `git grep` mit demselben Pathspec: ein Treffer außerhalb der Plan-Datei,
`harness/sensors/coverage-gate.md:204` — „21 Tests, 21 `SKIP`“ aus dem Lauf des
`postgressnapshot`-Pakets (Replication-Tests ohne DSN), **anderer Gegenstand**, kein Träger der Zahl
der `TestE2E*`-Funktionen. `make suchlauf-nachmessen PLAN=…` Exit 0, „9 Zeilen stimmen“.
**F-3 geschlossen.**

## Prüfpunkt 4 — Läufe

- Wächter am echten Skript, `-v`: `Runner-Vollständigkeit: 21 von 21 func TestE2E* in einem -run-Wert
  von tools/harness/run-integration-tests.sh erfasst` (`runner_vollstaendigkeit_test.go:64`);
  12 Unterfälle des Tabellentests PASS.
- `make test` Exit 0 (`test/integration` ok), `make fmt-check` Exit 0 (322 Dateien),
  `make kommentar-kennungen DIFF=fa7e07a9` Exit 0, `make suchlauf-nachmessen` Exit 0.
- `make docs-check`: siehe Commit-Abschnitt unten.

## Findings Lauf 2

### R-1 — Zwei Zweige von `ohneKommentar` sind zugesagt, aber nicht an der Eingabeseite gebunden

- `kategorie`: MEDIUM
- `quelle`: Maintainability; Skill „Zusage ohne Bindung an ihre Eingabeseite“ (die Klasse steht
  im Skill unter HIGH; hier trägt die Tabelle drei der fünf zugesagten Zweige, die Lücke ist auf
  zwei Zweige eines Hilfsparsers begrenzt, deshalb MEDIUM — Einstufung bei Widerspruch Architect)
- `pfad`: `test/integration/runner_vollstaendigkeit_test.go:35–36` und `:99–102` (Godoc: „außerhalb
  von Einfach- und Doppelanführungszeichen“, „Ein Backslash … schützt das folgende Zeichen“),
  Zweige `:112–113`, `:120–121`
- `befund`: mC2 (Doppel-Quote-Verfolgung entfernt) und mD (Backslash-Überspringen entfernt) färben
  weder den Wächter noch einen der zwölf Tabellenfälle rot; der Fall mit `#` in Anführungszeichen
  nutzt nur Einfachanführungszeichen, und kein Fall enthält `\#` oder `\"`. Die Godoc-Zusagen für
  Doppelanführungszeichen und Backslash sind damit grün ohne Aussage.
- `verifizierbar`: ja — mC2/mD, reproduzierbar an einer Kopie
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### R-2 — Godoc und Name von `TestRunnerLeserDreiZustaende` sagen „drei Zustände“, die Tabelle hat zwölf Fälle

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Zusage-Kommentar beschreibt, was da ist)
- `pfad`: `test/integration/runner_vollstaendigkeit_test.go:181–184`
- `befund`: Der Godoc zählt drei Zustände (vollständig, Name entfernt, Name nur im Kommentar) und
  „prüft jeweils den erwarteten Namen“; die Tabelle trägt zwölf Fälle, darunter Fehlerfälle
  (`-run` ohne Anführungszeichen, kein `-run`), bei denen der Fehlertext geprüft wird, nicht ein
  Name. Der Satz stimmt nur für einen Teil der Fälle; der Funktionsname trägt dieselbe Zahl.
  Keine semantische Wirkung, aber eine Zusage, die der Code nicht mehr trägt.
- `verifizierbar`: ja — Godoc gegen Tabelle gelesen
- `klasse`: Zusage-Kommentar veraltet

### R-3 — Formen außerhalb des Zustandsautomaten bleiben fail-open

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `test/integration/runner_vollstaendigkeit_test.go:103–127`
- `befund`: Ein `#` direkt nach `;`, `&`, `(` oder `|` (`true;# -run '^X$'`), ANSI-C-Quotes
  (`$'…\'…'`), Quote-Zustand über Zeilengrenzen und der Tabulator-Zweig des Ankers sind weder
  gelesen noch gebunden; jedes führt nur dazu, dass ein Kommentar-Rest als `-run` mitgelesen wird
  (fail-open, wie die benannte Grenze), nie dazu, dass ein reales `-run` entfällt. Keine Aktion
  erwartet.
- `verifizierbar`: nein — nur Lektüre
- `klasse`: Grenze des Sensors

## Verdikt

Lauf 1: F-1, F-2, F-3 sind mit der Fixrunde geschlossen (mA/mB/mC1 gesehen rot am Tabellentest,
echtes Skript `21 von 21`, Suchlauf Exit 0). Neu: **ein offenes MEDIUM (R-1)**, ein LOW (R-2), ein
INFO (R-3); F-4 bleibt INFO-Meldung. **Merge-blockierend: nein** (der Wächter trägt seine
Hauptzusage; R-1 ist eine Bindungslücke an zwei Hilfszweigen, das Verhalten am echten Skript ist
richtig) — **eine zweite, kleine Fixrunde ist nötig** (zwei Tabellenfälle für R-1, Godoc-Satz für
R-2). Wegen des offenen MEDIUM bleibt die DoD-Zeile „Review durchgeführt“ im Plan **offen**
(Skill §DoD-Checkbox-Nachzug); sie wird regulär bei Schritt 21 des Implementer-Ablaufs gezogen.

## Geprüft, ohne Befund

- `test/integration/` (`runner_vollstaendigkeit_test.go`: Kommentare, höchstens eine Kennung je
  Block, `make kommentar-kennungen DIFF=fa7e07a9` leer, `make fmt-check` sauber): Befunde nur R-1
  bis R-3.
- `tools/harness/run-integration-tests.sh` (nur gelesen, nicht im Diff): ohne Befund.
- Docker-only (`AGENTS.md` §3.1): der Test liest nur Repo-Dateien; ohne Befund.
- `docs/plan/planning/in-progress/` (Plan-Nachzug der Fixrunde, Suchlauf-Block): Zahlen und
  Ausschluss-Gründe stimmen; ohne Befund.
