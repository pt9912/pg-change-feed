# AGENTS.md — Briefing für AI-Coding-Agenten

---

## 1. Was diese Datei ist

Onboarding-Briefing für jede AI-Session, die in diesem Repo Code oder
Dokumentation ändert. Sie verweist auf die kanonischen Quellen und
formuliert die Hard Rules, die der Implementer-Agent immer
einhalten muss.

Regeln dieser Datei: Baseline-Regelwerk `modul-09-implementierung.md`
§Ziel-Form: AGENTS.md — sie trägt Hard Rules und Pointer auf kanonische
Quellen, sie dupliziert deren Inhalt nicht; sonst entsteht Drift.

**Bei Konflikt zwischen dieser Datei und einer kanonischen Quelle gilt
die kanonische Quelle** (Source Precedence — siehe
`harness/README.md`).

Strukturregeln (ID-Schemata, Verzeichniskonvention, Adaptionen ggü.
Baseline, Modus-Deklarationen pro Sub-Area, Zusatzklassen für
Sensors-Bindung) leben in
[`harness/conventions.md`](harness/conventions.md).

Das **Regelwerk der adoptierten Baseline** ist die **präsente,
nachschlagbare Vertiefung** zu diesem Briefing: ein self-navigierbares
**Modul-Bundle** (`README.md` = Index). Beim Bootstrap wird das
self-contained Release-ZIP
(<https://github.com/pt9912/ai-harness-course/releases/download/v6.5.0/lab-regelwerk.zip>)
**committet vendored** unter `.harness/baseline/<tag>/{regelwerk,templates}/`
(Regelwerk *und* Templates parallel, netzlos materialisiert samt `SHA256SUMS`
— Vorgehen siehe
Baseline-Regelwerk `modul-02-harness-bootstrap.md`;
Quelle/Stand in [`harness/conventions.md`](harness/conventions.md) §Baseline).

Die verkörperte Form (dieses Briefing, die Konventionen, deine
ausgefüllten Artefakte) **führt**; das Regelwerk wird **pro Entscheidung
nachgeschlagen, deren
operative Detailtiefe das Briefing nicht trägt** — Trigger-Klassen,
Sub-Area-Qualifikation, Carveout-vs-Reconciliation, Modus-Diagnose. Dabei
**nur den benötigten Abschnitt** laden (README ist der Index), **nicht das
ganze Regelwerk im Kontext halten**. Breiterer Pflicht-Blick bleibt bei:
Bootstrap, Änderung an [`harness/conventions.md`](harness/conventions.md)
(Adaptionen `MR-<NNN>`, Source-Precedence, ID-Schema), Drift-Audit gegen die
Baseline (Baseline-Regelwerk `modul-02-harness-bootstrap.md`
§Freshness-Audit der vendored Baseline — darunter die Stichprobe gegen
den Bestand, die auch bei aktuellem Pin läuft). Derivativ: bei Konflikt gelten die kanonischen Quellen.

Die **Skelett-Vorlagen** der Baseline liegen **vendored** unter
`.harness/baseline/<tag>/templates/` (aus demselben Baseline-Bundle) und
tragen zwei Rollen: als **Referenz-Form**, auf die das Regelwerk mit
`../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu
`regelwerk/` vendored), und als Vorlage, die beim Anlegen neuer Artefakte
(ADR, Slice, Welle, Carveout, Review-Report) **kopiert und ausgefüllt** wird
statt frei zu formulieren.

## 2. Kanonische Quellen (Source Precedence)

In dieser Reihenfolge:

1. [`spec/lastenheft.md`](spec/lastenheft.md) — vertraglich abnahmebindend.
2. [`spec/pflichtenheft.md`](spec/pflichtenheft.md) — technisch verbindlich, fortschreibbar.
3. [`spec/architecture.md`](spec/architecture.md) — Komponenten- und Sequenzsicht.
4. [`docs/plan/adr/`](docs/plan/adr/) — ADR-Verzeichnis und -Index.
5. [`docs/plan/planning/in-progress/roadmap.md`](docs/plan/planning/in-progress/roadmap.md) — Wellen-Sequenz.
6. `docs/user/*` *(falls vorhanden)* — Operations, Quality, Releasing. <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
7. [`README.md`](README.md) — Projekt-Überblick.
8. **AGENTS.md (diese Datei).**
9. [`harness/README.md`](harness/README.md) — Harness-Einstieg.

## 3. Harte Regeln

### 3.1 Docker-only

Kein lokales <venv/SDK/Toolchain-Install>. Alles läuft über `make`
(das Docker nutzt). Host braucht nur Docker und GNU `make`.

**Falsch:** <z.B. `pip install ...`>
**Richtig:** <z.B. `<make-target>`>

**Begründung:** Toolchain-Reproduzierbarkeit + Supply-Chain-Defense.

### 3.2 Suppression-Verbot

Dieses Repo hat **keinen Linter** — kein `.golangci.yml`, kein
`lint`-Target in `Makefile` oder `harness/mk/*.mk`. Inline-Suppression
(`//nolint`) ist deshalb **ausnahmslos verboten**, nicht weil eine
Ausnahmeliste sie einschränkt, sondern weil es kein Werkzeug gibt, dessen
Warnung sie unterdrücken könnte — ein `//nolint` ohne Linter ist entweder
tote Dekoration oder ein Vorgriff auf ein Profil, das noch niemand
entschieden hat.

Das Coverage-Gate (`make coverage-gate`, [ADR-0054](docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md))
hat aus demselben Grund **strukturell keinen** Ausnahme-Pfad: Es gibt kein
Coverage-Pendant zu `//nolint` — die Gesamt-Coverage besteht oder scheitert
als Zahl, keine Zeile lässt sich davon ausnehmen. Ist die Schwelle für den
aktuellen Ist-Stand unerreichbar, ist die Antwort ein Carveout (Modul 7)
oder ein bootstrap-aware Gate mit eigener, niedrigerer Einstiegsstufe
([ADR-0054](docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)) —
keine stille Ausnahmeliste.

**Falsch:** `//nolint:errcheck // reicht so` irgendwo im Code.
**Richtig:** kein `//nolint` — es gibt (noch) keinen Linter, dessen Warnung
es unterdrücken könnte.

**Begründung:** Ein Suppression-Verbot ohne Objekt ist eine Lücke, keine
Regel. Führt ein späterer Slice einen Linter ein, deklariert die
einführende ADR den Ausnahme-Ort dieses Abschnitts neu (zentrale
Konfigurationsdatei mit `Why:`-Begründung je Regel) — bis dahin bleibt
Inline-Suppression ohne Ausnahme verboten.

### 3.3 git mv + Inhaltsänderung = zwei Commits

Wenn eine Datei verschoben **und** der Inhalt umgeschrieben wird, sind das
zwei Commits — der Move-Commit bleibt rein (Git erkennt R-Rename). Welcher
zuerst kommt, sagt der Vorgang:

1. Regelfall: `git mv source target` → eigener Commit, dann Inhalt umschreiben.
2. Lifecycle-Übergang nach `done/`: erst der Inhalt (DoD-Häkchen,
   Closure-Notiz), dann der reine `git mv` — die Notiz ist die Bedingung für
   `done/`, nicht ihre Folge.

**Begründung:** Sonst fällt die Rename-Detection unter die 50%-
Similarity-Schwelle und `git log --follow` wird unzuverlässig.

### 3.4 Architektur ist sprach- und meilensteinfrei

`spec/architecture.md` darf Pfade zu **Code-Modulen** referenzieren
(`src/service/`), aber **keine** Wellen, Slices, Commit-Hashes oder
Closure-Daten. Die zeitliche Schicht lebt in
`docs/plan/planning/` und den späteren Closure-Notizen. Auch **keine
ADR-Bezüge**: Die Sicht steht im Stabilitäts-Rang über der ADR; welche ADR
eine Aussage verbindlich macht, deklariert die ADR in ihrem `Schärft:`-Feld.

Diese Regel ist *verkörpert*, nicht hier entschieden — sie folgt aus dem
Sicht-Stratum (Baseline-Regelwerk `modul-03-spec.md`
§Ziel-Form: Architektur-Sicht).

### 3.5 ADRs sind nach `Accepted` immutable

Eine ADR mit Status `Accepted` wird nicht inhaltlich überschrieben.
Korrekturen entstehen als neue ADR mit `Supersedes ADR-NNNN`.

### 3.6 Gates dürfen nicht ohne ADR gelockert werden

Jede Schwellen-Senkung (Coverage, Linter-Strenge, Architekturregel)
ist ein ADR, kein PR-Kommentar.

### 3.7 Ein Kommentar beschreibt, was da ist

Gilt für Code, Konfiguration und Skripte — und für Zustandsfelder (unten).
Ein Kommentar trägt eine dieser Klassen — **Zusage · Kopplung · Abgrenzung ·
Rang-Zeiger · Grenze** — und schreibt an den, der die Stelle *ändert*, nicht an den, der die
Entscheidung *trifft*. Regeln dieser Sektion: Baseline-Regelwerk
`grundlagen-harness-dateien.md` §Was ein Kommentar trägt.

**Falsch:** <z.B. „Ohne dieses Feld behauptete die Ausgabe eine Verteilung,
die nicht stattgefunden hat"> — Konjunktiv über die verworfene Alternative.
**Richtig:** <z.B. „Verteilt ist wahr, wenn die Splitting-Regel angewendet
werden konnte"> — Indikativ über den Zustand.

**Falsch:** <z.B. „die frühere Fassung prüfte nur die Länge"> — beschreibt
abwesenden Text.
**Richtig:** die geltende Zusage nennen; die vorige hält `git`.

**Zustandsfelder ebenso:** Eine `Stand`-/`Status`-Zelle in Roadmap,
Beobachtungs-Register oder Meilenstein-Tabelle nennt den Zustand und den Beleg
als auflösbaren Anker, nicht die Chronik; das Drift-Log der Roadmap trägt nur
Umplanungen, keine Schließungen und keine erreichten Meilensteine.

**Begründung:** Die Abwägung gehört in die ADR, die Historie in `git`, die
Herkunft in **ein** auflösbares Feld (`LH-*`, `ADR-*`, `· seit welle-<NN>`).
Was daneben steht, liest jeder Lauf mit und bezahlt es mit Kontext.

### 3.8 Action-Pinning (GitHub Actions)

Jede `uses:`-Zeile in `.github/workflows/*.yml` ist auf einen vollständigen
Commit-SHA gepinnt, mit einem Tag-Kommentar dahinter. Ein Tag allein ist
keine Pinnung: Ein Tag ist beweglich, sein Ziel-Commit kann sich ändern,
ohne dass die Workflow-Datei sich ändert — ein Workflow-Schritt läuft mit
den Zugangsdaten dieses Repositories (`GITHUB_TOKEN`, ggf.
Registry-Secrets), ein umgebogener Tag wäre ein stiller Angriffsvektor ohne
Diff im eigenen Repo.

**Falsch:** `uses: actions/checkout@v7`
**Richtig:** `uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`

**Begründung:** Supply-Chain-Härtung. Der Kommentar hält die Pinnung
lesbar und audit-/hebbar, ohne den SHA manuell gegen die Releases der
Action auflösen zu müssen ([`ADR-0051`](docs/plan/adr/0051-cicd-pipeline-github-actions.md)).

### 3.9 Exit-Code eines Gate-Laufs wird direkt geprüft, nie durch eine Pipe/einen Wrapper hindurch

Ein `make`-Gate-Aufruf (`make gates`, `make docs-check`,
`make test-integration`, …), dessen Ausgabe gefiltert wird (`| tail`,
`| grep`, …), oder der über einen Hintergrund-Task-Wrapper läuft, liefert
als Gesamt-Exit-Code den des **letzten** Pipe-Glieds bzw. den des
Wrappers — nicht den von `make` selbst. Schlägt `make` fehl, während
`tail`/`grep` erfolgreich terminiert (Regelfall) oder der Wrapper nur sein
eigenes Ende meldet, zeigt die Shell trotzdem Exit 0: ein rotes Gate
erscheint unbemerkt grün, und eine daran gehängte Folgehandlung
(`&& git push`, Closure, Merge) läuft auf Basis eines falschen Signals.

**Falsch:** `make gates | tail -15 && git push`
**Richtig:** `make gates` ungefiltert laufen lassen und seinen eigenen
Exit-Code unmittelbar danach feststellen, bevor irgendeine Filterung
dazwischentritt (z. B. `make gates; ec=$?; tail -15 <log>; test $ec -eq 0
&& git push` — oder `make gates > <log> 2>&1; ec=$?` mit separatem
`grep`/`tail` gegen die geschriebene Log-Datei). Bei Hintergrund-Tasks
ebenso: Der von einem Wrapper gemeldete „Exit"-Wert ist nicht automatisch
der von `make` — den `make`-eigenen Exit-Code separat sichern (z. B.
`make ...; echo $? > <datei>` als eigener, ungeketteter Schritt).

**Begründung:** Repo-weite Ausführungsdisziplin, die für jede Rolle gilt
(Implementer, Reviewer, Verifier, Planner/Architect) — kein
rollenspezifisches Problem, siehe
[`docs/reviews/architect-verdict-pipe-maskiert-make-exit-code.md`](docs/reviews/architect-verdict-pipe-maskiert-make-exit-code.md).
Dreifach real aufgetreten in `welle-15`
(`docs/plan/planning/observations/BEO-PGC/pipe-maskiert-make-exit-code`):
einmal vom Implementer selbst bemerkt und korrigiert, einmal beim Planner
bis zum realen `git push` durchgerutscht, einmal durch einen
Hintergrund-Task-Wrapper verschärft — jeweils ohne dass ein zweites,
unabhängig einsehbares Artefakt (Diff) den Fehler hätte fangen können, da
er in der Ausführung selbst liegt, nicht im committeten Ergebnis.

**Geschärft — Prüfung und Folgehandlung sind zwei Schritte, nicht einer.**
Ein korrekt (ungepiped) ermittelter roter Exit-Code schützt nur, wenn die
davon abhängige Folgehandlung (`git push`, Merge, Closure) tatsächlich auf
ihn wartet und ihn auswertet, bevor sie beauftragt wird — nicht schon
dadurch, dass der Wert irgendwo sichtbar im Output steht. Ein Gate-Lauf und
seine abhängige Folgehandlung dürfen deshalb nicht im selben
Werkzeug-Aufruf-Batch beauftragt werden: Dazwischen steht ein eigener,
expliziter Schritt, der den zuvor ermittelten Exit-Code auswertet, bevor
die Folgehandlung läuft. In einem einzelnen Shell-Aufruf erzwingt `&&`
diese Reihenfolge mechanisch (siehe die `test $ec -eq 0 && git push`-Form
oben); über mehrere, separat beauftragte Werkzeug-Aufrufe eines
Agenten-Laufs hinweg gibt es diese Kopplung nicht von selbst.

**Falsch:** `make gates` aufrufen und — auch mit korrekt (ungepiped)
ermitteltem, sichtbar rotem Exit-Code — im selben Arbeitsschritt-Batch
`git push` beauftragen; der sichtbare rote Wert verhindert die
Folgehandlung nicht von selbst, wenn sie nicht tatsächlich auf ihn
konditioniert ist.
**Richtig:** `make gates` laufen lassen, das Ergebnis in einem eigenen,
abgeschlossenen Schritt auswerten (Exit-Code lesen, bewusste Entscheidung
treffen), und erst danach — als eigene, separat beauftragte Handlung —
`git push` anstoßen.

Viertes reales Auftreten der zugrunde liegenden Beobachtung
(`docs/plan/planning/observations/BEO-PGC/report-nackte-id-ohne-link`,
`evidence/slice-063-blocker.md`): Der Exit-Code eines `make
gates`-Laufs wurde korrekt und ungepiped ermittelt (sichtbar rot), die
nachfolgende `git push`-Aktion lief aber im selben Arbeitsschritt-Batch,
bevor der bereits sichtbare rote Wert sie tatsächlich blockierte — eine
andere Fehlerklasse als die drei `welle-15`-Fälle oben (dort: Exit-Code
falsch *gemessen*; hier: Exit-Code korrekt gemessen, aber nicht
*wirksam*), siehe
[`docs/reviews/architect-verdict-report-nackte-id-ohne-link-4x.md`](docs/reviews/architect-verdict-report-nackte-id-ohne-link-4x.md).

### 3.10 Ein neuer oder strukturell geänderter GitHub-Actions-Workflow gilt erst nach einem realen, grünen Post-Push-Lauf als abgeschlossen

GitHub-Actions-Workflows laufen auf einem externen, gehosteten Runner.
Jede lokal mögliche Prüfung — YAML-Parse, `yamllint`, `actionlint`,
JSON-Schema-Validierung gegen die SchemaStore-Schemata, `make gates` —
bleibt strukturell statisch: Ob der Workflow auf dem echten Runner
tatsächlich grün läuft (Ressourcen-/Zeitverhalten,
Registry-Zugangsdaten-Pfade, Runner-spezifische
Umgebungsunterschiede), bleibt bis zum ersten realen Lauf unbewiesen.
Das ist keine Prozesslücke, sondern eine Eigenschaft des Werkzeugs
selbst — ein Docker-only-Sensor kann diese Prüfung strukturell nicht
leisten (§3.1).

**Falsch:** einen Slice, dessen Umfang einen neuen Workflow oder eine
strukturelle Änderung an einem bestehenden Workflow enthält (neue
Matrix-Achse, neuer Schritt, neue Job-Abhängigkeit), als
DoD-vollständig behandeln, weil `make gates` grün lief — ohne den
ersten realen Post-Push-Lauf auf GitHub geprüft zu haben (`gh run
view`/`gh run list` oder gleichwertig).
**Richtig:** das betroffene §6-Risiko (Modul 5) bleibt **weiter offen**,
bis der reale Lauf bestätigt ist; die Bestätigung — grüner Lauf oder
roter Befund mit Folgemaßnahme — wird explizit nachgetragen, bevor
Closure erfolgt.

**Begründung:** Dreifach real aufgetreten
(`docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal`,
`evidence/slice-039.md`, `evidence/slice-056.md`, `evidence/slice-064.md`)
— in allen drei Fällen korrekt als *weiter offen* im §6-Risiko geführt
und vor bzw. bei Closure real bestätigt, aber jedes Mal neu hergeleitet
statt einer bereits geltenden Regel gefolgt. Diese Regel macht die
bislang implizite, aber jedes Mal richtige Praxis explizit, damit
künftige Slice-/Wellen-Planungen sie nicht erneut herleiten müssen;
kein neuer Sensor, aus dem oben genannten strukturellen Grund — die
tragende Instanz bleibt die Planner-/Verifier-Disziplin beim
Risiko-Ausgang (Modul 5 §Offene Risiken werden bei Closure aufgelöst),
siehe
[`docs/reviews/architect-verdict-github-actions-unverifizierbar-lokal-3x.md`](docs/reviews/architect-verdict-github-actions-unverifizierbar-lokal-3x.md)
· seit welle-17.

## 4. Quality Gates

Regeln dieser Sektion: Nur Targets aufzählen, die im Makefile **existieren**.
Halluzinierte Gates sind die häufigste Form von Harness-Lüge
(Baseline-Regelwerk `modul-13-quality-gates.md`).

| Target | Zweck |
|---|---|
| `make baseline-verify` | vendored Baseline unverändert (Integrität + Vollständigkeit) |
| `make docs-check` | kaputte Referenzen in der Markdown-Doku (d-check) |
| `make a-check` | Hexagon-Schichten-Edges gegen `.a-check.yml` |
| `make commit-traceability` | Commit-Message-Traceability: je Message ≥ 1 `LH-*`/`ADR-*`, keine `SPEC-*`/`ARC-*` im Betreff; Standing-Gate über die letzten 5 Commits ([ADR-0045](docs/plan/adr/0045-commit-traceability-standing-gate.md)) |
| `make coverage-gate` | Go-Test-Coverage über `./internal/...`+`./cmd/...` gegen `THRESHOLD`; bootstrap-aware Gate, Einstiegsstufe 35 % → Endstufe 80 % ([ADR-0054](docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md), [`harness/sensors/coverage-gate.md`](harness/sensors/coverage-gate.md)) |
| `make gates` | alle inneren Gates (mandatory vor PR), Nachweis-Stempel zuletzt |
| `make image` | baut das OCI-Image, Image-Hash-Beleg (kein Gate) |
| `make image-stale` | advisory: Base-Image-Drift (kein Gate, braucht Netz) |
| `make proto-generate` | Protobuf-/gRPC-Go-Code aus `proto/cdc/stream/v1/changestream.proto` erzeugen; Docker-only über die gepinnte Dockerfile-Stufe `proto` ([ADR-0060](docs/plan/adr/0060-grpc-streaming-mechanismus.md), kein Gate) |
| `<make-target>` | volle Closure (vor Welle-Merge) |

Diese Tabelle **listet auf**; definiert wird hier nichts. Die *Bindung* eines
Targets — welche Anforderung oder Entscheidung es durchsetzt — steht in
`harness/README.md` §Sensors; von dort führt der Weg zur `LH-*`-ID, zur ADR
oder zum Carveout.

## 5. Dokumentations-Regeln

- **Anforderungs-IDs und ADR-Nummern** müssen in PRs/Commits referenziert
  sein — sie sagen, welche Zusage oder Entscheidung berührt ist. Struktur-IDs
  (`SPEC-<NNN>`, `ARC-<NNN>`) adressieren *innerhalb* der Spec und gehören
  nicht in die Commit-Message.
- Vergeben werden IDs beim Spec-/ADR-Schreiben nach dem in
  `harness/conventions.md` deklarierten ID-Schema (Default:
  `<PREFIX>-FA-<NN>` / `<PREFIX>-QA-<NN>` aus dem Lastenheft,
  `SPEC-<NNN>` im Pflichtenheft, `ARC-<NNN>` in der Sicht, ADR-Nummern
  über den ADR-Index) — nie ad hoc im PR.
- Neue ADRs müssen den ADR-Index aktualisieren.
- Roadmap/Status-Geschichte lebt in `docs/plan/planning/`, nicht in `spec/architecture.md`.

## 6. Minimal Agent Workflow

Pro Slice:

1. `harness/README.md` lesen.
2. Relevante kanonische Quelle lesen (Source Precedence beachten).
3. Betroffene Requirement-/ADR-IDs identifizieren.
4. Kleinste sinnvolle Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`) — Exit-Code direkt prüfen, nie durch eine Pipe/einen Wrapper hindurch (§3.9).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten — keine Erfolgsmeldung ohne Gate-Ausführung.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe `harness/README.md` §Guides) →
Verifier. Kein Self-Review — anderer Kontext findet andere Findings,
derselbe Kontext dieselben blinden Flecken (Baseline-Regelwerk
`modul-08-agentenrollen.md`).
