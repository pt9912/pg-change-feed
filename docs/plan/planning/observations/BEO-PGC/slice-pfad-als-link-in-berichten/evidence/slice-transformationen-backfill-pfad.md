**Vorgang:** slice-transformationen-backfill-pfad

**Fund:** Der Übergang `open` → `next` brach einen Inbound-Link: `state.md` von
`BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (Register-Datei ohne Überschrift, von den
zwei `structure`-Regeln nicht gedeckt) adressierte den Slice-Plan mit festem Lifecycle-Verzeichnis
(`open/`). `make docs-check` lief rot (`target-missing`); behoben durch Umstellung auf die Kennung
(Commit `13aaf752`, im Implementer-Lauf).

**Form (Ausprägung):** die in `.d-check.yml` benannte Grenze (3) — `state.md` und `evidence/*.md`
tragen keinen Abschnitts-Anker, dort bleibt Disziplin — ist eingetreten; der Fehler stammte aus einer
Register-Datei des Planners, gefunden hat ihn das Gate (`links`) am Move, vor dem Merge. Zwei weitere
Register-Dateien tragen denselben Linktyp als Anker eines Ausgangs *geplant* (`state.md` von
`BEO-PGC/adapter-fehler-ausgang` und `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung`, gemessen mit
`git grep -n -E '\]\([^)]*/(open|next|in-progress)/slice-' -- docs/plan/planning/observations`: drei
Treffer, einer davon die Illustration in einer `observation.md`).

Quelle: `docs/reviews/verifikation-slice-transformationen-backfill-pfad.md` (§8 V-6) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-transformationen-backfill-pfad.md` (F-7). <!-- d-check:status-provenance -->
