# Review-Report: slice-037 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-037`, §1/§2/§3/§4/§6/§8) und
`ADR-0050` (Accepted, zentraler Prüfmaßstab dieses Laufs) sowie `AGENTS.md`
§3 Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `2a771bc` (Kernimplementierung: neues Domänenmodell
`AdministrationRequest`, neuer Outbound Port `AdministrationRequestPort`,
Postgres-Adapter `AdministrationRequestAdapter` + `AdministrationListener`
mit eigener `LISTEN`-Verbindung und Selbst-Reconnect, `Assembler` bekommt
`sync.RWMutex` + `AddBinding`/`RemoveBinding`/`lookupBinding`,
`receive.Stream.Assembler()`-Accessor, `wiring.go`s
`activatedTableBindings`/`runAdministration`/`processAdministrationRequests`/
`applyAdministrationRequest`, `spec/architecture.md`s Sequenzdiagramm zu
`LH-FA-CFG-001.a` in CLI/SQL getrennt, neuer End-zu-End-Beleg in
`tools/harness/run-integration-tests.sh`, `Makefile`/`harness/README.md`
für `make test` mit Race-Detector, `make image` neu gebaut), `b763253`
(DoD-Häkchen und Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-037-administrations-goroutine-live-reload.md` (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8)
- `docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md` (vollständig — zentraler Prüfmaßstab)
- `spec/lastenheft.md` (`LH-FA-ADM-001`, `LH-FA-CFG-001`, `LH-FA-CFG-002`)
- Alle neuen/geänderten Go-Dateien vollständig gelesen, nicht nur die Implementer-Zusammenfassung: `internal/adapters/driving/replication/mapper/mapper.go` + `mapper_test.go`, `internal/adapters/driven/postgresstorage/administrationrequest.go` + `administrationrequest_test.go`, `internal/adapters/driven/postgresstorage/queries/queries.go`, `internal/application/port/outbound/administrationrequest.go`, `internal/domain/model/administrationrequest.go`, `internal/adapters/driving/replication/receive/receive.go`, `internal/bootstrap/wiring.go`
- Zum Vergleich herangezogen (Bestand, nicht Teil des Diffs): `internal/application/usecase/enable/service.go`, `internal/application/usecase/disable/service.go`, `internal/adapters/driven/postgresstorage/tableactivation.go` (Idempotenz-Prüfung), alle zehn übrigen `NewX`-Konstruktoren in `internal/domain/model/`, `internal/bootstrap/heartbeat_internal_test.go`, `internal/bootstrap/walretention_internal_test.go` (Whitebox-Testmuster-Vergleich)
- `spec/architecture.md`-Diff, `Makefile`-Diff, `tools/harness/run-integration-tests.sh`-Diff, `harness/README.md`-Diff (jeweils vollständig)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.2, §3.4, §3.7
- `docs/reviews/review-slice-036.md` (Format-Vorlage)

---

## Findings

- **[HIGH]** — `quelle`: `AGENTS.md` §3.2 (Suppression-Verbot) / Reviewer-Skill
  HIGH-Klasse „Suppression eines Gates (…) ohne ADR"
  `pfad`: `internal/adapters/driving/replication/mapper/mapper_test.go`
  (Zeile `xid := uint32(i + 1) //nolint:gosec // Testschleife, kein
  Sicherheitskontext`, in `TestAssemblerLiveReloadIsRaceFree`)
  `befund`: Dies ist die **einzige** Inline-Suppression im gesamten Repo
  (`grep -rn nolint` trifft nur diese eine Zeile) — ohne begleitenden
  ADR-Eintrag und ohne Eintrag in einer zentralen Ausnahme-Konfiguration.
  Mildernder, aber nicht entlastender Kontext: Das Repo führt aktuell
  **kein** aktives Lint-/`gosec`-Gate (kein `lint`-Target im `Makefile`,
  keine `.golangci.yml`) — die Suppression zielt auf ein Werkzeug, das in
  `make gates`/`make test` nicht läuft, ist also im Moment folgenlos.
  Trotzdem ist sie der erste Beleg dieses Musters im Repo und schafft
  einen Präzedenzfall für künftige stille Suppressions, sobald ein
  Lint-Gate eingeführt wird — genau die Drift, die die Norm verhindern
  soll, unabhängig davon, ob sie heute schon einen Sensor bricht.
  `verifizierbar`: ja (`grep -rn "nolint" --include="*.go"` — ein Treffer)
  `klasse`: „Inline-Suppression ohne ADR/zentrale Ausnahme" (erstes
  Auftreten dieser Klasse in diesem Skill-Lauf)

- **[MEDIUM]** — `quelle`: Maintainability (Reviewer-Skill §Klassifikation,
  „fehlende Negativtests bei neuem öffentlichem Vertrag")
  `pfad`: `internal/bootstrap/wiring.go` (`runAdministration`,
  `processAdministrationRequests`, `applyAdministrationRequest`,
  Typ `administrationListener`)
  `befund`: Der neu eingeführte Typ `administrationListener` trägt einen
  Kommentar, der exakt das im Repo etablierte Fake-basierte
  Whitebox-Testmuster verspricht („dasselbe Whitebox-Test-Muster wie
  `walRetentionMeasurer`: eine Fälschung belegt die Fallback-Poll-Schleife
  ohne reale PostgreSQL-Instanz") — analog zu den vorhandenen
  `heartbeat_internal_test.go`/`walretention_internal_test.go`. Eine
  solche Datei (z. B. `administration_internal_test.go`) existiert im
  Diff jedoch nicht; `grep -rln "administrationListener\|runAdministration\|
  processAdministrationRequests\|applyAdministrationRequest"
  internal/bootstrap/` trifft ausschließlich `wiring.go` selbst. Die
  Fehlerpfade — `MarkFailed`-Zweig, Listener-Fehler/Reconnect,
  ctx-Cancel-Austritt aus `runAdministration`, der `default`-Zweig bei
  ungültigem `Kind` — sind damit durch keinen automatisierten Test
  abgedeckt; nur die Happy-Paths laufen indirekt über das
  E2E-Integrationsskript.
  `verifizierbar`: ja (Abwesenheit der Datei; ein Coverage-Report von
  `internal/bootstrap` würde die ungetesteten Zweige zeigen)
  `klasse`: „fehlende Negativtests bei neuem öffentlichem Vertrag"
  (bereits benannte Skill-Klasse — 2. Auftreten nach `review-slice-036`,
  Schwelle 3× noch nicht erreicht)

- **[MEDIUM]** — `quelle`: Maintainability (Reviewer-Skill §Klassifikation,
  „unklare Fehlerbehandlung am Rand des Spec-Bereichs")
  `pfad`: `internal/adapters/driven/postgresstorage/administrationrequest.go`
  (`AdministrationListener.WaitForNotification`,
  `connectAdministrationListener`)
  `befund`: Bei dauerhaft unerreichbarem `cfg.AdminDSN` läuft der interne
  Reconnect-Versuch (`connectAdministrationListener(context.Background(),
  l.dsn)`) ohne Timeout und ohne Backoff. Schlägt der Reconnect sofort
  fehl (z. B. „connection refused"), kehrt `WaitForNotification` mit
  Fehler zurück, und `runAdministration`s Schleife ruft umgehend erneut
  `processAdministrationRequests` und danach erneut
  `WaitForNotification` auf — ohne Wartezeit zwischen den Versuchen. Kein
  Datenverlust und kein Zeitfenster mit unbemerkt bleibenden Anträgen
  (im Gegenteil: das Fallback-Polling läuft in diesem Fall sogar
  häufiger, nicht seltener, siehe Negativbefunde), aber eine potenzielle
  enge Wiederholschleife ohne Drosselung bei einem Fehlerfall, den weder
  `ADR-0050` noch der Slice-Plan explizit adressieren.
  `verifizierbar`: ja, nicht selbst reproduziert (Zeitaufwand in dieser
  Sitzung nicht eingeplant) — Repro-Vorschlag: `cfg.AdminDSN` auf einen
  stets „connection refused" liefernden Port zeigen lassen und
  CPU-/Log-Zeilen-Rate der Administrations-Goroutine beobachten.
  `klasse`: „unklare Fehlerbehandlung am Rand des Spec-Bereichs"

- **[MEDIUM]** — `quelle`: Maintainability (Konsistenz mit etabliertem
  Domain-Core-Muster, `ARC-001`)
  `pfad`: `internal/domain/model/administrationrequest.go`
  `befund`: Jeder der zehn bestehenden Typen in `internal/domain/model`
  (`Consumer`, `ConsumerPosition`, `ErrorClass`, `SourcePosition`,
  `Change`, `ChangeTransaction`, `RetentionPolicy`, `Source`,
  `SchemaVersion`, `SourceTable`, `TableSchema`) hat eine validierende
  `NewX(...) (X, error)`-Konstruktorfunktion, die Invarianten am
  Domain-Core-Rand durchsetzt. `AdministrationRequest`/
  `AdministrationRequestKind` hat keine solche Funktion — der Typ wird in
  `AdministrationRequestAdapter.ListPending`
  (`internal/adapters/driven/postgresstorage/administrationrequest.go`)
  als reines Struct-Literal befüllt. Die Prüfung der geschlossenen Menge
  `enable`/`disable` liegt stattdessen im `default`-Zweig von
  `applyAdministrationRequest` (`internal/bootstrap/wiring.go`,
  Bootstrap-Schicht `ARC-007`) — außerhalb der Domain-Schicht, in der laut
  Architektur-Sicht („Domänenobjekte und Invarianten … pur, ohne
  Treiber") genau diese Prüfung erwartet würde.
  `verifizierbar`: ja (Diff-Vergleich mit den zehn übrigen
  `NewX`-Konstruktoren in `internal/domain/model/*.go`)
  `klasse`: „Domäneninvariante außerhalb des Domain Core durchgesetzt"
  (erstes Auftreten dieser Klasse)

## Negativbefunde

- geprüft, ohne Befund: **Race-Freiheit real reproduziert, nicht nur aus
  dem Implementer-Bericht übernommen.** Sperren in `lookupBinding`/
  `AddBinding` (`internal/adapters/driving/replication/mapper/mapper.go`)
  in dieser Sitzung testweise entfernt: `go test -race -run
  TestAssemblerLiveReloadIsRaceFree` schlägt danach real mit `WARNING:
  DATA RACE` fehl, exakt an den Stellen `lookupBinding`
  (`mapper.go:378`, `mapassign_faststr`/`mapaccess2_faststr`) und
  `AddBinding` (`mapper.go:389`) — deckungsgleich mit der
  Implementer-Behauptung. Datei danach aus Backup wiederhergestellt
  (`git status` zeigt keinen Diff mehr an dieser Datei), derselbe Test
  läuft danach real wieder grün. Alle drei `a.tables`-Zugriffsstellen im
  Paket sind ausschließlich `lookupBinding`/`AddBinding`/`RemoveBinding`
  — kein ungesicherter Zugriff verblieben (`grep -n "a\.tables\["` trifft
  nur diese drei Methoden).
- geprüft, ohne Befund: **Regressionsschutz `observeRelation`
  (`slice-032`/`slice-033`) unverändert.** Die einzige Änderung an
  `change`/`observeRelation` ist der Ersatz des direkten Map-Zugriffs
  durch `lookupBinding`/`AddBinding` — keine sonstige Logik verändert.
- geprüft, ohne Befund: **Idempotenz bei Doppelverarbeitung eines Antrags
  über einen Prozess-Neustart hinweg** (`ADR-0050`s Konsequenz) — im
  Quellcode nachvollzogen, nicht nur dem Kommentar geglaubt:
  `EnableTableService.Enable`/`DisableTableService.Disable` sind beide
  idempotent (`TableActivationAdapter.Register`/`Unregister` prüfen den
  Bestand vorab). `applyAdministrationRequest` liest nach `Enable` die
  tatsächlich registrierte Bindung über `Registered`/`CurrentVersion`
  frisch zurück, statt der lokal berechneten `administrationTableID`/
  `administrationSchemaVersionID` blind zu vertrauen — die im DoD
  behauptete „vermeidet Stale-ID-Bugs bei wiederholten/idempotenten
  Anträgen"-Eigenschaft trägt real, auch für eine bereits über
  `CDC_TABLES` aktivierte Tabelle mit abweichender Bindungs-Kennung.
- geprüft, ohne Befund: **Kein Zeitfenster mit unbemerkt bleibenden
  Anträgen bei einem `LISTEN`-Verbindungsabbruch.**
  `runAdministration`s Schleife ruft in **jeder** Iteration zuerst
  `processAdministrationRequests` auf, unabhängig vom Zustand der
  `LISTEN`-Verbindung — ein gestörter Wecksignal-Kanal führt zu
  häufigerem, nicht selteneren Polling (siehe dazu das oben gemeldete
  MEDIUM-Finding zur fehlenden Drosselung, das die Kehrseite dieses
  Verhaltens betrifft).
- geprüft, ohne Befund: **Kein Cross-Goroutine-Race jenseits von
  `Assembler.tables`.** `schemaStore` und `activation` werden sowohl vom
  Capture-Stream als auch von der Administrations-Goroutine verwendet,
  tragen aber beide ausschließlich einen `*pgxpool.Pool` ohne
  zusätzliche unsynchronisierte In-Memory-Zustandshaltung — `pgxpool.Pool`
  ist laut eigenem Vertrag nebenläufigkeitssicher.
- geprüft, ohne Befund: **`spec/architecture.md`-Diagrammkorrektur
  konform zu Hard Rule 3.4** (sprach-/meilensteinfrei). Beide neuen
  Diagramme referenzieren ausschließlich `ARC-*`-IDs (inkl. korrekt
  gewähltem `ARC-007` für den Administrations-Hintergrundzug, da dieser
  in `internal/bootstrap` — Composition Root — lebt), keine Wellen-/
  Slice-/ADR-Bezüge im Diagrammtext selbst.
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7), keine
  Vorwärtsverweise auf `slice-038`.** `grep` über den gesamten Diff nach
  `slice-038` liefert keinen Treffer; kein Kommentar beschreibt eine
  verworfene Alternative, einen abwesenden Text oder bricht mitten im
  Satz ab.
- geprüft, ohne Befund: **`Makefile`/Race-Image sauber.**
  `TOOLCHAIN_RACE_IMAGE` ist digest-gepinnt
  (`golang:1.27@sha256:b475798fb…`), konsistent mit den übrigen
  Toolchain-Referenzen im Repo. Die Zielsignatur von `make test` bleibt
  unverändert (kein Bruch für bestehende Aufrufer), der reale Testlauf
  mit `-race` über die gesamte Suite ist grün — keine durch den
  Race-Detector neu aufgedeckte, unzusammenhängende Data Race.
- geprüft, ohne Befund: **`harness/image-hash.txt`-Aktualisierung
  konsistent** mit dem geänderten Build-Kontext (`internal/bootstrap/`
  geändert → `make image` neu gebaut, neuer Digest committet).
- geprüft, ohne Befund: **End-zu-End-Beleg „ohne Neustart" ist wörtlich
  zutreffend.** `tools/harness/run-integration-tests.sh`s neuer Abschnitt
  „SQL-Administration Live-Reload-Beleg" enthält kein `docker restart`;
  stattdessen prüft er explizit `docker inspect --format
  '{{.State.Running}}'` sowohl nach der Aktivierung als auch nach der
  Deaktivierung, dass der Feed-Container ununterbrochen weiterläuft.
- geprüft, ohne Befund: **DoD-Aktualisierung ehrlich.** Alle sechs mit
  `[x]` markierten Punkte sind durch den Diff und eigene Testläufe in
  dieser Sitzung gedeckt. Die vier offen gelassenen Punkte (Review,
  Closure-Notiz, Risiko-Ausgänge, drei Paarungen) sind korrekt
  unbeansprucht — Planner-Closure-Arbeit nach Modul 8, kein Self-Review.
- geprüft, ohne Befund: **§8 Sub-Area-Sichtung.** Einzige berührte
  Sub-Area `*`/`PGC`, korrekt als GF eingestuft. Der zitierte Treffer
  `BEO-PGC/verwaltung-keine-sql-administration` existiert im Register.
- geprüft, ohne Befund: **Hard Rules 3.1/3.5/3.6.** Jeder Rollout-/
  Testschritt läuft über `docker run`/`make`-Targets, kein lokales
  Toolchain-Install (3.1). Kein Accepted-ADR im Diff verändert, keine
  Gate-Schwelle gelockert (3.5/3.6).
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `ADR-0050`, kein `SPEC-*`/`ARC-*` im Betreff; `make commit-traceability`
  lief in dieser Sitzung über `HEAD~5..HEAD` grün (0 Befunde).
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf,
  eigenständig ausgeführt) — `baseline-verify` (54 Dateien), `d-check`
  (301 Dateien, 0 Befunde), `commit-traceability` (5 Commits, OK),
  `a-check` (0 Befunde) — alle grün.
- geprüft, ohne Befund: **`make test` mit Race-Detector**, real
  ausgeführt (diese Review-Sitzung) — alle Pakete grün, einschließlich
  `internal/bootstrap` und `internal/adapters/driving/replication/mapper`.
- geprüft, ohne Befund: **`make test-integration`, real dreimal
  ausgeführt** (diese Review-Sitzung, nicht aus dem Implementer-Bericht
  übernommen) — jedes Mal Exit 0, jedes Mal beide Log-Zeilen
  „SQL-Administration Live-Reload-Beleg (enable)"/„(disable)" vorhanden,
  jedes Mal `TestMVP*`-Suite vollständig grün. Nebenwirkung
  `tools/schema/plan.yaml` (von `schema-rollout` neu geschrieben) nach
  Abschluss auf HEAD zurückgesetzt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Inline-Suppression ohne ADR/zentrale
Ausnahme" (1×, erstes Auftreten), „fehlende Negativtests bei neuem
öffentlichem Vertrag" (1×, 2. Auftreten nach `review-slice-036` — Schwelle
3× noch nicht erreicht), „unklare Fehlerbehandlung am Rand des
Spec-Bereichs" (1×, erstes Auftreten), „Domäneninvariante außerhalb des
Domain Core durchgesetzt" (1×, erstes Auftreten) — kein Steering-Loop-
Eintrag fällig, keine Klasse erreicht 3×.

## Verdikt

**Merge-blockierend:** Das HIGH-Finding trägt **keinen Rollen-Widerspruch**
— es ist eine Beobachtung gegen eine Hard Rule, keine Implementer-Aussage,
der widersprochen wird. Die Architect-Sequenz aus Modul 8 greift bei
isoliertem HIGH ohne Rollen-Widerspruch nicht zwingend, ist aber die vom
Implementer zu treffende Entscheidung: Suppression entfernen (Kommentar
ohne `//nolint` trägt dieselbe Absicht — `xid` bleibt innerhalb der
Testschleifen-Grenze von 200 Iterationen weit unter jedem
`uint32`-Überlauf) oder die Ausnahme in einer zentralen Konfiguration mit
Begründung dokumentieren. Beides ist eine kleine, lokal begrenzte
Korrektur ohne Rollen-Konflikt.

**Zur zentralen Prüffrage dieses Laufs (Race-Freiheit, `ADR-0050`
Fitness Function Zeile 2):** Bestätigt — eigenständig reproduziert, nicht
nur aus dem Implementer-Bericht übernommen. Ohne die neuen Sperren
schlägt `go test -race` real fehl, exakt an den behaupteten Stellen; mit
ihnen ist die gesamte Suite (inklusive `internal/bootstrap`) grün.

**Zur Idempotenz-Frage** (`ADR-0050`s Konsequenz „ein offener Antrag …
bleibt pending, wird erneut abgeholt"): `EnableTableUseCase`/
`DisableTableUseCase` sind real idempotent, und `applyAdministrationRequest`
liest die Bindung nach `Enable` frisch aus der Datenbank zurück statt der
lokal berechneten Kennung blind zu vertrauen — die Stale-ID-Vermeidung
trägt.

**Zum End-zu-End-Beleg:** Dreifach grün, real ausgeführt, kein
`docker restart` im betreffenden Abschnitt — die DoD-Aussage „ohne
Neustart" ist zutreffend.

**Übergabe:** Ein HIGH-Finding ohne Rollen-Widerspruch, drei MEDIUM-
Findings. Der Implementer entscheidet über Annahme oder Begründung; die
Architect-Sequenz (Modul 8) ist bei diesem Befundbild nicht erforderlich,
kann aber vom Implementer gezogen werden, falls er die HIGH-Einstufung
für unzutreffend hält. Dieser Report ist ein Lauf-Beleg und wird über
Läufe hinweg nicht wieder gelesen — die Summary-Zeile speist bei Bedarf
den Closure-Eintrag (Modul 5). Er ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat.
