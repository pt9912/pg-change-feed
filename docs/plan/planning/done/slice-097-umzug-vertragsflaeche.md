# Slice slice-097: Umzug der Vertragsfläche — der Protobuf-Stub an einen öffentlichen Pfad

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die **Folgepflicht** aus
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §3,
nie geschnitten (Modul 6 §Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
§3 und seine Schnitt-Empfehlung **1** („**Der Umzug zuerst** … er soll allein
stehen") · [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (der
Stream und seine `.proto`) · [`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(die `tooling`-Kante, die dieser Slice zurücknimmt) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
und [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(der Messgegenstand — dieser Slice **bewegt seinen Nenner**) ·
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) (die **drei**
Sprachen) · `SPEC-018`, `SPEC-020`, `SPEC-022` (die Draht-Festlegungen, die der
Vertrag trägt) · `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (**2×**) und
`BEO-PGC/aufschub-adresse-verfaellt` (**3×**) — beide treffen die Adresse, die
dieser Slice einlöst.

**Berührte Spec-Stellen:** `SPEC-020` (der gRPC-Stream und die Form seiner
Stubs) · `ARC-005`/`ARC-007` (die Schicht-Edges, die der Umzug verschiebt).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

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

**Ziel:** Der aus der `.proto` erzeugte Go-Code zieht aus dem **privaten** Baum
an einen **öffentlichen** Pfad. Heute steht sein Ziel in der Quelle:

```text
option go_package = "github.com/pt9912/pg-change-feed/internal/adapters/driving/grpc/streamv1;streamv1";
```

Damit kann **kein** Programm außerhalb dieses Repositories den Draht ansprechen —
und **kein** Beispiel-Client im eigenen Haus ihn importieren, ohne die
`wrong-direction`-Kante zu reißen (`make a-check`, gemessen: Exit 2). **Das ist
der Grund, warum `examples/grpc-client/` seit `slice-095` nicht lieferbar ist.
**Und nur dafür:** der Umzug ist die Voraussetzung für **genau eine** der drei
`grpc`-Zellen der Matrix — die **Go**-Zelle. Die C#- und Kotlin-Clients stehen
**nicht** darauf: sie erzeugen ihre Stubs aus der `.proto`, die **nicht** umzieht,
und was ihnen fehlt, ist die `.proto` **in ihrem Bau-Kontext**
([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 2 —
gemessen: `buildx --build-context` erreicht sie, der heutige Sprach-Wurzelkontext
nicht).

**Der Rahmen (Nutzerentscheidung vom 2026-09-17): die volle Matrix.** Die
Beispiel-Clients werden über **(alle Clients) × (C#, Go, Kotlin)** geführt — die
vier Oberflächen `http`, `sse`, `grpc` und `nats` in **jeder** der drei Sprachen.
Das **erweitert** [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)
(Go die vier, C# und Kotlin bisher nur die HTTP-Familie); die Erweiterung ist eine
eigene Entscheidung und läuft als Folge-ADR. **Für diesen Slice ändert das nur
eine Pflicht:** der Umzug ist die Voraussetzung für **drei** der zwölf Programme
(`grpc` in jeder Sprache) und darf keinen von ihnen verbauen.

**Gemessene Folgepflichten** (vom Implementer von `slice-095` an beiden Routen
belegt, nicht geschätzt): die `go_package`-Zeile **und** der erzeugte Baum;
`make proto-generate` (es schreibt **in-place** in den Bind-Mount);
**vier** Importstellen (`internal/adapters/driving/grpc/server.go`, `server_test.go`,
`interceptor_test.go`, `tools/harness/grpcclient/main.go`); die
`.a-check.yml`-Gruppe `tooling` (sie darf den Stub heute ansprechen — `ADR-0068`)
und die **neue** Gruppe `contract` samt ihren Kanten; und die
**Coverage-Paketliste** (`go list ./internal/... ./cmd/...` — der Umzug nimmt
den erzeugten Code aus der **Liste**; der **Messgegenstand bleibt derselbe**,
weil `./gen/...` in der `go list`-Zeile der Stufe `coverage` wieder aufgenommen
wird — die Entscheidung ist **geschlossen**, kein offener Punkt, siehe LP3 und
Architect-Verdikt `docs/reviews/architect-verdict-slice-097-coverage-gegenstand.md`).
**Zwei Bau-Kontext-Träger, die diese Aufnahme voraussetzt** und die der
ursprüngliche Plan nicht nannte: `.dockerignore` (`!gen/`) und
`harness/sensors/generated-sync.md` §Vertrag (nennt den alten Pfad im
Präsens) — beide in §3.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der gRPC-Beispiel-Client selbst.** Er ist der **Nachfolger** und der Grund
  für diesen Slice — aber
  [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §3 sagt es selbst: der Umzug **soll allein stehen**. Ein Zug, der beides tut,
  kann bei einem roten `a-check` nicht sagen, welche Hälfte schuld ist.
- **Alle zwölf Clients der Matrix — sie gehören zum Ziel, nicht in diesen Zug.**
  Nutzerentscheidung: die Beispiel-Clients werden über **(alle Clients) ×
  (C#, Go, Kotlin)** geführt — vier Oberflächen in drei Sprachen. **Die Pflicht,
  die daraus hier erwächst:** dieser Umzug ist die Voraussetzung für die
  **Go**-Zelle der Matrix und darf sie nicht verbauen. Die C#-/Kotlin-Clients
  erzeugen ihre Stubs **aus der `.proto`** — die **nicht** umzieht —, und ihre
  Bau-Kontexte erreichen sie heute **nicht**; die Form dagegen steht gemessen in
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 2
  (ein benannter Zusatz-Kontext).
  **Die Grenze:** dieser Slice baut **keinen** Client — der Umzug soll nach
  [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §3 **allein stehen**, weil bei einem roten `a-check` sonst nicht zu sagen ist,
  welche Hälfte schuld ist.
- **Der Inhalt der `.proto`.** Der Draht-Vertrag ist
  [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)s Gegenstand; dieser
  Slice ändert **eine** Zeile an ihrem Ziel, keine Nachricht, kein Feld.
- **Ein zweites Go-Modul.** [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §Verglichene Alternativen hat es **gemessen** verworfen: unter dem Repo-Präfix
  bleibt der `internal/`-Import erlaubt, ein zweites `go.mod` ist **kein** Zaun.
- **Ein Gate für den Umzug.** Der Wächter ist `make a-check` (die Kanten) und
  `make test` (die Importe); ein eigenes Gate wäre Zeremonie.

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

- [x] **LP1 — der Stub liegt öffentlich und wird erzeugt.** Die
      `go_package`-Zeile zeigt auf einen Pfad **außerhalb** von `internal/`; der
      erzeugte Baum liegt dort; `make proto-generate` erzeugt ihn reproduzierbar
      (und `make generated-sync` bleibt grün — der Sync ist seit `slice-090` ein
      Gate). **Der neue Pfadname gehört dem umsetzenden Zug.**
- [x] **LP2 — alle Aufrufer ziehen mit.** Die vier Importstellen kompilieren
      (`make test`), und `.a-check.yml` trägt die **neue** Gruppe `contract`
      samt ihren Kanten — **und** die `tooling`-Gruppe spricht den Stub nicht
      mehr über `adapters` an, weil die Kante
      ([`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md))
      damit gegenstandslos wird. **Die Zurücknahme ist zu benennen**, nicht still
      zu lassen.
- [x] **LP3 — die Coverage-Fläche und ihre Träger sind nachgezogen.** Der
      erzeugte Code verlässt die **Paketliste**
      (`go list ./internal/... ./cmd/...`); er verlässt den **Messgegenstand
      nicht** — die tragende Regel ist die **Eigenschaft**, nicht die Liste
      ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
      Punkt 1), und
      [`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
      §Konsequenzen schreibt die Aufnahme vor (Architect-Verdikt,
      `docs/reviews/architect-verdict-slice-097-coverage-gegenstand.md`,
      Verdikt-Fall 1). Die Stufe `coverage` führt deshalb **`./gen/...`** in
      der `go list`-Zeile, mit unveränderter Filterregel. Der Nenner **bleibt**
      beim Wert des `slice-097`-Laufs — gemessen **1936** (`out` wäre
      **1850**, nicht mehr verwendet), Quote 83,3–83,4 % (Band 1 Statement) —
      und wird mit seinem Lauf im Bericht genannt. Nachzuziehen sind die
      Träger aus dem Verdikt §6 (siehe §3 dieser Datei).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: `docs/reviews/review-slice-097.md` (1 HIGH, keine Fixrunde).
- [x] Doku-Update für die öffentliche Vertragsfläche (`gen/cdc/stream/v1`) —
      berührte Träger: `harness/README.md` §Sensors (Pfadausdruck von
      `make coverage-gate`), `AGENTS.md` §4 (dieselbe Zeile), sowie
      `harness/sensors/{coverage-gate,generated-sync}.md` (§Vertrag/§Grenze/
      §Zählbasis) — alle in `57d2566`/`1c1d937` nachgezogen und vom Reviewer
      (`de0f356`) und Verifier (`b30f2f1`) unabhängig bestätigt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. **Entfällt** — Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF), keine `docs/plan/planning/reconciliation.md` vorhanden.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `proto/cdc/stream/v1/changestream.proto` | update | **Eine Zeile**: `option go_package`. Der Draht-Vertrag selbst bleibt. |
| `gen/cdc/stream/v1/**` (Name führt der umsetzende Zug) | neu | Das committete Erzeugnis am öffentlichen Pfad; der alte Baum entfällt. |
| `internal/adapters/driving/grpc/{server,server_test,interceptor_test}.go` | update | Drei Importstellen. |
| `tools/harness/grpcclient/main.go` | update | Die vierte Importstelle — der E2E-Belegträger des Streams. |
| `.a-check.yml` | update | Die neue Gruppe `contract`, ihre Kanten, und die **Rücknahme** der `tooling → adapters`-Berechtigung. |
| `Dockerfile`, Stufe `coverage` | update | `./gen/...` in die `go list`-Zeile **aufnehmen** (LP3, Verdikt „in"), Filter unverändert; Kommentarblock über der Stufe nachziehen. |
| `harness/sensors/coverage-gate.md` · `harness/mk/coverage.mk` | update | Nicht mehr bedingt — die Paketliste ändert sich sicher: §Vertrag (Pfadausdruck), §Grenze Punkt 1 (der `streamv1`-Absatz mit dem alten Pfad), §Zählbasis (Nenner **1936**, siehe LP3). |
| `harness/README.md` §Sensors · `AGENTS.md` §4 | update | Je Zeile `make coverage-gate` nennt den Pfadausdruck. |
| `.dockerignore` | update | **`!gen/` fehlt** — ohne die Zeile ist `gen/` nicht im Bau-Kontext der Stufen `coverage`/`build`, `COPY . .` legt es nicht nach `/src`. Ein Bau-Eingang derselben Art wie `!internal/`/`!cmd/` — **keine** Ausnahme nach [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) Festlegung 3 (die verlangt „gelesen, nicht gebaut" und „genau eine Datei"); ein Kommentar, der hier `ADR-0085` zitiert, wäre eine Fehl-Zitation. |
| `harness/image-hash.txt` | update, erwartet | `.dockerignore` ist eine Build-Kontext-Datei — `make image` läuft vor der Closure ([`ADR-0044`](../../adr/0044-image-beleg-semantik.md) Punkt 3); weicht der Digest ab, ist der Digest-Commit Teil des Slice. |
| `harness/sensors/generated-sync.md` §Vertrag | update | Nennt den alten Pfad im **Präsens** — nach dem Umzug falsch (§3.13-Träger). |
| `docs/user/benutzerhandbuch.md` | update, **falls** ein Pfad genannt wird | Prüfen, nicht annehmen. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`):
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
ist `Accepted`, `slice-096` liegt in `done/`, und `next/` trägt `slice-095`
(dessen §4 diesen Umzug als Bedingung nennt). Ohne Rückfrage feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn der Umzug
  **mehr als eine** Kante in `.a-check.yml` verschiebt oder ein Paket
  **außerhalb** `internal/`/`gen/` mitzieht. Dann ist der Schnitt falsch: dieser
  Slice bewegt **den erzeugten Code**, er baut die Schicht-Edges nicht um.
- `in-progress` → `open` (blockiert — Carveout?): wenn der erzeugte Code am
  öffentlichen Pfad **nicht** baubar ist, ohne den Modulpfad zu ändern. Dann ist
  das eine Entscheidung über das Modul (Architect), kein Umzug.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 sind real belegt — der Stub liegt
öffentlich **und** wird reproduzierbar erzeugt (`make generated-sync` grün), alle
Aufrufer kompilieren **und** `make a-check` ist grün mit der neuen Gruppe und der
benannten Rücknahme, und die Coverage-Fläche ist mit ihrem Lauf nachgezogen —
**und** `make gates` ist grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Naheliegender Kandidat: dieser Slice löst eine **Adresse ein**, die als
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (**2×**) und
`aufschub-adresse-verfaellt` (**3×**) im Register steht — ob daraus mehr wird als
eine Wiederholung, entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Nenner bewegt sich und die Quote mit ihm.** Der erzeugte Code
  (`streamv1`, **86** Statements, davon **76** gedeckt — Lauf `slice-091`) fällt
  aus dem Gegenstand; bei einer überdurchschnittlich gedeckten Teilmenge **sinkt**
  die Gesamtquote. Die Rampe steht seit `welle-20` auf der **Endstufe 80**. —
  **Ausgang: eingetreten, aufgefangen.** Der literale Mechanismus trat ein — der
  erzeugte Code verließ die Paketliste `./internal/... ./cmd/...` real (Umzug
  nach `gen/`). Die Konsequenz (Quoteneinbruch) trat **nicht** ein, weil LP3 —
  Teil desselben Slice, nicht eine nachträgliche Rettung — `./gen/...` in die
  `go list`-Zeile der Stufe `coverage` zurückholt (Architect-Verdikt „in",
  `docs/reviews/architect-verdict-slice-097-coverage-gegenstand.md`). Gemessen
  (Verifier, `b30f2f1`, #4/#5): `coverage-gate` druckt 83.40 % ≥ Schwelle 80 %.
  Kein Carveout, kein Folge-Slice — die Auffangmaßnahme ist bereits geliefert.
- **Die vier Importstellen sind nicht alle.** Gemessen sind vier; ob der Umzug
  weitere Aufrufer hat (Tests, Tools, Beispiele), weiß der Diff nicht. —
  **Ausgang: entfallen.** Der Verifier hat unabhängig per eigenem `grep` über
  das **ganze** Repo (nicht nur `internal/`+`cmd/`, auch `test/`, `examples/`)
  nach `gen/cdc/stream/v1` gesucht (`verify-slice-097.md` #10) und fand exakt
  vier Treffer — dieselben vier, keinen fünften. Die Unsicherheit, die das
  Risiko benannte, ist damit unabhängig geschlossen.
- **Der Umzug verbaut den C#-/Kotlin-Weg.** Ihre gRPC-Clients erzeugen ihre Stubs
  aus der `.proto`, die **nicht** umzieht — aber ihre Bau-Kontexte sind eigene
  ([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)), und ob sie
  die `.proto` erreichen, entscheidet sich dort. — **Ausgang: entfallen.**
  Geprüft (Reviewer- und Verifier-Negativbefund): die `.proto` selbst ändert
  sich um genau eine Zeile (`option go_package`), kein Nachrichten-/Feld-Inhalt;
  C#/Kotlin erzeugen ihre Stubs unverändert aus derselben Quelle und erreichen
  sie über einen eigenen Bau-Kontext
  ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 2,
  gemessen: `buildx --build-context`). Dieser Umzug hat daran nichts bewegt —
  die Entscheidungsfähigkeit von `ADR-0087`/`ADR-0090` bleibt unverändert.
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt den
  **Ort** des Erzeugnisses; welche Dokumente den **alten** Pfad nennen, weiß der
  Diff nicht (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, **3×**, seit
  `welle-20` eine Hard Rule mit `grep`-Pflicht). — **Ausgang: eingetreten.**
  Der Reviewer fand per eigenem `grep` (unabhängig vom Implementer-Suchlauf)
  einen realen Treffer:
  `docs/plan/planning/next/slice-095-beispiel-clients-drei.md:149–152`
  behauptete im Präsens, `gen/cdc/stream/v1` „existiere nicht" — seit diesem
  Umzug falsch. Direkt korrigiert in dieser Closure (siehe §7) — kein Carveout,
  kein Folge-Slice nötig, weil sofort behebbar. Zusätzlich als **fünftes**
  Auftreten der bereits verkörperten Klasse `arbeit-ueberholt-stehenden-traeger`
  vermerkt (§7) — der Implementer-Suchlauf dieses Slice fand diesen Treffer
  **nicht** (er lag außerhalb der von ihm durchsuchten Träger-Klassen), ein
  zweiter, unabhängiger Suchlauf (Reviewer) fand ihn.

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

- **Was hat funktioniert:** Die Rollen-Sequenz Architect → Planner
  (Plan-Diff) → Implementer → Reviewer → Verifier lief mit einem
  Übergabe-Artefakt je Übergang durch: Architect-Verdikt
  (`1eb3feb`) → Plan-Diff (`4f39cd5`) → drei Implementer-Commits
  (`1c1d937`, `57d2566`, `068f26e`) → Review (`de0f356`, 1 HIGH, kein
  Self-Review) → Verifikation (`b30f2f1`). Jede Rolle hat eigene, frische
  Messungen gefahren statt Vorgänger-Zahlen zu übernehmen (§3.12) — Reviewer
  und Verifier fanden je einen unabhängigen Fund (Commit-Struktur bzw.
  überholter Träger), den die vorherige Rolle nicht hatte.
- **Was ging anders als geplant:** LP3 stand im Ausgangsplan fälschlich als
  offene Entscheidung („im Gegenstand bleiben oder draußen"), obwohl
  `ADR-0071` Punkt 1 und `ADR-0076` §Konsequenzen sie bereits geschlossen
  hatten — ein Architect-Zug (Modul-8-Konflikt-Pfad, Verdikt-Fall 1) korrigierte
  den Plan vor der Implementierung. Der Implementer-Suchlauf für §3.13
  (bewegte Eigenschaft, Träger nachziehen) fand die sieben Träger aus dem
  Architect-Verdikt vollständig, verfehlte aber einen achten, außerhalb dieser
  Liste liegenden Träger (`next/slice-095`) — den fand erst ein zweiter,
  unabhängiger Suchlauf (Reviewer).
- **Steering-Loop-Eintrag:** `AGENTS.md` §3.3 (git mv + Inhaltsänderung = zwei
  Commits) wurde in Commit `1c1d937` verletzt (Umbenennung dreier erzeugter
  Dateien und Inhaltsänderung an zwei davon plus sechs weiteren Dateien im
  selben Commit) — vom Reviewer als HIGH F-1 erkannt, bewusst nicht
  fix-pflichtig gestellt (Rename-Detection griff trotzdem, `git log --follow`
  bleibt funktionsfähig), vom Verifier bestätigt. Für diese **bereits
  verkörperte** Regel gab es noch keinen Beobachtungs-Eintrag für den Fall
  „Verstoß gegen eine bestehende Regel" (anderer Beobachtungstyp als „Regel
  entsteht nach 3×") — neu angelegt als
  `BEO-PGC/git-mv-und-inhalt-in-einem-commit/`, Beleg
  `evidence/slice-097.md`. Zähler: **1×**, Zustand `offen` — mit diesem Slice
  wurde nichts verkörpert (der Normalfall bei 1×), die Teil-Zeile `— liegt
  in …` entfällt.
- **Beobachtungs-Register (`../observations/`):** Zwei Bewegungen.
  (1) `BEO-PGC/git-mv-und-inhalt-in-einem-commit/` **neu angelegt**, Beleg
  `evidence/slice-097.md` — siehe Steering-Loop-Eintrag oben.
  (2) `evidence/slice-097.md` in `BEO-PGC/arbeit-ueberholt-stehenden-traeger/`
  **ergänzt** — Zähler steht damit bei **5×** (bereits seit `welle-20` in
  `AGENTS.md` §3.13 verkörpert; dieser Fund ist ein weiteres, dokumentiertes
  Auftreten derselben Klasse, kein neuer Steering-Loop-Kandidat). Der
  überholte Satz in `docs/plan/planning/next/slice-095-beispiel-clients-drei.md`
  §4 ist in diesem Commit direkt korrigiert (Planner-Trägerpflege, kein
  fremder Code — der Melde-Schritt war bereits im Review-Report geschehen).
  Zwei im Plan zitierte Register-Einträge **geprüft, keine neue Instanz**:
  `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (weiter **2×**) — §1 dieses
  Slice schließt den gRPC-Client und die volle Matrix aus, ohne eine
  Folge-Slice-**Kennung** zu nennen (nur `ADR-0076`/`ADR-0090` als
  Ziel-Entscheidungen); das ist kein „Aufschub mit Adresse, die die Sendung
  nicht annimmt" im Sinn dieser Klasse, weil hier **keine** Adresse behauptet
  wird, die fehlschlagen könnte — die Klasse trifft eine falsche oder zu enge
  Adresse, nicht das ehrliche Fehlen einer Adresse für einen bewusst noch
  nicht geschnittenen Folge-Vorgang. `BEO-PGC/aufschub-adresse-verfaellt`
  (bereits **verkörpert** seit `slice-077`) — dieser Slice bindet keinen
  Aufschub an ein Ereignis; die drei Paarungen laufen in dieser wellenlosen
  Closure selbst (§7 unten), nicht an eine unsichere künftige Welle-Closure.
- **Folge-Slices:** keine — der gRPC-Beispiel-Client und die volle
  Client-Matrix bleiben Ziel (`ADR-0076` §3, `ADR-0090`), aber ohne
  Slice-Kennung, weil bewusst noch nicht geschnitten (siehe
  Beobachtungs-Register oben).
- **Risiken aus §6:** vier Risiken, vier Ausgänge — *eingetreten, aufgefangen*
  (Nenner/Quote, durch LP3 kompensiert) · *entfallen* (vier Importstellen,
  unabhängig verifiziert) · *entfallen* (C#-/Kotlin-Weg unberührt) ·
  *eingetreten* (überholter Träger `next/slice-095`, direkt korrigiert). Siehe
  §6 für die volle Begründung je Zeile.
- **Drei Paarungen:** Repo ohne Wellen-Betrieb — hier geprüft, nach dem `git
  mv` nach `done/` (siehe Commit 3 dieser Closure). Ergebnis wird nach
  Ausführung ergänzt: Anker-Paarung entfällt (kein `liegt in`-Feld, da nichts
  verkörpert wurde), Folge-Slice-Paarung entfällt (keine Folge-Slices
  genannt), Register-Paarung prüft `BEO-PGC/git-mv-und-inhalt-in-einem-commit`
  und `BEO-PGC/arbeit-ueberholt-stehenden-traeger` gegen ihre `evidence/`.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die Vertragsfläche
(`proto/`, `gen/`), `internal/adapters/driving/grpc`, `tools/harness/grpcclient`
und die zwei Träger — die repo-weite Default-Sub-Area `*`/`PGC` aus der
Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Vier
Treffer: `aufschub-adresse-nimmt-sendung-nicht-an` (**2×**, offen) und
`aufschub-adresse-verfaellt` (**3×**, verkörpert) — **beide treffen die Adresse,
die dieser Slice einlöst**: der Umzug ist seit `ADR-0076` benannt und nie
geschnitten. `arbeit-ueberholt-stehenden-traeger` (**3×**, verkörpert in
`AGENTS.md` §3.13) und `beleg-befehl-traegt-seinen-satz-nicht` (**5×**,
verkörpert) als Risiko in §6.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**. Der
Block pro Sub-Area entfällt; der **Abschnitt** bleibt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

### Sub-Area: <Name>

- **Modus:** GF | BF | Hybrid
- **Konventionen-Dichte:** <Beleg aus `harness/conventions.md`,
  Adaptions-Block oder Code>
- **Phase-Reife:** Phase 0–5 <Begründung gegen die Phase × Modus-Matrix>
- **Evidenz-/Diskrepanz-Risiko:** <bei BF/Hybrid: was kann die
  Inventur sichtbar machen? bei GF: meist niedrig>
- **Reconciliation-Aufwand:** <Slice-Schätzung;
  Graduation-/Folge-Slice-Trigger>
