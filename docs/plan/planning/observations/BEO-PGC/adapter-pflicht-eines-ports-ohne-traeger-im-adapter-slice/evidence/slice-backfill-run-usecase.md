**Vorgang:** slice-backfill-run-usecase (Review-Funde F-4, F-5)

**Fund:** Der Use Case ruft Rollback, `Finish` und das Schließen des Snapshots auf einem vom Abbruch gelösten Kontext (`context.WithoutCancel`) ohne Zeitgrenze; „die Dauer begrenzt der Adapter" stand allein im Plan dieses Slice (F-4, MEDIUM). Ebenso ließ die Port-Doku offen, was `Finish` für einen bereits beendeten Run liefert, obwohl der Use Case bei einem Commit mit unbekanntem Ausgang `Finish(failed)` ruft und auf einen wirkungslosen Erfolg baut (F-5, LOW). Der Snapshot-Adapter begrenzte seine Schließ-Dauer bereits selbst, die beiden anderen Adapter existieren erst mit `slice-backfill-run-store`. Behoben: die Pflichten stehen als Port-Vertrag in den Doc-Kommentaren und als DoD-Punkt „Adapter-Pflichten" samt §3-Zeile im Plan von `slice-backfill-run-store`; der Verifier bestätigte den Träger (+16 Zeilen, 0 gelöscht).

Quelle: `docs/reviews/review-slice-backfill-run-usecase.md` (F-4, F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-run-usecase.md` (§3, Zeile „Port-Verträge und ihr Träger"). <!-- d-check:status-provenance -->
