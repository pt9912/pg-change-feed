# Review-Report: slice-sdk-csharp-pack-werkzeug — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-pack-werkzeug.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `89f62131..HEAD` (Abschluss von
`slice-sdk-csharp-grpc-client-flaeche` bis Implementer-Commit dieses
Slice), Slice `slice-sdk-csharp-pack-werkzeug`, Welle
`welle-sdk-csharp-lh-fa-sst-009`. Vier Commits: `4a7b0d62` (open→next,
reiner Move, 0 Insertions/Deletions), `bb23e3d8` (Verantwortlich gesetzt,
nur Slice-Datei, 2 Insertions/1 Deletion), `8438ed95` (next→in-progress,
reiner Move, 0 Insertions/Deletions), `dcd05acc` (Inhalt: `harness/mk/sdk.mk`,
`tools/harness/sdk-pack-csharp.sh`, `sdks/csharp/Dockerfile`-Erweiterung,
`sdks/csharp/.gitignore`, `spec/pflichtenheft.md`-Träger-Nachzug,
`harness/README.md`-Werkzeugzeile).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. `AGENTS.md` §3.13 Träger-Nachzug, §3.9 Pipe-Disziplin).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-pack-werkzeug.md` (§1
  Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken, §8
  Sub-Area/Modus)
- `docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md` (Accepted) —
  Festlegung 4 (Docker-only bauen, kein Gate), §Konsequenzen
  Folgepflicht 2/3
- `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a`, §6 Externe Verträge, §7
  Historie
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.6 (Gates nur per ADR gelockert), §3.7 (Kommentar-Disziplin), §3.9
  (Pipe-Disziplin/Exit-Code), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug), §4 (kein Träger nennt ein nicht existentes Target)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `harness/README.md` §Sensors (Gate-Index), Vergleichszeile
  `make examples-csharp`
- `docs/reviews/review-slice-sdk-csharp-grpc-client-flaeche.md` — Formvorbild
  für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `grep -n "sdk.mk\|include harness/mk" Makefile` und
  `grep -n "GATE_CHECKS" harness/mk/*.mk a-check.mk` real ausgeführt:
  `harness/mk/sdk.mk` trägt **kein** `GATE_CHECKS +=` (im Gegensatz zu
  `baseline.mk`, `doc-gate.mk`, `generated-sync.mk`, `coverage.mk`,
  `a-check.mk`, die es alle tragen) — `make sdk-pack-csharp` ist real
  außerhalb der Gate-Kette, nicht nur dem Namen nach.
- `rm -rf sdks/csharp/dist && make sdk-pack-csharp` real ausgeführt:
  Exit-Code direkt `0`. `sdks/csharp/dist/PgChangeFeed.Client.0.1.0.nupkg`
  existiert danach real, `ls -la` zeigt 29947 Bytes — identisch mit dem in
  der Commit-Message (`dcd05acc`) genannten Wert, kein Drift. `file`
  bestätigt „Zip archive data" (ein `.nupkg` ist ein ZIP). `git status
  --porcelain` zeigt das Artefakt danach **nicht** als untracked —
  `sdks/csharp/.gitignore`s neue `dist/`-Zeile trägt real.
- **Mutations-Probe** (wie vom Implementer im Plan-Nachzug beschrieben,
  hier unabhängig wiederholt): `Assert.Equal("test-token", …)` in
  `PgChangeFeedClientOptionsTests.cs` auf `"mutated-wrong-token"`
  geändert, `rm -rf sdks/csharp/dist && make sdk-pack-csharp` erneut
  real ausgeführt — Bau bricht sichtbar an der `dotnet test`-Stufe ab
  (`Failed! - Failed: 1, Passed: 31, ...`, `docker build` Exit 1, `make`
  meldet `Error 1`/Gesamt-Exit `2`), **kein** `sdks/csharp/dist/`-
  Verzeichnis wird angelegt (die Zeile `mkdir -p sdks/csharp/dist` im
  Skript steht **nach** dem `docker build`-Aufruf und wird bei dessen
  Fehlschlag nie erreicht). Mutation danach exakt zurückgesetzt
  (`mv PgChangeFeedClientOptionsTests.cs.bak PgChangeFeedClientOptionsTests.cs`),
  `git diff --stat`/`git status --porcelain` beide leer bestätigt — kein
  Restartefakt.
- **Pipe-Disziplin real geprüft** (`AGENTS.md` §3.9): isolierter Test
  `bash -c 'set -o pipefail; docker run --rm --network none
  nonexistent-image-xyz-123 | tar -x -C /tmp; echo "PIPE_EXIT=$?"'` —
  Ergebnis `PIPE_EXIT=125` (Docker-Fehlercode), **nicht** `0` (den ein
  ungeschütztes `| tar -x` bei leerem Stream liefern könnte). Bestätigt:
  `set -o pipefail` in `tools/harness/sdk-pack-csharp.sh` (kombiniert mit
  `set -euo pipefail` und explizitem `bash`-Interpreter über
  `@bash tools/harness/sdk-pack-csharp.sh` im Makefile-Rezept) macht einen
  `docker run`-Fehlschlag sichtbar rot, exakt wie im
  `proto-generate.sh`-Vorbild.
- `make docs-check` real ausgeführt: Exit-Code direkt `0`
  (`d-check: 806 Datei(en) geprüft, 0 Befund(e)`) — die im Plan-Nachzug
  beschriebene Umgehung des `matrix`-Moduls (`spec → adr` verboten) trägt
  tatsächlich; kein `matrix-forbidden`-Befund in der neuen
  `LH-FA-SST-009.a`-Nachzugzeile oder der neuen `SPEC-026`-Zeile.
- `grep -rn "internal/\|cmd/pg-change-feed" tools/harness/sdk-pack-csharp.sh
  harness/mk/sdk.mk` — kein Treffer, kein Import-Grenzverstoß.
- `grep -c "SPEC-026" spec/*.md` — drei Treffer, alle in
  `spec/pflichtenheft.md` (§1-Fließtext, §6-Tabellenzeile, §7-Historie),
  `spec/lastenheft.md`/`spec/architecture.md` je `0` — keine Kollision mit
  einer bereits vergebenen Nummer, genau eine Definition.
- `spec/pflichtenheft.md` §3-Kopfregel gelesen („Regeln dieser Sektion: Die
  ADR, die einen Wert festlegt, deklariert das aufwärts in ihrem
  `Schärft:`-Feld — kein ADR-Rückzeiger hier") und §7-Kopfregel („kein
  ADR- und kein Slice-Verweis") — beide bestätigen: das Weglassen des
  `ADR-0106`-Verweises in §1/§6/§7 ist **keine** Ad-hoc-Notlösung wegen
  `matrix-forbidden`, sondern folgt einer bereits im Dokument etablierten,
  durchgängigen Struktur-Regel. Die neue `LH-FA-SST-009.a`-Nachzugzeile
  trägt die fachliche Aussage präzise („für C#/NuGet ist die Frage
  beantwortet … eine zweite Sprache oder ein zweiter Vertriebsweg bleibt
  offen — diese Kennung bleibt ihre Adresse") ohne inhaltliche Verwässerung.
- `SPEC-026`-Zeile gegen `SPEC-023` verglichen: Form (ID/System/Version/
  Vertrag-Datei) konsistent; die Vertrag-Datei-Spalte weicht bewusst vom
  „— (Vertrag steht in diesem Dokument, §2 …)"-Muster ab, weil `SPEC-026`
  **keinen** eigenen §2-Eintrag hat (explizit begründet: „das Package
  deckt bereits dokumentierte Drahtverträge") und stattdessen real auf
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` verweist —
  Datei real vorhanden, `<PackageId>PgChangeFeed.Client</PackageId>` und
  `<Version>0.1.0</Version>` real gelesen, deckt sich mit der Zeile.
- `sdks/csharp/Dockerfile` vollständig gelesen: `pack`-Stufe (`FROM build
  AS pack`, `--no-build` gegen den bereits in `build` kompilierten Stand)
  läuft **nach** den beiden `RUN dotnet build`/`RUN dotnet test`-Zeilen
  derselben Stufe `build` — ein roter Test in `build` verhindert
  strukturell, dass `pack` je erreicht wird (Docker-Layer-Abhängigkeit,
  nicht nur Doku-Behauptung) — durch die Mutations-Probe oben real
  bestätigt.
- `harness/README.md` §Werkzeuge, neue `make sdk-pack-csharp`-Zeile gegen
  die `make examples-csharp`-Zeile verglichen: gleiche Form (Bau-Kontext,
  Stufen, Testlauf-vor-Artefakt-Zusage, Export-Mechanismus, `.gitignore`-
  Hinweis, „Braucht Netz … Werkzeug statt Gate", `kein Gate,
  [ADR-Link] … · seit slice-<name>`"). Die Pipe im Fließtext ist korrekt
  als `\|` escaped (Tabellenzelle bleibt intakt).
- `git show 4a7b0d62 --stat`/`git show 8438ed95 --stat` geprüft: beide
  zeigen ausschließlich den Rename, 0 Insertions/Deletions — reine Moves,
  `AGENTS.md` §3.3 sauber eingehalten; `bb23e3d8` trägt ausschließlich das
  `Verantwortlich`-Feld der Slice-Datei.
- Neue/geänderte Dateien nach Slice-/Wellen-Chronik-Mustern durchsucht
  (`AGENTS.md` §3.7): Treffer in `Dockerfile`, `harness/mk/sdk.mk`,
  `tools/harness/sdk-pack-csharp.sh`, `.gitignore` — durchweg in der
  etablierten Herkunfts-Anker-Form („(slice-<name>, ADR-0106
  Festlegung …)"), indikativ über den aktuellen Zustand, keine
  Vorher/Nachher-Erzählung; kein Treffer in `.cs`-Produktionscode (dieser
  Slice ändert keine `.cs`-Dateien).
- Commit-Traceability real geprüft (`git log --oneline 89f62131..HEAD`):
  alle vier Commit-Betreffe tragen beide Kennungen aus dem Bezug-Feld dieses
  Slice-Plans, kein Struktur-ID im Betreff (siehe Negativbefund unten).

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings in diesem Lauf.

### F-1 — SPEC-026-Vertrag-Datei-Spalte weicht in der Form vom SPEC-023-Muster ab

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `spec/pflichtenheft.md:611`
- `befund`: Die `SPEC-026`-Zeile in §6 nennt in der Spalte „Vertrag-Datei"
  direkt einen Repo-Pfad (`sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`)
  statt, wie `SPEC-010`/`SPEC-017`/`SPEC-020`/`SPEC-023`/`SPEC-024`, mit
  „— (Vertrag steht in diesem Dokument, §2 …)" zu beginnen — sachlich
  begründet (kein eigener §2-Eintrag), aber eine neue Formvariante in
  dieser Tabelle.
- `verifizierbar`: nein — reine Formfrage, kein Gate prüft die Spaltenform.
- `klasse`: „Tabellenspalten-Formvariante ohne Vorbild in derselben Spalte"

## Negativbefunde

- geprüft, ohne Befund: `harness/mk/sdk.mk`/`Makefile` — `sdk-pack-csharp`
  hängt real nicht an `GATE_CHECKS` (direkte grep-Prüfung der Fragmente,
  nicht nur der Zielname).
- geprüft, ohne Befund: realer Pack-Lauf — `make sdk-pack-csharp` Exit `0`,
  `.nupkg` real erzeugt (29947 Bytes, ZIP-Archiv bestätigt via `file`),
  `.gitignore`-Wirkung real bestätigt (`git status --porcelain` leer).
- geprüft, ohne Befund: Mutations-Probe — roter Test bricht `docker build`
  real ab (Exit ≠ 0), kein `.nupkg`/kein `dist/`-Verzeichnis entsteht;
  Mutation danach exakt zurückgesetzt, `git diff --stat` leer.
- geprüft, ohne Befund: Pipe-Disziplin (`AGENTS.md` §3.9) — isolierter Test
  bestätigt, dass ein `docker run`-Fehlschlag unter `set -o pipefail`
  sichtbar rot bleibt (`PIPE_EXIT=125`), nicht hinter `tar`s Exit-Code
  verschwindet.
- geprüft, ohne Befund: `make docs-check` real gefahren, Exit-Code direkt
  `0` (806 Dateien, 0 Befunde) — die `matrix-forbidden`-Umgehung in
  `spec/pflichtenheft.md` §1/§6/§7 trägt real.
- geprüft, ohne Befund: Import-Grenze (`grep -rn "internal/\|cmd/pg-change-feed"
  tools/harness/sdk-pack-csharp.sh harness/mk/sdk.mk` — kein Treffer).
- geprüft, ohne Befund: `SPEC-026`-Nummernkollision — `grep -c "SPEC-026"
  spec/*.md` zeigt genau eine Definition (§1-Erwähnung, §6-Zeile,
  §7-Historie, alle in `spec/pflichtenheft.md`), keine Zweitvergabe.
- geprüft, ohne Befund: Spec-Stratum-Disziplin — das Weglassen des
  `ADR-0106`-Verweises in §1/§6/§7 folgt der bereits im Dokument
  etablierten Decken-Regel (§3-Kopfregel, §7-Kopfregel „kein ADR- und kein
  Slice-Verweis"), keine Ad-hoc-Ausnahme; die fachliche Aussage
  (C#/NuGet beantwortet, zweite Sprache/Vertriebsweg offen) bleibt
  präzise.
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — die neue Zeile
  ist formkonsistent mit `make examples-csharp`, referenziert real ein
  existierendes Target (`AGENTS.md` §4), Pipe im Fließtext korrekt als
  `\|` escaped.
- geprüft, ohne Befund: `sdks/csharp/Dockerfile` — `pack`/`pack-export`
  bauen strukturell (Docker-Layer-Abhängigkeit) auf einem bereits
  getesteten `build`-Stand auf, kein stiller Fallback bei rotem Test
  (Mutations-Probe bestätigt).
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als getrennte Commits
  (`AGENTS.md` §3.3) — beide Lifecycle-Moves (`4a7b0d62`, `8438ed95`)
  tragen 0 Insertions/Deletions; `Verantwortlich`-Feld eigener Commit
  (`bb23e3d8`).
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — kein
  Konjunktiv über eine verworfene Alternative, kein abwesender Text, keine
  narrative Slice-/Wellen-Chronik in Produktionscode (dieser Slice ändert
  keine `.cs`-Dateien); Slice-Nennungen in `Dockerfile`/`sdk.mk`/
  `sdk-pack-csharp.sh`/`.gitignore` folgen der etablierten
  Herkunfts-Anker-Form.
- geprüft, ohne Befund: Traceability — alle vier Commits im Diff-Bereich
  nennen `LH-FA-SST-009` und `ADR-0106` im Betreff, kein `SPEC-*`/`ARC-*`
  im Betreff.
- geprüft, ohne Befund: Byte-Größen-Konsistenz (`AGENTS.md` §3.12
  Instanz A) — die in der Commit-Message genannte Größe
  (29947 Bytes) stimmt mit der selbst gemessenen `ls -la`-Ausgabe überein,
  kein Drift.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Tabellenspalten-Formvariante ohne
Vorbild in derselben Spalte" (erstes Auftreten, INFO, kein Steering-Loop-
Zähler-Eintrag nötig unter 3x).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO. Keine
Fixrunde am Implementer nötig.

**DoD-Checkbox-Nachzug ohne Fixrunde** (Skill-Regel, `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde): Da dieses Verdikt zu keiner Fixrunde
führt, wird die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-pack-werkzeug.md`) im
selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen, mit
Verweis auf diesen Report-Pfad.

**Übergabe:** kein Rückgabe-Pfeil an den Implementer nötig. Dieser Report
ist ein Lauf-Beleg; er ersetzt keine Verifikation gegen die volle DoD —
das bleibt Verifier-Aufgabe (Modul 11).
