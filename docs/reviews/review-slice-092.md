# Review-Report: slice-092 — Coverage Cluster D1 (Anwendungs-Kern) · 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier, Modul 11),
§6-Risiko-Ausgänge, Beobachtungs-Register und die drei Paarungen
(Planner-Closure) sind **nicht** Gegenstand dieses Reports.

**Gegenstand:** `3de9547` (Parent `2b7c6ec`) — **ein** Commit, **15** Dateien,
874 hinzugefügte Zeilen: **13** `*_test.go` und **zwei** Planungs-Dokumente
(`slice-092-coverage-cluster-d1.md` §1/§3/§5/§6/§8, `welle-20.md` §4).
Gemessen: `git diff --name-only 3de9547^..3de9547` listet 15 Pfade, davon
**2** ohne die Endung `_test.go`. Produktcode (`internal/**`, `cmd/**`
außerhalb von Tests), `THRESHOLD`, ADRs und `harness/sensors/coverage-gate.md`
sind **unberührt**; Arbeitsbaum sauber (`git status --porcelain` leer).

**Skill:** `.harness/skills/reviewer.md` @ `3de9547` — die vier
repo-spezifischen HIGH-Unterpunkte (u. a. „Zusage ohne Bindung an ihre
Eingabeseite“, „Zahl im Träger“) gehören zum Prüfraster.
**Modell:** `deepseek-v4.1-flash:cloud[1m]` · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-092-coverage-cluster-d1` vollständig (§1–§8), die offene
  Welle `welle-20`
- `ADR-0082` (Schnittmaß, Cluster-Tabelle, §Konsequenzen „test-only“,
  §Re-Evaluierungs-Trigger), `ADR-0071` (Messgegenstand), `ADR-0077` (Rampe),
  `ADR-0024` (`LogPort`), `ADR-0055`/`ADR-0060` (die zwei Best-Effort-Pfade
  des Capture), `ADR-0028`/`ADR-0029` (die Use-Case- und Domänen-Zusagen),
  `ADR-0054` §(a)
- Berührte `LH-*`: `LH-FA-ADM-*`, `LH-FA-CFG-001`, `LH-FA-CFG-002`,
  `LH-FA-CFG-003`, `LH-FA-CFG-005`, `LH-FA-CON-001`, `LH-FA-CON-003`,
  `LH-FA-CON-004`, `LH-FA-CON-006`, `LH-FA-CAP-005`, `LH-FA-CAP-006`,
  `LH-FA-RET-002`, `LH-QA-REL-001.a` · `SPEC-001`
- `AGENTS.md` §3.1, §3.5, §3.6, §3.7, §3.9, §3.11, §3.12, §4, §6 ·
  `harness/conventions.md` (MR-000 ID-Schema)
- Beobachtungs-Register `BEO-PGC/:` `negativtest-ohne-bindung-an-seine-eingabe`
  (**5×**, Belege `slice-083/086/087/088/091`), `beleg-befehl-traegt-seinen-satz-nicht`
  (**3×**), `arbeit-ueberholt-stehenden-traeger` (**1×**),
  `zahl-in-traeger-driftet-gegen-die-messung` (**6×**) — mit `ls evidence/ | wc -l`
  nachgezählt, nicht aus Prosa übernommen
- Vorläufer am gleichen Modul: `review-slice-088` (F-2/F-3 Plan-Zahlen),
  `review-slice-089` (F-1/F-3 Zahl im Träger), `review-slice-091`
  (+ Delta) — keine Klasse doppelt zu diesem Diff

---

## Findings

### F-1 — Die Korrektur der Fehlzählung ist unvollständig und setzt für dieselbe Aussage eine dritte Zahl

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz A · Klasse
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (6×; Präzedenz für
  Plan-Zahlen: `review-slice-088` F-2/F-3, `review-slice-089` F-3 — je MEDIUM)
- `pfad`: `docs/plan/planning/in-progress/slice-092-coverage-cluster-d1.md:100`
  (§2 LP1) gegen `:49-50` (§1), `:141` (§3), `:187` (§5), `:211` (§6), `:291` (§8)
- `befund`: Der Commit erklärt „zehn Pakete“ zur Fehlzählung des Globs und
  berichtigt sie an fünf Stellen; **§2 LP1 trägt sie weiter** („Für die zehn
  Use-Case-Pakete und `domain/model` liegen Tests vor“) — gemessen führt der
  Glob **13** Pakete, und §1 nennt seit demselben Commit beide Zahlen. Im
  **selben Zug** wurde in §5/§6 „zehn“ durch „**dreizehn**“ ersetzt, obwohl
  beide Sätze das **Bewegen der Deckung** behaupten („bewegt eine gemessene
  Eigenschaft von dreizehn Paketen“ / „er bewegt die Deckung von dreizehn
  Paketen“): gemessen bewegt dieser Diff die Deckung von **11** Paketen
  (**10** Use-Case-Pakete + `domain/model`; `excludecolumn`, `includecolumn`,
  `readchanges` standen am Parent bereits vollständig gedeckt und bewegen
  sich nicht) — und §8, unverändert, behauptet für **dieselbe** Aussage
  „Deckung von **elf** Paketen“, also den gemessenen Wert. Dasselbe Dokument
  führt damit für eine Aussage zwei Zahlen (13 in §5/§6, 11 in §8) und für
  ihren Gegenstand vier (10 · 11 · 12 · 13) ohne Unterscheidung. Die
  Commit-Message nennt als berichtigte Stellen „§1/§3/§6/§7“ — geändert hat
  der Diff §1/§3/**§5**/§6/**§8** (§7 ist unberührt); die Liste trägt ihren
  Satz ebenfalls nicht (git-Historie, außerhalb des §3.12-Geltungsbereichs,
  hier nur als Beleg der Unvollständigkeit).
- `verifizierbar`: ja — `go list ./internal/application/usecase/... | wc -l`
  → **13**; Profil am Parent paketweise: 10 Use-Case-Pakete mit
  `uncovered > 0` + `domain/model` → **11** bewegte Pakete (Tabelle unten);
  `grep -n "die zehn\|dreizehn\|elf Paketen" <plan>` → die vier Stellen
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`

### F-2 — Das Acht-Läufe-Band trägt seinen Lauf, aber sein Qualifier zeigt auf den heutigen Gegenstand — und der ist 24 Statements weiter

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A · Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (1×, Erstauftreten `slice-091`)
  · Register `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `harness/sensors/coverage-gate.md:68-74` (Träger **außerhalb** des
  Diffs; §3-Zeile des Plans: „update, **nur falls** eine Zahl dort gegen die
  Messung driftet“)
- `befund`: Der Satz berichtet „Über **acht** Läufe **desselben, hier
  gegenständlichen** Produktionsstands … lag die gedeckte Zahl zwischen
  **1468** und **1471**, die gedruckte Zeile zwischen `77.1%` und `77.3%`
  (Lauf `slice-091`)“. Der Diff ist **test-only** — der Produktionsstand ist
  identisch geblieben. Gemessen an genau diesem Stand liegen sechs eigene
  Läufe jetzt bei **1493–1495** gedeckt (`78.5%`/`78.6%`): das Band ist um
  **24** Statements gewandert, und zwar nicht durch Code, sondern durch die
  Tests (die gedeckte Zahl hängt am *Baum*, nicht am „Produktionsstand“).
  Der Satz ist als **Lauf-Beleg nicht falsifiziert** (er nennt seinen Lauf;
  §Zählbasis erklärt die gedeckte Zahl ausdrücklich zu „kein Zustand“), und
  die von ihm abgeleitete **Breite** `3 Statements = 0,16 pp` gilt weiter
  (eigene Messung: Band von 3). Was nicht trägt, ist die Deixis: „derselbe,
  hier gegenständliche Produktionsstand“ lud schon bei seiner Einführung
  (`a7d7f7b`, D-2) dazu ein, das Band als das des *jeweiligen* Gegenstands zu
  lesen — und dieser Gegenstand zeigt heute ein anderes Band. Die
  §3-Bedingung des Plans ist damit im engen Sinn **nicht** erfüllt (kein
  Lauf-Wert ist falsch; ein Nachführen wäre eine bewegliche Zahl als
  Ist-Stand und damit gerade die verbotene Form) und im weiten Sinn
  **berührt** (die Aussage über den heutigen Gegenstand ist contra-gemessen).
  Der Rest ist eine Auffrischungs-Option für die Closure — dieselbe
  Behandlung, die `review-slice-091` der `streamv1`-Zeile derselben Datei gab.
- `verifizierbar`: ja — Profil der `coverage`-Stufe bei `3de9547` sechs Mal
  (`go test -count=1 -coverpkg=… -covermode=atomic`, Auswertung über die
  Block-Position): gedeckt ∈ {1493, 1495}, `go tool cover -func` druckt
  4 × `78.5%`, 2 × `78.6%`; `make gates` druckt `78.60%`
- `klasse`: `BEO-PGC/arbeit-ueberholt-stehenden-traeger`

### F-3 — Die zwei vorbestehenden Tests, die dieser Vorgang als Klassen-Fälle fand und band, sind nach der Register-Regel ein sechster Beleg-Vorgang

- `kategorie`: INFO
- `quelle`: `docs/plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md`
  („Zähler (abgeleitet): **5×**“) · die Regel, die
  `evidence/slice-091.md` selbst anwendet („**Gefunden** wurde er in
  `verify-slice-091`, und das ist das Vorkommen, das dieser Beleg dateit“)
- `pfad`: `internal/application/usecase/excludecolumn/service_test.go:71`
  und `internal/application/usecase/includecolumn/service_test.go:64` (Parent-Fassung)
- `befund`: Der Diff meldet im Commit-Text, er habe „zwei **VORBESTEHENDE**
  Tests gefunden, die genau die Klasse waren“, und bindet sie an die
  Spaltenadresse. Nachgemessen an der Parent-Fassung: der Fake lieferte
  `f.exists` (mit `exists: false`) unabhängig von der Abfrage — die
  Ablehnung war an keinen Eingabewert gebunden, die Form trifft also zu.
  Nach der Regel, die der Register-Eintrag für `slice-091` anwendet, zählt
  ein **im Vorgang gefundener** Fall als Vorkommen (Ursprung ≠ Vorkommen) —
  damit wäre `slice-092` der **sechste** Beleg-Vorgang dieser Klasse. Der
  Plan nennt in §8 den Stand **bei der Planung** (5×), das ist zu Recht so;
  die Dateiung des Belegs gehört der Slice-Closure. Kein Diff-Mangel, aber
  die Zahl bewegt sich durch diesen Diff, und der Eintrag steht bereits auf
  der Schwelle (Lese-Schritt: `welle-20`-Closure).
- `verifizierbar`: ja — Parent-Fassung der beiden Testdateien
  (`git show 3de9547^:internal/application/usecase/excludecolumn/service_test.go`)
  gegen die Mutationsprobe M-21 (unten); `ls …/evidence/ | wc -l` → 5
- `klasse`: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`

---

## Negativbefunde

- geprüft, ohne Befund: **die Zahlen des Commits — selbst nachgemessen.**
  Eigener Lauf über die Paketliste des Gates (`go list ./internal/... ./cmd/...`
  ohne `postgresstorage|postgresack|replication/receive`), Profil
  `-coverpkg=<alle> -covermode=atomic`, dedupliziert über die Block-Position
  (dieselbe Zählbasis wie `tools/harness/db-coverage.sh`), im netzlosen
  Toolchain-Container:

  | Stand | ungedeckt | Nenner | `go tool cover -func` total |
  |---|---|---|---|
  | `2b7c6ec` (Parent) | **432** | **1903** | **77.3 %** |
  | `3de9547` (Diff) | **410** | **1903** | **78.5 %** |

  Die Differenz der ungedeckten Block-Positionen ist **exakt** die des
  Commits: **24** neu gedeckte Blöcke = **24** Statements (**22** in
  `usecase/*/service.go`, **2** in `domain/model` — `source.go:31.2,31.40`
  und `transaction.go:29.3,30.1`), dagegen **1** Block = **2** Statements
  (`internal/bootstrap/wiring.go:991.5,992.13`, der Takt-Zweig von
  `runWALRetentionCheck`) neu ungedeckt — Netz **+22**, genau die erklärte
  Schwankung. `Nenner unveraendert 1903` hält, `THRESHOLD` ist unverändert
  **70**, `make gates` druckt in meinem Lauf `coverage-gate: OK — Coverage
  78.60% erfüllt Schwelle 70%`.
- geprüft, ohne Befund: **Cluster D1 steht real bei `uncovered = 0`.**
  Paketweise am Diff-Stand (dedupliziert): alle **13** Use-Case-Pakete
  `voll`, `domain/model` **110/110** `voll`. Die im Plan genannten Größen
  halten einzeln: **13** Use-Case-Pakete, **10** davon am Parent mit
  `uncovered > 0`, **12** der 13 Pakete mit geänderter Testdatei
  (`readchanges` unberührt, weil dort kein ungedecktes Statement lag),
  `welle-20.md:100` `(13 Pakete, 10 davon mit ungedeckten Statements) | 24`.
- geprüft, ohne Befund: **die Prüf-Kraft der neuen Zusagen — 32 eigene
  Mutationsproben, 32 × rot, kein Theater.** Alle **24** neu gedeckten
  Block-Positionen sind über eine Zusage gebunden, die an ihrer
  **Eingabeseite** rot werden kann; die Proben sind in der Tabelle unten
  einzeln ausgewiesen. Damit ist §6-Risiko 3 („Coverage-Theater“) für diesen
  Diff **nicht** eingetreten, und §6-Risiko 2 („Negativtest bindet an den
  Fake“) **nicht** an einer neuen Stelle: die Proben M-06/M-09/M-24/M-26
  mutieren die **Weitergabe der Kommando-Werte** (Adresse, Quelle,
  Publication) und färben die neuen Tests rot — wer nur den Fake mutiert,
  sieht sie grün, die Kette ist aber an beiden Enden gebunden.
- geprüft, ohne Befund: **die Klassen-Führung dieses Diff.** Nach der
  Mutations-Tabelle fügt der Diff **kein** neues Vorkommen von
  `negativtest-ohne-bindung-an-seine-eingabe` hinzu (er repariert zwei
  vorbestehende, siehe F-3); die zwei neuen Testkommentare, die eine
  Ablehnung zusichern, tragen die Gegenprobe jeweils im Testkörper
  (zweiter Aufruf mit getragenem Wert) — und die Gegenprobe ist
  mitgemessen, nicht behauptet.
- geprüft, ohne Befund: **keine Wirkungs-Behauptung ohne Deckung im Diff.**
  `grep` über alle `+`-Zeilen nach `mutier|färb|rot|grün|bliebe|würde|wäre`
  liefert **drei** Treffer, alle drei aus **vorbestehenden**, nicht neu
  geschriebenen Kommentarhälften („die er zuvor aus dem Kommando gebaut
  hat“, „die zuvor geschriebene Bindungs-Zeile“, „ohne dieses Sentinel
  bliebe…“). Der Diff nennt **keine** Mutation und behauptet keine
  Rot-Wirkung; die Belege stehen im Bericht — die Aussage des Commits hält.
  Auch die neu geschriebenen Kommentare tragen deckbare Sätze: die zwei
  Bindungs-Behauptungen mit „entscheidet die Policy des Kommandos“ (M-25b)
  und „über welchen Port … melden“ (M-02) sind gemessen.
- geprüft, ohne Befund: **`BEO-PGC/arbeit-ueberholt-stehenden-traeger` —
  die Prüfung des Commits hält.** `go list -f '{{.ImportPath}} Test={{len
  .TestGoFiles}} XTest={{len .XTestGoFiles}}'` über den Gegenstand liefert
  am Parent **und** am Diff **dieselbe** Datei (md5
  `adda7355479fc95c8414a9e9ff3947bf`) — „byte-identisch“ stimmt. Auch die
  Einzelzahlen von `harness/sensors/coverage-gate.md` §Grenze Punkt 1
  reproduzieren am Diff-Stand: **31** Pakete im Gegenstand, **genau vier**
  mit `Test=0 XTest=0` (die vier genannten), **23** mit `Test=0`, **19**
  davon mit `XTest>0`; `cmd/pg-change-feed` **0/49**; `streamv1` **76/86**;
  `postgresstorage/mapper` **20/20**; die drei Ausgenommenen **31/472**,
  **112/187**, **30/32** (die Rückrechnung in §Grenze Punkt 4 behält mit
  dem heutigen Zähler alle drei Ausgänge: 78,71 % grün · 76,79 % grün ·
  64,17 % rot).
- geprüft, ohne Befund: **die §3-Zeile `welle-20.md`.** Der Diff hat die
  Soll-Tabelle berichtigt (`13 Pakete, 10 davon mit ungedeckten Statements`),
  aber **keine** erreichte Zahl dort als Ist-Stand eingetragen. Auch sonst
  führt **kein** stehender Träger die erreichten Zahlen: `grep` über die
  `*.md` außerhalb von `docs/reviews/` und der vendored Baseline nach
  `78,5|78.5|78,6|78.6|1493|1495` liefert keinen Treffer in einem
  Doku-Träger. Die Zahlen leben bisher nur in der Commit-Message
  (git-Historie, außerhalb des §3.12-Geltungsbereichs) — §7 muss sie mit
  ihrem Lauf nennen (LP1 verlangt das ohnehin).
- geprüft, ohne Befund: **§3.7 und §3.11.** Keine Slice-/Wellen-Nummer in
  einem Kommentar (`grep -nE 'slice-[0-9]|welle-[0-9]'` über die
  `+`-Zeilen: leer), kein Vorher/Nachher-Vokabular, keine `//nolint`-Form;
  keine `+`-Zeile nennt einen host-lokalen absoluten Pfad (Präfix-Grep
  leer). `make docs-check` ist nach den Änderungen grün: „757 Datei(en)
  geprüft, 0 Befund(e)“.
- geprüft, ohne Befund: **`ADR-0082` Folgepflicht „test-only“.** 13 der 15
  Pfade enden auf `_test.go`; die zwei übrigen sind Planungs-Dokumente
  (`docs/plan/planning/**`), kein Produktcode. Kein `git mv`, keine Naht,
  kein Gate, keine Schwelle, keine Accepted-ADR berührt (§3.3, §3.5, §3.6
  unberührt); `harness/sensors/coverage-gate.md` wurde — wie die §3-Zeile es
  vorsieht — **nicht** angefasst (F-2, INFO).
- geprüft, ohne Befund: **Netzlosigkeit und Nebenläufigkeit.** `make test`
  (voller Unit-Lauf mit Race-Detector, `--network none`) → **EC=0**, 32 ×
  `ok`, 0 × `FAIL`; die sechs Profil-Läufe und die 32 Mutationsproben liefen
  ebenfalls netzlos (EC=0). Kein `t.Parallel` in den neuen Tests, keine
  paketweit veränderliche Doppel; die `…ErrFor`-Tabellen werden je Test
  instanziiert und nach der Konstruktion nur gelesen.
- geprüft, ohne Befund: **Traceability.** Betreff trägt `ADR-0082`, keine
  Struktur-ID (`SPEC-*`/`ARC-*`); `make commit-traceability` grün über die
  letzten fünf Commits. Keine neue Kennung, kein neues Präfix (MR-000
  unberührt). Keine neue Betreiber-Oberfläche (`CDC_*`, `cdc.*`, Endpunkt) —
  die beiden HIGH-Klassen Handbuch-Zug und Versionshistorie haben kein
  Objekt.
- geprüft, ohne Befund: **die vier Zusagen-Familien der neuen Tests, je
  einzeln.** `acknowledge`/`position`/`register`/`remove` binden den
  Port-Fehler an die Consumer-Kennung des Kommandos (M-01, M-12, M-13, M-14
  + die jeweilige Gegenprobe im Testkörper); `disable`/`enable`/`status`
  binden Existenz-, Bindungs- und Publication-Fehler an Adresse · Quelle ·
  Publication (M-03…M-09, M-18, M-19, M-24, M-28, M-30…M-32); `list` an
  Quelle und an Publication-plus-Adresse der gelesenen Zeile (M-10, M-11,
  M-26); `retention` an Quelle und an die **freigegebene Menge** (M-15,
  M-16, M-17, M-25b); `capture` an den injizierten `LogPort` samt Nachricht
  und Attribut (M-02, M-29).

---

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-only 3de9547^..3de9547` | **0** | 15 Pfade, **2** ohne `_test.go` (die Planungs-Dokumente) |
| Profil bei `2b7c6ec` (Container, netzlos) | **0** | **432** ungedeckt / **1903**; Cluster D1 = **24** (22 + 2) |
| Profil bei `3de9547` | **0** | **410** ungedeckt / **1903**, gedruckt `78.5 %`; D1 = **0** |
| Block-Differenz Parent → Diff | **0** | **+24** Statements gedeckt; **−2** (`wiring.go:991.5,992.13`) |
| 6 × Profil bei `3de9547` | **0** | gedeckt ∈ {**1493**, **1495**}, ungedeckt ∈ {410, 408}; gedruckt 4 × `78.5%`, 2 × `78.6%` |
| paketweise bei `3de9547` | **0** | 13 Use-Case-Pakete + `domain/model` `voll`; `mapper` 20/20; `streamv1` 76/86; `cmd` 0/49 |
| drei Ausgenommene, `-coverpkg` über alle | **0** | `postgresstorage` **31/472**, `receive` **112/187**, `postgresack` **30/32** |
| `go list -f '{{…Test…}}'` Parent vs. Diff | **0** | md5 **identisch** (`adda7355…`); 31 Pakete, 4 × `Test=0 XTest=0`, 23 × `Test=0`, 19 × `XTest>0` |
| `go list ./internal/application/usecase/...` | **0** | **13** Pakete |
| **32 Mutationsproben** (Tabelle unten) | **1** je Probe | **32 × zieltest rot** (1 Probe verworfen, s. u.) |
| `make gates` (Log, EC separat gelesen) | **0** | 6 Checks grün: `coverage-gate: OK — Coverage 78.60% erfüllt Schwelle 70%`, `baseline-verify` 54 Dateien OK, `d-check` 757 Dateien 0 Befunde, `commit-traceability` OK, `generated-sync` OK, `a-check` 0 Befunde; danach `git status --porcelain` leer |
| `make test` (voller Lauf, `-race`, `--network none`) | **0** | 32 × `ok`, 0 × `FAIL` |
| `grep` `+`-Zeilen: `slice-`/`welle-` | **1** (leer) | kein Treffer |
| `grep` `+`-Zeilen: host-lokale Pfade | **1** (leer) | kein Treffer |
| `grep` `*.md` (ohne `docs/reviews/`) nach den erreichten Zahlen | **1** (leer) | keine Zahl in einem stehenden Träger |

**Mutationsproben** (Produktions-Mutation → benannter Test; `-count=1`,
netzlos, jede Probe mit `diff -q` gegen die Variantendatei als
„angewandt“-Beleg):

| Probe | Produktionsmutation | Ergebnis |
|---|---|---|
| M-01 | `acknowledge/service.go:62` Fehler verworfen | **rot** |
| M-02 | `capture/service.go:77` `WithLog` setzt `NoopLog` | **rot** |
| M-03/M-04/M-05 | `disable/service.go:61`, `:72`, `:80` (Unpublish/Unregister/`default`) | **3 × rot** |
| M-06 | `disable/service.go:53` `TableExists` mit fester Adresse | **rot** |
| M-07/M-08 | `enable/service.go:68`, `:71` (Register/Publish) | **2 × rot** |
| M-09 | `enable/service.go:59` `TableExists` mit fester Adresse | **rot** |
| M-10/M-11 | `list/service.go:50`, `:56` (List/Published) | **2 × rot** |
| M-12/M-13/M-14 | `position:48`, `register:50`, `remove:50` | **3 × rot** |
| M-15/M-16/M-17 | `retention/service.go:60`, `:65`, `:78` | **3 × rot** |
| M-18 | `status/service.go:69` Published-Fehler verworfen | **rot** |
| M-19/M-28 | `status/service.go:60` feste Quelle, `:53` feste Adresse | **2 × rot** |
| M-20/M-21 | `exclude-/includecolumn/service.go:40` feste Spalte | **2 × rot** |
| M-22/M-23 | `model/source.go:31` konstanter Name, `transaction.go:28` Bedingung `&&` | **2 × rot** |
| M-24 | `disable/service.go:63` `Registered` mit fester Quelle | **rot** |
| M-25b | `retention/service.go:72` Freigabe unabhängig von der Policy | **rot** |
| M-26 | `list/service.go:54` `Published` mit fester Adresse | **rot** |
| M-27 | `enable/service.go:51` `NewSourceTable` mit fester Kennung | **rot** |
| M-29 | `capture/service.go:132` Warn-Nachricht vertauscht | **rot** |
| M-30/M-31/M-32 | `disable:55`, `enable:61`, `status:55` Fehler verworfen | **3 × rot** |

Verworfen und **nicht mitgezählt**: M-25 (`retention/service.go:72` →
`if true`) — die Probe bricht den Bau (`declared and not used`), ihr Rot ist
ein Compiler-Fehler, keine gefangene Mutation; ersetzt durch M-25b in einer
bauenden Form.

Der Gate-Lauf und seine Auswertung sind **zwei Schritte**: `make gates` lief
ungepiped in eine Log-Datei, der Exit-Code separat gelesen und erst danach in
einem eigenen Werkzeug-Aufruf ausgewertet; kein Folgekommando hing an einer
Pipe oder an einem Wrapper (§3.9).

---

## Antwort auf die Schwerpunkte

**(1) Prüf-Kraft statt Zeilen — trägt, ohne Fund.** 32 eigene
Mutationsproben, 32 × rot; jede der **24** neu gedeckten Block-Positionen ist
über eine Zusage gebunden, die an ihrer Eingabeseite färbt. Kein Test in
diesem Diff prüft nur „dass Code läuft“: die Zusagen sind Fehler-Klassen an
Use-Case-Grenzen, und jede Probe mit einer Produktionsmutation färbt ihren
benannten Test rot. §6-Risiko 3 ist damit beantwortet.

**(2) `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (5×) — erfüllt,
und die Klasse tritt nicht neu auf.** Die neuen Tests binden durchweg: die
Proben, die nur die **Eingabeseite** mutieren (M-06, M-09, M-24, M-26, M-27:
feste Adresse/Quelle/Publication/Kennung in der Weitergabe), färben die
neuen Tests rot — der häufigste blinde Fleck dieses Musters ist damit
geschlossen. Der Diff **repariert** zwei vorbestehende Fälle (F-3, INFO) und
fügt keinen neuen hinzu; der Zähler bleibt bei 5×, *bis* die Slice-Closure
den Reparatur-Vorgang dateit (nach der Regel aus `evidence/slice-091.md`
wäre das der sechste Beleg-Vorgang).

**(3) `beleg-befehl-traegt-seinen-satz-nicht` (3×) — kein Vorkommen im
Diff.** Die Testkommentare nennen keine Mutation und behaupten keine
Rot-Wirkung (`grep` über die `+`-Zeilen: kein Treffer für
`mutier|färb|rot|grün`); die drei Treffer für „zuvor“/„bliebe“ liegen in
vorbestehenden Kommentarhälften. Die zwei Behauptungen, die der Diff
**selbst** aufstellt, sind gemessen gedeckt: „über welchen Port die
Best-Effort-Pfade melden“ (M-02, M-29) und „ob sie das tut, entscheidet die
Policy des Kommandos“ (M-25b). Kein Rückgabe-Pfeil aus dieser Klasse.

**(4) DER GRENZFALL — der Satz trägt als Lauf-Beleg; die §3-Bedingung ist im
engen Sinn nicht erfüllt (F-2).** Gemessen: derselbe Produktionsstand, sechs
Läufe, gedeckt 1493–1495 (`78.5%`/`78.6%`) gegen das im Träger genannte Band
1468–1471 (`77.1%`–`77.3%`) — **24 Statements** Verschiebung, bei gleicher
Breite (**3**). Der Satz nennt seinen Lauf (`slice-091`) und ist als
Lauf-Beleg nicht falsifiziert; die von ihm weiter unten benutzte **Breite**
`3/1903 = 0,16 pp` gilt auch am heutigen Stand. Was altert, ist die Deixis
(„**denselben, hier gegenständlichen** Produktionsstands“): sie bindet das
Band an den *jeweiligen* Gegenstand, und der ist heute ein anderer Baum.
Nach der §3-Zeile („nur falls eine Zahl dort gegen die Messung driftet“) ist
die Datei deshalb **zu Recht** nicht angefasst — ein Nachführen des Bandes
wäre eine bewegliche Zahl als Ist-Stand und damit die Form, die §Zählbasis
selbst verbietet. Der Rest ist die Entscheidung der Closure: entweder das
Band ausdrücklich als Stand `slice-091` (Produktion **und** Tests) kennzeichnen
oder es als Auffrischungs-Option behandeln — dieselbe Behandlung wie die
`streamv1`-Zeile in `review-slice-091`. Ausdrücklich **verworfen**: den Satz
nachzuführen und dabei die Deixis zu behalten.

**(5) Kein Ist-Stand in einem stehenden Träger — geprüft, ohne Befund.**
`welle-20.md:100` trägt die **Soll**-Zahl (`24`) und die berichtigte
Paket-Angabe; die **erreichten** Zahlen (1493/1495, 408/410, `78.5 %`/`78.6 %`)
stehen in **keinem** `*.md` außerhalb von `docs/reviews/` und der vendored
Baseline. Damit ist die §3-Zeile „`welle-20.md` §4 = **nicht**“ eingehalten,
und die Ablage-Pflicht liegt unverändert bei der Closure-Notiz §7 (mit Lauf).

**(6) §3.12 — welche Zahl gehört wohin.** Gemessen und zuordenbar:
`22`/`2`/`24` = Planungsmessung (Parent-Profil, dedupliziert) — gehört mit
diesem Stand in §1/§7; `+24` Statements = Differenz zweier Läufe dieses
Vorgangs — gehört als **abgeleitet** in §7; `1493`/`1495` und `408`/`410`
= **gemessen**, Band von sechs Läufen — gehört mit Lauf in §7;
`78.5 %`/`78.6 %` (gedruckt) — dito, und `make gates` druckte in meinem Lauf
`78.60%`; `13`/`10`/`12` = gemessen (Glob bzw. Parent-Profil bzw. Datei-Liste)
— **gehören mit ihrem Zeitpunkt in §1/§2/§3**, und genau dort ist F-1
angesiedelt (die berichtigten Zahlen tragen die Herkunft ihrer Nachbarsätze,
aber die *Zählung* des Bewegten ist an drei Stellen widersprüchlich);
`21 Mutationen` = Angabe des Berichts, nicht des Diffs — nach meiner eigenen
Probe plausibel, aber nicht selbst nachgefahren (kein Diff-Träger, keine
Behauptung im Code). Die Zahl, die in §7 **nicht** fehlen darf, ist die
erreichte Quote mit ihrem Lauf (LP1 verlangt sie).

**Was nicht geprüft wurde:** DoD/LP1–LP3 (Verifier), die §6-Risiko-Ausgänge
und die drei Paarungen (Planner-Closure), die Dateiung des sechsten
Beleg-Vorgangs (F-3, Planner), und ob die 874 Zeilen neuer Tests die
CI-Laufzeit im Rahmen halten (dafür fehlt der reale Post-Push-Lauf;
`AGENTS.md` §3.10 gilt für den betroffenen Workflow, den dieser Diff
**nicht** ändert).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **0** |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(F-1: unvollständige Korrektur einer Fehlzählung + drei Zahlen für eine
Aussage — **derselbe Vorgang**; F-2 **ohne** neuen Beleg, der Träger liegt
außerhalb des Diffs) · `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (F-2,
zweiter Vorgang) · `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(F-3: sechster Beleg-Vorgang — der Vorgang **repariert**, er fügt nicht hinzu)

## Verdikt

**Merge-blockierend:** ja — **1 MEDIUM** (F-1). HIGH ist **keiner** dabei;
F-1 ist eine Text-Korrektur in einem Planungs-Dokument dieses Diffs, keine
Verletzung einer Hard Rule. Der Diff selbst — 13 Testdateien, kein
Produktcode — trägt: die Zahlen des Commits halten, Cluster D1 steht real bei
`uncovered = 0`, und 32 eigene Mutationsproben färben alle rot.

**Rückgabe-Pfeil Reviewer → Implementer:** **nötig** (F-1, eine Textstelle in
§2 LP1 plus die Zählung des Bewegten in §5/§6/§8 des Slice-Plans) — der
Verifier prüft die DoD **gegen diesen Plan**, die vier Stellen sollten
vorher stimmen. F-2 und F-3 liefern **keinen** Rückgabe-Pfeil: F-2s Träger
liegt außerhalb des Diffs und die §3-Bedingung ist zu Recht nicht gezogen,
F-3 gehört der Closure. Die Finding-Klassen gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler.

**DoD-Häkchen „Review durchgeführt, Report unter `docs/reviews/` liegt vor“:**
bleibt **offen** — der Slice braucht eine Fixrunde
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde greift
nicht); die Fixrunde deckt nach §2 des Plans ein **Delta-Review** ab.

**Übergabe:** Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-
Konformität prüft der Verifier separat (Modul 11).
