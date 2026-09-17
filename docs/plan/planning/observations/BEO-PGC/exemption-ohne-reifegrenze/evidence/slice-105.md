# Beleg: slice-105

Vorgang: `slice-105` — Baseline `v6.9.0` materialisiert, `.harness/baseline/v6.5.0/` <!-- d-check:ignore (historische Vor-Migrations-Erwähnung, kein lebender Pin) -->
entfernt, `.d-check.yml` `versions.exempt-paths` um `docs/reviews/**`
erweitert.

Fund: Die neue Ausnahme wird mit „dieselbe Begründung wie
`harness/conventions/done/**`" gerechtfertigt, ist aber breiter: Eine
Adaptions-Datei unter `conventions/done/**` wird erst durch ihren `git mv`
fixiert und bis dahin aktiv geprüft; `docs/reviews/**` ist ab dem ersten
Commit, der eine Datei dort anlegt, von der `versions`-Prüfung ausgenommen —
ohne Übergabe-Schwelle. Real reproduziert (Reviewer F-1, vom Verifier
unabhängig nachvollzogen): ein testweise unter `docs/reviews/**` eingefügter
`v6.5.0`-Pfad bleibt bei `make docs-check` grün (nur andere Module schlagen
ggf. an, `version-stale` nicht), derselbe Pfad außerhalb der Ausnahme
(`harness/README.md`) färbt korrekt rot.

Quelle: `docs/reviews/review-slice-105.md` F-1 (MEDIUM, Klasse „Exemption
ohne Reifegrenze") · `docs/reviews/verify-slice-105.md` §5 (eigene,
unabhängige Begründung, teilt die Einstufung) · `.d-check.yml:61`.
