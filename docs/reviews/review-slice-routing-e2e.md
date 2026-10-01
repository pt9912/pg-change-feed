# Review-Report: slice-routing-e2e — 2026-10-01

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10).
Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-routing-e2e`, Diff-Range `139b25c9~1..HEAD` (HEAD `577e76e7`).
Gegenstand sind die Commits `7ee258b1` (Testfunktionen, Runner-Phasen, Wegwerf-Clients,
Abdeckungs-Erzeugnis) und `577e76e7` (Plan-Nachzug, Suchlauf-Befunde, README-Zeilen,
Übergabe-Block im Plan `slice-routing-betriebsdoku`). Die dazwischen liegenden Commits des
Hauptlaufs (Lastenheft 0.14.0, Linkziele in Plänen und Reports) sind nicht Gegenstand; die
Lastenheft-Änderung liegt in einem eigenen Commit vor dem umsetzenden Test-Commit.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle
Mutationen liefen an `git archive HEAD`-Kopien im Scratchpad (Mutation per `sed … > Datei.new && mv`
innerhalb der Kopie, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei); `make test`,
`make fmt-check`, `make doc-trace`, `make docs-check`, `make suchlauf-nachmessen`,
`make kommentar-kennungen` und `make test-integration` liefen gegen den echten Arbeitsbaum,
ungefiltert, Exit-Code direkt gelesen ([`AGENTS.md`](../../AGENTS.md) §3.9). Keine Aktion wurde
von der Berechtigungsschicht verweigert.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-routing-e2e` (§1 Ziel und Abgrenzung, §2 DoD, §3 Plan samt Suchlauf-Feld,
  §6 Risiken) und Welle `welle-routing`
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) (Entscheidung 4,
  Folgepflicht 7),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
  [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  (Entscheidung 4, V3), [`ADR-0012`](../plan/adr/0012-at-least-once.md) und
  [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md) (nur Frage f)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md) (Lastenheft 0.14.0),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md), [`LH-QA-POR-001`](../../spec/lastenheft.md),
  [`LH-FA-REA-005`](../../spec/lastenheft.md), [`LH-FA-ADM-003`](../../spec/lastenheft.md)
- [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md), [`SPEC-032`](../../spec/pflichtenheft.md),
  [`SPEC-019`](../../spec/pflichtenheft.md), [`SPEC-024`](../../spec/pflichtenheft.md),
  [`SPEC-031`](../../spec/pflichtenheft.md), [`SPEC-008`](../../spec/pflichtenheft.md)
- [`AGENTS.md`](../../AGENTS.md) Hard Rules (§3.1, §3.7, §3.9, §3.12, §3.13)
- vorherige Reviews am Modul: `review-slice-capture-retry-realtest-belege-schaerfen`,
  `review-slice-wal-fehlerschwelle-ausgangsklasse`

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

| Lauf | Ergebnis |
|---|---|
| `make test` | Exit 0 (gefahren) |
| `make fmt-check` | Exit 0 (gefahren) |
| `make docs-check` | Exit 0 (gefahren) |
| `make doc-trace` | Exit 0, gedruckte Zeile `80 Anforderung(en), 0 Waise(n).` (gefahren) |
| `make suchlauf-nachmessen PLAN=…/slice-routing-e2e.md` | Exit 0, „10 Zeilen stimmen“ (gefahren) |
| `make kommentar-kennungen DIFF=139b25c9~1 TESTS=exclude` | Exit 0, kein Kandidat (gefahren) |
| `make kommentar-kennungen DIFF=139b25c9~1` | Exit 1 (über `make` Exit 2), vier Kandidaten, genau die Godocs der vier `TestE2ERouting*`-Funktionen (`routing_e2e_test.go:255–263`, `351–359`, `482–488`, `567–574`) (gefahren) |
| `make image`, dann `make test-integration` (PostgreSQL-Default des Legs 18) | Exit 0; Dauer von `make test-integration` 10 min 5 s (Differenz der Zeitstempel von `image.exit` und `ti.exit`, gemessen); gedruckte Zeilen unten |

Gedruckte Zeilen des eigenen `make test-integration`-Laufs (PostgreSQL 18.6):

```text
ROUTING-HAPPY-PATH Inhaltsregeln ["eu" "us" "" "" "eu" "eu"], Herkunftsregel ["herkunft" "herkunft" "herkunft"], Auflösung ["frueh" "spaet" "sonstige"]
ROUTING-REPLAY WAL-Ziele ["" "eu" "" "europa"], Backfill-Ziele zeilenweise 1="europa" 2="europa" 3="" 4="europa"
ROUTING-DELETE-MESSUNG PostgreSQL 18.6, Replica-Identität DEFAULT, Regeln eu+rest: INSERT="eu" UPDATE="eu" DELETE="sonstige" (gesetzt true), Alt-Bild des DELETE {"id": "1"}
ROUTING-DELETE-MESSUNG PostgreSQL 18.6, Replica-Identität DEFAULT, Regeln eu: INSERT="eu" UPDATE="eu" DELETE="" (gesetzt false), Alt-Bild des DELETE {"id": "1"}
ROUTING-DELETE-MESSUNG PostgreSQL 18.6, Replica-Identität FULL, Regeln eu+rest: INSERT="eu" UPDATE="eu" DELETE="eu" (gesetzt true), Alt-Bild des DELETE {"id": "1", "name": "Ada2", "region": "eu"}
run-integration-tests: Routing-Happy-Path (LH-FA-CFG-008) belegt — … 1 Dreiergruppe(n) (asia, us, eu) …
run-integration-tests: Routing-Neustart und Ausschluss-Sperre (LH-FA-CFG-008) belegt — …
run-integration-tests: Routing-Aktivierung über die API (LH-FA-CFG-008) belegt — … (Ziel je Tabelle: feed_e2e_route_api_grpc_rule: api; feed_e2e_route_api_grpc_plain: NULL; feed_e2e_route_api_http_rule: api; feed_e2e_route_api_http_plain: NULL;)
run-integration-tests: Routing-Nichtanwendbarkeit und Abhilfe (LH-FA-CFG-008) belegt — (b) … (0 Spaltenform-Zeilen nach der Aktivierung) und (c) … (0 Spaltenform-Zeilen) beendeten den Erfassungspfad real mit Klasse schema … (a) … (3 Spaltenform-Zeilen) …
run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen
```

Der Lauf an PostgreSQL 17 (17.11, laut Übergabe-Block im Plan `slice-routing-betriebsdoku`) ist
**übernommen**, nicht von mir gefahren; die PostgreSQL-18-Zeilen des Blocks (`DELETE="sonstige"`,
`DELETE=""`, `DELETE="eu"`, Alt-Bild `{"id": "1"}`) stimmen mit meinem Lauf überein.

Mutationen (alle drei an Kopien, am Runner gefahren, **gefahren**: 3):

| Nr. | Mutation | Stelle | Farbe |
|---|---|---|---|
| G1 | Gegenprobe ohne volle Replica-Identität: `c.identity == "FULL"` zu `"FULLX"` | `test/integration/routing_e2e_test.go:594` | rot: „DELETE-Messung an PostgreSQL 18.6 weicht von der Erwartung ab: [feed_e2e_route_del_full: DELETE="sonstige" … erwartet "eu" …]“ |
| G4 | `order` der Regel `frueh` von 5 auf 50 (Auflösung kippt) | `routing_e2e_test.go:335` | rot: „Auflösung bei zwei treffenden Regeln: ["spaet" "spaet" "sonstige"], erwartet ["frueh" "spaet" "sonstige"]“ |
| M1 | `grpcclient` verwirft das Ziel (`Target: (*target)[:0]`) | `tools/harness/grpcclient/main.go:69` | rot im Runner: „Routing-Happy-Path — grpc, Ziel eu — persistiertes Ziel der Change 2470-1 — erwartet 'eu', gelesen 'NULL'“ |

Ein erster Versuch von M1 (`Target: ""`) ließ eine Variable ungenutzt und scheiterte am Compiler;
er zählt nicht. Mutationen des Implementers: nicht vorgelegt, keine Zahl.

---

## Findings

### F-1 — „Ein auf ein Ziel gewählter Stream-Leser sieht nur dieses Ziel“ ist an den Stream-Wegen nur über die erste empfangene Change belegt

- `kategorie`: MEDIUM
- `quelle`: Skill-Klasse „Beleg trägt seinen Satz nicht“ (im Skill unter HIGH geführt; hier
  als MEDIUM eingestuft, Begründung im Befund); Plan §2 DoD (1) „Ziel A sieht nur A“;
  [`LH-FA-CFG-008`](../../spec/lastenheft.md) Happy Path;
  [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md) „Auswahl, keine Ausblendung“,
  [`SPEC-024`](../../spec/pflichtenheft.md)
- `pfad`: `tools/harness/run-integration-tests.sh:4680–4688` (Start der sechs Clients mit Ziel),
  `:4737–4751` (Prüfung `head -n1`); `tools/harness/grpcclient/main.go:41,79`,
  `tools/harness/sseclient/main.go:39,98`, `tools/harness/natsstreamsub/main.go:56,104`
  (`-count` Default 1)
- `befund`: Die sechs Clients mit Ziel laufen mit `-count 1` und enden nach der ersten
  empfangenen Change; der Runner prüft nur diese eine Zeile (Region und persistiertes Ziel
  gleich dem Ziel). Was der Client danach geliefert bekäme, wird nicht beobachtet, und kommt der
  Client nach der fremden Change einer Dreiergruppe in den Empfang (Registrierungsfenster), kann
  auch ein Client mit nicht greifendem Filter die erste Zeile treffen. Der Satz „sieht nur dieses
  Ziel“ steht in der `abdeckung_declare`-Zeile (damit in `docs/user/e2e-abdeckung.md`) und im DoD;
  getragen ist „der Filter weist die vorangestellte fremde Change ab“ (M1: rot, weil die erste
  Change der Gruppe `asia` ohne Ziel ist), nicht „keine fremde Change danach“. An den Pull-Wegen
  (`GET /changes?target=`, `ReadChanges`) ist „nur“ dagegen vollständig getragen (Gleichheit der
  Kennungsmengen mit `WHERE route_target`). Einstufung MEDIUM statt der im Skill geführten HIGH:
  der Filter ist je Change zustandslos, ein Leck erst hinter der ersten Treffer-Change ist keine
  plausible Fehlerform, und der Implementer nennt die Grenze selbst; die Aussage ist trotzdem breiter
  als ihre Messung.
- `verifizierbar`: ja — Mutation eines Filters, der erst ab der zweiten Change leckt, oder ein
  verspäteter Client; M1 belegt nur den Fall des vollständig verworfenen Ziels.
- `klasse`: Beleg trägt seinen Satz nicht (Aussage „nur“ breiter als die Messung „erste Change“)

### F-2 — `TestRunStreamWithRetrySlotStillActive` bindet „genau eine Lieferung im zweiten Versuch“ strenger, als die At-least-once-Zusage hält (Flake, Folge-Slice)

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0012`](../plan/adr/0012-at-least-once.md) (At-least-once),
  [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md) Fitness Function
  „Slot noch aktiv“, [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)
  (Leerlauf-Bestätigung); Maintainability
- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go` (Prüfung
  `len(delivered[2]) != 1`, „Lieferungen des zweiten Versuchs … erwartet genau eine“);
  `internal/adapters/driving/replication/mapper/mapper.go:206` (Position = `CommitLSN`);
  `internal/adapters/driving/replication/decode/decode.go:42–51`
- `befund`: Der Test wartet, bis `confirmed_flush_lsn` die Position der Halter-Transaktion
  erreicht (`flush >= haltedPosition`), und startet den Retry sofort; die bestätigte Position einer
  Transaktion ist ihre `CommitLSN` (Beginn des Commit-Records), nicht ihr `EndLSN`. Liegt der Slot-Stand
  beim Aufbau des zweiten Versuchs genau auf dieser Position, ist eine erneute Zustellung der bereits
  persistierten Halter-Transaktion eine legitime Wiederzustellung unter At-least-once (Persistierung
  idempotent, `count == 2` bleibt wahr); ob der Stand schon über sie hinaus gerückt ist, entscheidet,
  ob die Leerlauf-Bestätigung des Halters (auf ein Keepalive des Servers) vor dem Wartezug greift —
  Zeitverhalten, kein Vertrag. Das passt zum gemeldeten Rot auf PostgreSQL 17 („[94993648 95321768]“:
  zwei Positionen, vermutlich Halter und Retry) und zur grünen Wiederholung.
  **Herleitung, nicht reproduziert** (Code-Lesung plus die Annahme, dass der Server eine Transaktion
  mit Commit-Record-Beginn gleich dem Start-Stand nicht überspringt; nicht gegen den
  PostgreSQL-Quelltext dieses Laufs geprüft, kein eigener `make test-replication`-Lauf). Der Test
  kann bei verzögerter Leerlauf-Bestätigung (Last) rot werden, ohne dass At-least-once verletzt ist.
  Er beeinflusst die E2E-Belege dieses Slice nicht (anderer Tier, `make test-replication`, im Diff
  unverändert).
- `verifizierbar`: ja — `make test-replication` wiederholt (PostgreSQL 17 und 18), Lieferpositionen
  je Versuch gegen `haltedPosition` und `flushAtSecondStart` ausgeben; oder den Wartezug des Tests
  um eine Bestätigungs-Abwartung vor dem Aufbau ergänzen und die Rotquote vergleichen.
- `klasse`: Exaktheits-Aussage strenger als die At-least-once-Zusage (Flake)

### F-3 — Plan-DoD nennt die Position der Negative-Phase anders als Plan §3

- `kategorie`: LOW
- `quelle`: Skill-MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“ (hier LOW: §3
  benennt die Abweichung ausdrücklich)
- `pfad`: Plan `slice-routing-e2e` §2 DoD Negative (iii), letzter Satz („Die Phase läuft als letzte
  vor der Container-Ende-Grenze des Runners“) gegen §3, Zeile `run-integration-tests.sh`
  („Abweichung vom Plan-Wortlaut … letzter Rundlauf des Runners“)
- `befund`: Der DoD-Satz trägt weiter die Position, von der §3 abweicht; §3 verweist auf den
  „Plan-Wortlaut“, ohne die DoD-Zeile zu nennen. Die Abweichung selbst ist tragfähig (siehe
  Antwort c): die Transformations-Abhilfe steht nach derselben Grenze, und die Phase beendet den
  Container dreimal und lässt ihn nach Fall (a) stehen.
- `verifizierbar`: ja — `git grep -n "letzte" -- docs/plan/planning/in-progress/slice-routing-e2e.md`
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-4 — Die Träger der Aussagen „DELETE nicht gemessen“ und „Erreichbarkeit nicht gemessen“ stehen im Spec-Stratum und fehlen im Suchlauf-Feld

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.13 (Suchform: ganzer Baum, Hedge-Wörter) und §3.12
  (Ursprung); Träger außerhalb des Diffs, deshalb LOW mit Meldung an den Planner
- `pfad`: `spec/pflichtenheft.md:404–408` ([`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md): „Die
  Aussage zu `DELETE` ist … hergeleitet und nicht gegen PostgreSQL 17 und 18 gemessen“),
  `:452–460` („Offen und nicht gemessen: ob die Nichtanwendbarkeit der Regel am laufenden System
  entsteht …“), `:1374–1378` ([`SPEC-032`](../../spec/pflichtenheft.md) Abwesenheit, dieselbe
  DELETE-Aussage), `:1415–1418` ([`SPEC-032`](../../spec/pflichtenheft.md) Beispiel „Nicht
  anwendbare Regel“: „ist nicht gemessen“)
- `befund`: Die vier Stellen sind durch die Messung des Slice überholt: DELETE gemessen an
  PostgreSQL 18.6 (eigener Lauf) und 17.11 (übernommen), Nichtanwendbarkeit in (b) und (c) am System
  erzeugt, in (a) als inkompatible Schemaänderung ohne Zeile beendet (Teil „dort genügt das Entfernen
  der Regel nicht“ bleibt hergeleitet, Neustart nach (a) nicht gefahren). Das Suchlauf-Feld des Plans
  sucht nur in `harness docs/user README.md .d-check.yml` nach `Waise`/`CFG-008`/`CFG-007`/
  `abdeckung_declare`/`PG_TEST_IMAGE`, ohne Grund für die Einschränkung und ohne das Hedge-Wort
  „nicht gemessen“; der Übergabe-Block im Plan `slice-routing-betriebsdoku` trägt die Messung dagegen
  mit Version und Lauf (`ROUTING-DELETE-MESSUNG`, 18.6 und 17.11). Die Spec-Stellen sind nicht
  gemeldet. Wird der Nachzug im Spec-Stratum gemacht, trägt die gemessene Aussage ihren Ursprung
  (Version, Lauf, gedruckte Zeile) nach [`AGENTS.md`](../../AGENTS.md) §3.12.
- `verifizierbar`: ja — `git grep -n -E "nicht gemessen|nicht gegen PostgreSQL 17 und 18" -- spec/pflichtenheft.md`
- `klasse`: Träger-Nachzug außerhalb des Diffs (Suchraum des §3.13-Felds enger als der Baum)

### F-5 — Adresse des Zielnamens: Tabellenzelle und Prosa der Spec lesen sich verschieden (Spec-Nachzug)

- `kategorie`: LOW
- `quelle`: [`SPEC-019`](../../spec/pflichtenheft.md) (Fehlertext-Absatz, Routing-Tabelle),
  [`SPEC-032`](../../spec/pflichtenheft.md)
- `pfad`: `spec/pflichtenheft.md:935` (Zeile „`Zielname ist ungültig` | Zielname, wie beantragt“) und
  `:931` (Zeile „`Regelname ist ungültig` | Regelname, wie beantragt“) gegen `:921–923` (Prosa: „die
  Adresse … eines Zielnamens `schema.table.zielname`“)
- `befund`: Die Spalte „Adresse“ nennt „Zielname, wie beantragt“, die Prosa gibt die Form
  `schema.table.zielname` vor; Umsetzung und Test folgen der Prosa
  (`"Zielname ist ungültig: public.<Tabelle>.Eu.Bad"`, `routing_e2e_test.go:406`, im eigenen Lauf grün).
  „Wie beantragt“ meint nach `:860` die Zeichengenauigkeit, nicht die Adressform; das steht für die
  Routing-Tabelle nicht an der Zelle. Nicht Sache des Implementers.
- `verifizierbar`: nein — Lese-Handlung an der Spec; der Test bindet nur die Prosa-Form.
- `klasse`: Spec-Unschärfe zwischen Tabellenzelle und Prosa

### F-6 — Die Godocs der vier `TestE2ERouting*`-Funktionen sind Erzeugnis-Eingabe (Ausnahme greift)

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 „Ausnahme — Erzeugnis-Eingabe ist Abdeckung“
- `pfad`: `test/integration/routing_e2e_test.go:255–263`, `351–359`, `482–488`, `567–574`;
  `docs/user/e2e-abdeckung.md:33–36`
- `befund`: Die vier Kandidaten von `make kommentar-kennungen DIFF=139b25c9~1` sind genau diese Godocs.
  Die Kennungsmengen je Godoc ([`LH-FA-CFG-008`](../../spec/lastenheft.md) +
  [`SPEC-032`](../../spec/pflichtenheft.md); [`LH-FA-CFG-008`](../../spec/lastenheft.md) +
  [`LH-QA-SEC-004`](../../spec/lastenheft.md) + [`SPEC-019`](../../spec/pflichtenheft.md);
  [`LH-FA-CFG-008`](../../spec/lastenheft.md) + [`LH-FA-REA-005`](../../spec/lastenheft.md);
  [`LH-QA-POR-001`](../../spec/lastenheft.md) + [`LH-FA-CFG-008`](../../spec/lastenheft.md)) stehen
  jeweils vollständig in der Kennungsspalte der Zeilen 33–36 des Abdeckungs-Trägers; die Godocs
  tragen keine Kennung, die dort nicht ankommt, und das Subjekt ist der Test (Testfall-Provenienz).
  Mit `TESTS=exclude` meldet das Werkzeug keinen Kandidaten.
- `verifizierbar`: ja — `make kommentar-kennungen DIFF=139b25c9~1` und die Zeilen 33–36 des Trägers.
- `klasse`: Ausnahme Erzeugnis-Eingabe (Bestätigung, kein Fund)

### F-7 — Erreichbarkeit (b) und (c) teilen die Vorbedingung „keine Spaltenform nach der Aktivierung“

- `kategorie`: INFO
- `quelle`: [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  Entscheidung 4; Plan §2 DoD (ii)
- `pfad`: `tools/harness/run-integration-tests.sh` (Phase „Routing-Nichtanwendbarkeit und Abhilfe“,
  Fall (b) und (c), je `rn_shape_rows` gleich 0)
- `befund`: Beide Fälle laufen auf einer Tabelle ohne bekannte Spaltenform (gedruckt 0 Zeilen in
  `cdc.table_schema`); (c) belegt damit die Spaltenliste einer Publication bei Erstaktivierung, nicht
  den Fall bei bereits bekannter Spaltenform (dort ist nach der Herleitung die inkompatible
  Schemaänderung der erste Halt, vgl. Fall (a)). Die Messung (ii) ist erbracht; die Aussage „(c) ist
  unabhängig von (b) erreichbar“ ist es nicht und steht auch nirgends.
- `verifizierbar`: nein — Beleg bräuchte einen weiteren Runner-Fall.
- `klasse`: Erreichbarkeits-Messung deckt eine Vorbedingung, nicht beide Pfade

### F-8 — Datum der Historienzeile 0.14.0 liegt nach dem Commit-Datum

- `kategorie`: INFO
- `quelle`: Maintainability (außerhalb des Gegenstands: Commit `c16d408b` des Hauptlaufs)
- `pfad`: `spec/lastenheft.md` (Änderungshistorie, Zeile 0.14.0: „2026-10-02“)
- `befund`: Der Commit `c16d408b` trägt den Zeitstempel 2026-10-01 17:57, die Historienzeile 2026-10-02.
- `verifizierbar`: ja — `git log -1 --format=%ci c16d408b` gegen die Tabellenzeile.
- `klasse`: Datumsangabe im Träger gegen die Commit-Zeit

### F-9 — Laufzeit und PostgreSQL-17-Lauf sind für diesen Report teils übernommen

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12; Plan §6 Risiko „Laufzeit des Testpakets gegen das
  60-Minuten-Limit von `e2e.yml`“
- `pfad`: Plan `slice-routing-e2e` §6; Plan `slice-routing-betriebsdoku` (Übergabe-Block)
- `befund`: Gemessen (eigener Lauf, PostgreSQL 18.6) sind 10 min 5 s für `make test-integration`
  einschließlich der Routing-Phasen; der Lauf an 17.11 und die Mutationen des Implementers sind
  übernommen. Der Übergabe-Block nennt die Images als `postgres:18-alpine`/`postgres:17-alpine`
  (Tags) und nicht die Digests der Legs aus `.github/workflows/e2e.yml`; die Versionsangabe ist
  gedruckt (`ROUTING-DELETE-MESSUNG PostgreSQL <Version>`), der Digest-Bezug nicht.
- `verifizierbar`: ja — `make test-integration` mit `PG_TEST_IMAGE` auf den 17-Digest.
- `klasse`: Zahl im Träger — Ursprung teils übernommen

---

## Antworten zu den Schwerpunkten (a) bis (i)

- **(a) Trägt jeder DoD-Punkt seinen Satz?** Ja, mit F-1 als einziger Einschränkung.
  Happy Path: fünf Wege belegt; `cdc.changes`, `GET /changes?target=` und `ReadChanges` über Gleichheit
  der Kennungsmengen mit der SQL-Auswahl, die drei Stream-Wege über die persistierte Zeile derselben
  `change_id` (Region und Ziel gegen `cdc.changes`, nicht gegen die Erwartung), der ungefilterte
  Pull-Leser über alle Changes, die Clients ohne Ziel über die Menge {NULL, eu, us}. Boundary R1 bis R6
  samt Formzeilen: je `failed` mit exaktem Klartext und Adresse, Regelstand danach geprüft, Gegenprobe
  `applied`. Neustart: zwei reale `docker restart`, Startzeiten gedruckt. R3 in beide Richtungen und
  [`LH-QA-SEC-004`](../../spec/lastenheft.md): Ausschluss ohne Bedingung entfernt Schlüssel und Wert
  (vor und nach Neustart), Ausschluss der Bedingungsspalte `failed`, Regel auf ausgeschlossene Spalte
  `failed`. Replay: gleiche `change_id`, gleiches Label, vor der Regel erfasste Change `NULL`, Backfill
  trägt das Label (`origin = 'backfill'`); G1/G4 färben die Gegenproben. API-Aktivierung: gRPC und
  HTTP, je mit und ohne Regel. Negative nach [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md):
  (a) keine Zeile (Zählung `new_data->>'id' = '2'` gleich 0), (b)/(c) Spaltenform-Zeilen gedruckt 0,
  Sentinel „Routing-Regel auf die Änderung nicht anwendbar“ im Log, Klasse `schema` in
  `cdc.process_heartbeat`, Abhilfe (iii) mit `pending` bei stehendem Prozess und `applied` nach Neustart;
  die Reihenfolge „vor der ersten Transaktion“ ist am Ausgang belegt (Zeile id=1 ohne Ziel, kein zweiter
  `schema`-Fehler), nicht an einem eigenen Zeitpunkt. DELETE: PostgreSQL 18.6 gefahren, 17.11 übernommen.
- **(b) Grenzen des Belegs.** „Ziel A sieht nur A“ siehe F-1: nicht getragen über die erste Change
  hinaus (an den Stream-Wegen). Flakiness des Happy Path: kein Befund; alle neun Clients melden `READY`
  vor der ersten Einfügung, Dreiergruppen laufen bis zu sechsmal, jede mit 20 × 0,5 s Wartezeit; ein
  Scheitern wäre laut (`bf_fail` nach der sechsten Gruppe), nie ein stilles Grün. Im eigenen Lauf genügte
  eine Gruppe.
- **(c) Position der Negative-Phase.** Tragfähig: die Transformations-Abhilfe steht nach derselben
  Container-Ende-Grenze, die Routing-Phase beendet den Container dreimal und lässt ihn nach (a) stehen;
  vor der Grenze stünden alle späteren Phasen ohne Container. Die Reihenfolge im Runner ist
  Happy Path, Neustart, API-Aktivierung (vor „Prozessstart-Vorlauf-Frist“), Negative als letzter
  Rundlauf (`run-integration-tests.sh:5232`, danach nur die Abdeckungs-Zusammensetzung). Der DoD
  trägt den alten Wortlaut: F-3.
- **(d) Fehlertext-Unschärfe:** F-5 (Spec-Nachzug, nicht Implementer).
- **(e) Überholte „nicht gemessen“-Aussagen:** F-4 mit Zeilen. Die Messaussage im Übergabe-Block trägt
  Version (18.6, 17.11) und Lauf (gedruckte Zeile `ROUTING-DELETE-MESSUNG`); das Spec-Stratum trägt sie
  noch nicht.
- **(f) Flake `TestRunStreamWithRetrySlotStillActive`:** F-2. Bewertung: Test zu streng, nicht der
  Produktcode falsch; Folge-Slice, kein Einfluss auf die E2E-Belege. Die Aussage „genau eine Lieferung“
  hält nur, wenn die Leerlauf-Bestätigung des Halters vor dem Retry greift.
- **(g) `kommentar-kennungen`:** F-6, die Ausnahme greift.
- **(h) Eigene Läufe und Mutationen:** siehe oben; gefahren: sieben Gate- und Messläufe
  (`make test`, `make fmt-check`, `make docs-check`, `make doc-trace`, `make suchlauf-nachmessen`,
  zwei Läufe von `make kommentar-kennungen`) und `make test-integration` einmal (PostgreSQL 18.6);
  Mutationen: 3 gefahren, übernommen: keine.
- **(i) Umfang:** review-tragfähig; die Rückführung „(B) abtrennen“ ist nicht eingetreten, die Negative-Phase
  fügt sich ohne Sonderplatz in den Runner.

## Negativbefunde

- geprüft, ohne Befund: `test/integration/routing_e2e_test.go` (vier Funktionen; Kommentare ohne
  Vorher/Nachher-Sprache und ohne Konjunktiv über verworfene Alternativen; Fehlertexte der Spec
  stimmen mit [`SPEC-019`](../../spec/pflichtenheft.md) überein; G1 und G4 färben die Gegenproben)
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (vier neue Phasen mit je einem
  `abdeckung_declare`; Isolation durch eigene Tabellen je Phase; `-run`-Muster trägt alle vier Funktionen;
  Exit-Codes nicht durch Pipes maskiert; keine Host-Werkzeuge außerhalb der Klasse „ohne Installation“)
- geprüft, ohne Befund: `tools/harness/{httpclient,grpcclient,sseclient,natsstreamsub}` (Flags, Usage-Texte,
  Zählschleifen; Kommentare tragen höchstens eine Kennung; SSE-Scanner wird über Events hinweg
  wiederverwendet)
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` (Erzeugnis des Runners: 8 neue Zeilen für
  [`LH-FA-CFG-008`](../../spec/lastenheft.md), sonst nur verschobene Ortsangaben; Arbeitsbaum nach dem
  eigenen Lauf unverändert)
- geprüft, ohne Befund: `harness/README.md` (Zeilen `make doc-trace`, `make test-integration`; die
  Zahl `0 Waise(n)` stimmt mit der eigenen Messung überein; Parent-Zahl `1 Waise(n)` steht im Plan als
  gemessen)
- geprüft, ohne Befund: Plan-Nachzug `slice-routing-e2e` §3 (Suchlauf-Block: alle 10 Zeilen im eigenen Lauf
  ok; Befunde am Diff entsprechen den Messungen) und Übergabe-Block in `slice-routing-betriebsdoku`
  (Adresse nimmt den Gegenstand an: Kernbegriffe `ROUTING-DELETE-MESSUNG`, Alt-Bild, Abhilfe stehen als
  committeter Text im Plan der Adresse)
- geprüft, ohne Befund: Traceability (alle Commits des Range nennen
  [`LH-FA-CFG-008`](../../spec/lastenheft.md) oder
  [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md), keine
  `SPEC-*`-Kennung im Betreff) und Spec-Stratum (Lastenheft-Änderung `c16d408b` als eigener Commit
  vor den umsetzenden Test-Commit)
- geprüft, ohne Befund: Handbuch-Pflichten (kein Diff an `docs/user/benutzerhandbuch.md`, keine neue
  Betreiber-Oberfläche; Adresse der Betriebsdoku benannt und im Plan angenommen)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht (Aussage „nur“ breiter als die
Messung „erste Change“) · Exaktheits-Aussage strenger als die At-least-once-Zusage (Flake) · Nachzug
widerspricht dem Nachbarn im selben Träger · Träger-Nachzug außerhalb des Diffs · Spec-Unschärfe
zwischen Tabellenzelle und Prosa

## Verdikt

**Merge-blockierend: ja, wegen F-1 (MEDIUM).** Der DoD verlangt „kein offenes HIGH/MEDIUM“. F-1 ist eine
Frage der Belegtiefe an den drei Stream-Wegen, nicht des Produktverhaltens: alle Läufe sind grün, die
Mutationen färben, und die Pull-Wege tragen „nur“ vollständig. F-2 ist ein Befund an einem Test eines
abgeschlossenen Slice, hat keinen Einfluss auf die E2E-Belege dieses Slice und gehört in einen
Folge-Slice; er blockiert diesen Slice nicht. F-3 bis F-5 sind Nachzüge im Plan bzw. im Spec-Stratum
(Planner, Closure), F-6 bis F-9 erfordern keine Aktion.

**Fixrunde:** ja, am Implementer für F-1 (die vom Auftraggeber genannten Alternativen: Negativ-Zählung
mit Wartefenster, oder Zählung aller empfangenen Changes je Client). Die DoD-Zeile „Review durchgeführt“
bleibt deshalb offen; sie wird bei Schritt 21 des Implementer-Workflows nachgezogen.

**Architect-Fragen:**
1. [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md): Soll die Fitness Function
   „Slot noch aktiv“ „genau eine Lieferung im zweiten Versuch“ als Zusage tragen, oder genügt „die
   Retry-Change wird genau einmal persistiert, und eine Wiederzustellung der Halter-Transaktion ist
   zulässig“? Davon hängt ab, ob F-2 den Test lockert oder die Produktseite (Bestätigung vor dem Retry)
   schärft.
2. F-1: Reicht für die drei Stream-Wege „erste Change trifft“ plus die vollständige Pull-Gleichheit als
   Beleg für „sieht nur A“, oder verlangt die Auswahl-Zusage von [`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md)
   („Auswahl, keine Ausblendung“) die Beobachtung über mehrere Changes? Bleibt das Urteil bei MEDIUM, ist der
   Satz in der Abdeckungs-Zeile entsprechend zu verengen.
3. F-7: Soll die Erreichbarkeit der Nichtanwendbarkeit bei bekannter Spaltenform mit Publication-Spaltenliste
   (hergeleitet: erster Halt ist die inkompatible Schemaänderung) Gegenstand eines Belegs werden, oder
   genügt die Herleitung in [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)?

**Übergabe:** F-1 an den Implementer (Fixrunde); F-2 an den Planner als Folge-Slice-Kandidat mit der
Architect-Frage 1; F-3 bis F-5 an den Planner für Closure und Spec-Nachzug (die Messung als gemessen
mit Version 18.6 und 17.11 und gedruckter Zeile `ROUTING-DELETE-MESSUNG` eintragen); F-6 bis F-9 als
Kontext für die Closure-Notiz.
