# Architect-Verdikt — braucht die Bündel-Aufnahme von `tracked` eine eigene ADR?

**Datum:** 2026-09-17 · **Stand:** `550a959` · **Rolle:** Architect (eigener
Kontext; die Fakten dieses Zugs sind selbst nachgeprüft, nicht aus dem Review
übernommen — außer wo unten explizit als „übernommen" gekennzeichnet) ·
**Anlass:** `docs/reviews/review-slice-d-check-tracked-modul.md` F-2 (HIGH) und
F-3 (MEDIUM), Slice `slice-d-check-tracked-modul`, Welle `welle-d-check`.

## 1. Verdikt: **keine eigene ADR** — tragender Präzedenzfall ist `structure` (`f9e5a3c`), nicht `hostpaths` (`ADR-0072`/`ADR-0075`)

Die Aufnahme von `tracked` in die `modules:`-Liste von `.d-check.yml` braucht
**keine eigene ADR**. Das ist dieselbe Schlussfolgerung, die
`harness/sensors/docs-check.md` §Bindung bereits zieht — aber der dort
zitierte Beleg trägt sie nicht (Review F-2, bestätigt in §2 unten). Der
tragende Beleg ist ein anderer: der Commit `f9e5a3c`
(`docs(harness): structure-Modul aktiviert`, 2026-09-09), der `structure`
genau auf demselben Weg — einfacher Commit, keine ADR — in dieselbe
`modules:`-Liste aufgenommen hat, in der `tracked` jetzt steht.

## 2. Eigene Prüfung: F-2 des Reviews ist korrekt — `ADR-0072`/`ADR-0075` tragen die **gegenteilige** Aussage

Gelesen: `docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md`
§Entscheidung, `docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md`
§Kontext/§Entscheidung.

`ADR-0072` §Entscheidung Punkt 1 sagt wörtlich: *„Dieses Repo liest die
Verschärfung dennoch als **Entscheidung mit Träger** — sie prägt, welche
Zitate künftig zulässig sind und welchen Wortlaut Hard Rule §3.11 trägt, und
sie ist **ohne die Zitat-Korrektur der `Accepted`-ADRs nicht durchführbar**.
Der Träger ist diese ADR; das Dokument ist damit eine ADR, keine
Aktennotiz."* Drei Dinge zwingen dort zur ADR, die bei `tracked` alle drei
fehlen:

| Zwang bei `hostpaths` | Beleg in `ADR-0072`/`ADR-0075` | Bei `tracked` gegeben? |
|---|---|---|
| Reale Befunde, die eine Korrektur erzwingen | 31 Befunde `hostpath-forbidden` (§Kontext, Tabelle) | **Nein** — 0 Befunde, eigene Messung §3 |
| Korrektur an zwei `Accepted`-ADRs (§3.5-Ausnahme) | `ADR-0072` §Entscheidung Punkt 1/4, Zitat-Korrektur über `ADR-0073` | **Nein** — keine ADR wird berührt |
| Neue Hard Rule | `AGENTS.md` §3.11, Entwurf in `ADR-0072` §Entscheidung Punkt 5, neu gefasst durch `ADR-0075` | **Nein** — kein Vorschlag für eine neue Hard Rule |

`ADR-0075` entsteht zusätzlich ausschließlich aus einer **Nachfrage-Kette**
zu `ADR-0072` (Reichweite, Klassengrenze, Lokator-Disposition, siehe
`ADR-0075` §Kontext) — sie zeigt, dass selbst eine bereits geschriebene
ADR bei `hostpaths` noch Nacharbeit brauchte. Als Präzedenzfall für „braucht
keine ADR" tragen beide Dokumente ihren Satz nicht: Sie sind Belege für den
**gegenteiligen** Fall. Der Reviewer hat F-2 korrekt hergeleitet; ich
bestätige das eigenständig anhand desselben Wortlauts, nicht durch Übernahme
seiner Zusammenfassung.

## 3. Eigene Prüfung: Commit `f9e5a3c` — Charakterisierung des Reviews trifft zu

Gelesen: `git show f9e5a3c` (voller Diff und Commit-Message, nicht nur
`--stat`).

Bestätigt, alle drei Merkmale:

- **Kein ADR-Bezug.** Die Commit-Message (`docs(harness): structure-Modul
  aktiviert — Register-Tabellen unter Spalten-Invarianten`) nennt keine
  ADR-Nummer, referenziert keine. Der Diff selbst ist ausschließlich
  `.d-check.yml` (68 Zeilen, ein `structure:`-Konfigurationsblock plus die
  Aufnahme von `structure` in die `modules:`-Liste).
- **Keine neue Hard Rule.** Kein `AGENTS.md`-Bezug in der Message; keine
  begleitende Änderung an `AGENTS.md` im selben Commit oder — geprüft über
  den Diff-Umfang — an irgendeiner anderen Datei außer `.d-check.yml`.
- **Keine Korrektur an `Accepted`-Inhalten.** Die Message beschreibt einen
  eigenen Findungs-Zyklus („Erster Lauf meldete 14 überlange ID-Zellen …
  Grenzen an die beabsichtigte Form angepasst statt Inhalte zu kürzen") —
  eine Konfigurationsanpassung (Zellgrenzen im neuen `structure:`-Block),
  keine inhaltliche Korrektur an einer bestehenden `Accepted`-ADR oder einem
  anderen unveränderlichen Dokument.

Die Charakterisierung des Reviewers (F-4) ist damit durch eigene Lektüre
bestätigt, nicht nur übernommen.

## 4. Eigene Messung: reale Größenordnung der `tracked`-Aktivierung in diesem Repo

Selbst gemessen (nicht übernommen), Stand `550a959`:

```
make doc-tracked
d-check: 876 Datei(en) geprüft, 0 Befund(e)
```

Das Review nennt 875 Dateien zum Zeitpunkt seines eigenen Laufs (2026-09-17,
vor diesem Zug); die Differenz von einer Datei ist durch den seither
hinzugekommenen Review-Report selbst erklärt (`docs/reviews/review-slice-d-check-tracked-modul.md`
ist eine neue, getrackte Markdown-Datei mit Links). Die tragende Aussage —
**0 Befunde** — ist in beiden Läufen identisch und damit bestätigt, nicht nur
übernommen.

Diese Zahl ist die Voraussetzung für den Vergleich mit `structure`: `structure`
wurde ohne reale Befunde aktiviert und musste seinen Fehlschlag (14 zu lange
Zellen) durch eine Grenzen-Anpassung *vor* dem grünen Merge beheben — dieselbe
Form wie hier: `tracked` aktiviert grün, ohne eine einzige Korrektur am
bestehenden Doku-Bestand.

## 5. Warum `structure` der tragende Präzedenzfall ist, nicht `hostpaths`

Die vergleichende Eigenschaft ist nicht „ein Modul wird ins Bündel
aufgenommen" — das trifft auf beide Fälle zu und wäre zu grob, um zwischen
ihnen zu unterscheiden. Die tragende Eigenschaft ist: **löst die Aktivierung
eine Folgepflicht am bestehenden Doku-/Regel-Bestand aus, oder bleibt sie
in sich geschlossen?**

| Merkmal | `hostpaths` (`ADR-0072`/`075`) | `structure` (`f9e5a3c`) | `tracked` (dieser Slice) |
|---|---|---|---|
| Reale Befunde am Aktivierungs-Stand | 31 | 14 (Zellgrenzen, kein Inhalt) | 0 |
| Korrektur an bestehendem Doku-Inhalt nötig | ja (31 Zitate) | nein (nur Grenzen-Config) | nein |
| Korrektur an `Accepted`-ADR nötig | ja (2, über `ADR-0073`) | nein | nein |
| Neue `AGENTS.md`-Hard-Rule | ja (§3.11) | nein | nein |
| Träger | eigene ADR (`ADR-0072`, ergänzt `ADR-0075`) | einfacher Commit | — (diese Frage) |

`tracked` deckt sich in allen vier unterscheidenden Merkmalen mit `structure`
und in keinem mit `hostpaths`. Die Analogie zu `hostpaths`, die der
Implementer in `harness/sensors/docs-check.md` §Bindung gezogen hat, war die
falsche Analogie-Achse („beide sind Bündel-Aufnahmen") statt der richtigen
(„löst die Aktivierung Folgepflichten aus, die eine ADR bräuchte, oder
nicht"). Auf der richtigen Achse liegt `tracked` bei `structure`.

## 6. Verdikt-Typ (Modul 8, Konflikt-Pfad)

Dies ist keiner der drei Konflikt-Pfad-Fälle im engeren Sinn (keine
`Accepted`-Entscheidung steht gegen eine Plan-Behauptung) — der Slice-Plan
selbst hat die Frage bewusst offengelassen (§4 Trigger: „der Architect-Zug
… bestätigt vor Implementierungsbeginn, ob eine eigene ADR nötig ist").
Nächstliegend ist die **Rückkante Review → Plan-Defekt** (Modul 8): der
Implementer hat eine Schlussfolgerung getroffen, die der Architect-Rolle
vorbehalten war, und dabei den falschen Beleg zitiert. Das Verdikt korrigiert
den Beleg, nicht das Ergebnis — die Schlussfolgerung „keine eigene ADR" bleibt
richtig, ihre Begründung ändert sich vollständig.

## 7. F-3: dieses Artefakt ist der nachgetragene Architect-Zug

Der Slice-Plan (`docs/plan/planning/in-progress/slice-d-check-tracked-modul.md`
§4 Trigger, Abschnitt „Start") macht den Architect-Zug explizit zur
Start-Bedingung für `next → in-progress`. Der tatsächliche Übergang
(Commit `262bcda`, Message „WIP-Limit frei. Erster Slice der Welle
welle-d-check.") trug kein Architect-Verdikt und keinen Verweis auf eines —
das Review hat das korrekt als F-3 (MEDIUM, „Trigger ohne Übergabe-Artefakt
übersprungen") benannt.

**Dieses Dokument ist der fehlende Zug, nachträglich vollzogen.** Es schließt
sowohl §6 Risiko 1 des Slice-Plans („Die Aufnahme in `modules:` verlangt eine
eigene ADR" — **Ausgang: entfallen**, Begründung siehe §5/§6 dieses
Verdikts) als auch F-3 des Reviews (Übergabe-Artefakt liegt jetzt vor, wenn
auch nach statt vor Implementierungsbeginn — die Sequenz-Verletzung selbst
bleibt im Slice-Plan als Record stehen, sie wird durch diesen Nachtrag nicht
rückwirkend geheilt).

## 8. Folgearbeit — nicht Teil dieses Zugs

Dieses Verdikt ändert keinen Code und keine Konfiguration. Die nachfolgende
Implementer-Fixrunde zieht auf Basis dieses Artefakts nach:

- `harness/sensors/docs-check.md` §Bindung: der `tracked`-Eintrag zitiert
  künftig `f9e5a3c`/den `structure`-Präzedenzfall statt `ADR-0072`/`ADR-0075`
  (F-2-Fix).
- `harness/README.md` (Zeile 114, sieben statt acht Module),
  `.claude/agents/verifier.md:40`, `.claude/agents/implementer.md:44`
  (F-1-Fix, Träger-Nachzug §3.13).
- Slice-Plan §6 Risiko 1 bekommt den Ausgang „entfallen" mit Verweis auf
  dieses Dokument.
- DoD-Checkbox „Review durchgeführt" bleibt bis zur Fixrunde offen
  (Reviewer-Skill §DoD-Checkbox-Nachzug) — dieser Zug hebt sie nicht auf.

## 9. Gates

`make gates` am Stand nach dem Commit dieses Verdikts: Exit-Code direkt und
ungepiped geprüft (`AGENTS.md` §3.9) — siehe Commit-Historie für den
Lauf-Beleg dieses Zugs.

## 10. Was nicht getan wurde

Keine Änderung an `.d-check.yml`, `harness/sensors/docs-check.md`,
`harness/README.md`, `.claude/agents/*.md` oder dem Slice-Plan selbst — das
ist Aufgabe der Implementer-Fixrunde (§8). Keine neue ADR geschrieben (das
ist genau das Verdikt: keine wird gebraucht). Kein Produktionscode berührt.
