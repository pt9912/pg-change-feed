Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.5 (Absatz
„Ausnahme ist die Zitat-Korrektur“: unberührbar bleibt die **Aussage** von
§Entscheidung, §Konsequenzen und §Verglichene Alternativen; das Zitatgerüst
darin darf eine Zitat-Korrektur ändern, wenn der Referent gemessen gleich
bleibt) und
[`ADR-0157`](../../../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(Teil-Supersede von
[`ADR-0073`](../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Entscheidung 1; Entscheidungen 3 und 4 nehmen die MR-Einträge in die Klasse
auf) · seit slice-harness-baseline-v6-14-1. Der Herkunfts-Anker steht in
`AGENTS.md` §3.5 (Zeile „Herkunft: …“).

Zugewiesen im Lese-Schritt der Slice-Closure von
`slice-harness-baseline-v6-14-1` (wellenlos); Architect-Zug: `ADR-0157`
(`f1f6ae70`).

Vergleichseinheit und Normalisierung der Referent-Messung (Bedingung (b) von
`ADR-0157` Entscheidung 1): festgelegt in
[`ADR-0158`](../../../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
(Einheit je Verweisform) und
[`ADR-0159`](../../../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(Einheit einer HTML-`id`, Referent-Messung an MR-Pins, Tag-Paar, Befehlsform),
umgesetzt in `AGENTS.md` §3.5 (Satz „gemessen wird die Einheit, die der Verweis
adressiert“; `fc6d104c`, `fce2d159`). Die Befehlsform
als Skript hinter `make`: Folge-Slice `slice-zitat-vergleich-werkzeug`.

Zähler (abgeleitet): **3×** (evidence/slice-harness-baseline-v6-13-0.md,
evidence/slice-maintainer-ordner-releasing-verschieben.md,
evidence/slice-harness-baseline-v6-14-1.md).
