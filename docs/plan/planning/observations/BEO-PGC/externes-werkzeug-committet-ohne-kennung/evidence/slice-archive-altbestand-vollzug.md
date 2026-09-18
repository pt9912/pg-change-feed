# Beleg: slice-archive-altbestand-vollzug

Vorgang: `slice-archive-altbestand-vollzug` — realer `archive-welle`-Lauf
für die Schlüssel `altbestand` und `welle-d-check`.

Fund: Beide Läufe erzeugen je zwei Commits mit fest einprogrammierter
Message („archive-welle: <schlüssel> …"), keine trägt eine `LH-*`/`ADR-*`-
Kennung. Der lokale `commit-msg`-Hook lehnte den allerersten Versuch
(nur der Move-Teil, `f46526e`) ab; dieser wurde per `git revert`
zurückgenommen. Nach Nutzer-Entscheidung („Commit mit --no-verify") liefen
die zwei tatsächlich verbuchten Läufe (`altbestand`, `welle-d-check`) je in
einem Stück über eine prozess-scoped `GIT_CONFIG_COUNT=1
GIT_CONFIG_KEY_0=core.hooksPath GIT_CONFIG_VALUE_0=/dev/null`-Umgebungs-
variable vor dem Werkzeug-Aufruf — sie setzt `core.hooksPath` nur für den
einen Subprozessbaum, ohne die Repo-Konfiguration selbst zu ändern, und
lässt beide interne Commits (Move, Inhalt) je Lauf durchlaufen. Das reale
Standing-Gate (`make commit-traceability`) markierte die vier resultierenden
Commits danach im `HEAD~5..HEAD`-Fenster als `commit-untraceable` —
aufgelöst durch die nachfolgenden, Kennung-tragenden Closure-Commits dieses
Slices und der Welle.

Quelle: `docs/plan/planning/done/slice-archive-altbestand-vollzug.md` §2/§6 ·
Commits `b96e3e7`, `0caee7f`, `a39bbcd`, `a787203`.
