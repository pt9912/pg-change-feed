# Slice slice-098: C#-Sprach-Wurzel und HTTP-Client — die erste Zelle der Matrix

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). Dieses Repo führt derzeit **keine** offene Welle
(`docs/plan/planning/in-progress/roadmap.md` §Offene Wellen ist leer).

**Bezug:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(Festlegung 1 — Zelle `csharp`×`http` der Matrix; Festlegung 5 — je Sprache
ein `make`-Ziel, kein Gate; Festlegung 4 — Abhängigkeits-Pins; §Slice-Schnitt-
Empfehlung, Zeile 1: „**C#-Sprach-Wurzel + HTTP-Client** … keine
Abhängigkeit") · [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)
Festlegung 2/3 — **bestätigt, nicht superseded für diesen Teil**: der
Bau-Kontext ist das Sprach-Wurzelverzeichnis, ein eigenes digest-gepinntes
Dockerfile, kein Gate — dieser Slice ist der erste reale Bau, der diese Form
einlöst (derselbe ADR-Text selbst: eine Sprache, die zum ersten Mal auftritt,
**steht allein**, sie beweist die Form) ·
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(der nicht-blockierende Workflow als Träger-Muster, hier auf eine zweite
Sprache übertragen) · `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(**3×**, verkörpert) und `BEO-PGC/handbuch-versionshistorie-uebersprungen`
(**3×**, verkörpert) — beide treffen den Handbuch-Nachzug dieses Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-006` (die HTTP-/JSON-API) — dieser
Slice **zeigt** sie in einer weiteren Sprache, er ändert sie nicht.

**Verantwortlich:** — *(bis zur Priorisierung `open` → `next`)*.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `examples/csharp/` entsteht als eigenes Sprach-Wurzelverzeichnis mit
eigenem, digest-gepinntem Dockerfile, gepinntem Projekt-/Paket-Manifest und
dem `make`-Ziel `examples-csharp`; darin liegt der erste C#-Client,
`examples/csharp/http-client/` — ein Anfrage/Antwort-Aufruf gegen die
Verwaltungs-API (`GET /tables` mit dem `reader`-Token, Form-Vorbild
`examples/http-client` in Go), mit netzlos prüfbaren Tests. Ein
nicht-blockierender GitHub-Actions-Workflow fährt das Ziel, und Handbuch samt
`harness/README.md` §Werkzeuge tragen die neue Zeile.

**Warum dieser Slice zuerst und allein:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
§Slice-Schnitt-Empfehlung schneidet die beiden Sprach-Wurzeln (C#, Kotlin) vor
jede weitere Oberfläche — jede Sprache **steht allein** und beweist ihre Form,
bevor SSE-, NATS- oder gRPC-Clients auf ihr aufbauen
([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) Festlegung 3,
derselbe Satz für den ersten Zug einer Sprache).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der SSE-Client in C#.** Folge-Slice `slice-100` übernimmt ihn — zusammen
  mit dem Kotlin-Pendant, weil beide dieselbe Handbuch-Form (ein
  `**Beispiele:**`-Block) und dieselbe Abhängigkeit auf beide Sprach-Wurzeln
  tragen.
- **Der NATS-Client in C#.** Folge-Slice `slice-101` übernimmt ihn, samt dem
  gepinnten `NATS.Net`-Paket — aus demselben Grund wie oben.
- **Der gRPC-Client in C#, der benannte Zusatzkontext und der
  Protobuf-/gRPC-Generator.** Folge-Slice `slice-102` übernimmt sie. Die zwei
  Form-Fragen dieser Zelle — **Bau-Kontext, Generator** — stehen **nicht**
  hier: [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) weist sie
  ausdrücklich dem jeweiligen gRPC-Slice zu (§Slice-Schnitt-Empfehlung, „Was
  der Planner beim Anlegen beachten muss").
- **Die Kotlin-Sprach-Wurzel und ihr HTTP-Client.** Sie laufen **parallel** in
  `slice-099` — eigene Werkzeugkette, eigenes Dockerfile, keine Abhängigkeit
  in beide Richtungen (§Slice-Schnitt-Empfehlung, Zeile 2: „keine (parallel
  zu 1)"). Ein gemeinsamer Slice für beide Sprachen würde zwei
  Werkzeugketten in einer Sitzung prüfbar machen müssen — über der
  Drei-Liefer-Punkte-Grenze.
- **Der Go-gRPC-Client (`examples/grpc-client/`).** Er ist **kein** Teil
  dieser Matrix-Erweiterung um C#/Kotlin — er läuft als eigener Liefer-Punkt
  (LP2) von `slice-095` (aktuell `in-progress/`), gebunden an den Umzug der
  Vertragsfläche (`slice-097`, `done/`). Anderer Vorgang, andere Sprache,
  anderer Träger.
- **`.a-check.yml`.** Gemessen (`grep -n "languages" .a-check.yml`): die
  Gruppe `examples` trägt ausschließlich Go-Semantik, kein `csharp`-Schlüssel
  existiert und keiner entsteht hier — [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 6 hält das für die **gesamte** Matrix fest: es gibt in diesem
  Repo keine C#-Schicht, gegen die ein Import verstoßen könnte. Bestand
  bleibt bewusst stehen.
- **Ein eigenes Gate für das Sprachziel.** `examples-csharp` bleibt ein
  **Werkzeug** ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 5) — es braucht Netz (Paket-Bezug) und kann deshalb kein
  netzloses Gate werden; ein Gate wäre eine Zusage, die der Bau nicht halten
  kann.
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

- [ ] **LP1 — Sprach-Wurzel, Werkzeugkette, Ziel.** `examples/csharp/Dockerfile`
      mit digest-gepinnter .NET-SDK-Basis (Kandidat
      `mcr.microsoft.com/dotnet/sdk:10.0`, Existenz gemessen in
      [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung
      4 — Digest-Zeile gehört dem umsetzenden Zug), gepinntes
      Projekt-/Paket-Manifest, `make`-Ziel `examples-csharp` in
      `Makefile`/`harness/mk/*.mk` — der Bau kompiliert netzlos prüfbar (kein
      laufender Dienst nötig).
- [ ] **LP2 — der HTTP-Client.** `examples/csharp/http-client/` ruft real
      `GET /tables` mit dem `reader`-Token auf, gibt die Antwort aus, liest
      Adresse/Token aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER` mit
      Flag-Übersteuerung (Form-Vorbild: `examples/http-client` in Go); seine
      netzlos prüfbaren Teile (Aufbau der Anfrage, Fehlerpfad ohne
      erreichbaren Host) sind getestet und laufen über `examples-csharp`.
- [ ] **LP3 — die Träger samt Workflow.** `docs/user/benutzerhandbuch.md`
      §4 „Zugriff über die HTTP-/JSON-API": der bestehende `**Beispiel:**`-
      Absatz (nur Go) wird zu einem `**Beispiele:**`-Block mit **einer Zeile
      je Sprache** (Go, C#) — die im ADR festgelegte Ziel-Form
      ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
      Festlegung 7) — samt Änderungshistorie-Zeile; `harness/README.md`
      §Werkzeuge trägt `examples-csharp`; ein neuer, **nicht-blockierender**
      GitHub-Actions-Workflow (Vorschlag: `.github/workflows/examples.yml`)
      fährt das Ziel auf jeden PR/Push, ohne Required-Status-Check.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-098.md`
      liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) — *entfällt: Greenfield-
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
| `examples/csharp/Dockerfile` | neu | Eigener, digest-gepinnter Bau-Kontext = Sprach-Wurzel (`ADR-0087` Festlegung 3, bestätigt). |
| `examples/csharp/<Paket-Manifest>` (Name/Form gehört dem umsetzenden Zug, z. B. `Directory.Packages.props`) | neu | Gepinntes Projekt-/Paket-Manifest — Pin-Hebung ist ein bewusster Commit (`SPEC-023`). |
| `examples/csharp/http-client/**` | neu | Der Anfrage/Antwort-Client, Form-Vorbild `examples/http-client` (Go). |
| `Makefile` / `harness/mk/*.mk` | update | Neues Ziel `examples-csharp` (Werkzeug, kein Gate). |
| `.github/workflows/examples.yml` (Name Vorschlag) | neu | Nicht-blockierender Workflow für `examples-csharp` — kein Required-Status-Check. |
| `docs/user/benutzerhandbuch.md` | update | §4 „Zugriff über die HTTP-/JSON-API": `**Beispiel:**` → `**Beispiele:**`-Block mit Go- und C#-Zeile; Änderungshistorie-Zeile. |
| `harness/README.md` §Werkzeuge | update | Zeile für `examples-csharp`, kein Gate, nicht in `GATE_CHECKS`. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
ist `Accepted`, und `examples/csharp/` existiert nicht (gemessen: `ls
examples/csharp` → kein Treffer). Keine Abhängigkeit auf einen anderen Slice
dieser Matrix-Erweiterung. Ohne Rückfrage feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die
  C#-Werkzeugkette mehr als **ein** Projekt oder eine zweite Bau-Stufe
  braucht, um netzlos zu testen — dann trägt die Sprach-Wurzel-Form allein
  schon mehr als drei Liefer-Punkte.
- `in-progress` → `open` (blockiert — Carveout?): wenn die digest-gepinnte
  .NET-SDK-Basis nicht mehr auflösbar ist (Registry-Digest verschwunden) oder
  die Werkzeugkette eine nicht-öffentliche Quelle verlangt
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). Dann ist die Sprach-Form neu zu bewerten, kein
  stiller Workaround.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — `examples-csharp`
kompiliert und testet den Client netzlos, der nicht-blockierende Workflow ist
angelegt und läuft real (mindestens ein Post-Push-Lauf sichtbar, `AGENTS.md`
§3.10), beide Träger nennen die C#-Zeile — **und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: die erste reale Bauprobe der
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)-Sprachform
(bislang nur gemessen an C#-Basis-Existenz, nicht an einem echten Bau) — ob
sie unverändert trägt, entscheidet der Lauf.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die digest-gepinnte .NET-SDK-Basis ist zum Zeitpunkt des Baus nicht mehr
  auflösbar** (Tag zurückgezogen, Registry-Digest verschoben). —
  **Ausgang:** <…>
- **Die C#-Werkzeugkette verlangt eine nicht-öffentliche Quelle** (privater
  NuGet-Feed, Zugangsdaten in der Auflösung) —
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3. — **Ausgang:** <…>
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (Laufzeit
  über dem Runner-Budget, oder ein dauerhaft roter, überlesener Lauf) —
  §Re-Evaluierungs-Trigger 4, `BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit` (1×, offen), `BEO-PGC/github-actions-unverifizierbar-lokal`
  (5×, verkörpert in `AGENTS.md` §3.10). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-
  oberflaeche`, `handbuch-versionshistorie-uebersprungen`); sie ist der
  einzige Teil dieses Slice, den **kein** Kompilat erzwingt. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht). — **Ausgang:** <…>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `examples/csharp/**`,
`Makefile`/`harness/mk/*.mk`, `.github/workflows/**` und
`docs/user/benutzerhandbuch.md` — die repo-weite Default-Sub-Area `*`/`PGC`
aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine
zu grobe Sub-Area zu differenzieren; es gibt in diesem Repo keine eigene
C#-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Vier
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide als LP3 in die DoD gezogen; `nicht-blockierender-
workflow-alarmmuedigkeit` (**1×**, offen) und `github-actions-
unverifizierbar-lokal` (**5×**, verkörpert in `AGENTS.md` §3.10) — beide
treffen den neuen Workflow, als Risiko in §6. Kein Treffer zu C#/.NET/NuGet
selbst (gemessen: `grep -rli "dotnet\|csharp\|nuget" docs/plan/planning/observations/`
→ kein Fund).

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
  entschieden, dieser Slice liefert die erste reale Instanz.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; kein Bestand, der
  von der neuen Form abweichen könnte.
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
