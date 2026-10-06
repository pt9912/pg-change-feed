# BEO-PGC/werkzeugvertrag-zusage-ohne-testfall

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft neue
Harness-Werkzeuge mit Vertrag unter `harness/targets/` oder `harness/sensors/`
und ihren Tabellentest, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein neues Werkzeug kommt mit Vertrag und Tabellentest, und
der Test ist grün. Einzelne Zusagen des Vertrags oder Regeln des Codes haben
aber **keinen Fall**: eine Mutation an der Stelle, die die Zusage trägt, lässt
den ganzen Test grün. Sichtbar wird das nur durch Mutieren, nicht durch Lesen
des Tests.

**Abgrenzung.** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` beschreibt
einen **vorhandenen** Test, der an seiner Eingabeseite nicht rot werden kann;
hier **fehlt** der Fall (Reviewer-Skill, MEDIUM „fehlende Negativtests bei
neuem öffentlichem Vertrag“). Die Kategorie entscheidet über die Fixrunde.

Konkret (Erstauftreten): `make test-zitat-vergleich` — die Einzugsregel des
Fence, vier Zusagen des Vertrags §Einheit und das `10#` vor den
Lokator-Zahlen hatten je keinen Fall.
