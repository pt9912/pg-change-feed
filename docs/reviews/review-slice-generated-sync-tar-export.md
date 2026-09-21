# Review-Report: slice-generated-sync-tar-export — 2026-09-21

**Review-Art:** Code — geprüft gegen Plan (`docs/plan/planning/in-progress/slice-generated-sync-tar-export.md`)
+ `ADR-0084` + `ADR-0060` + `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Keine DoD-Verifikation — das ist Verifier-Aufgabe
(Modul 11); dieser Report prüft Maintainability/Konformität, nicht
Abnahme.

**Gegenstand:** Commit `93f5545b` (`tools(generated-sync): Bind-Mount
durch tar-Stream-Export ersetzen (ADR-0084)`) gegen Elternstand
`1cce85b5`.

**Skill:** `.harness/skills/reviewer.md` (Accepted, Schärfung 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-21

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-generated-sync-tar-export.md` (§1–§8, inkl. der bereits vom Planner geklärten Risiko-1-Entscheidung)
- `ADR-0084` (Sync-Gate für generierte Artefakte, Accepted, §Entscheidung/§Fitness-Function/§Konsequenzen vollständig gelesen)
- `ADR-0060` (gRPC-Streaming-Mechanismus, der Generator)
- `tools/harness/proto-generate.sh` (Referenzmuster der tar-Extraktion)
- `Dockerfile` Stufen `proto`/`proto-export` (aktueller Stand und Einführungs-Commit `338cfe00`, slice-104)
- `AGENTS.md` §3.1, §3.5, §3.6, §3.7, §3.9, §3.11, §4
- `harness/README.md` §Sensors, `harness/conventions.md`

---

## Findings

### F-1 — Skript-Kommentar narriert die verworfene Alternative statt nur die geltende Zusage zu nennen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Disziplin) / Reviewer-Skill-HIGH-Klasse „Kommentar trägt keine der Kommentar-Klassen"
- `pfad`: `tools/harness/generated-sync.sh:23-30`
- `befund`: Der Kopf-Kommentar beginnt mit „Frueher lief hier `docker build --target proto` gefolgt von einem manuellen `protoc`-Aufruf gegen einen … Bind-Mount. Der Mount war die Fehlerquelle auf …" — das ist wörtlich die von `AGENTS.md` §3.7 als **Falsch**-Beispiel genannte Form („die frühere Fassung prüfte nur die Länge" — beschreibt abwesenden Text). Die Abwägung (Bind-Mount vs. tar-Stream) gehört in ADR/Slice-Plan, die Historie in `git log`, nicht in den Skript-Kommentar; die geltende Zusage allein („Die tar-Stream-Extraktion braucht keinen Mount … und ist damit strukturell immun gegen [Colima `mounts: []` + `TMPDIR` außerhalb `$HOME`]") hätte gereicht (§3.7 „Richtig"-Muster). Der zweite Absatz (Zeilen 32–47, Modulpfad-/`.proto`-Cross-Check-Wegfall) ist davon **nicht** betroffen — er beschreibt eine aktuelle Grenze des Skripts und wo der Drift stattdessen auffällt, keine narrative Chronik einer entfernten Codezeile.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Klassen (Reviewer-Skill selbst benennt das als ungedeckte Klasse).
- `klasse`: „Kommentar trägt keine der Kommentar-Klassen" (Chronik/verworfene Alternative in Skript-Kommentar statt ADR/Slice-Plan)

### F-2 — Dieselbe narrative Struktur in `harness/sensors/generated-sync.md`, außerhalb des strikten §3.7-Anwendungsbereichs

- `kategorie`: INFO
- `quelle`: Maintainability (§3.7 gilt wörtlich für „Code, Konfiguration und Skripte", nicht ausdrücklich für Doku-Prosa)
- `pfad`: `harness/sensors/generated-sync.md:21-24`
- `befund`: „Bis `slice-generated-sync-tar-export` lief hier stattdessen `docker build --target proto` gefolgt von einem manuellen `protoc`-Aufruf gegen einen `docker run -v …`-Bind-Mount — der Mount scheiterte auf …" ist dieselbe Chronik-Struktur wie F-1, hier aber in einem Sensor-Dokument (Prosa, kein Inline-Kommentar) — der Rest des Repos hat kein Vorbild für dieses Muster in `harness/sensors/*.md` (Suchlauf negativ). Kein Gate fängt das, und §3.7s Wortlaut deckt Doku-Prosa nicht eindeutig; deshalb INFO statt HIGH — zur Aufmerksamkeit, falls die Klasse ein zweites Mal in Sensor-Doku auftritt (Steering-Loop-Schwelle).
- `verifizierbar`: nein
- `klasse`: „Kommentar trägt keine der Kommentar-Klassen" (Doku-Prosa-Variante, außerhalb des Kern-Anwendungsbereichs)

### F-3 — Reviewer-eigene Gegenprobe löst Risiko 4 stärker auf als der im Plan dokumentierte Ersatzbeleg

- `kategorie`: INFO
- `quelle`: Slice-Plan §6 Risiko 4 / DoD-Punkt 3 (Gegenprobe)
- `pfad`: `docs/plan/planning/in-progress/slice-generated-sync-tar-export.md:270-289`
- `befund`: Auf der Reviewer-Maschine liegt zufällig **dieselbe** Kombination vor, die dem Implementer fehlte (Colima `mounts: []`, aber `TMPDIR` zeigt auf ein Colima-eigenes Verzeichnis **innerhalb** `$HOME` (analog zur Implementer-Maschine)). Mit `TMPDIR=/tmp` (außerhalb `$HOME`, `/private/tmp`) reproduziert der **alte** Skript-Stand (`git show 1cce85b5:tools/harness/generated-sync.sh`) real den Auslöser (`Permission denied` beim Schreiben nach `/out`), und der **neue** Skript-Stand läuft unter identischem `TMPDIR=/tmp` grün (`generated-sync: OK`, `git status --porcelain` leer). Das ist ein stärkerer, direkter Beleg als der im Plan dokumentierte schwächere Ersatzbeleg (reine Code-Inspektion) — dem Verifier zur Aufnahme empfohlen, ändert aber nicht mein Verdikt: kein Hard-Rule-Verstoß, kein Fixrunde-Anlass.
- `verifizierbar`: ja — reproduzierbar mit `TMPDIR=/tmp bash <alter-Skript-Stand>` vs. `TMPDIR=/tmp bash tools/harness/generated-sync.sh` auf einer Colima-`mounts: []`-Maschine.
- `klasse`: „Beleg stärker als im Plan dokumentiert" (kein Fehlermuster, Positivbefund)

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/generated-sync.sh` — Kern-Mechanismus (kein `-v`, kein `--user`/`RUN_USER`, Extraktion in `mktemp -d "${TMPDIR:-/tmp}/…"`, nicht nach `.`) — real bestätigt per `make generated-sync` (Exit 0, `git status --porcelain` davor/danach leer) und per Mutation-Test (angehängte Zeile in `gen/cdc/stream/v1/changestream.pb.go` → `make generated-sync` real Exit 2 über `make`, Diff korrekt gemeldet, Datei restauriert).
- geprüft, ohne Befund: `ADR-0084` §Entscheidung/§Konsequenzen/§Verglichene Alternativen/§Status — unangetastet (keine ADR-Datei im Diff); die Fitness-Function-Tabelle nennt weiterhin „Dockerfile-Stufe `proto`" (jetzt stale, aber benannt statt verdeckt), diese Zeile und §Kontext (Zeilen 63/76) liegen **außerhalb** der in `AGENTS.md` §3.5 abschließend genannten unberührbaren Liste — eigene, unabhängige Lektüre der ADR bestätigt die Planner-Einschätzung „keine Supersede-ADR nötig".
- geprüft, ohne Befund: `Dockerfile` — `git diff --stat` leer, real bestätigt; kein geteiltes Stufen-Risiko (Risiko 3) eingetreten, `make proto-generate` real erneut gelaufen (Exit 0, `git status --porcelain` leer).
- geprüft, ohne Befund: Risiko-2-Behauptung (keine dynamische `.proto`-Erkennung in `proto-export` seit je) — verifiziert am Einführungs-Commit `338cfe00` (slice-104): die Stufe trug von Anfang an den hartcodierten Dateinamen `proto/cdc/stream/v1/changestream.proto` und den hartcodierten Modulpfad im `RUN`-Schritt, nie eine `find`-/`go.mod`-Ableitung.
- geprüft, ohne Befund: `harness/mk/generated-sync.mk` — Kommentar-Disziplin intakt (Zusage-/Grenze-Klasse, keine Chronik-Narration), `GENERATED_SYNC_MODULE`/`RUN_USER`-Override vollständig entfernt (Skript, Makefile-Fragment, README-Zelle — kein dangling Rest, `grep -rn "GENERATED_SYNC_MODULE\|RUN_USER"` bestätigt).
- geprüft, ohne Befund: `harness/sensors/generated-sync.md` §Ausgabe/§Overrides/§Grenze/§Sperren — Exit-Code-2-Satz über `make` real nachgemessen (GNU Make meldet für **jeden** fehlgeschlagenen Rezept-Exit-Code — 1, 2, 7 getestet — seinen eigenen Exit 2; die aktualisierte Formulierung „kommt der Exit-Code des Skripts (1 oder 2) als der Make-eigene Exit `2` an" bleibt korrekt, unabhängig von der internen Skript-Logik-Änderung).
- geprüft, ohne Befund: `harness/README.md` §Sensors (`generated-sync`-Zeile) — korrekt nachgezogen, Backtick-Parität aller vier berührten Markdown-Dateien (`harness/README.md`, `harness/sensors/generated-sync.md`, `docs/plan/planning/in-progress/slice-generated-sync-tar-export.md`) gerade.
- geprüft, ohne Befund: Träger-Nachzug-Suchlauf (`grep -rn "Bind-Mount\|bind-mount" harness/ tools/harness/`) — real wiederholt; Treffer außerhalb der vier berührten Dateien (`tools/harness/sdk-pack-*.sh`, `tools/harness/proto-generate.sh`, `tools/harness/run-integration-tests.sh`, `harness/README.md` Zeile 136 zu `make proto-generate`) betreffen andere Ziele und sind unverändert korrekt.
- geprüft, ohne Befund: `AGENTS.md` §3.11 (kein host-lokaler absoluter Pfad) — ein Suchlauf über alle fünf geänderten Dateien nach den in `.d-check.yml`s `hostpaths.prefixes` gelisteten Host-Wurzel-Mustern ohne Treffer; der vom Implementer berichtete Vorfall während der Arbeit ist im finalen Commit nicht mehr präsent.
- geprüft, ohne Befund: Out-of-Scope-Einhaltung — `git diff --name-only 1cce85b5..93f5545b` zeigt ausschließlich die fünf erwarteten Dateien; kein Umbau von `tools/harness/proto-generate.sh`, kein `Dockerfile`-Diff, keine ADR-Datei im Diff, kein Colima-/Host-Konfigurationsdateizugriff im Repo, `GATE_CHECKS`/Gate-Liste unverändert (kein neues Gate).
- geprüft, ohne Befund: Commit-Message-Traceability — Betreff nennt `ADR-0084`, keine `SPEC-*`/`ARC-*`-ID (Slice trägt keine `LH-*`-ID, konsistent mit `harness/conventions.md`); `make commit-traceability` innerhalb von `make gates` real grün.
- geprüft, ohne Befund: `make gates` — real ausgeführt (nicht gepiped, Exit-Code direkt geprüft): Exit 0, `d-check: 875 Datei(en) geprüft, 0 Befund(e)`, `a-check: gesamt: 0 Befund(e)`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kommentar trägt keine der Kommentar-Klassen (Chronik/verworfene Alternative) · Beleg stärker als im Plan dokumentiert (Positivbefund)

## Verdikt

**Merge-blockierend:** ja — ein HIGH-Finding (F-1), reguläre Fixrunde nötig.

**Übergabe:** F-1 (HIGH) geht an den Implementer zur Korrektur des
Kommentar-Kopfs in `tools/harness/generated-sync.sh` (Zeilen 23–30): die
geltende Zusage nennen, die Chronik des Bind-Mount-Mechanismus `git log`
überlassen (bereits im Slice-Plan/Commit-Message vollständig dokumentiert).
F-2 (INFO) ist eine Beobachtung für den nächsten Fall dieser Klasse in
Sensor-Doku (Steering-Loop-Zähler), kein Fixrunde-Anlass. F-3 (INFO) ist
eine Empfehlung an den Verifier, die eigene, stärkere Gegenprobe (reale
Reproduktion des Auslösers mit dem alten Skript-Stand und grüner Lauf des
neuen Stands unter identischer `TMPDIR`-außerhalb-`$HOME`-Bedingung) bei der
DoD-Bewertung von Risiko 4/DoD-Punkt 3 zu berücksichtigen — das ändert die
Kategorie dieses Reports nicht, weil es kein Fehlermuster ist.

Da eine Fixrunde nötig ist, wird die DoD-Checkbox „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan **nicht** vorab
nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde — Grenze:
nur bei 0 HIGH bzw. ohne Rückgabe-Pfeil). Der Nachzug erfolgt regulär nach
der Fixrunde (Implementer-Workflow Schritt 21).

Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).
