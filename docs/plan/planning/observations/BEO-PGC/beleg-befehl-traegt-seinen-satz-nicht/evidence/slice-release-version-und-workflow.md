# Beleg: slice-release-version-und-workflow

Vorgang: `slice-release-version-und-workflow` — zwei DoD-Checkboxen im
Slice-Plan (`docs/plan/planning/in-progress/release-version-und-workflow.md`)
verwiesen für ihre reale Verifikation (Zwei-lokale-Registries-Digest-
Vergleich; YAML-/`bash -n`-Syntaxprüfung) explizit auf „siehe §7".

Fund: §7 „Closure-Notiz" trug im selben Commit (`5e3c7b68`) weiterhin den
unausgefüllten Vorlagen-Platzhalter `*(wird bei Bearbeitung gefüllt.)*` —
eine **Adresse**, die auf ein zum Zeitpunkt des Verweises noch gar nicht
existierendes Ziel zeigte (dieselbe Unterklasse wie der vierte Beleg
dieser Beobachtung, `slice-093`s Verweis auf einen nicht existierenden
„Lauf-Bericht"). Die behaupteten Prüfungen selbst waren real durchgeführt
worden — nur ihr genannter Beleg-Anker war leer. Vom unabhängigen
Reviewer gefunden (`review-slice-release-version-und-workflow.md` F-1,
HIGH), nicht vom Implementer im selben Lauf.

Fix: die „siehe §7"-Verweise entfernt, die konkreten Belege direkt in die
Checkbox-Texte verschoben (kein Verweis auf eine erst später zu füllende
Sektion mehr).

Quelle: `docs/reviews/review-slice-release-version-und-workflow.md` <!-- d-check:status-provenance -->
F-1 · Fixrunden-Commit `31cf6e8a`.
