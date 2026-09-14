# ADR-0063: `LH-FA-SCH-003`-Testform-Korrektur — Supersedes ADR-0058 (nur Entscheidung 1)

**Status:** Accepted — Supersedes [`ADR-0058`](0058-testansatz-fuenf-luecken.md)
(nur deren Entscheidung 1, „`LH-FA-SCH-003` — Entfernte Spalten"; die
Entscheidungen 2–5 dieser ADR — `LH-FA-DAT-006`, `LH-QA-OPS-005`,
`LH-QA-POR-001`, `LH-QA-POR-002` — bleiben unverändert bestehen und werden
hier nicht wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-062` <!-- d-check:status-provenance -->,
der den Widerspruch fand und korrekt stoppte, statt ihn selbst
aufzulösen, Modul 8 §Konflikt-Pfad)

**Bezug:** [`LH-FA-SCH-003`](../../../spec/lastenheft.md) (korrigierte
Kennung), [`LH-FA-SCH-004`](../../../spec/lastenheft.md) (Konvergenz-Ziel),
[`ADR-0058`](0058-testansatz-fuenf-luecken.md) (korrigierte Entscheidung 1),
[`ADR-0059`](0059-spaltenauswahl-mechanismus.md) Teilfrage 4 (bereits
akzeptierter Präzedenzfall, selbes Datum — beantwortet die hier korrigierte
Frage inhaltlich bereits, nur für den Spaltenausschluss-Fall statt den
direkten Fall der realen Spaltenentfernung selbst),
`docs/plan/planning/done/slice-033-typauswertung-fehlerklasse-schema.md` <!-- d-check:status-provenance -->
(Bestandsbeleg: derselbe `relationOther`/`ErrIncompatibleSchemaChange`-Pfad
bereits real gegen PostgreSQL verifiziert, für den Typänderungs-Fall),
`internal/adapters/driving/replication/mapper/mapper.go`
(`classifyRelationColumns`, unverändert bestehender Code — nicht Gegenstand
dieser ADR), `docs/plan/planning/in-progress/slice-062-e2e-schema-drop-column-metadaten-erweiterbarkeit.md` <!-- d-check:status-provenance -->
(Umsetzungs-Slice, findet und meldet den Widerspruch)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0058`; trifft eine
Testmethoden-Korrektur, ändert keine Lastenheft-/Pflichtenheft-Zusage)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0058`](0058-testansatz-fuenf-luecken.md) (Accepted, 2026-09-14)
legt in Entscheidung 1 für `LH-FA-SCH-003`s Happy Path fest: „die danach
eingefügte Zeile erscheint im Row Image ohne die entfernte Spalte — kein
Platzhalter, kein Fehler." Der `slice-062`-Implementer-Lauf <!-- d-check:status-provenance -->
schrieb `TestE2ESchemaChangeDropColumn` exakt nach dieser Prämisse und
fand beim realen Testlauf gegen den Compose-Feed-Container das Gegenteil:
ein reales `ALTER TABLE … DROP COLUMN` löst `relationOther` aus
(`internal/adapters/driving/replication/mapper/mapper.go`,
`classifyRelationColumns`) und damit `mapper.ErrIncompatibleSchemaChange`
— derselbe sichtbare `schema`-Fehler, der den Erfassungspfad beendet
(`Assembler.Consume` → `receive.Stream.Run` → `bootstrap.Run` →
`os.Exit(1)`, `restart: "no"` in `compose.yaml`), den auch eine
inkompatible Typänderung auslöst (`LH-FA-SCH-004`, real getestet in
`TestE2ESchemaChangeIncompatibleTypeChange`).

Das ist **kein neu eingetretenes Verhalten**, sondern bereits akzeptierter
Bestand, an zwei Stellen unabhängig dokumentiert:

1. **`classifyRelationColumns` selbst** (unverändert seit `slice-033` <!-- d-check:status-provenance -->,
   `done/`) unterscheidet nicht zwischen „Spalte real gelöscht" und „Spalte
   umbenannt/Typ geändert" — jede bekannte Spalte, die in der eingehenden
   `decode.Relation` fehlt oder eine andere Typ-OID trägt, fällt in denselben
   `relationOther`-Zweig. Der Funktionskommentar in `mapper.go` benennt das
   wörtlich: „Spalte entfernt, Typ einer bestehenden Spalte geändert, Spalte
   umbenannt — jede nicht sicher als reine Obermengen-Erweiterung erkennbare
   Relation-Änderung." `slice-033` <!-- d-check:status-provenance --> hat diesen Pfad bereits real gegen
   PostgreSQL verifiziert (für den Typänderungs-Fall) — die Pfad-Semantik
   ist damit keine neue Erkenntnis, nur bislang nicht auf den direkten
   Spaltenentfernungs-Fall angewendet worden.
2. **`ADR-0059`** (Accepted, selbes Datum wie `ADR-0058`) beantwortet in
   Teilfrage 4 exakt diese Frage — für den Spaltenausschluss-Fall
   (`LH-FA-CFG-005`), aber mit einer Begründung, die den direkten Fall der
   realen Spaltenentfernung selbst einschließt: „eine real gelöschte …
   Spalte fehlt in der eingehenden `decode.Relation` genauso wie jede andere
   gelöschte Spalte und löst denselben `relationOther`-Pfad mit
   `ErrIncompatibleSchemaChange` aus wie jede reguläre Spaltenlöschung
   (`LH-FA-SCH-004.a`) — kein Sonderfall." `ADR-0059` Teilfrage 4 hat diese
   Konvergenz damit bereits akzeptiert, ohne dass `ADR-0058` Entscheidung 1
   — geschrieben im selben Architect-Lauf, am selben Tag — sie für ihre
   eigene Testform übernommen hätte. Die beiden Entscheidungen widersprechen
   sich, weil derselbe Fakt in zwei Teilen derselben Architect-Sitzung
   unterschiedlich behandelt wurde.

Nach `AGENTS.md` §3.5 / Modul 8 §Rollen-Regeln ist eine Accepted-ADR
immutable; der Implementer durfte `ADR-0058`s Klausel nicht stillschweigend
umschreiben, um seinen real roten Test grün zu bekommen — das wäre Drift,
kein pragmatisches Implementieren. Der korrekte Pfad ist Verdikt 2 aus
Modul 8 §Konflikt-Pfad: Folge-ADR mit `Supersedes`, hier bestätigt durch
einen unabhängigen Architect-Zug (anderer Kontext als der
Implementer-Lauf, der den Fund machte).

## Entscheidung

Wir korrigieren `ADR-0058`s Entscheidung 1, ausschließlich die
Happy-Path-Prämisse: **Der Happy Path für `LH-FA-SCH-003`s E2E-Beleg ist
nicht stille Spalten-Auslassung, sondern Konvergenz mit `LH-FA-SCH-004`** —
eine real gelöschte Spalte löst denselben sichtbaren `ErrIncompatibleSchemaChange`-
Fehlerpfad aus wie jede andere nicht sicher als Obermenge erkennbare
Relation-Änderung, bestätigt durch den bereits akzeptierten Präzedenzfall
`ADR-0059` Teilfrage 4. Die **Boundary-Klausel aus `ADR-0058`
Entscheidung 1 bleibt unverändert** — sie war nicht Teil des Widerspruchs
und bekommt hier ihren eigenen, bislang fehlenden Testbeitrag (siehe
§Testform unten).

### Trägt `LH-FA-SCH-003`s eigener Wortlaut die Konvergenz noch?

[`LH-FA-SCH-003`](../../../spec/lastenheft.md)s Happy Path verlangt
wörtlich nur: „ist ihr Verhalten definiert (dokumentiert)" — kein
bestimmtes Verhalten, nur *ein* definiertes. Die Konvergenz erfüllt das:
Das Verhalten bei einer real entfernten Spalte ist exakt definiert (fällt
in `relationOther`, meldet `ErrIncompatibleSchemaChange`, beendet den
Erfassungspfad) — bloß nicht das Verhalten, das `ADR-0058` fälschlich
angenommen hatte. Ein eigener, unterscheidbarer Test für den Happy Path
(über den bereits von `TestE2ESchemaChangeIncompatibleTypeChange`
belegten Mechanismus hinaus) ist **nicht nötig** — beide Fälle durchlaufen
denselben Code-Pfad mit derselben Konsequenz, ein zweiter Beleg desselben
Mechanismus würde nichts Neues prüfen.

Die **Boundary-Klausel** liegt anders: „ist definiert, ob und in welcher
Form ihre historischen Werte lesbar bleiben." `TestE2ESchemaChangeIncompatibleTypeChange`
prüft diesen Aspekt **nicht** — der Test verifiziert nur, dass die *nach*
der Typänderung eingefügte Zeile (id=11) nicht erfasst wird; er liest keine
*vor* der Änderung bereits erfasste Zeile erneut, um ihre unveränderte
Lesbarkeit samt historischem Wert zu bestätigen. `LH-FA-SCH-003`s
Boundary-Klausel bleibt damit ohne einen eigenständigen Testbeitrag durch
`TestE2ESchemaChangeDropColumn` faktisch unbelegt — dieser Teil der
Testfunktion trägt echten, nicht redundanten Wert.

## Testform: `TestE2ESchemaChangeDropColumn` (korrigiert)

Struktur bleibt wie im bisherigen Implementer-Diff (eigene, wegwerfbare
Spalte `removable`, getrennt von `amount`/`extra`), die Assertions ändern
sich:

```go
func TestE2ESchemaChangeDropColumn(t *testing.T) {
	env := newE2EEnv(t, "feed_e2e_schema")
	ctx := context.Background()

	if _, err := env.pool.Exec(ctx, "ALTER TABLE "+env.feed+" ADD COLUMN removable text"); err != nil {
		t.Fatalf("ALTER TABLE ADD COLUMN (removable): %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name, removable) VALUES (3, 'BeforeDrop', 'ToBeRemoved')"); err != nil {
		t.Fatalf("INSERT vor der Entfernung: %v", err)
	}
	beforeRows := awaitChangesViewRows(t, env, "3", 1)

	if _, err := env.pool.Exec(ctx, "ALTER TABLE "+env.feed+" DROP COLUMN removable"); err != nil {
		t.Fatalf("ALTER TABLE DROP COLUMN: %v", err)
	}
	if _, err := env.pool.Exec(ctx,
		"INSERT INTO "+env.feed+" (id, name) VALUES (4, 'AfterDrop')"); err != nil {
		t.Fatalf("INSERT nach der Entfernung: %v", err)
	}

	// LH-FA-SCH-003 Happy Path (ADR-0063, Supersedes ADR-0058 Entscheidung 1):
	// die reale Spaltenentfernung löst denselben ErrIncompatibleSchemaChange-
	// Pfad aus wie jede andere inkompatible Relation-Änderung (Konvergenz
	// mit LH-FA-SCH-004, bereits von ADR-0059 Teilfrage 4 akzeptiert) —
	// sichtbarer schema-Fehler, kein stilles Auslassen.
	if got := awaitHeartbeatErrorClass(t, env, "schema"); got != "schema" {
		t.Fatalf("cdc.heartbeat.error_class nach DROP COLUMN: %q, wollen \"schema\"", got)
	}
	rows, err := queryChangesView(ctx, env, "4")
	if err != nil {
		t.Fatalf("cdc.changes-Lesung für id=4: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("cdc.changes trägt id=4 nach dem gemeldeten schema-Fehler: %+v (Erwartung: der Erfassungspfad endete vor dem Commit dieser Transaktion)", rows)
	}

	// LH-FA-SCH-003 Boundary (unverändert ggü. ADR-0058): die vor der
	// Entfernung erfasste Change bleibt über cdc.changes unverändert lesbar,
	// inklusive des historischen Werts der entfernten Spalte.
	rereadBefore := awaitChangesViewRows(t, env, "3", 1)
	if string(rereadBefore[0].newData) != string(beforeRows[0].newData) {
		t.Fatalf("ältere Change nach DROP COLUMN: %s (vor der Entfernung: %s)", rereadBefore[0].newData, beforeRows[0].newData)
	}
	beforeImage := imageJSON(t, rereadBefore[0].newData)
	if beforeImage["removable"] != "ToBeRemoved" {
		t.Fatalf("ältere Change trägt den historischen Wert der entfernten Spalte nicht: %s", rereadBefore[0].newData)
	}
}
```

Kein Platzhalter-/Präsenz-Assertion mehr auf `afterImage["removable"]` —
diese Prämisse entfällt vollständig, weil es kein `afterImage` mehr gibt
(die Zeile id=4 wird nie erfasst).

### Eigene vs. gemeinsame Testfunktion

`TestE2ESchemaChangeDropColumn` bleibt eine **eigene** Testfunktion, wird
**nicht** mit `TestE2ESchemaChangeIncompatibleTypeChange` zusammengelegt
oder parametrisiert — aus zwei Gründen, nicht ohne das Duplikations-Risiko
zu verschweigen:

1. **Traceability pro Kennung.** `LH-FA-SCH-003` und `LH-FA-SCH-004` sind
   zwei getrennte Lastenheft-Kennungen mit eigenen Akzeptanzkriterien
   (`AGENTS.md` §5: „Neue oder geänderte Anforderungen brauchen einen
   Beleg"). Eine gemeinsame, parametrisierte Funktion verwischt die
   1:1-Zuordnung Kennung → Testfunktion, die `ADR-0058`s bisherige
   Fitness-Function-Tabelle bereits trägt — eine künftige Lockerung des
   einen Akzeptanzkriteriums (z. B. eine differenziertere Behandlung von
   `LH-FA-SCH-003`, s. `slice-033` <!-- d-check:status-provenance --> §1 Out-of-Scope) müsste sonst aus einer
   gemeinsamen Funktion herausgetrennt werden, statt eine bereits isolierte
   Funktion zu ändern.
2. **Der Boundary-Teil ist nicht redundant** (siehe oben) — er ist der
   einzige Testbeleg für `LH-FA-SCH-003`s Boundary-Klausel im gesamten
   Testbestand.

**Das Duplikations-Risiko ist real und wird nicht verschwiegen:** Der
Happy-Path-Assertion-Block (`awaitHeartbeatErrorClass`/`queryChangesView`-
Abwesenheitsprüfung) ist strukturell identisch zu
`TestE2ESchemaChangeIncompatibleTypeChange`s Fall 2 — dieselben zwei
Aufrufe, unterschiedliche auslösende DDL. Ändert sich der
Fehler-Propagationsmechanismus künftig (neue Fehlerklasse, anderer
Abbruchpfad), müssen beide Testfunktionen unabhängig nachgezogen werden.
Ein gemeinsamer Test-Helper für genau diesen wiederkehrenden Block (etwa
eine Funktion, die DSN, auslösende Transaktion und die nicht erfasste
ID entgegennimmt und Fehlerklasse + Abwesenheit prüft) ist eine sinnvolle
Implementierungs-Option für den umsetzenden Slice — diese ADR schreibt sie
nicht vor (Implementierungsdetail, kein Architektur-Constraint), benennt
aber das Duplikations-Risiko, das ein solcher Helper auflösen würde.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0058`s Happy-Path-Prämisse bleibt bestehen, `TestE2ESchemaChangeDropColumn` wird an das dokumentierte (aber real nie eintretende) Verhalten angepasst | keine neue ADR | die Testfunktion wäre garantiert rot — PostgreSQL/`classifyRelationColumns` liefern das beschriebene Verhalten nie; ignoriert den bereits akzeptierten Präzedenzfall `ADR-0059` Teilfrage 4, der denselben Fakt korrekt beschreibt |
| B — Implementer passt den Test eigenmächtig an die reale Erwartung an, ohne Folge-ADR | schnellste Lösung | widerspricht Modul 8 §Rollen-Regeln (Implementer darf höchstens Folge-ADR vorschlagen, niemals stillschweigend einer Accepted-ADR widersprechen) und `AGENTS.md` §3.5 (Immutabilität); genau das hat der Implementer-Lauf korrekt vermieden |
| C — vollständige Neufassung von `ADR-0058` (ganze ADR superseded) | ein einziges Nachfolgedokument | unverhältnismäßig — vier der fünf Teilentscheidungen (2–5) bleiben unverändert korrekt; dupliziert stabilen Text, entgegen dem bereits im Repo etablierten Muster enger Klausel-Korrekturen (`ADR-0048` Supersedes `ADR-0047`, `ADR-0062` Supersedes `ADR-0045`) |
| **D — enge Klausel-Korrektur per Folge-ADR: nur Entscheidung 1s Happy-Path-Prämisse wird korrigiert, Boundary und die Entscheidungen 2–5 bleiben unverändert (gewählt)** | löst den real gefundenen Widerspruch minimal-invasiv; bestätigt und verankert `ADR-0059` Teilfrage 4 als bereits vorweggenommenen Präzedenzfall statt ihn zu wiederholen; folgt dem etablierten Repo-Muster | zwei ADRs (`0058`+`0063`) müssen zusammengelesen werden, um `LH-FA-SCH-003`s volle Testform-Entscheidung zu verstehen |

## Konsequenzen

- Positiv: `LH-FA-SCH-003` bekommt eine testbare, mit dem real bestehenden
  Code konsistente Happy-Path-Erwartung — keine Testfunktion, die
  strukturell nie grün werden kann.
- Positiv: `ADR-0059` Teilfrage 4 wird als bereits vorweggenommener
  Präzedenzfall bestätigt, nicht revidiert — keine neue Architektur-
  Entscheidung, nur die bislang übersehene Konsequenz für `ADR-0058`s
  Testform.
- Negativ: `TestE2ESchemaChangeDropColumn` beendet den Erfassungspfad des
  Feed-Containers nun genauso dauerhaft wie `TestE2ESchemaChangeIncompatibleTypeChange`
  (`restart: "no"`, kein Neustart-Vertrag). `ADR-0058`s ursprüngliche
  Platzierungs-Annahme — „platziert nach `TestE2ESchemaChangeAddColumn`
  und vor `TestE2EHeartbeatHealthy`/`TestE2ESchemaChangeIncompatibleTypeChange`
  (Container-Ende-Grenze)" — trifft damit nicht mehr zu: Die Funktion
  gehört jetzt **hinter** die Container-Ende-Grenze, nicht davor. Zwei
  Testfunktionen, die den Container je unabhängig und dauerhaft beenden,
  können zudem nicht in derselben `go test`-Aufruf-Gruppe mit **einem**
  gemeinsamen laufenden Container ausgeführt werden, ohne dass die zweite
  Funktion nur noch den bereits vom ersten Fehler gesetzten
  `error_class`-Stand liest (falsch-positiver Beleg, kein echter Test der
  zweiten Funktion).
- Folgepflicht: `test/integration/integration_test.go`s
  `TestE2ESchemaChangeDropColumn` wird auf die oben angegebene Testform
  umgestellt (Implementer-Zug, `slice-062` <!-- d-check:status-provenance -->).
- Folgepflicht: `tools/harness/run-integration-tests.sh` braucht eine
  strukturelle Anpassung — `TestE2ESchemaChangeDropColumn` verschiebt sich
  aus dem vorderen (Container-lebt-noch)-`-run`-Muster in einen eigenen
  Testlauf **nach** der Container-Ende-Grenze, mit einem expliziten
  Container-Neustart (analog zum bereits bestehenden Muster für einen
  simulierten Neustart im Black-Box-CLI-Rundlauf dieses Skripts) und
  Health-Poll **zwischen** `TestE2ESchemaChangeDropColumn` und
  `TestE2ESchemaChangeIncompatibleTypeChange`, weil beide Funktionen
  unabhängig voneinander den Container dauerhaft beenden. Die konkrete
  Reihenfolge und der genaue Neustart-Mechanismus sind Implementierungsdetail
  des umsetzenden Slices (`slice-062` <!-- d-check:status-provenance -->), nicht
  dieser ADR — analog zur Delegation in `ADR-0059` §Konsequenzen.
- Folgepflicht: `slice-062` <!-- d-check:status-provenance -->s Kopf-Feld
  „Bezug" und Plan-Zeilen zu `ADR-0058` werden auf diese ADR nachgezogen
  (Implementer-Zug nach diesem Verdikt, nicht Teil dieses Architect-Zugs).
- Folgepflicht: `harness/README.md` §Sensors/§Werkzeuge, `make
  test-integration`-Zeile, trägt bereits einen Satz zu
  `TestE2ESchemaChangeDropColumn` (`LH-FA-SCH-003`) — der Implementer-Zug
  zieht den Wortlaut auf die korrigierte Testform nach, sobald die
  Umsetzung steht.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` (`TestE2ESchemaChangeDropColumn`) | Happy Path: `cdc.heartbeat.error_class` == `schema` nach `DROP COLUMN`, die danach eingefügte Zeile (id=4) fehlt in `cdc.changes`. Boundary: die vor der Entfernung erfasste Zeile (id=3) bleibt unverändert lesbar, inklusive des historischen Werts der entfernten Spalte | `make test-integration` (kein Gate, [`ADR-0030`](0030-testpyramide.md)) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die
Entscheidung unbefristet weiter, auch wenn ihre Voraussetzung weg ist
(Baseline-Regelwerk `modul-04-adrs.md` §Kernidee (Modul 4)).

Permanent — die Korrektur folgt aus bereits bestehendem, unveränderlichem
Code-Verhalten (`classifyRelationColumns`, seit `slice-033` <!-- d-check:status-provenance -->
real verifiziert) und einem bereits akzeptierten Präzedenzfall
(`ADR-0059` Teilfrage 4), keine externe Abhängigkeit. Ändert sich
`classifyRelationColumns`s Klassifikationslogik künftig selbst (z. B. eine
differenziertere Behandlung von `LH-FA-SCH-003` gegenüber `LH-FA-SCH-004`,
bereits in `slice-033` <!-- d-check:status-provenance --> §1 als eigener, hier
nicht geleisteter Vorgang benannt), wird auch diese Testform per
Folge-ADR erneut geprüft.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: `slice-062` <!-- d-check:status-provenance -->-Implementer-Lauf fand real einen Widerspruch zwischen `ADR-0058` Entscheidung 1 und dem bereits bestehenden, von `ADR-0059` Teilfrage 4 akzeptierten `relationOther`-Verhalten (`classifyRelationColumns`) und stoppte korrekt statt selbst aufzulösen (Modul 8 §Konflikt-Pfad, Verdikt 2); unabhängiger Architect-Zug bestätigt und korrigiert die betroffene Klausel | `slice-062` (in `in-progress/`) <!-- d-check:status-provenance -->, [`ADR-0059`](0059-spaltenauswahl-mechanismus.md) (Präzedenzfall), `docs/plan/planning/done/slice-033-typauswertung-fehlerklasse-schema.md` <!-- d-check:status-provenance --> (Bestandsbeleg) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0063` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
