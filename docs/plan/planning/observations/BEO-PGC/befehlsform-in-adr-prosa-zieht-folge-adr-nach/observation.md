# BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den Ort, an dem
eine Messvorschrift als ausführbare Form steht, keine eigene Sub-Area im Sinn
der Modus-Deklaration).

Die Beobachtung: Eine ADR trägt eine **ausführbare Befehlsform** (einen
`bash`-Block) als Verfahrensregel, und wer misst, zieht sie per `awk` aus der
ADR. Jeder Fehler der Form — ein Randfall, eine Shell-Option, eine Locale —
ist dann ein Fehler im Text einer `Accepted`-ADR. Eine Berichtigung geht nur
über eine Folge-ADR (`AGENTS.md` §3.5), und die nächste Prüfung der neuen Form
findet den nächsten Randfall. Die Kette Folge-ADR auf Folge-ADR entsteht nicht
aus einer falschen Entscheidung, sondern aus dem Träger: Code mit Randfällen
in einem unveränderlichen Dokument, ohne Tabellentest.

Konkret (Erstauftreten): `ADR-0158` Entscheidung 6 (Befehlsform) war unter
`nullglob`/`failglob` fail-open und verlor abschließende Leerzeilen (Review
F-3, MEDIUM; F-5, F-6, LOW); `ADR-0159` Entscheidung 4 ersetzte sie, und das
Re-Review fand drei weitere Formen, die sie nicht trägt (F-1 bis F-3, LOW)
sowie zwei unbenannte Ränder (F-5, F-6, INFO). `ADR-0159` sieht dafür
Re-Evaluierungs-Trigger (a) vor: Option E, ein Skript hinter `make`, statt
einer weiteren Folge-ADR.

Benannt, nicht gezählt: `ADR-0157` Entscheidung 4 trägt eine Schleife (MR-Datei-`cmp`),
deren Exit nicht gefärbt ist (Re-Review F-5 zu `slice-harness-baseline-v6-14-1`,
INFO); berichtigt wurde sie nicht, eine Folge-ADR entstand daraus nicht. Ältere
Folge-ADR-Ketten (`ADR-0062` → `ADR-0069` → `ADR-0070`, `ADR-0104` →
`ADR-0151` → `ADR-0155`) betreffen keine Befehlsform; sie sind eine andere
Ursache und hier nicht gezählt.
