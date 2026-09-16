# Verifikations-Report: slice-085 — 2026-09-16

**Verifikations-Art:** DoD-/Entscheidungs-Konformität (Modul 11) — geprüft gegen
**Plan, DoD und die im Plan/Report referenzierten Entscheidungen**. Der
Plan-vs-Code-Diff ist unten geführt. **Nicht** Gegenstand: die Validierung
(„bauen wir das Richtige?", repo-extern) und eine erneute Review des Diffs
(Modul 8 — anderer Eingabe-Kontext; `review-slice-085.md` ist **Eingang**, nicht
Prüfgegenstand).

**Gegenstand:** `slice-085`
(`docs/plan/planning/in-progress/slice-085-receive-naht.md`) — die Naht in
`internal/adapters/driving/replication/receive` (Treiber-Hülle + schmale,
fake-fähige Fläche, **im Paket**, `ADR-0080`).

**Verifikations-Range:** `95480a5..33440d3` (= `HEAD`). Lieferung: `610751c`
Naht · `506c811` Plan/Null-Befund · `1abc7c5` Digest · `7e0e937` Review ·
`dd29d83`/`61bb5cf` Fixrunde (F-1 Zahlenträger, F-2 Anteil) · `b83217c`
Fixrunde (Auftragserweiterung `coverage-gate.md`/`db-coverage.sh`) · `33440d3`
§3 „Grenze der Liste selbst".

**Eingang:** DoD-Bestätigung des Implementers · `review-slice-085.md`
(0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO) ·
[`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) ·
[`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md) ·
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) · `ADR-0071` Punkt 5.

**Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16 · Docker-only.

---

## 1. DoD-Konformität je Kriterium (§2, im heutigen Text)

| # | DoD-Kriterium | Status | Beleg (eigene Messung) |
|---|---|---|---|
| LP1.1 | Logik hängt an **adapter-eigener** Schnittstelle (**sieben** Operationen); der konkrete `*pgconn.PgConn` bleibt in Hülle und Dial | **bestätigt** | `seam.go`: `driverSession` mit genau sieben Operationen (`IdentifySystem`, `CreateReplicationSlot`, `StartReplication`, `SendStandbyStatusUpdate`, `Exec → []*pgconn.Result`, `ReceiveMessage → pgproto3.BackendMessage`, `Close`); `connSession{conn *pgconn.PgConn}` hält den Typ als **Feld**; `receive.go:99` `Stream.session driverSession`, `walretention.go:26` `WALRetentionChecker.session driverSession`; der konkrete Typ im Dial (`connectReplication`) |
| LP1.2 | **Kein Verhaltens-Change**: reale Verdrahtung geht unverändert durch `make test-replication` (Exit 0) | **bestätigt** | `make test-replication` **Exit 0** (unpip­ed) · `make test-store` Exit 0 · `make test` Exit 0 (31 Pakete `ok`, 0 `FAIL`); `stream_test.go` **byte-identisch** (§3) |
| LP2.1 | netzlos prüfbarer Teil als **reine Funktion** mit eigenen Tests | **bestätigt** | die sechs genannten Funktionen existieren und sind je mit eigenem Test belegt: `parseXLogData` (`TestParseXLogDataCarriesPayload`/`…RejectsShortMessage`), `parseKeepalive` (`TestParseKeepaliveReportsReplyRequested`/`…RejectsShortMessage`), `parseLSN` (`TestParseLSNCarriesCatalogValue`/`…RejectsUnreadableValue`), `firstRow` (`TestFirstRowTranslatesCatalogRow`/`…WithoutRowReportsAbsence`), `standbyStatus` (`TestStandbyStatusCarriesPositionInAllThreeLSNs`), `slotLSNQuery` (`TestSlotLSNQueryNamesTheSlot`) |
| LP2.2 | Fake fährt die **Verklebung** — **kein** Ersatz der realen Tests | **bestätigt** | `seam_test.go:1-8` benennt den Ausschluss wörtlich; `var _ driverSession = (*fakeSession)(nil)` (Test) und `var _ driverSession = connSession{}` (Produktion, `seam.go:72`); eigene Mutationen M1–M3 (§3) belegen, dass die Fake-Suite ihre Zusagen **scharf** prüft |
| LP3.1 | Null-Befund geführt: **(a)** `k_ab = 0` · **(b)** Paket-Diff ohne Trägerwechsel · **(c)** kein Verhalten verloren | **bestätigt** | alle drei Teile unabhängig nachgemessen — §2 |
| LP3.2 | `make gates` grün (Exit direkt, ungepiped) | **bestätigt** | `make gates` **Exit 0** — `baseline-verify: v6.5.0 OK — 54 Dateien` · d-check 713 Dateien / 0 Befunde · d-check `commits` 0 Befunde · `commit-traceability: OK — 5 Commit(s)` · `coverage-gate: OK — Coverage 72.00% erfüllt Schwelle 70%` · a-check 0 Befunde |
| LP3.3 | Review durchgeführt, Report unter `docs/reviews/` liegt vor (kein Self-Review) | **Substanz bestätigt, Kästchen offen** | `docs/reviews/review-slice-085.md` (Commit `7e0e937`) liegt vor, Summary **0 HIGH / 0 MEDIUM / 1 LOW / 3 INFO**; die zwei Fixrunden (`dd29d83`…`33440d3`) sind gelandet. Das Kästchen §2 ist **nicht** nachgezogen — der Report nennt den Grund selbst (Rückgabe-Pfeil F-1; Nachzug über den Implementer/Planner) — **Closure-Pflicht**, kein Liefer-Defekt |
| LP3.4 | *Falls* die Rampe bewegt: Transfer-Nachweis nachgezogen, **ohne** neue Schwellen-ADR | **bestätigt (Bedingung nicht eingetreten)** | kein Transfer (§2 (b)); `THRESHOLD ?= 70` und `DB_COVERAGE_THRESHOLD=${…:-70}` an **beiden** Enden der Range identisch (eigene Prüfung), Endstufen 80 % — keine Schwellen-ADR fällig |
| §2.5 | Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt / nicht fällig** | §7 ist im Baum **Template** (`<…>`, `<bei Closure>`); Planner-Arbeit **nach** der Verifikation |
| §2.6 | Reconciliation-Register `../reconciliation.md` fortgeschrieben | **n/a** | Datei existiert nicht (`ls` → „nicht gefunden") — Greenfield-Bootstrap, das Item entfällt nach seinem eigenen Wortlaut |
| §2.7 | Beobachtungs-Register fortgeschrieben | **nicht erfüllt (offen)** | `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code/evidence/` trägt **nur** `slice-084.md` (1×) · `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/` trägt `slice-081.md`,`slice-084.md` (2×) — `evidence/slice-085.md` fehlt in **beiden**; Planner-Pflicht |
| §2.8 | Jedes Risiko aus §6 mit genau **einem** Ausgang | **nicht erfüllt (offen)** | Risiko 1 = *entfallen — gestrichen mit Begründung* ✓; Risiken 2–4 stehen auf `<bei Closure>` — Planner-Pflicht |
| §2.9 | Drei Paarungen getragen | **n/a hier** | Repo führt **Wellen-Betrieb** (`roadmap.md` §Offene Wellen → `welle-20`); die Prüfung fällt der nächsten Wellen-Closure zu (Modul 6 Schritt 3c) — wie der Plan sagt |

**Ergebnis:** Sämtliche **Liefer**-Kriterien (LP1.1–LP3.4, `make gates`, Review)
sind im heutigen Text **erfüllt**; die vier offenen Kästchen sind
**Closure-Pflichten** des Planners (Notiz, Register, Risiko-Ausgänge; das
Review-Häkchen). Kein Liefer-Kriterium hängt an ihnen — dieselbe Lage wie
`slice-084`.

---

## 2. Der Null-Befund (a)/(b)/(c) — selbst nachgemessen

Alle Zahlen sind **eigene Messungen** (Dedup über die Block-Position, dieselbe
Regel wie `tools/harness/db-coverage.sh`, eigene `awk`), Exit-Codes ungepiped,
Docker-only. Die „vor"-Werte sind an einem **Wegwerf-Worktree** `95480a5`
gemessen, nicht übernommen.

| Gegenstand (Statements) | vor (`95480a5`) | nach (`HEAD`) | Δ |
|---|---|---|---|
| netzlos prüfbare Fläche (Unit, `coverage`-Stufe) | **1903** (1371 gedeckt, 72,04 %) | **1903** (1371 gedeckt, 72,04 %) | **0** |
| DB-Adapter-Gegenstand (gemergt) | **659** (491 gedeckt, **74,51 %**) | **691** (532 gedeckt, **76,99 %**) | **+32** |
| davon `postgresstorage` | 472 (347 gedeckt) | 472 (347 gedeckt) | 0 |
| davon `postgresack` | 32 (32 gedeckt) | 32 (32 gedeckt) | 0 |
| davon `replication/receive` | **155** (112 gedeckt) | **187** (153 gedeckt) | **+32** |

*(Unit: mein Lauf maß **1371** gedeckt — das andere der beiden dokumentierten
Enden 1369/1371 desselben Stands, ±2-Schwankung; der für `k_ab` tragende
**Nenner** ist beidseitig 1903.)*

### (a) `k_ab = 0` — bestätigt, direkt gemessen
- **Unit-Nenner 1903 auf beiden Bäumen.** Ich habe die `coverage`-Stufe des
  `Dockerfile` **am Parent `95480a5`** (eigener `docker build --target coverage`
  über den Worktree, `BUILD_EXIT=0`) und am `HEAD` gebaut, `/out/coverage.out`
  extrahiert und block-position-dedupliziert: **1371/1903** beidseitig;
  der **Per-Paket-Nenner-`diff` ist leer** — dieselbe Mengen.
- Das Unit-Profil enthält **keine Zeile** der drei exakten ausgenommenen Pakete
  (`…/postgresstorage`, `…/postgresack`, `…/replication/receive`) — eigene
  Verzeichnis-Prüfung; die einzigen `postgresstorage`-Zeilen sind die
  **Unterpakete** `…/postgresstorage/mapper` und `…/postgresstorage/sqlexec`,
  die zum Unit-Gegenstand gehören. `k_ab = 0` für die Unit-Fläche ist damit
  **direkt** belegt.
- Der DB-Nenner **wächst** um **+32**: `replication/receive` 155 → **187**
  (direkt über `go test -run='^$' -coverpkg=<drei Pakete>` auf beiden Bäumen
  gemessen, ohne DB); `postgresstorage` **472** und `postgresack` **32** sind
  **byte-stabil**. `k_ab = 0` ist damit die Null-Hälfte der Arithmetik, nicht
  eine Behauptung.

### (b) Paket-Diff — kein Trägerwechsel: **bestätigt**
- `git diff --name-only 95480a5..610751c -- internal/adapters/driving/replication/receive/`
  listet **vier** Dateien (`seam.go`, `seam_test.go`, `receive.go`,
  `walretention.go`); die **Gegenrichtung**
  `… -- internal/ ':!internal/adapters/driving/replication/receive/'` ist
  **leer** — nachgeprüft.
- **Gegenstands-Listen diff-frei:** `git diff --name-only 95480a5..HEAD --
  Dockerfile tools/harness/db-coverage.sh harness/mk/coverage.mk internal/bootstrap
  .a-check.yml spec harness/conventions` nennt **nur** `tools/harness/db-coverage.sh`
  — und dessen **volle** Range-Diff trägt **ausschließlich Kommentarzeilen**
  (`grep` auf Nicht-`#`-Änderungen: leer); `DB_COVERAGE_PKGS` und
  `DB_COVERAGE_THRESHOLD` sind an beiden Enden **zeichengleich**. Der in §3 als
  „Grenze dieser Liste selbst" benannte Fixrunden-Zugriff ist damit genau das,
  was der Plan sagt: Kommentar-/Beleg-Korrektur, **kein** Schnitt.
- **Kein Paket** wechselt zwischen den zwei Gegenständen; `.a-check.yml` und die
  Composition Root sind diff-frei.

### (c) Kein Verhalten verloren: **bestätigt**
- Die realen, dienst-gestützten Läufe sind grün: `make test-store` **Exit 0**
  (Teilzahl 347/472 = 73,52 %), `make test-replication` **Exit 0**
  (`DB-Adapter-Coverage: 76.99% (gedeckt 532 von 691)`, `db-coverage: OK`),
  `make test` **Exit 0**.
- **Kein Testfall entfernt, keine Zusicherung abgeschwächt:**
  `git diff --name-status 95480a5..HEAD -- '*_test.go'` zeigt genau **eine neue**
  Datei (`seam_test.go`), **keine** gelöschte; `stream_test.go` ist
  **byte-identisch** (SHA256 `7b5d6a04…` an beiden Enden, §3).
- **Der öffentliche Rand ist unverändert:** Mengenvergleich der exportierten
  Top-Level-Symbole **und** der Methodensignaturen (`Stream.*`,
  `WALRetentionChecker.*`) beider Bäume ist identisch
  (`NewStream`, `NewWALRetentionChecker`, `Config`, `Stream`,
  `WALRetentionChecker`, `ErrConfiguration`, `ErrReplication`); einzige
  Produktions-Neuzugänge sind **unexportiert** (`handleCopyData`, die Naht).

---

## 3. Kein Verhaltens-Change — eigene Mutationen an unberührten Zusagen

Die Zusage „verhaltensgleich" hängt an den Wächtern. Ich habe **drei frische
Mutationen** in einer Wegwerf-Kopie (`git archive HEAD` unter `/tmp`,
`--network none`) gefahren — an Zusagen, die **weder** der Implementer (Slot
ignoriert · Katalog-Wert · `WALApplyPosition` · neutrale Zerlegung ·
nicht-delegierendes `Exec`) **noch** der Reviewer (Keepalive invertiert ·
`parseXLogData` ohne Byte-ID · `slotLSNQuery` ohne Prädikat) mutiert haben:

| # | Mutation | netzlos Exit | Befund |
|---|---|---|---|
| V-M1 | `parseLSN`: Fehlerklasse `ErrReplication` → `ErrConfiguration` | **1** | `FAIL: TestParseLSNRejectsUnreadableValue` |
| V-M2 | `Run`: `case *pgproto3.CopyDone: return nil` → `return Fehler` | **1** | `FAIL: TestRunStopsOnCopyDone` · `TestRunIgnoresOtherBackendMessages` · `TestRunAnswersKeepaliveWithAcknowledgedPosition` · `TestRunKeepaliveWithoutReplyRequestedSendsNothing` |
| V-M3 | `ensurePublication`: Grenze invertiert (`if !exists` → `if exists`) | **1** | `FAIL: TestEnsurePublicationAcceptsExistingPublication` · `TestEnsurePublicationRejectsMissingPublication` |
| Kontrolle | unveränderte Kopie im **selben** Harness | **0** | `ok …/receive 0.003s` — die Fehlschläge sind den Mutationen zuzurechnen |

Alle drei Zusagen (Fehlerklasse der LSN-Übersetzung, reguläres `CopyDone`-Ende,
Publication-Bestand) sind **scharf**; die Kopie war danach byte-identisch
zurückgesetzt, der Arbeitsbaum nie mutiert.

**Byte-Beleg `stream_test.go`:** SHA256
`7b5d6a04b939de5c7667139cfc38981ce22a67ef9d485641fd229218a9cf3630` an **beiden**
Enden der Range (selbst gerechnet, `git show 95480a5:…` vs. heutiger Baum).

---

## 4. Die Form-Entscheidung der Fixrunde — trägt sie, und ist sie **vollständig** angewandt?

**Die Form** (Fixrunde): *jede Zahl eines Doku-Trägers trägt ihren Ursprung, und
ist sie eine Messung, zusätzlich den **Zeitpunkt**; der **Nenner** hängt am
Code-Stand, die **gedeckte** Zahl am Lauf.*

**Urteil: sie trägt in ihrer tragenden Hälfte — aber sie ist in beiden Dokumenten
NICHT vollständig durchgesetzt.** Ich habe gegen sie gesucht und **bewegliche
Zahlen ohne Zeitpunkt** gefunden.

**Was trägt (bestätigt):**
- `db-adapter-coverage.md` §Zählbasis: der **Nenner** 691 (mit den Anteilen
  472·32·187) ist als **Zustand** geführt und trägt den Lauf `slice-085`; die
  Vorläufer-Zahlen der ersten zwei Punkte sind ausdrücklich als
  **Kalibrierungs-Lauf**-Beispielwerte markiert; die **Kalibrierungs-Zelle**
  (477/650 = 73,38 %) nennt ihren Ursprung. Meine Messung deckt jede dieser
  Zahlen (691/532/76,99 % · 472 · 32 · 187 · 650 als Kalibrierungs-Nenner).
- `coverage-gate.md` §Zählbasis: Nenner **1903** und die gedeckte Zahl
  (1369/1371, 71,94 %/72,04 %) tragen ihren Lauf (`slice-084`, `slice-085`); die
  drei Ausgenommenen (`30/32`, `31/472`, `112/187`) tragen `Lauf slice-085`.
  Meine Messung deckt alle (30/32, 31/472, 112/187 — exakt).
- Die **Abgrenzung Nenner ↔ gedeckte Zahl** ist sachlich richtig: der Nenner ist
  über Bäume/Läufe stabil (ich maß 1903 auf **zwei** Bäumen), die gedeckte Zahl
  wandert (ich maß 1371, die Docs führen 1369 **und** 1371) — die Trennung ist
  die einzige Form, in der die Angabe nicht lügt.

**Was die Form NICHT deckt (die Lücke — eigene Prüfung, jede Zahl gemessen):**
1. **`coverage-gate.md` §Grenze Punkt 1** trägt **bewegliche Messwerte ohne
   jeden Ursprung/Zeitpunkt**: `postgresstorage/mapper` **20/16**,
   `grpc/streamv1` **86/61 (70,9 %)** und `cmd/pg-change-feed` **49/0**. Ich habe
   sie gemessen — sie sind **heute noch korrekt** (16/20 · 61/86 = 70,9 % ·
   0/49), aber sie **nennen weder Lauf noch Code-Stand**. Genau das ist die
   Klasse `zahl-in-traeger-driftet-gegen-die-messung`: sie **können** still
   driften, und kein Sensor fängt sie.
2. **Beide §Ausgabe-Abschnitte** (Rot-/Grün-Beleg) führen die eingefrorenen
   Messwerte **73,38 %** (`db-adapter-coverage.md`) bzw. **71,30 %**
   (`coverage-gate.md`) mit einem **Code-Stand**-Marker („Stand der Naht aus
   [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 5") — aber **ohne Lauf-/Datums-Anker**, obwohl beide Dokumente
   selbst die Regel „jeder Beleg nennt seinen Lauf" aufstellen (dieselbe Stelle,
   die `verify-slice-084` als **V-3** benannt und die Fixrunde stehengelassen
   hat).
3. **`db-adapter-coverage.md` §Zählbasis, Punkt 2** liest „darum trägt das Profil
   **132 Positionen × 2**" im **Indikativ als Gegenwartsaussage**; der
   Kalibrierungs-Vorbehalt steht erst am **Ende** des Punkts. Der Wert ist als
   Kalibrierungs-Zahl **korrekt** (eigene Messung am Kalibrierungs-Stand
   `fb6adf6`: **132** Positionen × 2 = 264 Zeilen), aber **heute** misst dasselbe
   Profil **168 Positionen × 2 = 336 Zeilen** (eigene Messung am `HEAD`).
4. **Nebenpunkt:** der Vorbehalt „Die Zahlen der **zwei ersten** Punkte …
   Beispielwerte, **nicht** die geltende Größe" trifft für die **472** des ersten
   Punkts **nicht** zu — 472 **ist** die geltende Größe (Teil der 691). Der
   Vorbehalt überzieht eine noch gültige Zahl.

**Fazit:** Der **Kern** der Form trägt (Nenner = Zustand mit Lauf, gedeckte Zahl
= Lauf-Beleg, Kalibrierungswerte als solche benannt); die **universale** Fassung
(„wie jede Zahl dieses Dokuments", „jeder Beleg nennt seinen Lauf") wird von den
Dokumenten **selbst nicht erfüllt**. Das ist **kein** DoD-Bruch (die Form ist
nicht DoD-Kriterium dieses Slice), aber ein **sachlicher Rest** der Fixrunde —
Finding **V-1**, LOW. Die Zahlen selbst stimmen (bis auf die „132"-Gegenwarts-
Lesart, die zur Kalibrierung passt, zur Gegenwart nicht).

---

## 5. Die zwei Befunde, die der Implementer **bewusst nicht** geändert hat

### (a) `db-adapter-coverage.md` „**132 Positionen × 2**" — heutiger Messwert **168**
- **Gemessen:** `replication.coverprofile` am `HEAD` hat **168** distinkte
  Block-Positionen (× 2 = **336** Zeilen) — eigene `awk` über das reale Profil
  des `make test-replication`-Laufs; am Kalibrierungs-Stand `fb6adf6` sind es
  **132** (× 2 = 264). Der Wert **war** also korrekt.
- **Urteil: Finding, INFO.** Die Zeile liegt in einem Punkt, den die Fixrunde
  ausdrücklich als **Kalibrierungs-Lauf** deklariert hat (Vorbehalt am
  Punktende) — die Herkunft ist damit **gegeben**. Was bleibt, ist die **Lesart**:
  der Satz ist Indikativ über den **Jetzt**-Zustand des Profils, und der ist
  heute 168. Keine der DoD-Kriterien hängt daran; der Fixrunde-Vorbehalt ist die
  schwächste vertretbare Auflösung, nicht die falsche.
- **DoD-Berührung: nein.**

### (b) `coverage-gate.md` §Grenze Punkt 1 — `go list -f '{{len .TestGoFiles}}'` trägt „fünf" nicht
- **Gemessen** (gepinntes Toolchain-Image, `--network none`, eigene Läufe):
  - `go list -f '{{if eq (len .TestGoFiles) 0}}{{.ImportPath}}{{end}}' $pkgs`
    → **25** Pakete mit `TestGoFiles == 0`;
  - `go list -f '{{if and (eq (len .TestGoFiles) 0) (eq (len .XTestGoFiles) 0)}}{{.ImportPath}}{{end}}' $pkgs`
    → **genau fünf**: `postgresstorage/queries`, `grpc/streamv1`,
    `application/port/inbound`, `domain/errors`, `cmd/pg-change-feed`.
  Der zitierte Befehl liefert mit dem gesetzten Prädikat **25**, nicht fünf — der
  Beleg-Befehl **trägt seinen Satz nicht**. Die **Liste** der fünf im Fließtext
  ist dagegen **korrekt**.
- **Urteil: Finding, LOW.** Dieselbe Klasse wie `verify-slice-084` **V-1**
  („Beleg-Befehl trägt seinen Satz nicht"). Der Satz wurde in `slice-079`
  (`65aead2`) eingeführt und von diesem Slice **nicht** berührt (`git blame`
  zeigt `65aead2`; §Grenze Punkt 1 ist im Fixrunden-Diff diff-frei) — also ein
  **Alt-Bestand** in einer Datei, die der Slice an **anderer** Stelle berührt
  hat. Kein Sensor fängt solche Beleg-Formeln (`d-check` hat kein
  Kommando-zu-Satz-Modul).
- **DoD-Berührung: nein.**

**Beide Befunde berühren die DoD nicht** — sie liegen in Sensor-Prosa, nicht in
einem Liefer-Kriterium. Sie gehören als Nachzug in die Closure-Notiz oder einen
eigenen kleinen Zug.

---

## 6. Der Trigger-Kandidat — die Zahl, und was daraus folgt

- **Gemessen:** die DB-Adapter-Coverage ist **76,99 %** (532/691, eigene Dedup
  des gemergten Profils) — **≥ 75 %**, der nächsten vollen 5-%-Stufe über der
  geltenden `DB_COVERAGE_THRESHOLD=70`. Die **reine Zahl** kreuzt die Stufe.
- **Der Anstieg ist Verdünnung, nicht Ausbau** — auch das ist eine Messung: der
  **netzlos** gedeckte Anteil der gedeckten Statements steigt von **41/491
  (8,35 %)** auf **142/532 (26,69 %)**; in `replication/receive` von **11/155
  (7,10 %)** auf **112/187 (59,89 %)** (eigene `go test -coverpkg`-Läufe ohne
  `CDC_REPLICATION_TEST_DSN`, beide Enden, Block-Position-dedupliziert). Von den
  153 gedeckten `receive`-Statements braucht der überwiegende Teil **keine**
  PostgreSQL-Instanz.
- **Folge:** Ob die Stufe **steigen soll**, ist **nicht** meine Entscheidung
  (Architect/Planner). Ob der Trigger **eingetreten** ist, ist messbar und
  doppeldeutig in genau der Weise, die der Plan benennt: nach dem **Buchstaben**
  („die nächste Ausbau-Stufe schließt die Lücke zur nächsten vollen 5-%-Stufe")
  ist 76,99 % ≥ 75 — der Trigger **liest sich als fällig**; nach dem, was die
  Rampe messen **soll** („Ausbau"), trägt der Anstieg **nicht** — er stammt aus
  demselben Verdünnungs-Phänomen, das `ADR-0080` als Trigger **(a)** und
  `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` als Klasse führt.
  Eine Anhebung auf 75 würde die Verdünnung **einfrieren**. Der Vorgang hat der
  Klasse mit diesem Slice den **zweiten** Beleg zu geben (der erste war
  `slice-084`).
- **Meine Rolle:** die **Zahl** ist bestätigt (76,99 %), die **Verdünnung** ist
  bestätigt (26,69 %); die **Entscheidung** gehört in die Closure (Risiko-/Trigger-Audit)
  und an den Architect.

---

## 7. `ADR-0044` — der Digest-Träger

**Bestätigt** (mit der von `ADR-0044` gesetzten Grenze):
- Der Zug ändert Build-Kontext-Dateien (`internal/**` liegt über `.dockerignore`
  im Kontext) — die Regel „`make image` vor der Closure" greift.
- `harness/image-hash.txt` ist der **richtige** Träger; der Wechsel
  `sha256:4bd43435…` → `sha256:84bdca56…` liegt als **eigener**
  `chore(image)`-Commit `1abc7c5` vor (eine Datei, Betreff mit `ADR-0044`).
- **Eigene Gegenprobe:** ein frischer
  `docker buildx build --load --metadata-file /tmp/…` **dieses** Baums endet
  **Exit 0** und liefert `containerimage.digest =` **`sha256:84bdca56…`** —
  identisch mit dem getrackten Träger; der Baum blieb unberührt (`git status`
  nach Läufen leer, §Lauf-Artefakte).
- **Lesart-Grenze:** das ist ein Befund im **selben Builder** — **kein** Inhalts-
  und **kein** Umgebungs-Beweis (`ADR-0044` §1/§3).

---

## 8. `AGENTS.md` §3.10 — greift **nicht**

**Geprüft, nicht angenommen:** der Slice berührt **keine** Workflow-Datei.
`git diff --name-only 95480a5..HEAD -- .github/workflows/` ist **leer**, und
`git diff --name-status 95480a5..HEAD` nennt für `.github/` **keinen** Eintrag.
Damit ist §3.10 („ein neuer oder strukturell geänderter Workflow gilt erst nach
einem realen grünen Post-Push-Lauf als abgeschlossen") **nicht anwendbar** — es
gibt aus diesem Slice **kein** Post-Push-Erfordernis.

---

## 9. Entscheidungs-Konformität

- **`ADR-0080` — konform.** Festlegung 1/5: Hülle + Schnittstelle, keine
  Zusicherung — `connSession` delegiert 1:1, `var _ driverSession = connSession{}`
  steht **im Produktionspfad** (`seam.go:72`, nicht im Test), die sieben
  Operationen sind wörtlich die der ADR. Festlegung 3: der konkrete Typ bleibt in
  Hülle und Dial. Festlegung 4: der Nachweis ist als **Null-Befund** geführt,
  `k_ab = 0` (§2). Festlegung 7: `.a-check.yml` diff-frei. Kein Unterpaket, kein
  Paketwechsel — der in §„Bestätigt im Buchstaben, widerlegt in der Folge"
  benannte Nicht-Blocker tritt ein.
- **`ADR-0071` Punkt 5 — konform.** Die Naht zieht **nicht** in die prüfbare
  Fläche (Unit-Nenner 1903 → 1903), sie macht die Logik nur netzlos **testbar**;
  der konkrete Typ bleibt hinter der Naht — genau die von `ADR-0080` §„Die
  benannte Grenze" beschriebene Wirkung.
- **`ADR-0078` — konform, als Null-Fall.** `k_ab = 0`, **kein** Subjekt-Transfer,
  **keine** Neu-Bemessung; der Regressions-Riegel (`k_auf < k_ab`, gesunkene
  Quote bei unverändertem Nenner) hat kein Objekt. Der Slice stützt sich **nicht**
  auf die entwertete aggregierte Summe, sondern auf den Paket-Diff (b) und die
  grünen realen Läufe (c).
- **`ADR-0044` — konform** (§7).
- **`AGENTS.md` §3.6** (Schwellen nur per ADR): keine Schwelle geändert
  (`THRESHOLD 70`, `DB_COVERAGE_THRESHOLD 70`, Endstufen 80 — beide Enden
  identisch) — keine Schwellen-ADR fällig. **§3.5** (Accepted-ADRs immutable):
  keine ADR inhaltlich überschrieben (`ADR-0071`/`0078`/`0080` diff-frei).
  **§3.11**: d-check inkl. `hostpaths` 713 Dateien / 0 Befunde.

---

## 10. Findings (Verifier)

### V-1 — Die Form-Universale ist in beiden Sensor-Dokumenten nicht vollständig durchgesetzt (LOW)
`coverage-gate.md` §Grenze Punkt 1 trägt bewegliche Messwerte **ohne**
Ursprung/Zeitpunkt (`mapper` 20/16 · `streamv1` 86/61 = 70,9 % · `cmd` 49/0 —
heute korrekt, aber unmarkiert, §4.1); beide §Ausgabe-Abschnitte führen die
Rot-/Grün-Werte (73,38 % / 71,30 %) mit Code-Stand-, aber **ohne**
Lauf-Anker (§4.2). Die Dokumente stellen die Regel „jeder Beleg nennt seinen
Lauf" selbst auf. Klasse `zahl-in-traeger-driftet-gegen-die-messung`. **Kein**
DoD-Kriterium berührt; Nachzug in die Closure-Notiz.

### V-2 — „132 Positionen × 2" liest als Gegenwart, gilt aber nur für den Kalibrierungs-Lauf (INFO)
Heutiger Messwert **168** (§5a). Die Herkunft ist durch den Kalibrierungs-Vorbehalt
gegeben; die **Lesart** ist der Rest. **Kein** DoD-Kriterium berührt.

### V-3 — §Grenze Punkt 1 zitiert einen `go list`-Befehl, der seinen Satz nicht trägt (LOW)
`{{len .TestGoFiles}}` liefert **25**, nicht fünf; erst `TestGoFiles` **und**
`XTestGoFiles` liefert genau die fünf genannten Pakete (§5b). **Alt-Bestand**
aus `slice-079` (`65aead2`), von diesem Slice nicht berührt. **Kein**
DoD-Kriterium berührt.

**Kein HIGH, kein MEDIUM.** V-1/V-3 sind Form-/Prosa-Reste in Sensor-Dokumenten
ohne Sensor; V-2 ist Lesart-Präzision. Keiner berührt den Null-Befund, die
Verhaltens-Change-Freiheit, die Nahtform oder eine Entscheidung.

---

## 11. Eigene Messungen (Exit-Codes ungepiped, Docker-only)

| Lauf | Exit | Ausgabe / Ergebnis |
|---|---|---|
| `make gates` (Log-Datei, Exit danach separat gelesen) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · d-check 713 Dateien / 0 Befunde (inkl. `hostpaths`) · d-check `commits` 0 · `commit-traceability: OK — 5 Commit(s)` · `coverage-gate: OK — Coverage 72.00% erfüllt Schwelle 70%` · a-check 0 Befunde |
| `make test` (Race, netzlos) | **0** | **31 Pakete `ok`**, 0 `FAIL` |
| `make test-store` | **0** | Teilzahl 347/472 = 73,52 % |
| `make test-replication` | **0** | `DB-Adapter-Coverage: 76.99% (gedeckt 532 von 691 Statements; Profile gemergt: store,replication)` · `db-coverage: OK … Schwelle 70%` |
| `make doc-commits RANGE=95480a5..HEAD` | **0** | d-check `commits` 713 Dateien / 0 Befunde |
| `make doc-immutable RANGE=95480a5..HEAD` | **0** | d-check `vcs` 713 Dateien / 0 Befunde (MR-Immutabilität: kein Befund) |
| Dedup `merged.coverprofile` **`HEAD`** (eigene `awk`, Block-Position) | — | `postgresack 32/32` · `postgresstorage 347/472` · `receive 153/187` · **GESAMT 532/691 = 76,99 %** |
| Dedup **am Parent `95480a5`** (eigener Worktree, `make test-store`+`make test-replication`) | 0 | `postgresack 32/32` · `postgresstorage 347/472` · `receive 112/155` · **GESAMT 491/659 = 74,51 %** |
| Dedup `/out/coverage.out` der `coverage`-Stufe **am Parent** (eigener `docker build --target coverage`) | 0 | **1371/1903 = 72,04 %**; Per-Paket-Nenner `diff` gegen HEAD **leer** |
| Dedup `/out/coverage.out` der `coverage`-Stufe **am HEAD** | 0 | **1371/1903 = 72,04 %** (gedruckt 72.00 %) — `k_ab = 0` beidseitig |
| Unit-Profil-Prüfung (exakte ausgenommene Paketverzeichnisse) | — | **keine** Zeile aus `postgresstorage`/`postgresack`/`replication/receive`; nur die **Unterpakete** `…/postgresstorage/{mapper,sqlexec}` |
| `go test -run='^$' -coverpkg=<3>` Parent/HEAD (Nenner, ohne DB) | 0 | Parent 472/32/**155** = **659** · HEAD 472/32/**187** = **691** (+32) |
| netzlos `-coverpkg`, Parent (`--network none`, ohne DSN) | 0 | `receive` **11/155 = 7,10 %** · `postgresack` 30/32 → Summe **41** |
| netzlos `-coverpkg`, HEAD | 0 | `receive` **112/187 = 59,89 %** · `postgresack` 30/32 → Summe **142** (→ 41/491 = 8,35 % · 142/532 = 26,69 %) |
| Block-Positionen `replication.coverprofile` HEAD | — | **168** (× 2 = **336** Zeilen) |
| Block-Positionen `replication.coverprofile` Kalibrierungs-Stand `fb6adf6` | 0 | **132** (× 2 = 264) |
| `go list -f '{{len .TestGoFiles}}'` (len==0) über den Gegenstand | — | **25** Pakete |
| `go list … and … .XTestGoFiles` (beide == 0) | — | **5** Pakete (die fünf genannten) |
| §Grenze-Punkt-1-Zahlen aus `/out/coverage.out` (HEAD) | — | `mapper` **16/20** · `streamv1` **61/86 = 70,9 %** · `cmd/pg-change-feed` **0/49** |
| **V-M1/M2/M3** (Wegwerf-Kopie, netzlos) | **1** | je scharf — §3 |
| Kontrolle zur Mutation (unveränderte Kopie) | 0 | `ok …/receive` |
| `docker buildx build --load --metadata-file /tmp/…` dieses Baums | **0** | `containerimage.digest = sha256:84bdca56…` — identisch mit `harness/image-hash.txt`; Träger unberührt |
| `git diff --name-only 95480a5..610751c -- internal/ ':!…/receive/'` | — | **leer** |

Alle Mutations-/Nenner-Läufe liefen in **Wegwerf-Kopien/Worktrees** außerhalb
des Baums (`--network none` bzw. eigene PostgreSQL-Testcontainer). Die zwei
Zusatz-Worktrees (`95480a5`, `fb6adf6`) wurden **entfernt** (`git worktree list`
= nur das Hauptverzeichnis).

### Lauf-Artefakte und Hygiene
- `tools/schema/plan.yaml` wurde nach den DB-Läufen **zurückgenommen**
  (`git checkout --`); `git status` ist seitdem **leer**.
- `git rev-parse --abbrev-ref HEAD` = `main`; Historie linear (kein Merge in der
  Range); keine Trailer in den Commit-Betreffen.
- Der getrackte Digest-Träger blieb beim Rebuild unberührt.

---

## 12. Was ich **nicht** prüfen konnte (Grenzen)

- **Der Digest-Beleg bleibt builder-/lauf-gebunden** (`ADR-0044` §3): mein Rebuild
  lief im **selben** Builder. Er belegt „der eingetragene Wert stammt aus einem
  realen Bau **dieses** Baums", **nicht** Inhalts- oder Umgebungs-Gleichheit.
- **Der E2E-Tier `make test-integration`** (Compose) ist **nicht** Teil meiner
  Werkzeuge und wurde nicht gefahren; `make test-store`/`make test-replication`
  decken den DB-gestützten Belegteil dieses Slice. §3.10 greift nicht (§8).
- **Die Closure-Notiz (§7), die §6-Risiko-Ausgänge und die Register-Belege** habe
  ich **gelesen** (sie sind im Baum **offen**), aber **nicht** als Lieferung
  bewertet — die Closure kommt **nach** mir. Die vier offenen DoD-Kästchen sind
  Planner-Pflicht.
- **Die Trigger-Entscheidung der DB-Rampe** (Stufe anheben oder nicht) habe ich
  **nicht** getroffen (§6) — nur die Zahl und die Verdünnung gemessen.
- **Die Form-Frage „welche Zahlen sind beweglich"** ist ein Urteil über Prosa;
  ich habe sie durch **Nachmessen** der Kandidaten und gegen die universale
  Fassung geprüft, nicht durch einen Sensor (es gibt keinen).

---

## 13. Verdikt

**Der Slice ist closure-fähig.**

- Jedes **Liefer**-Kriterium der DoD ist im heutigen Text erfüllt; die vier
  offenen Kästchen sind **Closure-Pflichten** des Planners (Notiz, Register,
  Risiko-Ausgänge; das Review-Häkchen) — kein Liefer-Kriterium hängt an ihnen.
- Der **Null-Befund (a)/(b)/(c)** ist in allen drei Teilen **unabhängig
  nachgemessen**, an **beiden** Enden der Range: `k_ab = 0` (Unit-Nenner 1903
  beidseitig, ausgenommene Pakete absent), DB **659 → 691 (+32)** mit
  `receive 155 → 187`, `postgresstorage` 472 / `postgresack` 32 byte-stabil; der
  Paket-Diff zeigt **keinen** Trägerwechsel (Range **und** Pathspec, beide
  Richtungen), die Gegenstands-Listen sind diff-frei.
- **Kein Verhaltens-Change** — die realen Tiere grün, `stream_test.go`
  byte-identisch, kein Testfall entfernt, öffentlicher Rand unverändert; durch
  **drei eigene** Mutationen an bisher unberührten Zusagen scharfgestellt.
- Die **Form** trägt in ihrer tragenden Hälfte, ist aber **nicht** vollständig
  durchgesetzt (**V-1**) — die Zahlen selbst stimmen.
- `ADR-0080`, `ADR-0071` Punkt 5, `ADR-0078` (Null-Fall) und `ADR-0044` sind
  eingehalten; §3.10 greift **nicht** (kein Workflow berührt).

**Mitgabe an die Closure (nicht blockierend):**
1. **Der DB-Hochschalt-Trigger** liest sich nach dem Buchstaben als fällig
   (76,99 % ≥ 75), ist aber **verdünnungs**-getrieben (netzlos gedeckter Anteil
   8,35 % → 26,69 %); die Entscheidung „Stufe heben" gehört in den
   Trigger-Audit/Architect — eine Anhebung würde die Verdünnung einfrieren.
2. **Register:** `evidence/slice-085.md` in
   `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (**2×**) **und** in
   `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**3×** — Lese-Schritt
   fällig) nachtragen.
3. **V-1/V-2/V-3** als Prosa-Nachzüge (Form-Universale, „132"-Lesart,
   `go list`-Formel).
4. Die §6-Risiken mit **einem** Ausgang je Risiko (1 ist gestrichen; 2–4 offen)
   und das Review-Häkchen nachziehen.

---

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, diese Messungen, dieses
Modell). Er ersetzt **keine** Validierung (Modul 11, repo-extern) und **keine**
Review — der Diff ist gegen Plan und Entscheidungen bereits geprüft
(`review-slice-085.md`).
