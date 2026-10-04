"""Quellseite der Python-Kompatibilitätsmessung (Schritt A5): liest die
Konstruktor-Signaturen der Fehlertypen beider Hierarchien unter der installierten
Bibliothek. Die 0.5.x-Aufrufform muss binden; ein dritter positionaler Wert
trifft kein Schlüsselwort-Argument und scheitert. Die Messung läuft unter der
Version, die der Treiber installiert hat."""

from __future__ import annotations

import inspect
import sys

import grpc
from importlib.metadata import version

from pgchangefeed import exceptions as ex

HTTP_TYPEN = [
    ex.PgChangeFeedError,
    ex.PgChangeFeedBadRequestError,
    ex.PgChangeFeedUnauthorizedError,
    ex.PgChangeFeedForbiddenError,
    ex.PgChangeFeedNotFoundError,
    ex.PgChangeFeedServerError,
    ex.PgChangeFeedUnexpectedStatusError,
    ex.PgChangeFeedMalformedResponseError,
]
GRPC_TYPEN = [
    ex.PgChangeFeedGrpcError,
    ex.PgChangeFeedGrpcInvalidArgumentError,
    ex.PgChangeFeedGrpcUnauthenticatedError,
    ex.PgChangeFeedGrpcPermissionDeniedError,
    ex.PgChangeFeedGrpcNotFoundError,
    ex.PgChangeFeedGrpcInternalError,
    ex.PgChangeFeedGrpcUnexpectedStatusError,
]


def bindet(signatur: inspect.Signature, *args: object) -> bool:
    try:
        signatur.bind(None, *args)
    except TypeError:
        return False
    return True


def main() -> int:
    installiert = version("pgchangefeed")
    rot = False
    gelesen = 0
    for typ in HTTP_TYPEN + GRPC_TYPEN:
        signatur = inspect.signature(typ.__init__)
        grpc_typ = typ in GRPC_TYPEN
        erste = grpc.StatusCode.NOT_FOUND if grpc_typ else 418
        alt_form = bindet(signatur, erste, "text")
        dritter = bindet(signatur, erste, "text", "PCF-E0001")
        parameter = signatur.parameters.get("message_code")
        schluesselwort = parameter is not None and parameter.kind is inspect.Parameter.KEYWORD_ONLY
        standard_leer = parameter is not None and parameter.default is None
        gelesen += 1
        urteil = "wie-erwartet"
        if installiert.startswith("0.5."):
            if not alt_form or parameter is not None:
                urteil = "ROT"
        elif not alt_form or dritter or not schluesselwort or not standard_leer:
            urteil = "ROT"
        if urteil == "ROT":
            rot = True
        print(
            f"KOMPAT python A5 signatur {installiert} {typ.__name__}: "
            f"0.5.x-Form bindet={alt_form}, drittes Positional bindet={dritter}, "
            f"message_code keyword-only={schluesselwort}, Standard None={standard_leer} {urteil}"
        )
    print(f"KOMPAT python A5 signatur {installiert}: {gelesen} Typen gelesen")
    return 1 if rot else 0


if __name__ == "__main__":
    sys.exit(main())
