# BEO-PGC/adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis zwischen
einem Use-Case-Slice und dem Slice, der den Adapter seiner Ports baut, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Use Case setzt am Port eine **Pflicht des Adapters** voraus — hier
die Zeitbegrenzung von `Finish`, `InterruptRunning`, `Rollback` und `Close` unter einem vom
Abbruch gelösten Kontext, und `Finish` auf einen beendeten Run als wirkungsloser Erfolg —,
und die Pflicht steht nur im Plan des Aufrufers. Weder die Port-Doku noch der Plan des
Slice, der den Adapter baut, nennt sie. Der Adapter-Slice hat damit keinen Anlass, sie zu
liefern, und kein Gate liest sie; die Pflicht bleibt liegen, bis ein Reviewer sie beim Lesen
des Aufrufers findet.

**Abgrenzung.** `BEO-PGC/adr-folgepflicht-ohne-traeger-slice` betrifft eine Pflicht, die eine
`Accepted`-ADR benennt, ohne ihr eine Kennung zuzuweisen; hier steht die Pflicht in einem
Slice-Plan und der Empfänger ist ein bekannter Folge-Slice.
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` betrifft eine Adresse, die die Sendung nicht
annimmt; hier fehlt die Adresse ganz.

**Träger, die die Pflicht tragen:** der Doc-Kommentar am Port (der Adapter-Autor liest ihn beim
Implementieren) **und** ein DoD-Punkt im Plan des Adapter-Slice mit Beleg (der Reviewer prüft
ihn); beide zusammen, nicht einer von beiden.

Deklaration: `slice-backfill-run-usecase` (Review F-4, F-5).
