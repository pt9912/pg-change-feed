# BEO-PGC/dod-kriterium-haengt-am-messhost

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form von
DoD-Kriterien und Closure-Triggern eines Slice-Plans, keine eigene Sub-Area im
Sinn der Modus-Deklaration).

Die Beobachtung: Ein DoD-Kriterium oder ein Closure-Trigger verlangt das
Ergebnis eines **Ganz-Targets** — hier „ein realer `make bench`-Lauf mit
Exit 0“ —, dessen Ausgang von einer **Eigenschaft des Messhosts** abhängt und
nicht von der Arbeit des Slice. Auf dem Messhost endete `make bench` am
vorhandenen Skript `tools/bench-source-impact.sh` oberhalb der 35-%-Schwelle
(die Ursache liegt laut Architect-Verdikt in der Latenz des Festschreibens auf
dem Datenträger des Hosts); das neue, vom Slice gelieferte Skript lief
einzeln mit Exit 0. Das Kriterium war damit nicht erfüllbar, ohne dass der
Slice etwas falsch gemacht hätte, und der Plan verlangte an zwei Stellen
(DoD-Punkt, Closure-Trigger) ein Ergebnis, das nur ein anderer Host liefert.

**Warum das schwer zu sehen ist:** Beim Schnitt gilt „`make bench` grün“ als
selbstverständliche Vorbedingung — die drei bestehenden Skripte liefen bei
ihrer Kalibrierung grün (die Zeile in `harness/README.md` nennt 25,0 %/28,2 %
ohne Host). Der Host steht in keinem Kriterium; er wird erst sichtbar, wenn ein
zweiter Host misst.

**Abgrenzung zu benachbarten Einträgen.**
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` beschreibt eine Zahl, die
gegen ihre Messung driftet; hier stimmt jede Zahl, aber ein **Kriterium** hängt
an einer Größe, die der Slice nicht bestimmt. `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
beschreibt einen genannten Beleg, der seinen Satz nicht trägt; hier trägt der
Beleg seinen Satz (das Skript ist rot), der Satz gehört nur nicht zum Slice.

**Die Antwort ist eine Formulierung, kein Werkzeug:** ein Kriterium nennt den
**Beleg des Gegenstands** (das gelieferte Skript einzeln, Exit 0, gedruckte
Zeilen mit Lauf) und führt das Ganz-Target-Ergebnis als **zur Kenntnis** mit
seinem Messhost, nicht als Bedingung.
