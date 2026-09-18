# Review-Report: Zitat-Korrektur Beobachtungsregister Batch 1 (`faccc97`) — 2026-09-18

**Review-Art:** Code — geprüft gegen [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) (Zitat-Korrektur an immutablen Dokumenten), `ADR-0094`/[`ADR-0097`](../plan/adr/0097-observation-matrixklasse-review-verboten.md) (Kontext der Batch-Aufteilung), `AGENTS.md` §3.5.

**Gegenstand:** Commit `faccc97` — "docs(observations): Zitat-Korrektur
Review-/Verify-/Verdikt-Adressen ([`ADR-0073`](../plan/adr))", 42 Beobachtungsregister-Dateien
unter `docs/plan/planning/observations/BEO-PGC/*`.

**Skill:** `.harness/skills/reviewer.md` @ Arbeitsbaum-Stand 2026-09-18
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) (Zitat-Korrektur-Erlaubnis und ihre Grenze)
- `ADR-0094`, [`ADR-0097`](../plan/adr/0097-observation-matrixklasse-review-verboten.md) (Kontext: warum 42 von insgesamt mehr Dateien in
  diesem Batch laufen)
- `AGENTS.md` §3.5 (ADR-Immutabilität, Zitat-Korrektur-Ausnahme)
- Commit `faccc97` selbst (Diff, Commit-Message)

---

## Findings

### F-1 — Evidence-Datei zum Review-Zug `slice-041`: Vorgang-Feld benennt jetzt einen Commit statt der Datei, Referent leicht mehrdeutig

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/observations/BEO-PGC/commit-traceability-kein-vorab-hook/evidence/` (Evidence-Datei zum Review-Zug `slice-041`), Zeile 3
- `befund`: Die Zeile „Vorgang: <Review-Report zu `slice-041`>
  (Reviewer-Lauf zu slice-041)" wurde zu "Vorgang: der Commit, der den
  Review-Report zu slice-041 abschloss (Reviewer-Lauf zu slice-041)".
  Git-Historie zeigt zwei Commits für diese Datei (`f98b0b3` erstellt sie,
  `4e4c7bb` — der im nachfolgenden "Fund"-Absatz beschriebene
  [SPEC-016](../../spec/pflichtenheft.md)-Verstoß — ergänzt nur Links). "der Commit, der … abschloss" ist
  dadurch nicht eindeutig demselben Commit zuordenbar, der im Fund-Absatz
  als Verstoß genannt wird — der ursprüngliche Dateipfad war als Referenz
  eindeutiger, auch wenn er als Adresse jetzt unerwünscht ist. Der Referent
  selbst (die Tatsache, dass `4e4c7bb` der [SPEC-016](../../spec/pflichtenheft.md)-Verstoß ist) bleibt in
  "Fund" unverändert korrekt.
- `verifizierbar`: nein — reine Lese-/Interpretationsfrage, kein Gate prüft
  Referenz-Eindeutigkeit in Prosa.
- `klasse`: Zitat-Korrektur führt neue Mehrdeutigkeit im Referenten ein

## Negativbefunde

- geprüft, ohne Befund: Stichprobe von 10 der 42 Batch-1-Dateien im Detail
  (`adapter-fehler-ausgang/evidence/slice-007.md`,
  `adr-folgepflicht-ohne-traeger-slice/evidence/slice-098.md`,
  `arbeit-ueberholt-stehenden-traeger/evidence/slice-095.md`,
  `arbeit-ueberholt-stehenden-traeger/evidence/slice-096.md`,
  `arbeit-ueberholt-stehenden-traeger/evidence/slice-097.md`,
  `beleg-befehl-traegt-seinen-satz-nicht/evidence/slice-084.md`,
  `beleg-befehl-traegt-seinen-satz-nicht/evidence/slice-085.md`,
  `dod-checkbox-nachzug/state.md`,
  `git-mv-und-inhalt-in-einem-commit/evidence/slice-097.md`,
  `mechanismus-erklaerung-ohne-werkzeugbeleg/state.md`) — Referent und
  Aussage in jedem Fall unverändert, nur Adressform (Basisname →
  Prosa-Adressierung, teils Markdown-Link → Fließtext) ersetzt.
- geprüft, ohne Befund: die vier Dateien, deren Dateiname selbst Gegenstand
  der jeweiligen Beobachtung ist — Evidence-Dateien zu den Review-Zügen
  `slice-041` (unter `commit-traceability-kein-vorab-hook`), `slice-070` und
  `slice-077` (unter `kommentar-behauptet-nicht-getragenen-fehlerpfad`) sowie
  `slice-080` (unter `mechanismus-erklaerung-ohne-werkzeugbeleg`)
  — Dateiname selbst unverändert (kein `git mv`, korrekt: der Dateiname ist
  keine "Adresse" im Sinne von `ADR-0073`, nur ihr Inhalt zitiert
  Basisnamen), Überschrift und Fließtext auf Prosa-Form umgestellt; einzige
  Auffälligkeit F-1 oben (Review-Zug `slice-041`).
- geprüft, ohne Befund: die zwei bewusst unveränderten
  `docs/reviews/**`-Erwähnungen
  (`arbeit-ueberholt-stehenden-traeger/evidence/slice-105.md`,
  `exemption-ohne-reifegrenze/evidence/slice-105.md`) — beide zitieren
  ausschließlich den `.d-check.yml`-Config-Glob-Wert `docs/reviews/**`
  (Verzeichnismuster, kein konkreter Report-Basisname); zu Recht nicht
  angefasst.
- geprüft, ohne Befund: Commit-Umfang gegen die Commit-Message — `git
  diff-tree --no-commit-id --name-only -r faccc97 -- 'docs/plan/planning/observations/*'`
  liefert exakt 42 Dateien, deckungsgleich mit der behaupteten Zahl.
- geprüft: `make gates` — Exit 0, direkt und ungepiped geprüft (kein
  `.d-check.yml`-Matrix- oder `hostpaths`-Befund gegen die geänderten
  Dateien).

## Beobachtungen außerhalb des engeren Diffs (nicht Teil der 42 Dateien)

**B-1 — Scope-Creep in `faccc97` selbst, bereits vor diesem Review
selbstständig korrigiert.** `git show faccc97 --stat` zeigt neben den 42
Beobachtungsregister-Dateien zwei weitere Änderungen, die nicht zur
"Zitat-Korrektur"-Aufgabe gehören: eine neue Datei
`docs/plan/adr/0097-observation-matrixklasse-review-verboten.md` und eine
Index-Zeile in `docs/plan/adr/README.md`. Der unmittelbar folgende Commit
`bc49fab` ("[`ADR-0097`](../plan/adr)-Index-Zeile und Datei aus letztem Commit entfernt")
nimmt beides zurück mit der Begründung: "Der vorige Commit hat versehentlich
zwei nicht zu diesem Auftrag gehoerende, parallel in Arbeit befindliche
Aenderungen mitgenommen … stammen aus einem parallelen Bearbeitungsschritt".
Ein weiterer Commit `a9508b3` fügt [`ADR-0097`](../plan/adr/0097-observation-matrixklasse-review-verboten.md) danach sauber isoliert wieder
ein. Der aktuelle Arbeitsbaum ist dadurch unauffällig ([`ADR-0097`](../plan/adr/0097-observation-matrixklasse-review-verboten.md) korrekt
vorhanden, keine Doppelung) — dies ist **kein** Befund am geprüften Commit
selbst mehr, sondern die reale Bestätigung des Risikos, vor dem die
Aufgabenstellung dieses Reviews warnt (`git add -A` in einem parallel
bearbeiteten Arbeitsbaum). Klasse: Rebase-/Interleaving-Risiko bei
parallelen Implementer-Batches im selben Arbeitsbaum — als INFO vermerkt,
keine Aktion nötig.

**B-2 — Repo-weite Gegenprobe zeigt deutlich mehr offene Zitate als in der
Auftragsbeschreibung erwartet.** `grep -rl` über
`docs/plan/planning/observations/` nach
`docs/reviews/\|review-slice-\|verify-slice-\|architect-verdict-\|architect-review-`
liefert **122 Dateien** (Basisnamen-Substring-Variante ohne
`.md`-Endungspflicht) bzw. **104 Dateien** bei strengerer Form mit
erzwungener `.md`-Endung nach dem Basisnamen — in beiden Fällen weit über
den ~41 erwarteten "Batch 2"-Dateien (83 laut `ADR-0097` §Kontext minus 42
hier). Nach Abzug der 42 Batch-1-Dateien bleiben davon exakt 2 Dateien mit
dem bekannten, bewusst unveränderten `docs/reviews/**`-Config-Glob-Muster
(s. o.) und **100+ weitere** Dateien mit noch offenen Basisnamen-Zitaten
(stichprobenartig gesichtet, z. B.
`a-check-null-abdeckung/evidence/slice-001.md`,
`arbeit-ueberholt-stehenden-traeger/evidence/slice-091.md`,
`architect-verdikt-ablageort-uneinheitlich/observation.md`). Das ist **kein
Mangel am geprüften Commit** — dessen eigene Zahlenbehauptung ("42 Dateien")
ist exakt zutreffend —, aber die "~41 verbleibend"-Erwartung aus der
Auftragsbeschreibung für Batch 2 trifft auf den realen Bestand nicht zu.
Plausible Erklärung: `ADR-0097`s "83 Dateien" bezog sich vermutlich auf die
für den konkreten `archive-welle`-Lauf tatsächlich betroffene Teilmenge
(Reports, die in dieser Welle archiviert werden), nicht auf den gesamten
Basisnamen-Zitat-Bestand im Register — das sollte vor der Planung von Batch
2 geklärt werden, sonst wird Batch 2 mit falscher Umfangserwartung geplant.
Klasse: Zahlenbehauptung im Umfeld ohne geprüften Ursprung (`AGENTS.md`
§3.12, hier: eine Erwartung des Auftraggebers, nicht des geprüften
Artefakts selbst — daher INFO, kein Finding gegen den Commit).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 (B-1, B-2 — außerhalb des engeren Diffs) |

**Finding-Klassen dieses Laufs:** Zitat-Korrektur führt neue Mehrdeutigkeit
im Referenten ein · Rebase-/Interleaving-Risiko bei parallelen
Implementer-Batches · Zahlenerwartung ohne geprüften Ursprung

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, ein einzelner LOW-Befund
ohne semantische Auswirkung (der Referent bleibt im Fund-Absatz eindeutig
belegt, nur das Vorgang-Feld ist eine Nuance schwächer als vorher). Für die
geprüften 42 Dateien gilt: Referent und Aussage sind in jedem gesichteten
Fall unverändert, nur die Adressform wurde ersetzt — die Commit-Message
("Referent und Aussage bleiben unverändert") trifft zu.

**Übergabe:** Kein Fixrunden-Bedarf am Implementer für diesen Batch. B-1 ist
bereits selbst aufgelöst (Folgecommits `bc49fab`/`a9508b3`) und braucht keine
weitere Handlung. B-2 sollte vor Planung von Batch 2 geklärt werden (welche
konkrete Dateimenge `ADR-0097`s "83" wirklich referenziert), damit Batch 2
nicht mit einer zu kleinen Umfangserwartung geplant wird — das ist eine
Empfehlung an Planner/Architect, kein Reviewer→Implementer-Rückgabe-Pfeil.
Dieser Report ist ein Lauf-Beleg; er ersetzt keine Verifikation (Modul 11).
