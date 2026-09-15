# BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Zahlen in
Kommentaren, Sensor-Dokumenten und Entscheidungen, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Beobachtung: Ein **Träger** nennt eine Zahl über den Gegenstand, die nicht
mehr stimmt — weil sich der Gegenstand bewegt hat und der Träger nicht mitkam.
Der Träger ist dabei *nicht* die Messung selbst, sondern ihre Beschreibung:
ein Kommentar in einem Messskript, ein Sensor-Dokument, eine
Architektur-Entscheidung. Kein Gate prüft Zahlen gegen die Messung; die
Aussage liest sich wie ein Beleg und ist einer geworden, der nicht mehr gilt.

Belegt an `slice-081`, **zweimal im selben Vorgang**, und die zwei Fundorte
zeigen die Spannweite:

- `tools/harness/db-coverage.sh:18` sagte „der Lauf ueber postgresstorage allein
  traegt **610** Statements" — gemessen **472**; dieselbe Aussage stand in
  `harness/sensors/db-adapter-coverage.md:43` bereits richtig. Zwei Träger
  derselben Zahl widersprachen sich im Baum, ohne deklarierten Gewinner.
- `ADR-0071` und `welle-20.md` nennen weiter **1679** (dort zusätzlich 2467,
  788, 69,74 %) für einen Gegenstand, der nach dem Subjekt-Transfer **1831**
  trägt. Lebende Träger außerhalb des Diffs — der Fund war ein INFO des Reviews,
  kein Defekt des Slice.

**Warum das zählt:** Die Zahl ist die Prüf-Form einer Zusage. Driftet sie,
prüft ein späterer Leser gegen einen Zustand, den es nicht mehr gibt — und er
kann es der Zahl **nicht ansehen**, weil sie ihren Ursprung nicht trägt. Die
Schwester-Klasse `dod-begruendung-unzutreffende-tatsachenbehauptung` betrifft
dieselbe Wurzel an einer anderen Stelle (der Begründung eines DoD-Kriteriums);
hier ist es der **Träger der Messung** selbst.
