# Verifikations-Report: slice-pin-digests-aktualisieren-2026-10 — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
([`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7) + Plan-vs-Code-Diff + Gates.
Review-Artefakt:
[`review-slice-pin-digests-aktualisieren-2026-10.md`](review-slice-pin-digests-aktualisieren-2026-10.md)
(Commit `3af01b63`; 0 HIGH/0 MEDIUM, F-1 LOW, F-2 bis F-4 INFO). Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md).

**Gegenstand:** Slice-Plan `slice-pin-digests-aktualisieren-2026-10` (wellenlos, Bezug
[`LH-QA-POR-001`](../../spec/lastenheft.md)), Diff-Range `8aeec6f2..HEAD` (`3af01b63`): Implementer-Commits
`cc93ca3b` (P1), `2d871258` (P3), `f69ce37e` (P4), `f3b904ec` (PostgreSQL-17-Pin), `79e40cbf` (P5), `7b5865f0` (P6),
`168db825` (P7), `0fc3bdd7` (Plan-Nachzug), Review `3af01b63`. Dieser Lauf ändert weder Code noch Plan noch
Doku; er schreibt nur diesen Report. Es wurden keine DoD-Häkchen gesetzt und nicht gepusht.

**Methodenvermerk (Scratchpad-Altlast):** das Scratchpad enthielt `*.exit`-/`ALLDONE`-Dateien früherer Läufe
unter denselben Namen wie meine Ablage (`integ.exit`, `store.exit`, …). Ich habe Exit-Werte deshalb nur gewertet,
wenn die Datei jünger war als der Start meines Laufs (`stat`-Zeitstempel 14:44 bis 15:03) bzw. mit `-nt` gewartet.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make image-stale` | Exit 0 | `OK golang:1.27-alpine == sha256:8a5910f3…3414`; `OK gcr.io/distroless/static-debian12:nonroot == sha256:afa5c872…` |
| `make pin-stale-race` | Exit 0 | `OK TOOLCHAIN_RACE_IMAGE (golang:1.27) == sha256:e0174e51…b190` |
| `make pin-stale-pgtest` | Exit 0 | `OK PG_TEST_IMAGE (postgres:18-alpine) == sha256:77f58511…1873` |
| `make pin-stale-dmigrate` | Exit 0 | `OK D_MIGRATE_IMAGE (…d-migrate:latest) == sha256:af9d3eb3…4212` |
| `make pin-stale-acheck` | Exit 0 | `OK A_CHECK_IMAGE (…a-check:latest) == sha256:e8208764…5da1` |
| `make pin-stale-dcheck` | Exit 0 | `OK DCHECK_DIGEST (…d-check:v0.79.0) == sha256:b4b8756b…`; `OK DCHECK_IMAGE Tag-Frische (v0.79.0) == neuester Release` |
| `make pin-stale-baseline` | Exit 0 | `OK Kurs-Baseline v6.13.0 == neuester Release` |
| `make pin-stale-actions` | Exit 0 | acht `OK`-Zeilen: Tag-Mutation und Tag-Frische je für `actions/checkout@v7.0.1`, `docker/setup-buildx-action@v4.4.1`, `docker/login-action@v4.6.0`, `astral-sh/setup-uv@v10.2.0` |
| PG-17-Pin von Hand | **bestätigt** | `docker buildx imagetools inspect postgres:17-alpine --format '{{.Manifest.Digest}}'` druckt den Index-Digest `sha256:b0f9560a…2b24` (nicht der Pin); der Einzel-Manifest-Eintrag `linux/amd64` ist `sha256:aa90e97e…63b3` (`application/vnd.oci.image.manifest.v1+json`, Annotation `17.11-alpine3.24` im Review) = Pin in `e2e.yml`. `docker manifest inspect postgres:17-alpine@sha256:aa90e97e…` endet mit „manifest verification failed“ (Eigenart des Befehls bei Kind-Digest, kein Befund); `imagetools inspect …@sha256:aa90e97e…` löst auf. Der Plan-Satz „Einzel-/Index-Form“ trägt |
| `make image` | Exit 0 | Lauf bis `#15 DONE`, `harness/image-hash.txt` gesetzt (lokal, nicht committet) |
| `make test` | Exit 0 | 51 `ok`-Zeilen, 0 `FAIL` |
| `make test-store` | Exit 0 | `ok …/postgresstorage 14.297s`; `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%` (gemergt store+replication) |
| `make test-replication` (PG 18) | Exit 0 | `PostgreSQL 18.6: Keepalive inmitten der Transaktion nach 11079 von 400000 Änderung…`; `db-coverage: OK … 83.03%` |
| `make test-replication` mit `PG_TEST_IMAGE=postgres:17-alpine@sha256:aa90e97e…` | Exit 0 | `PostgreSQL 17.11: Keepalive inmitten der Transaktion nach 11076 von 400000 Änderun…`; `db-coverage: OK … 83.03%` |
| `make schema-validate` | Exit 0 | `✓ Validation passed`, `Validation passed: 0 warning(s)` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | Exit 0 | `run-schema-rollout-guard-test: OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag v0.4.0, Negativ-Abbruch …, make-Exit 2/2 mit d-migrate-Exit 8)` (Zeile `rolloutguard: unbekannter Blocker DropColumn …` ist der erwartete Negativ-Beleg) |
| `make test-integration` | Exit 0 | 0 Zeilen `FAIL` im Log; Schlusszeilen `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand` und `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 54 Bash-Zeilen`; nach dem Lauf `docker ps -a` ohne Test-Container (nur das fremde `gitea`); Echtrepo `git status --short` leer |
| `make gates` | Exit 0 | `baseline-verify: v6.13.0 OK`; `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`; `d-check: 1621 Datei(en) geprüft, 0 Befund(e)` (zweimal); `commit-traceability: OK — 5 Commit(s)…`; `generated-sync: OK — … byte-gleich`; `gesamt: 0 Befund(e)` (a-check) |
| `make docs-check` | Exit 0 | `d-check: 1621 Datei(en) geprüft, 0 Befund(e)`, Image `d-check@sha256:b4b8756b…` |
| `make fmt-check` | Exit 0 | `fmt-check: 334 Go-Dateien geprüft, alle formatiert` |
| `make doc-immutable RANGE=8aeec6f2..HEAD` | Exit 0 | `d-check: 1621 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=8aeec6f2..HEAD` | Exit 0 | `0 Befund(e)` |
| `make commit-traceability` | Exit 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make suchlauf-nachmessen PLAN=…` | Exit 0 | `suchlauf-nachmessen: 16 Zeilen stimmen` (acht Parent-Stände, acht `diff`; `diff` 10/5/1/2/1/1/0/0) |
| `make doc-trace` | Exit 0 | `80 Anforderung(en), 0 Waise(n).` |
| Parent-Vergleich `8aeec6f2` (Klon im Scratchpad) | siehe §5 | `make a-check` Exit 0 `gesamt: 0 Befund(e)` (alter Digest `34d3dfb5…`); `make docs-check` Exit 0 `1620 Datei(en), 0 Befund(e)` (alter Digest `3f84502b…`) |

**Nicht gefahren:** nichts aus der Auftragsliste ausgelassen. Alle Läufe endeten mit Exit 0.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt: 6 `[x]`-Zeilen (Liefer-Punkte 1–3, `make gates`, Review, Doku-Update), 5 `[ ]`-Zeilen (Closure-Notiz,
Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen, dazu die Zeile Closure; korrekt offen — der Slice liegt in
`in-progress/`).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Messung P1–P9 (`[x]`) | **bestätigt, mit Einschränkung** | Die Parent-Drift-Zeilen in Plan §7 sind am Parent nicht mehr reproduzierbar (die Drift ist behoben); ihre Plausibilität tragen die alten Pins im Parent-Diff (`cf6fca66…`, `b475798f…`, `63bdc97d…`, `862dfb04…`, `34d3dfb5…`, `7456ef82…`) und die neuen Digests, die ich in der Registry als aktuell fand (§1). Die Zeilen sind als gemessen mit Parent-Kennung geführt. Der Plan-Satz „Blick in das Log des Laufs 37112891025 ist nicht erfolgt“ ist ehrlich als Erwartung gehalten |
| 2 | Hebung und Nachzug (`[x]`) | **bestätigt** | Je Achse ein Commit (7 Commits), Betreffs nennen [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), kein `SPEC-*`/`ARC-*`; Träger-Gleichheit §3; `diff`-Suchlauf zeigt keinen alten Digest in einem lebenden Träger |
| 3 | Belege (`[x]`) | **bestätigt** | `make image`, `make test-store`, `make test-replication` an PG 18.6 und PG 17.11 (Pin aus `e2e.yml`), `make test-integration`: alle von mir selbst gefahren, Exit 0 (§1). Vergleich gegen den Parent dort, wo ein Rot zu vergleichen gewesen wäre (a-check, docs-check): von mir wiederholt (§5) |
| 4 | `make gates`/`make test`/`make docs-check` (`[x]`) | **bestätigt** | eigene Läufe; die DoD-Zahlen `82.40%` und `gesamt: 0 Befund(e)` decken sich. `d-check: 1620` im Plan ist der Stand vor dem Review-Report; am HEAD sind es `1621` (neue Datei), kein Widerspruch |
| 5 | Review durchgeführt (`[x]`) | **bestätigt** | Report liegt vor; seine als „übernommen“ gekennzeichneten Läufe (F-4) habe ich nachgefahren und bestätigt (§1) |
| 6 | Doku-Update: kein öffentlicher Vertrag berührt (`[x]`) | **bestätigt** | `git diff 8aeec6f2 HEAD --stat` zeigt weder `spec/` noch `docs/user/` noch `README.md` noch `AGENTS.md`; kein Digest in diesen Trägern (§3) |
| 7 | Closure-Notiz mit Lerneintrag (`[ ]`) | **korrekt offen** | Plan §7 trägt „(mit der Closure)“ |
| 8 | Beobachtungs-Register (`[ ]`) | **korrekt offen** | kein Eintrag angelegt; Vorschläge §8 |
| 9 | Risiken §6 mit Ausgang (`[ ]`) | **korrekt offen** | alle Risiko-Zeilen tragen „offen bis Closure“ bzw. „weiter offen“; Vorschläge §8 |
| 10 | Drei Paarungen (`[ ]`) | **korrekt offen** | von der Slice-Closure selbst zu tragen |

Kein `[x]` ohne Beleg.

## 3. Träger-Gleichheit je Achse (`git grep`, §3.13)

Messung: `git grep -h -o -E '(golang:1\.27(-alpine)?|postgres:1[78]-alpine|d-migrate:latest|a-check:latest|d-check:v[0-9.]+)@sha256:[0-9a-f]+' -- . ':!docs' ':!.harness/baseline'`, Auszählung `sort | uniq -c` am HEAD.

| Achse | Treffer, verschiedene Werte | Wert |
|---|---|---|
| `golang:1.27-alpine` (P1) | 9 Zeilen, 1 Wert | `8a5910f3…` — `Dockerfile`, `Makefile`, `examples/Dockerfile`, sechs `tools/harness/*.sh` (`ci-matrix-abdeckung`, `lib-github-api`, `run-integration-tests`, `run-notify-tests`, `run-replication-tests`, `run-store-tests`), kein Mischstand |
| `golang:1.27` (P3) | 1, 1 | `e0174e51…` (`Makefile`) |
| `postgres:18-alpine` (P4) | 8, 1 | `77f58511…` — `Makefile`, `compose.yaml`, `examples/compose.yaml`, `e2e.yml`, `tools/bench-lib.sh`, `run-replication-tests.sh`, `run-schema-rollout-guard-test.sh`, `run-store-tests.sh` |
| `postgres:17-alpine` | 1, 1 | `aa90e97e…` (`e2e.yml`) |
| `d-migrate` (P5) | `Makefile`, `tools/bench-lib.sh`, 1 Wert | `af9d3eb3…` |
| `a-check` (P6) | `a-check.mk`, 1 Wert | `e8208764…` |
| `d-check` (P7) | `d-check.mk` `DCHECK_IMAGE` `v0.79.0` + `DCHECK_DIGEST`, 1 Wert | `b4b8756b…` |

Die Zeile „9 Treffer“ für `golang:1.27-alpine` in der Auszählung enthält zusätzlich einen Treffer `golang:1.27-alpine@sha256:cf6fca66` mit gekürztem Präfix: er steht in `harness/sensors/db-adapter-coverage.md` Zeile 300 (Messkontext eines vergangenen Laufs, siehe Review-Lese-Urteil) und in `docs/plan/adr/` (Accepted-Records). Alte Digests in lebenden Trägern: **keine**. Verbleibende alte Präfixe: ausschließlich `docs/plan/adr/0051`, `0071`, `0098`, `0128` (Accepted, unberührbar), die Plan-Messzeilen, `docs/reviews/`, `done/`, `observations/` und `db-adapter-coverage.md` Z. 300.

## 4. Plan-vs-Code-Diff (`git diff 8aeec6f2 HEAD`)

18 Dateien, davon 16 Pin-Träger und zwei Dokumente (Plan, Review). Plan §3 nennt Dockerfile, Makefile, `examples/Dockerfile`, fünf Skripte mit `golang`-Default, `compose.yaml`, `examples/compose.yaml`, `e2e.yml`, `tools/bench-lib.sh`, `run-schema-rollout-guard-test.sh`, `a-check.mk`, `d-check.mk` — der Diff trifft genau diese Menge, keine Zusatzdatei, kein Produktionscode, keine Konfiguration (`.d-check.yml`: nicht berührt, kein Anlass eingetreten).

`e2e.yml` (Auftrag 7): vier geänderte Stellen — zwei Kommentar-Datumszeilen (`2026-09-14` → `2026-10-03`), zwei Matrix-`pg_test_image`-Werte; keine Struktur, kein Trigger, keine Permission, keine Action. Matrix = Majors 17/18 (nach [`SPEC-012`](../../spec/pflichtenheft.md) `PG_MAJOR_VERSIONS`), PG-18-Eintrag byte-gleich `Makefile` `PG_TEST_IMAGE` (Z. 201), `compose.yaml` und `examples/compose.yaml`.

**F-1 bestätigt (Auftrag 7):** `e2e.yml` Z. 41 schreibt „zwei Matrix-Legs über die in [`SPEC-012`](../../spec/pflichtenheft.md) festgelegten Digests“; [`SPEC-012`](../../spec/pflichtenheft.md) (Zeile `PG_MAJOR_VERSIONS`, `17, 18`) legt Majors fest, keinen Digest (`git grep -n 'sha256:' spec/pflichtenheft.md` trifft die Zeile nicht). Z. 74 und 79 („PostgreSQL 17/18 (Kennung der Spec-Festlegung der Majors): Digest via …“) sind weniger falsch: die Klammer bindet an den Major, der Satz danach an den Digest — mehrdeutig, aber lesbar. **Vorschlag für den Nachzug (Planner/Implementer, eine Kommentarzeile, kein Pin berührt):** Z. 41 → „zwei Matrix-Legs über die in [`SPEC-012`](../../spec/pflichtenheft.md) festgelegten Major-Versionen (17, 18); die Digests pinnt dieser Workflow (postgres:17-alpine, postgres:18-alpine)“. Reine Kommentar-Korrektur, keine Zitat-Korrektur an einer `Accepted` ADR; `e2e.yml` ist kein Record.

## 5. Neue Tool-Versionen: Befund-Vergleich Parent gegen HEAD (Auftrag 5)

| Prüfwerkzeug | Parent `8aeec6f2` (Klon im Scratchpad, alter Pin) | HEAD (neuer Pin) |
|---|---|---|
| `make a-check` | Exit 0, `gesamt: 0 Befund(e)` (Image `a-check@sha256:34d3dfb5…`) | Exit 0, `gesamt: 0 Befund(e)` (`e8208764…`) |
| `make docs-check` | Exit 0, `d-check: 1620 Datei(en) geprüft, 0 Befund(e)` (`d-check@sha256:3f84502b…`, v0.77.0) | Exit 0, `1621 Datei(en)`, 0 Befund(e) (`b4b8756b…`, v0.79.0) |

Der Dateizähler differiert genau um die neue Review-Datei; keine neuen Befunde, keine Meldungsform-Änderung sichtbar. d-migrate (P5): `schema-validate` und der Rollout-Wächter-Test (sechs Läufe, Exit 0) tragen; die `nacharbeit-*.sql` bleiben tragend. Der Klon lag im Scratchpad (`git clone --no-hardlinks`), das Echtrepo blieb unberührt: `git status --short` leer nach allen Läufen.

## 6. §3.6 — Gate-Konfiguration (Auftrag 6)

`git diff 8aeec6f2 HEAD --stat -- .d-check.yml .a-check.yml harness spec 'tools/harness/*check*'` ist **leer**. Keine Schwelle, kein Modul, keine Regel, keine Spec berührt. Der Slice ändert weder Tag (`17`/`18`, `-alpine`, `golang:1.27`) noch Haupt- oder Minor-Version: `golang:1.27` meldet `go1.27.1` (Review nachgemessen), PostgreSQL bleibt 18.6 und 17.11 (eigene Zeilen oben). Entscheidungs-Konformität [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 („Digest-Hebung bleibt bewusster Commit“): erfüllt; keine neue ADR nötig — die Begründung im Plan §6 („solange die Messung keine Tag-/Hauptversions-Änderung zeigt“) trägt.

## 7. §3.10 — Post-Push-Läufe (Auftrag 8)

Der reale Lauf von `.github/workflows/e2e.yml` (beide Matrix-Legs, PG 17 und PG 18 mit den neuen Digests) und von `.github/workflows/upstream-drift.yml` (nach dem Push: P1 bis P9 grün erwartet) steht **aus**; der Plan führt das korrekt als weiter offen. Der Hauptlauf pusht und beobachtet (`gh run list --workflow e2e.yml`, `gh run list --workflow upstream-drift.yml`). Hinweis für den Beobachter: der PG-17-Pin ist ein amd64-Einzel-Manifest; der Runner läuft amd64, der Lauf trägt das ohne Anpassung — gemessen ist das nur lokal (PG 17.11 Replikationstest grün), nicht am Runner.

## 8. Vorschläge für Closure und Register (Auftrag 9, nur Vorschläge)

1. **PG-17-Sensor-Lücke.** Klasse: *Pin ohne Inventar-Eintrag driftet unsichtbar* — real belegt, nicht hergeleitet: der PG-17-Pin driftete (`7456ef82…` → `aa90e97e…`), `make pin-stale-pgtest` und `upstream-drift.yml` sahen ihn nicht (nur Variable `PG_TEST_IMAGE` = PG 18). Nächster Eintrag: neues `BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar` (Zähler 1, Beleg = dieser Slice). Nicht unter `regel-weiter-als-ihr-sensor` einsortieren, ohne dass der Planner dessen Definition gegenliest — dort geht es um eine Regel mit zu schmalem Sensor, hier um ein Inventar ohne diese Zeile. Prüfung der Eintragsform: Vorbild sind `observation.md` + `state.md` + `evidence/<slice>.md` (so führen es `nicht-blockierender-workflow-alarmmuedigkeit` und `d-migrate-nacharbeit`); der Zähler folgt aus den Dateien.
2. **`nicht-blockierender-workflow-alarmmuedigkeit`** steht bei Zähler 1 (`evidence/slice-056.md`). Der nächtliche `upstream-drift`-Lauf war rot bis zum Slice (gemessen am Plan-Anlass, Lauf-Kennung im Plan, nicht von mir wiederholt). Ob das ein zweiter Beleg ist, entscheidet der Planner gegen die Definition (Eintrag betrifft `e2e.yml`; hier ein anderer advisory-Workflow, aber dieselbe Wirkung: Rot ohne Reaktion über unbekannte Zeit); der Plan hat ihn bewusst nicht gezählt. Ich zähle ihn nicht.
3. **`d-migrate-nacharbeit`** (Zähler 8): am neuen d-migrate-Digest kein neuer Beleg; der Wächter-Test und `schema-validate` laufen grün, die Nacharbeit-Schritte bleiben tragend. Vorschlag: keine neue Evidenz-Datei, höchstens ein Satz in der Closure-Notiz („Nacharbeit trägt auch am Pin `af9d3eb3…`“).
4. **Stales Zitat in `harness/README.md`:** `make image-stale` trägt dort als Bindung `ADR-0039`; der [ADR-Index](../plan/adr/README.md) führt [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md) als `Superseded` (Paketstruktur, Nachfolger `ADR-0042`). Tragende Entscheidung für die Digest-Hebung ist [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7. Nachzug durch den Planner in einem Commit (eine Zelle, Bindung auf `ADR-0051` Entscheidung 7), kein Slice nötig.
5. **Folge-ADR-Frage an den Architect (nicht entschieden):** (a) PG-17-Pin als zehnte Achse im Pin-Inventar von [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (der Pin liegt im YAML-Matrix-Wert, `tools/harness/pin-stale.sh` liest Makefile-Variablen — die Form des Sensors ist offen); (b) das Inventar nennt Digest-Präfixe als Werte (`cf6fca66…`), die mit jeder Hebung veralten — eine Folge-ADR sollte auf die Variable verweisen, nicht auf den Wert.
6. **Out-of-Scope-Bestätigung (Auftrag 10):** gelesen per `git grep -E '^FROM .*@sha256|image: .*@sha256'` über `sdks/`, `examples/` (ohne `examples/Dockerfile`/`compose.yaml`), `tools/schema/`: gepinnt sind `mcr.microsoft.com/dotnet/sdk:10.0` (`60a2b223…`, `sdks/csharp`, `examples/csharp`), `mcr.microsoft.com/dotnet/runtime:10.0` (`e6541e52…`, `examples/csharp`, fünf Stufen), `eclipse-temurin:21-jdk` (`085eb93e…`, `sdks/kotlin`, `examples/kotlin`), `eclipse-temurin:21-jre` (`ca7551d4…`, fünf Stufen `examples/kotlin`), `python:3.14-slim` (`caaf356f…`, `sdks/python`); `tools/schema/Dockerfile` bezieht `D_MIGRATE_IMAGE` und `TOOLCHAIN_IMAGE` als Build-Argument (wird mit P1/P5 mit gehoben). Kein Sensor (`image-stale`, `pin-stale-*`) liest diese fünf Basis-Images; ihre Drift ist unbekannt und nicht gemessen. Der Plan benennt sie in §1 als „Out of Scope“ mit Verweis auf den Folge-Befund in §7; **in §7 steht dazu noch kein Eintrag** — Nachtrag durch den Planner (am besten im selben Register-Eintrag wie Punkt 1).
7. **Weitere Closure-Pflichten:** Lerneintrag — Kandidat „ein Pin außerhalb des Inventars driftet unbemerkt; die Messung von Hand ist einmalig, nicht wiederkehrend“; Risiken §6 (6 Zeilen) bekommen je einen Ausgang: Neue Prüfbefunde → *entfallen* (§5: keine); Tag-Frische ohne Digest → *entfallen* (`OK` `DCHECK_DIGEST` und Tag); Patch-Version ändert Testverhalten → *entfallen* (alle Läufe grün); CI-Matrix → *weiter offen* bis Post-Push (§7); `image-stale`/Rückfall → *entfallen* (§3: 1 Wert je Achse); Sensor-Lücke PG 17 → *weiter offen* mit BEO; stales Zitat → *eingetreten* als Nachzug; keine neue ADR → *entfallen*.

## 9. Abweichungen und Findings des Verifiers

| # | Klasse | Schwere | Befund |
|---|---|---|---|
| V-1 | Plan-Text ohne Anker | LOW | Plan §7 „Gemeldete Träger“ und §6 nennen die Out-of-Scope-Pins (fünf Basis-Images, §8 Punkt 6) zwar in §1, aber §7 hat den dort angekündigten „Folge-Befund“ nicht. Kein Sensor liest sie. Nachtrag durch den Planner |
| V-2 | Kommentar nennt falschen Träger | LOW | F-1 aus dem Review bestätigt (§4) |
| V-3 | Plan-Satz aus dem Anlass | INFO | „P9 war im Lauf rot“ ist als Erwartung formuliert und blieb ungeprüft (Log nicht gelesen); am Arbeitsstand und am HEAD `OK`. Ehrlich gekennzeichnet, kein Befund |
| V-4 | Übernommene Zahl | INFO | `d-check: 1620 Datei(en)` im Plan-DoD gilt am Implementer-Stand; am HEAD `1621` (mit Review-Datei) — kein Widerspruch, nur bei Wiederverwendung der Zahl in der Closure nicht mit „1620“ zitieren |

Keine HIGH-, keine MEDIUM-Befunde. Keine DoD-Verletzung.

## 10. Verdikt

**bestanden** — mit den folgenden Bedingungen für Closure (keine davon verlangt Code-Änderung am Pin-Bestand):

1. **Post-Push-Läufe lesen** (§3.10): `e2e.yml` beide Legs grün und `upstream-drift.yml` ohne `DRIFT` nach dem Push; Eintrag in Plan §7; bis dahin bleibt das CI-Matrix-Risiko *weiter offen*.
2. **F-1/V-2:** Kommentarzeile `e2e.yml` Z. 41 berichtigen (§4); optional Z. 74/79 mit.
3. **Stales Zitat `ADR-0039`** in `harness/README.md` nachziehen (§8.4).
4. **Register/§6:** BEO-Eintrag zur PG-17-Lücke samt Out-of-Scope-Pins (§8.1, §8.6); jedes Risiko mit Ausgang (§8.7).
5. **Architect:** Folge-ADR-Frage (§8.5).

Offene Punkte für den Planner: Punkte 1 bis 5 oben; die DoD-Haken der Closure-Zeilen setzt der Planner.
