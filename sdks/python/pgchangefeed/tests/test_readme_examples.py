"""Guards the package README against drifting from the code.

The README is the package description shown to users. This test reads it and
checks what it promises: every ```python block compiles and imports only names
that exist, every method the blocks call on ``client`` exists on one of the four
client classes, every method row of the API overview names an existing method
with the parameters the code has, the change-object table lists exactly the
fields of ``pgchangefeed.models.Change``, and every exception named in the error
table exists.

The README sits one level above ``pgchangefeed/`` in the repository and next to
``pyproject.toml`` inside the Docker build (sdks/python/Dockerfile copies it
there); both places are tried.
"""

from __future__ import annotations

import ast
import dataclasses
import importlib
import inspect
import re
from pathlib import Path

import pytest

import pgchangefeed
from pgchangefeed import (
    PgChangeFeedGrpcClient,
    PgChangeFeedHttpClient,
    PgChangeFeedNatsStreamClient,
    PgChangeFeedSseClient,
)
from pgchangefeed.models import Change

_CLIENTS = {
    cls.__name__: cls
    for cls in (
        PgChangeFeedHttpClient,
        PgChangeFeedGrpcClient,
        PgChangeFeedSseClient,
        PgChangeFeedNatsStreamClient,
    )
}


def _readme() -> str:
    here = Path(__file__).resolve()
    for candidate in (here.parents[1] / "README.md", here.parents[2] / "README.md"):
        if candidate.is_file():
            return candidate.read_text(encoding="utf-8")
    pytest.fail("README.md not found next to pyproject.toml or one level above pgchangefeed/")


def _python_blocks() -> list[str]:
    return re.findall(r"^```python\n(.*?)^```", _readme(), flags=re.DOTALL | re.MULTILINE)


def _section(title: str) -> str:
    match = re.search(rf"^## {re.escape(title)}\n(.*?)(?=^## |\Z)", _readme(), re.DOTALL | re.MULTILINE)
    assert match, f"README has no section '## {title}'"
    return match.group(1)


def _table_rows(section: str) -> list[list[str]]:
    rows = []
    for line in section.splitlines():
        if line.startswith("|") and not re.match(r"^\|[-| ]+\|$", line):
            rows.append([cell.strip() for cell in line.strip().strip("|").split("|")])
    return rows[1:] if rows else rows


def _code_spans(cell: str) -> list[str]:
    return re.findall(r"`([^`]+)`", cell)


def test_the_readme_has_python_examples() -> None:
    assert len(_python_blocks()) >= 5


@pytest.mark.parametrize("index", range(len(_python_blocks())))
def test_a_readme_example_compiles_and_imports_existing_names(index: int) -> None:
    tree = ast.parse(_python_blocks()[index])
    compile(tree, f"<README example {index}>", "exec")
    for node in ast.walk(tree):
        if isinstance(node, ast.Import):
            for alias in node.names:
                importlib.import_module(alias.name)
        elif isinstance(node, ast.ImportFrom):
            module = importlib.import_module(node.module)
            for alias in node.names:
                assert hasattr(module, alias.name), f"{node.module} has no {alias.name}"


def test_every_method_the_examples_call_on_client_exists() -> None:
    known = {name for cls in _CLIENTS.values() for name in dir(cls) if not name.startswith("_")}
    called: set[str] = set()
    for block in _python_blocks():
        for node in ast.walk(ast.parse(block)):
            if (
                isinstance(node, ast.Call)
                and isinstance(node.func, ast.Attribute)
                and isinstance(node.func.value, ast.Name)
                and node.func.value.id == "client"
            ):
                called.add(node.func.attr)
    assert called, "the examples call no method on `client`"
    assert called <= known, f"not a client method: {sorted(called - known)}"


def _parameters(signature_text: str) -> list[str]:
    inner = signature_text[signature_text.index("(") + 1 : signature_text.rindex(")")]
    return [part.split("=")[0].strip() for part in inner.split(",") if part.strip()]


def _code_parameters(target) -> list[str]:
    return [name for name in inspect.signature(target).parameters if name != "self"]


def test_the_api_overview_names_real_methods_with_their_parameters() -> None:
    checked = 0
    for row in _table_rows(_section("API overview")):
        first = _code_spans(row[0])
        if not first:
            continue
        if first[0].startswith("PgChangeFeed"):
            cls = _CLIENTS[first[0].split("(")[0]]
            assert _parameters(first[0]) == _code_parameters(cls.__init__), first[0]
            method = _code_spans(row[1])[0]
            assert _parameters(method) == _code_parameters(getattr(cls, method.split("(")[0])), method
        else:
            method = first[0]
            target = getattr(PgChangeFeedHttpClient, method.split("(")[0])
            assert _parameters(method) == _code_parameters(target), method
        checked += 1
    assert checked == 13


def test_the_change_table_lists_exactly_the_fields_of_change() -> None:
    documented = {name for row in _table_rows(_section("The change object")) for name in _code_spans(row[0])}
    assert documented == {field.name for field in dataclasses.fields(Change)}


def test_every_exception_of_the_error_table_exists() -> None:
    names = [name for row in _table_rows(_section("Error handling")) for name in _code_spans(row[0])]
    assert len(names) == 7
    for name in names:
        assert issubclass(getattr(pgchangefeed, name), pgchangefeed.PgChangeFeedError), name
