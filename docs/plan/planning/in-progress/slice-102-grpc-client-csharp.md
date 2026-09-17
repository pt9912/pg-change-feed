# Slice slice-102: gRPC-Client in C# — der benannte Zusatzkontext, real gebaut

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Reaktives auf Nutzerentscheidung, ein einzelner Slice
ohne Closure-Bedingung jenseits seiner DoD (Modul 6 §Wann Arbeit eine Welle
braucht). Dieses Repo führt derzeit **keine** offene Welle
(`docs/plan/planning/in-progress/roadmap.md` §Offene Wellen ist leer).

**Bezug:** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
(Festlegung 1 — Zelle `csharp`×`grpc`, „der benannte Zusatzkontext **und** der
Protobuf-/gRPC-Generator der Sprache — gemessen existent"; **Festlegung 2** —
der `.proto`-Weg: der Bau-Kontext bleibt das Sprach-Wurzelverzeichnis, ein
**zusätzlicher, benannter** Kontext trägt `proto/`, real gemessen per
`docker buildx build --build-context …` mit `COPY --from=<name> …`;
**Festlegung 3** — der erzeugte Stub entsteht **im Bau**, wird **nicht**
committet (`gen/**` bleibt Go-Bindung); Festlegung 4 — `Grpc.Net.Client`,
`Grpc.Tools`, `Google.Protobuf` (NuGet, Registry-Existenz gemessen
2026-09-17); §Slice-Schnitt-Empfehlung, Zeile 5: „**gRPC-Client in C#** …
Abhängigkeit 1") · [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(der Server-Stream und die `.proto`, deren Form der Client anspricht) ·
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(Form-Vorbild `examples/grpc-client` in Go, Liefer-Punkt LP2 von `slice-095`)
· `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(**3×**, verkörpert) und `BEO-PGC/handbuch-versionshistorie-uebersprungen`
(**3×**, verkörpert) — beide treffen den Handbuch-Nachzug dieses Slice.

**Berührte Spec-Stellen:** `LH-FA-SST-008` (der Live-Change-Stream, hier über
gRPC) — dieser Slice **zeigt** ihn in einer weiteren Sprache, er ändert ihn
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

**Ziel:** `examples/csharp/grpc-client/` öffnet real
`ChangeStream/StreamChanges` und gibt jede Nachricht aus (Form-Vorbild
`examples/grpc-client` in Go). Anders als die bisherigen C#-Zellen bringt
diese Zelle **zwei neue Form-Fragen** mit, beide in diesem Slice zu
beantworten — genau die zwei, die
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) dem
gRPC-Slice zuweist, nicht den HTTP-/SSE-/NATS-Slices:

1. **Der benannte Zusatzkontext.** `examples/csharp/Dockerfile` bekommt einen
   zusätzlichen, benannten Bau-Kontext, der `proto/` trägt
   (`docker buildx build --build-context <name>=<proto-Verzeichnis> …`),
   und kopiert die `.proto` daraus (`COPY --from=<name> …`) — real gemessen
   in [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) §Kontext
   (`buildx` v0.37.1, die Datei kommt an); der Bau-Kontext selbst bleibt das
   Sprach-Wurzelverzeichnis, unverändert.
2. **Der Generator.** `Grpc.Tools` (trägt `protoc` und das C#-Plugin) erzeugt
   den Stub **im Bau** aus der kopierten `.proto`; `Grpc.Net.Client` und
   `Google.Protobuf` tragen die Laufzeit. Der Stub liegt **nicht** im Baum
   (`gen/**` bleibt die Go-Bindung, [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
   Festlegung 3).

Der Handbuch-Abschnitt „Zugriff über den gRPC-Change-Stream" (bislang nur
Go, `slice-095`) bekommt die C#-Zeile im `**Beispiele:**`-Block.

**Warum dieser Slice vor `slice-103` (Kotlin-gRPC):** Die Form (Zusatzkontext
+ Generator) wird hier **zum ersten Mal real gebaut und gemessen** — bislang
nur mit einer isolierten `buildx`-Probe belegt
([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) §Kontext), noch
nicht im offiziellen Sprach-Bau. `slice-103` **überträgt** die hier gemessene
Form auf Kotlin, statt sie zweimal unabhängig zu erfinden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der gRPC-Client in Kotlin.** Folge-Slice `slice-103` übernimmt ihn — er
  überträgt die hier gemessene Form (Zusatzkontext, Generator-Stufe) auf die
  Kotlin-Werkzeugkette, statt sie in derselben Sitzung zweimal zu bauen und
  zu prüfen.
- **Der Go-gRPC-Client (`examples/grpc-client/`).** Er ist **kein** Teil
  dieser Matrix-Erweiterung um C#/Kotlin — er läuft als eigener Liefer-Punkt
  (LP2) von `slice-095` (aktuell `in-progress/`), gebunden an den Umzug der
  Vertragsfläche (`slice-097`, `done/`, nicht an diesen Slice). Anderer
  Vorgang, andere Sprache, anderer Träger.
- **Der Inhalt der `.proto`.** Der Draht-Vertrag ist
  [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)s Gegenstand;
  dieser Slice **liest** die `.proto` über den Zusatzkontext, ändert an ihr
  keine Nachricht, kein Feld — nur ein lesender Bauweg kommt hinzu
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 3,
  „dreimal gelesen, einmal geschrieben").
- **Ein committeter C#-Stub.** [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 3 entscheidet ausdrücklich gegen ein committetes Erzeugnis (kein
  drittes Sync-Gate, keine zweite Repräsentation des Vertrags). Bestand
  bleibt bewusst stehen — es gibt keinen Stub, den dieser Slice committen
  müsste.
- **`.a-check.yml`.** Unverändert — [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 6 hält fest, dass die `contract`-Gruppe (`gen/**`) **Go** bleibt;
  der erzeugte C#-Stub liegt nicht im Baum und wird von a-check ohnehin nicht
  gelesen. Bestand bleibt bewusst stehen.
- **Eine Änderung an `Dockerfile` oder `.dockerignore` der Wurzel.** Der
  ausgelieferte Bau bleibt unberührt
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) §8); der
  Zusatzkontext betrifft ausschließlich `examples/csharp/Dockerfile`.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **LP1 — der `.proto`-Weg im Bau.** `examples/csharp/Dockerfile` liest
      die `.proto` über einen zusätzlichen, benannten Bau-Kontext
      (`docker buildx build --build-context …`); `make examples-csharp` ruft
      den Bau mit diesem Zusatzkontext auf. Beleg: der Bau erzeugt den Stub
      real (kein Fallback, kein übersprungener Schritt bei fehlendem
      Kontext — der Bau bricht ohne ihn ab, wie
      [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
      §Fitness Function festhält).
- [ ] **LP2 — der Client.** `examples/csharp/grpc-client/` öffnet real den
      Server-Stream `ChangeStream/StreamChanges`, gibt jede Nachricht aus,
      liest Adresse/Token aus `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER`; netzlos
      prüfbare Teile sind getestet und laufen über `examples-csharp`.
- [ ] **LP3 — die Handbuch-Zeile.** `docs/user/benutzerhandbuch.md` §4
      „Zugriff über den gRPC-Change-Stream": der bestehende `**Beispiele:**`-
      Block (Go seit `slice-095`) bekommt die C#-Zeile, samt
      Änderungshistorie-Zeile.
- [ ] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/review-slice-102.md`
      liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-102.md`
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
| `examples/csharp/Dockerfile` | update | Generator-Stufe (`Grpc.Tools`) und der zusätzliche, benannte Bau-Kontext für `proto/`. |
| `examples/csharp/<Paket-Manifest>` | update | Pin von `Grpc.Net.Client`, `Grpc.Tools`, `Google.Protobuf`. |
| `examples/csharp/grpc-client/**` | neu | gRPC-Server-Stream-Client, Form-Vorbild `examples/grpc-client` (Go, `slice-095`). |
| `Makefile` / `harness/mk/*.mk` | update | `examples-csharp`-Ziel ruft den Bau mit dem benannten Zusatzkontext auf. |
| `docs/user/benutzerhandbuch.md` | update | §4 „Zugriff über den gRPC-Change-Stream": C#-Zeile im `**Beispiele:**`-Block; Änderungshistorie-Zeile. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-098` liegt in `done/` — die
C#-Sprachwurzel samt Werkzeugkette und `make`-Ziel existiert; die `.proto`
liegt unverändert unter `proto/cdc/stream/v1/`. **Keine** Abhängigkeit auf
`slice-099` oder `slice-101` (unabhängig von der Kotlin-Spur und vom
NATS-Client). Ohne Rückfrage feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn der benannte
  Zusatzkontext im realen Sprach-Bau **nicht** trägt (anders als die isolierte
  Probe in [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Kontext) und eine andere `.proto`-Zugangsform nötig wird — dann ist die
  Bau-Form neu zu bewerten, nicht in diesem Slice stillschweigend zu ändern.
- `in-progress` → `open` (blockiert — Carveout?): wenn `Grpc.Tools` den Stub
  nicht reproduzierbar aus der `.proto` erzeugt (nicht-deterministische
  Ausgabe, fehlendes C#-Plugin im gepinnten Image) oder eine der drei
  Bibliotheken nicht mehr auflösbar ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — der Bau erzeugt
den Stub über den benannten Zusatzkontext reproduzierbar, der Client öffnet
real den Server-Stream und gibt Nachrichten aus, die Handbuch-Zeile trägt
C# — **und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: die erste **reale** Bauprobe des benannten
Zusatzkontexts im offiziellen Sprach-Ziel (bislang nur eine isolierte
`docker buildx`-Probe, [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
§Kontext) — ob sie unverändert trägt oder eine Anpassung braucht, entscheidet
der Lauf und wird hier für `slice-103` (Kotlin) als übertragbare Form
festgehalten.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der benannte Zusatzkontext trägt im realen Sprach-Bau nicht** (nur die
  isolierte `buildx`-Probe aus [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  ist bislang gemessen, nicht der offizielle `examples-csharp`-Bau). —
  **Ausgang:** <…>
- **Eine der drei gRPC-/Protobuf-Bibliotheken ist zum Zeitpunkt des Baus
  nicht mehr auflösbar** (Paket zurückgezogen, Version gelöscht). —
  **Ausgang:** <…>
- **Die C#-gRPC-Werkzeugkette verlangt eine nicht-öffentliche Quelle**
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). — **Ausgang:** <…>
- **Der fremdsprachige gRPC-Bau kann seinen Stub nicht mehr aus der `.proto`
  erzeugen**, weil der Vertrag eine Einfuhr aus einem anderen Verzeichnis
  verlangt oder Codegen-Optionen den Vertrag selbst berühren müssten —
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 2. — **Ausgang:** <…>
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (ein
  vierter, teurerer Bau-Schritt je Sprache — Generator-Lauf braucht Zeit) —
  §Re-Evaluierungs-Trigger 4, `BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit` (1×, offen), `BEO-PGC/github-actions-unverifizierbar-lokal`
  (5×, verkörpert in `AGENTS.md` §3.10). — **Ausgang:** <…>
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht; besonders `slice-095`s §1, das den C#-/
  Kotlin-Weg als „ihre Bau-Kontexte erreichen sie heute nicht" beschreibt und
  mit diesem Slice teilweise überholt wird). — **Ausgang:** <…>

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
`proto/` (nur **lesend**, über den Zusatzkontext), `Makefile`/
`harness/mk/*.mk` und `docs/user/benutzerhandbuch.md` — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Keine zu
grobe Sub-Area zu differenzieren.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Fünf
Treffer: `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**,
verkörpert) und `handbuch-versionshistorie-uebersprungen` (**3×**,
verkörpert) — beide als LP3 in die DoD gezogen; `nicht-blockierender-
workflow-alarmmuedigkeit` (**1×**, offen) und `github-actions-
unverifizierbar-lokal` (**5×**, verkörpert in `AGENTS.md` §3.10) — beide
treffen den weiter wachsenden Workflow, als Risiko in §6;
`arbeit-ueberholt-stehenden-traeger` (verkörpert in `AGENTS.md` §3.13) — als
Risiko in §6, weil dieser Slice `slice-095`s Aussage zum C#-/Kotlin-Weg
teilweise überholt. Kein Treffer zu `Grpc.Tools`/`Grpc.Net.Client`/
`Google.Protobuf`/`buildx`-Zusatzkontext selbst (gemessen: `grep -rli
"grpc\.tools\|grpc\.net\|buildx" docs/plan/planning/observations/` → kein
Fund).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**.
Der Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

### Sub-Area: `*` (Default, PGC)

- **Modus:** GF
- **Konventionen-Dichte:** `harness/conventions.md` Modus-Deklaration setzt
  GF für das gesamte Repo (Doc führt, Code folgt).
- **Phase-Reife:** Phase 4 (Form entschieden, real gemessen als isolierte
  Probe — noch nicht im offiziellen Sprach-Bau belegt) — dieser Slice hebt
  die Zelle `csharp`×`grpc` auf Phase 5, sobald `examples-csharp` den Bau
  real trägt.
- **Evidenz-/Diskrepanz-Risiko:** niedrig bis mittel — GF, Doc führt; das
  einzige Diskrepanz-Risiko ist, dass die isolierte `buildx`-Probe nicht 1:1
  auf den offiziellen Sprach-Bau überträgt (siehe §6, erstes Risiko).
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
