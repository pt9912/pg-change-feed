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
(C#: `change_id 899-1`, Kotlin: `change_id 911-1`) — der Defekt betrifft
allein die **gemeldete Zahl**, nicht die Funktion. Zwei Deutungen sind
möglich und unentschieden: Der C#-Client misst etwas anderes als die
Payload-Länge (z. B. eine Puffer- oder Kopfgröße), oder seine
NATS-Bibliotheks-Anbindung liefert an dieser Stelle tatsächlich einen
gefüllten Rumpf, obwohl der Absender leer publiziert. Die zweite Deutung
wäre die ernstere und ist mit dem Bestand nicht ausgeschlossen.

Kein Gate liest die Ausgabe eines Beispiel-Clients (`make example-run-*` ist
Werkzeug, kein Gate; `make examples-*` fährt keinen echten NATS-Server) —
der Wächter ist das Lesen, kein Sensor.

Gefunden bei der Nachprobe des §6-Risikos „bestehende Wecksignal-Demo bricht
unter der neuen Server-Auth" in `slice-nats-drittstream-core`; nicht
Gegenstand dieses Slice (der Beleg dient der Auth-Frage, nicht der
Payload-Form).
