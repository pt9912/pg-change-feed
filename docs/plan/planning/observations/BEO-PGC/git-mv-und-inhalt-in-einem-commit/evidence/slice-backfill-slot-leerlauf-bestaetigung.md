**Vorgang:** slice-backfill-slot-leerlauf-bestaetigung (Review F-3, Verifikation §5 F-3)

**Fund:** Variante „Ersatz unter neuem Dateinamen“: Commit `00caec48` ersetzte `internal/bootstrap/walretention_endtoend_test.go` (Paket `bootstrap_test`, 249 Zeilen) durch `internal/bootstrap/walretention_endtoend_internal_test.go` (Paket `bootstrap`, 150 Zeilen) in einem Schritt; Git führt dort `A` und `D` statt `R`, `git log --follow` reißt an dieser Stelle. Der Reviewer stufte den Fund als LOW ein (kein Move im engen Sinn, ein Ersatz). Die Fixrunde trennt die spätere Umbenennung auf den Namen des Gegenstands sauber: `ec7e43dd` ist ein reiner Rename (100 %, 0 Einfügungen/0 Löschungen), der Inhalt folgt in `5a5d3422` (der Verifier las beide Commits). Der Bruch bei `00caec48` bleibt in der Historie.

Quelle: `docs/reviews/review-slice-backfill-slot-leerlauf-bestaetigung.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-slot-leerlauf-bestaetigung.md` (§5 F-3). <!-- d-check:status-provenance -->
