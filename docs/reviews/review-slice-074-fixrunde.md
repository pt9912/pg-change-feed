# Review-Report: slice-074 — Fixrunde — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports. Folgelauf zum Erstbefund
[`review-slice-074.md`](review-slice-074.md) — neue Datei, keine
Überschreibung (Modul 10 §Report je Lauf).

**Gegenstand:** die zwei Commits nach dem Erstbefund — `f4164a9` (Planner:
§1/§7/§8-Nachzug) und `e375209` (Implementer: Fixrunde zu F-1 und F-4).
`e375209` berührt vier Pfade (`.d-check.yml` · `harness/sensors/docs-check.md`
· `tools/harness/run-integration-tests.sh` · `docs/user/e2e-abdeckung.md`).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd`.
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext:** wie im Erstbefund, zusätzlich `docs/plan/planning/in-progress/slice-074-e2e-abdeckungstabelle.md`
in der Fassung `f4164a9`, die vier Dateien aus `e375209` und der Register-Stand
`BEO-PGC/test-runner-stiller-ausschluss` ·
`BEO-PGC/generierte-artefakte-ohne-sync-sensor` (beide `state.md`: 1×, offen).

---

## Status der Findings aus dem Erstbefund

| Erstbefund | Kategorie | Stand nach diesem Lauf |
|---|---|---|
| F-1 Zusage über die `structure`-Regelfamilie überdeckt den Config | MEDIUM | **geschlossen** — beide Sätze halten dem Config stand (§Frage 1) |
| F-2 §8-Sichtung lässt den Register-Eintrag der eigenen Sub-Area aus | MEDIUM (Planner) | **geschlossen** durch `f4164a9` |
| F-3 Beschreibungsspalte nicht kennungsfrei | LOW | **unverändert reproduziert** (`:20`, `:30`) — Einschätzung bestätigt |
| F-4 Stabilitäts-Zusage nennt eine der beiden Ursachen | LOW | **geschlossen** — der Kopf nennt jetzt beide, und die Zusage ist real eingetreten (§Frage 3) |
| F-5 Ableitungs-Scope enger als der Ausführungs-Scope | LOW | unberührt (kein Code-Pfad dieses Diffs angefasst) |
| F-6 Erzeugnis ohne eingehenden Verweis | INFO | **unverändert reproduziert** — Einschätzung bestätigt |
| F-7 Anker-Prüfung prüft Vorkommen, nicht Lage | INFO | unberührt |

Kein Finding des Erstbefunds ist offen geblieben.

## Findings dieses Laufs

### F-1 — Die neue Formulierung nennt für vier der fünf Tabellen-Regeln eine Teilmenge ihres Spalten-Mechanismus

- `kategorie`: INFO
- `quelle`: Maintainability · `AGENTS.md` §3.7
- `pfad`: `harness/sensors/docs-check.md:22-23`
- `befund`: Der Satz „die vier Tabellen-Regeln und die E2E-Abdeckungstabelle
  adressieren ihre Spalten über Mindestbreiten" ist wahr, aber nicht
  vollständig: die vier Tabellen-Regeln tragen **zusätzlich**
  `cell-max-chars` je Spalte (ADR-Index, Pflichtenheft-Defaults, beide
  Architektur-Tabellen), nur die E2E-Abdeckungstabelle trägt ausschließlich
  Mindestbreiten. Die Aussage behauptet keine Vollständigkeit und ist damit
  keine Überdeckung im Sinne des Erstbefunds F-1; sie nennt die Maxima nur
  nicht. Kein Verhalten, kein Gate berührt.
- `verifizierbar`: ja — `.d-check.yml` §structure, Regeln 1–4 gegen Regel 8.
- `klasse`: „Bündelnde Beschreibung nennt eine Teilmenge ohne
  Vollständigkeits-Signal" (1×)

### F-2 — Die Änderungsursachen des Kopfes sind aufgezählt; eine dritte Ursache existiert und trat in diesem Commit ein

- `kategorie`: INFO
- `quelle`: Maintainability · Slice-Plan §6 Risiko 4
- `pfad`: `tools/harness/run-integration-tests.sh:84-98` (`ABDECKUNG_KOPF`) ·
  `docs/user/e2e-abdeckung.md:8-14`
- `befund`: Der neue Kopf zählt zwei Ursachen auf — die Nachweis-Deklarationen
  und den Ort ihrer Quellen. Eine dritte, real auftretende Ursache ist der
  **Kopf-Text selbst**: `ABDECKUNG_KOPF` ist Teil des Erzeugnisses, eine
  Änderung an ihm schreibt die Datei, ohne dass eine Deklaration oder eine
  Quellzeile sich bewegt — in `e375209` genau so geschehen (der Kopf wuchs um
  drei Zeilen). Der Satz ist dadurch nicht falsch: das Schreiben „nur bei
  inhaltlicher Abweichung" deckt den Kopf mit ab; benannt ist er in der
  Aufzählung nicht. Die Umbenennung einer Quelle (Datei- oder Pfadwechsel)
  fällt dagegen unter den **Oberbegriff** „der Ort ihrer Quellen" — die
  Beispiel-Klausel der Zeile nennt nur Einfügungen.
- `verifizierbar`: ja — `git show e375209 -- docs/user/e2e-abdeckung.md` zeigt
  den Kopf-Text als geänderten Teil des Erzeugnisses.
- `klasse`: „Änderungsursachen-Aufzählung lässt den Erzeugnis-Kopf aus" (1×)

## Antworten auf die vier Fragen des Nachlaufs

1. **F-1 hält dem Config stand — ja.** Regeln selbst gezählt
   (`.d-check.yml` §structure, mechanisch aus der Datei gelesen): acht Regeln,
   **fünf** mit `table`-Knoten (ADR-Index · Pflichtenheft-Defaults · zwei
   Architektur-Tabellen · E2E-Abdeckungstabelle), **drei** ohne (Closure-Notiz
   · `docs/reviews/**` · `observations/**/observation.md`). Der Satz nennt
   „die vier Tabellen-Regeln und die E2E-Abdeckungstabelle" = fünf mit
   Mindestbreiten und „die Closure-Notiz-Regel und die beiden
   Verweisform-Regeln" = drei ohne Spalten-Knoten — beide Zahlen stimmen. Die
   Aufzählung der Muster-Keys
   (`non-empty`/`max-open-tasks`/`require-pattern`/`forbid-pattern`) ist
   vollständig für diese drei Regeln; das zusammenfassende Wort „Muster" ist
   ein weiter Mantel für `non-empty`/`max-open-tasks`, wird aber im selben Satz
   durch die exakte Key-Liste aufgelöst — keine Überdeckung. „Keine Regel
   zählt Zeilen" ist wahr: das Config kennt keinen Zeilen-/Zeilenzahl-Schlüssel.
   Der `.d-check.yml`-Kommentar nennt die Grenze jetzt richtig (Mindestbreiten
   fangen leere/zu kurze Zellen; einen fehlenden Link trägt `ids`), und seine
   Beispielrechnung (drei unverlinkte Kennungen = 43 Zeichen bei
   `cell-min-chars: 40` bleiben grün) habe ich **unabhängig** nachgestellt
   (§Messungen). Restpunkt, keine Überdeckung: Erstbefund F-1 dieses Laufs
   (Maxima ungenannt).
2. **Nichts Verhaltensrelevantes nebenbei — bestätigt.** `.d-check.yml`: die
   Nicht-`#`-Zeilen sind zwischen `972851e` und `HEAD` **unverändert**
   (`git diff 972851e..HEAD -- .d-check.yml | grep -v '^[+-][[:space:]]*#'`
   leer) — der Commit ist kommentar-only. `harness/sensors/docs-check.md` und
   `tools/harness/run-integration-tests.sh`: je **ein** Hunk, beide innerhalb
   von Prosa bzw. des single-quoted `ABDECKUNG_KOPF` (kein Kommando, kein
   `-run`-Muster, keine Deklaration angefasst);
   `test/integration/integration_test.go` ist von `e375209` **unberührt**.
   Das Erzeugnis ist **nicht handkorrigiert**: ein frischer Lauf der Go-Hälfte
   im Toolchain-Container plus Nachrechnung der Bash-Hälfte (24 Anker) plus
   Rendering ergibt **byte-gleich** dieselbe Datei (§Messungen). Die
   Nebenwirkung der Kopf-Erweiterung ist real und exakt so, wie der
   Implementer sie nennt: alle 24 Bash-`Ort`-Zellen wanderten um **+3**, die
   13 Go-`Ort`-Zellen blieben unverändert, keine Zeile kam hinzu oder fiel weg.
3. **F-4 ist vollständig bis auf eine Ursache — siehe Erstbefund F-2 dieses
   Laufs.** Zwei Ursachen sind jetzt benannt, die Klausel über Einfügungen ist
   wahr (in diesem Commit real eingetreten); die dritte (der Kopf-Text selbst)
   ist unbenannt, aber vom Schreib-Kriterium gedeckt. Umbenennungen fallen
   unter den Oberbegriff.
4. **Der Planner-Nachzug liest sich konsistent — ja, und er schreibt keine
   neue Unwahrheit über das Artefakt.** §8 führt
   `BEO-PGC/generierte-artefakte-ohne-sync-sensor` als Treffer und die
   Schlusszeile zählt beide Treffer; §7 erwartet `evidence/slice-074.md` in
   **beiden** Einträgen, und „Zähler dann je 2×" ist gegen die Register-Stände
   richtig (beide heute 1×, `evidence/slice-033.md` bzw. `slice-069.md`); der
   neue §1-Absatz benennt die dritte Richtung wahr (die `structure`-Regel
   sichert die **Existenz** — `section-missing` —, kein Sensor hält den Inhalt
   ohne den vollen E2E-Lauf gegen die Quelle). Die beiden Richtungen des
   ursprünglichen Satzes („die Tabelle behauptet einen Nachweis, den es nicht
   gibt" / „ein Nachweis existiert, fehlt aber in der Tabelle") bleiben
   daneben stehen, ohne dass der Nachzug sie relativiert — sie sind aber mit
   der neuen Einschränkung verträglich, wenn man sie auf den **erzeugten
   Stand** liest (gegeben ein Lauf), und genau diesen Fall fixiert die
   §5-Closure-Bedingung („der committete Stand entspricht der Ausgabe eines
   frischen Laufs"). Die Staleness entsteht danach, zwischen den Closures
   späterer Slices — das ist die neue dritte Richtung. Ich sehe darin keine
   Unwahrheit, sondern die Grenze des Satzes, eine Paragraphenzeile weiter
   benannt.

## Randposten — bestätigt oder verworfen

- **F-3 (Beschreibungsspalte nicht kennungsfrei): bestätigt, unverändert, LOW.**
  Die Datei trägt weiter `BEO-PGC/lese-doppelquelle` in der
  Beschreibungsspalte (`docs/user/e2e-abdeckung.md:20`) und „im Slice-Plan
  (§6) benannten Fällen" (`:30`) — die Zeilennummern wanderten um +3, der
  Inhalt nicht. Der Nachzug hat die Klasse nicht verschärft: `BEO-<KUERZEL>/<slug>`
  ist eine Kennung nach MR-000, `ids` prüft sie nicht, und `399370e` hat
  dieselbe Provenienz-Klasse (Testfall-Herkunft) auf der anderen Seite
  entfernt. Bleibt LOW, kein Gate berührt.
- **F-6 (kein eingehender Verweis): bestätigt, unverändert, INFO.** In keiner
  Doku verlinkt; die Treffer außerhalb der Datei sind Maschinerie
  (`run-integration-tests.sh`, `integration_test.go`), die Config-Zeile
  (`.d-check.yml:195`) und eine **Erwähnung** in der Grenzen-Liste
  (`harness/sensors/docs-check.md:62`) — kein Zug aus einem lesenden Knoten
  (`README.md`, Benutzerhandbuch, `harness/README.md` §Source precedence führt
  `docs/user/*` entlinkt). Kein Modul meldet verwaiste Dateien. Bleibt INFO.

## Negativbefunde

- geprüft, ohne Befund: **`docs/user/e2e-abdeckung.md`** — byte-gleiche
  Reproduktion eines frischen Laufs (§Messungen), keine Handkorrektur.
- geprüft, ohne Befund: **`.d-check.yml`** — kommentar-only, die neue
  Kommentar-Zusage ist gegen die acht Regeln wahr.
- geprüft, ohne Befund: **`harness/sensors/docs-check.md`** — die
  Regelfamilie-Zahlen stimmen; keine dritte überdeckende Aussage jenseits von
  F-1 dieses Laufs.
- geprüft, ohne Befund: **`tools/harness/run-integration-tests.sh`** — nur
  der Kopf-Text; Deklarationen, Anker-Lookup, Abbruch-Wächter und der
  Schreib-Pfad sind zeichengleich zum Stand des Erstbefunds.
- geprüft, ohne Befund: **Slice-Plan §1/§7/§8** nach `f4164a9` — siehe
  Frage 4.

## Messungen dieses Laufs (Exit-Code jeweils direkt, ungepiped)

| Messung | Ergebnis |
|---|---|
| `make docs-check` | Exit `0` — „601 Datei(en) geprüft, 0 Befund(e)" |
| `make gates` | Exit `0` (coverage 49,30 % ≥ 35 %, commit-traceability OK, a-check 0 Befunde) |
| `make test` | Exit `0` |
| Reproduktion des Erzeugnisses (Go-Hälfte real im Toolchain-Container, 13 Zeilen; Bash-Hälfte nachgerechnet, 24 Zeilen; Rendering) | sha `8c0c8cac2631259e2e759a04991ea985d543b9a1093968d2ed690887e9c2142d` — **IDENTISCH** mit der committeten Datei |
| Verschiebung der `Ort`-Zellen `972851e` → `HEAD` | 13 Go-Zellen **unverändert**, 24 Bash-Zellen **alle +3**, keine Zeile hinzugefügt/entfernt |
| Regeln des Configs gezählt (mechanisch aus `.d-check.yml` gelesen) | 8 Regeln: 5 × `table`/`column` (Mindestbreiten), 3 × Muster-Keys ohne Spalten-Knoten; kein Zeilenzahl-Schlüssel im Config |
| `.d-check.yml`-Diff `972851e..HEAD` ohne `#`-Zeilen | leer — kommentar-only |
| Hunk-Zahl je Datei in `e375209` | `.d-check.yml` 1 · `docs-check.md` 1 · `run-integration-tests.sh` 1 · Erzeugnis 2 (Kopf + verschobene Zeilen) |
| Arbeitsbaum nach allen Mutationen | `git status --porcelain` leer |

Nicht gefahren: der volle `make test-integration` (Compose-Kette). Die
Reproduktion oben ist die vom Erstbefund benannte Methode und **stärker** als
ein Lauf-Vergleich: sie zeigt, dass die committete Datei genau die Ausgabe der
Pipeline aus dem heutigen Quelltext ist.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Bündelnde Beschreibung nennt eine Teilmenge
ohne Vollständigkeits-Signal" (F-1, 1×) · „Änderungsursachen-Aufzählung lässt
den Erzeugnis-Kopf aus" (F-2, 1×). Beide sind Hinweise ohne erwartete Aktion;
die Klassen des Erstbefunds bleiben unverändert stehen und werden **nicht**
erneut gezählt (ein Vorgang zählt einmal, Modul 6 §Das Beobachtungs-Register).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Die beiden Fixes tragen: der
Config hält der neuen Formulierung stand (Zahlen unabhängig nachgezählt), und
der Kopf nennt beide praktisch relevanten Änderungsursachen, was die
Reproduktion zugleich belegt. Der Planner-Nachzug ist konsistent, seine
Zähler-Arithmetik stimmt gegen die Register-Stände.

**Übergabe:** keine Rückgabe an den Implementer — keine Fixrunde. F-1 und F-2
dieses Laufs sind Hinweise ohne erwartete Aktion; F-3 und F-6 bleiben wie im
Erstbefund benannt (LOW/INFO) und liegen beim Planner in der Closure; F-5 und
F-7 sind unberührte Grenzen des Erstdiffs.

**DoD-Nachzug:** **ja.** Nach diesem Lauf bleiben 0 HIGH, und kein Finding
trägt einen Reviewer→Implementer-Rückgabe-Pfeil. Die Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des Slice-Plans ist
deshalb in demselben Commit wie dieser Report auf `[x]` nachgezogen, mit
Verweis auf Erstbefund und diesen Bestätigungslauf
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Nur diese
eine Zeile — Closure-Notiz, Risiko-Ausgänge, Beobachtungs-Register und die drei
Paarungen bleiben Planner-Arbeit. Der Report ersetzt keine Verifikation:
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
