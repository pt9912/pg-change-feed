# Verifikations-Report: slice-capture-transient-wiederholung — 2026-09-30

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
ADR-Konformität (`ADR-0135`, `ADR-0136`) und Plan-vs-Code-Diff. Review-Artefakt:
[`review-slice-capture-transient-wiederholung.md`](review-slice-capture-transient-wiederholung.md)
(Fixrunde 2, Re-Review `73736af2`). Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md).

**Gegenstand:** Slice-Plan `slice-capture-transient-wiederholung`, Diff `50d9ecc4~1..HEAD`
(`73736af2`): 17 Dateien, davon 9 unter `internal/`. Dieser Lauf ändert weder Code noch Plan
(keine DoD-Häkchen gesetzt); er schreibt nur diesen Report. Die Mutation lief an einer
`git archive`-Kopie im Scratchpad.

## 1. Eigene Sensor-Belege (ungepiped, Exit-Code einzeln gesichert, §3.9)

| Sensor | Exit | Ausgabe |
|---|---|---|
| `make test` (Race-Detector) | 0 | alle Pakete `ok`, `internal/bootstrap` ok |
| `make test-replication` | 0 | `internal/bootstrap` ok in 13,8 s (Phase measure; der Nicht-DB-Lauf braucht 1,6 s, also lief der DSN-gebundene Realtest und übersprang nicht); Tier-Tests `TestSlotReserve…`, `TestWALRetention…`, `TestSourceKeepalive…` PASS |
| `make gates` | 0 | `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%`, `generated-sync: OK`, `sdk-public-doc-check` ok, `a-check: gesamt: 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=…` | 0 | 8 Zeilen stimmen (156/185, 2/1, 55/61, 99/112) |
| `make kommentar-kennungen DIFF=50d9ecc4~1` | 0 | keine Kandidaten (Form, nicht Wahrheit) |

Nicht gefahren: `-v`-Ausgabe des Realtests `TestRunStreamWithRetrySlotStillActive` (der Runner
druckt sie in der Measure-Phase nicht; Nachweis nur indirekt über die Laufzeit, siehe oben).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Wiederholung mit Backoff, Fortsetzung, Erschöpfung `transient` | **getragen** | `TestRunStreamWithRetry{BackoffFolge,Obergrenze,GesamtfensterFuenfMinuten,Stabilitaetsschwelle,AufbauGrenzeGesamtfenster}` gegen Fake-Uhr (Werte 2 s/30 s/5 min/Faktor 2 der ADR); Erschöpfung trägt `ErrTransientExhausted` (`wiring.go:1804` klassifiziert), Test Z. 171/303; Realtest „Slot noch aktiv“ grün. Fortsetzung/Persist-before-ACK: Capture-Service und ACK-Adapter sind im Diff **unberührt** (nur `receive.go`/`wiring.go`), ACK-Port wird je Zyklus neu gesetzt (`ackPort.set`) |
| 2 | Klassen getrennt, Negativtests | **getragen** | `TestRunStreamWithRetryKlassenEndenOhneWiederholung`: configuration, schema, Ordnungsverletzung, storage, permission (Sentinel und 42501), Server-Abweisung, je `waits == 0`; `serverfault_test.go` je SQLSTATE-Klasse |
| 3 | Träger (Kommentar `Run`, `compose.yaml`, Handbuch 1.83, `SPEC-008`) | **getragen** | im Diff vorhanden; `make docs-check` im `make gates` grün |
| 4 | `make gates` grün | **getragen** | Exit 0, eigener Lauf |
| 6 | §3.13-Suchlauf im Plan-Feld | **getragen** | Zahlen nachgemessen, beide Stände |
| 5, 8–13 | Review-Report, Closure, Register, Risiko-Ausgänge, Paarungen | **offen, gehört dem Planner** | Report liegt vor; Closure-Notiz und Risiko-Ausgänge §6/§7 stehen noch leer (erwartet vor `done/`) |

## 3. Plan-vs-Code-Diff

Alle in §3 des Plans genannten Dateien erscheinen im Diff; kein Code-Bestandteil ohne Plan-Zeile
(`walretention.go`: 1 Zeile, Anpassung an `serverFault`). Nicht-Ziele eingehalten: Backfill-Run,
`restart: "no"` und die übrigen Klassen sind unverändert.

## 4. Eigene Mutation

`wiring.go:1151`, Aufbau-Frist `streamRetryMaxDelay` durch `1 * time.Nanosecond` ersetzt
(Scratchpad-Kopie, `go test -race ./internal/bootstrap/`): **grün** (`ok … 1.639s`). N-5 ist damit
unabhängig reproduziert: `runStreamCycle` ist mit übergebenem Parameter getestet, der Aufrufwert in
`Run` hängt an keinem Test.

## 5. Bewertung der bekannten Punkte

- **N-5 (LOW):** bestätigt (Mutation grün). Verletzt die DoD nicht: die DoD verlangt Tests je Grenze
  Wiederholungszahl, Backoff-Folge, Erschöpfung — die Aufbau-Frist ist eine Zusatzgrenze aus
  `ADR-0136`. Schwäche der Bindung, kein DoD-Bruch.
- **F-8 (LOW):** bestätigt. Der Test erzwingt den Fehlschlag des ersten Versuchs indirekt (der Wartezug
  gibt den Slot erst nach einem Fehler frei; ohne Fehlschlag gäbe es keinen zweiten Versuch), prüft
  aber weder SQLSTATE 55006 noch dass die Ursache der Fehler „Slot aktiv“ ist; die Zählung `>= 2`
  schließt Duplikate nicht aus (at-least-once erlaubt sie). Die Fortsetzung an `confirmed_flush_lsn`
  ist damit nur als Lieferung der Folge-Change belegt, nicht als Start an der bestätigten Position.
  DoD-Zeile 1 nennt „Slot noch aktiv“ als Beispiel und verlangt den Realtest mit diesem Fall:
  getragen, aber schwach.
- **N-6 (LOW):** bestätigt. „Run unberührt von der Frist“ ist mit Fake-Stream belegt
  (`stream_cycle_internal_test.go`); im Quelltext liegt `stream.Run(ctx)` außerhalb von `setupCtx`
  (gelesen, `wiring.go:1906–1927`).

## 6. Verdikt

**DoD getragen: ja** für die Zeilen 1–4 und 6 (Code, Tests, Träger, Gates, Suchlauf). Keiner der drei
Punkte verletzt die DoD. Offen sind die Planner-Zeilen (Closure-Notiz, Risiko-Ausgänge §6, Register,
Paarungen).

**Nötiger Nachzug (optional, minimal, kein Blocker):**
1. N-5: in `Run` den Aufbau-Frist-Wert als benannte Konstante/Variable führen und einen
   Quelltext- oder Verdrahtungstest daran binden (oder die Mutation im Plan als bekannte Grenze benennen).
2. F-8: im Realtest die Ursache des ersten Fehlers prüfen (`errors.As` auf `*pgconn.PgError`,
   Code `55006`, im Zyklus vor dem Rückgabewert mitschneiden) und `count == 2` statt `>= 2`
   fordern bzw. die Start-LSN des zweiten Versuchs gegen `confirmed_flush_lsn` lesen.
3. Planner: Risiko „Grenzwerte sind Startwerte ohne Messung“ in §6 als gesetzt, nicht gemessen
   ausweisen.
