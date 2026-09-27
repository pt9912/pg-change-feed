# Review-Report: ADR-0129 (Quellseite der Leerlauf-Bestätigung als Test an beiden PostgreSQL-Pins) — 2026-09-27

**Review-Art:** Design — die ADR wird gegen `AGENTS.md` §3.12 „Verfasser einer ADR“ (Menge, Ursprung,
Mutationsangaben), §3.5 (Immutabilität von `ADR-0120`/`ADR-0121`), §3.7 (Ton) und gegen die gemessenen Belege
geprüft; nicht gegen die DoD (Verifier-Frage) und nicht gegen Code und Runner des Slice (bereits im Review-Report
`review-slice-capture-leerlauf-quellbelege` und im Verifikations-Report zum selben Slice behandelt).

**Gegenstand:** [`ADR-0129`](../plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md), Commit `4e654153`,
Diff `e4b77a05..4e654153` = die ADR (263 Zeilen, neu) plus eine Index-Zeile und die geänderte Titel-Zeile von `ADR-0121`
in `docs/plan/adr/README.md` (`git diff --stat`: 2 Dateien, +265/−1). DoD 3 und das Risiko „ADR behauptet mehr als der
Beleg“ (Plan §6) des Slice `slice-capture-leerlauf-quellbelege`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (repo-lokal, kein Tag-Bezug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-27

**Eingangs-Kontext:**

- [`ADR-0129`](../plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md), [`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md),
  [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md), [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) (nur zur Abgrenzung)
- Verifikations-Report [`verifikation-slice-capture-leerlauf-quellbelege`](verifikation-slice-capture-leerlauf-quellbelege.md) (§1 gedruckte Zeilen, §4 Mutationstabelle K1–K10, W, P1, P2) und
  Review-Report [`review-slice-capture-leerlauf-quellbelege`](review-slice-capture-leerlauf-quellbelege.md) (Mutationen M1–M10, nur als Gegenprobe)
- Test `internal/adapters/driving/replication/receive/sourcekeepalive_test.go`, Runner `tools/harness/run-replication-tests.sh`,
  Workflow `.github/workflows/e2e.yml`
- `LH-QA-REL-001`, `LH-FA-CAP-009`, `LH-QA-POR-001`, `SPEC-012`
- [`AGENTS.md`](../../AGENTS.md) §3.5, §3.7, §3.12, §3.13

**Ursprung dieses Berichts.** Was der Report als Zahl oder Farbe nennt, ist eines von drei: *gemessen* (mein Lauf,
mit Zeile), *nachgeschlagen* (GitHub-Abfrage, Git-Abfrage) oder *übernommen* (aus einem Report; dann so gekennzeichnet).

---

## Nachgemessen (Ursprung der ADR-Aussagen)

**Gedruckte Zeilen (ADR §Gemessen ↔ Verifikations-Report §1).** Wort für Wort verglichen: PostgreSQL 18.6, 11080 von
400000, `ServerWALEnd`/Commit-LSN/bestätigt `0/12786CB0`, 40.32 s; PostgreSQL 17.11, 11077 von 400000, `0/124C79A8`,
40.37 s — beide stehen so im Report. Das Test-Godoc, `t.Logf` (Zeile 257) und die Aufbau-Beschreibung der ADR
(38 s Stille, `proto_version 1`, Standard-`wal_sender_timeout`, eine Transaktion über 400.000 Einfügungen, Neustart ab
`confirmed_flush_lsn`) stimmen mit dem Test überein (gelesen). Die Pins: `e2e.yml` trägt beide Digests, der 18er steht
zugleich als `PG_TEST_IMAGE` im `Makefile` und im Runner (nachgeschlagen).

**Mutationszuordnung (ADR §Fitness Function ↔ Verifikations-Report §4).** Mutation 1 = K1 (17 und 18), Mutation 2 = K4
(17 und 18), Mutationen 3 bis 6 = K3, K5, K8, K10 (je nur 18), Mutationen 7 und 8 = K6 (grün), K7 (grün, ein Lauf),
Mutation 9 = W. Die Stellen (`confirmedPosition := …`, Test-Zeile 263), die Farben und die Texte der roten Zeilen
stimmen mit dem Report überein; „auch an 17 rot“ steht für 3 bis 6 als *hergeleitet*. Bei Mutation 9 nennt der Report
tatsächlich keine Version (W-Zeile gelesen); die ADR sagt das so. Kein Fall, in dem die ADR „erprobt“ schreibt, wo der
Verifier nicht gefahren hat.

**Eigene Mutationen an PostgreSQL 17** (Pin-Digest aus `e2e.yml`, Standard-`wal_sender_timeout`, gepinnter
Toolchain-Container in einem Wegwerf-Netz, Kopien der Testdatei im Scratchpad, je eine Änderung; die Version 17.11
druckte keiner der drei Läufe, sie endeten vor der Log-Zeile):

| # | Stelle | Mutation | Gesehene Farbe (gemessen) |
|---|---|---|---|
| R1 | Test-Zeile 229 | `serverWALEnd != finalLSN` → `!= finalLSN+8` (die Zeile, die Aussage (a) bindet und die keine Mutation der ADR setzt) | **rot** nach 38,17 s: `ServerWALEnd … des Keepalives ist nicht die Commit-LSN … der Transaktion` |
| R2 | Test-Zeile 250 | Neustart ab `confirmed_flush_lsn + 1` statt ab `confirmed_flush_lsn` (Gegenprobe zu „ungebunden“, Festlegung 2 Punkt 4) | **rot** nach 78,19 s: `Der Neustart des Streams lieferte die Transaktion nicht bis zum COMMIT (0 Änderungen gelesen) … context deadline exceeded` |
| R3 | Test-Zeile 220 | die Quelltransaktion braucht 45 s (`generate_series(…) g, (SELECT pg_sleep(45)) s`), also länger als die 38 s Stille — die Herleitung in Festlegung 2 Punkt 5 | **rot** nach 48,11 s: `COMMIT nach 400000 Änderungen ohne Keepalive zwischen BEGIN und COMMIT` |

R3 bestätigt die als *hergeleitet* geführte Aussage (Test endet laut statt mit falscher Aussage) an einem Lauf an
PostgreSQL 17; R1 zeigt, dass Aussage (a) an der Stelle gebunden ist, an der die ADR sie bindet — die ADR belegt das
nur nicht selbst (F-4). R2 widerspricht einem Teil von Festlegung 2 Punkt 4 (F-2).

**CI-Lage (GitHub, nachgeschlagen).** Lauf `36287009221` zum Commit `e4b77a05`: Versuch 1 — Leg PostgreSQL 18 `failure`
im Schritt „Compose-Integrationstest (Black-Box-E2E)“, Schritt „Replication-Tier“ `skipped`; das Log nennt in der Phase
„Leerlauf-Bestätigung“ `WAL-Rückstand 15238216 Bytes über Fehlerschwelle 8388608 Bytes`; Leg PostgreSQL 17 (Job
`108529548457`) alle Schritte `success`, beendet 02:11:03 UTC, also **nach** dem in der ADR genannten Stand („gegen 02:07 UTC“;
Commit der ADR: 02:07:45 UTC). Versuch 2 (`--failed`): Leg PostgreSQL 18 (Job `108531740887`) grün. In CI gedruckt: PostgreSQL 17.11,
Keepalive nach 11538 von 400000 Änderungen, `PASS` in 40.72 s (Job `108529548457`); PostgreSQL 18.6, nach 11541,
`PASS` in 43.20 s (Job `108531740887`). Die Angaben der ADR zum Stand vom ersten Versuch sind wahr; die ADR trägt keine
Versuchsnummer (F-1).

**Träger-Suchlauf der ADR (Folgepflicht 2).** Das Muster der ADR am Stand `e4b77a05`, ohne `docs/reviews` und
`.harness/baseline`: Treffer in `ADR-0119`, `ADR-0121`, `ADR-0125`, in `done/slice-backfill-slot-leerlauf-bestaetigung`, im
Plan des Slice und in drei Dateien des Registers `beleg-nur-als-einmalige-reviewer-messung` — genau die Liste der ADR;
`spec`, `docs/user`, `harness`, `README.md`, `AGENTS.md`: kein Treffer (Exit 1 von `git grep -c`). `ADR-0119` Zeile 48/89 und
`ADR-0125` Zeile 144 haben einen anderen Gegenstand (gelesen). Die Lokatoren zu `confirmIdle` (`receive.go:467-469`) und
`TestRunNoConfirmationInsideOpenTransaction` (`seam_test.go:694-697`) tragen am Stand `4e654153` den genannten
Begründungssatz. `make suchlauf-nachmessen` am Plan des Slice am HEAD: Exit 0, 10 Zeilen (Commit-Stände bereits gesetzt).

**Immutabilität.** `git diff --stat e4b77a05..4e654153` nennt `ADR-0120` und `ADR-0121` nicht: beide Dateien sind
unberührt (`Accepted`, §3.5). Die `Supersedes`-Kette (Kopf: Festlegung 2 und Konsequenz-Punkt „Negativ (Grenze, benannt)“
von `ADR-0121`; die dort wörtlich vorhandenen Stellen gelesen) ist eine teilweise Ablösung, wie sie `ADR-0118`,
`ADR-0119`, `ADR-0121` und `ADR-0126` vormachen; die Zitat-Korrektur (`ADR-0073`) ist nicht einschlägig (Aussagen ändern
sich).

**Index.** Zeile [`ADR-0129`](../plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md): Titel 73 Zeichen (≤ 80), Status/Datum/Datei stimmen. Die `ADR-0121`-Zeile trägt `; → ADR-0129` in
der Form der Zeilen zu `ADR-0091`/`ADR-0092` (`; → ADR-0101`). Ihr Titel wurde dabei gekürzt (kein „Bindung“ mehr) — der
Index-Titel ist keine Aussage der ADR.

## Findings

### F-1 — CI-Anker ohne Versuchsnummer; „weder belegt noch widerlegt“ steht in der Gegenwartsform

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz B (Beleg-Anker); Skill-HIGH „Beleg trägt seinen Satz nicht“ (herabgestuft, siehe `befund`)
- `pfad`: `docs/plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md:137-151` (Bullet „CI-Weg“)
- `befund`: Der Anker nennt „Lauf 36287009221 … gemessen mit `gh run view`“ ohne Versuchsnummer; `gh run view 36287009221`
  liefert heute `success` für beide Legs (Versuch 2), der Stand der ADR („Leg PostgreSQL 18 **rot**“, Replication-Tier
  `skipped`) ist nur mit `--attempt 1` auflösbar. Der Schlusssatz „Der Test ist damit in CI weder belegt noch widerlegt“ ist
  an „Stand … gegen 02:07 UTC“ gebunden und wahr (das PG17-Leg lief erst um 02:11 UTC durch), liest sich in einer
  immutablen ADR aber als bleibende Lage; belegt ist der Test in CI seit den grünen Läufen (siehe „Nachgemessen“).
  Auf MEDIUM statt HIGH gestellt: die Aussage ist zeitgestempelt und wahr, der Versuch 1 bleibt abfragbar, und die
  Entscheidung der ADR (Festlegungen 1 und 2) hängt nicht an der CI-Lage.
- `verifizierbar`: ja — `gh run view 36287009221` gegen `gh run view 36287009221 --attempt 1`
- `klasse`: Beleg-Adresse führt ohne Versuchsnummer zu einem anderen Stand

### F-2 — „Ungebunden“ (Festlegung 2 Punkt 4) ist breiter als die Messung

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 „Verfasser einer ADR“ (Aussage über alle Werte einer Menge; Verallgemeinerung von einer gefahrenen Instanz)
- `pfad`: `docs/plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md:165-169`
- `befund`: Festlegung 2 Punkt 4 nennt „die Position, ab der der Test den Stream neu startet“ ungebunden; gefahren ist ein
  Wert (`LSN(0)`, grün). R2 (`confirmed_flush_lsn + 1`) färbt den Test an PostgreSQL 17 rot (78,19 s; dieselbe Farbe
  meldet Reviewer-M10 an 17.11); die Position ist nach oben gebunden, nach unten (0) nicht. Ebenso ist „die Wartezeit auf die
  Inaktivität des Slots“ als ungebunden aus **einem** grünen Lauf der Streichung (K7) formuliert — ein Lauf zeigt bei einer
  Synchronisation keine Ungebundenheit. Die Aussage geht in die konservative Richtung (sie behauptet weniger Bindung, nicht
  mehr), deshalb LOW; die Klasse ist im Register `adr-aussage-breiter-als-ihre-messung` geführt (Zähler führt der Planner).
- `verifizierbar`: ja — Mutation R2 an einer Kopie von `sourcekeepalive_test.go:250`
- `klasse`: ADR-Aussage breiter als ihre Messung

### F-3 — „Zweite committende Transaktion ist nicht gemessen“ kollidiert mit der Messung in `ADR-0121`

- `kategorie`: LOW
- `quelle`: Skill-MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“ (hier: zwei `Accepted`-ADRs, deren Rest in Kraft bleibt) · `AGENTS.md` §3.12
- `pfad`: `docs/plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md:161-163`; Gegenstück `docs/plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md:37-48`
- `befund`: Festlegung 2 ersetzt die Grenze „einmalige Messung an PostgreSQL 18“ und schreibt: „eine zweite committende
  Transaktion sind nicht gemessen“. `ADR-0121` §Gemessen (nach `ADR-0129` „bleibt in Kraft“) berichtet genau diese Lage
  als einmalige Messung an PostgreSQL 18 (Stand-in hält die erste Transaktion, eine zweite committet 400.000 Änderungen). Der
  Satz ist wahr für **den Test**, steht aber ohne diese Einschränkung; der ersetzte Satz nannte die Messung, der neue nicht. Die
  Menge (nur der committete Test) fehlt im Satz.
- `verifizierbar`: ja — Lesen von Festlegung 2 Punkt 2 gegen `ADR-0121` §Gemessen
- `klasse`: Aussage ohne Bindung an ihre Menge

### F-4 — Aussage (a) „bindet“ der Test, aber keine Mutation der Fitness Function setzt an der bindenden Zeile an

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 „Verfasser einer ADR“ (Mutationsangabe nennt Stellen und Instanz); Skill-HIGH „Beleg trägt seinen Satz nicht“ (Probe an nicht erprobter Stelle)
- `pfad`: `docs/plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md:117-122` gegen `:232-239`
- `befund`: Festlegung 1 nennt (a) `ServerWALEnd` gleich Commit-LSN als vom Test gebunden. Die Zeile, die das prüft
  (`sourcekeepalive_test.go:229`), setzt keine der Mutationen 1 bis 9; die Mutation dazu (Reviewer-M9) schließt die ADR
  ausdrücklich als Beleg aus. Meine Probe an dieser Stelle (R1) ist an PostgreSQL 17 rot: die Aussage trägt; die Lücke
  liegt im Beleg-Text der ADR. Kein Widerspruch, keine Aktion erwartet.
- `verifizierbar`: ja — R1
- `klasse`: Bindung behauptet, Mutation der Stelle nicht genannt

### F-5 — Die Zahl „etwa 11,1 Tausend Änderungen“ ist über die vorliegenden Läufe nicht stabil

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahl mit Ursprung)
- `pfad`: `docs/plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md:77-78`
- `befund`: Die ADR führt die Zahl als „Messung dieser zwei Läufe, kein Vertrag“ (11080, 11077) — so gekennzeichnet und wahr.
  Die CI-Läufe drucken 11538 und 11541, K5 des Verifiers 11533; die Position des Keepalives hängt an Puffern und Takt des
  Hosts. Der Vorbehalt „kein Vertrag“ trägt; keine Aktion erwartet.
- `verifizierbar`: ja — CI-Logs der Jobs `108529548457`, `108531740887`
- `klasse`: Zahl über Läufe nicht stabil (Vorbehalt vorhanden)

## Negativbefunde

- geprüft, ohne Befund: **Menge und Ursprung** (ADR §Gemessen, §Festlegung 1) — die zwei Pins aus `SPEC-012` stehen mit Digest-Anker
  (`e2e.yml`, `Makefile`), Versionen 17.11/18.6 sind als Messung, nicht als Vertrag geführt; gedruckte Zeilen stimmen mit dem
  Verifikations-Report überein; jede Übernahme ist als *übernommen* gekennzeichnet.
- geprüft, ohne Befund: **Fitness-Function-Zeilen** — Stellen, Instanz und Farbe je Mutation stimmen mit K1–K10/W überein; die
  Verallgemeinerung „auch an 17 rot“ (3 bis 6) und die Version von Mutation 9 stehen als *hergeleitet*; „nicht gefahren“ ist
  benannt. Mutationsprobe an nicht erprobter Stelle (R1) und an der hergeleiteten Zeitaussage (R3): beide tragen.
- geprüft, ohne Befund: **Aussage über PostgreSQL 17 und den Stand-in-Aufbau** — die ADR schließt sie ausdrücklich aus
  („über den Aufbau des Stand-ins … wird nichts behauptet“); Festlegung 1 begrenzt die Gleichheit auf den Aufbau des Tests.
- geprüft, ohne Befund: **Alternativen A–D** — jede mit Pro/Contra; die Verwerfung von D trägt ihre Herleitung als *hergeleitet*
  und verweist die Adapter-Seite an den Unit-Test; die Bedingung des Adapters bleibt (Festlegung 1, letzter Punkt).
- geprüft, ohne Befund: **`Supersedes`/`Schärft`/Bezug/Geschichte** — Kopf nennt die zwei ersetzten Stellen wörtlich, `ADR-0120` und
  `ADR-0121` sind unberührt (`git diff --stat`), der eingetretene Trigger von `ADR-0121` ist benannt, der zweite bleibt; Index-Zeile
  und die `→ ADR-0129`-Kennzeichnung der `ADR-0121`-Zeile entsprechen der Form (Vorbild `ADR-0091`), Titel ≤ 80 Zeichen.
- geprüft, ohne Befund: **Grenzen** (`proto_version 1` — der Adapter nutzt ebenfalls `'1'`, `receive.go:160` —, ein Keepalive je
  Transaktion, keine zweite Quelltransaktion [aber F-3], zwei Nebenversionen, Zeitabhängigkeit) — mit R3 an einem Lauf bestätigt.
- geprüft, ohne Befund: **Folgepflichten mit Adresse** — Register, Träger `receive.go`/`seam_test.go` (Lokatoren am Stand stimmen)
  mit Adresse `slice-code-kommentare-bereinigung`, Slice-Plan (Planner), CI-Beleg (Planner); der Träger-Suchlauf deckt sich mit
  meiner Nachmessung; das Register trägt den Rerun bereits (`state.md`: „im ersten/zweiten Versuch“).
- geprüft, ohne Befund: **Ton und Kennungsform (§3.7)** — kein Konjunktiv über Verworfenes außerhalb der Alternativen-Tabelle, Indikativ
  über den Ist-Zustand; Kennungen stehen in Inline-Code oder als Link; die Regel „höchstens eine Kennung“ gilt für Go-Kommentare, nicht
  für ADR-Prosa.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (F-1) |
| LOW | 2 (F-2, F-3) |
| INFO | 2 (F-4, F-5) |

**Verdikt.** Keine Fixrunde an der ADR nötig; sie kann stehen bleiben. Die tragenden Aussagen (Menge = zwei Pins, gedruckte
Zeilen, Mutationszuordnung, Grenzen, Träger-Liste) halten der Nachmessung stand; die Findings sind Reichweite und Adressierung
am Rand (CI-Anker, zwei „ungebunden“/„nicht gemessen“-Sätze), keiner ändert eine Festlegung. Ein Zug am `Accepted`-Text wäre nur
als neue ADR mit `Supersedes` auf `ADR-0129` oder — für F-1 allein — als Zitat-Korrektur nach `ADR-0073` (Versuchsnummer im Anker; Entscheidung
und §Geschichte-Zeile beim Architect) zulässig und steht in keinem Verhältnis zum Gewicht der Findings. Stattdessen:

- **Planner (Closure des Slice):** den CI-Beleg mit **Versuchsnummer und Job-Kennungen** eintragen (Versuch 1: PG18 rot in „Leerlauf-Bestätigung“,
  PG17 grün mit `PASS` 40.72 s; Versuch 2: PG18 grün, `PASS` 43.20 s) — F-1; Register-Zähler der Klasse „ADR-Aussage breiter als ihre
  Messung“ für F-2; F-3 als Lesehilfe („nur der Test“) in der Closure-Notiz, wo Festlegung 2 Punkt 2 der ADR zitiert wird.
- **Architect:** bei einer späteren Berührung von `ADR-0129` (nächste Supersede oder Zitat-Korrektur) F-1 bis F-3 einarbeiten; sonst kein Zug.

Kein Rollen-Widerspruch; der Konflikt-Pfad (Modul 8) ist nicht ausgelöst.

**Hygiene.** Drei Mutationsläufe gegen einen Wegwerf-PostgreSQL 17 im eigenen Docker-Netz, sequenziell (`docker rm -fv`,
Netz entfernt); dangling Volumes vorher 36, nachher 36; kein `prune`; `git status --short` nach den Läufen leer. Die Mutationen
entstanden als `sed … Datei > Kopie` im Scratchpad, der Aufruf lief mit der Kopie als read-only-Mount über der Repo-Datei. Ein
versehentlicher Host-`python3`-Aufruf (`--version`) wurde vom Guard geblockt, ohne Wirkung. `make gates` mit diesem Report im Baum,
Exit-Code ungefiltert gelesen: Exit 0 (Lauf unmittelbar vor dem Commit dieses Reports).
