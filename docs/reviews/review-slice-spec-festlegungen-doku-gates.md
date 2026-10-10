# Review-Report: slice-spec-festlegungen-doku-gates — 2026-10-10

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `942ebf3f..64c68933` (5 Commits: `9e43f968` Festlegung
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
Vertrag und Gate-Index; `780cd2b5` und `ff391aa0` Plan; `deecc466`
[`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
vom Architect; `64c68933` dessen Folgepflicht)

**Skill:** `.harness/skills/reviewer.md` @ dc04087e
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-10

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

- Slice-Plan `slice-spec-festlegungen-doku-gates` (Stand `64c68933`), §2 DoD,
  §3 Gegenprobe A1–C24, Messungen M0–M23, Anschluss-Frage, Suchlauf, §4, §6, §8
- [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
  (Entscheidung 1–4, §Konsequenzen Folgepflicht 3 an den Reviewer, Fitness Function),
  dazu die Quellen
  [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
  [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md),
  [`ADR-0094`](../plan/adr/0094-review-matrixklasse-kennung-statt-adresse.md),
  [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md),
  [`ADR-0097`](../plan/adr/0097-observation-matrixklasse-review-verboten.md),
  [`ADR-0099`](../plan/adr/0099-slice-welle-review-regel-zurueckgenommen.md),
  [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md),
  [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
  (die zitierten Entscheidungen und Status-Zeilen)
- `.d-check.yml` (Eingabe des Laufs), `d-check.mk` (Digest v0.82.0)
- keine `LH-*`-Anforderung berührt (Slice-Kopf)
- `AGENTS.md` §3.1, §3.5, §3.9, §3.12, §3.13 (Hard Rules); `harness/conventions.md` MR-000

**Messungen des Reviewers.** d-check v0.82.0 über den Digest aus `d-check.mk`,
`docker run --rm --network none -v <Klon>:/repo:ro <Digest>`, an einem Klon von
`64c68933` im Scratchpad. Basislauf: `d-check: 1831 Datei(en) geprüft, 0 Befund(e)`,
Exit 0. Jede Mutation einzeln an der Kopie (Einfügen über `sed -n … > Kopie` und
`cp`), danach `git checkout -- .`, `git status --short` am Ende leer. Gedruckt ist
die Befundzeile bzw. die Schlusszeile.

| # | Mutation (Klon) | Ergebnis |
|---|---|---|
| R1 | `- [ ] …` und der `forbid-pattern`-Satz in §7 von `done/release-image-scan.md` (geschlossener Slice-Plan, Name ohne `slice-`) | `0 Befund(e)`, Exit 0 |
| R2 | `- [ ] …` in §7 von `done/slice-001-bootstrap.md` | `section-tasks-open`, Exit 1 |
| R3 | `[SPEC-040](../spec/lastenheft.md)`, `[SPEC-040](conventions.md)`, `[LH-QA-REL-001.a](sensors/docs-check.md)`, `[LH-QA-REL-001.a](../spec/lastenheft.md)` in `harness/README.md` | je `0 Befund(e)`, Exit 0 |
| R4 | `SPEC-040` nackt an derselben Stelle | `id-unlinked`, Exit 1 |
| R5 | Link auf die abgelöste ADR 0038 aus `harness/README.md`, aus `AGENTS.md` | je `0 Befund(e)`, Exit 0 |
| R6 | derselbe Link aus der ADR 0094, aus einer `observation.md`, aus dem Slice-Plan in `in-progress/` | je `matrix-inactive`, Exit 1 |
| R7 | Token `ADR-0038` (Inline-Code und nackt) in der ADR 0094 | `0 Befund(e)`, Exit 0 |
| R8 | Link, dessen **Ziel** über einen Zeilenumbruch geht: `[Text](gibt-es-nicht` / `.md)`; `[Text](conventions` / `.md)`; `[Text](` / `gibt-es-nicht.md)` | je `target-missing`, Exit 1 (auch beim existierenden Ziel `conventions.md`) |
| R9 | Link, dessen **Ziel-Anker** umbricht: `[Text](conventions.md#gibt-es` / `-nicht)` | `anchor-missing`, Exit 1 |
| R10 | Link, dessen **Linktext** umbricht, Ziel fehlt bzw. Anker fehlt | je `0 Befund(e)`, Exit 0 |
| R11 | `forbid-pattern`-Satz in §7 von `slice-001` in Inline-Code / im Fence / als Text | 0 / 0 / `section-forbidden` |
| R12 | jedes `Steering-Loop` in `slice-001` in Inline-Code gesetzt | `section-pattern-missing`, Exit 1 |
| R13 | `ADR-0094` und `slice-099` in einem Fence in `spec/pflichtenheft.md` | `0 Befund(e)`, Exit 0 |
| R14 | `slice-099` in Inline-Code mit `<!-- d-check:status-provenance -->` in der ADR 0094 (Regel adr → slice) | `0 Befund(e)`, Exit 0 |
| R15 | absoluter Host-Pfad in Prosa / im Fence in `harness/README.md` | `hostpath-forbidden` / `0 Befund(e)` |
| R16 | Link auf eine existierende, nicht getrackte Datei | `target-untracked`, Exit 1 |
| R17 | `versions.current-from` auf eine fehlende Datei; unbekannter Schlüssel unter `scan:` | je `d-check: error: …`, Exit 2 |
| R18 | `ADR-0163` Fitness Function Zeilen 2, 7, 8 nachgefahren (Token in ADR 0094; Marker in `spec/architecture.md`; Link in ADR 0039) | `matrix-forbidden` / 0 / 0 — wie die ADR |

`git grep` am Stand `ff391aa0` mit den Befehlen der Fitness Function, Zeilen 1 und
9: 488 Zeilen, 271 Dateien, 1 Zeile in den ADRs, 0 in `spec/`, 0 Treffer in den
ADRs 0039/0041 — wie die ADR. `make suchlauf-nachmessen PLAN=<Plan>` am Arbeitsbaum
`64c68933`: Exit 0, `suchlauf-nachmessen: 16 Zeilen stimmen`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Die Festlegung sagt, die `structure`-Regel gelte der Closure-Notiz §7 **jedes** geschlossenen Slice-Plans; das Werkzeug prüft nur `docs/plan/planning/done/slice-*.md`, zwölf geschlossene Slice-Pläne mit §7 in `done/` (Name ohne `slice-`) bleiben grün (R1, R2). Die Quelle A4 (Vertrag am Parent: „je `done/slice-*.md`“) war enger als der Satz, der mit Ausgang S aus ihr wurde. | [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 1 (weicht das Werkzeug ab, ist es falsch); `AGENTS.md` §3.12 Instanz B | `spec/pflichtenheft.md` · „der Closure-Notiz §7 jedes geschlossenen Slice-Plans“; ebenso `harness/sensors/docs-check.md` §Bindung · „Closure-Notiz je geschlossenem Slice-Plan“ | ja — d-check-Lauf an einem Klon mit Mutation R1 | Festlegung breiter als das Werkzeug |
| F-2 | HIGH | Die Zeile `ids` sagt, ein Treffer sei eine Kennung „ohne Link auf das Dokument, das das Muster ihr zuordnet“; das Werkzeug meldet nur die unverlinkte Kennung, ein Link auf ein beliebiges Ziel genügt (R3, R4). Damit trägt auch der Satz „gilt das erste der Liste; deshalb steht das Muster der Verfeinerungen … vor dem Basismuster“ keine beobachtbare Wirkung, und die Anschluss-Frage „`.a`-Verfeinerung verlangt den Link ins Pflichtenheft“ ist widerlegt. | [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 1; `AGENTS.md` §3.12 Instanz B | `spec/pflichtenheft.md` · „steht in Prosa ohne Link auf das Dokument, das das Muster ihr zuordnet“; Slice-Plan §3 · „`.a`-Verfeinerung verlangt den Link ins Pflichtenheft“ | ja — Mutation R3 | Festlegung breiter als das Werkzeug |
| F-3 | HIGH | Die Zeile `matrix-inactive` sagt „eine Datei verweist auf ein Ziel“ mit verbotenem Status; das Werkzeug meldet nur einen **Link** und nur aus einer Datei, die einer `matrix`-Klasse angehört — aus `harness/README.md` und `AGENTS.md` bleibt der Link grün, ein Token nirgends gemeldet (R5–R7). | [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 1; `AGENTS.md` §3.12 Instanz B | `spec/pflichtenheft.md` · „eine Datei verweist auf ein Ziel, dessen Status in `status.forbidden` steht“ | ja — Mutation R5 | Festlegung breiter als das Werkzeug |
| F-4 | HIGH | „Geht der Linktext oder das Ziel eines Links über einen Zeilenumbruch, melden `links` und `anchors` nichts“ gilt am gepinnten Werkzeug nur für den Linktext; ein umbrochenes Ziel meldet `target-missing` bzw. `anchor-missing`, auch wenn das zusammengesetzte Ziel existiert (R8–R10). Der Satz ist aus der Messung von `slice-077` an einem älteren d-check übernommen (A25, keine M-Messung), und der neu formulierte Vertrag Grenze 9 trägt denselben Satz. | [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 1; `AGENTS.md` §3.12 (übernommener Wert ohne Kennzeichnung) | `spec/pflichtenheft.md` · „Geht der Linktext oder das Ziel eines Links über einen Zeilenumbruch“; `harness/sensors/docs-check.md` Grenze 9 · „dessen Linktext oder Ziel über einen Zeilenumbruch geht“ | ja — Mutation R8 | Festlegung breiter als das Werkzeug |
| F-5 | HIGH | Messung M1 hat laut Auftraggeber je eine Zeile per `>>` an `README.md` und `.d-check.yml` im Arbeitsbaum angehängt (danach zurückgesetzt) — Schreiben in eine Repo-Datei per Umleitung. | `AGENTS.md` §3.1 (Docker-only, Text-Umschreiben); Skill HIGH „Docker-only-Verstoß“ | Slice-Plan §3 · „M1 lief versehentlich am Arbeitsbaum (`942ebf3f`)“ | nein — **übernommen** (Angabe des Auftraggebers; der Diff trägt keine Spur, `git status` des Arbeitsbaums heute leer) | Docker-only-Verstoß (Umleitung an Repo-Datei) |
| F-6 | MEDIUM | Der Plan trägt den Vorfall aus F-5 als Ortsversehen („lief versehentlich am Arbeitsbaum“), ohne den Weg (`>>`) und die verletzte Regel zu nennen; eine Closure-Notiz, die aus §3 schöpft, verliert so den Beleg für den Zähler von `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`. | `AGENTS.md` §3.1; Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Closure- und Lerneintrag-Regeln | Slice-Plan §3 · „die Mutationen sind mit `git checkout` zurückgesetzt (`git status` danach leer)“ | nein — Lese-Handlung | Vorfall ohne Weg und Regel im Träger |
| F-7 | MEDIUM | `ADR-0163` nimmt zwei Lücken ohne Sensor an (Marker in `spec/` hebt `matrix` auf; `matrix.exempt-paths` nimmt die ADRs 0039/0041 auch von den Regeln aus, „Wächter ist das Review“), und Entscheidung 4 weist Grenzen dem Vertrag zu; der Vertrag §Grenze nennt keine von beiden, ebenso wenig, dass ein Token in einem Fence für `matrix` kein Treffer ist (R13). Ein Reviewer, der den Vertrag als Liste dessen liest, was Grün nicht deckt, findet den Marker-Ausweg nicht. | [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 2 (Grenze), 4, §Konsequenzen „Negativ“ | `harness/sensors/docs-check.md` §Grenze · „Grenze — was das Grün nicht abdeckt“ (kein Punkt zu `status-provenance` oder `exempt-paths`) | nein — Lese-Handlung | Messwerkzeug-Grenze unbenannt (fail-open) |
| F-8 | MEDIUM | §3 führt die Architect-ADR weiter als „**offen**, Frage in §6“, während §2 sie abhakt und §6 den Ausgang *entfallen* mit `ADR-0163` trägt; C8 der Gegenprobe nennt „überholt durch diesen Slice … gemeldet“ ohne Verweis auf Entscheidung 4. Zwei Aussagen desselben Plans, keine verweist auf die andere. | Skill MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“ | Slice-Plan §3 Tabelle · „`Schärft:`-Kante (Bedingung vor der Closure) — **offen**, Frage in §6“ | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-9 | LOW | `ADR-0163` Entscheidung 2 bewegt die Eigenschaft „`make docs-check` Exit 0 belegt, dass die Spec keine ADR nennt“ (der Marker hebt `spec → adr` auf); der Suchlauf hat sie nicht gesucht. Sie steht als Beleg im eigenen §6 („Die Spec darf nicht auf ADRs zeigen … *Zu belegen durch:* `make docs-check` Exit 0“) und in fünf Plänen unter `open/` (`slice-spec-festlegungen-code-gates`, `-commit-baseline-gates`, `-coverage-gates`, `-kennungs-gates`, `-pruefer-hooks`) — die fünf sind fremde Träger, Meldung an den Planner, Frist: Closure dieses Slice. | `AGENTS.md` §3.13 | Slice-Plan §6 · „**Die Spec darf nicht auf ADRs zeigen**“ | ja — `git grep -c "Die Spec darf nicht auf ADRs zeigen" -- docs/plan/planning` (7 Dateien, davon 1 in `done/`) | Suchlauf ohne die von der ADR bewegte Eigenschaft |
| F-10 | LOW | Der Vertrag sagt, er wiederhole die Festlegung nicht, wiederholt aber Randformen und Ausgänge von [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) (Grenze 8: Fence, relative Pfade, Tilde mit Benutzername; Grenze 9; §Sperren: die Ursachen von Exit 2; „Die Farbe steht in der Zeile“); F-4 zeigt, dass der doppelte Satz schon gemeinsam falsch ist. | [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 4 | `harness/sensors/docs-check.md` §Vertrag · „diese Datei wiederholt sie nicht“ | nein — Lese-Handlung | Vertrag doppelt die Festlegung |
| F-11 | INFO | Frage des Implementers zur Historie §8: eine eigene Zeile ist nicht nötig — die beiden Sätze sind Teil der noch nicht gemergten Erstfassung derselben Kennung am selben Datum, und die Zeile nennt „Randformen der Verweis-Module“. Der Pilot hat eine Änderung nach seiner Architect-ADR in die Zeile desselben Datums aufgenommen (`d2a863a1`); dieselbe Form wäre hier eine Ergänzung der bestehenden Zeile um Marker und `matrix.exempt-paths`, keine neue Zeile. | Maintainability | `spec/pflichtenheft.md` §8 · „`SPEC-040` neu (`make docs-check`:“ | nein | — |
| F-12 | INFO | Umfang: `structure` war in dieser Sitzung prüfbar (R1, R2, R11, R12 und die Regelliste gegen `.d-check.yml` gelesen); die zweite Rückführungs-Bedingung aus §4 tritt von Seiten des Reviews nicht ein. | Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | Slice-Plan §4 · „der Review-Report meldet, dass er `structure` nicht in derselben Sitzung prüfen konnte“ | nein | — |

## Folgepflicht 3 aus `ADR-0163` — Ergebnis

- **Die zwei Sätze gegen Entscheidung 2 und 3.** Wortgleich mit Folgepflicht 1.
  Der Satz zu `matrix.exempt-paths` gibt Entscheidung 3 vollständig wieder (Quelle
  ganz aus, Ziel bleibt Gegenstand; R18, Zeile 8, und Fitness Zeile 2). Der
  Marker-Satz gibt Entscheidung 2 für das Werkzeug wieder (Token aufgehoben für
  jede Regel — von mir zusätzlich an adr → slice gemessen, R14; Link bleibt
  Befund; `d-check:ignore` wirkt auf `matrix` nicht). Die **Grenze** aus
  Entscheidung 2 (in `spec/` kein zulässiger Weg) steht nicht in der Festlegung —
  folgerichtig, weil die Festlegung das Werkzeug beschreibt und die Spec keine ADR
  nennen darf; sie fehlt aber auch im Vertrag (F-7).
- **Herkunftstabelle (Entscheidung 1) gegen die Gegenprobe §3.** Jede Zeile hat
  ihre Gegenprobe-Zeilen mit passendem Ausgang: Fläche ↔ C6, C24, A15; Marker ↔
  C15; Status-Prüfung, Lineage, `exempt-paths` ↔ B7, B8, C17; Klassen und Regeln
  ↔ C14, C18, C19, B6; `versions` ↔ C20, C22, C23; `hostpaths` ↔ C2, C10, C12,
  C13; `d-check:ignore` ↔ C2. Die Liste „Nicht Quelle“ deckt sich mit den
  E-Zeilen C1, C4, C5, C9, C11, C16, C21, C23 (zweiter Teil), C24 (Liste).
  Abweichung nur im Gehalt, nicht in der Zuordnung: die Zeile „diese ADR —
  keine Quell-ADR trägt sie“ stützt sich auf übernommene Messungen, und drei
  ihrer Sätze sind am Werkzeug nicht wahr (F-1, F-2, F-4).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/pflichtenheft.md` — Matrix-Regel Spec → ADR/Slice | geprüft, ohne Befund: kein `ADR-<NNNN>`- und kein `slice-<NNN>`-Token, kein Marker im Diff; `make gates` grün |
| `spec/pflichtenheft.md` — Stratum | geprüft, ohne Befund: Harness-Festlegung in §7, keine neue `LH-*`-Anforderung |
| `spec/pflichtenheft.md` — `hostpaths`, `versions`, `tracked`, Ausgänge | geprüft, ohne Befund außer F-1–F-4 (R15–R17 bestätigen die Sätze; Zeilenform mit Tabulatoren an R2 gesehen) |
| `docs/plan/adr/0163-…` — Instanz B (`AGENTS.md` §3.12) | geprüft, ohne Befund: Mengen-Aussagen gemessen (Zeilen 1–9) oder als *hergeleitet* bzw. *übernommen* gekennzeichnet; Zeilen 1, 2, 7, 8, 9 nachgemessen, gleich |
| `docs/plan/adr/0163-…` — Zitate in die Quell-ADRs | geprüft, ohne Befund: Absatz in der ADR 0094 (Z. 110), Satz in der ADR 0097 (Z. 57), Schlüsselname in der ADR 0073, Status-Zeile der ADR 0075, Folgepflicht 2 der ADR 0160 |
| `docs/plan/adr/README.md` | geprüft, ohne Befund: Zeile ADR 0163, Vermerke an 0094/0097 |
| `harness/README.md` | geprüft, ohne Befund: Bindung trägt Vertragsdatei und Festlegung |
| `harness/sensors/docs-check.md` — Verweise | geprüft, ohne Befund: alle Links auf §7 lösen auf; keine Grund-Codes außer im Reparatur-Pfad, keine Werte 220/120, keine Exit-Tabelle |
| Slice-Plan — Suchlauf-Zahlen | geprüft, ohne Befund: `make suchlauf-nachmessen` Exit 0, 16 Zeilen |
| Slice-Plan — Gegenprobe-Zählung | geprüft, ohne Befund: 36 + 16 + 24 = 76 Zeilen (*abgeleitet* aus den Zeilennummern) |
| Commit-Messages | geprüft, ohne Befund: jede nennt eine `ADR-*`-Kennung, keine `SPEC-*` im Betreff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 5 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Festlegung breiter als das Werkzeug ·
Docker-only-Verstoß (Umleitung an Repo-Datei) · Vorfall ohne Weg und Regel im
Träger · Messwerkzeug-Grenze unbenannt (fail-open) · Nachzug widerspricht dem
Nachbarn im selben Träger · Suchlauf ohne die von der ADR bewegte Eigenschaft ·
Vertrag doppelt die Festlegung

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-4 lassen die Festlegung, die nach
`ADR-0163` Entscheidung 1 über dem Werkzeug steht, an vier Stellen etwas anderes
sagen als das gepinnte Werkzeug; F-5 ist nicht mehr behebbar, verlangt aber die
Erfassung (F-6). Ob ein Satz der Festlegung an das Werkzeug angepasst wird oder das
Werkzeug abweicht (dann eigener Slice nach §6, dritter Punkt), entscheidet der
Implementer — bei Abweichung vom Wortlaut von `ADR-0163` als Frage an den
Architect (`AGENTS.md` §3.5). Die DoD-Zeile „Review durchgeführt“ bleibt offen:
eine Fixrunde folgt.

**Übergabe:** Findings gehen an den Implementer; F-9 (fremde Träger) zusätzlich
an den Planner. Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort
in den Zähler. Dieser Report ist ein **Lauf-Beleg**; er ersetzt keine
Verifikation.
