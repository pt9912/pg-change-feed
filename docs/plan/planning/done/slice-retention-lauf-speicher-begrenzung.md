# Slice retention-lauf-speicher-begrenzung: Der Retention-Lauf liest Kandidaten seitenweise ohne Row Images — der Speicher des Feed-Containers hängt an der Seitengröße, nicht an der Zahl der Changes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er hat keine Kante zu einer offenen Welle und ist
von den Slices der Welle [welle-transformationen](../welle-transformationen.md)
unabhängig.

**Bezug:** [`LH-FA-RET-002`](../../../../spec/lastenheft.md),
[`LH-FA-RET-003`](../../../../spec/lastenheft.md) und
[`LH-FA-RET-004`](../../../../spec/lastenheft.md) (Retention: Aufbewahrung,
Mindestalter, Consumer-Positionen),
[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des Bestands, der
`cdc.change` in einem Zug füllt),
[`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
(die Entscheidung dieses Slice: Kandidaten seitenweise ohne Row Images, `Accepted`),
[`ADR-0014`](../../adr/0014-retention-domain-policy.md) (Retention als Domain
Policy, unverändert in Kraft),
[`ADR-0009`](../../adr/0009-change-store-outbound-port.md) (ein Fähigkeits-Port
für den Change Store),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) (sein dritter
Re-Evaluierungs-Trigger ist der Anlass),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
(Richtgröße; sein Trigger „Eine Messung liegt vor“), Architect-Verdikt
[`architect-verdict-retention-lauf-speicher-begrenzung`](../../../reviews/architect-verdict-retention-lauf-speicher-begrenzung.md),
Messbericht
[`messbericht-slice-backfill-speicher-untersuchung`](../../../reviews/messbericht-slice-backfill-speicher-untersuchung.md).

**Berührte Spec-Stellen:** [`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md)
(Safe Watermark der Retention) — **geändert** im Träger-Nachzug (§3): Punkt 4
kommt hinzu; die Bereinigungsmenge bleibt dieselbe, geändert wird die Art, sie
zu bestimmen. [`SPEC-022`](../../../../spec/pflichtenheft.md) (Zeile „Noch nicht
begrenzt“: `GET /changes` ohne `limit` liest unbegrenzt) — gelesen, nicht
geändert: die Aussage gilt dem Lesezugriffsweg, nicht dem Retention-Pfad.
[`SPEC-005`](../../../../spec/pflichtenheft.md) (Zeile 72: große Transaktionen
nicht unbegrenzt im RAM) — gelesen: sie betrifft den Transaktionspuffer des
Capture-Pfads, nicht die Retention-Lesung.

**Verantwortlich:** Implementer-Agent, 2026-09-25.

**Autor:** Implementer-Agent, Ausgang von `slice-backfill-speicher-untersuchung`;
Nachzug auf das Architect-Verdikt: Planner-Agent. **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Speicher des Feed-Containers in Ruhe und nach einem Backfill hängt
nicht an der Zahl der Changes in `cdc.change`: der periodische Retention-Lauf
bestimmt die Bereinigungsmenge aus schmalen Kandidaten (Kennung, Commit-Position,
Commit-Zeitpunkt, kein Row Image) in Seiten fester Größe. Die Freigabe je Change
bleibt bei der Domain Policy.

**Lösung** ([`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md),
Festlegungen 1 bis 6):

- `ChangeStorePort` bekommt eine Methode
  `ReadRetentionCandidates(ctx, source, after, limit)` und den Transport-Typ
  `RetentionCandidate` (`ChangeID`, `Position`, `CommittedAt`); `ReadChanges`,
  `DeleteChanges` und ihre Aufrufer bleiben.
- `retention.PageSize` = 10.000 Kandidaten ist eine Konstante im Paket des Use
  Cases; der Schlüssel der Seite ist `change_id` (Index des Primärschlüssels, keine
  Sortierung), nicht die fachliche Ordnung der Lesewege.
- `Run` liest Consumer-Positionen und Uhr einmal und arbeitet dann je Seite:
  Kandidaten lesen, `AllowsDeletion` je Kandidat, ein `DeleteChanges` über die
  Freigaben **dieser** Seite (eigene Datenbank-Transaktion). Ein Lauf ist über die
  Seiten nicht atomar; ein Abbruch hinterlässt ein Präfix der Löschmenge, der nächste
  Takt setzt fort.
- Ein SQL-Prädikat, das `AllowsDeletion` nachbildet, ist ausgeschlossen
  ([`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md): „Retention ist eine Domain
  Policy; sie liegt nicht in einem Adapter“).

**Befund (Herkunft: Messbericht, gemessen; Code am Stand `5b1f7762` gelesen).**
`RunRetentionService.Run` (`internal/application/usecase/retention/service.go`)
ruft je Takt (`retentionInterval`, 10 s, `internal/bootstrap/wiring.go`)
`ChangeStorePort.ReadChanges` mit der Quelle als einzigem Filter auf; die Abfrage
`SelectChanges` (`internal/adapters/driven/postgresstorage/queries/queries.go`) trägt
`LIMIT NULL`, liefert je Change beide Row Images in fachlicher Ordnung und wird
vollständig in eine Liste gelesen, bevor `RetentionPolicy.AllowsDeletion` je Change
befragt wird. Die Freigabe liest je Change nur Commit-Zeitpunkt, Commit-Position und
Kennung (Schleife in `service.go`; Architect-Verdikt, Verdikt 1). Der Speicher des
Feed-Containers wächst deshalb mit der Zahl der Changes, unabhängig davon, ob ein
Backfill oder die laufende Erfassung sie geschrieben hat (Code gelesen; die
Live-Erfassung als Quelle ist nicht gemessen): etwa 1,0 bis 1,6 KiB je Change bei
schmalen Zeilen (abgeleitet, Messbericht Abschnitt 1, Ergebnis 3). Die Sicherheit der
Ursache ist hoch (Messbericht Abschnitt 4: Schalter „Bereinigung aus“ an fünf Reihen,
Schalter „`cdc.change` leeren“). Ab 2.000.000 Changes vor einem Run bleiben die
Bereinigungs-Takte in den Reihen B und C aus (beobachtet, Messbericht Abschnitt 3.3);
die Ursache ist nicht erklärt, die Sicherheit dort niedrig (Abschnitt 4).

**Reichweite.** Die Lesung besteht in `v0.1.0`, `v0.1.1` und `v0.1.2` unverändert
(gemessen am 2026-09-25): `git show <Tag>:internal/application/usecase/retention/service.go`
trägt in allen drei Tags die Zeile `ReadChanges(ctx, outbound.ChangeQuery{Source:
command.Source})`, `retentionInterval` ist in den drei `wiring.go` 10 s;
`git diff v0.1.2 HEAD -- internal/application/usecase/retention
internal/application/port/outbound/changestore.go` ist leer; die drei Tags enthalten
keine Datei unter `internal/application/usecase/backfill`. Die Grenze besteht damit in
der Live-Erfassung ohne Backfill; der Bedarf dort ist im Verdikt hergeleitet (Verdikt 6:
10 Changes/s halten 864.000 Changes, etwa 0,85 bis 1,3 GiB Spitze je Takt), nicht an
`v0.1.2` gemessen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Bytegrenze oder Zeilenbreiten-Wache im Backfill** (`NextBlock`,
  `DefaultBlockSize`). Der Messbericht belegt, dass der Speicher im Run selbst flach
  bleibt (Abschnitt 3.2 und 3.4); die Blockgrenze in Zeilen ist nicht die Ursache.
- **Kandidaten auf Transaktions-Ebene** ([`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  Alternative E): sie ändern den Lösch-Vertrag und das Verhalten für Transaktionen
  ohne Change; das Ziel dieses Slice braucht das nicht. Der Re-Evaluierungs-Trigger (a)
  der ADR benennt den Fall.
- **Konfigurierbarkeit von Takt, Mindestalter oder Seitengröße.** Die Seitengröße ist
  eine Konstante ([`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  Festlegung 3); Trigger (b) der ADR.
- **Andere Aufrufer von `ReadChanges` ohne Limit.** Der Architect hat sie gelesen
  (Verdikt 6): kein weiterer produktiver Aufrufer mit Wirkung im Feed-Container;
  `GET /changes` ohne `limit` bleibt eine bewusste Festlegung
  ([`SPEC-022`](../../../../spec/pflichtenheft.md)).
- **Ein Nachschärfen der Richtgröße von 4.000.000 geschätzten Zeilen in
  `internal/application/usecase/backfill/warn.go` ohne Anlass.** Ihre Bewertung nach
  der Nachmessung ist ein Liefer-Punkt (DoD); der Wert ändert sich nur, wenn die
  Bewertung es verlangt (Trigger „Eine Messung liegt vor“,
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)).
- **Ein Gate oder eine Schwelle für den Speicher.** Die Messung bleibt ohne
  Pass/Fail (`harness/targets/bench-backfill.md`); eine Schwelle brauchte eine ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6).

## 2. Definition of Done

- [x] Umsetzung nach [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md):
      `ReadRetentionCandidates` und `RetentionCandidate` am `ChangeStorePort`, Kommentar
      am Port nennt den Seitenvertrag (Ordnung des Schlüssels, leere Seite ist das
      Ende); `Run` arbeitet je Seite (`PageSize` = 10.000, Positionen und Uhr einmal
      je Lauf, ein `DeleteChanges` je gelesener Seite, `Deleted` zählt über alle
      Seiten); Adapter mit der Abfrage aus
      [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
      §Vorgabe der Abfrage, `limit` kleiner 1 endet als `outbound.ErrNonPositiveLimit`,
      eine leere Quelle als `ErrEmptyIdentifier`; kein Löschprädikat in SQL. *Zu
      belegen durch:* Lesen des Diffs gegen die sechs Festlegungen der ADR.
- [x] Unit-Tests (`make test`, `retention/service_test.go`, Fake mit Begrenzung je
      Aufruf und Aufruf-Zähler): (1) dieselbe freigegebene Menge bei Begrenzung 1, 2, 3,
      N und größer als N, verglichen mit einer **handgeschriebenen** Erwartung, nicht
      mit dem Ergebnis des Use Case; (2) jeder Lese-Aufruf trägt `limit` =
      `retention.PageSize`, `after` ist beim ersten Aufruf leer und danach die letzte
      Kennung der vorigen nichtleeren Seite; (3) jeder `DeleteChanges`-Aufruf trägt
      höchstens eine Seite und nur Kennungen dieser Seite; (4) der Lese-Fehler beim
      Aufruf `n` hinterlässt genau die Löschungen der Seiten davor und liefert den
      Fehler, der Lösch-Fehler bei Seite `n` ebenso; (5) ein Fake, der kürzere Seiten
      liefert als `limit`, beendet den Lauf nicht vorzeitig, nur die leere Seite; (6) die
      bestehenden Tests zu [`LH-FA-RET-002`](../../../../spec/lastenheft.md) bis
      [`LH-FA-RET-004`](../../../../spec/lastenheft.md) laufen mit angepasstem Fake ohne
      geänderte Erwartung — eine Erwartung, die sich ändern müsste, ist ein Befund im
      Bericht, keine stille Anpassung. Die Seitengröße ist eine Konstante; die
      Begrenzung auf 1 kommt vom Fake. *Zu belegen durch:* `make test`.
- [x] Store-Tier (`make test-store`, `postgresstorage`, reale PostgreSQL): (1) die
      Seiten einer Quelle decken genau ihre Changes ab (Vereinigung, keine Doppelten,
      keine Kennung einer anderen Quelle) bei Begrenzung 1, 2 und größer als die Menge;
      (2) `Position` und `CommittedAt` je Kandidat sind gleich denen des Datensatzes von
      `ReadChanges` für dieselbe Kennung, mit einer Kennung `0bf-…` (Backfill) und einer
      WAL-Kennung derselben Position; (3) ein Commit **hinter** dem Cursor zwischen zwei
      Seiten fehlt in diesem Durchlauf und steht im nächsten, ein Commit **davor** steht
      im laufenden; (4) ein Durchlauf über mehr als zwei echte Seiten (etwa 25.000
      Changes, eine `generate_series`-Anweisung) löscht über den echten Use Case
      dieselbe Menge, die eine unabhängige SQL-Zählung nach denselben Regeln nennt;
      (5) das Lesen gelingt unter einer `cdc_admin`-Login-Identität (die Zeile `GRANT
      SELECT, DELETE ON cdc.transaction, cdc.change TO cdc_admin` in
      `tools/schema/nacharbeit-roles.sql` trägt das Recht); (6) **Mutationen der
      Eingabeseite**, je mit gesehenem Rot: der Cursor `>=` statt `>` (die letzte Kennung
      käme doppelt), der Quellfilter entfernt (eine Fremdquelle erscheint), `LIMIT`
      entfernt (eine Seite liefert alles), `ORDER BY` entfernt (der Cursor überspringt
      Kennungen). *Zu belegen durch:* `make test-store` samt den gedruckten Mutationen im
      Bericht des Slice.
- [x] Gate-Zuordnung (`harness/sensors/db-adapter-coverage.md` §Gegenstand): die
      Übersetzungsfunktion liegt in `postgresstorage/sqlexec` (Gegenstand des
      Unit-Gates, `make coverage-gate`, netzlos über einen `Executor`-Fake geprüft), die
      Methode am Adapter in `postgresstorage` (DB-Adapter-Coverage, `make test-store`),
      die Abfrage-Konstante in `postgresstorage/queries` (Textkonstante ohne
      Statements). Kein Test ohne Verbindung im DB-Paket hebt die DB-Zahl.
- [x] Die Nachmessung liegt vor: `make image`, dann `tools/bench-backfill-memory.sh` mit
      `BENCH_MEM_STAGES=1000000` (Default-Image, drei Runs mit einem frischen
      Feed-Container je Run; Bestand in `cdc.change` vor den Runs 0, 1.000.000 und
      2.000.000, danach 3.000.000); die gedruckten Zeilen stehen im Bericht des Slice
      unter `docs/reviews/`, jede Zahl mit Host, Lauf und Zeile
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). *Erwartet* (Orientierung,
      keine Schwelle): die Spitze (`memory.peak`, gedruckt nach dem 60-s-Nachlauf) hängt
      nicht an der Zahl der Changes — die drei Werte liegen in der Größenordnung der
      Werte bei ausgeschalteter Bereinigung am Messhost des Messberichts (Reihe J:
      13,3, 13,2 und 13,7 MiB bei 1.000.000, 2.000.000 und 3.000.000 Changes) zuzüglich
      des Bedarfs einer Seite (hergeleitet etwa 1,5 MiB,
      [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
      §Was diese ADR nicht behauptet), und der Zähler „Bereinigung gelaufen“ wächst im
      Nachlauf weiter (hergeleitet aus dem Takt von 10 s: etwa fünf bis sechs je 60 s).
      Trifft das nicht zu — die Spitze wächst mit der Zahl der Changes oder die Takte
      bleiben aus —, steht der Befund im Bericht und der Slice geht nicht nach `done/`.
      *Ergebnis* (Messbericht Abschnitt 7.1, gemessen und abgeleitet je Zahl dort):
      Spitze flach im Bereich 14,9 bis 17,6 MiB bei 1.000.000 bis 3.000.000 Changes
      (zwei Läufe); Restanstieg von 2,4 und 2,5 MiB benannt (1,3 Bytes je Change,
      abgeleitet aus der Differenz der Endpunkte, 0,08 bis 0,12 % der 1,03 bis 1,59 KiB
      je Change vor der Änderung), Ursache nicht untersucht; die Takte laufen weiter
      (20 bis 23, im Review-Lauf 22 bis 24, keine fehlgeschlagene). Die Klausel ist der
      Sache nach nicht ausgelöst (Entscheidung des Auftraggebers; Begründung Messbericht
      7.1: Takte laufen, Größenordnung, Sprung 1.000.000 → 2.000.000 fällt mit dem
      Ausgangszustand des Containers zusammen, Streuung derselben Stufe über drei Läufe
      0,2 bis 0,8 MiB); dem Wortlaut nach steigen die Bereiche der zwei Läufe ohne
      Überlappung. Grenze: drei Stufen bis 3.000.000 Changes, schmale Zeilen, `n` = 2
      plus `n` = 1 (Review); darüber ist nichts gemessen.
- [x] Regressions-Beleg der Live-Last: der Skalierungs-Lauf (`tools/bench-scaling.sh`, das
      zweite Skript von `make bench`, Belegträger von
      [`LH-QA-PER-002`](../../../../spec/lastenheft.md)) läuft einmal; `cdc_capture_lag`
      je Lastenstufe liegt unter der bestehenden 60-s-Grenze
      ([`SPEC-013`](../../../../spec/pflichtenheft.md)), gedruckte Zeilen im Bericht.
      Kein neues Gate; der Vergleich gilt der Grenze, nicht einem Vorwert. Das Ergebnis
      der übrigen `make bench`-Skripte steht mit dem Messhost „zur Kenntnis“, kein
      Kriterium (`BEO-PGC/dod-kriterium-haengt-am-messhost`).
- [x] Der Retention-Lebenszyklus-Rundlauf von `make test-integration`
      (`tools/harness/run-integration-tests.sh`, Wartezeit über mehr als zwei Takte)
      bleibt ohne geänderte Erwartung grün; ein Ausfall wird gegen
      `BEO-PGC/test-integration-retention-timing-flake` gelesen, nicht stillschweigend
      wiederholt.
- [x] Träger nachgezogen, nach der Nachmessung und mit Ursprung je Zahl: das
      Benutzerhandbuch (Abschnitte „Aufbewahrung (Retention)“, „Bestand als Backfill
      überführen“, Absatz „Speicher des Feed-Containers“, „Grenzwerte“: Speicher-Bullet
      samt Bemessungsregel und `docker --memory`-Hinweis ersetzt durch die Nachmessung;
      Richtgröße-Bullet nach der Bewertung unten), `spec/pflichtenheft.md`
      ([`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) Punkt 4 und eine Zeile der
      Änderungshistorie), der Kommentar an `retentionInterval` in
      `internal/bootstrap/wiring.go`; Version und Änderungshistorie des Handbuchs tragen
      eine Zeile. *Zu belegen durch:* Lesen der Abschnitte und `make docs-check`.
- [x] Bewertung der Richtgröße nach dem Trigger von
      [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
      („Eine Messung liegt vor“): sie folgt der Kopierdauer (Messbericht Abschnitt 6) und
      bleibt vom Speicher der Bereinigung unberührt — *erwartet*, dass der Wert bleibt
      und das Speicher-Argument im Handbuch entfällt; die Bewertung und ihr Ergebnis
      (bleibt oder Nachschärfung der Konstante) stehen im Bericht und im Handbuch.
- [x] Die Aussage für die Release-Beschreibung von `v0.2.0` steht im Wortlaut in §7
      (Zeile „Release-Aussage“): der Defekt sitzt in `v0.1.0`, `v0.1.1` und `v0.1.2`
      (Belege: §1 Reichweite), `v0.2.0` trägt die Änderung; Vorschlag und Anker:
      Architect-Verdikt, Abschnitt „Aussage für die Release-Beschreibung“. Zahlen der
      Aussage tragen ihren Ursprung, der Bedarf in `v0.1.x` steht als hergeleitet.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); der Stand des Parent ist eingetragen,
      der Stand des Diff ist Aufgabe des Implementers.
- [x] `make gates` grün — Exit-Code ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8). Report
      `review-slice-retention-lauf-speicher-begrenzung` (0 HIGH, 3 MEDIUM): F-1 und F-2
      in der Fixrunde aufgelöst (Wortlaut in Handbuch 1.63, Pflichtenheft und Godoc;
      Auswertung der Klausel im Messbericht Abschnitt 7.1); F-3 und F-13: Architect-Verdikt
      `architect-verdict-retention-neue-consumer-und-seitengrenze` liegt vor, die
      Ausgänge stehen in §6, das Handbuch trägt die Folgen (1.64). Verifikation
      `verifikation-slice-retention-lauf-speicher-begrenzung`: DoD erfüllt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Slice-Closure selbst (§7): die Roadmap führt mit
      [welle-transformationen](../welle-transformationen.md) eine offene Welle, deren
      Closure hängt aber an zehn Slices in `open/` (gemessen am 2026-09-25 mit `ls
      docs/plan/planning/open`), die Ausführung dieses Slice hängt nicht an ihr, und ein
      Ereignis ohne gesichertes Eintreten ist keine Adresse.

**Umfang:** M — Schätzung, nicht gemessen; Grundlage sind die zwölf Zeilen der Tabelle in
§3 (ein Port, ein Use Case, ein Adapter über drei Dateien, vier Fakes, zwei Testebenen,
zwei Träger, eine Messung) ohne Vergleichs-Slice; die Laufzeit der Nachmessung bei
3.000.000 Changes ist ungemessen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/changestore.go` | update | `RetentionCandidate` und `ReadRetentionCandidates(ctx, source, after, limit)` am `ChangeStorePort`; der Kommentar nennt Seitenvertrag, Ordnung des Schlüssels und „leere Seite ist das Ende“ |
| `internal/application/usecase/retention/service.go` | update | Schleife über Seiten mit `PageSize` = 10.000; Paketkopf und `Run`-Kommentar tragen den Ablauf. Über den Plan hinaus: eine Seite, deren letzte Kennung der Cursor ist, endet als Fehler der Klasse `storage` statt als Endlosschleife (gefunden an der Mutation `>=`, bei der der Lauf ohne diese Prüfung nicht endete: Store-Tier-Lauf bis zum Go-Timeout von 10 Minuten) |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | Konstante `SelectRetentionCandidates` (Abfrage aus der ADR, `WHERE t.source_id = $1 AND c.change_id > $2 ORDER BY c.change_id LIMIT $3`) |
| `internal/adapters/driven/postgresstorage/sqlexec/translate.go` | update | Übersetzungsfunktion neben `ReadChanges`: Scan in drei Werte, `mapper.ToPosition(source, commitPosition)` |
| `internal/adapters/driven/postgresstorage/store.go` | update | Methode am Adapter; Validierung (`limit` ≥ 1, Quelle nicht leer) wie bei `ReadChanges` |
| `internal/bootstrap/heartbeat_internal_test.go`, `internal/application/usecase/readchanges/service_test.go`, `internal/application/usecase/capture/service_test.go` | update | drei Fakes tragen die Methode, damit sie den Port erfüllen (nur die Schnittstelle) |
| `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` | update | Unit-Test der Übersetzung (netzlos) |
| `internal/adapters/driven/postgresstorage/retentioncandidates_test.go` | neu (statt `store_test.go`) | Store-Tier-Tests der DoD (Punkte 1 bis 3, 5, Eingabe-Grenzen) samt Eingabeseiten-Mutationen im Testkommentar; eine eigene Datei, weil `store_test.go` und drei weitere Dateien das Schema per `DROP SCHEMA cdc CASCADE` neu aufsetzen und die Rollen-Grants des Rollouts verlieren — die Datei liegt in der Dateinamen-Ordnung davor und nutzt das ausgerollte Schema samt echten Login-Identitäten |
| `internal/bootstrap/retention_pages_test.go` | neu (statt `store_test.go`) | Store-Tier-Test 4 der DoD (etwa 25.000 Changes, echter Use Case, unabhängige SQL-Zählung) unter einer `cdc_admin`-Login-Identität; im Paket der Verdrahtung, weil `.a-check.yml` keine Kante `adapters` → `app` führt (ein Import des Use Cases im Adapter-Paket bräche `make a-check`) |
| `internal/application/usecase/retention/service_test.go` | update | prüfender Fake und Unit-Tests der DoD; in `TestRunStoreReadErrorFollowsSource` trägt der Fake einen Bestand, damit der Lauf der Quelle `src-1` eine nichtleere Seite und damit einen Lösch-Aufruf hat (Erwartung „Lösch-Zug“ unverändert; ohne Bestand endet der Lauf an der leeren ersten Seite ohne `DeleteChanges`, Befund im Bericht) |
| `internal/bootstrap/wiring.go` | update | Kommentar an `retentionInterval` („Jeder Takt liest alle Changes der Quelle …“) nennt die Seiten |
| `docs/user/benutzerhandbuch.md` (Retention, §Backfill Speicher-Absatz, §Grenzwerte, Version, Änderungshistorie) | update | Nachmessung statt Bemessungsregel; Richtgröße nach der Bewertung |
| `spec/pflichtenheft.md` ([`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md), Änderungshistorie) | update | Punkt 4 ohne ADR- und Slice-Bezug (das Doc-Gate verbietet die Kante Spec → ADR/Slice): „Die Bereinigungsmenge wird seitenweise bestimmt: der Arbeitsspeicher eines Bereinigungslaufs hängt an der Seitengröße (10.000 Kandidaten), nicht an der Zahl der gespeicherten Changes; eine Seite trägt je Change nur Kennung, Commit-Position und Commit-Zeitpunkt, keine Row Images.“ (Wortlaut-Vorschlag des Architects) |
| `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md` | neu | gedruckte Zeilen der Nachmessung, des Skalierungs-Laufs und der Mutationen des Store-Tiers (samt Mutationen des Use Cases und der Übersetzung, Befunde, Bewertung der Richtgröße); Form nach `messbericht-slice-backfill-speicher-untersuchung`, die gedruckten Zeilen im selben Dokument statt in einem Zeilen-Dokument |
| `harness/README.md` (Zeile `make test-store`) | update | Erwähnung der Kandidaten-Seiten-Tests (Entscheidung des Implementers laut Suchlauf-Tabelle) |
| `harness/sensors/db-adapter-coverage.md`, `harness/sensors/coverage-gate.md` | update | Nenner der Zählbasis nach dem Lauf des Slice (1064 bzw. 2601, je mit Lauf und gedruckter Zeile); §3.13-Träger der bewegten Statement-Zahlen |
| `harness/targets/bench-backfill.md` | update | der Halbsatz „Speicherbedarf des Feed-Containers bis in den GiB-Bereich“ entfällt (Träger der bewegten Eigenschaft, gefunden beim Lesen der Suchlauf-Treffer); in der Fixrunde zusätzlich die Zeile „Speicher“ (Grenzen des Skripts): „der Speicher des Feeds hängt an der Zahl der Changes“ ist nach der Messung nicht mehr wahr und nennt jetzt den Seiten-Anteil des Bereinigungslaufs und den Seiten-Cache in `memory.peak` (Fund der Fixrunde-Suche, Befehl B unten) |
| `internal/application/usecase/retention/service_test.go` (Fixrunde, F-4) | update | der Fake trägt ein Budget an Lese-Aufrufen (`maxReads`, Vorgabe `defaultMaxReads` = 20, Fehler `errReadBudget`): die Mutationen „Fortschrittsprüfung abgeschaltet“ und „Cursor bleibt leer“ färben Tests jetzt durch eine Assertion rot statt durch die Zeitüberschreitung des Testlaufs (Messbericht Abschnitt 5, Zeile „Eine Seite ohne Fortschritt“) |
| `internal/application/usecase/retention/service.go` (Fixrunde, F-1, F-8) | update | Godoc von `PageSize` nennt, was die Konstante tut (Kandidaten je Lese-Aufruf, eine Seite im Speicher), statt der Speicher-Aussage; Godoc von `Run` ohne den ausgeschlossenen Ausgang „nicht als Endlosschleife“ (nur Kommentare) |
| `docs/user/benutzerhandbuch.md`, `spec/pflichtenheft.md` (Fixrunde, F-1, F-6) | update | Wortlaut „hängt nicht an der Zahl der Changes“ auf „nicht mit nennenswertem Betrag“ mit gemessenem Bereich (14,9 bis 17,6 MiB bei 1.000.000 bis 3.000.000 Changes; darüber nichts gemessen); die 1,3 Bytes je Change als Differenz zweier Endpunkte gekennzeichnet; dritter Lauf (Review) mit Seiten-Cache-Anteil; Vorversionen im Ist-Zustand; Handbuch 1.63 mit Historienzeile; Punkt 4 von `LH-FA-RET-004.a` als Sachaussage ohne Zahl im Speicherbedarf |
| `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md` (Fixrunde, F-2, F-9 bis F-11) | update | Abschnitt 7 trägt die Auswertung der Plan-Klausel (7.1) und die Richtgröße (7.2); Grundlinie trennt den Seiten-Anteil (Abschnitt 3); Kennzahl `memory.peak` mit Cache-Anteil; zwei Zitier-Ungenauigkeiten berichtigt (Wert bei 3.000.000 Changes, Stand des Image) |

**Hinweis zur Bindung des 25.000-Changes-Tests** (Review F-7, INFO): in
`internal/bootstrap/retention_pages_test.go` ist die Consumer-Position gleich der
Change-Position an Position 3.500 jünger als die 24-h-Grenze der festen Uhr; die
Mutation `t.commit_position + 1` (S5) bleibt dort grün. Die Grenze „Consumer-Position
gleich Change-Position“ tragen der Einheitentest (`c4` an Position 200, Consumer 200,
`TestRunDistinguishesEligibleChangesFromMixedSet`) und der Adapter-Test gegen
`ReadChanges` (`TestRetentionCandidateCarriesPositionAndTimeOfReadChanges`, S5 rot);
der 25.000-Changes-Test bindet sie nicht.

Die vier Fakes (Zählung am Stand `5b1f7762`: `git grep -n -E 'func \(.*\) DeleteChanges'
-- '*.go'` nennt fünf Methoden, vier in Tests und eine am Adapter): der Fake in
`retention/service_test.go` prüft, die drei übrigen tragen nur die Schnittstelle.

**§3.13-Suchlauf** (committetes Feld — bewegte Eigenschaft: „die Retention liest alle
Changes der Quelle samt Row Images in den Speicher; der Speicher hängt an der Zahl der
Changes“; die Befehle stehen im Codeblock, der Parent ist `5b1f7762` (gemessen am
2026-09-25 mit `git grep <Befehl> 5b1f7762 -- <Wurzeln>`), den Stand des Diff trägt der
Implementer ein; Suchform nach [`AGENTS.md`](../../../../AGENTS.md) §3.13: Symbolname,
Zählwort, Beschreibung samt Hedge):

```text
git grep -n -E 'DeleteChanges|ReadChanges|SelectChanges|ChangeStorePort' -- internal test tools spec docs/user harness
git grep -n -E '1,03 bis 1,5[79]|2,65 bis 4,19|2 KiB je Change|4,5 KiB|7,7 GiB|8,9 bis 10,5|Bemessung des Speicherlimits|4\.000\.000' -- docs/user harness spec internal tools
git grep -n -i -E 'alle Changes der Quelle|liest alle|Bereinigungslauf|Retention-Lauf|Bereinigungs-Takt|retentionInterval|Speicher des Feed-Containers|Spitze im Run|Speicherlimit|unbegrenzt|nicht begrenzt' -- docs/user harness spec internal tools README.md docs/plan/adr
```

**Stand und Trefferzahl.** Parent `5b1f7762`: Befehl 1: 298 Zeilen (`internal` 263, `test`
24, `tools` 2, `spec` 7, `docs/user` 2, `harness` 0); Befehl 2: 12 (alle in `docs/user`,
sonst 0); Befehl 3: 80 (`docs/user` 18, `harness` 5, `spec` 5, `internal` 16, `tools` 9,
`README.md` 0, `docs/plan/adr` 27). Parent am Start der Arbeit (`8cd39719`, gemessen mit
`git grep -n … 8cd39719 -- <Wurzeln>`): Befehl 1: 298 (`internal` 263, `test` 24, `tools` 2,
`spec` 7, `docs/user` 2, `harness` 0); Befehl 2: 13 (alle `docs/user`); Befehl 3: 87
(`docs/user` 25, `harness` 5, `spec` 5, `internal` 16, `tools` 9, `README.md` 0,
`docs/plan/adr` 27; die Abweichung zu den 80 Zeilen für `5b1f7762` folgt aus den seither
geschriebenen Handbuch-Zeilen). Diff-Stand `b2484add` (Ende der Umsetzung, vor der
Fixrunde; `git grep -n … b2484add -- <Wurzeln>`): Befehl 1: 309 (`internal` 274, `test` 24,
`tools` 2, `spec` 7, `docs/user` 2, `harness` 0); Befehl 2: 6 (alle `docs/user`); Befehl 3: 85
(`docs/user` 18, `harness` 7, `spec` 7, `internal` 17, `tools` 9, `README.md` 0,
`docs/plan/adr` 27). Stand `bade3fef` (`HEAD` vor der Closure; die Fixrunde ergänzt Handbuch-
und Vertragszeilen, `git grep -n … bade3fef -- <Wurzeln>`, gemessen am 2026-09-25): Befehl 1:
309; Befehl 2: 7; Befehl 3: 88 (`docs/user` 20, `harness` 8, `spec` 7, `internal` 17, `tools` 9,
`README.md` 0, `docs/plan/adr` 27). Im Arbeitsbaum der Closure trägt Befehl 3 eine Zeile mehr
(89, `docs/user` 21: die Historienzeile 1.64 des Handbuchs), Befehl 1 und 2 sind unverändert.

**Ergebnis des Suchlaufs am Diff (Implementer).** *Gefunden und gezogen:* die Träger der
Tabelle in §3 (Port, Use Case, Adapter, vier Fakes, Kommentar an `retentionInterval`,
Handbuch an vier Stellen samt Version und Änderungshistorie, `LH-FA-RET-004.a` samt
Änderungshistorie); dazu drei Träger, die der Suchlauf des Planners nicht nannte:
`harness/sensors/db-adapter-coverage.md` (Nenner 1058 → 1064, mit Lauf),
`harness/sensors/coverage-gate.md` (Nenner 2567 → 2601, mit Lauf) und
`harness/targets/bench-backfill.md` (der Halbsatz „Speicherbedarf des Feed-Containers bis in
den GiB-Bereich“ beschrieb die bewegte Eigenschaft; gefunden beim Lesen der Treffer von
Befehl 3 an dieser Datei). *Ein Lokator, der mitgewandert wäre:* `coverage-gate.md` nennt
Zeilen in `internal/bootstrap/wiring.go` (`:1287.4,1288.1`, `:1183.5,1184.13`); der Kommentar
an `retentionInterval` liegt davor. Die Änderung des Kommentars ist deshalb zeilenneutral
(vier Zeilen vor, vier nach der Änderung, `git diff --stat` zeigt 4 Einfügungen, 4 Löschungen);
die Lokatoren stehen unverändert und sind am Diff geprüft (`git show 8cd39719:…` gegen den
Arbeitsbaum, Zeile 1287 `return`). *Nicht gefunden:* die Träger, die der Suchlauf des Planners
als „nicht zu ziehen“ führt, sind unverändert (Accepted-ADRs, Records, `SPEC-022`,
`tools/harness/run-integration-tests.sh`, `tools/schema/nacharbeit-roles.sql`); im Rest von
Befehl 2 stehen nur die Richtgröße (drei Zeilen), der Wert „8,9 bis 10,5“ des Handbuchs (Run
bei leerem `cdc.change`, unverändert wahr) und zwei Zeilen der Änderungshistorie (Records).
Die Grenze der Suchform bleibt (Zahlen und Lokatoren trifft `grep` nicht zuverlässig); die
Zeilen-Lokatoren wurden gezielt mit `git grep -n -E 'wiring\.go:[0-9]+|:1[0-9]{3}\.[0-9]'`
gesucht.

**Suchlauf der Fixrunde** (bewegte Eigenschaft: „der Speicher des Feed-Containers hängt
nicht an der Zahl der Changes“ wird zu „nicht mit nennenswertem Betrag“; Suchform
Symbolname/Wortlaut, Zählwort, Beschreibung; beide Stände gemessen: vor der Fixrunde
`c21c43b9` mit `git grep … c21c43b9 -- <Wurzeln>`, danach der Arbeitsbaum):

```text
git grep -n -i 'nicht an der Zahl\|hängt nicht' -- docs/user spec harness internal docs/plan/adr
git grep -n -i -E '(an|von|mit) der Zahl der (gespeicherten )?Changes' -- docs/user spec harness internal docs/plan/adr
```

Befehl A: `c21c43b9` 7 Zeilen (`docs/user` 3, `spec` 2, `docs/plan/adr` 2, `harness` 0,
`internal` 0), Arbeitsbaum 5 (`docs/user` 2, `spec` 1, `docs/plan/adr` 2). Befehl B:
`c21c43b9` 10 Zeilen (`docs/plan/adr` 4, `docs/user` 4, `harness` 1, `spec` 1), Arbeitsbaum
11 (`docs/plan/adr` 4, `docs/user` 5, `spec` 2). *Gefunden und gezogen:* die Handbuch-Zeilen
in §4 (Retention, Backfill) und §9 (Speicher-Bullet, Richtgröße-Bullet), Punkt 4 von
`LH-FA-RET-004.a` samt Änderungshistorie, der Godoc von `PageSize`, in `harness/targets/bench-backfill.md`
die Zeile „Speicher“ (von Befehl B getroffen, vom Befehl der Review-Fundstelle nicht). *Nicht
gefunden bzw. nicht gezogen:* `docs/plan/adr/0124-…md` Zeile 209 („nicht an der Zahl der
Changes“, `Accepted`, [`AGENTS.md`](../../../../AGENTS.md) §3.5 — eine inhaltliche Korrektur
wäre eine neue ADR mit `Supersedes`, keine Änderung dieser Zeile), `docs/plan/adr/0065-…md` (anderer Gegenstand: Ausschlussstand),
Handbuch Zeile „Dauerhaftigkeit“ (anderer Gegenstand), `spec/pflichtenheft.md` Zeile 672 (anderer
Gegenstand). Die Handbuch-Aussagen über die Seiten (10.000 je Seite, Positionen einmal je
Lauf) sind von der Formulierungsänderung nicht berührt.

Der Ausdruck `1,03 bis 1,5[79]` in Befehl 2 trifft beide Fassungen des Höchstwerts je Change:
das Handbuch trägt am Parent `5b1f7762` `1,03 bis 1,57`, am Stand `989beef3` `1,03 bis 1,59`
(Reihe B, Run 3 des Messberichts; Nachrechnung: 3.096,6 MiB × 1.024 / 2.000.000 = 1,585 KiB,
abgeleitet aus `verifikation-slice-backfill-speicher-untersuchung` §3, F-2). Gemessen mit
`git grep -c` (Befehl 2 mit dem Ausdruck `1,5[79]`): Parent 12 Zeilen, Stand `989beef3` 13
Zeilen; mit dem früheren Ausdruck `1,03 bis 1,57` am Stand `989beef3` 12 Zeilen, keine davon
die Höchstwert-Zeile des Handbuchs (Handbuch-Text dort `1,03 bis 1,59`).

| Träger | Befund | Behandlung |
|---|---|---|
| `internal/application/port/outbound/changestore.go`, `usecase/retention/service.go`, `postgresstorage/{store.go, queries/queries.go, sqlexec/translate.go}` | Code der Lesung und Löschung (Befehl 1); die Zeilen „ohne Limit unbegrenzt“ am Port (Kommentare zu `ChangeQuery`) beschreiben `ReadChanges` | zu ziehen laut Tabelle oben; die `ChangeQuery`-Kommentare bleiben wahr, `ReadChanges` bleibt |
| vier Test-Fakes (`retention/service_test.go`, `heartbeat_internal_test.go`, `readchanges/service_test.go`, `capture/service_test.go`) | tragen `DeleteChanges` und `ReadChanges` | zu ziehen: Methode ergänzen |
| `internal/bootstrap/wiring.go`, Kommentar an `retentionInterval` | „Jeder Takt liest alle Changes der Quelle und befragt `AllowsDeletion`“ (Befehl 3) | zu ziehen |
| `docs/user/benutzerhandbuch.md` (vier Stellen: Absatz „Speicher des Feed-Containers“ unter Backfill, Absatz unter Retention, Speicher-Bullet und Richtgröße-Bullet unter Grenzwerte; Befehle 2 und 3) | beschreiben „alle Changes der Quelle“, die Bemessung „2 KiB je Change plus 64 MiB“ und die 7,7 GiB der Richtgröße | zu ziehen nach der Nachmessung |
| `spec/pflichtenheft.md`, [`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) | trägt drei Punkte und den Satz „Retention ist eine Domain Policy“; keine Aussage zur Lesung | zu ziehen: Punkt 4 |
| `tools/schema/nacharbeit-roles.sql` (Kommentar und Grant von `cdc_admin`) | trägt `SELECT, DELETE` auf `cdc.transaction`/`cdc.change` | gelesen, **nicht zu ziehen**: das Lesen der Kandidaten braucht kein neues Recht |
| `harness/README.md`, Zeile `make test-store` | vom Suchlauf nicht getroffen (kein Symbolname); die Zeile führt die Inhalte des Laufs einzeln auf | Implementer entscheidet, ob die Kandidaten-Seiten-Tests eine Erwähnung tragen |
| [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Kontext, Teilfrage 7), [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md) („Ohne `limit` liest der Aufruf unbegrenzt“, „produktiver Aufrufer … Retention-Lauf“), [`ADR-0057`](../../adr/0057-http-grpc-api.md), [`ADR-0118`](../../adr/0118-backfill-umschreiben-im-snapshot-fenster.md), [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) | beschreiben den Stand ihrer Entscheidung (Befehle 1 und 3) | **nicht geändert**: `Accepted`-ADRs ([`AGENTS.md`](../../../../AGENTS.md) §3.5); `ADR-0111` ist der Anker der neuen ADR; `ADR-0124` nennt `1,03 bis 1,57 KiB` in den Zeilen 68, 191 und 199, der Messbericht den Höchstwert 1,59 — für die Entscheidung folgenlos (Faktor 6,9 bis 10,6 statt sieben bis zehn gegen etwa 0,15 KiB, abgeleitet aus `verifikation-slice-backfill-speicher-untersuchung` V-1); die Closure dieses Slice schreibt bei Bedarf die Berichtigung, eine Berichtigungs-ADR nur, wenn eine Entscheidung an der Zahl hinge |
| `spec/pflichtenheft.md`, [`SPEC-022`](../../../../spec/pflichtenheft.md) („Noch nicht begrenzt“) und `docs/user/benutzerhandbuch.md` („Ohne `limit` liest der Aufruf unbegrenzt“) | gelten `GET /changes`, nicht der Retention | **nicht geändert** |
| `tools/harness/run-integration-tests.sh`, `harness/targets/bench-backfill.md` | Kommentare zum Takt; die Zählung „Bereinigungs-Takte“ beschreibt die Zählung, nicht die Lesung | **nicht geändert** |
| `docs/reviews/**`, `docs/plan/planning/done/**`, Beobachtungs-Belege | tragen den Stand ihrer Zeit (Messbericht, Verdikt, Review-Reports) | **nicht geändert**: Records |
| **Nicht gefunden** | `spec/architecture.md` beschreibt die Retention-Lesung nirgends (gemessen: `git grep -n -i -E 'retention|bereinigung' spec/architecture.md` nennt sechs Zeilen, alle Komponenten- und Fehlerklassen-Zeilen); kein Träger in `README.md` | — |

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/` (WIP-Limit 1)
und die Entscheidung liegt vor — beobachtbar:

```text
ls docs/plan/planning/in-progress
git ls-files docs/plan/planning/done | grep slice-backfill-speicher-untersuchung
grep -m1 'Status:' docs/plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md
```

Der erste Befehl nennt nur `roadmap.md`, der zweite den Pfad von
`slice-backfill-speicher-untersuchung` in `done/`, der dritte `**Status:** Accepted`
(gemessen am 2026-09-25 nach dem Move der Closure von `slice-backfill-speicher-untersuchung`:
`in-progress/` trägt nur `roadmap.md`, der Slice liegt in `done/`, die ADR ist `Accepted` —
die Start-Bedingung ist erfüllt).

**Vorbedingung des Server-Release `v0.2.0`** (nicht Start dieses Slice; Entscheidung des
Nutzers, dokumentiert im Architect-Verdikt, Abschnitt „Anlass“): dieser Slice liegt in
`done/`, **bevor** `v0.2.0` veröffentlicht wird. Beobachtbar: `git tag -l 'v*'` nennt
solange kein Tag über `v0.1.2` (gemessen am 2026-09-25: `v0.1.0`, `v0.1.1`, `v0.1.2`).
Ein mechanischer Wächter für diese Bedingung existiert nicht; der Träger ist dieser
Abschnitt und die Release-Aussage in §7.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Port, Use Case, Adapter
  und die Nachmessung nicht in einem Review tragen — der abtrennbare Teil ist die
  Nachmessung samt Skalierungs-Lauf, Träger-Nachzug und Bewertung der Richtgröße (vier
  DoD-Punkte) als eigener Slice, Code und Tests bleiben.
- `in-progress` → `open` (blockiert): falls die Nachmessung an der Kapazität des
  Messhosts scheitert (3.000.000 Changes brauchen Platte und Zeit) oder eine Festlegung
  von [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  sich an der Umsetzung als unerfüllbar zeigt — dann folgt ein Architect-Zug (neue ADR
  mit `Supersedes`), keine stille Abweichung.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + die Nachmessung mit gedruckten Zeilen im Bericht +
Store-Tier-Tests samt gesehenem Rot der Eingabeseiten-Mutationen + Handbuch und
Pflichtenheft nachgezogen + Release-Aussage in §7 + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Bereinigungsmenge ändert sich** (Sicherheit: kein Change eines aktiven
  Consumers darf gelöscht werden, [`LH-FA-RET-004`](../../../../spec/lastenheft.md)).
  *Erwartet, zu belegen durch:* die Unit-Tests (Menge bei jeder Begrenzung gegen
  handgeschriebene Erwartung), der Store-Tier-Test über mehr als zwei Seiten gegen eine
  unabhängige SQL-Zählung und die unverändert grünen Tests der Anforderungen.
  *Befund des Implementers:* alle drei liegen vor und sind rot gesehen
  (Messbericht Abschnitt 5 und 6); eine bestehende Testerwartung blieb, ihr Fake trägt
  einen Bestand (Messbericht Abschnitt 6, Befund 1). **Ausgang:** *entfallen* — Beleg:
  Verifikation `verifikation-slice-retention-lauf-speicher-begrenzung` Abschnitt 5 (sieben
  Use-Case- und vier Store-Tier-Mutationen der Eingabeseite selbst rot gesehen, `make test` und
  `make test-store` Exit 0, `make gates` Exit 0).
- **Die Datenbank-Arbeit je Takt wächst mit der Zahl der Changes** (die Sortierung über
  die Menge entfällt; die Seiten-Abfrage bleibt linear: 0,38 bis 2,0 s je 1.000.000
  Changes, gemessen im Architect-Zug, `n` = 1, [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  §Konsequenzen). Das Ausbleiben der Bereinigungs-Takte ab 2.000.000 Changes ist
  beobachtet und nicht erklärt (Messbericht Abschnitt 3.3); ob die Ursache im
  Feed-Prozess liegt (Dekodierung und Garbage Collection großer Listen), ist
  **Hypothese**, nicht gemessen (Architect-Verdikt, Abschnitt „Befundlage“). Eine
  Gegenprobe nennt der Messbericht nicht: Reihe N (ein Container über alle Stufen,
  2.330.000 Changes vor dem Run) zeigt 60 s nach dem Run 2.448,4 MiB (`anon` nicht auf
  etwa 60 MiB gefallen), die Reihen B und C (frischer Container je Run) zeigen das
  Ausbleiben; der Unterschied ist nicht erklärt (übernommen aus dem Review-Report
  `review-slice-backfill-speicher-untersuchung`, F-11). *Erwartet, zu belegen durch:* die
  Nachmessung bei 1.000.000, 2.000.000 und 3.000.000 Changes — laufen die Takte dort bis
  zum Ende weiter, entfällt die Hypothese als Ursache des Ausbleibens; bleiben sie aus,
  steht der Befund im Bericht und der Slice geht nicht nach `done/`. *Befund des
  Implementers:* die Takte laufen in allen sechs Runs bis zum Ende weiter (20 bis 23 Zeilen
  „Bereinigung gelaufen“ seit dem Start, keine fehlgeschlagene; Garbage-Collection-Läufe im
  Nachlauf 560 bis 2.467); die Ausbleib-Beobachtung tritt nicht mehr auf (Messbericht
  Abschnitt 1, Punkt 3); der Verifier-Lauf bestätigt es (21 bis 24, keine fehlgeschlagene,
  Verifikation Abschnitt 4). Die Dauer eines Laufs ist nicht gemessen (Abschnitt 6, Befund 4).
  **Ausgang:** *weiter offen* — die Hypothese „Feed-Prozess als Ursache des Ausbleibens“
  entfällt als Ursache (die Takte laufen in allen Runs); die lineare Arbeit der Datenbank je
  Takt und die ungemessene Dauer eines Laufs am Feed bleiben. Adresse: Re-Evaluierungs-Trigger
  (a) von [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  (beobachtbar: die Zeilen „retention: Bereinigung gelaufen“ bleiben zwischen den Takten aus,
  oder die Dauer erreicht `retentionInterval`, 10 s; Kandidaten auf Transaktions-Ebene, Alternative E).
- **Ein Backfill oder die Erfassung committet während einer laufenden Bereinigung.** Die
  Seitengrenze darf eine Zeile in **diesem** Lauf überspringen (Commit hinter dem
  Cursor), sie löscht nie früher: die Freigabe hängt an dem, was `AllowsDeletion` für die
  gelesene Zeile entschied ([`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  Festlegung 6). Was nicht sein darf, ist ein **dauerhaftes** Überspringen. *Erwartet, zu
  belegen durch:* der Store-Tier-Test 3 („hinter dem Cursor: fehlt im laufenden
  Durchlauf, steht im nächsten“). *Befund des Implementers:* der Test liegt vor
  (`TestRetentionCandidateBehindCursorAppearsInNextRun`) und ist an `>=` und an einem
  ignorierten Cursor rot gesehen (Messbericht Abschnitt 5). **Ausgang:** *entfallen* — Beleg:
  Verifikation Abschnitt 5 (Mutationen S1 `>=` und S4 ohne `ORDER BY` färben
  `TestRetentionCandidateBehindCursorAppearsInNextRun` rot, selbst gefahren).
- **Ein Consumer bestätigt erstmals während eines längeren Laufs.** Die Positionen
  werden einmal je Lauf gelesen; ein Consumer, der erst danach seine erste Bestätigung
  setzt, ist im laufenden Lauf nicht berücksichtigt (er zählt erst ab seiner ersten
  Bestätigung, Handbuch, „Betriebs-Hinweis (Consumer-Bindung)“). Das Fenster besteht
  bereits zwischen Lesen der Positionen und Löschen und wird mit der Laufdauer breiter;
  [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  Festlegung 4 benennt nur die Vorwärtswanderung der Positionen eines **bekannten**
  Consumers. Die zweite Bedingung der Freigabe, das Mindestalter von 24 h
  (`retentionMinAge`), begrenzt die Kandidaten auf alte Changes. *Erwartet:* ein
  akzeptiertes Negativ ohne neuen Fall; **vom Planner beim Nachzug gefunden, nicht vom
  Architect bewertet** — der Implementer legt die Frage dem Architect vor, bevor er sie
  als akzeptiert schließt. *Befund des Implementers* (am Code, Messbericht Abschnitt 6,
  Befund 3): das Fenster besteht mit dem Lauf vor der Änderung ebenso (Positionen, dann
  Lesen aller Changes, dann Löschen der ganzen Menge); seine Breite ist die Laufdauer: im
  Regelbetrieb 0,38 bis 2,0 s, im Extremlauf (vollständige Löschung, der erste Lauf nach einem
  großen Bestand, der das Mindestalter erreicht) etwa 26 s statt etwa 3,4 s, das Siebenfache
  (abgeleitet aus den Einzelläufen von
  [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  §Gemessen, `n` = 1, nicht am Feed gemessen; Architect-Verdikt
  [`architect-verdict-retention-neue-consumer-und-seitengrenze`](../../../reviews/architect-verdict-retention-neue-consumer-und-seitengrenze.md));
  die Änderung ist keine Vertragsänderung, [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  Festlegung 4 ist wörtlich umgesetzt (Positionen einmal je Lauf,
  `TestRunReadsPositionsOncePerRun`). Eine Erneuerung der Positionen je Seite wäre eine
  Schärfung von Festlegung 4 (Folge-ADR mit `Supersedes`) und ist nicht umgesetzt; das Handbuch
  nennt die Lage im „Betriebs-Hinweis (Consumer-Bindung)“ samt der Folge für einen erstmals
  bestätigenden Consumer (Handbuch 1.64). **Ausgang:** *akzeptiertes Negativ* — Architect-Verdikt
  [`architect-verdict-retention-neue-consumer-und-seitengrenze`](../../../reviews/architect-verdict-retention-neue-consumer-und-seitengrenze.md)
  (F-3: ein Consumer ohne Positions-Zeile ist in jedem Lauf ungeschützt, das Fenster ist der
  letzte Ausschnitt dieser Lücke; keine Vertragsänderung, keine neue ADR).
- **Ein Abbruch zwischen zwei Seiten** hinterlässt ein Präfix der Löschmenge (Lauf nicht
  atomar). *Erwartet, zu belegen durch:* Unit-Test 4 (Fehler beim Aufruf `n` hinterlässt
  genau die Löschungen der Seiten davor). *Befund des Implementers:* Lese- und Lösch-Fehler
  ab Aufruf `n` sind geprüft und rot gesehen (`TestRunReadFailureLeavesPagesBefore`,
  `TestRunDeleteFailureStopsAtThatPage`, Messbericht Abschnitt 5). **Ausgang:** *akzeptiertes
  Negativ* — Architect-Verdikt
  [`architect-verdict-retention-neue-consumer-und-seitengrenze`](../../../reviews/architect-verdict-retention-neue-consumer-und-seitengrenze.md)
  (F-13 a: Alter und Commit-Position sind Eigenschaften der Transaktion, die Freigabe fällt für
  alle ihre Changes gleich aus; ein Consumer mit bestätigter Position sieht eine Transaktion nie
  teilweise gelöscht, ein Consumer ohne Bestätigung ist bereits ungeschützt); das Handbuch nennt
  die Nicht-Atomarität (1.64).
- **Die Größe eines Kandidaten im Go-Prozess ist hergeleitet, nicht gemessen** (etwa
  0,15 KiB, Seite etwa 1,5 MiB, [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  §Was diese ADR nicht behauptet); die Zahl ist eine Orientierung, keine Grenze
  (`BEO-PGC/geschaetzter-wert-als-grenze`). *Erwartet, zu belegen durch:* die Nachmessung
  (der Wert der Spitze steht neben der Herleitung, nicht statt ihrer). *Befund des
  Implementers:* die Spitze liegt bei 14,9 bis 17,6 MiB (Messbericht Abschnitt 3); der Bedarf
  einer Seite ist aus der Nachmessung nicht getrennt (keine Messung mit anderer Seitengröße);
  ein Anstieg der Spitze um 2,4 und 2,5 MiB zwischen 1.000.000 und 3.000.000 Changes ist
  gemessen und nicht erklärt (etwa 1,3 Bytes je Change, abgeleitet als Differenz der zwei
  Endpunkte, keine gemessene Steigung; Auswertung: Messbericht Abschnitt 7.1); der Verifier-Lauf
  liegt bei 15,5 bis 17,3 MiB (Verifikation Abschnitt 4). **Ausgang:** *weiter offen* — der Wert
  der Spitze steht neben der Herleitung, der Bedarf einer Seite ist nicht getrennt gemessen und
  der Restanstieg unerklärt. Adresse: Register `BEO-PGC/geschaetzter-wert-als-grenze` (die 1,5 MiB
  je Seite sind Orientierung, keine Grenze); Restanstieg und Seitenbedarf tragen als benannte
  Grenze Messbericht Abschnitt 7.1 und das Handbuch unter „Grenzwerte“.
- **Die Nachmessung hängt am Host** (Speicher, Platte und Laufzeit bei 3.000.000
  Changes). *Erwartet, zu belegen durch:* jede Zahl nennt Host und Lauf; die Erwartung ist
  als Form gefasst („hängt nicht an der Zahl der Changes“), die Vergleichswerte der Reihe J
  sind Orientierung am selben Host; kein `docker volume prune` und kein `docker system
  prune` im Lauf; der Bericht nennt die Aufräum-Schritte
  (`BEO-PGC/dod-kriterium-haengt-am-messhost`, 1×, offen). *Befund des Implementers:* jede
  Zahl nennt Host, Lauf und Zeile; die Läufe endeten ohne Platten- oder Zeitproblem, die
  Aufräum-Schritte stehen im Messbericht Abschnitt 6 („Aufgeräumt“). **Ausgang:** *weiter offen*
  → Register `BEO-PGC/dod-kriterium-haengt-am-messhost` (1×): die DoD ist hostunabhängig gefasst
  und erfüllt; alle Messwerte (Implementer, Review, Verifier) stammen von einem Host, eine
  zweite Umgebung ist nicht gemessen.
- **Der Slice läuft nicht vor dem Server-Release `v0.2.0`** (kein Wächter).
  *Erwartet, zu belegen durch:* die Vorbedingung in §4 und die Release-Aussage in §7; die
  Prüfung liegt beim Planner der Release-Vorbereitung. **Ausgang:** *entfallen* — Beleg: die
  Closure liegt mit ihrem Move-Commit vor dem Tag; `git tag -l 'v*'` nennt am 2026-09-25 vor dem
  Move `v0.1.0`, `v0.1.1` und `v0.1.2`; die Release-Aussage steht in §7.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Kette Untersuchung → Entscheidung → Umsetzung lief in getrennten
  Kontexten: die Untersuchung fand die Ursache, der Architect legte die Entscheidung
  ([`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)) mit
  Einzelmessungen fest, der Implementer setzte um und maß nach. Die Löschsicherheit trug an der
  Eingabeseite: sieben Use-Case-Mutationen und vier Store-Tier-Mutationen färben Tests per
  Assertion rot, der Verifier fuhr sie selbst (Verifikation Abschnitt 5); der 25.000-Changes-Test
  zählt gegen eine unabhängige SQL-Zählung. Der Review (0 HIGH · 3 MEDIUM · 3 LOW · 7 INFO) fand
  durch Nachmessen die unbedingte Speicher-Aussage gegen die gedruckten Zeilen (F-1) und die
  fehlende Auswertung der Plan-Klausel (F-2); der Verifier bestätigte die DoD mit eigenem
  `make gates`, `make test-store`, Skalierungs-Lauf und einer eigenen Nachmessung, deren Werte
  (15,5 bis 17,3 MiB) im Band der Implementer-Läufe liegen, und fand den offenen Nachzug des
  Architect-Verdikts (V-1).
- **Was ging anders als geplant:** (1) Die Fortschrittsprüfung des Use Case (eine Seite, deren
  letzte Kennung der Cursor ist, endet als Fehler der Klasse `storage`) stand nicht im Plan; die
  Mutation `>=` ließ den Lauf ohne sie nicht enden (§3, Zeile `service.go`). (2) Die zwei
  Store-Tier-Tests liegen in eigenen Dateien statt in `store_test.go` (§3, zwei Zeilen). (3) Die
  Plan-Klausel der Nachmessung („wächst die Spitze mit der Zahl der Changes“) ist dem Wortlaut
  nach an den zwei Läufen erfüllt (Bereiche ohne Überlappung) und der Sache nach nicht
  ausgelöst; die Lesart ist Entscheidung des Auftraggebers, ihre Begründung steht im Messbericht
  Abschnitt 7.1 (F-2). (4) Zwei Bewertungen des Reviews (F-3, F-13 a) gingen an den Architect;
  ihre Nachzüge trägt diese Closure (V-1). (5) Die Fitness-Function-Zeile 3 von
  [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) ist wörtlich
  nicht erfüllt (V-3), siehe unter „Lesart der Fitness-Function-Zeile“.
- **Bekannte `git bisect`-Falle:** der Commit `d49e248f` (Fortschrittsprüfung im Use Case)
  ist rot — `TestRunRejectsPageWithoutProgress` schlägt an ihm fehl (`gelöschte Menge = [c1 c1]`),
  der Folgecommit `7572051c` behebt es; `HEAD` ist grün, die Historie bleibt unverändert.
  Ein `git bisect` über diesen Bereich hält dort an; `git bisect skip d49e248f` übergeht ihn.
- **Steering-Loop-Eintrag (Lerneintrag):**
  *(a) Benannte Grenze und Lerneintrag zum Schnitt — die Ursache lag im Nachbarsystem.* Das
  Speicher-Verhalten der Retention zeigt sich in der Messung des Backfills; es besteht
  in `v0.1.0` bis `v0.1.2` ohne Backfill (§1 Reichweite) und ist ein vorbestehender Defekt.
  Der Untersuchungs-Slice trug im selben Zug Reichweite an den Tags, Vorbedingung im Trigger
  (§4) und Release-Aussage (unten): das ist die Anwendung des Kandidaten
  `BEO-PGC/vorbestehender-defekt-durch-messung-im-nachbarsystem-gefunden`, **kein zweites
  Auftreten** (derselbe Vorgang, kein Finding fand die Lücke); der Eintrag bleibt **1×, offen**,
  ohne Träger. Begründung gegen eine Regel bei 1×: sie wäre eine Formpflicht im Slice-Schnitt eines
  Untersuchungs-Slice, deren Nutzen ein zweiter, unabhängiger Vorgang zeigen muss; der Auslöser
  steht in der `state.md`. *(b) Neuer Sensor auf Testebene — Fake mit Budget.* Die
  Fortschrittsprüfung `last == after` war an ihrer Eingabeseite nur durch die Zeitüberschreitung
  des Testlaufs rot zu färben (Review F-4, LOW); der Fake trägt ein Budget an Lese-Aufrufen
  (`maxReads`, Fehler `errReadBudget`), die Mutationen „Prüfung abgeschaltet“ und „Cursor bleibt
  leer“ färben Tests per Assertion in 0,00 s (Verifikation Abschnitt 5, M1). Gelernt: ein Test einer
  Schleife mit Fortschrittsprüfung bindet die Zusage erst, wenn die Endlosigkeit im Fake eine
  Assertion ist, keine Uhr; die Klasse `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
  trägt sie (12×, verkörpert; **Deckel** — F-4 ist vor dem Merge vom Reviewer gefunden, LOW, bekannter
  Träger-Typ Fake im Use-Case-Test, daher keine weitere Datei, Finding-Kennung F-4 hier); kein
  neuer Träger, die Regel „mutiere den Eingabewert“ deckt den Fall. *(c) Benannte Spec-Lücke,
  geschlossen:* die Lesung der Retention hatte im Pflichtenheft keine Aussage;
  [`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) trägt Punkt 4 (Seiten, Speicher an der
  Seitengröße). *(d) Der Nachzug eines Architect-Verdikts braucht einen Träger je Aussage:* Plan
  §6, Plan-Prosa und Handbuch trugen Aussagen, die das Verdikt überholte (V-1);
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (11×, verkörpert, Deckel) trägt die Klasse — der
  Suchlauf in §3 ist die Regel, kein neuer Träger.
- **Lesart der Fitness-Function-Zeile 3 von [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
  (Verifikation V-3):** die Zeile erwartet die Spitze nach dem Run „im Bereich der Spitze im Run
  (8,9 bis 10,5 MiB) plus dem Bedarf einer Seite“; gemessen sind 14,9 bis 17,6 MiB (`memory.peak`,
  zwei Läufe) und 15,5 bis 17,3 MiB (Verifier-Lauf), der `anon`-Wert der Spitze im Run liegt im
  Verifier-Lauf bei 10,0 bis 13,2 MiB. Wörtlich ist die Erwartung nicht erfüllt, der Form nach ja:
  die Spitze ist flach über der Zahl der Changes und die Takte laufen weiter. Die Zeile ist eine
  Messung ohne Pass/Fail („Werkzeug, kein Gate“, [`AGENTS.md`](../../../../AGENTS.md) §3.6), ihre
  Erwartung ist hergeleitet. **Kein Berichtigungs-ADR:** keine der Festlegungen 1 bis 6 der ADR
  hängt an der Zahl (gemessen: `grep -n -E 'MiB|KiB'` über die Festlegungen im Abschnitt
  „Entscheidung“ der ADR, Zeilen 110 bis 186, trifft nichts; eine ADR ist nach `Accepted`
  unberührbar, [`AGENTS.md`](../../../../AGENTS.md) §3.5); die Lesart steht hier, im Messbericht
  Abschnitt 7.1 und im Register `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (6×). Dieselbe
  Entscheidung gilt für die Abweichung 1,57 gegen 1,59 KiB in der ADR (Zeilen 68, 191, 199):
  die Faktoren der Kandidaten-Größe werden 6,9 bis 10,6 statt sieben bis zehn (abgeleitet, §3
  Träger-Tabelle, Zeile zu den `Accepted`-ADRs), keine Festlegung hängt an ihr; die Release-Aussage
  nennt 1,03 bis 1,59.
- **Benannte Grenzen und Restrisiken (Verifikation):** *V-4* (INFO): `memory.peak` trägt einen
  Seiten-Cache-Anteil, der zwischen den Läufen streut (Review-Lauf: 31,8 MiB mit 16,5 MiB `file`
  bei 1.000.000 Changes; drei Runs des Verifiers und zwei Läufe des Implementers 0,0 MiB am
  Run-Ende); Handbuch, Messbericht und Vertrag `harness/targets/bench-backfill.md` nennen es,
  die Ursache ist nicht untersucht. *V-5* (INFO): die gedeckte Zahl des Unit-Gates streut um
  drei Statements zwischen den Läufen (2167 gegen 2170 bei Nenner 2601); die Zeile in
  `harness/sensors/coverage-gate.md` nennt ihren Lauf. *Restrisiken, weiter offen:* die Dauer
  eines Bereinigungslaufs am Feed ist ungemessen (Adresse: Trigger (a) von [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md), §6 zweites
  Risiko); die Live-Erfassung als Speicherquelle, breite Zeilen und Läufe über 3.000.000 Changes
  sind ungemessen; alle Messwerte stammen von einem Host, `n` = 2 bis 3 je Stufe (Register
  `BEO-PGC/dod-kriterium-haengt-am-messhost`); die Arbeit der Datenbank je Takt wächst weiter
  linear mit der Zahl der Changes (Trigger: Laufdauer ab 10 s, dann Kandidaten auf
  Transaktions-Ebene); der Restanstieg von 2,4 und 2,5 MiB zwischen 1.000.000 und 3.000.000
  Changes ist unerklärt (die Klausel ist der Sache nach nicht ausgelöst, Messbericht 7.1).
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-retention-lauf-speicher-begrenzung.md`, Zähler = Zahl der Dateien.
  *Neue Belege:* `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` **23×** (F-1 MEDIUM, F-2,
  V-2; verkörpert, Deckel bei 23×), `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` **11×**
  (F-3, V-1, beide MEDIUM; verkörpert, Deckel bei 10×), `BEO-PGC/adr-aussage-breiter-als-ihre-messung`
  **6×** (V-3; verkörpert). *Zustand nachgezogen, ohne neue Datei:*
  `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes` **2×** (Ausgang *verkörpert* → Use Case der
  Retention; der Gegenstand „Block zählt Zeilen“ bleibt in der Port-Doku, Zeilen im MB-Bereich
  bleiben ungemessen, Adresse: Trigger „Ausbaustufe für Durchsatz“ von
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)),
  `BEO-PGC/vorbestehender-defekt-durch-messung-im-nachbarsystem-gefunden` **1×** (Vorgang gelöst,
  Begründung oben), `BEO-PGC/dod-kriterium-haengt-am-messhost` **1×** und
  `BEO-PGC/geschaetzter-wert-als-grenze` **1×** (beide angewandt ohne Anfall).
  *Deckel-Fälle ohne Datei, Finding-Kennung hier:* F-4 (`negativtest-ohne-bindung-an-seine-eingabe`,
  Klasse siehe Lerneintrag b). *Kein eigener Register-Anfall* (Text-, Werkzeug- und
  Historien-Findings ohne wiederkehrende Klasse im Register): F-5 (roter Zwischencommit, benannt
  in „Bekannte `git bisect`-Falle“), F-6 (Vorher-Nachher-Sprache in Handbuch-Prosa, LOW, vor dem
  Merge behoben; Träger sind [`AGENTS.md`](../../../../AGENTS.md) §3.7 und der Reviewer-Skill),
  F-7 bis F-13, V-4, V-5. *Lese-Schritt der nächsten Welle-Closure (`welle-transformationen`):*
  kein Eintrag erreicht mit diesem Beleg neu 3× ohne Ausgang — die drei Einträge mit neuer Datei
  tragen bereits den Ausgang *verkörpert*, die übrigen stehen bei 1× bzw. 2× —, es entsteht kein
  Vermerk in einer `state.md`.
- **Folge-Slices:** keine angelegt. Was offen bleibt, trägt eine Adresse (Risiken unten:
  Trigger (a) von [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md), Register). Übergaben an offene Pläne
  ([`AGENTS.md`](../../../../AGENTS.md) §3.13, Suchlauf am Stand `bade3fef`):
  `slice-sdk-kotlin-cloudsmith` — die zwei Versionsbezüge auf das Handbuch (`Version:` 1.61 am
  Stand `989beef3`) tragen zusätzlich die Version 1.64 nach dieser Closure; dort stehende
  Zeilen-Lokatoren des Handbuchs sind mit ihrem Stand (`37e825e2`) gekennzeichnet und gelten
  nicht als nachzuziehen. Die Slices von `welle-transformationen` nennen keine Handbuch-Version
  und keinen Zeilen-Lokator des Handbuchs (gemessen: `git grep -n -E 'benutzerhandbuch.*(Z\.|Zeile|:[0-9]{3,4})'
  -- docs/plan/planning/open` trifft nur `slice-sdk-kotlin-cloudsmith`, dort mit Stand); ihr
  Bezug auf „retention“ gilt dem Register-Eintrag `test-integration-retention-timing-flake`
  (`git grep -n -i 'retention' -- docs/plan/planning/open`), nicht der Lesung; kein Nachzug. Die
  Handbuch-Zahlen „vor der Behebung“ aus der Closure von `slice-backfill-speicher-untersuchung`
  sind ersetzt: `grep -n -i 'vor der Behebung' docs/user/benutzerhandbuch.md` trifft nichts.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Die Bereinigungsmenge ändert
  sich · Ein Backfill oder die Erfassung committet während einer laufenden Bereinigung · Der Slice
  läuft nicht vor dem Server-Release `v0.2.0`. *Akzeptiertes Negativ* (zählt als entfallen mit
  Verdikt-Anker,
  [`architect-verdict-retention-neue-consumer-und-seitengrenze`](../../../reviews/architect-verdict-retention-neue-consumer-und-seitengrenze.md)):
  Ein Consumer bestätigt erstmals während eines längeren Laufs · Ein Abbruch zwischen zwei Seiten.
  *Weiter offen:* Die Datenbank-Arbeit je Takt wächst mit der Zahl der Changes (Trigger (a) von
  [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)) · Die Größe eines Kandidaten im Go-Prozess ist hergeleitet (Register
  `geschaetzter-wert-als-grenze`) · Die Nachmessung hängt am Host (Register
  `dod-kriterium-haengt-am-messhost`).
- **Release-Aussage (für die Beschreibung von `v0.2.0`):** kopierbare Kurzfassung (Anwender-Sprache;
  die Zahlen tragen ihren Ursprung, gemessen oder als hergeleitet gekennzeichnet; Anker:
  Architect-Verdikt, Abschnitt „Aussage für die Release-Beschreibung“, Messbericht dieses Slice
  Abschnitt 3, Verifikation Abschnitt 4 und 7):

  > **Behoben: Speicherbedarf des periodischen Bereinigungslaufs (Retention) des Feed-Containers.**
  > **Was war** (`v0.1.0`, `v0.1.1`, `v0.1.2`): Der Bereinigungslauf liest alle 10 Sekunden alle
  > gespeicherten Änderungen einer Quelle samt Zeilenabbildern in den Speicher; der Bedarf des
  > Feed-Containers wächst mit der Zahl der gespeicherten Änderungen (1,03 bis 1,59 KiB je Änderung
  > bei schmalen Zeilen, an einem Backfill gemessen; für die laufende Erfassung ohne Backfill nur
  > hergeleitet, nicht gemessen: bei 10 Änderungen je Sekunde und 24 Stunden Mindestalter etwa 0,85
  > bis 1,3 GiB Spitze je Takt).
  > **Was ist** (`v0.2.0`): Der Lauf liest die Kandidaten seitenweise (10.000 je Seite, ohne
  > Zeilenabbilder); die gemessene Spitze liegt bei 14,9 bis 17,6 MiB bei 1.000.000 bis 3.000.000
  > gespeicherten Änderungen (schmale Zeilen, ein Host, zwei Läufe; ein dritter, unabhängiger Lauf
  > 15,5 bis 17,3 MiB).
  > **Was bleibt:** ungemessen sind breite Zeilen, mehr als 3.000.000 Änderungen, die laufende
  > Erfassung als Speicherquelle und die Dauer eines Bereinigungslaufs; die Arbeit der Datenbank je
  > Takt wächst weiter mit der Zahl der Änderungen. Ein Lauf ist nicht atomar; ein Consumer ohne
  > bestätigte Position ist vor der Bereinigung nicht geschützt (Benutzerhandbuch, Abschnitt
  > Aufbewahrung).

  Belege der Reichweite: `git diff v0.1.2 8cd39719 --stat -- internal/application/usecase/retention
  internal/application/port/outbound/changestore.go` ist leer (Verifikation Abschnitt 7); `git show
  v0.1.0:…/service.go`, `v0.1.1:…` und `v0.1.2:…` tragen dieselbe Zeile
  `ReadChanges(ctx, outbound.ChangeQuery{Source: command.Source})`; die Rechnung 10 × 86.400 s =
  864.000 Änderungen, × 1,03 bis 1,59 KiB = 0,85 bis 1,31 GiB ist abgeleitet.
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt Anker,
  Folge-Slice und Register (Baseline-Regelwerk `modul-06-roadmap.md` §Was der wellenlose
  Betrieb selbst auslöst). *Anker:* die Zielorte der Lerneinträge existieren — Use Case
  `internal/application/usecase/retention/service.go` und `service_test.go` (Fake mit Budget),
  [`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) Punkt 4 — und tragen ihre Herkunft im
  Diff dieses Slice. *Folge-Slice:* keiner genannt, kein Versprechen offen. *Register:* jede
  genannte Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/`
  (geprüft mit `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence`). Die Prüfung läuft
  zusätzlich bei der Closure der nächsten Welle
  ([welle-transformationen](../welle-transformationen.md), die Roadmap führt sie unter *Offene
  Wellen*).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Use Case der Retention und Store-Adapter sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler gemessen
am 2026-09-25 nach der Closure von `slice-backfill-speicher-untersuchung` mit `ls evidence | wc -l` je Eintrag; Stand der Einträge aus ihrer
`state.md`):

- `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes` (2×, geplant → dieser Slice): der Befund
  der Untersuchung trägt diesen Slice; die `state.md` führt die Ursache als untersucht und
  belegt (Messbericht `messbericht-slice-backfill-speicher-untersuchung` §4).
- `BEO-PGC/vorbestehender-defekt-durch-messung-im-nachbarsystem-gefunden` (1×, offen): die
  Vorbedingung des Server-Release in §4 und die Release-Aussage in §7 sind der Träger des
  Release-Rahmens.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (22×, verkörpert, Deckel): die
  Zahlen des Handbuchs und der Bemessungsregel — jede Zahl der Nachmessung trägt Ursprung
  und Lauf; die Bemessungsregel ist abgeleitet oder entfällt.
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (32×, verkörpert): der Suchlauf in §3.
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (12×, verkörpert, Deckel): die
  Eingabeseiten-Mutationen des Store-Tiers in der DoD.
- `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (3×, verkörpert): die
  Gate-Zuordnung in der DoD.
- `BEO-PGC/dod-kriterium-haengt-am-messhost` (1×, offen): das Risiko der Nachmessung; das
  Ganz-Target `make bench` steht „zur Kenntnis“.
- `BEO-PGC/geschaetzter-wert-als-grenze` (1×, offen): die Seitengröße von 1,5 MiB ist eine
  Orientierung; die Erwartung der Nachmessung ist keine Schwelle.
- `BEO-PGC/test-integration-retention-timing-flake` (3×, verkörpert): der
  Retention-Lebenszyklus-Rundlauf als Regression der Verdrahtung.
- `BEO-PGC/aufschub-adresse-verfaellt` (3×, verkörpert): die Adresse der Paarungen ist die
  Slice-Closure selbst (DoD letzter Punkt).
- `BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (1×, offen) und
  `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (2×, offen): der Start-Trigger nennt
  die Entscheidung als Artefakt
  ([`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md),
  `Accepted`), die Rückführungs-Bedingung „der
  Vertrag verlangt eine ADR“ ist vor dem Start ausgewertet.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
