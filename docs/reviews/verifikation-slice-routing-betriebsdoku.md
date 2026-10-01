# Verifikations-Report: slice-routing-betriebsdoku — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 8 und Entscheidung 4,
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
[`ADR-0141`](../plan/adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)) und Plan-vs-Doku-Diff.
Review-Artefakt: [`review-slice-routing-betriebsdoku.md`](review-slice-routing-betriebsdoku.md).
Formvorbild: [`verifikation-slice-routing-e2e.md`](verifikation-slice-routing-e2e.md).

**Gegenstand:** Slice-Plan [`slice-routing-betriebsdoku`](../plan/planning/done/slice-routing-betriebsdoku.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-ADM-001`](../../spec/lastenheft.md); Welle
[`welle-routing`](../plan/planning/welle-routing.md)). Diff `2e68966f~1..HEAD` (`e465ca1c`):
`docs/user/benutzerhandbuch.md` (Version 1.84) und der Plan; Commits `758ae7a3` (Plan §3 konkretisiert),
`be178eb6` (Handbuch 1.84), `8cd71c20` (Review), `e465ca1c` (Fixrunde). Reine Doku, `git diff --stat 2e68966f~1 HEAD`
über Handbuch und Plan: 2 Dateien, +907/−58; `internal`, `cmd`, `proto`, `gen` unberührt.

Dieser Lauf ändert weder Handbuch noch Plan noch Code (keine DoD-Häkchen; `grep -c '\[x\]'` im Plan gleich 0); er
schreibt nur diesen Report. Mutationen liefen an einer Kopie im Scratchpad (`sed … > Kopie`, nie `-i`, nie eine
Umleitung auf eine Repo-Datei). Die Wegwerf-Umgebung ist abgebaut (`down -v`, kein `cdc-*`-Container, kein Netz
`cdc-feed-test`, `git status --short` leer). Keine verweigerte Aktion ([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | gedruckte Zeile |
|---|---|---|
| `make docs-check` | 0 | `d-check: 1515 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-betriebsdoku.md` | 0 | `suchlauf-nachmessen: 16 Zeilen stimmen` |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK`, `gesamt: 0 Befund(e)` (a-check), `d-check: 1515 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)` |
| `make commit-traceability RANGE=2e68966f~1..HEAD` | 0 | `OK — 7 Commit(s) in "2e68966f~1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-immutable RANGE=2e68966f~1..HEAD` | 0 | `d-check: 1515 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=2e68966f~1..HEAD` | 0 | `d-check: 1515 Datei(en) geprüft, 0 Befund(e)` |
| `make test-notify` (Nachmessung der Kosten-Aussage) | 0 | `natsstream-kosten: je 10000 Changes, Median von 5 Läufen: ohne Ziel 16.518416ms (605385 Changes/s), mit Ziel 19.762699ms (506004 Changes/s)`, `mit/ohne = 1.20` |

Vorbedingung der Wahrheitsprobe: `git diff d8371372 HEAD --stat -- internal cmd proto gen` ist leer (Exit 0); das
`:dev`-Image (angelegt 2026-10-01 18:00) passt damit zum Stand der Produktivdateien, ich habe nicht neu gebaut. Dass
das Image aus genau diesem Stand gebaut wurde, ist nicht nachgemessen (Annahme aus Zeitstempel und leerem Diff).

**Hinweis zur Messhygiene:** eine erste Diff-Ausgabe (`git diff 2e68966f~1 HEAD`, früh im Lauf in eine Datei
geschrieben) zeigte noch die Fassung vor der Fixrunde; ich habe sie verworfen, den Diff neu erzeugt und alle Aussagen
unten gegen das aktuelle Handbuch (`Read` der Zeilen 340–770 und neu erzeugter Diff) geprüft.

## 2. Wahrheitsprobe (Wegwerf-Umgebung, selbst gefahren)

Aufbau: `compose.yaml` mit Überschreibung im Scratchpad (Quelle `meine-quelle`, Tabelle `public.orders (id, name,
region)`, eigene Publication und eigener Slot), PostgreSQL 18 (Digest des Standard-Pins), Feed als Superuser, Schema-Rollout
über `make schema-rollout`, zwei Login-Rollen: `adm` (Mitglied `cdc_admin`) und `rdr` (Mitglied `cdc_reader`). Alle
`set_route`/`remove_route`-Aufrufe unter `adm`, alle Lesungen unter `rdr` bzw. dem `reader`-Token.

**Beispiel des Handbuchs, Poll nach dem Wortlaut des Handbuchs unter `adm`:** beide Anträge `applied`. Nach den Zeilen
1–3 (vor den Regeln) und 4–6 (danach) gedruckt unter `rdr`:

```text
 id | region | route_target
----+--------+--------------
 1  | eu     |
 2  | us     |
 3  |        |
 4  | eu     | eu
 5  | us     | sonstige
 6  |        | sonstige
```

`WHERE route_target = 'eu'` liefert die Zeile 4; `IS NULL` die Zeilen 1–3. **Identisch** mit dem Handbuch-Beispiel.

**`GET /changes?target=` (Token `reader`):**

| Aufruf | gedruckt |
|---|---|
| `target=eu` | `HTTP/1.1 200 OK`, eine Change (`"id":"4"`, `"origin":"wal"`), kein Feld `route_target` in der Antwort |
| `target=Gross` | `200`, `{"changes":[]}` |
| `target=Gross&limit=0` | `400 Bad Request` |
| `target=eu&schema=other` | `200`, `{"changes":[]}` |
| unbekannter Parameter `foo=1` | `400 Bad Request` |

**Fehlertexte (`error_message` unter `adm`) gegen die R-Tabelle und die Formprüfungen:**

```text
R1 => failed | Regelname bereits vergeben: public.orders.eu_orders
R2 => failed | order bereits vergeben: public.orders.10
R3 => failed | Spalte existiert nicht an der Quelle: public.orders.regio
R3 ausgeschlossen => failed | Spalte ist ausgeschlossen: public.orders.name
exclude_column auf region => failed | Spalte trägt eine Routing-Bedingung: public.orders.region
R4 => failed | Regel ohne when bereits vorhanden: public.orders.rest
R4 => failed | order liegt hinter der Regel ohne when: public.orders.rest
R5 => failed | Bedingung bereits vergeben: public.orders.region
R6 => failed | Regelname nicht geführt: public.orders.nope
unbekannter Schlüssel => failed | unbekannter Schlüssel in rule_spec: foo
order 0.5 => failed | rule_spec ist ungültig: public.orders.f2
order 2147483648 => failed | rule_spec ist ungültig: public.orders.f4
Zielname X => failed | Zielname ist ungültig: public.orders.X
Regelname "Bad Name" => failed | Regelname ist ungültig: public.orders.Bad Name
order 1e1 (= 10) => failed | order bereits vergeben: public.orders.10
```

Alle zehn Zeilen der R-Tabelle bis auf die Zeile „R4 — die Abschlussregel trägt die höchste `order`“ (Wortlaut
`Regel ohne when trägt nicht die höchste order`) sind am System nachgefahren; diese eine Zeile ist **nicht gefahren**
(der Text stammt aus der Spec; die Zeile des Reviews trug sie als gefahren — von mir nicht nachgemessen). `1e1` wird als
positive ganze Zahl angenommen (die Formprüfung passiert, der Antrag endet erst an R2).

**`DELETE` ohne volle Replica-Identität und Backfill (Aussagen des Hinweises):** Zeile 5 gelöscht: `DELETE | {"id": "5"} |
sonstige`. Backfill-Run (`applied`) mit dem Regelstand: Zeilen 1 und 4 `eu`, übrige `sonstige`, `origin = backfill`.

**Ursache 2 (die in der Fixrunde berichtigte Zuordnung):** `ALTER PUBLICATION … SET TABLE public.orders (id, name)` an der
Tabelle mit bekannter Spaltenform, eine eingefügte Zeile: Container `exit=1`, gedruckt
`Fehlerklasse schema: Relation-Änderung nicht sicher als Obermenge interpretierbar: public.orders`,
`cdc.heartbeat.error_class` = `schema`, die Zeile erscheint nicht in `cdc.changes`. Danach **beide Regeln entfernt**
(Anträge `pending` bei stehendem Prozess), `docker start`: beide `applied`, der Container endet erneut mit `exit=1` und
demselben Text, `error_class` bleibt `schema`. Auch nach `ALTER PUBLICATION … SET TABLE public.orders` (volle Spaltenliste)
und erneutem Start: `exit=1`, derselbe Text. Die Aussage „das Entfernen der Regel genügt dort nicht“ ist damit **von mir
am System bestätigt**, nicht mehr nur vom Reviewer übernommen. (Wiederhergestellt habe ich die Umgebung nur durch
`pg_drop_replication_slot` und Neustart; das Handbuch nennt diesen Weg nicht, und er verwirft die unbestätigte Position.)

**Ursache 1 samt Abhilfe:** Tabelle `public.rt2` über `cdc.enable_table` aktiviert, Regel `rt_eu` auf `region`
(`applied`), `ALTER PUBLICATION … DROP TABLE` und `ADD TABLE public.rt2 (id,name)`, eine Zeile eingefügt: Exit 1, gedruckt

```text
Fehlerklasse schema: Routing-Regel auf die Änderung nicht anwendbar: Regel "rt_eu" an public.rt2: Spalte der Routing-Regel fehlt in den Spalten der Änderung: region
```

identisch mit dem Block im Handbuch. `cdc.remove_route` bei stehendem Prozess: `pending`; nach `docker start`: `applied`,
Container `running=true health=healthy`, die Zeile erscheint `{"id": "1", "name": "B"}` mit `route_target` leer (NULL),
kein zweiter `schema`-Fehler. (Diese Spalten-Probe lief an einer Tabelle, deren Aktivierung und Publication-Liste ich
selbst setzte; die Zahl der Spaltenform-Zeilen nach der Aktivierung habe ich nicht gezählt — die erste Zählabfrage
scheiterte an einem falschen Tabellennamen, ich habe sie nicht wiederholt.)

**Kosten-Aussage:** mein Lauf `make test-notify` (siehe §1): Verhältnis 1,20, innerhalb der im Handbuch genannten
Spanne 0,92 bis 1,30; die Aussage „im Rauschen, ein Aufschlag nicht auflösbar“ widerspricht dem nicht. Das Handbuch
nennt sechs Läufe; mein siebter ist hier gemessen und liegt im Handbuch nicht.

## 3. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | (A) Abschnitt „Routing-Regel konfigurieren“ vollständig und wahr; Beispiele unter der Rolle gelaufen; R1–R6, Konsequenzen, Abhilfe, DELETE-Hinweis mit Ursprung und Menge | **getragen** | Abschnitt steht hinter „Transformationsregel konfigurieren“ (Voraussetzung, Vorgehen mit `rule_spec`-Tabelle, Beispiel, Entfernen, Fehlertabelle, Fehlerklasse `schema` mit zwei Ursachen, Ziel lesen, Hinweise). §2: Beispiel, R-Tabelle, Formprüfungen, `order`-Form (`10`, `10.0`/`1e1` Form-Annahme, Obergrenze), DELETE und Backfill am System bestätigt. DELETE-Hinweis nennt PostgreSQL 17.11 und 18.6 und die Fälle (Abschlussregel, ohne, `FULL`, `INSERT`/`UPDATE`). Eine R-Zeile (R4, „höchste order“) ist nicht von mir gefahren (§2) |
| 2 | (B) Lesewege: `target` je Weg, „bei sonst gültiger Anfrage“, Zählwort drei, Altserver-Aussage mit Ursprung | **getragen** | HTTP-, gRPC-Stream-, `ReadChanges`- und SSE-Abschnitt tragen `target` (Diff gelesen); „bei sonst gültiger Anfrage“ steht an HTTP und `ReadChanges` (dort: Lese-Kontrakt hat Vorrang); „drei“ an den Roh-Wegen (`make suchlauf-nachmessen`: `diff 2 … 'drei optionale'`); Altserver-Aussagen sind *hergeleitet* gekennzeichnet, nicht als gefahren. HTTP-Verhalten (200 leer, 400 mit `limit=0`, Konjunktion) von mir gefahren (§2). gRPC-/SSE-/NATS-Ziel-Auswahl am System: **übernommen** aus dem E2E-Verifikations-Report (verlinkt) |
| 3 | (B) NATS-Subjekt `cdc.route.<source_id>.<ziel>` und Kostenabsatz mit Ursprung | **getragen mit LOW (V-1)** | Subjekt, Fire-and-Forget, Lücke schließen über `route_target`/`GET /changes?target=`, `eu-west_1`, `.*`/`.>`, Wecksignal-Wurzel nur negativ belegt: alles mit Ursprung. Kostenabsatz trägt Methode, sechs Läufe, Spanne und Grenze; mein siebter Lauf 1,20 passt. Form des Ursprungs: V-1 |
| 4 | (B) Spalte `route_target` in „Änderungen lesen“ (NULL = nicht geroutet, `origin` davor) | **getragen** | SELECT-Liste und Satz: „`route_target` ist die letzte Spalte der View, `origin` steht davor“, „`NULL` heißt ‚nicht geroutet‘“; HTTP-Antwort ohne `route_target` (von mir gelesen: das Feld fehlt in der Antwort, §2); „Schema aktualisieren“ nennt beide Spalten; Verwaltungs-API-Feldliste: Ziel gehört nicht dazu |
| 5 | (C) Zeile `schema` der Fehlerklassen-Tabelle, „Neustart nach einem Fehler“, Rollen-Tabelle, Grenzwerte, Glossar | **getragen** | Zeile `schema` nennt Routing-Regel und den Run-Fall (Klasse, vor der ersten Zeile); „Neustart nach einem Fehler“ trennt die nicht anwendbare Regel (Abhilfe `remove_route` + Neustart) von der entfernten Spalte; Rollen-Tabelle: `cdc.set_route`/`cdc.remove_route`; Grenzwerte: Zielname-Alphabet `[a-z0-9][a-z0-9_-]{0,62}`, `order` bis 2147483647 als **Setzung**, keine Regelzahl-Grenze (die Spec nennt keine); Glossar: zwei Einträge |
| 6 | Version 1.84 und Historienzeile hinter 1.83 | **getragen** | Kopfzeile `Version: 1.84`; Zeile `| 1.84 | 2026-10-01 |` ist die letzte Zeile der Datei, direkt hinter `| 1.83 |`; Suchlauf: 1 Treffer `^Version: 1\.84$`, 5 Zeilen der Reihe 1.8x |
| 7 | Übergabe-Block aus allen Slices der Welle abgedeckt | **getragen, Gegenstand für Gegenstand** | Tabelle in §4 |
| 8 | `make gates` grün, Exit ungefiltert | **getragen** | §1, eigener Lauf, Exit 0 |
| 9 | Review durchgeführt, kein offenes HIGH/MEDIUM | **getragen in der Sache; Häkchen gehört dem Planner** | [Review-Report](review-slice-routing-betriebsdoku.md): 0 HIGH, 0 MEDIUM, 3 LOW, 2 INFO; F-1 bis F-5 siehe §5 |
| 10 | §3.13-Suchlauf, Feld mit Gefundenem und Nichtgefundenem, Nachmessen Exit 0 | **getragen** | Exit 0, 16 Zeilen; das Feld nennt je Träger Gefundenes und Nichtgefundenes (Transformations-Funktionen, `cdc.changes`, Version/Historie, Filter-Beschreibung, `README.md`/`examples`/`sdks` ohne Treffer) |
| 11 | Doku-Update, `README.md` nur bei Fund | **getragen** | Handbuch ist das Update; das Feld stellt für `README.md`, `examples/README.md`, SDK-READMEs 0 Treffer fest |
| 12–16 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter, §6 „bei der Closure einzutragen“ |

## 4. Übergabe-Block: Gegenstand für Gegenstand (`git grep` am Stand `HEAD`)

| Quell-Slice und Gegenstand | Fundstelle im Handbuch | Befund |
|---|---|---|
| [`slice-routing-kern-label`](../plan/planning/done/slice-routing-kern-label.md): Spalte `route_target`, `NULL` ≠ `wal` | „Änderungen lesen“ (Satz zur letzten Spalte, „anders als bei `origin` liest ein fehlender Wert nicht als Vorgabe“) | gedeckt |
| [`slice-routing-antragsweg`](../plan/planning/done/slice-routing-antragsweg.md): Funktionen, Rolle, `rule_spec`, R1–R6, R4-Abhilfe, `order`-Grenze und Zahlform | Rollen-Tabelle, Abschnitt Routing (`rule_spec`-Tabelle, R-Tabelle, „Lösung“) | gedeckt; Zahlform `10`/`10.0`/`1e1` und Grenze im Text, die Ablehnungsfälle (Bruchteil, 0, negativ, über Grenze) ebenso; `0.5` und `2147483648` von mir abgelehnt gesehen |
| [`slice-routing-backfill-pfad`](../plan/planning/done/slice-routing-backfill-pfad.md): Label des Regelstands zum Run, Stand-Wechsel (`configuration`), Snapshot ohne Spalte (`schema`, vor der ersten Zeile) | „Bestand als Backfill überführen“ (Aufzählung „Routing-Regeln“), Zeile `schema`, Abschnitt Routing („Im Backfill“), „Neustart nach einem Fehler“ | gedeckt; Lesefehler = Klasse der Ursache ([`ADR-0141`](../plan/adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md)) steht im Text; „durch Unit-Tests belegt, am System nicht gefahren“ ohne Test- oder Dateinamen |
| [`slice-routing-lesewege`](../plan/planning/done/slice-routing-lesewege.md): `target` je Weg, Konjunktion, Qualifier, Altserver, Auswahl ≠ Zugriffsschutz, Zählwort drei | HTTP, gRPC-Stream, `ReadChanges`, SSE; Hinweis „Auswahl, kein Zugriffsschutz“ | gedeckt; die drei SDK-Stellen (C#, Kotlin, Python) tragen weiter „zwei optionale Parameter `schema`/`table`“ und „das Feld `target` folgt mit dem Package“ (Suchlauf `diff 2 'zwei optionale'` plus die C#-Stelle, die der Zeilenmuster nie trifft) |
| [`slice-routing-nats-subjekt`](../plan/planning/done/slice-routing-nats-subjekt.md): Zusatz-Subjekt, Wecksignal-Wurzel, Unabhängigkeit der beiden Veröffentlichungen, Kosten | „Zugriff über das NATS-Wecksignal“ (Wurzel, nur negativ belegt), „Zugriff über den NATS-Vollinhalts-Stream“ (Subjekt, Kosten, Zustellsemantik mit Lücke-schließen) | gedeckt |
| [`slice-routing-e2e`](../plan/planning/done/slice-routing-e2e.md): DELETE-Messung 17.11/18.6, Erreichbarkeit (a)/(b)/(c), Abhilfe, Replay/Backfill | Hinweis `DELETE`, Ursache 1/2, Hinweis „Ziel zum Erfassungszeitpunkt fest“ | gedeckt; (a) als übernommen, die Hergeleitet-Marke für „Regel genügt nicht“ ist durch meine Probe belegt (§2) |
| [`slice-routing-spec-nachzug`](../plan/planning/done/slice-routing-spec-nachzug.md) / Entscheidung 4 | DELETE-Hinweis | gedeckt |

Nicht-Gefundenes (von mir gegrept): `git grep -n -E 'TestE2E|Test[A-Z][A-Za-z]+\(|_test\.go' -- docs/user/benutzerhandbuch.md` findet
nur die Historienzeile 1.68 (Bestand); `git grep -n -E 'cdc\.set_route|route_target' -- README.md examples sdks` 0 Treffer
(Bestätigung des Suchlauf-Felds).

## 5. Review-Findings F-1 bis F-5 — am Handbuch gelesen, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (LOW) Ursprung „gemessen im Lauf …“ ohne Anker | **geschlossen für die E2E-Werte, Rest V-1** | Alle E2E-Aussagen tragen die Form „*Ursprung:* übernommen aus dem E2E-Lauf von `make test-integration` (Verifikations-Report [`verifikation-slice-routing-e2e`](verifikation-slice-routing-e2e.md)) …“ (Neustart/Aktivierung, Ursache 1, Ursache 2 für die entfernte Spalte, Replay, DELETE, Happy-Path-gRPC/SSE/`ReadChanges`); der Link auf `e2e-abdeckung.md` ist weg; „gemessen an der Protobuf-Bibliothek“ ist ersetzt durch „übernommen aus den Tests der Lesewege“ mit Link auf [`verifikation-slice-routing-lesewege`](verifikation-slice-routing-lesewege.md); das NATS-Alphabet-Argument verweist auf [`verifikation-slice-routing-nats-subjekt`](verifikation-slice-routing-nats-subjekt.md); die `order`-Grenze trägt „gemessen in einer Compose-Umgebung (PostgreSQL 18.6)“ (von mir bestätigt, §2). Im Handbuch steht **nichts** als dort selbst gemessen, was nur der E2E-Lauf maß. Offen: der Kostenabsatz (V-1) und die Versionsangabe 18.6 (V-3) |
| F-2 (LOW) Ursache-Zuordnung | **geschlossen** | Ursache 1 enthält die Publication-Spaltenliste nur noch für die Tabelle ohne Spaltenform und verweist den Fall mit bekannter Form auf Ursache 2; Überschrift von Ursache 2 lautet „entfernt **oder fehlt in der Relation**“, nennt beide Anlässe und die gemessene Ausgabe; „Regel genügt nicht“ trägt jetzt eine Messung. Alles von mir am System bestätigt (§2). Ein Restpunkt: V-2 |
| F-3 (LOW) Transformations-Fehlerblock | **geschlossen** | Der Satz „Spaltenform noch nicht bekannt“ ist gestrichen; die Ursache nennt „`column` fehlt in der Relation der Change“ und sagt ausdrücklich, dass am System nur die Namenskollision gefahren ist; die entfernte Spalte ist als *hergeleitet* getrennt, mit Verweis auf den Routing-Abschnitt |
| F-4 (INFO) Adresse des unbekannten Schlüssels | **geschlossen** | Einleitungssatz: „bei einem unbekannten Schlüssel der bloße Schlüsselname ohne Tabelle“; Formprüfungen nennen `unbekannter Schlüssel in rule_spec: foo` (von mir gedruckt gesehen, auch ohne Tabelle) |
| F-5 (INFO) Test- und Dateinamen in Betreiber-Prosa | **geschlossen** | Kein Testname, keine Testdatei, kein `routing_test.go` mehr im Text (Grep §4); „durch Unit-Tests belegt, am System nicht gefahren“. Verbleibend: Phasen-Namen des Runners („Routing-Happy-Path (fünf Zustellwege)“) und `make test-notify`/`make test-integration` als Ursprung — Skript- bzw. Zielnamen, wie im Bestand üblich |

Keine Chronik-Sprache in neuer Prosa im engeren Sinn (kein „früher/vorher“); zwei Randformen: „Handbuch-Zug“ und
„Review, Fixrunde, Verifikation“ im Kostenabsatz (V-1).

## 6. Mutations-Analog (drei Tatsachenaussagen an einer Kopie verändert, Kopie im Scratchpad, `sed … > Kopie`)

| # | Mutation der Kopie | Welche Probe entlarvt sie | Gedruckte Gegenaussage |
|---|---|---|---|
| M1 | Fehlertext R1 `Regelname bereits vergeben` → `Regelname schon vergeben` | `set_route` mit vorhandenem Regelnamen, Poll auf `error_message` | `failed | Regelname bereits vergeben: public.orders.eu_orders` |
| M2 | Statuscode `?target=Gross&limit=0` `400` → `200` | `GET /changes?source=…&target=Gross&limit=0` | `HTTP/1.1 400 Bad Request` |
| M3 | `order` „bis 2147483647“ → „bis 2147483648“ | `set_route` mit `order` 2147483648 | `failed | rule_spec ist ungültig: public.orders.f4` |

`make docs-check` sähe keine der drei Mutationen (Sensor prüft Referenzen, nicht Tatsachen; an der Kopie nicht gefahren,
weil `d-check` das Repo liest — *hergeleitet*). Die Proben aus §2 sind die einzigen Falsifikatoren dieser Klasse; das
bestätigt, dass die Handbuch-Proben wiederholt und nicht übernommen werden müssen.

## 7. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 8 und Entscheidung 4 (Handbuch-Hinweis zu DELETE): getragen (DoD 1, §4).
- [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md): Trennung der zwei Ursachen im Handbuch umgesetzt; der Fall „Spaltenliste an Tabelle mit bekannter Form“ gehört zur Ursache 2 (Option 3 der ADR deckt ihn, keine neue ADR nötig); am System bestätigt, einschließlich „beide Regeln entfernen + Neustart bleibt rot“.
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)/[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)/[`ADR-0141`](../plan/adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md): `target` je Weg, Qualifier und Run-Fälle gedeckt (Backfill-Fälle als „durch Unit-Tests belegt“, nicht als gefahren).
- [`AGENTS.md`](../../AGENTS.md) §3.12: Ursprung je Aussage vorhanden bis auf die Form im Kostenabsatz (V-1). §3.13: Suchlauf Exit 0. §3.7: keine neue Kommentar-Form (reine Doku). §3.9: alle Exit-Codes einzeln gesichert. §3.5: keine `Accepted` ADR berührt (`make doc-immutable` Exit 0). §3.1: kein Host-Interpreter, kein `sed -i`, keine Umleitung auf eine Repo-Datei.
- Wahrheitsklasse „Handbuch führt nur gemessene Wirkung als Tatsache“: erfüllt; die Hergeleitet-Marken (Altserver, Schlüsselspalten-Bedingung bei DELETE, unverändertes TOAST, Abhilfe (a)) sind als solche gekennzeichnet.

## 8. Re-Review der Fixrunde (`e465ca1c`)?

**Re-Review: nein.** Begründung gegen das Register
[`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md),
engere Fassung „Produktionslogik oder Norm geändert, oder nach der Fixrunde hat kein anderer Kontext sie ausgeführt“:

1. **Keine Produktionslogik, keine Norm:** der Diff der Fixrunde berührt nur Handbuchtext; `internal`, `cmd`, `proto`, `gen`
   und die Spec bleiben unberührt (§1: leerer Diff gegen `d8371372`).
2. **Anderer Kontext hat sie ausgeführt:** die einzige neue Tatsachenaussage der Fixrunde (Spaltenliste an Tabelle mit
   bekannter Spaltenform: Exit 1, Klasse `schema`, „Regel entfernen genügt nicht“) habe ich in frischem Kontext am System
   gefahren und bestätigt, samt der Gegenprobe „beide Regeln entfernen + Neustart“ und „Publication zurücksetzen“ (§2).
   Der Implementer hatte sie nur vom Reviewer übernommen; diese Lücke ist hiermit geschlossen.
3. Die Ursprungs-Umformulierungen (F-1) sind Formänderungen an Aussagen, die ich inhaltlich gegen die verlinkten
   Berichte und gegen das System geprüft habe. Die zwei offenen Formpunkte (V-1, V-3) sind LOW/INFO und brauchen keinen
   zweiten Reviewer-Lauf, sondern einen Planner-Entscheid (Nachzug oder Hinnahme).

## 9. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | LOW | Der Kostenabsatz trägt seinen Ursprung in einer Form, die ein Betreiber nicht auflösen kann: „gemessen und abgeleitet in diesem Handbuch-Zug“, „übernommen aus den Berichten der Umsetzung und ihrer Prüfung“ und „(Review, Fixrunde, Verifikation)“ ohne Link auf einen Bericht; die gedruckte Zeile des Handbuch-Zugs (Verhältnis 1,02) steht in keinem committeten Report. Die Werte 1,05 und 0,92 stehen im [Review-Report](review-slice-routing-nats-subjekt.md) bzw. im [Verifikations-Report](verifikation-slice-routing-nats-subjekt.md) des Quell-Slice. Die Zahlen sind wahr (mein Lauf 1,20 liegt in der Spanne); beanstandet ist die Form (Prozess-Vokabular, fehlende Anker). Vorschlag: auf die Bandbreite und die beiden Berichte verlinken, „Handbuch-Zug“ streichen | `docs/user/benutzerhandbuch.md` Absatz „Kosten der zweiten Veröffentlichung“ |
| V-2 | LOW | Ursache 2 nennt die Abhilfe nur als „die der inkompatiblen Schemaänderung (`LH-FA-SCH-003`, `LH-FA-SCH-004`)“; das Handbuch beschreibt diese Abhilfe an keiner Stelle (das Stichwort „inkompatible Schemaänderung“ steht nur an den zwei Stellen dieses Zugs). Ein Betreiber im Fall „Publication mit Spaltenliste“ erfährt, was **nicht** hilft, aber nicht, was hilft. Am System erholte sich die Umgebung nur durch Verwerfen des Slots (verliert die unbestätigte Position); das ist nicht als Anleitung gemeint. Kein Fehler der Aussage, eine Lücke im Handbuch — Gegenstand für den Planner (Folge-Slice oder benannte Hinnahme) | `docs/user/benutzerhandbuch.md` „Ursache 2“, „Neustart nach einem Fehler“ |
| V-3 | INFO | Ursache 1 nennt „PostgreSQL 17.11 und 18.6“ und verweist auf den Verifikations-Report des E2E-Slice; dort ist die Phase „Routing-Nichtanwendbarkeit und Abhilfe“ nur an PostgreSQL 17 vom Verifier gedruckt, der 18-Lauf nach der Fixrunde ist dort ausdrücklich als vom Implementer übernommen geführt. Die Zusammenfassung „beide Versionen“ ist im Plan als gemessen geführt, der verlinkte Bericht trägt sie für die Phase nur halb | `docs/user/benutzerhandbuch.md` Ursache 1; [`verifikation-slice-routing-e2e`](verifikation-slice-routing-e2e.md) §1 |
| V-4 | INFO | Die R-Zeile „die Abschlussregel trägt die höchste `order`“ (Text `Regel ohne when trägt nicht die höchste order`) ist von mir nicht nachgefahren (die R4-Folgezeile „liegt hinter“ und „bereits vorhanden“ ja); die übrigen neun Zeilen und alle Formprüfungen sind gedruckt bestätigt | R-Tabelle |
| V-5 | INFO | Plan-Kopf „Berührte Spec-Stellen: —“ und Plan-Tabelle sind mit der Konkretisierung konsistent; Plan §7 und §6 sind noch Platzhalter (Closure steht aus) | Plan |

Kein HIGH, kein MEDIUM. Kein Befund am Produktionsverhalten (Produktivcode nicht im Diff).

## 10. Verdikt

**DoD getragen: ja für die Zeilen 1 bis 8, 10 und 11; Zeile 9 in der Sache getragen (Häkchen: Planner); Zeilen 12 bis 16 offen, gehören dem Planner.** Das
Handbuch 1.84 ist wahr gegen das System: Beispiel, `GET /changes?target=`, die R-Tabelle und die berichtigte Ursache 2
(einschließlich „beide Regeln entfernen + Neustart bleibt rot“) sind mit gedruckter Ausgabe von mir bestätigt. F-1 bis
F-5 sind geschlossen, mit zwei verbleibenden Formpunkten (V-1, V-3).

**Nötiger Nachzug (Planner):**

1. V-1 im Handbuch nachziehen (Kostenabsatz: auflösbare Anker, kein „Handbuch-Zug“) oder als LOW hinnehmen und in §7 des Plans nennen.
2. V-2 als Folge-Slice-Kandidat oder benannte Hinnahme führen (operative Abhilfe der inkompatiblen Schemaänderung im Handbuch).
3. §6-Ausgänge eintragen: „Aussage ohne Messung“ → ausgeführt, V-3 bleibt als Hinweis; „Beispiele unter der falschen Rolle“ → Beispiele unter `adm`/`rdr` gelaufen (gedruckt in §2 und im Review); „Zwischenzeit mit den SDK-Packages“ → weiter offen bis [`slice-routing-sdk-beispiel-target`](../plan/planning/done/slice-routing-sdk-beispiel-target.md).
4. Closure-Notiz mit Lerneintrag, Beobachtungs-Register (Eintrag `fixrunde-ohne-reviewer-lesung`: die enge Fassung trägt diesen Fall), drei Paarungen der Welle.

**Re-Review:** nein (§8).
