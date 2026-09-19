# Verifikationsbericht: slice-release-hub-description — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-hub-description.md` §2) und die
§6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-release-hub-description.md`](review-slice-release-hub-description.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** zwei Commits auf `main` (Welle
`docs/plan/planning/welle-release-pipeline-adr-0051.md`, `ADR-0051`
Entscheidung 8, Status `Accepted`):

- `36910422` — ursprünglicher Implementer-Commit (`.github/workflows/hub-description.yml`
  neu, `.github/workflows/release.yml` erweitert, `harness/README.md`-Update).
- `5cc7fa49` (= `HEAD`) — Fixrunde nach unabhängigem Review (2 HIGH F-1/F-2
  „Zitat nennt die falsche Stelle", 1 MEDIUM F-3, 1 LOW F-4, 2 INFO F-5/F-6),
  inklusive `tools/harness/dockerhub-token.sh` + `run-dockerhub-token-tests.sh`
  + `make test-dockerhub-token` (neu) und dem Review-Report selbst.

**Frischer Kontext:** Diese Sitzung hat Slice-Plan (vollständig, §1–§8),
Welle-Datei-Ausschnitt (§6 Out-of-Scope), `ADR-0051` (§Entscheidung,
Punkt 8 im Volltext), Review-Report (vollständig) und alle geänderten
Dateien (`hub-description.yml`, `release.yml`-Auszug,
`dockerhub-token.sh`, `run-dockerhub-token-tests.sh`) selbst gelesen.
Nichts aus Slice-Plan, Commit-Message oder Review-Report ungeprüft
übernommen: eigener `grep -n "^## "` gegen die volle Datei von [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), um
die zitierte Stelle selbst aufzuschlagen; eigener `grep -n` gegen die
Welle-Datei-Abschnitte; eigener `make test-dockerhub-token`-Lauf; zwei
eigene, selbst gewählte `stdin`-Eingaben direkt gegen `dockerhub-token.sh`
(unabhängig vom committeten Testskript); eigener ungepipter
`make gates`-Lauf mit direkter Exit-Code-Prüfung (`AGENTS.md` §3.9);
eigene isolierte `make doc-commits`/`make doc-immutable`-Läufe über den
exakten Slice-Commit-Bereich; eigene `git log -1 --format=%B`-Prüfung
beider Commit-Messages.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 „`hub-description.yml` existiert … synchronisiert … über `POST /v2/auth/token` + `PATCH /v2/repositories/…` … Token-Extraktion über `dockerhub-token.sh`"

Eigene Lektüre der Datei (77 Zeilen): `on: {workflow_dispatch:,
workflow_call: {secrets: {DOCKERHUB_USERNAME: {required: true},
DOCKERHUB_TOKEN: {required: true}}}}` — exakt wie behauptet. Der
`run:`-Block ruft `POST https://hub.docker.com/v2/auth/token` mit
`{identifier, secret}` (per `jq -n`), pipet die Antwort durch
`bash tools/harness/dockerhub-token.sh` (kein Inline-`jq` mehr für die
Extraktion, siehe 1.4 unten), bricht bei leerem Token mit `::error::` +
`exit 1` ab, baut sonst `{full_description: $(cat README.md)}` und ruft
`PATCH https://hub.docker.com/v2/repositories/pt9912/pg-change-feed` mit
`curl -fsS`. Kein gepinnter Action-Fork — einzige `uses:`-Zeile ist
`actions/checkout`. Die im DoD-Text behaupteten `400`/`401`-Funde stammen
aus dem Implementer- bzw. Reviewer-Bericht (reale Netzaufrufe gegen die
öffentliche API); diese Sitzung hat sie nicht erneut gegen das echte Netz
nachgefahren (nicht Teil des Prüfauftrags), aber die Fehlerpfad-Logik
selbst eigenständig mit Fake-Eingaben gegen `dockerhub-token.sh` bestätigt
(siehe 1.4). **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 „`release.yml` bekommt zusätzlichen Job `hub-description` mit `needs: release` … `uses: ./.github/workflows/hub-description.yml` mit `secrets: inherit`"

Eigener `grep`/Lektüre von `release.yml:108-113`:

```
hub-description:
  name: Docker-Hub-Beschreibung synchronisieren
  needs: release
  uses: ./.github/workflows/hub-description.yml
  secrets: inherit
```

Exakt wie behauptet. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 „`make gates` grün" — siehe eigenständiger Abschnitt 5 unten

### 1.4 „Review durchgeführt …" — siehe eigenständiger Abschnitt 6 unten

### 1.5 Die vier verbleibenden `[ ]`-Checkboxen (Doku-Update entfällt, Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiko-Ausgang, drei Paarungen)

Eigener `ls`-Beleg: Der Plan liegt weiterhin unter
`docs/plan/planning/in-progress/release-hub-description.md`, §7
„Closure-Notiz" trägt weiterhin ausschließlich
`*(wird bei Bearbeitung gefüllt.)*`. Für einen Slice, der noch nicht nach
`done/` gewandert ist, ist das korrekt — diese Punkte sind
Closure-Pflichten, keine Liefer-Punkte, ihr `[ ]`-Zustand widerspricht
sich nicht mit dem übrigen DoD-Bild.

**Eine Beobachtung, kein Blocker:** Die Checkbox „Jedes Risiko aus §6
trägt einen Ausgang" ist unverändert `[ ]`, obwohl inhaltlich **alle
vier** §6-Einträge bereits einen der drei zulässigen Ausgänge tragen
(siehe Abschnitt 4 unten — Inhalt ist voraus, Häkchen folgt erst bei
Closure). Kein DoD-Verstoß, da der Slice noch nicht geschlossen ist; dem
Planner zur Kenntnis für den formalen Nachzug bei Closure.

## 2. F-1-Fix real wirksam — Zitat in `hub-description.yml:7`

Vor der Fixrunde zitierte der Kopfkommentar „[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) §Konsistenz-Prüfung"
für den Satz „Präsentation, kein Bestandteil der Distributions-Zusage".
Eigene Lektüre der aktuellen Datei (Zeile 4–7):

```
# Docker-Hub-Beschreibungs-Sync (ADR-0051 Entscheidung 8): spiegelt den
# Inhalt von README.md als Docker-Hub-Repository-Beschreibung fuer
# pt9912/pg-change-feed. Praesentation, kein Bestandteil der
# Distributions-Zusage (ADR-0051 Entscheidung 8) — ...
```

Eigener `grep -n "^## "` gegen die volle Datei von [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) zeigt: es gibt
**keinen** Abschnitt „§Konsistenz-Prüfung" mit eigener Überschrift — der
zitierte Satz muss innerhalb der einzigen `## Entscheidung`-Sektion
liegen. Eigene Lektüre von `docs/plan/adr/0051-cicd-pipeline-github-actions.md`
Zeile 170–174 (Punkt 8 der nummerierten Liste unter `## Entscheidung`):

> 8. **`hub-description.yml`:** … Präsentation, kein Bestandteil der
> Distributions-Zusage; ein Fehlschlag macht das Release nicht rot.

Der Wortlaut steht wörtlich und ausschließlich dort — eigener `grep -c
"Distributions-Zusage"` gegen die volle Datei liefert `1`. Kein anderer
Abschnitt („Kontext", „Direkt ergänzte Hard Rule", „Verglichene
Alternativen", „Konsequenzen", „Fitness Function",
„Re-Evaluierungs-Trigger", „Geschichte") enthält diese Formulierung.
**Ergebnis: F-1 ist real und vollständig behoben — der zitierte Anker
trägt jetzt die Aussage, die er stützen soll.**

## 3. F-2-Fix real wirksam — Zitat in der DoD-Zeile des Slice-Plans

Vor der Fixrunde verwies die DoD-Zeile auf „§1 Abgrenzung der
Welle-Datei". Eigene Lektüre der aktuellen DoD-Zeile (§2, Zeile 58–59):

> `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN` referenziert, nicht angelegt —
> siehe §6 Out-of-Scope der Welle-Datei.

Eigener `grep -n "^## "` gegen `welle-release-pipeline-adr-0051.md` zeigt
sieben Abschnitte (`## 1.` … `## 7.`); Abschnitt 6 heißt „Out-of-Scope für
diese Welle" — kein „Abgrenzung"-Abschnitt existiert in dieser Datei.
Eigene Lektüre von §6, zweiter Aufzählungspunkt:

> **Anlage der `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`-Repository-Secrets**
> — eine externe, kontobezogene Handlung im GitHub-Repo-Settings, die nur
> der Auftraggeber selbst ausführen kann; kein technischer
> Slice-Gegenstand. `release.yml` referenziert die Secret-Namen (`ADR-0051`
> Entscheidung 2), ihr tatsächliches Vorhandensein bleibt bis zum ersten
> echten Release unbewiesen (`AGENTS.md` §3.10 analog).

Die zitierte Stelle trägt exakt die im Slice-Plan behauptete Aussage.
**Ergebnis: F-2 ist real und vollständig behoben.**

## 4. F-3-Fix real wirksam — `make test-dockerhub-token`, eigene stdin-Proben, Nutzung in `hub-description.yml`

**4.1 Eigener Lauf des committeten Testskripts (netzlos):**

```
$ make test-dockerhub-token
bash tools/harness/run-dockerhub-token-tests.sh
run-dockerhub-token-tests: alle Fälle bestanden
$ echo $?
0
```

Grün, netzlos (kein Docker, kein Netzzugriff — reines `bash`).

**4.2 Eigene, vom committeten Testskript unabhängige `stdin`-Proben gegen
`dockerhub-token.sh` direkt** (selbst gewählte Eingaben, nicht aus
`run-dockerhub-token-tests.sh` übernommen):

```
$ printf '{"access_token": "verifier-check-token-xyz"}' | bash tools/harness/dockerhub-token.sh; echo "exit=$?"
verifier-check-token-xyz
exit=0

$ printf '{"message":"unauthorized: incorrect username or password","errinfo":{}}' | bash tools/harness/dockerhub-token.sh; echo "exit=$?"
dockerhub-token: kein access_token in der Antwort ({"message":"unauthorized: incorrect username or password","errinfo":{}})
exit=1
```

Der gültige Fall liefert exakt den eingebetteten Tokenwert auf `stdout`
mit Exit `0`; der ungültige Fall (realistische 401-Fehlerform, in dieser
konkreten Formulierung nicht im committeten Testskript enthalten) liefert
Exit `1` mit einer erklärenden `stderr`-Meldung und **keiner** Ausgabe auf
`stdout`. Das Skript funktioniert unabhängig vom eigenen Testharness
korrekt.

**4.3 Nutzung in `hub-description.yml` — kein Inline-`jq` mehr:**

Eigener `grep -n "jq\|dockerhub-token.sh"` gegen die aktuelle Datei: Die
Token-Extraktion läuft über
`| bash tools/harness/dockerhub-token.sh` (Zeile 83); die beiden
verbliebenen `jq`-Aufrufe bauen ausschließlich Request-**Bodies**
(`jq -n --arg u … '{identifier: $u, secret: $p}'` und
`jq -n --arg desc … '{full_description: $desc}'`), keiner davon extrahiert
mehr aus einer Antwort. **Ergebnis: F-3 ist real und vollständig
behoben** — committete, netzlose Negativtestabdeckung existiert, ist grün,
und das Produktionsskript nutzt sie tatsächlich statt eigener,
unabhängig getesteter Inline-Logik.

## 5. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
5cc7fa4976bf3bf7d1a5cc02acf3f7b2e770d9a8
$ git status --short
(leer)
$ make gates > /tmp/gates_verifier_run.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 772 Datei(en) geprüft, 0 Befund(e)` (volle Modul-Liste) |
| `commit-traceability` | `d-check`-Modul `commits`: `772 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Nicht blockierende Beobachtung: Der Review-Report selbst maß auf dem
Vorgänger-Commit `36910422` noch `82.70 %`; diese Sitzung misst auf `HEAD`
(`5cc7fa49`, nach der Fixrunde) `82.80 %` — beides sind Messungen auf
verschiedenen Commits desselben Slice, keine Drift zwischen Behauptung
und Messung am selben Stand (`AGENTS.md` §3.12 eingehalten: beide Werte
sind als Messung ihres jeweiligen Laufs ausgewiesen, keiner wird als „der"
Ist-Stand ausgegeben).

**Commit-Traceability speziell für `5cc7fa49` geprüft** (Auftrag: „trägt
der Fixrunden-Commit eine `ADR-*`-ID?"):

```
$ git log -1 --format=%B 5cc7fa49 | grep -oE "ADR-[0-9]+|LH-[A-Z0-9-]+"
ADR-0051
$ git log -1 --format=%B 36910422 | grep -oE "ADR-[0-9]+|LH-[A-Z0-9-]+"
ADR-0051
$ git log -1 --format=%s 36910422
feat(release): Docker-Hub-Beschreibungs-Sync nach Release (ADR-0051)
$ git log -1 --format=%s 5cc7fa49
fix(release): Reviewer-Fixrunde für release-hub-description (ADR-0051, 2 HIGH, 1 MEDIUM, 1 LOW, 2 INFO)
```

Beide Commits tragen `ADR-0051` im Betreff, keine `SPEC-*`/`ARC-*`-Kennung
im Betreff. Zusätzlich isoliert über den exakten Slice-Commit-Bereich
bestätigt (Parent `869f8373` vor `next → in-progress` bis `5cc7fa49`):

```
$ make doc-commits RANGE=869f837345dda7f961a23d6f9bb29beb7ecc4971..5cc7fa4976bf3bf7d1a5cc02acf3f7b2e770d9a8
d-check: 772 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=869f837345dda7f961a23d6f9bb29beb7ecc4971..5cc7fa4976bf3bf7d1a5cc02acf3f7b2e770d9a8
d-check: 772 Datei(en) geprüft, 0 Befund(e)
```

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, `EXIT=0`, zusätzlich durch zwei
gezielte Modul-Läufe über den exakten Slice-Bereich bestätigt. Der
Fixrunden-Commit reißt kein Gate.

## 6. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout/Folge-Slice · *entfallen* →
gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | `AGENTS.md` §3.10 — realer Post-Push-Lauf unverifiziert | „weiter offen, strukturell (derselbe Fall wie `BEO-PGC/github-actions-unverifizierbar-lokal`)" | ✓ weiter offen | Register-Verzeichnis existiert real (eigener `ls`-Beleg unter `docs/plan/planning/observations/BEO-PGC/`), korrekt referenziert. |
| 2 | Scope-Risiko: Docker-Hub-API-Aufruf braucht ggf. anderen Token-Scope als der Image-Push | „eingetreten — im Schwester-Repo d-check real dokumentiert … `read/write/delete`-Scope … 403 Forbidden mit `read/write`" | ✓ eingetreten | Eigener `grep -n "read/write/delete\|403 Forbidden"` gegen `hub-description.yml` bestätigt: Zeile 37–41 trägt exakt diesen Scope-Hinweis im Kopfkommentar. Die zitierte externe Quelle (`d-check`s `packaging/dockerhub/README.md` §Transport) liegt außerhalb dieses Repos und war für diese Sitzung nicht nachprüfbar — die **Wirkung im eigenen Repo** (der Kommentar steht tatsächlich dort) ist bestätigt, die externe Quellenaussage selbst bleibt unverifiziert von hier aus. |
| 3 (F-5) | Erfolgspfad-Feldnamen bleiben bis zum ersten realen Lauf unbewiesen | „weiter offen, strukturell (dieselbe Kategorie wie `AGENTS.md` §3.10)" | ✓ weiter offen | Konsistent mit §3.10-Regel — kein Netzbeleg mit echten Secrets möglich, korrekt als offen geführt. |
| 4 (F-6) | Permissions-Vererbung bei `uses: ./…yml` ohne expliziten `permissions:`-Block im Aufrufer-Job | „weiter offen, strukturell (derselbe Fall wie `AGENTS.md` §3.10)" | ✓ weiter offen | Eigene Lektüre von `release.yml`s `hub-description`-Job bestätigt: kein eigener `permissions:`-Block dort, `hub-description.yml` selbst trägt `contents: read` auf Job-Ebene — die Unklarheit ist real vorhanden und nur durch einen echten Lauf klärbar. |

**Ergebnis: Alle vier §6-Risiken tragen eine der drei zulässigen
Ausgangsklassen**, inklusive des neu als „eingetreten" aufgelösten
Scope-Risikos und der beiden neuen F-5/F-6-Einträge — der im Kopfkommentar
behauptete Scope-Hinweis ist real im Code vorhanden, nicht nur im Plan
behauptet.

## 7. DoD-Checkbox „Review durchgeführt" — inhaltlich korrekt gegen den Report?

Review-Report-Summary-Tabelle (eigene Lektüre): HIGH 2, MEDIUM 1, LOW 1,
INFO 2 — Findings F-1 (HIGH, Zitat ADR-Stelle falsch), F-2 (HIGH, Zitat
Welle-Datei-Stelle falsch), F-3 (MEDIUM, fehlende Negativtests), F-4 (LOW,
`curl -f`-Asymmetrie nur in Commit-Message erklärt), F-5 (INFO,
Erfolgspfad unbewiesen), F-6 (INFO, Permissions-Vererbung ungeklärt).

DoD-Checkbox-Text (Plan §2, Zeile 70–79) nennt exakt dieselbe Verteilung
(„2 HIGH … 1 MEDIUM … 1 LOW … 2 INFO") mit denselben Kurzbeschreibungen je
Finding-ID — Zahlen und Zuordnung stimmen 1:1 mit dem Report überein.

Eigene Prüfung, ob „behoben bzw. dokumentiert" für jedes Finding zutrifft:

- **F-1** (HIGH): real behoben, siehe Abschnitt 2 oben.
- **F-2** (HIGH): real behoben, siehe Abschnitt 3 oben.
- **F-3** (MEDIUM): real behoben, siehe Abschnitt 4 oben.
- **F-4** (LOW): eigene Lektüre von `hub-description.yml:74-79` — der
  erklärende Kommentar zur `curl -f`-Asymmetrie steht jetzt direkt an der
  Stelle im Code, nicht mehr nur in der Commit-Message.
- **F-5** (INFO): als „weiter offen"-Risiko dokumentiert, siehe Abschnitt
  6, Zeile 3 oben — laut Review korrekt so vorgesehen (strukturell nicht
  code-behebbar vor einem realen Lauf).
- **F-6** (INFO): als „weiter offen"-Risiko dokumentiert, siehe Abschnitt
  6, Zeile 4 oben — ebenso strukturell erst durch einen realen Lauf
  klärbar.

**Kein offenes HIGH, kein unbehandeltes MEDIUM/LOW.** Die Fixrunde hat
keine neue, unentdeckte Verletzung erzeugt (Abschnitt 5 oben: `make gates`
bleibt grün, kein Gate reißt). **Ergebnis: Checkbox berechtigt auf `[x]`
gesetzt.**

---

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten in §2 sind bewusst noch offen, siehe 1.5). Alle sieben
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht aus
Bericht oder Commit-Message übernommen:

1. Beide `[x]`-Liefer-Punkt-Checkboxen (`hub-description.yml`,
   `release.yml`-Job) sind gegen reale Dateien geprüft und berechtigt
   gesetzt; die vier verbleibenden `[ ]`-Checkboxen sind konsistent mit
   dem `in-progress`-Zustand (Closure-Pflichten, nicht Liefer-Punkte) —
   eine davon (Risiko-Ausgang) ist inhaltlich bereits erfüllt, das Häkchen
   folgt formal erst bei Closure.
2. F-1 ist real und vollständig behoben — der Kopfkommentar von
   `hub-description.yml:7` zitiert jetzt „[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 8", und
   diese Stelle trägt den zitierten Satz wörtlich und exklusiv (eigener
   `grep -c` gegen die volle ADR-Datei: Treffer `1`).
3. F-2 ist real und vollständig behoben — die DoD-Zeile zitiert jetzt
   „§6 Out-of-Scope der Welle-Datei", und dieser Abschnitt trägt die
   Aussage über `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN` wörtlich.
4. F-3 ist real behoben: `make test-dockerhub-token` läuft grün und
   netzlos; zwei eigene, vom Testskript unabhängige `stdin`-Proben gegen
   `dockerhub-token.sh` bestätigen korrektes Verhalten bei gültiger und
   ungültiger Eingabe; `hub-description.yml` nutzt das Skript tatsächlich
   für die Token-Extraktion, kein Inline-`jq` dafür mehr.
5. `make gates` lief eigenständig, ungepiped, mit `EXIT=0` auf `HEAD`
   (`5cc7fa49`); beide Slice-Commits tragen `ADR-0051` im Betreff (keine
   `SPEC-*`/`ARC-*`-Kennung), der Fixrunden-Commit reißt kein Gate.
6. Alle vier §6-Risiken tragen eine der drei zulässigen Ausgangsklassen;
   der neu als „eingetreten" geführte Scope-Hinweis (`read/write/delete`
   statt `read/write`, `403 Forbidden`) steht real im Kopfkommentar von
   `hub-description.yml`.
7. Die DoD-Checkbox „Review durchgeführt" ist inhaltlich korrekt gegen
   den tatsächlichen Report (2 HIGH/1 MEDIUM/1 LOW/2 INFO, identische
   Zuordnung) und berechtigt gesetzt.

**Nicht blockierende Beobachtung:** Die externe Quellenaussage „im
Schwester-Repo d-check real dokumentiert (`packaging/dockerhub/README.md`
§Transport)" konnte von dieser Sitzung aus nicht nachgeprüft werden (Datei
liegt außerhalb dieses Repos); die **Wirkung im eigenen Repo** (Kommentar
im Kopf von `hub-description.yml`) ist real bestätigt. Dem Planner zur
Kenntnis, kein DoD-Blocker.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz, Reconciliation-Register-Prüfung,
Beobachtungs-Register-Nachzug — insbesondere der neue
`BEO-PGC/zitat-nennt-die-falsche-stelle`-Beleg für F-1/F-2 —, formaler
Risiko-Ausgangs-Häkchen-Nachzug in §2, die drei Paarungen bei der
nächsten Welle-Closure).
