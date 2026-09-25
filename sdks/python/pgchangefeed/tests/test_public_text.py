"""Guards the text this package ships to users.

Docstrings, comments and error messages in the package sources are read by
users (``help()``, IDE tooltips, tracebacks, the source distribution). They
describe what the code does in a user's words and carry no internal
identifiers of the project's specification, decision or requirement records.
"""

from __future__ import annotations

import re
from pathlib import Path

import pytest

import pgchangefeed

_INTERNAL_IDENTIFIER = re.compile(r"\b(?:SPEC|ADR|ARC)-\d+|\bLH-(?:FA|QA)-[A-Z]{3}-\d+")


def _package_sources() -> list[Path]:
    root = Path(pgchangefeed.__file__).resolve().parent
    return sorted(root.rglob("*.py"))


def test_the_package_has_sources_to_check() -> None:
    names = {path.name for path in _package_sources()}
    assert {"http_client.py", "models.py", "exceptions.py", "nats_stream_client.py"} <= names


@pytest.mark.parametrize("path", _package_sources(), ids=lambda path: path.name)
def test_package_sources_carry_no_internal_identifier(path: Path) -> None:
    hits = [
        f"{path.name}:{number}: {line.strip()}"
        for number, line in enumerate(path.read_text(encoding="utf-8").splitlines(), start=1)
        if _INTERNAL_IDENTIFIER.search(line)
    ]
    assert not hits, "\n".join(hits)
