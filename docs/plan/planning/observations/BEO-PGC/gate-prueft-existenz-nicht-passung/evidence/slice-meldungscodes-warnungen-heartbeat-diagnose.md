**Vorgang:** slice-meldungscodes-warnungen-heartbeat-diagnose (Review F-1; Verifikation §3 M1, §4)

**Fund:** Der Code `PCF-W1002` (Stream-Veröffentlichung im NATS-Adapter) stand an zwei Stellen: im
NATS-Publisher (`natsstream/publisher.go`) und in `capture/service.go`, wo der Fehlschlag des
prozessinternen Broadcasters gemeldet wird, der NATS nicht berührt. `make meldungscodes-check` blieb
Exit 0 (Code in Quelltext, Tabelle und Katalog vorhanden); der Reviewer fand die Fehlzuordnung durch
das Lesen der Emittenten gegen die Katalog-Bedeutung. Schwere MEDIUM, Träger-Typ Warn-Stelle im
Erfassungspfad. Behoben in `7fb3dfc9` (der Code entfällt an der Erfassungs-Stelle), gebunden an
`TestCaptureLoggtFehlschlaegeUeberDenInjiziertenPort`; der Verifier fuhr die Mutation (Code wieder
angehängt) rot und las alle 30 `W`-Zuordnungen ohne weiteren Fund. Ein Re-Review nach der Fixrunde
fand nicht statt.

Quelle: `docs/reviews/review-slice-meldungscodes-warnungen-heartbeat-diagnose.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-meldungscodes-warnungen-heartbeat-diagnose.md` (§3, §4, §9). <!-- d-check:status-provenance -->
