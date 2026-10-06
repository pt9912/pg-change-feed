# Review-Report: slice-abgeleitete-dokumente-vorlagen-nachzug — 2026-10-06

**Review-Art:** Code — gegen Plan, Architect-Verdikt, Entscheidungen und Hard Rules (Maintainability), nicht gegen die DoD

**Gegenstand:** Diff `3f51d3ed..c5b47c01`, die Commits d53b5985 (Slice-Pfad-Wächter), 9457df79 (ADR-Index, Carveout-Index), 37d86ec1 (Platzhalter in `harness/README.md` und `harness/conventions.md`) und c5b47c01 (Plan-Nachzug und Belege)

**Skill:** `.harness/skills/reviewer.md` @ 7fff2cac
**Modell:** Opus 5.5 (`claude-opus-5-5`) · **Datum:** 2026-10-06

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag +
> Pfad in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-abgeleitete-dokumente-vorlagen-nachzug` (§1 Abgrenzung, §2 Liefer-Punkte 1 bis 3, §3 Ergebnis je Liefer-Punkt, Suchlauf, gemeldeter Träger, §6 Risiken)
- [`architect-verdict-slice-pfad-waechter-namensform`](architect-verdict-slice-pfad-waechter-namensform.md) (Verdikt (a), Constraint 1 bis 3)
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md), [`MR-002`](../../harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md)
- `AGENTS.md` §3.1, §3.3, §3.5, §3.9, §3.12, §3.13
- Baseline `v6.14.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt, `regelwerk/modul-06-roadmap.md` §Wann Arbeit eine Welle braucht, die Vorlagen `templates/harness/README.template.md`, `templates/harness/conventions.template.md`, `templates/docs/plan/adr/README.template.md`, `templates/docs/plan/carveouts/README.template.md`
- berührte `LH-*`-IDs (nur als Quellen der Safety-Punkte): [`LH-FA-CAP-004`](../../spec/lastenheft.md), [`LH-FA-CAP-007`](../../spec/lastenheft.md), [`LH-FA-RET-001`](../../spec/lastenheft.md), [`LH-FA-CFG-006`](../../spec/lastenheft.md), [`LH-QA-SEC-002`](../../spec/lastenheft.md)

**Eigene Proben dieses Laufs (gemessen, 2026-10-06, Stand c5b47c01):**

1. **§Safety, Punkt für Punkt an der Quelle.** `spec/lastenheft.md` §MVP-Schnitt (Zeilen 101–117) führt die Abnahmekriterien „logische Reihenfolge“ (`LH-FA-CAP-004`), „Rollback-Änderungen werden nicht ausgeliefert“ (`LH-FA-CAP-007`) und „Neustart verliert keine dauerhaft erfassten CDC-Daten“ (`LH-FA-RET-001`); `LH-FA-CFG-006` „keine Änderungen am Anwendungscode“; `LH-QA-SEC-002` „getrennt berechtigbar“; `ADR-0047` „eine Verbindungs-DSN je PostgreSQL-Rolle“; die drei Kennungs-Gates stehen in §Sensors; `AGENTS.md` §3.1 trägt Docker-only. Kein Punkt erfunden. Punkt 4 siehe F-5.
2. **§Leseordnung** gegen die Baseline-Regel: vier geordnete Zeiger, der letzte „bei Bedarf“, jeder als Link; die Links lösen auf. Regel erfüllt; Plan-Begründung siehe F-4.
3. **Platzhalter, beide Formen** aus `harness/targets/pin-stale.md` (Bump-Ablauf, Schritt 2) an den vier Dokumenten: Form B 5 / 7 / 0 / 2, Form A 11 / 15 / 2 / 3 Trefferzeilen — gleich den Zahlen des Plans; jeder Treffer gelesen, alle sind Notation (Kennungsformen, `<tag>`, HTML-Kommentare der Vorlage, `<a id>`-Anker). Gliederungs-`diff` gegen die Vorlage (`<Projektname>` durch `PG Change Feed` ersetzt): Exit 0 für alle vier.
4. **`harness/conventions.md`:** `git log --follow --diff-filter=A -- harness/conventions.md` → `40c8c431 2026-09-09` (Status `A`, legt die Datei an) — das MR-000-Datum trägt. URL: `tools/harness/pin-stale-baseline.sh` fragt `repos/pt9912/ai-harness-course/releases/latest` — trägt. `harness/conventions/done/` existiert nicht — `— | —` trägt. Glossar: die Vorlage führt den Abschnitt als optional, `spec/lastenheft.md` §6 Glossar existiert — trägt. Zusatzklassen: Bindung-Spalte der Gate- und der Werkzeug-Tabelle ausgezählt (66 Zeilen) — siehe F-1.
5. **Gestrichene Musterzeilen der Sensors-Tabelle:** `grep -nE '^(fullbuild|ci|closure)[: ]' Makefile harness/mk/*.mk *.mk` Exit 1 (kein Treffer). Die Baseline verlangt keinen Closure-Eintrag: die Vorlage sagt „Nur eintragen, wenn das Ziel im Makefile existiert“ und „Weglassen ist nur für Targets richtig, die niemand braucht“; `regelwerk/modul-13-quality-gates.md` §Hard Rule verbietet behauptete, nicht existierende Befehle. Die Streichung verletzt keine Pflicht.
6. **`d-check:ignore` wirkungslos — nachgefahren** an einem Scratchpad-Klon auf `3f51d3ed` (Kommentar per `sed … > Kopie` entfernt, Kopie zurückkopiert), d-check `ghcr.io/pt9912/d-check@sha256:b4b8756b…`: Modulliste der `.d-check.yml` Exit 0, `1744 Datei(en) geprüft, 0 Befund(e)`; `--enable codepaths` mit und ohne Kommentar je Exit 1, `132 Befund(e)`, keiner in `docs/plan/planning/README.md`. Die Zeile entspricht danach wörtlich `templates/docs/plan/planning/README.template.md` Zeile 47. Bestätigt.
7. **Liefer-Punkt 3 — Muster und Mutation nachgefahren.** Beide `forbid-pattern` lauten wörtlich wie Constraint 1 des Verdikts; Grenze (5) und der Satz in `harness/sensors/docs-check.md` Grenze 6 sind gestrichen, `· seit slice-075` und (1) bis (4) stehen. Scratchpad-Klon auf c5b47c01: Bestand Exit 0, 0 Befunde; Mutation an **anderen** Stellen und Formen als der Implementer — ein Link mit Ziel `../plan/planning/open/slice-foo-bar.md` in `docs/reviews/review-slice-harness-fmt-check.md` und einer mit Ziel `../../../next/slice-x.md` in `observations/BEO-PGC/adr-aussage-breiter-als-ihre-messung/observation.md` — Exit 1, `4 Befund(e)`, davon je Datei ein `section-forbidden` (dazu zwei `target-missing`). Die Lifecycle-Formen `open/` und `next/` sind damit zusätzlich zu `in-progress/` erprobt.
8. `make suchlauf-nachmessen PLAN=<Plan>` Exit 0, `suchlauf-nachmessen: 10 Zeilen stimmen`. Zusätzlich eigene Suche nach der Kennungsform in Trägern der Regel (`git grep -n -E 'slice-NNN|slice-<NNN>'` ohne `docs/reviews`, `done/`, `.harness/baseline`) — siehe F-7.
9. `make kommentar-kennungen DIFF=3f51d3ed` Exit 0, kein Kandidat (der geänderte YAML-Kommentar in `.d-check.yml` ist im Gegenstand).
10. Commit-Form: vier Commits, je `ADR-0051` in der Message; kein `git mv` im Diff (§3.3 nicht berührt).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Die Zusatzklassen sind laut Plan „die Bindung-Spalte beider Tabellen …, ausgezählt“; die Zählung trägt das nicht: `LH-*` steht in **2** Bindung-Zellen (`make test-integration`, `.github/workflows/e2e.yml`), nicht in 3, und zwei Bindungsformen der Spalte fehlen in der Deklaration — der Hard-Rule-Link auf `AGENTS.md` §3.7 (Zeile `make kommentar-kennungen`) und die nackte Vorgangskennung ohne `seit` (`slice-sdk-kotlin-reale2e`, `slice-sdk-csharp-reale2e` in den Zeilen der SDK-Integrationstests). Failure-Szenario: die Tabelle gilt nach dem Plan als vollständig; die Baseline nennt eine undeklarierte Bindung eine stille Setzung, und ein Reviewer, der die Spalte gegen die Deklaration hält, findet zwei Formen ohne Klasse. | Skill HIGH „Beleg trägt seinen Satz nicht“ (genannte Zählung nachgefahren, anderes Ergebnis); `v6.14.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt (weitere Klassen im Konventionsdokument deklariert) | `harness/conventions.md` §Zusatzklassen · `die Adaption, die die Regel des Werkzeugs trägt` (letzte Zeile); Plan §3 · `(3 Zeilen)` | nein — kein Sensor hält die Bindung-Spalte gegen die Deklaration; Probe: Bindung-Zellen extrahieren, bekannte Klassen abziehen, Rest lesen | Auszählung trägt die abgeleitete Klassenliste nicht |
| F-2 | MEDIUM | Der neue Abschnitt `## Konventionen` sagt, ADRs seien nach `Accepted` immutable und Schärfungen entstünden als neue ADR, ohne Ausnahme; der Kopf-Absatz desselben Dokuments nimmt die Zitat-Korrektur nach `ADR-0073` ausdrücklich aus. Keine der beiden Stellen verweist auf die andere. Failure-Szenario: wer den Abschnitt liest, den die Vorlage als Regel-Ort setzt, hält eine zulässige Zitat-Korrektur an einer `Accepted`-ADR für verboten bzw. sieht zwei verschieden weite Regeln im Index. | Skill MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“; `AGENTS.md` §3.5 | `docs/plan/adr/README.md` §Konventionen · `ADRs sind nach` (erster Punkt) gegen den Kopf-Absatz · `Zitat-Korrektur` | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-3 | LOW | „Ohne Welle trägt die Slice-Closure diesen Trigger-Audit“ liest sich als Eigenschaft des Slice; die Baseline stellt die Achse voran: wellenlos ist eine Eigenschaft des Repos, und in einem Repo mit Wellen-Betrieb prüft die nächste Welle-Closure auch die Slices ohne Wellen-Zugehörigkeit. Dieses Repo führt Wellen (Roadmap §Abgeschlossene Wellen). | `v6.14.0` · `regelwerk/modul-06-roadmap.md` §Wann Arbeit eine Welle braucht („Achse zuerst“) | `docs/plan/carveouts/README.md` · `Ohne Welle trägt die Slice-Closure diesen Trigger-Audit` | nein — Lese-Handlung | Baseline-Regel enger zitiert als ihre Achse |
| F-4 | LOW | Die Plan-Quelle der Leseordnung heißt „Rang-Reihenfolge von `AGENTS.md` §2“, die gewählte Folge ist Rang 7, `AGENTS.md` §3, Rang 1 — die Quelle trägt die Auswahl der Zeiger, nicht ihre Reihenfolge. Die Leseordnung selbst erfüllt die Baseline-Regel (Probe 2). | `AGENTS.md` §3.12 Instanz B; [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | Plan §3 · `Rang-Reihenfolge von` | nein — Lese-Handlung | Quellenangabe trägt nur einen Teil der Aussage |
| F-5 | INFO | Safety-Punkt 4 sagt, die Nutzerdokumente unter `docs/user/` trügen keine interne Kennung; unter `docs/user/` liegen vier Runner-Erzeugnisse mit 114 Zeilen `LH-` (`grep -c 'LH-' docs/user/*.md`: 80 / 27 / 5 / 2), die das Gate bewusst ausnimmt. Die Formulierung übernimmt das Vokabular der Gate-Zeile und von `harness/sensors/handbuch-public-doc-check.md`; wörtlich gelesen ist sie weiter als das Gate. | Maintainability | `harness/README.md` §Safety · `Nutzerseitige Texte tragen keine interne Kennung` | ja — `make handbuch-public-doc-check` nennt „3 Nutzerdokumente“ | — |
| F-6 | INFO | Die Spalten-Begründung „eine Index-Spalte wäre eine zweite Quelle für denselben Bezug“ trüge ebenso gegen `Titel` und `Status`; die Vorlage erklärt den Index ausdrücklich als derivativ, und ihr `Bezug` ist eine Vertrags-Kennung (`<LH-FA-NN>`), das `**Schärft:**`-Feld eine Spec-Stelle (`SPEC-*`, `ARC-*`, `LH-FA-*.<Buchstabe>`). Tragend sind die Hälfte zur `structure`-Regel (`.d-check.yml`: `ID` 8–8 Zeichen) und die Messung 155 von 155 Dateien mit `**Schärft:**` (nachgemessen). | Maintainability | `docs/plan/adr/README.md` · `eine Index-Spalte wäre eine zweite Quelle für denselben Bezug` | nein | — |
| F-7 | INFO | Der gemeldete Träger ist richtig: `harness/sensors/docs-check.md` §Bindung nennt `BEO-PGC/slice-pfad-als-link-in-berichten` mit „3×“, `evidence/` des Eintrags trägt 5 Dateien. Die Datei liegt im Diff dieses Slice; „fremd“ ist hier der Gegenstand, nicht die Datei. Dazu ein weiterer Träger der Kennungsform, den M1 nicht trifft: `state.md` desselben Eintrags gibt den Reparatur-Pfad des `hint` mit der Nummern-Form `slice-NNN` wieder, der `hint` sagt `slice-<Kennung>`. Meldung an den Planner, Frist wie im Plan: die Closure dieses Slice. | `AGENTS.md` §3.13 | `docs/plan/planning/observations/BEO-PGC/slice-pfad-als-link-in-berichten/state.md` · `die Kennung zitieren` | nein — Träger außerhalb der Diff-Zeilen | — |
| F-8 | INFO | MR-000 führt im ID-Schema `BEO-<NNN>`; das Register und der Skill führen `BEO-<KUERZEL>/<slug>`, und `regelwerk/modul-06-roadmap.md` sagt „Eine fortlaufende Nummer gibt es nicht mehr“. Die Vorlage `templates/harness/conventions.template.md` Zeile 107 trägt dieselbe Form — die Abweichung kommt aus der Baseline, nicht aus diesem Diff. Hinweis an den Planner (Freshness-Audit, `AGENTS.md` §1). | `v6.14.0` · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register | `harness/conventions.md` MR-000 · `BEO-<NNN>` | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.d-check.yml` (Liefer-Punkt 3) | geprüft, ohne Befund — Muster wörtlich nach Verdikt-Constraint 1, Grenze (5) gestrichen, (1)–(4) und Anker stehen; Mutation nachgefahren (Probe 7) |
| `harness/sensors/docs-check.md` §Grenze 6 | geprüft, ohne Befund — Satz zur Nummern-Form gestrichen, Rest unverändert |
| `harness/README.md` §Sensors und Werkzeuge (gestrichene Musterzeilen) | geprüft, ohne Befund — keine Baseline-Pflicht verletzt (Probe 5) |
| `harness/README.md` §Safety, §Leseordnung | geprüft; Punkte 1, 2, 3, 5 tragen an ihrer Quelle; Punkt 4 F-5; Leseordnung regelkonform |
| `harness/conventions.md` §Adoptierte Konventions-Quellen, MR-000-Datum, §Aufgelöste Adaptionen, §Glossar | geprüft, ohne Befund (Probe 4) |
| `docs/plan/carveouts/README.md` | geprüft; Gliederung, Derivativ-Klausel und Konventionen der Vorlage vorhanden, Tabelle durch „(keine)“ samt Spalten-Satz ersetzt; F-3. `.gitkeep` daneben ist nicht mehr nötig (das Verzeichnis trägt eine getrackte Datei), das Repo hält `.gitkeep` aber auch neben Inhalt in `docs/plan/adr/`, `docs/reviews/`, `harness/conventions/` — kein Konventions-Anker, kein Befund |
| `docs/plan/planning/README.md` (`d-check:ignore`) | geprüft, ohne Befund — wirkungslos nachgemessen (Probe 6) |
| Plan §3 (Zahlen der Platzhalter-Prüfung, Suchlauf) | geprüft; Zahlen stimmen (Probe 3, 8); F-1, F-4 |
| Hard Rules §3.1 / §3.3 / §3.9 im Diff | geprüft, ohne Befund — keine Skript- oder Code-Änderung, keine Move-Commits |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Auszählung trägt die abgeleitete Klassenliste nicht · Nachzug widerspricht dem Nachbarn im selben Träger · Baseline-Regel enger zitiert als ihre Achse · Quellenangabe trägt nur einen Teil der Aussage

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) und F-2 (MEDIUM) gehen in eine Fixrunde
an den Implementer; die DoD-Zeile „Review durchgeführt“ bleibt deshalb offen
und wird bei Schritt 21 des Implementer-Ablaufs nachgezogen. F-3 und F-4 nach
Annahme oder Begründung des Implementers; F-7 und F-8 sind Meldungen an den
Planner.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report
ist ein **Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11).
