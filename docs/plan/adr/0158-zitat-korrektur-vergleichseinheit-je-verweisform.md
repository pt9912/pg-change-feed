# ADR-0158: Zitat-Korrektur — Vergleichseinheit und Normalisierung des gemessenen Referenten je Verweisform

**Status:** Accepted — Supersedes [`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(nur Bedingung (b) in deren Entscheidung 1, „der Referent ist gemessen
derselbe“; die Bedingungen (a) und (c), die Liste der unberührbaren Felder,
die Entscheidungen 2 bis 5 und die von `ADR-0157` unberührten Teile von
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) bleiben
unverändert in Kraft und werden hier nicht wiederholt)

**Datum:** 2026-10-06

**Autor:** pt9912 (Architect-Rolle, Modul 8; Architect-Zug zu
`slice-zitat-korrektur-vergleichseinheit` Liefer-Punkt 1, anderer Kontext als
der Reviewer-Lauf der Fixrunde von `slice-harness-baseline-v6-14-1`, dessen
Finding F-1 diese ADR beantwortet)

**Bezug:** [`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(Entscheidung 1 (a) bis (c), Entscheidung 4) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (die Klasse) ·
[`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md) (Hausform ohne
Adresse) · [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md)
(gemessen/hergeleitet) · `AGENTS.md` §3.5, §3.12 · Review
`review-slice-harness-baseline-v6-14-1-fixrunde` F-1.

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0073` und `ADR-0157`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0157` Entscheidung 1 lässt eine Zitat-Korrektur in jedem Abschnitt einer
`Accepted`-ADR zu, wenn (a) nur die Form eines Verweises geändert wird — Pfad,
Versions-Segment, Linkziel, Anker, Zeilen-Lokator, Form einer gebrochenen
Referenz —, (b) der Referent gemessen derselbe ist („`cmp` des Zielinhalts an
alter und neuer Adresse“) und (c) kein Wort der Aussage sich ändert.

Bedingung (b) sagt nicht, **welche Einheit** verglichen wird. Die einzige
ausgeführte Form ist der `cmp` ganzer Dateien. Für zwei der zugelassenen Formen
misst er etwas anderes als den Referenten (übernommen aus dem Re-Review F-1,
hier nachgemessen, siehe Fitness Function):

- Ein **Anker-Wechsel** in derselben Datei (`#teilfrage-2-…` → `#teilfrage-3-…`)
  vergleicht die Datei mit sich selbst; der `cmp` ist immer 0, obwohl der
  Verweis jetzt auf einen anderen normativen Abschnitt zeigt.
- Ein **verschobener Zeilen-Lokator** bei gleichem Inhalt: Die Datei hat sich
  verändert (deshalb wurde der Lokator nachgezogen); der Datei-`cmp` ist 1,
  obwohl die zitierten Zeilen gleich sind.

Auch die **Normalisierung** ist für Entscheidung 1 offen: (b) verlangt nur, dass
der Beleg „roh oder normalisiert“ nennt; welche Normalisierung zulässig ist,
legt `ADR-0157` nur für MR-Pins fest (Entscheidung 4).

## Entscheidung

Bedingung (b) aus `ADR-0157` Entscheidung 1 lautet ab hier: **Der Referent ist
gemessen derselbe — verglichen wird die Einheit, die der Verweis adressiert,
nicht die Datei, in der sie steht.**

1. **Vergleichseinheit je Verweisform.** Die Einheit folgt dem genauesten
   Lokator des Verweises; jede Seite (alte und neue Adresse) wird nach ihrem
   eigenen Lokator gelesen.

   | Form aus (a) | Vergleichseinheit |
   |---|---|
   | Pfad / Linkziel ohne Anker und Lokator | die ganze Zieldatei |
   | Anker auf ein Heading | der **Abschnittskörper**: die Zeilen nach dem Ziel-Heading bis vor das nächste Heading gleicher oder höherer Ebene; Unterabschnitte gehören dazu, Headings in Code-Fences zählen nicht. Die Heading-Zeile selbst gehört **nicht** dazu — ein umbenanntes Heading ist der typische Anlass einer Anker-Korrektur |
   | Anker als HTML-`id` außerhalb eines Headings | die Zeile, die die `id` trägt |
   | Zeilen-Lokator (eine Zeile oder ein Bereich) | die zitierten Zeilen |
   | Versions-Segment | die Einheit, die der übrige Verweis adressiert (Datei, Abschnitt oder Zeilen), mit der Normalisierung aus Punkt 2 |
   | Form einer gebrochenen Referenz | die Einheit der Form, die der Verweis nach der Korrektur trägt; die alte Adresse wird am Stand nach Punkt 3 gelesen |

   Tragen alte und neue Adresse verschiedene Lokator-Arten (etwa Zeilen alt,
   Anker neu), misst jede Seite ihre eigene Einheit; die Korrektur besteht nur,
   wenn beide Einheiten byte-gleich sind.

2. **Normalisierung.** Verglichen wird **roh**. Die einzige zulässige
   Normalisierung ist die eines Versions-Segments: auf beiden Seiten wird jedes
   Pfadsegment `/v<MAJOR>.<MINOR>.<PATCH>/` durch `/<tag>/` ersetzt (dieselbe
   Regel wie `ADR-0157` Entscheidung 4, auf jedes solche Segment ausgedehnt,
   auch auf Upstream-URLs im Text der Einheit). Sie ist nur zulässig, wenn die
   Korrektur ein Versions-Segment bewegt und der rohe Vergleich fällt. Keine
   Normalisierung von Leerraum, Zeilenenden oder Groß-/Kleinschreibung.

3. **Stände.** Die alte Adresse wird am jüngsten Stand gelesen, an dem sie
   auflöst — in der Regel der Parent des Korrektur-Commits; ist die Zieldatei
   umgezogen oder gelöscht, der Parent des Commits, der sie bewegt hat
   (`git log -1 --format=%h --diff-filter=DR -- <alter Pfad>`, dann `~1`). Die
   neue Adresse wird am Korrektur-Commit gelesen. Ist eine Einheit leer oder
   nicht lesbar (Anker unbekannt, Datei fehlt), besteht die Korrektur nicht.

4. **Nicht messbarer Referent.** Liegt die alte Adresse außerhalb des Repos
   (etwa ein host-lokaler Pfad, `AGENTS.md` §3.11) oder löste sie an keinem
   Stand auf, gibt es keinen gemessenen Referenten. Dann ist die Korrektur in
   §Entscheidung, §Konsequenzen und §Verglichene Alternativen **keine**
   Zitat-Korrektur, sondern Folge-ADR; in den übrigen Abschnitten gilt
   `ADR-0073` unverändert (Urteil am Diff), und der Beleg nennt „nicht messbar“
   mit Grund.

5. **Beleg.** Der Beleg einer Zitat-Korrektur nach `ADR-0157` Entscheidung 1
   nennt je Verweis: Form, Einheit, beide Stände, roh oder normalisiert, und
   die gedruckte Zeile des Vergleichs. Er steht in der Commit-Message oder im
   Slice-Plan; die §Geschichte-Zeile der ADR bleibt bei der Form aus
   `ADR-0073` Entscheidung 3.

6. **Befehlsform — eine Verfahrensregel, kein Werkzeug.** Gemessen wird mit
   dieser Shell-Form (`bash`, `git`, `awk`, `sed`, `cmp`; Host-Werkzeuge nach
   `AGENTS.md` §3.1). `vergleich` druckt je Aufruf eine Zeile mit dem
   `cmp`-Exit und gibt ihn zurück; eine leere oder nicht lesbare Einheit endet
   mit 2:

   ```bash
   # einheit <stand> <pfad> [<ref>]  ref: leer = Datei | #slug = Abschnittskörper | L<a>-<b> = Zeilen
   einheit() {
     local stand="$1" pfad="$2" ref="${3:-}" out
     out=$(git show "$stand:$pfad" 2>/dev/null) || { echo "einheit: $stand:$pfad nicht lesbar" >&2; return 2; }
     case "$ref" in
       "") printf '%s\n' "$out" ;;
       L*) local a=${ref#L}; local b=${a#*-}; a=${a%-*}
           out=$(printf '%s\n' "$out" | sed -n "${a},${b}p") ;;
       \#*) out=$(printf '%s\n' "$out" | awk -v want="${ref#\#}" '
             /^(```|~~~)/ { fence = !fence }
             !fence && /^#{1,6} / {
               lvl = match($0, /[^#]/) - 1
               if (inside && lvl <= ilvl) exit
               if (!inside) {
                 s = tolower(substr($0, lvl + 2)); gsub(/[^[:alnum:] _-]/, "", s); gsub(/ /, "-", s)
                 n = seen[s]++; if (n) s = s "-" n
                 if (s == want) { inside = 1; ilvl = lvl; next }
               }
             }
             inside { print }') ;;
     esac
     [ -n "$out" ] || { echo "einheit: leere Einheit $stand:$pfad$ref" >&2; return 2; }
     printf '%s\n' "$out"
   }
   # vergleich <alt-stand> <alt-pfad> <alt-ref> <neu-stand> <neu-pfad> <neu-ref> [norm]
   vergleich() {
     local n=cat; [ "${7:-}" = norm ] && n="sed -E s#/v[0-9]+\.[0-9]+\.[0-9]+/#/<tag>/#g"
     local x y; x=$(einheit "$1" "$2" "$3") || return 2; y=$(einheit "$4" "$5" "$6") || return 2
     cmp <(printf '%s\n' "$x" | $n) <(printf '%s\n' "$y" | $n); local ec=$?
     echo "vergleich ${7:-roh}: $1:$2$3 <-> $4:$5$6 cmp $ec"; return $ec
   }
   ```

   Der HTML-`id`-Fall aus Punkt 1 ist in dieser Form nicht ausgeführt; dort
   wird die Zeile per Zeilen-Lokator (`L<n>-<n>`) an beiden Ständen gelesen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; Datei-`cmp` für jede Form | kein Eingriff | ein Anker-Wechsel auf einen anderen Abschnitt derselben Datei besteht immer (gemessen: `cmp` 0), ein verschobener Lokator fällt immer; die Bedingung misst die zugelassene Form nicht |
| B — Datei-`cmp` bleibt, Anker- und Lokator-Wechsel in den Aussage-Abschnitten werden verboten | einfach, keine Abschnitts-Logik | verbietet den häufigsten legitimen Fall (Heading im Ziel umbenannt, Abschnitt gleich); Lokator-Korrekturen gingen nur noch per Folge-ADR — der Weg, den `ADR-0073` Option C als falsches Werkzeug verworfen hat |
| **C — Einheit je Verweisform, roh, Normalisierung nur für Versions-Segmente, Befehlsform in der ADR (gewählt)** | misst, was der Verweis adressiert; fängt den Anker-Wechsel (gemessen: `cmp` 1) und lässt den verschobenen Abschnitt und Lokator bestehen (gemessen: `cmp` 0); kein neues Werkzeug, keine neue Datei | die Abschnittsgrenze ist eine nachgebildete Slug-Regel, kein GitHub-Renderer (Grenze unten); die Befehlsform lebt als Text in einer ADR und hat keinen Tabellentest |
| D — wie C, aber als Skript hinter `make` mit Tabellentest | Slug-Regel und Fälle maschinell gesichert | eigener Slice (Plan §4 Rückführung) für eine Messung, die nur bei Zitat-Korrekturen anfällt — meist einmal je Baseline-Bump; der Aufwand trägt erst, wenn die Befehlsform sich als zu schwach zeigt (Re-Evaluierungs-Trigger (a)) |
| E — „gemessen“ streichen, der Referent wird am Diff beurteilt | kein Messaufwand | gibt die einzige maschinelle Hälfte von `ADR-0157` Entscheidung 1 auf; der Kanal in die Aussage-Abschnitte stünde allein auf Urteil |

## Konsequenzen

- Positiv: Bedingung (b) misst jetzt den Referenten jeder Form aus (a). Ein
  Anker-Wechsel auf einen anderen Abschnitt fällt, ein Wort im Referenten
  fällt, ein Umzug, ein verschobener Abschnitt und ein nachgezogener Lokator
  bestehen.
- Positiv: Die bisherigen Zitat-Korrekturen nach `ADR-0157` (`ADR-0095`,
  `eadf3054`; MR-Pins `5d8855d9`) sind Versions-Segmente auf ganze Dateien und
  bleiben gültig: dort ist Einheit gleich Datei, und der rohe bzw. nach
  `ADR-0157` Entscheidung 4 normalisierte `cmp` ist der hier verlangte.
- Negativ mit Grenze: Die Abschnittsgrenze bildet die Slug-Regel von GitHub
  nach (Kleinschreibung, Satzzeichen entfernt, Leerzeichen zu `-`, Zähler bei
  Dubletten). Gemessen trägt sie für die Headings unten (Umlaut, Gedankenstrich,
  Doppelpunkt, doppelter Bindestrich); für andere Zeichen ist sie hergeleitet.
  Weicht sie ab, ist die Einheit leer und der Vergleich endet mit 2 — die
  Abweichung macht den Lauf rot, nicht grün. `tolower` auf Mehrbyte-Zeichen ist
  an GNU awk gemessen; mit einem anderen `awk` hergeleitet.
- Negativ mit Grenze: Die Heading-Zeile gehört nicht zur Einheit. Ändert eine
  Korrektur nur, worauf der Verweis zeigt, weil das Ziel-Heading inhaltlich
  umbenannt wurde, der Körper aber gleich blieb, besteht sie; ob der neue Name
  dieselbe Aussage trägt, bleibt Urteil am Diff (`ADR-0157` Entscheidung 1 (c)).
- Negativ: Punkt 4 schließt die Aussage-Abschnitte für nicht messbare Referenten.
  Eine Korrektur eines host-lokalen Pfads in §Entscheidung braucht dort eine
  Folge-ADR; bisher gab es keinen solchen Fall nach `ADR-0157`.
- Folgepflicht (Implementer, `slice-zitat-korrektur-vergleichseinheit`
  Liefer-Punkt 2): `AGENTS.md` §3.5 nennt die Vergleichseinheit und verweist
  auf diese ADR; der Index trägt sie bereits.

## Fitness Function (falls maschinell prüfbar)

Alle Proben **gemessen** an einem Klon des Repos im Scratchpad (Stand
`b7d97cca`, GNU Awk 5.2.1), Instanz die Befehlsform aus Punkt 6; die
Probe-Commits existieren nur im Klon. Die Abschnitts-Proben laufen an
[`ADR-0100`](0100-nats-dritter-vollinhalts-zustellweg.md) (Headings
`### Teilfrage 1` bis `6`).

| Tooling | Regel | Make-Target |
|---|---|---|
| — | Entscheidung 1 (a) und (c) aus `ADR-0157` bleiben Urteil am Diff; maschinell ist nur der Vergleich der Einheit | — |
| Befehlsform Punkt 6 — Anker | **Lücke von A:** Datei-`cmp` von `ADR-0100` gegen sich selbst `cmp 0`. **Anker-Wechsel** `#teilfrage-2--subjekt--und-nachrichtenschema` → `#teilfrage-3--zustellsemantik-core-nats-vs-jetstream` an derselben Datei: `cmp 1`, fällt; Teilfrage 2 → `#entscheidung` (Elternabschnitt): `cmp 1`. **Verschoben:** Probe-Commit stellt Teilfrage 6 vor 5 und fügt drei Zeilen in §Kontext ein; `#teilfrage-5-…` alt gegen neu: `cmp 0`, besteht; der Datei-`cmp` derselben Änderung `cmp 1`. **Heading umbenannt** (Teilfrage 3 ohne „: Core NATS vs. JetStream“), Slug alt gegen Slug neu: `cmp 0`. **Unbekannter Anker:** `einheit: leere Einheit …`, Exit 2 | kein Make-Target (Lauf des Korrigierenden, Prüfung im Review) |
| Befehlsform — Zeilen-Lokator | Nach den drei eingefügten Zeilen: `L139-145` alt gegen `L142-148` neu `cmp 0`, besteht; Gegenprobe `L139-145` gegen `L139-145` (Lokator nicht nachgezogen) `cmp 1` | — |
| Befehlsform — Pfad und gebrochene Referenz | Probe-Commit `git mv docs/user/version.md docs/user/versionen.md`: alt gegen neu `cmp 0`; gegen eine andere Datei `cmp 1`. Die alte Adresse am Kopf des Klons: `einheit: … nicht lesbar`, Exit 2; am Stand nach Punkt 3 (`git log -1 --diff-filter=DR`, dann `~1`) gegen die neue Adresse `cmp 0` | — |
| Befehlsform — Versions-Segment | `11a5bac5` `v6.14.0/templates/.d-check.yml` gegen `b7d97cca` `v6.14.1/…` roh `cmp 0`; `regelwerk/modul-05-planning-harness.md#ziel-form-slice` roh `cmp 0`; `regelwerk/grundlagen-begriffe.md` roh `cmp 1` (Quelle-Zeile trägt den Tag in der URL), normalisiert `cmp 0`; `regelwerk/README.md` normalisiert `cmp 1` (Stand-Zeile ändert sich über den Tag hinaus) | — |
| Befehlsform — Mutation | **eine** Stelle: im Abschnittskörper von Teilfrage 2 ein Wort (`Contra` → `Kontra` im Tabellenkopf), Probe-Commit, gleicher Anker alt gegen neu: `cmp 1`, rot. Dass jede andere Byte-Änderung im Körper ebenso fällt, ist *hergeleitet* aus der Byte-Gleichheit des `cmp`. Der HTML-`id`-Fall, Headings in Code-Fences und der Dubletten-Zähler sind *hergeleitet*, nicht gefahren | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

**(a)** Ein Beleg weicht von der Befehlsform ab, weil sie eine Form nicht trägt
(Slug-Regel trifft ein Heading nicht, HTML-`id`, Fence) — dann Option D als
eigener Slice. **(b)** `d-check` bietet eine Abschnittsauflösung über seinen
Anker-Resolver als Kommando — dann ersetzt sie die `awk`-Nachbildung.
**(c)** Eine Korrektur besteht den Vergleich und erweist sich trotzdem als
inhaltlich — dann `ADR-0157` Re-Evaluierungs-Trigger (a). Andernfalls
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Accepted — Architect-Zug zu Re-Review F-1 der Fixrunde von `slice-harness-baseline-v6-14-1` | `slice-zitat-korrektur-vergleichseinheit` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0158` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
