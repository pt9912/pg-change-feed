# ADR-0075: `hostpaths`-Regel — Reichweite, §3.11-Entwurf und Lokator-Disposition

**Status:** Accepted — Supersedes [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
in genau **einer** Klausel: §Entscheidung Punkt 5 (der Entwurf der Hard Rule
§3.11). Alles Übrige aus `ADR-0072` — die Aktivierung ohne Ausschlussliste
(Punkt 1), kein Ventil (Punkt 2), die Korrektur aller Stellen (Punkt 4), der
Config-Block, die Grenzen, die Fitness Function und die Re-Evaluierungs-Trigger
— bleibt unverändert. §Entscheidung Punkt 3 aus `ADR-0072` ist durch
[`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md) supersedet; die
Disposition seiner Lokator-Klausel stellt Entscheidung Punkt 3 dieser ADR fest.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Implementer-Lauf,
der die Korrektur ausgeführt hat, und als der Reviewer-Lauf —
`modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz)

**Bezug:** [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) (in
einer Klausel supersedet), [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(die Zitat-Korrektur-Klasse), [`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md)
(die Hausform; Supersedes von `ADR-0072` Punkt 3) · `AGENTS.md` §3.5
(Accepted-ADRs immutable — Folge-ADR), §3.7 (Ist-Zustand), §3.11 (der Entwurf)
· [`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md)
(Modul-Semantik und Grenzen) · Review zu `slice-078` <!-- d-check:status-provenance -->
(Anlass) · der Architect-Verdikt dieses Zugs (Konflikt-Pfad `slice-078`) <!-- d-check:status-provenance -->

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie
[`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Regel „kein host-lokaler absoluter Pfad in der Doku" ist ausgeliefert: das
`hostpaths`-Modul ist ohne Ausschlussblock aktiv, die Korrektur ergibt 0
Befunde. Die Review des auslösenden Vorgangs
(Review zu `slice-078`) <!-- d-check:status-provenance --> stellt drei
Regelfragen, die keine Implementer-Frage sind.

**Erstens — die Reichweite.** Der Slice-Plan (ein Zeitdokument, Modul 5/6)
liest die Regel als „kein Host-Pfad im Repo, Fences eingeschlossen"; die Hard
Rule `AGENTS.md` §3.11 und ihr Entwurf in `ADR-0072` §Entscheidung Punkt 5
lesen dieselbe Regel als „nur Prosa + Inline-Code, Fenced-Beispiele erlaubt".
Zwei Reichweiten einer Regel; keine ist als Gewinner deklariert.

**Zweitens — der §3.11-Entwurf.** Der Unterabschnitt `Entwurf der Hard Rule
§3.11` liegt in `ADR-0072` §Entscheidung Punkt 5 und ist dort ausdrücklich
„Teil dieser Entscheidung". Ein korrigierender Commit (`df47282`) hat ihn
in-place geändert: das Falsch-Beispiel auf einen Platzhalter, das
Richtig-Beispiel auf die Hausform, den erklärenden Klammertext auf eine
Form-Beschreibung. `AGENTS.md` §3.5 und `ADR-0073` Punkt 1 nehmen
§Entscheidung von der Zitat-Klasse aus, und der Klammertext ist keine Form
eines Verweises — der Eingriff hat die Klassengrenze überschritten.

**Drittens — die Lokator-Klausel.** `ADR-0072` §Entscheidung Punkt 3 ordnete
an, Zeilen-/Bereichs-Lokatoren durch den stabilen benannten Anker zu ersetzen.
[`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md) hat Punkt 3 als
Ganzes supersedet; sein Ersatztext (die Hausform) restituiert die
Lokator-Klausel nicht und führt sie auch nicht unter „Der Rest bleibt" auf —
ihre Disposition ist nirgends ausgesprochen.

**Zählung.** `ADR-0072` nennt 31 Stellen (die Modul-Befunde zum
Entscheidungszeitpunkt); die Korrektur umfasst 42 Vorkommen (Fenced-,
Nicht-Markdown- und danach erzeugte Stellen). Die 31 sind die
Entscheidungs-Zahl der immutablen ADR, die 42 der ausgeführte Umfang.

## Entscheidung

Wir deklarieren die Reichweite der Regel, fassen den §3.11-Entwurf neu und
stellen die Lokator-Disposition fest.

1. **Reichweite — Fences eingeschlossen.** Die Regel gilt für die ganze
   Markdown-Fläche, Fenced-Code-Blöcke eingeschlossen. Der Ort der Regel ist
   die Hard Rule `AGENTS.md` §3.11; ein Slice-Plan ist ein Zeitdokument und
   trägt keine Regel. **Eine verbotene Form zeigt ein Dokument ausschließlich
   als Platzhalter** — das Wurzel-Segment als `<Host-Wurzel>`, kein reales
   Segment; die Richtig-Form ist die Hausform (`` `d-check`s `Dockerfile` ``,
   [`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md)); die reale
   Form steht nirgends, auch nicht im Fence.

   **Konsequenz — die Regel ist strenger als ihr Sensor.** Das
   `hostpaths`-Modul lässt Fenced-Blöcke frei (Modul-Design, ohne
   Opt-out-Marker); die Fenced-Fläche trägt die Regel voll, das Gate prüft sie
   nicht. Der Wächter dort ist das Review, kein Sensor. Diese Lücke ist
   **benannt**, nicht still — sie gehört in `AGENTS.md` §3.11 und in
   [`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md).

2. **§3.11-Entwurf — Fassung 2 (supersedet `ADR-0072` §Entscheidung
   Punkt 5).** Der Unterabschnitt wird mit der Reichweite aus Punkt 1 neu
   gefasst; der korrigierende Eingriff (`df47282`: Platzhalter-Beispiel,
   Hausform-Beispiel, Form-Beschreibung im Klammertext) ist damit ein
   **beschlossener** Text, keine still tolerierte In-place-Änderung. Der
   Wortlaut trägt:

   - **Aussage.** Kein Markdown-Dokument dieses Repos nennt an irgendeiner
     Stelle — Prosa, Inline-Code oder Fenced-Code-Block — einen host-lokalen
     absoluten Pfad (ein Wurzel-Segment eines Entwicklerrechners,
     Präfixliste `hostpaths.prefixes`, oder ein Windows-Laufwerks-/UNC-Muster);
     ein Schwester-Artefakt wird in der Hausform zitiert.
   - **Formzitat.** Eine verbotene Form wird nur als Platzhalter gezeigt
     (`<Host-Wurzel>`, kein reales Segment); das Richtig-Beispiel ist die
     Hausform. Der Klammertext, der sagt, die Regel decke auch dieses Beispiel,
     ist der richtige.
   - **Was der Sensor deckt — und was nicht.** Die durchsetzbare Hälfte trägt
     das `hostpaths`-Modul in `make docs-check` (`make gates`). Es deckt
     `.md`-Dateien unter `scan.roots` in Prosa und Inline-Code. Es deckt
     **nicht**: Fenced-Code-Blöcke, **relative** Pfade, **Nicht-Markdown**
     (`Makefile`, `tools/**`, `harness/mk/**`) und Dateien unter `scan.ignore`.
     **Die Fenced-Fläche deckt die Regel voll, der Sensor nicht** — das ist
     die benannte Lücke aus Punkt 1.

3. **Lokator-Disposition.** Die Lokator-Klausel in `ADR-0072` §Entscheidung
   Punkt 3 wurde mit Punkt 3 durch
   [`ADR-0074`](0074-zitationsform-schwester-repo-hausform.md) supersedet und
   nicht restituiert (sein Ersatztext fasst nur die Zitationsform; „Der Rest
   bleibt" führt die Lokator-Klausel nicht auf). Die Disposition lautet damit
   **supersedet und nicht restituiert** — es besteht **keine** Pflicht,
   Zeilen-/Bereichs-Lokatoren zu ersetzen. [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
   Punkt 1 nimmt Lokatoren weiter in die Zitat-Klasse auf; das ist eine
   **Erlaubnis**, kein Auftrag, und innerhalb von §Entscheidung bleibt auch sie
   unberührbar. Die noch stehenden Lokatoren (`ADR-0054`, ein
   `done/`-Zeitdokument) sind damit kein Defekt dieses Vorgangs.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; Fences frei, §3.11 unangetastet | kein Folge-ADR, kein Nachzug | die Zusage „alle entfernt, ohne Ausnahme" wäre auf die gescannte Fläche zurückgeschnitten; §3.11 bliebe doppelt lesbar (Regel-Satz gegen sein eigenes Beispiel); die Klassengrenz-Verletzung bliebe still toleriert |
| B — den §3.11-Hunk in `ADR-0072` zurücknehmen | stellt die vorige Fassung her | setzt eine **reale** Form in den Fence zurück: der repo-weite Grep wäre nicht mehr 0 — der Sensor fängt es nicht, der Grep schon; der Auftrag „entfernt" wäre verfehlt |
| C — die Zitat-Klasse auf „Formzitate innerhalb einer Entscheidung" erweitern | die Beispiel-Änderung wäre nachträglich in-class | öffnet den Kanal, den [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) eng hält („das Gerüst darf sich ändern, die Aussage nie" würde unschärfer); und sie heilt den **Klammertext** nicht — er ist keine Form eines Verweises, sondern Erklärung |
| D — §3.11-Entwurf per Folge-ADR neu fassen, Reichweite deklarieren — **gewählt** | die In-place-Änderung wird **beschlossen** statt toleriert; die Regel wird an einer Stelle festgezogen; die Entfernung bleibt erhalten | ein zweiter Träger für §3.11; der Leser muss die Überlagerung (`ADR-0072` Punkt 5 supersedet) auflösen |

**Fazit:** D. B setzt die entfernte Form zurück, C öffnet den Kanal, A lässt den
Auftrag zurückgeschnitten.

## Konsequenzen

- Positiv: Die Regel hat eine Reichweite und **einen** Ort — Fences
  eingeschlossen, an `AGENTS.md` §3.11; der Slice-Plan trägt keine zweite
  Fassung.
- Positiv: Die In-place-Änderung an `ADR-0072` §Entscheidung Punkt 5 ist als
  Folge-ADR nachgezogen (`modul-08-agentenrollen.md` §Konflikt-Pfad, Verdikt
  „Lockerung legitim, aber undokumentiert") — kein stiller Substanz-Edit.
- Positiv: Die Lokator-Klausel hat eine Disposition; kein Lokator-Ersatz ist
  offen.
- Negativ mit Grenze: Die Regel ist **strenger als ihr Sensor** — Fenced-Blöcke
  prüft kein Gate; der Wächter ist das Review. Benannt (`AGENTS.md` §3.11,
  [`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md)).
- Folgepflicht (Implementer-Zug, Doku): `AGENTS.md` §3.11 auf die Fassung 2
  (Reichweite Fences eingeschlossen, benannte Lücke; das Falsch/Richtig-Paar
  bleibt Platzhalter/Hausform); [`harness/sensors/docs-check.md`](../../../harness/sensors/docs-check.md)
  §Grenze um die Aussage „die Regel deckt Fences, das Modul nicht" ergänzen.
  Kein Produkt-Code, keine Config-Änderung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `hostpaths` (Digest aus `d-check.mk`) — aktiviert durch [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) | **0 Befunde** `hostpath-forbidden` in Prosa und Inline-Code; die Fenced-Fläche ist **nicht** maschinell geprüft (die Regel deckt sie voll) — der Wächter dort ist das Review | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** `d-check` ändert die Fence-Behandlung — das Modul
prüft Fenced-Blöcke —, dann schließt der Sensor die Lücke und Regel und Sensor
sind gleich weit; **(b)** ein Dokument braucht eine reale, verbotene Form als
Beispiel, weil der Platzhalter sie nicht trägt — dann braucht die Regel eine
zweite, eng benannte Beispiel-Form statt der Ausweitung; **(c)** die
Reichweite-Deklaration erweist sich an `AGENTS.md` §3.11 als zu weit — dann
wandert sie an einen engeren Ort zurück, nicht in den Plan. Andernfalls
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Review zu `slice-078` F-1 (Klassengrenze in §Entscheidung) und F-3 (Reichweite der Regel undeclared), F-2 (Disposition der Lokator-Klausel); Reichweite Fences eingeschlossen, §3.11-Entwurf Fassung 2, Lokator-Klausel supersedet-nicht-restitiert | der Architect-Verdikt dieses Zugs (Konflikt-Pfad `slice-078`) <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0075` — eine **Zitat-Korrektur**
([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)) ausgenommen
(Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
