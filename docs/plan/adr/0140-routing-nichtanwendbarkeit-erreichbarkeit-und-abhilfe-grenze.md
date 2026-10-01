# ADR-0140: Routing — Erreichbarkeit der Nichtanwendbarkeit und Grenze der Abhilfe-Zusage (V3; Schärft ADR-0137, ADR-0138, ADR-0139)

**Status:** Accepted

**Datum:** 2026-10-01

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch die
Architect-Frage zu F-3 im Review `review-slice-routing-kern-label`, Vorab-Bedingung
V3 in `welle-routing`)

**Bezug:** [`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Negative: sichtbarer
Fehlerzustand, kein stilles Standardziel),
[`LH-FA-ADM-003`](../../../spec/lastenheft.md) (Fehlerzustand erkennbar),
[`LH-FA-SCH-003`](../../../spec/lastenheft.md) (entfernte Spalten),
[`LH-FA-SCH-004`](../../../spec/lastenheft.md) (inkompatible Änderungen),
[ADR-0137](0137-routing-zustellziele-persistiertes-ziel-label.md) (Teilfrage 4 —
**nicht geändert**, nur geschärft),
[ADR-0138](0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
(Abschnitt „Offen — V3“, Festlegung 2 — **nicht geändert**),
[ADR-0139](0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[ADR-0112](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Teilfrage 4, Folgepflicht 5 — Gegenprobe),
[ADR-0117](0117-backfill-run-fehlerklasse-schema.md),
[ADR-0063](0063-lh-fa-sch-003-testform-korrektur.md)

**Schärft:** [`LH-FA-CFG-008.a`](../../../spec/pflichtenheft.md) (Absatz „Abhilfe
(Zusage)“: der Satz „Offen und nicht gemessen …“),
[`SPEC-032`](../../../spec/pflichtenheft.md) (Anwendbarkeit),
[`SPEC-008`](../../../spec/pflichtenheft.md) (Zeile `schema`, Absatz „Nicht
anwendbare Regel“)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0137` Teilfrage 4 sagt: Eine Routing-Regel ist nicht anwendbar, wenn
`when.column` in der Relation der Change fehlt; der Erfassungspfad endet mit
`ErrRoutingNotApplicable` (Klasse `schema`), Abhilfe ist `cdc.remove_route` und
Neustart. `ADR-0138` ließ offen, ob dieser Fall am laufenden System
erreichbar ist (V3), und verlangte bei einem unerreichbaren Fall ein
Architect-Verdikt vor dem Negative-Beleg.

### Befunde

- **Erprobt (Unit-Ebene, am Quelltext gelesen und als Test vorhanden).**
  `TestRoutingNotApplicableReachabilityThroughRelationCheck`
  (`internal/adapters/driving/replication/mapper/routing_test.go`, Instanz
  `Assembler` mit `fakeSchemaStore`; der Implementer hat ihn grün gefahren, der
  Reviewer nachgefahren, ich fahre ihn nicht): (a) bei bekannter Spaltenform der
  aktuellen Version meldet `Consume` eine Relation ohne die Bedingungsspalte als
  `ErrIncompatibleSchemaChange`, bevor eine Change assembliert wird; (b) trägt die
  aktuelle Version keine Spaltenform, registriert `observeRelation` die erste
  Relation ohne Vergleich, und die folgende Change endet als
  `ErrRoutingNotApplicable`.
- **Gelesen, `mapper.go`.** `observeRelation` vergleicht die Relation mit der
  gespeicherten Spaltenform der aktuellen Version (`classifyRelationColumns`);
  der Vergleich liest **keine Regel**. Eine Relation, die eine bekannte Spalte
  nicht mehr trägt, ist `relationOther`, mit oder ohne Routing- oder
  Transformationsregel.
- **Hergeleitet (nicht gefahren).** Der Server sendet beim Start einer
  Replikationssitzung die Relation-Nachricht vor der ersten Change der Tabelle
  erneut. Nach `cdc.remove_route` und Neustart träfe sie dieselbe gespeicherte
  Spaltenform, und der Prozess endete im Fall (a) an `ErrIncompatibleSchemaChange`
  erneut. Die Regel ist dafür nicht ursächlich; ihr Entfernen ändert nichts am
  Vergleich.
- **Gelesen, `pflichtenheft.md` (R3).** Ein `set_route` mit einer `when.column`,
  die an der Quelle nicht existiert, endet `failed` (`Spalte existiert nicht an
  der Quelle`). Eine Regel auf eine fehlende Spalte entsteht also nicht beim
  Antrag, sondern nur durch eine spätere Änderung an der Quelle oder an der
  Relation-Sicht.
- **Gegenprobe Bestand, `ADR-0112` Teilfrage 4 (gelesen).** Dort steht derselbe
  Befund: die spalten-entfernenden Fälle „enden dort schon vorher mit derselben
  Klasse“; erreichbar durch die Prüfung der Regeln ist die K3-Kollision nach
  kompatibler Erweiterung. `TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass`
  nutzt diesen Fall (gelesen, `test/integration/integration_test.go`). V3 war für
  Transformationen damit durch einen **erreichbaren Ersatzfall** gelöst, nicht
  durch die Abhilfe im Fall „Spalte entfernt“. Routing hat kein K3-Pendant
  (kein Zielname, der mit einer Spalte kollidieren kann).

### Konstraints

- Klassenmenge geschlossen (`ADR-0023`): keine achte Klasse.
- Dass ein Erfassungspfad an einer entfernten Spalte endet, ist `LH-FA-SCH-003`
  und `ADR-0063`; diese ADR ändert dieses Verhalten nicht.

## Entscheidung

Wir wählen **Option C — die Abhilfe-Zusage gilt genau für den Fall, in dem der
Erfassungspfad an `ErrRoutingNotApplicable` endet; `ErrRoutingNotApplicable`
bleibt als Sicherheitsnetz; der Weg „Spalte entfernt“ bleibt der Pfad der
inkompatiblen Schemaänderung, ohne dass eine Routing-Regel ihn verändert.**

1. **Option (1) wird nicht gewählt.** Die Relation-Prüfung berücksichtigt keine
   Regeln. Eine regelabhängige Meldung (`ErrRoutingNotApplicable` statt
   `ErrIncompatibleSchemaChange`) änderte die Klasse nicht und machte die
   Abhilfe nicht wahr: die gespeicherte Spaltenform bliebe inkompatibel. Eine
   Abhilfe, die den Pfad weiterlaufen lässt, verlangte, dass eine entfernte
   Spalte, die keine Regel nennt, den Erfassungspfad **nicht** beendet — das ist
   eine Änderung von `LH-FA-SCH-003`/`-004` und `ADR-0063`, nicht des Routings,
   und gehört in eine eigene Entscheidung, falls der Auftraggeber sie will.
2. **Geltung der Abhilfe.** „`cdc.remove_route` beantragen, Prozess neu starten“
   ist die Abhilfe für einen Erfassungspfad, der an `ErrRoutingNotApplicable`
   endete (`ADR-0137` Teilfrage 4, unverändert). Sie gilt **nicht** für einen
   Erfassungspfad, der an `ErrIncompatibleSchemaChange` endete; dort reicht das
   Entfernen der Regel nicht, und die Abhilfe ist die dieser Ursache
   (`LH-FA-SCH-003`/`-004`), die diese ADR nicht neu definiert. Der Betreiber
   entfernt eine Regel auf die entfernte Spalte zusätzlich, sonst endet der
   Erfassungspfad nach der Abhilfe der Schemaänderung an der dann fehlenden
   Bedingungsspalte (*hergeleitet*: gespeicherte Spaltenform ohne die Spalte,
   Regel noch geführt).
3. **`ErrRoutingNotApplicable` bleibt.** Das Sicherheitsnetz steht auf der
   Change-Prüfung (`checkRoutes`), vor der Sequenz und vor jeder Serialisierung
   (Review: M5 rot). Es bleibt erreichbar in diesen Fällen, geordnet nach der
   Belegtiefe:
   - (b) **Erstaktivierung ohne Spaltenform** — Unit gemessen (siehe Befund),
     am System **ungemessen**;
   - (c) **Publication mit Spaltenliste**, die die Bedingungsspalte nicht trägt:
     die Relation trägt sie dauerhaft nicht, die erste Relation legt diese
     Spaltenform an, der Vergleich meldet nichts, die Change endet als
     `ErrRoutingNotApplicable` — *hergeleitet*, nicht gefahren; die Abhilfe
     (Regel entfernen, Neustart) greift dort, weil die Spaltenform gleich bleibt
     — *hergeleitet*;
   - (d) **Backfill-Run** (`ADR-0138` Festlegung 2): `TableSnapshot.Columns()`
     liest die aktuelle Tabelle; eine entfernte Spalte mit weiter geführter Regel
     ergibt dort `failed`/`schema`, run-lokal, ohne dass der Erfassungspfad
     enden muss, solange keine Change der Tabelle eine neue Relation auslöst —
     *hergeleitet*; die Abhilfe ist ein neuer Antrag nach `cdc.remove_route`
     (`ADR-0117`, Handbuch-Muster der Transformationen).
4. **Beleg des Negative-Kriteriums von `LH-FA-CFG-008`.** Das Kriterium
   verlangt einen sichtbaren Fehlerzustand und kein stilles Standardziel oder
   Verwerfen; es verlangt keine bestimmte Sentinel-Fehlerart. Es wird getragen
   durch, in dieser Reihenfolge der Nähe zum System:
   - **E2E, Weg (a) mit aktiver Regel** (*Erwartung*, zu belegen in
     `slice-routing-e2e`): eine eigene Tabelle trägt eine Routing-Regel auf die
     Spalte `region`; `ALTER TABLE … DROP COLUMN region`; die nächste Change der
     Tabelle beendet den Erfassungspfad sichtbar mit der Klasse `schema`
     (`diagnose`, Heartbeat), und über `cdc.changes` erscheint für diese
     Change **keine** Zeile — weder mit `route_target = NULL` noch mit einem
     anderen Ziel. Dieser Beleg trägt *nur* das Negative-Kriterium, nicht die
     Abhilfe-Zusage (Entscheidung 2).
   - **Unit** für `ErrRoutingNotApplicable`: `TestConsumeRoutingRuleOnMissingColumnIsNotApplicable`
     (Prüfung vor Sequenz; M5 rot, am Reviewer-Lauf gelesen) und
     `TestRoutingNotApplicableReachabilityThroughRelationCheck` (Weg b),
     `classifyRunError` → `schema` (gelesen, `wiring.go`).
   - **Die Abhilfe-Zusage von `LH-FA-CFG-008.a`** („Regel entfernen, Neustart,
     zuvor nicht bestätigte Transaktion erscheint ohne zweiten `schema`-Fehler“)
     wird **nur** an einem Fall belegt, der `ErrRoutingNotApplicable` am System
     wirklich erzeugt. `slice-routing-e2e` misst als ersten Schritt, ob (b) oder
     (c) erzeugbar ist, und belegt die Abhilfe dort. Ist keiner erzeugbar, bleibt
     die Abhilfe-Zusage für den Erfassungspfad **Erwartung ohne Systembeleg**,
     steht im Bericht der Welle mit diesem Namen und trägt Unit (Persist-before-ACK
     und Sentinel) plus (d) am Backfill-Run — die Verengung wird benannt, nicht
     stillschweigend vollzogen.
5. **Konsistenz mit den Transformationen.** Gleiches Muster wie `ADR-0112`
   Teilfrage 4: der Weg „Spalte entfernt“ endet vorher in derselben Klasse, der
   Belegfall der Abhilfe ist ein Ersatzfall, der die Regel-Prüfung wirklich
   erreicht. Für Routing gibt es statt der K3-Kollision (b) und (c).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — (1): Relation-Prüfung berücksichtigt Regeln, Abhilfe auch für „Spalte entfernt“ | Zusage der Spec bliebe wörtlich wahr | die gespeicherte Spaltenform bleibt inkompatibel, das Entfernen der Regel läse die Relation nicht anders; wirksam nur, wenn ein Entfernen einer **nicht** von Regeln genannten Spalte den Pfad nicht mehr beendet — Eingriff in `LH-FA-SCH-003`/`-004`, `ADR-0063`, `ADR-0059` Teilfrage 4, alle `Accepted`; Produkt- statt Routing-Frage |
| B — (2): Abhilfe nur für die Erstaktivierung (b), Negative-Beleg auf (b)/Unit zugeschnitten | kleinster Wortlaut | verengt die Zusage auf einen Fall, der am System ungemessen ist; (c) und (d) bleiben unbenannt; der Betreiber im Weg (a), dem naheliegenden Weg, bekäme keine Aussage |
| **C — (3), geschärft (gewählt): Zusage = Fälle von `ErrRoutingNotApplicable`; Weg (a) benannt; Sicherheitsnetz bleibt; Belege getrennt** | ändert keinen `Accepted`-Text und kein Verhalten; sagt dem Betreiber, was im Weg (a) gilt; trennt Negative-Beleg (Weg a erreichbar) von Abhilfe-Beleg (nur an erzeugbarem Fall); gleiches Muster wie `ADR-0112` | die Abhilfe der Zusage ist am Erfassungspfad möglicherweise ohne Systembeleg; (c) und (d) sind hergeleitet |
| D — `ErrRoutingNotApplicable` entfernen, Fall (b) als „nie nötig“ behandeln | weniger Code | (b) und (c) wären dann ein stilles Standardziel (`route_target = NULL` bei fehlender Spalte) — `LH-FA-CFG-008` Negative widersprochen; `ADR-0138` Option E |

## Konsequenzen

- Positiv: der Betreiber im naheliegenden Weg (a) bekommt eine zutreffende
  Aussage; der Negative-Beleg hängt nicht an einem Fall, den das System
  vielleicht nicht erzeugt; kein Eingriff in `LH-FA-SCH-003`/`-004`.
- Negativ: im Weg (a) gibt es keine Routing-spezifische Abhilfe; die Abhilfe der
  Schemaänderung ist regelunabhängig und in der Spec nicht als Schrittfolge
  geführt (`ADR-0063` hält den Erfassungspfad dort für „dauerhaft beendet“).
  Wer mehr will (Fortsetzen nach entfernter, nicht von Regeln genannter Spalte),
  braucht die Folgeentscheidung aus Alternative A.
- Negativ: (c) und (d) sind hergeleitet; wird keiner gefahren, bleibt die
  Abhilfe am Erfassungspfad Erwartung.
- Beobachtung am Bestand (nicht Teil der Entscheidung): der Handbuch-Absatz
  „Fehler: Erfassung endet mit der Fehlerklasse `schema`“ (`docs/user/benutzerhandbuch.md`,
  Transformationen) nennt als Ursache „ihre `column` fehlt in der Relation“ und
  als Lösung „Regel entfernen oder ersetzen, Prozess neu starten“; für eine
  entfernte Spalte endet der Pfad nach `ADR-0112` Teilfrage 4 schon an
  `ErrIncompatibleSchemaChange`, wo das Entfernen der Regel nicht genügt. Gleicher
  Befund wie hier, *hergeleitet*; Träger beim Planner.

### Folgepflichten

1. **Pflichtenheft (`spec`-Slice, Planner).** `LH-FA-CFG-008.a`, Absatz „Abhilfe
   (Zusage)“: den Satz „Offen und nicht gemessen: …“ durch die geschärfte
   Aussage ersetzen (Abhilfe gilt für den Fall `ErrRoutingNotApplicable`; bei
   entfernter Spalte mit bekannter Spaltenform endet der Pfad zuerst an der
   inkompatiblen Schemaänderung, Entfernen der Regel genügt dort nicht; die
   Erreichbarkeit von (b)/(c) am System steht als Messung von `slice-routing-e2e`).
   `SPEC-032` (Anwendbarkeit): ein Satz, dass die Relation-Prüfung der Change-Prüfung
   vorausgeht und regelunabhängig ist. `SPEC-008`, Absatz „Nicht anwendbare Regel
   (Klasse `schema`)“: „Abhilfe nur bei nicht anwendbarer Regel“ ist dort schon
   so formuliert; ein Querverweis auf die Grenze genügt.
2. **`slice-routing-e2e`** (`docs/plan/planning/open/slice-routing-e2e.md`): das
   Kriterium „Negative (B)“ zerlegen in (i) Weg (a) mit aktiver Regel
   (Negative-Beleg, Entscheidung 4), (ii) Messung von (b)/(c) als erster Schritt,
   (iii) Abhilfe-Beleg am erzeugbaren Fall, sonst benannte Verengung im
   Closure-Bericht der Welle. Der Satz „Ist kein Fall erzeugbar: das
   Architect-Verdikt liegt vor“ wird auf dieses Verdikt (`ADR-0140`) verwiesen.
3. **`slice-routing-kern-label`** (`docs/plan/planning/in-progress/slice-routing-kern-label.md`
   §6, V3-Absatz und §7): „Ausgang bei der Closure einzutragen“ mit dem Verweis
   auf `ADR-0140`; kein Code- oder Testwechsel (Test und Sentinel bleiben).
4. **`slice-routing-backfill-pfad`**: Fall (d) erhält einen Unit-/Store-Beleg,
   soweit er dort ohnehin entsteht (Erwartung); kein neues Kriterium, falls der
   Slice die Prüfung der Run-Regeln bereits belegt.
5. **`welle-routing`** (`docs/plan/planning/welle-routing.md` §5): V3 wird
   „beantwortet durch `ADR-0140`“; A-3 bleibt beim Auftraggeber.
6. **Handbuch** (`slice-routing-betriebsdoku`): Fehler/Ursache/Lösung für
   Routing trennt die zwei Ursachen der Klasse `schema`; keine Aussage „Regel
   entfernen“ für die entfernte Spalte.
7. `AGENTS.md` §3.4/§3.5: `spec/architecture.md` bleibt unberührt (kein ADR-,
   Slice- oder Wellenbezug dort); die `Accepted` `ADR-0137`/`-0138`/`-0139`
   werden nicht verändert, ihr offener Punkt wird durch diese ADR geschärft.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (vorhanden) | `TestRoutingNotApplicableReachabilityThroughRelationCheck` trägt (a) und (b) auf Unit-Ebene; erprobt vom Implementer, vom Reviewer nachgefahren, von mir nicht gefahren | `make test` |
| Go-Test (vorhanden) | `TestConsumeRoutingRuleOnMissingColumnIsNotApplicable` trägt die Prüfung vor der Sequenz (Reviewer: Mutation M5 rot, eine Stelle, ein Test); Verallgemeinerung auf andere Stellen *hergeleitet* | `make test` |
| E2E (Erwartung) | Weg (a) mit aktiver Regel: `schema`-Ende, keine Zeile in `cdc.changes`; Messung von (b)/(c); Abhilfe-Beleg am erzeugbaren Fall | `make test-integration` (`slice-routing-e2e`) |

Eine maschinelle Prüfung der Aussage „die Relation-Prüfung liest keine Regel“
gibt es nicht; sie steht als gelesen (`observeRelation`).

## Re-Evaluierungs-Trigger

- Der Auftraggeber will, dass eine entfernte, von keiner Regel genannte Spalte
  den Erfassungspfad nicht beendet: Folge-ADR zu `LH-FA-SCH-003` (Supersedes
  `ADR-0063` teilweise); dann ist Option A neu zu bewerten.
- `slice-routing-e2e` misst, dass weder (b) noch (c) am System entsteht: die
  Abhilfe-Zusage der Spec wird auf Run-Fall (d) und Unit zurückgenommen
  (Folgepflicht 1, Spec-Slice); `ADR-0138` Re-Evaluierungs-Trigger („nie
  erreichbar“) tritt für den Erfassungspfad ein.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Accepted (Kurz-ADR auf Auftrag des Hauptlaufs; Schärft ADR-0137, ADR-0138, ADR-0139) | `welle-routing` V3, `review-slice-routing-kern-label` F-3 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen entstehen als neue ADR mit `Supersedes ADR-0140`.
