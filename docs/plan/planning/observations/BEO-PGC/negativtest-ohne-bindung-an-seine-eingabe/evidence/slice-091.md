# Beleg: slice-091

Vorgang: `slice-091` — Coverage Cluster C (Zustell- und Betriebs-Rand).

Fund: Im Verifikationsbericht zu `slice-091` **V-7** blieb ein **sechster** Test
derselben Form grün, wenn man seine Zusage zerstört:
`internal/adapters/driving/http/server_test.go:161`,
`TestRegisterConsumerUngueltigesJSONEndetMit400`. Er prüft, dass ein nicht
dekodierbarer Body mit `400` endet — aber mit abgeschaltetem Decode-Zweig
(`registerconsumer.go:38`, `err != nil` → `false && err != nil`) liefert er
weiterhin `400`, weil er die Ablehnung nicht an seinen Eingabewert bindet.

**Er ist nicht im Diff dieses Slice** — die fünf Geschwister
`…UngueltigesJSONEndetMit400`, die der Slice hinzufügt, binden (je durch Mutation
belegt); dieser eine war schon da. Nach der Bindung (vierte Runde, `W-3`) färbt
dieselbe Mutation ihn rot (`Status: 500 (Erwartung: 400 …)`, EC=1), Kontrolle
grün.

**Ursprung ≠ Vorkommen.** Eingeführt hat ihn `7b6b253` (`slice-061`, das
HTTP-Adapter-Grundgerüst) — das ist der **Ursprung** der Form, kein
Beleg-Vorgang. **Gefunden** wurde er im Verifikationsbericht zu `slice-091`, und das ist das
Vorkommen, das dieser Beleg dateit. Die Klasse trifft damit zum fünften Mal; dass
sie **im selben Zug** an einem neuen Test (F-1, dort als
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` geführt) und an diesem
vorbestehenden auftrat, ist der Grund, warum der Review sie als
„wahrscheinlichster Fehler dieses Slice" vorhergesagt hatte — er hatte recht.

**Warum das zählt:** Der Test ist grün und sieht aus wie seine fünf gebundenen
Geschwister; nur die Mutation trennt sie. Ein Bestand, der so aussieht wie die
neu geschriebene Hälfte, wird bei der nächsten Prüfung mitgezählt, ohne geprüft
zu werden.

Quelle: Verifikationsbericht zu `slice-091` (V-7) ·
Delta-Review zu `slice-091` (Negativbefunde, W-3) ·
`internal/adapters/driving/http/server_test.go:161` (berichtigt in `a7d7f7b`) ·
`git log --diff-filter=A -S 'TestRegisterConsumerUngueltigesJSONEndetMit400'`
(`7b6b253`, Ursprung).
