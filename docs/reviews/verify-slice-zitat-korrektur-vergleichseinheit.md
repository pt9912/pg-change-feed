# Verifikationsbericht: slice-zitat-korrektur-vergleichseinheit — 2026-10-06

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD (`slice-zitat-korrektur-vergleichseinheit` §2, Liefer-Punkte
1 und 2 und die Gate-Pflicht) und gegen die Entscheidungen
[`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md),
[`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
und
[`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md),
außerdem gegen [`AGENTS.md`](../../AGENTS.md) §3.5 und §3.12. Nicht geprüft
wird der Diff als Maintainability-Frage. Das ist Aufgabe des Reviewers:
[`review-slice-zitat-korrektur-vergleichseinheit`](review-slice-zitat-korrektur-vergleichseinheit.md)
und das Re-Review
[`review-slice-zitat-korrektur-vergleichseinheit-fixrunde`](review-slice-zitat-korrektur-vergleichseinheit-fixrunde.md).
Ebenfalls nicht geprüft wird der reale Bedarf. Das ist Aufgabe des Validators,
und dies ist kein MVP-Slice.

**Gegenstand:** Diff `b7d97cca..b2d5f060` mit 9 Commits. Darin sind die
Architect-Commits `07923146` (`ADR-0158`) und `b36f882d` (`ADR-0159`) enthalten.
Der Slice liegt in `in-progress/`; die Closure-Punkte stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat den Plan, beide Review-Reports,
`ADR-0157` Entscheidung 4, `ADR-0158`, `ADR-0159` und den Diff gelesen. Sie hat
keine Behauptung übernommen. Jede Zahl unten ist in diesem Lauf am Stand
`b2d5f060` gemessen (Arbeitsbaum sauber). Die Exit-Codes sind direkt und ohne
Pipe gesichert (`AGENTS.md` §3.9). Die Befehlsform lief an einem `git clone` im
Scratchpad (`<Scratchpad>/verify-vergleich/repo`). Probe-Commits existieren nur
dort, und ihre Dateien wurden per `sed … > <Scratchpad-Datei>` und `cp` in den
Klon geändert. Außer diesem Bericht wurde keine Repo-Datei geschrieben.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| Befehlsform aus `ADR-0159` Entscheidung 4, per `awk` wörtlich aus dem `bash`-Block gezogen (Einrückung von drei Zeichen entfernt), `bash -n` | 56 Zeilen, Syntax gültig; GNU Awk 5.2.1 | 0 |
| DoD-Proben (Abschnitt 2), je dreimal: ohne Option, mit `shopt -s nullglob`, mit `shopt -s failglob` | alle Fälle mit Soll-Exit; die drei Läufe ergeben dieselben Zeilen | je Soll |
| Mutationen der Eingabeseite und der Form (Abschnitt 2) | Farbe je Fall wie unten | je unten |
| Referent-Messung am Pin-Commit `5d8855d9` (Abschnitt 3) | drei Fälle `cmp 0`, `MR-001` nicht messbar | 0/0/0/2 |
| Datei-`cmp` am Pin-Commit `5d8855d9`, Schleife aus `ADR-0157` Entscheidung 4 | `MR-001` bis `MR-004` je Exit 0 | 0 je Datei |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-zitat-korrektur-vergleichseinheit.md` | `suchlauf-nachmessen: 12 Zeilen stimmen` (je Zeile `OK soll=ist`) | 0 |
| `make docs-check` (vor diesem Bericht) | `d-check: 1769 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=b7d97cca..HEAD` | `d-check: 1769 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-immutable RANGE=b7d97cca..HEAD` | `d-check: 1769 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make commit-traceability RANGE=b7d97cca..HEAD` | `commit-traceability: OK — 9 Commit(s) in "b7d97cca..HEAD", Betreffs ohne Struktur-ID` | 0 |
| `make gates` (nach dem Commit dieses Berichts) | siehe Abschnitt 7 | siehe dort |

Die Range enthält keinen Pin-Commit. `git diff --name-only b7d97cca..HEAD -- harness/conventions/`
ist leer, und `5d8855d9` ist Vorfahre von `b7d97cca`
(`git merge-base --is-ancestor`, Exit 0). Deshalb läuft `doc-immutable` über die
ganze Range. Teil-Ranges sind nicht nötig.

## 2. DoD gegen Belege

**Liefer-Punkt 1 (Entscheidung): bestätigt.** `ADR-0158` und `ADR-0159` legen
für jede Verweisform die Einheit und die Normalisierung fest. Nachgefahren mit
der Befehlsform aus `ADR-0159` Entscheidung 4. Die Datei
`F=docs/plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md` ist dieselbe wie
im Plan, die Probe-Commits sind eigene:

| Probe | Befehl (gekürzt) | gedruckte Zeile (Kern) | Exit | Soll |
|---|---|---|---|---|
| Anker-Wechsel auf anderen Abschnitt | `vergleich b2d5f060 $F '#teilfrage-2--subjekt--und-nachrichtenschema' b2d5f060 $F '#teilfrage-3--zustellsemantik-core-nats-vs-jetstream'` | `vergleich roh: …#teilfrage-2-… <-> …#teilfrage-3-… cmp 1` | 1 | fällt |
| Gegenprobe gleicher Anker | dieselbe Adresse an beiden Seiten | `cmp 0` | 0 | besteht |
| Lokator verschoben; Probe-Commit `4db2bc85` fügt vor Zeile 41 von `F` drei Zeilen ein | `vergleich b2d5f060 $F L139-145 4db2bc85 $F L142-148` | `vergleich roh: …mdL139-145 <-> 4db2bc85:…mdL142-148 cmp 0` | 0 | besteht |
| Gegenprobe Lokator nicht nachgezogen | `L139-145` an beiden Seiten | `cmp 1` | 1 | fällt |
| `#guard-haertung`, ein Wort im Körper; Probe-Commit `86316205` ändert in `modul-13-quality-gates.md` (`v6.14.1`) Zeile 272 „demselben“ zu „demselbigen“ | `vergleich b2d5f060 $M '#guard-haertung' 86316205 $M '#guard-haertung'` | `vergleich roh: …modul-13-quality-gates.md#guard-haertung <-> 86316205:…#guard-haertung cmp 1` | 1 | fällt |
| Gegenprobe Zeile der `id` (Einheit nach `ADR-0158`) | `L266-266` an beiden Seiten | `cmp 0` | 0 | Lücke aus Review F-1 sichtbar |
| Gegenprobe Änderung nur in `F` | `#guard-haertung` gegen `4db2bc85` | `cmp 0` | 0 | besteht |

**Mutation der Eingabeseite.** Jede ergibt Exit 2 mit einer gedruckten Zeile
und nie `cmp 0`:

| Mutation | gedruckte Zeile (Kern) | Exit |
|---|---|---|
| rechter Anker `#guard-haertungX` (gibt es nicht) | `vergleich: 86316205:…#guard-haertungX keine Einheit, Exit 2` | 2 |
| beide Anker unbekannt (`#nichtda`, `#auchnichtda`) | `vergleich: …#nichtda keine Einheit, Exit 2` | 2 |
| Lokator `L139-x` | `einheit: Lokator L139-x`, `vergleich: … keine Einheit, Exit 2` | 2 |
| Tag-Paar `v6.14.0-v6.14.1` | `vergleich: Tag-Paar v6.14.0-v6.14.1, Exit 2` | 2 |

**Mutation der Form.** Sie zeigt, dass die Probe an der Stelle hängt, die sie
belegen soll. Beide Fälle sind rot gegenüber dem Original:

| Mutation an einer Kopie der Form | Probe | Original | mutiert |
|---|---|---|---|
| m1: `hasid($0)` durch `0` ersetzt (`id`-Auflösung aus) | `#guard-haertung` gegen `86316205` | `cmp 1`, Exit 1 | `keine Einheit, Exit 2` |
| m2: Wächter `[ -n "$out" ] \|\| …` durch `:` ersetzt | `#nichtda` gegen `#auchnichtda` | Exit 2 | `cmp 0`, Exit 0 (fail-open) |

**Liefer-Punkt 2 (Träger nachgezogen): bestätigt, mit einem offenen Rest an
die Closure (Abschnitt 5, F-4).** Die Träger wurden gegen die ADR gelesen (Abschnitt 3). Der
Suchlauf stimmt: 12 von 12 Zeilen. Die Angabe „nicht gefunden“ in
`.harness/skills/` und `.claude/commands/` ist nachgemessen:
`git grep -niE 'Zitat-Korrektur|Referent|Vergleichseinheit|Tag-Paar|nicht messbar|ADR-0158|ADR-0159' -- .harness/skills .claude/commands`
liefert 0 Treffer.

**Gate-Pflicht:** `make gates` siehe Abschnitt 7. **Review:** Beide Reports
liegen vor. Das Re-Review hat 0 HIGH und 0 MEDIUM.

**Nicht abgehakt, korrekt offen:** Closure-Notiz, Register, Risiko-Ausgang zu
§6 und die drei Paarungen. Diese Punkte gehören zur Closure und nicht zur
Verifikation.

## 3. Entscheidungs-Konformität

- **[`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
  Entscheidung 2, Referent-Messung an `5d8855d9`.** Die bewegten Verweise stammen
  aus `git diff 5d8855d9~1 5d8855d9`. Gemessen wurde jeweils
  `vergleich 5d8855d9~1 .harness/baseline/v6.14.0/regelwerk/<datei> '<anker>' 5d8855d9 .harness/baseline/v6.14.1/regelwerk/<datei> '<anker>'`:

  | MR | Verweis | gedruckte Zeile (Kern) | Exit |
  |---|---|---|---|
  | `MR-002` | `grundlagen-source-precedence.md#vergabe-woher-die-nächste-kennung-kommt` | `vergleich roh: … cmp 0` | 0 |
  | `MR-003` | `grundlagen-durchsetzungsschicht.md#grenzen--ehrlich-benannt` | `vergleich roh: … cmp 0` | 0 |
  | `MR-004` | `modul-13-quality-gates.md#guard-haertung` | `vergleich roh: … cmp 0` | 0 |
  | `MR-001` | `grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument` | `vergleich: 5d8855d9~1:…#spec-straten-mehr-als-ein-spec-dokument keine Einheit, Exit 2`; am neuen Stand allein `einheit: leere Einheit`, Exit 2 | 2 — **nicht messbar** |

  Derselbe Anker in `grundlagen-referenz-richtung.md` löst auf: gegen sich
  selbst `cmp 0`. Das bestätigt die Aussage der ADR und des Plans §1. Der
  Datei-`cmp` nach `ADR-0157` Entscheidung 4 ergibt für alle vier MR-Dateien
  Exit 0. Damit sind Bedingung (a) und (c) maschinell belegt, und `MR-001` ist
  Urteil am Diff. Dieser Befund wird gemeldet (Review F-9) und ist der Closure
  zugewiesen.
- **`ADR-0159` Entscheidung 4, Zusagen der Befehlsform.** Für `nullglob` und
  `failglob` ist nachgefahren, dass das Ergebnis gleich bleibt. Für leere
  Einheit, ungültigen Lokator und ungültiges Tag-Paar ist nachgefahren, dass sie
  Exit 2 mit gedruckter Zeile ergeben. Nicht nachgefahren und deshalb
  **übernommen** aus der §Fitness Function: `set -euo pipefail`, die Leerzeilen
  am Ende und die Normalisierung des Tag-Paars.
- **Immutabilität (`AGENTS.md` §3.5).** `ADR-0157` ist in der Range nicht
  geändert. `ADR-0158` ist nur in ihrem Anlage-Commit `07923146` berührt
  (`git log b7d97cca..HEAD -- docs/plan/adr/`: nur `07923146` und `b36f882d`).
  Die Teilablösung läuft über `Supersedes` in `ADR-0159` und die Index-Vermerke
  `→ ADR-0158` und `→ ADR-0159`. `make doc-immutable` endet mit Exit 0.
- **`AGENTS.md` §3.5 gegen die ADR.** Der Kernsatz nennt Datei, Abschnitt
  hinter Heading oder HTML-`id`, die Tabellenzeile für eine `id` in einer
  Tabelle, zitierte Zeilen, roh und das bewegte Tag-Paar. Das stimmt mit
  `ADR-0158` Entscheidung 1 und 2 und `ADR-0159` Entscheidung 1 und 3 überein.
  Der Absatz „Beleg“ gibt `ADR-0158` Entscheidung 5 vollständig wieder (Form,
  Einheit, beide Stände, roh oder normalisiert, gedruckte Zeile) und nennt den
  Geltungsbereich. Zwei Ungenauigkeiten siehe Abschnitt 6 (INFO).
- **`.claude/agents/verifier.md` und `harness/targets/pin-stale.md`.** Beide
  nennen am Pin-Commit neben dem Datei-`cmp` die Referent-Messung je bewegtem
  Verweis. `verifier.md` verweist auf Entscheidung 4 (Befehlsform) und meldet
  „nicht messbar“. `pin-stale.md` verweist auf Entscheidung 2 und nennt den
  Slice-Plan des Bumps als Beleg-Ort. Beides entspricht `ADR-0159`
  Entscheidung 2 und der Folgepflicht unter §Konsequenzen.
- **`AGENTS.md` §3.12 (Herkunft).** Der Plan kennzeichnet die Proben des
  Architects als **übernommen** und die eigenen Proben als nachgefahren, mit
  Befehl, Stand und gedruckter Zeile. Die Ausgangslage aus dem Re-Review ist
  ausdrücklich „übernommen, nicht nachgemessen“. Die gemessenen Werte des Plans
  (Pflichtproben 1 und 2, `#guard-haertung`, die Messung an `5d8855d9`)
  stimmen an den Stellen, die dieser Lauf nachgefahren hat, in Farbe und Exit
  mit den eigenen Läufen überein. Die Commit-Kennungen der Proben weichen ab,
  weil sie an anderen Klonen entstanden sind.

## 4. Plan gegen Code-Diff

| Plan §3 | im Diff | Urteil |
|---|---|---|
| `ADR-0158` neu, Index (`07923146`) | ja | stimmt |
| `ADR-0159` neu, Index (`b36f882d`) | ja | stimmt |
| `AGENTS.md` §3.5 update | ja (`fc6d104c`, `fce2d159`) | stimmt |
| `verifier.md`, `pin-stale.md` update (Fixrunde) | ja (`fce2d159`) | stimmt |
| `implementer.md`, `.d-check.yml` nicht geändert | nicht im Diff | stimmt mit dem Plan überein; zur Begründung siehe Re-Review F-4 (Abschnitt 5) |
| Werkzeug hinter `make` nicht realisiert | kein `Makefile`/`tools/`-Diff | stimmt; die Rückführung aus §4 tritt nicht ein |
| Kein Produkt-Code (§1) | Diff berührt nur `.md` (9 Dateien) | stimmt |

Im Diff steht nichts, was der Plan nicht nennt. Plan und Reviews tragen die
beiden Review-Reports.

## 5. Offene Punkte aus den Reviews (gelesen, nicht behoben)

- Re-Review F-1 bis F-3 (LOW): gestapelte `id`-Zeilen, abweichende
  `id`-Schreibweisen und ein Fence aus vier Zeichen. Heute gibt es dafür keine
  Fundstelle. Laut Auftrag gehen sie an einen Skript-Slice; das passt zu
  `ADR-0159` Re-Evaluierungs-Trigger (a). Die Closure muss diesen Ausgang mit
  Kennung festhalten.
- Re-Review F-4 (LOW): `.claude/agents/implementer.md` Zeile 51 sagt, den
  Pin-Commit prüfe „der Verifier per `cmp`“. Plan §3 stuft die Datei weiter als
  „nicht geändert“ ein. Die Behebung liegt bei der Closure.
- Re-Review F-8 (INFO): `verifier.md` nennt „die sechs Gate-Ziele“. Der
  Gate-Lauf in Abschnitt 7 zeigt die tatsächliche Liste. Die Behebung liegt bei
  der Closure.
- Review F-9 bzw. `ADR-0159` Entscheidung 2: der gebrochene Anker in `MR-001`.
  Die Entscheidung trifft die Closure.
- Plan §3, fremde Träger: zwei `state.md` im Beobachtungs-Register beschreiben
  die Lücke noch als offen. Gemeldet, Frist ist die Closure.

## 6. Abweichungen und Hinweise

Keine DoD-Abweichung. Hinweise der Stufe INFO:

1. `AGENTS.md` §3.5, Absatz „Beleg“: Für „nicht messbar“ verweist er auf
   `ADR-0158` Entscheidung 4 und 5. Für einen **MR-Eintrag** trägt die Regel
   aber `ADR-0159` Entscheidung 2 (Zweig „übrige Abschnitte“). Im Ergebnis gilt
   dieselbe Regel, die Fundstelle ist jedoch unvollständig.
2. `AGENTS.md` §3.5, Kernsatz: „den Abschnitt hinter dem Anker (Heading oder
   HTML-`id`)“ nennt nicht die Block-Einheit einer `id` vor einem Absatz
   (`ADR-0159` Entscheidung 1, dritte und fünfte Tabellenzeile). Der Satz
   verweist auf die ADR; als Zusammenfassung ist er zu grob, aber nicht falsch
   gerichtet.
3. Außerhalb dieses Slice: `5d8855d9` ändert 12 Dateien, davon 4 MR-Dateien,
   und seine Message nennt `ADR-0051`, nicht `ADR-0073`. `ADR-0157`
   Entscheidung 4 und `verifier.md` beschreiben dagegen einen Pin-Commit, der
   „nur MR-Dateien ändert“ und `ADR-0073` nennt. Der Datei-`cmp` trägt trotzdem,
   weil er nur die MR-Dateien liest. Das ist eine Meldung an den Planner und
   betrifft den Record von `slice-harness-baseline-v6-14-1`.
4. `cmp` druckt bei Ungleichheit eine eigene Zeile auf stderr
   (`/dev/fd/63 /dev/fd/62 sind verschieden: …`). Diese Zeile kommt zusätzlich
   zu der von `vergleich` und stört die Zusagen nicht.

## 7. Gate-Lauf

Der erste `make gates`-Lauf nach dem Commit dieses Berichts (`4221a342`) endete
mit Exit 2. Grund war `docs-check`:
`d-check: 1770 Datei(en) geprüft, 2 Befund(e)`, zweimal `id-unlinked` in diesem
Bericht. Die beiden Index-Vermerke standen als nackte Kennung in
Anführungszeichen. Sie sind jetzt in Inline-Code gesetzt. Die übrigen Ziele des
Laufs waren bis dahin grün, darunter
`coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`. Das Ergebnis des
Laufs nach der Korrektur meldet der Verifier an den Aufrufer. Der Bericht friert
vor diesem Lauf ein und trägt die Zahl deshalb nicht selbst.

## Verdikt

**DoD bestätigt: ja** für Liefer-Punkt 1 und Liefer-Punkt 2. Die
Pflichtproben (Anker-Wechsel fällt, verschobener Lokator besteht,
`#guard-haertung` mit geändertem Wort fällt) sind mit der wörtlich gezogenen
Befehlsform aus `ADR-0159` Entscheidung 4 nachgefahren, ebenso je eine Mutation
der Eingabeseite und der Form mit gesehener Farbe. Die Referent-Messung an
`5d8855d9` bestätigt: drei Verweise bestehen roh, `MR-001` ist nicht messbar.
Für die Closure bleiben: Re-Review F-1 bis F-4 und F-8, `MR-001` (F-9), die
zwei `state.md` und der Risiko-Ausgang aus §6.
