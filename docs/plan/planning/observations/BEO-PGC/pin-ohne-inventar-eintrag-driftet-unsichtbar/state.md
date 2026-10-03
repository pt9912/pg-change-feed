Zustand: offen — Zähler unter der Schwelle (3×); Träger gebaut, Runner-Lauf gelesen.
Das Pin-Inventar von
[`ADR-0051`](../../../../adr/0051-cicd-pipeline-github-actions.md)
(`Accepted`) führte den PostgreSQL-17-Pin und fünf Basis-Images der
SDK-/Beispiel-Dockerfiles nicht; Sensor P10 prüft seit
[`ADR-0146`](../../../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
jede Digest-Referenz des Baums. Zähler (abgeleitet): 2× (evidence/slice-pin-digests-aktualisieren-2026-10.md,
evidence/slice-pin-digests-aktualisieren-2026-10-b.md).
Träger: [`ADR-0146`](../../../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Accepted; Pin-Inventar als Regel über alle Digest-Pins, Sensor P10). Die
Hebung der gedrifteten Pins trägt
[`slice-pin-digests-aktualisieren-2026-10-b`](../../../done/slice-pin-digests-aktualisieren-2026-10-b.md),
den Sensor
[`slice-pin-stale-alle-digest-pins`](../../../done/slice-pin-stale-alle-digest-pins.md);
der Sensor ist gebaut und auf dem Runner gelesen: `upstream-drift`, Lauf
`37144584033` (`workflow_dispatch`, Stand `eb4155e5`), Schritt P10 druckt
`pin-stale-all: 15 Referenzen — 15 OK, 0 DRIFT, 0 UNBESTIMMT`. Der erste
planmäßige Nachtlauf (`schedule`) mit P10 ist nicht gelesen.
