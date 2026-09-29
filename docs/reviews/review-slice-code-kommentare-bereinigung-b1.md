# Review-Report: slice-code-kommentare-bereinigung B1 (T1–T7) — 2026-09-29

**Review-Art:** Code — geprüft gegen den Slice-Plan
`slice-code-kommentare-bereinigung` (Lifecycle-Ort: `docs/plan/planning/`),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und die
Hard Rules `AGENTS.md` §3.7/§3.12/§3.13/§3.9 (Modul 10). Keine
DoD-Verifikation — das ist die Aufgabe des Verifiers (Modul 11).

**Gegenstand:** Diff `e63a1afd..3455a23c` — die neun Tranchen-Commits T1–T7
(`7b310270`, `a5280dc7`, `8416440f`, `73df6635`+`1163f55e`, `e9a398a9`,
`1be3fc5f`, `f1dd169f`+`195c3006`) samt Plan-Nachzug `3455a23c`. Im selben
Diff, aber nicht Prüfgegenstand dieses B1-Reviews: `47792d6e`
([`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md)),
`e658a1bd`, `5d572c44` (Plan-/Register-Züge).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1` (2026-09-28) ·
**Modell:** GLM (Claude-Agent-SDK, Typ `reviewer`) · **Datum:** 2026-09-29

**Eingangs-Kontext:**

- Slice-Plan @ `3455a23c` (inkl. Suchlauf-Feld und DoD-Liefer-Punkt-1-Haken)
- `AGENTS.md` §3.7 (Kommentar-Klassen, ein Anker), §3.12 (Zahlen-Ursprung),
  §3.13 (Suchlauf, Träger-Meldung), §3.9 (Exit-Disziplin)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `harness/sensors/kommentar-kennungen.md`, `harness/sensors/suchlauf-nachmessen.md`,
  `harness/sensors/coverage-gate.md` (Träger der gemeldeten Lokatoren)
- `harness/conventions.md` (MR-000) · Baseline `v6.13.0` ·
  `regelwerk/modul-10-review-harness.md`

---

## Findings

### F-1 — Kürzung lässt unauflösbare Herkunfts-Referenz zurück (`config_file.go`)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 („Herkunft im Go-Kommentar ist ein Feld, kein
  Absatz" — ein auflösbares Feld); Plan §1 „Regeln je Kandidat" (gekurzt wird
  auf den Anker, der die Norm trägt); Skill-HIGH „Kommentar trägt keine der
  Kommentar-Klassen" (Teilersetzung lässt den Rest stehen)
- `pfad`: `internal/bootstrap/config_file.go:6`
- `befund`: Der gekürzte Block ersetzt `ADR-0088` Festlegung 2 durch „die
  Decision Festlegung 2" — ein Verweis, den kein Leser auflösen kann: der
  Block-Anker ist `ADR-0052`, und `ADR-0052` trägt keine „Festlegung 2"
  (nachgelesen: die ADR kennt nur „Entscheidung"). Die Herkunft der
  zulässigen Feldmenge (normativ `ADR-0088` Festlegung 2, nachgelesen) ist
  damit aus dem Kommentar verloren; „die Decision" ist der Rest einer
  Teilersetzung. `make kommentar-kennungen` meldet die Stelle nicht (0
  Kandidaten in T6, nachgemessen) — kein Gate fängt das.
- `verifizierbar`: nein — Lese-Handlung; das Werkzeug zählt Kennungen, nicht
  hängende Verweise
- `klasse`: „Kürzung lässt hängende Herkunfts-Referenz zurück"

### F-2 — Hängende Teil-Referenz „Teilfrage 4/5" nach Anker-Entfernung

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7; Plan §1 „Regeln je Kandidat"
- `pfad`: `tools/harness/natssub/main.go:10`
- `befund`: Der Godoc-Block trug „seit `ADR-0100` … (Teilfrage 4/5)"; die
  Kürzung entfernte die ADR und ließ „(Teilfrage 4/5)" stehen. Der Block
  trägt nur noch `LH-FA-SST-007` — „Teilfrage 4/5" löst nicht auf (eine
  `LH-*`-Kennung hat keine Teilfragen; die Eltern-ADR
  [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md) ist
  aus dem Block weg). Dieselbe Klasse wie F-1, milde Form: der Satz trägt
  die Stelle, nur der Verweis ist tot.
- `verifizierbar`: nein — Lese-Handlung
- `klasse`: „Kürzung lässt hängende Herkunfts-Referenz zurück"

### F-3 — Hängende Teil-Referenz „(Festlegung 2)" resolviert falsch

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7; Plan §1 „Regeln je Kandidat"
- `pfad`: `examples/nats-client/main.go:62`
- `befund`: Der Block trug „(`SPEC-018`) … (`ADR-0079` Festlegung 2)"; die
  Kürzung entfernte `ADR-0079` und ließ „(Festlegung 2)" stehen. Der
  Block-Anker ist `SPEC-018` — und `SPEC-018` trägt keine Festlegungen
  (nachgelesen: 0 Treffer im Abschnitt). Der Leser resolviert die
  Teil-Referenz auf den falschen Anker oder gar nicht; die Herkunft
  (`ADR-0079` Festlegung 2) ist verloren. Dieselbe Klasse wie F-1/F-2.
- `verifizierbar`: nein — Lese-Handlung
- `klasse`: „Kürzung lässt hängende Herkunfts-Referenz zurück"

### F-4 — Umformulierung `rejectionMessage`: „Spalten-Ordnung der Antrags-Tabelle" trägt nicht

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7 (Spec-Wiedergabe in eigenen Worten); Plan §3
  Übergabe aus `slice-antragsqueue-lesefehler-failed`; Skill-MEDIUM „Herkunft
  … oder Spec-Wiederholung"
- `pfad`: `internal/adapters/driven/postgresstorage/sqlexec/translate.go:416`
- `befund`: Die plan-gemäße Umformulierung bindet die Reihenfolge Quelle,
  Schema, Tabelle, Antragsart, Spalte — korrekt am Switch und an der
  Verletzungstabelle von `SPEC-019` (nachgelesen, beide stimmen überein) —
  aber der eingeschobene Satz „bildet die Spalten-Ordnung der
  Antrags-Tabelle ab" trägt nicht: in `cdc.administration_request` liegt
  `column_name` an fünfter, `request_kind` an achter Stelle (nachgelesen in
  `tools/schema/schema.yaml`) — die Switch-Reihenfolge gibt das nicht wieder.
  Die Übergabe band die Reihenfolge an die Tabelle von `SPEC-019`; die
  Umformulierung nennt stattdessen eine Spalten-Ordnung, die sie nicht trägt.
- `verifizierbar`: nein — Lese-Handlung gegen Schema und Pflichtenheft
- `klasse`: „ungenau umgeformte Spec-Wiedergabe"

## Negativbefunde

- geprüft, ohne Befund: **Diff-Form aller neun Tranchen-Commits** — jede
  +Zeile beginnt nach Leerraum mit `//`, jede −Zeile ist Kommentar; keine
  Endkommentar-Zeile, keine Nicht-Go-Datei in den Tranchen-Commits
  (111 Go-Dateien, 820+/870− im Gesamt-Diff).
- geprüft, ohne Befund: **Zahlen-Ursprung** — Nachmessen an beiden Ständen
  über `make kommentar-kennungen COUNT=1 TESTS=exclude PATHS=…`: vor (Stand
  `e63a1afd`) 77/38/36/71/89/53/19, Summe 383, T8 188; nach je Tranche 0,
  Summe 0, T8 188 — die Zahlen des Implementers tragen, kein Drift. Die
  Plan-Zahlen (Stand `d13ab81e`, Summe 400, T8 197) sind als ersetzt
  deklariert und tragen ihren Ursprung.
- geprüft, ohne Befund: **B1/B2-Schnitt** — keine `*_test.go` im Diff; T8
  (188 Kandidaten) an beiden Ständen unverändert; die Plan-§4-Provision
  (Schnitt entlang T7|T8, B1 = T1–T7, B2 = T8 samt Restmenge und
  Erzeugnis-Lauf) lag vor dem Bereich committet vor.
- geprüft, ohne Befund: **Suchlauf-Feld** — `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-code-kommentare-bereinigung.md`:
  8/8 Zeilen OK, Exit 0 (4 Zeilen am Plan-Stand `7b70b34a`, 4 Zeilen am
  Arbeitsbaum); die Erklärung zur `ff.`=1-Zeile (Stringliteral im
  Tabellentestfall `tools/harness/kommentar-kennungen/main_test.go:88`,
  nachgelesen) ist zutreffend.
- geprüft, ohne Befund: **Träger-Meldung des Implementers** — der
  Lokator-Drift in `harness/sensors/coverage-gate.md` Z. 249 ist real
  gemessen: `main.go:44/62/86/101` (Healthcheck, RegisterConsumer,
  AcknowledgeConsumer, Diagnose) stehen jetzt bei 50/68/92/107 (je +6 durch
  T6). Die Meldung an den Planner mit Adresse und Frist (Closure) ist
  berechtigt; sie ist offen und muss bei der Closure eingelöst sein.
- geprüft, ohne Befund: **Chronik-/Vorher-Nachher-Sprache in +Zeilen** —
  Suche nach Chronik-/Vorher-Nachher-Mustern (`früher`, `bisher`, `nun`,
  `jetzt`, `nicht mehr`, `bislang`, `zuvor`, `vorher`, `slice-<NNN>`,
  `welle-<NN>`): drei Treffer, keiner eingeführt (mapper.go „nicht mehr
  getragene Bindung" = Zustandsbeschreibung; die „bislang anonyme"-Stellen
  in `wiring.go` und `natssub` sind Bestand bzw. Re-Wrap).
- geprüft, ohne Befund: **Stichproben gegen den Code je Tranche** (Zahl der
  Stichproben: 25 Blöcke/Stellen) — T1 `changestore.go` (5 Blöcke; Zusage
  „Persistenzfehler endet ohne Source-ACK" nachgefahren: `CaptureService.
  Capture` persistiert vor dem ACK), T2 `errors.go` (8 Blöcke), T3
  `acknowledge/backfill service.go` (4 Blöcke), T4 `receive.go`/`mapper.go`
  (`confirmIdle`-Übergabe aus `slice-capture-leerlauf-quellbelege` sauber
  umgesetzt — Begründungssatz entfernt, Regel und ein Anker bleiben),
  T5 `translate.go` (7 Blöcke — Befund F-4), T6 `main.go`/`wiring.go`/
  `backfill.go` (8 Blöcke — Befund F-1), T7 `natssub`/`nats-client`/
  `nats-stream-client` (5 Blöcke — Befunde F-2/F-3).
- geprüft, ohne Befund: **Sensors** — `make test` (Exit 0), `make fmt-check`
  (295 Go-Dateien, alle formatiert, Exit 0), `make a-check` (0 Befunde,
  Exit 0), `make kommentar-kennungen` (11 Läufe, s. o.).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Kürzung lässt hängende
Herkunfts-Referenz zurück" (3×) · „ungenau umgeformte Spec-Wiedergabe" (1×)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) und F-2/F-3/F-4 (MEDIUM) ziehen eine
Fixrunde am Implementer nach sich. Der Gegenstand ist schmal: je Stelle den
Eltern-Anker wieder einsetzen oder den Satz als Stellen-Beschreibung
umformen (das Muster liefert die `confirmIdle`-Stelle in T4 — dort ist die
Kürzung sauber gelungen). Die DoD-Zeile „Review durchgeführt" bleibt offen:
dieser Report ist der erste von zwei Reviews (T8/B2 steht aus), und die
Fixrunde läuft regulär über Schritt 21 des Implementer-Ablaufs.

**Übergabe:** Findings gehen an den Implementer (Fixrunde T1–T7, vor T8);
die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler (neue Klasse: „Kürzung lässt hängende Herkunfts-Referenz
zurück" — 3×, erstes Auftreten; „ungenau umgeformte Spec-Wiedergabe" — 1×,
erstes Auftreten). Dieser Report ist ein Lauf-Beleg.

**Verbleibende Risiken (nicht Findings):**

- `harness/sensors/coverage-gate.md` Z. 249 trägt die gedrifteten Lokatoren
  bis zum Planner-Nachzug (Frist: Closure dieses Slice; kein Gate liest die
  Stelle).
- Der `make test-integration`-Lauf in T8 schreibt
  `docs/user/e2e-abdeckung.md` neu — die Lokatoren der Datei verschieben
  sich mit den B2-Kürzungen (Erzeugnis, im Plan vorgesehen).
- F-4 zeigt, dass die Umformulierung von Sätzen (nicht nur das Kürzen von
  Ankern) eine Lese-Handlung braucht — für T8 (197→188 Kandidaten, fünf
  Teil-Commits) gilt dieselbe Probe.
