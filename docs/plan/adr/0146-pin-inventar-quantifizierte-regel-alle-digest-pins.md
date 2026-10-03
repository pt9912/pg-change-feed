# ADR-0146: Pin-Inventar als Regel über alle Digest-Pins statt als Liste (schärft ADR-0051 Entscheidung 7)

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-03

**Autor:** pt9912 (Architect-Rolle, Modul 8; anderer Kontext als der Planner- und Verifier-Lauf des Slice `pin-digests-aktualisieren-2026-10`, dessen Folge-ADR-Frage diese ADR beantwortet)

**Bezug:** [ADR-0051](0051-cicd-pipeline-github-actions.md) (Entscheidung 7: Pin-Inventar P1–P9 und `upstream-drift.yml`; bleibt unverändert, `AGENTS.md` §3.5), [ADR-0058](0058-testansatz-fuenf-luecken.md) (PostgreSQL-Versionsmatrix), [ADR-0083](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.9, §3.10, §3.12. Anlass: Verifikation des Slice `pin-digests-aktualisieren-2026-10` §8.5–§8.7 und die Register-Einträge `BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar` und `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`.

**Schärft:** [ADR-0051](0051-cicd-pipeline-github-actions.md) Entscheidung 7, Satz „gegen das vollständige, oben ermittelte Pin-Inventar (P1–P9)“: diese ADR legt fest, wie Vollständigkeit gelesen wird. Spec-Stelle: keine (Prozess-ADR ohne Spec-Stratum).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0051](0051-cicd-pipeline-github-actions.md) Entscheidung 7 führt das Pin-Inventar als **Liste** von neun Achsen und baut für jede Achse einen Sensor. Eine Liste ist vollständig, wenn sie geschrieben wird, und veraltet mit jedem neuen Pin. Gemessen am Arbeitsstand `28a6242a` (Befehle unten, alle von mir gefahren):

1. **Der PostgreSQL-17-Pin fehlt in der Liste.** Er steht als YAML-Matrixwert in `.github/workflows/e2e.yml`, `tools/harness/pin-stale.sh` liest nur Zeilen der Form `VAR ?= …` aus Makefile-artigen Dateien.
2. **Er ist nicht der einzige.** `git grep -ohE '[a-z0-9][a-z0-9./_-]*(:[A-Za-z0-9._-]+)?@sha256:[0-9a-f]{64}' -- . ':!docs' ':!.harness'` findet **15 verschiedene** Referenzen in 20 Dateien. Die Liste P1–P9 nennt davon: `golang:1.27-alpine`, `distroless/static-debian12:nonroot`, `golang:1.27`, `postgres:18-alpine`, `d-migrate`, `a-check`. Außerhalb stehen `postgres:17-alpine`, `nats:2-alpine` (`compose.yaml`, `examples/compose.yaml`, `tools/harness/run-notify-tests.sh`), `aquasec/trivy` (`Makefile`), `ghcr.io/astral-sh/uv:0.12.17` (`sdks/python/Dockerfile`) und die fünf Basis-Images der `sdks/`- und `examples/`-Dockerfiles (`dotnet/sdk:10.0`, `dotnet/runtime:10.0`, `eclipse-temurin:21-jdk`, `eclipse-temurin:21-jre`, `python:3.14-slim`).
3. **Diese Pins driften real.** Gemessen am 2026-10-03 mit `docker buildx imagetools inspect <tag> --format '{{.Manifest.Digest}}'` gegen den gepinnten Wert: `python:3.14-slim`, `dotnet/sdk:10.0`, `eclipse-temurin:21-jdk` und `eclipse-temurin:21-jre` weichen ab (`tools/harness/image-stale.sh` je Dockerfile gefahren: Exit 1); `nats:2-alpine` weicht ab (Digest `ac8f88a6…` statt `065e8355…`); `ghcr.io/astral-sh/uv:0.12.17` und `golang:1.27-alpine`, `postgres:18-alpine` stimmen überein. `dotnet/runtime:10.0`, `distroless`, `golang:1.27`, `d-migrate`, `a-check` und der Trivy-Vergleich gegen `:latest` (aktuell `af6acf9a…`, gepinnt `62b1e65e…`) habe ich nicht einzeln gelesen — Trivy weicht nach den gedruckten Werten ab, die übrigen sind nicht gemessen.
4. **Die Form des PostgreSQL-17-Pins passt nicht zu den übrigen.** Der gepinnte Wert `aa90e97e…` ist **kein** Index-Digest, sondern ein Manifest einer Einzelplattform: `docker buildx imagetools inspect postgres:17-alpine --raw | grep -c aa90e97ee862` liefert 1 (er ist Mitglied des aktuellen Index `b0f9560a…`), und ein Digest-Vergleich gegen den Index-Digest würde ihn **immer** als Drift melden (hergeleitet, an keiner Instanz gefahren). Der PostgreSQL-18-Pin ist ein Index-Digest (`… --raw | grep -c 77f585` liefert 0, `imagetools inspect postgres:18-alpine` druckt `77f58511…`). Die Kommentare „`docker manifest inspect … (amd64)`“ in `e2e.yml` und `Makefile` benennen den Weg, der den Einzelplattform-Digest liefert, nicht den Index-Digest, den alle Sensoren vergleichen.
5. **Der nächtliche Lauf war rot, bevor und ohne dass es jemand las.** `gh run list --workflow upstream-drift.yml --limit 100 --json conclusion,event` (2026-10-03): 15 Läufe, 14 `schedule` mit `failure`, 1 `workflow_dispatch` mit `success`. Die Ursachen der 14 roten Läufe habe ich nicht aus den Logs gelesen (nicht gemessen).

## Entscheidung

**Frage (a) und (b) — eine Regel statt zehn Achsen.** Das Pin-Inventar ist ab jetzt neben P1–P9 eine **quantifizierte Regel**: *jede* Referenz der Form `<image>[:<tag>]@sha256:<digest>` in einer getrackten Datei außerhalb von `docs/` und `.harness/` wird nächtlich gegen den Registry-Digest geprüft. Das ist **P10**. P1–P9 bleiben unverändert bestehen; Überschneidungen (derselbe Pin in P10 und in einer Achse davor) sind erlaubt und melden denselben Drift zweimal im selben roten Lauf — das ist billiger als eine Ausnahmeliste.

Festlegungen für P10:

1. **Aufzählung per `git grep`**, Muster wie in Kontext 2, Wertemenge dedupliziert. Eine Datei, ein YAML-Matrixwert, ein Shell-Default und ein `compose.yaml`-Default sind gleichwertige Fundorte; es gibt keine Liste von Fundorten, die zu pflegen wäre.
2. **Vergleichsmaßstab ist der Index-Digest** (`docker buildx imagetools inspect <ref ohne @digest> --format '{{.Manifest.Digest}}'`), wie bei `image-stale.sh` und `pin-stale.sh`. Ein Pin trägt den Index-Digest; ein Einzelplattform-Digest ist keine zulässige Pin-Form dieses Repos. Der PostgreSQL-17-Pin in `e2e.yml` wird mit der Sensor-Einführung auf den Index-Digest gehoben (die Kommentare dort und im `Makefile` nennen dann `docker buildx imagetools inspect`).
3. **Pins ohne Tag** (`d-migrate`, `a-check`, `aquasec/trivy`) werden gegen `:latest` verglichen, wie P5/P6 es heute tun.
4. **Fail-open wie [ADR-0051](0051-cicd-pipeline-github-actions.md) Entscheidung 7:** Registry nicht erreichbar heißt `UNBESTIMMT` je Referenz, nicht Rot der übrigen; ein Fund heißt `DRIFT`. Exit 1 bei mindestens einem `DRIFT`, 2 bei `UNBESTIMMT` ohne `DRIFT`, 0 sonst. Ein leerer Gegenstand (keine Referenz gefunden) ist Exit 2 — leer ist nicht bestanden.
5. **Ein Target, ein Schritt.** `make pin-stale-all` (Arbeitsname; der Implementer darf ihn umbenennen, solange er in `harness/README.md` §Sensors und im Schritt von `.github/workflows/upstream-drift.yml` identisch steht) ruft ein Skript unter `tools/harness/` auf, das `pin-stale.sh`s Digest-Vergleich wiederverwendet statt zu duplizieren. Der Workflow bekommt einen zehnten Schritt mit `if: always()`.
6. **Ein Tag-Wechsel ist nicht Gegenstand von P10.** Der Sensor meldet Digest-Drift am gepinnten Tag, nicht, dass ein neuerer Major existiert (`dotnet/sdk:10.0` → `11.0`, `eclipse-temurin:21` → `25`). Das ist die bestehende Grenze von `make image-stale` (nur der `<MAJOR>.<MINOR>`-Fall) und bleibt eine benannte Grenze, kein Vorgang.

**Frage (c) — Verweis statt Wert.** Wo eine Entscheidung, ein Sensor-Vertrag oder die Nachschlage-Doku einen Pin nennt, nennt sie die **Variable** (`PG_TEST_IMAGE`, `TOOLCHAIN_IMAGE`, …) oder den **Tag** — nie ein Digest-Präfix. Diese ADR trägt keinen Digest-Präfix als Soll-Wert. Die Präfixe in der Tabelle von [ADR-0051](0051-cicd-pipeline-github-actions.md) bleiben als Stand des 2026-09-13 stehen (Accepted, unberührbar); die tragende Quelle ist die Variable, und P10 liest sie aus dem Baum, nicht aus einer Tabelle. Eine Zahl dieser ADR (15 Referenzen, 20 Dateien) ist ein Messwert am genannten Stand und kein Soll.

**Frage (d) — Alarmmüdigkeit bleibt Beobachtung.** Diese ADR führt keine Pflicht zur Drift-Auflösung ein und baut keinen Alarm-Mechanismus. Begründung: die Ursache des Rauschens ist hier nicht die Sensor-Form, sondern dass ein Drift kein Eigentümer hat; das ist ein Zustand der Betriebsdisziplin, den der Register-Eintrag `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` bereits führt, und der Slice `pin-digests-aktualisieren-2026-10` hat ihn mit der Hebung real aufgelöst. Eine Pflicht ohne Leser wäre eine Formpflicht auf Prosa. **Eine Abhängigkeit gilt aber:** P10 färbt den Lauf am Tag seiner Einführung rot, solange die vier bis fünf in Kontext 3 gemessenen Drifts stehen. Die Reihenfolge ist deshalb Teil dieser Entscheidung: **erst die Hebung der gedrifteten Pins, dann der Sensor** — der erste Lauf mit P10 soll grün sein, sonst fügt die Einführung dem Register-Eintrag einen weiteren roten Lauf hinzu.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | kein Aufwand | der PostgreSQL-17-Pin und zehn weitere Referenzen bleiben unbeobachtet; vier der gemessenen Drifts blieben unsichtbar (Kontext 3) |
| B — Achsen P10, P11, … je fehlender Pin (YAML-Leseform für PG 17, Dockerfile-Läufe für die Basis-Images, je eine für `nats`, `trivy`, `uv`) | je Achse einfach | die Liste veraltet beim nächsten neuen Pin wieder, dieselbe Fehlerklasse; mindestens fünf neue Achsen und Make-Targets |
| C — nur den PG-17-Pin über eine Makefile-Variable als einzige Quelle führen (`e2e.yml` liest sie per Schritt in `$GITHUB_ENV`) | löst die Doppelquelle | strukturelle Änderung des Workflows (§3.10: realer Post-Push-Lauf nötig); löst die übrigen zehn Referenzen nicht |
| **D — eine Regel über alle Digest-Pins, ein Sensor (gewählt)** | deckt künftige Pins ohne Eintrag; kein Eintrag je Fundort; nutzt den vorhandenen Digest-Vergleich; ein Target, ein Workflow-Schritt | meldet Überschneidungen mit P1–P9 doppelt; Sensor liest die Form, nicht den Sinn (siehe Grenzen); der erste Lauf ist rot, wenn die Reihenfolge nicht eingehalten wird |

## Konsequenzen

- Positiv: das Inventar ist kein Pflegeobjekt mehr; ein neuer Pin an irgendeiner Stelle des Baums ist ab dem nächsten Nachtlauf beobachtet, ohne dass jemand ihn in eine Liste einträgt. Die Fehlerklasse von `BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar` ist damit strukturell geschlossen, nicht für einen Pin.
- Positiv: auch eine vergessene Kopie bei einer Hebung (derselbe Pin in Skript-Default und `Makefile`) fällt auf, weil ihr alter Digest als eigene Referenz driftet.
- Negativ: ein Pin, den die Muster-Form nicht trifft (Digest in einer Variablen, die erst zur Laufzeit zusammengesetzt wird, ein Pin ohne `@sha256:` wie ein bewegliches Tag), wird nicht gesehen. Das ist **akzeptiertes Negativ** und keine eigene Folgepflicht: ein Pin ohne Digest ist nach `AGENTS.md` §3.8 und Modul 14 ohnehin kein Pin.
- Negativ: ein Einzelplattform-Digest wird vom Vergleich als Drift gemeldet, solange er so gepinnt ist; das ist gewollt und der Grund für Festlegung 2.
- Folgepflicht (Reihenfolge): **Slice 1** `pin-digests-aktualisieren-2026-10-b` — die in Kontext 3 gemessenen Drifts heben (`nats:2-alpine` in drei Dateien, `dotnet/sdk` in zwei, `dotnet/runtime`, `eclipse-temurin` jdk/jre, `python:3.14-slim`, `aquasec/trivy`), PostgreSQL 17 auf den Index-Digest; die Hebungen verändern Bau-Images der SDKs und Beispiele, das Gate ist der jeweilige Bau (`make sdk-pack-*`, `make examples-*`) und die E2E-Matrix (§3.10: realer Post-Push-Lauf); Modus der bestehenden Slices dieser Art. **Slice 2** `pin-stale-alle-digest-pins` — Skript, Tabellentest, Make-Target, `upstream-drift.yml`-Schritt, `harness/README.md`-Zeile und Sensor-Vertrag unter `harness/sensors/`; Einführung erst nach Slice 1. Der Tabellentest deckt je Zweig von Festlegung 1 bis 4 einen Fall (Treffer, kein Treffer, Tag, kein Tag, Registry nicht erreichbar, leerer Gegenstand).
- Folgepflicht (Bestand): `docs/plan/planning/observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar` bekommt den Träger dieser ADR; die Alarmmüdigkeits-Beobachtung bleibt ohne Änderung.

## Fitness Function (falls maschinell prüfbar)

Der Sensor existiert noch nicht; alle Zeilen sind **erwartet**, nicht erprobt.

| Tooling | Regel | Make-Target |
|---|---|---|
| Tabellentest gegen das neue Skript (erwartet; Muster der anderen `test-*-check`-Läufe) | je Zweig von Festlegung 1–4 ein Fall: Referenz in `.yml`, Shell-Default und `compose.yaml` gefunden; tagloser Pin gegen `:latest`; `docs/` und `.harness/` ausgenommen; Registry-Ausfall `UNBESTIMMT`; leerer Gegenstand Exit 2 | `make test-pin-stale-all` (Arbeitsname) |
| `git grep` am Arbeitsstand (hergeleitet aus Kontext 2: die Menge ist am Stand `28a6242a` gemessen, nicht am Stand nach Slice 1) | die Referenz-Menge ändert sich mit Slice 1; eine Mutation (ein Pin ohne Digest) ist an **keiner** Stelle erprobt | — |
| `make pin-stale-all` im Nachtlauf (erwartet) | Exit 0 nach Slice 1; ein gemeinsam gehobener Pin mit vergessener Kopie → `DRIFT` für die Kopie — **hergeleitet**, an keiner Instanz gefahren | `make pin-stale-all` |

**Grenzen (benannt).** Der Sensor liest die Form `<ref>@sha256:<digest>`, nicht den Sinn: ob ein Pin gehoben werden *soll* (Major-Wechsel, Kompatibilität), entscheidet ein Mensch; ein Digest-Vergleich meldet Veränderung, nicht Gefahr. Er prüft nicht, ob mehrere Kopien desselben Pins **untereinander** gleich sind — nur, dass jede gegen die Registry stimmt.

## Re-Evaluierungs-Trigger

Wenn P10 in vier aufeinanderfolgenden Nachtläufen rot bleibt, ohne dass ein Hebe-Slice folgt (Kippen von Frage (d): dann ist die Alarmmüdigkeit nicht mehr Beobachtung, sondern ein Betriebsfehler, der eine eigene Entscheidung verlangt), oder wenn ein Pin gefunden wird, den die Muster-Form nicht trifft und der real driftet.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Accepted | Architect-Zug zur Folge-ADR-Frage der Verifikation des Slice `pin-digests-aktualisieren-2026-10` §8.5 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
