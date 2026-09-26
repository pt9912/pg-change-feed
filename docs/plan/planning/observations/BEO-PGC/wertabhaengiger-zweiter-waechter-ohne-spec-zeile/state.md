Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Benannte Spec-Lücke des
Lerneintrags von `slice-transformationen-kern-rename` (Adresse: dieser Eintrag). Trigger für die
Auflösung: die Umsetzung von `slice-transformationen-antragsweg-usecase` (K3 am Antrag, Test
„zwei Regeln mit gleichem Zielnamen“) macht die Konstellation im Betrieb unerreichbar; die Wahl
zwischen den zwei Auflösungen der `observation.md` steht an, sobald der Auftraggeber die Spec
zu `SPEC-030` erneut berührt oder ein zweiter Erzeuger von Regelständen den Wächter erreichen
lässt (`slice-transformationen-backfill-pfad` ruft `BuildRowImage` mit dem Regelstand je Zeile).
Stand nach `slice-transformationen-antragsweg-usecase`: K3 hält beim Antrag jeden Zielnamen
von den Spaltennamen der Tabelle und vom Zielnamen jeder anderen Regel fern; der Test „zwei
Regeln mit gleichem Zielnamen endet `failed`“ (Use-Case- und Verdrahtungs-Test) zeigt, dass der
Betrieb den Wächter in `model.BuildRowImage` nicht erreicht — der Trigger dieses Eintrags ist
für den Regelfall eingetreten. Zwei Wege bleiben, die den Wächter erreichen (Review F-7 und
Plan §6 des Slice, hergeleitet aus dem Quelltext, nicht erprobt): (1) ein fehlgeschlagener
Vermerk `applied` nach dem Nachtrag **und** ein kollidierender Folgeantrag im selben Durchlauf —
die laufende Bindung trägt beide Regeln bis zum Neustart, der Erfassungspfad endet in der
Fehlerklasse `schema`; (2) K3 prüft gegen die Spaltenliste zum Antragszeitpunkt, eine spätere
Spalten-Erweiterung lässt den Zielnamen kollidieren (der Fall, den `slice-transformationen-e2e-abhilfe`
real belegt). Die Wahl zwischen den zwei Auflösungen oben berücksichtigt (1): unter Auflösung (1)
prüfte die Anwendbarkeit die Zielnamen der Regelmenge unabhängig vom Wert; ob der
kollidierende Nachtrag dann beim Setzen oder erst beim ersten Bild endet, hängt von der
Ausführung der Auflösung ab (offen).
Zähler (abgeleitet): **1×** (evidence/slice-transformationen-kern-rename.md) — F-7 ist ein Weg zu
dieser Klasse, kein zweites Auftreten.
