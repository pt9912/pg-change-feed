# Beleg: slice-022 (Positions-Bestätigung über denselben Zugriffsweg)

Vorgang: slice-022 (zweiter Slice von `welle-6`).

Fund: `AcknowledgeConsumer` (`internal/bootstrap/wiring.go`) öffnet die
Consumer-State-Verbindung über dieselbe gemeinsame Instanz-DSN wie
`RegisterConsumer` (`slice-021`) und der ursprüngliche Bootstrap
(`slice-011`), nicht über die `cdc_admin`-Rolle — Folge des bestehenden,
unveränderten Verdrahtungsstands, keine Neuverschärfung durch diesen
Slice. `welle-6` §6 schließt eine Auflösung in dieser Welle ausdrücklich
aus. Dies ist der **dritte** unabhängige Vorgang, der dieses Muster
reproduziert — die Schwelle (3×) ist damit erreicht.

**Ausgang: weiter offen** — der Lese-Schritt (Ausgang-Zuweisung nach
Modul 6) läuft regulär bei der `welle-6`-Closure, da `slice-022` einer
Welle angehört.

Quelle: `docs/plan/planning/in-progress/slice-022-consumer-bestaetigung-zugriffsweg.md`
§6, `docs/reviews/verify-slice-022.md` V-1.
