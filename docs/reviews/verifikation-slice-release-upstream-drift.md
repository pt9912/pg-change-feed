# Verifikationsbericht: slice-release-upstream-drift — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-upstream-drift.md` §2) und die
§6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-release-upstream-drift.md`](review-slice-release-upstream-drift.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** zwei Commits auf `main` (Welle
`docs/plan/planning/welle-release-pipeline-adr-0051.md`, `ADR-0051`
Entscheidung 7, Status `Accepted`):

- `b4b45acc` — ursprünglicher Implementer-Commit (sieben neue
  `pin-stale-*`-Skripte/Targets, `.github/workflows/upstream-drift.yml`,
  `harness/README.md`-Update).
- `eb4e7d3f` — Fixrunde nach unabhängigem Review (1 HIGH F-1, 1 MEDIUM
  F-2, 2 LOW F-3/F-4, 1 INFO F-5), inklusive `tools/harness/lib-github-api.sh`
  (neu) und dem Review-Report selbst.

**Frischer Kontext:** Diese Sitzung hat Slice-Plan (vollständig, §1–§8),
Welle-Datei-Ausschnitt, Review-Report (vollständig), `ADR-0051`
(Pin-Inventar-Tabelle P1–P9, §Entscheidung 7, §Verglichene Alternativen)
und alle fünf neuen/geänderten Skripte (`pin-stale.sh`,
`pin-stale-dcheck.sh`, `pin-stale-baseline.sh`, `pin-stale-actions.sh`,
`lib-github-api.sh`) im Volltext selbst gelesen. Nichts aus Slice-Plan,
Commit-Message oder Review-Report wurde ungeprüft übernommen: eigener
realer Lauf aller sieben `make pin-stale-*`-Targets, eigener
`bash tools/harness/pin-stale.sh`-Lauf ohne Argumente, eigener `grep` über
alle `curl`-Vorkommen in den Skripten, eigenes Ruby-Stdlib-YAML-Parsing
von `upstream-drift.yml`, eigener `grep` gegen die referenzierten
Makefile-/`a-check.mk`-/`d-check.mk`-Variablendefinitionen, eigener
ungepipter `make gates`-Lauf mit direkter Exit-Code-Prüfung
(`AGENTS.md` §3.9), eigene isolierte `make doc-commits`/`make
doc-immutable`-Läufe über den exakten Slice-Commit-Bereich, eigene
`git log -1 --format=%B`-Prüfung beider Commit-Messages.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 „Sieben neue … Make-Targets … existieren" — real ausgeführt

Eigener Lauf aller sieben Targets auf dem aktuellen `HEAD` (`eb4e7d3f`):

```
$ make pin-stale-race
DRIFT  TOOLCHAIN_RACE_IMAGE (golang:1.27): gepinnt sha256:b475798f…, aktuell sha256:1cfcdb11…
$ make pin-stale-pgtest
DRIFT  PG_TEST_IMAGE (postgres:18-alpine): gepinnt sha256:63bdc97d…, aktuell sha256:6c538e72…
$ make pin-stale-dmigrate
DRIFT  D_MIGRATE_IMAGE (ghcr.io/pt9912/d-migrate:latest): gepinnt sha256:862dfb04…, aktuell sha256:af9d3eb3…
$ make pin-stale-acheck
OK     A_CHECK_IMAGE (ghcr.io/pt9912/a-check:latest) == sha256:34d3dfb5…
$ make pin-stale-dcheck
OK     DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.77.0) == sha256:3f84502b…
OK     DCHECK_IMAGE Tag-Frische (v0.77.0) == neuester Release
$ make pin-stale-baseline
OK     Kurs-Baseline v6.9.0 == neuester Release
$ make pin-stale-actions
OK     actions/checkout@v7.0.1 Tag-Mutation/Tag-Frische
OK     docker/setup-buildx-action@v4.4.1 Tag-Mutation/Tag-Frische
OK     docker/login-action@v4.6.0 Tag-Mutation/Tag-Frische
```

Deckt sich exakt mit der DoD-Behauptung: P3/P4/P5 zeigen echten Drift, P6/P7/P8/P9
sind aktuell (`OK`). Eigene Gegenprobe der referenzierten Variablen (nicht aus
dem Plan übernommen): `grep`-Abgleich gegen die tatsächlichen
`?=`-Definitionszeilen bestätigt exakte Übereinstimmung —
`TOOLCHAIN_RACE_IMAGE`/`PG_TEST_IMAGE`/`D_MIGRATE_IMAGE` in `Makefile:115-116,193`,
`A_CHECK_IMAGE` in `a-check.mk:9`, `DCHECK_IMAGE`/`DCHECK_DIGEST` in
`d-check.mk:6-7` — kein Skript liest die falsche Datei oder Variable.
**Ergebnis: Checkbox berechtigt auf `[x]`.**

Nebenbefund (kein DoD-Blocker, zur Kenntnis): Bei den drei DRIFT-Fällen meldet
GNU Make selbst `Error 1` und der `make <target>`-Aufruf endet mit `EXIT=2`
(nicht `1`) — das ist reguläres GNU-Make-Verhalten (Make selbst liefert bei
einem gescheiterten Rezept immer Exit 2, unabhängig vom Exit-Code des
darunterliegenden Kommandos), keine Eigenschaft dieses Slice und deckt sich
mit dem bereits dokumentierten Skript-eigenen Vertrag (`pin-stale.sh` selbst
liefert für DRIFT weiterhin Exit `1`, eigens nachgeprüft).

### 1.2 „`upstream-drift.yml` existiert … alle neun Achsen … `if: always()` … `schedule` + `workflow_dispatch`"

Eigene Lektüre plus eigenes Ruby-Stdlib-`YAML.load_file`-Parsing:

```
{"on"=>{"schedule"=>[{"cron"=>"37 3 * * *"}], "workflow_dispatch"=>nil}, …
 "steps"=>[Checkout, "P1/P2 …", "P3 …", "P4 …", "P5 …", "P6 …", "P7 …", "P8 …", "P9 …"]}
```

Acht Achsen-Schritte (P1/P2 gebündelt über `make image-stale` + sieben
einzelne `make pin-stale-*`-Aufrufe) decken die neun Pin-Achsen aus
`ADR-0051`s Inventar-Tabelle vollständig ab — jeder trägt `if: always()`.
`cron: '37 3 * * *'` liegt real 20 Minuten nach `image-scan.yml`s
`cron: '17 3 * * *'` (eigener `grep`-Vergleich beider Dateien) — „versetzt"
ist damit keine bloße Behauptung. Kein Pull-Request-/Push-Trigger. YAML
strukturell fehlerfrei (kein `"on"`-Boolean-Fallstrick, korrekt gequotet).
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 „`harness/README.md` §Werkzeuge trägt alle sieben neuen Targets und `upstream-drift.yml`"

Eigener `grep` über die aktuelle Datei: vier Zeilen (129–132) decken
P3–P6 (gebündelt), P7, P8, P9 einzeln ab, Zeile 142 beschreibt
`upstream-drift.yml` — Inhalt deckt sich mit der real beobachteten
Skript-/Workflow-Semantik (Cron-Versatz, fail-open, advisory, `ADR-0051`-Link
auflösbar, durch `docs-check` `links`/`ids` im vollen `make gates`-Lauf
mitbestätigt). **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 „`make gates` grün" — siehe eigenständiger Abschnitt 4 unten

### 1.5 „Review durchgeführt …" — siehe eigenständiger Abschnitt 6 unten

### 1.6 Die vier verbleibenden `[ ]`-Checkboxen (Doku-Update entfällt, Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiko-Ausgang, drei Paarungen)

Der Plan liegt weiterhin unter `docs/plan/planning/in-progress/` (eigener
`ls`-Beleg), §7 „Closure-Notiz" trägt weiterhin ausschließlich den
Platzhalter `*(wird bei Bearbeitung gefüllt.)*`. Für einen Slice, der noch
nicht nach `done/` gewandert ist, ist das korrekt — diese vier Punkte sind
Closure-Pflichten, keine Liefer-Punkte, und ihr `[ ]`-Zustand widerspricht
sich nicht mit dem übrigen DoD-Bild.

## 2. F-1-Fix real wirksam — `curl` jetzt ausschließlich containerisiert

Eigener `grep -n curl` über alle vier betroffenen Skripte plus die neue
Bibliothek:

```
tools/harness/pin-stale-actions.sh:10   (Kommentar)
tools/harness/pin-stale-baseline.sh:9   (Kommentar)
tools/harness/lib-github-api.sh:5       (Kommentar)
tools/harness/lib-github-api.sh:13      apk add --no-cache curl >/dev/null 2>&1 &&
tools/harness/lib-github-api.sh:14      curl -fsS -m 15 -H 'Accept: …' 'https://api.github.com/$1'
tools/harness/pin-stale-dcheck.sh:8     (Kommentar)
```

Die einzigen zwei echten `curl`-Aufrufzeilen (13/14) liegen beide
**innerhalb** von `docker run "$GITHUB_API_TOOLCHAIN_IMAGE" sh -c "…"` in
`github_api_get()` (`lib-github-api.sh:11-15`) — kein bare-Host-`curl` mehr
in `pin-stale-baseline.sh`, `pin-stale-dcheck.sh` oder `pin-stale-actions.sh`
selbst; alle drei binden die Funktion per `source` ein
(`. tools/harness/lib-github-api.sh`, eigene Lektüre bestätigt an je einer
Stelle je Datei) und rufen ausschließlich `github_api_get "…"` auf. Analog
zum bereits etablierten Vorbild `tools/harness/ci-matrix-abdeckung.sh`
(Kommentar-Text sogar wortgleich referenziert). `git ls-remote` in
`pin-stale-actions.sh` bleibt bewusst bare auf dem Host — das war laut
Review-Report nicht Teil von F-1 (eigenes Vorbild: `git rev-parse
--show-toplevel` in `ci-matrix-abdeckung.sh`), eigene Prüfung bestätigt
diese Abgrenzung als sachlich korrekt. **Ergebnis: F-1 ist real und
vollständig behoben, nicht nur teilweise oder nur behauptet.**

## 3. F-3-Fix real wirksam — `pin-stale.sh` liefert jetzt Exit 2 bei fehlendem Argument

Eigener Lauf, exakt wie im Auftrag verlangt:

```
$ bash tools/harness/pin-stale.sh
UNBESTIMMT  Aufruf — Pflichtargument fehlt (pin-stale.sh <datei> <variable> [vergleichs-tag])
$ echo $?
2
```

Vor der Fixrunde lieferte dieselbe Codezeile laut Review-Report Exit `1`
(`${1:?…}`-Bash-Parameter-Expansion, die einen Nutzungsfehler wie einen
echten `DRIFT`-Fund kodierte). Die jetzige Fassung
(`tools/harness/pin-stale.sh:15-17`) prüft `[ $# -lt 2 ]` explizit und
gibt `exit 2` zurück, bevor die Parameter-Expansion überhaupt greifen
könnte. **Ergebnis: F-3 ist real behoben — der im Kopfkommentar
dokumentierte Exit-Code-Vertrag („2 = … Aufruf-Fehler") gilt jetzt für
jeden Aufrufpfad, nicht nur für die verdrahteten Make-Targets.**

## 4. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
eb4e7d3f1b744c1bf1ad8d213565e445a5dda6b8
$ git status
nothing to commit, working tree clean
$ make gates > /tmp/verifier-gates.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 770 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check`-Modul `commits`: `770 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Die Coverage-Zahl (82.80 %) deckt sich exakt mit der im Review-Report
genannten Zahl für denselben Commit-Bereich — keine Drift wie beim
Nachbar-Slice `release-image-scan` (`AGENTS.md` §3.12 Herkunfts-Disziplin:
gemessen, nicht übernommen).

**Commit-Traceability speziell für `eb4e7d3f` geprüft** (Auftrag: „trägt der
Fixrunden-Commit eine `ADR-*`-ID?"):

```
$ git log -1 --format=%B eb4e7d3f | grep -oE "ADR-[0-9]+|LH-[A-Z0-9-]+"
ADR-0051
$ git log -1 --format=%B b4b45acc | grep -oE "ADR-[0-9]+|LH-[A-Z0-9-]+"
ADR-0051
```

Beide Commits tragen `ADR-0051` im Betreff. Zusätzlich isoliert über den
exakten Slice-Commit-Bereich bestätigt (Parent `58312105` vor
`next → in-progress` bis `eb4e7d3f`):

```
$ make doc-commits RANGE=5831210552cecfc1b9deb52446fa117d10562abe..eb4e7d3f
d-check: 770 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=5831210552cecfc1b9deb52446fa117d10562abe..eb4e7d3f
d-check: 770 Datei(en) geprüft, 0 Befund(e)
$ bash tools/harness/commit-traceability.sh "5831210552cecfc1b9deb52446fa117d10562abe..eb4e7d3f"
commit-traceability: OK — 2 Commit(s) in "…", Betreffs ohne Struktur-ID
```

**Anders als beim Nachbar-Slice `release-image-scan`** (dort hatte der
Fixrunden-Commit keine Struktur-ID und riss das Gate) **ist hier kein
Rotstand entstanden** — der Fixrunden-Commit selbst wurde korrekt mit
`ADR-0051` versehen. **Ergebnis: Die DoD-Checkbox „`make gates` grün." ist
berechtigt auf `[x]` gesetzt** — real, ungepiped, `EXIT=0`, zusätzlich durch
zwei gezielte Modul-Läufe und einen eigenständigen
`commit-traceability.sh`-Lauf über den exakten Slice-Bereich bestätigt.

## 5. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout/Folge-Slice · *entfallen* → gestrichen
mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | `AGENTS.md` §3.10 — realer Post-Push-Lauf unverifiziert | „weiter offen, strukturell (derselbe Fall wie `BEO-PGC/github-actions-unverifizierbar-lokal`)" | ✓ weiter offen | Register-Verzeichnis `docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal/` existiert real (eigener `ls`-Beleg) — kein neuer Eintrag nötig, korrekt referenziert. |
| 2 | P9 deckt ggf. noch nicht alle `uses:`-Zeilen künftiger Workflows | „entfallen … P9 ist ein Scan über den jeweils aktuellen Workflow-Baum" | ✓ entfallen | Begründung trägt: `pin-stale-actions.sh` liest `.github/workflows/*.yml` per Glob zur Laufzeit, kein statisches Inventar — eigene Lektüre bestätigt (`grep -hoE … .github/workflows/*.yml`). |
| 3 | Fail-open kann bei wiederholtem Werkzeugausfall unbemerkt bleiben | „weiter offen, dasselbe Muster wie `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×, noch unter der 3×-Schärfungsschwelle)" | ✓ weiter offen | Register-Verzeichnis existiert real (eigener `ls`-Beleg); Formulierung benennt explizit, dass dieser Slice keinen neuen, eigenständigen Beleg beisteuert (andere Fehlerklasse) — ehrlich, keine stille Zählung. |
| 4 (F-2) | `if: always()` unterscheidet nicht Drift-Fund von Werkzeugausfall auf Lauf-Statusebene | „eingetreten, Doku korrigiert (§1, §2, Workflow-Kopfkommentar) — kein Code-Fix möglich, ohne echte Funde zu verschlucken" | ✓ eingetreten | Eigene Prüfung: `upstream-drift.yml:11-23` trägt jetzt exakt diese präzisierte Aussage („`if: always()` unterscheidet dabei NICHT …", „`continue-on-error: true` wäre keine Abhilfe"); §1/§2 des Plans tragen dieselbe Präzisierung. Bemerkenswert (kein Blocker): `ADR-0051` selbst (§Verglichene Alternativen, Variante D) formuliert noch die ursprüngliche, jetzt widerlegte Fassung „Skip statt Rot" — die ADR bleibt als `Accepted`-Dokument unberührt (`AGENTS.md` §3.5), die Korrektur lebt zulässig im Slice-Plan/Workflow-Kommentar, nicht in der ADR selbst. |
| 5 (F-5) | P9 prüft nur Frische bereits §3.8-konformer Zeilen, keine Form-Konformität | „entfallen als Risiko … P9 war nie als Form-Prüfung geplant … aktuell real folgenlos" | ✓ entfallen | Eigene Stichprobe: `grep -n "uses:" .github/workflows/*.yml` liefert für alle Zeilen SHA+`# vX.Y.Z`-Form — Regex in `pin-stale-actions.sh:29-30` erfasst sie alle; Begründung deckt sich mit dem real beobachteten Zustand. |

**Ergebnis: Alle fünf §6-Risiken tragen eine der drei zulässigen
Ausgangsklassen**, inklusive der beiden neu von der Fixrunde hinzugefügten
Einträge (F-2, F-5) — beide inhaltlich nachvollziehbar und real geprüft,
nicht nur behauptet.

## 6. DoD-Checkbox „Review durchgeführt" — inhaltlich korrekt gegen den Report?

Review-Report-Summary-Tabelle (eigene Lektüre): HIGH 1, MEDIUM 1, LOW 2,
INFO 1 — Findings F-1 (HIGH, Docker-only-Verstoß), F-2 (MEDIUM,
fail-open-Eigenschaft überzeichnet), F-3 (LOW, Exit-Code-Vertrag), F-4
(LOW, doppelte Tabellenzeile), F-5 (INFO, P9-Formgrenze).

DoD-Checkbox-Text (Plan §2, Zeile 87–96) nennt exakt dieselbe Verteilung
(„1 HIGH … 1 MEDIUM … 2 LOW … 1 INFO") mit denselben Kurzbeschreibungen je
Finding-ID — Zahlen und Zuordnung stimmen 1:1 mit dem Report überein, keine
Verwechslung, keine Untertreibung.

Eigene Prüfung, ob „behoben bzw. dokumentiert" für jedes Finding zutrifft:

- **F-1** (HIGH): real behoben, siehe Abschnitt 2 oben.
- **F-2** (MEDIUM): nicht code-behebbar (Review selbst begründet das
  plausibel — `continue-on-error: true` würde echte Funde verschlucken),
  stattdessen dokumentiert — real bestätigt in §1/§2 des Plans und im
  Workflow-Kopfkommentar (Abschnitt 5, Zeile 4 oben).
- **F-3** (LOW): real behoben, siehe Abschnitt 3 oben.
- **F-4** (LOW): eigene Lektüre von Plan §3 (Zeilen 121–130) — genau eine
  Tabellenzeile für die `harness/README.md`-Änderung, keine Dopplung mehr
  auffindbar.
- **F-5** (INFO): als „entfallen"-Risiko dokumentiert, nicht code-behoben
  (laut Review korrekt so vorgesehen — P9 war nie als Form-Prüfung
  geplant), siehe Abschnitt 5, Zeile 5 oben.

**Kein offenes HIGH, kein unbehandeltes MEDIUM/LOW.** Anders als beim
Nachbar-Slice `release-image-scan` (dort riss der Fixrunden-Commit selbst
das `commit-traceability`-Gate) hat die Fixrunde dieses Slice keine neue,
unentdeckte Verletzung erzeugt (Abschnitt 4 oben) — die Checkbox ist damit
sowohl für ihren engen Gegenstand (die fünf Findings) **als auch** als
Stellvertreter für einen validen Gesamt-DoD-Zustand tragfähig.
**Ergebnis: Checkbox berechtigt auf `[x]` gesetzt.**

---

## Verdikt

**DoD erfüllt.** Alle sechs beauftragten Prüfpunkte wurden real und
unabhängig nachgemessen, nicht aus Bericht oder Commit-Message übernommen:

1. Jede der fünf `[x]`-Checkboxen in §2 ist gegen reale Läufe/Dateien
   geprüft und berechtigt gesetzt; die vier verbleibenden `[ ]`-Checkboxen
   sind konsistent mit dem `in-progress`-Zustand (Closure-Pflichten, nicht
   Liefer-Punkte).
2. F-1 ist real und vollständig behoben — kein bare-Host-`curl` mehr in
   den drei betroffenen Skripten, beide echten `curl`-Aufrufe liegen
   innerhalb eines `docker run`-Aufrufs in `lib-github-api.sh`.
3. F-3 ist real behoben — `bash tools/harness/pin-stale.sh` ohne Argumente
   liefert jetzt Exit `2` statt `1`.
4. `make gates` lief eigenständig, ungepiped, mit `EXIT=0`; beide
   Slice-Commits tragen `ADR-0051` in der Message (eigene `git log`-Prüfung),
   der Fixrunden-Commit `eb4e7d3f` reißt — anders als beim Nachbar-Slice
   `release-image-scan` — kein Gate.
5. Alle fünf §6-Risiken (inkl. der beiden von der Fixrunde neu
   hinzugefügten F-2/F-5-Einträge) tragen eine der drei zulässigen
   Ausgangsklassen, mit real nachvollziehbaren Begründungen.
6. Die DoD-Checkbox „Review durchgeführt" ist inhaltlich korrekt gegen den
   tatsächlichen Report (1 HIGH/1 MEDIUM/2 LOW/1 INFO, identische
   Zuordnung) und berechtigt gesetzt — kein offenes HIGH, kein
   unbehandeltes MEDIUM/LOW.

**Nicht blockierende Beobachtung:** `ADR-0051` selbst (§Verglichene
Alternativen, Variante D) trägt weiterhin die durch F-2 widerlegte
Formulierung „Skip statt Rot" — als `Accepted`-ADR zulässig unberührt
(`AGENTS.md` §3.5); die korrigierte Fassung lebt im Slice-Plan und im
Workflow-Kommentar. Kein DoD-Blocker, dem Planner zur Kenntnis für eine
mögliche künftige Zitat-Korrektur oder Folge-ADR, falls die Diskrepanz
zwischen ADR-Text und tatsächlichem Verhalten später verwirrt.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz, Reconciliation-Register-Prüfung,
Beobachtungs-Register-Nachzug, formaler Risiko-Ausgangs-Nachzug in §6, die
drei Paarungen bei der nächsten Welle-Closure).
