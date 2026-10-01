# Welle welle-routing — Routing auf Zustellziele: erfasste Changes tragen bei der Erfassung das durch geordnete Regeln bestimmte Ziel als persistiertes Label, konfiguriert über die SQL-Antrags-Queue und wählbar an jedem Lesezugriffsweg — Closure-Notiz

**Welle:** welle-routing
**Abschluss:** 2026-10-02
**Verantwortlich:** pt9912

## Was wurde geliefert?

- **[`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Happy Path, Boundary, Negative) ist
  am laufenden Feed-Container belegt.** Ein Antrag `cdc.set_route` legt eine geordnete
  Routing-Regel (Herkunfts- oder Inhaltsregel) auf eine aktivierte Tabelle; jede danach
  erfasste Change trägt das Ziel der ersten treffenden Regel als persistiertes Label
  `route_target` (leer, wenn keine trifft). Der Leser wählt das Ziel über `cdc.changes`,
  `GET /changes`, den RPC `ReadChanges`, den gRPC-Stream und den SSE-Stream
  (Parameter `target`) und über das Zusatz-Subjekt `cdc.route.<source_id>.<ziel>` des
  NATS-Vollinhalts-Wegs; ein ungefilterter Leser sieht alles. Zwei treffende Regeln
  löst die kleinere `order`; je eine Verletzung von R1 bis R6 endet `failed` mit dem
  Klartext der Spec und lässt den Regelstand stehen; eine am laufenden Prozess
  nicht anwendbare Regel beendet den Erfassungspfad sichtbar mit Fehlerklasse `schema`,
  und die Abhilfe (`cdc.remove_route` bei stehendem Prozess, Neustart) ist an den
  erzeugbaren Fällen belegt. Die Auswertung liegt an genau **einer** Stelle
  (`model.EvaluateRoute`), aufgerufen vom WAL-Pfad (`mapper.go`) **und** vom
  Backfill-Pfad (`usecase/backfill/service.go`) — Suchlauf im Abschnitt Verifikation.
- **Zehn Slices in `done/`:** `slice-routing-spec-nachzug` (Pflichtenheft,
  Architektur-Sicht, neue Regelform-Kennung [`SPEC-032`](../../../../spec/pflichtenheft.md),
  [`SPEC-019`](../../../../spec/pflichtenheft.md) mit neun Antragsarten),
  `slice-routing-kern-label` (Domäne, Auswertung im `Assembler`, Spalte `route_target`,
  View-Signatur von `cdc.changes` mit Rollout-Vorlauf), `slice-routing-antragsweg`
  (Antragsarten `set_route`/`remove_route`, zwei SQL-Funktionen allein für `cdc_admin`,
  R1 bis R6 einschließlich der Sperre gegen `exclude_column`, Regelstand-Ableitung),
  `slice-routing-backfill-pfad` (Label im Backfill-Run, Regelstand fail-closed,
  Nichtanwendbarkeit run-lokal), `slice-routing-lesewege` (`target` an `GET /changes`,
  `ReadChanges`, gRPC- und SSE-Stream, eine Filterfunktion), `slice-routing-nats-subjekt`
  (Zusatz-Subjekt), `slice-routing-e2e` (Runner-Phasen und vier Go-Testfunktionen am
  laufenden Feed-Container, DELETE-Messung, RTM-Träger), `slice-routing-betriebsdoku`
  (Benutzerhandbuch), `slice-routing-sdk-beispiel-target` (`target` in den drei
  SDK-Packages und den Beispiel-Clients), `slice-routing-sdk-realserver-e2e` (vier
  Routing-Phasen je SDK-Tier).
- **Fünf Entscheidungen:**
  [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  (Mechanismus: benannter Kanal je Quelle, Label im `Assembler` vor der
  Persistierung, zwei Antragsarten, R1–R6),
  [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  (V1/V2: Feld `target` an `ReadChanges`, Nichtanwendbarkeit im Run),
  [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  (Regelstand im Run fail-closed, `target` außerhalb des Alphabets liefert leer),
  [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  (V3: Erreichbarkeit der Nichtanwendbarkeit, Grenze der Abhilfe-Zusage) und
  [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md)
  (Lesefehler eines Run-Regelstands endet mit der Klasse der Ursache). Lastenheft 0.14.0
  präzisiert das Happy-Path-Kriterium („zugestellt“ heißt auswählbar unter dem Zustellziel).
- **Schema:** `cdc.change` trägt die nullable Spalte `route_target`, `cdc.changes` sie als
  letzte Spalte (Rollout über den Vorlauf `DROP VIEW`, [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md));
  zwei neue SQL-Funktionen (`cdc.set_route`, `cdc.remove_route`), `request_kind` wächst
  von sieben auf neun Werte.
- **Benutzerhandbuch:** Änderungshistorie 1.84 bis 1.86 (Routing-Regel konfigurieren,
  `target` an den Roh-Wegen und den SDKs, Zusatz-Subjekt, Spalte `route_target`,
  Fehlerklasse `schema` in zwei Ursachen).
- **RTM:** `make doc-trace` druckt `80 Anforderung(en), 0 Waise(n).` —
  [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) trägt die Nachweise `E2E, SDK-E2E`
  (Parent `30fd6cb5`: `80 Anforderung(en), 1 Waise(n)`, **gemessen** am 2026-10-02 aus
  `git archive 30fd6cb5` im Scratchpad mit `make doc-trace`, Exit `0`; die Zeile
  `LH-FA-CFG-008 … ADR-0112, ADR-0137 | — | — | WAISE` und `80 Anforderung(en), 1 Waise(n).` wurden gedruckt).

## Was hat funktioniert?

- **Das Muster der Transformations-Welle trug den Schnitt ohne Rückführung.**
  Antragsarten mit Guard-Liste und Alt-Tag-Lauf (`antragsweg`), Regelstand mit Sentinel
  und Mengenvergleich im Run (`backfill-pfad`), der Filter-Parameter neben `schema`/`table`
  (`lesewege`): kein Slice musste die in seinem Plan vorab benannte Teilungsnaht auslösen.
- **Eine Auswertungsstelle, zwei Aufrufer.** `EvaluateRoute` wird gerufen, nicht
  nachgebaut; der Backfill-Slice konnte deshalb die Parität über den Typ-Satz als Test
  fassen (`checkRouteParity`), statt zwei Implementierungen zu vergleichen.
- **Lücken im Text einer Entscheidung gingen als Architect-Frage vor die Spec-Setzung.**
  A-1/A-2 (Regelstand im Run, `target` außerhalb des Alphabets) wurden
  [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
  V3 wurde [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
  der Widerspruch zwischen ADR-Wortlaut und Bestand im Backfill-Review
  [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md); der
  Implementer musste keine dieser Setzungen raten.
- **Messen statt Herleiten bei den zwei offenen Erwartungen der Welle-Eröffnung.** V3
  war *hergeleitet* (die Nichtanwendbarkeit sei am System nicht erzeugbar); die Messung in
  `slice-routing-e2e` erzeugt sie an (b) Erstaktivierung ohne Spaltenform und (c)
  Publication-Spaltenliste, die Verengung nach `ADR-0140` Entscheidung 4 greift nicht. Die
  DELETE-Erwartung von `ADR-0137` Entscheidung 4 war *hergeleitet*; sie bestätigt sich an
  beiden PostgreSQL-Versionen (Abschnitt Verifikation).
- **Ein Re-Review nach einer Fixrunde, die Anweisungen änderte, lieferte jedes Mal einen
  Befund** (`antragsweg` F-N1 LOW, `backfill-pfad` F-N1 MEDIUM, `lesewege` F-N3 INFO); bei
  fünf Fixrunden, die weder Produktionslogik noch eine Norm änderten und die der Verifier
  ausführte, blieb ein Re-Review aus, ohne dass ein Fehler nachträglich auftrat (Lese-Schritt unten).
- **Eine Fixture-Quelle für drei Sprach-Tiers.** Vorbereitung und Phasenablauf der
  SDK-Routing-Phasen stehen einmal in `tools/harness/lib-sdk-route-fixture.sh`; dieselbe
  Eingabetabelle gilt je Sprache, die Schlusszeilen sind in Schwellen und Reihenfolge
  gleich gelesen.

## Was ging anders als geplant?

- **Alle zehn Slices brauchten eine Fixrunde** (§7 der zehn Slice-Pläne, **übernommen**).
  Drei davon änderten eine Norm oder Anweisungen mit Wirkung:
  `antragsweg` (HIGH F-1: die Annahme, `jsonb` bewahre die Schreibweise der Zahl,
  war eine Tatsachenbehauptung des Plans; gemessen an PostgreSQL 18 normalisiert `jsonb`
  `1e1` zu `10`), `backfill-pfad` (zwei Fixrunden; die erste folgte dem Wortlaut von
  `ADR-0139` und lief in die falsche Richtung, `ADR-0141` berichtigte den Wortlaut) und
  `lesewege` (Rangfolge „ungültiges Ziel gegen Fehler des Lese-Kontrakts“, vom Hauptlauf
  entschieden; der Slice änderte damit die Spec, die der Plan als unberührt führte).
- **Drei Folge-Entscheidungen entstanden während der Welle** (`ADR-0139`, `ADR-0140`,
  `ADR-0141`); die Eröffnung kannte nur `ADR-0138`.
- **Der Alt-Tag-Lauf verlor Vorbedingungen** (`kern-label`): der jüngste `v*`-Tag trägt die
  Transformations-Objekte bereits, die Belegkraft für sie sank von „Upgrade ergänzt“ auf
  „Upgrade erhält“. Der Lauf am Endstand dieser Closure fährt gegen `v0.4.0`.
- **Die Negative-Phase von `slice-routing-e2e` steht als letzter Rundlauf des Runners**
  statt vor der Container-Ende-Grenze (Abweichung im Plan, Rückführung „(B) abtrennen“ nicht
  eingetreten).
- **Das C#-Integrationsprojekt übersetzte seit dem Stream-Filter nicht** (`CS1503`, drei
  Aufrufe), unbemerkt, weil kein regelmäßiger Lauf die Stufe `integration` baut; der
  Fix liegt in `slice-routing-sdk-realserver-e2e`, die Lücke im Sensor-Satz bleibt
  (Register unten).
- **Hergeleitet statt gemessen bleibt** (Welle-Grenze, nicht Slice-Mangel): die Abhilfe am
  Fall (a) „Bedingungsspalte entfernt“; der Fall „Publication-Spaltenliste bei bereits
  bekannter Spaltenform“; das Ruhefenster der Negativ-Zählung (15 s, feste Menge NULL/asia/us/eu;
  der Beleg gilt „im Ruhefenster“, nicht „nie“); die Lesekosten der zwei Regelstände im
  Backfill-Run und die Last der Auswertung je Change (beide **nicht gemessen**, Abschnitt
  Offene Punkte); das Verhalten eines Packages gegen einen Server ohne den Parameter `target`.
- **Zwei Stellen der Welle-Datei trugen den Stand der Eröffnung und sind mit dieser Closure
  nachgezogen:** V3 („Systemmessung von (b)/(c) offen“, jetzt gemessen) und die DELETE-Aussage
  („bis dahin nur erwartet“, jetzt an beiden Versionen bestätigt).

## Steering-Loop-Einträge

Die zehn Slices tragen je einen Lerneintrag in ihrer Closure-Notiz (§7). Keiner
verkörpert eine Regel mit Feld `liegt in` (alle sind „geschärfte Handlung, kein neuer
Sensor“, Träger die Lese-Handlung einer Rolle); diese Closure fasst sie zusammen und
ergänzt sie um den Lese-Schritt des Registers.

- `spec-nachzug`: Ein Fixrunden-Bericht, der „X ergänzt“ sagt, ist eine Trägerzusage und wird
  gegen den **Diff** geprüft; eine Lücke im Text der Entscheidung wird vor dem Spec-Text zur ADR.
- `kern-label`: Eine Vorbedingung am jüngsten `v*`-Tag altert mit jedem Release und wird
  an eine feste Referenz gebunden oder als Risiko mit Wächter geführt; die Zahl einer
  Mutationsreihe trägt ihren Ursprung (Autor-Zahl **übernommen**, Zahl des unabhängigen
  Lesers **gemessen**).
- `antragsweg`: Ein Plan-Satz über das Verhalten des Fremdsystems ist eine
  Tatsachenbehauptung und gehört vor dem Schreiben an der gepinnten Instanz gemessen; ein
  Test fährt den Weg, den die Eingabe real geht.
- `backfill-pfad`: Erweist sich der Wortlaut einer `Accepted`-ADR am Bestand als unzutreffend,
  ist die Klärung ein Architect-Verdikt, nicht die Anpassung des Codes; vor einer Fixrunde,
  die dem Wortlaut folgt, die Nachbarstände samt Tests lesen.
- `lesewege`: Eine Fehlerrangfolge ist Vertrag des neuen öffentlichen Parameters und gehört
  vor den Code in die Zeile des Parameters („Was gewinnt bei gleichzeitigem Kontraktfehler?“).
- `nats-subjekt`: Eine Unabhängigkeits-Zusage zweier Seiteneffekte braucht beide
  Richtungen als Test; ein nur mit Fremdsystem laufender Test darf sich nicht still überspringen.
- `e2e`: Ein „nur“-Satz an einem Weg ohne Abschluss ist erst belegt, wenn der Client nach dem
  ersten Treffer weiterzählt; die Erwartung folgt der Zusage (At-least-once), nicht dem Normalfall.
- `betriebsdoku`: In einem Betreiber-Handbuch ist ein Ursprung nur so viel wert, wie der Leser
  ihn auflösen kann (verlinkter, committeter Report, sonst keine Einzelzahl).
- `sdk-beispiel-target`: Ein Docker-Cache kann Testläufe still überspringen; ein
  `make sdk-pack-*`-Lauf aus dem Cache belegt keine Testausführung (Register, 1×).
- `sdk-realserver-e2e`: Ein Literal im Ergebnisfeld einer Belegzeile ist kein Messwert; eine
  DoD-Mutation, die der Server ablehnt, ist keine Eingabeseiten-Mutation.
- **Welle-Befund (diese Closure):** die Fixrunde ist in dieser Welle die Regel (zehn von zehn
  Slices), und die Frage „Re-Review oder Gegenprüfung“ stellte sich in jedem. Die Lage ist
  nicht mehr „jede Fixrunde braucht einen Re-Review“, sondern: der Re-Review lohnte, wo die
  Fixrunde Anweisungen oder eine Norm änderte (drei von drei fanden etwas), und blieb
  entbehrlich, wo ein zweiter Kontext sie **ausgeführt** hat. Der Eintrag
  `BEO-PGC/fixrunde-ohne-reviewer-lesung` trägt den Wortlaut; sein Ausgang liegt im
  Lese-Schritt unten.

### Lese-Schritt des Beobachtungs-Registers

Das Register führt 140 Verzeichnisse (`ls docs/plan/planning/observations/BEO-PGC | wc -l`,
gemessen am Stand `84f5b8d2`), 52 davon mit `evidence/` ab 3 Dateien (Schleife über
`ls <Eintrag>/evidence | wc -l`, gemessen am selben Stand); 27 `evidence/`-Dateien dieses Standes
tragen den Namen eines Slices dieser Welle (`ls */evidence/slice-routing-* | wc -l`). Fünf Verzeichnisse
sind in der Welle entstanden (jede ihrer `evidence/`-Dateien trägt einen Routing-Slice-Namen):
`docker-cache-ueberspringt-tests-still`, `fixrunde-ohne-reviewer-lesung`,
`integrationsprojekt-uebersetzt-nicht-unbemerkt`, `suppression-ohne-linter-in-testcode` und
`test-strenger-als-die-zusage`.

**Bereits mit zugewiesenem Ausgang (verkörpert oder gestrichen)** — kein Zug nötig, der Zähler
steht in der Beobachtung. Neun dieser zehn Einträge sind in dieser Welle gewachsen (je 1 bis 3
`evidence/`-Dateien mit einem Routing-Slice-Namen, `ls <Eintrag>/evidence | grep -c slice-routing`,
gemessen am Stand `cff48b65`); `arbeit-ueberholt-stehenden-traeger` (34×) ist es **nicht** (derselbe
Befehl druckt `0`, der Eintrag steht unter einem Deckel ohne weitere Dateien):
`arbeit-ueberholt-stehenden-traeger` (34×),
`zahl-in-traeger-driftet-gegen-die-messung` (29×), `negativtest-ohne-bindung-an-seine-eingabe`
(23×), `beleg-befehl-traegt-seinen-satz-nicht` (21×), `nachzug-laesst-ueberholten-text-stehen`
(18×), `dod-begruendung-unzutreffende-tatsachenbehauptung` (13×),
`adr-aussage-breiter-als-ihre-messung` (11×), `kommentar-behauptet-nicht-getragenen-fehlerpfad`
(9×), `inplace-textwerkzeug-am-repo-trotz-nutzerregel` (9×),
`vorher-nachher-sprache-in-test-harness-kommentar` (8×).

**Sechs Einträge bei 3× oder darüber ohne Ausgang zum Zeitpunkt der Closure.** Die Ausgänge setzen
Regel-Wortlaut in Agenten-Dateien (`.claude/commands/`, `.harness/skills/`) oder in
[`AGENTS.md`](../../../../AGENTS.md); die Closure schrieb dort nichts, weil das Regeländerungen sind
und ein Planner, der sie selbst setzte, dieselbe Rolle wäre, die den eigenen Schnitt prüft
(Baseline-Regelwerk `modul-08-agentenrollen.md`). Die Ausgänge setzte das
[Architect-Verdikt](../../../reviews/architect-verdict-welle-routing-lese-schritt.md) (Auflage F-2 des
Closure-Note-Reviews); die Tabelle gibt sie wieder, die Register-`state.md` der Einträge tragen sie
(Kennung und Anker).

| Eintrag (`BEO-PGC/…`) | Zähler | Ausgang (Verdikt §0 und §3) | Kennung · Anker |
|---|---|---|---|
| `fixrunde-ohne-reviewer-lesung` | 3 (fünf weitere Gegenbelege ohne Datei) | **geplant**, **bedingt auf Freigabe V1**: Regel in der engeren Fassung (*Re-Review verlangen, wenn die Fixrunde Produktionslogik oder eine Norm ändert, oder wenn nach der Fixrunde kein anderer Kontext sie ausgeführt hat*), kein Sensor. Bei Ablehnung: *gestrichen* (akzeptiertes Negativ) | `slice-harness-lese-schritt-regeln-routing` — der Slice wird erst nach der Freigabe angelegt |
| `test-runner-stiller-ausschluss` | 3 | **geplant** (Sensor ohne Regeländerung): jede `func TestE2E*` steht in einem `-run`-Argument von `run-integration-tests.sh`; der Eintrag wird mit dem Slice *verkörpert* | `slice-harness-integration-runner-vollstaendigkeit` (Datei in `open/`) |
| `drei-sprachen-kopie-divergiert-am-randfall` | 3 | **gestrichen** (akzeptiertes Negativ): alle drei Belege LOW und vor dem Merge gefunden, Schaden klein, die Gegenmaßnahmen stehen in Code und Plan; Wiederaufnahme bei einem Divergenz-Fund nach dem Merge oder bei Schwere ab MEDIUM | Begründung in der `state.md` (Verdikt §3.3) |
| `zwei-quellen-drift-handbuch-gegen-pflichtenheft` | 3 | **verkörpert**, steht bereits; die Ergänzung des Beispiels wird nicht beauftragt | `.harness/skills/reviewer.md`, Punkt „Zwei-Quellen-Drift“ (Zeile 259); Anker `seit welle-routing` |
| `plan-zusage-erfuellung-ohne-committeten-anker` | 5 | **geplant**, **bedingt auf Freigabe V2**: ein Satz in `.claude/commands/implement-slice.md` Schritt 18. Bei Ablehnung: *gestrichen* (akzeptiertes Negativ), Trigger „ein Haken ohne Anker bis in `done/`“ | `slice-harness-lese-schritt-regeln-routing` — der Slice wird erst nach der Freigabe angelegt |
| `ein-instanz-annahme-ohne-erzwingung` | 3 | **verkörpert**, steht bereits | [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) §Re-Evaluierungs-Trigger; Anker `seit slice-backfill-run-store` |

**Unter der Schwelle, in dieser Closure nicht gelesen** (Adresse: Sichtungs-Schritt der nächsten
Slice-Planung, `plan-welle`): `suppression-ohne-linter-in-testcode` (1×, ein `# noqa` im
Python-Testcode; [`AGENTS.md`](../../../../AGENTS.md) §3.2 nennt allein `//nolint`),
`integrationsprojekt-uebersetzt-nicht-unbemerkt` (1×), `test-strenger-als-die-zusage` (1×, Adresse
`slice-capture-retry-realtest-lieferzahl-lockern`), `docker-cache-ueberspringt-tests-still` (1×),
`alt-tag-lauf-vorbedingung-am-juengsten-tag` (2×), `spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`
(2×), `kein-admin-weg-schema-fehler-recovery` (2×), `implementierung-weicht-von-adr-wortlaut-ab`
(2×, in der Welle kein Auftreten: die Abweichungen waren als Fragen geführt).

**Feststellung:** Diese Closure hat die sechs Einträge gelesen und je einen Vorschlag mit Adresse
festgehalten, aber keinem eine Kennung zugewiesen; das war die Auflage F-2 des
Closure-Note-Reviews (Modul 6: kein Eintrag überlebt eine Closure ohne Ausgang). Die Ausgänge setzte
anschließend das [Architect-Verdikt](../../../reviews/architect-verdict-welle-routing-lese-schritt.md)
(Tabelle oben); zwei davon (`fixrunde-ohne-reviewer-lesung`, `plan-zusage-erfuellung-ohne-committeten-anker`)
sind bis zur Freigabe des Auftraggebers *geplant*, nicht *verkörpert* (Offene Punkte). Alle übrigen
Einträge bei 3× oder darüber, die die Welle berührt, tragen einen Ausgang.

## Validator-Feststellung (Modul 8)

Die Welle liefert Endnutzer-Wert (ein Betreiber wählt Changes nach Zustellziel,
[`LH-FA-CFG-008`](../../../../spec/lastenheft.md)); der Validator-Schritt ist deshalb nicht „n/a“.
Die Notizen aller zehn Slices schreiben „entfällt — der Nutzer-Bedarf wird erst durch den
Wellen-Beleg validierbar“; dieser Abschnitt ist ihr Träger und gilt für alle zehn Slices (deren
Notizen sind Records).

- **Wer müsste validieren:** die Rolle `validator`
  ([`.claude/agents/validator.md`](../../../../.claude/agents/validator.md)) in frischem Kontext,
  nicht der Planner.
- **Was ist validierbar:** der Bedarf „der Betreiber kann erfasste Changes nach Ziel auswählen“
  gegen [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) in Lastenheft 0.14.0 — Happy Path
  (geordnete Regel, Label, Auswahl über die Lesewege), Boundary (zwei treffende Regeln, `order`;
  Verletzungen R1 bis R6 mit Klartext) und Negative (nicht anwendbare Regel beendet den
  Erfassungspfad mit Klasse `schema`, Abhilfe) am laufenden Feed-Container, über Handbuch und
  SDK-Packages als Betreiber-Oberfläche.
- **Stand:** der Validator-Lauf hat **nicht stattgefunden** und wird hier nicht nachgeholt (der
  Planner darf ihn nicht selbst fahren, kein Self-Review). Die Beleglage in dieser Notiz (e2e-Lauf
  `36933941965`, Verifikations-Reports der Slices) ist eine Belegsammlung des Planners und ersetzt
  keinen Validator; es wird kein Ergebnis behauptet.
- **Adresse:** Auftrag an den `validator`-Agenten nach dieser Closure; Frist: vor dem
  Release-Entscheid (Zeile „Release-Entscheid“ in „Offene Punkte“).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/)
(`BEO-PGC/<slug>/evidence/`); er wird nicht in dieser Notiz gepflegt. Was in dieser Welle
3× erreicht hat, steht oben im Lese-Schritt.

## Folge-Slices

- `slice-capture-retry-realtest-lieferzahl-lockern` — Datei in `docs/plan/planning/open/`;
  lockert `TestRunStreamWithRetrySlotStillActive` von „genau eine Lieferung“ auf „Retry-Change genau
  einmal persistiert“ (At-least-once), Adresse von `test-strenger-als-die-zusage`; aus
  `slice-routing-e2e` F-2, kein Gegenstand der Welle.
- `slice-sdk-sse-client-schema-table-filter` — Datei in `docs/plan/planning/open/`;
  die SSE-Clients der drei Packages und die SSE-Beispiele setzen `schema`/`table` noch nicht
  (der Server filtert seit [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md)); aus
  `slice-routing-sdk-beispiel-target`.

Die Entscheidungen des Lese-Schritts stehen im
[Architect-Verdikt](../../../reviews/architect-verdict-welle-routing-lese-schritt.md) (§5): ein Slice
liegt in `open/` (`slice-harness-integration-runner-vollstaendigkeit`, Ausgang für
`test-runner-stiller-ausschluss`); `slice-harness-lese-schritt-regeln-routing` (V1, V2, V3) entsteht
erst nach der Freigabe des Auftraggebers (Offene Punkte).

## Offene Punkte (mit Adresse)

| Punkt | Zustand | Adresse · Trigger |
|---|---|---|
| Lesekosten der zwei Regelstände im Backfill-Run | **nicht gemessen**; die Verdopplung der Lesungen ist *hergeleitet*, die Prozentzahlen des Architect-Verdikts der Transformations-Welle (≈8 % bei 10.000, ≈70 % bei 100.000 Zeilen der Queue) sind *übernommen* | Messträger `make bench` (`tools/bench-backfill.sh`); Trigger: Druck an der Blockdauer eines Runs bei aktiver Routing-Regel oder die erste Backfill-Messung nach einer Routing-Regel-Last |
| Last der Auswertung je Change (lineare Suche über die Regeln der Tabelle) | **nicht gemessen**, *erwartet* klein | Messträger ein Go-Benchmark im Paket `mapper` bei zehn und bei hundert Regeln; Trigger: erste beobachtete Verzögerung der Erfassung bei vielen Regeln je Tabelle oder Betriebsanforderung über zehn Regeln |
| Abhilfe am Fall (a) „Bedingungsspalte entfernt“ | *hergeleitet* (die Spec führt sie so), nicht gefahren; allgemeiner Recovery-Weg für Schema-Fehler | `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (2×), Out-of-Scope der Welle |
| Server ohne den Parameter `target` | *hergeleitet*; die READMEs sagen es mit dem Ursprung | kein Träger (der Realserver-Slice läuft gegen den Stand des Repos); Beleg wäre ein Lauf gegen ein Release ohne den Parameter |
| Validator-Lauf zum Bedarf (`LH-FA-CFG-008`) | **nicht gelaufen**, Feststellung im Abschnitt „Validator-Feststellung (Modul 8)“ | Auftrag an den `validator`-Agenten nach dieser Closure; Frist: vor dem Release-Entscheid |
| Regeländerung V1 (Re-Review nach einer Fixrunde): `.claude/commands/implement-slice.md` Schritt 21, neuer Absatz | **Freigabe des Auftraggebers ausstehend** (Empfehlung des Architects: zustimmen); ohne Zustimmung Ausgang von `fixrunde-ohne-reviewer-lesung` *gestrichen* | Auftraggeber; Wortlaut und Folge in §4 V1 des [Architect-Verdikts](../../../reviews/architect-verdict-welle-routing-lese-schritt.md); danach Slice `slice-harness-lese-schritt-regeln-routing` |
| Regeländerung V2 (Zusage-Anker beim Abhaken): `.claude/commands/implement-slice.md` Schritt 18, ein Satz am Absatz „DoD-Checkbox-Nachzug im selben Lauf“ | **Freigabe des Auftraggebers ausstehend** (Empfehlung: zustimmen); ohne Zustimmung Ausgang von `plan-zusage-erfuellung-ohne-committeten-anker` *gestrichen*, Trigger „ein Haken ohne Anker bis in `done/`“ | Auftraggeber; §4 V2 des Verdikts; danach derselbe Slice |
| Regeländerung V3 (Validator-Feststellung als Pflichtabschnitt der Results-Notiz): `.claude/commands/close-welle.md` Schritt 3, dazu ein Halbsatz in `implement-slice.md` Schritt 23; zugleich Pflege des veralteten Skill-Kopfs `.harness/skills/closure-note-reviewer.md` (F-8: nennt die fünfte `structure`-Regel „auskommentiert“, sie ist aktiv) | **Freigabe des Auftraggebers ausstehend** (Empfehlung: Variante a, kein Sensor; dritte Wiederholung derselben Klasse „Validator-Schritt ohne Träger“: Backfill-Welle, Transformations-Welle, diese Welle); ohne Zustimmung gilt im Einzelfall „akzeptiert, Review fängt es“ | Auftraggeber; §4 V3 und §6 des Verdikts; Quelle [`closure-note-review-welle-routing.md`](../../../reviews/closure-note-review-welle-routing.md) F-1 und F-8; danach derselbe Slice |
| Release-Entscheid | offen | Auftraggeber: Server-Release, SDK-Versionen; `target` ist in C# neuer letzter Parameter, in Kotlin neuer letzter Parameter mit Default `null` — quellkompatibel, **nicht** binärkompatibel; ein Java-Aufrufer von `PgChangeFeedSseClient.streamChanges` übersetzt nicht mehr ohne Argument |
| `AGENTS.md` §3.2 und Python-`# noqa` | Entscheidung offen | Auftraggeber (Regeländerung); Register `suppression-ohne-linter-in-testcode` (1×) |
| Übersetzen der Integrationsprojekte | Lücke im Sensor-Satz | `BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt` (1×); Aufnahme in `make gates` oder einen Workflow wäre eine Entscheidung (Netzbezug, [`AGENTS.md`](../../../../AGENTS.md) §3.10) |
| Lastenheft-Lesart „Zustellung“ (A-3) | offen, beim Auftraggeber | Review `review-slice-routing-spec-nachzug` F-5; Lastenheft 0.14.0 trägt die Präzisierung des Happy-Path-Kriteriums |

## Verifikation

Alle Zeilen am Stand `84f5b8d2`, sofern nicht anders angegeben; Gate-Exit-Codes ungefiltert gesichert
([`AGENTS.md`](../../../../AGENTS.md) §3.9). Ursprung je Zahl: **gemessen** in dieser Closure,
**übernommen** aus dem genannten Report (nicht nachgemessen), **abgeleitet** (gerechnet).

### Schritt 1 — Trigger

| Kriterium (Welle-Datei §3) | Beleg |
|---|---|
| Alle zehn Slices in `done/` | `ls docs/plan/planning/done \| grep -c '^slice-routing-'` druckt `10` (gemessen) |
| `make gates` grün | Exit-Code und Zeilen im Abschnitt „Closure-Lauf“ unten (gemessen) |
| Realer, grüner `make test-integration`-Lauf mit den Belegen | `e2e`-Lauf `36933941965` am Commit `84f5b8d2341946378da9af35c0a2d01633a278ae`, beide Legs `success` (`gh run list --commit 84f5b8d2341946378da9af35c0a2d01633a278ae`, gemessen 2026-10-02): `image + test-integration (PostgreSQL 17)` und `(PostgreSQL 18)`. Jedes Leg druckt die vier Routing-Zeilen „Routing-Happy-Path“, „Routing-Neustart und Ausschluss-Sperre“, „Routing-Aktivierung über die API“, „Routing-Nichtanwendbarkeit und Abhilfe“ und die vier Go-Testfunktionen `TestE2ERoutingSelectsTargetsOverChangesView`, `TestE2ERoutingConflictsFailWithSpecText`, `TestE2ERoutingReplayKeepsLabelAndBackfillCarriesIt`, `TestE2ERoutingDeleteWithoutFullReplicaIdentity` (`--- PASS`); der Lauf endet mit „Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen“ (gemessen im Log des Laufs) |
| DELETE-Messung an PostgreSQL 17 und 18 | gedruckt je Leg: `ROUTING-DELETE-MESSUNG PostgreSQL 17.11` bzw. `PostgreSQL 18.6`, Replica-Identität DEFAULT mit Regeln eu+rest: `INSERT="eu" UPDATE="eu" DELETE="sonstige"`, Alt-Bild des DELETE `{"id": "1"}`; DEFAULT mit Regel eu: `DELETE=""` (gesetzt false); FULL mit eu+rest: `DELETE="eu"`, Alt-Bild `{"id": "1", "name": "Ada2", "region": "eu"}` — beide Versionen liefern dieselben Ziele; `ADR-0137` Entscheidung 4 bestätigt (gemessen im Log von `36933941965`) |
| Laufzeit gegen das 60-Minuten-Limit | Leg PostgreSQL 18: Job 22:15:51 bis 22:35:45 = 19 min 54 s, Schritt „Compose-Integrationstest“ 22:16:26 bis 22:31:32 = 15 min 6 s; Leg PostgreSQL 17: Job 22:15:47 bis 22:40:56 = 25 min 9 s, Schritt „Compose-Integrationstest“ 22:16:30 bis 22:35:32 = 19 min 2 s (Zeitstempel aus `gh run view 36933941965 --json jobs` und dem Log, Differenzen **abgeleitet**); Limit `timeout-minutes: 60` in `.github/workflows/e2e.yml`; kein neuer Workflow, kein Workflow-Zug in der Welle ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht) |
| CI-Stand des Endstands | `ci` `36933941957` (Job `gates + test`), `e2e` `36933941965`, `examples` `36933942004`: alle `success` (gemessen; `gates + test` druckt `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK`, `gesamt: 0 Befund(e)`) |
| Reale, grüne Läufe der DB-Tiers | `e2e`-Lauf `36933941965`, beide Legs: die Schritte `make test-store`, `run-replication-tests.sh measure` und `run-replication-tests.sh tier` `success`; gedruckt `DB-Adapter-Coverage: 82.99% (gedeckt 961 von 1158 Statements; Profile gemergt: store,replication)` und `db-coverage: OK — … erfuellt Schwelle 80%` (je Leg, gemessen); `sourcekeepalive_test.go`-Zeilen nennen `PostgreSQL 17.11` und `PostgreSQL 18.6`. Die Läufe der Slices stehen in den Verifikations-Reports (`kern-label`, `antragsweg`, `backfill-pfad`, `lesewege`, **übernommen**) |
| `make schema-rollout` zweimal und Alt-Tag-Lauf | `bash tools/harness/run-schema-rollout-guard-test.sh`, Exit `0` (gemessen in dieser Closure, Arbeitsbaum `84f5b8d2`; im Nachzug zum Closure-Note-Review ([`closure-note-review-welle-routing.md`](../../../reviews/closure-note-review-welle-routing.md), dort reproduziert) am Stand `cff48b65` erneut gefahren: Exit `0`, dieselbe Zeile „Lauf 5 OK — Tag v0.4.0 …“ und dieselbe Schlusszeile „OK — alle Belege real erbracht“): „Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar, Soll-Signatur der View, route_target NULL (Spalte im Stand des Tags nicht vorhanden); … EXECUTE auf 5 Funktionen (… cdc.set_route(text, text, text, text, json), cdc.remove_route(text, text, text, text)) allein für cdc_admin (nicht PUBLIC)“ und „OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag v0.4.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)“ |
| Typ-Satz-Parität an PostgreSQL 17 (Ausgang aus `slice-routing-e2e` §6, vorher **übernommen**) | `e2e`-Lauf, Leg PostgreSQL 17: `ok … internal/adapters/driven/postgressnapshot 6.724s` bzw. `5.664s` (nicht übersprungen: ein Überspringen ohne DSN dauert 0,004 s, dieselbe Zeile im Log); der Aufruf von `checkRouteParity` liegt in einem Test dieses Pakets (`snapshot_test.go:741`); dass der Test mit Version 17 lief, folgt aus `PG_TEST_IMAGE` des Legs — *hergeleitet*, die Zeile nennt die Version nicht |
| Fitness-Function-Tests aus `ADR-0137` | `make test` und `make a-check`/`make generated-sync` laufen in `make gates` bzw. im `ci`-Lauf `36933941957` (`gates + test`, `success`); die Einzel-Tests (Determinismus, `ErrRoutingNotApplicable`, Eigenschaftstest R3 beide Richtungen, Store-Test) stehen in den Verifikations-Reports der Slices (**übernommen**) |
| `make doc-trace` | Exit `0` (gemessen in dieser Closure, Arbeitsbaum); gedruckt `\| LH-FA-CFG-008 \| Routing von Changes auf Zustellziele \| ADR-0112, ADR-0137, ADR-0138, ADR-0139, ADR-0140, ADR-0141 \| — \| E2E, SDK-E2E \| ok \|` und `80 Anforderung(en), 0 Waise(n).`; Parent `30fd6cb5` (`git archive 30fd6cb5` im Scratchpad, dort `make doc-trace`, Exit `0`, **gemessen** im Nachzug zum Closure-Note-Review): `\| LH-FA-CFG-008 \| … \| ADR-0112, ADR-0137 \| — \| — \| WAISE \|` und `80 Anforderung(en), 1 Waise(n).` |
| Auswertung an genau einer Stelle | `git grep -n 'EvaluateRoute(' -- 'internal/*.go' ':!*_test.go'` druckt drei Zeilen: `internal/domain/model/route.go:163` (Deklaration), `internal/adapters/driving/replication/mapper/mapper.go:325` (WAL-Pfad), `internal/application/usecase/backfill/service.go:634` (Backfill-Pfad) — keine zweite Konstruktionsstelle (gemessen am Arbeitsbaum) |
| Handbuch-Träger | Änderungshistorie 1.84 bis 1.86; `make sdk-public-doc-check` ist Teil von `make gates` (Exit unten) |
| Drei SDK-Realserver-Tiers | `make test-sdk-csharp-integration`, `…-kotlin-…`, `…-python-…` je Exit 0 mit zwölf Phasen und `foreign=0` in allen zwölf Routing-Phasen, Zustellung 97 bis 388 ms (Verifikations-Report `verifikation-slice-routing-sdk-realserver-e2e` §1 und §3, **übernommen**); `git diff --name-only db783913 HEAD` nennt außerhalb von `docs/` keine Datei (gemessen) — die Tiers laufen gegen unveränderten Code |
| Lese-Schritt des Beobachtungs-Registers | oben, Abschnitt „Lese-Schritt“ — gelaufen |
| Closure-Notiz | diese Datei |

### Schritt 2 — Trigger-Audit

- **Carveouts: 0 offen.** `find docs -iname "CO-*.md"` liefert keinen Treffer (gemessen); kein
  rotes Gate erreicht `done/`.
- **Bootstrap-aware Gates.** *Unit-Gate:* `harness/mk/coverage.mk` führt `THRESHOLD ?= 80`
  (Endstufe, Rampe ausgeschöpft); `ci`-Lauf `36933941957` druckt `Coverage 82.00%`. *DB-Adapter-Coverage:*
  Schwelle 80 (Endstufe), `e2e`-Lauf druckt `82.99%` je Leg. **0 offen.**
- **ADR-Re-Evaluierungs-Trigger** (gelesen: der Abschnitt jeder der fünf ADRs):
  - [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md): Entscheidungen
    1, 2, 5 unverändert; kein Operator über Gleichheit hinaus im Register verlangt und kein Bedarf
    an Umetikettieren gemeldet: `git grep -n -i 'umetikett\|relabel\|Vergleichsoperator' -- docs/plan/planning/observations`
    druckt `0` Treffer (gemessen am Stand `cff48b65`; ein Suchbegriff-Negativbefund, keine Aussage über
    Anliegen in anderer Wortwahl); die Messung der zweiten
    NATS-Veröffentlichung löst keinen Aufschlag am Publisher auf (Verhältnis mit/ohne 0,92 bis
    1,30 über fünf Läufe, **abgeleitet**, `slice-routing-nats-subjekt` §6 und §7; der
    Trigger „Messung zeigt Druck am Publisher“ ist nicht eingetreten). **0 offen.**
  - [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md):
    die Nichtanwendbarkeit ist am System erreichbar (Messung `slice-routing-e2e`, (b) und (c)) —
    Trigger „nie erreichbar“ nicht eingetreten; `ReadChangesRequest` trägt keinen weiteren
    Filter: die Nachricht in `proto/cdc/administration/v1/administration.proto` hat sieben Felder
    (`awk '/message ReadChangesRequest/,/^}/' … | grep -c '= [0-9]*;'` druckt `7` am Arbeitsbaum, `6` am
    Parent `30fd6cb5`), das neue siebte ist `string target = 7;` (gemessen). **0 offen.**
  - [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
    [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md): kein Run-Fehlerbild
    mit Konfigurationsbezug des Lesefehlers, keine neue Fehlerform für Filterwerte. Belegt ist nur die
    Hälfte „keine neue Klasse in `SPEC-008`“: `git grep -n -i 'nicht prüfbar' -- spec/pflichtenheft.md`
    druckt `0` Treffer (gemessen; Suchbegriff-Negativbefund). Für „kein Run-Fehlerbild mit
    Konfigurationsbezug“ und „keine Fehlerform für Filterwerte“ ist **kein Suchbefehl gelaufen**, der Satz
    stützt sich auf das Lesen der Verifikations-Reports der Slices `backfill-pfad` und `lesewege`
    (**übernommen**, **nicht belegt** als Suche). **0 offen** gilt unter diesem Vorbehalt.
  - [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md):
    Trigger 2 („Messung erzeugt weder (b) noch (c)“) ist **nicht** eingetreten — beide
    sind am System erzeugt, die Verengung greift nicht; Trigger 1 (Auftraggeber will, dass eine
    entfernte, von keiner Regel genannte Spalte den Pfad nicht beendet): in Reports und Register
    kein Hinweis darauf — **nicht belegt** durch einen Suchbefehl, Aussage des Planners nach dem
    Lesen der Reports. **0 offen.**
- **Beobachtungs-Register:** Lese-Schritt oben — sechs Einträge bei 3× oder darüber, Ausgänge durch
  das Architect-Verdikt gesetzt, zwei davon bedingt auf die Freigabe des Auftraggebers.

### Schritt 4 — Archivierung

Das Repo führt kein eigenes Archivierungs-Werkzeug (kein Make-Ziel in `Makefile` oder
`harness/mk/*.mk`, kein Skript unter `tools/harness/` mit Archiv-Namen, gemessen); das externe
Werkzeug `ai-harness-init archive-welle` liegt lokal im PATH, committet aber mit fest
einprogrammierten Messages ohne `LH-*`/`ADR-*`-Kennung
(`BEO-PGC/externes-werkzeug-committet-ohne-kennung`) — dieselbe Bedingung wie bei den
Vorgänger-Closures. Die Bedingung des Schrittes (ein im Repo geführtes Werkzeug) ist nicht
eingetreten; keine Handarchivierung. Ob der Lauf des externen Werkzeugs gleichwohl ausgeführt
wird, entscheidet der Betreiber.

### fünfte `structure`-Regel

Die Regel (Closure-Notiz-Gate auf `done/slice-*.md`) ist in `.d-check.yml` seit der ersten Closure
aktiv (Kommentar „AKTIVIERT mit der ERSTEN Closure“ im Abschnitt `structure`); diese Closure
aktiviert nichts neu. Die zehn Slice-Pläne dieser Welle liegen unter dieser Regel und
`make docs-check` ist grün (Abschnitt „Closure-Lauf“).

### Drei Paarungen

- **(a) Anker:** kein Steering-Loop-Eintrag dieser Welle trägt das Feld `liegt in` (alle zehn
  Lerneinträge sind geschärfte Handlungen mit Träger „Lese-Handlung einer Rolle“; die Ausgänge
  des Lese-Schritts sind Vorschläge, keine verkörperten Regeln). Nichts offen; geprüft mit
  `grep -n "liegt in" docs/plan/planning/done/slice-routing-*.md`: die Treffer sind
  Lifecycle-Prosa („liegt in `done/`“), kein Feld.
- **(b) Folge-Slice:** beide genannten Folge-Slices liegen als Datei in
  `docs/plan/planning/open/` (`ls` gemessen); die übrigen in den Closure-Notizen genannten
  Slices liegen in `done/`.
- **(c) Register:** alle 55 in den Plänen und der Welle-Datei genannten Kennungen
  `BEO-PGC/<slug>` existieren als Verzeichnis und tragen ein nicht leeres `evidence/` (Schleife
  über die aus den elf Dateien gezogene Liste, gemessen: keine Ausgabe). Zwei Verzeichnisse des
  Registers tragen **kein** `evidence/` (`architect-verdikt-ablageort-uneinheitlich`, gestrichen
  seit 2026-09-13; `limit-fortsetzung-innerhalb-einer-position`, offen seit 2026-09-23) — beide
  sind keine Kennung dieser Welle und stehen hier als Befund mit Adresse: Sichtungs-Schritt der
  nächsten Slice-Planung.
- Die DoD-Zeile „Die drei Paarungen“ der zehn Slice-Pläne ist mit dieser Closure abgehakt
  (Beleg-Anker: dieser Abschnitt).

### Closure-Lauf

Am Arbeitsbaum vor dem Self-Close-Commit (gemessen, 2026-10-02, Exit-Codes ungefiltert in
eigenen Aufrufen gesichert): `make docs-check` Exit `0` (`d-check: 1535 Datei(en) geprüft,
0 Befund(e)`); `make gates` Exit `0` (`baseline-verify: v6.13.0 OK — 54 Dateien`,
`coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK`,
`sdk-public-doc-check: keine interne Kennung unter sdks`, `gesamt: 0 Befund(e)` der
Architekturprüfung); `make doc-trace` Exit `0`. Der Gate-Stempel wird mit jedem Commit ungültig;
der Lauf nach den Move-Commits bestätigt ihn neu (Bericht der Closure).

### Docker-Lauf-Hygiene

Die Docker-Läufe dieser Closure (`make doc-trace`, der Schema-Rollout-Guard-Test,
`make gates`, `make docs-check`) liefen nacheinander, nicht gleichzeitig; jeder endete, bevor
der nächste startete.

### Ausgang des Closure-Note-Reviews

Zum Zeitpunkt des Self-Close lag **kein** Review-Lauf dieser Closure-Notiz in frischem Kontext
vor (kein Self-Review, Modul 8). Er ist danach gelaufen:
[`closure-note-review-welle-routing.md`](../../../reviews/closure-note-review-welle-routing.md)
(`.harness/skills/closure-note-reviewer.md`, kein HIGH, zwei MEDIUM, fünf LOW, ein INFO). Nachgezogen
in dieser Notiz: F-1 (Validator-Feststellung), F-3 (Zuschreibung), F-4 (Ursprungs-Etiketten), F-5
(Suchbefehle im Trigger-Audit); in den zehn Slice-Plänen: F-6 (Zustandsprosa). Offen mit
Adresse beim Architect: F-2 (Ausgänge des Lese-Schritts mit Kennung), F-7 (Zitat-Korrektur an
[`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)),
Pflege-Regel zu F-1. Das [Architect-Verdikt](../../../reviews/architect-verdict-welle-routing-lese-schritt.md)
beantwortet sie: F-2 durch die Ausgänge im Lese-Schritt, F-7 mit „stehen lassen“ (§3.7), die
Pflege-Regel zu F-1 als Vorlage V3 (Freigabe des Auftraggebers ausstehend). F-8 (INFO) ist ein
Hinweis: der Kopftext von `.harness/skills/closure-note-reviewer.md` nennt die fünfte
`structure`-Regel noch „auskommentiert“, sie ist seit der ersten Closure aktiv; die Pflege liegt beim
Besitzer der Skill-Datei und gehört zur Freigabe V3.

**Erwartbare Abweichung in den Slice-Plänen:** die `suchlauf`-Zeilen am Stand `diff` in
`slice-routing-kern-label` (4 von 12) und `slice-routing-e2e` (1 von 14) weichen bei
`make suchlauf-nachmessen` ab, weil spätere Slices dieselben Träger bewegten (`diff` bewegt sich
mit jedem Commit, Sensor-Vertrag); sie bleiben unverändert, weil die Pläne geschlossene Records sind.
