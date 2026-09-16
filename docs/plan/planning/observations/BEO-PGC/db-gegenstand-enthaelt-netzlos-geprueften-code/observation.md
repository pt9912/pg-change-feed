# BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den
**Messgegenstand** der DB-Adapter-Coverage, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Der Gegenstand der **DB-Adapter-Coverage** (`ADR-0071` Punkt 3)
enthält Code, der **ohne jede PostgreSQL-Verbindung** gedeckt wird. Seit
`slice-084` steht `postgresack` bei **32 von 32** gedeckten Statements — und
**30** davon deckt ein netzloser Test. Die Zahl dieses Pakets ist damit **nicht
mehr von einer echten Verbindung zu unterscheiden**: sie steigt, ohne dass
irgendetwas real geprüft wurde.

Das ist das **Spiegelbild** zu `slice-081`. Dort hat ein **Transfer** den
DB-Nenner *gedräniert* (138 Statements wanderten in den Unit-Gegenstand, die
Quote fiel 75,25 % → 73,38 %) — und das wurde **entschieden** (`ADR-0077`,
teilweise abgelöst durch `ADR-0078`: Subjekt-Transfer mit dreiteiligem
Nachweis). Hier geht die Bewegung in die **andere** Richtung: die Naht bringt
neuen, netzlos geprüften Code **in** das Paket, der DB-Nenner **wächst** (650 →
659) und die Quote **steigt** (73,38 % → 74,51 %) — ohne dass ein einziger
neuer DB-gestützter Beleg entstanden wäre.

**Gemeinsam ist beiden:** die Zahl der DB-Adapter-Coverage entfernt sich von
dem, was sie messen soll — einmal nach unten, einmal nach oben. `ADR-0078` hält
für die Abwärtsbewegung einen Riegel bereit (fällt der Nenner ohne Ankunft,
steht die Schwelle); für die **Aufwärtsbewegung** gibt es keinen.

**Warum das zählt:** Eine Messung, deren Zahl sich ohne den geprüften Gegenstand
bewegt, ist als Verdikt nur noch begrenzt lesbar — und ihre Rampe kann in beide
Richtungen irreführen. `ADR-0080` hat für dieses Paket eine **Verdünnung**
ausdrücklich als Trigger benannt; sie ist mit `slice-084` **eingetreten**.
