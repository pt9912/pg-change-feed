# Review-Report: slice-harness-baseline-v6-14-1, Fixrunde — 2026-10-06

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `ba0080d2..65c75a3c` (5 Commits): `f1f6ae70`
([`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
samt ADR-Index), `01399b76` (`AGENTS.md` §3.5), `fca136d4`
(`.claude/agents/verifier.md`, `harness/targets/pin-stale.md`), `891d81e7` und
`65c75a3c` (Slice-Plan, Abschnitt „Fixrunde“). Re-Review zu
[`review-slice-harness-baseline-v6-14-1`](review-slice-harness-baseline-v6-14-1.md).

**Skill:** `.harness/skills/reviewer.md` @ 65c75a3c
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

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

- Slice-Plan `slice-harness-baseline-v6-14-1` (Stand `65c75a3c`), Abschnitt „Fixrunde“ und „Verweise, Pins und Records“
- Erstes Review [`review-slice-harness-baseline-v6-14-1`](review-slice-harness-baseline-v6-14-1.md) (F-1 bis F-7)
- [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (neu, Supersedes [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) teilweise)
- [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) (Entscheidung 1–5)
- [`ADR-0063`](../plan/adr/0063-lh-fa-sch-003-testform-korrektur.md) / [`ADR-0058`](../plan/adr/0058-testansatz-fuenf-luecken.md) (Vorbild Teil-Supersede)
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 (P8), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md)
- `AGENTS.md` §3.1, §3.3, §3.5, §3.6, §3.9, §3.12, §3.13
- `.d-check.yml` Block `vcs:`; `harness/targets/pin-stale.md`; `.claude/agents/{verifier,implementer}.md`
- Register-Eintrag `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`
- keine `LH-*`-Kennung berührt

**Nachgefahren (Exit direkt, Ausgaben im Scratchpad des Laufs):**

- Symlink-Messung je Stand (`git ls-tree -r <Stand> | awk '$1==120000{print $3}'`,
  Blobs per `git cat-file -p`, `grep -c '/v6\.14\.0/'`): `11a5bac5` 7/4,
  `4045dc4f` 7/4, `b6c5b419` 7/0, `HEAD` 7/0 — gleich der Tabelle im Plan.
  Arbeitsbaum: vier Ziele unter `.harness/baseline/v6.14.1/regelwerk/`
  (`modul-01`, `-05`, `-06`, `-08`), drei auf `AGENTS.md`,
  `harness/conventions.md`, `harness/README.md`;
  `find . -path ./.git -prune -o -xtype l -print | wc -l` druckt 0. Der Befehl
  aus `harness/targets/pin-stale.md` Schritt 1 mit `grep -F '/v6.14.0/'`: Exit 1,
  kein Treffer; `git grep -c 'v6\.14\.0' 4045dc4f -- .claude/rules`: Exit 1
  (Grenze bestätigt).
- `make doc-immutable RANGE=11a5bac5..5d8855d9~1`: Exit 0,
  `d-check: 1753 Datei(en) geprüft, 0 Befund(e)`.
  `make doc-immutable RANGE=5d8855d9..HEAD`: Exit 0, dieselbe Zeile.
  Volle Range `11a5bac5..HEAD`: Exit 2, vier `core-drift-vcs` an `MR-001` bis
  `MR-004`. Range der Fixrunde `ba0080d2..HEAD`: Exit 0, 0 Befunde.
- Normalisierter `cmp` (Schleife aus `ADR-0157` Entscheidung 4, mit Ausgabe des
  `cmp`-Exits je Datei): `P=5d8855d9` vier Dateien je `cmp 0`; `P=484d20ec`
  vier Dateien je `cmp 0`.
- Mutation an einem Klon im Scratchpad: Commit mit Pin v6.14.1 → v6.14.2 und
  `Auflösungs-Trigger: permanent` → `nie` in `MR-001`; die Schleife druckt
  `kein reiner Pin: harness/conventions/MR-001-…` (cmp: Byte 1456, Zeile 25);
  der Exit der Schleife selbst ist 0 (F-5).
- Referent der `ADR-0095`-Korrektur: `templates/.d-check.yml` aus `11a5bac5`
  (v6.14.0) gegen v6.14.1, roh `cmp` Exit 0. §Geschichte-Zeile trägt `eadf305`.
- `git log --diff-filter=R -- 'harness/conventions/done/*'`: keine Ausgabe
  (Aussage in `ADR-0157` §Konsequenzen bestätigt).
- `make suchlauf-nachmessen PLAN=<Plan>`: Exit 0, `14 Zeilen stimmen`
  (Zeile 5 `diff 23` nachgemessen).
- `make kommentar-kennungen DIFF=ba0080d2`: Exit 0, kein Kandidat (der Diff
  ändert nur Markdown).
- `make doc-commits RANGE=ba0080d2..HEAD`: Exit 0, 0 Befunde.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Bedingung (b) bindet „Referent gemessen derselbe“ an einen `cmp` des Zielinhalts, Bedingung (a) lässt aber auch Anker und Zeilen-Lokator zu, und die ADR legt nicht fest, welche Einheit dabei verglichen wird. Die einzige ausgeführte Form ist der `cmp` ganzer Dateien (Fitness Function, Entscheidung 4, `ADR-0095`-Beleg). Bei einem Anker-Wechsel in derselben Zieldatei ist dieser `cmp` immer 0. Ein Link in §Entscheidung von `#entscheidung-2` auf `#entscheidung-3` erfüllt dann (a) bis (c), obwohl er den normativen Bezug ändert. Für einen verschobenen Zeilen-Lokator ist derselbe `cmp` dagegen 1, obwohl der Referent gleich ist. Die Normalisierung („roh oder normalisiert“) ist für Entscheidung 1 nicht festgelegt. | [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 1 · `AGENTS.md` §3.12 („Verfasser einer ADR“) | `docs/plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md` · „`cmp` des Zielinhalts an alter und neuer Adresse“ | nein — Urteil am Diff; Probe: `cmp` einer Datei gegen sich selbst ist unabhängig vom Anker 0 | Zulassungsbedingung misst nicht die zugelassene Form |
| F-2 | MEDIUM | Die Fixrunde trägt `ADR-0157` als Entscheidung zu F-2 ein und setzt im selben Slice `AGENTS.md` §3.5 um. Zwei Stellen im selben Plan widersprechen dem und verweisen nicht auf die Fixrunde. §1 schließt die Angleichung der Abschnitte-Liste weiter als „ein anderer Vorgang (Architect-Zug, Folge-ADR)“ aus. §6 erwartet den Ausgang „eingetreten, mit der Kennung des Folge-Slice“, den der Planner „spätestens bei der Closure“ anlegt, und §8 wiederholt „eigener Folge-Slice“. `ADR-0157` §Konsequenzen weist der Planner-Closure dagegen den Ausgang *verkörpert* zu, mit Zielort `AGENTS.md` §3.5. Wer die Closure nach §6 schreibt, legt einen Folge-Slice für eine geschlossene Lücke an oder trägt einen Ausgang ein, der der ADR widerspricht. | `AGENTS.md` §3.13 · `v6.14.1` · `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst | `docs/plan/planning/in-progress/slice-harness-baseline-v6-14-1.md` · „Erwarteter Ausgang: eingetreten, mit der Kennung des Folge-Slice“ (gleichartig „ein anderer Vorgang (Architect-Zug, Folge-ADR)“) | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-3 | LOW | Zwei Träger außerhalb des Diffs beschreiben die `vcs`-Prüfung noch ohne die Ausnahme aus `ADR-0157` Entscheidung 3 und 4. `.claude/agents/implementer.md` nennt `make doc-immutable RANGE=base..head` ohne Teil-Ranges um einen Pin-Commit; ein Implementer, der einen Bump nach diesem Briefing prüft, bekommt vier `core-drift-vcs` (volle Range nachgefahren: Exit 2). Der Kommentar im Block `vcs:` von `.d-check.yml` sagt zu: „Eintraege werden nie ueberschrieben“. Nach `ADR-0157` Entscheidung 3 ist eine Zitat-Korrektur an MR-Einträgen zulässig. Meldung an den Planner, Frist ist die Closure dieses Slice. | `AGENTS.md` §3.13 (Träger in fremder Datei) · `AGENTS.md` §3.7 (Zusage) | `.claude/agents/implementer.md` · „`make doc-immutable RANGE=base..head` bzw. `STAGED=1` — MR-Einträge“; `.d-check.yml` · „Eintraege werden nie ueberschrieben“ | ja — `make doc-immutable RANGE=11a5bac5..HEAD` Exit 2 | Arbeit überholt stehenden Träger |
| F-4 | INFO | Der Absatz „Beleg“ von `AGENTS.md` §3.5 nennt als Klassen ohne §Geschichte nur Records (`done/`, `docs/reviews/**`). Dass MR-Einträge nach `ADR-0157` Entscheidung 3 zur Klasse gehören und der Record-Belegform folgen, steht in §3.5 nicht. Die ADR verlangt das in ihrer Folgepflicht (1) nicht; es steht in `harness/targets/pin-stale.md` und `.claude/agents/verifier.md`. | [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 3 | `AGENTS.md` · „Records (`done/`, `docs/reviews/**`) tragen keine §Geschichte“ | nein | Norm-Reichweite nur im Folgeträger |
| F-5 | INFO | Die Schleife in Entscheidung 4 meldet einen unreinen Pin nur über die gedruckte Zeile „kein reiner Pin“; ihr Exit bleibt 0 (Mutation nachgefahren). Ändert `P` keine MR-Datei, etwa weil der falsche Commit gewählt ist, druckt sie nichts. `verifier.md` fordert „je Exit 0“, die Schleife liefert aber keinen Exit je Datei. Der Plan hat für seinen Beleg eine Variante benutzt, die den Exit je Datei druckt. | [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 4 | `.claude/agents/verifier.md` · „am Pin-Commit selbst je MR-Datei ein `cmp` nach Normalisierung des Tags (je“ | ja — Mutation im Klon, Schleifen-Exit 0 | Prüfschleife färbt den Exit nicht |

## Stand der Findings des ersten Reviews

| Finding | Behauptung „Fixrunde“ | Nachgeprüft |
|---|---|---|
| F-1 (HIGH) | behoben | stimmt: die Symlinks und `b6c5b419` stehen in Liefer-Punkt 2 und §3, die Messung je Stand ist nachgefahren und gleich, und die Grenze von `git grep` ist benannt. **Geschlossen.** |
| F-2 (MEDIUM) | vom Architect entschieden | stimmt: `ADR-0157` ist ein Teil-Supersede nach `AGENTS.md` §3.5 in der Form `ADR-0063` → `ADR-0058`. `ADR-0073` ist unverändert, der Index trägt `(→ ADR-0157)` und „Supers. `ADR-0073`, teilw.“, und `AGENTS.md` §3.5 entspricht Wort für Wort dem Wortlaut aus Entscheidung 2. Die `ADR-0095`-Korrektur erfüllt Entscheidung 1 (roh `cmp` 0, kein Wort der Aussage geändert). **Geschlossen**; neue Befunde F-1 und F-2 oben. |
| F-3 (MEDIUM) | behoben nach `ADR-0157` | stimmt: Teil-Ranges und `cmp` sind nachgefahren und gleich dem Plan, `verifier.md` und `pin-stale.md` folgen Entscheidung 4, und die Zeile nach Entscheidung 5 steht im Plan. **Geschlossen**; Restträger in F-3 oben. |
| F-4 (LOW) | behoben | stimmt: Schritt 1 enthält den Symlink-Befehl und die Prüfung `find -xtype l`; Gegenprobe nachgefahren. **Geschlossen.** |
| F-5 (LOW) | nicht geändert, Befund 3 | begründet: Das Wort „normalisiert“ ändert die Aussage einer `Accepted`-ADR. Der Weg ist eine Folge-ADR oder eine Lesart. Angenommen. |
| F-6 (INFO) | nicht geändert, Befund 4 | angenommen |
| F-7 (INFO) | keine Aktion | angenommen |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `ADR-0157` Form: Status mit Teil-Supersede und benannten Resten, Autor als Architect-Rolle mit eigenem Kontext, `Schärft: —`, drei Optionen je Frage mit „nichts tun“, Re-Evaluierungs-Trigger (a) bis (d), §Geschichte | geprüft, ohne Befund |
| `ADR-0157` §3.12: Getrennt gekennzeichnet sind übernommen (die drei Auftreten, d-check-Spezifikation v0.80.0 gegenüber Image v0.79.0), gemessen (Teil-Ranges, `cmp` an `5d8855d9`/`484d20ec`) und hergeleitet (Mutation an einer Stelle verallgemeinert, `git mv` eines MR-Eintrags). Die Fitness-Function-Zeilen sind an der Quelle und am Klon nachgefahren. | geprüft, ohne Befund über F-1 hinaus |
| `ADR-0157` Begründung: Asymmetrie „der Aussage nach“ in `ADR-0073`, Audit-Lesart, Option G als Trigger (b) | geprüft, ohne Befund |
| `AGENTS.md` §3.5 gegen `ADR-0157` Entscheidung 2 (Wortlaut, Link-Form) | geprüft, ohne Befund |
| `.claude/agents/verifier.md` gegen `ADR-0157` Entscheidung 4 (Definition des Pin-Commits, Teil-Ranges, Verweis auf den Befehl). Der Pin-Commit `5d8855d9` dieses Slice erfüllt die Definition nicht; Entscheidung 5 und die Zeile im Plan tragen ihn. | geprüft, ohne Befund über F-5 hinaus |
| `harness/targets/pin-stale.md` §Bump-Ablauf (Symlink-Schritt in Schritt 1, Absatz zum eigenen MR-Pin-Commit mit `ADR-0073`) | geprüft, ohne Befund |
| `docs/plan/adr/README.md` (Index-Zeilen `ADR-0073`/`ADR-0157`, Satz in §Konventionen) | geprüft, ohne Befund |
| §3.13-Suchlauf nach der alten Abschnitte-Liste (`git grep` nach „Unberührbar“, „§Entscheidung, §Konsequenzen“, „nie eine Zitat-Korrektur“, „Zitat-Korrektur“, „`ADR-0073`“, „doc-immutable“; ganzer Baum ohne `.harness/baseline/**`, `docs/reviews/**`, `done/**`). Treffer gibt es nur in `Accepted`-ADRs (`ADR-0073`, `-0090`, `-0101`, `-0102`, `-0157`), in Register-Records und in diesem Plan. Skills, Kommandos, `harness/sensors/docs-check.md` und `.claude/agents/architect.md` geben die Liste nicht wieder. | geprüft, ohne Befund über F-3 hinaus |
| Commit-Folge der Fixrunde: ADR und Index getrennt von den Norm-Nachzügen und vom Plan; jede Message nennt `ADR-0157` (`make doc-commits` 0 Befunde) | geprüft, ohne Befund |
| Plan-Abschnitte „Verweise, Pins und Records“ und „Fixrunde“ (Zahlen, Befehle, Stände nachgemessen; Suchlauf 14/14) | geprüft, ohne Befund über F-2 hinaus |
| Docker-only (`AGENTS.md` §3.1): Der Diff enthält keine Skripte; der Befehl in `pin-stale.md` benutzt nur `git`, `awk`, `readlink`, `grep` und `find`. Keine Umleitung schreibt in eine Repo-Datei. | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zulassungsbedingung misst nicht die zugelassene Form · Nachzug widerspricht dem Nachbarn im selben Träger · Arbeit überholt stehenden Träger · Norm-Reichweite nur im Folgeträger · Prüfschleife färbt den Exit nicht

## Verdikt

**Merge-blockierend:** nein. Kein HIGH ist offen, F-1 des ersten Reviews ist
geschlossen. F-1 dieses Laufs betrifft eine `Accepted`-ADR. Übergabe-Artefakt
an den Architect ist dieses Finding. Der Weg ist eine Folge-ADR oder eine
Lesart; über Trigger (a) entscheidet er. F-2 und F-3 gehen an den Planner mit
Frist Closure dieses Slice: §6 bekommt seinen Ausgang bei der Closure, F-3
betrifft fremde Dateien (`AGENTS.md` §3.13). F-4 und F-5 sind Hinweise.

**DoD-Haken „Review durchgeführt“:** nachgezogen im selben Commit wie dieser
Report (Skill §DoD-Checkbox-Nachzug ohne Fixrunde). Kein Finding geht als
Reviewer→Implementer-Rückgabe zurück. Die Notiz am Ende des Abschnitts
„Fixrunde“ im Plan („bleibt offen … ein Re-Review folgt“) ist nicht geändert;
der Skill erlaubt dem Reviewer nur die eine Checkbox.

**Übergabe:** F-1 an den Architect, F-2 und F-3 an den Planner. Die
Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg
und ersetzt die Verifikation nicht.

**Lauf-Notiz:** Der erste `make gates`-Lauf nach dem Report-Commit `26bfde80`
endete mit Exit 2: drei `id-unlinked` zu `ADR-0073` in diesem Report. Behoben
in diesem Report im Folge-Commit, danach `make gates` erneut gefahren.
