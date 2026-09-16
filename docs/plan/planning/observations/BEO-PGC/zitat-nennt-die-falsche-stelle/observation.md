# BEO-PGC/zitat-nennt-die-falsche-stelle

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in der
ein Träger eine fremde Stelle zitiert, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Träger **zitiert eine Stelle eines anderen Dokuments** —
einen Abschnitt, eine Nummer, einen Paragrafen — als Stütze einer Aussage, und
die zitierte Stelle **trägt sie nicht**: Sie sagt etwas anderes, oder sie ist
eine andere Stelle als die gemeinte. Wie bei
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` kann die **Aussage wahr** sein;
falsch ist die **Stütze**. Der Unterschied liegt im Gegenstand: dort trägt ein
**Befehl** seinen Satz nicht, hier ein **Verweis**.

**Der Mechanismus, der die Klasse schwer sichtbar macht:** Die Nummer stammt
aus einer **Zusammenfassung**. Wer einen Bericht liest, findet in seiner
Kopfzeile die Liste der Verträge, gegen die geprüft wurde — samt Kurzform ihrer
Inhalte. Ein solcher Kopf sieht aus wie eine Quelle und ist eine **Wiedergabe**;
wer aus ihm in einen **stehenden** Träger übernimmt, zitiert die Wiedergabe
statt des Originals. Das Original aufzuschlagen ist die billigste Messung, die
es gibt — und genau die, die unterbleibt, weil die Zusammenfassung schon
antwortet.

Belegt an einem abgeschlossenen Vorgang:

- **`slice-090`** (Delta-Review `review-slice-090-delta` D-1): Das in diesem Zug
  neu angelegte `harness/sensors/generated-sync.md` schrieb „(Festlegung 1 = der
  Baum wird nicht geschrieben; **Festlegung 2 = der Befund nennt Datei und
  Zeile**)". `ADR-0084` §Entscheidung führt **vier** Festlegungen: **1** „Der
  Protobuf-Code bekommt ein Gate" — die zwei Bedingungen (*Temp-Verzeichnis*,
  *Datei und Zeile*) sind ihre **Unterpunkte** —, **2** „Die E2E-Abdeckungstabelle
  bekommt heute kein Gate", **3** `plan.yaml`/`down.sql`, **4** `image-hash.txt`.
  Die zitierte Stelle trug die Aussage also nicht. Der Ursprung der falschen
  Nummer ist die **Kopfzeile** von `review-slice-090.md` (dort dieselbe
  Kurzform, ebenso in `verify-slice-090.md`) — ein Bericht ist **Lauf-Beleg**,
  kein Zitat-Träger; übernommen wurde er ungeprüft in einen **stehenden** Träger
  und von dort in einen zweiten.

**Grenze der Klasse — und die Abgrenzung zu den Nachbarn.** Das ADR zitiert sich
in seinen Re-Evaluierungs-Triggern **selbst richtig** (`:297` „das Gate aus
Festlegung 2" = die E2E-Tabelle) — das Original war die ganze Zeit korrekt, nur
die Wiedergabe nicht. **Nicht** dieselbe Klasse ist
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, wo ein **Wert** gegen eine
Messung driftet: Ein Abschnittsverweis ist keine Größe, er ist eine **Adresse**;
er wird nicht nachgemessen, sondern **aufgeschlagen**.
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` benennt den **Befehl**, der
seinen Satz nicht trägt — hier ist es der **Verweis**, der ihn nicht trägt.

**Warum das zählt:** Ein Verweis ist die Prüf-Form einer Aussage über ein
fremdes Dokument. Trägt er sie nicht, prüft ein späterer Leser an einer Stelle,
an der die Aussage nicht steht — und findet dort möglicherweise etwas, das
plausibel aussieht. Die Aussage selbst bleibt dabei unangetastet, sodass niemand
einen Anlass sieht, nachzusehen. **Träger des Fundes ist ein Leser, kein
Werkzeug:** Ob die zitierte Stelle die Aussage trägt, ist eine Lese-Handlung —
was fehlt, ist die Gegenprobe am Original, und die trägt der Reviewer-Skill
als Lese-Pflicht.
