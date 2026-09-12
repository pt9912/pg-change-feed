# Beleg: slice-021 (Consumer-Registrierung: Zugriffsweg + Verdrahtung)

Vorgang: slice-021 (erster Slice von `welle-6`).

Fund: Der neue CLI-Zugriffsweg (`register-consumer <name>`,
`internal/bootstrap/wiring.go::RegisterConsumer`) nutzt dieselbe
gemeinsame Instanz-DSN wie Store-/Aktivierungs-/Stream-Verbindung,
nicht die `cdc_admin`-Rolle — Folge des bestehenden, unveränderten
Verdrahtungsstands (siehe `observation.md`), keine Neuverschärfung
durch diesen Slice. `welle-6` §6 schließt eine Auflösung in dieser
Welle ausdrücklich aus.

**Ausgang: weiter offen.** Rollen-spezifische DSNs je Adapter/Zugriffsweg
zu verdrahten berührt Bootstrap und alle Driven-/Driving-Adapter-
Konstruktoren zugleich; kein Slice dafür existiert.

Quelle: `docs/plan/planning/in-progress/slice-021-consumer-registrierung-zugriffsweg.md`
§6, `docs/plan/planning/welle-6.md` §6.
