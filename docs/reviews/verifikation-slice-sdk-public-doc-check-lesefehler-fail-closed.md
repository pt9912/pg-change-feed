# Verifikations-Report: Slice sdk-public-doc-check-lesefehler-fail-closed ([ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)) — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD- und Entscheidungs-Konformität,
Plan-vs-Code-Diff, eigene Sensor-Läufe, Reproduktion des Defekts und Einzelmutationen, in frischem Kontext.

**Gegenstand:** `git diff 8e38389a HEAD` — zwei Commits: `8f6430a0` (Skript, Tabellentest,
Sensor-Vertrag, Plan-Feld) und `894fb7ce` (Review-Report).

**Eingang:**
[Plan](../plan/planning/done/slice-sdk-public-doc-check-lesefehler-fail-closed.md),
[ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) (`Accepted`, unberührbar),
[ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) (Schwester-Gate, Muster),
[Review-Report](review-slice-sdk-public-doc-check-lesefehler-fail-closed.md) (0 HIGH/MEDIUM/LOW, 3 INFO),
[Befund-Quelle](verifikation-slice-handbuch-public-doc-check-gate-und-skill.md) §7 Punkt 2.
Bezug: [LH-FA-SST-009](../../spec/lastenheft.md).

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Gedruckte Zeile aus meinem Lauf |
|---|---|---|
| `make test-sdk-public-doc-check` | Exit 0 | `run-sdk-public-doc-check-tests: alle 19 Fälle bestanden` (uid 1000, Lesefehler-Fälle gelaufen) |
| `make sdk-public-doc-check` | Exit 0 | `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `make gates` (in Log-Datei, Exit direkt gesichert) | Exit 0 | `baseline-verify: v6.13.0 OK`, `total: … 82.0%`, `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `a-check … gesamt: 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)` |
| `make test` | Exit 0 | alle Pakete `ok` |
| `make fmt-check` | Exit 0 | `323 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=8e38389a` | Exit 0 | keine Kandidaten |
| `make suchlauf-nachmessen PLAN=…` | Exit 0 | `suchlauf-nachmessen: 14 Zeilen stimmen` |
| `make docs-check` | Exit 0 | `d-check: 1590 Datei(en) geprüft, 0 Befund(e)` vor Anlage dieses Reports; `1591 Datei(en) geprüft, 0 Befund(e)` nach Anlage (Abschnitt 8) |
| `make commit-traceability` | Exit 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=8e38389a..HEAD` | Exit 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=8e38389a..HEAD` | Exit 0 | `0 Befund(e)` |
| `make -n sdk-pack-python` / `-csharp` / `-kotlin` | — | jeweils `bash tools/harness/sdk-public-doc-check.sh` als erste Zeile, dann das Pack-Skript |
| `make sdk-pack-python` | Exit 0 | erste gedruckte Zeile: `sdk-public-doc-check: keine interne Kennung unter sdks`; danach der Docker-Bau |
| `git status --short` (Echtrepo, nach allen Läufen und Mutationen) | leer | Mutationen nur auf Kopien im Scratchpad |

## 2. Reproduktion des Defekts (Scratchpad, uid 1000)

Altes Skript = `git show 8e38389a:tools/harness/sdk-public-doc-check.sh`, neues = Arbeitsbaum.

| Fall | Altes Skript | Neues Skript |
|---|---|---|
| `a/README.md` mit Kennung, `chmod 000` | Exit 0, „keine interne Kennung“ (grep-Meldung nur auf stderr) | Exit 2, `Lesefehler: <pfad>/a/README.md (grep: … Keine Berechtigung)` |
| nicht lesbares Unterverzeichnis | Exit 0, „keine interne Kennung“ | Exit 2, `Lesefehler beim Durchsuchen von <wurzel> (find: … Keine Berechtigung)` |
| nicht existierende Wurzel | Exit 0, „keine interne Kennung“ | Exit 2, `Lesefehler beim Durchsuchen …` |
| Pfad mit Leerzeichen (`a b/my file.md`) mit Kennung | nicht gefahren | Exit 1, Treffer mit vollem Pfad |
| NUL-Byte-Datei mit Kennung | nicht gefahren | Exit 1, Treffer `d.bin:2:SPEC-022` |
| saubere Datei | nicht gefahren | Exit 0 |

Der Defekt ist am Parent reproduziert (gemessen) und am HEAD behoben.

## 3. DoD Zeile für Zeile (keine Häkchen gesetzt)

| DoD-Zeile | Befund | Beleg |
|---|---|---|
| Liefer-Punkt 1 Skript: Lesefehler Exit 2, Meldung stderr, Exit 0/1 unverändert | erfüllt | Abschnitt 2; `git grep -n -F '\|\| true' -- tools/harness/sdk-public-doc-check.sh` leer (Exit 1); `make sdk-public-doc-check` Exit 0; Mutation M7 (`\|\| grc=$?` → `\|\| true`) rot (Abschnitt 4) |
| Liefer-Punkt 2 Tabellentest: Datei, Unterverzeichnis, NUL-Byte, Meldungsbindung, nur uid ≠ 0, Überspringen gemeldet | erfüllt | 19 Fälle grün; Skript-Zweig meldet das Überspringen unter root; je Fall eine rote Mutation (Abschnitt 4) |
| Liefer-Punkt 3 Sensor-Vertrag: Exit 2, Grenze der uid-0-Fälle | erfüllt | `git diff` der Datei: Ausgangstabelle Zeile 2, Grenze 4; `Lesefehler` 4 Treffer |
| Nur diese Pfade (`docs/plan/adr`, `.github`, `AGENTS.md`) | erfüllt | `git diff --name-only 8e38389a HEAD -- docs/plan/adr .github AGENTS.md harness/README.md` ist leer; Gesamt-Diff `git diff --name-only 8e38389a HEAD`: fünf Dateien (Plan, Review-Report, Sensor-Vertrag, Tabellentest, Skript) |
| `make gates` grün, Exit gesondert | erfüllt | Exit 0, ungepiped in Log-Datei |
| Review durchgeführt, kein offenes HIGH/MEDIUM | erfüllt (bereits `[x]`) | Review-Report: 0/0/0, 3 INFO |
| §3.13-Suchlauf, Gefundenes und Nichtgefundenes, Nachmessen Exit 0 | erfüllt | Feld in §3 des Plans; `14 Zeilen stimmen` |
| Doku-Update (`harness/README.md` nur falls Ausgänge aufgezählt) | erfüllt, kein Eingriff nötig | die Zeile `make sdk-public-doc-check` zählt keine Exit-Codes auf (gelesen) |
| Closure-Notiz, Register, §6-Ausgänge, drei Paarungen | offen — Planner-Aufgabe bei der Closure | Plan §6/§7 tragen Platzhalter |

## 4. Einzelmutationen auf Kopien (Stelle, Instanz, Farbe — [`AGENTS.md`](../../AGENTS.md) §3.12)

Instanz: `git archive HEAD` der beiden Dateien in das Scratchpad, Mutation per `awk … > Kopie` (kein `sed -i`),
Tabellentest an der Kopie. Ausgangslage grün (19 Fälle).

| # | Mutation (Stelle) | Farbe | Rote Fälle |
|---|---|---|---|
| M1 | grep-Lesefehler-Zweig (`if [ "$grc" -ge 2 ]`) → `if false` | rot | „nicht lesbare Datei“ (Exit 0 statt 2; Meldung) |
| M2 | find-Zweig (`if [ "$frc" -ne 0 ]`) → `if false` | rot | „nicht lesbares Unterverzeichnis“ (Exit und Meldung) |
| M3 | `grep -aHnE` → `grep -IHnE` | rot | „Kennung in Datei mit NUL-Byte“ (Exit 0 statt 1) |
| M4 | Prune `obj` entfernt | rot | „Kennung in einer Bau-Ausgabe (obj)“ |
| M5 | Prune `grpc_gen` entfernt | rot | „… erzeugten Python-Code (grpc_gen)“ |
| M6 | Prune `dist` entfernt | rot | „… Bau-Ausgabe (dist)“ |
| M7 | `\|\| grc=$?` → `\|\| true` (Plan-Vorgabe) | rot | „nicht lesbare Datei“ (Exit und Meldung) |
| M8 | Prune `build`, `.gradle`, `__pycache__`, `.pytest_cache`, `bin` je einzeln entfernt | **grün** (alle 19 bestanden) | kein Fall deckt diese fünf Prune-Einträge |

M8 ist keine Slice-Abweichung: die DoD verlangt „eine Prune-Zeile“ rot; Muster und Prune-Liste sind ausdrücklich
nicht Gegenstand (Plan §1). Es ist eine Altlücke des Tabellentests (er kannte schon vor dem Slice nur `obj`, `dist`,
`grpc_gen`; `*.egg-info` nicht mutiert).

## 5. Entscheidungs-Konformität

- [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md): Muster, Prune-Liste, Dateiauswahl und
  Gate-Einbindung unverändert (Suchlauf-Zeile 7 grün: `fail-closed|fail-open` in der ADR weiter 1 Zeile). Die
  Lesefehler-Semantik ist eine **Verschärfung** (keine §3.6-Senkung); der Hauptlauf-Entscheid „keine Folge-ADR,
  Sensor-Vertrag trägt sie“ ist nachvollziehbar, aber ein Vermerk: die `Accepted` ADR enthält die Exit-2-Semantik
  nicht, der Wortlaut liegt allein im Sensor-Vertrag (wie beim Handbuch-Gate,
  [ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md), dessen Verifikation das als
  „Deckung im Zweck, ohne Wortlaut in der ADR“ führte). Die Frage einer gemeinsamen Folge-ADR für beide Gates bleibt
  beim Planner.
- Muster des Schwester-Skripts (Scratch-Verzeichnis mit Trap, `find` in Datei, `grep -a` mit Exit-Auswertung):
  gelesen, übernommen.
- Code-Lesung des Skripts: `set -euo pipefail`; `|| frc=$?`/`|| grc=$?` vermeiden den `set -e`-Abbruch;
  `read -d ''` mit `-print0` und `grep … -- "$f"` tragen Leerzeichen-Pfade (gemessen, Abschnitt 2).
  Kosmetik (INFO): bei Treffern endet die Treffer-Ausgabe auf einer Leerzeile (`printf '%s\n' "$hits"` bei bereits
  `\n`-terminiertem Sammelwert, gemessen); kein Aufrufer parst sie.

## 6. Post-Push-CI ([`AGENTS.md`](../../AGENTS.md) §3.10)

Kein Workflow berührt; kein §3.10-Bedarf. Das Gate läuft unverändert in `make gates` (CI-Schritt „Gates“). Der
Tabellentest ist kein Gate-Bestandteil, die root-Grenze (Risiko 1) berührt CI daher nicht.

## 7. Für den Planner (gemeldet, nicht geändert)

1. **Register `intern-kennungen-in-ausgelieferten-texten`:** Der Eintrag nennt unter „Fangnetz für `docs/user/`“ die
   „offene Grenze des SDK-Wächters“ (Exit 0 bei Lesefehler) und als Träger der Behebung den Slice „in `open/`“.
   Beides ist am HEAD überholt: die Grenze ist am Parent reproduziert und am HEAD geschlossen (Abschnitt 2), der Plan
   liegt in `in-progress/` und wandert bei Closure nach `done/` (Zustandsfeld mit auflösbarem Anker,
   [`AGENTS.md`](../../AGENTS.md) §3.7). Zähler bleibt **1×** (kein neuer Vorgang der Klasse). `observation.md` und
   `state.md` habe ich nicht geändert.
2. **Risiko-Ausgänge (Vorschlag für §6):** (a) *Tabellentest unter uid 0* — weiter offen als dokumentierte Grenze
   (Sensor-Vertrag Grenze 4; das Überspringen wird gemeldet; das Gate selbst ist nicht betroffen). (b) *Falsch rot durch
   nicht lesbares Verzeichnis* — entfallen am echten Baum: `make sdk-public-doc-check` Exit 0, `make sdk-pack-python`
   Exit 0 (Gate als erste Zeile); die Richtung bleibt die sichere. (c) *Architect-Antwort bewegt den Umfang* —
   entfallen: Hauptlauf-Entscheid Verschärfung ohne Folge-ADR, Umfang unverändert.
3. **Folge-ADR-Frage** (Abschnitt 5) für beide Gates gemeinsam: Entscheidung des Planners/Architects, keine Bedingung
   dieses Slice.
4. **Altlücke M8** (fünf Prune-Einträge ohne Tabellenfall): optionaler Folge-Kandidat, kein Closure-Blocker.
5. **Closure-Pflichten:** Plan §7, §6-Ausgänge, übrige DoD-Häkchen, Register (Punkt 1) und der `git mv` nach `done/`
   (erst Inhalt, dann reiner Move, [`AGENTS.md`](../../AGENTS.md) §3.3). Ich habe keine Häkchen gesetzt.

## 8. Verdikt

**bestanden.** Alle Liefer-Punkte sind am Diff und an meinen Läufen belegt; der Defekt ist am Parent reproduziert
und am HEAD als Exit 2 mit Dateiname behoben; sieben Einzelmutationen färben den Tabellentest an der richtigen Stelle
rot; `make gates` und alle genannten Sensoren sind Exit 0; der Diff berührt weder ADR, `.github`, `AGENTS.md` noch
SDK-Code oder Versionen.

**Bedingungen:** keine technische. Ein weiterer Review-Durchgang ist nicht nötig. Vor dem Schluss (Planner): Punkte 1
und 5 aus Abschnitt 7.

`make docs-check` nach Anlage dieses Reports (vor dem Commit): Exit 0, `d-check: 1591 Datei(en) geprüft, 0 Befund(e)`.
