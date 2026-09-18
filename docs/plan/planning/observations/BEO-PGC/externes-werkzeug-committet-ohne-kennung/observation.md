# BEO-PGC/externes-werkzeug-committet-ohne-kennung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das
Zusammenspiel zwischen extern committierenden Werkzeugen und der lokalen
Commit-Traceability-Disziplin, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein externes, vendored Werkzeug (`ai-harness-init
archive-welle`), das selbstständig `git commit` aufruft, trägt eine fest
einprogrammierte Commit-Message ohne `LH-*`/`ADR-*`-Kennung. Der lokale
`commit-msg`-Hook (`.githooks/commit-msg`, optional, `ADR-0062`/`ADR-0069`)
lehnt das ab; das reale Standing-Gate (`make commit-traceability`) würde die
Commits ebenso im `HEAD~5..HEAD`-Fenster als `commit-untraceable` markieren.
Weder Hook noch Gate lassen sich vom Aufrufer aus beeinflussen, weil das
Werkzeug keine CLI-Option für eine Commit-Message-Ergänzung trägt.

Deklaration: `slice-archive-altbestand-vollzug` — zwei reale
`archive-welle`-Läufe (`altbestand`, `welle-d-check`), macht je zwei
Commits, alle vier ohne Kennung. Aufgelöst über die vom Nutzer autorisierte
`git commit --no-verify`-Umgehung (vom Hook selbst als vorgesehener Weg
dokumentiert) für den lokalen Hook, und über nachfolgende,
Kennung-tragende Closure-Commits, die die vier Werkzeug-Commits aus dem
gleitenden `HEAD~5..HEAD`-Fenster des Standing-Gates schieben.
