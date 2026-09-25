Zustand: **gestrichen** — Ausgang: **gestrichen** mit Begründung · seit
welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.2).

Begründung: die Klasse ist ein DoD-Wortlaut, der enger klingt als die tatsächliche
Kopplung (eine gemeinsame Docker-Stufe koppelt den benannten Bau-Kontext `proto` an
alle Programme). Der Fehler ist **laut**: ein Bau ohne den Kontext bricht an der
`COPY --from=proto`-Zeile ab, ein stiller Defekt ist nicht möglich; jeder der drei
`make sdk-pack-*`-Aufrufe trägt `--build-context proto=proto`. Der Aufwand einer
eigenen Regel übersteigt den möglichen Schaden.

Zähler: 3× (Dateien unter `evidence/`).
