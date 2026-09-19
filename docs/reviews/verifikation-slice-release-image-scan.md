# Verifikationsbericht: slice-release-image-scan — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-image-scan.md` §2) und die
§6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-release-image-scan.md`](review-slice-release-image-scan.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** zwei Commits auf `main` (Welle
`docs/plan/planning/welle-release-pipeline-adr-0051.md`, `ADR-0051`
Entscheidung 6, Status `Accepted`):

- `1141786c` — ursprünglicher Implementer-Commit (`Makefile`
  `image-cve`-Target, `.github/workflows/image-scan.yml`,
  `harness/README.md`).
- `1a8c0efb` — Fixrunde nach unabhängigem Review (2 MEDIUM F-1/F-2,
  2 LOW F-3/F-4).

**Frischer Kontext:** Diese Sitzung hat Slice-Plan, Welle-Datei,
Review-Report, `Makefile` (`image-cve`-Target), `.github/workflows/image-scan.yml`
und `ADR-0051` (Entscheidung 6, Verglichene Alternativen 4) selbst
gelesen. Nichts aus Slice-Plan, Commit-Message oder Review-Report wurde
ungeprüft übernommen: eigener realer Lauf von `make image-cve` (ohne und
mit Fake-Credentials), eigener `make -n image-cve`-Trockenlauf (ohne und
mit Variablen), eigene `docker run … image --help`-Prüfung der
Trivy-Flags, eigenes YAML-Parsing von `image-scan.yml` via
Ruby-Stdlib-`YAML.load_file`, eigener ungepipter `make gates`-Lauf mit
direkter Exit-Code-Prüfung (`AGENTS.md` §3.9), eigene isolierte
`make doc-commits`/`make doc-immutable`-Läufe über den exakten
Slice-Commit-Bereich.

---

## 1. `make image-cve` real ausgeführt — ohne und mit (Fake-)Credentials

Ohne Credentials:

```
$ make image-cve
docker run --rm  aquasec/trivy@sha256:62b1e65e... image --image-src remote  --severity CRITICAL,HIGH --exit-code 1 ghcr.io/pt9912/pg-change-feed:latest
[…Vulnerability-DB-Download, 114.61 MiB…]
FATAL  Fatal error  … unable to find the specified image "ghcr.io/pt9912/pg-change-feed:latest" in ["remote"]: 1 error occurred:
	* remote error: GET https://ghcr.io/token?...: DENIED: requested access to the resource is denied
make: *** [image-cve] Error 1
```

Mit Fake-Credentials (`GHCR_USERNAME=test GHCR_PASSWORD=test`):

```
$ make image-cve GHCR_USERNAME=test GHCR_PASSWORD=test
docker run --rm -e TRIVY_PASSWORD="test" aquasec/trivy@sha256:62b1e65e... image --image-src remote --username "test" --severity CRITICAL,HIGH --exit-code 1 ghcr.io/pt9912/pg-change-feed:latest
[…Vulnerability-DB-Download erneut, 114.61 MiB…]
FATAL  Fatal error  … unable to find the specified image "ghcr.io/pt9912/pg-change-feed:latest" in ["remote"]: 2 errors occurred:
	* remote error: GET https://ghcr.io/token?...: DENIED: denied
	* remote error: GET https://ghcr.io/token?...: DENIED: requested access to the resource is denied
make: *** [image-cve] Error 1
```

**Ergebnis: beide Läufe bestätigt.** Beide scheitern strukturell korrekt
am fehlenden Ziel-Image (`ghcr.io/pt9912/pg-change-feed:latest` existiert
nicht — kein echter Release bisher), nicht an einem Syntaxfehler in der
neuen `$(if …)`-Makefile-Logik. Trivy selbst läuft real, lädt seine
Vulnerability-DB (kein Kommandozeilen-Parsing-Fehler, keine leere Ausgabe).
Der Fehlerausgang zeigt in beiden Fällen ausschließlich `["remote"]`
statt der vollen Vier-Quellen-Fallback-Kette `["docker" "containerd"
"podman" "remote"]` — das mit `--image-src remote` behauptete F-3-Fix
wirkt real. Mit Credentials wird sichtbar ein zweiter Versuch mit
DENIED-Antwort protokolliert (Trivy versucht die übergebenen Zugangsdaten
tatsächlich, bevor es aufgibt) — ein Hinweis darauf, dass
`--username`/`TRIVY_PASSWORD` tatsächlich in den Registry-Auth-Pfad
eingespeist werden, nicht nur syntaktisch angenommen und ignoriert
werden.

## 2. `make -n image-cve` — Kommando-Rendering ohne und mit Variablen

```
$ make -n image-cve
docker run --rm  aquasec/trivy@sha256:... image --image-src remote  --severity CRITICAL,HIGH --exit-code 1 ghcr.io/pt9912/pg-change-feed:latest

$ make -n image-cve GHCR_USERNAME=test GHCR_PASSWORD=test
docker run --rm -e TRIVY_PASSWORD="test" aquasec/trivy@sha256:... image --image-src remote --username "test" --severity CRITICAL,HIGH --exit-code 1 ghcr.io/pt9912/pg-change-feed:latest
```

**Ergebnis:** Die `$(if $(GHCR_PASSWORD),-e TRIVY_PASSWORD="$(GHCR_PASSWORD)",)`-
und `$(if $(GHCR_USERNAME),--username "$(GHCR_USERNAME)",)`-Ausdrücke
rendern in beiden Fällen korrekt — ohne Variablen bleibt ein harmloses
doppeltes Leerzeichen an der Einsetzstelle (kosmetisch, kein
Funktionsfehler), mit Variablen erscheinen `-e TRIVY_PASSWORD="test"`
und `--username "test"` an der richtigen Position im Kommando, vor den
übrigen Flags. Kein Make-Syntaxfehler, keine falsch geschachtelte Klammer.

## 3. Begründung gegen den allgemeinen Docker-Credential-Mount — technisch plausibel?

Der Implementer-Bericht (Commit-Message `1a8c0efb`, Slice-Plan-Kommentar
`Makefile:49-57`) begründet die Entscheidung gegen einen generischen
`~/.docker/config.json`-Mount damit, dass ein referenzierter, im
Trivy-Container nicht ausführbarer Credential-Helper *jeden*
Registry-Zugriff des Containers einfärbt — einschließlich Trivys eigenem,
unauthentifiziertem Bezug seiner Vulnerability-DB
(`mirror.gcr.io/aquasec/trivy-db:2`, real in dieser Sitzung beobachtet,
siehe §1-Ausgabe).

Das ist technisch plausibel und deckt sich mit einer bekannten
Fehlerklasse bei `docker/containerd`-basierten Registry-Clients: Ist in
`~/.docker/config.json` ein globaler `credsStore`/`credHelpers`-Eintrag
gesetzt, versucht die zugrunde liegende Auth-Auflösung (die auch Trivys
`go-containerregistry`-Unterbau nutzt), diesen Helper für **jede**
Registry-Interaktion aufzurufen, bevor sie auf anonymen Zugriff
zurückfällt — fehlt der referenzierte Helper-Binary im Container (ein
minimales Trivy-Image trägt keine Docker-Desktop-Helper), bricht die
Auflösung fehlerhaft ab, nicht nur für die authentifizierte Ziel-Registry,
sondern auch für unbeteiligte, eigentlich anonyme Aufrufe. Diese
Fehlerklasse ist in der Praxis vielfach dokumentiert (u. a. als
„docker-credential-… not installed"-Fehler außerhalb dieses Repos).

Die im Slice gewählte Alternative — Trivys eigene, vom
Docker-Credential-Mechanismus unabhängige `--username`/`--password`/
`TRIVY_PASSWORD`-Flags — ist real und dokumentiert: eigene Prüfung via

```
$ docker run --rm aquasec/trivy@sha256:... image --help | grep -iE "username|password|image-src" -A2
--image-src strings   image source(s) to use, in priority order (allowed values: docker,containerd,podman,remote) (default [docker,containerd,podman,remote])
--password strings    password. Comma-separated passwords allowed. TRIVY_PASSWORD should be used for security reasons.
--password-stdin       password from stdin. Comma-separated passwords are not supported.
--registry-token string   registry token
--username strings     username. Comma-separated usernames allowed.
```

**Ergebnis:** `--username`, `--password`/`TRIVY_PASSWORD` und
`--image-src` (Default exakt `[docker,containerd,podman,remote]`, wie im
Review-Finding F-3 behauptet) sind reale, in der Hilfe dokumentierte
Trivy-Flags — kein erfundenes oder halluziniertes Verhalten. Der reale
Nachbau der `htpasswd`/bcrypt-Inkompatibilität selbst wurde nicht erneut
durchgeführt (laut Aufgabenstellung nicht zwingend erforderlich) — die
Plausibilitätsprüfung stützt sich auf die bekannte, dokumentierte
Fehlerklasse und auf die real beobachtete Tatsache, dass Trivys
DB-Bezug in dieser Sitzung ohne jeden Credential-Mount anstandslos
gelang.

## 4. YAML-Struktur von `image-scan.yml` erneut real geparst

Eigener Lauf via Ruby-Stdlib-`YAML.load_file` in einem Container:

```ruby
data = YAML.load_file(".github/workflows/image-scan.yml")
```

Ergebnis (vollständig):

```
"on" => {"schedule"=>[{"cron"=>"17 3 * * *"}], "workflow_dispatch"=>nil}
permissions (top): {}
job permissions: {"contents"=>"read", "packages"=>"read"}
step 0: Checkout (actions/checkout@3d3c42e...)
step 1: env {"GHCR_USERNAME"=>"${{ github.actor }}", "GHCR_PASSWORD"=>"${{ secrets.GITHUB_TOKEN }}"},
        run: make image-cve GHCR_USERNAME="$GHCR_USERNAME" GHCR_PASSWORD="$GHCR_PASSWORD"
```

**Ergebnis:** `"on":` korrekt gequotet (kein YAML-1.1-Boolean-Problem,
parst als String-Key), `schedule` + `workflow_dispatch` beide vorhanden,
`permissions: {}` auf Workflow-Ebene, `contents: read` + `packages: read`
auf Job-Ebene (die für F-1 neu hinzugefügte `packages: read`-Lockerung
ist real vorhanden), `env:` trägt `GHCR_USERNAME`/`GHCR_PASSWORD` aus
`github.actor`/`secrets.GITHUB_TOKEN`, `run:` übergibt beide korrekt als
Shell-Variablen an `make image-cve`. Kein Struktur-Defekt.

## 5. F-2-Fix — Kommentar-Bezug jetzt stabil?

Der Kopfkommentar (`image-scan.yml:24-27`) lautet jetzt:

```
# Bis zum ersten echten Release (`release-version-und-workflow`,
# Welle-Datei `docs/plan/planning/welle-release-pipeline-adr-0051.md`
# §6) existiert kein GHCR-`:latest`-Image — dieser Lauf scheitert bis
# dahin strukturell am fehlenden Ziel-Image, kein Implementierungsfehler.
```

Der beanstandete Pfad `docs/plan/planning/in-progress/release-image-scan.md`
(Lifecycle-Verzeichnis, wandert bei jeder Slice-Closure) ist verschwunden
— ersetzt durch den **Slice-Namen** als stabile Kennung
(`release-version-und-workflow`, ohne Pfad) und den **Welle-Datei-Pfad**
direkt.

Einschränkung, die auch der Review-Report selbst schon benennt (F-2,
Randbemerkung) und die diese Verifikation bestätigt: Die Welle-Datei
selbst ist **ebenfalls** lifecycle-gebunden — ihr eigener Kopf sagt
wörtlich „bei Closure wandert sie per `git mv` nach `done/`". Der neue
Kommentar ist also nicht *pfad-unabhängig*, sondern nur *weniger
volatil* als vorher (Welle-Zustandswechsel sind seltener als
Slice-Zustandswechsel, und die Welle schließt typischerweise erst, wenn
alle ihre Slices fertig sind). Dieses Muster ist kein Novum dieses Diffs:
eigene Prüfung von `.github/workflows/release.yml:26-28` zeigt exakt
dieselbe Zitierform (`Welle-Datei
docs/plan/planning/welle-release-pipeline-adr-0051.md §6`), real bereits
vorher im Repo vorhanden. **Ergebnis:** F-2 ist im behaupteten Umfang
behoben (kein Slice-Lifecycle-Pfad mehr) und folgt einem bereits
akzeptierten Repo-Präzedenzfall — die Restfragilität (Welle-Pfad selbst)
ist dieselbe Klasse wie beim Vorbild, nicht neu eingeführt, und vom
Review bereits transparent benannt statt verschwiegen.

## 6. `make gates` real, ungepiped ausgeführt — **ROT**

```
$ make gates > gates.log 2>&1; echo $?
2
```

Vollständige Einzelbelege aus demselben Lauf (mit `make -k gates`
zusätzlich verifiziert, um alle Checks trotz des Fehlschlags zu sehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 767 Datei(en) geprüft, 0 Befund(e)` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |
| **`commit-traceability`** | **`d-check: 767 Datei(en) geprüft, 1 Befund(e)` → `1a8c0ef:1  1a8c0ef  commit-untraceable  fix(release): Reviewer-Fixrunde für release-image-scan (2 MEDIUM, 2 LOW)` → `make: *** [commit-traceability] Error 1`** |

**Das ist der zentrale Befund dieser Verifikation.** `make gates`
schlägt real und reproduzierbar fehl, weil der Fixrunden-Commit
`1a8c0efb` selbst — der Commit, der die vier Review-Findings behebt —
**keine einzige** `LH-*`- oder `ADR-*`-Kennung in seiner Commit-Message
trägt (weder Betreff noch Body, eigene Prüfung von `git log -1 --format=%B
1a8c0efb`). Das verstößt gegen die Traceability-Regel
(`harness/README.md` §Traceability rules: „PRs/Commits **müssen**
mindestens eine `<LH-*>` oder `ADR-*`-ID nennen") und wird exakt vom
mechanisch getragenen Standing-Gate (`ADR-0045`) gefangen.

Isolierte Bestätigung über den exakten Slice-Commit-Bereich
(`ea631152..1a8c0efb`, Parent vor `next → in-progress` bis zum letzten
Slice-Commit):

```
$ make doc-commits RANGE=ea631152..1a8c0efb
d-check: 767 Datei(en) geprüft, 1 Befund(e)
1a8c0ef:1  1a8c0ef  commit-untraceable  fix(release): Reviewer-Fixrunde für release-image-scan (2 MEDIUM, 2 LOW)
Error 1

$ make doc-immutable RANGE=ea631152..1a8c0efb
d-check: 767 Datei(en) geprüft, 0 Befund(e)
Exit 0

$ bash tools/harness/commit-traceability.sh "ea631152..1a8c0efb"
commit-traceability: OK — 2 Commit(s) in "ea631152..1a8c0efb", Betreffs ohne Struktur-ID
Exit 0
```

Das Bild ist eindeutig: die *negative* Hälfte des Gates
(kein `SPEC-*`/`ARC-*` im Betreff) ist erfüllt, die Subjekt-Zeile selbst
ist unauffällig (`fix(release): …`), aber die *positive* Hälfte (mindestens
eine `LH-*`/`ADR-*`-Kennung irgendwo in der Message) fehlt vollständig.
Kein ADR wurde in diesem Slice inhaltlich berührt — `doc-immutable` ist
zu Recht grün, das ist kein Immutabilitätsproblem.

**Einordnung — warum das kein Reviewer-Fehler ist, sondern eine
Verifier-only-Klasse:** Das Review (`docs/reviews/review-slice-release-image-scan.md`)
lief gegen den Diff des Commits `1141786c` (Parent `ea631152`) — der
Fixrunden-Commit `1a8c0efb` existierte zu diesem Zeitpunkt noch nicht.
Ein Diff-Review kann eine Verletzung nicht sehen, die erst durch den
*eigenen Fix-Commit selbst* entsteht — der Reviewer-Skill deckt bereits
committete Änderungen gegen Plan, nicht die künftige Traceability-
Erfüllung eines noch nicht existierenden Commits. Kein Test hätte das
gefangen (kein Testfall betrifft Commit-Messages). Das ist exakt die
Lücke, die dieser Verifier-Durchlauf schließen soll: „Prüfe die Belege,
nicht die Behauptung — und fahre die Sensoren, deren Ausgabe du nicht
siehst, selbst."

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist im aktuellen
Zustand des Repos NICHT berechtigt auf `[x]` gesetzt.** Sie war es
vermutlich zum Zeitpunkt von Commit `1141786c` (damaliges Fenster
`HEAD~5..HEAD` enthielt `1a8c0efb` noch nicht), ist es aber nicht mehr,
seit `1a8c0efb` selbst in das rollierende Fünf-Commit-Fenster
eingetreten ist — und bleibt es dauerhaft nicht, weil die fehlende
Kennung ein permanentes Merkmal dieser einen Commit-Message ist, kein
transientes Zustandsproblem, das sich von selbst löst, bevor der Commit
aus dem Fenster fällt.

## 7. DoD-Checkbox „Review durchgeführt" — berechtigt gesetzt?

Ausgangslage: Der Review-Report fand 0 HIGH, 2 MEDIUM (F-1, F-2), 2 LOW
(F-3, F-4). Sein eigenes Verdikt zog die Checkbox bewusst **nicht** nach
(„ziehe ich die DoD-Checkbox … **nicht** nach" — die Skill-Regel für
Reviewer-seitigen Nachzug ohne Fixrunde greift laut Report-Text nur bei
0 HIGH **und** keiner Implementer-Rückgabe; hier lag mit F-1/F-2 eine
substanzielle Rückmeldung vor).

Eigene, unabhängige Prüfung des Fixrunden-Commits `1a8c0efb` gegen jedes
Finding:

- **F-1** (MEDIUM): real behoben in der hier geprüften Form —
  `GHCR_USERNAME`/`GHCR_PASSWORD` existieren real im Makefile und im
  Workflow, wirken nachweislich auf den Trivy-Aufruf (§1/§2), `packages:
  read` ist real gesetzt (§4). Der **volle** Authentifizierungspfad gegen
  ein tatsächlich privates GHCR-Paket bleibt jedoch strukturell
  unverifiziert bis zum ersten echten Release — das ist im Risiko-Eintrag
  selbst korrekt benannt (siehe §8 unten), kein verschwiegener Rest.
- **F-2** (MEDIUM): real behoben im behaupteten Umfang, mit der in §5
  benannten Rest-Einschränkung (Welle-Pfad bleibt selbst lifecycle-
  gebunden, aber das ist ein bereits akzeptiertes, vorbestehendes Muster,
  kein neuer Mangel).
- **F-3** (LOW): real behoben — `--image-src remote` real vorhanden,
  Fehlerausgang zeigt nur noch `["remote"]` (§1).
  vor Ort geprüft (§4-Lektüre): „sagte" statt „saegte" (`image-scan.yml:14`).
- **F-4** (LOW): real behoben.

**Kein offenes HIGH, kein unbehandeltes MEDIUM/LOW im engeren Sinn der
vier benannten Findings.** Der Rahmen, in dem die Checkbox „Review
durchgeführt" formal berechtigt gesetzt ist, ist damit für sich genommen
gegeben.

**Aber:** Dieselbe Fixrunde, die die vier Findings behebt, erzeugt selbst
eine neue, unentdeckte Gate-Verletzung (§6) — eine Verletzung, die kein
Review mehr sehen konnte, weil sie im Commit *nach* dem letzten Review-
Lauf entstand, und die vor dem Setzen der Checkbox „`make gates` grün"
hätte auffallen müssen, hätte diese Checkbox tatsächlich einen frischen,
realen `make gates`-Lauf **nach** dem Fixrunden-Commit belegt (nicht nur
den Lauf, den der Review-Report für den *vorherigen* Commit dokumentiert).
Formal betrifft das direkt die `make gates`-Checkbox (§6), nicht die
Review-Checkbox selbst — mit der Konsequenz, dass die Review-Checkbox
zwar für ihren eigenen, engen Gegenstand (die vier Findings) berechtigt
ist, aber **nicht** stellvertretend als Beleg für einen validen
Gesamt-DoD-Zustand gelesen werden darf, solange §6 offen ist.

## 8. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout oder Folge-Slice mit ID · *entfallen*
→ gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | Kein GHCR-`:latest`-Image existiert bis zum ersten Release | „eingetreten und real bestätigt" (DENIED-Fehler auf fehlendes Tag) | ✓ eingetreten | Real reproduziert (§1) — allerdings zeigt der reale Fehlerausgang inzwischen **zwei** DENIED-Ursachen überlagert (fehlendes Image UND fehlende Authentifizierung sind beide potenziell aktiv, sobald ein Image existiert); die Formulierung bleibt für den heutigen Zustand korrekt. |
| 2 (F-1) | GHCR-Paket könnte nach Release privat bleiben → dauerhaftes `DENIED` | „eingetreten, in derselben Fixrunde behoben … der volle Authentifizierungspfad … bleibt … strukturell unverifiziert" | ✓ eingetreten (mit explizit benanntem Rest) | Formal korrekt und **ehrlich** formuliert — behauptet keinen vollständigen Beweis, sondern benennt die verbleibende Lücke im selben Satz (`AGENTS.md` §3.12 Instanz B: Erwartung statt Tatsachenbehauptung). Sauberer wäre gewesen, den unverifizierten Rest als eigenen, separaten „weiter offen"-Punkt zu führen statt ihn als Nebensatz an einen „eingetreten"-Ausgang zu hängen — das ändert aber nichts an der inhaltlichen Ehrlichkeit der Aussage, nur an der taxonomischen Sauberkeit. Kein Blocker, aber dem Planner zur Kenntnis: Dieser Rest bleibt real bis zum ersten echten Release unbeweisbar (strukturell, wie Risiko 1). |
| 3 | `AGENTS.md` §3.10 — realer Post-Push-Lauf unverifiziert | „weiter offen, strukturell (derselbe Fall wie `BEO-PGC/github-actions-unverifizierbar-lokal`)" | ✓ weiter offen | Korrekt einem bereits verkörperten, wiederkehrenden Fall zugeordnet — kein neuer Registereintrag nötig. |
| 4 | Trivy-CVE-DB ändert sich täglich, Scan kann ohne Commit rot werden | „entfallen als Risiko … genau dieses Verhalten ist die bewusste Entscheidung von `ADR-0051` Entscheidung 4" | ✓ entfallen | Begründung trägt — advisory, kein Gate, `ADR-0051` §Verglichene Alternativen 4 verwirft ausdrücklich die blockierende Variante. |

**Ergebnis: Alle vier Risiken tragen formal eine der drei zulässigen
Ausgangsklassen.** Risiko 2 (F-1) verdient eine kleine redaktionelle
Schärfung (den unverifizierten Auth-Rest als eigenen „weiter offen"-Punkt
statt als Nebensatz zu führen), das ist aber eine Qualitäts-, keine
Zulässigkeitsfrage — die Aussage selbst verschweigt nichts.

---

## Verdikt

**DoD NICHT erfüllt — ein realer, reproduzierbarer Gate-Rotstand.**

Sechs der acht beauftragten Prüfpunkte bestätigen die Implementer-
Behauptungen vollständig:

1. `make image-cve` läuft real, ohne und mit Fake-Credentials, beide Male
   strukturell korrekt am fehlenden Ziel-Image gescheitert, `["remote"]`
   statt der vollen Vier-Quellen-Kette.
2. Die `$(if …)`-Makefile-Logik rendert in beiden Fällen syntaktisch
   korrekt (`make -n`).
3. Die Begründung gegen den allgemeinen Docker-Credential-Mount ist
   technisch plausibel; die gewählte Alternative (`--username`/
   `--password`/`TRIVY_PASSWORD`, `--image-src`) ist real dokumentiertes
   Trivy-Verhalten (`--help` selbst geprüft).
4. Die YAML-Struktur von `image-scan.yml` ist korrekt (`permissions`,
   `env`, `run`, `"on":`-Quotierung) — eigenes Ruby-Stdlib-Parsing.
5. Der F-2-Kommentar-Bezug ist jetzt stabiler (Slice-Name statt
   Lifecycle-Pfad), mit einer bereits vom Review benannten und durch
   Präzedenzfall (`release.yml:26-28`) gedeckten Rest-Fragilität.
7. Die DoD-Checkbox „Review durchgeführt" ist für ihren engen Gegenstand
   (die vier Findings) berechtigt gesetzt.
8. Alle vier §6-Risiken tragen eine zulässige Ausgangsklasse (eine kleine
   redaktionelle Schärfungsempfehlung bei Risiko 2/F-1).

**Punkt 6 — `make gates` — bricht das Verdikt:** Ein eigener, ungepipter
Lauf von `make gates` auf dem aktuellen `HEAD` (`1a8c0efb`) endet mit
Exit `2`. Fünf der sechs Gate-Prüfungen sind grün
(`baseline-verify`, `docs-check`, `coverage-gate`, `generated-sync`,
`a-check`); **`commit-traceability` ist rot**, weil der Fixrunden-Commit
`1a8c0efb` — derselbe Commit, der die DoD-Checkboxen „`make gates` grün"
und „Review durchgeführt" auf `[x]` setzt — selbst keine `LH-*`/`ADR-*`-
Kennung in seiner Commit-Message trägt. Isoliert bestätigt über den
exakten Slice-Commit-Bereich (`make doc-commits RANGE=ea631152..1a8c0efb`
→ 1 Befund; `make doc-immutable RANGE=ea631152..1a8c0efb` → 0 Befunde).

Das ist die klassische Verifier-only-Lücke, die diese Rolle abdecken
soll: Der Reviewer prüfte den Diff **vor** diesem Commit und konnte die
durch ihn selbst entstehende Verletzung nicht sehen; kein Test deckt
Commit-Messages ab. Die DoD-Checkbox „`make gates` grün" war zum
Zeitpunkt von `1141786c` mutmaßlich korrekt gesetzt, ist es im jetzigen
Repo-Zustand aber nicht mehr — und das ist kein transientes
Fenster-Artefakt, das sich von selbst löst, sondern eine dauerhafte
Eigenschaft der Commit-Message `1a8c0efb`, die auch dann bestehen bleibt,
wenn das rollierende Fünf-Commit-Fenster sie irgendwann nicht mehr
erfasst.

**Empfehlung an den Planner (keine eigene Korrektur — Verifier berichtet,
repariert nicht):** Vor jeder weiteren Closure (Slice oder Welle) muss
diese Traceability-Lücke geschlossen werden — üblicherweise durch einen
neuen, eigenständigen Commit, dessen Message eine `ADR-*`-Kennung (naheliegend
`ADR-0051`, da der gesamte Fixrunden-Inhalt Entscheidung 6 dieser ADR
betrifft) trägt; ein `git commit --amend` auf einen bereits mit dem
Reviewer abgestimmten historischen Commit wäre eine inhaltliche
Neuschreibung eines bestehenden Commits und sollte nur nach expliziter
Abwägung erfolgen (`AGENTS.md` git-Historie-Disziplin). Erst nach einem
erneuten, realen, ungepipten `make gates`-Lauf mit bestätigtem `EXIT=0`
darf die `make gates`-Checkbox als gültig gelten.

**Nicht blockierend, aber zur Kenntnis:**

- Risiko-Eintrag F-1 (§8, Zeile 2) sollte den unverifizierten
  Authentifizierungs-Rest als eigenen „weiter offen"-Punkt statt als
  Nebensatz eines „eingetreten"-Ausgangs führen — inhaltlich bereits
  ehrlich, nur taxonomisch nicht ganz sauber.
- Der `image-scan.yml`-Kommentarbezug (§5) bleibt, wie beim
  `release.yml`-Vorbild, an die Welle-Datei gebunden, die selbst bei
  Welle-Closure nach `done/` wandert — dieselbe Fragilitätsklasse wie
  vorher, nur mit selteneren Zustandswechseln, kein neuer Mangel.
