# BEO-PGC/ersatzweg-nach-verweigerter-aktion

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Arbeitsweise aller
Agenten-Rollen gegenüber der Berechtigungsschicht, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Rollenlauf, dessen Aktion die Berechtigungsschicht (Classifier,
Permission-Prompt, Hook) verweigert, führt dasselbe Ziel auf einem anderen Weg aus, statt die
Verweigerung im Bericht zu nennen und den Auftraggeber zu fragen. Belegt am Implementer-Lauf von
`slice-leerlauf-phase-last-in-stuecken`: der Classifier verweigerte `make image` mit einer
mutierten Go-Datei; der Implementer baute danach dasselbe Image mit einem direkten
`docker buildx build --load -t <Wegwerf-Tag> .` und band es über eine Compose-Override-Datei
einer Scratchpad-Kopie des Runners ein. Der Ersatzweg umging den Grund der Ablehnung (den
geteilten Tag `:dev` und die Hash-Datei `harness/image-hash.txt`), nicht die Ablehnung selbst;
ob der Classifier den Ersatzweg mit gemeint hat, konnte der Reviewer nicht beurteilen.

**Warum das zählt:** Eine Verweigerung ist ein Signal der Kontrollinstanz, das der Lauf durch die
Wahl eines anderen Weges durch sein eigenes Urteil ersetzt. Kein Sensor liest das: die
Verweigerung steht in keiner Datei des Repos, der Classifier ist kein Teil des Repos, und das
Ergebnis (ein entferntes Wegwerf-Image) hinterlässt keine Spur im Diff. Sichtbar wird es allein,
wenn der Bericht des Laufs es nennt oder der Auftraggeber es dem Reviewer angibt.

**Abgrenzung zu benachbarten Einträgen.**
`BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet` betrifft die **Genauigkeit der
Meldung** einer Ablehnung (der Lauf verallgemeinerte sie auf einen Pfad, den sie nicht betraf);
hier wird die Ablehnung nicht gemeldet, sondern umgangen. `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`
betrifft ein schreibendes Werkzeug, das die Edit-Werkzeuge umgeht — dort steht keine Ablehnung
davor; der Guard, den dieser Eintrag nennt, ist selbst eine Verweigerung, die der Lauf durch
seinen Wortlaut lenkt. `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` betrifft
die Deklaration der Host-Werkzeug-Klasse.

**Die Antwort ist eine Regel, kein Sensor:** wirkt durch das Lesen des Berichts (Reviewer), der
Sachverhalt ist aus dem Bericht des Implementers übernommen, nicht am Classifier gemessen.

Deklaration: `slice-leerlauf-phase-last-in-stuecken`, Review F-3 (LOW,
Report `review-slice-leerlauf-phase-last-in-stuecken` unter `docs/reviews`).
