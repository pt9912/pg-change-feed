# ADR-0135: Capture — `transient`-Fehler mit begrenztem Backoff am Stream-Zyklus wiederholen

**Status:** Accepted

**Datum:** 2026-09-29

**Autor:** pt9912

**Bezug:** [`LH-QA-REL-002`](../../../spec/lastenheft.md) (kontrollierter
Neustart), [`LH-QA-REL-001`](../../../spec/lastenheft.md) (keine stillen
Datenverluste), [`LH-QA-REL-003`](../../../spec/lastenheft.md) (sichtbarer
Unzuverlässigkeitszustand), [`LH-FA-ADM-003`](../../../spec/lastenheft.md)
(sichtbare Fehlerzustände), [ADR-0023](0023-fehlerklassifikation.md),
[ADR-0049](0049-replication-fehlerklassen-schwellen.md),
[ADR-0011](0011-persist-before-ack.md), [ADR-0012](0012-at-least-once.md),
[ADR-0026](0026-composition-root.md)

**Schärft:** [`SPEC-008`](../../../spec/pflichtenheft.md) (Zeile `transient`
— „Erneut versuchen mit begrenztem Backoff")

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Capture-Pfad endet heute auf jeden Adapter-Fehler mit Prozess-Ausgang 1:
`Run` in der Composition Root (`internal/bootstrap/wiring.go`) reicht jeden
Adapter-Fehler durch, `reportFault` trägt den Fehlerzustand in den Heartbeat
([`LH-FA-ADM-003`](../../../spec/lastenheft.md)), und die Fortsetzung trägt
der Prozess-Neustart durch den Aufrufer
([`LH-QA-REL-002`](../../../spec/lastenheft.md); Handbuch „Neustart nach
einem Fehler", `restart: "no"` im mitgelieferten `compose.yaml`). Die
`transient`-Aktion der Fehlerklassen-Tabelle ([`SPEC-008`](../../../spec/pflichtenheft.md):
„Erneut versuchen mit begrenztem Backoff") trägt dieser Pfad nicht — der
Kommentar an `Run` nennt sie ausdrücklich als nicht getragen. Die Klasse ist
im Domänenmodell als `ErrorClassTransient` angelegt, aber kein Pfad des
Capture-Laufs erzeugt sie.

[ADR-0049](0049-replication-fehlerklassen-schwellen.md) (a) hat innerhalb der
Klasse `replication` getrennt: Stream-Ordnungs-Verletzungen
(`mapper.ErrChangeWithoutBegin`, `mapper.ErrCommitWithoutBegin`,
`mapper.ErrBeginWithoutCommit`) bleiben hart abbrechend,
Transport-/Verbindungsstörungen (`receive.ErrReplication`,
`outbound.ErrReplication`) sind der Schwellen-Kandidat für die
WAL-Rückstand-Überwachung. Diese Festlegung bleibt unberührt; diese ADR
entscheidet, welcher Fehlerkreis die `transient`-Aktion der Klasse
`transient` bekommt und wo die Wiederholung liegt.

[ADR-0023](0023-fehlerklassifikation.md) ordnet die Strategie je Klasse der
Anwendung zu („Application entscheidet je Klasse … `transient` → Backoff");
die Adapter übersetzen technische Fehler nur in die sieben Kategorien.
[ADR-0128](0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
(Festlegung 1) verlegt den `START_REPLICATION`-Aufruf in `Stream.Run` — der
Fall „Slot noch aktiv" (SQLSTATE 55006) entsteht damit im Stream-Lauf, nicht
im Verbindungsaufbau.

Auslöser: das Architect-Verdikt
`architect-verdict-welle-backfill-bestand-lese-schritt` §8 (Slice C) und die
Beobachtung `BEO-PGC/adapter-fehler-ausgang` (3×). Diese ADR ist der benannte
Start-Trigger des Slices `slice-capture-transient-wiederholung`.

## Entscheidung

Wir treffen sechs Festlegungen:

**1. Ort — Composition Root; Wiederholeinheit ist der Stream-Zyklus.** Die
Composition Root ([ADR-0026](0026-composition-root.md)) umschließt den
Stream-Lauf mit einer begrenzten Wiederholungsschleife. Wiederholeinheit ist
der **Stream-Zyklus** — Neuaufbau von Stream-Verbindung, ACK-Adapter und
Bindung und ein neuer `stream.Run`-Aufruf. Der Start-Vorlauf der
Administrations-Anträge ([ADR-0128](0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md))
läuft einmal pro Prozess, nicht je Versuch. Die Adapter tragen keine
Wiederholung: der Stream-Adapter meldet Fehler, der Store-Adapter persistiert,
die Wiederholungspolitik liegt oberhalb der Adapter
([ADR-0023](0023-fehlerklassifikation.md)). Der Capture Service bleibt
unverändert — er sieht Verbindungsstörungen der Quelle strukturell nie, und
sein Aufrufvertrag (Persist-before-ACK) wird nicht je Versuch wiederholt.

**2. Backoff-Grenzen (Startwerte).** Anfangsverzögerung **2 s**, Verdopplung
je Fehlversuch, Obergrenze **30 s** je Warteschritt, Gesamtfenster
**5 Minuten** ab dem ersten wiederholten Fehler einer Episode. Ein
erfolgreicher Stream-Zyklus setzt die Episode zurück. Die Werte sind
**Setzungen ohne Messung** — derselbe Stil wie die Schwellenwerte in
[ADR-0049](0049-replication-fehlerklassen-schwellen.md) (b): Initialwerte,
über eine Folge-ADR schärfbar, kein neuer Konfigurationsschlüssel.

**3. Wiederholte Fehlermenge.** Wiederholt wird genau die
**Transport-/Verbindungsstörung am Quellzugriff** nach
[ADR-0049](0049-replication-fehlerklassen-schwellen.md) (a): ein Fehler, dessen
Kette `receive.ErrReplication` oder `outbound.ErrReplication` trägt —
Verbindungsaufbau, `START_REPLICATION` (einschließlich SQLSTATE 55006 „Slot
noch aktiv"), Katalog- und Slot-Abfragen, Keepalive, Nachrichtenempfang und
Quell-Bestätigung. **Nicht** wiederholt — der Prozess endet weiterhin mit
Ausgang 1 — werden:

- die Klasse `configuration` (Sentinels `ErrConfiguration` der Verdrahtung
  und des Stream-Adapters) — eine falsche Vorbedingung heilt nicht durch
  Warten;
- die Klasse `permission` — eine fehlende Berechtigung gewinnt Warten nicht;
- die Klasse `schema` (`decode.ErrSchema`, `mapper.ErrTruncateUnsupported`,
  `mapper.ErrIncompatibleSchemaChange`,
  `mapper.ErrTransformationNotApplicable`) — eine nicht sicher
  interpretierbare Änderung ist kein Zustand, in dem fortgesetzt wird;
- die Stream-Ordnungs-Verletzung (`mapper.ErrChangeWithoutBegin`,
  `mapper.ErrCommitWithoutBegin`, `mapper.ErrBeginWithoutCommit`,
  [ADR-0049](0049-replication-fehlerklassen-schwellen.md) a);
- die Klasse `storage` (Persistenzfehler im ChangeStore) — die Zusage
  „**Kein Source-ACK**" ([`LH-QA-REL-001`](../../../spec/lastenheft.md),
  [`SPEC-008`](../../../spec/pflichtenheft.md)) bleibt bindend, eine
  Wiederholung dort ist nicht Teil dieses Entscheidungsrahmens. Im MVP
  bleiben Quelle und CDC-Speicher dieselbe Instanz — eine Störung der
  Instanz zeigt sich an der Quellseite (Replication-Verbindung) und wird von
  der Wiederholung erfasst; eine Störung, die ausschließlich den
  Speicherzugriff träfe, setzt den heutigen Ausgang fort (Ausgang 1, kein
  ACK, Neustart des Aufrufers).

**4. Sichtbarkeit während der Wiederholung.** Je Fehlversuch trägt die
Composition Root einen strukturierten WARN-Log-Eintrag (Versuchszähler,
Warteschritt, Fehlertext), je erfolgreicher Fortsetzung einen INFO-Eintrag.
Der Heartbeat schreibt während der Wiederholung keinen Fehlerzustand — sein
periodisches Lebenszeichen bleibt unverändert, der `--healthcheck` stuft den
Container in dem Fenster nicht als unhealthy. Operative Sichtachsen in dem
Fenster sind Log, `cdc_capture_lag` und `cdc_wal_retention_bytes` (dasselbe
Muster wie die Warnstufe des WAL-Rückstands). Erschöpft die Grenze, schreibt
`reportFault` den Fehlerzustand der Klasse `transient` in den Heartbeat
([`LH-FA-ADM-003`](../../../spec/lastenheft.md)) — die Klasse wird durch
diese Entscheidung erstmals durch einen Lauf des Capture-Pfads erzeugt.

**5. Ausgang bei Erschöpfung.** Der Prozess endet mit Ausgang 1 und der
Fehlerklasse `transient` im Heartbeat-Fehlerzustand — wie der heutige
Adapter-Fehler. Die Fortsetzung trägt danach unverändert der Neustart des
Aufrufers ([`LH-QA-REL-002`](../../../spec/lastenheft.md)). Die Wiederholung
ersetzt den Aufrufer-Neustart nicht; sie verlängert die Toleranz kurzer
Störungen **vor** ihm.

**6. Verhältnis zu `restart: "no"`, Persist-before-ACK und Abgrenzung.**
`restart: "no"` bleibt im Container-Vertrag; diese Entscheidung ordnet keinen
Supervisor im Container an und ändert nichts am Neustart-Vertrag des
Aufrufers. Die Persist-before-ACK-Zusage
([ADR-0011](0011-persist-before-ack.md), [`LH-QA-REL-001`](../../../spec/lastenheft.md))
ist von der Wiederholung **nicht berührt**: Jeder Versuch liest die Quelle an
der serverseitig bestätigten Position (`confirmed_flush_lsn`,
[ADR-0012](0012-at-least-once.md)) — ein Vorgänger-Durchlauf, dessen
Persistierung scheiterte oder dessen Bestätigung nicht ankam, ist nicht
bestätigt, der Slot liefert die Transaktion erneut, und ein ACK liegt stets
hinter der Persistierung. Die Dopplungs-Semantik eines Versuchs ist damit
dieselbe wie die des heutigen Prozess-Neustarts (at-least-once). Die
Abgrenzung des Slices bestätigt sich: kein Restart-Supervisor, keine
`restart:`-Änderung, keine Wiederholung für `permission`/`configuration`/
`schema`/`storage`, keine Wiederholung im Backfill-Run
([ADR-0111](0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4/5 — ein
`transient`-Fehler des Runs endet ihn `failed`, run-lokal).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: jeder Adapter-Fehler endet mit Ausgang 1, der Aufrufer startet neu (Status quo) | kein Code; ein Mechanismus (Neustart) für alle Fälle | die `transient`-Aktion von [`SPEC-008`](../../../spec/pflichtenheft.md) bleibt unerfüllt; eine Störung von Sekunden bricht den Container und kostet den Neustart-Zyklus samt Start-Vorlauf ([ADR-0128](0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)) und Monitoring-Signal |
| B — Wiederholung im Stream-Adapter `receive` | kleinster Zug; die Verbindung gehört ohnehin dem Adapter | der Adapter entschied eine Politik, die [ADR-0023](0023-fehlerklassifikation.md) der Anwendung zuordnet; die Klassifikation und Sichtbarkeit bräuchten einen neuen Rückpfad aus dem Adapter in die Composition Root; die Speicherseite bliebe unerreichbar |
| C — Wiederholung im Capture Service (je Persist-Aufruf) | kleinster Eingriff an der Persistenzstelle | ein wiederholter Persist-Aufruf nach mehrdeutigem Ergebnis (Persistierung bestätigt, Antwort verloren) würde dieselbe Transaktion erneut persistieren, ohne dass der Slot sie neu liefert — genau die Persist-before-ACK-Rückführungsfrage; Verbindungsstörungen der Quelle sieht der Service strukturell nie |
| **D — Composition Root, Wiederholeinheit Stream-Zyklus (gewählt)** | ein Ort, der Klassifikation, Sichtbarkeit und Erschöpfung heute schon trägt (`classifyRunError`, `reportFault`); „kein ACK vor Persistierung" bleibt strukturell unangetastet, weil jeder Versuch am serverseitig bestätigten Stand beginnt; der Adapter bleibt politikfrei | der Zug reicht in die Verdrahtung (Stream-Neuaufbau je Versuch); die Grenzwerte sind Setzungen ohne Messung |

## Konsequenzen

- Positiv: [`SPEC-008`](../../../spec/pflichtenheft.md)s Zusage „Erneut
  versuchen mit begrenztem Backoff" für die Klasse `transient` bekommt einen
  konkreten, umsetzbaren Fehlerkreis (Transport-/Verbindungsstörung am
  Quellzugriff) und eine konkrete Wiederholeinheit (Stream-Zyklus) statt der
  Auslegung eines Implementers. Die Klasse `transient` wird im
  Heartbeat-Fehlerzustand erzeugbar. Kurze Störungen (Container-Tausch,
  Slot-Überlappung mit SQLSTATE 55006, Instanz-Neustart) beendet den
  Container nicht mehr beim ersten Fehler.
- Negativ: Ein Dauerfehler, der wie `transient` aussieht, wird bis zum
  Fensterablauf wiederholt statt sofort gemeldet — die Trennung der
  Fehlermenge (Festlegung 3) und die Negativtests je Klasse tragen dagegen.
  Die Grenzwerte sind Setzungen, keine Messwerte.
- Folgepflicht: die Träger folgen — der Kommentar an `Run` in
  `internal/bootstrap/wiring.go` nennt die `transient`-Aktion als getragen,
  die Container-Vertrags-Kommentar-Zeile in `compose.yaml`, der
  Handbuch-Abschnitt „Neustart nach einem Fehler" und die Zeile `transient`
  in [`SPEC-008`](../../../spec/pflichtenheft.md) (Bedingung der Festlegung
  3); je Backoff-Grenze ein Test. Die Umsetzung trägt
  `slice-capture-transient-wiederholung`.

## Fitness Function (falls maschinell prüfbar)

Die Zeilen sind **Zusagen** des Slices, keine Erprobungen — auf dem
Arbeitsbaum vor der Umsetzung existiert kein Test, der sie trägt; die
Mutationen sind an der Eingabeseite (verschobene Grenze) vorgesehen und der
Implementer fährt sie.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Fakes, deterministische Uhr) | je Grenze der Festlegung 2 (Anfangsverzögerung, Backoff-Folge, Erschöpfung) ein Test; ein `transient`-Fehler führt zu Wiederholungen mit Fortsetzung an `confirmed_flush_lsn`; Erschöpfung endet mit Klasse `transient` im Heartbeat-Fehlerzustand | `make test` (Race-Detector) |
| Go-Test | je Klasse der Festlegung 3 (`configuration`, `permission`, `schema`, Stream-Ordnungs-Verletzung, `storage`) ein Negativtest an seine Eingabe gebunden; Mutation „Klasse als wiederholbar behandeln" färbt rot | `make test` |
| Reale PostgreSQL | Fall „Slot noch aktiv" (SQLSTATE 55006): der Stream-Zyklus wiederholt statt zu beenden, die Fortsetzung liest an `confirmed_flush_lsn` | `make test-replication` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

- Für die Grenzwerte (Festlegung 2): **nicht permanent.** Zeigt der
  E2E-Lauf oder Betrieb, dass das Fenster zu grob oder zu fein greift (der
  Test scheitert an den Startwerten, oder reale Störungsdauern weichen
  deutlich ab), schärft eine Folge-ADR mit `Supersedes ADR-0135` die Werte.
- Für Ort und Wiederholeinheit (Festlegung 1) und für die Wiederholte
  Fehlermenge (Festlegung 3): **permanent, solange Quelle und CDC-Speicher
  dieselbe Instanz bleiben.** Trennt sich die Topologie (getrennte
  Instanzen) oder zeigt realer Betrieb eine isolierte Speicher-Störung
  (Quelle erreichbar, Speicher nicht), ist die Aufnahme der Speicherseite in
  die Wiederholungsmenge erneut zu entscheiden.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-29 | Accepted | `docs/plan/planning/open/slice-capture-transient-wiederholung.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0135` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
