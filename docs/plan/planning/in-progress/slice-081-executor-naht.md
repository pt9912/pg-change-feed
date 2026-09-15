# Slice slice-081: Executor-Naht — schmale Abhängigkeit statt konkreter Pool

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (eine
Struktur-Änderung mit eigenem Beleg), kein repo-weites *Mehr*. Ausdrücklich
**nicht** Teil von `welle-20`: deren §6 schließt sie als eigenen Vorgang aus.

**Bezug:** [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 (die Naht ist zulässig und **design**-begründet, nicht zahlen-getrieben);
`spec/architecture.md` (die Schichten- und Ports-Ordnung, die den Adapter als
austauschbare Schicht führt) · [`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md)
(die Schichten-Edges, die `.a-check.yml` prüft und die dieser Slice **nicht**
ändert: die Naht liegt innerhalb des driven Adapters) ·
[`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(die Neu-Bemessung der Rampe bei einem **Subjekt-Transfer** — er ist die
Entscheidung, die dieser Slice ausgelöst hat, und ihre Umsetzung gehört in ihn,
siehe §1 und §3) ·
[`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(der **Transfer-Nachweis** anstelle der Summen-Konstanz — Teil-Supersede von
`ADR-0077`, weil dessen tragende Bedingung gemessen nicht trug).

**Berührte Spec-Stellen:** — (die Sicht beschreibt Schichten und Ports; dieser
Slice ändert **innerhalb** des driven Adapters, ohne Vertrag oder Sicht zu
berühren).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die DB-Adapter hängen an einem **konkreten** `*pgxpool.Pool`. Eine
adapter-eigene Executor-Schnittstelle (`Query`/`Exec` plus `QueryRow`) macht
ihre Verklebung netzlos prüfbar: Fehlerklassifikation und Zeilen-Übersetzung
werden **reine Funktionen** mit eigenen Tests.

**Berichtigung des Vertrags** (Review `review-slice-081` F-2). Diese Sektion
sagte zuerst „plus ein **minimales** `Rows`" **und** „der reale Pool erfüllt
sie". Beides zugleich ist nicht zu haben — **compilerseitig gemessen**:
`pgx.Rows` trägt **10** Methoden, die deklarierte `sqlexec.Rows` fordert **4**;
`*pgxpool.Pool` erfüllt `Executor` **genau dann**, wenn `Query` `pgx.Rows`
liefert (mit der minimalen `Rows` bricht die Zusicherung ab — `wrong type for
method Query: have Query(...) (pgx.Rows, error), want Query(...) (minRows,
error)`). Die Verengung wäre über einen Vermittler zu haben, der `pgx.Rows` auf
vier Methoden eindampft; ihr Preis ist eine zusätzliche Schale **und** eine
`DB`-Zusicherung, die dann nicht mehr der Pool trägt. **Entschieden: die
Zusicherung bleibt** — der Zweck ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5) ist die **reine Funktion**, und die trägt; die Abhängigkeit auf den
**konkreten Pool** ist gelöst. **Benannte Grenze:** der Fake muss deshalb
**10** Methoden erfüllen statt vier (`translate_test.go`); die sechs
überzähligen (`CommandTag`, `FieldDescriptions`, `Values`, `RawValues`, `Conn`,
`TypeMap`) sind `nil`-Attrappen.

**Die Rechtfertigung ist Design, nicht die Zahl**
([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5): schmale Abhängigkeit statt konkreter Typ, Fehlerklassifikation als
reine Funktion. Dass der Coverage-Wert danach steigt, ist **Folge** — wer ihn
zum Zweck nimmt, hat den Gegenstand gewechselt.

**In den Slice aufgenommen — die Umsetzung von [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md).**
Die Naht **bewegt** die Messung: der Gegenstand der DB-Adapter-Coverage schrumpft
um 138 Statements (788 → 650), und die Quote der abfließenden Seite fällt
(75,25 % → 73,38 %), weil 116 der verlagerten überdurchschnittlich gedeckt waren.
**Der aufnehmende Gegenstand wächst aber um 152, nicht um 138** (1679 → **1831**):
die Naht *fügt* 14 Statements neuen Code hinzu (die generalisierte Übersetzung
samt `Classify`/`IsAbsent`) — die Summe ist damit **nicht** konstant
(2467 → 2481). Der Architect hat das als **Subjekt-Transfer** entschieden und die
Rampen neu bemessen (DB-Einstieg 75 → 70 %, Unit-Einstieg 65 → 70 %, Endstufen
unverändert 80 %); die Bindung, die das von einem Freibrief trennt: eine Quote,
die bei **unverändertem** Nenner fällt, ist eine Regression — dann steht die
Schwelle.

**Die Messung widerlegt die Voraussetzung dieser
Entscheidung** (§3, *Der Transfer-Nachweis*): `k_aufnehmend` 152 ≠
`k_abfließend` 138. **Entschieden durch [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md):**
die Summen-Konstanz war ein *hinreichender, aber weder notwendiger noch
hinreichender* Stellvertreter — sie ließe eine Löschung-mit-Ersatz durch (nicht
hinreichend) und verböte eine Generalisierung, die Code kostet (nicht
notwendig). An ihre Stelle tritt der **dreiteilige Transfer-Nachweis**
(§2, Liefer-Punkt 3). Der **Regressions-Riegel** bleibt unverändert scharf:
fällt die Quote bei **unverändertem** Nenner, oder fällt der abfließende Nenner
**ohne Ankunft** (`k_auf < k_ab`), steht die Schwelle.
**Berichtigung einer eigenen Zahl:** die hier zuvor stehende Angabe
`1679 → 1817` war die *abgeleitete* Summe `1679 + 138` — aus einem
Implementer-Bericht **übernommen und nicht gemessen**; sie ist über meinen
Auftrag an den Architect auch in `ADR-0077` §Kontext eingegangen und dort
ebenfalls falsch (siehe §7).

**Warum das hierher gehört und nicht in einen eigenen Slice:** diese Umsetzung
macht die DoD-Zeile „die reale Verdrahtung geht unverändert durch
`make test-store`/`make test-replication` (Exit 0)" erst **wahr** — ohne sie
bliebe der Slice dauerhaft rot an einer Kalibrierung, die seine eigene Bewegung
ausgelöst hat —, und das WIP-Limit lässt keinen zweiten Slice daneben zu.
**Wenn das Review das als Schnitt-Verstoß wertet**, ist der vorgesehene Weg die
Rückführung `in-progress → next` mit einer Zerlegung in *Naht* und
*Rampen-Nachzug*; §4 nennt die Rückführung vorab.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine coverage-getriebene Refaktorierung** — ein Umbau, dessen Zweck die
  Zahl ist, ist ausgeschlossen; die Naht muss aus sich heraus tragen.
- **Ein Umbau der SQL-Texte** (`queries` bleibt, wo es ist) — die Naht betrifft
  die **Ausführung**, nicht die Abfrage.
- **Eine neue Abhängigkeit** — die Schnittstelle ist adapter-eigen; `pgxpool`
  bleibt der reale Träger.
- **Verhaltensänderungen** — die reale Verdrahtung muss unverändert durch
  `make test-store`/`make test-replication` gehen.
- **Ein Fake als Ersatz der realen DB-Tests** — die Fake-Seite prüft die
  **Verklebung**; das SQL prüft weiterhin der echte PostgreSQL. Ein Fake, der
  die DB-Tests ersetzt, wäre eine Zusage, die er nicht hält.

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

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — die Naht existiert.**

- [x] Die Adapter hängen an einer **adapter-eigenen** Schnittstelle
      (`Query`/`Exec`/`QueryRow`) statt an `*pgxpool.Pool`; der reale Pool
      erfüllt sie (`var _ sqlexec.DB = (*pgxpool.Pool)(nil)`,
      `internal/adapters/driven/postgresstorage/sqlexec/seam.go`). **Die
      Zeilen-Schnittstelle bleibt `pgx.Rows`** — die Verengung auf vier
      Methoden verträgt sich nicht mit dieser Zusicherung (§1, compilerseitig
      gemessen).
- [x] **Kein Verhaltens-Change:** die reale Verdrahtung geht unverändert durch
      `make test-store` und `make test-replication` (Exit 0) — der Beleg, dass
      die Naht die Ausführung nicht verschiebt.

**Liefer-Punkt 2 — die reine Logik ist prüfbar.**

- [x] Fehlerklassifikation und Zeilen-Übersetzung sind **reine Funktionen** und
      netzlos geprüft (Fälle: Erfolg, Fehlerklasse, leeres Ergebnis).
- [x] Die Fake-Seite fährt die **Verklebung** (Query/Exec-Aufruf, Scan-Schleife,
      Fehlerpfad) — und ist ausdrücklich **kein** Ersatz der realen DB-Tests.

**Liefer-Punkt 3 — die Wirkung ist benannt, nicht angestrebt.**

- [x] Der Coverage-Effekt wird **als Folge** dokumentiert (Zahl vorher/nachher
      über der netzlos prüfbaren Fläche) — nicht als Zweck: **69,70 % → 71,30 %**
      (`make coverage-gate`, beide Läufe Exit 0; Gegenstand **1679 → 1831**
      Statements — §3).
- [x] **Die Neu-Bemessung aus [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
      in der Fassung von [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
      ist umgesetzt:** `DB_COVERAGE_THRESHOLD` 75 → 70
      (`tools/harness/db-coverage.sh`, Zeile der Vorgabe), `THRESHOLD` 65 → 70
      (`harness/mk/coverage.mk`), und die Träger-Doku
      (`harness/sensors/db-adapter-coverage.md`,
      `harness/sensors/coverage-gate.md`, `harness/README.md` §Sensors,
      `AGENTS.md` §4) nennt die neuen Stufen.
- [x] **Der Transfer-Nachweis liegt vollständig bei — alle drei Belege**
      (`ADR-0078` §Entscheidung): **(a) die Arithmetik** — `k_ab` 138, `k_auf`
      152, `k_auf ≥ k_ab`, die Differenz 14 ist **neuer** Code; **(b) der
      Paket-Diff**, der Abfluss ist auf den verlagerten Träger **isoliert**
      (übrige Pakete byte-identisch), der Zuwachs erscheint im neuen Paket
      `sqlexec`, und der verlagerte Code ist im abfließenden Gegenstand
      **vollständig abgegangen** — die aggregierte Summe genügt dafür
      ausdrücklich **nicht**; **(c) kein Verhalten verloren** — die realen,
      dienst-gestützten Läufe des abfließenden Gegenstands sind grün, kein
      Testfall entfernt. Träger: §3.
- [x] `make gates` grün (Exit direkt, ungepiped).

- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      `docs/reviews/review-slice-081.md` · Fixrunde
      `docs/reviews/review-slice-081-fixrunde.md` — 0 HIGH in der Fixrunde; die
      verbleibende LOW liegt im Plan-Text (§3/§1, Planner-Zug), deshalb ohne
      Rückgabe-Pfeil an den Implementer.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Zuschnitt des Implementer-Laufs** (die Liste nennt die Träger; der genaue
Schnitt entsteht hier):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgresstorage/sqlexec/**` (`seam.go`, `statement.go`, `errors.go`, `translate.go`) | neu | die Naht selbst — `Executor`/`DB` (`Query`/`Exec`/`QueryRow`; Zeilen als **`pgx.Rows`**, §1) — plus die von ihr getragene Zeilen-Übersetzung und Fehlerklassifikation. Die Naht liegt **außerhalb** der drei DB-Pakete und damit im Gegenstand des Unit-Gates (`ADR-0071` Punkt 5) |
| die sechs pool-tragenden Adapter-Dateien (`store`, `consumerstate`, `heartbeat`, `tableactivation`, `schemastore`, `administrationrequest`) · `schema.go` | refactor | sie hängen an `sqlexec.DB`/`sqlexec.Executor` statt an `*pgxpool.Pool`; `sqlexec.Classify`/`IsAbsent` ersetzen die inline gebildete Fehlerklasse und `errors.Is(err, pgx.ErrNoRows)` |
| `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` | neu | Fakes (`pgx.Rows`/`pgx.Row`/`Executor`) + die Fälle Erfolg · Fehlerklasse · Leerfall; **kein** Ersatz der realen DB-Tests |
| `tools/harness/db-coverage.sh` · `harness/mk/coverage.mk` | update | die beiden Zahlen der Neu-Bemessung aus `ADR-0077`: `DB_COVERAGE_THRESHOLD` 75 → 70, `THRESHOLD` 65 → 70. Die Endstufen bleiben 80 |
| `harness/sensors/db-adapter-coverage.md` · `harness/sensors/coverage-gate.md` · `harness/README.md` §Sensors · `AGENTS.md` §4 | update | sie tragen die geltenden Stufen und die Rampe; ohne sie stünde die Zahl im Werkzeug und eine andere in der Bindung |

**Abweichung von der Träger-Liste — `postgresack/ack.go` und
`receive/receive.go` bleiben unberührt.** Ihre Naht ist eine andere: beide
hängen an `*pgconn.PgConn`, nicht an `*pgxpool.Pool`, und `receive` ruft
`pglogrepl.StartReplication`/`CreateReplicationSlot` (Signatur: konkreter
`*pgconn.PgConn`) — die in §1/§2 genannte Schnittstelle (`Query`/`Exec` plus
minimales `Rows`) drückt diese Aufrufe nicht aus. Der Zug ist damit die
Rückführung aus §4 („nach Paket, je Adapter ein Slice"), nicht ein stilles
Weglassen: `postgresack` und `receive` gehören je in einen eigenen Vorgang —
**und sie existieren als Dateien**, damit die Folge-Slice-Paarung auflöst
(Review `review-slice-081` F-3):
`slice-084` (`../open/slice-084-postgresack-naht.md`) und
`slice-085` (`../open/slice-085-receive-naht.md`).
**Gemessen, nicht pauschal:** `postgresack` berührt **6**
`pgconn`/`pglogrepl`-Symbole, `receive` **18** — die Flächen sind eine andere
Größenordnung, deshalb zwei Vorgänge und nicht einer.
Berührte Träger dieses Laufs: **ein** Paket, **eine** Schicht; drei
Liefer-Punkte, kein vierter.

**Nicht angefasst:** `queries` (die SQL-Texte bleiben unberührt — die Naht
trägt sie als `sqlexec.Statement.SQL` durch), `.a-check.yml` (grün, 0 Befunde),
`spec/**`.

**Befund des ersten Laufs — die Naht bewegt die Messung.** Die Naht zieht die
Übersetzung aus `postgresstorage` heraus — und damit **aus dem Gegenstand der
DB-Adapter-Coverage** (`ADR-0071` Punkt 3), die genau dieses Paket misst. Der
reale Lauf: `make test-store` grün, `make test-replication` **rot**
(DB-Adapter-Coverage 75,25 % vorher → 73,38 % nachher, Stufe 75 %).
`ADR-0071` Re-Evaluierungs-Trigger (b) verlangt für das Ziehen der Naht die
Neu-Bemessung der Rampe als **Folge-ADR** (Architect); `AGENTS.md` §3.6
schließt eine Schwellen-Senkung durch den Implementer aus.

**Der Transfer-Nachweis — alle drei Belege** ([`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
§Entscheidung). Die Zahlen sind gemessen, nicht behauptet (gepinnte Toolchain-
und PostgreSQL-Images; Läufe dieses Zugs, Exit-Codes ungepiped):

| Gegenstand (Statements) | vor der Naht | nach der Naht | Δ |
|---|---|---|---|
| netzlos prüfbare Fläche (Unit) | 1679 (1171 gedeckt, 69,74 %) | 1831 (1306 gedeckt, 71,33 %) | **+152** |
| DB-Adapter-Gegenstand | 788 (593 gedeckt, 75,25 %) | 650 (477 gedeckt, 73,38 %) | **−138** |
| Summe (Lesehilfe) | 2467 | 2481 | +14 |

Die **gedeckte Zahl** des Unit-Gegenstands schwankt lauf-zu-lauf um wenige
Statements (beobachtet ±2; ein `ctx`-abhängiger Pfad in
`grpcstream/broadcaster.go` — hier 1306, im `make gates`-Lauf desselben Zugs
1308); Nenner (1831) und Abstand zur Schwelle sind davon unberührt.

**(a) Ankunft.** Der abfließende Gegenstand verliert `k_ab` **138** Statements
(`postgresstorage` 610 → 472; je Baum netzlos gemessen), der aufnehmende wächst
um `k_auf` **152** (`1679 → 1831`). `k_auf ≥ k_ab` ist erfüllt; die Differenz
**14** ist der **neue** Code des Zugs (die generalisierte Übersetzung samt
`Classify`/`IsAbsent`/`Statement.fail`), keine verlagerte Menge.

**(b) Paket-Granularitäts-Diff — nicht die Summe.** Der Abfluss ist auf **einen**
Träger **isoliert**: `git diff --name-only 4c7ea8a~1..8e9fe4f -- internal/` listet
ausschließlich Dateien unter `internal/adapters/driven/postgresstorage/`; kein
anderes Paket ist berührt, `postgresack` (23) und `replication/receive` (155)
bleiben unverändert. Der Zuwachs erscheint in **einem Träger, den es vorher
nicht gab** — das neue Paket `postgresstorage/sqlexec` (`errors.go` 2 ·
`statement.go` 3 · `translate.go` 147 = **152**); jedes andere Paket des
Unit-Gegenstands ist unverändert. Der verlagerte Code ist im abfließenden
Gegenstand **vollständig abgegangen** (kein Rest der Übersetzung: kein
`collectRecords`/`pgx.Rows` in den DB-Gegenstands-Paketen — kein Treffer). Die
**aggregierte Summe** (2467 → 2481) trägt diesen Nachweis **nicht**: sie
unterscheidet `−138 abgewandert / +14 neu` nicht von `−138 abgewandert / 0 neu`;
welcher Fall vorliegt, zeigt allein der Paket-Schnitt.

**(c) Kein Verhalten verloren.** Die realen, dienst-gestützten Läufe des
abfließenden Gegenstands sind **grün** (`make test-store` Exit 0,
`make test-replication` Exit 0, ungepiped), und **kein Testfall wurde entfernt**
(`git diff --name-status 4c7ea8a~1..8e9fe4f -- '*_test.go'` zeigt genau eine **neue**
Datei, `sqlexec/translate_test.go` — keine Löschung).

**Die Neu-Bemessung ist umgesetzt:** `DB_COVERAGE_THRESHOLD` 75 → 70
(`tools/harness/db-coverage.sh`), `THRESHOLD` 65 → 70 (`harness/mk/coverage.mk`);
die Träger-Doku nennt die neuen Stufen. Die **Rot-Gegenprobe** zeigt, dass die
neue Schwelle weiter prüft: `DB_COVERAGE_THRESHOLD=75` gegen denselben Stand
endet `db-coverage: FAIL — DB-Adapter-Coverage 73.38% unter Schwelle 75%`,
**Exit 1**. Der Regressions-Riegel (`ADR-0078` Entscheidung 2) bleibt scharf:
eine Quote, die bei **unverändertem** Nenner fällt, ist eine Regression.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 ist entschieden, `Verantwortlich:` gesetzt, WIP-Limit frei — und eine
Umgebung, in der `make test-store`/`make test-replication` laufen (die realen
Belege dieses Slice).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich der Umbau
  als breiter als **eine** Schicht oder als mehr als drei Liefer-Punkte (die Naht
  greift in die Adapterdateien mehrerer Pakete), gehört er zurück zur Zerlegung —
  **nach Paket** (je Adapter ein Slice), nicht nach Schichten.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass die Naht ohne
  **Verhaltensänderung** nicht zu haben ist (etwa weil der Pool Aufrufe trägt,
  die eine schmale Schnittstelle nicht ausdrückt), ist das ein Blocker mit
  Entscheidung: die Zusage dieses Slice ist „kein Verhaltens-Change".

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make test-store`/`make test-replication` grün (die reale
Verdrahtung) **und** `make gates` grün **und** der Coverage-Effekt **als Folge**
beziffert **und** die Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Umbau könnte ein Verhaltens-Change sein, der als Refactoring auftritt.**
  Wächter sind die **realen** DB-Tests (`make test-store`,
  `make test-replication`): sie müssen unverändert grün sein — sie sind der
  Beleg, nicht die neuen Fakes. — **Ausgang:** <bei Closure>
- **Der Fake könnte grün sein, ohne etwas zu prüfen.** Er prüft die
  **Verklebung** (Aufruf, Scan-Schleife, Fehlerpfad), nicht das SQL — ein Fake,
  der die DB-Tests ersetzt, hält die Zusage nicht. — **Ausgang:** <bei Closure>
- **Die Zahl könnte den Umbau ziehen.** `ADR-0071` Punkt 5 schließt das
  ausdrücklich aus: die Begründung muss **vor** dem Umbau stehen (Design), nicht
  danach (Wirkung). — **Ausgang:** <bei Closure>

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt die driven Adapter, ihre Tests
und die Schichten-Konfiguration in **einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; zwei Treffer:

- `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (2×, offen): **kein**
  Treffer — Kommentar-Wahrheit, nicht Adapter-Struktur; der Slice berührt
  Kommentare nur, soweit die Naht sie erklärt.
- `BEO-PGC/adapter-fehler-ausgang` (2×, offen): **kein** Treffer, aber
  **benachbart** — der Eintrag führt die Adapter-Grenze „kontrollierte
  Fortsetzung beim Aufrufer" und einen Retry-/Backoff-Aufschub. Dieser Slice
  ändert **nicht**, was der Aufrufer tut; er ändert, **woran** der Adapter hängt.

**Kein** Eintrag der Sichtung steht über der Schwelle, und **keiner** rückt mit
ihm auf 3×. Das ist der Stand **dieser Planung**, kein Versprechen über den
Lauf — erreicht einer die Schwelle doch, steht er in §7, und sein Ausgang fällt
dem Lese-Schritt der laufenden Welle-Closure zu (Modul 6), nicht dieser Closure.
Die Zählerstände der beiden Treffer sind am Register **nachgezählt**, nicht aus
diesem Text übernommen (je 2×, `ls evidence/`).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (`.a-check.yml` prüft die Schichten-Edges,
`spec/architecture.md` führt die Ports-und-Adapter-Ordnung), **Phase-Reife**
hoch, **Evidenz-/Diskrepanz-Risiko** **mittel** — der Umbau kann Verhalten
verschieben, und genau dagegen stehen die **realen** DB-Tests als Wächter (§6) —,
**Reconciliation-Aufwand** null.
