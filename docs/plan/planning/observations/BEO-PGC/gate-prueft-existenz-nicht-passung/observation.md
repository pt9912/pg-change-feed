# BEO-PGC/gate-prueft-existenz-nicht-passung

**Sub-Area:** `*` (`PGC`, Greenfield) — Gates, die Mengen von Kennungen vergleichen (Codes in
Quelltext, Tabelle und Katalog).

Die Beobachtung: Ein Gate prüft, dass eine Kennung existiert und in allen Trägern gleich steht,
nicht, dass sie zur Stelle **passt**. Der erste Vorgang:
`slice-meldungscodes-warnungen-heartbeat-diagnose` trug `PCF-W1002` („Stream-Veröffentlichung im
NATS-Adapter“) auch an der Stelle `capture/service.go`, die den prozessinternen Broadcaster
aufruft und NATS nicht berührt. `make meldungscodes-check` blieb grün (der Code ist in Quelltext,
Tabelle und Katalog vorhanden); der Reviewer fand die Fehlzuordnung durch das Lesen der Emittenten
je Code gegen die Bedeutung im Katalog.

**Abgrenzung:** keine Kennungs-Seite (`intern-kennungen-in-ausgelieferten-texten`: dort steht eine
interne Kennung in einem ausgelieferten Text) und keine Verhaltensänderung eines Umbaus
(`umbau-aendert-still-beobachtbares-verhalten`); hier ist die Kennung zulässig und vorhanden, ihre
**Zuordnung** zur Ursache der Stelle ist falsch.

**Warum das zählt:** Der Betreiber filtert nach `code=…` und liest die Maßnahme im Katalog; ein Code
an der falschen Stelle schickt ihn in die falsche Richtung, ohne dass ein Sensor es meldet. Die
Passung steht in keiner maschinenlesbaren Form.

Deklaration: `slice-meldungscodes-warnungen-heartbeat-diagnose`, Review F-1 (MEDIUM).
