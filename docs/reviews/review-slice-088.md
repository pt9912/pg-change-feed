# Review-Report: slice-088 — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `slice-088`, Diff `3de7fb1..7bf2aa5` (vier Commits, HEAD
`7bf2aa5`) — sieben Dateien: sechs Testdateien und der Slice-Plan selbst.
**Kein Produktionscode**: `git diff 3de7fb1..7bf2aa5 -- '*.go' ':!*_test.go'`
ist leer (Exit 0).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-088` (Cluster B) vollständig (§1–§8); offene Welle
  `welle-20`
- `ADR-0082` (das Schnittmaß: Tail statt Composition Root; die Zählbasis in
  §Kontext (2), die Cluster-Tabelle, die Folgepflicht „test-only" in
  §Konsequenzen), `ADR-0071` (Messgegenstand), `ADR-0080` (Nahtform der
  Adapter), `ADR-0059` Teilfrage 3 (Ausschlussstand), `ADR-0029` (Domänen-
  Invarianten), `ADR-0015` (dynamische Re-Versionierung), `ADR-0054`
  §(a) (Rampe)
- `LH-FA-CFG-001`, `LH-FA-CFG-005`, `LH-FA-SCH-004.a`, `LH-FA-CAP-008`,
  `LH-FA-DAT-003`, `SPEC-001`…`SPEC-003`, `SPEC-008`
- `AGENTS.md` (Hard Rules, insbesondere §3.7, §3.9), `harness/conventions.md`
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/`
  (Zähler an den drei in §8 genannten Einträgen)
- letzte Review-Reports desselben Gegenstandsbereichs: `review-slice-083`,
  `review-slice-086`, `review-slice-087`

---

## Findings

### F-1 — Zwei No-op-Zusagen sind nicht an ihre Eingabeseite gebunden

- `kategorie`: MEDIUM
- `quelle`: `LH-FA-CFG-001` („CDC erfasst nur aktivierte Tabellen") ·
  Maintainability · die Klasse `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
  (3×)
- `pfad`: `internal/adapters/driving/replication/mapper/mapper_test.go:1056-1064`
  (`TestConsumeRelationOnUnboundTableStaysNoop`), `:1095-1108`
  (`TestConsumeRelationWithoutRegisteredVersionStaysNoop`); Produktionsseite
  `internal/adapters/driving/replication/mapper/mapper.go:344-354`
- `befund`: Beide Tests behaupten die **Abwesenheit** eines Fehlers und
  stellen dafür gegen einen Stub, der seine Argumente ignoriert
  (`stubSchemaStore` liefert `current`/`schema` unabhängig von `SourceTableID`
  und `SchemaVersionID`). Gemessen: Mutation M4 (in `observeRelation` den
  Riegel `if !activated { return nil }` entfernt) und Mutation M5 (den Riegel
  `if !found { return nil }` entfernt) lassen
  `go test ./internal/... ./cmd/...` **je mit Exit 0** enden — die gesamte
  netzlos prüfbare Fläche bleibt grün, obwohl M4 real Verhalten ändert: mit
  einer Relation, deren Form die bekannte erweitert, endet `Consume` für eine
  **nicht aktivierte** Tabelle unter M4 mit einem sichtbaren Fehler
  (`Relation ausserhalb der Aktivierungen: leere Kennung`), wo der
  unveränderte Stand wirkungslos bleibt (Exit 0). Die beiden Tests können
  diesen Unterschied nicht sehen: ihr Eingabefall klassifiziert als
  „unverändert" (M4) bzw. läuft mit einer leeren bekannten Form weiter (M5),
  und keiner der beiden beobachtet, ob der Store überhaupt befragt wurde.
  Die in §6 als tragend benannte Klasse liegt damit **auch in diesem Slice**
  vor — an zwei der 16 neuen Tests. Die übrigen 14 neuen Tests dieses Pakets
  sind gebunden (Mutationen M1–M3 färben sie rot, siehe §Eigene Messungen).
- `verifizierbar`: ja — Mutationen M4/M5 dieses Laufs (Exit 0 statt rot) und
  die Gegenprobe (dieselbe Zusage an ihre Eingabe gebunden: Exit 1 unter M4,
  Exit 0 im unveränderten Stand)
- `klasse`: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`

### F-2 — §1 nennt fünf verlagerte Funktionen; gemessen sind es acht

- `kategorie`: MEDIUM
- `quelle`: Maintainability · die Klasse
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` („dort driftet eine
  **geschriebene Zahl** gegen die Messung", 3×)
- `pfad`: `docs/plan/planning/in-progress/slice-088-coverage-tail-uebersetzung.md:56-61`
- `befund`: Die Berichtigung sagt „eine frühere Fassung dieser Tabelle stellte
  **fünf Funktionen** ins falsche Paket" und zählt dabei sechs Namen auf, von
  denen einer (`ToChange`) gar nicht betroffen ist. Gegen den Diff der Tabelle
  (`0f21eb2..7bf2aa5`) gemessen haben **acht** Funktionen das Paket gewechselt:
  `observeRelation` (decode → mapper), `oldTupleValues` und `tupleValues`
  (mapper → decode), `JSONImage` (mapper → postgresstorage/mapper),
  `IncludeColumn`, `setSchemaVersion`, `removeExcluded`, `qualifiedNames`
  (sqlexec bzw. postgresstorage/mapper → mapper). Die Aufzählung nennt davon
  vier, nennt `ToChange` zusätzlich und lässt die drei Funktionen aus, die aus
  `replication/mapper` heraus- bzw. in `replication/decode` hineingewechselt
  sind. Die **Summen-Spalte** der neuen Tabelle selbst ist dagegen korrekt:
  je Paket 14 · 25 · 17 · 4 = **60**, und die Gesamt-Statements je Funktion
  stimmen mit dem Profil überein (`Decode` 50, `Consume` 32, `change` 27,
  `rowImage` 22, `observeRelation` 35, `JSONImage` 3, `ToChange` 6,
  `ReadChanges` 25, `ReadConsumerPositions` 20, … — vollständig in
  §Eigene Messungen).
- `verifizierbar`: ja — `git diff 0f21eb2..7bf2aa5 -- …slice-088…md` gegen die
  Paket-Zugehörigkeit der Funktionen im Baum; die Zahl ist mechanisch zählbar
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`

### F-3 — §8 zählt zwei Register-Einträge zu niedrig; einer davon steht auf der Schwelle

- `kategorie`: MEDIUM
- `quelle`: Maintainability · `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
  · Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Zwei Schritte vor der Modus-Begründung (Schritt 2 *Offene Beobachtungen
  sichten*)
- `pfad`: `docs/plan/planning/in-progress/slice-088-coverage-tail-uebersetzung.md:293-311`
- `befund`: §8 sagt „Zähler am Register **nachgezählt**" und nennt
  `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` mit **„1×"** und
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` mit **„2×, offen"**. Im
  Register liegen `evidence/slice-084.md` + `slice-085.md` (2×) bzw.
  `slice-081.md` + `slice-084.md` + `slice-085.md` (**3×, Schwelle erreicht**)
  — beide `state.md` tragen den Zähler selbst als abgeleitete Größe. Die
  Sichtung liest beim ersten Eintrag die **stehengebliebene Kopfzeile** des
  `state.md` („Zustand: offen (1×)" direkt über dem eigenen Satz „Zähler
  (abgeleitet): **2×**"), beim zweiten den Stand vor dem dritten Beleg. Folge:
  die Zeile „kein Eintrag rückt mit diesem Slice über die 3×-Schwelle"
  (§8 *Ergebnis*) steht neben einem Eintrag, der die Schwelle **bereits
  erreicht** hat und dessen Ausgang der Lese-Schritt der `welle-20`-Closure
  trägt. Das ist genau der Schritt, den §8 in einem Repo mit Wellen-Betrieb
  als einzigen Leser der Unter-Schwelle-Einträge führt — er hat hier die
  Zahl nicht aus der Ablage gelesen, sondern aus dem Prosa-Stand darüber.
- `verifizierbar`: ja — `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence/`
  gegen die in §8 genannten Zähler
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`

### F-4 — Die drei Zahlen der §3-Tabelle haben keine definierte Einheit

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-088-coverage-tail-uebersetzung.md:167-176`
- `befund`: Die Spalte „Trägt" nennt „die 17 Fehlerzweige" (sqlexec), „die 9
  erreichbaren Zweige von `Decode`" und „die 17 erreichbaren Zweige von
  `Consume`/`change`/`observeRelation`/…". Keine der beiden naheliegenden
  Lesarten geht für alle drei auf: Als **Statements des jeweiligen Pakets**
  gemessen sind es +17 (`sqlexec`), +13 (`replication/decode`) und +18
  (`replication/mapper`) — wobei je eine bzw. zwei Statements der Paket-Deltas
  auf die Whitebox-Dateien entfallen (`tupleValues`/`oldTupleValues` je 1,
  `setSchemaVersion` 1); als **neue Testfunktionen je Datei** gezählt sind es
  17 (`translate_test.go`), 9 (`decode_test.go`) und 16
  (`mapper_test.go`). Für `replication/mapper` geht keine der beiden Lesarten
  auf, für `replication/decode` nur die zweite. Eine Einheit nennt der Plan
  nirgends. Die drei Zahlen tragen keine Zusage (Liefer-Punkte, §1-Tabelle,
  §2-Wirkung und §6-Befund sind davon unberührt und gehen auf); sie sind
  beschreibender Text.
- `verifizierbar`: ja — die Deltas sind aus dem Profil der `coverage`-Stufe
  mechanisch bestimmbar (Werte in §Eigene Messungen), die Testzahlen aus dem
  Diff
- `klasse`: „Zahl ohne definierte Einheit"

## Negativbefunde

- geprüft, ohne Befund: **Produktionscode** — `git diff 3de7fb1..7bf2aa5 --
  '*.go' ':!*_test.go'` ist leer (Exit 0); die vier berührten
  Produktionsdateien sind byte-identisch zum Vorgänger-`HEAD`. Die
  Folgepflicht aus `ADR-0082` §Konsequenzen („test-only, keine Zeile
  `internal/**`/`cmd/**` außerhalb von Tests") ist eingelöst.
- geprüft, ohne Befund: **der Gate-Nenner** — 1903 in beiden eigenen Läufen
  (vorher wie nachher), dedupliziert über die Block-Position; keine
  Statement-Position ist hinzugekommen.
- geprüft, ohne Befund: **die acht als unerreichbar gemeldeten Statements** —
  die ungedeckten Positionen des Profils sind genau
  `decode.go:224.3` (1), `mapper.go:158.4/172.4/191.4/235.3/239.3/537.4/541.4`
  (7). Jede der acht Angaben trägt (Begründung in §Antwort auf Schwerpunkt 3).
- geprüft, ohne Befund: **die neue `sqlexec`-Testfläche** — 17 Statements
  geschlossen, keine verbleibende Lücke; alle vier Klassen-Zusagen
  (Klasse + Ursache + Aufruf des Übersetzungspunkts) sind an der Eingabeseite
  gebunden (Mutation M2 färbt rot).
- geprüft, ohne Befund: **die zwei neuen Whitebox-Dateien** — beide
  begründet (unexportierte Details, deren nil-Zweige der öffentliche Weg
  nicht herreicht; nachgeprüft an `pglogrepl` selbst) und beide tragen, was
  sie zusagen; keine Fabrikation verbotenen Objektzustands.
- geprüft, ohne Befund: **Hygiene** — Diff enthält keine Lauf-Artefakte
  (nur sechs `_test.go` und der Slice-Plan), keine Vorlagen-Reste über die
  repo-weit üblichen hinaus, keine unaufgelösten Kennungen außerhalb der
  bewusst offenen §6/§7-Closure-Felder; Historie linear (keine Merges), keine
  Trailer, keine `SPEC-*`/`ARC-*` in den Betreffs; `make commit-traceability`
  grün über die letzten 5 Commits.
- geprüft, ohne Befund: **keine neue Betreiber-Oberfläche** — der Diff führt
  keine `CDC_*`-Variable, keine `cdc.*`-Funktion und keinen Endpunkt ein; der
  Handbuch-Zug ist nicht ausgelöst.
- geprüft, ohne Befund: **`docs/user/**`, `spec/**`, `tools/schema/**`,
  `.a-check.yml`, `Makefile`, `harness/**`** — nicht berührt.
- geprüft, ohne Befund: **die Kommentare der sechs Testdateien** — Ist-Zustand,
  Test-Provenienz als Satzsubjekt, keine Slice-/Wellen-Chronik über einem
  Produktionscode-Pfad.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Messung | Kommando | Ergebnis |
|---|---|---|
| Produktionscode | `git diff 3de7fb1..7bf2aa5 -- '*.go' ':!*_test.go'` | leer, Exit 0 |
| Gate (nachher) | `make coverage-gate` | Exit **0** · `coverage-gate: OK — Coverage 74.70% erfüllt Schwelle 70%` |
| Zählbasis (nachher) | Dedup über die Block-Position aus `/out/coverage.out` | **1422 / 1903 = 74,72 %** |
| Gate (vorher) | `docker build --target coverage` über den Teststand `3de7fb1` (Arbeitskopie) | Exit **0** · `OK — Coverage 72.00%` |
| Zählbasis (vorher) | dito | **1371 / 1903 = 72,04 %** |
| Alle Gates | `make gates` | Exit **0** (`baseline-verify` 54 Dateien OK · `docs-check` 0 Befunde · `commit-traceability` OK · `a-check` 0 Befunde · `coverage-gate` OK) |
| Unit-Tests | `make test` | Exit **0** · 31 Pakete `ok`, 0 `FAIL` |

**Die Wirkung, je Paket und je Funktion gemessen** (Dedup über die
Block-Position, beide Profile):

| Größe | vorher | nachher | Δ |
|---|---|---|---|
| Nenner | 1903 | 1903 | 0 |
| `replication/decode` gedeckt | 75 | 88 | **+13** |
| `replication/mapper` gedeckt | 166 | 184 | **+18** |
| `postgresstorage/sqlexec` gedeckt | 135 | 152 | **+17** |
| `postgresstorage/mapper` gedeckt | 16 | 20 | **+4** |
| vier Trägerpakete zusammen | 392 | 444 | **+52** |
| `domain/model` (`NewSchemaVersion`) | 4/5 | 5/5 | **+1** |
| `internal/bootstrap` (`runWALRetentionCheck`) | 16/16 | 14/16 | **−2** |
| Gesamt gedeckt | 1371 | 1422 | +51 |

Die Differenz meiner Nachher-Zahl (74,7 %) zur gemeldeten (74,8 %) ist damit
**erklärt und nicht die Zählbasis**: mein Lauf liegt auf dem unteren Ende der
in `ADR-0082` §Kontext (2) dokumentierten Schwankung (`runWALRetentionCheck`
14/16 statt 16/16, −2), mein Vorher-Lauf auf dem oberen (1371). Auf
**gleichen Enden** gerechnet ist das Delta **+53** (1369 + 53 = 1422), und die
Aufteilung des Implementers — **+52** in den vier Trägerpaketen, **+1**
mittelbar in `domain/model` — geht positionsgenau auf.

**Eigene Mutationen an der Eingabeseite** (Arbeitskopie des Baums, gepinntes
Toolchain-Image, `go test ./internal/adapters/... ./internal/adapters/driven/...`;
Baseline der Kopie vor jeder Mutation: Exit 0):

| # | Mutation (Eingabeseite) | Ergebnis |
|---|---|---|
| M1 | `tupleValues`/`oldTupleValues`: den Arity-Riegel gegen `len(tuple.Columns)` entfernt | Exit **1** — 4 neue Tests rot |
| M2 | `ReadConsumerPositions`: den gescannten Offset ignoriert (`int64(1)` statt `offset`) | Exit **1** — u. a. `TestReadConsumerPositionsLeavesDomainFailureUnclassified` rot |
| M3 | `rowImage`: den Ausschlussstand ignoriert (`containsColumn(…)` → `false`) | Exit **1** — u. a. `TestIncludeColumnKeepsRemainingExclusions` und `TestIncludeColumnOnUnboundTableStaysNoop` rot |
| M4 | `observeRelation`: `if !activated { return nil }` entfernt | Exit **0** — gesamte Fläche `./internal/... ./cmd/...` grün (**F-1**) |
| M5 | `observeRelation`: `if !found { return nil }` entfernt | Exit **0** — Paket grün (**F-1**) |
| Probe zu M4 | dieselbe Zusage an ihre Eingabe gebunden (Store zählt seine Aufrufe; Relation mit einer Spalte über der bekannten Form) | unverändert: Exit **0** · unter M4: Exit **1** (`Relation ausserhalb der Aktivierungen: leere Kennung`) |

Die fünf Mutationen des Implementers sind im Diff nicht dokumentiert; geprüft
sind hier meine eigenen (M1–M3 bestätigen die Klasse an den neuen Tests, M4/M5
widerlegen sie an zweien — F-1).

## Antwort auf die Schwerpunkte

**(1) Kein Produktionscode — trägt.** Der Diff außerhalb der Tests ist leer;
geprüft am eigenen Kommando, nicht am Bericht des Implementers.

**(2) Der Nenner bleibt 1903 — trägt.** In beiden eigenen Profilen 1903,
dedupliziert über die Block-Position; das Delta von +53 geht auf gleichen
Enden der Schwankung positionsgenau auf und stammt **nicht** aus ihr — die
Schwankung liegt in `runWALRetentionCheck` und ist in meiner Rechnung separat
ausgewiesen.

**(3) „8 Statements sind strukturell unerreichbar" — trägt, jede der acht
Angaben einzeln.** Die ungedeckten Positionen des Profils sind genau die
gemeldeten acht; und für jede ist die Unerreichbarkeit nachgeprüft:
`decode.go:224` (der `default:`-Zweig) — `pglogrepl.Parse` kann nur zehn
Typen liefern (`Relation`, `Type`, `Insert`, `Update`, `Delete`, `Truncate`,
`LogicalDecoding` sowie über `getCommonDecoder` `Begin`, `Commit`, `Origin`)
und gibt für jeden anderen Nachrichtentyp einen Fehler statt einer Nachricht
zurück; die Fallunterscheidung führt genau diese zehn, der Rumpf ist damit
nicht auslösbar. `mapper.go:158/172/191` — `NewOpenTransaction` scheitert nur
an leerer Kennung, und die Kennung entsteht aus `strconv.FormatUint` (nie
leer), die Quelle ist im Konstruktor auf nicht-leer geprüft; `Commit`
scheitert nur an `committed` (nach dem Commit setzt `Consume` `a.open` auf
`nil`, eine zweite Commit-Nachricht läuft in `ErrCommitWithoutBegin`) oder an
einer fremden Quelle (die Position entsteht aus `a.source`, die Transaktion
ebenso); `AppendChange` scheitert nur an `committed` (dito), an fremder
Transaktions-Kennung (der Change trägt `a.open.tx.ID`) oder an doppelter
Sequenz (`a.open.sequence` zählt monoton). `mapper.go:235/239/537/541` —
`json.Marshal` eines `string` endet nie im Fehler; nachgeprüft mit einem
Programm gegen die gepinnte Toolchain (leerer String, Quoting, ungültiges
UTF-8, ` `, Länge): in allen sechs Fällen `err=<nil>`. **Die Aussage, die
acht seien für jeden Test unerreichbar, trägt.** Die Unterscheidung, die der
Implementer zieht, halte ich für konsistent und nicht für eine Bequemlichkeit:
Wo Whitebox-Zugriff legitim war (`tupleValues`/`oldTupleValues` mit `nil`;
`setSchemaVersion` auf einer nicht getragenen Bindung — beide Male ist der
Aufruf mit einem regulär konstruierten Objekt darstellbar, `RemoveBinding`
erzeugt den Zustand selbst), hat er ihn genommen; wo er den Objektzustand
hätte **fabrizieren** müssen, den der Konstruktor verbietet (`Assembler` mit
leerer Quelle; Transaktion in `committed`-Zustand), hat er es gelassen und den
Rest als Befund geführt. Der Befund ist außerdem nicht „grün gefärbt":
die acht Positionen bleiben im Profil ungedeckt.

**(4) Kein Test grün ohne Aussage — trägt nicht vollständig.** Drei eigene
Mutationen an der Eingabeseite färben die neuen Tests rot (M1–M3); **zwei**
Mutationen, die genau die Riegel zweier neuer Zusagen entfernen, bleiben
grün (M4/M5) — F-1. Die Klasse ist damit nicht getilgt, sondern an zwei
Stellen dieses Slice noch vorhanden.

**(5) Die §1-Berichtigung — die Tabelle trägt, die Zahl daneben nicht.**
Paket-Zuordnung und Summen sind korrekt (14/25/17/4 = 60, gegen den Baum
geprüft; alle 22 Funktions-Nenner stimmen mit dem Profil überein); die
Gesamt-Zahl „fünf" ist falsch (acht, F-2) und die Aufzählung nennt `ToChange`,
das nicht betroffen ist.

**(6) Hygiene — ohne Befund.** Keine Lauf-Artefakte, keine Trailer, lineare
Historie, Betreffs ohne Struktur-IDs; die zwei Whitebox-Dateien sind
begründet und tragen ihre Zusage (`pglogrepl` liefert für `Insert`/`Update`
nie ein `nil`-Tupel und für `Delete` nur `K`/`O`; der nil-Zweig ist über
`Decode` nicht erreichbar, der `setSchemaVersion`-Vertrag für eine nicht
getragene Bindung ebenso wenig).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (F-1) ·
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (F-2, F-3) ·
„Zahl ohne definierte Einheit" (F-4, INFO).

## Verdikt

**Merge-blockierend:** ja — 0 HIGH, 3 MEDIUM. F-1 liegt auf der
Implementer-Seite: zwei der 16 neuen Tests behaupten eine Abwesenheit, die
ihren Gegenstand nicht bindet; sie zu binden ist ein zweiter Implementer-Lauf
(dieselbe Datei, dieselbe Zusage), und der Slice benennt die Klasse in §6
selbst als tragend. F-2 und F-3 sind Plan-Defekte (Planner-Kante): eine
falsche Zahl in der Berichtigung selbst und ein nicht aus der Ablage gelesener
Register-Zähler — beide mechanisch prüfbar, beide ohne Wirkung auf Lieferung
oder Gate. Die vier Berichtigungen des Messgegenstands, die vier Liefer-Punkte
und die bezifferte Wirkung selbst sind **ohne Befund**: kein Produktionscode,
Nenner 1903 stabil, +53 auf gleichen Enden, `make gates` und `make test` grün.

**Übergabe:** **Rückgabe-Pfeil Reviewer → Implementer ist nötig** — für F-1
(Testbindung in `mapper_test.go`); F-2/F-3 gehen als **Rückkante
Review → Plan** an den Planner. Die Finding-Klassen gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein
**Lauf-Beleg** und ersetzt keine Verifikation (Modul 11): DoD-/Spec-Konformität
prüft der Verifier separat, §6-Risiko-Ausgänge und die drei Paarungen die
Planner-Closure.

**DoD-Häkchen „Review durchgeführt":** bleibt **offen** — der Slice braucht
eine Fixrunde, und der Nachzug läuft dann regulär über Schritt 21 des
Implementer-Workflows (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug
ohne Fixrunde, Grenze). Der Plan ist in diesem Commit **nicht** angefasst.
