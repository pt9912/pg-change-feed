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
und lässt den nachfolgenden `--execute`-Schritt in diesem Fall zusätzlich
mit `--allow-destructive` laufen, statt mit Exit 8 abzubrechen — ohne
echte destruktive Änderungen (Spalten-/Tabellenabbau, neue Fremdobjekte
außerhalb der bekannten Liste) ununterscheidbar durchzulassen.
**Fixrunde innerhalb dieses Slice (noch vor dem ersten Reviewer-Zug):**
der ursprünglich gebaute Entwurf — bei ausschließlich bekannten Blockern
den `--execute`-Schritt komplett überspringen — erwies sich als falsch:
er hätte eine echte, gleichzeitig anstehende, nicht-destruktive
Schema-Änderung (z. B. eine von `schema.yaml` weiterhin deklarierte, aber
außerhalb von d-migrate entfernte Spalte) stillschweigend verloren, weil
sie im selben Plan wie die sechs bekannten Blocker steht. Korrigiert:
`--execute` läuft **immer**; `--allow-destructive` wird nur zusätzlich
gesetzt, wenn der vorgelagerte `--plan-only`-Report bestätigt, dass
ausschließlich bekannte Fremdobjekt-Blocker vorliegen — jede echte
anstehende Änderung im selben Plan wird dadurch weiterhin angewendet. Der
Regressionsbeleg dafür ist Lauf 3 von
`tools/harness/run-schema-rollout-guard-test.sh` (siehe DoD).

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
      Lauf 1 (frisch) und Lauf 2 (--allow-destructive-Pfad, Ausgabe enthält
      „--execute laeuft mit --allow-destructive") beide Exit 0.
- [x] Eine echte, gleichzeitig anstehende, nicht-destruktive Schema-Änderung
      bleibt wirksam, auch wenn sie im selben Plan wie die sechs bekannten
      Blocker steht — der Regressionsbeleg für den in §1 beschriebenen,
      innerhalb dieses Slice gefundenen und korrigierten Fehler: eine von
      `schema.yaml` weiterhin deklarierte, nullable Spalte
      (`administration_request.error_message`) wird außerhalb von
      d-migrate per `ALTER TABLE … DROP COLUMN` entfernt; der nächste
      `make schema-rollout`-Lauf bringt sie real zurück — real geprüft
      (`tools/harness/run-schema-rollout-guard-test.sh` Lauf 3).
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
      4, `make: *** [schema-rollout] Error 8`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review —
      `docs/reviews/review-slice-schema-rollout-zentrale-idempotenz-wache.md`,
      2 HIGH (F-1, F-2) in der Fixrunde behoben, 1 MEDIUM (F-3) dokumentiert
      (kein Code-Fix möglich, siehe §6), 2 LOW (F-5, F-6) behoben/notiert.
- [x] `AGENTS.md` §3.14 angepasst oder gestrichen (Architect-Entscheidung
      bei Closure, siehe §1 Abgrenzung). **Fixrunde (Reviewer-Fund F-2):**
      die ursprüngliche ersatzlose Streichung verletzte `AGENTS.md` §3.13 —
      [`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
      §Teilfrage 4 (`Accepted`, unberührbar per §3.5) zitiert „`AGENTS.md`
      §3.14" namentlich als Analogie und wurde beim Nachzugs-Suchlauf
      übersehen. Da die Verglichene-Alternativen-Sektion einer `Accepted`
      ADR nicht editierbar ist (§3.5), trägt §3.14 jetzt einen reinen
      Rang-Zeiger-Absatz (keine Regel), der den Verweis auflösbar hält,
      statt die Nummer ersatzlos verwaist zu lassen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register
      (`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`)
      fortgeschrieben — Ausgang von `geplant` auf `verkörpert`.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen
      (Prüfung läuft nach dem `git mv` nach `done/`, siehe Template-Hinweis §7).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Makefile` (`schema-rollout`-Target) | update | Vorlauf-Schritt vor `--execute`: `--plan-only`-Lauf gegen dasselbe Ziel (eigener Report-Pfad `tools/schema/rollout-precheck.yaml`, damit der committete Pflicht-Report `tools/schema/plan.yaml` ausschließlich echte `--execute`-Läufe belegt), Guard-Aufruf, bedingter Fallback auf den echten `--execute`-Lauf. |
| `tools/schema/rolloutguard/{main.go,report.go,guard.go,guard_test.go}` | neu | **Plan-Nachzug** (Ort jetzt entschieden): Go statt eines neuen `jq`-artigen Werkzeugs — derselbe Docker-only-Weg wie `tools/harness/*` (gepinnter Toolchain-Container, `go run`), keine neue Werkzeugkette. Liest den `--plan-only`-Report, vergleicht jede Blocker-Operation (über `kind`/`objectType`/`path`, nicht über die inhaltsabhängigen `id`-Hashes) gegen die feste Liste der sechs bekannten Fremdobjekte — direkt neben den `nacharbeit-*.sql`-Dateien (Kolokation). Unit-getestet, netzlos, läuft unter `make test`. |
| `d-migrate`-Image (`D_MIGRATE_IMAGE`) | Recherche | Durchgeführt: `docker run … schema migrate --help` zeigt **keinen** gezielten Bestätigungs-Mechanismus je Operation (`--allow-destructive` bleibt pauschal). Das gewählte Verfahren braucht ihn aber nicht — `--plan-only` allein liefert die für die Klassifikation nötige Struktur; kein partieller Rollout nötig, weil der Skip-Fall (alle Blocker bekannt) laut Definition ohnehin nichts Neues auszurollen hat. |
| `AGENTS.md` §3.14 | streichen | Ursprünglich ersatzlos entfernt; in der Reviewer-Fixrunde (F-2) auf einen Rang-Zeiger-Absatz korrigiert, siehe DoD-Punkt „AGENTS.md §3.14" und die Fixrunden-Zeile unten. |
| `tools/harness/run-schema-rollout-guard-test.sh` | neu | Realer Drei-Läufe-Beleg (frisch → Skip-Pfad → Negativ-Abbruch) gegen eine eigenständige, abgeräumte Ziel-DB — kein Unit-Test, weil DB-/d-migrate-Zugriff nötig ist; läuft im selben Docker-only-Rahmen wie `tools/harness/run-store-tests.sh`. Kein neues Make-Target (Werkzeug, `bash`-Aufruf direkt, wie in der Pre-completion-Checkliste dieses Slice-Laufs dokumentiert). |
| `.gitignore` | update | **Plan-Nachzug**: `tools/schema/rollout-precheck.yaml` ist ein transientes Vorlauf-Artefakt, kein Bestand. |
| `Makefile`, `tools/schema/rolloutguard/{guard.go,main.go,report.go,guard_test.go}` | Fixrunde (Plan-Nachzug) | Der Mischfall aus §6 (dritter Risikopunkt) trat real ein, noch innerhalb dieses Slice, vor jedem Reviewer-Zug: der erste Entwurf überspringt `--execute` komplett bei ausschließlich bekannten Blockern und hätte dadurch eine echte, gleichzeitig anstehende, nicht-destruktive Änderung verloren (empirisch verifiziert). Korrektur: `--execute` läuft immer; `decide()`/`main.go` liefern jetzt `allowDestructive` statt `skip`, das Makefile-Target setzt `--allow-destructive` nur zusätzlich, wenn `decide()` es erlaubt. `guard_test.go`-Namen/-Kommentare entsprechend nachgezogen; kein Verhaltensunterschied für die fünf bestehenden Unit-Tests (reine Allowlist-Prüfung unverändert), aber neue Semantik der Rückgabe. Regressionsbeleg: Lauf 3 von `tools/harness/run-schema-rollout-guard-test.sh`. |
| `Makefile`, `tools/schema/rolloutguard/{guard.go,report.go}`, `AGENTS.md` | Fixrunde (Reviewer F-1, F-2, F-5) | F-1 (HIGH): Vorher/Nachher-Chronik-Prosa in beiden Kommentaren („statt beim bloßen Überspringen von --execute…", „ein reines Überspringen … ließ … nicht zurückkommen") auf reinen Ist-Zustand mit Anker auf den Regressionstest umgeschrieben (`BEO-PGC/slice-chronik-in-code-kommentar`, 5./6. Auftreten). F-2 (HIGH): `AGENTS.md` §3.14 trägt jetzt einen Rang-Zeiger-Absatz statt leer zu bleiben, weil [`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md) §Teilfrage 4 (`Accepted`) die Nummer namentlich zitiert — die ADR selbst bleibt unangetastet (§3.5, Verglichene-Alternativen-Sektion ist kein Zitat-Korrektur-Fall). F-5 (LOW): Grenze-Kommentar zu `foreignObject`s Schlüssel-Kollisionsannahme (Single-Schema, punktfreie Pfadsegmente) ergänzt. |

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
  nur pauschal `--execute`/kein `--execute` — **Ausgang: eingetreten**
  (Recherche in §3 Zeile 3 bestätigt: kein solcher Mechanismus existiert).
  Konsequenz, real gefunden im Review (F-3, MEDIUM): `--plan-only` und
  `--execute --allow-destructive` sind zwei unabhängige, sequenzielle
  Docker-Läufe gegen denselben lebenden Ziel-Zustand — ein zwischen beiden
  Läufen neu entstehender destruktiver Blocker würde vom Precheck nicht
  erfasst, liefe aber unter dem bereits gesetzten `--allow-destructive`
  durch (enges, aber reales Fenster, dokumentiert als Grenze-Kommentar im
  Makefile-Target). Kein Root-Cause-Fix möglich ohne ein d-migrate-Flag,
  das einen geprüften Plan zur Ausführung wieder einliest — dieses Flag
  existiert laut Recherche nicht; die `in-progress` → `next`-Rückführung
  aus §4 griff hier bewusst nicht, weil das Risiko klein und für den
  Single-Operator-/CI-Nutzungsfall dieses Repos hinnehmbar ist.
- Ein siebtes `nacharbeit-*.sql`-Skript wird angelegt, ohne die
  Kolokations-Konvention (Makefile-Zeile + Allowlist-Eintrag im selben
  Commit) einzuhalten — **Ausgang:** weiter offen, sollte über
  `commit-traceability`/Review auffallen (derselbe Aufrufer editiert
  ohnehin das Makefile für die neue Invocation-Zeile), aber kein Sensor
  erzwingt die Nähe der beiden Änderungen mechanisch.
- Die zentrale Wache maskiert einen Mischfall (bekannte Blocker neben
  einer echten anstehenden Änderung) unklar — **Ausgang: eingetreten**,
  noch innerhalb dieses Slice, vor dem ersten Reviewer-Zug. Der erste
  Entwurf (Skip von `--execute` bei ausschließlich bekannten Blockern)
  hätte eine echte, nicht-destruktive Änderung im selben Plan
  stillschweigend verloren; korrigiert auf „`--execute` läuft immer,
  `--allow-destructive` nur zusätzlich bei bestätigt bekannten Blockern"
  (siehe §1, §3 Fixrunden-Zeile). Der Fund kam aus einer bewussten
  Gegenprobe (Mischfall-Simulation), nicht aus einem Sensor — kein neuer
  Sensor daraus, der Regressionsbeleg ist Lauf 3 des Test-Skripts.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die im Beobachtungs-Register skizzierte dritte
  Option (zentrale, blocker-klassifizierende Wache statt Aufrufer-seitiger
  Existenz-Checks) war real umsetzbar — eine kleine, unit-getestete
  Go-Komponente (`tools/schema/rolloutguard`) neben den bestehenden
  `nacharbeit-*.sql`-Dateien, keine neue Werkzeugkette. Die vorab
  durchgeführte empirische Verifikation gegen das gepinnte d-migrate-Image
  (welcher Blocker real entsteht, bevor die Testform gebaut wurde)
  verhinderte einen falsch konstruierten DoD-Negativtest.
- **Was ging anders als geplant:** Der erste Entwurf (`--execute`
  komplett überspringen bei ausschließlich bekannten Blockern) erwies
  sich als korrektheitsverletzend — noch innerhalb dieses Slice per
  Gegenprobe gefunden und korrigiert, vor jedem Reviewer-Zug (§1, §3).
  Der Reviewer fand danach 2 HIGH (F-1 Chronik-Sprache in zwei
  Kommentaren, F-2 eine durch die §3.14-Streichung verwaiste
  Analogie-Referenz in der `Accepted`-`ADR-0100`) — beide in einer
  zweiten Fixrunde behoben, ohne die ADR selbst zu berühren (§3.5-konform:
  ein Rang-Zeiger-Absatz in `AGENTS.md` §3.14 statt eines ADR-Eingriffs).
  Der Verifier fand danach eine dritte, außerhalb des Review-Diffs
  liegende Instanz derselben Chronik-Sprache in `harness/README.md` —
  Doku-Prosa statt Produktionscode, außerhalb des wörtlichen Skopus des
  Reviewer-HIGH-Punkts (siehe Steering-Loop-Eintrag) — ebenfalls vor
  Closure behoben.
- **Steering-Loop-Eintrag:** Der Reviewer-Skill-HIGH-Punkt
  „Slice-/Wellen-Chronik in Produktionscode-Kommentar"
  (`BEO-PGC/slice-chronik-in-code-kommentar`, jetzt 7× belegt) fängt
  weiterhin jede Produktionscode-Instanz vor Merge — bestätigt, nicht
  geschärft. Neu benannt (Randbefund, noch nicht verkörpert, da 1×):
  dieselbe Chronik-Klasse tritt auch in Doku-Prosa (`harness/README.md`)
  auf, außerhalb des wörtlichen HIGH-Punkt-Skopus „Produktionscode-
  Kommentar" — dort fing sie der Verifier, nicht der Reviewer. Auslöser:
  `BEO-PGC/slice-chronik-in-code-kommentar`
  (`evidence/slice-schema-rollout-zentrale-idempotenz-wache.md`).
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-schema-rollout-zentrale-idempotenz-wache.md` in
  `BEO-PGC/slice-chronik-in-code-kommentar/` ergänzt (7. Beleg plus
  Doku-Randbefund). `BEO-PGC/schema-rollout-fremdobjekte` von `geplant`
  auf **verkörpert** fortgeschrieben — Träger ist dieser Slice
  (`Makefile`, `tools/schema/rolloutguard/`).
- **Folge-Slices:** keine — kein Punkt aus §1 „Ausdrücklich NICHT in
  diesem Slice" benennt eine konkrete Folge-Slice-Kennung; die
  Fremdobjekt-Überführung ins neutrale Modell bleibt an das
  `POST_EXECUTE_DRIFT`-Upstream-Verhalten gebunden
  (`BEO-PGC/d-migrate-nacharbeit`, unverändert offen, nicht Gegenstand
  dieses Slice).
- **Risiken aus §6:** drei Risiken, drei Ausgänge — (1) „kein gezielter
  Bestätigungs-Mechanismus" → **eingetreten**, Konsequenz ist das
  TOCTOU-Fenster aus Review-Fund F-3 (MEDIUM), als Grenze-Kommentar im
  Makefile-Target dokumentiert, kein Root-Cause-Fix möglich, Risiko im
  Nutzungskontext dieses Repos hinnehmbar; (2) „Kolokations-Konvention
  eines siebten `nacharbeit-*.sql`-Skripts" → **weiter offen**, kein
  Sensor erzwingt die Nähe der beiden Änderungen mechanisch, bleibt
  Review-/`commit-traceability`-Disziplin; (3) „Mischfall-Maskierung" →
  **eingetreten**, noch innerhalb dieses Slice gefunden und korrigiert
  (siehe „Was ging anders als geplant").
- **Drei Paarungen:** Anker — der Steering-Loop-Eintrag trägt keinen
  `liegt in`-Feld (kein neuer Sensor/keine neue Regel verkörpert, nur ein
  bestehender HIGH-Punkt bestätigt plus ein Randbefund benannt), also
  kein Anker zu prüfen. Folge-Slice — keine benannt (siehe oben), nichts
  zu prüfen. Register — `BEO-PGC/schema-rollout-fremdobjekte` und
  `BEO-PGC/slice-chronik-in-code-kommentar` existieren beide als
  Verzeichnis mit nicht-leerem `evidence/`; `BEO-PGC/d-migrate-nacharbeit`
  (§1 referenziert) existiert ebenfalls mit nicht-leerem `evidence/`. Alle
  drei tragen.

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
