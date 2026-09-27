# Welle welle-transformationen — Transformationen: erfasste Changes tragen vor der Persistierung die durch deklarative Regeln bestimmte Form, konfiguriert über die SQL-Antrags-Queue — Closure-Notiz

**Welle:** welle-transformationen
**Abschluss:** 2026-09-27
**Verantwortlich:** pt9912

## Was wurde geliefert?

- **[`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Happy Path, Boundary, Negative)
  ist am laufenden Feed-Container belegt.** Ein Antrag `cdc.set_transformation`
  legt eine `rename_column`- oder `map_value`-Regel auf eine aktivierte Tabelle;
  jede danach erfasste Change trägt die transformierte Form — nicht die
  Rohform — über `cdc.changes`, `GET /changes`, gRPC-Stream, SSE-Stream und
  NATS-Vollinhalts-Stream in derselben Gestalt; eine gegen K1–K4 verstoßende
  Regel endet `failed` mit Klartext, eine am laufenden Prozess nichtanwendbar
  gewordene Regel beendet den Erfassungspfad sichtbar mit Fehlerklasse `schema`
  (`diagnose`, Heartbeat), und das Abhilfe-Akzeptanzkriterium (a)–(d) aus
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 5 ist am realen Prozess-Neustart belegt (`cdc.remove_transformation`
  während des Stillstands beantragt, nach dem Neustart `applied` **vor** der ersten
  Transaktion der Tabelle, kein Datenverlust). Die Regelauswertung liegt an genau
  **einer** Stelle (`model.BuildRowImage`, `internal/domain/model/rowimage.go`),
  aufgerufen vom WAL-Pfad (`mapper.go`) **und** vom Backfill-Pfad
  (`usecase/backfill/service.go`) — Suchlauf und Belege im Abschnitt Verifikation.
- **Zehn Slices in `done/`:** `slice-transformationen-spec-nachzug` (Pflichtenheft
  §1/§2/§4/§7, Architektur-Sicht, [`SPEC-030`](../../../../spec/pflichtenheft.md)),
  `slice-transformationen-kern-rename` (Regeltyp `rename_column`, Regelstand,
  Set/Remove, Nichtanwendbarkeits-Prüfung, Fitness-Tests), `slice-transformationen-antragsweg-schema`
  (Antragsarten `set_transformation`/`remove_transformation`, zwei SQL-Funktionen,
  Grants, Idempotenz-Guard, Rollen-Test), `slice-transformationen-antragsweg-usecase`
  (Use Cases mit K1–K4, Regelstand-Port, Store-Adapter, Verdrahtung),
  `slice-transformationen-backfill-pfad` (Regelauswertung im Backfill-Run,
  Fail-closed um den Regelstand, Fehlerklasse `schema` im Run),
  `slice-transformationen-map-value` (Regeltyp `map_value`, ohne Änderung an
  Antragsweg oder Wirkort), `slice-transformationen-e2e-wirkung` (Happy Path,
  Boundary, Neustart, Ausschluss+Regel am laufenden Feed-Container),
  `slice-transformationen-start-reihenfolge` (offene Anträge vor `stream.Run`
  verarbeitet, deterministische Ordnung ohne Datenbank geprüft),
  `slice-transformationen-e2e-abhilfe` (Nichtanwendbarkeit und Abhilfe-Kriterium
  (a)–(d) am laufenden System), `slice-transformationen-betriebsdoku`
  (Benutzerhandbuch, SDK-Beleg).
- **Neun wellenlose Kanten-Slices der Welle in `done/`** (Kanten in Welle-Datei
  §5): `slice-harness-suchlauf-nachmessen`, `slice-capture-leerlauf-quellbelege`,
  `slice-code-kommentare-kennungen`, `slice-harness-fmt-check`,
  `slice-antragsqueue-lesefehler-failed`, `slice-sdk-regel-realserver-e2e`,
  `slice-leerlauf-phase-last-in-stuecken`, `slice-wal-fehlerschwelle-ausgangsklasse`,
  `slice-start-vorlauf-grenze`. Zwei weitere Kanten-Slices bleiben offen
  (Folge-Slices unten): `slice-code-kommentare-bereinigung`,
  `slice-capture-transient-wiederholung` (letzterer gehört zur Welle
  `welle-backfill-bestand`, nicht zu dieser).
- **Sechs Folge-Entscheidungen der Welle:**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  (Mechanismus: geschlossener Regelsatz, zwei Antragsarten, dauerhaft aus den
  `applied`-Zeilen abgeleitet), [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  (Fehlerklasse `schema` im Backfill-Run, geteilt mit `welle-backfill-bestand`),
  [`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md)
  (Regelform als `json`- statt `jsonb`-Parameter),
  [`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md)
  (Annahmemenge von `rule_spec`), [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)
  (`requested_at` ist der Aufrufzeitpunkt, Ordnung der Verarbeitung) und
  [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  (Frist des Vorlaufs, `START_REPLICATION` als erste Handlung von `Stream.Run`).
  Die Architect-Verdikte stehen unter `docs/reviews/`:
  `architect-verdict-welle-transformationen-offene-fragen`,
  `architect-verdict-wal-fehlerschwelle-ausgangsklasse` und
  `architect-verdict-leerlauf-bestaetigung-intermittenz`.
- **Benutzerhandbuch:** Änderungshistorie 1.44 bis 1.68 (25 Zeilen, davon 1.67/1.68
  unmittelbar zur Welle: der Abschnitt „Transformationsregel konfigurieren“ trägt
  Aufrufform, Regeltypen, K1–K4, Wirkung, Informationsverlust bei `map_value`,
  Verhältnis zum Spaltenausschluss, Reihenfolge der Aufrufe, Dauerhaftigkeit samt
  Grenze bei einem älteren Binärstand, Nichtanwendbarkeit samt Abhilfe im
  Erfassungspfad und im Backfill-Run, Backfill-Bezug, Wartezeit des Stream-Starts,
  Kosten der Lesung im Backfill, Form auf allen Zustellwegen; §2 nennt die beiden
  Funktionen, §6 die Zeile `schema`, §8 den Begriff „Transformationsregel“).
- **Schema:** `administration_request` trägt zwei neue nullable Spalten
  (`rule_name text`, `rule_spec jsonb`); zwei neue SQL-Funktionen
  (`cdc.set_transformation`, `cdc.remove_transformation`, `EXECUTE` allein für
  `cdc_admin`); `request_kind` wächst von fünf auf sieben Werte (belegt am
  Alt-Tag-Lauf gegen `v0.2.0`, Abschnitt Verifikation).
- **RTM:** `make doc-trace` druckt `80 Anforderung(en), 1 Waise(n).` —
  [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) trägt jetzt den Nachweis
  `E2E, SDK-E2E` und ist keine Waise mehr; die verbleibende Waise ist
  [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing, ausdrücklich
  Out-of-Scope dieser Welle, §6 der Welle-Datei).

## Was hat funktioniert?

- **Die eine gemeinsame Row-Image-Funktion (K1 aus `welle-backfill-bestand`)
  trug die ganze Welle ohne Bruch.** `model.BuildRowImage` blieb über zehn
  Slices die einzige Konstruktionsstelle für WAL- und Backfill-Pfad; kein
  Slice musste sie verzweigen.
- **Architect-Verdikte an den Stellen, an denen eine Messung oder ein Code-Fund
  dem Plan widersprach, hielten die Welle in Bewegung statt sie anzuhalten.**
  Die Abweichung vom ADR-Wortlaut (Start-Reihenfolge als eigener Slice statt
  Teil des E2E-Slice) wurde als Frage geführt und im Architect-Verdikt
  `architect-verdict-welle-transformationen-offene-fragen` entschieden; derselbe
  Verdikt entschied vier offene Fragen ((a) Host-Werkzeug/Umleitung, (b)
  Lesekosten des Backfill-Runs als akzeptiertes Negativ mit Messung, (d) Obergrenze
  der `map_value`-Paare als benannter Verzicht mit Trigger, (e) Zeitgrenze des
  Vorlaufs über `ADR-0128` und `slice-start-vorlauf-grenze`) und die Restfläche der
  Zustellwege (Neu-Bild einer INSERT auf fünf Wegen erprobt, der Rest hergeleitet).
  Zwei weitere Architect-Verdikte während der Welle
  (`architect-verdict-wal-fehlerschwelle-ausgangsklasse`,
  `architect-verdict-leerlauf-bestaetigung-intermittenz`) lieferten je einen
  Kanten-Slice, ohne die Reihenfolge der zehn Kern-Slices zu unterbrechen.
- **Der Alt-Tag-Lauf (`tools/harness/run-schema-rollout-guard-test.sh` Lauf 5)
  fängt eine Schema-Regression real, nicht nur am Diff.** Der Lauf gegen `v0.2.0`
  zeigt real die Vorher-Form (`rule_name`/`rule_spec` NULL, fünf `request_kind`-Werte)
  neben der Nachher-Form desselben Arbeitsbaums — ein Beleg, der nicht an einer
  Behauptung im Plan hängt.

## Was ging anders als geplant?

- **Die Welle wuchs um neun wellenlose Kanten-Slices**, die beim Öffnen (2026-09-23)
  noch nicht alle absehbar waren: `slice-wal-fehlerschwelle-ausgangsklasse` (ein
  Codefehler — Ausgangsklasse `storage` statt `replication`), `slice-start-vorlauf-grenze`
  (Umsetzung von `ADR-0128`, aus dem Re-Evaluierungs-Trigger 4 von `ADR-0112`
  gefolgert) und `slice-leerlauf-phase-last-in-stuecken` (ein intermittent rotes
  CI-Ergebnis der Leerlauf-Bestätigungs-Phase, Ursache im Testaufbau, nicht im
  Produkt). Konsequenz: die Roadmap trägt sechs Drift-Log-Zeilen zwischen
  2026-09-25 und 2026-09-27, die jede Kante einzeln benennen.
- **Die Restfläche der Zustellwege bleibt teilweise hergeleitet, nicht erprobt.**
  `make test-integration` belegt die Form „alle fünf Wege dieselbe“ für das
  Neu-Bild einer eingefügten Zeile; UPDATE, DELETE und das Alt-Bild sind über
  `cdc.changes` belegt und auf den vier weiteren Wegen hergeleitet (alle Wege
  lesen dieselbe Change, [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6 Option D). Kein Folge-Slice — Architect-Verdikt
  `architect-verdict-welle-transformationen-offene-fragen` §4 akzeptiert die
  gemessene Menge samt Kennzeichnung als hergeleitet für den Rest.
- **Die Lesekosten des Backfill-Runs sind ein akzeptiertes Negativ mit Messung,
  nicht mit einer Umsetzung geschlossen.** Bis rund 10.000 Zeilen der
  Antrags-Queue liegt der Mehraufwand je Lesung bei etwa 8 % einer Blockdauer,
  bei 100.000 Zeilen bei etwa 70 % (synthetische Messung, keine Queue realer
  Größe vorhanden); Trigger: mehr als 10.000 Zeilen in der Queue einer Quelle
  (Architect-Verdikt §3).
- **Zwei Einträge des Beobachtungs-Registers erreichten während der Welle die
  3×-Schwelle, ohne dass ein Architect-Zug ihren Ausgang schon entschieden
  hätte** — anders als die vier Einträge, die
  `architect-verdict-welle-transformationen-offene-fragen` §8 bereits abschließend
  liest. Diese Closure liest sie (unten, Lese-Schritt), entscheidet ihren Ausgang
  aber nicht selbst: beide state.md-Dateien benennen ausdrücklich einen
  Architect-Zug als nächsten Schritt („Planner → Architect → Planner“ bzw. „der
  Architect entscheidet … ob ein … Wächter … eine Entscheidung braucht“), und ein
  Planner, der diese Entscheidung an ihrer Stelle selbst träfe, wäre derselbe
  Kontext, der sie prüfen soll (Rollentrennung, Baseline-Regelwerk
  `modul-08-agentenrollen.md`). Konsequenz: beide bleiben mit Ausgang
  „gelesen, Architect-Entscheidung aussstehend“ offen, kein Folge-Slice, keine
  neue ADR in diesem Zug — Adresse ist ein eigener Architect-Zug nach dieser
  Closure (siehe Lese-Schritt unten und Bericht dieser Closure an den
  Auftraggeber).

## Steering-Loop-Einträge

Die vier Einträge, deren Ausgang `architect-verdict-welle-transformationen-offene-fragen`
§8 bereits liest, tragen ihren Anker dort und werden hier nicht wiederholt:
`test-name-behauptet-mehr-als-der-test-treibt` (gestrichen, akzeptiertes
Negativ), `gemeldete-ungenauigkeit-ohne-traeger` (verkörpert →
`.harness/skills/closure-note-reviewer.md`), `plan-zusage-erfuellung-ohne-committeten-anker`
(gestrichen, akzeptiertes Negativ), `adapter-unittest-verdeckt-bootstrap-luecke`
(geplant → `slice-start-vorlauf-grenze`, jetzt in `done/`).

Diese Closure selbst liefert keinen neuen verkörperten Eintrag — jeder von ihr
gelesene Eintrag, dessen Zähler während dieser Welle wuchs, trug bereits einen
zugewiesenen Ausgang aus einer früheren Welle (`nachzug-laesst-ueberholten-text-stehen`,
`kommentar-behauptet-nicht-getragenen-fehlerpfad`,
`db-gegenstand-enthaelt-netzlos-geprueften-code`, `rollen-test-abdeckungsluecken`,
u. a. — Liste im Lese-Schritt unten), außer den zwei Kandidaten, deren Ausgang
ausdrücklich einen eigenen Architect-Zug verlangt (`aufschub-adresse-nimmt-sendung-nicht-an`,
`ein-instanz-annahme-ohne-erzwingung`) und die deshalb **nicht** in diesem Zug
verkörpert werden.

### Lese-Schritt des Beobachtungs-Registers

Das Register führt 131 Verzeichnisse
(`ls docs/plan/planning/observations/BEO-PGC | wc -l`, gemessen am Stand
`ed1ef980`), 48 davon mit `evidence/` ab 3 Dateien. 33 Einträge tragen eine
Beleg-Datei eines Slices dieser Welle oder ihrer neun wellenlosen Kanten-Slices
(Suchlauf: `find … -iname "<slice>.md" -path "*/evidence/*"` je Slice, Ergebnis
zusammengefasst).

Bereits vor dieser Closure mit zugewiesenem Ausgang (unverändert, kein
Architect-Zug nötig): `adapter-unittest-verdeckt-bootstrap-luecke` (5×,
verkörpert/geplant, s. o.), `adr-aussage-breiter-als-ihre-messung` (10×,
verkörpert), `antrag-mit-leerem-regelnamen-stallt-die-queue` (2×, verkörpert),
`arbeit-ueberholt-stehenden-traeger` (33×, verkörpert, Deckel),
`beleg-befehl-traegt-seinen-satz-nicht` (17×, verkörpert),
`d-migrate-nacharbeit` (8×, verkörpert), `dod-begruendung-unzutreffende-tatsachenbehauptung`
(11×, verkörpert), `ersatzweg-nach-verweigerter-aktion` (1×, verkörpert),
`formatierungs-drift-ohne-gate` (3×, verkörpert), `gemeldete-ungenauigkeit-ohne-traeger`
(3×, verkörpert, s. o.), `host-werkzeug-jenseits-docker-und-make-ohne-deklaration`
(2×, verkörpert), `inplace-textwerkzeug-am-repo-trotz-nutzerregel` (8×,
verkörpert), `kommentar-behauptet-nicht-getragenen-fehlerpfad` (7×, verkörpert),
`kommentar-herkunft-als-kette` (2×, verkörpert), `nachzug-laesst-ueberholten-text-stehen`
(15×, verkörpert), `negativtest-ohne-bindung-an-seine-eingabe` (20×, verkörpert),
`plan-zusage-erfuellung-ohne-committeten-anker` (5×, gestrichen, s. o.),
`rollen-test-abdeckungsluecken` (4×, gestrichen), `slice-pfad-als-link-in-berichten`
(5×, verkörpert), `test-integration-retention-timing-flake` (5×, verkörpert),
`test-name-behauptet-mehr-als-der-test-treibt` (4×, gestrichen, s. o.),
`vorher-nachher-sprache-in-test-harness-kommentar` (7×, verkörpert),
`wartegrenze-ohne-zeitgrenze-im-startpfad` (2×, verkörpert),
`zahl-in-traeger-driftet-gegen-die-messung` (25×, verkörpert),
`zitat-nennt-die-falsche-stelle` (9×, verkörpert).

Unter der 3×-Schwelle, mit Adresse an eine spätere Closure statt an diese:
`alt-tag-lauf-vorbedingung-am-juengsten-tag` (1×), `implementierung-weicht-von-adr-wortlaut-ab`
(2×, die Fragen selbst sind bereits als Fragen geführt und beantwortet — die
Beobachtung über den Musterbruch bleibt eigenständig unter der Schwelle),
`queue-ordnung-transaktionsbeginn-und-zufallskennung` (1×),
`spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (1×),
`werkzeug-fuehrt-plan-inhalt-als-argument-aus` (1×),
`werkzeug-liest-nutzerkonfiguration-ohne-pin` (1×),
`wertabhaengiger-zweiter-waechter-ohne-spec-zeile` (1×). Ebenfalls unter der
Schwelle und in der Welle-Datei §6 benannt, unverändert: `test-runner-stiller-ausschluss`
(2×), `vorab-bedingung-nach-umsetzung-geprueft` (2×),
`kein-admin-weg-schema-fehler-recovery` (1×), `adr-folgepflicht-ohne-traeger-slice`
(1×, die Pflicht ist eingelöst).

**Zwei Einträge bei 3× oder darüber ohne Ausgang, den diese Closure selbst
zuweisen könnte:**

| Eintrag (`BEO-PGC/…`) | Zähler | Zustand | Adresse |
|---|---|---|---|
| `aufschub-adresse-nimmt-sendung-nicht-an` | 5 | geplant — Ausgang-Vorschlag liegt vor (zwei Zielorte: `.claude/commands/implement-slice.md` Schritt 17, `.claude/commands/plan-welle.md`, `.harness/skills/reviewer.md` HIGH „Neue Betreiber-Oberfläche ohne Handbuch-Zug“), der Architect entscheidet ihn nicht in diesem Zug | eigener Architect-Zug nach dieser Closure |
| `ein-instanz-annahme-ohne-erzwingung` | 3 | offen — der dritte Beleg (`slice-start-vorlauf-grenze`) verschiebt den Zeitpunkt, an dem eine zweite Instanz derselben Quelle am Slot scheitert; die `state.md` benennt ausdrücklich, dass der Architect am Lese-Schritt entscheidet, ob ein Ein-Instanz-Wächter für den Prozessstart eine Entscheidung braucht | eigener Architect-Zug nach dieser Closure |

**Feststellung:** Diese Closure hat beide Einträge gelesen und ihren Zustand
mit Zähler und Adresse festgehalten; sie hat **keinen** der beiden Ausgänge
selbst gesetzt, weil ihre eigenen `state.md`-Dateien einen Architect-Zug als
nächsten, noch ausstehenden Schritt benennen und ein Planner, der ihn selbst
setzte, dieselbe Rolle wäre, die die eigene Entscheidung prüfen soll. Alle
übrigen Einträge bei 3× oder darüber, die diese Welle berührt, tragen einen
Ausgang.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/)
(`BEO-PGC/<slug>/evidence/`); er wird nicht in dieser Notiz gepflegt. Was in
dieser Welle 3× erreicht und einen Ausgang bekommen hat, steht oben; die zwei
Einträge ohne zugewiesenen Ausgang stehen in der Tabelle des Lese-Schritts.

## Folge-Slices

Zwei Kanten-Slices der Welle bleiben offen (Datei in `docs/plan/planning/open/`):

- `slice-code-kommentare-bereinigung` — kürzt Kommentare in den Testdateien, die
  `e2e-wirkung`/`e2e-abhilfe` erweitert haben, und schreibt
  `docs/user/e2e-abdeckung.md` neu; die Welle wartet nicht auf ihn (Welle-Datei
  §5).
- `slice-capture-transient-wiederholung` — gehört zur Welle `welle-backfill-bestand`,
  nicht zu dieser; genannt, weil sein Start-Trigger (Architect-Entscheidung zur
  Wiederholungsform) unverändert offen ist.

Kein neuer Folge-Slice entsteht aus dieser Closure selbst — die zwei
Register-Einträge ohne Ausgang (Lese-Schritt oben) brauchen zunächst einen
Architect-Zug, keinen Slice-Schnitt.

## Verifikation

Alle Zeilen am Stand `ed1ef980` (Docker-Läufe dieser Closure, 2026-09-27),
sofern nicht anders angegeben; Gate-Exit-Codes ungefiltert gesichert
([`AGENTS.md`](../../../../AGENTS.md) §3.9).

### Schritt 1 — Trigger

| Kriterium (Welle-Datei §3) | Beleg |
|---|---|
| Alle zehn Slices in `done/` | `ls docs/plan/planning/done \| grep -c '^slice-transformationen-'` druckt `10`; die zehn Namen stehen im Abschnitt „Was wurde geliefert?“ |
| `make gates` grün | Exit `0`; `generated-sync: OK`, a-check `gesamt: 0 Befund(e)`, `coverage-gate: OK — Coverage 85.40% erfüllt Schwelle 80%`; Gate-Stempel `.harness/state/gates-passed.diffsha` gleicht `tools/harness/working-tree-hash.sh` |
| Realer, grüner `make test-integration`-Lauf mit den Transformations-Belegen | Ein Lauf, real und grün, bereits mit `set -euo pipefail` in `run-integration-tests.sh`: gedruckt „Transformationen-Happy-Path“ (fünf Zustellwege, gRPC/SSE/NATS je mit `change_id`), „Transformationen-Neustart und Ausschluss“ (zwei reale `docker restart`), „Leerlauf-Bestätigung“ (Backfill-Bestand über 30.000 Zeilen mit transformierter Form), „Prozessstart-Vorlauf-Frist“, `TestE2ETransformationRuleNotApplicableEndsCaptureWithSchemaClass` und „Transformationen-Nichtanwendbarkeit und Abhilfe“ (Antrag `pending` während des Stillstands, `applied` nach Neustart vor der ersten Transaktion, Rohform ohne zweiten `schema`-Fehler); Lauf endet mit „Lauf abgeschlossen — E2E-Abdeckungstabelle aus 17 Go-Zeilen und 41 Bash-Zeilen“, `docs/user/e2e-abdeckung.md` unverändert |
| CI-Stand des Push | `gh run list --limit 6`: `ci` (`36350963628`), `e2e` (`36350963623`), `examples` (`36350963622`) für den Push auf `ed1ef980` je `success` |
| `make test-store` grün | Exit `0`; `DB-Adapter-Coverage: 82.77% (gedeckt 889 von 1074 Statements; Profile gemergt: store,replication)`, `db-coverage: OK` |
| `make test-replication` grün | Exit `0`, Modi `measure` und `tier` beide durchgelaufen; `TestSourceKeepaliveInsideTransactionDeliversWholeTransaction` PASS (PostgreSQL 18.6), `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot` PASS, `TestSlotReserveExhaustedIsConfiguration` PASS |
| `make schema-rollout` zweimal und Alt-Tag-Lauf | `bash tools/harness/run-schema-rollout-guard-test.sh` Exit `0`; gedruckt: „Lauf 5 OK — Tag v0.2.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); … Spalten rule_name:text:YES,rule_spec:jsonb:YES (Alt-Zeile alttag-req: NULL), request_kind-Menge backfill,disable,enable,exclude_column,include_column,remove_transformation,set_transformation …“ und „OK — alle Belege real erbracht (…, Alt-Tag v0.2.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)“ |
| Fitness-Function-Tests (`ADR-0112`) grün | `make test` (`-race`, ganzer Baum) Exit `0`, keine `FAIL`-Zeile; `TestRenameColumnImagesAreDeterministic`, `TestMapValueImagesAreDeterministic`, `TestAssemblerTransformationsAreRaceFree`, `TestAssemblerMapValueIsRaceFree`, `TestAssemblerColumnExclusionIsRaceFree`, `TestAssemblerLiveReloadIsRaceFree` liefen im selben Lauf mit; `make a-check` (Teil von `make gates`) `gesamt: 0 Befund(e)` |
| `make doc-trace` | Exit `0`; Zeile `LH-FA-CFG-007 … ADR-0112, ADR-0117, ADR-0125, ADR-0126, ADR-0127, ADR-0128 … E2E, SDK-E2E … ok`, Zeile `LH-FA-CFG-008 … ADR-0112 … WAISE`, Schluss `80 Anforderung(en), 1 Waise(n).` |
| Restfläche der Zustellwege entschieden | `architect-verdict-welle-transformationen-offene-fragen` §4 — kein Folge-Slice, Aussage mit gemessener Menge (INSERT, Neu-Bild, fünf Wege) und „hergeleitet“ für den Rest |
| Regelauswertung an genau einer Stelle | `git grep -n 'BuildRowImage(' -- '*.go' \| grep -v _test.go` druckt drei Aufrufer-Zeilen (zwei in `mapper.go` für `new`/`old`, eine in `backfill/service.go`) und die Deklaration in `rowimage.go` — keine zweite Konstruktionsstelle |
| Lesekosten des Backfill-Runs entschieden | `architect-verdict-welle-transformationen-offene-fragen` §3 — akzeptiertes Negativ mit Messung (≈8 % bei 10.000 Zeilen, ≈70 % bei 100.000 Zeilen, synthetisch) und Trigger (>10.000 Zeilen in der Queue einer Quelle) |
| Zeitgrenze des Vorlaufs entschieden | [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md) `Accepted`; `slice-start-vorlauf-grenze` in `done/` |
| Handbuch-Träger | Änderungshistorie bis Version 1.68; „Transformationsregel konfigurieren“ trägt Aufrufform, K1–K4, Wirkung, Dauerhaftigkeit samt Grenze, Nichtanwendbarkeit samt Abhilfe im Erfassungspfad **und** im Backfill-Run, Backfill-Bezug; SDK-Beleg (Row-Image-Schlüssel in den drei SDK-Modellen opak) gestützt auf `slice-sdk-regel-realserver-e2e` (in `done/`) |
| Lese-Schritt des Beobachtungs-Registers | oben, Abschnitt „Steering-Loop-Einträge“ § Lese-Schritt — gelaufen, zwei Einträge ohne zugewiesenen Ausgang benannt und adressiert |
| Closure-Notiz | diese Datei |

### Schritt 2 — Trigger-Audit

- **Carveouts: 0 offen.** `find docs -iname "CO-*.md"` und eine Suche nach
  `carveout` außerhalb der vendored Baseline liefern keinen Treffer — dieses
  Repo hat aktuell keinen offenen Carveout.
- **Bootstrap-aware Gates.**
  - *Unit-Gate:* `harness/mk/coverage.mk` führt `THRESHOLD ?= 80` (Endstufe,
    Rampe ausgeschöpft), der Gate-Lauf druckt `coverage-gate: OK — Coverage
    85.40% erfüllt Schwelle 80%`. **0 offen.**
  - *DB-Adapter-Coverage:* `tools/harness/db-coverage.sh` führt
    `DB_COVERAGE_THRESHOLD=${DB_COVERAGE_THRESHOLD:-80}` (Endstufe, bereits mit
    `welle-backfill-bestand` hochgeschaltet); der gemergte Lauf dieser Closure
    druckt `82.77%` gegen Schwelle `80`. **0 offen.**
- **ADR-Re-Evaluierungs-Trigger** (gelesen: §Re-Evaluierungs-Trigger jeder der
  sechs Folge-ADRs dieser Welle):
  - [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md):
    kein Consumer verlangt eine abweichende Form je Zustellweg; kein dritter
    Regeltyp erreichte im Register 3× denselben Bedarf; die Routing-ADR
    (Teilfrage 7) existiert nicht; Trigger 4 (Abhilfe wirkt real nicht ohne
    Eingriff in die Startreihenfolge) **ist eingetreten und ist mit
    `slice-transformationen-start-reihenfolge` und
    [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)/`slice-start-vorlauf-grenze`
    abgearbeitet**. **0 offen.**
  - [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md): kein
    neuer Trigger seit der Closure von `welle-backfill-bestand` eingetreten
    (gelesen dort, unverändert). **0 offen.**
  - [`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md),
    [`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md),
    [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md):
    je gelesene Trigger-Zeile nicht eingetreten (kein d-migrate-Umbau auf
    `jsonb`-Abbau beobachtet, keine Betreiber-Meldung). **0 offen.**
  - [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md):
    kein Betreiber hat die Frist als zu lang/kurz gemeldet, keine weitere
    synchrone Wartestelle entstanden, `Stream.Run` sendet `START_REPLICATION`
    unverändert als erste Handlung. **0 offen.**
- **Beobachtungs-Register:** Lese-Schritt oben — zwei Einträge bei 3×/darüber
  ohne selbstgesetzten Ausgang, adressiert an einen eigenen Architect-Zug (kein
  stiller Verbleib: Zustand und Adresse stehen in der Tabelle oben).

### Schritt 4 — Archivierung

Das Repo führt kein eigenes Archivierungs-Werkzeug: kein Make-Ziel in
`Makefile` oder `harness/mk/*.mk` und kein Skript unter `tools/`. Das externe
Werkzeug `ai-harness-init archive-welle` liegt lokal vor (im PATH dieser
Umgebung); sein Vorschau-Lauf (`--vorschau welle-transformationen`, schreibt
nichts) meldete vor dieser Closure zwei Sperren — die fehlende Ergebnisnotiz
und den noch nicht verschobenen Welle-Plan, beide Vorbedingungen von Schritt 3
dieser Prozedur. Der schreibende Lauf ist nicht Teil dieser Closure: er
committet mit fest einprogrammierten Messages ohne `LH-*`/`ADR-*`-Kennung
(`BEO-PGC/externes-werkzeug-committet-ohne-kennung`), dieselbe Bedingung, die
die Closures von `welle-backfill-bestand` und `welle-nats-drittstream` als
nicht eingetreten festhalten. Die Bedingung des Schrittes (ein im Repo
geführtes Werkzeug) ist nicht eingetreten; keine Handarchivierung. Die
Entscheidung, den Lauf des externen Werkzeugs für diese Welle gleichwohl
auszuführen, liegt beim Betreiber.

### Docker-Lauf-Hygiene

Die Docker-Läufe dieser Closure (`make gates`, `make test-integration`,
`make test-store`, `make test-replication`, `make test`, der
Schema-Rollout-Guard-Test) liefen nacheinander, nicht gleichzeitig; jeder Lauf
endete grün, bevor der nächste startete.

### Ausgang des Closure-Note-Reviews

Ein eigener Review-Lauf dieser Closure-Notiz in frischem Kontext
(`.harness/skills/closure-note-reviewer.md`) liegt zum Zeitpunkt dieses
Schreibens **nicht** vor — anders als bei `welle-backfill-bestand`. Diese
Abweichung ist benannt, nicht verschwiegen: die Prüfung obliegt der nächsten
Rolle im Workflow, nicht diesem Planner-Zug selbst (kein Self-Review, Modul 8).
