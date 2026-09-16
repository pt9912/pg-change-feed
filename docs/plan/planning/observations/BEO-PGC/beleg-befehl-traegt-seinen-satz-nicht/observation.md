# BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in der
ein Beleg-Befehl seine Aussage trägt, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Träger **nennt einen Beleg** — einen Befehl, eine Abfrage,
einen Pfad — als Stütze einer Aussage, und der genannte Beleg **trägt diese
Aussage nicht**: Er misst etwas anderes, zählt etwas anderes oder liefert ein
anderes Ergebnis als das behauptete. Die Aussage kann dabei **wahr** sein; falsch
ist die **Stütze**. Das macht den Fall schwer zu sehen: Wer den Satz liest,
findet ihn plausibel; wer den Befehl ausführt, findet etwas anderes — und nur
der zweite Blick deckt es auf.

Belegt an zwei abgeschlossenen Vorgängen:

- **`slice-079`** (Verifikation `verify-slice-084` V-1): `harness/sensors/coverage-gate.md`
  §Grenze Punkt 1 zitiert `go list -f '{{len .TestGoFiles}}'` als Beleg für
  „fünf Pakete ohne Testdatei". Mit dieser Formel liefert der Lauf **25** — die
  externen Testpakete (`XTestGoFiles`) zählt sie nicht mit; erst
  `TestGoFiles` **und** `XTestGoFiles` ergeben genau die fünf genannten.
- **`slice-085`** (Verifikation `verify-slice-085` V-3): derselbe Satz, derselbe
  Befehl — beim zweiten Vorgang erneut als ungedeckter Beleg gefunden und
  **bewusst liegen gelassen** (mit Adresse), weil er nicht zum Zuschnitt gehörte.

**Warum das zählt:** Ein Beleg ist die Prüf-Form einer Aussage. Trägt er sie
nicht, prüft ein späterer Leser gegen etwas anderes als das Behauptete — und die
Aussage selbst bleibt dabei unangetastet, sodass niemand einen Anlass sieht,
nachzusehen. Die verwandte Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
trifft den **Wert**; hier trifft es den **Befehl**, der ihn stützen soll.
