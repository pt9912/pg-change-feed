# Review-Report: slice-spec-festlegungen-doku-gates, Re-Review — 2026-10-10

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `72e57998..2ac39c31` (3 Commits: `7f139836` Fixrunde;
`037a9d32` Architect-Verdikt `architect-verdict-matrix-inactive-nur-link` zu F-3;
`2ac39c31` Nachzug des Verdikts). Grundlage ist der erste Report
[`review-slice-spec-festlegungen-doku-gates`](review-slice-spec-festlegungen-doku-gates.md)
(`72e57998`, F-1 bis F-12).

**Form dieses Laufs.** Der Auftrag nannte einen Abschnitt „Re-Review“ im ersten
Report. Der Reviewer-Skill (§Output-Schema) und die Report-Vorlage verlangen einen
Report je Lauf, Folgeläufe als neue Datei; der Skill geht vor. Dieser Report ist der
Abschnitt „Re-Review“, als eigene Datei.

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

**Eingangs-Kontext:**

- Slice-Plan `slice-spec-festlegungen-doku-gates` (Stand `2ac39c31`), §2, §3
  (Messungen M24 bis M33, Gegenprobe, Anschluss-Frage, Suchlauf, Fixrunde), §6
- Architect-Verdikt `architect-verdict-matrix-inactive-nur-link` (§0 bis §4)
- [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
  Entscheidung 1 bis 4; `ADR-0095` §Kontext, §Entscheidung, Fitness Function
- `spec/pflichtenheft.md` §7 (`SPEC-040`) und §8; `harness/sensors/docs-check.md`
- `AGENTS.md` §3.1, §3.5, §3.9, §3.12, §3.13; `harness/conventions.md` MR-000
- keine `LH-*`-Anforderung berührt

**Messungen des Reviewers.** d-check v0.82.0 über den Digest aus `d-check.mk`,
`docker run --rm --network none -v <Klon>:/repo:ro <Digest>`, an einem Klon von
`2ac39c31` im Scratchpad. Basislauf: `d-check: 1833 Datei(en) geprüft, 0 Befund(e)`,
Exit 0. Jede Mutation einzeln an der Kopie (Text per `awk`/`cat … > Scratch-Datei`
und `cp`), danach `git checkout -- . && git clean -fdq`, `git status --short` leer.
Gedruckt ist die Befundzeile bzw. die Schlusszeile.

| # | Mutation (Klon `2ac39c31`) | Ergebnis | zu |
|---|---|---|---|
| Q1 | `- [ ] …` in §7 von `done/release-image-scan.md` | `0 Befund(e)`, Exit 0 | F-1 |
| Q2 | `- [ ] …` in §7 von `done/slice-001-bootstrap.md` | `section-tasks-open`, Exit 1 | F-1 |
| Q3 | Kopie von `slice-001` mit `- [ ] …` in §7 als neue, indizierte Datei `done/welle-10/slice-999-probe.md` | 20 × `target-missing` (Link-Tiefe der Kopie), **kein** `section-tasks-open` — die Unterverzeichnis-Hälfte des neuen Satzes hält | F-1 |
| Q4 | `[SPEC-040](harness/conventions.md)` und `[LH-QA-REL-001.a](spec/lastenheft.md)` in `README.md` | `0 Befund(e)`, Exit 0 | F-2 |
| Q5 | `SPEC-040` nackt in `README.md` | `id-unlinked`, Exit 1 | F-2 |
| Q6 | Link auf die ADR 0038 (Status `Superseded by …`) aus `AGENTS.md` | `0 Befund(e)`, Exit 0 | F-3 |
| Q7 | derselbe Link aus der ADR 0094 | `matrix-inactive`, Exit 1 | F-3 |
| Q8 | Token `ADR-0038` in Inline-Code und nackt in `done/welle-12.md` (Klasse `welle`) | nur `id-unlinked` (der nackte), **kein** `matrix-inactive` | F-3, Verdikt §4 („an `welle` *hergeleitet*“) |
| Q9 | Link auf die ADR 0038 aus `done/welle-12.md` | `matrix-inactive`, Exit 1 | F-3 |
| Q10 | `**Status:** Superseded` in `done/slice-001-bootstrap.md`, Link darauf aus `done/welle-12.md` | `matrix-inactive`, Exit 1 — Zielklasse `slice` | F-3, Verdikt §4 („als Ziel … `slice` *hergeleitet*“) |
| Q11 | Link, Ziel umbrochen: `[Text](harness/conventions` / `.md)` (zusammengesetzt existiert es) | `target-missing`, Exit 1 | F-4 |
| Q12 | Link, Anker umbrochen, Anker fehlt: `…#gibt-es` / `-nicht)` | `anchor-missing`, Exit 1 | F-4 |
| Q13 | Link, Anker umbrochen, zusammengesetzt existiert: `…#adaptions` / `-block)`; Kontrolle einzeilig | `anchor-missing`, Exit 1; Kontrolle `0 Befund(e)` | F-4 („auch wenn das zusammengesetzte Ziel existiert“, für `anchors`) |
| Q14 | Linktext umbrochen, Ziel fehlt / Anker fehlt / Ziel vorhanden | je `0 Befund(e)`, Exit 0 | F-4, Grenze 9 |
| Q15 | `ADR-0094` und der Marker, beide in Inline-Code, in einer Zeile von `spec/architecture.md`; Kontrolle ohne Marker | `0 Befund(e)`; Kontrolle `matrix-forbidden` | Grenze 12 (M33) |

Nachgemessen außerhalb des Klons (Arbeitsbaum `2ac39c31`, sauber):

- Verdikt §2 am Stand `7f139836` mit den drei Befehlen des Verdikts: 5 Dateien
  (0018, 0036, 0038, 0039, 0151), `122`, `27` — wie das Verdikt.
- `make suchlauf-nachmessen PLAN=<Plan>`: Exit 0, `suchlauf-nachmessen: 16 Zeilen stimmen`.
  Die Zuwächse „der Anker und `SPEC-040` je +2, beide im Vertrag“ per `git grep -c`
  an `72e57998` und `2ac39c31`: `harness/sensors/docs-check.md` 6 → 8 für beide
  Muster, sonst keine Datei bewegt; die zwei neuen Stellen sind Grenze 12 und §Sperren.
- §6 vierter Punkt: `git grep -n status-provenance -- spec/` trifft 2 Zeilen (2062,
  2208), die Pipeline mit `grep -cE 'ADR-[0-9]{4}|slice-[0-9]{3}'` druckt `0`; `git
  grep -nE 'ADR-[0-9]{4}|slice-[0-9]{3}' -- spec/` trifft nichts.
- Wortlaut Verdikt §3.1 und §3.2 gegen Spec und Vertrag, whitespace-normalisiert
  als Teilstring verglichen (Zellentext, Randformen-Satz, Grenze-12-Satz,
  Überschrift, Verweis-Ergänzung): alle fünf gefunden.

---

## Status der Findings aus dem ersten Report

| ID | Kat. | Status | Beleg |
|---|---|---|---|
| F-1 | HIGH | **geschlossen** | `SPEC-040` und §Bindung nennen `done/slice-*.md` direkt unter `done/`, Name und Unterverzeichnis ausgenommen; Q1–Q3 |
| F-2 | HIGH | **geschlossen** | Zeile `ids` und Randformen: ohne Link ein Treffer, das Linkziel ungeprüft; der Satz zum ersten Muster ist gestrichen, B2 auf E; Anschluss-Frage ersetzt; Q4, Q5 |
| F-3 | HIGH | **geschlossen** | Verdikt 1 aus Modul 8 (`037a9d32`); Wortlaut aus §3 in Spec, Historie und Grenze 12 übernommen (Vergleich oben); `ADR-0095` unverändert; Q6–Q10 |
| F-4 | HIGH | **geschlossen** | Randformen und Grenze 9 nach M27–M30; Q11–Q14 bestätigen beide Hälften, Q13 auch für den Anker |
| F-5 | HIGH | **erfasst** (nicht behebbar) | §3 trägt den Vorfall mit Weg (`>>` an `README.md`, `.d-check.yml`), Regel (`AGENTS.md` §3.1) und Rücksetzung; Ursache „`cd` schlug fehl“ ist eine Angabe des Implementers (*übernommen*) |
| F-6 | MEDIUM | **geschlossen** | wie F-5; Register-Kandidat `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` existiert und trägt die Umleitung in seinem `state.md` („`>`, `>>`, `tee`, Heredoc“), der Eintrag ist Planner-Arbeit bei der Closure |
| F-7 | MEDIUM | **geschlossen** | Grenze 12: Marker in `spec/`, Regel-Ausnahme der ADRs 0039/0041, Fence, dazu der Status nur am Link |
| F-8 | MEDIUM | **geschlossen** | §3-Zeile auf „erledigt“ mit `ADR-0163`, C8 mit Entscheidung 4 |
| F-9 | LOW | **geschlossen** für §6; fremde Träger **gemeldet** | §6 vierter Punkt mit Marker-Prüfung (nachgemessen oben); fünf Pläne in `open/` als Meldung an den Planner, Frist Closure — offen beim Planner. Zur Zusammenfassung dieses Punkts in §3: N-1 |
| F-10 | LOW | **geschlossen** | Grenze 8, §Ausgabe, §Sperren verweisen auf `SPEC-040` und tragen nur, was das Grün nicht deckt bzw. den Weg frei macht; `scan.ignore` steht weiter in Grenze 5 |
| F-11 | INFO | **geschlossen** | Historie-Zeile vom 2026-10-10 ergänzt, keine neue Zeile |
| F-12 | INFO | keine Handlung | — |

## Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| N-1 | HIGH | Die Fixrunde-Liste sagt, §6 werde „zusätzlich mit `git grep -c status-provenance -- spec/`“ belegt; der Befehl druckt `spec/pflichtenheft.md:2` — er zählt die Zeilen mit dem Marker und sagt nichts darüber, ob eine davon zugleich ein Token `adr`/`slice` trägt. §6 selbst nennt einen anderen Befehl (`git grep -n …` mit `grep -cE …`, druckt `0`), der den Satz trägt. | `AGENTS.md` §3.12 Instanz B; Skill HIGH „Beleg trägt seinen Satz nicht“ | Slice-Plan §3 Fixrunde · „belegt zusätzlich mit `git grep -c status-provenance -- spec/`“ | ja — den genannten Befehl ausführen | Beleg trägt seinen Satz nicht |
| N-2 | LOW | Die Überschrift von Grenze 12 sagt „zwei angenommene Lücken, der Fence und der Status nur am Link“; der Punkt trägt vier Lücken, und sein Schlusssatz nennt als „die zwei Lücken“ die aus `ADR-0163` Entscheidung 2 und 3 (Marker, `exempt-paths`). Die Apposition liest sich als Aufzählung der zwei. Der Wortlaut ist der aus Verdikt §3.2, wörtlich übernommen. | Maintainability; `AGENTS.md` §3.12 Instanz A (Zählwort) | `harness/sensors/docs-check.md` Grenze 12 · „zwei angenommene Lücken, der Fence und der Status nur am“ | nein — Lese-Handlung | Zählwort der Überschrift deckt den Inhalt nicht |
| N-3 | INFO | Verdikt §1 Punkt 1 sagt, die Fitness Function von `ADR-0095` (58 `matrix-inactive`) „bleibt wahr“, und nennt im selben Satz 57 Befunde aus M32. Die 58 sind der Lauf-Beleg am Stand des `ADR-0094`-Diffs (`ADR-0095` §Fitness Function, „real gemessen“), die 57 am heutigen Stand; kein Widerspruch, aber ohne Stand-Angabe zur 58. | `AGENTS.md` §3.12 Instanz A | Verdikt §1 · „Die Fitness Function (58 `matrix-inactive`“ | nein | — |
| N-4 | INFO | Q8 und Q10 messen zwei Punkte, die das Verdikt §4 als *hergeleitet* führt (Token in `welle` kein Treffer; Zielklasse `slice` mit verbotenem Status meldet). Beide wie hergeleitet; für Ziel `spec`, `review`, `observation` und Token in `spec` bleibt *hergeleitet* stehen. | Verdikt §4 | Verdikt §4 · „für `spec`, `slice`, `review` und `observation`“ | ja — Q8, Q10 | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Wortlaut Verdikt §3 gegen Spec und Vertrag | geprüft, ohne Befund: Zelle `matrix-inactive`, Randformen-Satz, Grenze-12-Satz, Überschrift, Verweis — wörtlich (Vergleich oben); Position nach dem Fence-Satz wie verlangt; §3.3 (Plan §6 *entfallen*, B7/C17 auf S, keine neue Historie-Zeile) umgesetzt |
| Verdikt — Mengen-Aussagen (`AGENTS.md` §3.12) | geprüft, ohne Befund außer N-3: §2 (5 ADRs, 122 Zeilen, 27 Dateien) nachgemessen, gleich; „alle sechs Klassen“ als Quelle mit §4 gedeckt (fünf gemessen, `review` *übernommen*); Ziel- und Token-Reichweite gekennzeichnet; Status-Regel *hergeleitet* gekennzeichnet |
| Vertrag Grenze 8 | geprüft, ohne Befund: drei Flächen benannt (Fence und Tilde mit Benutzername, Nicht-Markdown, relative Pfade), Ventil und Wächter stehen; `scan.ignore` in Grenze 5 |
| Vertrag Grenze 9 | geprüft, ohne Befund: Überschrift und Satz decken Q11–Q14; die Lücke ist der umbrochene Linktext |
| Vertrag Grenze 12 | geprüft, ohne Befund außer N-2: Marker-Satz an Q15, Regel-Ausnahme 0039/0041 (R18 des ersten Reports), Fence (M31), Status am Link (Q6–Q10) |
| Linkpflicht (`ids`) | geprüft, ohne Befund: neue Kennungen in Prosa verlinkt (`SPEC-040`, `ADR-0163` im Vertrag, `ADR-0163` und `AGENTS.md` im Plan), sonst in Inline-Code; Basislauf 0 Befunde |
| Matrix-Regel Spec → ADR/Slice | geprüft, ohne Befund: kein `ADR-<NNNN>`-/`slice-<NNN>`-Token und kein Link nach `docs/plan/adr/` oder `docs/reviews/` in `spec/`; Marker nur als Text (Zeilen 2062, 2208), keine davon mit Token |
| Suchlauf (`AGENTS.md` §3.13) | geprüft, ohne Befund: die Fixrunde bewegt keine Eigenschaft des Werkzeugs, sie zieht Sätze an das Werkzeug. Eigener `git grep` über den Baum (ohne `docs/reviews`, `done/`, `.harness/baseline`, den Plan) nach `mehrzeilig`, `Zeilenumbruch`, `umbr`, `erste(s) Muster`, `Basismuster`, `matrix-inactive`, `status.forbidden`, `je geschlossenem Slice`, `jedes geschlossenen Slice`, `Closure-Notiz je`, `done/slice-*`: kein überholter Träger. Treffer, die stehen bleiben dürfen: `ADR-0095` §Kontext („ausgehende `ADR-\d{4}`-Token“, Geschichte nach Verdikt §1); `ADR-0161` sagt „mehrzeiliger Linktext“ (stimmt mit Q14) |
| Plan §6 Frage zu F-3 | geprüft, ohne Befund: Fragetext bleibt, Ausgang *entfallen* mit Grund im selben Punkt |
| Plan §3 Vorfall M1 (F-5/F-6) | geprüft, ohne Befund: Weg, Regel, Rücksetzung, Register-Kandidat, Planner als Eintragender |
| Commit-Messages `7f139836`, `037a9d32`, `2ac39c31` | geprüft, ohne Befund: je eine `ADR-*`-Kennung, keine `SPEC-*` im Betreff |
| `docs/plan/adr/` | geprüft, ohne Befund: kein Diff, `ADR-0095` unverändert |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Zählwort der
Überschrift deckt den Inhalt nicht

## Verdikt

**Merge-blockierend:** ja, wegen N-1 — eine Zeile der Fixrunde-Liste nennt einen
Beleg, der seinen Satz nicht trägt. F-1 bis F-4 sind am gepinnten Werkzeug
nachgemessen geschlossen, F-3 über das Verdikt; F-5 ist erfasst, F-6 bis F-11 sind
geschlossen. Für die Nachprüfung von N-1 genügt die Lesung der geänderten Zeile;
einen weiteren vollen Review-Lauf braucht es nicht. Die DoD-Zeile „Review
durchgeführt“ bleibt offen, weil eine Fixrunde folgt.

**Übergabe:** N-1 an den Implementer. N-2 betrifft den Wortlaut aus Verdikt §3.2;
eine Änderung daran ist eine Frage an den Architect (`AGENTS.md` §3.5,
sinngemäß), keine Fixrunde-Pflicht des Implementers. F-9 (fünf fremde Pläne) bleibt
beim Planner mit Frist Closure. Die **Finding-Klassen** gehen in die
Slice-Closure §7. Dieser Report ist ein **Lauf-Beleg**; er ersetzt keine
Verifikation.
