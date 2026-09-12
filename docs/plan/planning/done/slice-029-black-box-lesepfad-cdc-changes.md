# Slice slice-029: Black-Box-Lesepfad über `cdc.changes` — Vertragstest gegen `BEO-PGC/lese-doppelquelle`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-9`](welle-9.md) — der Nachweis, dass Lesen über die
externe SQL-Sicht real funktioniert und mit dem internen Go-Lesepfad
übereinstimmt, ist welle-9s Closure-Trigger (§3), kein Einzel-Slice-DoD.

**Bezug:** [`LH-FA-REA-002`](../../../../spec/lastenheft.md) (Lesen über
den ausgelieferten Zugriffsweg) — dieser Slice erweitert die Testabdeckung
für einen bestehenden Vertrag, ändert ihn nicht. Kein aktives ADR wird
geändert.

**Berührte Spec-Stellen:** — (reine Testinfrastruktur, keine
Verhaltensänderung an einer Spec-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Compose-Integrationstest bekommt einen neuen Testfall, der
Changes **ausschließlich über die SQL-Sicht `cdc.changes`** liest (rohes
`docker exec … psql`/SQL, kein `postgresstorage.NewChangeStoreAdapter`/
`ReadChanges`), für denselben realen Rundlauf (INSERT/UPDATE/DELETE), den
der bestehende `TestMVPCaptureFlow` bereits weiß-box über den Go-Adapter
liest. Beide Lesungen desselben Datensatz werden verglichen (Reihenfolge,
Feldinhalt). Das ist der fehlende Vertragstest für
`BEO-PGC/lese-doppelquelle`: driftet die SQL-View-Semantik von der
Go-Use-Case-Semantik, schlägt dieser Test künftig sichtbar fehl, statt
unbemerkt zu bleiben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ersetzung der bestehenden Go-Adapter-Lese-Tests** — Bestand bleibt
  bewusst stehen: Der Go-Adapter bleibt der primäre interne Lesezugriffsweg
  (z. B. für künftige Consumer-Zugriffswege); dieser Slice ergänzt einen
  externen Vertragstest, ersetzt keinen bestehenden.
- **Schema-Änderungen** — Folge-Slice `slice-030` übernimmt das.
- **CDC-Verwaltung/Observability black-box** — bereits als eigene, spätere
  Welle vorgemerkt (Roadmap *Nächste Wellen*).
- **`BEO-PGC/lese-doppelquelle` für alle vier SQL-Views verallgemeinern**
  (`active_tables`, `consumer_status` sind nicht Teil dieses Slices) —
  anderer Vorgang: Der Vertragstest deckt hier gezielt `cdc.changes` ab,
  die View mit der komplexesten Lese-Semantik (`LH-FA-REA-001…006`); die
  übrigen Views haben einfachere, weniger driftanfällige Projektionen.

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

- [x] `LH-FA-REA-002` erfüllt: neuer Testfall liest Changes real über
      `cdc.changes` (rohes SQL, kein Go-Adapter-Import) für denselben
      INSERT/UPDATE/DELETE-Rundlauf, den `TestMVPCaptureFlow` bereits über
      `ReadChanges` liest. Beleg: `TestMVPChangesViewMatchesReadChanges` in
      `test/integration/integration_test.go` (Commit `7b7ca3c`).
- [x] Vertragstest belegt Übereinstimmung: dieselbe Reihenfolge
      (`commit_position`, `sequence`), derselbe Feldinhalt (`operation`,
      `old_data`/`new_data`, `schema_version`) zwischen SQL-View-Lesung und
      Go-Adapter-Lesung desselben Datensatzes. Beleg: drei grüne
      `make test-integration`-Läufe (kein Feldunterschied); zusätzlich real
      als roter Fund verifiziert — `old_data`/`new_data` in der
      `changes`-View-Spaltenprojektion testweise vertauscht, derselbe
      Testfall lief rot (einzig er, die übrigen vier blieben grün), Mutation
      danach zurückgesetzt (`tools/schema/schema.yaml` unverändert laut
      `git diff`).
- [x] `make gates` grün, `make test-integration` dreimal in Folge grün.
      Beleg: `make gates` grün nach dem Test-Commit; `make test-integration`
      dreimal in Folge grün im Anschluss an die Mutationsprobe (drei weitere
      Läufe auf dem zurückgesetzten Stand).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-029.md`](../../../reviews/review-slice-029.md)
      (1 LOW, F-1 behoben in Commit `f96eff3`); Verifikation in
      [`docs/reviews/verify-slice-029.md`](../../../reviews/verify-slice-029.md)
      (Mutationsprobe eigenständig reproduziert, DoD-Konformität bestätigt).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird — keiner
      berührt: reine Testabdeckung eines bestehenden Lesezugriffswegs, kein
      Guide/Sensor/Vertrag geändert. Kein Doku-Update vorgenommen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` §Modus-Deklaration), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Beleg: `evidence/slice-029.md` in `BEO-PGC/lese-doppelquelle/` angelegt
      — Zähler damit bei 3×. Ausgang **verkörpert**: der neue Testfall ist
      der dauerhafte Sensor, den die Beobachtung vermisste (real
      mutations-geprüft, siehe Verifikationsbericht). Siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der `welle-9`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | neuer Testfall `TestMVPChangesViewMatchesReadChanges`: SQL-Lesung gegen `cdc.changes` vs. bestehende `ReadChanges`-Lesung desselben Rundlaufs |

**Plan-Nachzug (Abweichungen, nach Bestandsprüfung entschieden):**

- **Go-Testfall statt Bash/`docker exec … psql`.** Der bestehende
  `mvpEnv.pool` (pgx) verbindet bereits gegen dieselbe Compose-Instanz; ein
  direktes `pool.Query` gegen `cdc.changes` liest den externen
  SQL-Lesezugriffsweg genauso „von außen" wie ein `docker exec … psql`
  (kein Import von `postgresstorage.NewChangeStoreAdapter`/`ReadChanges`),
  ohne eine zweite Sprache (Bash-Heredoc/`psql`-Parsing) im selben Vergleich
  zu tragen. `run-integration-tests.sh` bleibt unverändert.
- **Eigener, isolierter Rundlauf statt Wiederverwendung von
  `TestMVPCaptureFlow`s Zeilen.** Eine feste Zeilen-ID (`id=60` auf
  `feed_mvp_full`, das über den ganzen Testlauf aktiviert bleibt) macht den
  Vergleich unabhängig von der Ausführungsreihenfolge der übrigen
  Testfälle dieser Datei — anders als eine Kopplung an `TestMVPCaptureFlow`s
  Zustand, die bei künftigen Testfall-Einfügungen zwischen beiden brechen
  könnte.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  eine reale Drift zwischen SQL-View und Go-Adapter existiert, die eine
  Korrektur an einer der beiden Implementierungen braucht (mehr als
  Testinfrastruktur), gehört das zurück zur Zerlegung — ein Fund, kein
  Scheitern.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
dreimal in Folge grün **und** Closure-Notiz geschrieben (inkl.
`BEO-PGC/lese-doppelquelle`-Ausgang).

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Vertragstest könnte eine reale, bisher unbemerkte Drift zwischen
  SQL-View und Go-Adapter aufdecken (z. B. unterschiedliche Behandlung von
  `NULL`-Werten in `old_data`/`new_data`) — das wäre ein echter Fund, kein
  Testinfrastruktur-Defekt, und bräuchte eine Architect-Entscheidung, welche
  Semantik korrekt ist. **Ausgang: entfallen** — keine Drift gefunden;
  SQL-View- und Go-Adapter-Lesung stimmen für den realen INSERT/UPDATE/
  DELETE-Rundlauf real überein (dreimal grün, zusätzlich durch die
  Mutationsprobe als scharfer Vergleich bestätigt, nicht nur trivial grün).
- Der Vergleich zweier Lesungen (SQL-View, Go-Adapter) im selben Testlauf
  könnte durch Timing/Nebenläufigkeit einer laufenden Erfassung flaky
  werden, wenn nicht sauber auf „alle erwarteten Changes erfasst" gewartet
  wird. **Ausgang: entfallen** — `awaitChangesViewRows` wartet mit
  Deadline/Poll-Intervall auf die vollständige Zeilenzahl, bevor verglichen
  wird; Implementer- und Verifier-Läufe (insgesamt acht `make
  test-integration`-Durchläufe über beide Rollen) zeigten keine Flakiness.

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

- **Was hat funktioniert:** Der Verifier hat die vom Reviewer offen
  gelassene Prüfgrenze (Mutationsbeleg im Repo-Zustand nicht nachvollziehbar)
  nicht als Lücke stehen gelassen, sondern die Mutationsprobe selbst
  reproduziert — inklusive der Erkenntnis, dass eine reine Spalten-
  *Reihenfolge*-Vertauschung wirkungslos gewesen wäre (die generierte
  `CREATE VIEW`-SQL trägt keine explizite Spaltenliste) und stattdessen ein
  `AS`-Alias-Tausch nötig war. Damit ist der zentrale Sicherheitsbeleg dieses
  Slices unabhängig von Implementer und Reviewer ein drittes Mal bestätigt.
- **Was ging anders als geplant:** Keine wesentliche Abweichung vom
  Slice-Ziel; die im Plan-Nachzug (§3) begründeten zwei Entscheidungen
  (Go/pgx statt `docker exec … psql`, eigener isolierter Rundlauf statt
  Wiederverwendung von `TestMVPCaptureFlow`) trugen. Der Reviewer fand ein
  LOW-Finding (F-1: Kopplungs-Kommentar nannte nur zwei von drei
  koexistierenden ID-Gruppen) — direkt behoben (Commit `f96eff3`), kein
  Sachdefekt.
- **Steering-Loop-Eintrag:** `BEO-PGC/lese-doppelquelle` erreicht mit
  diesem Slice 3× und geht direkt auf Ausgang *verkörpert* — der neue
  Testfall ist der dauerhafte Sensor, den die Beobachtung seit `slice-010`
  vermisste, läuft mit jedem `make test-integration` mit und wurde real
  mutations-scharf verifiziert. Verkörpert in
  `test/integration/integration_test.go` (`TestMVPChangesViewMatchesReadChanges`)
  — liegt in `test/integration/integration_test.go`. Auslöser: `BEO-PGC/lese-doppelquelle`
  (slice-010, slice-011, slice-029 — 3×).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-029.md`
  in `BEO-PGC/lese-doppelquelle/` ergänzt — Zähler steht damit bei 3×,
  Ausgang **verkörpert** (`state.md` fortgeschrieben, Anker `seit slice-029`).
- **Folge-Slices:** keine — `slice-030` (Black-Box-E2E für
  Schema-Änderungen) ist bereits als Datei in `open/` vorhanden, unabhängig
  von diesem Slice geplant, kein neuer Folge-Slice dieses Vorgangs.
- **Risiken aus §6:** beide *entfallen* — siehe §6 für Begründung.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-9` offen) —
  Prüfung läuft bei der `welle-9`-Closure.

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
Treffer für `PGC`: `BEO-PGC/lese-doppelquelle` (2×, weiter offen — dieser
Slice liefert den fehlenden Vertragstest), `BEO-PGC/test-isolation-geteilter-zustand`
(1×, nicht einschlägig — anderer Test-Layer). Keiner der übrigen Treffer
erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
