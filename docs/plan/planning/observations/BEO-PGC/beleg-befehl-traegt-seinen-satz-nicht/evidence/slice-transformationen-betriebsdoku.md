**Vorgang:** slice-transformationen-betriebsdoku (Review F-1, HIGH)

**Fund:** Der neue Handbuch-Abschnitt zitierte als „gemessen im Review-Report“ (ein Link auf
den Review-Report des Slice `antragsweg-schema`) den Wortlaut
`` `function cdc.set_transformation(text, text, text, text, jsonb) does not exist` ``. Der
zitierte Report misst an der genannten Stelle einen anderen Wortlaut: die vier
String-Literal-Parameter tragen dort den Typ `unknown` (PostgreSQL typisiert nicht gecastete
Literale beim Signatur-Fehler so, nicht als `text`), der reale Satz lautet „function
cdc.set_transformation(unknown, unknown, unknown, unknown, jsonb) does not exist“. Eine
repoweite Suche fand die im Handbuch verwendete Zeichenkette nirgends sonst — auch die als
zweite Quelle genannte `ADR-0125` Festlegung 1 führt an dieser Stelle nur die abgekürzte Form.
Ein Betreiber, der den realen Fehler erhält und den zitierten Wortlaut zum Abgleich heranzieht,
sieht einen anderen Text als die Handbuch-Zeile behauptet.

**Form (Ausprägung):** die Form **Zitat/Befehl** wie bei den bisherigen Belegen dieser Klasse:
eine Aussage nennt einen konkreten, ausführbaren/nachschlagbaren Beleg (hier: eine
Zeichenkette aus einem realen SQL-Fehler, gemessen in einem genannten Review-Report), und der
genannte Beleg trägt den zitierten Wortlaut nicht — die Aussage „gemessen im Review-Report“
war korrekt in der Herkunfts-Kennzeichnung (`AGENTS.md` §3.12),
falsch im übernommenen Zeichen-für-Zeichen-Wortlaut selbst. Schwere HIGH (der Reviewer-Skill
führt „Beleg trägt seinen Satz nicht“ unter HIGH), daher eine Datei trotz Deckel bei 14×; vor
dem Merge vom Reviewer gefunden, in der Fixrunde behoben (Wortlaut auf `unknown` korrigiert)
und vom Verifier Zeichen-für-Zeichen gegen die Quelle nachgemessen. Ausgang bleibt
**verkörpert** — kein neuer Sensor nötig, die Reviewer-Skill-Probe („den Befehl ausführen, den
Beleg lesen“) hat den Fund getragen.

Quelle: `docs/reviews/review-slice-transformationen-betriebsdoku.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-transformationen-betriebsdoku.md` (§3, §4). <!-- d-check:status-provenance -->
