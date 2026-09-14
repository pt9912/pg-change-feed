# Verifikationsbericht: slice-071 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-071` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5
Closure-Trigger, §6 Risiken, §8 Sub-Area) und die bindenden Entscheidungen
([`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
— `Accepted`, superseded **nur** die Änderungs-Ausnahmeklausel von
[`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md);
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) Teilfragen 3/4/5/6
und §Fitness Function;
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
§Entscheidung Festlegung 4, Fire-and-Forget ohne Replay). Nicht gegen den Diff
als solchen (Reviewer-Aufgabe; [`review-slice-071`](review-slice-071.md) in
beiden Läufen vollständig gelesen, nur als Kontext) und nicht gegen realen
Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8,
Stand `HEAD = b540449`), die vollständige
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md),
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) und
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
das Architect-Verdikt
[`architect-verdict-a-check-composition-root-tools`](architect-verdict-a-check-composition-root-tools.md),
`welle-19.md`, den vollständigen Review-Report (beide Läufe), `.a-check.yml`,
`tools/harness/grpcclient/main.go`, den gRPC-Block in
`tools/harness/run-integration-tests.sh`, `compose.yaml`, die
`harness/README.md`-Zelle und `LH-FA-SST-008` im Lastenheft. Alle Sensoren und
Stichproben wurden **eigenständig real ausgeführt** (Exit-Code je in einem
eigenen, ungepipten, mechanisch konditionierten Schritt, `AGENTS.md` §3.9) —
kein Implementer-/Reviewer-/Architect-Beleg ungeprüft übernommen. Der reale
E2E-Rundlauf und die Mutation sind **eigene** Läufe.

**Gegenstand:** `docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md`
zum Stand `HEAD = b540449`. Der Elter-`3425bfb` ist der reine
`next → in-progress`-Move (`git show 3425bfb --stat`: `{next => in-progress}`).
Der Slice-Diff `3425bfb..HEAD` nennt elf Pfade; die Slice-eigenen Code-Artefakte
sind `tools/harness/grpcclient/main.go` (neu, `b835dde`),
`tools/harness/run-integration-tests.sh` (gRPC-Block),
`compose.yaml` (`CDC_GRPC_ADDR`), `harness/README.md` (Sensor-Zelle) und
`.a-check.yml` (`tooling`-Gruppe + Kante, Architect-Zug `e051071`). Die übrigen
Pfade sind Plan-/ADR-/Review-Artefakte (Planner-/Architect-/Reviewer-Arbeit).

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — die implementierungs-/
reviewbezogenen Zeilen sind Prüfgegenstand, die Closure-Zeilen **müssen** offen
bleiben, bis der Planner sie schließt.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | [`LH-FA-SST-008`](../../spec/lastenheft.md) Happy Path/Negative erfüllt, `make test-integration` — Client empfängt eine reale committed Änderung über den laufenden Feed-Container, **und** ein Öffnungsversuch ohne gültiges Token wird real abgelehnt | **erfüllt, selbst real reproduziert** | Eigener grüner `make test-integration`-Lauf (Exit **0**, §2): RECEIVED-Zeile `change_id=964-1 table=feed_e2e_full operation=INSERT new_image={"id":"261","name":"GrpcStreamE2ESentinel"}`, danach `REJECTED code=Unauthenticated` (§5). Beides durch **zwei** eigene Mutationen rot-fähig belegt — die `change_id`-Bindung rot (§3), die Assertion-Kette am realen Stream (§5) |
| 2 | `compose.yaml` exponiert `CDC_GRPC_ADDR`; `run-integration-tests.sh` fährt den gRPC-Rundlauf als Teil des Compose-Integrationstests | **erfüllt** | `compose.yaml` trägt `CDC_GRPC_ADDR: ":9090"` in derselben Doppelpunkt-Form wie `CDC_HTTP_ADDR`; der gRPC-Block `run-integration-tests.sh:1757-1900` läuft real im selben Skript — **eigener** Lauf (§2) hat ihn ausgeführt (§5) |
| 3 | `make gates` grün | **erfüllt, selbst ausgeführt** | Eigener Lauf Exit **0** (§2) |
| 4 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | [`review-slice-071`](review-slice-071.md) vorhanden (zwei Läufe, Fixrunden-Vermerk; zweiter Lauf 0 HIGH / 0 MEDIUM); Reviewer-Zug in eigenem Kontext, Architect-Zug (`ADR-0068`) in drittem Kontext |
| 5 | Doku-Update `harness/README.md` §Sensors (gRPC-Rundlauf-Satz) | **erfüllt** | `harness/README.md`-Zelle `make test-integration` trägt den gRPC-Satz mit `· seit slice-071`; die Zeile nennt Tabelle/Operation/Spaltenwert + `change_id` gegen `cdc.changes` und verweist die Feldvollständigkeit an `server_test.go` — deckungsgleich mit der Assertion (§5) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<…>`-Platzhalter (Volltext gelesen). Planner-Arbeit nach diesem Bericht |
| 7 | Reconciliation-Register — entfällt (GF-Repo) | **korrekt offen, Entfall-Vermerk trägt** | `docs/plan/planning/reconciliation.md` existiert real nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `grep -rln "slice-071" docs/plan/planning/observations/` → nur der bereits bestehende `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`-Treffer der §8-Sichtung; **kein** neues `evidence/slice-071.md`. Die Eintragung (F-1-Klasse, §7) gehört in die Closure |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Einträge tragen wörtlich „**Ausgang:** wird bei Closure zugewiesen"; keiner ist vorzeitig geschlossen |
| 10 | Die drei Paarungen | **korrekt offen** | Slice liegt real in `in-progress/`; `Welle: welle-19`-Feld vorhanden; die DoD-Zeile verweist korrekt auf die `welle-19`-Closure |

**Ergebnis §1:** Alle fünf implementierungs-/reviewbezogenen DoD-Punkte (1–5)
sind real erfüllt. Punkt 1 wurde **selbst real reproduziert** — der reale
E2E-Rundlauf ist grün, und die `change_id`-Bindung ist über eine eigene Mutation
als **rot** belegt (§3). Die fünf Closure-Punkte (6–10) sind korrekt offen,
keine ist vorweggenommen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien`; d-check 562 Dateien/0 Befunde (links/anchors/ids/matrix/versions/structure), `commits`-Modul `HEAD~5..HEAD` 0 Befunde; `commit-traceability: OK — 5 Commit(s)`, Betreffe ohne Struktur-ID; a-check `gesamt: 0 Befund(e)`; `coverage-gate: OK — Coverage 48.50% erfüllt Schwelle 35%` (bootstrap-aware Einstiegsstufe [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)) |
| `make test` | **0** | 28 Pakete `ok`, **kein `FAIL`** (`grpc`, `grpcstream`, `capture`, `bootstrap` grün) |
| `make a-check` | **0** | `gesamt: 0 Befund(e)` — **kein** Abdeckungs-Hinweis (die drei Client-`.go` liegen jetzt in `tooling`, §4) |
| `make doc-commits RANGE=3425bfb..HEAD` | **0** | 562 Dateien, 0 Befunde — jeder Commit des Slice-Fensters ist kennungstragend |
| `make doc-immutable RANGE=3425bfb..HEAD` | **0** | 0 Befunde — `ADR-0041`s Datei ist im Fenster unberührt (`AGENTS.md` §3.5) |
| `make test-integration` (unverändert, Clean-Baseline) | **0** | alle Rundläufe grün, inkl. des gRPC-Blocks (§5) |
| `make test-integration` + Mutation (§3) | **2** | rot genau an der `change_id`-Bindung |

`git status --porcelain` nach allen Sensor-, Mutations- und Probenläufen:
**leer** (§3/§4).

## 3. Mutation — Kernzusage selbst rot/grün gesehen (Verifier-only)

Der einzige reale Beleg dieses Slice ist der Compose-E2E-Lauf; ich habe die
tragende Zusage **selbst** mutiert und **selbst** rot gesehen.

**Mutation:** in `tools/harness/grpcclient/main.go:73` die gedruckte
`change_id` durch das Literal `"FALSCH-260"` ersetzt (die RECEIVED-Zeile trägt
damit eine `change_id`, die in `cdc.changes` nicht existiert; Tabelle,
Operation und das Row Image bleiben unverändert — die Mutation trifft **nur**
die Identitäts-Bindung).

| Schritt | Ergebnis |
|---|---|
| `make test-integration` mit der Mutation | **Exit 2**. Der Lauf passiert READY, RECEIVED, die Assertion auf Tabelle/Operation/`new_image` und den Negative-Beleg und stoppt **genau** an der erweiterten Assertion: „`gRPC-Stream-Rundlauf — die über den Stream empfangene Änderung (change_id=FALSCH-260, GrpcStreamE2ESentinel) ist nicht real über cdc.changes lesbar (count=0)`" |
| Rücknahme | `git checkout -- tools/harness/grpcclient/main.go` |
| Blatt-Identität | `git hash-object` = `c4ed07e8e025ce0c6762277bfc0a524e6280cd64` = `git rev-parse HEAD:tools/harness/grpcclient/main.go` |
| Arbeitsbaum | `git status --porcelain` leer |

**Ergebnis §3:** Die Assertion bindet die **empfangene** `change_id` an
`cdc.changes` — sie prüft nicht nur die Existenz einer Sentinell-Zeile (die
Mutation ließ Tabelle, Operation und Row Image intakt und wurde trotzdem
gefangen). Die Clean-Baseline läuft grün (Exit 0), die Mutation rot (Exit 2),
die Rücknahme ist blatt-identisch. Damit ist der Slice-Beleg **real
reproduziert**, nicht übernommen.

## 4. Die neue Kante (`tooling` → `adapters`) gegen `ADR-0068`

**Deklaration (`HEAD`):** `composition_root: ["internal/bootstrap/**", "cmd/**",
"test/integration/**"]` — **kein** Werkzeug-Pfad (`ADR-0068` Festlegung 1 erfüllt).
`layers` führt `tooling: ["tools/harness/**"]` und **genau eine** Kante
`{from: tooling, to: adapters}` (`grep -c "from: tooling"` → **1**; `ADR-0068`
Festlegung 2 erfüllt).

**Zwei eigene Proben** (Exit-Code je eigener, ungepipter Schritt; danach
zurückgenommen):

| # | Probe | Erwartung | Ergebnis |
|---|---|---|---|
| A | `tools/harness/grpcclient/probe_scope.go` mit Import `internal/application/usecase/capture` | `wrong-direction`, Exit 2 | **Exit 2**: `tools/harness/grpcclient/probe_scope.go:3: wrong-direction: tooling -> app (…/usecase/capture)` — die Kante gewährt nur `adapters` |
| B | `tools/probe_scope/main.go` (außerhalb `tools/harness/**`) mit Import `…/streamv1` | wieder rot, Exit 2 | **Exit 2**: `wrong-direction: (ohne Schicht) -> adapters (…/streamv1)` **plus** Abdeckungs-Hinweis — die Gruppe reicht nicht über `tools/harness/**` hinaus |

**Rücknahme und Blatt-Identität:** beide Probe-Pfade gelöscht,
`git hash-object .a-check.yml` = `04415977511a49d93ac2783f59109c621ffc148f` =
`git rev-parse HEAD:.a-check.yml`; anschließend `make a-check` Exit **0**
(`gesamt: 0 Befund(e)`), `git status --porcelain` leer.

**Ergebnis §4:** `ADR-0068`s Festlegungen 1 und 2 sind am deklarierten Stand
und am Lauf real belegt. Die Probe A trifft `tooling -> app` — dieselbe Aussage,
die `ADR-0068` §Fitness Function und das Architect-Verdikt nennen; Probe B zeigt
die Gegenrichtung (außerhalb der Gruppe fällt der Stub-Import auf „(ohne
Schicht)" zurück). Die in `ADR-0068` §Konsequenzen **benannte** Grenze (die
Kante erlaubt die ganze `adapters`-Schicht, nicht nur `streamv1`) habe ich nicht
separat nachgestellt — sie ist als benannte Grenze geführt, kein
Verifikations-Gegenstand dieses Berichts.

## 5. Rundlauf-Ordnung, Assertion-Bindung, Negative-Beleg

- **Ordnung im Skript — vollständig vor beiden Grenzen.** Der gRPC-Block liegt
  `run-integration-tests.sh:1757-1900`. Der Upgrade-Sicherheits-Container-Tausch
  (`$COMPOSE up -d --force-recreate`) beginnt `:1902`; die
  Container-Ende-Grenze von `TestE2ESchemaChangeDropColumn` folgt `:1999-2013`,
  die von `TestE2ESchemaChangeIncompatibleTypeChange` `:2091-2108`. Der Rundlauf
  liegt damit **vollständig vor** dem Container-Tausch und vor beiden
  Schema-Negative-Aufrufen — die Bedingung, die `harness/README.md` für den
  HTTP-Beleg formuliert, ist identisch erfüllt und im Blockkommentar
  (`:1768-1774`) benannt.
- **Assertion bindet die empfangene `change_id`.** Der Lauf extrahiert die
  `change_id` aus der RECEIVED-Zeile (`grep -oE 'RECEIVED change_id=[^ ]+' |
  head -n1 | cut -d= -f2`) und hält sie über
  `SELECT count(*) FROM cdc.changes WHERE source_id='src-e2e' AND change_id='…'
  AND table_name=… AND new_data->>'name'=…` gegen den Lesezugriffsweg. Die
  Mutation (§3) beweist, dass diese Hälfte — nicht nur die Sentinell-Existenz —
  lasttragend ist. `change_id` ist in `cdc.changes` eine eigene, per
  `primary_key` eindeutige Spalte (`tools/schema/schema.yaml`), die Bindung ist
  kein Filter, der immer trifft.
- **Negative-Beleg im realen Lauf sichtbar.** Der Clean-Lauf zeigt die
  Client-Ausgabe `READY` → `RECEIVED …` → `REJECTED code=Unauthenticated`. Der
  Runner verlangt zusätzlich den Client-Ausgang 0; der Client unterscheidet
  Status-Codes und beendet jeden anderen Ausgang mit Exit 1. Der Wächter kann
  nicht leerlaufen; die Ablehnung ist real sichtbar (nicht nur die Client-Ausgabe
  wird gesucht, der Prozess-Endzustand wird ebenfalls geprüft).

**Ergebnis §5:** Auftrags-Punkte 4 und 5 sind bestätigt — Ordnung korrekt,
Assertion real bindend, Negative-Beleg sichtbar.

## 6. DoD-Häkchen-Trennung und der Restpunkt „zweiter Test-Zeiger"

- **Häkchen-Trennung — korrekt.** Der DoD-Block §2 führt fünf Zeilen auf `[x]`
  (Happy Path/Negative, `compose.yaml`+Runner, `make gates`, Review,
  README-Doku) und fünf auf `[ ]` (Closure-Notiz, Reconciliation,
  Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen). Der Reviewer hat in
  seiner Fixrunde genau die Review-Zeile nachgezogen
  (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde) — keine
  Closure-Zeile vorweggenommen.
- **Restpunkt „zweiter Test-Zeiger" — Verdikt: Planner-Ein-Zeiler, kein
  DoD-Mangel.** Die DoD-Zeile §2 nennt als einzigen Zeiger „Test referenziert:
  `make test-integration`" und wiederholt zugleich die
  Anforderungsformulierung „mit **vollständigem** Inhalt". Seit der Fixrunde
  trägt `make test-integration` Transport und Identität, die
  Feldvollständigkeit des Nachrichtenschemas (alle zehn Felder, beide Row Images)
  trägt `make test` (`internal/adapters/driving/grpc/server_test.go`). Der
  benannte E2E-Zeiger deckt die tragende Hälfte, aber nicht die
  Vollständigkeits-Hälfte.

  Das ist **kein** DoD-Mangel: die Anforderung ([`LH-FA-SST-008`](../../spec/lastenheft.md)
  Happy Path) ist real erfüllt, und ihr Beleg existiert vollständig — nur
  **über zwei Tiers**. Die DoD-Zeile ist damit geprüft erfüllt; ihr Zeiger ist
  eng, nicht falsch. Es fehlt keine Zusage, sondern ein Verweis. Ein
  „DoD-Mangel" wäre die Unerfüllbarkeit oder materielle Unbelegtheit eines
  Punktes — beides liegt nicht vor. Konsequenz: ein **Planner-Ein-Zeiler** bei
  der Closure (den zweiten Zeiger `make test` / `server_test.go` ergänzen), kein
  Merge- und kein Closure-Hindernis. Ich schließe mich dem Reviewer-Verdikt an,
  mit derselben Begründung.

## 7. `welle-19`-Stand — was nach dieser Closure noch fehlt

Der Closure-Trigger (`welle-19.md` §3) hat fünf Bedingungen; ich prüfe jede
gegen den Stand `HEAD`:

| Bedingung | Stand meiner Prüfung | Erfüllt? |
|---|---|---|
| Alle vier Slices (`slice-069`…`slice-072`) in `done/` | `slice-069`, `slice-070` in `done/`; `slice-071` in `in-progress/` (diese Closure); `slice-072` in `open/` | **nein** |
| `make gates` grün | eigener Lauf Exit **0** | **ja** |
| Realer gRPC-E2E-Rundlauf (`slice-071`, `make test-integration`) | **eigener grüner Lauf** dieses Berichts (Clean-Baseline Exit 0; §2/§5); Mutation rot | **ja** |
| Realer SSE-E2E-Rundlauf (`slice-072`, `make test-integration`) | `slice-072` in `open/`, nicht begonnen; kein Beleg | **nein** |
| Closure-Notiz in `welle-19-results.md` | Datei existiert real nicht; `welle-19.md` §7 noch Platzhalter | **nein** |

**Ergebnis §7:** Nach dieser Closure fehlen zur Welle noch **`slice-072`** (der
SSE-Endpunkt samt realem SSE-E2E-Beleg — ein zweiter, vom gRPC-Beleg
unabhängiger Nachweis über denselben `Broadcaster`) und die Wellen-Closure
(`welle-19-results.md` + `git mv` der Welle-Datei nach `done/`). Der gRPC-E2E-Beleg
liegt mit diesem Bericht **real grün** vor; der `make gates`-Beleg ist grün. Die
Welle ist mit dem gRPC-Beleg allein **nicht** schließbar — `slice-072` und der
SSE-Beleg stehen aus.

## Finding (Verifier → Planner, kein Blocker): `welle-19` §6 Out-of-Scope ist durch den ausgelieferten Stand widerlegt

**Beobachtung:** `welle-19.md` §6 (Out-of-Scope) führt als letzten Punkt:
„**`.a-check.yml`** — beide ADRs stellen fest, dass keine Änderung nötig ist
(bestehende Globs decken beide neuen Pakete ab); ein Slice, der dennoch eine
Änderung vornimmt, hat den Plan geändert, nicht nur ergänzt." Der ausgelieferte
`slice-071` nimmt genau diese Änderung vor (`.a-check.yml`: Rücknahme von
`tools/**` aus `composition_root`, neue Gruppe `tooling` + Kante, Architect-Zug
`e051071`, Träger [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)).
`welle-19.md` ist seit seinem Eröffnungs-Commit `2621d89` **unverändert** — der
Satz steht also weiter im Plan, obwohl der gelieferte Stand ihn widerlegt.

**Bewertung:** Kein DoD-Bezug des Slice (`slice-071`s §2 enthält die
Welle-Out-of-Scope nicht); die Zusage des Slice ist erfüllt. Nach der
Out-of-Scope-Disziplin (Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice, Klasse 1/2 — „Wer später mitnimmt, was hier ausgeschlossen
war, hat den Plan **geändert**, nicht nur ergänzt") ist das aber ein
**Plan-vs-Code-Diff**, den die `welle-19`-Closure tragen muss: entweder den
§6-Punkt auf den ausgelieferten Stand nachziehen (die Begründung „keine Änderung
nötig" ist durch den Client-Stub-Import falsifiziert) oder die Abweichung im
Welle-Plan benennen. Der Slice selbst hat seinen Teil getan (die §3-Zeile des
Slice-Plans ist nachgezogen), die Welle-Ebene nicht.

**Verdikt:** Planner-Aufgabe bei der `welle-19`-Closure (Modul 6, Eröffnung
Schritt 1 / Closure Schritt 6). Kein Implementer-Rückweg, kein Blocker für die
`slice-071`-Closure — aber eine benannte Drift, keine stille.

## Negativbefunde

- **geprüft, ohne Befund: `compose.yaml`-Erweiterung.** `CDC_GRPC_ADDR: ":9090"`
  steht in derselben Doppelpunkt-Form wie `CDC_HTTP_ADDR: ":8090"` und
  `CDC_NATS_URL`; der Kommentar benennt Form und Grund (Netz-Alias, Bindung auf
  allen Interfaces) und dieselben zwei Token-Klassen wie der HTTP-Adapter. Kein
  Widerspruch zu [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  Teilfrage 6 („Aktivierung optional über `CDC_GRPC_ADDR`").
- **geprüft, ohne Befund: `harness/README.md`-Zelle.** Der gRPC-Satz steht in der
  `make test-integration`-Zelle, trägt `· seit slice-071` in der Form der
  Nachbarsätze und ändert die Bindungsspalte nicht („kein Gate, `ADR-0030`,
  `LH-QA-POR-003`") — keine neue Gate-Behauptung, kein neues Target. Der Satz ist
  deckungsgleich mit der Assertion (Tabelle/Operation/Spaltenwert + `change_id`
  gegen `cdc.changes`; Feldvollständigkeit an `server_test.go` verwiesen).
- **geprüft, ohne Befund: Kommentar-Disziplin.** Kein `//nolint`/`#noqa`
  (`AGENTS.md` §3.2), kein Build-/Test-Pfad außerhalb von `make` (§3.1); der
  Client läuft per `go run` im gepinnten Toolchain-Container mit derselben
  Modul-Cache-Volume wie der bestehende `httpclient`. In den **hinzugefügten**
  Zeilen von `tools/harness/grpcclient/main.go` und dem gRPC-Block des Runners:
  kein `slice-`/`welle-`-Chronik-Treffer, kein `SPEC-`/`ARC-`-Struktur-ID-Treffer
  (`BEO-PGC/slice-chronik-in-code-kommentar`).
- **geprüft, ohne Befund: Traceability und ID-Schema.** `make doc-commits
  RANGE=3425bfb..HEAD` Exit 0; kein `SPEC-*`/`ARC-*` in den Betreffs;
  `make doc-immutable` Exit 0 (`ADR-0041` nicht in-place überschrieben,
  `AGENTS.md` §3.5; die Nachfolge-Linie steht im ADR-Index, Zeilen für
  [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
  („→ `ADR-0068`, teilweise") und
  [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)).
- **geprüft, ohne Befund: Scope-Fidelity.** Kein Pfad unter `internal/**` im
  Slice-Code-Diff; `test/integration/**` unberührt; der Client importiert genau
  ein internes Paket (`…/streamv1`) und liegt außerhalb jeder Produktionsschicht.
  `tools/harness/grpcclient/**` ist über `.dockerignore` vom Build-Kontext
  ausgeschlossen — ein `harness/image-hash.txt`-Nachzug ist für diesen Diff nicht
  fällig.
- **geprüft, ohne Befund: `ADR-0060`s Fitness-Function-Zeile (`.a-check`
  unverändert).** Sie ist auf den Gegenstand von `ADR-0060` bezogen (die zwei
  neuen Pakete liegen in `adapters`/`ports`, deren Globs unberührt sind); die
  tatsächliche `.a-check.yml`-Änderung trägt
  [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  als Träger. Keine stille Gate-Lockerung — real rot-fähig belegt (§4).

## Summary

| Prüfung | Ergebnis |
|---|---|
| DoD-Punkte 1–5 (implementierungs-/reviewbezogen) | **erfüllt**, Punkt 1 selbst real reproduziert |
| DoD-Punkte 6–10 (Closure) | korrekt offen, keine vorweggenommen |
| Sieben Sensoren (`gates`, `test`, `a-check`, `doc-commits`, `doc-immutable`, Integration clean+mutation) | je Exit 0 bzw. Mutation **2** |
| Mutation (gefälschte `change_id` in der RECEIVED-Zeile) | **rot** (Exit 2) genau an der `change_id`-Bindung; Rücknahme + Blatt-Identität verifiziert |
| a-check-Probe A (`tools/harness/**` → Use-Case) | **Exit 2** `wrong-direction: tooling -> app` |
| a-check-Probe B (außerhalb `tools/harness/**` → Stub) | **Exit 2** `wrong-direction: (ohne Schicht) -> adapters` |
| `composition_root` ohne Werkzeug-Pfad · genau eine `tooling`-Kante | **bestätigt** (3 Träger; `from: tooling` = 1) |
| Rundlauf-Ordnung (vor Ende-Grenze und vor Upgrade-Tausch) | **bestätigt** |
| Negative-Beleg `REJECTED code=Unauthenticated` im realen Lauf | **sichtbar** |
| DoD-Häkchen-Trennung | **korrekt** |
| Restpunkt „zweiter Test-Zeiger" | **Planner-Ein-Zeiler, kein DoD-Mangel** |
| `welle-19` §6 Out-of-Scope `.a-check.yml` | **widerlegt — Planner-Drift, kein Blocker** |
| `welle-19`-Closure-Reife | **noch nicht schließbar** — `slice-072` + SSE-E2E-Beleg + `welle-19-results.md` fehlen |

**Verdikt:** `slice-071` erfüllt seine DoD. Die Kernzusage
([`LH-FA-SST-008`](../../spec/lastenheft.md) Happy Path + Negative) ist **real
reproduziert**: der eigene `make test-integration`-Lauf ist grün (Exit 0), die
RECEIVED-Zeile trägt die Änderung mit Tabelle, Operation und Row Image, ihre
`change_id` ist gegen `cdc.changes` gebunden, und der tokenlose
Öffnungsversuch wird sichtbar mit `Unauthenticated` abgelehnt. Eine **eigene**
Mutation (gefälschte `change_id`) färbt genau die Identitäts-Bindung rot
(Exit 2, `count=0`), mit blatt-identischer Rücknahme. Die neue Kante
`tooling → adapters` steht deckungsgleich mit
[`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(`composition_root` ohne Werkzeug-Pfad, genau eine Kante) und ist über zwei
eigene Proben in **beide** Richtungen rot-fähig. DoD-Häkchen sind korrekt
getrennt; die Closure-Punkte sind offen. Alle sieben selbst ausgeführten
Sensoren laufen Exit 0 (Mutation Exit 2).

**Einzige Rest-Beobachtungen ohne Zwang:** (1) der DoD-Zeiger §2 darf bei der
Closure den zweiten Träger `make test`/`server_test.go` bekommen
(Planner-Ein-Zeiler; kein DoD-Mangel); (2) die `welle-19` §6 Out-of-Scope-Zeile
zu `.a-check.yml` ist durch den gelieferten Stand widerlegt und ist bei der
Wellen-Closure nachzuziehen (Planner).

**Übergabe an den Planner (Verifier → Planner):** Der Slice ist DoD-konform;
**kein** Implementer-Rückweg nötig. Für die `slice-071`-Closure bleiben:
Closure-Notiz §7 mit Steering-Loop-Eintrag (inkl. der F-1-Klasse aus dem Review
und der Zuordnung zu `BEO-PGC/a-check-null-abdeckung`), die zwei
§6-Risiko-Ausgänge (heute „wird bei Closure zugewiesen"), die
Beobachtungs-Register-Entscheidung, der optionale zweite Test-Zeiger und der
`git mv` nach `done/`. Für die danach folgende `welle-19`-Closure: dieser Lauf
ist der repo-weite Verifikations-Beleg aus Closure-Schritt 1 (Modul 6/8) für die
gRPC-Hälfte; die Welle ist erst schließbar, wenn `slice-072` in `done/` liegt,
der reale SSE-E2E-Beleg grün vorliegt, die §6-Out-of-Scope-Drift nachgezogen ist
und `welle-19-results.md` geschrieben ist.
