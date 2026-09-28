# Review-Report: fix(bootstrap) EnableTable/DisableTable-Assembler-Sync — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (kritischer Datenverlust-Fix),
Hard Rules und Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `a40b4809` — `fix(bootstrap): EnableTable/DisableTable
über HTTP/gRPC aktualisieren den laufenden Assembler (LH-FA-CFG-001,
LH-FA-CFG-002)`, 8 Dateien.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09-Schärfung (Commit
`17b5efd4`-Baum)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext:**

- `docs/reviews/review-example-kotlin-grpc-client-verbmatrix.md` INFO-1 (die
  Originalmeldung, real reproduziertes Protokoll des Bugs)
- `internal/application/usecase/enable/service.go`,
  `internal/application/usecase/disable/service.go`
- `internal/adapters/driving/replication/mapper/mapper.go`
  (`AddBinding`/`RemoveBinding`/`lookupBinding`/`change`)
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
  [`ADR-0057`](../plan/adr/0057-http-grpc-api.md),
  [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)
- `AGENTS.md` §3 (insb. §3.1, §3.7, §3.9, §3.12, §3.13)
- `.harness/skills/reviewer.md` vollständig
- [`LH-FA-CFG-001`](../../spec/lastenheft.md), [`LH-FA-CFG-002`](../../spec/lastenheft.md)

---

## Findings

### F-1 — Vorher/Nachher-Chronik in Produktionscode-Kommentar (`syncAssemblerAddBinding`)

- `kategorie`: HIGH
- `quelle`: Reviewer-Skill HIGH „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar" (`AGENTS.md` §3.7)
- `pfad`: `internal/bootstrap/assemblersync.go:18-23`
- `befund`: Der Godoc-Kommentar über der Produktionsfunktion
  `syncAssemblerAddBinding` erzählt explizit die Vorgeschichte des Bugs in
  der Vergangenheitsform statt den heutigen Zustand zu beschreiben: „…
  aufrufen — **beide riefen bislang** denselben `EnableTableUseCase` auf,
  aber nur der Antragsqueue-Pfad trug den Nachtrag in den laufenden Prozess;
  eine über HTTP oder gRPC aktivierte Tabelle **blieb bis zum nächsten
  Neustart unerfasst**". Das Wort „bislang" ist genau der von diesem Repo
  wiederholt benannte Vorher/Nachher-Marker; Subjekt des Satzes ist der
  Produktionscode-Pfad selbst („beide", „eine … Tabelle"), nicht ein
  Testfall — die im Skill genannte Ausnahme für Test*-Godocs greift hier
  nicht. Die am Satzende zitierte `LH-FA-CFG-001` rettet den Fund nicht: die
  Kennung steht als Anhängsel hinter der Chronik-Erzählung, nicht an ihrer
  Stelle. Eine schwächere zweite Instanz desselben Musters steht im
  Kommentar über `enableTableWithAssemblerSync` (Zeile 68: „ohne diesen
  Nachtrag verliert eine über diesen Weg aktivierte Tabelle jede Änderung
  bis zum nächsten Prozess-Neustart") — dort in Präsens/Indikativ als
  Konsequenz-Zusage formuliert, näher an der zulässigen Form, aber mit
  identischer „Ohne-X-Konsequenz"-Rhetorik, die `AGENTS.md` §3.7 als
  Falsch-Beispiel führt. Eine dritte, schwächste Instanz steht als
  Testabdeckungs-Beschreibung in `tools/harness/run-integration-tests.sh`
  (`abdeckung_declare`-Aufruf für „Direkter Zugriffsweg Live-Reload (gRPC
  enable)"): „… jetzt für den **bis dahin nicht nachgetragenen** direkten
  Weg belegt" — dieselbe Marker-Klasse, in einem Runner-Skript (Skopus
  ausdrücklich eingeschlossen), aber niedrigeren Gewichts als die erste
  Instanz, weil sie eine Testaussage statt eine Produktionscode-Semantik
  beschreibt.
- `verifizierbar`: nein — kein Gate/Sensor liest Vorher/Nachher-Sprache
  (`make kommentar-kennungen DIFF=17b5efd4` prüft nur Kennungsform/-kette,
  lief mit 0 Kandidaten); die Probe ist ausschließlich das Lesen des
  Kommentars gegen sein Satzsubjekt.
- `klasse`: Slice-/Wellen-Chronik in Produktionscode-Kommentar
  (Vorher/Nachher-Sprache)

### F-2 — Drei Doc-Kommentare außerhalb des Diffs beschreiben jetzt einen falschen Alleinstellungsanspruch

- `kategorie`: LOW (Träger außerhalb des Diffs — `AGENTS.md` §3.13)
- `quelle`: `AGENTS.md` §3.13 „Eine Arbeit, die eine beschriebene
  Eigenschaft bewegt, zieht ihre Träger nach"
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:438-445`
  (`AddBinding`), `internal/adapters/driving/replication/mapper/mapper.go:642-647`
  (`RemoveBinding`), `internal/adapters/driving/replication/receive/receive.go:234-241`
  (`Stream.Assembler()`)
- `befund`: Alle drei — keiner davon Teil dieses Diffs — behaupten
  wörtlich, dass `AddBinding`/`RemoveBinding` ausschließlich „nachdem eine
  über SQL beantragte Aktivierung/Deaktivierung real ausgeführt wurde"
  aufgerufen werden bzw. dass die Administrations-Verarbeitung der
  Composition Root der alleinige Aufrufer sei. Seit diesem Fix stimmt das
  nicht mehr: `enableTableWithAssemblerSync`/`disableTableWithAssemblerSync`
  rufen dieselben Methoden auf demselben `*mapper.Assembler` jetzt auch aus
  dem direkten HTTP-/gRPC-Zugriffsweg auf, ohne die SQL-Antragsqueue zu
  durchlaufen — die dokumentierte Eigenschaft „wie diese Methoden erreicht
  werden" hat sich bewegt, die drei Kommentarstellen sind nicht nachgezogen.
  Eine vierte, schwächere Instanz mit derselben Wurzel steht im
  Testfile-Kommentar `internal/bootstrap/administration_internal_test.go:220-224`
  (`fakeTableActivationPort`).
- `verifizierbar`: nein — Lese-Handlung, kein Sensor deckt Kommentar-Genauigkeit.
- `klasse`: Träger-Nachzug außerhalb des Diffs (§3.13)

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/assemblersync.go` — Mechanismus,
  Fehlerpfade, Synchronisations-Grenzen (`tablesMu` in `mapper.go` deckt
  beide Aufrufpfade ab, keine neue Race-Möglichkeit)
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` — Extract-Refactor von
  `applyAdministrationRequest`s Enable-/Disable-Zweig ist verhaltensgleich;
  `enableTableAPI`/`disableTableAPI` wrappen dieselben, bereits verdrahteten
  `enableTables`/`disableTables`-Instanzen, keine doppelte Registrierung
  möglich; `qualified`-Variable bleibt für die übrigen Zweige (Exclude-/
  IncludeColumn) in Gebrauch; der Prozessstart-Seed (`enableTables.Enable`
  am CDC_TABLES-Vorlauf) bleibt bewusst undekoriert, weil
  `activatedTableBindings` den vollständigen DB-Stand danach ohnehin
  frisch liest
- geprüft, ohne Befund: Adapter-Grenzen — `internal/adapters/driving/http`
  und `internal/adapters/driving/grpc` importieren `replication/mapper`
  nicht (`git grep` bestätigt 0 Treffer); `make a-check` bestätigt 0 Befunde
- geprüft, ohne Befund: `internal/bootstrap/assemblersync_internal_test.go`
  — 5 Tests, eingangsseitig gebunden; Mutationsprobe bestätigt: das
  Entfernen des `syncAssemblerAddBinding`-Aufrufs aus `Enable` (Reproduktion
  des Original-Bugs) lässt `TestEnableTableWithAssemblerSyncAddsBindingOnSuccess`
  und `TestEnableTableWithAssemblerSyncFailsWhenSyncLookupFails` real rot
  werden (`go test -run TestEnableTableWithAssemblerSync -v` im
  Toolchain-Container), Rücknahme per `git checkout` verifiziert sauber
- geprüft, ohne Befund: `internal/bootstrap/administration_internal_test.go`
  — unverändert, keine bestehende Whitebox-Assertion angepasst
- geprüft, ohne Befund: `tools/harness/httpclient/main.go`,
  `tools/harness/grpcenabletableclient/main.go` — Kommentare je Block
  höchstens eine Kennung, Aufrufform konsistent mit bestehenden
  Wegwerf-Clients
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` — die drei
  neuen Phasen (gRPC enable, gRPC disable, HTTP enable) real durchlaufen,
  Erfassung/Nicht-Erfassung und Feed-Container-Weiterlauf jeweils geprüft
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version 1.77→1.78,
  neue Zeile in `### Änderungshistorie`, beide neuen Hinweisabsätze sachlich
  im Präsens, kein ADR-/Review-Verweis im Fließtext, keine Chronik-Sprache
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` — Diff konsistent mit
  den in `run-integration-tests.sh` neu deklarierten `abdeckung_declare`-Zeilen,
  Zeilennummern-Verschiebung durch die neu eingefügten Zeilen plausibel
- geprüft, ohne Befund: `make kommentar-kennungen DIFF=17b5efd4` — 0
  Kandidaten (Kennungsform/-kette; deckt keine Vorher/Nachher-Sprache, s. F-1)
- geprüft, ohne Befund: `make gates`, `make generated-sync`, `make a-check`,
  `make fmt-check`, `make docs-check`, `make test` — alle Exit 0, direkt
  geprüft (kein Pipe/Wrapper dazwischen)
- geprüft, ohne Befund: `make test-integration` (voller, realer Lauf) — alle
  Phasen inkl. der drei neuen bestätigt, kein unerwarteter Fehlschlag, Lauf
  endet mit „Lauf abgeschlossen"

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Slice-/Wellen-Chronik in
Produktionscode-Kommentar (Vorher/Nachher-Sprache) · Träger-Nachzug
außerhalb des Diffs (§3.13)

## Restrisiko 2 (Implementer-Meldung) — eigene Bewertung

Der Implementer hat gemeldet: Scheitert `syncAssemblerAddBinding` (die vier
DB-Lesungen `Registered`/`CurrentVersion`/`ExcludedColumns`/
`TransformationRules`), *nachdem* der innere `EnableTableUseCase.Enable`
bereits erfolgreich committet hat, liefert `enableTableWithAssemblerSync`
einen Fehler (`500`/`Internal`) zurück, obwohl die Tabelle in der DB schon
aktiviert ist — ein ausbleibender Wiederholungsaufruf hinterließe denselben
stillen Datenverlust wie der ursprüngliche Bug, nur auf ein schmaleres
Zeitfenster verengt.

**Eigene Einschätzung, nicht übernommen:** Dieses Risiko ist real, aber
**kein neues Risiko dieses Fixes** — es ist exakt dieselbe Fehlersemantik,
die der SQL-Antragsqueue-Pfad (`applyAdministrationRequest`) bereits seit
`ADR-0050` trägt: derselbe extrahierte Code (`syncAssemblerAddBinding`)
lief dort vorher inline, mit identischem Fehlerpfad, und endet dort
ebenfalls im `failed`-Vermerk ohne automatischen Retry. Der Unterschied
liegt allein in der Sichtbarkeit des Fehlschlags: Am Antragsqueue-Pfad
bleibt ein gescheiterter Antrag als `failed`-Zeile in
`cdc.administration_request` sichtbar für Monitoring/Betreiber; am direkten
HTTP-/gRPC-Pfad erhält der aufrufende Client einen synchronen `500`/
`Internal`-Fehler, aber es gibt kein serverseitiges Äquivalent zur
`failed`-Zeile, falls der Client diesen Fehler nicht auswertet oder nicht
wiederholt. Die vier Lese-Operationen laufen unmittelbar nach einem soeben
erfolgreichen Schreibzugriff auf dieselbe Datenbank — das Fenster für einen
Fehlschlag ist schmal (transiente Netzwerk-/Verbindungsstörung im
Millisekunden-Bereich zwischen Commit und den vier Folge-Lesungen), aber
nicht null.

**Einstufung:** Kein blockierendes HIGH-Finding (es reproduziert einen
bestehenden, bereits akzeptierten Fehlermodus, führt keinen neuen
Fehlerpfad ein, und der Aufrufer bekommt ein explizites Fehlersignal statt
eines stillen Erfolgs). Es qualifiziert sich als MEDIUM „unklare
Fehlerbehandlung am Rand des Spec-Bereichs" — dem direkten Zugriffsweg
fehlt ein Betreiber-sichtbares Äquivalent zum `failed`-Vermerk der
Antragsqueue für genau diesen (seltenen) Fehlerfall, und
`docs/user/benutzerhandbuch.md` dokumentiert nicht, dass ein `500` auf
`/tables/enable` einen Wiederholungsaufruf erfordert (der dank
`AlreadyEnabled=true` sicher idempotent ist). Da der Implementer dieses
Risiko selbst bereits benannt hat und es hier lediglich bewertet statt neu
entdeckt wird, führe ich es **nicht** als eigenständiges nummeriertes
Finding, sondern als benannte, nicht blockierende Beobachtung — eine
Behebung (Betreiber-Hinweis im Handbuch, oder ein Log-Eintrag mit
eindeutigem Sentinel bei diesem Fehlerpfad) ist ein sinnvoller, aber
optionaler Folge-Schritt.

## Verdikt

**Merge-blockierend:** **Ja, für F-1 (HIGH)** — nach Reviewer-Skill-Regel
blockieren HIGH-Findings typischerweise. Der Fund ist eine reine
Kommentar-Textkorrektur (kein Verhaltens-, Test- oder Wiring-Fix nötig) und
damit ohne Risiko für die bereits real verifizierte Korrektheit des Fixes
selbst zu beheben — dennoch bleibt die Klassifikation HIGH nach der
expliziten Skill-Regel, unabhängig vom Umfang der Korrektur.

**Fachliche Korrektheit des Fixes selbst:** Vollständig bestätigt, ohne
Abkürzung: Mechanismus ist ein einzelner geteilter Codepfad (`git grep`
bestätigt genau 2 tatsächliche `assembler.AddBinding(`/
`assembler.RemoveBinding(`-Aufrufstellen, beide in `assemblersync.go`); der
SQL-Antragsqueue-Pfad ist ein reiner, verhaltensgleicher Extract-Refactor
(bestätigt durch unverändert grünen `administration_internal_test.go`);
keine neue Adapter-zu-Adapter-Kante (`a-check`: 0 Befunde, `git grep`: 0
Treffer für `replication/mapper` in den HTTP-/gRPC-Adaptern); die 5 neuen
Unit-Tests sind eingangsseitig gebunden (eigene Mutationsprobe bestätigt
rot bei entferntem Nachtrag, sauber zurückgenommen); der volle reale
`make test-integration`-Lauf bestätigt alle drei neuen E2E-Phasen
(gRPC-Enable/-Disable, HTTP-Enable) mit den exakt vom Implementer zitierten
Log-Zeilen, ohne dass ein Neustart nötig war; `make gates` (inkl.
Coverage-Gate bei 80,50 % gegen Schwelle 80 %), `make generated-sync`,
`make a-check`, `make fmt-check`, `make docs-check` und `make test` liefen
alle mit Exit 0, jeweils direkt geprüft.

**Übergabe:** F-1 (und optional F-2 im selben Zug) gehen an den Implementer
zur Fixrunde — reine Kommentar-Korrektur an vier Stellen (zwei in
`assemblersync.go`, eine in `run-integration-tests.sh`, optional die drei
Fremd-Doc-Stellen in `mapper.go`/`receive.go`/`administration_internal_test.go`
laut §3.13-Meldung). Restrisiko 2 bleibt eine benannte, nicht blockierende
Beobachtung. Dieser Report ist ein Lauf-Beleg; DoD-/Spec-Konformität prüft
der Verifier separat (Modul 11).
