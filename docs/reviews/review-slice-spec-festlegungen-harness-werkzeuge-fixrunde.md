# Review-Report: slice-spec-festlegungen-harness-werkzeuge, Fixrunde — 2026-10-07

**Review-Art:** Code — enges Re-Review der Fixrunde gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `01264f8d..8ed4eb7f` (3 Commits: die Fixrunde `941b1207`,
das Architect-Verdikt `86bff10f`, die Festlegung nach dem Verdikt `8ed4eb7f`)

**Skill:** `.harness/skills/reviewer.md` @ 8ed4eb7f
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Erstes Review [`review-slice-spec-festlegungen-harness-werkzeuge`](review-slice-spec-festlegungen-harness-werkzeuge.md)
  (`01264f8d`; 0 HIGH, 3 MEDIUM, 5 LOW, 5 INFO)
- Slice-Plan `slice-spec-festlegungen-harness-werkzeuge`, §3 Absatz „Fixrunde“
  (Stand `8ed4eb7f`)
- [Architect-Verdikt zur Lesung von Exit 2 bei mehrdeutiger Stellung](architect-verdict-zitat-vergleich-mehrdeutig.md)
  (`86bff10f`, Lesung (a), ohne ADR)
- [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
  Entscheidung 3 und 4 sowie Status-Zeile;
  [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
  Status-Zeile, Entscheidung 2, §Konsequenzen;
  [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
  Status-Zeile;
  [`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
  Entscheidung 1 (Tabelle, (c)), 2 und 4
- [`SPEC-038`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
  (§7 vollständig am Stand `8ed4eb7f`), Vertrag `harness/targets/zitat-vergleich.md`
  vollständig, Plan `slice-spec-festlegungen-pruefer-hooks` §1
- `AGENTS.md` §3.5, §3.9, §3.12, §3.13; `v6.14.0` · `regelwerk/modul-08-agentenrollen.md`
  §Konflikt-Pfad als Rollen-Sequenz

**Eigene Läufe dieses Reviews** (Exit je direkt gesichert, `AGENTS.md` §3.9):

| Lauf | Ergebnis |
|---|---|
| `make suchlauf-nachmessen PLAN=<Plan>` | Exit 0, gedruckt `suchlauf-nachmessen: 13 Zeilen stimmen` (darunter `OK  soll=15 ist=15` für den Anker `7-festlegungen-der-harness-werkzeuge`) |
| `make test-zitat-vergleich` | Exit 0, gedruckt `run-zitat-vergleich-tests: 297 Fälle bestanden (je Runde 99, Runden: ohne Option, nullglob, failglob)` |
| `make kommentar-kennungen DIFF=01264f8d` | Exit 0, keine Ausgabe (kein Kandidat; der Diff trägt keinen Code) |
| `git grep -c formnorm -- docs/plan/planning/open/slice-spec-festlegungen-pruefer-hooks.md` | Exit 0, `…pruefer-hooks.md:2` |
| Stichprobe zu Verdikt §3: `git grep -n -o '<a id="[^"]*"' 86bff10f~1 -- '*.md' \| wc -l`; Link-Muster auf `attr`, `x`, `X`, `qa` bis `qf` am Stand `941b1207`; `tools/harness/zitat-vergleich.sh` je `#attr`, `#x`, `#qa` gegen sich selbst in `docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md` und `harness/targets/zitat-vergleich.md` | 96 Fundstellen; 0 Link-Zeilen; alle sechs Läufe Exit 2, drei davon mit `mehrdeutig` auf stderr — stützt die Aussage des Verdikts, dass heute kein Verweis auf einen Anker in mehrdeutiger Stellung zeigt |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | INFO | Träger außerhalb des Diffs (`AGENTS.md` §3.13, Meldung an den Planner, Frist: Closure dieses Slice): das Beobachtungs-Register `BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open` führt die mehrdeutige Stellung in `observation.md` mit der Glosse („nicht messbar, Urteil am Diff“), die das Verdikt verwirft; `state.md` verankert das Prinzip im Vertrag „§Einheit, „Mehrdeutig endet mit Exit 2““, und §Einheit trägt den Satz nicht mehr (er steht seit dem Nachzug in `SPEC-038`). `observation.md` ist ab Anlage unveränderlich, `state.md` nicht. | `AGENTS.md` §3.13; Verdikt §0 | `docs/plan/planning/observations/BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open/state.md` · „Vertrag `harness/targets/zitat-vergleich.md` §Einheit,“ | ja — `git grep -n 'Mehrdeutig endet' -- harness/targets/zitat-vergleich.md` (kein Treffer) | Arbeit überholt stehenden Träger |

### Prüfpunkt 1 — sind die Findings des ersten Reviews gelöst?

| Finding | Behauptung im Plan | Ergebnis |
|---|---|---|
| F-1 (MEDIUM) | Satz aus `ADR-0158` E3 in „Stände“; mehrdeutige Stellung nach Verdikt | **gelöst.** „Ist eine Einheit leer oder nicht lesbar (Anker unbekannt, Datei fehlt), besteht die Korrektur nicht.“ steht wortgleich zu `ADR-0158` Entscheidung 3 im Absatz „Stände“. Der Satz „Mehrdeutig heißt „nicht messbar““ ist durch den Wortlaut aus Verdikt §2 ersetzt, Zeichen für Zeichen gleich; die Historie-Zeile trägt den Anhang aus Verdikt §2 wörtlich. Der Absatz „nicht messbarer Referent“ nennt nur die Adresse außerhalb des Repos und die an keinem Stand auflösende, und der neue Satz schließt die mehrdeutige Stellung ausdrücklich aus. Kein Träger außerhalb der Records bindet „mehrdeutig“ noch an „nicht messbar“ außer dem Register (F-1 oben) |
| F-2 (MEDIUM) | ein Maßstab im Vertrag | **gelöst.** §Vertrag: Block der ADR „historische Fassung … und kein Maßstab“, Maßstab „allein die Festlegung“ mit Zeiger auf `ADR-0162` E2. §Abweichungen: „Alle zehn Punkte sind Festlegung“; die Tabelle hat zehn Zeilen, deckungsgleich mit der Aufzählung in `ADR-0162` E2. „Auslegung im Sinn von Verdikt §2“, „die drei letzten Zeilen“ und „außerhalb der Punkte“ sind weg (`grep` ohne Treffer) |
| F-3 (MEDIUM) | Übergabe `formnorm` als Text im Folge-Plan | **gelöst.** §1 des Plans von `slice-spec-festlegungen-pruefer-hooks` nimmt die Sendung an: Gegenstand (Vergleich je MR-Datei, eine Zeile je Datei, `cmp 0`, ohne Zeile falsch bestimmt), Quelle `ADR-0161` E4 zweiter Spiegelstrich, Übergabe `ADR-0162` E4, Zielort (eigene Zeile in §7, Vertrag `pin-stale.md`). `git grep -c formnorm` = 2 (oben) |
| F-4 (LOW) | als Befund vermerkt | **vermerkt** im Plan §3 Fixrunde, Altfall für die Closure |
| F-6 (LOW) | `sed`/`awk`-Fehler in die Exit-Tabelle | **gelöst.** Die Zeile Exit 2 von `SPEC-038` nennt „`sed` am Zeilen-Lokator oder `awk` am Anker endet mit Fehler“; der Satz im Vertrag ist entfernt |
| F-7 (LOW) | §Grenzen verweist direkt auf `SPEC-038` | **gelöst.** „(„mehrdeutige Stellung“ in `SPEC-038`)“ mit Link auf §7; der Anker-Zählwert 15 im Suchlauf deckt diese Zeile (gemessen, oben) |
| F-8 (LOW) | als Befund vermerkt | **vermerkt** im Plan §3 Fixrunde, Lese-Hinweis für die Closure |
| F-12 (INFO) | Message-Pflicht in `pin-stale.md` | **gelöst.** Der Absatz „Umzug eines abgelösten Eintrags“ nennt „die aufgelöste Kennung und `ADR-0162`“ wie die drei anderen Träger |

### Prüfpunkt 2 — widerspricht das Verdikt einer `Accepted`-ADR, trägt „ohne ADR“?

Kein Widerspruch, „ohne ADR“ trägt; daher kein Finding an den Architect.

- `ADR-0158` Entscheidung 3 ist in Kraft: `ADR-0159` löst aus `ADR-0158` nur die
  HTML-`id`-Zeile von E1, E2 und E6 ab und lässt „die Entscheidungen 3 bis 5“
  ausdrücklich stehen; `ADR-0161` berührt `ADR-0158` nicht. `ADR-0162`
  Entscheidung 1 nennt E3 als Quelle von „Stände“.
- `SPEC-038` führt die mehrdeutige Stellung schon vor der Fixrunde unter „Keine
  Einheit haben“ und in der Exit-Tabelle unter „eine Seite hat keine Einheit“;
  `ADR-0162` Entscheidung 2 nimmt „mehrdeutige Stellung endet mit Exit 2“ als
  fail-closed an. Lesung (a) wendet E3 auf eine Form an, die die Festlegung
  bereits als „keine Einheit“ führt; sie ändert keine Folge einer Entscheidung.
- `ADR-0158` Entscheidung 4 und `ADR-0162` Entscheidung 1 (c) binden „nicht
  messbar“ an eine Adresse außerhalb des Repos bzw. an keinem Stand auflösend;
  keine der Quell-ADRs nennt die mehrdeutige Stellung als Fall davon. Der alte
  Vertrag (`abbe11b4`) stützte „Urteil am Diff“ auf den Zweig von `ADR-0159` E2,
  der dieselbe Bedingung trägt — die Lesung (b) hatte keinen Träger.
- `ADR-0162` Entscheidung 2: „Ändert sich die Festlegung, wird `SPEC-038`
  fortgeschrieben … Eine Semantik-Änderung gegen eine `Accepted`-ADR bleibt eine
  Folge-ADR.“ Die Berichtigung bringt die Festlegung an die ADR heran, nicht von
  ihr weg — Verdikt 1 aus Modul 8 („ADR gilt, der Text hat falsch behauptet“),
  Übergabe-Artefakt das Verdikt selbst mit Rolleninhaber.
- Eine Lesart im Verdikt ist Urteil, kein Messwert: §1 Punkt 1 liest die Klammer
  „(Anker unbekannt, Datei fehlt)“ als Beispielliste. Die Lesung trägt ohne diesen
  Punkt ebenso (Punkt 2 bis 4); kein Befund.

### Prüfpunkt 3 — neue Widersprüche?

Keine zwischen `SPEC-038`, Vertrag, Verdikt und ADR-Kette. Im Einzelnen: der
neue Satz verweist auf „Stände“, und „Stände“ trägt den Satz aus E3; die
Lesung von Exit 2 im Adaptions-Durchgang („Prüfauftrag“, Absatz „Ausgänge“)
gilt dem Durchgang ohne Korrektur und widerspricht „besteht nicht“ nicht; der
Vertrag verweist für die Ausgänge und die mehrdeutige Stellung nur auf
`SPEC-038` und nennt keine eigene Lesung; der Vertrag §Grenzen „Was nicht im
Werkzeug steht“ (Befehlsform `formnorm` bleibt in der ADR) deckt sich mit dem
Absatz im Folge-Plan („bis dahin bleibt die Befehlsform `formnorm` in der ADR“).
Der einzige verbliebene Träger der verworfenen Lesung liegt außerhalb des Diffs
(F-1).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/pflichtenheft.md` §7 (`SPEC-038` vollständig, `SPEC-039` unverändert) und §8 Zeile 2026-10-07 gegen Verdikt §2 und `ADR-0158` E3 | geprüft, ohne Befund |
| `harness/targets/zitat-vergleich.md` vollständig: Maßstab, §Abweichungen (zehn Zeilen gegen `ADR-0162` E2), §Ausgaben, §Grenzen | geprüft, ohne Befund |
| `harness/targets/pin-stale.md` Absatz „Umzug eines abgelösten Eintrags“ | geprüft, ohne Befund |
| `docs/plan/planning/open/slice-spec-festlegungen-pruefer-hooks.md` §1 Übergabe (Adresse nimmt an) | geprüft, ohne Befund |
| `docs/plan/planning/in-progress/slice-spec-festlegungen-harness-werkzeuge.md` §3 Suchlauf (Zahlen nachgemessen) und Absatz „Fixrunde“ gegen die Commits | geprüft, ohne Befund |
| `docs/reviews/architect-verdict-zitat-vergleich-mehrdeutig.md`: Verdikt-Typ, Rolleninhaber, Stand, gemessen/hergeleitet getrennt (§3), Messung als Stichprobe nachgefahren | geprüft, ohne Befund |
| Träger der verworfenen Lesung im Baum außer Records (`git grep -i 'nicht messbar'`, `git grep -i mehrdeutig` ohne `docs/reviews`, `done/`, `.harness/baseline`, `docs/plan/adr`) | geprüft — einzige Fundstelle F-1; `AGENTS.md` §3.5 nennt „nicht messbar“ nur mit `ADR-0158` E4, ohne Bindung an die Stellung |
| Hard Rules §3.5 (keine `Accepted`-ADR im Diff geändert), §3.1 (kein Host-Werkzeug im Diff), §3.7 (kein Code-Kommentar im Diff) | geprüft, ohne Befund |
| Traceability der drei Commit-Messages (je `ADR-0158`/`ADR-0162`, keine `SPEC-*` im Betreff) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Arbeit überholt stehenden Träger

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. F-1, F-2, F-3, F-6, F-7 und
F-12 des ersten Reviews sind gelöst, F-4 und F-8 als Befund vermerkt; F-5, F-9,
F-10, F-11 und F-13 gehen wie im Plan vermerkt an die Closure. Keine weitere
Fixrunde: die DoD-Zeile „Review durchgeführt“ ist im Slice-Plan im selben
Commit wie dieser Report gesetzt (Skill §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** F-1 ist Planner-Sache (Träger außerhalb des Diffs, Frist: Closure
dieses Slice) und gehört mit der Klasse in die Closure-Notiz §7. Kein HIGH mit
Rollen-Widerspruch, kein Konflikt-Pfad. Dieser Report ersetzt keine
Verifikation (Modul 11).
