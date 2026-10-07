# Review-Report: slice-harness-baseline-v6-16-0 — 2026-10-07

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `1009269e..1b7a486a` (11 Commits, darunter `d9c8a5ff` mit
[`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
vom Architect, der Form-Commit `3f88e0c2` und das Entfernen von v6.14.1 in
`9a7da482`)

**Skill:** `.harness/skills/reviewer.md` @ 1b7a486a
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

- Slice-Plan `slice-harness-baseline-v6-16-0` (Stand `1b7a486a`)
- [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
  (neu, Entscheidung 1 bis 7, Folgepflicht 1 bis 8) und die Vorgabe des
  Architect an den Orchestrator (vom Auftraggeber wörtlich übernommen, nicht
  als Artefakt im Repo)
- [`ADR-0156`](../plan/adr/0156-versions-gate-nimmt-done-records-aus.md),
  [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md),
  [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md),
  [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
  (teilweise abgelöst durch `ADR-0161`),
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
  [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P8)
- `harness/targets/pin-stale.md` §Bump-Ablauf
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.5, §3.6, §3.7, §3.9, §3.12, §3.13, §3.15)
- Vorlagen `v6.16.0` · `templates/AGENTS.template.md`,
  `templates/harness/README.template.md`,
  `templates/harness/conventions/MR-NNN-titel.template.md`
- keine `LH-*`-Kennung berührt

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.16.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Klammer „für einen MR-Eintrag `ADR-0159` Entscheidung 2“ ist gestrichen statt ersetzt; der Satz davor verlangt für eine Zitat-Korrektur an einem MR-Eintrag weiter die gedruckte Zeile von `make zitat-vergleich`, während `harness/targets/zitat-vergleich.md` Zeile 10 die Beleg-Rolle jetzt auf die `Accepted` ADR beschränkt und die Vorgabe dort „an einem MR-Eintrag“ lautete. Zwei Träger sagen Verschiedenes über den Beleg eines Form-Commits `F`; der Verweis auf `formnorm`-`cmp` (Entscheidung 4) fehlt in §3.5, und der Beleg von `3f88e0c2` im Plan trägt nur `formnorm`, keine `zitat-vergleich`-Zeile vor/nach `F`. Szenario: der nächste Bump-Implementer liest §3.5 und belegt den Form-Commit mit der falschen Messung, oder der Verifier meldet einen fehlenden Beleg, den `zitat-vergleich.md` nicht verlangt. | [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) Folgepflicht 6 („ersetzt“) und 7, Entscheidung 4 | `AGENTS.md` · „Entscheidung 4, 5); für Records bleibt die Commit-Kennung der Beleg.“; `harness/targets/zitat-vergleich.md` · „die eine Zitat-Korrektur an einer `Accepted` ADR belegt“ | nein — Lese-Handlung, kein Gate | Implementierung weicht von ADR-Wortlaut ab |
| F-2 | MEDIUM | Der Plan zieht DoD Liefer-Punkt 2 und §3 auf `ADR-0161` nach, lässt im selben Dokument aber die überholten Aussagen stehen: §1 **Ziel** „die MR-Pins sind in einem eigenen Pin-Commit umgestellt“, §1 **Bezug** `ADR-0157` „(MR-Pins im eigenen Commit)“ und `ADR-0159` „(Referent-Messung …)“ ohne `ADR-0161`, §1 Ausschluss Records „Referent je Verweis mit `make zitat-vergleich` gemessen (Liefer-Punkt 2)“ (für `ce045921` nicht gemessen und nicht mehr verlangt), §3 Symlinks „Die Prüfung auf hängende Symlinks nach dem Löschen steht noch aus“ gegen den Beleg „druckt 0“ in der DoD. Szenario: Verifier oder Closure-Reviewer prüft gegen das Ziel und findet keinen Pin-Commit. | Maintainability; `AGENTS.md` §3.13 | `docs/plan/planning/in-progress/slice-harness-baseline-v6-16-0.md` · „MR-Pins sind in einem eigenen Pin-Commit umgestellt“; · „Die Prüfung auf hängende Symlinks nach dem Löschen steht noch aus.“ | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-3 | LOW | Die Vorgabe lautete „stattdessen liest der Verifier den Durchgangs-Beleg im Bump-Plan“; der Agent-Text macht daraus ein Nachfahren der Referent-Messung („du fährst sie nach“, „der Verifier fährt beide nach“). Das ist eine Verifier-Pflicht, die `ADR-0161` Entscheidung 5 (Prüfauftrag des Durchgangs, nicht Bedingung eines Commits) nicht nennt; nach dem Löschen des alten Tags ist sie nur an einem historischen Stand fahrbar. | [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) Entscheidung 5, Folgepflicht 5 | `.claude/agents/verifier.md` · „Adaptions-Durchgangs im Plan des Bumps (Entscheidung 5), du fährst sie nach“; `.claude/agents/implementer.md` · „fährt beide nach“ | nein | Implementierung weicht von ADR-Wortlaut ab |
| F-4 | LOW | Im Block `versions:` steht der Kommentar zu `docs/plan/planning/done/**` jetzt über `current-from:` statt über `exempt-paths:`, den er erklärt; im Block `vcs:` steht „Der Lauf umgeht diesen Commit …“ zwischen `paths:` und `immutable-when:`. Die Trennung hält je Block eine Kennung, rückt die Kommentare aber von der Stelle ab, die sie beschreiben. | `AGENTS.md` §3.7 | `.d-check.yml` · „current-from: harness/conventions.md#baseline“ (Zeile 123, direkt unter dem Kommentar, der mit „(ADR-0156).“ endet) | nein | Kommentar steht nicht an seiner Stelle |
| F-5 | INFO | Der `versions`-Kommentar trägt die Vorgabe sinngemäß, ohne den Grund „ab Accepted bzw. Datum immutabel; ein version-stale entsteht durch einen Bump außerhalb der Datei“ und mit dem Zusatz „am Pfad, nicht am Status“ aus Entscheidung 2. Keine Bedeutungsänderung. | `ADR-0161` Folgepflicht 1 | `.d-check.yml` · „ADR-Dateien und MR-Eintraege nennen den Stand ihrer Abfassung“ | nein | — |
| F-6 | INFO | `MR-005` trägt den ADR-Titel statt „… frieren ein“ und einen breiteren Geltungsbereich als die Vorgabe (zusätzlich Block `vcs:` und `AGENTS.md` §3.5). Feld „Ersetzt-Baseline-Regel“, Adaption, Begründung und Auflösungs-Trigger entsprechen Entscheidung 1 bis 6; der Link ist einzeilig und wird von `links`/`anchors` gelesen. | `ADR-0161` Entscheidung 6 | `harness/conventions/MR-005-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md` · „Block `vcs:`, die Baseline-Pins“ | nein | — |
| F-7 | INFO | Die Ausgänge R2–R10 und T1–T10 nennen „Folge-Slice Neu 1“ bzw. „Neu 2“ ohne Kennung; in `open/` und `next/` liegt keine Datei mit dem Gegenstand (`git grep` nach `harness/mk/<werkzeug>` und „Festlegungen der Harness-Werkzeuge“ dort ohne Treffer). Die Probe „nimmt die Adresse die Sendung an“ ist erst nach Anlage möglich; Frist laut DoD ist die Closure (Planner). | Modul 5 `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice, Klasse 1 | Plan · „Für Neu 1 und Neu 2 ist „Folge-Slice““ | nein | — |
| F-8 | INFO | Die übernommene Klausel „TARGET-ZELLE = NACKTER NAME“ entspricht der Vorlage `v6.16.0` wörtlich; ihr Beispiel `make verify-slice` ist kein Ziel dieses Makefiles, und der Sensor, für den die Zeile „unsichtbar“ würde (`targets`), ist hier nicht aktiv (Neu 1). Im Kommentarblock unschädlich. | Vorlage `v6.16.0` · `templates/harness/README.template.md` | `harness/README.md` · „`make verify-slice SLICE=<id>` in der Code-Span“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Vendoring `.harness/baseline/v6.16.0/` | geprüft, ohne Befund — Asset per `gh release download v6.16.0` in `<Scratchpad>/review-6160/`, `sha256sum` `feb4d744…` gleich dem Release-`SHA256SUMS` und dem Plan; im Container entpackt (`unzip`, `--network none`), `diff -r` gegen `regelwerk/` und `templates/` Exit 0, 54 Dateien, selbst gebildete Summen gleich dem vendored `SHA256SUMS`; `make baseline-verify` Exit 0, gedruckt `baseline-verify: v6.16.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)`; nur ein Tag-Verzeichnis |
| Delta v6.14.1 → v6.16.0 (Plan „Bump-Ablauf — Belege“) | geprüft, ohne Befund — nachgemessen aus `git archive bee507d7`, normalisiert: `templates/` 103 Zeilen, 8 Dateien, 70 `<`/`>`; `regelwerk/` 131, 6, 104 — gleich der Tabelle; jede Zeile R1–R10, T1–T10 trägt einen Ausgang |
| Adaptions-Durchgang `MR-001`–`MR-004` | geprüft, ohne Befund — `make zitat-vergleich` am Stand `bee507d7` nachgefahren: `MR-001` Exit 1 (`cmp 1`), `MR-002`–`MR-004` je `cmp 0`, gleich dem Plan |
| Form-Commit `3f88e0c2` | geprüft, ohne Befund — ändert nur die vier MR-Dateien, Message nennt `ADR-0073` und `ADR-0161`, liegt vor `9a7da482`; `formnorm`-`cmp` nach Entscheidung 4: vier Zeilen, je `cmp 0` |
| `teilrange` um `F` (`ADR-0160` Entscheidung 1) | geprüft, ohne Befund — `1009269e..3f88e0c2~1` „enthält 6 Commit(s), Lauf“, `0 Befund(e)`, Exit 0; `3f88e0c2..1b7a486a` „enthält 4 Commit(s), Lauf“, `0 Befund(e)`, Exit 0; Gegenprobe volle Range `make doc-immutable RANGE=1009269e..1b7a486a`: `4 Befund(e)` `core-drift-vcs` je MR-Datei, make-Exit 2 |
| Record-Korrektur `ce045921` | geprüft, ohne Befund — eigener Commit, nur Form (Link → Inline-Code, Pfad und Anker gleich), Message nennt `ADR-0073`, `ADR-0156` |
| Symlinks `.claude/rules/` | geprüft, ohne Befund — vier Ziele auf v6.16.0, `find . -path ./.git -prune -o -xtype l -print \| wc -l` druckt 0 |
| `harness/targets/pin-stale.md` | geprüft, ohne Befund — Absatz durch Entscheidung 1, 4, 5 ersetzt, Schritt 1 trägt die Referent-Messung als Prüfauftrag; Abschnitt „Bump eines Gate-Werkzeugs“ unverändert |
| `harness/targets/zitat-vergleich.md` Zeile 60 und 239 | geprüft, ohne Befund — `<stand>` mit beiden Bäumen, Verweis auf `ADR-0161` Entscheidung 4 (Zeile 10: F-1) |
| `AGENTS.md` §3.6, §5 und `harness/README.md` Kommentar-Block | geprüft — Wortlaut gleich der Vorlage `v6.16.0`; §5-Tabelle mit Datei-Spalte „—“ wie im Vorlagen-Kommentar; `<PREFIX>-RB-<NN>` steht im ID-Schema von MR-000 (F-8 INFO) |
| `docs/plan/adr/README.md` | geprüft, ohne Befund — Index-Zeile `ADR-0161`, Rückverweise an 0157/0159/0160, Konventions-Zeile |
| `ADR-0161` (Architect) | geprüft, ohne Befund — Stichproben: 2 `Proposed`-ADRs (`0021`, `0022`), aus `docs/plan/adr/[0-9]*.md` kein Markdown-Link in die Baseline außer dem Formbeispiel im Fence; Pins nach `pin-pattern` nur `ADR-0095`, `ADR-0160` |
| Suchlauf | geprüft, ohne Befund — `make suchlauf-nachmessen PLAN=…` Exit 0, „18 Zeilen stimmen“; Verteilung der `diff`-Zeile 1 nachgezählt (ADRs 17, `ADR-0161` 10, MR 4, Register 3, Fixture 18, `zitat-vergleich.md` 1) |
| Kommentar-Herkunft | geprüft, ohne Befund — `make kommentar-kennungen DIFF=1009269e` Exit 0, kein Kandidat |
| Commit-Folge und Push-Stand | geprüft, ohne Befund — jede Message trägt `ADR-*`; `origin/main` = `1009269e`, der per `--amend` ersetzte Commit `3dd74886` war nicht gepusht |
| Docker-only, §3.15 | geprüft, ohne Befund — Kopie aus dem Klon per `cp -r`, Normalisierung auf Scratchpad-Kopien; keine Verweigerung im Bericht |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Implementierung weicht von ADR-Wortlaut ab · Nachzug widerspricht dem Nachbarn im selben Träger · Kommentar steht nicht an seiner Stelle

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM) brauchen eine Fixrunde am
Implementer; F-3 und F-4 können mitgehen. Die DoD-Zeile „Review durchgeführt“
bleibt offen (Nachzug bei Schritt 21 des Implementer-Ablaufs).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Kein HIGH mit
Rollen-Widerspruch, kein Konflikt-Pfad über den Architect nötig — F-1 und F-3
messen am Wortlaut von `ADR-0161` und der Vorgabe des Architect, nicht an einer
neuen Entscheidung. Dieser Report ersetzt keine Verifikation (Modul 11).
