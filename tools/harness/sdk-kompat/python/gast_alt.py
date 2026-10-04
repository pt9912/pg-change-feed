"""Gast-Programm der Python-Kompatibilitätsmessung: benutzt jede
Konstruktor- und Lesefläche der 0.5.x-Fehlertypen beider Hierarchien in der Form,
die ein gegen 0.5.0 geschriebener Aufrufer trägt (positionale Argumente)."""

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


class GastFehler(ex.PgChangeFeedError):
    """Fremde Unterklasse der HTTP-Basis ohne eigenen Konstruktor."""


class GastGrpcFehler(ex.PgChangeFeedGrpcError):
    """Fremde Unterklasse der gRPC-Basis ohne eigenen Konstruktor."""


def lauf() -> None:
    http_typen = [
        ex.PgChangeFeedError,
        ex.PgChangeFeedBadRequestError,
        ex.PgChangeFeedUnauthorizedError,
        ex.PgChangeFeedForbiddenError,
        ex.PgChangeFeedNotFoundError,
        ex.PgChangeFeedServerError,
        ex.PgChangeFeedUnexpectedStatusError,
        ex.PgChangeFeedMalformedResponseError,
        GastFehler,
    ]
    for typ in http_typen:
        fehler = typ(418, "text")
        pruefe(fehler.status_code == 418, f"{typ.__name__}(int, str) status_code")
        pruefe(str(fehler) == "text", f"{typ.__name__}(int, str) str")
        pruefe(fehler.args == ("text",), f"{typ.__name__}(int, str) args")

    grpc_typen = [
        ex.PgChangeFeedGrpcError,
        ex.PgChangeFeedGrpcInvalidArgumentError,
        ex.PgChangeFeedGrpcUnauthenticatedError,
        ex.PgChangeFeedGrpcPermissionDeniedError,
        ex.PgChangeFeedGrpcNotFoundError,
        ex.PgChangeFeedGrpcInternalError,
        ex.PgChangeFeedGrpcUnexpectedStatusError,
        GastGrpcFehler,
    ]
    for typ in grpc_typen:
        fehler = typ(grpc.StatusCode.NOT_FOUND, "text")
        pruefe(fehler.code == grpc.StatusCode.NOT_FOUND, f"{typ.__name__}(code, str) code")
        pruefe(str(fehler) == "text", f"{typ.__name__}(code, str) str")

    gefangen = False
    try:
        raise ex.PgChangeFeedNotFoundError(404, "wurf")
    except ex.PgChangeFeedError as fehler:
        gefangen = isinstance(fehler, ex.PgChangeFeedNotFoundError)
    pruefe(gefangen, "except PgChangeFeedError")

    gefangen = False
    try:
        raise ex.PgChangeFeedGrpcNotFoundError(grpc.StatusCode.NOT_FOUND, "wurf")
    except ex.PgChangeFeedGrpcError as fehler:
        gefangen = isinstance(fehler, ex.PgChangeFeedGrpcNotFoundError)
    pruefe(gefangen, "except PgChangeFeedGrpcError")


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
