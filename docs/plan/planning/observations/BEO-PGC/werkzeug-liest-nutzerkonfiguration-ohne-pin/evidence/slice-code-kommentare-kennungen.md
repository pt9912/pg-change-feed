**Vorgang:** slice-code-kommentare-kennungen (Review F-1, MEDIUM)

**Fund:** Der `DIFF`-Pfad von `make kommentar-kennungen` zerlegte die `+++`-Zeilen des `git diff`
mit festem Präfix `b/`; der Aufrufer pinnte `--no-color --no-ext-diff`, nicht die
Präfix-Optionen. Gemessen vom Reviewer: `DIFF=HEAD~40 COUNT=1` meldete 39 Kandidaten, mit
`diff.mnemonicPrefix=true` in der Git-Konfiguration 0 bei Exit 0; der Tabellentest hatte keinen
Fall für die Naht, weil der Diff-Strom dort als Literal kam. Die Fixrunde pinnte die Form des
Stroms, ließ das Programm bei fremdem Präfix mit Exit 2 enden und band den Strom im Test des
Aufrufers an eine fremde Git-Konfiguration; der Verifier maß den Lauf danach mit und ohne
`diff.mnemonicPrefix=true` (je 39) und sah die Mutation „Präfix-Pinnung entfernt“ rot.

Quelle: Review-Report `review-slice-code-kommentare-kennungen` (F-1) ·
Verifikations-Report `verifikation-slice-code-kommentare-kennungen` (§3, §4 S3, §7).
