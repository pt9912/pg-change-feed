Zustand: offen (**3×**) — Schwelle erreicht mit `slice-backfill-snapshot-reader`;
Ausgang noch **nicht** zugewiesen, er gehört dem Lese-Schritt der Closure von
`welle-backfill-bestand` (Modul 6). Der dritte Beleg
(`evidence/slice-backfill-snapshot-reader.md`) trägt eine andere Ausprägung: netzlos
prüfbare Logik lag im ausgenommenen Paket `postgressnapshot`, ohne dass ein
netzloser Test die DB-Zahl hob; der Slice löste sie durch das Unterpaket
`snapshotlogic` im Unit-Gegenstand (Nenner 2040 → 2082, 42 von 42 gedeckt,
`ADR-0080` §Kontext Punkt 5 führt Unterpakete ausgenommener Pakete im
Unit-Gegenstand). Der Slice benennt diesen Weg im Plan vorab als Mittel des Risikos;
ob ein Unterpaket als **Regel** für neue DB-Pakete gelten soll, ist Teil der
Entscheidung. Bis dahin: ein Träger
ist **nicht** vorgeschlagen: die naheliegende Antwort wäre eine Neudefinition
des Gegenstands (nur Code, dessen Test eine Verbindung braucht), und das ist
eine **Entscheidung** über eine Messfläche — `ADR-0071` Punkt 3 trägt sie, und
`ADR-0080` hat die Verdünnung als Trigger benannt. Sie gehört als
Architect-Frage behandelt, nicht als Notiz.

Zähler (abgeleitet): **3×** (evidence/slice-084.md, evidence/slice-085.md,
evidence/slice-backfill-snapshot-reader.md) — Schwelle erreicht. **Ein zweiter Slug wäre die verbotene Umformulierung:** es
ist **eine** Beobachtung, nur die Größenordnung hat sich geändert (netzlos
gedeckter Anteil im Gegenstand **8,35 % → 26,69 %**).

**Der Trigger (a) aus [`ADR-0080`](../../../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
ist damit materiell geworden — und der DB-Hochschalt-Trigger liest sich als
fällig** (76,99 % ≥ 75 %). Nach dem Buchstaben ja, nach der Property nein: der
Anstieg ist Verdünnung, kein Ausbau. **Diese Entscheidung gehört dem
Trigger-Audit der Wellen-Closure** (Modul 6 Schritt 2), nicht einem Slice.

**Verwandt, aber entschieden:** die **Dränage**-Hälfte (`slice-081`, ein
Transfer *aus* dem Gegenstand) ist mit `ADR-0077`/`ADR-0078` beantwortet und
deshalb **kein** Beleg dieses Eintrags — sie steht hier als die andere Richtung
derselben Bewegung. **Nicht zu verwechseln** mit
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (dort driftet eine
**geschriebene Zahl** gegen die Messung; hier driftet die **Messung** selbst
gegen ihren Zweck) und nicht mit
`BEO-PGC/endstufe-unter-eigennem-messgegenstand-unerreichbar` (dort ist ein
Ziel unerreichbar, hier ist ein erreichtes Ziel wenig wert).
