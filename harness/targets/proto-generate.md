# `make proto-generate` — Go-Code aus den `.proto`-Quellen

Ausführliche Fassung der Index-Zeilen aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## `make proto-generate`

erzeugt den Go-Code des Live-Change-Streams aus `proto/cdc/stream/v1/changestream.proto` und den Go-Code des `Administration`-Service aus `proto/cdc/administration/v1/administration.proto` (`ADR-0130`) — Docker-only: die Dockerfile-Stufe `proto` trägt `protoc` und die gepinnten Plugins `protoc-gen-go`/`protoc-gen-go-grpc`; die Folgestufe `proto-export` (seit slice-104) kopiert beide `.proto`-Quellen per `COPY` hinein und ruft `protoc` **zur Build-Zeit** in einem einzigen Aufruf für beide auf (kein Bind-Mount, kein `--user`-Workaround), ihr `ENTRYPOINT` gibt das Erzeugnis als `tar`-Stream über stdout aus; `tools/harness/proto-generate.sh` baut die Stufe und extrahiert **host-seitig** über `docker run --rm --network none <image> | tar -x -C .` — die Pipe läuft explizit unter `bash` mit `set -o pipefail`, damit ein `docker run`-Fehlschlag den Lauf sichtbar rot macht statt hinter `tar`s Exit-Code zu verschwinden (`AGENTS.md` §3.9; das Makefile-Rezept selbst liefe unter `/bin/sh`, dort `dash` ohne `pipefail`) — die extrahierten Dateien gehören dadurch automatisch dem aufrufenden Nutzer. Die erzeugten Dateien liegen committet im Baum und werden von `make test`/`make image` mitkompiliert; der Generator läuft nur, wenn sich eine der beiden `.proto`-Quellen ändert

**Bindung:** kein Gate, [`ADR-0060`](../../docs/plan/adr/0060-grpc-streaming-mechanismus.md), [`ADR-0130`](../../docs/plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)
