"""Rohdraht-Probe der HTTP-Seite (Schritt B1 von make test-sdk-altserver): ein
POST /tables/enable auf eine nicht vorhandene Tabelle ohne jedes SDK. Gedruckt
werden der rohe Körper und der Fehlertext; trägt der Körper das Feld `code`,
endet die Probe mit Exit 1 (der Server sendet dann Meldungscodes und ist kein
Server vor ihrer Einführung). Läuft per `python -` im Python-Image des SDK
innerhalb des Compose-Netzes; Eingaben als Umgebungsvariablen."""

from __future__ import annotations

import json
import os
import sys
import urllib.error
import urllib.request

adresse = os.environ["ALTSERVER_HTTP_ADDR"]
token = os.environ["ALTSERVER_ADMIN_TOKEN"]
anfrage = {
    "source": os.environ["ALTSERVER_SOURCE_ID"],
    "schema": "public",
    "table": "sdk_error_code_missing_table",
    "table_id": "sdk-error-code",
    "schema_version_id": "sdk-error-code-v1",
    "version": 1,
    "publication": os.environ["ALTSERVER_PUBLICATION"],
}
request = urllib.request.Request(
    adresse + "/tables/enable",
    data=json.dumps(anfrage).encode(),
    headers={"Authorization": "Bearer " + token, "Content-Type": "application/json"},
    method="POST",
)
try:
    urllib.request.urlopen(request, timeout=30)
except urllib.error.HTTPError as fehler:
    status = fehler.code
    koerper = fehler.read().decode()
else:
    print("ALTSERVER B1 FEHLER: die Aktivierung einer fehlenden Tabelle endete nicht mit einem Fehlerstatus")
    sys.exit(1)

print(f"BODY {koerper}")
print(f"STATUS {status}")
try:
    objekt = json.loads(koerper)
except ValueError:
    print("ALTSERVER B1 FEHLER: der Fehlerkörper ist kein JSON")
    sys.exit(1)
if status != 404:
    print("ALTSERVER B1 FEHLER: erwartet Status 404")
    sys.exit(1)
if "code" in objekt:
    print("ALTSERVER B1 FEHLER: der Fehlerkörper trägt das Feld code — der Server sendet Meldungscodes")
    sys.exit(1)
print(f"TEXT {objekt['error']}")
