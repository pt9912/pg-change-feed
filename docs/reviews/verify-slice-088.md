# Verifikationsbericht: slice-088 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-088` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0082`](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
§Konsequenzen (Folgepflicht „test-only"), §Fitness Function und §„Was daraus
für die Slices folgt" (Cluster B),
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand), [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a) (die Rampe)); `AGENTS.md` §3.6, §3.9, §3.10, §3.11. **Nicht** gegen den
Diff als solchen (Reviewer-Aufgabe) und **nicht** gegen realen Bedarf
(Validator, hier nicht ausgelöst).

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan der Fassung `HEAD`
gelesen, dazu `ADR-0082`/`0071`/`0054`, die zwei Review-Reports und die
berührten Artefakte. **Alle** Zahlen dieses Berichts stammen aus eigenen, in
dieser Sitzung gefahrenen Läufen; keine Zahl des Implementers, des Reviewers
oder des Planners wurde übernommen. Die beiden Review-Reports wurden als
**Kontext** gelesen und ihre Findings **nicht** nachgeprüft (Auftrag) — die
eine DoD-Zusage, die auf ihnen aufsetzt („F-1 ist behoben"), ist dagegen
**eigenständig** reproduziert (§4). Exit-Codes je ungepiped und in eigenem
Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung getrennt
beauftragt. Die Mutations-Kopien lagen unter `/tmp`, nicht im Arbeitsbaum;
`git status --porcelain` ist am Ende dieses Laufs **leer** (kein
Lauf-Artefakt, keine Mutation zurückgeblieben).

**Gegenstand:** `slice-088` (Cluster B des Schnittmaßes), geprüfter Stand
`HEAD` = `8506539`, Zweig `main`, Baum sauber. Der Slice liegt weiterhin in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt.
Slice-Range: `3de7fb1..8506539` (elf Dateien: sechs Testdateien, der
Slice-Plan, zwei Review-Reports, zwei `state.md`-Köpfe).

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `git diff 3de7fb1..8506539 -- '*.go' ':!*_test.go'` (`> datei`, Größe geprüft) | **0** | **0 Bytes** — der Diff außerhalb der Testdateien ist leer |
| 2 | `make coverage-gate` (ungepiped, Exit direkt gelesen) | **0** | `coverage-gate: OK — Coverage 74.70% erfüllt Schwelle 70%` |
| 3 | eigenes Profil des Nachher-Stands aus der `coverage`-Stufe gezogen, dedupliziert über die Block-Position | — | Nenner **1903** · gedeckt **1422** · **74,7241 %** |
| 4 | **eigener Vorher-Lauf**: `git archive 3de7fb1` → `docker build --target coverage` → Profil dedupliziert | **0** | Nenner **1903** · gedeckt **1371** · **72,0441 %** |
| 5 | Positions-Vergleich vorher/nachher (Mengendifferenz beider Profil-Schlüssel) | — | **1252** Positionen in *beiden*, Mengen **identisch** — keine Statement-Position kommt hinzu oder fällt weg |
| 6 | `runWALRetentionCheck` (`wiring.go:982…1005`) in beiden Profilen | — | vorher **16/16**, nachher **14/16** — genau die in `ADR-0082` §Kontext (2) dokumentierte ±2-Schwankung |
| 7 | Datei-Deltas gedeckter Statements | — | `replication/decode/decode.go` 75/89 → 88/89 (**+13**) · `replication/mapper/mapper.go` 166/191 → 184/191 (**+18**) · `sqlexec/translate.go` 130/147 → 147/147 (**+17**) · `postgresstorage/mapper/mapper.go` 16/20 → 20/20 (**+4**) = **+52** · `domain/model/schema_version.go` 4/5 → 5/5 (**+1**) · `bootstrap/wiring.go` 195/521 → 193/521 (**−2**) |
| 8 | ungedeckte Statements in den **vier Trägerdateien**, dedupliziert und auf `stmts > 0` gefiltert | — | vorher **60** · nachher **8**, an **genau** den acht genannten Positionen (Liste §3) |
| 9 | `make gates` | **0** | baseline-verify `v6.5.0` OK · docs-check 0 Befunde · a-check 0 Befunde · commit-traceability OK · coverage-gate OK |
| 10 | `make test` (`--network none`, `-race`) | **0** | **31** Pakete `ok`, 0 `FAIL` |
| 11 | `make doc-commits RANGE=3de7fb1..8506539` | **0** | 718 Dateien geprüft, **0** Befunde |
| 12 | `make doc-immutable RANGE=3de7fb1..8506539` | **0** | 718 Dateien geprüft, **0** Befunde |
| 13 | `make doc-planning` (Zusatz, nicht in `make gates`) | **0** | 718 Dateien geprüft, 0 Befunde — Lifecycle Roadmap ↔ `in-progress` konsistent |
| 14 | Mutationskopie des `HEAD`-Stands, Baseline-Lauf (decode + mapper) | **0** | `ok` / `ok` — Ausgangspunkt unmutiert grün |
| 15 | **eigene Mutation A** an `tupleValues`s `default:`-Zweig (Binär-Tupel wie Text behandelt) | **1** | genau `TestDecodeInsertRejectsBinaryTupleValue` rot |
| 16 | **eigene Mutation B** an `qualifiedNames` (`relation.Name` statt `QualifiedName()`) | **1** | genau `TestConsumeTruncateNamesQualifiedRelations` rot — Meldung `TRUNCATE an feed, orders` statt `public.feed, public.orders` |
| 17 | **M4-Wiederholung** (Riegel `if !activated { return nil }` entfernt) | **1** | genau `TestConsumeRelationOnUnboundTableStaysNoop` rot — `Relation außerhalb der Aktivierungen: leere Kennung` |
| 18 | **M5-Wiederholung** (Riegel `if !found { return nil }` entfernt) | **1** | genau `TestConsumeRelationWithoutRegisteredVersionStaysNoop` rot — `TableSchema [""], RegisterVersion [{tbl-1-v1 tbl-1 1}]` |
| 19 | `go tool cover -func` der `coverage`-Stufe, je genannte Funktion | — | alle in §2 genannten Funktionen **100,0 %**, außer `Decode` 98,0 %, `Consume` 90,6 %, `change` 92,6 %, `rowImage` 90,9 % — je genau die in §6 als unerreichbar gemeldeten Positionen |
| 20 | `git status --porcelain` am Ende | — | leer |

**Die zwei Enden der Schwankung, ausdrücklich:** mein Vorher-Lauf liegt auf
dem **oberen** Ende (1371), mein Nachher-Lauf auf dem **unteren** (1422). Der
gleichende Abstand ist in **beiden** Paarungen **+53** (1369 + 53 = 1422 ·
1371 + 53 = 1424). Der Plan berichtet das **obere** Paar (1371 → 1424,
72,0 % → 74,8 %); das Paar ist in sich konsistent, und das Delta hängt
**nicht** an der Schwankung — genau das war zu prüfen (Auftrag Punkt 3).

---

## 2. DoD-Konformität, Kriterium für Kriterium

### Liefer-Punkt 1 — die Stream-Übersetzung ist gedeckt

- **bestätigt** — `replication/decode` und `replication/mapper` haben Tests für
  die gemessenen Lücken. Beleg: Messung 8 (die 14 bzw. 25 vorher ungedeckten
  Statements beider Pakete sind geschlossen, soweit erreichbar) und Messung 19
  (`tupleValues`, `oldTupleValues`, `observeRelation`, `classifyRelationColumns`,
  `IncludeColumn`, `removeExcluded`, `setSchemaVersion`, `qualifiedNames` je
  **100,0 %**; `Decode`, `Consume`, `change`, `rowImage` je 100 % bis auf die in
  §6 benannten, nachgemessen unerreichbaren Positionen).
- **bestätigt** — `JSONImage` (vorher **3/3 ungedeckt**) ist gedeckt:
  `postgresstorage/mapper/mapper.go:96` steht bei **100,0 %**, das Paket bei
  **20/20**. Die Zuordnung zu Liefer-Punkt 3, die der Aufzählungspunkt
  behauptet, trifft zu (die Funktion liegt in `postgresstorage/mapper`).
- **bestätigt** — „netzlos, ohne externe Dienste": Messung 10 fährt die volle
  Testfläche in einem Container mit `--network none` und `-race`, Exit 0.

### Liefer-Punkt 2 — die SQL-Übersetzung ist gedeckt

- **bestätigt** — die sieben genannten Funktionen (`ReadConsumerPositions`,
  `ReadExcludedColumns`, `ReadTableSchema`, `ReadSourceTables`,
  `ReadPendingRequests`, `ReadChanges`, `ReadConsumerPosition`) stehen je bei
  **100,0 %** (Messung 19); `sqlexec/translate.go` ist **147/147** (Messung 7).
- **bestätigt** — „über den **bestehenden** Fake der Naht aus `ADR-0080`": die
  Fakes (`fakeRows`, `fakeExecutor`, der Block `--- Fakes der Naht ---`) liegen
  bereits im Stand `3de7fb1` (`git show 3de7fb1:…sqlexec/translate_test.go`);
  der Slice ergänzt Nutzungen der vorhandenen `scanErrs`/`iterErr`-Felder und
  legt **keinen** zweiten Fake an.

### Liefer-Punkt 3 — die Wirkung ist gemessen, nicht angestrebt

- **bestätigt** — `postgresstorage/mapper`s Lücken (`JSONImage`, `ToChange`)
  sind gedeckt: beide **100,0 %**, Paket **20/20**.
- **bestätigt mit Reserve** — „Der Effekt ist beziffert: die Gate-Zahl
  vorher/nachher, gemessen über `make coverage-gate` — erwartet ≈60
  Statements, und der Nenner bleibt **1903**". Der Nenner **ist** 1903, in
  meinem Vorher- **und** Nachher-Profil (Messungen 3/4), und die
  Positionsmengen sind **identisch** (Messung 5) — die Aussage „dieser Slice
  fügt keinen Produktionscode hinzu" ist damit nicht nur am Diff, sondern an
  der Messung selbst bestätigt. Das Delta ist **+53** (Messungen 7/8 und die
  Enden-Rechnung in §1) und die Aufteilung **+52** in den vier Trägerpaketen,
  **+1** mittelbar in `domain/model` geht positionsgenau auf. Die **Reserve**:
  die im Plan berichtete Nachher-Zahl **74,8 % / 1424** ist das *obere* Ende
  der dokumentierten Schwankung; mein eigener Nachher-Lauf druckt **74,70 % /
  1422** und liegt damit am *unteren* (Messung 6). Das Delta ist von der Wahl
  des Endes unabhängig — die Zahl selbst ist es nicht. Siehe §5 Punkt 3.
- **bestätigt** — `make gates` grün (Messung 9). Anmerkung: die
  `coverage`-Stufe des Gates fährt `go test` über die netzlos prüfbare Fläche
  und bricht bei einem roten Test ab — „`make gates` grün" trägt deshalb auch
  die Testfläche ohne Race-Detector; `make test` (mit `-race`) ist zusätzlich
  grün (Messung 10).

### Review-Zeile und Closure-Posten

- **bestätigt** — „Review durchgeführt, Report unter `docs/reviews/` liegt
  vor … F-1 im Fixrunden-Commit behoben": beide Reports liegen vor (Erstlauf
  und Delta); die Behauptung „F-1 behoben" ist **eigenständig** reproduziert
  (Messungen 17/18, §4).
- **offen, und zwar richtig** — Closure-Notiz mit Lerneintrag. Planner-Pflicht
  beim Übergang nach `done/`.
- **entfällt** — „Reconciliation-Register (`../reconciliation.md`)
  fortgeschrieben, **falls** dieser Slice einen Inventur-Fund auflöst".
  `docs/plan/planning/reconciliation.md` existiert **nicht**
  (`ls` auf `docs/plan/planning/` → `done/ in-progress/ next/ observations/
  open/ README.md welle-20.md`); das Repo ist Greenfield (Modus-Deklaration
  `*`/`PGC`), also greift die im Plan selbst genannte Ausnahme. Das Item ist
  **korrekt unbehakt**.
- **offen, und zwar richtig** — Beobachtungs-Register fortgeschrieben;
  Planner-Pflicht bei der Closure (§7).
- **offen, und zwar richtig** — jedes Risiko aus §6 mit einem Ausgang; alle
  drei stehen auf `<bei Closure>`.
- **offen, und zwar richtig** — die drei Paarungen. Dieses Repo **hat**
  Wellen-Betrieb (`welle-20` ist offen), also trägt sie der Lese-Schritt der
  `welle-20`-Closure, nicht dieser Slice.

**Kein DoD-Kriterium ist „nicht bestätigt".** Das einzige Kriterium mit einer
Einschränkung ist die bezifferte Wirkung; die Einschränkung ist benannt und
ändert das Kriterium nicht (§5 Punkt 3).

---

## 3. Der Befund „8 Statements unerreichbar" — zwei Angaben selbst nachgeprüft

**Die Positionenliste trägt exakt.** Messung 8 liefert nach der Deduplizierung
und dem Filter `stmts > 0` in den vier Trägerdateien genau acht ungedeckte
Block-Positionen, und zwar dieselben acht, die §6 nennt:

| Position | Funktion | Begründung des Plans | Meine Prüfung |
|---|---|---|---|
| `decode.go:224.3` | `Decode` | der `default:`-Zweig | **trägt — selbst nachgeprüft** |
| `mapper.go:158.4` | `Consume` | `NewOpenTransaction` | trägt (strukturell, §unten) |
| `mapper.go:172.4` | `Consume` | `Commit` | trägt (strukturell) |
| `mapper.go:191.4` | `Consume` | `AppendChange` | trägt (strukturell) |
| `mapper.go:235.3` | `change` | `rowImage` Neu-Image | trägt (folgt aus 537/541) |
| `mapper.go:239.3` | `change` | `rowImage` Alt-Image | trägt (folgt aus 537/541) |
| `mapper.go:537.4` | `rowImage` | `json.Marshal(column.Name)` | **trägt — selbst nachgeprüft** |
| `mapper.go:541.4` | `rowImage` | `json.Marshal(*values[i])` | **trägt — selbst nachgeprüft** |

**(a) `decode.go:224` — der `default:`-Zweig.** Nachgeprüft an der gepinnten
Quelle des Treibers, nicht an der Aussage des Implementers: `pglogrepl.Parse`
(`message.go:675`) wählt über einen `switch` genau sieben Dekoder
(`Relation`, `Type`, `Insert`, `Update`, `Delete`, `Truncate`, `LogicalDecoding`)
und fällt sonst auf `getCommonDecoder`, das genau drei weitere kennt (`Begin`,
`Commit`, `Origin`); ist der Dekoder danach `nil`, endet `Parse` **vor** der
Rückgabe mit `errMsgNotSupported`. Die Fallunterscheidung in `Decode`
(`decode.go:142-222`) führt **genau diese zehn** Typen und leitet jeden anderen
Fall über den `err != nil`-Riegel bei `:139/140` ab. Der `default:`-Rumpf ist
damit über die öffentliche Fläche nicht auslösbar. **Trägt.**

**(b) `mapper.go:537/541` — `json.Marshal` eines `string`.** Empirisch gegen
die gepinnte Toolchain geprüft (`golang:1.27-alpine`, derselbe Digest wie
`TOOLCHAIN_IMAGE`): zwölf Eingaben — leerer String, Anführungszeichen,
Zeilenumbruch/Tabulator, **ungültiges UTF-8**, eingebettetes NUL,
Surrogat-Byte-Folge, 1 MiB Länge, Steuerzeichen, DEL, HTML-Zeichen, Emoji —
ergeben **in allen zwölf Fällen `err=<nil>`**. `column.Name` ist ein `string`,
`*values[i]` dereferenziert `[]*string` ebenfalls zu `string`. **Trägt.**

**(c) `mapper.go:158/172/191` — strukturell nachgeprüft.** `NewOpenTransaction`
(`domain/model/transaction.go:24`) scheitert nur an leerer Kennung; die Kennung
entsteht aus `strconv.FormatUint` (nie leer), und `a.source` ist im Konstruktor
(`mapper.go:127`) auf nicht-leer geprüft. `Commit` scheitert nur an
`t.committed` (nach dem Commit setzt `Consume` `a.open` auf `nil`; eine zweite
COMMIT-Nachricht läuft in `ErrCommitWithoutBegin`) oder an fremder Quelle (die
Position entsteht aus `a.source`). `AppendChange` scheitert nur an
`t.committed`, an fremder Transaktions-Kennung (der Change trägt
`a.open.tx.ID`, `mapper.go:244`) oder an doppelter Sequenz
(`a.open.sequence` zählt monoton, `:218`). **Trägt.**

**Keine der acht Angaben ist in Wahrheit erreichbar** — es fehlt also **kein**
Test. Zusätzlich: die acht Positionen bleiben im Nachher-Profil **ungedeckt**
(Messung 8) — der Befund ist nicht „grün gefärbt". **Das Urteil des Reviews
zum 8-Statements-Befund halte ich für bestätigt.**

---

## 4. Die Testbindung — F-1s Fix und zwei eigene Mutationen

**F-1s Fix hängt jetzt an der Eingabeseite.** Gelesen (nicht nur behauptet):
`stubSchemaStore` (`mapper_test.go:940-966`) protokolliert die **empfangenen
Kennungen** in `readTables`/`readVersions`/`registrations`; die beiden
No-op-Zusagen lesen sie aus —
`TestConsumeRelationOnUnboundTableStaysNoop` verlangt
`len(readTables) == 0 && len(readVersions) == 0 && len(registrations) == 0`
(`:1088`), `TestConsumeRelationWithoutRegisteredVersionStaysNoop` verlangt
`readTables == ["tbl-1"]` **und** `readVersions`/`registrations` leer
(`:1135/:1138`). Beide Zusagen binden damit *ob* und *womit* der Store befragt
wurde — die Lücke, die der Erstlauf beschrieb, ist an der Form geschlossen.

**Eigenständig reproduziert (Messungen 17/18):** die beiden Mutationen, die den
Erstlauf noch grün ließen, färben jetzt **genau je einen** der beiden Tests
rot, mit der jeweils vorhergesagten Meldung. Die DoD-Zeile „F-1 im
Fixrunden-Commit behoben" trägt damit auf eigener Messung, nicht auf dem
Delta-Report.

**Zwei eigene Mutationen an Zusagen, die *niemand* mutiert hat** (Auftrag
Punkt 5) — beide an der **Eingabeseite**, beide an einem Zweig außerhalb der
Mutationsliste des Implementers und des Reviewers:

| # | Mutation (Eingabeseite) | Exit | Was rot wird |
|---|---|---|---|
| **A** | `tupleValues`: der `default:`-Zweig behandelt ein **Binär-Tupel** wie Text (`values[i] = &value` statt Fehler, `decode.go:273-274`) | **1** | genau `TestDecodeInsertRejectsBinaryTupleValue` — die neue Zusage bindet den eingehenden Tupel-Typ |
| **B** | `qualifiedNames`: `relation.Name` statt `relation.QualifiedName()` (`mapper.go:562`) | **1** | genau `TestConsumeTruncateNamesQualifiedRelations` — `TRUNCATE an feed, orders` statt `public.feed, public.orders`; die Zusage bindet den qualifizierten Namen an die Eingabe-Relation |

Beide sind **saubere Einzelpunkt-Sonden**: die Mutation ändert genau eine
Zeile, es fällt genau der vorhergesagte Test. **Präzisierung zu A:** die
Mutation färbt `TestDecodeDeleteRejectsBinaryKeyValue` *nicht* rot — das ist
**korrekt**, nicht eine zweite Lücke: der `K`-Tupel-Pfad läuft über
`oldTupleValues` (`decode.go:309`) und damit über eine **andere**
Statement-Position, die von A nicht berührt wird.

**Nebenprodukt — die Klasse ist nicht nur an den zwei nachgebesserten
Stellen getilgt:** §6 führt
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (3×) als tragende Klasse.
Über die Mutationen dieses Laufs (M4/M5 für die zwei nachgebesserten Zusagen,
A/B für zwei weitere) und die des Deltas (M6/M7 für Aufruf- und
Argument-Bindung) tragen die geprüften neuen Zusagen der zwei
`replication`-Pakete ihre Aussage an der Eingabeseite. **Der Ausgang dieser
Klasse gehört gleichwohl dem Lese-Schritt der `welle-20`-Closure**, nicht
diesem Bericht (Modul 6).

---

## 5. Entscheidungs-Konformität

**1. `ADR-0082` §Konsequenzen, Folgepflicht „test-only".** Der Diff
`3de7fb1..8506539` berührt **keine** Produktionsdatei (Messung 1: 0 Bytes) —
und über die Messung hinaus auch keine, die außerhalb `*_test.go` läge: die
elf geänderten Dateien sind sechs `_test.go`, der Slice-Plan, zwei
Review-Reports und zwei `state.md`-Köpfe. **Eingelöst.**

**2. `ADR-0082` §„Was daraus für die Slices folgt" — Cluster B und seine
Zahl.** Cluster B führt vier Pakete mit **60** ungedecktem; mein Vorher-Profil
zählt in genau diesen vier Trägerdateien **exakt 60** ungedeckte Statements
(Messung 8, dedupliziert, `stmts > 0`). Die Zahl des Schnittmaßes ist
**nachgemessen**, nicht übernommen. **Eingelöst.**

**3. `ADR-0082` §Fitness Function, Zeile 3 — „Nenner und Zähler sind
Zustandsgrößen".** Der Nenner ist 1903 in beiden Läufen (Messungen 3/4); der
Zähler schwankt um 2 über denselben Quelltext, und die Schwankung liegt in
genau der von der ADR benannten Funktion (`runWALRetentionCheck`, Messung 6).
**Eine Reserve zur Zahl des Plans, kein Verstoß:** `ADR-0082` §Kontext (2)
verlangt, wo beide Läufe auseinandergehen, **beide Enden** zu nennen. Der
Slice-Plan nennt für den **Vorher**-Lauf beide Enden („72,0 %" und „der
frühere Lauf desselben Stands druckte 71,9 %"), für den **Nachher**-Lauf nur
das obere (74,8 %) — mein Lauf reproduziert das untere (74,70 %). Das
**Delta** ist von der Endenwahl unabhängig (+53, §1) und damit die belastbare
Aussage; die berichtete Nachher-Zahl ist es nicht. Das ist die Klasse
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**3×**, Schwelle
erreicht) — der **Lese-Schritt der `welle-20`-Closure** sollte diesen Punkt
mitlesen. Ich melde ihn **nicht** als DoD-Verletzung: das Kriterium verlangt
„beziffert" und „Nenner 1903", und beides trägt.

**4. `ADR-0082` Festlegungen 1/4 — der Composition Root bleibt unangetastet.**
`internal/bootstrap/wiring.go` ist im Diff **nicht** enthalten (Messung 1,
Dateiliste). Die **−2** an `runWALRetentionCheck` (Messung 7) ist die
dokumentierte Mess-Schwankung, keine Code-Änderung — bei
positionsgleichem Profil (Messung 5) kann sie keine sein. **Eingelöst.**

**5. `ADR-0071` (Messgegenstand) und `ADR-0054` §(a) (Rampe) — unberührt.**
`Dockerfile`, `harness/mk/coverage.mk`, `tools/coverage-gate.sh` und
`.a-check.yml` sind **nicht** im Diff; `THRESHOLD` steht unverändert bei 70.
Keine Schwellen- oder Gegenstandsänderung, also kein `AGENTS.md` §3.6-Fall.
Die vier Liefer-Punkte und §1 der Abgrenzung („Der Messmechanismus. Gegenstand,
Endstufe und Rampe sind unberührt") treffen zu. **Eingelöst.**

**6. `AGENTS.md` §3.9 (Exit ungepiped).** Alle 20 Messungen dieses Berichts
wurden mit getrennter Exit-Code-Lesung gefahren (`… > log 2>&1` oder
`…; ec=$?`); kein Gate-Aufruf lief durch eine Pipe; Gate-Lauf und
Folgehandlung sind getrennte Werkzeug-Aufrufe. **Eingelöst.**

**7. `AGENTS.md` §3.11 (kein host-lokaler Pfad).** Der Slice fügt keine Doku
hinzu, die eine solche Form tragen könnte; das `hostpaths`-Modul ist Teil von
`make docs-check` und lief grün (Messung 9). **Eingelöst.**

---

## 6. Die Register-Köpfe

Der Planner hat in `8506539` **zwei** `state.md`-Köpfe berichtigt. Ich habe
**alle vier** Einträge nachgezählt, die dieser Slice berührt — die drei in §8
genannten und den in `8506539` mitberichtigten:

| Eintrag (`BEO-PGC/…`) | Belege in `evidence/` | Kopfzeile `state.md` | Übereinstimmung |
|---|---|---|---|
| `negativtest-ohne-bindung-an-seine-eingabe` | 3 (`slice-083`, `086`, `087`) | „offen — **3× erreicht**" | **ja** |
| `db-gegenstand-enthaelt-netzlos-geprueften-code` | 2 (`slice-084`, `085`) | „offen (**2×**)" — in `8506539` von „1×" berichtigt | **ja** |
| `zahl-in-traeger-driftet-gegen-die-messung` | 3 (`slice-081`, `084`, `085`) | „offen — **3× erreicht**" | **ja** |
| `dod-begruendung-unzutreffende-tatsachenbehauptung` | 4 (`slice-036`, `081`, `082`, `083`) | „offen — **4× erreicht**" — in `8506539` von „3×" berichtigt | **ja** |

Beide Berichtigungen sind **mechanisch richtig** (`ls evidence/` gegen die
Kopfzeile), und die §8-Zähler des Plans (**3× / 2× / 3×**) stimmen mit den
Köpfen und den Beleglisten überein. Die frühere Divergenz — Kopfzeile gegen
den eigenen Satz „Zähler (abgeleitet)" in **derselben** Datei — ist an beiden
Stellen geschlossen.

**Randbefund außerhalb des DoD:** `BEO-PGC/architect-verdikt-ablageort-uneinheitlich`
hat **kein** `evidence/`-Verzeichnis, während die maschinelle Hälfte der
Register-Paarung (Modul 6) „jedes Verzeichnis ein nicht leeres `evidence/`"
verlangt. Der `state.md` **dokumentiert** das ausdrücklich („Zähler
(abgeleitet): 0× — kein abgeschlossener Vorgang … die Behebung lief wellenlos
direkt aus der Nutzerfrage"), ist also **benannt, nicht still** — und der
Eintrag ist weder von diesem Slice berührt noch zitiert. Kein Befund zu
diesem Slice; ein Objekt des Lese-/Paarungs-Schritts der `welle-20`-Closure.

---

## 7. `AGENTS.md` §3.10

**§3.10 greift nicht.** `git diff --name-only 3de7fb1..8506539 -- .github/workflows/`
ist **leer** (Exit 0, keine Ausgabe); keine der elf geänderten Dateien liegt
unter `.github/workflows/`, und `compose.yaml`, `Dockerfile` und `ci.yml` sind
unberührt. Der Slice führt **keinen** neuen Workflow und **keine** strukturelle
Änderung an einem bestehenden ein — die Hard Rule hat hier **kein Objekt**.
**Bestätigt, nicht angenommen.**

---

## 8. Was ich nicht prüfen konnte

- **Den Inhalt der beiden Review-Reports** habe ich **nicht** nachgeprüft
  (Auftrag); ich habe sie als Kontext gelesen und **eine** ihrer tragenden
  Aussagen — „F-1 ist behoben" — eigenständig reproduziert (§4). F-4 (INFO)
  habe ich übernommen, nicht bewertet.
- **Die zwei Weißbox-Dateien** (`tuplevalues_internal_test.go`,
  `schemaversion_internal_test.go`) habe ich **gelesen**, aber **nicht**
  mutiert. Ihre Legitimation (unexportierte nil-Zweige, die der öffentliche
  Weg nicht herreicht) und ihre Zusage konnte ich damit nur strukturell, nicht
  durch eigene Mutation bestätigen. **Grenze dieses Laufs.**
- **Die drei `sqlexec`-Klassen-Zusagen** (§„Klasse + Ursache + Aufruf des
  Übersetzungspunkts") habe ich über die 100-%-Coverage und die unveränderte
  Fake-Herkunft bestätigt, **nicht** durch eigene Mutation nachgefahren — der
  Erstlauf hat das mit M2 getan, und der Auftrag verlangte eine Mutation an
  einer *anderen* Stelle (geleistet: A/B).
- **Der Post-Push-Lauf auf dem GitHub-Runner** (`AGENTS.md` §3.10): für diesen
  Slice **nicht einschlägig** — kein Workflow berührt (§7). Die strukturelle
  Grenze (ein Docker-only-Sensor kann sie nicht leisten) bleibt bestehen, hat
  hier aber kein Objekt.
- **Die §6-Risiko-Ausgänge, die Closure-Notiz, das Beobachtungs-Register und
  die drei Paarungen** sind **Planner**-Pflichten beim Übergang nach `done/`
  und hier bewusst noch offen. Der **Lese-Schritt** der `welle-20`-Closure
  (welcher 3×-Eintrag bekommt welchen Ausgang) ist **nicht** mein Gegenstand —
  ich habe nur die Zähler-Stände nachgezählt (§6).
- **Die `welle-20`-Closure-Belege selbst** (Gate grün bei `THRESHOLD=80`, rot
  bei `THRESHOLD=85`) sind Wellen-Arbeit; ich habe nur den gegenwärtigen
  Stand (74,70 % bei 70 %) gemessen.
- **Das DB-gestützte Tier** (`make test-store`, `make test-replication`,
  `make test-integration`) habe ich **nicht** gefahren: der Slice fügt nur
  netzlose Unit-Tests hinzu, und seine DoD fordert dieses Tier nicht. Dass die
  reale SQL-Deckung dort unberührt bleibt, ist damit **nicht** von mir belegt —
  aber auch nicht zugesagt (der Test-Kommentar in `translate_test.go:4-6`)
  sagt es selbst ausdrücklich ab.

---

## 9. Verdikt

**Der Slice ist closure-fähig.**

Alle sechs prüfbaren DoD-Kriterien aus §2 tragen; das siebte
(Reconciliation-Register) entfällt nachweislich, und die fünf übrigen sind
naturgemäß Planner-Pflichten beim `git mv` nach `done/`. Getragen ist alles,
was der Slice zusagt:

- **kein Produktionscode** — der Diff außerhalb `*_test.go` ist **0 Bytes**
  (Messung 1), und die Messung selbst bestätigt es: die Positionsmenge des
  Profils ist vorher und nachher **identisch** (Messung 5);
- **der Nenner bleibt 1903** — in **beiden** eigenen Läufen (Messungen 3/4);
- **die Wirkung ist beziffert und reproduziert** — **+52** in den vier
  Trägerpaketen (4 + 17 + 13 + 18) und **+1** mittelbar in `domain/model`
  (Messung 7); das Delta **+53** hängt nicht an der ±2-Schwankung
  (Messung 6; die Enden-Rechnung in §1). Die berichtete Nachher-Zahl **74,8 %**
  ist das obere Ende der Schwankung, mein Lauf druckt **74,70 % / 1422** —
  benannte Reserve, kein Verstoß (§5 Punkt 3);
- **der 8-Statements-Befund trägt** — genau acht ungedeckte Positionen über
  (Messung 8), zwei davon **selbst nachgeprüft** (`decode.go:224` an der
  Treiberquelle, `mapper.go:537/541` empirisch), keine davon in Wahrheit
  erreichbar (§3);
- **die Testbindung trägt** — F-1s Fix ist eigenständig reproduziert
  (Messungen 17/18), und **zwei eigene Mutationen an bisher unmutierten
  Zusagen** färben je genau ihren Test rot (Messungen 15/16). Es fehlt damit
  **kein** Test, und die tragende Klasse ist an den geprüften neuen Zusagen
  gebunden;
- **die Register-Köpfe stimmen** — vier Einträge nachgezählt, Kopf = Belegliste
  (§6);
- **`make gates` grün** (Messung 9), **`make test` grün** (`--network none`,
  `-race`, 31 Pakete, Messung 10), `doc-commits`/`doc-immutable`/`doc-planning`
  über die Slice-Range grün (Messungen 11-13);
- **`§3.10` nicht ausgelöst** — kein Workflow berührt (§7).

**Was offen bleibt:** (1) die Planner-Closure-Posten (§7-Notiz, Register,
Risiko-Ausgänge, die drei Paarungen über die `welle-20`-Closure);
(2) die **Nachher-Zahl des Plans** trägt nur als Enden-Aussage, nicht als
Einzelwert — für den Lerneintrag der Closure relevant, **nicht**
closure-blockierend; (3) `BEO-PGC/architect-verdikt-ablageort-uneinheitlich`
als Randbefund außerhalb dieses Slice.

**Kein Rückgabe-Pfeil an den Implementer und keiner an den Planner.**

---

**Übergabe:** Dieser Bericht geht an den **Planner** (Closure-Entscheidung;
§5 Punkt 3 und §6 als Material für den Lerneintrag und den Register-Lese-Schritt).
Er ist ein **Verifikations**-Artefakt (Modul 11) und ersetzt weder den
Review-Report (Diff gegen Plan/ADR) noch die Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst).
