# Review-Report: slice-018 — 2026-09-12

**Review-Art:** Code — geprüft gegen Slice-Plan
(`docs/plan/planning/in-progress/slice-018-commit-zeitstempel-store-adapter.md`)
und [`ADR-0040`](../plan/adr/0040-clockport.md) (Maintainability).

**Gegenstand:** Commit `f35a94d`
(`feat(postgresstorage): committed_at trägt den realen Quell-Commit-Zeitpunkt
(LH-FA-ADM-004)`) —
`internal/adapters/driven/postgresstorage/mapper/{mapper.go,mapper_test.go}`,
`internal/adapters/driven/postgresstorage/queries/queries.go`,
`internal/adapters/driven/postgresstorage/{store.go,store_test.go}`,
`tools/schema/schema.yaml` (Kommentar).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-018-commit-zeitstempel-store-adapter.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8
  Sub-Area-Prüfung)
- [`ADR-0040`](../plan/adr/0040-clockport.md) (ClockPort — `internal/domain`
  und `internal/application/usecase/*` importieren `time` nicht; Fitness
  Function laut ADR selbst "Gate geplant", aktuell nicht in `.a-check.yml`
  verdrahtet)
- [`ADR-0011`](../plan/adr/0011-persist-before-ack.md) (Persist-before-ACK,
  Idempotenz-Semantik für den §6-Risiko-2-Kontext)
- [`LH-FA-ADM-004`](../../spec/lastenheft.md) (messbarer CDC-Abstand)
- `internal/domain/model/transaction.go` (`ChangeTransaction.commitPosition`
  / `.committed` / `.sourceCommittedAt`, `SourceCommittedAt()`,
  `CommitPosition()`) und `internal/domain/model/timepoint.go`
  (`TimePoint.UnixNanos`, `NewTimePoint`)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.3 Move/Inhalt-Trennung,
  §3.7 Kommentar-Klassen) · `harness/conventions.md` (MR-000, genau eine
  Sub-Area `PGC`, Greenfield)
- Commit-Traceability (`AGENTS.md` §5, `ADR-0045`)
- Beobachtungs-Register `docs/plan/planning/observations/` (elf
  BEO-PGC-Verzeichnisse, insbesondere `cdc-capture-lag-real`)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-017.md`
  (F-1 Zeitzonen-Begründung, F-2 Mutationsproben-Bericht) — beide betreffen
  Decoder/Mapper/Domäne, nicht den Store-Adapter dieses Slice

---

## Prüfungen im Detail

**1. Diff gegen Plan/[`ADR-0040`](../plan/adr/0040-clockport.md).** `internal/domain/model/` und
`internal/application/usecase/*` importieren weiterhin kein `"time"`
(eigener `grep -rn '"time"'` über beide Verzeichnisse: keine Treffer). Die
Konvertierung `time.Unix(0, sourceCommittedAt.UnixNanos).UTC()` steht
ausschließlich in `mapper.go` (`internal/adapters/driven/postgresstorage/mapper`,
Driven-Adapter-Schicht) — nicht in der Domäne und nicht im
Driving-seitigen Mapper (`internal/adapters/driving/replication/mapper`,
slice-017). `TimePoint` (`internal/domain/model/timepoint.go:10`) trägt das
Feld exakt als `UnixNanos int64` — die Implementer-Angabe ist korrekt, eigene
Prüfung des Feldnamens bestätigt es. `time.Unix(0, ns)` ist die
stdlib-korrekte Umkehrung, da `TimePoint` selbst Unix-Nanosekunden trägt
(`NewTimePoint(unixNanos int64) TimePoint`); `.UTC()` fixiert die Location
für die weitere Ausgabe (Spalte `timezone: true` in `schema.yaml`), ändert
den absoluten Zeitinstant nicht.

**2. Verworfenes zweites Rückgabe-Bool von `SourceCommittedAt()`.**
Eigene Prüfung der Domäne (`transaction.go:40-55`): `CommitPosition()` und
`SourceCommittedAt()` geben beide denselben Struct-Feldwert `t.committed`
als zweites Ergebnis zurück — es gibt kein zweites, unabhängiges
Committed-Flag, das auseinanderlaufen könnte. `store.go:82-92`
(`PersistTransaction`) liest `committed` über `CommitPosition()` und kehrt
bei `!committed` sofort zurück, **bevor** `NewTransactionRow` (und damit
`SourceCommittedAt()`) überhaupt aufgerufen wird. Zum Zeitpunkt des Aufrufs
ist `t.committed` also bereits als `true` bestätigt — der verworfene zweite
Rückgabewert kann an dieser Aufrufstelle nicht von `true` abweichen. Kein
Lücken-Pfad gefunden; die Implementer-Begründung ("analog zu
`CommitPosition()`") trägt.

**3. Realer Test `TestPersistCarriesSourceCommittedAtNotPersistenceTime`.**
Existiert (`store_test.go`, neu). Eigener Lauf: `make test-store` in dieser
Sitzung selbst ausgeführt (nicht nur den Bericht übernommen) — reales
PostgreSQL im Testcontainer, Paket `postgresstorage` `ok` in 4.284s (klar
über der reinen SQL-Laufzeit, konsistent mit der eingebauten 1,5s-Sleep).
Der Test setzt `sourceCommittedAt := time.Now().Add(-2 * time.Hour)`,
committet die Transaktion domänenseitig damit, schläft 1,5s, persistiert,
liest `committed_at` zurück und prüft zwei Bedingungen: (a) Abstand zu
`sourceCommittedAt` < 1ms, (b) `committed_at` liegt mehr als eine Sekunde
vor dem Persistenz-Aufruf. Unter der alten DEFAULT-Semantik (`committed_at`
ungeschrieben, Spalte zieht `current_timestamp`) hätte (a) einen Abstand von
~2 Stunden ergeben und (b) wäre verletzt gewesen — der Test unterscheidet
also tatsächlich zwischen den beiden Verhalten, nicht nur zufällig über eine
Toleranzbreite. Eine genauere Prüfung der eingebauten Verzögerung zeigt
allerdings: Der 2-Stunden-Versatz allein trägt bereits beide Assertions;
die 1,5s-Sleep trägt zu keiner der beiden Bedingungen bei (siehe F-1).

**4. Idempotenz-Interaktion (`ON CONFLICT ... DO NOTHING`).** Die
bestehende `TestPersistTransactionIsIdempotent` (`store_test.go:194`)
persistiert **dieselbe** `tx`-Instanz zweimal — beide Aufrufe tragen
denselben `sourceCommittedAt`-Wert, die Zeile prüft also nicht das im
Slice-Plan §6 benannte Szenario (Retry mit *abweichendem* Zeitstempel).
Das ist konsistent mit der Selbstauskunft des Implementers und mit §6 des
Slice-Plans, das dieses Risiko ausdrücklich als "ohne expliziten Test
bislang unbelegt. Wird bei Closure bewertet" führt. Sachlich verankert
(`ADR-0011`, Persist-before-ACK): ein Retry nach Crash vor ACK wiederholt
dieselbe WAL-Position und damit denselben Quell-Commit-Datensatz — der
Quell-Commit-Zeitpunkt wäre in der Praxis bei einem echten Retry identisch,
nicht abweichend; ein abweichender Wert setzte einen Decoder-/Domänenfehler
voraus, der außerhalb dieses Slice liegt. Die Einstufung als offenes,
unverankertes Risiko statt als zu bauender Test ist dadurch nachvollziehbar
begründbar; siehe F-2 für die Einordnung, was für die Closure fehlt.

**5. `tools/schema/plan.yaml`-Reset.** `git show f35a94d --stat` führt die
Datei nicht; `git log -3 -- tools/schema/plan.yaml` zeigt nur ältere
Commits vor slice-018. Eigener `make test-store`-Lauf in dieser Sitzung hat
die Datei tatsächlich als Nebenprodukt regeneriert (Ziel-DSN unterscheidet
sich je Testlauf); nach dem Lauf per `git checkout -- tools/schema/plan.yaml`
zurückgesetzt. Die Implementer-Aussage ist damit bestätigt: kein
Nebenprodukt-Leck im finalen Commit.

**6. Hard Rules.** Siehe Negativbefunde unten (Traceability, Docker-only,
Suppression, Kommentar-Klassen, git-mv-Trennung) — je einzeln geprüft, kein
Verstoß gefunden.

**7. `make gates` / `make test-store` unabhängig nachvollzogen.** Beide in
dieser Sitzung selbst ausgeführt: `make test-store` → alle Pakete `ok`
(inkl. `postgresstorage` und `postgresstorage/mapper`); `make gates` →
`baseline-verify` (v6.5.0, 54 Dateien), `docs-check`/d-check (188 Dateien, 0
Befunde, zweimal — Standardlauf und `commits`-Modul über `HEAD~5..HEAD`),
`commit-traceability` (5 Commits, Betreffs ohne Struktur-ID), `a-check`
(0 Befunde) — alle grün, deckt sich mit dem Implementer-Bericht.

**8. §8-Angaben des Slice-Plans stichprobenhaft geprüft.** Das
Beobachtungs-Register trägt tatsächlich elf `BEO-PGC/*`-Verzeichnisse
(eigenes `find`); `cdc-capture-lag-real/evidence/` trägt genau eine Datei
(`slice-013.md`), `state.md` bestätigt "Zähler (abgeleitet): 1×" — die
Plan-Aussage "1× — dieser Slice bringt den Wert real in die DB, liefert
aber noch keinen Register-Beleg" ist damit konsistent mit dem tatsächlichen
Register-Stand.

---

## Findings

### F-1 — Künstliche Testverzögerung trägt nicht zur Unterscheidungskraft des Tests bei

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/adapters/driven/postgresstorage/store_test.go`
  (`TestPersistCarriesSourceCommittedAtNotPersistenceTime`, `time.Sleep(1500
  * time.Millisecond)`-Zeile und ihr Kommentar)
- `befund`: Der Testkommentar begründet die 1,5s-Sleep damit, dass "die
  Instanzzeit beim Exec ... danach klar nach sourceCommittedAt" liege. Der
  Zeitstempel ist aber bereits beim `Commit()`-Aufruf auf 2 Stunden in der
  Vergangenheit gesetzt — dieser Versatz allein reicht für beide
  Assertions (Abstand-zu-`sourceCommittedAt` und
  Abstand-zu-`persistCallTime`) aus, unabhängig von der Sleep-Dauer. Die
  Sleep verlangsamt den Testlauf messbar (Teil der beobachteten 4,284s
  Paketlaufzeit), ohne die Trennschärfe des Tests zu erhöhen; der Kommentar
  suggeriert eine kausale Notwendigkeit, die nicht besteht.
- `verifizierbar`: ja — Sleep-Zeile entfernen und Test erneut laufen lassen
  (kein Gate prüft das automatisiert).
- `klasse`: Testkommentar behauptet Notwendigkeit einer induzierten
  Verzögerung, die die eigentliche Testbedingung nicht trägt.

### F-2 — Idempotenz-Risiko (§6-Risiko-2) bleibt ohne Ausgang und ohne Test — konsistent mit Plan, aber offen für die Closure

- `kategorie`: INFO
- `quelle`: Maintainability (Slice-Plan §6, Risiko 2)
- `pfad`:
  `docs/plan/planning/in-progress/slice-018-commit-zeitstempel-store-adapter.md`
  §6, zweiter Punkt
- `befund`: Das Risiko "`ON CONFLICT ... DO NOTHING` behält bei
  abweichendem Retry-Zeitstempel den zuerst geschriebenen Wert" ist im Plan
  benannt, aber weder mit Test belegt noch mit einem der drei
  Closure-Ausgänge (eingetreten/entfallen/weiter offen) versehen — zum
  Prüfzeitpunkt erwartungsgemäß, da die Datei noch in `in-progress/` liegt.
  Die bestehende `TestPersistTransactionIsIdempotent` deckt nur den
  Gleich-Zeitstempel-Fall ab. Kein Diff-Defekt, sondern ein Hinweis für die
  Closure-Notiz (§7).
- `verifizierbar`: nein (Planungs-/Closure-Frage, kein Gate-Gegenstand)
- `klasse`: Offenes Risiko ohne Ausgang zum Review-Zeitpunkt — erwartungsgemäß vor Closure.

---

## Negativbefunde

- geprüft, ohne Befund: **`ADR-0040`-Grenze (Domäne bleibt time-frei)** —
  kein `"time"`-Import in `internal/domain` oder
  `internal/application/usecase/*`; die einzige `time.Unix(...)`-Konvertierung
  steht im Store-Adapter-Mapper, korrekt platziert.
- geprüft, ohne Befund: **`TimePoint.UnixNanos`-Feldname** — eigene Prüfung
  von `internal/domain/model/timepoint.go:10` bestätigt die
  Implementer-Angabe exakt.
- geprüft, ohne Befund: **Verworfenes zweites Rückgabe-Bool von
  `SourceCommittedAt()`** — `CommitPosition()` und `SourceCommittedAt()`
  spiegeln denselben Feldwert `t.committed`; `PersistTransaction` prüft ihn
  vor dem Aufruf von `NewTransactionRow`. Kein Lückenpfad gefunden.
- geprüft, ohne Befund: **`TestPersistCarriesSourceCommittedAtNotPersistenceTime`
  real** — eigener `make test-store`-Lauf bestätigt Existenz, echte
  PostgreSQL-Persistenz und dass die zwei Assertions unter der alten
  DEFAULT-Semantik fehlgeschlagen wären (siehe F-1 für einen Nebenpunkt zur
  Sleep-Notwendigkeit, ohne die Aussagekraft des Tests selbst zu berühren).
- geprüft, ohne Befund: **`queries.go`/`store.go`-Verdrahtung** — vierter
  SQL-Parameter, vierter Exec-Parameter und `TransactionRow.CommittedAt`
  bilden eine durchgängige Kette ohne Bruch.
- geprüft, ohne Befund: **`tools/schema/plan.yaml`** — nicht Teil von
  `f35a94d`; eigener `make test-store`-Lauf hat die Datei als Nebenprodukt
  regeneriert (andere Ziel-DSN), danach zurückgesetzt.
- geprüft, ohne Befund: **`make test-store`/`make gates`** — in dieser
  Sitzung unabhängig ausgeführt (nicht nur den Bericht übernommen): beide
  grün, 0 Befunde.
- geprüft, ohne Befund: **Docker-only** — `make test-store`/`make gates`
  laufen über `docker run` mit gepinnten Digests; kein lokales
  Toolchain-Install im Diff.
- geprüft, ohne Befund: **Suppression-Verbot** — kein
  `nolint`/`noqa`/`SuppressMessage`/`d-check:ignore`-Marker im Diff.
- geprüft, ohne Befund: **Kommentar-Klassen** — alle neuen/geänderten
  Kommentare (`mapper.go` `TransactionRow`/`NewTransactionRow`,
  `queries.go` `InsertTransaction`, `schema.yaml` `committed_at`) sind
  präsentisch/indikativ, tragen `LH-FA-ADM-004`/`seit slice-018` als
  Herkunfts-Anker, keine Konjunktiv-Rede über eine verworfene Alternative,
  kein abgebrochener Satz.
- geprüft, ohne Befund: **Traceability der Commit-Message** — Betreff trägt
  `LH-FA-ADM-004`, keine `SPEC-*`/`ARC-*`-Kennung; `make
  commit-traceability` lief in dieser Sitzung gegen `HEAD~5..HEAD` grün.
- geprüft, ohne Befund: **Hard Rule 3.3 (git mv + Inhalt)** — kein `git mv`
  in diesem Commit (Slice bleibt in `in-progress/`), reine
  Inhaltsänderung ohne Move ist hier korrekt.
- geprüft, ohne Befund: **Zwei-Quellen-Drift Plan/Verzeichnis** — DoD-Häkchen
  im Slice-Plan sind vollständig unchecked, Datei liegt in `in-progress/`;
  kein Widerspruch zwischen Statustext und Verzeichnis-Position.
- geprüft, ohne Befund: **§8-Angaben des Slice-Plans** — Registerzahl
  (elf `BEO-PGC/*`) und `cdc-capture-lag-real`-Zählerstand (1×) eigenständig
  gegen den tatsächlichen Register-Bestand geprüft, beide korrekt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Kein HIGH-Finding, kein MEDIUM-Finding, kein Rollen-Widerspruch.** Beide
Findings sind isoliert (F-1 LOW: Testkommentar/Sleep-Redundanz ohne
Auswirkung auf die Testaussage selbst; F-2 INFO: erwartungsgemäß noch
offenes Risiko vor Closure). Der Implementer hat keinem Befund
widersprochen — der Konflikt-Pfad aus Modul 8 (Sequenz mit
Übergabe-Artefakten über den Architect) ist **nicht** ausgelöst.

**Finding-Klassen dieses Laufs:** Testkommentar behauptet Notwendigkeit
einer induzierten Verzögerung, die die eigentliche Testbedingung nicht
trägt · Offenes Risiko ohne Ausgang zum Review-Zeitpunkt (erwartungsgemäß
vor Closure)

## Verdikt

**Merge-blockierend:** nein — kein offenes HIGH- oder MEDIUM-Finding; F-1
ist eine redaktionelle Testkommentar-/Sleep-Korrektur ohne
Korrektheitsauswirkung, F-2 ist der planmäßig noch fehlende Closure-Schritt
für ein bereits im Plan benanntes Risiko.

**Übergabe:** F-1 geht an den Implementer zur optionalen Bereinigung des
Sleep-Kommentars (bzw. Reduktion/Entfernung der Sleep-Dauer) vor der
nächsten Berührung dieser Datei — kein Blocker für diesen Slice. F-2 geht
an den Planner: Bei der Closure-Notiz (§7) braucht §6-Risiko-2 einen der
drei Ausgänge (voraussichtlich "weiter offen" mit Verweis auf die fehlende
Retry-mit-abweichendem-Zeitstempel-Testabdeckung, oder "entfallen" mit der
in Prüfung 4 genannten Begründung über die WAL-Determinismus-Eigenschaft
von `ADR-0011`-Retries — das Urteil, ob diese Begründung trägt, liegt beim
Planner/Architect, nicht bei diesem Review). Beide Findings gehen
zusätzlich als Finding-Klassen in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler. Dieser Report selbst ist ein **Lauf-Beleg**
und wird über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).

**Was für die Closure noch fehlt** (nur benannt, nicht bewertet — Planner-
Arbeit): §7 Closure-Notiz ist noch nicht geschrieben (alle
Platzhalter-`<…>` stehen noch); beide §6-Risiken tragen noch keinen der
drei Ausgänge (Risiko 1 — bestehende Tests implizit auf Persistenzzeit
vertrauend — legt nach dieser Review-Sitzung "entfallen"/"nicht
eingetreten" nahe, da `make test-store` vollständig grün lief; Risiko 2
siehe F-2); das Beobachtungs-Register ist noch nicht fortgeschrieben (der
Plan sieht dafür laut §8 explizit **keinen** neuen Beleg in diesem Slice
vor — das gehört als "keine Beobachtung angefallen" oder mit Bezug auf
`cdc-capture-lag-real` in §7 benannt); alle DoD-Häkchen in §2 sind noch
unchecked; die Datei liegt weiterhin in `in-progress/` (kein `git mv` nach
`done/`); die drei Paarungen sind laut Slice-Plan korrekt der
`welle-5`-Closure zugeordnet, nicht dieser Slice-Closure.

---

**Gate-Beleg:** `make test-store` und `make gates` in dieser
Review-Sitzung real ausgeführt (beide grün, `make gates` mit 0 Befunden
über `baseline-verify`/`docs-check`/`commit-traceability`/`a-check`). Der
neue Test `TestPersistCarriesSourceCommittedAtNotPersistenceTime` lief real
gegen die PostgreSQL-Testinstanz im Rahmen dieses Laufs mit (Paket
`postgresstorage` `ok`, 4,284s). `tools/schema/plan.yaml`, das der eigene
Testlauf als Nebenprodukt regenerierte, wurde danach zurückgesetzt
(`git status` am Ende dieser Review-Sitzung sauber gegenüber Commit
`f35a94d`).
