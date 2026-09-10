# Harness

---

## Purpose

Dieser Harness verbindet bestehende Spezifikationen, ADRs,
Planning-Dokumente und Gates. Er ist **kein Ersatz** für `spec/` oder
`docs/`, sondern ein **Einstiegspunkt** für Menschen und AI-Code-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, **gewinnt die
kanonische Quelle**, und diese Datei wird angepasst.

Strukturregeln (Verzeichniskonvention, ID-Schemata, Modus-Deklarationen
pro Sub-Area, Zusatzklassen für Sensors-Bindung) sowie Adaptionen ggü.
der adoptierten Baseline leben in [`conventions.md`](conventions.md).
Diese Datei dupliziert sie nicht.

## Source precedence

| Rang | Datei | Charakter |
|---|---|---|
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | vertraglich abnahmebindend |
| 2 | [`spec/pflichtenheft.md`](../spec/pflichtenheft.md) | technisch fortschreibbar |
| 3 | [`spec/architecture.md`](../spec/architecture.md) | Komponenten/Sequenzen, meilensteinfrei |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Wellen-Sequenz |
| 6 | `docs/user/*` *(falls vorhanden)* | Operations, Quality, Releasing | <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
| 7 | [`README.md`](../README.md) | Projekt-Überblick |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 9 | diese Datei | Harness-Einstieg |

> Die Ränge 1–3 sind die **drei Spec-Straten** — Vertrag, Technik, Sicht —,
> und alle drei sind obligatorisch (Baseline-Regelwerk
> `grundlagen-referenz-richtung.md` §Spec-Straten). **Adaption ist die
> Zwei-Straten-Form, nicht die Drei-Straten-Form**: Wer Rang 2 streicht,
> deklariert das als `MR-<NNN>` in [`conventions.md`](conventions.md) und
> nummeriert neu (dann acht Ränge).

## Guides (Feedforward-Quellen)

<!--
Was lenkt den Agenten *vor* der Handlung? Pointer, kein Inhalt.
-->

| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen, IDs, Akzeptanzkriterien |
| [`spec/pflichtenheft.md`](../spec/pflichtenheft.md) | technische Details, Defaults |
| [`spec/architecture.md`](../spec/architecture.md) | Komponenten, Schichten, Constraints |
| [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |
| `.harness/skills/reviewer.md` | Reviewer-Skill: HIGH-Liste, Kategorien-Regeln, Negativbefund-Pflicht, Output-Schema (Modul 10) — nächste Rolle nach Schritt 8 des Minimal Agent Workflow, nicht Teil der Implementer-Eingabe |
| `.harness/baseline/<tag>/regelwerk/` (vendored; `README.md` = Index) | adoptiertes Betriebsregelwerk in Agenten-Kurzform — **präsente nachschlagbare Vertiefung**, pro Entscheidung abschnittsweise (siehe [`AGENTS.md`](../AGENTS.md) §1); derivativ, Stand/Tag siehe [`conventions.md`](conventions.md) §Baseline |
| `.harness/baseline/<tag>/templates/` (vendored, parallel) | Referenz-Form der Skelette, auf die das Regelwerk mit `../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu `regelwerk/`); Vorlagen zum Kopieren-und-Ausfüllen |

## Sensors (Feedback-Gates)

<!--
WICHTIG: Nur Befehle aufzählen, die im Makefile *existieren*.
Halluzinierte Gates sind die häufigste Form von Harness-Lüge (Modul 13).

Drei Spalten — kein Lauf-Status:
- Target:  der Make-Befehl.
- Vertrag: was prüft das Gate (was wäre verletzt, wenn es rot wird).
- Bindung: strukturelle Referenzen — Carveout-ID (`CO-<NNN>`),
  Slice-ID, Schwelle, Image-Hash, ADR-ID. NICHT der Lauf-Status,
  sondern was das Gate *strukturell trägt*.

Lauf-Wahrheit pro Commit liegt in CI (Badge/Dashboard), nicht hier
(`harness/README.md` ist Rang 9 in der Source Precedence).
Strukturell rote Gates (dauerhaft rot) bekommen einen Carveout in
`docs/plan/carveouts/CO-<NNN>-…` mit Auflösungs-Trigger und Folge-Slice
(Modul 7); die Bindung-Spalte verweist auf die `CO-<NNN>`-ID, die
Begründung lebt im Carveout, nicht hier.

Bei d-check-Einsatz (≥ v0.73.0) deckt Modul `reviews` (Ziel `doc-reviews`,
Review-Report-Deckung für `done/`-Slices mit Review-DoD-Haken) die
Code→Review-Kante ab, Modul `planning` (Ziel `doc-planning`,
Planning-Lifecycle-Konsistenz) die Verify→Closure-Kante — beide aus dem
Lebenszyklus-Diagramm in Modul 1. Nur eintragen, wenn das Ziel im
Makefile existiert (siehe oben).

NICHT-GATES: Ein Target, das der Agent braucht, aber das nichts über den
Zustand des Repos urteilt — es *bewegt* (Slice-Move), *misst* (Latenz) oder
*sagt*, was ein schreibender Lauf täte — steht in der zweiten Tabelle und
trägt `kein Gate` IN DER ZEILE SELBST, in der Spalte, die hier die Bindung
führt — nicht in Prosa daneben. Weglassen ist nur für Targets richtig, die
niemand braucht (Modul 13 §Vorhanden ≠ behauptet).

WÄCHST DIE SEKTION: Die Tabelle bleibt klein, die Prosa darunter nicht.
Braucht ein Gate mehr als EINEN SATZ — Deckungsgrenze, Ausgabe-Bedeutung,
Exit-Codes, Abbruch-Bedingungen —, wandert das nach. Ob der Überhang schon
unter der Tabelle steht oder in die Zelle gedrängt wurde, ist dieselbe Sache:
eine Zelle, die zum Absatz geworden ist, ist der Fund, nicht die Ausnahme.
`harness/sensors/<target>.md`, und die **Target-Zelle wird zum Link darauf**
— wie die `MR`-Zelle im Adaptions-Block. Der Link ist kein Komfort: Er ist die
einzige Fassung dieser Zuordnung, die der Link-Sensor prüft. Eine bloße
Namenskonvention (`make X` -> `sensors/X.md`) bleibt still grün, wenn die
Datei verschwindet und die Zeile stehen bleibt. Seine Grenze: Er prüft EINE
Richtung — ob das Ziel existiert; eine Datei ohne Index-Zeile und eine Zeile
auf die falsche Datei bleiben still grün, und geprüft wird nur, wo ein
Link-Sensor über `harness/` läuft. Kein `sensors/done/`:
ein retiriertes Gate verschwindet, `git` hält seine Geschichte. Was das
Werkzeug selbst deckt (welcher Test welche Hälfte trägt), gehört NICHT
dorthin, sondern in seine ADR/Spec-Zeile/seinen Skriptkopf.
-->

| Target | Vertrag | Bindung |
|---|---|---|
| `make baseline-verify` | verifiziert die vendored Baseline netzlos (Integrität + Vollständigkeit gegen `SHA256SUMS`) | [`harness/sensors/baseline-verify.md`](sensors/baseline-verify.md) |
| `make docs-check` | kaputte Referenzen in der Markdown-Doku (links, anchors, ids, matrix, versions, structure) | [`harness/sensors/docs-check.md`](sensors/docs-check.md) |
| `make a-check` | prüft die Hexagon-Schichten-Edges aus `.a-check.yml` gegen den Go-Baum — netzlos, read-only, digest-gepinntes Release-Image | [`ADR-0041`](../docs/plan/adr/0041-a-check-maschinenform-architekturpruefung.md) · [`harness/sensors/a-check.md`](sensors/a-check.md) |
| `make commit-traceability` | Commit-Message-Traceability: je Message der Range ≥ 1 `LH-*`-/`ADR-*`-Kennung (d-check Modul `commits`, Befund `commit-untraceable`) und keine `SPEC-*`/`ARC-*`-Kennung im Betreff (`tools/harness/commit-traceability.sh`); Standing-Gate über die letzten 5 Commits (`RANGE=base..head` überschreibt) | [`ADR-0045`](../docs/plan/adr/0045-commit-traceability-standing-gate.md) · seit slice-006 |
| `make gates` | alle inneren Gates (baseline-verify, docs-check, a-check, commit-traceability), Nachweis-Stempel zuletzt | — |
| `<make-target>` | volle Closure | Image-Hash `sha256:…` (Modul 14) |

**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate:**

| Target | Tut was | Bindung |
|---|---|---|
| `make image` | baut das OCI-Image (Multi-Stage, digest-gepinnt); **Image-Digest** nach `harness/image-hash.txt` — der Digest ist der **Lauf-Beleg** des letzten offiziellen `make image`-Laufs (builder- und lauf-gebunden, kein Inhalts-Fingerabdruck); ein Digest-Vergleich über Umgebungen/Läufe ist kein Staleness-Beweis, Inhalts-Streits entscheidet der sha256 des extrahierten Binaries (Container-Export). Ein Zug, der Build-Kontext-Dateien ändert, läuft `make image` vor seiner Closure; der Digest-Commit entfällt bei unverändertem Digest | kein Gate, [`ADR-0044`](../docs/plan/adr/0044-image-beleg-semantik.md), [`ADR-0042`](../docs/plan/adr/README.md) |
| `make image-stale` | meldet Base-Image-Drift: FROM-Digests des Dockerfile gegen die aktuellen Registry-Digests derselben Tags, plus Existenz des nächsten Major-Tags; braucht Netz | kein Gate, [`ADR-0039`](../docs/plan/adr/README.md) (Update = bewusster Digest-Commit, Modul 14) |
| `make test` | Unit-Tests im gepinnten Toolchain-Container (netzlos; Modul-Cache im Docker-Volume, Vorbereitung: `make mod-download`) | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) |
| `make test-store` | Adapter-Tests gegen reale PostgreSQL im Testcontainer (gepinnte Digests: Toolchain + PostgreSQL; DB-Daten im Container) | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) |
| `make test-replication` | Replication-Stream-Tests gegen reale PostgreSQL mit Publication und Logical Replication Slot (`wal_level=logical`, gepinnte Digests: Toolchain + PostgreSQL; DB-Daten im Container) | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) |
| `make test-integration` | MVP-Integrationstest gegen die Compose-Umgebung (`compose.yaml`): Umgebung frisch hochfahren → Schema-Rollout über d-migrate ([`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)) → Toolchain-Container gegen das Compose-Netz; gepinnte Digests, DB-Daten im Container | kein Gate, [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md), [`LH-QA-POR-003`](../spec/lastenheft.md) |
| `make schema-validate` | prüft das neutrale Schema-YAML mit d-migrate (netzlos; Vorlauf vor jedem `generate`/`migrate`) | kein Gate, [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md) |
| `make schema-rollout` | rollt das Schema mit d-migrate aus — `schema migrate --execute` mit Pflicht-Report (`tools/schema/plan.yaml`) und Rollback-Artefakt (`tools/schema/down.sql`); braucht DB-Zugang | kein Gate, [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md) |
| `make <mover>` | bewegt <…>, prüft nichts | kein Gate |
| `make <messung>` | misst <…> gegen <Schwelle> | kein Gate, ADR-<NNNN> |
| `make <vorschau>` | sagt, was <schreibender Lauf> täte; Ausgänge und Sperren in der verlinkten Datei | kein Gate |

**Aktueller Lauf-Status:** CI-Badge bzw. lokal `make help` / `make gates`.
**Rote Gates:** Begründung im verlinkten `CO-<NNN>` (siehe Bindung-Spalte), Modul 7.
**Nicht behauptet** (geplant): `make image-cve` (CVE-Scan des gebauten
Images, advisory) — die Aktivierungsbedingung ist seit slice-001 eingetreten
(der erste `make image`-Lauf läuft grün); das Target ist noch nicht
implementiert, bis dahin bleibt der Scan unausgesprochen.

<!-- Domänenspezifische Gates ergänzen, je nach Repo-Klasse: -->

## Traceability rules

- PRs/Commits **müssen** mindestens eine `<LH-*>` oder `ADR-*`-ID nennen.
- Mechanisch getragen (Standing-Gate, [`ADR-0045`](../docs/plan/adr/0045-commit-traceability-standing-gate.md)): `make commit-traceability` prüft je Message der letzten 5 Commits (`RANGE=base..head` überschreibt) beide Grenzen — positive Hälfte via d-check Modul `commits` (Befund `commit-untraceable`), Betreff-Grenze via `tools/harness/commit-traceability.sh` · seit slice-006 (ausgelöst durch das dritte Auftreten der Verstoß-Klasse, review-slice-006 F-2).
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Neue ADRs müssen im ADR-Index ergänzt werden.
- Änderungen an Planning-Dokumenten müssen die Lifecycle-Regeln beachten (open → next → in-progress → done; reine `git mv`-Commits siehe AGENTS.md §3.3).

## Safety and scope boundaries

<!--
Repo-spezifisch formulieren. Beispiele:

Für ein Referenz-Repo:
- Dies ist kein produktiver Service.
- Externer Cloud-Zugriff darf nicht für lokale Demo-Abnahme vorausgesetzt werden.
- Determinismus und Replayability sind Kernverträge.

Für ein Safety/Control-Repo:
- Markt-/Optimierungs-Output muss durch Statemachine, Constraint-Limiter, Ramp-Limiter fließen.
- Software-Stop ersetzt keine Hardware-Sicherheitsfunktionen.
- Produktion-Profile müssen fail-closed sein.

Für ein Policy/Compliance-Repo:
- Dieses Werkzeug ist keine Rechts-/Steuer-/Fachberatung.
- KI-Funktionen liefern Vorschläge, keine verbindlichen Entscheidungen.
-->

- <…>
- <…>

## Minimal agent workflow

1. Diese Datei lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe `harness/README.md` §Guides) →
Verifier. Kein Self-Review — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).

## Leseordnung

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt — die Menschen-Hälfte des Einstiegs:
drei bis fünf **geordnete** Zeiger, was ein neuer Mensch zuerst liest und was
bei Bedarf; eine Leseordnung, die alles nennt, ist keine.

1. <zuerst — z. B. `AGENTS.md` §Hard Rules>
2. <dann — z. B. `spec/lastenheft.md`>
3. <bei Bedarf — z. B. `harness/conventions.md`>
