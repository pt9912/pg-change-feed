# BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Reihenfolge
zwischen einer §4-Rückführungs-Bedingung des Slice-Plans und der Umsetzung,
keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine §4-Bedingung, die die Klärung einer Frage **vor** die
Umsetzung legt, wird erst bei der Closure ausgewertet — geschrieben wird die
Umsetzung vorher. Ist die Bedingung dann eingetreten, steht der Slice zwischen
zwei Ausgängen: der benannte Weg (`in-progress` → `next`, zurück zur Zerlegung)
ist nicht mehr gangbar, weil das Artefakt bereits existiert; der tatsächliche
Weg ist der Konflikt-Pfad. Belegt an `slice-073`: §4 nannte die
Kongruenz-Frage wörtlich vorab („… bevor der Hook geschrieben wird"), die
Bedingung trat ein, der Slice folgte ihr nicht durch einen Neuschnitt, und der
Preis war der volle Konflikt-Pfad (zwei HIGH, zwei Review-Runden, zwei
Folge-ADRs).
