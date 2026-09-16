# Slice slice-091: Coverage Cluster C — Zustell- und Betriebs-Rand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-20` — Coverage 80 % über der netzlos prüfbaren Fläche. Cluster
**C** ist der zweite der vier; Cluster B (`slice-088`) liegt in `done/`.

**Bezug:** [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(das **Schnittmaß** — die Cluster-Tabelle und die Zahl `62`) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand) · [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(die Rampe 70 → 80) · [`ADR-0057`](../../adr/0057-http-grpc-api.md) und
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (die zwei
Zustell-Oberflächen, deren Pakete hier liegen) ·
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md)
(`natsnotify`) · `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` und
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Zähler-Stände führen die
Einträge selbst — ein Verweis braucht keine Zahl, die altern kann).

**Berührte Spec-Stellen:** `LH-FA-SST-006`, `LH-FA-SST-007`, `LH-FA-SST-008` —
die Zusagen, die die vier Pakete tragen; dieser Slice prüft sie zusätzlich, er
ändert sie nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-16.

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

**Ziel:** Die vier Pakete des Clusters **C** — `driving/http`, `driving/grpc`,
`driven/natsnotify`, `driving/grpc/streamv1` — werden über **netzlose** Tests
gedeckt. Das Maß ist das der Welle ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)):
**62 ungedeckte Statements, davon ≈52 netzlos erreichbar.** Die restlichen
≈10 sind **nicht** erreichbar und werden **benannt**, nicht kaschiert.

**Was dieser Slice liefert:** Tests. Er ändert **keinen** Produkt-Code, außer
einer Naht, die ein Test wirklich braucht — und dann nach
[`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md), nicht
frei gewählt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Das Anheben von `THRESHOLD`.** Die Rampe 70 → 80 ist ein eigener Schritt
  der **Wellen-Closure** ([`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  §(a), [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md));
  ein Slice, der seine eigene Messlatte mitzieht, könnte sein Grün nicht mehr
  belegen. Die Welle verlangt dafür **zwei** Belege: grün bei 80 **und** rot bei
  85.
- **Die Cluster D und A** ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)):
  eigene Slices. Die Welle schneidet **nach dem Maß, nicht auf Vorrat** (Modul 5
  — Plan und Implementation alternieren); C ist der zweite, D und A entstehen,
  wenn es hier liegt.
- **Die DB-Adapter-Coverage** ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3). Das ist ein **anderer** Messgegenstand; die zwei Zahlen
  **partitionieren** denselben Code, und ein Test, der hier zählt, verändert dort
  nichts. Wer beides in einem Zug bewegt, kann nicht mehr sagen, welche Zahl
  sich warum bewegt hat.
- **Eine Änderung am Test-Runner** (`tools/harness/run-integration-tests.sh`).
  Ein anderer Vorgang: dort geht es um **stille Ausschlüsse** von
  Integrationstests, nicht um die netzlose Fläche.
- **Tests, die die Zahl heben, ohne eine Zusage zu prüfen.** Sie sind der
  naheliegende Fehler dieses Slice und stehen als Risiko in §6, nicht als
  Liefer-Punkt.

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

- [x] **LP1 — die Tests existieren und sind netzlos grün.** Für die vier Pakete
      des Clusters C liegen Tests vor, die der Gate-Lauf **wirklich fährt**
      (`make test`, netzlos); `make gates` ist grün. **Der Zuwachs wird als Zahl
      mit ihrem Lauf genannt**, nicht als „deutlich besser" — ein DoD-Kriterium,
      das eine ungemessene Zahl behauptet, ist der Fehler aus
      `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (4×).
- [x] **LP2 — die Negativtests binden ihre Ablehnung an die Eingabe.** Wo ein
      Test eine Ablehnung prüft (`401`, `Unauthenticated`, verweigerte
      Publikation), ist sie an **den Eingabewert** gebunden, der sie auslösen
      soll — nicht an einen Fake, der sie unabhängig von der Abfrage liefert.
      Auslöser: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**4×**).
- [x] **LP3 — die unerreichbaren Statements sind benannt.** Jedes der vier
      Pakete, das danach noch ungedeckte Statements hat, nennt sie **einzeln mit
      dem Grund**, warum sie netzlos nicht erreichbar sind (lebender Dienst,
      Zeitabhängigkeit). „Rest nicht erreichbar" ohne Namen gilt als **nicht
      erfüllt** ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
      nennt für C 62 ungedeckt gegen ≈52 erreichbar — die Differenz ist die
      benannte Lücke, nicht die stille).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-091.md`
      liegt vor (Modul 11, frischer Kontext).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — dieses Repo führt Wellen-Betrieb; die Prüfung fällt der `welle-20`-Closure zu.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/http/**` | Test neu | Cluster C, Teil 1 — die HTTP-Zustell-Fläche (`ADR-0057`). |
| `internal/adapters/driving/grpc/**` | Test neu | Cluster C, Teil 2 — die gRPC-Fläche (`ADR-0060`). |
| `internal/adapters/driven/natsnotify/**` | Test neu | Cluster C, Teil 3 — der Wecksignal-Adapter (`ADR-0055`). |
| `internal/adapters/driving/grpc/streamv1/**` | Test neu, **falls** dort prüfbarer Code liegt | Cluster C, Teil 4. Der generierte Stub ist **nicht** Gegenstand — sein Sync trägt seit `slice-090` `make generated-sync`; hier zählt nur, was ein Test wirklich prüft. |
| `docs/plan/planning/welle-20.md` §4 | **nicht** | Die Cluster-Tabelle dort trägt die **Soll**-Zahlen aus [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md). Die **erreichte** Zahl gehört in die Closure-Notiz dieses Slice, mit ihrem Lauf — eine zweite Fassung in der Welle wäre genau die driftende Kopie. |
| `harness/sensors/coverage-gate.md` | update, **nur falls** eine Zahl dort gegen die Messung driftet | Der Sensor-Träger des Messgegenstands. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-088` (Cluster B) liegt in `done/` und
`docs/plan/planning/welle-20.md` ist offen. Beides ist ohne Rückfrage
feststellbar — die Welle schneidet **nach dem Maß, nicht auf Vorrat**.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn ein Paket eine
  **Naht im Produkt-Code** verlangt, die über die Form aus
  [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) hinausgeht.
  Dann ist der Schnitt falsch: Ein Test-Slice, der die Produktion umbaut, prüft
  nicht mehr, sondern ändert — und das ist ein eigener Vorgang mit eigener
  Entscheidung.
- `in-progress` → `open` (blockiert — Carveout?): wenn die **netzlos
  erreichbare** Fläche real weit unter den ≈52 Statements liegt, die
  [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  für C ausweist. Dann trägt Cluster C die Welle nicht, und die Antwort ist eine
  **Neu-Bemessung des Schnittmaßes** (Architect), kein Carveout: Der Carveout
  setzte eine Schwelle herab, hier wäre die Schwelle richtig und nur die
  Rechnung falsch.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 aus §2 sind real belegt — die Tests
laufen **im Gate** (nicht nur lokal daneben), die Negativtests binden ihre
Ablehnung an die Eingabe, und die verbleibenden ungedeckten Statements der vier
Pakete sind **einzeln benannt** — **und** `make gates` ist grün, mit dem
Zuwachs als Zahl samt ihrem Lauf.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Der naheliegende Kandidat ist die Zahl selbst: Diese Welle misst ihren
Fortschritt an einer **beweglichen** Quote, und
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` steht bei **6×** — ob dieser
Slice dazu etwas Neues beiträgt oder die bestehende Regel nur erneut anwendet,
entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die ≈52 sind eine Über-Schätzung** — dann liefert C weniger als die Welle
  veranschlagt, und das Budget aus B+C+D hat nur **13 Statements Puffer**
  ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)).
  — **Ausgang: eingetreten — und begrenzt.** Real erreichbar waren **47** statt
  ≈52; die fünf Differenz-Statements sind **einzeln benannt** (drei tote
  defensive Doppelprüfungen in `driving/http`: `readchanges.go:181.3,182.1` und
  `retention.go:47.4,49.1`, die eine Invariante erneut prüfen, die der Aufrufer
  schon erzwungen hat; zwei im Publish-Erfolgspfad von `natsnotify`,
  `notify.go:136.2,137.12`). Die Rückführung `in-progress → open` greift
  **nicht**: §4 bindet sie an „real weit unter ≈52", 47 von 52 sind **90 %**.
  Wirkung auf die Welle: der Puffer B+C+D sinkt von 13 auf **8** Statements —
  eine Planner-/Wellen-Entscheidung, kein Befund dieses Slice.
- **Ein Test wird zeitabhängig** — dann flappt die Zahl, und der Flap ist für
  Tests und diff-skopiertes Review unsichtbar
  (`BEO-PGC/test-integration-retention-timing-flake`, 2×; der `slice-090`-Beleg
  fand denselben Gegenstand an der Coverage-Zahl). — **Ausgang: entfallen für
  die Tests dieses Slice — und bestätigt für den Bestand.** Gemessen:
  `go test -race -count=20` über die vier Pakete, 4 × `ok`, **keine** Frist
  gefeuert; die vier neuen Fristen sind als **Hänge-Schutze** gesetzt (30 s
  gegen ≈0,29 s je Iteration) und nicht als Zusicherung über eine Uhr. Der
  **vorbestehende** Flake hat sich in diesem Vorgang trotzdem gezeigt — zwischen
  zwei Läufen **desselben** Commits (77,20 gegen 77,30 %; 2/1903 = 0,105 pp) —;
  er ist nicht Gegenstand dieses Slice, sondern wandert als Beleg in den
  Register-Eintrag.
- **Ein Negativtest bindet die Ablehnung an den Fake statt an die Eingabe** —
  der Test ist grün, egal was der Adapter mit dem Wert macht, und die Zahl
  steigt trotzdem. Das ist die Klasse mit **4×** und der wahrscheinlichste
  Fehler dieses Slice. — **Ausgang: eingetreten — und behoben.** Die Verifikation
  fand die Form an einem **vorbestehenden** Test (`server_test.go:161`, der
  sechste `…UngueltigesJSONEndetMit400`: grün bei zerstörter Zusage); er ist
  gebunden und durch Mutation belegt (rot) mit Kontrolle (grün).
  **Nicht dieses Risiko, sondern sein Nachbar** ist der zweite Fund derselben
  Runde: der `LogPort`-Test (`notify_test.go`) behauptete mehr, als er band —
  dort trägt der **Beleg** seinen Satz nicht
  (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`), nicht die Ablehnung ihren
  Eingabewert. Die zwei Klassen liegen dicht beieinander und wurden deshalb
  getrennt geführt, nicht als ein Fall gezählt.
- **Coverage-Theater** — Tests, die Statements durchlaufen, ohne eine Zusage zu
  prüfen: Die Zahl steigt, die Prüf-Kraft nicht, und die Welle hätte ihr Ziel
  formal erreicht. — **Ausgang: entfallen — gemessen, nicht behauptet.** Der
  Review hat **24** Mutationsproben gefahren (20 rot), der Verifier **13** an
  der Produktionsseite (**13 rot**); jede neue Zusage dieses Slice färbt bei
  Zerstörung ihrer Produktionsseite rot. Ein Test, der nur die Zahl hebt, wäre
  in dieser Probe grün geblieben.

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

- **Was hat funktioniert:** Die **Mutation** war das Rückgrat dieses Slice — nicht
  der Zähler. Der Review hat **24** Proben gefahren, der Verifier **13** an der
  Produktionsseite; diese Methode hat (a) die zwei Bindungs-Lücken gefunden
  (einen neuen und einen vorbestehenden Negativtest, die grün blieben, wenn man
  ihre Zusage zerstört), (b) das vierte §6-Risiko *widerlegt* statt es zu
  beschwichtigen, und (c) die Reproduzierbarkeit der Mutationsangaben erzwungen
  — was zwei Kommentare als unwahr entlarvte. Zweitens: die Trennung von
  **Zustand** und **Lauf** hat getragen — Cluster C ist **62 → 15** ungedeckt
  (eine Statement-Differenz, stabil), während die Prozentzahl ein **Band** ist
  (77,1–77,3 %; die Enden aus dem Lauf `slice-091` und den 47 Läufen der
  Verifikation), und beide Zahlen tragen ihren Lauf. Die **untere Kante** des
  Bands ist ein **seltener** Fall (1 von 47 Läufen der Verifikation, 0 von 54 im
  Delta-Review); die typischen Enden sind 77,2 und 77,3 %.
- **Was ging anders als geplant:** **Vier** Runden statt einer, und der Grund
  liegt nicht im Slice, sondern in seiner Umgebung. Die schärfste Beobachtung:
  **dieser Slice hat drei Sätze in einem Dokument falsch gemacht, das er nie
  angefasst hat.** Er gab `streamv1` eine Testdatei — damit wurde die stehende
  Liste „Fünf Pakete … haben keine Testdatei", ihr Schlusssatz (`coverage: 0.0%`)
  und die Gruppierung selbst falsch, ohne dass jemand diese Datei im Diff hatte.
  Gefunden hat das erst der **Delta-Review**, nach zwei Runden; die *Behebung*
  hat dann zunächst eine **falsche Zählung** tragend gemacht (D-1), weil sie die
  Gruppe mit `TestGoFiles` begründete, das 23 von 31 Paketen trifft. Zweitens:
  die ≈52 aus [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  waren real **47**. Drittens: die Ersetzung der falschen Flap-Ursache hat einen
  **zweiten** Flapper gefunden, der das dritte, unkolokalisierte Ende der
  Verifikation erklärt — der Fund kam erst durch die Korrektur.
  **Und eine Form-Lücke, die dieser Slice nicht mehr schließt:** `a7d7f7b` und
  `af3ea9f` sind von **keinem** Delta-Review gedeckt. Der Report über die zwei
  Fixrunden (`review-slice-091-delta`) prüfte `f9cd5e4` und `a7d7f7b`; die
  vierte Runde (`af3ea9f`) führt wörtlich die vom Delta-Review vorgeschriebene
  Ein-Stelle-Korrektur aus, und ihre Zahlen sind von der Verifikation geprüft —
  die **Form** des zweiten Delta-Pfeils fehlt trotzdem. Das ist dieselbe Lücke,
  die `slice-090` als V-2 führte und dort geschlossen hat; hier bleibt sie
  **benannt** statt still.
- **Steering-Loop-Eintrag:** **kein neuer Träger — die Leser-Hälfte hat
  getragen.** Die Regel, die alle vier Runden deckt, steht
  (`AGENTS.md` §3.12 Instanz A und B), und ihre durchsetzenden Leser sind
  Reviewer und Verifier. Was hier **neu** hinzukommt, ist eine Beobachtung ohne
  Zielort: **die Arbeit überholt einen Träger, den sie nicht anfasst.** Wer eine
  *gemessene Eigenschaft* eines Gegenstands bewegt (hier: ob ein Paket eine
  Testdatei hat), macht die Träger falsch, die diese Eigenschaft **beschreiben**
  — und die stehen nicht im Diff und werden von keinem Sensor gelesen. Sie steht
  als `BEO-PGC/arbeit-ueberholt-stehenden-traeger` im Register (1×), **nicht**
  als verkörperte Regel: ein Gate müsste dafür wissen, welche Sätze von welcher
  Eigenschaft abhängen, und das weiß es nicht.
  Auslöser: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (`slice-091` — 1×).
- **Beobachtungs-Register (`../observations/`):** **ein Verzeichnis neu
  angelegt** (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, 1×) und **drei
  Belege** ergänzt: `test-integration-retention-timing-flake` → **3×**
  (Schwelle erreicht), `beleg-befehl-traegt-seinen-satz-nicht` → **3×**
  (Schwelle erreicht), `negativtest-ohne-bindung-an-seine-eingabe` → **5×**.
  **Kein Zähler wird gesetzt** — jeder folgt aus der Zahl der Dateien unter
  `evidence/`. Die zwei neuen Schwellen-Einträge weist der **Lese-Schritt der
  `welle-20`-Closure** zu (Modul 6), nicht dieser Slice.
- **Folge-Slices:** keine Datei in `open/` — die Cluster **D** und **A** sind
  die nächsten Schnitte **derselben Welle** ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md));
  sie entstehen nach dem Maß, wenn dieser liegt. `ADR-0082`s
  Folge-Slice-Vorschlag (eine Quelle für Erzeuger und Prüfer der
  E2E-Abdeckungstabelle) bleibt **unadressiert** — die ADR nennt ihn ohne
  Kennung, und sein Re-Evaluierungs-Trigger ist nicht eingetreten.
- **Risiken aus §6:** vier, je ein Ausgang — R1 *eingetreten und begrenzt*
  (47 statt ≈52, Puffer 13 → 8), R2 *entfallen für die Tests dieses Slice*
  (`-race -count=20`, keine Frist gefeuert), R3 *eingetreten und behoben*
  (ein vorbestehender Test; der zweite Fund derselben Runde gehört zur
  Nachbar-Klasse), R4 *entfallen* (gemessen über 37 Mutationsproben — die Zahl
  ist eine **Untergrenze**: sie zählt die zwei großen Berichte, die
  Delta-Runden prüften acht weitere Proben).
- **Drei Paarungen:** dieses Repo führt **Wellen-Betrieb**; die Prüfung fällt der
  `welle-20`-Closure zu (Modul 6 Schritt 3c). Vorab geprüft: die vier
  Register-Adressen dieses Slice existieren als Verzeichnis, und **jedes** der
  vier ergänzten führt ein nicht leeres `evidence/`.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die vier
Adapter-Pakete des Clusters C und ihre Tests — durchweg die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Eine feinere
Sub-Area ist dort nicht deklariert, und die vier Pakete als eigene Sub-Areas
auszudifferenzieren wäre falsch: Sie teilen denselben Modus, dieselbe
Konventionen-Dichte und dieselbe Inventur-Lage — die Schwelle „≥ 2 von 3 Achsen"
wäre für jede erfüllt, ohne etwas zu trennen. Die Deklaration `*`/`PGC` gilt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`observations/BEO-PGC/`). Für die Fläche dieses Slice — **Test-Arbeit an
Adapter-Paketen** — sechs Treffer. Die genannten Zähler-Stände sind der Stand
**bei dieser Planung**; sie sind nach diesem Slice teilweise höher (der Zähler
folgt den `evidence/`-Dateien, nicht dieser Zeile).

- `negativtest-ohne-bindung-an-seine-eingabe` — **4×**, `offen`. **Der
  schärfste Treffer:** er beschreibt genau die Form, die ein Test in `driving/grpc`
  und `driving/http` naheliegend annimmt (Ablehnung gegen einen Fake statt gegen
  den Eingabewert). Er ist als **LP2** in die DoD gezogen, nicht nur als Risiko.
- `zahl-in-traeger-driftet-gegen-die-messung` — **6×**, `offen`. Betrifft die
  Zahlen **dieses Plans** (die `62`/`≈52` aus
  [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md))
  und die Zahl, die dieser Slice in §7 schreiben wird: sie trägt ihren **Lauf**.
- `dod-begruendung-unzutreffende-tatsachenbehauptung` — **4×**, `offen`. Betrifft
  LP1 direkt: keine Zahl im DoD-Kriterium, die nicht gemessen ist.
- `test-integration-retention-timing-flake` — **2×**, `offen`. Betrifft die
  **Zeitabhängigkeit** neuer Tests: der `slice-090`-Beleg fand sie an einer
  Coverage-Zahl, bei durchweg grünen Tests. Steht als Risiko in §6.
- `a-check-null-abdeckung` — **3×**, `offen`. Grenzt die Fläche ab: er betrifft
  die **Layer-Globs** der Architekturprüfung, nicht die netzlose Quote — dieser
  Slice bewegt ihn nicht, und ein grüner a-check-Lauf sagt über ihn nichts.
- `db-gegenstand-enthaelt-netzlos-geprueften-code` — **2×**, `offen`. **Nicht
  getroffen, und das ist die Antwort:** dort geht es um **netzlose** Tests, die
  in den **DB-Adapter**-Gegenstand zählen; dieser Slice zählt in den netzlosen
  Gegenstand und berührt keinen Store-Adapter. Die beiden Gegenstände bleiben
  getrennt — §1 schließt den Übergang ausdrücklich aus.

**Keine Neuanlage durch die Sichtung.** Ob dieser Slice selbst eine Beobachtung
erzeugt, entscheidet der Lauf, nicht der Plan.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**
(Greenfield-Default `*`/`PGC`, siehe Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md)). Der Block pro
Sub-Area entfällt damit; der **Abschnitt** bleibt, weil die zwei vorgelagerten
Prüfungen oben in jedem Slice-Plan laufen.
