# Welle routing: Routing auf Zustellziele — erfasste Changes tragen bei der Erfassung das durch geordnete Regeln bestimmte Ziel als persistiertes Label, konfiguriert über die SQL-Antrags-Queue und wählbar an jedem Lesezugriffsweg (`LH-FA-CFG-008`, `ADR-0137`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-routing-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (Rolleninhaber der Implementer-Rolle je Slice, gesetzt
beim Übergang `open` → `next`; geschnitten aus
[`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
§Folgepflichten durch den Planner). **Datum:** 2026-10-01.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Eine Tabelle trägt eine geordnete Liste von Routing-Regeln, und **jede** danach
erfasste Change trägt das Ziel der ersten treffenden Regel als Label
`route_target` (leer, wenn keine Regel trifft); der Leser wählt das Ziel an
**jedem** Weg, der Changes liefert: SQL-Lesezugriff `cdc.changes`,
`GET /changes`, gRPC-Stream und SSE-Stream über den Parameter `target`, der
NATS-Vollinhalts-Weg über das zusätzliche Subjekt
`cdc.route.<source_id>.<ziel>`; ein ungefilterter Leser sieht weiterhin alle
Changes; eine auf eine Change nicht anwendbare Regel endet sichtbar statt still.
Das ist die Aussage der drei Akzeptanzkriterien von
[`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Happy Path, Boundary, Negative).
Der Mechanismus ist mit
[`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
entschieden: ein Zustellziel ist ein benannter Kanal je Quelle, das Label wird im
`Assembler` vor der Persistierung bestimmt (ein Log, ein Auswertungsort),
konfiguriert über zwei neue Antragsarten der Antrags-Queue (`cdc.set_route`,
`cdc.remove_route`), Mehrdeutigkeit über `order` geordnet und über sechs
Invarianten R1–R6 statisch ausgeschlossen.

Das *Mehr* gegenüber den neun Slice-DoDs: keiner der Slices belegt eines der drei
Kriterien allein. Der Happy Path braucht die Domäne, den `Assembler`, die
Spalte samt View, den Antragsweg mit Dauerhaftigkeit und alle Lesewege am
laufenden Feed-Container; die Boundary (Reihenfolge entscheidet, R1–R6) liegt
im Use Case und wird erst am komponierten System als `failed` mit Text sichtbar;
die Negative braucht die Prüfung im `Assembler`, die Klassen-Abbildung,
`diagnose`, die Startreihenfolge und die Abhilfe. Dazu kommt der Träger-Zustand,
den kein Einzel-DoD beobachtet: die Auswertung liegt an genau **einer** Stelle
(aufgerufen vom WAL-Pfad **und** vom Backfill-Pfad), die View-Signatur von
`cdc.changes` wächst über einen Alt-Bestand, und
[`LH-FA-CFG-008`](../../../spec/lastenheft.md) ist im RTM-Lauf (`make doc-trace`)
nicht mehr Waise.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  trägt den Status `Accepted` — **bereits erfüllt**, geprüft an der Zeile
  `ADR-0137` im ADR-Index ([`docs/plan/adr/README.md`](../adr/README.md),
  Status-Spalte `Accepted`, Datum 2026-10-01).
- Das Lastenheft trägt [`LH-FA-CFG-008`](../../../spec/lastenheft.md) als eigene
  Anforderung — **bereits erfüllt**, geprüft an `spec/lastenheft.md` (Abschnitt
  `LH-FA-CFG-008`, Boundary-Kriterium nennt die Regel-Reihenfolge).
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden. Die
  Abhängigkeiten der Slices untereinander (§5) sind Start-Trigger einzelner
  Slices, kein Trigger der Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle neun Slices in `done/` (die Zahl ist ein Plan-Stand:
  [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) nennt
  acht Folgepflichten, Folgepflicht 8 ist in `betriebsdoku` und
  `sdk-beispiel-target` geteilt, §4 Abweichung 1).
- `make gates` grün — der Exit-Code des Laufs wird ungefiltert gesichert und
  gesondert ausgewertet ([`AGENTS.md`](../../../AGENTS.md) §3.9).
- **Ein realer, grüner `make test-integration`-Lauf** mit den Belegen am
  laufenden Feed-Container — Happy Path (eine Herkunfts- und eine Inhaltsregel:
  Ziel A sieht nur A, ein ungefilterter Leser sieht alles, über `cdc.changes`,
  `GET /changes`, gRPC, SSE und das NATS-Subjekt `cdc.route…`), Boundary (zwei
  Regeln treffen, `order` entscheidet; je eine Verletzung R1–R6 endet `failed`
  mit Text, der Regelstand bleibt), Neustart-Festigkeit (`docker restart`),
  Ausschluss-Sperre R3 in beiden Richtungen
  ([`LH-QA-SEC-004`](../../../spec/lastenheft.md)), Replay (erneutes Lesen
  derselben Change liefert dasselbe Label), ein Backfill-Bestand mit Label und
  die **Negative** (nicht anwendbare Regel endet sichtbar mit Klasse `schema`).
  Für die Negative gilt die Vorab-Bedingung V3 (§5): Kann die Nichtanwendbarkeit
  am komponierten System nicht erzeugt werden, steht das Verdikt dazu im
  Closure-Bericht; sie wird nicht stillschweigend auf einen Unit-Beleg
  verengt. Kein Gate, aber Pflichtbeleg dieser Welle.
- **Die Messung des DELETE-Verhaltens an beiden PostgreSQL-Versionen**
  ([`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Entscheidung 4, Folgepflicht 7): eine Inhaltsregel auf eine Nicht-Schlüsselspalte
  ohne volle Replica-Identität wirkt bei DELETE nicht; die Aussage gilt erst als
  belegt, wenn ein Lauf gegen PostgreSQL 17 **und** 18 sie trägt
  (`e2e.yml`-Matrix, beide Legs `success`; die Lauf-ID steht im Closure-Bericht).
  Bis dahin steht sie im Handbuch nur als erwartet.
- **Reale, grüne Läufe der DB-Tiers:** `make test-store` (Persistenz und Lesen
  von `route_target`, Antrags-Funktionen, Grants, Regelstand-Ableitung),
  `make test-replication` (der `Assembler` am realen Stream; die DB-Adapter-Coverage
  prüft `make test-replication` gegen ihre Schwelle) und `make schema-rollout`
  zweimal hintereinander gegen dieselbe Ziel-Datenbank (Idempotenz). Dazu der
  **Alt-Tag-Lauf** von `tools/harness/run-schema-rollout-guard-test.sh`
  ([`ADR-0114`](../adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 7): das Schema des jüngsten `v*`-Tags ausrollen, danach den
  Arbeitsbaum — Exit 0 zweimal, der Datenstand über `cdc.changes` lesbar. Anders
  als bei der Transformations-Welle **ändert diese Welle eine bestehende View**
  (`cdc.changes` erhält `route_target` als letzte Spalte): der Lauf zeigt den
  Vorlauf `DROP VIEW cdc.changes` der Wache über einen Alt-Bestand (erwartet, am
  Lauf zu belegen; `slice-routing-kern-label` trägt den Beleg).
- **Die Fitness-Function-Tests aus
  [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  sind grün und tragen ihre Eingabe:** Determinismus und Reihenfolge der
  Auswertung (Domäne), `ErrRoutingNotApplicable` bei fehlender Spalte und
  Nicht-Treffer bei abwesendem Wert (`mapper`), der Eigenschaftstest R3 in beide
  Richtungen, der Store-Test (`route_target` überlebt Persistierung und Lesen,
  NULL liest als „nicht geroutet"), `make a-check` (die Regel in
  `internal/domain/**`), `make generated-sync` (Proto-Änderung).
- `make doc-trace` führt [`LH-FA-CFG-008`](../../../spec/lastenheft.md) nicht
  mehr unter den Waisen: die Closure-Messung ist die gedruckte Zeile am
  Arbeitsbaum der Closure und nennt `0 Waise(n)`. **Parent-Stand** (gemessen,
  `make doc-trace` am Arbeitsbaum auf `30fd6cb5`, Exit 0): `80 Anforderung(en),
  1 Waise(n)`. Der Träger der Deckung ist die Zeile in
  [`docs/user/e2e-abdeckung.md`](../../user/e2e-abdeckung.md), die der Runner von
  `make test-integration` schreibt (die Datei ist ein Erzeugnis, kein
  Lauf-Beleg); mit der Deckung zieht `slice-routing-e2e` die Aussage über die
  Waisen in [`harness/README.md`](../../../harness/README.md) (Zeile
  `make doc-trace`) nach.
- Die Auswertung liegt an genau einer Stelle: der Suchlauf über `internal/**`
  nach den Aufrufern der Routing-Auswertung findet eine Auswertungsstelle, die
  WAL-Pfad und Backfill-Pfad gemeinsam aufrufen (Befehl und Fundstellen im
  Closure-Bericht).
- Das Benutzerhandbuch trägt den Abschnitt zur Routing-Konfiguration samt der
  Konsequenzen (Label fest zum Erfassungszeitpunkt, keine Rückwirkung, Change
  ohne Treffer nur ungefiltert und über `route_target IS NULL` sichtbar,
  DELETE und Inhaltsregeln), Abhilfe-Prozedur und Fehlerklassen-Zeile — jede Zahl
  und Wirkungs-Aussage mit ihrem Ursprung ([`AGENTS.md`](../../../AGENTS.md)
  §3.12); die drei SDK-Packages und die Beispiel-Clients tragen den Parameter
  `target` (`slice-routing-sdk-beispiel-target`), ohne interne Kennungen
  (`make sdk-public-doc-check` grün).
- Der **Lese-Schritt** des Beobachtungs-Registers ist gelaufen (Einträge bei 3×
  oder darüber, Modul 6).
- Closure-Notiz in `welle-routing-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-routing-spec-nachzug | Pflichtenheft (`LH-FA-CFG-008.a` beantwortet, `SPEC-019`, neue Regelform-Kennung, `SPEC-001`/`-002`, `SPEC-008`, Filterparameter und Zusatz-Subjekt in `SPEC-020`/`-021`/`-022`/`-024`) auf den beschlossenen Stand ziehen — ohne ADR-/Slice-Bezug | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-CFG-008.a`](../../../spec/pflichtenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 1 |
| slice-routing-kern-label | Domäne: Routing-Regel und reine Auswertung, `Change.RouteTarget`, Auswertung im `Assembler`, `ErrRoutingNotApplicable` und Klassen-Abbildung; Persistenz der Spalte `route_target`, View-Spalte von `cdc.changes` mit Rollout-Vorlauf | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-DAT-006`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 2, [`ADR-0114`](../adr/0114-schema-rollout-vorlauf-view-signatur.md) |
| slice-routing-antragsweg | Antragsarten `set_route`/`remove_route`, SQL-Funktionen, Use Cases mit R1–R6 (einschließlich der Sperre gegen `exclude_column`), Regelstand-Ableitung, Verdrahtung in `applyAdministrationRequest`, Aktivierungs-Zweig, Prozessstart | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-ADM-001`](../../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 3, [`ADR-0065`](../adr/0065-spaltenausschluss-dauerhafter-traeger.md) |
| slice-routing-backfill-pfad | Auswertung im Backfill-Run: Backfill-Changes tragen das Label des Regelstands zum Run; Fail-closed um den Regelstand, Nichtanwendbarkeit im Run | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 6 |
| slice-routing-lesewege | Parameter `target` an `GET /changes`, am gRPC-Stream (Proto-Feld, Generierung) und am SSE-Stream; gemeinsame Filterfunktion mit dem Tabellenfilter | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 4, [`ADR-0133`](../adr/0133-tabellen-granulare-filterung-grpc-sse.md) |
| slice-routing-nats-subjekt | zusätzliche Veröffentlichung jeder gerouteten Change auf `cdc.route.<source_id>.<ziel>` im NATS-Vollinhalts-Weg | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 5, [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md) |
| slice-routing-e2e | `make test-integration`: Regel per SQL, alle Lesewege, Boundary R1–R6, Negative, Neustart, Ausschluss-Sperre, Replay, Backfill-Bestand, DELETE-Messung an PostgreSQL 17 und 18; RTM-Träger für `LH-FA-CFG-008` | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../../spec/lastenheft.md), [`LH-FA-ADM-003`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 7 |
| slice-routing-betriebsdoku | Benutzerhandbuch: Routing-Konfiguration, Konsequenzen, DELETE und Inhaltsregeln, Abhilfe, Fehlerklasse | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 8 (Handbuch-Hälfte) |
| slice-routing-sdk-beispiel-target | Parameter `target` in den drei SDK-Packages (HTTP, gRPC, SSE, NATS) und den Beispiel-Clients; SDK-Dokumentation ohne interne Kennungen | [`LH-FA-CFG-008`](../../../spec/lastenheft.md), [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 8 (SDK-Hälfte) |

**Reihenfolge:** sequentiell in der Tabellen-Reihenfolge (WIP-Limit 1 je
Rolleninhaber, Baseline-Regelwerk `modul-05-planning-harness.md`); jeder Slice
startet, sobald sein Start-Trigger gilt — die Kanten stehen in §5. Die
Reihenfolge der Tabelle folgt einer Regel: **erst die Wirkung, dann der Beleg,
dann die Doku** — und der Backfill-Pfad **unmittelbar nach** dem Antragsweg, damit
das Fenster, in dem Regeln setzbar sind, ein Backfill-Bestand aber noch kein
Label trägt, ein Slice lang ist (§5).

**Abweichungen vom Schnitt-Vorschlag** in
[`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
§Folgepflichten (dort ausdrücklich „der Planner formt Welle und Slices"),
je mit Grund:

1. **Folgepflicht 8 geteilt, neun statt acht Slices.** Die Pflicht nennt zwei
   Gegenstände mit verschiedener Abhängigkeit und verschiedenem Gate: das
   Handbuch darf nur beschreiben, was `slice-routing-e2e` belegt hat
   (`betriebsdoku`, Server-Seite, `make docs-check`), die SDK- und
   Beispiel-Parameter hängen allein an den Wire-Formen aus `lesewege` und
   `nats-subjekt`, bauen in drei fremden Bau-Kontexten (`sdks/`, `examples/`) und
   unterliegen der Regel „keine interne Kennung unter `sdks/`"
   (`make sdk-public-doc-check`). Die Pflicht selbst nennt die SDKs als eigenen
   Folge-Schritt „wie bei `ADR-0133`" (dort je Sprache ein Slice).
2. **Reihenfolge der Folgepflichten 4, 5, 6.** Die ADR lässt 4, 5 und 6
   untereinander unabhängig. Der Backfill-Pfad (6) steht hier vor den Lesewegen
   (4, 5): sonst wäre ein Regelstand setzbar (nach `antragsweg`), während ein
   Backfill-Run Changes ohne Label liefert, über die Dauer von zwei Slices
   statt einem — dieselbe Regel wie in der Welle
   [welle-transformationen](done/welle-transformationen.md) §4 Abweichung 5.
3. **Folgepflicht 3 bleibt ein Slice** (anders als in der Welle
   [welle-transformationen](done/welle-transformationen.md), die ihn teilte): die
   Spalten `rule_name`/`rule_spec` und das Muster der Funktionen stehen seit den
   Transformationen, der Schema-Anteil ist deshalb klein (zwei Funktionen, die
   `request_kind`-Menge, Grants, Guard). Der Slice trägt die Teilungsnaht in
   seiner Rückführung (`antragsweg-schema` / `antragsweg-usecase`); er ist mit drei
   Liefer-Punkten der größte der Welle.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Wird blockiert von:** keiner Welle;
  [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) ist
  `Accepted`, der Trigger (§2) ist erfüllt. Offene Slices anderer Wellen gibt es
  nicht (`docs/plan/planning/open/` und `in-progress/` tragen vor dieser
  Eröffnung keinen Slice-Plan, gemessen am Arbeitsbaum auf `30fd6cb5`); WIP-Limit 1
  gilt trotzdem repo-weit.
- **Blockiert:** keine Welle.
- **Interne Kanten** (`A → B` heißt: `B` setzt `A` voraus, jede Stufe braucht das
  lauffähige System der Vorstufe): `spec-nachzug` → jeder übrige Slice (die Spec
  führt) · `kern-label` → `antragsweg` → `backfill-pfad` · `kern-label` →
  `lesewege` · `kern-label` → `nats-subjekt` · `antragsweg`, `backfill-pfad`,
  `lesewege`, `nats-subjekt` → `e2e` → `betriebsdoku` → `sdk-beispiel-target`.
  Gemessen gegen die ADR-Abhängigkeiten (1 vor allen; 2 vor 3–6; 3 vor 6 und 7; 4, 5,
  6 untereinander unabhängig; 7 nach 2–6; 8 nach 7): die Welle erfüllt sie bis
  auf die Wahl der Reihenfolge innerhalb der unabhängigen Gruppe 4–6 (§4,
  Abweichung 2). Die ADR führt die Abhängigkeiten selbst als „hergeleitet, nicht
  gegen einen Plan geprüft"; die Kanten oben sind die geprüfte Fassung.
- **Spec-Kollision.** `spec-nachzug` vergibt die nächste freie `SPEC`-Kennung für
  die Regelform; die Kennungsvergabe im Pflichtenheft zählt fortlaufend je Datei.
  Stand am Parent `30fd6cb5` (gemessen, `git grep -h -o 'SPEC-0[0-9][0-9]'
  30fd6cb5 -- spec/pflichtenheft.md | sort -u | tail -1`): `SPEC-031`; der Slice
  misst an seinem Start neu.
- **Vorab-Bedingungen mit Adresse** (Lücken im Text von
  [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md), die
  der Planner nicht entscheidet — jede steht als Start-Trigger im Slice, der sie
  braucht, und als Frage im Eröffnungs-Bericht; die Antwort ist ein
  Architect-Verdikt bzw. eine Folge-ADR, keine Auslegung durch den Implementer):
  - **V1 — gRPC-`ReadChanges`.** Die Tabelle in Teilfrage 5 nennt `GET /changes`,
    den `StreamChangesRequest`, SSE und NATS, nicht den zehnten RPC des
    `Administration`-Service
    ([`ADR-0131`](../adr/0131-grpc-readchanges-zehnter-rpc.md): derselbe Use Case
    wie `GET /changes`, Request mit `source`, `schema`, `table`, `from`, `to`,
    `limit`). Ohne `target` dort ist die Gleichwertigkeit der Zugriffswege
    ([`LH-FA-SST-006`](../../../spec/lastenheft.md)) für das Routing verletzt.
    Adresse: `slice-routing-lesewege` (Start-Trigger).
  - **V2 — Nichtanwendbarkeit im Backfill-Run.**
    [`ADR-0117`](../adr/0117-backfill-run-fehlerklasse-schema.md) setzt die Klasse
    `schema` für eine im Run nicht anwendbare **Transformationsregel** fest;
    [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
    Teilfrage 6 sagt nur, dass der Backfill dieselbe Auswertung durchläuft. Ob
    dieselbe Run-Behandlung (Run `failed`, Klasse `schema`, run-lokal, ohne Change)
    für eine Routing-Regel gilt, steht nirgends. Adresse:
    `slice-routing-backfill-pfad` (Start-Trigger: Kurzverdikt).
  - **V3 — Erreichbarkeit der Nichtanwendbarkeit.** Eine Routing-Regel ist nicht
    anwendbar, wenn `when.column` in der Relation der Change fehlt
    ([`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
    Teilfrage 4). Die naheliegende Ursache — Spalte entfernt — endet nach dem
    Bestand im Pfad der inkompatiblen Schemaänderung (`ErrIncompatibleSchemaChange`,
    Fehlerklasse `schema`, `LH-FA-SCH-003`), **bevor** eine Change assembliert wird
    (*hergeleitet* aus dem Aufbau von `Assembler.observeRelation`, am Code nicht für
    diesen Fall geprüft). Welche Ursache `ErrRoutingNotApplicable` am laufenden
    System erzeugt — und ob der Abhilfe-Weg „`cdc.remove_route` und Neustart"
    dann greift — ist offen. Adresse: Messung im ersten Schritt von
    `slice-routing-kern-label` (Unit-Ebene) und von `slice-routing-e2e`
    (Systemebene); bei unerreichbarem Fall ein Architect-Verdikt vor dem
    Negative-Beleg.
  - **V4 — Pflichtenheft-Kennung der Regelform.** Kein Entscheidungsbedarf: die
    Regelform bekommt eine eigene Kennung neben `SPEC-030`
    ([`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
    Folgepflicht 1);
    `spec-nachzug` entscheidet die Aufteilung zwischen `SPEC-019` und der neuen
    Kennung wie `slice-transformationen-spec-nachzug`.
- **Benanntes Zwischenzustands-Fenster.** Ab `antragsweg` sind Routing-Regeln
  setzbar und `cdc.changes` trägt die Spalte; ein Backfill-Run vor `backfill-pfad`
  lieferte `route_target = NULL`, ein Leser vor `lesewege` kann das Label nur über
  SQL auswählen. Das Backfill-Fenster ist ein Slice lang (§4 Reihenfolge);
  Risiko und Ausgang tragen `antragsweg` §6 und `backfill-pfad` §6.
- **Benannte Kopplung: View-Signatur.** `kern-label` fügt `route_target` als
  letzte Spalte der bestehenden View `cdc.changes` hinzu; der Rollout läuft über
  den Vorlauf von [`ADR-0114`](../adr/0114-schema-rollout-vorlauf-view-signatur.md).
  Der Slice trägt den Alt-Tag-Lauf (§3); ein späterer Slice, der die View ändert,
  kollidiert mit seinem Stand — keiner ist geplant.
- **Ereignis-Adresse.** „Die Closure dieser Welle" ist eine Adresse, die
  eintreten kann: die Roadmap führt diese Welle unter *Offene Wellen*, und die
  Welle hat einen Closure-Trigger (§3).

**Träger der Folgepflichten** — jede Pflicht aus
[`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
§Folgepflichten hat einen Träger oder eine benannte Adresse
(`BEO-PGC/adr-folgepflicht-ohne-traeger-slice`):

| Folgepflicht | Träger |
|---|---|
| 1 Spec-Nachzug | `slice-routing-spec-nachzug` |
| 2 Kern: Domäne, `Assembler`, Persistenz, View | `slice-routing-kern-label` |
| 3 Antragsweg und Dauerhaftigkeit | `slice-routing-antragsweg` |
| 4 Lesewege (`GET /changes`, gRPC, SSE) | `slice-routing-lesewege`; Vorab-Bedingung V1 für den RPC `ReadChanges` |
| 5 NATS-Zusatz-Subjekt | `slice-routing-nats-subjekt` |
| 6 Backfill-Pfad | `slice-routing-backfill-pfad`; Vorab-Bedingung V2 |
| 7 E2E-Belege samt DELETE-Messung | `slice-routing-e2e`; Vorab-Bedingung V3 |
| 8 Betriebsdokumentation und SDK-/Beispiel-Parameter | `slice-routing-betriebsdoku` (Handbuch) und `slice-routing-sdk-beispiel-target` (SDK, Beispiele) |

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Aktive Senken** (Webhook, Queue, Dead-Letter, Secrets) und die
  **Consumer-Bindung** an ein Ziel —
  [`ADR-0137`](../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Teilfrage 1 (Optionen B, D) und Entscheidung 1: Re-Evaluierungs-Trigger mit
  `Supersedes`, kein Teil dieser Welle.
- **Mehrere Ziele je Change** (Fan-out) — Entscheidung 5; ein Ziel genügt.
- **Weitere Operatoren** (Präfix, Menge, Vergleich) — der Satz ist `equals`; ein
  weiterer braucht eine Folge-Entscheidung (Trigger: dreimal im
  Beobachtungs-Register).
- **Umetikettieren und Rückwirkung** — Entscheidung 6: das Label ist zum
  Erfassungszeitpunkt fest; der Altbestand wird über einen Backfill-Run mit dem
  aktuellen Regelstand neu erzeugt, nicht umgeschrieben.
- **Warnung oder Sicht für Changes ohne Treffer** — Konsequenz in der ADR:
  `route_target IS NULL` ist die Sicht des Betreibers; eine Warnung ist nicht Teil.
- **Änderung des Nachrichtenschemas** — die Live-Nachrichten
  ([`SPEC-020`](../../../spec/pflichtenheft.md)/[`SPEC-021`](../../../spec/pflichtenheft.md)/[`SPEC-024`](../../../spec/pflichtenheft.md))
  tragen zehn Felder, die HTTP-Antwort von
  [`SPEC-022`](../../../spec/pflichtenheft.md) dreizehn; das Label ist nicht Teil
  der Nachrichten, sondern über die SQL-Sicht lesbar. Das Wecksignal
  ([`SPEC-017`](../../../spec/pflichtenheft.md)) bleibt unverändert. Die
  Proto-Änderung beschränkt sich auf das Anfrage-Feld `target`
  (`make generated-sync`, `slice-routing-lesewege`).
- **Realserver-E2E der SDK-Packages für das Routing** — die drei SDK-Realserver-Tiers
  (`make test-sdk-*-integration`) bekommen in dieser Welle keine Routing-Phase;
  `slice-routing-sdk-beispiel-target` belegt den Parameter mit Unit-Tests und dem
  `sdk-public-doc-check`. Ein Realserver-Beleg wäre ein weiterer Slice (Frage an den
  Auftraggeber im Eröffnungs-Bericht).
- **Ein allgemeiner Recovery-Weg für Schema-Fehler** —
  `BEO-PGC/kein-admin-weg-schema-fehler-recovery`: die Welle liefert die Abhilfe
  für die Nichtanwendbarkeit einer Routing-Regel (soweit V3 sie erreichbar macht),
  nicht die Wiederinbetriebnahme nach einer entfernten Spalte.
- **Kein neuer GitHub-Actions-Workflow und keine strukturelle Workflow-Änderung** —
  [`AGENTS.md`](../../../AGENTS.md) §3.10 greift nicht; `e2e.yml` fährt
  `make test-integration` unverändert weiter (die Matrix über PostgreSQL 17 und 18
  ist der Träger der DELETE-Messung, §3).
- **Kein Server-Release und kein SDK-Release** — `docs/user/version.md` und die
  Package-Versionen der SDKs bleiben unberührt; ein Release braucht eine neue
  Freigabe des Auftraggebers. Die Änderungshistorie des Benutzerhandbuchs trägt
  Zeilen.
- **Keine Schwellen-Senkung und keine neue Gate-Klasse** — neue Use-Case-Pakete
  liegen in der netzlos gemessenen Fläche des Coverage-Gates; eine Senkung bliebe
  per [`AGENTS.md`](../../../AGENTS.md) §3.6 ADR-pflichtig.
- **Lastenheft unverändert** — die Welle schärft Techniken im Pflichtenheft,
  nicht Anforderungen.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt 2):**
Das Register (`docs/plan/planning/observations/README.md`) wurde durchgesehen.
Die Zähler sind die Zahl der Dateien unter `evidence/` (ausgezählt am 2026-10-01
an `docs/plan/planning/observations/BEO-PGC/*/evidence/`). Einträge, die mit
dieser Welle den Stand 3× erreichen könnten, liegen bei 2×; Einträge ab 3× sind
bereits verkörpert, gestrichen oder (siehe unten) offen. Treffer, mit dem Slice,
der sie als Risiko trägt (Detail je Slice in §8):

- `BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (offen, 3×, Ausgang
  nicht zugewiesen — Lese-Schritt der Closure) — **einschlägig:** `spec-nachzug`,
  `betriebsdoku` und `sdk-beispiel-target` beschreiben dieselben Parameter und
  Funktionen in Pflichtenheft, Handbuch und SDK-README.
- `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×) — **einschlägig:**
  die Abhilfe nach einer nicht anwendbaren Regel (V3), nicht die allgemeine
  Recovery (§6).
- `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 2×) —
  **einschlägig:** V1 bis V3 sind Stellen, an denen der ADR-Wortlaut die Umsetzung
  nicht trägt; sie sind als Fragen geführt, nicht als stille Abweichung.
- `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (verkörpert, 5×) und
  `BEO-PGC/aufschub-adresse-verfaellt` (verkörpert, 3×) — jeder Aufschub in den
  Plänen dieser Welle trägt eine Adresse, deren §2 den Gegenstand deckt;
  „die Closure dieser Welle" ist eine Adresse, die eintreten kann (§5).
- `BEO-PGC/adr-folgepflicht-ohne-traeger-slice` (offen, 1×) — die Träger-Tabelle
  in §5 gibt jeder Folgepflicht einen Träger.
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 34×,
  [`AGENTS.md`](../../../AGENTS.md) §3.13) — betrifft jeden Slice; der Suchlauf
  steht als committetes Feld je Slice-Plan §3.
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×) —
  `spec-nachzug`, `betriebsdoku`.
- `BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×),
  `BEO-PGC/d-migrate-nacharbeit` (verkörpert, 8×) — `kern-label`, `antragsweg`:
  neue Spalte, View-Spalte, CHECK-Menge, zwei Funktionen, Guard-Eintrag, Grants.
- `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (verkörpert, 4×) —
  `lesewege`: Proto-Änderung, `make generated-sync`.
- `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (offen, 1×),
  `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (verkörpert, 3×),
  `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
  `antragsweg`, `kern-label`: der Regelstand ist ein dauerhafter, tabellen-scoped
  Träger, der `Assembler`-Cache ist Laufzeit; neuer Code verschiebt Nenner und
  Zähler der Coverage-Messungen.
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
  `BEO-PGC/test-isolation-geteilter-zustand` (offen, 2×),
  `BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×),
  `BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 5×) — `e2e`.
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×),
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 18×),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 28×),
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 11×),
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 10×) — Slices mit
  Messungen und Negativ-Belegen (`e2e`, `kern-label`, `lesewege`,
  `nats-subjekt`).
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert, 3×), `BEO-PGC/handbuch-versionshistorie-uebersprungen`
  (verkörpert, 3×) — `betriebsdoku`: Handbuch und Änderungshistorie; die
  Betreiber-Oberfläche (zwei SQL-Funktionen) entsteht in `antragsweg`, das Handbuch
  folgt fünf Slices später (§4 Abweichung 1 und Aufschub-Adresse in
  `antragsweg` §2).
- `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 2×),
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert, 3×),
  `BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen, 1×),
  `BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen` (gestrichen, 2×) —
  `sdk-beispiel-target`: drei Sprachen, öffentliche Texte, ein neues Anfrage-Feld.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 8×) — **nicht
  einschlägig:** kein Workflow-Zug in dieser Welle.
- Gesichtet, ohne Bezug zu dieser Welle: die übrigen Einträge des Registers.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-routing-results.md`
Zähler: `../observations/` (Beobachtungs-Register)
