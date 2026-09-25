# ADR-0121: Capture — Bindung der Bedingung „keine offene Transaktion“ der Leerlauf-Bestätigung (Supersedes ADR-0120, teilweise)

**Status:** Accepted — Supersedes [`ADR-0120`](0120-capture-slot-leerlauf-bestaetigung.md)
in genau zwei Stellen: die Store-Zeile der Fitness Function (der Satzteil „Mutation
(Bestätigung inmitten der Transaktion) → Change fehlt“) und der Begründungssatz in
Festlegung 1 Punkt 3 („`ServerWALEnd` liegt dann hinter Nachrichten, die noch nicht
gespeichert sind“). Alles Übrige von `ADR-0120` bleibt in Kraft, insbesondere die
**Regel** von Festlegung 1 Punkt 3 (inmitten einer Quelltransaktion bestätigt der
Adapter nicht) und die übrigen Zeilen der Fitness Function.

**Datum:** 2026-09-25

**Autor:** Planner-Agent (Modul 8), Closure von `slice-backfill-slot-leerlauf-bestaetigung`;
jede Tatsachenaussage trägt ihren Beleg-Anker (siehe §Kontext).

**Bezug:** [`LH-QA-REL-001`](../../../spec/lastenheft.md) (kein Datenverlust,
Persist-before-ACK; Haupt-Bezug), [`LH-FA-CAP-009`](../../../spec/lastenheft.md),
[`ADR-0120`](0120-capture-slot-leerlauf-bestaetigung.md) (teilweise superseded —
Haupt-Bezug), [`ADR-0011`](0011-persist-before-ack.md) (nur zur Abgrenzung)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0120`](0120-capture-slot-leerlauf-bestaetigung.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Zwei seiner Aussagen tragen eine andere Reichweite
als die Messungen am realen Stream; die Berichtigung ändert den Referenten und ist
keine Zitat-Korrektur.

### Gemessen

**Position eines Keepalive inmitten einer Transaktion (PostgreSQL 18).** Ein
blockierender `Capture`-Stand-in hält die erste Transaktion offen, eine zweite
committet währenddessen 400.000 Änderungen auf die veröffentlichte Tabelle, nach
55 s wird der Stand-in freigegeben; Instanz mit Standard-`wal_sender_timeout`,
Bedingung „keine offene Transaktion“ durch Mutation entfernt. Gedruckt:
`idle P=210904992 open=true confirmed_flush before cancel=210904992`;
`confirmed_flush_lsn` nach dem Abbruch `0/C9227A0`; der Neustart des Streams
liefert `commit=210904992 changes=400000`. Ein Keepalive inmitten der Transaktion
tritt also auf, seine Position ist die Commit-LSN dieser Transaktion, und die
Quelle liefert die Transaktion nach der Bestätigung dieser Position vollständig
(Review-Report `review-slice-backfill-slot-leerlauf-bestaetigung`, F-1 und Abschnitt
„Store-Experiment zu F-1“; der Wegwerf-Test liegt nicht im Repository).

**Bindung der Bedingung im Store-Tier (PostgreSQL 18).** Mutation „Prüfung
`TransactionOpen()` in `confirmIdle` entfernt“, `make test-replication`: Exit 2,
gedruckt allein `--- FAIL: TestRunNoConfirmationInsideOpenTransaction`; alle
Store-Tests einschließlich des Sicherheits-Tests bleiben grün. Mutation „gemeldete
Position `+ (1<<30)`“: `--- FAIL: TestStreamIdleConfirmationKeepsOpenTransactionDeliverable`
(„kein CaptureCommand innerhalb 19.99999958s“) und
`--- FAIL: TestStreamRestartsOnExistingSlot` (Verifikations-Report
`verifikation-slice-backfill-slot-leerlauf-bestaetigung`, §4 S1 und S2; der
Review-Report trägt dieselbe Aussage für S1).

Reports (Verzeichnis `docs/reviews/`): `review-slice-backfill-slot-leerlauf-bestaetigung`,
`verifikation-slice-backfill-slot-leerlauf-bestaetigung`.

### Konstraints

- Persist-before-ACK ([`ADR-0011`](0011-persist-before-ack.md)): eine bestätigte
  Position verdeckt keinen Change, der nicht gespeichert ist; die Bedingung „keine
  offene Transaktion“ hält diese Invariante auf der Seite des Adapters.
- Eine `Accepted`-ADR ist Beleg für nachgeordnete Träger; ihre Fitness-Function-Zeilen
  sind erfüllbar (`AGENTS.md` §3.12 Instanz B).

## Entscheidung

Wir wählen **eine neue ADR mit teilweisem `Supersedes`**, die die zwei Aussagen auf
die gemessene Lage setzt. Zwei Festlegungen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Plan und Report tragen die Messung | keine Änderung an der ADR | die Fitness Function von `ADR-0120` verlangt eine Store-Mutation, die die Messung als nicht erfüllbar zeigt; jeder umsetzende Slice liest sie wörtlich |
| B — die Zeile in `ADR-0120` in-place ändern | ein Ort | ändert §Entscheidung und Fitness Function, keine Zitat-Korrektur (`AGENTS.md` §3.5, [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)) |
| C — einen Store-Test bauen, der die Zeile wörtlich erfüllt | Aussage und Beleg decken sich ohne Berichtigung | nach der Messung liefert die Quelle bei Gleichheit weiter: die Mutation verliert im Store-Tier keinen Change, ein solcher Test ist nicht herstellbar |
| **D — neue ADR mit teilweisem `Supersedes` (gewählt)** | die Aussagen tragen ihre Messung; `ADR-0120` bleibt unberührt | eine weitere ADR für zwei Aussagen |

### Festlegung 1 — Bindung der Bedingung „keine offene Transaktion“

Die Regel von Festlegung 1 Punkt 3 in `ADR-0120` gilt unverändert: inmitten einer
Quelltransaktion bestätigt der Adapter nicht. Ihre Begründung lautet: die Bedingung
hält die Invariante „eine bestätigte Position ist nie größer als der Stand des
Gespeicherten“ auf der Seite des Adapters, unabhängig vom Verhalten der Quelle. Die
Position eines Keepalive inmitten einer Transaktion ist gemessen die Commit-LSN
dieser Transaktion, und die Quelle liefert bei dieser Gleichheit weiter (PostgreSQL
18); die Datensicherheit hängt dort am Verhalten der Quelle, nicht an der Bedingung.
Diese Festlegung ersetzt den Begründungssatz („`ServerWALEnd` liegt dann hinter
Nachrichten, die noch nicht gespeichert sind“) in Festlegung 1 Punkt 3.

### Festlegung 2 — Träger der Bindung je Tier

Die Bedingung „keine offene Transaktion“ trägt der Unit-Test mit der Fake-Sitzung
(`TestRunNoConfirmationInsideOpenTransaction`, Paket `receive`, `make test`);
ihre Mutation färbt genau diesen Test rot. Das Store-Tier bindet die **Größe** der
bestätigten Position: eine Position hinter dem WAL-Ende lässt den Change nach dem
Neustart des Streams ausbleiben (Sicherheits-Test, Mutation „Position `+ 1 GiB`“ rot).
Diese Festlegung ersetzt den Satzteil „Mutation (Bestätigung inmitten der
Transaktion) → Change fehlt“ in der Store-Zeile der Fitness Function von `ADR-0120`.

## Konsequenzen

- Positiv: die Fitness Function trägt nur Zeilen, die eine Messung erfüllt; jede
  Aussage nennt das Tier, das sie bindet.
- Negativ: `ADR-0120` trägt die berichtigten Sätze weiter im Text; wer sie liest,
  liest diese ADR über den Index (Titel „Supers. ADR-0120, teilw.“).
- Negativ (Grenze, benannt): die Aussage zur Position eines Keepalive inmitten einer
  Transaktion stützt sich auf eine einmalige Messung an PostgreSQL 18; kein
  committeter Test bindet die Quellseite, für PostgreSQL 17 liegt die Messung nicht
  vor (der Sicherheits-Test und die Form X1 laufen dort grün,
  `verifikation-slice-backfill-slot-leerlauf-bestaetigung` §1).
- Folgepflicht: keine — Plan und Closure von `slice-backfill-slot-leerlauf-bestaetigung`
  tragen die Aussage.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Fake-Sitzung), Paket `receive` | Keepalive **inmitten** einer Transaktion → keine Bestätigung; Mutation (Prüfung `TransactionOpen()` entfernt) → `TestRunNoConfirmationInsideOpenTransaction` rot | `make test` |
| Go-Test, reale PostgreSQL (Tier `test-replication`) | Sicherheits-Test: die Größe der bestätigten Position; Mutation (gemeldete Position `+ 1 GiB`) → der Change wird nach dem Neustart nicht geliefert | `make test-replication` |

Die übrigen Zeilen der Fitness Function von `ADR-0120` bleiben.

## Re-Evaluierungs-Trigger

- **Eine neue PostgreSQL-Hauptversion oder `pgoutput` mit Streaming großer
  Transaktionen:** die Position eines Keepalive inmitten einer Transaktion neu
  messen (Festlegung 1).
- **Ein committeter Test der Quellseite entsteht:** Festlegung 2 um dessen Tier
  ergänzen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-25 | Accepted — Berichtigung zweier Aussagen von `ADR-0120` an den Messungen des Review- und des Verifikations-Reports zu `slice-backfill-slot-leerlauf-bestaetigung` | [`LH-QA-REL-001`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0121` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
