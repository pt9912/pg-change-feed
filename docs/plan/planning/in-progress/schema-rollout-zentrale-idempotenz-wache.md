# Slice schema-rollout-zentrale-idempotenz-wache: Idempotenz-Wache zentral im `schema-rollout`-Makefile-Target

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — kein Closure-Kriterium jenseits der eigenen DoD.

**Bezug:** [ADR-0043](../../adr/0043-schemamigrationen-mit-d-migrate.md) (nur *aktive* ADR; dieser Slice ändert keine
`Accepted`-ADR-Entscheidung, sondern liefert Werkzeug-Verhalten innerhalb
der dort getroffenen Entscheidung für d-migrate).

**Berührte Spec-Stellen:** — (reine Werkzeug-/Tooling-Robustheit, kein
Lastenheft-/Pflichtenheft-/Architektur-Bezugspunkt; `make schema-rollout`
selbst trägt keine `SPEC-*`/`ARC-*`-Kennung).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** Architect (Fixrunde zu Reviewer-Fund F-1 des Architect-Verdikts
zum Lese-Schritt gegen Commit `7d0bf05`, siehe
`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/state.md`).
**Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make schema-rollout` erkennt selbst — zentral, für jeden
Aufrufer gleich —, dass ein zweiter Lauf gegen ein bereits migriertes Ziel
ausschließlich auf den bekannten Fremdobjekten blockiert (die beiden seit
`slice-016` verbliebenen Views-Reste sind bereits gelöst; real offen sind
aktuell sechs Objekte: vier SQL-Funktionen aus
`nacharbeit-administration.sql`, siehe
`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`),
und behandelt diesen Fall als bereits erfüllt statt mit Exit 8
abzubrechen — ohne echte destruktive Änderungen (Spalten-/Tabellenabbau,
neue Fremdobjekte außerhalb der bekannten Liste) ununterscheidbar
durchzulassen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Fremdobjekte ins neutrale Modell (`tools/schema/schema.yaml`)
  überführen — bleibt für die SQL-Funktionen an `POST_EXECUTE_DRIFT`
  (Exit 5, Upstream-Verhalten von d-migrate) gebunden und ist kein
  gangbarer Weg, solange dieses Verhalten besteht
  (`docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`). Anderer
  Vorgang, an ein Upstream-Fix gebunden.
- Ein generisches `--allow-destructive`-Handling, das jede Art destruktiver
  Operation durchlässt — bewusst verworfen (state.md dieser Beobachtung),
  bleibt verworfen: zu grobkörnig gegenüber dem eng umrissenen Problem.
- Rückbau der bestehenden Aufrufer-seitigen Existenz-Checks
  (`examples/bootstrap.sh`, `AGENTS.md` §3.14) — bleibt bestehen, bis
  dieser Slice tatsächlich liefert; §3.14 wird erst mit der Closure dieses
  Slice angepasst oder gestrichen, nicht vorher.
- Ein Mechanismus, der auch zukünftige, noch unbekannte Fremdobjekt-Klassen
  automatisch erkennt — dieser Slice deckt die **heute bekannte** Liste
  (sechs Objekte); ein siebtes `nacharbeit-*.sql`-Skript braucht eine
  eigene, kleine Nachtrags-Änderung an derselben Stelle (siehe §3), keine
  neue Architektur.

## 2. Definition of Done

- [x] `make schema-rollout` läuft **zweimal hintereinander** gegen dieselbe
      real migrierte Ziel-DB durch (Exit 0 bei beiden Läufen) — real
      geprüft mit dem vollen `nacharbeit-*.sql`-Bestand (aktuell: roles,
      observability, heartbeat, administration): `tools/harness/run-schema-rollout-guard-test.sh`
      Lauf 1 (frisch) und Lauf 2 (Skip-Pfad, Ausgabe enthält
      „--execute uebersprungen") beide Exit 0.
- [x] Ein **künstlich eingefügtes**, nicht auf der bekannten Liste
      stehendes destruktives Ziel bricht weiterhin mit Exit 8 ab — der
      Beleg, dass die Wache nicht pauschal durchlässt. **Plan-Korrektur:**
      Das ursprünglich skizzierte Beispiel (eine per `ALTER TABLE … DROP
      COLUMN` entfernte, von `schema.yaml` weiterhin deklarierte Spalte)
      erzeugt real **keinen** destruktiven Blocker, sondern ein
      nicht-destruktives `AddColumn` (d-migrate plant, die fehlende Spalte
      wiederherzustellen) — empirisch gegen das gepinnte Image verifiziert,
      bevor die Testform gebaut wurde. Der tatsächlich destruktive Fall ist
      die Umkehrung: eine per `ALTER TABLE … ADD COLUMN` real
      hinzugefügte, **nicht** in `schema.yaml` deklarierte Spalte erzeugt
      einen unbekannten `DropColumn`-Blocker neben den sechs bekannten —
      real geprüft (`tools/harness/run-schema-rollout-guard-test.sh` Lauf
      3, `make: *** [schema-rollout] Error 8`).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [x] `AGENTS.md` §3.14 angepasst oder gestrichen (Architect-Entscheidung
      bei Closure, siehe §1 Abgrenzung). — gestrichen (ersatzlos): der
      Negativtest (siehe oben) belegt, dass die zentrale Wache auch den
      Mischfall (bekannt + unbekannt) korrekt auflöst, keine Restlücke für
      einen dünneren Hinweis.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register
      (`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`)
      fortgeschrieben — Ausgang von `geplant` auf den tatsächlichen
      Liefer-Zustand.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Makefile` (`schema-rollout`-Target) | update | Vorlauf-Schritt vor `--execute`: `--plan-only`-Lauf gegen dasselbe Ziel (eigener Report-Pfad `tools/schema/rollout-precheck.yaml`, damit der committete Pflicht-Report `tools/schema/plan.yaml` ausschließlich echte `--execute`-Läufe belegt), Guard-Aufruf, bedingter Fallback auf den echten `--execute`-Lauf. |
| `tools/schema/rolloutguard/{main.go,report.go,guard.go,guard_test.go}` | neu | **Plan-Nachzug** (Ort jetzt entschieden): Go statt eines neuen `jq`-artigen Werkzeugs — derselbe Docker-only-Weg wie `tools/harness/*` (gepinnter Toolchain-Container, `go run`), keine neue Werkzeugkette. Liest den `--plan-only`-Report, vergleicht jede Blocker-Operation (über `kind`/`objectType`/`path`, nicht über die inhaltsabhängigen `id`-Hashes) gegen die feste Liste der sechs bekannten Fremdobjekte — direkt neben den `nacharbeit-*.sql`-Dateien (Kolokation). Unit-getestet, netzlos, läuft unter `make test`. |
| `d-migrate`-Image (`D_MIGRATE_IMAGE`) | Recherche | Durchgeführt: `docker run … schema migrate --help` zeigt **keinen** gezielten Bestätigungs-Mechanismus je Operation (`--allow-destructive` bleibt pauschal). Das gewählte Verfahren braucht ihn aber nicht — `--plan-only` allein liefert die für die Klassifikation nötige Struktur; kein partieller Rollout nötig, weil der Skip-Fall (alle Blocker bekannt) laut Definition ohnehin nichts Neues auszurollen hat. |
| `AGENTS.md` §3.14 | streichen | Ersatzlos entfernt (kein dünnerer Hinweis nötig, siehe DoD-Punkt 2 oben). |
| `tools/harness/run-schema-rollout-guard-test.sh` | neu | Realer Drei-Läufe-Beleg (frisch → Skip-Pfad → Negativ-Abbruch) gegen eine eigenständige, abgeräumte Ziel-DB — kein Unit-Test, weil DB-/d-migrate-Zugriff nötig ist; läuft im selben Docker-only-Rahmen wie `tools/harness/run-store-tests.sh`. Kein neues Make-Target (Werkzeug, `bash`-Aufruf direkt, wie in der Pre-completion-Checkliste dieses Slice-Laufs dokumentiert). |
| `.gitignore` | update | **Plan-Nachzug**: `tools/schema/rollout-precheck.yaml` ist ein transientes Vorlauf-Artefakt, kein Bestand. |

**Real gefundenes, außerhalb des Scopes liegendes Risiko:** `tools/schema/plan.yaml`/`down.sql` sind geteilte, feste Schreibziele — **jeder** Aufrufer von `make schema-rollout` (auch gegen eine völlig unabhängige Test-/Scratch-DB) überschreibt den committeten Pflicht-Report der letzten echten Produktions-/CI-Rollout mit seinem eigenen Ergebnis. Real aufgetreten: ein Lauf von `tools/harness/run-schema-rollout-guard-test.sh` hinterließ eine geänderte `tools/schema/plan.yaml` im Arbeitsbaum (zurückgesetzt, nicht committet). Vorbestehend (jeder heutige Aufrufer, u. a. `tools/schema/apply-rollout.sh`, trägt dasselbe Risiko) und nicht Gegenstand dieses Slice — siehe Beobachtungs-Register.

## 4. Trigger

**Start** (`next` → `in-progress`): sobald priorisiert; kein externer
Trigger, reine Robustheits-/Tooling-Arbeit ohne Abhängigkeit zu einer
laufenden Welle.

**Rückführungen:**

- `in-progress` → `next` (zu groß): falls die d-migrate-Recherche ergibt,
  dass ein partieller Rollout (bekannte Blocker durchlassen, echte neue
  Inhalte trotzdem anwenden) nur mit erheblichem Zusatzaufwand geht —
  dann Zerlegung in "voller Skip bei rein bekannten Blockern" (kleiner
  Liefer-Schnitt) und "partieller Rollout" (Folge-Slice).
- `in-progress` → `open` (blockiert): falls d-migrate keinerlei
  maschinenlesbaren Weg bietet, Blocker-Objekte von echten Root-Cause-Fixes
  zu unterscheiden (unwahrscheinlich — `plan.yaml` ist strukturiertes JSON
  mit `objectType`/`path`/`kind` je Operation, siehe
  `tools/schema/plan.yaml` als Referenzform eines unblockierten Laufs).

## 5. Closure-Trigger

DoD vollständig + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- d-migrate bietet keinen gezielten Bestätigungs-Mechanismus je Operation,
  nur pauschal `--execute`/kein `--execute` — **Ausgang:** weiter offen,
  klärt sich in der Recherche von §3 Zeile 3; falls zutreffend, greift die
  `in-progress` → `next`-Rückführung aus §4.
- Ein siebtes `nacharbeit-*.sql`-Skript wird angelegt, ohne die
  Kolokations-Konvention (Makefile-Zeile + Allowlist-Eintrag im selben
  Commit) einzuhalten — **Ausgang:** weiter offen, sollte über
  `commit-traceability`/Review auffallen (derselbe Aufrufer editiert
  ohnehin das Makefile für die neue Invocation-Zeile), aber kein Sensor
  erzwingt die Nähe der beiden Änderungen mechanisch.
- Die zentrale Wache maskiert einen Mischfall (ein Teil der Blocker ist
  bekannt, ein Teil ist eine echte neue destruktive Änderung) unklar —
  **Ausgang:** weiter offen, DoD-Punkt 2 oben ist der Testfall dafür; wird
  bei Umsetzung real geklärt, nicht vorab angenommen.

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt — siehe Template-Hinweis, entfällt beim
Kopieren nicht, weil dieser Plan direkt in `open/` geschrieben wurde statt
über die Vorlage kopiert; die Struktur folgt trotzdem
`.harness/baseline/v6.9.0/templates/docs/plan/planning/slice.template.md`.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „Schemamigration"
(d-migrate, Makefile-Target `schema-rollout`) — bereits mehrfach berührt
(`slice-015`, `slice-016`, `slice-036`, `slice-063`, `slice-066`,
`slice-beispiele-compose-bootstrap`), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/schema-rollout-fremdobjekte` (3×, Ausgang jetzt `geplant`, Träger
dieser Slice) und `BEO-PGC/d-migrate-nacharbeit` (6×, dritte Objektklasse
offen, unabhängig von diesem Slice) sind die Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Modus-Deklaration
`harness/conventions.md`, Default `PGC`/Greenfield).
