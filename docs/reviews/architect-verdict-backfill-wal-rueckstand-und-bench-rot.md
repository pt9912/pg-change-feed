# Architect-Verdikt: Backfill — WAL-Rückstand des Capture-Slots und rotes `make bench` (`LH-QA-PER-001`)

**Rolle:** Architect (Modul 8)

**Anlass:** Zwei Befunde aus dem Abschnitt „Befunde der Messung“ von
`slice-backfill-bench-richtgroesse`, vom Implementer als Übergabe an Planner und
Architect gereicht (kein Umfang des Slice): (1) ein Backfill-Run treibt den
WAL-Rückstand des Capture-Slots; ein dritter Run über 1.000.000 Zeilen im selben
CDC-Speicher überschritt die Fehlerschwelle von 1 GiB, der Feed-Container endete
mit Ausgang 1. (2) `make bench` endet am Parent-Stand `933ab054` bei
`tools/bench-source-impact.sh` (`LH-QA-PER-001`) mit 87,5 % / 93,9 % / 95,8 %
Overhead gegen die 35-%-Schwelle; der Implementer vermutete Host-Last, ohne es
zu untersuchen.

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der Implementer-Lauf)

**Datum:** 2026-09-25

**Bezug:** [`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`LH-QA-REL-001`](../../spec/lastenheft.md),
[`LH-QA-REL-003`](../../spec/lastenheft.md),
[`LH-QA-PER-001`](../../spec/lastenheft.md),
[`SPEC-009`](../../spec/pflichtenheft.md),
[`SPEC-013`](../../spec/pflichtenheft.md),
[`SPEC-025`](../../spec/pflichtenheft.md),
[`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md),
[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md),
[`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md); die Slices sind
als Kennung genannt, nicht als Pfad-Link (ein Slice wechselt die Lifecycle-Ablage,
ein Pfad-Link bräche mit)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (Accepted,
  ergänzt `ADR-0007` und `ADR-0049`, überschreibt keine `Accepted`-ADR); Index-Zeile in
  `docs/plan/adr/README.md`
- dieses Dokument

Die Pläne, die Spec, das Handbuch und der Code bleiben in diesem Zug unberührt.

---

## Verdikt

**1 — Befund 1 ist ein Defekt im Capture-Pfad, kein Fehler der Richtgröße; vor der
Welle-Closure zu beheben (Verdikt 2/3: die Entscheidungslage hatte eine Lücke, die
Folge-Entscheidung `ADR-0120` schließt sie).**

- **Ursache (gemessen):** der Slot bestätigt nur persistierte Commits; WAL ohne Inhalt für die
  Publication erreicht den Feed nie als Commit, `confirmed_flush_lsn` bleibt stehen, der
  Rückstand (WAL-Ende minus `confirmed_flush_lsn`) wächst um jedes solche WAL. Der Backfill ist
  nur die Last, die es sichtbar macht: ein Schreiber auf eine **nicht aktivierte** Tabelle erzeugt
  ohne Backfill dieselbe Wirkung (X1/X2, §Befundlage). Der Standard-Empfänger von PostgreSQL
  (Abonnement) bestätigt dasselbe WAL: Rückstand seines Slots 0,0 MiB in jeder Zeile.
- **Nicht ausschlaggebend (gemessen):** die Größe der Schreibtransaktion. Eine Transaktion (X1)
  und 200 Transaktionen (X2) über dieselben 200.000 Zeilen liefern 175 und 177 B/Zeile und denselben
  Zuwachs. Ausschlaggebend ist die Slot-Bestätigung.
- **Vertragsbruch:** nicht die Warnung. Die Warnung (1) von [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3 gilt der Kopierdauer und dem Snapshot; ihre Zusage („Warnung vor
  Überschreitung“) ist über die Zeit definiert und wird nicht gebrochen. Gebrochen ist der
  Konstraint aus [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md): „Der
  Capture-kritische Pfad bleibt unberührt; der Backfill läuft daneben“ — ein Run beendet den
  Capture-Prozess. Ein Betreiber, der die Warnung beachtet, läuft trotzdem (bei Breite 74 B ab
  etwa 1,5 Mio. Zeilen, abgeleitet) in den Abbruch; bei breiten Zeilen früher (≈ 0,66 Mio. bei
  1.008 B).
- **Abhilfe (tragend): Option H, Leerlauf-Bestätigung** (`ADR-0120`): im Leerlauf des Streams
  (keine Transaktion offen) bestätigt die Application das `ServerWALEnd` der Keepalive-Nachricht,
  ohne Persistenz; die Entscheidung bleibt in der Application (`ADR-0007`).
  `SPEC-013`-Schwellen, Metrik-Formel und Richtgröße bleiben.
- **Verworfen (mit Messung):** (a) Richtgröße an die WAL-Schwelle koppeln — die Bytes je Zeile
  hängen an Breite (644 → 1.620 B) und Bestand (643 → 730 B), sind beim Antrag geschätzt und
  behebt nur den Backfill; (b) blockweise Commits — ohne Wirkung auf den Rückstand (X1/X2) und
  bricht die Atomarität; (c) Schwelle für Runs anheben — Lockerung ohne Ursachenbehebung
  (`AGENTS.md` §3.6); (d) Vorabprüfung im Antrag — dieselbe Schätzung wie (a); (e) nur
  dokumentieren — bricht den Konstraint von `ADR-0111`.

**2 — Befund 2 ist eine Eigenschaft dieses Hosts, kein Regress und keine CPU-Last: die
Entscheidung `ADR-0104` gilt, die Vermutung „Host-Last“ trifft die Ursache nicht
(Verdikt 1: der Plan hat falsch behauptet).**

- **Ursache (gemessen):** die Latenz des Festschreibens (`fsync`) auf dem Datenträger dieses Hosts.
  `pg_test_fsync` im PostgreSQL-Container: `fdatasync` 2.956 µs/Operation; die 5.000
  Einzeltransaktionen des Skripts ohne CDC brauchen 15.870 ms, also 3,17 ms je Commit
  (abgeleitet) — und mit CDC 30.885 ms (94,6 %). Der Feed schreibt je Quelltransaktion eine
  zweite, fsync-gebundene Transaktion in dieselbe Instanz; bei 3 ms je Flush verdoppelt sich die
  Latenz, bei 0,2 bis 0,3 ms (die Zeitbasis der 25,0 %/28,2 % aus `ADR-0104`, abgeleitet aus
  „200–300 ms für 1000 Transaktionen“) bleibt der Aufschlag im Band.
- **Kontrolle (gemessen):** dasselbe Skript in einer Wegwerf-Kopie mit
  `synchronous_commit=off` als einzigem Unterschied: ohne CDC 223 ms, mit CDC 208 ms
  (−6,7 %), Exit 0. Ohne Flush-Wartezeit ist kein Aufschlag durch den Feed messbar.
- **Host-Last widerlegt (gemessen):** `top` vor dem Lauf 91,5 % idle, Load Average 1,87 auf 20
  CPU; `docker ps` ohne Container; ein Fremdprozess (Suricata) belegt einen Kern.
- **Kein Bisect** gegen den Stand von `ADR-0104`; ein Beitrag späterer Code-Änderungen ist
  damit nicht ausgeschlossen, zur Erklärung aber nicht nötig (die Kontrolle erklärt 87 bis 96 %
  vollständig aus der Flush-Latenz). Die drei Läufe des Implementers (übernommen) und der Lauf
  dieses Zugs liegen zwischen 87,5 % und 95,8 %: ein stabiles Band dieses Hosts, kein
  Ausreißer.
- **Aussagekraft der 35-%-Schwelle:** sie ist eine Aussage über die Flush-Latenz des
  Messhosts, nicht nur über die Software. `SPEC-025` und `ADR-0104` bleiben; **keine**
  Skript-Änderung, **keine** Schwellen-Änderung (`AGENTS.md` §3.6: nie ohne ADR), **keine**
  neue ADR. `make bench` ist kein Gate (`ADR-0054`); das rote Ergebnis blockiert nichts.
- **Kein Register-Eintrag:** der Befund verschwindet nicht spurlos — dieses Verdikt hält ihn —
  und ist ohne Handlungsbedarf, solange `make bench` kein Gate ist. **Re-Evaluierungs-Trigger:**
  eine zweite Umgebung (CI, anderer Host) führt `make bench` aus und ist rot oder
  ungewöhnlich grün → Folge-ADR zur Bezugsgröße der Schwelle (etwa relativ zur gemessenen
  Flush-Latenz).

**3 — Zuschnitt:**

- **Befund 1:** **eigener Slice** in `welle-backfill-bestand`
  (`slice-backfill-slot-leerlauf-bestaetigung`, Arbeitsname; Namensregel
  [`MR-002`](../../harness/conventions.md)), Start nach dem `done/` von
  `slice-backfill-bench-richtgroesse` (WIP-Limit 1), **vor** der Welle-Closure. Keine Fixrunde im
  laufenden Slice: der Eingriff berührt die Empfangs-Schleife, den `CaptureService` und einen
  Port des Capture-kritischen Pfads und braucht Review und Verifier für sich;
  `slice-backfill-bench-richtgroesse` ändert keine dieser Stellen und bleibt wie geplant.
  **Neue ADR: ja, `ADR-0120`** — eine Ergänzung, keine Ablösung: `ADR-0007` (Application
  entscheidet) und `ADR-0049` (Schwellen) bleiben in Kraft, `ADR-0111` und `ADR-0113` werden nicht
  geändert.
- **Befund 2:** kein Slice, keine ADR. Der Plan von `slice-backfill-bench-richtgroesse`
  bindet den DoD-Punkt „realer `make bench`-Lauf“ an das, was erreichbar ist: der Lauf von
  `tools/bench-backfill.sh` einzeln (Exit 0, gedruckte Zahlen) trägt die Messung, der `make
  bench`-Lauf endet mit Exit 2 an `LH-QA-PER-001` mit der Ursache aus diesem Verdikt (das ist die
  Feststellung des Planners bei der Closure, keine Änderung dieses Zugs).

---

## Befundlage

### Befund 1 — Messung

Wegwerf-Skript, nicht committet (der Slice wiederholt jede Zeile als Test). Umgebung: Linux
6.8.0-139-generic, Docker 29.8.1, 20 CPU; PostgreSQL 18.6 (Pin `PG_TEST_IMAGE` im `Makefile`,
`wal_level=logical`, sonst die Startparameter von `tools/bench-lib.sh`); Feed-Image `:dev` aus
`b98b5af7`; Container samt Volumes am Lauf-Ende entfernt (`docker rm -fv`). Rückstand =
`pg_current_wal_lsn()` minus `confirmed_flush_lsn`. Die Kontrolle ist ein **Abonnement**
(`CREATE SUBSCRIPTION`) mit eigenem Slot auf einer anderen Tabelle derselben Instanz. Die
gedruckten Zeilen stehen im Lauf `wal-verdikt[20260925T004731Z]`; die vollständigen Tabellen und
die Zusatzläufe (Leerlauf `idle-verdikt[20260925T005308Z]`, Spill `spill-verdikt[20260925T005729Z]`)
stehen in `ADR-0120` §Gemessen.

| Zeile | Schreiblast | WAL der Instanz | Rückstand des Feed-Slots | Abonnement-Slot |
|---|---|---|---|---|
| X1 | 200.000 Zeilen, **nicht aktivierte** Tabelle, eine Transaktion | 33,5 MiB (175 B/Zeile) | +33,5 MiB, nach 40 s unverändert | 0,0 MiB |
| X2 | 200.000 Zeilen, nicht aktivierte Tabelle, 200 Transaktionen | 33,8 MiB (177 B/Zeile) | +33,8 MiB | 0,0 MiB |
| — | ein Commit auf einer aktivierten Tabelle | — | 0,0 MiB nach 3 s | — |
| R1/R2/R3 | Backfill-Run, 200.000 Zeilen à 74 B, gestapelt | 122,8 / 122,6 / 139,2 MiB (644 / 643 / 730 B/Zeile) | Zuwachs gleich dem WAL des Runs, nach dem Run unverändert | 0,0 MiB |
| W | Backfill-Run, 20.000 Zeilen à 1.008 B | 30,9 MiB (1.620 B/Zeile) | Zuwachs gleich dem WAL des Runs | 0,0 MiB |

**Was sich daraus tragen lässt, jeweils mit Beleg:**

- *Slot-Bestätigung, nicht Transaktionsgröße:* X1 gegen X2; der Commit auf der aktivierten
  Tabelle senkt den Rückstand sofort (Zeile „—“; dasselbe der Befund des Implementers, 776 →
  2 MiB, übernommen).
- *Nicht an den Backfill gebunden:* X1/X2 laufen ohne Run.
- *Bytes je Zeile hängen an Breite und Bestand:* 644 B (Breite 74 B) gegen 1.620 B (Breite
  1.008 B); 643 B (Bestand 200.000) gegen 730 B (Bestand 400.000). Die Handbuch-Zahl 735 B/Zeile
  (übernommen) liegt in diesem Band.
- *Leerlauf ist kein Problem:* 48,6 B/s (abgeleitet aus 538.280 → 546.112 B über 161 s).
- *Spill des Walsenders:* 79 MB `spill_bytes` bei einem Run über 200.000 Zeilen (≈ 414 B/Zeile,
  abgeleitet) — eine Eigenschaft der Ein-Transaktions-Form, nicht dieser Abhilfe
  (`ADR-0111`-Trigger).
- *Nicht gemessen:* der Neustart-Kreislauf des Containers nach dem Abbruch (hergeleitet);
  die Wirkung der Abhilfe (Code fehlt, das ist der Slice); PostgreSQL 17.

### Optionen (a) bis (e) der Auftragsfrage, an der Messung

| Option | Trägt sie? | Beleg |
|---|---|---|
| (a) Richtgröße an die WAL-Grenze koppeln | nein | Bytes je Zeile 643 bis 1.620 B je nach Breite und Bestand; beim Antrag nur geschätzt (`BEO-PGC/geschaetzter-wert-als-grenze`); X1 bleibt |
| (b) Block-Commits | nein | X1 gegen X2: kein Unterschied; bricht `ADR-0111` Teilfrage 4 |
| (c) Schwelle für Runs anheben/ausnehmen | nein | Lockerung (`AGENTS.md` §3.6); X1 bleibt |
| (d) Vorabprüfung im Antrag | nein | dieselbe Schätzung wie (a); die freie Reserve hängt an der Last anderer Schreiber |
| (e) nur dokumentieren | nein | ein Run beendet den Capture-Prozess; Konstraint aus `ADR-0111` |
| **(f) Leerlauf-Bestätigung** (nicht in der Auftragsliste) | **ja** | Abonnement-Slot in jeder Zeile 0,0 MiB; `ADR-0120` |

**Bis der Slice geliefert ist:** zwei vorhandene Abhilfen für den Betreiber — ein Commit auf
einer aktivierten Tabelle und die Fehlerschwelle `wal_retention_error_bytes` der
Konfigurationsdatei (`SPEC-013`-Override, [`Pflichtenheft`](../../spec/pflichtenheft.md)
Konfigurationstabelle). Der Planner nimmt den Hinweis in das Handbuch von
`slice-backfill-bench-richtgroesse` auf (ein Satz im Punkt „WAL-Rückstand des Capture-Slots“,
kein neuer Umfang).

### Befund 2 — Messung

| Größe | Wert | Ursprung |
|---|---|---|
| `tools/bench-source-impact.sh`, N=5.000, 5 Läufe je Phase, Lauf dieses Zugs | ohne CDC 15.870 ms, mit CDC 30.885 ms, 94,6 % | gemessen, Exit 1 des Skripts |
| dieselben Läufe des Implementers | 87,5 % / 93,9 % / 95,8 % | übernommen (Plan §2 DoD) |
| `pg_test_fsync -s 2` im `postgres:18-alpine`-Container | `fdatasync` 2.956 µs, `open_datasync` 5.485 µs, `fsync` 8.619 µs | gemessen |
| Commit-Latenz ohne CDC | 15.870 ms / 5.000 = 3,17 ms | abgeleitet |
| Zeitbasis von `ADR-0104` | 200–300 ms je 1.000 Transaktionen = 0,2–0,3 ms | übernommen aus `ADR-0104` §Kontext (2), abgeleitet |
| Kontrolle `synchronous_commit=off` (Wegwerf-Kopie des Skripts) | 223 ms gegen 208 ms (−6,7 %), Exit 0 | gemessen |
| Host vor dem Lauf | 91,5 % idle, Load Average 1,87, kein Container | gemessen (`top`, `docker ps`) |

**Speicher und Volumes** (`free -m`, Zeile „Speicher“): benutzt 11103 MB (Beginn) → 11081 MB
(nach dem letzten Lauf); dangling Volumes (`docker volume ls -qf dangling=true | wc -l`): 34
vor, 34 nach den Läufen; kein `prune`.

---

## Folge-Arbeit

Jede Pflicht hat einen Träger; die Pläne ändert dieses Dokument nicht — der Planner zieht sie
nach (`AGENTS.md` §3.13).

| Träger | Nachzug | Herkunft |
|---|---|---|
| `slice-backfill-slot-leerlauf-bestaetigung` (neu, Planner) | Umsetzung nach `ADR-0120` Festlegung 1 bis 3; Start nach `done/` von `slice-backfill-bench-richtgroesse`; in die Slice-Tabelle und die Abhängigkeiten von `welle-backfill-bestand`; Closure-Trigger der Welle setzt den Slice voraus | Verdikt 1, 3 |
| dieser Slice, Tests | **Unit** (`make test`): Paket `receive` mit Fake-Sitzung — Keepalive im Leerlauf mit `ServerWALEnd` hinter `lastAcked` bestätigt (auch ohne `ReplyRequested`), Keepalive inmitten einer Transaktion nicht, kein Rückschritt; Mutation „Bedingung offene Transaktion entfernt“ rot; Paket `usecase/capture`: nur `ReplicationAckPort`, nie `PersistTransaction`, Portfehler als Klasse `replication`. **Store** (`make test-replication`): X1-Form gegen den laufenden Stream, Rückstand unter einem Zehntel der Last, auch bei Standard-`wal_sender_timeout`; **Sicherheits-Test:** offene Transaktion auf einer veröffentlichten Tabelle, Leerlauf-Bestätigung, Commit, Stream-Neustart → der Change wird geliefert, Mutation (Bestätigung inmitten der Transaktion) → Change fehlt. **E2E** (`make test-integration`): Run über mehr WAL, als die Fehlerschwelle (Override) trägt → Feed läuft, Run `completed`. **Bench** (`make bench`, kein Gate): Rückstand je Run gedruckt, Spitze der größten Stufe unter der Warnschwelle | `ADR-0120` Fitness Function |
| dieser Slice, Spec | [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) Schritt „ACK Source“ und `SPEC-009` (Zeile `cdc_wal_retention_bytes`) je ein Satz; Architektur-Sicht: der Leerlauf-Weg ohne ADR-/Slice-Bezug (`AGENTS.md` §3.4) | `ADR-0120` Folgepflicht 2 |
| dieser Slice, Handbuch | „WAL-Rückstand prüfen“, §4 Punkt „WAL-Rückstand des Capture-Slots“, §Grenzwerte Punkt „Backfill, WAL-Rückstand des Capture-Slots“ (der Satz „Diese Schwelle kann bei weniger Zeilen greifen als die Richtgröße“ entfällt); gehaltenes WAL und Spill als Grenze der Ein-Transaktions-Form | `ADR-0120` Folgepflicht 3 |
| dieser Slice, Bench | `tools/bench-backfill.sh`: die Zeile „WAL-Rückstand-Schwellen (abgeleitet)“ entfällt; `harness/targets/bench-backfill.md` | `ADR-0120` Folgepflicht 4 |
| dieser Slice, Suchlauf | `ADR-0080` §Entscheidung Punkt 5 und Kommentar/Test von `handleCopyData`/`standbyStatus` („Keepalive-Antwort mit der letzten bestätigten Position“) | `AGENTS.md` §3.13, `ADR-0120` Folgepflicht 5 |
| `slice-backfill-bench-richtgroesse` (Plan) | §3 „Befunde der Messung“: Punkt 1 auf „beantwortet durch `ADR-0120` und den Slice `slice-backfill-slot-leerlauf-bestaetigung`“, Punkt 3 auf „Ursache Flush-Latenz des Hosts, Verdikt“ ziehen; DoD-Punkt „realer `make bench`-Lauf“ auf den erreichbaren Beleg (Skript einzeln, Exit 0; `make bench` Exit 2 mit Ursache) setzen; Handbuch: ein Satz Abhilfe | Verdikt 2, 3 |
| Spec / Handbuch dieses Zugs | keine Änderung; die Nachzüge tragen die genannten Träger | — |

**Reihenfolge:** `slice-backfill-bench-richtgroesse` schließt wie geplant →
`slice-backfill-slot-leerlauf-bestaetigung` → Welle-Closure; `slice-backfill-sdk-origin` ist von
beidem unabhängig.

---

## Auftraggeber-Frage

**Keine offen.** `ADR-0120` steht `Accepted` (Vollmacht zum Architect-Zug, jede Tatsachenaussage
an einer gedruckten Messzeile oder als hergeleitet gekennzeichnet, `AGENTS.md` §3.12). Hält der
Auftraggeber statt des eigenen Slice eine Fixrunde für richtig, ist das eine Planner-Entscheidung
ohne Wirkung auf die ADR; die Sicherheits-Tests der Fitness Function gelten dann unverändert.
