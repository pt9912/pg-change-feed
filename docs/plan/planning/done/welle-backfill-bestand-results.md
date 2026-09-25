# Welle welle-backfill-bestand — Backfill des Bestands: der Tabellenbestand einer aktivierten Tabelle wird als Backfill erkennbar über den bestehenden Lesezugriffsweg lesbar — Closure-Notiz

**Welle:** welle-backfill-bestand
**Abschluss:** 2026-09-25
**Verantwortlich:** pt9912

## Was wurde geliefert?

- **[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Happy Path, Boundary,
  Negative) ist am laufenden Feed-Container belegt.** Ein Antrag
  `cdc.backfill_table` legt einen Run an (`queued`); der Worker liest den Bestand
  im `REPEATABLE READ`-Snapshot eines je Run angelegten temporären logischen Slots
  und committet ihn in **einer** Store-Transaktion als `INSERT` mit
  `origin = 'backfill'`; der Bestand ist über `cdc.changes` und `GET /changes`
  lesbar, unterscheidbar von WAL-Changes, und an den WAL-Pfad lückenlos
  angeschlossen ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)).
  Zustand und Warnungen des Runs stehen in `cdc.backfill_status` und `diagnose`.
- **Elf Slices in `done/`:** `slice-backfill-spec-nachzug` (Pflichtenheft,
  Architektur-Sicht, [`SPEC-029`](../../../../spec/pflichtenheft.md)),
  `slice-backfill-row-image-gemeinsam` (eine gemeinsame Row-Image-Funktion),
  `slice-backfill-change-origin` (Feld `origin` in Domäne, Store, View,
  `GET /changes`), `slice-backfill-snapshot-reader` (Slot mit Export-Snapshot,
  Cursor-Blöcke), `slice-backfill-run-usecase` (Run, Use Case, Fail-closed-Prüfung),
  `slice-backfill-run-store` (Tabelle `cdc.backfill_run`, atomarer Schreiber),
  `slice-backfill-sql-administration` (Antragsart, Worker, Start-Abgleich, View,
  `diagnose`, Idempotenz-Guard), `slice-backfill-e2e` (Happy Path, Boundary,
  Negative in `make test-integration`), `slice-backfill-bench-richtgroesse`
  (Bench der Kopierdauer, Warn-Richtgröße, Warnungen im Use Case),
  `slice-backfill-slot-leerlauf-bestaetigung` (Capture-Slot bestätigt im Leerlauf),
  `slice-backfill-sdk-origin` (`origin` in den drei SDK-HTTP-Lesemodellen).
- **Folge-Entscheidungen der Welle:**
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Rollenschnitt, Aufnahme, Warnkriterium),
  [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) (Vorlauf bei View-Signatur-Änderung),
  [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) (Spaltenwerte im Text-Ergebnisformat),
  [`ADR-0116`](../../adr/0116-backfill-schema-version-referenz-reichweite.md) (Reichweite der Schema-Version-Referenz),
  [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (Fehlerklasse `schema` im Run),
  [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md) (Umschreiben im Snapshot-Fenster),
  [`ADR-0119`](../../adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md) (Wirkung der Lesesperre, Berichtigung),
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) (Slot-Bestätigung im Leerlauf),
  [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md) (Store-Bindung der Leerlauf-Bedingung, Berichtigung) und
  [`ADR-0122`](../../adr/0122-backfill-replay-invariante-e2e-tier.md) (Tier der Replay-Invariante).
  Die Architect-Verdikte stehen unter `docs/reviews/`:
  `architect-verdict-backfill-schema-klasse-rollen`,
  `architect-verdict-backfill-tabellen-rewrite-im-fenster`,
  `architect-verdict-backfill-wal-rueckstand-und-bench-rot`,
  `architect-verdict-schema-rollout-view-signatur` und
  `architect-verdict-welle-backfill-bestand-lese-schritt`.
- **Benutzerhandbuch:** Änderungshistorie 1.44 bis 1.59 (16 Zeilen; Ursprung:
  die Zeilen `| 1.44 |` bis `| 1.59 |` der Änderungshistorie, Stand `32028d4f`).
  Der Abschnitt „Bestand als Backfill überführen“ trägt Auslösung,
  Sichtbarkeits-Grenze, Startposition eines neuen Consumers, Lese-Regel ohne
  `LIMIT`, Warnungen und die Bedeutung des WAL-Rückstands; „Grenzwerte“ trägt
  Richtgröße, Toleranz und die Grenze der Ein-Transaktions-Form.
- **RTM:** `make doc-trace` druckt `80 Anforderung(en), 2 Waise(n).` — die zwei
  Waisen sind [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) und
  [`LH-FA-CFG-008`](../../../../spec/lastenheft.md); [`LH-FA-CAP-009`](../../../../spec/lastenheft.md)
  trägt den Nachweis `E2E` über die Zeilen von `docs/user/e2e-abdeckung.md`
  (gemessen am Stand `32028d4f`).
- **Vorbereitung der Welle Transformationen:** die Kanten K1 bis K3 aus
  `welle-backfill-bestand` §5 sind erfüllt (K1: `model.BuildRowImage` liegt vor;
  K2: [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) ist der
  Start-Trigger von `slice-transformationen-backfill-pfad`; K3: die
  `request_kind`-Menge trägt fünf Werte, gemessen mit
  `sed -n 62p tools/schema/nacharbeit-administration.sql`).

## Was hat funktioniert?

- **Der Schnitt in elf Slices mit benannten Kanten** (Welle-Datei §4 und §5)
  trug die Welle bis zur Closure; die Kanten K1 bis K3 zur Welle Transformationen
  sind erfüllt (siehe „Was wurde geliefert?“). Der elfte Slice entstand aus einer
  Messung (siehe unten).
- **Die Rollenteilung fängt die Fehler vor dem Merge.** Der Lese-Schritt maß, dass
  die vier Mega-Einträge (`arbeit-ueberholt-stehenden-traeger`,
  `zahl-in-traeger-driftet-gegen-die-messung`, `beleg-befehl-traegt-seinen-satz-nicht`,
  `negativtest-ohne-bindung-an-seine-eingabe`) in dieser Welle ausnahmslos von
  einem nachmessenden Leser (Reviewer oder Verifier) gefunden wurden, zwei davon
  als HIGH (`review-slice-backfill-row-image-gemeinsam` F-1,
  `review-slice-backfill-sdk-origin` F-1); Quelle:
  `architect-verdict-welle-backfill-bestand-lese-schritt` §3.1 und §3.2.
- **Architect-Verdikte an den Stellen, an denen eine Messung dem Plan
  widersprach.** Drei Befunde der Welle (Umschreiben der Tabelle im
  Snapshot-Fenster, WAL-Rückstand eines Runs, Fehlerklasse und Rollen des Runs)
  sind je in einer ADR entschieden
  ([`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md),
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md),
  [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)); der jeweilige
  Folge-Slice startet auf dem `Accepted`-Stand.
- **Belege an beiden PostgreSQL-Versionen im CI.** Der Lauf `36108615045` von
  `e2e.yml` (Stand `c82d3333`, beide Legs) druckt die Backfill-Phasen (Zeilen im
  Abschnitt Verifikation).

## Was ging anders als geplant?

- **Umschreiben der Tabelle im Snapshot-Fenster.** Ein `ALTER COLUMN TYPE` oder
  `DROP COLUMN` zwischen Snapshot-Export und Snapshot-Import beendet den Run mit
  `failed` ohne Change (E2E-Phase „Backfill-DDL-Fenster“); die Wirkung der
  Lesesperre steht in
  [`ADR-0119`](../../adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md)
  berichtigt (eine Aussage von [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md)
  war breiter als ihre Messung). Konsequenz: [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md).
- **WAL-Rückstand ist ein Defekt des Capture-Pfads, kein Backfill-Effekt.** Eine
  Messung sah einen Run die Fehlerschwelle des WAL-Rückstands überschreiten; die
  Ursache liegt in der Empfangs-Schleife des Capture-Slots, der im Leerlauf nicht
  bestätigte. Konsequenz: der elfte Slice
  `slice-backfill-slot-leerlauf-bestaetigung`
  ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md), Welle-Datei §4
  „Herkunft des elften Slice“); die Store-Zeile der Fitness Function war nicht
  erfüllbar und ist in
  [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
  berichtigt.
- **`make bench` als Ganzes endet am Messhost mit Exit 2** an
  [`LH-QA-PER-001`](../../../../spec/lastenheft.md) (Ursache laut
  `architect-verdict-backfill-wal-rueckstand-und-bench-rot`, übernommen: Flush-Latenz
  des Hosts; die Kontrolle mit `synchronous_commit=off` senkte den Overhead um
  6,7 %). Die Welle bindet ihr Kriterium deshalb an `tools/bench-backfill.sh`
  (einzeln); der Re-Evaluierungs-Trigger von
  [`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md) ist nicht
  eingetreten (Trigger-Audit unten).
- **Die Replay-Invariante liegt in einem anderen Tier als die Fitness-Function-Zeile
  sagt.** [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt
  `make test-replication`, der Beleg liegt in `make test-integration`
  (`TestE2EBackfillReplayInvariant`). Konsequenz:
  [`ADR-0122`](../../adr/0122-backfill-replay-invariante-e2e-tier.md).
- **Speicher des Feed-Containers ist bei großen Runs ungeklärt.** Das Handbuch
  („Grenzwerte“) nennt für die Stufe mit 200.000 Zeilen eine Spitze im Run von 401,7 bis
  467 MiB (401,7 und 467 MiB übernommen, im Repository nicht auflösbar; 416,5 MiB im
  Review-Report von `slice-backfill-bench-richtgroesse`) — die Spitze im Run ist eine Untergrenze: der höchste
  gemessene Wert der Stufe ist 641,7 MiB, 20 s nach dem letzten Run — und für zwei Runs
  über je 1.000.000 Zeilen 435 MiB und 1.544 MiB (übernommen, im Repository nicht
  auflösbar); der Lauf `20260925T084803Z` dieser
  Closure (gemessen) zeigt für die Stufe mit 200.000 Zeilen eine Spitze von
  433,0 MiB (Ruhe davor 170,8 MiB). Eine Blockgröße zählt Zeilen, nicht Bytes.
  Konsequenz: `slice-backfill-speicher-untersuchung`.
- **Die Warn-Richtgröße streut mit der Rate.** Der Lauf `20260925T084803Z`
  (gemessen, gedruckt) rechnet 6.470 Zeilen/s (Median, Stufe 200.000) mal 600 s =
  3.882.000 Zeilen, abgerundet 3.000.000; die Konstante im Code ist 4.000.000. Das
  Handbuch („Grenzwerte“) nennt die Spanne der bisherigen Läufe (2.000.000 bis
  5.000.000 nach Rundung) und die Konstante als Startwert; der Wert 3.000.000 liegt in
  dieser Spanne, die Konstante bleibt.
- **Die DB-Tier-Läufe hinterlassen `tools/schema/plan.yaml` verändert.**
  `make test-store`, `make test-replication` und `tools/bench-backfill.sh`
  (gemessen: `git status --short` nach jedem Lauf zeigt `M tools/schema/plan.yaml`);
  die Closure nimmt die Datei per `git checkout` zurück. Der Träger ist
  `slice-harness-suchlauf-nachmessen` (Klasse
  `test-schreibt-in-committete-datei`, 4×); dessen Plan nennt
  `tools/bench-lib.sh` und `tools/harness/run-integration-tests.sh` als weitere
  direkte Aufrufer von `make schema-rollout`, so dass der Zuschnitt der Rücknahme
  (in `apply-rollout.sh` oder im Target) dort zu entscheiden ist.

## Steering-Loop-Einträge

Die Quelle aller Einträge ist der Lese-Schritt dieser Closure
(`architect-verdict-welle-backfill-bestand-lese-schritt`); jede Regel liegt an
ihrem Zielort und trägt den Anker `seit welle-backfill-bestand`. Die Zähler sind
die Zahl der Dateien unter `evidence/` (`ls evidence | wc -l`, gemessen am Stand
`32028d4f`). Die Steering-Loop-Einträge je Slice stehen in den elf Slice-Notizen als
Sammelabsatz je Notiz; diese Notiz trägt die Regeln der Welle einzeln (im Verdikt
R1 bis R7, §7).

- **`AGENTS.md` §3.13 (Suchform, Frist der Meldung)** geschärft: Suchraum ist der
  ganze Baum (Ausnahmen namentlich), das Muster trägt drei Arten (Symbolname,
  Zählwort, Beschreibung samt Hedge), jeder Befehl steht im Codeblock, der Parent
  als Commit-Kennung, die Meldung eines Trägers in einer fremden Datei nennt ihre
  Frist (die Closure des meldenden Slice) — liegt in `AGENTS.md §3.13`.
  Auslöser: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (31×, Ausgang bleibt
  verkörpert, geschärft) und `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (8×,
  Ausgang verkörpert).
- **`.harness/skills/reviewer.md`, MEDIUM-Punkt „Nachzug widerspricht dem Nachbarn
  im selben Träger“** ergänzt (Probe: Kontext um jede hinzugefügte Zeile lesen,
  `git diff -U20`) — liegt in `.harness/skills/reviewer.md`. Auslöser:
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (8×; fünf Belege nennen ihre
  Schwere uneinheitlich als MEDIUM, LOW, INFO).
- **`.harness/skills/reviewer.md`, HIGH-Punkt „Kommentar trägt keine der
  Kommentar-Klassen“** um zwei Klauseln geschärft: *Skopus* (gilt für Tests und
  Runner `tools/harness/*.sh`) und *Zusage* (ein zugesichertes Verhalten wird im
  Code nachgefahren; liegt es in einem anderen Slice, trägt der Kommentar einen
  Rang-Zeiger) — liegt in `.harness/skills/reviewer.md`. Auslöser:
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (3×) und
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×).
- **`AGENTS.md` §3.12, Absatz „Verfasser einer ADR“** ergänzt (Instanz B für den
  Architect: eine Aussage über alle Werte einer Menge nennt die Menge, eine
  Fitness-Function-Zeile den tragenden Test, erprobt oder als hergeleitet
  gekennzeichnet), Zeiger in `.claude/agents/architect.md` — liegt in
  `AGENTS.md §3.12`. Auslöser: `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (5×:
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
  [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md),
  [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md),
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)).
- **`harness/sensors/db-adapter-coverage.md` §Gegenstand** ergänzt (Unterpaket-Regel:
  netzlos prüfbare Logik eines neuen DB-Pakets liegt von Beginn an in einem
  Unterpaket im Gegenstand des Unit-Gates, die Zuordnung steht im Slice-Plan vor
  dem Start) — liegt in `harness/sensors/db-adapter-coverage.md §Gegenstand`.
  Auslöser: `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (3×).
- **`docs/plan/planning/observations/README.md`, Deckel für verkörperte Einträge ab
  10×** — liegt in `docs/plan/planning/observations/README.md`. Auslöser: die vier
  Mega-Einträge (31×, 21×, 13×, 11×), deren Auftreten der Reviewer oder Verifier vor
  dem Merge fängt; ihre `state.md` nennen „Deckel bei N× (seit welle-backfill-bestand)“.
- **Sensor (geplant): Nachmess-Werkzeug für das Suchlauf-Feld** —
  `make suchlauf-nachmessen PLAN=<Datei>` führt die Suchzeilen eines Slice-Plans
  aus und druckt Soll und Ist; kein Gate, kein Sensor auf Prosa
  ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) bleibt
  unberührt) — liegt in `slice-harness-suchlauf-nachmessen`. Auslöser:
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (21×; jede der acht
  Belegdateien der Welle nennt ein Suchlauf-Feld, `grep -l -i suchlauf` gegen
  `ls`: 8 von 8). Derselbe Slice trägt die Rücknahme von `plan.yaml` und
  `down.sql` (Auslöser: `BEO-PGC/test-schreibt-in-committete-datei`, 4×).
- **Benannte Spec-Lücke:** [`SPEC-008`](../../../../spec/pflichtenheft.md) nennt für die
  Klasse `transient` „Erneut versuchen mit begrenztem Backoff“, kein Element des
  Erfassungspfads trägt die Aktion; der Pfad endet auf jeden Adapter-Fehler mit
  Prozess-Ausgang 1. Zustand: **offen** — die Entscheidung zur Wiederholungsform
  (Ort, Grenzen, Sichtbarkeit) trifft der Architect in einer ADR; sie liegt nicht vor
  (das ADR-Verzeichnis endet bei
  [`ADR-0122`](../../adr/0122-backfill-replay-invariante-e2e-tier.md)). Adresse: der
  Start-Trigger von `slice-capture-transient-wiederholung` (Plan §4) und
  `BEO-PGC/adapter-fehler-ausgang` (geplant). [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  Folgepflicht 1 betrifft die Zeile `schema` von
  [`SPEC-008`](../../../../spec/pflichtenheft.md), nicht `transient`.
  Auslöser: `BEO-PGC/adapter-fehler-ausgang` (3×: `evidence/slice-007.md`,
  `evidence/slice-038.md`, `evidence/slice-backfill-snapshot-reader.md`).
- **Ausdrücklich verworfen:** ein Mutations-Harness für
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (Docker-only-Pinnung samt
  Laufzeit je Paket, äquivalente Mutanten; die Regel wirkt, Begründung in
  `architect-verdict-welle-backfill-bestand-lese-schritt` §3.4) und ein Sensor
  „jede Zahl trägt ihren Ursprung“
  ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) §Entscheidung 4).

### Lese-Schritt des Beobachtungs-Registers

Das Register führt 113 Verzeichnisse (`ls docs/plan/planning/observations/BEO-PGC | wc -l`,
gemessen am Stand `32028d4f`), 39 davon mit `evidence/` ab 3 Dateien. Diese Closure
liest die **17** Einträge, deren Zähler in der Welle wuchs oder deren Ausgang nicht
zugewiesen war; die übrigen 22 Einträge ab 3× wachsen in der Welle nicht und tragen
einen zugewiesenen Ausgang (Liste in `architect-verdict-welle-backfill-bestand-lese-schritt`
§2). Jeder Ausgang steht in der `state.md` des Eintrags, die Verzeichnisse bleiben
liegen.

| Eintrag (`BEO-PGC/…`) | Zähler | Ausgang | Anker |
|---|---|---|---|
| `arbeit-ueberholt-stehenden-traeger` | 31 | verkörpert, geschärft; Deckel | `AGENTS.md §3.13` |
| `zahl-in-traeger-driftet-gegen-die-messung` | 21 | verkörpert; Suchlauf-Anteil geplant; Deckel | `slice-harness-suchlauf-nachmessen` |
| `beleg-befehl-traegt-seinen-satz-nicht` | 13 | verkörpert; Deckel | `AGENTS.md §3.13`, `slice-harness-suchlauf-nachmessen` |
| `negativtest-ohne-bindung-an-seine-eingabe` | 11 | verkörpert; Deckel; Mutations-Harness verworfen; die drei Kandidaten der Slice-Notizen (Fake ab Aufruf n, wertfreie Konstanten-Bindung, Filter-Eingabe): zwei gestrichen, eines akzeptiertes Negativ | `.harness/skills/reviewer.md`, `state.md` des Eintrags |
| `zitat-nennt-die-falsche-stelle` | 8 | verkörpert (unverändert) | `.harness/skills/reviewer.md` |
| `nachzug-laesst-ueberholten-text-stehen` | 8 | verkörpert | `AGENTS.md §3.13`, `.harness/skills/reviewer.md` |
| `github-actions-unverifizierbar-lokal` | 8 | verkörpert (unverändert) | `AGENTS.md §3.10` |
| `dod-begruendung-unzutreffende-tatsachenbehauptung` | 7 | verkörpert (unverändert) | `AGENTS.md §3.12` |
| `d-migrate-nacharbeit` | 7 | verkörpert (unverändert) | `harness/README.md` §Sensors |
| `adr-aussage-breiter-als-ihre-messung` | 5 | verkörpert | `AGENTS.md §3.12` |
| `test-schreibt-in-committete-datei` | 4 | geplant | `slice-harness-suchlauf-nachmessen` |
| `vorher-nachher-sprache-in-test-harness-kommentar` | 3 | verkörpert | `.harness/skills/reviewer.md` |
| `kommentar-behauptet-nicht-getragenen-fehlerpfad` | 3 | verkörpert | `.harness/skills/reviewer.md` |
| `rollen-test-abdeckungsluecken` | 3 | gestrichen (Punkt 1 geschlossen, Punkt 2 mit Begründung) | `state.md` des Eintrags |
| `db-gegenstand-enthaelt-netzlos-geprueften-code` | 3 | verkörpert | `harness/sensors/db-adapter-coverage.md` |
| `adapter-fehler-ausgang` | 3 | geplant | `slice-capture-transient-wiederholung` |
| `zusatzkontext-kopplung-breiter-als-dod-wortlaut` | 3 | gestrichen (der Fehler ist laut, ein Bau ohne den Kontext bricht ab) | `state.md` des Eintrags |

**Feststellung:** kein Eintrag über der Schwelle ohne Ausgang. Die 39 Einträge ab 3×
tragen 34 × verkörpert, 2 × geplant, 2 × gestrichen und 1 × eingetreten
(`rollen-verdrahtung`); Ursprung: die erste Ausgangs-Zeile jeder `state.md`
(`grep -m1 -o -i -E 'verkörpert|geplant|gestrichen|offen'`, Summe 39).

Einträge unter 3× mit Adresse an diese Closure — Ausgang in der `state.md`:
`sdk-decoder-verhalten-am-neuen-feld-ungemessen` (2×, gestrichen: Bibliothekssemantik;
Beleg ist eine Quelltext-Suche, kein Dekoder-Lauf — `git grep` nach strikten
Dekoder-Mustern druckt 0 Treffer an `sdk-csharp-v0.1.0` und an `sdk-python-v0.1.0`;
Restgrenze: die `0.1.0`-Packages sind nicht erneut ausgeführt, sie tragen dieselben
Dekoder-Zeilen; Muster und Quelle in der `state.md` des Eintrags und in
`architect-verdict-welle-backfill-bestand-lese-schritt` §5 (a)),
`lesesperre-ohne-zeitgrenze` (1×, gestrichen als akzeptiertes Negativ),
`setzpfad-einer-kennzeichnung-nur-im-unit-test-belegt` (1×, gestrichen als akzeptiertes
Negativ), `beleg-nur-als-einmalige-reviewer-messung` (1×, geplant →
`slice-capture-leerlauf-quellbelege`), `blockgroesse-zaehlt-zeilen-nicht-bytes` (1×,
geplant → `slice-backfill-speicher-untersuchung`),
`backfill-schema-version-hinter-snapshot-spalten` (1×, verkörpert →
[`ADR-0116`](../../adr/0116-backfill-schema-version-referenz-reichweite.md)),
`run-fehlerklasse-schema-im-transformations-backfill` (1×, verkörpert →
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)); offen bleiben
`architect-verdikt-rollen-scope-luecke`, `dod-kriterium-haengt-am-messhost` (Trigger
„zweite Umgebung führt `make bench` aus“ nicht eingetreten) und
`backfill-adapter-startwerte-ohne-messung`.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`docs/plan/planning/observations/`](../observations/)
(`BEO-PGC/<slug>/evidence/`); er wird nicht in dieser Notiz gepflegt. Was in
dieser Closure 3× erreicht hat, steht oben unter *Steering-Loop-Einträge*. Die
Adresse der Cursor-Form (`LIMIT` innerhalb einer Position) ist
`BEO-PGC/limit-fortsetzung-innerhalb-einer-position` (offen, „benannt, nicht
gezählt“, kein `evidence/`) und der Re-Evaluierungs-Trigger von
[`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md).

## Folge-Slices

Alle vier liegen als Datei in `docs/plan/planning/open/`:

- `slice-harness-suchlauf-nachmessen` — Nachmess-Werkzeug für das Suchlauf-Feld,
  Rücknahme von `plan.yaml`/`down.sql` in den DB-Tier-Läufen; geht
  `slice-transformationen-kern-rename` voraus.
- `slice-capture-leerlauf-quellbelege` — Keepalive inmitten einer Transaktion an
  PostgreSQL 17 als committeter Test im Tier `make test-replication`, E2E
  „Fehlerschwelle erreicht → Container endet“; geht
  `slice-transformationen-e2e-abhilfe` voraus; endet mit einer Ergänzung von
  [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md).
- `slice-capture-transient-wiederholung` — Wiederholung der Klasse `transient`;
  Start-Trigger: Architect-Entscheidung zur Wiederholungsform.
- `slice-backfill-speicher-untersuchung` — Untersuchung der Speicher-Spitze großer
  Runs, vor der ersten Veröffentlichung eines Server-Releases mit Backfill.

Die Welle Transformationen (`welle-transformationen`, zehn Slices in `open/`) startet
mit ihren Kanten K1 bis K3; zwei Kanten zu den Folge-Slices oben stehen in ihrem
§5.

## Validator-Feststellung (Modul 8)

Die Welle liefert Endnutzer-Wert (der Bestand einer aktivierten Tabelle wird
übernommen, [`LH-FA-CAP-009`](../../../../spec/lastenheft.md)); der
Validator-Schritt ist deshalb nicht „n/a“. Die Notizen von
`slice-backfill-spec-nachzug`, `-change-origin`, `-snapshot-reader`,
`-run-usecase`, `-run-store` und `-sql-administration` verschieben die
Validierbarkeit des Bedarfs auf den Wellen-Beleg; dieser Abschnitt ist ihr Träger
und gilt für alle elf Slices (deren Notizen sind Records).

**Feststellung: der Bedarf ist belegt, soweit die Belege unten reichen; der Rest ist
benannt und hat einen Träger.** Ein eigener Lauf der Validator-Rolle in frischem
Kontext fand nicht statt; die Feststellung ist die Belegsammlung des Planners auf den
vorhandenen Läufen und Reports und ersetzt keinen Anwender.

Belegt (Ursprung je Punkt gemessen am 2026-09-25):

- **Bestandsübernahme am laufenden Feed-Container, beide PostgreSQL-Versionen.**
  `docs/user/e2e-abdeckung.md` trägt acht Zeilen mit
  [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (`grep -c 'LH-FA-CAP-009'`
  druckt `8`): Happy Path, Schema-Version, Startposition eines neuen Consumers,
  Boundary, DDL-Fenster, Negative (`docker kill`, `queued`-Aufnahme),
  Leerlauf-Bestätigung und `TestE2EBackfillReplayInvariant`. Die Läufe `36108615045`
  (Push `c82d3333`) und `36116700952` (Push `eb3ba91e`) von `e2e.yml` sind je Leg
  „PostgreSQL 17“ und „PostgreSQL 18“ `success` (`gh run view`).
- **Der Bestand ist lesbar, unterscheidbar und lückenlos angeschlossen.** Die Phase
  „Backfill-Happy-Path“ liest fünf Bestandszeilen als `INSERT` mit
  `origin = 'backfill'` über `cdc.changes` und `GET /changes`, die Replay-Invariante
  vergleicht das Log mit dem Quellstand (Zeilen im Abschnitt Verifikation).
- **Das Benutzerhandbuch-Kapitel „Bestand als Backfill überführen“ ist ausführbar,
  soweit die Reports es gefahren haben.** Der Review von `slice-backfill-sql-administration`
  führte die SQL-Beispiele (Auslösung, Antrags-Abfrage, Status-Abfrage,
  Schlüsselvergleich) gegen eine Wegwerf-Datenbank unter `cdc_admin` und `cdc_reader`
  aus (`review-slice-backfill-sql-administration`, Abschnitt „Handbuch-SQL gegen eine
  Wegwerf-DB ausgeführt“); die dabei gefundene Status-Abfrage unter `cdc_admin` (F-4)
  ist im Handbuch berichtigt (`verifikation-slice-backfill-sql-administration`, F-4).
  Die Startposition und die Sperr-Warteschlange sind im Handbuch mit gemessenem
  Lauf-Ursprung geführt (`verifikation-slice-backfill-e2e`, Zeile 3 und F-1). Weitere
  Beispiele des Kapitels führen die Reports nicht als ausgeführt.
- **Richtgrößen.** `tools/bench-backfill.sh`, Lauf `20260925T084803Z` (gemessen,
  gedruckt): 6.470 Zeilen/s (Median, Stufe 200.000), Spitze 433,0 MiB; die
  Warn-Richtgröße ist gegen diese Zahl bewertet (Abschnitt „Was ging anders als
  geplant?“).

Nicht validiert — je mit Träger:

- **Kein Anwender-Realbetrieb.** Kein Server-Tag trägt den Backfill: der höchste
  Server-Tag ist `v0.1.2`, und `git tag -l 'v*' --contains 510b247e` (erster Commit
  des Backfill-Use-Case) druckt keine Zeile. Der Bedarf beim Anwender bleibt die
  Annahme von [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) selbst; er wird erst mit
  einem Server-Release und der Rückmeldung eines Betreibers geprüft (kein Slice).
- **Breite Zeilen und große Werte, Speicher großer Runs.** Die Bench-Zeilen sind etwa
  74 Bytes breit (Handbuch „Grenzwerte“, übernommen); die Spaltenwerte einschließlich
  `jsonb` sind auf Bild-Parität belegt
  ([`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)), eine
  Last mit breiten Zeilen ist ungemessen. Träger: `slice-backfill-speicher-untersuchung`.
- **Keepalive inmitten einer Transaktion an PostgreSQL 17 an der Quelle und
  „Fehlerschwelle erreicht → Container endet“.** Träger: `slice-capture-leerlauf-quellbelege`.
- **Wiederholung der Klasse `transient` im Erfassungspfad.** Träger:
  `slice-capture-transient-wiederholung`.

**Was den Rest vor einem Server-Release schließt:** `slice-backfill-speicher-untersuchung`
liegt in `done/`, bevor ein Server-Tag entsteht, dessen Commit den Backfill trägt
(Start-Trigger im Plan des Slice, §4). Ein mechanischer Wächter dafür existiert nicht;
die Prüfung liegt beim Planner der Release-Vorbereitung. Die Pläne der drei anderen
Slices nennen keine Bedingung vor einem Server-Release (`grep -n -i release` über
sie druckt einen Treffer: den Namen eines Testskripts in
`slice-harness-suchlauf-nachmessen`).

## Übergabe an den Release-Zug

Die drei SDK-Packages tragen `0.2.0` (`<Version>` in
`sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`, `version` in
`sdks/python/pgchangefeed/pyproject.toml`, `version` in
`sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`, gemessen am Stand `32028d4f`).
Die drei Tags `sdk-csharp-v0.2.0`, `sdk-python-v0.2.0` und `sdk-kotlin-v0.2.0` stehen
an `eb3ba91e` (`git tag -l`, gemessen am 2026-09-25); der Tag-Push ist
Betreiber-Handlung und nicht Teil der Welle. Die drei Publish-Läufe
(`36117929191` C#, `36117929298` Python, `36117929552` Kotlin, je `event` `push`
auf dem Tag) sind `success` (`gh run view`, gemessen am 2026-09-25) — der reale
Post-Push-Beleg nach [`AGENTS.md`](../../../../AGENTS.md) §3.10; das Erscheinen der
Packages in NuGet.org, PyPI und GitHub Packages ist nicht gemessen.
Beleg-Anker der Versionsentscheidung: §3 „Versionsentscheidung“
des Plans `slice-backfill-sdk-origin`. Ein Server-Release ist nicht Teil der Welle
(`docs/user/version.md` bleibt unberührt).

## Verifikation

Alle Zeilen am Stand `32028d4f` (Docker-Läufe dieser Closure, 2026-09-25), sofern
nicht anders angegeben; Gate-Exit-Codes ungefiltert gesichert
([`AGENTS.md`](../../../../AGENTS.md) §3.9).

### Schritt 1 — Trigger

| Kriterium (Welle-Datei §3) | Beleg |
|---|---|
| Alle elf Slices in `done/` | `ls docs/plan/planning/done \| grep -c '^slice-backfill-'` druckt `11`; die elf Namen stehen im Abschnitt „Was wurde geliefert?“ |
| `make gates` grün | Exit `0`; `baseline-verify: v6.9.0 OK — 54 Dateien`, `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%`, `d-check: 1137 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`, `generated-sync: OK`, a-check `gesamt: 0 Befund(e)`; ein zweiter Lauf am Arbeitsbaum dieser Closure (Exit `0`) druckt `coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%` und `d-check: 1138 Datei(en) geprüft, 0 Befund(e)` — die gedeckte Zahl streut zwischen Läufen ohne Code-Änderung |
| Realer, grüner `make test-integration`-Lauf mit den Backfill-Belegen | Lauf `36108615045` von `e2e.yml` (Push auf `c82d3333`): beide Legs „image + test-integration (PostgreSQL 17)“ und „(PostgreSQL 18)“ `success`; gedruckt je Leg die Zeilen „Backfill-Happy-Path“ (5 Bestandszeilen, 6 gelesene Changes, davon 1 mit `origin=wal`), „Backfill-Schema-Version-Beleg“, „Backfill-Startposition“, „Backfill-Boundary (leere Tabelle, zweiter Antrag)“, `TestE2EBackfillReplayInvariant` (PG 17: 55 Backfill-Changes, WAL-Changes davor 412 und dahinter 131, 54 Zeilen im Quellstand; PG 18: 56, 406, 166, 56), „Backfill-DDL-Fenster“, „Backfill-Negative (docker kill, queued-Aufnahme)“, „Leerlauf-Bestätigung“. Zwischen `c82d3333` und dem Stand dieser Closure ändert kein Commit `internal/`, `cmd/`, `test/`, `proto/`, `Dockerfile`, `compose.yaml`, `go.mod`, `tools/schema/`, `sdks/`, `examples/` oder `.github/` (`git diff --name-only c82d3333..HEAD -- <Pfade> \| wc -l` druckt `0`); ein lokaler Wiederholungslauf ist nicht gefahren, weil das Kriterium einen realen, grünen Lauf verlangt und dieser vorliegt |
| CI-Stand des Push | `gh run list --commit c82d3333…`: `ci` (`36108615036`), `examples` (`36108615027`), `e2e` (`36108615045`) je `success` |
| `make test-replication` grün | Exit `0` an PostgreSQL 18 (Pin von `PG_TEST_IMAGE`) und an PostgreSQL 17 (`postgres:17-alpine@sha256:7456ef82…`, Digest aus `e2e.yml`); beide Phasen (`measure`, `tier`) laufen durch |
| `make test-store` grün | Exit `0` an PostgreSQL 18 und an PostgreSQL 17; Teilzahl `DB-Adapter-Coverage (Teilzahl, nur store-Bestand): 78.72% (gedeckt 540 von 686 Statements; Profil: store)` |
| `make schema-rollout` zweimal und Alt-Tag-Lauf | `bash tools/harness/run-schema-rollout-guard-test.sh` Exit `0`; Lauf 2 fährt den zweiten Rollout gegen dasselbe migrierte Ziel; gedruckt: `Lauf 5 OK — Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; 19 Tabellen-/View-Rechte der drei Rollen …` und `run-schema-rollout-guard-test: OK — alle Belege real erbracht (…)` |
| `make doc-trace` | Exit `0`; die Zeile `LH-FA-CAP-009 … E2E … ok`, Schluss `80 Anforderung(en), 2 Waise(n).` |
| Kein Ende des Capture-Prozesses über den WAL-Rückstand | `bash tools/bench-backfill.sh` (einzeln, nach `make image`, Exit `0`), Lauf `20260925T084803Z`: `WAL-Rückstand der größten Stufe (200000 Zeilen) — höchste Spitze im Run 0 MiB, unter der Warnschwelle von 100 MiB (SPEC-013)`; die Ursache (Slot bestätigt im Leerlauf, [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)) belegt die E2E-Phase „Leerlauf-Bestätigung“ des Laufs oben |
| Row-Image-Konstruktion an genau einer Stelle | `grep -rn -E 'json\.(Marshal\|NewEncoder)' internal --include=*.go` (ohne Testdateien): die Bild-Erzeugung `json.Marshal(column)` und `json.Marshal(*values[i])` steht in `internal/domain/model/rowimage.go` (`BuildRowImage`); die übrigen Treffer sind Wire-Nachrichten (`http/sse.go`, `natsstream/publisher.go`) und HTTP-Antworten. Aufrufer: `mapper.go` (WAL-Pfad, zwei Aufrufe für `new` und `old`) und `usecase/backfill/service.go` (Backfill-Pfad) — `grep -rn 'BuildRowImage' internal --include=*.go` |
| Handbuch-Träger | Änderungshistorie bis Version 1.59; „Bestand als Backfill überführen“ trägt die Sichtbarkeits-Grenze (Position `X`, Consumer hinter `X`), „Startposition eines neuen Consumers“ (mit dem gemessenen Lauf-Beleg), „Lesen“ (ohne `LIMIT`), Warnungen, WAL-Rückstand; „Grenzwerte“ trägt Toleranz („Startwert, Setzung ohne Messung“), Richtgröße (mit Ursprung je Lauf) und die Grenze der Ein-Transaktions-Form (gehaltenes WAL, Spill) |
| Closure-Notiz | diese Datei |

### Schritt 2 — Trigger-Audit

- **Carveouts: 0 offen.** `ls -la docs/plan/carveouts` zeigt genau die Datei
  `.gitkeep`.
- **Bootstrap-aware Gates.**
  - *Unit-Gate:* `harness/mk/coverage.mk` führt `THRESHOLD ?= 80` (Endstufe), der
    Gate-Lauf druckt `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%`;
    Rampe ausgeschöpft, **0 offen**.
  - *DB-Adapter-Coverage:* der Hochschalt-Trigger war fällig
    (`DB_COVERAGE_THRESHOLD` stand auf 70). Die gemergten Läufe von `make test-store`
    und `make test-replication` (je eigenes `DB_COVERAGE_DIR`) druckten an
    PostgreSQL 18 **und** an PostgreSQL 17 je `DB-Adapter-Coverage: 82.51% (gedeckt
    873 von 1058 Statements; Profile gemergt: store,replication)`; der CI-Lauf
    `36108615045` druckte dieselbe Zeile in beiden Legs. Die Bewertung „Ausbau, nicht
    Verdünnung“ steht in `architect-verdict-welle-backfill-bestand-lese-schritt` §6.
    Die Stufe steht auf 80 (Endstufe, `tools/harness/db-coverage.sh`); der Lauf
    endet bei Schwelle 80 mit `db-coverage: OK — DB-Adapter-Coverage 82.51% erfuellt
    Schwelle 80%` (Exit 0), bei `DB_COVERAGE_THRESHOLD=85` mit `db-coverage: FAIL —
    DB-Adapter-Coverage 82.51% unter Schwelle 85%` (Exit 1). Eine Hebung ist keine
    Senkung ([`AGENTS.md`](../../../../AGENTS.md) §3.6); eine ADR ist nicht nötig. Die
    Wirkung auf den Workflow `e2e.yml` (dessen Schritt „DB-Adapter-Coverage —
    Replication-Teil, Merge + Schwelle“ trägt die Stufe) ist ein Wert, kein
    strukturelles Workflow-Element; der Lauf des ersten Push mit Stufe 80 belegt sie
    im CI (Beleg im Abschnitt „Nachlauf“ unten). Die gedeckte Zahl streut von Lauf zu Lauf.
- **ADR-Re-Evaluierungs-Trigger.**
  - [`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md): der Trigger
    „wiederholt Rot ohne reale Regression (Rauschen)“ ist nicht eingetreten — die
    Ursache ist ein gemessener, systematischer Host-Befund; der Trigger „zweite
    Umgebung führt `make bench` aus“ ist nicht eingetreten
    (`grep -rl 'make bench' .github/workflows \| wc -l` druckt `0`). **0 offen.**
  - [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md):
    „Eine Messung liegt vor“ ist eingetreten und verbraucht (die Messung trägt die
    Richtgröße; die Toleranz bleibt Startwert). **0 offen.**
  - [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md): Betriebs-
    und Versions-Ereignisse (Kopierdauer über der Toleranz, Live-Weg, Auslösung ohne
    SQL, `EXPORT_SNAPSHOT`) sind nicht eingetreten; `make pin-stale-dmigrate`
    (Netz, advisory) druckt
    `DRIFT       D_MIGRATE_IMAGE (ghcr.io/pt9912/d-migrate:latest): gepinnt sha256:862dfb04…, aktuell sha256:af9d3eb3…`
    (make-Exit `2`): der Tag `:latest` trägt einen anderen Bau als der gepinnte
    Digest. Der Pin bleibt, ein Digest-Wechsel ist ein bewusster Commit
    ([`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7);
    ob die aktuelle d-migrate-Version CHECK-Änderungen konvergiert (Trigger „d-migrate
    konvergiert CHECK-Änderungen“), ist aus dem Digest nicht ablesbar — Anker:
    `BEO-PGC/d-migrate-nacharbeit` (7×, Klassen Funktion und CHECK offen). Kein
    Blocker.
  - [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) bis
    [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md):
    kein Trigger eingetreten (gelesen im Verdikt §6); `ADR-0121`s Trigger „Ein
    committeter Test der Quellseite entsteht“ tritt mit
    `slice-capture-leerlauf-quellbelege` ein.
    [`ADR-0122`](../../adr/0122-backfill-replay-invariante-e2e-tier.md): 0 offen.

### Schritt 4 — Archivierung

Das Repo führt kein Archivierungs-Werkzeug: kein Make-Ziel in `Makefile` oder
`harness/mk/*.mk` und kein Skript unter `tools/` (`grep -rn archiv harness Makefile
tools` findet keinen Aufruf). Das externe Werkzeug `ai-harness-init archive-welle`
liegt lokal (`.harness/state/bin/`, nicht im Repo, `.gitignore`); sein Vorschau-Lauf
(`--vorschau welle-backfill-bestand`, schreibt nichts) meldet 11 Mitglieder, 1 wellenlosen
Slice, 40 fremde und zwei Sperren (Ergebnisnotiz und Welle-Plan in `done/`), die mit
Schritt 3 und Schritt 5 entfallen. Der schreibende Lauf ist nicht Teil dieser Closure:
er committet mit festen Messages ohne `LH-*`/`ADR-*`-Kennung
(`BEO-PGC/externes-werkzeug-committet-ohne-kennung`), und die Closure von
`welle-nats-drittstream` hält dieselbe Bedingung als nicht eingetreten fest. Die
Bedingung des Schrittes (Werkzeug im Repo) ist nicht eingetreten; keine
Handarchivierung. Die Entscheidung, den Lauf des externen Werkzeugs für diese Welle
gleichwohl auszuführen, liegt beim Betreiber.

### Prüfung in frischem Kontext

Die Prüfung dieser Closure-Notiz durch die Rolle in frischem Kontext
(`.harness/skills/closure-note-reviewer.md`) liegt als
`review-closure-notes-welle-backfill-bestand` vor; der Planner prüft sie nicht selbst.
Der Ausgang je Finding steht unter „Ausgang des Closure-Note-Reviews“.

### Paarungen (Modul 6)

- **(a) Anker.** Jede Regel trägt an ihrem Zielort `seit welle-backfill-bestand`:
  `AGENTS.md` (§3.12 Absatz „Verfasser einer ADR“ und §3.13 Suchform),
  `.harness/skills/reviewer.md` (MEDIUM-Punkt und die zwei Klauseln des
  Kommentar-HIGH), `harness/sensors/db-adapter-coverage.md` §Gegenstand,
  `docs/plan/planning/observations/README.md` (Deckel) — je
  `grep -n -B1 'welle-backfill-bestand' <Datei>` findet die Zeile; der Zeiger in
  `.claude/agents/architect.md` verweist auf `AGENTS.md` §3.12.
- **(b) Folge-Slice.** Die vier Slices stehen als Datei in `docs/plan/planning/open/`
  (`ls` nennt sie).
- **(c) Register.** Jede in dieser Notiz genannte Kennung `BEO-PGC/<slug>` besitzt
  ein Verzeichnis mit nicht leerem `evidence/` — mit einer Ausnahme:
  `limit-fortsetzung-innerhalb-einer-position` (offen, „benannt, nicht gezählt“, in
  §Beobachtungs-Register als Adresse genannt) trägt kein `evidence/`; im Register
  tragen zwei von 113 Verzeichnissen keins (`architect-verdikt-ablageort-uneinheitlich`,
  gestrichen, und diese; `ls` je Verzeichnis).

### Docker-Lauf-Hygiene

Die Docker-Läufe dieser Closure (`make gates`, `make image`, `make test-store` und
`make test-replication` je an PostgreSQL 18 und 17, der Guard-Test, `bench-backfill`)
liefen nacheinander, nicht gleichzeitig. Die Zahl der `dangling` Docker-Volumes ist
vor und nach den Läufen `34` (`docker volume ls -q -f dangling=true \| wc -l`); die
Läufe hinterlassen keine Container.

### Nachlauf

- Der erste CI-Lauf von `e2e.yml` nach dem Push dieser Closure trägt
  `DB_COVERAGE_THRESHOLD` 80 im Schritt „DB-Adapter-Coverage — Replication-Teil,
  Merge + Schwelle“ (beide Legs). Beleg ([`AGENTS.md`](../../../../AGENTS.md) §3.10,
  gemessen am 2026-09-25): der Lauf `36116700952` (Push `eb3ba91e`) ist `success`,
  beide Legs; der Schritt druckt je Leg `DB-Adapter-Coverage: 82.51% (gedeckt 873
  von 1058 Statements; Profile gemergt: store,replication)` und `db-coverage: OK —
  DB-Adapter-Coverage 82.51% erfuellt Schwelle 80%` (`gh run view 36116700952 --log`).
  Die Workflows `ci` (`36116701060`) und `examples` (`36116700976`) desselben Push
  sind `success`.

### Ausgang des Closure-Note-Reviews

Der Review in frischem Kontext liegt vor (`review-closure-notes-welle-backfill-bestand`:
0 HIGH, 1 MEDIUM, 5 LOW, 3 INFO). Jedes Finding hat einen Ausgang; die elf
Slice-Notizen sind Records und bleiben unverändert, diese Notiz trägt die
Feststellungen für alle.

| Finding | Ausgang | Träger |
|---|---|---|
| F-1 Validator-Schritt ohne Träger | erfüllt: explizite Feststellung „belegt so weit, Rest benannt“ | Abschnitt „Validator-Feststellung (Modul 8)“ oben |
| F-2 zwei gemeldete Ungenauigkeiten ohne Träger | Kommentar `allowDestructive` in `tools/schema/rolloutguard/guard.go` → Plan `slice-transformationen-antragsweg-schema` (ändert dieselbe Datei, §3-Zeile); Kommentar an der gRPC-Nachricht in `proto/cdc/stream/v1/changestream.proto` → Register-Eintrag, kein Slice passt (kein offener Slice ändert die `.proto`) | `BEO-PGC/gemeldete-ungenauigkeit-ohne-traeger` (offen, 1×) |
| F-3 Spec-Lücke „aufgelöst“ | berichtigt: offen, Adresse benannt | „Benannte Spec-Lücke“ unter Steering-Loop-Einträge |
| F-4 Streichung des Decoder-Eintrags | Formulierung präzisiert: Quelltext-Suche, Restgrenze genannt | Abschnitt „Einträge unter 3×“ oben, `state.md` des Eintrags |
| F-5 drei Kandidaten ohne Ausgang | zwei gestrichen, eines akzeptiertes Negativ | `state.md` von `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` |
| F-6 Lerneintrag im Sammelabsatz | keine Änderung an den Records; diese Notiz führt die Einträge einzeln | Einleitung der Steering-Loop-Einträge |
| F-7 Kosten der Ausschluss-Lesung je Block | keine Aktion: der Bench fährt den Use Case im Feed-Container und misst `finished_at − started_at` des Runs, die Lesung ist in der Kopierdauer enthalten (Stufe 200.000 Zeilen bei Blockgröße 1.000: 200 Blöcke, abgeleitet); die 1.000.000-Stufe liest `slice-backfill-speicher-untersuchung` mit | Plan `slice-backfill-speicher-untersuchung` §1 |
| F-8 Speicher-Zahlen enger als das Handbuch | berichtigt: Untergrenze, höchster Wert 641,7 MiB, zwei Runs mit 435 und 1.544 MiB | Abschnitt „Was ging anders als geplant?“ |
| F-9 Zeitformen in Records | der Post-Push-Beleg steht im Abschnitt „Nachlauf“; die Stand-Aussagen der Slice-Notizen sind Records | Abschnitt „Nachlauf“ |
