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
| 6 | [`docs/user/`](../docs/user/) | Operations, Quality, Version; [`docs/maintainer/`](../docs/maintainer/) — Releasing |
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
| `.harness/skills/nutzerdoku-schreiben.md` | Skill „Nutzerdokumentation schreiben“: Ist-Zustand, Betreibersicht, keine internen Kennungen, Versionshistorie-Zeile — Feedforward für Schritt 17 des Implementer-Ablaufs; das Gate `make handbuch-public-doc-check` ist das Fangnetz ([`ADR-0143`](../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)) |
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
| `make docs-check` | kaputte Referenzen in der Markdown-Doku (links, anchors, ids, matrix, versions, structure, hostpaths, tracked) | [`harness/sensors/docs-check.md`](sensors/docs-check.md) |
| `make a-check` | prüft die Hexagon-Schichten-Edges aus `.a-check.yml` gegen den Go-Baum — netzlos, read-only, digest-gepinntes Release-Image | [`ADR-0041`](../docs/plan/adr/0041-a-check-maschinenform-architekturpruefung.md) · [`harness/sensors/a-check.md`](sensors/a-check.md) |
| `make commit-traceability` | Commit-Message-Traceability: je Message der Range ≥ 1 `LH-*`-/`ADR-*`-Kennung, keine `SPEC-*`/`ARC-*`-Kennung im Betreff; Standing-Gate über die letzten 5 Commits | [`ADR-0045`](../docs/plan/adr/0045-commit-traceability-standing-gate.md) · [`harness/sensors/commit-traceability.md`](sensors/commit-traceability.md) · seit slice-006 |
| `make coverage-gate` | Go-Test-Coverage über die netzlos prüfbare Fläche gegen `THRESHOLD` — bootstrap-aware Gate, geltende Stufe in [`harness/mk/coverage.mk`](mk/coverage.mk) | [`ADR-0071`](../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) · [`ADR-0054`](../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) · [`ADR-0077`](../docs/plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md) · [`harness/sensors/coverage-gate.md`](sensors/coverage-gate.md) · seit slice-049 |
| `make generated-sync` | prüft, dass der committete Protobuf-/gRPC-Code byte-gleich der Ausgabe des **gepinnten** Generators aus den committeten `.proto`-Quellen ist | [`harness/sensors/generated-sync.md`](sensors/generated-sync.md) · seit slice-090 |
| `make sdk-public-doc-check` | prüft mit `grep`, dass keine Datei unter `sdks/` eine interne Kennung oder einen Slice-/Welle-Namen trägt | [`ADR-0134`](../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) · [`harness/sensors/sdk-public-doc-check.md`](sensors/sdk-public-doc-check.md) · seit slice-sdk-public-doc-check-gate |
| `make handbuch-public-doc-check` | prüft mit `grep`, dass die Nutzerdokumente unter `docs/user/` keine interne Kennung und keinen Link nach `docs/plan/` oder `docs/reviews/` tragen | [`ADR-0143`](../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) · [`harness/sensors/handbuch-public-doc-check.md`](sensors/handbuch-public-doc-check.md) · seit slice-handbuch-public-doc-check-gate-und-skill |
| `make ausgabe-kennungen-check` | prüft mit `grep`, dass kein Ausgabe-Literal des Servers und seiner Betriebs-Skripte eine interne Kennung trägt | [`ADR-0144`](../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10 · [`harness/sensors/ausgabe-kennungen-check.md`](sensors/ausgabe-kennungen-check.md) · seit slice-meldungscodes-kennungsfreie-ausgaben |
| `make meldungscodes-check` | gleicht die Meldungscodes `PCF-<S><NNNN>` in Quelltext und Handbuch mit der Code-Tabelle `internal/domain/messagecode/codes.go` ab | [`ADR-0144`](../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 5 · [`harness/sensors/meldungscodes-check.md`](sensors/meldungscodes-check.md) · seit slice-meldungscodes-registry-fehlerkopf |
| `make gates` | alle inneren Gates dieser Tabelle, Nachweis-Stempel zuletzt | [`harness/sensors/gates.md`](sensors/gates.md) |
| `<make-target>` | volle Closure | Image-Hash `sha256:…` (Modul 14) |

**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate:**

| Target | Tut was | Bindung |
|---|---|---|
| `make image` | baut das OCI-Image (Multi-Stage, digest-gepinnt); **Image-Digest** nach `harness/image-hash.txt` | kein Gate, [`harness/targets/image.md`](targets/image.md), [`ADR-0044`](../docs/plan/adr/0044-image-beleg-semantik.md), [`ADR-0103`](../docs/plan/adr/0103-image-hash-lokal-statt-committet.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0042`](../docs/plan/adr/README.md) |
| `make image-stale` | meldet Base-Image-Drift der FROM-Digests des Dockerfile gegen die Registry; braucht Netz | kein Gate, [`harness/targets/image.md`](targets/image.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 (Update = bewusster Digest-Commit, Modul 14) |
| `make image-cve` | Trivy CRITICAL/HIGH gegen das publizierte GHCR-`:latest`-Image | kein Gate, [`harness/targets/image.md`](targets/image.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 6 |
| `make pin-stale-race`/`make pin-stale-pgtest`/`make pin-stale-dmigrate`/`make pin-stale-acheck` | Upstream-Pin-Freshness P3–P6: Digest einer Makefile-/`a-check.mk`-Variable gegen den Registry-Digest | kein Gate, [`harness/targets/pin-stale.md`](targets/pin-stale.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz |
| `make pin-stale-dcheck` | Upstream-Pin-Freshness P7: Digest-Drift und Tag-Frische des d-check-Pins | kein Gate, [`harness/targets/pin-stale.md`](targets/pin-stale.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz |
| `make pin-stale-baseline` | Upstream-Pin-Freshness P8: die adoptierte Kurs-Baseline-Version gegen den neuesten Release | kein Gate, [`harness/targets/pin-stale.md`](targets/pin-stale.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz |
| `make pin-stale-actions` | Upstream-Pin-Freshness P9: jede SHA-gepinnte `uses:`-Zeile gegen Tag-Mutation und Tag-Frische | kein Gate, [`harness/targets/pin-stale.md`](targets/pin-stale.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7, braucht Netz |
| `make pin-stale-all` | Upstream-Pin-Freshness P10: jede Digest-Pin-Referenz außerhalb von `docs/` und `.harness/` gegen den Index-Digest | kein Gate, [`ADR-0146`](../docs/plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md), [`harness/sensors/pin-stale-all.md`](sensors/pin-stale-all.md), braucht Netz · seit slice-pin-stale-alle-digest-pins |
| `make doc-tracked` | isolierter Einzel-Lauf des `tracked`-Moduls, Diagnose-Werkzeug neben `make docs-check` | kein Gate, [`harness/sensors/docs-check.md` §`make doc-tracked`](sensors/docs-check.md#make-doc-tracked) · seit slice-d-check-tracked-modul |
| `make doc-trace` | gibt die Requirements Traceability Matrix (RTM) auf stdout aus (advisory) | kein Gate, wie `make image-stale`, [`harness/targets/doc-trace.md`](targets/doc-trace.md), [`harness/sensors/docs-check.md`](sensors/docs-check.md) §Grenze · seit slice-d-check-trace-rtm |
| `make doc-ci-matrix` | schreibt `docs/user/ci-matrix-abdeckung.md` aus dem letzten erfolgreichen Lauf; braucht Netz | kein Gate, [`harness/targets/doc-trace.md`](targets/doc-trace.md), braucht Netz (wie `make image-stale`), [`ADR-0105`](../docs/plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md) |
| `make proto-generate` | erzeugt den Go-Code der beiden gRPC-Services aus den `.proto`-Quellen (Docker-only) | kein Gate, [`harness/targets/proto-generate.md`](targets/proto-generate.md), [`ADR-0060`](../docs/plan/adr/0060-grpc-streaming-mechanismus.md), [`ADR-0130`](../docs/plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) |
| `make test` | Unit-Tests mit Race-Detector im gepinnten Toolchain-Container (netzlos) | kein Gate, [`harness/targets/tier-tests.md`](targets/tier-tests.md), [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md) |
| `make test-store` | Adapter-Tests gegen reale PostgreSQL im Testcontainer | kein Gate, [`harness/targets/tier-tests.md`](targets/tier-tests.md), [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md), [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md); DB-Adapter-Coverage ([`sensors/db-adapter-coverage.md`](sensors/db-adapter-coverage.md), [`ADR-0071`](../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) |
| `make test-replication` | Replication-Stream-Tests gegen reale PostgreSQL mit Publication und Logical Replication Slot | kein Gate, [`harness/targets/tier-tests.md`](targets/tier-tests.md), [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md); DB-Adapter-Coverage ([`sensors/db-adapter-coverage.md`](sensors/db-adapter-coverage.md), [`ADR-0071`](../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) |
| `make test-notify` | `natsnotify`- und `natsstream`-Tests gegen einen echten NATS-Server im Testcontainer | kein Gate, [`harness/targets/tier-tests.md`](targets/tier-tests.md), [`ADR-0055`](../docs/plan/adr/0055-nats-change-notification-wecksignal.md), [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md), [`ADR-0137`](../docs/plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) · seit slice-052, erweitert seit slice-routing-nats-subjekt |
| `make test-integration` | Compose-Integrationstest: E2E-Rundläufe gegen den laufenden Feed-Container | kein Gate, [`harness/targets/test-integration.md`](targets/test-integration.md), [`ADR-0030`](../docs/plan/adr/0030-testpyramide.md), [`LH-QA-POR-003`](../spec/lastenheft.md) |
| `.github/workflows/upstream-drift.yml` | automatisiert die Achsen des Pin-Inventars nächtlich (advisory, je Achse fail-open) | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 |
| `.github/workflows/e2e.yml` | automatisiert `make image` und `make test-integration` auf jeden Pull Request und Push (nicht blockierend) | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md), [`LH-QA-POR-001`](../spec/lastenheft.md), [`LH-QA-POR-003`](../spec/lastenheft.md) |
| `.github/workflows/release.yml` | Release-Workflow auf `v*`-Tags: Image bauen, nach GHCR und Docker Hub pushen, GitHub-Release | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) |
| `.github/workflows/hub-description.yml` | synchronisiert die Docker-Hub-Beschreibung von `pt9912/pg-change-feed` mit `README.md` | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 8 |
| `.github/workflows/sdk-csharp-release.yml` | NuGet.org-Publish-Workflow für `PgChangeFeed.Client`, Trigger `sdk-csharp-v*`-Tag | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0106`](../docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md) Festlegung 4 · seit slice-sdk-csharp-publish-workflow |
| `.github/workflows/sdk-python-release.yml` | PyPI-Publish-Workflow für `pgchangefeed`, Trigger `sdk-python-v*`-Tag | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0107`](../docs/plan/adr/0107-python-pypi-zweites-sdk-package.md) Festlegung 5, [`ADR-0108`](../docs/plan/adr/0108-python-sdk-uv-statt-build-twine.md) §Entscheidung Festlegung 1 · seit slice-sdk-python-publish-workflow |
| `.github/workflows/sdk-kotlin-release.yml` | Publish-Workflow für `pgchangefeed-kotlin` (GitHub Packages, Cloudsmith), Trigger `sdk-kotlin-v*`-Tag | kein Gate, [`harness/targets/workflows.md`](targets/workflows.md), [`ADR-0109`](../docs/plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Festlegung 2/5, [`ADR-0123`](../docs/plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) · seit slice-sdk-kotlin-publish-workflow, erweitert seit slice-sdk-kotlin-cloudsmith |
| `make suchlauf-nachmessen` | misst die `suchlauf`-Blöcke eines Slice-Plans nach; netzlos | kein Gate, [`harness/sensors/suchlauf-nachmessen.md`](sensors/suchlauf-nachmessen.md), [`ADR-0083`](../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md) · seit slice-harness-suchlauf-nachmessen |
| `make test-suchlauf-nachmessen` | Tabellentest gegen `tools/harness/suchlauf-nachmessen.sh`; netzlos | kein Gate, [`harness/sensors/suchlauf-nachmessen.md`](sensors/suchlauf-nachmessen.md) · seit slice-harness-suchlauf-nachmessen |
| `make kommentar-kennungen` | listet Kommentarblöcke, die ihre Herkunft nicht als ein auflösbares Feld tragen; Docker-only | kein Gate, [`harness/sensors/kommentar-kennungen.md`](sensors/kommentar-kennungen.md), [`ADR-0083`](../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md), [`AGENTS.md`](../AGENTS.md) §3.7 · seit slice-code-kommentare-kennungen, Nicht-Go-Formen seit slice-kommentar-kennungen-skripte |
| `make test-kommentar-kennungen` | Tabellentests des Programms `tools/harness/kommentar-kennungen/` und seines Aufrufers; Docker-only | kein Gate, [`harness/sensors/kommentar-kennungen.md`](sensors/kommentar-kennungen.md) · seit slice-code-kommentare-kennungen |
| `make fmt-check` | meldet jede Go-Datei, die `gofmt -l` nicht als formatiert führt; Docker-only | kein Gate, [`harness/sensors/fmt-check.md`](sensors/fmt-check.md), `BEO-PGC/formatierungs-drift-ohne-gate` · seit slice-harness-fmt-check |
| `make test-pin-stale-all` | Tabellentest gegen `tools/harness/pin-stale-all.sh`; netzlos | kein Gate, [`harness/sensors/pin-stale-all.md`](sensors/pin-stale-all.md) · seit slice-pin-stale-alle-digest-pins |
| `make test-fmt-check` | Tabellentest gegen `tools/harness/fmt-check.sh`; Docker-only | kein Gate, [`harness/sensors/fmt-check.md`](sensors/fmt-check.md) · seit slice-harness-fmt-check |
| `make image-mutation` / `make image-mutation-rm` | baut bzw. entfernt ein Mutations-Image `pg-change-feed-mutation:<TAG>`, getrennt vom Lauf-Beleg-Pfad | kein Gate, [`harness/targets/image-mutation.md`](targets/image-mutation.md) · seit slice-harness-mutationsbild-und-verweigerte-aktion |
| `make test-image-mutation` | Tabellentest gegen `tools/harness/image-mutation.sh`; netzlos | kein Gate, [`harness/targets/image-mutation.md`](targets/image-mutation.md) · seit slice-harness-mutationsbild-und-verweigerte-aktion |
| `make test-command-guard` | Tabellentest gegen den PreToolUse-Guard `.claude/hooks/pretooluse-command-guard.sh`; netzlos | kein Gate (ein Wächter verhindert eine Handlung, er prüft kein Ergebnis), [`harness/targets/command-guard.md`](targets/command-guard.md), [`MR-003`](conventions/MR-003-guard-inplace-textwerkzeug.md), [`MR-004`](conventions/MR-004-guard-host-python-am-kopf.md) · seit slice-harness-guard-inplace-textwerkzeug, erweitert seit slice-harness-guard-blocked-python |
| `make schema-validate` | prüft das neutrale Schema-YAML mit d-migrate (Vorlauf vor jedem `generate`/`migrate`) | kein Gate, [`harness/targets/schema-validate.md`](targets/schema-validate.md), [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md), [`ADR-0142`](../docs/plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md) |
| `make bench` | Performance-Benchmarks: vier Bench-Skripte, drei mit Schwelle, eines ohne | kein Gate, [`harness/targets/bench.md`](targets/bench.md), [ADR-0054](../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b), [ADR-0104](../docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md) |
| `make schema-rollout` | rollt das neutrale Schema mit d-migrate aus, ohne Bind-Mount des Arbeitsbaums; braucht DB-Zugang | kein Gate, [`harness/targets/schema-rollout.md`](targets/schema-rollout.md), [`ADR-0043`](../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md), [`ADR-0114`](../docs/plan/adr/0114-schema-rollout-vorlauf-view-signatur.md), [`ADR-0142`](../docs/plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md) |
| `make examples-csharp` | baut die vier Werkzeugketten-Images der C#-Beispiele und fährt darin ihre Tests; braucht Netz | kein Gate, [`harness/targets/examples.md`](targets/examples.md), [`ADR-0087`](../docs/plan/adr/0087-beispiel-clients-csharp-kotlin.md), [`ADR-0090`](../docs/plan/adr/0090-beispiel-clients-volle-matrix.md) · seit slice-098, erweitert seit slice-100, seit slice-101, seit slice-102 |
| `make examples-kotlin` | baut die vier Werkzeugketten-Images der Kotlin-Beispiele und fährt darin ihre Tests; braucht Netz | kein Gate, [`harness/targets/examples.md`](targets/examples.md), [`ADR-0087`](../docs/plan/adr/0087-beispiel-clients-csharp-kotlin.md), [`ADR-0090`](../docs/plan/adr/0090-beispiel-clients-volle-matrix.md) · seit slice-099, erweitert seit slice-100, seit slice-101, seit slice-103 |
| `make example-run-go` | baut bei Bedarf und startet eines der vier Go-Beispiele (Pflicht-Argument `SURFACE`) | kein Gate, [`harness/targets/examples.md`](targets/examples.md), [`ADR-0098`](../docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) |
| `make example-run-csharp`/`make example-run-kotlin` | startet eines der vier bereits gebauten C#- bzw. Kotlin-Beispiele (Pflicht-Argument `SURFACE`) | kein Gate, [`harness/targets/examples.md`](targets/examples.md), [`ADR-0098`](../docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) · seit slice-beispiele-csharp-kotlin-start-target |
| `make example-demo-up`/`make example-demo-down` | fährt die Demo-/Quickstart-Umgebung unter `examples/` hoch bzw. räumt sie ab | kein Gate, [`harness/targets/examples.md`](targets/examples.md), [`ADR-0098`](../docs/plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) · seit slice-beispiele-compose-bootstrap |
| `make example-transformation-demo` | zeigt die Wirkung einer Transformationsregel live an der laufenden Demo-Umgebung | kein Gate, [`harness/targets/examples.md`](targets/examples.md), [`ADR-0112`](../docs/plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) |
| `make test-sdk-public-doc-check` | fährt den Tabellentest gegen `tools/harness/sdk-public-doc-check.sh` | kein Gate, [`ADR-0134`](../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) · [`harness/sensors/sdk-public-doc-check.md`](sensors/sdk-public-doc-check.md) · seit slice-sdk-readme-nutzerdoku |
| `make test-ausgabe-kennungen-check` | fährt den Tabellentest gegen `tools/harness/ausgabe-kennungen-check.sh`; netzlos | kein Gate, [`ADR-0144`](../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 10 · [`harness/sensors/ausgabe-kennungen-check.md`](sensors/ausgabe-kennungen-check.md) · seit slice-meldungscodes-kennungsfreie-ausgaben |
| `make test-meldungscodes-check` | fährt den Tabellentest gegen `tools/harness/meldungscodes-check.sh`; netzlos | kein Gate, [`ADR-0144`](../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 5 · [`harness/sensors/meldungscodes-check.md`](sensors/meldungscodes-check.md) · seit slice-meldungscodes-registry-fehlerkopf |
| `make test-handbuch-public-doc-check` | fährt den Tabellentest gegen `tools/harness/handbuch-public-doc-check.sh`; netzlos | kein Gate, [`ADR-0143`](../docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) · [`harness/sensors/handbuch-public-doc-check.md`](sensors/handbuch-public-doc-check.md) · seit slice-handbuch-public-doc-check-gate-und-skill |
| `make sdk-pack-csharp` | baut, testet und paketiert das C#-SDK-Package `PgChangeFeed.Client`; braucht Netz | kein Gate, [`harness/targets/sdk-pack.md`](targets/sdk-pack.md), [`ADR-0106`](../docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md) Festlegung 4 · seit slice-sdk-csharp-pack-werkzeug |
| `make sdk-pack-python` | baut, testet und paketiert das Python-SDK-Package `pgchangefeed`; braucht Netz | kein Gate, [`harness/targets/sdk-pack.md`](targets/sdk-pack.md), [`ADR-0107`](../docs/plan/adr/0107-python-pypi-zweites-sdk-package.md) Festlegung 5, [`ADR-0108`](../docs/plan/adr/0108-python-sdk-uv-statt-build-twine.md) §Entscheidung Festlegung 1/3, [`ADR-0110`](../docs/plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) · seit slice-sdk-python-pack-werkzeug |
| `make sdk-pack-kotlin` | baut, testet und paketiert das Kotlin-SDK-Package `pgchangefeed-kotlin`; braucht Netz | kein Gate, [`harness/targets/sdk-pack.md`](targets/sdk-pack.md), [`ADR-0109`](../docs/plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Festlegung 5, [`ADR-0123`](../docs/plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) (Probe) · seit slice-sdk-kotlin-pack-werkzeug, erweitert seit slice-sdk-kotlin-cloudsmith |
| `make test-sdk-kotlin-integration` | Realserver-Integrationstest der Kotlin-SDK-Zustellweg-Flächen gegen die `compose.yaml`-Umgebung | kein Gate, [`harness/targets/sdk-integration.md`](targets/sdk-integration.md), slice-sdk-kotlin-reale2e · seit slice-sdk-kotlin-reale2e, erweitert seit slice-sdk-regel-realserver-e2e und slice-routing-sdk-realserver-e2e und slice-sdk-sse-client-schema-table-filter-realserver und slice-sdk-sse-filter-phase-verbindung-haertung und slice-sdk-meldungscodes-in-fehlertypen und slice-sdk-tls-optionen |
| `make test-sdk-csharp-integration` | Realserver-Integrationstest der C#-SDK-Zustellweg-Flächen gegen die `compose.yaml`-Umgebung | kein Gate, [`harness/targets/sdk-integration.md`](targets/sdk-integration.md), slice-sdk-csharp-reale2e · seit slice-sdk-csharp-reale2e, erweitert seit slice-sdk-regel-realserver-e2e und slice-routing-sdk-realserver-e2e und slice-sdk-sse-client-schema-table-filter-realserver und slice-sdk-sse-filter-phase-verbindung-haertung und slice-sdk-meldungscodes-in-fehlertypen und slice-sdk-tls-optionen |
| `make test-sdk-python-integration` | Realserver-Integrationstest der Python-SDK-Zustellweg-Flächen gegen die `compose.yaml`-Umgebung | kein Gate, [`harness/targets/sdk-integration.md`](targets/sdk-integration.md), [`ADR-0110`](../docs/plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung Festlegung 2 · seit slice-sdk-python-grpc-client-flaeche, erweitert seit slice-sdk-regel-realserver-e2e und slice-routing-sdk-realserver-e2e und slice-sdk-sse-client-schema-table-filter-realserver und slice-sdk-sse-filter-phase-verbindung-haertung und slice-sdk-meldungscodes-in-fehlertypen und slice-sdk-tls-optionen |
| `make test-sdk-kompat` | misst die Kompatibilität der SDK-Packages des Arbeitsstands gegenüber 0.5.0; braucht Netz | kein Gate, [`harness/targets/sdk-kompat.md`](targets/sdk-kompat.md), [`ADR-0145`](../docs/plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) · seit slice-sdk-0-6-kompatibilitaet-messen |
| `make test-sdk-altserver` | fährt die Fehlerfälle der SDKs gegen einen Server vor 0.6.0; braucht DB-Zugang und Netz | kein Gate, [`harness/targets/sdk-altserver.md`](targets/sdk-altserver.md), [`ADR-0145`](../docs/plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md), [`ADR-0146`](../docs/plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) (Tag mit Digest), [`ADR-0148`](../docs/plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) (Phase U) · seit slice-sdk-0-6-kompatibilitaet-messen, erweitert seit slice-upgrade-versionswechsel-alt-image |
| `make <mover>` | bewegt <…>, prüft nichts | kein Gate |
| `make <messung>` | misst <…> gegen <Schwelle> | kein Gate, ADR-<NNNN> |
| `make <vorschau>` | sagt, was <schreibender Lauf> täte; Ausgänge und Sperren in der verlinkten Datei | kein Gate |

**Aktueller Lauf-Status:** [![ci](https://github.com/pt9912/pg-change-feed/actions/workflows/ci.yml/badge.svg)](https://github.com/pt9912/pg-change-feed/actions/workflows/ci.yml)
bzw. lokal `make help` / `make gates`. Der Workflow
[`.github/workflows/ci.yml`](../.github/workflows/ci.yml) automatisiert
diesen Gate-Lauf (plus `make test`) auf jeden Pull Request und Push — kein
neues Gate, nur die Automatisierung des bestehenden
([`ADR-0051`](../docs/plan/adr/0051-cicd-pipeline-github-actions.md)).
Unmittelbar nach dem Checkout, vor dem `Gates`-Schritt, prüft ein eigener
Workflow-Schritt `uname -s`/`go env GOOS` gegen `Linux`/`linux` mit
sichtbarem Fehlschlag bei Abweichung — der sichtbare Log-Beleg für
[`LH-QA-POR-002`](../spec/lastenheft.md)
([`ADR-0058`](../docs/plan/adr/0058-testansatz-fuenf-luecken.md)
Entscheidung 5, kein Gate).
**Rote Gates:** Begründung im verlinkten `CO-<NNN>` (siehe Bindung-Spalte), Modul 7.

<!-- Domänenspezifische Gates ergänzen, je nach Repo-Klasse: -->

## Traceability rules

- PRs/Commits **müssen** mindestens eine `<LH-*>` oder `ADR-*`-ID nennen.
- Mechanisch getragen (Standing-Gate, [`ADR-0045`](../docs/plan/adr/0045-commit-traceability-standing-gate.md)): `make commit-traceability` prüft je Message der letzten 5 Commits (`RANGE=base..head` überschreibt) beide Grenzen — positive Hälfte via d-check Modul `commits` (Befund `commit-untraceable`), Betreff-Grenze via `tools/harness/commit-traceability.sh` · seit slice-006 (ausgelöst durch das dritte Auftreten der Verstoß-Klasse, review-slice-006 F-2).
- Optional, **nicht-durchsetzend** ([`ADR-0062`](../docs/plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md), Zusage einseitig nach [`ADR-0069`](../docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md)): der lokale `commit-msg`-Hook unter `.githooks/commit-msg` meldet vorab in den Klassen, in denen er liest — er fängt **nicht** jeden Verstoß, den `make commit-traceability` fängt, und weist keinen Commit zurück, den das Gate zulässt. Aktivierung ist ein einmaliger, lokaler Schritt — `git config core.hooksPath .githooks` (`core.hooksPath` ist lokale Konfiguration, sie reist nicht mit dem Klon; `git commit --no-verify` umgeht den Hook). Er ersetzt `make commit-traceability` nicht: er läuft ausschließlich lokal vor `git commit` und ist kein Gate-Ziel · seit slice-073.
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
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`) — Exit-Code direkt prüfen, nie durch eine Pipe/einen Wrapper hindurch (`AGENTS.md` §3.9).
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
