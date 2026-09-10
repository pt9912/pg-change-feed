# Beleg: slice-013

Vorgang: slice-013 (Fehlerzustände sichtbar, CDC-Abstand messbar,
Welle 4).

Fund: die vorab benannte §4-Rückführungs-Bedingung trat ein — reales
`cdc_capture_lag` bräuchte Zeitstempel-Wiring durch vier Schichten
zugleich. Der Implementer präzisierte die Schätzung (`pglogrepl`
trägt den Commit-Zeitstempel bereits, keine neue Wire-Verbindung
nötig) und lieferte stattdessen eine dokumentierte Näherung. Der
Planner benannte das reale Maß als Folge-Slice-Ausschluss statt einer
offenen Rückführung.

Quelle: slice-013-Plan §1/§4/§6.
