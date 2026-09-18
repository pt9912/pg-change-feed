# Beleg: slice-archive-altbestand-vollzug

Vorgang: `slice-archive-altbestand-vollzug` — realer `archive-welle`-Lauf
für die Schlüssel `altbestand` und `welle-d-check`.

Fund: Beide Läufe erzeugen je zwei Commits mit fest einprogrammierter
Message („archive-welle: <schlüssel> …"), keine trägt eine `LH-*`/`ADR-*`-
Kennung. Der lokale `commit-msg`-Hook lehnte den ersten Versuch ab; nach
Nutzer-Entscheidung wurde `git commit --no-verify` für die interne
Commit-Logik des Werkzeugs verwendet (über eine prozess-scoped
`GIT_CONFIG_*`-Umgebungsvariable, die `core.hooksPath` nur für den einen
Werkzeug-Subprozess auf einen leeren Pfad setzt, ohne die Repo-Konfiguration
selbst zu ändern). Das reale Standing-Gate (`make commit-traceability`)
markierte die vier Commits danach im `HEAD~5..HEAD`-Fenster als
`commit-untraceable` — aufgelöst durch die nachfolgenden, Kennung-tragenden
Closure-Commits dieses Slices und der Welle.

Quelle: `docs/plan/planning/done/slice-archive-altbestand-vollzug.md` §2/§6 ·
Commits `b96e3e7`, `0caee7f`, `a39bbcd`, `a787203`.
