# Welle welle-20: Coverage 80 % über der netzlos prüfbaren Fläche

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<NN>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** pt9912. **Datum:** 2026-09-15.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**„Die Coverage" dieses Repos erreicht 80 %** — die Gate-getragene Unit-Zahl
über der **netzlos prüfbaren Fläche** ([`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md);
und „die Coverage" ohne Subjekt-Zusatz meint genau diese Zahl, nie die
DB-Adapter-Coverage). **Der Nenner ist eine Zustandsgröße, die sich mit jedem
Zug bewegt:** `slice-079` hat ihn auf **1679** geschnitten, seither ist er durch
die Nähte und Endpunkte der folgenden Slices auf **1903** gewachsen (zuletzt
`slice-084`, +9; eigener Messstand: Lauf `slice-089`). Wer den Fortschritt dieser Welle liest, liest die **Quote**,
nicht eine eingefrorene Statement-Zahl — die Klasse
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` steht genau dafür.
Der Weg zum Ziel hat zwei Hälften: die **Präzisierung des Gegenstands**
(`slice-079`: die drei Pakete, deren Tests ohne externen Dienst überspringen,
verlassen den Nenner — der Ist-Stand sprang real von 49,3 % auf ~69,7 %,
**ohne eine Zeile Test**) und danach **Test-Arbeit** an dem, was ungedeckt
bleibt (Abschnitt *Slices in dieser Welle*: `internal/bootstrap` 327 ungedeckt,
`cmd/pg-change-feed` 49, Rest-Tail). Getragen wird das Maß von
`make coverage-gate` gegen `THRESHOLD`; die Schwelle wandert nach dem
unveränderten bootstrap-aware-Mechanismus ([`ADR-0054`](../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a)) stufenweise bis 80. **Die gelten Stufen heute:** Einstieg **70 %**
(angehoben durch den Subjekt-Transfer aus `slice-081`) → Endstufe **80 %**.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `slice-079` (Scope-Schnitt und Neukalibrierung des Coverage-Gates) liegt in
  `done/` — **erfüllt.** Ohne ihn misst die Welle gegen einen Nenner, den
  [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  ersetzt hat, und ihr eigener Fortschritt wäre nicht belegbar.
- [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  ist `Accepted` — **erfüllt** (der Gegenstand ist entschieden).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices dieser Welle liegen in `done/`.
- `make coverage-gate` ist **real grün bei `THRESHOLD=80`** über der netzlos
  prüfbaren Fläche, mit eigenem Grün-Beleg dieses Laufs. Das ist das *Mehr*:
  keine einzelne Slice-DoD kann eine repo-weite Schwelle belegen.
- Ein **Rot-Beleg unmittelbar über der Endstufe** ist geführt (`THRESHOLD=85`,
  Exit ≠ 0) — sonst prüfte die Stufe nicht real.
- `make gates` grün.
- Closure-Notiz in `welle-20-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| `slice-079` | Coverage-Gate: Scope-Schnitt und Neukalibrierung | [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) |
| `slice-088` | Coverage-Tail „Reine Übersetzung" — Cluster B | [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) |
| `slice-091` | Coverage Cluster C — Zustell- und Betriebs-Rand | [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) |
| `slice-092` | Coverage Cluster D1 — Anwendungs-Kern | [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) |
| `slice-093` | Coverage Cluster D2 — Bootstrap-Rest und Telemetrie | [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) · [`ADR-0085`](../adr/0085-build-kontext-ausnahme-test-only-zweck.md) |
| `slice-094` | Coverage Cluster A — Prozess-Rand | [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) |

**Die Liste ist vollständig, und sie wächst mit dem Schnitt** — jede Zeile entsteht
mit ihrer Slice-Datei, im selben Zug. Sie ist keine zweite Zustandsquelle (§Lifecycle:
der Zustand bleibt das Verzeichnis) und trägt darum **keinen** Status; sie ist der
Überblick, welche Slices diese Welle beansprucht. **Der Nachzug dieser vier Zeilen
war überfällig** — die Liste stand seit Cluster C bei den ersten zwei, während die
Cluster gearbeitet wurden.

**Das Schnittmaß steht — und es hat den ersten Vorschlag dieser Welle
widerlegt** ([`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)).
Diese Sektion nannte `internal/bootstrap` „den Hebel" (591 Statements, 327
ungedeckt, „ggf. zwei Slices"). **Gemessen ist das falsch:** bootstrap hat 336
ungedeckte Statements, davon sitzen **308 in fünf Funktionen, die kein
Test-Slice bewegt** — netzlos prüfbar sind dort nur **26–28**. Und `Run` (198
Statements) ist **netzlos gar nicht prüfbar**: 0 von 198 im Gate-Profil, weil
jeder `postgresstorage.New*` ein `pgxpool.New` **plus `pool.Ping`** ist und
`receive.NewStream`/`nats.Connect` ebenso einen lebenden Dienst verlangen.

**Das Maß ist der Tail** — vier Cluster, je ein Slice, geschnitten **nach der
Messung** (`ADR-0082` §Schnittmaß):

| Cluster | Träger | ungedeckt |
|---|---|---|
| **B — Reine Übersetzung** | `replication/decode`, `replication/mapper`, `postgresstorage/sqlexec`, `postgresstorage/mapper` | 60 |
| **C — Zustell- und Betriebs-Rand** | `driving/http`, `driving/grpc`, `driven/natsnotify`, `driving/grpc/streamv1` | 62 |
| **D — Anwendungs-Kern und Bootstrap-Rest** | Use-Cases, `domain/model`, `telemetry`, `bootstrap`-Rest | 53–55 |
| ↳ **D1 — Anwendungs-Kern** | `application/usecase/*` (13 Pakete, 10 davon mit ungedeckten Statements), `domain/model` | 24 |
| ↳ **D2 — Bootstrap-Rest und Telemetrie** | `bootstrap`-Rest, `adapters/driven/telemetry` | 30 |
| **A — Prozess-Rand (Puffer)** | `cmd/pg-change-feed` (`main`-Dispatch), `bootstrap` (`Run`-Fehlerpfad) | ≈54 |
| | **Summe** | **229–231** |

**B + C + D tragen die 154 fehlenden Statements** (Ziel 1523 von 1903) — aber
nur mit **13 Statements Puffer**; **A ist kein Beiwerk, sondern der Puffer** (mit
A: 63 Statements, 3,3 pp). Die Decke ohne die sechs unbeweglichen Funktionen
liegt bei **81,2 %**, mit den Präfixen bei **84,1 %** — die Endstufe 80 % ist
über dem **unveränderten** Gegenstand erreichbar, und der unten verlangte
Rot-Beleg bei `THRESHOLD=85` ist damit **notwendig** rot.

**Herkunft der Zahlen dieses Abschnitts** (`AGENTS.md` §3.12 gilt für die
Zahlen dieses Plans mit): die **ungedeckten** Statement-Zahlen der
Cluster-Tabelle (**62**, **53–55**, **≈54** und ihre Summe **229–231**) stehen
in [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
§*Was daraus für die Slices folgt*; die Paket- und Funktions-Werte (`336`,
`308`, `26–28`, `198`, `0 von 198`) und der Nenner `1903` stehen in §Kontext
(2)–(5) — dort **gemessen** in zwei Läufen über denselben Quelltext,
dedupliziert über die Block-Position, ADR-Stand 2026-09-16. Die Zahl **60** der
ersten Cluster-Zeile ist die **ungedeckte** (`60`), während die Angabe
„**60 von 60** erreichbar" weiter unten die **zweite** Spalte derselben Tabelle
liest (`netzlos erreichbar`) — zwei verschiedene Größen, beide aus derselben
Quelle. Der Vorgänger-Nenner `1679` steht in §1 und stammt aus `slice-079`.
Die **Summen, Differenzen und Prozente** dieses Abschnitts (`1523`, `154`,
`13`, `63`, `3,3 pp`, `81,2 %`, `84,1 %`) sind **gerechnet**, nicht gemessen —
dieselbe Zahl kann in [`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
§Kontext (5) stehen, gerechnet ist sie dort ebenso. Die zwei Werte des
**widerlegten** Vorschlags im Satz darüber (`591`, `327`) stehen als Zitat
dieses Vorschlags — **nicht** als Messung;
[`ADR-0082`](../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
widerlegt sie.

**Cluster D ist geteilt — der Trigger der ADR ist eingetreten.** `ADR-0082`
sieht vor: „bei Schicht-Überschreitung wird Cluster D an seiner Grenze geteilt
(fünf)". D trägt vier Träger in **vier Schichten** (Application, Domain, Driven
Adapter, Composition Root) — gemessen bei der Planung von `slice-092`, über die
Block-Position dedupliziert: `application/usecase/*` **22**,
`domain/model` **2**, `adapters/driven/telemetry` **2**, `bootstrap`-Rest
**≈28**, Summe **54** (die `53–55` der ADR hält). Die Naht liegt dort, wo auch
der **Test-Stil** wechselt: D1 fährt einen Fake-Port, D2 baut die Verdrahtung.

**Geschnitten wird nach dem Maß, nicht auf Vorrat** (Modul 5: Plan und
Implementation alternieren): Cluster B zuerst — er ist der reinste (60 von 60
erreichbar) und trägt am wenigsten Kopplung; die folgenden entstehen, wenn er
liegt. **Ihre Adresse ist dieses Maß** (`ADR-0082`), ihre Kennungen entstehen mit
ihren Dateien — kein Name ohne Adresse, die die Sendung annimmt (die Klasse aus
`BEO-PGC/aufschub-adresse-verfaellt`).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Wird blockiert von:** `slice-079` (der Scope-Schnitt) — ohne ihn misst die
  Welle gegen einen Nenner, den [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  ersetzt hat.
- **Blockiert:** nichts.
- **Verwandt, aber nicht Teil:** die DB-gestützte Messung (`ADR-0071` Punkt 3)
  und die Executor-Naht (Punkt 5) — beide eigene Vorgänge. Die Naht würde die
  Decke heben, ohne das Verfahren zu ändern; sie ist deshalb ausdrücklich
  **nicht** der Weg zu diesen 80 %.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Die Endstufe 80 % selbst** — sie steht ([`ADR-0054`](../adr/0054-coverage-gate-und-benchmark-infrastruktur.md),
  [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md));
  sie zu ändern wäre eine Schwellen-Senkung und nach `AGENTS.md` §3.6
  ADR-pflichtig. Diese Welle **erreicht** sie, sie verhandelt sie nicht.
- **Die DB-gestützte Ebene** (`ADR-0071` Punkt 3): eigene, subjekt-qualifizierte
  Messung im nicht-blockierenden `e2e.yml`, eigener Vorgang. `make gates` bekommt
  in dieser Welle **keinen** Container.
- **Die Executor-Naht** (`ADR-0071` Punkt 5): design-begründeter eigener Vorgang,
  nicht Beigabe. Ein Slice dieser Welle ändert **kein Produktionsverhalten**, um
  Coverage zu gewinnen — er prüft Verhalten.
- **Der Messmechanismus** (`-coverpkg`, Docker-Stage, Gate-Skript): mit
  `slice-079` eingestellt. Wer ihn weiter ändern will, braucht eine Entscheidung,
  keinen Test-Slice.
- **Die Integrations-/E2E-Fläche** (`test/integration/**`, die Compose-Kette):
  sie hat ihren eigenen Beleg-Träger (`make test-integration`), und ihre Coverage
  ist nicht Gegenstand dieses Maßes (`ADR-0054` §(a), `ADR-0071` Punkt 1).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-<NN>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
