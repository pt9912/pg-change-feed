# Review-Report: slice-harness-grep-pipe-sigpipe-unter-pipefail — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `harness-grep-pipe-sigpipe-unter-pipefail`, Diff `git diff 3bcafd41 f8aca691` (Implementer-Commit `f8aca691`; der dazwischenliegende Commit `7d1c4611` (Handbuch) ist fremd;
seine `docs/user/benutzerhandbuch.md`-Zeilen im Diffstat sind fremd und ausgeschlossen). Geprüft: sechs Skripte unter
`tools/harness/`, `docs/user/e2e-abdeckung.md`, der Plan.

**Skill:** `.harness/skills/reviewer.md` @ `f8aca691`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage und Arbeitsweise:** Report mit dem Write-Werkzeug geschrieben (kein Edit-Werkzeug im Lauf); die DoD-Zeile im Plan wurde über eine Kopie im Scratchpad
(`sed … > Kopie`, kein `-i`) und `cp` nachgezogen. Das Wegwerf-Skript der Reproduktion liegt im Scratchpad, nicht im Repo. Keine Verweigerung der
Berechtigungsschicht im Lauf (`AGENTS.md` §3.15 nicht ausgelöst). `make image` nicht aufgerufen (kein Server-Code im Diff); `:dev` lag als
`ghcr.io/pt9912/pg-change-feed:dev` geladen vor.

**Eingangs-Kontext:**

- Slice-Plan `harness-grep-pipe-sigpipe-unter-pipefail` (§1, §2 DoD, §3 Tabelle, Suchlauf-Feld und Befunde)
- [`ADR-0030`](../plan/adr/0030-testpyramide.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-QA-POR-003`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.7, §3.9, §3.12, §3.13, §3.15), [`harness/conventions.md`](../../harness/conventions.md)
- Herkunft: `review-slice-sdk-sse-filter-phase-verbindung-haertung` F-1; Register `BEO-PGC/runner-grep-pipe-verfehlt-zeile`

**Eigene Messungen** (Exit-Codes je als eigener Schritt ausgewertet):

- **Reproduktion (eigenes Skript, Docker 29.8.2 build 7fc2dff):** beendeter `busybox`-Container mit Treffer-Zeile `REJECTED code=401` und 20 Folgezeilen,
  `set -euo pipefail`, je 1000 Aufrufe in einer `if`-Schleife. Gedruckt: `form=q calls=1000 failures=16` (alte Form `docker logs … | grep -qF`) und
  `form=nq calls=1000 failures=0` (neue Form `… | grep -F … >/dev/null`). Der Implementer nennt 10 bzw. 9 von 1000 vor dem Fix; meine Rate (1,6 %) liegt in
  derselben Größenordnung, die Richtung ist eindeutig. Die Rate ist lauf- und lastabhängig (Zufallsprozess).
- **Zeile für Zeile des Diffs** (`git diff --word-diff`): alle 35 Ersetzungen (21 `run-integration-tests.sh`, je 3 in den drei `run-sdk-*`, 3 in `lib-sdk-filter-fixture.sh`,
  2 in `lib-sdk-route-fixture.sh`) haben dieselbe Form: aus `grep -qF|-qE M` wird `grep -F|-E M >/dev/null`. Die Regex-Wahl blieb in allen 35 erhalten (die drei
  `-E`-Stellen in den Fixtures `SEEN`/`SEEN_SECOND` mit Zeilenanker bleiben `-E`; alle anderen waren `-F` und bleiben `-F`), Muster und Quoting (`"$reject_marker"`,
  `"REJECTED code=Unauthenticated"`) byte-gleich, `if`/`until`-Kontext unverändert, keine `!`-Negation unter den 35 (die zwei Negativ-Prüfungen stehen als
  `if … ; then <Fehler> exit 1` ohne `!`). Exit-Status: Treffer 0 / kein Treffer 1 in beiden Formen; unter `set -e` im `if`-Kopf ohne Wirkung; `>/dev/null` ändert nur die
  Ausgabe. Der Pipeline-Status unter `pipefail` ist jetzt der von `grep` (0/1), solange `docker logs` selbst mit 0 endet.
- `bash -n` über die sechs Skripte: ok. `make fmt-check` Exit 0. `make kommentar-kennungen DIFF=3bcafd41` Exit 0, kein Kandidat.
  `make suchlauf-nachmessen PLAN=<Plan>` Exit 0, „14 Zeilen stimmen“. `make docs-check` Exit 0 (1569 Dateien, 0 Befunde) vor Anlage dieses Reports.
  `make test` Exit 0.
- `docs/user/e2e-abdeckung.md`: 53 Zeilen geändert, ausschließlich der Zeilen-Lokator `run-integration-tests.sh:<N>`, jeweils um genau +4 (nachgemessen über
  Paarbildung alt/neu: 53 Differenzen, alle 4; keine weitere Textänderung, keine `.go`-Lokatoren berührt). Das entspricht den vier neuen Kommentarzeilen im Kopf des Runners.
- Tier-Lauf unmutiert: `make test-sdk-python-integration` Exit 0 (Schlusszeilen: Regel-, Routing- und Filter-Belege grün, gedruckt u. a.
  `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`). `make test-integration` habe ich **nicht** gefahren (Aufwand); die 21 Stellen in
  `run-integration-tests.sh` tragen daher nur die Zeilenprüfung und die Reproduktion, nicht einen eigenen Gesamtlauf. Der Lauf des Implementers
  („Exit 0, 21 Go-Zeilen und 53 Bash-Zeilen“) ist **übernommen**.
- Mutationsprobe: keine zusätzliche Kopie-Mutation; die Reproduktion der alten Form (16 von 1000) gegen die neue (0 von 1000) in einem Skript ist die
  Unterscheidungsprobe der Schleife, die Mutation des Implementers (9 von 1000) ist **übernommen**.

---

## Findings

### F-1 — Stichprobe der 79 unveränderten Pipes: Erzeuger Variable, keine Prozess-Pipe; Restgrenze 64 KiB bleibt an der Variablengröße der Test-Container-Logs

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Ursprung hergeleitet vs. gemessen)
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh:272` (ebenso `…kotlin…:270`, `…python…:306`, `lib-sdk-filter-fixture.sh:190`, `lib-sdk-route-fixture.sh:168`)
- `befund`: Ich habe `git grep -E '\| *grep +-[a-zA-Z]*q'` auf Zeilen ohne `printf`/`echo` gefiltert: es bleiben nur `examples/compose.yaml:60` (`wget -q -O - … | grep -q ok`, `sh` ohne `pipefail`)
  und `run-fmt-check-tests.sh:175` (`grep -A1 … | grep -Fxq`, zwei Zeilen); alle übrigen sind `printf '%s' "$var" | grep -q…` aus einer Variable, in Übereinstimmung mit dem Plan. Die größte
  Variable ist `test_output=$(docker logs … 2>&1)` des Test-Containers (Gradle/dotnet/pytest); der Abstand zur Schwelle (zwischen 48 und 200 KiB, Pipe-Puffer 64 KiB **übernommen**)
  ist mit 7,5–8,1 KB (Messung des Implementers) etwa der Faktor 8 und hängt an der Ausführlichkeit des Test-Frameworks, nicht an der Laufzeit (der Container endet, bevor gelesen wird; der Lauf-Log
  von 38 214 Byte ist eine Summe vieler Variablen, keine einzelne). Ein Überlauf färbte die `!`-Prüfungen rot (kein falsches Grün); nur `nats_reconnect_before_output`
  (`run-integration-tests.sh:1970`, Negativ-Prüfung) färbte falsch grün, deren Variable ist eine kleine Subscriber-Log-Ausgabe.
- `verifizierbar`: nein — kein Sensor liest Variablengrößen; Nachmessen je Variable wäre ein eigener Lauf.
- `klasse`: Restgrenze einer „bleibt“-Entscheidung

### F-2 — Negativ-Prüfungen (Filter-Rundläufe): SIGPIPE-Pfad geschlossen, Ausfall von `docker logs` selbst bleibt unbemerkt

- `kategorie`: INFO
- `quelle`: Maintainability; `AGENTS.md` §3.9 sinngemäß
- `pfad`: `tools/harness/run-integration-tests.sh:2574` und `:3185` (Parent 2570, 3181)
- `befund`: Die Aussage „vorher falsch grün statt rot“ trägt: bei einem vorhandenen Treffer (Filter-Leck) mit Folgezeilen verwarf `grep -q` die Pipeline per SIGPIPE unter `pipefail`, der
  `if` blieb falsch, und das Leck wäre unbemerkt geblieben; die neue Form liest bis Dateiende und schließt den Fall (Mechanismus mit der Reproduktion oben **gemessen**, die
  Negativ-Stelle selbst **hergeleitet**, nicht mutiert). Unverändert bleibt: schlägt `docker logs` selbst fehl (Container entfernt, `2>/dev/null`), endet die Pipeline mit ≠ 0 und die Prüfung meldet kein Leck.
- `verifizierbar`: nein — der Fehlerpfad bräuchte eine Mutation (Container vor der Prüfung entfernen); nicht Teil dieses Slice.
- `klasse`: Negativ-Prüfung ohne Bindung an den Ausfall ihres Erzeugers

### F-3 — Gemeldete Fremdstellen `pin-stale-*`: Wert nach Herleitung nicht betroffen

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/pin-stale-actions.sh:43`, `pin-stale-baseline.sh:26`, `pin-stale-dcheck.sh:31`
- `befund`: `latest=$(github_api_get … | grep -m1 '"tag_name"' | sed …)` unter `set -uo pipefail` (ohne `-e`, Zeile 19 der ersten Datei). `grep -m1` gibt die Treffer-Zeile aus, bevor es endet; `sed`
  bekommt sie also, der gelesene Wert stimmt, und der Pipeline-Status wird nicht ausgewertet (kein `-e`, kein `$?`-Test). Die Fehlwirkung unterscheidet sich damit von `-q` mit `!`-Test:
  höchstens eine Fehlermeldung des Erzeugers auf stderr (z. B. beim Schreibfehler), keine falsche Verzweigung. **Hergeleitet**, nicht gemessen; die Stellen sind Netz-Werkzeuge ohne Gate. Die Meldung des Implementers
  an den Planner ist angemessen.
- `verifizierbar`: nein — benötigt Netz und GitHub-API.
- `klasse`: Nachbar-Form derselben Mechanik, anderer Folgeeffekt

### F-4 — Kopf-Kommentare: Form in Ordnung

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7
- `pfad`: `tools/harness/run-integration-tests.sh:87-90`, `lib-sdk-filter-fixture.sh:24-27`, `lib-sdk-route-fixture.sh:16-19`, `run-sdk-{csharp,kotlin,python}-integration-tests.sh` (je vier Zeilen vor `set -euo pipefail`)
- `befund`: Sechs Absätze im Indikativ, keine Kennung, keine Chronik, keine verworfene Alternative; sie tragen eine Kopplung (Prüfform ↔ `pipefail`). Die Wendung „ohne `-q`“ nennt eine Abwesenheit als Eigenschaft
  der Form, nicht als früheren Text; kein Befund der Klasse „Kommentar trägt keine der Klassen“.
- `verifizierbar`: ja — `make kommentar-kennungen DIFF=3bcafd41` Exit 0.
- `klasse`: —

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (21 Ersetzungen Zeile für Zeile; Muster, `-F`/`-E`, Quoting, Kontext)
- geprüft, ohne Befund: `tools/harness/run-sdk-{csharp,kotlin,python}-integration-tests.sh` (je 3 Ersetzungen, `$reject_marker` unverändert gequotet)
- geprüft, ohne Befund: `tools/harness/lib-sdk-filter-fixture.sh`, `lib-sdk-route-fixture.sh` (3 bzw. 2 Ersetzungen, beide `-E`-Anker erhalten)
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` (53 Lokatoren, jeweils +4, sonst unverändert)
- geprüft, ohne Befund: Plan §3 (Suchlauf-Feld, `make suchlauf-nachmessen` Exit 0; Summe 57+3+2+4+3+3+4+1+1+1 = 79 stimmt mit der Messung)
- geprüft, ohne Befund: Hard Rules §3.1 (keine Host-Werkzeuge, Scratchpad-Skript außerhalb des Repos), §3.2 (keine Suppression), §3.9 (Exit-Codes separat), §3.15 (keine Verweigerung berichtet)
- nicht geprüft: Gesamtlauf `make test-integration` und die beiden Tier-Läufe C#/Kotlin (nur Python gefahren); Realdaten des ersten Post-Push-Laufs von `e2e` (Plan §2 Liefer-Punkt 2, offen bis nach dem Push)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Verdikt:** keine Fixrunde am Implementer nötig. Merge-blockierend: nein. Der Plan-Punkt „Review durchgeführt“ ist mit diesem Report nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug).
