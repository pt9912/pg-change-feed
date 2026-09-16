# Beleg: slice-092

Vorgang: `slice-092` — Coverage Cluster D1 (Anwendungs-Kern).

Fund: **Zwei vorbestehende** Tests banden ihre Ablehnung nicht an die Eingabe.
`TestExcludeColumnRejectsMissingSourceColumn` und
`TestIncludeColumnRejectsMissingSourceColumn`
(`internal/application/usecase/{exclude,include}column/service_test.go`) standen
gegen einen Fake, der `exists = false` **unabhängig von der Abfrage** lieferte:
die Suite blieb grün, wenn der Use Case die **falsche** Spalte prüfte (gemessen
an der Parent-Fassung: die Mutation „Use Case prüft `secret` statt der
Kommando-Spalte" → grün; nach der Bindung → rot).

**Gefunden hat sie die Mutationsprobe des Implementers**, der sie — als Teil
seines eigenen Zuschnitts (`slice-092` schreibt die Spalten-Use-Cases an) — auf
die Spaltenadresse gebunden hat. Die Klasse hat damit einen Ort erreicht, an dem
sie **nicht** vermutet wurde: nicht in einem neuen Test, sondern in einem
bestandenen, der wie seine Nachbarn aussah.

**Ursprung ≠ Vorkommen.** Eingeführt wurden die zwei Tests von einem früheren
Vorgang als diesem (die Spaltenausschluss-Use-Cases entstanden mit
`LH-FA-CFG-005`); das ist der **Ursprung** der Form. **Gefunden** wurden sie in
`slice-092`s Implementer-Lauf und in der Verifikation `verify-slice-092` (V-2
P-7, dort am Parent gegengeprüft); das ist das Vorkommen, das dieser Beleg
dateit.

**Warum das zählt:** Der Test ist grün, sieht aus wie seine fünf Geschwister und
trägt eine echte Zusage. Nur die Mutation trennt ihn von ihnen — und die Zusage
„diese Spalte wird abgelehnt" hatte bis dahin keinen Träger.

Quelle: `docs/reviews/verify-slice-092.md` (P-7, am Parent gegengeprüft) ·
`docs/reviews/review-slice-092.md` · Commit `3de9547` (die Bindung) ·
`internal/application/usecase/excludecolumn/service_test.go`,
`internal/application/usecase/includecolumn/service_test.go`.
