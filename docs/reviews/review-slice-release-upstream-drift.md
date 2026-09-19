# Review-Report: release-upstream-drift — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/release-upstream-drift.md`), `ADR-0051`
Entscheidung 7 (Pin-Inventar-Tabelle P1–P9) und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `b4b45accd701f4047a04ff760c029ec1cabfd201` (Parent
`5831210552cecfc1b9deb52446fa117d10562abe`), Slice
`release-upstream-drift`, Welle `welle-release-pipeline-adr-0051`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-upstream-drift.md` (Slice-Plan,
  §1–§3, DoD und Plan-Nachzug-Tabelle)
- `docs/plan/planning/welle-release-pipeline-adr-0051.md` (Welle-Datei,
  §4/§5/§6)
- `ADR-0051` (Accepted) — Entscheidung 7, Pin-Inventar-Tabelle P1–P9,
  §Konsistenz-Prüfung gegen AGENTS §3.1/§3.6
- `docs/reviews/review-slice-release-image-scan.md` — Vorgänger-Slice
  derselben Welle, als Sorgfalts-Referenz
- `AGENTS.md` §3.1 (Docker-only), §3.6 (Gates nicht ohne ADR lockern),
  §3.7 (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9
  (Exit-Code-Disziplin), §3.10 (GH-Actions-Workflow braucht realen Lauf),
  §3.12 (Herkunft von Aussagen)
- `tools/harness/image-stale.sh` und `tools/harness/ci-matrix-abdeckung.sh`
  als etablierte Vorbild-Skripte im selben Verzeichnis
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)

---

## Findings

### F-1 — Drei der vier neuen Skripte rufen `curl` direkt auf dem Host auf, im Widerspruch zum bereits etablierten Container-Muster für GitHub-API-Abfragen in diesem Verzeichnis

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.1 (Docker-only)
- `pfad`: `tools/harness/pin-stale-baseline.sh:21`,
  `tools/harness/pin-stale-dcheck.sh:25`,
  `tools/harness/pin-stale-actions.sh:35`
- `befund`: Alle drei Zeilen rufen `curl -fsS -m 15
  https://api.github.com/...` direkt im Skript auf, ohne Docker-Kapselung.
  Dasselbe Verzeichnis trägt bereits ein Vorbild für exakt diese Aufgabe
  (GitHub-API-Abfrage aus einem `tools/harness/*.sh`-Skript):
  `tools/harness/ci-matrix-abdeckung.sh` kapselt seinen `curl`+`jq`-Aufruf
  bewusst in einem Container und begründet das explizit im Kopfkommentar:
  „Die Abfrage selbst laeuft in einem Container (AGENTS.md §3.1) — kein
  gh/curl auf dem Host noetig." Ein repo-weiter `grep -rl curl
  tools/harness/*.sh Makefile harness/mk/*.mk` liefert genau diese vier
  Treffer — `ci-matrix-abdeckung.sh` (containerisiert) und die drei neuen
  Skripte dieses Diffs (nicht containerisiert). `git ls-remote` in
  `pin-stale-actions.sh:19` ist davon nicht betroffen — bare `git` auf dem
  Host hat mit `ci-matrix-abdeckung.sh:10` (`git rev-parse
  --show-toplevel`) ein eigenes, unstrittiges Vorbild; nur die
  `curl`-Aufrufe weichen vom etablierten Muster ab. Praktisch läuft der
  Workflow auf `ubuntu-latest` trotzdem (dort ist `curl` vorinstalliert),
  aber „Host braucht nur Docker und GNU make" (`AGENTS.md` §3.1) gilt
  damit für drei von vier neuen Skripten nicht mehr wörtlich, und die
  Inkonsistenz zum eigenen, [`ADR-0105`](../plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md)-gebundenen
  Vorbild im selben Verzeichnis ist unbegründet (kein Kommentar, kein
  Plan-Eintrag erklärt die Abweichung).
- `verifizierbar`: nein — kein Sensor prüft, ob ein `tools/harness/*.sh`-Skript
  Netzwerkwerkzeuge auf dem Host statt in einem Container aufruft;
  `docs-check` prüft Markdown, nicht Skript-Inhalte.
- `klasse`: Docker-only-Verstoß (Host-`curl` statt Container-Kapselung)

### F-2 — `if: always()` allein liefert nicht die im DoD behauptete Eigenschaft „Werkzeug-/Netzausfall führt nicht zu Rot des Gesamtlaufs"

- `kategorie`: MEDIUM
- `quelle`: `ADR-0051` Entscheidung 7 / Plan §2 DoD (zweiter Punkt) /
  Maintainability
- `pfad`: `.github/workflows/upstream-drift.yml:44-73` (alle acht
  `if: always()`-Schritte) gegen
  `docs/plan/planning/in-progress/release-upstream-drift.md:69-76`
- `befund`: `if: always()` bewirkt ausschließlich, dass ein Schritt auch
  dann läuft, wenn ein vorheriger Schritt fehlgeschlagen ist — es
  unterscheidet nicht zwischen Exit 1 (echter `DRIFT`-Fund, gewünscht rot)
  und Exit 2 (`UNBESTIMMT`, Registry/API nicht erreichbar). GitHub Actions
  markiert einen Schritt bei jedem Nicht-Null-Exit als „failure" und
  damit den gesamten Job/Lauf als „failure", unabhängig vom konkreten
  Exit-Code — nur `continue-on-error: true` würde das verhindern, fehlt
  hier aber (und würde umgekehrt auch echte Funde verschlucken). Ein
  einzelner transienter Netz-Hänger bei irgendeiner der acht Achsen färbt
  den gesamten nächtlichen Lauf exakt so rot wie ein echter Fund — die
  DoD-Formulierung „Werkzeug-/Netzausfall einer Achse führt zu Skip
  dieser Achse, nicht zu Rot des Gesamtlaufs" (fast wortgleich aus
  `ADR-0051` Entscheidung 7 übernommen) trifft für die Lauf-Statusanzeige
  nicht zu; sie trifft nur auf „nachfolgende Achsen laufen trotzdem" zu.
  Das allgemeine Alarmmüdigkeits-Risiko ist in Plan §6 bereits als „weiter
  offen" benannt (`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit"),
  diese Beobachtung ist die genauere technische Fassung davon: schon ein
  einziger Ausfall reicht, nicht erst wiederholter.
- `verifizierbar`: nein — nur über einen realen Lauf mit induziertem
  Werkzeug-/Netzausfall (analog `AGENTS.md` §3.10) oder über die
  dokumentierte GitHub-Actions-Job-Konklusions-Semantik feststellbar,
  kein Gate deckt das.
- `klasse`: Beleg trägt seinen Satz nicht (fail-open-Mechanismus deckt
  nur einen Teil der behaupteten Eigenschaft)

### F-3 — `pin-stale.sh` verletzt seinen eigenen Exit-Code-Vertrag bei fehlendem Pflichtargument

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/pin-stale.sh:15-17` (`${1:?Datei fehlt}` /
  `${2:?Variablenname fehlt}`)
- `befund`: Der Kopfkommentar dokumentiert „Exit: 0 = keine Drift · 1 =
  Drift gefunden · 2 = Registry nicht erreichbar oder Variable nicht
  gefunden (Ergebnis unbestimmbar)". Real getestet (`bash
  tools/harness/pin-stale.sh` ohne Argumente) liefert aber Exit 1 mit der
  Meldung „Datei fehlt" — ein Nutzungsfehler wird also wie ein echter
  `DRIFT`-Fund kodiert, nicht wie „unbestimmbar" (2). Unschädlich für den
  gelieferten Umfang: alle verdrahteten Aufrufstellen
  (`Makefile:53-63`, `tools/harness/pin-stale-dcheck.sh:24`) übergeben
  beide Argumente immer korrekt, der Pfad ist über `make pin-stale-*`
  nicht erreichbar.
- `verifizierbar`: ja — `bash tools/harness/pin-stale.sh` ohne Argumente
  ausführen, Exit-Code prüfen.
- `klasse`: Exit-Code-Vertrag im Kommentar stimmt nicht für jeden Pfad

### F-4 — Doppelte, nicht konsolidierte Tabellenzeile für `harness/README.md §Werkzeuge` im Plan-Nachzug

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/release-upstream-drift.md:115-116`
- `befund`: Die Plan-Nachzug-Zeile (Zeile 115, „update (Plan-Nachzug)",
  acht neue Zeilen im Detail) und die bereits vor Code existierende,
  gröbere Zeile (Zeile 116, „update", sieben neue Targets +
  `upstream-drift.yml`) beschreiben dieselbe Datei-Änderung redundant
  nebeneinander, ohne dass die ältere entfernt oder verschmolzen wurde.
  Kein Widerspruch in der Aussage, nur unbereinigte Redundanz.
- `verifizierbar`: nein — reine Doku-Struktur, kein Gate prüft
  Tabellen-Redundanz innerhalb eines Plan-Dokuments.
- `klasse`: Wiederholung eines Musters, das schon zweimal LOW war (Grenzfall — erstes Auftreten dieser konkreten Klasse)

### F-5 — `pin-stale-actions.sh` deckt nur bereits korrekt gepinnte `uses:`-Zeilen; eine noch nicht §3.8-konforme Zeile bleibt für P9 unsichtbar statt als Befund zu erscheinen

- `kategorie`: INFO
- `quelle`: Maintainability / Abgrenzung zu `AGENTS.md` §3.8
- `pfad`: `tools/harness/pin-stale-actions.sh:29-30` (Eingangs-`grep`)
- `befund`: Der einleitende `grep -hoE 'uses: [A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+@[0-9a-f]{40}
  # v[0-9][^[:space:]]*'` erfasst ausschließlich bereits korrekt
  SHA-gepinnte, mit `# vX.Y.Z`-Kommentar versehene Zeilen. Eine künftige
  `uses:`-Zeile, die §3.8 verletzt (floatender Tag, SHA ohne
  Versionskommentar), erzeugt in `pin-stale-actions.sh` weder `DRIFT`
  noch `UNBESTIMMT` — sie taucht im Skript schlicht nicht auf, und der
  Lauf bleibt bei `rc=0` grün. P9 prüft also nur die **Frische** bereits
  konformer Pins, nicht die **Form**-Konformität selbst; diese bleibt
  weiterhin allein Review-Aufgabe (kein Gate deckt §3.8 heute). Reale
  Stichprobe: aktuell entsprechen alle neun `uses:`-Zeilen über
  `.github/workflows/*.yml` dem erwarteten Muster (`grep -rn "uses:"`
  real ausgeführt), das Verhalten ist also aktuell folgenlos.
- `verifizierbar`: ja, durch Lesen des Regex gegen eine hypothetische
  nicht-konforme `uses:`-Zeile.
- `klasse`: Sensor deckt nur die bereits korrekt gepinnte Form, keine Formprüfung selbst

## Negativbefunde

- geprüft, ohne Befund: alle sieben `make pin-stale-*`-Targets real
  ausgeführt (nicht nur `make -n`) — Ergebnisse decken sich mit der
  Implementer-Behauptung: P3 (`TOOLCHAIN_RACE_IMAGE`), P4
  (`PG_TEST_IMAGE`), P5 (`D_MIGRATE_IMAGE`) zeigen `DRIFT`; P6
  (`A_CHECK_IMAGE`), P7 (`DCHECK_DIGEST`+Tag-Frische), P8
  (Kurs-Baseline `v6.9.0`), P9 (drei distinct Repo@Tag-Paare) zeigen
  `OK`.
- geprüft, ohne Befund: Variablennamen/Dateipfade — `TOOLCHAIN_RACE_IMAGE`/
  `PG_TEST_IMAGE`/`D_MIGRATE_IMAGE` in `Makefile`, `A_CHECK_IMAGE` in
  `a-check.mk`, `DCHECK_IMAGE`/`DCHECK_DIGEST` in `d-check.mk` — real per
  `grep` gegen die tatsächlichen Definitionszeilen abgeglichen, exakte
  Übereinstimmung mit den Aufrufen in den vier neuen Skripten.
- geprüft, ohne Befund: `make gates` real ausgeführt, Exit-Code direkt
  (ungepiped) geprüft: `0`. Alle Segmente grün — `baseline-verify`
  (v6.9.0, 54 Dateien), `docs-check`/d-check (769 Dateien, 0 Befunde,
  volle Modul-Liste inkl. `hostpaths`/`ids`/`links`), `commit-traceability`
  (5 Commits im Range, alle mit Struktur-ID), `generated-sync`
  (Proto-Erzeugnis byte-gleich), `coverage-gate` (82,80 % ggü. Schwelle
  80 %), `a-check` (0 Befunde).
- geprüft, ohne Befund: GHCR-Pakete `pt9912/d-migrate`, `pt9912/a-check`,
  `pt9912/d-check` sind öffentlich pullbar — real nach `docker logout
  ghcr.io` gegen alle drei getestet, kein Registry-Login im Workflow
  nötig für P5–P7.
- geprüft, ohne Befund: Dedup-Logik in `pin-stale-actions.sh` — real
  bestätigt: `docker/login-action` ist zweifach in `release.yml`
  referenziert, taucht im Lauf nur einmal auf.
- geprüft, ohne Befund: YAML-Struktur von `upstream-drift.yml` — real via
  Ruby-Stdlib-`YAML.load_file` geparst, kein Fehler; `"on":` korrekt
  gequotet, `schedule` (`37 3 * * *`, 20 Minuten versetzt zu
  `image-scan.yml`s `17 3 * * *`) + `workflow_dispatch` beide vorhanden,
  `permissions: {}` auf Workflow-, `permissions: {contents: read}` auf
  Job-Ebene — konsistent mit `image-scan.yml`.
- geprüft, ohne Befund: Action-Pinning (`AGENTS.md` §3.8) —
  `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`
  identisch mit dem bereits in `ci.yml`/`e2e.yml`/`release.yml`/
  `examples.yml`/`image-scan.yml` verwendeten Pin.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) in allen
  vier neuen Skript-Kopfkommentaren, im Makefile-Kommentarblock
  (`Makefile:44-51`) und im Workflow-Kopfkommentar — durchweg indikativ,
  kein Konjunktiv über verworfene Alternativen, keine Slice-/Wellen-Zahl
  als alleiniger Begründungsanker.
- geprüft, ohne Befund: `harness/README.md` — alle acht neuen Zeilen
  vorhanden, `ADR-0051`-Links korrekt aufgelöst (durch `docs-check`
  `ids`/`links`-Module im vollen `make gates`-Lauf mitbestätigt).
- geprüft, ohne Befund: Plan-Deckung — die in
  `docs/plan/planning/in-progress/release-upstream-drift.md` §3
  gelistete Datei-Menge deckt sich 1:1 mit den im Commit tatsächlich
  geänderten Dateien; keine Scope-Erweiterung, keine fehlende Datei.
- geprüft, ohne Befund: Traceability — Commit-Message nennt `ADR-0051`,
  kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `harness/conventions.md` MR-002
  (Slice-Kennungen sind Namen) — `release-upstream-drift` ist bereits
  ein Name, keine Nummer.
- geprüft, ohne Befund: keine neue Betreiber-Oberfläche eingeführt
  (kein `CDC_*`, keine `cdc.*`-SQL-Funktion) — Handbuch-Nachzugsregeln
  (Skill-HIGH-Bullets) greifen nicht.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Docker-only-Verstoß (Host-`curl` statt
Container-Kapselung) · Beleg trägt seinen Satz nicht (fail-open-Mechanismus
deckt nur einen Teil der behaupteten Eigenschaft) · Exit-Code-Vertrag im
Kommentar stimmt nicht für jeden Pfad · Wiederholung eines Musters, das
schon zweimal LOW war · Sensor deckt nur die bereits korrekt gepinnte
Form, keine Formprüfung selbst

## Verdikt

**Merge-blockierend:** ja — F-1 ist ein Hard-Rule-Verstoß (`AGENTS.md`
§3.1) mit konkretem, in-repo-eigenem Gegenbeispiel
(`tools/harness/ci-matrix-abdeckung.sh`); die Kern-Mechanik des Slice
bleibt davon unberührt (`make gates` bleibt grün, alle sieben Targets
liefern real die behaupteten Ergebnisse), aber die Inkonsistenz zum
etablierten Muster im selben Verzeichnis braucht eine bewusste
Implementer-Korrektur oder eine benannte, begründete Ausnahme, keine
stille Duldung.

**Übergabe:** F-1 (HIGH) geht als Rückmeldung an den Implementer
(Fixrunde: die drei `curl`-Aufrufe analog `ci-matrix-abdeckung.sh`
containerisieren, oder eine explizite, im Skript-Kommentar begründete
Ausnahme setzen). Da mit F-1 ein HIGH-Finding mit Rollen-Rückgabe
vorliegt, ziehe ich die DoD-Checkbox „Review durchgeführt" im Slice-Plan
**nicht** nach (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift
nur bei 0 HIGH oder wenn alle Findings ohne Implementer-Rückgabe
weitergereicht werden — hier liegt beides nicht vor). F-2 bis F-5 gehen
als zusätzliche, nicht blockierende Beobachtungen mit in dieselbe
Rückmeldung. Die Finding-Klassen gehen bei Slice-/Welle-Closure ins
Beobachtungs-Register. Dieser Report ersetzt keine Verifikation gegen die
DoD — das bleibt Verifier-Aufgabe (Modul 11).
