# Slice slice-067: Assembler-Filterung und Live-Reload-Verdrahtung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-18`](../welle-18.md) — baut auf `slice-066` auf (die
Antrags-Verarbeitung ruft die hier gebaute Live-Reload-Methode auf);
`slice-068` (E2E-Beleg) baut auf diesem Slice auf.

**Bezug:** [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Haupt-Bezug,
insbesondere die Boundary-Klausel zu `LH-FA-SCH-003`),
[`LH-FA-SCH-003`](../../../../spec/lastenheft.md) (Konvergenz-Verhalten bei
realer Spaltenlöschung), [`LH-FA-DAT-005`](../../../../spec/lastenheft.md)
(Negative-Abwesenheitssemantik, auf der die Filterung aufbaut),
[`ADR-0059`](../../../../docs/plan/adr/0059-spaltenauswahl-mechanismus.md)
(nur umgesetzt — keine aktive ADR wird geändert, `ADR-0059` bleibt
`Accepted`), [`ADR-0050`](../../../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(bestehender Live-Reload-Mechanismus, `tablesMu`-Synchronisation).

**Berührte Spec-Stellen:** [`ARC-004`](../../../../spec/architecture.md)
(Assembler-Zuständigkeit erweitert um Filterung).

**Verantwortlich:** —.

**Autor:** pt9912. **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `TableBinding`
(`internal/adapters/driving/replication/mapper/mapper.go`) um
`ExcludedColumns` erweitern; `rowImage`/`change` filtern ausgeschlossene
Spalten vor jeder Serialisierung (nutzt die bestehende `nil`-Abwesenheits-
Kodierung, kein neuer Platzhalter — `LH-FA-DAT-005`); eine neue, unter
`tablesMu` synchronisierte `Assembler`-Methode für den Live-Reload-Nachtrag
eines Ausschluss-Standes. **Zentrale Konsequenz aus `ADR-0059` Teilfrage 3,
explizit zu lösen:** Der bestehende Schema-Bump-Pfad
(`observeRelation`/`AddBinding(relation.QualifiedName(), TableBinding{TableID:
..., SchemaVersion: nextID})`, Zeile ~376) überschreibt heute die komplette
`TableBinding` und würde einen zuvor gesetzten Ausschlussstand
stillschweigend zurücksetzen — dieser Slice baut das so um, dass ein
Schema-Versions-Nachtrag den bestehenden `ExcludedColumns`-Stand einer
Bindung **erhält**, kein blindes vollständiges Überschreiben mehr.
Unit-Tests: ein ausgeschlossener Spaltenschlüssel erscheint nie in
`old_data`/`new_data` (Happy Path **und** nach einem Schema-Bump derselben
Tabelle), plus ein Konvergenz-Test für eine real gelöschte, zuvor
ausgeschlossene Spalte — derselbe `ErrIncompatibleSchemaChange`-Pfad wie
jede andere Spaltenlöschung (`LH-FA-SCH-003`/`LH-FA-SCH-004.a`,
`ADR-0059` Teilfrage 4: kein Sonderfall, reine Schichtung reicht).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **SQL-Funktionen und Antrags-Verarbeitung** (`cdc.exclude_column`/
  `cdc.include_column`, `applyAdministrationRequest`-`case`-Zweige) —
  Vorgänger-Slice `slice-066` liefert das bereits; dieser Slice setzt nur
  voraus, dass ein verarbeiteter Antrag existiert, der die hier gebaute
  Methode aufrufen kann.
- **E2E-Beleg am laufenden Feed-Container** — eigener Slice `slice-068`.
- **Ein eigenes Diagnose-Metadatum, das „ausgeschlossen" von „real
  gelöscht" unterscheidet** — `ADR-0059` Teilfrage 4 entscheidet explizit
  gegen eine dritte Zustandskategorie (Option B verworfen): die Boundary-
  Klausel verlangt ausdrücklich nur „gilt das Verhalten aus
  `LH-FA-SCH-003`", keine eigene Diagnose-Kategorie.
- **Quellen-/musterweiter Ausschluss über mehrere Tabellen hinweg** —
  Bestand bleibt bewusst außen vor, `ADR-0059` Teilfrage 2 schließt Option
  C bewusst aus.

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

- [ ] `TableBinding` um `ExcludedColumns` erweitert; `rowImage`/`change`
      filtern ausgeschlossene Spaltenschlüssel real heraus, bevor das Row
      Image gebaut wird. Beleg: Unit-Test in
      `internal/adapters/driving/replication/mapper/mapper_test.go`, der
      belegt, dass ein in `TableBinding.ExcludedColumns` geführter
      Spaltenname nie als Schlüssel im resultierenden `old_data`/`new_data`
      erscheint.
- [ ] Neue synchronisierte `Assembler`-Methode für den Live-Reload-Nachtrag
      eines Ausschluss-Standes; der Schema-Bump-Pfad
      (`observeRelation`/`AddBinding`) erhält den bestehenden
      `ExcludedColumns`-Stand einer Bindung, statt ihn zurückzusetzen.
      Beleg: Unit-Test, der einen Ausschluss setzt, danach einen
      Schema-Bump derselben Tabelle simuliert, und erneut prüft, dass die
      ausgeschlossene Spalte weiterhin gefiltert wird; `go test -race`
      grün (neue Methode greift wie `AddBinding`/`RemoveBinding` unter
      `tablesMu`).
- [ ] Konvergenz-Test: eine real gelöschte, zuvor ausgeschlossene Spalte
      löst denselben `ErrIncompatibleSchemaChange`-Pfad aus wie jede andere
      Spaltenlöschung — kein Sonderfall (`ADR-0059` Teilfrage 4).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: keiner erwartet (kein neuer öffentlicher Vertrag; die
      Folgepflichten aus `ADR-0059` — `spec/architecture.md`-Korrektur,
      neuer `SPEC-*`-Eintrag — sind bereits `slice-066` zugeordnet);
      Implementer bestätigt oder begründet Abweichung im Plan-Nachzug.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb (`welle-18` offen) — Prüfung läuft bei der `welle-18`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `TableBinding.ExcludedColumns`, Filterung in `rowImage`/`change`, neue Live-Reload-Methode, Schema-Bump-Erhalt |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | update | Filter-, Schema-Bump-Erhalt- und Konvergenz-Tests |
| `internal/bootstrap/wiring.go` | update (kleiner Nachtrag) | ruft die neue Live-Reload-Methode aus `slice-066`s `applyAdministrationRequest`-Zweigen auf |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-066` liegt in `done/` (die
Antrags-Verarbeitung muss existieren, damit die hier gebaute Methode einen
Aufrufer hat), WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Filterung, Schema-Bump-Erhalt und Konvergenz-Test zusammen mehr als drei
  Liefer-Punkte oder mehr als zwei Schichten in einer Review-Sitzung nicht
  mehr prüfbar machen, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Der Schema-Bump-Erhalt
  verlangt einen tieferen Umbau von `TableBinding`s Konstruktions-Stellen,
  als hier angenommen (mehrere weitere Aufrufer, die `TableBinding`
  literal statt über eine Update-Methode konstruieren) — dann Carveout
  statt eines unvollständigen Umbaus.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün (inklusive `go test -race`) **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Schema-Bump-Erhalt könnte weitere, hier nicht erfasste Konstruktions-
  Stellen von `TableBinding` berühren (z. B. Erstaktivierung über
  `NewAssembler`) und dort eine andere Semantik brauchen (leerer
  Ausschluss bei Erstaktivierung ist korrekt, ein Bump muss dagegen
  erhalten) — Verwechslungsgefahr zwischen beiden Pfaden. — **Ausgang:**
  <bei Closure zuzuweisen>
- `go test -race` könnte einen Data Race an der neuen Methode aufdecken,
  wenn ein Aufrufer außerhalb von `tablesMu` auf `ExcludedColumns`
  zugreift (dieselbe Fitness-Function-Anforderung wie in `ADR-0059`
  §Fitness Function benannt). — **Ausgang:** <bei Closure zuzuweisen>

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

- **Was hat funktioniert:** <bei Closure>
- **Was ging anders als geplant:** <bei Closure>
- **Steering-Loop-Eintrag:** <bei Closure>
- **Beobachtungs-Register (`../observations/`):** <bei Closure>
- **Folge-Slices:** keine — `slice-068` (E2E-Beleg) steht bereits in
  `welle-18` §4 als vorgesehener nächster Slice.
- **Risiken aus §6:** <bei Closure>
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-18` offen) —
  Prüfung läuft bei der `welle-18`-Closure.

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
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC` mit Bezug zu `Assembler`/`TableBinding`/Schema-Evolution:
`BEO-PGC/schema-evolution-nicht-dynamisch` (2×, direkt aufgelöst/verkörpert
seit `slice-033` — dieser Slice berührt denselben Code-Bereich
(`observeRelation`/`classifyRelationColumns`-Nachbarschaft), verändert aber
nicht die dort verkörperte Re-Versionierungslogik selbst, sondern nur, wie
eine bestehende `TableBinding` beim Bump aktualisiert wird — kein Rückfall
in den ursprünglichen Fund zu erwarten). Keine weiteren Treffer für
`Assembler`/`TableBinding`/Live-Reload über die bereits in `slice-066`
gesichteten hinaus (`BEO-PGC/adapter-fehler-ausgang`,
`BEO-PGC/rollen-verdrahtung`, `BEO-PGC/verwaltung-keine-sql-administration`
— alle ohne neuen Bezug zu diesem Slice). Kein Treffer erreicht mit diesem
Slice 3× erstmals.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
