**Vorgang:** slice-sdk-kotlin-pack-werkzeug (Review-Fund F-1, Planner-Closure)

**Fund:** `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a` trägt seit dem
Python-Nachzug (`slice-sdk-python-pack-werkzeug`) einen stale gewordenen
C#-Absatz: „Eine zweite Sprache oder ein zweiter Vertriebsweg bleibt offen"
— real ist inzwischen die dritte (Python) und mit diesem Slice die vierte
(Kotlin) Sprache beantwortet. Der Python-Nachzug hätte diese C#-Schlusszeile
bereinigen müssen (sie widersprach bereits seit Python dem tatsächlichen
Stand), tat es aber nicht — dieselbe Fehlerklasse wie beim bereits
belegten Fall, nur an einer anderen Stelle, in einem anderen Dokument
(`spec/pflichtenheft.md` Fließtext statt Slice-Plan-§3) und mit **Ursprung
und Vorkommen getrennt**: Ursprung ist `slice-sdk-python-pack-werkzeug`
(dort entstand die Staleness), gefunden wurde sie erst beim Review dieses
Kotlin-Slice (F-1, LOW), das selbst den Python-Absatz korrekt bereinigt
hat (eigener Negativbefund im Review-Report), aber den C#-Nachbarabsatz
bewusst außerhalb seines DoD-Umfangs beließ. Erster Beleg dieser Klasse an
einem Pflichtenheft-Fließtext-Träger statt einem Slice-Plan-Dokument.

**Ausgang:** kein Fix in diesem Slice (bewusst, siehe Slice-Plan §3) — der
nächste Nachzug, der `LH-FA-SST-009.a` ohnehin anfasst (z. B. eine vierte
Sprache/ein vierter Vertriebsweg), bereinigt den C#-Absatz mit.
