# ADR-0072: `hostpaths`-Modul aktiviert — kein host-lokaler Pfad in der Doku, ohne Ausnahme

**Status:** Accepted

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; der Planner-Zug hat die 31 Befunde gemessen,
dieser Zug hat sie **nachgemessen** und die Entscheidung geschrieben —
`modul-08-agentenrollen.md` §Rollen-Regeln: „ADR-Änderung: Architect
schreibt; Reviewer prüft auf Konsistenz; Implementer liest als Constraint")

**Bezug:** `AGENTS.md` §3.1 (Docker-only), §3.5 (Accepted-ADRs immutable —
nachgezogen durch [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)),
§3.6 (Gates ohne ADR nicht lockern), §3.7 (Ist-Zustand), §3.11 (neu — der
Entwurf dieser ADR) · [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(der Träger, der die Korrektur der zwei `Accepted`-ADRs legitimiert) ·
`harness/sensors/docs-check.md` (Sensor-Vertrag und seine Grenzen) ·
`.d-check.yml` (`modules`) · `d-check.mk` (`DCHECK_DIGEST`) · gepinntes
Modul-Image `pt9912/d-check` `v0.75.0`
(`internal/hexagon/core/rules/hostpaths.go` — Modul-Semantik) ·
der Architect-Verdikt dieses Zugs (hostpaths-Aktivierung ohne Ausnahme)

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md) und
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Auftraggeber hat angeordnet, das `hostpaths`-Modul in `.d-check.yml` zu
aktivieren, und den Auftrag in einem zweiten Schritt verschärft: **ohne jede
Ausnahme** — 31 Befunde → 0, keine Dokumentklasse wird ausgenommen. Das ist
eine Gate-**Verschärfung**; `AGENTS.md` §3.6 behandelt im Wortlaut die
*Lockerung* („Jede Schwellen-Senkung … ist ein ADR, kein PR-Kommentar"),
nicht die Verschärfung. Ob auch die Verschärfung einen Träger braucht,
entscheidet die Lesart in §Entscheidung Punkt 1.

### Das Modul

`hostpaths` meldet host-lokale **absolute** Pfade in **Prosa und
Inline-Code** — das Maschinen-Layout eines Entwicklerrechners, das in
zitierten Beispielen stehen bleibt. Fenced-Code-Blöcke sind ausgenommen (dort
gehören bewusste Beispiel-Pfade hin), und es gibt **keinen Opt-out-Marker**:
die einzigen konfigurierbaren Achsen sind die Präfixliste (`prefixes`,
Default `Development`, `home`, `Users`, `Volumes`, `mnt`, `media`) und der
generische modul-lokale `scope`. Windows-Laufwerks- und UNC-Muster sind
**fest** (nicht konfigurierbar). Das Modul liest ausschließlich
`.md`-Dateien (`walkMarkdown` filtert auf `.md`); die Skriptkommentare in
`Makefile`, `tools/**`, `harness/mk/**` erreicht es nicht.

### Der gemessene Befund

Eigene Messung dieses Zugs am gepinnten Digest (`pt9912/d-check` `v0.75.0`,
`--enable hostpaths --disable <alle übrigen Module>`) über
`scan.roots: ["."]` (616 Dateien): **31 Befunde** `hostpath-forbidden`. Die
Verteilung, unabhängig reproduziert:

| Dokumentklasse | Dateien | Befunde | Was an ihr eine Entscheidung verlangt |
|---|---|---|---|
| `Accepted` ADRs | `docs/plan/adr/0051-…md`, `docs/plan/adr/0054-…md` | 11 | `AGENTS.md` §3.5 verbietet die In-place-Korrektur — bis [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) |
| Zeitdokumente in `done/` | `slice-049-…md`, `slice-050-…md`, `welle-14.md`, `welle-14-results.md` | 11 | Closure-Record — die Konvention friert den Beleg, nicht sein Zitat <!-- d-check:status-provenance --> |
| Lauf-Belege | Reviews zu `slice-036`, `-039`, `-049`, `-073` | 7 | Review-Report ist Lauf-Beleg (Modul 10) — dito <!-- d-check:status-provenance --> |
| lebende Doku | `harness/sensors/coverage-gate.md` | 2 | frei korrigierbar |

Die 31 Stellen nennen drei Schwester-Repos über ihren Pfad auf einem
Entwicklerrechner: `d-check`, `d-migrate`, `ai-harness-init` — je Fundstelle
ein Artefakt darin (`Dockerfile`, `tools/coverage-gate.sh`, `Makefile`,
`.github/`, `.golangci.yml`). Der Referent ist unverändert; die
Datei:Zeile-Position zeigt der aktivierte Modul-Lauf, nicht dieses Dokument.

### Der Zwang

„Ohne Ausnahme" ist nur erreichbar, wenn (a) `AGENTS.md` §3.5 für eine eng
umrissene Klasse nachzieht — die 11 in den zwei `Accepted` ADRs — und (b) auch
die Zeitdokumente und Lauf-Belege ihr Zitat korrigieren dürfen. Beides ist
eine **Regel**-Entscheidung, keine Config-Änderung: sie bestimmt, welche
Zitate künftig zulässig sind und worauf die neue Hard Rule §3.11 sich stützt.

## Entscheidung

Wir aktivieren `hostpaths` **ohne Ausschlussliste** und korrigieren alle 31
Stellen auf **eine** host-pfad-freie Zitationsform.

1. **Aktivierung — eine Verschärfung mit Träger.** `hostpaths` kommt in die
   `modules:`-Liste der `.d-check.yml`. Die Aktivierung **fügt eine Prüfung
   hinzu**, senkt keine Schwelle und nimmt keine Ausnahme: §3.6 ist auf sie
   nicht anwendbar (er regelt die Lockerung). Dieses Repo liest die
   Verschärfung dennoch als **Entscheidung mit Träger** — sie prägt, welche
   Zitate künftig zulässig sind und welchen Wortlaut Hard Rule §3.11 trägt,
   und sie ist ohne die Zitat-Korrektur der `Accepted`-ADRs nicht
   durchführbar. Der Träger ist diese ADR; das Dokument ist damit eine ADR,
   keine Aktennotiz.

2. **Keine Ausnahme, kein Ventil.** Kein `hostpaths.scope`, kein
   `scope.ignore`, kein `exempt-paths`, kein Zeilen-Marker — `hostpaths`
   kennt **keine feine Ventil-Achse**: `exempt-paths` fehlt dem Modul, und
   den Zeilen-Marker `d-check:ignore` kennen nur `codepaths`, `ids`,
   `versions` und `diagrams`. Der Geltungsbereich ist der globale Scan
   (`scan.roots: ["."]` minus `scan.ignore`), also die gesamte gescannte
   Markdown-Fläche. Der Config-Block ist deshalb **eine Zeile** (§Config).

   **Benannter Ausgang für einen nicht korrigierbaren Befund.** Wäre eine der
   31 Stellen *nicht* als Zitat-Korrektur behandelbar, hätte sie im Modul
   **kein** Ventil; ihre Behandlung wäre eine eigene, **benannte Folge-ADR**
   — nie eine Scope-Ausnahme, die der Auftrag ausschließt. Nach der Messung
   dieses Zugs trifft das auf **keine** der 31 Stellen zu: alle 31 sind
   Zitat-Korrekturen im Sinne von
   [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md).

3. **Zitationsform — eine für alle 31.** Ein Schwester-Artefakt wird
   **besitzer-qualifiziert** zitiert: `pt9912/<repo>` (der `owner/repo`-Slug),
   optional gefolgt vom Pfad **innerhalb** dieses Repos —
   `pt9912/d-check/Dockerfile`. Wo ein anklickbarer Verweis gewünscht ist,
   zeigt die Kennung auf `https://github.com/pt9912/<repo>` (bzw.
   `/blob/main/<pfad>`). Kein führendes `/`, kein Host-Segment. Zeilen-/
   Bereichs-Lokatoren (`Zeilen 69–93`) werden durch den **stabilen benannten
   Anker** ersetzt (`Stage coverage`, `bench:`-Target) — der Referent bleibt
   derselbe. Die Form ist im Repo bereits geübt: `docs/user/benutzerhandbuch.md`
   zitiert `pt9912/d-migrate` als GitHub-Link.

4. **Alle 31 werden korrigiert, nicht ausgenommen** — die 2 in der lebenden
   Doku, die 7 in `docs/reviews/**`, die 11 in `done/` und die 11 in
   `docs/plan/adr/0051-…`/`0054-…`. Für die zwei `Accepted` ADRs trägt die
   Korrektur [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md);
   für Zeitdokumente und Lauf-Belege trägt sie dieselbe Klasse.

5. **Hard Rule §3.11** — der Wortlaut unten ist Teil dieser Entscheidung, nicht
   ein Nachtrag: die Verkörperung (Modul 6 Schritt 3b) braucht wie die
   Entscheidung die Architect-Rolle. Die Änderung an `AGENTS.md` selbst ist
   Folgearbeit (Planner/Implementer), damit der Wortlaut geprüft und nicht
   nachgereicht wird.

### Config-Block — der ganze Eingriff

`.d-check.yml`, ein Token in der bestehenden Liste (kein `hostpaths:`-Abschnitt:
es gibt nichts zu konfigurieren — keine Präfix-Änderung, kein Scope, keine
Ausnahme):

```yaml
modules: [links, anchors, ids, matrix, versions, structure, hostpaths]
```

Kein `hostpaths:`-Knoten. Die vorhandenen Einzel-Targets
(`doc-immutable`, `doc-commits`, `doc-planning`, `doc-tracked`,
`doc-targets`, `doc-structure`) führen `--disable hostpaths` bereits und sind
von der Aktivierung **nicht** betroffen.

### Entwurf der Hard Rule `AGENTS.md` §3.11 — Kein host-lokaler absoluter Pfad in der Doku

**Aussage.** Kein Markdown-Dokument dieses Repos nennt einen host-lokalen
absoluten Pfad — Wurzel-Segment `Development`, `home`, `Users`, `Volumes`,
`mnt`, `media` oder ein Windows-Laufwerks-/UNC-Muster — in Prosa oder
Inline-Code; ein Schwester-Artefakt wird besitzer-qualifiziert zitiert
(`pt9912/d-check/Dockerfile`), nicht über seinen Pfad auf einem
Entwicklerrechner.

**Was der Sensor deckt — und was nicht.** Die durchsetzbare Hälfte trägt das
`hostpaths`-Modul in `make docs-check` (`make gates`). Es deckt: `.md`-Dateien
unter `scan.roots`, in Prosa und Inline-Code. Es deckt **nicht**:
Fenced-Code-Blöcke (dort sind Beispiel-Pfade erlaubt — Modul-Design, ohne
Opt-out-Marker), **relative** Pfade (ein relatives Ziel, das auf einen
Nachbarn zeigt, wird nicht geprüft), **Nicht-Markdown-Dateien** — die
Skriptkommentare in `Makefile`, `tools/**` und `harness/mk/**` liest der Scan
nicht —, und Dateien unter `scan.ignore` (`.harness/**` vendored Baseline,
`**/*.template.md`). Die Zusage „keine absoluten Pfade in diesem Repo" trägt
damit genau die **gescannte Markdown-Fläche**; Skript-Kommentare und
Fenced-Beispiele sind benannte, nicht gedeckte Ränder.

**Falsch** (Form: der Pfad beginnt mit dem Wurzel-Segment eines
Entwicklerrechners; der Platzhalter steht für dieses Segment, weil die Regel
auch das Beispiel deckt):

```text
Real geprüftes Vorbild: <Host-Wurzel>/d-check/Dockerfile (Stage `coverage`)
Die Musterquelle (<Host-Wurzel>/d-check/.github/, gelesen, nicht kopiert) ...
```

**Richtig:**

```text
Real geprüftes Vorbild: `d-check`s `Dockerfile` (Stage `coverage`)
Die Musterquelle (`d-check`s `.github/`, gelesen, nicht kopiert) ...
```

**Begründung.** Ein host-lokaler absoluter Pfad ist eine Aussage über einen
Rechner, nicht über das Repo: nicht portabel, für Mitlesende unauflösbar, und
er verrät das Maschinen-Layout. Die Regel und die Zusage sind identisch; das
`hostpaths`-Modul ist ihre durchsetzbare Hälfte — der Rand (Fenced-Code,
relative Pfade, Nicht-Markdown, `scan.ignore`) ist **benannt**, nicht
überzogen.

**Träger und Anker.** Die Regel wirkt über das `hostpaths`-Modul in
`make docs-check` (`make gates`); die Entscheidung trägt
[`ADR-0072`](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md). Die
Modul-Semantik (Präfixliste, Fence-Ausnahme, Windows-Muster) steht **einmal**
in `harness/sensors/docs-check.md` — §3.11 nennt sie nicht erneut, sonst
entstünde eine zweite Quelle für dieselbe Aussage (`AGENTS.md` §3.7).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nicht aktivieren (Status quo) | kein Korrekturaufwand, keine §3.5-Frage | der Auftrag bleibt unerfüllt; host-lokale Pfade bleiben unbemerkt im Repo; eine Zusage ohne Prüfung ist eine Absichtserklärung |
| B — aktivieren **mit** Ausschlussliste (`scope.ignore` für die zwei `Accepted` ADRs, `done/`, `docs/reviews/**`) | nur 2 lebende Stellen zu korrigieren; §3.5 und die Beleg-Konvention bleiben unangetastet | **vom Auftraggeber verworfen** („ohne Ausnahme"); die Ausnahme wäre eine zweite Quelle für „wo gilt die Zusage nicht" und altert; ihr einziger Grund ist §3.5, den [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) sauber nachzieht |
| C — aktivieren, nur die lebende Doku korrigieren, den Rest ausnehmen | kleinster Eingriff | wie B, aber inkonsistent: die Zusage gälte für `spec/`/`harness/`, nicht für `done/` — dieselbe Willkür, nur ungleichmäßig |
| D — aktivieren, **alle 31** korrigieren, §3.5 klassenweise nachziehen ([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)) — **gewählt** | die Zusage gilt repo-weit und ausnahmefrei; **eine** Zitationsform für alle 31; der eine verbotene Fall (`Accepted` ADR) wird über eine benannte, enge Klasse gelöst statt über ein Ventil | [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) öffnet einen schmalen Kanal in `Accepted` ADRs — der Preis ist benannt und begrenzt (nur Zitat-Gerüst, nie Normatives) |
| E — aktivieren und `ADR-0051`/`0054` per Folge-ADR `Supersedes` „bereinigen" | kein Eingriff in `Accepted` Text | **entfernt die Pfade nicht**: der alte, host-lokale Text bleibt in der abgelösten ADR stehen und wird weiter gescannt — der Auftrag „alle entfernt" bliebe unerfüllt; und nichts Normatives ist falsch, `Supersedes` ist das Werkzeug für *unwahre* Klauseln ([`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md)/[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)), nicht für eine Zitat-Form |

**Fazit:** D. B und C sind die Ausnahmelösung, die der Auftrag ausschließt; E
erreicht das Ziel nicht.

## Konsequenzen

- Positiv: Die Zusage „kein host-lokaler Pfad in der Doku" gilt repo-weit und
  **ohne Ausnahme** — eine Prüfung, ein Ergebnis, keine Ausnahmeliste.
- Positiv: **eine** Zitationsform für alle 31 Stellen; die drei
  Schwester-Artefakte werden besitzer-qualifiziert und bleiben für Mitlesende
  auflösbar.
- Positiv: §3.6 bleibt unberührt — die Aktivierung verschärft, sie lockert
  nichts.
- Negativ mit Grenze: Der Sensor deckt nur die **gescannte Markdown-Fläche**.
  Die zwei bekannten Nicht-Markdown-Stellen (`Makefile`, `tools/coverage-gate.sh`
  — je ein Skriptkommentar mit demselben Pfad) werden in derselben Folgearbeit
  mitkorrigiert, aber **kein Gate hält sie**; unentdeckt bleiben außerdem
  Fenced-Beispiele, relative Pfade und alles unter `scan.ignore`. Das ist
  benannt, nicht weggedefiniert.
- Negativ (Übergangszustand): Die Aktivierung macht `make docs-check` ab dem
  Commit der Config-Änderung rot, bis die 31 Stellen korrigiert sind. Die
  Reihenfolge — **erst korrigieren, dann aktivieren** (oder beides in einem
  Zug) — trägt die Closure; ein grüner Zwischenstand ist nicht behauptet.
- Folgepflicht (Implementer-Zug, kein Produkt-Code): die 31 Korrekturen in der
  Form aus §Entscheidung Punkt 3; `AGENTS.md` §3.11 (Wortlaut oben);
  `harness/sensors/docs-check.md` §Vertrag um `hostpath-forbidden` und §Grenze
  um die Fence-/Relativ-/Nicht-Markdown-Ränder ergänzen; `harness/README.md`
  §Sensors-Zeile; die zwei Skriptkommentare (`Makefile`,
  `tools/coverage-gate.sh`) in derselben Form mitkorrigieren. Beide
  `Accepted` ADRs bekommen ihre §Geschichte-Zeile aus
  [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md).
- Hinweis: Die Aktivierung ändert **nicht** die bestehenden Einzel-Targets
  (`doc-immutable`, `doc-commits`, …) — sie führen `--disable hostpaths`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `hostpaths` (Digest aus `d-check.mk`), aktiviert über die `modules:`-Liste | kein Dokument unter `scan.roots` nennt in Prosa oder Inline-Code einen host-lokalen absoluten Pfad: **0 Befunde** `hostpath-forbidden` | `make docs-check` (in `make gates`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** `d-check` ändert das Modul `hostpaths` — die
Default-Präfixliste, die festen Windows-Muster, die Fence-Behandlung, oder es
entsteht Ventil/Marker —, dann neu messen und diese Entscheidung nachziehen;
**(b)** `scan.roots`/`scan.ignore` ändern sich so, dass eine neue Dateiklasse
in den Scan tritt (etwa `.harness/**` nicht mehr ausgenommen) — dann neu
messen und die neuen Stellen korrigieren; **(c)** es entsteht ein zitiertes
Schwester-Artefakt **ohne** host-/owner-qualifizierbare Identität (kein
`owner/repo`-Slug) — dann braucht die Zitationsform einen zweiten Fall;
**(d)** `ADR-0051`/`0054` werden später supersedet und ihr host-lokaler Text
fällt weg — dann ist der Anlass der
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)-Klasse für
diese zwei entfallen, die Klasse selbst bleibt. Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Anordnung „`hostpaths` einbauen, ohne jede Ausnahme"; 31 Befunde eigenständig nachgemessen; Aktivierung ohne Ausschlussliste, Zitationsform, Hard-Rule-§3.11-Entwurf | der Architect-Verdikt dieses Zugs (hostpaths-Aktivierung ohne Ausnahme) |
| 2026-09-15 | Zitat-Korrektur — host-lokale Pfade ersetzt (`ADR-0073`) | `df47282` |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0072` — eine **Zitat-Korrektur** ausgenommen
([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md))
(Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
