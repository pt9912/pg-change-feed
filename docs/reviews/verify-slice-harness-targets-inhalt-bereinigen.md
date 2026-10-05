# Verifikationsbericht: slice-harness-targets-inhalt-bereinigen — 2026-10-05

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD (`slice-harness-targets-inhalt-bereinigen` §2, Liefer-Punkte 1
und 2 und Gate-Pflicht), gegen
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) und
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und gegen die
Hard Rules [`AGENTS.md`](../../AGENTS.md) §3.7, §3.10 und §3.12. Nicht geprüft wird
der Diff als Maintainability-Frage. Das hat der Reviewer getan
([`review-slice-harness-targets-inhalt-bereinigen.md`](review-slice-harness-targets-inhalt-bereinigen.md)).
Der reale Bedarf ist Aufgabe des Validators; dies ist kein MVP-Slice.

**Gegenstand:** Diff `c38d3a6b..204db336`.
- Implementierung: `0cdeb963` (Sensor- und Target-Dateien, `harness/README.md`,
  Kopfkommentar `.github/workflows/ci.yml`), `e6e336bc` (Plan: Anker-Gegenprobe).
- Review: `204db336`.

Der Slice liegt in `in-progress/`. Die Closure-Punkte (Doku-Update, Closure-Notiz,
Register, Risiko-Ausgänge, Paarungen) stehen noch aus und sind nicht Gegenstand.

**Frischer Kontext:** Diese Sitzung hat Plan, Review-Report und Diff gelesen. Das
Skript `gegenprobe.sh` des Implementers ist nicht benutzt; die Gegenprobe unten
ist eine eigene Methode. Jede Zahl unten ist in diesem Lauf am Stand `204db336`
gemessen, außer wo **übernommen** steht. Die Exit-Codes sind direkt und ohne
Pipe gesichert (`AGENTS.md` §3.9). Kopien und Mutationen liegen in einem eigenen
Unterverzeichnis `verifier-targets/` des Scratchpads; am Arbeitsbaum ist keine
Repo-Datei außer diesem Bericht geschrieben.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| `git grep -l '^## Fassung im Gate-Index' -- harness/` | keine Zeile (am Parent `c38d3a6b`: 15 Dateien, 11 `sensors/`, 4 `targets/`) | 1 (kein Treffer) |
| Gegenprobe Wort-Deckung je Satz, eigene Methode (§2) | alle 15 Fassungen, Summe 24.253 Byte, gleich der Startmessung im Plan | — |
| `make test-handbuch-public-doc-check` | `run-handbuch-public-doc-check-tests: alle 41 Fälle bestanden` | 0 |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-targets-inhalt-bereinigen.md` | `suchlauf-nachmessen: 16 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1738 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=c38d3a6b..HEAD` | `d-check: 1738 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-immutable RANGE=c38d3a6b..HEAD` | `d-check: 1738 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make pin-stale-race` (Netz) | `OK          TOOLCHAIN_RACE_IMAGE (golang:1.27) == sha256:e0174e51…` | 0 |
| `make image-cve` (Netz) | `ghcr.io/pt9912/pg-change-feed:latest (debian 12.15) │  debian  │        0        │` | 0 |
| `make baseline-verify` | `baseline-verify: v6.14.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` | 0 |
| `gh run view 37258852180 --json jobs` | `release`, `v0.7.0`; Job „Docker-Hub-Beschreibung synchronisieren / …“ `success` | 0 |
| `gh run list --workflow hub-description.yml --limit 50` | **sechs** Läufe, alle `success`, `workflow_dispatch`, zuletzt 35483322058 | 0 |
| Anker-Gegenprobe README-Link (§4) | ohne Mutation 0 Befunde; mit Mutation 1 Befund `anchor-missing` | 0 / 1 |

`make gates` läuft nach dem Commit dieses Berichts; sein Exit steht in der
Rückmeldung an den Auftraggeber, nicht hier.

## 2. Liefer-Punkt 1 — Fassungen zusammengeführt

**Eigene Methode.** Je Datei ist die alte Fassung aus `git show c38d3a6b:<pfad>`
von der Überschrift `## Fassung im Gate-Index` bis zum Dateiende ausgeschnitten.
Ein `awk`-Skript zerlegt sie in Sätze und misst je Satz den Anteil seiner
Inhaltswörter (ab sechs Zeichen, normalisiert), der in der neuen Datei steht.
Jeder Satz unter 100 % ist danach von Hand gegen den neuen Vertrag gelesen. Die
Methode ist eine andere als die des Implementers (Backtick-Spans, Zahlen und
Zahlwörter): sie trifft auch Prosa-Aussagen ohne Eigennamen.

**Ergebnis.** In jeder der 15 Dateien liegen die ersten beiden Sätze bei 0–58 %:
der Herkunfts-Satz („wortgleich umgezogene Index-Text“) und der Vorrang-Satz. Die
DoD streicht beide ausdrücklich. Jeder weitere Satz unter 100 % ist am neuen Text
gelesen; die fehlenden Wörter sind Wortformen, nicht Aussagen. Die Fundstellen:

| Datei | gelesene Aussage der Fassung | steht im neuen Vertrag |
|---|---|---|
| `sensors/handbuch-public-doc-check.md` | vier `*-abdeckung.md` ausgenommen, Maintainer-Doku außerhalb; jede andere `*.md` Exit 2; Chronik-Sprache grün, Lese-Hälfte Skill und Reviewer | Tabelle „ausgenommen“ (Zeile 32), Zeilen 35–36, 46, Grenze 1 und 3 |
| `sensors/kommentar-kennungen.md` | Host-Werkzeuge `bash`, `git`, `docker`; Exit 1, über `make` Exit 2; Aufrufer Schritt 20 und Reviewer; Argument-Zerlegung, Diff-Strom unter fremder Git-Konfiguration | Zeilen 92, 108–117, §Wer es aufruft, §Test Zeilen 171–177 |
| `sensors/sdk-public-doc-check.md` | Wheel, `.nupkg`, Sources-Jar, `sdist`; Ausnahmen bis `grpc_gen`; `test_public_text.py`; Vorgänger-Kante der drei `sdk-pack-*` | Zeilen 18–27, neuer Abschnitt §Vorstufe der Paket-Bauten |
| `sensors/pin-stale-all.md` | `git grep`, dedupliziert, erster Fundort, fail-open, leerer Gegenstand Exit 2, Tag-Wechsel und Laufzeit-Digest außerhalb, `PROG=<Datei>` | Zeilen 15–21, 38–48, Grenze 2, §Overrides, §Tabellentest |
| `sensors/generated-sync.md` | Build-Zeit-Erzeugung, host-seitige `tar`-Extraktion ins Temp-Verzeichnis, kein `--user`-Workaround, byte-gleich | §Vertrag Zeilen 5–19 und 51; byte-gleich am Skript gemessen (`cmp -s`, `tools/harness/generated-sync.sh` Zeile 115), beide `gen/`-Bäume am `Dockerfile` (Stufe `proto-export`, beide `.proto`) |
| `targets/image-mutation.md` | Bau- und Entfern-Befehl, kein `-f`/`prune`, Eingabeprüfung, `$(error …)`, Hash-Dateien und `:dev` unberührt, Fallliste | Zeilen 10–36, 50, 65–66, §Test Zeilen 124–133 |
| `targets/sdk-altserver.md` | B0 bis U, Index-Digest, Override, Exit 1 ohne `:dev`, „schreibt nichts in den Arbeitsbaum“, Werkzeug statt Gate | Zeilen 12, 30–37, 54, 74–85 |
| `targets/sdk-kompat.md` | A1 bis A5, Binärdateien, Quelltext, `docker run` statt `RUN`, `SDK_KOMPAT_NEU`/`REGISTRY_VERSION`, Netz | Zeilen 23–40, 76, 107, 119–120, Grenze 5 |
| `targets/schema-rollout.md` | Precheck, Wache, View-Vorlauf, `plan.yaml`/`down.sql`, vier Nacharbeit-Schritte, stdin, ohne Bind-Mount, sechs Läufe, Unit-Tests der Wache | §Ablauf Schritte 1–6 samt Tabelle, §Erzeugnisse, §Belege Zeilen 233–247 |

Die übrigen sechs Dateien (`ausgabe-kennungen-check`, `coverage-gate`,
`docs-check`, `fmt-check`, `meldungscodes-check`, `suchlauf-nachmessen`) erreichen
je Satz 75–100 %; ihre Lücken (`Kommentare`, `Docker-only`, `Katalog-Zeilen`,
`Plan-Datei`) stehen in anderer Wortform im Vertrag (Zeilen gelesen). Die
Streich-Liste des Plans (7 Einträge) ist an den genannten Ersatzstellen gelesen.

**Kein Inhaltsverlust gefunden.** Liefer-Punkt 1 ist bestätigt: 0 Treffer, jede
Aussage steht im Vertrag oder ist begründet gestrichen.

**Widerspruch „vier“ gegen „fünf“ am Gegenstand.** `excluded=` in
`tools/harness/handbuch-public-doc-check.sh` Zeilen 31–32 trägt vier Namen;
`docs/user/` führt genau die vier `*-abdeckung.md`; der Tabellentest trägt vier
Fälle `expect 0 "… ausgenommen"` (Zeilen 90–93) und einen Exit-2-Fall „fehlende
ausgenommen-genannte Datei“ (Zeile 116), der kein Ausnahme-Fall ist. Der Vertrag
sagt an allen drei Stellen „vier“ (Zeilen 64, 97, 98). Entscheidung „vier“ trägt.

**Verweis.** Die README-Bindung `make doc-tracked` zeigt auf
`harness/sensors/docs-check.md#make-doc-tracked`; der Abschnitt steht in Zeile 199.

## 3. Liefer-Punkt 2 — Befund → Ausgang → Beleg

| Punkt | Nachgefahren | Ergebnis |
|---|---|---|
| 2 `make image-cve` | gemessen (Netz), dazu `--exit-code 1` und `GHCR_USERNAME`/`GHCR_PASSWORD` am Ziel im `Makefile` Zeilen 110–111 | trägt; Zeile und Exit gleich dem Vertrag |
| 3 `hub-description.yml` | gemessen (`gh run view`, `gh run list`) | trägt in der Sache; **Zahl weicht ab**: `gh run list` druckt sechs Läufe `success`, der Plan nennt „fünf Läufe“, der Vertrag „die fünf letzten Läufe“ (Befund V-1) |
| 4 Artefakt-Namen | `ls sdks/*/dist` | trägt; fünf Dateien je `0.6.1` |
| 5 `make pin-stale-*` | `make pin-stale-race` gemessen; `-pgtest`, `-dmigrate` **übernommen** aus dem Plan, `-acheck` **übernommen** aus dem Review (Probe 3) | trägt für `-race` |
| 6 Kopierfehler SDK-Integration | Marker je Runner gelesen (`run-sdk-csharp-…` 537, `-kotlin-…` 532, `-python-…` 611) | trägt |
| 7 vier/fünf | §2 | trägt |
| 8 Positionsverweis | Link auf `harness/targets/schema-rollout.md`; „elf“ dort Zeilen 114, 238 | trägt |
| 10a `baseline-verify` | gemessen; Grenze 3 zeigt auf `make pin-stale-baseline`, der Anker löst auf (`make docs-check` Exit 0) | trägt |
| 10b `ci.yml` | siehe §5 | trägt |
| 11 Chronik | Wort-Diff von `targets/test-integration.md`, `targets/image.md`, `targets/examples.md` gelesen | trägt für die Befund-Stellen; verbleibende Chronik in anderen Absätzen ist Reviewer-F-4, außerhalb der Befund-Liste |

## 4. Anker-Gegenprobe README-Link

Zwei Kopien von `HEAD` per `git archive` im Scratchpad. In einer ist der Link in
`harness/README.md` Zeile 139 per `sed … > Kopie` auf `#make-doc-trackedX`
geändert (`diff` zeigt genau diese Zeile). d-check mit dem gepinnten Image,
`--disable tracked` (die Kopien tragen kein `.git`):

- ohne Mutation: Exit 0, `0 Befund(e)`;
- mit Mutation: Exit 1, `1 Befund(e)`, Zeile
  `harness/README.md:139	sensors/docs-check.md#make-doc-trackedX	anchor-missing`.

Der Anker ist von `anchors` gedeckt; der Sensor fiele bei einem Bruch rot.

## 5. Entscheidungs- und Hard-Rule-Konformität

| Regel | Prüfung | Ergebnis |
|---|---|---|
| `AGENTS.md` §3.10 / [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) | `git diff -U0 c38d3a6b..HEAD -- .github/workflows/ci.yml`: 6 Zeilen entfernt, 6 hinzugefügt, alle Kommentarzeilen (Filter auf Nicht-Kommentar: 0 Zeilen); `tags-ignore:` weiter Zeile 40, das Zitat „ci.yml Zeile 40f.“ in `sdk-csharp-release.yml` löst auf | keine strukturelle Änderung, §3.10 nicht ausgelöst |
| `AGENTS.md` §3.12 / [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | neue Zahlen in den Verträgen: „54 Dateien“, `debian 12.15` mit `0`, die vier `OK`-Zeilen, „vier“, Lauf `37258852180` | tragen Lauf und Datum, die zwei „real geprüft“-Sätze in `make image` sind als übernommen gekennzeichnet; Ausnahme V-1 |
| `AGENTS.md` §3.7 | hinzugefügte Zeilen nach Chronik-Wörtern durchsucht | an den Befund-Stellen Ist-Zustand; Treffer nur in vorbestehenden Absätzen außerhalb der Befund-Liste (Reviewer F-4) |
| Reviewer-Commit `204db336` | `git show 204db336 -- docs/plan/` | am Plan nur die DoD-Zeile „Review durchgeführt“ (Häkchen und Beleg-Satz) |
| Traceability, Immutabilität | `make doc-commits`, `make doc-immutable` über `c38d3a6b..HEAD` | je Exit 0 |

## 6. Plan gegen Code-Diff

Geänderte Dateien (`git diff --stat c38d3a6b..HEAD`): die 15 Dateien der Fassungen,
die Dateien der Befund-Liste (`targets/image.md`, `workflows.md`, `sdk-pack.md`,
`pin-stale.md`, `sdk-integration.md`, `examples.md`, `test-integration.md`,
`sensors/baseline-verify.md`, `.github/workflows/ci.yml`), `harness/README.md` (eine
Bindung-Zelle), Plan und Review. `generated-sync.md` und `sdk-public-doc-check.md`
stehen zusätzlich im Plan-Nachzug §3 mit Begründung. Kein Skript, kein Gate, keine
`Accepted`-ADR geändert (§1 Abgrenzung eingehalten). Keine Datei im Diff ohne
Plan-Zeile.

## 7. Abweichungen

| ID | Klasse | Befund | Beleg |
|---|---|---|---|
| V-1 | LOW | Die Zahl der `hub-description.yml`-Läufe weicht von der Messung ab: Plan §3 Tabelle Punkt 3 nennt „fünf Läufe `success`“, `harness/targets/workflows.md` „die fünf letzten Läufe über `workflow_dispatch`“. `gh run list --workflow hub-description.yml --limit 50` druckt sechs, alle `success`, alle `workflow_dispatch`, zuletzt 35483322058. Der Vertragssatz ist wahr, aber die Zahl ist nicht die gedruckte (§3.12, Instanz A). | `gh run list` dieses Laufs |
| V-2 | LOW (übernommen aus Review F-1, bestätigt) | Die Meldung fremder Träger ist unvollständig: `.github/workflows/hub-description.yml` Zeile 29, `sdk-csharp-release.yml` Zeilen 38–39 und `sdk-python-release.yml` Zeilen 52–53 tragen „unbewiesen“. Dazu, über F-1 hinaus: `.github/workflows/e2e.yml` Zeile 35 („Ob das Zeitlimit real ausreicht, bleibt bis zum ersten echten Lauf offen“), während `gh run list --workflow e2e.yml` grüne Läufe zeigt (37301088944, 37296720796). Frist: Closure (`AGENTS.md` §3.13). | `git grep -n -E 'unbewiesen\|ersten echten (Lauf\|Tag-Push\|Release)' -- .github harness tools` |

Reviewer F-2 (Zeiger „Grenze 7“ auf eine Doppelnummer) ist am Stand unverändert und
geht wie vereinbart in die Closure; kein eigener Befund.

## 8. Verdikt

**DoD bestätigt** für Liefer-Punkt 1, Liefer-Punkt 2 und die DoD-Zeile „Review
durchgeführt“: ja. `make docs-check` Exit 0 (gemessen). Die DoD-Zeile `make gates`
bleibt offen bis zum Lauf am Endstand; die Closure-Punkte stehen aus. V-1 und V-2
sind LOW und blockieren nicht; die Closure trägt sie (V-1 berichtigen oder als
Ausschnitt kennzeichnen, V-2 in der Träger-Meldung nachtragen).

**Übergabe:** an den Planner zur Closure.
