# Review-Report: slice-harness-baseline-v6-16-0, Fixrunde — 2026-10-07

**Review-Art:** Code, enges Re-Review der Fixrunde gegen Plan, Entscheidungen
und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `85ef5f35..a445afa7` (Commits `c88a766b`, `a445afa7`),
Antwort auf [`review-slice-harness-baseline-v6-16-0.md`](review-slice-harness-baseline-v6-16-0.md)
(F-1 bis F-8)

**Skill:** `.harness/skills/reviewer.md` @ a445afa7
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
> Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-baseline-v6-16-0`, Abschnitt „Fixrunde“ (Stand `a445afa7`)
- [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
  (Entscheidung 3 bis 5, Folgepflicht 1 bis 8) und die Vorgabe des Architect zur
  Fixrunde (vom Auftraggeber wörtlich übernommen, nicht als Artefakt im Repo):
  §3.5-Klammer „für die Form-Korrektur an einem MR-Eintrag `ADR-0161`
  Entscheidung 4“ samt Einfrier-Satz, `zitat-vergleich.md` Zeile 10 „an einem
  MR-Eintrag“, Verifier liest den Durchgangs-Beleg und misst den
  `formnorm`-`cmp`, `versions`-Kommentar mit der Immutabilität
- `AGENTS.md` §3.5, §3.7, §3.9, §3.13
- keine `LH-*`-Kennung berührt

---

## Findings

Spalten nach dem §Output-Schema des Reviewer-Skills
(`v6.16.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| R-1 | LOW | §3.5 verlangt für die Zitat-Korrektur an einem MR-Eintrag je Verweis die `zitat-vergleich`-Zeile **und** den `formnorm`-`cmp`; im selben Absatz steht „die Referent-Messung je MR-Eintrag gehört in den Adaptions-Durchgang“. Die Bump-Träger `harness/targets/pin-stale.md` (Form-Commit vor dem Löschen) und `.claude/agents/implementer.md` nennen am Form-Commit nur `teilrange` und `formnorm`-`cmp`. Ob die Messung derselben Adresse an `F~1` und `F` Pflicht ist, liest ein Implementer des nächsten Bumps nur aus §3.5; dieser Plan trägt beide Belege, ein Schaden ist in diesem Slice nicht eingetreten. | `AGENTS.md` §3.5; [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) Entscheidung 4, §Konsequenzen („Referent-Messung als Bedingung … entfallen“) | `AGENTS.md` · „dazu je MR-Datei die Zeile des `formnorm`-`cmp` am“; `harness/targets/pin-stale.md` · „und am Form-Commit je MR-Datei den `formnorm`-`cmp`“ | nein — Lese-Handlung | Implementierung weicht von ADR-Wortlaut ab |
| R-2 | INFO | Die vier neuen `zitat-vergleich`-Zeilen messen dieselbe Adresse an `3f88e0c2~1` und `3f88e0c2`; `F` ändert nur MR-Dateien, die Einheit unter `.harness/` ist an beiden Ständen dieselbe Datei. `cmp 0` folgt aus der Konstruktion und belegt darüber hinaus, dass Adresse und Anker auflösen (kein Exit 2). Nachgefahren: drei der vier Zeilen, je `cmp 0`, make-Exit 0. | `ADR-0161` Entscheidung 3 | Plan · „`make zitat-vergleich` mit derselben Adresse an `3f88e0c2~1`“ | ja — `make zitat-vergleich` | — |
| R-3 | INFO | Der Fixrunde-Eintrag F-1 sagt „Beide Träger verlangen für den Form-Commit dieselben zwei Belege“; `harness/targets/zitat-vergleich.md` ist der Vertrag eines Werkzeugs und nennt den `formnorm`-`cmp` nicht. Der Widerspruch aus F-1 ist trotzdem aufgelöst (beide nennen den MR-Eintrag als Gegenstand der Messung). | Maintainability | Plan · „Beide Träger verlangen für den Form-Commit dieselben“ | nein | — |
| R-4 | INFO | Die Datei-Tabelle des Plans nennt für `AGENTS.md` §3.5, `zitat-vergleich.md`, `verifier.md`/`implementer.md` und `.d-check.yml` nur `fff016fb`; die Fixrunde-Änderungen aus `c88a766b` stehen allein im Abschnitt „Fixrunde“. | Maintainability | Plan · „Folgepflicht 4 bis 7 von `ADR-0161`“ (Zeile der Datei-Tabelle) | nein | — |

## Stand der Findings aus dem ersten Review

| Finding | Stand |
|---|---|
| F-1 (MEDIUM) | gelöst — §3.5 ersetzt den Verweis (Folgepflicht 6) mit der Klammer der Architect-Vorgabe und dem Einfrier-Satz; `zitat-vergleich.md` Zeile 10 „an einer `Accepted` ADR und an einem MR-Eintrag“. Die beiden Träger widersprechen sich beim Beleg nicht mehr. Die vier `zitat-vergleich`-Zeilen stehen im Plan; Rest-Unschärfe der Bump-Träger: R-1 |
| F-2 (MEDIUM) | gelöst — §1 Bezug mit `ADR-0161`, Klammer zu `ADR-0157` ohne Pin-Commit, §1 Ziel und Ausschluss Records nach `ADR-0161`/`ADR-0156` Entscheidung 3, Symlink-Absatz mit gemessenem Ergebnis; `grep -c` im Plan: „steht noch aus“ 0, „Pin-Commit“ 5 — alle abgrenzend („Form-Commit statt Pin-Commit“, „Kein Pin-Commit“, „an Stelle des Pin-Commits“), im Halt-Abschnitt vor der Auflösung oder im `ADR-0160`-Messbeleg |
| F-3 (LOW) | gelöst — `verifier.md` „du liest ihren Beleg im Plan des Bumps, du fährst sie nicht nach“, `implementer.md` „fährt den `formnorm`-`cmp` nach und liest den Beleg des Durchgangs“, gleich der Vorgabe |
| F-4 (LOW) | gelöst — der Kommentar steht direkt über `exempt-paths:`, der `vcs`-Satz im Block über `paths:`; je Block eine Kennung (`ADR-0073`, `ADR-0161`, `ADR-0161`) |
| F-5 (INFO) | gelöst — Begründung mit der Immutabilität und „durch einen Bump außerhalb der Datei“ |
| F-6 (INFO) | begründet nicht geändert (Core-Immutabilität von `MR-005`); angenommen |
| F-7, F-8 (INFO) | an die Closure; unverändert offen wie benannt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `zitat-vergleich` am Form-Commit | geprüft, ohne Befund — `make zitat-vergleich` für `MR-002`, `MR-003`, `MR-004` (`3f88e0c2~1` gegen `3f88e0c2`, roh) je gedruckt `… cmp 0`, make-Exit 0; Anker gleich den Inline-Code-Pfaden in den MR-Dateien an `3f88e0c2` (R-2) |
| `.d-check.yml` Blöcke `versions:`/`vcs:` | geprüft, ohne Befund — Stellung und je Block eine Kennung (§3.7); Kommentar trägt die Immutabilität nach Vorgabe |
| `.claude/agents/{verifier,implementer}.md` | geprüft, ohne Befund außer R-1 |
| Plan, Abschnitt „Fixrunde“ und §1 | geprüft, ohne Befund außer R-3, R-4; kein neuer Widerspruch zu DoD oder §6-Ausgängen gefunden |
| Suchlauf | geprüft, ohne Befund — `make suchlauf-nachmessen PLAN=…` Exit 0, „18 Zeilen stimmen“ |
| Kommentar-Herkunft | geprüft, ohne Befund — `make kommentar-kennungen DIFF=85ef5f35` Exit 0, kein Kandidat |
| Commit-Messages der Fixrunde | geprüft, ohne Befund — beide nennen `ADR-0161` |
| Docker-only, §3.15 | geprüft, ohne Befund — kein Host-Werkzeug am Repo; das Werkzeug Edit ist in diesem Lauf nicht verfügbar, Textänderungen liefen über Write |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Implementierung weicht von ADR-Wortlaut ab

## Verdikt

**Merge-blockierend:** nein. F-1 bis F-6 sind gelöst oder begründet; R-1 und
die INFO-Findings gehen an die Closure und lösen keine weitere Fixrunde aus.
Die DoD-Zeile „Review durchgeführt“ ist im selben Commit auf `[x]` gesetzt
(Skill §DoD-Checkbox-Nachzug ohne Fixrunde). Kein Konflikt-Pfad über den
Architect nötig. Dieser Report ersetzt keine Verifikation (Modul 11).
