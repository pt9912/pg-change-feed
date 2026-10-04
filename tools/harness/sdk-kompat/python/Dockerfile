# tools/harness/sdk-kompat/python/Dockerfile — Gast-Image der Python-Kompatibilitäts-
# messung (make test-sdk-kompat). Die Basis ist das lokal getaggte Image der
# Stufe `build` von sdks/python/Dockerfile (BASE, vom Runner übergeben): dieses
# Dockerfile trägt kein eigenes Digest-Literal. Der Bau legt nur die Pakete beider
# Versionen samt Abhängigkeiten als Wheel-Sammlung ab; die Läufe der Schritte
# geschehen beim `docker run` (run.sh) in einer frischen virtuellen Umgebung,
# nie in einer RUN-Schicht: der Schicht-Cache überspringt einen Lauf im Bau
# still.
ARG BASE=pg-change-feed:sdk-kompat-base-python
FROM ${BASE}
WORKDIR /kompat

# Das benannte Verzeichnis `neu` trägt die Artefakte von make sdk-pack-python
# (Modus dist) oder ist leer (Modus registry).
COPY --from=neu / /neu

RUN pip download --no-cache-dir -d /wheels/05 pgchangefeed==0.5.0
RUN set -eu; \
    if ls /neu/pgchangefeed-*.whl >/dev/null 2>&1; then \
        pip download --no-cache-dir -d /wheels/06 /neu/pgchangefeed-*.whl; \
        cp /neu/pgchangefeed-*.whl /wheels/06/; \
    else \
        pip download --no-cache-dir -d /wheels/06 pgchangefeed==0.6.0; \
    fi

COPY gast_alt.py gast_neu.py quelle.py run.sh ./
ENTRYPOINT ["/bin/bash", "/kompat/run.sh"]
