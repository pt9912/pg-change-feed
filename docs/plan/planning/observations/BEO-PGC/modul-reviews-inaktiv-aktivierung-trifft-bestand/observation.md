# BEO-PGC/modul-reviews-inaktiv-aktivierung-trifft-bestand

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis zwischen dem
gepinnten Gate-Werkzeug d-check und einem Modul, das in der Modul-Liste von `.d-check.yml` und in
jedem `--enable` von `d-check.mk` fehlt, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Pin-Bump hebt ein Gate-Werkzeug, dessen Changelog eine „nicht rein additive“
Änderung an einem Modul nennt, das in diesem Repo **nicht läuft**. Das Reichweiten-Urteil „keine
Reichweiten-Änderung eines Ziels“ ist wahr, trägt aber die latente Folge nicht: Wer das Modul
später aktiviert, trifft die Änderung am gewachsenen Bestand, ohne dass ein Lauf sie je gezeigt
hat. Der Plan des Bumps setzte für dieses Risiko vorab den Ausgang *entfallen* („läuft heute nicht
mit“); die Messung widerlegte ihn.

**Erster Vorgang:** `slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1` (Review F-1, MEDIUM). Das
Modul `reviews` (Review-Report-Deckung für `done/`-Slices mit Review-DoD-Haken; Soll-Kante
Code→Review in `harness/README.md`, Ziel `doc-reviews` existiert nicht) meldet am Bestand
`--enable reviews`, gemessen am Stand `d496d0ee`:

| d-check | Exit | Befunde |
|---|---|---|
| 0.82.0 | 0 | 0 |
| 0.84.0 | 0 | 0 |
| 0.85.0 | 1 | **132** × `review-missing` |
| 0.86.1 | 1 | **132** × `review-missing` |

Die 132 gehen damit auf 0.85.0 zurück (zwei Defaults von `reviews` ändern sich); der Anteil der
Änderung an `reviews.match` in 0.86.0 ist nicht gemessen.

**Abgrenzung:** `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` beschreibt ein abgeschaltetes Modul,
dessen Regel das Repo ausdrücklich will und die der Bestand verletzt (dort `hostpaths`, aktiviert);
hier ist das Modul ebenfalls aus, aber die **Werkzeug-Version** hat die Folge der Aktivierung
verschoben, und es gibt keinen Beschluss, es zu aktivieren.
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` beschreibt eine Erweiterung eines **aktiven**
Moduls.

**Warum das zählt:** Aktiviert jemand `reviews` später, ist die Aktivierung ein Vorgang mit 132
Befunden am Bestand statt eines grünen Einschaltens. Die Zahl steht hier, weil der
Reichweiten-Abschnitt des Plans mit dem Slice in `done/` liegt.
