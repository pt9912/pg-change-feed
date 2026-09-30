# Verifikations-Report: slice-capture-retry-aufbau-frist-bindung — 2026-09-30

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
ADR-Konformität ([`ADR-0136`](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
Festlegung 2) und Plan-vs-Code-Diff. Review-Artefakt:
[`review-slice-capture-retry-aufbau-frist-bindung.md`](review-slice-capture-retry-aufbau-frist-bindung.md).
Formvorbild: [`verifikation-slice-capture-transient-wiederholung.md`](verifikation-slice-capture-transient-wiederholung.md).

**Gegenstand:** Slice-Plan `slice-capture-retry-aufbau-frist-bindung`, Diff `b56b0e87..HEAD`
(`26fd3c55`): 4 Dateien, davon 2 unter `internal/` (`wiring.go` +4/−1, ein Test +32). Dieser Lauf
ändert weder Code noch Plan (keine DoD-Häkchen); er schreibt nur diesen Report. Mutationen liefen an
einer `git archive`-Kopie im Scratchpad (Kopie per `sed … > Datei`, kein `-i`), Race-Image
`golang:1.27@sha256:b475…` mit `--network none`.

## 1. Eigene Sensor-Belege (ungepiped, Exit-Code einzeln gesichert, §3.9)

| Sensor | Exit | Ausgabe |
|---|---|---|
| `make test` (Race-Detector) | 0 | alle Pakete `ok`, `internal/bootstrap` ok (1.691s) |
| `make gates` | 0 | `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung`, `a-check gesamt: 0 Befund(e)` (Tail; Exit in Datei gesichert) |
| `make suchlauf-nachmessen PLAN=…` | 0 | 6 Zeilen stimmen (4/3/0/8/3/3, beide Stände) |

Nicht gefahren: `make test-replication` (Änderung berührt keinen DB-Pfad; Konstante mit unverändertem
Wert 30 s).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Wert als benannte Größe, durch Test gebunden; Mutation färbt `make test` rot | **getragen** | `streamSetupTimeout = 30 * time.Second` (`wiring.go:1846`), Aufruf `wiring.go:1151`; Mutationen §4 |
| 2 | `make gates` grün | **getragen** | Exit 0, eigener Lauf |
| 3 | Review-Report | **getragen** | liegt vor (Häkchen bereits gesetzt) |
| 4–7 | Closure-Notiz, Register, Risiko-Ausgang §6, drei Paarungen | **offen, gehört dem Planner** | §6/§7 tragen noch Platzhalter (erwartet vor `done/`) |

## 3. Plan-vs-Code-Diff

Beide Plan-Zeilen (§3) erscheinen im Diff, nichts darüber hinaus. Wert unverändert 30 s (Nicht-Ziel
eingehalten; `streamRetryMaxDelay` = 30 s, `streamRetryStableAfter` unberührt). ADR-0136 Festlegung 2
(Aufbau-Frist 30 s) konform. Suchlauf-Feld im Plan stimmt mit der Nachmessung überein.

## 4. Eigene Mutationen (Scratchpad-Kopie, `go test -race ./internal/bootstrap/`)

| # | Mutation | Ergebnis |
|---|---|---|
| M1 | Wert: `streamSetupTimeout = 31 * time.Second` | **rot**, `TestRunSourceTextPassesTheSetupTimeoutToTheStreamCycle`: „streamSetupTimeout = 31s, wollen 30s" |
| M2 | Aufruf in `Run`: `runStreamCycle(attemptCtx, 10*streamRetryWindow, …)` | **rot**, derselbe Test: „zweites Argument … ist nicht streamSetupTimeout" |
| M3 | lokales Shadowing `streamSetupTimeout := 10 * streamRetryWindow` direkt vor dem Aufruf | **grün** (`ok … 1.650s`) |

## 5. Bewertung F-1 (LOW)

M3 reproduziert F-1: der Quelltext-Test liest nur den Bezeichner des zweiten Arguments, nicht seine
Auflösung. **Die DoD wird nicht verletzt.** Sie verlangt „die Mutation des Werts im Aufruf in `Run`
färbt `make test` rot"; die beiden Lesarten — Wert der Größe ändern (M1) und im Aufruf einen anderen
Wert einsetzen (M2) — färben rot, mit genannter Stelle, Instanz und Farbe. Shadowing ist eine dritte,
absichtsvoll umgehende Form, die der Test-Kommentar als Grenze benennt („liest die Gestalt des
Quelltexts, kein Lauf"; Plan §6 benennt dieselbe Schwäche). Das Risiko aus Plan §6 ist damit belegt
(Wert, nicht Name, färbt rot), aber dessen Ausgang bleibt die Bindung der Schreibweise. Empfehlung
ohne Blockerwirkung: in der Closure-Notiz/§6 M3 als bekannte Grenze ausweisen oder den Test um die
Prüfung ergänzen, dass `Run` keine lokale Deklaration des Namens `streamSetupTimeout` enthält.

## 6. Verdikt

**DoD getragen: ja** für die Zeilen 1–3 (Code, Bindung, Gates, Review). Kein Blocker.

**Nötiger Nachzug (Planner):**
1. DoD-Zeilen 1 und 2 abhaken (Belege: §1, §4); Closure-Notiz mit Lerneintrag, Register-Vermerk
   (`negativtest-ohne-bindung-an-seine-eingabe` ist einschlägig), Risiko-Ausgang §6 schreiben.
2. Im Risiko-Ausgang §6 M3 (Shadowing bleibt grün) als bekannte Grenze der Quelltext-Bindung benennen,
   oder als optionaler Folge-Schritt den Shadowing-Test ergänzen.
