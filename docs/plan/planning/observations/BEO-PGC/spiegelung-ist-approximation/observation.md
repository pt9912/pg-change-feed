# BEO-PGC/spiegelung-ist-approximation

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis
zwischen einer gate-getragenen Regel und einem außerhalb des Gates
nachgebauten Artefakt, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Artefakt, das eine bereits gate-getragene Regel außerhalb
des Gates nachbaut, beansprucht Fidelität, die es strukturell nicht liefern
kann — es liest einen anderen Text als sein Vorbild: der lokale
`commit-msg`-Hook die **rohe** Message-Datei, das d-check-Modul `commits` die
**bereinigte**. Die daraus abgeleitete Zusage („spiegelt exakt" / „meldet
denselben Verstoß") ist damit messbar unwahr, und das Artefakt gibt in den
Klassen, in denen es blind ist, ein lokales Grün, das das Gate danach
widerruft — die Vorab-Meldung, für die es gebaut wurde, bleibt dort aus.
Belegt an `.githooks/commit-msg` gegen das Modul `commits`: drei laxere
Klassen (Kennung nur in einer `#`-Kommentarzeile · nur hinter der
scissors-Zeile · Struktur-ID auf der Fortsetzungszeile des ersten Absatzes),
am gepinnten Image und an je einem realen Commit mit aktivem Hook gemessen.
