# Beleg: slice-084

Vorgang: `slice-084` — die Naht in `postgresack`.

Fund: Die zwei Sensor-Dokumente trugen nach dem Zug **vier** Zahlen, die gegen
die Messung drifteten — gefunden vom Reviewer als F-1 (LOW), nachgemessen vom
Implementer:

| Träger | stand dort | gemessen |
|---|---|---|
| `harness/sensors/db-adapter-coverage.md` (§Zählbasis, §Kalibrierung) | DB-Nenner 650 · `postgresack` 23 · 477 gedeckt · 73,38 % | **659 · 32 · 491 · 74,51 %** |
| `harness/sensors/coverage-gate.md` (§Zählbasis) | Unit-Nenner 1831 | **1903** |
| `harness/sensors/coverage-gate.md` (§Grenze Punkt 4) | `2/23` | **30/32** |
| `harness/sensors/coverage-gate.md` (§Grenze Punkt 1) | `mapper` 15/12 | **20/16** |

**Zwei davon sind nicht diesem Zug zuzurechnen** (`mapper` wuchs durch einen
anderen Vorgang; die `1831`-Drift stammt aus Vorgängen davor) — der Implementer
hat sie trotzdem mitgezogen und das im Commit eigens benannt, weil er sonst in
derselben Datei eine bekannt-falsche Gegenwartsaussage stehen ließe.

**Die Form-Entscheidung des Fixrunden-Laufs ist der eigentliche Ertrag:**
Die **gedeckte** Zahl ist lauf-gebunden (derselbe Stand, derselbe Befehl maß
`71.90 %` und `72.00 %`), der **Nenner** ist es nicht. Also steht der Nenner als
**Zustand** und die gedeckte Zahl als **Beleg eines konkreten Laufs** — sie
nennt ihren Lauf und nie „der Ist-Stand".

**Und die ehrliche Grenze:** für Zahlen in Doku-Trägern gibt es **keinen
Sensor**. „Welche Änderung macht das rot?" hat hier kein Objekt — die
verfügbare Falsifikation ist die **Messung** selbst, und sie hat vier Werte im
Baum widerlegt. Das ist die Klasse dieses Eintrags in einem Satz.

Quelle: `docs/reviews/review-slice-084.md` (F-1) · Implementer-Bericht der
Fixrunde (Commit `db42231`, mit der Form-Begründung) ·
`harness/sensors/db-adapter-coverage.md`, `harness/sensors/coverage-gate.md`.
