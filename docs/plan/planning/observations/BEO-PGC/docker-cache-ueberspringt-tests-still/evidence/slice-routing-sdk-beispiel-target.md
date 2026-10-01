**Vorgang:** slice-routing-sdk-beispiel-target

**Fund:** `make sdk-pack-python`, `make sdk-pack-csharp`, `make sdk-pack-kotlin`, `make examples-csharp`
und `make examples-kotlin` endeten bei Reviewer und Verifier mit Exit 0 aus dem Docker-Schicht-Cache und
druckten keine Testzeile. Der Reviewer nannte das als Grenze seiner Aussage „Tests grün“ (Review F-5);
der Verifier fuhr die Test-Stufe jedes Packages und der C#- und Kotlin-Beispiele mit `docker build
--no-cache` und las dort die gedruckten Zahlen (Python `144 passed`, C# `Passed: 137`; Kotlin ohne
Zahl im grünen Lauf, `> Task :test` ausgeführt) und fuhr zehn Mutationen, die rot wurden.

Quelle: `docs/reviews/review-slice-routing-sdk-beispiel-target.md` (Messungen, F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-sdk-beispiel-target.md` (§1 Cache-Disziplin, V-1). <!-- d-check:status-provenance -->
