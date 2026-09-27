# Architect-Verdikt — Ausgangs-Klasse der Fehlerschwelle des WAL-Rückstands

**Datum:** 2026-09-27 · **Stand:** `3f688677` (Baum sauber, Implementer-Commits lokal,
nicht gepusht) · **Rolle:** Architect (frischer Kontext) · **Anlass:** Befund des
Implementers in
[`slice-capture-leerlauf-quellbelege`](../plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md)
§3 „Befund der Erprobung“ — die Runner-Phase „Fehlerschwelle beendet den Container“
läuft grün, der Fehlerzustand trägt aber die Klasse `storage` statt `replication`.
**Vollmacht:** Empfehlung statt Optionsliste, keine Rückfragen an den Auftraggeber.

**Bezug:** [`LH-QA-REL-001`](../../spec/lastenheft.md),
[`LH-QA-REL-003`](../../spec/lastenheft.md),
[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md) (Folgepflicht),
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md),
[`AGENTS.md`](../../AGENTS.md) §3.5, §3.7, §3.12, §3.13.

---

## 0. Ergebnis in einem Blick

- **(1) Codefehler, keine Text-Korrektur.** `ADR-0049` und `SPEC-008` legen für die
  Fehlerschwelle die Klasse `replication` fest, das Benutzerhandbuch sagt sie an vier
  Stellen zu. Die Priorität in `mergeStreamAndWALFaultOutcome` („Stream-Fehler jeder
  Klasse zuerst“) ist breiter als `ADR-0049` (a) verlangt und verdeckt den WAL-Fehler
  durch den Kontext-Abbruch, den die Schwellen-Prüfung selbst auslöst. Keine neue ADR:
  die Entscheidung gilt, der Stand weicht ab (Konflikt-Pfad, erstes Verdikt).
- **(2) Ein kleiner Umsetzungs-Slice** `slice-wal-fehlerschwelle-ausgangsklasse`
  (wellenlos, Umfang S — Schätzung), direkt hinter `slice-capture-leerlauf-quellbelege`
  und vor `slice-start-vorlauf-grenze`. Das Handbuch bleibt unberührt.
- **(3) DoD-Punkt 2 des Slice schließt mit dem Beleg von Ende, Ausgang, Abbruch-Zeile und
  sichtbarem Fehlerzustand**; die Klasse steht als benannte Grenze mit Träger, nicht als
  Zusage (Wortlaut §4).
- **(4) Die Klassenabweichung ist für den Hauptfall repräsentativ** (Persistierung
  stockt), nicht für jeden Stand des Streams: die Klasse hängt heute davon ab, wo der
  Stream im Moment des Abbruchs steht — ein Rennen, kein Vertrag (§5, hergeleitet aus
  dem Quelltext).

Keine `Accepted`-ADR wird berührt; geschrieben sind dieses Verdikt und ein
Register-Beleg. Zwei Nachzüge liegen bei Planner bzw. Auftraggeber (§7).

---

## 1. Was gelesen und was gemessen ist

Alle Zeilen am Stand `3f688677`, 2026-09-27, Docker-frei (Lesen mit `git grep`/`sed -n`).
Keine eigene Docker-Messung: die Klasse `storage` ist eine **übernommene** Messung des
Implementers (§3 des Slice-Plans: fünf Läufe, davon ein Diagnose-Lauf mit vier; nicht
nachgemessen). Alles unten Genannte ist **gelesen** oder ausdrücklich *hergeleitet*.

| Nr | Fundstelle | Lesung |
|---|---|---|
| L1 | `ADR-0049` §Konsequenzen, Folgepflicht | „oberhalb der Fehlerschwelle → kontrollierter Abbruch mit `replication`-Klassifikation“ |
| L2 | `ADR-0049` §Entscheidung (a) | nur die drei Mapper-Sentinels (Stream-Ordnungs-Verletzung) bleiben unverändert hart abbrechend, „unabhängig vom WAL-Rückstand“ |
| L3 | `spec/pflichtenheft.md` Klassen-Tabelle, Zeile `replication` | Transport-/Verbindungsstörung: „Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte Fortsetzung“ — die Klasse der Schwellen-Kette ist `replication` |
| L4 | `internal/bootstrap/wiring.go` `runWALRetentionCheck` (Zeilen 1219–1224) | setzt `fault` mit `outbound.ErrReplication`, ruft danach `stopStream()` |
| L5 | `wiring.go` `mergeStreamAndWALFaultOutcome` (Zeilen 1178–1183) | `streamErr != nil` → `streamErr`, sonst der WAL-Fehler; der Kommentar sagt „jede Klasse“ |
| L6 | `wiring.go` `classifyRunError` (Zeilen 1683–1692) | `ErrReplication` und die drei Mapper-Sentinels → `replication`; `ErrStorage` → `storage` |
| L7 | `internal/application/usecase/capture/service.go` (Zeilen 115–118) | ein Fehler von `PersistTransaction` geht unverändert zurück |
| L8 | `internal/adapters/driven/postgresstorage/store.go` `storageFailure` und `sqlexec/errors.go` `Classify` | Treiber-Fehler → `fmt.Errorf("%w: %w", ErrStorage, cause)`; die Ursache bleibt in der Kette |
| L9 | `internal/adapters/driving/replication/receive/receive.go` (Zeilen 396–402) | `ReceiveMessage` mit beendetem Kontext → `return nil` (regulärer Abschluss); `process` (Zeilen 551–556) gibt einen Capture-Fehler unverändert zurück |
| L10 | `wiring.go` `reportFault` und `outbound.HeartbeatPort.Fault` (`heartbeat.go` Zeile 43) | der Fehlerzustand trägt **nur die Klasse**, keinen Text |
| L11 | `docs/user/benutzerhandbuch.md`, Zeilen 464–466, 866–868, 1526, 1614–1615 | vier Stellen sagen `replication` für die Fehlerschwelle zu |
| L12 | `internal/bootstrap/walretention_internal_test.go`, `walretention_slotgrowth_internal_test.go` | die Kette ist in Teilen belegt: Sentinel-Priorität (Ordnungs-Verletzung gewinnt) und `merge(nil, fault)` → `replication`; **kein** Test führt einen Stream-Fehler der Klasse `storage` neben einem gesetzten WAL-Fehler |

---

## 2. Verdikt (1) — Codefehler; die Entscheidung gilt

**Verdikt:** Die Klasse `storage` ist in diesem Aufbau ein Fehler des Codes, nicht des
Textes. Der Plan des Slice hat nichts falsch behauptet; die Träger (Handbuch, Kommentare)
geben `ADR-0049` richtig wieder, der Code weicht von ihr ab. Adresse: neuer Slice (§3).

**Begründung.**

1. Die Norm ist eindeutig (L1, L3): dieselbe Ursache — WAL-Rückstand über der
   Fehlerschwelle — führt zur Klasse `replication`. Die Klasse ist der Schlüssel, nach
   dem ein Betreiber alarmiert (der Fehlerzustand trägt nur sie, L10); zwei Klassen für
   dieselbe Ursache machen den Alarm zufällig.
2. Die Abweichung ist ein Rennen, kein Zustand: der Kontext-Abbruch, den die
   Schwellen-Prüfung selbst auslöst (L4), erreicht den Stream je nach Standort als
   `nil` (in `ReceiveMessage`, L9 → `replication`), als Klasse `storage` (in
   `PersistTransaction`, L7/L8) oder als `replication` (in der Bestätigung, `ErrReplication`
   aus dem ACK-Adapter). Ein Vertrag „die Klasse hängt vom Standort des Streams ab“
   wäre kein Vertrag.
3. Die Text-Korrektur wäre die teurere Lösung: `ADR-0049`s Folgepflicht,
   `SPEC-008`, vier Handbuch-Stellen und die Kommentare müssten die Zufälligkeit
   beschreiben; das Handbuch trägt Versionshistorie. Die Code-Korrektur ist eine Regel in
   einer kleinen Funktion (L5), und alle bestehenden Aussagen bleiben wahr.
4. Die Priorität „Stream-Fehler jeder Klasse zuerst“ (L5) ist breiter als die Norm (L2):
   `ADR-0049` (a) schützt die **Stream-Ordnungs-Verletzung**, nicht jeden Stream-Fehler.
   Ein Stream-Fehler, der die **Folge** des Abbruchs durch die Schwellen-Prüfung ist,
   trägt keine eigene Information.

**Die Regel für den Umsetzungs-Slice** (Vertrag, nicht Code): Ist ein WAL-Schwellen-Fehler
gesetzt, so ist der Ausgang von `Run`

| Stream-Ausgang | Ausgang von `Run` | Stand |
|---|---|---|
| `nil` | WAL-Fehler | unverändert, belegt (L12) |
| Ordnungs-Sentinel (`mapper.ErrChangeWithoutBegin`, `…CommitWithoutBegin`, `…BeginWithoutCommit`) | der Sentinel | unverändert, belegt (L12) |
| Fehler, dessen Kette `context.Canceled` trägt (Folge des Abbruchs; heute `storage`) | WAL-Fehler | **neu** |
| jeder andere Stream-Fehler (echter Persistenzfehler ohne Abbruch-Folge, Ack-Fehler) | der Stream-Fehler | unverändert |

Ohne gesetzten WAL-Fehler ändert sich nichts: die Klasse eines echten Persistenzfehlers
bleibt `storage`, `SPEC-008` (`storage`: kein Source-ACK) und `ADR-0023` bleiben
unberührt. Der Fehlertext des Stream-Endes („Stream beendet mit Fehler … context
canceled“) bleibt im Log neben der Abbruch-Zeile stehen: der Betreiber liest die
**Klasse** im Fehlerzustand, die **Ursache** im Log (beide Zeilen tragen sie); mehr
Beobachtbarkeit braucht der Fall nicht.

*Zur Kette `context.Canceled`* (**hergeleitet**, nicht erprobt): `Classify` verwendet
`%w: %w` (L8), und die Meldung des Befunds endet auf „context canceled“; dass die
Treiber-Kette `errors.Is(err, context.Canceled)` an der echten Stelle trägt, ist nicht
gelesen. Der Slice belegt es zweifach: ein Unit-Test mit einer Kette, die die Form der
echten trägt, und die Runner-Phase (sie färbt sich rot, wenn die Regel die Klasse nicht
erreicht). Trägt die Prüfung die Kette nicht, ist `streamCtx.Err() != nil` bei lebendem
Prozess-Kontext der zulässige Ersatz derselben Regel — kein weiterer Architect-Zug.

---

## 3. Verdikt (2) — Slice-Vorschlag an den Planner

**Kennung:** `slice-wal-fehlerschwelle-ausgangsklasse` (wellenlos, Umfang S —
Schätzung: eine Funktion, ein Test, eine Assertion im Runner).

**Start-Trigger:** `slice-capture-leerlauf-quellbelege` liegt in `done/`. Grund: beide
ändern die Runner-Phase „Fehlerschwelle beendet den Container“ in
`tools/harness/run-integration-tests.sh` (WIP-Limit 1); die Phase ist die Falsifikation
des Slice.

**Reihenfolge:** `slice-capture-leerlauf-quellbelege` →
`slice-wal-fehlerschwelle-ausgangsklasse` → `slice-start-vorlauf-grenze` →
`slice-transformationen-e2e-abhilfe` → `slice-transformationen-betriebsdoku`. Vor
`slice-start-vorlauf-grenze`, weil (a) der Runner-Beleg frisch und der Umfang klein ist,
(b) beide `wiring.go` und den Runner berühren, aber an entgegengesetzten Enden —
`slice-start-vorlauf-grenze` `Stream.Run`/`START_REPLICATION` und den Start-Pfad, dieser
Slice `mergeStreamAndWALFaultOutcome` nach der Rückkehr von `stream.Run` — und der
kleinere zuerst die Prüfspur des größeren nicht verstellt. Die Kante trägt der Planner als
Start-Trigger in `slice-start-vorlauf-grenze` §4 und in
[`welle-transformationen`](../plan/planning/welle-transformationen.md) §5 nach; die
Trigger von `e2e-abhilfe` bleiben (er wartet auf `slice-start-vorlauf-grenze`).

**DoD-Kern:**

1. `mergeStreamAndWALFaultOutcome` folgt der Regel aus §2; ein Tabellentest über alle
   vier Zeilen, dazu ein Fall mit **echt gewrappter** Kette
   (`fmt.Errorf("%w: %w", outbound.ErrStorage, <Fehler mit context.Canceled>)`) und ein
   Fall „`storage` ohne Abbruch-Folge neben gesetztem WAL-Fehler bleibt `storage`“.
   Mutation: die alte Priorität (`streamErr != nil` zuerst) zurück — der Fall
   „Abbruch-Folge → WAL-Fehler“ färbt sich rot (*erwartet*; Stelle, Instanz und Farbe
   trägt der Implementer nach, `AGENTS.md` §3.12).
2. Die Runner-Phase „Fehlerschwelle beendet den Container“ prüft die Klasse
   `replication` als **Zusage** (`cdc.process_heartbeat.error_class`); die
   Abdeckungs-Zeile in `docs/user/e2e-abdeckung.md` (Erzeugnis des Runners) nennt sie.
   Mutation: die Regel im Code zurücknehmen — die Phase färbt sich rot (*erwartet, zu
   erproben*: eine Stelle, ein Lauf).
3. Kommentare an ihre Träger: der Kommentar an `mergeStreamAndWALFaultOutcome`
   (Zeilen 1168–1177: „jede Klasse“, „nur zum Zug, wenn der Stream-Lauf regulär endete“),
   der Kommentar von `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` und der
   Kommentar in `walretention_slotgrowth_internal_test.go` (Zeile 27, die Kette „Fehler bei
   Stream-Ende“ trägt das Prozessende) — je höchstens eine Kennung, kein Konjunktiv
   (`AGENTS.md` §3.7).
4. Die benannte Grenze in `harness/README.md` §Sensors bei `make test-integration`
   entfällt; das Handbuch bleibt unberührt (**gelesen**: L11 sagt `replication` zu, jede
   Stelle bleibt wahr).
5. §3.13-Suchlauf im Plan, `make gates` grün, Review.

**Abgrenzung:** keine Änderung an den Schwellen, an `classifyRunError`, an der
Sentinel-Trennung oder an der Klassen-Tabelle; kein Betreiber-Weg. Keine neue ADR:
`ADR-0049` wird eingelöst, nicht geändert.

---

## 4. Verdikt (3) — DoD-Punkt 2 von `slice-capture-leerlauf-quellbelege`

**Verdikt:** Der Punkt schließt mit dem, was die Phase trägt; die Klasse steht als
benannte Grenze mit Adresse. Die vorhandene Klausel „Der Ausgang ‚Grenze bleibt‘ ist
zulässig“ deckt das (sie nennt die Grenze und ihren Träger in `harness/README.md`); der
Punkt braucht nur den geänderten Wortlaut. Wortlaut zur Übernahme durch den Planner:

> - [x] Der Beleg „Fehlerschwelle erreicht → Container endet“ steht als Runner-Phase
>   „Fehlerschwelle beendet den Container“ in `make test-integration`: bei gehaltener
>   Persistierung und WAL ohne Inhalt für die Publication über der Fehlerschwelle endet
>   der Feed-Container mit Ausgang 1, das Log trägt die Abbruch-Zeile mit einem
>   Rückstand über der Fehlerschwelle, und `cdc.process_heartbeat` trägt einen
>   Fehlerzustand. Die **Klasse** des Ausgangs ist nicht Teil der Zusage der Phase: sie
>   ist `storage` gemessen (fünf Läufe, Ausgabezeile der Phase), `ADR-0049` legt
>   `replication` fest — Codefehler, Träger `slice-wal-fehlerschwelle-ausgangsklasse`
>   (Architect-Verdikt `architect-verdict-wal-fehlerschwelle-ausgangsklasse`; den Link
>   setzt der Planner in der Form seiner Plan-Datei).
>   *Zu belegen durch:* ein realer, grüner `make test-integration`-Lauf mit der
>   Abdeckungs-Zeile in `docs/user/e2e-abdeckung.md` und die benannte Grenze in
>   `harness/README.md` §Sensors; die Mutation „Aufruf von `stopStream` in der
>   Schwellen-Prüfung entfernt“ färbt die Phase rot — nur dann als **erprobt** zu führen,
>   wenn der Implementer sie gefahren hat (Stelle, Instanz, gesehene Farbe), sonst als
>   *hergeleitet* (der Container endet dann nicht; die Phase wartet 90 s und meldet den
>   Fehler).

Der **Träger der Grenze** ist `harness/README.md` §Sensors bei `make test-integration`
und der Slice aus §3, dessen DoD sie streicht. Die Ausgabezeile des Runners, die die
Klasse nennt, bleibt: sie ist die Messung, die die Grenze trägt. Der Text der Grenze in
`harness/README.md` nennt heute „der Befund liegt beim Architect“; **Nachzug** (§7):
ersetzen durch „Codefehler, Träger `slice-wal-fehlerschwelle-ausgangsklasse`“.

Für die Closure-Notiz von `slice-capture-leerlauf-quellbelege` (Vorschlag): §6-Risiko
„Aufbau instabil“ — **entfallen** (Aufbau stabil, fünf von fünf Läufe); *Was ging anders*
— die Kette Fehlerschwelle → Prozessende war in Teilen belegt (Unit-Tests mit Fakes), die
Klasse erschien erst am komponierten Prozess; Lerneintrag — §6, Register-Beleg
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`. Die ADR-Ergänzung von `ADR-0121`
Festlegung 2 (Keepalive-Tier, PostgreSQL 17) hängt an diesem Befund **nicht** und bleibt
im geplanten Zug nach der Verifikation.

---

## 5. Verdikt (4) — Messgrundlage und Repräsentativität

**Verdikt:** Fünf von fünf Läufen (ein Diagnose-Lauf: vier) belegen die Klasse `storage`
**in dem Aufbau, den die Phase fährt**: der Stream steht in `PersistTransaction`, die an
der exklusiven Sperre auf `cdc.change` wartet. Der Aufbau ist kein Kunstgriff: die reale
Ursache eines wachsenden Rückstands ist eine stockende oder zu langsame Persistierung
(Sperre, Wartung, Lastspitze) — dann steht der Stream nahezu dauernd in `Capture`, der
Hauptfall trifft die Klasse `storage` (*hergeleitet*). **Nicht** repräsentativ ist der
Aufbau für jeden Stand des Streams; die folgende Abbildung ist aus dem Quelltext
**hergeleitet**, außer der ersten Zeile:

| Stand des Streams beim Abbruch | Ausgang des Stream-Laufs | Klasse heute |
|---|---|---|
| in `PersistTransaction` (wartet an Sperre) | Fehler mit `ErrStorage` (L7/L8) | `storage` — **erprobt**, 5 von 5 im Aufbau der Phase |
| in `Acknowledge` | `ErrReplication` aus dem ACK-Adapter (`replicationFailure`) | `replication` |
| in `ReceiveMessage` (Leerlauf, offene Quelltransaktion) | `nil` (L9) | `replication` |
| in der Leerlauf-Bestätigung | `ErrReplication` (`receive.go` Zeile 483) | `replication` |
| in `Assembler.Consume` → `observeRelation` (Schema-Store mit `ctx`, `mapper.go` Zeile 388) | Fehler des Schema-Stores, Klasse ungelesen | offen (*hergeleitet*: `storage`-nah) |

Ein parallel laufender Backfill-Run ändert den Standort des Stream-Laufs nicht (eigene
Goroutine, eigene Verbindung); wächst der Rückstand dort bei ruhendem Stream, endet der
Lauf über `ReceiveMessage` mit `replication`. Die Abbildung zeigt: die heutige Klasse
folgt dem Standort — genau die Abhängigkeit, die die Regel aus §2 beseitigt. Eine weitere
Messung ist unnötig; die Erprobung der Regel trägt die Runner-Phase des Slice (§3, Punkt
2), und sie prüft den Fall, der zuvor `storage` war.

**Akzeptiertes Negativ:** die Zeile „`observeRelation`“ wird nicht gesondert erprobt; sie
fällt unter die Regel, sobald ihre Kette `context.Canceled` trägt (Zeile 3 der Tabelle in
§2), und ein Abbruch dort ist ein Randfall des Randfalls.

---

## 6. Register

`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (Klasse „Unit-Test mit Fakes verdeckt
eine Lücke der realen Verdrahtung“) trägt den Fall als fünften Beleg: die Sentinel-
Priorität und die Klassen-Abbildung sind je mit Fakes belegt (L12), die Ausgangs-Klasse
am komponierten Prozess erschien erst in der Runner-Phase. Der Träger dieser Klasse ist
nach dem Eintrag die Phase in `make test-integration` im Slice, der die Eigenschaft
einführt — hier trägt sie der Slice aus §3. Der Ausgang des Eintrags (`geplant` →
`slice-start-vorlauf-grenze`) bleibt unverändert. Kein zweiter Eintrag:
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (die Kommentare an
`mergeStreamAndWALFaultOutcome` und in Zeile 27 des Tests) wird durch DoD-Punkt 3 des Slice
getragen; ein siebter Beleg brächte dem Zähler nichts, weil der Implementer den Befund
selbst gefunden hat.

## 7. Was nicht getan wurde — und die Nachzüge

Kein Produktionscode, keine Änderung an `ADR-0049`/`ADR-0121`, an
`docs/user/benutzerhandbuch.md`, an `spec/**` und am Plan des Slice. Geschrieben: dieses
Verdikt und der Register-Beleg.

Offen bei anderen Rollen:

1. **Planner:** die DoD-Zeile aus §4 in `slice-capture-leerlauf-quellbelege`, die
   Kanten-Nachzüge aus §3, das Anlegen des Slice-Plans (Adresse; Frist: Closure von
   `slice-capture-leerlauf-quellbelege`).
2. **Text der benannten Grenze in `harness/README.md`** (§4): das Werkzeug dieses Zugs
   erlaubte keine Teil-Änderung der Datei; der Ersatz steht in §4 und wird mit dem
   Nachzug des Planners oder des Auftraggebers eingesetzt.
