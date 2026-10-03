# Verifikations-Report: slice-pin-stale-alle-digest-pins — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich, Konformität mit
[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Festlegung 1 bis 6) und Plan-vs-Code-Diff. Review-Artefakt:
[`review-slice-pin-stale-alle-digest-pins.md`](review-slice-pin-stale-alle-digest-pins.md)
(Commit `9650f4da`; 0 HIGH, F-1 MEDIUM, F-2 LOW, F-3/F-4 INFO; Fixrunde `66fdc6e2`).
Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md).

**Gegenstand:** Slice-Plan `slice-pin-stale-alle-digest-pins`, Diff `3eb60de6..HEAD`
(`8abd3d81`), vier Commits: Implementierung `90bb3f88`, Review `9650f4da`, Fixrunde
`66fdc6e2`, Review-Haken `8abd3d81`. Dieser Lauf ändert weder Code noch Plan noch
Doku; er schreibt nur diesen Report. Alle Mutationen liefen an Kopien im Scratchpad
(`PROG=<Kopie>`, Kopie per `sed … > Kopie` bzw. `awk`, nie `-i`, nie auf eine
Repo-Datei); `git status --short` war vor dem Report leer.

**Fixrunde und Produktionscode:** `git show --stat 66fdc6e2` nennt genau drei Dateien
(Plan, Sensor-Vertrag, Tabellentest); weder `tools/harness/pin-stale-all.sh` noch
`tools/harness/lib-pin-compare.sh` sind darin. Der Produktionscode des Sensors ist seit
`90bb3f88` unverändert; die Fixrunde ist voll mitgeprüft (Test und Vertrag gelesen,
Test gefahren, Mutation d unten).

## 1. Eigene Sensor-Belege (ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Gedruckte Zeile |
|---|---|---|
| `make gates` | 0 | `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`; `generated-sync: OK`; `commit-traceability: OK — 5 Commit(s)`; `handbuch-public-doc-check: keine interne Kennung in 3 Nutzerdokumenten`; `meldungscodes-check: 97 Codes …`; a-check `gesamt: 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1634 Datei(en) geprüft, 0 Befund(e)` |
| `make test` | 0 | alle Pakete `ok` (Race-Detektor) |
| `make fmt-check` | 0 | `334 Go-Dateien geprüft, alle formatiert` |
| `make doc-immutable RANGE=3eb60de6..HEAD` | 0 | `0 Befund(e)` |
| `make doc-commits RANGE=3eb60de6..HEAD` | 0 | `0 Befund(e)` |
| `make commit-traceability` | 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make kommentar-kennungen DIFF=3eb60de6` | 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | 0 | `suchlauf-nachmessen: 6 Zeilen stimmen` |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |
| `make test-pin-stale-all` | 0 | `run-pin-stale-all-tests: alle 47 Prüfungen bestanden` |

## 2. DoD — Verdikt je Zeile

Am Plan gezählt: 5 Zeilen `[x]`, 5 Zeilen `[ ]`. Ich setze keine Häkchen.

| DoD-Zeile | Stand im Plan | Verifier-Befund |
|---|---|---|
| Liefer-Punkt 1 Sensor | `[x]` | **bestätigt**: Aufzählung per `git grep` mit dem Muster der ADR, Wertemenge dedupliziert, Index-Digest, `:latest` für Pins ohne Tag, fail-open, Exit 1/2/0, leerer Gegenstand Exit 2. Vergleich in `lib-pin-compare.sh` geteilt, nicht dupliziert. Keine Vollreferenz als Literal (siehe §5). Refactor-Gleichheit gefahren (§4). |
| Liefer-Punkt 2 Tabellentest und Mutationsbeleg | `[x]` | **bestätigt**: 47 Prüfungen grün; die drei benannten Mutationen des Plans (a, b, c) von mir nachgefahren, alle rot (§3). |
| Liefer-Punkt 3 Verdrahtung und realer Lauf | `[ ]` | **teilweise belegt, offen**: Make-Ziele, README-Zeilen, Vertragsdatei, zehnter Schritt, Kopf- und Jobname-Kommentare vorhanden und gelesen (§6); `make pin-stale-all` Exit 0 von mir gemessen (§5). Offen und dem Hauptlauf nach dieser Verifikation zugewiesen: der Post-Push-`workflow_dispatch`-Lauf von `upstream-drift.yml` (zehn Schritte, `success`; [`AGENTS.md`](../../AGENTS.md) §3.10). Er ist nicht belegt und wird hier nur benannt. |
| `make gates`, `make docs-check` | `[x]` | **bestätigt** (§1). |
| Review, Report liegt vor | `[x]` | **bestätigt**: Report vorhanden; F-1 durch Tabellenfall und Gegenfall behoben (Mutation d belegt, dass der Fall trägt). |
| Doku-Update | `[x]` | **bestätigt**: `harness/README.md` (Zeilen zum Sensor, zum Tabellentest und zu `upstream-drift.yml`), `docs/maintainer/releasing.md` §5 und Version 1.14, Zeile der Versionshistorie. |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | `[ ]` | korrekt offen (Planner-Arbeit, Slice liegt in `in-progress/`). |

## 3. Mutationen (Instanz: `tools/harness/run-pin-stale-all-tests.sh`, Prüfling per `PROG=<Kopie im Scratchpad>`, je eine Stelle)

| # | Zusage | Stelle der Mutation | Gesehene Farbe |
|---|---|---|---|
| a | `docs/`-Ausschluss (Festlegung 1) | Pathspec `':!docs'` aus dem `git grep`-Aufruf entfernt | **rot**, Exit 1, vier `FEHLER`-Zeilen (Fall „docs/ und .harness/ ausgenommen“ Exit 2 statt 0, Referenz aus `docs/` gemeldet, Zusammenfassung, Meldung „leerer Gegenstand unter docs/“) |
| b | leerer Gegenstand ist Exit 2 (Festlegung 4) | `exit 2` hinter der Meldung „leerer Gegenstand“ auf `exit 0` | **rot**, Exit 1, zwei Zeilen („keine Referenz im Baum“ und „einzige Referenz unter docs/“ Exit 0 statt 2) |
| c | `:latest` für Pins ohne Tag (Festlegung 3) | `target="$image:latest"` auf `target="$image"` | **rot**, Exit 1, vier Zeilen |
| d | Zeitlimit je Registry-Aufruf (Fixrunde F-1) | `limit=(timeout "$PIN_COMPARE_TIMEOUT")` auf `limit=()` in `lib-pin-compare.sh` | **rot**, Exit 1, vier Zeilen (Exit 1 statt 2, `DRIFT` statt `UNBESTIMMT`, Zusammenfassung, Lauf dauerte 6 s) |
| e | Deduplizierung (Festlegung 1) | `if [ -z "${first_at[$ref]+x}" ]` auf `if true` | **rot**, Exit 1, zwei Zeilen (Eintrag mit Zähler `+1 weitere`, ein Registry-Aufruf je Referenz) |
| f | Exit bei `DRIFT` ist 1 (Festlegung 4) | `exit 1` hinter `n_drift -gt 0` auf `exit 2` | **rot**, Exit 1, drei Zeilen (Digest weicht ab, Einzelplattform-Digest, `DRIFT` neben `UNBESTIMMT`) |

Sechs Mutationen an sechs Stellen, alle rot; der Mutationsbeleg des Plans (drei Stellen
plus Fixrunde) ist damit reproduziert und um e, f erweitert. **Hergeleitet, nicht
gefahren** ([`AGENTS.md`](../../AGENTS.md) §3.12): der `.harness/`-Ausschluss einzeln
(Mutation a entfernt nur `docs`), die Pin-ohne-Tag-Behandlung bei Registry-Host, die
Fehlerzweige „Lesefehler von `git grep`“ (Exit ≠ 0 und ≠ 1) und `UNBESTIMMT` ohne
`DRIFT` (Exit 2 am Ende des Skripts). Die Instanz ist ein Stub-`docker`; eine Mutation
an der echten Registry-Antwort ist nirgends erprobt.

## 4. Refactor von `pin-stale.sh` (Ausgabe und Exit unverändert)

Skript-Diff `3eb60de6..HEAD` für `tools/harness/pin-stale.sh` gelesen: nur der
Block ab dem Registry-Aufruf (zwölf Zeilen) ist durch den Aufruf von
`pin_compare_digest` ersetzt, `exit $?` reicht den Rückgabewert durch; Eingabe-Parsing
(Aufruf-Fehler, Variable fehlt, Wert ohne Digest, Vergleichs-Ziel) steht unverändert davor.
Die Funktion setzt `PIN_COMPARE_TIMEOUT` nur, wenn die Variable gesetzt ist; `pin-stale.sh`
setzt sie nicht (unbegrenzt wie vorher; wer die Variable exportiert, ändert das Verhalten
bewusst).

**Stub-Vergleich, selbst gefahren:** Fassung `git show 3eb60de6:tools/harness/pin-stale.sh`
gegen die Fassung am `HEAD`, Stub-`docker` im Scratchpad, neun Fälle (OK, DRIFT,
Registry-Ausfall, Wert ohne Digest, Variable fehlt, Vergleichs-Tag als drittes Argument
mit OK und mit DRIFT, Pin ohne Tag ohne Argument-Tag mit Host-Namen, Pflichtargument
fehlt): Ausgabe und Exit je Fall **byte-gleich** (`SAME` in allen neun).
`tools/harness/image-stale.sh` ist im Diff nicht enthalten (unverändert).

**Reale Läufe:** nach Ablauf des Docker-Hub-Limits `make image-stale`,
`pin-stale-race`, `-pgtest`, `-dmigrate`, `-acheck`, `-dcheck`: alle Exit 0 mit
Zeilen der Form `OK … == sha256:…`. Einen Vorher/Nachher-Lauf an der echten Registry
habe ich nicht gefahren (der Parent-Stand wurde nicht ausgecheckt); die Gleichheit trägt
der Stub-Vergleich.

## 5. Sensor gegen ADR-0146 (Festlegung 1 bis 6) und realer Lauf

| Festlegung | Befund |
|---|---|
| 1 `git grep`, dedupliziert, keine Fundort-Liste | `pin-stale-all.sh` Zeile 22/24: Muster wörtlich das der ADR, Pathspecs `':!docs' ':!.harness'`; Dedup über Assoziativ-Array (Mutation e rot). **erfüllt** |
| 2 Index-Digest | `lib-pin-compare.sh`: `docker buildx imagetools inspect … --format '{{.Manifest.Digest}}'`. Der Pin von PostgreSQL 17 in `.github/workflows/e2e.yml:80` ist Index-Digest `b0f9560a…`; die Kommentare in `e2e.yml` (Zeilen 50, 77) und `Makefile` (Zeile 193) nennen `docker buildx imagetools inspect`. **erfüllt** |
| 3 Pins ohne Tag gegen `:latest` | Zeilen 53 bis 56; Mutation c rot; im realen Lauf `aquasec/trivy:latest`, `ghcr.io/pt9912/d-migrate:latest`, `ghcr.io/pt9912/a-check:latest`. **erfüllt** |
| 4 fail-open, Exit 1/2/0, leer ist Exit 2 | je Referenz `UNBESTIMMT` ohne Abbruch (Tabellenfall und reale 429-Läufe); Mutationen b, f rot. **erfüllt** |
| 5 ein Target, ein Schritt, Name identisch | `make pin-stale-all` in `Makefile` Zeile 91, in `harness/README.md` Zeile 138 und im Schritt von `upstream-drift.yml`; `pin-stale.sh` wird wiederverwendet über `lib-pin-compare.sh`. **erfüllt** |
| 6 Tag-Wechsel nicht Gegenstand | Vertrag Grenze 2; das Skript prüft keinen Major. **erfüllt** |

**Klausel (keine Vollreferenz als Literal):** `git grep -n -E '@sha256:[0-9a-f]{64}' -- . ':!docs' ':!.harness'`
liefert 48 Zeilen; verschiedene Referenzen (`-ohE … | sort -u | wc -l`) **15**, Dateien
(`-lE`) **20** — gleich den Messwerten der ADR. Jede der 15 ist ein echter Pin (die
Zeilen des realen Laufs nennen sie). Weder Skript noch Tabellentest noch Vertrag noch
README tragen eine Vollreferenz: die Treffer in diesen Dateien sind Muster- und
Prosa-Fragmente; der Tabellentest erzeugt Digests zur Laufzeit.

**Realer Lauf `make pin-stale-all`:**

1. Erster Lauf, Exit **2**, `real 0m17,161s`, Schlusszeile
   `pin-stale-all: 15 Referenzen — 6 OK, 0 DRIFT, 9 UNBESTIMMT`; Ursache Docker-Hub-Abruflimit
   (`429 Too Many Requests`, am Messhost bekannt; `docker buildx imagetools inspect golang:1.27`
   zeigte den Statuscode). **UNBESTIMMT, kein Beleg für oder gegen Drift.**
2. Nach Ablauf des Limits erneut, Exit **0**, `real 0m34,461s`, Schlusszeile
   `pin-stale-all: 15 Referenzen — 15 OK, 0 DRIFT, 0 UNBESTIMMT`; unter anderem
   `.github/workflows/e2e.yml:80 (postgres:17-alpine) == sha256:b0f9560a…` und
   `compose.yaml:49 +2 weitere (nats:2-alpine) == sha256:ac8f88a6…`.

**Gegenprobe zur Reichweite des Musters (Frage 11), eigener `git grep`:**
(a) `sha256:<64 Hex>` ohne vorangehendes `@`: Treffer nur `d-check.mk:7` (`DCHECK_DIGEST`,
von P7 `make pin-stale-dcheck` gelesen, zur Make-Zeit aus zwei Variablen zusammengesetzt),
Messkommentare `-> sha256:…` in fünf Dockerfiles (keine Referenzen) und Fixture-Werte in
`tools/harness/run-dockerfile-from-tests.sh`. (b) Referenz mit Registry-Port vor dem Namen:
kein Treffer im Baum außerhalb des Tabellentests. (c) Großbuchstaben im Image-Namen:
kein Treffer (die Treffer meiner groß/klein-unempfindlichen Probe waren Prosa-Wörter).
(d) andere Hash-Algorithmen (`sha512`): kein Treffer. **Befund:** am Stand gibt es
**keine** Referenz der Form `<image>@sha256:`, die der Sensor nicht erfasst; die einzige
Digest-Quelle außerhalb des Musters ist `DCHECK_DIGEST`, und sie hat eine eigene Achse.

## 6. Verdrahtung, Vertrag, §3.6/§3.8/§3.13

- `GATE_CHECKS`: `git grep GATE_CHECKS` zeigt keinen Eintrag für `pin-stale-all`
  (Zuweisungen nur in `baseline.mk`, `coverage.mk`, `doc-gate.mk`, `generated-sync.mk`,
  `sdk.mk`); `make gates` Exit 0 ohne Registry-Zugriff des Sensors: **advisory, unberührt**.
- `upstream-drift.yml`: zehnter Schritt `P10 — … (make pin-stale-all)` mit `if: always()` und
  `run: make pin-stale-all`; im Diff keine neue `uses:`-Zeile (`git diff … -- .github | grep uses:` leer);
  `timeout-minutes: 15` unverändert; Kopf-Kommentar und Jobname auf den Umfang gezogen.
  §3.10: **nicht** am Runner gelesen.
- Vertrag `harness/sensors/pin-stale-all.md`, Grenzen 1 bis 8 gelesen: 1, 2, 7, 8 folgen der
  ADR; 3 (Port) ist am Tabellenfall belegt, die Antwort einer realen Registry bleibt
  ausdrücklich *hergeleitet*; 4 deckt die Klausel; 5 (Abruflimit, Exit 2) ist von mir
  reproduziert; 6 (15 × 60 s = 15 min, *abgeleitet*, nicht unter dem Job-Limit) ist
  arithmetisch wahr und ehrlich als Obergrenze benannt. **Wahr, vollständig**; eine
  Ergänzung siehe V-1.
- `docs/maintainer/releasing.md`: Version 1.14 und Historienzeile vorhanden, §5 nennt
  neun benannte Achsen plus zehnte; kein Verstoß gegen §3.5 (kein Accepted-Dokument
  inhaltlich geändert: `git diff --stat` enthält keine ADR).
- §3.13: `git grep -n -i "neun-achsen|P1.P9|neun achsen"` im lebenden Baum (ohne Reviews,
  `done/`, Register, Baseline, ADRs): einziger Treffer ist die Plan-Datei selbst
  (Suchlauf-Zeilen); `make suchlauf-nachmessen`: 6 Zeilen stimmen. Weitere Fundstellen
  der Zählung („neun“ in Prosa) stehen nur in den nachgezogenen Trägern und sind dort
  korrekt.
- §3.6/§3.8: keine Schwellensenkung, keine Gate-Änderung, keine SDK-/Beispiel-Dateien im
  Diff (`git diff --stat 3eb60de6 HEAD -- sdks examples compose.yaml Dockerfile` leer).

## 7. Plan-vs-Code-Diff

Alle Zeilen der Tabelle in §3 des Plans sind geliefert: `pin-stale-all.sh`, `lib-pin-compare.sh`,
Refactor `pin-stale.sh`, `run-pin-stale-all-tests.sh`, `Makefile`, `upstream-drift.yml`,
`harness/README.md`, `harness/sensors/pin-stale-all.md`, `releasing.md`. Abweichung, vom
Plan benannt: die Make-Ziele liegen im `Makefile` (kein `*.mk`). Nicht geliefert, korrekt
Closure-Arbeit: die Zustandsänderung des Register-Eintrags
`pin-ohne-inventar-eintrag-driftet-unsichtbar` (Zustand dort bereits „weiter offen, Träger
[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)“).

## 8. Verdikt

**Bestanden**, mit einer Bedingung, die der Hauptlauf nach dieser Verifikation erfüllt:

1. **Bedingung:** realer `workflow_dispatch`-Lauf von `upstream-drift.yml` nach dem Push,
   gelesen (zehn Schritte, Gesamtlauf `success`, [`AGENTS.md`](../../AGENTS.md) §3.10) —
   bis dahin bleibt Liefer-Punkt 3 offen.

Keine HIGH-, keine MEDIUM-Verletzung. Die Fixrunde ändert keinen Produktionscode.

## 9. Findings und offene Punkte für den Planner

- **V-1 (MEDIUM, Risiko für die Bedingung):** neun der 15 Referenzen liegen bei Docker Hub
  (`postgres` zweifach, `golang` zweifach, `nats`, `eclipse-temurin` zweifach, `python`,
  `aquasec/trivy`). Am Messhost führte das anonyme Abruflimit bereits einmal zu
  `9 UNBESTIMMT`, Exit 2. Ob die gehosteten Runner dasselbe Limit treffen, ist nicht gemessen
  (Vertrag Grenze 5 benennt es). Der erste `workflow_dispatch`-Lauf kann deshalb am
  Schritt P10 rot enden, ohne dass ein Pin driftet. §6 des Plans führt dieses Risiko nur
  als „erster Lauf rot durch Drift“; ein Ausgang muss die 429-Ursache mitnennen, und der
  Re-Evaluierungs-Trigger von
  [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  (vier rote Nachtläufe in Folge) würde sie als Drift missdeuten. Keine Änderung am Code
  verlangt.
- **V-2 (LOW):** die Ausgabe unterscheidet 429 nicht von „nicht erreichbar“ (der Sensor druckt
  nur `Registry nicht erreichbar`; `2>/dev/null` verwirft den Statuscode). Wer einen roten
  Nachtlauf liest, sieht die Ursache nicht. Kandidat für eine Beobachtung, kein Mangel
  gegenüber der ADR.
- **V-3 (INFO):** Mutationen sind an einer Stub-Instanz gefahren; kein Lauf belegt das Verhalten
  gegenüber einer realen Registry für Port-Adressen (Grenze 3, *hergeleitet*).

**Vorschläge für Register und §6-Ausgänge der Closure (Planner entscheidet):**

- `pin-ohne-inventar-eintrag-driftet-unsichtbar`: Zustand nach dem ersten gelesenen grünen
  Nachtlauf; heute belegt: Sensor deckt die 15 Referenzen am Stand (15 OK, §5), Gegenprobe
  ohne Lücke außerhalb von `DCHECK_DIGEST`.
- `nicht-blockierender-workflow-alarmmuedigkeit`: V-1 ist ein weiteres Evidenzstück (ein roter
  Schritt ohne Drift-Ursache), erst nach dem Runner-Lauf zu zählen.
- `gate-prueft-existenz-nicht-passung`: passt als Querbezug — der Sensor meldet Digest-Gleichheit
  gegen die Registry, nicht ob der Pin zur Verwendung passt (Grenze 1, Mehrfachkopien).
- §6-Ausgänge: erster Lauf rot — offen bis Runner-Lauf, mit V-1; Laufzeit — entfallen,
  gemessen 34,461 s bei 15 Referenzen und 17,161 s bei 429; Muster zu eng/zu weit — entfallen
  am Stand (Gegenprobe §5), akzeptiertes Negativ bleibt; Verhaltensänderung `pin-stale.sh` —
  entfallen (Stub-Vergleich byte-gleich, reale Läufe Exit 0); Mutationsbeleg zu schmal —
  entfallen für die benannten Stellen (sechs rot), Rest *hergeleitet*; Workflow lokal
  unbeweisbar — weiter offen bis zum Lauf
  ([`BEO-PGC/github-actions-unverifizierbar-lokal`](../plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md));
  dauerhafter Einzelplattform-Pin — entfallen am Stand (der PostgreSQL-17-Pin ist Index-Digest, §5).
- Release-Folge: keine (kein Produkt-, SDK- oder Image-Pin berührt).
