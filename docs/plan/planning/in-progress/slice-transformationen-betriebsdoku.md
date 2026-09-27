# Slice transformationen-betriebsdoku: Betriebsdokumentation und SDK-Beleg — Transformationsregeln im Benutzerhandbuch, Row-Image-Schlüssel in den SDK-Modellen als undurchsichtig belegt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md),
[`LH-FA-ADM-001`](../../../../spec/lastenheft.md),
[`LH-FA-SST-009`](../../../../spec/lastenheft.md) (die drei SDK-Packages),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 6 (SDK-/Doku-Beleg) und Teilfrage 8 (keine Änderung des
Nachrichtenschemas).

**Berührte Spec-Stellen:** [`SPEC-026`](../../../../spec/pflichtenheft.md),
[`SPEC-027`](../../../../spec/pflichtenheft.md),
[`SPEC-028`](../../../../spec/pflichtenheft.md) (die drei SDK-Packages),
[`SPEC-018`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) (Drahtverträge) — gelesen,
nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Zwei Belege, beide erst jetzt wahr formulierbar: (1) die drei
SDK-Packages setzen keine bestimmte Menge von Row-Image-Schlüsseln voraus — zu
belegen, nicht als geprüft behauptet (§3.12 Instanz B), erwartet wird keine
Code-Änderung
([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 8); der Realserver-Beleg liegt bei `slice-sdk-regel-realserver-e2e`, dieser
Slice führt den Suchlauf über den Produktivcode und nennt jenen Beleg; (2) das Benutzerhandbuch beschreibt die Betreiber-Oberfläche der
Transformationsregeln, nachdem `e2e-wirkung` und `e2e-abhilfe` die Wirkung
belegt haben. Das ist der **aufgeschobene Gegenstand** von `antragsweg-schema`
(Adresse, siehe dessen §1) und die Adresse aller Slices, die Handbuch-Anteile
an dieses Slice gemeldet haben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine SDK-Code-Änderung oder ein SDK-Release** — erwartet keine; findet der
  Suchlauf einen Zugriff auf feste Schlüssel im Produktivcode eines SDK, ist
  das ein Plan-Nachzug (Rückführung §4) samt Versionshebung, kein stiller
  Zusatz. Ein Tag-Push bleibt Betreiber-Handlung außerhalb der Welle.
- **Ein SDK-Realserver-E2E mit aktiver Regel** — trägt
  `slice-sdk-regel-realserver-e2e` (Start nach `e2e-wirkung`, vor diesem Slice):
  die drei SDK-Tiers fahren dort je vier Phasen gegen eine aktive
  `rename_column`-Regel; dieser Slice führt den Beleg im Bericht und im Handbuch
  als erprobt.
- **Aussagen, die kein Beleg trägt** — das Handbuch nennt nur, was
  `e2e-wirkung` und `e2e-abhilfe` real belegt haben; die Aussage „die Abhilfe
  wirkt“ steht erst mit dem Beleg.
- **Routing** ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) — keine
  Handbuch-Aussage; das Routing hat keine Entscheidung.
- **Ein allgemeiner Recovery-Weg für Schema-Fehler** — das Handbuch beschreibt
  die Abhilfe für die Nichtanwendbarkeit einer Transformationsregel, nicht mehr
  (`BEO-PGC/kein-admin-weg-schema-fehler-recovery`, offen, 1×).

## 2. Definition of Done

- [x] SDK-Beleg: die Modelle der Row Images in den drei Packages behandeln
      `old_image`/`new_image` undurchsichtig (Anker am Start gelesen:
      `sdks/csharp/PgChangeFeed.Client/…/Change.cs`,
      `sdks/kotlin/pgchangefeed-kotlin/…/Change.kt`,
      `sdks/python/pgchangefeed/src/pgchangefeed/models.py`), und **kein
      Produktivcode** unter `sdks/*/` greift auf einen festen Bild-Schlüssel zu
      (Suchlauf, Befund und Nichtbefund im Bericht); die Test-Schlüssel der
      SDK-Tests sind Fixture-Daten. *Zu belegen durch:* der Suchlauf in §3, Lesen
      der Modelle und der Bericht von `slice-sdk-regel-realserver-e2e` (die drei
      Tiers fahren je vier Phasen gegen eine aktive Regel; der Slice liegt am
      Start in `done/`); dieser Slice ändert keinen SDK-Code und führt kein
      `make sdk-*` aus, die Nicht-Ausführung wird im Bericht begründet.
- [x] Betriebsdokumentation im Benutzerhandbuch
      ([`docs/user/benutzerhandbuch.md`](../../../user/benutzerhandbuch.md)):
      ein neuer Abschnitt in §4 (Aufgaben) „Transformationsregel konfigurieren“
      mit dem Aufruf von `cdc.set_transformation`/`cdc.remove_transformation`
      (Beispiele `rename_column`, `map_value`; Aufrufform der Regelform: der
      Parameter `rule_spec` ist `json`, ein Literal oder ein `::json`-Wert wird
      angenommen, ein `::jsonb`-Wert nicht — „function cdc.set_transformation(…,
      jsonb) does not exist“, gemessen im Review
      `review-slice-transformationen-antragsweg-schema`,
      [`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md)
      Folgepflicht 2; Annahmemenge nach
      [`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md)
      Festlegung 1: SQL-`NULL`, JSON-`null`, ein Wert ohne Objekt und ein
      doppelter Schlüssel werden angenommen und in Go geprüft, `\u0000`, eine
      Zahl außerhalb des Zahlbereichs von `numeric` und ein Syntaxfehler enden
      beim Aufruf ohne Antrags-Zeile), der Rolle `cdc_admin`, dem
      Status `pending`/`applied`/`failed` und den Fehlertexten von K1–K4, der
      Wirkung („ab `applied` für künftige Changes, nicht rückwirkend; die
      Rohform wird nicht gespeichert“, der Informationsverlust bei
      `map_value`, Übergabe aus `slice-transformationen-map-value` unten), dem
      Verhältnis zum Spaltenausschluss (der Ausschluss gilt
      zuerst), der Dauerhaftigkeit (überlebt Neustart und
      Deaktivierung/Aktivierung; Grenze aus `antragsweg-usecase`: eine vermerkte Regel, die
      ein älterer Binärstand nicht lesen kann — etwa eine `map_value`-Regel unter einem Stand
      vor `map-value` —, hält Prozessstart und jeden Regel-Antrag der Quelle an, hergeleitet
      aus dem Quelltext, nicht erprobt), den Zeilen, die kein Antrag sind
      (Übergabe aus `antragsqueue-lesefehler-failed`,
      [`SPEC-019`](../../../../spec/pflichtenheft.md) Absatz „Zeilen, die kein
      Antrag sind“: die SQL-Funktionen prüfen Quelle, Schema, Tabelle und Spalte
      nicht; eine `pending`-Zeile, die der Antrags-Konstruktor verwirft — leeres
      Schema, leerer Tabellenname, `exclude_column`/`include_column` mit leerer
      Spalte —, endet `failed` mit dem Klartext und der Antrags-Kennung in
      `error_message` und hält die Zeilen dahinter nicht an; eine Zeile ohne
      Kennung bleibt `pending` mit einer Warnung im Log; belegt durch den
      Store-Test `TestAdministrationRequestListPendingPassesRejectedRowsThrough`
      und den Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins`;
      die Gründe „Quelle ist leer“ und „Antragsart ist unbekannt“ entstehen über
      die SQL-Funktionen nicht — Fremdschlüssel und `CHECK` — und sind nur an
      einem Fake-Test belegt, sie gehören nicht ins Handbuch), der
      Nichtanwendbarkeit samt Abhilfe-Prozedur
      (`cdc.remove_transformation` beantragen, Prozess starten — so, wie
      `e2e-abhilfe` sie belegt hat) **und ihrer Form im Backfill-Run**
      ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
      Folgepflicht 4, am laufenden System belegt von der Phase
      „Backfill-Regelstand“: der Run endet `failed` mit der Klasse `schema` ohne
      Change, wenn eine Regel auf die Spalten des Snapshots nicht anwendbar ist —
      Spalte fehlt, Zielname kollidiert —; er ist run-lokal, der Erfassungspfad
      läuft weiter; die Abhilfe ist die Regel zu entfernen und danach einen
      **neuen** `cdc.backfill_table`-Antrag zu stellen, **kein** Prozessneustart,
      der `failed`-Run wird nicht fortgesetzt und sperrt den neuen Antrag nicht;
      wechselt der Regelstand während des Runs, endet er mit `configuration`),
      der Reihenfolge der Aufrufe (Aufrufe einer
      Transaktion werden in Aufrufreihenfolge verarbeitet; Anträge auf dieselbe
      Regel oder Spalte nicht aus überlappenden Transaktionen absetzen —
      [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)
      Folgepflicht 4 und Festlegung 3 Punkt 1; die Aussage gilt für alle sieben
      Funktionen und gehört in den Transformations-Abschnitt in §4) und dem
      Backfill-Bezug (Regeln gelten auch für einen Backfill-Run: der Bestand
      trägt die transformierte Form, belegt von der Phase „Backfill-Regelstand“
      des Runners von `make test-integration`; der Regelstand gilt ab dem Öffnen des
      Snapshots des Runs — Festlegung des Implementers von
      `slice-transformationen-backfill-pfad`, von Review und Verifikation als konform
      zur Entscheidung gelesen, ohne Spec-Satz zum Zeitpunkt; im Handbuch als Zusage
      des Ist-Verhaltens, nicht als Spec-Aussage); dazu die Zeile
      `schema` in §6 Fehlerklassen, die Rollen-Beschreibung in §2, das Glossar
      und die Änderungshistorie. *Zu belegen durch:* Review des Abschnitts
      gegen die Belege von `e2e-wirkung`/`e2e-abhilfe` und `make docs-check`.
- [x] Jede Zahl und jede Wirkungs-Aussage trägt ihren Ursprung
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12): gemessen (Lauf), übernommen
      (Beleg-Anker) oder abgeleitet; eine Aussage, die nur die ADR trägt, steht
      als Zusage bzw. „erwartet“. Die Änderungshistorie trägt eine Zeile
      (`BEO-PGC/handbuch-versionshistorie-uebersprungen`, verkörpert, 3×). *Zu
      belegen durch:* Review; die Suchläufe in §3 belegen, dass keine
      Handbuch-Stelle mit einer Aufzählung der Antragsarten oder der Klasse
      `schema` übersehen ist.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). Report:
      [`review-slice-transformationen-betriebsdoku.md`](../../../reviews/review-slice-transformationen-betriebsdoku.md)
      (1 HIGH, 1 MEDIUM); Fixrunde behebt F-1 (Fehlertext-Zitat korrigiert)
      und F-2 (drittes Beispiel „Zeilen, die kein Antrag sind“ ergänzt) —
      kein offenes HIGH.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: der Slice **ist** das Doku-Update (Handbuch,
      Änderungshistorie); `harness/README.md` ist unberührt — die Sensor-Zeilen
      tragen `e2e-wirkung` und `e2e-abhilfe`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

**Übergabe aus `slice-transformationen-map-value`** (gemeldet, kein zusätzlicher
Umfang; Herkunft: Review-Report Finding F-6, Verifikations-Report V-3 und §7,
gelesen am Stand `8b9b9c5b`; Frist der Meldung: die Closure des meldenden Slice,
gezogen). Der Handbuch-Abschnitt zu `map_value` trägt vier Aussagen, jede mit
ihrem Ursprung ([`AGENTS.md`](../../../../AGENTS.md) §3.12):

- **Informationsverlust.** Bilden mehrere Quellwerte auf denselben Zielwert ab,
  meldet das System nichts — kein Fehler, keine Warnung; die Abbildung eines
  Werts auf sich selbst ist zulässig, ein leeres `values` wird abgelehnt (Tests
  der Domäne, gelesen im Review). Die Rohform wird nicht gespeichert, die
  Abbildung ist nicht umkehrbar
  ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Konsequenzen; im Handbuch als Zusage der Entscheidung, nicht als Messung).
- **Vergleich.** Der Wert einer Spalte wird als Zeichenkette zeichengenau mit
  den Schlüsseln von `values` verglichen, ohne Normalisierung: die
  Groß-/Kleinschreibung zählt, das Präfix eines Schlüssels ist kein Schlüssel,
  `é` als ein Zeichen und `e` mit U+0301 sind verschiedene Schlüssel (Test
  `TestBuildRowImageMapValue`, Fälle „Groß-/Kleinschreibung zählt“ und „Präfix
  eines Schlüssels ist kein Schlüssel“; das Zeichenpaar: Negativbefund des
  Reviews, **übernommen**).
- **Größenordnung der Suche.** Die Suche einer Zuordnung ist linear in der Zahl
  der Paare (hergeleitet aus dem Quelltext von `lookupMappedValue`). Die Zahlen
  sind **gemessen** im Architect-Zug (Verdikt
  [`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md)
  §6 und §1 Zeile M1; Wegwerf-Benchmark in einer Kopie von `internal/domain`,
  Schlüssel fehlt als ungünstigster Fall, ohne `-race`, `-benchtime 300x
  -count 3`, i9-13900H, Bereiche über die drei Läufe): 54 bis 123 ns bei 10
  Paaren, 6,5 bis 6,7 µs bei 1000 Paaren, 65 bis 80 µs bei 10 000 Paaren, 0,66
  bis 0,70 ms bei 100 000 Paaren je Wert und Regel; die Kosten fallen je Wert
  einer Spalte mit Regel an, je Zeile eines Backfill-Blocks und je Change im
  Erfassungspfad. Die früheren Werte des Verifiers (80 ns, 6,8 µs, 0,72 ms bei
  10, 1000 und 100 000 Paaren) sind **übernommen** und bestätigen die
  Größenordnung; das Handbuch nennt die Zahlen des Verdikts mit diesen
  Bedingungen. **Ob die Zahl der Paare eine Obergrenze braucht, ist
  entschieden:** benannter Verzicht, keine Änderung an
  [`SPEC-030`](../../../../spec/pflichtenheft.md) (Verdikt §6; Trigger: ein
  Betreiber meldet ein Wachstum von `cdc_capture_lag` mit einer aktiven
  `map_value`-Regel, oder eine Regel mit mehr als 10 000 Paaren steht in der
  Queue). Das Handbuch führt keine Obergrenze; die abgeleitete Last der Stufe
  „groß“ (1000 Changes/s, zwei Suchen je Change: 1,3 % eines Kerns bei 1000
  Paaren, 13 bis 16 % bei 10 000, über 100 % bei 100 000) steht dort nur mit dem
  Vermerk **abgeleitet**.
- **Rückfall auf einen älteren Binärstand** (die Grenze steht oben im
  Handbuch-Punkt zur Dauerhaftigkeit). Erprobt ist der erste Schritt der
  Herleitung (Verifikations-Report §7, **übernommen**): am Parent-Stand von
  `slice-transformationen-map-value` liefert `FoldTransformations` für eine
  `applied`-Zeile mit `map_value` den Fehler `Regelstand: Regel "r": unbekannter
  Regeltyp: map_value`. Dass dieser Fehler Prozessstart und jeden Regel-Antrag der
  Quelle anhält, bleibt hergeleitet; der Prozessstart ist nicht gefahren.

**Übergabe aus dem Architect-Verdikt** [`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md)
(§9.3, Planner-Nachzug vom 2026-09-27, Frist: vor dem Start dieses Slice,
gezogen). Der Handbuch-Abschnitt trägt fünf Aussagen, jede mit ihrem Ursprung
([`AGENTS.md`](../../../../AGENTS.md) §3.12); der Ursprung ist der Anker des
Verdikts, dieser Plan erweitert ihn nicht:

- **Wartezeit des Stream-Starts.** Ein hängender Antrag hält den Stream-Start
  höchstens 30 s an (Frist des Vorlaufs,
  [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  Festlegung 2). Der Wert 30 s ist **hergeleitet** (die Hälfte der
  Fehlergrenze von 60 s des Capture-Abstands), nicht gemessen; die Aussage
  „höchstens 30 s“ steht im Handbuch mit der **Messung des Rundlaufs**, den
  `slice-start-vorlauf-grenze` liefert (gedruckte
  Zeile des Laufs), nicht als Herleitung — liegt dieser Beleg am Start nicht
  vor, steht die Aussage als Zusage der ADR. Der Healthcheck bleibt über die
  Wartezeit gesund (Exit 0, Heartbeat-Alter höchstens 7,1 s): **gemessen** im
  Architect-Zug an drei Läufen, **übernommen** (Anker: [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  §Gemessen, Läufe A bis C); eine Anzeige in `diagnose` gibt es nicht, der Warn-Eintrag im
  Log trägt den Ablauf der Frist. Bei Ablauf bleibt der Antrag `pending`
  (Festlegung 3); eine andere Antragsart als `enable` ist in den Läufen nicht
  gefahren, die Gleichheit für jede Art ist **hergeleitet**.
- **Größenordnung der Suche einer Zuordnung:** siehe den Punkt „Größenordnung
  der Suche“ oben (Architect-Werte **gemessen**, Verifier-Werte **übernommen**,
  Obergrenze: benannter Verzicht).
- **Kosten der Lesung im Backfill.** Ein Backfill-Run liest den Regelstand und
  den Ausschlussstand aus den `applied`-Zeilen aller Tabellen der Quelle, je
  Block und vor dem Commit. **Gemessen** (Verdikt §3, Messung M2, PostgreSQL 18
  im Container ohne Netz, synthetische Queue ohne Index außer dem
  Primärschlüssel): die SQL-Lesung kostet 0,26 bis 0,31 ms bei 100 Zeilen der
  Queue und 186 bis 201 ms bei 1 000 000 Zeilen; die Faltung der Regeln kostet
  94 ms bei 100 000 Anträgen (M1), die Werte bei 30 000 und 300 000 gelesenen
  Anträgen sind daraus **linear hochgerechnet**. **Abgeleitet:** etwa 8 % einer
  Blockdauer von 0,12 bis 0,13 s bei 10 000 Zeilen der Queue einer Quelle, etwa
  70 % bei 100 000 Zeilen (die Blockdauer ist aus dem Handbuch-Abschnitt
  Backfill **übernommen**). Eine Queue realer Größe gibt es nicht (kein
  Server-Tag trägt die Antragsarten): die Messung ist synthetisch, und die
  Aussage endet dort, wo sie gemessen ist. Betreiber-Hinweis: die Zahl der
  Antrags-Zeilen ist nicht begrenzt; ab rund 10 000 Zeilen einer Quelle kostet
  die Lesung je Block messbar (Trigger und Behebung: Verdikt §3).
- **„Alle Wege dieselbe Form“ nur mit der gemessenen Menge.** Gemessen ist:
  eine per `cdc.set_transformation` beantragte `rename_column`- und
  `map_value`-Regel prägt eine danach **eingefügte Zeile** (INSERT, Neu-Bild) auf
  den fünf Wegen (Abdeckungstabelle
  [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md), Phase
  „Transformationen-Happy-Path (fünf Zustellwege)“); UPDATE, DELETE und das
  Alt-Bild sind über `cdc.changes` belegt
  (`TestE2ETransformationRulesShapeBothImages`). Für UPDATE, DELETE und das
  Alt-Bild auf den vier weiteren Wegen steht das Handbuch als **hergeleitet**
  (zwei Glieder: die Regel wirkt im `Assembler` auf beide Bilder; die Adapter
  bilden die Bilder unverändert ab, Verdikt §4 — am NATS-Publisher ist das
  zweite Glied nur gelesen, kein Unit-Test mit einem Alt-Bild, das nicht `null`
  ist).
- **Bezugspunkt des Regelstands im Run.** Der Regelstand gilt ab dem Öffnen des
  Snapshots des Runs — schon im Handbuch-Punkt zur Backfill-Beziehung oben
  geführt (Festlegung des Implementers von
  `slice-transformationen-backfill-pfad`, kein Spec-Satz), hier nicht
  erweitert.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | neuer Abschnitt in §4, §2 Rollen, §6 Fehlerklassen (Zeile `schema`), Glossar, Änderungshistorie. |
| `sdks/**` (Modelle, Tests) | lesen | Beleg-Suchlauf; Änderung nur nach Plan-Nachzug. |
| SDK-READMEs (`sdks/csharp/README.md`, `sdks/python/README.md`, `sdks/kotlin/pgchangefeed-kotlin/README.md`) | prüfen | Aussagen über Row Images (jede README trägt einen Abschnitt zum Change-Objekt mit `old_image`/`new_image`, Suchlauf am Start); die README ist die Paketbeschreibung auf PyPI und NuGet und erscheint dort erst mit einer neuen Package-Version, eine Änderung wird im Bericht benannt. Jede Textdatei unter `sdks/` trägt keine interne Kennung (`SPEC-`/`ADR-`/`ARC-`/`LH-FA-`/`LH-QA-`, Slice-/Welle-Name): `make sdk-public-doc-check` färbt sonst jedes `make sdk-pack-*`. |
| `docs/user/benutzerhandbuch-standard.md` | prüfen | Standard-Form des Handbuchs — der neue Abschnitt folgt ihr. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Menge der
Antragsarten“, „die Sätze über Row-Image-Schlüssel und Fehlerklasse `schema`“,
„die Aussagen der SDK-Modelle über Row Images“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Zugriffe auf feste Bild-Schlüssel im Produktivcode der SDKs | `grep -rn 'old_data\|new_data\|oldData\|newData\|old_image\|new_image' sdks --include=*.py --include=*.cs --include=*.kt`, dann Fundstellen in `src/main`/`src/pgchangefeed`/`PgChangeFeed.Client/` (ohne Tests) lesen | **Nichtbefund.** 137 Treffer an beiden Ständen (Parent `b80f45cf`, Diff — `sdks/` in diesem Slice unberührt); jede Fundstelle in Produktivcode ist eine Feld-Deklaration/-Zuweisung des obersten Nachrichtenfelds (`old_image`/`new_image` als `JsonElement?` in C#/Kotlin, `Any \| None` in Python, `bytes` im gRPC-Client) — kein Zugriff auf einen benannten Schlüssel **innerhalb** eines Row Image. Geprüft: `sdks/csharp/PgChangeFeed.Client/{Http,Sse,Nats}/Models/*.cs`, `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../{http,sse,nats}/model/*.kt` + `grpc/PgChangeFeedGrpcClient.kt`, `sdks/python/pgchangefeed/src/pgchangefeed/models.py`. Alle übrigen Treffer liegen in Tests/Integrationstests (Fixture-Daten) oder Docstrings, die die zehn Nachrichtenfelder aufzählen. | keine — Erwartung aus `ADR-0112` Teilfrage 8 bestätigt, keine SDK-Code-Änderung |
| Aufzählungen der Antragsarten im Handbuch | `grep -rn 'exclude_column\|enable_table' docs/user` | Parent `b80f45cf`: 6 Treffer (Spaltenausschluss-Abschnitt allein). Diff (Feature-Commit `97de226c`): 9 Treffer — die drei neuen Treffer liegen im neuen Abschnitt „Transformationsregel konfigurieren“ (Reihenfolge-der-Aufrufe-Absatz nennt `exclude_column`/`include_column`/`disable_table`/`enable_table` als Beispiel derselben Regel für alle sieben Funktionen; „Zeilen, die kein Antrag sind“ nennt keine Aufzählung der Antragsarten selbst). Keine Aufzählung der sieben Antragsarten wurde durch diesen Slice unvollständig gemacht — der neue Abschnitt fügt eine achte/neunte Art (`set_transformation`/`remove_transformation`) hinzu, ohne eine bestehende Aufzählungsstelle zu berühren. **Fixrunde** (Review F-2, drittes Beispiel „Zeilen, die kein Antrag sind“ ergänzt): 11 Treffer — die zwei neuen Treffer sind der ergänzte Absatz selbst (`exclude_column`/`include_column` mit leerer Spalte) und die Änderungshistorie-Zeile 1.68, die dieselben zwei Funktionsnamen als Beleg nennt; keine bestehende Aufzählungsstelle berührt. | keine |
| Sätze über die Fehlerklasse `schema` | `grep -rn 'schema' docs/user/benutzerhandbuch.md` | Parent `b80f45cf`: 40 Treffer (überwiegend `schema_name`/`schema_version`/`CDC_SCHEMA`, ein Treffer die Fehlerklassen-Zeile). Diff: 49 Treffer — die neun neuen Treffer liegen im neuen Abschnitt (Regeltypen-Tabelle, Nichtanwendbarkeits-Absatz, Backfill-Form) und in der erweiterten Fehlerklassen-Zeile selbst; „Container startet nicht“ und „Neustart nach einem Fehler“ bleiben unverändert und nennen die Nichtanwendbarkeit bewusst nicht (die Abhilfe steht im neuen Abschnitt, die Recovery wird nicht verallgemeinert, siehe §1 Abgrenzung). | keine |
| Beschreibung des SQL-Lesezugriffs `cdc.changes` und `GET /changes` im Handbuch | `grep -rn 'Row Image\|row_image\|new_data' docs/user` | Parent `b80f45cf`: 22 Treffer. Diff: 24 Treffer — die zwei neuen Treffer liegen im neuen Abschnitt (Einleitungssatz „wirkt auf das Row Image“, Backfill-Form). Die bestehenden Sätze über die Bildform (§4 „Änderungen lesen“, „Bestand als Backfill überführen“) sind von diesem Slice unverändert; sie nannten den Regelstand vor diesem Slice nicht, weil es keinen Regelstand gab — kein Nachzugsbedarf an ihnen, die neue Aussage „die Schlüsselmenge folgt dem Regelstand zum Erfassungszeitpunkt“ steht im neuen Abschnitt selbst (Zeile „Ergebnis“/„Wirkung“). | keine |
| Übernahme aus Nachbar-Dokumenten | Lesen der `harness/README.md`-Zeile `make test-integration` | Gelesen (kein Zählmaß): die Zeile nennt die vier Transformations-E2E-Rundläufe (Happy Path, Bilder aller Operationen, Boundary, Neustart mit Ausschluss) und den Backfill-Regelstand-Rundlauf — dieselben Belege, auf die der neue Abschnitt verweist (`TestE2ETransformationRulesShapeBothImages`, `TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass`, Phase „Backfill-Regelstand“, Phase „Transformationen-Happy-Path (fünf Zustellwege)“); keine Abweichung gefunden. | keine |

```suchlauf
b80f45cf 137 -E 'old_data|new_data|oldData|newData|old_image|new_image' -- ':(glob)sdks/**/*.py' ':(glob)sdks/**/*.cs' ':(glob)sdks/**/*.kt'
diff 137 -E 'old_data|new_data|oldData|newData|old_image|new_image' -- ':(glob)sdks/**/*.py' ':(glob)sdks/**/*.cs' ':(glob)sdks/**/*.kt'
b80f45cf 6 -E 'exclude_column|enable_table' -- docs/user
diff 11 -E 'exclude_column|enable_table' -- docs/user
b80f45cf 40 'schema' -- docs/user/benutzerhandbuch.md
diff 49 'schema' -- docs/user/benutzerhandbuch.md
b80f45cf 22 -E 'Row Image|row_image|new_data' -- docs/user
diff 24 -E 'Row Image|row_image|new_data' -- docs/user
```

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-e2e-abhilfe`
in `done/` liegt (die Wirkung und die Abhilfe sind belegt, bevor das Handbuch
sie beschreibt; sein eigener Start-Trigger setzt `slice-start-vorlauf-grenze`
in `done/` voraus, also liegt auch dessen Rundlauf zur Wartezeit des
Stream-Starts vor), `slice-sdk-regel-realserver-e2e` in `done/` liegt (Kante: der
SDK-Beleg von DoD Punkt 1 ist dann am realen Server erprobt) und kein anderer
Slice in `in-progress/` liegt (WIP-Limit 1). Letzter Slice der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  Doku-Zug mit einem Such-Beleg; sprengte er den Umfang, wäre der SDK-Beleg
  (erster Liefer-Punkt) der abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls der SDK-Suchlauf einen Zugriff auf
  feste Bild-Schlüssel im Produktivcode findet (Plan-Nachzug: SDK-Änderung samt
  Versionshebung, eine Architect-Frage nach
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 8 wird geführt) oder falls ein Beleg von
  `e2e-wirkung`/`e2e-abhilfe` eine Handbuch-Aussage nicht trägt (dann steht die
  Aussage nicht im Handbuch).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Das Handbuch behauptet mehr, als die Belege tragen**
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B), etwa „die Abhilfe
  wirkt“ vor dem realen Beleg. *Erwartet, zu belegen durch:* Review liest jede
  Wirkungs-Aussage gegen den Beleg-Anker; der Start-Trigger stellt den Beleg
  voran. **Ausgang:** *(bei Closure)*
- **Das Handbuch übersieht eine Stelle**, die die Antragsarten oder die Klasse
  `schema` aufzählt (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, offen,
  2×; `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  verkörpert, 3×). *Erwartet, zu belegen durch:* die Suchläufe in §3.
  **Ausgang:** *(bei Closure)*
- **Der SDK-Suchlauf findet einen Schlüssel-Zugriff**, den die Erwartung
  („keine Code-Änderung“) ausschließt. *Erwartet, zu belegen durch:* der
  Suchlauf mit Befund und Nichtbefund im Bericht. **Ausgang:** *(bei Closure)*
- **Die Abhilfe-Prozedur wird zu einer allgemeinen Recovery** verallgemeinert
  (`BEO-PGC/kein-admin-weg-schema-fehler-recovery`, offen, 1×). *Erwartet, zu
  belegen durch:* Review des Abschnitts auf die benannte Ursache. **Ausgang:**
  *(bei Closure)*
- **Zahlen und Fristen im Handbuch ohne Ursprung** (etwa eine Wartezeit nach
  dem Neustart) — jede Zahl trägt Lauf und Ursprung oder entfällt. *Erwartet,
  zu belegen durch:* Review. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); das Handbuch und die SDK-Verzeichnisse (`sdks/*/`, nur
gelesen) sind keine Sub-Areas dieses Slice — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, je 3×,
einschlägig — Kern dieses Slice),
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×, einschlägig —
dieser Slice ist die Adresse von `antragsweg-schema`; sein §2 deckt den
aufgeschobenen Gegenstand), `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
(offen, 2×), `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×,
Abgrenzung §1), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×) und
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 6×,
Risiko §6), `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (verkörpert,
3×, gesichtet — SDK-READMEs werden nur gelesen),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
