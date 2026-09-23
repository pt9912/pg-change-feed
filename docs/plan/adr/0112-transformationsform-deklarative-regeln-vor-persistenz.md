# ADR-0112: Transformationen — geschlossener Satz deklarativer Regeln, per Antrags-Queue konfiguriert, vor der Persistierung im Assembler ausgewertet

**Status:** Accepted

**Datum:** 2026-09-23

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-23)

**Bezug:** [`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Haupt-Bezug —
Transformationen zwischen Erfassung und Zustellung; Out-of-Scope
schließt eine vollständige Transformationssprache und ein
Skripting-/Plugin-Modell aus),
[`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Routing auf
Zustellziele — eigene Anforderung, in dieser ADR nicht umgesetzt),
[`LH-FA-CFG-005`](../../../spec/lastenheft.md) (Spaltenauswahl — die reine
Auswahl bleibt dort),
[`LH-QA-SEC-004`](../../../spec/lastenheft.md) (ausgeschlossene
Spaltenwerte erscheinen nicht in den Changes),
[`LH-FA-ADM-003`](../../../spec/lastenheft.md) (sichtbare Fehlerzustände),
[`LH-FA-SCH-004`](../../../spec/lastenheft.md) (nicht sicher
interpretierbare Schemaänderung endet sichtbar),
[`LH-FA-REA-005`](../../../spec/lastenheft.md) (erneutes Lesen innerhalb der
Aufbewahrung),
[`LH-FA-DAT-005`](../../../spec/lastenheft.md) (Abwesenheit erkennbar),
[`LH-FA-SST-008`](../../../spec/lastenheft.md) (Live-Streaming),
[ADR-0050](0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue mit Live-Reload — Erweiterungsbasis),
[ADR-0059](0059-spaltenauswahl-mechanismus.md) (Spaltenausschluss —
Muster für Mechanismus, Wirkort und Fehlerpfad),
[ADR-0065](0065-spaltenausschluss-dauerhafter-traeger.md) (dauerhafter,
tabellen-scoped Träger — Muster für den Regelstand),
[ADR-0060](0060-grpc-streaming-mechanismus.md),
[ADR-0061](0061-http-sse-zusaetzlich-zu-grpc.md),
[ADR-0100](0100-nats-dritter-vollinhalts-zustellweg.md) (Zustellwege auf
demselben `Broadcaster`),
[ADR-0046](0046-sql-driving-adapter-lese-schreib-trennung.md) (keine
Domänenlogik in SQL),
[ADR-0023](0023-fehlerklassifikation.md) (sieben Fehlerklassen),
[ADR-0052](0052-optionale-yaml-konfigurationsdatei.md) (Konfigurationsdatei —
geprüfte, verworfene Alternative)

**Schärft:** [`LH-FA-CFG-007.a`](../../../spec/pflichtenheft.md)
(Konfigurationsmechanismus, Ausdrucksform, Auswertungsreihenfolge und
Verhältnis zum Spaltenausschluss),
[`SPEC-019`](../../../spec/pflichtenheft.md) (Antrags-Datensatz — zwei
weitere Antragsarten),
[`ARC-005`](../../../spec/architecture.md) (Antragsqueue-Pfad trägt zwei
weitere Antragsarten),
[`ARC-002`](../../../spec/architecture.md) (zwei neue Use Cases),
[`ARC-004`](../../../spec/architecture.md) (ein neuer Outbound Port)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-CFG-007`](../../../spec/lastenheft.md) fordert, erfasste Changes vor
der Zustellung gezielt zu transformieren (Feld-Umbenennung, wertbasierte
Ableitung über die reine Spaltenauswahl hinaus); das Routing auf
unterschiedliche Zustellziele ist
[`LH-FA-CFG-008`](../../../spec/lastenheft.md). Die Anforderung nennt ihre offenen
Fragen selbst: die Ausdrucksform sei „Architektur- (ADR) bzw.
Spezifikationsfrage", eine vollständige Transformationssprache oder ein
Skripting-/Plugin-Modell (vergleichbar Kafka-Connect-SMTs) sei
ausdrücklich nicht gefordert. Das Pflichtenheft führt die Frage als
[`LH-FA-CFG-007.a`](../../../spec/pflichtenheft.md): Konfigurationsmechanismus,
Ausdrucksform, Auswertungsreihenfolge bei mehreren zutreffenden Regeln und
das Verhältnis zur Spaltenausschluss-Antragsart
([`SPEC-019`](../../../spec/pflichtenheft.md)).

Der Code trägt heute keinen Transformationsmechanismus. Fünf am Code
geprüfte Eigenschaften legen den Lösungsraum fest (Beleg-Anker je Punkt):

1. **Ein Erzeugungspfad, ein `model.Change`, viele Leser.** Der einzige
   Pfad, der `model.Change`-Werte erzeugt, ist `Assembler.change`
   (`internal/adapters/driving/replication/mapper/mapper.go`).
   `CaptureService.Capture` (`internal/application/usecase/capture/service.go`)
   persistiert die Transaktion und veröffentlicht danach dieselben
   Change-Werte über `ChangeStreamPort.Publish` an den `Broadcaster`;
   der gRPC-Server, der SSE-Handler und der NATS-Vollinhalts-Publisher
   abonnieren ihn (`ADR-0060`/`ADR-0061`/`ADR-0100`). Der Lesezugriff über SQL
   projiziert `c.old_data`/`c.new_data` aus `cdc.change`
   (`tools/schema/schema.yaml`, `views: changes`); `GET /changes` liest
   dieselbe Projektion (Kommentar an `queries.go`, Zeile ~39). Das
   NATS-Wecksignal (`SPEC-017`) trägt keinen Inhalt.
2. **Der Wirkort für „Wert nie serialisiert" existiert bereits.**
   `rowImage` (`mapper.go`) baut das Row Image spaltenweise aus der
   Relation-Nachricht, überspringt `nil`-Werte und Spalten aus
   `TableBinding.ExcludedColumns` (`containsColumn`) und schreibt jeden
   verbleibenden Wert als JSON-String (`json.Marshal(*values[i])`) —
   der Text-Stand der Quelle, ohne Typ-Interpretation (`ADR-0016`).
3. **Der Live-Reload- und Dauerhaftigkeits-Pfad existiert bereits.**
   Die Antrags-Queue (`cdc.administration_request`, `SPEC-019`) trägt vier
   Antragsarten; `applyAdministrationRequest`
   (`internal/bootstrap/wiring.go`) ruft je Art einen Use Case und trägt
   danach die laufende `Assembler`-Bindung nach; der dauerhafte
   Ausschlussstand wird aus den `applied`-Zeilen abgeleitet und beim
   Anlegen jeder Bindung mitgeführt (`activatedTableBindings`,
   Aktivierungs-Zweig, `ColumnExclusionPort.ExcludedColumns`, `ADR-0065`).
4. **Ein nicht interpretierbares Schema endet sichtbar.** Jede Relation-
   Änderung außer der reinen Spalten-Erweiterung meldet
   `ErrIncompatibleSchemaChange` (`observeRelation`, `relationOther`);
   `classifyRunError` (`wiring.go`) bildet sie auf die Fehlerklasse
   `schema` ab, `reportFault` schreibt sie in den Heartbeat, `diagnose`
   zeigt sie an (`LH-FA-ADM-003`). Die Fehlerklassen sind eine
   geschlossene Menge von sieben (`model.ErrorClass`, `ADR-0023`).
5. **Es gibt kein Zustellziel-Modell.** gRPC und SSE liefern ungefiltert
   alle Changes (`SPEC-020`: „tabellen-granulare Filterung ist nicht Teil
   dieser Version"); die einzige adressierbare Zielform ist das
   NATS-Subjekt `cdc.stream.<source_id>.<schema>.<table>` (`SPEC-024`,
   `Accepted`, unberührbar nach `AGENTS.md` §3.5).

Die Clients behandeln die Row Images als undurchsichtiges JSON
(`JsonElement?` in den gelesenen C#- und Kotlin-Modellen, `Any | None` in
`sdks/python/pgchangefeed/src/pgchangefeed/models.py`) — eine geänderte
Schlüsselmenge eines Images ändert das Nachrichtenschema nicht.

## Entscheidung

Wir wählen: **einen geschlossenen Satz deklarativer Regeln, konfiguriert
über zwei neue Antragsarten der Antrags-Queue aus `ADR-0050`, dauerhaft
nach dem Muster von `ADR-0065`, ausgewertet im `Assembler` bei der
Row-Image-Konstruktion — nach dem Spaltenausschluss, vor jeder
Serialisierung und Persistierung.** Routing (`LH-FA-CFG-008`) kommt nicht in diesem Zug
(Teilfrage 7). Acht Teilfragen, acht Festlegungen.

### Teilfrage 1 — Konfigurationsmechanismus

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | kein Aufwand | verletzt eine aktive, vertraglich abnahmebindende Anforderung |
| B — YAML-Konfigurationsdatei (`ADR-0052`/`SPEC-016`), Regeln je Tabelle | strukturierte Regeln (Wertzuordnungen) sind in YAML natürlich ausdrückbar; Start-Validierung wäre fail-fast | Boot-Zeit-Konfiguration ohne Reload (dieselbe Neustart-Einschränkung wie `ADR-0059` Teilfrage 1 Option C); kein Rückkanal für einen Antrag gegen eine nicht existierende Spalte; eine zweite Quelle der Wahrheit neben dem SQL-verwalteten Ausschlussstand derselben Tabelle — zwei Konfigurationswege für zwei Aspekte derselben Bindung |
| C — Umgebungsvariable | kein neuer Code-Pfad im Store | Wertzuordnungen als Zeichenkette unhandlich; Boot-Zeit-Konfiguration; kein Rückkanal |
| **D — zwei Antragsarten auf der Antrags-Queue (gewählt)** | ein Konfigurationsweg für alle administrativen Fähigkeiten, die den laufenden Prozess ohne Neustart erreichen müssen; Live-Reload und Rückkanal (`pending`/`applied`/`failed`, `error_message`) sind vorhanden; Validierung liegt im Go-Use-Case, nicht in SQL (`ADR-0046`) | erweitert `chk_administration_request_kind` und die Spaltenform von `cdc.administration_request` um zwei nullable Spalten (additive Migration); ein JSON-Parameter je Regel ist in SQL nur als `jsonb` durchreichbar, die Prüfung liegt vollständig in Go |

Festlegung: die SQL-Funktionen `cdc.set_transformation(source_id,
schema_name, table_name, rule_name, rule_spec jsonb)` und
`cdc.remove_transformation(source_id, schema_name, table_name, rule_name)`
schreiben ausschließlich eine Antrags-Zeile (`request_kind`
`set_transformation` bzw. `remove_transformation`) und senden
`pg_notify` — wie `cdc.exclude_column`. Beide Funktionen sind
ausschließlich der Rolle `cdc_admin` ausführbar (`LH-QA-SEC-002`, wie die
bestehenden vier). Neue Antrags-Spalten: `rule_name text`,
`rule_spec jsonb`, beide nullable; für die beiden neuen Antragsarten
Pflicht (`rule_spec` nur für `set_transformation`). Granularität:
ausschließlich je Tabelle (Begründung wie `ADR-0059` Teilfrage 2). Das
Pflichtenheft (`SPEC-019`) führt die Feldform nach.

### Teilfrage 2 — Ausdrucksform: ein geschlossener Regelsatz, keine Sprache

| Option | Pro | Contra |
|---|---|---|
| A — freier Ausdruck/Skript (z. B. Go-Template, SQL-Ausdruck, JavaScript) | beliebig ausdrucksstark | ist genau die vom Lastenheft ausgeschlossene Transformationssprache/das Skripting-Modell; ein SQL-Ausdruck verlagert Auswertung und Fehlerbehandlung in die Datenbank (`ADR-0046`); unbegrenzte Fehlerfläche und Sicherheitsfläche im Capture-Pfad |
| B — Regelsatz aus der Konfiguration nachladbar (Plugin) | Erweiterbarkeit ohne Release | Plugin-Modell, ausgeschlossen |
| **C — geschlossener Satz deklarativer Regeltypen (gewählt)** | jede Regel hat eine feste, vollständig testbare Semantik; Konfliktfreiheit ist statisch prüfbar (Teilfrage 3); keine Sprache | neue Bedürfnisse brauchen einen neuen Regeltyp und damit eine Folge-Entscheidung |

Festlegung: zwei Regeltypen. `rule_spec` ist ein JSON-Objekt mit dem
Pflichtschlüssel `kind`; unbekannte Schlüssel und unbekannte `kind`-Werte
enden als `failed`.

| `kind` | Schlüssel | Wirkung auf das Row Image (`old_image` und `new_image`) |
|---|---|---|
| `rename_column` | `column`, `to` | der Schlüssel `column` heißt im Image `to`; der Wert bleibt unverändert |
| `map_value` | `column`, `values` (Objekt `alt → neu`, Zeichenketten) | ist der Wert von `column` als Zeichenkette ein Schlüssel von `values`, steht der zugeordnete Wert im Image; jeder andere Wert bleibt unverändert |

Ein Wert, der im Image fehlt (NULL, unverändertes TOAST,
ausgeschlossene Spalte), bleibt für `rename_column` und `map_value`
abwesend — derselbe Abwesenheits-Vertrag wie in `rowImage`
(`LH-FA-DAT-005` Negative), kein Platzhalter. Ein fehlendes Bild bleibt
fehlend (`LH-FA-CAP-008`). Ein Regeltyp, der ein konstantes Zusatzfeld
ergänzt, gehört nicht zum Satz: `LH-FA-CFG-007` nennt Umbenennung und
wertbasierte Ableitung, keine Ergänzung; ein weiterer Regeltyp braucht
eine Folge-Entscheidung (Re-Evaluierungs-Trigger). Regeln wirken ausschließlich auf Schlüssel und Werte der
Row Images; `change_id`, `transaction_id`, `source_table_id`, `sequence`,
`operation`, `schema_version`, `schema` und `table` bleiben unverändert —
insbesondere bleibt die Tabellen-Identität, aus der das NATS-Subjekt
(`SPEC-024`) und der Tabellenfilter von `GET /changes` folgen, die
Quell-Identität.

### Teilfrage 3 — Auswertungsreihenfolge und Auflösung bei Mehrdeutigkeit

Festlegung: **Mehrdeutigkeit wird statisch ausgeschlossen, nicht zur
Laufzeit aufgelöst.** Vier Konfliktfreiheits-Invarianten je Tabelle; ein
Antrag, der eine verletzen würde, endet `failed` mit Fehlertext
(`error_message`), der Regelstand bleibt unverändert:

- **K1** — `rule_name` ist je Tabelle eindeutig; ein `set_transformation`
  gegen einen vergebenen Namen endet `failed` (erst entfernen, dann neu
  setzen).
- **K2** — jede Quellspalte trägt höchstens eine Spaltenregel
  (`rename_column` oder `map_value`).
- **K3** — die Zielnamen (`to` von `rename_column`) sind untereinander
  verschieden und verschieden von jedem Spaltennamen der Quelltabelle
  (ausgeschlossene Spalten eingeschlossen).
- **K4** — `column` existiert an der Quelle (`ColumnExists`, Muster von
  `ErrSourceColumnMissing`, `ADR-0059` Teilfrage 5); `remove_transformation`
  gegen einen nicht geführten Namen endet `failed`.

Unter K1–K4 wirken die Regeln unabhängig voneinander; es gibt keinen
Fall, in dem zwei zutreffende Regeln denselben Schlüssel oder denselben
Wert beanspruchen. Die Auswertung ist trotzdem vollständig geordnet
(`LH-FA-CFG-007` Boundary: definiert, nicht stillschweigend eines):
(1) Spaltenausschluss, (2) Spaltenregeln in Relation-Spaltenreihenfolge —
der umbenannte oder abgebildete Schlüssel behält die Position seiner
Quellspalte. Die Reihenfolge bestimmt die Schlüsselreihenfolge des Images
und ist deterministisch; der Regelstand selbst entsteht aus den
Antrags-Zeilen in der Ordnung `requested_at`, bei gleichem Zeitstempel nach
`administration_request_id` (dieselbe Ordnung wie `ADR-0065`).

### Teilfrage 4 — Sichtbarer Fehlerzustand bei nicht anwendbarer Regel

Eine Regel ist auf eine Change nicht anwendbar, wenn ihre `column` in der
Relation der Change fehlt oder ein Zielname (K3) mit einer Spalte der
Relation kollidiert (etwa nach einer kompatiblen Spalten-Erweiterung).

| Option | Pro | Contra |
|---|---|---|
| A — Regel überspringen, Change roh ausliefern, Fehler nur protokollieren | Erfassung läuft weiter | die Change trägt die Rohform, die der Betreiber gerade nicht ausliefern wollte (z. B. den internen Spaltennamen) — die „still unveränderte Auslieferung", die `LH-FA-CFG-007` Negative ausschließt, nur mit Log-Zeile; bei K3-Kollision zwei gleichnamige JSON-Schlüssel |
| B — Change verwerfen, Erfassung läuft weiter | kein Rohdatenaustritt | stiller Datenverlust nach ACK-Fortschritt; verletzt `LH-QA-REL-001.a`-Geist (Persist-before-ACK) |
| **C — Erfassungspfad endet sichtbar mit Fehlerklasse `schema` (gewählt)** | dieselbe, bereits getragene Behandlung wie `LH-FA-SCH-004` (`ErrIncompatibleSchemaChange`): kein ACK für die Transaktion, kein Datenverlust, Heartbeat-Fehlerzustand und `diagnose` machen es erkennbar und unterscheidbar (`LH-FA-ADM-003`); keine achte Fehlerklasse (`ADR-0023`) | die Erfassung der gesamten Quelle steht bis zur Abhilfe (ein Stream, ein Slot); Abhilfe ist `cdc.remove_transformation` bzw. eine passende Regel (Akzeptanzkriterium in Folgepflicht 5) |

Festlegung: `Assembler.change` prüft die Regeln der Bindung gegen
`event.Relation.Columns`, bevor ein Wert serialisiert wird, und meldet
bei Nichtanwendbarkeit `mapper.ErrTransformationNotApplicable`;
`classifyRunError` bildet ihn auf `model.ErrorClassSchema` ab. Die
Transaktion wird nicht persistiert und nicht bestätigt; nach Abhilfe
setzt die Erfassung ab der bestätigten Position fort (`ADR-0012`).
Eine Änderung der Relation, die die Regel unanwendbar macht, erreicht
diesen Pfad nur, wenn `observeRelation` sie nicht bereits als
`relationOther` meldet — die spalten-entfernenden Fälle enden dort
schon vorher mit derselben Klasse; die K3-Kollision nach kompatibler
Erweiterung ist der Fall, den erst diese Prüfung fängt.

### Teilfrage 5 — Verhältnis zum Spaltenausschluss (`SPEC-019`)

Festlegung: **Der Ausschluss gilt zuerst und ist durch die Struktur der
Auswertung nicht unterlaufbar.** Die Regelauswertung sitzt in derselben
Schleife wie der Ausschluss (`rowImage`) und sieht einen Wert erst, nachdem
der Ausschluss ihn dort zugelassen hat: eine ausgeschlossene Spalte wird
nie gelesen, weder von `rename_column` noch von `map_value`; ihr Schlüssel
erscheint weder unter dem Quell- noch unter einem Zielnamen
(`LH-QA-SEC-004`). K3 hält zusätzlich jeden Zielnamen von den Namen
ausgeschlossener Spalten fern. Ein `exclude_column` gegen eine Spalte mit
Regel und ein `include_column` bleiben zulässig; der Ausschlussstand
entscheidet immer vor dem Regelstand.

### Teilfrage 6 — Wirkort: vor der Persistierung, im `Assembler`

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | kein Aufwand | siehe Teilfrage 1 |
| B — Zustellseite je Weg (Transformation nach der Persistierung in Lese-Use-Case, Stream-Vorstufe, NATS-Publisher) | Rohform bleibt gespeichert; Regeländerung wirkt rückwirkend; je Weg unterschiedliche Form möglich | der SQL-Lesezugriff `cdc.changes` (`LH-FA-SST-002`) kann ohne Domänenlogik in SQL (`ADR-0046`) nicht transformieren und lieferte weiter die Rohform — Happy Path verfehlt, Rohform bleibt an einem Weg lesbar; dieselbe Change läge live in alter, beim erneuten Lesen (`LH-FA-REA-005`) in neuer Form vor, sobald sich Regeln ändern; vier bis fünf Auswertungsstellen statt einer |
| C — Rohform und transformierte Form beide speichern | Rückwirkung und Replay der Rohform möglich | doppelte Speicher- und Schreiblast je Change; bewahrt die Rohform, die die Konfiguration ausblenden soll, auf Platte; zusätzliche Spalte in `cdc.change` und eine Auswahlregel je Leser |
| **D — im `Assembler` bei der Row-Image-Konstruktion (gewählt)** | ein Auswertungsort für alle Wege — Speicher, SQL-Sicht, `GET /changes`, gRPC, SSE, NATS-Vollinhalt sehen dieselbe Form (Punkt 1 im Kontext); dieselbe Change ist live und beim erneuten Lesen identisch (`LH-FA-REA-005`); Ausschluss und Regel liegen in einer Schleife (Teilfrage 5); kein Eingriff in `CaptureService`, Ports, Broadcaster oder Adapter der Zustellwege | die Rohform ist nicht mehr vorhanden; eine Regeländerung wirkt nur auf künftige Changes; ein Konfigurationsfehler ist für bereits gespeicherte Changes nicht rückholbar |

Festlegung: D. `schema_version` einer Change referenziert weiterhin die
Struktur der Quelltabelle (`LH-FA-SCH-005`); die Schlüssel des Images
folgen zusätzlich dem Regelstand zum Erfassungszeitpunkt. Die Domäne
trägt die Regeltypen und ihre Auswertung als reine Funktion
(`internal/domain/model`, `ARC-001`); der `Assembler` hält den Regelstand
als Laufzeit-Cache der Bindung neben `ExcludedColumns`, ersetzt ihn wie
diese unter `tablesMu` als unveränderliche Liste und erhält ihn beim
`AddBinding`-Merge und bei `setSchemaVersion`. Jeder künftige Pfad, der
`model.Change`-Werte erzeugt, trägt dieselbe Auswertung (Folgepflicht 7).

**Dauerhaftigkeit:** Der Regelstand hat denselben Träger-Typ wie der
Ausschlussstand (`ADR-0065`): die `applied`-Zeilen der beiden neuen
Antragsarten in `cdc.administration_request`, in der Reihenfolge
`requested_at`, bei gleichem Zeitstempel nach `administration_request_id`
ausgewertet; `set_transformation` trägt die Regel ein,
`remove_transformation` nimmt sie heraus. Jeder Pfad, der eine Bindung
anlegt (Prozessstart, Aktivierungs-Zweig), trägt den abgeleiteten Stand
mit. Die Antrags-Zeilen sind dadurch tragend. Der aktive Regelstand ist
über diese Tabelle lesbar; eine eigene Sicht ist nicht Teil dieser
Entscheidung.

### Teilfrage 7 — Routing (`LH-FA-CFG-008`): eigener Folgezug

| Option | Pro | Contra |
|---|---|---|
| A — Routing im selben Zug | Transformation und Routing in einem Zug | verlangt ein Zielmodell, das nicht existiert (Kontext Punkt 5): entweder eine Änderung des `Accepted`-Subjekt-Schemas von `SPEC-024` (`ADR-0100`) oder Filterparameter an gRPC/SSE (`SPEC-020`); die Änderung von Zustellwegen mehrerer Sprach-SDKs kommt hinzu; blockiert die Transformation |
| **B — Transformation jetzt, Routing als eigene Anforderung mit eigener Folge-ADR (gewählt)** | `LH-FA-CFG-007` ist ohne Zielmodell vollständig erfüllbar; Routing bekommt eine Entscheidung mit eigener Abwägung, statt hier geraten zu werden | `LH-FA-CFG-008` bleibt bis zur Routing-ADR ohne Umsetzung und ist in `make doc-trace` als Waise sichtbar |
| C — Routing nicht umsetzen | kein Aufwand | die Anforderung `LH-FA-CFG-008` bliebe unerfüllt ohne Entscheidung darüber |

Festlegung: B. Das Boundary-Kriterium von `LH-FA-CFG-007` wird durch
Teilfrage 3 beantwortet; das von `LH-FA-CFG-008` bleibt dem Routing-Zug
vorbehalten. Als Leitplanken für die Folge-ADR gelten: höchstens ein Ziel
je Change, Auswertung in definierter Reihenfolge (erster Treffer), keine
Route führt auf das bestehende Ziel — alles Weitere entscheidet sie.

### Teilfrage 8 — Auswirkung auf SDKs und Nachrichtenschema

Festlegung: **keine Änderung des Nachrichtenschemas.** Die Nachrichten von
`SPEC-020`/`SPEC-021`/`SPEC-024`/`SPEC-022` tragen dieselben zehn Felder;
die Row Images bleiben JSON-Objekte mit String-Werten, nur ihre
Schlüsselmenge folgt dem Regelstand. Die drei SDK-Packages
(`SPEC-026`/`SPEC-027`/`SPEC-028`) brauchen deshalb erwartungsgemäß keine
Code-Änderung; das zu belegen ist Teil der Folgepflicht (6).

## Konsequenzen

- Positiv: `LH-FA-CFG-007` bekommt für Transformationen einen
  vollständigen Umsetzungspfad ohne neue Sprache, ohne zweiten
  Konfigurationsweg und ohne Eingriff in `CaptureService`, Ports oder
  Zustellwege; eine Auswertungsstelle statt fünf; Live-Reload, Rückkanal
  und Neustart-Festigkeit kommen aus den bestehenden Mechanismen
  (`ADR-0050`/`ADR-0065`); keine achte Fehlerklasse.
- Positiv: Mehrdeutigkeit ist ausgeschlossen statt aufgelöst — es gibt
  keinen datenabhängigen Laufzeit-Konflikt, dessen Auflösung ein Betreiber
  nachvollziehen müsste.
- Negativ: die Rohform geht verloren. Eine Regeländerung wirkt nicht
  rückwirkend; eine fehlerhafte Regel hinterlässt fehlerhaft geformte
  gespeicherte Changes, und ein Consumer erkennt am einzelnen Change nicht,
  unter welchem Regelstand er entstand (`schema_version` bleibt die
  Quellstruktur). `rename_column` verliert keinen Wert; `map_value` verliert
  Information, sobald mehrere Quellwerte auf denselben Zielwert abgebildet
  werden, und ist dann nicht umkehrbar. Eine Rekonstruktion der Rohform je
  Change ist nicht zugesagt: die Antrags-Zeilen tragen `requested_at`, ordnen
  einer Change aber keinen Regelstand zu. Korrekturweg ist
  `remove_transformation` und ein neuer Antrag; bereits gespeicherte Changes
  behalten ihre Form. K1–K4 fangen formale Fehler beim Antrag, keinen
  inhaltlich falschen Wertabgleich. Für Sicherheits-Ausschluss ist der
  Verlust die gewollte Eigenschaft (`ADR-0059` Teilfrage 3 Option D), für
  Umbenennung und Wertabbildung ein Preis, den der Nutzer mit dieser ADR
  annimmt.
- Negativ: die Erfassung der gesamten Quelle steht, wenn eine Regel
  unanwendbar wird (Teilfrage 4), bis die Abhilfe verarbeitet ist; die
  Abhilfe im gescheiterten Prozess ist ein Akzeptanzkriterium (Folgepflicht 5).
- Negativ: das Pflichtenheft (`SPEC-019`, §4 Fehlerklassen-Zeile `schema`)
  und der Antrags-Schema-Rollout (`tools/schema/`) wachsen; die
  Antrags-Tabelle bekommt zwei nullable Spalten.
- Erwartet, nicht am Code belegt: die Abhilfe wirkt in einem an
  `ErrTransformationNotApplicable` gescheiterten Prozess, weil
  `runAdministration` vor `stream.Run` startet und jeder Durchlauf mit
  `processAdministrationRequests` beginnt (`wiring.go`); eine
  Reihenfolge-Garantie folgt daraus nicht.

## Folgepflichten

Umsetzende Arbeitspakete, in dieser Reihenfolge schneidbar; der Planner
formt daraus Welle und Slices. Jedes Paket ist für sich lauffähig und
belegbar.

1. **Spec-Nachzug (Pflichtenheft).** `SPEC-019` um die beiden
   Antragsarten, die Spalten `rule_name`/`rule_spec`, die Regeltypen
   (Teilfrage 2), K1–K4 und die Fehlertexte erweitern; die Zeile
   `schema` der Fehlerklassen-Tabelle um die Nichtanwendbarkeit;
   `LH-FA-CFG-007.a` auf den beantworteten Stand ziehen; Routing steht als
   offene Adresse `LH-FA-CFG-008.a`. Voraussetzung ist die Lastenheft-Version,
   die `LH-FA-CFG-007` auf Transformationen begrenzt und `LH-FA-CFG-008`
   (Routing) führt — sie steht vor `Accepted` dieser ADR in einem eigenen
   Commit. Reihenfolge: zuerst.
2. **Kern: Domäne und `Assembler`, Regeltyp `rename_column`.**
   `model`-Regeltyp mit Konstruktor-Invarianten und reiner
   Auswertungsfunktion; `TableBinding.Transformations`;
   `Assembler.SetTransformation`/`RemoveTransformation`; Auswertung in
   `rowImage` nach dem Ausschluss; Nichtanwendbarkeits-Prüfung gegen
   `event.Relation.Columns`; `ErrTransformationNotApplicable` und
   `classifyRunError`. Belegt durch `make test` (Unit- und
   Nebenläufigkeitstests wie bei `ExcludeColumn`).
3. **Antragsweg und Dauerhaftigkeit.** Schema-Rollout (Antragsart-CHECK,
   zwei Spalten, zwei SQL-Funktionen, `GRANT` nur `cdc_admin`),
   `AdministrationRequestKind`-Erweiterung, Use Cases
   `SetTransformation`/`RemoveTransformation` mit K1–K4, Outbound Port für
   Regelstand-Ableitung und Spaltenliste, Verdrahtung in
   `applyAdministrationRequest`, `activatedTableBindings` und
   Aktivierungs-Zweig. Belegt durch `make test-store`.
4. **Regeltyp `map_value`.** Domäne, Spec-Zeile, Tests — ohne Änderung an
   Antragsweg oder Wirkort. Jeder weitere Regeltyp danach braucht eine
   Folge-ADR.
5. **E2E-Belege in `make test-integration`.** Happy Path am laufenden
   Feed-Container (Regel per SQL beantragt, `applied`, danach erfasste
   Change über `cdc.changes` und die Zustellwege mit umbenanntem
   Schlüssel), Neustart-Festigkeit, Ausschluss+Regel (`LH-QA-SEC-004`: weder
   Quell- noch Zielname noch Wert im Image) und die Nichtanwendbarkeit.
   **Akzeptanzkriterium der Abhilfe** (zu belegen durch den E2E-Lauf; der
   Slice schließt erst, wenn es belegt ist): (a) eine Regel wird
   nichtanwendbar, der Prozess endet sichtbar mit Fehlerklasse `schema`
   (`diagnose`, Heartbeat); (b) `cdc.remove_transformation` wird beantragt,
   während der Prozess steht (Antrag bleibt `requested`); (c) nach dem
   Neustart ist der Antrag `applied`, **bevor** die erste Transaktion der
   Tabelle assembliert wird; (d) die zuvor nicht bestätigte Transaktion
   erscheint danach über `cdc.changes` — kein Datenverlust, keine zweite
   Neustart-Schleife. Trägt die heutige Startreihenfolge (`runAdministration`
   startet vor `stream.Run`, ohne Reihenfolge-Garantie) (c) nicht, gehört
   „offene Anträge vor `stream.Run` verarbeiten" in denselben Slice.
6. **SDK-/Doku-Beleg.** Bestätigen, dass die SDK-Tests keine
   Image-Schlüssel voraussetzen (erwartet: keine Code-Änderung);
   Betriebsdokumentation der SQL-Funktionen unter `docs/user/`.
7. **Bindung künftiger Erzeugungspfade.** Ein Backfill-Pfad
   (`LH-FA-CAP-009.a`) oder ein anderer Change-Erzeuger trägt dieselbe
   Regelauswertung; die zugehörige ADR nennt sie.
8. **Routing-ADR** (eigene Entscheidung, Teilfrage 7) — Vorbedingung, damit
   `LH-FA-CFG-008` belegt werden kann.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Eigenschaftstest im `mapper`-Paket) | für jeden Regeltyp × jede ausgeschlossene Spalte trägt das Image weder den Quellschlüssel noch einen Zielnamen noch den Quellwert (`LH-QA-SEC-004`) | `make test` |
| Go-Test | Auswertung ist deterministisch: gleiche Regelmenge und Relation → byte-gleiches Image; die Regelliste eines Lesers bleibt unter gleichzeitigem `SetTransformation` stabil (`-race`) | `make test` |
| a-check | die Regeltypen liegen in `internal/domain/**` und importieren aus keiner anderen Schicht (Kanten in `.a-check.yml` führen nur nach innen) | `make a-check` |
| Realer Rundlauf | Regel wirkt am laufenden Prozess, nach Neustart und auf allen Zustellwegen mit derselben Form; Nichtanwendbarkeit endet sichtbar in `diagnose` | `make test-integration` |

## Re-Evaluierungs-Trigger

- Ein Consumer braucht eine abweichende Form je Consumer oder Zustellweg,
  oder die Rohform wird gebraucht (Reprocessing, Backfill-Entscheidung zu
  `LH-FA-CAP-009.a`): Folge-ADR mit `Supersedes` für Teilfrage 6.
- Ein dritter Regeltyp wird verlangt (im Beobachtungs-Register dreimal
  auf denselben Bedarf getroffen): Folge-ADR, die den geschlossenen Satz
  erweitert.
- Die Routing-ADR (Teilfrage 7) ändert Leitplanken oder Zustellwege.
- Die Abhilfe nach `ErrTransformationNotApplicable` wirkt real nicht ohne
  Eingriff in die Startreihenfolge (Folgepflicht 5, Kriterium (c)).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-23 | Proposed | Architect-Lauf zu `LH-FA-CFG-007.a` |
| 2026-09-23 | Proposed, überarbeitet: Regelsatz auf `rename_column`/`map_value`, Rohform-Konsequenz, Abhilfe-Kriterium, Routing als `LH-FA-CFG-008` | Auftraggeber-Rückmeldung |
| 2026-09-23 | Accepted — Annahme durch den Auftraggeber samt Rohform-Konsequenz, Regelsatz und Abhilfe-Kriterium | [`LH-FA-CFG-007.a`](../../../spec/pflichtenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
