# Verifikations-Bericht: slice-d-check-trace-rtm — 2026-09-17

**Rolle:** Verifier (Baseline-Regelwerk `v6.9.0`
`regelwerk/modul-11-verification.md`) — „Bauen wir es richtig?" gegen **DoD und
Entscheidungen**, nicht gegen Plan-Maintainability (das war der Reviewer).

**Gegenstand:** `slice-d-check-trace-rtm` (`d-check`s `trace:`-Block
verdrahtet: `trace.requirements.id-pattern` + `trace.coverage`, Welle
`welle-d-check`, zweiter und letzter Slice der Welle). Diff
`9562042..HEAD`, neun Commits: `d6f7ccc` (Ruhe-Marker zurück) · `0e77982`
(open→next) · `9f6a47d` (Verantwortlich gesetzt) · `31f6e56` (next→
in-progress) · `459d132` (Ruhe-Marker entfernt) · `f02a7b5` (`trace:`-Block
verdrahtet) · `b6daee6` (Träger-Nachzug) · `fce746c` (DoD-Checkboxen
LP1–LP3 + Doku-Update) · `d2ee520` (Review, 0 HIGH/0 MEDIUM/0 LOW/1 INFO).

**Plan:** `docs/plan/planning/in-progress/slice-d-check-trace-rtm.md`
(vollständig, §1–§8, unverändert seit `fce746c`).

**Entscheidungen:** keine eigene ADR — der Slice ist bewusst ADR-frei
angelegt (§`Bezug`: „reine `d-check`-Konfigurationserweiterung ohne
Vertrags- oder Entscheidungsgegenstand"). Kein Architect-Zug in diesem
Slice nötig oder erfolgt.

**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17.

**Eigener Eingabe-Kontext** (frisch, ohne Implementer-/Reviewer-Zusage):
Plan §1–§8 · Review-Report
([`review-slice-d-check-trace-rtm.md`](review-slice-d-check-trace-rtm.md))
· `AGENTS.md` §3.3, §3.9, §3.12 · `.d-check.yml` · `harness/README.md`
§Werkzeuge · `harness/sensors/docs-check.md` · `git log`/`git show` über
`9562042..HEAD`. Alle Messungen dieses Berichts sind **eigene** Läufe
(Exit-Codes ungepiped, `AGENTS.md` §3.9), keine Übernahme der Implementer-/
Reviewer-Behauptung.

---

## A. Eigene Sensoren (Belege, keine Behauptungen)

| Sensor | Kommando | Ergebnis | Exit |
|---|---|---|---|
| `make doc-trace` (Bestand, mit `trace.coverage`) | `make doc-trace` (ungefiltert, Exit-Code direkt geprüft) | volle RTM, 76 Zeilen; **76 Anforderung(en), 7 Waise(n)** | **0** |
| `make doc-trace` (temporär ohne `trace.coverage`-Block, danach vollständig restauriert) | `.d-check.yml`s `coverage:`-Unterblock per Skript entfernt, Lauf, dann Original-Datei zurückkopiert | **76 Anforderung(en), 9 Waise(n)** — Differenzmenge zu oben: genau `LH-FA-CAP-002`, `LH-FA-CAP-003` wechseln von `WAISE` zu `ok` | **0** |
| `git status --porcelain` nach Restaurierung | — | leer — kein Rest-Diff durch die Probe | — |
| `grep -c '^### LH-' spec/lastenheft.md` / Muster-Gegenprobe | `grep -oE '^### LH-[A-Za-z0-9.-]+' spec/lastenheft.md \| grep -vE '^### LH-(FA\|QA)-[A-Z]{3}-[0-9]{3}$'` | 76 Treffer gesamt, **0** Nicht-Treffer gegen `id-pattern` | — |
| `make doc-complete` (Docker-Tool direkt, ohne Make-Wrapper) | `docker run … --trace --require-complete` | Exit des **Tools** | **1** |
| `make doc-complete` (Make-Wrapper) | `make doc-complete` | GNU-Make-Fehlschlag-Exit des **Wrappers** | **2** |
| 7 reale Waisen gegen `docs/user/e2e-abdeckung.md` | `grep -c "<ID>" docs/user/e2e-abdeckung.md` je ID | 0 Treffer für alle sieben (`LH-FA-CFG-006`, `LH-FA-CON-002`, `LH-FA-DAT-002`, `LH-FA-DAT-003`, `LH-FA-SST-001`, `LH-FA-SST-005`, `LH-QA-REL-004`) — echte Lücken, keine übersehene Zitierung | — |
| `make doc-commits RANGE=9562042..HEAD` | `--enable commits` | 885 Dateien, 0 Befunde | **0** |
| `make doc-immutable RANGE=9562042..HEAD` | `--enable vcs` | 885 Dateien, 0 Befunde | **0** |
| Commit-Betreffs aller neun Slice-Commits, einzeln gelesen (`git show -s --format=%s`) | — | alle neun tragen `ADR-0045`, keine `SPEC-*`/`ARC-*`-Kennung | — |
| `git show --stat` je `git mv`-Commit (`0e77982`, `31f6e56`) | — | beide 0 Insertions/Deletions — reine Moves | — |
| `grep -n "GATE_CHECKS" harness/mk/*.mk Makefile` und `grep -n "doc-trace\|doc-complete" d-check.mk` | — | `doc-trace`/`doc-complete` sind Targets in `d-check.mk`, an **keiner** Stelle an `GATE_CHECKS` gehängt | — |
| `make gates` (repo-weit, HEAD `d2ee520`) | `make gates` (ungefiltert, Exit-Code direkt geprüft) | `baseline-verify` 54 Dateien OK · `docs-check` 885 Dateien/0 Befunde · `commit-traceability` 5 Commits OK · `coverage-gate` 83,40 % ≥ 80 % · `generated-sync` byte-gleich · `a-check` 0 Befunde | **0** |
| `git status --porcelain` nach `make gates` | — | leer | — |

---

## B. DoD-Konformität Punkt für Punkt

### LP1 — `trace.requirements.id-pattern` verdrahtet, `make doc-trace` läuft real

**Erfüllt.** `.d-check.yml:250-251`: `trace.requirements.id-pattern:
'LH-(FA|QA)-[A-Z]{3}-\d{3}'`. Eigener Lauf am HEAD (nicht die
Planungs-Vorabmessung übernommen, `AGENTS.md` §3.12): `make doc-trace`
Exit 0, volle RTM auf stdout, **76 Anforderung(en)**. Das Pattern deckt
alle 76 `### LH-*`-Überschriften in `spec/lastenheft.md` ohne
Nicht-Treffer — eigene, unabhängige Gegenprobe (§A) bestätigt die
Plan-Behauptung exakt.

### LP2 — `trace.coverage` ergänzt, Waisenzahl sinkt real geprüft

**Erfüllt.** `.d-check.yml:252-257`: `trace.coverage: [{files:
[docs/user/e2e-abdeckung.md], label: E2E}]`. Eigene, isolierte
Doppelmessung (§A): mit Block **7 Waisen**, nach temporärer Entfernung des
Blocks (danach vollständig restauriert, `git status` bestätigt leeren
Diff) **9 Waisen** — die Differenzmenge sind exakt `LH-FA-CAP-002` und
`LH-FA-CAP-003`, beide mit `E2E`-Nachweis in der RTM-Spalte. Der
Vergleich beider Läufe ist damit nicht nur behauptet, sondern selbst
reproduziert und die Ursache (E2E-Coverage-Dimension) bestätigt.

### LP3 — Träger-Nachzug (`harness/README.md`, `harness/sensors/docs-check.md`)

**Erfüllt.** `harness/README.md:129` trägt eine neue Zeile für
`make doc-trace` in der „kein Gate"-Werkzeuge-Tabelle (nicht in der
Gates-Tabelle), mit den real gemessenen Zahlen (76/9 ohne, 76/7 mit
`trace.coverage`), den sieben benannten Waisen, dem expliziten Vermerk
„`make doc-complete`/`--require-complete` bleibt bewusst **nicht** in
`GATE_CHECKS`/`make gates`" und dem Anker `· seit
slice-d-check-trace-rtm`. `harness/sensors/docs-check.md:135-146` (Punkt
10) hält fest, dass `trace:` kein Modul ist und nicht in
`docs-check`/`make gates` läuft, samt der `id-pattern`-Abweichung vom
Default. Eigene Prüfung von `d-check.mk` und allen `harness/mk/*.mk`
bestätigt: `doc-trace`/`doc-complete` hängen an **keiner** Stelle an
`GATE_CHECKS` (§A) — die Plan-Aussage „`doc-complete` ist nicht in
`GATE_CHECKS`/`make gates` aufgenommen" ist am Code verifiziert, nicht nur
behauptet.

### Review durchgeführt, Report liegt vor

**Erfüllt.** Report `d2ee520` liegt unter `docs/reviews/
review-slice-d-check-trace-rtm.md` vor: 0 HIGH/0 MEDIUM/0 LOW, 1 INFO,
nicht merge-blockierend, DoD-Checkbox-Nachzug ohne Fixrunde (Reviewer-Skill
§DoD-Checkbox-Nachzug). Der Report belegt eigene Gegenproben (isolierte
`trace.coverage`-Entfernung, `grep`-Gegenproben, `make gates`) — eigene
Nachmessung dieses Berichts (§A) reproduziert dieselben Zahlen unabhängig
und bestätigt sie exakt, kein Abweichungsfall gefunden.

### Doku-Update

**Erfüllt.** Deckt sich mit LP3 (dieselben zwei Träger). Kein weiterer
öffentlicher Vertrag berührt — `.d-check.yml`s eigener Kommentarblock
(`:245-249`) dokumentiert den Zweck des `trace:`-Blocks inline.

### `make gates` grün

**Unchecked im Plan — Faktenlage bestätigt „grün".** Die DoD-Checkbox
selbst ist im committeten Plan-Stand (`fce746c`) **nicht** angehakt; die
Commit-Message von `fce746c` benennt dies ausdrücklich als bewusst offen
gelassen für „nachgelagerte Rollen". Eigener, unabhängiger Lauf am HEAD
`d2ee520`: **Exit 0**, alle sechs Gate-Bestandteile grün (§A) — die
zugrunde liegende Tatsache trifft zu, die Checkbox ist nur noch nicht
formal nachgezogen. Kein DoD-Verstoß: Der Plan selbst weist diesen Punkt
(zusammen mit Closure-Notiz, Beobachtungs-Register, Risiko-Ausgängen,
drei Paarungen) den nachgelagerten Rollen zu — die welle-weite
Zusatz-Bedingung (kombinierter Lauf mit `slice-d-check-tracked-modul`,
§Welle im Plankopf) ist ohnehin nicht Teil dieser Einzel-DoD, sondern der
Welle-Closure vorbehalten.

### Closure-Notiz mit Steering-Loop-Lerneintrag

**Offen — korrekt unchecked.** §7 trägt ausschließlich Platzhalter
(`<…>`). Erwartungsgemäß, da der Slice erst jetzt zur Closure ansteht.

### Beobachtungs-Register fortgeschrieben

**Offen — korrekt unchecked.** Eigene, unabhängige Wiederholung des im
Plan §8 dokumentierten Such-Laufs: `grep -rli
"trace\|rtm\|traceability.matrix\|waise" docs/plan/planning/observations/`
liefert zehn Treffer, ausnahmslos False Positives (Datei- oder
Verzeichnisnamen anderer Beobachtungen wie
`spiegelung-ist-approximation`, `commit-traceability-kein-vorab-hook`,
Substring-Zufallstreffer) — **kein** Treffer zu genau diesem Gegenstand
(RTM-Konfiguration, Waisen-Anforderungen). Die im Plan §2/§8 gestellte
Frage („eigene Beobachtung für die 7 realen Waisen erwägen") ist noch
nicht beantwortet — Planner-Wertung, hier nicht vorweggenommen.

### Jedes Risiko aus §6 trägt einen Ausgang

**Offen — korrekt unchecked, siehe §C unten.** Alle drei Risiken tragen im
committeten Plan-Stand weiterhin ihre Platzhalter-Form
(`<weiter offen: … / entfallen: …>`). Die Fakten für alle drei liegen
jedoch bereits vor (§C).

### Drei Paarungen

**Korrekt unchecked, delegiert.** Plan-Text weist dies explizit der
Welle-Closure `welle-d-check` zu — konsistent mit dem repo-weiten
Wellen-Betrieb, wie beim Schwester-Slice.

---

## C. §6-Risiko-Ausgänge — eigene Prüfung

### Risiko 1 — „`trace.coverage`s `files:`-Liste veraltet bei Umbenennung/Verschiebung von `docs/user/e2e-abdeckung.md`"

**Faktenlage stützt „weiter offen" (Kopplungs-Vermerk), nicht „entfallen".**
`docs/user/e2e-abdeckung.md` ist laut `harness/README.md` §Werkzeuge
weiterhin das **Erzeugnis** von `make test-integration`
(`tools/harness/run-integration-tests.sh` schreibt es) — der Pfad bleibt
eine Kopplung an einen anderen Lauf, keine vom Plan selbst kontrollierte
Konstante. Ein Umbenennungs-Anlass ist zwar aktuell nicht erkennbar, aber
das ist keine strukturelle Garantie gegen künftige Drift. Plan-Text ist
noch Platzhalter.

### Risiko 2 — „Waisenzahl (9→7) verschiebt sich beim vollen Implementer-Lauf"

**Faktenlage stützt „entfallen".** Eigene Messung (§A): der reale
Implementer-/Verifier-Lauf am vollen `.d-check.yml`-Bestand liefert exakt
dieselben Zahlen wie die isolierte Planungsmessung (76/9 ohne, 76/7 mit
`trace.coverage`) — keine Abweichung durch andere Module (`matrix`, `ids`
etc.) im vollen Bündel. Plan-Text ist noch Platzhalter.

### Risiko 3 — „Advisory-Status von `--trace` wird künftig als durchgesetztes Gate fehlgelesen"

**Faktenlage stützt „entfallen".** `harness/README.md:129` trägt den
expliziten Vermerk „kein Gate, wie `make image-stale`" plus die konkrete
Nennung, dass `doc-complete`/`--require-complete` nicht in
`GATE_CHECKS`/`make gates` steht; `harness/sensors/docs-check.md` Punkt 10
wiederholt dies strukturell gleichlautend an zweiter Stelle. Zwei
unabhängige, disambiguierende Träger — LP3s Vermerk trägt die
Risiko-Bedingung. Plan-Text ist noch Platzhalter.

**Einordnung:** Alle drei offenen Platzhalter sind konsistent mit der
unchecked DoD-Checkbox „Jedes Risiko aus §6 trägt einen Ausgang" — kein
Widerspruch zwischen Plan-Zustand und Checkbox-Zustand. Closure-Arbeit,
kein Verifikations-Defekt.

---

## D. Plan-vs-Code-Diff über den gesamten Slice-Lauf

- `.d-check.yml`-Diff (`f02a7b5`) beschränkt sich exakt auf den in §1
  angekündigten `trace:`-Block (13 Zeilen, ausschließlich Ergänzung) — kein
  Touch an `modules:`, `GATE_CHECKS`, `d-check.mk`s Digest-Pin
  (`sha256:18e9c…`, `v0.75.0` unverändert im Diff).
- `harness/README.md`/`harness/sensors/docs-check.md`-Diff (`b6daee6`)
  entspricht exakt dem in §3 des Plans vorgesehenen Umfang — je eine neue
  Zeile/ein neuer Punkt, kein Umschreiben bestehender Aussagen.
- Keine Änderung an `--require-complete`/`doc-complete`-Aktivierung,
  `trace.requirements.modality`, `trace.cross-consistency` — eigene
  Prüfung von `.d-check.yml` und `d-check.mk` bestätigt: `doc-complete`
  bleibt ein reines `d-check.mk`-Target ohne `GATE_CHECKS`-Bindung (§A).
- Keine Sanierung der sieben realen Waisen — kein Diff an
  `spec/lastenheft.md` oder `docs/user/e2e-abdeckung.md` im gesamten
  Slice-Bereich (`git diff 9562042..HEAD -- spec/lastenheft.md
  docs/user/e2e-abdeckung.md` liefert leer, eigene Prüfung).
- Kein Versions-Bump des `d-check`-Digests: `d-check.mk`s
  `DCHECK_DIGEST`/`DCHECK_TAG` unverändert im Diff (identisch mit dem im
  Plan zitierten `v0.75.0`/`sha256:18e9c…`, bestätigt durch den `docker
  run`-Aufruf-Log jedes eigenen `make doc-trace`/`make gates`-Laufs in
  §A).
- Review-Findings: **keine** HIGH/MEDIUM/LOW im Ausgangs-Review — kein
  Fixrunden-Zyklus nötig, keine Rückkante zu prüfen. Das einzige Finding
  (F-1, INFO) ist unten separat bewertet (§E).
- **Keine über LP1–LP3/Träger-Nachzug hinausgehende Abweichung** zwischen
  Plan und committetem Code gefunden.

---

## E. Bewertung F-1 (Exit-Code-Verwechslung `d-check` vs. `make`-Wrapper)

**Rein informativ, keine DoD-Relevanz.** Eigene Reproduktion (§A) bestätigt
exakt die Review-Feststellung: Der Docker-Aufruf selbst liefert bei
`--trace --require-complete` gegen den Bestand **Exit 1** — deckungsgleich
mit der Plan-Aussage in §1 („real gemessen: Exit 1 bei den verbleibenden 7
Waisen"), die sich explizit auf das zugrunde liegende Werkzeug bezieht,
nicht auf `make doc-complete` selbst. Der Make-Wrapper liefert separat
**Exit 2** (GNU-Make-Standardverhalten bei gescheitertem Rezept). Die
Plan-Aussage ist damit **korrekt**, nicht fehlerhaft — F-1 markiert eine
Verwechslungsgefahr für künftige Leser, keinen bestehenden Fehler. Kein
DoD-Kriterium hängt an diesem Exit-Code (`--require-complete` ist
ausdrücklich nicht Teil des Slice-Umfangs, §1 Abgrenzung); die
Advisory-Eigenschaft von `--trace` selbst (Exit 0) ist unabhängig davon
korrekt und eigenständig verifiziert (§A, `make doc-trace` Exit 0).

---

## F. Entscheidungs-Konformität

- Keine `Accepted`-ADR wird durch diesen Slice inhaltlich verändert —
  `git diff 9562042..HEAD -- docs/plan/adr/` ist leer (eigene Prüfung).
- `make doc-immutable RANGE=9562042..HEAD` → Exit 0 (§A): keine
  `MR-*`-Kerndatei verändert.
- `make doc-commits RANGE=9562042..HEAD` → Exit 0 (§A): alle Commits
  traceable.
- Kein Architect-Zug nötig oder erfolgt — der Slice trägt bewusst keine
  eigene ADR-Entscheidungsfrage (§`Bezug` im Plankopf), anders als der
  Schwester-Slice `slice-d-check-tracked-modul`.

---

## G. Offene Closure-Obliegenheiten (benannt, nicht ausgeführt)

1. **`make gates` grün — Checkbox formal nachziehen.** Faktenlage liegt vor
   (§A, §B): Exit 0, alle sechs Gate-Bestandteile grün.
2. **§6-Risiko-Ausgänge (alle drei)** — Faktenlage liegt vor (§C):
   Risiko 1 „weiter offen" (Kopplungs-Vermerk), Risiko 2 und Risiko 3
   „entfallen".
3. **Beobachtungs-Register-Entscheidung** — ob die 7 realen Waisen eine
   eigene, neue Beobachtung werden, oder keine — Planner-Wertung, hier
   nicht vorweggenommen.
4. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag — noch nicht
   geschrieben.
5. **`git mv` nach `done/`** (eigener Commit nach dem Inhalts-Commit,
   `AGENTS.md` §3.3) — noch nicht erfolgt (Datei liegt weiterhin unter
   `in-progress/`).
6. **Drei Paarungen** — von der Welle-Closure `welle-d-check` zu prüfen,
   zusammen mit dem Schwester-Slice `slice-d-check-tracked-modul`
   (bereits `done/`) — die welle-weite Bedingung aus dem Plankopf
   (kombinierter `make gates`-Lauf mit beiden Slices gleichzeitig aktiv in
   derselben `.d-check.yml`) ist durch diesen Verifikations-Lauf faktisch
   bereits erfüllt: `.d-check.yml` trägt zum Zeitpunkt dieses Berichts
   sowohl den `tracked`-Modul-Eintrag (`modules:`-Liste) als auch den
   `trace:`-Block gleichzeitig, und `make gates` lief genau darüber grün
   (§A) — formale Feststellung bleibt der Welle-Closure vorbehalten.

---

## H. Verdikt

**DoD-Konformität (Implementierungs-Teil, LP1–LP3 + Review + Doku-Update):
erfüllt.** Alle fünf `[x]`-Checkboxen sind am Artefakt belegt, nicht nur
behauptet: `.d-check.yml` trägt `trace.requirements.id-pattern` und
`trace.coverage` korrekt, beide Effekte sind durch eine eigene, isolierte
Doppelmessung reproduziert (76/9 ohne, 76/7 mit — exakte Differenzmenge
`LH-FA-CAP-002`/`LH-FA-CAP-003` bestätigt), die zwei Träger
(`harness/README.md`, `harness/sensors/docs-check.md`) sind vollständig
nachgezogen, und der Review schließt mit 0 HIGH/0 MEDIUM/0 LOW ohne
Fixrunde.

**`make gates`: real grün (Exit 0), Checkbox im Plan noch nicht formal
angehakt.** Kein DoD-Verstoß — die Commit-Message von `fce746c` weist dies
bewusst den nachgelagerten Rollen zu; die zugrunde liegende Tatsache ist
durch einen eigenen, ungepipten Lauf bestätigt.

**Commit-Traceability: erfüllt.** Alle neun Commits tragen `ADR-0045`,
keine `SPEC-*`/`ARC-*`-Kennung im Betreff; `doc-commits`/`doc-immutable`
über die volle Slice-Range beide Exit 0; beide `git mv`-Commits sind reine
Moves (`AGENTS.md` §3.3).

**Out-of-Scope-Disziplin: eingehalten.** Kein Touch an
`--require-complete`/`doc-complete`-Aktivierung, `trace.requirements.
modality`, `trace.cross-consistency`, keine Sanierung der 7 realen Waisen,
kein Digest-Bump — jeweils eigenständig am Diff und an `.d-check.yml`/
`d-check.mk` geprüft, nicht übernommen.

**F-1 (INFO): keine DoD-Relevanz.** Eigene Reproduktion bestätigt: die
Plan-Aussage ist korrekt (bezieht sich auf den Tool-Exit-Code 1, nicht den
Make-Wrapper-Exit-Code 2); reiner Lese-Hinweis ohne erwartete Aktion.

**Nicht closure-bereit — erwartungsgemäß.** Fünf DoD-Checkboxen bleiben
unchecked: `make gates` grün (formal), Closure-Notiz,
Beobachtungs-Register-Entscheidung, vollständige §6-Risiko-Ausgänge (3 von
3 fehlen als Plan-Text), `git mv` nach `done/`. Das ist **keine
DoD-Verletzung**, sondern der erwartete Zustand unmittelbar vor der
Closure — dieser Bericht ist der Verifikations-Schritt, der dem
Planner-Closure-Zug vorausgeht (Modul 8). Die Faktenlage für `make gates`
und alle drei Risiko-Ausgänge liegt bereits vollständig vor (§A, §C, §G).

**Übergabe an den Planner:** §G-Obliegenheiten, insbesondere die drei noch
nicht geschriebenen Risiko-Ausgänge (Fakten liegen vor), das formale
Nachziehen der `make gates`-Checkbox, die Closure-Notiz, sowie —
welle-weit — die Prüfung der drei Paarungen für `welle-d-check` zusammen
mit dem bereits geschlossenen Schwester-Slice
`slice-d-check-tracked-modul`.

---

Der Bericht ist ein **Lauf-Beleg** dieses Verifier-Laufs (dieser Diff,
diese Sensoren, dieses Modell) und ersetzt weder Review noch Closure.
