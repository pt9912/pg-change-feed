**Vorgang:** slice-068
**Fund:** Der `open→next`-Übergang von `slice-068` brach einen
Markdown-Link mit festem Verzeichnis: das Architect-Verdikt zur
Spaltenausschluss-Dauerhaftigkeit adressierte den Slice-Plan mit
`../plan/planning/open/slice-068-…md`. `make gates` lief real rot (Exit 2,
`docs-check`, Befund `target-missing`); der Exit-Code war korrekt ungepiped
ermittelt und sichtbar, die Folgehandlung (`git commit`) lief trotzdem,
weil sie nicht an ihn konditioniert war — dieselbe Sequenzierungs-Klasse
wie im Beleg `slice-063-blocker`: gemessen richtig, nicht wirksam. Behoben
durch Auflösen des Links in eine Kennungs-Zitierung.
