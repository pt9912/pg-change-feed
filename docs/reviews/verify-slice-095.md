# Verifikationsbericht: slice-095 (2. Durchlauf — LP2 `examples/grpc-client/` + Rest-LP3) — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-095` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(die drei Clients, Import-Grenze §2, Träger `.a-check.yml`/Handbuch §4/§7) ·
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (Server-Stream,
Fire-and-Forget ohne Replay) · [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
(volle Matrix — §6 „examples-Gruppe bleibt eine Go-Aussage", Slice-Schnitt-
Empfehlung Zeile 7) · [`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
(Form-Vorbild) — sowie die Hard Rules `AGENTS.md` §3.3, §3.9, §3.12, §3.13.
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-095.md`](review-slice-095.md) abgeschlossen) und **nicht**
gegen realen Bedarf (Validator, nicht ausgelöst).

**Gegenstand dieses Durchlaufs:** nur LP2 (`examples/grpc-client/`) und der
Rest von LP3 (`.a-check.yml`, Handbuch) — Commits `d953356` (grpc-client),
`1812e36` (Träger-Nachzug). LP1 (die zwei HTTP-Familien-Clients) wurde im
ersten Durchlauf geliefert, akzeptiert und **hier nur gegengecheckt**, nicht
erneut vollständig verifiziert.

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan am Stand `HEAD`
vollständig gelesen (§1–§8), dazu `review-slice-095.md`, `ADR-0076`,
`ADR-0060`, `ADR-0090` (vollständig, inkl. §Slice-Schnitt-Empfehlung),
`examples/grpc-client/{main.go,format.go}`, `.a-check.yml`,
`docs/user/benutzerhandbuch.md` (Abschnitt „Zugriff über den
gRPC-Change-Stream" und Versionshistorie). Zahlen und Befunde aus dem
Review waren **Kontext**, nicht übernommen — jede Aussage dieses Berichts
stammt aus einem hier selbst gefahrenen Lauf.

**Gegenstand/Umgebung.** `HEAD` lag zu Beginn dieser Sitzung bei `117b9d9`
(Review-Commit). Während der Sitzung hat ein **paralleler, unabhängiger
Agenten-Lauf** (separater Planungs-Auftrag: `slice-098`…`103`,
Sprachwurzel-Clients C#/Kotlin, `ADR-0090`) weitere Commits auf `main`
gelandet, darunter einen Fix an `docs/reviews/review-slice-095.md` selbst
(`e7e7e80`, unverlinkte `ADR-0076`/`ADR-0090`-Erwähnungen gebacktickt). Das
ist **nicht Gegenstand dieser Verifikation** (Aufgabenstellung); ein erster
`make docs-check`-Lauf dieser Sitzung fing einen **transienten** Zustand
dieses parallelen Laufs ab (4 Befunde in `review-slice-095.md`, obwohl die
Datei zu diesem Zeitpunkt bereits die Backtick-Korrektur trug) — ein
zweiter, unmittelbar folgender Lauf war bereits wieder grün (0 Befunde,
siehe unten). Vermerkt, nicht blockierend: kein Befund dieses zweiten Laufs
betrifft `slice-095`s eigene Dateien.

**Beleg-Lage.** Jeder Gate-Exit ist **ungepiped** ermittelt und in einem
eigenen, abgeschlossenen Schritt aus einer separaten Log-Datei gelesen
(`AGENTS.md` §3.9); Gate-Lauf und Auswertung waren getrennt beauftragt. Kein
host-lokaler absoluter Pfad in diesem Bericht.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make test` (netzlos, `-race`) | **0** | 36 × `ok`, 0 × `FAIL`, 0 × `DATA RACE`; `ok github.com/pt9912/pg-change-feed/examples/grpc-client 1.035s` |
| 2 | `make a-check` | **0** | „gesamt: 0 Befund(e)" |
| 3 | `make docs-check` (1. Lauf, überlappte mit dem parallelen Agenten-Commit) | **2** | 7 Befunde — 3 in `docs/plan/planning/open/slice-098-…md`/`slice-101-…md` (fremder, nicht-committeter bzw. zwischenzeitlich committeter Bestand des Parallel-Laufs), 4 in `docs/reviews/review-slice-095.md:42/55/71` (`id-unlinked` für `ADR-0076`/`ADR-0090`) — als **transient** bewertet, siehe unten |
| 4 | `make docs-check` (2. Lauf, sofort danach, keine Änderung meinerseits dazwischen) | **0** | „812 Datei(en) geprüft, 0 Befund(e)" — dieselbe Datei, derselbe Inhalt, jetzt sauber. Diagnose: Lauf #3 hat einen Zwischenzustand des parallelen Commits `e7e7e80` erwischt |
| 5 | `make commit-traceability RANGE=15715e9..HEAD` | **0** | 9 Commit(s) (umfasst zwischenzeitlich auch die Parallel-Commits), Betreffs ohne Struktur-ID |
| 6 | `make gates` (Log in Datei, Exit danach aus eigener Datei) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 812 Datei(en) geprüft, 0 Befund(e)` (2×, docs-check + Teil des Bündels) · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` |
| 7 | `make doc-commits RANGE=15715e9..1812e36` | **0** | 812 Dateien, 0 Befunde — scharf auf diesen Durchlauf geschnitten (vor den Parallel-Commits) |
| 8 | `make doc-immutable RANGE=15715e9..1812e36` | **0** | 812 Dateien, 0 Befunde |
| 9 | `grep -rn "internal/" examples/grpc-client/*.go` | Exit 1 (kein Treffer) | **kein** `/internal/`-Import — eigener Lauf, nicht aus dem Review übernommen |
| 10 | `grep -n "google.golang.org/grpc\|google.golang.org/protobuf" go.mod` | — | beide bereits vor diesem Durchlauf gepinnt (`v1.83.2`/`v1.36.12`, aus `ADR-0060`); `go.mod`/`go.sum` unverändert in diesem Diff (bestätigt: `git diff --stat 15715e9..1812e36` führt sie nicht) |
| 11 | Lektüre `examples/grpc-client/main.go` | — | `grpc.NewClient` → `streamv1.NewChangeStreamClient` → `client.StreamChanges(ctx, &streamv1.StreamChangesRequest{})` → `for { stream.Recv() }`; Token über `metadata.AppendToOutgoingContext` (`authorization: Bearer <token>`); Doc-Kommentar nennt explizit „der Stream kennt kein Replay (`ADR-0060`)" |
| 12 | Lektüre `.a-check.yml` | — | Gruppe `examples` nennt jetzt alle vier Clients (http, sse, grpc, nats) im Kommentar; neue Kante `{from: examples, to: contract}`, begründet mit dem seit `slice-097` öffentlichen Pfad `gen/cdc/stream/v1` |
| 13 | Lektüre `docs/user/benutzerhandbuch.md` Z. 700–704 | — | `**Beispiel:**`-Absatz unter „Zugriff über den gRPC-Change-Stream" nennt `examples/grpc-client`, Startbefehl, `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER` |
| 14 | `grep -n "^| 1\." docs/user/benutzerhandbuch.md` (Versionshistorie) | — | lückenlose Folge `…1.16, 1.17, 1.18, 1.19` — `1.19` trägt genau diesen Durchlauf |
| 15 | Beobachtungs-Register: vier im Plan zitierte Pfade | — | `BEO-PGC/{handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche,handbuch-versionshistorie-uebersprungen,dod-begruendung-unzutreffende-tatsachenbehauptung,arbeit-ueberholt-stehenden-traeger}/` existieren, je mit nicht leerem `evidence/` |
| 16 | `ADR-0090` §Slice-Schnitt-Empfehlung, Zeile 7 gelesen | — | „**Go-gRPC-Client** … LP2 die Kante in `.a-check.yml`" — bestätigt: die breitere „Go-Aussage"-Reformulierung (§6) ist Folgepflicht des **Träger-Zugs**, nicht an Zeile 7/diesen Slice gebunden |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — die zwei HTTP-Familien-Clients (Gegencheck, nicht erneut vollständig verifiziert)

| Kriterium (§2) | Befund |
|---|---|
| `examples/http-client/`, `examples/sse-client/` liegen vor, kompiliert und getestet | **weiterhin erfüllt** — beide in `make test`s 36 `ok`-Zeilen enthalten (#1); unverändert seit dem ersten Durchlauf |
| Handbuch-Zeilen (Version 1.17) | **weiterhin vorhanden** (#14) |

**LP1: weiterhin erfüllt, gegengecheckt.**

### LP2 — der gRPC-Client

| Kriterium (§2) | Befund |
|---|---|
| `examples/grpc-client/` liegt vor, kompiliert und getestet über `make test` | **erfüllt** — #1 |
| Öffnet den **Server-Stream** `ChangeStream/StreamChanges` | **erfüllt** — #11, `client.StreamChanges(...)` gefolgt von `stream.Recv()` in Schleife, exakt die in `ADR-0060` Teilfrage 1/3 festgelegte Server-Streaming-/Fire-and-Forget-Semantik |
| Import über den öffentlichen `gen/cdc/stream/v1`-Pfad, kein `internal/`-Import | **erfüllt** — #9 (eigener `grep`, kein Treffer), #11 (Importzeile `streamv1 "github.com/pt9912/pg-change-feed/gen/cdc/stream/v1"`) |

**LP2: erfüllt, gemessen.**

### LP3 — die zwei Träger sind nachgezogen

| Kriterium (§2) | Befund |
|---|---|
| `.a-check.yml` nennt alle vier Clients | **erfüllt** — #12 |
| Neue Kante `{from: examples, to: contract}` mit Objekt | **erfüllt** — #12, #2 (`make a-check` Exit 0) |
| `docs/user/benutzerhandbuch.md` führt alle vier Clients namentlich an ihren Zugriffs-Abschnitten, mit Startform | **erfüllt** — #13 (grpc); http/sse/nats bereits aus früheren Durchläufen bestätigt (#14 Versionshistorie zeigt 1.16/1.17 dafür) |
| Versionshistorie lückenlos, mit Zeile für diesen Durchlauf | **erfüllt** — #14, `1.19` |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt** — #6, Exit 0, sechs Checks |
| Review durchgeführt, Report unter `docs/reviews/` | **erfüllt** — `review-slice-095.md` (`117b9d9`), 0 HIGH/MEDIUM, 1 LOW ohne Fixrunde |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-Arbeit bei Closure, kein Verifikations-Fehler (Aufgabenstellung Punkt 5) |

**Ergebnis §2:** LP1 weiterhin erfüllt, LP2 und LP3 dieses Durchlaufs
gemessen erfüllt. `make gates` grün.

---

## 3. Risiken aus §6 — Status und Materialisierung

Alle vier Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler. Geprüft wurde, ob eines
davon *bereits jetzt* so schwer eingetreten ist, dass es die
DoD-Konformität selbst infrage stellt:

| Risiko | Materialisiert? |
|---|---|
| Ein Client braucht mehr als die Standardbibliothek | **eingetreten, aber gedeckt** — `examples/grpc-client` importiert `google.golang.org/grpc` und dessen Unterpakete; das ist durch `ADR-0076` §2 Festlegung ausdrücklich als zulässiges „öffentliches Fremdmodul" benannt (die Zulässigkeitsregel nennt `google.golang.org/grpc`/`google.golang.org/protobuf` namentlich), keine neue Abhängigkeit (#10, bereits seit `ADR-0060` gepinnt), kein `go.mod`/`go.sum`-Diff. Kein DoD-relevantes Risiko |
| Ein Client ist ohne laufenden Dienst nicht kompilierbar | **nicht eingetreten** — `make test` kompiliert `examples/grpc-client` ohne laufenden Dienst (#1); der Import zeigt auf den öffentlichen, statisch vorhandenen Stub, nicht auf einen internen, Bootstrap-ziehenden Pfad |
| Der Handbuch-Nachzug wird vergessen | **nicht eingetreten** — #13/#14, Absatz und Versionshistorie-Zeile beide vorhanden |
| Ein Träger wird überholt, den dieser Slice nicht anfasst | **nicht eingetreten für diesen Durchlauf** — der Review hat bereits `next/slice-095` selbst (den eigenen Plan) als frühere Fundstelle benannt (aus dem `slice-097`-Umzug); ein weiterer `grep`-Fund dieses Laufs (#9) beschränkt sich auf den eigenen, bereits erwarteten Scope |

**Keines der vier Risiken stellt die DoD-Konformität dieses Durchlaufs
infrage.**

---

## 4. Entscheidungs-Konformität

| Zusage | Befund |
|---|---|
| [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §2 (Import-Grenze: Standardbibliothek + öffentliche Fremdmodule, kein `internal/`) | **eingehalten** — #9, #11; `google.golang.org/grpc`/`credentials/insecure`/`metadata` sind exakt die von §2 genannte Klasse |
| [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §1/§4 (Startform, `.a-check.yml`-Träger statt zweitem Modul) | **eingehalten** — Startform `go run ./examples/grpc-client`, Flags `-addr`/`-token`, ENV-Fallback (#11); keine zweite `go.mod` |
| [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (Server-Stream, Fire-and-Forget ohne Replay, `reader`/`admin`-Token über Metadata) | **eingehalten** — #11: `StreamChanges` als Server-Stream-RPC, Token über `authorization`-Metadata mit `Bearer`-Präfix, kein Replay-Code im Client (Doc-Kommentar spricht es aus) |
| [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) §6 (examples-Gruppe bleibt Go-Aussage) | **eingehalten** — der aktuelle `.a-check.yml`-Kommentar beschreibt weiterhin nur Go-Verhalten der vier vorhandenen Clients; die von §6 verlangte explizite Umformulierung auf „Go-Aussage" ist per §Slice-Schnitt-Empfehlung Zeile 7 (#16) **nicht** dieser Slice zugeordnet, sondern dem Träger-Zug, der mit dem ersten Fremdsprachen-Client fällig wird — Review-Negativbefund dazu ist beim eigenen Nachvollzug bestätigt, nicht nur übernommen |
| [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) (kein neuer `languages`-Schlüssel, keine neue `resolution`) | **eingehalten** — #12, keine solchen Schlüssel in `.a-check.yml` |
| `AGENTS.md` §3.3 (git mv + Inhaltsänderung = zwei Commits) | **nicht einschlägig für diesen Durchlauf** — `d953356` fügt ausschließlich drei neue Dateien hinzu (keine Umbenennung), `1812e36` ist eine reine Inhaltsänderung ohne `git mv`; kein Rename-Vorgang in diesem Durchlauf |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten in diesem Bericht** — alle Läufe in #1 einzeln, ungepiped, Exit-Code direkt gelesen |
| `AGENTS.md` §3.12/§3.13 (Herkunft von Zahlen, Träger-Nachzug) | **eingehalten** — die Coverage-Zahl (83,40 %) ist eine eigene Messung dieses Laufs, nicht übernommen; die zwei fälligen Träger (`.a-check.yml`, Handbuch) sind nachgezogen (LP3) |

---

## 5. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans, beschränkt auf die in diesem
Durchlauf gelieferten Zeilen (LP2/LP3).

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/grpc-client/**` — neu | `main.go`, `format.go`, `format_test.go` | **Plan eingehalten** |
| `.a-check.yml` — update (vier Clients im Kommentar, Kante mit Objekt) | genau das, plus die begründete Kommentarzeile zur Kante | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update (Zugriffs-Abschnitt + Versionshistorie) | `**Beispiel:**`-Absatz unter „Zugriff über den gRPC-Change-Stream", Zeile `1.19` in der Versionshistorie, Versionskopf hochgezählt | **Plan eingehalten** |
| — | `docs/reviews/review-slice-095.md` (Reviewer-Übergabe-Artefakt) | **kein Plan-Bruch** |
| — | DoD-Checkbox-Nachzug LP1/LP2/LP3 in der Plan-Datei selbst | **kein Plan-Bruch** — Skill-konforme Checkbox-Setzung im selben Commit wie der Review |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht als eigene Zeile
nennt:** die Plan-§3-Begründungszeile „keine Kante — `ADR-0079` Festlegung 1"
blieb stehen, obwohl in diesem Durchlauf tatsächlich eine Kante geliefert
wurde — bereits vom Review als F-1 (LOW) benannt, kein Code-Bruch, reine
Plan-Text-Inkonsistenz.

---

## 6. Verdikt

**Konform.** Beide in diesem Durchlauf fälligen Liefer-Punkte (LP2, LP3)
sind gemessen erfüllt, LP1 bleibt beim Gegencheck bestätigt. `make gates`
läuft grün (Exit 0, sechs Checks, alle ungepiped selbst geprüft). Keines der
referenzierten `Accepted`-ADRs (`ADR-0076`, `ADR-0060`, `ADR-0090`,
`ADR-0079`) wird durch den gelandeten Code verletzt — die einzige
scheinbare Abweichung (gRPC-Fremdmodul-Import) ist durch `ADR-0076` §2
namentlich gedeckt, und die scheinbar fehlende „Go-Aussage"-Reformulierung
in `.a-check.yml` ist laut `ADR-0090`s eigener Slice-Zuordnung (Zeile 7)
nicht dieser Slice, sondern dem Träger-Zug der Fremdsprachen-Clients
zugeordnet. Keines der vier §6-Risiken ist so eingetreten, dass es die
Lieferung infrage stellt. Das einzige Review-Finding (F-1, LOW: stehen
gebliebene Plan-Begründung) ist merge-nicht-blockierend und ohne
Auswirkung auf den gelieferten Code.

**Umgebungs-Hinweis (nicht Teil des Verdikts):** Ein erster
`make docs-check`-Lauf dieser Sitzung fing einen transienten Zwischenstand
eines parallel laufenden, unabhängigen Agenten-Vorgangs ab (4 scheinbare
`id-unlinked`-Befunde in einer bereits korrekten Datei); ein sofort
folgender Lauf war grün. Kein Befund betraf `slice-095`s eigene Lieferung.

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die verbleibenden Closure-Checkboxen (Closure-Notiz,
   Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) sind noch nicht
   angehakt.
2. §6 — alle vier Risiken brauchen ihren Ausgang; keines davon ist
   *eingetreten* im schädlichen Sinn (siehe §3 dieses Berichts), aber jedes
   braucht die formale Zuweisung (voraussichtlich: erstes → *entfallen* mit
   Begründung [durch `ADR-0076` §2 gedeckt], zweites → *entfallen*, drittes
   → *entfallen*, viertes → *entfallen* mit Verweis auf den bereits im
   Review benannten historischen Fund).
3. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag; die vier im Kopf
   zitierten Register-Pfade existieren bereits (#15) und sind zitierfähig;
   naheliegender Kandidat für den Eintrag: die zwei Register-Treffer mit
   **je 3×**, die dieser Slice einlöst (`handbuch-nicht-nachgezogen-bei-
   neuer-betreiber-oberflaeche`, `handbuch-versionshistorie-uebersprungen`).
4. F-1 aus dem Review (stehengebliebene §3-Begründung „keine Kante") — bei
   Closure formlos nachziehen oder stehenlassen (Review-Disposition).
5. Der `git mv` nach `done/` selbst.
