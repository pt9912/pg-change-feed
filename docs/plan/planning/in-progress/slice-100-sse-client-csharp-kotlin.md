# Slice slice-100: SSE-Client in C# und Kotlin

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). Dieses Repo führt derzeit **keine** offene Welle
(`docs/plan/planning/in-progress/roadmap.md` §Offene Wellen ist leer).

**Bezug:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(Festlegung 1 — Zellen `csharp`×`sse` und `kotlin`×`sse`, „**keine neue
Abhängigkeit** — beide Runtimes tragen Streaming-HTTP in ihrer
Standardausstattung"; §Slice-Schnitt-Empfehlung, Zeile 3: „**SSE-Client in C#
und Kotlin** … Abhängigkeit 1, 2") · [`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
(der SSE-Endpunkt `GET /changes/stream`, dessen Form der Client anspricht) ·
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(Form-Vorbild `examples/sse-client` in Go) · `BEO-PGC/handbuch-nicht-
nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**, verkörpert) und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (**3×**, verkörpert) —
beide treffen den Handbuch-Nachzug dieses Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-008` (der Live-Change-Stream, hier über
SSE) — dieser Slice **zeigt** ihn in zwei weiteren Sprachen, er ändert ihn
nicht.

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

**Ziel:** `examples/csharp/sse-client/` und `examples/kotlin/sse-client/`
entstehen — beide öffnen `GET /changes/stream` mit dem `reader`-Token und
geben jedes Event aus (Form-Vorbild `examples/sse-client` in Go), mit
netzlos prüfbaren Tests (Zerlegen eines SSE-Frames). **Keine neue
Abhängigkeit:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
Festlegung 1 hält fest, dass beide Runtimes Streaming-HTTP in ihrer
Standardausstattung tragen (C#: der `HttpClient` der Runtime samt dem
SSE-Parser `System.Net.ServerSentEvents`; Kotlin/JVM: `java.net.http.HttpClient`
mit einem streamenden Body-Handler seit JDK 11) — welcher der beiden Wege je
Sprache gewählt wird, ist Detail dieses Zuges. Beide Handbuch-Zeilen (der
zweite Zugriffs-Abschnitt, der einen `**Beispiele:**`-Block bekommt) werden
nachgezogen.

**Warum beide Sprachen in einem Slice:** Beide Programme teilen dieselbe
Handbuch-Form (`**Beispiele:**`-Block, ein Zugriffs-Abschnitt) und dieselbe
Voraussetzung (beide Sprach-Wurzeln existieren) — kein Codegenerator, keine
neue Abhängigkeit, kein neuer Bau-Kontext auf beiden Seiten. Zwei Programme,
zwei Test-Suiten, zwei Handbuch-Zeilen bleiben bei drei Liefer-Punkten.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der NATS-Client in C#/Kotlin.** Folge-Slice `slice-101` übernimmt ihn —
  eigene Abhängigkeit (gepinnte Client-Bibliothek je Sprache), eigener
  Zugriffs-Abschnitt im Handbuch.
- **Der gRPC-Client in C#/Kotlin.** Folge-Slices `slice-102` (C#) und
  `slice-103` (Kotlin) übernehmen ihn — er trägt den benannten Zusatzkontext
  und den Protobuf-/gRPC-Generator, zwei Form-Fragen, die dieser Slice nicht
  öffnet.
- **Der Go-SSE-Client (`examples/sse-client`).** Er existiert bereits
  (`slice-095`, `in-progress/`) und ist das **Form-Vorbild**; ihn umzubauen
  wäre Arbeit an Bestand ohne Adresse.
- **Eine dritte, gemeinsame SSE-Bibliothek für beide Sprachen.** Es gibt in
  diesem Repo keine sprachneutrale Werkzeugketten-Form
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 1); jede Sprache trägt ihre eigene
  Standardbibliotheks-Lösung. Anderer Vorgang, träte er ein.
- **Eine Änderung am SSE-Endpunkt oder seinem Nachrichtenschema.** Der Client
  **benutzt** `GET /changes/stream`; verlangt er eine Vertragsänderung, ist
  das eine Spec-Änderung, kein Beispiel-Umbau (`SPEC-023`).
- **`.a-check.yml`.** Unverändert aus denselben Gründen wie in `slice-098`/
  `slice-099`: keine C#-/Kotlin-Schicht in diesem Repo, kein neuer Schlüssel.
  Bestand bleibt bewusst stehen.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1 — der C#-SSE-Client.** `examples/csharp/sse-client/` öffnet real
      `GET /changes/stream` mit dem `reader`-Token, gibt jedes Event aus,
      liest Adresse/Token aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`; seine
      netzlos prüfbaren Teile (Zerlegen eines SSE-Frames) sind getestet und
      laufen über `examples-csharp`.
- [x] **LP2 — der Kotlin-SSE-Client.** `examples/kotlin/sse-client/` — dieselbe
      Zusage, über `examples-kotlin`.
- [x] **LP3 — die zwei Handbuch-Zeilen.** `docs/user/benutzerhandbuch.md` §4
      „Zugriff über Server-Sent-Events": der bestehende `**Beispiele:**`-
      Block (Go seit `slice-095`) bekommt zwei weitere Zeilen (C#, Kotlin),
      samt Änderungshistorie-Zeile.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-100.md`
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
| `examples/csharp/sse-client/**` | neu | SSE-Client, Standardbibliothek (`System.Net.ServerSentEvents` oder `HttpClient`-Streaming), Form-Vorbild `examples/sse-client` (Go). |
| `examples/kotlin/sse-client/**` | neu | SSE-Client, Standardbibliothek (`java.net.http.HttpClient` mit streamendem Body-Handler). |
| `docs/user/benutzerhandbuch.md` | update | §4 „Zugriff über Server-Sent-Events": zwei weitere Zeilen (C#, Kotlin) im `**Beispiele:**`-Block; Änderungshistorie-Zeile. |
| `examples/csharp/Dockerfile` (Plan-Nachzug) | update | Ohne Fremdmodul (Entscheidung dieses Zuges, §5) hätte der SSE-Client keinen eigenen Bau-Zusatz gebraucht — die zweite Zelle (`sse-client`) läuft aber im selben Container-Kontext wie `http-client`, deshalb erweitert die gemeinsame `build`-Stufe um Restore/Build/Test/Publish des zweiten Programms; zwei neue `runtime-*`-Stufen (`runtime-sse` vor der unverändert letzten `runtime`) halten den bestehenden `pg-change-feed-examples:csharp`-Tag byte-gleich. |
| `examples/kotlin/Dockerfile` (Plan-Nachzug) | update | Dasselbe Muster: gemeinsame `build`-Stufe um `:sse-client:test`/`:sse-client:installDist` erweitert, neue `runtime-sse`-Stufe vor der unverändert letzten `runtime`-Stufe. |
| `examples/kotlin/settings.gradle.kts` (Plan-Nachzug) | update | `include("sse-client")` — Multi-Modul-Registrierung des zweiten Gradle-Moduls, nicht im ursprünglichen Plan einzeln benannt. |
| `harness/mk/examples.mk` (Plan-Nachzug) | update | Beide Ziele bauen jetzt **zwei** Runtime-Images je Sprache (zweiter `docker build --target runtime-sse`-Aufruf) statt eines. |
| `harness/README.md` §Sensors (Plan-Nachzug, §3.13-Fund) | update | Der §3.13-Suchlauf fand zwei Sätze, die dieser Zug falsch macht: beide Zeilen beschrieben „das Werkzeugketten-Image" (Singular) je Sprache — jetzt zwei Images je Sprache, Zeilen korrigiert und Image-Tags benannt. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-098` **und** `slice-099` liegen in
`done/` — die C#- und Kotlin-Sprachwurzel samt Werkzeugkette und `make`-Ziel
existieren, ohne sie ist kein Bau-Kontext für den SSE-Client vorhanden. Ohne
Rückfrage feststellbar (Verzeichnis-Position der beiden Vorgänger-Slices).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn eine der
  beiden Runtimes entgegen der Annahme in
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 1
  doch eine externe Abhängigkeit für SSE-Streaming braucht — dann ist die
  Größenannahme („keine neue Abhängigkeit") falsch und der Slice ist neu zu
  bewerten (ggf. Trennung in zwei Slices je Sprache).
- `in-progress` → `open` (blockiert — Carveout?): wenn `slice-098` oder
  `slice-099` zum Zeitpunkt des Starts entgegen der Prüfung doch nicht in
  `done/` liegen (WIP-Limit-Verstoß) oder eine der beiden Sprach-Wurzeln
  einen inkompatiblen Bau-Zustand hinterlassen hat.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — beide Clients
öffnen real den SSE-Stream und geben Events aus, beide sind über ihr
jeweiliges Sprachziel kompiliert und getestet, die Handbuch-Zeilen tragen
beide Sprachen — **und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: ob die Annahme „keine neue Abhängigkeit" für
beide Runtimes real trägt, oder ob eine der beiden doch ein Fremdmodul
braucht (`ADR-0090` Festlegung 4 nennt `System.Net.ServerSentEvents`/
`com.squareup.okhttp3:okhttp-sse` als **Kandidaten**, nicht als
Festlegung) — der Lauf entscheidet, welcher Weg gewählt wurde und ob er trägt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Eine der beiden Runtimes braucht doch ein Fremdmodul für SSE** (z. B.
  `System.Net.ServerSentEvents` bei einer älteren .NET-Basis als angenommen,
  oder `okhttp-sse` statt `java.net.http`) — Registry-/Digest-Pin dieses
  Fremdmoduls nicht mehr auflösbar wäre die Folge-Ausprägung. —
  **Ausgang:** <…>
- **Die Werkzeugkette einer Sprache verlangt eine nicht-öffentliche Quelle**
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). — **Ausgang:** <…>
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (zwei
  weitere Bau-Ziele im selben Workflow, Laufzeit über dem Runner-Budget) —
  §Re-Evaluierungs-Trigger 4, `BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit` (1×, offen), `BEO-PGC/github-actions-unverifizierbar-lokal`
  (5×, verkörpert in `AGENTS.md` §3.10). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen; sie ist der einzige Teil dieses Slice, den **kein**
  Kompilat erzwingt. — **Ausgang:** <…>
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
`examples/kotlin/**` und `docs/user/benutzerhandbuch.md` — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine zu
grobe Sub-Area zu differenzieren.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Vier
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide als LP3 in die DoD gezogen; `nicht-blockierender-
workflow-alarmmuedigkeit` (**1×**, offen) und `github-actions-
unverifizierbar-lokal` (**5×**, verkörpert in `AGENTS.md` §3.10) — beide
treffen den weiter wachsenden Workflow, als Risiko in §6.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**.
Der Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

### Sub-Area: `*` (Default, PGC)

- **Modus:** GF
- **Konventionen-Dichte:** `harness/conventions.md` Modus-Deklaration setzt
  GF für das gesamte Repo (Doc führt, Code folgt).
- **Phase-Reife:** Phase 5 (etabliert) — beide Sprach-Wurzeln existieren
  bereits (`slice-098`, `slice-099`); dieser Slice fügt ein zweites Programm
  je Sprache ohne neue Bauform hinzu.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; kein Bestand, der
  von der neuen Form abweichen könnte.
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
