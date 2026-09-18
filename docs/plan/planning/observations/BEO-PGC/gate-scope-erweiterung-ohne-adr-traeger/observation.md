# BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Scope-Deklaration der Gates, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Gate-Scope wächst, weil ein Slice eine Ausnahme
braucht — und zwar **durch die Konfigurationsdatei selbst** (`.a-check.yml`),
also ohne den Träger, den die zuständige ADR für Änderungen an der
Architekturregel benennt (`AGENTS.md` §3.6: ADR, kein PR-Kommentar). Die
Ausnahme fällt dabei nicht als Regeländerung auf, weil sie wie der
normale Pflege-Vorgang „deklarativer Stand nachziehen" aussieht; der
Unterschied — **Verfeinerung** (ADR-frei) gegen **Erweiterung**
(ADR-pflichtig) — war vor `ADR-0068` nirgends gezogen.

Deklaration: `slice-071` (Review-Finding F-1,
Review zu `slice-071`; entschieden über
`ADR-0068`, `Supersedes ADR-0041` nur die Änderungs-Ausnahmeklausel).
Verwandt, aber Gegenrichtung: `BEO-PGC/a-check-null-abdeckung` (Dateien in
**keiner** Schicht) — hier ist es eine Gruppe, die **zu weit** reicht.
