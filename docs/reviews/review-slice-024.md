# Review-Report: slice-024 — 2026-09-12

**Review-Art:** Design — geprüft gegen Plan (`slice-024`) + neu geschriebene
`ADR-0049` (Accepted) + Konventionen (`AGENTS.md` §3, `harness/README.md`
§Referenz-Richtung/Matrix). Kein Code-Diff: dieser Slice liefert eine
Entscheidung und zwei Spec-Präzisierungen, keine Implementierung.

**Gegenstand:** Commits `741f845` (ADR-0049 neu), `52e207d`
(`spec/pflichtenheft.md` `SPEC-008`/`SPEC-013` geschärft), `988ab9a`
(Slice-Plan-Nachzug: DoD-Häkchen, Plan-Nachzug-Tabelle)

**Skill:** `.harness/skills/reviewer.md` @ `e9159ff` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-024-adr-fehlerklassen-schwellen-praezisierung.md`
  §1–§8 (inkl. Plan-Nachzug)
- `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md` (Accepted)
- `docs/plan/adr/0023-fehlerklassifikation.md` (Accepted, geprüft auf
  Widerspruch/Immutabilität)
- `spec/pflichtenheft.md` §3 (`SPEC-013`), §4 (`SPEC-008`), §5 (Metriken-Prosa)
- `spec/lastenheft.md` (`LH-QA-REL-003`, `LH-FA-CAP-004`)
- `docs/plan/planning/welle-7.md` (§3 Ende-zu-Ende-Beleg, referenziert vom
  Slice-Kopf)
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/`
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.4, §3.5, §3.7)
- `harness/conventions.md` (MR-000/MR-001)
- `.d-check.yml` (Matrix-Regel `{from: spec, to: adr, allow: false}`)
- Code (als Sachverhaltsprüfung, nicht als Review-Gegenstand):
  `internal/adapters/driving/replication/mapper/mapper.go`,
  `internal/adapters/driving/replication/receive/receive.go`,
  `internal/adapters/driven/postgresack/ack.go`,
  `internal/application/port/outbound/replicationack.go`,
  `internal/bootstrap/wiring.go` (`classifyRunError`)

---

## Findings

### F-1 — Transport-/Verbindungsstörung-Kategorie fasst im Kontext-Abschnitt heterogenere Fehlerursachen, als die Aufzählung nahelegt

- `kategorie`: INFO
- `quelle`: Maintainability (Präzisions-Hinweis, kein Mangel)
- `pfad`: `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md`
  (Abschnitt Kontext, zweiter Aufzählungspunkt) vs.
  `internal/adapters/driving/replication/receive/receive.go:313,319`
  (`leeres CopyData`, `ParseXLogData`-Fehler) und
  `internal/adapters/driven/postgresack/ack.go:64,98` (`nil`-Verbindung,
  Bestätigung mit Zero-Value-Position)
- `befund`: Die ADR-Aufzählung nennt „Verbindungsaufbau, Start des Streams,
  Katalogabfragen und Slot-Verwaltung, Keepalive-Empfang/-Antwort,
  Nachrichtenempfang, Quell-Bestätigung" als Beispiele für
  Transport-/Verbindungsstörung. Tatsächlich trägt derselbe Sentinel
  (`receive.ErrReplication`/`outbound.ErrReplication`) auch reine
  Eingabe-/Zustands-Guards, die keine Verbindungsstörung im engeren Sinn
  sind (leeres `CopyData`, ein Parse-Fehler der Protokoll-Hülle, ein
  `nil`-Constructor-Argument, eine Bestätigung mit Zero-Value-Position).
  Die **Vollständigkeit der Zweiteilung ist davon nicht berührt** — bei
  Sentinel-Granularität (fünf Go-Sentinel-Variablen) ist die Zuordnung
  exhaustiv und korrekt, belegt durch die `switch`-Anweisung in
  `classifyRunError` (`internal/bootstrap/wiring.go:400-421`), die alle
  fünf und keine weiteren Sentinels auf `ErrorClassReplication`
  abbildet. Nur die Kontext-*Erzählung* ist enger als der reale
  Aufrufer-Kreis des Sentinels.
- `verifizierbar`: ja — `grep -n "ErrReplication" internal/adapters/driving/replication/receive/receive.go internal/adapters/driven/postgresack/ack.go` zeigt alle Aufrufstellen; `internal/bootstrap/wiring.go:400-421` zeigt die vollständige, geschlossene Sentinel-Liste.
- `klasse`: „ADR-Kontext-Aufzählung enger als der reale Aufrufer-Kreis des zitierten Sentinels"

## Negativbefunde

- geprüft, ohne Befund: **Sentinel-Zuordnung, sachliche Korrektheit** —
  `mapper.ErrChangeWithoutBegin`/`ErrCommitWithoutBegin`/`ErrBeginWithoutCommit`
  entstehen ausschließlich, nachdem der Assembler ein bereits vollständig
  decodiertes `pgoutput`-Ereignis in eine widersprüchliche Reihenfolge
  einordnen müsste (`internal/adapters/driving/replication/mapper/mapper.go:104,114,129`)
  — echte Stream-Ordnungs-Verletzung, kein Übertragungsproblem.
  `receive.ErrReplication`/`outbound.ErrReplication` entstehen an der
  Verbindung/Quelle (Connect, `START_REPLICATION`, Katalog, Slot,
  Keepalive, `ReceiveMessage`, `Acknowledge`) — die ADR-Behauptung trägt
  (siehe F-1 für die eine Präzisions-Nuance, die den Befund nicht kippt).
  Kein Sentinel bleibt unklassifiziert; `classifyRunError` ist die
  geschlossene Referenzliste und deckt sich mit der ADR-Aufzählung.
- geprüft, ohne Befund: **ADR-Immutabilität / Widerspruch zu `ADR-0023`** —
  `ADR-0023` legt nur die sieben Klassennamen fest (`transient`,
  `configuration`, `permission`, `schema`, `storage`, `replication`,
  `internal`) und trifft keine Aussage über Sub-Verhalten innerhalb einer
  Klasse. `ADR-0049` verändert weder Klassennamen noch
  Klassen-Zugehörigkeit einzelner Sentinels (alle fünf bleiben
  `replication`), sondern präzisiert nur die Handlungs-Konsequenz
  innerhalb dieser einen Klasse. Kein `Supersedes ADR-0023` nötig, keiner
  behauptet. `ADR-0023` selbst wurde nicht editiert (`git show 741f845
  --stat` enthält keine Änderung an `0023-fehlerklassifikation.md`).
- geprüft, ohne Befund: **Referenz-Richtung (SDP) / `.d-check.yml`
  Matrix-Regel** — `git show 52e207d` enthält im geänderten Text keinen
  `ADR-\d{4}`-Token; `spec/pflichtenheft.md` nennt „ADR" nur zweimal in
  Meta-Prosa außerhalb der geänderten Zeilen („kein ADR-Rückzeiger",
  „kein ADR- und kein Slice-Verweis"), keine konkrete ADR-Kennung. `make
  gates` (d-check, Modul `matrix`) läuft in dieser Sitzung grün mit
  (231 Dateien, 0 Befunde) — die Sperre `{from: spec, to: adr, allow:
  false}` ist nicht verletzt.
- geprüft, ohne Befund: **Spec-interne Konsistenz `SPEC-013`** — §5 verwies
  bereits vor diesem Slice per Prosa auf `SPEC-013` als Initialwert-Quelle
  „für WAL-Rückstand und Capture-Lag" (Zeile nach der Metriken-Tabelle),
  obwohl die `SPEC-013`-Zeile in §3 bis dahin nur `CDC_LAG_THRESHOLDS`
  (Capture-Lag) trug. Der Commit `52e207d` erweitert exakt diese
  bestehende Zeile (Name-Feld verallgemeinert auf `CDC_THRESHOLDS`,
  Wert-Spalte um die WAL-Rückstand-Schwellen ergänzt) statt eine neue
  §5-Zeile anzulegen — löst die im ADR-Kontext benannte Inkonsistenz
  auf, wie der Plan-Nachzug in `slice-024` §3 behauptet.
- geprüft, ohne Befund: **Kein Chronik-Sprachgebrauch in
  `spec/pflichtenheft.md`** — der geänderte Text enthält keine
  `seit slice-*`/`seit welle-*`-Referenz; beide geänderten Zeilen
  beschreiben den geltenden Zustand indikativ.
- geprüft, ohne Befund: **Schwellenwert-Begründung und
  Re-Evaluierungs-Trigger** — die ADR benennt eine nachvollziehbare
  Stil-Analogie (Faktor 10 zu `SPEC-013`s bestehendem Faktor-12-Muster bei
  `cdc_capture_lag`) und einen realen Risikobezug (Disk-Erschöpfung durch
  unbegrenztes WAL-Wachstum bei inaktivem Slot, mangels
  `max_slot_wal_keep_size`). Der Wert ist explizit **nicht** als
  Produktionsmessung deklariert (`<!-- d-check:status-provenance -->`
  markiert), und der Re-Evaluierungs-Trigger für (b) ist korrekt
  **nicht permanent** gesetzt (Ende-zu-Ende-Test aus `welle-7` §3 oder
  späterer Produktionsbetrieb widerlegt die Werte → Folge-ADR mit
  `Supersedes ADR-0049`); Trigger für (a) korrekt **permanent** (folgt der
  Adapterstruktur, unabhängig von Zahlenwerten). Das im Slice-Plan §6 als
  Risiko benannte „Stil-Analogie statt Produktionsdaten" ist damit in der
  ADR selbst durch den Trigger gedeckt.
- geprüft, ohne Befund: **Provenienz / Template-Konformität** — `ADR-0049`
  führt exakt die sieben Abschnitte der vendored Ziel-Form
  (`.harness/baseline/v6.5.0/templates/docs/plan/adr/NNNN-titel.template.md`):
  Kontext, Entscheidung, Verglichene Alternativen (zwei Tabellen mit je
  ≥ 3 Optionen inkl. „nichts tun"), Konsequenzen, Fitness Function,
  Re-Evaluierungs-Trigger, Geschichte — keine freihändig erfundene
  Gliederung. Kopfzeilen-Felder (`Status`/`Datum`/`Autor`/`Bezug`/
  `Schärft`/`Regeln`) vollständig; `Bezug` nennt `LH-QA-REL-003` und
  `LH-FA-CAP-004` (beide existieren in `spec/lastenheft.md`, Zeilen 325
  und 1025) sowie `ADR-0023` als verwandte, nicht ersetzte Entscheidung.
  Das Fehlen von `#anker`-Suffixen an `Bezug`/`Schärft` ist im Plan-Nachzug
  begründet (Bestandskonsistenz mit allen 48 vorherigen ADRs, keine
  aktivierte Anker-Reifestufe) — nachvollziehbar, kein Verstoß.
- geprüft, ohne Befund: **ADR-Index** — `docs/plan/adr/README.md` trägt
  die neue Zeile `ADR-0049 | ... | Accepted | 2026-09-12 |
  [0049-...md](0049-...md)`, Nummer `0049` korrekt als nächste freie nach
  `0048` vergeben (Plan-Nachzug ersetzt den Platzhalter `00NN` konsistent).
- geprüft, ohne Befund: **Folge-Slice-Existenz** — `slice-025` und
  `slice-026`, in `ADR-0049` §Konsequenzen als Folgepflicht benannt,
  existieren real als Dateien in `docs/plan/planning/open/`.
- geprüft, ohne Befund: **DoD-Häkchen-Ehrlichkeit (`988ab9a`)** — alle
  vier neu auf `[x]` gesetzten Punkte (ADR Accepted, `SPEC-008`-Zeile,
  `SPEC-013`-Eintrag, `make gates` grün) sind durch die tatsächlich
  gelieferten Artefakte gedeckt; das fünfte (`Kein Doku-Update jenseits
  von spec/pflichtenheft.md und dem ADR-Index nötig`) trifft zu — einzig
  `spec/pflichtenheft.md`, das neue ADR-Dokument, der ADR-Index und der
  Slice-Plan selbst wurden geändert. Die weiterhin offen gelassenen
  Punkte (Review, Closure-Notiz, Reconciliation-/Beobachtungs-Register,
  Risiken-Ausgänge, drei Paarungen) sind konsistent mit dem
  Lifecycle-Zustand `in-progress` — Closure ist nicht Teil dieser drei
  Commits.
- geprüft, ohne Befund: **Beobachtungs-Register unverändert wie
  angekündigt** — `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md`
  trägt weiterhin `weiter offen`, 1× (`evidence/slice-020.md`); der
  Plan-Nachzug kündigt genau das an (Fortschreibung ist
  Closure-/Planner-Arbeit, nicht Teil dieses Architect-Laufs).
- geprüft, ohne Befund: **Out-of-Scope-Disziplin §1** — die drei
  benannten Ausschlüsse (Implementierung → `slice-025`/`slice-026`,
  Verhalten der Stream-Ordnungs-Verletzung bleibt Bestand, keine
  Lastenheft-Change-Request-Zeremonie) sind je mit Begründung versehen
  und werden von den drei Commits eingehalten — keine Implementierung,
  kein Verhaltenswechsel für die Mapper-Sentinels, kein Lastenheft-Edit.
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs
  nennen `ADR-0049` bzw. sind im Range der letzten 5 Commits erfasst;
  kein `SPEC-*`/`ARC-*` im Betreff. `make commit-traceability` lief in
  dieser Sitzung grün mit (5 Commits, Betreffs ohne Struktur-ID).
- geprüft, ohne Befund: **Sub-Area-/Modus-Prüfung (§8 des Slice-Plans)**
  — einzige berührte Sub-Area ist die Default-`PGC` (Greenfield); die
  Sichtung der offenen Beobachtungen benennt korrekt die drei
  Treffer (`spec008-replication-luecke`, `walsender-wirksamkeit`,
  `rollen-test-abdeckungsluecken`) und stuft zwei davon zutreffend als
  nicht einschlägig ein, ohne 3×-Schwelle zu erreichen.
- geprüft, ohne Befund: **`make gates`** (in dieser Sitzung selbst
  ausgeführt) — `baseline-verify` (54 Dateien OK), `docs-check`
  (231 Dateien, 0 Befunde, inkl. `matrix`/`ids`/`anchors`), `commits`
  (5 Commits, Betreffs ohne Struktur-ID), `commit-traceability.sh` OK,
  `a-check` (0 Befunde) — alle vier inneren Gates grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** ADR-Kontext-Aufzählung enger als der
reale Aufrufer-Kreis des zitierten Sentinels

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM, kein LOW. Die
einzige Einstufung (F-1, INFO) beschreibt eine Präzisionsnuance in der
Kontext-Erzählung der ADR, die die Vollständigkeit und Korrektheit der
tatsächlichen Sentinel-Zuordnung nicht berührt — kein Nacharbeits-Bedarf.

**Übergabe:** Findings gehen an den Implementer/Architect. Die
Finding-Klasse aus F-1 geht zusätzlich in die Slice-Closure §7 und von
dort in den Zähler (`docs/plan/planning/observations/`). Dieser Report
ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD-/Spec-Konformität
(inkl. der noch offenen DoD-Punkte: Closure-Notiz, Beobachtungs-/
Reconciliation-Register, Risiken-Ausgänge, drei Paarungen) prüft der
Verifier separat.
