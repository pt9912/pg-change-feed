**Vorgang:** slice-meldungscodes-kennungsfreie-ausgaben (Review F-1; Verifikation §3 M1, §4)

**Fund:** Die Zusage „eine Shell-Kommentarzeile mit `echo` und Kennung ist keine Ausgabe und
trifft nicht“ (die Ausnahme `skip_sh` in `tools/harness/ausgabe-kennungen-check.sh`) war im
Tabellentest ohne Eingabe gebunden, die sie ausübt: der einzige Fall zur Shell-Kommentar-Ausnahme
trug kein `echo`/`printf` und wurde schon vom Shell-Muster nicht getroffen. Der Reviewer setzte
`skip_sh` an einer Kopie im Scratchpad auf ein nie passendes Muster — der Tabellentest blieb bei
43 von 43 Fällen grün; die Go-Seite (`skip_go`) war gebunden (gleiche Mutation färbte den Test
rot). Die Fixrunde `1cead04f` ergänzte Fälle je Skript-Wurzel (`examples`, `tools/schema`),
eingerückt und nicht eingerückt; der Verifier fuhr dieselbe Mutation erneut — genau vier Fälle
rot, `skip_go` rot, das Gate-Skript selbst in der Fixrunde byte-gleich. Schwere HIGH,
Träger-Typ Wächter-Skript mit Tabellentest: eine Datei trotz Deckel. Ausprägung: bei einem
Wächter mit mehreren Ausnahme-Zweigen ist die Eingabeseite **jedes Zweigs** eine Zeile, die
ohne den Zweig träfe.

Quelle: `docs/reviews/review-slice-meldungscodes-kennungsfreie-ausgaben.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-meldungscodes-kennungsfreie-ausgaben.md` (§3 M1, §4). <!-- d-check:status-provenance -->
