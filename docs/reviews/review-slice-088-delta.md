# Review-Report: slice-088 — Delta-Nachlauf zur Fixrunde — 2026-09-16

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** die Fixrunde **`d522266`** (nur
`internal/adapters/driving/replication/mapper/mapper_test.go` und die eine
DoD-Zeile des Slice-Plans) und der Plan-Nachzug **`11ba45b`**
(F-2/F-3 des Erstlaufs). Beide als **frisches Artefakt** gelesen, nicht als
Bestätigung des Erstlaufs; der Erstlauf-Report `review-slice-088` wurde nur als
Verweis benutzt. Kein neuer Vollauf: die drei übrigen Testdateien und die
Produktionsseite sind unberührt.

**Produktionsseite:** leer. `git diff b477f0f..HEAD -- '*.go' ':!*_test.go'`
endet ohne Ausgabe (Exit 0); über beide Commits hinweg ändern sich nur zwei
Dateien (`--stat`: Slice-Plan + `mapper_test.go`). Der Nenner **1903** kann
damit nicht gewandert sein; die Gate-Zahl des Erstlaufs gilt unverändert.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` ·
**Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext:**

- Erstlauf-Report `docs/reviews/review-slice-088.md` (0 HIGH, F-1 MEDIUM,
  F-2/F-3 MEDIUM, F-4 INFO)
- Slice-Plan `slice-088` §1/§2/§3/§8 im Stand `eb68126`
- `ADR-0082`, `ADR-0080`, `ADR-0029`, `ADR-0059` Teilfrage 3;
  `LH-FA-CFG-001`, `LH-FA-CFG-005`
- `AGENTS.md` §3.7, §3.9; Beobachtungs-Register `BEO-PGC/` (Zähler-Stände)
- Form-Vorlage des Nachlaufs: `docs/reviews/review-slice-077-delta.md`

---

## Findings

### D-1 — Der berichtigte Absatz schließt mit einer Chronik-Zeile

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (*Ein Kommentar beschreibt, was da ist* — die
  Regel für Zustandsfelder und Prosa: „die Historie [hält] `git`") ·
  Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-088-coverage-tail-uebersetzung.md:66-68`
- `befund`: Der Absatz endet mit „**Die Zahl „fünf" in einer früheren Fassung
  dieses Absatzes war selbst wieder eine ungezählte Übernahme** — sie nannte
  sechs Namen und traf acht; sie steht hier als berichtigte Zahl." Der Satz
  nennt keinen Ist-Zustand der Sache, sondern den Zustand eines **früheren
  Textes** (was er nannte, was er traf, dass er berichtigt wurde). Der Satz
  direkt darüber trägt dagegen eine Funktion: er begründet, *warum* die Spalte
  „Summe" als Gegenprobe existiert und warum die Zuordnung nachgemessen ist —
  das ist die geforderte Begründung des §1. Der neue Schlusssatz trägt keine
  vergleichbare Funktion; kein anderes der ~80 Slice-Pläne dieses Repos führt
  eine solche Vorher/Nachher-Zeile (repo-weit nur `slice-088`).
- `verifizierbar`: ja — `grep` über `docs/plan/planning/` auf dieselbe Form;
  die Form-Regel ist Prosa, also kein Gate
- `klasse`: „Vorher/Nachher-Sprache in Doku-Prosa" (im Register noch **nicht**
  geführt; Erstauftreten dieses Laufs)

### D-2 — LP1 beansprucht `JSONImage`, das zu LP3 gehört (Bestand, nicht aus der Fixrunde)

- `kategorie`: INFO
- `quelle`: Maintainability · `ADR-0082` §Was daraus für die Slices folgt
  (Cluster-Tabelle)
- `pfad`: `docs/plan/planning/in-progress/slice-088-coverage-tail-uebersetzung.md:112-113`
- `befund`: Liefer-Punkt 1 trägt die Träger `replication/decode` und
  `replication/mapper`; sein zweiter Aufzählungspunkt sagt von `JSONImage`,
  es sei „**dabei** der erste Fall" — `JSONImage` liegt aber in
  `postgresstorage/mapper` und gehört damit zu Liefer-Punkt 3, der ihn auch
  selbst führt (`:121`). Der Satz steht unverändert seit der Plan-Anlage; die
  Fixrunde hat die Liste unmittelbar darüber korrigiert (`eb68126`) und ihn
  damit als einzigen Träger-Verweis des Abschnitts stehen gelassen, der eine
  Funktion dem falschen Liefer-Punkt zuordnet — dasselbe Muster, dessen
  Berichtigung Anlass des Nachzugs war. Keine Wirkung auf die Lieferung: die
  drei Liefer-Punkte sind vollständig und die Funktion ist gedeckt.
- `verifizierbar`: ja — Paket-Zugehörigkeit von `JSONImage` im Baum gegen die
  Träger-Nennung des Liefer-Punkts
- `klasse`: „Funktion am falschen Liefer-Punkt geführt"

### D-3 — Die Quelle des F-3-Fehlleses steht weiter im Register

- `kategorie`: INFO
- `quelle`: Maintainability · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register („Der Zähler wird
  abgeleitet, nicht geführt")
- `pfad`: `docs/plan/planning/observations/BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code/state.md:1`
  gegen `:8` derselben Datei
- `befund`: Zeile 1 trägt „Zustand: offen **(1×)** — unter der Schwelle …",
  Zeile 8 derselben Datei „Zähler (abgeleitet): **2×** (evidence/slice-084.md,
  evidence/slice-085.md)". Zwei Zahlen für denselben Zähler in **einer** Datei;
  die Kopfzeile ist die Stelle, aus der §8 des Plans den Stand „1×" gelesen
  hat (Erstlauf F-3). Der Plan ist berichtigt, die Fundstelle nicht — und der
  Lese-Schritt der `welle-20`-Closure liest genau diese Datei. **Außerhalb
  dieses Deltas** (die Fixrunde berührt `observations/` nicht) und daher hier
  als Hinweis, nicht als Mangel des geprüften Artefakts.
- `verifizierbar`: ja — `ls` der `evidence/`-Dateien gegen beide Zahlen der
  Datei
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`

## Negativbefunde

- geprüft, **ohne Befund: F-1 ist behoben** — die zwei No-op-Zusagen sind an
  ihre Eingabeseite gebunden. Eigene Reproduktion (Arbeitskopie, gepinntes
  Toolchain-Image, `go test -count=1` über das Paket; Baseline des fixierten
  Stands Exit **0**):

  | Mutation (Eingabeseite) | Exit | Aussage |
  |---|---|---|
  | M4 — `if !activated { return nil }` → `_ = activated` | **1** | genau `TestConsumeRelationOnUnboundTableStaysNoop` (Relation mit einer Spalte über der bekannten Form endet sichtbar) |
  | M5 — `if !found { return nil }` → `_ = found` | **1** | genau `TestConsumeRelationWithoutRegisteredVersionStaysNoop` (`TableSchema [""]`, `RegisterVersion [{tbl-1-v1 tbl-1 1}]`) |
  | M6 — Store **vor** dem Riegel befragt, ohne Fehler | **1** | **nur** die Aufruf-Zusage (`Store-Zugriffe außerhalb der Aktivierungen: CurrentVersion [""], TableSchema ["sv-1"], RegisterVersion []`) |
  | M7 — `CurrentVersion` mit `QualifiedName()` statt `TableID` | **1** | unter 7 Fehlschlägen auch der neue Test (`CurrentVersion-Aufrufe = ["public.feed"], wollen genau die gebundene Tabellen-Kennung tbl-1`) |

  M6 belegt, dass die **Aufruf**-Zusage trägt, M7 die **Argument**-Bindung —
  wie im Fixrunden-Bericht angegeben. Einschränkung zu M7: es ersetzt den
  einzigen Aufrufpunkt und färbt deshalb sechs **vor** diesem Slice
  bestehende Tests mit rot; als Beleg für die *neue* Bindung trägt es nur über
  die Fehlermeldung des neuen Tests, nicht als saubere Einzelpunkt-Sonde.
- geprüft, ohne Befund: **die Fidelity-Angabe des Implementers** — die Riegel
  sind als `_ = activated`/`_ = found` entfernt, nicht gelöscht; beide Formen
  kompilieren, der Exit 1 stammt nachweislich aus den Testaussagen
  (`mapper_test.go:1086` bzw. `:1139`), nicht vom Compiler.
- geprüft, ohne Befund: **der Stub als Träger** — die drei Aufzeichnungen
  hängen nicht am Fehlerpfad; `-race` über das Paket endet Exit **0** (die
  beiden nebenläufigen Tests des Pakets nutzen den anderen Fake,
  `stubSchemaStore` wird ausschließlich einläufig verwendet).
- geprüft, ohne Befund: **die Plan-Berichtigungen F-2/F-3** — `11ba45b` nennt
  jetzt **acht** Paketwechsel und zählt genau die acht gemessenen Funktionen
  auf (`observeRelation`, `oldTupleValues`, `tupleValues`, `JSONImage`,
  `IncludeColumn`, `setSchemaVersion`, `removeExcluded`, `qualifiedNames`);
  die Summen-Gegenprobe 14/25/17/4 = 60 steht daneben; die Zähler nennen den
  Stand **3× / 2× / 3×** und sind gegen `evidence/` nachgezählt; das
  §8-Ergebnis („zwei stehen bereits darüber") stimmt.
- geprüft, ohne Befund: **die DoD-Zeile** — `[x]` mit Report-Pfad, nachgezogen
  im Fixrunden-Commit (regulärer Pfad `implement-slice` Schritt 21, nicht der
  Nachzug des Reviewers). Nur diese eine Zeile; die übrigen Closure-Items
  bleiben offen.
- geprüft, ohne Befund: **F-4 (INFO)** — unformuliert gelassen. Als INFO
  verlangt es keine Aktion; die drei Zahlen der §3-Tabelle bleiben allerdings
  ohne Einheit und gehen unter beiden Lesarten nicht auf (Statements: +17/+13/+18;
  neue Testfunktionen je Datei: 17/9/16). Der Hinweis bleibt damit **offen**,
  nicht erledigt.
- geprüft, ohne Befund: **Hygiene des Deltas** — zwei Commits, keine
  Lauf-Artefakte, keine Trailer, keine `SPEC-*`/`ARC-*` in den Betreffs,
  Historie linear; die Änderung an `mapper_test.go` bleibt test-only.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Nachlaufs:** „Vorher/Nachher-Sprache in Doku-Prosa"
(D-1, Erstauftreten) · „Funktion am falschen Liefer-Punkt geführt" (D-2) ·
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (D-3).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. F-1 ist mit eigener Messung als
behoben bestätigt: die beiden Mutationen, die den Erstlauf noch grün ließen,
färben jetzt genau je einen der beiden Tests rot, und zwei darüber
hinausgehende Mutationen belegen Aufruf- und Argument-Bindung. Die
Plan-Defekte F-2/F-3 sind berichtigt; die Nachprüfung findet die acht
Paketwechsel und die drei Zähler-Stände.

**Übergabe:** **kein Rückgabe-Pfeil.** D-1 kann ohne zweiten Implementer-Lauf
gehen — es ist Form, keine Aussage; D-2/D-3 sind Hinweise an Planner bzw.
Wellen-Closure. Die Klassen gehen mit der Slice-Closure §7 in den Zähler.

**DoD-Häkchen „Review durchgeführt":** bleibt auf `[x]` — bereits im
Fixrunden-Commit nachgezogen; dieser Nachlauf ändert den Plan **nicht**.

**Aus Sicht dieses Nachlaufs ist der Slice verifikationsreif.** Offen bleiben
allein die Verifier-Strecke (DoD/Spec), die §6-Risiko-Ausgänge, das
Beobachtungs-Register und die drei Paarungen — Planner-Closure, Modul 11/5.
