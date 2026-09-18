# Lang laufende Demo-Umgebung verpasst eine neue Betreiber-Variable eines Vorgänger-Slices

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Demo-/Quickstart-Umgebung unter `examples/`, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Demo-Umgebung (`make example-demo-up`, `examples/compose.yaml`) startet
den Feed-Container einmalig mit dem zu diesem Zeitpunkt aktuellen Image
(`ghcr.io/pt9912/pg-change-feed:dev`, per `make image` geladen) und der zu
diesem Zeitpunkt aktuellen `examples/.env`. Liefert ein späterer Slice eine
**neue, additive** Umgebungsvariable (Zwei-Bedingungen-Aktivierung o. ä.)
und trägt `examples/.env` sie bereits nach, bleibt ein bereits laufender,
nicht neu gestarteter Feed-Container trotzdem auf dem alten Verdrahtungs-
Stand — `docker inspect` zeigt die alte `Config.Env`-Liste, nicht die neue
Datei. Ein Folge-Slice, der real gegen diese Umgebung testet, sieht dadurch
scheinbar ein defektes Feature (kein Empfang auf dem neuen Zustellweg),
obwohl Code und Konfiguration beide korrekt sind — der einzige Defekt ist
die veraltete Laufzeit-Instanz. Weder `make example-demo-up` noch
`make example-run-go` warnen davor; ein Container mit stunden- oder
tagealtem Uptime sieht auf den ersten Blick „gesund" aus (`docker ps`
zeigt `healthy`).
