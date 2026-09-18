# Beleg: slice-090

Vorgang: `slice-090` — das Sync-Gate des generierten Protobuf-Codes.

Fund: Im Delta-Review zu `slice-090` **D-2** behauptete das neu
angelegte Sensor-Dokument, die Zeilennummer des Befunds sei „der des
**Kontext**-Diffs **eine Zeile voraus**", und der Nachbarsatz nannte den
Abstand „systematisch ein bis drei Zeilen". Gemessen (Einfügungs- **und**
Änderungs-Zweig, Lagen 1, 2, 3, 4, 5, 40, 226): der Abstand ist **0** in
Zeile 1, **1** in Zeile 2, **2** in Zeile 3, **3** ab Zeile 4. Die Aussage ist
damit an zwei der sieben Lagen falsch und an keiner allgemein.

**Besonderheit dieses Vorkommens: die Zahl stand in einer Korrektur.** Der Satz
wurde geschrieben, um F-2 des Reviews zu `slice-090` zu beheben (die Zeilenangabe war
der Hunk-Anfang statt der Abweichung) — er ersetzte eine falsche Angabe durch
eine andere. Die Klasse hat hier also ihren **eigenen** Behebungs-Vorgang
getroffen, nicht einen Altbestand: der `diff -U0`-Fix war richtig, der
erklärende Satz darüber war es nicht.

**Ein Vorgang, eine Zählung.** D-1 und D-2 desselben Reports sind zwei Klassen
und je ein Vorkommen; innerhalb desselben Vorgangs zählt jede Klasse einmal
(Modul 6). Die Zahl `74,80 %`/`74,70 %` aus der Verifikation zu `slice-090` V-1 ist **kein**
Vorkommen dieser Klasse — dort ist die *Messung* nicht reproduzierbar, nicht die
Aussage falsch; sie steht in
`BEO-PGC/test-integration-retention-timing-flake`.

Quelle: Delta-Review zu `slice-090` (D-2) ·
`harness/sensors/generated-sync.md` §Ausgabe und Ausgänge (der Satz über den
Abstand, berichtigt in `81f1fff`).
