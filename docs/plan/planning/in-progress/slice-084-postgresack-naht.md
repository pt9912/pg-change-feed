# Slice slice-084: postgresack-Naht — schmale Abhängigkeit statt konkreter pgconn-Verbindung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (eine
Struktur-Änderung mit eigenem Beleg), kein repo-weites *Mehr*.

**Bezug:** [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
(die **Nahtform** für die `*pgconn.PgConn`-förmigen Pakete: Treiber-Hülle plus
schmale, fake-fähige Fläche, **im Paket**) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 (schmale Abhängigkeit statt konkreter Typ, Fehlerklassifikation als
reine Funktion) · [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(der Nachweis, den diese Naht als **Null-Befund** führt) ·
`slice-081` (die Geschwister-Naht im pool-förmigen Adapter; aus ihrer Abweichung
ist dieser Vorgang entstanden).

**Berührte Spec-Stellen:** — (die Naht liegt **innerhalb** des driven Adapters,
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

**Ziel:** `internal/adapters/driven/postgresack` hängt an einer **konkreten**
`*pgconn.PgConn` — gemessene Fläche (**6** Symbole): `pgconn.ConnectConfig`,
`pgconn.ParseConfig`, `pgconn.PgConn` sowie `pglogrepl.LSN`,
`pglogrepl.SendStandbyStatusUpdate`, `pglogrepl.StandbyStatusUpdate`. Die Naht
macht seine Verklebung netzlos prüfbar.

**Die Nahtform ist entschieden** ([`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)):
eine **adapter-eigene Treiber-Hülle** plus eine **schmale, fake-fähige Fläche**
— **im Paket**, ohne Paketwechsel. Die Schnittstelle trägt **eine** Methode
(`SendStandbyStatusUpdate(ctx, pglogrepl.StandbyStatusUpdate) error`); die Hülle
(`connSender{conn *pgconn.PgConn}`) hält den konkreten Typ. `New(conn, opts…)`
samt nil-Grenze und Composition-Root-Verdrahtung bleiben **unverändert**, dazu
kommt ein paket-interner Einstieg für die netzlosen Tests.

**Gemessen, nicht behauptet:** `pglogrepl.SendStandbyStatusUpdate` verlangt
`conn *pgconn.PgConn` **in der Signatur** — der Typ erfüllt es also nicht; und
die einzige direkt erfüllbare `Exec`-Form (`*pgconn.MultiResultReader`) ist
**nicht fake-bar** (kein exportierter Konstruktor). „Fake kann den Treiber
spielen" und „der konkrete Typ bleibt draußen" sind deshalb **nur über eine
Hülle** zugleich zu haben (positiv belegt: Hülle und Fake erfüllen dieselbe
Schnittstelle). **`slice-081`s Zusicherungsform trägt hier nicht** — sie setzte
voraus, dass der reale Träger die Schnittstelle **selbst** erfüllt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Design-Vorgriff.** Die Naht-Form ist eine Architect-Entscheidung
  (§1); dieser Plan nennt den Gegenstand, nicht den Schnitt.
- **`receive`** — eigene Naht, eigene Datei: `slice-085`. Dort ist die Fläche
  eine andere Größenordnung (18 `pgconn`/`pglogrepl`-Symbole gegen 6 hier).
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
      (eine Methode); der **konkrete** `*pgconn.PgConn` bleibt in der Hülle und
      im Dial. **Nicht** „statt der konkreten Verbindung": der Typ bleibt, aber
      **hinter** der Naht (`ADR-0080`).
- [x] **Kein Verhaltens-Change:** die reale Verdrahtung geht unverändert durch
      `make test-replication` (Exit 0).

**Liefer-Punkt 2 — die reine Logik ist prüfbar.**

- [x] Der netzlos prüfbare Teil (Aufbau der Standby-Status-Meldung, LSN-Form)
      liegt als **reine Funktion** mit eigenen Tests vor.
- [x] Die Fake-Seite fährt die **Verklebung** — und ist ausdrücklich **kein**
      Ersatz der realen Tests.

**Liefer-Punkt 3 — der Nachweis ist geführt, und er ist ein Null-Befund.**

- [x] Der Nachweis aus [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
      wird als **Null-Befund** geführt (`ADR-0080`): die Naht bleibt **im**
      Paket, es gibt **keinen Subjekt-Transfer** — `k_ab = 0`, der abfließende
      Nenner ist unberührt, und es gibt **keine Neu-Bemessung**; die Rampen
      (70 / 70, Endstufen 80 %) bleiben unverändert. Nachzuweisen sind die drei
      Teile dennoch: **(a)** die Arithmetik zeigt `k_ab = 0`; **(b)** der
      **Paket-Diff** zeigt **keinen** Trägerwechsel; **(c)** kein Verhalten
      verloren. **Benannte Grenze:** die neuen netzlosen Tests verdünnen den
      DB-Nenner leicht — mit Trigger, nicht still.
- [x] `make gates` grün (Exit direkt, ungepiped).

- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-084.md`, `.harness/skills/reviewer.md`) —
      Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6),
      kein Self-Review (Modul 8). **Fixrunde gelaufen:** 0 HIGH, 0 MEDIUM;
      F-1 und F-3 geschlossen, F-2 vom Planner (`a3cb2c4`), F-4 kein Fix.
- [x] **Falls dieser Zug die Rampe bewegt:** der Transfer-Nachweis ist in
      `harness/sensors/db-adapter-coverage.md` bzw.
      `harness/sensors/coverage-gate.md` nachgezogen — **ohne** neue
      Schwellen-ADR (`ADR-0078`). *(Kein Transfer — die Naht bleibt im Paket;
      die Bedingung ist nicht eingetreten, siehe §3.)*
- [x] Verifikation durchgeführt, Report unter
      `docs/reviews/verify-slice-084.md` liegt vor (Modul 11, frischer Kontext) —
      **closure-fähig**; alle Liefer-Kriterien bestätigt. Der Verifier hat den
      **Null-Befund am Parent-Commit** nachgebaut (Wegwerf-Worktree auf
      `fb6adf6`: derselbe Unit-Nenner **1903**, **null**
      `postgresack`/`receive`-Zeilen in beiden Profilen), die Zahlenträger
      einzeln gegen seine eigene Messung gehalten, die neue Beobachtung
      bestätigt (`postgresack` **30/32** netzlos) und zwei eigene Mutationen rot
      gesehen (`ackLSN` Offset+1, doppelter Fehler-Log). Seine **V-1** (mein
      Beleg (b) trug ohne Pathspec nicht) ist berichtigt und **nachgeprüft**.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) — **entfällt**: dieses
      Repo führt die Datei nicht (Greenfield-Bootstrap, kein Inventur-Fund).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — ein Beleg
      ergänzt (`zahl-in-traeger-driftet-gegen-die-messung`, damit **2×**) und
      **ein Verzeichnis neu angelegt**
      (`db-gegenstand-enthaelt-netzlos-geprueften-code`, **1×**), **kein Zähler
      gesetzt**.
- [x] Jedes Risiko aus §6 trägt einen Ausgang — alle drei *entfallen,
      gestrichen mit Begründung*.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) — dieses Repo führt
      **Wellen-Betrieb**; die Prüfung fällt der `welle-20`-Closure zu, hier nicht
      geprüft und hier nicht fällig.

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
| `internal/adapters/driven/postgresack/seam.go` | neu | die Naht (`ADR-0080`): `standbySender` (eine Methode), die Treiber-Hülle `connSender` über `*pgconn.PgConn`, die Kompilier-Zusicherung `var _ standbySender = connSender{}` |
| `internal/adapters/driven/postgresack/ack.go` | refactor | die Logik hängt an der Naht; `ackLSN` (Null-Positions-Grenze, LSN-Form), `standbyStatus` (Standby-Status-Form) und `replicationClass` (Fehlerklassen-Wrapping) sind reine Funktionen; `New` behält seinen Signatur-Vertrag und reicht die Hülle durch, `newOnSender` ist der **paket-interne** Einstieg der netzlosen Tests |
| `internal/adapters/driven/postgresack/seam_test.go` | neu | Fake und Log-Träger; die Verklebung netzlos: abgesetzte Meldung, Null-Positions-Grenze ohne Absetzen, Fehlerpfad der Naht — der Fake erfüllt dieselbe Schnittstelle wie die Hülle |
| `harness/image-hash.txt` | update | der Zug ändert Build-Kontext-Dateien; `make image` stempelt den Digest des Laufs (`ADR-0044`) — `sha256:f4475e1e…` → `sha256:4bd43435…` |

**Kein Paketwechsel, kein Unterpaket** (`ADR-0080`) — deshalb auch **keine**
Änderung am `Dockerfile`-Filter oder an `DB_COVERAGE_PKGS`: der DB-Gegenstand
bleibt, wie er ist.

**Nicht angefasst:** `internal/bootstrap/wiring.go` (die Composition Root ruft
`postgresack.New(stream.Conn(), …)` unverändert), `Dockerfile` (Stufe
`coverage`, Paket-Filter), `tools/harness/db-coverage.sh` (`DB_COVERAGE_PKGS`),
`harness/mk/coverage.mk` (`THRESHOLD`), `harness/sensors/**`, `.a-check.yml`,
`spec/**`, `internal/adapters/driving/replication/receive/**` (→ `slice-085`).
`ack_test.go` (die realen Tests) ist **unverändert** — kein Testfall entfernt.

**Nicht in dieser Liste:** `internal/adapters/driving/replication/receive/**`
(→ `slice-085`), `spec/**`, `.a-check.yml`.

**Der Null-Befund — alle drei Teile, mit den Zahlen dieses Laufs**
([`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
§Entscheidung 4, [`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
§Entscheidung 1). Gemessen in den gepinnten Images, Exit-Codes ungepiped; die
Zahlen je Gegenstand aus den beiden `-coverprofile` des Laufs nach der
Dedup-Regel aus `tools/harness/db-coverage.sh`:

| Gegenstand (Statements) | vor dem Zug | nach dem Zug | Δ |
|---|---|---|---|
| netzlos prüfbare Fläche (Unit) | 1903 (1369 gedeckt, `coverage-gate` 71,90 %) | 1903 (1371 gedeckt, `coverage-gate` 72,00 %) | **0** |
| DB-Adapter-Gegenstand | 650 (477 gedeckt, 73,38 %) | **659** (491 gedeckt, 74,51 %) | **+9** |
| davon `postgresack` | 23 (18 gedeckt) | 32 (32 gedeckt) | **+9** |
| davon `postgresstorage` | 472 | 472 | 0 |
| davon `replication/receive` | 155 | 155 | 0 |

Die **gedeckte** Zahl der Unit-Fläche bewegt sich lauf-zu-lauf (`internal/bootstrap`
262 → 264, beobachtet ±2 — `slice-081` §3 beschreibt dieselbe Schwankung);
Nenner (1903) und jedes Paket-Statement sind davon unberührt.

**(a) `k_ab = 0`.** Kein Gegenstand gibt Statements ab: der Unit-Nenner steht
unverändert bei 1903 (die Naht liegt **im** ausgenommenen Paket), und innerhalb
des DB-Gegenstands sind `postgresstorage` (472) und `replication/receive` (155)
byte-stabil. Der DB-Nenner **wächst** um **+9** — die neuen Statements des
paket-internen Einstiegs, der reinen Funktionen und der Hülle; `k_ab = 0` ist
damit die Null-Hälfte der Arithmetik, nicht eine Behauptung.

**(b) Paket-Diff — kein Trägerwechsel.** Die Änderung ist auf **ein** Paket
isoliert, und der Beleg braucht **beides — Range *und* Pathspec**:
`git diff --name-only fb6adf6..4035ee7 -- internal/adapters/driven/postgresack/`
listet die drei geänderten Dateien; und die **Gegenrichtung** trägt den Rest:
`git diff --name-only fb6adf6..4035ee7 -- internal/ ':!internal/adapters/driven/postgresack/'`
ist **leer**. Die Gegenstands-Listen sind unberührt
(`git diff fb6adf6..4035ee7 -- Dockerfile tools/harness/db-coverage.sh
harness/mk/coverage.mk` ist leer), ebenso `.a-check.yml` und die Composition
Root. Kein Paket wechselt zwischen den zwei Gegenständen.
**Beides gehört in den Beleg** (Review `review-slice-084` F-2, Verifikation
desselben Slice **V-1**): ein `git diff` **ohne Range** ist auf sauberem Baum
leer, und eines **ohne Pathspec** listet mehr als das Behauptete — die erste
Fassung dieses Belegs trug aus beiden Gründen nicht.

**(c) Kein Verhalten verloren.** Die realen, dienst-gestützten Läufe sind grün
(`make test-store` Exit 0, `make test-replication` Exit 0, `make test`
Exit 0 — je ungepiped), und **kein Testfall wird entfernt**: `ack_test.go`
ist unverändert, `git diff --name-status fb6adf6..4035ee7 -- '*_test.go'` zeigt
genau **eine neue** Datei (`seam_test.go`).

**Rot-Gegenprobe an der Zusage — einmal gesehen.** Auf einer Wegwerf-Kopie
brachen zwei Mutationen je ihre Prüfung: `WALApplyPosition: 0` in
`standbyStatus` färbt `TestStandbyStatusCarriesPositionInAllThreeLSNs` und
`TestAcknowledgeSendsStandbyStatus` rot (`go test` Exit 1 — die Form-Zusage),
und `replicationClass` ohne die Wrappung (`return cause`) färbt
`TestReplicationClassWrapsCause` und `TestAcknowledgeWrapsSenderFailure` rot
(die Fehlerklassen-Zusage).

**Benannte Grenze:** die neuen netzlosen Tests laufen im Messlauf der
DB-Adapter-Coverage mit und decken dort Statements, die keine PostgreSQL-Instanz
berührt hat — `postgresack` steht deshalb bei 32 von 32 gedeckten Statements.
Die Zahl heißt weiterhin richtig „Coverage des DB-Gegenstands"; die Verdünnung
ist benannt und hat den Trigger aus `ADR-0080`
(§„Die benannte Grenze"). **Keine Rampe bewegt:** `DB_COVERAGE_THRESHOLD`
bleibt 70, `THRESHOLD` bleibt 70, die Endstufen bleiben 80 % — eine
Schwellen-ADR wird nicht fällig.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): das Naht-Design ist architect-entschieden,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich die Naht
  als breiter als **ein** Paket, gehört sie zurück zur Zerlegung — **nach
  Aufruf**, nicht als Sammelumbau.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass die Naht ohne
  **Verhaltensänderung** nicht zu haben ist, ist das ein Blocker mit
  Entscheidung.

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

- **Der Umbau könnte ein Verhaltens-Change sein, der als Refactoring auftritt.**
  Wächter sind die **realen** Tests, nicht die neuen Fakes. — **Ausgang:**
  *entfallen — gestrichen mit Begründung*: `make test-store` und
  `make test-replication` sind grün, `ack_test.go` ist **byte-identisch**, kein
  `t.Skip`, kein `|| true`, keine Zusicherung abgeschwächt. Und der Wächter hat
  **gebissen**: die Mutation „nicht delegierende Hülle" bleibt netzlos grün und
  wird vom **realen Tier** gefangen (`TestAcknowledgeOnClosedConnection`,
  Exit 1) — die Wächter-Rolle liegt genau dort, wo der Plan sie verortet.
- **Der Fake könnte grün sein, ohne etwas zu prüfen.** Er prüft die
  **Verklebung**, nicht das Protokoll. — **Ausgang:** *entfallen — gestrichen
  mit Begründung*: **sechs** Mutationen wurden rot gesehen — zwei des
  Implementers (`WALApplyPosition: 0`, `replicationClass` ohne Wrapping) und
  **vier** eigene des Reviewers (nil-Grenze, Aufrufzähler, `New`-nil, die nicht
  delegierende Hülle). Die letzte zeigt zusätzlich die **Arbeitsteilung**: was
  der Fake nicht fangen kann, fängt der reale Lauf.
- **Die Naht könnte die Messung bewegen** — wie bei `slice-081`. Dann gilt
  `ADR-0078`s Transfer-Nachweis; eine Schwellen-Frage entsteht daraus **nicht**
  erneut (§3). — **Ausgang:** *entfallen — gestrichen mit Begründung*: der
  **Null-Befund** ist geführt (`k_ab = 0`, Unit-Nenner 1903 → 1903, Paket-Diff
  nur `postgresack/**`) — **kein** Subjekt-Transfer, **keine** Neu-Bemessung.
  **Aber die Zahl hat sich bewegt** (73,38 % → 74,51 %, der Nenner **wuchs** um
  9) — das ist **keine** Schwellen-Frage, sondern eine **neue Beobachtung**:
  der DB-Gegenstand enthält jetzt netzlos geprüften Code
  (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`, **neu**).

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

- **Was hat funktioniert:** **Die Nahtform aus [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
  hat getragen** — und der Slice hat sie nicht behauptet, sondern **ausgeführt**:
  der `Null-Befund` ist geführt (`k_ab = 0`, Unit-Nenner 1903 → 1903, Paket-Diff
  nur `postgresack/**`), nicht zugesagt.
  Und **die Arbeitsteilung der Wächter** ist sichtbar geworden: die Mutation
  „nicht delegierende Hülle" bleibt **netzlos grün** und wird erst vom **realen
  Tier** gefangen (`TestAcknowledgeOnClosedConnection`, Exit 1). Das ist der
  Beweis, dass „der Fake prüft die Verklebung" nicht „der Fake ersetzt die
  realen Tests" heißt — genau der Satz, den §1 als Ausschluss führt.
- **Was ging anders als geplant:** (1) Die Fixrunde musste **vier** driftende
  Zahlen in **zwei** Sensor-Dokumenten nachziehen; **zwei** davon stammten aus
  **anderen** Vorgängen und wurden mitgezogen, weil sie sonst in derselben Datei
  gegen die eigene Messung stünden. (2) Daraus entstand eine **Form**, die der
  Slice vorher nicht hatte (§3): der **Nenner** ist Zustand, die **gedeckte
  Zahl** ist Beleg eines konkreten Laufs. (3) **Der wichtigste Befund liegt
  neben dem Slice:** von den **32** Statements des Pakets sind **30** netzlos
  gedeckt — der DB-Gegenstand enthält jetzt Code, der keine Verbindung braucht.
  Das ist das **Spiegelbild** zu `slice-081` (dort *dräniert* ein Transfer den
  Nenner, hier *wächst* er) und eine **neue** Register-Klasse. (4) **Mein
  Fehler, zweimal derselbe:** mein Auftrag hat dem Implementer den Slice-Plan
  entzogen und damit die DoD-Box „Review durchgeführt" mitbetroffen, die
  `implement-slice` Schritt 21 der Fixrunde zuweist. Er hat es **beide Male
  gemeldet statt still abzuweichen**; ab jetzt schreibe ich die Randbedingung
  enger.
- **Lerneintrag (geschärfte Regel):** *Eine Messung, deren Zahl sich bewegt,
  ohne dass sich ihr Gegenstand bewegt, ist nur bedingt lesbar — und die
  Bewegung kann in **beide** Richtungen irreführen.* `slice-081` hat das nach
  **unten** gezeigt (Transfer, 75,25 % → 73,38 %), dieser Slice nach **oben**
  (Netzlos-Code, 73,38 % → 74,51 %). `ADR-0078` hält für die Abwärtsbewegung
  einen Riegel bereit; für die Aufwärtsbewegung gibt es keinen — sie ist als
  Beobachtung geführt.
  **Dazu die Form-Regel für Zahlenträger:** der **Nenner** steht als Zustand,
  die **gedeckte Zahl** als **Beleg eines konkreten Laufs** — sie nennt ihren
  Lauf, nie „der Ist-Stand". Beide Regeln sind **Verkörperungs-Kandidaten** (die
  zweite ist die *Benennung* der Herkunfts-Regel, die `ADR-0078` formuliert und
  niemand gebaut hat); geführt, nicht behauptet.
- **Steering-Loop-Eintrag:** **keine Verkörperung durch diesen Slice.** Zwei
  Registerbewegungen: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
  +`slice-084` → **2×**; **neu angelegt**
  `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (**1×**).
- **Beobachtungs-Register (`../observations/`):** zwei Belege ergänzt
  (`zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-084.md`) und **ein
  Verzeichnis neu angelegt**
  (`db-gegenstand-enthaelt-netzlos-geprueften-code/`, mit
  `evidence/slice-084.md`), **kein Zähler gesetzt**.
- **Mitgaben der Verifikation (INFO, nicht behoben):** **V-2** — eine
  zweideutige Klammer in `harness/sensors/db-adapter-coverage.md` („die geltende
  Größe des Gegenstands: §Zählbasis") direkt neben dem eingefrorenen „477 von
  650". **V-3** — zwei Rot-/Grün-Belege ebenda **ohne Lauf-Benennung**, während
  die von diesem Slice eingeführte Form sie verlangt. **Beide benannt, nicht
  behoben**: sie sind INFO und nicht blockierend, und die zwei Belege sind
  ausdrücklich **eingefrorene Einzelläufe**. Sie durchzusetzen heißt, die
  Beleg-Zeilen anzufassen — ein eigener kleiner Zug, kein Anhang dieser Closure.
- **Folge-Slices:** `slice-085` (die Geschwister-Naht in `replication/receive`)
  — liegt als Datei in `open/`.
- **Risiken aus §6:** alle drei *entfallen, gestrichen mit Begründung* — siehe §6.
- **Drei Paarungen:** nicht hier — dieses Repo führt **Wellen-Betrieb**, die
  Prüfung fällt der `welle-20`-Closure zu (Modul 6 Schritt 3c, auch für Slices
  ohne Wellen-Zugehörigkeit).
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
  die Klasse führt die Adapter-Grenze „kontrollierte Fortsetzung beim Aufrufer";
  dieser Slice ändert, **woran** der Adapter hängt, nicht was der Aufrufer tut.
- `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (2×, offen):
  **kein Treffer** — sie betrifft ungeprüft übernommene Tatsachenbehauptungen in
  Begründungen; die Zahlen dieses Plans sind am Gegenstand gemessen bzw. als
  gemessen zitiert.

**Ergebnis** (Stand: Anlage dieses Plans — kein Versprechen über den Lauf):
kein Eintrag steht über der Schwelle, und keiner rückt mit diesem Slice auf 3×.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Schichten-Edges prüft `.a-check.yml`, die
Ports-und-Adapter-Ordnung führt `spec/architecture.md`), **Phase-Reife** hoch,
**Evidenz-/Diskrepanz-Risiko** **mittel** (der Umbau kann Verhalten verschieben,
und genau dagegen stehen die realen Tests als Wächter), **Reconciliation-Aufwand**
null.
