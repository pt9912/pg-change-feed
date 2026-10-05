# Review-Report: slice-harness-readme-zellen-kuerzen — 2026-10-05

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (Maintainability), nicht gegen die DoD

**Gegenstand:** Diff `b61412ac..76c3b9e0`, die Commits c967b51b (Umzug nach `harness/targets/` und `harness/sensors/`), 3f900161 (Index-Zeilen gekürzt), 10ad71b0 (`structure`-Regel in `.d-check.yml`) und 76c3b9e0 (Plan: Startmessung, Belege, Befunde)

**Skill:** `.harness/skills/reviewer.md` @ 675246dd
**Modell:** Opus 5.5 (`claude-opus-5-5`) · **Datum:** 2026-10-05

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag +
> Pfad in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-readme-zellen-kuerzen` (§1 Abgrenzung, §3 Umsetzung mit Plan-Nachzug, Befund-Liste und Träger-Meldungen, §6 Risiken)
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `AGENTS.md` §3.1, §3.3, §3.5, §3.7, §3.9, §3.12, §3.13, §3.14, §4
- `v6.14.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice (Out-of-Scope)
- keine `LH-*`-ID berührt

**Eigene Proben dieses Laufs (gemessen, 2026-10-05, Stand 76c3b9e0):**

1. **Gegenprobe Wortlaut** — das Skript aus Plan §3 unverändert aus einer Datei im Scratchpad gefahren: 63 Zeilen `OK`, letzte Zeile `geaenderte Zeilen: 63, Fehler: 0`, Exit 0; die gedruckten Werte stimmen mit dem Plan (`OK 147 targets/test-integration.md alt=26149 form=26314`, `OK 124 sensors/gates.md alt=230 form=230`).
2. **Rückrichtung, andere Methode** — Zelle `make test-integration` am Parent gegen die Absätze 7–53 von `harness/targets/test-integration.md`, mit Leerzeichen verbunden und die Link-Präfixe zurückgeschrieben: `cmp` Exit 0, je 26.464 Byte, 2.694 Wörter auf beiden Seiten. Zusätzlich für alle 15 neuen Dateien und die hinzugefügten Zeilen der 15 ergänzten Dateien: jeder Absatz außer Überschrift, Einleitungssatz und `**Bindung:**` steht als Teilstring in einer alten Zelle — kein Treffer außer dem Liefer-Punkt-3-Text in `harness/sensors/docs-check.md` (nichts erfunden).
3. **Mutationsprobe der Regel** — Klon des Stands 76c3b9e0 im Scratchpad, je `make docs-check` im Klon: ohne Mutation `0 Befund(e)`, Exit 0. Mutiert: `make test-sdk-altserver` (letzte echte Werkzeug-Zeile) `Tut was` 121 Zeichen → `section-cell-oversized … hat 121 Zeichen, erlaubt sind 120`, Exit 2; Platzhalter `make <vorschau>` 121 → rot, Exit 2; `make gates` `Vertrag` 221 → `… hat 221 Zeichen, erlaubt sind 220`, Exit 2; Platzhalter `<make-target>` `Vertrag` 221 → rot, Exit 2; 120 × `ä` (240 Byte) → grün, Exit 0; 121 × `ä` → rot (gezählt werden Zeichen, nicht Byte); ein Link `[x](b…)` mit 130 Zeichen Rohtext → `… hat 130 Zeichen`, rot (gezählt wird der Markdown-Rohtext). Beide Tabellen sind erfasst, auch an ihren Randzeilen.
4. `make suchlauf-nachmessen PLAN=<Plan>` — `suchlauf-nachmessen: 12 Zeilen stimmen`, Exit 0.
5. Anlass- und Ergebnis-Zahlen des Plans nachgemessen: Zeile 147 am Parent 26.286 Zeichen; README 34.679 Byte, 265 Zeilen; längste `Vertrag` 162 (`make commit-traceability`), längste `Tut was` 113 (`make pin-stale-all`); Tabellenzeilen 70 an Parent und Endstand; Werkzeuge-Tabelle am Parent 58 Zeilen, 22 mit Datei-Link — alle wie im Plan.
6. Träger-Suche: `git grep` nach `harness/README.md` in `harness/`, `tools/`, `examples/`, `.github/`, `.claude/`, `AGENTS.md` und nach Positionsverweisen (`Zeile oben`, `Zeile darunter`, `siehe … Zeile`) in den hinzugefügten Zeilen; jede Trefferzeile gelesen.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die 15 ergänzten Dateien tragen am Ende `## Fassung im Gate-Index` neben ihrem eigenen Vertrag, ohne dass einer der beiden Teile den anderen als maßgeblich nennt oder auf ihn verweist; in `harness/sensors/handbuch-public-doc-check.md` steht damit in derselben Datei „Die vier Erzeugnisse sind ausgenommen“ (Grenze 3) neben „die fünf ausgenommenen Dateien“ (Abschnitt `make test-handbuch-public-doc-check`). Failure-Szenario: der nächste Slice, der den Sensor ändert, zieht den Vertrag nach und lässt die Fassung stehen — kein Leser der README sieht sie mehr, und die neue Regel begrenzt nur die README-Zelle, also wächst die Abweichung unbemerkt in einer Datei. Die Befund-Liste (Punkte 7 und 9) nennt den Zustand ehrlich; im Träger selbst steht er nicht. | `AGENTS.md` §3.13 (Abgrenzung Nachzug im selben Träger); Skill MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“ | `harness/sensors/handbuch-public-doc-check.md` · `die fünf ausgenommenen Dateien` | nein — kein Gate liest Widersprüche zwischen Abschnitten | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-2 | LOW | Die Träger-Meldungen in §3 nennen `tools/coverage-gate.sh`, drei Test-Kommentare, `AGENTS.md` §4 und den Bump-Slice, aber zwei gelesene Treffer des Suchlaufs (Zeilen 5/6) fehlen: `AGENTS.md` §3.14 verweist für die zentrale Wache auf die README-Zeile (`tools/schema/rolloutguard`, siehe … `make schema-rollout`-Zeile), die die Wache nicht mehr nennt (sie steht jetzt in `harness/targets/schema-rollout.md`); `harness/mk/coverage.mk` nennt `harness/README.md` §Sensors als Kalibrierungs-Bindung und sagt, die übrigen Träger nennen die Rampe (70 % → 80 %) — die README-Zelle nennt sie nicht mehr. | `AGENTS.md` §3.13 | `AGENTS.md` · `siehe \`harness/README.md\` §Sensors,` und `harness/mk/coverage.mk` · `Kalibrierungs-Bindung (harness/README.md §Sensors` | nein — Lese-Handlung; `make suchlauf-nachmessen` prüft Zahlen, nicht die Auswertung | Träger-Meldung unvollständig |
| F-3 | LOW | Plan-Nachzug 1 (ganze Zelle statt „soweit fehlt“) ist begründet und trägt — der Grund, dass „schon da“ am Diff nicht nachprüfbar ist, ersetzt ein Urteil durch einen byte-genauen Vergleich, und Neuschnitt fand nicht statt (nichts ersetzt, Proben 1 und 2). Aber §1 („die Mega-Zelle wird dort nur ergänzt, soweit ihr Inhalt dort noch fehlt“) steht unverändert und ohne Verweis auf den Nachzug; nach der eigenen Regel des §1 ist das eine Planänderung, die der Abgrenzungs-Abschnitt nicht zeigt. | `v6.14.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice (Out-of-Scope: „hat den Plan geändert, nicht nur ergänzt“) | `docs/plan/planning/in-progress/slice-harness-readme-zellen-kuerzen.md` · `nur ergänzt, soweit ihr Inhalt dort noch fehlt` | nein | Abgrenzung nicht mitgezogen bei Planänderung |
| F-4 | INFO | Die Bindung-Zelle `make doc-tracked` verweist mit „§Bindung“ in `harness/sensors/docs-check.md`; die umgezogene Zelle steht dort unter `## Fassung im Gate-Index`, und §Bindung nennt `doc-tracked` nicht (`git grep doc-tracked` in der Datei: Zeilen 55 und 224). Der Abschnitts-Zeiger stand schon am Parent so; der Umzug hätte ihn auflösen können. | Skill HIGH „Beleg trägt seinen Satz nicht“ (Nachbar-Form Zitat; hier Bestand) | `harness/README.md` · `(sensors/docs-check.md) §Bindung · seit slice-d-check-tracked-modul` | ja — Aufschlagen der Stelle | Abschnitts-Zeiger zeigt auf die falsche Stelle |
| F-5 | INFO | Die Regel begrenzt `Vertrag` und `Tut was`, nicht `Bindung`; diese Spalte trägt bis zu sieben Herkunfts-Anker (`make test-sdk-*-integration`). Grenze 11 in `harness/sensors/docs-check.md` nennt die Bindung-Spalte nicht als ungemessen, und ihr Satz „eine neue Tabelle … mit einem anderen Spaltennamen fällt aus der Regel“ ist hergeleitet, nicht gemessen; die gemessenen Teile der Grenze hat dieser Lauf unabhängig bestätigt (Probe 3). | `AGENTS.md` §3.12 Instanz B | `harness/sensors/docs-check.md` · `Eine neue Tabelle im Abschnitt mit einem` | ja — Mutation im Klon | Grenze ohne Ursprungsangabe |
| F-6 | INFO | Die Befund-Liste nennt Chronik-Sprache nur bei `make image-cve`; dieselbe Klasse steht wortgleich umgezogen auch in `harness/targets/test-integration.md` („inhaltlich über den ursprünglichen MVP-Zuschnitt hinausgewachsen“) und `harness/targets/image.md` („jetzt als Multi-Arch-Manifestliste“). Kein Fehler des Umzugs (§1 verbietet die Berichtigung), nur eine unvollständige Liste. | `AGENTS.md` §3.7 | `harness/targets/test-integration.md` · `über den ursprünglichen MVP-Zuschnitt hinausgewachsen` | nein | Befund-Liste unvollständig |
| F-7 | INFO | In `make doc-trace` steht der neue Datei-Link zwischen „kein Gate,“ und „wie `make image-stale`“; der Vergleich bezieht sich jetzt dem Wortlaut nach auf den Link. Lesbarkeit, kein Fehler. Für den Planner: der Kopf `**Bezug:**` des Plans verlinkt `ADR-0051` auf den ADR-Index und beschreibt sie als „Harness-Pflege und Pin-Inventar“, die ADR heißt „CI/CD-Pipeline über GitHub Actions“ (nicht im Diff). | Maintainability | `harness/README.md` · `(targets/doc-trace.md), wie \`make image-stale\`` | nein | Einfüge-Position im Satz |

## Geprüfte Schwerpunkte — Ergebnis

- **Wortlaut-Erhalt (Risiko 1):** trägt. Gegenprobe nachgefahren (Probe 1), unabhängige Rückrichtung byte-gleich für die 26-KB-Zelle und ohne erfundenen Absatz für alle 30 Dateien (Probe 2). Stichprobe Wort für Wort zusätzlich gelesen: `harness/targets/bench.md`, `harness/targets/sdk-integration.md` (Kotlin-Abschnitt), `harness/sensors/coverage-gate.md` (Fassung), `harness/sensors/gates.md` — mit der Zelle am Parent gleich bis auf `\|` und Link-Präfixe.
- **Tragen die gekürzten Sätze die Zeile (Risiko 4):** ja. Jede der 63 Zellen ist ein ganzer Satz bzw. eine ganze Nominalphrase mit Gegenstand und Wirkung, keine zerhackten Satzreste; die längsten liegen bei 162 bzw. 113 Zeichen, also mit Abstand unter den Schwellen — kein Druck auf Pflichterfüllung sichtbar. `make gates` wurde vom Aufzählen auf „alle inneren Gates dieser Tabelle“ umformuliert; die Tabelle führt genau die zehn Gates der Aufzählung (Plan-Nachzug 2).
- **Neue `structure`-Regel:** beide Tabellen erfasst, Rand- und Platzhalterzeilen eingeschlossen, gezählt werden Zeichen des Rohtexts (Probe 3); der Kommentar in `.d-check.yml` trägt Zusage und Grenze, kein Vorher/Nachher.
- **Abweichung „ganze Zelle in bestehende Dateien“:** kein Verstoß gegen den Ausschluss „Neuschnitt“ (nichts ersetzt); der Nachzug trägt, §1 ist nicht mitgezogen (F-3), die Folge im Träger ist F-1.
- **Suchlauf und Meldungen:** Zahlen stimmen (Probe 4); zwei Träger fehlen in den Meldungen (F-2); die Positionsverweise der umgezogenen Zellen sind genau die zwei aus Befund 8 (`siehe Zeile darunter`, `siehe \`make schema-rollout\`-Zeile oben`), kein weiterer gefunden.
- **Link-Form und Kennungen:** jede geänderte Bindung-Zelle verlinkt ihre Datei in Link-Form; die umgezogenen Links sind relativ zum neuen Ort umgeschrieben und lösen auf (`make docs-check` Exit 0, Probe 3); Einleitungssatz verlinkt `../README.md#sensors-feedback-gates`. Kennungen in den neuen Dateien stehen in der Form der alten Zelle.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `harness/README.md` (Source precedence, Guides, Prosa-Abschnitte) | geprüft, ohne Befund — `git diff` zeigt nur die zwei Tabellen-Hunks |
| `harness/targets/` (13 neue, 4 ergänzte Dateien) | geprüft, ohne Befund außer F-1, F-6 |
| `harness/sensors/` (2 neue, 11 ergänzte Dateien, `docs-check.md`) | geprüft, ohne Befund außer F-1, F-4, F-5 |
| `.d-check.yml` | geprüft, ohne Befund (Probe 3) |
| Commit-Messages c967b51b, 3f900161, 10ad71b0, 76c3b9e0 | geprüft, ohne Befund — je `ADR-0051`, keine `SPEC-*`/`ARC-*` im Betreff |
| `AGENTS.md` §3.1 (Docker-only, kein in-place-/Umleitungs-Schreiben) | geprüft, ohne Befund — der Diff zeigt keine Formatspur eines Host-Werkzeugs; die Gegenprobe des Plans nutzt nur Lese-Werkzeuge und `mktemp` |
| `AGENTS.md` §3.5 (Accepted-ADRs) | geprüft, ohne Befund — keine ADR im Diff; die ADR-Zitate der README-Zeilen sind korrekt als eingefroren gemeldet |
| Plan §3 Zahlen (§3.12) | geprüft, ohne Befund — alle nachgemessen (Probe 5) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Nachzug widerspricht dem Nachbarn im selben Träger · Träger-Meldung unvollständig · Abgrenzung nicht mitgezogen bei Planänderung · Abschnitts-Zeiger zeigt auf die falsche Stelle · Grenze ohne Ursprungsangabe · Befund-Liste unvollständig · Einfüge-Position im Satz

## Verdikt

**Merge-blockierend:** ja, für F-1 (MEDIUM) — eine Fixrunde am Implementer; F-2 und F-3 gehen in dieselbe Runde, die INFO-Punkte sind Hinweise ohne erwartete Aktion. Die DoD-Zeile „Review durchgeführt“ bleibt deshalb offen (Nachzug in der Fixrunde, Skill §DoD-Checkbox-Nachzug).

**Übergabe:** Findings an den Implementer; die Finding-Klassen gehen in die Closure §7. Kein Rollen-Widerspruch, kein Konflikt-Pfad. Dieser Report ersetzt keine Verifikation (Modul 11).
