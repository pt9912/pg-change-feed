# Beleg: slice-028 (Compose-Integrationstest-Nachzug — Rollen-DSN-Trennung, MVP-Umbenennung)

Vorgang: slice-028.

Fund: `slice-028` baute eine Rollen-Verifikation im Compose-Integrationstest
(`tools/harness/run-integration-tests.sh`, Abschnitt „Rollen-DSN-
Verifikation gegen den Compose-Stack") — real gegen die laufende
PostgreSQL-Instanz, mit eigens angelegten, ephemeren Login-Test-Identitäten.
Sie beweist: ein `cdc_reader`-gebundener Login wird bei einem schreibenden
Aufruf abgelehnt, ein Login ohne `REPLICATION`-Attribut wird beim
Verbindungsaufbau für die Replikation abgelehnt, ein Login mit dem
Attribut gelingt. Reviewer und Verifier bestätigten unabhängig (jeweils
eigene Lektüre des Skripts): Die Prüfung ruft **nicht** die tatsächlichen
Replication-Stream-/ACK-Adapter (`receive.NewStream`/`postgresack.New`)
auf und bleibt vom laufenden Feed-Container entkoppelt (der ohnehin mit
Superuser-DSNs verdrahtet ist, unverändert seit `slice-023`).

**Ausgang: weiter offen.** Punkt (2) der Beobachtung
(„Replication-Stream- und ACK-Adapter … sind gegen Rollen-Vertauschung
nicht testgesichert") ist damit nur auf PostgreSQL-Server-Ebene
geschlossen, nicht auf Adapter-Ebene — der ursprünglich angestrebte
vollständige Abschluss von Punkt (2) wurde nicht erreicht. Punkt (1)
(Heartbeat-Grant-Test liest nie den tatsächlichen
`nacharbeit-roles.sql`-Inhalt) blieb unverändert außerhalb des Scopes
dieses Slices.

Quelle: `docs/reviews/review-slice-028.md` F-1, `docs/reviews/verify-slice-028.md`
(unabhängig bestätigt).
