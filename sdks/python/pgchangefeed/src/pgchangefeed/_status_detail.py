"""Reads the message code out of the ``grpc-status-details-bin`` trailer.

The trailer carries a serialized ``google.rpc.Status``; its repeated field 3
holds ``google.protobuf.Any`` entries, and an entry whose type is
``google.rpc.ErrorInfo`` carries the message code as ``reason`` (field 1) next
to its ``domain`` (field 2). The package depends on no ``google.rpc`` classes,
so this module reads exactly these three levels of fields with the protobuf
runtime that is already a dependency. It never raises: whatever it cannot read
yields ``None``.
"""

from __future__ import annotations

from typing import Iterable

from google.protobuf import any_pb2

_TRAILER_KEY = "grpc-status-details-bin"
_ERROR_INFO_TYPE_SUFFIX = "google.rpc.ErrorInfo"
_SERVER_DOMAIN = "pg-change-feed"

_STATUS_DETAILS_FIELD = 3
_ERROR_INFO_REASON_FIELD = 1
_ERROR_INFO_DOMAIN_FIELD = 2

_WIRE_VARINT = 0
_WIRE_FIXED64 = 1
_WIRE_LENGTH_DELIMITED = 2
_WIRE_FIXED32 = 5


def message_code_from_trailers(trailers: Iterable[tuple[str, str | bytes]] | None) -> str | None:
    """Returns the ``reason`` of the first ``ErrorInfo`` of domain
    ``pg-change-feed`` in the status detail trailer, or ``None``."""
    try:
        for key, value in trailers or ():
            if key.lower() == _TRAILER_KEY and isinstance(value, bytes):
                return _message_code_from_status(value)
    except (TypeError, ValueError, AttributeError):
        return None
    return None


def _message_code_from_status(status: bytes) -> str | None:
    try:
        for number, payload in _fields(status):
            if number != _STATUS_DETAILS_FIELD or not isinstance(payload, bytes):
                continue
            code = _message_code_from_any(payload)
            if code is not None:
                return code
    except (ValueError, IndexError):
        return None
    return None


def _message_code_from_any(payload: bytes) -> str | None:
    try:
        detail = any_pb2.Any.FromString(payload)
    except Exception:  # protobuf raises DecodeError, a subclass of Exception, for unreadable bytes
        return None
    if not detail.type_url.endswith(_ERROR_INFO_TYPE_SUFFIX):
        return None
    reason = ""
    domain = ""
    try:
        for number, field in _fields(detail.value):
            if not isinstance(field, bytes):
                continue
            if number == _ERROR_INFO_REASON_FIELD:
                reason = field.decode("utf-8")
            elif number == _ERROR_INFO_DOMAIN_FIELD:
                domain = field.decode("utf-8")
    except (ValueError, IndexError):
        return None
    if domain != _SERVER_DOMAIN or reason == "":
        return None
    return reason


def _fields(data: bytes) -> Iterable[tuple[int, bytes | int]]:
    """Yields ``(field number, payload)`` for each field of a serialized
    message; payloads of length-delimited fields are ``bytes``, others ``int``
    or the raw fixed-width bytes. Raises ``ValueError`` on truncated input."""
    position = 0
    while position < len(data):
        tag, position = _varint(data, position)
        number, wire_type = tag >> 3, tag & 0x7
        if wire_type == _WIRE_VARINT:
            value, position = _varint(data, position)
            yield number, value
        elif wire_type == _WIRE_LENGTH_DELIMITED:
            length, position = _varint(data, position)
            end = position + length
            if end > len(data):
                raise ValueError("truncated field")
            yield number, data[position:end]
            position = end
        elif wire_type in (_WIRE_FIXED64, _WIRE_FIXED32):
            width = 8 if wire_type == _WIRE_FIXED64 else 4
            end = position + width
            if end > len(data):
                raise ValueError("truncated field")
            position = end
        else:
            raise ValueError("unsupported wire type")


def _varint(data: bytes, position: int) -> tuple[int, int]:
    result = 0
    shift = 0
    while True:
        if position >= len(data) or shift > 63:
            raise ValueError("truncated varint")
        byte = data[position]
        position += 1
        result |= (byte & 0x7F) << shift
        if not byte & 0x80:
            return result, position
        shift += 7
