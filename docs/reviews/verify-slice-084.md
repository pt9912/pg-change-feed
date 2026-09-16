# Verifikations-Report: slice-084 — 2026-09-16

**Verifikations-Art:** DoD-/Entscheidungs-Konformität (Modul 11) — geprüft gegen
**Plan, DoD und die im Plan/Report referenzierten Entscheidungen**. Der
Plan-vs-Code-Diff ist unten geführt. **Nicht** Gegenstand: die Validierung
(„bauen wir das Richtige?", repo-extern) und eine erneute Review des Diffs
(Modul 8 — anderer Eingabe-Kontext, `review-slice-084.md` ist Eingang, nicht
Prüfgegenstand).

**Gegenstand:** `slice-084`
(`docs/plan/planning/in-progress/slice-084-postgresack-naht.md`) — die Naht in
`internal/adapters/driven/postgresack`.

**Verifikations-Range:** `fb6adf6..ea971d3` (Lieferung: `eecb1d9` Naht ·
`8e53d2a` Plan-Nachzug · `26d4ef0` Digest · `4035ee7` §3 · `5ff1ffa` Review ·
`a3cb2c4` F-2 · `db42231` Fixrunde F-1/F-3 · `ea971d3` DoD-Haken). **Der Kopf
bewegte sich während der Verifikation:** auf `main` erschien zusätzlich
`d9a45d7` (`welle-20.md`, Start-Trigger und Nenner) aus einem parallelen
Planner-Zug — sie berührt keinen Träger dieses Slice (nur `welle-20.md`) und
wurde in den Range-Aussagen mitgeführt. Der Arbeitsbaum trägt daneben die
**uncommittete Closure-Arbeit des Planners** (§6/§7 des Plans,
Register-Belege); sie ist **kein** Gegenstand dieses Reports und wurde **nicht**
mitgenommen.

**Eingang:** DoD-Bestätigung des Implementers · `review-slice-084.md`
(0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO; Fixrunde `db42231`) ·
[`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) ·
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 ·
[`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
· [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md).

**Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16 · Docker-only.

---

## 1. DoD-Konformität je Kriterium (§2, im heutigen Text)

| # | DoD-Kriterium | Status | Beleg (eigene Messung) |
|---|---|---|---|
| LP1.1 | Logik hängt an **adapter-eigener** Schnittstelle (eine Methode); konkreter `*pgconn.PgConn` bleibt in Hülle und Dial | **bestätigt** | `seam.go`: `standbySender` mit **einer** Methode `SendStandbyStatusUpdate(ctx, pglogrepl.StandbyStatusUpdate) error`; `connSender{conn *pgconn.PgConn}` hält den Typ als **Feld**; `ack.go` ruft `a.sender.SendStandbyStatusUpdate(...)`. `New(conn *pgconn.PgConn, …)` unverändert |
| LP1.2 | **Kein Verhaltens-Change**: reale Verdrahtung unverändert durch `make test-replication` (Exit 0) | **bestätigt** | `make test-replication` **Exit 0** (unpip­ed); zusätzlich `make test-store` Exit 0, `make test` Exit 0 (31 Pakete `ok`, kein `FAIL`/`SKIP`) |
| LP2.1 | netzlos prüfbarer Teil (Standby-Status-Form, LSN-Form) als **reine Funktion** mit eigenen Tests | **bestätigt** | `ackLSN` (Null-Positions-Grenze + LSN-Form), `standbyStatus` (drei LSN-Felder), `replicationClass` (Wrapping) in `ack.go`; Tests `TestAckLSNRejectsZeroPosition`, `TestAckLSNCarriesOffset`, `TestStandbyStatusCarriesPositionInAllThreeLSNs`, `TestReplicationClassWrapsCause` in `seam_test.go` |
| LP2.2 | Fake-Seite fährt die **Verklebung** — **kein** Ersatz der realen Tests | **bestätigt** | `fakeSender` erfüllt dieselbe Schnittstelle (`var _ standbySender = (*fakeSender)(nil)`); `seam_test.go:1-7` benennt den Ausschluss; die realen Tests bleiben unverändert (s. u.) |
| LP3.1 | Null-Befund geführt: **(a)** `k_ab = 0` · **(b)** Paket-Diff ohne Trägerwechsel · **(c)** kein Verhalten verloren | **bestätigt** (Substanz; zu **(b)** siehe V-1) | eigene Messung §2 |
| LP3.2 | `make gates` grün (Exit direkt, ungepiped) | **bestätigt** | `make gates` **Exit 0** — `coverage-gate: OK — Coverage 71.90% erfüllt Schwelle 70%` · d-check 709 Dateien / 0 Befunde · `commit-traceability: OK` · a-check 0 · `baseline-verify: v6.5.0 OK — 54 Dateien` |
| LP3.3 | Review durchgeführt, Report unter `docs/reviews/` liegt vor; Fixrunde 0 HIGH/0 MEDIUM | **bestätigt, mit Anmerkung** | `docs/reviews/review-slice-084.md` (Commit `5ff1ffa`), Summary 0/0/2/2. F-1 **geschlossen** (Zahlen decken sich mit meiner Messung, §2), F-3 **geschlossen** (`go list -f '{{.Doc}}'` zeigt nur den `ack.go`-Text; Leerzeile vor `package`), F-2 **weitgehend** (`a3cb2c4` trägt die Range nach — der Satz daneben stimmt mit seinem eigenen Befehl noch nicht überein, **V-1**) |
| LP3.4 | *Falls* die Rampe bewegt: Transfer-Nachweis in `harness/sensors/…` nachgezogen, **ohne** neue Schwellen-ADR | **bestätigt (Bedingung nicht eingetreten)** | kein Transfer (§2); `THRESHOLD ?= 70` und `DB_COVERAGE_THRESHOLD=${…:-70}` unverändert, Endstufen je 80 % — keine Schwellen-ADR fällig. Die Sensor-Docs wurden dennoch nachgezogen (F-1), **ohne** Schwelle zu ändern |
| §2.5 | Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt / nicht fällig** | Closure-Pflicht des Planners; §7 ist im Baum **uncommittet** vorbereitet — per Auftrag **nicht** Prüfgegenstand |
| §2.6 | Reconciliation-Register `../reconciliation.md` fortgeschrieben | **n/a** | Datei existiert nicht (`docs/plan/planning/reconciliation.md` fehlt) — Repo ohne Brownfield-Bootstrap; das DoD-Item entfällt nach seinem eigenen Wortlaut |
| §2.7 | Beobachtungs-Register fortgeschrieben | **erfüllt, uncommittet** | `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code/` neu (observation/state/evidence) · `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-084.md` neu; `state.md` auf **2×**. Kein Zähler gesetzt (abgeleitet) — bestätigt |
| §2.8 | Jedes Risiko aus §6 mit genau **einem** Ausgang | **erfüllt, uncommittet** | alle drei §6-Risiken: *entfallen — gestrichen mit Begründung* (einer der drei Ausgänge der geschlossenen Menge) |
| §2.9 | Drei Paarungen getragen | **n/a hier** | Repo führt **Wellen-Betrieb**; die Prüfung fällt der `welle-20`-Closure zu (Modul 6 Schritt 3c) — wie der Plan sagt |

**Ergebnis:** Sämtliche **Liefer**-Kriterien (LP1–LP3, `make gates`, Review) sind
im heutigen Text erfüllt. Die vier offenen Kästchen sind **Closure-Pflichten**
des Planners (Notiz, Register, Risiko-Ausgänge, Paarungen); sie sind teils im
Baum vorbereitet, aber per Auftrag nicht mein Gegenstand. Kein Liefer-Kriterium
hängt an ihnen.

---

## 2. Der Null-Befund (a)/(b)/(c) — selbst nachgemessen

**Alle Zahlen sind eigene Messungen** (Dedup über die Block-Position, dieselbe
Regel wie `tools/harness/db-coverage.sh`), Exit-Codes ungepiped.

| Gegenstand (Statements) | mein Messwert „nach" | mein Messwert „vor" (`fb6adf6`) | Δ |
|---|---|---|---|
| netzlos prüfbare Fläche (Unit) | **1903** (zwei Läufe desselben Baums: 1369 → 71.90 %, 1371 → 72.00 %) | **1903** (1371 gedeckt, 72,04 %) | **0** |
| DB-Adapter-Gegenstand (gemergt) | **659** (491 gedeckt, 74,51 %) | 650 (Plan/ADR-Referenz) | **+9** |
| davon `postgresack` | **32** (32 gedeckt) | **23** (direkt gemessen) | **+9** |
| davon `postgresstorage` | 472 (347 gedeckt) | 472 | 0 |
| davon `replication/receive` | 155 (112 gedeckt) | 155 | 0 |

### (a) `k_ab = 0` — bestätigt, **direkt** gemessen
- Der **Unit-Nenner ist 1903 auf beiden Bäumen.** Ich habe dazu die
  `coverage`-Stufe des Dockerfiles **am Parent-Commit `fb6adf6`** in einem
  Wegwerf-Worktree gebaut und ihr `/out/coverage.out` dedupliziert
  (`1371/1903`) sowie gegen das Profil des heutigen Baums (`1369/1903`)
  gestellt: **die Mengen der Paket-Statements sind zeichengleich** (`diff` der
  Per-Paket-Nenner leer), der Nenner ist 1903 beidseitig.
- Das Unit-Profil enthält **keine einzige** `postgresack`- oder
  `replication/receive`-Zeile (Grep = 0 auf beiden Bäumen) — die Naht bleibt im
  ausgenommenen Paket.
- Der DB-Nenner **wächst** um **+9**: `postgresack` **23 → 32**, direkt über
  `go test -run='^$' -coverpkg=./internal/adapters/driven/postgresack` auf
  beiden Bäumen gemessen (ohne DB, nur Instrumentierung). `postgresstorage` 472
  und `replication/receive` 155 auf beiden Seiten unverändert (die Änderung ist
  ausschließlich `postgresack/**`, s. (b)).
- **Die ±2-Schwankung ist real und unabhängig bestätigt — auf demselben Stand:**
  zwei `make gates`-Läufe **desselben** Baums druckten `71.90 %` (dedup
  **1369/1903**) und `72.00 %` (dedup **1371/1903**); der Lauf am Parent-Commit
  maß ebenfalls 1371. Genau die zwei Werte, die beide Sensor-Docs als beobachtete
  Spanne führen — **auf demselben Baum beobachtet**. Der Nenner (1903) ist davon
  unberührt.

### (b) Paket-Diff — kein Trägerwechsel: **Substanz bestätigt**, Belegform mit Rest (V-1)
- `git diff --name-status fb6adf6..HEAD` (= `…ea971d3` plus `d9a45d7`) listet
  **nur**: `internal/adapters/driven/postgresack/{ack.go,seam.go,seam_test.go}`,
  `harness/image-hash.txt`, die zwei Sensor-Docs, den Slice-Plan, den
  Review-Report und `welle-20.md`. **Kein** Produktcode außerhalb `postgresack`.
- **Die Gegenstands-Listen sind diff-frei** — eigene Prüfung:
  `git diff fb6adf6..HEAD -- Dockerfile tools/harness/db-coverage.sh harness/mk/coverage.mk internal/bootstrap .a-check.yml spec harness/conventions`
  ist **leer**. Kein Paket wechselt zwischen den zwei Gegenständen.
- **V-1 (LOW):** der zitierte Beleg-Befehl in §3(b) trägt seinen Satz noch nicht
  (siehe Findings).

### (c) Kein Verhalten verloren: **bestätigt**
- Die realen, dienst-gestützten Läufe sind grün: `make test-store` **Exit 0**
  (Teilzahl 347/472 = 73,52 %), `make test-replication` **Exit 0**
  (`DB-Adapter-Coverage: 74.51% (gedeckt 491 von 659)`, `db-coverage: OK`),
  `make test` **Exit 0**.
- **Kein Testfall entfernt, keine Zusicherung abgeschwächt:**
  `git diff fb6adf6..HEAD -- internal/adapters/driven/postgresack/ack_test.go`
  ist **leer** (byte-identisch); der `*_test.go`-Diff zeigt genau **eine neue**
  Datei (`seam_test.go`), **keine** gelöschte; im Diff kein neues `t.Skip`,
  `//nolint`, `#noqa` oder `|| true` (Muster-Grep nur auf Prosa des Reports).

---

## 3. Kein Verhaltens-Change — eigene Gegenprobe an der Zusage

Die Zusage „der Umbau ist verhaltensgleich" hängt an den Wächtern, nicht an der
Behauptung. Ich habe **zwei frische Mutationen** gefahren — an Zusagen, die
**weder** der Implementer (`WALApplyPosition: 0`, `replicationClass` ohne
Wrapping) **noch** der Reviewer (nil-Grenze, Aufrufzähler, `New`-nil,
nicht-delegierende Hülle) mutiert haben (Wegwerf-Kopien, `-network none`):

| # | Mutation | Exit | Befund |
|---|---|---|---|
| V-M1 | `ackLSN`: `pglogrepl.LSN(position.Offset)` → `… + 1` (Offset wandert **verändert** in die LSN) | **1** | `FAIL: TestAckLSNCarriesOffset` (LSN 4712, erwartet 4711) **und** `FAIL: TestAcknowledgeSendsStandbyStatus` (101/101/101) |
| V-M2 | `replicationFailure`: der Fehler-Log wird **zweimal** abgesetzt | **1** | `FAIL: TestAcknowledgeWrapsSenderFailure` (Fehler-Log-Aufrufe 2, erwartet 1) |

Beide Zusagen sind **scharf** — die netzlosen Tests prüfen sie real. Zusammen mit
den vier Reviewer- und zwei Implementer-Mutationen aus dem Eingang deckt das die
tragenden Zusagen des Slices ab, ohne die bereits genutzten zu wiederholen.

---

## 4. Die Zahlenträger — stimmen sie, und trägt die Form?

**Zahlen stimmen** (jede gegen meine eigene Messung, §2):
`db-adapter-coverage.md`: Nenner **659** (472/32/155) ✓ · gedeckt **491** ✓ ·
**74,51 %** ✓ · `postgresack` **32** ✓ · `replication/receive` **112/155** ✓.
`coverage-gate.md`: Nenner **1903** ✓ · gedeckt **1369** ✓ · **20/16** `mapper`
✓ · `30/32` ✓ · `31/472` ✓ · `11/155` ✓ · die Rückrechnung
`(1369+30)/(1903+32)=72,30 %`, `1400/2375=58,95 %`, `1380/2058=67,06 %` ✓
(nachgerechnet). Der F-1-Befund ist damit **substantiv geschlossen**.

**Form-Urteil: die Form trägt — sie verschiebt die Drift nicht, sie reduziert
sie.** Der Kern der gewählten Form ist richtig und wird von meiner Messung
gedeckt:
- Der **Nenner** (1903 / 659) ist über Läufe und Bäume **stabil** — ich habe
  1903 auf **zwei** Bäumen und 659 gegen die Per-Paket-Zerlegung gemessen; er
  ist zu Recht als **Zustand** geführt.
- Die **gedeckte Zahl** (1369 / 491) schwankt real (ich maß selbst 1369 **und**
  1371 auf demselben Stand) — sie zu Recht als **lauf-gebundenen Beleg** zu
  führen ist keine Kosmetik, sondern die einzige Form, in der die Angabe nicht
  lügt. Ohne diese Trennung wäre jede gedeckte Zahl eine ablaufende
  Gegenwartsaussage — genau die Klasse des Registers.
- Die **Grenze ist ehrlich benannt:** beide Docs sagen selbst, dass es für
  Zahlen in Trägern **keinen Sensor** gibt; die Falsifikation ist die Messung.
  Das deckt sich mit meinem Befund: die vier driftenden Werte waren nur durch
  Nachmessen zu finden.

**Zwei Reste (INFO, nicht blockierend, V-2/V-3):**
- **V-2:** die Einstiegs-Zelle der Kalibrierungs-Bindung in
  `db-adapter-coverage.md:81` stellt **477 von 650** neben „(die geltende Größe
  des Gegenstands: §Zählbasis)" — während §Zählbasis **659** als
  „Zustandsgröße" führt. Auf der intendierten Lesart („die geltende Größe steht
  in §Zählbasis") trägt der Satz; ein unvorsichtiger Leser kann die Klammer
  aber als Aussage „650 ist geltend" nehmen, und dann widersprechen sich zwei
  Stellen derselben Datei. Ein Wort („…: siehe §Zählbasis") entfernt die
  Zweideutigkeit.
- **V-3:** die beiden eingefrorenen **Rot-/Grün-Belege** (`db-adapter-coverage.md`
  §Ausgabe, `coverage-gate.md` §Ausgabe) drucken weiter **73,38 %** bzw.
  **71,30 %** und nennen — anders als die von denselben Dokumenten aufgestellte
  Form-Regel („jeder Beleg nennt seinen Lauf") — keinen Lauf/Datumsanker. Als
  eingefrorener Lauf-Beleg ist der Wert zulässig; die fehlende Lauf-Benennung
  macht ihn aber von einer Zustandsaussage ununterscheidbar.

**Fazit:** Die Form ist ein **echter Fortschritt**, kein Platzwechsel der Drift:
sie macht die stabile Größe zur Aussage und die wandernde zur datierten
Beobachtung. Die zwei Reste sind Prosa-Präzision ohne Sensor — sie gehören in
die Closure-Notiz oder einen Nachzug, sind aber **kein** Grund, den Slice
zurückzuhalten.

---

## 5. Die neue Beobachtung — Zahl und Einordnung

**Zahl trägt.** Eigener voller `-coverpkg`-Lauf über **alle** Pakete, netzlos
(`-network none`, `-coverpkg=./...`, Exit 0), dedupliziert:

| Paket (netzlos) | gedeckt/Statements |
|---|---|
| `internal/adapters/driven/postgresack` | **30/32** |
| `internal/adapters/driven/postgresstorage` | 31/472 |
| `internal/adapters/driving/replication/receive` | 11/155 |
| `…/postgresstorage/mapper` | 16/20 |

**30 von 32** ist **bestätigt** — und die drei Zahlen daneben (31/472, 11/155,
16/20) sind genau die, die `coverage-gate.md` §Grenze Punkt 4 nun führt. Die
Verdünnung ist real: `postgresack` steht im DB-Gegenstand bei **32/32**, und
**30** davon deckt ein Test, der keine PostgreSQL-Instanz berührt.

**Einordnung als Spiegelbild zu `slice-081` — trägt.** `slice-081` hat den
DB-Nenner **gedräniert** (ein Transfer: `788 → 650`, −138, Quote
`75,25 % → 73,38 %`, [`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
§Kontext (1)); hier **wächst** der Nenner (`650 → 659`, +9) und die Quote
**steigt** (`73,38 % → 74,51 %`) **ohne** einen neuen DB-gestützten Beleg. Beide
Bewegungen entfernen die Zahl von dem, was sie messen soll — einmal nach unten,
einmal nach oben; `ADR-0078` hält nur für die Abwärtsbewegung einen Riegel
bereit. Die Metapher trägt, weil sie an den gemessenen Zahlen hängt.

**Klasse ist eigenständig.** Die drei benachbarten Registereinträge prüfen
anderes: `zahl-in-traeger-driftet-gegen-die-messung` = eine **geschriebene
Zahl** driftet (hier driftet die **Messung** selbst gegen ihren Zweck);
`endstufe-unter-eigenem-messgegenstand-unerreichbar` = ein **Ziel** ist
unerreichbar (hier ist ein erreichtes Ziel wenig wert);
`adapter-unittest-verdeckt-bootstrap-luecke` = eine **Verdrahtungs**-Lücke.
Die Neuanlage ist gerechtfertigt; Sub-Area `*`/`PGC` stimmt mit der
Modus-Deklaration (`harness/conventions.md`) überein; `state.md` = `offen (1×)`
mit abgeleitetem Zähler ist form-korrekt.

**Eine Nuancierung:** `observation.md` sagt, „`ADR-0080` hat die Verdünnung …
als Trigger benannt; sie ist mit diesem Vorgang **eingetreten**." Die
**Verdünnung** ist eingetreten (30/32 gemessen); ob `ADR-0080` Trigger (a)
(„wird sie **materiell**") **ausgelöst** ist, ist ein **Urteil**, und die
Eintragung entscheidet es zu Recht **nicht** — sie führt es als Architect-Frage
(`state.md`). Ein Leser könnte „sie ist eingetreten" auf den **Trigger** beziehen;
die Eintragung sollte zwischen „Phänomen eingetreten" und „Triggerschwelle
entschieden" nicht verschwimmen. **Nicht blockierend.**

---

## 6. `ADR-0044` — der Digest-Träger

**Bestätigt.**
- Der Zug ändert **Build-Kontext** (`internal/**` liegt über `.dockerignore` im
  Kontext: `*` mit `!internal/`, `!cmd/`, `!go.mod`, `!go.sum`) — die
  welle-1-Regel („`make image` vor der Closure") greift, `ADR-0044` §3.
- `harness/image-hash.txt` ist der **richtige Träger** (`ADR-0044` §1: Digest als
  **Lauf-Beleg**, nicht als Inhalts-Fingerabdruck). Der Wechsel
  `sha256:f4475e1e…` → `sha256:4bd43435…` liegt als **eigener** `chore(image)`-Commit
  `26d4ef0` vor (eine Datei, Betreff mit `ADR-0044`/`ADR-0080`).
- **Eigene Gegenprobe:** ein frischer `docker buildx build --load` dieses Baums
  (Metadaten nach `/tmp`, der getrackte Träger blieb unberührt — `git status`
  leer) liefert **denselben** Digest `sha256:4bd43435…`. Das ist der Befund
  „unveränderter Digest im **selben** Builder" — **kein** Inhalts- und **kein**
  Umgebungs-Beweis (`ADR-0044` §1/§3). Genau so, und nicht stärker, ist er zu
  lesen.

---

## 7. `AGENTS.md` §3.10 — greift nicht

**Geprüft, nicht angenommen:** der Slice berührt **keine** Workflow-Datei.
`git diff --name-only fb6adf6..HEAD -- .github/workflows/` ist **leer**, und
`git diff --name-status fb6adf6..HEAD` nennt für `.github/` keinen Eintrag. Damit
ist §3.10 („ein neuer oder strukturell geänderter Workflow gilt erst nach einem
realen grünen Post-Push-Lauf als abgeschlossen") **nicht anwendbar** — es gibt
kein Post-Push-Erfordernis aus diesem Slice.

---

## 8. Entscheidungs-Konformität

- **[`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) —
  konform.** Festlegung 1/5: **Hülle + Schnittstelle, keine Zusicherung** —
  `connSender` delegiert 1:1, `var _ standbySender = connSender{}` steht **im
  Paket** (nicht im Test), Fake und Hülle erfüllen dieselbe Schnittstelle.
  Festlegung 3: der konkrete Typ bleibt in Hülle und Dial. Festlegung 4: der
  Nachweis ist als **Null-Befund** geführt, `k_ab = 0` (§2). Festlegung 7:
  `.a-check.yml` unverändert (diff-frei). Kein Unterpaket, kein Paketwechsel —
  der in §„Bestätigt im Buchstaben, widerlegt in der Folge" benannte
  Nicht-Blocker tritt ein.
- **[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 5 — konform.** Die Naht zieht **nicht** in die prüfbare Fläche (Unit-Nenner
  1903 → 1903), sie macht die Logik nur netzlos **testbar**; der konkrete Typ
  bleibt hinter der Naht. Genau die von `ADR-0080` §„Die benannte Grenze"
  beschriebene Wirkung.
- **[`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  — konform, als Null-Fall.** `k_ab = 0`, **kein** Subjekt-Transfer, **keine**
  Neu-Bemessung; der Regressions-Riegel (`k_auf < k_ab`, gesunkene Quote bei
  unverändertem Nenner) hat kein Objekt. Die aggregierte Summe ist als
  Allein-Beleg entwertet — dieser Slice stützt sich nicht auf sie, sondern auf
  den Paket-Diff (b) und die grünen realen Läufe (c).
- **[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) — konform** (§6).
- **`AGENTS.md` §3.6** (Schwellen nur per ADR): keine Schwelle geändert
  (`THRESHOLD 70`, `DB_COVERAGE_THRESHOLD 70`, Endstufen 80) — keine Schwellen-ADR
  fällig. **§3.5** (Accepted-ADRs immutable): keine ADR inhaltlich überschrieben
  (`ADR-0071`/`0078`/`0080` diff-frei bzw. unberührt). **§3.11**: d-check
  `hostpaths` über 709 Dateien 0 Befunde.

---

## 9. Findings (Verifier)

### V-1 — §3(b) zitiert einen Befehl, dessen Ausgabe den Satz daneben nicht deckt (LOW, **Rest von F-2**)

- `pfad`: `docs/plan/planning/in-progress/slice-084-postgresack-naht.md` §3(b)
- `befund`: Der Satz lautet: „`git diff --name-only fb6adf6..4035ee7` listet
  **ausschließlich** Dateien unter `internal/adapters/driven/postgresack/`".
  Der Befehl liefert **fünf** Pfade;
  gemessen:
  `docs/plan/planning/in-progress/slice-084-postgresack-naht.md`,
  `harness/image-hash.txt` und die drei `postgresack`-Dateien. Die **Aussage**
  (kein Trägerwechsel) ist wahr und von mir unabhängig bestätigt (§2 (b)) — die
  **zitierte Form** trägt sie aber nicht. `a3cb2c4` hat die **Range**
  nachgetragen (der Kern von F-2), aber nicht den **Pathspec** (`--
  internal/adapters/driven/postgresack/`) und den Satz nicht auf den
  Code-Umfang eingeschränkt. F-2 ist damit **weitgehend**, nicht vollständig
  geschlossen.
- `verifizierbar`: ja — `git diff --name-only fb6adf6..4035ee7` (5 Pfade) gegen
  denselben Befehl mit Pathspec (3 Pfade).

### V-2 — Zweideutige Klammer in der Kalibrierungs-Zelle (INFO)

Die Einstiegs-Zelle stellt **477 von 650** neben „(die geltende Größe des
Gegenstands: §Zählbasis)", während §Zählbasis **659** führt (§4). Ein Wort
(„…: siehe §Zählbasis") entfernt die Zweideutigkeit.

### V-3 — Eingefrorene Rot-/Grün-Belege ohne Lauf-Benennung (INFO)

Siehe §4. Zulässig als Lauf-Beleg; die fehlende Lauf-/Datums-Angabe widerspricht
der selbst aufgestellten Form-Regel „jeder Beleg nennt seinen Lauf".

**Kein HIGH, kein MEDIUM.** V-1 ist ein Zitat-/Form-Rest (der Beleg-Befehl, nicht
die Aussage), V-2/V-3 sind Prosa-Präzision. Keiner berührt die Substanz des
Null-Befunds, des Verhaltens-Change-Freiheits-Belegs oder der Nahtform.

---

## 10. Was ich **nicht** prüfen konnte (Grenzen)

- **Die Zahl `postgresack` 23 (18 gedeckt) „vor dem Zug"** (Plan §3, DB-Zeile):
  den **Nenner** 23 habe ich direkt gemessen, den **gedeckten** DB-Wert 18 nicht
  — das verlangte einen `make test-store`/`make test-replication`-Lauf am
  Parent-Commit (reale PostgreSQL). Für den Null-Befund irrelevant: die
  tragende Aussage ist der **Nenner** (+9), und der ist gemessen.
- **`make doc-commits` / `make doc-immutable`** existieren in diesem Repo
  **nicht** (Makefile und `harness/mk/*.mk` haben kein solches Target; das
  `vcs`-Modul ist in `.d-check.yml` nicht in `modules` aktiviert). Ersatzweise
  habe ich die **repo-eigenen** Äquivalente gefahren: `make commit-traceability`
  (Exit 0; `RANGE=fb6adf6..HEAD` → 9 Commits, Betreffs ohne Struktur-ID) und den
  vollen `make docs-check`-Modulsatz (inkl. `matrix`-Status) aus `make gates`.
  **MR-Immutabilität** hat hier **kein Objekt**: die Range berührt
  `harness/conventions/**` nicht.
- **Der Digest-Beweis bleibt lauf-/builder-gebunden** (`ADR-0044` §3): mein
  Rebuild lief im **selben** Builder. Er belegt „der eingetragene Wert stammt
  aus einem realen Bau **dieses** Baums", **nicht** Inhalts- oder
  Umgebungs-Gleichheit.
- **Der Tier-Lauf `test/integration`** (Compose/E2E) ist nicht Teil meiner
  Werkzeuge und wurde nicht gefahren; `make test-replication`/`make test-store`
  decken den DB-gestützten Belegteil dieses Slice.
- **Die Closure-Notiz (§7) und die §6-Ausgänge** habe ich **gelesen** (sie
  liegen uncommittet im Baum), aber **nicht** als Lieferung bewertet — die
  Closure kommt nach mir.

---

## 11. Lauf-Artefakte und Hygiene

- `tools/schema/plan.yaml` wurde nach `make test-store`/`make test-replication`
  **zurückgenommen** (`git checkout`); der Baum trägt danach **nur** die
  uncommittete Planner-Closure (`slice-084`-Plan §6/§7, `state.md`,
  Beleg-Verzeichnisse).
- Alle Läufe Docker-only, Exit-Codes **direkt und ungepiped** ermittelt. Die
  Wegwerf-Kopien lagen außerhalb des Baums; das Wegwerf-Worktree am
  Parent-Commit wurde entfernt (`git worktree list` = nur das Hauptverzeichnis).
- Der getrackte Digest-Träger blieb beim Rebuild unberührt (`git status` leer
  für `harness/image-hash.txt`).

---

## 12. Verdikt

**Der Slice ist closure-fähig.**

- Jedes **Liefer**-Kriterium der DoD ist im heutigen Text erfüllt; die vier
  offenen Kästchen sind Closure-Pflichten des Planners.
- Der **Null-Befund (a)/(b)/(c)** ist in allen drei Teilen **unabhängig
  nachgemessen** — `k_ab = 0` ist direkt belegt (Unit-Nenner 1903 auf beiden
  Bäumen, `postgresack` 23 → 32), der Paket-Diff zeigt keinen Trägerwechsel, kein
  Testfall ist entfernt, die realen Läufe sind grün.
- **Kein Verhaltens-Change** — zusätzlich durch zwei **eigene** Mutationen an
  bisher unberührten Zusagen scharfgestellt.
- Die **Zahlenträger** stimmen mit meiner Messung; die gewählte **Form**
  (Nenner als Zustand, gedeckte Zahl als Lauf-Beleg) **trägt** und verschiebt die
  Drift nicht, sie reduziert sie — mit zwei benannten Prosa-Resten (V-2/V-3).
- Die **neue Beobachtung** ist gemessen (30/32 netzlos) und als eigenständige
  Klasse korrekt eingeordnet; die Spiegelbild-Metapher zu `slice-081` trägt.
- `ADR-0080`, `ADR-0071` Punkt 5, `ADR-0078` (Null-Fall) und `ADR-0044` sind
  eingehalten; §3.10 greift **nicht** (kein Workflow berührt).

**Mitgabe an die Closure (nicht blockierend):** **V-1** (den §3(b)-Satz auf
seinen Code-Umfang einschränken oder den Pathspec zitieren — der letzte Rest von
F-2) sowie **V-2/V-3** als Prosa-Nachzug. Ferner: die Eintragung
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` sollte „Phänomen
eingetreten" und „Triggerschwelle entschieden" sprachlich trennen (§5).

---

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, diese Messungen, dieses
Modell). Er ersetzt **keine** Validierung (Modul 11, repo-extern) und **keine**
Review — der Diff ist gegen Plan und Entscheidungen bereits geprüft
(`review-slice-084.md`).
