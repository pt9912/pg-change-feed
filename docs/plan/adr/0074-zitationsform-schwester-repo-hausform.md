# ADR-0074: Zitationsform eines Schwester-Repos — die Hausform — Supersedes ADR-0072 (eine Klausel)

**Status:** Accepted — Supersedes [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
in genau **einer** Klausel: §Entscheidung Punkt 3 (die Zitationsform). Alles
Übrige aus `ADR-0072` — die Aktivierung ohne Ausschlussliste (Punkt 1), kein
Ventil (Punkt 2), die Korrektur **aller** 31 Stellen (Punkt 4), der Entwurf
der Hard Rule §3.11 (Punkt 5), der Config-Block, die Grenzen, die Fitness
Function und die Re-Evaluierungs-Trigger — bleibt unverändert bestehen und
wird hier nicht wiederholt.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Planner-Zug, der
die 31 Befunde gemessen hat — die Hausform wurde am **Bestand** gemessen,
`modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz)

**Bezug:** [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) (in
einer Klausel korrigiert), [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(die Zitat-Korrektur-Klasse), `AGENTS.md` §3.5 (Accepted-ADRs immutable —
Folge-ADR), §3.7 (Ist-Zustand) ·
`docs/plan/planning/done/welle-19-results.md` (Belegstelle der Hausform) ·
`docs/plan/planning/observations/BEO-PGC/plan-vorlagen-defekt/state.md`
(Belegstelle) · `docs/user/benutzerhandbuch.md` (Link-Variante mit
`d-migrate`) · `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md`
(Verdikt dieses Zugs)

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie
[`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) (Accepted,
2026-09-15) setzt in §Entscheidung Punkt 3 die Zitationsform eines
Schwester-Repos als **besitzer-qualifiziert** — `pt9912/<repo>`, optional mit
dem Pfad darin — und belegt sie mit **einem** Fund (`docs/user/benutzerhandbuch.md`
zitiert `d-migrate` als GitHub-Link). Die Korrektur der 31 Stellen soll nach
**einer** Form laufen, damit sie mechanisch wird statt 31 Einzelfälle.

**Der Bestand widerspricht der gesetzten Form.** Eine Messung über die
`*.md`-Dateien dieses Repos (ohne die vendored Baseline `.harness/**`) zeigt:
das Schwester-Repo wird **überwiegend als blankes Inline-Code-Wort** genannt —
`` `d-check` `` erscheint **89-mal** —, und sein Artefakt wird **relativ
darin** benannt, ohne Host- und ohne Owner-Präfix. Belegstellen:

| Belegstelle | Form |
|---|---|
| `docs/plan/planning/done/welle-19-results.md` | `` `d-check`s `structure`-Regel `` · `` `d-check`s `commits`-Modul `` · `` `d-check` 574 Dateien / 0 Befunde `` |
| `docs/plan/planning/observations/BEO-PGC/plan-vorlagen-defekt/state.md` | `` (kein d-check-Modul erkennt Zeilen-Duplikate) `` |
| `docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md` | `` `d-check`-Modul `commits` `` |
| `docs/plan/adr/0070-supersede-reichweite-und-klassengrenze.md` | `` `d-check` Modul `commits` `` |
| `docs/user/benutzerhandbuch.md` | Link-Variante: `[`d-migrate`](https://github.com/pt9912/d-migrate)` |

Die von `ADR-0072` gewählte Form ist damit **nicht falsch** (sie ist
host-pfad-frei und auflösbar), aber sie ist **nicht die Hausform**: sie führte
für 31 Stellen eine **zweite** Konvention neben 89 bestehenden Nennungen ein.
Die Regel „eine Form für alle 31" ist richtig; die *gewählte* Form war es
nicht. Der Bestand liefert die richtige — die Form ist **belegbar**, nicht zu
setzen. Das ist der Anlass zu dieser Folge-ADR (`AGENTS.md` §3.5).

**Und die Hausform selbst geprüft.** Kein Vorkommen der Hausform steht *neben*
einem host-lokalen Pfad, das nicht schon unter den 31 gemessen wäre — die
Zählung bleibt **31**; es kommt keine 32. Stelle hinzu.

## Entscheidung

Wir wählen die **Hausform** — und superseden damit `ADR-0072` in genau der
Zitations-Klausel.

1. **Kanonische Form — das blanke Repo-Wort.** Ein Schwester-Repo wird als
   **blankes Inline-Code-Wort** genannt: `` `d-check` ``, `` `d-migrate` ``,
   `` `ai-harness-init` ``. Das Artefakt darin wird **relativ** angehängt, ohne
   Host- und ohne Owner-Präfix: `` `d-check`s `tools/coverage-gate.sh` ``,
   `` im Repo `d-check`: `Dockerfile` ``. Kein führendes `/`, kein
   Host-Segment.

2. **Link-Variante.** Wo ein anklickbarer Verweis gewünscht ist, zeigt er auf
   `https://github.com/pt9912/<repo>` (bzw. `/blob/main/<pfad>`) — die in
   `docs/user/benutzerhandbuch.md` mit `d-migrate` geübte Form. Der
   Owner-Präfix bleibt damit **verfügbar**, aber **nicht verlangt**; er ist
   die Ausnahme (Link), nicht die Regel (Prosa).

3. **Der Rest aus [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
   bleibt.** Aktivierung ohne Ausschlussliste, kein Ventil, Korrektur **aller**
   31 Stellen, Hard-Rule-§3.11-Entwurf, Config-Block, Grenzen, Fitness
   Function und Re-Evaluierungs-Trigger gelten unverändert.

**Beispiel (Fence, weil die Regel sonst ihre eigene Aussage verletzte):**

```text
vorher (ADR-0072-Form):  Real geprüftes Vorbild: `pt9912/d-check/Dockerfile` (Stage `coverage`)
nachher (Hausform):      Real geprüftes Vorbild: `d-check`s `Dockerfile` (Stage `coverage`)
```

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; die `ADR-0072`-Form (`pt9912/<repo>`) bleibt | kein Folge-ADR, kein Doku-Nachzug | führt für 31 Stellen eine **zweite** Konvention neben 89 bestehenden Nennungen ein; die Form wäre gesetzt statt belegbar — genau der Vorwurf, den `ADR-0070` an eine „gesetzte" Aussage richtet |
| B — `pt9912/<repo>/<pfad>` beibehalten, aber als *Empfehlung* neben der Hausform | beide Formen „erlaubt" | zwei Formen für **eine** Sache — die Doppelquelle, die `ADR-0072` gerade vermeiden wollte; der Implementer müsste 31-mal wählen |
| C — **Hausform: blankes Repo-Wort + relativer Pfad (gewählt)** | belegbar aus dem Bestand (89 Nennungen); **eine** Form, keine 31 Einzelfälle; konsistent mit den bereits vorhandenen Zitaten | der Repo-Name ist ohne Owner nicht global eindeutig — in diesem Repo unkritisch (`d-check`/`d-migrate`/`ai-harness-init` kollidieren nicht); die Link-Variante trägt den Owner dort, wo Auflösbarkeit zählt |
| D — überall die volle GitHub-URL | eindeutig, auflösbar | in Fließtext-Provenienz („Vorbild: …") übermäßig lang, verdrängt den Satz; entspricht **keiner** der beiden geübten Formen |

**Fazit:** C.

## Konsequenzen

- Positiv: Die Zitationsform ist **belegbar aus dem Bestand** statt gesetzt —
  der Implementer folgt der Hausform, die 89-mal bereits steht.
- Positiv: **eine** Form für alle 31 Stellen; die Korrektur wird mechanisch.
- Positiv: `ADR-0072`s übrige Klauseln (Aktivierung, kein Ventil, §3.11,
  Grenzen) bleiben in Kraft — kein Nachzug nötig.
- Negativ mit Grenze: Das blanke Repo-Wort ist ohne Owner nicht **global**
  eindeutig; wer das Repo nicht kennt, kann `d-check` nicht auflösen, ohne den
  Kontext zu lesen. Die Link-Variante ist die Antwort dort, wo Auflösbarkeit
  zählt — die Regel ist also „Hausform in Prosa, Link mit Owner wo ein Link
  gewünscht ist", keine owner-freie Beliebigkeit.
- Negativ: Die 89 bestehenden Nennungen sind **nicht** nachzuziehen (sie sind
  bereits in der Hausform); nur die 31 Stellen werden korrigiert.
- Folgepflicht (Implementer-Zug, Doku): die 31 Korrekturen in der Hausform aus
  §Entscheidung Punkt 1 statt in der `ADR-0072`-Form; im Übrigen gilt die
  Folgepflicht-Liste aus [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
  unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `hostpaths` (Digest aus `d-check.mk`) — aktiviert durch [`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) | **0 Befunde** `hostpath-forbidden`; die Hausform (`d-check`, `d-migrate`, `ai-harness-init` als blanke Wörter) ist host-pfad-frei und wird nicht gemeldet — der Bestand ist der Beleg, nicht der Sensor | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** ein Schwester-Repo-Name wird mehrdeutig (zwei
Repos unter demselben blanken Namen werden zitiert) — dann braucht die Form
den Owner-Präfix als Pflicht, nicht als Link-Variante; **(b)** der Bestand
wechselt die Hausform (die Mehrheit der Nennungen wird besitzer-qualifiziert)
— dann neu messen und diese Entscheidung nachziehen. Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Messung am Bestand (89 blanke `d-check`-Nennungen, Artefakte relativ benannt); korrigiert `ADR-0072` §Entscheidung Punkt 3 (Zitationsform) | `docs/reviews/architect-verdict-hostpaths-aktivierung-ohne-ausnahme.md`, eigene Bestands-Messung |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0074` — eine **Zitat-Korrektur**
([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)) ausgenommen
(Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
