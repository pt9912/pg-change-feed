**Vorgang:** slice-antragsqueue-lesefehler-failed (Review F-11, INFO; Verifikation §4 und §6)

**Fund:** Zwei Stellen, keine mit Wirkung auf das Repo.

- Der Bericht des Implementers nannte einen versehentlichen `python3`-Aufruf mit leerem Heredoc; der Reviewer fand im Diff keine Spur eines Host-Interpreters (`git diff … | grep -nE 'sed -i|perl -pi|awk -i'` ohne Treffer, `make fmt-check` Exit 0) und nannte es Kenntnis ohne Verstoß am Erzeugnis (F-11).
- Der Verifier führte 23 Mutationsläufe mit einem Host-`python3`-Skript im Scratchpad aus, das ein Textstück ersetzt; sein Bericht nennt als Ziel zugleich „die Arbeitskopie“ mit `git checkout -- <Datei>` als Rücknahme (§4) und „Kopien der Arbeitsdateien“ (§5). Nach jeder Mutation und am Ende stand `git status --short` leer. Die Regel in `AGENTS.md` §3.1 verlangt für die Mutationsprobe eine Kopie im Scratchpad und nennt einen Host-Interpreter auf der Datei verboten; ob das Skript eine Repo-Datei oder eine Kopie schrieb, lässt sich dem Bericht nicht entnehmen.
- Der Planner der Closure-Sitzung dieses Vorgangs rief `python3 -` mit leerem Heredoc auf, ohne Wirkung auf das Repo (Angabe des Auftraggebers im Auftrag zur Closure von `slice-harness-guard-inplace-textwerkzeug`; übernommen, kein Report als Anker, nicht nachgemessen).

**Form (Ausprägung):** derselbe Fehlgriff wie in den drei bisherigen Vorgängen, diesmal in **drei Rollen** (Implementer, Verifier, Planner) und mit einem Host-Interpreter statt `sed -i`; der Reviewer desselben Vorgangs arbeitete auf einer Kopie im Scratchpad (`git archive`). Der Guard `.claude/hooks/pretooluse-command-guard.sh` liest Sprach-Interpreter nicht; die Durchsetzung bleibt das Review, die Adresse des Guard-Ausbaus ist `slice-harness-guard-inplace-textwerkzeug`.

Quelle: `docs/reviews/review-slice-antragsqueue-lesefehler-failed.md` (F-11) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-antragsqueue-lesefehler-failed.md` (§4 Kopf, §5, §6). <!-- d-check:status-provenance -->
