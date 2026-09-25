# Ein Beleg lebt nur als einmalige Messung eines Lesers, kein committeter Wächter trägt ihn

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Dauerhaftigkeit von
Belegen für Sicherheits- und Ende-zu-Ende-Aussagen, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Eine tragende Aussage eines Capture-kritischen Zuges ist am realen System
gemessen, aber nur durch einen Wegwerf-Test oder eine einmalige Mutation des Reviewers; im
Repository liegt kein Test, der sie bei einer Regression rot färbt. Der Slice benennt die
Grenze ehrlich (Plan §6), der Verifier bestätigt sie — und trotzdem bleibt der Beleg ein
Ereignis dieses einen Laufs.

**Abgrenzung.** `BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (1×) trägt dieselbe Form
für den Boundary-Fall einer Metrik-View; hier betrifft sie das Verhalten der Quelle
(PostgreSQL) und die Seite „Schwelle erreicht → Container endet“. Ob beide Einträge eine
Klasse sind, entscheidet der Lese-Schritt der Welle-Closure von `welle-backfill-bestand`;
gezählt wird bis dahin getrennt.
