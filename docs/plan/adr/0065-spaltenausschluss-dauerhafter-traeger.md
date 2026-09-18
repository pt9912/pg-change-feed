# ADR-0065: Spaltenausschluss — dauerhafter Träger des Ausschlussstandes

**Status:** Accepted — Supersedes [`ADR-0059`](0059-spaltenauswahl-mechanismus.md)
(nur deren Dauerhaftigkeits-Aussage: der Pro-Satz in Teilfrage 1 Option D
(„löst die Neustart-Einschränkung strukturell (derselbe Live-Reload-Pfad wie
Tabellen-Aktivierung)") und der Schlusssatz von §Bestätigung — Live-Reload-
Konsistenz („exakt wie eine Tabellen-Aktivierung heute schon"); die übrigen
Festlegungen dieser ADR — Konfigurationsmechanismus (Teilfrage 1),
Granularität (Teilfrage 2), Wirkort im `Assembler` (Teilfrage 3), Konvergenz
mit [`LH-FA-SCH-003`](../../../spec/lastenheft.md) (Teilfrage 4) und
Rückkanal über den Antrags-Status (Teilfrage 5) — bleiben unverändert
bestehen und werden hier nicht wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-067` <!-- d-check:status-provenance -->,
der die Grenze selbst benannte, und als der Reviewer-Lauf, der sie als
Entscheidung weiterreichte; Modul 8 §Konflikt-Pfad, Planner → Architect →
Planner)

**Bezug:** [`LH-FA-CFG-005`](../../../spec/lastenheft.md) (Haupt-Bezug —
Spaltenauswahl und ihr Rückkanal), [`LH-QA-SEC-004`](../../../spec/lastenheft.md)
(Ausschluss sensibler Spalten — ihre Zusage hängt an der Dauerhaftigkeit des
Ausschlussstandes), [`ADR-0059`](0059-spaltenauswahl-mechanismus.md)
(Mechanismus, Granularität, Wirkort und Rückkanal — superseded nur in der
Dauerhaftigkeits-Aussage), [`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue als einziger Schreibpfad administrativer Zustände),
[`ADR-0034`](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt — Grundlage der
neuen Lesefähigkeit), [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md)
(Schemamigrationen — nur berührt, falls der umsetzende Slice ein
Schema-Objekt wählt), [`ADR-0030`](0030-testpyramide.md) (Test-Tier des
E2E-Belegs), Review zu `slice-067` (Anlass: F-1/F-2), <!-- d-check:status-provenance -->
`docs/plan/planning/open/slice-075-dauerhafter-ausschlussstand-wiedereinspielung.md` <!-- d-check:status-provenance -->
(umsetzender Slice), `docs/plan/planning/observations/BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/`
(Registereintrag der Beobachtung)

**Schärft:** [`SPEC-019`](../../../spec/pflichtenheft.md) (die Bedeutung des
`status`-Wertes `applied` für die beiden Spalten-Antragsarten),
[`ARC-004`](../../../spec/architecture.md) (der Spalten-Prüfungs-Port bekommt
eine Lesefähigkeit für den Ausschlussstand)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0059`](0059-spaltenauswahl-mechanismus.md) (Accepted, 2026-09-14)
entscheidet den Spaltenausschluss in fünf Teilfragen: zwei neue Antragsarten
auf der Queue aus [`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md),
Granularität pro Tabelle, Wirkort im `Assembler` bei der
Row-Image-Konstruktion (vor jeder Serialisierung), Konvergenz mit
[`LH-FA-SCH-003`](../../../spec/lastenheft.md) bei realer Spaltenlöschung und
den Antrags-Status als einzigen Rückkanal. Der umsetzende `slice-067` <!-- d-check:status-provenance -->
trägt den Ausschlussstand als `TableBinding.ExcludedColumns`; der Reviewer-Lauf
hat drei Eigenschaften dieses Standes am Code festgestellt, die die
Entscheidung nicht adressiert:

1. **Der Prozessstart kennt keinen Ausschlussstand.** `activatedTableBindings`
   (`internal/bootstrap/wiring.go`) baut die Bindungen aus den
   Aktivierungs-Zeilen (`TableActivationPort.List`) und der aktuellen
   Schema-Version auf; die Spaltennamen der Spalten-Antragsarten liest dieser
   Pfad nicht. Ein Prozess-Neustart erfasst eine zuvor ausgeschlossene Spalte
   wieder.
2. **Ein Bindungs-Zyklus löscht den Stand.** `RemoveBinding` (Deaktivierung)
   entfernt den Karteneintrag samt `ExcludedColumns`; `AddBinding`
   (Aktivierung) legt ihn über die Bindung des Aufrufers an. Der Merge in
   `AddBinding` erhält den Stand nur, solange ein Vorgänger-Eintrag liegt —
   nach einer Deaktivierung liegt keiner. Dieser Auslöser braucht keinen
   Prozess-Neustart.
3. **Der einzige dauerhafte Beleg meldet Erfolg.** Die `applied`-Zeile in
   `cdc.administration_request` trägt den verarbeiteten Antrag;
   `processAdministrationRequests` liest `pending`-Zeilen. Nach (1) oder (2)
   liest der Antrags-Datensatz weiter `applied`, während die Spalte wieder
   erfasst wird.

Damit hängt die Zusage von [`LH-QA-SEC-004`](../../../spec/lastenheft.md)
(„ausgeschlossene Spaltenwerte erscheinen nicht in den Changes", Messmethode
über [`LH-FA-CFG-005`](../../../spec/lastenheft.md)) an der Prozesslebensdauer
— und ihr Verlust ist still: kein Gate, kein Statuswert, kein Log meldet ihn.
Ein bereits persistierter Wert ist nicht zurückholbar; genau den schließt
[`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 3 Option D als
Wirkort aus („der sensible Wert wird **nie** serialisiert und nie an die
Persistenzschicht übergeben").

**Was [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) zusagt und was nicht:**
Ihre Live-Reload-Aussage richtet sich auf das *Wirksamwerden ohne
Prozess-Neustart* — die Contra-Spalte von Option B/C benennt die
„Neustart-Einschränkung" als „wirkt nur beim Prozessstart", und die
Entscheidung sucht den Zug auf einen *laufenden* Prozess. Die Dauerhaftigkeit
des Standes ist damit nicht zugesagt. Die zwei Stellen, die eine Parität zur
Tabellen-Aktivierung nahelegen („derselbe Live-Reload-Pfad wie
Tabellen-Aktivierung", „exakt wie eine Tabellen-Aktivierung heute schon"),
vergleichen den *Pfad* des Wirksamwerdens und treffen mit
`cdc.source_table` auf eine Fähigkeit, deren Zustand dauerhaft liegt — der
Vergleich trägt die Dauerhaftigkeit nicht aus. Diese ADR zieht die
Konsequenz: Sie benennt die beiden Stellen als superseded und entscheidet
die Frage, die dort offen blieb.

**F-1 und F-2 sind eine Frage.** Der stille Erfolg eines `exclude_column`-
Antrags gegen eine Tabelle ohne laufende Bindung (Review-Finding F-2) hat
dieselbe Wurzel: der Stand hängt an einer Bindung, nicht an der Tabelle. Ein
Statuswechsel allein (`failed` statt `applied`) träfe das Symptom: der Antrag
bliebe ohne dauerhaften Träger, und eine Tabelle ohne laufende Bindung hätte
weiterhin keinen Ort für ihren Ausschlussstand.

## Entscheidung

Wir wählen: **Der Ausschlussstand bekommt einen dauerhaften, tabellen-scoped
Träger; die `Assembler`-Bindung ist sein Laufzeit-Cache.** Vier Festlegungen:

1. **Tabellen-scoped statt bindungs-scoped.** Der Ausschlussstand gehört zu
   einer Tabelle der Quelle (qualifizierter Name), nicht zu einer konkreten
   `TableBinding`. Jeder Pfad, der eine Bindung anlegt — der Prozessstart
   (`activatedTableBindings`) und der Aktivierungs-Zweig der
   Antrags-Verarbeitung (`AddBinding`) — trägt den dauerhaften Stand der
   Tabelle mit.
2. **Abgeleitet aus der Antrags-Historie.** Quelle des Standes sind die
   `applied`-Zeilen der beiden Spalten-Antragsarten in
   `cdc.administration_request`, in der Reihenfolge `requested_at`, bei
   gleichem Zeitstempel deterministisch nach `administration_request_id`
   ausgewertet: `exclude_column` trägt den Spaltennamen ein,
   `include_column` nimmt ihn heraus. Kein zweiter Konfigurationsweg, keine
   zweite Zustandstabelle — derselbe Mechanismus wie in
   [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 1 Option D,
   jetzt mit dauerhafter Wirkung. Die konkrete Abfrageform (Sicht, Index,
   Materialisierung) und die Form der neuen Lesefähigkeit am Outbound Port
   sind Gegenstand des umsetzenden Slices und — soweit sie einen öffentlichen
   Vertrag bilden — des Pflichtenhefts.
3. **`applied` heißt: dauerhaft vermerkt.** Der Antrags-Status bleibt der
   Rückkanal ([`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 5).
   Ein Antrag gegen eine existierende Spalte einer Tabelle ohne laufende
   Bindung endet **nicht** `failed`: mit dem dauerhaften Stand wirkt er,
   sobald die Tabelle erfasst wird; sein Zielzustand ist über diesen
   Mechanismus erreichbar, er ist nur aufgeschoben. Der Negative-Pfad für
   eine an der Quelle nicht existierende Spalte bleibt unverändert
   (`ErrSourceColumnMissing`, `failed` samt Fehlertext).
4. **Wirkort und Rückkanal bleiben, wie [`ADR-0059`](0059-spaltenauswahl-mechanismus.md)
   sie entschieden hat.** Die Filterung liegt weiter im `Assembler` bei der
   Row-Image-Konstruktion (Teilfrage 3), der Zustand wird weiter über die
   Antrags-Queue geschrieben (Teilfrage 1), und die beiden Erhalt-Punkte im
   laufenden Prozess (`AddBinding`-Merge, `setSchemaVersion`) bleiben als
   In-Prozess-Schutz bestehen. Diese ADR ändert allein, **woher** der Stand
   beim Anlegen einer Bindung kommt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; der Stand bleibt prozessgebunden, die Grenze wird nur benannt (§6 des Slice, `welle-18` §6, Registereintrag) | kein Code, kein Schema-Objekt, keine Dauerhaftigkeits-Zusage verletzt; [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) sagt Live-Reload zu, nicht Dauerhaftigkeit | die Grenze ist **still**: der Antrags-Datensatz liest `applied`, während der Wert einer als sensibel ausgeschlossenen Spalte wieder erfasst und persistiert wird — genau der Zustand, den [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 3 Option D als Wirkort ausschließt; ein persistierter Wert ist nicht zurückholbar, und `LH-QA-SEC-004`s Zusage hängt unbemerkt an der Prozesslebensdauer; ein benannter Träger allein macht den Verlust nicht unmöglich, nur lesbar |
| B — Grenze benennen und laut machen (Start-Warnung oder Statuszeile über nicht wiederhergestellte Ausschlüsse, Registereintrag, Folge-Slice-Vorbehalt) | billig, ehrlich, kein Schema-Objekt; macht den Verlust beobachtbar statt still | löst die Sache nicht: die Erfassung persistiert den Wert weiterhin; eine Warnung stellt die Zusage von `LH-QA-SEC-004` nicht her, sie verschiebt nur den Zeitpunkt, zu dem der Betreiber den Verlust bemerkt — und der bereits persistierte Wert bleibt |
| C — Ausschluss zusätzlich als Publication-Spaltenliste (`ALTER PUBLICATION … SET TABLE t (…)`); die Quelle filtert dann selbst | dauerhaft ohne Go-seitigen Zustand, und der Wert verlässt die Quelle nie — stärker als `LH-QA-SEC-004` verlangt | verlegt den Wirkort aus dem `Assembler` in die SQL-Administration und damit gegen [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 3 Option D; `cdc.enable_table`/`cdc.disable_table` tragen die Publication bereits, die Spaltenliste müsste bei jedem Aktivierungs-Zyklus neu gesetzt werden; die eingehende `decode.Relation` trägt dann eine andere Spaltenform, was die Konvergenz der Teilfrage 4 mit `LH-FA-SCH-003` berührt — eine eigene, größere Entscheidung, die diese Frage nicht braucht |
| D — eigener Zustand neben `cdc.source_table` (z. B. `cdc.excluded_column` mit Quelle/Tabelle/Spalte als Zeilen) | dauerhafter Stand symmetrisch zur Tabellen-Aktivierung; ein Leseort je Tabelle | zweiter Schreibpfad auf einen administrativen Zustand (die SQL-Funktion schreibt weiterhin nur den Antrag; der Go-Use-Case müsste den Stand zusätzlich pflegen) und eine **zweite Quelle der Wahrheit** neben der Antrags-Historie, für zwei Aspekte derselben Bindung; neue Schema-Objekte samt Migration ([`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md)), ohne dass ein Akzeptanzkriterium sie verlangt; der Antrag selbst bliebe unpfleglich — die Ableitung ist der kleinere Träger |
| **E — dauerhafter, tabellen-scoped Stand, aus der Antrags-Historie abgeleitet (gewählt)** | keine zweite Quelle: der `applied`-Antrag selbst wird wieder ausgewertet — dieselbe Herkunft, die der Rückkanal schon trägt; kein neues Schema-Objekt, kein zweiter Schreibpfad; **derselbe** Mechanismus deckt Prozess-Neustart und Bindungs-Zyklus, statt zwei Sonderfälle zu bauen; der stille Erfolg (F-2) löst sich mit derselben Festlegung auf: der Antrag wirkt, sobald die Tabelle erfasst wird | die Antrags-Historie wird **tragend**: sie darf nicht bereinigt werden, ohne den Stand zu verlieren (Re-Evaluierungs-Trigger unten); der Startpfad liest die Spalten-Antragsarten zusätzlich (eine Abfrage je Quelle, nicht je Tabelle) |

**Zur Alternative B im Verhältnis zu dieser Entscheidung:** Der Reviewer hat
einen Träger für die Entscheidung „zulässige Grenze oder Lücke" verlangt.
Diese ADR entscheidet „Lücke" und benennt mit `slice-075` <!-- d-check:status-provenance -->
den Träger ihrer Auflösung. Die Grenze bleibt zwischen diesem Verdikt und der
Umsetzung bestehen und wird dort sichtbar geführt (§Konsequenzen) — sie ist
damit benannt **und** adressiert, nicht nur benannt.

## Konsequenzen

- Positiv: `LH-QA-SEC-004`s Zusage hängt nicht mehr an der Prozesslebensdauer.
  Ein ausgeschlossener Wert erreicht die Persistenzschicht weder im ersten
  Prozess noch nach einem Neustart noch nach einem Bindungs-Zyklus.
- Positiv: **ein** Mechanismus für beide Auslöser. Prozess-Neustart und
  `disable`/`enable`-Zyklus sind dieselbe Frage („woher kommt der Stand beim
  Anlegen einer Bindung") und bekommen dieselbe Antwort, statt einer
  Neustart-Wiederherstellung plus einer Bindungs-Sonderregel.
- Positiv: der Antrags-Status wird wieder wahr. `applied` trägt den
  dauerhaften Stand; die Konstellation, in der der gewählte Rückkanal Erfolg
  meldete, ohne dass etwas wirken konnte, hat mit dem dauerhaften Träger kein
  Objekt mehr.
- Negativ: die Antrags-Historie ist **tragend**. Eine Bereinigung oder
  Alters-Retention auf `cdc.administration_request` verlöre den Stand; sie ist
  damit eine Entscheidung, keine Aufräumarbeit (Re-Evaluierungs-Trigger).
- Negativ: der Startpfad liest eine weitere Quelle je Quelle (die
  Spalten-Antragsarten); der Outbound Port, der die Spaltenexistenz prüft,
  bekommt eine Lesefähigkeit für den Stand
  ([`ADR-0034`](0034-ports-nach-faehigkeiten.md) — Fähigkeit am bestehenden
  Zuschnitt, kein zweiter Port).
- Negativ: die Ableitung setzt voraus, dass die Antrags-Zeilen die einzige
  Herkunft des Standes bleiben. Die in
  [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 1 offen gelassene
  Boot-Zeit-Konfiguration für initial ausgeschlossene Spalten bekäme mit
  dieser Entscheidung einen zweiten Herkunftsort und bräuchte eine eigene
  Festlegung, wie beide zusammengeführt werden — sie bleibt, wie in
  [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) festgehalten, nicht
  Gegenstand dieser Entscheidung (Re-Evaluierungs-Trigger).
- Folgepflicht: `slice-075` <!-- d-check:status-provenance --> (`docs/plan/planning/open/`)
  setzt die Entscheidung um: Ableitung des Standes aus den `applied`-Zeilen,
  Lesefähigkeit am Outbound Port, Wiedereinspielung im Prozessstart **und** im
  Aktivierungs-Zweig, Unit-Test (Neuaufbau der Bindungen und
  `disable`/`enable`-Zyklus) sowie ein E2E-Beleg über den bestehenden
  Neustart-Rundlauf in `tools/harness/run-integration-tests.sh`. Die
  endgültige Planung führt der Planner.
- Folgepflicht: `SPEC-019` ([`SPEC-019`](../../../spec/pflichtenheft.md))
  beschreibt in seinem Fließtext, was `applied` für die beiden
  Spalten-Antragsarten bedeutet (dauerhaft vermerkter Stand, Wirksamkeit
  sobald die Tabelle erfasst wird) — Planner-/Architect-Zug, sobald die
  Umsetzung steht.
- Folgepflicht: `slice-067` <!-- d-check:status-provenance --> §6 (Planner-Zug)
  trägt für das nachgetragene Risiko den Ausgang **eingetreten → `slice-075`** <!-- d-check:status-provenance -->
  und benennt in seinem Text den zweiten, neustart-unabhängigen Auslöser
  (`RemoveBinding`/`AddBinding`); §7 verweist auf den Registereintrag unten.
- Folgepflicht: `welle-18` §6 (Planner-Zug) benennt die Grenze, die bis zur
  Umsetzung von `slice-075` <!-- d-check:status-provenance --> gilt — der
  Ausschlussstand wirkt über die Lebensdauer eines Prozesses und über einen
  Bindungs-Zyklus hinweg neu beantragt werden, ein Neustart stellt ihn nicht
  her. Der Closure-Trigger der Welle bleibt davon unberührt (§Re-Evaluierungs-
  Trigger dieser ADR, Verdikt-Abschnitt des Review-Nachzugs).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` (`internal/bootstrap`) | Der Neuaufbau der Bindungen (`activatedTableBindings`) und ein `disable`/`enable`-Zyklus führen denselben Ausschlussstand zurück, der aus den `applied`-Zeilen abgeleitet ist; ein Change danach trägt den Spaltennamen nicht im Row Image | `make test` (bestehendes Ziel) |
| `run-integration-tests.sh` | Ein Container-Neustart lässt einen beantragten Ausschluss wirksam: ein danach eingefügter Change trägt den ausgeschlossenen Spaltenwert nicht in `cdc.changes` | `make test-integration` (kein Gate, [`ADR-0030`](0030-testpyramide.md)) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Zwei Trigger, beide benannt statt „permanent"; die Festlegungen 1, 3 und 4
sind von keinem der beiden betroffen:

1. **Die Antrags-Historie wird bereinigt.** Führt ein Slice eine
   Alters-Retention, eine Archivierung oder eine Löschung auf
   `cdc.administration_request` ein, verliert die Ableitung (Festlegung 2)
   ihren Gegenstand — dann entscheidet eine Folge-ADR zwischen einem eigenen
   Zustand (Option D oben) und einem Bereinigungs-Ausschluss für die beiden
   Spalten-Antragsarten.
2. **Ein Boot-Zeit-Ausschlussfeld entsteht.** Bekommt `CDC_TABLES` oder die
   optionale YAML-Konfigurationsdatei
   ([`ADR-0052`](0052-optionale-yaml-konfigurationsdatei.md)) ein Feld für
   initial ausgeschlossene Spalten, hat der Stand einen zweiten Herkunftsort;
   die Zusammenführung beider Herkünfte wird dann per Folge-ADR entschieden.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: `slice-067` <!-- d-check:status-provenance -->-Review führt F-1 (Ausschlussstand ohne dauerhaften Träger, über Prozess-Neustart **und** `disable`/`enable`-Zyklus) und F-2 (stiller Erfolg eines nie wirksamen Antrags) als Entscheidungen an Planner → Architect (Modul 8 §Konflikt-Pfad); unabhängiger Architect-Zug entscheidet den dauerhaften Träger und korrigiert die Dauerhaftigkeits-Aussage in [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) | Review zu `slice-067`, der Architect-Verdikt zur Spaltenausschluss-Dauerhaftigkeit |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0065` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
