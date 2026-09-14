# BEO-PGC/generierte-artefakte-ohne-sync-sensor

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft erzeugte,
aber committete Artefakte und ihre Bindung an die Quelle, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Teil des Repos besteht aus **erzeugten Artefakten, die
committet im Baum liegen** — `tools/schema/plan.yaml` und
`tools/schema/down.sql` (d-migrate-Rollout), `harness/image-hash.txt`
(`make image`), seit `slice-069` zusätzlich der Protobuf-Code
(`internal/adapters/driving/grpc/streamv1/*.pb.go` aus
`proto/cdc/stream/v1/changestream.proto`). Für keines davon existiert ein
Sensor, der das committete Artefakt gegen seine Quelle hält: Ändert jemand
die Quelle und vergisst den Generator, bleibt der Baum konsistent
kompilierbar, aber Quelle und Erzeugnis laufen auseinander — und `make
gates` bliebe grün. Die Bindung ist bisher Disziplin (der erzeugende Lauf
wird vor der Closure gefahren), nicht Mechanik.

Deklaration: `slice-069` (Review-Finding F-6, `docs/reviews/review-slice-069.md`;
der Verifier teilte die Einordnung als „Klasse gehört ins Register",
`docs/reviews/verify-slice-069.md`).
