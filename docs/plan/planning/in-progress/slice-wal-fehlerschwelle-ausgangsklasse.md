# Slice wal-fehlerschwelle-ausgangsklasse: Fehlerschwelle des WAL-Rückstands — der Ausgang von `Run` trägt die Klasse `replication`, auch wenn der Abbruch den Stream mit einem Kontext-Fehler beendet

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er startet nach `slice-capture-leerlauf-quellbelege`
und `slice-leerlauf-phase-last-in-stuecken`
und geht `slice-start-vorlauf-grenze` voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5, Kante zu
`slice-transformationen-e2e-abhilfe` über `slice-start-vorlauf-grenze`).

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein
Datenverlust; der Abbruch ist sichtbar),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand sichtbar,
Schwellen),
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Folgepflicht:
oberhalb der Fehlerschwelle kontrollierter Abbruch mit `replication`-Klassifikation;
Entscheidung (a): nur die drei Stream-Ordnungs-Sentinels bleiben unabhängig vom
WAL-Rückstand hart abbrechend),
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (Fehlerklassen),
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) (Leerlauf-Bestätigung;
Kontext der Fehlerschwellen-Kette), Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
§2 bis §4.

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklassen; Zeile `replication`) und
[`SPEC-013`](../../../../spec/pflichtenheft.md) (Schwellen des WAL-Rückstands) —
gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent, Planner-Zug nach dem Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md).
**Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Erreicht der WAL-Rückstand die Fehlerschwelle, trägt der Fehlerzustand des
Feed-Containers die Klasse `replication` — unabhängig davon, wo der Stream im Moment
des Abbruchs steht: `mergeStreamAndWALFaultOutcome` (`internal/bootstrap/wiring.go`)
gibt bei gesetztem WAL-Fehler den WAL-Fehler zurück, wenn der Stream-Fehler die Folge
des Abbruchs ist, den die Schwellen-Prüfung selbst auslöst
([`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md), Verdikt §2).

**Kernaussage und ihre Herkunft** (**übernommen** aus Verdikt §1 und §5, nicht vom
Planner gemessen): bei gehaltener Persistierung endet der Feed-Container mit Ausgang 1
und der Klasse `storage` statt `replication` — gemessen vom Implementer von
`slice-capture-leerlauf-quellbelege` in fünf von fünf Läufen (drei Läufe der Phase in
einem Wegwerf-Aufbau, ein vollständiger `make test-integration`, ein Diagnose-Lauf der
Vorfassung; Anker: Plan dort §3 „Befund der Erprobung“, Ausgabezeile der Runner-Phase
„Fehlerschwelle beendet den Container“). Die Abbildung „Stand des Streams beim Abbruch →
Klasse“ (Verdikt §5) und die Ursache — der Kontext-Abbruch lässt `Capture` mit einem
Fehler der Klasse `storage` zurückkehren, und die Priorität „Stream-Fehler jeder Klasse
zuerst“ gibt ihn vor dem WAL-Fehler zurück — sind aus dem Quelltext **hergeleitet**, nicht
erprobt. Dass `errors.Is(err, context.Canceled)` an der echten Stelle trägt, ist
**hergeleitet** (die Meldung des Befunds endet auf „context canceled“, `Classify` wrappt
mit `%w: %w`) und nicht gelesen; der Slice belegt es zweifach (§2, Punkt 1 und 2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung an den Schwellen, an `classifyRunError`, an der Sentinel-Trennung oder an
  der Klassen-Tabelle.** `ADR-0049` wird eingelöst, nicht geändert; die Norm ist
  eindeutig, der Stand weicht ab (Verdikt §2, Konflikt-Pfad). Keine neue ADR.
- **Das Benutzerhandbuch.** Die vier Stellen sagen `replication` für die Fehlerschwelle
  zu und bleiben mit der Korrektur wahr (§3, Suchlauf: sechs Treffer gelesen); eine
  Text-Korrektur wäre die teurere Lösung gewesen (Verdikt §2, Begründung 3).
- **Ein Betreiber-Weg, die Persistierung anzuhalten.** Der Beleg nutzt den Aufbau
  der Runner-Phase; kein Produktpfad entsteht.
- **Die Zeile „`observeRelation`“ der Abbildung in Verdikt §5** (Abbruch im
  Schema-Store). Akzeptiertes Negativ des Verdikts: sie fällt unter die Regel, sobald
  ihre Kette `context.Canceled` trägt, und wird nicht gesondert erprobt.
- **Die Klassen- und Wiederholungsfrage bei `transient`.**
  `slice-capture-transient-wiederholung`;
  die Schwellen-Kette ist eine andere Ursache.
- **Der Start-Pfad.** `slice-start-vorlauf-grenze`
  ändert `Stream.Run`/`START_REPLICATION` und den Vorlauf; dieser Slice ändert die
  Rückgabe-Priorität nach der Rückkehr von `stream.Run`. Beide berühren `wiring.go` und
  den Runner an entgegengesetzten Enden (Verdikt §3, Reihenfolge).

## 2. Definition of Done

Jedes Kriterium trägt „Zu belegen durch:“; jede Aussage über eine Mutation ist eine
**Erwartung**, bis der Implementer sie gefahren hat (Stelle, Instanz, gesehene Farbe;
[`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).

- [x] **Die Regel steht im Code.** `mergeStreamAndWALFaultOutcome` folgt der
      Vertragstabelle des Verdikts (§2): ist ein WAL-Schwellen-Fehler gesetzt, ist der
      Ausgang von `Run` bei Stream-Ausgang `nil` der WAL-Fehler; bei einem
      Ordnungs-Sentinel (`mapper.ErrChangeWithoutBegin`, `…CommitWithoutBegin`,
      `…BeginWithoutCommit`) der Sentinel; bei einem Fehler, dessen Kette
      `context.Canceled` trägt, der WAL-Fehler (neu); bei jedem anderen Stream-Fehler der
      Stream-Fehler. Ohne gesetzten WAL-Fehler bleibt der Stream-Ausgang unverändert; ein
      echter Persistenzfehler ohne Abbruch-Folge bleibt `storage`. Tabellentest über die
      vier Zeilen, dazu ein Fall mit **echt gewrappter** Kette
      (`fmt.Errorf("%w: %w", outbound.ErrStorage, <Fehler mit context.Canceled>)`) und ein
      Fall „`storage` ohne Abbruch-Folge neben gesetztem WAL-Fehler bleibt `storage`“.
      *Belegt durch:* `make test` (Race-Detector) grün — `go test -race ./...` im
      gepinnten `TOOLCHAIN_RACE_IMAGE`, Exit 0, alle Pakete `ok`, `internal/bootstrap`
      1,252 s. Mutation **erprobt**: die alte Priorität (`streamErr != nil` liefert den
      Stream-Fehler ohne die `context.Canceled`-Prüfung) auf einer `git archive`-Kopie
      wiederhergestellt — der Fall
      `TestMergeStreamAndWALFaultOutcomeAbortDerivedStreamErrorYieldsFault/Fehlerkette_trägt_context.Canceled_(…)`
      färbte sich rot (`mergeStreamAndWALFaultOutcome(… context canceled) = …, wollen den
      WAL-Fehler`), der zweite Fall (echter Persistenzfehler ohne Abbruch-Folge) blieb
      grün — die Mutation trifft genau den beabsichtigten Fall. Die Prüfung
      `errors.Is(err, context.Canceled)` trägt die Kette in der Test-Instanz (echt
      gewrappte Kette `fmt.Errorf("%w: %w", outbound.ErrStorage, fmt.Errorf(...:
      %w, context.Canceled))`); ob sie an der realen Treiber-Kette
      (`postgresstorage`/`sqlexec`) ebenso trägt, belegt DoD 2 am komponierten Prozess.
      **Fixrunde (Review F-1, MEDIUM):** Der Reviewer band mit zwei eigenen Mutationen
      (Sentinel-Zweig auf einen von drei Sentinel-Werten verkürzt; die beiden
      `case`-Zeilen Sentinel ↔ `context.Canceled` vertauscht) eine reale Testlücke:
      keiner der bisherigen Fälle unterscheidet, ob eine Fehlerkette *sowohl* einen
      Ordnungs-Sentinel *als auch* `context.Canceled` trägt — beide Mutationen
      blieben grün. Heute strukturell unerreichbar (die drei Mapper-Sentinels sind
      unverkettete `errors.New`-Werte, `mapper.go:178,188,203`), aber eine Priorität,
      die `ADR-0049`(a) verlangt und kein Test band. Neuer Testfall
      `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled`
      (`walretention_internal_test.go`) konstruiert die Kette künstlich
      (`fmt.Errorf("%w: %w", <Sentinel>, context.Canceled)`, je einmal für alle drei
      Sentinels) und bindet die Priorität „Sentinel vor `context.Canceled`“ an einen
      Test. *Belegt durch:* `make test` (Race-Detector) grün, Exit 0. Mutation
      **erprobt** an beiden vom Reviewer genannten Stellen, je auf einer eigenen
      `tar`-Kopie des Arbeitsbaums gegen den gepinnten `golang:1.27`-Container:
      Sentinel-Zweig auf `ErrChangeWithoutBegin` verkürzt — zwei der drei Subtests
      (`…CommitWithoutBegin`, `…BeginWithoutCommit`) färbten sich rot, der dritte blieb
      grün, wie erwartet; die beiden `case`-Zeilen vertauscht — alle drei Subtests
      färbten sich rot. Beide Mutationen treffen genau den vom Reviewer benannten
      Fall.
      **Fixrunde (Verifikation V-1, MEDIUM):** Der Verifier band mit einer eigenen
      Mutation (früher Guard `if walErr == nil { return streamErr }` vollständig
      entfernt, Mutation M5) eine reale Testlücke: ein regulärer Prozess-Abbruch
      (SIGTERM) während `PersistTransaction`, dessen Kette `context.Canceled` trägt,
      **ohne** dass die WAL-Schwelle je erreicht wurde (`walErr == nil`), lief am
      kompletten `TestMergeStreamAndWALFaultOutcome*`-Lauf grün durch — kein Testfall
      unterschied „`context.Canceled`-Kette bei gesetztem WAL-Fault" von „dieselbe
      Kette bei leerem WAL-Fault". Der amtierende Code behandelt den Fall korrekt
      (der Guard gibt `streamErr` unverändert zurück, bevor der `switch` läuft); nur
      die Testbindung fehlte. Neuer Testfall
      `TestMergeStreamAndWALFaultOutcomeLeavesStreamErrorUnchangedWithoutWALFault`
      (`walretention_internal_test.go:396-414`) konstruiert dieselbe gewrappte Kette
      (`ErrStorage` + eingebetteter `context.Canceled`) mit **leerem** `walErr` und
      erwartet `streamErr` unverändert. *Belegt durch:* `make test` (Race-Detector)
      grün, Exit 0. Mutation **erprobt** an einer `git archive`-Kopie: der Guard
      entfernt (drei Zeilen) — genau ein Test färbte sich rot
      (`…LeavesStreamErrorUnchangedWithoutWALFault`, „= `<nil>`, wollen den
      Stream-Fehler unverändert"), alle übrigen blieben `PASS`; die zweite Verifikation
      hat dieselbe Mutation unabhängig reproduziert und exakt dieselbe Farbe gesehen
      (`verifikation-slice-wal-fehlerschwelle-ausgangsklasse-2` §2, „V-1 ist real
      geschlossen") und den Guard zusätzlich mit drei weiteren eigenen Mutationen
      (Bedingung invertiert, Rückgabewert geändert, Vergleichsvariable vertauscht;
      ebd. §4b N1–N4) gebunden — alle vier färbten sich wie erwartet rot.
- [x] **Die Runner-Phase trägt die Klasse als Zusage.** Die Phase „Fehlerschwelle beendet
      den Container“ in `make test-integration` prüft `cdc.process_heartbeat.error_class`
      auf `replication` (statt sie nur auszugeben); die Abdeckungs-Zeile in
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) (Erzeugnis des
      Runners) nennt die Klasse; die Ausgabezeile der Phase bleibt. *Belegt durch:* ein
      realer, grüner `make test-integration`-Lauf (`bash tools/harness/run-integration-tests.sh`,
      `set -euo pipefail` an der Spitze — die letzte Zeile `Lauf abgeschlossen — …`
      erreicht, kein `bf_fail` griff) mit der Ausgabezeile „…, cdc.process_heartbeat trug
      die Klasse replication; …“ und der Abdeckungs-Zeile in
      `docs/user/e2e-abdeckung.md` Zeile 68 („… meldet einen Fehlerzustand der Klasse
      replication über cdc.process_heartbeat …“, `Ort` `tools/harness/run-integration-tests.sh:3580`).
      Mutation **erprobt** (Eingabeseite: das Image, nicht der Runner) — `make image-mutation
      SRC=<git-archive-Kopie von HEAD, Commit fe9d0afa> TAG=wal-old-priority` (die alte,
      unbehobene Priorität ist exakt der vorige Commit), Bindung über eine
      Compose-Override-Datei im Scratchpad (`COMPOSE="docker compose -f compose.yaml -f
      <Override>"`, `image: pg-change-feed-mutation:wal-old-priority`) und derselbe (neue)
      Runner-Lauf: die Phase färbte sich rot exakt an der neuen Assertion —
      „`Fehlerschwelle beendet den Container — Fehlerklasse in cdc.process_heartbeat —
      erwartet 'replication', gelesen 'storage'`“ (`bf_fail`, Exit 1, `cleanup`-Trap fuhr die
      Umgebung herunter) — derselbe Aufbau (gehaltene Persistierung), der vor diesem Slice
      `storage` maß. Mutations-Image und Compose-Override danach entfernt
      (`make image-mutation-rm TAG=wal-old-priority`, `docker compose … down -v`).
- [x] **Kommentare und die benannte Grenze sind nachgezogen.** Der Kommentar an
      `mergeStreamAndWALFaultOutcome` (`wiring.go`: „jede Klasse“, „nur zum Zug, wenn der
      Stream-Lauf regulär endete“), der Kommentar von
      `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` und der Kommentar in
      `walretention_slotgrowth_internal_test.go` (die Kette „Fehler bei Stream-Ende“ trägt
      das Prozessende) tragen die Regel im Indikativ, je höchstens eine Kennung, ohne
      Konjunktiv über die frühere Fassung ([`AGENTS.md`](../../../../AGENTS.md) §3.7); die
      benannte Grenze zur Klasse in [`harness/README.md`](../../../../harness/README.md)
      §Sensors bei `make test-integration` entfällt, ihr Träger dort war dieser Slice.
      **Plan-Nachzug (über die drei genannten Stellen hinaus):** der Kommentar an
      `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`
      (`walretention_slotgrowth_internal_test.go` Zeilen 18–35) trug vier Kennungen
      (`ADR-0049`, `SPEC-013`, `LH-QA-REL-003`, `ADR-0120`) — mein Edit an Zeile 27
      überlappte diesen Block, `make kommentar-kennungen DIFF=fe9d0afa` meldete ihn als
      Kandidat; auf `ADR-0049` als einzigen Anker reduziert, die drei übrigen Kennungen
      durch beschreibende Prosa ohne Kennung ersetzt. Der Kommentar vor der Runner-Phase
      (`tools/harness/run-integration-tests.sh`, Zeilen 3488–3497) trug drei
      Träger-Übergaben aus `slice-capture-leerlauf-quellbelege` (Review F-4, Verifikation
      V-5, gemeldet in dessen §3): die Allaussage „der Rückstand erreicht die
      Fehlerschwelle nur, wenn der Slot nichts bestätigt“ ohne Anker (jetzt auf die
      Gegenseite der vorigen Phase bezogen statt als Allaussage), die Klasse ohne
      Rang-Zeiger (jetzt `` `ADR-0049` `` an der Stelle, wo der Fehlerzustand die Klasse
      trägt) und die Kopplung an die vorige Phase ohne benannte Symbole (jetzt
      `WAL_WARN_BYTES`/`WAL_ERROR_BYTES`, `bf_wal_hold`, `wal_feed_started` genannt) —
      alle drei in diesem Lauf behoben.
      *Belegt durch:* `make kommentar-kennungen DIFF=fe9d0afa` — Exit 0, kein Kandidat
      (nach dem Nachzug oben; ein Zwischenlauf vor dem Nachzug meldete den einen
      Kandidaten oben). `make fmt-check` — 260 Go-Dateien geprüft, alle formatiert; kein
      Chronik- oder Konjunktiv-Kandidat im diff-skopierten Lauf (Schritt 20; die drei
      Treffer auf „statt“ sind Mutationsbeschreibungen in Test-Godocs, zulässig; ein
      Treffer auf „slice-026“ ist eine Testfall-Provenienz-Zitierform in unverändertem
      Bestandscode, zulässig). Review folgt separat.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9): ein erster Lauf
      färbte `docs-check` rot (sechs `target-missing`-Befunde — der `next → in-progress`-Move
      dieses Slice hinterließ Links mit fest verdrahtetem `open/`-Pfad in `roadmap.md`
      (Zeilen 220, 310), im eigenen Plan (Zeilen 83, 85, 313) und in
      `slice-start-vorlauf-grenze.md` Zeile 289 — `BEO-PGC/slice-pfad-als-link-in-berichten`);
      alle sechs auf Kennungs-Zitat ohne Link umgestellt. Zweiter Lauf: `docs-check` 1352
      Dateien, 0 Befunde; `generated-sync` OK; `a-check` 0 Befunde; `coverage-gate` 85,30 %
      ≥ 80 %; `commit-traceability` OK (5 Commits, keine Struktur-ID im Betreff);
      `record-gates` schrieb den Arbeitsbaum-Hash — gegen `bash
      tools/harness/working-tree-hash.sh` danach gegengeprüft (byte-gleich), kein Commit/Move
      dazwischen (`git status --short` unverändert seit dem Lauf).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). Report:
      [`review-slice-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/review-slice-wal-fehlerschwelle-ausgangsklasse.md)
      — Verdikt 0 HIGH/1 MEDIUM/1 LOW/1 INFO, nicht merge-blockierend, keine zwingende
      Fixrunde (Checkbox-Nachzug durch den Reviewer selbst, Skill §DoD-Checkbox-Nachzug ohne
      Fixrunde).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13). **Fixrunden-Korrektur
      (Verifikation V-2, MEDIUM):** der ursprünglich eingetragene `diff`-Wert (56) traf nach
      der vorigen Fixrunde nicht mehr zu — der Verifier maß `ist=59`, weil
      `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled` den `diff`-Suchlauf
      nicht erneut gefahren hatte
      ([`verifikation-slice-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md)
      V-2). Diese Fixrunde hat den Wert erneut am jetzigen Stand (inklusive des neuen
      V-1-Testfalls) gemessen und auf `61` korrigiert (§3, Fixrunden-Nachmessung), von der
      zweiten Verifikation eigen bestätigt (Exit 0, „9 Zeilen stimmen", zusätzlich per
      unabhängiger `git grep`-Zählung). **Planner-Closure-Nachmessung:** Der Zug dieser
      Closure hat selbst zwei weitere Treffer auf `mergeStreamAndWALFaultOutcome` erzeugt
      (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/state.md`, das die bewegte Stelle
      als Text der Beobachtung nennt) — die Lehre der eigenen Closure-Notiz (§7,
      Steering-Loop-Eintrag) unmittelbar am eigenen Zug angewandt: `make
      suchlauf-nachmessen` vor dem Commit erneut gefahren, Abweichung `soll=61 ist=63`
      gesehen, den Wert auf `63` korrigiert (§3), danach erneut mit Exit 0 bestätigt.
- [x] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors (Zeile
      `make test-integration`: die Grenze entfällt, die Phase nennt die Klasse als
      Zusage); das Benutzerhandbuch bleibt unberührt (§1, real gegengeprüft: `git diff`
      trägt keine Änderung an `docs/user/benutzerhandbuch.md`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke) — §7.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben —
      `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` erhält eine neue
      `evidence/`-Datei (25., §7); `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`
      bleibt beim Ausgang „geplant → `slice-start-vorlauf-grenze`" (durch diesen Slice
      nicht zusätzlich setzbar, von der zweiten Verifikation bestätigt, §8); kein
      weiterer Anfall (§7 trägt beide).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen) — §6 trägt sie, §7 fasst sie zusammen.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure dieses Slice (§7) und zusätzlich der Closure der nächsten Welle
      (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten: ihre Closure liegt nach
      `slice-transformationen-betriebsdoku`, der nach diesem Slice startet; ein Slice
      ohne Welle wird von ihr mitgeprüft).

**Umfang:** S — Schätzung, nicht gemessen: eine Regel in einer kleinen Funktion, ein
Tabellentest, eine Assertion im Runner, drei Kommentare und ein Satz in
`harness/README.md` (Verdikt §3); ein Runner-Lauf trägt die Zeit von `make
test-integration` (Risiko §6).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (`mergeStreamAndWALFaultOutcome`, Kommentar; `walRetentionFault`-Kommentar liest der Implementer) | update | die Regel des Verdikts §2; der Kommentar trägt die Regel im Indikativ |
| `internal/bootstrap/walretention_internal_test.go` (`TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`, `…FallsBackToFaultOnRegularStreamEnd`) | update | Tabellentest über die vier Zeilen, gewrappter Fall, `storage`-ohne-Abbruch-Fall (`LH-QA-REL-003`, `ADR-0049`); Kommentar des ersten Tests nachgezogen |
| `internal/bootstrap/walretention_slotgrowth_internal_test.go` (Kommentar, Zeile 27) | update | die Kette „Fehler bei Stream-Ende“ trägt das Prozessende — Kommentar nachgezogen |
| `tools/harness/run-integration-tests.sh` (Phase „Fehlerschwelle beendet den Container“: `abdeckung_declare`-Text, Prüfung von `wal_stop_class`, Kommentar davor) | update | die Klasse `replication` als Zusage der Phase (Eingabeseiten-Mutation: Regel zurücknehmen); der Kommentar davor trägt drei Übergaben aus `slice-capture-leerlauf-quellbelege` (Review F-4, Verifikation V-5): die Allaussage „der Rückstand erreicht die Fehlerschwelle nur, wenn der Slot nichts bestätigt“ ohne Anker (die Phase belegt einen Fall), die Grenze der Klasse ohne Rang-Zeiger (`ADR-0049`, dieser Slice) und die Kopplung an die Phase „Leerlauf-Bestätigung“ davor (`WAL_*`, `bf_wal_hold`, `wal_feed_started`), im Kommentar nur für die Konfigurationsdatei genannt |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner; die Ort-Zeilen verschieben sich |
| `harness/README.md` §Sensors (Zeile `make test-integration`) | update | die benannte Grenze entfällt; ändern nur der genau benannte Satzteil (die Zeilen sind sehr lang) |
| `internal/bootstrap/walretention_slotgrowth_internal_test.go` (Kommentar an `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`, Zeilen 18–35) | update (Plan-Nachzug) | über den Plan hinaus: mein Edit an Zeile 27 überlappte einen Bestandsblock mit vier Kennungen (`ADR-0049`, `SPEC-013`, `LH-QA-REL-003`, `ADR-0120`); `make kommentar-kennungen DIFF=fe9d0afa` meldete ihn — auf `ADR-0049` als einzigen Anker reduziert (§3.7) |
| `docs/plan/planning/in-progress/roadmap.md` (Zeilen 220, 310), `docs/plan/planning/open/slice-start-vorlauf-grenze.md` (Zeile 289), diese Plan-Datei selbst (Zeilen 83, 85, 313) | update (Plan-Nachzug) | über den Plan hinaus: der `next → in-progress`-Move dieses Slice ließ sechs Markdown-Links mit fest verdrahtetem `open/`-Pfad zurück (`BEO-PGC/slice-pfad-als-link-in-berichten`); der erste `make gates`-Lauf färbte `docs-check` mit sechs `target-missing`-Befunden rot — alle sechs auf Kennungs-Zitat ohne Link umgestellt |

**Ansatz (Liste):**

- Die Regel ist eine Funktion mit vier Fällen (Verdikt §2), kein Umbau: `Run` ruft sie
  unverändert an einer Stelle (`wiring.go`, Rückgabe nach dem Herunterfahren der
  Goroutinen).
- Der Test des gewrappten Falls trägt die Form der echten Kette (`ErrStorage` und ein
  Fehler mit `context.Canceled`), nicht einen einfachen `context.Canceled`; der Fall
  „`storage` ohne Abbruch-Folge“ hält die Grenze nach der anderen Seite.
- Die Runner-Phase bleibt der Beleg am komponierten Prozess (`ADR-0030`): die Klasse
  erschien erst dort (`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`, Verdikt §6).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Klasse des Ausgangs bei
Erreichen der Fehlerschwelle“). Suchraum: der ganze Baum ohne die drei Ausnahmen von
[`AGENTS.md`](../../../../AGENTS.md) §3.13 (`docs/reviews/**`, Records unter `done/`,
`.harness/baseline/**`), keine weitere Einschränkung; die Plan-Datei schließt das Werkzeug
aus. Stand ist `7305b578` (der Commit vor der Anlage dieses Slice, vom Planner am
2026-09-27 gemessen; ein Stand ist eine Commit-Kennung, nie `HEAD`); der Implementer misst
am Parent seiner Arbeit neu und trägt die `diff`-Zeilen ein (`make suchlauf-nachmessen
PLAN=docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md`). Die Zahlen
des Standes `7305b578` sind mit `git grep -n` gemessen, nicht übernommen. **Implementer-
Nachmessung am eigenen Parent `fe9d0afa`** (Commit unmittelbar vor der ersten
Produktionsänderung dieses Laufs, nach dem `next → in-progress`-Move): zwischen `7305b578`
und `fe9d0afa` liegen acht fremde Slice-Läufe (`slice-harness-guard-blocked-python`,
`slice-leerlauf-phase-last-in-stuecken`, `slice-harness-mutationsbild-und-verweigerte-aktion`,
u. a.); `slice-capture-leerlauf-quellbelege` wanderte dabei nach `done/` (Ausschluss aus dem
Suchraum) und `harness/README.md`/`welle-transformationen.md`/Roadmap erhielten Träger-
Nachzüge — die Zahlen bewegen sich dadurch, ohne dass dieser Slice etwas geändert hätte
(Gefundenes: die Bewegung ist real und dokumentiert, nicht mein Zutun). Die Werte bei
`7305b578` bleiben unverändert wahr (Commit-Stand ist unveränderlich); `fe9d0afa` ist die für
diesen Lauf richtige Vergleichsbasis.**

Symbolnamen der bewegten Stelle:

```suchlauf
7305b578 58 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
fe9d0afa 51 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':(exclude,glob)**/slice-wal-fehlerschwelle-ausgangsklasse.md'
diff 63 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':(exclude,glob)**/slice-wal-fehlerschwelle-ausgangsklasse.md'
```

**Fixrunden-Nachmessung des `diff`-Werts (Verifikation V-2, dann eigener V-1-Testfall):**
der Verifier maß am Stand `8823f818` `ist=59` statt der zu diesem Zeitpunkt im Plan
eingetragenen `56` — Ursache: die vorige Fixrunde (`a7d27ddc`) hatte
`TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled` ergänzt (drei zusätzliche
lauffähige Treffer auf `mergeStreamAndWALFaultOutcome`), ohne den `diff`-Suchlauf danach
erneut zu fahren
([`verifikation-slice-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md)
V-2). Diese Fixrunde ergänzt zusätzlich
`TestMergeStreamAndWALFaultOutcomeLeavesStreamErrorUnchangedWithoutWALFault` (V-1 desselben
Reports) — zwei weitere Treffer (Aufruf, `t.Fatalf`-Meldung) heben den Wert auf `61`, von der
zweiten Verifikation bestätigt. Beide Bewegungen sind Nichtgefundenes im Sinne von
[`AGENTS.md`](../../../../AGENTS.md) §3.13: kein neuer Träger außerhalb der Tabellentestdatei
selbst ist entstanden, die Zahl wächst ausschließlich mit dem Testfall-Bestand derselben
Datei, die die Regel bindet.

**Planner-Closure-Nachmessung (dritte Bewegung, dieselbe Ursachen-Klasse, diesmal außerhalb
von Code):** Diese Closure legt die Beleg-Datei
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-wal-fehlerschwelle-ausgangsklasse.md`
an und schreibt die zugehörige `state.md` fort — beide nennen die bewegte Stelle
(`mergeStreamAndWALFaultOutcome`) als Text der Beobachtung selbst, zwei zusätzliche
Treffer in `state.md`. Vor dem Commit erneut mit `make suchlauf-nachmessen PLAN=…`
gemessen: Abweichung `soll=61 ist=63`, auf `63` korrigiert (oben), danach erneut mit
Exit 0 bestätigt (9 Zeilen stimmen). Das ist dieselbe Ursachen-Klasse wie die beiden
vorigen Bewegungen — ein neuer Text, der die bewegte Stelle beim Namen nennt —, hier
aber im Register statt im Code, und noch **innerhalb** der Closure gefunden statt von
einer nachfolgenden Rolle.

Beschreibung und Zählwort der Klassen-Aussage samt Hedge (die Klasse des Ausgangs, die
Priorität, das reguläre Stream-Ende):

```suchlauf
7305b578 24 -n -E 'Ausgangs-Klasse|Klasse des Ausgangs|mit der Klasse .replication.|beendet sich der Feed-Container|nur zum Zug|Fehler bei Stream-Ende|regulär endete|jede Klasse|Stream-Fehler jeder Klasse' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
fe9d0afa 20 -n -E 'Ausgangs-Klasse|Klasse des Ausgangs|mit der Klasse .replication.|beendet sich der Feed-Container|nur zum Zug|Fehler bei Stream-Ende|regulär endete|jede Klasse|Stream-Fehler jeder Klasse' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':(exclude,glob)**/slice-wal-fehlerschwelle-ausgangsklasse.md'
diff 18 -n -E 'Ausgangs-Klasse|Klasse des Ausgangs|mit der Klasse .replication.|beendet sich der Feed-Container|nur zum Zug|Fehler bei Stream-Ende|regulär endete|jede Klasse|Stream-Fehler jeder Klasse' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':(exclude,glob)**/slice-wal-fehlerschwelle-ausgangsklasse.md'
```

Die Differenz `fe9d0afa` (20) → `diff` (18): die zwei fremden Kommentar-Klauseln „jede Klasse“
und „nur zum Zug“ (`wiring.go`) sowie die Phrase „Fehler bei Stream-Ende“
(`walretention_slotgrowth_internal_test.go`) sind mit der Regel nachgezogen und tragen die
alten Formulierungen nicht mehr (DoD 3, Nichtgefundenes: kein neuer Träger dieser drei
Phrasen ist entstanden — die Differenz ist ausschließlich Abbau).

Die Klassen-Zusagen des Benutzerhandbuchs (der Träger, der die Klasse den Betreibern
nennt):

```suchlauf
7305b578 6 -n -E 'Feed-Container mit der Klasse .replication.|klassifiziert den Lauf als .replication.|Transport-/Verbindungsstörung' -- docs/user/benutzerhandbuch.md
fe9d0afa 6 -n -E 'Feed-Container mit der Klasse .replication.|klassifiziert den Lauf als .replication.|Transport-/Verbindungsstörung' -- docs/user/benutzerhandbuch.md
diff 6 -n -E 'Feed-Container mit der Klasse .replication.|klassifiziert den Lauf als .replication.|Transport-/Verbindungsstörung' -- docs/user/benutzerhandbuch.md
```

Unverändert bei 6 (Nichtgefundenes: das Benutzerhandbuch ist in diesem Lauf unberührt
geblieben, wie §1 zusagt — alle vier Stellen bleiben mit der Code-Korrektur wahr).

| Träger | Befund (Stand `7305b578`, vom Planner gelesen) | Behandlung |
|---|---|---|
| Kommentar an `mergeStreamAndWALFaultOutcome` (`wiring.go`, Zeilen 1168–1177) | trägt „jede Klasse“ und „nur zum Zug, wenn der Stream-Lauf regulär endete“; beide gelten nach der Regel nicht mehr | Implementer zieht nach (DoD 3) |
| Kommentar von `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` und `…FallsBackToFaultOnRegularStreamEnd` (`walretention_internal_test.go`) | der erste beschreibt „eine Stream-Ordnungs-Verletzung erreicht `Run`s Rückgabewert“ als einzige Priorität-Zusage, der zweite „erst ein regulärer Stream-Abschluss … lässt den WAL-Fehler durch“; beides wird mit der Regel enger/falsch | Implementer zieht nach (DoD 3) |
| Kommentar in `walretention_slotgrowth_internal_test.go` (Zeile 27) und die Assertion in Zeile 141 („bei regulärem Stream-Ende“) | nennt „Fehler bei Stream-Ende“ als Träger des Prozessendes; die Assertion prüft `merge(nil, fault)` und bleibt wahr | Kommentar nachziehen; Assertion lesen |
| `harness/README.md` §Sensors, Zeile `make test-integration` | trägt den Satz „**Benannte Grenze:** … der Befund liegt beim Architect“ | entfällt mit DoD 2 und 3; bis dahin nennt der Satz als Träger dieses Slice (Nachzug des Planners, Stand dieser Anlage) |
| Beschreibung der Phase im Runner (`abdeckung_declare`, Kommentar, Ausgabezeile) und `docs/user/e2e-abdeckung.md` Zeile 68 | die Ausgabezeile nennt die Klasse (`$wal_stop_class`), die Zusage nicht; die Abdeckungs-Zeile nennt sie nicht | Runner-Phase und Erzeugnis ziehen mit (DoD 2) |
| `docs/user/benutzerhandbuch.md` (Zeilen 466, 866, 871, 1526, 1614, 1833) | vom Planner gelesen: vier Stellen (466, 866–868, 1526, 1614–1615) sagen `replication` für die Fehlerschwelle bzw. die Unterart Transport-/Verbindungsstörung zu; 871 und 1833 beschreiben die Reichweite der Schwellen und die Versionshistorie. **Alle bleiben mit der Korrektur wahr; kein Träger nennt `storage` für diese Kette** (Muster `Klasse .storage.` im Handbuch: zwei Treffer, Zeilen 578/581, betreffen den Backfill-Run im Snapshot-Fenster, einen anderen Gegenstand) | keine Änderung; der Implementer misst neu |
| `spec/pflichtenheft.md` (Klassen-Tabelle, Zeile `replication`) und `ADR-0049` | die Norm, gelesen (Verdikt L1 bis L3); ihr Wortlaut ist der Gegenstand der Korrektur, kein nachzuziehender Träger | unberührt |
| **Fremde Träger:** `slice-capture-leerlauf-quellbelege` §2 und §3 (DoD 2, Meldung der Träger); das Register `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` und `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (`state.md`) | der Plan dort nennt dieses Slice als Adresse (Nachzug des Planners); die Registereinträge tragen den Fall (Verdikt §6) | Register: **gemeldet, nicht mitgeändert**; Frist: die Closure dieses Slice, der Planner der Closure setzt die Ausgänge |
| Beschreibung in `docs/plan/planning/` (Plan von `slice-capture-leerlauf-quellbelege`, Welle-Plan, Roadmap) | die Aussage zur Klasse `storage` bei gehaltener Persistierung steht im Plan des Quellbeleg-Slice (§3 „Befund der Erprobung“) und beschreibt die Messung jenes Slice | bleibt dort stehen; die Träger-Meldung und DoD 2 dort nennen diesen Slice als Adresse |

## 4. Trigger

**Start** (`next` → `in-progress`): eine **Vorab**-Bedingung, kein Nachweis nach der
Umsetzung (`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen, 2×):
`slice-capture-leerlauf-quellbelege` liegt in `done/` (beide Slices ändern die
Runner-Phase „Fehlerschwelle beendet den Container“ in
`tools/harness/run-integration-tests.sh`; die Phase ist die Falsifikation dieses Slice;
WIP-Limit 1), `slice-leerlauf-phase-last-in-stuecken`
liegt in `done/` (beide Slices ändern denselben Runner an verschiedenen Stellen — jener die Last
der Phase „Leerlauf-Bestätigung“, dieser die Phase „Fehlerschwelle beendet den Container“ —, und ein
Rot in der Leerlauf-Phase lässt diese Phase ungelaufen;
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§3), und kein weiterer Slice liegt in `in-progress/`. Der Slice muss `done`
sein, **bevor** `slice-start-vorlauf-grenze` startet
(Start-Trigger dort): beide ändern `wiring.go` und den Runner; der kleinere Slice
zuerst verstellt die Prüfspur des größeren nicht (Verdikt §3). Der Übergangs-Commit
`next` → `in-progress` nennt
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md).

**Kante zu `slice-leerlauf-phase-last-in-stuecken` (beauftragt).** Die Phase
„Leerlauf-Bestätigung“ vor der Runner-Phase dieses Slice war im ersten
`e2e.yml`-Lauf nach dem Push von `slice-capture-leerlauf-quellbelege` einmal rot (Lauf
36287009221, Leg PostgreSQL 18); der Architect hat die Ursache im Testaufbau bestimmt (der
Stoß einer einzelnen Anweisung gegen den Prüf-Takt) und den Slice beauftragt
([`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§3; Register `BEO-PGC/test-integration-retention-timing-flake`). Er geht diesem Slice
**voraus**: ein Rot dort lässt die Phasen dahinter ungelaufen, auch die Mutation dieses Slice
an der Phase „Fehlerschwelle beendet den Container“.

**Verifier-Hinweis** (Verdikt §4): ein Rot mit der Signatur „Fehlerklasse `replication` …
WAL-Rückstand … über Fehlerschwelle“ in der Phase „Leerlauf-Bestätigung“ ist, solange
`slice-leerlauf-phase-last-in-stuecken` nicht in `done/` liegt, weder Beleg noch
Widerlegung dieses Slice — der Verifier wiederholt den Lauf (`gh run rerun <Lauf> --failed`)
und nennt Lauf, Versuchsnummer und Job-Kennungen beider Versuche; nach diesem Slice ist
dasselbe Rot ein Befund und ein Architect-Zug. Ein Rot anderer Signatur oder in einer anderen
Phase gilt nicht als dieser Fall.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht zu erwarten (Umfang S);
  sollte die Runner-Phase einen zweiten, unabhängigen Befund zeigen, ist dieser der
  abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls die Regel die Klasse in der Runner-Phase
  nicht erreicht, obwohl der Unit-Test grün ist — die Kette trägt `context.Canceled`
  dann an der echten Stelle nicht, und der zulässige Ersatz (`streamCtx.Err() != nil`)
  trägt ebenfalls nicht: eine Architect-Frage, kein stiller Umbau.

## 5. Closure-Trigger

DoD vollständig (die Regel im Code, die Klasse `replication` als Zusage der Runner-Phase
am komponierten Prozess belegt, Kommentare und die benannte Grenze nachgezogen) + `make
gates` grün + ein realer, grüner `make test-integration`-Lauf + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Fall, den die Regel abbildet, bleibt unerreichbar: `errors.Is(err,
  context.Canceled)` trägt an der echten Stelle nicht.** *Hergeleitet, nicht erprobt*
  (Verdikt §2: `Classify` wrappt mit `%w: %w`, die Meldung endet auf „context canceled“;
  nicht gelesen). *Erwartet, zu belegen durch:* der Tabellentest mit gewrappter Kette
  **und** die Runner-Phase mit der Klasse `replication`; der zulässige Ersatz ist
  `streamCtx.Err() != nil` bei lebendem Prozess-Kontext. **Ausgang: nicht eingetreten** —
  der reale, grüne `make test-integration`-Lauf des Reviewers zeigt die Klasse
  `replication` am komponierten Prozess (nicht nur im Unit-Test); die Kette trägt
  `context.Canceled` an der echten `postgresstorage`/`sqlexec`-Stelle — vom Verdikt-Status
  „hergeleitet" auf „erprobt" gehoben (Review-Report §„Eigene Messungen"; zweite
  Verifikation §7 Zeile 1).
- **Die Regel verdeckt einen echten Persistenzfehler.** Ein Fehler der Klasse `storage`
  ohne Abbruch-Folge neben einem gesetzten WAL-Fehler würde nach der Regel weiter den
  Stream-Fehler zurückgeben; ein Fehler, dessen Kette `context.Canceled` trägt **und**
  der eine echte Ursache hat, wird zum WAL-Fehler — die Ursache steht im Log neben der
  Abbruch-Zeile (Verdikt §2). *Erwartet, zu belegen durch:* der Fall „`storage` ohne
  Abbruch-Folge bleibt `storage`“ im Tabellentest; ohne gesetzten WAL-Fehler ändert sich
  nichts. **Ausgang: nicht eingetreten** — der Fall ist durch den Tabellentest und die
  Reviewer-Mutation M7 (`default: return nil` statt `return streamErr`, korrekt rot)
  gebunden; keine Verdeckung eines echten Persistenzfehlers beobachtet.
- **Die Mutation färbt nichts rot** (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  verkörpert): der Test der Regel und die Runner-Phase binden sich an die Eingabeseite
  (der Stream-Fehler der Klasse `storage` neben dem gesetzten WAL-Fehler, die Klasse aus
  `cdc.process_heartbeat`), nicht nur an die Ausgabe. *Erwartet, zu belegen durch:* je eine
  Mutation, deren Farbe der Implementer **gesehen** und im Bericht genannt hat (Stelle,
  Instanz, Farbe). **Ausgang: überwiegend nicht eingetreten, mit einer real gefundenen und
  geschlossenen Lücke** — Implementer und Reviewer zusammen 6 Mutationen (alle wie
  erwartet rot); die erste Verifikation fand mit einer eigenen, im Auftrag nicht genannten
  Mutation (M5, früher Guard `if walErr == nil { return streamErr }` entfernt) eine reale
  Lücke (V-1, MEDIUM: der komplette Testlauf blieb grün), die Fixrunde band sie mit einem
  neuen Testfall, die zweite Verifikation reproduzierte die Mutation unabhängig (genau ein
  Test rot) und ergänzte vier weitere eigene Mutationen (N1–N4, alle rot) — insgesamt 14
  unabhängige Mutationen über beide Verifikationsrunden, keine färbt mehr unerwartet grün.
- **Der Beleg ist zeitabhängig und verlängert `make test-integration` nicht**: die
  Phase besteht bereits; ihre Laufzeit ist unverändert bis auf die Klassen-Prüfung
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert). *Erwartet, zu belegen
  durch:* die gedruckte Laufzeit des Laufs ([`AGENTS.md`](../../../../AGENTS.md) §3.12
  Instanz A). **Ausgang: plausibel, unbelegt** — weder Implementer noch Reviewer noch
  Verifier nennen eine Vorher-/Nachher-Laufzeit der Phase explizit (beide Verifikationen
  bestätigen das ausdrücklich, zweite Runde §1 „nicht erneut gefahren"); die Änderung ist
  eine reine Assertion auf einen bereits gelesenen Wert (`$wal_stop_class`), kein neuer
  Wartezyklus — plausibel ohne Laufzeitwirkung, aber keine gedruckte Zeile belegt es
  (`AGENTS.md` §3.12 Instanz A: als Erwartung, nicht als Messung geführt).
- **Die neue Assertion fällt still aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×). *Erwartet, zu belegen durch:* die Zeile in `docs/user/e2e-abdeckung.md`
  nennt die Klasse, und die Phase färbt sich rot, wenn die Regel zurückgenommen ist.
  **Ausgang: nicht eingetreten** — `docs/user/e2e-abdeckung.md` Zeile 68 nennt die Klasse;
  die Eingabeseiten-Mutation des Implementers (`make image-mutation` mit der alten,
  unbehobenen Priorität, Commit `fe9d0afa`) färbt die Phase real rot, exakt an der neuen
  Assertion.
- **Ein Kommentar sagt die Regel breiter oder schmaler, als der Code sie trägt**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, verkörpert, 6×). *Erwartet,
  zu belegen durch:* der Reviewer fährt die zugesagten Pfade der drei Kommentare im Code
  nach ([`AGENTS.md`](../../../../AGENTS.md) §3.7). **Ausgang: nicht eingetreten** — Review
  und beide Verifikationen lesen alle vier Stellen (`wiring.go`, zwei Testkommentare,
  `harness/README.md` §Sensors) eigen gegen den Code — konform, kein Informationsverlust,
  keine Übertreibung (Review „Negativbefunde"; zweite Verifikation §7 Zeile 6).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Vertragstabelle des Verdikts (§2) war klein und
  eindeutig genug, um sie als vier `switch`-Fälle mit einem frühen Guard direkt
  umzusetzen — keine Architect-Rückfrage nötig. Der reale, grüne
  `make test-integration`-Lauf hob die Kernaussage vom Verdikt-Status
  „hergeleitet" auf „erprobt": `errors.Is(err, context.Canceled)` trägt real an der
  `postgresstorage`/`sqlexec`-Kette, die Klasse ist am komponierten Prozess real
  `replication`, nicht `storage`. Zwei unabhängige Review-/Verifikationsrunden mit
  insgesamt 14 eigenen Mutationen (Review 6, erste Verifikation 9, zweite
  Verifikation 5 weitere/reproduzierte) fanden zusammen zwei reale, aber enge
  Testlücken (F-1: strukturell unerreichbare Sentinel/`context.Canceled`-Interaktion;
  V-1: der ungebundene frühe Guard bei real erreichbarem Fall) und schlossen beide
  — kein HIGH über den ganzen Vorgang.
- **Was ging anders als geplant:** Die Implementierung bündelte Produktionscode,
  Tests, Runner-Skript und zwei Doku-Erzeugnisse/-Träger in einem Commit
  (`e7210994`) statt der im Plan §3 nahegelegten Datei-Granularität — ein
  vorzeitiger `git add -A` (Review F-2, LOW, kein Traceability-Schaden). Der erste
  `make gates`-Lauf färbte `docs-check` rot: der `next → in-progress`-Move dieses
  Slice hatte sechs Markdown-Links mit fest verdrahtetem `open/`-Pfad
  zurückgelassen (`BEO-PGC/slice-pfad-als-link-in-berichten`) — als Plan-Nachzug
  behoben, bevor der zweite Lauf grün war. Das committete §3.13-Suchlauf-Feld
  driftete **zweimal** durch dieselbe Ursache — eine Fixrunde, die einen neuen
  Testfall auf die bewegte Stelle selbst schrieb, ohne den `diff`-Suchlauf im
  selben Zug erneut zu fahren: einmal unbehandelt (von der ersten Verifikation als
  V-2, MEDIUM, gefunden: `soll=56 ist=59`), einmal — bei der Folge-Fixrunde für
  V-1 — korrekt behandelt (der Implementer maß vor der eigenen Meldung neu und
  trug `61` ein, von der zweiten Verifikation unabhängig bestätigt). Die zweite
  Verifikation fand zusätzlich einen rein redaktionellen Fund (V-4, LOW): die
  DoD-Zeile 1 dokumentierte die V-1-Fixrunde nicht an ihrer eigenen Stelle,
  sondern nur beiläufig im Suchlauf-Absatz — bei dieser Closure nachgezogen (§2,
  DoD-Zeile 1).
- **Steering-Loop-Eintrag (Lerneintrag):** **Geschärfte Regel:** Eine Fixrunde,
  die einen neuen Testfall auf eine §3.13-bewegte Stelle schreibt, fährt
  `make suchlauf-nachmessen` und trägt eine ggf. geänderte `diff`-Zahl **im
  selben Commit** nach — nicht erst, wenn eine spätere Rolle sie beim Nachmessen
  findet. Dieser Slice zeigt beide Seiten am eigenen Leib: die erste Fixrunde
  (`a7d27ddc`, Testfall gegen F-1) tat es nicht und produzierte einen realen
  MEDIUM-Fund (V-2); die zweite Fixrunde (`2e0623bf`, Testfall gegen V-1) tat es
  — die Lehre aus dem ersten Auftreten griff beim zweiten bereits innerhalb
  desselben Slice. **Neuer Sensor:** keiner — `make suchlauf-nachmessen`
  existiert bereits und deckt die Zahlen; was fehlte, war die Disziplin, es nach
  jeder eigenen Änderung an der bewegten Stelle erneut zu fahren, nicht ein
  fehlendes Werkzeug (`AGENTS.md` §3.13 benennt diese Grenze bereits: das
  Werkzeug prüft Zahlen und Stände, nicht ob der Implementer es vor dem Commit
  ein letztes Mal laufen ließ — dieselbe Klasse wie der vierundzwanzigste Beleg
  in `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`). **Benannte
  Spec-Lücke:** keine — `ADR-0049` ist vollständig eingelöst, kein Konflikt mit
  `SPEC-008`/`SPEC-013` gefunden (Review/beide Verifikationen: kein Diff an
  Schwellen, `classifyRunError`, Sentinel-Trennung oder Klassen-Tabelle).
- **Beobachtungs-Register (`../observations/`):**
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` erhält den 25. Beleg
  (`evidence/slice-wal-fehlerschwelle-ausgangsklasse.md`) — die zweifache
  Fixrunden-Erfahrung mit dem §3.13-`diff`-Wert (einmal unbehandelt, V-2 MEDIUM;
  einmal korrekt behandelt), Zähler 24× → 25×, Ausgang unverändert *verkörpert*.
  `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` bleibt unverändert beim
  Ausgang „geplant → `slice-start-vorlauf-grenze`" — der Trigger der
  Neubewertung ist nicht eingetreten (von der zweiten Verifikation bestätigt,
  §8); nichts für diese Closure zu setzen. Kein weiterer Anfall.
- **Folge-Slices:** keiner neu angelegt. Die Kette bleibt:
  `slice-start-vorlauf-grenze` (nächster Slice, `open/`, trägt die Kopplung als
  Start-Trigger, §4 dort) → `slice-transformationen-e2e-abhilfe` →
  `slice-transformationen-betriebsdoku`.
- **Risiken aus §6:** sechs Ausgänge gesetzt — vier „nicht eingetreten"
  (Risiko 1: die Kette trägt real an der echten Stelle, vom Verdikt-Status
  „hergeleitet" auf „erprobt" gehoben; Risiko 2: `storage`-ohne-Abbruch-Fall
  gebunden; Risiko 5: Assertion nicht still ausgeschlossen; Risiko 6: Kommentare
  konform), ein „überwiegend nicht eingetreten, mit einer real gefundenen und
  geschlossenen Lücke" (Risiko 3: V-1 fand die eine Lücke, die Fixrunde und 14
  Mutationen über beide Verifikationsrunden schlossen sie), ein „plausibel,
  unbelegt" (Risiko 4: keine gedruckte Vorher-/Nachher-Laufzeit der Phase,
  aber keine strukturelle Wirkung der Änderung — eine Assertion auf einen
  bereits gelesenen Wert, kein neuer Wartezyklus). Kein Risiko blieb „weiter
  offen".
- **Drei Paarungen:** dieser Slice hat keine Welle; die Roadmap führt
  [welle-transformationen](../welle-transformationen.md) unter *Offene
  Wellen*, das Ereignis kann eintreten: die Closure dieser Welle prüft die
  Paarungen mit. Die Slice-Closure trägt sie zusätzlich jetzt: *Anker:* die
  Regel existiert als committeter Text (`mergeStreamAndWALFaultOutcome` in
  `internal/bootstrap/wiring.go`, die Runner-Phase und ihre Assertion in
  `tools/harness/run-integration-tests.sh`), `harness/README.md` §Sensors
  nennt den Beleg mit der Zeile „Ausgangs-Klasse seit
  slice-wal-fehlerschwelle-ausgangsklasse" (`git grep -o
  'slice-wal-fehlerschwelle-ausgangsklasse' -- harness/README.md` trifft 1).
  *Folge-Slice:* `slice-start-vorlauf-grenze` existiert als Datei in `open/`
  und trägt die Kopplung als Start-Trigger (§4 dort, Zeile 289/294). *Register:*
  jede in dieser Notiz genannte Kennung `BEO-PGC/<slug>` existiert als
  Verzeichnis mit nicht leerem `evidence/` (`zahl-in-traeger-driftet-gegen-die-messung`
  25 Dateien, `adapter-unittest-verdeckt-bootstrap-luecke` 5 Dateien).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Composition Root und Test-Runner sind keine eigenen Sub-Areas —
kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-27 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (5×, Ausgang *geplant*:
`slice-start-vorlauf-grenze`; der fünfte Beleg ist dieser Fall, Verdikt §6 — der Träger
der Klasse ist die Phase in `make test-integration` im Slice, der die Eigenschaft
einführt, hier dieser Slice; kein weiterer Eintrag, der Ausgang bleibt),
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (6×, verkörpert, Risiko §6 und
DoD 3: die Kommentare an `mergeStreamAndWALFaultOutcome` und in Zeile 27 des Tests),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (19×, verkörpert, Risiko §6),
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` (8×, verkörpert; die Aussage über den
Stand des Streams in Verdikt §5 ist als hergeleitet gekennzeichnet, Kernaussage §1),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, Risiko §6),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×, Risiko §6),
`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×, Start-Trigger §4),
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 5×: Repo-Dateien
entstehen über Edit/Write, nie über Umleitung, [`AGENTS.md`](../../../../AGENTS.md) §3.1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
