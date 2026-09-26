**Vorgang:** slice-harness-suchlauf-nachmessen (Review F-1, HIGH; Verifikation §5, Zeile F-1)

**Fund:** Der Zerleger von `tools/harness/suchlauf-nachmessen.sh` reichte jedes Wort einer `suchlauf`-Zeile unverändert an `git grep` weiter. Der Reviewer legte mit einer Plan-Zeile `diff 0 -O'touch <Datei>;true' -e … -- …` die Datei an; das Werkzeug meldete `OK  soll=0 ist=0` und Exit 0. Dasselbe gilt für `--open-files-in-pager=<Kommando>`. Die Zusage „eine Plan-Zeile führt keinen Shell-Code aus“ stand im Plan (§3, Form der Zeile), ohne Bindung an feindliche Eingabe. Die Fixrunde begrenzt Optionen auf eine Allow-List, Pathspec-Magic auf eine geschlossene Liste und lehnt leere Argumente ab (Exit 2); der Verifier sah die Mutation „Optionsprüfung entfernt“ rot (Marker-Datei entsteht).

Quelle: `docs/reviews/review-slice-harness-suchlauf-nachmessen.md` (F-1, Abschnitt „Injektion über Plan-Inhalt“) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-harness-suchlauf-nachmessen.md` (§4 MA, §5 F-1). <!-- d-check:status-provenance -->
