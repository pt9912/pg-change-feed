# Review-Report: slice-079 (dritte Fixrunde) — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität, §6-Risiko-Ausgänge, Register und die drei Paarungen
sind **nicht** Gegenstand dieses Reports (Verifier bzw. Planner-Closure,
Modul 11/6).

**Gegenstand:** `slice-079`, dritte Fixrunde — `65aead2`, Diff
`c1a3353..65aead2` (**genau eine** Datei: `harness/sensors/coverage-gate.md`),
dazu der während dieses Laufs entstandene Planner-Zug `9671b8f`
(Register-Beleg, außerhalb des benannten Diffs, mitgeprüft). Frisch geprüft,
nicht als Bestätigung der Vorläufe; besonderes Augenmerk auf der **neuen**
§Zählbasis-Sektion. Vorgänger: `review-slice-079.md`,
`review-slice-079-fixrunde.md` (eigene Datei je Lauf).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` · **Modell:**
deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext:**

- Fixrunde `65aead2` vollständig, dazu `9671b8f` und der Stand `c1a3353..HEAD`
- Die beiden Vorläufe dieses Slice (0 HIGH / 3 LOW / 5 INFO bzw. 0 HIGH /
  2 LOW / 1 INFO)
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Entscheidung 1–3, §Fitness Function · `AGENTS.md` §3.1, §3.5, §3.7, §3.9, §4
- `harness/mk/coverage.mk` (`THRESHOLD ?= 65`) · `harness/sensors/coverage-gate.md`

---

## Eigene Messungen dieses Laufs (Stand `65aead2`/`9671b8f`)

Alle Zahlen aus eigenen Läufen und eigener Profil-Deduplizierung; Exit-Codes
ungepiped gelesen (`AGENTS.md` §3.9).

| Lauf / Probe | Ergebnis |
|---|---|
| `make coverage-gate` (Stufe 65) | **Exit 0** — `coverage-gate: OK — Coverage 69.90% erfüllt Schwelle 65%`; die Stage druckt `total: (statements) 69.9%` |
| `make gates` | **Exit 0** — coverage-gate OK **69,70 %**; d-check 649 Dateien/0 Befunde; a-check 0; baseline-verify `v6.5.0` OK |
| Dedup-Regel der Datei („gedeckt = mindestens ein Vorkommen `count > 0`“), eigene Implementierung | 1679 Statements / **1173** gedeckt = **69,863 %** → druckt **69,9 %** — deckt sich mit der gedruckten Zeile **desselben** Laufs |
| derselbe Lauf, `go tool cover -func` gegengelesen | 69,9 % — Regel und Werkzeug stimmen überein |
| beschriebener **alter** Fehler (gedeckt nur, wenn das erste Vorkommen `0` war und ein späteres `> 0`) | 1679 / **1169** = 69,625 % — reproduziert die alte Basis 1167 im Rahmen der Lauf-zu-Lauf-Schwankung (2 Statements) |
| Zählung **pro Profilzeile** („jede Block-Position so oft, wie sie vorkommt“) | 1601 / 40431 = **3,96 %** |
| `go tool cover -func` auf Original-Profil vs. summierendem Dedup-Profil | **69,9 % / 69,9 %** — identisch |
| Rollen-Probe der fünf Pakete ohne Testdatei (`go list -f '{{len .TestGoFiles}}'`) | `queries`, `port/inbound`, `domain/errors`: **0 Profilzeilen**; `grpc/streamv1` 86 / 61 = 70,9 %; `cmd/pg-change-feed` 49 / **0** |
| `postgresstorage/mapper` | 15 Statements / 12 gedeckt (80,0 %) |

---

## Findings

### F-1 — Die §Zählbasis erklärt die gedruckte Prozentzeile mit einem Mechanismus, den das Werkzeug nicht hat; die zwei „Größen“ sind eine

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (eine Aussage beschreibt, was
  da ist) · `harness/sensors/coverage-gate.md` §Zählbasis (neu in diesem Diff)
- `pfad`: `harness/sensors/coverage-gate.md:54-58` (Bullet 2) gegen `:32`
  (Einstieg-Zeile, die auf §Zählbasis verweist) und `:52-53` (Bullet 1)
- `befund`: Bullet 2 sagt, `go tool cover` führe die gedruckte Prozentzeile
  „über das gemergte Profil, in dem jede Block-Position so oft zählt, wie sie
  vorkommt“, das Verhältnis sei dadurch dasselbe, die absoluten Zahlen nicht.
  Gemessen ist das nicht so: zählt man pro Profilzeile, ergibt sich
  **1601/40431 = 3,96 %**, nicht 69,9 %; und `go tool cover -func` liefert auf
  dem Original-Profil und auf einem aufsummierten Dedup-Profil **denselben**
  Wert (69,9 %) — die gedruckte Zahl ruht also auf genau der deduplizierten
  Basis (1171/1679), und die Duplikate sind für sie ohne Wirkung. Damit ist die
  gedruckte Zeile **keine eigene Größe** gegenüber Bullet 1: die beiden
  genannten Zahlen `69,70 %` und `69,74 %` sind **eine** Messung, getrennt
  allein durch die Ausgabepräzision (`go tool cover` druckt eine
  Nachkommastelle: 69,74 % → `69.7%`; `tools/coverage-gate.sh` formatiert sie
  mit `%.2f` → `69.70%`). Der Abschnitt verfehlt damit genau die Beziehung, für
  die er angelegt wurde, und die Einstieg-Zeile verweist für „die Größe dieser
  Zeile“ in ihn.
- `verifizierbar`: ja — Profil ziehen, einmal pro Zeile und einmal dedupliziert
  rechnen, `go tool cover -func` auf beide Profile legen
- `klasse`: „Erklärter Mechanismus deckt sich nicht mit dem Werkzeugverhalten“

### F-2 — Der Register-Beleg führt zwei Zählbasen in einem Satz (Planner-Datei, aus `9671b8f`)

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (ein Wert, ein Ort) ·
  **zweites Vorkommen** der Klasse aus `review-slice-079-fixrunde.md` F-2
  (dort in der Sensor-Doku, jetzt geschlossen)
- `pfad`: `docs/plan/planning/observations/BEO-PGC/regel-weiter-als-ihr-sensor/evidence/slice-079.md:7-12`
  gegen `harness/sensors/coverage-gate.md:109-111`
- `befund`: Der Nachzug `9671b8f` hat in diesem Beleg die Formelzeile auf die
  neue Basis gehoben (`(1171 + 2) / (1679 + 23) = 68,92 %`), die beiden
  Rot-Fälle im selben Satz aber auf dem alten Stand gelassen:
  `postgresstorage (52,3 %)`, `replication/receive (64,2 %)`. Die Sensor-Doku
  führt dieselben zwei Fälle jetzt mit **52,51 % / 64,45 %**. Zwei Artefakte
  desselben Slice nennen dieselbe Rückrechnung damit 0,2 Prozentpunkte
  auseinander; welcher Stand gilt, ist dem Beleg allein nicht zu entnehmen.
- `verifizierbar`: ja — beide Dateien gegeneinander lesen
  (1171/1679 = Basis, die drei Ausgänge 68,92 / 52,51 / 64,45)
- `klasse`: „Rechnung nennt eine Basis, die die eigene Ist-Stand-Zeile nicht trägt“

## Negativbefunde

- geprüft, ohne Befund: **§Zählbasis Bullet 1** — die Regel ist die richtige:
  eigene Deduplizierung über die Block-Position mit „gedeckt = mindestens ein
  Vorkommen `count > 0`“ liefert für **denselben** Lauf genau die gedruckte
  Prozentzeile (1679/1173 = 69,863 % → 69,9 %). Der Vorläufer-Fund (Basis 1167
  nicht herleitbar) ist damit **ursächlich aufgeklärt** und geschlossen: die
  alte Basis entsteht, wenn man einen Block nur dann als gedeckt zählt, wenn
  sein **erstes** Vorkommen `0` war — meine Nachbildung dieser Regel ergibt
  1169 und liegt damit im Schwankungsband der 1167. Die Zahl `1171` der Datei
  ist die Kalibrierungs-Messung (`1171/1679 = 69,74 %`, identisch mit
  `ADR-0071`); mein heutiger Lauf liegt mit 1173 zwei Statements darüber, was
  die Datei als Lauf-zu-Lauf-Schwankung selbst benennt (§Grenze Punkt 4).
- geprüft, ohne Befund: **§Zählbasis Bullet 3** — die Herkunft der Zahlen der
  drei Ausgenommenen ist richtig verortet: mein eigener Nachbau des
  Vorher-Laufs (Arbeitsbaum-Kopie bei `35b9a00`, alter Umfang) ergibt
  `2467 → 1679` mit 1215 gedeckt und **44** in den drei Paketen (2/23, 31/610,
  11/155) — dieselben Werte, die Bullet 3 und §Grenze Punkt 4 nennen.
- geprüft, ohne Befund: erster Vorläufer F-1 (Paket-Aufzählung) — geschlossen.
  Die Pakete **ohne Testdatei** im Gegenstand sind genau fünf, und die drei
  Rollen sind erschöpfend (keine Statements / einige Statements ohne eigene
  Testdatei / kein gedecktes Statement): die drei erstgenannten haben real
  **null** Profilzeilen, `grpc/streamv1` 86 Statements mit 61 gedeckt (über
  fremde Testpakete), `cmd/pg-change-feed` 49 mit 0. Der Superlativ steht jetzt
  unter der Lesart, die der Satz selbst trägt („das einzige Paket des
  Gegenstands ohne ein einziges gedecktes Statement“) und ist in dieser Form
  wahr; „die beiden Pakete mit Statements“ ist die richtige Zahl. Ein sechstes
  Paket ohne Testdatei und eine vierte Rolle gibt es nicht.
- geprüft, ohne Befund: **§Grenze Punkt 4** — die Rückrechnung ist jetzt
  herleitbar: Basis und ihre Herkunft („§Zählbasis: 1171 gedeckt von 1679
  Statements im Gegenstand; 2/23, 31/610, 11/155“), als Rechnung ausgewiesen
  („kein eigener Lauf“). Nachgerechnet stimmen alle drei Ausgänge:
  `1173/1702 = 68,92 %` · `1202/2289 = 52,51 %` · `1182/1834 = 64,45 %`; die
  Abstände zur Schwelle (+3,9 / −12,5 / −0,6 Prozentpunkte) sind korrekt
  gerechnet und tragen die Aussage, dass die Schwankung die Ausgänge nicht
  umkehrt (0,6 pp ≈ 11 Statements bei `receive`, die gemessene Schwankung liegt
  bei wenigen Statements).
- geprüft, ohne Befund: die §Zählbasis führt **keinen** beweglichen Wert —
  `THRESHOLD` kommt in ihr nicht vor (er bleibt in `harness/mk/coverage.mk`);
  sie erklärt die Herkunft von Messgrößen, nicht die geltende Stufe. Die
  Figuren, die sie und §Grenze Punkt 4 beide nennen, stehen dort als Definition
  und als Anwendung, und die Anwendung nennt ihre Quelle („§Zählbasis: …“) —
  keine zweite Führung desselben Zustands im Sinne von `AGENTS.md` §3.7. Die
  Drift-Gefahr dieser Doppelung ist real und bereits eingetreten (in der
  Planner-Datei, s. F-2), innerhalb der Sensor-Doku aber nicht.
- geprüft, ohne Befund: `9671b8f` (außerhalb des benannten Diffs) — inhaltlich
  der Nachzug der Zählbasis in den Register-Beleg; die Formelzeile ist auf 1171
  gehoben, die zwei Rot-Werte des Satzes nicht (F-2).
- geprüft, ohne Befund: neuer Text ohne Host-Präfix, ohne Slice-/Wellen-Chronik;
  der `welle-20`-Zeiger in §Grenze Punkt 1 bleibt ein Adress-Zeiger in der
  Kennungs-Form der Datei. `git diff c1a3353..65aead2` umfasst genau eine Datei;
  beide Betreffe nennen `ADR-0071`, keine `SPEC-*`/`ARC-*`, kein
  Attributions-Trailer.
- geprüft, ohne Befund: **Abnahme unverändert** — `make coverage-gate` Exit 0
  (69,90 %) und `make gates` Exit 0 (dort 69,70 %); beide über 65. Der Diff ist
  rein dokumentarisch, Rezept und Profil sind unberührt.

## Was ausdrücklich trägt

- **Die Zählbasis-Idee trägt** — eine Stelle, die einmal erklärt, woher jede
  Zahl kommt, ist die richtige Antwort auf die zwei Vorläufe; Bullet 1 (Regel
  und Größe) und Bullet 3 (Herkunft der Vorher-Zahlen) sind gemessen richtig,
  und die Datei reproduziert `ADR-0071`s `1171/1679 = 69,74 %` mit eigener
  Deduplizierung. Zu korrigieren ist Bullet 2 (F-1) — und zwar, weil die
  gedruckte Zeile **dieselbe** Größe in anderer Ausgabepräzision ist, nicht
  eine eigene.
- **Die Paket-Rollen tragen** (Vorläufer F-1 geschlossen), **die Rückrechnung
  trägt** (Vorläufer F-2 in der Sensor-Doku geschlossen), **die Abnahme
  trägt** (beide Läufe grün, Exit ungepiped gelesen).
- Der Register-Beleg ist inhaltlich richtig aufgebaut; offen ist dort allein
  die Zahlen-Konsistenz (F-2) — Planner-Datei, kein Rückweg zum Implementer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Erklärter Mechanismus deckt sich nicht mit
dem Werkzeugverhalten“ (neu) · „Rechnung nennt eine Basis, die die eigene
Ist-Stand-Zeile nicht trägt“ (**2. Vorkommen**, Vorgänger
`review-slice-079-fixrunde.md` F-2; dort die Sensor-Doku, jetzt geschlossen)

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Kein Finding bestreitet eine
Regel, ein Verdikt, eine Zahl der operativen Aussagen (Schwelle, Rampe, Nenner,
Rückrechnung) oder die Abnahme.

**Übergabe:** F-1 nimmt die **Rückkante Reviewer → Implementer** — ein
Bullet in `harness/sensors/coverage-gate.md`; inhaltlich ist es eine
Korrektur der Erklärung, keine Zahlenänderung. F-2 nimmt die **Rückkante
Review → Planner** (`observations/` ist seine Datei). **Zur Proportion, wie
erbeten:** F-1 ist ein Satz ohne Wirkung auf Gate, Schwelle oder Zahlen — ich
halte einen vierten, kleinen Nachzug für vertretbar, aber **nicht** für
blockierend; wird er stattdessen als benannte Kleinigkeit in die Closure
genommen, ist das ein zulässiger Ausgang (die Klasse ist neu, nicht die dritte
Wiederholung derselben). F-2 ist reine Zahlen-Konsistenz in einer Planner-Datei
und läuft mit dessen Closure ohnehin durch die Hand.

**DoD-Checkbox-Nachzug:** **nein** — F-1 trägt einen Reviewer → Implementer-
Pfeil, also greift `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde nicht; die Zeile „Review durchgeführt, Report unter `docs/reviews/`
liegt vor“ bleibt in §2 des Slice-Plans offen. Entscheidest du, F-1 in die
Closure zu nehmen statt in eine vierte Runde, fährt der Lauf, der ihn tilgt,
den Häkchen-Nachzug mit.

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-Konformität, die
§6-Risiko-Ausgänge und die drei Paarungen prüft der Verifier bzw. die
Planner-Closure separat (Modul 11/6/5).
