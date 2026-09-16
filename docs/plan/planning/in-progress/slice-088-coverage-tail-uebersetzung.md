# Slice slice-088: Coverage-Tail „Reine Übersetzung" — Decode, Mapper, sqlexec

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-20`](../welle-20.md) („die Coverage erreicht 80 % über der
netzlos prüfbaren Fläche") — **Cluster B** ihres Schnittmaßes.

**Bezug:** [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(das Schnittmaß: der Tail trägt die Endstufe, nicht der Composition Root) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand) · [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a) (die Rampe) · [`ADR-0039`](../../adr/README.md) (Mapper übersetzen,
Schichten-Ordnung).

**Berührte Spec-Stellen:** — (der Slice fügt **Tests** hinzu; er ändert keinen
Vertrag).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-16.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die **reine Übersetzung** des Repos netzlos decken — die 60
ungedeckten Statements aus Cluster B des Schnittmaßes ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)),
gemessen über die vier Trägerpakete. Es sind **21 Funktionen mit 1 bis 8
ungedeckten Statements** — **flach und breit, nicht tief**; kein Umbau, sondern
Tests an Stellen, die noch keinen haben.

**Die Gegenstände, gemessen** (ungedeckt / gesamt):

| Paket | Funktionen |
|---|---|
| `replication/decode` | `Decode` 8/50 · `observeRelation` 8/35 |
| `replication/mapper` | `Consume` 6/32 · `change` 4/27 · `JSONImage` **3/3** · `oldTupleValues` 3/17 · `tupleValues` 3/11 · `rowImage` 2/22 · `classifyRelationColumns` 1/11 |
| `postgresstorage/sqlexec` | `ReadConsumerPositions` 4/20 · `ReadExcludedColumns` 3/22 · `ReadTableSchema` 3/19 · `ReadSourceTables` 3/19 · `ReadPendingRequests` 2/19 · `ReadChanges` 1/25 · `ReadConsumerPosition` 1/10 · `IncludeColumn` 1/7 · `setSchemaVersion` 1/7 · `removeExcluded` 1/5 |
| `postgresstorage/mapper` | `ToChange` 1/6 · `qualifiedNames` 1/6 |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Umbau der Übersetzung.** Der Slice fügt **Tests** hinzu; ändert sich
  eine Funktion als nicht-prüfbar, ist das ein **Befund** (und die Naht aus
  `ADR-0080` ihr Träger), kein Auftrag.
- **Der Composition Root.** `bootstrap.Run` ist **netzlos nicht prüfbar** — 0
  von 198 Statements im Gate-Profil (`ADR-0082`); er ist **nicht** Gegenstand
  dieser Welle, und der frühere Vorschlag, dort „den Hebel" zu suchen, ist
  widerlegt.
- **Die anderen drei Cluster** (C, D, A) — je ein eigener Slice, geschnitten
  **nach** diesem (Modul 5: Plan und Implementation alternieren).
- **Der Messmechanismus.** Gegenstand, Endstufe und Rampe sind unberührt
  (`AGENTS.md` §3.6); dieser Slice **hebt die Zahl**, er ändert nicht, was sie
  misst. Die code-granulare Neufassung des Gegenstands ist als Trigger (c) in
  `ADR-0082` vertagt.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — die Stream-Übersetzung ist gedeckt.**

- [x] `replication/decode` und `replication/mapper` haben Tests für die
      gemessenen Lücken (`Decode`, `observeRelation`, `Consume`, `change`,
      `JSONImage`, `oldTupleValues`, `tupleValues`, `rowImage`,
      `classifyRelationColumns`) — **netzlos**, ohne externe Dienste.
- [x] `JSONImage` (heute **3/3 ungedeckt**) ist dabei der erste Fall: eine
      Funktion ohne jede Abdeckung.

**Liefer-Punkt 2 — die SQL-Übersetzung ist gedeckt.**

- [x] `postgresstorage/sqlexec` hat Tests für die gemessenen Lücken
      (`ReadChanges`, `ReadConsumerPosition(s)`, `ReadSourceTables`,
      `ReadExcludedColumns`, `ReadTableSchema`, `ReadPendingRequests`,
      `IncludeColumn`, `removeExcluded`, `setSchemaVersion`) — über den
      **bestehenden** Fake der Naht aus [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md).

**Liefer-Punkt 3 — die Wirkung ist gemessen, nicht angestrebt.**

- [x] `postgresstorage/mapper`s Lücken (`ToChange`, `qualifiedNames`) sind
      gedeckt.
- [x] **Der Effekt ist beziffert:** die Gate-Zahl vorher/nachher, gemessen über
      `make coverage-gate` — **erwartet ≈60 Statements**, und der Nenner bleibt
      **1903** (dieser Slice fügt **keinen** Produktionscode hinzu; wächst er,
      ist etwas anderes passiert und gehört in den Bericht).
- [x] `make gates` grün (Exit direkt, ungepiped).

**Die bezifferte Wirkung, gemessen** (`make coverage-gate`, Exit 0 in beiden
Läufen): **72,0 % → 74,8 %**, Nenner **1903** unverändert; gedeckt
**1371 → 1424** (+53). Davon **+52** in den vier Trägerpaketen (die 60
gemessenen Statements dieses Slice, **8** davon sind über die öffentliche
Fläche **nicht erreichbar** — siehe §6) und **+1** mittelbar in
`domain/model` (`NewSchemaVersion`); der frühere Lauf desselben Stands druckte
71,9 % — die von [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
§Kontext (2) dokumentierte ±2-Schwankung von `runWALRetentionCheck`.

- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/decode/*_test.go` · `.../mapper/*_test.go` | neu | Tests für die gemessenen Lücken; beide Pakete sind netzlos prüfbar |
| `internal/adapters/driven/postgresstorage/sqlexec/*_test.go` | update/neu | Tests über den **bestehenden** Fake der Naht (`slice-081`) |
| `internal/adapters/driven/postgresstorage/mapper/*_test.go` | neu | `ToChange`, `qualifiedNames` |

**Der genaue Datei-Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste
nennt die Träger. **Produktionscode wird nicht geändert**; erweist sich eine
Funktion als netzlos nicht prüfbar, ist das ein **Befund** und gehört in den
Bericht (§6), nicht in einen Umbau.

**Der Zuschnitt dieses Laufs** — sechs Dateien, alle **Tests**:

| Datei | Änderungs-Art | Trägt |
|---|---|---|
| `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` | update | die 17 Fehlerzweige der Zeilen-Übersetzung |
| `internal/adapters/driven/postgresstorage/mapper/mapper_test.go` | update | `JSONImage`, `ToChange` |
| `internal/adapters/driving/replication/decode/decode_test.go` | update | die 9 erreichbaren Zweige von `Decode` |
| `internal/adapters/driving/replication/decode/tuplevalues_internal_test.go` | neu (Whitebox) | die Abwesenheits-Grenze von `tupleValues`/`oldTupleValues` |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | update | die 17 erreichbaren Zweige von `Consume`/`change`/`observeRelation`/`IncludeColumn`/`removeExcluded`/`qualifiedNames` |
| `internal/adapters/driving/replication/mapper/schemaversion_internal_test.go` | neu (Whitebox) | der idempotente Vertrag von `setSchemaVersion` |

**Nicht in dieser Liste:** `internal/bootstrap/**`, `cmd/**`, `internal/driving/http/**`,
`.../grpc/**`, `.../natsnotify/**` (Cluster A, C, D — eigene Slices); der
Messmechanismus (`Dockerfile`-Filter, `coverage.mk`); `spec/**`.

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
liegt `Accepted` vor — **erfüllt** —, `welle-20` ist offen, `Verantwortlich:`
gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich der Tail
  als **tiefer** als gemessen (ein Paket trägt viel mehr, als die genannten
  Funktionen vermuten lassen), gehört er zurück zur Zerlegung — **nach Paket**,
  nicht als Sammel-Zug.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass eine der
  Funktionen **netzlos nicht prüfbar** ist (sie braucht einen lebenden Dienst),
  ist das ein Blocker mit Entscheidung — `ADR-0082` hat für `Run` genau diesen
  Fall entschieden; ein **zweiter** Fall gehört als eigener Befund geführt, nicht
  stillschweigend übergangen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make coverage-gate` real grün mit **höherer** Zahl
**und** der **Nenner unverändert** (1903 — dieser Slice fügt keinen
Produktionscode hinzu) **und** `make gates` grün **und** die Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Eine Funktion könnte netzlos nicht prüfbar sein.** Dann trägt sie zum
  Gegenstand bei, ohne für dieses Ziel erreichbar zu sein — [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  hat das für `Run` gemessen. — **Ausgang:** <bei Closure>
  **Beobachtet in diesem Lauf** (gemessen am Profil der `coverage`-Stufe,
  dedupliziert über die Block-Position): **keine** der 21 Funktionen ist
  unprüfbar — jede trägt Tests. **8 der 60 Statements** sind aber über die
  öffentliche Fläche nicht **erreichbar**, bleiben also ungedeckt:
  `decode` `Decode` 224 (der `default:`-Zweig — `pglogrepl.Parse` liefert nur
  Nachrichtentypen, die die Fallunterscheidung führt), `mapper` `Consume`
  158/172/191 (`NewOpenTransaction`/`Commit`/`AppendChange` — ihre Grenzen
  sind an einem Assembler, den `NewAssembler` erzeugt, nicht herstellbar),
  `mapper` `change` 235/239 und `rowImage` 537/541
  (`json.Marshal` eines `string` endet nie im Fehler).** Der Zuschnitt bleibt:
  kein Umbau (Modul 9, §1 dieser Datei); die Zahl der gedeckten Statements
  dieser vier Pakete ist damit **52**, nicht 60.
- **Der Slice könnte den Produktionscode anfassen.** Die Zusage ist „Tests,
  kein Umbau"; ein Umbau wäre ein **anderer Vorgang**. — **Ausgang:** <bei Closure>
- **Die Tests könnten die Zahl heben, ohne etwas zu prüfen.** Ein Test, der eine
  Funktion aufruft und nichts behauptet, ist grün ohne Aussage — die Klasse
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**3×**). — **Ausgang:** <bei Closure>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt die vier Trägerpakete und die
Testpyramide in **einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen; Zähler am Register **nachgezählt**:

- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**3×**, Schwelle erreicht
  — Ausgang beim Lese-Schritt der `welle-20`-Closure): **Treffer, und ein
  genauer.** Dieser Slice besteht **aus** Zusagen an ihre Eingabeseite: er fügt
  Tests hinzu, deren einziger Zweck es ist, eine Aussage an einen Eingabewert zu
  binden. Die Klasse ist damit **Auftrag** dieses Slice, nicht nur Beobachtung —
  *ein Test, der eine Funktion aufruft und nichts behauptet, ist grün ohne
  Aussage.*
- `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (1×, offen): **kein**
  Treffer — sie betrifft den **DB**-Gegenstand, dieser Slice den **Unit**-Gegenstand.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (2×, offen):
  **benachbart** — die Zahl dieses Slice (≈60) ist eine **Erwartung** aus einer
  Messung; die Closure führt den gemessenen Ist-Wert daneben, nicht an ihrer
  Stelle.

**Ergebnis** (Stand: Anlage dieses Plans): kein Eintrag rückt mit diesem Slice
über die 3×-Schwelle.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Testpyramide ist über `ADR-0030` verankert, die
Schichten-Ordnung über `ADR-0039`, der Messgegenstand über `ADR-0071`),
**Phase-Reife** hoch, **Evidenz-/Diskrepanz-Risiko** **niedrig** — die Lücken
sind **gemessen**, die Träger sind netzlos —, **Reconciliation-Aufwand** null.
