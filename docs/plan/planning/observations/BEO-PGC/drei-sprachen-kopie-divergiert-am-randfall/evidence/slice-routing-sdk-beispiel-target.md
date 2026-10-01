**Vorgang:** slice-routing-sdk-beispiel-target (Review F-3, LOW)

**Fund:** Die Python-README sagte für den NATS-Stream „A target that is blank or contains `.`, `*`,
`>` or whitespace raises `ValueError`“, in der Form der Schwester-READMEs von C# und Kotlin. Das
Verhalten ist dort ein anderes: `stream_changes(target="")` liest das leere Ziel als „kein Ziel“ und
abonniert den Namensraum der Quelle, während `BuildTargetSubject("")` in C# und Kotlin wirft. Dieselbe
Aussage über den Randfall „leeres Ziel“ stand in drei Sprachen mit zwei Verhalten. Die Verhaltens-
Differenz selbst folgt dem Parameter-Schnitt (Python hat kein Subjekt-Argument) und war im Plan benannt;
der Satz der README war enger als der Code. Behoben in der Fixrunde: die Python-README und der Docstring
nennen das leere Ziel als „kein Ziel“ und den Fehlerzeitpunkt (erster `next()`), die C#- und
Kotlin-READMEs nennen, dass ein leeres Ziel dort ungültig ist. HTTP, SSE und gRPC lesen die Eingaben
`null`/`""`/`eu`/`a&b=c` in allen drei Sprachen gleich (Tests mit derselben Eingabetabelle).

Quelle: `docs/reviews/review-slice-routing-sdk-beispiel-target.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-sdk-beispiel-target.md` (§5 Zeile F-3). <!-- d-check:status-provenance -->
