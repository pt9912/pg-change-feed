# ADR-0159: Zitat-Korrektur — Einheit einer HTML-`id`, Referent-Messung an MR-Pins und gehärtete Befehlsform

**Status:** Accepted — Supersedes [`ADR-0158`](0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
(nur: in Entscheidung 1 die Zeile „Anker als HTML-`id` außerhalb eines
Headings“; Entscheidung 2, Normalisierung; Entscheidung 6, Befehlsform samt
dem Satz zum HTML-`id`-Fall darunter; und unter §Konsequenzen die Aussagen zu
den MR-Pins `5d8855d9` und zum Anker-Wechsel. Die übrigen Zeilen von
Entscheidung 1, die Entscheidungen 3 bis 5 und alles, was `ADR-0158` von
[`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) und
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) in Kraft lässt,
bleiben unverändert und werden hier nicht wiederholt)

**Datum:** 2026-10-06

**Autor:** pt9912 (Architect-Rolle, Modul 8; Architect-Zug zum Review
`review-slice-zitat-korrektur-vergleichseinheit`, anderer Kontext als der
Reviewer-Lauf und als der Architect-Lauf von `ADR-0158`)

**Bezug:** [`ADR-0158`](0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
(Entscheidung 1, 2, 4, 6) ·
[`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(Entscheidung 1, 3, 4) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) ·
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) (gemessen/hergeleitet) ·
`AGENTS.md` §3.1, §3.5, §3.12 · Review
`review-slice-zitat-korrektur-vergleichseinheit` F-1 bis F-3, F-5 bis F-9.

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0073`, `ADR-0157` und
`ADR-0158`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0158` legt fest, welche Einheit eine Zitat-Korrektur je Verweisform
vergleicht, und gibt die Befehlsform dafür. Das Review
`review-slice-zitat-korrektur-vergleichseinheit` hat an ihr gemessen
(übernommen aus dem Report; die tragenden Fälle sind hier nachgemessen, siehe
Fitness Function):

- **HTML-`id` (F-1).** Für eine `id` außerhalb eines Headings ist die Einheit
  die Zeile mit der `id`. Von den acht `id`-Ankern im Regelwerk von `v6.14.1`
  stehen sechs allein auf einer Zeile vor einem Heading, einer vor einem
  Absatz, einer in einer Heading-Zeile (gemessen mit `grep -n '<a id='` in
  `regelwerk/`, ohne die zwei Nennungen in Inline-Code); dazu tragen vier
  Tabellenzeilen in `harness/conventions.md` je eine `id`. Die Zeile
  `<a id="…"></a>` vor einem Heading hat keinen Inhalt und bleibt an jedem
  Stand gleich. Der reale Verweis `MR-004` →
  `modul-13-quality-gates.md#guard-haertung` bestünde jeden Bump, gleich was im
  Abschnitt „Guard-Härtung“ steht.
- **MR-Pins (F-2).** Die Konsequenz zu `5d8855d9` nennt die vier MR-Pins
  „Versions-Segmente auf ganze Dateien“. Alle vier tragen einen Anker; nach
  Entscheidung 1 ist ihre Einheit ein Abschnitt. Der `cmp` aus `ADR-0157`
  Entscheidung 4 liest die MR-Datei, nicht den Referenten. Ein Anker
  (`MR-001` → `grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument`)
  löst an keinem Stand auf (F-9).
- **Befehlsform (F-3, F-5, F-6).** Die Normalisierung steht ungequotet
  (`$n`); unter `shopt -s nullglob` oder `failglob` scheitert `sed`, beide
  Seiten sind leer, und zwei verschiedene Einheiten ergeben `cmp 0`. Unter
  `set -e` fehlt beim roten Vergleich die gedruckte Zeile. `$(…)` schneidet
  abschließende Leerzeilen ab, der Vergleich ist also nicht roh.
- **Normalisierung (F-7).** Sie ersetzt jedes `/v<X.Y.Z>/`, auch einen fremden
  Werkzeug-Pin.
- **Gleicher Körper (F-8).** Ein Anker-Wechsel zwischen zwei Abschnitten mit
  gleichem Körper besteht.

## Entscheidung

1. **Einheit einer HTML-`id`: der Abschnitt, den sie markiert.** Die Zeile
   „Anker als HTML-`id` außerhalb eines Headings | die Zeile, die die `id`
   trägt“ aus `ADR-0158` Entscheidung 1 wird ersetzt durch die Einheit nach
   der Stellung der `id`:

   | Stellung von `<a id="…">` | Vergleichseinheit |
   |---|---|
   | in einer Heading-Zeile | der Abschnittskörper dieses Headings (wie ein Anker auf ein Heading) |
   | allein auf ihrer Zeile, die nächste nicht leere Zeile ist ein Heading | der Abschnittskörper dieses Headings |
   | allein auf ihrer Zeile, die nächste nicht leere Zeile ist kein Heading | der Block ab dieser nächsten Zeile bis vor das nächste Heading beliebiger Ebene oder bis zum Dateiende |
   | in einer Tabellenzeile (die Zeile beginnt mit `\|`) | die Tabellenzeile |
   | in einer anderen Zeile mit Text | der Block ab dieser Zeile bis vor das nächste Heading beliebiger Ebene oder bis zum Dateiende |

   Aufgelöst wird die **erste** Fundstelle in Dateireihenfolge, Heading-Slug
   oder `id`. Eine `id` in einem Code-Fence und eine `id`, vor der unmittelbar
   ein Backtick steht (Inline-Code), zählen nicht. Gelesen wird nur
   `<a id="…">`; `name=` und eine `id` an einem anderen Element ergeben eine
   leere Einheit (Exit 2; heute 0 Fundstellen beider Formen in getrackten
   `.md`-Dateien, gemessen mit `git grep` am Stand `625ddbef`). Der
   Heading-Slug lässt HTML-Tags im Heading-Text weg.

2. **MR-Pins: Datei-`cmp` und Referent-Messung, beides.** Für eine
   Pin-Umstellung an MR-Einträgen (`ADR-0157` Entscheidung 3, 4) sind zwei
   Messungen nötig, und sie belegen Verschiedenes:
   - Der `cmp` der **MR-Datei** nach `ADR-0157` Entscheidung 4 belegt die
     Bedingungen (a) und (c) aus `ADR-0157` Entscheidung 1 maschinell: außer
     dem Tag hat sich in der Datei nichts bewegt. Er sagt nichts über den
     Referenten.
   - Die **Referent-Messung** je Verweis, dessen Versions-Segment der
     Pin-Commit bewegt, belegt Bedingung (b): `vergleich` aus Entscheidung 4
     mit der Einheit des Verweises (Datei, Abschnitt oder `id`-Einheit nach
     Punkt 1), alte Adresse am Parent des Pin-Commits, neue am Pin-Commit,
     roh, und nur wenn roh fällt, mit dem Tag-Paar.

   Die bewegten Verweise liest der Messende aus `git diff "$P~1" "$P" -- <MR-Datei>`.
   Belegt wird im Slice-Plan des Bumps (`ADR-0158` Entscheidung 5); der
   Verifier fährt beide Messungen nach.

   **Nicht messbarer Referent an einem MR-Eintrag.** `ADR-0158` Entscheidung 4
   unterscheidet die Aussage-Abschnitte einer ADR von den übrigen. Ein
   MR-Eintrag hat diese Abschnitte nicht; für ihn gilt der Zweig der übrigen
   Abschnitte: Löst die alte Adresse an keinem Stand auf, ist die Korrektur
   Urteil am Diff nach `ADR-0073`, und der Beleg nennt „nicht messbar“ mit
   Grund. Bedingung (c) bleibt gemessen durch den Datei-`cmp`.

   **Die Konsequenz zu `5d8855d9` aus `ADR-0158` lautet berichtigt:** Die
   Zitat-Korrektur `ADR-0095` (`eadf3054`) ist ein Versions-Segment auf eine
   ganze Datei und bleibt gültig (roh `cmp 0`, gemessen in `ADR-0158`). Die
   MR-Pins `5d8855d9` bewegen vier Verweise mit Anker. Drei davon bestehen die
   Referent-Messung roh (`MR-002`, `MR-003`, `MR-004`; siehe Fitness
   Function). Der vierte, `MR-001` →
   `grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument`,
   ist nicht messbar: Der Anker löst weder an `5d8855d9~1` noch an `5d8855d9`
   auf, das Heading „Spec-Straten: mehr als ein Spec-Dokument“ steht in
   `grundlagen-referenz-richtung.md`. Die Pin-Umstellung an `MR-001` ist damit
   Urteil am Diff (Datei-`cmp` normalisiert Exit 0, nur der Tag bewegt). Der
   gebrochene Anker ist ein **Befund** (Review F-9); er bestand schon am Parent
   des Bumps. Ob und wie er behoben wird, entscheidet die Closure von
   `slice-zitat-korrektur-vergleichseinheit`. Eine Behebung ist eine Korrektur
   der Form einer gebrochenen Referenz, deren alte Adresse an keinem Stand
   auflöst, also ebenfalls „nicht messbar“ mit Grund.

3. **Normalisierung: nur der bewegte Tag.** `ADR-0158` Entscheidung 2 Satz 2
   („jedes Pfadsegment `/v<MAJOR>.<MINOR>.<PATCH>/` …, auch auf Upstream-URLs“)
   wird ersetzt: Die Normalisierung nimmt das **Tag-Paar** des Bumps
   (`<alt-tag>:<neu-tag>`). Auf der alten Seite wird nur das Segment
   `/<alt-tag>/` durch `/<tag>/` ersetzt, auf der neuen nur `/<neu-tag>/`. Das
   trifft den Baseline-Pfad (`.harness/baseline/<tag>/`) und die Kurs-URLs, die
   denselben Tag tragen (`ai-harness-course/blob/<tag>/`); ein fremder Pin mit
   einer anderen Version bleibt roh. Die übrigen Sätze von Entscheidung 2
   gelten weiter (roh zuerst; Normalisierung nur, wenn die Korrektur ein
   Versions-Segment bewegt und roh fällt; keine Normalisierung von Leerraum,
   Zeilenenden oder Groß-/Kleinschreibung). Der MR-Datei-`cmp` aus `ADR-0157`
   Entscheidung 4 behält seine eigene Normalisierung (`.harness/baseline/v…/`).

4. **Befehlsform — ersetzt `ADR-0158` Entscheidung 6.** Gemessen wird mit
   dieser Shell-Form (`bash`, `git`, `awk`, `sed`, `cmp`; Host-Werkzeuge nach
   `AGENTS.md` §3.1). Sie bleibt eine Verfahrensregel in dieser ADR, kein
   Werkzeug. Zusagen der Form:
   - Jede Expansion steht in Anführungszeichen; Shell-Optionen zur
     Pfadnamen-Expansion (`nullglob`, `failglob`) ändern das Ergebnis nicht.
   - Eine leere oder nicht lesbare Einheit, ein ungültiger Lokator und ein
     ungültiges Tag-Paar enden auf **jeder** Seite mit Exit 2 und einer
     gedruckten Zeile, nie mit `cmp 0`.
   - `vergleich` druckt seine Zeile vor dem `return`; auch unter `set -e` steht
     sie beim roten Vergleich da, danach bricht das Skript ab.
   - Roh heißt byte-gleich, abschließende Leerzeilen eingeschlossen; die Form
     hängt dazu an jede Kommando-Substitution ein Wächter-Zeichen und schneidet
     es danach ab.

   ```bash
   # einheit <stand> <pfad> [<ref>]  ref: leer = Datei | #anker = Heading-Abschnitt oder HTML-id | L<a>-<b> = Zeilen
   einheit() {
     local stand="$1" pfad="$2" ref="${3:-}" out a b
     out=$(git show "$stand:$pfad" 2>/dev/null && printf .) || { echo "einheit: $stand:$pfad nicht lesbar" >&2; return 2; }
     out="${out%.}"
     case "$ref" in
       "") ;;
       L*) a="${ref#L}"; b="${a#*-}"; a="${a%%-*}"
           [[ "$a" =~ ^[0-9]+$ && "$b" =~ ^[0-9]+$ ]] || { echo "einheit: Lokator $ref" >&2; return 2; }
           out=$(printf '%s' "$out" | sed -n "${a},${b}p" && printf .); out="${out%.}" ;;
       \#*) out=$(printf '%s' "$out" | awk -v want="${ref#\#}" '
             function slug(t, s) { s = tolower(t); gsub(/<[^>]*>/, "", s); gsub(/[^[:alnum:] _-]/, "", s); gsub(/ /, "-", s); return s }
             function hasid(l) { return index(l, "<a id=\"" want "\"") && !index(l, "`<a id=\"" want "\"") }
             done { next }
             /^(```|~~~)/ { fence = !fence }
             !fence && /^#{1,6} / {
               lvl = match($0, /[^#]/) - 1
               if ((mode == "abs" && lvl <= ilvl) || mode == "blk") { done = 1; next }
               if (mode == "vor") { mode = "abs"; ilvl = lvl; next }
               if (mode == "") {
                 s = slug(substr($0, lvl + 2)); n = seen[s]++; if (n) s = s "-" n
                 if (s == want || hasid($0)) { mode = "abs"; ilvl = lvl; next }
               }
             }
             mode == "vor" && /^[[:space:]]*$/ { next }
             mode == "vor" { mode = "blk" }
             mode == "abs" || mode == "blk" { print; next }
             mode == "" && !fence && hasid($0) {
               if ($0 ~ /^[[:space:]]*[|]/) { print; done = 1; next }
               r = $0; gsub(/<a id="[^"]*"><\/a>/, "", r)
               if (r ~ /^[[:space:]]*$/) { mode = "vor"; next }
               mode = "blk"; print
             }' && printf .); out="${out%.}" ;;
       *) echo "einheit: Lokator $ref" >&2; return 2 ;;
     esac
     [ -n "$out" ] || { echo "einheit: leere Einheit $stand:$pfad$ref" >&2; return 2; }
     printf '%s' "$out"
   }
   # tagnorm <tag>  ersetzt nur das Segment /<tag>/
   tagnorm() { sed -E "s#/${1//./\\.}/#/<tag>/#g"; }
   # vergleich <alt-stand> <alt-pfad> <alt-ref> <neu-stand> <neu-pfad> <neu-ref> [<alt-tag>:<neu-tag>]
   vergleich() {
     local x y ec art=roh
     x=$(einheit "$1" "$2" "$3" && printf .) || { echo "vergleich: $1:$2$3 keine Einheit, Exit 2"; return 2; }
     y=$(einheit "$4" "$5" "$6" && printf .) || { echo "vergleich: $4:$5$6 keine Einheit, Exit 2"; return 2; }
     x="${x%.}"; y="${y%.}"
     if [ -n "${7:-}" ]; then
       [[ "$7" =~ ^v[0-9]+\.[0-9]+\.[0-9]+:v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo "vergleich: Tag-Paar $7, Exit 2"; return 2; }
       art="norm $7"
       x=$(printf '%s' "$x" | tagnorm "${7%%:*}" && printf .); x="${x%.}"
       y=$(printf '%s' "$y" | tagnorm "${7#*:}" && printf .); y="${y%.}"
     fi
     if cmp <(printf '%s' "$x") <(printf '%s' "$y"); then ec=0; else ec=$?; fi
     echo "vergleich $art: $1:$2$3 <-> $4:$5$6 cmp $ec"
     return "$ec"
   }
   ```

5. **Benannte Grenzen der Einheit.**
   - **Gleicher Körper (F-8).** Die Heading-Zeile gehört nicht zur Einheit
     (`ADR-0158` Entscheidung 1). Ein Anker-Wechsel zwischen zwei Abschnitten
     mit byte-gleichem Körper besteht deshalb den Vergleich. Dass der neue
     Anker auf dieselbe Aussage zeigt, bleibt in diesem Fall Urteil am Diff
     (`ADR-0157` Entscheidung 1 (c)). Heute gibt es keinen solchen Abschnitt
     (übernommen aus dem Review, Probe P3: 0 im Baum ohne Records).
   - **Schluss-Zeilenumbruch.** Im Abschnitts- und Block-Modus schreibt `awk`
     jede Zeile mit Zeilenumbruch, auch die letzte Zeile einer Datei ohne
     Schluss-Umbruch. Ein Unterschied allein im Schluss-Umbruch der Datei
     fällt dort nicht auf; im Datei- und im Zeilen-Modus fällt er (gemessen,
     siehe Fitness Function).
   - **Slug mit HTML im Heading.** Steht vor einem Tag am Zeilenende ein
     Leerzeichen, endet der nachgebildete Slug auf `-` und trifft den Slug
     des Renderers nicht; die Einheit ist leer, der Vergleich endet mit 2
     (gemessen an `#die-22-matrix-böckeler`). Adressiert wird dann über die
     `id`.
   - **NUL-Bytes.** Eine Shell-Variable trägt kein NUL-Byte; die Form ist für
     Textdateien gemacht. Heute trägt keine getrackte `.md`-Datei eines
     (gemessen mit `grep -lP '\x00'` über `git ls-files '*.md'`, 0 Treffer).
   - **Tag gleich fremder Version.** Trägt ein fremder Pin zufällig genau den
     alten oder neuen Baseline-Tag, normalisiert die Form auch ihn. Heute steht
     `/v6.14.1/` außerhalb von Baseline-Pfad und Kurs-URL an einer Stelle, in
     einem Record unter `done/` (gemessen mit `git grep '/v6\.14\.1/'` am
     Stand `625ddbef`).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0158` bleibt | kein Eingriff | der reale Verweis `MR-004` → `#guard-haertung` besteht jeden Bump ohne Messung des Abschnitts; die Konsequenz zu `5d8855d9` bleibt falsch; die Befehlsform ist unter `nullglob`/`failglob` fail-open (Review F-3) |
| B — HTML-`id` als eigene Zeile belassen und für die Aussage-Abschnitte verbieten | einfache Einheit | verbietet den realen Fall `MR-004` und die acht `id`-Anker des Regelwerks als Ziel einer Zitat-Korrektur; verschiebt die Arbeit auf Folge-ADRs, ohne dass sich etwas am Referenten bewegt hat |
| C — Normalisierung nur auf `.harness/baseline/v…/` | genau der Pfad, den `ADR-0157` Entscheidung 4 normalisiert | alle 26 Regelwerk-Dateien von `v6.14.1` tragen den Tag in der Kurs-URL ihrer Quelle-Zeile (`ai-harness-course/blob/v…/`), in 8 davon steht sie in Zeile 2 und damit im Körper des ersten Abschnitts (gemessen mit `grep -l` und `head -1`/`sed -n 2p` am Stand `625ddbef`); der Vergleich der ganzen Datei und dieser Abschnitte fiele an jedem Bump, ohne dass sich der Referent bewegt hat (gemessen an `grundlagen-begriffe.md`, siehe Fitness Function) |
| **D — Einheit je Stellung der `id`, beide Messungen an MR-Pins, Normalisierung nur des bewegten Tags, gehärtete Befehlsform in der ADR (gewählt)** | misst an `#guard-haertung` den Abschnitt; fail-closed unter jeder geprüften Shell-Option; Leerzeilen roh; fremde Pins roh; kein neues Werkzeug | die Befehlsform ist länger und hat weiter keinen Tabellentest; die `id`-Regel ist eine Nachbildung des Renderers, keine Auflösung durch ihn |
| E — Befehlsform als Skript hinter `make` mit Tabellentest | maschinell gesicherte Fälle, kein Herausziehen aus der ADR | eigener Slice für eine Messung, die je Baseline-Bump einmal anfällt; in diesem Zug nicht gewählt, bleibt Re-Evaluierungs-Trigger (a) |

## Konsequenzen

- Positiv: Ein Verweis auf eine HTML-`id` vor einem Heading misst den
  Abschnitt; eine Änderung im Abschnitt „Guard-Härtung“ macht den Vergleich
  rot (gemessen). Eine `id` in einer Tabellenzeile misst die Zeile.
- Positiv: An MR-Pins ist klar, was welcher `cmp` belegt; der Verifier fährt
  beide.
- Positiv: Die Befehlsform endet bei leerer Einheit auf jeder Seite mit 2,
  auch unter `nullglob` und `failglob`, und vergleicht abschließende
  Leerzeilen roh (gemessen).
- Negativ mit Grenze: Ein Anker-Wechsel auf einen anderen Abschnitt fällt,
  **wenn** dessen Körper sich in mindestens einem Byte unterscheidet; bei
  gleichem Körper bleibt es Urteil (Entscheidung 5).
- Negativ mit Grenze: Die Block-Einheit einer `id` vor einem Absatz reicht bis
  zum nächsten Heading; eine Änderung in einem späteren Absatz desselben
  Blocks macht den Vergleich rot, auch wenn der markierte Absatz gleich blieb.
  Das ist fail-closed gewählt: der Fehler ist eine Folge-ADR zu viel, nie eine
  stille Inhaltsänderung.
- Negativ: `MR-001` trägt einen Anker, der an keinem Stand auflöst (Befund,
  Entscheidung 2). Die Behebung ist offen bis zur Closure des Slice.
- Folgepflicht (Implementer, `slice-zitat-korrektur-vergleichseinheit`):
  `AGENTS.md` §3.5 nennt die `id`-Einheit, die Normalisierung des bewegten
  Tags und den Geltungsbereich des Belegs, mit Verweis auf diese ADR;
  `.claude/agents/verifier.md` und `harness/targets/pin-stale.md` nennen am
  Pin-Commit neben dem Datei-`cmp` die Referent-Messung je bewegtem Verweis.
  Der Index trägt diese ADR.

## Fitness Function (falls maschinell prüfbar)

Alle Proben **gemessen** an einem Klon des Repos im Scratchpad (Stand
`625ddbef`, GNU Awk 5.2.1, GNU bash). Instanz ist die Befehlsform aus
Entscheidung 4, per `awk` wörtlich aus dem `bash`-Block dieser Datei gezogen
(56 Zeilen, Einrückung von drei Leerzeichen entfernt) und per `source`
geladen. Probe-Commits existieren nur im Klon: `e613170b` (nur die
Heading-Zeile „Guard-Härtung“ geändert), `d02575d2` (ein Wort im Körper von
„Guard-Härtung“), `76282e13` (ein Wort in der Tabellenzeile `MR-004` von
`harness/conventions.md`), `a080f7df` (ein Wort in der Zeile `MR-003`),
`4593941c` (ein Wort im Absatz hinter
`<a id="jedes-artefakt-hat-einen-konsumenten">`), `dd0d32fd` (ein Wort im
Abschnitt mit `<a id="2x2-matrix">` im Heading), `70861bb8`, `a1521a75` und
`f7b23eeb` (Probe-Datei: Baseline-Pfad und Kurs-URL von `v6.14.0` auf
`v6.14.1`, danach ein fremder Pin `v0.79.0` → `v0.80.0`; `x\n` gegen
`x\n\n\n\n`), `d95da84d` (drei Zeilen in `ADR-0100` §Kontext), `9cf54985`
(`git mv docs/user/version.md docs/user/versionen.md`), `904f0ae7` (nur der
fehlende Schluss-Umbruch von `ADR-0001` ergänzt). Der Probenlauf (37 Fälle mit
Soll-Exit) ist dreimal gefahren: ohne Option, mit `shopt -s nullglob`, mit
`shopt -s failglob`; je Lauf alle 37 mit dem Soll-Exit, die Ausgaben der drei
Läufe bis auf die `shopt`-Kopfzeile byte-gleich (`diff`). Gedruckte Zeilen
gekürzt auf die Kernaussage.

| Tooling | Regel | Make-Target |
|---|---|---|
| Befehlsform — HTML-`id` (Entscheidung 1) | `<a id="guard-haertung">` vor dem Heading: nur die Heading-Zeile geändert (`e613170b`) `cmp 0`; ein Wort im Körper (`d02575d2`) `cmp 1`, und an demselben Probe-Commit die Zeile der `id` (`L266-266`, Einheit nach `ADR-0158`) `cmp 0` — die Lücke aus Review F-1. Die `id`-Einheit ist byte-gleich der Einheit des Heading-Slugs `#guard-härtung-wächter-reifen-in-wellen-modul-13` (`cmp` Exit 0). Tabellenzeile `#mr-004`: Wort in der eigenen Zeile `cmp 1`, Wort in der Nachbarzeile `MR-003` `cmp 0`. `id` vor einem Absatz (`#jedes-artefakt-hat-einen-konsumenten`): Wort im Absatz `cmp 1`. `id` in der Heading-Zeile (`#2x2-matrix`): Wort im Abschnitt `cmp 1`, Änderung in einer anderen Datei `cmp 0`. `id` nur in Inline-Code (`#mr-<NNN>` in `grundlagen-harness-dateien.md`): `keine Einheit, Exit 2`. *Hergeleitet:* die Stellung „`id` in einer Prosazeile mit Text“ (heute ohne Fundstelle) und die Regel der ersten Fundstelle bei gleichem Slug und gleicher `id` | kein Make-Target (Lauf des Korrigierenden, Prüfung im Review und durch den Verifier) |
| Befehlsform — MR-Pins `5d8855d9` (Entscheidung 2) | je Verweis `vergleich 5d8855d9~1 …/v6.14.0/regelwerk/<datei> '<anker>' 5d8855d9 …/v6.14.1/regelwerk/<datei> '<anker>'`: `MR-002` `#vergabe-woher-die-nächste-kennung-kommt` roh `cmp 0`; `MR-003` `#grenzen--ehrlich-benannt` roh `cmp 0`; `MR-004` `#guard-haertung` (Abschnitt nach Entscheidung 1) roh `cmp 0`; `MR-001` `#spec-straten-mehr-als-ein-spec-dokument` an `5d8855d9~1` `keine Einheit, Exit 2`, an `5d8855d9` `einheit: leere Einheit`, Exit 2 — nicht messbar; derselbe Anker in `grundlagen-referenz-richtung.md` löst auf (`cmp 0` gegen sich selbst). Datei-`cmp` aus `ADR-0157` Entscheidung 4 an `5d8855d9`: alle vier MR-Dateien je Exit 0 | — |
| Befehlsform — fail-closed (Entscheidung 4) | verschiedene Einheiten (`probe/p.md#leer` gegen `ADR-0100`) `cmp 1` in allen drei Läufen, auch unter `nullglob`/`failglob` (Review F-3: dort `cmp 0`). Fehlende Datei, unbekannter Anker links, unbekannter Anker rechts, Lokator `X1-2`, Tag-Paar `v6.14.0-v6.14.1` und `v6.14.0:v6.14.1x`: je eine gedruckte Zeile `vergleich: … Exit 2`, Exit 2 | — |
| Befehlsform — `set -euo pipefail` (Entscheidung 4) | roter Vergleich in einem Skript mit `set -euo pipefail`: die Zeile `vergleich roh: … cmp 1` steht da, danach Exit 1, die Folgezeile des Skripts wird nicht erreicht (Review F-5: Zeile fehlte). Fehlende Datei: `vergleich: … keine Einheit, Exit 2`, Exit 2. Grüner Vergleich am frühen Abschnitt `#1-einleitung` von `docs/user/benutzerhandbuch.md` (232049 Byte) dreimal in Folge unter `pipefail`: Exit 0 | — |
| Befehlsform — roh (Entscheidung 4, Review F-6) | `x\n` gegen `x\n\n\n\n` als ganze Datei `cmp 1`; eine Leerzeile als Zeilen-Lokator (`L2-2` gegen `L3-3`) ist eine Einheit, `cmp 0`, nicht „leer“. Grenze aus Entscheidung 5: nur der ergänzte Schluss-Umbruch von `ADR-0001` (`904f0ae7`) — Abschnitt `#geschichte` `cmp 0`, ganze Datei `cmp 1` | — |
| Befehlsform — Normalisierung (Entscheidung 3, Review F-7) | Probe-Datei, Abschnitt `#pins` mit Baseline-Pfad und Kurs-URL `v6.14.0` → `v6.14.1`: roh `cmp 1`, mit `v6.14.0:v6.14.1` `cmp 0`; zusätzlich fremder Pin `v0.79.0` → `v0.80.0` mit demselben Tag-Paar `cmp 1` (Review F-7: `cmp 0`). `grundlagen-begriffe.md` `11a5bac5` gegen `625ddbef`: ganze Datei roh `cmp 1`, mit Tag-Paar `cmp 0`; erster Abschnitt `#kernbegriffe-und-trennschärfen` mit Tag-Paar `cmp 0`. Alternative C (`sed` aus `ADR-0157` Entscheidung 4, nur `.harness/baseline/v…/`): ganze Datei `cmp 1`, erster Abschnitt `cmp 1` | — |
| Befehlsform — Regression `ADR-0158` §Fitness Function | an `ADR-0100`: Anker-Wechsel Teilfrage 2 → 3 `cmp 1`; gleicher Anker `cmp 0`; → `#entscheidung` `cmp 1`; nach drei eingefügten Zeilen (`d95da84d`) `L139-145` gegen `L142-148` `cmp 0`, nicht nachgezogen `cmp 1`, `#teilfrage-5--verbindungs-wiederverwendung-und-aktivierung` alt gegen neu `cmp 0`. `git mv` (`9cf54985`): alte Adresse am Commit `keine Einheit, Exit 2`, an `9cf54985~1` gegen die neue `cmp 0`. Versions-Segment `11a5bac5` gegen `625ddbef`: `templates/.d-check.yml` roh `cmp 0`; `modul-05-planning-harness.md#ziel-form-slice` roh `cmp 0`; `regelwerk/README.md` mit Tag-Paar `cmp 1` | — |
| Mutation | **Stellen:** je `id`-Stellung aus Entscheidung 1 außer der Prosazeile **eine** Stelle (ein Wort, Probe-Commits oben), Instanz die Befehlsform als Shell-Lauf, gesehen `cmp 1`, rot; die Gegenproben an der Nachbarzeile bzw. an einer anderen Datei grün. Dass jede andere Byte-Änderung in der Einheit ebenso fällt, ist *hergeleitet* aus der Byte-Gleichheit des `cmp`. Die Befehlsform hat keinen Go- oder Tabellentest; dass der Verifier sie nachfährt, ist eine Erwartung | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

**(a)** Die Befehlsform braucht eine weitere Berichtigung, oder ein Beleg
weicht von ihr ab, weil sie eine Form nicht trägt (Slug-Regel, `id`-Stellung,
Fence) — dann Option E als eigener Slice, statt einer weiteren Folge-ADR zur
Befehlsform. **(b)** `d-check` bietet eine Abschnittsauflösung über seinen
Anker-Resolver als Kommando — dann ersetzt sie die `awk`-Nachbildung.
**(c)** Eine Korrektur besteht den Vergleich und erweist sich trotzdem als
inhaltlich — dann `ADR-0157` Re-Evaluierungs-Trigger (a). Andernfalls
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Accepted — Architect-Zug zum Review `review-slice-zitat-korrektur-vergleichseinheit` | `slice-zitat-korrektur-vergleichseinheit` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0159` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
