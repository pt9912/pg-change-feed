# Verifikations-Report: slice-backfill-speicher-untersuchung — 2026-09-25

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich, Entscheidungs-Konformität, Plan-vs-Code-Diff,
Gates und Nachmessung der Fixrunde. Review-Artefakt des Reviewers:
[`review-slice-backfill-speicher-untersuchung.md`](review-slice-backfill-speicher-untersuchung.md); Formvorbild dieses
Reports: [`verifikation-slice-sdk-readme-nutzerdoku.md`](verifikation-slice-sdk-readme-nutzerdoku.md). Die vendored
Baseline trägt kein eigenes Verifikations-Template (`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur
`review-report.template.md`), ein Skill `.harness/skills/verifier.md` liegt nicht vor; der Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-backfill-speicher-untersuchung` (ohne Welle, Lifecycle `in-progress`), Stand `HEAD` =
`aff5f418` (gepusht), Diff-Range `493a28ad..HEAD` (17 Dateien, +3434/−51). Slice-Inhalt sind `fb04c86b` (Werkzeuge),
`b9cf9729` (Messberichte), `7459afbb` (Handbuch), `5a141059` (Folge-Slice), `28768e7e` und `4f94f900` (Plan-Nachzug),
`aff5f418` (Fixrunde) sowie die reinen `git mv`-Commits `446e3ca6` und `37e825e2`. Nicht Slice-Inhalt sind `ec7f76e2`
(`slice-sdk-kotlin-cloudsmith`), `3ec90077` (Architect-Verdikt und ADR), `5b1f7762` (Review-Report) und `3405c9a1`
(Plan-Nachzug des Folge-Slice). Bezug: [`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md).
Dieser Lauf ändert weder Code noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen liefen als
Kopie bzw. im Arbeitsbaum und sind zurückgenommen (`git status --short` leer); nichts gepusht, nichts getaggt.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei des Scratchpads; der Exit-Code wurde gesondert gesichert und danach gelesen (der
Exit-Wert eines Hintergrund-Wrappers ist nicht der von `make`, gemessen wurde die geschriebene Exit-Datei). Schwere
Docker-Läufe liefen nacheinander; `free -m` (verfügbar) vor den Läufen 13,1 bis 13,4 GB.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` (Stand `aff5f418`) | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%` · `d-check: 1165 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make docs-check` | **EXIT=0** | `d-check: 1165 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=493a28ad..HEAD` | **EXIT=0** | `d-check: 1165 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=493a28ad..HEAD` | **EXIT=0** | `d-check: 1165 Datei(en) geprüft, 0 Befund(e)` (der Aufruf ohne `RANGE` endet mit `flag needs an argument: --range`, Aufrufform, kein Befund) |
| `bash -n tools/bench-*.sh` (sechs Skripte) | je Exit 0 | `bench-backfill-memory.sh`, `bench-backfill.sh`, `bench-batch-vs-single.sh`, `bench-lib.sh`, `bench-scaling.sh`, `bench-source-impact.sh` |
| `git diff --stat 493a28ad..HEAD -- internal cmd proto gen tools/schema` | leer | kein Produktionscode im Diff; die 17 Dateien liegen in `docs/`, `harness/`, `tools/` |

Hygiene: dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**, nach allen Läufen
**34**; kein `prune`; kein Container und kein Netz `pgc-bench-*` zurück; die zwei laufenden Container `waza:local`
gehören nicht zu diesem Lauf; `tools/schema/plan.yaml` und `tools/schema/down.sql` nach jedem Lauf per `git checkout`
zurückgenommen; die Mutations-Kopie eines Skripts entfernt.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Messreihe mit gedruckten Zeilen; Stufen bis 1.000.000 (`--full`); Host, Lauf, Zeile je Zahl; Messung von Übernahme getrennt | **erfüllt** | Zeilen-Dokument enthält die Reihen A bis K, L1, L2, M, N, O und P1 bis R (Überschriften gezählt); Reihe N ist `tools/bench-backfill.sh --full` (Stufen 10.000, 100.000, 1.000.000, drei Runs); die Kopfzeile jeder Reihe nennt Host und Image-ID; übernommen ist nur 1.544 MiB (Lauf `20260924T233628Z`), daneben die vereinbare Messung 1.538,7 MiB (Reihe B, Run 2, Zeile „memory.peak … 1538.7 MiB“) |
| 2 | Ursache benannt und belegt (Schalter an/aus) | **erfüllt** | Code an `v0.1.0`, `v0.1.2` und `HEAD` gelesen (§5); Schalter „Bereinigung aus“ an den Reihen E, F, G, I, J nachgerechnet: `anon` im Run 8,9 bis 10,6 MiB (n = 12, E und J), `memory.peak` J 13,3/13,2/13,7 MiB bei 1.000.000/2.000.000/3.000.000 Changes gegen 1.082,7 MiB (Reihe B, Run 1); der Schalter „Bereinigung aus/an“ ist vom Reviewer an eigenen Läufen reproduziert (übernommen aus dem Review-Report: 13,0 / 13,4 / 12,3 gegen 135,4 / 242,2 / 352,7 MiB), von mir an den gedruckten Zeilen nachgerechnet, nicht neu gefahren. Das Ausbleiben der Takte ab 2.000.000 Changes bleibt als *nicht belegt* im Bericht und im Handbuch |
| 3 | Ausgang gesetzt (Änderungs-Slice in `open/` oder Handbuch-Grenze; Richtgröße bewertet; Version und Historie) | **erfüllt** | Datei `docs/plan/planning/open/slice-retention-lauf-speicher-begrenzung.md` existiert; Handbuch `Version: 1.61`, Zeilen 1.60 und 1.61 in `### Änderungshistorie`; Richtgröße im Messbericht §6 bewertet (7.715 × 600 = 4.629.000, nachgerechnet); `make docs-check` Exit 0 |
| 4 | `make gates` grün, Exit ungefiltert gesichert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, kein offenes HIGH/MEDIUM nach der Fixrunde | **erfüllt** | Review-Report vorhanden (2 HIGH, 3 MEDIUM, 5 LOW, 4 INFO); F-1 bis F-5 am Zeilen-Dokument nachgemessen behoben (§3); die Fixrunde hat keinen zweiten Reviewer-Durchgang (V-5) |
| 6 | §3.13-Suchlauf, Gefundenes und Nichtgefundenes, beide Stände | **erfüllt, mit V-1, V-3** | Zahlen an beiden Ständen nachgefahren (§6): alle Trefferzahlen der Fixrunde stimmen; die Aussage „Nicht gefunden“ der Fixrunde trägt nicht (V-1); das erste Zählfeld ist am `HEAD` veraltet (V-3) |
| 7 | Doku-Update (Handbuch §Grenzwerte, Vertrag `harness/targets/bench-backfill.md`, Grundlinie) | **erfüllt** | Handbuch §9, §4, „Aufbewahrung (Retention)“ gelesen (§7); Vertrag führt Grundlinie, zweites Skript, Kill-Zusage und cgroup-Vorbedingung; `harness/README.md` nennt das zweite Skript |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ |
| 9 | Reconciliation-Register | **korrekt offen** | Zeile „entfällt“ trägt die Begründung; das Häkchen setzt der Planner |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht; `state.md` des Eintrags trägt noch „geplant“ (Planner) |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | §6 trägt fünf Risiken mit „*(bei Closure)*“ |
| 12 | Drei Paarungen | **korrekt offen** | hängen an der nächsten Welle-Closure |

`[x]` sind sieben Zeilen, `[ ]` fünf, zusammen zwölf (Zeilen des Plans gezählt). Kein `[x]` ohne Beleg; kein `[ ]`, das
über die Rollen-Sequenz hinaus belegt wäre.

## 3. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) `GOGC=25` | Reihe K gegen Reihe A aus dem Zeilen-Dokument, `awk`: 224,5 gegen 273,9 MiB = −18,0 %, 415,5 gegen 468,6 = −11,3 %, 477,4 gegen 605,9 = −21,2 %; Handbuch („11 bis 21 %“, drei Paare) und Messbericht §3.6 und §4 tragen genau diese Spanne | **behoben** |
| F-2 (HIGH) Höchstwert je Change | 3.096,6 × 1.024 / 2.000.000 = 1,585 (B, Run 3) und 2.751,7 → 1,409 (C, Run 3); `n` = 15 Werte gezählt, Minimum 1,03 (A, 600.000), Maximum 1,59; Handbuch „1,03 bis 1,59 KiB“ (zwei Treffer in `docs/user harness`), die Runs 3 stehen in einer eigenen Zeile „2.000.000 vor dem Run, 3.000.000 danach (Spitze im Run)“; Folgezahlen 4.000.000 × 1,03/1,59/2,65/4,19 KiB = 3,93/6,07/10,11/15,98 GiB nachgerechnet | **behoben** |
| F-3 (MEDIUM) Handbuch §4 gegen §9 | §4-Absatz „Speicher des Feed-Containers“ nennt den flachen Speicher nur bei **leerem** `cdc.change` und die Spitzen im Run bei nicht leerem `cdc.change` (142,1 MiB bei 100.000 Changes davor, Reihe A, Stufe 100.000, Run 2; 1.460,8 MiB bei 1.000.000 davor, Reihe B, Run 2 — beide Zeilen gelesen, `anon` Spitze); verweist auf „Grenzwerte“; „nach dem Run, nicht während er läuft“ steht nicht mehr im Baum (§6, Befehl 2) | **behoben** |
| F-4 (MEDIUM) Kill in `bench-backfill.sh` | eigener Lauf (§4): Exit 1 mit „Feed-Container läuft nicht mehr … OOMKilled true“; Mutation ohne Prüfung: „nicht beendet nach 60 s“ ohne Zustand | **behoben** |
| F-5 (MEDIUM) Trigger gegen Empfehlung | Folge-Slice §4 trägt „Vorbedingung des Server-Release `v0.2.0`“ mit `git tag -l 'v*'` als beobachtbarem Anker; Messbericht §7 „Ausgang“ trägt Reichweite `v0.1.0` bis `v0.1.2` und dieselbe Reihenfolge; beide Stellen widersprechen sich nicht | **behoben** |
| F-6 (LOW) `gc_summary` | die Funktion mit einer Attrappe von `docker` gefahren: fünf Werte von `BENCH_FEED_ENV` drucken exakt die Zeilen der Reihe R („gesetzt“ nur bei den zwei Werten mit `GODEBUG=gctrace=1`, `XGODEBUG=gctrace=1` zählt nicht); mit GC-Zeile Zusammenfassung | **behoben** |
| F-7 (LOW) Chronik im Handbuch | Suche in den hinzugefügten Handbuch-Zeilen nach früher, bisher, jetzt, zuvor, wurde, „bis dahin“: keine Prosa-Treffer (nur Historienzeilen und Dateinamen in Linkzielen); keine `ADR-`/`LH-`/`SPEC-`/`slice-`-Kennung im Fließtext | **behoben** |
| F-8 (LOW) Ursache im Plan des Folge-Slice | Folge-Slice §6 nennt „ab 2.000.000 Changes“, „beobachtet und nicht erklärt“, „Hypothese, nicht gemessen“ und die Gegenprobe Reihe N | **behoben** (Plan-Nachzug `3405c9a1`) |
| F-9 (LOW) Ursprung der Live-Aussage | Handbuch: „aus dem Code gelesen, nicht gemessen; gemessen ist der Backfill als Quelle“; Messbericht §1 und §9 gleich | **behoben** |
| F-10 (LOW) Fehlerpfad-Reste | (a) eigener Lauf mit `BENCH_MEM_RUN_TIMEOUT_S=1` und leerem `TMPDIR`: Exit 1, `TMPDIR` danach leer; (b) eigener Lauf mit `--memory 64m`: Exit 1, „Feed-Container nicht mehr lesbar … OOMKilled true“, **keine** Zeile „Zeile 58 … Datei oder Verzeichnis nicht gefunden“ | (a) **behoben**, (b) nicht reproduzierbar, Nicht-Änderung plausibel |
| F-11 (INFO) Gegenprobe Reihe N | Messbericht §3.3 „Gegenprobe (Reihe N)“ und Handbuch §9 nennen 2.330.000 Changes davor, 2.448,4 MiB nach 60 s (Zeile der Reihe N gelesen) und den unerklärten Unterschied | **behoben** |
| F-12 (INFO) Tag-Rahmen | siehe §5 | **behoben** |
| F-13 (INFO) Vorbedingung im Vertrag | Vertrag nennt cgroup v2 mit systemd-Treiber, den Pfad und Exit 1 mit der Meldung; die Exit-Codes der Mutationen stehen im Messbericht §8 (P2, P3, Reihe L) | **behoben** |
| F-14 (INFO) Prozess-Hinweis | keine Aktion vorgesehen; `git status` sauber, `bash -n` ohne Meldung | **kein Rückstand** |

## 4. Mutationen und Nachläufe der Eingabeseite (dieser Lauf)

| # | Lauf | Ergebnis (gedruckt) |
|---|---|---|
| A | `tools/bench-backfill.sh` mit `BENCH_FEED_DOCKER_ARGS="--memory 64m"`, Stufe 100.000, 3 Runs (Lauf `20260925T171638Z`) | **Exit 1**; Zeile „Feed-Container läuft nicht mehr, Run … von bench_backfill_s100000 nicht beendet (Run-Status running): Status exited, Exit 137, OOMKilled true“; vorher Run 1 mit 11,0 MiB und Run 2 mit 63,8 MiB Spitze |
| B | dieselbe Eingabe, Kopie des Skripts ohne die Zeile 181 (Prüfung in `run_backfill`), `BENCH_BACKFILL_RUN_TIMEOUT_S=60` (Lauf `20260925T171900Z`) | **Exit 1**, aber „Run … von bench_backfill_s100000 nicht beendet nach 60 s“ — ohne Zustand des Containers: die Prüfung ist es, die den Kill benennt (Rot der Mutation gesehen) |
| C | `tools/bench-backfill-memory.sh`, `BENCH_MEM_RUN_TIMEOUT_S=1`, leeres `TMPDIR` (Lauf `20260925T172158Z`) | Exit 1, „Run … nicht beendet nach 1 s“; `ls -A` des `TMPDIR` danach leer |
| D | `tools/bench-backfill-memory.sh` mit `--memory 64m` und `GODEBUG=gctrace=1`, 100.000 Zeilen, 1 Run (Lauf `20260925T172325Z`) | Exit 1, „Feed-Container nicht mehr lesbar: Status exited, Exit 137, OOMKilled true“; keine Zeile „Zeile 58“ |
| E | `gc_summary` gegen eine Attrappe von `docker`, fünf Werte von `BENCH_FEED_ENV` | Ausgabe wie Reihe R |

Nicht gefahren: die Live-Phase-Mutation (die Prüfung in Zeile 322 entfernen, Reihe P3) — sie läuft ohne Ende (Reihe P1:
mehr als 600 s ohne Zeile) und verlangt einen Abbruch von Hand; die Zeile der Reihe P3 wird gelesen, nicht nachgefahren.
Nicht gefahren: die Reihen A bis K, N (Stunden Laufzeit); `make bench` als Ganzes (Auftrag).

## 5. Reichweite: `v0.1.0` bis `v0.1.2`

`git grep -n ReadChanges <Tag> -- internal/application/usecase/retention` (ohne Tests): `v0.1.0`, `v0.1.1`, `v0.1.2`
und `HEAD` tragen je Zeile 63 `s.store.ReadChanges(ctx, outbound.ChangeQuery{Source: command.Source})`;
`retentionInterval = 10 * time.Second` steht in `internal/bootstrap/wiring.go` Zeile 172 der drei Tags, der Aufruf
`runRetentionCleanup(…, retentionInterval, …)` in Zeile 779; `SelectChanges` trägt in `v0.1.0` und `v0.1.2` `old_data`,
`new_data` und `LIMIT $6`; `git diff --stat v0.1.2 HEAD -- internal/application/usecase/retention
internal/application/port/outbound/changestore.go` ist **leer**. Die Aussage „v0.1.0 bis v0.1.2 tragen denselben
Defekt“ trägt (Quelltext verglichen; die Handbuch-Formulierung „nicht am Image der Version gemessen“ ist ehrlich).

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

| Plan-Zeile | Meine Messung (`git grep`, Parent `493a28ad` bzw. `5b1f7762`, Diff = `HEAD`) | Ergebnis |
|---|---|---|
| Befehl 1 (erstes Feld) | Parent **242** (`docs/user` 58, `harness` 13, `spec` 36, `internal` 111, `tools` 24); Stand `4f94f900` **280** (71 in `docs/user`); `HEAD` **288** (`docs/user` 79, `harness` 23, `spec` 36, `internal` 111, `tools` 39) | Parent bestätigt; „Diff 280 (71 …)“ gilt für `4f94f900`, nicht für `HEAD` (V-3) |
| Befehl 2 (Symbolnamen) | Parent **11**, `HEAD` **11** | bestätigt |
| Alte Zahlen / „nicht untersucht“ | Parent **7** und **1**, `HEAD` **2** und **0** | bestätigt |
| Fixrunde Befehl 1 bis 5 | Parent/Diff **16/7**, **2/1**, **2/0**, **17/27**, **8/13** | alle fünf bestätigt |
| Die sieben Treffer von Fixrunde-Befehl 1 im Diff | zwei Tabellenzeilen des Messberichts (1,51 aus Reihe D, 1,57 aus Reihe C Run 2), zwei Treffer „16,0 GiB“, Versionszeile 1.57, Historienzeile 1.61, Zeile 293 des Folge-Slice | bestätigt; alle wahr |
| Zeile 293 des Folge-Slice-Plans | trägt `1,03 bis 1,57`; `git grep '1,03 bis 1,57' HEAD -- docs/user harness` **0** Treffer, `'1,03 bis 1,59'` **2** | Meldung des Implementers bestätigt (V-2) |
| „Nicht gefunden: kein weiterer Träger außerhalb der Suchwurzeln nennt den Höchstwert 1,57 …“ | `git grep -n -E '1,03 bis 1,57' HEAD` über den ganzen Baum: `docs/plan/adr/0124-…md` Zeilen 68, 191, 199; Record `architect-verdict-retention-lauf-speicher-begrenzung.md` Zeilen 152, 155, 167; Review-Report Zeilen 112, 113, 165 | Aussage **falsch** für die ADR und das Verdikt (V-1) |

## 7. Träger, Zahlen, Sprache

- **Handbuch** (Abschnitte §9 „Grenzwerte“, §4 „Bestand als Backfill überführen“, „Aufbewahrung (Retention)“): Zahlen der
  Tabelle nachgerechnet und gegen das Zeilen-Dokument gehalten (23,9; 133,8 bis 151,3 mit n = 4; 245,0 und 273,9; 605,9;
  1.082,7 und 1.273,5; 2.269,2 und 3.058,0; 3.096,6 und 2.751,7 MiB — alle wie gedruckt); breite Zeilen 361,0 bis 1.550,6 MiB,
  n = 6; „31 GiB RAM“ (33.362.599.936 Byte / 2^30 = 31,07); Bemessung 2 KiB (schmal) und 4,5 KiB (breit) plus 64 MiB liegt
  über 1,59 und 4,19 KiB; 4.000.000 × 2 KiB + 64 MiB = 7,69 GiB („etwa 7,7 GiB“); `--memory 64m` gemessen (Reihen L1, L2,
  P2, P3 und mein Lauf A).
- **Messbericht** (§1 bis §9): Tabellen §3.3 (alle 15 Werte je Change, Ableitung Spitze × 1.024 / Changes), §3.4
  (Blockgröße: Median 8,5 / 9,9 / 24,4 MiB bei 100 / 1.000 / 10.000; rund 1,65 KiB je zusätzliche Zeile), §3.5 (Reihe H:
  5,07; 3,70; 3,71 und 4,19; 3,53; 2,71; 2,65), §3.7 (Reihe D: 55,7 / 84,1 / 82,9 % nach 120 s), §3.8 (Reihe N: Ruhe 6,1 / 41,7 /
  275,7 MiB, Spitzen 25,5 / 231,9 / 2.738,2, 20 s / 60 s 349,6 / 243,6 und 2.737,2 / 2.448,4, Changes 30.000 / 330.000 /
  3.330.000, Kopierdauer 109.110 bis 133.442 ms), Grundlinie: 19 Einheiten, `anon` 5,0 bis 5,5 MiB, `memory.current`
  15 × 5,8 bis 6,8 MiB und 4 × 18,9 bis 22,4 MiB (aus den Zeilen extrahiert und gezählt). Kein Zahlenbefund; Herkunft
  (gemessen / abgeleitet / übernommen) je Zahl gekennzeichnet, übernommen nur 1.544 MiB (und die 1,9 bis 2,0 s des
  Architect-Verdikts, als solche benannt).
- **Vertrag und Werkzeuge:** `harness/targets/bench-backfill.md` stimmt mit den Skripten überein (Kill-Zusage für zwei
  Skripte, cgroup-Vorbedingung, Zeitreihen-Datei); die Reihen-Aufzählung des Messberichts (A bis O, P1 bis R) deckt die
  Überschriften des Zeilen-Dokuments vollständig.
- **Sprache:** keine Chronik/Forensik in Handbuch-Prosa, Vertrag oder Skript-Kommentaren (Suche der hinzugefügten Zeilen);
  keine internen Kennungen im Handbuch-Fließtext; kein host-lokaler absoluter Pfad (`make docs-check`).

## 8. Folge-Slice `slice-retention-lauf-speicher-begrenzung` (Plan in `open/`)

DoD tragfähig: jede Zeile trägt „Zu belegen durch“ oder „Erwartet“ mit Beleg-Weg (Unit-Test gegen handgeschriebene
Erwartung, Store-Tier mit vier Eingabeseiten-Mutationen, Nachmessung mit `tools/bench-backfill-memory.sh`, Skalierungs-Lauf,
Handbuch-Nachzug, Richtgröße-Bewertung, Release-Aussage). Trigger beobachtbar: Start über `ls docs/plan/planning/in-progress`,
`git ls-files docs/plan/planning/done` und den `Status`-Kopf der ADR; die Vorbedingung des Server-Release `v0.2.0` über
`git tag -l 'v*'` (gemessen: `v0.1.0`, `v0.1.1`, `v0.1.2`), als „ohne mechanischen Wächter“ benannt. Aussage zu `v0.1.x`
belegt (§5). Erwartung „Reihe J 13,3 / 13,2 / 13,7 MiB“ stimmt mit dem Zeilen-Dokument. Eine Nachzug-Stelle: V-2.

## 9. Harte Regeln

- **§3.1** — nur `make`, Repo-Skripte und `docker`; awk/`git`-Auswertungen auf dem Host, kein Host-Go/-Python, kein `sed -i`.
- **§3.2, §3.3, §3.5** — kein `nolint` im Diff; die zwei `git mv`-Commits sind rein; `make doc-immutable` Exit 0.
- **§3.9** — Exit-Codes nie durch eine Pipe; Gate-Lauf und Folgehandlung getrennt.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Befehl oder Lauf genannt) oder als **abgeleitet** gekennzeichnet
  (Rechnungen aus gedruckten Zeilen); übernommen sind nur die Findings-Anzahlen des Review-Reports.
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§6).
- **Commit-Traceability** — `make doc-commits` Exit 0 über 15 Commits.

## 10. Befunde und Restrisiken

Kein Befund blockiert die Closure-Bedingungen der DoD; V-1 (MEDIUM) und V-2 bis V-3 (LOW) sind Nachzüge für Planner und
Architect, keine Implementer-Arbeit am Code.

- **V-1 (MEDIUM) — Zahl im Träger gegen die Messung driftend, Negativaussage des Suchlaufs falsch.** Die
  [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) (`Accepted`) nennt in drei Zeilen (68, 191,
  199) „1,03 bis 1,57 KiB je Change (n = 15)“, der korrigierte Höchstwert ist 1,59; das Architect-Verdikt trägt sie in den Zeilen
  152, 155, 167. Das Plan-Feld §3 (Fixrunde) sagt „Nicht gefunden: kein weiterer Träger außerhalb der Suchwurzeln nennt den
  Höchstwert 1,57“ — die Suchwurzeln enthalten `docs/plan/adr` nicht. Die Entscheidung bleibt tragfähig („Faktor sieben bis
  zehn“ gegen etwa 0,15 KiB: 1,03 bis 1,59 durch 0,15 = 6,9 bis 10,6; die Bedarfsrechnung 864.000 × 1,03 bis 1,59 KiB ergibt
  0,85 bis 1,31 GiB), die ADR ist nicht in-place änderbar ([`AGENTS.md`](../../AGENTS.md) §3.5: Zahl, kein Zitat-Gerüst).
  Zuständig: Architect (Umgang mit der ADR), Planner (Plan-Feld berichtigen).
- **V-2 (LOW) — Suchausdruck im Plan des Folge-Slice.** Zeile 293 trägt `1,03 bis 1,57` als Befehl; er trifft im Handbuch
  nicht mehr (dort 1,59, §6). Nachzug an den Planner, nicht selbst geändert.
- **V-3 (LOW) — Zählwort des ersten Suchlauf-Feldes.** „Diff 280 (docs/user 71)“ ist der Stand von `4f94f900`; die Fixrunde
  hat 8 Zeilen im Handbuch ergänzt, am `HEAD` sind es 288 (79). Dasselbe Feld nennt „15 Reihen“ des Zeilen-Dokuments, mit den
  Reihen P1 bis R sind es 21. Instanz A: Zustandsgröße mit Stand nennen.
- **V-4 (INFO) — Reproduzierbarkeit der Schalter.** Die Varianten-Images (Reihen E, F, G, I, J) sind aus dem Arbeitsbaum mit
  einer geänderten Konstante gebaut und nicht committet; das Rezept steht als Prosa im Messbericht §2. Der Reviewer und ich
  haben die Ordnung an den lokal vorhandenen Images bzw. an der Code-Lesung getragen (§2 Zeile 2); ein Dritter baut die
  Images nach dem Prosa-Rezept. Die Aussage „Nicht reproduziert“ für F-10 (b) stützt der Plan auf einen Lauf `20260925T165631Z`,
  der im Zeilen-Dokument nicht steht; mein Lauf D bestätigt sie.
- **V-5 (INFO) — Fixrunde ohne zweiten Reviewer-Durchgang.** Die Fixrunde (`aff5f418`: Skripte, Vertrag, Handbuch, Berichte)
  ist von mir nachgemessen (§3, §4, §7), nicht vom Reviewer gelesen; das ersetzt keine Reviewer-Lesung der neuen Texte.
- **Restrisiko (weiter offen, im Folge-Slice getragen):** das Ausbleiben der Bereinigungs-Takte ab 2.000.000 Changes ist
  beobachtet und unerklärt; die Live-Erfassung als Speicherquelle ist aus dem Code gelesen, nicht gemessen; ein Host, n = 2
  bis 3 je Bedingung; die Release-Reihenfolge (Änderungs-Slice vor `v0.2.0`) hat keinen mechanischen Wächter; die Zahlen des
  Handbuchs §9 sind Zahlen vor der Behebung und werden im Folge-Slice ersetzt.

## Verdikt

**DoD des Plans: erfüllt** — alle sieben Liefer- und Belegzeilen (`[x]`) tragen einen realen Beleg (eigene Läufe §1 und §4,
Nachrechnung §3 und §7), die fünf offenen Zeilen sind korrekt offen (Closure durch den Planner). **Entscheidungs-Konformität:**
konform mit [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) und
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (kein Produktionscode, keine Schwelle, Richtgröße
bewertet, nicht geändert), [`AGENTS.md`](../../AGENTS.md) §3.5 und §3.9. **Plan-vs-Code:** jede der 17 Dateien gehört zu einer
Zeile der Plan-Tabellen (§3 und Fixrunde) oder ist fremde Arbeit (`ec7f76e2`, `3ec90077`, `5b1f7762`, `3405c9a1`).

**Übergabe an den Planner:** V-2 (Suchausdruck Zeile 293 im Plan des Folge-Slice auf `1,03 bis 1,59`), V-3 (Zählwörter), das
Plan-Feld „Nicht gefunden“ (V-1), dann Closure-Notiz; V-1 zusätzlich an den Architect (Umgang mit den drei Zeilen in
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)). Dieser Report ist ein Lauf-Beleg; er ändert
nichts.
