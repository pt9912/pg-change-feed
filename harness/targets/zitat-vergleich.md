# `make zitat-vergleich` / `make test-zitat-vergleich` — Referent-Messung einer Zitat-Korrektur

## Vertrag

`make zitat-vergleich ARGS="…"` vergleicht die **Einheit**, die eine alte
Adresse an ihrem Stand adressiert, mit der Einheit der neuen Adresse an deren
Stand: roh (byte-gleich, abschließende Leerzeilen eingeschlossen) oder, mit
Tag-Paar, nach Normalisierung nur des bewegten Baseline-Tags
(`tools/harness/zitat-vergleich.sh`; netzlos, kein Docker). Es ist die Messung,
die eine Zitat-Korrektur an einer `Accepted` ADR belegt
([`AGENTS.md`](../../AGENTS.md) §3.5) und im Adaptions-Durchgang eines
Baseline-Bumps den Referenten je MR-Eintrag zwischen den Tags des Bumps misst
([`ADR-0161`](../../docs/plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
Entscheidung 5).

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
  Heading-Slug oder einer HTML-`id` · `L<a>-<b>` = Zeilen `a` bis `b` mit
  `1 <= a <= b`. Jede andere Form endet mit Exit 2 (auch `L7`, `L5-3`, `L0-1`,
  `L1-2-3`).
- **Tag-Paar** (optional, siebtes Argument): `v<X.Y.Z>:v<X.Y.Z>`, nur wenn roh
  fällt und die Korrektur ein Versions-Segment bewegt.
- Genau 6 oder 7 Argumente, sonst Exit 2 mit Gebrauchszeile. Ohne `ARGS` bricht
  das Make-Ziel mit `$(error …)` ab.

**Quotierung in `ARGS`.** `ARGS` wird von der Shell des Rezepts zerlegt. Ein
ungequotetes `#…` beginnt dort einen Kommentar und kürzt die Argumente stumm;
die Zahl der Argumente fällt dann unter 6, und der Lauf endet mit Exit 2. Eine
leere Referenz steht als `''`:

```text
make zitat-vergleich ARGS="<stand> .harness/baseline/<alt-tag>/regelwerk/modul-13-quality-gates.md '#guard-haertung' <stand> .harness/baseline/<neu-tag>/regelwerk/modul-13-quality-gates.md '#guard-haertung'"
make zitat-vergleich ARGS="<alt-stand> .harness/baseline/<alt-tag>/regelwerk/grundlagen-begriffe.md '' <neu-stand> .harness/baseline/<neu-tag>/regelwerk/grundlagen-begriffe.md '' <alt-tag>:<neu-tag>"
```

`<stand>` ist ein Stand eines Baseline-Bumps, an dem beide Tag-Verzeichnisse
im Baum liegen (Adaptions-Durchgang, `ADR-0161` Entscheidung 5); gemessene
Aufrufe stehen im Plan von `slice-zitat-vergleich-werkzeug` §2 (dort am
Pin-Commit gemessen) und von `slice-harness-baseline-v6-16-0`.

## Einheit je Referenz

- **Ganze Datei** und **Zeilen** (`sed -n "<a>,<b>p"`): roh, ein Unterschied im
  Schluss-Umbruch der Datei fällt.
- **Heading-Slug:** der Abschnittskörper (ohne die Heading-Zeile) bis vor das
  nächste Heading gleicher oder höherer Ebene. Ein Heading beginnt am
  Zeilenanfang mit 1 bis 6 `#` und einem Leerzeichen. Der Slug ist der
  Heading-Text ohne schließende `#`-Folge (CommonMark: ` ##` am Zeilenende),
  klein geschrieben, ohne HTML-Tags, ohne Zeichen außer Buchstaben, Ziffern,
  Leerzeichen, `_` und `-`, Leerzeichen zu `-`; ein doppelter Slug bekommt
  `-1`, `-2`, … . Headings in einem Code-Fence zählen nicht. Aufgelöst wird die
  erste Fundstelle in Dateireihenfolge, Slug oder `id`.
- **HTML-`id`** — gelesen wird nur das öffnende Tag `<a id="X">` außerhalb von
  Fence, Inline-Code (vor der Fundstelle steht in der Zeile eine ungerade Zahl
  Backticks), eingerücktem Code (die
  Zeile beginnt mit vier Leerzeichen oder einem Tab; das trifft auch einen
  Fence in einem Listenpunkt) und HTML-Kommentar (`<!--` bis `-->`, auch über
  Zeilen; ein `<!--` hinter einer ungeraden Zahl Backticks der Zeile steht in
  Inline-Code und öffnet keinen). Die Einheit folgt der Stellung:

  | Stellung | Einheit |
  |---|---|
  | in einer Heading-Zeile | der Abschnittskörper dieses Headings |
  | in einer Zeile ohne Inhalt, die nächste Zeile mit Inhalt ist ein Heading | der Abschnittskörper dieses Headings |
  | in einer Zeile ohne Inhalt, die nächste Zeile mit Inhalt ist kein Heading | der Block ab dieser Zeile bis vor das nächste Heading beliebiger Ebene oder bis zum Dateiende |
  | in einer Tabellenzeile (die Zeile beginnt mit `\|`) | die Tabellenzeile |
  | in einer anderen Zeile mit Text | der Block ab dieser Zeile bis vor das nächste Heading beliebiger Ebene oder bis zum Dateiende |

  Eine **Zeile ohne Inhalt** trägt nach dem Entfernen aller `<a id="…"></a>`
  und aller HTML-Kommentar-Teile nur Leerraum. Gestapelte `id`-Zeilen vor einem
  Heading adressieren deshalb alle den Abschnittskörper dieses Headings, vor
  einem Absatz den Block ab dem Absatz; eine Zeile `<!-- … --> <a id="X"></a>`
  vor einem Heading ebenso den Abschnitt. Ein Fence direkt nach einer
  `id`-Zeile beginnt den Block; seine Öffnungszeile gehört zur Einheit.

  **Mehrdeutig endet mit Exit 2.** Das Werkzeug parst Code-Spans nicht. Trägt
  eine Zeile außerhalb eines Fence `<a id="X"` (nicht eingerückt) oder eine
  Kommentar-Grenze (`<!--`, `-->`) und dazu einen Backtick-Lauf ab Länge 2,
  eine ungerade Zahl Backticks, oder steht sie in einem Absatz, dessen
  Vorzeilen eine ungerade Zahl Backticks tragen (offener Code-Span), so ist
  nicht sicher, ob `id` oder Kommentar in Inline-Code stehen. Eine solche
  `id`-Zeile und jede Fundstelle von `X` nach einer solchen Kommentar-Grenze
  enden mit Exit 2 („mehrdeutig“); ebenso eine `id`-Zeile mit einem Einzug aus
  Leerzeichen und Tab. Exit 2 heißt hier „nicht messbar, Urteil am Diff“
  ([`ADR-0159`](../../docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
  Entscheidung 2). Einzelne Backticks in gerader Zahl ohne offenen Span aus
  der Vorzeile gelten als sicher gelesen: Eine `id` oder ein `<!--`, vor dem in
  der Zeile eine ungerade Zahl Backticks steht, liegt in Inline-Code und wird
  nicht gelesen, auch mit Text zwischen Backtick und Tag
  (`` `siehe <a id="x"></a>` ``); bei gerader Zahl liegt die Fundstelle
  außerhalb. Ein escapter Backtick trägt diese Lesung nicht (§Grenzen).
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
vergleich norm v6.14.1:v6.16.0: … cmp 1
vergleich: <stand>:<pfad><ref> keine Einheit, Exit 2
```

Bei verschiedenen Einheiten druckt `cmp` davor die Stelle des ersten
Unterschieds; Meldungen von `einheit` (nicht lesbar, Lokator, leere Einheit,
`id` in anderer Form, awk-Fehler) stehen auf stderr.

| Exit | Bedeutung |
|---|---|
| 0 | die Einheiten sind gleich (`cmp 0`) |
| 1 | die Einheiten sind verschieden (`cmp 1`) |
| 2 | eine Seite hat keine Einheit (Datei an dem Stand nicht lesbar, Lokator ungültig, Einheit leer, `id` in anderer Form, `id` in mehrdeutiger Stellung), das Tag-Paar ist ungültig, die Zahl der Argumente ist nicht 6 oder 7, oder `awk` besteht die Fähigkeitsprobe nicht |

**Die Farbe steht in der Zeile.** Über `make` kommt jeder Exit ungleich 0 als
der Make-eigene Exit `2` an; ein Beleg zitiert deshalb die gedruckte Zeile
(`cmp 0`, `cmp 1`, `… Exit 2`), nicht den Exit von `make`. Wer den Exit des
Skripts braucht, ruft es direkt auf.

## Abweichungen von der Befehlsform im Block der ADR

Gewollt und im Tabellentest gebunden
([Architect-Verdikt](../../docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md) §2 bis §5;
die drei letzten Zeilen aus der Review-Fixrunde, als Auslegung im Sinn von
Verdikt §2: dieselbe Erkennungsregel wie „`id` in Code-Fence oder Inline-Code
zählt nicht“, ein strikter Lokator nach „ungültiger Lokator endet mit Exit 2“
und der Slug des Renderers):

| Punkt | Werkzeug | Block der ADR |
|---|---|---|
| gestapelte `id` | Zeile ohne Inhalt wird übersprungen, die Einheit ist der Abschnitt bzw. Block dahinter | die Einheit der oberen `id` ist die Zeile der unteren |
| `id` in anderer Form (`<a id="X" class="…">`, `<a id="X"/>`) | Exit 2 mit einer Zeile, die die Form nennt, auch wenn ein Heading-Slug `X` trifft | als Zeile mit Text gelesen, die Einheit ist nur die `id`-Zeile |
| Fence | CommonMark: Einzug 0 bis 3, Zeichen und Länge gemerkt | jede Zeile, die am Zeilenanfang mit ```` ``` ```` oder `~~~` beginnt, schaltet um |
| Tag-Paar | nur mit verschiedenen Tags, `.harness/baseline/<alt-tag>` am alten und `.harness/baseline/<neu-tag>` am neuen Stand; ein leeres siebtes Argument endet mit Exit 2 | jedes Paar der Form `v…:v…`; ein leeres siebtes Argument heißt roh |
| Locale | setzt `LC_ALL=C.UTF-8` selbst und prüft `printf 'Ä' \| awk '{print tolower($0), length($0)}'` gegen `ä 1`, sonst Exit 2 | hängt an der Locale des Aufrufers; unter `LC_ALL=C` endet ein Slug mit Umlaut als leere Einheit |
| `id` in eingerücktem Code und in HTML-Kommentar | nicht gelesen | gelesen; `harness/conventions.md#mr-<NNN>` löst auf den Kommentar auf |
| Lokator | nur `L<a>-<b>` mit `1 <= a <= b` | `L7` als `L7-7`, `L5-3` als Zeile 5 |
| schließende `#`-Folge eines Headings | gehört nicht zum Slug (`## Eins ##` → `#eins`) | gehört zum Slug (`#eins-`) |
| mehrdeutige Stellung (Backtick-Lauf ab 2, ungerade Backticks, offener Code-Span, Einzug aus Leerzeichen und Tab) | Exit 2, „mehrdeutig“ | gelesen wie jede andere Zeile |
| HTML-Kommentar in einer `id`-Zeile | gilt als ohne Inhalt | gilt als Text, Einheit ist der Block |

Daneben prüft das Werkzeug die Zahl der Argumente und den Exit von `sed` am
Zeilen-Lokator und von `awk` am Anker; beides endet mit Exit 2.

## Host-Werkzeuge

`bash`, `git`, `awk`, `sed`, `cmp` ([`AGENTS.md`](../../AGENTS.md) §3.1,
Klasse „Host-Werkzeug ohne Installation“); der Tabellentest zusätzlich `env`,
`grep`, `mktemp` und `mv` (im Temp-Verzeichnis). `awk` muss Multibyte unter `C.UTF-8` können; das prüft die
Fähigkeitsprobe vor der ersten Messung (fail-closed). Gemessen ist das
Werkzeug an GNU Awk 5.2.1 und GNU bash 5.2; GNU Awk wird nicht nach Name
verlangt, und das Skript benutzt keine Funktion, die nur gawk kennt.

## Grenzen

- **Nachbildung, keine Auflösung durch den Renderer.** Slug und Fence sind
  nachgebildet. Steht vor einem HTML-Tag am Heading-Ende ein Leerzeichen,
  endet der nachgebildete Slug auf `-` und trifft den Slug des Renderers
  nicht; die Einheit ist leer, Exit 2 — adressiert wird dann über die `id`.
  Ein Heading mit Einzug, ein Setext-Heading und ein Tab im Fence-Einzug
  werden nicht gelesen. Ein Heading in einem HTML-Kommentar zählt als Heading
  (nur die `id` wird dort ausgeblendet).
- **Code-Spans werden nicht geparst.** Ob ein `<!--` oder eine `id` in
  Inline-Code steht, liest das Werkzeug nur bei einzelnen Backticks in gerader
  Zahl ohne offenen Span aus der Vorzeile; jede andere Backtick-Lage an einer
  `id`-Zeile oder Kommentar-Grenze endet mit Exit 2 („mehrdeutig“, siehe
  §Einheit). Ein `<!--` in einer eingerückten Code-Zeile öffnet einen
  Kommentar; spätere `id`s bis zum nächsten `-->` sind dann ausgeblendet, die
  Einheit ist leer (Exit 2) oder die eines gleichnamigen Heading-Slugs.
- **`<!--` in eingerücktem Code nach einer `id`-Zeile (fail-open).** Folgt auf
  eine `id`-Zeile ohne Inhalt eingerückter Code, der mit `<!--` beginnt, so gilt
  der Code als Kommentar und damit als „ohne Inhalt“; die Einheit wird der
  Abschnitt des nächsten Headings statt des Code-Blocks, und eine Änderung im
  Code-Block ergibt `cmp 0`. Heute ohne Fundstelle (übernommen aus dem
  Re-Review zu Fixrunde 2); der Messende prüft diese Stellung am Diff.
- **Eingerückter Code in einem Blockzitat (fail-open).** Container werden
  nicht gelesen; die Einzugsregel liest nur den Zeilenanfang. Eine `id` in
  `>     <a id="x"></a>` gilt als Anker, und mit einer echten `id` gleichen
  Namens danach misst der Vergleich die falsche Stelle (`cmp 0` möglich).
  Heute ohne Fundstelle (übernommen aus dem Re-Review zu Fixrunde 2); der
  Messende prüft diese Stellung am Diff.
- **Escapter Backtick vor einem Code-Span (fail-open).** Die Zählung kennt
  keinen escapten Backtick: ein `` \` `` zählt mit und verschiebt die Parität.
  Steht er vor einem Code-Span mit `<a id="x"></a>` und ist die Gesamtzahl
  gerade, gilt die `id` im Code-Span als Anker; mit einer echten `id` gleichen
  Namens danach misst der Vergleich die falsche Stelle (`cmp 0` möglich). Die
  Gegenrichtung endet mit Exit 2 (fail-closed). Heute ohne Fundstelle
  (übernommen aus dem Re-Review zu Fixrunde 3); der Messende prüft diese
  Stellung am Diff.
- **Eingerückter Absatz in einem Listenpunkt.** Eine `id` in einer Zeile mit
  vier Leerzeichen Einzug wird auch dann nicht gelesen, wenn die Zeile nach
  CommonMark ein Absatz eines Listenpunkts ist; der Verweis endet mit Exit 2
  (fail-closed), adressiert wird dann über ein Heading oder eine nicht
  eingerückte `id`.
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
- **Was nicht im Werkzeug steht:** die Datei-Schleife des `formnorm`-`cmp` am
  Form-Commit aus
  [`ADR-0161`](../../docs/plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
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
nach `id`-Zeile, Einzug 3 ist ein Fence, Einzug 4 keiner) · gestapelte `id`
(vor Heading und vor Absatz) · `id` in anderer Form (Attribut,
selbstschließend, neben gleichnamigem Heading-Slug, in Inline-Code) · `id` in
eingerücktem Code und in HTML-Kommentar (mehrzeilig, einzeilig, `<!--` in
Inline-Code) · Normalisierung nur am Segment (`tool-v6.14.0` bleibt roh, Punkt
im Tag ist kein Platzhalter) · Slug (ohne HTML-Tags, Dubletten `-1`/`-2`,
Ebene höchstens 6, schließende `#`-Folge) · mehrdeutige Stellung mit Exit 2
(Kommentar hinter einzelnem Backtick, neben Doppel-Backtick-Span, im Absatz
mit offenem Code-Span, `id` nach unsicherer Kommentar-Grenze, Einzug aus
Leerzeichen und Tab) und ihre Gegenprobe (Tabellenzeile mit Inline-Code
bleibt gelesen) · Kommentar in der `id`-Zeile vor einem Heading · Lokator mit
führenden Nullen (`L010-012`) · `id` in einem Code-Span mit Text davor
(Absatz, mit Kommentar, Tabellenzelle, nur im Code-Span) und die Gegenprobe
`id` nach geschlossenem Code-Span · Lokator strikt (`L7`, `L1-2-3`,
`L5-3`, `L0-1`) · Locale (`LC_ALL=C` beim Aufrufer) und Stub-`awk` ohne
Multibyte · Argumentzahl · roter Vergleich unter `set -euo pipefail` beim
Aufrufer (über `SHELLOPTS`) · `source` lässt die Shell-Optionen des Aufrufers
unverändert. Der Prüfling ist per `PROG=<Datei>` übersteuerbar
(Mutationsläufe an Kopien).
