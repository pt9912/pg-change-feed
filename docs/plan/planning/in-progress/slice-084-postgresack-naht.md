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

- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] **Falls dieser Zug die Rampe bewegt:** der Transfer-Nachweis ist in
      `harness/sensors/db-adapter-coverage.md` bzw.
      `harness/sensors/coverage-gate.md` nachgezogen — **ohne** neue
      Schwellen-ADR (`ADR-0078`). *(Kein Transfer — die Naht bleibt im Paket;
      die Bedingung ist nicht eingetreten, siehe §3.)*
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
isoliert (`git diff --name-only` listet ausschließlich Dateien unter
`internal/adapters/driven/postgresack/`); die Gegenstands-Listen sind
unberührt (`git diff` gegen `Dockerfile` und `tools/harness/db-coverage.sh` ist
leer), ebenso `.a-check.yml` und die Composition Root. Kein Paket wechselt
zwischen den zwei Gegenständen.

**(c) Kein Verhalten verloren.** Die realen, dienst-gestützten Läufe sind grün
(`make test-store` Exit 0, `make test-replication` Exit 0, `make test`
Exit 0 — je ungepiped), und **kein Testfall wird entfernt**: `ack_test.go`
ist unverändert, `git diff --name-status -- '*_test.go'` zeigt genau **eine
neue** Datei (`seam_test.go`).

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
  Wächter sind die **realen** Tests, nicht die neuen Fakes. — **Ausgang:** <bei Closure>
- **Der Fake könnte grün sein, ohne etwas zu prüfen.** Er prüft die
  **Verklebung**, nicht das Protokoll. — **Ausgang:** <bei Closure>
- **Die Naht könnte die Messung bewegen** — wie bei `slice-081`. Dann gilt
  `ADR-0078`s Transfer-Nachweis; eine Schwellen-Frage entsteht daraus **nicht**
  erneut (§3). — **Ausgang:** <bei Closure>

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
