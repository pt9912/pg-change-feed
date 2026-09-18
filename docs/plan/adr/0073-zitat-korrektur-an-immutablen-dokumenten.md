# ADR-0073: Zitat-Korrektur an immutablen und historischen Dokumenten — die Klasse (zieht `AGENTS.md` §3.5 nach)

**Status:** Accepted

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Planner-Zug, der
die 31 Befunde gemessen hat, und als die Implementer-/Reviewer-/Verifier-Läufe
der betroffenen Slices — `modul-08-agentenrollen.md` §Rollen-Regeln)

**Bezug:** `AGENTS.md` §3.5 (die nachgezogene Hard Rule), §3.7 (Ist-Zustand),
§5 (Traceability) · [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
(die Aktivierung, die diese Klasse braucht) · [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(Commit-Kennung als Beleg) · Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule
für Accepted-ADRs · `modul-10-review-harness.md` (Review-Reports sind
Lauf-Belege) · `modul-06-roadmap.md`/`modul-05-planning-harness.md`
(Zeitdokumente) · `docs/plan/adr/README.md` (Index-Kopf trägt denselben Satz)
· der Architect-Verdikt dieses Zugs (hostpaths-Aktivierung ohne Ausnahme)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0045`, `ADR-0069`,
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Auftraggeber hat angeordnet, die 31 `hostpath-forbidden`-Befunde **alle**
zu entfernen, ohne Ausnahme — 11 davon sitzen in zwei `Accepted` ADRs
(`ADR-0051`, `ADR-0054`). `AGENTS.md` §3.5 sagt: „Eine ADR mit Status
`Accepted` wird nicht **inhaltlich** überschrieben. Korrekturen entstehen als
neue ADR mit `Supersedes ADR-NNNN`." Bis heute hat dieses Repo auf den
Konflikt mit **Ausnahmen** geantwortet: `matrix.status.exempt-paths` in
`.d-check.yml` existiert nur, weil In-place-Korrekturen an `Accepted` ADRs als
verboten gelten (Grandfather für `ADR-0039`/`ADR-0041`).

„Ohne Ausnahme" ist deshalb nur erreichbar, wenn §3.5 **selbst** für eine eng
umrissene Klasse nachzieht. Die zweite Hälfte der Befunde sitzt in `done/`
(Zeitdokumente, Modul 5/6) und `docs/reviews/**` (Lauf-Belege, Modul 10) —
nicht durch §3.5 geschützt, aber durch die **Konvention** „Lauf-Belege und
Zeitdokumente werden nicht umgeschrieben".

**Der Konflikt ist Umfang, nicht Rang.** §3.5 (Immutabilität) und die neue
Zusage „kein host-lokaler Pfad in der Doku" ([`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md))
stehen nicht in einem Rangverhältnis; sie überschneiden sich in der Frage,
was das Wort **„inhaltlich"** umfasst. Eine Rangordnung zwischen Hard Rules
wäre Seniorität in Tabellenform und ist als Auflösungsmittel ausgeschlossen
(`modul-08-agentenrollen.md` §Konflikt-Pfad, „Senioritäts-Verbot"). Die
Lösung ist eine **Klassen-Eingrenzung**, keine Priorisierung.

## Entscheidung

Wir führen eine **Zitat-Korrektur** als benannte, eng umrissene Klasse ein und
ziehen `AGENTS.md` §3.5 auf sie nach.

1. **Die Klasse — die Zitat-Korrektur.** Eine Änderung ist eine
   *Zitat-Korrektur*, wenn sie **ausschließlich die Form eines Verweises auf
   einen unveränderten Referenten** ändert: host-lokale absolute Pfade,
   Linkziele, Zeilen-/Bereichs-Lokatoren, die Form einer gebrochenen Referenz.
   Sie ist **keine** Zitat-Korrektur, wenn sie einen der folgenden Punkte
   berührt:

   - bei einer ADR: §Entscheidung, §Konsequenzen *der Aussage nach*,
     §Verglichene Alternativen, §Status, die `Supersedes`-Kette, die
     Fitness-Function-Regeln, die Re-Evaluierungs-Trigger, `Datum`/`Autor`,
     die Aussage-Semantik von §Bezug/§Schärft;
   - bei einem Record (`done/`-Slice/Wellen, `docs/reviews/**`): die **Funde,
     Beobachtungen, Closure-Aussagen und DoD-Haken**.

   Kurzform: **das Gerüst darf sich ändern, die Aussage nie; der Referent
   bleibt derselbe.**

2. **`AGENTS.md` §3.5 wird nachgezogen** (Wortlaut der Folgearbeit). Der
   Grundsatz bleibt unangetastet; nur die Reichweite des Wortes „inhaltlich"
   wird benannt — §3.5 lautet künftig: „Eine ADR mit Status `Accepted` wird
   nicht **inhaltlich** überschrieben. Korrekturen entstehen als neue ADR mit
   `Supersedes ADR-NNNN`. Eine **Zitat-Korrektur** (das Zitat- und
   Verweisgerüst — host-lokale Pfade, Linkziele, Zeilen-Lokatoren — bei
   unverändertem Referenten) ist kein inhaltliches Überschreiben und in-place
   zulässig, wenn sie `ADR-0073` genügt; Entscheidung, Konsequenzen,
   Alternativen, Status und `Supersedes`-Kette bleiben unberührbar." Das ist
   eine **Schärfung der Reichweite**, keine Aufhebung der Immutabilität.

3. **Belegform der Korrektur.** (a) Die korrigierende Commit-Message nennt
   `ADR-0073` — die Commit-Kennung ist der Traceability-Beleg (`ADR-0045`).
   (b) Eine betroffene `Accepted` ADR erhält **eine** Zeile ihrer eigenen
   §Geschichte-Tabelle: Datum, „Zitat-Korrektur — host-lokale Pfade ersetzt
   (`ADR-0073`)", Verweis auf den Commit. Die §Geschichte-Zeile ist das
   **in-Format-Provenienzfeld** der ADR und **selbst** eine Zitat-Korrektur
   (Form, nicht Aussage); sie macht den Eingriff **in der ADR sichtbar**,
   statt ihn nur im `git`-Diff zu lassen. Records (`done/`,
   `docs/reviews/**`, lebende Doku) tragen keine §Geschichte und bekommen
   **keine** Zeile — bei ihnen ist der Commit der Beleg.

   Begründung der Zeilen-Pflicht: Eine ADR wird als **Historie ohne `git`
   gelesen** (sie ist die immutabile Entscheidungsschicht); eine stille
   In-place-Änderung wäre dort von einer substanziellen Änderung nicht zu
   unterscheiden. Die Zeile ist der minimale, format-eigene Träger, der diese
   Unterscheidung im Dokument selbst hält. Ein Commit allein genügt **nicht**,
   weil er die Leseordnung (ADR vor `git`) nicht bedient.

4. **Was bei Records erlaubt ist, das bei ADRs nicht gilt (und umgekehrt).**
   Records sind **keine** `Accepted`-Dokumente: ihre Einfrierung ist
   *zeitlich* (sie beschreiben einen vergangenen Lauf), nicht *Status*-basiert.
   Sie brauchen deshalb **keine** §Geschichte-Zeile und **keinen** neuen
   `Accepted`-Träger je Datei; der eine Commit genügt. Umgekehrt gilt für sie
   dieselbe Grenze wie für ADRs: ihr **Record-Inhalt** (Funde, Beobachtungen,
   Closure-Aussagen) ist unantastbar — **nur das Zitat-Gerüst** darf sich
   ändern. Der Träger beider Hälften ist **dieselbe** Entscheidung (diese ADR);
   die Klasse ist eine, die Dokumentarten sind zwei.

5. **Die `Supersedes`-Alternative wird verworfen** (§Verglichene Alternativen C).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; §3.5 unangetastet, die zwei `Accepted` ADRs behalten ihre Pfade | keine Regel-Änderung | der Auftrag „alle entfernt, ohne Ausnahme" bleibt unerfüllt, oder es bräuchte doch eine Ausnahme — genau die zweite Quelle, die `ADR-0072` ausschließt |
| B — Ausnahmeliste im Scope für die zwei `Accepted` ADRs (`matrix`/`hostpaths`-`exempt-paths`) | kein Eingriff in `Accepted` Text | **vom Auftraggeber verworfen**; eine Ausnahmeliste altert, ist eine zweite Quelle für „wo gilt die Zusage nicht", und der Anlass ist mit einer engen Klasse sauber lösbar |
| C — Folge-ADR mit `Supersedes`, die die betroffenen Klauseln host-pfad-frei **neu fasst** | kein In-place-Eingriff; das Werkzeug, das dieses Repo für *unwahre* `Accepted`-Klauseln nutzt ([`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md)/[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) | **entfernt die Pfade nicht**: der alte, host-lokale Text bleibt in der abgelösten ADR stehen und wird weiter gescannt — der Auftrag bleibt unerfüllt (im Umfang „Entfernen" ist C gleich A); zudem ist **nichts Normatives falsch**, also ist `Supersedes` das falsche Werkzeug |
| D — Zitat-Korrektur als Klasse einführen, §3.5 klassenweise nachziehen — **gewählt** | „ohne Ausnahme" wird erreichbar; die Immutabilität bleibt für jeden Inhalt voll erhalten; **eine** Klasse für alle drei Dokumentarten (Accepted ADR, Record, lebende Doku) statt 31 Einzelfälle | es entsteht ein **schmaler Kanal** in `Accepted` ADRs (Konsequenzen); der Wächter ist ein Urteil, kein Sensor |

**Fazit:** D. C erreicht das Ziel nicht; B verfehlt den Auftrag; A lässt ihn
unerfüllt.

## Konsequenzen

- Positiv: „ohne Ausnahme" ist erreichbar, **ohne** eine Ausnahmeliste — die
  Immutabilität bleibt im Grundsatz (jede Aussage) unangetastet.
- Positiv: **eine** Klasse für alle drei Dokumentarten; kein Sonderfall je
  Datei, kein 31-facher Einzelfall.
- Positiv: Der Korrektur-Beleg steht **im Dokument** (die §Geschichte-Zeile)
  und **im `git`** (die Commit-Kennung) — die ADR-Historie bleibt ohne `git`
  lesbar und unterscheidbar.
- Negativ mit Grenze: Es entsteht ein **schmaler Kanal** in `Accepted` ADRs.
  Der Wächter ist ein **Urteil** („ändert die Korrektur nur die Zitat-Form?"),
  kein Sensor: Das `hostpaths`-Modul prüft nur das Ergebnis (0 Befunde), nicht
  die Klassenzugehörigkeit. Ein Review, das den Diff nicht liest, kann den
  Kanal für Substanz öffnen. Das ist benannt, nicht wegdefiniert — der
  Re-Evaluierungs-Trigger (a) greift genau hier.
- Negativ: Die §Geschichte-Zeile ist selbst ein Eingriff in die `Accepted`
  ADR; sie ist bewusst als Zitat-Korrektur definiert (Provenienzfeld, nicht
  Aussage), aber sie ist die einzige Stelle, an der diese Entscheidung **sich
  selbst** anwendet — der Grund, sie eng zu fassen.
- Folgepflicht (Implementer-Zug): `AGENTS.md` §3.5 Wortlaut (Punkt 2);
  `AGENTS.md` §3.11 (Wortlaut aus [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md));
  `docs/plan/adr/README.md` Kopf-Satz („Korrekturen als Folge-ADR mit
  `Supersedes`" um „eine Zitat-Korrektur ausgenommen" ergänzen); die 31
  Korrekturen aus [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md);
  je betroffener ADR eine §Geschichte-Zeile.
- Hinweis (die gestellte Frage): **Keine Priorität/Rangordnung nötig** — der
  Konflikt wird durch die klassenweise Eingrenzung von §3.5 gelöst. `ADR-0051`
  und `ADR-0054` bleiben in **jeder** ihrer Entscheidungen unberührt; geändert
  wird an ihnen ausschließlich die Zitat-Form.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | Die **Klasse selbst ist ein Urteil**, nicht maschinell prüfbar: kein Sensor entscheidet, ob eine Korrektur nur das Zitat-Gerüst ändert. Maschinell prüfbar ist nur das **Ergebnis** dieser Entscheidung | — |
| `d-check` Modul `hostpaths` (Digest aus `d-check.mk`) — aktiviert durch [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) | nach den Korrekturen **0 Befunde** `hostpath-forbidden`; der Kanal ist am Ergebnis **nicht** erkennbar, sondern an seinem Diff — der Wächter ist das Review | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** eine als Zitat-Korrektur deklarierte Änderung
erweist sich als inhaltlich (der Kanal hat Substanz durchgelassen) — dann
sofort Folge-ADR: Klasse verschärfen oder schließen; **(b)** `d-check` liefert
einen Marker, mit dem der Befund zeilenweise ausnehmbar wird — dann kann die
Zitat-Korrektur durch den Marker ersetzt werden (Klasse entlastet);
**(c)** `matrix.status.exempt-paths` (Grandfather `ADR-0039`/`ADR-0041`) kann
enger werden, wenn eine Folge-ADR die dort gezeigten Zeiger auflöst; **(d)**
`ADR-0051`/`0054` werden supersedet und ihr host-lokaler Text fällt weg —
dann ist der Anlass der Klasse für diese zwei entfallen, die Klasse selbst
bleibt. Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Auftrag „alle host-lokalen Pfade entfernen, ohne Ausnahme"; 11 Befunde in zwei `Accepted` ADRs; Zitat-Korrektur als Klasse, §3.5 klassenweise nachgezogen | der Architect-Verdikt dieses Zugs (hostpaths-Aktivierung ohne Ausnahme) |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0073` — eine **Zitat-Korrektur** ausgenommen (diese ADR selbst)
(Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
