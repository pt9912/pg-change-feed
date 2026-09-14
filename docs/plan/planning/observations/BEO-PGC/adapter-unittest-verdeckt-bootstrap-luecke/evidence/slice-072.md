**Vorgang:** slice-072
**Fund:** Der neue `503`-Pfad des SSE-Endpunkts ist über den regulären
Bootstrap-Pfad strukturell **unerreichbar** (sobald `CDC_HTTP_ADDR` gesetzt
ist, wird der `Broadcaster` garantiert verdrahtet) — der Unit-Test erreicht
ihn nur per direkter `Config`-Konstruktion. Genau das Muster dieses
Eintrags: die Whitebox-Testform sieht die reale Verdrahtung nicht. Zweites,
unabhängiges Auftreten.
