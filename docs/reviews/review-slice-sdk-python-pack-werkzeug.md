# Review-Report: slice-sdk-python-pack-werkzeug — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-python-pack-werkzeug.md`),
`ADR-0107` (Accepted, Festlegung 5) und `ADR-0108` (Accepted, §Entscheidung
Festlegung 1/2/3) sowie `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `19d20c7b..HEAD` (Abschluss von
`slice-sdk-python-http-client-flaeche` bis Implementer-Commit dieses
Slice), Slice `slice-sdk-python-pack-werkzeug`, Welle
`welle-sdk-python-lh-fa-sst-009`. Vier Commits: `2184fafd` (open→next,
reiner Move, 0 Insertions/Deletions), `8bd536c1` (Verantwortlich gesetzt,
nur Slice-Datei, 2 Insertions/1 Deletion), `b43142ca` (next→in-progress,
reiner Move, 0 Insertions/Deletions), `0e42b801` (Inhalt: `harness/mk/sdk.mk`,
`tools/harness/sdk-pack-python.sh` neu, `sdks/python/Dockerfile`-Erweiterung
(Stufen `pack`/`pack-export`), `spec/pflichtenheft.md`-Träger-Nachzug,
`harness/README.md`-Werkzeugzeile).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. `AGENTS.md` §3.13 Träger-Nachzug, §3.9 Pipe-Disziplin).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-python-pack-werkzeug.md` (§1
  Ziel/Abgrenzung, §2 DoD, §3 Plan, §6 Risiken samt Ausgängen, §7
  Closure-Notiz, §8 Sub-Area/Modus)
- `docs/plan/adr/0107-python-pypi-zweites-sdk-package.md` (Accepted) —
  Festlegung 5 (Docker-only bauen, kein Gate), §Konsequenzen
  Folgepflicht 2/3
- `docs/plan/adr/0108-python-sdk-uv-statt-build-twine.md` (Accepted,
  Supersedes `ADR-0107` in einer Klausel) — §Entscheidung Festlegung 1
  (`uv build --no-sources`), Festlegung 2 (Backend unverändert), Festlegung 3
  (exakter `COPY --from=`-Wortlaut)
- `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a`, §6 Externe Verträge, §7
  Historie
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.6 (Gates nur per ADR gelockert), §3.7 (Kommentar-Disziplin), §3.9
  (Pipe-Disziplin/Exit-Code), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug), §4 (kein Träger nennt ein nicht existentes Target)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `harness/README.md` §Werkzeuge (Gate-Index/Werkzeug-Index), Vergleichszeile
  `make sdk-pack-csharp`
- `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` — Formvorbild für
  diesen Report und Strukturvorbild für den geprüften Diff selbst

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `grep -n "GATE_CHECKS" Makefile harness/mk/*.mk` real ausgeführt:
  `harness/mk/sdk.mk` trägt **kein** `GATE_CHECKS +=` (im Gegensatz zu
  `baseline.mk`, `doc-gate.mk`, `generated-sync.mk`, `coverage.mk`), sowohl
  für die bestehende `sdk-pack-csharp`- als auch die neue
  `sdk-pack-python`-Zeile — `make sdk-pack-python` ist real außerhalb der
  Gate-Kette.
- `rm -rf sdks/python/dist && make sdk-pack-python` real ausgeführt:
  Exit-Code direkt `0`. `sdks/python/dist/pgchangefeed-0.1.0-py3-none-any.whl`
  (9330 Bytes) und `sdks/python/dist/pgchangefeed-0.1.0.tar.gz`
  (11290 Bytes) existieren danach real — identisch mit den in Slice-Plan §6
  und Commit-Message (`0e42b801`) genannten Werten, kein Drift. `file`
  bestätigt „Zip archive data" für das `.whl` und „gzip compressed data" für
  das `.tar.gz`. `git status --porcelain` zeigt beide Artefakte danach
  **nicht** als untracked — `sdks/python/.gitignore`s bereits bestehende
  `dist/`-Zeile trägt real.
- **Mutations-Probe** (unabhängig durchgeführt): in
  `sdks/python/pgchangefeed/tests/test_options.py` eine Assertion auf einen
  falschen erwarteten Wert geändert (`assert options.address ==
  "MUTATED-WRONG-VALUE"`), `rm -rf sdks/python/dist && make
  sdk-pack-python` erneut real ausgeführt — Bau bricht sichtbar an der
  `RUN pytest`-Stufe ab (`1 failed, 22 passed`, Docker meldet „process
  \`/bin/sh -c pytest\` did not complete successfully: exit code: 1",
  `make` meldet `Error 1`, Gesamt-Exit `2`), **kein**
  `sdks/python/dist/`-Verzeichnis wird angelegt (die `pack`-Stufe, die
  `uv build` ausführt, liegt strukturell hinter der bereits gescheiterten
  `pytest`-Zeile derselben `build`-Stufe — Docker-Layer-Abhängigkeit, nicht
  nur Doku-Behauptung). Mutation danach exakt zurückgesetzt (Backup-Datei
  zurückbenannt), `git diff --stat`/`git status --porcelain` beide leer
  bestätigt — kein Restartefakt. Artefakte danach real neu erzeugt
  (identische Bytegrößen wie oben).
- **Digest-Reverifikation** (`ADR-0108` §Entscheidung Festlegung 3):
  `docker buildx imagetools inspect ghcr.io/astral-sh/uv:0.12.17` real
  ausgeführt — Index-Digest
  `sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc`,
  identisch zum in `ADR-0108` und `sdks/python/Dockerfile:62` genannten
  Digest. `grep -n "COPY --from=ghcr.io/astral-sh/uv"` gegen beide Dateien
  verglichen — die Zeile ist **byte-identisch** (kein Unterschied in
  Leerzeichen, Groß-/Kleinschreibung, Pfad-Endung `/bin/`).
- **Pipe-Disziplin real geprüft** (`AGENTS.md` §3.9): isolierter Test
  `bash -c 'set -o pipefail; docker run --rm --network none
  nonexistent-image-xyz-python-test | tar -x -C /tmp; echo
  "PIPE_EXIT=$?"'` — Ergebnis `PIPE_EXIT=125` (Docker-Fehlercode), **nicht**
  `0`. Bestätigt: `set -o pipefail` in `tools/harness/sdk-pack-python.sh`
  (kombiniert mit `set -euo pipefail` und explizitem `bash`-Interpreter über
  `@bash tools/harness/sdk-pack-python.sh` im Makefile-Rezept) macht einen
  `docker run`-Fehlschlag sichtbar rot, exakt wie im
  `sdk-pack-csharp.sh`-Vorbild.
- `pyproject.toml`s `[build-system]`-Block gelesen: `requires =
  ["setuptools>=68"]`, `build-backend = "setuptools.build_meta"` —
  unverändert gegenüber dem Vor-Slice-Stand, kein Diff an dieser Datei im
  gesamten Diff-Bereich (`ADR-0108` §Entscheidung Festlegung 2 bestätigt).
- `make docs-check` real ausgeführt: Exit-Code direkt `0`
  (`d-check: 831 Datei(en) geprüft, 0 Befund(e)`) — die
  `matrix`-Modul-Umgehung (kein `spec → adr`-Verweis in den neuen
  Nachzugzeilen) trägt tatsächlich.
- `grep -c "SPEC-027" spec/*.md` — drei Treffer, alle in
  `spec/pflichtenheft.md` (§1-Fließtext, §6-Tabellenzeile, §7-Historie),
  `spec/lastenheft.md`/`spec/architecture.md` je `0` — keine Kollision mit
  einer bereits vergebenen Nummer, genau eine Definition.
- `spec/pflichtenheft.md` §7-Kopfregel gelesen („kein ADR- und kein
  Slice-Verweis") — beide neuen Historie-Zeilen (`SPEC-027`,
  `LH-FA-SST-009.a`-Nachzug) tragen keine `ADR-*`/`slice-*`-Referenz.
- `SPEC-027`-Zeile gegen `SPEC-026` verglichen: Form (ID/System/Version/
  Vertrag-Datei) konsistent — beide verweisen direkt auf die
  Metadaten-Quelldatei statt auf einen eigenen §2-Eintrag, mit derselben
  Begründungsformel (kein eigener §2-Eintrag, das Package deckt bereits
  dokumentierte(n) Drahtvertrag/-verträge); `pgchangefeed`s Zeile nennt
  korrekt nur `SPEC-018` (Singular, HTTP-API), nicht zusätzlich `SPEC-020`
  (gRPC) — deckt sich mit der Beschreibung in `pyproject.toml` (HTTP-API,
  kein gRPC-Umfang) und der Nachzugzeile in §1, die für Python explizit nur
  HTTP-API nennt statt HTTP-API und gRPC-Stream wie bei C#.
- `grep -rn "internal/\|cmd/pg-change-feed" tools/harness/sdk-pack-python.sh
  harness/mk/sdk.mk` — kein Treffer, kein Import-Grenzverstoß.
- `sdks/python/Dockerfile` vollständig gelesen: `pack`-Stufe (`FROM build
  AS pack`) läuft **nach** der `RUN pip install`/`RUN pytest`-Zeile
  derselben Stufe `build` — ein roter Test verhindert strukturell, dass
  `pack` je erreicht wird (Docker-Layer-Abhängigkeit, durch die
  Mutations-Probe oben real bestätigt, nicht nur Doku-Behauptung).
- `harness/README.md` §Werkzeuge, neue `make sdk-pack-python`-Zeile gegen
  die `make sdk-pack-csharp`-Zeile verglichen: gleiche Form (Bau-Kontext,
  Stufen, Testlauf-vor-Artefakt-Zusage, `uv`-Bezugsmechanismus statt
  `pip install uv`, Export-Mechanismus, `.gitignore`-Hinweis, der Satz
  „Braucht Netz … Werkzeug statt Gate" und der Schluss `kein Gate,
  ADR-Verweis … seit slice-<name>` — dieselbe Form Ende zu Ende. Die Pipe
  im Fließtext ist korrekt als `\|` escaped (Tabellenzelle bleibt intakt).
- `git show 2184fafd --stat`/`git show b43142ca --stat` geprüft: beide
  zeigen ausschließlich den Rename, 0 Insertions/Deletions — reine Moves,
  `AGENTS.md` §3.3 sauber eingehalten; `8bd536c1` trägt ausschließlich das
  `Verantwortlich`-Feld der Slice-Datei.
- Neue/geänderte Dateien nach Slice-/Wellen-Chronik-Mustern durchsucht
  (`AGENTS.md` §3.7): Treffer in `Dockerfile`, `harness/mk/sdk.mk`,
  `tools/harness/sdk-pack-python.sh` — durchweg in der bereits im
  `sdk-pack-csharp`-Pendant etablierten Herkunfts-Anker-Form
  („(slice-<name>, ADR-<NNNN> …)"), indikativ über den aktuellen Zustand,
  keine Vorher/Nachher-Erzählung; kein Treffer in Produktionscode (dieser
  Slice ändert keine `.py`-Quelldateien).
- Commit-Traceability real geprüft (`git log --oneline 19d20c7b..HEAD`):
  alle vier Commit-Betreffe tragen die beiden Kennungen aus dem Bezug-Feld
  dieses Slice-Plans — `LH-FA-SST-009` sowie `ADR-0107` —, kein
  `SPEC-*`/`ARC-*` im Betreff
  (`tools/harness/commit-traceability.sh`-Regex `(SPEC|ARC)-[0-9]{3}` gegen
  `0e42b801`s Betreff geprüft — kein Treffer).

---

## Findings

Keine HIGH- oder MEDIUM-Findings in diesem Lauf.

### F-1 — Override-Variablennamen zwischen den beiden `sdk-pack-*`-Skripten uneinheitlich

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/sdk-pack-python.sh:38` (`SDK_PACK_PYTHON_IMAGE`)
  vs. `tools/harness/sdk-pack-csharp.sh:33` (`SDK_PACK_IMAGE`, ohne
  Sprachsegment)
- `befund`: Das C#-Skript nennt seine Override-Variable `SDK_PACK_IMAGE`
  (kein Sprachsegment), das neue Python-Skript `SDK_PACK_PYTHON_IMAGE` (mit
  Sprachsegment) — funktional beide korrekt, aber kein einheitliches
  Namensschema über die beiden Sprach-Pendants hinweg.
- `verifizierbar`: nein — reine Namensfrage, kein Gate prüft
  Variablennamen-Konsistenz.
- `klasse`: „Override-Variablenname ohne einheitliches Schema über
  Sprach-Pendants"

### F-2 — Slice-Namen in Produktionsnahen Kommentaren neben ADR-Referenz genannt

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.7)
- `pfad`: `sdks/python/Dockerfile:2,8,49`, `tools/harness/sdk-pack-python.sh:4`
- `befund`: Mehrere Kommentare nennen `slice-sdk-python-pack-werkzeug` als
  Herkunftsangabe direkt neben der ADR-Referenz, z. B.
  `(slice-sdk-python-pack-werkzeug, ADR-0107 §Konsequenzen Folgepflicht 1,
  ADR-0108 §Entscheidung Festlegung 1/3)` — indikativ über den aktuellen
  Zustand, keine Vorher/Nachher-Erzählung, und exakt die Form, die bereits
  im `sdk-pack-csharp`-Pendant (dort review-geprüft, kein HIGH-Finding)
  steht. Kein neuer Verstoß der HIGH-Klasse „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar" — reines Bestandsmuster, hier zur
  Vollständigkeit dokumentiert.
- `verifizierbar`: nein — Lese-Urteil, kein Gate fängt die Klasse.
- `klasse`: „Slice-Name als Herkunftsangabe neben ADR-Referenz (etabliertes
  Bestandsmuster)"

## Negativbefunde

- geprüft, ohne Befund: `harness/mk/sdk.mk`/`Makefile` — `sdk-pack-python`
  hängt real nicht an `GATE_CHECKS` (direkte grep-Prüfung der Fragmente,
  nicht nur der Zielname).
- geprüft, ohne Befund: realer Pack-Lauf — `make sdk-pack-python` Exit `0`,
  `.whl` (9330 Bytes, ZIP bestätigt) und `.tar.gz` (11290 Bytes, gzip
  bestätigt) real erzeugt, `.gitignore`-Wirkung real bestätigt
  (`git status --porcelain` leer).
- geprüft, ohne Befund: Mutations-Probe — roter Test bricht `docker build`
  an der `RUN pytest`-Zeile real ab (Exit ≠ 0), kein `dist/`-Verzeichnis
  entsteht; Mutation danach exakt zurückgesetzt, `git diff --stat` leer.
- geprüft, ohne Befund: `uv`-Digest-Reverifikation (`ADR-0108`
  §Entscheidung Festlegung 3) — Registry-Digest identisch zu ADR und
  Dockerfile, `COPY --from=`-Zeile byte-identisch zwischen ADR und
  Dockerfile.
- geprüft, ohne Befund: `pyproject.toml` `[build-system]` — unverändert
  (`setuptools.build_meta`), kein Diff an dieser Datei.
- geprüft, ohne Befund: Pipe-Disziplin (`AGENTS.md` §3.9) — isolierter Test
  bestätigt, dass ein `docker run`-Fehlschlag unter `set -o pipefail`
  sichtbar rot bleibt (`PIPE_EXIT=125`), nicht hinter `tar`s Exit-Code
  verschwindet.
- geprüft, ohne Befund: `make docs-check` real gefahren, Exit-Code direkt
  `0` (831 Dateien, 0 Befunde) — die `matrix-forbidden`-Umgehung in
  `spec/pflichtenheft.md` §1/§6/§7 trägt real.
- geprüft, ohne Befund: Import-Grenze (`grep -rn "internal/\|cmd/pg-change-feed"
  tools/harness/sdk-pack-python.sh harness/mk/sdk.mk` — kein Treffer).
- geprüft, ohne Befund: `SPEC-027`-Nummernkollision — `grep -c "SPEC-027"
  spec/*.md` zeigt genau eine Definition (§1-Erwähnung, §6-Zeile,
  §7-Historie, alle in `spec/pflichtenheft.md`), keine Zweitvergabe.
- geprüft, ohne Befund: Spec-Stratum-Disziplin — das Weglassen des
  `ADR-0107`/`ADR-0108`-Verweises in §1/§6/§7 folgt der bereits im Dokument
  etablierten Decken-Regel (§7-Kopfregel „kein ADR- und kein
  Slice-Verweis"), keine Ad-hoc-Ausnahme; die fachliche Aussage
  (Python/PyPI beantwortet, dritte Sprache/Vertriebsweg offen) bleibt
  präzise und deckt korrekt nur den HTTP-API-Umfang.
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — die neue Zeile
  ist formkonsistent mit `make sdk-pack-csharp`, referenziert real ein
  existierendes Target (`AGENTS.md` §4), Pipe im Fließtext korrekt als
  `\|` escaped.
- geprüft, ohne Befund: `sdks/python/Dockerfile` — `pack`/`pack-export`
  bauen strukturell (Docker-Layer-Abhängigkeit) auf einem bereits
  getesteten `build`-Stand auf, kein stiller Fallback bei rotem Test
  (Mutations-Probe bestätigt).
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als getrennte Commits
  (`AGENTS.md` §3.3) — beide Lifecycle-Moves (`2184fafd`, `b43142ca`)
  tragen 0 Insertions/Deletions; `Verantwortlich`-Feld eigener Commit
  (`8bd536c1`).
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — kein
  Konjunktiv über eine verworfene Alternative, kein abwesender Text, keine
  narrative Slice-/Wellen-Chronik in Produktionscode (dieser Slice ändert
  keine `.py`-Dateien); Slice-Nennungen in `Dockerfile`/`sdk-pack-python.sh`
  folgen der bereits etablierten Herkunfts-Anker-Form (siehe F-2, INFO).
- geprüft, ohne Befund: Traceability — alle vier Commits im Diff-Bereich
  nennen `LH-FA-SST-009` und `ADR-0107` im Betreff, kein `SPEC-*`/`ARC-*`
  im Betreff (Regex-Probe gegen `commit-traceability.sh`s Muster
  ausgeführt).
- geprüft, ohne Befund: Byte-Größen-Konsistenz (`AGENTS.md` §3.12
  Instanz A) — die im Slice-Plan/in der Commit-Message genannten Größen
  (9330/11290 Bytes) stimmen mit der selbst gemessenen `ls -la`-Ausgabe
  überein, kein Drift.
- geprüft, ohne Befund: `docs/plan/planning/observations/` — keine neue
  Beleg-Datei im Diff, Slice-Plan §7 begründet dies explizit („keine
  Beobachtung angefallen"); kein Widerspruch zu den eigenen Findings
  dieses Reports (F-1/F-2 sind beide unterhalb der Beobachtungs-Schwelle:
  LOW ohne Wiederholung, INFO über ein Bestandsmuster).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Override-Variablenname ohne
einheitliches Schema über Sprach-Pendants" (erstes Auftreten) · „Slice-Name
als Herkunftsangabe neben ADR-Referenz (etabliertes Bestandsmuster)"
(kein neuer Fehlerklassen-Zähler, siehe F-2).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 1 LOW, 1 INFO. Keine
Fixrunde am Implementer nötig; F-1 ist eine stilistische Beobachtung ohne
semantische Auswirkung (beide Variablen funktionieren korrekt und
unabhängig voneinander).

**DoD-Checkbox-Nachzug ohne Fixrunde** (Skill-Regel, `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde): Da dieses Verdikt zu keiner Fixrunde
führt, wird die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-sdk-python-pack-werkzeug.md`) im
selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen, mit
Verweis auf diesen Report-Pfad.

**Übergabe:** kein Rückgabe-Pfeil an den Implementer nötig. Dieser Report
ist ein Lauf-Beleg; er ersetzt keine Verifikation gegen die volle DoD —
das bleibt Verifier-Aufgabe (Modul 11).
