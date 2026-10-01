# Slice routing-e2e: E2E-Belege — Routing am laufenden Feed-Container, auf allen Lesewegen, über PostgreSQL 17 und 18

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (ausgeschlossene
Spaltenwerte erscheinen nicht in den Changes),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbare Fehlerzustände),
[`LH-FA-REA-005`](../../../../spec/lastenheft.md) (erneutes Lesen),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Gleichwertigkeit der Zugriffswege),
[`LH-QA-POR-001`](../../../../spec/lastenheft.md) (PostgreSQL 17 und 18),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 7 (E2E-Belege — Träger dieses Slice) und Entscheidung 4 (DELETE),
[`ADR-0030`](../../adr/0030-testpyramide.md) (E2E-Tier),
[`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md) (Testansatz, zwei
PostgreSQL-Versionen).

**Berührte Spec-Stellen:** — (kein SPEC-/ARC-Eintrag; der Slice belegt die Zusagen
von `slice-routing-spec-nachzug`, er ändert sie nicht; eine Abweichung der Messung
von einer Zusage ist ein Befund an die Spec, kein stiller Nachzug).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make test-integration` belegt am laufenden Feed-Container, dass eine per
SQL beantragte Routing-Regel die danach erfasste Change auf allen Lesewegen
auswählbar macht, dass die Auflösung bei Mehrdeutigkeit definiert ist, dass eine
nicht anwendbare Regel sichtbar endet, und dass `LH-FA-CFG-008` im RTM nicht mehr
Waise ist. Drei Liefer-Punkte:

- (A) **Wirkung:** Happy Path (Herkunfts- und Inhaltsregel; Ziel A sieht nur A, ein
  ungefilterter Leser sieht alles — über `cdc.changes`, `GET /changes`, gRPC, SSE und
  das NATS-Subjekt), Boundary (zwei treffende Regeln, `order` entscheidet; je eine
  Verletzung R1 bis R6 endet `failed` mit Klartext der Spec, der Regelstand bleibt),
  Neustart-Festigkeit (`docker restart`), Ausschluss-Sperre R3 in beide Richtungen,
  Replay (erneutes Lesen liefert dasselbe Label, auch nach einer Regeländerung),
  Backfill-Bestand mit Label;
- (B) **Negative und Abhilfe:** eine auf eine Change nicht anwendbare Regel endet
  sichtbar mit der Klasse `schema` (`diagnose`, Heartbeat); die Abhilfe über
  `cdc.remove_route` und Neustart — **soweit die Vorab-Bedingung V3 der Welle den
  Fall am System erzeugbar macht**;
- (C) **DELETE-Messung und RTM:** die Messung des Verhaltens einer Inhaltsregel auf
  eine Nicht-Schlüsselspalte bei DELETE ohne volle Replica-Identität gegen
  PostgreSQL 17 **und** 18 ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Entscheidung 4), und die RTM-Deckung: `make doc-trace` führt
  [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) nicht mehr unter den Waisen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktivcode** — alle Wirkung steht in den Slices davor; ein im Lauf gefundener
  Fehler wird gemeldet und in einem eigenen Slice behoben (anderer Vorgang, kein
  Verdecken durch Testanpassung).
- **Das Handbuch** — `slice-routing-betriebsdoku`: es beschreibt, was dieser Slice
  gemessen hat; die DELETE-Aussage steht dort erst mit der Messung dieses Slice.
- **SDK-Realserver-Tiers** (`make test-sdk-*-integration`) — `slice-routing-sdk-beispiel-target`
  trägt die SDK-Parameter mit Unit-Tests; ein Realserver-Beleg der drei SDKs für das
  Routing wäre ein weiterer Slice (Welle §6 führt ihn als ausgeschlossen; offene
  Frage an den Auftraggeber im Eröffnungs-Bericht).
- **Der `e2e.yml`-Workflow** — keine strukturelle Änderung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht); die beiden Legs der
  Matrix fahren das erweiterte Testpaket, ihr Lauf nach dem Push ist der Beleg des
  Closure-Kriteriums zur DELETE-Messung (Welle §3).
- **Eine allgemeine Recovery nach Schema-Fehlern** — `BEO-PGC/kein-admin-weg-schema-fehler-recovery`
  (offen, 1×); die Abhilfe gilt der nicht anwendbaren Routing-Regel.

## 2. Definition of Done

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Happy Path und Boundary (A):
      ein realer, grüner `make test-integration`-Lauf trägt (1) eine Herkunftsregel
      (ohne `when`) und eine Inhaltsregel (`when`), beide per
      `SELECT cdc.set_route(...)` beantragt und per Poll auf `applied` bestätigt; die
      danach erfasste Change trägt das Ziel der Regel, und jedes Ziel ist an `cdc.changes`
      (`WHERE route_target`), `GET /changes?target=`, dem gRPC-Stream, dem SSE-Stream
      und dem NATS-Subjekt `cdc.route.<source_id>.<ziel>` auswählbar — jedes gelesene
      Ziel gegen die persistierte Zeile derselben `change_id` gehalten; Ziel A sieht nur A,
      ein ungefilterter Leser sieht jede Change (geroutete eingeschlossen), eine Change ohne
      Treffer hat `route_target IS NULL` und erscheint nur ungefiltert; (2) zwei treffende
      Regeln: die kleinere `order` gewinnt; (3) je eine Verletzung von R1 bis R6 endet
      `failed` mit dem Klartext der Spec, der Regelstand bleibt, die Gegenprobe ohne die
      Verletzung endet `applied`; (4) nach einem realen `docker restart` leitet der
      Prozessstart den Regelstand aus den `applied`-Zeilen ab, die danach erfasste Change
      trägt das Ziel; (5) R3 in beide Richtungen (Regel gegen ausgeschlossene Spalte,
      `exclude_column` gegen Routing-Spalte) und [`LH-QA-SEC-004`](../../../../spec/lastenheft.md):
      nach dem Ausschluss einer Spalte erscheinen weder ihr Wert noch ein Ziel, das ihn
      verriete; (6) Replay: dieselbe `change_id` liefert nach einer Regeländerung dasselbe
      Label über `cdc.changes` und `GET /changes`, eine vor der Regel erfasste Change bleibt
      `NULL` (keine Rückwirkung), ein neuer Backfill-Run mit dem aktuellen Regelstand trägt
      das Label des Bestands (`origin = 'backfill'`). *Zu belegen durch:* `make test-integration`
      (Phasen im Runner `tools/harness/run-integration-tests.sh`, Testfunktionen in
      `test/integration/routing_e2e_test.go`; die gedruckten Zeilen je Phase stehen im
      Bericht). Die Wegwerf-Clients unter `tools/harness/` (`httpclient`, `grpcclient`,
      `sseclient`, `natsstreamsub`) erhalten die Auswahl des Ziels als Flag.
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Negative (B), nach
      Vorab-Bedingung V3: **erster Schritt des Slice** ist die Messung, welche Ursache
      `ErrRoutingNotApplicable` am laufenden System erzeugt (Kandidaten: Entfernen der
      Spalte — nach dem Bestand der Pfad der inkompatiblen Schemaänderung —, eine
      Publication mit Spaltenliste, deren Relation die Bedingungsspalte nicht trägt;
      *hergeleitet*, keiner geprüft); die gedruckte Messung steht im Bericht. Ist ein Fall
      erzeugbar: der Erfassungspfad endet sichtbar mit der Klasse `schema` samt
      Log-Sentinel (`diagnose`/Heartbeat), `cdc.remove_route` wird beantragt, während der
      Prozess steht (Antrag `pending`), nach einem realen Neustart ist der Antrag `applied`,
      bevor die erste Transaktion der Tabelle assembliert wird, und die zuvor nicht
      bestätigte Transaktion erscheint über `cdc.changes` ohne Label und ohne zweiten
      `schema`-Fehler. Ist kein Fall erzeugbar: das Architect-Verdikt (Welle §5 V3) liegt
      vor, und der Closure-Bericht der Welle nennt, welcher Beleg das Negative-Kriterium
      von [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) trägt — die Verengung auf
      einen Unit-Beleg steht dort mit Namen, nicht stillschweigend. Die Phase läuft als
      letzte vor der Container-Ende-Grenze des Runners. *Zu belegen durch:*
      `make test-integration`; `TestE2E…`-Funktion mit den Kennungen im Godoc
      (Erzeugnis-Eingabe der E2E-Abdeckung).
- [ ] [`LH-QA-POR-001`](../../../../spec/lastenheft.md) und RTM (C): die Aussage zu
      DELETE ist gemessen — eine Inhaltsregel auf eine Nicht-Schlüsselspalte, DELETE einer
      Zeile ohne volle Replica-Identität: das Verhalten (erwartet: Wert abwesend, Regel
      Nicht-Treffer, nächste Regel bzw. `NULL`) mit der Gegenprobe unter voller
      Replica-Identität (Treffer) — **an PostgreSQL 17 und an 18**: `make test-integration`
      lokal je einmal mit `PG_TEST_IMAGE` auf den Digest des Legs (die Digests stehen in
      `.github/workflows/e2e.yml`), gedruckt je Lauf die PostgreSQL-Version. Weicht die
      Messung von der Erwartung ab, steht der Befund im Bericht und die Aussage in der ADR
      ist Gegenstand eines Architect-Zugs vor dem Handbuch. `make doc-trace` führt
      [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) nicht mehr unter den Waisen: die
      gedruckte Zeile am Arbeitsbaum nennt `0 Waise(n)`; der Parent-Stand (`30fd6cb5`,
      gemessen, Exit 0) lautet `80 Anforderung(en), 1 Waise(n)`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-e2e.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` (Zeile `make doc-trace`: die Aussage über die
      Waisen-Zahl und `LH-FA-CFG-008`; Zeile `make test-integration`: die Routing-Belege)
      und `docs/user/e2e-abdeckung.md` (Erzeugnis des Runners, nicht von Hand); das
      Benutzerhandbuch bleibt unberührt (Adresse: `slice-routing-betriebsdoku`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/routing_e2e_test.go` (neu) | neu | Testfunktionen `TestE2E…` (Namen des Implementers) mit den Kennungen im Godoc — die Godocs sind Erzeugnis-Eingabe der E2E-Abdeckungstabelle ([`AGENTS.md`](../../../../AGENTS.md) §3.7, Ausnahme). Formvorbild `test/integration/transformation_e2e_test.go` (Stand `30fd6cb5`: drei Dateien, 22 `func TestE2E*`). |
| `tools/harness/run-integration-tests.sh` | update | neue Phasen (Happy Path samt Lesewegen, Boundary R1–R6, Neustart, Ausschluss-Sperre, Replay, Backfill-Bestand, Negative); `abdeckung_declare` je Phase (Parent: 52 Aufrufe). Die Negative-Phase endet den Container und steht hinter allen Phasen, die einen laufenden Container brauchen. |
| `tools/harness/httpclient/`, `tools/harness/grpcclient/`, `tools/harness/sseclient/`, `tools/harness/natsstreamsub/` | update | Flag für die Auswahl des Ziels (bei `natsstreamsub`: Abonnement des Ziel-Subjekts); Wegwerf-Clients des E2E-Tiers, keine SDK-Pakete. |
| `tools/harness/lib-sdk-rule-fixture.sh` | lesen | Formvorbild der Regel-Vorbereitung per SQL (`cdc.set_transformation`), nicht Teil. |
| `docs/user/e2e-abdeckung.md` | update (Erzeugnis des Runners) | die Zeilen für `LH-FA-CFG-008`; wird vom Lauf geschrieben und committet, nicht von Hand. |
| `harness/README.md` | update | Zeile `make doc-trace` (Waisen-Aussage) und Zeile `make test-integration` (Beschreibung der Belege). |
| `compose.yaml`, `.github/workflows/e2e.yml`, `.d-check.yml` | lesen | erwartet unverändert: `trace.coverage` liest `docs/user/e2e-abdeckung.md` bereits (Stand `30fd6cb5`). |

**Aufbau der Messungen** (Festlegungen des Plans, jede mit Beleg des Bestands am
Start vom Implementer zu prüfen): eine Tabelle für die Herkunftsregel und eine für die
Inhaltsregeln (eine Spalte mit zwei Werten und einem dritten ohne Treffer); der
DELETE-Fall braucht eine Tabelle ohne volle Replica-Identität (Default) und, als
Gegenprobe, `REPLICA IDENTITY FULL` an derselben Tabelle; die Phasen nutzen eigene,
wegwerfbare Tabellen (Isolation: kein Zustand wandert zwischen Phasen,
`BEO-PGC/test-isolation-geteilter-zustand`).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CFG-008` ist
Waise im RTM"; „der E2E-Lauf belegt die Zustellwege der Transformationen"; Parent ist
`30fd6cb5`; der Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und
Nichtgefundenes ein):**

```suchlauf
30fd6cb5 1 -n -E 'Waise' -- harness docs/user README.md
30fd6cb5 1 -n -E 'CFG-008' -- harness docs/user README.md .d-check.yml
30fd6cb5 14 -n -E 'CFG-007' -- docs/user
30fd6cb5 52 -n -E 'abdeckung_declare' -- tools/harness
30fd6cb5 7 -n -E 'PG_TEST_IMAGE' -- compose.yaml .github/workflows/e2e.yml tools/harness/run-integration-tests.sh
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aussage über die Waisen | Zeile 1: 1 Zeile (`harness/README.md:135`, Zeile `make doc-trace`: „80 Anforderungen, 1 Waise", Messung vom 2026-09-26) | nach dem Lauf neu messen (`make doc-trace`) und die Zeile auf die gemessene Zahl ziehen; die Zeile nennt die Zahl mit ihrem Ursprung und Lauf; Befund am Diff: einzutragen |
| Nennung von `LH-FA-CFG-008` außerhalb von Spec und Records | Zeile 2: 1 Zeile (dieselbe, `harness/README.md:135`) | im Diff trägt `docs/user/e2e-abdeckung.md` die neuen Zeilen; Befund: einzutragen |
| Zeilen der Transformationen im Abdeckungs-Träger als Formvorbild | Zeile 3: 14 Zeilen mit `CFG-007` in `docs/user` | die Routing-Zeilen folgen derselben Form; Befund: einzutragen |
| Phasen-Deklarationen des Runners | Zeile 4: 52 Aufrufe | jede neue Phase deklariert sich; kein stiller Ausschluss (`BEO-PGC/test-runner-stiller-ausschluss`); Befund: einzutragen |
| PostgreSQL-Image des Tiers | Zeile 5: 7 Zeilen (`compose.yaml`, `e2e.yml`, Runner) | die lokale Übersteuerung auf PostgreSQL 17 läuft über `PG_TEST_IMAGE`; Befund: einzutragen |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-antragsweg`,
`slice-routing-backfill-pfad`, `slice-routing-lesewege` und `slice-routing-nats-subjekt`
liegen in `done/` (Welle §5: alles, was der E2E belegt, steht) und kein anderer Slice
liegt in `in-progress/` (WIP-Limit 1). Am Start: `make image` ist ausgeführt
(`compose.yaml` referenziert das lokal gebaute Image,
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); das Architect-Verdikt zu V3
liegt vor, **wenn** die erste Messung des Slice den Fall nicht erzeugen kann.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): (B) trennt sich als
  `slice-routing-e2e-abhilfe` ab (wie in der Transformations-Welle: zweite
  Container-Ende-Grenze, eigener Runner-Platz); (A) und (C) liefern dann.
- `in-progress` → `open` (blockiert): ein im Lauf gefundener Produktivfehler (z. B.
  ein Lesepfad liefert ein anderes Ziel als `cdc.changes`) — der Fehler wird in einem
  eigenen Slice behoben, der Beleg wartet; ein Rot geht nie als Anpassung der
  Erwartung in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), je ein realer, grüner `make test-integration`-Lauf gegen PostgreSQL 18
und gegen PostgreSQL 17 (gedruckte Version im Bericht), `make doc-trace` mit
`0 Waise(n)`, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Nichtanwendbarkeit ist am System nicht erzeugbar (V3).** Entfernt der
  Betreiber die Bedingungsspalte, endet der Prozess nach dem Bestand im Pfad der
  inkompatiblen Schemaänderung, bevor die Regel zählt (*hergeleitet*); die Abhilfe
  `cdc.remove_route` + Neustart ist dann nicht der Weg aus diesem Zustand. — **Ausgang:**
  bei der Closure einzutragen (gemessene Ursache oder Verdikt; Welle-Closure-Bericht).
- **Die DELETE-Erwartung kann falsch sein.** „Wert abwesend, Regel Nicht-Treffer"
  ist *hergeleitet* aus dem Abwesenheits-Vertrag von `LH-FA-DAT-005` und der
  Transformations-ADR, an PostgreSQL 17 und 18 nicht gemessen. — **Ausgang:** bei der
  Closure einzutragen (Messung, gedruckte Zeilen je Version).
- **Laufzeit des Testpakets gegen das 60-Minuten-Limit von `e2e.yml`.** Jede Phase
  verlängert den Lauf; ob das Limit reicht, ist bis zum Lauf auf dem gehosteten Runner
  offen (`BEO-PGC/github-actions-unverifizierbar-lokal`, verkörpert, 8×; kein
  Workflow-Zug in dieser Welle, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht,
  das Laufzeit-Risiko bleibt). — **Ausgang:** bei der Closure einzutragen (gemessene
  lokale Laufzeit vor und nach; Lauf-ID der Legs, sobald vorhanden).
- **Isolation und stiller Ausschluss.** Neue Phasen teilen Replikationsslot, Queue und
  Container mit den bestehenden (`BEO-PGC/test-isolation-geteilter-zustand`, offen,
  2×; `BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). — **Ausgang:** bei der
  Closure einzutragen (gedruckte Zeile je Phase; Phasen mit eigenen Tabellen).
- **Ein Beleg, der nur den Reviewer überzeugt.** Eine Boundary-Aussage, die allein
  der Reviewer an der Zahl gemessen hat (`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`,
  offen, 1×). — **Ausgang:** bei der Closure einzutragen.
- **Timing-Flake der Phase „Leerlauf-Bestätigung".** Ein Rot dort lässt alle Phasen
  dahinter ungelaufen (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert,
  5×); ein Rot dieser Signatur ist weder Beleg noch Widerlegung. — **Ausgang:** bei der
  Closure einzutragen.
- **Erzeugnis `docs/user/e2e-abdeckung.md`.** Der Runner schreibt die Datei; eine
  Hand-Änderung geht beim nächsten Lauf verloren
  (`BEO-PGC/test-schreibt-in-committete-datei`, verkörpert, 4×). — **Ausgang:** bei der
  Closure einzutragen (Datei kommt aus dem Lauf, `git diff` zeigt nur Routing-Zeilen).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `test/integration/`, `tools/harness/` und `docs/user/e2e-abdeckung.md` — eine
Sub-Area (Testschicht).

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/test-isolation-geteilter-zustand` (offen, 2×),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×),
`BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 5×),
`BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 18×),
`BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 8×). Mit diesem Slice
erreicht `test-isolation-geteilter-zustand` oder `test-runner-stiller-ausschluss` bei
einem weiteren Auftreten 3× und braucht dann einen Ausgang.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

