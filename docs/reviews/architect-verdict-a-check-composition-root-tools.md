# Architect-Verdikt: `tools/**` in `composition_root` — Gate-Scope-Erweiterung ohne Träger

**Rolle:** Architect (Modul 8)

**Anlass:** Finding F-1 (HIGH) aus
Review zu `slice-071`. Der Reviewer hat die Änderung
nicht als Code-Defekt, sondern als **Entscheidung** gereicht, die in
`slice-071` fällig ist — drei rückfragefrei entscheidbare Antworten sind in
seinem Abschnitt *Zum `.a-check.yml`-Verdikt* ausformuliert. Damit läuft der
Rollenwechsel Reviewer → Architect → Planner/Implementer (Modul 8
§Konflikt-Pfad; Verdikt 3: „Erweiterung zulässig, aber falsch zugeschnitten").

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der
Implementer-Lauf von `slice-071`, der `.a-check.yml` erweiterte, und als der
Reviewer-Lauf)

**Datum:** 2026-09-14

**Bezug:** [`AGENTS.md`](../../AGENTS.md) §3.6 (Schwellen-Senkung =
ADR) · [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
(§Entscheidung, §Re-Evaluierungs-Trigger — in der Ausnahmeklausel
superseded) · [`ADR-0026`](../plan/adr/0026-composition-root.md) (Composition
Root) · [`ADR-0002`](../plan/adr/0002-abhaengigkeitsrichtung.md) ·
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) ·
[`spec/architecture.md`](../../spec/architecture.md) (Zeilen
[`ARC-005`](../../spec/architecture.md)/[`ARC-007`](../../spec/architecture.md)) ·
[`LH-FA-SST-008`](../../spec/lastenheft.md) ·
`.a-check.yml` (Stand `b835dde`) · `tools/harness/grpcclient/main.go` ·
`test/integration/integration_test.go` ·
`docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md` (§3) ·
`harness/sensors/a-check.md`

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  (Accepted, `Supersedes ADR-0041` **nur** die Änderungs-Ausnahmeklausel),
  samt Index-Zeile und Nachfolge-Vermerk in
  [`docs/plan/adr/README.md`](../plan/adr/README.md)
- `.a-check.yml` — Rücknahme von `tools/**` aus `composition_root`, Aufnahme
  der Gruppe `tooling: ["tools/harness/**"]` + Kante
  `{from: tooling, to: adapters}`, korrigierter Kommentar
- `harness/sensors/a-check.md` — Bindung um `ADR-0068` ergänzt, Grenze 1
  geschärft (Verfeinerung vs. Erweiterung)
- diese Verdikt-Datei

**Nicht erzeugt:** ein **neuer Folge-Slice**. Der Fix ist ein Eingriff in
`slice-071` — der Slice ist noch **nicht** in `done/`, F-2/F-3/F-4 liegen
ohnehin als Rückkante an, und der `.a-check.yml`-Eintrag ist Teil seines
eigenen Diffs. Ein Folge-Slice würde einen offenen Liefer-Punkt doppeln.

---

## Verdikt

**`tools/**` gehört nicht in `composition_root` — in keiner Weite. Der Client
ist kein Composition Root: er verdrahtet nichts und importiert genau ein
Paket.** Die Rolle des Schlüssels ist „verdrahtet konkret"
([`ADR-0026`](../plan/adr/0026-composition-root.md),
[`ARC-007`](../../spec/architecture.md)); `composition_root` listet seine drei
Träger (`internal/bootstrap/**`, `cmd/**`, `test/integration/**`) — alle drei
verdrahten real. Ein Client, der den Protokoll-Stub über das Netz konsumiert,
tut das nicht.

**Die zutreffende Form ist eine begrenzte Kante, keine unbeschränkte
Ausnahme.** `.a-check.yml` führt künftig eine Gruppe
`tooling: ["tools/harness/**"]` **und** genau eine Kante
`{from: tooling, to: adapters}`. Der Client darf den erzeugten Stub
importieren; ein Import aus `app`, `ports` oder `domain` in
`tools/harness/**` bleibt `wrong-direction`. Damit ist der Review-Befund
„unbeschränktes Importrecht für den ganzen Baum" ausgeräumt — real rot-fähig
belegt (§Befundlage).

**Das braucht eine ADR.** `ADR-0041` nennt `composition_root` in der
**Definition** des Standes und lässt ihn in **beiden** Wiederholungen der
ADR-freien Änderungsklasse aus; die Aufnahme ist eine **Erweiterung der
Architekturregel**, keine Abbildung von §2, und fällt unter
[`AGENTS.md`](../../AGENTS.md) §3.6. Träger ist die Folge-ADR
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md).
Sie schärft zusätzlich die Klausel selbst: ADR-frei sind nur Änderungen, die
den §2-Komponentensatz und seine Constraints **unverändert abbilden**
(Verfeinerung); das Aufnehmen eines Nicht-§2-Bereichs **mit
Import-Berechtigung** — über `layers`/`edges` **oder** `composition_root` —
ist ADR-pflichtig. Das ist die eine Hälfte, die die Klausel heute offen lässt;
die andere Hälfte (`composition_root` fehlt in der Ausnahmeliste) fällt damit
mit.

**Träger:** Fix in `slice-071` (offen), **kein Folge-Slice**. Umsetzung des
`.a-check.yml`-Teils in diesem Zug; die Plan-Korrektur ist Rückkante zum
Implementer-/Planner-Lauf (§Frage 4).

---

## Befundlage (eigener Code-Gang, real ausgeführt)

Alle Läufe ungefiltert in eigene Log-Dateien, Exit-Code danach in einem
**eigenen, ungeketteten** Schritt ermittelt ([`AGENTS.md`](../../AGENTS.md) §3.9).

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make a-check` (Stand `HEAD`, `tools/**` in `composition_root`) | **0** | `gesamt: 0 Befund(e)` — der Baum ist mit dem implementierten Stand grün |
| `make a-check` + Probe-Datei `tools/harness/grpcclient/probe_scope.go` (Import `internal/application/usecase/capture`), Stand `HEAD` | **0** | grün **ohne jeden Hinweis** — die Aufnahme gibt dem ganzen Baum unbeschränktes Importrecht. F-1 bestätigt, unabhängig reproduziert |
| `make a-check` + Probe derselben Datei, mit **diesem** Zug (`tooling`-Kante) | **2** | `tools/harness/grpcclient/probe_scope.go:3: wrong-direction: tools -> app (…/usecase/capture)` — das Recht ist auf `adapters` begrenzt |
| `make a-check` (Stand dieses Zugs) | **0** | `gesamt: 0 Befund(e)`; der Stub-Import bleibt grün, der Abdeckungs-Hinweis „Dateien ohne Schicht" ist weg (die drei Client-`.go` liegen jetzt in `tooling`) |

**Die Rücknahme und Blatt-Identität:** die Probe-Datei ist gelöscht,
`git status --porcelain` wird vor den Commits geprüft. Die drei
`.go`-Dateien der Clients sind unberührt (`git diff` nennt nur `.a-check.yml`,
`harness/sensors/a-check.md`, die ADR-Datei und den Index).

---

## Frage 1 — Weite und Rolle

Der Client importiert **ein** Paket (`internal/adapters/driving/grpc/streamv1`,
`main.go:25`) und verdrahtet nichts. `test/integration/*.go` importiert
**sieben** interne Pakete (`postgresstorage`, `inbound`, `outbound`, vier
Use-Cases, `domain/model`) und verdrahtet real — der Vergleich im
Nachzug-Kommentar trägt sachlich nicht.

| Weite | Trägt? | Warum |
|---|---|---|
| `tools/**` (implementiert) | **nein** | nimmt `tools/schema/**` und die Bench-Skripte mit, die der Kommentar nicht meint; gibt dem ganzen Baum unbeschränktes Importrecht (Probe grün); Rolle falsch |
| `tools/harness/**` | **nein** | verengt den **Baum**, nicht das **Recht** — die Use-Case-Probe bliebe grün; die Rolle „Composition Root" trägt weiter nicht |
| nur der konkrete Client | **nein** | Per-Client-Ausnahme proliferiert; behebt keinen der drei Review-Punkte; Recht bleibt unbeschränkt |
| **Gruppe `tooling` + Kante `to: adapters` (gewählt)** | **ja** | begrenzt das Recht auf `adapters`; Probe wird rot; Abhängigkeit steht als **Kante** im Modell, nicht als Ausnahme davon |

Die vollständige Alternativen-Tabelle (A–F, inkl. „nichts tun" und „Client nach
`test/` verschieben") steht in
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
§Verglichene Alternativen; die Weite-Entscheidung ist dort Festlegung 2.

---

## Frage 2 — ADR ja/nein

**Ja.** `ADR-0041` §Entscheidung: „Änderungen an Schichten-Globs und Edges sind
Änderungen dieser Datei, keine neuen ADRs, solange die §2-Constraints
unverändert abgebildet bleiben." Der §Re-Evaluierungs-Trigger wiederholt genau
„Schichten-Globs und Edges". `composition_root` steht **in der Definition des
Standes** und in **keiner** der beiden Ausnahmen. Die Aufnahme ist zudem keine
Abbildung von §2 — sie erteilt einer Pfad-Menge, die §2 nicht als Komponente
führt, eine Import-Berechtigung, und §2s `ARC-007`-Zeile („kennt als **einzige**
Komponente konkrete Adapter") ist die, deren Maschinenform `composition_root`
ist. Damit ist es eine Erweiterung der Architekturregel → ADR
([`AGENTS.md`](../../AGENTS.md) §3.6).

Die Klausel ist **an beiden Enden** unzutreffend gefasst: sie nimmt
`composition_root` nicht in die ADR-freie Klasse auf (zu eng für den einen
Leser) **und** sie nimmt „Schichten-Globs und Edges" pauschal aus, ohne
Verfeinerung von Erweiterung zu trennen (zu weit für den anderen). Beides
schließt [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
Festlegung 3. Der Reviewer-Vorschlag „Antwort 1" (die Klausel nachziehen, damit
`composition_root` in beide Sätze gehört) ist damit **nicht** gewählt: er hätte
`composition_root` in die ADR-freie Klasse gezogen; diese ADR zieht ihn
ausdrücklich heraus.

---

## Frage 3 — Ausgelieferter Zustand

**Die Zeile bleibt nicht als `composition_root`-Eintrag; der Client bleibt an
seinem Ort.** Konkret:

- `tools/**` verlässt `composition_root` (Rücknahme, kein Ersatz dort).
- Der Client behält seinen Stub-Import — es wird **nicht** verlangt, die
  Protobuf-Form handzurollen, und **nicht**, den Client umzuziehen.
- Die Berechtigung steht als Kante `tooling → adapters`.

Der Client hätte auch **ohne** jede Berechtigung einen Weg gebraucht; die drei
vom Reviewer benannten Auswege sind in
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
§Verglichene Alternativen geprüft: Import des Stubs vermeiden (D/Handrollen —
verworfen, brüchig), Client unter `test/integration/**` verschieben (D —
verworfen als Ablage-Dodge), Glob anders schneiden (B/C — verworfen, weil sie
das Recht nicht verengen). Gewählt ist die begrenzte Kante (E).

`make a-check` ist mit dem neuen Stand **grün** (Exit 0, §Befundlage);
`make docs-check` deckt die neuen Markdown-Verweise (Gate-Beleg unten).

---

## Frage 4 — Träger

- **Fix in `slice-071`, kein Folge-Slice.** `.a-check.yml` ist Teil des
  Slice-Diffs (`b835dde`); F-1 ist merge-blockierend; der Slice ist in
  `in-progress/`. Der Fix bleibt innerhalb des bestehenden Liefer-Punkts
  („E2E-Beleg") — kein Liefer-Punkt-Zuwachs, kein Re-Cut.
- **Plan-Korrektur als Rückkante** (Implementer-/Planner-Zug): §3
  (Tabellenzeile `.a-check.yml`) und der Absatz *Implementer-Entscheidungen*
  verweisen auf die zurückgenommene `composition_root`-Erweiterung und die
  sachlich falsche `test/integration`-Gleichsetzung. Sie werden auf
  [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  und die Kante `tooling → adapters` nachgezogen — zusammen mit F-2/F-3/F-4,
  die bereits an dieselbe Rolle gehen. Die Arbeit am Plan ist nicht dieser
  Zug (Modul 8).
- **Kein neuer Registereintrag in diesem Zug.** Die Finding-Klasse zu F-1
  („Gate-Scope-Erweiterung ohne die Rolle, die der Schlüssel deklariert")
  gehört in die **Slice-Closure** §7 von `slice-071` und von dort in den
  Zähler; der verwandte, **verkörperte** Eintrag ist
  `BEO-PGC/a-check-null-abdeckung` (dessen Gegenrichtung — der Hinweis
  verstummt, statt zu greifen — dieser Befund trifft). Ein Eintrag aus dem
  Architect-Zug heraus wäre eine zweite Schreibstelle für denselben Zähler.

---

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `.a-check.yml` — `tools/**` aus `composition_root` entfernt; Gruppe
  `tooling: ["tools/harness/**"]` und Kante `{from: tooling, to: adapters}`
  ergänzt; Kommentar auf die Rolle „testseitiger Protokoll-Konsument"
  korrigiert (die falsche `test/integration`-Gleichsetzung entfällt)
- `docs/plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md` (neu,
  Accepted) und `docs/plan/adr/README.md` (Zeile `ADR-0068`; Zeile `ADR-0041`
  mit Nachfolge-Vermerk)
- `harness/sensors/a-check.md` — Bindung um `ADR-0068` ergänzt, Grenze 1 um
  die Verfeinerung/Erweiterung-Unterscheidung geschärft
- diese Verdikt-Datei

**Nicht geändert:**

- `internal/**` — kein Code-Eingriff.
- `tools/harness/**` — kein Umzug, kein Import-Eingriff; die Clients bleiben
  wie sie sind.
- `test/integration/**` und `spec/architecture.md` — die Sicht trägt keine
  `tooling`-Zeile; `test/integration/**` bleibt `composition_root`-Mitglied
  (bestätigt, nicht geändert).
- `docs/plan/planning/in-progress/slice-071-…md` — Plan-Nachzug ist
  Planner-/Implementer-Arbeit (§Frage 4).

## Offen für den Implementer-/Planner-Zug (`slice-071`)

1. §3-Tabellenzeile `.a-check.yml` und der Absatz *Implementer-Entscheidungen*
   auf die Kante `tooling → adapters` und
   [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
   nachziehen; die `test/integration`-Gleichsetzung streichen.
2. F-2 (Erfolgsmeldung breiter als die Assertion), F-3 (`SPEC-019` → `SPEC-020`
   im Kopf), F-4 (Aufschub-Adresse `slice-072` → `slice-077`) — unverändert
   aus dem Review.
3. Ein weiterer Review-Lauf des geänderten `.a-check.yml` ist zulässig; der
   Skill verlangt ihn nicht von selbst, aber F-1 ist mit diesem Zug nicht
   „weg", sondern **entschieden** — der neue Stand ist zu prüfen.

Nach `Accepted` der [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
wird die alte Klausel nicht nachgebessert: die Nachfolge trägt der Index
([`AGENTS.md`](../../AGENTS.md) §3.5).
