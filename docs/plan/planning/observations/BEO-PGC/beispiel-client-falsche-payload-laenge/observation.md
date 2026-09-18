# BEO-PGC/beispiel-client-falsche-payload-laenge

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Beispiel-Clients unter `examples/`, keine eigene Sub-Area im Sinn der
Modul-Deklaration).

Die Beobachtung: Der C#-Wecksignal-Beispiel-Client meldet für das **leere**
Wecksignal (`SPEC-017`: leerer Payload) eine Länge von 37 Byte, während der
Kotlin-Client desselben Beispiels und der Harness-Testclient `natssub` auf
demselben Subjekt 0 Byte messen.

Real gemessen am 2026-09-18 gegen die laufende Demo-Umgebung
(`make example-demo-up`), je eine real eingefügte Zeile in `public.orders`:

```
C#:      nats-client: Weckruf auf cdc.changes.demo-source.public.orders (Payload 37 Byte) — hole die Aenderung ueber HTTP
Kotlin:  nats-client: Weckruf auf cdc.changes.demo-source.public.orders (Payload 0 Byte) — hole die Aenderung ueber HTTP
```

Beide Clients holten die Änderung anschließend real über `GET /changes`
(C#: `change_id 899-1`, Kotlin: `change_id 911-1`) — der Defekt betraf
allein die **gemeldete Zahl**, nicht die Funktion.

Aufgeklärt am 2026-09-18: es war keine Transport-Abweichung, sondern ein
**Etikett**. `NATS.Client.Core` v3.2.0 füllt `NatsMsg.Size` aus dem
30-Bit-Größenfeld des NATS-Protokoll-Headers — das ist die Rahmen-Größe
(Subjekt **plus** Nutzdaten), nicht der Nutzdaten-Anteil; die gemessenen
37 Byte sind exakt die Zeichenzahl des Subjekts
`cdc.changes.demo-source.public.orders`. Kotlin (`msg.data.size`) und Go
(`len(msg.Data)`) druckten dagegen schon die Nutzdaten-Länge. Behoben durch
Vereinheitlichung auf die Nutzdaten-Länge (`msg.Data?.Length ?? 0`); real
nachgemessen: der C#-Client meldet danach `Payload 0 Byte` auf demselben
leeren Wecksignal.

Kein Gate liest die Ausgabe eines Beispiel-Clients (`make example-run-*` ist
Werkzeug, kein Gate; `make examples-*` fährt keinen echten NATS-Server) —
der Wächter ist das Lesen, kein Sensor.

Gefunden bei der Nachprobe des §6-Risikos „bestehende Wecksignal-Demo bricht
unter der neuen Server-Auth" in `slice-nats-drittstream-core`; nicht
Gegenstand dieses Slice (der Beleg dient der Auth-Frage, nicht der
Payload-Form).
