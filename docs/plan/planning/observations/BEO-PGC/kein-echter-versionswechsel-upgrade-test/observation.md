# BEO-PGC/kein-echter-versionswechsel-upgrade-test

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft
`LH-QA-OPS-005`s Testbeleg, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: `LH-QA-OPS-005`s realer E2E-Beleg
([`ADR-0064`](../../../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md),
Supersedes [`ADR-0058`](../../../../adr/0058-testansatz-fuenf-luecken.md)
Entscheidung 3) tauscht die Feed-Container-Instanz real aus
(`$COMPOSE up -d --force-recreate --no-deps`), aber „alt" und „neu" sind
dasselbe `:dev`-Image — kein echter Versionswechsel. Ohne Git-Tags/
Releases (`ADR-0051` Folgepflicht `slice-040`, noch nicht angelegt) fehlt
ein definierter Gegenstand für „Vorgängerversion".

Deklaration: `slice-063` §6, unverändert aus `ADR-0058`/`ADR-0064`
übernommen.
