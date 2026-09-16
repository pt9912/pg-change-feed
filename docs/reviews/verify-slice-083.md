# Verifikationsbericht: slice-083 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-083` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
— der vierte Beispiel-Client, sechs Festlegungen —,
[`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
— die Form: öffentlicher Draht-Vertrag, kein `/internal/`, Kompilier-Bindung
statt Lauf-Beleg —,
[`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
— der Endpunkt, den der Client benutzt —,
[`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)
/ [`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md) —
leerer Payload, tabellen-granulares Subjekt —; `SPEC-017`, `SPEC-022`;
`AGENTS.md` §3.5, §3.9, §3.10).
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe; der Report
`review-slice-083.md` wurde als Kontext gelesen, **nicht** als Beleg
übernommen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Frischer Kontext:** Dieser Lauf hat den vollständigen Slice-Plan in der
Fassung von `HEAD` gelesen (einschließlich des berichtigten LP1-Textes aus
F-1), dazu die vier ADRs, `SPEC-017`/`SPEC-022`, den ausgelieferten Endpunkt
(`readchanges.go`, `server.go`), das Handbuch und den Review-Report.
**Alle** Zahlen dieses Berichts stammen aus eigenen, in dieser Sitzung
gefahrenen Läufen; kein Beleg des Implementers oder des Reviewers wurde
übernommen. Exit-Codes je in eigenem, ungepiptem Schritt gelesen
(`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung getrennt beauftragt
(`make gates` als eigener Aufruf, `make test` als eigener Aufruf, danach
getrennt ausgewertet). Die Mutations-Proben liefen über eine **Kopie** des
Baums außerhalb des Arbeitsbaums (`git archive` von `HEAD`); der Arbeitsbaum
ist nach allen Läufen unverändert (`git status --short` leer).

**Gegenstand:** `slice-083`, geprüfter Stand `HEAD` = `56f5aee`
(„slice-083 F-1 berichtigt — die drei anderen Clients sind entschieden,
nicht gebaut"). Der Slice liegt in `in-progress/`; der `git mv` nach `done/`
ist **nicht** erfolgt.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (ungepiped, Ausgabe in Logdatei, Exit direkt gelesen) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 701 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) … Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 71.90% erfüllt Schwelle 70%` · `a-check … gesamt: 0 Befund(e)` |
| 2 | `make test` (ungepiped, Exit direkt) | **0** | `ok github.com/pt9912/pg-change-feed/examples/nats-client 1.014s` — das Beispiel liegt real in `go test -race ./...`; kein `FAIL` |
| 3 | `grep -rn "internal/" examples/` und `grep -rn "pg-change-feed/gen" examples/` | — | **leer** — kein `/internal/`- und kein `/gen/`-Import im Beispiel |
| 4 | Nicht-Stdlib-Importe in `examples/` | — | genau ein Treffer: `github.com/nats-io/nats.go` (`main.go:26`) |
| 5 | `git diff --name-only ce84f7f..HEAD` | — | `go.mod`/`go.sum` **nicht** enthalten; `nats.go v1.53.1` steht schon in `f6c74e4^:go.mod` (Zeile 8) |
| 6 | `ls examples/` + `git log --diff-filter=A -- examples/` | — | nur `nats-client`; genau **ein** Commit (`f6c74e4`) hat je eine Datei unter `examples/` angelegt |
| 7 | Suche nach `examples/http-client`/`sse-client`/`grpc-client` in der Planung | — | **keine** Slice-Datei nennt sie; `open/`/`next/` sind leer; die Treffer `slice-061`/`slice-071` sind die E2E-Belegclients unter `tools/harness/**` (`done/`) |
| 8 | `.a-check.yml`-Diff | — | `examples: ["examples/**"]` aufgenommen; **keine** `edges`-Zeile mit `examples` |
| 9 | **Mutation P1** (Kopie) — `os.Getenv("CDC_API_TOKEN_READER")` → `CDC_API_TOKEN_ADMIN` in `parseFlags` | **0** | Suite grün — **kein** Test bindet die Token-Herkunft |
| 10 | **Mutation P2** (Kopie) — der `if cfg.addr == ""`-Guard in `main` entfernt | **0** | Suite grün — **kein** Test bindet die `CDC_HTTP_ADDR`-Pflicht |
| 11 | **Mutation P3** (Kopie, Positiv-Kontrolle) — `url.URL{Scheme: "http"` → `"https"` in `ChangesURL` | **1** | `TestChangesURLCarriesFiltersOnReadPath` **und** `TestChangesURLKeepsHostBoundary` rot — der `Aufbau der Abfrage` **ist** gebunden |
| 12 | **F-2-Gegenprobe** (Kopie) — `ChangesURL("https://feed:8080", "quelle-1", "public", "orders")` | **0** | `= "http://https:%2F%2Ffeed:8080/changes?schema=public&source=quelle-1&table=orders"` — die Funktion setzt ihr Schema **unbedingt** |
| 13 | Beispiel-Binary, `CDC_HTTP_ADDR` ungesetzt (Kopie, `--network none`) | prog. **2** | stderr: `nats-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder -http-addr) ist nötig …` |
| 14 | Beispiel-Binary, nur `CDC_HTTP_ADDR` gesetzt | prog. **2** | stderr-Klartext zu `CDC_NATS_URL` |
| 15 | Beispiel-Binary, `CDC_HTTP_ADDR`+`CDC_NATS_URL`, ohne `-source/-schema/-table` | prog. **2** | stderr: `-source, -schema und -table sind Pflicht` |
| 16 | Beispiel-Binary, alles außer Token | prog. **2** | stderr-Klartext zu `CDC_API_TOKEN_READER` |
| 17 | Beispiel-Binary, NATS nicht erreichbar (`nats://127.0.0.1:1`) | **1** | stderr: `Verbindung zu NATS … fehlgeschlagen: nats: no servers available for connection` |
| 18 | `.github/workflows/**` im Slice-Bereich (`ce84f7f..HEAD`) | — | **0** Dateien, **0** Commits |
| 19 | `Version:`/Historie des Handbuchs | — | `Version: 1.16`; genau **eine** `1.16`-Zeile (`:961`) neben `1.15` (`:960`) |

> Messungen 13–17: `go run` gibt als Wrapper `exit status 2` aus und
> terminiert selbst mit 1; der **Programm**-Exit ist die genannte 2. Der
> Pfad ist vor jedem Verbindungsversuch abgeschlossen (der `nats.Connect`
> folgt erst nach allen vier Guards).

---

## 2. DoD-Konformität, Kriterium für Kriterium

Gelesen ist der Text in der Fassung `56f5aee` (die Datei wurde nach dem
Review geändert — F-1 — geprüft wurde der committete Stand, nicht eine
Erinnerung).

| # | DoD-Zeile (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| LP1-K1 | Eigenständiges CLI (eigenes `main`), nur öffentlicher Draht-Vertrag, **kein** `/internal/`-Import; **es ist das erste Beispiel** — die drei `ADR-0076`-Clients sind entschieden, nicht gebaut, ihre Slices nicht geschnitten | **erfüllt** | Messung 3/4: kein `/internal/`-, kein `/gen/`-Import; `package main` mit eigenem `main` in `main.go`, `subject.go` ist `package main`. Messung 6/7: `examples/` führt nur `nats-client`, genau ein Anlege-Commit, **keine** Slice-Datei baut `http-client`/`sse-client`/`grpc-client`. Der berichtigte Text (F-1) sagt die Wahrheit. |
| LP1-K2 | `github.com/nats-io/nats.go` (bereits Modul-Abhängigkeit) + Standardbibliothek — **keine** neue Abhängigkeit | **erfüllt** | Messung 4: einziger Nicht-Stdlib-Import ist `nats.go`. Messung 5: `go.mod`/`go.sum` sind im Diff-Bereich **nicht** enthalten; `nats.go v1.53.1` stand schon **vor** `f6c74e4` in `go.mod`. |
| LP2-K1 | Das **Lauschen**: das Subjekt wird aus `<source_id>.<schema>.<table>` **abgeleitet**, nicht handgetippt; **keine** Daten aus dem Payload | **erfüllt** | `Subject(sourceID, schema, table)` (`subject.go:11-13`) verkettet die drei Eingaben; `main` ruft es mit den Flag-Werten und abonniert es (`main.go:73,82`). Aus dem Payload kommt kein Datum: die Ausgabe nennt ausschließlich `len(msg.Data)` (`main.go:104-105`); kein Lesezugriff auf `msg.Data`-Inhalt. |
| LP2-K2 | Die **HTTP-Abfrage**: beim Weckruf holt der Client die Änderung **real** über die HTTP-API und gibt sie aus; ungesetztes `CDC_HTTP_ADDR` → **sichtbares Scheitern** | **erfüllt — mit benannter Grenze (Leitung ohne Sensor)** | Der Weckruf (`sub.NextMsg(0)`) löst `fetchChanges` aus (`main.go:96,107`), das `GET /changes` gegen `ChangesURL(cfg.addr, …)` mit `Authorization: Bearer` baut und den Body ausgibt. Endpunkt real vorhanden (Messung: `server.go:102`, Rechtsklasse `reader`). Das sichtbare Scheitern ist real belegt (Messung 13: Klartext auf stderr, Programm-Exit 2). Die **Leitung selbst** (löst ein realer Weckruf die Abfrage gegen einen laufenden Feed aus) hat **keinen** Sensor — `ADR-0076` Festlegung 5 / `ADR-0079` Festlegung 3, siehe §9. |
| LP2-K3 | Die netzlos prüfbare Hälfte — Subjekt-Ableitung und Aufbau der Abfrage — als **reine Funktion** mit eigenen Tests | **erfüllt** | `Subject`/`ChangesURL` sind reine Funktionen (`subject.go`), vier Tests in `subject_test.go`; `make test` grün (Messung 2). Eigene Positiv-Kontrolle (Messung 11): eine Mutation der Eingabeseite (`Scheme` → `https`) macht zwei Tests rot — die Bindung ist real, nicht behauptet. |
| LP3-K1 | `docs/user/benutzerhandbuch.md` trägt einen Abschnitt zum NATS-Wecksignal und **nennt das Beispiel beim Namen**, in der Form der drei bestehenden Zugriffs-Abschnitte | **erfüllt** | Vierter `### Zugriff über das NATS-Wecksignal` (`:723`), gleiche Gliederung (`###`-Kopf, **Erreichbarkeit**, **Subjekt**, **Payload**, **Zustellsemantik**, zweiseitiger Ablauf) wie `:587`/`:649`/`:695`. Nennt `examples/nats-client` mit Startbefehl (`:762-764`); Abgrenzungssatz zu `tools/harness/` (`:765-767`). Alle fünf `ADR-0079`-Festlegung-4-Aussagen vorhanden (Erreichbarkeit samt Asymmetrie, Subjekt-Schema mit beiden Wildcards, leerer Payload, Zustellsemantik, zweiseitiger Ablauf). |
| LP3-K2 | `make gates` grün (Exit direkt, ungepiped) | **erfüllt** | Messung 1: **Exit 0**, alle fünf inneren Gates grün. |
| LP3-K3 | Review durchgeführt, Report unter `docs/reviews/review-slice-083.md` liegt vor; 0 HIGH | **erfüllt** | `review-slice-083.md` (0 HIGH / 1 MEDIUM / 2 LOW / 2 INFO) liegt real vor; das MEDIUM ist F-1 und liegt laut Report im **Plan-Text**, nicht im Code (`:355-363`). F-1 ist in `56f5aee` berichtigt. |
| LP3-K4 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **korrekt offen — Planner** | §7 trägt noch `<…>`-Platzhalter; der `git mv` steht aus. |
| LP3-K5 | Reconciliation-Register fortgeschrieben, falls ein Inventur-Fund aufgelöst wird; **entfällt**, falls die Datei fehlt | **korrekt offen — Planner** | `docs/plan/planning/reconciliation.md` existiert real **nicht** (kein Brownfield-Bootstrap) — das Item entfällt. |
| LP3-K6 | Beobachtungs-Register fortgeschrieben, **kein** Zähler gesetzt | **korrekt offen — Planner** | von der Closure zu tragen (Modul 6); nicht Gegenstand dieses Berichts. |
| LP3-K7 | Jedes Risiko aus §6 trägt **einen** Ausgang (eingetreten / entfallen / weiter offen) | **nicht erfüllt — V-1 (Planner)** | R2/R3/R4 stehen auf `<bei Closure>` (erwartet). **R1 trägt einen veralteten Ausgang**: „*eingetreten* … der Slice geht deshalb **zurück nach `open/`**, nicht nach `done/`" — das widerspricht §3 („der Blocker ist gelöst") und dem tatsächlichen Lifecycle. Siehe §9/V-1. |
| LP3-K8 | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen, hier nicht zuständig** | Das Repo **hat** Wellen (`roadmap.md` liegt in `in-progress/`, `welle-17` ist geschlossen); die Paarungen trägt die nächste Welle-Closure (Modul 6 Schritt 3c, „auch für Slices ohne Wellen-Zugehörigkeit") — die §2-Zeile sagt genau das. |

**Ergebnis §2:** Die drei Liefer-Punkte tragen. Alle **Code- und Doku-Träger**
(LP1-K1/K2, LP2-K1/K2/K3, LP3-K1/K2/K3) sind bestätigt. LP2-K2 ist bestätigt
mit einer benannten, von `ADR-0076`/`ADR-0079` **entschiedenen** Grenze (der
Leitungslauf hat keinen Sensor). **Ein** Kriterium (LP3-K7) ist **nicht**
erfüllt: der Ausgangstext von §6-R1 ist veraltet und widerspricht §3.

---

## 3. Der Nachzug gegen den ausgelieferten Endpunkt

Der Client lief vor `slice-086`; der Vertrag ist nachgezogen. Geprüft an den
**ausgelieferten** Artefakten, nicht am Plan-Text:

| Merkmal | Client | Ausgelieferter Endpunkt | Verdikt |
|---|---|---|---|
| Pfad | `ChangesURL`: `Path: "/changes"` (`subject.go:20`) | `mux.Handle("GET /changes", …)` (`server.go:102`) | **deckungsgleich** |
| Methode | `http.MethodGet` (`main.go:133`) | `GET` | **deckungsgleich** |
| Filter-Parameter | `source`, `schema`, `table` (`subject.go:22-24`) | `SPEC-022`: `source` Pflicht, `schema`/`table` optional (`readchanges.go:23-25`, `:154-157`) | **deckungsgleich** |
| Rechtsklasse | `Authorization: Bearer <CDC_API_TOKEN_READER>` (`main.go:121,137`) | `withToken(cfg.TokenReader, …, roleReader)` (`server.go:102`); `SPEC-022` „`reader` oder `admin`" | **deckungsgleich** |

Zusätzlich: der Client setzt **kein** `from`/`to`/`limit` — beide sind laut
`ADR-0081` §Teilfrage 2 weglassbar, der Client braucht **keinen** Nachzug
(dort ausdrücklich so festgelegt). Die Handbuch-Bump-Hälfte ist form-korrekt:
`Version: 1.15` → `1.16` und **genau eine** neue Historie-Zeile (`:961`), die
den neuen Abschnitt inhaltlich benennt (Messung 19). Der Review hat dieselben
Formen als Negativbefund geführt (`review-slice-083.md` §Negativbefunde,
Antwort 3) — ich habe sie an den Artefakten **selbst** nachvollzogen, nicht
übernommen.

---

## 4. Gegenprobe an einer von beiden Läufen unberührten Zusage

Weder Implementer (falscher Pfad/Präfix) noch Reviewer (Präfix,
`table`→`tables`, Wildcard-Subjekt, `Authorization` entfernt) haben die
folgenden Stellen mutiert; ich habe sie gewählt (Messungen 8–11):

- **Token-Herkunft (`CDC_API_TOKEN_READER`)** — die Zusage „das Token kommt
  aus derselben Variable, die das Handbuch führt". Mutation des Namens auf
  `CDC_API_TOKEN_ADMIN` lässt die **ganze Suite grün** (Messung 9).
- **`CDC_HTTP_ADDR`-Pflicht** — die Zusage „scheitert sichtbar bei
  ungesetztem `CDC_HTTP_ADDR`". Entfernen des Guards lässt die Suite **grün**
  (Messung 10).
- **Positiv-Kontrolle** — Mutation `Scheme: "http"` → `"https"` macht zwei
  Tests **rot** (Messung 11). Die pure-Funktion-Hälfte ist also gebunden.

**Deutung:** Die zwei ungebundenen Zusagen sind **kein DoD-Verstoß** — LP2-K3
verlangt Tests ausdrücklich nur für *Subjekt-Ableitung und Aufbau der
Abfrage*. Beide sind jedoch **nur** über einen Lauf belegbar (die
Konfigurations-Guards führt kein Test aus); der Reviewer hat sie per
Laufzeit-Probe geprüft, und ich habe sie **unabhängig** reproduziert
(Messungen 13–17): der Guard-Pfad ist real sichtbar (stderr, Exit 2), die
Token-Quelle ist am Code und am Klartext von Messung 16 belegt. Das ist die
im nächsten Abschnitt benannte Grenze, keine DoD-Verletzung.

---

## 5. Urteil zu Review-F-2

**F-2 berührt kein DoD-Kriterium — es gehört in die Closure (§7/Register),
nicht in die DoD.**

Der Review beschreibt F-2 (LOW) korrekt: der Kommentar `subject_test.go:37-39`
behauptet, `ChangesURL` setze „kein eigenes Schema auf einen bereits absoluten
Wert"; gemessen setzt sie `Scheme: "http"` **unbedingt**. **Eigene Gegenprobe
(Messung 12):** `ChangesURL("https://feed:8080", …)` = `"http://https:%2F%2Ffeed:8080/changes?…"`
— der Kommentar sagt das Gegenteil dessen, was die Funktion tut
(`AGENTS.md` §3.7: „Ein Kommentar beschreibt, was da ist").

**DoD-Bezug:** LP2-K3 verlangt „Subjekt-Ableitung und Aufbau der Abfrage als
reine Funktion mit eigenen Tests" — beide Funktionen und beide Tests
**existieren** und tragen (Messung 11). Der fehlerhafte Kommentar ändert
weder das Vorhandensein noch die Bindungs-Kraft der Tests; `CDC_HTTP_ADDR`
ist per Handbuch §5 `host:port`, ein absoluter `base` liegt außerhalb des
Vertrags, **kein** Verhaltensfehler. F-2 ist damit eine
Kommentar-Genauigkeits-Klasse (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
ist im Register als **weiter** geführt), die — wie der Review selbst routet —
mit der Slice-Closure §7 in den Zähler geht. **Kein** Rückgabe-Pfeil an den
Implementer nötig.

**Bemerkenswert für die Closure (nicht DoD):** Der Kommentar ist im selben
Zug **falsch und unbemerkt geblieben** — weder Test noch Gate liest
Test-Kommentare. Das ist genau die Verifier-only-Klasse; die Zuordnung zum
Register trifft der Lese-Schritt, nicht dieser Bericht.

---

## 6. `AGENTS.md` §3.10 — greift nicht, bestätigt

Der Slice hat **keine** Datei unter `.github/workflows/**` berührt.
Eigener `git diff --name-only ce84f7f..HEAD -- .github/` ist **leer**, und
`git log --oneline ce84f7f..HEAD -- .github/` liefert **0** Commits
(Messung 18). Die acht Slice-Dateien sind `examples/nats-client/{main,subject,subject_test}.go`,
`.a-check.yml`, `docs/user/benutzerhandbuch.md`, der Slice-Plan und der
Review-Report. §3.10 („neuer oder strukturell geänderter Workflow gilt erst
nach realem Post-Push-Lauf") ist damit **nicht anwendbar** — bestätigt durch
die Namensliste, nicht angenommen.

---

## 7. Plan-vs-Code-Diff

**Zuschnitt (§3):** Der Plan nennt als Träger `examples/nats-client/main.go`
(neu), `subject.go` + `_test.go` (neu), `docs/user/benutzerhandbuch.md`
(update), `.a-check.yml` (update). Der Diff `ce84f7f..HEAD` berührt **genau**
diese Dateien — plus den Slice-Plan und den Review-Report (Vorgangs-Artefakte,
kein Code).

**Die „Nicht in dieser Liste"-Zeile trägt:** `internal/**`, `tools/harness/**`,
`spec/**`, `gen/**`, `go.mod`/`go.sum` und jede neue `.a-check.yml`-**Kante**
sind in **keinem** Commit des Slice enthalten (Messungen 3–8). Insbesondere
ist **keine** Kante gezogen — die `examples`-Gruppe steht **ohne** `edges`-Zeile,
wie §3 und `ADR-0079` Festlegung 1 es verlangen („die Kante erst mit ihrem
Objekt").

**Die `.a-check.yml`-Gruppe ist kein toter Eintrag:** Sie nimmt keine
Import-Berechtigung auf (keine Kante) und lässt `wrong-direction` real stehen;
`ADR-0068` Festlegung 3 (ADR-Pflicht bei Aufnahme **mit** Import-Berechtigung)
ist damit **nicht einmal ausgelöst**, und die Aufnahme ist ohnehin durch
`ADR-0076` Festlegung 4 / `ADR-0079` Festlegung 1 gedeckt. `make a-check` ist
grün (Messung 1).

**Kein stiller Umfangszuwachs:** Der Plan-Nachzug (§3) sagt, nachzuziehen
waren allein zwei Handbuch-Stellen; das deckt sich mit dem Diff — `.a-check.yml`
ist wie geplant **ohne** Kante, das Beispiel ist byte-gleich zum Branch-Stand
(der Review hat das über `git diff 550cff2 f6c74e4` belegt; ich habe es nicht
nachgemessen, weil es für die DoD nicht tragend ist).

---

## 8. Entscheidungs-Konformität

| Entscheidung | Prüfung | Verdikt |
|---|---|---|
| [`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md) Festlegung 1 | `examples/nats-client/`, eigenes `main`, kein `/internal/`-Import, Startform `go run ./examples/nats-client` (Doc-Kommentar `main.go:7`); keine `.a-check.yml`-Kante | **konform** |
| `ADR-0079` Festlegung 2 | Subjekt abgeleitet; keine Payload-Daten; `GET /changes` beim Weckruf; sichtbares Scheitern bei leerem `CDC_HTTP_ADDR`; Adresse/Token aus `CDC_NATS_URL`/`CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`, per Flag übersteuerbar; Quelle/Schema/Tabelle per Flag | **konform** — `parseFlags` (`main.go:117-127`) bildet exakt diese Form ab |
| `ADR-0079` Festlegung 3 | die zwei Entscheidungsschritte als reine Funktionen mit Tests im `make test`-Pfad; **kein** neuer Lauf-Beleg, kein neues Gate; `examples/**` außerhalb der Coverage-Fläche | **konform** — `examples/` erscheint nicht in der `go list`-Zeile des Coverage-Laufs (Messung 1); kein neues Make-Target im Diff |
| `ADR-0079` Festlegung 4 | vierter `### Zugriff über das NATS-Wecksignal`, fünf Aussagen, Beispiel beim Namen, Abgrenzungssatz | **konform** — alle fünf Aussagen und beide Zitate vorhanden (`:723-767`) |
| `ADR-0079` Festlegung 5 | `tools/harness/natssub` bleibt Orakel, unverändert | **konform** — `tools/harness/**` im Diff unberührt |
| `ADR-0079` Festlegung 6 | kein payload-lesendes Beispiel, keine Änderung an den drei, keine zweite NATS-Abhängigkeit, kein neuer Lauf-Beleg, kein Vertragsumzug, keine Coverage-Frage | **konform** |
| [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) Festlegung 2 | nur öffentlicher Draht-Vertrag; kein `/internal/`-Pfad (Grenze: `composition_root` bleibt von a-check ungefasst — benannte Lücke dort) | **konform** |
| `ADR-0076` Festlegung 4 | `examples: ["examples/**"]` ohne Kante; der `contract`-Umzug bleibt dem eigenen Slice | **konform** — die `contract`-Gruppe/`tooling→adapters`-Rücknahme sind als „anderer Vorgang" (§1) ausgeschlossen, `ADR-0079` Festlegung 1 bestätigt das |
| [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) | `GET /changes`, `source`/`schema`/`table`, `reader`-Klasse; „der Client braucht keinen Nachzug" | **konform** (§3 dieses Berichts) |
| [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)/[`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md) + `SPEC-017` | leerer Payload, tabellen-granulares Subjekt; Handbuch-Aussage deckt sich mit `wiring.go:591-612` (gesetzte `CDC_NATS_URL` = Start-Vorbedingung, Klasse `configuration`) | **konform** |
| `AGENTS.md` §3.5 (Accepted-ADRs immutable) | `ADR-0079` supersedet `ADR-0076` in **einer** Klausel (Umfang drei → vier), `ADR-0081` `ADR-0057` teilweise; keine inhaltliche Überschreibung | **konform** |

---

## 9. Befunde dieses Laufs

Verifier-only-Klasse: unsichtbar für Tests und für das Review; betrifft die
Zusage (DoD/Plan-Text), nicht den Diff.

### V-1 — §4 und §6-R1 tragen die überholte Blocker-Erzählung (Adressat: Planner)

- `pfad`: `docs/plan/planning/in-progress/slice-083-nats-beispielclient.md`
  §4 (Rückführung `in-progress → open`, Zeilen 242–250) und §6 (Risiko 1,
  Zeilen 269–277) gegen §3 („Plan-Nachzug dieses Laufs — der Blocker ist
  gelöst", Zeilen 194–207)
- `befund`: §4 sagt im Präsens der Zustandsmaschine „**EINGETRETEN — der
  Grund, mit dem dieser Slice zurückgeht** … **die HTTP-API hat keinen
  Changes-Lese-Endpunkt** … real antwortet der Abruf `404`". §6-R1 schreibt
  als Ausgang: „*eingetreten* … der Slice geht deshalb **zurück nach `open/`**,
  nicht nach `done/`". Beides widerspricht §3 desselben Dokuments, das
  denselben Blocker als **gelöst** führt (`ADR-0081` entschieden, `slice-086`
  geliefert, Beispiel unverändert richtig) und widerspricht dem tatsächlichen
  Lifecycle (der Slice steht in `in-progress/`, der Review ist geschlossen,
  F-1 berichtigt). Das Dokument erzählt zwei Geschichten über denselben
  Vorgang; nur eine kann gelten.
- **DoD-Bezug:** §2-Item LP3-K7 („Jedes Risiko aus §6 trägt einen Ausgang")
  ist deshalb **nicht** erfüllbar, solange R1s Ausgang „zurück nach `open/`"
  lautet — er beschreibt einen Zustand, der nicht (mehr) vorliegt.
- `kategorie`: „Plan-Aussage von der eigenen Nachzug-Aussage widerlegt"
  (Mechanismus der Klasse `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`)
- `verifizierbar`: ja — `git show 56f5aee:…/slice-083-nats-beispielclient.md`
  §3/§4/§6 nebeneinander; kein Gate liest Plan-Prosa.
- `grenze`: **Kein Code-Defekt** — das Beispiel, die Tests, die Gate-Gruppe
  und das Handbuch sind vollständig und richtig. Betroffen ist der
  Plan-Text (Adressat Planner), nicht der Diff.

---

## 10. Was ich strukturell nicht prüfen konnte (benannte Grenzen)

- **Der zweiseitige Ablauf hat keinen Sensor** — das ist die von
  `ADR-0076` Festlegung 5 und `ADR-0079` Festlegung 3 **entschiedene**
  Grenze (Kompilier-Bindung statt Lauf-Beleg). Ich habe geprüft, dass die
  Zusage des Beispiels **am Code** richtig ist (Subjekt-Ableitung, `GET /changes`
  mit den drei Filtern, `len(msg.Data)`) und dass die **Konfigurations-Pfade**
  real sichtbar scheitern (Messungen 13–17). **Nicht** prüfbar ohne neuen
  Belegträger: ob ein **realer Weckruf** eines laufenden Feed-Containers
  tatsächlich die Abfrage auslöst und die abgerufene Änderung ausgibt. Der
  Rest-Fall ist von `ADR-0079` §Re-Evaluierungs-Trigger 1 an einen eigenen
  Smoke-Lauf gebunden.
- **Weist der Slice sie als Grenze aus, statt sie als Beleg zu verkaufen?**
  **Ja.** §1 schließt „ein neuer Lauf-Beleg-Träger" ausdrücklich aus
  (Verweis `ADR-0076`); §5 fordert als Closure-Trigger nur „das Beispiel
  übersetzt im Kompilierpfad" und „die netzlos prüfbare Hälfte läuft real
  grün" — **keine** Behauptung eines Leitungslaufs; §8 nennt die
  „Kompilier-Lücke"; der Doc-Kommentar `main.go:11-15` spricht aus, dass der
  E2E-Belegträger `tools/harness/natssub` ist. Die Grenze ist **benannt**,
  nicht verschwiegen.
- **Die Leitung der zwei ungebundenen Zusagen** (Token-Quelle,
  `CDC_HTTP_ADDR`-Guard): nur über Laufzeit-Proben belegbar (Messungen 9/10
  zeigen, dass kein Test sie bindet; Messungen 13–17 belegen sie real). Das
  ist konsistent mit LP2-K3, das Tests nur für die zwei reinen Funktionen
  verlangt — ich benenne es, statt es zu einem DoD-Verstoß zu machen.
- **Den Post-Push-Lauf** gibt es hier nicht (§6 dieses Berichts); §3.10 ist
  nicht anwendbar, nicht „unbestätigt".
- **Die Vollständigkeit des Handbuch-Abschnitts gegen künftige Drift** —
  geprüft ist der gelieferte Stand, nicht ob der Abschnitt bei einer
  künftigen NATS-Änderung mitzieht.

---

## 11. Verdikt

**Der Slice ist in der Sache closure-fähig — mit einer Plan-Textkorrektur,
die vor dem `git mv` nach `done/` stehen sollte.**

Getragen ist alles, was der Slice zusagt:

- Das Beispiel ist das **erste** `examples/`-Programm, eigenständiges
  `main`, ausschließlich öffentlicher Draht-Vertrag (kein `/internal/`,
  kein `/gen/`), `nats.go` + stdlib, **keine** neue Abhängigkeit — die
  berichtigte LP1-Aussage (F-1) ist am Dateisystem und an `git` bestätigt.
- Der zweiseitige Ablauf ist **am Code** vollständig: abgeleitetes
  tabellen-granulares Subjekt, keine Payload-Daten, `GET /changes` mit den
  drei Filtern beim Weckruf, sichtbares Scheitern bei ungesetztem
  `CDC_HTTP_ADDR` (eigene Laufzeit-Probe, Exit 2, Klartext auf stderr).
- Die netzlos prüfbare Hälfte ist reine Funktion mit eigenen Tests; die
  Bindung ist mit einer **eigenen** Positiv-Kontrolle bestätigt (nicht nur
  behauptet).
- Der Vertrag ist gegen den **ausgelieferten** Endpunkt nachgezogen (Pfad,
  Parameter, Rechtsklasse, Token-Quelle), der Handbuch-Bump ist form-korrekt
  (`1.16`, genau eine Historie-Zeile).
- `make gates` und `make test` grün (eigene, ungepipete Läufe); die
  `.a-check.yml`-Gruppe ist da und **ohne** Kante; die im Slice
  referenzierten ADRs sind konform (§8).
- §3.10 ist **nicht anwendbar** — kein Workflow berührt (bestätigt, nicht
  angenommen).

**Was fehlt:** §4 und §6-R1 tragen die **überholte** Blocker-Erzählung
(„EINGETRETEN … zurück nach `open/`"), die §3 desselben Dokuments widerlegt
(V-1). Damit ist das DoD-Item LP3-K7 (Risiko-Ausgänge) derzeit **nicht**
erfüllbar. Das ist keine Kosmetik: der Text ist der Träger der
Lifecycle-Aussage und des Risiko-Ausgangs; er gehört in den Plan-Text,
bevor er nach `done/` wandert. Der **Adressat ist der Planner** — der
Verifier repariert nicht, er berichtet. Die übrigen Planner-Posten
(K4–K8, `git mv`, Paarungen) sind im Text geführt bzw. der nächsten
Welle-Closure zugewiesen.

**Nicht closure-blockierend** und ausdrücklich als getragen gewertet: der
Diff und die Handbuch-Hälfte; die benannte Grenze des Leitungslaufs (§10
dieses Berichts); F-2 als Kommentar-Genauigkeit (gehört in §7/Register, kein
DoD-Kriterium, kein Rückgabe-Pfeil).

---

**Übergabe:** Dieser Bericht geht an den **Planner** (V-1: §4/§6-R1
Plan-Text; die Closure-Entscheidung), an den **Implementer** als Entlastung
(der Diff ist nachweislich tragend, kein Rückgabe-Pfeil) und an den
**Validator** nur, falls dieser Slice als MVP-Slice validiert wird — der
reale Bedarf liegt hier repo-extern nicht vor. Der Bericht ist ein
**Lauf-Beleg** (dieser Stand `56f5aee`, diese Läufe, dieses Verdikt) und
ersetzt weder das Review noch die Planner-Closure.
