# Review-Report: slice-090 — Sync-Gate des Protobuf-Codes · 2026-09-16

**Review-Art:** Code + Tooling — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier, Modul 11), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `9029c05` (= `HEAD`, ein Commit, Parent `cf3f887`) — fünf Dateien,
`git diff --name-status 9029c05^..9029c05`: `A harness/mk/generated-sync.mk` (22 Z.),
`A tools/harness/generated-sync.sh` (120 Z.), `M harness/README.md` (+2/−1),
`M AGENTS.md` (+1), `M docs/plan/planning/in-progress/slice-090-sync-gate-protobuf.md`
(+2/−2). Kein Produktionscode, `Makefile` und `Dockerfile` unberührt, Baum sauber.

**Skill:** `.harness/skills/reviewer.md` @ `9029c05` — die fünf repo-spezifischen
HIGH-Unterpunkte (u. a. Zahl im Träger, Zusage ohne Bindung an ihre Eingabeseite)
gehören zum Prüfraster · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-090` vollständig (§1–§8), die offene Welle `welle-20`
- `ADR-0084` (Festlegung 1 = Baum nicht schreiben; Festlegung 2 = Befund nennt
  Datei **und** Zeile; §Was das Gate nicht fangen kann; §Re-Evaluierungs-Trigger),
  `ADR-0060` (der Generator und seine Folgepflicht), `ADR-0044` (Lauf-Beleg),
  `ADR-0045`/`ADR-0041` (Haus-Muster der `GATE_CHECKS`-Bindung)
- `AGENTS.md` §3.1, §3.6, §3.7, §3.9, §3.11, §3.12, §4, §6 ·
  `harness/README.md` §Sensors **samt** dem Kommentar-Block über der Tabelle
  (Z. 60–109) · `harness/conventions.md` (MR-000 ID-Schema)
- Letzte fünf Reviews am gleichen Modul: `review-slice-085`/`-086`/`-087`/`-088`/`-089`
  (Klassen: Träger-Zahl driftet · Fitness-Function-Zusage ohne Träger · Beleg ohne
  Eingabeseiten-Bindung · Zahl in Träger) — keine Klasse doppelt zu diesem Diff
- Beobachtungs-Register `BEO-PGC/` (Zähler mit `ls evidence/` nachgezählt, nicht
  aus Prosa übernommen)

---

## Findings

### F-1 — Die neue Sensors-Zelle ist zum Absatz geworden; die Sektion nennt dafür einen anderen Ort

- `kategorie`: MEDIUM
- `quelle`: `harness/README.md` §Sensors, Kommentar-Block über der Tabelle:
  Z. 93–98 („Braucht ein Gate mehr als **EINEN SATZ** — Deckungsgrenze,
  Ausgabe-Bedeutung, Exit-Codes, Abbruch-Bedingungen —, wandert das nach. …
  eine Zelle, die zum Absatz geworden ist, ist der Fund, nicht die Ausnahme.
  `harness/sensors/<target>.md`, und die **Target-Zelle wird zum Link darauf**“) ·
  Z. 106–108 („Was das Werkzeug selbst deckt … gehört NICHT dorthin, sondern in
  seine ADR/Spec-Zeile/seinen Skriptkopf“)
- `pfad`: `harness/README.md:118` (gegen `:93-98`, `:108`, gegen die Nachbarzeilen)
- `befund`: Die neue Vertrags-Zelle trägt **440 Zeichen** und drei der vier
  namentlich genannten Überhang-Arten in einer Zelle: Deckungsgrenze (welche
  Paarung verglichen wird), Ausgabe-Bedeutung (Temp-Verzeichnis, `:ro`-Bind-Mount)
  und Befund-Form (Datei und Zeile). Die Sektion nennt als Ort für genau diesen
  Überhang `harness/sensors/<target>.md` mit der Target-Zelle als Link; ein
  `harness/sensors/generated-sync.md` legt der Diff nicht an. Gemessene
  Zell-Längen derselben Spalte: `generated-sync` 440 · `commit-traceability` 295 ·
  `coverage-gate` 286 (Letztere hat ein Sensor-Dokument) · `a-check` 123 ·
  `docs-check` 101 · `baseline-verify` 91. Die drei Mechanik-Aussagen stehen
  daneben bereits im Skriptkopf (`generated-sync.sh:11-20`) und in `ADR-0084`
  §Entscheidung 1 — der Träger führt sie also doppelt; heute stimmen sie überein
  (selbst nachgemessen), morgen ist der Zell-Text die zweite Fassung ohne Gewinner.
  Gegen-Prezedenz, offen benannt: die Zeile `commit-traceability` (295 Zeichen,
  kein Sensor-Dokument, seit slice-006) steht in derselben Spannung und wurde nie
  als Fund geführt.
- `verifizierbar`: nein — kein Gate prüft Zell-Länge oder die Existenz eines
  Sensor-Dokuments je Target; der Link-Sensor prüft nur vorhandene Links, und
  `make docs-check` ist grün (`735 Datei(en), 0 Befund(e)`)
- `klasse`: „Sensors-Zelle trägt ihren Überhang statt eines Sensor-Dokuments“

### F-2 — Die genannte „erste Abweichung“ ist der Hunk-Anfang, nicht die abweichende Zeile

- `kategorie`: MEDIUM
- `quelle`: `ADR-0084` §Entscheidung 1, zweite Bedingung („Der Befund nennt den
  Diff. Ein rotes Gate zeigt, **welche** Datei und welche **Zeile** abweichen“) ·
  `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist) · §3.12 Instanz B
- `pfad`: `tools/harness/generated-sync.sh:75-77` (gegen den Skriptkopf `:22-23`)
- `befund`: Die Zeile wird aus dem **Hunk-Kopf** des Unified-Diffs gezogen
  (`sed -n 's/^@@ -\([0-9]\{1,\}\).*/\1/p'`), also aus dem Kontext-Anfang bei
  `diff -u` (drei Zeilen Kontext) — nicht aus der ersten abweichenden Zeile.
  Gemessen (Mutation einer einzigen Zeile im committeten Erzeugnis): die
  abweichende Zeile ist **102** (`grep -n 'GetSequenceMUTATED'` → `102:`), der
  eigene Nachbau des Generatorlaufs bestätigt das unabhängig mit
  `diff --unified=0` → `@@ -102 +102 @@`; die Ausgabe des Gates behauptet
  „erste Abweichung: committete **Zeile 99**“. Der Skriptkopf sagt „je abweichender
  Datei die committete Zeile **der ersten Abweichung**“ — gemessen ist es der
  Hunk-Anfang; die Abweichung liegt systematisch ein bis drei Zeilen darunter
  (gleich ist sie nur, wenn die Änderung auf Zeile 1 liegt). Die formale Hälfte
  der Festlegung („Datei **und** Zeile“) ist damit erfüllt, der Diff folgt
  unmittelbar dahinter — die Aussage über *welche* Zeile ist es aber nicht.
- `verifizierbar`: ja — eine Ein-Zeilen-Mutation im committeten `.pb.go` und
  `diff --unified=0` gegen den Generatornachbau
- `klasse`: „Zeilenangabe nennt den Hunk-Anfang, nicht die Abweichung“

### F-3 — Der Quellverzeichnis-Override kann ein gedriftetes `proto/` grün stellen und steht in keiner Vertragsstelle

- `kategorie`: LOW
- `quelle`: Maintainability / Haus-Muster der Vertragszelle
  (`harness/README.md:116`, Zeile `commit-traceability`, die ihren
  Bereichs-Override `RANGE=base..head` **in der Vertrags-Zelle nennt**)
- `pfad`: `tools/harness/generated-sync.sh:35`, `:54`, `:66` ·
  `harness/mk/generated-sync.mk:14-19`
- `befund`: `GENERATED_SYNC_SOURCE_DIR` ersetzt die verglichene Quelle. Gemessen:
  bei einer real gedrifteten `proto/cdc/stream/v1/changestream.proto` läuft der
  Aufruf ohne Override rot (`EC=1`, Zeilenbefund), mit
  `GENERATED_SYNC_SOURCE_DIR=<unveränderte Quelle im Baum>` grün (`EC=0`). Die
  Erfolgszeile nennt die benutzte Quelle (`generated-sync:   Quelle: …`) — der
  Pfad ist also **nicht still**, und das Fragment reicht die Variable nur aus der
  Umgebung durch (`GENERATED_SYNC_SOURCE_DIR=$(GENERATED_SYNC_SOURCE_DIR)`, im
  Regelfall leer). Ungenannt ist der Override in beiden Vertrags-Zellen
  (`harness/README.md:118`, `AGENTS.md:439`), während der Skriptkopf `:28-29`
  vier Overrides auflistet.
- `verifizierbar`: ja — der Aufruf mit gesetztem Override auf einem Baum mit
  gedrifteter `proto/`
- `klasse`: „Override kann den Befund grün stellen und steht nicht im Vertragsträger“

### F-4 — Eine zweite Enumeration der Gate-Liste ist stehen geblieben (Träger außerhalb des Diffs)

- `kategorie`: INFO
- `quelle`: `.harness/skills/reviewer.md` HIGH-Unterpunkt „Zahl im Träger ohne
  Ursprung — oder gegen die Messung driftend“, *Träger außerhalb des Diffs
  bleiben INFO* · Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- `pfad`: `.github/workflows/ci.yml:10` und `:65` (gegen `harness/README.md:119`,
  gegen `make -p`)
- `befund`: Der Diff zieht die Enumeration in `harness/README.md:119` auf sechs
  Namen nach (nachgemessen: deckungsgleich mit `GATE_CHECKS`). Die zweite Fassung
  derselben Liste im Workflow nennt weiterhin vier — `make gates` „deckt
  baseline-verify, docs-check, a-check und commit-traceability ab“ (Z. 10) und der
  Schrittname `Gates (baseline-verify, docs-check, a-check, commit-traceability)`
  (Z. 65). Es fehlen `coverage-gate` (seit slice-049) und `generated-sync` (dieser
  Zug); die Drift ist damit älter als dieser Diff. `ci.yml` weist in derselben
  Zeile auf `harness/README.md` §Sensors als Quelle — der Gewinner ist deklariert,
  es ist eine stehen gebliebene Zweitfassung, kein unentschiedener
  Zwei-Quellen-Fall.
- `verifizierbar`: ja — `make -p | grep '^GATE_CHECKS'` gegen die beiden
  `ci.yml`-Zeilen
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Träger außerhalb
  des Diffs)

### F-5 — Grenze: liegt `TMPDIR` im Baum, wird der Lauf falsch-rot

- `kategorie`: INFO
- `quelle`: `ADR-0084` Festlegung 1 („Das Gate schreibt den Arbeitsbaum nicht“) ·
  `AGENTS.md` §3.7 (Grenze)
- `pfad`: `tools/harness/generated-sync.sh:60` gegen `:110-111`
- `befund`: Das Temp-Verzeichnis folgt `${TMPDIR:-/tmp}`; Richtung 2 sucht die
  gekennzeichneten Erzeugnisse mit `grep -rl … "$repo_root"`. Gemessen mit
  `TMPDIR` auf ein Verzeichnis **innerhalb** des Baums: `EC=1`, und die beiden
  Rot-Zeilen nennen die eigenen Temp-Dateien
  (`scratch/pg-change-feed-generated-sync.<X>/internal/…/changestream.pb.go ist
  als Erzeugnis gekennzeichnet …`). Der Lauf schreibt in diesem Fall während des
  Laufs in den Baum; der `trap` räumt ihn danach wieder weg, `git status
  --porcelain` ist anschließend leer (die Zusage aus LP2 hält also auch dann).
  Die Richtung des Fehlers ist die **sichere** (rot, nicht grün).
- `verifizierbar`: ja — `TMPDIR=<Verzeichnis im Baum> bash tools/harness/generated-sync.sh`
- `klasse`: „Grenze: TMPDIR innerhalb des Baums macht den Lauf falsch-rot“

---

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/generated-sync.sh:63-69` (Generatoraufruf) —
  der Baum hängt als `-v "$repo_root":/src:ro`, geschrieben wird nur nach `/out`
  (Temp-Verzeichnis); zwei Läufe hintereinander in einem Klon: je `EC=0`, danach
  `git status --porcelain` **leer** in beiden Fällen. Ebenso im Original-Repo:
  Lauf davor und danach leerer Status. `ADR-0084` Festlegung 1 ist real getragen,
  nicht behauptet. Die Modul-Layout-Relation hält: `--go_out=/out` +
  `--go_opt=module=<Modulpfad aus go.mod>` legt die Ausgabe unter genau den
  Pfaden ab, unter denen das Erzeugnis im Baum liegt
  (`internal/adapters/driving/grpc/streamv1/changestream[.pb.go|_grpc.pb.go]`) —
  belegt durch den grünen Lauf und durch den unabhängigen Generatornachbau.
- geprüft, ohne Befund: **beide Vergleichsrichtungen.** Richtung 1: committete
  Datei gelöscht → `FAIL — … changestream_grpc.pb.go fehlt im Baum (der Generator
  erzeugt sie)`, `EC=1`. Richtung 2: eine zusätzliche, als `protoc-gen-go`
  gekennzeichnete Datei in den Baum gelegt → `FAIL — … ghost.pb.go ist als
  Erzeugnis gekennzeichnet (protoc-gen-go), der Lauf erzeugt sie aber nicht`.
  Die Paarung ist über beide Richtungen geschlossen. (`gelesen, nicht gemessen`:
  beide Suchen laufen in Prozess-Substitutionen, deren Exit-Status `set -e` nicht
  sieht — ein fehlschlagender `grep` würde Richtung 2 still verkürzen; Richtung 1
  deckt dann noch Inhalt und Vorhandensein jeder erzeugten Datei ab.)
- geprüft, ohne Befund: **`set -euo pipefail` und `trap`.** `trap 'rm -rf …' EXIT`
  wird erst nach `mktemp` gesetzt; die beiden frühen Ausgänge (kein Modulpfad,
  keine Quelle) liegen davor — ein `rm -rf` auf leerer Variable ist damit
  strukturell nicht erreichbar. Beide Ausgänge laufen **laut** rot: `go.mod`-Zweig
  `EC=2` mit Klartextmeldung; leeres und nicht existierendes Quellverzeichnis
  `EC=2` mit `FAIL — keine .proto-Quelle unter …`. Bash 5.2 auf dem Aufrufer,
  `bash tools/…` in Rezept und Aufruf identisch zum Haus-Muster der übrigen
  `tools/harness/*.sh`.
- geprüft, ohne Befund: **der Bild-Override öffnet keinen Grün-Pfad.**
  `docker build --target proto -t "$GENERATED_SYNC_IMAGE" .` baut die Stufe
  **immer** aus dem lokalen `Dockerfile` und taggt sie neu — ein gesetzter
  `GENERATED_SYNC_IMAGE` kann keinen ungepinnten Generator unterschieben. Die
  Pinnung selbst stimmt mit dem Skriptkopf überein (`protobuf-dev=31.1-r1`,
  `protoc-gen-go@v1.36.12`, `protoc-gen-go-grpc@v1.6.2` gegen `Dockerfile:33-36`).
  `GENERATED_SYNC_MODULE`: ein **nicht** passender Präfix bricht den Generator mit
  genau der im Kommentar zitierten Meldung ab
  (`--go_out: …: generated file does not match prefix "example.com/other"`,
  `EC=1`); ein kürzerer, noch passender Präfix verschiebt die Ausgabe und wird
  von beiden Richtungen rot gemeldet (`EC=1`) — die Zusage „landet dann auf keinem
  falschen Pfad, sondern der Lauf ist an dieser Stelle rot“ hält in beiden Fällen.
  `GENERATED_SYNC_RUN_USER` ist über `make generated-sync GENERATED_SYNC_RUN_USER=…`
  erreichbar (gemessen: GNU make exportiert Kommandozeilen-Variablen in die
  Rezept-Umgebung), obwohl das Fragment nur drei der vier Übergaben explizit
  durchreicht.
- geprüft, ohne Befund: **die Zusage des Quellverzeichnis-Kommentars** (`:52-53`):
  eine zusätzliche `.proto` unter `proto/` wird real mitgenommen
  (`FAIL — … streamv1probe/probe.pb.go fehlt im Baum`, `EC=1`).
- geprüft, ohne Befund: **der reale Drift-Fall.** Eine Feldänderung in der `.proto`
  bei stehen gelassenem committetem `.pb.go` → `EC=1`, Zeilenbefund und
  Unified-Diff mit allen vier Hunks. Das Gate fängt genau die Klasse, für die es
  gebaut ist.
- geprüft, ohne Befund: **die Bindung ist vollständig und einfach deklariert.**
  `make -p` löst auf: `GATE_CHECKS := baseline-verify coverage-gate docs-check
  commit-traceability generated-sync a-check` (6 Einträge); `record-gates:
  $(GATE_CHECKS)` ist die Ordnungskante im Aggregator. `grep -rn '^generated-sync:'`
  und `grep -rn GENERATED_SYNC --include='*.mk' --include=Makefile` liefern genau
  **einen** Ort (`harness/mk/generated-sync.mk`); der `Makefile` ist unberührt, wie
  §3 des Plans es festhält — keine zweite Deklaration.
- geprüft, ohne Befund: `AGENTS.md:439` — ein Satz, nennt Vertrag (`byte-gleich`),
  Quelle und die zwei ADRs; kein halluziniertes Target, kein Erweiterungs-Satz,
  der im Pflichtenheft stehen müsste.
- geprüft, ohne Befund: **§3.11** — keine host-lokalen absoluten Pfade in den
  hinzugefügten Zeilen. **§3.12 Instanz A** — die geänderten Doku-Zeilen führen
  **keine** Messzahl; die Zahlen im Umfeld (`5 Commits`, Rampe `70 → 80 %`,
  `74.70 %`) wurden nachgemessen und halten. Ebenso die Pinnings.
- geprüft, ohne Befund: **§3.7, die Kommentare.** Skriptkopf (`:2-29`) und
  Fragmentkopf tragen Zusage, Kopplung (Bind-Mount, Aufrufer-uid wie
  `PROTO_RUN_USER`, die Modul-Layout-Relation), Abgrenzung („Nicht Gegenstand:
  dass der erzeugte Code kompiliert oder sich richtig verhaelt“) und Rang-Zeiger
  (`ADR-0084` Festlegung 1, `ADR-0060`). Keine Chronik, kein Vorher/Nachher, keine
  Slice-/Wellen-Nummer im Mechanik-Kommentar. Die Begründung, warum die Stufe ohne
  `--no-cache-filter` läuft, hält: das Urteil fällt im Vergleich des Laufs,
  außerhalb der Stufe. Als **Fund** steht die *Zusage* über die Abweichungszeile in
  F-2, nicht die Klasse der Kommentare.
- geprüft, ohne Befund: `harness/mk/generated-sync.mk` (22 Zeilen) — `.PHONY`
  gesetzt, Haus-Form (`?=`-Default, `##`-Hilfezeile für `make help`), keine
  Lockerung (§3.6: ein neues Gate ist keine Schwellen-Senkung), kein
  Docker-only-Verstoß (§3.1).
- geprüft, ohne Befund: **Plan-vs-Code** — §3 des Plans sagt `Makefile` = „**nicht**“,
  und der Diff hält das ein; der Zielname ist der im Plan genannte; die Doku-Punkte
  des §3 sind beide geliefert. Kein `git mv`, damit §3.3 unberührt.
- geprüft, ohne Befund: **Register-Zähler des Plans §8** — `ls evidence/`:
  `generierte-artefakte-ohne-sync-sensor` **4×** (`slice-069/-074/-082/-086`),
  `regel-weiter-als-ihr-sensor` **2×**, `gate-scope-erweiterung-ohne-adr-traeger`
  **1×** — alle drei wie im Plan genannt. Dass `state.md` dieses Slice **nicht**
  schreibt, ist die richtige Disziplin (Ausgang = Lese-Schritt der
  `welle-20`-Closure, Modul 6).
- geprüft, ohne Befund: Hygiene des Ranges — keine hinzugefügten Lauf-Artefakte,
  keine Vorlagen-Reste in den neuen Dateien, Betreff mit `ADR-*` und ohne
  Struktur-ID, Arbeitsbaum sauber, `harness/sensors/` unverändert.
- geprüft, ohne Befund: `make gates` real — **Exit 0** aus separater Datei gelesen
  (§3.9), alle sechs Checks gefahren, danach leerer `git status --porcelain`. Die
  „billig“-Zusage aus `ADR-0084` §Konsequenzen hält warm:
  `time make generated-sync` → **0,9 s**.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-status 9029c05^..9029c05` | **0** | 5 Dateien (2 `A`, 3 `M`) |
| `make -p \| grep '^GATE_CHECKS'` | **0** | 6 Einträge, `generated-sync` enthalten |
| `make gates` (Log in Datei, Exit danach separat gelesen) | **0** | alle sechs Checks; danach `git status --porcelain` leer |
| `make generated-sync` (Original-Baum) | **0** | `OK — byte-gleich`; Status leer |
| `time make generated-sync` (warmer Cache) | **0** | `real 0m0,900s` |
| Klon: zwei Läufe, je `git status --porcelain` | **0 / 0** | beide leer → Festlegung 1 real getragen |
| Klon: Zeile mutiert; `grep -n 'GetSequenceMUTATED'` | **0** | abweichende Zeile ist **102** |
| Klon: Generatornachbau + `diff --unified=0` | **0** | `@@ -102 +102 @@` — unabhängige Bestätigung |
| Klon: `generated-sync.sh` bei dieser Mutation | **1** | Befund nennt „Zeile **99**“ → **F-2** |
| Klon: zusätzliche `ghost.pb.go` | **1** | Richtung 2, `FAIL — … ist als Erzeugnis gekennzeichnet` |
| Klon: committete `changestream_grpc.pb.go` entfernt | **1** | Richtung 1, `FAIL — … fehlt im Baum` |
| Klon: Feld `review_probe = 11` in die `.proto` | **1** | Zeilenbefund + 4 Hunks — realer Drift gefangen |
| Klon: zusätzliche `probe.proto` unter `proto/` | **1** | neue Quelle tritt von selbst ein |
| `GENERATED_SYNC_MODULE=example.com/other …` | **1** | `generated file does not match prefix` |
| `GENERATED_SYNC_SOURCE_DIR=<leer/nicht vorhanden>` | **2 / 2** | `FAIL — keine .proto-Quelle unter …` — kein stiller Grün-Pfad |
| Klon: `proto/` gedriftet + Quell-Override | **0** | `OK …` + `Quelle: …` (Override) → **F-3** |
| Klon: `TMPDIR=<Verzeichnis im Baum>` | **1** | zwei Rot-Zeilen über die eigenen Temp-Dateien → **F-5** |
| Zell-Längen der Vertragsspalte | **0** | 440 · 295 · 286 · 123 · 101 · 91 → **F-1** |
| `ls harness/sensors/` | **0** | kein `generated-sync.md` → **F-1** |
| `ls observations/BEO-PGC/<slug>/evidence/` | **0** | 4 · 2 · 1 — deckt §8 des Plans |
| `grep -rn '^generated-sync:'` / `GENERATED_SYNC` | **0** | nur `harness/mk/generated-sync.mk` — eine Deklaration |
| `git show 9029c05 \| grep '^+'` gegen die Präfixliste `hostpaths.prefixes` | **1** | keine Ausgabe — §3.11 gehalten (die Präfixe selbst stehen hier nicht, §3.11) |
| `git show 9029c05 -- <Doku-Träger> \| grep -nE '[0-9]+'` | **0** | nur Kennungen — keine Messzahl (§3.12 A) |

Der Gate-Lauf und seine Auswertung sind **zwei Schritte**: `make gates` lief
ungepiped in eine Log-Datei, der Exit-Code separat gelesen und erst danach in
einem eigenen Werkzeug-Aufruf ausgewertet; kein Folgekommando hing an einer Pipe
oder an einem Wrapper (§3.9).

## Antwort auf die Schwerpunkte

**(1) Das Skript — trägt, mit einer Aussage-Drift (F-2).** Der Generatoraufruf
bildet die Relation von `make proto-generate` nach, die Quellpfade sind relativ und
korrekt gequotet (`mapfile` + `"${proto_sources[@]}"`), beide Vergleichsrichtungen
sind real rot-fähig, der `trap` sitzt nach `mktemp` und vor jedem frühen Ausgang,
und jede der vier `GENERATED_SYNC_*`-Variablen wurde auf Grün-Pfad-Tauglichkeit
geklopft: Das **Bild**-Override kann den Generator nicht austauschen, **Modul** und
**Quelle** führen in jeder geprüften Stellung zu einem lauten Ergebnis — nur bei
der Quelle in eine Richtung, die grün ist und dabei auf eine andere Quelle zeigt
(F-3, im Erfolgstext sichtbar). Die Herleitung der ersten Abweichungszeile ist der
einzige Punkt, an dem die Ausgabe etwas anderes sagt als die Messung (F-2).

**(2) `ADR-0084` Festlegung 1 und 2 — beide real getragen.** Festlegung 1: Baum
`:ro`, Schreiben nach `/out`, zwei Läufe hintereinander mit leerem `git status` —
gemessen, nicht gelesen. Festlegung 2: Datei und Zeile kommen; nur die Zeile ist
der Hunk-Anfang (F-2), der Diff folgt unmittelbar.

**(3) §3.7 — die Kommentare tragen Klassen, kein Chronik-Ton.** Zusage ·
Kopplung · Abgrenzung · Rang-Zeiger, je Stelle geprüft. Der einzige Kommentar, der
gegen die Messung steht, ist die Zusage über die „erste Abweichung“ (F-2).

**(4) §3.12 — kein Fund in Instanz A, das Prüfraster greift bei der Ausgabe.**
Die geänderten Doku-Zeilen führen keine Zahl; die Zahlen in ihrem Umfeld wurden
**nachgemessen** und halten. Instanz-B-relevant ist die Ausgabe-Aussage des
Skripts, als F-2 geführt — der durchsetzende Leser für Instanz B bleibt der
Verifier.

**(5) Vollständigkeit der Bindung — trägt.** `GATE_CHECKS += generated-sync` steht
im Fragment, `include harness/mk/*.mk` zieht es ein: `make -p` löst sechs Einträge
auf, `make -n record-gates` zeigt den Aufruf als Prüfschritt vor dem Stempel, und
es gibt **keine** zweite Deklaration. Die `make gates`-Zeile in
`harness/README.md:119` ist deckungsgleich; die stehen gebliebene zweite Fassung
liegt **außerhalb** des Diffs (F-4).

**Was nicht geprüft wurde:** DoD/LP1–LP3 (Verifier), die §6-Risiko-Ausgänge und
die drei Paarungen (Planner-Closure), und ob der neue Gate-Schritt im kalten
CI-Lauf innerhalb des Timeouts bleibt (dafür fehlt der reale Post-Push-Lauf;
`AGENTS.md` §3.10 gilt für den betroffenen Workflow, den dieser Diff **nicht**
ändert).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Sensors-Zelle trägt ihren Überhang statt eines
Sensor-Dokuments“ (F-1) · „Zeilenangabe nennt den Hunk-Anfang, nicht die
Abweichung“ (F-2) · „Override kann den Befund grün stellen und steht nicht im
Vertragsträger“ (F-3) · `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (F-4,
Träger außerhalb des Diffs) · „Grenze: TMPDIR innerhalb des Baums macht den Lauf
falsch-rot“ (F-5)

## Verdikt

**Merge-blockierend:** nein — **0 HIGH**.

**Auf der Implementer-Seite:** F-1 und F-2 sind Ein-Stellen-Korrekturen in Trägern
dieses Diffs (`harness/README.md` ist eines der zwei Doku-Ziele, die der Slice
selbst benennt; `tools/harness/generated-sync.sh` ist der Kern des Lieferwerts).
Sie erfordern eine **Rückgabe an den Implementer**; F-2 betrifft die Ausgabe, die
`ADR-0084` Festlegung 2 trägt, F-1 die Form des Sensors-Trägers.

**Planner-Seite:** F-4 ist eine Adresse, die der Diff nicht schuldet
(`.github/workflows/ci.yml`, Drift älter als dieser Zug); sie gehört als Fundstelle
benannt, nicht in diesen Zug hinein.

**INFO:** F-5 ist die benannte Grenze des neuen Gates (Fehlerrichtung rot); F-3
ist die Frage, ob der Quell-Override in der Vertragszelle oder im Sensor-Dokument
genannt wird — sie hängt an derselben Entscheidung wie F-1.

**Übergabe:** Rückgabe-Pfeil Reviewer → Implementer ist **nötig** (F-1, F-2); die
Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation
(Modul 11).

**DoD-Häkchen „Review durchgeführt, Report unter `docs/reviews/` liegt vor“:**
bleibt **offen** — der Slice braucht eine Fixrunde.
