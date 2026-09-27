# Verifikations-Report: slice-leerlauf-phase-last-in-stuecken — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt:
[`review-slice-leerlauf-phase-last-in-stuecken.md`](review-slice-leerlauf-phase-last-in-stuecken.md)
(Commit `cea198fb`; 0 HIGH, 1 MEDIUM, 2 LOW, 3 INFO, keine Fixrunde). Kontext-Verdikt:
[`architect-verdict-leerlauf-bestaetigung-intermittenz.md`](architect-verdict-leerlauf-bestaetigung-intermittenz.md).
Formvorbild dieses Reports:
[`verifikation-slice-capture-leerlauf-quellbelege.md`](verifikation-slice-capture-leerlauf-quellbelege.md).

**Gegenstand:** Slice-Plan `slice-leerlauf-phase-last-in-stuecken` (wellenlos), Diff-Range
`6a976f58..38981f07` — 8 Commits, 10 Dateien, +136/−43 (`git diff --shortstat`). Inhalt: Lifecycle
(`ca93e19c` open→next, `c9a0308c` Verantwortlich, `828e1625` next→in-progress — die beiden Moves
rein, `git show --stat -M`: je 0 Zeilen), Implementer-Lauf (`2d41eb50` Runner-Phase
„Leerlauf-Bestätigung“ + `docs/user/e2e-abdeckung.md`, `e3ec5319` `harness/README.md` + Handbuch
Version 1.66, `3471c20d` Verweisform in fremden Plan-Dateien, `a1a665c2` Plan-Nachzug), DoD-Haken
`make gates` (`38981f07`). Review-Commit `cea198fb` liegt **hinter** dieser Range (Reviewer-Rolle,
nicht Implementer-Gegenstand). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur
diesen Report. Alle Mutationen liefen an Scratchpad-Kopien des Runners (Zusammensetzung der
Runner-Bereiche 1–248, 330–395, 2718–2862 und 3369–3484, wie beim Reviewer, gegen eine
Wegwerf-Compose-Umgebung mit den bestehenden Container-/Netz-Namen — sequentiell, ein Lauf zugleich,
`docker compose down -v` durch den eigenen Cleanup-Trap nach jedem Lauf); nie `sed -i`, nie ein
Host-Interpreter, keine Umleitung auf eine Repo-Datei. `git status --short` war nach jedem Lauf
leer.

**Repo-Zustand bei Beginn (nicht mein Zutun):** `origin/main` == lokales `HEAD` == `40e46221`
(`git rev-list --count origin/main..HEAD` = 0) — der Stand war bereits vor dieser Verifikation
gepusht, einschließlich zweier Planner-Commits hinter dieser Range
(`slice-harness-guard-blocked-python`, `slice-harness-mutationsbild-und-verweigerte-aktion`,
beide außerhalb dieses Gegenstands). Ich habe in diesem Lauf **nicht gepusht**.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert und
danach gelesen. `free -m` vor dem ersten schweren Lauf: 1,9 GB frei / 15,1 GB verfügbar (Puffer/Cache
15,4 GB) — gemessen, ein Wert, kein Engpass beobachtet. Je ein schwerer Docker-Lauf zugleich.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make image` | **Exit 0** | Digest `sha256:a8ef48aab93a…`, unverändert gegenüber `harness/image-hash.txt` vor dem Lauf — der Diff dieses Slice berührt keinen Produktionscode (§3) |
| `make test-integration` (ein voller Lauf, 07:57:18–08:03:51, 393 s) | **Exit 0** | Zeile der Phase, gedruckt: „… endete der Backfill-Run … über 30000 Zeilen completed und erzeugte 22188928 B WAL, ein Schreiber auf die nicht aktivierte Tabelle feed_e2e_wal_foreign erzeugte in 6 Stücken zu je 10000 Zeilen zusammen 16019584 B WAL (höchstes Stück 2676872 B, unter der Warnschwelle; nach jedem Stück erreichte confirmed_flush_lsn die Position hinter dem Stück, Frist 30 s); der Feed-Container lief über beide Lasten weiter (je 12 s beobachtet, kein Neustart), der Bestand ist über cdc.changes lesbar …, liegt bei 0 B“; Abschluss: `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; `git status --short` danach leer (keine committete Datei bewegt) |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md` | **Exit 0** | „11 Zeilen stimmen“ (Stände `d078700d`, `6a976f58`, `diff`) |
| `make commit-traceability RANGE=6a976f58..38981f07` | **Exit 0** | „OK — 8 Commit(s) …, Betreffs ohne Struktur-ID“ |
| `make doc-immutable RANGE=6a976f58..38981f07` | **Exit 0** | `d-check: 1339 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | „OK — 0 Commit(s) …“ (Repo-Zustand oben — bereits gepusht) |
| `make doc-immutable RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1339 Datei(en) geprüft, 0 Befund(e)` |
| `make gates` (vor diesem Report) | **Exit 0** | `generated-sync: OK`, `a-check: gesamt: 0 Befund(e)` (letzte Zeilen); vollständiger Lauf ungefiltert, Exit separat gesichert |
| **acht eigene Mutationen** an der Phase (§4) | sechs rot, zwei grün mit benannter Bedeutung | siehe Tabelle |

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **36**,
nach allen Läufen und Mutationen **36**; kein `prune`, kein `system prune`; nach jedem Lauf kein
Container (`docker ps -a` leer) und kein Netz `cdc-*`. CI-Läufe nachgelesen mit `gh api`/`gh run
view --log`, keine Schreibhandlung auf GitHub.

Nicht gefahren: ein zweiter voller `make test-integration`-Lauf (Laufzeit-Streuung nicht gemessen);
`make test-store`/`make test-replication`/`make bench` (außerhalb des Gegenstands); jede
Docker-Image-Mutation mit verändertem Produktionscode (§5, F-3/§8).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 8 `[x]`-Zeilen (1–3, Gates, Review, Suchlauf, Doku), 5 `[ ]`-Zeilen (DoD-5 e2e-Lauf,
Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen — sechs Zeilen
in Summe `[ ]`, siehe Plan §2). Ich setze keinen Haken (Planner).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Last in Stücken, Wartebedingung, zwei Wächter (`[x]`) | **bestätigt** | Diff der Phase (`tools/harness/run-integration-tests.sh:3369-3465`, §3 unten); mein voller `make test-integration`-Lauf grün mit exakt der zugesagten Form (6 Stücke, höchstes Stück unter Warnschwelle, Summe über Fehlerschwelle, Frist 30 s eingehalten); acht eigene Mutationen bestätigen jede Klausel unabhängig (§4) — inklusive der vom Review benannten (a)/(b)/(c)/(d) und einer eigenen, neuen (H: Ein-Stück-Rückfall auf 60.000 Zeilen, rot am Wächter je Stück). **DoD-Beleg-Anker fehlt weiterhin im Plan-Text selbst** (Review F-1, unten bestätigt) |
| 2 | Träger nachgezogen (`[x]`) | **bestätigt** | `harness/README.md` Zeile 141 (nur der Satzteil zur Phase geändert, Zeile zu `make test-replication` unberührt), Handbuch Version 1.66 mit zwei geänderten Stellen + neuem Absatz + Änderungshistorie-Zeile, `docs/user/e2e-abdeckung.md` (Erzeugnis, in meinem Lauf „unverändert“ bestätigt); `make suchlauf-nachmessen` und `make docs-check` (Teil von `make gates`) beide Exit 0 |
| 3 | Backfill-Hälfte gemessen (`[x]`) | **übernommen, nicht selbst nachgemessen** | Implementer-Messung im Plan (1500 Proben, größter Abstand 23 ms, Spitze 741.816 B = 17,7 % der Hälfte der Fehlerschwelle); vom Reviewer bereits mit unabhängigen WAL-Werten aus drei Läufen (lokal 19.826.288 B, CI 22.064.160 B/22.186.272 B, mein eigener Lauf 22.188.928 B) plausibilisiert — alle unter dem gemessenen Backfill-Volumen, keiner löste in irgendeinem Lauf ein Rot der Leerlauf-Phase aus. Ich habe die 50-ms-Sondenmessung selbst nicht wiederholt (Aufwand/Nutzen: vier unabhängige Läufe ohne Rot sind ein stärkeres Signal als eine fünfte Einzelmessung) |
| 4 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf Exit 0 (§1) |
| 5 | Erster grüner `e2e.yml`-Lauf mit beiden Legs (`[ ]`) | **bestätigt** (Beleg unten) | Lauf **36295604780**, `run_attempt` **1**, `conclusion` `success`; Jobs **108553753483** (PostgreSQL 18) und **108553753527** (PostgreSQL 17), je `success`; Job-Logs (`gh run view --job --log`) tragen die Phasen-Ausgabezeile mit „6 Stücken zu je 10000 Zeilen“ (§3 unten). Verifier-Regel aus Verdikt §4 greift nicht — kein Rot der Signatur zu diesem Stand |
| 6 | Review durchgeführt (`[x]`) | **bestätigt** | Report liegt vor, Commit `cea198fb`, 0 HIGH/1 MEDIUM/2 LOW/3 INFO |
| 7 | §3.13-Suchlauf im committeten Feld (`[x]`) | **bestätigt** | `make suchlauf-nachmessen` Exit 0, 11 Zeilen (§1) |
| 8 | Doku-Update `harness/README.md` + Handbuch (`[x]`) | **bestätigt** | wie Zeile 2 |
| 9 | Closure-Notiz mit Lerneintrag (`[ ]`) | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 10 | Reconciliation-Register — entfällt (`[ ]`) | **korrekt offen** | keine Reconciliation-Datei im Repo (Greenfield); Haken bei Closure |
| 11 | Beobachtungs-Register fortgeschrieben (`[ ]`) | **korrekt offen, teilweise bereits geführt** | `BEO-PGC/test-integration-retention-timing-flake/state.md` führt den Ausgang des Leerlauf-Falls bereits mit Ursache und Kette (Zeilen 33–61, s. §6); `BEO-PGC/ersatzweg-nach-verweigerter-aktion` führt F-3 mit eigener Evidence-Datei; der Implementer-Aufruf `python3 -c 1` steht in `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md` Zeile 58 genannt, **ohne eigene Evidence-Datei** — Planner-Pflicht bei Closure (bestätigt, siehe Auftragstext) |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang (`[ ]`) | **korrekt offen** | alle sechs Zeilen tragen „**Ausgang:** *(bei Closure)*“ |
| 13 | Drei Paarungen (`[ ]`) | **korrekt offen** | hängt an der Closure von `welle-transformationen` |

Kein `[x]` ohne Beleg. DoD 1 trägt die von Review F-1 benannte Lücke weiter (Substanz belegt, Anker
im Plan-Text fehlt) — siehe §5.

## 3. Plan-vs-Code-Diff

**§1 Ziel und „Ausdrücklich NICHT“:** kein Produktionscode berührt. Nachgemessen:
`git diff --name-only 6a976f58..38981f07 -- internal cmd gen test compose.yaml Makefile .github
docs/plan/adr spec` ist leer. Weder Schwellen noch `runWALRetentionCheck`, `ADR-0049`, `ADR-0120`
noch die Phase „Fehlerschwelle beendet den Container“ liegen im Diff.

**Diff der Phase** (`tools/harness/run-integration-tests.sh:3369-3465`, gelesen): `abdeckung_declare`
um den Stücke-Satz erweitert; Kopf-Kommentar um vier Sätze zur Stück-/Wartelogik ergänzt;
`WAL_FOREIGN_CHUNKS=6`, `WAL_FOREIGN_CHUNK_ROWS=10000`, `WAL_FOREIGN_CONFIRM_SECONDS=30` neu; der
einzelne `INSERT … generate_series(1, 60000)` ersetzt durch eine `for`-Schleife über 6 Stücke mit
disjunkten Schlüsselbereichen, Wächter je Stück (`bf_fail` bei `wal_chunk_bytes` ≥ `WAL_WARN_BYTES`),
`bf_await_sql` auf `confirmed_flush_lsn` je Stück, Summen-Wächter nach der Schleife (`bf_fail` bei
Summe ≤ `WAL_ERROR_BYTES`); Ausgabezeile um Stückzahl/-größe/höchstes Stück/Frist erweitert.

| Plan-Zeile (§3) | Ist im Diff |
|---|---|
| Runner-Phase (Schreiber in Stücken, Wartebedingung, zwei Wächter; Kommentar, `abdeckung_declare`, Ausgabezeile) | `M`, wie oben beschrieben; entspricht dem Plan wörtlich |
| `docs/user/e2e-abdeckung.md` (Erzeugnis) | `M`, in meinem Lauf „unverändert“ bestätigt (Erzeugnis am Ist-Quelltext) |
| `harness/README.md` §Sensors (Zeile `make test-integration`) | `M`, nur der Satzteil zur Phase; Nachbar-Zeile `make test-replication` (Form X1) unberührt |
| `docs/user/benutzerhandbuch.md` (zwei Stellen, neuer Absatz, Version/Historie) | `M`, Version 1.65→1.66, Stand 2026-09-25→2026-09-27, neue Zeile in „### Änderungshistorie“ |
| Verweisform in fremden Plan-Dateien (Roadmap 2×, `welle-transformationen`, zwei Register-`state.md`/Folge-Pläne 3×) | `M`, `git diff --word-diff` zeigt ausschließlich `[Link](../pfad.md)` → `` `Kennung` `` — keine Aussage geändert (nachgemessen für alle fünf Dateien) |
| Lifecycle (3 Commits) | `A`/`M` wie Plan; beide Moves rein (§ oben) |

Kein Betreiber-Weg, kein Gate-/Workflow-Struktur-Diff (`.github` nicht berührt,
[`AGENTS.md`](../../AGENTS.md) §3.10 greift nicht — Verdikt §3 Abgrenzung).

## 4. Mutationen (dieser Lauf, eigene Auswahl, unabhängig vom Review)

Acht Mutationen an Scratchpad-Kopien der Phase (Zusammensetzung wie Reviewer, s. o.), sequentiell
gegen eine frische Wegwerf-Compose-Umgebung, je ein Lauf. Baseline (unmutiert, dieselbe
Zusammensetzung, endet nach der Phase mit `VLP-PHASE-DONE exit0`): **grün**, 65 s, Ausgabezeile
identisch zur Zusage (Backfill-Run 19.826.288 B, Stücke zusammen 16.019.504 B, höchstes Stück
2.676.792 B — **exakter Bytewert wie im Grundlauf des Reviewers**, weil die Last deterministisch
über feste Zeichenketten/Zeilenzahlen erzeugt wird).

| # | Mutation | Ort | Ergebnis (gesehene Zeile) |
|---|---|---|---|
| A | `WAL_FOREIGN_CHUNK_ROWS` 10.000 → 30.000 | `:3389` | **rot**, 48 s: „Stück 1 von 6 erzeugte 8013480 B WAL, nicht weniger als die Warnschwelle 4194304 B: das Stück ist ein Stoß, keine Bestätigungs-Runde“ |
| B | `WAL_FOREIGN_CHUNKS` 6 → 2 | `:3388` | **rot**, 46 s: „die 2 Stücke erzeugten zusammen 5345184 B WAL, nicht mehr als die Fehlerschwelle 8388608 B: die Last trägt den Beleg nicht“ |
| C | Position `+ 1073741824` (1 GiB) hinter dem Stück (`::pg_lsn + …`) | `:3461` | **rot** nach 30 s, 76 s gesamt: „confirmed_flush_lsn des Slots erreichte die Position hinter Stück 1 von 6 nicht innerhalb von 30 s — erwartet 't', gelesen 'f'“ |
| D | `docker pause` des Feed-Containers vor dem ersten Stück | vor `:3450` | **rot** nach 30 s, an derselben Zeile wie C |
| E | Position **vor** dem Stück (`wal_chunk_before` statt `_after`) | `:3461` | **grün**, 66 s — Ausgabezeile unverändert; Grenze, kein Fehler (deckt sich mit Review m4) |
| F | Wartebedingungs-Zeile ganz gestrichen | `:3461` | **grün**, 66 s — dieselbe Ausgabezeile; Grenze, kein Fehler (deckt sich mit Review m8) |
| G | Slot-Filter `slot_name = '$SLOT'` → `'kein_slot'` | `:3461` | **rot** nach 30 s, 75 s: „… erwartet 't', gelesen 'leer' (… slot_name = 'kein_slot')“ — nicht trivial erfüllt |
| H | **eigene, im Review nicht gefahrene Mutation:** `WAL_FOREIGN_CHUNKS=1`, `WAL_FOREIGN_CHUNK_ROWS=60000` (Rückfall auf die alte Ein-Stoß-Form, mit dem neuen Wächter) | `:3388-3389` | **rot**, 46 s: „Stück 1 von 1 erzeugte 16019264 B WAL, nicht weniger als die Warnschwelle 4194304 B“ — der Wächter fängt einen Rückfall auf die alte Form unabhängig von der Wartebedingung |

Hygiene: `docker ps -a`/`docker network ls | grep cdc` nach jeder Mutation und am Ende leer;
dangling Volumes vor/nach **36**; `git status --short` durchgehend leer.

**Ergebnis:** sechs von acht Mutationen bestätigen die Fehlerseite unabhängig von den zehn
Reviewer-Läufen (A/B/C/D/G decken sich mit dessen m1/m2/c/d/m3; H ist eine eigene, zusätzliche
Bindung: der Wächter je Stück verhindert auch einen stillen Rückfall auf einen einzelnen Stoß).
E und F reproduzieren unabhängig die vom Review benannte Grenze (F-2): die Wartebedingung ist an
ihrer Fehlerseite gebunden, ihre Schutzwirkung (dass sich Stücke nicht stapeln) ist es nicht — ohne
eine künstliche Verzögerung der Bestätigung liefert der Aufbau keinen Fall, der das Warten selbst
braucht.

Nicht gefahren: Mutation (e)/(f) des Plans (Produktionscode-Mutation, `ConfirmIdle`/`confirmIdle`) —
siehe §5, §8 (VERWEIGERTE-AKTION-REGEL: kein freigegebener Weg für ein Mutations-Image mit eigenem
Tag; **übernommen** aus Review-Messung, nicht selbst gefahren). Kein Docker-Image gebaut oder
mutiert in diesem Lauf.

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) DoD 1 ohne committeten Beleg-Anker im Plan-Text | Plan §2 DoD 1 gelesen: die Mutationen (a)–(f) stehen mit gedruckter Zeile im Plan, der grüne Volllauf selbst (Dauer, Ausgabezeile, Stand) nicht. Mein eigener Volllauf (§1) und die zwei CI-Leg-Zeilen (§2 Zeile 5) liefern den fehlenden Anker, stehen aber ebenfalls nur in diesem Report, nicht im Plan | **bestätigt**, Lücke besteht fort — Planner-Pflicht bei Closure |
| F-2 (LOW) Wartebedingung an Fehlerseite gebunden, Schutzwirkung nicht | eigene Mutationen E und F (§4) reproduzieren m4/m8 unabhängig: beide grün, dieselbe statische Ausgabezeile | **bestätigt**, unabhängig reproduziert |
| F-3 (LOW) Ersatzweg nach verweigerter Aktion, Route nicht im Plan | Sachverhalt außerhalb des Diffs, nicht durch mich beobachtbar (kein Image gebaut, kein `docker buildx`-Aufruf in diesem Lauf); Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion` führt den Fall bereits mit eigener Evidence-Datei und einer Nutzer-Entscheidung („ja“, 2026-09-27) zur Regel; der Folge-Slice `slice-harness-mutationsbild-und-verweigerte-aktion` liegt bereits in `open/` | **bestätigt** (Registerstand), **an den Auftraggeber** — Frage unten |
| F-4 (INFO) Mutation (f) im Unit-Tier gebunden | nicht selbst nachgefahren (Produktionscode-Mutation, kein Toolchain-Bau in diesem Lauf); Reviewer-Zahl (fünf rote Tests im Paket `receive`) plausibel anhand der Testnamen im Diff (`TestRunConfirmsIdleWithoutReplyRequested` u. a., bereits vor diesem Slice bestehend, nicht Gegenstand des Diffs) | **übernommen**, nicht nachgemessen |
| F-5 (INFO) `ADR-0120`/`ADR-0129` tragen den breiteren Satz | `git grep -n 'jeden Schreiber auf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'` — ein Treffer in `ADR-0120:243` (`Accepted`, unberührbar); der Verdikt (§2.4) benennt die Präzisierung als bewusste Entscheidung ohne Folge-ADR | **bestätigt**, kein Verstoß gegen §3.5 (keine inhaltliche Änderung an der ADR versucht) |
| F-6 (INFO) Begriff „Bestätigungs-Runde“ ohne Einführung im Handbuch | `grep -n 'Bestätigungs-Runde' docs/user/benutzerhandbuch.md` — ein Treffer, Zeile 887 | **bestätigt**, redaktioneller Hinweis ohne DoD-Wirkung |

Kein offenes HIGH. Das MEDIUM (F-1) besteht als reine Text-/Ankerlücke fort, keine DoD-Verletzung
in der Substanz.

## 6. Register und Träger — Lese-Prüfung

- **`BEO-PGC/test-integration-retention-timing-flake/state.md`:** führt den Ausgang des
  Leerlauf-Falls bereits vollständig (Zeilen 33–61): Ursache (Stoß, 1,9-fache Fehlerschwelle),
  Abhilfe (Stücke, lokal erprobt 32 % der Fehlerschwelle), Kette der Folge-Slices, Trigger für eine
  Neubewertung samt der zwei `gh`-Befehle. Der Zähler bleibt bei 4× (Resthälfte 2×: `slice-057` ohne
  Ursache, der Leerlauf-Fall mit Ursache und Ausgang) — konsistent mit dem, was der Architect dort
  bereits eingetragen hat; keine Korrektur nötig.
- **`BEO-PGC/ersatzweg-nach-verweigerter-aktion`:** Zustand „geplant“, Ausgang an
  `slice-harness-mutationsbild-und-verweigerte-aktion` (Liefer-Punkt 3) gebunden; eigene
  Evidence-Datei für diesen Slice liegt vor (`evidence/slice-leerlauf-phase-last-in-stuecken.md`).
  Nichts nachzutragen.
- **`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md` Zeile 58:** nennt den
  Implementer-Aufruf `python3 -c 1` in diesem Slice **ohne eigene Evidence-Datei** — bestätigt fehlend
  (`find docs/plan/planning/observations/BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/evidence`
  zeigt keine Datei zu diesem Slice). Planner-Pflicht bei Closure, wie im Auftrag benannt.
- **Kommentar „nur, wenn der Slot nichts bestätigt“** (`tools/harness/run-integration-tests.sh`,
  Phase „Fehlerschwelle beendet den Container“, Kopf davor): bereits mit Adresse im Plan §3 dieses
  Slice geführt (Frist: Start von `slice-wal-fehlerschwelle-ausgangsklasse`) — `git grep -n
  'nur, wenn der Slot nichts bestätigt' -- tools/harness` bestätigt eine Fundstelle, unverändert seit
  `8611185b` laut Review F-5b. Kein zusätzlicher Träger gefunden.
- **`ADR-0120` Zeile 243** („dasselbe gilt für jeden Schreiber auf Tabellen ohne Publication-Bezug“):
  bestätigt als Rand (s. §5, F-5) — der Architect hat den Gewinner im Verdikt §2.4 bewusst deklariert,
  keine Folge-ADR nötig; ein Leser der ADR allein sieht den breiteren Satz ohne Zeiger dorthin. Kein
  Verstoß, nur benannt wie im Auftrag verlangt.

## 7. Entscheidungs-Konformität

- **[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md):** unberührt (`Accepted`),
  Regel und Fitness-Function-Zeilen bleiben; die Phase belegt jetzt einen Fall, den der Mechanismus
  trägt, statt einen, den er nicht trägt (Verdikt §2.3/§2.4). Konform.
- **[`ADR-0129`](../plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md):** nur als
  Beleg referenziert (Zeilen 145/218 unverändert laut Plan §3), nicht im Diff. Konform.
- **[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md):** die Gegenseite
  („Fehlerschwelle beendet den Container“) bleibt unberührt und ist ausdrücklich abgegrenzt
  (Verdikt §3 Abgrenzung, Plan §1); die neuen Wächter der Leerlauf-Phase ändern keine Schwelle,
  sondern die **Last-Form**. Konform.
- **[`ADR-0030`](../plan/adr/0030-testpyramide.md):** die Phase bleibt im E2E-Tier
  (`make test-integration`, Black-Box über Docker); kein neuer Mock. Konform.
- **[`AGENTS.md`](../../AGENTS.md) §3.5:** `ADR-0120`/`ADR-0129` inhaltlich unverändert
  (`make doc-immutable` Exit 0, §1); die Präzisierung aus Verdikt §2.4 lebt außerhalb der ADR
  (Handbuch/README), keine Zitat-Korrektur versucht.
- **§3.12/§3.13:** Zahlen im Plan tragen ihren Ursprung (**übernommen**/**gemessen**/**hergeleitet**
  konsequent markiert, z. B. Verdikt-Zahlen als „**übernommen**“); Suchlauf-Feld trägt beide Stände
  (Parent `d078700d`, `6a976f58`, `diff`) und Gefundenes/Nichtgefundenes je Träger. Konform.
- **§3.1 Docker-only:** meine eigenen Mutationen liefen ausschließlich über Scratchpad-Kopien des
  Runners mit `docker exec`/`docker compose`/`docker pause`, kein Host-Interpreter, kein `sed -i`.
  Für den Implementer-Sachverhalt (F-3) siehe §8 — außerhalb meines Diffs, nicht durch mich
  reproduzierbar.
- **§3.10:** greift nicht (`.github` nicht im Diff, Verdikt §3 Abgrenzung).

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **DoD 1 trägt den Beleg-Anker weiterhin nicht im Plan-Text** (F-1 bestätigt, unverändert seit dem Review). Substanz belegt (mein Volllauf, meine acht Mutationen, die zwei CI-Legs); der Plan-Text selbst verweist noch nicht auf einen committeten Lauf-Beleg — dieser Report kann als Anker dienen | Plan §2 DoD 1 | ja — Lesen |
| V-2 | INFO | **Mutation H (eigene, unabhängig vom Review) bestätigt eine zusätzliche Bindung:** der Wächter je Stück fängt auch einen stillen Rückfall auf die alte Ein-Stoß-Form (`WAL_FOREIGN_CHUNKS=1`, 60.000 Zeilen), unabhängig von der Wartebedingung. Kein Finding gegen den Slice — stärkt die DoD-1-Substanz zusätzlich | §4, Mutation H | ja — eigener Lauf |
| V-3 | INFO | **Grenze F-2 unabhängig reproduziert** (Mutationen E/F): die Schutzwirkung der Wartebedingung (Stücke stapeln sich nicht) ist am gesunden Aufbau nicht bindbar, nur ihre Fehlerseite. Kein Widerspruch zum Ziel aus Plan §1; deckt sich exakt mit Review m4/m8 | §4, Mutationen E/F | ja — eigener Lauf |
| V-4 | INFO | **Register-Lücke bestätigt:** `python3 -c 1` (Implementer) in `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel/state.md` Zeile 58 genannt, ohne eigene Evidence-Datei — Planner-Pflicht bei Closure, wie im Auftrag benannt | §6 | ja — `find` im Evidence-Verzeichnis |
| V-5 | INFO | **Repo-Zustand bereits gepusht** (nicht mein Zutun): `origin/main` == `HEAD` bei Beginn dieses Laufs, einschließlich zweier Planner-Commits außerhalb des Gegenstands. Kein Finding gegen den Implementer (Push ist nicht Teil seines Auftrags); zur Kenntnis für die Closure | Kopfzeile dieses Reports | ja — `git rev-list --count origin/main..HEAD` = 0 |

Kein HIGH, kein MEDIUM in dieser Verifikation selbst (das MEDIUM des Reviews, F-1/V-1, besteht als
Text-/Ankerlücke fort, keine DoD-Verletzung in der Substanz).

## 9. Verdikt

**DoD bestätigt:** ja, in der Substanz — DoD 1 (Last in Stücken, Wartebedingung, zwei Wächter: eigener
Volllauf + acht eigene Mutationen, sechs rot/zwei grün mit benannter Grenze), DoD 2/8 (Träger
nachgezogen, `make suchlauf-nachmessen`/`make docs-check` grün), DoD 3 (Backfill-Hälfte,
**übernommen**, durch vier unabhängige rotfreie Läufe plausibilisiert), DoD 4 (`make gates` Exit 0,
eigener Lauf), DoD 5 (**erster grüner `e2e.yml`-Lauf**: 36295604780, Versuch 1, beide Legs
`success`, Phasen-Zeile mit „6 Stücken zu je 10000 Zeilen“ in beiden Job-Logs bestätigt), DoD 6–8
(Review vorliegend, Suchlauf-Feld Exit 0, Doku-Update). DoD 9–13 korrekt offen (Closure-Pflicht des
Planners). **Plan-vs-Code:** keine unbenannte Abweichung, kein Produktionscode im Diff.
**Entscheidungs-Konformität:** `ADR-0120`, `ADR-0129`, `ADR-0049`, `ADR-0030` konform; `ADR-0120`
Zeile 243 als bewusst unaufgelöster Rand bestätigt (kein Verstoß). **Review-Findings:** F-1 bis F-6
gemessen, F-1/F-2 unabhängig durch eigene Mutationen bestätigt, F-3/F-4 übernommen (außerhalb meines
Diffs bzw. Produktionscode-Mutation ohne freigegebenen Weg), F-5/F-6 bestätigt. **Mutationen:** acht
eigene, unabhängig von den zehn des Reviews (sechs Überschneidungen im Ergebnis, zwei
Grenz-Reproduktionen, eine neue Bindung H). **Gates:** `make gates`, `make commit-traceability`
(sowohl über `6a976f58..38981f07` als auch `origin/main..HEAD`), `make doc-immutable`,
`make suchlauf-nachmessen` — alle Exit 0 im eigenen Lauf.

### Verweigerte Aktion — keine in diesem Lauf

Ich habe in diesem Verifikationslauf **keine** Docker-Image-Mutation mit verändertem
Produktionscode versucht (VERWEIGERTE-AKTION-REGEL: kein freigegebener Weg, `make image-mutation`
existiert noch nicht). Mutation (e)/(f) des Plans bleiben **übernommen** aus Review-Messung, nicht
selbst gefahren. Kein Aufruf wurde von der Berechtigungsschicht in diesem Lauf verweigert.

### Übergabe an den Planner

1. **DoD-Haken:** DoD 1, 2, 3, 4, 5, 6, 7, 8 können gesetzt werden (Substanz belegt, s. Tabelle §2);
   der Plan-Text von DoD 1 sollte diesen Report (oder den Volllauf/die CI-Zeile) als Anker nennen,
   um F-1/V-1 endgültig zu schließen.
2. **Register:** `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` — Evidence-Datei für den
   `python3 -c 1`-Aufruf dieses Slice anlegen (V-4); die übrigen drei Register-Einträge
   (`test-integration-retention-timing-flake`, `ersatzweg-nach-verweigerter-aktion`,
   „Kommentar davor“ an `slice-wal-fehlerschwelle-ausgangsklasse“) sind bereits geführt, keine
   Nacharbeit nötig.
3. **Risiken §6:** alle sechs Zeilen können auf einen Ausgang gesetzt werden — „Bestätigung
   bestätigt nichts färbt nicht an der Wartebedingung“ (eingetreten, im Plan bereits mit (e)
   dokumentiert und durch Review-Mutation gefahren), „Stücke über/unter Schwelle“ (entfallen, meine
   acht Mutationen und der reale Lauf zeigen die Wächter greifen), „Backfill-Hälfte ist ein Stoß“
   (entfallen, DoD 3), „Rot mit Signatur im Zyklus vor `done`“ (entfallen, DoD 5 grün im ersten
   Versuch), „Wartebedingung fällt still aus dem Runner“ (entfallen, Ausgabezeile + `e2e-abdeckung.md`
   tragen sie), „Träger sagt Zusage breiter“ (eingetreten und behandelt, F-5/§5).
4. **Drei Paarungen:** bleiben an der Closure von `welle-transformationen` hängen (Plan §7,
   korrekt).
5. **Für den Nutzer vorzulegen:** die offene Frage aus Review F-3 — welcher Weg für ein
   Mutations-Image mit eigenem Tag gilt (`make image` aus `git archive`-Kopie mit Wiederherstellung
   von `:dev`, oder ein benannter Wegwerf-Weg über `slice-harness-mutationsbild-und-verweigerte-aktion`)
   — ist bereits als offener Slice in `open/` angelegt; keine zusätzliche Eskalation nötig.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch
Closure.
