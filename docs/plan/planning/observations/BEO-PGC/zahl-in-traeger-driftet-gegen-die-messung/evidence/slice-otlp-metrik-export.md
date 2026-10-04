# Beleg: slice-otlp-metrik-export

Vorgang: `slice-otlp-metrik-export` (Teil a) — Review F-1 (HIGH,
`review-slice-otlp-metrik-export.md`), Klasse „Zahl im Träger driftet gegen die
Messung", **neue Träger-Form**: die Zahl einer Probe-Messung trägt eine
Entscheidungsbegründung.

Quelle: `docs/reviews/review-slice-otlp-metrik-export.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-otlp-metrik-export.md` (§5 und §6). <!-- d-check:status-provenance -->

Fund: Der Bibliothek-Prüfpunkt des Plans maß den Fußabdruck der Typen-Bibliothek an
einem **Probe-Aufruf** (ein Wegwerf-Programm mit einem Aufruf der Nachricht): Binary
19677344 → 20013216 Byte (+335872 Byte, +1,7 %). Die Entscheidung stützte sich auf
„Fußabdruck unter 2 %“, eine Schwelle, die die Entscheidungsregel des Plans nicht
kennt (Eigenzusatz des Implementers). Der fertige Code am `HEAD` maß 20512928 Byte
(+835584 Byte, +4,2 %); die Probe-Zahl blieb als Begründung stehen, nach der
Verdrahtung gegen die Messung am Endstand falsch. Der Reviewer fand es durch
Nachmessen des Binaries an Parent und `HEAD`; der Verifier maß unabhängig dieselben
beiden Zahlen (19677344 und 20512928).

Gezogen: Fixrunde `9812f141` — der Plan trägt beide Zahlen mit Ursprung (die Probe als
„Probe-Aufruf“ gekennzeichnet, die Messung am fertigen Code mit Befehl und gedruckten
Zeilen), die Schwelle „unter 2 %“ ist als Eigenzusatz benannt und entfällt als
Begründung; die Entscheidung trägt jetzt allein der Vergleich mit dem Rückfall (SDK:
+13,9 %, an beiden gemessenen Größen schlechter, im Probe-Bau gemessen).

Einordnung: Ausgang bleibt **verkörpert** (`AGENTS.md` §3.12 Instanz A, Reviewer-Skill
HIGH-Punkt „Zahl im Träger ohne Ursprung — oder gegen die Messung driftend“); HIGH →
Datei unabhängig vom Deckel, und die Form ist neu (nicht Drift durch spätere Arbeit
an einem Träger, sondern eine Messung am **Vorläufer** des Gegenstands, die als
Messung des Gegenstands gelesen wurde). Geschärfte Anwendung: eine Probe-Messung trägt
ihr Etikett „Probe“ bis zur Nachmessung am fertigen Code; ein Entscheidungsgrund
mit Zahl wird nach der Verdrahtung neu gemessen, bevor der Plan ihn als Begründung
stehen lässt.

Querbezug (keine zweite Datei): die erfundene Schwelle „unter 2 %“ als eigene
Begründungs-Zeile liegt in der Nähe der Schwester-Klasse
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`, hat aber eine andere Form
(eine Schwelle ohne Rückhalt, keine falsche Tatsachenbehauptung über den Gegenstand);
der Mechanismus „Zahl aus der Probe als Grund“ ist hier erfasst. Ein erneutes
Auftreten der Eigenzusatz-Schwelle in einer Entscheidungsregel wäre der Anlass für
einen eigenen Eintrag.
