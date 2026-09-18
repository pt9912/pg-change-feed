**Vorgang:** slice-nats-drittstream-example-go
**Fund:** Der reale Closure-Rundlauf (`make example-run-go
SURFACE=nats-stream` gegen die Demo-Umgebung) empfing zunächst kein Event
trotz korrekt eingefügter Zeile in `public.orders`; `docker inspect
cdc-examples-feed --format '{{range .Config.Env}}{{println .}}{{end}}' |
grep NATS` zeigte weder `CDC_NATS_URL` mit eingebettetem Token noch
`CDC_NATS_STREAM_TOKEN` — der seit ~5 Stunden laufende Container war mit
der `examples/.env` von vor `slice-nats-drittstream-core` gestartet worden.
Erst `make image` (Image neu geladen) gefolgt von `make example-demo-down`
+ `make example-demo-up` (Container neu erzeugt) stellte den erwarteten
Empfang her (`change_id=804-1`, dokumentiert im Slice-Plan §2 DoD-Punkt 1).
