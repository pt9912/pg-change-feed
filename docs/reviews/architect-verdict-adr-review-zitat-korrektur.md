# Architect-Verdikt: `Accepted`-ADRs zitieren `docs/reviews/**` per Live-Link — Zitat-Korrektur statt Adresse

**Rolle:** Architect (Modul 8)

**Anlass:** Auftraggeber-Anweisung: 25 `Accepted`-ADRs verlinken direkt und
live in `docs/reviews/` (Review-/Verify-Reports, Architect-Verdikte).
Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur behandelt
Review-Reports als **Lauf-Belege ohne eigene Identität jenseits ihres
Slice** — sie wandern bei Archivierung nach `done/<welle-id>/archiv.zip`
und bekommen **keinen** Stub (anders als Slices/Wellen). Ein lebender
Markdown-Link aus einer permanent im Repo stehenden `Accepted`-ADR macht den
zitierten Report faktisch nie archivierbar, ohne die ADR zu brechen — real
reproduziert: `.harness/state/bin/ai-harness-init archive-welle --vorschau
welle-d-check` bricht mit einer `[haenger]`-Sperre ab, weil u. a.
`ADR-0062`, `-0069`, `-0070`, `-0075`, `-0083`, `-0084`, `-0086`, `-0087`,
`-0089`, `-0091`, `-0092`, `-0093` Review-Reports zitieren, die der Lauf
verschieben würde.

Zwei Nutzer-Nachträge während dieses Zugs verschärfen den Auftrag
gegenüber der ursprünglichen Fragestellung („ist es eine Zitat-Korrektur"):
(1) die Grundsatz-Position „diese Links sind nicht regelwerkkonform,
**alle** entfernen", (2) der Auftrag, den Befund **mechanisch** über eine
neue `matrix`-Regel in `.d-check.yml` durchzusetzen. Beide sind unten als
eigene Abschnitte beantwortet — mit eigener Architect-Begründung, nicht als
Übernahme der vorgeschlagenen Formulierung.

**Rolleninhaber:** pt9912 (Architect-Rolle; anderer Kontext als die
Implementer-/Reviewer-/Verifier-Läufe der 25 betroffenen ADRs — Modul 8
§Rollen-Regeln)

**Datum:** 2026-09-18

**Bezug:** `AGENTS.md` §3.5 (ADR-Immutabilität, Zitat-Korrektur-Ausnahme) ·
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(die Klasse „Zitat-Korrektur" und ihre Grenze) · Baseline-Regelwerk
`modul-06-roadmap.md` §Wellen-Closure-Prozedur (Review-Reports ohne Stub) ·
Baseline-Regelwerk `modul-08-agentenrollen.md` §Rollen-Regeln (Architect
schreibt, Accepted-ADRs überschreibt niemand) · `/Development/KI/ai-harness-init`
`internal/archive/scan.go` (der Hänger-Mechanismus, gelesen in diesem Zug) ·
`/Development/d-check` `docs/user/benutzerhandbuch.md` §4.7/§5 und
`spec/spezifikation.md` `DC-FA-MTX-001`…`003` (Modul `matrix`, gelesen in
diesem Zug)

**Geltungsbereich dieses Verdikts:** Es trifft **kein** ADR-Textänderung.
Es liefert der nachfolgenden Implementer-Fixrunde die vollständige
Bestandsaufnahme, das Verdikt je Fundstelle und die Ziel-Zitierform. Es
vergibt keine Slice-Kennung.

---

## 1. Der Bestand — vollständig, selbst verifiziert

`grep -rln "reviews/review-\|reviews/verify-\|reviews/architect-verdict-\|reviews/architect-review-" docs/plan/adr/*.md`
liefert **25** Treffer (Zahl bestätigt). Je ADR: Fundort-Klasse — **(a)**
`§Bezug`/`§Kontext`/Fließtext (Auslöser der Entscheidung, inkl. der
Analyse-Tabellen), **(b)** `§Geschichte`-Zeile (Ereignis-Text und/oder
`Verweis`-Spalte) — und die aktuelle Zitierform.

| ADR | Fundort(e) | Aktuelle Form (Kurzfassung) |
|---|---|---|
| [`0044`](../plan/adr/0044-image-beleg-semantik.md) | (a) `§Bezug`, `§Kontext` (2×), `§Entscheidung` Punkt 4 (2×) · (b) Geschichte-Ereignistext (2 Links) | Markdown-Links `[docs/reviews/verify-slice-004.md](...)`, `[docs/reviews/review-slice-005.md](...)`, `[review-slice-005 F-7](...)`, `[verify-slice-004.md](...)` (relativer Pfad `../../../docs/reviews/...`) |
| [`0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) | (a) `§Kontext` · (b) Geschichte-Verweis-Spalte | Markdown-Link `[docs/reviews/review-slice-010.md](../../reviews/review-slice-010.md)`, je mit `<!-- d-check:status-provenance -->` |
| [`0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) | (a) `§Bezug` | Markdown-Link `[architect-review-welle-6.md](../../reviews/architect-review-welle-6.md)` + Text „Zug 2" |
| [`0053`](../plan/adr/0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md) | (a) `§Bezug` (2×), `§Kontext` (Prosa, kein Link) | Markdown-Links `[architect-verdict-retention-loeschausfuehrung.md](...)`, `[architect-verdict-slice-044-rollen-grant.md](...)`; in `§Kontext` derselbe Pfad als Inline-Code-Bare-Pfad `(docs/reviews/architect-verdict-retention-loeschausfuehrung.md)` |
| [`0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md) | (a) `§Kontext` (2×, Inline-Code-Bare-Pfad) · (b) Geschichte-Verweis-Spalte (Inline-Code-Bare-Pfad) | `` `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md` `` (Bare-Pfad, kein Link), Geschichte nennt eine **abweichende** Datei `-gegengeprueft.md` |
| [`0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (2 Bare-Pfade) | `` `docs/reviews/review-slice-067.md` ``, `` `docs/reviews/architect-verdict-spaltenausschluss-dauerhaftigkeit.md` `` |
| [`0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md) | (a) `§Bezug` (Bare-Pfad), `§Konsequenzen` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (2 Bare-Pfade) | `` `docs/reviews/review-slice-069.md` ``, `` `docs/reviews/architect-verdict-publish-blockiert-capture-pfad.md` `` |
| [`0067`](../plan/adr/0067-capture-publish-einbindung-fitness-function-korrektur.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/review-slice-070.md` `` |
| [`0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md) | (a) `§Bezug` (Bare-Pfad), `§Kontext` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/review-slice-071.md` `` |
| [`0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md) | (a) `§Bezug` (Link + Bare-Pfad) · (b) Geschichte-Ereignistext (Link) und Verweis-Spalte (Link + Bare-Pfad) | `[docs/reviews/review-slice-073.md](../../reviews/review-slice-073.md)`, `` `docs/reviews/architect-verdict-commit-msg-hook-einseitige-zusage.md` `` |
| [`0070`](../plan/adr/0070-supersede-reichweite-und-klassengrenze.md) | (a) `§Bezug` (2 Bare-Pfade), `§Kontext` (2 Bare-Pfade) · (b) Geschichte-Verweis-Spalte (2 Links) | `` `docs/reviews/review-slice-073-fixrunde.md` ``, `` `docs/reviews/verify-slice-073.md` ``; Geschichte-Zeile nutzt Links auf dieselben Dateien |
| [`0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/architect-verdict-coverage-gate-messgegenstand.md` `` |
| [`0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) | (a) `§Bezug` (Bare-Pfad), `§Kontext`-Messtabelle (komprimierte Bare-Pfad-Notation „`review-slice-036/039/049/073.md`") · (b) Geschichte-Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md` `` |
| [`0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md` `` — **Vorsicht:** diese ADR selbst definiert die Zitat-Korrektur-Klasse; ihre eigene Korrektur ist ein Selbstanwendungsfall (s. u. §4) |
| [`0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md` `` |
| [`0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) | (a) `§Bezug` (2 Links), `§Kontext` (Link) · (b) Geschichte-Ereignistext (Link) und Verweis-Spalte (Link) | `[review-slice-078](../../reviews/review-slice-078.md)`, `[architect-verdict-slice-078-konfliktpfad.md](...)` |
| [`0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | (a) `§Bezug` (4 Bare-Pfade) · (b) Geschichte-Ereignistext (2 Bare-Pfade) und Verweis-Spalte (3 Bare-Pfade) | `` `docs/reviews/review-slice-081.md` ``, `` `-084.md` ``, `` `verify-slice-085.md` ``, `` `review-slice-088.md` ``, `` `review-slice-088-delta.md` `` |
| [`0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md) | (a) `§Bezug` (3 Bare-Pfade) · (b) Geschichte-Ereignistext/Verweis-Spalte (2 Bare-Pfade) | `` `docs/reviews/review-slice-069.md` ``, `` `verify-slice-069.md` ``, `` `review-slice-074.md` `` |
| [`0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Ereignistext/Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/review-slice-093.md` `` |
| [`0086`](../plan/adr/0086-herkunft-aussagen-schwere-folgt-der-konsequenz.md) | (a) `§Bezug` (7 Bare-Pfade) · (b) Geschichte-Ereignistext/Verweis-Spalte (5 Bare-Pfade) | `` `docs/reviews/review-slice-081.md` `` … `` `review-slice-094.md` `` (sieben Reports) |
| [`0087`](../plan/adr/0087-beispiel-clients-csharp-kotlin.md) | (b) Geschichte-Ereignistext, zweite Zeile (Digest-Korrektur) | `docs/reviews/review-slice-099.md F-1` (Bare-Pfad, kein Backtick) |
| [`0089`](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Ereignistext/Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/review-slice-096.md` `` |
| [`0091`](../plan/adr/0091-zugangsdaten-klasse-sechs-schluessel.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Ereignistext/Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/verify-slice-096.md` `` |
| [`0092`](../plan/adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md) | (a) `§Bezug` (Bare-Pfad) · (b) Geschichte-Ereignistext/Verweis-Spalte (Bare-Pfad) | `` `docs/reviews/verify-slice-096.md` `` |
| [`0093`](../plan/adr/0093-digest-korrektur-adr-0087-kotlin-basis-image.md) | (a) `§Autor`-Feld (Bare-Pfad), `§Kontext` (Bare-Pfad) · (b) Geschichte-Ereignistext (Bare-Pfad) | `` `docs/reviews/review-slice-099.md` `` |

**Ergebnis der Zählung:** 25/25 tragen Fundort (b) — jede betroffene ADR
zitiert den auslösenden Review in ihrer eigenen `§Geschichte`-Zeile (Beleg
für die Statusänderung). 22/25 tragen zusätzlich Fundort (a) — nur
`ADR-0087` trägt den Review-Bezug **ausschließlich** in einer *zweiten*
Geschichte-Zeile (Digest-Korrektur, kein `§Bezug`/`§Kontext`-Vorkommen).
**Keine ADR trägt ausschließlich (a)** — der Beleg in der
`§Geschichte`-Zeile ist durchgängig vorhanden, weil `AGENTS.md` §3.5 (über
die Baseline-ADR-Vorlage) für jede Statusänderung einen Beleg verlangt.

## 2. Technischer Befund — was der Hänger-Scan tatsächlich erkennt

`internal/archive/scan.go` (`/Development/KI/ai-harness-init`) läuft **rein
über Basisnamen-Teilstring-Vergleich**, nicht über Link-Auflösung:

```go
func treffer(inhalt, datei string, ziele []string) []string {
    for _, ziel := range ziele {
        if strings.Contains(inhalt, filepath.Base(ziel)) {
            out = append(out, datei+" -> "+ziel)
        }
    }
    ...
}
```

`Haenger()` läuft über `Suchraum(dateien)` — **nicht** über
`SuchraumNachzug(dateien)` — und `docs/plan/adr` ist in
`AusgenommenePfade()` **nicht** enthalten (nur in
`AusgenommenePfadeNachzug()`, die ausschließlich die beiden
NACHZUG-Leser `VerweisFund`/`Nachziehen` betrifft, die eine Accepted-ADR
ohnehin nie schreibend anfassen dürfen). Das bedeutet konkret:

1. **Jede** Erwähnung des Datei-Basisnamens zählt — Markdown-Link,
   Inline-Code-Bare-Pfad, nackter Fließtext. Der Scan unterscheidet nicht
   zwischen „verlinkt" und „nur erwähnt". Delinken allein (Klammern
   entfernen, Pfad-Text als Prosa stehen lassen) **löst den Hänger nicht**
   — der Basisname bleibt Teilstring des Dateiinhalts.
2. Der Suchraum für `Haenger` deckt **alle** getrackten Dateien
   (`git ls-files`), nicht nur `.md` — es gibt keine Dateityp-Achse
   (Kommentar in `scan.go` Zeile 21–24 nennt explizit Shell-Hooks,
   Go-Kommentare, bats-Dateien).
3. Für die 25 ADRs folgt daraus zwingend: Die Ziel-Zitierform darf den
   **exakten Datei-Basisnamen des Reviews an keiner Stelle mehr enthalten**
   — auch nicht als reine Erwähnung ohne `[]()`-Klammern. Ersetzt werden
   muss durch eine Kennung, die den Basisnamen nicht reproduziert (siehe
   §5).
4. Wichtig für den Implementer: Dies **entsperrt** `archive-welle` für die
   25 ADRs, aber **nicht notwendig vollständig** — der Scan läuft
   repo-weit über alle getrackten Dateien. Bare Erwähnungen desselben
   Review-Basisnamens in anderen `.md`-Dateien (`done/`-Records,
   `docs/reviews/**` selbst, `harness/**`) oder in Nicht-Markdown-Dateien
   bleiben unberührt von dieser Fixrunde und müssten separat geprüft
   werden, bevor `archive-welle --vorschau welle-d-check` tatsächlich grün
   läuft. Dieser Zug deckt ausschließlich die 25 ADRs.

## 3. Die Grenze nach `ADR-0073` — Zitat-Korrektur oder nicht

`ADR-0073` §Entscheidung Punkt 1 definiert die Klasse eng: eine Änderung ist
eine Zitat-Korrektur, wenn sie **ausschließlich die Form eines Verweises
auf einen unveränderten Referenten** ändert (Pfade, Linkziele,
Lokatoren, Form einer gebrochenen Referenz). Unberührbar: §Entscheidung,
§Konsequenzen *der Aussage nach*, §Verglichene Alternativen, §Status, die
`Supersedes`-Kette, Fitness-Function-Regeln, Re-Evaluierungs-Trigger,
`Datum`/`Autor`, die **Aussage-Semantik** von §Bezug/§Schärft.

**Verdikt für alle 25 Fälle: Zitat-Korrektur ist durchgängig zulässig —
für beide Fundort-Klassen (a) und (b), für alle 25 ADRs, ohne
Ausnahme.** Begründung:

- **Der Referent ändert sich nirgends.** In jedem der 25 Fälle bleibt
  klar, *welcher* Review/Verdikt gemeint ist — nur die Zitierform wechselt
  von „Adresse" (Pfad/Link) auf „Kennung" (Slice-Nummer + Finding-ID bzw.
  thematische Kurzform).
- **Die Aussage-Semantik ändert sich nirgends.** In allen 25 Fällen steht
  der substanzielle Inhalt des zitierten Findings bereits **im
  Fließtext der ADR selbst** — als Beispiel `ADR-0044`: „Binary an beiden
  Range-Grenzen bit-identisch (`43c3aec0…`), Digest wechselte
  (`447eab36…` → `9ac4a9fb…`)" steht bereits ausformuliert neben dem Link;
  der Link fügt dieser Aussage keine Information hinzu, die nach seiner
  Entfernung fehlen würde. Dasselbe Muster (Finding-ID + Kurzcharakterisierung
  bereits in Prosa) gilt für alle 25 ADRs — durchgängig, weil das
  MADR-Kontext-Muster dieses Repos verlangt, den Anlass **auszuformulieren**,
  nicht nur zu verlinken (Baseline-Regelwerk `modul-04-adrs.md`
  §Ziel-Form: ADR).
- **Kein Fall verlangt neuen Inhalt.** Die eingangs vermutete
  Grenzverletzung — „eine ADR verlinkt NUR, ohne den Inhalt
  zusammenzufassen, und die Korrektur müsste die Zusammenfassung als
  neuen Inhalt hinzufügen" — tritt **nicht** ein. Selbst die knappsten
  Belege (`ADR-0047` `§Bezug`: „Architect-Verdikt
  [`architect-review-welle-6.md`](...) Zug 2 (`BEO-PGC/rollen-verdrahtung`,
  3×, Ausgang `geplant` → `slice-023`)") tragen die Kennung („Zug 2" +
  BEO-Registereintrag) bereits **neben** dem Link — die Korrektur streicht
  den Pfad, nicht die Aussage.
- **Der `§Geschichte`-Beleg verlangt keine Adresse.** `AGENTS.md` §3.5
  (über die Baseline-ADR-Vorlage) verlangt einen Beleg je Statusänderung —
  nicht zwingend einen Pfad-Link. Präzedenz **im eigenen Bestand**:
  `ADR-0072`s zweite Geschichte-Zeile trägt als `Verweis` bereits den
  bloßen Commit-Hash `df47282`, keinen Pfad. Ein Beleg als Kennung
  (`slice-<NNN>-Review, Finding F-<N>`) oder Commit-Hash ist eine
  gleichwertige, im Repo bereits geübte Form — dieselbe „Kennung, nicht
  Adresse"-Regel wie in
  `.harness/baseline/v6.9.0/templates/docs/plan/planning/archiv-stub-slice.template.md`
  §Zitier-Form.
- **Sonderfall `ADR-0073` selbst** (self-application, §Geschichte Zeile 1):
  Diese ADR **definiert** die Zitat-Korrektur-Klasse und wendet sie in
  ihrer eigenen zweiten Geschichte-Zeile (Zeile 282, `df47282`) bereits auf
  sich selbst an — ein Präzedenzfall, kein Hindernis. Die Korrektur ihrer
  ersten Geschichte-Zeile (Zeile 192, zitiert den hostpaths-Verdikt) folgt
  demselben Muster wie jede andere Zitat-Korrektur hier; es ändert nichts
  an der **normativen** Aussage der ADR (§Entscheidung, §Konsequenzen).

**Es gibt in diesem Bestand keinen echten Grenzfall.** Die eingangs
befürchtete Kategorie — Zitat ohne begleitende Zusammenfassung, sodass die
Korrektur neuen Inhalt einführen müsste — kommt unter den 25 ADRs nicht
vor. Zwei Stellen verdienen dennoch redaktionelle Sorgfalt beim
Umschreiben (kein Grenzfall, aber eine Fußnote):

1. **`ADR-0062`** (§Geschichte): Die Geschichte-Zeile nennt eine
   **andere** Datei (`architect-verdict-…-gegengeprueft.md`) als der
   `§Bezug`/`§Kontext` (`architect-verdict-…-kein-vorab-hook.md`, ohne
   `-gegengeprueft`). Das ist ein bereits **bestehender** inhaltlicher
   Unterschied (zwei verschiedene Architect-Verdikt-Dokumente im selben
   Konflikt-Pfad), keine Folge der Zitat-Korrektur — der Implementer hält
   die Unterscheidung in der Kennung-Form aufrecht (z. B. „Architect-Verdikt,
   erster Zug" vs. „…, Gegenprüfung"), ändert an der Tatsache selbst
   nichts.
2. **`ADR-0072`** (§Kontext-Tabelle, Zeile 66): Die komprimierte Notation
   „`review-slice-036/039/049/073.md`" ist bereits eine Abkürzung für vier
   Dateien (keine reale Pfadangabe). Die Zitat-Korrektur ersetzt sie durch
   die vier Slice-Kennungen (`slice-036`, `-039`, `-049`, `-073`) ohne
   `review-`-Präfix und ohne `.md`-Endung — dieselbe Information (vier
   Review-Reports, sieben Befunde), andere Form.

## 4. Was NICHT angefasst werden darf

Unverändert bleiben in allen 25 ADRs: `§Entscheidung`, `§Konsequenzen`
(der Aussage nach), `§Verglichene Alternativen`, `§Status`, die
`Supersedes`-Kette, jede `Fitness Function`-Regel, jeder
`Re-Evaluierungs-Trigger`, `Datum`/`Autor`. Die **Zahlen, Funde und
Schiedssprüche**, die aus den zitierten Reviews in die ADR-Prosa
übernommen wurden (z. B. `43c3aec0…`, „sechs Zugangsdaten-Schlüssel",
„siebzehn von neun Schlüsseln") bleiben **wortgleich** stehen — nur der
Pfad/Link, der auf die Quelldatei zeigt, wird durch die Kennung ersetzt.
Der `<!-- d-check:status-provenance -->`-Marker, der in den meisten
Fundstellen unmittelbar neben einem `slice-<NNN>`-Token steht, bleibt an
diesem Token erhalten — er gehört zur bestehenden `matrix`-Regel
`{from: adr, to: slice, allow: false}` (s. §6) und ist von dieser
Zitat-Korrektur unberührt; er verschwindet nur dort, wo mit dem Pfad auch
der einzige slice-Token der Zeile verschwindet (dann ist der Marker
gegenstandslos und wird mit entfernt, nicht stehengelassen als
Deko-Kommentar — §3.7).

## 5. Ziel-Zitierform

**Regel:** Kein Fließtext, keine `§Bezug`-Zeile und keine
`§Geschichte`-Zeile einer `Accepted`-ADR nennt den Datei-Basisnamen (mit
oder ohne `docs/reviews/`-Präfix, mit oder ohne `.md`-Endung, mit oder
ohne Link-Klammern) eines Review-/Verify-Reports oder Architect-Verdikts.
Stattdessen:

- **Reports mit Slice-Bezug** (`review-slice-NNN.md`,
  `verify-slice-NNN.md`, `review-slice-NNN-delta.md`,
  `review-slice-NNN-fixrunde.md`): zitiert als **„Review zu `slice-NNN`,
  Finding F-`N`"** bzw. „Verifikationsbericht zu `slice-NNN`, Befund
  V-`N`" — der `slice-NNN`-Teil trägt (wie bisher) den
  `<!-- d-check:status-provenance -->`-Marker, wo die bestehende
  `matrix`-Regel ihn verlangt.
  Beispiel (`ADR-0044`, `§Kontext`): statt
  „[`review-slice-005 F-7`](../../../docs/reviews/review-slice-005.md)
  widerlegt die Formel …" neu: „Review zu `slice-005`, Finding F-7,
  widerlegt die Formel …".
- **Architect-Verdikte ohne Slice-Nummer**
  (`architect-verdict-<thema>.md`, `architect-review-welle-N.md`):
  zitiert über ihr **Thema in Prosa**, ggf. ergänzt um den zugehörigen
  Commit-Hash oder BEO-Registereintrag, wenn einer bereits im Text steht
  (kein neuer Inhalt — nur die bereits vorhandene Kennung bleibt, der Pfad
  entfällt). Beispiel (`ADR-0053`, `§Bezug`): statt „Architect-Verdikt
  [`architect-verdict-retention-loeschausfuehrung.md`](...)  (Frage 1 —
  prüfte …)" neu: „der vorausgehende Architect-Verdikt zur
  Retention-Löschausführung (Frage 1 — prüfte …)".
- **`§Geschichte`-Verweis-Spalte:** Kennung statt Pfad, wahlweise ergänzt
  um den Commit-Hash der Korrektur, exakt wie in `ADR-0072`s eigener
  zweiten Zeile bereits geübt (`df47282` ohne Pfad).

Diese Form fügt **keinen** Inhalt hinzu, den die ADR nicht bereits trägt —
sie entfernt die Adresse und behält die Kennung, die in praktisch allen 25
Fällen ohnehin schon unmittelbar neben dem Link steht (Finding-ID,
Slice-Nummer, Themen-Slug).

## 6. Beleg der Korrektur

Wie in `ADR-0073` Punkt 3 vorgeschrieben: (a) die korrigierende
Commit-Message nennt `ADR-0073`; (b) jede der 25 betroffenen `Accepted`
ADRs erhält **eine** neue Zeile ihrer eigenen `§Geschichte`-Tabelle:
„Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt
(`ADR-0073`)" mit dem Commit-Hash als `Verweis` — dieselbe Form, die
`ADR-0072` bereits für seine eigene Pfad-Korrektur nutzt (Zeile 282,
`df47282`). Records (hier: keine betroffen, alle 25 Fälle sind ADRs)
bekämen keine solche Zeile.

## 7. Mechanisierung — neue `matrix`-Regel in `.d-check.yml`

**Auftrag:** Diese Fehlerklasse (ADR verlinkt live in `docs/reviews/`)
künftig mechanisch über `make docs-check` erkennen, nicht nur einmalig von
Hand bereinigen.

### 7.1 Modul-Semantik geprüft

Gelesen: `/Development/d-check/docs/user/benutzerhandbuch.md` §4.7/§5,
`/Development/d-check/spec/spezifikation.md` `DC-FA-MTX-001`…`003`. Ergebnis:

- Eine `matrix.classes[]`-Klasse **ohne `token`** ist zulässig — bereits
  Präzedenz in diesem Repo (`spec`-Klasse in `.d-check.yml` trägt keinen
  `token`) — und deckt dann **nur echte Markdown-Links** (`[text](pfad)`).
- Eine Klasse **mit `token`** deckt zusätzlich **bare Text-Vorkommen**
  (auch in Inline-Code, aber **nicht** in Fenced-Code-Blöcken oder
  innerhalb eines bereits gezählten Links) — genau die Form, in der die
  meisten der 25 Fundstellen tatsächlich stehen (Inline-Code-Bare-Pfad,
  kein `[]()`-Link, siehe §1-Tabelle: 18 von 25 ADRs nutzen überwiegend
  Bare-Pfade statt Links).
- Der Provenance-Marker `<!-- d-check:status-provenance -->` ist der
  bestehende Escape-Hatch für ein **erlaubtes**, aber sonst durch eine
  `token`-Regel gefangenes Vorkommen (aktuell für die `slice`-Klasse
  genutzt, sichtbar an den vielen Markern neben `slice-<NNN>`-Token in
  genau den 25 ADRs dieses Bestands).

### 7.2 Pfad-Deckung geprüft

`find docs/reviews -type d` → **keine Unterverzeichnisse** (flach, ein
Level). `ls docs/reviews | grep -vE
'^(review-|verify-|architect-verdict-|architect-review-)'` → genau **eine**
Ausnahme, `blocker-slice-063.md` (kein Review-Report im engeren Sinn,
sondern ein Blocker-Protokoll). Ein `paths`-Glob `["docs/reviews/*.md"]`
(ohne Präfix-Enumeration) erfasst **alle** 291 Dateien des Verzeichnisses
flach und robust gegen künftige Namensvarianten — dieselbe Glob-Form wie
die bestehenden `spec`/`adr`/`slice`-Klassen. Eine Enumeration der vier
Präfixe wäre enger, aber unnötig eng: `blocker-slice-063.md` ist ebenfalls
ein Lauf-Beleg ohne Stub-Anspruch und gehört in dieselbe Archivierungs-Klasse.

### 7.3 Konkreter Diff-Vorschlag (nicht angewendet — Implementer-Fixrunde)

```yaml
matrix:
  classes:
    - name: spec
      paths: [spec/lastenheft.md, spec/pflichtenheft.md, spec/architecture.md]
    - name: adr
      paths: ["docs/plan/adr/[0-9]*.md"]
      token: 'ADR-\d{4}'
    - name: slice
      paths: ["docs/plan/planning/**/slice-*.md"]
      token: 'slice-\d{3}'
    - name: review                                    # NEU
      paths: ["docs/reviews/*.md"]                     # NEU
      token: 'docs/reviews/[\w.-]+\.md'                 # NEU — fängt auch Bare-Pfad-Erwähnungen ohne Link
  rules:
    - {from: spec, to: adr, allow: false}
    - {from: spec, to: slice, allow: false}
    - {from: adr, to: slice, allow: false}
    - {from: adr, to: review, allow: false}            # NEU
```

**Bewusst keine neue Regel für:**

- `slice → review`: Ein Slice zitiert seinen eigenen Review-Report
  routinemäßig (DoD-Haken, Closure) — beide werden **gemeinsam**
  archiviert (Slice bekommt einen Stub, sein Review-Report wandert
  zeitgleich mit ihm ins Archiv), kein Hänger-Risiko.
- `review → review`: Delta-Reviews verweisen auf ihren Ausgangs-Report
  (`review-slice-088-delta.md` → `review-slice-088.md`); beide liegen in
  derselben Archiv-Einheit.
- `adr → adr` bleibt unverändert erlaubt (Supersedes-Kette,
  Bezug-Verweise).

**Bewusst kein Provenance-Marker-Escape für `adr → review`:** Anders als
bei der `slice`-Klasse (deren Mitglieder beim Archivieren einen Stub
behalten und daher als reine Provenance-Kennung dauerhaft auflösbar
bleiben) gibt es für `review`-Klassen-Mitglieder **keine** dauerhafte
Adresse — genau das ist der Grund für diese ganze Fixrunde (§Wellen-Closure-
Prozedur, kein Stub). Ein Escape-Hatch würde die Klasse aushöhlen, die
dieser Zug gerade schließt: Jede künftige `adr → review`-Referenz, auch
eine „nur als Provenance gemeinte", soll über die Kennung-Form (§5) laufen,
nicht über den Marker.

### 7.4 Reihenfolge — zwingend

Erst die 25 ADRs bereinigen (diese Fixrunde), **dann** die `matrix`-Regel
aktivieren. In umgekehrter Reihenfolge liefert `make docs-check` sofort 25+
`matrix-forbidden`-Befunde und das Gate wird rot, bevor der Implementer
überhaupt zu schreiben beginnt. Die Implementer-Fixrunde führt beide
Schritte **in derselben Runde**, aber innerhalb ihrer eigenen Reihenfolge:
(1) alle 25 ADR-Texte umschreiben, (2) `.d-check.yml` um die `review`-Klasse
und -Regel ergänzen, (3) `make docs-check`/`make gates` grün, erst dann
Commit/Closure.

### 7.5 Braucht diese Matrix-Erweiterung eine eigene ADR?

**Verdikt: Ja — abweichend von der eingangs vermuteten Einschätzung.**
Begründung, eigenständig hergeleitet, nicht übernommen: `AGENTS.md` §3.6
verlangt eine ADR nur für **Schwellen-Senkungen** — diese Änderung senkt
nichts, sie verschärft. Trotzdem zeigt der **Bestand dieses Repos**
durchgängig ein anderes, stärkeres Muster: **jede** bisherige Aktivierung
eines neuen `d-check`-Moduls oder einer neuen `matrix`/Gate-Regel in
diesem Repo trägt ihre eigene ADR — `ADR-0045` (Modul `commits`),
`ADR-0072` (Modul `hostpaths`), `ADR-0084` (Sync-Gate) — unabhängig davon,
ob die Änderung eine Lockerung oder eine Verschärfung war. Der Grund ist
in `ADR-0073` selbst benannt (§Konsequenzen, „Negativ mit Grenze"): die
Klassengrenze einer Zitat-Korrektur ist „ein Urteil, kein Sensor" — die
neue `matrix`-Regel **mechanisiert genau diesen Teil** des Urteils (das
Vorhandensein einer Adresse), nicht das ganze Urteil (ob eine Umformulierung
die Aussage-Semantik wahrt, bleibt Review-Prüfpflicht, `AGENTS.md` §3.12
Instanz B). Eine Fitness-Function-Ergänzung, die eine bisher nur als
Review-Prüfpflicht deklarierte Grenze teilweise in ein Gate hebt, ist
strukturell derselbe Vorgang wie `ADR-0072`s Aktivierung von `hostpaths` —
und verdient denselben Träger: eine schlanke Folge-ADR zu `ADR-0073`
(`Supersedes` nicht nötig, da `ADR-0073`s eigene Klausel unverändert
bleibt — die neue ADR ergänzt nur deren Fitness-Function-Zeile, die bisher
„—" lautet). Diese ADR ist **nicht** Gegenstand dieses Verdikts; sie ist
eine Folgepflicht an die nächste Architect-Rolle, spätestens in derselben
Fixrunde wie die `.d-check.yml`-Änderung.

## 8. Zusammenfassung für den Implementer

1. 25 `Accepted`-ADRs (Liste §1) — in **allen** 25 ist eine Zitat-Korrektur
   nach `ADR-0073` zulässig, für **beide** Fundort-Klassen (a) und (b).
   Kein echter Grenzfall; zwei redaktionelle Fußnoten (§3).
2. Zielform: Kennung (`slice-NNN` + Finding-ID, oder Themen-Slug für
   Architect-Verdikte ohne Slice-Nummer) statt Pfad/Link/Bare-Pfad — der
   Datei-Basisname darf **an keiner Stelle** mehr im ADR-Text stehen
   (§2, §5).
3. Unberührbar bleiben §Entscheidung, §Konsequenzen (Aussage), §Alternativen,
   §Status, Supersedes-Kette, Fitness Function, Re-Evaluierungs-Trigger,
   Datum/Autor, sowie jede übernommene Zahl/Finding-Aussage (§4).
4. Jede korrigierte ADR bekommt eine neue `§Geschichte`-Zeile
   („Zitat-Korrektur … `ADR-0073`", Commit-Hash als Verweis, §6); die
   Commit-Message nennt `ADR-0073`.
5. `.d-check.yml`: neue `review`-Klasse (`paths: ["docs/reviews/*.md"]`,
   `token: 'docs/reviews/[\w.-]+\.md'`) und Regel `{from: adr, to: review,
   allow: false}` — **nach** Schritt 1–4, sonst rotes Gate (§7.3, §7.4).
   Kein Escape-Marker für diese Regel (§7.3).
6. Folgepflicht: eine schlanke Folge-ADR zu `ADR-0073` für die
   Matrix-Mechanisierung (§7.5) — nicht Teil dieser Fixrunde, aber vor
   deren Abschluss zu benennen.
7. Diese Fixrunde entsperrt `archive-welle` für die 25 ADRs, aber nicht
   notwendig repo-weit — andere Träger (`done/`-Records, `docs/reviews/**`
   selbst, Nicht-Markdown) können denselben Review-Basisnamen noch
   erwähnen und bleiben außerhalb dieses Verdikts (§2 Punkt 4).

## Fitness Function (dieses Verdikt)

| Prüfung | Ergebnis |
|---|---|
| `grep -rln "reviews/review-\|reviews/verify-\|reviews/architect-verdict-\|reviews/architect-review-" docs/plan/adr/*.md \| wc -l` | 25 (verifiziert in diesem Zug) |
| `find docs/reviews -type d` | keine Unterverzeichnisse (verifiziert) |
| `ls docs/reviews \| grep -vE '^(review-\|verify-\|architect-verdict-\|architect-review-)'` | genau `blocker-slice-063.md` (verifiziert) |
| `internal/archive/scan.go` `Haenger`/`treffer` | reiner Basisnamen-Teilstring-Vergleich, kein Link-Parsing, Suchraum ohne ADR-Ausnahme (gelesen, §2) |
| `.d-check.yml` matrix-Erweiterung angewendet | **nein** — Vorschlag für die Implementer-Fixrunde (§7.3) |
