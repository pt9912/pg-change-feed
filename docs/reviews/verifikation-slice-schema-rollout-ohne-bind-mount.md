# Verifikations-Report: slice-schema-rollout-ohne-bind-mount — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). Frischer Kontext. Frage: „Bauen wir es richtig?“
gegen Plan, DoD und Entscheidung.

**Gegenstand:** `git diff fcbb2044 HEAD` — drei Commits: Implementer `fb269dc3`, Review `01551df9`
([Review-Report](review-slice-schema-rollout-ohne-bind-mount.md): 0 HIGH, 1 MEDIUM F-1, 4 LOW),
Fixrunde `786e5b3f`. Plan:
[`slice-schema-rollout-ohne-bind-mount`](../plan/planning/in-progress/slice-schema-rollout-ohne-bind-mount.md).
Entscheidungen: [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md),
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
[`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md) Trigger (b); Verdikt
[`architect-verdict-schema-rollout-ohne-bind-mount`](architect-verdict-schema-rollout-ohne-bind-mount.md).
Anforderung: [`LH-QA-OPS-005`](../../spec/lastenheft.md).

**Repo-Zustand:** `HEAD` = `786e5b3f`, nicht gepusht von mir. Mutationen liefen an Klonen im
Scratchpad, erzeugt mit `sed '…' Datei > Kopie` (kein `sed -i`, keine Umleitung auf eine Repo-Datei).
`git status --short` im Echtrepo war nach jedem Lauf leer — mit einer Ausnahme (§6, Eigenfehler,
per `git checkout` zurückgenommen).

## 1. Sensor-Belege (eigener Lauf, Exit-Codes ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | gedruckter Beleg |
|---|---|---|
| `make gates` | **0** | `baseline-verify` v6.13.0 OK 54 Dateien; `d-check: 1555 Datei(en) geprüft, 0 Befund(e)` (docs-check); `gesamt: 0 Befund(e)` (a-check); `commit-traceability: OK — 5 Commit(s)`; `coverage-gate: OK — Coverage 82.00 % erfüllt Schwelle 80 %`; `generated-sync: OK`; `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | **0** | Lauf 1 OK (22126 B `plan.yaml`, 7890 B `down.sql` in `.tmp/schema-rollout`, `git status --short` 0 Zeilen vor/nach); Lauf 1b OK (Exit 2, Prüfsumme unverändert); Läufe 2–4, 5 (Alt-Tag `v0.4.0`: Exit 0/0/0), 6a, 6b; Schlusszeile `OK — alle Belege real erbracht (… make-Exit 2/2 mit d-migrate-Exit 8)` |
| `make schema-validate` | 0 | `Validation passed`; warm 1,4 s (Bauzeit-Messung für Risiko „Cache-Kosten“) |
| `make image` | 0 | `:dev` war nicht geladen; ein Lauf am unmutierten Baum, `git status --short` danach leer |
| `make test-store` | **0** | `DB-Adapter-Coverage: 82.99 %`, `db-coverage: OK`; `git status --short` danach leer |
| `make test-replication` | **0** | `--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction`; `git status --short` danach leer |
| `make test-integration` | **0** | `E2E-Abdeckungstabelle unverändert`, `Lauf abgeschlossen`; die Zeile `schema-rollout: Erzeugnisse (…) in .tmp/schema-rollout` steht im Log; `git status --short` danach leer (auch `docs/user/e2e-abdeckung.md` unverändert) |
| `make example-demo-up` / `-down` | **0** / **0** | Rollout-Zeile im Log; `git status --short` danach leer |
| `make test` | **0** | |
| `make fmt-check` | 0 | 323 Go-Dateien, alle formatiert |
| `make kommentar-kennungen DIFF=fcbb2044` | 0 | kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | 0 | `15 Zeilen stimmen`, alle `OK` (u. a. `diff 0 … CURDIR[^:]*:/work`, `diff 0 … rollout.sh` ohne Mount) |
| `make doc-commits RANGE=fcbb2044..HEAD` | 0 | 0 Befunde |
| `make doc-immutable RANGE=fcbb2044..HEAD` | 0 | 0 Befunde |

Keine Verweigerung durch die Berechtigungsschicht (`make image` lief), also nichts nach
[`AGENTS.md`](../../AGENTS.md) §3.15 zu berichten.

## 2. Fixrunde `786e5b3f` — Code-Lesung und Einzelmutationen

Kein Re-Review hat stattgefunden; ich habe `tools/schema/rollout.sh` vollständig gelesen und die
Behauptungen mutiert.

**(1) F-1 geschlossen — bestätigt.** `--execute` exportiert in `$DIR/.stage.XXXXXX` (Zeile 143, 151–152) und
verschiebt `plan.yaml`/`down.sql` erst danach per `mv -f` (Zeile 153). Am Lauf-Anfang entfernt Zeile 112
nur `rollout-precheck.yaml`. Der `trap` (Zeile 54–64) löscht im Staging-Verzeichnis namentlich die zwei
Dateien (`rm -f`, danach `rmdir`), kein `rm -rf`. Fehlerpfade: Precheck-Exit ≠ 0/8 bricht nicht ab und exportiert
nichts, `--execute` läuft dennoch (wie im Vertrag „Der Precheck-Exit wird gelesen“); der Execute-Exit wird
unverändert weitergegeben (Zeile 148–150, make zeigt `Fehler 2`/`Fehler 8`); jeder Exportfehler endet mit Exit 2;
nach `mv` ist das Staging-Verzeichnis leer und wird entfernt. Dateirechte im Ziel: `-rw-r--r-- db db`
(gelesen, `ls -l`). `.stage.*` liegt unter `.tmp/` und ist durch `.tmp/` in `.gitignore` ausgenommen
(`git check-ignore` nicht eigens gefahren — hergeleitet aus dem Eintrag; `git status --short` blieb leer,
auch nach einem Lauf).

Konsistenz mit Entscheidung und Vertrag: [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
Festlegung 1 sagt „je Lauf überschrieben“ und „druckt den Pfad am Ende des erfolgreichen Laufs“; ein
Best-Effort-Export eines gescheiterten `--execute` ist dort nirgends zugesagt (`git grep` nach „best-effort“
in der ADR: 0 Treffer), sein Entfall widerspricht der ADR also nicht. Der Vertrag
([`schema-rollout.md`](../../harness/targets/schema-rollout.md) §Erzeugnisse) nennt den Entfall ausdrücklich
(„der Report eines gescheiterten `--execute` wird nicht abgelegt“) und stimmt mit dem Skript überein.

**(2) AQ-2 — bestätigt.** `build_stage guard` steht nur im Zweig `plan_exit = 8` (Zeile 128); der Vertrag
(§Voraussetzungen) sagt dasselbe. Kein Lauf mit Precheck-Exit 0 baut das Wache-Image.

**(3) Vertragstexte.** Der Bauabschnitt nennt „zur Laufzeit keinen Netzzugang“, belegt mit
`docker build --no-cache --network none --target guard` (F-3 gelöst). Die Grenzen nennen feste Image-Tags
(als hergeleitet gekennzeichnet) und den DSN in gestoppten Containern (als hergeleitet gekennzeichnet) —
beides entspricht dem Review (F-5, F-6). Zwei Randbefunde: V-1 und V-2 unten.

**Einzelmutationen** (je ein Klon, Guard-Test unverändert, gesehene Farbe):

| # | Mutation (Stelle) | Instanz | Farbe |
|---|---|---|---|
| M-mv | `mv -f … "$DIR/"` → `"$DIR/nodir/"` (`rollout.sh:153`) | Guard-Test Lauf 1 | **rot** — `make: *** [Makefile:314: schema-rollout] Fehler 2`, Test-Exit 2 |
| M-del | Löschzeile am Anfang wieder eingefügt (`rm -f "$DIR/plan.yaml" "$DIR/down.sql"` an Zeile 112) | Guard-Test Lauf 1b | **rot** — `FEHLER — Lauf 1b: plan.yaml/down.sql haben sich durch einen Lauf verändert, der vor --execute scheiterte`, Test-Exit 1 (Lauf 1 blieb grün: die Mutation trennt genau 1b) |
| M1 | Export-Pfad `/work/down.sql` → `/work/downX.sql` (`rollout.sh:152`) | Guard-Test Lauf 1 | **rot** — `Could not find the file /work/downX.sql in container …`, Fehler 2, Test-Exit 2 |
| M2 | leerer stdin für die Wache (`</dev/null`, `rollout.sh:129`) | Guard-Test Lauf 2 | **rot** — Test-Exit 2 (`--execute` ohne `--allow-destructive` endet mit d-migrate-Exit 8) |
| M3 | `-i` fehlt am Nacharbeit-`docker run` (`rollout.sh:158`) | Guard-Test Lauf 2 | **rot** — `FEHLER — Lauf 2 hat den --allow-destructive-Pfad nicht genommen`, Test-Exit 1 (Lauf 1 blieb grün: der stille No-op fällt erst in Lauf 2 auf) |

M4 (vertauschte Reihenfolge/Pfad der Nacharbeit-Dateien) habe ich nicht gefahren; der Review hat sie
gefahren und F-7 vermerkt (Aussage gilt nur für eine frische Instanz). Verallgemeinerung von den fünf
Instanzen auf „jede“ Stelle ist **hergeleitet**, nicht erprobt.

## 3. DoD- und Entscheidungs-Konformität (Plan vs. Code)

| DoD / Festlegung | Befund |
|---|---|
| LP 1: `tools/schema/Dockerfile` (+ `.dockerignore` Allow-Liste, Stufen `rollout`/`guard-build`/`guard`, `FROM scratch`, `CGO_ENABLED=0`, kein eigener Digest, Pins als `--build-arg`) | **erfüllt** (Diff gelesen) |
| LP 1: `rollout.sh` unter bash mit `pipefail`, `create`/`start`/`cp`/`rm`, `tar -x --no-same-owner`, Wache über stdin, Nacharbeit über stdin mit `-i`, `D_MIGRATE_RUN_USER` entfällt | **erfüllt**; Abweichung vom Plan-Wortlaut: `docker cp`-Export geht jetzt über Staging (Fixrunde, Plan-DoD nennt das nicht, [ADR-0142](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md) verträgt es) |
| LP 1 (a) Guard-Test sechs Läufe + Lauf-1-Zusatz (Erzeugnisse, `git status`) | **erfüllt** (§1) |
| LP 1 (b) `plan.yaml`/`down.sql` in `SCHEMA_ARTEFACT_DIR`, Pfad gedruckt, `git status` leer | **erfüllt** |
| LP 1 (c) Suchlauf „kein Bind-Mount der Rezeptur“ Soll 0 | **erfüllt** (`diff 0` OK) |
| LP 1 (d) `git diff <Parent> -- Dockerfile .dockerignore` leer | **erfüllt** (`git diff fcbb2044 HEAD --stat -- Dockerfile .dockerignore`: 0 Zeilen) |
| LP 1 (e) `schema-validate` Exit 0, `--network none`, ohne `-v` | **erfüllt** (`rollout.sh:101`, Lauf Exit 0) |
| LP 1 (f) Mutationsproben | M1–M3 von mir erprobt (§2); M4 vom Review |
| LP 2: sieben Aufrufer + Eigensicherung nachgezogen, `rollout-restore.sh`, Test, Ziel, `plan.yaml`/`down.sql` entfallen | **erfüllt** (`git ls-files | grep -c rollout-restore` = 0; `tools/schema/` ohne die zwei Dateien) |
| LP 2: Läufe `test-store`, `test-replication`, `test-integration`, `example-demo-up/-down` mit leerem `git status` | **erfüllt** (§1) |
| DoD „Doku-Update“ (`harness/README.md`, Vertrag, Makefile-Kommentar, Handbuch 1.88; `AGENTS.md` unverändert) | **erfüllt**; `AGENTS.md` unverändert, nur §3.14 nennt `schema-rollout` als Rang-Zeiger — keine Änderung nötig (melden) |
| DoD „§3.13-Suchlauf Feld im Plan“ | **erfüllt** (15 Zeilen stimmen) |
| [ADR-0142](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md) Festlegung 1–4 | **eingehalten** (Ort, `git rm`, kein Mount, Wache über stdin); Geltung „nur diese zwei Ziele“: `Makefile` zeigt die drei `:ro`-Mounts unverändert |
| Kein Schema-/Release-/Versions-Diff | **bestätigt**: Diff berührt `schema.yaml`, `nacharbeit-*.sql`, `version.md`, Package-Versionen nicht; Handbuch-Version 1.87 → 1.88 ist Doku-Version, keine Paketversion |

DoD-Häkchen setze ich nicht (Auftrag).

## 4. Träger und Reste (ganzer Baum, außerhalb ADR-0142, Register und Records)

`git grep -E 'rollout-restore|test-rollout-restore|D_MIGRATE_RUN_USER|tools/schema/(plan\.yaml|down\.sql)'`
außerhalb `docs/reviews`, `done/`, `.harness/baseline`: Treffer nur in `ADR-0084`, `ADR-0125`, `ADR-0142`
(`Accepted`, unberührbar), im Slice-Plan selbst und im Beobachtungs-Register — **kein stehender Träger** in
Makefile, `harness/`, `tools/`, `examples/`, `docs/user/`, `compose.yaml`, `.github/`.

## 5. Findings

| Nr | Kategorie | Befund |
|---|---|---|
| V-1 | INFO | Ein `SIGKILL` lässt `.stage.XXXXXX` mit den halb exportierten Dateien unter `SCHEMA_ARTEFACT_DIR` zurück (`trap` greift nicht); der Vertrag nennt in §Grenzen nur Container und DSN, nicht das Staging-Verzeichnis. Zeigt ein Betreiber die Variable auf einen versionierten Pfad, wäre das Verzeichnis nicht ausgenommen. **Hergeleitet** aus dem Skript, nicht gefahren. Blockiert nicht. |
| V-2 | INFO | `harness/README.md` (Zeile `make schema-validate`) und Vertrag belegen „Basis-Images lokal vorausgesetzt“ mit `--target guard`; `schema-validate` baut die Stufe `rollout`. Die Messung trägt die Aussage nur für die Stufe `guard`; für `rollout` **hergeleitet**. Blockiert nicht. |
| V-3 | INFO | Plan-Risiken §6 (Bauzeit, `/dev/stdin`, `-i`-No-op, Besitzer-Zuordnung, Gleichzeitigkeit, Verhaltensänderung) tragen noch „zu entscheiden bei Closure“. Meine Messungen dafür: warm 1,4 s (`schema-validate`); `/dev/stdin` trägt die Wache (Läufe 2–6, M2 färbt rot); `-i`-No-op von M3 gefangen; Besitzer der Dateien `db:db` (aufrufender Nutzer, Host nicht root). Kalter Bau und root-Runner nicht gemessen. |

## 6. CI und Eigenfehler

- Push von `786e5b3f`: Lauf `ci` (36995923719) **success**, `examples` (36995923571) **success**. Lauf `e2e`
  (36995923595) war bei Abschluss dieses Berichts **in_progress**: in beiden Legs (PostgreSQL 17 und 18) sind
  `Image bauen` und `Compose-Integrationstest (Black-Box-E2E)` — der Schritt mit dem Rollout-Pfad (F-8) —
  **success**; `DB-Adapter-Coverage — Replication-Teil` und `Replication-Tier` standen noch aus. Das Endergebnis
  des Laufs ist damit nicht belegt ([`AGENTS.md`](../../AGENTS.md) §3.10: erst nach Grün abgeschlossen).
- Der `upstream-drift`-Lauf 36993077277 (schedule, failure) ist advisory und nicht Gegenstand dieses Slice.
- Eigenfehler: Ich habe versehentlich `make doc-ci-matrix` gefahren; es überschrieb
  `docs/user/ci-matrix-abdeckung.md` (2 Zeilen). Per `git checkout --` zurückgenommen, `git status --short`
  danach leer; kein Commit enthält die Änderung.

## 7. Register-Fortschreibung für die Closure (melden, nicht geändert)

1. `BEO-PGC/test-schreibt-in-committete-datei/state.md`: Träger `rollout-restore.sh` und
   `make test-rollout-restore` entfallen; Zustand nachtragen.
2. `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (`observation.md`, `state.md`): „`ADR-0084` Trigger (b)
   eingelöst durch `ADR-0142`“; beide führen `plan.yaml`/`down.sql` als Kandidaten.
3. `BEO-PGC/schema-rollout-fremdobjekte/state.md`: verweist auf `Makefile` (`schema-rollout`-Target) und
   „Bericht (`tools/schema/plan.yaml`)“; Träger ist jetzt `tools/schema/rollout.sh`, Ort `SCHEMA_ARTEFACT_DIR`
   (der Text zählt außerdem „sechs bekannte Fremdobjekte“, vorbestehend elf).
4. `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`-Evidence und weitere Evidence-Dateien nennen die alten Pfade
   als Records; keine Änderung nötig.

## 8. Verdikt

**bestanden, mit Bedingungen.** Der Code erfüllt DoD und [ADR-0142](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md); F-1 ist geschlossen und durch zwei
Einzelmutationen (M-mv, M-del) an den richtigen Läufen rot gesehen; alle Sensoren und Aufrufer-Läufe sind grün.

Bedingungen vor der Closure:

1. Endergebnis des `e2e`-Laufs 36995923595 (beide PostgreSQL-Legs) prüfen (`gh run view`), Ausgang in Plan §6
   (F-8) nachtragen ([`AGENTS.md`](../../AGENTS.md) §3.10).
2. Die zwei Register-Nachträge aus §7 (Punkt 1 und 2) sowie Punkt 3 in der Closure.
3. Plan-Risiken §6 mit Ausgang versehen (V-3); V-1/V-2 nach Ermessen des Planners als Grenze im Vertrag ergänzen
   oder verwerfen.

Offene Punkte für den Planner: keine HIGH/MEDIUM; V-1 bis V-3 sind INFO. Änderung an `AGENTS.md`: keine nötig.
