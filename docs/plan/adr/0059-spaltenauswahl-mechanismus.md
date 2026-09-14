# ADR-0059: Spaltenauswahl — Mechanismus, Granularität und Wirkort

**Status:** Accepted

**Datum:** 2026-09-14

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-14)

**Bezug:** [`LH-FA-CFG-005`](../../../spec/lastenheft.md) (Haupt-Bezug —
Spaltenauswahl, seit Lastenheft-Version 0.7.0 aktive Anforderung statt
dauerhaftem Ausschluss), [`LH-FA-CFG-001`](../../../spec/lastenheft.md)
(Tabellen-Aktivierung — die Bindung, an die eine Spaltenauswahl anknüpft),
[`LH-FA-SCH-003`](../../../spec/lastenheft.md) (entfernte Spalten — die
eigene Boundary-Klausel von `LH-FA-CFG-005` verweist explizit hierher),
[`LH-FA-DAT-005`](../../../spec/lastenheft.md) (relevante Datenwerte —
trägt bereits die Negative-Abwesenheits-Semantik, auf der diese
Entscheidung aufbaut), [`LH-QA-SEC-004`](../../../spec/lastenheft.md)
(Ausschluss sensibler Spalten — Messmethode verweist auf
`LH-FA-CFG-005`), [ADR-0050](0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue mit LISTEN/NOTIFY + Live-Reload — Erweiterungsbasis dieser
Entscheidung), [ADR-0046](0046-sql-driving-adapter-lese-schreib-trennung.md)
(Lese-/Schreib-Trennung — Präzedenzfall für die Wirkort-Frage),
[ADR-0028](0028-inbound-use-cases.md) (Inbound Use Cases, denen die neuen
Fähigkeiten folgen), [ADR-0034](0034-ports-nach-faehigkeiten.md)
(Port-Zuschnitt — Begründung für einen neuen statt eines erweiterten
Ports), [ADR-0052](0052-optionale-yaml-konfigurationsdatei.md) (optionale
YAML-Konfigurationsdatei — geprüfte, verworfene Alternative für den
Mechanismus)

**Schärft:** [`ARC-005`](../../../spec/architecture.md) (die dort bereits
zwei getrennte SQL-Administrationspfade — CLI-Direktaufruf und
Antragsqueue — zeigende Sequenzsicht bekommt zwei weitere Antragsarten
auf demselben Pfad), [`ARC-002`](../../../spec/architecture.md) (zwei
neue Use Cases), [`ARC-004`](../../../spec/architecture.md) (ein neuer
Outbound Port)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`LH-FA-CFG-005`](../../../spec/lastenheft.md) wurde mit Lastenheft-Version
0.7.0 (Commit `97de1e9`, 2026-09-14) von einer dauerhaft ausgeschlossenen,
perspektivischen Anforderung zu einer aktiven Anforderung umgestellt. Die
Anforderung selbst benennt ihre offenen Fragen wörtlich in ihrem eigenen
Out-of-Scope: „Konfigurationsmechanismus, Konfigurationsgranularität
(Tabelle vs. Quelle) und Zusammenspiel mit der SQL-Administration sind
Architektur- (ADR) bzw. Spezifikationsfragen (`SPEC-*`), keine
Lastenheft-Festlegung." Diese ADR beantwortet genau diese drei Fragen,
zusammen mit zwei daran hängenden technischen Detailfragen (Konvergenz mit
`LH-FA-SCH-003` bei realer Spaltenlöschung, Fehlerpfad für eine
nicht existierende Spalte).

Der Code trägt heute **keinen** Mechanismus für Spaltenausschluss. Die
Bausteine, an die eine Lösung anknüpfen kann, existieren aber bereits:

- **`internal/adapters/driving/replication/mapper/mapper.go`** baut in
  `Assembler.change`/`rowImage` das Row Image jeder Änderung aus dem
  dekodierten `pgoutput`-Ereignis. `rowImage` überspringt bereits einen
  Spaltenschlüssel, wenn der zugehörige Wert `nil` ist (NULL oder
  unverändertes TOAST) — dieselbe Abwesenheits-Kodierung, die
  `LH-FA-DAT-005`s Negative-Akzeptanzkriterium fordert („Abwesenheit
  erkennbar, nicht mit einem Platzhalter überdeckt").
- **`TableBinding`** (`mapper.go`) trägt heute `TableID` und
  `SchemaVersion` je aktivierter Tabelle, indiziert nach qualifiziertem
  Namen in `Assembler.tables`. `Assembler.tables` ist seit `ADR-0050`
  über `tablesMu` (`sync.RWMutex`) synchronisiert, weil zwei Goroutinen
  darauf zugreifen: der Capture-Stream (Lese-Zugriff) und die
  Administrations-Goroutine (`runAdministration`/`applyAdministrationRequest`
  in `internal/bootstrap/wiring.go`, Schreib-Zugriff über
  `AddBinding`/`RemoveBinding`).
- **`ADR-0050`** hat für „schreibende SQL-Administration, die den
  laufenden Prozess ohne Neustart erreichen muss" bereits eine Lösung
  etabliert: eine schreibende SQL-Funktion (`cdc.enable_table(...)`,
  `cdc.disable_table(...)`) schreibt ausschließlich einen
  Antrags-Datensatz in `cdc.administration_request`
  (`tools/schema/schema.yaml`) und sendet `pg_notify`; der laufende
  Capture-Prozess verarbeitet offene Anträge in einer eigenen
  Hintergrund-Goroutine über den zuständigen Inbound Port
  (`EnableTableUseCase`/`DisableTableUseCase`) und trägt bei Erfolg die
  `Assembler`-Bindung im selben Prozess nach (`AddBinding`/`RemoveBinding`)
  — Live-Reload ohne Neustart, ohne zweiten Schreibpfad auf
  `cdc.source_table`/die Publication.
- **`ADR-0046`** hat für den SQL-Driving-Adapter bereits entschieden:
  reine Lese-Views dürfen direkt gegen gespeicherte Tabellen lesen,
  schreibende/aktionsauslösende Zugriffe laufen ausschließlich über
  Inbound Ports — keine Domänenlogik-Duplikation in SQL.
- **`ADR-0052`** hat eine optionale YAML-Konfigurationsdatei
  (`CDC_CONFIG_FILE`, `SPEC-016`) eingeführt — additiv zu den
  Umgebungsvariablen, mit Feld-für-Feld-Precedence, **geladen einmalig
  beim Prozessstart** (`wiring.go`, dasselbe Boot-Zeitfenster wie
  `CDC_TABLES`). Sie trägt heute keinen Laufzeit-Reload-Mechanismus.

`internal/application/port/inbound/verwaltung.go` trägt bereits das
Muster für einen expliziten Negative-Fehlerpfad bei einem nicht
existierenden Ziel: `ErrSourceTableMissing`, geprüft über
`TableActivationPort.TableExists` und ausgelöst von `EnableTableService`/
`DisableTableService` (`LH-FA-CFG-001`/`002` Negative-Akzeptanzkriterium).

## Entscheidung

Wir wählen: **Spaltenausschluss/-einschluss als zwei neue Antragsarten
auf derselben Antrags-Queue aus `ADR-0050`** (SQL-Funktionen
`cdc.exclude_column(schema, table, column)`/`cdc.include_column(schema,
table, column)`), **Granularität ausschließlich pro Tabelle**, **Wirkort
im `Assembler` bei der Row-Image-Konstruktion** (vor jeder Serialisierung,
nicht bei Schreiben oder Lesen der Speichertabelle), mit **Konvergenz**
auf `LH-FA-SCH-003`s bestehendem Verhalten bei realer Spaltenlöschung und
**demselben Fehlerpfad-Muster** wie `cdc.enable_table` für eine nicht
existierende Tabelle. Fünf Teilfragen, fünf Festlegungen:

### Teilfrage 1 — Konfigurationsmechanismus

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (`LH-FA-CFG-005` bleibt ohne Mechanismus) | kein Aufwand | verletzt eine seit Lastenheft-Version 0.7.0 aktive, vertraglich abnahmebindende Anforderung; `LH-QA-SEC-004` bleibt strukturell unbelegbar (ihre Messmethode verweist ausschließlich auf `LH-FA-CFG-005`) |
| B — Umgebungsvariable (`CDC_EXCLUDED_COLUMNS`, Format analog zu `CDC_TABLES`) | kein neuer Mechanismus, kein neuer Code-Pfad im Store | inkonsistent mit dem für die verwandte Fähigkeit „Tabellen-Aktivierung" bereits etablierten SQL-Administrationsweg (`ADR-0046`/`ADR-0050`) — zwei Konfigurationswege für zwei eng verwandte administrative Fähigkeiten; wirkt nur beim Prozessstart, widerspricht `LH-FA-CFG-005`s Happy-Path-Prämisse „CDC ist für `t` **bereits aktiviert**, when Spalte `c` vom Ausschluss konfiguriert wird" (eine bereits laufende Aktivierung reagieren lassen, ohne Neustart) und kann `LH-FA-CFG-005`s Negative-Akzeptanzkriterium (expliziter Fehlerpfad bei nicht existierender Spalte) nicht als Aufrufer-Feedback liefern — eine ENV-Variable hat keinen Rückkanal |
| C — optionale YAML-Konfigurationsdatei (`ADR-0052`/`SPEC-016`) um ein `excluded_columns`-Feld je Tabelle erweitern | nutzt einen bereits etablierten strukturierten Konfigurationsweg, additiv zu ENV | `ADR-0052` ist bewusst auf **Boot-Zeit-Konfiguration** zugeschnitten (Feld-für-Feld-Precedence-Merge beim Prozessstart, kein Reload-Mechanismus) — dieselbe Neustart-Einschränkung wie Option B; führt zusätzlich eine zweite Quelle der Wahrheit neben der DB-gestützten Tabellen-Aktivierung ein (Aktivierungszustand lebt in `cdc.source_table`, Ausschlusszustand läge in einer Datei) für zwei Aspekte derselben Bindung |
| **D — zwei neue Antragsarten auf der Antrags-Queue aus `ADR-0050` (gewählt)** | ein einziger, bereits etablierter Konfigurationsweg für alle administrativen Fähigkeiten, die den laufenden Prozess ohne Neustart erreichen müssen; löst die Neustart-Einschränkung strukturell (derselbe Live-Reload-Pfad wie Tabellen-Aktivierung); erfüllt den Negative-Rückkanal über den bereits etablierten Antrags-Status (`pending`/`applied`/`failed`, `error_message`) | erweitert `cdc.administration_request`s `request_kind`-CHECK-Constraint und (für die Spalte) ihre Spaltenform — Migrationsaufwand, aber kein neuer Mechanismus |

`CDC_TABLES` bzw. die YAML-Datei bleiben unverändert **Erstaktivierungs-
Seed** für eine leere Datenbank (wie in `ADR-0050` für die
Tabellen-Aktivierung selbst entschieden); ein Spaltenausschluss, der von
Anfang an gelten soll, wird nach der Erstaktivierung per
`cdc.exclude_column(...)` beantragt — ein zusätzliches Boot-Feld für
initial ausgeschlossene Spalten ist durch diese Entscheidung nicht
ausgeschlossen, aber nicht ihr Gegenstand.

### Teilfrage 2 — Konfigurationsgranularität

| Option | Pro | Contra |
|---|---|---|
| A — ausschließlich pro Quelle (ein globaler Ausschluss über alle Tabellen hinweg, z. B. nach Spaltennamens-Konvention) | ein einziger Ausschluss deckt eine Sicherheits-Konvention wie „nie eine Spalte `password`" über den gesamten Bestand ab, ohne je Tabelle wiederholt zu werden | widerspricht `LH-FA-CFG-005`s eigener Akzeptanzkriterien-Formulierung wörtlich („Given CDC ist für `t` aktiviert, when Spalte `c` [...]" — eindeutig tabellenbezogen); ein gleichnamiges Feld kann in verschiedenen Tabellen unterschiedliche Sensibilität tragen (eine global fixierte Regel kann nicht selektiv eine Tabelle ausnehmen); eine „Spalte existiert nicht"-Negative-Prüfung (`LH-FA-CFG-005` Negative) ist bei einem quellenweiten Namens-Muster nicht sinnvoll formulierbar (die Spalte muss nicht in jeder Tabelle existieren) |
| **B — ausschließlich pro Tabelle (gewählt)** | deckt `LH-FA-CFG-005`s Akzeptanzkriterien exakt (Tabelle `t`, Spalte `c`); dieselbe Granularität wie die bestehende Tabellen-Bindung (`TableBinding`), kein neues Konzept; die Negative-Prüfung „Spalte existiert nicht" bindet eindeutig an eine konkrete Tabelle; erfüllt `LH-QA-SEC-004`s Messmethode vollständig, da diese ausschließlich auf `LH-FA-CFG-005` verweist | eine Sicherheits-Konvention über viele Tabellen hinweg (z. B. „jede Spalte `password_hash`") verlangt einen Aufruf je Tabelle statt eines einzigen — Betriebsaufwand, keine Korrektheitslücke |
| C — beides: pro Tabelle **und** zusätzlich ein quellenweiter Muster-Ausschluss, der mit der Tabellen-Konfiguration kombiniert wird | deckt `LH-QA-SEC-004`s im Auftrag genanntes Beispiel „niemals irgendeine Spalte namens `password`" unmittelbar ab | verdoppelt den Umfang dieser Entscheidung: zwei Persistenzmodelle, zwei Antragsarten, und eine neue Konfliktregel bei widersprüchlicher Konfiguration (globaler Ausschluss vs. tabellenspezifischer Einschluss — welcher gewinnt?), ohne dass ein Akzeptanzkriterium diesen Umfang verlangt (`LH-QA-SEC-004`s Messmethode prüft ausschließlich über `LH-FA-CFG-005`, die selbst tabellenbezogen formuliert ist) |

Die globale Variante (Option C) bleibt damit **bewusst ausgeschlossen** —
nicht verworfen, sondern vertagt: Erst ein konkret benannter Bedarf (ein
Beobachtungs-Register-Eintrag, der dreimal auf denselben
Wiederholungsaufwand trifft, oder eine geschärfte `LH-QA-SEC-004`) macht
sie zu einer eigenen Folge-ADR-würdigen Entscheidung; bis dahin deckt
Option B jedes vorliegende Akzeptanzkriterium.

### Teilfrage 3 — Wirkort im Row-Image-Aufbau

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (kein Filter-Mechanismus) | kein Aufwand | siehe Teilfrage 1, Option A |
| B — Filterung an der Lese-View `cdc.changes` (`ADR-0046`-konform als reine Projektion) | zentraler Punkt, an dem jeder Consumer liest; kein neuer Go-Code am Capture-Pfad | der sensible Wert wird trotzdem **in `cdc.change.new_data`/`old_data` auf Platte persistiert** — nur die Anzeige über die View wird gefiltert, nicht die Erfassung selbst; ein Direktzugriff auf die Speichertabelle, ein Backup oder ein künftiger zweiter Lesekanal umgeht den Schutz vollständig; `LH-QA-SEC-004` verlangt Ausschluss „von der CDC-**Erfassung**", nicht nur von einer bestimmten Lese-View |
| C — Filterung beim Schreiben nach `cdc.change` (SQL-Trigger/Funktion vor `INSERT`) | verhindert wie D die Persistenz des sensiblen Werts | dupliziert die Auswertung der Ausschlussliste in SQL — ein zweiter Ort neben dem Go-seitigen `TableBinding`-Zustand, der mit ihm synchron gehalten werden müsste, obwohl der Live-Reload-Mechanismus (`ADR-0050`) ausschließlich den Go-Prozess-Zustand aktualisiert; verlangt, dass das volle (ungefilterte) Row Image erst als JSON gebaut und dann von der SQL-Ebene erneut interpretiert und bereinigt wird — unnötiger Zusatzschritt pro Change, den der Assembler mit vollständigem Kontext (Spaltennamen aus der Relation-Nachricht) bereits ohne Umweg leisten kann |
| **D — Filterung im `Assembler` bei der Row-Image-Konstruktion (`rowImage`/`change`) (gewählt)** | der sensible Wert wird **nie** serialisiert und nie an die Persistenzschicht übergeben — `LH-QA-SEC-004` vollständig erfüllt, nicht nur eine Sichtbarkeits-Filterung; ein einziger Auswertungsort für die Ausschlussliste (`TableBinding`, dieselbe Struktur wie die Tabellen-Aktivierung selbst); nutzt denselben Live-Reload-Pfad wie `ADR-0050` (`AddBinding`-Mechanismus, `tablesMu`) ohne neuen Synchronisationsmechanismus; die bereits etablierte Abwesenheits-Kodierung von `rowImage` (ein übersprungener Spaltenschlüssel bei `nil`-Wert) trägt denselben Vertrag für einen ausgeschlossenen Wert — `LH-FA-DAT-005`s Negative-Kriterium („Abwesenheit erkennbar, nicht mit einem Platzhalter überdeckt") ist damit ohne neue Kodierung erfüllt | `TableBinding` und `Assembler` bekommen eine weitere Verantwortung (Filterung neben Übersetzung/Bindung); die Live-Reload-Schreibpfade (`AddBinding` für neue Bindung/Schema-Bump, eine neue Methode für Ausschluss-Änderung) müssen den Ausschlussstand einer bestehenden Bindung beim Nachtragen einer neuen Schema-Version *erhalten*, statt ihn stillschweigend zurückzusetzen — ein Implementierungsdetail, das der umsetzende Slice sauber lösen muss (kein blindes vollständiges Überschreiben von `TableBinding` mehr) |

### Teilfrage 4 — Konvergenz mit `LH-FA-SCH-003` bei realer Spaltenlöschung

`LH-FA-CFG-005`s Boundary-Akzeptanzkriterium sagt wörtlich: „Given ein
Spaltenausschluss ist konfiguriert, when die Spalte gelöscht wird, dann
gilt das Verhalten aus [`LH-FA-SCH-003`](../../../spec/lastenheft.md)" — die Anforderung selbst verlangt
also keinen Sonderfall, sondern ausdrücklich die Anwendung des
bestehenden Verhaltens.

| Option | Pro | Contra |
|---|---|---|
| A — unspezifiziert lassen (kein explizites Konvergenz-Verhalten) | kein Aufwand | verletzt die Boundary-Akzeptanzkriterium-Pflicht, ein definiertes Verhalten zu benennen |
| B — getrennter Zustand: ein eigenes Metadatum (bzw. eine eigene Fehlermeldung) unterscheidet „Spalte ausgeschlossen" von „Spalte real gelöscht" im weiteren Verlauf | maximale Diagnostizierbarkeit für einen Administrator, der den Unterschied nachvollziehen will | verlangt eine dritte Zustandskategorie neben der bestehenden Abwesenheits-/Fehler-Semantik, ohne dass ein Akzeptanzkriterium sie verlangt — die Boundary-Klausel fordert ausdrücklich nur „gilt das Verhalten aus `LH-FA-SCH-003`", keine eigene Kategorie |
| **C — Konvergenz auf `LH-FA-SCH-003`s bestehendem Verhalten (gewählt)** | der Ausschluss wirkt (Teilfrage 3) ausschließlich auf der Row-Image-Konstruktionsebene, **unterhalb** der Schema-Vergleichsebene (`classifyRelationColumns`/`observeRelation`); eine real gelöschte, zuvor ausgeschlossene Spalte fehlt in der eingehenden `decode.Relation` genauso wie jede andere gelöschte Spalte und löst denselben `relationOther`-Pfad mit `ErrIncompatibleSchemaChange` aus wie jede reguläre Spaltenlöschung (`LH-FA-SCH-004.a`) — kein Sonderfall, keine Zusatzlogik nötig, die Boundary-Klausel ist damit durch reine Schichtung erfüllt | ein Administrator kann aus dem Fehler allein nicht ablesen, dass die gelöschte Spalte zuvor ausgeschlossen war — dieselbe Grenze, die auch für jede andere gelöschte Spalte gilt, kein neues Diagnose-Defizit |

### Teilfrage 5 — Fehlerpfad bei Ausschlussantrag gegen eine nicht existierende Spalte

| Option | Pro | Contra |
|---|---|---|
| A — stiller Erfolg (Antrag wird vermerkt, wirkt nie, weil die Spalte nie auftritt) | kein Prüfaufwand | verletzt `LH-FA-CFG-005`s Negative-Akzeptanzkriterium wörtlich („folgt ein expliziter Fehlerpfad") |
| B — ein neuer, vom Antrags-Status unabhängiger Fehlermechanismus (z. B. die SQL-Funktion selbst prüft `information_schema.columns` und wirft eine synchrone SQL-Exception) | Fehler wäre sofort beim `SELECT cdc.exclude_column(...)`-Aufruf sichtbar | bricht mit `ADR-0050`s etabliertem Muster: die schreibende SQL-Funktion darf ausschließlich einen Antrags-Datensatz schreiben, keine Prüf- oder Domänenlogik enthalten (`ADR-0018`/`ADR-0046`); die Spaltenexistenz-Prüfung würde in SQL dupliziert, obwohl der Go-Anwendungsdienst sie ohnehin treffen muss (Analogie zu `TableExists`) |
| **C — derselbe Fehlerpfad wie `cdc.enable_table` gegen eine nicht existierende Tabelle (gewählt)** | ein neuer Sentinel-Fehler `ErrSourceColumnMissing` (Muster von `ErrSourceTableMissing`, `internal/application/port/inbound/verwaltung.go`), geprüft über eine neue `ColumnExists`-Fähigkeit am Outbound Port und ausgelöst vom neuen `ExcludeColumnService`/`IncludeColumnService` bei der asynchronen Antrags-Verarbeitung (`applyAdministrationRequest`); der Fehlschlag landet als `failed` mit Fehlertext im Antrags-Datensatz — dieselbe Sichtbarkeit, die jeder andere Verarbeitungsfehler dieser Queue bereits trägt (`ADR-0050`) | asynchron: der Fehler ist erst über den Antrags-Status sichtbar, nicht als Rückgabewert des `SELECT`-Aufrufs selbst — dieselbe bereits in `ADR-0050` akzeptierte Asynchronitäts-Eigenschaft gilt hier fort |

### Bestätigung — Live-Reload-Konsistenz

Der bestehende Live-Reload-Mechanismus aus `ADR-0050` (Antrags-Queue,
verarbeitet von der bereits laufenden Administrations-Goroutine im
Capture-Prozess selbst) **reicht ohne Erweiterung des Grundmechanismus**
aus: `runAdministration`/`applyAdministrationRequest`
(`internal/bootstrap/wiring.go`) bekommen zwei weitere `case`-Zweige
(`exclude_column`/`include_column`) neben den bestehenden
(`enable`/`disable`), die nach erfolgreicher Persistierung des
Ausschlussstands eine neue, unter `tablesMu` synchronisierte
`Assembler`-Methode aufrufen (Update der `ExcludedColumns` einer
bestehenden `TableBinding`), statt eine neue Wecksignal-Verbindung oder
eine zweite Goroutine einzuführen. Ein bereits aktiver Ausschluss wird
damit ohne Prozess-Neustart wirksam, exakt wie eine Tabellen-Aktivierung
heute schon.

## Konsequenzen

- Positiv: `LH-FA-CFG-005` und `LH-QA-SEC-004` bekommen einen
  vollständigen, entscheidbaren Umsetzungspfad, ohne einen der drei
  bestehenden Grundsätze zu verletzen: keine Domänenlogik-Duplikation in
  SQL (`ADR-0018`/`ADR-0046`), ein einziger Schreibpfad auf
  administrative Zustände über Inbound Ports (`ADR-0028`/`ADR-0034`),
  Live-Reload ohne Neustart über die bereits etablierte Antrags-Queue
  (`ADR-0050`).
- Positiv: Die gewählte Wirkort-Entscheidung (Assembler, vor jeder
  Serialisierung) erfüllt `LH-QA-SEC-004` im engeren, sicherheitsrelevanten
  Sinn — ein ausgeschlossener Wert erreicht nie die Persistenzschicht,
  nicht nur eine Lese-View.
- Negativ: `TableBinding` und die `Assembler`-Live-Reload-Schreibpfade
  (`AddBinding` für neue Bindung/Schema-Bump) müssen so umgebaut werden,
  dass ein Nachtrag einer Schema-Version den Ausschlussstand einer
  bestehenden Bindung nicht stillschweigend zurücksetzt — ein
  Implementierungsdetail, das der umsetzende Slice sauber lösen muss
  (kein blindes vollständiges Überschreiben mehr, siehe Teilfrage 3).
- Negativ: `cdc.administration_request` (`tools/schema/schema.yaml`)
  braucht eine Schema-Erweiterung (mindestens eine zusätzliche,
  für zwei der vier Antragsarten genutzte Spalte für den Spaltennamen,
  und eine erweiterte `request_kind`-CHECK-Klausel) — konkrete
  Spaltenform, Name und Migrationsschritt sind Gegenstand des
  Pflichtenhefts und der umsetzenden Slices, nicht dieser Entscheidung
  (dieselbe Delegation wie in `ADR-0050`).
- Folgepflicht: `spec/architecture.md`s Sequenzdiagramm zu
  [`ARC-005`](../../../spec/architecture.md) (Administrator über SQL →
  Antragsqueue → Administrations-Hintergrundzug → `EnableTableUseCase`)
  zeigt heute nur die Antragsarten `enable`/`disable`; eine
  Planner-/Architect-Korrektur ergänzt die beiden neuen Antragsarten auf
  demselben Pfad.
- Folgepflicht: `spec/pflichtenheft.md` bekommt einen neuen `SPEC-*`-Eintrag
  für die konkrete Feldform des erweiterten Antrags-Datensatzes (Name der
  Spalte, Antragsarten-Erweiterung) sowie — falls Teilfrage 2 später per
  Folge-ADR auf eine globale Granularität erweitert wird — deren eigene
  `SPEC-*`-Festlegung.
- Folgepflicht: **Empfohlener Slice-Schnitt** (Empfehlung an den Planner,
  keine Festlegung dieser ADR):
  1. **SQL-Funktionen + Antrags-Verarbeitung** — Schema-Erweiterung von
     `cdc.administration_request` (Spalte, CHECK-Erweiterung),
     `cdc.exclude_column`/`cdc.include_column`, neuer Outbound Port
     (`ColumnExclusionPort` o. ä., `ColumnExists`-Fähigkeit analog
     `TableActivationPort.TableExists`), neue Inbound Ports
     (`ExcludeColumnUseCase`/`IncludeColumnUseCase`, `ADR-0028`/`ADR-0034`),
     `ErrSourceColumnMissing`, Erweiterung von
     `applyAdministrationRequest` um die zwei neuen Antragsarten.
  2. **Assembler-Filterung + Live-Reload-Verdrahtung** — `TableBinding`
     um `ExcludedColumns` erweitern, `rowImage`/`change` filtern lassen,
     Live-Reload-Nachtrag über eine neue synchronisierte
     `Assembler`-Methode, ohne den Ausschlussstand bei einem
     Schema-Versions-Bump zu verlieren (Konsequenz oben); Unit-Tests, die
     belegen, dass ein ausgeschlossener Spaltenschlüssel nie im
     resultierenden Row Image erscheint (weder `old_data` noch
     `new_data`), inklusive Konvergenz-Test für eine real gelöschte,
     zuvor ausgeschlossene Spalte (Teilfrage 4).
  3. **E2E-Beleg** — Erweiterung von `make test-integration`
     (`tools/harness/run-integration-tests.sh`), analog zum bestehenden
     Schema-Evolution-Rundlauf: `SELECT cdc.exclude_column(...)` gegen
     eine bereits aktivierte Tabelle im laufenden Feed-Container, Poll auf
     `status = 'applied'`, dann realer Beleg, dass künftige **und**
     — für den bereits vor Slice 2 gedeckten Fall — historische Changes
     gemäß `LH-FA-CFG-005`s Akzeptanzkriterien den ausgeschlossenen Wert
     nicht tragen; zusätzlich ein Negative-Beleg (`cdc.exclude_column`
     gegen eine nicht existierende Spalte, Antrag landet `failed` mit
     Fehlertext).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (Review-Prüfpflicht, wie `ADR-0018`/`ADR-0046`/`ADR-0050`) | `cdc.exclude_column`/`cdc.include_column` schreiben ausschließlich einen Antrags-Datensatz und senden `pg_notify` — kein direkter Schreibzugriff auf eine Ausschluss-Zustandstabelle oder `cdc.source_table` | kein Gate — `.a-check.yml` kennt nur Go-Globs, SQL-Layer-Edges sind ihm unsichtbar |
| — (Unit-Test, Umsetzungs-Detail des Slices) | Ein in `TableBinding.ExcludedColumns` geführter Spaltenname erscheint nie als Schlüssel im resultierenden `old_data`/`new_data`-JSON eines `Change` | `make test` (bestehendes Ziel) |
| — (Code-Review) | Die neue Live-Reload-Methode für Ausschluss-Änderungen greift wie `AddBinding`/`RemoveBinding` unter `tablesMu` zu — ohne sie ist der gleichzeitige Zugriff aus zwei Goroutinen ein Data Race (`go test -race`, `ADR-0030`) | `make test` (Race-Detector, bestehendes Ziel) |

## Re-Evaluierungs-Trigger

Zwei unabhängige Trigger, keiner davon macht die Entscheidung insgesamt
`permanent`:

- Derselbe Trigger wie in `ADR-0050`: Wird eine technische Brücke
  eingeführt, die SQL-Objekten einen echten, synchronen Port-Aufruf
  erlaubt (FDW, `dblink`, Extension mit Prozess-/Socket-Zugriff), wird
  re-evaluiert, ob die Antrags-Queue für die hier entschiedenen
  Antragsarten noch nötig ist.
- Granularität (Teilfrage 2): Sobald ein konkret benannter Bedarf für
  einen quellen-/musterweiten Spaltenausschluss über mehrere Tabellen
  hinweg entsteht — belegt entweder durch einen
  Beobachtungs-Register-Eintrag, der die Schwelle 3× erreicht, oder durch
  eine geschärfte `LH-QA-SEC-004` — wird Teilfrage 2 per Folge-ADR
  erneut entschieden; die übrigen vier Teilfragen bleiben davon
  unberührt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Architect-Entscheidung vor dem Schneiden der umsetzenden Slices; löst die drei von `LH-FA-CFG-005` selbst benannten offenen Fragen (Lastenheft-Version 0.7.0) | [`spec/lastenheft.md`](../../../spec/lastenheft.md) Commit `97de1e9` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0059` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
