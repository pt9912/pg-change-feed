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
**Coverage-Messfläche** (`go list ./internal/... ./cmd/...` — der Umzug nimmt
den erzeugten Code **aus** dem Gegenstand, der **Nenner** bewegt sich).

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

- [ ] **LP1 — der Stub liegt öffentlich und wird erzeugt.** Die
      `go_package`-Zeile zeigt auf einen Pfad **außerhalb** von `internal/`; der
      erzeugte Baum liegt dort; `make proto-generate` erzeugt ihn reproduzierbar
      (und `make generated-sync` bleibt grün — der Sync ist seit `slice-090` ein
      Gate). **Der neue Pfadname gehört dem umsetzenden Zug.**
- [ ] **LP2 — alle Aufrufer ziehen mit.** Die vier Importstellen kompilieren
      (`make test`), und `.a-check.yml` trägt die **neue** Gruppe `contract`
      samt ihren Kanten — **und** die `tooling`-Gruppe spricht den Stub nicht
      mehr über `adapters` an, weil die Kante
      ([`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md))
      damit gegenstandslos wird. **Die Zurücknahme ist zu benennen**, nicht still
      zu lassen.
- [ ] **LP3 — die Coverage-Fläche und ihre Träger sind nachgezogen.** Der
      erzeugte Code verlässt den Messgegenstand (`go list ./internal/...
      ./cmd/...`): **der Nenner bewegt sich** — die erreichte Quote **mit ihrem
      Lauf und ihrem Band** nennen, die betroffenen Träger
      (`harness/sensors/coverage-gate.md`, `harness/mk/coverage.mk`) prüfen und
      **die Entscheidung** festhalten, ob der erzeugte Code **im** Gegenstand
      bleibt (dann nimmt die Paketliste ihn wieder auf) oder **draußen**.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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
| `harness/sensors/coverage-gate.md` · `harness/mk/coverage.mk` | update, **falls** eine Zahl driftet oder die Paketliste sich ändert | Der Messgegenstand bewegt sich mit dem Umzug. |
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
  **Ausgang:** <…>
- **Die vier Importstellen sind nicht alle.** Gemessen sind vier; ob der Umzug
  weitere Aufrufer hat (Tests, Tools, Beispiele), weiß der Diff nicht. — **Ausgang:** <…>
- **Der Umzug verbaut den C#-/Kotlin-Weg.** Ihre gRPC-Clients erzeugen ihre Stubs
  aus der `.proto`, die **nicht** umzieht — aber ihre Bau-Kontexte sind eigene
  ([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)), und ob sie
  die `.proto` erreichen, entscheidet sich dort. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt den
  **Ort** des Erzeugnisses; welche Dokumente den **alten** Pfad nennen, weiß der
  Diff nicht (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, **3×**, seit
  `welle-20` eine Hard Rule mit `grep`-Pflicht). — **Ausgang:** <…>

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
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
