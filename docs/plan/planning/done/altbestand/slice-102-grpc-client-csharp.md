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
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-102.md`
      liegt vor (Modul 11, frischer Kontext) — DoD-konform, mit zwei benannten,
      nicht-blockierenden Einschränkungen.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) — *entfällt: Greenfield-
      Bootstrap (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF),
      keine Datei vorhanden.*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im
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
  **Ausgang: entfallen.** Implementer und Verifier bestätigen unabhängig
  voneinander, dass die isolierte Probe unverändert trägt: der Bau bricht
  ohne den benannten Zusatzkontext real mit Exit 1 ab
  (`docs/reviews/verify-slice-102.md` #1) und erzeugt den Stub mit ihm real,
  auch bei einer unabhängigen `--no-cache`-Gegenprobe (#3).
- **Eine der drei gRPC-/Protobuf-Bibliotheken ist zum Zeitpunkt des Baus
  nicht mehr auflösbar** (Paket zurückgezogen, Version gelöscht). —
  **Ausgang: entfallen.** Alle drei Pakete lösten in der `--no-cache`-Probe
  des Verifiers real gegen NuGet auf und sind zusätzlich die aktuell
  neuesten stabilen Versionen, exakt die gepinnten
  (`docs/reviews/verify-slice-102.md` #3/#5).
- **Die C#-gRPC-Werkzeugkette verlangt eine nicht-öffentliche Quelle**
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 3). — **Ausgang: entfallen.** Beide Bauproben des
  Verifiers liefen ausschließlich gegen die öffentliche NuGet-Registry, ohne
  Zugangsdaten (`docs/reviews/verify-slice-102.md` #3, #5).
- **Der fremdsprachige gRPC-Bau kann seinen Stub nicht mehr aus der `.proto`
  erzeugen**, weil der Vertrag eine Einfuhr aus einem anderen Verzeichnis
  verlangt oder Codegen-Optionen den Vertrag selbst berühren müssten —
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md)
  §Re-Evaluierungs-Trigger 2. — **Ausgang: entfallen.** Der Bau erzeugt den
  Stub real und reproduzierbar; die unabhängige `--no-cache`-Gegenprobe des
  Verifiers lief fehlerfrei durch, alle sieben `GrpcClient.Tests` real neu
  gelaufen (`docs/reviews/verify-slice-102.md` #3).
- **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (ein
  vierter, teurerer Bau-Schritt je Sprache — Generator-Lauf braucht Zeit) —
  §Re-Evaluierungs-Trigger 4, `BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit` (1×, offen), `BEO-PGC/github-actions-unverifizierbar-lokal`
  (7×, verkörpert in `AGENTS.md` §3.10). — **Ausgang: weiter offen**, ohne
  neue Evidenz in einem der beiden zitierten Register-Einträge. Implementer
  und Verifier bestätigen unabhängig: `.github/workflows/examples.yml`
  selbst ist von diesem Slice **nicht** strukturell geändert (Implementer:
  „`.github/workflows/examples.yml` itself was not touched"), `AGENTS.md`
  §3.10 wird nicht neu ausgelöst — wie bereits bei `slice-100`/`slice-101`.
  `nicht-blockierender-workflow-alarmmuedigkeit` bleibt bei 1×,
  `github-actions-unverifizierbar-lokal` bleibt bei real 7× (die im Plan
  selbst zitierte „5×" war eine bei Niederschrift bereits veraltete
  Übernahme, siehe Korrektur oben und
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-102.md`
  — kein Drift durch die Arbeit *dieses* Risikos).
- **Der Handbuch-Nachzug wird vergessen** — die Klasse mit **je 3×** in zwei
  Registereinträgen. — **Ausgang: entfallen.** Reviewer (Negativbefund) und
  Verifier (`docs/reviews/verify-slice-102.md` #9) bestätigen unabhängig den
  `**Beispiele:**`-Block mit der C#-Zeile und die Änderungshistorie-Zeile
  1.24. Beide Registereinträge
  (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  `handbuch-versionshistorie-uebersprungen`) sind **eingelöst**, nicht
  verletzt — keine neue Evidenzdatei, beide bleiben bei 3×.
- **Ein Träger wird überholt, den dieser Slice nicht anfasst**
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, verkörpert in `AGENTS.md`
  §3.13 — Suchlauf-Pflicht; besonders `slice-097`s §1, das den C#-/
  Kotlin-Weg als „ihre Bau-Kontexte erreichen sie heute nicht" beschreibt und
  mit diesem Slice teilweise überholt wird). — **Ausgang: eingetreten, direkt
  behoben — mit einem zweiten, unterschiedenen Fund derselben Klasse.** (a)
  Der Implementer-eigene §3.13-Suchlauf fand eine stale Aussage in
  `harness/README.md` §Sensors („drei Images" → „vier Images" je Sprachziel,
  ausgelöst durch das vierte Runtime-Image im C#-Bau) und korrigierte sie im
  selben Commit (`a71b425`). (b) `slice-097`s §1 (`done/`, immutabel) trägt
  seit diesem Slice eine **teilweise** überholte Aussage — „ihre
  Bau-Kontexte erreichen sie heute nicht" gilt seit `slice-102` für C#
  nicht mehr, für Kotlin weiterhin. Das ist kein Zitat-Korrektur-Fall nach
  `ADR-0073` (die Aussage war bei ihrer Niederschrift wahr, wird erst durch
  spätere Arbeit überholt) und keine Editier-Gelegenheit — `AGENTS.md`
  §3.13 verlangt hier **Melden, nicht Ändern**. Beide Funde gehören zum
  selben Vorgang (`slice-102`) und damit zu **einer** Evidenzdatei
  (`evidence/slice-102.md`), Zähler **8× → 9×** (kein neuer
  Schwellen-Übertritt, die Regel steht bereits seit `welle-20`).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Die isolierte `buildx`-Probe aus
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) §Kontext trug
  unverändert im offiziellen `examples-csharp`-Sprachbau — sowohl die
  negative Probe (Bau bricht ohne den benannten Zusatzkontext real ab,
  `docs/reviews/verify-slice-102.md` #1) als auch die positive, cache-lose
  Gegenprobe (Stub-Erzeugung + alle vier Testsuiten real neu gelaufen, #3).
  Das bestätigt die Form für `slice-103`, der sie unverändert auf Kotlin
  überträgt.
- **Was ging anders als geplant:** Der benannte Zusatzkontext wurde nicht nur
  für den `grpc`-Aufrufpfad zwingend, sondern strukturell für **alle vier**
  `docker build`-Aufrufe von `examples-csharp` (Review F-1) — eine
  Nebenwirkung der gemeinsamen `build`-Stufe, die der Implementer innerhalb
  des vorgegebenen Dockerfile-Musters nicht vermeiden konnte, ohne die Stufe
  je Programm aufzuspalten. Reviewer und Verifier stufen das als akzeptabel
  und transparent dokumentiert ein, es ist aber eine Abweichung von der
  engeren Wortlaut-Erwartung „für grpc zusätzlich" in `ADR-0090` §Fitness
  Function.
- **Steering-Loop-Eintrag:** Die Docker-Fehlermeldung bei fehlendem
  Zusatzkontext („pull access denied … `docker.io/library/proto:latest`")
  ist irreführend (Review F-2) — sie deutet ohne Vorwissen auf ein
  Registry-/Auth-Problem statt auf einen fehlenden Bau-Kontext. Entscheidung:
  **keine neue Hard Rule/kein neuer Sensor** — die Milderung existiert
  bereits als erklärender Kommentar direkt über der `COPY`-Zeile in
  `examples/csharp/Dockerfile` sowie in `harness/mk/examples.mk`/
  `harness/README.md`, und dieses Muster überträgt sich naturgemäß, sobald
  `slice-103` dieselbe Dockerfile-Form auf Kotlin kopiert. Bleibt eine
  **Notiz für `slice-103`** (dort ergänzt, siehe Folge-Slices), kein
  Registereintrag — die Fehlermeldung selbst ist Docker-Verhalten, nicht
  reparierbar, nur dokumentierbar.
- **Beobachtungs-Register (`../observations/`):** Drei Einträge
  fortgeschrieben: `arbeit-ueberholt-stehenden-traeger` jetzt **9×**
  (`evidence/slice-102.md`, zwei Funde desselben Vorgangs — (a)
  `harness/README.md` „drei"→„vier" Images, direkt behoben in `a71b425`,
  (b) `slice-097`s §1, teilweise überholt, gemeldet statt geändert, da
  `done/`-immutabel; beide zählen als **eine** Gelegenheit); `zahl-in-
  traeger-driftet-gegen-die-messung` jetzt **11×** (`evidence/slice-102.md`,
  die im eigenen Plan-Kopf/§6/§8 bei Niederschrift bereits veraltete „5×"
  für `github-actions-unverifizierbar-lokal`, real 7×, vom Verifier
  gefunden und in diesem Zug korrigiert); `zitat-nennt-die-falsche-stelle`
  jetzt **2×** (`evidence/slice-102.md`, Review-F-3: Plan zitierte
  `slice-095` statt `slice-097` als Beleg für „Bau-Kontexte erreichen sie
  heute nicht" — weiterhin unter der 3×-Schwelle, bleibt `offen`). Zwei
  Einträge eingelöst, keine neue Evidenzdatei:
  `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (bleibt 3×)
  und `handbuch-versionshistorie-uebersprungen` (bleibt 3×). Zwei Einträge
  bewusst **nicht** fortgeschrieben: `github-actions-unverifizierbar-lokal`
  (bleibt real 7×, nur der Plan zitierte veraltet — siehe oben) und
  `nicht-blockierender-workflow-alarmmuedigkeit` (bleibt 1×) — Begründung in
  §6, Risiko 5: dieser Slice ändert `.github/workflows/examples.yml` nicht
  strukturell.
- **Folge-Slices:** keine neuen. `slice-103` existierte bereits vor dieser
  Closure in `open/`; sein Plan enthält bereits die konkrete
  Docker-Befehlsform (§1, Bullet 1: `docker buildx build --build-context
  <name>=<proto-Verzeichnis> …`, `COPY --from=<name> …`) und wurde um einen
  knappen Hinweis ergänzt, den erklärenden Dockerfile-Kommentar über der
  `COPY`-Zeile (siehe Steering-Loop-Eintrag oben) beim Übertragen auf Kotlin
  mitzuführen.
- **Risiken aus §6:** sieben Zeilen, sieben Ausgänge — fünf **entfallen**
  (Zusatzkontext trägt real; alle drei Bibliotheken auflösbar; keine
  nicht-öffentliche Quelle; Stub-Erzeugung real und reproduzierbar;
  Handbuch-Nachzug gemessen vollständig), eines **weiter offen**
  (nicht-blockierender Workflow — ohne neue Registerevidenz, da dieser
  Slice den Workflow nicht strukturell ändert) und eines **eingetreten,
  direkt behoben, mit einem zweiten, gemeldeten (nicht geänderten) Fund**
  derselben Klasse (`harness/README.md` sowie `slice-097`s §1 →
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, jetzt 9×). Details je Zeile
  in §6.
- **Drei Paarungen:** Anker-Paarung entfällt (keine neue Verkörperung durch
  diesen Slice — alle drei fortgeschriebenen Register-Einträge waren entweder
  bereits verkörpert oder bleiben unter der Schwelle). Folge-Slice-Paarung
  entfällt (keine neuen Folge-Slices; `slice-103` existierte bereits vor
  dieser Closure in `open/`). Register-Paarung **grün** — alle sieben
  zitierten Beobachtungs-Verzeichnisse
  (`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  `handbuch-versionshistorie-uebersprungen`,
  `nicht-blockierender-workflow-alarmmuedigkeit`,
  `github-actions-unverifizierbar-lokal`,
  `arbeit-ueberholt-stehenden-traeger`,
  `zahl-in-traeger-driftet-gegen-die-messung`, `zitat-nennt-die-falsche-
  stelle`) existieren mit nicht leerem `evidence/`. Letztes DoD-Häkchen
  bestätigt gegen den mv-Commit `564ac52`.

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
unverifizierbar-lokal` (**7×** — bei Niederschrift dieses Plans fälschlich
als „5×" zitiert, dieselbe Klasse wie `slice-101`, siehe §6 letzter
Risiko-Punkt sowie
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-102.md`;
verkörpert in `AGENTS.md` §3.10) — beide
treffen den weiter wachsenden Workflow, als Risiko in §6;
`arbeit-ueberholt-stehenden-traeger` (verkörpert in `AGENTS.md` §3.13) — als
Risiko in §6, weil dieser Slice `slice-097`s Aussage zum C#-/Kotlin-Weg
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
- **Phase-Reife:** Phase 5 (etabliert) — bei Plan-Niederschrift Phase 4 (nur
  isolierte Probe); dieser Slice hat die Zelle `csharp`×`grpc` mit dem
  realen `examples-csharp`-Bau (`docs/reviews/verify-slice-102.md` #1/#3)
  auf Phase 5 gehoben.
- **Evidenz-/Diskrepanz-Risiko:** niedrig bis mittel — GF, Doc führt; das
  einzige Diskrepanz-Risiko ist, dass die isolierte `buildx`-Probe nicht 1:1
  auf den offiziellen Sprach-Bau überträgt (siehe §6, erstes Risiko).
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bootstrap).
