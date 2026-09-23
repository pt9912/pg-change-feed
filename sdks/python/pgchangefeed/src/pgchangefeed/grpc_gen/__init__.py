"""Package marker for the gRPC stub modules generated at Docker build time.

The stub modules (``changestream_pb2.py``/``changestream_pb2_grpc.py``;
protoc names them after the ``.proto`` file) are generated inside the Docker
build (``sdks/python/Dockerfile``) from
``proto/cdc/stream/v1/changestream.proto`` via the additional, named build
context ``proto`` (``--build-context proto=proto``) and are not committed to
this tree (same pattern as C#/Kotlin) — ``gen/**`` remains the Go binding.
This file itself is committed so the generated modules form a regular
subpackage of ``pgchangefeed`` (packaged by ``setuptools`` into the wheel,
importable as ``pgchangefeed.grpc_gen``).
"""