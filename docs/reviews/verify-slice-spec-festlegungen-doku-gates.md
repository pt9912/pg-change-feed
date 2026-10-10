# Verifikationsbericht: slice-spec-festlegungen-doku-gates — 2026-10-10

**Rolle:** Verifier (Modul 11). Die Frage lautet: „Bauen wir es richtig?“. Geprüft
wird gegen die DoD von `slice-spec-festlegungen-doku-gates` (§2, Liefer-Punkte 1
und 2, die Bedingung `Schärft:`-Kante und die Gate-Pflicht), gegen die neue
Festlegung [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
gegen [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
(Entscheidungen 1 bis 4, Folgepflicht 1 und 2), gegen das Architect-Verdikt
[`architect-verdict-matrix-inactive-nur-link`](architect-verdict-matrix-inactive-nur-link.md)
§3 und gegen [`AGENTS.md`](../../AGENTS.md) §3.5, §3.9, §3.12 und §3.13. Den Diff als
Maintainability-Frage prüft der Reviewer:
[`review-slice-spec-festlegungen-doku-gates`](review-slice-spec-festlegungen-doku-gates.md)
und das Re-Review
[`review-slice-spec-festlegungen-doku-gates-re-review`](review-slice-spec-festlegungen-doku-gates-re-review.md)
(1 HIGH N-1, 1 LOW N-2, 2 INFO). Den realen Bedarf prüft der Validator. Dies ist
kein MVP-Slice.

**Gegenstand:** Diff `942ebf3f..3122ff63`, 11 Commits, lokal und nicht gepusht.
Darunter sind der Liefer-Commit `9e43f968`, der Architect-Commit `deecc466`
(`ADR-0163`, ADR-Index), die Folgepflicht `64c68933`, die Fixrunde `7f139836`, das
Verdikt `037a9d32`, dessen Nachzug `2ac39c31` und der N-1-Fix `3122ff63`. Der Slice
liegt in `in-progress/`. Die Closure-Punkte der DoD stehen noch aus.

**Frischer Kontext:** In dieser Sitzung habe ich den Plan vollständig gelesen,
außerdem `ADR-0163`, das Verdikt, beide Review-Reports, den Spec-Diff (§7 `SPEC-040`
und §8), den Vertrag `harness/sensors/docs-check.md` im Stand `3122ff63`, die Diffs
von `harness/README.md` und ADR-Index sowie die Nicht-Kommentar-Zeilen von
`.d-check.yml`. Keine Behauptung wurde übernommen, außer wo „übernommen“ steht.
Gemessen wurde am Stand `3122ff63` bei sauberem Arbeitsbaum. Jeden Exit-Code habe ich
direkt gesichert (`echo $?` unmittelbar nach dem Aufruf, bei langer Ausgabe über eine
Log-Datei im Scratchpad; `AGENTS.md` §3.9). Proben liefen nur an einem Klon im
Scratchpad (`<Scratchpad>/vk`, Stand `3122ff63`). Je Mutation wurde danach
`git reset --hard && git clean -fd` ausgeführt, und `git status --porcelain` war leer.
Außer diesem Bericht wurde keine Repo-Datei geschrieben.

---

## 1. DoD je Punkt

| DoD-Punkt | Verdikt | Beleg (selbst gefahren) |
|---|---|---|
| **LP 1**: Festlegung in §7 | **nicht vollständig bestätigt**: Abweichung A-1 (Absatz „Ausgänge“) | §7-Tabelle trägt die Zeile `SPEC-040` mit acht Modulen. Die Absätze Fläche, Befund je Modul, Randformen, `structure`, `hostpaths` und Ausgänge sind vorhanden, ebenso die Historie-Zeile vom 2026-10-10. Die nächste freie Kennung ist belegt: `git grep -ohE 'SPEC-[0-9]{3}' 942ebf3f \| sort -u \| tail -1` ergibt `SPEC-039`. Der Liefer-Stand ist belegt: d-check am Klon von `9e43f968` ergibt Exit 0 und `d-check: 1830 Datei(en) geprüft, 0 Befund(e)`, wie im Plan. Die Gegenprobe hat 76 Zeilen (A1–A36, B1–B16, C1–C24, ausgezählt). 44 Proben am Werkzeug (Abschnitt 3) folgen der Festlegung in Befund und Exit; vier Zusatzmessungen zur Ausgabe-Form (A1a–A1d) weichen vom Absatz „Ausgänge“ ab (A-1) |
| **LP 2**: Vertrag und Index verweisen | **bestätigt**, mit Lese-Hinweis A-2 | §Vertrag „Wer was trägt“, Grenze 8, 9, 11 und 12, §Ausgabe, §Sperren und §Bindung verweisen auf `SPEC-040`. `grep -nE '220\|120\|Exit [01]'` auf den Vertrag findet keine Schwelle und keine Exit-Tabelle. Die Befund-Code-Tabelle ist entfallen. `harness/README.md` Zeile `make docs-check`, Spalte Bindung, enthält: Vertragsdatei · `SPEC-040` |
| `Schärft:`-Kante | **bestätigt** | `ADR-0163` hat Status `Accepted`, `Schärft:` zeigt auf `SPEC-040`, und die Index-Zeile steht im Diff. Die Quell-ADRs sind unverändert: `git diff --stat 942ebf3f..3122ff63 -- docs/plan/adr/` zeigt nur `0163-…` und `README.md`. Folgepflicht 1 ist wörtlich in den Randformen umgesetzt (Abschnitt 4) |
| `make gates` grün | **bestätigt** | eigener Lauf am Stand `3122ff63`: Exit 0 (Abschnitt 6); ein zweiter Lauf nach dem Commit dieses Berichts steht in der Rückmeldung |
| Review durchgeführt | **materiell bestätigt**, Häkchen fehlt | Beide Reports liegen vor. N-1 ist in `3122ff63` behoben: Die Fixrunde-Zeile zu F-9 nennt jetzt `git grep -n status-provenance -- spec/ \| grep -cE 'ADR-[0-9]{4}\|slice-[0-9]{3}'`. Selbst gefahren druckt der Befehl `0` (2 Marker-Zeilen, 2062 und 2208, keine mit Token), und der Befehl trägt den Satz. N-2 (LOW) steht offen für die Closure. Das DoD-Häkchen setzt der Plan noch nicht |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen** (nicht Gegenstand dieses Laufs) | §7 des Plans trägt noch die Platzhalter |

**Bewusstes Brechen der DoD-Testbehauptung** (Modul 11). Die Aussage „die Spec zeigt
auf keine ADR, belegt durch `make docs-check` Exit 0“ (Plan §6, vierter Punkt) habe
ich am Klon gebrochen. Ein Link aus `spec/pflichtenheft.md` auf `ADR-0163` ergibt
`matrix-forbidden`, Exit 1 (V4d). Mit dem Marker `d-check:status-provenance` bleibt
derselbe Link rot (V4f). Ein Token mit Marker bleibt dagegen grün (V4e). Das ist die
akzeptierte Lücke aus `ADR-0163` Entscheidung 2, und die zusätzliche
`git grep`-Messung aus §6 deckt sie zu Recht ab. Das Gate wird also aus dem richtigen
Grund rot, und die Lücke ist benannt.

## 2. Plan gegen Code-Diff

| Plan §3 | Diff | Ergebnis |
|---|---|---|
| `spec/pflichtenheft.md` §7 und §8 Historie | `9e43f968`, `64c68933`, `7f139836`, `2ac39c31` | wie geplant |
| `harness/sensors/docs-check.md` (Verweis, §Sperren, §Bindung) | `9e43f968`, `7f139836`, `2ac39c31` | wie geplant; Grenze 12 ist neu (Review F-7, Verdikt §3.2) |
| `harness/README.md` §Sensors | `9e43f968` | wie geplant, eine Zelle |
| ADR und ADR-Index | `deecc466` | wie geplant; die Index-Zeilen `ADR-0094` und `ADR-0097` tragen den Teil-Supersede-Hinweis (Hinweis A-4) |
| `AGENTS.md` §3.11, `.d-check.yml` nicht geändert | kein Diff an beiden | wie geplant |
| Reports und Verdikt unter `docs/reviews/` | `72e57998`, `037a9d32`, `b36af1f1` | Artefakte anderer Rollen, keine Plan-Zeile nötig |

Kein Diff-Teil ist ohne Plan-Zeile, und kein Produkt-Code ist berührt. Die neun
Dateien aus `git diff --stat` liegen unter `spec/`, `harness/` und `docs/`.

## 3. `SPEC-040` gegen das Werkzeug (Klon im Scratchpad)

Gefahren wurde d-check v0.82.0 über den Digest aus `d-check.mk` (`sha256:d28e9437…`),
`docker run --rm --network none -v <Klon>:/repo:ro <Digest>`. Basislauf am Klon
`3122ff63`: Exit 0, `d-check: 1834 Datei(en) geprüft, 0 Befund(e)`. Jede Mutation ist
eine Änderung an einer Datei und wurde danach zurückgesetzt. Gedruckt ist der
Grund-Code der Befundzeile. Host-lokale Pfade stehen hier nur als Platzhalter
(`AGENTS.md` §3.11).

| # | Modul | Mutation | Erwartet nach `SPEC-040` | Ergebnis |
|---|---|---|---|---|
| V1 | `links` | Link auf eine fehlende Datei in `README.md` | `target-missing` | `target-missing`, Exit 1 |
| V2a | `anchors` | Link auf ein fehlendes Heading | `anchor-missing` | `anchor-missing`, Exit 1 |
| V2b | `anchors` | Link auf die HTML-`id` `mr-002` (Tabellenzeile in `harness/conventions.md`) | kein Treffer | Exit 0 |
| V2c | `links` | Link mit URL-Schema auf ein Ziel, das nicht existiert | kein Gegenstand | Exit 0 |
| V3a | `ids` | `SPEC-040` nackt in Prosa | `id-unlinked` | `id-unlinked`, Exit 1 |
| V3b | `ids` | dasselbe mit Zeilen-Marker `d-check:ignore` | kein Treffer | Exit 0 |
| V3c | `ids` | `MR-002` und `BEO-PGC/xyz` nackt | Kennung ohne Muster, kein Treffer | Exit 0 |
| V3d | `ids` | erfundene Kennung `SPEC-999`, verlinkt | kein Treffer, Linkziel ungeprüft | Exit 0 |
| V4a | `matrix` | Token `ADR-0094` in Inline-Code in `spec/architecture.md` | `matrix-forbidden` | `matrix-forbidden`, Exit 1 |
| V4b | `matrix` | `ADR-0094 slice-099` in einem Fence derselben Datei | kein Treffer | Exit 0 |
| V4c | `matrix` | V4a mit `d-check:ignore` | wirkt auf `matrix` nicht | `matrix-forbidden`, Exit 1 |
| V4d | `matrix` | Link aus `spec/pflichtenheft.md` auf `ADR-0163` | `matrix-forbidden` | `matrix-forbidden`, Exit 1 |
| V4e | `matrix` | Token `ADR-0163` mit `d-check:status-provenance` im Pflichtenheft | kein Treffer (Lücke) | Exit 0 |
| V4f | `matrix` | V4d mit `d-check:status-provenance` | Link bleibt Befund | `matrix-forbidden`, Exit 1 |
| V5a | `versions` | Pin auf ein älteres Baseline-Tag in Inline-Code in `README.md` | `version-stale` | `version-stale`, Exit 1 |
| V5b | `versions` | dasselbe in `ADR-0163` | ausgenommen (`docs/plan/adr/[0-9]*.md`) | Exit 0 |
| V5c | `versions` | dasselbe im ADR-Index | Index bleibt geprüft | `version-stale`, Exit 1 |
| V5d | `versions` | dasselbe in `MR-002` | ausgenommen | Exit 0 |
| V5e | `links` | Link in ein älteres Tag aus `ADR-0163` | `links` prüft weiter | `target-missing`, Exit 1 |
| V6a | Fläche | kaputter Link in `.tmp/probe.md` | `scan.ignore` | Exit 0 |
| V6b | Fläche | kaputter Link in einer indizierten `docs/probe.template.md` | `scan.ignore` | Exit 0 |
| V7a | `structure` | Gate-Index, Zelle `Vertrag` mit 220 Zeichen | kein Treffer | Exit 0 |
| V7b | `structure` | dieselbe Zelle mit 221 Zeichen | `section-cell-oversized` | `section-cell-oversized` („221 Zeichen, erlaubt sind 220“), Exit 1 |
| V7c | `structure` | 219 Zeichen in Backticks (221 im Quelltext) | gezählt mit Markdown-Syntax | `section-cell-oversized` („221 Zeichen“), Exit 1 |
| V7d | `structure` | dieselbe Zelle leer | `section-cell-undersized` | `section-cell-undersized`, Exit 1 |
| V7e | `structure` | Kopfzelle `Vertrag` der ersten Tabelle umbenannt | `section-column-missing` | `section-column-missing`, Exit 1 |
| V7f | `structure` | `forbid-pattern` der §7-Regel als Prosa in `done/slice-001-bootstrap.md` §7 | `section-forbidden` | `section-forbidden`, Exit 1 |
| V7g | `structure` | dasselbe in Inline-Code | bereinigter Text, kein Treffer | Exit 0 |
| V7h | `structure` | Review-Bericht mit Markdown-Link auf `open/slice-…` | Verweisform-Regel | `section-forbidden` und `target-missing`, Exit 1 |
| V7i | `structure` | ADR-Index, Status-Zelle 18 Zeichen | `section-cell-oversized` | `section-cell-oversized` („erlaubt sind 10“), Exit 1 |
| V8a–V8e | `hostpaths` | `<Host-Wurzel>/…` in Prosa, dasselbe in Inline-Code, Windows-Laufwerksmuster, `~/<Verzeichnis>` ohne Schrägstrich, `~/<Datei>` | `hostpath-forbidden` | je `hostpath-forbidden`, Exit 1 |
| V8f–V8h | `hostpaths` | `<Host-Wurzel>/…` im Fence; relativer Pfad; `<Host-Wurzel>/…` mit `d-check:ignore` | kein Treffer, kein Treffer, Treffer (kein Zeilen-Marker) | Exit 0, Exit 0, `hostpath-forbidden` Exit 1 |
| V9a/V9b | `tracked` | Link auf eine existierende, ungetrackte Datei bzw. auf das ignorierte `harness/image-hash.txt` | `target-untracked` | je `target-untracked`, Exit 1 |
| V10a | Ausgänge | unbekannter Schlüssel in `.d-check.yml` | Exit 2, kein Rückfall | Exit 2 |
| V10b | Ausgänge | `versions.current-from` auf eine fehlende Datei | Exit 2 | Exit 2, `d-check: error: versions.current-from nicht lesbar …` |
| V10c | Ausgänge | unbekanntes Modul in `modules:` | Exit 2 (Konfigurationsfehler) | Exit 2, `d-check: error: unbekanntes Modul …` |
| V11 | Ausgänge | Baum ohne `.git` (`git archive`-Export von `3122ff63`) | Exit 2 | Exit 2, `d-check: error: kein lesbares git-Repository unter /repo …` |

**Bei Befund und Exit stimmen alle 44 Proben über alle acht Module mit `SPEC-040`
überein.** Davon abgesetzt sind vier Zusatzmessungen zum Absatz „Ausgänge“, die
nicht übereinstimmen (A-1).

| # | Lauf | Gedruckt |
|---|---|---|
| A1a | V1 direkt, stdout und stderr getrennt | stdout: `README.md:65⇥nix-da.md⇥target-missing⇥Linkziel existiert nicht`; stderr: `d-check: 1834 Datei(en) geprüft, 1 Befund(e)` |
| A1b | V1 über `make docs-check`, `2>&1` in eine Datei | Reihenfolge: Schlusszeile, dann die Befundzeile, dann `make: *** […] Fehler 1`; Make-Exit 2 |
| A1c | V10a direkt, stdout und stderr getrennt | stdout leer; stderr zwei Zeilen: `d-check: error: .d-check.yml: yaml: unmarshal errors:` und `  line 349: field unbekannt_schluessel not found` |
| A1d | V8 (`hostpaths`), stdout | `README.md:66⇥<Pfad>⇥hostpath-forbidden`: drei Felder, kein Text |

## 4. ADR- und Verdikt-Konformität

| Prüfpunkt | Ergebnis |
|---|---|
| Matrix-Regel: Spec verweist auf keine ADR | erfüllt. `git grep -nE 'ADR-[0-9]{4}\|slice-[0-9]{3}\|docs/plan/adr\|docs/reviews' -- spec/` trifft eine Zeile (2093, „Berichten unter `docs/reviews/`“, Verzeichnisname, kein Token der Klasse `review`). Die zwei Marker-Zeilen tragen kein Token. `make docs-check` ergibt 0 Befunde. Gebrochen in V4d und V4f |
| `ADR-0163` Entscheidung 1 (Herkunft je Teil) | Jeder Teil der Tabelle hat in `SPEC-040` einen Absatz: Fläche, Befund je Modul, `matrix` samt Status, Marker und `exempt-paths`, `versions`, `hostpaths`, `d-check:ignore`. Die Ausgänge stammen nach der Tabelle aus „diese ADR“, gegeben durch das Werkzeug. Dort liegt A-1 |
| Entscheidung 2 (Marker) | Der Satz in den Randformen ist wörtlich der aus Folgepflicht 1; V4e und V4f bestätigen ihn |
| Entscheidung 3 (`matrix.exempt-paths`) | Der Satz ist wörtlich aus Folgepflicht 1. Die Wirkung ist am Werkzeug bestätigt: V5e (Ziel bleibt geprüft für `links`), außerdem M16 und Zeile 8 der ADR (*übernommen*) |
| Entscheidung 4 (Ort der Modul-Semantik) | Vertrag §Vertrag „Wer was trägt“ nennt `SPEC-040`, und Grenze 8 verweist für die Formen dorthin |
| Verdikt §3.1 und §3.2 | Zelle `matrix-inactive`, Randformen-Satz, Überschrift und Satz von Grenze 12 sowie die Verweis-Ergänzung habe ich gelesen und gegen das Verdikt gehalten: wörtlich |
| Entscheidung 1, letzter Absatz („weicht das Werkzeug ab, ist das Werkzeug falsch“) | Das setzt eine richtige Festlegung voraus (Verdikt §1 Punkt 3). A-1 ist ein Fall, in dem die Festlegung falsch behauptet, nicht das Werkzeug |

## 5. Weitere Läufe

| Lauf | Exit | Gedruckt |
|---|---|---|
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-spec-festlegungen-doku-gates.md` | 0 | `suchlauf-nachmessen: 16 Zeilen stimmen` (u. a. `diff 35 -F 'SPEC-040'` mit Soll 35 und Ist 35) |
| `make doc-commits RANGE=942ebf3f..3122ff63` | 0 | `d-check: 1834 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=942ebf3f..3122ff63` | 0 | `d-check: 1834 Datei(en) geprüft, 0 Befund(e)`. Kein Commit der Range berührt `harness/conventions/` (`git diff-tree` je Commit, 11-mal 0). Es gibt also weder Form- noch Umzugs-Commit, und keine Teilung ist nötig |

## 6. Gate-Lauf

`make gates > <Scratchpad>/gates1.log 2>&1; echo $? > …` am Stand `3122ff63` bei
sauberem Baum ergab **Exit 0**. Gedruckt wurden u. a.
`coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`, zweimal
`d-check: 1834 Datei(en) geprüft, 0 Befund(e)`,
`commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`, `generated-sync: OK` und
`gesamt: 0 Befund(e)`. Den Lauf nach dem Commit dieses Berichts kann der Bericht nicht
tragen. Er steht in der Rückmeldung an den Planner.

## 7. Abweichungen

| ID | Klasse | Befund | Bewertung |
|---|---|---|---|
| A-1 | DoD (LP 1), Festlegung behauptet mehr als das Werkzeug | `SPEC-040` §Ausgänge sagt: „Je Befund steht eine Zeile `<Datei>:<Zeile>`, Ziel, Grund-Code und Text, getrennt durch Tabulatoren; die letzte Zeile lautet `d-check: <N> Datei(en) geprüft, <M> Befund(e)`“ und „ob ein Befund oder ein Fehler vorliegt, sagt die letzte Zeile (`<M> Befund(e)` oder `d-check: error: …`)“. Am Werkzeug gemessen gilt dreierlei. **(a)** Die Befundzeilen gehen nach stdout, die Schlusszeile nach stderr (A1a). Über `make` mit `2>&1` steht die Schlusszeile vor der Befundzeile, und die letzte Zeile ist die von `make` (A1b). **(b)** Ein Konfigurationsfehler aus dem YAML-Parser endet mit einer Fortsetzungszeile `  line 349: …`, nicht mit `d-check: error: …` (A1c). **(c)** Eine `hostpath-forbidden`-Zeile hat drei Felder, ohne Text (A1d) | Die Exit-Codes stimmen, aber die Ausgabe-Form ist falsch festgelegt. Nach `ADR-0163` Entscheidung 1 hieße diese Abweichung „das Werkzeug ist falsch“, und das trifft nicht zu. Eine Berichtigung im Sinne von Verdikt 1 aus Modul 8 („die Festlegung hat falsch behauptet“) wäre etwa: „stderr endet mit …“, „eine Fehlermeldung beginnt mit `d-check: error:` und kann mehrzeilig sein“, „Text kann fehlen“. Das ist eine Fixrunde vor der Closure. Beide Reviews und die Messungen M0/M1a/M1b des Plans haben die Streams nicht getrennt. Die Anschluss-Frage im Plan („es gilt die Schlusszeile“) und der Satz im Vertrag §Ausgabe („ein Beleg zitiert die Schlusszeile“) bleiben richtig |
| A-2 | Beleg breiter als seine Messung (`AGENTS.md` §3.12 Instanz B) | DoD LP 2 sagt: „die Grund-Codes, die Zellgrenzen 220/120 und die Exit-Tabelle stehen nicht mehr in ihm“. Gemessen ist das nur für `höchstens (220\|120)` und „Deklaration dieses Vertrags“. Der Vertrag nennt weiterhin `id-unlinked` (Grenze 7, Zeile 69), `matrix-inactive` (Grenze 12) und `id-unlinked`/`target-missing` (§Ausgabe, `make doc-repair`) | Kein DoD-Bruch: Die Code-Tabelle ist entfallen, und die Restnennungen sind Kontext. Der Satz des Belegs sollte „die Befund-Code-Tabelle“ heißen. Lese-Hinweis für die Closure |
| A-3 | durchgereicht | Re-Review N-2 (Zählwort „zwei angenommene Lücken“ in Grenze 12, Wortlaut des Verdikts) und F-9 (fünf Pläne in `open/`, Frist Closure) sind offen | Ausgang bei der Closure, Planner bzw. Architect |
| A-4 | INFO | Die Index-Zeile `ADR-0097` hat den Titelzusatz `(ergänzt ADR-0094)` gegen `(→ 0099/0163 teilw.)` getauscht und nennt die Ergänzung nicht mehr | kein DoD-Bezug; die ADR-Datei selbst ist unverändert |
| A-5 | INFO | `SPEC-040` §Ausgänge nennt drei Beispiele für Exit 2. Ein unbekanntes Modul in `modules:` ist ein vierter Fall (V10c) und fällt unter „Konfigurationsfehler“ | kein Befund, wenn die Liste als Beispielliste gelesen wird |

## Verdikt

**Nicht bestanden, wegen A-1.** LP 2, die `Schärft:`-Kante, die Gate-Pflicht und das
Review (N-1 behoben, selbst nachgefahren) sind bestätigt. LP 1 trifft bei Befund und
Exit an 44 Proben aus allen acht Modulen, auch `structure`, das Werkzeug. Der Absatz
„Ausgänge“ von `SPEC-040` legt aber eine Ausgabe-Form fest, die das gepinnte
Werkzeug nicht liefert (Schlusszeile auf stderr, mehrzeilige Fehlermeldung, Text
fehlt bei `hostpaths`). Die Berichtigung ist eine Fixrunde des Implementers am
Absatz und an der Historie-Zeile. Danach genügt eine Nachprüfung von A-1. A-2 bis
A-5 gehen an die Closure.

**Übergabe:** A-1 an den Implementer über den Planner; A-2 bis A-5 an die Closure.
Dieser Bericht ist ein Lauf-Beleg. Er ersetzt keine Validierung.
