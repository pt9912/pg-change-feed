# `make zitat-vergleich` / `make test-zitat-vergleich` — Referent-Messung einer Zitat-Korrektur

## Vertrag

`make zitat-vergleich ARGS="…"` vergleicht die **Einheit**, die eine alte
Adresse an ihrem Stand adressiert, mit der Einheit der neuen Adresse an deren
Stand: roh (byte-gleich, abschließende Leerzeilen eingeschlossen) oder, mit
Tag-Paar, nach Normalisierung nur des bewegten Baseline-Tags
(`tools/harness/zitat-vergleich.sh`; netzlos, kein Docker). Es ist die Messung,
die eine Zitat-Korrektur an einer `Accepted` ADR und an einem MR-Pin belegt
([`AGENTS.md`](../../AGENTS.md) §3.5).

**Wer was trägt.** Die **Semantik** — welche Einheit ein Verweis hat, was
normalisiert wird, was „roh“ heißt, wann der Lauf mit 2 endet — trägt
[`ADR-0159`](../../docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(Entscheidung 1 bis 5, mit den Teilen von
[`ADR-0158`](../../docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
die sie in Kraft lässt), ausgelegt durch das
[Architect-Verdikt](../../docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md) §3.
Die **Messung** trägt dieses Ziel; Beleg ist seine gedruckte Zeile. Der
`bash`-Block in Entscheidung 4 der ADR ist die historische Fassung der
Messform. Weicht dieses Werkzeug außerhalb der Punkte unter §Abweichungen von
ihm ab, ist das ein Fehler des Werkzeugs; weicht es von der Semantik ab, ist
das Werkzeug falsch.

**Kein Gate.** Die Messung fällt je Baseline-Bump und je Zitat-Korrektur einmal
an und hat keinen stehenden Stand; ob die Korrektur nur das Gerüst berührt,
bleibt Urteil am Diff. Das Ziel steht in keinem Gate-Bündel (`make gates`);
eine Aufnahme braucht eine ADR ([`AGENTS.md`](../../AGENTS.md) §3.6).

## Aufruf

```text
make zitat-vergleich ARGS="<alt-stand> <alt-pfad> '<alt-ref>' <neu-stand> <neu-pfad> '<neu-ref>' [<alt-tag>:<neu-tag>]"
bash tools/harness/zitat-vergleich.sh <alt-stand> <alt-pfad> <alt-ref> <neu-stand> <neu-pfad> <neu-ref> [<alt-tag>:<neu-tag>]
```

- **Stand:** eine Commit-Kennung (oder ein anderer von `git show` lesbarer
  Name); **Pfad:** relativ zur Repo-Wurzel. Das Skript wechselt in die Wurzel
  des Repos, in dem es aufgerufen wird.
- **Referenz:** `''` = ganze Datei · `'#anker'` = Abschnitt hinter einem
  Heading-Slug oder einer HTML-`id` · `L<a>-<b>` = Zeilen `a` bis `b`. Jede
  andere Form endet mit Exit 2.
- **Tag-Paar** (optional, siebtes Argument): `v<X.Y.Z>:v<X.Y.Z>`, nur wenn roh
  fällt und die Korrektur ein Versions-Segment bewegt.
- Genau 6 oder 7 Argumente, sonst Exit 2 mit Gebrauchszeile. Ohne `ARGS` bricht
  das Make-Ziel mit `$(error …)` ab.

**Quotierung in `ARGS`.** `ARGS` wird von der Shell des Rezepts zerlegt. Ein
ungequotetes `#…` beginnt dort einen Kommentar und kürzt die Argumente stumm;
die Zahl der Argumente fällt dann unter 6, und der Lauf endet mit Exit 2. Eine
leere Referenz steht als `''`:

```text
make zitat-vergleich ARGS="<P>~1 .harness/baseline/<alt-tag>/regelwerk/modul-13-quality-gates.md '#guard-haertung' <P> .harness/baseline/<neu-tag>/regelwerk/modul-13-quality-gates.md '#guard-haertung'"
make zitat-vergleich ARGS="<alt-stand> .harness/baseline/<alt-tag>/regelwerk/grundlagen-begriffe.md '' <neu-stand> .harness/baseline/<neu-tag>/regelwerk/grundlagen-begriffe.md '' <alt-tag>:<neu-tag>"
```

`<P>` ist der Pin-Commit eines Baseline-Bumps; die gemessenen Aufrufe am
realen Bump stehen im Plan von `slice-zitat-vergleich-werkzeug` §2.

## Einheit je Referenz

- **Ganze Datei** und **Zeilen** (`sed -n "<a>,<b>p"`): roh, ein Unterschied im
  Schluss-Umbruch der Datei fällt.
- **Heading-Slug:** der Abschnittskörper (ohne die Heading-Zeile) bis vor das
  nächste Heading gleicher oder höherer Ebene. Der Slug ist der Heading-Text
  klein geschrieben, ohne HTML-Tags, ohne Zeichen außer Buchstaben, Ziffern,
  Leerzeichen, `_` und `-`, Leerzeichen zu `-`; ein doppelter Slug bekommt
  `-1`, `-2`, … . Headings in einem Code-Fence zählen nicht. Aufgelöst wird die
  erste Fundstelle in Dateireihenfolge, Slug oder `id`.
- **HTML-`id`** — gelesen wird nur das öffnende Tag `<a id="X">` außerhalb von
  Fence und Inline-Code (vor dem Tag steht kein Backtick); die Einheit folgt
  der Stellung:

  | Stellung | Einheit |
  |---|---|
  | in einer Heading-Zeile | der Abschnittskörper dieses Headings |
  | in einer Zeile ohne Inhalt, die nächste Zeile mit Inhalt ist ein Heading | der Abschnittskörper dieses Headings |
  | in einer Zeile ohne Inhalt, die nächste Zeile mit Inhalt ist kein Heading | der Block ab dieser Zeile bis vor das nächste Heading beliebiger Ebene oder bis zum Dateiende |
  | in einer Tabellenzeile (die Zeile beginnt mit `\|`) | die Tabellenzeile |
  | in einer anderen Zeile mit Text | der Block ab dieser Zeile bis vor das nächste Heading beliebiger Ebene oder bis zum Dateiende |

  Eine **Zeile ohne Inhalt** trägt nach dem Entfernen aller `<a id="…"></a>`
  nur Leerraum. Gestapelte `id`-Zeilen vor einem Heading adressieren deshalb
  alle den Abschnittskörper dieses Headings, vor einem Absatz den Block ab dem
  Absatz. Ein Fence direkt nach einer `id`-Zeile beginnt den Block; seine
  Öffnungszeile gehört zur Einheit.
- **Fence** (CommonMark): öffnet mit 0 bis 3 Leerzeichen Einzug und mindestens
  drei gleichen Zeichen `` ` `` oder `~` (ein Backtick-Fence trägt im
  Info-String keinen Backtick) und schließt nur mit demselben Zeichen in
  mindestens derselben Länge, danach höchstens Leerraum. Ein nicht
  geschlossener Fence reicht bis zum Dateiende.
- **Normalisierung mit Tag-Paar:** auf der alten Seite wird nur das Segment
  `/<alt-tag>/` durch `/<tag>/` ersetzt, auf der neuen nur `/<neu-tag>/`; ein
  fremder Pin mit anderer Version bleibt roh.

## Ausgabe und Ausgänge

Je Lauf eine Zeile auf stdout, gedruckt vor dem Ende:

```text
vergleich roh: <alt-stand>:<alt-pfad><alt-ref> <-> <neu-stand>:<neu-pfad><neu-ref> cmp 0
vergleich norm v6.14.0:v6.14.1: … cmp 1
vergleich: <stand>:<pfad><ref> keine Einheit, Exit 2
```

Bei verschiedenen Einheiten druckt `cmp` davor die Stelle des ersten
Unterschieds; Meldungen von `einheit` (nicht lesbar, Lokator, leere Einheit,
`id` in anderer Form, awk-Fehler) stehen auf stderr.

| Exit | Bedeutung |
|---|---|
| 0 | die Einheiten sind gleich (`cmp 0`) |
| 1 | die Einheiten sind verschieden (`cmp 1`) |
| 2 | eine Seite hat keine Einheit (Datei an dem Stand nicht lesbar, Lokator ungültig, Einheit leer, `id` in anderer Form), das Tag-Paar ist ungültig, die Zahl der Argumente ist nicht 6 oder 7, oder `awk` besteht die Fähigkeitsprobe nicht |

**Die Farbe steht in der Zeile.** Über `make` kommt jeder Exit ungleich 0 als
der Make-eigene Exit `2` an; ein Beleg zitiert deshalb die gedruckte Zeile
(`cmp 0`, `cmp 1`, `… Exit 2`), nicht den Exit von `make`. Wer den Exit des
Skripts braucht, ruft es direkt auf.

## Abweichungen von der Befehlsform im Block der ADR

Gewollt und im Tabellentest gebunden
([Architect-Verdikt](../../docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md) §2 bis §5):

| Punkt | Werkzeug | Block der ADR |
|---|---|---|
| gestapelte `id` | Zeile ohne Inhalt wird übersprungen, die Einheit ist der Abschnitt bzw. Block dahinter | die Einheit der oberen `id` ist die Zeile der unteren |
| `id` in anderer Form (`<a id="X" class="…">`, `<a id="X"/>`) | Exit 2 mit einer Zeile, die die Form nennt, auch wenn ein Heading-Slug `X` trifft | als Zeile mit Text gelesen, die Einheit ist nur die `id`-Zeile |
| Fence | CommonMark: Einzug 0 bis 3, Zeichen und Länge gemerkt | jede Zeile, die am Zeilenanfang mit ```` ``` ```` oder `~~~` beginnt, schaltet um |
| Tag-Paar | nur mit verschiedenen Tags, `.harness/baseline/<alt-tag>` am alten und `.harness/baseline/<neu-tag>` am neuen Stand; ein leeres siebtes Argument endet mit Exit 2 | jedes Paar der Form `v…:v…`; ein leeres siebtes Argument heißt roh |
| Locale | setzt `LC_ALL=C.UTF-8` selbst und prüft `printf 'Ä' \| awk '{print tolower($0), length($0)}'` gegen `ä 1`, sonst Exit 2 | hängt an der Locale des Aufrufers; unter `LC_ALL=C` endet ein Slug mit Umlaut als leere Einheit |

Daneben prüft das Werkzeug die Zahl der Argumente und den Exit von `sed` am
Zeilen-Lokator und von `awk` am Anker; beides endet mit Exit 2.

## Host-Werkzeuge

`bash`, `git`, `awk`, `sed`, `cmp` ([`AGENTS.md`](../../AGENTS.md) §3.1,
Klasse „Host-Werkzeug ohne Installation“); der Tabellentest zusätzlich `env`
und `mktemp`. `awk` muss Multibyte unter `C.UTF-8` können; das prüft die
Fähigkeitsprobe vor der ersten Messung (fail-closed). Gemessen ist das
Werkzeug an GNU Awk 5.2.1 und GNU bash 5.2; GNU Awk wird nicht nach Name
verlangt, und das Skript benutzt keine Funktion, die nur gawk kennt.

## Grenzen

- **Nachbildung, keine Auflösung durch den Renderer.** Slug und Fence sind
  nachgebildet. Steht vor einem HTML-Tag am Heading-Ende ein Leerzeichen,
  endet der nachgebildete Slug auf `-` und trifft den Slug des Renderers
  nicht; die Einheit ist leer, Exit 2 — adressiert wird dann über die `id`.
  Ein Heading mit Einzug, ein Setext-Heading und ein Tab im Fence-Einzug
  werden nicht gelesen.
- **Schluss-Umbruch im Abschnitts- und Block-Modus.** `awk` schreibt jede Zeile
  mit Zeilenumbruch; ein Unterschied allein im Schluss-Umbruch der Datei fällt
  dort nicht auf.
- **Gleicher Körper.** Ein Anker-Wechsel zwischen zwei Abschnitten mit
  byte-gleichem Körper besteht; ob der neue Anker dieselbe Aussage meint,
  bleibt Urteil am Diff.
- **Block-Einheit.** Der Block einer `id` vor einem Absatz reicht bis zum
  nächsten Heading; eine Änderung in einem späteren Absatz desselben Blocks
  fällt auch, wenn der markierte Absatz gleich blieb (fail-closed).
- **`id` in anderer Form nur für den gesuchten Namen.** Eine andere `id` in
  Attributform in einer gestapelten Zeile gilt als Inhalt.
- **NUL-Bytes** trägt eine Shell-Variable nicht; das Werkzeug ist für
  Textdateien gemacht.
- **Fremder Pin mit dem Bump-Tag.** Trägt ein fremder Pin zufällig genau den
  alten oder neuen Baseline-Tag, normalisiert das Werkzeug auch ihn.
- **Was nicht im Werkzeug steht:** die Datei-Schleife des MR-Datei-`cmp` aus
  [`ADR-0157`](../../docs/plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
  Entscheidung 4; sie bleibt die Befehlsform der ADR, und der Verifier liest
  ihre Ausgabe.

## Test

`make test-zitat-vergleich` (`tools/harness/run-zitat-vergleich-tests.sh`,
netzlos, kein Docker) baut ein Wegwerf-Repo im Temp-Verzeichnis mit Probe-Commits
und fährt je Fall einen Aufruf mit Soll-Exit und Muster der Ausgabe; die
Fall-Tabelle läuft dreimal, ohne Shell-Option und mit `nullglob` bzw.
`failglob` (über `BASHOPTS` des Prüflings). Fallgruppen: HTML-`id` je Stellung
(vor Heading, Heading-Zeile, Absatz, Tabellenzeile, nur Inline-Code) ·
fail-closed (verschiedene Einheiten, fehlende Datei, unbekannter Anker links
und rechts, Lokator, Tag-Paar-Form) · roh (abschließende Leerzeilen,
Leerzeile als Lokator, Schluss-Umbruch) · Normalisierung und Tag-Paar (Bump,
fremder Pin, fremdes Paar ohne Baseline-Baum, gleiche Tags, fehlender Baum je
Seite, leeres Paar) · Anker-Wechsel, Lokator, `git mv` · Heading mit
Inline-Code · Fence nach CommonMark (vier Backticks mit einem und mit zwei
inneren Fences, `~~~` in einem Backtick-Fence, Backtick im Info-String, Fence
nach `id`-Zeile) · gestapelte `id` (vor Heading und vor Absatz) · `id` in
anderer Form (Attribut, selbstschließend, neben gleichnamigem Heading-Slug, in
Inline-Code) · Locale (`LC_ALL=C` beim Aufrufer) und Stub-`awk` ohne Multibyte ·
Argumentzahl · roter Vergleich unter `set -euo pipefail` beim Aufrufer (über
`SHELLOPTS`). Der Prüfling ist per `PROG=<Datei>` übersteuerbar
(Mutationsläufe an Kopien).
