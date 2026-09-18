# Architect-Verdikt: Spaltenausschluss — Dauerhaftigkeit des Ausschlussstandes und Erfolgsmeldung ohne Wirkung

**Rolle:** Architect (Modul 8)

**Anlass:** Die zwei MEDIUM-Findings F-1 und F-2 aus
dem Review zu `slice-067`. Der Reviewer hat
sie ausdrücklich als **Entscheidungen** an Planner/Architect gereicht (kein
Implementer-Pfeil, keine Fixrunde) — damit läuft der Rollenwechsel
Planner → Architect → Planner aus Modul 8 §Konflikt-Pfad, Verdikt 2/3
(Folge-ADR statt stiller Lockerung).

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der
Implementer-Lauf von `slice-067` und als der Reviewer-Lauf)

**Datum:** 2026-09-14

**Bezug:** [`LH-FA-CFG-005`](../../spec/lastenheft.md) (Happy Path, Boundary,
Negative), [`LH-QA-SEC-004`](../../spec/lastenheft.md) (Ausschluss sensibler
Spalten — Messmethode über `LH-FA-CFG-005`),
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (Mechanismus,
Granularität, Wirkort, Konvergenz, Rückkanal — in der
Dauerhaftigkeits-Aussage durch die Folge-ADR korrigiert),
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue, Live-Reload), [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md)
(Port-Zuschnitt), [`SPEC-019`](../../spec/pflichtenheft.md) (Feldform und
Bedeutung des Antrags-Datensatzes),
`docs/plan/planning/welle-18.md` (§3
Closure-Trigger, §4 Slices, §6 Out-of-Scope), `slice-068`
(E2E-Beleg des laufenden Pfads; als Kennung zitiert, nicht als Pfad-Link —
ein Slice wechselt die Lifecycle-Ablage und ein Pfad-Link bräche mit)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  (Accepted, `Supersedes ADR-0059` — nur die Dauerhaftigkeits-Aussage aus
  Teilfrage 1 Option D und §Bestätigung), samt Index-Zeile in
  [`docs/plan/adr/README.md`](../plan/adr/README.md)
- ``slice-075``
  (`docs/plan/planning/open/`) — Adresse der Umsetzung
- Registereintrag
  [`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`](../plan/planning/observations/BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/observation.md)
  mit Beleg im Review zu `slice-067`

---

## Verdikt

**F-1 — Lücke, keine zulässige Grenze.** Der Ausschlussstand braucht einen
dauerhaften Träger. `ADR-0059` sagt ihn nicht zu (sie sagt Live-Reload zu),
aber die zwei Stellen, die eine Parität zur Tabellen-Aktivierung nahelegen,
vergleichen mit einer dauerhaft liegenden Fähigkeit — und für ein
sicherheitsrelevantes Ausschlussmerkmal (`LH-QA-SEC-004`) ist ein Verlust,
den weder ein Statuswert noch ein Log noch ein Gate meldet, kein
Betriebsdetail. Die Entscheidung ist mit `ADR-0065` gefallen, die Umsetzung
hat mit `slice-075` eine Adresse, die Beobachtung steht im Register.

**F-2 — kein `failed`-Pfad, sondern derselbe Träger.** Der stille Erfolg hat
dieselbe Wurzel wie F-1: der Stand hängt an einer Bindung statt an der
Tabelle. Mit dem dauerhaften, tabellen-scoped Stand wirkt der Antrag, sobald
die Tabelle erfasst wird — der Antrags-Status bleibt der einzige Rückkanal
(`ADR-0059` Teilfrage 5) und wird wieder wahr. Ein Statuswechsel allein
träfe das Symptom und ließe den Zielzustand unerreichbar.

**Kein Code-Eingriff, keine Fixrunde.** Keines der beiden Findings ist eine
Korrektur am Diff von `slice-067`; beide sind Entscheidungen. `slice-067`
bleibt ohne Änderung an `internal/**`.

**Konsequenz für `welle-18`:** Der Closure-Trigger der Welle bleibt
erfüllbar. Sie schließt über den realen E2E-Beleg aus `slice-068`; die
Grenze bis zur Umsetzung von `slice-075` wird benannt geführt (§Frage 3).

---

## Befundlage (eigener Code-Gang)

Drei Eigenschaften des ausgelieferten Standes, am Code nachvollzogen:

1. **Der Prozessstart trägt keinen Ausschlussstand.**
   `activatedTableBindings` (`internal/bootstrap/wiring.go`) baut je
   registrierter Tabelle `mapper.TableBinding{TableID, SchemaVersion}` aus
   `TableActivationPort.List` und `SchemaStorePort.CurrentVersion`; die
   Spaltennamen der Antragsarten `exclude_column`/`include_column` liest
   dieser Pfad nicht.
2. **Ein Bindungs-Zyklus löscht den Stand.** `Assembler.RemoveBinding`
   (`internal/adapters/driving/replication/mapper/mapper.go`) entfernt den
   Karteneintrag; `Assembler.AddBinding` übernimmt `ExcludedColumns` nur aus
   einem **vorhandenen** Eintrag. Der Deaktivierungs-Zweig ruft
   `RemoveBinding`, der Aktivierungs-Zweig `AddBinding` — die Reihenfolge
   disable → enable hinterlässt einen leeren Stand. Dieser Auslöser braucht
   keinen Prozess-Neustart.
3. **Der Rückkanal meldet weiter Erfolg.** `processAdministrationRequests`
   liest `pending`-Zeilen und setzt `applied`; die `applied`-Zeile in
   `cdc.administration_request` ist der einzige dauerhafte Beleg des
   Antrags. Nach (1) oder (2) trägt sie Erfolg, während die Spalte wieder
   erfasst wird.

Der Reviewer-Befund ist damit in beiden Hälften bestätigt: der
Neustart-Auslöser (vom Implementer benannt) und der Bindungs-Zyklus (im
Review ergänzt, in Zeilen, die der Diff von `slice-067` anfasst).

---

## Frage 1 — Fehlende Dauerhaftigkeit: zulässige Grenze oder Lücke?

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: Stand bleibt prozessgebunden, Grenze wird nur benannt (§6 des Slice, `welle-18` §6, Registereintrag) | kein Code, kein Schema-Objekt; `ADR-0059` sagt keine Dauerhaftigkeit zu, also kein ADR-Verstoß | die Grenze ist **still**: der Antrags-Datensatz liest `applied`, während der Wert einer als sensibel ausgeschlossenen Spalte wieder persistiert wird — genau der Zustand, den `ADR-0059` Teilfrage 3 Option D als Wirkort ausschließt; ein persistierter Wert ist nicht zurückholbar; ein benannter Träger ohne Adresse macht den Verlust lesbar, nicht unmöglich |
| B — Grenze benennen **und laut machen** (Start-Warnung oder Statuszeile über nicht wiederhergestellte Ausschlüsse, Registereintrag, Folge-Slice-Vorbehalt) | billig, ehrlich, kein Schema-Objekt; der Verlust wird beobachtbar | stellt die Zusage nicht her: die Erfassung persistiert den Wert weiterhin; die Warnung verschiebt nur den Zeitpunkt, zu dem der Betreiber den bereits eingetretenen Verlust bemerkt |
| C — Ausschluss zusätzlich als Publication-Spaltenliste (`ALTER PUBLICATION … SET TABLE t (…)`) | dauerhaft ohne Go-seitigen Zustand; der Wert verlässt die Quelle nie | verlegt den Wirkort gegen `ADR-0059` Teilfrage 3 Option D in die SQL-Administration, verdoppelt die Auswertung, muss bei jedem Aktivierungs-Zyklus neu gesetzt werden und ändert die eingehende Spaltenform — berührt damit die Konvergenz der Teilfrage 4 (`LH-FA-SCH-003`); eine eigene, größere Entscheidung |
| D — eigener Zustand neben `cdc.source_table` (`cdc.excluded_column` je Quelle/Tabelle/Spalte) | dauerhafter Stand symmetrisch zur Tabellen-Aktivierung | zweiter Schreibpfad auf einen administrativen Zustand und **zweite Quelle der Wahrheit** neben der Antrags-Historie, für zwei Aspekte derselben Bindung; neue Schema-Objekte samt Migration (`ADR-0043`), ohne dass ein Akzeptanzkriterium sie verlangt |
| **E — dauerhafter, tabellen-scoped Stand, aus den `applied`-Zeilen der Spalten-Antragsarten abgeleitet (gewählt in `ADR-0065`)** | keine zweite Quelle: der Antrag selbst wird wieder ausgewertet; kein neues Schema-Objekt, kein zweiter Schreibpfad; **ein** Mechanismus für Prozess-Neustart **und** Bindungs-Zyklus; löst zugleich F-2 auf | die Antrags-Historie wird tragend (Bereinigung = Entscheidung, nicht Aufräumarbeit); der Startpfad liest die Spalten-Antragsarten zusätzlich |

**Fazit F-1:** Lücke. Die Alternative „benannte, aber unveränderte Grenze"
(A/B) trägt nicht, weil sie den Verlust des Schutzes nur beschreibt statt
ihn zu verhindern — bei einem Ausschlussmerkmal, dessen Zweck die
Nicht-Persistenz des Werts ist. Entschieden ist **E** in
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md);
`ADR-0059` bleibt in allen fünf Teilfragen gültig und ist nur in ihrer
Dauerhaftigkeits-Aussage (Teilfrage 1 Option D, §Bestätigung) superseded —
die Stellen, an denen der Vergleich mit der dauerhaften
Tabellen-Aktivierung in die Irre führt.

## Frage 2 — Stiller Erfolg für einen nie wirksamen Antrag

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: der No-op bleibt, der Plan-Nachzug und zwei Tests pinnen ihn | kein Eingriff; `LH-FA-CFG-005`s Negative-Kriterium ist nicht verletzt (die Spalte existiert) | der Zielzustand ist über diesen Mechanismus unerreichbar: eine spätere Aktivierung läuft über `AddBinding` ohne Vorbindung, der Antrag wird nicht erneut gelesen; der einzige gewählte Rückkanal meldet Erfolg für einen Antrag, der nie wirkt |
| B — `failed` mit Fehlertext für eine Tabelle ohne laufende Bindung (neuer Sentinel, Muster von `ErrSourceTableMissing`) | der Status sagt die Wahrheit über die Wirkung im laufenden Prozess; erzwingt die in `ADR-0059` Teilfrage 1 vorgesehene Reihenfolge (Ausschluss nach der Erstaktivierung) | verbietet eine Vor-Konfiguration dauerhaft und trifft das Symptom: der Stand hätte weiterhin keinen Ort bei der Tabelle; mit dem dauerhaften Träger (Frage 1) wird die Ablehnung gegenstandslos — sie wäre eine zweite, strengere Regel, die kein Akzeptanzkriterium verlangt |
| C — eigener Status oder Hinweisfeld (z. B. `deferred` bzw. ein Warnungs-Feld) | macht den aufgeschobenen Fall unterscheidbar sichtbar | erweitert die geschlossene Statusmenge `pending`/`applied`/`failed` und damit `SPEC-019` um einen vierten Wert, für eine Konstellation, die mit E keinen Sonderzustand mehr braucht; mehr Oberfläche als die Frage trägt |
| **D — derselbe dauerhafte Träger: `applied` heißt „dauerhaft vermerkt"; der Antrag wirkt, sobald die Tabelle erfasst wird (gewählt in `ADR-0065`)** | löst beide Findings mit einer Festlegung; der Antrag bekommt einen erreichbaren Zielzustand; Statusmenge und Rückkanal bleiben unverändert (`ADR-0059` Teilfrage 5) | verlangt einen nachgezogenen `SPEC-019`-Fließtext, der die Bedeutung von `applied` für die beiden Spalten-Antragsarten ausspricht (Folgepflicht, `ADR-0065`) |

**Fazit F-2:** Der stille Erfolg ist nicht durch einen Statuswechsel
aufzulösen, sondern durch den Träger. Die Konstellation, in der der
Rückkanal Erfolg meldete, ohne dass etwas wirken konnte, hat mit dem
dauerhaften Stand kein Objekt mehr: der Ausschluss gehört zur Tabelle und
greift bei der ersten Erfassung. B bleibt als geprüfte und verworfene
Option dokumentiert; ein `failed` für eine noch nicht erfasste Tabelle wäre
eine Regel, die kein Akzeptanzkriterium fordert.

## Frage 3 — Konsequenz für `welle-18` und `slice-068`

- **Der Closure-Trigger der Welle bleibt erfüllbar.** `welle-18` §3 schließt
  über `slice-066`/`slice-067`/`slice-068` in `done/`, `make gates` grün und
  den realen E2E-Beleg aus `slice-068`. Alle drei Bedingungen sind von
  diesem Verdikt unberührt; `slice-075` ist **nicht** Teil des Triggers.
  Wird er dennoch in die Welle gezogen, ist das eine Planner-Umplanung
  (§4/§5 der Welle nachziehen, Drift-Log-Eintrag, Modul 6).
- **`slice-068` bleibt, was es ist.** Sein Rundlauf belegt den laufenden
  Pfad (SQL-Antrag → Live-Reload → gefiltertes Row Image) und läuft ohne
  Neustart. Er belegt damit ausdrücklich **keine** Dauerhaftigkeit; wer ihn
  als Neustart-Beleg liest, liest über `ADR-0065` hinaus. Der Neustart-Beleg
  gehört zu `slice-075` (sein E2E-Teil).
- **Der Risiko-Ausgang in `slice-067` §6 lautet *eingetreten → `slice-075`*,
  nicht „weiter offen"** (Planner-Zug). Der Reviewer hat diesen Ausgang an
  eine vorher zugewiesene Entscheidung gebunden; die Entscheidung liegt mit
  `ADR-0065` vor. Der Risiko-Text bekommt den zweiten, neustart-unabhängigen
  Auslöser (`RemoveBinding`/`AddBinding`) dazu.
- **`welle-18` §6 benennt die Grenze** (Planner-Zug): bis zur Umsetzung von
  `slice-075` gilt der Ausschlussstand über die Lebensdauer eines Prozesses;
  ein Neustart und ein Bindungs-Zyklus stellen ihn nicht her, der Antrag
  bleibt `applied`.
- **Wellen-Zuordnung von `slice-075`:** Planner-Entscheidung. Das Kopf-Feld
  der Datei trägt `—` mit dem Verweis auf dieses Verdikt, statt eine
  Zuordnung zu behaupten.

## Frage 4 — Registereintrag: Kennung und Klasse

Die vorgeschlagene Kennung `BEO-PGC/spaltenausschluss-nur-prozesslebensdauer`
benennt die **Instanz** (den Spaltenausschluss). Die Register-Konvention
dieses Repos benennt die **Klasse** (`pipe-maskiert-make-exit-code`,
`slice-chronik-in-code-kommentar`, `report-nackte-id-ohne-link`). Die Klasse
ist hier weiter als der Spaltenausschluss: *ein administrativer Zustand, den
ein `applied`-Antrags-Datensatz als verarbeitet ausweist, aber nur im
Prozessspeicher lebt*. Angelegt ist deshalb
[`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`](../plan/planning/observations/BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/observation.md)
— ein enger gefasster Slug hätte einen künftigen zweiten Fall derselben
Klasse nicht aufnehmen können und die Beobachtung gespalten (Modul 6: eine
Kennung zitieren, nicht neu formulieren).

Kein bestehender Eintrag nimmt den Fund auf: `BEO-PGC/schema-evolution-nicht-dynamisch`
trägt die dynamische Re-Versionierung (verkörpert seit `slice-033`),
`BEO-PGC/kein-admin-weg-schema-fehler-recovery` einen fehlenden
Administrationsweg nach einem `schema`-Abbruch,
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` eine Testform-Lücke —
drei andere Klassen. Der Zustand steht auf **geplant** (die Entscheidung ist
gefallen, der Träger ist benannt); die Beobachtung braucht die 3×-Schwelle
nicht, dieselbe direkte Auflösung unter der Schwelle wie
`BEO-PGC/schema-evolution-nicht-dynamisch` und
`BEO-PGC/retention-keine-loeschausfuehrung`. Der Beleg trägt den Vorgang
des Reviews zu `slice-067`; `slice-067` und sein Review sind **eine** Gelegenheit —
ein zweiter Beleg für `slice-067` bei der Slice-Closure wäre ein zweites
Auftreten derselben Gelegenheit und gehört nicht angelegt.

---

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `docs/plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md` (neu,
  Accepted) und `docs/plan/adr/README.md` (Zeile `ADR-0065`; Zeile
  `ADR-0059` mit dem Supersedes-Vermerk)
- `docs/plan/planning/open/slice-075-dauerhafter-ausschlussstand-wiedereinspielung.md`
  (neu — Adresse der Verdikt-Auflage; die endgültige Planung führt der
  Planner)
- `docs/plan/planning/observations/BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/`
  (neu: `observation.md`, `state.md`, Beleg im Review zu `slice-067`)
- diese Verdikt-Datei

**Nicht geändert:** `internal/**` (kein Code-Eingriff), `spec/**`,
`docs/plan/planning/in-progress/slice-067-*.md`,
`docs/plan/planning/welle-18.md`, `docs/plan/planning/open/slice-068-*.md`
und die übrigen Planungsdokumente — sie führt der Planner-Zug nach diesem
Verdikt (Risiko-Ausgang, §6-Wortlaut der Welle, `SPEC-019`-Fließtext).

**Offen für den Planner-Zug:** Risiko-Ausgang in `slice-067` §6
(*eingetreten → `slice-075`*), Grenz-Benennung in `welle-18` §6,
`SPEC-019`-Fließtext zur Bedeutung von `applied` (mit der Umsetzung), sowie
die Wellen-Zuordnung von `slice-075`.
