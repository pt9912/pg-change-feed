# Beleg: welle-d-check-verkoerperung

Vorgang: Delta-Review der Verkörperungs-Fixrunde (Commit `4da1ad8`), die
diesen Registereintrag selbst in `.harness/skills/reviewer.md` einbettet
(Architect-Verdikt `architect-verdict-welle-d-check-lese-schritt.md`).

Fund (Review-Finding, HIGH,
`docs/reviews/review-welle-d-check-verkoerperung.md`): Der neue Grenz-Absatz
in `AGENTS.md` §3.13 zitierte `docs/reviews/review-slice-096.md` F-2 mit
einer Zählung, die das Original nicht trägt — F-2 berichtet **zwei**
Zeilen-Lokatoren, der Satz sprach von „**einen** Zeilen-Lokator". Der
verwendete Begriff „vierter Anker" kam in F-2 gar nicht vor; er stammte aus
`evidence/slice-096.md` (dort dem Architect zugeschrieben), während die
Closure-Notiz von `slice-096` ihn dem Review zuschreibt — zwei Quellen
widersprachen sich.

**Vierter Vorgang, direkt nach Embodiment.** Korrigiert durch Gegenprobe an
der primären, autoritativen Quelle
(`docs/plan/planning/done/slice-096-konfigurationsdatei-nachzug.md` §7,
dritter Punkt: „der §3.13-Lauf fand den Symbolnamen, der Review zwei
Lokatoren und einen vierten Anker, und der Architect-Zug fand die zwei
falschen Zahlen in den ADRs selbst") statt der zweiten Hand
(`evidence/slice-096.md`). Der Fund bestätigt den Wert der gerade
verkörperten Regel unmittelbar: Sie hätte diesen eigenen Fehler gefangen,
wäre sie beim Schreiben angewendet worden.

Quelle: `docs/reviews/review-welle-d-check-verkoerperung.md` ·
`AGENTS.md:451-455` (korrigiert) ·
`docs/plan/planning/done/slice-096-konfigurationsdatei-nachzug.md` §7.
