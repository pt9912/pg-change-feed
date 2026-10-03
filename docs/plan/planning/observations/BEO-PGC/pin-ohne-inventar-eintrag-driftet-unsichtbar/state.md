Zustand: offen — Ausgang: **weiter offen**. Die Aufnahme des PostgreSQL-17-Pins
in das Pin-Inventar von
[`ADR-0051`](../../../../adr/0051-cicd-pipeline-github-actions.md)
(`Accepted`) ist eine Folge-ADR des Architects; die Form des Sensors ist offen
(`pin-stale.sh` liest Makefile-Variablen, der Pin ist ein YAML-Matrix-Wert).
Die fünf Basis-Images der SDK-/Beispiel-Dockerfiles sind weder gemessen noch
im Inventar. Zähler (abgeleitet): 2× (evidence/slice-pin-digests-aktualisieren-2026-10.md,
evidence/slice-pin-digests-aktualisieren-2026-10-b.md).
Träger: [`ADR-0146`](../../../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Accepted; Pin-Inventar als Regel über alle Digest-Pins, Sensor P10). Die
Hebung der gedrifteten Pins trägt
[`slice-pin-digests-aktualisieren-2026-10-b`](../../../done/slice-pin-digests-aktualisieren-2026-10-b.md),
den Sensor
[`slice-pin-stale-alle-digest-pins`](../../../open/slice-pin-stale-alle-digest-pins.md);
der Ausgang bleibt **weiter offen**, bis der Sensor im Nachtlauf läuft.
