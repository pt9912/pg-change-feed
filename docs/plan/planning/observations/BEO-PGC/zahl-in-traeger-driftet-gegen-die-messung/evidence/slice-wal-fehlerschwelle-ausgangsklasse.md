**Vorgang:** slice-wal-fehlerschwelle-ausgangsklasse (Verifikation V-2, Fixrunde, zweite
Verifikation §3)

**Fund:** Das committete §3.13-Suchlauf-Feld des Plans (DoD-Zeile 6, Symbolnamen-Zeile
`mergeStreamAndWALFaultOutcome|stopStream`) behauptete für den `diff`-Stand „56 Treffer,
9 Zeilen stimmen, Exit 0". Die erste Verifikation maß eigen `ist=59` (V-2, MEDIUM): die
vorige Fixrunde (`a7d27ddc`) hatte einen neuen Testfall
(`TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled`) ergänzt, der die
bewegte Stelle real drei weitere Male nennt, ohne den `diff`-Suchlauf danach erneut zu
fahren und das Plan-Feld nachzutragen. Dieselbe Ursachen-Klasse trat **im selben Slice
ein zweites Mal** ein, diesmal aber korrekt behandelt: die Fixrunde, die V-1 (ungetesteter
Guard) behob, ergänzte einen weiteren Testfall
(`TestMergeStreamAndWALFaultOutcomeLeavesStreamErrorUnchangedWithoutWALFault`, zwei
weitere Treffer) und hat den `diff`-Suchlauf diesmal **vor** der eigenen Meldung erneut
gefahren und auf `61` korrigiert (Plan §3, „Fixrunden-Nachmessung des `diff`-Werts") — die
zweite Verifikation bestätigte `61` unabhängig sowohl über `make suchlauf-nachmessen` als
auch über eine vom Werkzeug unabhängige `git grep`-Zählung (§3 dort). Die Lehre: eine
Fixrunde, die einen neuen Testfall auf die bewegte Stelle selbst schreibt, zieht das
committete Suchlauf-Feld **im selben Commit** nach, statt es der nächsten Rolle zu
überlassen — die erste Fixrunde tat das nicht, die zweite tat es.

Quelle: `docs/reviews/verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md` (V-2, §6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-wal-fehlerschwelle-ausgangsklasse-2.md` (§3). <!-- d-check:status-provenance -->
