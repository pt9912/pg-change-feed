# Review-Report: slice-harness-guard-inplace-textwerkzeug — 2026-09-26

**Review-Art:** Code — der Diff verschärft den PreToolUse-Guard (`.claude/hooks/pretooluse-command-guard.sh`)
um die in-place Textwerkzeuge und Host-`python`/`perl` auf Repo-Pfaden, liefert den Tabellentest
`make test-command-guard`, den Adaptions-Eintrag `MR-003`, den Absatz „Durchsetzung“ in `AGENTS.md` §3.1,
eine Zeile in `harness/README.md` und `harness/conventions.md` sowie die Schritt-20-Ergänzung in
`.claude/commands/implement-slice.md`; geprüft gegen Plan, `MR-003`, `AGENTS.md` §3.1/§3.6/§3.7/§3.9/§3.12/§3.13,
die Hard Rules und `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung (Modul 10 §Drei Review-Arten).
Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-harness-guard-inplace-textwerkzeug` (ohne Welle, Harness-Querschnitt), Diff-Range
`e98d419c..a0472421` (6 Commits, 9 Dateien, +559/−32; Baum sauber, nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ in der Form des Diff-Stands
`a0472421`. **Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-guard-inplace-textwerkzeug` (§1 Ziel und Abgrenzung, §2 Liefer-Punkte als Bezug,
  §3 Plan, Suchlauf-Feld, Träger-Tabelle und Abweichungen, §4 Rückführungen, §6 Risiken 1–7 und die zwei
  Nutzer-Fragen); Register-Eintrag `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (`observation.md`,
  `state.md`, vier Beleg-Dateien)
- [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen in Trägern)
- `harness/conventions/MR-003-guard-inplace-textwerkzeug.md`, `harness/conventions.md` (MR-Block, Index-Zeile),
  `MR-000`/`MR-001`/`MR-002`
- `AGENTS.md` (Hard Rules §3.1 mit dem Absatz „Durchsetzung“, §3.6, §3.7, §3.9, §3.12, §3.13, §4)
- Baseline `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung, `regelwerk/grundlagen-durchsetzungsschicht.md`
  §Grenzen — ehrlich benannt, `templates/harness/conventions/MR-NNN-titel.template.md`
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-harness-fmt-check.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; Scratchpad der Sitzung; Mutationen
ausschließlich an Kopien, erzeugt mit `sed` ohne `-i` nach stdout — der Guard war dabei in der Sitzung aktiv):

- **Läufe am Stand `a0472421`** (einer zugleich; dangling Volumes `docker volume ls -qf dangling=true | wc -l`
  34 vor und 34 nach allen Läufen; kein `prune`; `git status --short` nach den Läufen leer):
  `make test-command-guard` Exit 0, gedruckt „run-command-guard-tests: alle 123 Fälle bestanden“;
  `make docs-check` Exit 0, „1278 Datei(en) geprüft, 0 Befund(e)“; `make gates` (einmal) Exit 0, „gesamt: 0
  Befund(e)“; `make commit-traceability RANGE=e98d419c..a0472421` Exit 0, „6 Commit(s) … Betreffs ohne
  Struktur-ID“; `make doc-immutable RANGE=e98d419c..a0472421` Exit 0, „0 Befund(e)“;
  `make kommentar-kennungen DIFF=87491cca COUNT=1` Exit 0, Ausgabe **0** (das Werkzeug liest nur Go; die
  Shell-Kommentare von Hand: unten).
- **Menge der Fälle:** `grep -cE '^(block|pass) '` der Tabelle 119, dazu vier Fälle über `run_raw`/PATH-Umbau
  (abgeschnittenes JSON, leere Eingabe, `\u`-Escape, fehlendes `awk`) — 123, wie „alle 123 Fälle“ im Lauf.
- **51 rote Fälle gegen den Guard am Parent:** `git show e98d419c:.claude/hooks/pretooluse-command-guard.sh` in
  eine Scratchpad-Datei, `GUARD=<Datei> bash tools/harness/run-command-guard-tests.sh` Exit 1, 51 Zeilen
  `FEHLER:` von 123 — die Zahl des Implementers stimmt. Darunter zwei Fälle mit der Bezeichnung „Bestand“
  (`ls | xargs -n1 pip`, `env -i pip install x`), die der Guard am Parent nicht blockt (F-3).
- **Eingabeseiten-Mutationen am Guard** (41 Kopien, je `GUARD=<Kopie> bash tools/harness/run-command-guard-tests.sh`,
  Exit und Zahl der `FEHLER:`-Zeilen; das Original blieb unverändert):

  | Nr. | Mutation an der Kopie | Ergebnis |
  |---|---|---|
  | M1 | sed-Bündel-Regel entfernt | rot, 22 Fälle |
  | M2 | sed nur nacktes `-i` | rot, 3 |
  | M3 | jedes `sed` blockt | rot, 12 |
  | M4 | Anführungszeichen-Token als Flag gelesen | rot, 2 |
  | M5 | `-exec`-Kopf nicht gelesen | rot, 4 |
  | M6 | perl-Bündel mit beliebigem `i` | rot, 1 |
  | M7 | awk `-i` ohne `inplace` | rot, 3 |
  | M8 | Block-Ausgabe entfernt | rot, 66 |
  | M9 | Host-Interpreter: nur das Segment statt des ganzen Befehls | rot, 4 |
  | M10 | Pfadzeichen-Bedingung entfernt | rot, 5 |
  | M11 | `./` nicht erlaubt | rot, 1 |
  | M12 | Datei-Namen der obersten Ebene nicht gelesen | rot, 3 |
  | M13 | Wrapper-Optionen nicht übersprungen | rot, 4 |
  | M14 | defektes JSON nicht fail-closed | rot, 3 |
  | M15 | Tiefe 4 erlaubt | rot, 1 |
  | M16 | `sudo` kein Präfix | rot, 2 |
  | M17 | `perl` nicht im Interpreter-Kopf | rot, 2 |
  | M18 | in-place-Block mit dem Grund der Klasse `interp` | rot, 36 |
  | M19 / M20 | `-iinplace` bzw. `--include=inplace` entfernt | rot, je 1 |
  | M21 | `python3.N` nicht erkannt | rot, 1 |
  | M29 / M36 / M38 / M39 | perl ohne Ziffern · Pfadklasse ohne Ziffern · Datei-Name-Nachbedingung · find-Ende | rot, je 1 |
  | M22–M26 | sed-Klasse `[nEsrzu]` ohne `s`, `r`, `z`, `u` (einzeln; M22 zusammen: `[nE]`) | **grün** (F-1) |
  | M27, M28, M30–M32 | perl-Klasse `[0-9lanpsw]` ohne `l a n s w` zusammen (M27) bzw. ohne `w`, `l`, `a`, `s` einzeln | **grün** (F-1) |
  | M33–M35, M37 | Pfadzeichen-Klasse ohne `~`, `-`, `_` bzw. ohne Großbuchstaben | **grün** (F-1) |
  | M-A | Prüfung `command -v awk` allein entfernt | grün — äquivalent (Doppelschutz durch den Exit-Code der `awk`-Extraktion; beide Zeilen zusammen entfernt: die vier Fail-closed-Fälle rot) |

  Die 21 vom Plan genannten Mutationen (Regel entfernt, nacktes `-i`, jedes `sed`, Anführungszeichen-Bereinigung,
  `-exec`, perl-Bündel, awk-Wert, Block-Ausgabe, Segment statt Befehl, Pfadzeichen, `./`, Datei-Name) sind an den
  genannten Fällen rot gesehen.
- **Alltags- und Umgehungsformen** (rund 165 Kommandostrings über einen kleinen Treiber gegen den Guard am
  Stand `a0472421` und gegen den Guard am Parent; Ausgabe je Zeile „BLOCK-<Klasse>“ oder „pass“). Blocken
  (gewollt): `sed -ni`, `sed -E -i`, `sed -e … -i f`, `sed … -i f`, `sed -s -i`, `sed -z -i`, `sed -i --`,
  `command sed -i`, `env`/`nice`/`time`/`sudo -E`/`xargs -0`/`xargs -P4 -n1 sed -i`, `(sed -i …)`,
  `$(sed -i …)`, Backticks, `eval "sed -i …"`, `sh`/`dash`/`zsh`/`bash -c`, `/bin/sed -i`, Tabulator, `{ sed -i …; }`,
  `find … | xargs sed -i`, `perl -p -i -e`, `perl -Mstrict -pi`, `perl -CS -pi`, `perl -wpi`, `awk -v x=1 -i inplace`,
  `awk -iinplace`, `python3 ./docs/x` — alle wie zugesagt. Passieren: die benannten Grenzen der Tabelle sowie die
  Formen unter F-2, F-4 und F-6.
- **Fehler- und Ausgabeform:** je Klasse eine Meldung mit dem Ersatzweg; `jq` liest die Ausgabe aller drei Klassen
  als gültiges JSON (`decision` `block`, `reason` ohne `"` und `\`), Exit 0 wie im Bestand; im Pass-Fall keine Ausgabe;
  Latenz je Aufruf rund 11 ms (zehn Läufe 0,117 s, der Guard am Parent 0,110 s).
- **Live-Wirkung in der Sitzung des Reviewers:** der Hook lieferte in dieser Sitzung viermal einen Block —
  zweimal die Klasse `inplace`, zweimal `interp` — je bei einem **lesenden** `git grep`/`grep` (F-2); die
  Ausgabeform (Meldung mit Ersatzweg) kommt im Werkzeug-Ergebnis an. Der Wortlaut steht im Guard
  (`git grep -n 'rewrite repo files without a trace' -- .claude tools harness`: ein Treffer,
  `pretooluse-command-guard.sh:55`) und stimmt mit der gelieferten Meldung überein. Einen Live-Aufruf mit
  `python3` habe ich nicht gefahren (§3.1); die Klasse `interp` ist an der Hook-Schnittstelle durch die 11
  Treffer- und 13 Nicht-Treffer-Fälle des Tabellentests belegt, deren Mutationen M9–M12, M17, M21 rot färben.
- **Suchlauf-Feld:** `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-guard-inplace-textwerkzeug.md`
  Exit 0, gedruckt „suchlauf-nachmessen: 15 Zeilen stimmen“. Von Hand mit `git grep` nachgefahren (Ist gleich
  Soll): Zeile 1 (`89d427e0`, `pretooluse-command-guard`) 3; Zeile 2 (`89d427e0`, Beschreibung) 6; Zeile 3
  (`89d427e0`, Zählwort, mit `-e` je Alternative) 9; Zeile 4 (`89d427e0`, „in keinem committeten Text“) 1;
  Zeile 5 (`89d427e0`, Register-Kennung in `open`/`next`/`welle-transformationen`) 3; Zeile 6 0;
  Stand `diff` Zeile 1 11, Zeile 3 54, Zeile 7 (`test-command-guard`) 7 — jeweils ohne die Plan-Datei.
  Selbst gesucht (Nichtgefundenes): kein weiterer Träger unter `docs/user/`, `spec/`, `.github/`,
  `harness/sensors/`, `harness/targets/`, `.claude/agents/`, der den Guard beschreibt; die Träger-Meldungen des
  Plans (`welle-transformationen` §6 (a), `state.md` des Register-Eintrags) treffen zu.
- **Register-Zahlen:** `ls evidence` des Eintrags nennt vier Dateien; `MR-003` und Plan zählen vier (fünf
  `sed -i`-Läufe in drei Vorgängen, ein Host-Python-Heredoc, zwei Host-`python3`-Aufrufe im vierten Vorgang) —
  stimmt gegen die Beleg-Dateien.
- **Hard Rules:** `git diff e98d419c..a0472421 | grep` auf `sed -i`/`perl -pi`/`awk -i` trifft nur Text, Tests und
  Kommentare (kein Aufruf im Diff); Shell-Kommentare des Diffs von Hand auf Konjunktiv/Chronik gelesen
  (`grep -nE` auf `wäre|waere|würde|wuerde|sonst|statt|früher|zuvor|jetzt|Slice|welle` über die hinzugefügten
  `#`-Zeilen: keine Treffer); Kopfkommentar und Testkopf im Indikativ, keine Slice-Nummer; die Lifecycle-Commits
  `05273900` und `87491cca` sind reine Renames (`git show -M --stat`: 0 Zeilen); das Benutzerhandbuch bleibt
  unberührt und muss es (keine Betreiber-Oberfläche).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Die Zeichenklassen der Regeln sind nicht an ihre Eingabe gebunden: das sed-Bündel `-[nEsrzu]*i` ist nur mit `n` und `E` belegt (`-ni`, `-Ei`), das perl-Bündel `-[0-9lanpsw]*i` nur mit `p` und Ziffern (`-pi`, `-0777pi`), die Pfadzeichen-Klasse `[A-Za-z0-9_./~-]` nur mit `/`, `.`, Ziffern; die Mutationen auf `s`, `r`, `z`, `u`, `l`, `a`, `s`, `w`, `~`, `-`, `_` und die Großbuchstaben einzeln sowie auf `n` gemeinsam mit `l a s w` (M22–M28, M30–M35, M37) lassen alle 123 Fälle grün. `MR-003` sagt, der Tabellentest binde „jede Zusage an ihre Eingabe“; Plan und `MR-003` nennen die Klassen wörtlich als Zusage. | Reviewer-Skill HIGH „Zusage ohne Bindung an ihre Eingabeseite“; `AGENTS.md` §3.12 Instanz B (Aussage über eine Menge nennt die Menge, an der sie geprüft ist) | `tools/harness/run-command-guard-tests.sh:138-149,205-225`; `.claude/hooks/pretooluse-command-guard.sh:108,111,145,149`; `harness/conventions/MR-003-guard-inplace-textwerkzeug.md:46-49` | ja — `sed` ohne `-i` nach stdout erzeugt die Kopie mit gekürzter Klasse, `GUARD=<Kopie> make test-command-guard` Exit 0 | Zusage ohne Bindung an ihre Eingabeseite |
| F-2 | MEDIUM | Ein `\|` in einem Such-Muster startet ein Segment, dessen Kopf die in-place Form oder `python`/`perl` trägt, und der Guard blockt einen lesenden Aufruf: `git grep -E 'sed -i\|perl -pi\|awk -i' -- .` (das Muster der Suchlauf-Zeile 3 des Plans, das Zählwort der bewegten Eigenschaft; das mittlere Glied `perl -pi` trägt kein Anführungszeichen), `grep -E 'a\|perl -pi\|b'`, `git grep -E 'sed -i\|awk -i inplace\|x'`; für die Klasse `interp` blockt `grep -E 'sed -i\|python' <Pfad unter docs/>`, sobald ein Repo-Pfad irgendwo im Befehl steht. Die Tabelle bindet nur zweigliedrige Muster (Nicht-Treffer) und den Rand mit Leerzeichen vor dem Anführungszeichen; das mittlere und das letzte Glied stehen weder als Rand noch in `MR-003`. Der Aufruf war in dieser Sitzung viermal live geblockt (Handnachfahren der Suchlauf-Zeile 3, ein `grep` im Register, ein `git grep` auf den Skill, ein `grep` auf `MR-003` mit `python3` als letztem Glied). Damit tritt die Rückführung §4 (a) des Plans ein („Alltags-Aufruf der Rollen, zulässig, geblockt; der Weg dahin verlangt Quote-Bewusstsein“): Architect-Frage. | Plan §4 Rückführung (a), §6 Risiko 1; `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung (Denylist um den Interpreter „blockiert legitime Shell-Arbeit“) | `.claude/hooks/pretooluse-command-guard.sh:164-173,189-205`; `tools/harness/run-command-guard-tests.sh:194-196,231` | ja — `bash`-Treiber mit dem Guard und den vier Mustern (Ausgabe „BLOCK-inplace“); live: der Aufruf selbst | Falsch-Positiv-Rand ungebunden und nicht benannt |
| F-3 | MEDIUM | Die Wrapper-Optionen werden für **alle** Klassen übersprungen, nicht nur für in-place: `command -v pip`, `command -v npm`, `command -v yarn` (die Existenzprüfung eines Werkzeugs) blocken jetzt als Paketmanager, am Parent passieren sie; ebenso `time -p pip`. Der Plan (Abweichungen (a)) nennt die Erweiterung nur für `xargs -r sed -i`; DoD und `MR-003` sagen „Bestandsregeln bleiben unverändert“, und zwei Fälle der Gruppe „Bestand“ (`ls \| xargs -n1 pip`, `env -i pip install x`) sind am Parent grün, also keine Bestandsfälle. | Plan §2 Liefer-Punkt 1 („Bestandsregeln … unverändert“); `AGENTS.md` §3.12 Instanz B | `.claude/hooks/pretooluse-command-guard.sh:179-183`; `tools/harness/run-command-guard-tests.sh:113-114`; `harness/conventions/MR-003-guard-inplace-textwerkzeug.md:46-48` | ja — `git show e98d419c:…` in eine Scratchpad-Datei, Treiber mit `command -v pip` gegen beide Guards | Bestandsregel verändert, als unverändert ausgewiesen |
| F-4 | MEDIUM | Kopf-Wörter der Shell-Kontrollstrukturen werden nicht übersprungen: `for f in a b; do sed -i s/a/b/ "$f"; done`, `while read f; do sed -i …; done`, `if …; then sed -i …; fi`, `… ; else sed -i …`, `! sed -i …`, `case x in x) sed -i …;; esac` passieren (Kopf ist `do`, `then`, `else`, `!`, `x)`). Das ist die typische Form einer Massenänderung mit `sed -i`; sie steht in keiner Grenz-Zeile (`MR-003`, Kopfkommentar, `AGENTS.md` §3.1) und in keinem Tabellenfall (die Brace-Group ist behandelt, `do`/`then` nicht). Das Verhalten ist Bestand für Paketmanager, gilt aber jetzt als Zusage für die in-place Klasse. | `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung (die Grenz-Zeile wird mitgezogen; eine Lücke, die nicht dasteht, ist eine Harness-Lüge); Plan §6 Risiko 2 | `.claude/hooks/pretooluse-command-guard.sh:176-188`; `harness/conventions/MR-003-guard-inplace-textwerkzeug.md:51-75` | ja — Treiber mit den Zeilen gegen den Guard (alle „pass“) | Grenz-Zeile nennt eine Lücke der Zusage nicht |
| F-5 | MEDIUM | Die Blockmeldung der Klasse `interp` und `MR-003` (Begründung) verweisen auf `python3 /path/to/scratch/mutate.py /path/to/scratch/copy` als zulässigen Weg; `AGENTS.md` §3.1 führt `python` unter den Werkzeugen, die im Container laufen (Zeile 88), und nennt für die Mutation auf der Kopie nur Edit/Write und `sed … Datei > Kopie` (Zeile 98-100). Der Plan behauptet ohne Anker, die Python-Ersetzung auf einer Kopie sei „der zulässige Weg nach §3.1“. Die Meldung erscheint in jeder Sitzung und lenkt zu einem Host-Interpreter-Aufruf, den die Regel im selben Slice nicht deckt. | `AGENTS.md` §3.1; Reviewer-Skill MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“; `AGENTS.md` §3.12 Instanz B | `.claude/hooks/pretooluse-command-guard.sh:56`; `harness/conventions/MR-003-guard-inplace-textwerkzeug.md:77-82`; `AGENTS.md:88,98-101`; Plan §1 (Ausgangslage, Zeile 62) | ja — Lesen der vier Stellen gegeneinander | Meldung und Regel nennen verschiedene zulässige Wege |
| F-6 | LOW | Weitere Formen, die der Guard nicht liest und die keine Grenz-Zeile nennt: `sed -i''` und `perl -i'' -pe` (das Token mit Anführungszeichen ist kein Flag; die Grenz-Zeile nennt nur `sed '-i'`), `/usr/bin/env python3 tools/x.py` und `/usr/bin/sudo …` (der Wrapper mit absolutem Pfad), `busybox sed -i`, `\sed -i`, `gsed -i`, die Abkürzung `sed --in-p`, `find … -exec sh -c 'sed -i …' {} \;`, ein `bash -c` mit maskierten Anführungszeichen in der zweiten Ebene, `uv run python`/`node -e`/`ruby -e` auf Repo-Pfaden (nur `python`/`perl` sind Interpreter der Regel); umgekehrt liest die Flag-Suche bis zum Segmentende: `perl /tmp/x.pl -input a` und `sed -n 1p -input.txt` blocken (Skript- und Dateiargumente `-i…`). Die Grenz-Zeile steht außerdem nicht „gleichlautend“ in den drei Trägern, wie `MR-003` sagt: `AGENTS.md` nennt „Sprach-Toolchains“, die `MR-003` nicht führt; der Kopfkommentar führt „Flag in Anführungszeichen“ nicht als Grenze. | `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung (Grenz-Zeile mitziehen); `AGENTS.md` §3.12 Instanz B | `harness/conventions/MR-003-guard-inplace-textwerkzeug.md:51-75`; `.claude/hooks/pretooluse-command-guard.sh:35-43`; `AGENTS.md:117-121` | ja — Treiber mit den Zeilen | Grenz-Zeile unvollständig oder ungleichlautend |
| F-7 | LOW | Der Konjunktiv-Kandidatenlauf in Schritt 20 (`git diff -U0 <Basis> -- '*.go' \| grep -nE '^\+.*//…(wäre\|würde\|…)'`) liest nur Go-Kommentare mit `//` und Umlaut-Schreibweise; der Skopus des Reviewer-Skills für die Klasse „Kommentar trägt keine der Kommentar-Klassen“ umfasst Skripte, Tests und Runner (`tools/harness/*.sh`, `.claude/hooks/*.sh`), deren Kommentare in diesem Repo transliteriert sind (`waere`, `entkaeme`). Der Bestand trägt eine solche Stelle (`pretooluse-command-guard.sh:184-185`, „Kopf waere sonst `{`“, nicht im Diff, kein Gegenstand dieses Reviews). | Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen“ (Skopus); `AGENTS.md` §3.7 | `.claude/commands/implement-slice.md:243-248` (Diff `c625b572`) | ja — `git grep -n -E 'waere\|entkaeme' -- .claude/hooks` (ein Treffer) gegen den Befehl des Schritts (kein Treffer) | Selbstprüfung enger als der Prüfpunkt |
| F-8 | INFO | Baseline-Bezug für den Architect: (a) `MR-003` trägt für `perl -i` und `awk -i inplace` keinen Beleg (`MR-003` und Plan §6 Risiko 5 sagen es); die Baseline nennt eine Wächter-Regel ohne Sensor-Evidenz „Aufwand ohne Begründung“ — beide Regeln kosten je eine Erkennung und sind an der Tabelle gebunden (bis auf F-1). (b) Liefer-Punkt 2 (Host-`python`/`perl` mit Repo-Pfad-Muster über den ganzen Befehlsstring) ist der Interpreter-Denylist näher als der „Zerlegung“, die die Baseline als Härtungsrichtung nennt; `MR-003` begründet die Abgrenzung, F-2 zeigt die Falsch-Positiv-Fläche. (c) Das Feld „Ersetzt-Baseline-Regel“ nennt einen Satz, der nach dem Eintrag weiter wahr bleibt („die Grenze rückt, sie fällt nicht“); die Vorlage verlangt genau eine Regel, die der Eintrag ersetzt. | `v6.9.0` · `regelwerk/modul-13-quality-gates.md` §Guard-Härtung; `regelwerk/grundlagen-durchsetzungsschicht.md` §Grenzen — ehrlich benannt | `harness/conventions/MR-003-guard-inplace-textwerkzeug.md:10-24,77-82` | nein | Wächter-Regel ohne Evidenz / Interpreter-Denylist |
| F-9 | INFO | Der Fall „fehlendes awk (fail-closed)“ wird von einer Einzelmutation nicht rot (M-A: die Zeile `command -v awk` allein entfernen): der Exit-Code der `awk`-Extraktion blockt dann dieselbe Eingabe; erst beide Zeilen zusammen färben ihn und die drei Parse-Fälle. Die Doppelsicherung ist gewollt; die Aussage des Implementers („Einzelmutation äquivalent“) stimmt. Die Live-Wirkung des Hooks ist in dieser Sitzung an beiden Klassen `inplace` und `interp` gesehen (unbeabsichtigt, F-2); der Tabellentest trägt die Hook-Schnittstelle und genügt für Liefer-Punkt 2 ohne Live-Aufruf mit `python3`. | Maintainability | `.claude/hooks/pretooluse-command-guard.sh:226,232`; `tools/harness/run-command-guard-tests.sh:131-135` | ja — M-A | Äquivalente Mutation durch Doppelschutz |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.claude/hooks/pretooluse-command-guard.sh` — Erkennung der drei Klassen `pkg`/`inplace`/`interp`, Positionen (Kopf, `&&`, `;`, Pipe, `xargs`, `sudo`, Zuweisung, absoluter Werkzeugpfad, `bash -c` Tiefe 3, `find -exec`/`-execdir`/`-ok`/`-okdir` mit Ende `+`/`\`/`;`), rohe Flag-Tokens, fail-closed | geprüft, ohne Befund über F-1 bis F-4, F-6 hinaus: 21 Mutationen der Zusagen rot (M1–M21), die genannten Umgehungen `sed -ni`, `sed -E -i`, `-e … -i`, Flag hinter Argument, `command`/`env`/`nice`/`time`/`eval`, Sub-Shells, Backticks, `$(…)`, Brace-Group blocken; `sed -n`, `sed --debug`/`--version`, `awk` lesend, `perl -e`/`-MList::Util`/`-I`, `find -iname`/`-exec grep -i`, `xargs -I{}`, `cp -i`, `grep -i`, `git commit -F`, `git commit -m "… sed -i …"`, `git log --grep='sed -i'` passieren; das Ausgabeformat ist gültiges JSON der Bestandsform, Exit 0, im Pass-Fall keine Ausgabe |
| `.claude/hooks/pretooluse-command-guard.sh` — Kopfkommentar | geprüft, ohne Befund über F-6 hinaus: Indikativ, eine Herkunftszeile (`AGENTS.md`), keine Slice-/Wellen-Nummer, Zusage · Abgrenzung · Grenze tragen; die Zusagen (`-[nEsrzu]*i`, `-[0-9lanpsw]*i`, `-i inplace`, Wrapper-Optionen, Repo-Pfad-Muster) stimmen mit dem Code an den Zeilen 93–158 überein; kein Konjunktiv über eine verworfene Alternative in den hinzugefügten Zeilen |
| `tools/harness/run-command-guard-tests.sh` und `make test-command-guard` | geprüft, ohne Befund über F-1, F-3, F-9 hinaus: Wegwerf-Repo im Temp-Verzeichnis, schreibt nur dorthin (`trap` räumt auf), `GUARD=` übersteuert den Prüfling, je Block prüft der Test Begründung der Klasse und JSON-Form, Nicht-Treffer neben jedem Treffer, benannte Ränder und Grenzen als Fälle, 123 Fälle gezählt, 51 rot gegen den Guard am Parent nachgefahren; Kopf nennt die Host-Werkzeuge (coreutils-Basis); Makefile-Ziel mit `.PHONY` und Hilfe-Zeile, kein Eintrag in `GATE_CHECKS` |
| `harness/conventions/MR-003-…md`, `harness/conventions.md` (Index-Zeile, Anker `mr-003`) | geprüft, ohne Befund über F-5, F-6, F-8 hinaus: per `cp` aus der Vorlage, Pflichtfelder Datum · Geltungsbereich · Ersetzt-Baseline-Regel (Link mit Anker, `make docs-check` löst ihn) · Adaption · Begründung · Auflösungs-Trigger vorhanden, Grenz-Zeile führt die Punkte aus Plan §6 Risiko 2, Auslöser stimmt mit dem Register (vier Beleg-Dateien, fünf `sed -i`-Läufe); append-only (`make doc-immutable` Exit 0), kein früherer Wächter-`MR`, den er schärfen könnte |
| `AGENTS.md` §3.1 (Mutationsprobe, „Durchsetzung“) | geprüft, ohne Befund über F-5, F-6 hinaus: Aussagen zu `sed -i`/`perl -i`/`awk -i inplace`, unabhängig vom Ziel, hinter `find -exec`, Host-`python`/`perl` mit Repo-Pfad und die genannten Nicht-Lesungen sind wahr gegen den Guard (Treiber); der Satz zur Mutationsprobe nennt den Stdout-Weg, und `sed … Datei > Kopie` passiert den Guard; Herkunft ein Feld (`· seit slice-harness-guard-inplace-textwerkzeug`) |
| `harness/README.md` (Zeile `make test-command-guard`) | geprüft, ohne Befund: „kein Gate“ mit Begründung (Wächter verhindert eine Handlung), Verweis auf `MR-003`, Host-Werkzeuge genannt, Aufruf `GUARD=<Datei>` |
| `.claude/commands/implement-slice.md` Schritt 20 (`c625b572`) | geprüft, ohne Befund über F-7 hinaus: Urteilsregel „Zusage im Indikativ zulässig / verworfene Alternative umformulieren“, Ausnahmen (Mutationsbeschreibung in Test-Godocs, normale Zweige) stimmen mit dem Reviewer-Skill überein |
| Plan `slice-harness-guard-inplace-textwerkzeug` (Nachzug, Abweichungen, Suchlauf-Feld, Träger-Tabelle, DoD-Haken) | geprüft, ohne Befund über F-3, F-5 hinaus: 15 Suchlauf-Zeilen stimmen (Werkzeug und fünf Zeilen von Hand), Nichtgefundenes und Meldungen mit Frist stehen; die DoD-Haken (Liefer-Punkt 1 und 3, Suchlauf, Doku-Update, `make gates`) stimmen mit dem Diff, Liefer-Punkt 2, Review, Closure, Register, Risiko-Ausgänge und Paarungen bleiben `[ ]`; keine Wellen-/Spec-Änderung |
| Hard Rules, Commit-Struktur, Handbuch | geprüft, ohne Befund: Docker-only (keine Host-Toolchain, kein in-place-Aufruf im Diff), keine Suppression, kein Gate gelockert (`make gates` Exit 0, kein neues Gate), sechs Betreffs mit `ADR-0083` ohne `SPEC-`/`ARC-`, Lifecycle-Moves rein, Benutzerhandbuch unberührt (keine Betreiber-Oberfläche), `spec/` unberührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite · Falsch-Positiv-Rand ungebunden und
nicht benannt · Bestandsregel verändert, als unverändert ausgewiesen · Grenz-Zeile nennt eine Lücke der Zusage
nicht · Meldung und Regel nennen verschiedene zulässige Wege · Grenz-Zeile unvollständig oder ungleichlautend ·
Selbstprüfung enger als der Prüfpunkt · Wächter-Regel ohne Evidenz / Interpreter-Denylist · Äquivalente Mutation
durch Doppelschutz

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (HIGH nach der Liste des Reviewer-Skills: die Zeichenklassen der sed- und
perl-Regel und die Pfadzeichen-Klasse färben keinen Fall rot, obwohl `MR-003` sagt, der Tabellentest binde jede
Zusage an ihre Eingabe) und der vier MEDIUM F-2 bis F-5. Der Kern trägt: die 21 vom Plan genannten Mutationen
färben rot, die 51 roten Fälle gegen den Guard am Parent stimmen, Ausgabeform, Exit-Semantik und Latenz sind wie im
Bestand, `make test-command-guard`, `make docs-check`, `make gates` sind grün, die Träger nennen Ursprung und
Grenze. Zu klären sind: F-1 (je ein Fall pro Klassenmitglied), F-3 (Wrapper-Optionen nur für die in-place Klasse
lesen oder die Änderung ausweisen; `command -v`), F-4 und F-6 (Grenz-Zeile in `MR-003`, Kopfkommentar und
`AGENTS.md` §3.1 nachziehen und gleichlautend halten), F-5 (Meldung und Regel auf einen zulässigen Weg bringen).
F-2 berührt die Rückführung §4 (a) des Plans: ob ein Muster mit `|` zulässig geblockt bleibt oder der Guard
Quote-Bewusstsein bekommt, ist eine Architect-Frage; bis zur Antwort steht die Fläche als benannter Rand in
`MR-003` und als Tabellenfall (mittleres und letztes Glied). F-7 bis F-9 sind Hinweise ohne erwartete Aktion in
diesem Slice.

**Übergabe:** Findings gehen an den Implementer (Fixrunde; F-2 zusätzlich an den Architect über den
Konflikt-Pfad des Plans §4); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Die DoD-Zeile „Review durchgeführt“ im Slice-Plan bleibt offen (eine Fixrunde folgt; sie wird bei Schritt
21 des Implementer-Ablaufs nachgezogen). Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser
Skill, dieses Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der
Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes
Prüf-Artefakt, anderer Eingabe-Kontext).
