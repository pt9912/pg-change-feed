# Review-Report: slice-091 — **Delta** (Fixrunden `f9cd5e4` und `a7d7f7b`) · 2026-09-16

**Review-Art:** Delta-Review **zweier** Fix-Commits — geprüft gegen **Plan und
Entscheidungen** sowie gegen die Findings der ersten Runde und die Adressen des
Verifikationsberichts (Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-10-review-harness.md` §Drei Review-Arten). DoD-/Spec-Konformität
(Verifier), §6-Risiko-Ausgänge, Beobachtungs-Register und die drei Paarungen
(Planner-Closure) sind **nicht** Gegenstand.

**Gegenstand** (Stand `HEAD` = `a7d7f7b`, Arbeitsbaum vor **und** nach jedem Lauf
sauber):

| Commit | Umfang (gemessen) |
|---|---|
| `f9cd5e4` | 4 Dateien, +64/−23: `natsnotify/notify_test.go`, `grpc/interceptor_test.go`, `grpc/server_test.go`, `http/sse_test.go` |
| `a7d7f7b` | 2 Dateien, +45/−15: `harness/sensors/coverage-gate.md`, `http/server_test.go` |
| Vorgang `f90c3f4^..a7d7f7b` | **12 Pfade**: 9 `*_test.go`, 2 Review-Berichte, 1 Sensor-Doku — **kein** Produktcode, kein `git mv`, `THRESHOLD`/`coverage.mk` unberührt |

**Skill:** `.harness/skills/reviewer.md` @ `HEAD` — die vier repo-spezifischen
HIGH-Unterpunkte (u. a. „Zahl im Träger", „Zusage ohne Bindung an ihre
Eingabeseite") gehören zum Prüfraster · **Datum:** 2026-09-16.

**Eingangs-Kontext:** das Review zu `slice-091` (F-1…F-4) · der Verifikationsbericht zu `slice-091`
(V-1…V-8) · Slice-Plan `slice-091-coverage-cluster-c` §1/§3 (Zeile 148:
`harness/sensors/coverage-gate.md` „update, **nur falls** eine Zahl dort gegen
die Messung driftet") · `ADR-0082` §Kontext (2), §Konsequenzen, §Fitness
Function `:350`, §Re-Evaluierungs-Trigger (a)–(e) · `ADR-0071`, `ADR-0077`,
`ADR-0055`, `ADR-0057`/`ADR-0060`, `ADR-0024` · `AGENTS.md` §3.1, §3.3, §3.5,
§3.6, §3.7, §3.9, §3.11, §3.12 · `harness/conventions.md` (MR-000) ·
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` · Nachbar-Delta,
der Delta-Review zu `slice-090`.

---

## Findings

### D-1 — Die neue `streamv1`-Begründung trägt ihre Gruppe nicht; die Nachbarsätze desselben Aufzählungspunkts sind mit diesem Slice falsch geworden

- `kategorie`: **MEDIUM**
- `quelle`: `AGENTS.md` §3.12 Instanz B (Tatsachenbehauptung ohne tragenden
  Beleg-Anker) · Source Precedence (der Lauf ist der Beleg, nicht die Zählung) ·
  Klasse `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
- `pfad`: `harness/sensors/coverage-gate.md:118-125` (neu) gegen `:110-112`
  (Überschrift) und `:131-133` (Schlusssatz)
- `befund`: Der Zusatz begründet die Zugehörigkeit von `streamv1` zu den „**Fünf
  Pakete[n] des Gegenstands [ohne] Testdatei" mit der Zählung
  `go list -f '{{len .TestGoFiles}}'`. Gemessen liefert diese Zählung über den
  Gegenstand **23** Pakete mit `0` (darunter `natsnotify` 0/1, `port/outbound`
  0/3 und **13** `usecase/*` 0/1) — sie ist damit keine Gruppierungsregel,
  sondern für **jedes** Paket dieses Repos mit externem Testpaket wahr; die fünf
  genannten Pakete sind umgekehrt unter **keiner** der beiden mechanischen
  Lesarten reproduzierbar („keine Testdatei" = 4 Pakete `0/0`; „kein
  Test-Binary im Tor" = dieselben 4 zzgl. `cmd/pg-change-feed`). Der neue Satz
  *ist* unter der zitierten Zählung wörtlich wahr — er macht aber genau die
  Zählung tragend, die die Liste nicht erzeugt. Dazu zwei Sätze desselben
  Aufzählungspunkts, die dieser Slice falsch gemacht hat und die der Fix stehen
  ließ: die Überschrift `:110` („haben keine Testdatei" — `streamv1` hat seit
  `f90c3f4` eine, die Datei dieses Vorgangs) und `:131-133` („Die beiden Pakete
  mit Statements weist der Lauf … als `coverage: 0.0% of statements` aus" — der
  Lauf weist für `streamv1` `ok … coverage: 2.4%` aus, die `0.0%`-Zeile trägt
  heute nur noch `cmd/pg-change-feed`).
- `verifizierbar`: ja — `go list -f '{{.ImportPath}} Test={{len .TestGoFiles}}
  XTest={{len .XTestGoFiles}}'` über den Gegenstand; die Ausgabezeilen des
  `coverage`-Laufs; `git log --diff-filter=A --
  internal/adapters/driving/grpc/streamv1/changestream_test.go` → `f90c3f4`
- `klasse`: „Gruppenzugehörigkeit mit einer Zählung begründet, die die Gruppe
  nicht erzeugt" — **derselbe Vorgang** wie F-1/F-2 der ersten Runde (Modul 6:
  ein Vorgang zählt einmal; kein neuer Zählerstand)

### D-2 — Das Band `1468 … 1471` ist in 54 eigenen Läufen nicht reproduzierbar (Mechanik dagegen bestätigt); zwei Stände stehen ohne Stand-Marke im selben Absatz

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (die Zahl trägt ihren Lauf; die Probe ist
  das Nachmessen)
- `pfad`: `harness/sensors/coverage-gate.md:68-74` (neu), benutzt in `:151-153`
  (geändert)
- `befund`: **Trägt:** die zwei Blöcke und ihre Zuordnung — 24 isolierte
  `internal/bootstrap`-Läufe zeigen als **einzige** unruhige Blöcke
  `wiring.go:991.5,992.13` (2 Statements) und `:1091.4,1092.1` (1 Statement);
  der zweite fiel real auf `count = 0` (1 von 24) — das „dritte Ende" der
  Verifikation (V-1: 435 ungedeckt) ist damit arithmetisch **exakt** erklärt:
  denselben Lauf synthetisch um diesen Block bereinigt druckt `go tool cover
  -func` **77.1 %** statt 77.2 %. **Trägt (bei mir) nicht:** das genannte Band als
  Acht-Läufe-Befund — in **54** eigenen Läufen der vollen Fläche (48 sequenziell
  + 6 **parallel unter Last**) lag die gedeckte Zahl ausschließlich bei **1469**
  (41×, gedruckt `77.2%`) oder **1471** (13×, `77.3%`); der
  `runAdministration`-Block fiel in **keinem** dieser 54 Läufe auf `0`. Die
  Aussage ist damit **nicht widerlegt** (ihr unteres Ende ist erreichbar und ihr
  Lauf ist ihrer), aber sie beschreibt einen ~1 %-Fall als Beobachtung eines
  Acht-Läufe-Fensters; mit der Verifikation (1 von 47) sind es **1 von 101**
  Läufen der vollen Fläche. Zweiter Punkt: der Absatz stellt nun zwei Bänder
  **desselben Nenners 1903** nebeneinander — 1369/1371 (Lauf `slice-084`/`-085`)
  und 1468/1471 (Lauf `slice-091`) — beide mit „demselben Stand"/„demselben
  Produktionsstands" eingeleitet; nur die Lauf-Marker unterscheiden sie. **Kein
  Reparatur-Ort in diesem Diff, kein Blocker:** die Verwendung in `:151-153` ist
  die konservative Richtung (breiteres Band = größere Marge).
- `verifizierbar`: ja — 8/24/54 Läufe mit derselben Zählbasis nachfahren und die
  zwei Block-Positionen je Lauf auswerten; ein Profil mit dem
  `:1091.4,1092.1`-Vorkommen auf `0` durch `go tool cover -func` geben
- `klasse`: „Lauf-Angabe im Träger — Beleg nicht reproduzierbar, Substanz
  bestätigt"

### D-3 — „zwölf Zeilen darüber" in der Commit-Message benennt keine auflösende Stelle

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 §Geltungsbereich (git-Historie ist **kein**
  Doku-Träger, `ADR-0083`) — deshalb INFO und **ohne** erwartete Aktion
- `pfad`: Commit `a7d7f7b`, Message-Absatz „Drei Doku-Zeilen statt zwei"
- `befund`: Der Satz sagt, §Grenze Punkt 4 habe „nach der Korrektur zwölf Zeilen
  **darüber**" widersprochen. Gemessen liegt die geänderte Stelle in `:152`, die
  Band-Aussage in `:68-74` — Abstand **78 Zeilen**; zwölf Zeilen *über* `:152`
  steht die Überschrift von §Grenze Punkt 4 selbst (`:140`), nicht die
  Band-Aussage. Die **Substanz** des Satzes trägt (der alte Maßstab „±2"
  widersprach tatsächlich der korrigierten Band-Angabe); falsch ist nur der
  Anker.
- `verifizierbar`: ja — `grep -n` über `harness/sensors/coverage-gate.md`
- `klasse`: „Abstands-Angabe über den eigenen Diff ohne auflösende Stelle"

---

## Negativbefunde

Jeder Block mit dem gefahrenen Kommando; Exit-Codes ungepiped, je in eigenem
Schritt gelesen (§3.9). Mutationsproben liefen auf einer Arbeitsbaum-Kopie
**außerhalb** des Repos (`git archive HEAD` → `/tmp/pgc-mut`); die Kopie ist nach
den Proben byte-gleich zum Stand (`diff -r` leer).

**`internal/adapters/driven/natsnotify/` — F-1 trägt, gemessen: der Test ist jetzt gebunden, nicht der Kommentar zurückgenommen.**

| Probe | Mutationsstelle | Ergebnis |
|---|---|---|
| Kontrolle | — | `ok` **EC=0** |
| M-A (im Kopf genannt) | `notify.go:81` `log: o.log` → `outbound.NoopLog` | **EC=1** — `notify_test.go:224: Aufzeichnung des Adapterfelds: [] …` |
| M-B (im Kopf genannt) | `newOptions` Options-Schleife entfernt | **EC=1** — `notify_test.go:216: Aufzeichnung des Konstruktionsaufrufs: [] …` |
| M-C (eigene Zusatzprobe, drittes Kettenglied) | `notify.go:134` `notifyFailure(ctx, a.log, err)` → `outbound.NoopLog` | **EC=1** — dieselbe Zeile `:224` |

Das Paket als Ganzes unter M-A: **EC=1**. Die im Kommentar behauptete Wirkung
(„dann fehlt die Fehler-Aufzeichnung und dieser Test färbt rot") tritt ein; die
Kette ist über **drei** Glieder gebunden (Options → Adapterfeld → Aufrufseite).
Der `recordingLog`-Zusatz (`errors`) bricht keinen anderen Test des Pakets. Keine
`t.Parallel`-Nutzung, kein geteilter Zustand neu.

**`internal/adapters/driving/grpc/` — F-2 trägt in allen drei genannten Formen.**

| Probe | Mutationsstelle | Ergebnis |
|---|---|---|
| Kontrolle | — | `ok` **EC=0** |
| gültiger Wert (als rot benannt) | `interceptor.go:71` → `"reader-token"` | **EC=1** — `interceptor_test.go:75: Status ohne eingehende Metadata: OK (Erwartung: Unauthenticated)` |
| unbekannter Wert (als grün benannt) | `interceptor.go:71` → `"probe"` | **EC=0** |
| `!ok`-Zweig entfernt (als grün benannt) | `md, _ := metadata.FromIncomingContext(ctx)`, `if !ok` gestrichen | **EC=0** |

Genau die Verteilung, die der neue Kommentar behauptet — die als „nicht rot
färbend" benannten Formen färben real nicht rot.

**`internal/adapters/driving/http/` — W-3 trägt, gemessen (und war vorher grün).**

| Probe | Mutationsstelle | Ergebnis |
|---|---|---|
| Kontrolle | — | `ok` **EC=0** |
| Decode-Zweig abgeschaltet | `registerconsumer.go:38` `… ; false && err != nil {` | **EC=1** — `server_test.go:177: Status: 500 (Erwartung: 400 für einen nicht dekodierbaren Body)` |
| dasselbe, ganzes Paket | wie oben | **EC=1** |

Damit ist die von V-7 gemessene Lücke (derselbe Test blieb vorher **grün**)
geschlossen. Der Gegenproben-Fake liefert für einen gültigen Body real `500`
(`registerconsumer.go:46-53` kennt nur `ErrEmptyIdentifier` als `400`);
`newFakeRegisterConsumerUseCase` bleibt andernorts benutzt (kein toter Helfer).

**Die vier Fristen (F-3) — Zurückhaltung belegt, Marge gemessen.**
`3 * time.Second` je Datei und Stand, `grep -c`: `http/sse_test.go` 6 → **9**
(`f90c3f4`) → **6** (`f9cd5e4`) → 6; `grpc/server_test.go` 3 → **4** → **3** → 3.
Die vier **neuen** Fristen liegen jetzt auf `haengeFrist = 30 * time.Second` (je
**eine** Deklaration pro Paket, `:356` bzw. `:287`; vier Verwendungen); die
**neun** vorbestehenden 3-s-Fristen sind unangetastet. `go test -race -count=20`
über die vier Pakete: **EC=0**, keine Frist gefeuert (`http` 5,588 s · `grpc`
1,269 s · `streamv1` 1,018 s · `natsnotify` 1,040 s für alle zwanzig
Iterationen) — die Bezugsgröße der neuen Kommentare („in wenigen Sekunden")
hält. Die Klasse aus F-3 bleibt der **Form** nach bestehen (die Frist ist in drei
Fällen zugleich die Beobachtung „der Handler endete"); sie ist mit 30 s gegen
≈0,29 s je Iteration um zwei Größenordnungen abgesetzt und im Kommentar benannt —
**kein neuer Befund über F-3 hinaus**, und die vorbestehenden Stellen sind nicht
Gegenstand des Diffs.

**`harness/sensors/coverage-gate.md` — W-1: die falsche Ursache ist ersetzt, nicht umgedreht.**
Kein Rest der alten Zuordnung (`grep -n "broadcaster\|±\|schwank" …` → nur die
neue „2 Statements"-Zeile `:65`); keine Vorher/Nachher-Sprache, kein abwesender
Text (§3.7) — der Kommentar der Message („sie wird weggenommen, nicht als
‚früher' geführt") ist im Text eingelöst. Alle **neuen** Zahlen sind
nachgemessen: Block-Positionen und Statement-Zahlen aus dem eigenen Profil
(`:991.5,992.13` = **2**, `:1091.4,1092.1` = **1**), die Bandbreite arithmetisch
konsistent (`1468 → 77.1 %`, `1471 → 77.3 %`), `3/1903 = 0,16 pp` korrekt,
`§Zählbasis` und `ADR-0082 §Kontext (2)` lösen auf, und die zitierte Wendung
(„der Takt-Zweig feuert nur, wenn der Tick vor dem Kontext-Ende liegt") **ist**
der Wortlaut der ADR. Der Mechanismus selbst ist gegen den Code geprüft:
`fakeWALRetentionMeasurer.Measure` gibt nach der Wertfolge `ctx.Err()` zurück
(`walretention_internal_test.go:112-125`) — der Block ist der Messfehler-Zweig
**innerhalb** des Takt-Zweigs, genau der Wettlauf, den die ADR beschreibt. Die
übrigen drei Stellen mit der alten Ursache sind **Records**
(das Review zu `slice-081`, `done/slice-081-executor-naht.md`) und damit
außerhalb der Träger-Pflege; **kein stehender Träger** außer dieser Datei führt
sie.

**`ADR-0082` wird von der Korrektur nicht verletzt.** Ihre Aussagen („genau
**eine** Funktion", „±2") sind ausdrücklich auf *ihre zwei Läufe* bezogen
(`§Kontext (2)`; `§Fitness Function :350` „zum Stand dieser ADR"). Keiner der
fünf Re-Evaluierungs-Trigger (a)–(e) wird durch die Korrektur ausgelöst; die
13-Statements-Reserve auf 80 % bleibt bei einem 3-Statements-Band erhalten. Die
neue Messung **erweitert** das Bild (zweiter, seltener Träger) — eine Erweiterung
braucht keinen `supersedes` (§3.5 unberührt, keine ADR angefasst).

**Plan-Konformität des Sensor-Doku-Zugs.** §3 Zeile 148 erlaubt „update, nur
falls eine Zahl dort gegen die Messung driftet" — der Fix zieht **drei** Stellen
nach (Zählbasis-Ursache, `streamv1`-Bild, §Grenze-4-Maßstab); alle drei waren
adressiert (V-2, V-6) bzw. durch die Korrektur selbst widersprüchlich geworden.
Kein `THRESHOLD`, kein Gate, keine Schwelle berührt (§3.6 unberührt), kein
`git mv` (§3.3), keine Betreiber-Oberfläche und `docs/user/` unberührt, keine
neue Kennung (MR-000 unberührt).

**§3.7/§3.11 in den `+`-Zeilen beider Commits.** Keine Slice-/Wellen-Nummer als
Begründung (Treffer nur als **Lauf-Marker** „Lauf `slice-091`" — die Hausform der
Datei), kein host-lokaler absoluter Pfad, keine abwesende Sprache, keine mitten
im Absatz gerissenen Kommentare. `make docs-check` grün.

**Die Mutationsangaben der beiden neuen Testkommentare halten wörtlich** (F-1,
F-2, W-3 — Tabellen oben); die Aussage „der Status folgt dem Body, nicht dem
Fake" ist über die Gegenprobe mitgemessen (gültiger Body → `500`, nicht
dekodierbarer Body → `400`, Mutation dreht die zweite auf `500`).

**Hygiene und Gates.** `make gates`: **EC=0**: `baseline-verify v6.5.0 OK — 54
Dateien` · `coverage-gate: OK — Coverage 77.20% erfüllt Schwelle 70%` ·
`d-check: 749 Datei(en), 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)`,
Betreffs ohne Struktur-ID · `generated-sync: OK` · `a-check: gesamt: 0
Befund(e)`; danach `git status --porcelain` leer. Beide Betreffe nennen
`ADR-0082`; der Vorgang berührt **keinen** Produktionscode, die `test-only`-
Folgepflicht der ADR hält.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-status f90c3f4^..a7d7f7b` | **0** | 12 Pfade: 9 `*_test.go`, 2 Berichte, 1 Sensor-Doku |
| Profil der Gegenstands-Fläche, **54 Läufe** (48 sequenziell + 6 parallel) | **0** (je) | gedeckt **1469** (41×) / **1471** (13×); gedruckt `77.2%`/`77.3%`; Nenner 1903 |
| dieselben 54, Block-Auswertung | **0** | einziger unruhiger Block `wiring.go:991.5,992.13` (2 Statements, `count>0` in 13/54); `:1091.4,1092.1` in **0/54** auf `0` |
| `internal/bootstrap` isoliert, **24 Läufe** | **0** (je) | genau **zwei** unruhige Blöcke: `:991.5,992.13` (8/24 gedeckt), `:1091.4,1092.1` (**1/24** auf `0`) |
| `go tool cover -func` auf synthetischem Profil (Admin-Block auf `0`) | **0** | `77.1%` statt `77.2%` — das dritte Ende ist genau dieser Block |
| `go list -f '{{len .TestGoFiles}} {{len .XTestGoFiles}}'` über den Gegenstand (31 Pakete) | **0** | `streamv1` = **0/1** ✓; **23 von 31** Paketen haben `TestGoFiles=0` → **D-1** |
| Profil-Auswertung `streamv1` | **0** | **86** Statements, **76** gedeckt (88,37 % → „88,4 %") ✓ |
| Lauf-Ausgabe (`coverage`-Stufe) | **0** | `streamv1`: `ok … coverage: 2.4% …`; `cmd/pg-change-feed`: `coverage: 0.0% of statements` → **D-1** |
| Mutationsproben (6, Baum-Kopie außerhalb) | **1/0** | F-1: 3×**rot**, Kontrolle 0; F-2: **rot/green/green** wie behauptet; W-3: **rot**, Kontrolle 0 |
| `go test -race -count=20` (vier Pakete) | **0** | 4 × `ok`, keine Frist gefeuert |
| `grep -c "3 \* time.Second"` je Stand | **0** | `sse` 6/9/6/6 · `grpc` 3/4/3/3 — nur die vier neuen gehoben |
| `make gates` (Log in Datei, EC separat) | **0** | sechs Checks grün; danach Status leer |
| `diff -r` Kopie ↔ Baum | **0** | leer (alle Mutationsproben zurückgesetzt) |

Der Gate-Lauf und seine Auswertung sind **zwei Schritte**: `make gates` lief
ungepiped in eine Log-Datei, der Exit-Code in eine eigene Datei, gelesen in einem
eigenen Werkzeug-Aufruf; kein Folgekommando hing an Pipe oder Wrapper (§3.9).

## Antwort auf die Schwerpunkte

**(1) W-1 — die gefährlichste Stelle: trägt die Ersatz-Ursache?** **Ja, gemessen
— und besser belegt als die ersetzte.** Die zwei Blöcke existieren mit den
genannten Statement-Zahlen, sind in **24** isolierten Läufen die **einzigen**
unruhigen Blöcke des Pakets, und der zweite erklärt das von V-1 nicht
lokalisierte dritte Ende **arithmetisch exakt** (`1469 − 1 = 1468 → 77.1 %`). Die
*Quantifizierung* („Band 1468–1471 über acht Läufe, der zweite Block fiel in
einem davon auf `0`") konnte ich mit **54** eigenen Läufen **nicht**
reproduzieren (mein Band: 1469–1471; der zweite Block fiel in 0/54, dagegen 1/24
isoliert) — **D-2, INFO**: die Aussage ist nicht widerlegt, aber ihre untere
Kante ist ein seltener Fall, und mein Nachmessen trägt nur die *Mechanik*, nicht
die behauptete Acht-Läufe-Beobachtung. Kein Fall von „falsche Ursache durch
andere falsche Ursache": die neue Zuordnung hält der eigenen Messung stand.

**(2) Hat die Korrektur etwas Neues falsch gemacht?** **Ein Punkt: ja.**
`grep`-fest, gemessen und mit Zeile: **D-1** (die neue Begründung der
`streamv1`-Gruppe über eine Zählung, die 23 Pakete trifft, während der
Aufzählungspunkt „fünf" sagt — und zwei Nachbarsätze desselben Punkts, die dieser
Slice falsch gemacht hat). Alles andere Neue hält: keine ungemessene Zahl in der
Sensor-Doku, kein Chronik-Ton, kein abwesender Text; in den vier Testdateien
halten die Mutationsangaben wörtlich, die beiden neuen `haengeFrist`-Kommentare
beschreiben den Code und ihre Bezugsgröße ist gemessen. **D-3** (INFO) ist eine
falsche Abstandsangabe ausschließlich in der Commit-Message — außerhalb des
Geltungsbereichs von §3.12, deshalb keine Aktion.

**(3) F-1 und W-3 per Mutation nachgeprüft.** Beide Behauptungen halten:
`notify.go:81` → **EC=1** (und zusätzlich das dritte Kettenglied `notify.go:134`
→ EC=1); `registerconsumer.go:38` `false && err != nil` → **EC=1**
(`server_test.go:177`, „500 statt 400"), ganzes `http`-Paket → EC=1. Kontrollen
bei unmutiertem Stand je **EC=0**.

**(4) F-3 — nur die vier neuen Fristen gehoben.** Bestätigt: die neun
vorbestehenden 3-s-Fristen sind unangetastet (`grep -c` vor/nach dem Fix-Commit
identisch 6 und 3), die vier neuen stehen auf `haengeFrist` = 30 s. Die
Begründung trägt in der Sache (Marge ≈100× gegenüber ≈0,29 s je Iteration,
`-race -count=20` grün); die vorbestehenden Stellen sind nicht Gegenstand dieses
Diffs, und die Reviewer-Regel „kein Refactoring-Vorschlag über den Diff hinaus"
schließt einen Angleichungsvorschlag aus — die Asymmetrie (3 s und 30 s für
dieselbe Bauform in einer Datei) ist Folge der Schnitt-Disziplin, nicht Fund.

**(5) W-2 nachgerechnet.** **Zahlen halten:** 86 Statements, 76 gedeckt = 88,37 %
(„88,4 %"), `XTestGoFiles` = 1 ✓, ein eigenes externes Testpaket
(`package streamv1_test`) existiert ✓. **Die Begründung hält nicht:** die Zählung
`{{len .TestGoFiles}}` trifft 23 Pakete des Gegenstands, nicht die genannten
fünf, und die zwei Nachbarsätze des Punkts sind mit diesem Slice falsch geworden
→ **D-1 (MEDIUM)**.

**Was nicht geprüft wurde:** DoD/LP1–LP3 (Verifier), die §6-Risiko-Ausgänge, die
Register-Zählung und die drei Paarungen (Planner-Closure), die Frage, ob V-7s
„kein Reparatur-Ort" durch die nachgezogene Bindung berührt wird (Modul 6,
Lese-Schritt), der reale Post-Push-Lauf (§3.10: dieser Vorgang ändert keinen
Workflow), und ob die Frist-Marge im kalten CI-Lauf trägt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **0** |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Gruppenzugehörigkeit mit einer Zählung
begründet, die die Gruppe nicht erzeugt" (D-1) · „Lauf-Angabe im Träger — Beleg
nicht reproduzierbar, Substanz bestätigt" (D-2) · „Abstands-Angabe über den
eigenen Diff ohne auflösende Stelle" (D-3).

**Register-Hinweis (Modul 6):** D-1 fällt in **dieselbe Klasse** wie F-1/F-2 der
ersten Runde (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`) und in
**denselben Vorgang** (`slice-091`) — der Zähler bewegt sich durch diesen Delta
**nicht** (ein Vorgang zählt einmal). D-2 berührt sachlich
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (dort seit `slice-089`
verkörpert in `AGENTS.md` §3.12); ob das ein Auftreten oder eine zweite
Beobachtung ist, entscheidet der Lese-Schritt der `welle-20`-Closure.

## Verdikt

**Merge-blockierend:** ja — **1 MEDIUM** (D-1). Kein HIGH.

**Tragen die zwei Fixrunden?** **Ja, in der Substanz — mit einer Korrektur an
W-2.**

- **`f9cd5e4` (F-1, F-2, F-3) trägt vollständig, gemessen.** F-1: der Test ist
  real **gebunden** (drei Kettenglieder, drei rote Proben); F-2: die drei
  genannten Mutationsformen verhalten sich genau wie beschrieben; F-3: nur die
  vier **neuen** Fristen gehoben, 30 s, `-race -count=20` grün ohne
  Frist-Feuerung. **V-4 ist damit erledigt** — dieser Report deckt den Stand
  `f9cd5e4` mit ab.
- **`a7d7f7b` trägt in drei von vier Aussagen.** **W-1 trägt** (die Ursache ist
  gemessen ersetzt — die einzige mir bekannte Stelle in diesem Repo, an der das
  gelungen ist; der zweite Block erklärt V-1s drittes Ende exakt); **W-3 trägt**
  (Mutation rot, Kontrolle grün — vorher grün, V-7 geschlossen); **W-2s Zahlen
  tragen**, seine **Begründung der Gruppenzugehörigkeit nicht** (D-1).

**Dritte Runde: ja, aber eine sehr kurze.** Für D-1 ist ein Rückgabe-Pfeil
**Reviewer → Implementer** nötig (MEDIUM; Ein-Stelle-Korrektur in
`harness/sensors/coverage-gate.md`: die Gruppe auf ihren mechanischen Bestand
bringen — vier Pakete ohne jede Testdatei plus `streamv1` in einer eigenen Rolle
— und den `0.0%`-Schlusssatz auf `cmd/pg-change-feed` begrenzen). D-2 und D-3
sind INFO **ohne** Rückgabe-Pfeil.

**DoD-Häkchen „Review durchgeführt, Report unter `docs/reviews/` liegt vor":**
bleibt **offen** — der Slice braucht eine (kurze) Fixrunde
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde greift nicht).

**Übergabe:** Rückgabe-Pfeil Reviewer → Implementer mit **D-1**; die
Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser
Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation.
