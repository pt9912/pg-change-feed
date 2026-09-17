# Review-Report: slice-095 (2. Durchlauf: LP2 `examples/grpc-client/` + Rest-LP3) — 2026-09-17

**Review-Art:** Code — geprüft gegen Plan + ADR + Konventionen (Modul 10
§Drei Review-Arten).

**Gegenstand:** `git diff 15715e9..HEAD` (Commits `d953356`,
`1812e36`) — Umfang: `examples/grpc-client/{main.go,format.go,format_test.go}`,
`.a-check.yml`, `docs/user/benutzerhandbuch.md`,
`docs/plan/planning/in-progress/slice-095-beispiel-clients-drei.md`
(nur die DoD-Checkbox-Zeilen). LP1 (http-client, sse-client) wurde im ersten
Durchlauf geliefert und akzeptiert und ist **nicht** erneut geprüft.

**Skill:** `.harness/skills/reviewer.md` @ Accepted, Schärfung 2026-09-09
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-095-beispiel-clients-drei.md` (vollständig)
- `ADR-0076` (die drei Clients, Startform, Doc-Kommentar, Import-Grenze)
- `ADR-0060` (der gRPC-Server-Stream)
- `ADR-0090` (die volle Matrix — §6 „examples-Gruppe bleibt eine Go-Aussage",
  Festlegung 2/3)
- `ADR-0079` (Form-Vorbild NATS-Client, Festlegung 1 — keine Kante)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.13 Träger-Nachzug), §5
- `harness/conventions.md` (MR-000/MR-001)
- `LH-FA-SST-008` (Live-Change-Stream)

---

## Findings

### F-1 — Plan §3-Begründung „keine Kante" widerspricht der tatsächlich gelieferten Kante

- `kategorie`: LOW
- `quelle`: Maintainability (Plan-Konsistenz)
- `pfad`: `docs/plan/planning/in-progress/slice-095-beispiel-clients-drei.md:127`
- `befund`: §3 der Plan-Datei begründet die geplante Änderung an
  `.a-check.yml` weiterhin mit „keine Kante — `ADR-0079` Festlegung 1"; die
  in diesem Durchlauf tatsächlich gelieferte Änderung fügt genau eine neue
  Kante hinzu (`{from: examples, to: contract}`), korrekt begründet mit dem
  seit `slice-097` öffentlichen Pfad `gen/cdc/stream/v1` und konsistent mit
  `ADR-0076`/`ADR-0090`. Die Zeile in §3 wurde bei der Rückführung/dem zweiten
  Anlauf nicht nachgezogen und ist jetzt eine stehengebliebene, falsche
  Tatsachenbehauptung im Plan-Dokument selbst (nicht im gelieferten Code).
- `verifizierbar`: nein — reine Lese-Prüfung des Plan-Texts, kein Gate-Bezug.
- `klasse`: „Plan-Begründung nicht nach Rückführung nachgezogen"

## Negativbefunde

- geprüft, ohne Befund: `examples/grpc-client/main.go` — Import ausschließlich
  `context`, `flag`, `fmt`, `os`, `google.golang.org/grpc`,
  `google.golang.org/grpc/credentials/insecure`,
  `google.golang.org/grpc/metadata`, `gen/cdc/stream/v1` (eigener `grep -rn
  "internal/" examples/grpc-client/*.go` liefert keinen Treffer). Kein
  `/internal/`-Import, `ADR-0076` §2 Festlegung eingehalten.
- geprüft, ohne Befund: repo-weiter `grep -rn
  "internal/adapters/driving/grpc/streamv1"` — alle verbleibenden Treffer
  liegen in `.a-check.yml`-Kommentaren („umgezogen von …"),
  `harness/sensors/coverage-gate.md` („vormals …"), `done/`-Records,
  immutablen ADRs und Lauf-Belegen (`docs/reviews/*`) — durchweg historisch
  korrekt formuliert, kein lebender Code-Pfad zeigt mehr dorthin (deckt sich
  mit `verify-slice-097.md` #11).
- geprüft, ohne Befund: `.a-check.yml` — die neue Kante
  `{from: examples, to: contract}` ist begründet kommentiert; der
  `examples`-Gruppen-Kommentar nennt jetzt konsistent alle vier Clients
  (http, sse, grpc, nats) und ihre gemeinsame Import-Grenze; `make a-check`
  läuft mit „gesamt: 0 Befund(e)".
- geprüft, ohne Befund: `ADR-0090` §6/Festlegung 2/3 — die dort verlangte
  Umformulierung des `examples`-Gruppenkommentars auf eine explizite
  „Go-Aussage" (nicht-Go-Quellen mechanisch ungeprüft, Fremdsprachen-Stubs
  nicht im Baum) ist laut `ADR-0090`s eigener Slice-Schnitt-Empfehlung (Zeile
  7, LP2) **nicht** Teil dieses Slices — Zeile 7 verlangt für den
  Go-gRPC-Client nur „die Kante in `.a-check.yml` mit ihrem Objekt", die
  breitere Reformulierung ist an die Slices gebunden, die tatsächlich einen
  ersten nicht-Go-Client einführen (Zeilen 1–6). Festlegung 2/3
  (`.proto`-Weg, erzeugte Stubs) berührt Go nicht: `examples/grpc-client`
  nutzt weiterhin den bereits committeten Stub unter `gen/cdc/stream/v1`,
  keine Build-Kontext-Änderung nötig.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — der neue
  `**Beispiel:**`-Absatz sitzt unter `### Zugriff über den gRPC-Change-Stream`
  (Zeile 655), nach der Zustellsemantik/Nachvollziehbarkeit-Prosa, an
  derselben Stelle wie die Beispiel-Absätze der drei anderen Clients;
  Startbefehl `go run ./examples/grpc-client`, Umgebungsvariablen
  `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER` korrekt. Die neue
  Versionshistorie-Zeile trägt `1.19`, fortlaufend nach `1.18` (kein
  Sprung — eigener `grep -n "^| 1\." docs/user/benutzerhandbuch.md`
  bestätigt die lückenlose Folge `…1.16, 1.17, 1.18, 1.19`); der
  Versionskopf ist ebenfalls auf `1.19` hochgezählt. Keine neue Instanz von
  `BEO-PGC/handbuch-versionshistorie-uebersprungen`.
- geprüft, ohne Befund: `examples/grpc-client/format_test.go` — beide Fälle
  sind aussagekräftig, kein Kosmetik-Test: der Happy-Path füllt Schema,
  Tabelle, Operation, `change_id` und ein reales JSON-`new_image` und prüft
  die exakte Ausgabezeile per String-Vergleich; der zweite Fall bindet
  gezielt den `DELETE`-Grenzfall (leeres `NewImage`-Feld) und prüft, dass die
  Ausgabe einen leeren Wert trägt statt eines erfundenen Platzhalters — die
  Eingabeseite (Struct-Felder) wird direkt mutiert, nicht nur ein Fake.
- geprüft, ohne Befund: Out-of-Scope-Disziplin (§1 des Plans) — keine
  Änderung am Stream-Vertrag (`gen/cdc/stream/v1` unverändert, `make
  generated-sync` grün), kein eigenes Gate, kein DB-/SQL-Beispiel, `go.mod`/
  `go.sum` unverändert (die genutzten Module `google.golang.org/grpc` u. a.
  sind bereits vorhandene, gepinnte Abhängigkeiten).
- geprüft, ohne Befund: Commit-Traceability — beide Commit-Betreffs
  (`feat(examples): grpc-client…`, `docs(examples): Traeger für
  grpc-client…`) tragen `ADR-0076`/`ADR-0060`-Bezug im Betreff, kein
  `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability` grün (5 Commits,
  0 Befund).
- geprüft, ohne Befund: DoD-Checkbox-Setzung LP1/LP2/LP3/`make gates` —
  eigener Nachvollzug: `make test` kompiliert und testet
  `examples/grpc-client` (`ok … examples/grpc-client 1.019s`), `make
  a-check`, `make generated-sync`, `make baseline-verify`,
  `make commit-traceability` liefen alle grün. `make docs-check` meldet
  einen Befund (`id-unlinked` in `docs/plan/planning/open/slice-098-…md`) —
  dieser liegt in einer **unversionierten** Datei eines parallel laufenden
  Agenten (nicht Teil dieses Diffs, nicht committet); mit dieser Datei
  temporär entfernt läuft `make docs-check` mit „0 Befund(e)". Die
  Checkbox-Setzung ist damit für den tatsächlichen Diff-Umfang dieses Slice
  gedeckt.
- geprüft, ohne Befund: Doc-Kommentar von `examples/grpc-client/main.go` —
  nennt Schnittstelle (gRPC-Server-Stream), Handbuch-Abschnitt, `LH-FA-SST-008`,
  `ADR-0060`/`ADR-0076`, spricht aus, dass `tools/harness/grpcclient` der
  E2E-Belegträger ist — Form deckungsgleich mit `examples/http-client` und
  `examples/sse-client`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Plan-Begründung nicht nach Rückführung nachgezogen

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM; das einzige LOW-Finding
betrifft einen stehengebliebenen Begründungssatz im Plan-Dokument selbst
(keine semantische Auswirkung auf den gelieferten Code, der ADR-konform
ist), keine Fixrunde am Implementer nötig.

**Übergabe:** Kein Reviewer→Implementer-Rückgabe-Pfeil (0 HIGH, das LOW kann
bei Closure formlos nachgezogen oder stehengelassen werden). Die DoD-Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" wird im
selben Commit wie dieser Report auf `[x]` gezogen (Skill §DoD-Checkbox-Nachzug
ohne Fixrunde). Verifikation (DoD-/Spec-Konformität) bleibt eigener,
separater Schritt (Modul 11).
