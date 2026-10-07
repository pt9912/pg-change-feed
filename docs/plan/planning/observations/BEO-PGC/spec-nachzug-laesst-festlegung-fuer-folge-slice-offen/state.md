Zustand: offen — **3× erreicht**, Ausgang noch nicht zugewiesen: der Lese-Schritt
(wellenlos, ausgelöst von der Closure von `slice-spec-festlegungen-harness-werkzeuge`)
braucht einen Architect-Zug und steht aus. Zu prüfen ist, ob die
Anschlussfähigkeits-Frage des Verifiers („kann ein Implementer den ersten Folge-Slice aus der
Spec allein umsetzen — welche Festlegung fehlt je Schicht?") als Prüfpunkt in den
Verifier-Auftrag für Spec-Nachzug-Slices gehört, oder als Pflicht-Zeile in den Plan eines
Spec-Nachzugs (Festlegungs-Tabelle mit Bestands-Beleg je Zeile, Vorbild:
`slice-transformationen-spec-nachzug` §3).

Zweites Auftreten (`slice-routing-lesewege`): die fehlende Festlegung war eine
Fehlerrangfolge am Lesepfad — gewinnt ein ungültiges Ziel gegen einen Fehler des
Lese-Kontrakts? Prüffrage, die den Fund beim Schreiben der Zeile eines neuen
Filter-Parameters fängt: „was gewinnt bei gleichzeitigem Kontraktfehler?“ Der
Trigger für eine Folge-ADR zur Rangfolge steht im Plan des Slice (der nächste
Lese-Weg mit einem Filter-Parameter, der ein Alphabet prüft).

Drittes Auftreten (`slice-spec-festlegungen-harness-werkzeuge`): die fehlende Festlegung war
die Folge eines Ausgangs (Exit 2 einer Werkzeug-Festlegung: besteht die Korrektur?) beim
Zusammenziehen von vier `Accepted`-ADRs und einem Vertrag in `SPEC-038`; die Prüffrage, die
es fing, war das Gegenlesen jedes Satzes der Festlegung gegen die Quell-ADR, die die
Herkunftstabelle von `ADR-0162` nennt. Für die fünf Folge-Slices der Festlegungen
(`slice-spec-festlegungen-doku-gates`, `-kennungs-gates`, `-code-gates`, `-coverage-gates`,
`-pruefer-hooks`) gilt dieselbe Form.

Zähler (abgeleitet): 3× (evidence/slice-transformationen-spec-nachzug.md,
evidence/slice-routing-lesewege.md, evidence/slice-spec-festlegungen-harness-werkzeuge.md).
