# `make schema-validate` — Schema-YAML mit d-migrate prüfen

Ausführliche Fassung der Index-Zeilen aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## `make schema-validate`

prüft das neutrale Schema-YAML mit d-migrate (der Lauf selbst braucht keinen Netzzugang, `--network none`, und keinen Bind-Mount: das YAML reist per `COPY` in das Image der Stufe `rollout` von `tools/schema/Dockerfile`; dafür läuft ein `docker build`, der die gepinnten Basis-Images lokal voraussetzt, gemessen: `docker build --no-cache --network none --target guard` endet grün, für die Stufe `rollout` hergeleitet, nicht gefahren; Vorlauf vor jedem `generate`/`migrate`)

**Bindung:** kein Gate, [`ADR-0043`](../../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md), [`ADR-0142`](../../docs/plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
