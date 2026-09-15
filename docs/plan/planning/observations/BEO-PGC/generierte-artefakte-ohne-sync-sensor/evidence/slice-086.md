# Beleg: slice-086

Vorgang: `slice-086` — das Changes-Lesen über die HTTP-API.

Fund: Der Vorgang berührt **zwei** der namentlich geführten Erzeugnisse dieses
Eintrags und behandelt beide weiterhin von Hand:

- `docs/user/e2e-abdeckung.md` — der Slice **erzeugt** die Tabelle neu (das ist
  Teil seiner DoD: der volle `make test-integration`-Lauf), weil seine
  Änderungen an Client und Runner die Zeilennummern verschieben. Das Erzeugnis
  folgt damit weiterhin **freiwilligen** Läufen; kein Sensor hält es gegen seine
  Quelle.
- `tools/schema/plan.yaml` — `make test-store`/`make test-integration` schreiben
  den Pflicht-Report; der Implementer hat ihn **vor jedem Commit von Hand**
  zurückgenommen. Er ist in **keinem** der Commits.

Dazu `harness/image-hash.txt`: der Digest ist nach `ADR-0044` der **Lauf-Beleg**
eines `make image`-Laufs, kein Inhalts-Fingerabdruck — er wird committet, prüft
aber nichts gegen die Quelle. Er gehört in dieselbe Familie, ist aber nicht
einer der drei ursprünglich benannten Träger.

Quelle: Implementer-Bericht `slice-086` ·
`docs/plan/planning/in-progress/slice-086-changes-lesen-api.md` §7 ·
`tools/harness/run-integration-tests.sh:73-189` (der Schreiber vergleicht und
schreibt nur bei inhaltlicher Abweichung) · `ADR-0044` (Digest-Semantik).
