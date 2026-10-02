# ADR-0143: `handbuch-public-doc-check` — Gate für kennungsfreie Nutzerdokumentation (Schärft ADR-0087, ADR-0088, ADR-0090)

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-02

**Autor:** pt9912 (Architect-Rolle, Modul 8; anderer Kontext als der Planner-Lauf
des Slice `handbuch-public-doc-check-gate-und-skill`, dessen Fragen 1 bis 6 diese
ADR beantwortet)

**Bezug:** [`LH-QA-OPS-001`](../../../spec/lastenheft.md) (Dokumentation für
Betreiber; die Nutzerdokumentation ist die Betreiber- und Integrator-Sicht),
[ADR-0134](0134-sdk-public-doc-check-gate-make-gates.md) (Muster: netzloses
`grep`-Gate „keine interne Kennung“, Aufnahme in `GATE_CHECKS`),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen: der
Skill führt das Ursprungs-Wort statt einer Berichts-Kennung),
`AGENTS.md` §3.6 (Träger-Pflicht für ein neues Gate), §3.5 (Accepted-ADRs sind
unberührbar), §3.1 (Docker-only; Host-Klasse `bash`/`git`/`grep`);
Slice-Plan `docs/plan/planning/open/slice-handbuch-public-doc-check-gate-und-skill.md`.

**Schärft:** [ADR-0087](0087-beispiel-clients-csharp-kotlin.md) (Abschnitt 6 /
Folgepflicht Handbuch-Zug), [ADR-0088](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)
(Folgepflicht Träger-Zug) und [ADR-0090](0090-beispiel-clients-volle-matrix.md)
(Festlegung 7 / Folgepflicht Handbuch-Zug) im **Wortlaut der Historie-Zeile**: die
Forderung „Zeile in der Änderungshistorie des Handbuchs im selben Zug“ bleibt
unverändert gültig; diese ADR ergänzt, dass die Zeile **ohne interne Kennung und in
Betreibersicht** geschrieben ist. Die drei Accepted-ADRs werden nicht geändert
(`AGENTS.md` §3.5); [`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md)
(Beispiel-Zeile „die Änderungshistorie des Handbuchs trägt die Zeile“) bleibt wahr.
Keine Spec-Stelle wird geschärft (Prozess-ADR ohne Spec-Stratum).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Anlass.** Das Benutzerhandbuch trug interne Kennungen und Chronik; die
Bereinigung ist committet (`7d1c4611`, `de0d899a`, `983e209a`). Ein Wächter, der den
Rückfall fängt, fehlte: bisher fing allein der Reviewer ihn durch Lesen.

**(2) Stand der Nutzerdokumente.** Gemessen am Stand `cfb68a31`, 2026-10-02, mit
`git grep -cE "$P" -- <Datei>` (Muster `$P` aus Festlegung 3), Trefferzeilen je Datei:
`benutzerhandbuch.md` 0, `benutzerhandbuch-standard.md` 0, `version.md` 0,
`releasing.md` 26, `bench-abdeckung.md` 4, `ci-matrix-abdeckung.md` 3,
`e2e-abdeckung.md` 74, `sdk-e2e-abdeckung.md` 21. Dieselben acht Zahlen liefert das
Muster des Planners (Slice-Plan §3, ERE ohne Wortrand-Schärfung), gemessen im selben
Lauf. Links nach `docs/plan/` oder `docs/reviews/` (Muster `$L`, Festlegung 3):
`releasing.md` 7, `bench-abdeckung.md` 1, `ci-matrix-abdeckung.md` 1, alle anderen 0
(gemessen, `git grep -cE "$L"`).

**(3) Warum die Dateien außerhalb liegen.** `releasing.md` ist Maintainer-Doku zum
Release-Prozess; ihre Kennungen (Entscheidungs-Links, Änderungshistorie mit
Review-Findings) sind dort gewollt. Die vier `*-abdeckung.md` sind Erzeugnisse der
Runner (`make test-integration`, `make bench`, `make doc-ci-matrix`, die SDK-
Integrationsziele); die Kennung je Zeile ist ihr Inhalt und die Eingabe von
`make doc-trace` (`.d-check.yml`, `trace.coverage`).

**(4) Der Tabellentest und das Muster.** `ADR-0134` zeigt die Form (Skript,
Sensor-Vertrag, Tabellentest als Werkzeug, `GATE_CHECKS +=`); hier ist der Gegenstand
eine Datei-Auswahl statt eines Baums, und das Muster braucht eine Wortrand-Schärfung
(Festlegung 3).

## Entscheidung

Wir wählen **ein netzloses Gate `handbuch-public-doc-check` in `GATE_CHECKS`**: ein
`grep`-Wächter über eine **benannte Liste** der Nutzerdokumente, ergänzt um eine
**Vollständigkeitsprüfung** des Verzeichnisses, plus einen kleinen Skill statt einer
Rolle. Das Kapitel `### Änderungshistorie` des Handbuchs **bleibt**; es trägt keine
Kennungen.

### Teilfrage 1 — Reichweite des Gates (Planner-Frage 1 und 2)

| Option | Fehlerbild | Pro | Contra |
|---|---|---|---|
| A — Dateiliste im Skript | eine **neue** Nutzerdatei bleibt **unbewacht** (still grün) | eine Erzeugnis-Datei färbt nie falsch rot | der Fehler ist lautlos; eine Umbenennung der genannten Datei lässt `grep` ins Leere laufen |
| B — Verzeichnis `docs/user/` minus Ausnahmeliste | eine **neue Erzeugnis-Datei** färbt das Gate **falsch rot** (laut, sichere Richtung); `releasing.md` stünde dauerhaft in der Ausnahmeliste | neue Nutzerdateien sind sofort bewacht | die Ausnahmeliste ist eine Schuldenliste (`releasing.md`), genau die „stille Liste“, die Festlegung 5 vermeiden will |
| **C — zwei benannte Listen plus Vollständigkeitsprüfung (gewählt)** | eine Datei unter `docs/user/*.md`, die in keiner der beiden Listen steht, färbt den Lauf rot (Exit 2, Meldung „unklassifiziert“); eine genannte, aber fehlende Datei ebenfalls | beide Fehlerbilder von A und B werden laut: eine neue Datei **erzwingt die Klassifikation** (geprüft oder ausgenommen); die Ausnahmen sind benannt, im Sensor-Vertrag begründet und stehen im Skript | etwa zehn Zeilen mehr Skript; jede neue Datei unter `docs/user/` kostet einen Eintrag |

### Festlegungen

1. **Reichweite.** Geprüft werden `docs/user/benutzerhandbuch.md`,
   `docs/user/benutzerhandbuch-standard.md` und `docs/user/version.md`. `version.md`
   gehört dazu: eine Zeile, aber das Handbuch verweist auf sie und der Release-
   Workflow liest sie; die Aufnahme kostet nichts (0 Treffer, gemessen, Kontext 2).
   **Ausgenommen** sind `docs/user/releasing.md` (Maintainer-Doku, Kontext 3) und die
   vier Erzeugnisse `bench-abdeckung.md`, `ci-matrix-abdeckung.md`,
   `e2e-abdeckung.md`, `sdk-e2e-abdeckung.md`. Das Skript klassifiziert jede
   `*.md`-Datei unter `docs/user/` (Wurzel überschreibbar wie bei `ADR-0134`): geprüft,
   ausgenommen oder **unklassifiziert** (Exit 2, Dateiname auf stderr); eine in
   einer Liste genannte, nicht vorhandene Datei ist ebenfalls Exit 2. Nicht-`.md`-
   Dateien (Bilder) sind nicht Gegenstand.
2. **`releasing.md` dauerhaft außerhalb (Planner-Frage 2).** Entschieden: kein
   Folge-Slice für eine Bereinigung. Begründung: die Datei ist Maintainer-Doku mit
   bewussten Entscheidungs-Links; eine Bereinigung wäre eine Inhaltsänderung ohne
   Nutzen für den Betreiber-Leser, und ein Gate, das sie einschlösse, verlangte sie
   zuerst. Die Ausnahme steht benannt in der Ausnahmeliste, nicht als Schuld
   (Re-Evaluierungs-Trigger). Die Zahl 26 ist gemessen (Kontext 2).
3. **Muster (Planner-Frage 4).** Zeilenweises `grep -InE`, zwei Muster:

   ```text
   P = \b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]
   L = docs/(reviews|plan)/|\.\./(reviews|plan)/
   ```

   Gegenüber dem Muster des Planners sind zwei Stellen geschärft: `MR-`/`CO-`
   verlangen **drei** Ziffern (das ID-Schema `MR-<NNN>`/`CO-<NNN>`,
   `harness/conventions.md` MR-000), und `slice-`/`welle-` verlangen am linken Rand
   Zeilenanfang oder ein Zeichen **außerhalb** von Buchstabe, Ziffer, `_` und `-`
   (damit trifft `byte-slice-x` nicht; `\b` allein trüge das nicht, der Bindestrich
   ist Wortgrenze). Das Muster ist **ERE**, kein PCRE: `grep -P` gehört nicht zur
   Host-Klasse aus `AGENTS.md` §3.1. Groß-/Kleinschreibung zählt (die Kennungen sind
   so geschrieben; `Slice-1` trifft nicht, benannte Grenze im Sensor-Vertrag).
   Befunde:

   | Zeile | Ergebnis mit `P` (`grep -nE`, gemessen 2026-10-02 an Eingabezeilen aus `printf`) |
   |---|---|
   | `slice-foo`, `welle-x`, `(slice-1)`, `a ADR-0134`, `CO-123`, `MR-001`, `LH-FA-SST-009`, `LH-FA-X`, `ARC-001`, `BEO-PGC/x` | Treffer |
   | `byte-slice-x`, `my-slice-1`, `CO-2 Emissionen`, `ALCO-123`, `x_CO-123`, `ZMR-100`, `sdk-csharp-v1`, `NR-0001`, `Slice-1`, `ADR-` | kein Treffer |

   Am echten Baum liefert `P` dieselben acht Zeilenzahlen wie in Kontext 2
   (gemessen, derselbe Lauf). Verallgemeinerung auf künftige Handbuch-Texte ist
   *hergeleitet*, nicht erprobt.
4. **Ausnahmen (Planner-Frage 3): keine.** Weder im Fence noch für Befehlsbeispiele;
   `grep` liest Zeilen. Ein Beispiel mit Kennung im Handbuch ist selbst der Fehler.
   Befehls-Tag-Muster wie `sdk-csharp-v*` treffen keines der Muster (gemessen:
   0 Treffer im Handbuch, Kontext 2). Eine Ausnahmeform (Opt-out-Marker) braucht eine
   eigene ADR. Folge für den Implementer: ein Satz in `benutzerhandbuch-standard.md`
   über Kennungsfreiheit **nennt keine Beispiel-Kennung** („keine internen
   Kennungen“, nicht „z. B. `ADR-0001`“), sonst färbt er das Gate rot.
5. **Gate-Aufnahme und Verdrahtung (Planner-Frage 6).** Name `handbuch-public-doc-check`
   (deckt Handbuch, Standard und `version.md`; `nutzerdoku-` wäre ein Missname, weil
   `releasing.md` ebenfalls Nutzerdokumentation unter `docs/user/` ist). Skript
   `tools/harness/handbuch-public-doc-check.sh`, Tabellentest
   `tools/harness/run-handbuch-public-doc-check-tests.sh`, Sensor-Vertrag
   `harness/sensors/handbuch-public-doc-check.md`. Makefile-Heimat
   `harness/mk/doc-gate.mk` (die Dokumentations-Gates; `sdk.mk` ist die Heimat der
   SDK-Ziele): Ziel `handbuch-public-doc-check` samt `GATE_CHECKS +=`, Werkzeug-Ziel
   `test-handbuch-public-doc-check` ohne Gate-Aufnahme (`ADR-0134` Teilfrage 2,
   Option A, gleiche Begründung). Netzlos, nur `bash`/`git`/`grep`; das Skript steigt
   über `git rev-parse --show-toplevel` auf die Wurzel, der optionale erste Parameter
   ist die Wurzel für den Tabellentest. `ci.yml` fährt `make gates` unverändert, kein
   Workflow ändert sich (`AGENTS.md` §3.10 greift nicht; am Start des Slice zu
   messen).
6. **Schärfung der drei ADRs (Planner-Frage 5).** Siehe `Schärft:` im Kopf. Die
   Änderungshistorie des Handbuchs trägt je Version eine Zeile, in Betreibersicht
   (was ändert sich für den Betreiber), **ohne** `LH-*`, `ADR-*`, `SPEC-*`, `ARC-*`,
   Slice-/Welle-/BEO-Namen. Wer die Zeile schreibt, kann ihren Anlass in Commit und
   Plan finden; das Handbuch braucht ihn nicht.
7. **Regeltexte und Standard (Auftraggeber-Entscheidung).** Die Pflicht „`Version:`
   hochzählen **und** eine Zeile in `### Änderungshistorie`“ bleibt in
   `.claude/commands/implement-slice.md` Schritt 17 und im HIGH-Punkt „Handbuch-
   Versionshistorie nicht fortgeschrieben“ von `.harness/skills/reviewer.md`; sie wird
   um „ohne Kennungen, Betreibersicht“ und den Verweis auf Skill und Gate ergänzt,
   nicht abgeschwächt. `docs/user/benutzerhandbuch-standard.md` (Aufzählung
   „Changelog oder Änderungshistorie“ und Kapitel „11. Änderungshistorie“) trägt die
   Kennungsfreiheit; das Kapitel bleibt. Diese Änderungen sind Implementer-Arbeit des
   Slice, nicht dieser ADR.
8. **Skill statt Rolle (Planner-Frage 6, Skill-Ort).** `.harness/skills/nutzerdoku-
   schreiben.md`, in der Form der Vorgänger `reviewer.md` und
   `closure-note-reviewer.md` (gemessen: `ls .harness/skills` nennt genau diese beiden
   Dateien); Eintrag in `harness/README.md` §Guides, Verweis aus Schritt 17 und dem
   Reviewer-Skill. Inhalt laut Slice-Plan §1 (Ist-Zustand, keine Kennungen, keine
   Links nach `docs/plan/`/`docs/reviews/`, Betreibersicht, Ursprungs-Wort ohne
   Berichts-Kennung, Historie-Zeile, Verweis aufs Gate als Fangnetz). Der Skill ergänzt
   keine Hard Rule in `AGENTS.md`. Eine Rolle „Handbuch-Schreiber“ entsteht nur bei
   erneutem Drift trotz Gate und Skill (Auftraggeber-Entscheidung).
9. **Milderung braucht einen Träger.** Wie bei `ADR-0134` Festlegung 4: weniger Muster
   oder eine Verschiebung einer Datei von „geprüft“ nach „ausgenommen“ ist eine
   Umfangsänderung dieses Gates im Sinn von `AGENTS.md` §3.6 und braucht eine eigene
   ADR. Eine **neue** Datei in einer der beiden Listen zu führen, ist keine Milderung
   des Bestands und der erwartete Weg der Klassifikation.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Reviewer-Lesen allein | kein Aufwand | der Rückfall war bereits eingetreten; Lesen ist die einzige Falsifikation, sie fehlt am Push |
| B — eigene Rolle „Handbuch-Schreiber“ | trennt Schreiben und Prüfen | Auftraggeber: nur bei erneutem Drift; ein Gate plus Skill ist der kleinere Eingriff |
| C — Gate nur auf `benutzerhandbuch.md` | minimal | der Standard (Form des Handbuchs) und `version.md` blieben ungeschützt, und die Frage der neuen Datei bliebe ungelöst |
| **D — Gate mit Klassifikation, Skill, Regeltext-Ergänzung (gewählt)** | fängt den Rückfall bei jedem Push; beide Fehlerbilder der Datei-Auswahl werden laut; Pflicht „Historie“ bleibt | zehn Zeilen Skript mehr als nötig; jede neue Datei unter `docs/user/` braucht einen Listeneintrag |
| E — Historie-Kapitel streichen | nichts zu prüfen | widerspricht dem Auftraggeber („Die History kann ohne Kennung auskommen“) und den drei ADRs |

## Konsequenzen

- Positiv: Ein Rückfall in Kennungen oder in Links nach `docs/plan/`/`docs/reviews/`
  fällt bei Push und PR auf (`ci.yml` fährt `make gates`).
- Positiv: Eine neue Datei unter `docs/user/` zwingt zur Entscheidung „Nutzerdoku oder
  Erzeugnis“, statt still unbewacht oder falsch rot zu sein.
- Negativ (Grenze): Das Gate liest Kennungen, nicht Sinn. Chronik-Sprache ohne
  Kennung („früher … jetzt …“) bleibt grün; diese Hälfte tragen Skill und Reviewer.
  Groß geschriebene Namen (`Slice-1`) und Kennungsformen außerhalb der Muster bleiben
  grün.
- Negativ (akzeptiertes Negativ): Die Wächter-Logik ist im Gate-Lauf nicht gegen den
  Tabellentest gesichert (Werkzeug, `ADR-0134` Teilfrage 2); der Review-Diff liest eine
  Änderung am Wächter gegen den Sensor-Vertrag.
- Negativ: `releasing.md` und die Erzeugnisse sind von der Kennungsfreiheit
  ausgenommen; eine Kennung dort fängt nur der Reviewer.
- Folgepflicht (Implementer-Zug des Slice): Skript, Tabellentest, Sensor-Vertrag,
  `doc-gate.mk`-Ziele, `harness/README.md` (§Sensors, §Guides), Skill, die zwei
  Regeltexte, der Standard. Kein Folge-Slice für `releasing.md` (Festlegung 2).
- Folgepflicht: `spec/`, `spec/architecture.md`, `AGENTS.md` und die drei Accepted-ADRs
  bleiben unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Muster `P` an Eingabezeilen (erprobt) | zehn Zeilen mit Kennung treffen, zehn ohne treffen nicht, darunter der Wortrand `byte-slice-x` und `CO-2 Emissionen` (Tabelle in Festlegung 3; `grep -nE` an `printf`-Zeilen, gemessen 2026-10-02) | — |
| Muster `P` am echten Baum (erprobt) | `benutzerhandbuch.md`, `-standard.md`, `version.md` je 0 Trefferzeilen; `releasing.md` 26 und die vier Erzeugnisse 4/3/74/21 außerhalb (`git grep -cE`, gemessen 2026-10-02, Stand `cfb68a31`) | `make handbuch-public-doc-check` (nach der Umsetzung) |
| Gate-Verdrahtung (hergeleitet) | `handbuch-public-doc-check` steht in `GATE_CHECKS`; `make gates` fährt es | `make gates` |
| Mutation: Kennung (hergeleitet; **nicht erprobt**) | Eine in eine Kopie von `benutzerhandbuch.md` im Scratchpad eingefügte Kennung färbt das Skript mit Wurzel = Kopie-Verzeichnis rot (Exit 1); ebenso ein Link nach `docs/plan/`. Stelle, Instanz und Farbe trägt der Implementer in seinen Bericht ein; „der Implementer fährt sie“ ist eine Erwartung | `make test-handbuch-public-doc-check` |
| Klassifikation (hergeleitet; **nicht erprobt**) | eine zusätzliche `.md`-Datei im Wurzel-Verzeichnis des Tabellentests färbt Exit 2; eine genannte, fehlende Datei ebenfalls; die drei geprüften Dateien und die fünf ausgenommenen sauber/ausgenommen Exit 0 | `make test-handbuch-public-doc-check` |

## Re-Evaluierungs-Trigger

Ein erneuter Rückfall in Nutzerdokumentation trotz Gate und Skill öffnet die Frage nach
einer Rolle „Handbuch-Schreiber“ (Auftraggeber-Entscheidung). `releasing.md` wird neu
bewertet, wenn sie sich an Betreiber statt an Maintainer richtet. Ändert sich die
Kennungs-Syntax (`harness/conventions.md`), ist das Muster `P` nachzuziehen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-02 | Accepted — Architect-Entscheidung als Vorstufe des Umsetzungs-Slices | `docs/plan/planning/open/slice-handbuch-public-doc-check-gate-und-skill.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0143` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
