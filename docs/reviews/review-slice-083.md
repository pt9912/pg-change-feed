# Review-Report: slice-083 — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier, Modul 11),
§6-Risiko-Ausgänge, Beobachtungs-Register und die drei Paarungen
(Planner-Closure) sind **nicht** Gegenstand dieses Reports.

**Gegenstand:** `slice-083`, Diff `ce84f7f..e02b00a` — vier Commits:
`f6c74e4` (Beispielclient, `.a-check.yml`-Gruppe), `095bda2`
(Handbuch-Abschnitt), `1707dab` (Handbuch-Nachzug: `GET /changes` und
Versionshistorie), `e02b00a` (Plan-Nachzug und DoD-Häkchen). Sechs Dateien:
`examples/nats-client/{main,subject,subject_test}.go` (neu),
`.a-check.yml`, `docs/user/benutzerhandbuch.md`, der Slice-Plan. Kein
`internal/**`, kein `spec/**`, kein `go.mod`/`go.sum`, kein
`Makefile`/`harness/mk/**`.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-083-nats-beispielclient.md`
  vollständig (§1–§8), einschließlich des §3-Nachtrags in `e02b00a`
- [`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
  (der vierte Client, die sechs Festlegungen),
  [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  (Form der Beispiele, Kompilier-Bindung, Gate-Scope-Gruppen),
  [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  (der Changes-Lese-Endpunkt),
  [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)
  (leerer Payload), [`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md)
  (tabellen-granulares Subjekt),
  [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  Festlegung 3 (ADR-Pflicht bei Aufnahme mit Import-Berechtigung),
  [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
- [`SPEC-017`](../../spec/pflichtenheft.md) (Subjekt- und Nachrichtenform),
  [`SPEC-018`](../../spec/pflichtenheft.md) (Authn-Header, Aktivierung),
  [`SPEC-022`](../../spec/pflichtenheft.md) (`GET /changes`),
  [`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-SST-007`](../../spec/lastenheft.md)
- `internal/adapters/driving/http/server.go`, `.../readchanges.go` (der
  ausgelieferte Endpunkt), `tools/harness/natssub` (der E2E-Belegträger)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Ist-Zustand statt Chronik), §3.9
  (Exit-Code und Folgehandlung), §3.11 (host-lokale Pfade), §5
- `harness/conventions.md` (MR-000 ID-Schema), `.harness/skills/reviewer.md`
- Vorherige Läufe am gleichen Gegenstand: `docs/reviews/review-slice-086.md`
  (F-2, derselbe Mechanismus), `docs/reviews/review-slice-087.md` (F-1),
  `docs/reviews/review-slice-082.md` (F-1, Plan-Tatsachenbehauptung)

---

## Eigene Messungen dieses Laufs

Alle Läufe sind **eigene** Läufe (Docker-only, Exit-Code ungepiped gelesen,
`AGENTS.md` §3.9); die Mutationen und die `a-check`-Proben liefen über eine
**Kopie** des Baums außerhalb des Arbeitsbaums (`git archive` des
geprüften Stands), der Arbeitsbaum ist nach allen Läufen unverändert.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Stand `e02b00a`) | **0** | `d-check: 700 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK`; `coverage-gate: OK — Coverage 72.00% erfüllt Schwelle 70%`; `a-check … gesamt: 0 Befund(e)` |
| 2 | `make test` | **0** | `ok …/examples/nats-client 1.014s` — das Beispiel liegt real in `go test -race ./...` |
| 3 | **Mutation A** — `Subject`-Präfix `cdc.changes.` → `cdc.change.`, nur `TestSubjectKeepsDistinctTables` | **0** | dieser eine Test bleibt grün |
| 4 | **Mutation A** — dieselbe Änderung, ganze Paket-Suite | **1** | `TestSubjectDerivesTableGranularSubject` rot (`Subject = "cdc.change.quelle-1.public.orders"`) |
| 5 | **Mutation B** — Query-Parameter `table` → `tables` | **1** | beide `TestChangesURL*` rot |
| 6 | **Mutation C** — Wildcard-Form des Subjekts (`cdc.changes.<source_id>.>`) | **1** | `TestSubjectDerives…` **und** `TestSubjectKeepsDistinctTables` rot |
| 7 | **Mutation D** — `Authorization`-Header aus `fetchChanges` entfernt | **0** | Suite grün — die von [`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md) Festlegung 3 **benannte** Grenze |
| 8 | Probe: `examples/probe` importiert `internal/adapters/driving/http` | **1** | `wrong-direction: examples -> adapters` |
| 9 | Probe: `examples/probe3` importiert `internal/domain/model` | **1** | `wrong-direction: examples -> domain` |
| 10 | Probe: `examples/probe2` importiert `internal/bootstrap` | **0** | grün — die benannte Grenze aus [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) Festlegung 2 |
| 11 | Probe: `tools/probe4` **ohne** `layers`-Eintrag importiert `internal/adapters/…` | **1** | `wrong-direction: (ohne Schicht) -> adapters` **plus** Abdeckungs-Hinweis „liegen in keiner Schicht und bleiben ungeprüft" |
| 12 | Beispiel-Binary, `CDC_HTTP_ADDR` ungesetzt | **2** | `nats-client: keine HTTP-Adresse gesetzt — CDC_HTTP_ADDR (oder -http-addr) ist nötig …` auf stderr |
| 13 | Beispiel-Binary, nur `CDC_HTTP_ADDR` gesetzt | **2** | Abbruch mit Klartext zu `CDC_NATS_URL` |
| 14 | Beispiel-Binary, alles außer Token | **2** | Abbruch mit Klartext zu `CDC_API_TOKEN_READER` |
| 15 | Beispiel-Binary, NATS nicht erreichbar (`nats://127.0.0.1:1`) | **1** | `Verbindung zu NATS … fehlgeschlagen: nats: no servers available for connection` |
| 16 | `gofmt -l examples/` | **0** | leer |
| 17 | Rebase-Nachweis: `git diff 550cff2 f6c74e4 -- examples/ .a-check.yml` | **0** | leer — Beispiel und Gruppe **byte-gleich** zum Vor-Rebase-Stand der Branch-Spitze `7314669` |
| 18 | Vor-Rebase-Stand `7314669`: `Version:` des Handbuchs | — | `1.14` — die Branch-Fassung trug **weder** Bump **noch** Historie-Zeile |

---

## Findings

### F-1 — Das DoD-Kriterium vergleicht mit drei Beispiel-Clients, die es nicht gibt

- `kategorie`: **MEDIUM**
- `quelle`: Maintainability · `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice (das
  DoD-Kriterium ist die Prüf-Form des Liefer-Punkts)
- `pfad`: `docs/plan/planning/in-progress/slice-083-nats-beispielclient.md:102-104`
  (§2, Liefer-Punkt 1, erstes Kriterium) · `:44-53` (§1 Ziel) · `:12-16`
  (§Bezug) · `:61-64` (§1, zweiter Ausschluss)
- `befund`: Das in `e02b00a` gesetzte Häkchen wird mit „eigenes `main`, **wie
  die drei anderen**" begründet; gemeint sind `examples/http-client`,
  `examples/sse-client` und `examples/grpc-client` aus
  [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §Entscheidung Festlegung 1. Sie existieren nicht: `examples/` führt genau
  `nats-client`, und `git log --diff-filter=A -- examples/` zeigt für den
  ganzen Baum **einen** Commit (`f6c74e4`) — dieses Beispiel **ist** der
  erste. Kein Slice in `open/`, `next/`, `in-progress/` oder `done/` baut
  die drei; sie stehen nur im ADR-Text. Der Vergleich hat damit keinen
  Bezugspunkt, und dieselbe Formulierung trägt §1 („das einzige der
  Beispiele") und der §Bezug („in unveränderter Form der drei anderen").
- `verifizierbar`: ja — Dateisystem und `git log --diff-filter=A -- examples/`;
  kein Gate prüft es
- `klasse`: „DoD-Kriterium beruft sich auf nicht existierende Nachbarartefakte"
  (Mechanismus der Klasse
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`, 3×, Ausgang
  offen — Zuordnung entscheidet der Lese-Schritt, nicht dieser Report)
- `hinweis`: Der **Code** ist davon nicht betroffen; die Kriterien-Substanz
  (eigenes `main`, kein `/internal/`-Import) ist real erfüllt (Messungen 1/2
  und `grep`). Der Defekt liegt im Plan — Träger ist der Planner, nicht der
  Implementer.

### F-2 — Der Grenz-Kommentar über `TestChangesURLKeepsHostBoundary` beansprucht eine Aussage, die ihr Input nicht ausübt

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 (ein Kommentar beschreibt, was
  da ist) · [`SPEC-018`](../../spec/pflichtenheft.md) (Form der Horch-Adresse)
- `pfad`: `examples/nats-client/subject_test.go:37-39` (`subject.go:17-26`
  betroffen)
- `befund`: Der Kommentar sagt: „die Adresse ist der Horch-Host aus
  `CDC_HTTP_ADDR`; die Funktion setzt kein eigenes Schema auf einen bereits
  absoluten Wert." `ChangesURL` setzt `Scheme: "http"` **unbedingt**; gemessen
  liefert `ChangesURL("https://feed:8080", …)` den Wert
  `http://https:%2F%2Ffeed:8080/changes?…` — die Funktion **setzt** ihr
  Schema auf einen absoluten Wert auf, sie repariert ihn nicht. Der Test übt
  diesen Fall nicht aus: sein einziger Input ist `localhost:9090`, ein
  schema-loser Host. Die Aussage hat damit keinen Träger — weder im Code noch
  in der Assertion.
- `verifizierbar`: nein — kein Gate liest Test-Kommentare; die Gegenprobe ist
  ein Aufruf der reinen Funktion mit einem absoluten `base`
- `klasse`: „Aussage nicht an ihre Eingabeseite gebunden (Träger: Kommentar)"
  — Mechanismus der Klasse
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (2×, offen). Der
  Register-Eintrag ist laut seinem `state.md` ausdrücklich weiter als sein
  Name: er zählt den Mechanismus, nicht die Negativtest-Form.
- `grenze`: Es ist **kein** Verhaltensfehler — `CDC_HTTP_ADDR` ist per
  Handbuch §5 `host:port`, ein absoluter Wert liegt außerhalb des Vertrags.
  Betroffen ist die Aussage des Kommentars, nicht die Funktion.

### F-3 — Das Handbuch wird in `095bda2` inhaltlich geändert, ohne `Version:`-Bump und ohne Historie-Zeile im selben Commit

- `kategorie`: LOW
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation, HIGH-Punkt
  „Handbuch-Versionshistorie nicht fortgeschrieben" — hier **commit**-skopiert,
  nicht slice-skopiert (siehe `grenze`)
- `pfad`: Commit `095bda2` (`docs/user/benutzerhandbuch.md`, 44 hinzugefügte
  Zeilen; `Version: 1.15` unverändert, `### Änderungshistorie` ohne neue
  Zeile) — der Bump `1.15` → `1.16` samt Zeile kommt erst in `1707dab`
- `befund`: Der Abschnitt `### Zugriff über das NATS-Wecksignal` landet in
  einem Commit, der den Versionskopf und die Änderungshistorie nicht
  fortschreibt; beide werden erst im **folgenden** Commit desselben Laufs
  nachgetragen. Dieselbe Form trug schon die Vor-Rebase-Fassung der Branch-
  Spitze `7314669` (`Version: 1.14`, keine Zeile) — der Rebase hat sie nicht
  geheilt, erst der Nachzug-Commit. Die Selbstprüfung des Implementer-Laufs
  hat die Klasse im ersten Commit also nicht gefeuert; sichtbar wird die
  Stelle über die Nachweistechnik der Klasse selbst
  (`git log -p -- docs/user/benutzerhandbuch.md`).
- `verifizierbar`: ja — `git show 095bda2:docs/user/benutzerhandbuch.md`
  (Kopfzeile und Tabellenende)
- `klasse`: „Handbuch-Versionshistorie übersprungen (commit-skopiert)"
- `grenze`: **Der ausgelieferte Stand ist form-korrekt.** Der Diff-Bereich als
  Ganzes trägt `Version: 1.16` und **genau eine** neue Zeile `1.16`, die den
  Abschnitt inhaltlich benennt (Handbuch §Anhang). Der HIGH-Punkt des Skills
  greift gegen den *Slice*-Diff deshalb **nicht**; benannt ist die
  Commit-Form, damit der Lese-Schritt entscheiden kann, ob sie in den Zähler
  gehört.

### F-4 — Verwaister lokaler Branch `slice-081-executor-naht` neben geschlossenem `slice-081`

- `kategorie`: INFO
- `quelle`: Maintainability (Hygiene des Arbeitsbaums)
- `pfad`: `git branch -a` (lokale Refs außerhalb des Diffs)
- `befund`: Der Branch des geprüften Slice ist **weg** (`slice-083-nats-beispielclient`
  existiert nur noch im Reflog); daneben führt der Klon einen lokalen Branch
  `slice-081-executor-naht`, dessen Vorgang mit `done/slice-081-executor-naht.md`
  abgeschlossen ist. Er liegt außerhalb dieses Diffs und außerhalb des
  Gegenstands.
- `verifizierbar`: ja — `git branch -a`
- `klasse`: „verwaister lokaler Branch nach Slice-Closure"

### F-5 — Der Selbst-Befund des Implementers trifft die Klasse nicht

- `kategorie`: INFO
- `quelle`: Maintainability · `.harness/skills/reviewer.md` §Klassifikation
- `pfad`: `examples/nats-client/subject_test.go:17-26`
- `befund`: Der Implementer hat `TestSubjectKeepsDistinctTables` als „nicht
  diskriminierend" benannt, weil er bei einer falschen Präfix-Mutation grün
  blieb. **Unabhängig nachgemessen und nicht bestätigt:** die Mutation macht
  die **ganze Paket-Suite** rot (Messung 4) — `TestSubjectDerivesTableGranularSubject`
  pinnt die vollständige Zeichenkette `cdc.changes.quelle-1.public.orders`
  mit drei verschiedenen Werten und trägt damit die Präfix-Aussage. Der
  Nachbartest bleibt nur **in Isolation** grün; sein eigener Anspruch
  („zwei Tabellen derselben Quelle ⇒ zwei verschiedene Subjekte") ist an
  seine Eingabe gebunden (Messung 6 macht ihn bei einer die Tabelle
  verschluckenden Mutation rot). Die Klasse ist hier **nicht** ausgelöst: die
  Aussage ist grün **mit** Träger, nicht grün ohne Aussage.
- `verifizierbar`: ja — Mutationsläufe 3, 4 und 6 dieses Reports
- `klasse`: „Selbstbefund des Implementers, am Gegenstand nicht bestätigt"

---

## Negativbefunde

- geprüft, ohne Befund: **öffentlicher Draht-Vertrag** — `grep` über
  `examples/` liefert keinen `/internal/`- und keinen `/gen/`-Import; die
  einzigen Nicht-Stdlib-Importe sind `github.com/nats-io/nats.go` (Messung:
  `go.mod:8`, `v1.53.1`, **vor** diesem Diff). `go.mod`/`go.sum` sind im
  Diff-Bereich unverändert — die Zusage „keine neue Abhängigkeit" trägt.
- geprüft, ohne Befund: **Vertragstreue des Clients gegen den ausgelieferten
  Endpunkt** — Pfad `/changes` (`server.go:102`), Query-Parameter
  `source`/`schema`/`table` (Rechtsklasse `reader`), Header
  `Authorization: Bearer <token>` mit `CDC_API_TOKEN_READER` stimmen mit
  [`SPEC-022`](../../spec/pflichtenheft.md), [`SPEC-018`](../../spec/pflichtenheft.md)
  und `readchanges.go:22-43`, `:132-162` überein. Die Aussage des Vorgänger-
  Laufs („unverändert lauffähig") ist bestätigt: der Branch-Stand ist
  **byte-gleich** übernommen (Messung 17).
- geprüft, ohne Befund: **die fünf Aussagen des Handbuch-Abschnitts**
  ([`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
  Festlegung 4) — Erreichbarkeit samt Asymmetrie zu `CDC_HTTP_ADDR`, beide
  Wildcards, leerer Payload, Zustellsemantik/Nachvollziehbarkeit, zweiseitiger
  Ablauf; jede deckt sich mit [`SPEC-017`](../../spec/pflichtenheft.md) und
  der §5-Zeile `CDC_NATS_URL`. Der Abschnitt nennt das Beispiel mit Pfad und
  Startbefehl, steht in §4 in der Form der drei bestehenden
  `### Zugriff über …`-Abschnitte (`:587`, `:649`, `:695`) und trägt den
  Abgrenzungssatz zu `tools/harness/`. **Kein** neues `CDC_*`-Feld — die
  Klasse „neue Betreiber-Oberfläche ohne Handbuch-Zug" ist nicht ausgelöst.
- geprüft, ohne Befund: **der zweiseitige Ablauf und die `CDC_HTTP_ADDR`-Zusage**
  — Lauschen auf das **abgeleitete** Subjekt (`Subject`), keine Payload-Lesung
  (die Ausgabe nennt nur `len(msg.Data)`), Abfrage über `ChangesURL`. Bei
  ungesetztem `CDC_HTTP_ADDR` bricht das Programm sichtbar ab (Messung 12,
  Exit 2, Klartext auf stderr) — die §2-Zusage trägt.
- geprüft, ohne Befund: **`.a-check.yml`** — die Gruppe
  `examples: ["examples/**"]` ist **ohne** Kante aufgenommen und `make a-check`
  bleibt grün (Messung 1). Sie ist kein toter Eintrag: ein `examples`-Import in
  `adapters` **und** in `domain` ist `wrong-direction` (Messungen 8/9), und
  eine Gruppe ohne Kante nimmt auch **keine** Berechtigung auf — ein Pfad ohne
  jeden `layers`-Eintrag wird weiterhin als `(ohne Schicht)` befundet, nur
  zusätzlich als „ungeprüft" ausgewiesen (Messung 11). Nach
  [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  Festlegung 3 ist die ADR-Pflicht an eine Aufnahme **mit Import-Berechtigung**
  gebunden — diese Hälfte liegt hier nicht vor; gedeckt ist die Aufnahme
  ohnehin durch [`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  Festlegung 4 und [`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
  Festlegung 1.
- geprüft, ohne Befund: **die netzlos prüfbare Hälfte** — `Subject` und
  `ChangesURL` sind reine Funktionen mit vier Tests; drei unabhängige
  Mutationen der Eingabeseite (Präfix, Parameter-Name, Wildcard-Form) enden
  rot (Messungen 4–6), die Suite übersetzt das Beispiel mit (Messung 2).
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7)** — kein
  `slice-`/`welle-`-Verweis und keine Vorher/Nachher-Sprache in den
  hinzugefügten Go-Zeilen; die Herkunfts-Anker im Handbuch (`slice-083` in der
  Historie-Zeile, `ADR-0056` in der `CDC_NATS_URL`-Zeile) stehen in der
  zulässigen Form.
- geprüft, ohne Befund: **Hygiene** — kein `//nolint`/`#noqa`/`[SuppressMessage]`
  im Diff; `gofmt -l examples/` leer; keine Lauf-Artefakte (`image-hash.txt`,
  `plan.yaml`, `down.sql`) in einem der vier Commits; keine Vorlagen-Reste und
  keine unaufgelösten `<…>`-Platzhalter in den neuen Dateien (die `<…>` in
  `main.go:3` und `subject.go:8` sind das Subjekt-Schema, kein Platzhalter);
  kein host-lokaler absoluter Pfad.
- geprüft, ohne Befund: **Traceability** — alle vier Betreffs nennen
  `LH-*`/`ADR-*` und keine Struktur-ID; kein Trailer in einem der vier
  Commit-Texte; `make commit-traceability` grün (Messung 1). `main` ist
  **linear** geblieben: `git log --merges` ist leer, der Reflog führt
  `merge slice-083-nats-beispielclient: Fast-forward`, und
  `slice-083-nats-beispielclient` ist als Ref verschwunden.

---

## Antworten auf die Schwerpunkte

**1 — Ist das Beispiel ein Vorbild?** Ja. Nur öffentlicher Draht-Vertrag:
`flag`, `fmt`, `io`, `net/http`, `os`, `time`, `net/url`, `testing` und
`github.com/nats-io/nats.go` — **kein** `/internal/`-Import (Negativbefund
oben), **keine** neue Abhängigkeit (`go.mod` im Diff unverändert). Die
Aufnahme der `examples`-Gruppe ist keine stille Berechtigung: sie nimmt keine
Kante auf und lässt `wrong-direction` real stehen (Messungen 8/9).

**2 — Trägt der zweiseitige Ablauf?** Ja. Das Subjekt wird abgeleitet
(`Subject`), nicht getippt; aus dem Payload kommt **kein** Datum (nur seine
Länge wird ausgegeben); der Weckruf löst `GET /changes` mit den drei Filtern
aus dem Subjekt aus. Die Zusage aus §2 „scheitert sichtbar bei ungesetztem
`CDC_HTTP_ADDR`" ist real belegt (Messung 12: Exit 2, Klartext auf stderr),
ebenso die drei übrigen Konfigurationspfade (13–15). Was **keinen** Sensor
hat, ist die Leitung selbst — genau die von
[`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
Festlegung 3 benannte Grenze; die Prüfung dieser Hälfte lag bei diesem Lauf
(Lesen des Codes plus Laufzeitproben) und ist grün.

**3 — Stimmt der Nachzug gegen den ausgelieferten Endpunkt?** Ja, in beiden
Hälften. (a) Pfad, Parameter und Rechtsklasse des Clients decken sich mit
[`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md),
[`SPEC-022`](../../spec/pflichtenheft.md) und `server.go`/`readchanges.go`
(Negativbefund oben); die Behauptung „unverändert übernommen" ist über den
Reflog und `git diff 550cff2 f6c74e4 -- examples/ .a-check.yml` **byte-genau**
bestätigt (Messung 17) — nicht nur plausibel. (b) Der Handbuch-Bump ist
form-korrekt: `Version: 1.15` → `1.16`, genau **eine** neue Zeile `1.16` in
`### Änderungshistorie`, die den neuen Abschnitt inhaltlich benennt. Die
Klasse `BEO-PGC/handbuch-versionshistorie-uebersprungen` ist im Register als
**verkörpert** geführt und für den gelieferten Stand **nicht** ausgelöst; die
commit-skopierte Abweichung ist F-3.

**4 — Ist `TestSubjectKeepsDistinctTables` ein Finding?** **Nein** — der
Verdacht ist am Gegenstand widerlegt (F-5 mit Messungen 3, 4, 6). Der Test
ist in Isolation schwach, aber die Aussage, die er trägt, ist eine
**relationale** und an ihre Eingabe gebunden; die Präfix-Aussage trägt der
Nachbartest, und die Mutation des Implementers macht die Suite rot. Die
Klasse `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` ist an dieser
Stelle **nicht** ausgelöst. Ein Finding aus derselben Klasse steht aber
daneben, an einem anderen Gegenstand: **F-2**, wo eine Aussage des Kommentars
weder im Code noch in der Assertion einen Träger hat.

**5 — Ist die `.a-check.yml`-Aufnahme gedeckt?** Ja — und sie ist enger, als
die Regel vermutet. Sie erfüllt die ADR-Pflicht nach
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
Festlegung 3 nicht nur, sie **löst sie nicht einmal aus**: eine Gruppe **ohne**
Kante vergibt keine Import-Berechtigung (Messungen 8/9/11), sie verschiebt die
Dateien von „ohne Schicht, ungeprüft" nach „in der Gruppe `examples`,
befundet". Träger der Aufnahme ist
[`ADR-0076`](../plan/adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
Festlegung 4, wiederholt in
[`ADR-0079`](../plan/adr/0079-nats-beispielclient-vierter-examples-client.md)
Festlegung 1 („die Kante erst mit ihrem Objekt"). `make a-check` ist grün
(Messung 1); die neue Kante ist korrekt **nicht** gezogen.

**6 — Hygiene.** Keine Lauf-Artefakte in einem Commit; keine Vorlagen-Reste
in den neuen Dateien; alle im Plan genannten Kennungen lösen auf (`slice-081`,
`slice-086`, `slice-087` liegen in `done/`); keine host-lokalen Pfade; die
vier Betreffs tragen `LH-*`/`ADR-*` und **keine** `SPEC-*`/`ARC-*`, kein
Trailer. `main` ist linear geblieben (Fast-Forward, `git log --merges` leer),
und `slice-083-nats-beispielclient` ist weg — nur der Reflog führt ihn noch.
Der verwaiste Nachbarbranch `slice-081-executor-naht` steht außerhalb dieses
Gegenstands (F-4).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „DoD-Kriterium beruft sich auf nicht
existierende Nachbarartefakte" · „Aussage nicht an ihre Eingabeseite gebunden
(Träger: Kommentar)" · „Handbuch-Versionshistorie übersprungen
(commit-skopiert)" · „verwaister lokaler Branch nach Slice-Closure" ·
„Selbstbefund des Implementers, am Gegenstand nicht bestätigt"

## Verdikt

**Merge-blockierend:** nein — **0 HIGH**. F-1 ist MEDIUM, liegt aber nicht im
Code: das Beispiel, die Tests, die Gate-Gruppe und das Handbuch sind
vollständig und richtig, und die überzeichnete Vergleichsaussage steht im
**Plan** (Planner). Das entspricht dem Muster von
`docs/reviews/review-slice-082.md` F-1. F-2 und F-3 sind LOW, F-4/F-5 sind
INFO; keines davon ist ein Verhaltens- oder Vertragsfehler, und keines verlangt
einen zweiten Implementer-Lauf. Nicht still entschieden: F-1 braucht einen
Ausgang (Plan-Korrektur oder bewusste Annahme), bevor der Slice nach `done/`
geht.

**Übergabe:** **kein Rückgabe-Pfeil** Reviewer → Implementer. F-1 geht als
Rückkante **Review → Plan** an den Planner (Plan-Defekt, Modul 10
§Übergabe-Artefakt dieses Reports); wäre die drei-Beispiele-Referenz
bedeutungstragend, wäre der Träger dort. F-2 und F-3 reisen mit der
Slice-Closure §7, ohne Rückkante; F-4 liegt außerhalb des Gegenstands. Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler — die Zuordnung zu `BEO-<KUERZEL>` trifft der Lese-Schritt, nicht
dieser Report. Der Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff,
dieser Skill, dieses Modell, dieses Verdikt) und ersetzt keine Verifikation
(Modul 11) — die Prüfung der Liefer-Punkte gegen die DoD bleibt beim Verifier,
die §6-Risiko-Ausgänge und die drei Paarungen bei der Planner-Closure.

**DoD-Häkchen „Review durchgeführt":** wird **in diesem Commit** auf `[x]`
nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde), mit Verweis auf diesen Report — kein HIGH, und kein MEDIUM/LOW
verlangt einen zweiten Implementer-Lauf. Nur diese eine Zeile des Plans ist
geändert; DoD-Substanz und §7 bleiben offen (Verifier/Planner).
