# Beleg: slice-nats-drittstream-core

Vorgang: `slice-nats-drittstream-core` — §6-Risiko „Verhaltensänderung eines
bereits ausgelieferten Merkmals": die neue `CDC_NATS_STREAM_TOKEN` verlangt
dem NATS-Server einen serverweiten Token ab, auch für die bislang anonyme
Wecksignal-Verbindung; die bestehende Demo-Umgebung (`examples/compose.yaml`)
schaltet dieselbe Auth scharf, womit die bereits ausgelieferten
Wecksignal-Clients betroffen wären.

Nachprobe (eigener Lauf, nicht aus dem Bericht übernommen):

```
make example-demo-up                       # Exit 0
make example-run-csharp SURFACE=nats ARGS="--source demo-source --schema public --table orders"
make example-run-kotlin SURFACE=nats ARGS="--source demo-source --schema public --table orders"
make example-demo-down                     # Exit 0, keine Container/Netzwerk zurückgelassen
```

Je Client eine real eingefügte Zeile in `public.orders` während des
Lauschens. Beide verbanden sich über die URL-eingebettete Kennung
(`nats://demo-nats-stream-token@nats:4222` gegen Server-Auth
`--auth demo-nats-stream-token`), meldeten `lauscht auf
cdc.changes.demo-source.public.orders`, empfingen den Weckruf und holten die
Änderung real über `GET /changes` (C#: `change_id 899-1`, Kotlin:
`change_id 911-1`). Das Auth-Risiko ist damit für beide Clients eingelöst.

Nebenbefund derselben Probe: Die gemeldete Payload-Länge weicht zwischen den
Clients ab (C# 37 Byte, Kotlin 0 Byte auf demselben leeren Wecksignal) —
eigener Register-Eintrag:
`docs/plan/planning/observations/BEO-PGC/beispiel-client-falsche-payload-laenge/`.

Quelle: eigener Lauf am 2026-09-18 (Logs des Vorgangs) ·
`docs/plan/planning/in-progress/slice-nats-drittstream-core.md` §6 ·
Verifikationsbericht zu `slice-nats-drittstream-core` (W-1: die
Run-Behauptung der Fixrunde war ohne Repo-Artefakt, diese Datei trägt sie).
