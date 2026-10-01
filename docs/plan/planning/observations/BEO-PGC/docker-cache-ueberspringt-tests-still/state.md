Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen.

Zähler (abgeleitet): **1×** (evidence/slice-routing-sdk-beispiel-target.md). Der Reviewer
nannte die Grenze im Review-Report (F-5, Grenze zu den Sensor-Läufen), der Verifier erzwang die
Ausführung mit `docker build --no-cache` an der Test-Stufe jedes Packages und mit Mutationen.
Die verfügbare Falsifikation ist die gedruckte Testzeile im Bau ohne Cache oder eine rote
Mutation; ein Sensor, der „Bau ohne Cache“ erzwingt, ist nicht gebaut.

**Hinweis, kein Auftreten (slice-routing-sdk-realserver-e2e):** die Mutationsläufe des Verifiers
überschrieben die Tag-Images der Tier-Bauten (`pg-change-feed:sdk-<sprache>-integration`, nicht
`:dev`). Die Tier-Läufe dieses Slice sind keine `make sdk-pack-*`-Bauten; ob die Stufe
`integration` bei unveränderter Eingabe aus dem Cache kommt und dann keine Testzeile druckt, ist
hier weder gemessen noch Teil des Belegs (die Läufe druckten die `ROUTE_RESULT`-Zeilen).
