import pytest

from pgchangefeed import ClientOptions


def test_client_options_holds_address_and_token() -> None:
    options = ClientOptions(address="https://example.invalid", api_token="secret")
    assert options.address == "https://example.invalid"
    assert options.api_token == "secret"


def test_client_options_rejects_empty_address() -> None:
    with pytest.raises(ValueError):
        ClientOptions(address="", api_token="secret")


def test_client_options_rejects_empty_token() -> None:
    with pytest.raises(ValueError):
        ClientOptions(address="https://example.invalid", api_token="")
