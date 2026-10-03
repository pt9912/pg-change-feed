# Review-Report: slice-pin-digests-aktualisieren-2026-10 — 2026-10-03

**Review-Art:** Code/Konfiguration — geprüft gegen Plan, Entscheidungen und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `pin-digests-aktualisieren-2026-10` (wellenlos), Diff-Range `8aeec6f2..HEAD` (`HEAD` = `0fc3bdd7`,
8 Commits, 17 Dateien). Hebungs-Commits `cc93ca3b` (P1), `2d871258` (P3), `f69ce37e` (P4), `f3b904ec` (PostgreSQL-17-Pin
der e2e-Matrix), `79e40cbf` (P5), `7b5865f0` (P6), `168db825` (P7), Plan-Nachzug `0fc3bdd7`. Berührt sind `Dockerfile`,
`Makefile`, `a-check.mk`, `d-check.mk`, `compose.yaml`, `examples/Dockerfile`, `examples/compose.yaml`,
`.github/workflows/e2e.yml`, neun Skripte unter `tools/`, der Plan. Kein Produktionscode.

**Skill:** `.harness/skills/reviewer.md`. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03.

**Eingangs-Kontext:** Slice-Plan; [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7;
[`LH-QA-POR-001`](../../spec/lastenheft.md); [`SPEC-012`](../../spec/pflichtenheft.md); `AGENTS.md` §3.1, §3.6, §3.9,
§3.10, §3.12, §3.13; `harness/conventions.md`.

## Eigene Messungen

Alle Exit-Codes ungefiltert gesichert (Ausgabe in Datei, `$?` direkt).

- **Registry-Auflösbarkeit (Schwerpunkt 1):** `make image-stale`, `make pin-stale-race`, `-pgtest`, `-dmigrate`,
  `-acheck`, `-dcheck` je Exit 0, je `OK … == sha256:<voller Wert>` mit den Werten aus dem Diff (P1 `8a5910f3…`, P3
  `e0174e51…`, P4 `77f58511…`, P5 `af9d3eb3…`, P6 `e8208764…`, P7 `b4b8756b…`, dazu `OK DCHECK_IMAGE Tag-Frische
  (v0.79.0)`). `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.79.0` druckt denselben Digest wie
  `d-check.mk`.
- **PostgreSQL-17-Pin:** `docker buildx imagetools inspect postgres:17-alpine` listet
  `postgres:17-alpine@sha256:aa90e97ee862e558111d34cfb8b2c4bec768c2b039fb791341686928560263b3` als Eintrag
  `application/vnd.oci.image.manifest.v1+json`, Plattform `linux/amd64`, Annotation `17.11-alpine3.24`; der Index-Digest
  (`b0f9560a…`) ist nicht der gepinnte Wert. Die Form (amd64-Einzel-Manifest) stimmt mit der Aussage des Plans §7 und der
  Kommentarform in `e2e.yml` („amd64“) überein. Der PG-18-Eintrag der Matrix ist byte-gleich zu `Makefile`
  `PG_TEST_IMAGE`.
- **Träger-Gleichheit (Schwerpunkt 2, §3.13):** `git grep -ohE` je Achse außerhalb `docs/` und `.harness/`:
  `golang:1.27-alpine@sha256:` 9 Treffer, **1** verschiedener Wert (`8a5910f3…`); `golang:1.27@sha256:` 1/1;
  `postgres:18-alpine@sha256:` 8/1 (`77f58511…`); `postgres:17-alpine@sha256:` 1/1; `d-migrate@sha256:` 2/1;
  `a-check@sha256:` 1/1; `DCHECK_DIGEST` nur `d-check.mk`. Kein Mischstand. Alte Präfixe (`cf6fca66`, `b475798f`,
  `63bdc97d`, `7456ef82`, `862dfb04`, `34d3dfb5`, `3f84502b`) stehen nur in `docs/plan/adr/`, `docs/reviews/`,
  `docs/plan/planning/done/`, `observations/`, im Plan selbst (Messzeilen) und in `harness/sensors/db-adapter-coverage.md`
  Zeile 300. Lese-Urteil Zeile 300: der Satz benennt den Rot-/Grün-Beleg eines Laufs von `slice-081` („real, gepinntes
  Toolchain-Image `golang:1.27-alpine@sha256:cf6fca66…`, … Lauf `slice-081`“) samt dessen Messzahl 73,38 % — Messkontext
  eines vergangenen Laufs, kein geltender Pin; kein Nachzug nötig. Ein Treffer `golang:1.27-alpine` mit gekürztem
  Präfix ist genau diese Zeile.
- **`e2e.yml` (Schwerpunkt 3):** `git diff 8aeec6f2 HEAD -- .github` zeigt vier geänderte Stellen: zwei Datumszeilen in
  Kommentaren (`2026-09-14` → `2026-10-03`) und die zwei Matrix-`pg_test_image`-Werte; keine Struktur-, Trigger-,
  Permissions- oder Action-Änderung. Matrix bleibt auf den Majors 17/18 ([`SPEC-012`](../../spec/pflichtenheft.md) Zeile
  `PG_MAJOR_VERSIONS` 17, 18).
- **Gate-Lockerung (Schwerpunkt 4, `AGENTS.md` §3.6):** `git diff 8aeec6f2 HEAD --stat -- .d-check.yml .a-check.yml
  harness tools/harness/*check*.sh spec` leer; keine Schwelle, kein Modul, keine Regel berührt.
- **Nachmessen Plan §7 (Schwerpunkt 5):** `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-pin-digests-aktualisieren-2026-10.md`
  Exit 0, „16 Zeilen stimmen“ (Stände `e5820a030a…` und `diff`). Gemessen nachgefahren: `go version go1.27.1` an beiden
  neuen `golang`-Digests (Behauptung „kein Tag-Wechsel, `go1.27.1`“ trägt); PG-17-Handmessung und DCHECK_DIGEST-Zeile
  (oben) tragen die Plan-Zahlen; `make gates` am Endstand druckt `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle
  80%`, `d-check: 1620 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check mit dem neuen Digest) — wie im
  Plan.
- **Gefahrene Läufe (Schwerpunkt 6):** `make gates` Exit 0; `make docs-check` Exit 0 (`d-check: 1620 Datei(en) geprüft, 0
  Befund(e)`, Image `d-check@sha256:b4b8756b…` = neuer Pin); `make fmt-check` Exit 0 („334 Go-Dateien geprüft, alle
  formatiert“); `make test` Exit 0 (51 `ok`, kein `FAIL`); `make test-replication
  PG_TEST_IMAGE=postgres:17-alpine@sha256:aa90e97e…` Exit 0 (`PostgreSQL 17.11: Keepalive inmitten der Transaktion …`,
  `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%`).
- **Übernommen, nicht nachgefahren:** `make image`, `make test-store`, `make test-replication` an PostgreSQL 18,
  `make schema-validate`, `run-schema-rollout-guard-test.sh`, `make test-integration`, die Parent-Vergleichsläufe von
  `make a-check`/`make docs-check` (Befund-Gleichheit am alten Digest), `make pin-stale-baseline`/`-actions` (P8/P9).
  Diese Angaben stehen im Plan §7 als gemessen; für diesen Report sind sie **übernommen**.
- **Docker-only (Schwerpunkt 10):** im Diff keine Host-Toolchain, kein in-place-Werkzeug, keine Umleitung auf Repo-Dateien
  (nur Pin-Werte in vorhandenen Variablen/`FROM`-Zeilen). `make gates` wurde mit Exit direkt gelesen (§3.9), nicht durch
  eine Pipe. Commit-Betreffs der acht Commits nennen `ADR-0051` und `LH-QA-POR-001`, kein `SPEC-*`/`ARC-*`.
- **Eigener Verfahrensvermerk:** Das Edit-Werkzeug stand im Lauf nicht zur Verfügung. Eine Zeilennummer in diesem Report
  wurde einmal über `sed … > Scratchpad-Kopie` und anschließendes `cp` auf den Report korrigiert (kein `sed -i`); der
  Plan-Haken und die Berichtigung der Report-Kennungen laufen über Write.

## Findings

Keine HIGH-, keine MEDIUM-Findings.

### F-1 — Kommentar in `e2e.yml` nennt die Spec-Festlegung der PostgreSQL-Majors als Träger der Digests

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist), §3.12 Instanz B (Tatsachenbehauptung ohne tragenden Beleg)
- `pfad`: `.github/workflows/e2e.yml:41` (zusätzlich `:74`, `:79`)
- `befund`: Der Kommentar spricht von „den in `SPEC-012` festgelegten Digests“, und die Matrix-Zeilen tragen
  `# PostgreSQL 17 (SPEC-012): Digest via …`; [`SPEC-012`](../../spec/pflichtenheft.md) legt nur die Majors 17 und 18 fest,
  keinen Digest (`git grep -n 'sha256:' spec/pflichtenheft.md` ohne Treffer in der Zeile von
  [`SPEC-012`](../../spec/pflichtenheft.md); der Plan §Berührte Spec-Stellen sagt dasselbe). Die Zeilen sind im Diff nur im
  selben Kommentarblock berührt (Datum), die Formulierung ist Bestand.
- `verifizierbar`: ja — Lesen von `spec/pflichtenheft.md` Zeile `PG_MAJOR_VERSIONS` gegen den Kommentar
- `klasse`: Zitat nennt die falsche Stelle

### F-2 — PostgreSQL-17-Pin und Zitat `ADR-0039`: gemeldet, nicht gefixt (Schwerpunkt 8)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (fremde Träger melden), Plan §6
- `pfad`: Plan §6 „Offene Frage — Sensor-Lücke PG 17“ und „stale Zitat“; `harness/README.md` Zeile
  `make image-stale` (Zitat `ADR-0039`)
- `befund`: Die Sensor-Lücke (PG-17-Pin fehlt im Pin-Inventar von
  [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md)) und das veraltete Zitat sind benannt, die fremde Zeile
  in `harness/README.md` ist unberührt (Diff ohne `harness/`). Das Beobachtungs-Register ist nicht fortgeschrieben; der
  Plan führt das als Closure-Arbeit des Planners (§7). Keine Aktion im Review.
- `verifizierbar`: ja — `git diff 8aeec6f2 HEAD -- harness` leer
- `klasse`: Träger-Meldung

### F-3 — §3.10: Post-Push-Lauf von `e2e.yml` steht aus (Schwerpunkt 9)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.10
- `pfad`: Plan §6 Risiko „CI-Matrix `e2e.yml` PG 17/18“, §7 „Post-Push-Lauf“
- `befund`: Der Diff berührt `e2e.yml` nur in Pin-Werten und Kommentar-Datum, keine strukturelle Änderung; der Plan führt
  den realen Post-Push-Lauf beider Legs trotzdem als weiter offen und trägt ihn nicht als bestätigt ein. Konsistent mit
  der Regel. Closure bleibt an diesem Lauf gebunden.
- `verifizierbar`: nein — Post-Push-Lauf auf dem Runner nicht lokal messbar
- `klasse`: externer Lauf offen

### F-4 — Parent-Vergleichsläufe und weitere Läufe übernommen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Übernommenes als übernommen kennzeichnen)
- `pfad`: Plan §7 „Läufe am Diff“
- `befund`: Die in „Eigene Messungen“ als übernommen gelisteten Läufe sind im Plan als gemessen mit gedruckter Zeile
  geführt; die geprüften Zeilen (Gates, Docs-Check, `test-replication` PG 17, Digest-Zeilen) stimmen mit meinen Messungen
  überein, was die übrigen plausibel macht, aber nicht belegt.
- `verifizierbar`: ja — Nachfahren der genannten Läufe (Verifier)
- `klasse`: übernommene Messung

## Negativbefund (geprüft, ohne Befund)

- `Dockerfile`, `examples/Dockerfile`: ein Digest je `golang:1.27-alpine`-`FROM`, Wert gleich Makefile — ohne Befund
- `Makefile`, `a-check.mk`, `d-check.mk`: Pin-Werte vollständig und in der Registry auflösbar, Tag/Digest von d-check
  gemeinsam gehoben — ohne Befund
- `compose.yaml`, `examples/compose.yaml`: PG-18-Digest gleich `Makefile` — ohne Befund
- `.github/workflows/e2e.yml`: nur Pin-Werte und Kommentar-Datum, beide Legs konsistent (abgesehen von F-1) — ohne
  Befund zur Struktur
- `tools/` (neun Skripte, `bench-lib.sh`): je Achse derselbe Digest wie die Hauptträger, keine Host-Werkzeug-Aufrufe
  hinzugefügt — ohne Befund
- `.d-check.yml`, `.a-check.yml`, `harness/`, `spec/`, Gate-Skripte: unberührt, keine Gate-Lockerung — ohne Befund
- Handbuch/Betreiber-Oberfläche: nicht berührt (`docs/user/` nicht im Diff) — ohne Befund
- Plan `docs/plan/planning/in-progress/slice-pin-digests-aktualisieren-2026-10.md`: Messzeilen mit gedruckter Zeile,
  Ursprung gekennzeichnet (gemessen/Erwartung/übernommen), 16 Suchlauf-Zeilen stimmen — ohne Befund

## Verdikt

0 HIGH, 0 MEDIUM, 1 LOW (F-1, Bestandskommentar), 3 INFO. **Nicht merge-blockierend.** Keine Fixrunde am Implementer
nötig; F-1 kann der Planner mitnehmen. Der Haken „Review durchgeführt“ im Plan ist mit diesem Report gezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug); andere DoD-Punkte unberührt. Offen bleibt der reale Post-Push-Lauf
von `e2e.yml` (F-3).
