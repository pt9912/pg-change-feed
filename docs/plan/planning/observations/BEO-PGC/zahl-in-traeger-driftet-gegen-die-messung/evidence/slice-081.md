# Beleg: slice-081

Vorgang: `slice-081` — die Executor-Naht und die Rampen-Neu-Bemessung.

Fund: **Zwei** Träger nannten eine Zahl, die gegen die Messung driftete — beide
im selben Vorgang, an verschiedenen Orten.

1. `tools/harness/db-coverage.sh:18` — „der Lauf ueber postgresstorage allein
   traegt **610** Statements". Gemessen sind **472** (Summe der `numStmts` über
   die `postgresstorage/<datei>.go`-Positionen in `store.coverprofile`).
   Dieselbe Aussage stand in `harness/sensors/db-adapter-coverage.md:43`
   **bereits richtig** — zwei Träger derselben Zahl widersprachen sich im Baum.
   Gefunden als Review-Finding **F-1 (HIGH)**; berichtigt mit dem gemessenen
   Wert.
2. `docs/plan/adr/0071-…md` und `docs/plan/planning/welle-20.md` nennen weiter
   **1679** Statements für die netzlos prüfbare Fläche (in `ADR-0071` zusätzlich
   2467, 788 und 69,74 %) — nach dem Subjekt-Transfer trägt sie **1831**.
   Lebende Träger **außerhalb** des Slice-Diffs; gefunden als Review-Finding
   **F-6 (INFO)**, ausdrücklich als „Architect-/Planner-Frage, kein
   Slice-081-Defekt" eingeordnet und **nicht** in diesem Slice berichtigt.

Beide Fundstellen zeigen dieselbe Ursache: der Träger nennt eine Zahl, aber
nicht ihren **Ursprung** — ob sie gemessen, übernommen oder abgeleitet ist. Die
Zahl `610` war eine Messung, die der Gegenstand überholt hat; die `1679` war
eine Messung, die der Transfer überholt hat.

Quelle: `docs/reviews/review-slice-081.md` (F-1, F-6, mit eigenen Messungen des
Reviewers) · `harness/sensors/db-adapter-coverage.md:43` (der richtige Träger
derselben Zahl) · `internal/adapters/driven/postgresstorage/sqlexec/` (der
Gegenstand, der sich bewegt hat).
