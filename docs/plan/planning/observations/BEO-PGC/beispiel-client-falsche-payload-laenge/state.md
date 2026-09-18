Zustand: **verkörpert** — Ausgang: behoben am 2026-09-18, im ersten Vorgang
(1×), auf Nutzerentscheidung („Bitte vereinheitlichen"). Der C#-Client druckt
nicht mehr `NatsMsg.Size`, sondern `msg.Data?.Length ?? 0`.

Die Ursache war keine Transport-Abweichung, sondern ein **Etikett**: `Size`
ist in `NATS.Client.Core` v3.2.0 die Protokoll-Größe des Rahmens (30-Bit-Feld
des NATS-Headers, Subjekt **plus** Nutzdaten), während Kotlin `msg.data.size`
und Go `len(msg.Data)` druckten. Die gemessenen 37 Byte sind exakt die
Zeichenzahl des Subjekts `cdc.changes.demo-source.public.orders`
(12+11+1+6+1+6). Alle drei Clients melden jetzt dieselbe Größe — die
Nutzdaten-Länge; real nachgemessen am 2026-09-18: nach der Änderung meldet
auch der C#-Client `Payload 0 Byte` auf dem leeren Wecksignal.

Bleibt als Eintrag stehen, weil die Klasse („Beispiel-Client meldet eine
andere Größe als ihre zwei Geschwister") jederzeit erneut auftreten kann;
ein Zähler wird nicht gesetzt, er folgt aus den Dateien.
