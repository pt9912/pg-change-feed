# ADR-0137: Routing auf Zustellziele — ein Zustellziel ist ein benannter Kanal; das Ziel wird bei der Erfassung als Label an der Change persistiert

**Status:** Accepted

**Datum:** 2026-10-01

**Herkunft der Festlegungen:** Entscheid des Auftraggebers vom 2026-10-01
(„Alles wie empfohlen“ auf die sechs Fragen, nach Vorlage der Konsequenzen und
Vorschläge in der Sitzung); die Entscheidungen stehen in §Entscheidungen des
Auftraggebers.

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-10-01)

**Bezug:** [`LH-FA-CFG-008`](../../../spec/lastenheft.md) (Haupt-Bezug —
Routing von Changes auf Zustellziele; Out-of-Scope: keine vollständige
Routing-Sprache, Ausdrucksform und Zielmodell sind ADR-/Spec-Frage),
[`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Transformationen — die
Umformung des Inhalts bleibt dort),
[`LH-FA-ADM-003`](../../../spec/lastenheft.md) (sichtbare Fehlerzustände),
[`LH-FA-SST-006`](../../../spec/lastenheft.md) (Gleichwertigkeit der
Zugriffswege),
[`LH-FA-REA-005`](../../../spec/lastenheft.md) (erneutes Lesen innerhalb der
Aufbewahrung),
[`LH-FA-DAT-006`](../../../spec/lastenheft.md) (Erweiterbarkeit der
Change-Metadaten),
[`LH-QA-SEC-004`](../../../spec/lastenheft.md) (ausgeschlossene Spaltenwerte
erscheinen nicht in den Changes),
[ADR-0112](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Teilfrage 7 benennt dieses Routing als Folge-ADR und gibt drei
Leitplanken; Muster für Antragsweg, Wirkort, Fehlerpfad),
[ADR-0011](0011-persist-before-ack.md) und
[ADR-0012](0012-at-least-once.md) (Persist-before-ACK, At-Least-Once),
[ADR-0013](0013-consumer-domainkonzept.md) (Consumer mit eigener Position),
[ADR-0050](0050-sql-administration-antragsqueue-und-live-reload.md) und
[ADR-0130](0130-grpc-verwaltungs-api-neun-rpcs.md) (Antragsweg),
[ADR-0055](0055-nats-change-notification-wecksignal.md) (Wecksignal),
[ADR-0060](0060-grpc-streaming-mechanismus.md),
[ADR-0061](0061-http-sse-zusaetzlich-zu-grpc.md),
[ADR-0100](0100-nats-dritter-vollinhalts-zustellweg.md) (Zustellwege),
[ADR-0133](0133-tabellen-granulare-filterung-grpc-sse.md) (Filterparameter
`schema`/`table` — Formvorbild für den Parameter `target`),
[ADR-0057](0057-http-grpc-api.md) und
[ADR-0081](0081-changes-lesen-ueber-die-http-api.md) (HTTP),
[ADR-0111](0111-backfill-bestand-snapshot-bulk-copy.md) (Spalte `origin`
als Präzedenz einer nullablen Metadaten-Spalte an `cdc.change`),
[ADR-0114](0114-schema-rollout-vorlauf-view-signatur.md) (View-Signatur-
Vorlauf), [ADR-0023](0023-fehlerklassifikation.md) (sieben Fehlerklassen),
[ADR-0046](0046-sql-driving-adapter-lese-schreib-trennung.md) (keine
Domänenlogik in SQL)

**Schärft:** [`LH-FA-CFG-008.a`](../../../spec/pflichtenheft.md)
(Modell der Zustellziele, Konfigurationsmechanismus, Ausdrucksform und
Auflösung), [`SPEC-019`](../../../spec/pflichtenheft.md) (Antrags-Datensatz —
zwei weitere Antragsarten), [`SPEC-020`](../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../spec/pflichtenheft.md),
[`SPEC-022`](../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../spec/pflichtenheft.md) (je ein Filterparameter bzw. ein
Zusatz-Subjekt)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-CFG-008`](../../../spec/lastenheft.md) fordert, erfasste Changes
anhand von Herkunft oder Inhalt auf unterschiedliche Zustellziele zu routen.
Die Anforderung lässt drei Dinge offen: das Modell der Zustellziele, die
Ausdrucksform der Regeln und die Auflösung bei Mehrdeutigkeit. Sie ist die
einzige Waise im RTM-Lauf (`make doc-trace`, Messung der Beobachtung
`arbeit-ueberholt-stehenden-traeger`, Beleg `slice-transformationen-e2e-wirkung`:
80 Anforderungen, 1 Waise).

Ausgangslage, **übernommen** aus [ADR-0112](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Kontext und [ADR-0133](0133-tabellen-granulare-filterung-grpc-sse.md)
Kontext (nicht neu am Code gemessen):

1. **Es gibt kein serverseitiges Zielmodell.** gRPC, SSE und
   NATS-Vollinhalt liefern jedem Abonnenten alle Changes aus einem
   gemeinsamen `Broadcaster`; der Consumer wählt selbst per `schema`/`table`
   (gRPC/SSE, `ADR-0133`; `GET /changes`, `ADR-0081`) oder per
   NATS-Subjekt `cdc.stream.<source_id>.<schema>.<table>` (`SPEC-024`). Keine
   Regel des Betreibers lenkt eine Change an ein bestimmtes Ziel.
2. **Ein Erzeugungspfad, ein Wirkort.** `Assembler.change` erzeugt jedes
   `model.Change`; `CaptureService` persistiert und veröffentlicht danach;
   alle Leser (SQL-Sicht `cdc.changes`, `GET /changes`, drei Live-Wege) sehen
   dasselbe persistierte Bild. Ein Backfill-Pfad trägt dieselbe Auswertung
   (`ADR-0112` Folgepflicht 7).
3. **Der Lesezugriff über SQL kann nicht rechnen** (`ADR-0046`); ein
   Routing nur auf der Zustellseite erreicht ihn nicht.
4. **Eine nullable Metadaten-Spalte an `cdc.change` ist erprobte Praxis**
   (`origin`, `ADR-0111`; Rollout mit View-Vorlauf `ADR-0114`).
5. **Die Leitplanken aus `ADR-0112` Teilfrage 7:** höchstens ein Ziel je
   Change, Auswertung in definierter Reihenfolge (erster Treffer), keine Route
   führt auf das bestehende Ziel.

Diese ADR schreibt der **Architect**, ohne dass ein konkreter Nutzer des
Routings benannt wäre; die Anforderung kommt aus dem Lastenheft, nicht aus
einem beobachteten Betriebsbedarf. Die sechs offenen Fragen hat der
Auftraggeber am 2026-10-01 entschieden (§Entscheidungen des Auftraggebers);
die gewählte Form ist die kleinste, die alle drei Akzeptanzkriterien trägt,
und lässt die teureren Lesarten (aktive Senken, Consumer-Bindung,
Mehrfach-Ziele) als benannte Re-Evaluierung offen.

## Entscheidung

Wir wählen: **Ein Zustellziel ist ein benannter Kanal je Quelle. Das Ziel
einer Change wird bei der Erfassung im `Assembler` aus einer geordneten
Liste einfacher Regeln je Tabelle bestimmt, als nullables Label
`route_target` an der Change persistiert (ein Log, keine Fragmentierung) und
von allen Lesewegen über einen optionalen Parameter `target` ausgewählt.**
Konfiguriert wird über zwei neue Antragsarten der Antrags-Queue. Sieben
Teilfragen.

### Teilfrage 1 — Was ist ein Zustellziel?

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | kein Aufwand | verletzt eine aktive, abnahmebindende Anforderung; Waise im RTM |
| B — Ziel = benannter **Consumer** (`ADR-0013`): ein Consumer trägt eine Regelmenge und sieht nur passende Changes | Consumer hat Position, Ack, Retention-Bezug; at-least-once ist vorhanden | die Live-Wege (gRPC, SSE, NATS) kennen keine Consumer-Identität — Fire-and-Forget ohne Position (`SPEC-024`); ein Stream-Abonnent müsste sich identifizieren (neuer Authentifizierungs-/Identitätsweg, Änderung dreier `Accepted` Zustellwege); koppelt Routing an Consumer-Lebenszyklus (Registrierung) |
| **C — Ziel = benannter Kanal (Label) je Quelle; Zustellung = Abruf oder Abonnement des Kanals über jeden vorhandenen Weg (gewählt)** | passt auf alle Wege mit demselben Parameter (SQL, `GET /changes`, gRPC, SSE) bzw. einem Zusatz-Subjekt (NATS); Consumer-Positionen bleiben Log-Positionen; kein neuer Identitäts- oder Zustellmechanismus | „Zustellung" ist Bereitstellung zum Abruf, kein aktives Schieben an externe Senken |
| D — Ziel = benannte **Senke** mit aktiver Zustellung (Webhook, Queue) | entspricht der Lesart „an ein Ziel zustellen" am nächsten | neue Zustell-Infrastruktur (Wiederholung, Ack, Dead-Letter, Secrets) ohne benannten Bedarf; sprengt die Out-of-Scope-Linie „keine vollständige Routing-Sprache" im Geist; eigener ADR-Umfang |
| E — Ziel nur als NATS-Subjekt | NATS hat Subjekte bereits | HTTP, gRPC, SSE und SQL bleiben ohne Ziel; verletzt `LH-FA-SST-006` (Gleichwertigkeit) |

Festlegung: C. Der Name eines Ziels folgt `[a-z0-9][a-z0-9_-]{0,62}`; das
Alphabet enthält weder Punkt noch NATS-Wildcards, damit der Name ein einzelnes
Subjekt-Token sein kann (*hergeleitet* aus der Subjekt-Syntax, nicht gegen
einen NATS-Server geprüft). Eine Change trägt **höchstens ein** Ziel
(Leitplanke aus `ADR-0112`); „keine Route führt auf das bestehende Ziel" heißt
hier: es gibt kein benanntes Standardziel — eine Change ohne Regeltreffer hat
`route_target = NULL` und bleibt nur über die ungefilterten Wege sichtbar
(Teilfrage 4).

### Teilfrage 2 — Ort der Auswertung

| Option | Pro | Contra |
|---|---|---|
| A — pro Zustellweg nach der Persistierung | Regeländerung wirkt rückwirkend | vier bis fünf Auswertungsstellen; SQL-Sicht kann Inhaltsregeln nicht rechnen (`ADR-0046`); eine Change hätte live ein anderes Ziel als beim erneuten Lesen, sobald sich Regeln ändern (`LH-FA-REA-005`) |
| B — getrennte physische Logs je Ziel | strikte Trennung | zerschneidet das Log: Replay-Invariante (Log ab Anfang ergibt den Quellstand), globale Ordnung, Consumer-Positionen und Retention (`ADR-0014`) bräuchten Neuentwurf; verletzt den Kern von `ADR-0011` |
| **C — im `Assembler`, vor der Persistierung, als Label an der Change (gewählt)** | ein Auswertungsort für alle Wege (wie `ADR-0112` Teilfrage 6); ein Log, ein `change_id`-Raum; live und beim erneuten Lesen identisch; Persist-before-ACK unberührt: die Auswertung erzeugt kein Verwerfen, nur ein Label; Fehler enden vor dem ACK | das Label ist zum Erfassungszeitpunkt festgeschrieben; eine Regeländerung wirkt nur auf künftige Changes (Konsequenzen) |

Festlegung: C. Die Domäne trägt die Regel und ihre Auswertung als reine
Funktion; `model.Change` erhält das Feld `RouteTarget` (leer = NULL);
`cdc.change` erhält die nullable Spalte `route_target` (kein DEFAULT, kein
CHECK, dieselbe Form wie `origin`, `ADR-0111`); die View `cdc.changes`
erhält sie als letzte Spalte (View-Signatur-Vorlauf `ADR-0114`).
Die Auswertung geschieht je Change, nicht je Transaktion. Zusatz-Spalten an
`cdc.change` sind laut `LH-FA-DAT-006` vorgesehen.

### Teilfrage 3 — Regelform (minimal)

Zwei neue Antragsarten `set_route` und `remove_route` auf der Antrags-Queue
(`ADR-0050`), Parameter in den bestehenden nullablen Spalten `rule_name` und
`rule_spec`; SQL-Funktionen `cdc.set_route(source_id, schema_name,
table_name, rule_name, rule_spec jsonb)` und `cdc.remove_route(source_id,
schema_name, table_name, rule_name)`, nur der Rolle `cdc_admin` ausführbar
(`LH-QA-SEC-002`). Die Validierung liegt in Go (`ADR-0046`).

`rule_spec` ist ein JSON-Objekt, unbekannte Schlüssel enden `failed`:

| Schlüssel | Pflicht | Bedeutung |
|---|---|---|
| `target` | ja | Name des Ziels (Alphabet Teilfrage 1) |
| `order` | ja | positive ganze Zahl; kleinere `order` wird zuerst geprüft |
| `when` | nein | Objekt `column`, `equals` (Zeichenkette): Bedingung auf einen Spaltenwert der Quelle |

Herkunft = die Tabelle der Antragszeile (Regel ohne `when` lenkt alle
Changes der Tabelle); Inhalt = `when` (Gleichheit des Textwerts, wie das Row
Image ihn trägt, `ADR-0115`). Weitere Operatoren (Präfix, Menge, Vergleich)
sind nicht Teil des Satzes; ein weiterer braucht eine Folge-Entscheidung.
Die Bedingung liest den **Quellwert vor jeder Transformation** (vor
`rename_column`/`map_value`, `ADR-0112`); Bildbasis ist das Neu-Bild bei
INSERT und UPDATE, das Alt-Bild bei DELETE.

**Auflösung bei Mehrdeutigkeit.** Mehrdeutigkeit ist über `order`
vollständig geordnet (`LH-FA-CFG-008` Boundary: Regel-Reihenfolge); der
erste Treffer in aufsteigender `order` bestimmt das Ziel. Zusätzlich werden
Fälle, in denen eine Regel nie wirken könnte, statisch ausgeschlossen; ein
Antrag, der eine Invariante verletzte, endet `failed` mit Fehlertext, der
Regelstand bleibt:

- **R1** — `rule_name` ist je Tabelle unter den Routing-Regeln eindeutig
  (eigener Namensraum, getrennt von Transformationsregeln).
- **R2** — `order` ist je Tabelle eindeutig.
- **R3** — `when.column` existiert an der Quelle und ist keine
  ausgeschlossene Spalte; umgekehrt endet `exclude_column` gegen eine Spalte,
  die eine Routing-Bedingung trägt, `failed`. Grund: das Ziel verriete
  sonst eine Eigenschaft des ausgeschlossenen Werts (`LH-QA-SEC-004`).
- **R4** — höchstens eine Regel ohne `when` je Tabelle, und sie trägt die
  höchste `order` der Tabelle (sonst wären alle Regeln dahinter wirkungslos).
- **R5** — das Paar (`column`, `equals`) kommt je Tabelle höchstens einmal
  vor.
- **R6** — `remove_route` gegen einen nicht geführten Namen endet `failed`.

### Teilfrage 4 — Nichtanwendbarkeit, kein Treffer, Sichtbarkeit

Eine Regel ist auf eine Change **nicht anwendbar**, wenn `when.column` in der
Relation der Change fehlt (z. B. nach einer Schemaänderung). Dann meldet der
`Assembler` `ErrRoutingNotApplicable`, den `classifyRunError` auf die
Fehlerklasse `schema` abbildet (keine achte Klasse, `ADR-0023`): die
Transaktion wird weder persistiert noch bestätigt, Heartbeat und `diagnose`
zeigen den Zustand (`LH-FA-ADM-003`); Abhilfe ist `cdc.remove_route` und
Neustart — derselbe Pfad wie `ADR-0112` Teilfrage 4 (Startreihenfolge dort
als Erwartung geführt; sie wird hier **nicht** als belegt vorausgesetzt, der
E2E-Slice belegt sie für diese Antragsart erneut).

Zwei Fälle sind **keine** Nichtanwendbarkeit, weil die Regel anwendbar ist
und verneint:

- Der Wert ist im Bild abwesend (NULL, unverändertes TOAST, bei DELETE ohne
  volle Replica-Identität jede Nicht-Schlüsselspalte): `when` trifft nicht
  zu, die nächste Regel wird geprüft.
- Keine Regel trifft: `route_target = NULL`. Die Change geht **nicht**
  verloren und **nicht** an ein Standardziel; sie ist im Log, in der SQL-Sicht
  mit `route_target IS NULL` und über jeden ungefilterten Weg sichtbar.

Beide Fälle sind eine Auslegung von „Negative", die der Auftraggeber
entschieden hat (Entscheidungen 3 und 4). Das Handbuch weist auf DELETE und
Inhaltsregeln hin (Folgepflicht 8).

### Teilfrage 5 — Zustellung: ein Parameter je Weg

| Weg | Form |
|---|---|
| SQL-Sicht `cdc.changes` | Spalte `route_target`; Filter `WHERE route_target = '<ziel>'` |
| `GET /changes` (`SPEC-022`) | optionaler Query-Parameter `target`; unbekannte Parameter bleiben `400` |
| gRPC (`SPEC-020`) | `StreamChangesRequest` Feld 3 `string target`; leer = kein Filter (wie `schema`/`table`, `ADR-0133` Teilfrage 2) |
| SSE (`SPEC-021`) | Query-Parameter `target`; derselbe `400`-Pfad |
| NATS-Vollinhalt (`SPEC-024`) | zusätzliche Veröffentlichung jeder Change mit `route_target` auf `cdc.route.<source_id>.<ziel>` (Payload wie `cdc.stream…`); `cdc.stream.<source_id>.<schema>.<table>` bleibt unverändert |
| NATS-Wecksignal (`SPEC-017`) | unverändert: leerer Payload, Tabellen-Kardinalität, kein Ziel |

`target` kombiniert sich mit `schema`/`table` als Konjunktion. Die
Filterprüfung ist wie in `ADR-0133` eine reine Funktion im jeweiligen
Driving-Handler; `ChangeStreamPort` und `Broadcaster` bleiben unverändert.
Das Nachrichtenschema der Wege (zehn Felder) bleibt unverändert: das Label ist
über die SQL-Sicht lesbar, nicht Teil der Stream-Nachrichten. Ein
ungefilterter Leser sieht **weiterhin alle** Changes, geroutete eingeschlossen
(Auswahl, keine Ausblendung; Entscheidung 2).

**Consumer.** Positionen bleiben Log-Positionen (`ADR-0013`). Ein Consumer, der
nur ein Ziel liest, bestätigt die höchste gesehene `change_id`; die
Changes anderer Ziele davor überspringt er verlustfrei bezüglich seines
Ziels. Die Retention bleibt am niedrigsten Consumer-Stand (*hergeleitet*, am
Code nicht geprüft). Eine feste Bindung Consumer ↔ Ziel ist nicht Teil
dieser Entscheidung.

### Teilfrage 6 — Verhältnis zu Transformationen, Spaltenausschluss, Backfill

- **Transformation** (`LH-FA-CFG-007`, `ADR-0112`): Die Routing-Bedingung
  liest den Quellwert vor der Transformation; das Ziel ist Metadatum der
  Change und ändert das Bild nicht. Die Transformation ändert das Ziel nicht.
  Beide Regelwerke sind unabhängig, getrennte Namensräume und Stände.
- **Spaltenausschluss:** R3 (gegenseitige Sperre) hält ausgeschlossene
  Werte aus dem Routing heraus.
- **Backfill** (`LH-FA-CAP-009`): Backfill-Changes durchlaufen dieselbe
  Auswertung (Bindung künftiger Erzeugungspfade, Folgepflicht 6) und erhalten
  das Label des Regelstands zum Run; die Replay-Invariante ist unberührt,
  weil das Label nicht zum Zeilenzustand gehört.
- **Dauerhaftigkeit:** Regelstand aus den `applied`-Zeilen der beiden
  Antragsarten, Ordnung `requested_at`, dann `administration_request_id`
  (Muster `ADR-0065`/`ADR-0112`), in jedem Pfad, der eine Bindung anlegt.
- **Fehlerklasse:** `schema` (Teilfrage 4).

### Teilfrage 7 — Zerlegung: eigener Zug, siehe Folgepflichten

Die Umsetzung ist größer als `ADR-0133` (ein Slice) und kleiner als die
Transformationen; Schnitt in den Folgepflichten.

## Verglichene Alternativen

Die Optionstabellen stehen je Teilfrage (Teilfrage 1 Zielmodell, Teilfrage 2
Ort der Auswertung); „nichts tun" ist jeweils Option A. Die Regelform
(Teilfrage 3) wurde gegen den Satz „freier Ausdruck/Skript" (vom Lastenheft
ausgeschlossen, wie `ADR-0112` Teilfrage 2) und gegen eine Konfiguration
über die YAML-Datei (Boot-Zeit, kein Rückkanal, `ADR-0112` Teilfrage 1 B)
verworfen.

## Konsequenzen

- Positiv: `LH-FA-CFG-008` bekommt einen vollständigen Umsetzungspfad ohne
  neue Sprache, ohne neuen Zustellmechanismus und ohne Eingriff in
  `CaptureService`, Ports oder `Broadcaster`. Ein Log, ein Auswertungsort,
  gleiche Sicht live und beim erneuten Lesen.
- Positiv: Mehrdeutigkeit ist vollständig geordnet (`order`), tote Regeln sind
  statisch ausgeschlossen, Nichtanwendbarkeit ist sichtbar (`schema`).
- Negativ: das Label ist zum Erfassungszeitpunkt fest. Eine Regeländerung wirkt
  nicht rückwirkend; ein neues Ziel sieht ältere Changes nicht unter seinem
  Namen. Es gibt kein Umetikettieren in dieser Entscheidung.
- Negativ: „Zustellung" ist Abruf/Abonnement des Kanals, kein Schieben an
  externe Senken (Teilfrage 1 D); wer das braucht, braucht eine eigene
  Entscheidung.
- Negativ: eine Change ohne Treffer ist für einen Leser, der nur ein Ziel
  abruft, unsichtbar. Der Betreiber sieht sie über `route_target IS NULL`;
  eine Warnung dafür ist nicht Teil dieser Entscheidung.
- Negativ: Inhaltsregeln auf Nicht-Schlüsselspalten wirken bei DELETE ohne
  volle Replica-Identität nicht (Wert abwesend, Teilfrage 4); *hergeleitet*
  aus dem Abwesenheits-Vertrag von `LH-FA-DAT-005` und `ADR-0112`, die
  Messung gegen PostgreSQL 17 und 18 gehört in den E2E-Slice (Folgepflicht 7);
  vorher gilt die Aussage nicht als belegt.
- Negativ: zweite NATS-Veröffentlichung je gerouteter Change (zusätzliche
  Last am Publisher); Kosten *nicht gemessen*.
- Negativ: Schema-Rollout wächst (Spalte, View-Spalte, zwei Antragsarten,
  zwei Funktionen); SDKs brauchen für die Stream-Wege einen optionalen
  Parameter `target` (je eigener Folge-Schritt wie bei `ADR-0133`).

## Folgepflichten

Vorschlag zur Zerlegung; der Planner formt Welle und Slices. Jeder Slice ist
für sich lauffähig. **Aussage über eine Menge:** die Zahl acht ist ein
Plan-Entwurf, nicht gemessen. **Abhängigkeiten** (*hergeleitet*, nicht gegen
einen Plan geprüft): 1 vor allen; 2 vor 3 bis 6; 3 vor 6 und 7, weil Regeln nur
über den Antragsweg entstehen; 4, 5 und 6 untereinander unabhängig; 7 nach 2
bis 6; 8 nach 7.

1. **Spec-Nachzug (zuerst, eigener Commit):** `SPEC-019` um die Antragsarten
   `set_route`/`remove_route`; ein neuer SPEC-Eintrag für die Routing-Regel
   (Schlüssel, R1–R6, Fehlertexte) neben `SPEC-030`; `LH-FA-CFG-008.a` auf den
   beantworteten Stand; `SPEC-008` Zeile `schema`; `SPEC-020`/`-021`/`-022`/
   `-024` je Parameter bzw. Subjekt; `SPEC-001`/`SPEC-002` Spalte
   `route_target`. `spec/architecture.md` bleibt frei von Slices und
   ADR-Bezügen (`AGENTS.md` §3.4); die Sicht braucht höchstens einen Satz zur
   Metadaten-Spalte, falls `ARC-*` die Change-Felder aufzählt (dort zu prüfen).
2. **Kern:** Domänenregel + reine Auswertung, `model.Change.RouteTarget`,
   `TableBinding`-Routen, Auswertung im `Assembler`, `ErrRoutingNotApplicable`
   + `classifyRunError`; Persistenz der Spalte `route_target`, View-Spalte
   (Rollout mit `ADR-0114`-Vorlauf). Belegt durch `make test`,
   `make test-store`.
3. **Antragsweg und Dauerhaftigkeit:** Antragsarten, SQL-Funktionen, Use Cases
   mit R1–R6, Regelstand-Ableitung, Verdrahtung in
   `applyAdministrationRequest`, Aktivierungs-Zweig, Prozessstart. Belegt
   durch `make test-store`.
4. **Lesewege:** `GET /changes`-Parameter, gRPC-Feld (Proto-Generierung,
   `make generated-sync`), SSE-Parameter, gemeinsame Filterfunktion mit
   `ADR-0133`. Belegt durch `make test`.
5. **NATS-Zusatz-Subjekt** `cdc.route.<source_id>.<ziel>`. Belegt durch
   `make test-notify` bzw. den NATS-Phasen-Test.
6. **Backfill-Pfad** trägt die Auswertung (Bindung künftiger Erzeugungspfade,
   `ADR-0112` Folgepflicht 7).
7. **E2E-Belege in `make test-integration`** (gleichzeitig der Nachweis, der
   `LH-FA-CFG-008` aus den Waisen nimmt, `docs/user/e2e-abdeckung.md`): Happy
   Path (Herkunfts- und Inhaltsregel; Ziel A sieht nur A, ungefilterter Leser
   alles), Boundary (zwei Regeln treffen; `order` entscheidet; R1–R6 je eine
   Verletzung `failed`), Negative (Nichtanwendbarkeit endet `schema`, Abhilfe
   über `remove_route` + Neustart), Neustart-Festigkeit, Ausschluss-Sperre
   (R3), Replay (erneutes Lesen liefert dasselbe Label). Zusätzlich die
   Messung des DELETE-Verhaltens (Inhaltsregel auf eine Nicht-Schlüsselspalte
   ohne volle Replica-Identität) gegen PostgreSQL 17 und 18, bevor die Aussage
   in Teilfrage 4 als belegt gilt (Entscheidung 4).
8. **Betriebsdokumentation und SDK-/Beispiel-Parameter** `target` (eigener
   Folge-Schritt, nicht Gegenstand dieser ADR); das Handbuch weist auf DELETE
   und Inhaltsregeln hin (Entscheidung 4).

## Fitness Function (falls maschinell prüfbar)

Regeln dieser Sektion: Erwartungen an den umsetzenden Slice, **nicht erprobt**
(`AGENTS.md` §3.12; der Architect testet nicht). Keine Mutation ist gefahren.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Domäne, erwartet) | für jede Regelmenge und jede Change ist das Ziel deterministisch; bei zwei treffenden Regeln gewinnt die kleinere `order` | `make test` |
| Go-Test (`mapper`, erwartet) | ein `when` auf eine in der Relation fehlende Spalte liefert `ErrRoutingNotApplicable`; abwesender Wert liefert Nicht-Treffer ohne Fehler | `make test` |
| Go-Test (Eigenschaftstest, erwartet) | kein Zielname und keine Routing-Bedingung nennt eine ausgeschlossene Spalte (R3, beide Richtungen) | `make test` |
| Store-Test (erwartet) | `route_target` überlebt Persistierung und Lesen über `cdc.changes`; NULL liest als „nicht geroutet" | `make test-store` |
| a-check | Routing-Regel liegt in `internal/domain/**` und importiert aus keiner anderen Schicht | `make a-check` |
| Realer Rundlauf (erwartet) | Regel wirkt am laufenden Prozess, nach Neustart und auf allen Lesewegen mit demselben Ergebnis | `make test-integration` |

## Entscheidungen des Auftraggebers

Herkunft: Antwort des Auftraggebers vom 2026-10-01 („Alles wie empfohlen“) auf
die sechs Fragen des Entwurfs, die diese ADR bis dahin als `Proposed` führte.
Der Grund für `Proposed` war, dass die Entscheidung keinen benannten Nutzer
hat; die Anforderung kommt aus dem Lastenheft.

| # | Frage | Entscheidung | Konsequenz |
|---|---|---|---|
| 1 | Bedarf und Form der „Zustellung“ | Ziel = Kanal, den ein Consumer abruft oder abonniert (Option C) | aktive Senken (Option D) und Consumer-Bindung (Option B) bleiben Re-Evaluierung; das Label ist die Grundlage für eine spätere ADR zu aktiven Senken |
| 2 | Exklusivität | geroutete Changes bleiben für ungefilterte Leser sichtbar | `cdc.changes` und die Wege der Altclients ändern sich nicht; Routing ist Auswahl, keine Ausblendung |
| 3 | Kein Treffer | `route_target = NULL`, kein Standardziel | eine Regel ohne `when` mit höchster `order` bleibt als Abschlussregel möglich |
| 4 | Abwesender Wert (NULL, TOAST, DELETE ohne volle Replica-Identität) | Nicht-Treffer | Handbuch-Hinweis zu DELETE und Inhaltsregeln; Messung gegen PostgreSQL 17 und 18 im E2E-Slice |
| 5 | Ein Ziel je Change | ein Ziel genügt | kein Fan-out; mehrere Ziele wären eine Zielmenge statt einer Spalte |
| 6 | Label nicht rückwirkend | bei der Erfassung festgeschrieben, kein Umetikettieren | Altbestand über einen Backfill-Run mit dem aktuellen Regelstand |

## Re-Evaluierungs-Trigger

- Eine der Entscheidungen 1, 2 oder 5 wird anders getroffen (insbesondere
  aktive Senken, Option D): Folge-ADR mit `Supersedes` für die betroffene
  Teilfrage.
- Ein Operator über Gleichheit hinaus wird dreimal im Beobachtungs-Register
  verlangt: Folge-ADR zum Regelsatz.
- Ein Bedarf an Umetikettieren bereits erfasster Changes: Folge-ADR, die den
  Ort der Auswertung (Teilfrage 2) neu bewertet.
- Messung der zweiten NATS-Veröffentlichung zeigt Druck am Publisher.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-01 | Proposed — Architect-Lauf zu `LH-FA-CFG-008.a`; sechs Fragen an den Auftraggeber | [`LH-FA-CFG-008`](../../../spec/lastenheft.md) |
| 2026-10-01 | Accepted — sechs Entscheidungen des Auftraggebers wie im Entwurf | [`LH-FA-CFG-008`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
