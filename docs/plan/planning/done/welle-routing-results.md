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
  (Parent `30fd6cb5`: `80 Anforderung(en), 1 Waise(n)`).

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

**Bereits mit zugewiesenem Ausgang (verkörpert oder gestrichen), in dieser Welle gewachsen** —
kein Zug nötig, der Zähler steht in der Beobachtung: `arbeit-ueberholt-stehenden-traeger` (34×),
`zahl-in-traeger-driftet-gegen-die-messung` (29×), `negativtest-ohne-bindung-an-seine-eingabe`
(23×), `beleg-befehl-traegt-seinen-satz-nicht` (21×), `nachzug-laesst-ueberholten-text-stehen`
(18×), `dod-begruendung-unzutreffende-tatsachenbehauptung` (13×),
`adr-aussage-breiter-als-ihre-messung` (11×), `kommentar-behauptet-nicht-getragenen-fehlerpfad`
(9×), `inplace-textwerkzeug-am-repo-trotz-nutzerregel` (9×),
`vorher-nachher-sprache-in-test-harness-kommentar` (8×).

**Vier Einträge bei 3× ohne Ausgang, den diese Closure zuweisen kann**, und zwei Einträge,
die der Lese-Schritt der Vorgängerwelle schon adressiert hat. Die Ausgänge setzen
Regel-Wortlaut in Agenten-Dateien (`.claude/commands/`, `.harness/skills/`) oder in
[`AGENTS.md`](../../../../AGENTS.md); diese Closure schreibt dort nichts, weil das Regeländerungen sind
und ein Planner, der sie selbst setzte, dieselbe Rolle wäre, die den eigenen Schnitt prüft
(Baseline-Regelwerk `modul-08-agentenrollen.md`). Jede Zeile trägt einen begründeten
Ausgang-Vorschlag und die Adresse der Entscheidung.

| Eintrag (`BEO-PGC/…`) | Zähler | Befund beim Lesen | Ausgang-Vorschlag · Adresse |
|---|---|---|---|
| `fixrunde-ohne-reviewer-lesung` | 3 (fünf weitere Gegenbelege ohne Datei) | Die drei Auftreten: `spec-nachzug` (Fixrunde am Spec-Text; der Verifier fand danach V-1/V-2), `kern-label` (Kommentare und Skript-Kopftext), `antragsweg` (Parser; das nachgeholte Re-Review fand F-N1). Die Gegenbelege zeigen, wo ein Re-Review Befunde lieferte (Fixrunde mit geänderter Anweisung: `antragsweg`, `backfill-pfad`, `lesewege`) und wo er entbehrlich war (der Verifier hat die Fixrunde ausgeführt: `nats-subjekt`, `e2e`, `betriebsdoku`, `sdk-beispiel-target`, `sdk-realserver-e2e`). Der Kandidat-Wortlaut „sobald die Fixrunde Anweisungen ändert“ ist zu weit (ein Skript mit Exit-Code ist eine Anweisung, die ein zweiter Kontext ausführte). | **Verkörpern**, engere Fassung: *Re-Review verlangen, wenn die Fixrunde Produktionslogik oder eine Norm ändert, oder wenn nach der Fixrunde kein anderer Kontext sie ausgeführt hat.* Träger: eine Zeile im Handoff nach der Fixrunde in `.claude/commands/implement-slice.md` und der Prüfpunkt des Verifiers; kein Sensor (ob eine Fixrunde Logik ändert, ist eine Lese-Frage). Architect-Zug, danach Planner/Implementer schreibt die Zeile. |
| `test-runner-stiller-ausschluss` | 3 | Drittes Auftreten (`nats-subjekt` F-3) ist in seinem Skript **geschlossen** (`run-notify-tests.sh` fährt das ganze Paket und endet bei `--- SKIP` mit Exit 1, der Verifier fuhr es in drei Zuständen). Für `run-integration-tests.sh` bleibt die Lücke: ein `-run`-Muster (Zeile 466 und folgende) erfasst nur die Funktionen, die es nennt; die Deklarations-Hälfte deckt `TestAbdeckungstabelleZeilen`, die Vollständigkeits-Hälfte kein Wächter. | **Bleibt offen**, Teil geschlossen. Entscheidung, die aussteht: ob ein Vollständigkeits-Sensor gebaut wird (jede `func TestE2E*` wird von mindestens einem `-run`-Muster oder einer Runner-Phase erfasst). Architect-Zug; Sensor wäre ein eigener Slice. Trigger bis dahin: das nächste Auftreten bei `run-integration-tests.sh`. |
| `drei-sprachen-kopie-divergiert-am-randfall` | 3 | Das dritte Auftreten ist ein Satz der Python-README, enger als der Code; das Verhalten (leeres Ziel am NATS-Stream) ist eine gewollte API-Form-Differenz. Wirksam waren zwei Gegenmaßnahmen: derselbe Pflicht-Eingabesatz je Sprache im Plan und eine Fixture-Quelle für Testvorbereitung; sie stehen nicht als Regel. | **Verkörpern** als Plan-Pflicht für Drei-Sprachen-Slices: der Plan nennt die Eingabetabelle (hier `null`/`""`/`eu`/`a&b=c`) und eine README-Aussage über einen Randfall wird am Code, nicht an der Schwester-README gelesen. Träger: `.claude/commands/plan-welle.md` (Pflichtangabe) und Punkt „Form-Vorbild“ in `.harness/skills/reviewer.md`. Architect-Zug (Alternative: akzeptiertes Negativ, weil jedes Auftreten vor dem Merge gefunden wurde). |
| `zwei-quellen-drift-handbuch-gegen-pflichtenheft` | 3 | Seit der Vorgänger-Closure kein neues Auftreten; in dieser Welle beschrieben Pflichtenheft, Handbuch und drei SDK-READMEs dieselben Parameter und Fehlertexte, und die Reviews fanden keinen Drift der Klasse (Spec-Nachzug: „kein neues Auftreten“). Was wirkte: das Handbuch wiederholt Wortlaut mit Adresse und verweist je Sachverhalt auf die führende Stelle; neun von zehn Fehlertexten hat der Verifier am System gegen den Text gehalten. | **Verkörpern** im bestehenden Punkt „Zwei-Quellen-Drift“ von `.harness/skills/reviewer.md` (Handbuch-gegen-Pflichtenheft als benanntes Beispiel samt der Gegenmaßnahme „Verweis auf die führende Stelle statt Wiederholung“). Architect-Zug oder Auftraggeber. |
| `plan-zusage-erfuellung-ohne-committeten-anker` | 5 | Nicht Gegenstand dieser Welle (kein Auftreten in den zehn Slices); „offen, Neubewertung fällig“ seit `welle-transformationen`. | unverändert: Architect-Zug aussteht. |
| `ein-instanz-annahme-ohne-erzwingung` | 3 | Nicht Gegenstand dieser Welle; das Architect-Verdikt der Vorgänger-Closure steht (Trigger nicht ausgelöst, Ausgang „weiter offen“). | unverändert. |

**Unter der Schwelle, in dieser Closure nicht gelesen** (Adresse: Sichtungs-Schritt der nächsten
Slice-Planung, `plan-welle`): `suppression-ohne-linter-in-testcode` (1×, ein `# noqa` im
Python-Testcode; [`AGENTS.md`](../../../../AGENTS.md) §3.2 nennt allein `//nolint`),
`integrationsprojekt-uebersetzt-nicht-unbemerkt` (1×), `test-strenger-als-die-zusage` (1×, Adresse
`slice-capture-retry-realtest-lieferzahl-lockern`), `docker-cache-ueberspringt-tests-still` (1×),
`alt-tag-lauf-vorbedingung-am-juengsten-tag` (2×), `spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`
(2×), `kein-admin-weg-schema-fehler-recovery` (2×), `implementierung-weicht-von-adr-wortlaut-ab`
(2×, in der Welle kein Auftreten: die Abweichungen waren als Fragen geführt).

**Feststellung:** Diese Closure hat die vier Einträge bei 3× gelesen und je Ausgang-Vorschlag und
Adresse festgehalten; sie hat **keinen** der Ausgänge selbst gesetzt (kein `liegt in`-Feld,
daher keine Anker-Paarung offen). Alle übrigen Einträge bei 3× oder darüber, die die Welle
berührt, tragen einen Ausgang.

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

Die Entscheidungen des Lese-Schritts (Tabelle oben) sind Architect-Fragen, keine Slices; ein
Slice entsteht erst mit einer Entscheidung.

## Offene Punkte (mit Adresse)

| Punkt | Zustand | Adresse · Trigger |
|---|---|---|
| Lesekosten der zwei Regelstände im Backfill-Run | **nicht gemessen**; die Verdopplung der Lesungen ist *hergeleitet*, die Prozentzahlen des Architect-Verdikts der Transformations-Welle (≈8 % bei 10.000, ≈70 % bei 100.000 Zeilen der Queue) sind *übernommen* | Messträger `make bench` (`tools/bench-backfill.sh`); Trigger: Druck an der Blockdauer eines Runs bei aktiver Routing-Regel oder die erste Backfill-Messung nach einer Routing-Regel-Last |
| Last der Auswertung je Change (lineare Suche über die Regeln der Tabelle) | **nicht gemessen**, *erwartet* klein | Messträger ein Go-Benchmark im Paket `mapper` bei zehn und bei hundert Regeln; Trigger: erste beobachtete Verzögerung der Erfassung bei vielen Regeln je Tabelle oder Betriebsanforderung über zehn Regeln |
| Abhilfe am Fall (a) „Bedingungsspalte entfernt“ | *hergeleitet* (die Spec führt sie so), nicht gefahren; allgemeiner Recovery-Weg für Schema-Fehler | `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (2×), Out-of-Scope der Welle |
| Server ohne den Parameter `target` | *hergeleitet*; die READMEs sagen es mit dem Ursprung | kein Träger (der Realserver-Slice läuft gegen den Stand des Repos); Beleg wäre ein Lauf gegen ein Release ohne den Parameter |
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
| `make schema-rollout` zweimal und Alt-Tag-Lauf | `bash tools/harness/run-schema-rollout-guard-test.sh`, Exit `0` (gemessen in dieser Closure, Arbeitsbaum `84f5b8d2`): „Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar, Soll-Signatur der View, route_target NULL (Spalte im Stand des Tags nicht vorhanden); … EXECUTE auf 5 Funktionen (… cdc.set_route(text, text, text, text, json), cdc.remove_route(text, text, text, text)) allein für cdc_admin (nicht PUBLIC)“ und „OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag v0.4.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)“ |
| Typ-Satz-Parität an PostgreSQL 17 (Ausgang aus `slice-routing-e2e` §6, vorher **übernommen**) | `e2e`-Lauf, Leg PostgreSQL 17: `ok … internal/adapters/driven/postgressnapshot 6.724s` bzw. `5.664s` (nicht übersprungen: ein Überspringen ohne DSN dauert 0,004 s, dieselbe Zeile im Log); der Aufruf von `checkRouteParity` liegt in einem Test dieses Pakets (`snapshot_test.go:741`); dass der Test mit Version 17 lief, folgt aus `PG_TEST_IMAGE` des Legs — *hergeleitet*, die Zeile nennt die Version nicht |
| Fitness-Function-Tests aus `ADR-0137` | `make test` und `make a-check`/`make generated-sync` laufen in `make gates` bzw. im `ci`-Lauf `36933941957` (`gates + test`, `success`); die Einzel-Tests (Determinismus, `ErrRoutingNotApplicable`, Eigenschaftstest R3 beide Richtungen, Store-Test) stehen in den Verifikations-Reports der Slices (**übernommen**) |
| `make doc-trace` | Exit `0` (gemessen in dieser Closure); gedruckt `\| LH-FA-CFG-008 \| Routing von Changes auf Zustellziele \| ADR-0112, ADR-0137, ADR-0138, ADR-0139, ADR-0140, ADR-0141 \| — \| E2E, SDK-E2E \| ok \|` und `80 Anforderung(en), 0 Waise(n).` (Parent `30fd6cb5`: `80 Anforderung(en), 1 Waise(n)`) |
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
    1, 2, 5 unverändert; kein Operator über Gleichheit hinaus im Register verlangt (`BEO-PGC`
    nennt keinen Eintrag dazu); kein Bedarf an Umetikettieren gemeldet; die Messung der zweiten
    NATS-Veröffentlichung löst keinen Aufschlag am Publisher auf (Verhältnis mit/ohne 0,92 bis
    1,30 über fünf Läufe, **abgeleitet**, `slice-routing-nats-subjekt` §6 und §7; der
    Trigger „Messung zeigt Druck am Publisher“ ist nicht eingetreten). **0 offen.**
  - [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md):
    die Nichtanwendbarkeit ist am System erreichbar (Messung `slice-routing-e2e`, (b) und (c)) —
    Trigger „nie erreichbar“ nicht eingetreten; `ReadChangesRequest` trägt keinen weiteren
    Filter. **0 offen.**
  - [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
    [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md): kein Run-Fehlerbild
    mit Konfigurationsbezug des Lesefehlers, keine neue Fehlerform für Filterwerte. **0 offen.**
  - [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md):
    Trigger 2 („Messung erzeugt weder (b) noch (c)“) ist **nicht** eingetreten — beide
    sind am System erzeugt, die Verengung greift nicht; Trigger 1 (Auftraggeber will, dass eine
    entfernte, von keiner Regel genannte Spalte den Pfad nicht beendet): in Reports und Register
    kein Hinweis darauf. **0 offen.**
- **Beobachtungs-Register:** Lese-Schritt oben — vier Einträge bei 3× ohne selbstgesetzten
  Ausgang, je mit Vorschlag und Adresse.

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

Ein Review-Lauf dieser Closure-Notiz in frischem Kontext
(`.harness/skills/closure-note-reviewer.md`) liegt zum Zeitpunkt dieses Schreibens **nicht**
vor. Die Abweichung ist benannt, nicht verschwiegen: die Prüfung obliegt der nächsten Rolle im
Workflow, nicht diesem Planner-Zug (kein Self-Review, Modul 8).
