# BEO-PGC/report-nackte-id-ohne-link

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Ausführungsdisziplin jeder Rolle beim Verfassen von Review-/
Verifikations-Reports, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein neu geschriebener Review- oder Verifikations-Report
(`docs/reviews/*.md`) trägt eine nackte `LH-*`-/`ADR-*`-Kennung ohne Link
auf ihre Definition — genau die Regel, die `docs-check`s `ids`-Prüfung
durchsetzt und die jede Rolle explizit angewiesen wird, vor dem eigenen
Commit selbst zu prüfen. Der Fund tritt trotz dieser wiederholten
expliziten Anweisung real auf: einmal drang er real bis zum gepushten
Commit durch (Planner, `slice-054`), zweimal fanden ihn Reviewer/Verifier
selbst — einmal noch vor dem eigenen Commit (`slice-055`), einmal erst
durch den nachfolgenden `make gates`-Lauf des Planners (`slice-056`).
