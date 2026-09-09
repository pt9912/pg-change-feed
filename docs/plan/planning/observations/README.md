# Beobachtungs-Register

Steht nur diese `README.md` in der Ablage, ist *nichts beobachtet* die
Aussage — die leere Ablage ist sie, und jedes Repo fängt mit ihr an
(Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register).

Regeln dieser Ablage: Je Beobachtung ein Verzeichnis
`BEO-<KUERZEL>/<slug>/` — `<KUERZEL>` aus der Modus-Deklaration in
`harness/conventions.md` nachgeschlagen, nicht erfunden; `<slug>` lowercase
Kebab-Case. Darin drei Dateien mit drei Lebensdauern: `observation.md`
(unveränderlich ab Anlage), `state.md` (veränderlich: offen oder einer der
drei Ausgänge — verkörpert · geplant · gestrichen — mit auflösbarem Anker),
`evidence/<vorgangs-id>.md` (eine je Auftreten, formgebunden).

- **Geschrieben** wird bei der Slice-Closure: Kennung zitieren (eine
  weitere `evidence/`-Datei) oder neu anlegen — neu formulieren spaltet die
  Klasse. Der Zähler ist die Zahl der `evidence/`-Dateien, er wird nie
  geschrieben.
- **Gelesen** wird an zwei Stellen: die Welle-Closure liest, was 3× erreicht
  hat (Lese-Schritt); die Slice-Planung sichtet, was darunter steht (§8
  des Slice-Plans, `docs/plan/planning/`).
- Gestrichen heißt nicht gelöscht: `state.md` trägt den Ausgang mit
  Begründung.