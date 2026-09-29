# Review-Report: slice-code-kommentare-bereinigung B2 (T8) — 2026-09-29

**Review-Art:** Code — geprüft gegen den Slice-Plan
[`slice-code-kommentare-bereinigung`](../plan/planning/in-progress/slice-code-kommentare-bereinigung.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und die
Hard Rules `AGENTS.md` §3.7/§3.12/§3.13/§3.9 (Modul 10). Keine
DoD-Verifikation — das ist die Aufgabe des Verifiers (Modul 11).

**Gegenstand:** Diff `89d347db..d6accc38` — die fünf T8-Teil-Commits
(`d1fc13da` driven, `56486d09` bootstrap+cmd, `c550329e` domain+application,
`c7e9517d` driving, `793bb71b` test/integration samt Restmenge), der
Suchlauf-Nachzug `532093c8` und der Erzeugnis-Commit `d6accc38`. Vorlage
dieses Reviews: die Finding-Klasse „Kürzung lässt hängende
Herkunfts-Referenz zurück" aus dem B1-Report
(`review-slice-code-kommentare-bereinigung-b1.md`, 1 HIGH + 3 MEDIUM,
behoben in `89d347db`).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1` (2026-09-28) ·
**Modell:** GLM (Claude-Agent-SDK, Typ `reviewer`) · **Datum:** 2026-09-29

**Eingangs-Kontext:**

- Slice-Plan @ `d6accc38` (inkl. Suchlauf-Feld am Stand `793bb71b` und
  Liefer-Punkt-2-Haken)
- `AGENTS.md` §3.7 (Kommentar-Klassen, ein Anker), §3.12 (Zahlen-Ursprung),
  §3.13 (Suchlauf), §3.9 (Exit-Disziplin), §3.11 (host-lokaler Pfad)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- Generator-Vertrag des Erzeugnisses:
  `TestAbdeckungstabelleZeilen`/`abdeckungsZeile` in
  `test/integration/integration_test.go`
- `harness/sensors/kommentar-kennungen.md`,
  `harness/sensors/suchlauf-nachmessen.md`,
  `harness/sensors/db-adapter-coverage.md` (Nenner-Behauptung Teil 1)
- `harness/conventions.md` (MR-000) · Baseline `v6.13.0` ·
  `regelwerk/modul-10-review-harness.md`

---

## Findings

### F-1 — Kürzung lässt das hängende Fragment „(Konvergenz)" zurück

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7; Plan §1 „Regeln je Kandidat"; Skill-HIGH
  „Kommentar trägt keine der Kommentar-Klassen" (Teilersetzung lässt den
  Rest stehen); dieselbe Klasse wie F-1/F-2/F-3 des B1-Reports
- `pfad`: `test/integration/integration_test.go:1008`
- `befund`: Der Inline-Kommentar trug „(Konvergenz mit
  `LH-FA-SCH-004`, bereits von `ADR-0059` Teilfrage 4 akzeptiert)"; die
  Kürzung entfernte beide Referenzen und ließ „(Konvergenz)" stehen —
  ein Wort, dessen Bezug (mit was? wogegen?) nicht mehr auflösbar ist.
  Der Satz davor trägt die Substanz („denselben ErrIncompatibleSchemaChange-Pfad
  aus wie jede andere inkompatible Relation-Änderung"); das Fragment ist
  der Rest der Teilersetzung. `make kommentar-kennungen` meldet die
  Stelle nicht (0 Kandidaten in T8 außerhalb der 14 Restblöcke,
  nachgemessen) — kein Gate fängt das.
- `verifizierbar`: nein — Lese-Handlung; das Werkzeug zählt Kennungen,
  nicht hängende Fragmente
- `klasse`: „Kürzung lässt hängende Herkunfts-Referenz zurück"

### F-2 — „die Konvergenz" ohne ihren Partner

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7; Plan §1 „Regeln je Kandidat"
- `pfad`: `internal/adapters/driving/replication/mapper/mapper_test.go:647`
- `befund`: Der Block trug „die Konvergenz mit `LH-FA-SCH-003`
  (`ADR-0059` Teilfrage 4)"; die Kürzung entfernte den Partner und ließ
  „trägt die Konvergenz (`ADR-0059` Teilfrage 4)" stehen. Der Anker
  löst auf (Festlegung 4 von `ADR-0059` trägt die Konvergenz-Frage,
  nachgelesen), und der Satztext trägt die Substanz — milde Form
  derselben Klasse wie F-1.
- `verifizierbar`: nein — Lese-Handlung
- `klasse`: „Kürzung lässt hängende Herkunfts-Referenz zurück"

### F-3 — „der korrigierte Ursprungstext" ohne seinen Korrektur-Anker

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7; Plan §1 „Regeln je Kandidat"
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:187`
- `befund`: Die Rot-färbende-Mutation-Liste trug „(der von
  `ADR-0048` korrigierte Ursprungstext)"; die Kürzung ließ „(der
  korrigierte Ursprungstext)" stehen — das Attribut „korrigiert" hat
  seinen Bezug verloren (dieselbe Datei behält `ADR-0048` an anderer
  Stelle im Diff korrekt als Block-Anker; nur hier bleibt es hängen).
  Milde Form derselben Klasse.
- `verifizierbar`: nein — Lese-Handlung
- `klasse`: „Kürzung lässt hängende Herkunfts-Referenz zurück"

### F-4 — Die Markierung der Restmenge grenzt breiter als die Maschinen-Wirklichkeit

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Grenze/Abgrenzung trägt, was die Stelle
  leistet); Plan §1 „begründete, benannte Restmenge mit Adresse";
  Generator-Vertrag (`abdeckungsAdressiert`)
- `pfad`: `test/integration/integration_test.go:202` (Stellvertreter; die
  Formulierung steht an allen 14 Restblöcken)
- `befund`: Der Schluss-Absatz „die Kennungen dieses Blocks trägt die
  E2E-Abdeckungstabelle je Testzeile" stimmt pauschal nicht für den
  ADR-/ARC-Anteil der Mengen: der Generator liest nur
  `LH-*`/`SPEC-*` in die Kennungsspalte
  (`abdeckungsAdressiert`); seine Design-Kommentierung sagt dasselbe
  („die ADR- und Sicht-Kennungen … fallen deshalb aus der
  Beschreibungsspalte heraus"). Neun der 14 Blöcke tragen ADR-Kennungen
  (z. B. `ADR-0016`, `ADR-0028`, `ADR-0063`/`ADR-0058`/`ADR-0059`) —
  für sie würde Kürzen auf einen Anker **keine** Abdeckungszeile
  löschen, und ein künftiger Editor liest die Grenze als unveränderten
  Schutz der vollen Menge. Die Begründung der Restmenge (Erhalt der
  Abdeckung) trägt so für einen Teil der Menge nicht.
- `verifizierbar`: nein — Lese-Handlung gegen den Generator; die Zahl
  14 ist nachgemessen, die Teilmenge der ADR-geführten Blöcke aus der
  Werkzeug-Ausgabe abgezählt
- `klasse`: „unpräzise Grenz-Formulierung der Restmenge"

### F-5 — Die Restmenge ist ein Fall für die §3.7-Konkretisierung; die Meldung liegt noch nirgends

- `kategorie`: INFO
- `quelle`: Plan §2 Liefer-Punkt 3 („ein Kandidat mit zwei Ankern, den der
  Implementer als konform begründet, ist ein Fall für die Regel — Meldung
  an den Planner"); `.harness/skills/reviewer.md` §Pflege
- `pfad`: `docs/plan/planning/in-progress/slice-code-kommentare-bereinigung.md:125`
- `befund`: Die 14 gekennzeichneten Blöcke sind genau der Fall, den der
  Plan als Meldung an den Planner (Konkretisierung von `AGENTS.md`
  §3.7) vorsieht: die Kennungs-Menge ist konform begründet
  (Erzeugnis-Eingabe), nicht als Ausnahme stillgestellt. Die Meldung
  selbst ist noch nirgends als Übergabe-Artefakt sichtbar — der Plan
  verortet sie bei der Closure (§7). Hinweis an den Planner, keine
  Aktion an diesem Diff.
- `verifizierbar`: nein — Lese-Handlung; bei Closure prüfen
- `klasse`: „Restmengen-Begründung braucht §3.7-Konkretisierung"

## Negativbefunde

- geprüft, ohne Befund: **Diff-Form aller fünf Teil-Commits** — jede
  +Zeile und jede −Zeile in `*.go` beginnt nach Leerraum mit `//`
  (mechanisch gezählt je Commit: 0 Abweichungen in 68 Dateien); keine
  Endkommentar-Zeile; die Nicht-Go-Dateien sind Plan (Suchlauf-Feld,
  DoD-Haken) und Erzeugnis (`docs/user/e2e-abdeckung.md`).
- geprüft, ohne Befund: **Zahlen-Ursprung je Teil-Commit** — die fünf
  Messages tragen Befehl und beide Stände (68→0 am Stand `89d347db`,
  34→0 am `d1fc13da`, 34→0 am `56486d09`, 31→0 am `c550329e`, 21→14 am
  `c7e9517d`; Summe vor = 188). Nachmessen: `make kommentar-kennungen
  COUNT=1 TESTS=exclude` → 0 (Exit 0); Gesamtauflistung → genau 14
  Kandidaten, alle `test/integration`, alle Godocs von `func TestE2E*`.
- geprüft, ohne Befund: **Suchlauf-Feld** — `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-code-kommentare-bereinigung.md`:
  8/8 Zeilen OK, Exit 0 (4 Zeilen am Stand `7b70b34a`, 4 Zeilen am Stand
  `793bb71b`: 1732/147/1/18). Die Erklärung zur `ff.`=1-Zeile
  (Stringliteral im Tabellentestfall
  `tools/harness/kommentar-kennungen/main_test.go`) ist zutreffend.
- geprüft, ohne Befund: **Restmengen-Behandlung gegen den
  Generator-Vertrag** — nachgelesen: `abdeckungsZeile` bricht ab, wenn
  das Godoc einer `func TestE2E*` keine Spec-Kennung trägt, und speist
  die Kennungsspalte aus `funktion.Doc`; Kürzen auf einen Anker würde
  Abdeckungszeilen löschen. Die 14 Marker stehen genau an den 14
  Werkzeug-Kandidaten (13 + 1 gezählt), der Marker-Text landet nicht in
  der Kurzbeschreibung (der Generator liest nur den ersten Absatz,
  `abdeckungsKommentarAbsatz`), und der Generator sonder-case-t den
  Marker nicht — kein Hidden-Coupling.
- geprüft, ohne Befund: **Marker gegen §3.7 und gegen das Werkzeug** —
  der Marker trägt keine Kennung (die §3.7-Herkunfts-Regel ist nicht
  verletzt), ist der Klasse Grenze/Abgrenzung zuzuordnen (er schreibt an
  den, der die Stelle ändert), und `make kommentar-kennungen` zählt die
  14 Blöcke unverändert — keine Ausnahmeliste, kein Marker-Feature im
  Werkzeug (Plan-Risiko „Die Restmenge wird zur Ausnahmeliste": nicht
  eingetreten).
- geprüft, ohne Befund: **Chronik-/Vorher-Nachher-Sprache in +Zeilen** —
  Suche über alle +Zeilen des Diffs; ein Treffer („eine frühere Position
  endet über die Invariante", `consumerstate_test.go`) ist
  Zustandsbeschreibung der Regel, keine Chronik.
- geprüft, ohne Befund: **Erzeugnis-Diff** —
  `docs/user/e2e-abdeckung.md` ändert in 15 Zeilen ausschließlich die
  Ort-Spalte (`backfill_e2e_test.go:270→274`,
  `integration_test.go:202→205` … `1235→1285`); Kennungs- und
  Beschreibungsspalten sind byte-gleich. `make test-integration` ist
  nicht selbst gefahren — der Lauf wird als gemeldet übernommen, der
  Erzeugnis-Diff ist eigenständig nachgemessen.
- geprüft, ohne Befund: **Stichproben gegen den Code/ADRs je Teil**
  (Zahl der Stichproben: 22 Blöcke/Stellen) — Teil 1: `snapshot_test.go`
  (Nenner-Behauptung gegen `harness/sensors/db-adapter-coverage.md`
  Z. 29 wahr; `ADR-0115` Festlegung 4 trägt das Zitat), `roles_test.go`,
  `administrationrequest_test.go`, `backfilladmission_test.go`,
  `consumerstate_test.go`; Teil 2: `changestream_internal_test.go`
  (`ADR-0100` Teilfrage 5 trägt die Aktivierungs-Frage, nachgelesen),
  `roles_rollout_file_internal_test.go`, `roles_wiring_test.go`,
  `config_file_rest_internal_test.go`; Teil 3: `backfill/service_test.go`
  (`ADR-0113` Festlegung 3 trägt „Kopierdauer = `finished_at −
  started_at`", Z. 281-282 nachgelesen), `capture/service_test.go`
  (`ADR-0060` Teilfrage 2), `change_test.go`, `validation_test.go`;
  Teil 4: `mapper_test.go`, `seam_test.go`, `stream_test.go`,
  `http/server_test.go`, `decode_test.go`; Teil 5: Generator-Vertrag,
  Marker-Zählung, Erzeugnis-Diff. Drei Fundstellen (F-1 bis F-3) — alle
  übrigen Kürzungen lassen einen vollständigen Satz und den tragenden
  Anker stehen.
- geprüft, ohne Befund: **Sensors** — `make test` (Exit 0, direkt
  gesichert), `make fmt-check` (295 Go-Dateien, alle formatiert, Exit 0),
  `make a-check` (0 Befunde, Exit 0), `make kommentar-kennungen`
  (Gesamt- und Exclude-Lauf, s. o.), `make suchlauf-nachmessen`
  (8/8, s. o.).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Kürzung lässt hängende
Herkunfts-Referenz zurück" (3×; Bestands-Klasse aus B1) · „unpräzise
Grenz-Formulierung der Restmenge" (1×, erstes Auftreten) ·
„Restmengen-Begründung braucht §3.7-Konkretisierung" (1×, erstes
Auftreten)

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) zieht eine Fixrunde am
Implementer nach sich; F-2/F-3 (LOW) sind in derselben Runde je nach
Implementer-Entscheid mitzunehmen (dieselbe Klasse, je Stelle ein Wort
bzw. ein Anker). Der Gegenstand ist schmal: je Stelle den Satz
vollständig machen oder das Fragment streichen (das Muster liefern die
sauber gekürzten Blöcke desselben Diffs). Die Fixrunde berührt
Kommentare in `test/integration` — die Lokatoren des Erzeugnisses
verschieben sich, d. h. der Erzeugnis-Lauf (`make test-integration`)
und `make suchlauf-nachmessen` sind nach der Fixrunde erneut zu fahren
(Liefer-Punkt-3-Feld, §2). Die DoD-Zeile „Review durchgeführt" bleibt
offen (Fixrunde läuft regulär über Schritt 21 des Implementer-Ablaufs).

**Übergabe:** Findings gehen an den Implementer (Fixrunde T8); die
Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler („Kürzung lässt hängende Herkunfts-Referenz zurück" — 3× in
diesem Lauf, 6× über B1+B2; „unpräzise Grenz-Formulierung der
Restmenge" — 1×, erstes Auftreten; „Restmengen-Begründung braucht
§3.7-Konkretisierung" — 1×, erstes Auftreten). F-5 geht an den Planner
(Konkretisierung `AGENTS.md` §3.7 bei der Closure). Dieser Report ist
ein Lauf-Beleg.

**Verbleibende Risiken (nicht Findings):**

- `harness/sensors/coverage-gate.md` trägt die in B1 gemeldeten
  gedrifteten Lokatoren (`main.go:44/62/86/101` → 50/68/92/107) weiterhin
  bis zum Planner-Nachzug (Frist: Closure dieses Slice; kein Gate liest
  die Stelle). Die Meldung aus B1 ist offen.
- Der Suchlauf-Stand `793bb71b` driftet mit der Fixrunde (1732
  Kennungs-Zeilen); die vier Zeilen am Stand `793bb71b` sind dann am
  neuen Stand neu zu messen und das Feld nachzutragen.
- F-4 zeigt: die Grenz-Formulierung der Restmenge schützt mehr, als die
  Maschine liest — künftige Kürzungen an den 14 Blöcken können den
  ADR-Anteil ohne Abdeckungsverlust prüfen, sobald die §3.7-Konkretisierung
  (F-5) die Frage entschieden hat.
