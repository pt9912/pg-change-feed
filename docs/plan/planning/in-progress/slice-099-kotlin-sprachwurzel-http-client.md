# Slice slice-099: Kotlin-Sprach-Wurzel und HTTP-Client — die zweite Zelle der Matrix

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). Dieses Repo führt derzeit **keine** offene Welle
(`docs/plan/planning/in-progress/roadmap.md` §Offene Wellen ist leer).

**Bezug:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(Festlegung 1 — Zelle `kotlin`×`http` der Matrix; Festlegung 5 — je Sprache
ein `make`-Ziel, kein Gate; Festlegung 4 — Abhängigkeits-Pins; §Slice-Schnitt-
Empfehlung, Zeile 2: „**Kotlin-Sprach-Wurzel + HTTP-Client** … dieselbe Form,
eigene Werkzeugkette … keine (parallel zu 1)") ·
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) Festlegung 2/3
— **bestätigt, nicht superseded für diesen Teil**: der Bau-Kontext ist das
Sprach-Wurzelverzeichnis, ein eigenes digest-gepinntes Dockerfile, kein Gate —
dieser Slice ist der erste reale Bau, der diese Form für Kotlin einlöst
([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) selbst: eine
Sprache, die zum ersten Mal auftritt, **steht allein**, sie beweist die Form)
· [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(der nicht-blockierende Workflow als Träger-Muster) ·
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `BEO-PGC/handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide treffen den Handbuch-Nachzug dieses Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-006` (die HTTP-/JSON-API) — dieser
Slice **zeigt** sie in einer weiteren Sprache, er ändert sie nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `examples/kotlin/` entsteht als eigenes Sprach-Wurzelverzeichnis mit
eigenem, digest-gepinntem Dockerfile (JVM-Basis, Kandidat
`eclipse-temurin:21-jdk`, Existenz gemessen in
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 4),
gepinntem Build-Werkzeug samt Prüfsumme (Gradle- oder Maven-Wrapper — Form
gehört dem umsetzenden Zug) und dem `make`-Ziel `examples-kotlin`; darin
liegt der erste Kotlin-Client, `examples/kotlin/http-client/` — ein
Anfrage/Antwort-Aufruf gegen die Verwaltungs-API (`GET /tables` mit dem
`reader`-Token, Form-Vorbild `examples/http-client` in Go), mit netzlos
prüfbaren Tests. Ein nicht-blockierender GitHub-Actions-Workflow fährt das
Ziel, und Handbuch samt `harness/README.md` §Werkzeuge tragen die neue Zeile.

**Warum dieser Slice parallel zu `slice-098`, nicht danach:**
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
§Slice-Schnitt-Empfehlung führt die zwei Sprach-Wurzeln als **unabhängige**
Slices: eigene Werkzeugkette, eigenes Dockerfile, keine Abhängigkeit in beide
Richtungen. Ein gemeinsamer Slice für C# und Kotlin würde zwei
Werkzeugketten in einer Sitzung prüfbar machen müssen — über der
Drei-Liefer-Punkte-Grenze.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der SSE-Client in Kotlin.** Folge-Slice `slice-100` übernimmt ihn —
  zusammen mit dem C#-Pendant, weil beide dieselbe Handbuch-Form (ein
  `**Beispiele:**`-Block) und dieselbe Abhängigkeit auf beide Sprach-Wurzeln
  tragen.
- **Der NATS-Client in Kotlin.** Folge-Slice `slice-101` übernimmt ihn, samt
  dem gepinnten `io.nats:jnats`-Paket — aus demselben Grund wie oben.
- **Der gRPC-Client in Kotlin, der benannte Zusatzkontext und der
  Protobuf-/gRPC-Kotlin-Generator.** Folge-Slice `slice-103` übernimmt sie.
  Die zwei Form-Fragen dieser Zelle — **Bau-Kontext, Generator** — stehen
  **nicht** hier: [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  weist sie ausdrücklich dem jeweiligen gRPC-Slice zu (§Slice-Schnitt-
  Empfehlung, „Was der Planner beim Anlegen beachten muss"); die Form ist
  zudem bereits in `slice-102` (C#) gemessen und wird dort **übertragen**,
  nicht neu erfunden.
- **Die C#-Sprach-Wurzel und ihr HTTP-Client.** Sie läuft **parallel** in
  `slice-098` — eigene Werkzeugkette, eigenes Dockerfile, keine Abhängigkeit
  in beide Richtungen.
- **Der Go-gRPC-Client (`examples/grpc-client/`).** Er ist **kein** Teil
  dieser Matrix-Erweiterung um C#/Kotlin — er läuft als eigener Liefer-Punkt
  (LP2) von `slice-095` (aktuell `in-progress/`), gebunden an den Umzug der
  Vertragsfläche (`slice-097`, `done/`). Anderer Vorgang, andere Sprache,
  anderer Träger.
- **`.a-check.yml`.** Gemessen (`grep -n "languages" .a-check.yml`): die
  Gruppe `examples` trägt ausschließlich Go-Semantik, kein `kotlin`-Schlüssel
  existiert und keiner entsteht hier — [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 6 hält das für die **gesamte** Matrix fest. Bestand bleibt
  bewusst stehen.
- **Ein eigenes Gate für das Sprachziel.** `examples-kotlin` bleibt ein
  **Werkzeug** ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 5) — es braucht Netz (Paket-Bezug) und kann deshalb kein
  netzloses Gate werden.
- **Eine Änderung an der HTTP-API selbst.** Der Client **benutzt** `GET
  /tables`; verlangt er eine Vertragsänderung, ist das eine Spec-Änderung,
  kein Beispiel-Umbau (`SPEC-023`).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1 — Sprach-Wurzel, Werkzeugkette, Ziel.** `examples/kotlin/Dockerfile`
      mit digest-gepinnter JVM-Basis (Kandidat `eclipse-temurin:21-jdk`,
      Existenz gemessen in [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
      Festlegung 4 — Digest-Zeile gehört dem umsetzenden Zug), gepinntes
      Build-Werkzeug samt Prüfsumme, `make`-Ziel `examples-kotlin` in
      `Makefile`/`harness/mk/*.mk` — der Bau kompiliert netzlos prüfbar (kein
      laufender Dienst nötig).
- [x] **LP2 — der HTTP-Client.** `examples/kotlin/http-client/` ruft real
      `GET /tables` mit dem `reader`-Token auf, gibt die Antwort aus, liest
      Adresse/Token aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER` mit
      Flag-/Umgebungs-Übersteuerung (Form-Vorbild: `examples/http-client` in
      Go); seine netzlos prüfbaren Teile sind getestet und laufen über
      `examples-kotlin`.
- [x] **LP3 — die Träger samt Workflow.** `docs/user/benutzerhandbuch.md`
      §4 „Zugriff über die HTTP-/JSON-API": der `**Beispiele:**`-Block (seit
      `slice-098` mit Go- und C#-Zeile) bekommt die **dritte** Zeile (Kotlin)
      samt Änderungshistorie-Zeile; `harness/README.md` §Werkzeuge trägt
      `examples-kotlin`; der nicht-blockierende Workflow aus `slice-098`
      (`.github/workflows/examples.yml`) bekommt einen **zweiten Job** für
      `examples-kotlin` — der Carrier wächst mit seinem Umfang
      ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
      Festlegung 5).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/review-slice-099.md`
      liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-099.md`
      liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) — *entfällt: Greenfield-
      Bootstrap (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF),
      keine Datei vorhanden.*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im
      Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/kotlin/Dockerfile` | neu | Eigener, digest-gepinnter Bau-Kontext = Sprach-Wurzel (`ADR-0087` Festlegung 3, bestätigt); zwei Stufen (`build`: `eclipse-temurin:21-jdk`, `runtime`: `eclipse-temurin:21-jre`), Digests real gemessen (`docker manifest inspect`, amd64/linux). |
| `examples/kotlin/{gradlew,gradlew.bat,gradle/wrapper/*}` | neu | Gradle-Wrapper 8.14 (erzeugt über den gepinnten `gradle:8.14-jdk21`-Image-Digest), `distributionSha256Sum` in `gradle-wrapper.properties` ergänzt — Pin-Hebung ist ein bewusster Commit (`SPEC-023`). |
| `examples/kotlin/{settings.gradle.kts,build.gradle.kts}` | neu | Wurzelprojekt (ein Modul `http-client`), Kotlin-Gradle-Plugin-Version 2.4.20 einmal gepinnt. |
| `examples/kotlin/http-client/build.gradle.kts` | neu | Programm-Manifest — kein Fremdmodul für den Client selbst (`java.net.http`), gepinnte Testabhängigkeiten (`kotlin-test-junit5:2.4.20`, `junit-platform-launcher:6.1.3`). |
| `examples/kotlin/http-client/src/**` | neu | Der Anfrage/Antwort-Client (`Main.kt`, `Config.kt`, `Cli.kt`, `TablesClient.kt`, `TablesUrlBuilder.kt`) und seine netzlos prüfbaren Tests, Form-Vorbild `examples/http-client` (Go) und `examples/csharp/http-client` (C#). |
| `Makefile` / `harness/mk/examples.mk` | update | Neues Ziel `examples-kotlin` (Werkzeug, kein Gate) neben dem bestehenden `examples-csharp`. |
| `.github/workflows/examples.yml` | update | Zweiter Job `examples-kotlin`, neben dem in `slice-098` angelegten C#-Job (`examples-csharp` umbenannt für die Job-Symmetrie, kein Verhaltenswechsel). |
| `docs/user/benutzerhandbuch.md` | update | §4 „Zugriff über die HTTP-/JSON-API": dritte Zeile (Kotlin) im `**Beispiele:**`-Block; Änderungshistorie-Zeile 1.21. |
| `harness/README.md` §Werkzeuge | update | Zeile für `examples-kotlin`, kein Gate, nicht in `GATE_CHECKS`. |
| `spec/pflichtenheft.md` `SPEC-023` „Sprachen und Umfang“ | update (Plan-Nachzug) | Nicht im ursprünglichen Plan-Umfang, aber als offene, adresslose `ADR-0090`-Folgepflicht seit `slice-098` gemessen (`BEO-PGC/adr-folgepflicht-ohne-traeger-slice`, 1×, offen) — in diesem Zug auf die volle Matrix nachgezogen, samt §7-Historie-Zeile; Entscheidung im Bericht begründet. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
ist `Accepted`, und `examples/kotlin/` existiert nicht (gemessen: `ls
examples/kotlin` → kein Treffer). Keine Abhängigkeit auf `slice-098` oder
einen anderen Slice dieser Matrix-Erweiterung — parallel lieferbar. Ohne
Rückfrage feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die
  Kotlin-Werkzeugkette mehr als **ein** Projekt oder eine zweite Bau-Stufe
  braucht, um netzlos zu testen — dann trägt die Sprach-Wurzel-Form allein
  schon mehr als drei Liefer-Punkte.
- `in-progress` → `open` (blockiert — Carveout?): wenn die digest-gepinnte
  JVM-Basis nicht mehr auflösbar ist (Registry-Digest verschwunden) oder die
  Werkzeugkette eine nicht-öffentliche Quelle verlangt
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). Dann ist die Sprach-Form neu zu bewerten, kein
  stiller Workaround.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — `examples-kotlin`
kompiliert und testet den Client netzlos, der nicht-blockierende Workflow trägt
den zweiten Job und läuft real (mindestens ein Post-Push-Lauf sichtbar,
`AGENTS.md` §3.10), beide Träger nennen die Kotlin-Zeile — **und**
`make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: ob die zwei Sprach-Wurzel-Formen (C#, Kotlin)
nach zwei realen Bauten dasselbe Muster tragen oder eine dritte,
sprachneutrale Form nahelegen — ob daraus mehr wird als eine Wiederholung,
entscheidet der Lauf.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die digest-gepinnte JVM-Basis ist zum Zeitpunkt des Baus nicht mehr
  auflösbar** (Tag zurückgezogen, Registry-Digest verschoben). —
  **Ausgang:** <…>
- **Die Kotlin-Werkzeugkette verlangt eine nicht-öffentliche Quelle**
  (privates Maven-Repository, Zugangsdaten in der Auflösung) —
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3. — **Ausgang:** <…>
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (zwei
  Jobs statt einem, Laufzeit über dem Runner-Budget, oder ein dauerhaft
  roter, überlesener Lauf) — §Re-Evaluierungs-Trigger 4,
  `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×, offen),
  `BEO-PGC/github-actions-unverifizierbar-lokal` (5×, verkörpert in
  `AGENTS.md` §3.10). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-
  oberflaeche`, `handbuch-versionshistorie-uebersprungen`); sie ist der
  einzige Teil dieses Slice, den **kein** Kompilat erzwingt. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht; besonders der Workflow aus `slice-098`, den
  dieser Slice erweitert). — **Ausgang:** <…>

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `examples/kotlin/**`,
`Makefile`/`harness/mk/*.mk`, `.github/workflows/**` und
`docs/user/benutzerhandbuch.md` — die repo-weite Default-Sub-Area `*`/`PGC`
aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine
zu grobe Sub-Area zu differenzieren; es gibt in diesem Repo keine eigene
Kotlin-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Vier
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide als LP3 in die DoD gezogen; `nicht-blockierender-
workflow-alarmmuedigkeit` (**1×**, offen) und `github-actions-
unverifizierbar-lokal` (**5×**, verkörpert in `AGENTS.md` §3.10) — beide
treffen den erweiterten Workflow, als Risiko in §6. Kein Treffer zu
Kotlin/JVM/Gradle/Maven selbst (gemessen: `grep -rli
"kotlin\|maven" docs/plan/planning/observations/` → kein Fund).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**.
Der Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

### Sub-Area: `*` (Default, PGC)

- **Modus:** GF
- **Konventionen-Dichte:** `harness/conventions.md` Modus-Deklaration setzt
  GF für das gesamte Repo (Doc führt, Code folgt) — die drei Spec-Straten
  standen vor dem ersten Code-Commit.
- **Phase-Reife:** Phase 5 (etabliert) — das Muster „Sprach-Wurzel + eigenes
  Dockerfile + `make`-Ziel + nicht-blockierender Workflow" ist mit
  [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) bereits
  entschieden; dieser Slice überträgt es auf die zweite Sprache.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; kein Bestand, der
  von der neuen Form abweichen könnte.
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
