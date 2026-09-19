# Verifikationsbericht: slice-release-doku-releasing — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-doku-releasing.md` §2) und die
§6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-release-doku-releasing.md`](review-slice-release-doku-releasing.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** zwei Commits auf `main`, **letzter** Slice der Welle
`docs/plan/planning/welle-release-pipeline-adr-0051.md`:

- `b019b05823bfdc63cba376fbe746fff9198c839c` — ursprünglicher
  Implementer-Commit (`docs/user/releasing.md` neu, `harness/README.md`
  §Source-Precedence Zeile 6 und `AGENTS.md` §2 Zeile 6 auf echten Link
  umgestellt).
- `e045bd78dfc53d48942198e0326522ac2880917e` (= `HEAD`) — Fixrunde nach
  unabhängigem Review (1 HIGH F-1 „Zitat nennt die falsche Stelle" — `§3`
  statt `§4`, 1 INFO F-2), inklusive Review-Report selbst und neuem Beleg
  im Beobachtungs-Register
  (`BEO-PGC/zitat-nennt-die-falsche-stelle/evidence/slice-release-doku-releasing.md`).

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`, `AGENTS.md`,
`harness/conventions.md`, den vollständigen Slice-Plan (§1–§8), die
vollständige Welle-Datei, den vollständigen Review-Report und
`docs/user/releasing.md` selbst (beide Abschnitte §3 und §4 im Volltext)
gelesen. Nichts aus Slice-Plan, Commit-Message oder Review-Report
ungeprüft übernommen: eigene Lektüre von `releasing.md` §3/§4 im Original
zur F-1-Bestätigung; eigener `git show`-Diff beider Commits; eigene,
unabhängige Läufe aller sechs Tag-Beispiele gegen
`tools/harness/release-tag-info.sh`; eigene Lektüre von `release.yml`,
`hub-description.yml`, `image-scan.yml`, `upstream-drift.yml`,
`ci.yml`/`e2e.yml` (`tags-ignore`); eigener `ls`/`cat` gegen alle in
`releasing.md` genannten Dateien; eigener ungepipter `make gates`-Lauf
mit direkter Exit-Code-Prüfung (`AGENTS.md` §3.9); eigene isolierte
`make doc-commits`/`make doc-immutable`-Läufe über den exakten
Slice-Commit-Bereich; eigene `git tag -l`-Prüfung (kein realer Tag);
eigener `ls`-Beleg über `docs/user/` (sieben Dateien) und über
`docs/plan/planning/done/`+`in-progress/` (Welle-Vollständigkeits-Check).

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 „`docs/user/releasing.md` existiert … beschreibt real existierende Artefakte … alle sechs Tag-Beispiele verifiziert … Hinweis auf noch keinen realen Tag"

- **Datei existiert**, trägt Kopf `Version: 1.0` / `Stand: 2026-09-19`,
  strukturell analog `benutzerhandbuch.md` (`Version:`/`Stand:`-Zeilen,
  `### Änderungshistorie`-Tabelle mit denselben Spaltenköpfen
  `Version | Datum | Änderung`). Das Fehlen eines eigenen
  `Software-Version:`-Felds ist sachlich begründet — der Release-Prozess
  ist an keine einzelne Software-Version gebunden — und keine
  DoD-Abweichung, da die Zeile nur „analog" verlangt, nicht identisch.
- **Alle genannten Dateien/Targets real gegen den Baum geprüft** (eigene
  Prüfung, nicht übernommen):
  - `docs/user/version.md` existiert, Inhalt `0.1.0`.
  - `tools/harness/release-tag-info.sh` existiert.
  - `.github/workflows/{release,image-scan,upstream-drift,hub-description}.yml`
    existieren alle.
  - `README.md` existiert.
  - `make test-release-tag-info` und `make image-cve` sind reale
    Makefile-Ziele (`grep -n` gegen `Makefile` bestätigt beide Zeilen).
- **Alle sechs Tag-Beispiele eigenständig erneut gefahren** (nicht aus
  Review-Report übernommen):

  ```
  $ bash tools/harness/release-tag-info.sh "v0.1.0"       → version=0.1.0 latest=true, exit=0
  $ bash tools/harness/release-tag-info.sh "v1.2.0-rc.1"  → version=1.2.0-rc.1 latest=false, exit=0
  $ bash tools/harness/release-tag-info.sh "v2.0.0+build.5" → version=2.0.0+build.5 latest=true, exit=0
  $ bash tools/harness/release-tag-info.sh "v1.0.0-01"    → Fehler, exit=1
  $ bash tools/harness/release-tag-info.sh "v1.0"         → Fehler, exit=1
  $ bash tools/harness/release-tag-info.sh "1.0.0"        → Fehler, exit=1
  ```

  Deckt sich exakt mit dem in `releasing.md` §3 behaupteten Ergebnis (drei
  gültig mit `version=`/`latest=`, drei ungültig mit Exit 1).
- **Kein-Release-Hinweis:** `git tag -l` liefert eigenständig geprüft
  keine Ausgabe — kein realer Tag existiert. `releasing.md` §1 trägt
  wörtlich „Zum Zeitpunkt dieses Dokuments wurde noch kein realer
  Release-Tag gesetzt" — zutreffend.
- Inhaltlicher Abgleich §4/§5 gegen die vier realen Workflow-Dateien
  (eigene vollständige Lektüre, nicht nur Ausschnitte): Trigger
  (`push: tags: ['v*']`), Reihenfolge (Tag-Validierung →
  `version.md`-Abgleich → `make image VERSION=…`/`LATEST=…` →
  GitHub-Release mit Digest → `hub-description`-Job über
  `needs: release`/`uses: ./.github/workflows/hub-description.yml`/
  `secrets: inherit`), Docker-Hub-Token-Scope-Hinweis
  (`read/write/delete`, `403 Forbidden` bei `read/write`) und die
  advisory-Eigenschaft von `image-scan.yml`/`upstream-drift.yml` stimmen
  jeweils exakt überein. `ci.yml`/`e2e.yml` tragen real `tags-ignore:
  ['**']` — die Behauptung „kein Doppellauf" ist zutreffend.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 „`harness/README.md` Source-Precedence Zeile 6 … entfernt und auf echten Link umgestellt … `AGENTS.md` §2 Zeile 6 im selben Zug mitgezogen"

Eigener `git show b019b058 -- harness/README.md AGENTS.md`:

```
- | 6 | `docs/user/*` *(falls vorhanden)* | Operations, Quality, Releasing | <!-- d-check:ignore … -->
+ | 6 | [`docs/user/`](../docs/user/) | Operations, Quality, Releasing |
```

und analog in `AGENTS.md` §2 mit `[docs/user/](docs/user/)`. Eigener
`ls docs/user/` bestätigt sieben Dateien (`bench-abdeckung.md`,
`benutzerhandbuch-standard.md`, `benutzerhandbuch.md`,
`ci-matrix-abdeckung.md`, `e2e-abdeckung.md`, `releasing.md`,
`version.md`) — die im entfernten Kommentar genannte Bedingung „im
frischen Repo selten vorhanden" ist damit widerlegt, wie die DoD-Zeile
behauptet. Der Plan-Nachzug für `AGENTS.md` §2 steht sichtbar in Plan §3
als eigene Tabellenzeile, bevor die DoD-Zeile auf erledigt gesetzt wurde.
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 „`make gates` grün" — siehe eigenständiger Abschnitt 4 unten

### 1.4 „Review durchgeführt …" — siehe eigenständiger Abschnitt 6 unten

### 1.5 Die verbleibenden `[ ]`-Checkboxen (Doku-Update entfällt, Closure-Notiz, Reconciliation, Risiko-Ausgang, drei Paarungen)

Eigener Beleg: Der Plan liegt weiterhin unter
`docs/plan/planning/in-progress/release-doku-releasing.md`, §7
„Closure-Notiz" trägt weiterhin ausschließlich
`*(wird bei Bearbeitung gefüllt.)*`. Für einen Slice, der noch nicht nach
`done/` gewandert ist, ist das korrekt — diese Punkte sind
Closure-Pflichten, keine Liefer-Punkte; ihr `[ ]`-Zustand widerspricht
sich nicht mit dem übrigen DoD-Bild. Das Beobachtungs-Register-Item ist
demgegenüber bereits `[x]` (siehe Abschnitt 5 unten) — die Fixrunde hat es
vorgezogen erledigt, das ist zulässig und kein Widerspruch.

**Eine Beobachtung, kein Blocker:** Die Checkbox „Jedes Risiko aus §6
trägt einen Ausgang" ist unverändert `[ ]`, obwohl der einzige
§6-Risiko-Eintrag bereits seit Plan-Erstellung (Commit `06442e62`, vor
jeder Implementierung) den Text „**Ausgang:** weiter offen" trägt (siehe
Abschnitt 5 unten). Kein DoD-Verstoß, da der Slice noch nicht geschlossen
ist; dem Planner zur Kenntnis für den formalen Häkchen-Nachzug bei
Closure.

## 2. F-1-Fix real wirksam — `docs/user/releasing.md` §2 zitiert jetzt §4 statt §3

Eigene Lektüre des Originals **vor** dem Aufschlagen des Diffs, um die
Behauptung unabhängig zu prüfen:

- **§3 „Einen Release auslösen"** (Zeilen 41–60): beschreibt
  ausschließlich die Tag-Form `v<SemVer>`, ihre SemVer-2.0-Validierung
  über `release-tag-info.sh` und wer taggen darf. Kein Wort zu
  `docs/user/version.md` oder einem Abgleich dagegen.
- **§4 „Was beim Release automatisch passiert"** (Zeilen 62–107), Punkt 2:
  „`docs/user/version.md` gegen den Tag abgleichen — Abbruch bei
  Abweichung (§2)." — genau der Vorgang, den §2 Zeile 32 zitiert.

Eigener `git show e045bd78 -- docs/user/releasing.md`:

```
- Release-Workflow (§3) gleicht sie gegen den gesetzten Git-Tag ab und
+ Release-Workflow (§4) gleicht sie gegen den gesetzten Git-Tag ab und
```

Der zitierte Anker zeigt jetzt korrekt auf §4, die einzige Stelle, die den
Abgleich tatsächlich beschreibt. **Ergebnis: F-1 ist real und vollständig
behoben.**

## 3. F-2 — zur Kenntnis genommen, kein Code-Fix nötig

F-2 (INFO) verlangt keine Korrektur, sondern hält fest, dass
`docs/user/releasing.md` in der „Siehe:"-Liste von `README.md` fehlt.
Eigene Lektüre von `README.md` Zeilen 20–24 bestätigt: Die Liste verlinkt
weiterhin nur `benutzerhandbuch.md`, `spec/lastenheft.md`,
`spec/pflichtenheft.md`, `docs/plan/adr/` — `releasing.md`,
`version.md`, `bench-abdeckung.md`, `ci-matrix-abdeckung.md` und
`e2e-abdeckung.md` fehlen alle gleichermaßen. Vorbestehende Lücke,
korrekt als solche eingeordnet, kein neues Muster. **Ergebnis: F-2
zutreffend zur Kenntnis genommen, keine offene Handlung.**

## 4. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
e045bd78dfc53d48942198e0326522ac2880917e
$ git status --short
(leer)
$ make gates > /tmp/verifier-gates-release-doku-releasing.log 2>&1; ec=$?; echo "EXIT_CODE=$ec"
EXIT_CODE=0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 777 Datei(en) geprüft, 0 Befund(e)` (volle Modul-Liste) |
| `commit-traceability` | `d-check`-Modul `commits`: `777 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Commit-Traceability speziell für `e045bd78` geprüft** (Auftrag: „trägt
der Fixrunden-Commit eine `ADR-*`-ID?"):

```
$ git log -1 --format=%s e045bd78
fix(release): Reviewer-Fixrunde für release-doku-releasing (ADR-0051, 1 HIGH, 1 INFO)
$ git log -1 --format=%s b019b058
docs(user): releasing.md für den realen Release-Prozess (ADR-0051)
```

Beide Commits tragen `ADR-0051` im Betreff, keine `SPEC-*`/`ARC-*`-Kennung.
Zusätzlich isoliert über den exakten Slice-Commit-Bereich bestätigt
(Parent `8e0d8d0b` — letzter Commit vor `open → next` dieses Slice — bis
`e045bd78`):

```
$ make doc-commits RANGE=8e0d8d0b..e045bd78
d-check: 777 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=8e0d8d0b..e045bd78
d-check: 777 Datei(en) geprüft, 0 Befund(e)
```

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, `EXIT_CODE=0`, zusätzlich durch zwei
gezielte Modul-Läufe über den exakten Slice-Bereich bestätigt. Der
Fixrunden-Commit reißt kein Gate und trägt eine gültige `ADR-*`-Kennung.

## 5. §6-Risiko — zulässiger Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout/Folge-Slice · *entfallen* →
gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

Der einzige §6-Eintrag (Plan bereits seit Erstellungs-Commit `06442e62`,
vor jeder Implementierung, unverändert): „Eine Doku, die einen Prozess
beschreibt, der real noch nie durchlaufen wurde …" mit „**Ausgang:**
weiter offen, löst sich mit dem ersten echten Release — dann
Nachtrag/Korrektur als eigener kleiner Folge-Vorgang, kein Blocker für
diesen Slice."

**Ergebnis: Der einzige §6-Risiko-Eintrag trägt eine der drei zulässigen
Ausgangsklassen** (*weiter offen*) — konsistent mit der strukturell
identischen Klasse `AGENTS.md` §3.10 (Post-Push-Verifikationspflicht) und
mit dem analogen Risiko-Umgang in den vier bereits geschlossenen
Geschwister-Slices dieser Welle. Formaler Häkchen-Nachzug in §2 steht
noch aus (siehe Abschnitt 1.5), inhaltlich ist die Anforderung erfüllt.

## 6. DoD-Checkbox „Review durchgeführt" — inhaltlich korrekt gegen den Report?

Review-Report-Summary-Tabelle (eigene Lektüre): HIGH 1, MEDIUM 0, LOW 0,
INFO 1 — Findings F-1 (HIGH, Zitat §3 statt §4), F-2 (INFO, `releasing.md`
fehlt in README-„Siehe:"-Liste).

DoD-Checkbox-Text (Plan §2, Zeile 76–83) nennt exakt dieselbe Verteilung
(„1 HIGH … und 1 INFO …") mit denselben Kurzbeschreibungen je Finding —
Zahlen und Zuordnung stimmen 1:1 mit dem Report überein.

Eigene Prüfung, ob „behoben bzw. zur Kenntnis genommen" für jedes Finding
zutrifft:

- **F-1** (HIGH): real behoben, siehe Abschnitt 2 oben.
- **F-2** (INFO): korrekt als vorbestehende, nicht zu diesem Slice
  gehörende Lücke zur Kenntnis genommen, siehe Abschnitt 3 oben — der
  Plan (§1) nennt `README.md` nicht als Änderungsziel dieses Slice, kein
  Blocker.

**Kein offenes HIGH.** Die Fixrunde hat kein Gate gerissen (Abschnitt 4).
Der Review-Report selbst dokumentiert im Verdikt korrekt, dass er die
DoD-Checkbox wegen des HIGH-Fundes nicht selbst nachgezogen hat (Skill-
Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift nur bei 0 HIGH) — der
Implementer hat sie nach der Fixrunde manuell auf `[x]` gesetzt, was
diese Verifikation als inhaltlich berechtigt bestätigt.
**Ergebnis: Checkbox berechtigt auf `[x]` gesetzt.**

## 7. DoD-Checkbox „Beobachtungs-Register fortgeschrieben" — Evidence real, Zähler korrekt?

Eigener `find`-Beleg: die Datei
`docs/plan/planning/observations/BEO-PGC/zitat-nennt-die-falsche-stelle/evidence/slice-release-doku-releasing.md`
existiert real und beschreibt inhaltlich exakt F-1 (§3 → §4-Fehlzitat,
dritter Fundort dieser Klasse in drei aufeinanderfolgenden Slices
derselben Welle).

Eigene Lektüre von `state.md`: Zähler „**7×**" mit expliziter Aufzählung
aller sieben Belege:

1. `evidence/slice-090.md`
2. `evidence/slice-102.md`
3. `evidence/slice-d-check-tracked-modul.md`
4. `evidence/welle-d-check-verkoerperung.md`
5. `evidence/slice-release-hub-description.md` (zwei Fundstellen F-1/F-2,
   zählt als ein Beleg-Dateieintrag)
6. `evidence/slice-release-doku-releasing.md` (dieser Slice)

Nachrechnung: vier Einzelbelege (1–4) + zwei Fundstellen im
`release-hub-description`-Beleg (5) + ein Fund hier (6) = **7** — der
Zähler ist korrekt nachgezogen, nicht nur behauptet. Der Beobachtungs-
Status bleibt „**verkörpert**" (bereits bei `welle-d-check`-Closure
zugewiesen) — dieser Beleg löst keine erneute Verkörperung aus, was
`state.md` explizit so festhält. **Ergebnis: Checkbox berechtigt auf
`[x]` gesetzt.**

## 8. Wellen-Vollständigkeit nach dieser Closure

Eigener `ls`-Beleg:

```
$ ls docs/plan/planning/done | grep -i release
release-hub-description.md
release-image-scan.md
release-upstream-drift.md
release-version-und-workflow.md
$ ls docs/plan/planning/in-progress | grep -i release
release-doku-releasing.md
$ ls docs/plan/planning/next | grep -i release
(leer)
```

Vier der fünf in der Welle-Datei §4 gelisteten Slices
(`release-version-und-workflow`, `release-image-scan`,
`release-upstream-drift`, `release-hub-description`) liegen bereits in
`done/`; `release-doku-releasing` ist der einzige verbleibende, aktuell
`in-progress`. Nach seiner Closure (`git mv` nach `done/`) lägen **alle
fünf** Slices der Welle in `done/` — der Closure-Trigger der Welle-Datei
§3 („Alle fünf Slices in `done/`.") wäre dann erfüllbar, vorbehaltlich der
noch offenen Closure-Pflichten dieses Slice selbst (§7 Closure-Notiz,
Risiko-Häkchen, drei Paarungen — letztere ausdrücklich „von der nächsten
Welle-Closure" laut DoD-Zeile, was hier korrekt die anstehende
Wellen-Closure selbst ist).

---

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten in §2 sind bewusst noch offen, siehe 1.5). Alle sieben
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht aus
Bericht oder Commit-Message übernommen:

1. Beide `[x]`-Liefer-Punkt-Checkboxen (`releasing.md`,
   `harness/README.md`/`AGENTS.md`-Link-Umstellung) sind gegen reale
   Dateien und eigenständig wiederholte Sub-Prüfungen (sechs
   Tag-Beispiele, vier Workflow-Dateien Satz für Satz, `docs/user/`
   sieben Dateien) geprüft und berechtigt gesetzt.
2. F-1 ist real und vollständig behoben — `docs/user/releasing.md:32`
   zitiert jetzt „§4", die einzige Stelle, die den `version.md`-Abgleich
   tatsächlich beschreibt (§3 enthält keine `version.md`-Erwähnung,
   eigenständig durch Lektüre beider Abschnitte im Original bestätigt).
3. Alle in `releasing.md` genannten Dateien/Targets existieren real
   (`version.md`, `release-tag-info.sh`, vier Workflow-Dateien, `README.md`,
   `make test-release-tag-info`, `make image-cve`); alle sechs
   Tag-Beispiele eigenständig erneut gegen `release-tag-info.sh` gefahren,
   Ergebnis deckt sich exakt.
4. `make gates` lief eigenständig, ungepiped, mit `EXIT_CODE=0` auf `HEAD`
   (`e045bd78`); Commit-Traceability speziell für den Fixrunden-Commit
   bestätigt `ADR-0051` im Betreff, kein `SPEC-*`/`ARC-*`.
5. Der einzige §6-Risiko-Eintrag trägt die zulässige Ausgangsklasse
   *weiter offen*.
6. Die DoD-Checkbox „Review durchgeführt" ist inhaltlich korrekt gegen
   den tatsächlichen Report (1 HIGH/1 INFO, identische Zuordnung) und
   berechtigt gesetzt.
7. Die DoD-Checkbox „Beobachtungs-Register fortgeschrieben" ist korrekt:
   die Evidence-Datei existiert real, der Zähler in `state.md` ist
   nachweislich auf 7× nachgezogen (eigene Nachrechnung aller
   Einzelbelege).

Zusätzlich (Auftrag „letzter Slice der Welle"): Nach Closure dieses Slice
lägen real alle fünf Welle-Slices in `done/` — die
Wellen-Closure-Voraussetzung „Alle fünf Slices in `done/`" ist damit
erreichbar.

**Nicht blockierende Beobachtung:** Die DoD-Checkbox „Jedes Risiko aus §6
trägt einen Ausgang" bleibt formal `[ ]`, obwohl der Ausgang inhaltlich
bereits seit Plan-Erstellung feststeht — Häkchen-Nachzug ist eine
Closure-Formalie, kein DoD-Verstoß im aktuellen `in-progress`-Zustand.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz mit Steering-Loop-Lerneintrag, formaler
Risiko-Ausgangs-Häkchen-Nachzug in §2) — anschließend ist die Welle
`welle-release-pipeline-adr-0051` mit allen fünf Slices in `done/`
closure-bereit.
