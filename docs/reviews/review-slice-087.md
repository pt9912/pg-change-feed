# Review-Report: slice-087 — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `slice-087`, Commit `93ac92e` (Vorgänger `5d9ef32`) — fünf
Dateien: `tools/harness/httpclient/main.go`,
`tools/harness/run-integration-tests.sh`, `docs/user/e2e-abdeckung.md`,
`harness/README.md`, der Slice-Plan selbst. Ein Commit, kein
`internal/**`, kein `spec/**`, kein `tools/schema/**`, keine
`.a-check.yml`-/`Makefile`-Änderung. Eingelöst wird die Fitness-Function-Zeile
`make test-integration` aus [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
§Fitness Function.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-087-changes-e2e-beleg.md`
  vollständig (§1–§8)
- [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  (Gegenstand: Fitness Function, Teilfragen 2 und 4),
  [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) (der bestehende
  HTTP-Rundlauf), [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  (begrenzte Import-Berechtigung der Wegwerf-Clients),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md) (E2E-Tier),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Digest-Beleg),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (View-Direktzugriff), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (Schema-Rollout des Laufs)
- [`LH-FA-SST-006`](../../spec/lastenheft.md) (Hauptbezug),
  [`LH-FA-REA-001`](../../spec/lastenheft.md)…`006` (Bereichs-/Ordnungs-/
  Limit-Zusagen des Leseports), [`LH-QA-POR-003`](../../spec/lastenheft.md)
- `internal/application/port/inbound/readchanges.go` (Start inklusiv / End
  exklusiv), `internal/adapters/driven/postgresstorage/queries/queries.go`
  (`commit_position >= $2 AND < $3`), `tools/schema/schema.yaml`
  (`cdc.changes`-Projektion), `internal/adapters/driving/http/readchanges.go`
  (Response-Form)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.11, §5 · `harness/conventions.md`
  (MR-000 ID-Schema) · `harness/README.md` §Sensors
- Vorherige Läufe am gleichen Modul: `docs/reviews/review-slice-086.md`
  (F-1 „Fitness-Function-Zusage ohne Träger" — der Anlass dieses Slice),
  `review-slice-081.md`, `review-slice-082.md`

---

## Eigene Messungen dieses Laufs

Docker-only; jeder Exit-Code **ungepiped** und in einem eigenen Schritt
gelesen (`AGENTS.md` §3.9). Nach jedem Lauf `git status --porcelain` geprüft;
die Mutationen wurden per `cp` aus einer Kopie des Originals zurückgenommen.
Der Baum ist vor diesem Commit leer.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` | **0** | fünf innere Gates grün; `coverage-gate: OK — 72.00% erfüllt Schwelle 70%`; `d-check: 697 Dateien, 0 Befunde`; `a-check: gesamt 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s)`, „Betreffs ohne Struktur-ID" |
| 2 | `make test-integration` (unmutiert) | **0** | 3m06s. `READ changes=1 table=feed_e2e_full schema=public change_id=965-1 operation=INSERT commit_position=30901328 new_image={"id":"285","name":"HttpChangesReadE2ESentinel"}`; Bereich `[30901328,30901329)`; danach `git status --porcelain` **leer** |
| 3 | **Mutation A** — Client endet still mit `os.Exit(1)` vor jedem Request (`main.go`) | **2** | `… HTTP-API-Rundlauf (httpclient) endete mit Ausgang 1: exit status 1` — die Phase endet rot **mit** Meldung; der Ausgang wird ausgewertet, nicht verschluckt |
| 4 | **Mutation B** — Client sendet `table=…_gegenprobe` (nicht treffender Filter, Endpunkt antwortet `200`) | **2** | `httpclient: GET /changes (reader) fehlgeschlagen: leere Changes-Liste für den erwarteten Bestand (public.feed_e2e_full, Bereich [30901520,30901521)): {"changes":[]}` — rot **aus der Inhaltsprüfung**, nicht aus dem HTTP-Status |
| 5 | **Mutation C** — Client lässt die Filterachse (`schema`/`table`) aus der Abfrage weg | **0** | **grün**, `READ changes=1 table=feed_e2e_full schema=public … change_id=965-1` — der Beleg unterscheidet „Filter wirkt" nicht von „Filter fehlt" → **F-1** |
| 6 | Anker-Nachrechnung (24 Bash-Zeilen) | — | eigener `awk`-Nachbau der Runner-Logik (`index($0,anker)` ab `NR > Deklarationszeile`, Backtick-Escapes entfernt): **24/24** `Ort`-Werte der Tabelle stimmen mit der aufgelösten Zeile überein, inkl. `HTTP-API-Rundlauf → 2061` |
| 7 | Go-Zeilen der Tabelle | — | `grep -n` gegen `test/integration/integration_test.go`: `194`/`425`/`1065` = Tabellenwerte; 13 Go- + 24 Bash-Zeilen = 37 Zeilen |
| 8 | Umgebung nach den Läufen | — | keine `cdc-test-*`-Container, kein `cdc-*`-Netz zurückgeblieben; `tools/schema/plan.yaml` und `down.sql` (getrackt) unverändert → keine Lauf-Artefakte zu entfernen |

---

## Findings

### F-1 — Der E2E-Beleg bindet die Filterachse nicht an ihre Eingabe

- `kategorie`: INFO
- `quelle`: [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
  Teilfrage 4 („das Ignorieren eines Filters ist keine Auslassung, sondern
  eine stille Änderung des Ergebnisstands") · Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:1993` (Bereich
  `[pos, pos+1)`), `tools/harness/httpclient/main.go:88-94`
- `befund`: Die Phase ruft `GET /changes` in einem Bereich der Breite **1**
  auf; der Endpunkt kann in diesem Bereich nur die eine Sentinel-Change
  liefern, deren `schema`/`table` ohnehin die gefilterten Werte sind. Damit
  bleibt die Phase **grün**, wenn die Filterachse (`schema`/`table`) gar nicht
  wirkt — gemessen: Mutation C (Filterparameter entfallen, sonst unverändert)
  endet mit Exit **0**. Die Aussage der Phase „mit Filter gelesen" hat in
  dieser Aufruf-Form keinen Träger; die Filtertreue selbst ist nur
  Unit-getragen (`internal/adapters/driven/postgresstorage`,
  `internal/adapters/driving/http`).
- `verifizierbar`: ja — Mutation C dieses Laufs (Exit 0 statt rot).
  Das Ende-seitige Gegenstück (Endpunkt ignoriert den Filter) ist aus
  demselben Grund grün; es wurde **nicht** gemessen (hätte einen Image-Neubau
  verlangt), sondern nur aus der Bereichsbreite geschlossen.
- `klasse`: „Beleg ohne Bindung an seine Eingabeseite"

### F-2 — Reihenfolge- und Limit-Zusage des Clients sind in dieser Aufruf-Form unausübbar

- `kategorie`: INFO
- `quelle`: [`LH-FA-REA-004`](../../spec/lastenheft.md),
  [`LH-FA-REA-003`](../../spec/lastenheft.md) · Maintainability
- `pfad`: `tools/harness/httpclient/main.go:138-141` (Ordnungsprüfung),
  `main.go:112-118` (Doc-Kommentar der Funktion),
  `tools/harness/run-integration-tests.sh:1974` (`HTTP_READ_LIMIT=100`)
- `befund`: Die Ordnungsprüfung vergleicht `commit_position` **zwischen**
  Einträgen und die Limit-Angabe wirkt erst oberhalb der Trefferzahl; die
  Phase ruft mit Bereichsbreite 1 und `limit=100` auf und liefert gemessen
  `changes=1`. Beide Prüfungen können in dieser Aufruf-Form nicht feuern, und
  der Doc-Kommentar nennt die Positionsprüfung „die deterministische Ordnung
  des Endpunkts (`LH-FA-REA-004`)", obwohl nur der erste der drei
  Sortierschlüssel (`commit_position`, `transaction_id`, `sequence`)
  verglichen wird.
- `verifizierbar`: nein — es ist keine Gegenprobe möglich, die die beiden
  Prüfungen hier zum Tragen brächte; die volle Ordnung trägt die
  Unit-Ebene.
- `klasse`: „Assertion, die die Aufruf-Form nicht ausüben kann"

### F-3 — Anker der Abdeckungs-Deklaration bricht mitten in der Phrase ab

- `kategorie`: LOW
- `quelle`: Maintainability (Form der 24 übrigen `abdeckung_declare`-Aufrufe)
- `pfad`: `tools/harness/run-integration-tests.sh:1948`
- `befund`: Der vierte Argument-Wert („Anker") lautet
  `GET /changes real per HTTP mit reader-Token (die eigens eingefügte Zeile`
  und endet mit einer offenen Klammer, während die 23 übrigen Deklarationen
  einen abgeschlossenen Satzteil tragen. Der Anker löst auf (gemessen: Zeile
  2061) — die Abweichung ist Lesbarkeit, keine Wirkung.
- `verifizierbar`: ja — Messung 6 dieses Laufs (Anker-Nachrechnung),
  `make test-integration` bricht bei nicht auflösendem Anker ab.
- `klasse`: „Anker bricht mitten in der Phrase ab"

## Antworten auf die Schwerpunkte

**1 — Beweist der Beleg etwas — und prüft er den Endpunkt, nicht den
Prozess?** Ja. Der Client prüft inhaltlich (nicht leere Liste, `change_id`,
`operation` ∈ {INSERT,UPDATE,DELETE}, `commit_position ≥ 1`, Filtertreue,
nicht absteigende Position) und der Runner hält die `change_id` der
sentinel-tragenden READ-Zeile unabhängig gegen `cdc.changes`
(`run-integration-tests.sh:2043-2054`). Gemessen in meinem eigenen Lauf:
`changes=1`, `new_image={"id":"285","name":"HttpChangesReadE2ESentinel"}` —
der Sentinel reist **im Antwortkörper** des Endpunkts, und die
`change_id` bindet ihn an den Store. Mutation B zeigt, dass die
**Inhaltsprüfung** rot wirkt (Endpunkt antwortete `200`, die leere Liste
brach ab) — der Beleg hängt also nicht am HTTP-Status allein. Die rote
Gegenprobe des Implementers (Pfad mutiert → `404`, Exit 2) ist **ehrlich** im
Sinne von „derselbe Lauf, nur der Pfad mutiert"; sie übt aber nur die
Erreichbarkeits-Hälfte. Die Inhalts-Hälfte habe ich mit Mutation B selbst
rot gesehen. Die Grenze des Belegs ist F-1.

**2 — Die `set +e`-Änderung: richtige Stelle, richtige Form?** Ja, und sie
verdeckt nichts. Das Fenster umschließt **genau** die Zuweisung
(`run-integration-tests.sh:2000-2010`); unmittelbar danach steht die
unbedingte Prüfung `if [ "$http_status" -ne 0 ]` mit eigener Fehlerzeile. Die
Lockerung ändert damit die **Diagnose**, nicht die Rot-Grün-Aussage: unter
`set -e` wäre derselbe Ausfall ebenfalls rot gewesen, nur ohne Meldung.
Gemessen mit Mutation A (Client endet still mit 1, keine Ausgabe): Exit **2**,
Meldung `endete mit Ausgang 1: exit status 1` — ein stiller Nicht-Null-Ausgang
bleibt sichtbar rot. Das Muster ist zudem Bestand (sechs frühere
`set +e`/`set -e`-Paare im selben Skript, `262/273/284/986/1231/1288`). Die
zitierte `AGENTS.md` §3.9 ist die Regel des **Rolleneingriffs** auf
`make`-Kommandos, nicht eine Regel über `set -e` in einem Skript; das ist ein
etwas weiter gefasster Rang-Zeiger, aber nicht falsch — kein Finding.

**3 — Der Tabellen-Diff trägt den Nachweis — und die Zeilen-Drift ist
mitgezogen?** Ja, beides gemessen. (a) Die Zeile `HTTP-API-Rundlauf` trägt in
der Beschreibung das Lesen und als `Ort`
`tools/harness/run-integration-tests.sh:2061`; das ist die Echo-Zeile der
Phase, wie bei den 23 übrigen Zeilen. (b) Unabhängig nachgerechnet: mein
`awk`-Nachbau der Runner-Logik löst **24/24** Anker auf genau die
Tabellenwerte auf; die Go-Zeilen stimmen mit der Testdatei. Die fünf
verschobenen Zeilen wanderten je um **+67** (die Einfügung liegt hinter der
NATS-Negative-Deklaration, alle Zeilen davor bleiben stehen). (c) Der Beleg
ist ein **Erzeugnis**, keine Behauptung: nach dem vollen Lauf ist
`git status --porcelain` leer — die committete Datei ist byte-gleich zu dem,
was der Lauf schreibt (der Runner schreibt nur bei Abweichung).

**4 — Kein zweiter Client, kein neuer Endpunkt, `internal/**` unberührt?**
Bestätigt. `tools/harness/` trägt vier Clients (`grpcclient`, `httpclient`,
`natssub`, `sseclient`); geändert ist **einer**. Der Commit berührt fünf
Dateien, keine unter `internal/**`, `spec/**`, `tools/schema/**`,
`.a-check.yml`. Der neue Import ist `net/url` (Standardbibliothek) unter dem
bestehenden Glob `tooling: ["tools/harness/**"]`; `make a-check` grün
(Messung 1). `ADR-0068` unberührt.

**5 — Die Zeitachse: deterministisch oder bestandsabhängig?** Deterministisch.
Die Phase erzeugt ihr Subjekt **selbst** (`INSERT … id=285`, Zeile 1976-1978)
und wartet dessen Erfassung über `cdc.changes` ab, **bevor** sie liest; sie
liest danach nicht „irgendeinen Bestand", sondern den Bereich
`[pos, pos+1)`. Die End-exklusiv-Semantik des `to` ist keine Annahme des
Runners, sondern die des Ports („Start inklusiv, End exklusiv",
`internal/application/port/inbound/readchanges.go:16-19`) und der Abfrage
(`t.commit_position < $3`); `from`/`to` tragen dieselbe Größe wie die
SQL-Spalte `cdc.changes.commit_position` (`t.commit_position` in beiden
Projektionen). `id=285` kollidiert weder mit den übrigen Runner-Ids auf
`feed_e2e_full` (200/210/230/231/240/241) noch mit dem Testpaket. Belegt
durch vier vollständige Läufe dieses Reviews (einer grün ohne Mutation, drei
mit) — kein Lauf hing an zufälligem Bestand. Die Aktivierung von
`feed_e2e_full` über den ganzen Lauf trägt der Publication-Entzug-Beleg auf
einer **eigenen** Wegwerf-Tabelle (`feed_e2e_walsender_timing`), nicht auf
dieser.

**6 — Hygiene.** Sauber. Keine host-lokalen Pfade (§3.11) in den geänderten
Dateien; kein `//nolint` (§3.2); kein Trailer, Betreff
`test(e2e): … (LH-FA-SST-006, ADR-0081)` ohne `SPEC-*`/`ARC-*`; Lauf-Artefakte
in keinem Commit und nach vier Läufen nicht einmal im Baum; keine
Container-/Netz-Reste. `ADR-0044` ist **nicht** berührt: `.dockerignore`
schließt `tools/**` aus dem Build-Kontext aus (nur `cmd/`, `internal/`,
`go.mod`, `go.sum`, `tools/coverage-gate.sh`) — dieses Diff kann den
Image-Digest nicht ändern, ein Digest-Commit entfällt zu Recht.

## Negativbefunde

- geprüft, ohne Befund: `internal/**` (unberührt; der Endpunkt ist mit
  `slice-086` abgenommen)
- geprüft, ohne Befund: `tools/schema/**` · `spec/**` · `.a-check.yml` ·
  `Makefile`/`harness/mk/**` (kein Eingriff, keine neue Schichtkante)
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` (Erzeugnis echt,
  37 Zeilen, Drift vollständig)
- geprüft, ohne Befund: `harness/README.md` (der Belegsatz deckt sich mit dem
  gemessenen Phasen-Ausgang; `· seit slice-087` als Herkunfts-Anker)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (keine neue
  Betreiber-Oberfläche — kein neuer Endpunkt, keine neue `CDC_*`-Variable)
- geprüft, ohne Befund: ID-Schema/Traceability (`LH-FA-SST-006`, `ADR-0081`
  im Betreff; keine Struktur-ID; keine erfundene Kennung)
- geprüft, ohne Befund: Kommentar-Klassen (§3.7) in beiden geänderten
  Code-Dateien — alle neuen Kommentare nennen Zusage/Kopplung/Rang-Zeiger,
  keiner beschreibt abwesenden Text; keine Slice-/Wellen-Chronik in einem
  Produktionscode-Pfad (kein solcher Pfad berührt)
- geprüft, ohne Befund: Poll-Idiom `… 2>/dev/null || true` (Zeile 1983) —
  Bestandsmuster im selben Skript (Zeile 332), keine Abweichung
- geprüft, ohne Befund: `ADR-0030`/`ADR-0057` (kein neues Gate, kein neuer
  Test-Tier, derselbe Wegwerf-Client)
- geprüft, ohne Befund: Vorlagen-Reste und Behandlungs-Hinweise (die
  `BEDIENHINWEIS`-Blöcke bleiben repo-weit in 81 Plänen stehen)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Beleg ohne Bindung an seine Eingabeseite" ·
„Assertion, die die Aufruf-Form nicht ausüben kann" · „Anker bricht mitten in
der Phrase ab"

*Benachbart, ohne Zählung:* F-1 ist die positive Schwester von
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (1×, `slice-086`) — dort
band ein **Negativ**test seine Eingabe nicht, hier bindet der **positive**
E2E-Beleg seine Filter-Eingabe nicht; sichtbar wird beides erst durch
Mutation der Eingabeseite. Ob die Closure die Kennung zitiert oder eine neue
anlegt, entscheidet der Lese-Schritt, nicht dieser Report.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Die Fitness-Function-Zeile
`make test-integration` aus `ADR-0081` hat mit diesem Commit einen realen
Träger: der Endpunkt wird über den Netzwerkweg angesprochen, die Antwort
inhaltlich geprüft und ihr Change unabhängig gegen `cdc.changes` gehalten;
die rote Seite ist doppelt belegt (Gegenprobe des Implementers für die
Anfrage, meine Mutation B für den Inhalt). F-1 und F-2 sind Grenzen der
**Aussagekraft** des Belegs, keine Verhaltens- oder Vertragsfehler; F-3 ist
eine Formulierung an einem String-Argument. Keine davon verlangt eine
Fixrunde am Implementer — sie können mit der Closure reisen.

**Übergabe:** kein **Rückgabe-Pfeil** Reviewer → Implementer. Die drei
Findings gehen ohne Rückkante in die Slice-Closure §7 und von dort als
Finding-Klassen in den Zähler. Der Report selbst ist Lauf-Beleg und wird über
Läufe hinweg nicht gelesen; er ersetzt keine Verifikation (Modul 11) — die
Prüfung der Liefer-Punkte gegen die DoD bleibt beim Verifier, die
§6-Risiko-Ausgänge und die drei Paarungen bei der Planner-Closure.

**DoD-Häkchen „Review durchgeführt":** wird **in diesem Commit** auf `[x]`
nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde), mit Verweis auf diesen Report — kein HIGH/MEDIUM und kein
LOW/MEDIUM-Finding verlangt einen zweiten Implementer-Lauf. Nur diese eine
Zeile des Plans ist geändert; DoD-Substanz und §7 bleiben offen
(Verifier/Planner).
