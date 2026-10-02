# Review-Report: slice-handbuch-public-doc-check-gate-und-skill — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan, [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md),
das Architect-Verdikt `architect-verdict-handbuch-public-doc-check-gate-und-skill`,
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) (Muster) und
[`AGENTS.md`](../../AGENTS.md) Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-handbuch-public-doc-check-gate-und-skill` (wellenlos), Diff-Range
`05fe2400..HEAD`: Implementierung `caab92ac` (Skript, Tabellentest, Sensor-Vertrag, Make-Ziele,
Skill, `harness/README.md`, Regeltexte), Nachzug `56af618f` (Prüfpunkte, Standard),
Plan-Suchlauf `89915e3a`. 10 geänderte Dateien.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt (letzte Änderung im Diff selbst: HIGH-Punkt Handbuch-Versionshistorie).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle Mutationen
liefen an `git archive HEAD`-Kopien im Scratchpad (Ausgabe per `awk … > Kopie`, nie `sed -i`, nie
eine Umleitung auf eine Repo-Datei, kein `rm -rf` mit Variablen); `git status --short` im Echtrepo
blieb leer. `make gates` und alle weiteren Läufe liefen gegen den echten Arbeitsbaum, ungefiltert,
Exit-Code direkt gelesen ([`AGENTS.md`](../../AGENTS.md) §3.9). Die DoD-Checkbox „Review durchgeführt“
im Slice-Plan wird mit diesem Commit nachgezogen (Skill §DoD-Checkbox-Nachzug ohne Fixrunde). Ein
`rm -rf "$S/$1"` in einem ersten Mutations-Anlauf wurde von der Berechtigungsschicht verweigert; der
Aufruf wurde gestrichen und durch `mkdir` neuer Verzeichnisse ersetzt (kein Ersatzweg zum selben
Ziel, kein geteilter Zustand berührt; [`AGENTS.md`](../../AGENTS.md) §3.15).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-handbuch-public-doc-check-gate-und-skill` (§1 Skill-Inhalt, §2 DoD, §3 Suchlauf, §6 Risiken)
- [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) (Festlegungen 1 bis 9),
  [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) (Vorbild `tools/harness/sdk-public-doc-check.sh`),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-QA-OPS-001`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.6, §3.7, §3.9, §3.12, §3.13, §3.15, §4

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

- **Muster wörtlich verglichen:** `pat_id` und `pat_link` im Skript sind zeichengleich zu `P` und `L`
  in `ADR-0143` Festlegung 3 und im Sensor-Vertrag; ERE, `grep -InE -e … -e …`, kein `grep -P`.
- **Gate und Tabellentest:** `make test-handbuch-public-doc-check` Exit 0 („alle 31 Fälle bestanden“),
  `make handbuch-public-doc-check` Exit 0 („keine interne Kennung in 3 Nutzerdokumenten unter docs/user/“).
- **`make gates`**, ungefiltert in eine Log-Datei, Exit gesondert gelesen: Exit 0; die Lauf-Ausgabe
  trägt die Zeile des neuen Gates; `make -n gates` nennt `bash tools/harness/handbuch-public-doc-check.sh`.
- **Weitere Läufe, je Exit 0:** `make test`, `make fmt-check` (323 Go-Dateien), `make kommentar-kennungen DIFF=05fe2400`
  (kein Kandidat), `make suchlauf-nachmessen PLAN=…` (28 Zeilen stimmen), `make docs-check` (0 Befunde),
  `make sdk-public-doc-check`, `make test-sdk-public-doc-check`, `make commit-traceability`.
- **Pfad-Grenze:** `git diff --name-only 05fe2400 HEAD` berührt weder `AGENTS.md` noch `spec/`, `.github/`,
  `docs/user/benutzerhandbuch.md`, noch `ADR-0087`/`-0088`/`-0090`.
- **Mutationen am Skript, einzeln, Tabellentest gegen die Kopie** (Instanz: Skript-Kopie mit
  `git init`-Wurzel, Tabellentest unverändert):

  | Mutation | Ergebnis |
  |---|---|
  | M1: `pat_link` durch ein nie treffendes Muster ersetzt | rot (3 Fälle: Link plan, Link reviews, relativ) |
  | M2: Wortrand `(^\|[^[:alnum:]_-])(slice\|welle)` durch `\b(slice\|welle)` | rot (Fall `byte-slice`) |
  | M3: Klassifikationsbedingung `known -eq 0` durch `false` | rot (beide Fälle unklassifiziert) |
  | M4: Ausnahmeliste um `benutzerhandbuch.md` erweitert (aus `checked` entfernt) | rot (Fenced-Kennung, beide Link-Fälle) |
  | M5: „genannte Datei fehlt“ ohne `class_err=1` | rot (beide Fälle fehlende Datei) |
  | M6: Klassifikations-Exit 2 zu Exit 1 | rot (drei Fälle) |
  | M7: `LH-(FA\|QA\|RB)-` zu `LH-(FA\|QA)-` | **grün** (siehe F-3) |
  | M8: Meldungstexte („unklassifizierte Datei“, „genannte Datei fehlt“, Sammelzeile) ersetzt | **grün** (siehe F-2) |

- **Eigene Gegenbeispiele** (Skript gegen eine Wegwerf-Wurzel, `benutzerhandbuch.md` mit einer Zeile):
  Treffer (Exit 1): `siehe ADR-0143 hier`, `make x LH-FA-ADM-001` (Befehlsbeispiel), `slice-x` am Zeilenanfang
  und eingerückt, `a/slice-x/b`, `"welle-1"`, `docs/plan/` im Fließtext, `../plan/x`,
  `https://github.com/o/r/blob/main/docs/plan/x.md`, `LH-RB-001`, `CO-1234`; kein Treffer (Exit 0):
  `adr-0134`, `ADR-` ohne Ziffer, `Welle-3`, `https://example.com/plan/x`, `cmd --slice-x`,
  `Tag sdk-python-v0.2.0`, `LH-FA` ohne Rest, `ADR 0143`, `ADR` mit Gedankenstrich U+2013 oder U+2011.
  Dateiname mit Leerzeichen und Unterverzeichnis: Exit 2, Name korrekt ausgegeben.
- **Fehlerfälle** (Exit 0 trotz Lesefehler gemessen, siehe F-1): `benutzerhandbuch.md` mit `chmod 000`;
  Unterverzeichnis `docs/user/sub` mit `chmod 000`; Datei mit eingebettetem NUL-Byte und `ADR-0001`.
- **Regeltexte:** `.claude/commands/implement-slice.md` 373 auf 369 Zeilen (`wc -l` an Parent und HEAD).
  Schritt 17 behält „`Version:`-Kopf hochzählen **und** neue Zeile in `### Änderungshistorie`“,
  ergänzt „in Betreibersicht und ohne Kennungen“ und die Verweise auf Skill und Gate;
  der Reviewer-Skill behält die Pflicht im HIGH-Punkt und ergänzt dieselben Verweise. Der Kandidatenlauf
  steht vollständig im Skill (beide Befehle samt Auswertungsregel).

---

## Findings

### F-1 — Gate meldet „keine interne Kennung“ nach einem Lesefehler

- `kategorie`: LOW
- `quelle`: Maintainability; [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Festlegung 1 (Klassifikationsfehler sind laut)
- `pfad`: `tools/harness/handbuch-public-doc-check.sh:62` und `:54`
- `befund`: `h=$(grep … || true)` verschluckt auch grep-Exit 2 (nicht lesbare Datei): gemessen mit `chmod 000` auf
  `benutzerhandbuch.md` druckt grep den Fehler, das Skript aber „keine interne Kennung in 3 Nutzerdokumenten“ und Exit 0;
  `find … 2>/dev/null` verbirgt ein nicht lesbares Unterverzeichnis ebenso (Exit 0), und `grep -I` überspringt eine Datei mit
  NUL-Byte stillschweigend (gemessen: `ADR-0001` darin bleibt unentdeckt). Das Vorbild `sdk-public-doc-check.sh` trägt dieselbe `|| true`-Form.
- `verifizierbar`: ja — `make test-handbuch-public-doc-check` mit einem Fall „nicht lesbare geprüfte Datei“.
- `klasse`: Gate fail-open bei Lesefehler

### F-2 — Tabellentest bindet nur Exit-Codes, nicht den Meldungstext

- `kategorie`: LOW
- `quelle`: `.harness/skills/reviewer.md` HIGH-Klasse „Zusage ohne Bindung an ihre Eingabeseite“ (hier an der Ausgabeseite, deshalb LOW); Sensor-Vertrag `harness/sensors/handbuch-public-doc-check.md` §Ausgabe und Ausgänge
- `pfad`: `tools/harness/run-handbuch-public-doc-check-tests.sh:27-29`
- `befund`: `run()` verwirft stdout und stderr; der Vertrag nennt drei Meldungstexte (Sammelzeile, `unklassifizierte Datei`, `genannte Datei fehlt`)
  und das Zeilenformat `docs/user/<datei>:<zeile>:<text>`. Mutation M8 (alle Meldungstexte ersetzt) lässt alle 31 Fälle grün; ebenso wäre ein Vertauschen der beiden Exit-2-Ursachen nicht sichtbar.
- `verifizierbar`: ja — Tabellentest um Ausgabe-Assertions ergänzen.
- `klasse`: Test bindet Ausgabeseite nicht

### F-3 — `LH-RB-`-Zweig des Musters P ohne Tabellenfall

- `kategorie`: LOW
- `quelle`: Maintainability; [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Festlegung 3 (Tabelle der Muster-Befunde)
- `pfad`: `tools/harness/run-handbuch-public-doc-check-tests.sh:51-52`
- `befund`: Die Fälle decken `LH-FA-` und `LH-QA-`, nicht `LH-RB-`; Mutation M7 (`RB` aus dem Muster entfernt) lässt den Tabellentest grün. Im Repo kommt `LH-RB-` weder in `spec/` noch in den ADRs vor (`git grep -c "LH-RB-" -- spec docs/plan/adr` leer); der Zweig stammt aus dem ID-Schema der Konventionen.
- `verifizierbar`: ja — ein Tabellenfall `LH-RB-001`.
- `klasse`: Test bindet Musterzweig nicht

### F-4 — Plan widerspricht sich zum Kandidatenlauf in Schritt 17

- `kategorie`: LOW
- `quelle`: Maintainability (Nachzug widerspricht dem Nachbarn im selben Träger)
- `pfad`: `docs/plan/planning/done/slice-handbuch-public-doc-check-gate-und-skill.md:212` gegen `:273`; `.claude/commands/implement-slice.md:124-134`
- `befund`: DoD-Punkt 3 sagt „Der Kandidatenlauf in Schritt 17 bleibt“, die Träger-Tabelle (Zeile 273) „Inhalt zum Skill verlagern“. Umgesetzt ist die Verlagerung: Schritt 17 nennt Pflicht, Betreibersicht, Verweis und Grenze, die beiden Befehle des Kandidatenlaufs stehen nur im Skill; Schritt 17 allein ist für den Lauf nicht ausreichend, der Skill ist vollständig.
- `verifizierbar`: nein — Lese-Handlung; der Verifier liest den DoD-Wortlaut gegen die Umsetzung.
- `klasse`: Plan-interner Widerspruch DoD gegen Träger-Tabelle

### F-5 — Sensor-Vertrag nennt weniger Host-Werkzeuge, als das Skript ruft

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.1 (Host-Klasse; ein Vertrag nennt die Werkzeuge seines Skripts, soweit sie über `bash` und `git` hinausgehen)
- `pfad`: `harness/sensors/handbuch-public-doc-check.md` §Sperren; Skriptkopf Zeile 17
- `befund`: Vertrag und Kopf sagen „nur bash/git/grep“; das Skript ruft zusätzlich `find`, `sort` und `sed` (alle in der POSIX-/coreutils-Klasse, keine Regelverletzung), der Tabellentest `mktemp` und `rm`.
- `verifizierbar`: nein
- `klasse`: Werkzeugliste im Vertrag unvollständig

### F-6 — Muster-Grenzen jenseits der benannten (gemessen)

- `kategorie`: INFO
- `quelle`: Sensor-Vertrag §Grenze Punkt 2
- `pfad`: `harness/sensors/handbuch-public-doc-check.md` §Grenze
- `befund`: Neben Groß-/Kleinschreibung bleiben auch `ADR 0143` (Leerzeichen), `ADR` mit Gedankenstrich und `--slice-x` grün (gemessen); umgekehrt trifft Muster L jeden Link mit `docs/plan/` im Pfad, auch eine legitime GitHub-URL auf das Repo (gemessen, Exit 1). Beides ist mit der Entscheidung „keine Ausnahme“ vereinbar, steht aber nicht im Vertrag.
- `verifizierbar`: ja
- `klasse`: Gate-Grenze unvollständig benannt

### F-7 — Tabellentest ruft `rm -rf "$tmp/root"` ohne Absicherung gegen leeres `tmp`

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-handbuch-public-doc-check-tests.sh:10,16`
- `befund`: Das Skript läuft mit `set -uo pipefail` ohne `-e`; schlägt `mktemp -d` fehl, ist `tmp` leer und `fresh` entfernt `/root`. Dieselbe Form trägt `run-sdk-public-doc-check-tests.sh`.
- `verifizierbar`: nein
- `klasse`: Aufräumpfad ohne leeres-Variable-Guard

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/handbuch-public-doc-check.sh` Muster, Reichweite (Hybrid), Exit-Semantik 0/1/2 gegen Vertrag; Dateinamen mit Leerzeichen und Unterverzeichnisse; `grep` liest den Arbeitsbaum (Vertrag benennt es), kein Pipe-`grep -q` mit weiterschreibendem Erzeuger.
- geprüft, ohne Befund: `harness/mk/doc-gate.mk` (`GATE_CHECKS +=`, `.PHONY`, Werkzeug-Ziel ohne Gate-Aufnahme) und Einbindung; `make -n gates` und die Lauf-Ausgabe von `make gates` nennen das Gate.
- geprüft, ohne Befund: `harness/README.md` (Gate-Zeile, Werkzeug-Zeile des Tabellentests, `make gates`-Aufzählung, Skill in §Guides) gegen das Makefile (§4).
- geprüft, ohne Befund: `.harness/skills/nutzerdoku-schreiben.md` (vollständig laut Plan §1, im Text keine Kennung als Beispiel, nennt `P`/`L` nicht, Verweise aus `implement-slice.md` und `reviewer.md`).
- geprüft, ohne Befund: `.claude/commands/implement-slice.md` und `.harness/skills/reviewer.md` (Pflicht „Version hochzählen und Zeile in Änderungshistorie“ unverändert, „Neue Betreiber-Oberfläche zieht das Handbuch mit“ unverändert, Zeilenzahl nicht größer).
- geprüft, ohne Befund: `docs/user/benutzerhandbuch-standard.md` (Kennungsfreiheit der Historie an Zeile 229 und Kapitel 11, ohne Beispiel-Kennung, Gate grün).
- geprüft, ohne Befund: Kommentare in Skript, Tabellentest und Makefile-Fragment (§3.7: Zusage, Abgrenzung, Rang-Zeiger; höchstens eine Kennung je Block); Suppression (§3.2) nicht vorhanden; Gate-Lockerung (§3.6) nicht vorhanden; Herkunft der Messwerte im Plan (§3.12) durch `make suchlauf-nachmessen` getragen (28 Zeilen).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Gate fail-open bei Lesefehler · Test bindet Ausgabeseite nicht · Test bindet Musterzweig nicht · Plan-interner Widerspruch DoD gegen Träger-Tabelle · Werkzeugliste im Vertrag unvollständig · Gate-Grenze unvollständig benannt · Aufräumpfad ohne leeres-Variable-Guard

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM; die vier LOW gehen ohne Rückgabe-Pfeil an den Planner bzw. Implementer eines Folgezugs.

**Fragen an den Architect:**

1. F-1: Soll `ADR-0143` für das Gate Lesefehler als Exit 2 festlegen (dann auch für `sdk-public-doc-check`, [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)), oder bleibt Vorbild-Parität (`|| true`) die Entscheidung?
2. F-4: Gilt DoD-Punkt 3 „Kandidatenlauf in Schritt 17 bleibt“ oder die Träger-Tabelle „Inhalt zum Skill verlagern“?

**Übergabe:** Die DoD-Zeile „Review durchgeführt“ im Slice-Plan ist mit diesem Commit auf `[x]` gezogen (kein offenes HIGH/MEDIUM, keine Fixrunde nötig). Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation ([`AGENTS.md`](../../AGENTS.md) §6).
