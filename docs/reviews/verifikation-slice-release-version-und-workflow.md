# Verifikationsbericht: slice-release-version-und-workflow — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-version-und-workflow.md` §2) und
die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, mit
[`review-slice-release-version-und-workflow.md`](review-slice-release-version-und-workflow.md)
abgeschlossen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Gegenstand:** zwei Commits auf `main` (Welle
`docs/plan/planning/welle-release-pipeline-adr-0051.md`, `ADR-0051`
Entscheidung 1/2/3, Status `Accepted`):

- `5e3c7b68` — ursprünglicher Implementer-Commit (`version.md`, `make
  image VERSION=`, `release.yml`).
- `31cf6e8a` — Fixrunde nach unabhängigem Review (1 HIGH F-1, 4 MEDIUM
  F-2…F-5, 1 LOW F-6).

**Frischer Kontext:** Diese Sitzung hat Slice-Plan, Review-Report,
`release.yml`, `tools/harness/release-tag-info.sh`,
`tools/harness/run-release-tag-info-tests.sh`, den `Makefile`-`image`-
Ausschnitt, `ADR-0103` (vollständiger §Re-Evaluierungs-Trigger-Text) und
`ADR-0051` (Entscheidung 1–4, Verglichene Alternativen 1–2) selbst
gelesen. Nichts aus Slice-Plan, Commit-Message oder Review-Report wurde
ungeprüft übernommen: eigener Testlauf des bestehenden Zwölf-Fälle-Tests,
eine eigene, unabhängig entworfene Zusatz-Testdatei mit zwanzig weiteren
SemVer-2.0-Fällen (offizielle Spezifikationsbeispiele plus eigene
Randfälle), eigene Lektüre von `release.yml` gegen den behaupteten
Skript-Aufruf, eigener `make doc-commits`/`make doc-immutable`-Lauf über
den Slice-Commit-Bereich, eigener ungepipter `make gates`-Lauf mit
direkter Exit-Code-Prüfung (`AGENTS.md` §3.9).

---

## 1. `bash tools/harness/run-release-tag-info-tests.sh` real ausgeführt + eigene Zusatzfälle

Eigener Lauf des bestehenden Tests:

```
$ bash tools/harness/run-release-tag-info-tests.sh
run-release-tag-info-tests: alle Fälle bestanden
```

Exit `0`, alle zwölf Fälle bestanden (sechs `assert_ok`, sechs
`assert_invalid`) — deckungsgleich mit der DoD-Behauptung „real gegen
zwölf Fälle".

**Eigene, unabhängig entworfene Gegenprobe** (nicht aus dem bestehenden
Testfile übernommen, eigene Skriptdatei im Scratchpad, zwanzig Fälle):

| Klasse | Fälle | Ergebnis |
|---|---|---|
| Offizielle valide SemVer-2.0-Beispiele (semver.org) | `1.0.0-alpha`, `1.0.0-alpha.1`, `1.0.0-0.3.7`, `1.0.0-x.7.z.92`, `1.0.0-x-y-z.--`, `1.0.0-alpha+001`, `1.0.0+20130313144700`, `1.0.0-beta+exp.sha.5114f85`, `1.0.0+21AF26D3---117B344092BD`, `1.0.0-alpha.beta`, `1.0.0-alpha0.valid`, `1.0.0-rc.1+build.1` | alle zwölf korrekt als gültig akzeptiert, `latest` korrekt gesetzt (`false` bei Prerelease, `true` bei reiner Build-Metadata) |
| Vom Auftrag vorgegebene Randfälle | `v1.0.0-00` (rein-numerischer Identifier mit führender Null), `v1.0.0-a.01` (zweiter Identifier rein-numerisch mit führender Null), `v1.0.0-0` (einzelnes `0`, kein Leading-Zero-Fall) | `v1.0.0-00` und `v1.0.0-a.01` korrekt **abgelehnt**; `v1.0.0-0` korrekt **akzeptiert** (`latest=false`) |
| Weitere eigene Randfälle | `v1.0.0-alpha..1` (leeres Identifier-Segment), `v1.0.0-alpha_beta` (Unterstrich), `v1.0.0-` (Bindestrich ohne Identifier), `v1.0.0+` (Plus ohne Build-Metadata), `v1.0.0+01` (Build-Metadata mit führender Null — laut Spec erlaubt, nur Prerelease-Identifier verbieten sie), `v1.0.0-01` (F-2-Gegenprobe), `v1.0.0+exp-sha.5114f85` (F-3-Gegenprobe), `v1.0.0+exp-sha.5114f85-01` (Bindestrich **und** Ziffernfolge in der Build-Metadata, kein Prerelease) | alle acht korrekt bewertet |

Alle zwanzig eigenen Zusatzfälle liefen wie erwartet — Skript-Exit `0`
(„ALLE ZUSATZ-FAELLE BESTANDEN"). **Ergebnis: Die Korrektur trifft nicht
nur die zwei ursprünglich im Review gefundenen Fälle
(`1.0.0-01`/`1.0.0+exp-sha.5114f85`), sondern die zugrunde liegende
SemVer-2.0-Regel selbst** — sowohl die „keine führende Null bei
rein-numerischen Prerelease-Identifiern"-Regel (§9) als auch die
„Build-Metadata zählt nie als Prerelease, unabhängig von ihrem Inhalt
(auch Bindestriche, auch Ziffernfolgen)"-Regel (§10) sind an mehreren,
unabhängig gewählten Stellen bestätigt, nicht nur an den zwei
Fixrunden-Testfällen.

`make test-release-tag-info` wurde ebenfalls geprüft (`Makefile:46-48`,
`.PHONY: test-release-tag-info` → `bash
tools/harness/run-release-tag-info-tests.sh`) — ruft exakt dasselbe
Skript auf, kein zweiter Codepfad.

## 2. `release.yml` ruft real `tools/harness/release-tag-info.sh` auf

Eigene Lektüre der aktuellen Datei (`.github/workflows/release.yml:54-55`):

```yaml
- name: Tag validieren, Version und Stabilitaet ermitteln (tools/harness/release-tag-info.sh)
  run: bash tools/harness/release-tag-info.sh "${GITHUB_REF_NAME}" >> "$GITHUB_ENV"
```

Keine Inline-SemVer-Regex, kein Inline-`case`-Ausdruck mehr im Workflow
selbst — die gesamte Tag-Validierung/Stabilitäts-Ermittlung ist
vollständig in das externe, unter Punkt 1 geprüfte Skript ausgelagert.
Die nachfolgenden Schritte (`docs/user/version.md`-Abgleich, `make image
VERSION="$version" LATEST="$latest"`, GitHub-Release-Anlage) lesen `$version`/
`$latest` ausschließlich aus `$GITHUB_ENV`, das genau dieser Schritt
befüllt — kein zweiter, konkurrierender Ermittlungspfad im Diff
gefunden. **Ergebnis: F-2/F-3 sind strukturell behoben** (nicht nur
patched an der alten Stelle, sondern die fehlerhafte Logik existiert im
Workflow gar nicht mehr).

## 3. F-1s Behebung — DoD-Checkboxen tragen jetzt in sich geschlossene Belege; §7 weiterhin korrekt Platzhalter

Eigene Lektüre von
`docs/plan/planning/in-progress/release-version-und-workflow.md` §2
(aktueller Stand, `31cf6e8a`):

- Kein „siehe §7" mehr in den beiden beanstandeten Checkbox-Zeilen. Die
  `make image`-Checkbox trägt den Beleg jetzt direkt im Text: „real gegen
  zwei lokale Registry-Container geprüft (`registry:2`, zwei Ports):
  derselbe Manifest-Digest auf allen vier Tag/Registry-Kombinationen
  (`docker buildx imagetools inspect`, Beleg im Fixrunden-Commit)".
- Die `release.yml`-Checkbox trägt ebenfalls einen in sich geschlossenen
  Beleg: „YAML-Struktur und jeder `run:`-Schritt syntaktisch geprüft
  (Ruby-Stdlib-YAML-Parser bzw. `bash -n`) — der reale Post-Push-Lauf
  bleibt unverifiziert (`AGENTS.md` §3.10, siehe §6)". Der einzige
  verbliebene Verweis auf eine andere Sektion (`§6`) betrifft explizit
  **nicht** die hier behauptete Verifikation selbst (die steht direkt im
  Satz davor), sondern die separat und korrekt als offen deklarierte
  Post-Push-Frage — das ist kein Fall des F-1-Musters „Beleg trägt seinen
  Satz nicht", weil der Satz, auf den verwiesen wird, nichts über die
  hier behauptete Prüfung selbst aussagt, sondern eine andere,
  ausdrücklich offene Eigenschaft.
- §7 „Closure-Notiz" trägt weiterhin ausschließlich `*(wird bei
  Bearbeitung gefüllt.)*` — für diesen Zwischenstand (Slice weiterhin
  `in-progress`, keine Closure) korrekt: die drei zugehörigen
  Closure-Pflicht-Checkboxen (§2, „Closure-Notiz", „Beobachtungs-Register",
  „Jedes Risiko aus §6 trägt einen Ausgang", „Die drei Paarungen") stehen
  konsistent ebenfalls noch auf `[ ]`.

**Ergebnis: F-1 ist real behoben, nicht nur behauptet behoben.**

## 4. F-5s Behandlung in §6 — Argumentation gegen `ADR-0103`s eigenen Wortlaut geprüft

Eigene Lektüre von `ADR-0103` §Re-Evaluierungs-Trigger im Volltext:

> „Beobachtbarer Trigger: **ein Lauf-Zweig braucht einen historischen,
> über Commits hinweg nachschlagbaren Image-Beleg** — etwa ein `docker
> push`-Workflow (Archiv-Form wird real erfüllt) oder ein
> Replay-Feature nach Modul 12 (…). Dann wird der Beleg als Folge-ADR mit
> `supersedes` wieder gehoben (…). Sonst `permanent`."

Zwei tragende Beobachtungen aus eigener Lektüre, die die Slice-Argumentation
stützen:

1. **Der neue `--push`-Codepfad selbst braucht keinen
   über-Commits-nachschlagbaren Beleg.** Der `ifdef VERSION`-Zweig in
   `Makefile:32-37` schreibt `harness/image-hash.txt` genauso lokal und
   nicht-committet wie der bestehende `:dev`-Pfad (`Makefile:39-40`) — der
   Digest-Beleg für einen `VERSION`-Lauf ist der GitHub-Release-
   Beschreibungstext (`release.yml:84-98`), ein Mechanismus **außerhalb**
   von `git`-Commit-Historie. Es entsteht durch diesen Slice kein
   Code-Zweig, der tatsächlich `harness/image-hash.txt` über Commits
   hinweg zurückverfolgen *müsste* — das ist exakt die Eigenschaft, die
   `ADR-0103` ursprünglich für unnötig erklärt hatte, und sie bleibt für
   diesen neuen Pfad ebenso unnötig.
2. **Die Welle-Datei selbst trennt Mechanismus von Anwendung so scharf,
   dass „Archiv-Form wird real erfüllt" konsequent erst mit einem
   tatsächlichen Push eintreten kann:** `welle-release-pipeline-adr-0051.md`
   §6 sagt wörtlich „Diese Welle liefert den **Mechanismus**, nicht seine
   erste Anwendung" — dieselbe Unterscheidung, die die Slice-Argumentation
   für den Trigger heranzieht (,Fähigkeit existiert' vs. ,braucht real'),
   ist bereits an anderer Stelle der Welle unabhängig etabliert, nicht neu
   erfunden für F-5.

Beide Punkte stützen die im Slice-Plan gewählte Lesart: „braucht real"
liest sich stimmig als *ein tatsächlich durchgeführter Push*, nicht als
*die bloße Existenz eines Workflows, der pushen könnte* — der Trigger-Text
selbst enthält keine Formulierung wie „ein Workflow **existiert**", sondern
„ein Lauf-Zweig **braucht**" und „Archiv-Form wird real **erfüllt**", beides
Zustands- nicht Fähigkeits-Aussagen.

**Zugleich bleibt dies — wie der Review-Report selbst sagt — eine
auslegungsbedürftige Formulierung, keine mechanisch entscheidbare Frage.**
Der Slice-Plan tut das Richtige: Er löst die Frage nicht eigenmächtig auf,
benennt sie als offene Architect-Frage mit einer konkreten, eigenen
Fälligkeit („spätestens vor dem ersten echten Release") und nennt zwei
gangbare Auflösungen (Folge-ADR mit `supersedes`, oder eine explizite
Begründung, warum der GitHub-Release-Text als Beleg im Sinn des Triggers
ausreicht). Das ist die korrekte Verifier-/Architect-Arbeitsteilung
(„ich flagge den Befund, ohne ihn selbst aufzulösen", Review-Report F-5).

**Ergebnis: Die Argumentation ist in sich schlüssig, textlich gegen
`ADR-0103`s eigenen Wortlaut und gegen die Welle-Datei konsistent, und die
Einordnung als „weiter offen" statt Blocker ist angemessen** — ein Blocker
würde bedeuten, dass kein Slice dieser Welle vor einer Architect-Antwort
schließbar wäre, obwohl der Trigger nach der plausibelsten Lesart erst mit
einem tatsächlichen, außerhalb dieser Welle liegenden Release feuert.

## 5. `make gates` real, ungepiped ausgeführt

Eigener Lauf auf aktuellem `HEAD` (`31cf6e8a`), Ausgabe in eine Datei
umgeleitet, Exit-Code separat und direkt geprüft (`AGENTS.md` §3.9):

```
$ make gates > gates.log 2>&1
$ echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 764 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check`-Modul `commits`: `764 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Anmerkung zur Coverage-Zahl (Herkunfts-Disziplin, `AGENTS.md` §3.12):**
Der Review-Report nennt für denselben Commit `31cf6e8a` `82.80%`, diese
eigene, unabhängige Messung liefert `82.70%` — beide Werte liegen über der
80-%-Schwelle, das Gate ist in beiden Fällen grün. Die kleine Differenz
(0,1 Prozentpunkte) ist nicht durch eine Code-Änderung erklärbar (kein
Commit zwischen den beiden Läufen), plausibelste Ursache ist eine
Emulations-Randbedingung des lokalen Build-Hosts (`WARNING: the requested
image's platform (linux/amd64) does not match the detected host platform
(linux/arm64/v8)`, beide Läufe liefen unter dieser Warnung). **Kein
DoD-Blocker**, da beide Werte die Schwelle erfüllen; als Beobachtung
festgehalten, nicht als Befund.

Zusätzlich eigenständig ausgeführt, wie in der Rollenvorgabe verlangt:

```
$ make doc-commits RANGE=75d17987..31cf6e8a
d-check: 764 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=75d17987..31cf6e8a
d-check: 764 Datei(en) geprüft, 0 Befund(e)
```

Beide isolierten Modul-Läufe (`commits`, `vcs`) über exakt den
Slice-Commit-Bereich (Parent `75d17987` vor `docs/plan/planning`-Move bis
`31cf6e8a`) bestätigen: keine Commit-Traceability-Lücke, keine
ADR-/Doc-Immutabilitätsverletzung in diesem Slice-Umfang.

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, mit `EXIT=0` bestätigt, zusätzlich
durch zwei gezielte Modul-Läufe über den Slice-Bereich ergänzt.

## 6. DoD-Checkbox „Review durchgeführt" — berechtigt gesetzt?

Ausgangslage: Der Review-Report fand 1 HIGH (F-1) und 4 MEDIUM
(F-2…F-5); sein eigenes Verdikt zog die Checkbox bewusst **nicht** nach
(„Da mindestens ein HIGH-Finding eine Implementer-Fixrunde auslöst, ziehe
ich die DoD-Checkbox … nicht nach").

Eigene, unabhängige Prüfung des Fixrunden-Commits `31cf6e8a` gegen jedes
Finding:

- **F-1** (HIGH): real behoben, siehe §3 oben — kein „siehe §7"-Verweis
  mehr, Belege direkt im Checkbox-Text.
- **F-2/F-3** (MEDIUM): real behoben, siehe §1/§2 oben — nicht nur an den
  zwei ursprünglich gefundenen Stellen, sondern strukturell (eigenständiges
  Skript, gegen zwanzig eigene Zusatzfälle bestätigt).
- **F-4** (MEDIUM): real behoben — `tools/harness/run-release-tag-info-tests.sh`
  existiert, deckt zwölf Fälle ab, läuft grün (§1).
- **F-5** (MEDIUM): bewusst nicht in der Fixrunde aufgelöst, sondern
  korrekt als offenes, architect-pflichtiges Risiko in §6 übernommen
  (siehe §4 oben) — das ist die vom Review selbst vorgesehene Behandlung
  für diesen Punkt, kein übersehenes MEDIUM.
- **F-6** (LOW): real behoben — die `harness/README.md`-Checkbox steht auf
  `[x]`, und `harness/README.md` trägt tatsächlich die neuen `make
  image`/`release.yml`-Zeilen (eigene `grep`-Bestätigung,
  `harness/README.md:126`/`:138`).

**Kein offenes HIGH, kein unbehandeltes MEDIUM.** Die Fixrunde behandelt
alle fünf adressierbaren Findings (F-1 bis F-4, F-6) durch reale
Behebung und das eine Architect-Finding (F-5) durch korrekte
Risiko-Übernahme statt stiller Auflösung.

**Ergebnis: Die Checkbox ist berechtigt auf `[x]` gesetzt.**

## 7. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout oder Folge-Slice mit ID · *entfallen*
→ gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | `release.yml` real erst nach echtem Tag-Push beweisbar | „weiter offen, strukturell (…), dokumentiert analog `BEO-PGC/github-actions-unverifizierbar-lokal`" | ✓ weiter offen | Deckt sich mit `AGENTS.md` §3.10 (dreifach bereits belegtes Muster) — korrekt kein neuer Registereintrag, sondern derselbe bestehende Fall. |
| 2 | `--push` ändert den Digest-Ermittlungspfad ggü. `ADR-0103`/`ADR-0044` | „eingetreten, real geprüft und unauffällig" mit Beleg (`--metadata-file` liefert in beiden Modi denselben Digest) | ✓ eingetreten | Klasse formal korrekt; Begründung nennt einen konkreten Prüfweg (`docker buildx imagetools inspect` gegen vier Tag/Registry-Kombinationen). Nicht Bestandteil dieses Verifikationsauftrags, eigenständig nachzurechnen (Registry-Container liefen zum Zeitpunkt dieser Verifikation nicht mehr) — die Begründung ist in sich plausibel und deckt sich mit derselben Verifikationsklasse wie F-1s ursprünglich beanstandeter (jetzt behobener) Beleg-Anker. |
| 3 | Docker Hub braucht andere Namensform als GHCR | „eingetreten, real geprüft: `docker.io/<user>/<repo>` … löst identisch … auf" | ✓ eingetreten | Klasse formal korrekt, Begründung plausibel (Docker-Tooling normalisiert `docker.io/`-Präfix bekanntermaßen); reale Kontoexistenz auf Docker Hub bleibt explizit und korrekt als separat unbewiesen benannt. |
| 4 (F-5) | `ADR-0103`-Re-Evaluierungs-Trigger sachlich berührt | „weiter offen, bewusst nicht in diesem Slice aufgelöst … Architect-Entscheidung … vor dem ersten echten Release" | ✓ weiter offen | Siehe eigene, ausführliche Prüfung in §4 oben — Argumentation schlüssig, Einordnung als offen statt Blocker angemessen. |

**Ergebnis: Alle vier Risiken tragen eine der drei zulässigen
Ausgangsklassen**, mit nachvollziehbaren, für diesen Slice-Typ (Kern-
Mechanismus ohne realen Release) stimmigen Begründungen.

---

## Verdikt

**DoD erfüllt.** Alle sieben beauftragten Prüfpunkte wurden real und
unabhängig nachgemessen, nicht aus Bericht oder Commit-Message
übernommen:

1. `bash tools/harness/run-release-tag-info-tests.sh` lief real, alle
   zwölf Fälle bestanden; zwanzig eigene Zusatzfälle (offizielle
   SemVer-2.0-Spec-Beispiele plus eigene Randfälle) bestätigen, dass die
   Korrektur die zugrunde liegende SemVer-2.0-Regel trifft, nicht nur die
   zwei ursprünglich gefundenen Instanzen.
2. `release.yml` ruft real `tools/harness/release-tag-info.sh` auf; kein
   Rest der alten Inline-Logik im Diff auffindbar.
3. F-1 ist real behoben — beide DoD-Checkboxen tragen jetzt in sich
   geschlossene Belege statt eines Verweises auf die leere §7; §7 bleibt
   für diesen Zwischenstand korrekt der unausgefüllte Platzhalter.
4. F-5s Behandlung in §6 ist gegen `ADR-0103`s eigenen Trigger-Wortlaut
   und die Welle-Datei geprüft und schlüssig; die Einordnung als „weiter
   offen" statt Blocker ist angemessen, mit korrekt benannter Fälligkeit
   vor dem ersten echten Release.
5. `make gates` lief eigenständig, ungepiped, mit `EXIT=0`; zusätzlich
   `make doc-commits`/`make doc-immutable` über den exakten
   Slice-Commit-Bereich, beide `0 Befund(e)`.
6. Die DoD-Checkbox „Review durchgeführt" ist berechtigt gesetzt: alle
   fünf adressierbaren Findings real behoben, das eine Architect-Finding
   (F-5) korrekt als Risiko übernommen statt stillschweigend aufgelöst,
   kein offenes HIGH.
7. Alle vier §6-Risiken tragen eine der drei zulässigen Ausgangsklassen.

**Nicht blockierende Beobachtung:** Die Coverage-Zahl driftet minimal
zwischen dem Review-Lauf (`82.80%`) und dieser eigenen Messung
(`82.70%`) auf demselben Commit, plausibel durch Plattform-Emulation
(`arm64`-Host, `amd64`-Image) erklärt — beide Werte liegen über der
80-%-Schwelle, kein Gate-Risiko.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die bereits im Plan selbst als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz, Beobachtungs-Register, formaler
Risiko-Ausgangs-Nachzug, die drei Paarungen) sowie — spätestens vor dem
ersten tatsächlichen Release, nicht vor dieser Slice-Closure — die
Architect-Entscheidung zu F-5/`ADR-0103`.
