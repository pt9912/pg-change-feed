# Beleg: slice-maintainer-ordner-releasing-verschieben

Vorgang: `slice-maintainer-ordner-releasing-verschieben` — zweites Auftreten der
Klasse, befundet als V-3 im Verifikations-Report
`verifikation-slice-maintainer-ordner-releasing-verschieben.md` (Übergabe an den
Planner mit der Frage an den Architect).

Fund: Die Zitat-Korrektur an
[`ADR-0123`](../../../../../adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
(Pfadwechsel von `releasing.md`) änderte drei Verweise: Linkziel und Linktext
je Stelle. Gemessen mit `grep -n '^## '` auf der Datei: Zeile 48 liegt in
§Kontext (Zeilen 39–96), Zeile 139 in §Entscheidung (97–236), Zeile 304 in
§Konsequenzen (283–315). Zwei der drei Stellen sitzen damit in Abschnitten, die
[`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Entscheidung 1 als nicht korrigierbar listet (§Entscheidung, §Konsequenzen);
nach der Kurzform („das Gerüst darf sich ändern, die Aussage nie; der Referent
bleibt derselbe") ist die Änderung zulässig: Linkziel und Linktext bezeichnen
dieselbe Datei an ihrem neuen Ort, der umgebende Satz (Abschnitt „Beschreibung",
Tabelle „Benötigte Secrets", Secret-Tabelle) steht unverändert. Beleg-Form
erfüllt: Commit-Messages nennen
[`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(`36b311b4`, `1c90030f`, `df503ad5`), genau eine §Geschichte-Zeile in der ADR.

Unterschied zum Erstauftreten (`evidence/slice-harness-baseline-v6-13-0.md`):
dort änderte sich ein Versions-Segment bei gleichem Referenten, hier der **Ort**
des Referenten (Datei verschoben). Die Spanne zwischen Abschnitte-Liste und
Kurzform bleibt in beiden Formen dieselbe: die Liste nennt §Entscheidung und
§Konsequenzen als unberührbar, die Kurzform lässt dort ein Gerüst-Korrektur zu.

Behandlung (Planner, kein Architect-Zug): die Lesart der Kurzform trägt wie im
Audit §4.2/§7 des Erstauftretens; keine Änderung an
[`ADR-0123`](../../../../../adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
über das Verweisgerüst hinaus. Ob eine Folge-ADR die Abschnitte-Liste von
[`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Entscheidung 1 an die Kurzform angleicht (Gerüst-Verweise, auch auf einen
verschobenen Referenten, sind in allen Abschnitten zulässig; die Aussage bleibt
unberührbar), ist eine Entscheidung des Architect — **benannt, nicht angelegt**.
Der Zähler steht damit bei 2×; ein drittes Auftreten ist der übliche Anlass,
den Architect-Zug zu beauftragen.
