# Review-Report: slice-092 — **Delta** (Fixrunden `ad95444` und `a81a6ea`) · 2026-09-16

**Review-Art:** Delta-Review **zweier** Fix-Commits — geprüft gegen **Plan und
Entscheidungen**, gegen die Findings der ersten Runde (`review-slice-092.md`
F-1…F-3) und gegen die Adressen des Verifikationsberichts (`verify-slice-092.md`
V-1…V-3). Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten. DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Register-Zähler und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand.

**Gegenstand** (Stand `HEAD` = `a81a6ea`, Arbeitsbaum vor und nach jedem Lauf
sauber):

| Commit | Umfang (gemessen) |
|---|---|
| `ad95444` | 2 Dateien, **+13/−7**: `slice-092-coverage-cluster-d1.md` (§1, §2/LP1, §5, §6, §8), `harness/sensors/coverage-gate.md` (Deixis) |
| `a81a6ea` | 1 Datei, **+4/−1**: `harness/sensors/coverage-gate.md` (der „abgeleitet"-Marker) |
| Vorgang `3de9547..a81a6ea` | **4 Pfade**: 2 Träger, 2 Berichte — **kein** Produktcode, `THRESHOLD` unberührt |

**Skill:** `.harness/skills/reviewer.md` @ `HEAD` · **Datum:** 2026-09-16.

---

## Findings

### D-1 — Die neue Überschrift nennt **drei** Zahlen; die Liste darunter trägt **vier** — und der eigene Commit-Text sagt „vier"

- `kategorie`: **MEDIUM** · `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `slice-092-coverage-cluster-d1.md:51` gegen `:52-55`
- `befund`: Die von `ad95444` eingeführte Passage beginnt „**Drei Zahlen, drei
  Dinge — sie sind nicht austauschbar:**" und führt darunter **vier** Größen mit
  je eigener Aussage: **13** (Glob), **10** (trugen ungedeckte Statements),
  **12** (haben Tests bekommen), **11** (Deckung bewegt). Die Überschrift zählt
  eins zu wenig; der **Commit-Text desselben Zugs** sagt „das Dokument trug
  **vier** Zahlen (10, 11, 12, 13)". Die Stelle, deren einziger Zweck die
  Unterscheidung von Zählungen ist, trägt damit selbst eine falsche Zählung —
  und sie ist die **dritte** Runde an dieser Stelle.
- *Abgrenzung:* Gelesen als „drei im ersten Semikolon-Grupp (13 · 10 · 12) und
  die 11 danach" wäre die Überschrift rettbar — die Lesart trägt aber nicht, weil
  die Überschrift die **Verschiedenheit** zusichert („nicht austauschbar") und
  genau die 11 eine der vier ist.
- `verifizierbar`: ja — Überschrift gegen die eigene Liste halten; `git show
  ad95444 -- <Plan>`

### D-2 — Der Herkunfts-Marker ist invertiert, und die „Deutung" ist gegen die Messung: die zwei Bänder liegen auf **demselben** Code-Stand

- `kategorie`: **MEDIUM** · `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `quelle`: `AGENTS.md` §3.12 Instanz A („ein **abgeleiteter** Wert (Summe,
  Produkt, **Differenz**, Prozent) wird als **abgeleitet** gekennzeichnet — nie
  als gemessen ausgegeben") · *Bewegliche Zahlen* („die gedeckte Zahl … nennt
  ihn [den Lauf] und nie ‚der Ist-Stand'")
- `pfad`: `harness/sensors/coverage-gate.md:68-73` (neu, `a81a6ea`), benutzt in
  `:151-153`
- `befund`: Die eingefügte Klausel lautet „… **abgeleitet** aus der Differenz der
  gedeckten Zahlen (1468/1471 gegen 1369/1371), ~99–100 Statements auseinander;
  **gemessen ist die Differenz**, die Aussage über die Code-Stände ist ihre
  Deutung". Drei Beobachtungen:
  1. **Die zwei Etiketten sind vertauscht.** Gemessen sind die vier *gedeckten
     Zahlen* (je mit Lauf); die **Differenz** ~99–100 ist nach §3.12 Instanz A
     ausdrücklich **abgeleitet** — der Satz nennt sie „gemessen".
  2. **Die „Deutung" ist contra-gemessen.** Eigene Profile: der **Nenner ist
     1903 an beiden Bänder-Ständen** (`3d5b60d`: 1903/1369; `2b7c6ec`:
     1903/1469), und zwischen den beiden Bänder-Ständen hat sich im
     Messgegenstand **kein** Produktcode bewegt (`git diff --name-only
     7835b6e..f90c3f4 -- internal cmd` ohne `_test.go` und ohne das
     ausgenommene `replication/receive` → **leer**). Nach der Definition
     **derselben Datei** (§Zählbasis: „Der **Nenner** … hängt am **Code-Stand**")
     sind die zwei Bänder damit Stände **desselben** Code-Stands; die ~99–100
     liegen **vollständig in den gedeckten** Zahlen, also am **Test**-Baum. Die
     Klausel „die beiden Bänder gehören zu zwei verschiedenen Code-Ständen" ist
     unter der Lesart, die der Satz mit „Produktionsstand" selbst vorgibt,
     **falsch**.
  3. **Die Spanne nennt eine Paarung, nicht die Spanne.** „~99–100" ist die
     Gleich-End-Paarung; die Kreuz-Paarungen ergeben **97–102**.
  Was **nicht** neu ist: die Klausel „zwei verschiedene Code-Stände" stammt aus
  `ef493ed` (`slice-091`) und war schon Gegenstand von F-2/V-1. Neu ist, dass
  der Fix sie **durch eine Ableitung stützt** und mit einer falsch gesetzten
  Herkunft versieht — aus einem unbelegten Satz wird ein belegter, der nicht
  trägt.
- `verifizierbar`: ja — eigene netzlose Profile der Stände `3d5b60d`, `2b7c6ec`,
  `a81a6ea`; `git diff --name-only 7835b6e..f90c3f4 -- internal cmd`

### D-3 — Der Hauptsatz ist an der neu eingefügten Klausel gerissen; die Beleg-Klammer hängt jetzt am Satz der *Deutung*

- `kategorie`: LOW · `klasse`: „Satzbau: Rest einer Teilersetzung in einem Doku-Träger"
- `pfad`: `harness/sensors/coverage-gate.md:68-75` (neu, `a81a6ea`)
- `befund`: Der Satz beginnt „Über **acht** Läufe **desselben** Produktionsstands
  …" und hat sein Hauptverb erst in `:75` („**lag** die gedeckte Zahl zwischen
  1468 und 1471"). Dazwischen steht nach dem Semikolon ein **abgeschlossener
  Teilsatz** („gemessen ist die Differenz, die Aussage über die Code-Stände ist
  ihre Deutung"), danach die Beleg-Klammer (`go test …`, `Lauf slice-091`),
  danach das gestrandete „lag". Der Lauf-Beleg, der **die acht Läufe** bezeugt,
  steht damit am Ende des Deutungs-Satzes. Die Substanz trägt weiter (ein
  sorgfältiger Leser hängt „lag" zurück), deshalb LOW — und **nicht**
  `beleg-befehl-traegt-seinen-satz-nicht`: der genannte Befehl trägt seinen Satz,
  er steht nur eine Klausel zu früh.
- `verifizierbar`: ja — `git show a81a6ea`; die Zeilen `68-79` als ein Satz gelesen

### D-4 — Die zwei Commits fassen einen Träger an, den §3 nur **auf eine Bedingung** freigibt; die Ausweitung ist in keinem der beiden Commits benannt

- `kategorie`: INFO · `quelle`: Plan §3 `:149` · `verify-slice-092.md` **V-1** (erste Hälfte)
- `pfad`: `harness/sensors/coverage-gate.md` in `ad95444` **und** `a81a6ea`
- `befund`: Gemessen driftet **keine** Zahl der Datei (Nenner 1903 an allen drei
  gemessenen Ständen; die untere Kante `1369` exakt reproduziert). Beide Commits
  ändern die Datei trotzdem. `ad95444` nennt als Anlass F-2 (Prosa-Alterung),
  `a81a6ea` nennt `V-1` — aber **keiner** sagt, dass damit eine **Prosa-** statt
  einer **Zahl**-Änderung der §3-Anlass war, und §7 trägt dazu keine Zeile.
  Kein Rückgabe-Pfeil: die Entscheidung liegt bei der Closure.

### D-5 — „die zweithäufigste Klasse des Registers" trifft nur, wenn man „des Registers" als „der sieben in §8 genannten" liest

- `kategorie`: INFO · `pfad`: `slice-092-coverage-cluster-d1.md:115-116`
- `befund`: Gemessen: **drei** Klassen stehen bei **6×**, **drei** bei **5×**.
  Die mit **5×** genannte Klasse ist im Register **geteilt-vierte**, nicht
  zweite. Trägt die Formulierung als Aussage über die **sieben in §8 gelisteten**
  Einträge (6 · 5 · 4 · 3 · 3 · 2 · 1) — dort ist sie zweite. Die Wendung „des
  Registers" lädt die erste Lesart ein; die Zahl selbst (**5×**) ist richtig.
- `verifizierbar`: ja — `ls <BEO>/evidence | wc -l` über alle Einträge

---

## Negativbefunde

**geprüft, ohne Befund: die vier Größen — jede einzeln selbst nachgezählt.**

| Größe | Aussage | Eigener Lauf | Ergebnis |
|---|---|---|---|
| **13** | der Glob führt 13 Use-Case-Pakete | `ls -d internal/application/usecase/*/ \| wc -l` → **13** | **hält** |
| **10** | 10 davon trugen ungedeckte Statements | Profil Parent `2b7c6ec`: genau **10** mit `uncovered > 0` | **hält** |
| **22** | „(die 22)" | Summe jener 10 = **22**; `domain/model` **2** → **24** | **hält, exakt** |
| **12** | 12 haben Tests bekommen | `git diff --name-only 3de9547^..3de9547 \| grep -c "usecase/.*_test.go"` → **12** | **hält** |
| **11** | Deckung bewegt bei 11 Paketen | Parent → Head: die 10 auf `0`, `domain/model` 2 → 0; drei bleiben `0 → 0` | **hält** |

**geprüft, ohne Befund: die Platzierungen nehmen je die richtige Zahl.**
§1 `:52-55`, §2/LP1 `:106`, §3 `:147`, §5 `:193`, §6 `:217`, §8 `:297` — keine
Stelle nimmt die falsche der vier. Der Widerspruch aus F-1 (13 in §5/§6 gegen 11
in §8) ist **weg**: alle drei Bewegungs-Aussagen stehen auf **11**.

**geprüft, ohne Befund: „readchanges war bereits vollständig gedeckt"** — am
Parent `readchanges` `0/10`; die Klammer erklärt die Lücke **13 → 12** und steht
an der 12, nicht an der 10. Kein Finding.

**geprüft, ohne Befund: §3.7/§3.11 in den hinzugefügten Zeilen.** Die zwei
Vergangenheits-Formen stehen in **Plan-Prosa** (Unterscheidung der
Planungs-Größen), nicht in einem Kommentar und nicht in einem Zustandsfeld;
kein Slice-/Wellen-Vokabular als Begründung, kein host-lokaler Pfad, **keine**
neue `§`-Angabe. `make docs-check` **EC=0**: 759 Dateien, 0 Befunde.

**geprüft, ohne Befund: F-2 ist in der Sache erfüllt.**
`grep -rn "gegenständlichen"` → kein Treffer; die `+`-Zeilen tragen kein „hier",
„aktuell", „dieser Stand".

**geprüft, ohne Befund: Hygiene und Traceability.** 4 Pfade, kein Produktcode,
kein `THRESHOLD`, keine ADR, kein `git mv`; `make commit-traceability` **EC=0**,
beide Betreffe nennen `ADR-0082` ohne Struktur-ID; Arbeitsbaum leer.

**geprüft, ohne Befund: die Zahlenpaare des Trägers reproduzieren.**

| Stand | `go tool cover -func` | Nenner | gedeckt | ungedeckt |
|---|---|---|---|---|
| `3d5b60d` (Stand `slice-085`) | `71.9%` | **1903** | **1369** | 534 |
| `2b7c6ec` (Parent `slice-092`) | `77.2%` | **1903** | **1469** | 434 |
| `a81a6ea` (dieser Stand) | `78.5%` | **1903** | **1493** | 410 |

Die untere Kante des Trägers (**1369 von 1903**) ist damit **exakt
reproduziert**, der Nenner an **allen drei** Ständen **1903** — die Grundlage von
D-2. `TESTEXIT=0` in allen drei Läufen.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git rev-parse HEAD` · `git status --porcelain` | **0** | `a81a6ea`; Baum **leer** |
| `git diff --name-only 3de9547..a81a6ea` | **0** | **4** Pfade |
| `ls -d internal/application/usecase/*/ \| wc -l` | **0** | **13** |
| `git diff --name-only 3de9547^..3de9547 \| grep -c "usecase/.*_test.go"` | **0** | **12** |
| Profil `3d5b60d` / `2b7c6ec` / `a81a6ea` (Container, netzlos) | **0** (je) | 1903/1369 · 1903/1469 · 1903/1493 |
| `git diff --name-only 7835b6e..f90c3f4 -- internal cmd \| grep -v _test.go \| grep -v replication/receive` | **0** | **leer** — kein Produktcode zwischen den zwei Bänder-Ständen |
| `grep -rn "gegenständlichen" harness/sensors/coverage-gate.md` | **1** | kein Treffer |
| `ls <BEO>/evidence \| wc -l` über alle Einträge | **0** | 6× · 6× · 6× — dann 5× · 5× · 5× (D-5) |
| `make docs-check` (Log in Datei, EC danach gelesen) | **0** | 759 Dateien, 0 Befunde |
| `make commit-traceability` | **0** | OK, 5 Commits |

Nicht gefahren: `make gates`/`make test` für `a81a6ea` — der Vorgang ist
**doku-only** (4 Pfade, kein Produktcode, kein Workflow); das relevante Gate ist
`docs-check`, und es ist grün.

## Antwort auf die Schwerpunkte

**(1) Tragen die zwei Korrektur-Commits? — `ad95444`: in der Sache ja, mit einem
neuen Zählfehler im Kopf der Korrektur (D-1).** Die Zahlen-Nachzüge halten alle,
die Deixis ist fort, §5/§6/§8 widersprechen sich nicht mehr — **F-1 ist an
seinen Platzierungen erledigt**. Was der Fix selbst erzeugt hat, ist die
Überschrift „Drei Zahlen, drei Dinge" über vier Größen. Das ist das Muster, das
diese Slice-Folge kennt (`slice-090`, `slice-091`) — hier zum dritten Mal in
Folge und wieder an einem **Zähl-Wort**.

**(2) `a81a6ea`: der Marker ist gesetzt — und invertiert (D-2, MEDIUM).** Der Fix
benennt das **Falsche** als gemessen (die Differenz — §3.12 zählt „Differenz" zu
den abgeleiteten Werten) und schützt mit dem Wort „Deutung" eine Aussage, die
die Messung **ausschließt**: die zwei Bänder liegen auf demselben Code-Stand, die
~99–100 liegen in den **gedeckten** Zahlen, also am Test-Baum. Der Fix hat eine
alternde Form durch eine statische Falschaussage mit falsch gesetztem
Herkunfts-Etikett ersetzt.

**(3) „Hält kein Leser die Stelle mehr für eine Aussage über heute?" — nicht
vollständig.** Der Satz trägt **zwei** Träger: die *acht Läufe* (sauber
markiert) und die *Aussage über die zwei Code-Stände*. Die zweite steht **flach**
und wird erst eine Klausel später als „Deutung" markiert — und sie ist gegen die
Messung falsch.

**(4) Ist „alternde Prosa ohne Zahl-Drift" ein echter Fall von
`BEO-PGC/arbeit-ueberholt-stehenden-traeger`? Meine Antwort: Überdehnung.** Drei
Gründe:
1. **Die Grenze des Eintrags greift nicht.** Dort steht: „§3.12 hilft nicht: die
   Aussage trug ihren Ursprung korrekt — sie war zum Zeitpunkt ihrer
   Niederschrift wahr." Hier **hilft §3.12**: die Wendung „desselben, hier
   gegenständlichen Produktionsstands" hat eine Lauf-Größe als „der Ist-Stand"
   vorgeführt — genau die Form, die §3.12 bei Niederschrift verbietet.
2. **„Falsch durch diese Arbeit" trifft nicht zu.** `slice-092` ist test-only;
   der Produktionsstand der Aussage ist gemessen **derselbe** geblieben.
3. **Die Form des Erstauftretens ist nicht reproduziert.** Bei `slice-091`s
   `streamv1`-Liste war der Satz bei Niederschrift **wahr** und wurde durch die
   Arbeit **falsch** — das ist der BEO-Fall. Hier war nichts bei Niederschrift
   wahrheitsfähig gebunden, was die Arbeit umstoßen konnte.
**Konsequenz:** würde die `welle-20`-Closure diese Stelle als **zweites
Vorkommen** dateien, zählte sie einen **Formulierungs**-Defekt als Fall von
„Arbeit überholt stehenden Träger" — und der Eintrag warnt in seinem eigenen Text
genau davor. Zur **Weite der §3-Bedingung**: sie ist tatsächlich zu eng (sie
kennt nur den Zahl-Drift), aber was der Fall belegt, ist der **engere** Fall: die
Deixis-Form ist **ohne** Alterung prüfbar.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **0** |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen:** `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (D-1,
D-2 — **derselbe Vorgang** wie F-1/F-2: ein Vorgang zählt einmal) · „Satzbau:
Rest einer Teilersetzung in einem Doku-Träger" (D-3, neu) · §3-Bedingung nicht
benannt (D-4) · „Zähl-Wort einer Rangfolge" (D-5).

**Register-Hinweis (Modul 6):** Für `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
**kein** zweiter Beleg aus diesem Vorgang; die Stelle trägt ihren Fall in §3.12
(Instanz A, bereits verkörpert). Den Zähler setzt die Closure.

## Verdikt

**Merge-blockierend:** ja — **2 MEDIUM** (D-1, D-2). Kein HIGH.

**Tragen die zwei Korrektur-Commits? — Teilweise.** `ad95444` trägt in der
Substanz und erzeugt einen neuen Zählfehler (**D-1**). `a81a6ea` **trägt
nicht**: der Marker ist formal gesetzt und inhaltlich vertauscht, die „Deutung"
ist contra-gemessen (**D-2**), die Klausel syntaktisch gerissen (**D-3**), V-1s
erste Hälfte unbenannt (**D-4**).

**Rückgabe-Pfeil Reviewer → Implementer: nötig, kurz** — zwei Stellen:
1. `slice-092-coverage-cluster-d1.md:51` — das Zähl-Wort auf **vier** (D-1).
2. `harness/sensors/coverage-gate.md:68-73` — die Etiketten richtig setzen (die
   vier gedeckten Zahlen sind gemessen, die Differenz **abgeleitet**), die
   Aussage über „zwei verschiedene Code-Stände" auf das bringen, was gemessen ist
   (gleicher Nenner 1903, kein Produktcode geändert — die Verschiebung liegt an
   den **Tests**), und den Satz wieder schließen (D-2, D-3 in einem Zug).

D-4 und D-5 liefern **keinen** Rückgabe-Pfeil (Closure/§7 bzw. Träger außerhalb
des Deltas). **V-2 aus `verify-slice-092.md` ist mit diesem Report erfüllt** — er
deckt beide Fix-Commits ab.

**DoD-Häkchen „Review durchgeführt":** bleibt **offen** — der Slice braucht die
kurze Fixrunde.

**Was nicht geprüft wurde:** DoD/LP1–LP3 (Verifier) · die §6-Ausgänge, die
Register-Dateiung und die drei Paarungen (Planner-Closure) · der **volle**
Gate-Lauf für `a81a6ea` (doku-only; nur `docs-check` und `commit-traceability`
gefahren) · die obere Kante `1371` des Trägers · der reale Post-Push-Lauf
(§3.10: kein Workflow geändert).
