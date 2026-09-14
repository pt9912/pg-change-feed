# Verifikationsbericht: slice-069 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-069`
§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5 Closure-Trigger, §6
Risiken, §8 Sub-Area) und die bindenden Entscheidungen
([`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) — umgesetzt,
in der Puffer-Klausel durch die Folge-ADR abgelöst;
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
— `Accepted`, `Supersedes` [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
**nur** die Puffer-Klausel, Prüfgrundlage dieses Laufs; das Architect-Verdikt
[`architect-verdict-publish-blockiert-capture-pfad`](architect-verdict-publish-blockiert-capture-pfad.md)
vollständig gelesen). Nicht gegen den Diff als solchen (Reviewer-Aufgabe;
[`review-slice-069`](review-slice-069.md) vollständig gelesen, inkl. Fixrunden-
Vermerk, aber nur als Kontext) und nicht gegen realen Bedarf (Validator — hier
nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8,
Stand `HEAD = b23ac0e`), die vollständige
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md), die vollständige
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
das vollständige Architect-Verdikt, den vollständigen Review-Report (beide
Läufe), `SPEC-020` (`spec/pflichtenheft.md:339-360`), `broadcaster.go`,
`broadcaster_test.go`, `interceptor.go`, `interceptor_test.go`,
`server_test.go`, `changestream.go`, `welle-19.md` und den ADR-Index. Die
Sensoren wurden in dieser Sitzung **eigenständig real ausgeführt** (Exit-Code
je in einem eigenen, ungepipten, mechanisch konditionierten Schritt,
`AGENTS.md` §3.9) — kein Implementer- oder Reviewer-Beleg ungeprüft
übernommen. Beide Mutationsstichproben sind **eigene**, selbst gesetzt.

**Gegenstand:** `docs/plan/planning/in-progress/slice-069-grpc-streaming-adapter-grundgeruest.md`
zum Stand `HEAD = b23ac0e`. Diff-Fenster `95ca52c..HEAD`: drei Commits
(`0867835` Broadcaster nicht-blockierend mit begrenzter Empfangswarteschlange,
`2f6882c` `SPEC-020`-Präzisierung, `b23ac0e` Fixrunden-Verdikt). Das
Diff-Fenster enthält **nicht** die Anlage von `ADR-0066` und das Verdikt — die
liegen bereits in `95ca52c` (`249dad5`/`4df23a3`) und sind Kontext, nicht
Gegenstand.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — die implementierungs-/
reviewbezogenen Zeilen sind Prüfgegenstand, die Closure-Zeilen **müssen** offen
bleiben, bis der Planner sie schließt.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `LH-FA-SST-008` (Authn/Authz-Boundary): Öffnung ohne/mit unbekanntem `authorization` endet `Unauthenticated`; gültiges `reader`/`admin` öffnet | **erfüllt, selbst nachvollzogen** | `interceptor.go:91-98` (`authStreamInterceptor` → `Unauthenticated` bei `roleNone`); `server_test.go:125-156` drei Negative, `:162-208` zwei Erfolgsfälle mit vollständigem Change-Inhalt; eigener `make test`-Lauf grün (§2). |
| 2 | `ADR-0060` Teilfrage 1/2/4/5/6: Port `ChangeStreamPort`, `Broadcaster`, gRPC-Server `+` Auth-Interceptor, additive `CDC_GRPC_ADDR`-Verdrahtung, Docker-only Proto-Stufe | **erfüllt** — mit einer Text-Rest-Beobachtung (unten) | `changestream.go:30-46` (Port `+` Sentinel); `broadcaster.go` (Broadcaster); `server.go`/`interceptor.go` (Server/Interceptor); `wiring.go:115,678-700` (`CDC_GRPC_ADDR`, No-Op bei leer); `Dockerfile:33-36` (Stufe `proto`), `Makefile:102` (`proto-generate`); `proto/cdc/stream/v1/changestream.proto` `+` committeter `streamv1/*.pb.go`. `grpcstream.New()` steht **nur** im `apigrpc.Config` (`wiring.go:691`) — kein `ChangeStreamPort` am `CaptureService` (Capture-Integration aus §1 ausgeschlossen). |
| 3 | Fixrunde (F-1/F-2, `ADR-0066`): `Publish` nicht-blockierend `+` begrenzte Empfangswarteschlange je Abonnent (Drop-Newest, Kanal nie geschlossen), Godoc nachgezogen, zwei Regressionstests | **erfüllt, selbst reproduziert** | `broadcaster.go:36` (`queueCapacity = 64`), `:67` (begrenzter Kanal), `:110-115` (`select` mit `default`); Godoc `:1-12,25-36,56-65,85-95`; Port-Godoc `changestream.go:31-44`. Beide Kern-Wächter **selbst rot gesehen** (§3, Mutation A), der Leer-Token-Test **selbst als paartragend** belegt (§3, Mutationen B/C). |
| 4 | `spec/pflichtenheft.md` erhält `SPEC-020` (Protobuf-Schema, RPC-Name, Stream-Semantik) | **erfüllt** | `spec/pflichtenheft.md:339-360`, neun Merkmals-Zeilen inkl. *Zustellgarantie* (`:356`) und *Erzeuger-Blockade* (`:357`); deckungsgleich mit `changestream.proto` (Dienst/RPC) und `server.go` (Felder). |
| 5 | `make gates` grün | **erfüllt, selbst ausgeführt** | Eigener Lauf Exit **0** (§2). |
| 6 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | [`review-slice-069`](review-slice-069.md) vorhanden (2 Läufe, Fixrunden-Vermerk); 0 HIGH / 0 MEDIUM im zweiten Lauf, deckungsgleich mit §2 der DoD-Zeile. |
| 7 | Doku-Update `harness/README.md` §Werkzeuge `+` `AGENTS.md` §4 (`make proto-generate`) | **erfüllt** | `harness/README.md:127` (Werkzeug-Zeile, „kein Gate"), `AGENTS.md:315` (§4-Tabelle); das Target existiert real (`Makefile:102`, `.PHONY`). Keine halluzinierte Zeile. |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<…>`-Platzhalter (Volltext gelesen). Planner-Arbeit nach diesem Bericht. |
| 9 | Reconciliation-Register — entfällt (GF-Repo) | **korrekt offen, Entfall-Vermerk trägt** | `docs/plan/planning/reconciliation.md` existiert real nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/` real gelistet: **kein** `slice-069`-Beleg (`grep -rl "slice-069" observations/` leer). Die §8-Sichtung nennt die Treffer; die Eintragung gehört in die Closure. |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Drei §6-Einträge tragen wörtlich „Ausgang: wird bei Closure zugewiesen"; das vierte (Broadcaster-Blockade) trägt bereits **entfallen** (`:259-268`) und ist über [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md) belegt. Kein Risiko vorzeitig geschlossen; die Zeile bleibt korrekt offen, weil drei Ausgänge fehlen. |
| 12 | Die drei Paarungen | **korrekt offen** | Slice liegt real in `in-progress/`; `Welle: welle-19`-Feld vorhanden; die DoD-Zeile verweist korrekt auf die `welle-19`-Closure. |

**Ergebnis §1:** Alle sieben implementierungs-/reviewbezogenen DoD-Punkte (1–7)
sind real erfüllt; 1, 3 und 5 wurden **selbst reproduziert**, die Kernzusage
zusätzlich über eigene Mutationen als tragend belegt. Die fünf Closure-Punkte
(8–12) sind korrekt offen und nicht vorweggenommen.

**Rest-Beobachtung (Textpräzision, kein DoD-Verstoß).** Die DoD-Zeile 2
beschreibt den `Broadcaster` weiterhin mit dem Wort „**kein Puffer**", während
die gleichfalls abgehakte DoD-Zeile 3 (Fixrunde, `ADR-0066`) „begrenzte
Empfangs-Warteschlange je Abonnent (Drop-Newest)" verlangt und der ausgelieferte
Code genau das tut (`broadcaster.go:67`, Kapazität 64). Der Widerspruch ist in
§3 ausdrücklich als „korrigiert durch Architect-Verdikt,
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)"
dokumentiert (`:182-190`), und
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
löst die „kein Puffer"-Klausel ab — die spätere, speziellere DoD-Zeile 3 trägt
den Punkt, das Häkchen steht zu Recht. Nicht deckungsgleich ist allein der
Wortlaut von Zeile 2: Wer sie isoliert liest, sieht eine zugesagte Eigenschaft
(„pufferlos"), die das Artefakt nicht erfüllt. Der Reviewer hat das nicht
geführt (sein F-1 galt `SPEC-020`/Godoc, nicht dem DoD-Text). Empfehlung ohne
Zwang: die Zeile 2 beim Closure-Nachzug an Zeile 3 angleichen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien`; d-check 543 Dateien/0 Befunde (links/anchors/ids/matrix/versions/structure); d-check `--enable commits --range HEAD~5..HEAD` 0 Befunde; `commit-traceability: OK — 5 Commit(s)`, Betreffe ohne Struktur-ID; a-check gesamt 0 Befunde (2 unverschichtete Werkzeug-Dateien, bekannte Hinweise); `coverage-gate: OK — Coverage 48.30% erfüllt Schwelle 35%` (bootstrap-aware Einstiegsstufe `ADR-0054`) |
| `make test` | **0** | 28 Pakete `ok`, **kein `FAIL`** im Lauf; `grpcstream` und `grpc` grün |
| `make a-check` | **0** | 0 Befunde gegen den vollen Baum; `.a-check.yml` unverändert (§6) |
| `make doc-commits RANGE=4085f18..HEAD` | **0** | 543 Dateien, 0 Befunde — jeder Commit des Slice-Fensters ist kennungstragend |
| `make doc-immutable RANGE=4085f18..HEAD` | **0** | 0 Befunde — keine `Accepted`-ADR überschrieben; `docs/plan/adr/`-Diff im Fenster: nur `0066-…md` neu `+` `README.md`-Index, **`ADR-0060`-Datei unberührt** |

`git status --porcelain` nach allen Sensor- und Mutationsläufen: **leer**
(§3).

## 3. Mutations-Stichprobe — Kernzusage und F-2-Paar selbst rot/grün gesehen (Verifier-only-Nachweis)

Der Reviewer nennt vier Mutationen mit Exit-Codes. Ich habe die tragenden
davon **selbst gestellt**, jede in einem eigenen `make test`-Aufruf, danach per
`git checkout --` zurückgenommen. Der Blob-Hash der zurückgenommenen Datei
stimmt mit `HEAD` überein, `git status --porcelain` ist nach jeder Rücknahme
leer.

| # | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| A | Sende-Zweig in `Publish` blockierend: `select { case changes <- change: default: }` → `changes <- change` (`broadcaster.go:110-115`) | rot genau an der Kernzusage `ADR-0066` Festlegung 1 | **rot** (Exit 2; genau `TestNichtLesenderAbonnentHaeltPublishNichtAn` und `TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen`, beide `2.00s`-Frist, `broadcaster_test.go:123`/`:146`) |
| C | Nur der Leer-Token-Frühausstieg entfernt: `if token == "" { return roleNone }` (`interceptor.go:52-54`) | grün — die Restschranke allein trägt | **grün** (Exit 0) — bestätigt F-10 |
| B | Frühausstieg **und** beide `!= ""`-Wächter entfernt (`interceptor.go:51-60`) | rot in **beiden** Hälften des F-2-Tests | **rot** (Exit 2; `TestClassifyTokenLeereKonfigurationTrifftKeinToken` und `TestStreamChangesLeereTokenKonfigurationEndetMitUnauthenticated`) |

**Rücknahme und Blatt-Identität:** nach `git checkout --` stimmt `git
hash-object` beider Dateien mit `git rev-parse HEAD:<pfad>` überein
(`broadcaster.go` `36672aa9…`, `interceptor.go` `49f98fee…`),
`git status --porcelain` ist leer, kein Rest im Baum.

**Der Wächter pulsiert real über die Kapazität — eigenständig geprüft.** Der
Test `TestNichtLesenderAbonnentHaeltPublishNichtAn` looped bis
`queueCapacity+1` (`broadcaster_test.go:122`) — **nicht** bis zu einer
festen 64. Er ist damit strukturell an die Konstante gekoppelt: Wächst die
Kapazität, wächst die Pulszahl mit, und der Test kann nicht blind werden (ein
einzelner Change passte auch in einen blockierenden Send auf einen 64er-Puffer
und ließe ihn grün). Mutation A fällt folgerichtig erst nach der 65. Übergabe in
die Frist. `TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen` prüft
zusätzlich die Reihenfolge der 64 zugestellten Changes und die Nicht-Zustellung
des 65. — die **Richtung** (Drop-Newest, der eintreffende fällt weg) ist damit
testgetragen, nicht nur kommentiert. `ADR-0066` Festlegung 1/2/3 ist real
testgedeckt.

**Ergebnis §3:** Die Kernzusage aus `ADR-0066` ist **einzeln** tragend und
**rot gesehen** — genau an der zugesagten Stelle. Die Leer-Token-Grenze ist
**paartragend**: einzeln (nur Frühausstieg entfernt) bleibt der Lauf grün, erst
das Entfernen **beider** Hälften färbt beide Testhälften rot. F-10s Befund ist
damit unabhängig bestätigt.

## 4. `ADR-0066` gegen den Code — Deckungsgleichheit Punkt für Punkt

`broadcaster.go` vollständig gelesen. Je Festlegung:

- **Festlegung 1 — Senke hält Quelle nie auf, strukturell.** `Publish` kopiert
  die Ziele unter `mu` (`:103-108`) und bedient jeden mit `select { case …:
  default: }` (`:110-115`) — kein `done`, kein `ctx.Done()` im Handoff; ein
  bereits beendetes `ctx` wird vorab geprüft (`:100-102`). Der Aufruf kehrt
  unabhängig vom Lese-Zustand des Abonnenten zurück. **Deckungsgleich.**
- **Festlegung 2 — begrenzte Warteschlange je Abonnent.** `Subscribe` legt
  `make(chan *model.Change, queueCapacity)` an (`:67`), `queueCapacity = 64`
  als benannte Paket-Konstante (`:36`); Signatur `Subscribe() (<-chan
  *model.Change, func())` unverändert, Kapazität für den Abonnenten unsichtbar.
  **Deckungsgleich.**
- **Festlegung 3 — Drop-Newest.** Ist die Warteschlange voll, greift der
  `default`-Zweig; der **eintreffende** Change fällt weg, die eingereihten
  bleiben unverdrängt (`:110-115`). Keine Umsortierung. **Deckungsgleich**
  (`TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen` belegt es real).
- **Festlegung 3 — Kanal nie geschlossen, bei `cancel` verworfen, kein Replay,
  keine Position.** `cancel` löscht nur den Map-Eintrag (`:74-81`, idempotent
  über `sync.Once`); `changes` wird an **keiner** Stelle geschlossen (Grep über
  die Datei: kein `close(`). Nach `cancel` bedient kein `Publish` den Kanal
  mehr, die Warteschlange ist mit dem Eintrag verworfen. Der `Broadcaster`
  führt **keine** Position, keine Abonnenten-Identität über die Subskription
  hinaus; ein neues `Subscribe` liefert einen frischen, leeren Kanal
  (`TestAbmeldenIstIdempotent` prüft den entfernten Empfänger). **Deckungsgleich
  — nicht die verworfene Option B/C** (kein Nachliefern).
- **Festlegung 4 — Wecksignal-Charakter.** Kein Log, keine Metrik, kein
  Zustellversprechen; die Nachvollziehbarkeit bleibt beim Store. Im Code: kein
  Log-Aufruf im Drop-Pfad. **Deckungsgleich.**

## 5. `SPEC-020`-Zeilen und Port-Godoc gegen den Code

- **Zeile *Zustellgarantie*** (`spec/pflichtenheft.md:356`): „nicht verbunden
  **oder langsamer lesend** … verpasst ersatzlos — je Abonnent … begrenzte
  Empfangs-Warteschlange, deren Überlauf verworfen wird". Deckt beide Hälften
  (verwerfende **und** entkoppelnde), wie das Verdikt es verlangt. **Deckungsgleich**
  — mit der Richtungs-Rest-Beobachtung F-9 (unten §7).
- **Zeile *Erzeuger-Blockade*** (`:357`): „`Publish` blockiert nie auf einen
  Abonnenten … Der Aufruf liefert nur bei ungültigem Aufruf (**fehlender Change
  oder bereits beendeter Aufruf-Kontext**) einen Fehler". Der Code hat genau
  **zwei** Fehlerausgänge: `change == nil` → `fmt.Errorf("%w: …",
  ErrChangeStream)` (`:97-99`) und `ctx.Err()` → roher Kontext-Fehler
  (`:100-102`). Die Zeile nennt **beide**. **Deckungsgleich.**
- **Port-Godoc `Publish`** (`changestream.go:31-44`): nennt „nicht-blockierend",
  „begrenzte Empfangs-Warteschlange", „Drop-Newest, nicht nachgeliefert", die
  fehlende Zustellgarantie und beide Fehlerausgänge („nur aus einem ungültigen
  Aufruf (`ErrChangeStream`) oder einem bereits beendeten `ctx`").
  **Deckungsgleich** mit dem Code.

**Rest-Beobachtung (kein DoD-Verstoß, benannt).**
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
Festlegung 1 schließt mit „Der `Publish`-Fehler entsteht weiter nur aus einem
ungültigen Aufruf (`ErrChangeStream`)". Der ausgelieferte Code liefert für den
beendeten `ctx` den **rohen** Kontext-Fehler (nicht in `ErrChangeStream`
gehüllt) — die spätere `SPEC-020`-Zeile (`:357`) und der Port-Godoc nennen
diesen zweiten Ausgang ausdrücklich und ordnen ihn korrekt zu (der Aufruf bricht
**vor** jeder Verteilung ab), sodass die **ausgelieferte Semantik vollständig
gepinnt** ist. Die ADR-Klammer „(nur … `ErrChangeStream`)" ist enger als der
Code; da die ADR `Accepted` und nach `AGENTS.md` §3.5 nicht nachbesserbar ist,
bleibt das eine benannte Textspanne, kein Umsetzungsfehler.

## 6. Unverändert-Prüfungen (Auftrags-Punkt 6) und DoD-Häkchen-Trennung

- **`ChangeNotificationPort`/`natsnotify`: kein Zeilendiff.** `git diff
  --name-only 95ca52c..HEAD -- internal/application/port/outbound/changenotification.go
  internal/adapters/driven/natsnotify` liefert **keine** Ausgabe. `ADR-0060`
  Teilfrage 2 Option A (Port erweitern) ist nachweislich unausgeführt; der neue
  `ChangeStreamPort` ist ein getrenntes Interface (`changestream.go`).
- **`.a-check.yml`: kein Zeilendiff.** `git diff --name-only 95ca52c..HEAD --
  .a-check.yml` leer; `make a-check` 0 Befunde. `ADR-0060` Teilfrage 5 Punkt 2
  eingehalten.
- **`internal/adapters/driving/http`: nur der Gegenverweis-Kommentar.**
  Der einzige Diff in `middleware.go` ist der Kommentar-Block (`:26-32`):
  `slice-059`-Chronik-Verweis durch die Invariante ersetzt und der Gegenverweis
  auf die zweite `classifyToken`-Fassung ergänzt — **die Funktion selbst ist
  zeichengleich**. Deckt F-3 und `AGENTS.md` §3.7 (Kommentar beschreibt, was da
  ist). `.a-check.yml` führt keine `adapters→adapters`-Kante, der Gegenverweis
  ist damit die zulässige Träger-Form der Duplikation.
- **DoD-Häkchen-Trennung — korrekt.** Der DoD-Block ist in drei Zustände
  getrennt: Zeilen 1–7 (Implementierung, Fixrunde, `SPEC-020`, Gate, Review,
  Doku) `[x]`; Zeile 9 (Reconciliation) `[x]` mit Entfall-Vermerk im selben
  Satz; Zeilen 8, 10, 11, 12 (Closure-Notiz, Register, Risiko-Ausgänge,
  Paarungen) `[ ]`. Der Reviewer hat in der Fixrunde **genau** die
  Review-Zeile nachgezogen (`.harness/skills/reviewer.md`
  §DoD-Checkbox-Nachzug ohne Fixrunde), keine Closure-Zeile vorweggenommen.

## 7. Urteil zu den zwei neuen INFO-Findings (F-9, F-10)

**F-9 („`SPEC-020` nennt die Verwerfungsrichtung nicht") — Closure-Notiz, kein
Blocker.** Die Zeile *Zustellgarantie* sagt „deren Überlauf verworfen wird",
ohne Drop-Newest vs. Drop-Oldest zu benennen. Die Richtung ist an drei anderen
Stellen gepinnt: `ADR-0066` Festlegung 3, der Port-Godoc (`changestream.go:34-36`,
„Drop-Newest, nicht nachgeliefert") und der Paket-Kommentar
(`broadcaster.go:5-7`); der Test belegt sie real. Die DoD verlangt für
`SPEC-020` „Stream-Semantik", nicht die Richtung im Klartext. Die Zeile ist
**nicht falsch**, nur weniger präzise als die übrigen Träger. Für die
Slice-Closure gehört die Klasse („Verwerfungsrichtung im Spec-Text nicht
benannt") ins Register; ein Implementer-Rückweg ist nicht nötig.

**F-10 („Mutations-Beleg greift die redundante Hälfte") — Closure-Notiz, kein
Blocker.** Der Befund ist eine **Methoden**-Beobachtung über den Nachweis, kein
Code-Defekt. Ich habe ihn unabhängig bestätigt (§3): der Frühausstieg allein
entfernt → grün (Exit 0); Frühausstieg **und** beide `!= ""`-Wächter entfernt →
rot (Exit 2). Die Verteidigung im Code ist doppelt gedeckt (redundant, aber
harmlos — Defense in depth), und der **tragende** Fall ist das Paar, das der
F-2-Test real übt. Ein Beleg, der nur die redundante Hälfte mutiert, weist
allerdings keine Sensitivität nach — das ist der Punkt, und er verlangt keine
Code-Änderung. Klasse ins Register bei der Closure.

**Beide INFO-Findings sind Closure-Notizen, keine Blocker.** Kein HIGH, kein
MEDIUM in der Fixrunde; mein eigenes Verdikt deckt sich mit dem des Reviewers
(„keine Fixrunde mehr nötig").

## 8. `welle-19`-Stand — was nach dieser Closure noch fehlt

`ls docs/plan/planning/done/ | grep -E "slice-069"` → **kein Treffer**;
`slice-069` liegt real in `in-progress/`. `ls docs/plan/planning/done/ | grep -E
"slice-07[0-2]"` → **keine** Treffer: `slice-070`, `slice-071`, `slice-072`
liegen noch nicht in `done/`.

Der Closure-Trigger (`welle-19.md` §3) hat fünf Bedingungen; ich prüfe jede
gegen den Stand:

| Bedingung | Stand meiner Prüfung | Erfüllt? |
|---|---|---|
| Alle vier Slices (`slice-069`…`slice-072`) in `done/` | `slice-069` in `in-progress/`; `slice-070`/`071`/`072` nicht geschlossen | **nein** |
| `make gates` grün | eigener Lauf Exit **0** | **ja** |
| Realer gRPC-E2E-Rundlauf (`slice-071`, `make test-integration`) | `slice-071` nicht begonnen; kein Beleg | **nein** |
| Realer SSE-E2E-Rundlauf (`slice-072`, `make test-integration`) | `slice-072` nicht begonnen; kein Beleg | **nein** |
| Closure-Notiz in `welle-19-results.md` | Datei existiert real nicht; `welle-19` §7 noch Platzhalter | **nein** |

**Ergebnis §8:** Nach der `slice-069`-Closure fehlen zur Welle noch **drei
Slices** (`slice-070` Capture-Integration, `slice-071` gRPC-Beispiel-Client/E2E,
`slice-072` HTTP/SSE-Endpunkt) sowie die beiden realen E2E-Belege (gRPC **und**
SSE) und die Wellen-Closure-Notiz — das *Mehr* der Welle, das keine Einzel-DoD
für sich beansprucht. `make gates` ist bereits grün, `slice-069` selbst ist der
erste von vier Schritten. Nebenbei: `welle-19.md:10` führt `Verantwortlich: —`
(deklarativ, kein Sensor — Modul 5 §„Ein Sensor prüft das Feld nicht"); kein
Verifikationsgegenstand.

## Negativbefunde

- **geprüft, ohne Befund: DoD-Punkt-Trennung und -Vollständigkeit.** Zwölf
  Punkte, sieben abgehakt (Implementierung, Fixrunde, `SPEC-020`, Gate, Review,
  Doku, Reconciliation-Entfall), fünf offen (Closure-Notiz, Register,
  Risiko-Ausgänge, Paarungen) — keine vorweggenommene Closure-Zeile. Einzige
  Rest-Beobachtung: der „kein Puffer"-Wortlaut in DoD-Zeile 2 (§1).
- **geprüft, ohne Befund: §5-Closure-Trigger des Slice.** „DoD vollständig `+`
  PR gemerged `+` `make gates` grün `+` Closure-Notiz" — die Sensor-Hälfte ist
  grün, die Closure-Notiz ist die noch fehlende Planner-Leistung.
- **geprüft, ohne Befund: §3-Plan-Nachzüge (kein Scope-Creep).** Kein
  Liefer-Punkt-Zuwachs; die Fixrunde bleibt innerhalb des Liefer-Punkts
  „Adapter-Grundgerüst". Der einzige nicht im Plan-Platzhalter genannte Pfad
  (`streamv1/`) ist als Implementer-Entscheidung geführt (`:172-175`).
- **geprüft, ohne Befund: `ADR-0060` unverändert.** `git diff --name-only
  4085f18..HEAD -- docs/plan/adr/` nennt nur `0066-…md` (neu) `+` `README.md`
  (Index-Zeile); `ADR-0060`s Textkörper ist unberührt. `make doc-immutable
  RANGE=4085f18..HEAD` 0 Befunde. Der `Supersedes`-Vermerk steht am Index
  (`docs/plan/adr/README.md:73,79`), die Klausel-Supersession ist eng gefasst.
- **geprüft, ohne Befund: Capture-Integration bleibt ausgeschlossen.** `grep`
  über den Produktionscode: `grpcstream.New()` steht **nur** im
  `apigrpc.Config` (`wiring.go:691`); kein `WithChangeStream`, kein
  `ChangeStreamPort` am `CaptureService`; `internal/application/usecase/**` ist
  im Diff nicht enthalten.
- **geprüft, ohne Befund: Schichtung und Importe.** Der `grpc`-Adapter
  importiert **kein** Adapter-Paket; `make a-check` 0 Befunde. Das lokale
  `changeSubscriber`-Interface wird strukturell von `*grpcstream.Broadcaster`
  erfüllt; die Zuordnung prüft der Compiler an der Verdrahtungsstelle.
- **geprüft, ohne Befund: keine Chronik in Produktionscode-Kommentaren.** Die
  neuen Godocs (`broadcaster.go`, `changestream.go`, `interceptor.go`,
  `middleware.go`) nennen Zusage/Kopplung/Abgrenzung und tragen `ADR-*`-Zeiger;
  der einzige `slice-`-Treffer ist die **vorbestehende**, von diesem Diff
  geänderte `slice-059`-Chronik in `middleware.go`, die die Fixrunde **entfernt**
  hat (§6).
- **geprüft, ohne Befund: kein Suppression-Pfad.** Kein `nolint`/`noqa`/
  `SuppressMessage` in den neuen Paketen und dem Port (`AGENTS.md` §3.2).
- **geprüft, ohne Befund: Traceability.** `make doc-commits RANGE=4085f18..HEAD`
  Exit 0; die drei Fixrunden-Commits tragen `LH-FA-SST-008`/`ADR-0066` im
  Betreff, kein `SPEC-*`/`ARC-*`.
- **geprüft, ohne Befund: Proto-Toolchain bleibt Docker-only.** `protoc` und die
  beiden Plugins kommen ausschließlich aus der `Dockerfile`-Stufe `proto`; kein
  Host-`protoc`/`buf` (`AGENTS.md` §3.1). Die Stufe steht vor
  `coverage`/`build`/`runtime`, der Default-Build trifft weiterhin `runtime`.

## Summary

| Prüfung | Ergebnis |
|---|---|
| DoD-Punkte 1–7 (implementierungs-/reviewbezogen) | **erfüllt**, 1/3/5 selbst reproduziert |
| DoD-Punkte 8–12 (Closure) | korrekt offen, keine vorweggenommen |
| Fünf Sensoren (`gates`, `test`, `a-check`, `doc-commits`, `doc-immutable`) | **je Exit 0** |
| Mutation A (Sende-Zweig blockierend) | **rot** an beiden Kern-Wächtern (2.00s-Frist), Rücknahme + Blatt-Identität verifiziert |
| Mutation C (nur Frühausstieg) / B (Paar) | **grün** (Exit 0) / **rot** (Exit 2) — F-10 bestätigt |
| Wächter pulsiert über die Kapazität | **bestätigt** — `queueCapacity+1`, strukturell an die Konstante gekoppelt |
| `ADR-0066` gegen den Code (4 Festlegungen) | **deckungsgleich** |
| `SPEC-020`-Zeilen `+` Port-Godoc (inkl. beider Fehlerausgänge) | **deckungsgleich**; eine ADR-Textspanne benannt (§5) |
| Unverändert: `natsnotify`/`ChangeNotificationPort`/`.a-check.yml`/`http` | **bestätigt** |
| DoD-Häkchen-Trennung | korrekt |
| F-9 / F-10 | **Closure-Notizen, keine Blocker** |
| `welle-19`-Closure-Reife | **noch nicht schließbar** — drei Slices offen, zwei E2E-Belege fehlen |

**Verdikt:** `slice-069` erfüllt seine DoD. Die Kernzusage aus
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(Erzeuger hält nie an, begrenzte Empfangswarteschlange, Drop-Newest, Kanal nie
geschlossen) ist am Code deckungsgleich umgesetzt und durch eine **eigene**
Mutation als **rot** belegt — an genau der zugesagten Stelle, mit einem Wächter,
der strukturell über die Kapazität hinaus pulsiert und deshalb nicht blind ist.
Die Leer-Token-Grenze ist **paartragend**; das habe ich unabhängig bestätigt
(nur der Frühausstieg: grün; das Paar: rot in beiden Testhälften). Die
`SPEC-020`-Zeilen und der Port-Godoc decken die Semantik **vollständig**,
inklusive beider Fehlerausgänge von `Publish`. Die Unverändert-Zusagen
(`natsnotify`/`ChangeNotificationPort`, `.a-check.yml`, HTTP-Adapter außer dem
Gegenverweis-Kommentar) halten einem eigenen Diff-Vergleich stand. Alle fünf
Sensoren, die ich selbst ausgeführt habe, laufen Exit 0.

**Zwei Rest-Beobachtungen ohne Zwang** (kein Implementer-Rückweg): der
„kein Puffer"-Wortlaut in DoD-Zeile 2 ist durch Zeile 3/`ADR-0066` überholt und
kann beim Closure-Nachzug angeglichen werden; die Klammer in `ADR-0066`
„(nur … `ErrChangeStream`)" ist enger als der Code (der beendete `ctx` liefert
den rohen Kontext-Fehler), was `SPEC-020` und der Godoc bereits korrekt pinnen.

**Übergabe an den Planner (Verifier → Planner):** Der Slice ist DoD-konform;
über diese Verifikation hinaus ist **kein** Implementer-Rückweg nötig. Für den
`slice-069`-Closure-Schritt bleiben: Closure-Notiz mit Steering-Loop-Eintrag,
die drei §6-Risiko-Ausgänge (die noch „wird bei Closure zugewiesen" tragen), die
Beobachtungs-Register-Entscheidung (Sichtung nennt Treffer; kein `slice-069`-
Beleg), die F-9/F-10-Klassen im Register und der `git mv` nach `done/`. Für die
danach folgende `welle-19`-Closure: mein Lauf dieses Berichts ist der repo-weite
Verifikations-Beleg aus Closure-Schritt 1 (Modul 6/8); die Welle ist erst
schließbar, wenn `slice-070`…`slice-072` in `done/` liegen und beide realen
E2E-Belege (gRPC über `slice-071`, SSE über `slice-072`) grün vorliegen.
