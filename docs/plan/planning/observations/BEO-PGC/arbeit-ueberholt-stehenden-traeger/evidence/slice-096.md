# Beleg: slice-096

Vorgang: `slice-096` — Konfigurationsdatei-Nachzug.

Fund: **Vier Ankerstellen, die eine Umbenennung ungültig gemacht hat — und die
Regel, die sie finden sollte, hat nur eine gefunden.**

Der Implementer hat `forbiddenFileDSNKeys` → `forbiddenFileCredentialKeys`
umbenannt und dabei den **`grep`-Lauf nach der bewegten Eigenschaft** geführt, den
`AGENTS.md` §3.13 seit der `welle-20`-Closure verlangt. Er fand den
**Symbolnamen** in `ADR-0088` §Bezug und meldete ihn — richtig, und die Regel
sagt ausdrücklich: einen fremden Träger **melden**, nicht still mitändern.

**Was der Lauf nicht fand, fand der Review:** **zwei Zeilen-Lokatoren** in
`ADR-0088` §Kontext (`wiring.go:276-280` → heute `282-286`,
`config_file.go:148-193` → heute `167-222`). Der Suchlauf sucht die bewegte
**Eigenschaft** — und findet **Symbolnamen**; **Zahlen** fallen durch, obwohl sie
in dieselbe Klasse gehören (`ADR-0073`: ein Zeilen-Lokator bei unverändertem
Referenten ist Zitat-Korrektur). Den **vierten** Anker fand der Architect beim
Nachziehen: §Kontext (4).

**Und der Lauf traf die Träger-Aussagen selbst:** `ADR-0088` nannte die
Zugangsdaten-Klasse „fünf" statt sechs, `ADR-0089`s Ersatztext war für die
Feldmenge zu weit — beide mussten als Folge-ADR (`0091`, `0092`) korrigiert
werden. Dieselbe Eigenschaft, die die Umbenennung bewegte, steckte auch in den
Zahlen der Dokumente, die sie beschreiben.

**Die Klasse, vierter Vorgang.** Die ersten drei (`slice-091`, `-093`, `-094`)
trafen je Träger in einem Dokument; dieser trifft **vier** Stellen in **zwei**
ADRs und dazu zwei falsche Zahlen in eben diesen ADRs.

Quelle: `docs/reviews/verify-slice-096.md` (V-5) ·
`docs/reviews/review-slice-096.md` (F-2) ·
`docs/plan/adr/0091-zugangsdaten-klasse-sechs-schluessel.md` ·
`docs/plan/adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md` ·
`AGENTS.md` §3.13.
