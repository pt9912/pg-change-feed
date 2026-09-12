# Review-Report: slice-017 — 2026-09-12

**Review-Art:** Code — geprüft gegen Slice-Plan +
[`ADR-0040`](../plan/adr/0040-clockport.md) (Maintainability).

**Gegenstand:** Commit `3fa7e5d`
(`feat(replication): Commit-Zeitstempel durch Decoder, Mapper, Domäne
spiegeln (LH-FA-ADM-004)`) — `internal/adapters/driving/replication/decode/{decode.go,decode_test.go}`,
`internal/adapters/driving/replication/mapper/{mapper.go,mapper_test.go}`,
`internal/domain/model/{transaction.go,transaction_test.go}`,
`internal/bootstrap/heartbeat_internal_test.go`,
`internal/adapters/driven/postgresstorage/store_test.go`,
`internal/application/usecase/capture/service_test.go`, sowie der
Slice-Plan selbst
(`docs/plan/planning/in-progress/slice-017-commit-zeitstempel-decoder-mapper-domaene.md`,
§3 Plan-Nachzug). **Zum Prüfzeitpunkt bereits committet** — entgegen der
Ankündigung im Auftrag ("Implementer hat … noch NICHT committet") lag der
Arbeitsbaum sauber (`git status`: nichts zu committen) und der Diff als
Commit `3fa7e5d` vor; geprüft wurde `git show 3fa7e5d`, nicht der
Arbeitsbaum.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-017-commit-zeitstempel-decoder-mapper-domaene.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Nachzug, §4 Trigger, §6
  Risiken, §8 Sub-Area-Prüfung)
- [`ADR-0040`](../plan/adr/0040-clockport.md) (ClockPort — `internal/domain`
  importiert `time` nicht; Fitness Function laut ADR selbst "Gate geplant",
  aktuell nicht in `.a-check.yml` verdrahtet)
- [`LH-FA-ADM-004`](../../spec/lastenheft.md) (messbarer CDC-Abstand)
- `internal/domain/model/timepoint.go` (`TimePoint`/`NewTimePoint`-Vertrag)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.3 Move/Inhalt-Trennung,
  §3.7 Kommentar-Klassen) · `harness/conventions.md` (MR-000, genau eine
  Sub-Area `PGC`, Greenfield)
- Commit-Traceability (`AGENTS.md` §5, `ADR-0045`)
- Externe Quelle (Diff-Verweis geprüft, nicht Bestandteil des Repos):
  `github.com/jackc/pglogrepl@4ae5c490f7ce97fb0b918eed3023b35c16990f5a`
  (gepinnte Version aus `go.mod`), Dateien `message.go`, `pglogrepl.go`
- Vorherige Findings am gleichen Modul: keine (erster Diff an
  `decode`/`mapper`/`transaction.go` seit Beginn des Beobachtungs-Registers)

---

## Prüfungen im Detail

**1. `TimePoint`-Verwendung / `ADR-0040`.** `internal/domain/model/transaction.go`
importiert weiterhin kein `"time"` (eigener `grep` über das gesamte
`internal/domain`- und `internal/application/usecase`-Verzeichnis: keine
Treffer). `model.NewTimePoint(unixNanos int64) TimePoint` (`timepoint.go:15`)
nimmt exakt den Typ, den `mapper.go:120`
(`model.NewTimePoint(event.CommitTime.UnixNano())`) liefert — Konstruktor
und Aufruf passen zusammen, kein Bruch der `ADR-0040`-Grenze.

**2. Plan-Nachzug-Bewertung (vier Test-Call-Sites).** Eigene Prüfung der vier
Stellen (`heartbeat_internal_test.go:134`, `store_test.go:117,439`,
`capture/service_test.go:105`): jede Änderung ist ein einzeiliger
Parameter-Zusatz (`tx.Commit(position)` → `tx.Commit(position,
model.NewTimePoint(1))`), keine Design- oder Assertion-Änderung. Gemessen an
§Ziel-Form: Slice (Baseline-Regelwerk `modul-05-planning-harness.md`,
"zu groß, wenn … mehr als drei Liefer-Punkte … mehrere Schichten betroffen …
nicht in einer Review-Sitzung prüfbar"): Der Nachzug fügt keinen
Liefer-Punkt hinzu (reiner Compile-Zwang derselben Signatur-Änderung),
berührt keine neue Schicht über die bereits geplanten hinaus (Domäne bleibt
die einzige Schicht mit Verhaltensänderung, die vier Dateien sind reine
Test-Fixtures in bereits vom Signatur-Wechsel betroffenen Paketen), und war
in dieser Sitzung in unter einer Minute vollständig nachvollziehbar. **Ich
teile die Einstufung des Implementers**: kein Rückführungs-Trigger nach §4
("unerwartet viele Call-Sites … sprengt den Umfang"). Kein
Rollen-Widerspruch, keine Sequenz nach Modul 8 nötig.

**3. `TestDecodeFlowToCapture` real?** Ja — eigener Testlauf
(`go test ./internal/adapters/driving/replication/decode/... -run
DecodeFlowToCapture -v`, PASS) bestätigt: Der Test baut vier
Binär-Payloads (`beginPayload`, `relationPayload`, `insertPayload`,
`commitPayload`) von Hand zusammen, füttert sie einzeln durch
`decode.Decoder.Decode` → `mapper.Assembler.Consume`, und prüft am Ende
`command.Transaction.SourceCommittedAt()` gegen den aus `commitMicros`
erwarteten Wert. Das ist ein echter Byte-Durchlauf durch alle drei
Schichten, keine Attrappe.

**4. Zeitzonen-Behauptung (§6-Risiko 1).** Geprüft gegen den gepinnten
`pglogrepl`-Quellcode (`pglogrepl.go:845-850`):
`func pgTimeToTime(microsecSinceY2K int64) time.Time { … return
time.Unix(0, microsecSinceUnixEpoch*1000) }`. Der Aufruf `time.Unix(0,
nsec)` ist damit belegt. **Die Zusatz-Behauptung "UTC-verankert" ist so
nicht gedeckt**: Go-Stdlib-Semantik von `time.Unix` liefert einen Wert mit
`Location = Local` (System-Zeitzone), nicht `UTC` — die Rückgabe ist gerade
*nicht* UTC-verankert. Was tatsächlich trägt, ist eine andere Eigenschaft:
`time.Time.UnixNano()` (verwendet in `mapper.go:120`) und `time.Time.Equal()`
(verwendet in den Tests) sind beide *Location-unabhängig*, weil sie auf dem
absoluten Zeitinstant rechnen, nicht auf der gespeicherten `Location`. Der
Code ist dadurch korrekt — aber aus einem anderen Grund als dem berichteten.
Siehe F-1.

**5. Bestehende Domain-Invarianten.** `transaction_test.go`: die beiden
Fälle "doppelter Commit" (`ErrTransactionAlreadyCommitted`) und "Commit an
fremder Quelle" (`ErrSourceMismatch`) sind in ihrer Assertion unverändert;
einzige Änderung ist der zusätzliche `NewTimePoint(...)`-Parameter am
`Commit()`-Aufruf. `TestLHFACAP006OpenTransactionIsNotConsumable` bekam eine
zusätzliche Prüfung (`SourceCommittedAt()` an offener Transaktion liefert
`committed=false`), die bestehende Prüfungen unverändert lässt.

**6. Hard Rules.** Siehe Negativbefunde unten (Traceability, Docker-only,
Suppression, Kommentar-Klassen) — je einzeln geprüft, kein Verstoß
gefunden.

**7. Eigene Reproduktion der Mutation-Test-Behauptung.** In dieser Sitzung
selbst nachvollzogen (Datei danach zurückgesetzt, `git status` sauber):

- Mapper-Zeile auf `model.NewTimePoint(0)` gepatcht → `go test ./decode/...
  ./mapper/...` liefert **zwei** rote Tests (`TestDecodeFlowToCapture`,
  `TestConsumeFullTransaction`) — deckt sich mit dem Bericht.
- Domain-Zuweisung `t.sourceCommittedAt = sourceCommittedAt` in
  `transaction.go` entfernt → `go test ./model/... ./decode/...
  ./mapper/...` liefert **drei** rote Tests
  (`TestCommitCarriesSourceCommittedAt`, `TestDecodeFlowToCapture`,
  `TestConsumeFullTransaction`), **nicht nur einen** wie berichtet ("ein
  Test lief rot"). Siehe F-2.

**8. `make gates` / `make test` unabhängig nachvollzogen.** Beide in dieser
Sitzung selbst ausgeführt (nicht nur die Implementer-Behauptung
übernommen): `make test` → alle Pakete `ok`; `make gates` →
`baseline-verify` (v6.5.0, 54 Dateien), `docs-check`/d-check (181 Dateien, 0
Befunde, zweimal — Standard-Lauf und `commits`-Modul über `HEAD~5..HEAD`),
`commit-traceability` (5 Commits, Betreffs ohne Struktur-ID),
`a-check` (0 Befunde) — alle grün, deckt sich mit dem Bericht.

---

## Findings

### F-1 — Risiko-Begründung "UTC-verankert" für `pglogrepl.CommitMessage.CommitTime` ist technisch ungenau

- `kategorie`: MEDIUM
- `quelle`: Maintainability (§6-Risiko-1-Begründung des Slice-Plans, noch
  nicht in §7 geschrieben, aber laut Bericht die vorgesehene
  Closure-Begründung)
- `pfad`: extern — `github.com/jackc/pglogrepl@4ae5c490f7ce97` `pglogrepl.go:847-850`
  (`pgTimeToTime`); im Repo betroffen: die noch zu schreibende §7-Notiz zu
  §6-Risiko 1 in
  `docs/plan/planning/in-progress/slice-017-commit-zeitstempel-decoder-mapper-domaene.md`
- `befund`: `pgTimeToTime` liefert `time.Unix(0, nsec)`; laut Go-Stdlib-
  Semantik trägt das Ergebnis die `Local`-Location, nicht `UTC` — die
  Formulierung "UTC-verankert" behauptet das Gegenteil dessen, was die
  Bibliothek tut. Der Code selbst ist trotzdem korrekt, weil ausschließlich
  `.UnixNano()` (Mapper) und `.Equal()` (Tests) verwendet werden — beide
  Methoden sind Location-unabhängig, unabhängig davon, ob der Zeitpunkt als
  UTC oder Local markiert ist. Träfe die falsche Begründung unverändert in
  die Closure-Notiz ein, stünde dort eine Invariante ("UTC-verankert"), auf
  die sich künftiger Code fälschlich verlassen könnte (z. B. ein direktes
  `.Format()` auf `CommitTime`/`SourceCommittedAt` ohne vorheriges `.UTC()`,
  im Vertrauen auf die hier behauptete Verankerung).
- `verifizierbar`: ja — durch Lesen des gepinnten `pglogrepl`-Quellcodes und
  der Go-Stdlib-Dokumentation zu `time.Unix`; kein Gate prüft das
  automatisiert.
- `klasse`: Technische Begründung einer Risiko-Closure benennt die falsche
  Ursache für ein korrektes Verhalten (Ort-/Zeitzonen-Semantik einer
  externen Bibliotheksfunktion falsch charakterisiert).

### F-2 — Mutation-Test-Bericht unterschätzt den Blast-Radius der Domain-Mutation

- `kategorie`: LOW
- `quelle`: Maintainability (Berichtsgenauigkeit)
- `pfad`: `internal/domain/model/transaction.go:69`
  (`t.sourceCommittedAt = sourceCommittedAt` in `Commit`)
- `befund`: Der Bericht nennt für die entfernte Domain-Zuweisung "ein Test
  lief rot". Eigene Reproduktion in dieser Sitzung zeigt **drei** rote Tests
  (`TestCommitCarriesSourceCommittedAt` im Domain-Paket,
  `TestDecodeFlowToCapture` im Decode-Paket, `TestConsumeFullTransaction`
  im Mapper-Paket) — die Mutation propagiert über alle drei Schichten, weil
  der Ende-zu-Ende-Test und der Mapper-Test denselben Domänenzustand
  prüfen. Die Unterschätzung ist in diese Richtung folgenlos (die reale
  Abdeckung ist besser als berichtet), macht den Bericht aber ungenau.
- `verifizierbar`: ja — in dieser Sitzung reproduziert, Arbeitsbaum danach
  wieder sauber.
- `klasse`: Mutationstest-Bericht unterschätzt Blast-Radius.

---

## Negativbefunde

- geprüft, ohne Befund: **`TimePoint`/`ADR-0040`-Grenze** — kein
  `"time"`-Import in `internal/domain` oder
  `internal/application/usecase/*`; `NewTimePoint`/Konstruktor-Signatur
  passt zum Mapper-Aufruf.
- geprüft, ohne Befund: **Plan-Nachzug-Größenbewertung** — vier
  Test-Call-Sites sind mechanisch, ohne neue Schicht- oder Liefer-Punkt-
  Berührung; Einstufung des Implementers ("kein Rückführungs-Trigger")
  geteilt, kein Rollen-Widerspruch.
- geprüft, ohne Befund: **`TestDecodeFlowToCapture` real** — eigener
  Testlauf bestätigt echten Byte-Durchlauf durch Decoder, Mapper und
  Domäne inklusive Zeitstempel-Prüfung.
- geprüft, ohne Befund: **Domain-Invarianten unverändert** —
  `ErrTransactionAlreadyCommitted`/`ErrSourceMismatch`-Tests unverändert in
  ihrer Aussage, nur um den neuen Parameter ergänzt.
- geprüft, ohne Befund: **`make test`/`make gates`** — in dieser Sitzung
  unabhängig ausgeführt (nicht nur den Bericht übernommen): beide grün, 0
  Befunde.
- geprüft, ohne Befund: **Mutationsproben real reproduziert** — beide vom
  Implementer berichteten Rot-Gegenproben laufen tatsächlich rot (Blast-
  Radius bei einer davon abweichend, siehe F-2); Arbeitsbaum nach beiden
  Reproduktionen wieder sauber.
- geprüft, ohne Befund: **Docker-only** — `make test`/`make gates` laufen
  über `docker run` mit gepinnten Digests; kein lokales Toolchain-Install
  im Diff. Die zusätzlich berichteten `gofmt -l .`/`go vet ./...`-Läufe
  hängen an keinem Make-Target (kein Gate wrappt sie) — eigene Reproduktion
  über denselben gepinnten Toolchain-Container bestätigt aber dieselben
  Ergebnisse (ein vorbestehender, von diesem Diff nicht berührter Treffer
  `internal/application/port/outbound/log_test.go` bei `gofmt -l`; `go vet`
  sauber) — kein Docker-only-Verstoß feststellbar, aber unverifizierbar
  über ein Gate (siehe Summary).
- geprüft, ohne Befund: **Suppression-Verbot** — kein
  `nolint`/`noqa`/`SuppressMessage`/`d-check:ignore`-Marker im Diff.
- geprüft, ohne Befund: **Kommentar-Klassen** — alle neuen/geänderten
  Kommentare (`decode.go`, `mapper.go`, `transaction.go`,
  `transaction_test.go`, `decode_test.go`) tragen Zusage- oder
  Rang-Zeiger-Charakter, präsentischer Zustand, keine Konjunktiv-Rede über
  verworfene Alternativen, kein abgebrochener Satz.
- geprüft, ohne Befund: **Traceability der Commit-Message** — Betreff
  trägt `LH-FA-ADM-004`, keine `SPEC-*`/`ARC-*`-Kennung; `make
  commit-traceability` lief in dieser Sitzung gegen `HEAD~5..HEAD` grün.
- geprüft, ohne Befund: **Hard Rule 3.3 (git mv + Inhalt)** — kein `git mv`
  in diesem Commit (Slice bleibt in `in-progress/`, kein
  Lifecycle-Übergang); reine Inhaltsänderung ohne Move ist hier korrekt.
- geprüft, ohne Befund: **Doku-Update-Item (§2 DoD)** — `ChangeStorePort`
  und die in `spec/architecture.md`/`spec/pflichtenheft.md` dokumentierten
  öffentlichen Verträge sind unverändert; die Begründung "reine interne
  Signatur-Änderung, Item entfällt" ist durch eigene Grep-Prüfung gedeckt.
- geprüft, ohne Befund: **Zwei-Quellen-Drift Plan/Verzeichnis** — DoD-Häkchen
  im Slice-Plan sind vollständig unchecked, Datei liegt in `in-progress/`;
  kein Widerspruch zwischen Statustext und Verzeichnis-Position.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Kein HIGH-Finding.** Insbesondere: **kein Rollen-Widerspruch** zur
Plan-Nachzug-Bewertung — ich teile die Einstufung des Implementers ("vier
Test-Call-Sites sind bounded/mechanisch, kein Rückführungs-Trigger") nach
eigener Prüfung der vier Stellen und der §4-Schwelle. Der einzige
inhaltliche Dissens zum Implementer-Bericht betrifft F-1 (technische
Begründung, nicht die Bewertung von Umfang/Rückführung) und F-2
(Berichts-Genauigkeit einer Mutationsprobe) — beide sind isolierte
MEDIUM-/LOW-Findings ohne Rollen-Widerspruch (der Implementer hat diesen
Befunden noch nicht widersprochen, und beide sind ohne neue Entscheidung
korrigierbar). Der Konflikt-Pfad aus Modul 8 (Sequenz mit
Übergabe-Artefakten über den Architect) ist damit **nicht** ausgelöst.

**Finding-Klassen dieses Laufs:** Technische Begründung einer
Risiko-Closure benennt falsche Ursache für korrektes Verhalten ·
Mutationstest-Bericht unterschätzt Blast-Radius

## Verdikt

**Merge-blockierend:** nein — kein offenes HIGH-Finding; F-1 (MEDIUM)
betrifft eine noch nicht geschriebene §7-Formulierung und ist vor der
Closure-Notiz korrigierbar, F-2 (LOW) ist eine Berichtskorrektur ohne
Code-Auswirkung.

**Übergabe:** F-1 geht an den Implementer/Planner zur Korrektur der
§6-Risiko-1-Begründung vor dem Schreiben der §7-Closure-Notiz — die
tragende Aussage ist nicht "UTC-verankert", sondern "`.UnixNano()`/`.Equal()`
sind Location-unabhängig, unabhängig von der (hier: `Local`) Location, die
`pglogrepl.pgTimeToTime` zurückgibt". F-2 ist eine redaktionelle Korrektur
des Commit-Berichts ohne Code-Änderung, kein Blocker. Beide Findings gehen
zusätzlich als Finding-Klassen in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler. Dieser Report selbst ist ein **Lauf-Beleg**
und wird über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).

**Was für die Closure noch fehlt** (nur benannt, nicht bewertet — Planner-
Arbeit): §7 Closure-Notiz ist noch nicht geschrieben (alle
Platzhalter-`<…>` stehen noch); beide §6-Risiken tragen noch keinen der
drei Ausgänge (eingetreten/entfallen/weiter offen) — Risiko 2 (Call-Sites)
legt nach dieser Review-Sitzung "eingetreten, im Rahmen absorbiert" nahe,
Risiko 1 (Zeitzonen) hängt an der F-1-Korrektur; das
Beobachtungs-Register (`../observations/`) ist noch nicht fortgeschrieben
(weder neuer Eintrag noch "keine Beobachtung angefallen"-Vermerk); alle
DoD-Häkchen in §2 sind noch unchecked; die Datei liegt weiterhin in
`in-progress/` (kein `git mv` nach `done/`); die drei Paarungen sind laut
Slice-Plan korrekt der `welle-5`-Closure zugeordnet, nicht dieser
Slice-Closure.

---

**Gate-Beleg:** `make test` und `make gates` in dieser Review-Sitzung real
ausgeführt (beide grün, `make gates` mit 0 Befunden über
`baseline-verify`/`docs-check`/`commit-traceability`/`a-check`); zusätzlich
`gofmt -l .` und `go vet ./...` über den gepinnten Toolchain-Container
eigenständig reproduziert (ein vorbestehender, unberührter `gofmt`-Treffer,
`go vet` sauber). Beide vom Implementer berichteten Mutationsproben
(Mapper-Zeitstempel auf `0`, Domain-Zuweisung entfernt) in dieser Sitzung
selbst nachgestellt und wieder zurückgesetzt (`git status` danach jeweils
sauber). Arbeitsverzeichnis am Ende dieser Review-Sitzung unverändert
gegenüber Commit `3fa7e5d`.
