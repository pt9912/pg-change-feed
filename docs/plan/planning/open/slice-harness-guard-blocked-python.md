# Slice harness-guard-blocked-python: Der PreToolUse-Guard sperrt Host-`python` und `python3` am Kopf eines Kommando-Segments unbedingt — Fragment `tools/harness/blocked/python`, Härtung als `MR-004`, Tabellentest `make test-command-guard`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — Harness-Querschnitt: der Slice trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre. Er berührt Guard und
Tabellentest, nicht den Runner und nicht `internal/`; eine technische Kante zu
einem Slice der [welle-transformationen](../welle-transformationen.md) hat er nicht.
Die empfohlene Position steht in §4.

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.1 (Docker-only, Verbot des
Host-Interpreters auf einer Repo-Datei, Absatz „Durchsetzung“), §3.6 (die Sperre
verschärft, sie lockert kein Gate: keine ADR nötig), §3.7 (Kommentare), §3.9,
§3.12 und §3.13;
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft der
Aussagen im Plan; Grenze: kein Sensor über Dateiinhalte);
`MR-003` ([`harness/conventions.md`](../../../../harness/conventions.md), Anker
`mr-003`: Vertrag des Guards, Grenz-Zeile, Auflösungs-Trigger der Kopf-Liste);
Beobachtungs-Register `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`
(`state.md`: Adresse dieses Slice) und
`BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`; Baseline-Regelwerk
`modul-13-quality-gates.md` §Guard-Härtung (jede Härtung landet als neuer `MR-<NNN>`
mit Grenz-Zeile).

**Herkunft der Entscheidung:** Nutzer-Entscheidung „Weg 3“ in der Sitzung
(2026-09-27); Frage und Wege stehen in der Closure-Notiz von
`slice-transformationen-map-value` (§7, Punkt „Frage an den Nutzer“).

**Berührte Spec-Stellen:** — (Harness-Wächter; keine Spec-Stelle).

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Auftrag des Auftraggebers (Nutzer-Entscheidung „Weg 3“).
**Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`) blockt
`python` und `python3` am Kopf eines Kommando-Segments **unbedingt** — auch hinter
`cd <Repo> &&`, auch auf einem Pfad ohne Repo-Namen, auch als `python3 --version`.
Die Wortliste steht in einem Fragment `tools/harness/blocked/python`, das der Guard
über seine bestehende Fragment-Ladung liest; die Härtung landet als `MR-004`, und
`make test-command-guard` bindet die Zusage an ihre Eingabe.

**Ausgangslage — Beleg-Anker je Aussage.**

- *Ladung (gelesen, Stand `cea198fb`):* der Guard liest jede Datei unter
  `tools/harness/blocked/*` als Wortliste und hängt sie an `BLOCKED` an
  (`.claude/hooks/pretooluse-command-guard.sh`, Zeilen 79–90); der Boden trägt
  15 Paketmanager-Namen (gezählt). Der Kopf eines Segments wird nach den Übersprüngen
  (Zuweisungen, Wrapper samt Optionen, `\`, Schlüsselwörter, Basename) gegen `BLOCKED`
  gelesen. Das Verzeichnis `tools/harness/blocked/` existiert nicht (`ls tools/harness/blocked`
  meldet „nicht gefunden“, gemessen am 2026-09-27, Stand `cea198fb`).
- *Heute (gelesen in `MR-003`, Adaption, und im Guard, Zeile 251):* `python`,
  `python3`, `python3.<N>` und `perl` blocken nur, wenn der **ganze Befehlsstring** einen
  Repo-Pfad nennt. Die Grenz-Zeile von `MR-003` nennt zwei Ränder, die das offen lässt: `cd <Repo> && python3 x.py`
  und ein Aufruf auf einem Pfad ohne Repo-Namen. Beide sind an der Hook-Schnittstelle des
  echten Guards gemessen (Hook-JSON auf stdin, Stand `cea198fb`): `python3 --version` und
  `cd <Repo-Wurzel> && python3 --version` enden mit Exit 0 ohne Ausgabe, also Pass.
- *Register (gezählt, `ls evidence | wc -l` am 2026-09-27):* `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`
  trägt 5 Beleg-Dateien; die Neubewertung der Kopf-Liste hat ihr zweites Trigger-Kriterium
  seit der Closure von `slice-transformationen-map-value` erfüllt (`MR-003`, Auflösungs-Trigger).
- *Aufrufe seit der Frage (**übernommen** aus Agenten-Berichten, Angabe des Auftraggebers;
  nicht am Guard gemessen):* ein Planner-Aufruf `cd <Repo> && python3 --version` passierte den
  Guard, ein Planner-Aufruf `python3 --version` wurde geblockt, ein Planner-Aufruf `python3` mit
  leerem Heredoc bei der Slice-Anlage wurde geblockt, ein Architect-Aufruf `python3 --version`
  ([`ADR-0129`](../../adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md), kein Repo-Pfad) passierte den Guard,
  ein Reviewer-Aufruf `python3 --version` im selben Vorgang wurde geblockt, ein Implementer-Aufruf `python3 -c 1` in
  `slice-leerlauf-phase-last-in-stuecken` trägt im Bericht keinen Guard-Ausgang. Keiner hatte
  Wirkung auf eine Repo-Datei. Das Register trägt sie als Belege (§3, Träger-Tabelle).
- *Entscheidung:* der Nutzer hat „Weg 3“ gewählt — die kleinere Liste, nur `python python3`,
  „die einzigen Interpreter mit Belegen“ (Wortlaut der Wege: `slice-transformationen-map-value`
  §7). Die übrigen Namen der ursprünglichen Liste (`go gofmt node dotnet java gradle uv`)
  bleiben ungesperrt, weil für sie kein Beleg vorliegt.
- *Wirkung auf den Bestand des Tabellentests (gemessen, Stand `cea198fb`):* der Tabellentest zählt
  317 Fälle, 311 davon Tabellenzeilen (gedruckt „alle 317 Fälle bestanden“;
  `grep -cE '^(block|pass) '` über `tools/harness/run-command-guard-tests.sh`). Gegen eine Kopie des
  Guards im Scratchpad, in der `BLOCKED` um `python python3` erweitert ist (Emulation des Fragments,
  `GUARD=<Kopie> bash tools/harness/run-command-guard-tests.sh`), färben **52** Fälle rot
  (`grep -c '^FEHLER:'`): 24 Block-Fälle der Klasse `interp` blocken jetzt mit der Klasse `pkg`,
  27 Pass-Fälle blocken, und der Fall „Meldung interp: Weg nach AGENTS.md 3.1“ findet eine andere
  Meldung. Die 52 sind Fälle der Gruppen „Host-Interpreter auf Repo-Pfaden“, „Pfadzeichen davor/dahinter“,
  „Falsch-Positiv-Ränder“ und „Grenzen“; sie sind der Umbau-Umfang des Tabellentests (DoD 2).

**Kernaussage der Wirkung (hergeleitet aus dem Guard-Code, erprobt im Tabellentest, DoD 2):** die
Sperre gilt für den Kopf des Segments. `make …` und `docker run … python3 …` bleiben frei (Kopf ist
`make` bzw. `docker`), ebenso `command -v python3` und jedes `python3` als Argument, in
Anführungszeichen oder als Teil eines Suchmusters. `python3 --version` blockt jetzt — ausdrücklich
gewollt (Nutzer-Entscheidung).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die übrigen Namen der ursprünglichen Liste** (`go`, `gofmt`, `node`, `dotnet`, `java`, `gradle`,
  `uv`). Für sie liegt kein Beleg vor (der Plan von `slice-harness-guard-inplace-textwerkzeug` §1
  nennt die drei Vorfälle ohne Host-`go`-Aufruf, **übernommen**; für die übrigen Namen ist der Bestand
  nicht gesucht); der Nutzer hat sie ausgenommen. Die Neubewertung hat einen Trigger in `MR-004`
  (Auflösungs-Trigger); die Ergänzung der Liste ist eine Nutzer-Entscheidung.
- **`python3.<N>`, `python2`, `pypy`, `ipython`, `uv run python`.** Der Guard vergleicht das Kopf-Wort
  ganz gegen die Liste; die Nutzer-Entscheidung nennt `python` und `python3`. `python3.<N>` bleibt unter
  der Repo-Pfad-Regel von `MR-003` (blockt nur mit einem Repo-Pfad im Befehlsstring). Die Namen gehören
  als benannte Grenze in `MR-004` und als Pass-Fälle in den Tabellentest.
- **`perl`.** Die Repo-Pfad-Regel und die in-place-Form (`perl -pi`) bleiben unverändert; `perl` steht
  nicht in der Liste (Nutzer-Entscheidung, kein Beleg).
- **Umleitungen und flaglose Schreibwege, ein Skript, das ein anderer Interpreter liest, ein `cd`
  vor einem nicht gelisteten Kopf.** Grenz-Zeile von `MR-003`, unverändert; ein Sensor über Dateiinhalte
  ist ausgeschlossen ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md), Grenze; ein
  Werkzeugaufruf hinterlässt in der Datei keine Signatur).
- **Eine Scratchpad-Ausnahme oder eine Ausnahmeliste.** Für `sed -i` entschieden (keine Ausnahme,
  `slice-harness-guard-inplace-textwerkzeug` §6, Frage 2); dieselbe Entscheidung gilt hier. Der zulässige
  Weg steht in `AGENTS.md` §3.1: Edit/Write, ein Repo-Werkzeug hinter `make` (gepinntes Docker-Image),
  `sed … Datei > Kopie` nach stdout.
- **Eine ADR und die Aufnahme in `make gates`.** Ein Gate braucht eine ADR (`AGENTS.md` §3.6, §4); ein
  Wächter gehört nicht in die Gate-Tabelle (Baseline `modul-13-quality-gates.md` §Guard-Härtung: er
  verhindert eine Handlung). Die Sperre verschärft; `make test-command-guard` bleibt ein Werkzeug ohne
  Gate-Bindung.
- **Ein Live-Aufruf mit Host-`python` als Beleg.** Er verstieße gegen `AGENTS.md` §3.1; der Beleg läuft
  über die Hook-Schnittstelle des echten Guards (DoD 1).
- **Eine Änderung der Regel in `AGENTS.md` §3.1 selbst.** Die Regel steht (Verbotssatz, Host-Werkzeug-Klasse);
  nachgezogen wird nur der Absatz „Durchsetzung“.

## 2. Definition of Done

Jedes Kriterium trägt „Zu belegen durch:“; jede Aussage über eine Mutation ist eine
**Erwartung**, bis der Implementer sie gefahren hat (Stelle, Instanz, gesehene Farbe;
[`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).

- [ ] **Liefer-Punkt 1 — die Sperre.** Das Fragment `tools/harness/blocked/python` trägt genau die
      zwei Wörter `python python3` in einer Zeile (Name nach der Konvention des Guards,
      `blocked/<sprache>`, Kopfkommentar Zeile 80; der Guard liest jede Datei des Verzeichnisses, der
      Name ist keine Vorgabe des Guards). Der Guard blockt ein Kommando-Segment, dessen Kopf —
      nach den Übersprüngen des Bestands (Zuweisungen; Wrapper `sudo env command exec nice time xargs
      eval busybox` samt Optionen; führendes `\`; absoluter Pfad, Basename; Schlüsselwörter; der Rest hinter
      `eval`, `bash -c`, `-exec`/`-execdir`/`-ok`/`-okdir`) — `python` oder `python3` ist, **unabhängig vom
      Befehlsstring**: auch hinter `cd <Repo> &&`, auf einem Pfad ohne Repo-Namen, bei `--version`.
      Ausgabe wie im Bestand (`"decision": "block"`, Exit 0, gültiges JSON, Klasse `pkg`). Die Blockmeldung,
      die ein Rollenlauf liest, nennt die Wege nach `AGENTS.md` §3.1 (Edit/Write, Repo-Werkzeug hinter
      `make`, `sed … > Kopie`); trägt die bestehende `pkg`-Meldung sie nicht, ergänzt der Implementer sie
      (kein `"` und kein `\` im Text, damit die Ausgabe gültiges JSON bleibt). Frei bleiben: Kopf `make`,
      `docker`, `git`, `grep`, `echo`; `command -v|-V python3`; `python3` als Argument, in
      Anführungszeichen, als Musterteil, als Pfadteil; jeder nicht gelistete Name. Der Code des Guards
      ändert sich nicht, außer der Beleg verlangt es (Erwartung: die Fragment-Ladung trägt); der Kopfkommentar
      und die Grenz-Zeile im Kopfkommentar tragen den neuen Stand im Indikativ, höchstens eine Kennung je
      Kommentar (`AGENTS.md` §3.7). *Zu belegen durch:* der Tabellentest (Liefer-Punkt 2) und ein
      **Beleg an der Hook-Schnittstelle des echten Guards mit dem echten Fragment**: das Hook-JSON für
      `python3 --version` auf der stdin von `bash .claude/hooks/pretooluse-command-guard.sh` (kein Aufruf
      von Host-`python`) endet mit der Block-Ausgabe, deren Wortlaut im Bericht steht; das Hook-JSON für
      `make test` endet ohne Ausgabe.
- [ ] **Liefer-Punkt 2 — der Tabellentest.** `tools/harness/run-command-guard-tests.sh` legt die Fragmente
      der Quelle `BLOCKED_DIR` (Vorgabe: `tools/harness/blocked` des Repos, übersteuerbar für Mutationsläufe
      an Kopien, wie `GUARD` und `MASKER`) in das Wegwerf-Repo (`tools/harness/blocked/`) und bindet:
      (a) **Block-Fälle je Position** für `python` und `python3` (Kopf, nach `&&`, hinter `cd <Repo> &&`
      und `cd /tmp/x &&`, nach `;`, `||`, Pipe, Zeilenumbruch, in `$(…)` und Backticks, in `( … )`,
      Zuweisungs-Präfix, `sudo`, `env`, `env -i`, `time -p`, `xargs`, absoluter Pfad `/usr/bin/python3`,
      `\python3`, `bash -c`, `eval`, `find -exec`, `for … do`, `if … then`, Heredoc mit `python3 -` und
      `python3 -c`); (b) **Pass-Fälle neben jedem Block-Fall**: Kopf `make`, `docker run … python3 …`,
      `command -v python3`, `command -V python`, `type`/`which`, `echo python3`, `git grep -n python3`,
      `grep -E 'a|python3' f`, `git commit -m "… python3 …"`, `ls /usr/bin/python3`, ein Pfad mit
      `python` als Verzeichnisname (`cat sdks/python/pyproject.toml`), `python3-config` und `ipython` als
      nicht gelistete Köpfe; (c) **benannte Falsch-Positiv-Ränder mit erwartetem Block**: `python3 --version`
      (ausdrücklich gewollt), `python3 -c 'print(1)'`, eine Heredoc-Zeile, die mit `python3` beginnt;
      (d) **benannte Grenzen mit erwartetem Pass** — die Nutzer-Entscheidung „Weg 3“ gebunden: je ein Fall für
      `go`, `gofmt`, `node`, `dotnet`, `java`, `gradle`, `uv` und für `python3.12` und `python2` ohne
      Repo-Pfad; `perl -e 1` ohne Repo-Pfad; (e) **die 52 Fälle des Bestands** (Ausgangslage) sind je einem
      Ziel zugeordnet, keiner still gestrichen: entweder auf den Kopf `python3.12` umgestellt, wo sie die
      Repo-Pfad-Regel und ihre Zeichenklassen binden (der Zweig `python|python[0-9]*|perl` des Guards bleibt
      für `python3.<N>` und `perl` in Kraft, Zeile 251), oder auf die neue Erwartung (Block, Klasse `pkg`)
      umgestellt, oder mit Grund gestrichen; die Zahl der Fälle vorher (317, gemessen) und nachher steht im
      Bericht. *Zu belegen durch:* `make test-command-guard` Exit 0 mit gedruckter Zahl; **Mutationen an
      Kopien** (`GUARD=<Kopie>`, `BLOCKED_DIR=<Kopie>`), je Zeile Stelle, Instanz und gesehene Farbe im
      Bericht — *erwartet, zu erproben*: Fragment fehlt (leeres `BLOCKED_DIR`) → die Block-Fälle rot, auch die
      auf einem Repo-Pfad (die Repo-Pfad-Regel fängt sie mit der Klasse `interp` weiter, die Klassen-Prüfung
      des Tests färbt sie trotzdem rot); Fragment leer (Datei vorhanden, 0 Byte) → dieselben rot; Liste um
      `python3` gekürzt → nur die `python3`-Fälle rot; um `python` gekürzt → nur die `python`-Fälle rot;
      Liste um `go` erweitert → der Pass-Fall `go version` rot; Kopf-Erkennung geändert (Basename am Kopf
      entfernt → `/usr/bin/python3` rot; `env` aus den Wrapper-Präfixen → `env python3` rot; `bash -c`-Rekursion
      entfernt → `bash -c "python3 x"` rot; `command -v`-Ausnahme entfernt → `command -v python3` rot; `&&` nicht
      als Trenner → der Fall hinter `cd <Repo> &&` rot); Fragment-Ladung liest nur eine feste Datei `go` → die
      Block-Fälle rot. Menge der Erprobung: die Fälle des Tabellentests.
- [ ] **Liefer-Punkt 3 — die Träger.** (a) `harness/conventions/MR-004-guard-host-python-am-kopf.md` per `cp` aus
      `.harness/baseline/v6.9.0/templates/harness/conventions/MR-NNN-titel.template.md` (byte-Gleichheit mit
      `diff -q` vor dem Füllen), in place gefüllt, plus die Zeile in `harness/conventions.md` §Aktive
      Adaptionen (Anker `mr-004`); Inhalt siehe §3 „Ansatz“ (Grenz-Zeile: die Liste ist die benannte Grenze).
      `MR-003` bleibt unverändert (`Accepted`, immutable). (b) [`AGENTS.md`](../../../../AGENTS.md) §3.1 Absatz
      „Durchsetzung“ nennt, was der Guard jetzt zusätzlich blockt (`python`/`python3` am Kopf, unbedingt), was
      unter der Repo-Pfad-Regel bleibt (`python3.<N>`, `perl`), und dass die übrigen Toolchain-Namen ungelesen
      bleiben; Herkunft `· seit slice-harness-guard-blocked-python`, Verweis auf `MR-004` neben `MR-003`.
      (c) Die Zeile `make test-command-guard` in [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Fragment, `BLOCKED_DIR`, die neuen Gruppen, Bindung `MR-004`), die Hilfe-Zeile des Ziels im
      `Makefile` (`BLOCKED_DIR`), die Kopfkommentare von Guard und Tabellentest im Indikativ. (d) Das Register:
      `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md` (Ausgang der Kopf-Liste: verkörpert durch
      Fragment, `MR-004`, Tabellentest; die Nennung „`tools/harness/blocked/go`“ trägt den Verweis auf das
      Fragment `blocked/python`), `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration/state.md`
      (Zeiger auf die Durchsetzung für `python`/`python3`). *Zu belegen durch:* Lesen der Stellen,
      `make docs-check` (Links, Anker), `make doc-immutable RANGE=<Parent>..HEAD` Exit 0 (`MR-003` unverändert)
      und der Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/open/slice-harness-guard-blocked-python.md`
      läuft nach jeder Fixrunde mit den `diff`-Zeilen des Implementers durch
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors (Liefer-Punkt 3c);
      das Benutzerhandbuch bleibt unberührt (keine Betreiber-Oberfläche).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang der Kopf-Liste in
      `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (Liefer-Punkt 3d); je ein Anfall
      im Lauf dieses Slice ist eine weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine Antwort und
      wird in §7 notiert (§8 nennt die Einträge, die dieser Slice trägt).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** S — Schätzung, nicht gemessen: ein Fragment mit einer Zeile, ein Kopfkommentar im Guard, ein
Tabellentest-Umbau von 52 bestehenden Fällen plus die neuen Gruppen (a) bis (d), ein `MR`, vier kurze
Träger-Änderungen und zwei Register-Zeilen. Der Tabellentest-Umbau ist der größte Teil.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/blocked/python` | neu | Liefer-Punkt 1: die Wortliste `python python3`, eine Zeile; der Guard liest sie über seine bestehende Ladung |
| `.claude/hooks/pretooluse-command-guard.sh` | update (Kommentar; Code nur, wenn der Beleg es verlangt) | Kopfkommentar Zeilen 3–4 („Host-python/-perl auf Repo-Pfaden“), 39–42 („Host-Interpreter: … blockt, wenn der GANZE Befehlsstring …“), 46–55 („ein Skript, das ein Interpreter liest (`python3 x.py` mit Text ohne Repo-Namen)“) und 79–83 (Fragmente) tragen den neuen Stand im Indikativ; die `pkg`-Meldung nennt die Wege nach §3.1, falls sie sie nicht trägt |
| `tools/harness/run-command-guard-tests.sh` | update | Liefer-Punkt 2: `BLOCKED_DIR`, Fragment-Lieferung ins Wegwerf-Repo, neue Gruppen, Umbau der 52 Fälle, Kopfkommentar (Gruppen, Übersteuerung) |
| `Makefile` | update | Hilfe-Zeile von `test-command-guard` nennt `BLOCKED_DIR` neben `GUARD`/`MASKER` |
| `harness/conventions/MR-004-guard-host-python-am-kopf.md` | neu, per `cp` | Liefer-Punkt 3a |
| `harness/conventions.md` | update | Zeile `MR-004` in §Aktive Adaptionen (Anker `mr-004`) |
| `AGENTS.md` §3.1 | update | Liefer-Punkt 3b |
| `harness/README.md` §Sensors | update | Liefer-Punkt 3c |
| `docs/plan/planning/observations/BEO-PGC/…/state.md` (zwei Einträge) | update | Liefer-Punkt 3d |

**Ansatz (Liste):**

- **Fragment und Ladung.** Der Guard lädt `tools/harness/blocked/*` bereits; das Fragment ist eine Datei mit einer
  Zeile. Ein fehlendes oder leeres Verzeichnis lässt den Boden unberührt (der Guard ist nie fail-open, Kopfkommentar
  Zeile 81–82); ohne das Fragment fiele `python3` still auf die Repo-Pfad-Regel zurück. Deshalb liest der
  Tabellentest das Fragment **aus dem Repo** und liefert es ins Wegwerf-Repo: das Löschen oder Leeren der Datei
  färbt ihn rot.
- **Klasse und Meldung.** Ein Treffer der Wortliste trägt die `pkg`-Klasse (der Guard setzt `REASON` erst
  bei den in-place- und Interpreter-Regeln). Der Tabellentest prüft je Block die Klasse; deshalb färbt die Mutation
  „Fragment fehlt“ auch die Fälle rot, die die Repo-Pfad-Regel weiter mit der Klasse `interp` blockt (Doppelschutz
  verdeckt die Mutation nicht; *hergeleitet*, im Lauf zu bestätigen).
- **Der Zweig `python|python[0-9]*|perl` (Zeile 251) bleibt.** `python`/`python3` erreichen ihn nie mehr; `python3.<N>`
  und `perl` erreichen ihn. Der Kommentar dort nennt das, damit der Zweig nicht als toter Code gilt; die Fälle
  mit Kopf `python3.12` halten ihn.
- **`MR-004` (Inhalt, der Implementer füllt aus der Vorlage).** *Titel:* der Guard sperrt Host-`python`/`python3` am Kopf
  unbedingt. *Datum:* das des Slice. *Geltungsbereich:* das Fragment, der Guard, der Tabellentest, `AGENTS.md` §3.1
  „Durchsetzung“. *Ersetzt-Baseline-Regel:* genau eine — der Punkt „Gehärtet wird die Zerlegung, nicht die Denylist“
  in [`modul-13-quality-gates.md` §Guard-Härtung](../../../../.harness/baseline/v6.9.0/regelwerk/modul-13-quality-gates.md#guard-haertung)
  (Anker `guard-haertung` vom Implementer mit `make docs-check` zu prüfen). *Adaption:* der Guard sperrt zwei
  Kopf-Wörter unbedingt; die Repo-Pfad-Regel von `MR-003` bleibt für `python3.<N>` und `perl`. *Begründung:* die
  Baseline nennt als Schaden der Denylist, dass sie `make` blockiert; die Kopf-Liste des Guards liest nur den Kopf, `make`
  und `docker` bleiben frei (*hergeleitet* aus dem Guard, *belegt* durch die Pass-Fälle von Liefer-Punkt 2). Die zwei Ränder
  der Grenz-Zeile von `MR-003` schließt eine Zerlegung nicht: ein Pfad ohne Repo-Namen ist keine Zerlegungs-Frage, ein
  `cd` im selben Kommando ist Shell-Zustand hinter der Obergrenze der Quote-Lesung (`MR-003`, Adaption); die unbedingte
  Liste schließt beide. Auslöser: der Auflösungs-Trigger von `MR-003` (zweites Kriterium) und die Nutzer-Entscheidung
  „Weg 3“; der Register-Eintrag mit seinen Beleg-Dateien (Zahl aus dem `state.md`, übernommen). *Grenz-Zeile — was der Guard
  nicht kann:* jeder nicht gelistete Name (`go`, `gofmt`, `node`, `dotnet`, `java`, `gradle`, `uv`, `python3.<N>`, `python2`,
  `pypy`, `ipython`, `uv run python`); ein `python3` aus einer Variablen oder einem Alias; Umleitungen und flaglose
  Schreibwege; ein Skript, das ein anderer Interpreter liest; die Rezepte hinter `make` und die Docker-Bauten (der Guard
  liest die Bash-Aufrufe des Laufs); die Heredoc-Zeile, die mit `python3` beginnt, blockt (Falsch-Positiv, Quote-Lesung
  ohne Heredoc-Kenntnis); `python3 --version` blockt (gewollt). *Schärft* `MR-003` in genau diesen Sätzen: „`cd <Repo> && python3 x.py` geht
  durch“, „ein Skript, das der Interpreter liest (`python3 /tmp/x.py`) geht durch“ und der Falsch-Positiv-Rand „ein
  `python3`-Aufruf, dessen Text einen Repo-Namen nennt, blockt, auch wenn er nichts schreibt“ gelten für `python`/`python3`
  nicht mehr. *Auflösungs-Trigger:* permanent; die Neubewertung der nicht gelisteten Namen folgt einer weiteren Beleg-Datei
  des Register-Eintrags, die einen Host-Aufruf eines nicht gelisteten Namens nennt; die Ergänzung der Liste ist eine
  Nutzer-Entscheidung.
- **Der zulässige Weg für einen legitimen Python-Bedarf:** Edit/Write, ein Repo-Werkzeug hinter `make` (gepinntes Image,
  `AGENTS.md` §3.1), `sed … Datei > Kopie` nach stdout. Ein Bedarf, der in keinen Weg passt, ist ein neues
  `make`-Ziel (eigener Slice) oder eine Beleg-Datei im Register, keine Ausnahme im Guard.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „was der PreToolUse-Guard für Host-`python` sperrt und
was er nicht liest“).** Suchraum: der ganze Baum ohne die drei Ausnahmen von
[`AGENTS.md`](../../../../AGENTS.md) §3.13 (`docs/reviews/**`, Records unter `done/`, `.harness/baseline/**`), keine
weitere Einschränkung; die Plan-Datei schließt das Werkzeug aus. Stand ist `cea198fb` (der Commit vor der Anlage
dieses Slice, vom Planner am 2026-09-27 gemessen; ein Stand ist eine Commit-Kennung, nie `HEAD`); der Implementer
misst am Parent seiner Arbeit neu und trägt die `diff`-Zeilen ein (`make suchlauf-nachmessen
PLAN=docs/plan/planning/open/slice-harness-guard-blocked-python.md`). Die Zahlen des Standes `cea198fb` sind mit
`git grep -n` gemessen, nicht übernommen. Drei Arten des Musters: der Symbolname (Fragment und Kopf-Liste), die
Beschreibung (Host-`python`, Host-Interpreter, Sprach-Toolchains) samt dem Hedge des offenen Punkts, das Zählwort
(die Namensliste):

```suchlauf
cea198fb 5 -n -E 'blocked/go|blocked/python|Kopf-Liste|Host-Toolchain-Sperre' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cea198fb 37 -n -E 'Host-.?python|Host python|Host-Interpreter|Sprach-Toolchains' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cea198fb 2 -n -E 'gofmt python python3|go gofmt|node dotnet java gradle uv' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cea198fb 2 -n -E 'Offene Nutzer-Entscheidung|Entscheidung des Nutzers steht|Neubewertung der Kopf|Re-Evaluierung der Host' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cea198fb 5 -n -E 'Skript, das (ein|der) Interpreter liest|jeder andere Interpreter|andere Interpreter' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cea198fb 6 -n -E 'python3? --version' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
cea198fb 2 -n -i -E '\.claude/hooks|pretooluse|harness/conventions' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/in-progress docs/plan/planning/welle-transformationen.md
```

| Träger | Befund (Stand `cea198fb`, vom Planner gelesen) | Behandlung |
|---|---|---|
| `.claude/hooks/pretooluse-command-guard.sh`, Kopfkommentar | Zeilen 3–4, 39–42, 46–55 nennen Host-`python` nur auf Repo-Pfaden und „ein Skript, das ein Interpreter liest (`python3 x.py` mit Text ohne Repo-Namen)“ als ungelesen | Liefer-Punkt 1 und 3c; Implementer zieht nach |
| `tools/harness/run-command-guard-tests.sh`, Kopfkommentar und 52 Fälle | die Gruppe „Host-Interpreter auf Repo-Pfaden“ und ihre Nachbargruppen tragen `python3` als Kopf; Zeilen 14–16 beschreiben die Gruppen | Liefer-Punkt 2 (Umbau) und 3c (Kommentar) |
| [`AGENTS.md`](../../../../AGENTS.md) §3.1 Absatz „Durchsetzung“ (Zeilen 122–129) | „einen Host-`python`-/`perl`-Aufruf, dessen Befehlsstring einen Repo-Pfad nennt“, „ein Skript, das ein Interpreter liest“ und „Sprach-Toolchains (`go`, `gofmt`, …)“ als ungelesen, „Ein Host-`python`/`perl` auf einem Pfad ohne Repo-Namen … passiert den Guard“ | Liefer-Punkt 3b; Zeilen 75, 88, 94 (die Regel, die Klasse) bleiben wahr |
| [`harness/README.md`](../../../../harness/README.md) §Sensors, Zeile `make test-command-guard` | „Host-`python`/`perl` auf Repo-Pfaden“, Übersteuerung `GUARD`/`MASKER` | Liefer-Punkt 3c |
| `harness/conventions/MR-003-guard-inplace-textwerkzeug.md` | die Grenz-Zeile nennt `cd <Repo> && python3 x.py` und `python3 /tmp/x.py` als Durchgang, der Auflösungs-Trigger nennt `tools/harness/blocked/go` | **unverändert** (`Accepted`, immutable, `make doc-immutable`); `MR-004` schärft und nennt die betroffenen Sätze |
| `.claude/commands/implement-slice.md` (Zeile 30), `.harness/skills/reviewer.md` (Zeile 236) | nennen „Host-Interpreter“ als Teil des Verbots und verweisen auf `AGENTS.md` §3.1 | bleiben wahr; Implementer liest, ändert nicht |
| **Fremde Datei:** `docs/plan/planning/open/slice-code-kommentare-bereinigung.md` §6 | „der PreToolUse-Guard blockt … einen Host-Interpreter auf Repo-Pfaden, die Grenze steht in `MR-003`“ | gemeldet, nicht mitgeändert; Frist: die Closure dieses Slice, der Planner der Closure zieht nach oder benennt den Träger mit Adresse |
| Beobachtungs-Register `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md` | nennt `tools/harness/blocked/go` als offene Nutzer-Entscheidung | Nachzug des Planners bei der Anlage dieses Slice (Ausgang der Frage, Adresse dieser Slice); Ausgang „verkörpert“ trägt die Closure |
| Beobachtungs-Register `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration/state.md` | nennt die Durchsetzung für Host-Werkzeuge außerhalb der Klasse nur als Reviewer-HIGH | Nachzug des Planners bei der Anlage (Zeiger); Ausgang trägt die Closure |
| `docs/plan/planning/welle-transformationen.md` §5 (a) | nennt die Guard-Grenze für die Umleitung in `MR-003` als „unverändert“ | bleibt wahr (die Umleitung liest der Guard weiter nicht); Implementer liest, ändert nicht |
| Beschreibung in `done/` (`slice-transformationen-map-value` §7, `slice-harness-guard-inplace-textwerkzeug` §1, §6) | Records der Frage und ihrer Wege | bleiben stehen (Record-Einfrierung); die Entscheidung steht in diesem Plan und im Register |

**Nicht durch Suchlauf fassbar:** ein Zeilen-Lokator auf den Guard oder auf `AGENTS.md` §3.1 — der Suchlauf trifft
Symbolnamen, keine Zahlen ([`AGENTS.md`](../../../../AGENTS.md) §3.13, Grenze); der Reviewer liest sie gegen den Diff.

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1; beobachtbar:
`ls docs/plan/planning/in-progress` nennt nur `roadmap.md`). Eine Nutzer-Entscheidung ist nicht mehr abzuwarten: „Weg 3“
liegt vor (Herkunft oben). Die Änderung am Guard wirkt in jeder Sitzung, die das Repo öffnet, ab dem Zeitpunkt, an dem das
Fragment im Arbeitsbaum liegt (der Guard liest das Verzeichnis bei jedem Aufruf, gelesen, Zeilen 85–90); ein zweiter
Implementer-Lauf in derselben Zeit liefe unter einem Guard, der sich unter ihm ändert — das WIP-Limit deckt das.

**Reihenfolge (Empfehlung an den Orchestrator):** nach `slice-leerlauf-phase-last-in-stuecken` (er liegt in
`in-progress/`; das WIP-Limit ordnet das) und **vor** `slice-wal-fehlerschwelle-ausgangsklasse`. Eine technische Kante gibt es
nicht: der Slice berührt weder den Runner noch `internal/bootstrap/wiring.go` (die Dateien der Runner-Slices), und kein
anderer offener Plan ändert `.claude/hooks/`, den Guard oder `harness/conventions` (die letzte Suchlauf-Zeile oben trifft
zwei Beschreibungen — `slice-code-kommentare-bereinigung` §6 und `welle-transformationen` §5 (a) —, keine Änderung, Träger-Tabelle).
Gründe für „vor“ statt „nach“, ohne Bindung: jede weitere Sitzung ohne die Sperre ist ein weiterer möglicher
Anfall der Klasse (die Aufrufe seit der Frage stehen in Berichten von vier Rollen — Planner, Architect, Reviewer, Implementer —,
übernommen), die drei folgenden Runner-Slices
laufen alle unter dem Guard und ihre Implementer, Reviewer und Verifier fahren Mutationen an Kopien — die Umstellung auf Edit/Write
am kleinen Slice zuerst ist die billigste; der Preis ist ein S-Slice Verzögerung auf dem Weg der Welle (nicht gemessen). Wählt der
Orchestrator „nach `slice-wal-fehlerschwelle-ausgangsklasse`“, ändert sich kein Plan: keine Kante, kein Start-Trigger eines
anderen Plans nennt diesen Slice, und das WIP-Limit ordnet den Rest.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht zu erwarten (Umfang S); trägt der Umbau der 52 Fälle
  mehr als erwartet, ist der Umbau der Repo-Pfad-Fälle auf den Kopf `python3.12` der abtrennbare Teil.
- `in-progress` → `open` (blockiert): ein Alltags-Aufruf der Rollen, der nach den Regeln blockt, obwohl er zulässig ist (kein
  zulässiger Weg im Sinn von `AGENTS.md` §3.1 fällt unter den Kopf `python`/`python3`; ein Beispiel wäre ein `make`-Rezept, das der
  Guard als Kopf `python3` läse) — dann Architect-Frage, keine Ausnahme im Guard.

## 5. Closure-Trigger

DoD vollständig (die Sperre, der Tabellentest mit gesehenen Mutationen und dem Umbau der 52 Fälle, die Träger nachgezogen) +
`make gates` grün + Review-Report ohne offenes HIGH oder MEDIUM + Verifikation, dass die DoD trägt (`make test-command-guard`
grün, die Mutationen rot gesehen, der Hook-Beleg mit dem Wortlaut der Blockmeldung im Bericht) + Closure-Notiz mit
Lerneintrag geschrieben (geschärfte Regel: [`AGENTS.md`](../../../../AGENTS.md) §3.1 Durchsetzung und `MR-004`; neuer Sensor:
keiner — `make test-command-guard` wird erweitert, ohne Gate).

## 6. Risiken und offene Punkte

- **1. Falsch-positive Blockaden harmloser Aufrufe.** `python3 --version` blockt jetzt (gewollt); eine Heredoc-Zeile, die mit
  `python3` beginnt, blockt (die Quote-Lesung liest keine Heredocs, `MR-003` Obergrenze); ein `python3` als Argument, in
  Anführungszeichen, als Musterteil oder als Pfadteil blockt nicht (der Kopf entscheidet). *Erwartet, zu belegen durch:* die
  Pass-Fälle neben jedem Block-Fall und die benannten Ränder von Liefer-Punkt 2 (b), (c). **Ausgang:** *(bei Closure)*
- **2. Die Tabellentest-Bindung färbt nur an der Ausgabeseite** (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  verkörpert). Der Test liest das Fragment aus dem Repo; die Eingabeseite ist die Wortliste (fehlt, leer, gekürzt,
  erweitert), die Ladung und die Kopf-Erkennung. *Erwartet, zu belegen durch:* je eine Mutation, deren Farbe der Implementer
  **gesehen** und im Bericht genannt hat (Stelle, Instanz, Farbe; Liste in DoD 2). **Ausgang:** *(bei Closure)*
- **3. Ein Rollenlauf braucht legitim Python.** Zulässiger Weg: Edit/Write, ein Repo-Werkzeug hinter `make`, `sed … > Kopie`
  (`AGENTS.md` §3.1); `docker run … python3 …` passiert den Guard (Kopf `docker`), ist aber nur mit einem gepinnten Image
  Weg nach §3.1 — der Slice führt keinen neuen Weg ein. Ein Bedarf ohne Weg ist eine Beleg-Datei im Register oder ein `make`-Ziel,
  keine Ausnahme im Guard. *Erwartet, zu belegen durch:* die Blockmeldung nennt die Wege (DoD 1); ein Review-Lauf und ein
  Verifikations-Lauf unter dem Guard melden, ob ein Bedarf ohne Weg auftrat. **Ausgang:** *(bei Closure)*
- **4. Wirkung auf laufende Rollen-Aufträge und Sitzungen.** Der Guard wirkt ab dem Fragment im Arbeitsbaum, nicht erst ab dem
  Commit; ein Auftrag, der `python3` nennt, läuft in einen Block. Die Repo-Texte verweisen auf `AGENTS.md` §3.1 statt den
  `--version`-Satz zu wiederholen (`.claude/commands/implement-slice.md` Zeile 30 und `.harness/skills/reviewer.md`
  Zeile 236, gelesen). *Erwartet, zu belegen durch:* das WIP-Limit (§4) und die Träger-Tabelle (§3): kein Repo-Text nennt
  `python3` als zulässigen Weg. **Ausgang:** *(bei Closure)*
- **5. Der Zweig der Repo-Pfad-Regel für `python`/`python3` wird als toter Code gestrichen.** Er bleibt für `python3.<N>` und
  `perl` in Kraft. *Erwartet, zu belegen durch:* die Fälle mit Kopf `python3.12` und ihre Zeichenklassen (Liefer-Punkt 2 e)
  und der Kommentar am Zweig. **Ausgang:** *(bei Closure)*
- **6. Die Sperre verspricht mehr, als sie kann** (`BEO-PGC/regel-weiter-als-ihr-sensor`, verkörpert). Nicht gelesen:
  jeder nicht gelistete Name, ein Name aus einer Variablen oder einem Alias, Umleitungen, ein Skript eines anderen
  Interpreters. *Erwartet, zu belegen durch:* die Grenz-Zeile von `MR-004`, die Kurzform im Kopfkommentar und in
  `AGENTS.md` §3.1, die benannten Grenz-Fälle im Tabellentest (Liefer-Punkt 2 d); der Reviewer liest die drei Stellen
  gegeneinander. **Ausgang:** *(bei Closure)*
- **7. Die neuen Fälle laufen nicht** (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen durch:* die
  gedruckte Zahl von `make test-command-guard` (Instanz A, `AGENTS.md` §3.12) und die Farbe der Mutationen. **Ausgang:**
  *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure trägt sie
  selbst, nach dem `git mv` nach `done/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`/`PGC`
(Greenfield); `.claude/hooks`, `tools/harness` und `harness/` sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler gemessen am
2026-09-27 mit `ls evidence | wc -l` je Eintrag, Stand `cea198fb`) —
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 5×: der Gegenstand; Adresse der
Kopf-Liste, Risiken 1, 3, 4, 6; der Nachzug des Planners bei der Anlage fügt Beleg-Dateien hinzu),
`BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (verkörpert, 2×: Zeiger,
Liefer-Punkt 3d), `BEO-PGC/regel-weiter-als-ihr-sensor` (verkörpert, 4×, Risiko 6),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (19×, Deckel; Risiko 2, die Mutationen je Eingabeseite),
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (11×, Deckel: DoD und Ausgangslage nennen ihren
Beleg-Anker oder sind als Erwartung formuliert), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (23×, Deckel: die
Zahlen tragen Ursprung und Stand), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (33×, Deckel: Suchlauf-Feld in §3),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, Risiko 7),
`BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (7×) und `BEO-PGC/kommentar-herkunft-als-kette` (2×): die
Kopfkommentare von Guard und Tabellentest sind Kommentare in Skripten (Indikativ, höchstens eine Kennung, keine
Slice-Nummer); übrige Einträge gesichtet, kein Bezug.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
