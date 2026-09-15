# Review-Report: slice-079 (Fixrunde) — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität, §6-Risiko-Ausgänge, Beobachtungs-Register und die drei
Paarungen sind **nicht** Gegenstand dieses Reports (Verifier bzw.
Planner-Closure, Modul 11/6).

**Gegenstand:** `slice-079`, Fixrunde — Diff `a34f7c6..5d3ac6d`, genau zwei
Dateien: `docs/plan/…/slice-079-…md` (aus `0aa4d37`, Planner) und
`harness/sensors/coverage-gate.md` (aus `5d3ac6d`, Implementer). **Frisch
geprüft, nicht als Bestätigung des ersten Reports gelesen** — jede Aussage der
Fixrunde ist an der Messung nachgeprüft. Vorgänger: `review-slice-079.md`
(derselbe Lauf, eigene Datei je Lauf).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` · **Modell:**
deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext:**

- Fixrunde `5d3ac6d` und Plan-Fix `0aa4d37` vollständig, dazu der
  Plan-Stand `a34f7c6..5d3ac6d`
- Erster Report dieses Slice (`review-slice-079.md`, 0 HIGH / 3 LOW / 5 INFO)
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Entscheidung Punkte 1–3, §Fitness Function, §Re-Evaluierungs-Trigger ·
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(a)
- `AGENTS.md` §3.1, §3.7, §3.9, §4 · `harness/mk/coverage.mk` (`THRESHOLD ?= 65`)
- `docs/plan/planning/welle-20.md` §1 (Träger der Test-Arbeit)
- Klassen-Vorgänger: `docs/reviews/review-slice-076.md` F-1
  (dokument-interner Zahlen-Widerspruch) und F-4 (benannte Lücke)

---

## Eigene Messungen dieses Laufs (Stand `5d3ac6d`)

| Lauf / Probe | Exit bzw. Ergebnis |
|---|---|
| `make coverage-gate` (Stufe 65) | **0** — `coverage-gate: OK — Coverage 69.90% erfüllt Schwelle 65%` |
| `make gates` | **0** — coverage-gate OK **69,70 %** ≥ 65; d-check 648 Dateien/0 Befunde; a-check 0; baseline-verify `v6.5.0` OK |
| `go tool cover -func` auf dem extrahierten Profil | `total: (statements) 69.9%` — Statement-Summe aus dem Profil (Unique Blocks): **1679** |
| `go list -f '{{len .TestGoFiles}}'` über `./internal/... ./cmd/...` | **fünf** Pakete ohne Testdatei: `postgresstorage/queries`, `grpc/streamv1`, `application/port/inbound`, `domain/errors`, `cmd/pg-change-feed` |
| `go list -f '{{.GoFiles}} {{.TestGoFiles}}'` für `grpc/streamv1` | `[changestream.pb.go changestream_grpc.pb.go]` / `[]` |
| Profil je Paket (eigene Deduplizierung) | `cmd/pg-change-feed` 49 Statements / **0** gedeckt — das einzige Paket mit 0 gedeckten Statements; `grpc/streamv1` 86 / 61 (70,9 %) |
| Ausgabeform im Stage-Log | `cmd/pg-change-feed` **und** `grpc/streamv1` je als `<pkg>\t\tcoverage: 0.0% of statements`; die drei übrigen als `?  <pkg> [no test files]` |

---

## Findings

### F-1 — Die neue Aussage „`cmd/pg-change-feed` ist der einzige ganz ungetestete Gegenstand“ hält nur unter einer Lesart, die ihre eigene Begründung nicht nennt; `grpc/streamv1` fehlt weiterhin

- `kategorie`: LOW
- `quelle`: Maintainability · Slice-Plan §3 (die Sensor-Doku trägt den
  Messgegenstand) · `AGENTS.md` §3.7 · **zweites Vorkommen** der Klasse aus
  `docs/reviews/review-slice-079.md` F-1 (dort LOW, an derselben Stelle)
- `pfad`: `harness/sensors/coverage-gate.md:60-64`
- `befund`: Die Pakete **ohne jede Testdatei** im Messbereich sind fünf, nicht
  vier: `postgresstorage/queries`, `application/port/inbound`,
  `domain/errors`, `cmd/pg-change-feed` **und**
  `internal/adapters/driving/grpc/streamv1` (nur
  `changestream.pb.go`/`changestream_grpc.pb.go`, `TestGoFiles` leer). Der
  Absatz nennt für `cmd` als Begründung „hat keine Testdatei“ — eine
  Eigenschaft, die jenes fünfte Paket teilt — und listet daneben „die drei
  `[no test files]`-Pakete“; `grpc/streamv1` steht nirgends. Wahr ist die
  Spitzenaussage nur unter der Lesart „ungetestet = kein Statement gedeckt“
  (dort ist `cmd` mit 0 von 49 das einzige Paket), und diese Lesart trägt der
  Satz nicht: die drei genannten Pakete tragen **keine ausführbaren**
  Statements, `grpc/streamv1` trägt **86** (61 davon gedeckt, über fremde
  Testpakete) — der Leser schließt aus dem Absatz, im Messbereich sei jedes
  Paket entweder getestet, `cmd`, oder ohne ausführbare Statements.
- `verifizierbar`: nein — kein Gate deckt Doku-Aussagen; mechanisch
  nachprüfbar mit `go list -f '{{len .TestGoFiles}}'` in der gepinnten
  Toolchain-Stage
- `klasse`: „Messgegenstands-Aussage übertrifft die reale Paketlage“

### F-2 — Die Arithmetik der neuen §Grenze-Punkte ist für ihre eigene Basis korrekt, die Basis selbst ist aus dem Dokument nicht herleitbar

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (ein Wert, ein Ort; Ist-Zustand)
  · verwandte Klasse: `docs/reviews/review-slice-076.md` F-1
- `pfad`: `harness/sensors/coverage-gate.md:78-84` gegen `:32` (Ist-Stand) und
  `:80-85` (Beleg-Zeile)
- `befund`: Nachgerechnet stimmen alle drei Verhältnisse für die Basis 1167 —
  `1169/1702 = 68,68 %` · `1198/2289 = 52,34 %` · `1178/1834 = 64,23 %`. Die
  Basis 1167 widerspricht aber den eigenen Zahlen derselben Sektion: die
  Beleg-Zeile führt **69,70 %**, was über 1679 Statements auf 1170 bzw. 1171
  gedeckte Statements führt (`1167/1679 = 69,50 %` — ein solcher Lauf hätte
  `69.50%` gedruckt), und die Bezifferung (`2467 − 788 = 1679`, vorher 1216
  gedeckt, davon 44 wegfallend) führt auf 1172. Mit meinem eigens gemessenen
  Profil bei `5d3ac6d` (1173 gedeckt ⇒ 69,9 %) lesen sich die drei Fälle
  **69,04 % · 52,60 % · 64,56 %**; mit dem schwächeren meiner beiden Läufe
  (1170 ⇒ 69,7 %) **68,86 % · 52,48 % · 64,40 %**. Die drei **Ergebnisse**
  (grün · rot · rot) sind in jeder Rechnung identisch — nicht identisch sind
  die Zahlen, die das Dokument als Rechenweg führt.
- `verifizierbar`: ja — aus dem eigenen Profil nachrechnen
  (`go tool cover -func` plus Deduplizierung über die Block-Position)
- `klasse`: „Rechnung nennt eine Basis, die die eigene Ist-Stand-Zeile nicht trägt“

### F-3 — Der neue DoD-Wortlaut ist in der Sache exakt; die Klammer untertreibt die Ausgabemenge

- `kategorie`: INFO
- `quelle`: Maintainability · Rollen-Zeiger: die DoD-Bewertung liegt beim
  **Verifier** (Modul 11) · Vorgänger dieses Punktes:
  `docs/reviews/review-slice-079.md` F-8 (INFO, geschlossen)
- `pfad`: `docs/plan/planning/in-progress/slice-079-coverage-scope-schnitt.md:93-97`
- `befund`: „das Profil führt real **1679** Statements als Nenner
  (dedupliziert über die Block-Position — die Stage druckt nur die
  `total:`-Prozentzeile)“ — die Zahl ist exakt reproduziert (eigene Probe:
  Profil-Summe 1679, `go tool cover -func` 69,9 %), es wird **nicht**
  überzogen. Die Klammer untertreibt: die Stage druckt die gesamte
  `-func`-Ausgabe (213 Funktionszeilen) plus die `coverage-gate:`-Meldung;
  gemeint und zutreffend ist, dass **keine** Zeile die Statement-Zahl nennt.
- `verifizierbar`: ja — Log der `coverage`-Stufe
- `klasse`: „Klammer untertreibt die Ausgabemenge (benannt)“

## Negativbefunde

- geprüft, ohne Befund: erster Report F-1 **(b)** — die Erklärung des
  Log-Artefakts trägt: der Stage-Lauf weist `cmd/pg-change-feed` real als
  `<pkg>\t\tcoverage: 0.0% of statements` aus, nicht als `[no test files]`
  (Log-Zeile der Stufe; dieselbe Form trägt `grpc/streamv1`) — das ist die
  Aufrufform mit `-coverpkg`, und die Probe, die den Fehler in F-1 sichtbar
  macht, ist mit dieser Erklärung im Dokument selbst angelegt.
- geprüft, ohne Befund: erster Report F-1 **(c)** — kein lebender Träger führt
  die alte Behauptung mehr (repo-weiter `grep` über `*.md` außerhalb
  `.harness/baseline/`). Im eingefrorenen Closure-Artefakt
  `docs/plan/planning/done/slice-049-test-coverage-gate.md:133` steht der alte
  Wortlaut weiter („Kein Paket ganz ohne Testdatei mit Fachlogik: die drei
  …“) — `done/`-Bestand, kein Träger dieses Slice, nicht korrigiert; benannt
  für die Closure.
- geprüft, ohne Befund: erster Report F-2 — die zwei Grenzen stehen jetzt als
  nummerierte Punkte 4 und 5 in der §Grenze, in der Form, die die Sektion sonst
  führt (nummeriert, fett gesetzte Wächter-Aussage), mit zutreffender
  Wächter-Benennung: Punkt 4 erklärt die Prozent-Schwelle ausdrücklich zum
  alleinigen Wächter und sagt, dass die Gegenstands-Hälfte der Fitness
  Function damit unvollständig getragen ist; Punkt 5 nennt „keiner“ als
  Wächter. Die Verortung trägt, und die Begründung gegen eine ADR-Umdeutung
  ist sachlich richtig: `ADR-0071`s Re-Evaluierungs-Trigger (a) beschreibt
  Kommen oder Gehen eines Pakets; ihn auf die *Rücknahme* eines bereits
  ausgenommenen umzudeuten, wäre eine Änderung einer `Accepted`-ADR
  (`AGENTS.md` §3.5) für einen Punkt, der als Sensor-Grenze vollständig
  beschreibbar ist. Der Vorbehalt dieses Laufs zu denselben Zeilen betrifft
  allein die Zahlen-Basis (F-2).
- geprüft, ohne Befund: erster Report F-3 — geschlossen durch `0aa4d37`: der
  doppelte „Keine Mindestzahl“-Block ist entfernt (jetzt 1×), der
  Vorlagen-Satz „Pflicht, sobald mindestens eine berührte Sub-Area BF oder
  Hybrid ist“ ist weg (0×), und §8 endet mit der repo-eigenen Antwort
  („Alle berührten Sub-Areas GF …“ samt den vier Pflichtkriterien, 1×). Kein
  Abschnitts- oder Inhaltsverlust: alle acht Sektionen vorhanden, sieben
  gesetzte und sechs offene DoD-Zeilen wie zuvor.
- geprüft, ohne Befund: der neue `welle-20`-Zeiger in der Sensor-Doku
  (`:61-62`) — ein Adress-Zeiger in der Kennungs-Form, die dieselbe Datei
  bereits führt (`seit slice-049`, `roadmap.md` §Nächste Wellen), inhaltlich
  gedeckt (`welle-20` §1 trägt die Test-Arbeit an dem, was ungedeckt bleibt),
  kein Host-Pfad, kein Link auf einen Lifecycle-Pfad.
- geprüft, ohne Befund: neue Zeilen — kein Host-Präfix, keine Slice-/Wellen-
  *Chronik*; `git diff a34f7c6..5d3ac6d` umfasst genau zwei Dateien, die
  Betreffe nennen beide `ADR-0071`, keine `SPEC-*`/`ARC-*`, kein
  Attributions-Trailer.
- geprüft, ohne Befund: Abnahme unverändert — `make coverage-gate` Exit 0
  (69,90 %) und `make gates` Exit 0 (dort 69,70 %); beide Werte liegen über 65,
  die Fixrunde berührt weder Rezept noch Profil (zwei Doku-Dateien).
- geprüft, ohne Befund: erster Report F-4 bis F-7 (INFO) — unverändert,
  ohne Rückgabe-Pfeil; sie gehen in die Closure.
- geprüft, ohne Befund: §2 Paarungs-Zeile (`:129`) — die Aussage „im Repo
  **mit** Wellen von der nächsten Welle-Closure“ ist für diesen Stand
  zutreffend: `welle-20` ist offen (Datei flach unter
  `docs/plan/planning/`), `slice-079` nennt sie als seine Welle, und Modul 6
  lässt die Welle-Closure alles seit der letzten Welle Geschlossene prüfen —
  auch Slices ohne Wellen-Zugehörigkeit. Kein Finding.

## Was ausdrücklich trägt

- **F-1 (a)** — die Kernaussage ist wahr: `cmd/pg-change-feed` liegt in
  `-coverpkg` und in der Testpaket-Liste, hat keine Testdatei und trägt 49
  Statements, alle mit `count = 0`; es steht im Nenner der 1679.
- **F-1 (b)** — die Artefakt-Erklärung ist gemessen richtig (s. Negativbefunde).
- **F-1 (c)** — kein lebender Träger der alten Aussage; der `done/`-Bestand
  ist benannt, nicht stillschweigend mitgezogen.
- **F-2** — Verortung, Form und die Begründung gegen die ADR-Umdeutung tragen.
- **F-8** — der DoD-Wortlaut beschreibt jetzt, was der Lauf hergibt; die
  Überziehung ist weg (F-3 dieses Reports nennt die eine verbleibende
  Ungenauigkeit, ohne Wirkung auf die Aussage).
- **F-3** — geschlossen, ohne Nebenverlust (Plan-Integrität nachgezählt).
- **Abnahme** — grün auf beiden Wegen, dieselbe Schwelle 65, dieselben 1679
  Statements; der Schnitt selbst ist unberührt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Messgegenstands-Aussage übertrifft die reale
Paketlage“ (**2. Vorkommen**, Vorgänger `docs/reviews/review-slice-079.md`
F-1) · „Rechnung nennt eine Basis, die die eigene Ist-Stand-Zeile nicht trägt“
(verwandt `docs/reviews/review-slice-076.md` F-1) · „Klammer untertreibt die
Ausgabemenge (benannt)“

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Kein Finding bestreitet eine
Regel, ein Verdikt oder die Abnahme; kein Konflikt-Pfad (Modul 8).

**Übergabe:** F-1 und F-2 nehmen die **Rückkante Reviewer → Implementer** —
beide liegen in `harness/sensors/coverage-gate.md`, der Datei, die diese
Fixrunde geändert hat, und beide ändern, was die §Grenze über den Gegenstand
und seine Zahlen behauptet; der Umfang ist zwei Stellen. F-1 ist das
**zweite** Vorkommen derselben Klasse in derselben Datei — für den
Steering-Loop-Zähler daher als Wiederholung zu notieren (die Klasse hat die
3×-Schwelle noch nicht erreicht). F-3 geht **ohne Rückkante** in die Closure.
Die Klassen dieses Reports und die unveränderten Klassen des ersten Reports
(F-4 bis F-7) sind der Übergabepunkt in den Register-Zähler (§7).

**DoD-Checkbox-Nachzug:** **nein** — F-1/F-2 tragen einen
Reviewer → Implementer-Pfeil, also greift
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde nicht; die
Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor“ bleibt in
§2 des Slice-Plans offen und wird regulär bei Schritt 21 des
Implementer-Workflows nachgezogen. Nimmt der Implementer F-1/F-2 an oder
begründet sie (isolierte LOWs, Modul 8), greift der Nachzug in dem Lauf, der
den DoD-Nachtrag ohnehin fährt.

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität, die
§6-Risiko-Ausgänge und die drei Paarungen prüft der Verifier bzw. die
Planner-Closure separat (Modul 11/6/5).
