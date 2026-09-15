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
siehe §1 und §3).

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
schmale, adapter-eigene Executor-Schnittstelle (`Query`/`Exec` plus ein
minimales `Rows`) macht ihre Verklebung netzlos prüfbar: Fehlerklassifikation
und Zeilen-Übersetzung werden **reine Funktionen** mit eigenen Tests.

**Die Rechtfertigung ist Design, nicht die Zahl**
([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5): schmale Abhängigkeit statt konkreter Typ, Fehlerklassifikation als
reine Funktion. Dass der Coverage-Wert danach steigt, ist **Folge** — wer ihn
zum Zweck nimmt, hat den Gegenstand gewechselt.

**In den Slice aufgenommen — die Umsetzung von [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md).**
Die Naht **bewegt** die Messung: 138 Statements wandern aus dem Gegenstand der
DB-Adapter-Coverage in den Unit-Gegenstand (788 → 650 bzw. 1679 → 1817), und die
Quote der abfließenden Seite fällt (75,25 % → 73,38 %), weil 116 der verlagerten
Statements überdurchschnittlich gedeckt waren. Der Architect hat das als
**Subjekt-Transfer** entschieden und die Rampen neu bemessen (DB-Einstieg
75 → 70 %, Unit-Einstieg 65 → 70 %, Endstufen unverändert 80 %); die Bindung,
die das von einem Freibrief trennt: eine Quote, die bei **unverändertem** Nenner
fällt, ist eine Regression — dann steht die Schwelle.

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
      (`Query`/`Exec` plus ein minimales `Rows`) statt an `*pgxpool.Pool`;
      der reale Pool erfüllt sie (`var _ sqlexec.DB = (*pgxpool.Pool)(nil)`,
      `internal/adapters/driven/postgresstorage/sqlexec/seam.go`).
- [ ] **Kein Verhaltens-Change:** die reale Verdrahtung geht unverändert durch
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
      (`make coverage-gate`, beide Läufe Exit 0; Gegenstand 1679 → 1817
      Statements).
- [ ] **Die Neu-Bemessung aus [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
      ist umgesetzt:** `DB_COVERAGE_THRESHOLD` 75 → 70
      (`tools/harness/db-coverage.sh`, Zeile der Vorgabe), `THRESHOLD` 65 → 70
      (`harness/mk/coverage.mk`), und die Träger-Doku
      (`harness/sensors/db-adapter-coverage.md`,
      `harness/sensors/coverage-gate.md`, `harness/README.md` §Sensors,
      `AGENTS.md` §4) nennt die neuen Stufen. **Der Nenner-Nachweis** — die
      Statement-Summe beider Gegenstände bleibt über den Zug hinweg konstant —
      liegt als Beleg bei; er ist die Bedingung, unter der die Senkung trägt.
- [ ] `make gates` grün (Exit direkt, ungepiped).

- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
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
| `internal/adapters/driven/postgresstorage/sqlexec/**` (`seam.go`, `statement.go`, `errors.go`, `translate.go`) | neu | die Naht selbst — `Executor`/`DB` mit minimalem `Rows` — plus die von ihr getragene Zeilen-Übersetzung und Fehlerklassifikation. Die Naht liegt **außerhalb** der drei DB-Pakete und damit im Gegenstand des Unit-Gates (`ADR-0071` Punkt 5) |
| die sechs pool-tragenden Adapter-Dateien (`store`, `consumerstate`, `heartbeat`, `tableactivation`, `schemastore`, `administrationrequest`) · `schema.go` | refactor | sie hängen an `sqlexec.DB`/`sqlexec.Executor` statt an `*pgxpool.Pool`; `sqlexec.Classify`/`IsAbsent` ersetzen die inline gebildete Fehlerklasse und `errors.Is(err, pgx.ErrNoRows)` |
| `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` | neu | Fakes (`Rows`/`Row`/`Executor`) + die Fälle Erfolg · Fehlerklasse · Leerfall; **kein** Ersatz der realen DB-Tests |
| `tools/harness/db-coverage.sh` · `harness/mk/coverage.mk` | update | die beiden Zahlen der Neu-Bemessung aus `ADR-0077`: `DB_COVERAGE_THRESHOLD` 75 → 70, `THRESHOLD` 65 → 70. Die Endstufen bleiben 80 |
| `harness/sensors/db-adapter-coverage.md` · `harness/sensors/coverage-gate.md` · `harness/README.md` §Sensors · `AGENTS.md` §4 | update | sie tragen die geltenden Stufen und die Rampe; ohne sie stünde die Zahl im Werkzeug und eine andere in der Bindung |

**Abweichung von der Träger-Liste — `postgresack/ack.go` und
`receive/receive.go` bleiben unberührt.** Ihre Naht ist eine andere: beide
hängen an `*pgconn.PgConn`, nicht an `*pgxpool.Pool`, und `receive` ruft
`pglogrepl.StartReplication`/`CreateReplicationSlot` (Signatur: konkreter
`*pgconn.PgConn`) — die in §1/§2 genannte Schnittstelle (`Query`/`Exec` plus
minimales `Rows`) drückt diese Aufrufe nicht aus. Der Zug ist damit die
Rückführung aus §4 („nach Paket, je Adapter ein Slice"), nicht ein stilles
Weglassen: `postgresack` und `receive` gehören je in einen eigenen Vorgang.
Berührte Träger dieses Laufs: **ein** Paket, **eine** Schicht; drei
Liefer-Punkte, kein vierter.

**Nicht angefasst:** `queries` (die SQL-Texte bleiben unberührt — die Naht
trägt sie als `sqlexec.Statement.SQL` durch), `.a-check.yml` (grün, 0 Befunde),
`spec/**`.

**Befund dieses Zugs, der nicht dem Implementer gehört:** Die Naht zieht die
Übersetzung aus `postgresstorage` heraus — und damit **aus dem Gegenstand der
DB-Adapter-Coverage** (`ADR-0071` Punkt 3), die genau dieses Paket misst. Der
reale Lauf: `make test-store` grün, `make test-replication` **rot**
(DB-Adapter-Coverage 75,25 % vorher → 73,38 % nachher, Stufe 75 %). Die Naht
ist ohne Verhaltens-Change zu haben — die **Messung** aber nicht ohne
Entscheidung: `ADR-0071` Re-Evaluierungs-Trigger (b) verlangt für das Ziehen
der Naht die Neu-Bemessung der Rampe als **Folge-ADR** (Architect), und
`AGENTS.md` §3.6 schließt eine Schwellen-Senkung durch den Implementer aus.
Der Slice schließt mit diesem Zug deshalb **nicht**; die DoD-Häkchen unten
stehen nur, wo ein Beleg vorliegt.

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
