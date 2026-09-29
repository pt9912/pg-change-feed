# Beleg: slice-sdk-kotlin-grpc-administration-flaeche

Vorgang: `slice-sdk-kotlin-grpc-administration-flaeche` — viertes Auftreten
der Klasse, erstes in einer **Test-Fixture-Zeichenkette** statt in
KDoc/Plan-Prosa (Kotlin-Review F-1, HIGH,
`docs/reviews/review-sdk-kotlin-grpc-administration-flaeche.md`).

Fund: Die neue Testklasse
`sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientErrorMappingTest.kt`
trug in Zeile 51 und 73 die beiden deutschen Fixture-Strings
`"Rechtsklasse unzureichend für diese RPC"` (`PermissionDenied`) und
`"interner Fehler"` (`Internal`) mitten in der sonst durchgängig
englischen Datei — wortgleich aus dem C#-Formvorbild
`PgChangeFeedAdministrationClientErrorMappingTests.cs` (Zeile 44, 64)
übernommen. Der Review betont die Richtung der Übernahme: die wortgleiche
Übereinstimmung mit dem Vorbild ist gerade nicht die Bestätigung der
Sprachreinheit, sondern der Weg, auf dem das Fragment weiterwandert —
und der C#-Altbestand war im Geschwister-Review
(`review-sdk-csharp-grpc-administration-flaeche.md`) nicht als Befund
dieser Klasse erkannt worden.

Gezogen: Fixrunde `bd10c391` übersetzt beide Kotlin-Strings
(`"insufficient role for this rpc"`, `"internal error"`); der
C#-Altbestand folgte im eigenen Zug `48e04899` mit denselben
Zielformen — die vom Kotlin-Review offene Anmerkung zum Altbestand ist
damit ebenfalls erledigt (`grep` nach den deutschen Fragmenten in
`sdks/` liefert 0 Treffer). Die Assertion-Aussage der Fixrunde (der
Text wird in keinem Test assertiert, beide Tests prüfen nur den
Statuscode) ist am Quelltext belegt. Nachgemessen vom Verifier
(`verifikation-slice-sdk-kotlin-grpc-administration-flaeche.md` §3:
`git show bd10c391`, grep Exit 1, No-Cache-Testlauf 19 Suites/89
Testfälle/0 failures).

Einordnung: dieselbe Klasse, vierter Form-Teil — KDoc-Kommentar
(Erstauftreten C#, wortgleiche Übernahme Kotlin), Plan-Prosa (drittes
Auftreten, siehe `evidence/slice-sdk-python-http-reale2e.md`), jetzt
die Fixture-Zeichenkette eines Unit-Tests. Kein neuer
Schwellen-Übertritt: der Ausgang (verkörpert, Reviewer-Skill
HIGH-Punkt) ist unverändert zuständig und hat den Fund real getragen.
