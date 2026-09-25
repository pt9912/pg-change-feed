"""Guards the package README against drifting from the code.

The README is the package description shown to users. This test reads it and
checks what it promises:

- every ```python block compiles and imports only names that exist;
- every call in a block -- a client method, a constructor, or a function of an
  imported module -- binds against the real signature: no unknown keyword, no
  missing required argument, no surplus positional argument;
- every public method of ``PgChangeFeedHttpClient`` and every live stream
  client has a row in the API overview, with the parameters the code has;
- the change table lists exactly the fields of ``pgchangefeed.models.Change``,
  and the stream paragraph lists exactly the fields of ``StreamChange`` and of
  the generated gRPC ``Change`` message;
- the error table names every exception of the package, with the status code
  the client maps it to.

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
from pgchangefeed.grpc_gen import changestream_pb2
from pgchangefeed.http_client import _STATUS_TO_ERROR
from pgchangefeed.models import Change, StreamChange

_CLIENTS = {
    cls.__name__: cls
    for cls in (
        PgChangeFeedHttpClient,
        PgChangeFeedGrpcClient,
        PgChangeFeedSseClient,
        PgChangeFeedNatsStreamClient,
    )
}

_NUMBER_WORDS = {
    word: number
    for number, word in enumerate(
        "zero one two three four five six seven eight nine ten eleven twelve "
        "thirteen fourteen fifteen sixteen seventeen eighteen nineteen twenty".split()
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


def _namespace(tree: ast.AST) -> dict[str, object]:
    """The names a README block imports, resolved by running only its imports."""
    namespace: dict[str, object] = {}
    for node in ast.walk(tree):
        if isinstance(node, (ast.Import, ast.ImportFrom)):
            exec(compile(ast.Module([node], []), "<README imports>", "exec"), namespace)
    return namespace


def _resolve(func: ast.expr, namespace: dict[str, object]) -> object | None:
    if isinstance(func, ast.Name):
        return namespace.get(func.id)
    if isinstance(func, ast.Attribute):
        base = _resolve(func.value, namespace)
        return getattr(base, func.attr, None) if base is not None else None
    return None


def _binds(target: object, call: ast.Call, *, unbound_method: bool) -> bool:
    """Whether ``call`` fits the signature of ``target`` (values are ignored)."""
    assert not any(isinstance(arg, ast.Starred) for arg in call.args), ast.dump(call)
    assert all(keyword.arg is not None for keyword in call.keywords), ast.dump(call)
    placeholders = [None] * (len(call.args) + (1 if unbound_method else 0))
    keywords = {keyword.arg: None for keyword in call.keywords}
    try:
        inspect.signature(target).bind(*placeholders, **keywords)
    except TypeError:
        return False
    return True


def test_every_client_has_a_readme_example() -> None:
    constructed = set()
    for block in _python_blocks():
        for node in ast.walk(ast.parse(block)):
            if isinstance(node, ast.Call) and isinstance(node.func, ast.Name):
                constructed.add(node.func.id)
    assert set(_CLIENTS) <= constructed, sorted(set(_CLIENTS) - constructed)


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


@pytest.mark.parametrize("index", range(len(_python_blocks())))
def test_every_call_of_a_readme_example_fits_the_real_signature(index: int) -> None:
    tree = ast.parse(_python_blocks()[index])
    namespace = _namespace(tree)
    assigned = {
        _CLIENTS[node.value.func.id]
        for node in ast.walk(tree)
        if isinstance(node, ast.Assign)
        and any(isinstance(target, ast.Name) and target.id == "client" for target in node.targets)
        and isinstance(node.value, ast.Call)
        and isinstance(node.value.func, ast.Name)
        and node.value.func.id in _CLIENTS
    }
    client_classes = assigned or set(_CLIENTS.values())
    checked = 0
    for node in ast.walk(tree):
        if not isinstance(node, ast.Call):
            continue
        func = node.func
        if isinstance(func, ast.Attribute) and isinstance(func.value, ast.Name) and func.value.id == "client":
            candidates = [getattr(cls, func.attr) for cls in client_classes if hasattr(cls, func.attr)]
            assert candidates, f"client has no method {func.attr}"
            assert any(_binds(method, node, unbound_method=True) for method in candidates), (
                f"client.{func.attr} called with arguments its signature does not accept: "
                f"{ast.unparse(node)}"
            )
            checked += 1
            continue
        target = _resolve(func, namespace)
        if target is not None and callable(target):
            assert _binds(target, node, unbound_method=False), (
                f"call does not fit the signature of {ast.unparse(func)}: {ast.unparse(node)}"
            )
            checked += 1
    assert checked, f"README example {index} contains no checkable call"


def _parameters(signature_text: str) -> list[str]:
    call = ast.parse(signature_text, mode="eval").body
    assert isinstance(call, ast.Call), signature_text
    return [arg.id for arg in call.args] + [keyword.arg for keyword in call.keywords]


def _code_parameters(target) -> list[str]:
    return [name for name in inspect.signature(target).parameters if name != "self"]


def _documented_api() -> dict[str, list[str]]:
    """Method name (qualified by class for the live streams) to documented parameters."""
    documented: dict[str, list[str]] = {}
    for row in _table_rows(_section("API overview")):
        first = _code_spans(row[0])
        if not first:
            continue
        if first[0].startswith("PgChangeFeed"):
            cls_name = first[0].split("(")[0]
            assert _parameters(first[0]) == _code_parameters(_CLIENTS[cls_name].__init__), first[0]
            method = _code_spans(row[1])[0]
            documented[f"{cls_name}.{method.split('(')[0]}"] = _parameters(method)
        else:
            documented[f"PgChangeFeedHttpClient.{first[0].split('(')[0]}"] = _parameters(first[0])
    return documented


def test_the_api_overview_names_real_methods_with_their_parameters() -> None:
    for qualified, parameters in _documented_api().items():
        cls_name, method = qualified.split(".")
        assert parameters == _code_parameters(getattr(_CLIENTS[cls_name], method)), qualified


def test_the_api_overview_constructors_match_the_real_ones() -> None:
    spans = _code_spans(_section("API overview"))
    for cls_name, cls in _CLIENTS.items():
        signatures = [span for span in spans if span.startswith(f"{cls_name}(")]
        assert signatures, f"the API overview shows no constructor of {cls_name}"
        assert _parameters(signatures[0]) == _code_parameters(cls.__init__), signatures[0]


def test_the_api_overview_covers_every_public_method() -> None:
    expected = {
        f"{cls.__name__}.{name}"
        for cls in _CLIENTS.values()
        for name, member in inspect.getmembers(cls, inspect.isfunction)
        if not name.startswith("_")
    }
    assert set(_documented_api()) == expected


def test_the_change_table_lists_exactly_the_fields_of_change() -> None:
    documented = {name for row in _table_rows(_section("The change object")) for name in _code_spans(row[0])}
    assert documented == {field.name for field in dataclasses.fields(Change)}


def test_the_stream_paragraph_lists_exactly_the_fields_of_the_stream_message() -> None:
    section = _section("The change object")
    match = re.search(r"with (\w+) fields: ([^.]*)\.", section)
    assert match, "the change-object section has no 'with <n> fields: ...' sentence"
    listed = _code_spans(match.group(2))
    stream_fields = {field.name for field in dataclasses.fields(StreamChange)}
    assert set(listed) == stream_fields
    assert len(listed) == len(stream_fields)
    assert _NUMBER_WORDS[match.group(1)] == len(stream_fields)
    assert set(listed) == {field.name for field in changestream_pb2.Change.DESCRIPTOR.fields}
    missing = re.search(r"They carry no ([^.]*)\.", section)
    assert missing, "the change-object section does not name the fields the streams lack"
    assert set(_code_spans(missing.group(1))) == {field.name for field in dataclasses.fields(Change)} - stream_fields


def _error_classes() -> list[type[pgchangefeed.PgChangeFeedError]]:
    return [
        getattr(pgchangefeed, name)
        for name in pgchangefeed.__all__
        if isinstance(getattr(pgchangefeed, name), type)
        and issubclass(getattr(pgchangefeed, name), pgchangefeed.PgChangeFeedError)
        and getattr(pgchangefeed, name) is not pgchangefeed.PgChangeFeedError
    ]


def test_the_error_table_names_every_exception_with_its_status_code() -> None:
    rows = {
        _code_spans(row[0])[0]: row[1] for row in _table_rows(_section("Error handling")) if _code_spans(row[0])
    }
    assert set(rows) == {cls.__name__ for cls in _error_classes()}
    status_of = {cls: status for status, cls in _STATUS_TO_ERROR.items()}
    for cls in _error_classes():
        code = re.match(r"(\d{3}) ", rows[cls.__name__])
        if cls in status_of:
            assert code and int(code.group(1)) == status_of[cls], f"{cls.__name__}: {rows[cls.__name__]}"
        else:
            assert code is None, f"{cls.__name__} maps no status code: {rows[cls.__name__]}"
