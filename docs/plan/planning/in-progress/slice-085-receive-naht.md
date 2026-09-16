# Slice slice-085: receive-Naht — schmale Abhängigkeit statt konkreter pgconn-Verbindung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (eine
Struktur-Änderung mit eigenem Beleg), kein repo-weites *Mehr*.

**Bezug:** [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
(die **Nahtform**: Treiber-Hülle plus schmale, fake-fähige Fläche, **im Paket** —
sie beantwortet den Verdacht aus §6) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 (schmale Abhängigkeit statt konkreter Typ) · [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(der Nachweis, den diese Naht als **Null-Befund** führt) · `slice-084` (die
Geschwister-Naht in `postgresack`) ·
[`ADR-0006`](../../adr/0006-replication-stream-driving-adapter.md) (der
Replication-Stream als **driving** Adapter — die Schicht, in der diese Naht
liegt) · [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
(die Fehlerklassen und WAL-Rückstand-Schwellen dieses Empfangs).

**Berührte Spec-Stellen:** — (die Naht liegt **innerhalb** des driving Adapters,
ohne Vertrag oder Sicht zu berühren).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

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

**Ziel:** `internal/adapters/driving/replication/receive` hängt an einer
**konkreten** `*pgconn.PgConn` — gemessene Fläche (**18** Symbole):
`pgconn.ConnectConfig`, `pgconn.ErrorResponseToPgError`, `pgconn.ParseConfig`,
`pgconn.PgConn` sowie **14** `pglogrepl`-Symbole (`StartReplication`,
`CreateReplicationSlot`, `IdentifySystem`, `ParseXLogData`,
`ParsePrimaryKeepaliveMessage`, `SendStandbyStatusUpdate`, …). Die Naht macht
ihre Verklebung netzlos prüfbar.

**Die Nahtform ist entschieden** ([`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)):
dieselbe Form wie in `slice-084`, **im Paket** — eine **Treiber-Hülle**
(`connSession{conn *pgconn.PgConn}`) plus eine **schmale, fake-fähige Fläche**
mit **sieben** Operationen (`IdentifySystem`, `CreateReplicationSlot`,
`StartReplication`, `SendStandbyStatusUpdate`, `Exec → []*pgconn.Result`,
`ReceiveMessage → pgproto3.BackendMessage`, `Close`). Der öffentliche Rand
(`NewStream`, `Stream`, `WALRetentionChecker`) bleibt **unverändert**.

**Der Verdacht dieses Plans ist gemessen — bestätigt im Buchstaben, widerlegt in
seiner Folge.** §6 führte, `pglogrepl.StartReplication` und
`CreateReplicationSlot` verlangten den konkreten Typ, die Naht trage deshalb
nicht: das stimmt **wörtlich** (vier Paketfunktionen tragen `conn *pgconn.PgConn`
in der Signatur) — die **Folge** tritt aber **nicht** ein, weil die Hülle den
konkreten Typ **innen** hält und die Fläche nach außen fake-fähig macht. Der
§4-Blockerfall ist damit **beantwortet**. **`slice-081`s Zusicherungsform trägt
hier nicht** — der Typ erfüllt die Schnittstelle nicht selbst.
`replication/decode` und `.../mapper` bleiben **draußen**.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Design-Vorgriff.** Die Naht-Form ist eine Architect-Entscheidung.
- **`postgresack`** — eigene Naht, eigene Datei: `slice-084`.
- **Verhaltensänderungen** — die reale Verdrahtung muss unverändert durch
  `make test-replication` gehen; die realen Tests sind der Wächter.
- **Ein Fake als Ersatz der realen DB-Tests** — er prüft die Verklebung, nicht
  das Protokoll.
- **Eine neue Abhängigkeit** — `pgconn`/`pglogrepl` bleiben die realen Träger.

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

**Liefer-Punkt 1 — die Naht existiert.**

- [x] Die **Logik** des Adapters hängt an einer **adapter-eigenen** Schnittstelle
      (sieben Operationen); der **konkrete** `*pgconn.PgConn` bleibt in der Hülle
      und im Dial. **Nicht** „statt der konkreten Verbindung": der Typ bleibt,
      aber **hinter** der Naht (`ADR-0080`). (`driverSession`/`connSession`,
      `seam.go`; `Stream.session`, `WALRetentionChecker.session`.)
- [x] **Kein Verhaltens-Change:** die reale Verdrahtung geht unverändert durch
      `make test-replication` (Exit 0).

**Liefer-Punkt 2 — die reine Logik ist prüfbar.**

- [x] Der netzlos prüfbare Teil (Meldungs-Zerlegung, LSN-Form,
      Keepalive-Behandlung) liegt als **reine Funktion** mit eigenen Tests vor —
      **soweit er nicht schon in `replication/decode`/`mapper` liegt**; was
      dorthin gehört, entscheidet der Architect (§1). (`parseXLogData`,
      `parseKeepalive`, `parseLSN`, `firstRow`, `standbyStatus`, `slotLSNQuery`
      — je mit eigenem Test in `seam_test.go`.)
- [x] Die Fake-Seite fährt die **Verklebung** — und ist ausdrücklich **kein**
      Ersatz der realen Tests.

**Liefer-Punkt 3 — der Nachweis ist geführt, und er ist ein Null-Befund.**

- [x] Der Nachweis aus [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
      wird als **Null-Befund** geführt (`ADR-0080`): die Naht bleibt **im**
      Paket, **kein Subjekt-Transfer** — `k_ab = 0`, der abfließende Nenner
      (155 Statements) unberührt, **keine Neu-Bemessung**, Rampen unverändert.
      Nachzuweisen sind dennoch alle drei Teile: **(a)** `k_ab = 0`; **(b)** der
      **Paket-Diff** zeigt **keinen** Trägerwechsel; **(c)** kein Verhalten
      verloren. **Benannte Grenze:** die neuen netzlosen Tests verdünnen den
      DB-Nenner leicht — mit Trigger, nicht still. (Alle drei Teile mit den
      Zahlen dieses Laufs in §3.)
- [x] `make gates` grün (Exit direkt, ungepiped).

- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] **Falls dieser Zug die Rampe bewegt:** der Transfer-Nachweis ist in
      `harness/sensors/db-adapter-coverage.md` bzw.
      `harness/sensors/coverage-gate.md` nachgezogen — **ohne** neue
      Schwellen-ADR (`ADR-0078`). *(Kein Transfer — die Naht bleibt im Paket;
      die Bedingung ist nicht eingetreten, siehe §3.)*
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(**Entfällt**: dieses Repo führt die Datei nicht — Greenfield-Bootstrap, kein Inventur-Fund.)*
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

**Zuschnitt des Implementer-Laufs** (die Liste nennt die Träger; der genaue
Schnitt entsteht hier):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/receive/seam.go` | neu | die Naht (`ADR-0080`): `driverSession` (sieben Operationen), die Treiber-Hülle `connSession` über `*pgconn.PgConn`, die Kompilier-Zusicherung `var _ driverSession = connSession{}` |
| `internal/adapters/driving/replication/receive/receive.go` | refactor | die Logik hängt an der Naht (`Stream.session`); reine Funktionen `firstRow` (Katalog-Zeilen-Übersetzung), `parseLSN` (LSN-Form), `standbyStatus` (Standby-Form), `parseXLogData` und `parseKeepalive` (Meldungs-Zerlegung) sowie `slotLSNQuery`; die Verklebung `querySingle`, `ensureSlot`, `ensurePublication` und `handleCopyData` fährt der Fake; `newStreamOnSession` ist der **paket-interne** Einstieg; `NewStream`, `Stream.Conn/Run/Assembler/BindCapture` unverändert |
| `internal/adapters/driving/replication/receive/walretention.go` | refactor | `WALRetentionChecker.session` statt `conn`; `Measure`, `Close` und `reconnectAfterError` laufen über die Naht; `newWALRetentionCheckerOnSession` als **paket-interner** Einstieg; `NewWALRetentionChecker` unverändert |
| `internal/adapters/driving/replication/receive/seam_test.go` | neu | der Fake und die netzlosen Tests: Meldungs-Zerlegung, Empfangs-Schleife, Slot-/Publication-Auflösung, Katalog-Zeilen-Übersetzung, Rückstands-Messung — der Fake erfüllt **dieselbe** Schnittstelle wie die Hülle |
| `harness/image-hash.txt` | update | der Zug ändert Build-Kontext-Dateien; `make image` stempelt den Digest des Laufs (`ADR-0044`) — `sha256:4bd43435…` → `sha256:84bdca56…` |

**Kein Paketwechsel, kein Unterpaket** (`ADR-0080`) — deshalb auch **keine**
Änderung am `Dockerfile`-Filter oder an `DB_COVERAGE_PKGS`: der DB-Gegenstand
bleibt, wie er ist.

**Nicht angefasst:** `internal/bootstrap/wiring.go` (die Composition Root ruft
`receive.NewStream`/`stream.Conn()`/`receive.NewWALRetentionChecker` unverändert),
`Dockerfile` (Stufe `coverage`, Paket-Filter), `tools/harness/db-coverage.sh`
(`DB_COVERAGE_PKGS`), `harness/mk/coverage.mk` (`THRESHOLD`),
`harness/sensors/**`, `.a-check.yml`, `spec/**`. `stream_test.go` (die realen
Tests) ist **unverändert** — kein Testfall entfernt.

**Nicht in dieser Liste:** `internal/adapters/driven/postgresack/**`
(→ `slice-084`); `internal/adapters/driving/replication/decode/**` und
`.../mapper/**` bleiben **außerhalb** (entschieden, `ADR-0080`); `spec/**`,
`.a-check.yml`.

**Der Null-Befund — alle drei Teile, mit den Zahlen dieses Laufs**
([`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
§Entscheidung 4, [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
§Entscheidung 1). Gemessen in den gepinnten Images, Exit-Codes ungepiped; die
Zahlen je Gegenstand aus den beiden `-coverprofile` des Laufs nach der
Dedup-Regel aus `tools/harness/db-coverage.sh`:

| Gegenstand (Statements) | vor dem Zug | nach dem Zug | Δ |
|---|---|---|---|
| netzlos prüfbare Fläche (Unit) | 1903 (1369 gedeckt, 71,90 %) | 1903 (1369 gedeckt, 71,90 %) | **0** |
| DB-Adapter-Gegenstand | 659 (491 gedeckt, 74,51 %) | **691** (532 gedeckt, 76,99 %) | **+32** |
| davon `postgresstorage` | 472 | 472 | 0 |
| davon `postgresack` | 32 | 32 | 0 |
| davon `replication/receive` | 155 | **187** | **+32** |

**(a) `k_ab = 0`.** Kein Gegenstand gibt Statements ab: der Unit-Nenner steht
unverändert bei 1903 (die Naht liegt **im** ausgenommenen Paket — der Lauf
deckt dieselben 1369 Statements), und innerhalb des DB-Gegenstands sind
`postgresstorage` (472) und `postgresack` (32) **byte-stabil**. Der DB-Nenner
**wächst** um **+32** — die neuen Statements der sieben Hüllen-Methoden, der
reinen Funktionen und des paket-internen Einstiegs; `k_ab = 0` ist damit die
Null-Hälfte der Arithmetik, nicht eine Behauptung.

**(b) Paket-Diff — kein Trägerwechsel.** Die Änderung ist auf **ein** Paket
isoliert, und der Beleg braucht **beides — Range *und* Pathspec**:
`git diff --name-only 95480a5..610751c -- internal/adapters/driving/replication/receive/`
listet **vier** geänderte Dateien (`seam.go`, `seam_test.go`, `receive.go`,
`walretention.go`); und die **Gegenrichtung** trägt den Rest:
`git diff --name-only 95480a5..610751c -- internal/ ':!internal/adapters/driving/replication/receive/'`
ist **leer**. Die Gegenstands-Listen sind unberührt
(`git diff 95480a5..610751c -- Dockerfile tools/harness/db-coverage.sh harness/mk/coverage.mk`
ist leer), ebenso `.a-check.yml` und die Composition Root. Kein Paket wechselt
zwischen den zwei Gegenständen. **Beides gehört in den Beleg**
(`review-slice-084` F-2, Verifikation desselben Slice **V-1**): ein `git diff`
**ohne Range** ist auf sauberem Baum leer, und eines **ohne Pathspec** listet
mehr als das Behauptete.

**(c) Kein Verhalten verloren.** Die realen, dienst-gestützten Läufe sind grün
(`make test-replication` Exit 0, `make test` Exit 0 — je ungepiped), und **kein
Testfall wird entfernt**: `stream_test.go` ist unverändert,
`git diff --name-status 95480a5..610751c -- '*_test.go'` zeigt genau **eine
neue** Datei (`seam_test.go`).

**Rot-Gegenprobe an der Zusage — einmal gesehen.** Auf einer Wegwerf-Kopie
brachen **vier** Mutationen je ihre Prüfung: die Slot-Auflösung ohne den
bestehenden Slot (`if exists && false`) färbt
`TestEnsureSlotResumesExistingSlot` rot; die Katalog-Zeilen-Übersetzung mit
einem verfälschten Wert färbt `TestFirstRowTranslatesCatalogRow` rot; die
Standby-Form ohne den Apply-Stand färbt
`TestStandbyStatusCarriesPositionInAllThreeLSNs` **und**
`TestRunAnswersKeepaliveWithAcknowledgedPosition` rot; die neutralisierte
Meldungs-Zerlegung färbt die drei Keepalive-/XLogData-Tests rot. Und **die
Hülle selbst fängt nicht der Fake, sondern der reale Tier**: eine
nicht-delegierende `connSession.Exec` bleibt in den neuen netzlosen Tests
**grün** und färbt `make test-replication` rot (sieben reale Fälle,
Exit 2) — die Wächter-Rolle liegt genau dort, wo §1 sie verortet.

**Benannte Grenze:** die neuen netzlosen Tests laufen im Messlauf der
DB-Adapter-Coverage mit und decken dort Statements, die keine
PostgreSQL-Instanz berührt hat — `replication/receive` steht deshalb bei
**153 von 187** gedeckten Statements (vorher 112/155), und `stream_test.go`
ist byte-identisch: der Zuwachs stammt ausschließlich aus `seam_test.go`. Die
Verdünnung ist damit keine Randgröße: **ohne** jede PostgreSQL-Verbindung
gedeckt sind in `replication/receive` **112 von 187 (59,89 %)** statt vorher
**11 von 155 (7,10 %)**; im Gegenstand insgesamt steigt der netzlos gedeckte
Anteil der gedeckten Statements von **41 von 491 (8,35 %)** auf **142 von 532
(26,69 %)** — je ein `go test -coverpkg` ohne `CDC_REPLICATION_TEST_DSN` an
beiden Enden der Range, Block-Position-dedupliziert. Das ist die Größe des
Triggers aus
[`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
(§„Die benannte Grenze") — die Zahl, die seinen Buchstaben entscheidet, ist
der Anteil, nicht der Nenner. Die Zahl heißt weiterhin
richtig „Coverage des DB-Gegenstands". **Keine Rampe bewegt:**
`DB_COVERAGE_THRESHOLD` bleibt 70, `THRESHOLD` bleibt 70, die Endstufen bleiben
80 % — eine Schwellen-ADR wird nicht fällig.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): das Naht-Design ist architect-entschieden,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich die Naht
  als **Paket-/API-Umzug** — muss also etwas **aus** dem Paket heraus oder die
  öffentliche API sich ändern —, gehört sie zurück zur Zerlegung, nicht in einen
  wachsenden Sammelumbau. **`decode`/`mapper` sind entschieden draußen**
  (`ADR-0080`); dieser Anlass ist verbraucht, ein neuer müsste benannt werden.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass die Naht ohne
  **Verhaltensänderung** nicht zu haben ist (die `pglogrepl`-Aufrufe verlangen
  den konkreten `*pgconn.PgConn`), ist das ein Blocker mit Entscheidung — genau
  der Fall, an dem `slice-081`s Abweichung diesen Vorgang ausgelöst hat.
  **Vorab beantwortet:** `ADR-0080` hat ihn gemessen; die Hülle löst ihn.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** die reale Verdrahtung unverändert grün **und** der
Transfer-Nachweis geführt **und** `make gates` grün **und** die Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Naht könnte hier nicht tragen.** `pglogrepl.StartReplication` und
  `CreateReplicationSlot` verlangen den **konkreten** `*pgconn.PgConn` — eine
  schmale Schnittstelle drückt das womöglich nicht aus. — **Ausgang:**
  *entfallen — gestrichen mit Begründung*, **noch vor dem Start**: `ADR-0080`
  hat den Verdacht gemessen. Er stimmt im **Buchstaben** (vier Paketfunktionen
  tragen den konkreten Typ in der Signatur), tritt in seiner **Folge** aber
  nicht ein — die Hülle hält den Typ **innen**, die Fläche ist nach außen
  fake-fähig. Der §4-Blockerfall ist beantwortet, bevor er eintreten konnte.
- **Der Umbau könnte ein Verhaltens-Change sein, der als Refactoring auftritt.**
  Wächter sind die **realen** Tests, nicht die neuen Fakes. — **Ausgang:** <bei Closure>
- **Der Fake könnte grün sein, ohne etwas zu prüfen.** — **Ausgang:** <bei Closure>
- **Die Naht könnte die Messung bewegen** — wie bei `slice-081`. Dann gilt
  `ADR-0078`s Transfer-Nachweis; eine Schwellen-Frage entsteht daraus **nicht**
  erneut. — **Ausgang:** <bei Closure>

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `internal/adapters/**` in
**einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen; die Zählerstände sind am Register
**nachgezählt** (`ls evidence/`), nicht aus diesem Text übernommen:

- `BEO-PGC/adapter-fehler-ausgang` (2×, offen): **benachbart, kein Treffer** —
  sie führt die Adapter-Grenze „kontrollierte Fortsetzung beim Aufrufer" und
  einen Retry-/Backoff-Aufschub; dieser Slice ändert, **woran** der Adapter
  hängt, nicht was der Aufrufer tut.
- `BEO-PGC/spec008-replication-luecke` (2×, offen): **kein Treffer** — sie
  betrifft eine Spec-Abdeckungslücke des Streams, nicht seine Verklebung.

**Ergebnis** (Stand: Anlage dieses Plans — kein Versprechen über den Lauf):
kein Eintrag steht über der Schwelle, und keiner rückt mit diesem Slice auf 3×.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Schichten-Edges prüft `.a-check.yml`, die
Ports-und-Adapter-Ordnung führt `spec/architecture.md`; `ADR-0006` und
`ADR-0049` verankern diesen Adapter), **Phase-Reife** hoch,
**Evidenz-/Diskrepanz-Risiko** **niedrig** — der Verdacht aus §6 ist mit
[`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) gemessen
und in seiner Folge widerlegt; was bleibt, ist ein Umbau innerhalb **eines**
Pakets —,
**Reconciliation-Aufwand** null.
