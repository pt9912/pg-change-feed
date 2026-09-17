# BEO-PGC/konfigurierter-pfad-ohne-umbenennungs-schutz

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft
Konfigurationsdateien, die einen anderen Repo-Pfad hartcodiert referenzieren,
keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine Konfigurationsdatei (hier `.d-check.yml`) referenziert
einen anderen Pfad im Repo als **Literal** — kein Anker, keine ID, kein
Sensor, der die Referenz gegen den Ist-Bestand hält. Wird die referenzierte
Datei umbenannt oder verschoben, bemerkt das kein Gate: Der referenzierende
Config-Block bleibt syntaktisch gültig, verweist aber ins Leere oder — falls
der Vertrag "Datei existiert" statt "Datei existiert unter genau diesem
Pfad" lautet — auf gar nichts Prüfbares mehr.

**Abgrenzung zu `BEO-PGC/generierte-artefakte-ohne-sync-sensor`:** Dort läuft
erzeugter **Inhalt** gegen seine Quelle auseinander (Protobuf-Code vs.
`.proto`-Datei) — ein Sync-Problem. Hier ist der Gegenstand ein **Pfad-Zeiger**,
dessen Ziel selbst unverändert bleiben kann; das Problem ist die
Umbenennung/Verschiebung des Ziels, nicht dessen Inhalt.

Deklaration: `slice-d-check-trace-rtm` (§6 Risiko 1): `.d-check.yml`s
`trace.coverage: [{files: [docs/user/e2e-abdeckung.md], ...}]` referenziert
das Erzeugnis von `make test-integration`
(`tools/harness/run-integration-tests.sh`) als hartcodierten Pfad. Aktuell
kein Umbenennungs-Anlass erkennbar, aber keine strukturelle Garantie gegen
künftige Drift — der Pfad ist eine Kopplung an einen anderen Lauf, keine vom
`trace.coverage`-Block selbst kontrollierte Konstante.
