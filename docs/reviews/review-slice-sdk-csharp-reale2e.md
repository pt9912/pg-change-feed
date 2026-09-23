# Review-Report: slice-sdk-csharp-reale2e — 2026-09-23

**Review-Art:** Code — geprüft gegen Plan
(`../plan/planning/done/slice-sdk-csharp-reale2e.md`, §2 DoD, §3 Plan +
Plan-Nachzug + §3.13-Suchlauf-Feld, §6 Risiken), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Accepted) §Entscheidung Festlegung 2/Folgepflicht 1,
[`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md) (Ort
`sdks/csharp/`, Import-Grenze), [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
(benannter Bau-Kontext `proto`), `AGENTS.md` §3 Hard Rules (§3.2, §3.7,
§3.9, §3.11, §3.12, §3.13). Kein DoD-Abgleich als solcher — Verifier-Aufgabe
(Modul 11); wo eine DoD-Zeile eine **Zahl** behauptet, ist das Nachmessen
dieser Zahl trotzdem Reviewer-Aufgabe (`AGENTS.md` §3.12 Instanz A,
Reviewer-Skill HIGH „Zahl im Träger … driftend").

**Gegenstand:** Diff-Range `fce7af10..HEAD` — Commits `ecff7370` +
`11e42f7b` des Slices `slice-sdk-csharp-reale2e` (13 Dateien, 822
Insertions / 10 Deletions): Runner-Skript, Docker-Stufe `integration`,
Integrations-Testprojekt (4 Phasenklassen + PhaseEnvironment + csproj),
Make-Target, Abdeckungs-Träger `docs/user/sdk-e2e-abdeckung.md`,
`.d-check.yml`-Eintrag (Label `SDK-E2E`), `harness/README.md`-Zeile,
Plan-Update (DoD-Häkchen, Plan-Nachzug, §3.13-Feld).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09".
**Modell:** glm-5.3-flash (Claude-Agent-SDK-Subagent) · **Datum:** 2026-09-23.

**Eingangs-Kontext:**

- Plan (§1–§8), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 2 (Mechanik-Klasse: reale Server-Instanz je neuer Fläche,
  Folgepflicht 1) und die dort beibehaltenen Festlegungen 3–5,
  [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md)
  Festlegung 2/4, [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 2, [`welle-sdk-reale2e.md`](../plan/planning/welle-sdk-reale2e.md)
  (§1 Träger-Erzeugnis, §4 Reihenfolge, §6 Beobachtungen-Sicht)
- `spec/lastenheft.md` [`LH-FA-SST-009`](../../spec/lastenheft.md),
  [`LH-FA-SST-008`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md),
  [`LH-FA-CON-001`](../../spec/lastenheft.md); `spec/pflichtenheft.md`
  `SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-024`/`SPEC-027` (gelesen)
- Vorbilder (gelesen, nicht importiert):
  `tools/harness/run-sdk-python-integration-tests.sh` (unverändert als
  Quelle), `sdks/python/Dockerfile`, `docs/user/e2e-abdeckung.md`
  (Träger-Form), `docs/reviews/review-slice-sdk-python-nats-stream-client-flaeche.md`
  (Lernklassen-Kette F-2/F-4/F-5/F-6/F-10)

**Eigenständig durchgeführte Prüfungen (nicht aus dem Plan-Text übernommen):**

- **Realer Lauf des Targets** `make test-sdk-csharp-integration`
  (2026-09-23, Ausgang **0**, ungepiped ermittelt — §3.9): alle vier
  Phasen grün gegen die reale Compose-Umgebung. Die DoD-Zahlen sind dabei
  **byte-identisch reproduziert** (deterministische Bring-up-Kette):
  gRPC `change_id=804-1`, SSE `807-1`, NATS `810-1` — je Phase mit
  SQL-Gegenprüfung grün; `consumer_id` ist laufgebunden (Zeitstempel-Form:
  eigener Lauf `csharp-sdk-e2e-20260923080736`, DoD nennt
  `csharp-sdk-e2e-20260923075217` — erwartbar, Form identisch). Der Runner
  meldete „Abdeckungs-Traeger unveraendert" — der committete Träger ist
  byte-identisch dem Runner-Erzeugnis, der Arbeitsbaum blieb danach sauber
  (Idempotenz real bestätigt). Der NATS-Reject-Lauf endete grün — die
  Ursachen-Bindung („Authorization" über die Ausnahme-Kette) trägt am
  realen Server (F-6-Lernklasse des Vorgänger-Slices gebunden).
- `make docs-check` → 0 Befunde (927 Dateien); `make commit-traceability`
  → OK (5 Commits, Betreffs ohne Struktur-ID); `make doc-trace`
  (advisory) → exit 0, RTM trägt die Spalte `SDK-E2E` für
  `LH-FA-SST-008`/`LH-FA-SST-006`/`LH-FA-CON-001`;
  `LH-FA-SST-009` → **WAISE** (F-2).
- **§3.13-Suchlauf über beide Stände** (`git grep` gegen `fce7af10` und
  `HEAD`): die vier Feld-Zeilen des Plans einzeln bestätigt (README-Zeile
  fehlte am Parent, existiert am HEAD; `sdks/csharp/README.md` 0 Treffer
  „Realserver" beide Stände; Benutzerhandbuch 0 Treffer;
  `spec/pflichtenheft.md` im Diff unverändert) — und drei weitere
  Fundstellen, die das Feld nicht trägt (F-1).
- **Suppressions:** `grep` über `sdks/csharp/PgChangeFeed.Client.Integration/`
  nach `#pragma`/`nolint`/`noqa`/`[SuppressMessage]`/`TODO` → 0 Treffer
  (`#noqa`-Klasse des Python-Vorgänger-Slices nicht wiederholt).
- **Import-Grenze:** einziges `ProjectReference` ist
  `../PgChangeFeed.Client/PgChangeFeed.Client.csproj`; `using`-Liste
  ausschließlich Client-Namespaces + xunit — kein Bezug auf
  `internal/**`/`cmd/**`/`gen/**`.
- **Eingabeseiten-Bindung je Phase:** Token-Form (gRPC/SSE reader-Token,
  HTTP admin/reader-Paar, NATS connection-token — je Env-gebunden über
  `PhaseEnvironment.Required`), Header-Form (in der Client-Assembly:
  `authorization`/`Bearer <token>` gRPC, `Authorization: Bearer` HTTP/SSE —
  geprüft), Subjekt-Form (`BuildSourceSubject(sourceId)` über Env), und
  die HTTP-Eingaben `PGCHANGEFEED_SOURCE_ID`/`PGCHANGEFEED_HTTP_PUBLICATION`
  (Plan-Nachzug F-6-Disziplin des Vorgänger-Slices) — gebunden.
  Mutationsprobe gedacht, nicht ausgeführt: ein korrekter Token färbt den
  Reject-Test über `Contains("Authorization")` bzw. `RpcException`-Status
  real rot (Ausnahme-Klasse gelesen, Negativseite über `test_exit != 0`
  zusätzlich gebunden).
- **Kommentar-Sichtung je übernommenem Form-Teil** (§6-Risiko 1 des Plans,
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`, Zusage):
  Runner-Kopf, Phasen-Prosa, Dockerfile-Kommentar und Make-Target-Kommentar
  je separat gegen das Python-Vorbild gesichtet — kein unübersetztes
  deutsches Wortfragment in englischen Kommentaren (BEO-Zähler bleibt 2×,
  kein drittes Auftreten); gefunden stattdessen der eigene Grammatik-Slip
  F-3 und die ASCII-Mischform F-4. Die deutschen Assert-Meldungen der
  Integrationstests sind Vorbild-Form (`sdks/python/pgchangefeed/integration/`
  trägt dieselben Meldungen wortgleich) — kein Befund.

---

## Findings

### F-1 — §3.13-Suchlauf-Feld unvollständig: drei weitere Fundstellen, von der Suchform nicht getroffen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug)
- `pfad`: `docs/plan/planning/done/slice-sdk-csharp-reale2e.md`
  (§3.13-Feld), `harness/README.md:134`, `harness/sensors/docs-check.md:143`,
  `docs/plan/planning/welle-sdk-reale2e.md:311`
- `befund`: Das committete Suchlauf-Feld trägt vier Träger-Zeilen — alle
  vier einzeln bestätigt. Ein eigener, breiterer Lauf (`grep` nach
  `trace.coverage`/„abdeckung", beide Stände) findet drei Fundstellen, die
  das Feld nicht nennt: (a) `harness/README.md:134` — die
  `make doc-trace`-Zeile zählt die `trace.coverage`-Coverage-Dimensionen
  wörtlich auf („liest zusätzlich `e2e-abdeckung.md`, `bench-abdeckung.md`
  **und** `ci-matrix-abdeckung.md`"); genau diese Aufzählung wurde durch
  den `.d-check.yml`-Eintrag dieses Diffs um einen vierten Eintrag
  überholt — der Träger war vor dem Diff exakt richtig und ist durch die
  eigene Arbeit falsch geworden, nicht gezogen, nicht gemeldet. (b)
  `harness/sensors/docs-check.md:143` nennt nur `docs/user/e2e-abdeckung.md`
  als kuratierte Coverage-Dimension (vorgelagert bereits für Bench/
  CI-Matrix überholt — der neue Eintrag verbreitert die Lücke; dieselbe
  Grep-Form findet sie). (c) `welle-sdk-reale2e.md:311` behauptet über den
  hier erzeugten Träger, er „trägt `Datei:Zeile`-Orte (Muster
  `e2e-abdeckung.md`)" — der reale Träger trägt Testklassen-/Runner-Orte
  **ohne** Zeilenanker; die Abweichung von der im Plan §3 deklarierten
  Form („Form nach `docs/user/e2e-abdeckung.md`") ist weder im
  Plan-Nachzug noch im Suchlauf-Feld dokumentiert.
- `verifizierbar`: ja — die Fundstellen sind einzeln per `grep` aufzulösen
  (alle in diesem Lauf gegen beide Stände ausgeführt); `make doc-trace`
  belegt die vierte Coverage-Dimension real.
- `klasse`: „Arbeit überholt stehenden Träger"
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`; Suchlauf-Feld zu schmal —
  Kalibrierung: Python-NATS-Review F-4, MEDIUM)

### F-2 — Coverage-Zeilen nennen `LH-FA-SST-009` nicht: die Kopf-Anforderung des DoD bleibt RTM-Waise

- `kategorie`: LOW
- `quelle`: Maintainability · Plan §2 DoD („`LH-FA-SST-009`-Beleg-Stand
  stärker") · [`welle-sdk-reale2e.md`](../plan/planning/welle-sdk-reale2e.md)
  §1 (Zweck des Trägers: RTM-Sichtbarkeit der zwölf Belege)
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh`
  (`abdeckung_csharp_abschnitt`), `docs/user/sdk-e2e-abdeckung.md:10–13`
- `befund`: Die vier neuen Coverage-Zeilen citen `LH-FA-SST-008`,
  `LH-FA-SST-006` und `LH-FA-CON-001`, nicht aber `LH-FA-SST-009` —
  jene Anforderung, deren Beleg-Stand die DoD-Kopfzeile dieses Slices
  stärkt und deren Happy Path („einbinden → verbinden und Changes
  empfangen, ohne das Protokoll selbst zu implementieren") die vier
  Phasen wörtlich demonstrieren. Nachgemessen (`make doc-trace`):
  `LH-FA-SST-009` steht in der RTM als WAISE mit leerer Abdeckungs-Spalte
  — auch am Parent (in keiner Coverage-Datei citiert, also vorgelagert
  bestehend, nicht durch diesen Diff verursacht). Der Welle-Plan schiebt
  die Matrix-Prüfung bewusst an die Welle-Closure („kein Einzel-DoD trägt
  diese Prüfung") — solange aber keine der zwölf geplanten Zeilen die
  Kennung trägt, bleibt die Kopf-Anforderung der Welle trotz zwölf Belegen
  RTM-blind; die Welle-Closure-Prüfung wird genau diese Lücke sehen.
- `verifizierbar`: ja — `make doc-trace` (advisory, exit 0, WAISE-Zeile
  `LH-FA-SST-009`).
- `klasse`: „Abdeckungszuordnung nennt ihre Anforderung nicht" (neue
  Kurzform — Vorlauf auf die Welle-Closure-Prüfung)

### F-3 — „die vier Flaeche": Grammatik-Slip im Runner-Kopf

- `kategorie`: LOW
- `quelle`: Maintainability (§6-Risiko 1 des Plans: übernommene Form-Teile
  je separat gesichtet)
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh:5`
- `befund`: „die vier Flaeche des Packages" — Singular statt Plural
  („Flächen"), in ASCII-Transliteration; die Kopfzeile davor („Zustellweg-
  Flaechen") ist korrekt, der Slip stammt aus der Neufassung, nicht
  wortgleich aus dem Python-Vorbild. Kein semantischer Effekt; kein
  unübersetztes deutsches Wortfragment im BEO-Sinne (BEO-Zähler bleibt 2×).
- `verifizierbar`: nein — Text lesen; kein Sensor für Prosa-Grammatik.
- `klasse`: „Form-Fragment im Vorbild-Kopf" (einmalig)

### F-4 — ASCII-Transliteration im erzeugten Träger; im Runner selbst beide Orthografien gemischt

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/user/sdk-e2e-abdeckung.md:3–9`,
  `tools/harness/run-sdk-csharp-integration-tests.sh` (Kopf-Kommentare
  ASCII, Fehlermeldungen mit echten Umlauten, z. B. „trägt die
  Verdrahtung nicht")
- `befund`: Der erzeugte Doku-Träger schreibt durchgängig „ueber/traegt/
  gruen/zustaendige", das Schwester-Träger-Vorbild `docs/user/e2e-abdeckung.md`
  schreibt „über/trägt/grün/zuständige"; im Runner selbst koexistieren
  beide Formen. Orthografie-Inkonsistenz ohne semantische Auswirkung —
  `docs-check` prüft sie nicht.
- `verifizierbar`: nein — Form-Vergleich, kein Sensor.
- `klasse`: „Orthografie-Split in der Träger-Familie"

### F-5 — Träger-Intro nennt als Erzeuger beide Targets; heute schreibt nur das C#-Target

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 (eine Aussage, die als Beleg gelesen wird)
- `pfad`: `docs/user/sdk-e2e-abdeckung.md:3`
- `befund`: „Erzeugt von `make test-sdk-*-integration` über die
  Runner-Skripte" nennt die Glob-Form über beide bestehenden Targets —
  `make test-sdk-python-integration` existiert real und läuft grün, aber
  sein Runner schreibt (noch) keinen Abschnitt dieser Datei; der
  Python-Abschnitt kommt erst mit `slice-sdk-python-http-reale2e`. Der
  Satz trägt seinen heutigen Beleg nur zur Hälfte; ein Zeitpunkt fehlt
  (§3.12: „Der Schreiber setzt den Zeitpunkt").
- `verifizierbar`: ja — `grep -c "abdeckung" tools/harness/run-sdk-python-integration-tests.sh`
  → 0 Treffer; der C#-Runner trägt den Schreiber.
- `klasse`: „Beleg trägt seinen Satz nicht" (schmale Form: Erzeuger-Glob)

### F-6 — `CSPROC_TEST_NAME`: kryptischer Variablenname, bricht die Namens-Reihe

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh:62`
- `befund`: Die drei übrigen Phasen-Variablen heißen `SSE_TEST_NAME`,
  `NATS_TEST_NAME`, `HTTP_TEST_NAME`; die gRPC-Phase heißt
  `CSPROC_TEST_NAME` — eine unlesbare Abkürzung ohne Entsprechung im
  Vorbild (`GRPC_TEST_FILE`) oder im Plan.
- `verifizierbar`: nein — Name lesen.
- `klasse`: „Namens-Reihenbruch in der Phasen-Tabelle" (einmalig)

### F-7 — Writer-Form regeneriert alles oberhalb des eigenen begin-Markers — Reihenfolge-Constraint für die Folge-Runner

- `kategorie`: INFO
- `quelle`: Maintainability (Hinweis an die Folge-Slices
  `slice-sdk-kotlin-reale2e`/`slice-sdk-python-http-reale2e`)
- `pfad`: `tools/harness/run-sdk-csharp-integration-tests.sh`
  (`abdeckung_schreiben`)
- `befund`: Der Writer erzeugt Kopf + Tabellenkopf neu und erhält nur den
  Inhalt **hinter** dem eigenen end-Marker — für den ersten Abschnitt der
  Datei korrekt und idempotent (real bestätigt). Kopiert ein Folge-Runner
  dieselbe Form wörtlich, verwirft sein Re-Run die Abschnitte oberhalb
  seines eigenen begin-Markers; die Plan-Nachzug-Formel „jeder Runner
  ersetzt nur seinen eigenen Abschnitt" trägt die Symmetrie über die drei
  Runner nicht von selbst. Die beiden Folge-Pläne deklarieren
  „idempotent, Marker-gegrenzt" ohne diese Form-Randbedingung.
- `verifizierbar`: ja — `awk`-Rest-Extraktion lesen; Gegenprobe erst mit
  dem ersten fremden Abschnitt real möglich.
- `klasse`: „Form-Vorgabe für Folge-Runner" (Zustands-Hinweis)

### F-8 — DoD-Lauf-Bezug „erneuten Lauf bestätigt" ohne eigenen Anker — durch Wiederholungsmessung entkräftet

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahlenwert in einem Träger)
- `pfad`: `docs/plan/planning/done/slice-sdk-csharp-reale2e.md`
  (§2 DoD, Sensor-Beleg)
- `befund`: Die DoD-Zeile nennt das Target, aber keinen Lauf-Anker und
  führt zwei Läufe („in einem erneuten Lauf bestätigt"), ohne zu sagen,
  welcher Lauf die quotierten `change_id`s trug. Die Nachmessung entkräftet
  den Bedenken-Grund real: die Bring-up-Kette ist deterministisch, ein
  erneuter Lauf (dieser Review, 2026-09-23) reproduzierte gRPC `804-1`,
  SSE `807-1`, NATS `810-1` byte-identisch — die Zahlen sind nicht
  lauf-gebunden im driftenden Sinne, sondern an der aktuellen Kette
  reproduzierbar; jede Änderung an der Bring-up-Kette (Schema-Rollout,
  Bootstrap-Aktivität) re-basiert sie. Kein Befund, Hinweis für den
  Verifier: die Messung ist wiederholbar, bis die Bring-up-Kette sich
  ändert.
- `verifizierbar`: ja — `make test-sdk-csharp-integration` erneut fahren;
  `change_id`-Gleichheit ist bis zur nächsten Bring-up-Änderung der
  erwartbare Ausgang.
- `klasse`: „Zahl im Träger — Ursprung getragen, Lauf nicht adressiert"
  (durch Messung entkräftet)

---

## Negativbefunde (geprüft, ohne Befund)

- `sdks/csharp/PgChangeFeed.Client.Integration/**` (5 Testdateien + csproj):
  keine Inline-Suppression (`#pragma warning`/`// nolint`/`#noqa`/
  `[SuppressMessage]` — 0 Treffer); keine `internal/**`/`cmd/**`/`gen/**`-
  Importe (einziger Projektverweis `PgChangeFeed.Client`); je Phase an der
  Eingabeseite gebunden (Token-Form, Header-Form in der Client-Assembly
  `authorization`/`Bearer`, Subjekt-Form, HTTP-Quelle/Publikation als Env);
  die NATS-Reject-Bindung trägt die Ursache real — der eigene Lauf
  bestätigte die Assertion gegen den echten Server (F-6-Lernklasse des
  Python-Vorgänger-Slices nicht wiederholt); deutsche Assert-Meldungen sind
  Vorbild-Form (identisch in `sdks/python/pgchangefeed/integration/`).
- `tools/harness/run-sdk-csharp-integration-tests.sh`: `set -euo pipefail`,
  Cleanup-Trap (EXIT), explizite Phasen-Auswahl je Testklasse (kein stiller
  Ausschluss), SQL-Gegenprüfung je Phase (`cdc.changes`/`cdc.consumer`),
  Fire-and-Forget-Fenster mit fünf Versuchen, Make-Aufruf ungepiped (§3.9)
  — realer Lauf exit 0, vier Phasen, Träger „unveraendert".
- `sdks/csharp/Dockerfile` (integration-Stufe): `FROM build` — kein neuer
  Pin (Wiederverwendung `mcr.microsoft.com/dotnet/sdk:10.0@sha256:60a2…`),
  `:?`-Guard gegen blanken Aufruf, `--no-build`-Kette konsistent; Bau im
  eigenen Lauf cache-grün.
- `harness/mk/sdk.mk`: Target + `.PHONY` + `##`-Hilfe im Kommentar-Stil der
  bestehenden Werkzeug-Targets.
- `.d-check.yml`: Label `SDK-E2E` aktiv — `make doc-trace` trägt die
  Spalte für `LH-FA-SST-008`/`LH-FA-SST-006`/`LH-FA-CON-001` (eigener Lauf).
- `docs/user/sdk-e2e-abdeckung.md`: Marker-gegrenzt, idempotent,
  byte-identisch dem Runner-Erzeugnis (eigener Lauf: „unveraendert");
  Linktiefe `../../spec/lastenheft.md` korrekt aufgelöst (docs-check 0
  Befunde); alle citen Kennungen matchen `LH-(FA|QA)-[A-Z]{3}-\d{3}`.
- `harness/README.md` neue Werkzeuge-Zeile: jede Verhaltens-Aussage der
  Zeile im eigenen Lauf bestätigt (Bring-up, vier Phasen, Reject-Formen
  gRPC `Unauthenticated`/HTTP `401`/NATS-Abweisung, SQL-Gegenprüfung,
  Träger-Schreibung, `:?`-Guard) — `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
  ohne Befund; kein Handbuch-Zug nötig (keine neue Betreiber-Oberfläche —
  `PGCHANGEFEED_*` sind Test-Container-ENV, keine `CDC_*`-Fläche).
- Plan-vs-Code: alle 7 §3-Zeilen und alle 6 Plan-Nachzug-Zeilen im Diff
  vertreten (Testklassen-Form, `PhaseEnvironment.cs`, HTTP-Env-Bindung,
  ID-Bereiche 400/410/420/430 je Phase, Träger-Erhalt hinter dem
  end-Marker, DTO je Fläche); kein still gestrichener Planpunkt; die drei
  §6-Risiken tragen je einen Ausgang (Marker fristnah, beide Token-Klassen
  gefahren, ID-Bereiche kollisionsfrei — real im Lauf).
- §3.13-Feld-Zeilen 2–4: `sdks/csharp/README.md` (0 „Realserver"-Treffer
  beide Stände), `docs/user/benutzerhandbuch.md` (keine E2E-Beleg-Aussage),
  `spec/pflichtenheft.md` (`SPEC-027`-Aussage im Diff unverändert) —
  bestätigt.
- Kommentar-Klassen (§3.7): Runner-/Dockerfile-/Make-Kommentare tragen
  Zustands-/Kopplungs-/Grenze-Form; Slice-Namen als etablierte
  Herkunfts-Anker-Form (Vorbild-Konform, keine Vorher/Nachher-Sprache);
  keine Slice-/Wellen-Chronik im Satzsubjekt Produktionscode-Pfad.
- Hard Rules §3.11 (Hostpfade): 0 Treffer im Diff (alle Pfade relativ).
- Gates, die die Diff-Fläche lesen: `make docs-check` exit 0,
  `make commit-traceability` exit 0 (eigene Läufe, ungepiped).

---

## Gesamt-Verdikt

**0 HIGH · 1 MEDIUM · 4 LOW · 2 INFO.** Die Mechanik ist [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)-konform
(Stufe auf `build`, kein neuer Pin, benannter `proto`-Kontext), der Runner
ist laufbestätigt grün (eigenes Nachmessen, exit 0, deterministisch
reproduziert die DoD-Zahlen), der Träger ist byte-stabil idempotent, die
Suppressions-/Import-Grenzen halten, und die F-6-Lernklasse des
Vorgänger-Slices (Reject-Ursache binden) ist an der Eingabeseite gebunden
und real getragen. Die Fix-Fläche ist klein: F-1 (drei Träger nachziehen,
Suchlauf-Feld ergänzen) ist Implementer-Substanz, F-2 bis F-6 sind
Einzelzeilen-Form.

**Fixrunde am Implementer: ja** (F-1 zwingend; F-2/F-3/F-4/F-5/F-6
dieselbe Runde, Einzelzeilen). Die DoD-Checkbox „Review durchgeführt" bleibt
deshalb offen — der Nachzug erfolgt regulär bei Schritt 21 des
Implementer-Workflows (`BEO-PGC/dod-checkbox-nachzug`); der
Selbst-Nachzug-Regelfall (Review ohne Fixrunde) greift hier nicht. Dieser
Report ist mit der Fixrunde committet (`86892bdb`).