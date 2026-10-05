**Vorgang:** slice-harness-targets-inhalt-bereinigen (Review F-4, INFO; vom Implementer in §3 des Plans als „bemerkt, nicht geändert“ geführt)

**Fund:** Vorher-Nachher-Sätze in der Doku-Prosa zweier Harness-Verträge, außerhalb der Befund-Liste des Slice:

- `harness/targets/sdk-integration.md`, alle drei Abschnitte: „baut die neue `integration`-Docker-Stufe“, „Pins aus der bestehenden `build`-Stufe wiederverwendet“, „Eigentum des bestehenden Runners“, „die HTTP-Fläche … nachgezogen“.
- `harness/sensors/generated-sync.md` §Vertrag: „Bis `slice-generated-sync-tar-export` lief hier stattdessen …“, „Zwei Eigenschaften des Bind-Mount-Mechanismus entfallen mit dem Stufen-Wechsel“, „Vorher leitete das Skript …“, „trägt … jetzt fest“.

Die Planner-Closure setzte beide in den Ist-Zustand: in `sdk-integration.md` ohne die Wörter „neue“, „bestehende“, „nachgezogen“; in `generated-sync.md` ein Satz zur `tar`-Stream-Extraktion (Herkunft *übernommen*) und die zwei entfallenen Eigenschaften als Grenze 4 und 5.

**Form (Ausprägung):** dieselbe Klasse in Doku-Prosa eines **Sensor- bzw. Target-Vertrags** — die Doku-Prosa-Variante, die der sechste Beleg (`slice-harness-guard-inplace-textwerkzeug`) an Plan und `MR`-Adaption zeigte. In `generated-sync.md` ist es der zweite Fund derselben Stelle (der Beleg `slice-generated-sync-tar-export` in `BEO-PGC/slice-chronik-in-code-kommentar` nennt sie als Reviewer-INFO, damals stehen gelassen). Schwere INFO, vor dem Merge gefunden; Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-harness-targets-inhalt-bereinigen.md` (F-4). <!-- d-check:status-provenance -->
