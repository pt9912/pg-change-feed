# ADR-0129: Capture — die Quellseite der Leerlauf-Bestätigung trägt ein committeter Test an beiden PostgreSQL-Pins (Supersedes ADR-0121, teilweise)

**Status:** Accepted — Supersedes [`ADR-0121`](0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
in genau zwei Stellen: Festlegung 2 („Träger der Bindung je Tier“; ersetzt durch
Festlegung 1 dieser ADR, die ihren Inhalt übernimmt und um die Quellseite ergänzt)
und den Konsequenz-Punkt „Negativ (Grenze, benannt)“ (ersetzt durch Festlegung 2
dieser ADR). Alles Übrige von `ADR-0121` bleibt in Kraft, insbesondere Festlegung 1
(Begründung der Bedingung „keine offene Transaktion“) und die Zeilen der Fitness
Function. Der erste Re-Evaluierungs-Trigger von `ADR-0121` („Ein committeter Test der
Quellseite entsteht“) ist mit dieser ADR eingetreten und eingelöst; der zweite
(Streaming großer Transaktionen) bleibt.

**Datum:** 2026-09-27

**Autor:** Architect-Agent (Modul 8), Architect-Zug zu DoD 3 von
`slice-capture-leerlauf-quellbelege`; jede Aussage über eine Menge nennt die Menge,
an der sie geprüft ist, jede Zahl ihren Ursprung (gemessen · übernommen · hergeleitet,
`AGENTS.md` §3.12). Vom Architect selbst gefahren wurde nichts; alle Messungen stammen
aus dem Verifikations-Report und sind als **übernommen** gekennzeichnet.

**Bezug:** [`LH-QA-REL-001`](../../../spec/lastenheft.md) (kein Datenverlust,
Persist-before-ACK; Haupt-Bezug), [`LH-FA-CAP-009`](../../../spec/lastenheft.md),
[`LH-QA-POR-001`](../../../spec/lastenheft.md) (beide PostgreSQL-Versionen),
[`ADR-0121`](0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md) (teilweise
superseded — Haupt-Bezug), [`ADR-0120`](0120-capture-slot-leerlauf-bestaetigung.md),
[`ADR-0030`](0030-testpyramide.md) (Tier des Tests)

**Schärft:** — (Prozess-Ergänzung der Belegform; keine Spec-Stelle ändert ihren Wortlaut)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0121` `Accepted` (unberührbar, `AGENTS.md` §3.5) führt die Aussage „die Position
eines Keepalive inmitten einer Transaktion ist die Commit-LSN dieser Transaktion, die
Quelle liefert bei dieser Gleichheit weiter“ als einmalige Messung an PostgreSQL 18
(Aufbau: ein `Capture`-Stand-in hält die erste Transaktion offen, eine zweite committet
währenddessen) und nennt als Grenze: kein committeter Test, für PostgreSQL 17 keine
Messung. Der Slice `slice-capture-leerlauf-quellbelege` hat den committeten Test
geliefert; die Grenze trifft damit für den Aufbau des Tests nicht mehr zu. Die
Berichtigung ändert §Konsequenzen und eine Festlegung von `ADR-0121`, ist also keine
Zitat-Korrektur ([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)).

### Gemessen

Alle Zeilen stammen aus dem Verifikations-Report
`verifikation-slice-capture-leerlauf-quellbelege` (Verzeichnis `docs/reviews/`) §1 und
§4 (Stand `65a63968` der Arbeit, Läufe des Verifiers); der Architect hat sie
**nicht wiederholt** (übernommen). Für diese ADR zählen nur die gedruckten Zeilen des
Verifiers; die Läufe von Implementer und Reviewer (Review-Report
`review-slice-capture-leerlauf-quellbelege`, ebenfalls `docs/reviews/`) sind hier nicht
als Beleg geführt.

**Aufbau des Tests** (`internal/adapters/driving/replication/receive/sourcekeepalive_test.go`,
`TestSourceKeepaliveInsideTransactionDeliversWholeTransaction`): ein roher
Protokoll-Client (kein Stream-Adapter), `proto_version 1` (kein Streaming), gegen eine
Instanz mit Standard-`wal_sender_timeout` (60 s) und ohne fremden Schreiber. Der
Client startet den Stream, liest 38 s nichts und antwortet nichts, während **eine**
Transaktion über 400.000 Einfügungen committet; der Walsender sendet seinen Keepalive
zwischen BEGIN und COMMIT. Der Client bestätigt dessen `ServerWALEnd`, schließt die
Verbindung inmitten der Transaktion und startet den Stream ab `confirmed_flush_lsn`
neu.

**Gedruckt, `make test-replication` (Standard-Pin, PostgreSQL 18):**
`PostgreSQL 18.6: Keepalive inmitten der Transaktion nach 11080 von 400000 Änderungen,
ServerWALEnd 0/12786CB0 gleich Commit-LSN 0/12786CB0, bestätigt 0/12786CB0; Neustart ab
0/12786CB0 lieferte 400000 Änderungen mit Commit-LSN 0/12786CB0` ·
`--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction (40.32s)`.

**Gedruckt, derselbe Lauf mit `PG_TEST_IMAGE` auf den 17er-Digest aus
`.github/workflows/e2e.yml`:** `PostgreSQL 17.11: Keepalive inmitten der Transaktion
nach 11077 von 400000 Änderungen, ServerWALEnd 0/124C79A8 gleich Commit-LSN 0/124C79A8,
bestätigt 0/124C79A8; Neustart ab 0/124C79A8 lieferte 400000 Änderungen mit Commit-LSN
0/124C79A8` · `--- PASS: … (40.37s)`. Die Zahl „nach etwa 11,1 Tausend Änderungen“
ist eine Messung dieser zwei Läufe, kein Vertrag.

**Mutationen** (Verifikations-Report §4; Instanz: Kopie der Testdatei, Aufruf des
Tests im gepinnten Toolchain-Container gegen eine eigene Instanz, Farbe die gelesene
Log-Zeile): siehe Fitness Function.

### Konstraints

- Persist-before-ACK ([`ADR-0011`](0011-persist-before-ack.md)): eine bestätigte
  Position verdeckt keinen Change, der nicht gespeichert ist.
- Ein Beleg, der nur als Ereignis eines Laufs lebt, färbt keine Regression rot
  (Register `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung`).

## Entscheidung

Wir wählen **eine neue ADR mit teilweisem `Supersedes`**, die die Träger der Bindung
je Tier um den committeten Test der Quellseite ergänzt und die Grenze „PostgreSQL 17
nicht gemessen“ für den Aufbau des Tests ersetzt. Zwei Festlegungen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | keine Änderung | `ADR-0121` nennt „kein committeter Test“ und „PostgreSQL 17 nicht gemessen“, beides trifft für den Aufbau des Tests nicht mehr zu; der nächste Leser hält die Quellseite für unbelegt |
| B — `ADR-0121` in-place ändern | ein Ort | ändert Festlegung und Konsequenz, keine Zitat-Korrektur (`AGENTS.md` §3.5) |
| **C — neue ADR mit teilweisem `Supersedes` (gewählt)** | Aussage und Grenze tragen den Stand; `ADR-0121` bleibt unberührt | eine weitere ADR |
| D — den Aufbau des Stand-ins (Adapter mit blockierendem `Capture`, zweite committende Transaktion) als committeten Test an beiden Pins nachbauen | derselbe Aufbau wie die Messung von `ADR-0121` | mehr Aufbau ohne zusätzliche Protokoll-Aussage (*hergeleitet*: der Keepalive-Text kommt vom Walsender, nicht vom Adapter); die Seite des Adapters trägt der Unit-Test |

### Festlegung 1 — Träger der Bindung je Tier (ersetzt `ADR-0121` Festlegung 2)

- **Unit-Tier** (`make test`): die Bedingung „keine offene Transaktion“ trägt
  `TestRunNoConfirmationInsideOpenTransaction` (Paket `receive`, Fake-Sitzung); ihre
  Mutation färbt genau diesen Test rot. Unverändert.
- **Store-Tier, Größe der bestätigten Position** (`make test-replication`): eine
  Position hinter dem WAL-Ende lässt den Change nach dem Neustart des Streams
  ausbleiben (Sicherheits-Test, Mutation „Position `+ 1 GiB`“ rot). Unverändert.
- **Store-Tier, Quellseite** (`make test-replication`, Phase `tier`, eigener Lauf auf
  der Standard-Instanz mit `CDC_SOURCE_KEEPALIVE_TEST_DSN`, `--- PASS`-Wächter in
  `tools/harness/run-replication-tests.sh`): `TestSourceKeepaliveInsideTransactionDeliversWholeTransaction`
  bindet **im oben beschriebenen Aufbau** drei Aussagen: (a) der Keepalive inmitten
  einer Quelltransaktion trägt als `ServerWALEnd` die Commit-LSN dieser Transaktion;
  (b) der Slot nimmt die bestätigte Position an (`confirmed_flush_lsn` erreicht sie);
  (c) nach dem Schließen der Verbindung inmitten der Transaktion liefert der
  Neustart ab `confirmed_flush_lsn` die Transaktion vollständig, mit derselben
  Commit-LSN.
- **Die Menge** dieser Aussage sind die **zwei Pins der Testdatenbank**
  ([`SPEC-012`](../../../spec/pflichtenheft.md): Major-Versionen 17 und 18; die
  Digests stehen in `.github/workflows/e2e.yml`, der 18er zugleich als `PG_TEST_IMAGE`
  im `Makefile`), gedruckt als PostgreSQL 17.11 und 18.6. Belegt (übernommen, siehe
  §Gemessen): (a), (b) und (c) an beiden. **Keine** Aussage gilt damit für eine andere
  Nebenversion, und über den Aufbau des Stand-ins von `ADR-0121` an PostgreSQL 17 wird
  nichts behauptet.
- `ADR-0121` Festlegung 1 nennt die Gleichheit „Keepalive-Position gleich Commit-LSN
  und die Quelle liefert weiter“ für PostgreSQL 18 im Stand-in-Aufbau; im Aufbau des
  Tests gilt sie jetzt an beiden Pins. Die Regel von `ADR-0120` Festlegung 1
  Punkt 3 und die Begründung von `ADR-0121` Festlegung 1 (Invariante auf der Seite
  des Adapters, unabhängig vom Verhalten der Quelle) bleiben unverändert: der Test
  belegt das Verhalten der Quelle, er macht die Bedingung des Adapters nicht
  entbehrlich.
- **CI-Weg (hergeleitet, nicht belegt):** `.github/workflows/e2e.yml` führt den Schritt
  „Replication-Tier (go test ./... und Slot-Reserve)“
  (`bash tools/harness/run-replication-tests.sh tier`) je Matrix-Leg (PostgreSQL 17
  und 18, `PG_TEST_IMAGE` auf Job-Ebene) **hinter** dem Schritt „Compose-Integrationstest
  (Black-Box-E2E)“ aus; ist dieser rot, steht der Replication-Tier-Schritt des Legs auf
  `skipped`. Belegt ist der Test in CI erst durch einen grünen Lauf dieses Schritts je
  Leg. Stand des Ankers, des `e2e`-Laufs 36287009221 zum Commit `e4b77a05` (gemessen mit
  `gh run view`, 2026-09-27 gegen 02:07 UTC): Leg PostgreSQL 18 **rot** im
  Compose-Integrationstest, Phase „Leerlauf-Bestätigung“ (eine Phase, die
  `slice-capture-leerlauf-quellbelege` nicht ändert; der Feed-Container endete mit
  „WAL-Rückstand 15238216 Bytes über Fehlerschwelle 8388608 Bytes“), Schritt
  „Replication-Tier“ dort `skipped`; Leg PostgreSQL 17 im Compose-Integrationstest
  grün, die Schritte dahinter zu diesem Zeitpunkt nicht abgeschlossen. Der Test ist
  damit in CI weder belegt noch widerlegt. Der Workflow ändert sich nicht (`AGENTS.md`
  §3.10 greift nicht); `e2e.yml` ist nicht blockierend.

### Festlegung 2 — Reichweite und Grenzen (ersetzt den Konsequenz-Punkt „Negativ (Grenze, benannt)“ von `ADR-0121`)

Der committete Test bindet die Quellseite in diesem Aufbau; die Grenze ist nicht mehr
„einmalige Messung“ und „PostgreSQL 17 nicht gemessen“, sondern:

1. **`proto_version 1`, kein Streaming großer Transaktionen.** Mit Streaming (`pgoutput`
   Protokoll ab Version 2) ist die Lage nicht gemessen; der zweite Trigger von
   `ADR-0121` bleibt.
2. **Ein Keepalive je Transaktion im Fenster** des Tests und **keine gleichzeitige
   zweite Quelltransaktion.** Ein zweiter Keepalive und eine zweite committende
   Transaktion sind nicht gemessen.
3. **Nur die zwei Nebenversionen 17.11 und 18.6.**
4. **Ungebunden** (gemessen: die Mutation bleibt grün, Verifikations-Report §4): die
   Position, ab der der Test den Stream neu startet (`START_REPLICATION` ab 0 nimmt
   den Stand des Slots) und die Wartezeit auf die Inaktivität des Slots vor dem
   Neustart. Die Aussage (c) hängt an `confirmed_flush_lsn`, nicht an der Startposition
   des zweiten Laufs.
5. **Zeitabhängigkeit.** Der Test wartet 38 s ab `START_REPLICATION`; ein Runner, auf
   dem die 400.000 Einfügungen länger als diese Zeit brauchen, beendet den Test laut
   („COMMIT … ohne Keepalive“), nicht mit einer falschen Aussage (*hergeleitet* aus
   dem Code des Tests; gemessen ist nur die Gegenrichtung: Stille 70 s statt 38 s
   färbt den Test rot, PostgreSQL 18, Verifikations-Report §4).

## Konsequenzen

- Positiv: die Quellseite hat einen committeten Wächter an beiden Pins; der Beleg lebt
  nicht mehr als Ereignis eines Laufs, und jede Grenze steht mit ihrer Menge.
- Negativ: Laufzeit. Der Test braucht je Lauf etwa 40 s (gemessen, 40,32 s und
  40,37 s, Verifikations-Report §1); in `e2e.yml` kommt dieser Anteil je Matrix-Leg
  hinzu (hergeleitet aus der gemessenen Testdauer).
- Negativ: `ADR-0121` trägt die ersetzten Sätze weiter im Text; wer sie liest, liest
  diese ADR über den Index (Titel „→ ADR-0129“).
- **Folgepflichten** (mit Adresse; keine neue Entscheidung nötig):
  1. **Register** `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung`: der Ausgang des
     ersten Liefer-Punkts („Quellseite als committeter Test an beiden Pins“) ist
     lokal erreicht (§Gemessen); `state.md` trägt beim Planner-Nachzug der Closure
     `Ausgang` und Zähler (Planner, Closure von `slice-capture-leerlauf-quellbelege`).
  2. **Träger, die die alte Grenze noch nennen** (gemessen mit
     `git grep -n -E 'einmalige Messung|PostgreSQL 17 (ist )?nicht gemessen|für PostgreSQL 17 liegt'`
     ohne `docs/reviews` und `.harness/baseline` am Stand `e4b77a05`): `ADR-0121`
     §Konsequenzen (ersetzt durch Festlegung 2 dieser ADR, bleibt als `Accepted`-Text
     stehen); die Register-Dateien des Eintrags (`observation.md`, `state.md`,
     `evidence/`) — Planner bei der Closure; der Plan dieses Slice (§3 Suchlauf-Feld,
     §6) — Planner; die Records unter `done/` (`slice-backfill-slot-leerlauf-bestaetigung`,
     `welle-backfill-bestand*`) — bleiben (Records). **Kein** Träger in `spec/`,
     `docs/user/`, `harness/`, `README.md`, `AGENTS.md` (`git grep -c` mit demselben
     Muster an demselben Stand ohne Treffer). `ADR-0119` Zeilen 48 und 89
     („PostgreSQL 17: nicht gemessen“) und `ADR-0125` Zeile 144 haben einen anderen
     Gegenstand (Wirkung der Lesesperre bzw. d-migrate) und bleiben draußen.
  3. **Code-Kommentare** mit dem Begründungssatz, den `ADR-0121` Festlegung 1 ersetzt
     hat („das WAL-Ende liegt dann hinter Nachrichten, die noch nicht gespeichert
     sind“): das Godoc von `confirmIdle` in
     `internal/adapters/driving/replication/receive/receive.go` (Zeilen 467 bis 469)
     und das Godoc von `TestRunNoConfirmationInsideOpenTransaction` in
     `internal/adapters/driving/replication/receive/seam_test.go` (Zeilen 694 bis 697).
     Adresse: `slice-code-kommentare-bereinigung` (Tranche Paket `receive`; die neue
     Fassung nennt die Regel und `ADR-0121` als einen Anker, nicht die Begründung in
     eigenen Worten). Diese ADR ändert keinen Produktionscode.
  4. **Slice-Plan** `slice-capture-leerlauf-quellbelege` (Planner): DoD 1 auf den
     Aufbau des Tests ziehen und den Verifikations-Report als „Zu belegen durch“
     nennen; DoD 3 auf diese ADR und ihre Index-Zeile; die `diff`-Stände des
     Suchlauf-Felds auf Commit-Stände setzen, sobald diese ADR committet ist (Muster 1
     trifft sie: „Quellseite“, „inmitten … Transaktion“, „einmalige Messung“).
  5. **CI-Beleg und Befund:** das Ergebnis des Schritts „Replication-Tier“ je Leg trägt
     der Planner in die Closure ein (Anker: Festlegung 1, CI-Weg). Der rote Leg
     PostgreSQL 18 in der Phase „Leerlauf-Bestätigung“ ist ein Befund an einer
     Phase, die dieser Slice nicht ändert; Zuordnung (Register
     `BEO-PGC/test-integration-retention-timing-flake` als Kandidat, oder ein
     eigener Vorgang) entscheidet der Planner mit den Läufen dieses Slice als Beleg.
     Diese ADR trifft keine Aussage darüber, ob die Phase auf gehosteten Runnern
     zeitabhängig ist.

## Fitness Function (falls maschinell prüfbar)

Die Zeilen von `ADR-0121` bleiben; diese ADR fügt hinzu. **Ursprung der Farben:** alle
Mutationen sind vom Verifier gefahren und aus seinem Report **übernommen** (der
Architect hat keine gefahren); „hergeleitet“ steht, wo die Aussage über die gefahrene
Instanz hinausgeht.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL 17.11 und 18.6 (Tier `test-replication`, Phase `tier`), `TestSourceKeepaliveInsideTransactionDeliversWholeTransaction` | Menge: die zwei Pins (§Festlegung 1). Der Test läuft grün an beiden (gedruckt, §Gemessen); welche Verletzungen er rot färbt, sagen die Mutationen 1 bis 6 (die gefahrenen Instanzen, nicht „jede“); ein Lauf ohne `--- PASS` des Tests färbt das Tier rot (Wächter im Runner) | `make test-replication` |
| dieselbe Zeile, Mutation 1: bestätigte Position `+ (1 << 30)` (Zeile `confirmedPosition := …`, **eine** Stelle) | erprobt an **PostgreSQL 18 und 17**; Instanz: Kopie der Testdatei, Go-Test im Toolchain-Container gegen eine eigene Instanz; Farbe **rot** an beiden (je 78,18 s): „Der Neustart des Streams lieferte die Transaktion nicht bis zum COMMIT (0 Änderungen gelesen) … context deadline exceeded“ | `make test-replication` |
| Mutation 2: der Vergleich der gelieferten Commit-LSN mit der des Keepalives (Test-Zeile 263) um eins verschoben (**eine** Stelle; Wortlaut der Mutation im Report: `commitLSN != finalLSN+1`) | erprobt an **PostgreSQL 18 und 17**, Instanz wie Mutation 1; Farbe **rot** an beiden: „Neustart lieferte die Transaktion mit BEGIN-Kopf … und COMMIT …, erwartet die Commit-LSN …“ | `make test-replication` |
| Mutationen 3 bis 6: Prüfung der Änderungszahl auf `sourceKeepaliveRows-1`; Schranke `changesSoFar >= 5000` statt `>= sourceKeepaliveRows`; Stille 70 s statt 38 s; `confirmRaw(…, confirmedPosition-1)` (je **eine** Stelle) | erprobt an **PostgreSQL 18** allein (Instanz wie Mutation 1); Farbe **rot** je Mutation („400000 von 400000“ / „liegt nicht inmitten der Transaktion“ / „unexpected EOF“ / „confirmed_flush_lsn … erreicht die bestätigte Position … nicht“). An PostgreSQL 17 sind sie **nicht** gefahren; die Aussage „auch an 17 rot“ ist *hergeleitet* (derselbe Test, dieselbe Zeile), nicht erprobt | `make test-replication` |
| Mutationen 7 und 8: Startposition des Neustarts `LSN(0)`; Aufruf `awaitSlotInactive` gestrichen | erprobt an **PostgreSQL 18** (Mutation 8: ein Lauf); Farbe **grün**: der Test bindet beide **nicht** (Festlegung 2, Punkt 4); keine Aussage dieser ADR hängt an ihnen | `make test-replication` |
| Runner `tools/harness/run-replication-tests.sh tier`, Mutation 9: `-run`-Muster um ein Zeichen verfälscht | erprobt (Kopie des Repos, Runner-Aufruf `tier`; Version des Laufs im Report nicht genannt, der Default-Pin 18 ist *hergeleitet*); Farbe **rot**: „ist nicht als PASS gelaufen“, Exit 1 | `make test-replication` |

Nicht in der Fitness Function, weil nicht gefahren: Mutationen 3 bis 9 an PostgreSQL 17;
jeder grüne Lauf des Tests auf GitHub (Festlegung 1, CI-Weg).

## Re-Evaluierungs-Trigger

- **Eine neue PostgreSQL-Hauptversion, ein Wechsel eines Pins oder `pgoutput` mit
  Streaming großer Transaktionen:** der Test läuft im Tier am neuen Pin; ist er rot,
  ist Festlegung 1 (a) bis (c) für diesen Pin neu zu messen. Die Versionsnummern
  17.11 und 18.6 sind Messung, kein Vertrag.
- **Ein committeter Test des Aufbaus „Stand-in mit zweiter committender Transaktion“
  oder mit Streaming entsteht:** Festlegung 2, Punkte 1 und 2 um dessen Tier ergänzen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-27 | Accepted — Ergänzung von `ADR-0121` Festlegung 2 nach der Verifikation von `slice-capture-leerlauf-quellbelege` | [`LH-QA-REL-001`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0129` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
