# Welle 6 — Consumer-Zugriffsweg — Closure-Notiz

**Welle:** welle-6
**Abschluss:** 2026-09-12
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-CON-001.a`](../../../../spec/pflichtenheft.md) (Registrierung:
  Zugriffsweg offen) und [`LH-FA-CON-004.a`](../../../../spec/pflichtenheft.md)
  (Bestätigung: Zugriffsweg offen, Vorwärts-Invariante nicht durchgesetzt)
  sind erfüllt: ein externer Consumer kann sich über den neuen
  CLI-Unterbefehl `register-consumer` registrieren (`slice-021`) und über
  `acknowledge-consumer` eine Position bestätigen (`slice-022`), ohne die
  CDC-Speichertabellen direkt zu schreiben; die Vorwärts-Invariante
  (`ErrPositionRegression`) ist damit für jeden externen Aufruf real
  durchgesetzt.
- Der volle Rundlauf (registrieren → lesen → bestätigen → Neustart
  simulieren → an der bestätigten Position fortsetzen), ausschließlich über
  den neuen Zugriffsweg, ist real gegen PostgreSQL belegt
  (`internal/bootstrap/welle6_endtoend_test.go::TestWelle6ConsumerFullCycleEndToEnd`)
  — das war das *Mehr* dieser Welle gegenüber den einzelnen Slice-DoDs und
  wurde von keinem der beiden Slices allein geliefert (siehe „Was ging
  anders als geplant?").

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der Architect-Verdikt-Weg (Verdikt statt neues ADR, wenn bereits
  Accepted-ADRs den Optionsraum abdecken) trug zwei Slices ohne
  Folge-ADR-Overhead: `slice-021`s Verdikt (`ADR-0019`/`0020`/`0046`
  schließen den Zugriffsweg erschöpfend) wurde für `slice-022`
  kanalgenerisch weiterverwendet, ohne einen zweiten Architect-Rundlauf für
  dieselbe Frage.
- Reviewer und Verifier haben in beiden Slices unabhängig voneinander
  eigene Mutationstests gegen die Vorwärts-Invariante gefahren (nicht nur
  Implementer-Behauptungen übernommen) und jedes Mal real dieselbe rote
  Stelle reproduziert — doppelte, unabhängige Bestätigung.
- Der Ende-zu-Ende-Rundlauf hat genau die Lücke gefunden, für die er als
  *Mehr* benannt war (Welle-Ziel §1): Keine der beiden Slice-DoDs allein
  bewies, dass Registrierung und Bestätigung tatsächlich zusammen einen
  funktionierenden Consumer-Zyklus ergeben.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Der Welle-Plan hatte den vollen Ende-zu-Ende-Rundlauf `slice-022`s
  eigenem DoD zugeschrieben (§3 der Welle); tatsächlich deckte `slice-022`
  nur die Vorwärts-Invariante ab, nicht das Lesen und Fortsetzen nach
  Neustart. Das war eine Planungslücke, keine Umsetzungslücke der beiden
  Slices — Konsequenz: der Rundlauf wurde als eigener, repo-weiter
  Verifikations-Beleg (Modul 8 §Rollen-Sequenz für eine Welle,
  Closure-Schritt 1) direkt vor der Welle-Closure nachgeholt, statt
  stillschweigend als „erledigt" angenommen zu werden.
- Zwei Lastenheft-CRs (`LH-FA-SST-006`, `LH-FA-SST-007`) landeten während
  laufender Slice-Commit-Fenster auf `main`, ohne Bezug zu dieser Welle —
  von Reviewer/Verifier korrekt als fremd erkannt und nicht dem jeweiligen
  Slice angelastet. Konsequenz: `ADR-0020` wurde beim Trigger-Audit dieser
  Welle re-evaluiert (Ausgang: bestätigt, kein Folge-ADR nötig — siehe
  Steering-Loop-Einträge).
- `BEO-PGC/rollen-verdrahtung` erreichte während dieser Welle die
  3×-Schwelle (`slice-021` und `slice-022` reproduzieren dasselbe
  Verdrahtungsmuster wie `slice-011`) — Konsequenz: Folge-Slice `slice-023`
  angelegt (siehe Folge-Slices unten).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Rollen-spezifische DSN-Verdrahtung** — Ausgang **geplant** (keine
  geschärfte Regel, da eine physische Verdrahtungslücke im Bootstrap-Code
  vorliegt, kein Workflow-Disziplin-Defizit; Begründung:
  [`architect-review-welle-6.md`](../../adr/architect-review-welle-6.md) Zug 2)
  — Folge-Slice `slice-023` · seit welle-6.
  Auslöser: `BEO-PGC/rollen-verdrahtung` (slice-011, slice-021,
  slice-022 — 3×).
- **`ADR-0020` re-evaluiert:** Trigger („beobachtbarer API-Consumer-Bedarf")
  ist gegen `LH-FA-SST-006` wörtlich eingetreten — Verdikt: bestätigt,
  kein Folge-ADR (Entscheidungssatz widerspruchsfrei mit `LH-FA-SST-006`;
  künftige Protokollwahl ist eine neue, parallele ADR, keine Supersession).
  Kein Register-Auslöser (ADR-Trigger-Audit, Modul 6 Closure-Schritt 2, kein
  3×-Beobachtungs-Eintrag).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Ein
Eintrag mit neuem Ausgang in dieser Welle: `rollen-verdrahtung` (3×,
**geplant** → `slice-023` — siehe Steering-Loop-Eintrag oben). Unter der
Schwelle, zur Sichtung bei der nächsten Slice-Planung: `adapter-fehler-ausgang`
(1×), `lese-doppelquelle` (2×, in dieser Welle geprüft und nicht berührt),
`schema-rollout-fremdobjekte` (1×), `walsender-wirksamkeit` (1×),
`test-isolation-geteilter-zustand` (1×, neu in dieser Welle). Bereits
verkörpert, unverändert in dieser Welle: `a-check-null-abdeckung` (3×),
`d-migrate-nacharbeit` (4×), `dod-checkbox-nachzug` (3×), `plan-nachzug`
(2×), `plan-vorlagen-defekt` (3×). Bereits eingetreten, unverändert: `cdc-capture-lag-real` (2×),
`health-endpoint-heartbeat` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- `slice-023` (Rollen-spezifische DSN-Verdrahtung) — wellenlos, liegt in
  `open/`.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Beide Slices (`slice-021`, `slice-022`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Ende-zu-Ende-Rundlauf (repo-weiter Verifikations-Beleg, Modul 8
  Closure-Schritt 1): `internal/bootstrap/welle6_endtoend_test.go`,
  real gegen PostgreSQL — Registrieren, Lesen, Bestätigen, Neustart
  simulieren, Fortsetzen ab der bestätigten Position, ausschließlich über
  den externen Zugriffsweg, kein Direktschreiben der CDC-Speichertabellen.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig"/bestätigt
  (`docs/plan/adr/architect-review-welle-6.md`).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die beiden Slice-Dateien, ihre Review-/
Verifier-Reports und Architect-Verdikte sowie dieser Welle-Plan bleiben
vollständig in `done/`. Vor der ersten tatsächlichen Archivierung gilt die
Prüfpflicht aus dem Closure-Command: Geltungsbereich der Sensoren
(`structure`-Regel 5 keilt auf `done/slice-*.md` — bei Stubs in einem
Unterverzeichnis anzupassen) und Link-/ID-Pflichten im Stub prüfen.
