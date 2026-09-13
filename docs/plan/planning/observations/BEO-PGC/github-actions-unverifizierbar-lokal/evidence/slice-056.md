**Vorgang:** slice-056
**Fund:** Wie schon bei `slice-039` musste `e2e.yml` bis zum ersten realen
Push auf `main` ohne Docker-only-Beleg auskommen (statische Prüfung war
das einzig lokal Mögliche). Der reale Beleg fiel diesmal positiv aus (vier
unabhängige grüne Läufe, von Reviewer und Verifier unabhängig per
`gh run view`/`gh run list` gegengeprüft) — das ändert nichts daran, dass
das strukturelle Verifikationsgrenze-Muster (kein lokaler Beleg vor dem
ersten realen Push möglich) erneut auftrat. Zweites, unabhängiges
Auftreten.
