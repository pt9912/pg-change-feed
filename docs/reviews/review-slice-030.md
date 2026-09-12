# Review-Report: slice-030 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-030`, §1/§2/§3) und
`AGENTS.md` §3 Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `287b621` (zwei Testfälle + dritte
Test-Tabellenbindung), `79fa509` (DoD-/Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-030-black-box-e2e-schema-aenderungen.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken, §8 Register-Sichtung)
- `spec/lastenheft.md` (`LH-FA-SCH-001`…`005`, Original-Wortlaut, nicht nur die Zusammenfassung des Implementers)
- `spec/pflichtenheft.md` (`SPEC-004` TableSchema/SchemaVersion)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, permanent — Schärft `LH-FA-SCH-004.a`, `SPEC-004`)
- `internal/adapters/driving/replication/mapper/mapper.go` (`Assembler.Consume`, `TableBinding`, `rowImage`)
- `internal/adapters/driving/replication/decode/decode.go` (`tupleValues`)
- `internal/bootstrap/wiring.go:320-336`, `internal/adapters/driving/replication/receive/receive.go:55-63`, `internal/application/port/inbound/verwaltung.go:18-32` (referenzierter „Metadata-Pfad")
- `test/integration/integration_test.go` (Gesamtdatei, für Isolations-Prüfung), `tools/harness/run-integration-tests.sh`, `compose.yaml`
- `docs/plan/planning/observations/BEO-PGC/*/observation.md` (Register-Sichtung gegenprüfen)
- `AGENTS.md` §3 Hard Rules

---

## Findings

Ein HIGH-Finding (Kern dieses Reviews), keine weiteren HIGH/MEDIUM. Ein
LOW-Finding zur DoD-Formulierung.

### F-1 — ADR-0015 (Accepted, permanent) ist nicht implementiert; die neuen Tests beweisen es erstmals real

- `kategorie`: HIGH
- `quelle`: `ADR-0015` (Schema Evolution, Accepted, `permanent`) — Entscheidung
  Option C (TableSchema-/SchemaVersion-Modelle je Change, dynamisch,
  inkompatible Änderungen sichtbar gemeldet), Folgepflicht „SchemaStorePort als
  Outbound Port; Fehlerklasse `schema` für nicht sicher interpretierbare
  Änderungen"
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:53,53-58,146-156` (`TableBinding.SchemaVersion` statisch, `Consume`-Default-Zweig verwirft `*decode.Relation`), `internal/adapters/driving/replication/decode/decode.go:251-270` (`tupleValues`, reiner Text-Passthrough ohne Typ-Auswertung), `internal/bootstrap/wiring.go:328-332` (Kommentar referenziert „Metadata-Pfad", der nirgends existiert) — bestätigend `test/integration/integration_test.go:642-754` (`TestMVPSchemaChangeAddColumn`, `TestMVPSchemaChangeIncompatibleTypeChange`)
- `befund`: Ich habe beide vom Implementer gemeldeten Funde eigenständig am
  Code verifiziert (nicht nur aus dem Bericht übernommen — Zeilenzitate
  geprüft, `grep -rln SchemaStorePort` über das ganze Repo: **0 Treffer**).
  (1) `mapper.TableBinding.SchemaVersion` wird einmalig in `NewAssembler`
  aus der `tables`-Map gesetzt und danach nie mehr verändert;
  `Assembler.Consume` fällt für jedes `*decode.Relation`-Ereignis in den
  `default`-Zweig des Switch (`case decode.Begin/Commit/Change/Truncate`
  sind die einzigen behandelten Typen) und gibt `nil, nil` zurück, ohne die
  Bindung anzufassen — `LH-FA-SCH-005`s Boundary
  („zwei Changes vor und nach einer Schemaänderung … lassen sich
  unterscheiden") ist damit strukturell nicht erfüllbar, nicht nur im
  Testfall zufällig nicht getroffen. (2) `decode.tupleValues` liest jeden
  Spaltenwert ausschließlich aus `pglogrepl.TupleDataTypeText` als String;
  es gibt keine Stelle im Decoder- oder Mapper-Pfad, die eine
  Spalten-`Oid`/einen Typ auswertet — `LH-FA-SCH-004`s Negative-Fall
  („inkompatible Typänderung … erkennbar gemeldet") kann aus diesem Pfad
  nicht entstehen, weil er keinen Typ interpretiert, den er verlieren
  könnte. **Das ist kein bloßer Spec-Rückstand, sondern ein Verstoß gegen
  eine bereits `Accepted` und `permanent` gesetzte ADR:** ADR-0015 hat
  explizit Option A („Relation Metadata 1:1 durchreichen") gegen Option C
  abgewogen und **verworfen**, u. a. weil Option A „keine stabile
  historische Interpretation" trägt — der heutige Code verhält sich exakt
  wie das verworfene Option A, nicht wie das beschlossene Option C. Die
  in ADR-0015 als „Folgepflicht" benannte `SchemaStorePort`-Abstraktion
  existiert im gesamten Repo nicht (0 Treffer). Die vier Code-Kommentare,
  die einen „Metadata-Pfad" referenzieren (`wiring.go:331`, `mapper.go:53`,
  `receive.go:62`, `verwaltung.go:26`), verweisen auf einen Mechanismus,
  der nirgends implementiert ist — sie sind zwar als Rang-Zeiger auf
  `LH-FA-SCH-004.a` formal zulässig (Kommentar-Klasse „Kopplung"/Verweis
  auf eine Spec-Stelle), lesen sich aber wie ein bereits vorhandener Pfad
  und nicht wie eine offene Folgepflicht — das begünstigt genau die
  Fehleinschätzung, die dieser Slice jetzt korrigiert.
  **Wichtig für die Einordnung:** Dieser Zustand ist nicht durch den
  vorliegenden Diff eingeführt — beide betroffenen Dateien
  (`mapper.go`, `decode.go`, `wiring.go` etc.) sind in `287b621`/`79fa509`
  unverändert. Der Diff **beweist** den Verstoß aber erstmals real
  (Compose-Stack, keine Behauptung aus Code-Lektüre) und dokumentiert ihn
  ehrlich statt ihn zu verdecken — das ist die Qualität, die diesen Fund
  überhaupt sichtbar macht. Die Einordnung als HIGH bezieht sich auf den
  **Sachverhalt**, den der Diff aufdeckt, nicht auf die Qualität des Diffs
  selbst (dazu siehe Negativbefunde).
- `verifizierbar`: ja — `grep -rln SchemaStorePort` (0 Treffer, repo-weit);
  `go test ./test/integration/...` kompiliert die neuen Assertions;
  `make test-integration` reproduziert beide Testfälle real gegen den
  Compose-Stack (von mir nicht erneut ausgeführt, siehe Negativbefunde/
  Prüfgrenze).
- `klasse`: „ADR-Verstoß durch schwarzbox-getriebenen Fund — Option C
  beschlossen, Option A implementiert"

**Empfehlung zur Weiterleitung (keine Reviewer-Entscheidung, sondern
Übergabe):** Dieser Fund ist kein Merge-Blocker für den vorliegenden
Test-Diff, aber ein **Closure-Blocker** für `slice-030`: Zwei DoD-Punkte
(`LH-FA-SCH-004`, `LH-FA-SCH-005`) bleiben bei aktuellem Systemstand
strukturell unerfüllbar innerhalb des Slice-Zuschnitts (§1 schließt eine
Änderung des Erkennungsmechanismus explizit aus). Das ist mehr als ein
gewöhnliches „offenes Risiko" nach §6 — es ist ein bestätigter
ADR-Verstoß gegen eine `Accepted`/`permanent` Entscheidung. Nach Modul 8
§Konflikt-Pfad braucht ein HIGH-Finding dieser Art keinen
Rollen-*Widerspruch*, um die Sequenz auszulösen, wenn ein Architect-Verdikt
über den Weg entscheidet, der die DoD-Lücke schließt; ich empfehle dem
Planner, hier den Architect einzubeziehen, bevor er §6/§7 der Closure
schreibt — die drei plausiblen Verdikte sind (a) Folge-Slice(s) mit ADR-
Bezug, der `SchemaStorePort` und die dynamische Re-Versionierung sowie die
Typ-Prüfung tatsächlich liefert, (b) ein Carveout mit Folge-Slice-ID, falls
die Lücke vorerst bewusst offen bleiben soll, oder (c) eine Folge-ADR, die
ADR-0015 zurücknimmt oder abschwächt, falls sich die Anforderung inzwischen
als nicht mehr tragfähig erweist (unwahrscheinlich, da Lastenheft und
Pflichtenheft unverändert `LH-FA-SCH-004`/`005`/`SPEC-004` fordern). Meine
Einordnung ist **HIGH**, nicht „außerhalb des Reviewer-Kontexts": Der
Reviewer prüft explizit gegen ADR (Modul 1 Kernidee, Modul 8 Kernidee),
und die HIGH-Kategorie „ADR-Verstoß" aus dem Reviewer-Skill trifft
unmittelbar zu — unabhängig davon, dass die *Verifikation* der DoD-Punkte
selbst Verifier-Aufgabe bleibt.

### F-2 — DoD-Teilerfüllungs-Vermerk verwendet kein im Slice-Template vorgesehenes Formelement

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-030-black-box-e2e-schema-aenderungen.md:84-101`
- `befund`: Die Checkbox-Items für `LH-FA-SCH-001`/`002`/`005` und
  `LH-FA-SCH-004` tragen einen freitextigen **„Teilweise:"**-Absatz
  innerhalb eines weiterhin offenen `- [ ]`-Punkts. Das Slice-Template
  (`../templates/docs/plan/planning/slice.template.md`) kennt dafür kein
  eigenes Formelement; die Lösung ist inhaltlich korrekt (Checkbox bleibt
  konsequent offen, da nicht vollständig erfüllt) und exakt durch den Diff
  gedeckt, aber die Konvention ist ad hoc. Kein Wiederholungsfall bislang
  bekannt (erstes Auftreten dieser konkreten Form), daher LOW statt MEDIUM.
- `verifizierbar`: nein — Formfrage, kein Gate prüft DoD-Prosa-Konventionen.
- `klasse`: „DoD-Teilerfüllung ohne Template-Formelement"

## Negativbefunde

- geprüft, ohne Befund: **Testqualität `TestMVPSchemaChangeAddColumn`** —
  real black-box: ausschließlich `env.pool.Exec`/`awaitChangesViewRows`
  gegen SQL-DDL bzw. `cdc.changes`, kein interner Paket-Aufruf. Vorher/
  Nachher an derselben aktivierten Tabelle (`feed_mvp_schema`, per
  `newMVPEnv(t, "feed_mvp_schema")` mit eigenem `tableID`-Filter),
  physisch und über `source_table_id` von `feed_mvp_flow`/`feed_mvp_full`
  getrennt — keine Verwechslungsgefahr.
- geprüft, ohne Befund: **Testqualität
  `TestMVPSchemaChangeIncompatibleTypeChange`** — kein Zeremonie-Test: Fall
  1 belegt real, dass PostgreSQL eine inkompatible Typänderung selbst
  ablehnt (Boundary aus §6); Fall 2 belegt real den in F-1 beschriebenen
  Fund (Text-Passthrough ohne Fehlermeldung) als **Tatsachenbeleg**, nicht
  als stillschweigend erwartetes Verhalten — der Funktionskommentar sagt
  explizit, dass der eigentlich intendierte Negative-Fall nicht erreichbar
  ist. Die Abhängigkeit von Bestandsdaten aus `TestMVPSchemaChangeAddColumn`
  (id=1/id=2) ist im Doc-Kommentar offen benannt („nach
  TestMVPSchemaChangeAddColumn") und folgt derselben bereits im Repo
  etablierten Konvention wie `TestMVPDisableRetainedState` („Der Test läuft
  nach den Capture-Läufen … Quell-Reihenfolge"); Go führt Tests einer
  Datei ohne `-run`-Filter und ohne `t.Parallel()` in Deklarationsreihen-
  folge aus, `run-integration-tests.sh` ruft `go test -v
  ./test/integration/...` ohne Filter auf — die Reihenfolge ist damit
  deterministisch gegeben, kein latenter Defekt.
- geprüft, ohne Befund: **Isolation der dritten Tabellenbindung** —
  `compose.yaml`s `CDC_TABLES` trägt `feed_mvp_schema=tbl-mvp-schema:
  sv-mvp-schema`, eigene Tabellen-/Schema-Version-Kennung, keine
  Überschneidung mit `tbl-mvp-flow`/`tbl-mvp-full`; die Aktivierung läuft
  über denselben `EnableTable`-Verdrahtungsweg wie die bestehenden zwei
  Tabellen (`wiring.go`), keine Sonderbehandlung nötig.
- geprüft, ohne Befund: **DoD-Aktualisierung ehrlich, weder über- noch
  unterbehauptet** — `LH-FA-SCH-001`/`002` sind durch
  `TestMVPSchemaChangeAddColumn` tatsächlich real belegt (neue Spalte im
  Row Image danach erfasster Changes, ältere Change unverändert/ohne die
  neue Spalte) und werden korrekt als „teilweise" erfüllt markiert (die
  Checkbox bleibt wegen `LH-FA-SCH-005` insgesamt offen, da der Punkt alle
  drei IDs bündelt). Keine Überbehauptung, keine Unterbehauptung gefunden.
- geprüft, ohne Befund: **Doku-Update-Begründung** — `docs/user/
  benutzerhandbuch.md` nennt `CDC_TABLES` nur generisch (Beispiele
  `orders`/`customers`), keine der Compose-Tabellennamen; die Behauptung
  „kein öffentlicher Vertrag berührt" ist zutreffend.
- geprüft, ohne Befund: **Register-Sichtung §8** — keine der 16 bestehenden
  `BEO-PGC/*`-Beobachtungen (`schema-rollout-fremdobjekte`,
  `spec008-replication-luecke` u. a. geprüft) betrifft die dynamische
  Schema-Version-Bindung oder die Typprüfung im Decoder-/Mapper-Pfad — die
  Aussage „keine Treffer" im Slice-Kopf §8 stimmt. **Dieser Review-Lauf
  selbst erzeugt keinen neuen Registereintrag** — das ist Planner-Arbeit
  bei der Closure (§7), nicht Teil dieses Reviews; ich weise nur darauf
  hin, dass F-1 dort einen neuen `BEO-PGC/`-Eintrag begründet.
- geprüft, ohne Befund: **Kompilierbarkeit** — `make test` (`go test
  ./...`, netzloser Toolchain-Container) lief grün einschließlich
  `test/integration` (0.002s — beide neuen Tests übersprungen mangels
  `CDC_INTEGRATION_DSN`, aber kompiliert); `gofmt -l
  test/integration/integration_test.go` im selben Toolchain-Image: keine
  Ausgabe (bereits formatiert). Beide Läufe eigenständig in dieser
  Review-Sitzung ausgeführt.
- geprüft, ohne Befund: **Traceability** — beide Commit-Betreffs tragen
  `LH-FA-SCH-*`-Kennungen, keine `SPEC-*`/`ARC-*`-Kennung im Betreff;
  `make commit-traceability` lief in dieser Sitzung über `HEAD~5..HEAD`
  grün (0 Befunde).
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — keine lokale
  Toolchain-Installation; `compose.yaml`/`run-integration-tests.sh`-Änderung
  bleibt innerhalb des bestehenden Docker-/Compose-Musters.
- geprüft, ohne Befund: **Hard Rule 3.2** (Suppression-Verbot), **3.5**
  (ADR-Immutabilität), **3.6** (Gate-Lockerung) — keine Berührung; kein
  `#noqa`/`//nolint`, kein ADR-Inhalt editiert (ADR-0015 bleibt
  unangetastet — der Diff zitiert sie nur implizit über den Code-Befund),
  keine Gate-Schwelle geändert.
- geprüft, ohne Befund: **Hard Rule 3.7** (Kommentar-Klassen) — die beiden
  neuen Funktionskommentare beschreiben durchweg den geltenden Zustand
  indikativ (Grenze-/Zusage-Klasse: „bleibt dabei unbelegt", „ist … nicht
  real herstellbar"), keine verworfene Alternative, kein abgebrochener
  Satz, keine Chronik. Die vier **vorbestehenden** „Metadata-Pfad"-
  Kommentare (`wiring.go:331` u. a., nicht Teil dieses Diffs) sind unter
  F-1 separat gewürdigt, aber nicht als eigener Hard-Rule-3.7-Verstoß
  gezählt, da sie außerhalb dieses Diffs liegen.
- geprüft, ohne Befund: **Hard Rule 3.3** (git-mv/Inhalt-Trennung) — nicht
  einschlägig, beide Commits sind reine Content-Commits ohne Umbenennung.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify` (v6.5.0, 54 Dateien OK), `d-check` (259 Dateien, 0
  Befunde), `commit-traceability` (OK, `HEAD~5..HEAD`), `a-check` (0
  Befunde) — alle grün.
- **Prüfgrenze — nicht reproduzierbar in dieser Review-Sitzung:** Der
  Implementer-Bericht behauptet drei aufeinanderfolgende grüne
  `make test-integration`-Läufe („Commit folgt"). `make test-integration`
  ist kein Gate, braucht eine laufende Compose-Umgebung mit echter
  logischer Replikation und wurde in dieser Review-Sitzung nicht
  eigenständig reproduziert (Ressourcen-/Zeitaufwand, außerdem
  Verifier-Aufgabe nach Modul 8/11). Stattdessen unabhängig geprüft:
  Kompilierbarkeit (`make test`), `gofmt`, sowie beide gemeldeten Funde
  eigenständig am Code nachvollzogen (nicht nur aus dem Bericht
  übernommen) — siehe F-1. Konsistent mit `review-slice-029.md`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** ADR-Verstoß durch schwarzbox-getriebenen
Fund — Option C beschlossen, Option A implementiert · DoD-Teilerfüllung
ohne Template-Formelement

## Verdikt

**Merge-blockierend:** nein, für den vorliegenden Test-Diff selbst — die
beiden neuen Testfälle sind real black-box, sauber isoliert, ehrlich
dokumentiert und decken die tatsächliche Systemgrenze auf, statt sie zu
verdecken; DoD-Häkchen sind exakt durch den Diff gedeckt (keine Über- oder
Unterbehauptung). F-2 ist eine kleine Formfrage ohne Merge-Blockade.

**Closure-blockierend:** ja, im engeren Sinn — `slice-030` kann nach
aktuellem Plan-Zuschnitt (§1 schließt eine Änderung des
Erkennungsmechanismus aus) nicht mit vollständiger DoD nach `done/`
wandern, weil `LH-FA-SCH-004`/`005` strukturell unerfüllt bleiben (F-1).
Das ist keine Reviewer-Entscheidung über den *Weg* (Carveout,
Folge-Slice(s), Folge-ADR) — das bleibt Planner/Architect vorbehalten
(Modul 5 §Offene Risiken werden bei Closure aufgelöst, Modul 8
§Konflikt-Pfad) —, aber die **Schwere** des Funds selbst (ADR-Verstoß
gegen eine `Accepted`/`permanent` Entscheidung, zwei unerfüllte
Lastenheft-Akzeptanzkriterien, bisher unentdeckt) ordne ich als HIGH ein,
nicht als INFO/„außerhalb des Reviewer-Kontexts" — der Reviewer prüft
explizit gegen ADR und Plan, und genau das trifft hier zu.

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** Beide gemeldeten Funde am Code selbst nachvollzogen
(`mapper.go`, `decode.go`, `wiring.go`/`receive.go`/`verwaltung.go`),
inklusive repo-weiter Suche nach `SchemaStorePort` (0 Treffer) und ADR-0015
im Original gelesen (Optionen-Vergleich, Folgepflicht); Kompilierbarkeit
(`make test`, `gofmt`); Traceability; Register-Sichtung gegen alle 16
bestehenden `BEO-PGC/*`-Einträge; `make gates` komplett neu ausgeführt.
Die 3×-`make test-integration`-Behauptung bleibt als Prüfgrenze explizit
unbestätigt (siehe Negativbefunde).

**Übergabe:** F-1 geht — über den Implementer hinaus — an Planner/Architect
zur Entscheidung über den Weg vor der Closure (Carveout, Folge-Slice(s)
mit ADR-Bezug, oder Folge-ADR); F-2 geht an den Implementer/Planner zur
freien Entscheidung (LOW). Die **Finding-Klassen** gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers.
Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität
(einschließlich der Reproduktion von `make test-integration`) prüft der
Verifier separat.
