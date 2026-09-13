# ADR-0051: CI/CD-Pipeline über GitHub Actions

**Status:** Accepted

**Datum:** 2026-09-13

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-13)

**Bezug:** [`ADR-0044`](0044-image-beleg-semantik.md) (Digest-Beleg-Semantik
des lokalen `make image`-Laufs — diese ADR ergänzt sie um
Registry-Publishing, siehe Entscheidung 2; kein Widerspruch, siehe
Abgrenzung dort), [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(Commit-Traceability-Standing-Gate — bindet die Form der
Dependabot-Commit-Messages, siehe Entscheidung 8),
[`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md) (a-check-Image
als einer der Pin-Inventar-Einträge in Entscheidung 5),
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (d-migrate-Image
ebenso Pin-Inventar-Eintrag in Entscheidung 5)

**Schärft:** — *(Prozess-ADR ohne Spec-Stratum; die getragenen Regeln stehen
in [`AGENTS.md`](../../../AGENTS.md) §3.8 (Action-Pinning) und in den noch
zu schreibenden `.github/workflows/*.yml` — Folge-Slice `slice-039` <!-- d-check:status-provenance -->
und Folge-Slice `slice-040` <!-- d-check:status-provenance -->, Planning-Stratum.
Wer diese ADR ändert, zieht §3.8 und die Workflow-Dateien nach.)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

pg-change-feed hat aktuell **kein** `.github/`-Verzeichnis, keine
Git-Tags, keine `version.md` und keine Registry-Veröffentlichung. Es gibt
nur einen rein lokalen `make image`-Build
(`ghcr.io/pt9912/pg-change-feed:dev`, Digest-Beleg in
`harness/image-hash.txt`, Semantik in [`ADR-0044`](0044-image-beleg-semantik.md))
und lokale Gates (`make gates`: `baseline-verify`, `docs-check`, `a-check`,
`commit-traceability` — [`harness/README.md`](../../../harness/README.md)
§Sensors). `harness/README.md` nennt bereits, unausgesprochen bis zur
Umsetzung, ein advisory `make image-cve` als „nicht behauptet (geplant)"
mit eingetretener Aktivierungsbedingung (erster grüner `make image`-Lauf,
seit slice-001 <!-- d-check:status-provenance -->) — diese ADR löst genau
dieses Vorhaben ein, zusammen mit
vier weiteren Teilentscheidungen, die der Auftraggeber (pt9912) explizit
angefordert hat: eine GitHub-Actions-Pipeline analog zur Fünf-Workflow-Form
in `d-check` (`ci`, `release`, `image-scan`, `upstream-drift`,
`hub-description`, plus `dependabot.yml`), mit vollem Umfang — **beide**
Registries (GHCR und Docker Hub) als Release-Ziel.

Die Musterquelle (`/Development/d-check/.github/`, gelesen, nicht kopiert)
zeigt eine bewährte Form: Action-Pinning per SHA mit Tag-Kommentar,
`permissions: {}` auf Workflow-Ebene mit gezielter Lockerung je Job,
fail-open bei Nachtläufen (Netz-/Werkzeugausfall → Skip statt Rot; ein
tatsächlich gefundener Drift bleibt sichtbar), advisory statt Gate für
CVE-Scan und Pin-Freshness. Diese Form ist Vorbild für **Stil und
Disziplin**, nicht zum 1:1-Kopieren — d-check hat andere Pins (eigener
Go-/Lint-Stand), eine andere Registry-Historie (`ADR-0002`/`ADR-0014` dort)
und ein anderes Gate-Set (`make ci` dort existiert in diesem Repo nicht;
hier heißt das Äquivalent `make gates`).

Fünf Teilentscheidungen sind zu treffen, bevor Implementierung beginnt
(Architect vor Code, Modul 8):

1. Versionierungsschema (`version.md`-Ort, Git-Tag-Form).
2. Registry-Ziele (GHCR **und** Docker Hub, `:latest`-Semantik).
3. Action-Pinning (neue Hard Rule).
4. CVE-Scan (Trivy, advisory).
5. Upstream-Pin-Freshness (fail-open, gegen das **tatsächliche**
   Pin-Inventar dieses Repos).

**Das tatsächliche Pin-Inventar** (ermittelt, nicht übernommen von
d-check):

| # | Pin | Fundort | Achse(n) |
|---|---|---|---|
| P1 | `golang:1.27-alpine@sha256:cf6fca66…` | `Dockerfile:19` (deps-Stage), `Makefile` `TOOLCHAIN_IMAGE` (identisch) | Digest-Drift (Tag fix) |
| P2 | `gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872…` | `Dockerfile:29` (runtime-Stage) | Digest-Drift (Tag fix) |
| P3 | `golang:1.27@sha256:b475798f…` | `Makefile` `TOOLCHAIN_RACE_IMAGE` (Debian-Variante für `-race`) | Digest-Drift (Tag fix) |
| P4 | `postgres:18-alpine@sha256:63bdc97d…` | `Makefile` `PG_TEST_IMAGE` | Digest-Drift (Tag fix) |
| P5 | `ghcr.io/pt9912/d-migrate@sha256:862dfb04…` | `Makefile` `D_MIGRATE_IMAGE` | Digest-Drift (Tag fix), [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) |
| P6 | `ghcr.io/pt9912/a-check@sha256:34d3dfb5…` | `a-check.mk` `A_CHECK_IMAGE` | Digest-Drift (Tag fix), [`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md) |
| P7 | `ghcr.io/pt9912/d-check:v0.75.0@sha256:18e9cd85…` | `d-check.mk` `DCHECK_IMAGE`/`DCHECK_DIGEST` | **zwei** Achsen — Tag-Frische (neuerer `d-check`-Release?) **und** Digest-Drift (derselbe Tag, anderer Bau?) |
| P8 | Kurs-Baseline `v6.5.0` | `harness/conventions.md` §Baseline | Tag-Frische (neuere Kurs-Welle?) |
| P9 | GitHub-Action-Pins (neu, aus Entscheidung 3) | `.github/workflows/*.yml` | Tag-Frische je SHA |

P1 und P2 (Dockerfile-`FROM`) sind bereits **teilweise** durch
`make image-stale` abgedeckt (advisory, braucht Netz,
[`harness/README.md`](../../../harness/README.md) §Werkzeuge). P3–P9 haben
**keine** Frische-Beobachtung — diese ADR schließt die Lücke für das
**gesamte** Inventar, nicht nur für den Dockerfile-Ausschnitt.

**Konsistenz-Prüfung gegen bestehende Hard Rules/ADRs** (vor der
Entscheidung, nicht danach):

- **AGENTS §3.1 (Docker-only):** Jeder Workflow-Schritt ruft ein
  bestehendes `make <target>` auf; kein Workflow enthält Inline-Shell-Logik,
  die ein Gate umgeht oder dupliziert. Wo ein Schritt reine Workflow-Mechanik
  ist (Tag-Parsing, Registry-Login, Digest-Vergleich zwischen zwei
  Registries), bleibt das Workflow-Skript — es ersetzt kein `make`-Target,
  es orchestriert sie.
- **AGENTS §3.6 (Gates dürfen nicht ohne ADR gelockert werden):** Diese ADR
  **lockert kein bestehendes Gate**. `make gates` bleibt inhaltlich
  unverändert; der CI-Workflow automatisiert nur seinen Aufruf. CVE-Scan und
  Pin-Freshness werden **advisory** eingeführt (kein Gate), nicht als Ersatz
  für ein bestehendes Gate.
- **[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (d-migrate):**
  berührt diese ADR nur über das Pin-Inventar (P5); die Rollout-Semantik
  (Pflicht-Report, Rollback-Artefakt) ist unverändert und nicht Gegenstand
  dieser Entscheidung.
- **[`ADR-0044`](0044-image-beleg-semantik.md) (Image-Beleg-Semantik):**
  behandelt den **lokalen** `make image`-Digest (`ghcr.io/pt9912/pg-change-feed:dev`,
  builder- und lauf-gebunden, kein Inhalts-Fingerabdruck über Läufe/Umgebungen
  hinweg). Diese ADR behandelt das **Registry-Publishing** eines mit
  demselben `make image`-Mechanismus, aber mit Versions-Tag gebauten Images.
  Kein Widerspruch: ADR-0044s Aussage („Digest ist Lauf-Beleg, kein
  Inhalts-Fingerabdruck über Umgebungen hinweg") gilt für **jeden**
  `make image`-Lauf unverändert, auch den Release-Lauf; diese ADR fügt
  hinzu, **wohin** das gebaute Image zusätzlich zum lokalen Tag gepusht
  wird und **welcher Tag-Name** es dort trägt. Sie ändert nichts an der
  Digest-Semantik von `harness/image-hash.txt`.

## Entscheidung

Wir wählen die folgende Pipeline-Architektur:

1. **Versionierung:** `docs/user/version.md` trägt die aktuelle Version als
   einzeiligen SemVer-String (`MAJOR.MINOR.PATCH[-PRERELEASE]`, ohne
   führendes `v`). Git-Tags der Form `v<SemVer>` sind der Release-Auslöser.
   Der Release-Workflow validiert den Tag strikt gegen die SemVer-2.0-Form
   (fail-fast vor Login/Build/Push, wie im d-check-Vorbild) **und** prüft
   `docs/user/version.md` am getaggten Commit gegen die Tag-Version — bei
   Abweichung bricht der Release-Lauf ab. `version.md` wird als eigener,
   traceability-pflichtiger Commit **vor** dem Tag aktualisiert; der Tag
   selbst bleibt die verbindliche Quelle der tatsächlich veröffentlichten
   Versionsnummer, `version.md` der menschlich lesbare Spiegel für den
   Repo-Zustand ohne `git tag`-Zugriff.
2. **CI-Workflow (`ci.yml`):** läuft auf jeden Pull Request und Push
   (`tags-ignore: ['**']`, damit Tag-Pushes ausschließlich `release.yml`
   auslösen). Ruft `make gates` (Docker-only, netzlos) auf. Ergänzend, wo
   sinnvoll und ohne ein Gate zu duplizieren, `make test` (netzlos,
   Race-Detector) — die DB-gebundenen Testziele (`test-store`,
   `test-replication`, `test-integration`) bleiben kein Gate und laufen
   nach Bedarf, nicht zwingend bei jedem PR (Umfang: Folge-Slice).
3. **Release-Workflow (`release.yml`):** Trigger `push: tags: ['v*']`.
   Baut über `make image` (Docker-only, AGENTS §3.1 — kein separater
   `docker build`-Aufruf außerhalb `make`) mit der validierten Version;
   `make image` braucht dafür einen Versions-/Tag-Parameter zusätzlich zum
   bisherigen festen `:dev`-Tag — diese Erweiterung ist Implementierungsumfang
   der Folge-Slices, nicht dieser ADR. Push nach **GHCR und Docker Hub**
   (Entscheidung 2), danach GitHub-Release mit Digest-Pin im Beschreibungstext
   (Vorbild: d-check `release.yml`).
4. **Registry-Ziele und `:latest`-Semantik:** siehe Entscheidung 2 unten.
5. **Action-Pinning:** neue Hard Rule `AGENTS.md` §3.8 — siehe unten,
   direkt ergänzt als Teil dieses Architect-Zugs.
6. **CVE-Scan (`image-scan.yml`):** Trivy gegen das publizierte
   GHCR-`:latest`-Image, nächtlich (`schedule`) + `workflow_dispatch`,
   **advisory** (kein Gate, kein Abbruch von `make gates`/`ci.yml`).
7. **Upstream-Pin-Freshness (`upstream-drift.yml`):** ein **fail-open**
   Nachtlauf (`schedule` + `workflow_dispatch`) gegen das vollständige,
   oben ermittelte Pin-Inventar (P1–P9). „Fail-open" heißt: Netz-,
   Werkzeug- oder Registry-Ausfall bei einer einzelnen Achse führt zu
   **Skip dieser Achse**, nicht zu Rot des gesamten Laufs (`if: always()`
   je Achsen-Schritt, wie im d-check-Vorbild) — ein tatsächlich gefundener
   Drift bleibt sichtbar (roter Lauf in der Actions-Übersicht), bleibt
   aber advisory und blockiert kein Gate. P1/P2 nutzen das bestehende
   `make image-stale`; P3–P9 brauchen neue, analog benannte Make-Targets
   (Namensschema und genaue Skript-Logik: Implementierungsumfang der
   Folge-Slices).
8. **`hub-description.yml`:** eigener, `workflow_dispatch`-fähiger
   Workflow, der die Docker-Hub-Beschreibung nach einem erfolgreichen
   Release synchronisiert (aufgerufen als `needs: release`-Job aus
   `release.yml`, wie im Vorbild) — Präsentation, kein Bestandteil der
   Distributions-Zusage; ein Fehlschlag macht das Release nicht rot.
9. **`dependabot.yml`:** zwei Ecosysteme, `gomod` und `github-actions`,
   wöchentlich, je mit `commit-message.prefix` `[ADR-0051]` — diese ADR ist
   der ehrliche Anker für die Commits, die daraus entstehen (Vorbild:
   d-checks eigene `ADR-0067`-Selbstreferenz), notwendig, damit
   Dependabot-Commits das Commit-Traceability-Standing-Gate
   ([`ADR-0045`](0045-commit-traceability-standing-gate.md)) erfüllen.
   Kein `docker`-Ecosystem (Digest-Hebung bleibt bewusster Commit, analog
   zur bestehenden `make image-stale`-Disziplin dieses Repos — kein
   Auto-Update von Docker-Basis-Images).

**Nicht Gegenstand dieser ADR:** die konkreten SHA-Pins der einzusetzenden
Actions, die konkrete erste Versionsnummer in `version.md`, die genauen
Make-Target-Namen für P3–P9, und der genaue Aufteilungsschnitt zwischen
`slice-039` und `slice-040` <!-- d-check:status-provenance --> — das ist
Implementierungs- bzw.
Planungsdetail, keine Architekturentscheidung.

## Direkt ergänzte Hard Rule — `AGENTS.md` §3.8

Als Teil dieses Architect-Zugs wurde `AGENTS.md` um §3.8 (Action-Pinning)
ergänzt (siehe Entscheidung 5): jede `uses:`-Zeile in
`.github/workflows/*.yml` SHA-gepinnt mit `# vX.Y.Z`-Kommentar.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra je
Teilentscheidung** — „nichts tun" ist eine davon. Eine ADR ohne
Alternativen ist ein Postulat, kein Entscheidungsprotokoll, und im Review
nicht verteidigbar (Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR)).

### 1. Versionierungsschema — Ort von `version.md`

| Option | Pro | Contra |
|---|---|---|
| A — Repo-Wurzel (`VERSION.md`) | verbreitete OSS-Konvention, an bekannter Stelle für Tools/Menschen | Repo-Wurzel führt heute nur Struktur-Dateien (`Makefile`, `Dockerfile`, `a-check.mk` …), keine Zustands-Doku; ein neuer Root-Zustands-Doc hat keinen Rang in der Source-Precedence-Tabelle ([`harness/README.md`](../../../harness/README.md) §Source precedence) und müsste dort neu eingeordnet werden |
| B — `docs/plan/planning/version.md` | räumliche Nähe zu Roadmap/Wellen (zeitliche Schicht) | `version.md` ist kein Slice/keine Welle und unterliegt keiner Lifecycle-Zustandsmaschine (Modul 5); in `docs/plan/planning/` liegend, würde sie versehentlich wie ein Lifecycle-Artefakt behandelt (z. B. von Tools, die dieses Verzeichnis nach `slice-*`/`welle-*`-Mustern durchsuchen) |
| **C — `docs/user/*` (gewählt)** | Verzeichnis existiert bereits (`benutzerhandbuch.md`, `benutzerhandbuch-standard.md`) und ist als Rang 6 in der Source-Precedence-Tabelle geführt: „Operations, Quality, Releasing"; Versionsstand ist exakt operative Metadaten, passt zur bestehenden Rolle des Verzeichnisses; deckt sich mit dem d-check-Vorbild (`docs/user/releasing.md`) | weitet die bisher rein prosaische Rolle von `docs/user/` erstmals um eine maschinenlesbare Ein-Zeilen-Zustandsdatei; keine inhaltliche Härte (kein neuer Rang nötig — Rang 6 trägt es bereits) |

### 2. Registry-Ziele und `:latest`-Semantik

| Option | Pro | Contra |
|---|---|---|
| A — nur GHCR, `:latest` bei jedem Tag | einfachster Pfad, ein Secret-Paar weniger | widerspricht dem expliziten Nutzerwunsch (beide Registries); `:latest` bei Prerelease-Tags zeigt Konsumenten instabilen Code |
| B — nur Docker Hub | breitere Standard-Reichweite für Docker-Nutzer ohne GHCR-Login | kein GHCR-Bezug, der GitHub-Release-Beleg (Digest-Pin) verliert seinen natürlichen Registry-Anker; widerspricht dem expliziten Nutzerwunsch |
| C — beide Registries, aber je ein **eigener Build** je Ziel | organisatorisch trivial (zwei unabhängige Push-Jobs) | verdoppelt die Baudauer und riskiert Inhalts-Drift zwischen den Registries (zwei Builds desselben Commits sind nicht zwingend byte-identisch, [`ADR-0044`](0044-image-beleg-semantik.md)); widerspricht dem Docker-only-Prinzip eines einzigen `make image`-Laufs (AGENTS §3.1) |
| **D — beide Registries, ein Build, Content-Mirror; `:latest` nur bei stabilen (nicht-Prerelease) Tags, auf beiden (gewählt)** | ein `make image`-Lauf, Ergebnis wird auf beide Registries getaggt/gepusht (kein zweiter Build); Inhalts-Gleichheit über Config-Digest verifizierbar (Manifest-Digest ist registry-lokal, Config-Digest nicht — d-check-Vorbild); `:latest` zeigt nie auf einen instabilen Stand | ein Job mehr im Release-Workflow (Docker-Hub-Login + Push + Gleichheits-Check); zwei Secret-Paare (GHCR über `GITHUB_TOKEN`, Docker Hub über `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`) statt eines |

### 3. Action-Pinning

| Option | Pro | Contra |
|---|---|---|
| A — floatende Major-Tags (`uses: actions/checkout@v7`) | am wenigsten Pflege, Standard-GitHub-Beispielform | Tag ist beweglich — sein Ziel-Commit kann sich ändern, ohne dass die Workflow-Datei sich ändert; ein Workflow-Schritt läuft mit den Zugangsdaten dieses Repos (`GITHUB_TOKEN`, Registry-Secrets), ein umgebogener Tag ist ein stiller Angriffsvektor ohne eigenen Diff |
| B — SHA-Pinning ohne Kommentar | schließt den Angriffsvektor von A | ein nackter SHA ist für Menschen nicht auf einen Blick versionierbar; Audit und Hebung (welche Version ist das?) verlangen manuelles Auflösen des SHA gegen die Releases der Action |
| **C — SHA-Pinning mit `# vX.Y.Z`-Tag-Kommentar (gewählt)** | schließt denselben Angriffsvektor wie B; der Kommentar hält die Pinnung lesbar und audit-/hebbar ohne SHA-Auflösung; etabliertes Muster (d-check-/u-boot-Konvention) | ein Kommentar kann veralten, wenn der SHA gehoben, der Kommentar aber vergessen wird — bleibt menschliche Disziplin, mangels eines heute existierenden Sensors (Grenze, benannt) |
| D — Dependabot-only, kein statisches Pinning-Gebot | Dependabot hebt automatisch, kein manueller Pflegeaufwand für die Kommentar-Form | verhindert nicht, dass ein ungepinnter `uses:`-Eintrag je landet — kein Enforcement gegen den Erstfall, nur gegen das Altern eines bereits gepinnten Standes; löst A/B nicht, ergänzt nur C |

### 4. CVE-Scan

| Option | Pro | Contra |
|---|---|---|
| A — kein CVE-Scan | kein zusätzlicher Workflow, kein Netzbedarf | blind gegenüber CVEs, die ohne Commit in diesem Repo auftauchen (Basis-Images altern zwischen Releases); widerspricht der in `harness/README.md` bereits angekündigten Aktivierungsbedingung |
| B — CVE-Scan als blockierendes Gate in `make gates`/`ci.yml` | maximale Durchsetzung, kein PR mit bekannten CVEs mergbar | verletzt die Netzlos-Eigenschaft der bestehenden Gates (`baseline-verify`, `docs-check`, `a-check`, `commit-traceability` sind alle netzlos, [`harness/README.md`](../../../harness/README.md) §Sensors); ein extern getriebener, jederzeit veränderlicher CVE-Feed als PR-Blocker widerspricht AGENTS §3.6 im Ergebnis (eine Gate-Schwelle, die sich ändert, ohne dass der Code sich ändert) |
| **C — Trivy gegen das publizierte GHCR-`:latest`-Image, nächtlich + `workflow_dispatch`, advisory (gewählt)** | deckt sich mit der bereits in `harness/README.md` angekündigten, aber unimplementierten Zusage; prüft den tatsächlich veröffentlichten Stand statt eines lokalen Dev-Builds; kein Gate, kein PR-Blocker | Befunde werden nicht automatisch behoben — Hebung bleibt ein bewusster Akt (Dockerfile-Basis-Update), das Scannen allein schließt keine Lücke |
| D — Scan bei jedem PR gegen das lokal gebaute Dev-Image | frühestmögliches Signal | das Dev-Image ist nie veröffentlicht — ein Scan-Ergebnis dagegen sagt nichts über das aus, was tatsächlich ausgeliefert wird; verletzt zusätzlich die Netzlos-Eigenschaft der PR-Gates (Trivy-DB-Pull) |

### 5. Upstream-Pin-Freshness

| Option | Pro | Contra |
|---|---|---|
| A — keine Pin-Freshness-Beobachtung | kein neuer Workflow | das Repo führt neun Pin-Achsen (P1–P9); ohne Beobachtung altert es unbemerkt, bis ein CVE oder ein Migrationsbedarf es erzwingt |
| B — fail-**closed** nächtlicher Workflow (jede Abweichung, auch Netz-/Werkzeugausfall, färbt rot) | maximale Sichtbarkeit, keine Sonderfälle in der Logik | ein Nachtlauf, der wegen transienter Registry-/Netzausfälle regelmäßig rot wird, verliert sein Signal (Alarmmüdigkeit) — genau das Muster, das die eigene Diagnose in `d-check/upstream-drift.yml` benennt: „ein dauerroter Nachtlauf ist derselbe verwaiste Sensor" |
| C — nur `make image-stale` erweitern (Dockerfile-`FROM`, P1/P2) | kleinster Aufwand, nutzt Bestehendes vollständig | deckt nur 2 von 9 Pin-Achsen ab; P3–P9 (Race-Toolchain, PG-Testcontainer, d-migrate, a-check, d-check selbst, Kurs-Baseline, Action-Pins) blieben unbeobachtet — widerspricht der Anforderung, das **tatsächliche** Inventar dieses Repos abzudecken |
| **D — ein Workflow, alle neun Achsen, fail-open bei Werkzeug-/Netzausfall (Skip statt Rot), gefundener Drift bleibt sichtbar (gewählt)** | deckt das vollständige Inventar; unterscheidet „Achse konnte nicht geprüft werden" von „Achse zeigt Drift" — nur Letzteres ist ein Befund; bewährtes Muster (d-check-Vorbild, dort seit mehreren Wellen im Betrieb) | mehr Wartungsaufwand (neun Einzel-Achsen statt einer); die genauen Make-Targets für P3–P9 existieren noch nicht und sind Implementierungsumfang der Folge-Slices |

## Konsequenzen

- Positiv: Die Pipeline-Architektur ist vor Implementierung entschieden und
  gegen bestehende Hard Rules/ADRs geprüft (§3.1, §3.6, ADR-0043, ADR-0044).
  `version.md` bekommt einen begründeten, in die Source-Precedence
  eingeordneten Ort. Beide gewünschten Registries sind mit einer
  Inhalts-Konsistenzgarantie (Config-Digest-Vergleich) statt eines zweiten,
  potenziell driftenden Builds abgedeckt. Die Action-Pinning-Regel schließt
  einen realen Supply-Chain-Angriffsvektor, bevor der erste Workflow
  entsteht. CVE-Scan und Pin-Freshness lösen eine bereits in
  `harness/README.md` angekündigte, aber unimplementierte Zusage ein und
  decken — anders als das d-check-Vorbild unreflektiert übernommen hätte —
  das **tatsächliche** Neun-Achsen-Pin-Inventar dieses Repos.
- Negativ: Kein bestehendes Gate wird härter oder weicher; die neuen
  advisory-Läufe (CVE-Scan, Pin-Freshness) brauchen Netz und liefern damit
  keine deterministische Grün/Rot-Aussage wie die netzlosen Gates. Zwei
  Registry-Secret-Paare statt eines erhöhen die Angriffsfläche der
  Repository-Secrets geringfügig. Die genauen Make-Targets für P3–P9 und
  der `make image`-Versionsparameter existieren noch nicht — bis zur
  Umsetzung ist der Beleg für „vollständige Pin-Abdeckung" diese ADR, kein
  laufender Sensor.
- Folgepflicht: **Zwei Folge-Slices** setzen diese Entscheidung um
  (`slice-039`, `slice-040` <!-- d-check:status-provenance -->,
  Planning-Stratum — die genaue Aufteilung ist
  Planungsentscheidung, nicht Teil dieser ADR): `docs/user/version.md`
  anlegen, `.github/workflows/{ci,release,image-scan,upstream-drift,hub-description}.yml`
  und `.github/dependabot.yml` schreiben, `make image` um einen
  Versions-/Tag-Parameter erweitern, die Make-Targets für P3–P9
  ergänzen, `harness/README.md` §Sensors/§Werkzeuge um die neuen
  (advisory) Targets nachziehen, sobald sie existieren (kein
  Halluzinieren nicht existierender Targets, Modul 13).

## Fitness Function (falls maschinell prüfbar)

Heute **keine** maschinell geprüfte Regel für diese ADR selbst — sie ist
eine Architektur-/Prozessentscheidung, kein Code-Constraint. Nach Umsetzung
der Folge-Slices gilt:

| Tooling | Regel | Make-Target |
|---|---|---|
| `make gates` (unverändert) | bestehende Gates laufen im CI-Workflow, keine neue/gelockerte Schwelle | `make gates` |
| Trivy | CRITICAL/HIGH gegen das publizierte GHCR-`:latest`-Image, advisory | `make image-cve` (noch nicht implementiert — siehe Folgepflicht) |
| Pin-Freshness-Skripte | je Pin-Achse P1–P9 Tag-/Digest-Vergleich gegen Upstream, fail-open | `make image-stale` (P1/P2, existiert) + neue Targets für P3–P9 (noch nicht implementiert) |

Bis zur Implementierung ist die Action-Pinning-Regel (§3.8) **textlich
verkörpert, nicht maschinell bewacht** — wie bei ADR-0044s Digest-Semantik.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die
Entscheidung unbefristet weiter, auch wenn ihre Voraussetzung weg ist
(Baseline-Regelwerk `modul-04-adrs.md` §Kernidee (Modul 4)).

| Teilentscheidung | Trigger |
|---|---|
| 1 — Versionierung, Ort `docs/user/` | `docs/user/` wird als Konvention aufgegeben oder umbenannt (MR-Eintrag in `harness/conventions.md`) → Folge-ADR mit `supersedes` |
| 2 — Registry-Ziele | Docker Hub wird als Ziel nicht mehr gewünscht oder das Konto entfällt → Folge-ADR mit `supersedes` (nicht: stilles Entfernen des Push-Jobs) |
| 3 — Action-Pinning | GitHub Actions verändert das SHA-Referenzierungsmodell grundlegend (derzeit nicht absehbar) → Re-Evaluierung |
| 4 — CVE-Scan | dreimaliges Auftreten eines Werkzeug-Ausfalls von Trivy im Nachtlauf (Steering-Loop-Schwelle, Beobachtungs-Register) → Werkzeugwechsel prüfen |
| 5 — Pin-Freshness | ein Lauf-Zweig braucht aus einer der neun Achsen ein **blockierendes** Signal (z. B. eine Sicherheitsrichtlinie verlangt harte Durchsetzung) → Folge-ADR, das die betroffene Achse von advisory auf Gate hebt (AGENTS §3.6, Gate-Verschärfung ist ebenfalls ADR-pflichtig) |

Sonst **permanent** — die Pipeline-Architektur wächst mit dem Repo, ohne
dass die Grundentscheidungen erneut zu treffen sind.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Anlass: Auftraggeber-Anforderung (pt9912) für eine GitHub-Actions-CI/CD-Pipeline analog zu `d-check`, mit GHCR **und** Docker Hub als Release-Ziele; Architect-Lauf vor jeder Implementierung (Modul 8) | Folge-Slices `slice-039`, `slice-040` (Planning-Stratum) <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
