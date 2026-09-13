# Architect-Review slice-014 — Verdikt zu F-1 (globaler `slog`-Singleton vs. `ADR-0024`)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-10.
**Eingang:** [`docs/reviews/review-slice-014.md`](../../reviews/review-slice-014.md)
F-1 (HIGH, Rollen-Widerspruch) · Slice-Plan
[`docs/plan/planning/done/slice-014-strukturiertes-logging.md`](../planning/done/slice-014-strukturiertes-logging.md)
§1/§3 · [`ADR-0024`](0024-observability-ausserhalb-der-domain.md) (Accepted,
permanent) · [`ADR-0026`](0026-composition-root.md) (Accepted, permanent) ·
`spec/architecture.md` `ARC-011` (Telemetrie-Backend) · `spec/lastenheft.md`
`LH-QA-OPS-003`/`LH-QA-OPS-004` · Präzedenzfall
[`architect-review-slice-011.md`](architect-review-slice-011.md) (Heartbeat-
Persistenz als Erweiterung derselben `ADR-0024`-Fläche) · Baseline-Regelwerk
`modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz.
**Ausgang:** Übergabe-Artefakt an Planner/Implementer — **kein Folge-ADR,
kein Carveout**; Verdikt 1 aus dem Konflikt-Pfad (`ADR-0024` gilt
uneingeschränkt für strukturiertes Logging, die Implementierung war falsch).
Fix-Zug liegt beim Implementer.
**Harte Regel eingehalten:** keine Accepted-ADR in-place geändert; keine
Folge-ADR-Datei angelegt (die Prüfung unten verneint den Bedarf); kein
ADR-Index-Eintrag (dieses Dokument ist kein ADR); der Slice-Plan wird von
diesem Lauf **nicht** editiert — Fix-Umsetzung ist Implementer-Sache im
nächsten Zug (Modul 8 §Konflikt-Pfad, Sequenz).

---

## Kernfrage — deckt `ADR-0024` reines strukturiertes Logging ab, oder nur Metriken/Events?

**Kurzform: ja, eindeutig — und zwar mehrfach unabhängig belegt, nicht nur
durch eine Lesart des Entscheidungssatzes.**

Vier voneinander unabhängige Textstellen tragen dieselbe Aussage:

1. **`ADR-0024` `Bezug:`-Feld** nennt explizit **beide** Anforderungen:
   [`LH-QA-OPS-003`](../../../spec/lastenheft.md) (Metriken) **und**
   [`LH-QA-OPS-004`](../../../spec/lastenheft.md) (strukturiertes Logging).
   Ein ADR-Kopf, der nur Metriken entscheiden wollte, hätte `LH-QA-OPS-004`
   nicht im Bezug geführt.
2. **`ADR-0024` Kontext** formuliert wörtlich: „Das Lastenheft fordert
   maschinenlesbare Metriken **und strukturiertes Logging** (…). Die Domain
   soll davon frei bleiben (…), während Betriebsinformationen trotzdem aus
   allen Schichten ankommen müssen." — Logging ist Teil der Prämisse, nicht
   nachträglich hineingelesen.
3. **`ADR-0024` Entscheidung** sagt: „**Logging-/Metrics-Frameworks** bleiben
   Infrastruktur. MetricsPort und EventSinkPort werden durch Driven Adapters
   implementiert." Der erste Satz benennt Logging-Frameworks ausdrücklich als
   von der Entscheidung erfasst; der zweite Satz benennt den Mechanismus
   (Port + Driven Adapter) für die aus dem ersten Satz gemeinsam behandelte
   Klasse. Dass nur zwei Portnamen fallen (`MetricsPort`, `EventSinkPort`) und
   kein explizites `LogPort`, ist eine Frage der **Umsetzung** dieses
   Mechanismus, keine Einschränkung seines **Geltungsbereichs** — dieselbe
   Lesart, die schon `architect-review-slice-011.md` für die
   Heartbeat-Persistenz trug („dieselbe Klasse trägt eine
   Heartbeat-Persistenz", kein neuer Adapter-Typ nötig).
4. **`ADR-0024` `Schärft:`-Feld** nennt `ARC-011` ausdrücklich. `ARC-011`
   selbst (`spec/architecture.md` §3) sagt: „Telemetrie-Backend (…) |
   Metriken **und strukturierte Logs** | über den Outbound Port
   substituierbar; Frameworks bleiben Infrastruktur." Das ist die
   **Sicht-Schicht** (Rang 3 der Source Precedence, `harness/README.md`),
   die über der ADR (Rang 4) steht — eine ADR schärft die Sicht, sie
   schränkt sie nicht ein. Selbst eine enger gelesene `ADR-0024` könnte die
   in `ARC-011` bereits gesetzte Pflicht „strukturierte Logs (…) über den
   Outbound Port substituierbar" nicht aufheben; dafür bräuchte es eine
   Änderung von `spec/architecture.md` selbst, die hier weder beantragt noch
   angezeigt ist.

Diese vier Stellen ziehen unabhängig voneinander dieselbe Grenze. Es gibt
keine belastbare Lesart, in der `ADR-0024` (und die von ihr geschärfte
`ARC-011`) reines strukturiertes Logging ausnimmt.

## Prüfung — verletzt die vorliegende Implementierung diese Entscheidung?

Ja, direkt und vollständig, unabhängig am Code nachvollzogen (nicht aus dem
Review-Report übernommen):

- `internal/bootstrap/wiring.go:213` setzt `slog.SetDefault(newLogger(cfg.LogLevel))`
  — ein **globaler** Paket-Zustand der stdlib, gesetzt einmal in der
  Composition Root.
- `internal/adapters/driven/postgresack/ack.go:37,50`,
  `internal/adapters/driven/postgresstorage/{store,heartbeat,consumerstate,tableactivation}.go`
  und `internal/adapters/driving/replication/receive/receive.go` rufen
  danach die **paketweiten** `slog.Info`/`slog.Error`/`slog.InfoContext`/
  `slog.ErrorContext`-Funktionen **direkt** auf — jeder Adapter importiert
  `log/slog` selbst und schreibt gegen den globalen Default, statt einen
  injizierten Port zu bedienen.
- `grep -rn "EventSinkPort\|MetricsPort" internal/` liefert **keinen**
  Treffer — weder für Logging noch für Metriken existiert im Repo aktuell
  ein Outbound Port; die einzigen bestehenden Outbound Ports sind
  `changestore`, `clock`, `consumerstate`, `heartbeat`, `replicationack`.

Das ist exakt Alternative B aus `ADR-0024` („globale Telemetrie-
Singletons"), und die dort für B notierte Contra-Begründung trifft
wortwörtlich zu: „versteckte Abhängigkeit; Verdrahtung und Testdoubles
intransparent" — jeder Adapter hängt am globalen `slog`-Zustand statt an
einer injizierten Port-Grenze; ein Testdouble für Logging-Verhalten lässt
sich an keiner der Aufrufstellen ohne globalen State-Umbau einsetzen.

## Einordnung nach Modul 8 §Konflikt-Pfad

Von den drei legitimen Verdikten trägt **Verdikt 1**: *die Entscheidung
gilt und die Implementierung/Selbstverteidigung hat falsch adressiert.*

- **Nicht Verdikt 2** (Folge-ADR mit `supersedes`, Grenze neu ziehen): Es
  gibt keine unklare oder widersprüchliche Textstelle, die eine Klärung
  bräuchte — vier unabhängige Belege (oben) ziehen dieselbe Grenze
  übereinstimmend. Eine Lockerung würde außerdem nicht nur `ADR-0024`,
  sondern auch `ARC-011` berühren, was eine ADR nicht leisten kann
  (Sicht steht über ADR in der Source Precedence). Eine Lockerung wäre
  hier keine Klarstellung, sondern eine Senkung einer bestehenden,
  kohärenten Anforderung ohne erkennbaren neuen Grund — `AGENTS.md` §3.6
  verlangt dafür ohnehin ein ADR, keinen Implementer-Alleingang.
- **Nicht Verdikt 3** (legitime, aber undokumentierte Lockerung, per
  Folge-ADR nachgezogen): Die Commit-Begründung (`e499e42`) verteidigt die
  Wahl nicht als bewusste Abweichung von `ADR-0024`, sondern adressiert die
  **falsche** ADR (`ADR-0026`, Composition Root — dort geht es um die
  fachliche Verdrahtung von Driving Adapters/Inbound Ports/Application
  Services/Outbound Ports/Driven Adaptern als Objekte im
  Abhängigkeitsgraph, nicht um Cross-Cutting-Infrastruktur wie Logging).
  `ADR-0026` ist durch `slog.SetDefault` tatsächlich **nicht** verletzt —
  der Reviewer hat das in seinem Report bereits korrekt festgestellt, und
  dieses Verdikt bestätigt das unabhängig. Es gibt also keine bewusste,
  begründete Entscheidung *gegen* `ADR-0024`, die man nachträglich als
  Lockerung anerkennen könnte — nur eine Prüfung gegen die falsche ADR.
- **Verdikt 1 trägt:** `ADR-0024` gilt unverändert (Accepted, permanent,
  Re-Evaluierungs-Trigger „permanent — die Domain-Freiheit von Telemetrie
  trägt die Abhängigkeitsregel der Architektur-Sicht"). Die
  Implementierung weicht davon ab; die Selbstprüfung des Implementers hat
  die falsche ADR herangezogen (`ADR-0026` statt `ADR-0024`) und dadurch
  den eigentlichen Konflikt nicht gesehen.

## Disposition

| Frage | Ergebnis |
|---|---|
| Deckt `ADR-0024` strukturiertes Logging ab? | Ja — `Bezug`, Kontext, Entscheidungssatz, `Schärft:ARC-011` tragen übereinstimmend |
| Verletzt der globale `slog`-Singleton `ADR-0024`? | Ja — Alternative B, wortwörtlich mit der dort verworfenen Begründung |
| Verletzt der globale `slog`-Singleton `ADR-0026`? | Nein — bestätigt, `ADR-0026` betrifft fachliche Verdrahtung, nicht Cross-Cutting-Infrastruktur |
| Konflikt-Pfad-Verdikt | **1** — ADR gilt, Implementierung/Selbstprüfung war falsch |
| Folge-ADR nötig? | Nein |
| Carveout nötig? | Nein — kein rotes Gate, `a-check` modelliert diese Pflicht nicht (Review-Negativbefund „verifizierbar: nein") |

**F-1 bleibt HIGH und merge-blockierend.** Der `in-progress → done`-Übergang
von slice-014 bleibt gesperrt, bis der Fix vorliegt und erneut geprüft ist.

## Fix-Zug (Implementer-Sache, nicht Teil dieses Verdikts)

Strukturierte Log-Aufrufe in den betroffenen Driven-/Driving-Adaptern
(`postgresack`, `postgresstorage/*`, `replication/receive`) und im
Bootstrap müssen über einen Outbound Port + Driven Adapter laufen, nicht
über den globalen `slog`-Default. `ADR-0024` autorisiert den Mechanismus
bereits (Port + Driven Adapter); **welche konkrete Portform** — Erweiterung
von `EventSinkPort`, ein eigenständiger neuer Port für Log-Zeilen, oder ein
gemeinsamer Telemetrie-Port für Metriken und Logs — ist eine
**Umsetzungsfrage unter bereits Entschiedenem**, keine neue
Architekturentscheidung (Präzedenz: `architect-review-slice-011.md`,
Heartbeat-Persistenz als Erweiterung derselben `ADR-0024`-Fläche ohne neuen
Adapter-Typ). Sollte der Refactor zusammen mit einer neuen Port-Definition
die verbliebene Slice-Größe sprengen (Port-Definition in
`application/port/outbound`, Driven-Adapter-Implementierung,
Bootstrap-Neuverdrahtung, Umstellung aller bisherigen Aufrufstellen —
mehrere Schichten gleichzeitig), ist das nach den Größenkriterien in Modul 5
ein regulärer Kandidat für `in-progress → next`; diese Einschätzung liegt
beim Implementer/Planner im nächsten Zug, nicht bei diesem Verdikt.

---

## Beleg-Anker (Kurzfassung)

| Aussage | Tragende ADR/Spec/Code-Stelle |
|---|---|
| `ADR-0024` nennt Logging explizit im Bezug | `docs/plan/adr/0024-observability-ausserhalb-der-domain.md` `Bezug:`-Feld |
| `ADR-0024`-Kontext behandelt Metriken und Logging gemeinsam | ebd., Abschnitt „Kontext" |
| `ADR-0024`-Entscheidung nennt „Logging-/Metrics-Frameworks" gemeinsam | ebd., Abschnitt „Entscheidung" |
| `ARC-011` verlangt strukturierte Logs „über den Outbound Port substituierbar" | `spec/architecture.md` §3, Zeile `ARC-011` |
| Globaler Singleton + paketweite Aufrufe, kein Port | `internal/bootstrap/wiring.go:213`, `internal/adapters/driven/postgresack/ack.go:37,50`, `internal/adapters/driven/postgresstorage/{store,heartbeat,consumerstate,tableactivation}.go`, `internal/adapters/driving/replication/receive/receive.go` |
| Kein `EventSinkPort`/`MetricsPort` im Repo | `grep -rn "EventSinkPort\|MetricsPort" internal/` → kein Treffer |
| `ADR-0026` betrifft fachliche Verdrahtung, nicht Cross-Cutting-Infra | `docs/plan/adr/0026-composition-root.md` Kontext/Entscheidung/Alternative B |
| Präzedenz „Port-Erweiterung statt neuer Adapter-Typ" | `docs/plan/adr/architect-review-slice-011.md` §Prüfung: Schreib-Seite |
