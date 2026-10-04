"""Gegenrichtung der Python-Kompatibilitätsmessung: ein gegen 0.6.0 geschriebener
Aufrufer übergibt den Meldungscode als Schlüsselwort-Argument und liest die
Eigenschaft. Unter der Bibliothek 0.5.0 nimmt der Konstruktor das Schlüsselwort
nicht an."""

from __future__ import annotations

import sys

import grpc

from pgchangefeed import exceptions as ex

aufrufe = 0


def pruefe(ok: bool, was: str) -> None:
    global aufrufe
    if not ok:
        raise AssertionError(f"Abweichung bei {was}")
    aufrufe += 1


def lauf() -> None:
    http = ex.PgChangeFeedBadRequestError(400, "bad", message_code="PCF-E0001")
    pruefe(http.message_code == "PCF-E0001", "HTTP message_code")

    rpc = ex.PgChangeFeedGrpcNotFoundError(grpc.StatusCode.NOT_FOUND, "m", message_code="PCF-E0002")
    pruefe(rpc.message_code == "PCF-E0002", "gRPC message_code")


def main() -> int:
    schritt = sys.argv[1] if len(sys.argv) > 1 else "?"
    try:
        lauf()
    except Exception as ausnahme:
        print(f"KOMPAT python {schritt}: Ausnahme {type(ausnahme).__name__}: {ausnahme}")
        return 1
    print(f"KOMPAT python {schritt}: {aufrufe} Aufrufe ok")
    return 0


if __name__ == "__main__":
    sys.exit(main())
