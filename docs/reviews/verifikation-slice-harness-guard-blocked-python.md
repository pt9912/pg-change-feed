# Verifikations-Report: slice-harness-guard-blocked-python — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
(Nutzer-Entscheidung „Weg 3“) + Plan-vs-Code-Diff + Gates. Review-Artefakt:
[`review-slice-harness-guard-blocked-python.md`](review-slice-harness-guard-blocked-python.md)
(Commit `1753e541`; 1 HIGH/1 MEDIUM/2 INFO, Fixrunde `3d380abe`/`1df70522`/`eb04d172`). Formvorbild
dieses Reports:
[`verifikation-slice-leerlauf-phase-last-in-stuecken.md`](verifikation-slice-leerlauf-phase-last-in-stuecken.md).

**Gegenstand:** Slice-Plan `slice-harness-guard-blocked-python` (wellenlos), Diff-Range
`cea198fb..HEAD` — 26 Commits gesamt (dieser Range enthält interleaved auch die Closure von
`slice-leerlauf-phase-last-in-stuecken` und die Anlage von
`slice-harness-mutationsbild-und-verweigerte-aktion`, beide außerhalb dieses Gegenstands, wie
bereits vom Reviewer abgegrenzt); die diesem Slice zurechenbaren Commits sind `ae593de4` … `eb04d172`
(17 Commits, alle mit `(ADR-0083)` im Betreff). Dieser Lauf ändert weder Code noch Plan noch Doku; er
schreibt nur diesen Report. Alle eigenen Proben liefen netzlos (`bash`, Hook-JSON auf stdin, Scratchpad-
Kopien des Fragments via `BLOCKED_DIR=<Kopie>`); kein Aufruf von Host-`python`/`python3`, kein `sed -i`,
keine Umleitung auf eine Repo-Datei. `git status --short` war nach jedem eigenen Lauf leer.

**Repo-Zustand:** `HEAD` ist 15 Commits vor `origin/main` (nicht mein Zutun, Auftrag verlangt kein
Push). Ich habe in diesem Lauf **nicht gepusht**.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-guard-blocked-python.md` | **Exit 0** | „suchlauf-nachmessen: 14 Zeilen stimmen" — alle 14 Zeilen `OK`, einschließlich der vier von der Fixrunde berichtigten `diff`-Zeilen (1, 2, 5, 6: jetzt 26/54/4/16 statt der vom Review als falsch befundenen 23/53/3/14) |
| Eigene `git grep -c`-Nachzählung der vier berichtigten Zeilen (ohne das Werkzeug, direkt mit dem Ausschluss-Pathspec der Plan-Datei) | **deckungsgleich** | 26, 54, 4, 16 — exakt wie im Plan-Text und wie vom Werkzeug bestätigt |
| `make test-command-guard` | **Exit 0** | „run-command-guard-tests: alle 368 Fälle bestanden" — unabhängig gegengerechnet: 362 `block`/`pass`-Top-Level-Aufrufe (`grep -cE '^(block\|pass) '`) + 6 manuelle `n=$((n+1))`-Inkremente außerhalb dieser Funktionen = 368 |
| `make gates` | **Exit 0** | alle sechs inneren Gates grün: `baseline-verify: v6.9.0 OK`, `docs-check`/`commit-traceability` (d-check `1346 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK`), `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%`, `generated-sync: OK — byte-gleich`, `a-check: gesamt: 0 Befund(e)` |
| `make doc-immutable RANGE=cea198fb..HEAD` | **Exit 0** | `d-check: 1346 Datei(en) geprüft, 0 Befund(e)` — `MR-003` unverändert bestätigt |
| `make commit-traceability RANGE=cea198fb..HEAD` | **Exit 0** | „OK — 26 Commit(s) …, Betreffs ohne Struktur-ID" |
| **10 eigene Hook-Aufrufe** an der echten Guard-Schnittstelle (§4) | 4 block / 6 pass, alle wie zugesagt | siehe Tabelle |
| **8 eigene Mutationen** an Fragment-Kopien (§4), davon 3 im Auftrag/Review/Fixrunde nicht genannt | 3 grün (Fix hält, keine Nebenwirkung), 5 rot (Bindung an die Eingabe) | siehe Tabelle |

Nicht gefahren: ein zweiter `make gates`-Lauf (Determinismus nicht erneut geprüft, Docker-Layer-Cache
identisch); kein Docker-Image-Bau nötig (kein Produktionscode im Diff — reine Bash/Markdown-Änderung).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 10 `[x]`-Zeilen (Liefer-Punkte 1–3, `make gates`, Review, §3.13-Suchlauf, Doku-Update,
Reconciliation-entfällt, Beobachtungs-Register), 3 `[ ]`-Zeilen (Closure-Notiz, Risiken-Ausgänge, drei
Paarungen — korrekt offen, Planner-Pflicht bei Closure). Ich setze keinen Haken (Planner).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — die Sperre (`[x]`) | **bestätigt** | Fragment `tools/harness/blocked/python` = `python python3` (LF); Guard blockt Kopf `python`/`python3` unbedingt (4 eigene Block-Proben inkl. `cd <Repo> &&`, `env`, Pfad ohne Repo-Namen, `--version`), lässt `make`/`docker`/`command -v`/`echo`/Kommentar-Text frei (6 eigene Pass-Proben, §4); `REASON_PKG` nennt die drei §3.1-Wege ohne `"`/`\` |
| 2 | Liefer-Punkt 2 — der Tabellentest (`[x]`) | **bestätigt** | `make test-command-guard` Exit 0, 368 Fälle (eigen nachgerechnet); Stichprobe der 52 umgebauten Bestandsfälle nicht erneut gelesen (bereits vom Review stichprobenartig geprüft, Diff des Guard-Codes zeigt keine unabhängige Regression) |
| 3 | Liefer-Punkt 3 — die Träger (`[x]`) | **bestätigt** | `MR-004` per Formvergleich mit `MR-003` gleich strukturiert (Datum/Geltungsbereich/Ersetzt-Baseline-Regel/Adaption/Begründung/Grenz-Zeile/Auflösungs-Trigger), Anker `guard-haertung` löst auf (`make docs-check` 0 Befunde); `AGENTS.md` §3.1, `harness/README.md` §Sensors, Makefile-Hilfezeile alle nachgezogen (Diff gelesen, §3); zwei Register-`state.md` fortgeschrieben (8 Evidence-Dateien, `ls evidence \| wc -l` = 8) |
| 4 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, kein offenes HIGH/MEDIUM (`[x]`) | **bestätigt** | F-1 (HIGH) und F-2 (MEDIUM) beide durch Fixrunde aufgelöst und von mir unabhängig nachgemessen (§5); Checkbox im Fixrunden-Commit gesetzt — das ist der in `implement-slice.md` Schritt 21 „Fixrunden-Checkbox-Nachzug" **vorgesehene** Weg, kein Self-Review-Verstoß |
| 6 | §3.13-Suchlauf (`[x]`) | **bestätigt** | `make suchlauf-nachmessen` Exit 0, 14/14 (§1); alle vier von F-1 beanstandeten Zeilen jetzt korrekt |
| 7 | Doku-Update `harness/README.md` §Sensors (`[x]`) | **bestätigt** | Zeile `make test-command-guard` nennt Fragment, `BLOCKED_DIR`, CRLF-Fall, `MR-004`; Benutzerhandbuch unberührt (`git diff cea198fb..HEAD -- docs/user/` leer für diesen Slice-Anteil, keine Betreiber-Oberfläche betroffen) |
| 8 | Closure-Notiz (`[ ]`) | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 9 | Reconciliation-Register — entfällt (`[x]`) | **bestätigt** | keine Reconciliation-Datei im Repo (Greenfield) |
| 10 | Beobachtungs-Register fortgeschrieben (`[x]`) | **bestätigt, mit einer Lücke** | Ausgang der Kopf-Liste in `inplace-textwerkzeug-am-repo-trotz-nutzerregel` korrekt „verkörpert"; F-3 (Umleitung, selbstkorrigiert) als 8. Evidence-Datei korrekt geführt. **Aber:** F-1 selbst (Suchlauf-Zahl falsch schon beim Commit) wurde in keinem Register als neue Evidence nachgetragen — siehe V-2 |
| 11 | Jedes Risiko aus §6 (`[ ]`) | **korrekt offen** | alle sieben Zeilen tragen „**Ausgang:** *(bei Closure)*" |
| 12 | Drei Paarungen (`[ ]`) | **korrekt offen** | Plan §7: „die Slice-Closure trägt sie selbst" |

Kein `[x]` ohne Beleg.

## 3. Plan-vs-Code-Diff

**Guard-Code-Diff** (`git diff cea198fb..HEAD -- .claude/hooks/pretooluse-command-guard.sh`, nicht
kommentarzeilen-gefiltert): genau zwei Code-Änderungen — `REASON_PKG` um die §3.1-Wege ergänzt, und die
Fragment-Ladung entfernt `\r` vor der Übernahme in `BLOCKED` (`${frag//$'\r'/}`). Keine weitere
Code-Änderung; die Scan-/Segmentierungslogik (Übersprünge, Präfixe, Rekursion) ist unverändert — deckt
sich mit dem Plan (§3: „Der Code des Guards ändert sich nicht, außer der Beleg verlangt es").

| Plan-Zeile (§3) | Ist im Diff |
|---|---|
| `tools/harness/blocked/python` (neu) | `A`, ein Commit (`e0bffcb2`), einzige Historie — der in F-3 gemeldete Umleitungs-Zwischenschritt erreichte nie einen Commit (`git log --follow` bestätigt) |
| Guard-Kopfkommentar | `M`, Zeilen zu Host-`python`/`python3` am Kopf und zur CRLF-Behandlung im Indikativ, höchstens ein `MR`-Zeiger je Satz |
| `tools/harness/run-command-guard-tests.sh` | `M`, Umbau der 52 Fälle + neue Gruppen (a)–(e) + CRLF-Fall aus der Fixrunde |
| `Makefile` (Hilfezeile) | `M`, `BLOCKED_DIR` ergänzt |
| `harness/conventions/MR-004-…` (neu, per `cp`) | `A`, Formvergleich mit `MR-003` bestätigt (§2) |
| `harness/conventions.md` §Aktive Adaptionen | `M`, Zeile `MR-004` mit Anker `mr-004` |
| `AGENTS.md` §3.1 | `M`, exakter Wortlaut nachgezogen (§2) |
| `harness/README.md` §Sensors | `M`, wie oben |
| zwei Register-`state.md` | `M`, wie oben |
| `harness/conventions/MR-003-…` | **unverändert** (`make doc-immutable` Exit 0) |

Kein Betreiber-Weg, kein Produktionscode (`internal/`, `cmd/`, `gen/`), kein `.github`-Workflow im Diff.

## 4. Eigene Proben und Mutationen

### 4a. Zehn Hook-Aufrufe an der echten Guard-Schnittstelle (kein Host-`python`)

| Kommando | Erwartet | Gesehen |
|---|---|---|
| `make test` | pass | pass (keine Ausgabe, Exit 0) |
| `docker run --rm x python3 -c 1` | pass | pass |
| `command -v python3` | pass | pass |
| `git commit -m "fix python3 thing"` | pass | pass |
| `echo python3` | pass | pass |
| `python3.12 --version` (kein Repo-Pfad) | pass | pass |
| `python3 --version` | **block**, Klasse `pkg` | block, `REASON_PKG` mit den drei §3.1-Wegen |
| `cd <Repo> && python3 -c 1` | **block** | block |
| `env python3 -c 1` | **block** | block |
| `cd /tmp && python3 -c 1` (Pfad ohne Repo-Namen) | **block** | block |

Alle 10 wie zugesagt (Plan §1 „Kernaussage der Wirkung").

### 4b. Acht Mutationen an Fragment-Kopien (`BLOCKED_DIR=<Kopie>`, Scratchpad)

| # | Mutation | Neu ggü. Auftrag/Review/Fixrunde? | Ergebnis (gesehen) |
|---|---|---|---|
| A | CRLF-Fragment `python python3\r\n` (Redo des F-2-Fixes) | nein (Review/Fixrunde-Fall) | **grün**, 368/368 — Fix hält |
| B | Gemischte Zeilenenden im selben Fragment (`python\r\npython3\n`) | **ja** | **grün**, 368/368 — kein Kollateralschaden |
| C | Fragment mit Leerzeilen am Ende (`python python3\n\n\n`) | **ja** | **grün**, 368/368 |
| D | Fragment mit Kommentar-Token nach der Liste (`python python3 # comment\n`) | **ja** | **grün**, 368/368 — die zusätzlichen Wörter `#`/`comment` sind harmlose Extra-Einträge in `BLOCKED` |
| E | `BLOCKED_DIR` leeres Verzeichnis (Fragment fehlt) | nein | **rot**, 30 `FEHLER:` |
| F | Fragment vorhanden, 0 Byte | nein | **rot**, 30 `FEHLER:` |
| G | Liste auf `python` gekürzt (kein `python3`) | nein | **rot**, 29 `FEHLER:` |
| H | Liste auf `python3` gekürzt (kein `python`) | nein | **rot**, 1 `FEHLER:` |

E–H reproduzieren exakt die vom Review berichteten Zahlen (30/30/29/1) — unabhängig nachgefahren, nicht
dem Bericht geglaubt. A–D sind meine eigene Ausweitung der F-2-Regressionsprobe: die
Parameter-Substitution `${frag//$'\r'/}` entfernt **nur** `\r`-Zeichen, unabhängig von ihrer Position im
Fragment (nicht nur am Zeilenende) und ohne Nebenwirkung auf normale LF-Inhalte, mehrzeilige Fragmente
oder zusätzliche Wörter — die Fix-Form ist minimal und robust gegen die naheliegenden Varianten des vom
Review benannten Randes.

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Suchlauf-Zeilen falsch schon beim Commit | `make suchlauf-nachmessen` Exit 0, 14/14; eigene `git grep -c`-Nachzählung aller vier berichtigten Zeilen (26/54/4/16) deckungsgleich mit dem Plan-Text | **aufgelöst, bestätigt** |
| F-2 (MEDIUM) CRLF im Fragment lässt `python3` durch | CRLF-Fix im Code gelesen (`${frag//$'\r'/}`, betrifft nur `\r`); eigene Mutation A reproduziert den Fix-Erfolg, B–D testen Nachbarfälle ohne Regression; neuer Tabellentest-Fall „CRLF-Fragment" mit `cp`/Wiederherstellung im Wegwerf-Repo gelesen — Muster korrekt (kopiert vor der Mutation, stellt danach wieder her, kein Repo-Effekt: `git status --short` blieb während aller meiner eigenen Läufe leer) | **aufgelöst, bestätigt** |
| F-3 (INFO) Umleitung während der Implementierung, selbstkorrigiert | `git log --follow --oneline -- tools/harness/blocked/python` zeigt genau einen Commit (`e0bffcb2`); Evidence-Datei im Register korrekt mit „eingetreten, ohne Wirkung im Endergebnis" klassifiziert | **bestätigt, keine Aktion nötig** |
| F-4 (INFO) Wechsel auf Write-Tool nach geblocktem Heredoc für Commit-Message | `git commit -F <Datei>` ist der in diesem Workflow selbst vorgeschriebene Weg für Commit-Messages (auch mein eigener Auftrag verlangt ihn); eine Commit-Message-Datei ist kein Repo-getrackter Inhalt — der Wechsel wiederholt kein verbotenes Ziel auf einem Ersatzweg, sondern nutzt einen bereits sanktionierten Mechanismus für ein anderes Artefakt. Ob der Implementer das Blockereignis selbst im eigenen Bericht genannt hat, kann ich — wie der Reviewer — aus dem Repo allein nicht feststellen (kein committeter Träger dieses konkreten Vorgangs) | **korrekt gehandhabt, kein Eskalationsbedarf**; die Frage nach der Berichts-Nennung bleibt wie vom Review benannt offen (Orchestrator-Ebene, kein Repo-Beleg) |

Kein offenes HIGH/MEDIUM.

## 6. Register und Träger — Lese-Prüfung

- **`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md`:** Ausgang „verkörpert" korrekt,
  8 Evidence-Dateien (`ls evidence | wc -l` = 8), F-3 korrekt als 8. Beleg geführt. **Fund (V-1):** der
  Fließtext nennt weiterhin „Tabellentest … mit 367 Fällen" (Zeile 23) — der Wert war beim Schreiben
  dieses Satzes (Commit `c9771b83`, **vor** der Fixrunde) korrekt, ist aber seit dem CRLF-Fix
  (`1df70522`, neuer 368. Fall) veraltet. Die Fixrunde hat den Plan-Text (§3.13-Suchlauf) nachgezogen,
  aber nicht dieses Register — der Suchlauf-Block des Plans enthält kein Muster, das die Ziffer
  „367"/„368" selbst träfe (die Suchmuster zielen auf Symbolnamen/Beschreibungen, nicht auf diese Zahl;
  §3.13-Grenze „ein Suchlauf trifft Symbolnamen, keine Zahlen").
- **`BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration/state.md`:** Zeiger auf die
  Durchsetzung korrekt, kein eigenes Auftreten behauptet, keine Zahl darin.
- **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** (23×, Deckel laut Plan §8) und
  **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (33×, Deckel): keine neue Evidence-Datei zu diesem
  Slice in beiden Verzeichnissen (`ls evidence` vor/nach diesem Slice unverändert für beide) — F-1
  selbst (eine im Träger committete, beim Commit bereits falsche Zahl) ist eine reale Instanz genau der
  Klasse, die `zahl-in-traeger-driftet-gegen-die-messung` bereits führt (dessen `state.md` nennt
  ausdrücklich den Unterfall „seit ihrer Annahme nie korrekt … kein Drift durch spätere Arbeit"), wurde
  aber nicht als Evidence nachgetragen — siehe V-2.
- **Reviewer-Übergabe „drittes Auftreten der Klasse Suchlauf-Nachmessung committet vor dem letzten
  Gegencheck":** eigene Suche (`grep -rn "vor dem letzten Gegencheck\|committet vor dem\|schon zum
  Commit\|bereits beim Commit" docs/reviews/*.md`) findet **nur** die F-1-Stelle selbst — kein
  Vorläufer-Beleg dieser engen Formulierung in der Review-Historie. Die weitere, breitere Klasse
  (Zahl falsch, nicht zwingend im Suchlauf-Mechanismus) ist mit 23 Belegen bereits ein Deckel-Eintrag.
  Einschätzung: kein **neuer** Registereintrag fällig — der bestehende Eintrag
  `zahl-in-traeger-driftet-gegen-die-messung` ist die passende Adresse für eine 24. Evidence-Datei;
  der engere Rahmen „ausgerechnet der §3.13-Suchlauf, der genau das verhindern soll, versagt an sich
  selbst" ist noch keine dritte Instanz und braucht (noch) keinen eigenen Eintrag. Nur benannt, nicht
  selbst angelegt (Auftrag).

## 7. Entscheidungs-Konformität

- **Nutzer-Entscheidung „Weg 3":** Fragment trägt exakt `python python3`, keine weiteren Namen (`go`,
  `gofmt`, `node`, `dotnet`, `java`, `gradle`, `uv` bleiben ungesperrt — 6 eigene Pass-Proben und
  `MR-004` §Grenz-Zeile bestätigen das). Konform.
- **`MR-003`:** unverändert (`Accepted`, immutable), `make doc-immutable` Exit 0. Konform.
- **`MR-004`:** vollständiges Feld-Gerüst (Datum/Geltungsbereich/Ersetzt-Baseline-Regel/Adaption/
  Begründung/Grenz-Zeile/Auflösungs-Trigger), Anker `guard-haertung` löst real auf. Konform.
- **`AGENTS.md` §3.6:** keine ADR nötig (Härtung, keine Gate-Lockerung) — kein Gate berührt. Konform.
- **`AGENTS.md` §3.7:** Kopfkommentare im Indikativ, höchstens ein `MR`-Zeiger je Satz, keine
  Konjunktiv-Prosa über verworfene Alternativen. Konform.
- **`AGENTS.md` §3.9:** `make gates` in diesem Lauf ungefiltert, Exit-Code direkt geprüft. Konform.
- **§3.12/§3.13:** Zahlen tragen ihren Ursprung, der Plan-Text ist nach der Fixrunde korrekt; die eine
  gefundene Ausnahme (V-1, Register-Fließtext „367") ist ein Nachbar-Träger außerhalb des vom Suchlauf
  abgedeckten Musters, keine Verletzung dieses Slice-Plans selbst.
- **§3.1 Docker-only:** kein Host-Interpreter, keine Umleitung im finalen Diff (F-3 betrifft nur einen
  zurückgenommenen Zwischenschritt ohne Commit).

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **Register-Fließtext nennt eine veraltete Zahl:** `inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md:23` sagt „Tabellentest … mit 367 Fällen" — seit dem CRLF-Fix (368. Fall) veraltet. War beim Schreiben korrekt (vor der Fixrunde), die Fixrunde hat es nicht nachgezogen, weil der §3.13-Suchlauf keine Zahlen trifft. Kein DoD-Bruch (die Zahl steht nicht im Plan selbst), aber eine reale kleine Drift | `docs/plan/planning/observations/BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md:23` | ja — Lesen + `make test-command-guard` |
| V-2 | LOW | **F-1 selbst wurde nicht als Register-Evidence nachgetragen:** die im Review gefundene Instanz „Zahl im Träger schon beim Commit falsch" passt in die bestehende Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (23×, Deckel), hat dort aber keine 24. Evidence-Datei bekommen; auch keine neue, engere Klasse „Suchlauf-Mechanismus versagt an sich selbst" wurde angelegt (wäre erst beim dritten Auftreten fällig, dieses ist das erste laut eigener Suche, §6) | §6 | ja — `ls evidence`, `grep` über `docs/reviews/` |
| V-3 | INFO | **3 der 8 eigenen Mutationen (B/C/D) sind neu ggü. Auftrag/Review/Fixrunde** (gemischte Zeilenenden, Leerzeilen am Ende, Kommentar-Token) — alle grün: der CRLF-Fix ist minimal, trifft nur `\r`, und hat keine Nebenwirkung auf Nachbarfälle | §4b | ja — eigener Lauf |
| V-4 | INFO | **Der Prozessschritt „Review-Checkbox im Fixrunden-Commit selbst setzen"** (`eb04d172`) ist der in `implement-slice.md` Schritt 21 ausdrücklich vorgesehene Weg, keine Abkürzung — die unabhängige Bestätigung liefert regulär diese Verifikation | `.claude/commands/implement-slice.md:269-272` | ja — Lesen |

Kein HIGH, kein MEDIUM in dieser Verifikation.

## 9. Verdikt

**DoD bestätigt:** ja, in der Substanz — alle zehn `[x]`-Zeilen tragen einen realen, von mir
unabhängig nachgemessenen Beleg; die drei `[ ]`-Zeilen sind korrekt der Planner-Closure vorbehalten.
**Plan-vs-Code:** keine unbenannte Abweichung; Guard-Code-Diff beschränkt sich exakt auf die zwei im
Plan zugesagten Änderungen (§3.1-Wege in `REASON_PKG`, CRLF-Bereinigung der Fragment-Ladung).
**Entscheidungs-Konformität:** Nutzer-Entscheidung „Weg 3" exakt umgesetzt (nur `python`/`python3`),
`MR-003` unverändert, `MR-004` vollständig und formkonform, kein Gate gelockert.
**Review-Findings:** F-1 (HIGH) und F-2 (MEDIUM) beide unabhängig nachgemessen und bestätigt aufgelöst
(14/14 Suchlauf-Zeilen, CRLF-Fix hält gegen 4 eigene Regressionsproben); F-3/F-4 (INFO) bestätigt ohne
Aktionsbedarf. **Mutationen:** 10 eigene Hook-Proben + 8 eigene Fragment-Mutationen, alle wie zugesagt,
3 davon eigenständig über den bisherigen Prüfumfang hinaus. **Gates:** `make gates`,
`make commit-traceability RANGE=cea198fb..HEAD`, `make doc-immutable RANGE=cea198fb..HEAD`,
`make suchlauf-nachmessen`, `make test-command-guard` — alle Exit 0 im eigenen Lauf.

Zwei LOW-Findings (V-1, V-2) betreffen ausschließlich Register-/Nachbar-Träger außerhalb des
Plan-Texts selbst und sind keine DoD-Verletzung dieses Slice — sie sind Übergabe-Punkte an die
Planner-Closure.

### Übergabe an den Planner

1. **DoD-Haken:** die zehn `[x]`-Zeilen sind belegt (Tabelle §2); die drei `[ ]`-Zeilen (Closure-Notiz,
   Risiken-Ausgänge, drei Paarungen) bleiben regulär bei der Closure zu setzen.
2. **Register-Nachzug (V-1):** `inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md:23` von „367"
   auf „368 Fälle" korrigieren (oder allgemeiner formulieren, um künftige Drift durch weitere
   Testfälle zu vermeiden).
3. **Register-Nachzug (V-2):** eine 24. Evidence-Datei für `BEO-PGC/zahl-in-traeger-driftet-gegen-die-
   messung` anlegen, die F-1 dieses Reviews referenziert (Zahl im Träger bereits beim Commit falsch,
   nicht erst durch spätere Drift) — passt in die bestehende, bereits gedeckelte Klasse; kein neuer,
   engerer Eintrag „Suchlauf-Mechanismus versagt an sich selbst" nötig, da dies die erste beobachtete
   Instanz dieser engen Formulierung ist (eigene Suche über `docs/reviews/*.md`, §6).
4. **Risiken §6:** alle sieben Zeilen können auf einen Ausgang gesetzt werden — Risiko 1
   (Falsch-positive Blockaden: eingetreten wie gewollt, `python3 --version` blockt, Pass-Fälle bleiben
   frei, durch meine 10 Hook-Proben bestätigt), Risiko 2 (Tabellentest-Bindung an die Eingabe:
   entfallen/erfüllt, 8 eigene Mutationen plus die 6 des Reviewers bestätigen es), Risiko 3 (legitimer
   Python-Bedarf: nicht eingetreten, kein Fall ohne Weg berichtet), Risiko 4 (Wirkung auf laufende
   Rollen-Aufträge: eingetreten in dokumentierter Form — die vier `--version`-Vorfälle aus Berichten
   in §1 des Plans, kein Repo-Text nennt `python3` als zulässigen Weg), Risiko 5 (toter Code am
   `python|python[0-9]*|perl`-Zweig: entfallen, `python3.12`/`python2`-Fälle halten ihn), Risiko 6
   (Sperre verspricht mehr als sie kann: entfallen/erfüllt, `MR-004` §Grenz-Zeile benennt es), Risiko 7
   (neue Fälle laufen nicht: entfallen, `make test-command-guard` druckt 368 und lief real).
5. **Drei Paarungen:** bleiben an der Slice-Closure selbst hängen (Plan §7, korrekt).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
