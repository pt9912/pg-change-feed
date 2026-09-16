# Verifikationsbericht: slice-090 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-090` §2, LP1–LP3), die im Slice referzierte Entscheidung
[`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md)
(**Festlegung 1** = der Baum wird nicht geschrieben · **Festlegung 2** = der
Befund nennt Datei **und** Zeile · §Was das Gate nicht fangen kann ·
§Re-Evaluierungs-Trigger) sowie die Hard Rules `AGENTS.md` §3.7, §3.9, §3.10,
§3.11, §3.12, §4. **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-090.md`](review-slice-090.md) abgeschlossen) und **nicht** gegen
realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan am Stand `HEAD` gelesen,
dazu [`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md),
`ADR-0060`, den Review-Report, das neue Sensor-Dokument und den Skriptkopf.
Review und Report wurden als **Kontext** gelesen, ihre Findings **nicht**
übernommen: jede Aussage unten stammt aus einem hier selbst gefahrenen Lauf
oder aus einem Zitat, das ich gegen seine Quelle gehalten habe. Exit-Codes je
ungepiped und in eigenem Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und
Folgehandlung getrennt beauftragt. Kein Gate-Exit-Code dieses Berichts ist über
eine Pipe oder einen Wrapper gelaufen; Logs liegen **außerhalb** des Baums.

**Gegenstand:** `slice-090`, Stand `HEAD` = `c2bc08f`, Zweig `main`, Baum
sauber (`git status --porcelain` leer — vor **und** nach jedem Lauf dieses
Berichts). Die drei Commits des Vorgangs: `9029c05` (Implementierung),
`0122eb1` (Review-Report), `c2bc08f` (Fixrunde). Der Slice liegt in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt. Der Range
`9029c05^..c2bc08f` umfasst **8 Dateien** (4 `A`, 4 `M`) — kein Go-Code, kein
`Makefile`, kein `Dockerfile`, kein `proto/`, keine `.pb.go`.

**Eine Zahl zur Schreibweise:** das Gate-Skript endet im Rot-Fall mit **Exit
1**, `make` kapselt den Rezept-Fehlschlag zu **Exit 2**. Alle Rot-Läufe dieses
Berichts sind über `make` gefahren und stehen deshalb als **2**; das ist die
Zahl, die ein Prüfer an der Shell sieht.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify v6.5.0 OK` (54 Dateien) · `coverage-gate: OK — Coverage 74.80% erfüllt Schwelle 70%` · `d-check: 737 Datei(en), 0 Befund(e)` · `commit-traceability: OK` (5 Commits) · `generated-sync: OK` · `a-check` (Digest-Image) · `gesamt: 0 Befund(e)`; danach `git status --porcelain` leer |
| 2 | `make -p` → `^GATE_CHECKS` | **0** | `baseline-verify coverage-gate docs-check commit-traceability generated-sync a-check` — **sechs** Einträge, das neue Ziel enthalten |
| 3 | `make -n record-gates` | **0** | der `generated-sync`-Aufruf steht als Prüfschritt **vor** `record-gates.sh` |
| 4 | zwei Läufe `make generated-sync` + Status je danach | **0 / 0** | beide `OK — das committete Erzeugnis ist byte-gleich …`; `git status --porcelain` danach je **0 Bytes** → LP2 |
| 5 | Dauer der beiden Läufe (`date +%s%N`) | — | **968 ms** / **999 ms** (warmer Layer-Cache) |
| 6 | Mutation Zeile 226 (Ein-Zeilen-Änderung), Gate | **2** | Befund nennt **Zeile 226**; Kontext-Diff-Kopf im selben Log `@@ -223,7 +223,7 @@` → die **alte** Herleitung hätte 223 genannt |
| 7 | Mutation Zeile 102 (`GetSequence`-Zeile — der F-2-Fall), Gate | **2** | Befund nennt **Zeile 102**; Kontext-Diff-Kopf `@@ -99,7 +99,7 @@` → alte Herleitung **99**; eigenes `diff -U0` gegen den mutierten Stand: `@@ -102 +102 @@` |
| 8 | Blocklöschung (Original-Zeilen 231–233), Gate | **2** | Befund nennt **231**; eigenes `diff -U0` → `@@ -230,0 +231,3 @@` — genau die im Skriptkopf/Sensor-Dokument dokumentierte `N+1`-Regel |
| 9 | `proto/` real gedriftet (Feld umbenannt), Gate ohne Override | **2** | `FAIL — … changestream.pb.go weicht … ab (erste abweichende Stelle …: Zeile 37)` + Diff |
| 10 | derselbe Drift-Stand, `GENERATED_SYNC_SOURCE_DIR=<Kopie im Baum>` | **0** | `OK — …`, `Quelle: …/proto-orig/cdc/stream/v1/changestream.proto` — der Grün-Pfad des Overrides ist **real** |
| 11 | derselbe Override auf ein Verzeichnis **außerhalb** des Baums | **2** | `Could not make proto path relative: …: No such file or directory` → laut rot, kein stiller Grün-Pfad |
| 12 | `TMPDIR=<Verzeichnis im Baum>` | **2** | zwei Rot-Zeilen über die **eigenen Temp-Dateien** (`… ist als Erzeugnis gekennzeichnet (protoc-gen-go), der Lauf erzeugt sie aber nicht`); der `trap` räumt, Status danach leer → Sensor-Doku §Grenze 3 |
| 13 | `make coverage-gate` zweimal (eigene Aufrufe) | **0 / 0** | **74.70 %** / **74.70 %** — gegen **74.80 %** aus dem Aggregat-Lauf (#1) |
| 14 | `internal/bootstrap` 8× mit `-coverpkg`, in einem Container | **0** (8×) | `runWALRetentionCheck 87.5 %` in **7** Läufen, `100.0 %` in **1** Lauf — die Ursache des Flaps aus #13 |
| 15 | Zell-Längen der Vertragsspalte in `harness/README.md` §Sensors | **0** | `generated-sync` **227** · `commit-traceability` 295 · `coverage-gate` 286 · `gates` 135 · `a-check` 123 · `docs-check` 101 · `baseline-verify` 91 |
| 16 | `.github/workflows/ci.yml` Z. 10–11 und Z. 66 gegen `GATE_CHECKS` | **0** | beide nennen **sechs** Namen (baseline-verify, docs-check, a-check, commit-traceability, coverage-gate, generated-sync) — deckungsgleich |
| 17 | `docker build --no-cache-filter proto --target proto` (kalte Stufe, Probe-Tag danach entfernt) | **0** | `#9 DONE 8.5s` — `apk add` 117 Pakete + zwei `go install` **mit Downloads**; die Zahlen der Sensor-Doku §Grenze 4 reproduziert |
| 18 | `grep -rn '^generated-sync:\|GENERATED_SYNC' --include='*.mk' --include=Makefile` | **0** | **ein** Ort: `harness/mk/generated-sync.mk` — keine zweite Deklaration |
| 19 | `git diff --name-only 9029c05^..c2bc08f -- '*.go'` | **0** | **0** Dateien — kein Go-Code im Range |
| 20 | `ls docs/reviews/ \| grep 090` | **0** | nur `review-slice-090.md` — **kein** Delta-Report |
| 21 | DoD-Checkboxen des Slice-Plans (`grep '^ *- \['`) | **0** | **12× `[ ]`**, keine einzige gesetzt |
| 22 | Dockerfile-Stufe `proto` gelesen (`:34–36`) | **0** | `protobuf-dev=31.1-r1`, `protoc-gen-go@v1.36.12`, `protoc-gen-go-grpc@v1.6.2` — genau die Pins, die Skriptkopf und Sensor-Doku zitieren |
| 23 | `ls docs/plan/planning/observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/evidence/` | **0** | 4 Dateien (`slice-069`, `-074`, `-082`, `-086`) — der Stand, den §8 des Plans nennt; im Range **keine** Register-Datei geändert |
| 24 | `ls docs/plan/planning/reconciliation.md` | **1** | Datei existiert nicht (Greenfield) → das §2-Item „entfällt" ist nachweislich richtig |

---

## 2. DoD-Konformität, Kriterium für Kriterium

**Liefer-Punkt 1 — das Gate existiert und hängt am Aggregat.**

| Kriterium (§2) | Befund |
|---|---|
| ein `make`-Ziel läuft den gepinnten Generator gegen `proto/cdc/stream/v1/changestream.proto` | **erfüllt** — `harness/mk/generated-sync.mk` + `tools/harness/generated-sync.sh`; Skriptkopf, Dockerfile (`:34–36`) und Sensor-Doku nennen **dieselben** Pins (#22) |
| vergleicht mit dem committeten Erzeugnis | **erfüllt** — der Erfolgstext nennt beide geprüften Dateien; der Rot-Fall ist real gesehen (#6–#9) |
| über `GATE_CHECKS` in `make gates` eingebunden | **erfüllt** — `make -p` löst sechs Einträge auf, `generated-sync` darunter (#2); `record-gates` fährt es vor dem Stempel (#3); **eine** Deklaration (#18) |
| `make gates` ist grün | **erfüllt** — Exit **0** aus separater Datei gelesen (#1) |

**Liefer-Punkt 2 — das Gate schreibt den Arbeitsbaum nicht.**

| Kriterium (§2) | Befund |
|---|---|
| erzeugt in ein Temp-Verzeichnis | **erfüllt** — `mktemp -d` + `trap … EXIT` vor jedem frühen Ausgang; Baum als `:ro`-Bind-Mount (#4, #6–#9: jeder Rot-Lauf ließ den Baum sauber) |
| `git status --porcelain` nach dem Lauf leer — **auch beim zweiten Lauf** | **erfüllt** — zwei Läufe hintereinander, Status danach je **0 Bytes** (#4). Das ist der Nachweis, den die DoD-Zeile verlangt: das Gate schreibt nicht „erst und liest dann" |

**Liefer-Punkt 3 — ein roter Befund nennt die Abweichung.**

| Kriterium (§2) | Befund |
|---|---|
| der Rot-Fall ist **real gesehen**, nicht behauptet | **erfüllt** — vier eigene Rot-Läufe (#6, #7, #8, #9) |
| die Ausgabe nennt **Datei und Zeile** | **erfüllt** — `generated-sync: FAIL — <Datei> weicht von der Generatorausgabe ab (erste abweichende Stelle im committeten Erzeugnis: Zeile N)`, darunter der Unified-Diff |
| **und die Zeile ist die abweichende** (`ADR-0084` Festlegung 2, Review-F-2) | **erfüllt** — die mutierte Zeile **102** wird als 102 genannt (#7); die alte Herleitung hätte 99 genannt, gemessen am Kontext-Diff-Kopf desselben Logs. Zweite, unabhängige Probe: Mutation auf Zeile 226 → 226 (#6). Ein Befund „nicht synchron" ohne Datei existiert nicht |
| danach wird die Änderung zurückgenommen | **erfüllt** — nach jeder Mutation `git checkout --`, Status danach leer; im Range keine `proto/`- oder `.pb.go`-Änderung (Umfeld von #19) |

**Die Closure-Pflichten aus §2 (die neun unteren Zeilen).**

| Kriterium (§2) | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#1) |
| Review durchgeführt, Report unter `docs/reviews/`, kein Self-Review | **formell erfüllt, Substanz offen** — `review-slice-090.md` liegt vor und hat den Stand `9029c05` geprüft; die **Fixrunde `c2bc08f` hat keine Review** (kein Delta-Report, #20) → **V-2** |
| Doku-Update: Ziel in `harness/README.md` §Sensors **und** `AGENTS.md` §4; `proto-generate` bleibt **Werkzeug** | **erfüllt** — `harness/README.md:118` (Vertragszelle **227** Zeichen, #15, Target-Zelle ist der Link auf `harness/sensors/generated-sync.md`), `:119` sechs Namen; `AGENTS.md:439` eine Zeile; `harness/README.md:128` und `AGENTS.md:443` führen `make proto-generate` unverändert als Werkzeug |
| Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-090.md` | **erfüllt** mit diesem Bericht; das Häkchen ist offen (#21) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt** — §7 (`:231–242`) trägt in **9 von 10** Inhaltszeilen Platzhalter |
| Reconciliation-Register fortgeschrieben *(entfällt …)* | **entfällt nachweislich** (#24) |
| Beobachtungs-Register fortgeschrieben | **nicht erfüllt** — im Range **keine** Register-Datei (#23); die Zeile selbst bleibt offen. Der Träger `BEO-PGC/generierte-artefakte-ohne-sync-sensor/` steht mit **4** Belegen, `state.md` ist unberührt (wie der Plan es verlangt) |
| Jedes Risiko aus §6 trägt einen Ausgang | **nicht erfüllt** — alle drei stehen auf `— **Ausgang:** <…>` (`:205`, `:208`, `:212`) |
| Die drei Paarungen sind getragen | **nicht erfüllt** — fällt laut Zeile der `welle-20`-Closure zu |

**Ergebnis:** Die **drei Liefer-Punkte**, die Gate-Zeile und die Doku-Zeile
**tragen** — gemessen, nicht gelesen. Offen sind die fünf regulären
Closure-Pflichten (Review-/Verifikations-Häkchen, §7, Register-Zeile, §6-Ausgänge,
Paarungen) sowie die zwei Punkte, die nur diese Rolle sieht: **V-1** und **V-2**.

---

## 3. `ADR-0084` Festlegung 1 und 2 — eigene Messung

**Festlegung 1 (der Baum wird nicht geschrieben) — trägt.** Zwei Läufe
hintereinander in **diesem** Baum (#4): je Exit 0, `git status --porcelain`
danach je 0 Bytes. Der Generator läuft mit `--network none`, der Baum ist
`:ro`-Bind-Mount, geschrieben wird nur nach `/out` (Temp-Verzeichnis). Die
Doppelung ist der eigentliche Beleg: hätte das Ziel den Baum geschrieben, wäre
der zweite Lauf grün geworden — er ist es aus dem richtigen Grund (kein
Schmutz, `cmp` byte-gleich).

**Die Grenze dazu, selbst gemessen (#12):** mit `TMPDIR` auf ein Verzeichnis
**innerhalb** des Baums wird der Lauf **falsch-rot** — die zweite
Vergleichsrichtung sucht die gekennzeichneten Erzeugnisse ab dem Baum-Wurzel
und findet dann die eigenen Temp-Dateien. Der `trap` räumt sie weg, danach ist
`git status --porcelain` leer: die Zusage aus Festlegung 1 hält auch in diesem
Fall, und die Fehlerrichtung ist die sichere. Das Sensor-Dokument §Grenze 3
nennt genau das — die Aussage ist gedeckt.

**Festlegung 2 (der Befund nennt Datei und Zeile) — trägt, nachgemessen am
F-2-Fall selbst.** Der Review hatte gefunden, dass die Zeile aus dem
Hunk-**Anfang** des Kontext-Diffs kam (99 statt 102). Reproduktion mit eigener
Mutation an genau derselben Zeile (#7): Mutation auf `changestream.pb.go:102`,
Lauf, Befund nennt **102**; der Kontext-Diff-Kopf im selben Log nennt
`@@ -99,7 +99,7 @@` — die **alte** Herleitung wäre 99 gewesen. Die Abweichung ist
also real behoben, nicht behauptet. Zwei weitere Stellungen sichern die Regel
gegen Zufall ab:

- Mutation auf Zeile **226** → Befund nennt **226** (#6);
- **Blocklöschung** (Original-Zeilen 231–233) → Kopf des kontextfreien Diffs
  `@@ -230,0 +231,3 @@`, Befund nennt **231** — genau die im Skriptkopf
  (`:24–26`) und Sensor-Dokument §Ausgabe **dokumentierte** `N+1`-Regel für
  den Einfüge-Fall.

Damit ist die Zeile in **beiden** Zweigen der Herleitung die abweichende bzw.
die dokumentiert nächste — Festlegung 2 hält.

**Festlegung 2, zweite Hälfte („… und nur das") — keine Über-Erklärung
gefunden.** Das Grün behauptet nichts über Kompilierbarkeit oder
Draht-Vertrag (Skriptkopf `:6–9`, Sensor-Dokument §Nicht Gegenstand) — der
Range enthält keinen Go-Code (#19); die Aussage, dass `make test`/`make image`
dafür zuständig sind, ist unberührt. Die Erfolgszeile nennt Quelle und
geprüfte Dateien — ein Grün ist damit lesbar lokalisiert.

---

## 4. Plan-vs-Code-Diff

**Verglichen wurde gegen den Plan-Stand *vor* der Arbeit (`b3542c1`) — nicht
gegen die Fassung am `HEAD`.** Der Grund ist eine Eigenschaft dieses Runs: der
Plan wurde **von den Implementer-Commits selbst** nachgezogen (`9029c05`
+2/−2, `c2bc08f` +4/−2 in §3). Ein Vergleich gegen `HEAD` wäre zirkulär, weil
die Liste dann die gelieferten Dateien bereits enthielte.

| §3-Zeile (Stand `b3542c1`) | Geliefert | Urteil |
|---|---|---|
| `harness/mk/generated-sync.mk` — **neu** („Name des Ziels führt der Implementer") | ja, Ziel heißt `generated-sync` | **Plan eingehalten**, Name wie in der ADR vorgeschlagen |
| `tools/harness/generated-sync.sh` — **neu** | ja | **Plan eingehalten** |
| `Makefile` — **update** („Das neue Ziel deklarieren") | **nicht angefasst**; die Zeile wurde in `9029c05` auf **nicht** korrigiert | **Abweichung — begründet.** `include harness/mk/*.mk` zieht das Fragment ein; Ziel **und** `GATE_CHECKS`-Anhang stehen in **einer** Datei (Haus-Muster `coverage.mk`). Faktisch nachgeprüft: **eine** Deklaration (#18), `Makefile` nicht im Range. Die ursprüngliche Zeile hätte eine zweite Deklarationsstelle erzeugt |
| `harness/README.md` §Sensors — **update** | ja, Zelle + `make gates`-Zeile | **Plan eingehalten** |
| `AGENTS.md` §4 — **update** | ja, eine Zeile | **Plan eingehalten** |
| `…/generierte-artefakte-ohne-sync-sensor/state.md` — **nicht** | nicht angefasst | **Plan eingehalten** (#23) |
| `docs/user/e2e-abdeckung.md` — **nicht** | nicht angefasst | **Plan eingehalten** |
| — (im Plan **nicht** als Zeile geführt) | `harness/sensors/generated-sync.md` **neu**, `.github/workflows/ci.yml` **update** | **Nachtrag, keine Abweichung im Ergebnis**: beide Zeilen wurden in der Fixrunde ergänzt (`c2bc08f`), zusammen mit der Arbeit, die sie beschreiben. Der Plan-Nachzug inkl. **Nicht-Realisierungen** im selben Lauf ist die verkörperte Regel (`.claude/commands/implement-slice.md` Schritt 14, `BEO-PGC/plan-nachzug`, seit slice-009) — der Nachzug ist damit die vorgeschriebene Form, nicht das still gestrichene Weglassen |

**Was die Nachträge ausgelöst hat:** `harness/sensors/generated-sync.md` ist die
Antwort auf Review-F-1 (die Sektion §Sensors schreibt für einen Überhang die
Sensor-Datei mit der Zelle als Link vor), `ci.yml` die Antwort auf F-4. Beiden
ist eigen, dass der **Reviewer sie nicht im Zug erwartete** — F-4 hatte der
Report ausdrücklich der Planner-Seite zugewiesen („gehört als Fundstelle
benannt, nicht in diesen Zug hinein"). Der Zug hat sie trotzdem mitgenommen und
im Plan sichtbar gemacht; das ist eine **Umfangs-Erweiterung gegenüber dem
Review-Verdikt**, nicht gegenüber dem Plan (weil der Plan sie mitführt).

**Nicht im Plan, nicht im Code:** die beiden von der Fixrunde geänderten
Träger haben **keine** eigene §1-Out-of-Scope-Zeile und keinen Folge-Slice
gebraucht — die Änderung ist eine Textkorrektur an einer bereits bestehenden
zweiten Fassung, keine Mechanik. Kein Befund; festgehalten wird nur, dass die
Zell-Länge und die Formfrage der Sensors-Zelle damit **durch den Zug selbst**
entschieden worden sind (F-1 war als „nein, unverifizierbar" geführt), und der
Träger jetzt dem Haus-Muster folgt.

---

## 5. Die Zahlen, die in §7 landen werden — nachgemessen (`AGENTS.md` §3.12 Instanz A)

| Zahl | Wo sie heute steht | Mein Lauf | Befund |
|---|---|---|---|
| Vertragszelle `generated-sync` = **227** Zeichen | Commit `c2bc08f` (Message) | 227 (#15) | **hält** — und liegt wieder in der Größenordnung der Nachbarzeilen (91–295) |
| `make gates`-Enumeration = **sechs** Namen | `harness/README.md:119`, `ci.yml:10–11`, `:66` | sechs, deckungsgleich mit `GATE_CHECKS` (#2, #16) | **hält** — die stehen gebliebene Zweitfassung (F-4) ist nachgezogen |
| **Coverage „74.80 %"** | Review-Report nennt 74.70 %, `verify-slice-089` 74.70 %, die Aufgabe nennt 74.80 % | Aggregat-Lauf gibt **74.80 %** (#1), `make coverage-gate` allein gibt **74.70 %** (#13, zweimal) | **hält nicht reproduzierbar** → **V-1** |
| Gate-Dauer **„0,9 s"** warm | `ADR-0084` §Konsequenzen („billig"), Review-Messung 0,900 s | **968 ms / 999 ms** (#5); Sensor-Doku nennt 941/928 ms | **hält** (Größenordnung; eine Zahl für §7 trägt ihren Lauf) |
| `proto`-Stufe kalt **8,5 s**, braucht Netz | Sensor-Dokument §Grenze 4 | **8.5 s** reproduziert, mit sichtbaren Downloads (#17) | **hält** — und belegt zugleich §6-Risiko 2 (der **Build** braucht Netz, der **Lauf** ist `--network none`) |
| Register-Stand **4×** / **2×** / **1×** | §8 des Plans | 4 Belege im Eintrag (#23) | **hält** |

---

## 6. Findings

### V-1 — Die Coverage-Zahl des Standes ist **nicht reproduzierbar** (Flap 74,70 ↔ 74,80 %), und der Träger nennt sie als Messung

- `kategorie`: MEDIUM
- `pfad`: `harness/mk/coverage.mk:12` (`coverage-gate` gegen `THRESHOLD ?= 70`)
  gegen `internal/bootstrap/walretention_internal_test.go:175`, `:219`;
  betroffene Messstelle `internal/bootstrap/wiring.go:981`
- `befund`: Derselbe Befehl am **selben** Stand liefert verschiedene Zahlen —
  `make gates` (Aggregat) **74.80 %** (#1), `make coverage-gate` zweimal
  **74.70 %** (#13). Der Grund ist lokalisiert: `go tool cover -func` gibt für
  `runWALRetentionCheck` (`wiring.go:981`) in **7 von 8** Läufen **87.5 %**, in
  **einem** Lauf **100.0 %** (#14) — eine nicht abgedeckte Anweisung im
  `select`/Ticker-Pfad, deren Erreichung von der Laufzeit abhängt. Der
  Test selbst ist in **allen** Läufen grün (`rc=0`), die Unit-Tests melden also
  nichts: **diese Klasse ist für Tests unsichtbar und für ein diff-skopiertes
  Review nur über eine Wiederholungsmessung sichtbar.** Folge für §3.12: die
  Zahl trägt mit der Nennung ihres Laufs einen Ursprung — sie ist aber als
  **Ergebnis einer Wiederholung** nicht stabil, und ein späterer Leser kann sie
  nicht nachfahren. Solange `THRESHOLD = 70 %` und die gemessene Zahl ~4,8 pp
  darüber liegt, ist keine Gate-Entscheidung gefährdet; bei der Rampe nach
  `ADR-0077` (Endstufe 80 %) rückt die gemessene Zahl an die Schwelle, und ein
  Flap von 0,1 pp wird dann zur **flappenden Gate-Entscheidung**.
- `verifizierbar`: ja — `make coverage-gate` zweimal; `docker run --rm
  --network none pg-change-feed:coverage bash -c 'for i in 1 ··· 8; do go test
  -count=1 -coverpkg=./internal/bootstrap/ -coverprofile=c.out
  -covermode=atomic ./internal/bootstrap/; go tool cover -func=c.out |
  grep runWALRetentionCheck; done'` (das Profil ist der Container-interne
  Ablagepfad)
- `urteil`: **vor der Closure zu adressieren, nicht zu beheben.** Die Antwort
  ist keine Änderung an diesem Slice (das Gate ist nicht sein Gegenstand): die
  §7-Zeile nennt ihre Zahl **mit ihrem Lauf** und dem Hinweis auf die
  Flap-Quelle, und die Beobachtung gehört **benannt** in das
  Beobachtungs-Register. Zwei Adressen existieren bereits und sind zu prüfen:
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**5×**, verkörpert in
  `AGENTS.md` §3.12 — hier die neue Nuance „gemessen, aber nicht
  reproduzierbar") und `BEO-PGC/test-integration-retention-timing-flake`
  (1×, offen — derselbe **Test-Gegenstand** `WAL-Retention`, andere Träger:
  dort der Integrations-Exit, hier die Coverage-Zahl). Ob die Zuordnung eine
  **zweite Beobachtung** oder ein **Auftreten** des bestehenden Eintrags ist,
  entscheidet das Register-Urteil (Modul 6: „Mensch urteilt, Maschine prüft
  Deckung"); ein Beleg-Vorgang liegt mit diesem Slice erst nach seiner Closure
  vor.

### V-2 — Die Fixrunde `c2bc08f` ist von **keiner** Review gedeckt

- `kategorie`: MEDIUM
- `pfad`: `docs/reviews/` (nur `review-slice-090.md`) gegen
  `docs/reviews/review-slice-090.md` §Verdikt und
  `tools/harness/generated-sync.sh` (Stand der Fixrunde)
- `befund`: Der Review-Report beurteilt ausdrücklich den Stand `9029c05` und
  schließt mit „**Rückgabe an den Implementer** … F-1 und F-2 … der Slice
  braucht eine Fixrunde". Die Fixrunde hat danach den **Kern des Lieferwerts**
  geändert (die Herleitung der Befundzeile), dazu die Vertragszelle in
  `harness/README.md`, ein **neues** Sensor-Dokument und
  `.github/workflows/ci.yml`. Für diesen Stand existiert **kein**
  Review-Artefakt (#20) — gemessen: `ls docs/reviews/ | grep 090` → nur der
  Report zum Vorstand. Der Haus-Präzedenzfall `slice-089` hatte für seine
  Fixrunde einen eigenen Delta-Report (`review-slice-089-delta.md`), der die
  Korrekturen gegen die Findings hielt. Die **Substanz** von F-2 habe ich in
  diesem Bericht unabhängig bestätigt (§3: 102 statt 99), ebenso die
  Formhälfte von F-1 (Zelle 227, Dokument vorhanden) — meine Rolle ist aber
  **Verifikation gegen DoD/ADR**, nicht Review gegen Plan/ADR/Hard Rules
  (Modul 8: „wer geschrieben hat, reviewt nicht"); §3.7-Kommentarklassen, die
  Haus-Form des neuen Sensor-Dokuments und die `ci.yml`-Änderung hat kein
  Reviewer gesehen.
- `verifizierbar`: ja — `ls docs/reviews/ | grep 090`;
  `git diff --name-status c2bc08f^..c2bc08f`
- `urteil`: **eine Entscheidung ist zu treffen, bevor der Slice schließt.**
  Entweder ein Delta-Review auf `c2bc08f` (Präzedenz `slice-089`), oder eine
  im §7 **ausdrücklich festgehaltene** Entscheidung, dass die Fixrunde ohne
  Delta schließt und warum. Beides ist Planner-/Reviewer-Sache; still bleiben
  darf es nicht, weil der Report selbst die Rückgabe als **nötig** bezeichnet
  hat.

### V-3 — Die DoD-Häkchen aus Schritt 18/21 sind nicht gesetzt, obwohl die Punkte materiell erfüllt sind

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-090-sync-gate-protobuf.md:96`,
  `:100`, `:105`, `:111`, `:112` gegen `.claude/commands/implement-slice.md:167–173`
  (Schritt 18) und `:246–249` (Schritt 21)
- `befund`: **Zwölf** DoD-Zeilen stehen auf `[ ]`, keine einzige ist gesetzt
  (#21). Materiell erfüllt sind davon LP1, LP2, LP3 und `make gates` (alle
  hier gemessen, §2), und nach der Fixrunde auch die Review-Zeile — Schritt 21
  sagt dafür wörtlich: *„Löst eine Fixrunde nach Reviewer-Findings einen bislang
  offenen DoD-Punkt auf (typischerweise ‚Review durchgeführt …'), wird die
  zugehörige Checkbox **im Fixrunden-Commit** mitgesetzt — nicht erst bei der
  Planner-Closure nachgetragen."* Die Klasse ist **verkörpert**
  (`BEO-PGC/dod-checkbox-nachzug`, seit welle-5) — der Fixrunden-Commit
  `c2bc08f` hat sie nicht ausgeführt, und auch `9029c05` hat die drei
  Liefer-Punkte nicht gesetzt, obwohl sie zu diesem Zeitpunkt belegt waren.
  Damit ist der Zustand nicht mehr von „noch nicht gearbeitet" unterscheidbar,
  und die nächste Rolle muss die Belege erneut zusammentragen.
- `verifizierbar`: ja — `grep -n '^ *- \[' docs/plan/planning/in-progress/slice-090-sync-gate-protobuf.md`
- `urteil`: **vor der Closure nachzuziehen** (LP1–LP3, `make gates`,
  Review-Zeile, Verifikations-Zeile). Kein Liefer-Defekt: die Häkchen tragen
  hier die Belege dieses Berichts. Die übrigen Zeilen (Closure-Notiz,
  Register, §6-Ausgänge, Paarungen) bleiben zu Recht offen — sie sind
  Planner-Arbeit.

### V-4 — Der Quell-Override greift **nur** für Pfade im Baum; außerhalb ist der Lauf laut rot — die Mount-Grenze steht nicht im Sensor-Dokument

- `kategorie`: INFO
- `pfad`: `harness/sensors/generated-sync.md:67` (Override-Tabelle) und
  `:78–79` (Grenze 2) gegen `tools/harness/generated-sync.sh:67–68`
  (`-v "$repo_root":/src:ro`)
- `befund`: Gemessen in beiden Stellungen: in-Tree-Kopie als Quelle →
  **`EC=0`** und `Quelle: <Pfad>` im Erfolgstext (#10, der Grün-Pfad ist real,
  die Grenze stimmt); Verzeichnis **außerhalb** des Baums → **`EC=2`** mit
  `Could not make proto path relative: …: No such file or directory` (#11).
  Der Grund ist strukturell: nur der Baum wird als `/src` in den Container
  gemountet. Die Fehlerrichtung ist die **sichere**, und der Erfolgstext macht
  den Override sichtbar — die Doku-Aussage („der Ausgang ist damit nicht still")
  hält. Unbenannt bleibt, dass die verglichene Quelle **im Baum** liegen muss.
- `verifizierbar`: ja — `make generated-sync GENERATED_SYNC_SOURCE_DIR=<Pfad
  außerhalb> · <Pfad im Baum>`
- `urteil`: kein Blocker. Eine Zeile in §Overrides oder §Grenze („die Quelle
  muss im Baum liegen, sonst bricht der Lauf") würde die Doku deckungsgleich
  mit dem Verhalten machen — Formfrage der Fixrunden-Nacharbeit, kein
  DoD-Bezug.

### V-5 — `AGENTS.md` §3.10 ist dem Buchstaben nach **nicht** ausgelöst — die CI-Wirksamkeit des neuen Gates bleibt trotzdem unbelegt

- `kategorie`: INFO
- `pfad`: `.github/workflows/ci.yml:10–11`, `:66`, `:50` gegen die neue
  `GATE_CHECKS`-Zeile
- `befund`: Die Fixrunde hat `ci.yml` angefasst (Kommentarzeile und
  Schrittname) — **keine** Stufe, keine Matrix, keine Abhängigkeit, kein neuer
  Schritt; §3.10 greift damit nicht, und der Review-Satz „dieser Diff ändert
  den Workflow nicht" ist für den Vorstand richtig, für `c2bc08f` überholt.
  Die **Wirkung** bleibt aber lokal ungeprüft: der Schritt `Gates …` fährt
  `make gates` und damit einen neuen Build-Schritt, der auf kaltem Cache
  **Netz** braucht und **8,5 s** kostet (#17). Der Job trägt
  `timeout-minutes: 30` (`ci.yml:50`) — gegen gemessene +8,5 s plus Bau-Zeit
  ist das kein Timeout-Risiko; ob der Lauf auf dem gehosteten Runner grün ist,
  kann nur ein realer Post-Push-Lauf zeigen.
- `verifizierbar`: ja — `grep -n 'timeout-minutes\|uses:\|name:' .github/workflows/ci.yml`;
  `git diff --name-status c2bc08f^..c2bc08f`
- `urteil`: **Adresse statt Neuanlage.** Genau dieser Gegenstand ist §6-Risiko 2
  („Der Generator-Lauf braucht Netz") — mein Lauf (#17) liefert dafür den Beleg
  und die Präzisierung: der **Build** braucht Netz, der **Lauf** ist
  `--network none`. Der Risiko-Ausgang (Planner) entscheidet, ob der reale
  Post-Push-Lauf vor der Closure abgewartet wird.

---

## 7. Negativbefunde

- **geprüft, ohne Befund: die „Baum nicht schreiben"-Zusage über den ganzen
  Rot-Pfad hinweg.** Jeder der vier Rot-Läufe (#6–#9) und der TMPDIR-Fall (#12)
  hinterließ `git status --porcelain` **leer**; das Temp-Verzeichnis wird vom
  `trap` entfernt, auch wenn der Vergleich rot endet. Das Gate kann den Baum
  nicht schmutzig zurücklassen — auch nicht im Fehlerfall.
- **geprüft, ohne Befund: die Pinnung.** Skriptkopf (`:12`), Dockerfile
  (`:34–36`) und der kalte Build (#17) nennen dieselben drei Pins
  (`protobuf-dev=31.1-r1`, `protoc-gen-go@v1.36.12`,
  `protoc-gen-go-grpc@v1.6.2`). Das Grün vergleicht gegen einen **gepinnten**
  Generator — die Voraussetzung, auf der Festlegung 2 überhaupt ruht.
- **geprüft, ohne Befund: die Bindungs-Deklaration.** Genau **eine** Stelle
  deklariert Ziel und `GATE_CHECKS`-Anhang (#18); `Makefile` und `Dockerfile`
  sind unberührt; `make -p` und `make -n record-gates` lösen beides auf (#2,
  #3). Kein Target in `AGENTS.md` §4 oder `harness/README.md` §Sensors ist
  halluziniert: alle sechs genannten Namen existieren als Targets (`make -p`).
- **geprüft, ohne Befund: die beiden Doku-Träger.** `harness/README.md:118`
  trägt **227** Zeichen, und die Target-Zelle ist der Link auf das vorhandene
  Sensor-Dokument; `harness/sensors/generated-sync.md` führt Deckungsgrenze,
  beide Vergleichsrichtungen, Exit-Codes, Overrides, **fünf** benannte Grenzen
  und die Bindung. §3.11-Probe: `make docs-check` über 737 Dateien **0
  Befunde** (#1) — kein `hostpath-forbidden`; die hinzugefügten Zeilen des
  Ranges tragen keine host-lokale absolute Pfadform.
- **geprüft, ohne Befund: §3.7-Klassen in den neuen Trägern.** Skriptkopf,
  Fragmentkopf, Sensor-Dokument: Zusage · Kopplung (Bind-Mount, Aufrufer-uid,
  Modul-Layout-Relation) · Abgrenzung („nicht Gegenstand: kompiliert/verhält
  sich richtig") · Rang-Zeiger (`ADR-0084` Festlegung 1, `ADR-0060`) — Indikativ
  über den Ist-Zustand, kein Vorher/Nachher, keine Slice-Nummer als Begründung;
  als Herkunft genau die Hausform (`· seit slice-090`).
- **geprüft, ohne Befund: kein Go-Code im Range** (#19). Die Zusage „das Gate
  prüft nicht, dass der Code kompiliert" ist damit auch nicht im Range
  gebrochen, und die Coverage-Zahl bezieht sich auf einen unveränderten
  Messgegenstand — was **V-1** umso schärfer macht.

---

## 8. Was ich nicht prüfen konnte

- **Der reale Post-Push-Lauf des Workflows.** §3.10 greift dem Buchstaben nach
  nicht (V-5); ob `make gates` mit dem neuen Build-Schritt auf dem gehosteten
  Runner grün ist, bleibt bis zum ersten realen Lauf unbelegt. Der Beleg ist
  über `gh run view`/`gh run list` nachziehbar — außerhalb dieses Laufs.
- **Der kalte Build im netzlosen Umfeld.** #17 zeigt den Build **mit** Netz
  (Downloads sichtbar); dass er **ohne** Netz scheitert, ist die Aussage des
  Sensor-Dokuments §Grenze 4 und von §6-Risiko 2 — ich habe sie nicht gefahren
  (ein netzloses `docker build` gegen einen leeren Cache ist nicht ohne
  Verwerfen des Layer-Caches herstellbar, und ein Cache-Eingriff wäre kein
  Verifikations-, sondern ein Umgebungs-Eingriff).
- **Die Coverage-Zahl als stabile Größe.** V-1 zeigt das Gegenteil; einen Beleg
  für **eine** verbindliche Zahl kann es am heutigen Stand nicht geben.
- **Die Fortgeltung in die Zukunft.** Ob das Gate bei einer **bewussten**
  Hand-Nachbearbeitung des Erzeugnisses rot wird, ist der gewünschte Befund
  (`ADR-0084` §Was das Gate nicht fangen kann); ob eine **falsche** `.proto`
  grün bleibt, ist ebenfalls dort benannt (Unit-/Integration-Tier). Beides ist
  Zusage, nicht Beleg dieses Laufs.

---

## 9. Verdikt

**Die drei Liefer-Punkte aus §2 tragen — gemessen.** LP1 (Ziel existiert,
hängt an `GATE_CHECKS`, `make gates` **Exit 0**), LP2 (zwei Läufe, danach
zweimal leerer Status), LP3 (vier Rot-Läufe, Datei **und** die **abweichende**
Zeile, danach zurückgenommen). `ADR-0084` **Festlegung 1** hält real,
einschließlich ihrer selbst gemessenen `TMPDIR`-Grenze (falsch-rot, sichere
Richtung), und **Festlegung 2** ist mit eigener Mutation am F-2-Fall selbst
bestätigt: der Befund nennt **102** — die Zeile, die die alte Herleitung als
99 verfehlte. Der Plan-vs-Code-Vergleich ergibt **keine unbegründete
Abweichung**: die zwei Nachträge und die `Makefile`-Nicht-Realisierung sind die
verkörperte Plan-Nachzug-Form, und der Verzicht auf die `Makefile`-Zeile ist
sachlich der bessere Weg (eine Deklaration statt zwei).

**`done/`-fähig ist der Slice damit noch nicht — fünf Punkte fehlen:**

1. **V-2** — die Fixrunde `c2bc08f` hat **keine** Review; der Report hatte die
   Rückgabe selbst als nötig bezeichnet, ein Delta-Review (Präzedenz
   `slice-089`) oder eine ausdrückliche §7-Entscheidung fehlt.
2. **V-3** — die Häkchen LP1–LP3, `make gates`, Review, Verifikation sind nach
   Schritt 18/21 nachzuziehen (verkörperte Regel `BEO-PGC/dod-checkbox-nachzug`).
3. Die **drei §6-Risiko-Ausgänge** (`:205`, `:208`, `:212` stehen auf `<…>`) —
   Risiko 2 hat mit #17 seinen Beleg und mit V-5 seine Adresse.
4. Die **Closure-Notiz §7** mit Lerneintrag (9 von 10 Zeilen Platzhalter) und
   die **Register-Zeile** (im Range keine Register-Datei geändert).
5. **V-1** — wird in §7 eine Coverage-Zahl genannt, trägt sie ihren Lauf **und**
   den Hinweis, dass dieselbe Messung am selben Stand 74,70 % oder 74,80 %
   ergeben kann; die Beobachtung selbst gehört mit ihren zwei Adressen
   (`zahl-in-traeger-driftet-gegen-die-messung` 5×,
   `test-integration-retention-timing-flake` 1×) benannt.

Nach 1, 2 — und der Behandlung von 3–5 als regulärer Closure-Arbeit — ist der
Slice `done/`-fähig. Kein Liefer-Defekt und kein rotes Gate: die offenen
Punkte sind Übergebenes, keine gebrochene Zusage.

---

**Beleg-Lage dieses Reports:** Jede Zahl dieses Berichts stammt aus einem der
24 Läufe in §1, je in eigener Werkzeug-Beauftragung gefahren; der Gate-Lauf
(#1) und seine Auswertung waren **zwei** Schritte, sein Exit-Code wurde aus
einer separaten Datei gelesen, nie durch eine Pipe (§3.9). Der Baum ist nach
diesem Bericht sauber bis auf diese Datei; **kein** Commit, keine Änderung an
Artefakten des Slice.
