# Slice slice-087: E2E-Beleg für `GET /changes` — der Träger der Fitness Function

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (ein realer
Beleg), kein repo-weites *Mehr*.

**Bezug:** [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)
§Fitness Function (die Zeile, die dieser Slice einlöst) und §Testabdeckung
(dort als „Erwartung, keine abschließende Festlegung" geführt) ·
[`ADR-0057`](../../adr/0057-http-grpc-api.md) (der bestehende HTTP-Rundlauf,
den dieser Slice erweitert) · [`LH-FA-SST-006`](../../../../spec/lastenheft.md)
(deren Happy-Path-Akzeptanzkriterium einen **Netzwerk-Client** verlangt) ·
`ADR-0068` (die Wegwerf-Clients und ihre begrenzte Import-Berechtigung) ·
[`ADR-0030`](../../adr/0030-testpyramide.md) (das E2E-Tier).

**Berührte Spec-Stellen:** — (dieser Slice baut einen **Beleg**; er ändert
keinen Vertrag).

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

**Ziel:** `GET /changes` bekommt **denselben Beleg wie seine neun Geschwister**:
einen realen Aufruf gegen den **laufenden** Feed-Container im E2E-Rundlauf
(`make test-integration`), der das Lesen über den Netzwerkweg beweist — und der
damit die Fitness-Function-Zeile aus [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)
**einlöst**. Heute hat sie **keinen Träger**: `grep -rn '"/changes"' tools/ test/`
findet nur `/changes/stream` (Review `review-slice-086` F-1).

**Warum ein eigener Slice und nicht ein Zusatz zu `slice-086`:** `slice-086`
liefert den **Endpunkt**, dieser Slice seinen **Beleg** — zwei Lieferwerte. Und
`slice-086` ist bereits einmal über seine §3-Liste hinausgewachsen; der Grund,
warum fremde Arbeit nicht in einen laufenden Slice gehört, steht in dessen §7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein neuer Wegwerf-Client.** Der bestehende `tools/harness/httpclient` wird
  **erweitert**, nicht verdoppelt — `ADR-0068` gibt die begrenzte
  Import-Berechtigung je Client, und ein zweiter Client für dieselbe
  Schnittstelle wäre eine zweite Drift-Fläche.
- **Ein neuer Endpunkt oder eine Vertragsänderung.** Der Endpunkt existiert
  (`slice-086`); dieser Slice **belegt** ihn nur.
- **Die Adapter-Tests** (`internal/adapters/driving/http/**`) — sie tragen die
  Feld- und Fehlerform und sind mit `slice-086` abgenommen; dieser Slice
  wiederholt sie nicht.
- **Ein Gate.** `make test-integration` ist kein Gate (`ADR-0030`); daran
  ändert dieser Slice nichts.

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

**Liefer-Punkt 1 — der Client kann das Lesen.**

- [x] `tools/harness/httpclient` ruft zusätzlich `GET /changes` gegen die
      Basis-URL — mit `source` und optional `schema`/`table`, `from`/`to`,
      `limit` — und prüft die Antwort **inhaltlich**, nicht nur ihren Status.
      — Beleg: `readChanges` (`tools/harness/httpclient/main.go`) baut die
      Abfrage aus `source`/`schema`/`table`/`from`/`to`/`limit` (leerer Wert
      lässt den Parameter weg) und prüft: gesetzte, nicht leere
      `changes`-Liste; je Eintrag nicht leere `change_id`, bekannter
      Operationswert (`INSERT`/`UPDATE`/`DELETE`), `commit_position` ≥ 1 und
      die zum Filter passende Klartext-Identität (`schema`/`table`); die
      Einträge in nicht absteigender `commit_position`.
- [x] Er bleibt **ein** Client (`ADR-0068`): kein zweites Programm, keine
      zweite Import-Berechtigung, kein neuer Eintrag in `.a-check.yml`.
      — Beleg: derselbe `tools/harness/httpclient`; der neue `net/url`-Import
      ist Standardbibliothek unter dem bestehenden Glob
      `tooling: ["tools/harness/**"]`, `.a-check.yml` unberührt.

**Liefer-Punkt 2 — der Rundlauf beweist es.**

- [x] Die HTTP-Rundlauf-Phase des Runners ruft ihn real gegen den **laufenden**
      Feed-Container und wertet das Ergebnis aus; scheitert er, endet die Phase
      **rot**. — Beleg: der volle `make test-integration`-Lauf ist auf Exit 0
      gelaufen; die Phase fügt eine eigene Zeile (`id=285`,
      `HttpChangesReadE2ESentinel`) ein, wartet ihre Erfassung über
      `cdc.changes` ab und liest sie dann über `GET /changes`
      (`READ changes=1 table=feed_e2e_full schema=public operation=INSERT …`),
      deren `change_id` unabhängig gegen `cdc.changes` gehalten wird. `set +e`
      um den Client-Aufruf macht den Ausgang auswertbar — sonst beendete ein
      fehlgeschlagenes `docker run` die Zuweisung selbst und die Phase endete
      rot **ohne Ausgabe**.
- [x] `abdeckung_declare` trägt den erweiterten Nachweis — Anker und
      Kurzbeschreibung nennen das Lesen. — Beleg: Anker ist die Echo-Zeile der
      Phase (`… GET /changes real per HTTP mit reader-Token (die eigens
      eingefügte Zeile …`), die Kurzbeschreibung nennt das Lesen; die erzeugte
      Zeile steht in `docs/user/e2e-abdeckung.md`.
- [x] **Eine rote Gegenprobe:** der Beleg wird **rot gesehen**, wenn der
      Endpunkt nicht antwortet (etwa ein falscher Pfad) — der Beweis, dass der
      Beleg den **Endpunkt** prüft und nicht sich selbst. — Beleg: mit
      `GET /changes` → `/changes-gegenprobe` im Client endete der volle Lauf
      rot (Exit 2; `httpclient: GET /changes (reader) fehlgeschlagen: status
      404 (erwartet 200): 404 page not found` → `endete mit Ausgang 1`),
      danach zurückgenommen und grün wiederholt.

**Liefer-Punkt 3 — der Beleg ist sichtbar.**

- [x] Die E2E-Abdeckungstabelle (`docs/user/e2e-abdeckung.md`) ist durch ein
      volles `make test-integration` real **neu erzeugt** und trägt den
      Nachweis; die Zeilen-Drift ist mitgezogen. — Beleg: der finale Lauf
      (Exit 0) hat die Tabelle geschrieben (`13 Go-Zeilen und 24
      Bash-Zeilen`); die HTTP-Zeile zeigt auf
      `tools/harness/run-integration-tests.sh:2061` und ihre
      Kurzbeschreibung nennt das Lesen.
- [x] `harness/README.md` nennt den Beleg in der Aufzählung des
      `make test-integration`-Ziels — dort steht heute **jeder** E2E-Nachweis.
      — Beleg: die Zeile des `make test-integration`-Ziels trägt den
      Lese-Beleg des Wegwerf-Clients (`· seit slice-087`).
- [x] `make gates` grün (Exit direkt, ungepiped). — Beleg: `make gates`
      Exit 0 (baseline-verify, docs-check, a-check, commit-traceability,
      coverage-gate 72,00 % ≥ 70 %, Nachweis-Stempel).

- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-087.md` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      0 HIGH/0 MEDIUM, kein Rückgabe-Pfeil (F-1/F-2 INFO, F-3 LOW).
- [x] Verifikation durchgeführt, Report unter
      `docs/reviews/verify-slice-087.md` liegt vor (Modul 11, frischer Kontext) —
      **closure-fähig**; `ADR-0081`s Fitness-Function-Zeile **eingelöst**. Der
      Verifier hat die `READ`-Zeile aus **eigenem** vollem Lauf gesehen, die
      Abdeckungstabelle **byte-gleich** nach dem Lauf bestätigt und eine eigene
      Mutation gefahren (fabrizierte `change_id` → rot, gefangen von der Bindung
      gegen `cdc.changes`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) — **entfällt**: dieses
      Repo führt die Datei nicht (Greenfield-Bootstrap, kein Inventur-Fund).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — ein Beleg
      ergänzt (`negativtest-ohne-bindung-an-seine-eingabe/evidence/slice-087.md`,
      Zähler damit **2×**), **kein** neues Verzeichnis, **kein Zähler gesetzt**.
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

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/httpclient/main.go` | update | der zusätzliche `GET /changes`-Aufruf samt **inhaltlicher** Auswertung; es bleibt **ein** Client (`ADR-0068`) |
| `tools/harness/run-integration-tests.sh` | update | die HTTP-Rundlauf-Phase ruft ihn, wertet real aus, und `abdeckung_declare` trägt den Nachweis |
| `docs/user/e2e-abdeckung.md` · `harness/README.md` | update | Erzeugnis des Laufs (Zeilennummern bzw. die Beleg-Aufzählung) |

**Der genaue Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste nennt
die Träger. Wer sie erweitert, prüft die Größenregel (≤ 3 Liefer-Punkte).

**Nicht in dieser Liste:** `internal/**` (der Endpunkt ist mit `slice-086`
abgenommen), `spec/**` (kein Vertrag berührt), `tools/schema/**`,
`.a-check.yml` (keine neue Gruppe, keine neue Kante).

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-086` liegt in `done/` — **der
Endpunkt muss existieren**, sonst prüfte der Beleg nichts —,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich, dass der
  Beleg eine **eigene Phase** braucht (etwa weil der bestehende Rundlauf zu
  einem Zeitpunkt läuft, an dem kein Change mehr erzeugt wird), gehört er zurück
  zur Zerlegung — nicht in eine wachsende Phase.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass `GET /changes`
  im Compose-Betrieb **nicht erreichbar** ist (Verdrahtung, Token-Klasse), ist
  das ein Blocker mit Entscheidung — `slice-086` hat den Endpunkt verdrahtet,
  dieser Slice prüft es real.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** der Beleg läuft real **grün** **und** er ist **rot
gesehen**, wenn der Endpunkt fehlt **und** die Abdeckungstabelle ist aus einem
vollen Lauf neu erzeugt **und** `make gates` grün **und** die Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Beleg könnte sich selbst prüfen.** Prüft er nur seinen eigenen
  Exit-Code, ist er grün ohne Aussage — dieselbe Klasse wie ein Fake, der nichts
  prüft. Wächter: Liefer-Punkt 2 verlangt die **rote Gegenprobe**.
  — **Ausgang:** *entfallen — gestrichen mit Begründung*: **zwei unabhängige**
  Rot-Belege wurden gesehen. Der Implementer hat den Pfad mutiert (`404 page
  not found`, Exit 2) — und im roten Lauf registrierte und listete der Client
  **erfolgreich**, der Rot-Beleg trifft also das **Lesen**, nicht den Prozess.
  Der Verifier hat eine **fabrizierte `change_id`** in die READ-Zeile gesetzt →
  Exit 2, gefangen von der unabhängigen Bindung gegen `cdc.changes`
  („nicht real über cdc.changes lesbar"). Die Zeile ist damit nicht
  selbstgenügsam; der Netzwerk-Leseweg hängt am Store.
- **Er könnte an der Zeitachse scheitern.** Der HTTP-Rundlauf läuft **spät** im
  E2E-Ablauf (nach dem NATS-Negative-Beleg); erzeugt zu diesem Zeitpunkt niemand
  mehr einen Change, findet das Lesen nichts. Der Beleg muss das einplanen —
  eine eigene Zeile erzeugen oder einen belegten Bestand lesen.
  — **Ausgang:** *entfallen — gestrichen mit Begründung*: die Phase erzeugt ihr
  Subjekt **selbst** (`INSERT` der Sentinel-Zeile) und pollt auf ihren
  `commit_position`, bevor sie mit `[from, to)` genau darauf liest. Vier volle
  Läufe, keine Bestandsabhängigkeit.
- **Die Abdeckungstabelle könnte erneut driften.** Sie bindet an `Datei:Zeile`,
  und dieser Slice verschiebt **beides** (Client und Runner). Deshalb ist der
  volle Lauf Teil der DoD, nicht ein `cmp` gegen die alte Datei.
  — **Ausgang:** *entfallen — gestrichen mit Begründung*: nach einem **vollen**
  Lauf ist die committete Datei **byte-gleich** (sha256 unverändert, `git status`
  leer) — sie **ist** das Generator-Erzeugnis, keine Behauptung darüber. Die
  `Ort`-Werte wurden zusätzlich unabhängig nachgerechnet (24/24).

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

- **Was hat funktioniert:** **Eine Zusage in einem unveränderlichen Dokument
  hat jetzt ihren Träger.** Der Weg dorthin ist die eigentliche Leistung: der
  Reviewer hat die ungetragene Fitness-Function-Zeile in
  [`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md) gefunden, der
  Planner hat sie **nicht** in den laufenden Slice gezogen, sondern als eigenen
  geschnitten (`slice-086` war schon einmal gewachsen), und der Implementer hat
  sie gebaut **mit** roter Gegenprobe.
  Und: der Verifier hat „die Tabelle ist das Erzeugnis" zu einer **Messung**
  gemacht statt zu einer Aussage — nach einem vollen Lauf ist die Datei
  **byte-gleich**, und die 24 Anker hat er unabhängig nachgerechnet.
- **Was ging anders als geplant:** (1) Der Implementer hat unterwegs einen
  **echten Fehler** gefunden und behoben: die Phase rief `x=$(docker run …)`
  unter `set -e` — ein fehlgeschlagener `docker run` beendete die **Zuweisung**
  selbst, die Statusprüfung lief nie, und der Lauf endete **rot ohne Ausgabe**.
  Aufgedeckt hat das die Gegenprobe, nicht das Lesen. (2) Der Review hat mit
  einer eigenen Mutation zwei **Grenzen der Aussagekraft** gefunden (F-1: der
  Beleg bindet die **Filterachse** nicht — bei Bereichsbreite 1 ist „Filter
  wirkt" von „Filter fehlt" ununterscheidbar; F-2: Ordnungs- und Limit-Prüfung
  sind in dieser Aufruf-Form strukturell unausübbar). Beide **benannt**, keiner
  closure-blockierend. (3) Die DoD nennt `coverage-gate 72,00 %`, der Verifier
  maß **71,90 %** auf demselben Commit — beide über der Schwelle; die exakte
  Zahl ist **lauf-gebunden**, und das ist benannt statt eingeebnet.
- **Lerneintrag (geschärfte Regel):** *Ein Beleg, der seine Aussage nicht an
  seine **Eingabeseite** bindet, ist grün ohne Aussage.* Das ist die **zweite**
  Gelegenheit in zwei Tagen an einem **anderen Gegenstand** — `slice-086`: ein
  **Negativtest**, dessen Ablehnung an keinem Eingabewert hing; `slice-087`:
  ein **E2E-Beleg**, dessen Filterprüfung nicht am Filter hing. Beide Male hat
  es **kein Lesen** gefunden, sondern nur das Mutieren der Eingabeseite. Die
  Regel aus `slice-086` steht damit **über Vorgänge belegt**: jede Zusage wird
  an ihrer Eingabeseite mutiert, nicht nur an ihrer Ausgabeseite.
  **Verkörperung offen** (Träger käme im Reviewer-Skill in Frage); als Kandidat
  geführt, nicht behauptet.
- **Steering-Loop-Eintrag:** **keine Verkörperung durch diesen Slice.**
  Eine Registerbewegung: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
  bekommt einen **zweiten Beleg** (`slice-087`) → **2×**. Dabei **benannt**:
  der Eintrag ist **enger benannt als sein Gegenstand** — sein Name und seine
  `observation.md` sprechen von einem *Negativtest*, der zweite Fall ist ein
  *E2E-Beleg*. Der Zähler zählt den **Mechanismus**, die Benennung stammt vom
  Erstauftreten und ist `observation.md`-unveränderlich; die Klarstellung steht
  in `state.md`.
- **Beobachtungs-Register (`../observations/`):** ein Beleg ergänzt
  (`negativtest-ohne-bindung-an-seine-eingabe/evidence/slice-087.md`, Zähler
  **2×**), **kein** neues Verzeichnis, **kein** Zähler gesetzt.
  **Benannt, nicht gezählt:** die unter `set -e` unerreichbare Fehlerprüfung
  (`slice-087`, im Slice gefunden und behoben, Erstauftreten) — sie gehört als
  eigener Kandidat in die Closure, nicht in diesen Zähler.
- **Folge-Slices:** keiner. F-1 und F-2 sind **Grenzen der Aussagekraft**, keine
  fehlende Lieferung — sie sind benannt und gehen über die Finding-Klassen in
  den Zähler.
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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt den Runner, den Wegwerf-Client,
das Erzeugnis und die README-Zeile in **einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen; die Zählerstände sind am Register
**nachgezählt** (`ls evidence/`), nicht aus diesem Text übernommen:

- `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (**3×**, Schwelle erreicht —
  Ausgang beim Lese-Schritt der `welle-20`-Closure): **Treffer, und ein
  genauer.** `docs/user/e2e-abdeckung.md` ist eines der drei namentlich
  geführten Erzeugnisse dieses Eintrags; dieser Slice **erzeugt sie neu** und
  zieht die Drift nach. Er behebt damit **nicht** die Bindung — die bleibt
  Handarbeit —, aber er ist der Vorgang, an dem die Klasse zuletzt auftrat.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (1×, offen):
  **benachbart** — die Tabelle trägt Zeilennummern, also Zahlen gegen einen
  beweglichen Gegenstand; die Zuordnung trifft die Closure.

**Ergebnis** (Stand: Anlage dieses Plans): kein Eintrag rückt mit diesem Slice
auf 3× — `generierte-artefakte-ohne-sync-sensor` steht bereits darüber.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (der E2E-Rundlauf ist über `ADR-0030` und
`ADR-0057` verankert, die Client-Grenze über `ADR-0068`), **Phase-Reife** hoch,
**Evidenz-/Diskrepanz-Risiko** **niedrig** — der Endpunkt ist mit `slice-086`
abgenommen und verdrahtet —, **Reconciliation-Aufwand** null.
