# Beleg: slice-105

Vorgang: `slice-105` — Baseline von `v6.5.0` auf `v6.9.0` gehoben,
`.harness/baseline/v6.5.0/` per `git rm -r` entfernt. <!-- d-check:ignore (historische Vor-Migrations-Erwähnung, kein lebender Pin) -->

Fund: Die Arbeit bewegt eine beschriebene Eigenschaft — welche
Baseline-Version vendored ist — und überholt damit Aussagen, die das nicht
selbst anfassen: drei Markdown-Links in `docs/reviews/**`
(`architect-verdict-aufschub-adresse-verfaellt.md`,
`review-slice-041-fixrunde-2.md`) zeigten real auf
`.harness/baseline/v6.5.0/...` und wurden durch das Entfernen zu <!-- d-check:ignore (historische Vor-Migrations-Erwähnung, kein lebender Pin) -->
`target-missing`. Gefunden vom Implementer selbst über den vorgeschriebenen
§3.13-Suchlauf (repo-weiter `grep` auf den bewegten Pfad, nicht nur über den
eigenen Diff), im selben Zug per Zitat-Korrektur nach `ADR-0073` behoben
(Referent unverändert, nur das Versions-Pfadsegment korrigiert). 21 weitere
reine Text-Erwähnungen (kein Link) desselben Pfads in `docs/reviews/**`
wurden über eine neue `.d-check.yml`-`versions.exempt-paths`-Zeile
aufgefangen statt einzeln editiert — Records sind unveränderliche
Lauf-Belege (`AGENTS.md` §3.5 analog). Reviewer und Verifier bestätigen
unabhängig voneinander: kein verbliebener kaputter Link außerhalb der vier
erwarteten Klassen (eigene Plan-Datei, `docs/reviews/**` jetzt exempt,
`done/`-Records, `ADR-0051` P8-Pininventur).

Quelle: `docs/plan/planning/in-progress/slice-105-baseline-v6.9.0-materialisieren.md`
§6 (nachträglich ergänztes viertes Risiko) · `docs/reviews/review-slice-105.md` ·
`docs/reviews/verify-slice-105.md` §7 · Commit `94de7bb`.
