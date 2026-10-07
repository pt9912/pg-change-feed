# Review-Report: slice-spec-festlegungen-harness-werkzeuge — 2026-10-07

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `abbe11b4..ccaf329f` (10 Commits: die fünf Folge-Slices in
`open/`, der Neuschnitt auf Zuschnitt A, §7 im Pflichtenheft, der Nachzug der
Verträge, der Umzugs-Commit `8e00e831`, `MR-006`, der Plan-Nachzug,
[`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
vom Architect und dessen Folgepflicht in `d2a863a1`)

**Skill:** `.harness/skills/reviewer.md` @ ccaf329f
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

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

- Slice-Plan `slice-spec-festlegungen-harness-werkzeuge` (Stand `ccaf329f`,
  Planänderung auf Zuschnitt A vom 2026-10-07; Zuschnitt und Nachfolger `MR-006`
  laut Auftraggeber freigegeben — **übernommen**, kein Artefakt im Repo)
- [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
  (Kante, Entscheidung 2, Umzug), dazu die Quellen
  [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
  [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md),
  [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md),
  [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
  (Entscheidungen und Status-Zeilen)
- Architect-Verdikt zum Werkzeug von `slice-zitat-vergleich-werkzeug` §2 bis §5
- Baseline `v6.16.0` · `regelwerk/grundlagen-referenz-richtung.md`
  §Spec-Straten: mehr als ein Spec-Dokument · `regelwerk/modul-03-spec.md`
  §Ziel-Form: Spezifikation · `templates/spec/spezifikation.template.md` §7/§8 ·
  `templates/harness/sensors/gate.template.md` · `templates/harness/README.template.md` ·
  `templates/docs/plan/adr/README.template.md` ·
  `templates/harness/conventions/MR-NNN-titel.template.md`
- Alte Fassung des Vertrags `harness/targets/zitat-vergleich.md` am Stand
  `abbe11b4`, Zeile für Zeile gegen die Festlegung
  [`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
  gelesen
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.5, §3.7, §3.9, §3.12, §3.13)
- Keine `LH-*`-Anforderung berührt (Plan, Bezug).

**Eigene Läufe dieses Reviews** (Exit je direkt gesichert, `AGENTS.md` §3.9):

| Lauf | Ergebnis |
|---|---|
| `make suchlauf-nachmessen PLAN=<Plan>` | Exit 0, gedruckt `suchlauf-nachmessen: 13 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=abbe11b4` | Exit 0, keine Ausgabe (kein Kandidat; der Diff trägt keinen Code) |
| `make test-zitat-vergleich` | Exit 0, gedruckt `297 Fälle bestanden (je Runde 99, Runden: ohne Option, nullglob, failglob)` |
| `make doc-immutable RANGE=fcff30ec..HEAD` | Exit 2, `1 Befund(e)`, `core-drift-vcs` an `harness/conventions/MR-001-…` |
| `make doc-immutable RANGE=fcff30ec..8e89831d` (= `B..M~1`) | Exit 0, `0 Befund(e)` |
| `make doc-immutable RANGE=8e00e831..HEAD` (= `M..H`) | Exit 0, `0 Befund(e)` |
| `umzug` wörtlich aus [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 3, per `awk` in eine Datei im Scratchpad gezogen, an einem Klon im Scratchpad | `8e00e831` Exit 0; `dc04087e` Exit 1; eigene Mutationen als Probe-Commits: Umzug von `MR-002` plus neue Datei Exit 1; Nicht-MR-Datei nach `done/` Exit 1; **Merge-Commit Exit 2** (in der ADR *hergeleitet*, hier gefahren), sein erster Parent (reiner Umzug von `MR-003`) Exit 0; Umzug von `MR-004` **mit Moduswechsel** `100644 → 100755` **Exit 0** (F-9); `8e00e831` unter `color.ui=always` und unter `diff.renames=false` je Exit 0 |
| `tools/harness/zitat-vergleich.sh` an einer Probe-Datei im Klon (`id`-Zeile ohne Inhalt Zeile 5, Absatz Zeilen 7–9, Heading Zeile 10) | `#x` gegen `L7-9` `cmp 0`, gegen `L5-9` und `L6-9` je `cmp 1` — die Wortlaut-Berichtigung aus Entscheidung 2 stimmt mit dem Werkzeug |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Festlegung lässt die Lesung von Exit 2 bei einer gewöhnlichen Zitat-Korrektur offen: der Satz „Ist eine Einheit leer oder nicht lesbar …, besteht die Korrektur nicht“ aus [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 3 fehlt, obwohl [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 1 diese Entscheidung als Quelle von „Stände“ nennt; und „Mehrdeutig heißt „nicht messbar““ ist an keinen Zweig gebunden, weil der Absatz „nicht messbarer Referent“ nur die Adresse außerhalb des Repos und die an keinem Stand auflösende nennt (der alte Vertrag band „mehrdeutig“ an „Urteil am Diff“). Ein Implementer mit Exit 2 kann daraus „Urteil am Diff, besteht“ lesen, wo die Quell-ADR „besteht nicht“ sagt. | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 3; [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 1 | `spec/pflichtenheft.md` · „Mehrdeutig heißt „nicht messbar“.“ (§7, [`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)) | nein — Lese-Handlung gegen die Quell-ADR | Festlegung lässt eine Randform offen (`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`) |
| F-2 | MEDIUM | Der Vertrag führt zwei Maßstäbe: er nennt die Punkte aus der Review-Fixrunde weiter „Auslegung im Sinn von Verdikt §2“ (und „die drei letzten Zeilen“ — es sind fünf Zeilen, nicht die letzten drei) und hält „Weicht dieses Werkzeug außerhalb der Punkte unter §Abweichungen von ihm ab, ist das ein Fehler des Werkzeugs“ fest; [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 2 nimmt alle zehn Punkte als Entscheidung an und lässt Verdikt §2 nur noch mit „weicht es von [`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) ab, ist das Werkzeug falsch“ gelten. Wer am Werkzeug ändert, liest im Vertrag den historischen Block als Maßstab und die fünf Punkte als zurücknehmbar. | [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 2 | `harness/targets/zitat-vergleich.md` · „die drei letzten Zeilen aus der Review-Fixrunde, als Auslegung im Sinn von“ und „Weicht dieses Werkzeug außerhalb der“ | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`) |
| F-3 | MEDIUM | Die Übergabe des `formnorm`-`cmp` an `slice-spec-festlegungen-pruefer-hooks` ([`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 4, Plan §1 Ausschluss) hat keine Annahme in der Adresse: der Plan des Folge-Slice (im Diff, `cc55be92`) nennt `formnorm` an keiner Stelle, und sein Liefer-Punkt 1 zählt nur Make-Ziele, Guard und Hook auf. Nach dessen Closure stünde die Festlegung weiter nur in der Befehlsform von [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md). | [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 4; Skill §HIGH, Probe „deckt die Adresse den Gegenstand?“ | `docs/plan/planning/open/slice-spec-festlegungen-pruefer-hooks.md` · „je eine Zeile mit eigener `SPEC-<NNN>` für `make suchlauf-nachmessen`“ | ja — `git grep -c formnorm -- docs/plan/planning/open/slice-spec-festlegungen-pruefer-hooks.md` (Exit 1, kein Treffer, gemessen am Stand `ccaf329f`) | Aufschub-Adresse nimmt Sendung nicht an (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`) |
| F-4 | LOW | Die Message des Umzugs-Commits nennt `MR-001` und [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md), nicht [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md), wie Entscheidung 3 Spiegelstrich 2 verlangt; der Commit liegt vor der ADR (`f2d6f194`), `umzug` prüft die Message nicht (Exit 0 gemessen), und Plan §3 benennt den Altfall. Einschätzung zu Punkt 3 unten. | [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 3 | Commit `8e00e831` · „chore(harness): MR-001 nach conventions/done/“ | ja — `git log -1 --format=%B 8e00e831` | Regel-Altfall vor der Regel |
| F-5 | LOW | `make doc-immutable` (d-check Modul `vcs`) hat keine Folge-Adresse: [`SPEC-039`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) legt nur den Leer-Test davor fest, `slice-spec-festlegungen-doku-gates` nimmt „alle Module der Liste `modules:`“ (dort steht `vcs` nicht), und keiner der fünf Folge-Pläne nennt `doc-immutable` oder `vcs` (0 Treffer); Plan §6 sagt dagegen, „die übrigen Gates, Prüfer und Hooks“ gingen an die fünf. | Plan §6 „Lücke in §7 bis zu den Folge-Slices“; `v6.16.0` · `regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten | `docs/plan/planning/in-progress/slice-spec-festlegungen-harness-werkzeuge.md` · „die übrigen Gates, Prüfer und Hooks führen Schwelle und“ | ja — `git grep -n -i 'doc-immutable\|vcs' -- docs/plan/planning/open/` (0 Zeilen am Stand `ccaf329f`) | Abgrenzung ohne Adresse |
| F-6 | LOW | Der Vertrag trägt weiter einen Exit-2-Fall, den die Exit-Tabelle der Festlegung nicht nennt: der Fehler von `sed` am Zeilen-Lokator und von `awk` am Anker; die Vorlage sagt für den Vertrag „Keine Schwelle und keine Randform hier“, Liefer-Punkt 2 „trägt … Ausgänge nicht mehr selbst“. | `v6.16.0` · `templates/harness/sensors/gate.template.md` §Vertrag; Plan Liefer-Punkt 2 | `harness/targets/zitat-vergleich.md` · „Zeilen-Lokator und von `awk` am Anker; beides endet mit Exit 2“ | nein — Lese-Handlung | Randform bleibt im Vertrag |
| F-7 | LOW | §Grenzen verweist für „mehrdeutig“ auf „§Einheit“; dieser Abschnitt trägt die Mehrdeutigkeit nach dem Nachzug nicht mehr und verweist selbst nur weiter auf die Festlegung (ein Sprung zu viel, innerhalb derselben Datei). | Maintainability | `harness/targets/zitat-vergleich.md` · „(„mehrdeutig“, siehe §Einheit)“ | nein | Zitat nennt die falsche Stelle (`BEO-PGC/zitat-nennt-die-falsche-stelle`, innerhalb einer Datei) |
| F-8 | LOW | `MR-006` nimmt in die Adaption den Satz „Inhalt und Struktur folgen der Spezifikations-Vorlage der Baseline: Abschnitte 1–8 …“ auf — die Art Bestandsangabe, deren Drift `MR-001` zur Ablösung zwang; der Eintrag ist unter `vcs` immutabel, eine Änderung der Gliederung der Vorlage erzwingt wieder einen Nachfolger. | Maintainability; `MR-006` „Löst auf“ (Begründung der Ablösung) | `harness/conventions/MR-006-technik-dokument-heisst-pflichtenheft.md` · „Spezifikations-Vorlage der Baseline: Abschnitte 1–8“ | nein | Adaption trägt beweglichen Baseline-Bestand |
| F-9 | INFO | `umzug` gibt einen Umzug mit Moduswechsel (`100644 → 100755`) als reinen Umzug aus (Exit 0, `R100`); Entscheidung 4 sagt, die Randformen stünden „vollständig oben“, der Moduswechsel steht dort nicht. Gemessen am Klon im Scratchpad; der Merge-Commit (in der ADR *hergeleitet*) endet gemessen mit Exit 2. | [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 3/4 | `docs/plan/adr/0162-…` · „ihre Randformen stehen vollständig oben“ | ja — Probe-Commit mit `chmod +x` am Klon | ADR-Aussage breiter als ihre Messung (Nachbar von `BEO-PGC/adr-aussage-breiter-als-ihre-messung`) |
| F-10 | INFO | Träger außerhalb des Diffs (`AGENTS.md` §3.13, Meldung an den Planner, Frist: Closure dieses Slice): der Kommentar am Block `vcs:` in `.d-check.yml` nennt als umgangenen Commit nur die Form-Korrektur; seit [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 3 umgeht die Teil-Range auch den Umzugs-Commit. | `AGENTS.md` §3.13 | `.d-check.yml` · „Teil-Range samt Leer-Test umgeht“ | nein | Arbeit überholt stehenden Träger |
| F-11 | INFO | Die Fitness-Function-Zeile `matrix`/`ids` nennt ihr Ergebnis „im Architect-Bericht“, also außerhalb des Repos; der Gate-Lauf dieses Reviews (unten) trägt dieselbe Aussage für den Stand mit Report. | `AGENTS.md` §3.12 | `docs/plan/adr/0162-…` · „Ergebnis im Architect-Bericht“ | ja — `make docs-check` | Beleg außerhalb des Repos |
| F-12 | INFO | Drei Träger nennen für den Umzugs-Commit die Message-Pflicht („Message nennt die Kennung und `ADR-0162`“: `harness/conventions.md`, Implementer- und Verifier-Briefing), der Absatz „Umzug eines abgelösten Eintrags“ in `harness/targets/pin-stale.md` nicht. | [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md) Entscheidung 3 | `harness/targets/pin-stale.md` · „**Umzug eines abgelösten Eintrags.**“ | nein | — |
| F-13 | INFO | Außerhalb des Diffs: `make doc-immutable` steht in keiner Zeile von `harness/README.md` §Sensors, obwohl [`SPEC-039`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) und mehrere Verträge es als Lauf des Verifiers nennen (Ziel existiert, `d-check.mk`). Planner-Sache, gehört zu F-5. | `AGENTS.md` §4 | `harness/README.md` §Sensors (kein Treffer `doc-immutable`) | ja — `grep -c doc-immutable harness/README.md` | — |

### Einschätzung zu Punkt 3 — Umzugs-Commit `8e00e831` ohne `ADR-0162`

**Benannter Altfall, kein Neuschnitt** (F-4, LOW). Gründe:

1. Die Pflicht entstand mit [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
   (`f2d6f194`), zwei Commits **nach** `8e00e831`. Eine Regel bindet Commits,
   die nach ihr entstehen; der Commit verletzt keine Regel, die zu seiner Zeit
   galt.
2. Der Zweck der Message-Pflicht — dass ein Leser den Commit als
   Umzugs-Commit erkennt — ist hier anders getragen, und zwar stärker: die ADR
   selbst nennt `8e00e831` in ihrer Fitness Function als gemessenen
   Umzugs-Commit (Exit 0), der Plan §3 nennt ihn mit Altfall-Vermerk und Beleg.
   `make commit-traceability` ist grün (die Message nennt `ADR-0161`).
3. Der Exit-1-Zweig von Entscheidung 3 („wird dann neu geschnitten“) gilt dem
   **nicht reinen Umzug**; die Message ist kein Kriterium von `umzug`. Gemessen:
   `8e00e831` Exit 0.
4. Ein Neuschnitt schriebe die sechs Commits ab `8e00e831` um; dabei änderten
   sich die Kennungen `8e00e831`, `dc04087e`, `54b0d398`, die die
   `Accepted`-ADR in ihrer Fitness Function als Messstände nennt. Deren
   Berichtigung wäre eine Zitat-Korrektur an einer `Accepted`-ADR
   (`AGENTS.md` §3.5) für einen Gewinn, den Grund 2 schon trägt.

Der Altfall gehört in die Closure-Notiz §7 (Plan §3 trägt ihn bereits).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/pflichtenheft.md` §7: Struktur gegen `templates/spec/spezifikation.template.md` §7 (Regel-Absatz, Tabelle `ID · Werkzeug · Festlegung`), kein ADR-/Slice-Verweis im Körper, Kennungen `SPEC-038`/`SPEC-039` fortlaufend | geprüft, ohne Befund |
| `spec/pflichtenheft.md` §1-Verweis „§2 bis §7“, Umnummerierung §8 Historie, Historie-Zeile (Inhalt, ohne ADR/Slice) | geprüft, ohne Befund |
| alte Fassung `harness/targets/zitat-vergleich.md` (`abbe11b4`) gegen [`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge): Referenz-Formen, Heading/Slug/Fence, HTML-`id` je Stellung, Zeile ohne Inhalt, mehrdeutige Stellung, roh, Schluss-Umbruch, Tag-Paar, fremder Pin, Exit-Tabelle, alle zehn Punkte der Abweichungs-Tabelle | geprüft — jede Randform steht in der Festlegung oder bleibt im Vertrag; einzige Lücke die Folge von Exit 2 (F-1), der Rest F-6 |
| Wortlaut-Berichtigung der Stellungs-Tabelle gegen [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 und das Werkzeug | geprüft, ohne Befund (gemessen, s. o.) |
| [`SPEC-039`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) gegen [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md) Entscheidung 1 und [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) Entscheidung 4 erster Spiegelstrich, Umzugs-Commit nach Entscheidung 3 von `ADR-0162` | geprüft, ohne Befund |
| Referenz-Richtung: Spec ohne Abwärts-Kante, Kante aufwärts über `Schärft:` von `ADR-0162`; Link im ADR-Index auf §7 (Richtung ADR → Spec) | geprüft, ohne Befund |
| `ADR-0162` als Ergänzung ohne `Supersedes`: keine der zehn Festlegungen widerspricht einer `Accepted`-Entscheidung (Lokator, Slug, Erkennungsregel, fail-closed, Kommentar in `id`-Zeile je gegen die Quell-ADR gelesen); der Umzug ergänzt die Teilung von `ADR-0161` | geprüft, ohne Befund |
| `ADR-0162` Herkunftstabelle und Entscheidung 1 (a) gegen die Status-Zeilen von `ADR-0159` bis `ADR-0161` (abgelöst: `ADR-0158` E2/E6, `ADR-0159` E2, `ADR-0160` E2; `ADR-0159` E3 trägt die bleibenden Sätze von `ADR-0158` E2; der Satz „leer oder grün“ folgt aus `ADR-0160` E1 (b)/(c)) und gegen Verdikt §2 bis §4 (fünf Punkte aus §3/§4) | geprüft — stimmt; einzige Abweichung: die Zeile „Stände“ nennt `ADR-0158` E3, dessen letzter Satz in der Festlegung fehlt (F-1). Die vom Auftraggeber genannten „vier Korrekturen des Architects“ liegen nicht als Artefakt im Repo vor; geprüft wurde die Tabelle selbst |
| `ADR-0162` gemessen/hergeleitet getrennt (Option (b), Merge-Commit *hergeleitet*; Mutationen mit Stelle und Instanz) | geprüft — getrennt; F-9, F-11 |
| `MR-006` gegen `MR-001` (Adaption, Geltungsbereich, Begründung, Trigger gleich) und gegen `templates/harness/conventions/MR-NNN-titel.template.md` (Felder `Löst auf` auf die Index-Zeile, `Ausgelöst durch Baseline-Stand`, Template-Block entfernt, Link mit Anker in die vendored Fassung) | geprüft, ohne Befund außer F-8 |
| `harness/conventions.md` Index: `MR-006` unter *Aktive*, `MR-001` mit mitgezogenem Anker `mr-001` unter *Aufgelöste*; Regel zum Auflösen im Adaptions-Block | geprüft, ohne Befund |
| Verweise auf den alten Pfad von `MR-001` außerhalb von `docs/reviews/**`, `done/`-Records und `.harness/baseline/**` (`git grep conventions/MR-001`) | geprüft, ohne Befund — nur Messausgaben in `ADR-0162` und Plan |
| Briefings `.claude/agents/{architect,planner,reviewer}.md`, `.harness/skills/reviewer.md` auf `MR-006`; `implementer.md`, `verifier.md` mit Umzugs-Commit | geprüft, ohne Befund |
| `harness/README.md` (T6 Zeile `make zitat-vergleich`, T7 Kommentar wörtlich wie Vorlage, Bindungs-Liste „Spec-Kennung“), ADR-Index (T4, Zeile `ADR-0162`) gegen die Vorlagen `v6.16.0` | geprüft, ohne Befund |
| `harness/targets/pin-stale.md` (Adaptions-Durchgang und Leer-Test verweisen auf die Festlegungen, Befehlsform bleibt mit ADR-Zeiger) | geprüft, ohne Befund außer F-12 |
| fünf Folge-Slices in `open/` (je Werkzeug-Gruppe, Übergabe R2/R7/R9/T6/T8 in §1) | geprüft — F-3, F-5 |
| Suchlauf §3 des Plans: Zahlen (`make suchlauf-nachmessen`) und Muster (zusätzlich groß/klein-unabhängig `Pflichtenheft … §7`, `Abschnitte 1-7` mit Bindestrich, `§1 bis §7`) | geprüft, ohne Befund — keine weitere Fundstelle |
| Hard Rules: §3.3 (Umzug rein, `R100`, eigener Commit; Index im Folge-Commit), §3.5 (keine `Accepted`-ADR geändert), §3.1 (keine Host-Toolchain im Diff), §3.7 (kein Code-Kommentar im Diff) | geprüft, ohne Befund |
| Traceability der Commit-Messages `abbe11b4..ccaf329f` (je eine `ADR-*`, keine `SPEC-*` im Betreff) | geprüft, ohne Befund außer F-4 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 5 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Festlegung lässt eine Randform offen
(`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` — mit diesem
Slice das dritte Auftreten, falls die Closure F-1 zählt) · Nachzug widerspricht
dem Nachbarn im selben Träger · Aufschub-Adresse nimmt Sendung nicht an ·
Regel-Altfall vor der Regel · Abgrenzung ohne Adresse · Randform bleibt im
Vertrag · Zitat nennt die falsche Stelle · Adaption trägt beweglichen
Baseline-Bestand · ADR-Aussage breiter als ihre Messung · Arbeit überholt
stehenden Träger · Beleg außerhalb des Repos

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-3 (MEDIUM) brauchen eine Fixrunde am
Implementer; F-1 betrifft genau das §6-Risiko „Festlegung lässt eine Randform
offen“. Die LOW-Findings können mitgehen, F-4 ist ein benannter Altfall ohne
Neuschnitt. Die DoD-Zeile „Review durchgeführt“ bleibt offen (Nachzug bei
Schritt 21 des Implementer-Ablaufs).

**Übergabe:** Findings gehen an den Implementer; F-5, F-10 und F-13 sind
Planner-Sache (Träger außerhalb des Diffs bzw. Adresse eines Folge-Slice). Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Kein HIGH mit Rollen-Widerspruch: F-1 und F-2 messen am Wortlaut von
`ADR-0158` Entscheidung 3 und `ADR-0162` Entscheidung 2, nicht an einer neuen
Entscheidung; ein Konflikt-Pfad über den Architect ist nicht nötig, solange der
Implementer F-1 durch Fortschreiben der Festlegung nach der Quell-ADR löst.
Würde F-1 statt dessen eine andere Lesung von Exit 2 festlegen als
`ADR-0158`, wäre das eine Semantik-Änderung und ginge an den Architect. Dieser
Report ersetzt keine Verifikation (Modul 11).
