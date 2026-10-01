**Vorgang:** slice-routing-sdk-realserver-e2e (Review F-1, MEDIUM)

**Fund:** `sdks/python/pgchangefeed/integration/route_scenario.py` trug im Sammler
`except BaseException as exc:  # noqa: BLE001 - kept for the scenario`. Im Repo gibt es keinen
Linter, dessen Warnung das unterdrückt hätte; es war die einzige `noqa`-Zeile außerhalb von Doku
und Skill. Der Reviewer ordnete sie als MEDIUM ein (kein Gate unterdrückt, daher nicht HIGH). Die
Fixrunde (`03a1a20a`) fängt `Exception` und trägt einen Zusage-Kommentar; ein `git grep` nach
`noqa|nolint|type: ignore|pylint:|# pragma` unter `sdks` und `tools` trifft danach nichts.

**Form (Ausprägung):** Inline-Suppression in einer Sprache, die der Regelwortlaut nicht nennt.
Schwere MEDIUM; vor dem Merge vom Reviewer gefunden, in der Fixrunde geschlossen.

Quelle: `docs/reviews/review-slice-routing-sdk-realserver-e2e.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-sdk-realserver-e2e.md` (§6 Zeile F-1). <!-- d-check:status-provenance -->
