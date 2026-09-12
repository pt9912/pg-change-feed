# Review-Report: Kommentar-Bereinigung `199396a`/`ca61cfb` (slice-018-Umfeld) — 2026-09-12

**Review-Art:** Code — gegen Plan/Konventionen (Hard Rule 3.7, `AGENTS.md`
§3.7 „Ein Kommentar beschreibt, was da ist"), nachträglich ausgelöst durch
Verifier-Finding VF-2 (`docs/reviews/verify-slice-018.md`): `ca61cfb` lag
nach Review-Schluss von `review-slice-018.md` und wurde von keinem
Reviewer-Auge geprüft, bevor es in Richtung Closure ging.

**Gegenstand:** Commit `199396a` (3 Dateien) und Commit `ca61cfb`
(18 Dateien) — beide bereits auf `main`, beide entstanden im Umfeld von
slice-018 (`docs/plan/planning/in-progress/slice-018-commit-zeitstempel-store-adapter.md`),
aber nicht Teil von dessen ursprünglicher Diff-Range (`f35a94d`).

**Skill:** `.harness/skills/reviewer.md` @ Accepted (Repo-Stand 2026-09-12)
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Diff `git show 199396a`, `git show ca61cfb` (vollständig gelesen, nicht nur Stat)
- `docs/reviews/verify-slice-018.md` (VF-2 — Auslöser dieses Nachtrags-Reviews)
- `internal/adapters/driven/postgresstorage/queries/queries.go` (`InsertTransaction`),
  `store.go`, `mapper/mapper.go` (realer Code-Stand, gegen den die korrigierte
  Aussage in `nacharbeit-observability.sql` geprüft wurde)
- `AGENTS.md` §3.7 (Kommentar-Klassen), §6 Minimal Agent Workflow
- `LH-FA-ADM-004` (referenzierte Anforderung in beiden Commit-Messages)
- `harness/conventions.md` (MR-000 ID-Schema, für die Traceability-Prüfung)

---

## Prüfschritte und Ergebnis

**1. Sind die Commits rein kommentar-/prosa-ändernd?**
Eigener Zeilen-Scan (`git show --unified=0`, jede `+`/`-`-Zeile gegen
Kommentar-Präfix `//` (Go), `--` (SQL) bzw. `#` (YAML) geprüft, kein
Sampling): **0 Nicht-Kommentar-Zeilen** in beiden Commits — bestätigt
unabhängig vom Verifier-Befund. Kein `CREATE`/`ALTER`/`INSERT`/`SELECT`,
keine Go-Signatur, kein Funktionskörper betroffen.

**2. Fachliche Korrektheit der korrigierten Aussage
(`tools/schema/nacharbeit-observability.sql`)**
Die alte Fassung behauptete, `committed_at` trage nur die
Persistenz-Zeit der Instanz (DB-DEFAULT) und `cdc_capture_lag_approx`
messe deshalb nur „Pipeline-Frische", nicht den realen
Quell-Commit-Abstand. Gegen den aktuellen Code geprüft:
`queries.go` `InsertTransaction` übergibt `committed_at` als expliziten
Parameter (`$4`); `store.go` Zeile 111 füllt ihn aus
`transactionRow.CommittedAt`; `mapper.go` `TransactionRow.CommittedAt`
trägt laut Dokumentation und `mapper_test.go`
(`TestTransactionRowCarriesSourceCommittedAt`) den realen
Quell-Commit-Zeitpunkt aus dem Domänenobjekt, nicht die Mapper-Instanzzeit.
Die neue Formulierung („`committed_at` trägt den realen
Quell-Commit-Zeitpunkt aus dem WAL … der Wert unten misst damit
tatsächlich den Abstand zwischen Quelländerung und CDC-Verfügbarkeit.
Offen ist allein der Name …") deckt sich mit diesem Code-Stand. Die
DEFAULT-Klausel in `schema.yaml`/`queries.go`-Kommentar bleibt konsistent
dazu: DEFAULT bleibt als Absicherung außerhalb des Anwendungspfads
bestehen, greift im regulären Pfad nicht. Fachlich korrekt.

**3. Vollständigkeit der Slice-Referenz-Entfernung**
`grep -rn "slice-[0-9]\{3\}" internal/ cmd/ tools/schema/` und
`grep -rn "review-slice\|architect-review-slice" internal/ cmd/ tools/schema/`
— beide leer (Exit 1, keine Treffer). Kein Rest einer
Slice-/Review-Report-Referenz in Code oder Schema-Kommentaren geprüft. Die
laut Aufgabenstellung ausgenommenen Pfade (`docs/plan/*`,
`harness/README.md`, `AGENTS.md`, `.harness/skills/*`) waren nicht
Gegenstand dieser Prüfung.

**4. Kommentar-Klassen nach Hard Rule 3.7 (Stichprobe)**
4 Blöcke gelesen: `heartbeat.go` (`NewHeartbeat`-Dok), `queries.go`
(`UpsertHeartbeat`/`UpsertHeartbeatFault`), `log_test.go` (`fakeLog`),
`healthcheck_test.go` (`captureStderr`) — alle vier Indikativ über den
Ist-Zustand, keine Konjunktiv-Aussage über eine verworfene Alternative,
kein Verweis auf abwesenden Text, kein mitten im Satz abbrechender Rest.
`slog_internal_test.go` enthält ein „bliebe" (Konjunktiv II) — das
beschreibt aber eine Mutationstest-Hypothese („würde der Code diese
Zeile nicht mehr durchreichen, bliebe der Test trotzdem grün"), keine
verworfene Alternative im Sinne der Regel; unauffällig.

**5. `make gates` und `make test`**
Beide real ausgeführt (nicht angenommen): `make gates` — alle vier
inneren Gates grün (`baseline-verify`, `docs-check`, `commit-traceability`,
`a-check`, je 0 Befunde). `make test` — alle Pakete `ok`, keine
Fehlschläge.

---

## Findings

### F-1 — Unwrapped Kommentarzeile nach Zusammenführung zweier Sätze

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/bootstrap/wiring.go:306`
- `befund`: Beim Entfernen der Chronik-Klausel („, slice-012 §6-Risiko).")
  wurden zwei Kommentarzeilen zu einer zusammengeführt
  (`// Capture-Persist-ACK-Pfads (LH-QA-REL-001.a). heartbeatCtx
  endet spätestens mit stream.Run; das`, 108 Zeichen), während die
  umgebenden Zeilen im selben Absatz auf ~70–75 Zeichen umbrochen sind.
  Rein optisch, keine semantische Auswirkung.
- `verifizierbar`: nein — kein Line-Length-Gate im Repo (`.golangci*`
  nicht vorhanden, `gofmt` erzwingt keine Breite).
- `klasse`: „Kommentar-Zusammenführung ohne Rewrap"

### F-2 — Vorbestehende Konjunktiv-Formulierung außerhalb des geprüften Diffs

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/schema/nacharbeit-observability.sql:40`
- `befund`: „eine reine Lese-View auf bereits persistierten Zustand
  konnte den Lauf-Zustand des CDC-Prozesses selbst nicht bezeugen" nutzt
  Präteritum/Modalverb-Vergangenheit für eine zeitlose Aussage. Die Zeile
  liegt **außerhalb** des von `199396a`/`ca61cfb` geänderten Bereichs
  (beide Commits ließen diese Zeile unverändert) und ist damit nicht
  Gegenstand dieses Nachtrags-Reviews — hier nur zur Kenntnis, falls eine
  künftige Bereinigungsrunde dieselbe Datei erneut anfasst.
- `verifizierbar`: ja — Diff-Vergleich zeigt, dass die Zeile in keinem der
  beiden geprüften Commits berührt wurde.
- `klasse`: „Präteritum für zeitlose Aussage" (INFO, kein Steering-Loop-Zähler-Kandidat für diesen Lauf)

Kein HIGH-, kein MEDIUM-Finding.

## Negativbefunde

- geprüft, ohne Befund: `git show 199396a` — 0 Nicht-Kommentar-Zeilen (eigener Zeilen-Scan, nicht nur Verifier-Aussage übernommen)
- geprüft, ohne Befund: `git show ca61cfb` — 0 Nicht-Kommentar-Zeilen (eigener Zeilen-Scan)
- geprüft, ohne Befund: fachlicher Inhalt der korrigierten `committed_at`-Aussage in `tools/schema/nacharbeit-observability.sql` gegen `queries.go`/`store.go`/`mapper.go` — deckungsgleich
- geprüft, ohne Befund: `grep -rn "slice-[0-9]\{3\}" internal/ cmd/ tools/schema/` — leer
- geprüft, ohne Befund: `grep -rn "review-slice\|architect-review-slice" internal/ cmd/ tools/schema/` — leer
- geprüft, ohne Befund: Stichprobe von 5 Kommentarblöcken auf Konjunktiv-über-verworfene-Alternative / abwesenden Text — keiner betroffen
- geprüft, ohne Befund: Commit-Traceability beider Commits — `LH-FA-ADM-004` im Betreff, kein `SPEC-*`/`ARC-*` im Betreff
- geprüft, ohne Befund: `make gates` — grün (real gelaufen)
- geprüft, ohne Befund: `make test` — grün (real gelaufen)
- geprüft, ohne Befund: beide Commits bereits auf `main` (`git branch --contains`), kein offener PR-Zustand

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Kommentar-Zusammenführung ohne Rewrap · Präteritum für zeitlose Aussage

## Verdikt

**Merge-blockierend:** entfällt — beide Commits liegen bereits auf `main`
(retrospektive Prüfung, ausgelöst durch VF-2, kein offener PR). Kein
gefundenes HIGH oder MEDIUM, das eine Rückführung erzwingen würde.

**Übergabe:** F-1 und F-2 gehen als Kenntnisnahme an den Implementer/Planner
für die slice-018-Closure (§7-Notiz „was ging anders als geplant" ist der
passende Ort für den Prozess-Hinweis aus VF-2, nicht dieser Report). Beide
Findings sind isoliert LOW/INFO — kein Rollen-Konflikt, keine
Konflikt-Sequenz nach Modul 8 nötig. Dieser Report selbst ist Lauf-Beleg
und wird über Läufe hinweg nicht wieder gelesen; die Finding-Klassen gehen
bei Bedarf in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler. Der Report ersetzt keine Verifikation — die
DoD-/Spec-Konformität von slice-018 selbst ist bereits in
`docs/reviews/verify-slice-018.md` geprüft; dieser Lauf deckt ausschließlich
die von VF-2 benannte Lücke (Review-Abdeckung von `ca61cfb`).
