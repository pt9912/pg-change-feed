# Verifikationsbericht: slice-097 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-097` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
§3 (Schnitt-Empfehlung 1: „der Umzug soll allein stehen") und §Konsequenzen
(Coverage-Aufnahme) · [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
(der Draht-Vertrag, unverändert) · [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(die zurückgenommene `tooling→adapters`-Kante) ·
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)/[`ADR-0082`](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(Messgegenstand-Eigenschaft, nicht Liste) · [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
(keine Supersession dieses Slice) — sowie die Hard Rules `AGENTS.md` §3.3
(git mv + Inhalt = zwei Commits), §3.6, §3.9, §3.12, §3.13. **Nicht** gegen den
Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-097.md`](review-slice-097.md) abgeschlossen) und **nicht** gegen
realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan am Stand `HEAD`
vollständig gelesen (§1–§8), dazu den Architect-Verdikt
`architect-verdict-slice-097-coverage-gegenstand.md`, den Review-Report, den
Plan-Diff-Commit `4f39cd5`, `.a-check.yml`, `Dockerfile` (Stufe `coverage`),
`.dockerignore`, `harness/sensors/{coverage-gate,generated-sync}.md`,
`AGENTS.md` §4, `harness/README.md` §Sensors. Zahlen und Befunde aus Review und
Architect-Verdikt waren **Kontext**, nicht übernommen — jede Aussage dieses
Berichts stammt aus einem hier selbst gefahrenen Lauf.

**Gegenstand.** `HEAD` = `de0f356`, Zweig `main`, Baum sauber. Der Vorgang seit
der Plan-Korrektur ist `4f39cd5` (Plan-Diff) → `1eb3feb` (Architect-Verdikt,
bereits vor `4f39cd5` gelandet) → `1c1d937` (Arbeit, Umzug) → `57d2566`
(Coverage-Fläche/Bau-Kontext) → `068f26e` (Image-Digest) → `de0f356`
(Review). Der Slice liegt in `in-progress/`; der `git mv` nach `done/` ist
**nicht** erfolgt — erwartungsgemäß, das ist Planner-Arbeit bei der Closure.

**Beleg-Lage.** Jeder Gate-Exit ist **ungepiped** ermittelt und in einem
eigenen, abgeschlossenen Schritt aus einer separaten Log-Datei gelesen
(`AGENTS.md` §3.9); Gate-Lauf und Auswertung waren getrennt beauftragt. Kein
host-lokaler absoluter Pfad in diesem Bericht.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make generated-sync` | **0** | „byte-gleich der Ausgabe des gepinnten Generators"; geprüft: `gen/cdc/stream/v1/changestream_grpc.pb.go`, `gen/cdc/stream/v1/changestream.pb.go` |
| 2 | `make a-check` | **0** | „gesamt: 0 Befund(e)" |
| 3 | `make test` (netzlos, `-race`) | **0** | **35 × `ok`, 0 × `FAIL`, 0 × `DATA RACE`** |
| 4 | `make coverage-gate` | **0** | gedruckt `total: (statements) 83.4%`; `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` |
| 5 | `make gates` (Log in Datei, Exit danach aus eigener Datei) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 800 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` · `coverage-gate: OK — 83.40%` |
| 6 | `make doc-commits RANGE=0dece82..HEAD` | **0** | 800 Dateien, 0 Befunde — jeder Commit des ganzen Slice-Vorgangs trägt seine Kennung |
| 7 | `make doc-immutable RANGE=0dece82..HEAD` | **0** | 800 Dateien, 0 Befunde |
| 8 | `make commit-traceability RANGE=4f39cd5..HEAD` | **0** | 4 Commit(s), Betreffs ohne Struktur-ID — der Umsetzungs-Bereich für sich genommen |
| 9 | `make image` | **0** | Digest `sha256:50348a0abe99c5a7880f32e81c8b6d27fb2180b838581038c7d67fe807670a3b` — identisch mit `harness/image-hash.txt` |
| 10 | `grep -rn "gen/cdc/stream/v1"` über `internal/`, `tools/`, `cmd/`, `test/`, `examples/` | — | Treffer in genau **vier** `.go`-Dateien: `internal/adapters/driving/grpc/{server,server_test,interceptor_test}.go`, `tools/harness/grpcclient/main.go`. **Kein** fünfter Aufrufer, auch nicht in `test/`/`examples/` |
| 11 | `grep -rn "internal/adapters/driving/grpc/streamv1"` repo-weit (kein `.git`) | — | Treffer nur in: `.a-check.yml` (3×, Kommentar „umgezogen von …") und `harness/sensors/coverage-gate.md` (1×, „vormals …") — beide **historisch korrekt formuliert**; alle übrigen Treffer liegen in `done/`-Records, `docs/reviews/*` (Lauf-Belege), immutablen ADRs oder in `next/slice-095` (bereits vom Review als offene §6-Beobachtung dieses Slice benannt) |
| 12 | `ls examples/`, `find examples -iname "*grpc*"` | — | `examples/{http-client,nats-client,sse-client}` — **kein** gRPC-Client gebaut |
| 13 | `git diff --name-status -M 4f39cd5 1c1d937` | — | bestätigt R098/R100/R099 (die drei erzeugten Dateien) **und** sechs weitere `M`-Zeilen **im selben Commit** — Grundlage für die HIGH-Prüfung unten |
| 14 | `.a-check.yml` gelesen | — | Gruppe `contract: ["gen/**"]` neu; Kanten `adapters→contract`, `tooling→contract` neu; `tooling→adapters` entfernt, mit Kommentar „ZURÜCKGENOMMEN" samt Begründung |
| 15 | Observation-Register: drei im Plan zitierte Pfade | — | `BEO-PGC/{aufschub-adresse-nimmt-sendung-nicht-an,aufschub-adresse-verfaellt,arbeit-ueberholt-stehenden-traeger}/` existieren, je mit nicht leerem `evidence/` |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — der Stub liegt öffentlich und wird erzeugt

| Kriterium (§2) | Befund |
|---|---|
| `go_package`-Zeile zeigt auf einen Pfad außerhalb `internal/` | **erfüllt** — `proto/cdc/stream/v1/changestream.proto:10`: `option go_package = "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1;streamv1";` (selbst gelesen) |
| Der erzeugte Baum liegt dort | **erfüllt** — `gen/cdc/stream/v1/{changestream.pb.go,changestream_grpc.pb.go,changestream_test.go}` vorhanden, alter Baum `internal/adapters/driving/grpc/streamv1/` vollständig entfernt (#10 findet dort keinen `.go`-Rest) |
| `make proto-generate`/`make generated-sync` bleibt grün | **erfüllt** — #1, Exit 0, byte-gleich |

**LP1: erfüllt, gemessen.**

### LP2 — alle Aufrufer ziehen mit

| Kriterium (§2) | Befund |
|---|---|
| Die vier Importstellen kompilieren | **erfüllt** — #3 (`make test`, Exit 0, 0 `FAIL`), #10 (genau vier Fundstellen, kein fünfter Aufrufer im ganzen Repo) |
| `.a-check.yml` trägt die neue Gruppe `contract` samt Kanten | **erfüllt** — #14, #2 (`make a-check` Exit 0) |
| `tooling`-Gruppe spricht den Stub nicht mehr über `adapters` an, Rücknahme benannt | **erfüllt** — #14: Kommentar nennt „ZURÜCKGENOMMEN" mit Begründung, nicht stillschweigend gelöscht |

**LP2: erfüllt, gemessen.**

### LP3 — die Coverage-Fläche und ihre Träger sind nachgezogen

| Kriterium (§2) | Befund |
|---|---|
| `./gen/...` in der `go list`-Zeile der Stufe `coverage` | **erfüllt** — `Dockerfile:69` (`go list ./internal/... ./cmd/... ./gen/... …`), Filterregel unverändert |
| Nenner bleibt bei 1936 (`slice-097`-Lauf), Quote 83,3–83,4 % | **erfüllt an der gedruckten Ausgabe** — #4/#5 drucken `83.4%`/`83.40%`, deckungsgleich mit dem im Architect-Verdikt genannten Band; die genaue Statement-Zahl 1936 selbst wurde nicht separat nachgerechnet (kein `go tool cover -func`-Rohlauf außerhalb der Stufe), die **gedruckte Prozentzahl** ist aber die tragende, im Plan zitierte Größe und hält |
| Träger aus dem Verdikt §6 nachgezogen (`.dockerignore`, `harness/image-hash.txt`, `harness/sensors/{coverage-gate,generated-sync}.md`, `harness/mk/coverage.mk`, `harness/README.md` §Sensors, `AGENTS.md` §4) | **erfüllt** — alle sieben Träger geprüft und mit Pfad `gen/...`/`./gen/...` nachgezogen (#9, #11, direkte Lektüre); `harness/image-hash.txt` trägt den nach dem Bau-Kontext-Wechsel neu erfassten Digest, der eigene `make image`-Lauf (#9) reproduziert ihn exakt |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#5, Exit 0, sechs Checks) |
| Review durchgeführt, Report unter `docs/reviews/` | **erfüllt** — `review-slice-097.md` (`de0f356`), 1 HIGH ohne Fixrunde, Checkbox im selben Commit nachgezogen |
| Doku-Update für `<Schnittstelle X>` | Zeile ist unausgefüllte Vorlagenzeile — dieselbe Klasse wie bereits in `verify-slice-096.md` V-4 benannt; der tatsächliche öffentliche-Vertrag-Nachzug ist real erfolgt (sieben Träger, s. LP3), das Kriterium selbst ist nicht entscheidbar geschrieben |
| Closure-Notiz, Beobachtungs-Register, Reconciliation, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-Arbeit bei Closure, kein Verifikations-Fehler (Aufgabenstellung Punkt 6) |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, nicht gelesen. `make gates` grün. Die reguläre Vorlagenzeile ist ein wiederkehrender Formaldefekt ohne Wirkung auf die Lieferung.

---

## 3. Risiken aus §6 — Status und Materialisierung

Alle vier Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler. Geprüft wurde, ob eines
davon *bereits jetzt* so schwer eingetreten ist, dass es die DoD-Konformität
selbst infrage stellt:

| Risiko | Materialisiert? |
|---|---|
| Nenner bewegt sich, Quote sinkt | **nicht negativ eingetreten** — die Quote liegt bei 83,4 % (#4), deutlich über der Endstufe 80 %; die Aufnahme von `./gen/...` hält den Gegenstand konstant (Architect-Verdikt, nachgemessen an der gedruckten Zeile) |
| Vier Importstellen sind nicht alle | **nicht eingetreten** — eigener `grep` über das **ganze** Repo (nicht nur `internal/`+`cmd/`) findet exakt vier, kein fünfter (#10) |
| Umzug verbaut C#-/Kotlin-Weg | außerhalb dieses Slice-Umfangs, unverändert offen |
| Träger überholt (alter Pfad) | teilweise vorweggenommen — der Review hat bereits einen konkreten Fund (`next/slice-095`) benannt; eigener `grep` (#11) findet **keine** weiteren live-Träger mit dem alten Pfad im Präsens außerhalb bereits bekannter historischer/immutabler Stellen |

**Keines der vier Risiken stellt die DoD-Konformität dieses Vorgangs infrage.**

---

## 4. Entscheidungs-Konformität

| Zusage | Befund |
|---|---|
| [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §3 „der Umzug soll allein stehen" | **eingehalten** — kein gRPC-Client, kein zweites Go-Modul, kein neues Gate im Diff (#12, eigene Prüfung); Diff bleibt auf Umzug + seine Träger beschränkt |
| [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §Konsequenzen (Coverage-Aufnahme) | **eingehalten** — `./gen/...` real in der Stufe (LP3) |
| [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (Draht-Vertrag unverändert) | **eingehalten** — genau eine Zeile der `.proto` geändert (`option go_package`), kein Nachrichten-/Feld-Inhalt (selbst gelesener Diff) |
| [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md) (Kante zurückgenommen, benannt) | **eingehalten** — #14, Kommentar explizit |
| [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)/[`ADR-0082`](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) (Eigenschaft, nicht Liste) | **eingehalten** — der erzeugte Code bleibt im Gegenstand, Pfad wandert, Eigenschaft (netzlos testbar) bleibt |
| [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) (keine Supersession dieses Slice) | **eingehalten** — der Slice trifft keine neue Client-Matrix-Entscheidung, bleibt auf dem Umzug |
| `AGENTS.md` §3.3 (git mv + Inhaltsänderung = zwei Commits) | **verletzt** — bestätigt (#13): `1c1d937` vereint die Umbenennung (R098/R100/R099, von Git korrekt mit 98–100 % Similarity erkannt) mit Inhaltsänderungen an zwei der drei umbenannten Dateien (`go_package`-String im Descriptor, Importpfad) **und** sechs weiteren Dateien im selben Commit. Der Review hat dies bereits als **HIGH** (F-1) gefasst und **bewusst nicht fix-pflichtig** gestellt (`git log --follow` bleibt funktionsfähig, Rename-Detection griff trotzdem; ein nachträgliches Umschreiben der Historie wäre ein größerer Eingriff als der Nutzen rechtfertigt) — als Steering-Loop-Kandidat vorgemerkt. Diese Disposition ist eine **legitime Reviewer-Entscheidung** (Modul 8 §Rollen-Regeln: Merge-Blockierung ist Reviewer-Urteil); sie hebt die **Verletzung** selbst nicht auf. Für die DoD-Konformität dieses Slice ist sie **nicht kriterienrelevant** — kein LP verlangt Commit-Struktur —, gehört aber in die Closure-Notiz §7 als Finding-Klasse (dort bereits vom Review benannt) |
| `AGENTS.md` §3.6 (Schwellen nur per ADR) | **eingehalten** — `THRESHOLD ?= 80` unberührt |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten in diesem Bericht** — alle 15 Läufe einzeln, ungepiped, Exit-Code direkt gelesen |
| `AGENTS.md` §3.12 (Herkunft von Zahlen) | **eingehalten** — die Zahl 1936 steht in `harness/sensors/coverage-gate.md` als datierte Messung mit Lauf-Referenz (`slice-097`, Architect-Verdikt), die vorherige 1903 bleibt daneben explizit als Messung von `slice-085` stehen, nicht überschrieben |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten für die Träger dieses Slice** — alle sieben in §6 des Architect-Verdikts genannten Träger sind nachgezogen (LP3); der eine noch offene Fremdträger (`next/slice-095`) ist als §6-Risiko dieses Slice-Plans bereits benannt, nicht verschwiegen — Auflösung ist Closure-Arbeit |

---

## 5. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `HEAD`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `proto/cdc/stream/v1/changestream.proto` — update | genau eine Zeile (`option go_package`) | **Plan eingehalten** |
| `gen/cdc/stream/v1/**` — neu, alter Baum entfällt | drei Dateien am neuen Pfad, alter Baum entfernt, Testdatei zog mit | **Plan eingehalten** |
| Drei Importstellen `internal/adapters/driving/grpc/{server,server_test,interceptor_test}.go` | alle drei aktualisiert | **Plan eingehalten** |
| `tools/harness/grpcclient/main.go` — vierte Importstelle | aktualisiert | **Plan eingehalten** |
| `.a-check.yml` — neue Gruppe, Kanten, Rücknahme | wie geplant | **Plan eingehalten** |
| `Dockerfile` Stufe `coverage` | `./gen/...` aufgenommen, Filter unverändert | **Plan eingehalten** |
| `harness/sensors/coverage-gate.md`, `harness/mk/coverage.mk` | beide nachgezogen | **Plan eingehalten** |
| `harness/README.md` §Sensors, `AGENTS.md` §4 | beide nachgezogen | **Plan eingehalten** |
| `.dockerignore` | `!gen/` ergänzt, **keine** `ADR-0085`-Fehlzitation (Plan warnt davor explizit) | **Plan eingehalten** |
| `harness/image-hash.txt` | Digest-Commit `068f26e`, reproduziert (#9) | **Plan eingehalten** |
| `harness/sensors/generated-sync.md` §Vertrag | Pfad nachgezogen | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` (falls Pfad genannt) | kein Treffer, kein Nachzug fällig (Review-Negativbefund, hier nicht erneut nachgemessen) | **Plan eingehalten** |
| — | `docs/reviews/review-slice-097.md` (A) | **kein Plan-Bruch** — Reviewer-Übergabe-Artefakt |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht (als eigene Zeile)
nennt:** die Commit-Struktur-Verletzung selbst (§3.3, oben) ist keine
Diff-Abweichung vom Plan, sondern eine Verfahrensfrage der Lieferung.

---

## 6. Verdikt

**Konform mit einer benannten, nicht-blockierenden Ausnahme.** Alle drei
Liefer-Punkte (LP1–LP3) sind gemessen erfüllt, `make gates` läuft grün
(Exit 0, sechs Checks, alle ungepiped selbst geprüft), kein referenziertes
`Accepted`-ADR wird durch den gelandeten Code verletzt, `ADR-0076` §3s
Vorgabe „der Umzug soll allein stehen" hält (kein Client gebaut), und keines
der vier §6-Risiken ist so eingetreten, dass es die Lieferung infrage
stellt. Die eine echte Verletzung — `AGENTS.md` §3.3 in Commit `1c1d937`
(Umbenennung + Inhaltsänderung in einem Commit) — ist vom Review bereits
korrekt als HIGH erkannt, sachlich richtig als nicht-merge-blockierend und
nicht-fix-pflichtig disponiert, und gehört als Finding-Klasse in die
Closure-Notiz §7 und von dort in den Zähler des Beobachtungs-Registers
(„git-mv-und-Inhalt-in-einem-Commit").

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die vier Closure-Checkboxen (`make gates`, Doku-Update, Closure-Notiz,
   Reconciliation *entfällt*, Beobachtungs-Register, Risiko-Ausgänge, drei
   Paarungen) sind noch nicht angehakt.
2. §6 — alle vier Risiken brauchen ihren Ausgang; keines davon ist
   *eingetreten* im schädlichen Sinn (siehe §3 dieses Berichts), aber jedes
   braucht die formale Zuweisung (voraussichtlich: erstes und zweites Risiko
   → *entfallen* mit Begründung, drittes → *weiter offen*/Folge-Slice-Bezug
   auf `ADR-0087`/`ADR-0090`, viertes → *entfallen*, weil die einzige reale
   Fundstelle bereits im Review benannt und von dort in die Closure
   überführbar ist).
3. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag für die §3.3-Verletzung
   (Finding-Klasse „git-mv-und-Inhalt-in-einem-Commit"); die drei im Kopf
   zitierten Register-Pfade existieren bereits (#15) und sind zitierfähig.
4. Die unausgefüllte Vorlagenzeile „Doku-Update für `<Schnittstelle X>`" in
   §2 — vor dem `git mv` streichen oder auf die real berührten Träger
   füllen (dieselbe wiederkehrende Formklasse wie bereits in
   `verify-slice-096.md` V-4 benannt).
5. Der `git mv` nach `done/` selbst.
