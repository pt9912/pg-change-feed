# Slice wal-fehlerschwelle-ausgangsklasse: Fehlerschwelle des WAL-Rückstands — der Ausgang von `Run` trägt die Klasse `replication`, auch wenn der Abbruch den Stream mit einem Kontext-Fehler beendet

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er startet nach `slice-capture-leerlauf-quellbelege`
und `slice-leerlauf-phase-last-in-stuecken`
und geht `slice-start-vorlauf-grenze` voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5, Kante zu
`slice-transformationen-e2e-abhilfe` über `slice-start-vorlauf-grenze`).

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein
Datenverlust; der Abbruch ist sichtbar),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand sichtbar,
Schwellen),
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Folgepflicht:
oberhalb der Fehlerschwelle kontrollierter Abbruch mit `replication`-Klassifikation;
Entscheidung (a): nur die drei Stream-Ordnungs-Sentinels bleiben unabhängig vom
WAL-Rückstand hart abbrechend),
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (Fehlerklassen),
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) (Leerlauf-Bestätigung;
Kontext der Fehlerschwellen-Kette), Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
§2 bis §4.

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklassen; Zeile `replication`) und
[`SPEC-013`](../../../../spec/pflichtenheft.md) (Schwellen des WAL-Rückstands) —
gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Planner-Zug nach dem Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md).
**Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Erreicht der WAL-Rückstand die Fehlerschwelle, trägt der Fehlerzustand des
Feed-Containers die Klasse `replication` — unabhängig davon, wo der Stream im Moment
des Abbruchs steht: `mergeStreamAndWALFaultOutcome` (`internal/bootstrap/wiring.go`)
gibt bei gesetztem WAL-Fehler den WAL-Fehler zurück, wenn der Stream-Fehler die Folge
des Abbruchs ist, den die Schwellen-Prüfung selbst auslöst
([`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md), Verdikt §2).

**Kernaussage und ihre Herkunft** (**übernommen** aus Verdikt §1 und §5, nicht vom
Planner gemessen): bei gehaltener Persistierung endet der Feed-Container mit Ausgang 1
und der Klasse `storage` statt `replication` — gemessen vom Implementer von
`slice-capture-leerlauf-quellbelege` in fünf von fünf Läufen (drei Läufe der Phase in
einem Wegwerf-Aufbau, ein vollständiger `make test-integration`, ein Diagnose-Lauf der
Vorfassung; Anker: Plan dort §3 „Befund der Erprobung“, Ausgabezeile der Runner-Phase
„Fehlerschwelle beendet den Container“). Die Abbildung „Stand des Streams beim Abbruch →
Klasse“ (Verdikt §5) und die Ursache — der Kontext-Abbruch lässt `Capture` mit einem
Fehler der Klasse `storage` zurückkehren, und die Priorität „Stream-Fehler jeder Klasse
zuerst“ gibt ihn vor dem WAL-Fehler zurück — sind aus dem Quelltext **hergeleitet**, nicht
erprobt. Dass `errors.Is(err, context.Canceled)` an der echten Stelle trägt, ist
**hergeleitet** (die Meldung des Befunds endet auf „context canceled“, `Classify` wrappt
mit `%w: %w`) und nicht gelesen; der Slice belegt es zweifach (§2, Punkt 1 und 2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung an den Schwellen, an `classifyRunError`, an der Sentinel-Trennung oder an
  der Klassen-Tabelle.** `ADR-0049` wird eingelöst, nicht geändert; die Norm ist
  eindeutig, der Stand weicht ab (Verdikt §2, Konflikt-Pfad). Keine neue ADR.
- **Das Benutzerhandbuch.** Die vier Stellen sagen `replication` für die Fehlerschwelle
  zu und bleiben mit der Korrektur wahr (§3, Suchlauf: sechs Treffer gelesen); eine
  Text-Korrektur wäre die teurere Lösung gewesen (Verdikt §2, Begründung 3).
- **Ein Betreiber-Weg, die Persistierung anzuhalten.** Der Beleg nutzt den Aufbau
  der Runner-Phase; kein Produktpfad entsteht.
- **Die Zeile „`observeRelation`“ der Abbildung in Verdikt §5** (Abbruch im
  Schema-Store). Akzeptiertes Negativ des Verdikts: sie fällt unter die Regel, sobald
  ihre Kette `context.Canceled` trägt, und wird nicht gesondert erprobt.
- **Die Klassen- und Wiederholungsfrage bei `transient`.**
  [`slice-capture-transient-wiederholung`](slice-capture-transient-wiederholung.md);
  die Schwellen-Kette ist eine andere Ursache.
- **Der Start-Pfad.** [`slice-start-vorlauf-grenze`](slice-start-vorlauf-grenze.md)
  ändert `Stream.Run`/`START_REPLICATION` und den Vorlauf; dieser Slice ändert die
  Rückgabe-Priorität nach der Rückkehr von `stream.Run`. Beide berühren `wiring.go` und
  den Runner an entgegengesetzten Enden (Verdikt §3, Reihenfolge).

## 2. Definition of Done

Jedes Kriterium trägt „Zu belegen durch:“; jede Aussage über eine Mutation ist eine
**Erwartung**, bis der Implementer sie gefahren hat (Stelle, Instanz, gesehene Farbe;
[`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).

- [ ] **Die Regel steht im Code.** `mergeStreamAndWALFaultOutcome` folgt der
      Vertragstabelle des Verdikts (§2): ist ein WAL-Schwellen-Fehler gesetzt, ist der
      Ausgang von `Run` bei Stream-Ausgang `nil` der WAL-Fehler; bei einem
      Ordnungs-Sentinel (`mapper.ErrChangeWithoutBegin`, `…CommitWithoutBegin`,
      `…BeginWithoutCommit`) der Sentinel; bei einem Fehler, dessen Kette
      `context.Canceled` trägt, der WAL-Fehler (neu); bei jedem anderen Stream-Fehler der
      Stream-Fehler. Ohne gesetzten WAL-Fehler bleibt der Stream-Ausgang unverändert; ein
      echter Persistenzfehler ohne Abbruch-Folge bleibt `storage`. Tabellentest über die
      vier Zeilen, dazu ein Fall mit **echt gewrappter** Kette
      (`fmt.Errorf("%w: %w", outbound.ErrStorage, <Fehler mit context.Canceled>)`) und ein
      Fall „`storage` ohne Abbruch-Folge neben gesetztem WAL-Fehler bleibt `storage`“.
      *Zu belegen durch:* `make test` (Race-Detector) grün; die Mutation „die alte
      Priorität (`streamErr != nil` zuerst) zurück“ färbt den Fall „Abbruch-Folge →
      WAL-Fehler“ rot (*erwartet, zu erproben*; Stelle, Instanz und Farbe trägt der
      Implementer nach). Trägt die Prüfung `errors.Is(err, context.Canceled)` die Kette
      an der echten Stelle nicht, ist `streamCtx.Err() != nil` bei lebendem
      Prozess-Kontext der zulässige Ersatz derselben Regel (Verdikt §2, kein weiterer
      Architect-Zug).
- [ ] **Die Runner-Phase trägt die Klasse als Zusage.** Die Phase „Fehlerschwelle beendet
      den Container“ in `make test-integration` prüft `cdc.process_heartbeat.error_class`
      auf `replication` (statt sie nur auszugeben); die Abdeckungs-Zeile in
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) (Erzeugnis des
      Runners) nennt die Klasse; die Ausgabezeile der Phase bleibt. *Zu belegen durch:*
      ein realer, grüner `make test-integration`-Lauf mit der Ausgabezeile (Klasse
      `replication`) und der Abdeckungs-Zeile; die Mutation „die Regel im Code
      zurücknehmen“ (Image neu gebaut) färbt die Phase rot — *erwartet, zu erproben*:
      eine Stelle, ein Lauf. Der Aufbau ist der der Phase (gehaltene Persistierung), also
      der Fall, der zuvor `storage` war.
- [ ] **Kommentare und die benannte Grenze sind nachgezogen.** Der Kommentar an
      `mergeStreamAndWALFaultOutcome` (`wiring.go`: „jede Klasse“, „nur zum Zug, wenn der
      Stream-Lauf regulär endete“), der Kommentar von
      `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` und der Kommentar in
      `walretention_slotgrowth_internal_test.go` (die Kette „Fehler bei Stream-Ende“ trägt
      das Prozessende) tragen die Regel im Indikativ, je höchstens eine Kennung, ohne
      Konjunktiv über die frühere Fassung ([`AGENTS.md`](../../../../AGENTS.md) §3.7); die
      benannte Grenze zur Klasse in [`harness/README.md`](../../../../harness/README.md)
      §Sensors bei `make test-integration` entfällt, ihr Träger dort war dieser Slice.
      *Zu belegen durch:* `make kommentar-kennungen DIFF=<Basis>` ohne Kandidat in den
      geänderten Blöcken (Form, nicht Wahrheit) und Review; der Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); `make suchlauf-nachmessen
      PLAN=docs/plan/planning/open/slice-wal-fehlerschwelle-ausgangsklasse.md` läuft mit
      den `diff`-Zeilen des Implementers durch.
- [ ] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors (Zeile
      `make test-integration`: die Grenze entfällt, die Phase nennt die Klasse als
      Zusage); das Benutzerhandbuch bleibt unberührt (§1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert (§8 nennt die Einträge, die dieser Slice trägt).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten: ihre Closure liegt nach
      `slice-transformationen-betriebsdoku`, der nach diesem Slice startet; ein Slice
      ohne Welle wird von ihr mitgeprüft).

**Umfang:** S — Schätzung, nicht gemessen: eine Regel in einer kleinen Funktion, ein
Tabellentest, eine Assertion im Runner, drei Kommentare und ein Satz in
`harness/README.md` (Verdikt §3); ein Runner-Lauf trägt die Zeit von `make
test-integration` (Risiko §6).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (`mergeStreamAndWALFaultOutcome`, Kommentar; `walRetentionFault`-Kommentar liest der Implementer) | update | die Regel des Verdikts §2; der Kommentar trägt die Regel im Indikativ |
| `internal/bootstrap/walretention_internal_test.go` (`TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`, `…FallsBackToFaultOnRegularStreamEnd`) | update | Tabellentest über die vier Zeilen, gewrappter Fall, `storage`-ohne-Abbruch-Fall (`LH-QA-REL-003`, `ADR-0049`); Kommentar des ersten Tests nachgezogen |
| `internal/bootstrap/walretention_slotgrowth_internal_test.go` (Kommentar, Zeile 27) | update | die Kette „Fehler bei Stream-Ende“ trägt das Prozessende — Kommentar nachgezogen |
| `tools/harness/run-integration-tests.sh` (Phase „Fehlerschwelle beendet den Container“: `abdeckung_declare`-Text, Prüfung von `wal_stop_class`, Kommentar davor) | update | die Klasse `replication` als Zusage der Phase (Eingabeseiten-Mutation: Regel zurücknehmen); der Kommentar davor trägt drei Übergaben aus `slice-capture-leerlauf-quellbelege` (Review F-4, Verifikation V-5): die Allaussage „der Rückstand erreicht die Fehlerschwelle nur, wenn der Slot nichts bestätigt“ ohne Anker (die Phase belegt einen Fall), die Grenze der Klasse ohne Rang-Zeiger (`ADR-0049`, dieser Slice) und die Kopplung an die Phase „Leerlauf-Bestätigung“ davor (`WAL_*`, `bf_wal_hold`, `wal_feed_started`), im Kommentar nur für die Konfigurationsdatei genannt |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner; die Ort-Zeilen verschieben sich |
| `harness/README.md` §Sensors (Zeile `make test-integration`) | update | die benannte Grenze entfällt; ändern nur der genau benannte Satzteil (die Zeilen sind sehr lang) |

**Ansatz (Liste):**

- Die Regel ist eine Funktion mit vier Fällen (Verdikt §2), kein Umbau: `Run` ruft sie
  unverändert an einer Stelle (`wiring.go`, Rückgabe nach dem Herunterfahren der
  Goroutinen).
- Der Test des gewrappten Falls trägt die Form der echten Kette (`ErrStorage` und ein
  Fehler mit `context.Canceled`), nicht einen einfachen `context.Canceled`; der Fall
  „`storage` ohne Abbruch-Folge“ hält die Grenze nach der anderen Seite.
- Die Runner-Phase bleibt der Beleg am komponierten Prozess (`ADR-0030`): die Klasse
  erschien erst dort (`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`, Verdikt §6).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Klasse des Ausgangs bei
Erreichen der Fehlerschwelle“). Suchraum: der ganze Baum ohne die drei Ausnahmen von
[`AGENTS.md`](../../../../AGENTS.md) §3.13 (`docs/reviews/**`, Records unter `done/`,
`.harness/baseline/**`), keine weitere Einschränkung; die Plan-Datei schließt das Werkzeug
aus. Stand ist `7305b578` (der Commit vor der Anlage dieses Slice, vom Planner am
2026-09-27 gemessen; ein Stand ist eine Commit-Kennung, nie `HEAD`); der Implementer misst
am Parent seiner Arbeit neu und trägt die `diff`-Zeilen ein (`make suchlauf-nachmessen
PLAN=docs/plan/planning/open/slice-wal-fehlerschwelle-ausgangsklasse.md`). Die Zahlen des
Standes `7305b578` sind mit `git grep -n` gemessen, nicht übernommen.**

Symbolnamen der bewegten Stelle:

```suchlauf
7305b578 58 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Beschreibung und Zählwort der Klassen-Aussage samt Hedge (die Klasse des Ausgangs, die
Priorität, das reguläre Stream-Ende):

```suchlauf
7305b578 24 -n -E 'Ausgangs-Klasse|Klasse des Ausgangs|mit der Klasse .replication.|beendet sich der Feed-Container|nur zum Zug|Fehler bei Stream-Ende|regulär endete|jede Klasse|Stream-Fehler jeder Klasse' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Die Klassen-Zusagen des Benutzerhandbuchs (der Träger, der die Klasse den Betreibern
nennt):

```suchlauf
7305b578 6 -n -E 'Feed-Container mit der Klasse .replication.|klassifiziert den Lauf als .replication.|Transport-/Verbindungsstörung' -- docs/user/benutzerhandbuch.md
```

| Träger | Befund (Stand `7305b578`, vom Planner gelesen) | Behandlung |
|---|---|---|
| Kommentar an `mergeStreamAndWALFaultOutcome` (`wiring.go`, Zeilen 1168–1177) | trägt „jede Klasse“ und „nur zum Zug, wenn der Stream-Lauf regulär endete“; beide gelten nach der Regel nicht mehr | Implementer zieht nach (DoD 3) |
| Kommentar von `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` und `…FallsBackToFaultOnRegularStreamEnd` (`walretention_internal_test.go`) | der erste beschreibt „eine Stream-Ordnungs-Verletzung erreicht `Run`s Rückgabewert“ als einzige Priorität-Zusage, der zweite „erst ein regulärer Stream-Abschluss … lässt den WAL-Fehler durch“; beides wird mit der Regel enger/falsch | Implementer zieht nach (DoD 3) |
| Kommentar in `walretention_slotgrowth_internal_test.go` (Zeile 27) und die Assertion in Zeile 141 („bei regulärem Stream-Ende“) | nennt „Fehler bei Stream-Ende“ als Träger des Prozessendes; die Assertion prüft `merge(nil, fault)` und bleibt wahr | Kommentar nachziehen; Assertion lesen |
| `harness/README.md` §Sensors, Zeile `make test-integration` | trägt den Satz „**Benannte Grenze:** … der Befund liegt beim Architect“ | entfällt mit DoD 2 und 3; bis dahin nennt der Satz als Träger dieses Slice (Nachzug des Planners, Stand dieser Anlage) |
| Beschreibung der Phase im Runner (`abdeckung_declare`, Kommentar, Ausgabezeile) und `docs/user/e2e-abdeckung.md` Zeile 68 | die Ausgabezeile nennt die Klasse (`$wal_stop_class`), die Zusage nicht; die Abdeckungs-Zeile nennt sie nicht | Runner-Phase und Erzeugnis ziehen mit (DoD 2) |
| `docs/user/benutzerhandbuch.md` (Zeilen 466, 866, 871, 1526, 1614, 1833) | vom Planner gelesen: vier Stellen (466, 866–868, 1526, 1614–1615) sagen `replication` für die Fehlerschwelle bzw. die Unterart Transport-/Verbindungsstörung zu; 871 und 1833 beschreiben die Reichweite der Schwellen und die Versionshistorie. **Alle bleiben mit der Korrektur wahr; kein Träger nennt `storage` für diese Kette** (Muster `Klasse .storage.` im Handbuch: zwei Treffer, Zeilen 578/581, betreffen den Backfill-Run im Snapshot-Fenster, einen anderen Gegenstand) | keine Änderung; der Implementer misst neu |
| `spec/pflichtenheft.md` (Klassen-Tabelle, Zeile `replication`) und `ADR-0049` | die Norm, gelesen (Verdikt L1 bis L3); ihr Wortlaut ist der Gegenstand der Korrektur, kein nachzuziehender Träger | unberührt |
| **Fremde Träger:** `slice-capture-leerlauf-quellbelege` §2 und §3 (DoD 2, Meldung der Träger); das Register `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` und `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (`state.md`) | der Plan dort nennt dieses Slice als Adresse (Nachzug des Planners); die Registereinträge tragen den Fall (Verdikt §6) | Register: **gemeldet, nicht mitgeändert**; Frist: die Closure dieses Slice, der Planner der Closure setzt die Ausgänge |
| Beschreibung in `docs/plan/planning/` (Plan von `slice-capture-leerlauf-quellbelege`, Welle-Plan, Roadmap) | die Aussage zur Klasse `storage` bei gehaltener Persistierung steht im Plan des Quellbeleg-Slice (§3 „Befund der Erprobung“) und beschreibt die Messung jenes Slice | bleibt dort stehen; die Träger-Meldung und DoD 2 dort nennen diesen Slice als Adresse |

## 4. Trigger

**Start** (`next` → `in-progress`): eine **Vorab**-Bedingung, kein Nachweis nach der
Umsetzung (`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen, 2×):
`slice-capture-leerlauf-quellbelege` liegt in `done/` (beide Slices ändern die
Runner-Phase „Fehlerschwelle beendet den Container“ in
`tools/harness/run-integration-tests.sh`; die Phase ist die Falsifikation dieses Slice;
WIP-Limit 1), `slice-leerlauf-phase-last-in-stuecken`
liegt in `done/` (beide Slices ändern denselben Runner an verschiedenen Stellen — jener die Last
der Phase „Leerlauf-Bestätigung“, dieser die Phase „Fehlerschwelle beendet den Container“ —, und ein
Rot in der Leerlauf-Phase lässt diese Phase ungelaufen;
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§3), und kein weiterer Slice liegt in `in-progress/`. Der Slice muss `done`
sein, **bevor** [`slice-start-vorlauf-grenze`](slice-start-vorlauf-grenze.md) startet
(Start-Trigger dort): beide ändern `wiring.go` und den Runner; der kleinere Slice
zuerst verstellt die Prüfspur des größeren nicht (Verdikt §3). Der Übergangs-Commit
`next` → `in-progress` nennt
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md).

**Kante zu `slice-leerlauf-phase-last-in-stuecken` (beauftragt).** Die Phase
„Leerlauf-Bestätigung“ vor der Runner-Phase dieses Slice war im ersten
`e2e.yml`-Lauf nach dem Push von `slice-capture-leerlauf-quellbelege` einmal rot (Lauf
36287009221, Leg PostgreSQL 18); der Architect hat die Ursache im Testaufbau bestimmt (der
Stoß einer einzelnen Anweisung gegen den Prüf-Takt) und den Slice beauftragt
([`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§3; Register `BEO-PGC/test-integration-retention-timing-flake`). Er geht diesem Slice
**voraus**: ein Rot dort lässt die Phasen dahinter ungelaufen, auch die Mutation dieses Slice
an der Phase „Fehlerschwelle beendet den Container“.

**Verifier-Hinweis** (Verdikt §4): ein Rot mit der Signatur „Fehlerklasse `replication` …
WAL-Rückstand … über Fehlerschwelle“ in der Phase „Leerlauf-Bestätigung“ ist, solange
`slice-leerlauf-phase-last-in-stuecken` nicht in `done/` liegt, weder Beleg noch
Widerlegung dieses Slice — der Verifier wiederholt den Lauf (`gh run rerun <Lauf> --failed`)
und nennt Lauf, Versuchsnummer und Job-Kennungen beider Versuche; nach diesem Slice ist
dasselbe Rot ein Befund und ein Architect-Zug. Ein Rot anderer Signatur oder in einer anderen
Phase gilt nicht als dieser Fall.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht zu erwarten (Umfang S);
  sollte die Runner-Phase einen zweiten, unabhängigen Befund zeigen, ist dieser der
  abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls die Regel die Klasse in der Runner-Phase
  nicht erreicht, obwohl der Unit-Test grün ist — die Kette trägt `context.Canceled`
  dann an der echten Stelle nicht, und der zulässige Ersatz (`streamCtx.Err() != nil`)
  trägt ebenfalls nicht: eine Architect-Frage, kein stiller Umbau.

## 5. Closure-Trigger

DoD vollständig (die Regel im Code, die Klasse `replication` als Zusage der Runner-Phase
am komponierten Prozess belegt, Kommentare und die benannte Grenze nachgezogen) + `make
gates` grün + ein realer, grüner `make test-integration`-Lauf + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Fall, den die Regel abbildet, bleibt unerreichbar: `errors.Is(err,
  context.Canceled)` trägt an der echten Stelle nicht.** *Hergeleitet, nicht erprobt*
  (Verdikt §2: `Classify` wrappt mit `%w: %w`, die Meldung endet auf „context canceled“;
  nicht gelesen). *Erwartet, zu belegen durch:* der Tabellentest mit gewrappter Kette
  **und** die Runner-Phase mit der Klasse `replication`; der zulässige Ersatz ist
  `streamCtx.Err() != nil` bei lebendem Prozess-Kontext. **Ausgang:** *(bei Closure)*
- **Die Regel verdeckt einen echten Persistenzfehler.** Ein Fehler der Klasse `storage`
  ohne Abbruch-Folge neben einem gesetzten WAL-Fehler würde nach der Regel weiter den
  Stream-Fehler zurückgeben; ein Fehler, dessen Kette `context.Canceled` trägt **und**
  der eine echte Ursache hat, wird zum WAL-Fehler — die Ursache steht im Log neben der
  Abbruch-Zeile (Verdikt §2). *Erwartet, zu belegen durch:* der Fall „`storage` ohne
  Abbruch-Folge bleibt `storage`“ im Tabellentest; ohne gesetzten WAL-Fehler ändert sich
  nichts. **Ausgang:** *(bei Closure)*
- **Die Mutation färbt nichts rot** (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  verkörpert): der Test der Regel und die Runner-Phase binden sich an die Eingabeseite
  (der Stream-Fehler der Klasse `storage` neben dem gesetzten WAL-Fehler, die Klasse aus
  `cdc.process_heartbeat`), nicht nur an die Ausgabe. *Erwartet, zu belegen durch:* je eine
  Mutation, deren Farbe der Implementer **gesehen** und im Bericht genannt hat (Stelle,
  Instanz, Farbe). **Ausgang:** *(bei Closure)*
- **Der Beleg ist zeitabhängig und verlängert `make test-integration` nicht**: die
  Phase besteht bereits; ihre Laufzeit ist unverändert bis auf die Klassen-Prüfung
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert). *Erwartet, zu belegen
  durch:* die gedruckte Laufzeit des Laufs ([`AGENTS.md`](../../../../AGENTS.md) §3.12
  Instanz A). **Ausgang:** *(bei Closure)*
- **Die neue Assertion fällt still aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×). *Erwartet, zu belegen durch:* die Zeile in `docs/user/e2e-abdeckung.md`
  nennt die Klasse, und die Phase färbt sich rot, wenn die Regel zurückgenommen ist.
  **Ausgang:** *(bei Closure)*
- **Ein Kommentar sagt die Regel breiter oder schmaler, als der Code sie trägt**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, verkörpert, 6×). *Erwartet,
  zu belegen durch:* der Reviewer fährt die zugesagten Pfade der drei Kommentare im Code
  nach ([`AGENTS.md`](../../../../AGENTS.md) §3.7). **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle
  ([welle-transformationen](../welle-transformationen.md), offen).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Composition Root und Test-Runner sind keine eigenen Sub-Areas —
kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-27 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (5×, Ausgang *geplant*:
`slice-start-vorlauf-grenze`; der fünfte Beleg ist dieser Fall, Verdikt §6 — der Träger
der Klasse ist die Phase in `make test-integration` im Slice, der die Eigenschaft
einführt, hier dieser Slice; kein weiterer Eintrag, der Ausgang bleibt),
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (6×, verkörpert, Risiko §6 und
DoD 3: die Kommentare an `mergeStreamAndWALFaultOutcome` und in Zeile 27 des Tests),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (19×, verkörpert, Risiko §6),
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` (8×, verkörpert; die Aussage über den
Stand des Streams in Verdikt §5 ist als hergeleitet gekennzeichnet, Kernaussage §1),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, Risiko §6),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×, Risiko §6),
`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×, Start-Trigger §4),
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 5×: Repo-Dateien
entstehen über Edit/Write, nie über Umleitung, [`AGENTS.md`](../../../../AGENTS.md) §3.1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
