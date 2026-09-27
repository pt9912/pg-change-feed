# Review-Report: slice-harness-guard-blocked-python — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, Nutzer-Entscheidung „Weg 3“ und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-harness-guard-blocked-python` (wellenlos), Diff-Range `cea198fb..HEAD` insgesamt 22
Commits/26 Dateien; der für diesen Slice relevante Teil sind die vom Auftrag benannten Commits: Lifecycle `1b37cb01`,
`3036046b`, `ca97b802`; Implementierung `e0bffcb2` (Guard + Fragment `tools/harness/blocked/python`), `470c36a1`
(Tabellentest-Umbau + Makefile-Hilfezeile), `d5afaaa8` (`MR-004` + Index), `3287b044` (`AGENTS.md` + `harness/README.md`),
`c9771b83` (Beobachtungs-Register), `cb1df743` (Link-Fix in fremder Datei), `42b5a9ca` (DoD-Haken, Suchlauf-Zeilen),
`fa1e95de` (restliche DoD-Haken). Kein Go-Code im Diff (Guard, Tabellentest, Fragment sind Bash; Träger sind Markdown).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen
ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle Mutationen liefen an
Scratchpad-Kopien des Guards/Tabellentests (eigene Kopien im Scratchpad, `BLOCKED_DIR=<Kopie>`/`GUARD=<Kopie>`
übersteuert; nie `sed -i`, nie eine Umleitung auf eine Repo-Datei). Jeder Hook-Probe-Aufruf lief über ein eigenes
Scratchpad-Skript (`json_escape` in reinem Bash, kein Host-`python`) direkt gegen `bash
.claude/hooks/pretooluse-command-guard.sh`; kein Aufruf von Host-`python`/`python3` fand statt.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-harness-guard-blocked-python` (§1 Ziel/Ausgangslage, §2 DoD, §3 Plan mit Suchlauf-Feld, §6 Risiken)
- `harness/conventions/MR-003-guard-inplace-textwerkzeug.md` (unverändert, `Accepted`, immutable),
  `harness/conventions/MR-004-guard-host-python-am-kopf.md` (neu)
- `AGENTS.md` §3.1 (Docker-only, Durchsetzung), §3.6, §3.7, §3.9, §3.11 (kein host-lokaler absoluter Pfad — diese
  Datei ist selbst danach korrigiert, siehe unten), §3.12, §3.13
- `.harness/baseline/v6.9.0/regelwerk/modul-13-quality-gates.md` §Guard-Härtung (Ersetzt-Baseline-Regel von `MR-004`)
- Beobachtungs-Register `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`,
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`
- Nutzer-Entscheidung „Weg 3“ (Closure-Notiz `slice-transformationen-map-value` §7)
- Vorherige Findings am Modul: `review-slice-harness-guard-inplace-textwerkzeug.md` (kein offenes HIGH/MEDIUM)

**Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren):**

- **Guard-Verhalten, 20 eigene Hook-Aufrufe** gegen den echten Guard mit dem echten Fragment (Hook-JSON auf stdin,
  kein Host-`python`): block — `python3 --version`, `python` (bare), `cd <Repo-Wurzel> && python3 -c 1`,
  `cd /tmp && python3 -c 1` (Pfad ohne Repo-Namen), `env python3 -c 1`, `bash -c "python3 -c 1"`,
  `/usr/bin/python3 -c 1`, `echo x | xargs -n1 python3`, `find . -name '*.py' -exec python3 {} \;`; pass — `make test`,
  `docker run --rm x python3 -c 1`, `command -v python3`, `echo python3`, `grep -E 'a|python3' f`,
  `git commit -m "fix python3 thing"`, `# python3 -c 1` (Kommentarzeile), `mypython3script --run` (Teil eines
  längeren Tokens), `cat sdks/python3/foo.txt` (Pfad mit `python3` im Verzeichnisnamen, nicht im Kopf),
  `python3.12 --version`, `python2 --version`. Alle 20 Fälle wie in Plan/DoD zugesagt.
- **`make test-command-guard`**: Exit 0, „alle 367 Fälle bestanden“ (Plan-Zusage 317 → 367 bestätigt).
- **`BLOCKED_DIR`-Override**: fehlendes Verzeichnis → Exit 2 mit Meldung; leeres Verzeichnis (Fragment fehlt) → 30
  `FEHLER:`; Fragment mit CRLF-Zeilenende → 29 `FEHLER:` (eigene Mutation, siehe F-2); Fragment mit zusätzlichem
  Leerraum/Leerzeilen → 367/367 grün (Ladung robust gegen Whitespace-Varianz).
- **52-vs-50-Diskrepanz aus §1 der Plan-Datei aufgelöst:** die „52 Fälle“ sind eine live gefahrene Emulation — Kopie
  des **alten** Guards (`cea198fb`) mit `python`/`python3` direkt im gebackenen `BLOCKED`-String (nicht über ein
  Fragment) gegen den **alten** Tabellentest (317 Fälle) laufen lassen: reproduziert exakt 52 `FEHLER:` (24× „Klasse
  interp fehlt“, 27× „Pass-Fall erwartet“, 1× „Meldung interp“ — Zahlen einzeln nachgezählt, deckungsgleich mit der
  Planaussage). Die „50“ ist eine andere, ebenfalls korrekte Zahl: `367 − 317 = 50`, die Netto-Differenz der
  **fertig umgebauten** Tabellenzeilen (`git diff cea198fb..HEAD -- tools/harness/run-command-guard-tests.sh`:
  50 `^-block|^-pass`-Zeilen entfernt, 100 `^+block|^+pass`-Zeilen hinzugefügt). Beide Zahlen sind reale, verschiedene
  Messungen an verschiedenen Objekten (Ist-Analyse vor dem Umbau vs. Zeilenbilanz nach dem Umbau) und widersprechen
  sich nicht; §3.12 ist hier eingehalten. Stichprobe von 20 der 52 betroffenen Namen gelesen: jeder ist entweder auf
  Kopf `python3.12` umgestellt (bindet weiter die Repo-Pfad-Regel, u. a. alle 16 „Pfadzeichen davor/dahinter“-Fälle,
  die 5 Positionsfälle `env`/Backslash/`find -exec`/`eval`/Wrapper) oder auf die neue `pkg`-Erwartung reassigniert
  (u. a. „Rand: cd im selben Kommando“ → `block pkg "hinter cd /tmp/x &&, mit Repo-Pfad-Text"`); keiner ist ersatzlos
  gestrichen.
- **11 behauptete Mutationen, 6 selbst nachgefahren** (mindestens 3 gefordert): Fragment fehlt → 30 (bestätigt);
  Fragment leer (0 Byte) → 30 (bestätigt); Liste um `python3` gekürzt (nur `python`) → 29 (bestätigt); Liste um
  `python` gekürzt (nur `python3`) → 1 (bestätigt, einziger Fall mit bloßem Kopf `python`: „Kopf: python“); Liste um
  `go` erweitert → 1 (bestätigt: „Grenze: go bleibt ungelesen“); `env` aus `PREFIXES` entfernt → 8 (bestätigt, exakte
  Fälle gelistet). Alle sechs stimmen mit den im Auftrag genannten Zahlen überein.
- **`make doc-immutable RANGE=cea198fb..HEAD`**: Exit 0, 1345 Dateien, 0 Befunde — `MR-003` unverändert bestätigt.
- **`make docs-check`**: Exit 0, 1345 Dateien, 0 Befunde (Anker `guard-haertung` in `MR-004` löst auf).
- **`make commit-traceability RANGE=cea198fb..HEAD`**: OK, 22 Commits, keine Struktur-ID im Betreff.
- **`make kommentar-kennungen DIFF=cea198fb`**: Exit 0, kein Kandidat (kein Go im Diff).
- **`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-guard-blocked-python.md`**: **4 von 14
  Zeilen weichen ab** (siehe F-1) — sowohl am aktuellen `HEAD` als auch, per `git worktree add --detach 42b5a9ca`
  nachgestellt, exakt am Commit, der die Nachmessung geschrieben hat.
- Lifecycle-Moves `1b37cb01`/`ca97b802` sind reine Renames (0 Zeilen, `git show --stat`); `3036046b`
  („Verantwortlich gesetzt“) ist ein eigener Inhalts-Commit; `cb1df743` ändert in der fremden Datei
  `slice-harness-mutationsbild-und-verweigerte-aktion.md` ausschließlich die eine Verweiszeile (Link → zitierte
  Kennung), sonst nichts (`git show` gelesen).

**Eigener Fehlgriff, selbst korrigiert:** der erste Entwurf dieses Reports nannte in der Aufzählung der 20
Hook-Proben den host-lokalen absoluten Pfad des Arbeitsverzeichnisses (`AGENTS.md` §3.11 — kein host-lokaler
absoluter Pfad in der Doku, auch nicht im Fließtext). `make docs-check` (Modul `hostpaths`) fing das vor dem Commit
(1 Befund, `hostpath-forbidden`). Korrigiert auf `<Repo-Wurzel>` in der Aufzählung oben; erneuter `make docs-check`-Lauf
danach grün (siehe „Eigene Messungen“).

## Findings

### F-1 — Der §3.13-Suchlauf behauptet vier Zeilen, die nicht reproduzieren, schon am Commit, der sie schreibt

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz A/B, §3.13; Skill-HIGH „Zahl im Träger ohne Ursprung — oder gegen die Messung
  driftend“; Skill-HIGH „Beleg trägt seinen Satz nicht“ (der DoD-Haken nennt `make suchlauf-nachmessen` als seinen
  eigenen Beleg-Anker)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-guard-blocked-python.md:311,312,315,316` (die vier
  abweichenden `diff`-Zeilen des `suchlauf`-Blocks), `:320-355` (Prosa „Nachmessung des Implementers“), `:213-217`
  (DoD-Haken „§3.13-Suchlauf“, mit `[x]` markiert)
- `befund`: `make suchlauf-nachmessen PLAN=…` meldet am aktuellen `HEAD` (`fa1e95de`) vier Abweichungen: Zeile 1
  (Symbolname Fragment/Kopf-Liste) soll `18`, ist `23`; Zeile 2 (Host-Interpreter-Beschreibung) soll `52`, ist `53`;
  Zeile 5 (Skript-liest-Interpreter) soll `3`, ist `4`; Zeile 6 (`python3? --version`) soll `14`, ist `16`. Exit 1,
  „4 von 14 Zeilen weichen ab“. Um auszuschließen, dass die Abweichung erst durch spätere, unabhängige Slices
  im breiteren Diff-Bereich `cea198fb..HEAD` entstanden ist, habe ich per `git worktree add --detach 42b5a9ca` genau
  den Commit nachgestellt, der die Nachmessungs-Prosa und die Zahlen selbst geschrieben hat (danach folgt nur noch
  `fa1e95de`, das ausschließlich zwei DoD-Häkchen setzt, keine Prosa ändert): dieselben vier Zeilen weichen dort
  identisch ab (`18`/`23`, `52`/`53`, `3`/`4`, `14`/`16`). Die Abweichung ist also nicht Drift nach der Niederschrift,
  sondern die Zahlen waren **bereits beim Commit, der sie als „gemessen mit `make suchlauf-nachmessen`“ ausgibt,
  falsch** — der DoD-Haken „§3.13-Suchlauf“ ist mit `[x]` markiert, obwohl sein eigener benannter Beleg-Anker beim
  Nachfahren nicht grün läuft. Das ist genau die Klasse, die dieser Sensor fangen soll (`harness/sensors/suchlauf-
  nachmessen.md`): eine im Träger committete Zahl, die sich als „gemessen“ ausgibt und die Messung nicht trägt.
  Zehn der vierzehn Zeilen (alle sieben `cea198fb`-Baseline-Zeilen und drei der sieben `diff`-Zeilen) stimmen; der
  Fehler betrifft ausschließlich vier `diff`-Zeilen mit unspezifischen Suchmustern (z. B. `Host-.?python|Host
  python|Host-Interpreter|Sprach-Toolchains`), die zusätzliche Treffer in Nachbar-Dateien der laufenden Session
  eingesammelt haben, ohne dass die Prosa das vor dem Commit ein letztes Mal gegengeprüft hätte.
- `verifizierbar`: ja — `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-guard-blocked-python.md`
  (Exit 1, druckt Soll/Ist je Zeile); reproduzierbar am Commit `42b5a9ca` über `git worktree add --detach 42b5a9ca`
- `klasse`: Zahl im Träger driftet gegen die eigene Messung (Suchlauf-Feld)

### F-2 — Die Fragment-Ladung tokenisiert roh auf Whitespace; ein CRLF-Zeilenende im Fragment lässt das letzte Wort der Zeile unbemerkt durch

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“; `AGENTS.md` §3.1 (die Zusage lautet
  „blockt IMMER … unabhängig vom Befehlsstring“); MR-004 §Grenz-Zeile (nennt viele Grenzen, diese nicht)
- `pfad`: `.claude/hooks/pretooluse-command-guard.sh:92-98` (`BLOCKED="$BLOCKED $(cat "$bf")"`, keine
  Zeilenenden-Normalisierung); `harness/conventions/MR-004-guard-host-python-am-kopf.md:57-74` (§Grenz-Zeile — nennt
  nicht gelistete Namen, Variablen/Aliase, Umleitungen, Heredoc-Zeilen, aber keine Zeilenenden-Form des Fragments)
- `befund`: Eigene, im Plan/Bericht nicht genannte Mutation: Fragment-Datei mit `python python3\r\n` (Windows-
  Zeilenende) statt `python python3\n`. Bash-Wortsplitting trennt auf `IFS` (Leerzeichen, Tab, Newline) — `\r` ist
  keines davon, es bleibt am letzten Wort der Zeile hängen (`python3\r`). Das Wort matcht `in_set` dann nicht mehr
  gegen das Literal `python3`, während `python` (nicht das letzte Wort der Zeile) unverändert matcht. Ergebnis: 29 von
  367 Fällen färben rot — alle `python3`-Kopf-Fälle (Block-Fälle fallen ganz durch bzw. auf die schwächere
  Repo-Pfad-Regel zurück, `python3 --version`/`-c 'print(1)'` passieren unbemerkt), während `python`-Kopf-Fälle
  weiter korrekt blocken. Belegt auch direkt an der Hook-Schnittstelle: `python3 -c 1` gegen ein CRLF-Fragment endet
  ohne Ausgabe (Pass), obwohl die Zusage „blockt IMMER“ lautet. Die aktuelle Fragment-Datei im Repo ist LF-terminiert
  (`cat -A tools/harness/blocked/python` → `python python3$`, kein `^M`) — der Fehler ist **latent, nicht aktiv**;
  Risiko ist gering (Docker-only-/Linux-Pfad des Repos, `LH-QA-POR-002`), aber die Zusage „unbedingt … unabhängig vom
  Befehlsstring“ gilt nur so lange, wie niemand das Fragment mit einem CRLF-Editor anfasst — ein Fall, den weder der
  367-Fälle-Tabellentest noch die Grenz-Zeile von `MR-004` (die sonst sehr genau ist, u. a. Heredoc-Zeilen benennt)
  abdeckt.
- `verifizierbar`: ja — Mutation `printf 'python python3\r\n' > <Kopie>` gegen `BLOCKED_DIR=<Kopie>`; 29 `FEHLER:`
- `klasse`: Fragment-Ladung ungetestet gegen Zeilenenden-Form

### F-3 — Selbstbefund einer Umleitung während der eigenen Implementierung: korrekt dokumentiert, ohne Wirkung im Diff

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.1 „Text-Umschreiben im Repo ist Sache der Datei-Werkzeuge des Laufs“; Skill-HIGH
  „Docker-only-Verstoß“ (Klausel Umleitung) — geprüft, im **Endergebnis** nicht erfüllt
- `pfad`: `docs/plan/planning/observations/BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/evidence/slice-harness-guard-blocked-python.md`
  (achter Beleg des Registers); `git log --follow -- tools/harness/blocked/python` (ein einziger Commit, `e0bffcb2`)
- `befund`: Der Implementer meldet selbst, das Fragment zunächst per `mkdir -p … && printf … > tools/harness/blocked/
  python` (Umleitung) angelegt zu haben, die Verletzung im selben Werkzeugaufruf-Batch bemerkt, die Datei sofort
  entfernt und über das Write-Tool neu angelegt zu haben, bevor ein Folgeschritt darauf aufbaute. Verifiziert:
  `tools/harness/blocked/python` hat in der gesamten Historie genau einen Commit (`e0bffcb2`) — die per Umleitung
  entstandene Fassung hat nie einen Commit, geschweige denn den finalen Diff, erreicht. Der Selbstbefund ist präzise
  klassifiziert („eingetreten, ohne Wirkung im Endergebnis“, Register-`state.md` fortgeschrieben, Zähler korrekt von 7
  auf 8 erhöht). Keine Aktion nötig — die Regel wirkt genau wie vorgesehen (Review als Wächter dessen, was der Guard
  selbst nicht liest, `AGENTS.md` §3.1 „Was der Guard nicht liest, bleibt Sache des Reviews“), hier bereits durch
  Selbstkorrektur im selben Lauf, bevor ein Review überhaupt nötig wurde.
- `verifizierbar`: ja — `git log --follow --oneline -- tools/harness/blocked/python` (ein Commit)
- `klasse`: Selbstbefund, korrekt dokumentiert, keine Aktion

### F-4 — Wechsel auf das Write-Tool nach einem vom Guard geblockten Heredoc für die Commit-Message ist ein vorgesehener, kein ausweichender Weg

- `kategorie`: INFO
- `quelle`: VERWEIGERTE-AKTION-ZWISCHENREGEL (Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion`, noch kein
  Repo-Träger); Aufgabenstellung des Auftrags
- `pfad`: keine Repo-Fundstelle (das Ereignis selbst ist nicht committet; es steht nur in der Aufgabenstellung
  dieses Reviews)
- `befund`: Laut Auftrag blockte der PreToolUse-Guard einen Bash-Heredoc, der eine Commit-Message-Datei füllen sollte
  (der Aufruf enthielt sehr wahrscheinlich eine Zeile, die selbst wie ein Guard-Treffer aussieht — dieser Plan zitiert
  an vielen Stellen wörtlich Kommandostrings wie `python3 - <<'EOF'`); der Implementer wechselte auf das Write-Tool,
  um dieselbe Datei zu erzeugen, und rief anschließend `git commit -F <Datei>` auf. Einordnung: `git commit -F
  <Datei>` ist der im Auftrag selbst vorgeschriebene, reguläre Weg für Commit-Messages in diesem Workflow (nicht erst
  als Ausweichroute erfunden), und eine Commit-Message-Datei ist kein Repo-getrackter Inhalt, den Edit/Write
  „umgehen“ würde — sie ist ohnehin nur über eine Datei zu übergeben. Der Wechsel wiederholt die geblockte Handlung
  nicht auf einem anderen Weg zum selben (verbotenen) Ziel, sondern wechselt auf einen bereits sanktionierten,
  andersartigen Mechanismus für ein anderes Artefakt (Commit-Message-Datei statt Repo-Quelldatei). Ich sehe darin
  keinen Verstoß gegen die Zwischenregel und keinen Eskalationsbedarf; eine abschließende Bewertung, ob der
  Implementer das Ereignis auch im eigenen Bericht genannt hat (wie die Zwischenregel verlangt), kann ich aus dem
  Repo allein nicht führen — dafür fehlt ein committeter Träger dieses konkreten Vorgangs (anders als beim in F-3
  behandelten Umleitungs-Fund, der einen eigenen Registereintrag bekam).
- `verifizierbar`: teilweise — der Zielmechanismus (`git commit -F`) ist im Auftrag und in `AGENTS.md` §3.1 als Weg
  belegt; das Blockereignis selbst ist außerhalb des Diffs und nicht nachvollziehbar
- `klasse`: Ersatzweg-Bewertung ohne Repo-Beleg

## Negativbefunde

- geprüft, ohne Befund: **Guard-Kern (Segmentierung, Präfixe, in-place-Formen)** — die einzige Code-Änderung im
  Guard ist der Text von `REASON_PKG` (neuer Satz mit den §3.1-Wegen, weder `"` noch `\` enthalten, alle 20 eigenen
  Hook-Aufrufe liefern gültiges JSON); die Scan-Logik selbst ist unverändert (`git diff` zeigt ausschließlich
  Kommentar-Zeilen als Code-Diff außer der einen `REASON_PKG`-Zeile).
- geprüft, ohne Befund: **`python3.<N>`/`perl` bleiben unter der Repo-Pfad-Regel** — eigene Proben `python3.12
  --version` und `python2 --version` ohne Repo-Pfad passieren den Guard (Pass), wie zugesagt; der Zweig
  `python|python[0-9]*|perl` (Zeile 258–262) bleibt für diese Namen erreichbar, der Kommentar dort nennt das korrekt.
- geprüft, ohne Befund: **Tabellentest-Umbau, keine stille Streichung** — Stichprobe von 20 der 52 betroffenen
  Bestandsfälle gelesen: jeder ist auf `python3.12` umgestellt (bindet weiter die Repo-Pfad-Regel und ihre
  Zeichenklassen) oder auf die neue `pkg`-Erwartung reassigniert; keine ersatzlose Streichung gefunden.
- geprüft, ohne Befund: **`MR-004`** — per Formvergleich mit `MR-003` (gleiches Feld-Gerüst: Datum, Geltungsbereich,
  Ersetzt-Baseline-Regel, Adaption, Begründung, Grenz-Zeile, Auflösungs-Trigger); `Ersetzt-Baseline-Regel` verlinkt
  mit Anker `guard-haertung`, der über `make docs-check` auflöst; die Begründung „Gehärtet wird die Zerlegung, nicht
  die Denylist“ ist stichhaltig, weil die Kopf-Liste nur den Kopf liest (`make`/`docker` bleiben frei, eigene Proben
  #5/#6 bestätigen das) — die zwei von `MR-003` offen gelassenen Ränder (`cd`-Präfix, Pfad ohne Repo-Namen) sind
  keine Zerlegungsfrage, weil die unbedingte Liste den Rest des Befehlsstrings gar nicht mehr liest; `MR-003` selbst
  unverändert (`make doc-immutable` Exit 0).
- geprüft, ohne Befund: **Index-Zeile `harness/conventions.md`** — `MR-004` mit Anker `mr-004`, korrekter Link,
  korrekte Geltungsbereich-Spalte.
- geprüft, ohne Befund: **Träger** — `AGENTS.md` §3.1 Absatz „Durchsetzung“ nennt `MR-004` neben `MR-003`, die neue
  Grenze korrekt (Kopf unbedingt vs. Rest unter Repo-Pfad-Regel); `harness/README.md:155` zieht Fragment,
  `BLOCKED_DIR` und beide `MR`-Links nach; Makefile-Hilfezeile nennt `BLOCKED_DIR`; beide Kopfkommentare (Guard,
  Tabellentest) im Indikativ, ohne Konjunktiv über verworfene Alternativen, höchstens Rang-Zeiger auf `MR-003`/
  `MR-004` als Kopplung (kein Go-Kommentar, die „eine Kennung“-Regel aus §3.7 ist dort nicht einschlägig).
- geprüft, ohne Befund: **Die zwei Register-`state.md`** — `inplace-textwerkzeug-am-repo-trotz-nutzerregel`: Ausgang
  „verkörpert“ korrekt begründet, achte Beleg-Datei (F-3) korrekt geführt, Zähler stimmt (8, `ls evidence | wc -l`);
  `host-werkzeug-jenseits-docker-und-make-ohne-deklaration`: Zeiger auf die Durchsetzung korrekt, kein eigenes
  Auftreten behauptet.
- geprüft, ohne Befund: **Link-Fix in fremder Datei** (`cb1df743`) — `git show` zeigt ausschließlich den Wechsel
  wandernder Link → zitierte Kennung in genau einer Zeile von `slice-harness-mutationsbild-und-verweigerte-aktion.md`,
  keine inhaltliche Änderung.
- geprüft, ohne Befund: **Traceability, ID-Schema, Moves** — alle 10 relevanten Betreffe tragen `(ADR-0083)`, keine
  `SPEC-*`/`ARC-*`; `make commit-traceability RANGE=cea198fb..HEAD` OK (22 Commits im Gesamt-Range); die drei
  Lifecycle-Commits `1b37cb01`/`ca97b802` sind reine Renames, `3036046b` ein eigener Inhalts-Commit; kein `//nolint`
  (kein Go im Diff).
- geprüft, ohne Befund: **Docker-only im finalen Diff** — kein Host-Toolchain-Install, kein `sed -i`/in-place, keine
  Umleitung auf eine Repo-Datei im committeten Ergebnis (F-3 betrifft nur einen zurückgenommenen Zwischenschritt).
- geprüft, ohne Befund: **Spec-Stratum, Zwei-Quellen-Drift, Sicherheit, kritischer Pfad** — `spec/` nicht berührt
  (Harness-Wächter, keine Spec-Stelle laut Plan-Kopf); keine neue `Accepted`-ADR verändert; kein Geheimnis im Diff;
  kein Produktionscode/`internal/`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zahl im Träger driftet gegen die eigene Messung (Suchlauf-Feld) · Fragment-Ladung
ungetestet gegen Zeilenenden-Form · Selbstbefund, korrekt dokumentiert, keine Aktion · Ersatzweg-Bewertung ohne
Repo-Beleg

## Verdikt

**Merge-blockierend: ja, F-1 (HIGH).** Der Guard-Mechanismus selbst ist solide: 20 eigene Hook-Proben, 6 eigene
Mutationen und der volle Tabellentest-Lauf (367/367) bestätigen die zugesagte Wirkung ausnahmslos, einschließlich der
im Plan explizit benannten Ränder (`cd`-Präfix, Pfad ohne Repo-Namen, `--version`). Der Befund betrifft nicht den
Guard, sondern die **Prosa des Slice-Plans**: der §3.13-Suchlauf, der genau diese Art Drift verhindern soll, trägt
selbst vier falsche Zahlen — nachweislich schon zum Zeitpunkt des Commits, der sie schreibt, nicht erst durch
spätere Drift. Der DoD-Haken „§3.13-Suchlauf“ ist zu Unrecht auf `[x]`.

**Fixrunde nötig:** Implementer/Planner läuft `make suchlauf-nachmessen PLAN=…` am aktuellen Stand erneut, korrigiert
die vier abweichenden Zeilen im `suchlauf`-Block (Zeilen 311/312/315/316) und den zugehörigen Prosa-Absatz
„Nachmessung des Implementers“ (Zeilen 320–355) auf die reproduzierbaren Werte, bevor der DoD-Haken erneut auf `[x]`
gesetzt wird. F-2 (MEDIUM) ist eine Ergänzung der Grenz-Zeile in `MR-004` (ein Punkt: Fragment-Ladung normalisiert
keine Zeilenenden) plus optional eine neue Zeile im Tabellentest — kein Code-Fix am Guard zwingend, da der aktuelle
Fragment-Inhalt LF-terminiert ist und keine akute Wirkung hat; die Ergänzung dokumentiert die Grenze ehrlich, wie es
die übrigen Punkte von `MR-004` bereits vorbildlich tun. F-3 und F-4 sind INFO ohne erwartete Aktion.

Die DoD-Zeile „Review durchgeführt“ bleibt **offen** (kein Nachzug in diesem Commit) — die Fixrunde an F-1 führt
regulär über Schritt 21 des Implementer-Ablaufs zurück zum Reviewer.

**Übergabe:** F-1 an den Implementer/Planner (Suchlauf-Zeilen und Prosa korrigieren, DoD-Haken danach erneut setzen);
F-1 zusätzlich als Beobachtung notierbar, sollte ein drittes Auftreten der Klasse „Suchlauf-Nachmessung committet vor
dem letzten Gegencheck“ auffallen (Steering-Loop, Skill §Pflege). F-2 an den Implementer/Planner für `MR-004`
§Grenz-Zeile (optional: ein Tabellentest-Fall mit CRLF-Fragment als benannte Grenze statt aktivem Fix). F-3 und F-4
ohne erwartete Aktion; F-4 bleibt ohne Repo-Beleg des konkreten Vorgangs — sollte der Implementer-Bericht das
Blockereignis nicht genannt haben, ist das eine Frage an die Orchestrator-Ebene, nicht an dieses Review. Der Report
ist ein Lauf-Beleg und ersetzt keine Verifikation.
