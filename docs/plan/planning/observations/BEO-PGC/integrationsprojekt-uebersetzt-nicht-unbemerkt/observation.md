# BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Realserver-Integrationsprojekte
der SDK-Tiers: C#-Stufe `integration` in `sdks/csharp/Dockerfile`, Kotlin-Aufgabe
`integrationTest`).

Die Beobachtung: Ein Test-Tier, der nur auf Anforderung läuft (`make test-sdk-*-integration`,
Netzbezug, nicht in `make gates`, in keinem Workflow), übersetzt sein Projekt nur in diesem Lauf.
Eine Signaturänderung im Package bricht das Integrationsprojekt ohne jedes rote Signal bis zum
nächsten Tier-Lauf: `StreamChangesAsync(schema, table, cancellationToken)` (Stream-Filter) machte
drei positionale Token-Aufrufe im C#-Integrationsprojekt zu `CS1503`, und das blieb unbemerkt,
bis der nächste Slice den Tier fuhr. `make sdk-pack-csharp` und der Release-Workflow bauen nur
`PgChangeFeed.Client` und `PgChangeFeed.Client.Tests`.

**Warum das zählt:** Ein Bruch im Test-Tier wird erst beim Slice sichtbar, der den Tier braucht, und
wird dort als Nebenbefund in den Plan gezogen. Der Sensor „Übersetzen der Integrationsprojekte“
existiert nicht; er braucht Netz (NuGet-/Maven-Restore) und passt daher nicht in `make gates`. Ob er
einen Träger bekommt (ein Bau-Schritt in `make sdk-pack-*`, ein advisory Workflow, ein Hinweis in
`make gates`-Nähe) oder der Bruch bis zum nächsten Tier-Lauf akzeptiert bleibt, ist eine eigene
Entscheidung.

**Abgrenzung.** `BEO-PGC/test-runner-stiller-ausschluss` und `BEO-PGC/test-methode-lauft-still-nicht`
zählen Tests, die laufen sollten und es still nicht tun; `BEO-PGC/docker-cache-ueberspringt-tests-still`
den Cache, der eine vorhandene Stufe überspringt. Hier übersetzt das Projekt nicht, und kein Lauf
sieht es.

## Benannt, nicht gezählt

Keine.
