# Beleg: slice-rtm-letzte-zwoelf-tag-only

Vorgang: `slice-rtm-letzte-zwoelf-tag-only` — Review-Finding F-1 präzisierte
den `LH-FA-CON-002`-Kommentar über
`TestE2ERetentionBlockersViewShowsFurthestBehindConsumer`
(`test/integration/integration_test.go`), netto +3 Zeilen.

Fund: Derselbe Fehlermodus wie im unmittelbar vorangegangenen Slice
`rtm-reste-sst-cfg-por` (siehe
`evidence/slice-rtm-reste-sst-cfg-por.md`), diesmal an der Zahlen-Hälfte:
die Fixrunde (Commit `d1136723`) verlängerte den Doc-Kommentar über der
Testfunktion, ohne `docs/user/e2e-abdeckung.md` (generiert von
`make test-integration`) danach neu zu erzeugen — dessen
`LH-FA-CON-002`-Zeile zitierte weiterhin `integration_test.go:511`,
während die Funktion ab diesem Commit real bei Zeile 514 beginnt. Von
`AGENTS.md` §3.13 §Grenze exakt vorhergesagt: ein Zeilen-Lokator
verschiebt sich mit jeder Nachbaränderung, ohne im Text eine
wiederholbare Spur zu hinterlassen, nach der ein `grep`-Suchlauf zuverlässig
suchen könnte.

Nicht vom Implementer im selben Lauf gefunden, sondern vom unabhängigen
Verifier — beiläufig beim Code-Abgleich der eigentlichen sechs
Prüfpunkte aufgefallen, nicht Teil des beauftragten Prüfumfangs. Der
Drift war zum Fundzeitpunkt bereits **inzidentell behoben**: ein
paralleler Fixrunden-Commit (`89e9922f`) für den Nachbar-Slice
`rtm-reste-sst-cfg-por` hatte `docs/user/e2e-abdeckung.md` aus
unabhängigem Anlass bereits komplett neu erzeugt und dabei den
korrekten Stand (`:514`) mitgebracht — real durch einen erneuten
`make test-integration`-Lauf bestätigt („E2E-Abdeckungstabelle
unverändert").

Quelle: `docs/reviews/verifikation-slice-rtm-letzte-zwoelf-tag-only.md` <!-- d-check:status-provenance -->
§Zusätzlicher eigener Fund · Commit `d1136723` (Ursache) ·
Commit `89e9922f` (inzidentelle Behebung, anderer Slice).
