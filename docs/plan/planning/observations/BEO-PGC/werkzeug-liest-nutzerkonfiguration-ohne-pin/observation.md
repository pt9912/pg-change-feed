# BEO-PGC/werkzeug-liest-nutzerkonfiguration-ohne-pin

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Harness-Werkzeuge, die die
Ausgabe eines Fremdwerkzeugs parsen, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Werkzeug ruft ein Fremdwerkzeug auf (`git diff`) und parst dessen Ausgabe mit
einer festen Form (Zieldatei-Zeile `+++ b/<Pfad>`). Die Form hängt an der **Konfiguration des
Aufrufers** (`diff.mnemonicPrefix`, `diff.noprefix`, `diff.external`, `color.diff`): unter einer
gängigen Einstellung liest das Werkzeug keinen Pfad, findet keinen Kandidaten und endet mit
Exit 0 — **grün, ohne etwas gelesen zu haben**. Gemessen im Review: derselbe Lauf meldete
39 Kandidaten, mit `diff.mnemonicPrefix=true` **0** bei Exit 0.

**Warum das schwer zu sehen ist:** Der Tabellentest des Programms bekam den Diff-Strom als
Literal; die Naht zwischen `git` und Programm — die Stelle, an der die Konfiguration wirkt — lag
außerhalb jedes Falls. Ein grüner Lauf und ein Lauf, der nichts las, sind an der Ausgabe nicht zu
unterscheiden.

**Abgrenzung.** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` betrifft eine Zusage, die
kein Test an ihre Eingabeseite bindet; hier ist die Eingabeseite eine **Konfiguration außerhalb
des Repos**, die kein Fall im Repo mitbringt. `BEO-PGC/regel-weiter-als-ihr-sensor` betrifft das
Muster, das den Träger nicht trifft; hier trifft das Werkzeug den Träger nicht, weil sein Strom
anders aussieht. `BEO-PGC/werkzeug-fuehrt-plan-inhalt-als-argument-aus` betrifft Text aus einem
Dokument als Kommando-Argument; `make kommentar-kennungen` nimmt seine Argumente aus
Make-Variablen, nicht aus Plan-Text, und bindet die Argument-Zerlegung an eine Marker-Datei.

**Die Antwort ist eine Pinnung am Werkzeug:** der Aufrufer fährt `git` mit gepinnter Form des
Stroms (`--src-prefix=a/ --dst-prefix=b/ --no-ext-diff --no-textconv --no-color`,
`-c core.quotePath=false`), das Programm endet bei einer Zieldatei-Zeile ohne Präfix `b/` mit
Exit 2 statt „kein Kandidat“, und der Test des Aufrufers fährt den Strom unter fremder
Git-Konfiguration (`tools/harness/kommentar-kennungen.sh`,
`tools/harness/run-kommentar-kennungen-tests.sh`).

Deklaration: `slice-code-kommentare-kennungen`, Review F-1 (MEDIUM).
