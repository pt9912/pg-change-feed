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
(<https://github.com/pt9912/ai-harness-course/releases/download/v6.9.0/lab-regelwerk.zip>)
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

Ausnahme ist die **Zitat-Korrektur**: eine Änderung ausschließlich am Zitat-
und Verweisgerüst — host-lokale Pfade, Linkziele, Zeilen-Lokatoren, die Form
einer gebrochenen Referenz — bei unverändertem Referenten. Sie ist **kein**
inhaltliches Überschreiben und in-place zulässig, wenn sie
[`ADR-0073`](docs/plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
genügt. **Unberührbar** bleiben §Entscheidung, §Konsequenzen, §Verglichene
Alternativen, §Status und die `Supersedes`-Kette; ihre Änderung ist eine neue
ADR mit `Supersedes ADR-NNNN`, nie eine Zitat-Korrektur.

**Beleg:** Die Commit-Message nennt
[`ADR-0073`](docs/plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md);
jede betroffene `Accepted` ADR
erhält **eine** Zeile ihrer §Geschichte-Tabelle (Datum, Ereignis,
Commit-Kennung). Records (`done/`, `docs/reviews/**`) tragen keine §Geschichte —
bei ihnen ist die Commit-Kennung der Beleg.

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

### 3.11 Kein host-lokaler absoluter Pfad in der Doku

Kein Markdown-Dokument dieses Repos nennt an irgendeiner Stelle — Prosa,
Inline-Code oder Fenced-Code-Block — einen host-lokalen absoluten Pfad: ein
Wurzel-Segment eines Entwicklerrechners (Präfixliste `hostpaths.prefixes`) oder
ein Windows-Laufwerks-/UNC-Muster. Ein Schwester-Artefakt wird in der Hausform
zitiert (`d-check`s `Dockerfile`), nicht über seinen Pfad auf einem Rechner.
Eine verbotene Form zeigt ein Dokument **nur als Platzhalter** — `<Host-Wurzel>`
für das Wurzel-Segment; die reale Form steht nirgends, auch nicht im Fence.

**Falsch** (der Pfad beginnt mit dem Wurzel-Segment eines Entwicklerrechners;
der Platzhalter steht für dieses Segment, weil die Regel auch dieses Beispiel
deckt):

```text
Real geprüftes Vorbild: <Host-Wurzel>/d-check/Dockerfile (Stage `coverage`)
```

**Richtig:**

```text
Real geprüftes Vorbild: `d-check`s `Dockerfile` (Stage `coverage`)
```

**Was der Sensor deckt — und was nicht.** Die durchsetzbare Hälfte trägt das
`hostpaths`-Modul in `make docs-check` (`make gates`, Modulliste in
`.d-check.yml`). Es deckt `.md`-Dateien unter `scan.roots` in Prosa und
Inline-Code. Es deckt **nicht**: Fenced-Code-Blöcke (Modul-Design, ohne
Opt-out-Marker), **relative** Pfade, **Nicht-Markdown** (die Skriptkommentare
in `Makefile`, `tools/**` und `harness/mk/**` liest der Scan nicht), und
Dateien unter `scan.ignore`. **Die Fenced-Fläche deckt die Regel voll, der
Sensor nicht** — diese Lücke ist benannt, nicht still; der Wächter dort ist
das Review, kein Gate.

**Begründung:** Ein host-lokaler absoluter Pfad ist eine Aussage über einen
Rechner, nicht über das Repo: nicht portabel, für Mitlesende unauflösbar, und
er verrät das Maschinen-Layout.

**Träger und Anker:** Die Regel wirkt über das `hostpaths`-Modul in
`make docs-check`; ihre Reichweite trägt
[`ADR-0075`](docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md), die
Aktivierung
[`ADR-0072`](docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
die Zitationsform eines Schwester-Repos
[`ADR-0074`](docs/plan/adr/0074-zitationsform-schwester-repo-hausform.md).
Dieser Abschnitt trägt die **Regel und ihre Reichweite**; die Präfixliste
führt allein das Modul (`hostpaths.prefixes`), und was der Sensor deckt und
was nicht, führt
[`harness/sensors/docs-check.md`](harness/sensors/docs-check.md) aus seiner
Sicht · seit slice-078.

### 3.12 Eine Aussage, die als Beleg gelesen wird, trägt ihren Ursprung

Regeln dieser Sektion: [`ADR-0083`](docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
— dort die Entscheidung, ihre Begründung, ihre Grenze und ihr Geltungsbereich.

**Instanz A — Zahlenwert in einem Träger.**

> **Aussage.** Jede Zahl eines Doku-Trägers trägt ihren **Ursprung**: ob sie
> **gemessen**, **übernommen** oder **abgeleitet** ist. Ist sie eine Messung,
> trägt sie zusätzlich den **Lauf**, aus dem sie stammt (die gedruckte Zeile).
> Ein aus Bericht oder Nachbardokument übernommener Wert wird nachgemessen
> oder als **übernommen** gekennzeichnet; ein **abgeleiteter** Wert (Summe,
> Produkt, Differenz, Prozent) wird als **abgeleitet** gekennzeichnet — nie
> als gemessen ausgegeben.
>
> **Bewegliche Zahlen.** Der **Nenner** hängt am Code-Stand und ist eine
> **Zustandsgröße**; die **gedeckte** Zahl hängt am Lauf und ist der **Beleg
> eines konkreten Laufs** — sie nennt ihn und nie „der Ist-Stand". Der
> Schreiber setzt den Zeitpunkt; der Leser kann ihn nicht erraten.

**Instanz B — Tatsachenbehauptung in Plan oder Begründung.**

> **Aussage.** Eine Begründung, die eine **Tatsache über den Gegenstand**
> behauptet — ein DoD-Kriterium („beweist …"), ein Plan-Satz („wie die drei
> anderen"), eine Konsequenz —, nennt den **Beleg-Anker**, an dem sie geprüft
> wurde (Befehl, Datei, Abfrage), **oder** sie ist als **erwartet** formuliert
> („zu belegen durch …"). Eine ungeprüfte Übernahme steht nie als geprüfte
> Aussage im Text: was aus Bericht, Nachbardokument oder Erinnerung stammt,
> wird nachgemessen oder als **übernommen** gekennzeichnet. Ein Kriterium, das
> erst nach der Arbeit belegt werden kann, ist eine **Zusage** — es ist als
> solche formuliert und nicht als Feststellung.

**Was diese Regel nicht hat — die benannte Grenze.** Ein Sensor ist **nicht**
Teil der Entscheidung: verlangte er, dass jede Zahl ihren Ursprung trägt, wäre
er eine **Formpflicht auf Prosa**, erzeugte Pflichterfüllung und hätte genau
die Klasse, die er prüfen soll. Die verfügbare **Falsifikation ist die Messung
selbst**. Die durchsetzenden Leser existieren bereits und sind Rollen, kein
Werkzeug: der **Reviewer** für Instanz A (diff-skopiert; eigener HIGH-Unterpunkt
in [`.harness/skills/reviewer.md`](.harness/skills/reviewer.md) — die vier
belegten Fälle stehen in
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`,
er hat sie durch eigenes Nachmessen gefunden), der **Verifier**
für Instanz B („Prüfe die **Belege**, nicht die Behauptung") — und **der
Planner als Verfasser** der Begründung selbst, der sie beim Schreiben an ihren
Anker bindet. Für Träger **außerhalb** des Diffs bleibt als Leser allein die
Messung.

**Benachbarte Regel — und die Abgrenzung zu ihr.** §3.7 bleibt die Regel für
den **Kommentar** und das **Zustandsfeld**; ihre Kernaussage ist eine
geschlossene Liste von Kommentar-Klassen, die der Reviewer-Skill namentlich
adressiert. Dieser Abschnitt gilt für die **Aussage** in einem Doku-Träger und
ist §3.7s Geschwister-Ort, nicht seine Erweiterung. Der **Geltungsbereich** —
Doku-Träger samt den Begründungen in ihnen, **nicht** Testausgaben, Lauf-Logs
und `git`-Historie — steht in
[`ADR-0083`](docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md).

**Träger und Anker:** Diese Regel wirkt durch **Lesen** (Reviewer · Verifier),
nicht durch ein Gate; ihre Begründung führt
[`ADR-0083`](docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md) · seit slice-089.

### 3.13 Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach

**Aussage.** Wer eine Änderung landet, die eine **Eigenschaft eines Gegenstands
bewegt**, die anderswo **beschrieben** steht — eine Deckungszahl, eine Anzahl,
eine Liste, eine Gruppierung —, sucht diese Träger und zieht sie nach. Der
Suchlauf ist eine **Handlung**, keine Erinnerung: `grep` über die Träger nach
der **bewegten Eigenschaft** und nicht über den eigenen Diff, **beide Stände
gemessen** (Parent und Diff). Sein Ergebnis steht im Bericht des Slice — mit
dem, was er fand, **und** mit dem, was er nicht fand. Ein Träger, der eine
fremde Datei betrifft, wird **gemeldet** statt still mitgeändert: die Meldung
ist das Übergabe-Artefakt, die stille Mitänderung der blinde Übergang
(Baseline-Regelwerk `modul-08-agentenrollen.md` §Die neun Übergaben).

**Warum diese Regel einen Träger braucht und keinen Vorsatz.** Die Träger, die
eine bewegte Eigenschaft beschreiben, stehen **nicht im Diff**, und **kein
Sensor liest sie**: `docs-check` prüft Referenzen, `coverage-gate` eine Zahl
gegen eine Schwelle, `generated-sync` Bytes. Falsch werden sie trotzdem, und
zwar in dem Moment, in dem die Arbeit landet. Der billigste Wächter ist die
Frage „welche Träger beschreiben das, was ich hier gerade geändert habe?".

**Benachbarte Regel — und die Abgrenzung zu ihr.** §3.12 bleibt die Regel für
die **Aussage**, die ihren Ursprung trägt: dort fehlt der Ursprung, oder ein
Wert driftet gegen die Messung. Hier trug die Aussage ihren Ursprung
**korrekt** und war zum Zeitpunkt ihrer Niederschrift **wahr** — falsch wird
sie erst **durch diese Arbeit**, und die Arbeit selbst war richtig. Die Grenze
zwischen beiden Regeln ist der **Zeitpunkt**, nicht die Form des Satzes. Der
Grund, den diese Regel trägt, ist zugleich der einzige, den
[`ADR-0085`](docs/plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md)
§Entscheidung 2 außerhalb der zwei Bäume und der drei Messgrößen zulässt.

**Kein Sensor.** Ein Sensor müsste wissen, welche Sätze von welcher Eigenschaft
abhängen; die Abhängigkeit steht in Prosa. Die verfügbare Falsifikation ist die
**Messung** an beiden Ständen.

**Belegte Fälle (Herkunft).** `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(3×, `slice-091`/`-093`/`-094`) — die Fundstellen, ihre Messungen an beiden
Ständen und der abgelehnte Kandidat stehen in den Beleg-Dateien des Eintrags.

**Wer sie liest:** der **Implementer** führt den Suchlauf und trägt sein
Ergebnis in den Bericht; der **Planner** benennt ihn dort, wo die bewegte
Eigenschaft schon bei der Planung bekannt ist; der **Reviewer** und der
**Verifier** prüfen das **berichtete Ergebnis**, nicht die Behauptung, gesucht
zu haben.

**Träger und Anker:** Diese Regel wirkt durch **Ausführen** und **Messen**,
nicht durch ein Gate · seit welle-20.

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
| `make coverage-gate` | Go-Test-Coverage über die netzlos prüfbare Fläche (`./internal/...`+`./cmd/...`+`./gen/...` ohne die Pakete, deren Testlauf einen externen Dienst voraussetzt) gegen `THRESHOLD`; bootstrap-aware Gate, Einstiegsstufe 70 % → Endstufe 80 % ([ADR-0071](docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), [ADR-0054](docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md), [ADR-0077](docs/plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md), [`harness/sensors/coverage-gate.md`](harness/sensors/coverage-gate.md)) |
| `make generated-sync` | Sync-Gate: der committete Protobuf-/gRPC-Code ist byte-gleich der Ausgabe des gepinnten Generators (Dockerfile-Stufe `proto`) aus der committeten `.proto`-Quelle ([ADR-0084](docs/plan/adr/0084-sync-gate-fuer-generierte-artefakte.md), [ADR-0060](docs/plan/adr/0060-grpc-streaming-mechanismus.md)) |
| `make gates` | alle inneren Gates (mandatory vor PR), Nachweis-Stempel zuletzt |
| `make image` | baut das OCI-Image, Image-Hash-Beleg (kein Gate) |
| `make image-stale` | advisory: Base-Image-Drift (kein Gate, braucht Netz) |
| `make proto-generate` | Protobuf-/gRPC-Go-Code aus `proto/cdc/stream/v1/changestream.proto` erzeugen; Docker-only über die gepinnten Dockerfile-Stufen `proto`/`proto-export` — Erzeugung zur Build-Zeit, Host-Extraktion aus einem `tar`-Stream über stdout (seit slice-104, kein Bind-Mount mehr) ([ADR-0060](docs/plan/adr/0060-grpc-streaming-mechanismus.md), kein Gate) |
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
