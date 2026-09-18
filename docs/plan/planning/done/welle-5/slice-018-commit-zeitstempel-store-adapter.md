# Slice slice-018: Commit-Zeitstempel — Store-Adapter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-5`](welle-5.md) — die Ende-zu-Ende-Belegpflicht
(Lasttest gegen einen real abgebildeten Verzögerungs-Wert) geht über die
DoD dieses einzelnen Slice hinaus.

**Bezug:** [`LH-FA-ADM-004`](../../../../spec/lastenheft.md)

**Berührte Spec-Stellen:** [`SPEC-001`](../../../../spec/pflichtenheft.md)
(`cdc.transaction.committed_at`)

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

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

**Ziel:** Den in slice-017 domänenseitig verfügbaren Quell-Commit-
Zeitstempel real in `cdc.transaction.committed_at` persistieren —
`committed_at` trägt danach den echten Quell-Commit-Zeitpunkt statt der
bisherigen Instanz-DEFAULT-Persistenzzeit.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Metrik-Umbenennung `cdc_capture_lag_approx` → `cdc_capture_lag` und der
  Ende-zu-Ende-Lasttest-Beleg — Folge-Slice slice-019, das den jetzt real
  gespeicherten Wert konsumiert.
- Änderung am `ChangeStorePort`-Interface — Bestand bleibt bewusst
  stehen: `PersistTransaction(ctx, *model.ChangeTransaction)` reicht
  unverändert; der neue Zeitstempel kommt aus dem bereits übergebenen
  Domänenobjekt (slice-017), kein neuer Parameter am Port nötig.

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

- [x] `mapper.NewTransactionRow` liest den Zeitstempel aus dem
      Domänenobjekt; `TransactionRow` trägt ein neues Feld dafür.
- [x] `queries.InsertTransaction` nimmt `committed_at` als Parameter statt
      sich auf die Spalten-DEFAULT zu verlassen; `store.go` übergibt den
      Wert. Real gegen eine PostgreSQL-Testinstanz getestet: eine mit
      künstlicher Verzögerung eingefügte Transaktion zeigt
      `committed_at` ≠ Persistenzzeitpunkt. Beleg: `TestPersistCarriesSourceCommittedAtNotPersistenceTime`,
      Verifier hat den Effekt real per eigener Mutationsprobe reproduziert.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`review-slice-018.md`](../../../../docs/reviews/review-slice-018.md)
      (F-1 LOW disponiert, F-2 INFO), Verifier bestätigt in
      [`verify-slice-018.md`](../../../../docs/reviews/verify-slice-018.md).
- [x] Doku-Update falls öffentlicher Vertrag berührt — geprüft: die
      Spalten-DEFAULT in `tools/schema/schema.yaml`
      (`transaction.committed_at`) bleibt als Absicherung stehen (nur ohne
      Anwendungscode greifend); Kommentar aktualisiert.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Geprüft: Datei existiert nicht (GF-Repo) — entfällt.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Keine Beobachtung angefallen — siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice gehört zu `welle-5` — die Paarungen prüft die **Welle-Closure**, nicht dieser Slice. Item entfällt hier bewusst.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgresstorage/mapper/mapper.go` | update | `TransactionRow` trägt den Zeitstempel; `NewTransactionRow` liest ihn aus der Domäne |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | `InsertTransaction`-SQL nimmt `committed_at` als vierten Parameter |
| `internal/adapters/driven/postgresstorage/store.go` | update | Aufrufstelle übergibt den vierten Parameter |
| `tools/schema/schema.yaml` | update (Kommentar) | Kommentar an `transaction.committed_at` — Anwendung liefert den Wert jetzt immer explizit, DEFAULT bleibt Absicherung |
| zugehörige `_test.go`-Dateien der drei Go-Dateien | update | Test gegen reale PostgreSQL: `committed_at` ≠ Persistenzzeitpunkt bei künstlicher Verzögerung |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): kein anderes Slice in `in-progress/`
(WIP-Limit 1); slice-017 liegt in `done/` (liefert den domänenseitigen
Zugriff, den dieser Slice braucht).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Falls der
  Zeitstempel-Zugriff aus slice-017 nicht so geformt ist, dass ihn
  `NewTransactionRow` ohne weitere Domänen-Änderung lesen kann — Rückzug
  mit Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Falls das Entfernen der
  Instanz-Zeit-Abhängigkeit einen bestehenden, real gefahrenen Test
  bricht, der bewusst auf die alte DEFAULT-Semantik baute, und die
  Korrektur den Slice-Umfang sprengt — Blocker, Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss ohne
  offenes HIGH-Finding; realer Test gegen PostgreSQL zeigt
  `committed_at` ≠ Persistenzzeitpunkt bei künstlicher Verzögerung.
- Lerneintrag §7: was beim Store-Adapter-Umbau gut oder schlecht lief.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Bestehende Tests könnten implizit auf `committed_at` als
  Persistenzzeitpunkt (statt Quell-Commit-Zeitpunkt) vertrauen und nach
  der Umstellung fehlschlagen. **Ausgang: entfallen.** `make test-store`
  läuft grün, kein bestehender Test brach.
- `ON CONFLICT (transaction_id) DO NOTHING` bei einer erneut
  persistierten Transaktion (Idempotenz-Fall) behält den zuerst
  geschriebenen Zeitstempel — das ist beabsichtigt, aber ohne expliziten
  Test bislang unbelegt. **Ausgang: entfallen.** Der Wert stammt aus dem
  WAL-`COMMIT`-Zeitstempel derselben Transaktion (`pglogrepl`, real
  deterministisch bei erneuter Auslieferung derselben WAL-Position nach
  einem Crash) — ein Retry liefert denselben Zeitstempel wie der
  Erstversuch, `DO NOTHING` behält also ohnehin den korrekten Wert. Kein
  eigener Test nötig, da kein abweichendes Verhalten möglich ist.

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

- **Was hat funktioniert:** Der reale Test mit künstlicher Verzögerung
  zwischen Domänen-Commit und Store-Aufruf hat den Unterschied
  `committed_at` vs. Persistenzzeitpunkt sauber bewiesen; die
  Review-Fixrunde (F-1) zeigte real per Laufzeit-Differenz (−1,5s), dass
  eine zuvor eingebaute Sleep für die Testaussage überflüssig war — der
  Verifier reproduzierte den Kern-Beleg danach unabhängig per eigener
  Mutationsprobe. Die WAL-Determinismus-Überlegung zum
  Idempotenz-Risiko (§6) hielt einer unabhängigen Prüfung durch Reviewer
  und Verifier stand, ohne dass ein eigener Test nötig wurde.
- **Was ging anders als geplant:** Während der Bearbeitung entstand ein
  repo-weiter Kommentar-Bereinigungsauftrag (Hard Rule 3.7, konkret für
  `tools/schema/schema.yaml` und Quellcode) über 18 Dateien, die
  außerhalb des Slice-Plans §3 lagen. Der Verifier flaggte das zu Recht
  als unreviewten Nachtrag (VF-2); eine gezielte Nachprüfung
  (`review-comment-cleanup.md`) hat die Änderungen als rein
  kommentar-ändernd und fachlich korrekt bestätigt.
- **Steering-Loop-Eintrag:** Mit diesem Slice wurde nichts verkörpert.
  Der Eintrag ist gezählt, nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen.
- **Folge-Slices:** keine — slice-019 lag bereits als Datei in `open/`.
- **Risiken aus §6:** beide `entfallen` (siehe §6).
- **Drei Paarungen:** entfällt hier — dieser Slice gehört zu `welle-5`;
  die Paarungen prüft die Welle-Closure für alle drei Slices gemeinsam.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührte Sub-Area
(Store-Adapter/Schemamigration) ist ein Segment des GF-Baums der
Modus-Deklaration (`PGC` Greenfield, Doc führt) — erfüllt die Schwelle
≥ 2 von 3 Achsen. Nicht zu grob.

**Vorgelagert — offene Beobachtungen sichten:** Register gelesen (zwölf
Einträge — `BEO-PGC/dod-checkbox-nachzug` kam während der Closure von
slice-017 neu hinzu): `cdc-capture-lag-real` 1× — dieser Slice bringt den
Wert real in die DB, liefert aber noch keinen Register-Beleg (der
Metrik-Bezug/die Auflösung folgt in slice-019). Übrige elf ohne Bezug —
siehe slice-017 §8 für die vollständige Liste der zehn ursprünglichen,
plus `dod-checkbox-nachzug` (3×, Lese-Schritt fällt der
`welle-5`-Closure zu). Kein Eintrag erreicht mit diesem Slice neu 3× —
keine Lücke.

**Modus-Begründungsblock — Umfang.** Reiner GF-Hinweis genügt (siehe oben);
kein Sub-Area-Block.
