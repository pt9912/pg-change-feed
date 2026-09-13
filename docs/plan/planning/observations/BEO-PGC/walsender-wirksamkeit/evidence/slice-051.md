# Beleg: slice-051

Vorgang: slice-051 (Publication-Entzug-Wirksamkeit — isolierter Beleg).

Fund: der isolierte Testabschnitt in `tools/harness/run-integration-tests.sh`
(dedizierte Tabelle `feed_mvp_walsender_timing`, aktiviert über die reguläre
`cdc.enable_table`-Antragskette, Publication-Mitgliedschaft direkt per
`ALTER PUBLICATION ... DROP TABLE` als `$PG_USER` entzogen — **ohne**
`cdc.disable_table` — sodass die App-seitige `Assembler`-Bindung
unverändert aktiv bleibt) zeigt real und reproduzierbar (drei
`make test-integration`-Läufe, identisches Ergebnis): eine nach dem
Entzug eingefügte Zeile erscheint **nicht** in `cdc.changes`. PostgreSQLs
bereits laufende Decoding-Session filtert eine per `ALTER PUBLICATION`
entzogene Tabelle sofort aus, ohne Neuaufbau der Session — die
Verzögerungs-Annahme der Beobachtung trifft für die geprüfte
PostgreSQL-Version (18-alpine, `compose.yaml`) nicht zu.

Quelle: `docs/reviews/architect-verdict-walsender-wirksamkeit.md`
(Testspezifikation), `tools/harness/run-integration-tests.sh` (Abschnitt
„Publication-Entzug-Wirksamkeit — isolierter Beleg"), Implementer-Bericht
slice-051 (drei reale `make test-integration`-Läufe, konsistentes
Ergebnis).
