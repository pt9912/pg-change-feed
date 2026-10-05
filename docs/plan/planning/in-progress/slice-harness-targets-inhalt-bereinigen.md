# Slice harness-targets-inhalt-bereinigen: die Sensor- und Target-Dateien tragen ihren Inhalt einmal und im Ist-Zustand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (CI/CD-Pipeline über GitHub Actions, mit Pin-Inventar),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen: eine berichtigte Zahl trägt ihren Ursprung). Keine `LH-*`-Anforderung
ist berührt: der Slice ändert Harness-Dokumente, nicht das Produkt.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).

**Autor:** pt9912 (Planner, Closure von `slice-harness-readme-zellen-kuerzen`).
**Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**.

**Anlass.** `slice-harness-readme-zellen-kuerzen` zog die ausführlichen Zellen
der `harness/README.md` wortgleich in Sensor- und Target-Dateien um und durfte
keinen Wortlaut ändern. Seine Closure-Notiz (§7, Befund-Liste Punkte 2 bis 11)
nennt, was dabei als doppelt, widersprüchlich oder veraltet auffiel; 15 Dateien
tragen am Ende einen Abschnitt `## Fassung im Gate-Index` neben dem eigenen
Vertrag (Review-Befund F-1, `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`).

**Ziel:** Jede Sensor- und Target-Datei trägt ihren Vertrag einmal, ohne
`## Fassung im Gate-Index`, und jede Aussage aus der Befund-Liste ist nachgemessen,
als übernommen gekennzeichnet oder gestrichen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Index-Zeilen der `harness/README.md`.** Sie sind kurz und von der
  `structure`-Regel gehalten; ändert sich ein Kurzsatz durch eine Berichtigung,
  wird er mitgezogen, sonst bleibt die README unberührt. Die Vorlagen-Platzhalter
  der README (§Sensors-Musterzeile, Werkzeug-Musterzeilen, §Safety, §Leseordnung)
  trägt `slice-abgeleitete-dokumente-vorlagen-nachzug`.
- **Neue Gates oder Änderungen an Skripten.** Die Berichtigung betrifft Prosa;
  stellt sich eine Aussage als Defekt eines Werkzeugs heraus, ist das ein anderer
  Vorgang mit eigenem Slice.
- **Die `Accepted`-ADRs, die eine Zeile der alten Tabelle zitieren.** Sie frieren
  ihren Stand ein (`AGENTS.md` §3.5).

**Keine Mindestzahl.** Die vier Klassen des Ausschlusses sind ein Suchraster,
keine Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [x] **Fassungen zusammenführen (Liefer-Punkt 1).** In den 15 Dateien mit
      `## Fassung im Gate-Index` (`git grep -l '^## Fassung im Gate-Index' --
      harness/`) ist der Inhalt der Fassung in die Abschnitte davor eingearbeitet,
      der Abschnitt samt Vorrang-Satz entfernt; ein Widerspruch zwischen beiden
      (z. B. „die vier“ gegen „die fünf“ ausgenommenen Dateien in
      `harness/sensors/handbuch-public-doc-check.md`) ist am Gegenstand gemessen
      entschieden, und Verweise auf `§Fassung im Gate-Index` (README-Bindung
      `make doc-tracked`) zeigen auf den neuen Ort. Zu belegen durch: der
      `git grep` oben liefert 0 Treffer; je Datei nennt der Bericht, welche Sätze
      zusammengeführt und welche gestrichen wurden. — Beleg: `git grep -l
      '^## Fassung im Gate-Index' -- harness/ | wc -l` druckt `0` (2026-10-05);
      Tabelle je Datei, Gegenprobe und Widerspruch in §3 Umsetzung.
- [x] **Veraltete Aussagen berichtigen (Liefer-Punkt 2).** Jeder Punkt 2 bis 8,
      10 und 11 der Befund-Liste von `slice-harness-readme-zellen-kuerzen` §7
      ist erledigt: nachgemessen (Lauf genannt), als übernommen gekennzeichnet
      oder gestrichen; Chronik-Sprache steht im Ist-Zustand (`AGENTS.md` §3.7).
      Zu belegen durch eine Tabelle Punkt → Ausgang → Beleg im Bericht. — Beleg:
      Tabelle „Liefer-Punkt 2“ in §3 Umsetzung.
- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: Träger, die `§Fassung im Gate-Index` oder eine berichtigte
      Zahl zitieren, sind nachgezogen (Suchlauf, `AGENTS.md` §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| die 15 Dateien unter `harness/sensors/` und `harness/targets/` mit `## Fassung im Gate-Index` | update | Liefer-Punkt 1 |
| `harness/targets/image.md`, `workflows.md`, `sdk-pack.md`, `pin-stale.md`, `sdk-integration.md`, `examples.md`, `test-integration.md`; `harness/sensors/baseline-verify.md`; Kopfkommentar `.github/workflows/ci.yml` | update | Liefer-Punkt 2 (Befunde 2 bis 6, 8, 10, 11) |
| `harness/README.md` | update, nur falls ein Kurzsatz sich ändert | Bindung `make doc-tracked`, §1 |

Der Suchlauf (§3.13) wird beim Start gesetzt; bewegte Eigenschaften: die Menge
der Dateien mit `## Fassung im Gate-Index` und die berichtigten Zahlen.

### §3 Umsetzung (Implementer, 2026-10-05)

**Startmessung** am Parent `c38d3a6b`: `git grep -l '^## Fassung im Gate-Index' -- harness/`
liefert 15 Dateien (11 unter `sensors/`, 4 unter `targets/`), wie im Plan; die
Fassungen zusammen 24.253 Byte (*abgeleitet*: Summe der 15 je Datei gemessenen
Abschnittsgrößen, `awk` von der Überschrift bis Dateiende, `wc -c`). Der Umfang trägt ein Review: je Datei ein Abschnitt
von einem bis zwei Absätzen, die meisten ohne neue Aussage. Keine Rückführung
`in-progress → next`.

**Plan-Nachzug** (über die Tabelle oben hinaus):

| Datei | Änderungs-Art | Begründung |
|---|---|---|
| `harness/sensors/generated-sync.md` | update | Vertrag nannte nur `gen/cdc/stream/v1/*.pb.go`; die Stufe `proto-export` erzeugt auch `gen/cdc/administration/v1/` (`Dockerfile` Zeilen 66–67), der Vergleich ist byte-gleich (`cmp -s`, `tools/harness/generated-sync.sh` Zeile 115) — beim Zusammenführen der Fassung („byte-gleich“) gemessen |
| `harness/sensors/sdk-public-doc-check.md` | update | zwei neue Abschnitte `§Vorstufe der Paket-Bauten` und `§Tabellentest` aus der Fassung; die Fallliste des Tabellentests am Skript nachgemessen (die Fassung nannte nur vier Klassen, der Test trägt dazu Wortrand-, NUL-Byte- und Lesefehler-Fälle) |
| `harness/targets/sdk-pack.md` | update | „das bestehende Backend … bleibt unverändert“ (Chronik) steht als Ist-Zustand (`build-backend = "setuptools.build_meta"`, `pyproject.toml` Zeile 3) |
| `harness/README.md` | update | Bindung `make doc-tracked` zeigt auf `harness/sensors/docs-check.md` §`make doc-tracked` (Link mit Anker statt `§Fassung im Gate-Index`); kein Kurzsatz geändert |

**Liefer-Punkt 1 — je Datei.** Die Gegenprobe ist ein Skript im Scratchpad
(`gegenprobe.sh`): je Datei jeder Backtick-Span, jede Zahl und jedes Zahlwort der
alten Fassung muss danach in der Datei stehen oder auf der Streich-Liste. Endstand:
`Fehlend: 0`. Gegenbeispiel rot gesehen: dieselbe Probe gegen die Verträge **ohne**
Fassung und ohne Einarbeitung (Stand Parent) druckt `Fehlend: 29`. Die
Streich-Liste hat 7 Einträge (Spalte „gestrichen“).

| Datei | zusammengeführt | gestrichen (doppelt oder veraltet) |
|---|---|---|
| `sensors/ausgabe-kennungen-check.md` | Kennungsliste `LH-…`/`ADR-…`/`SPEC-…`/`ARC-…` in §Vertrag; „netzlos und schnell“ in §Sperren; neuer Abschnitt §Tabellentest (Fallliste, „bleibt Werkzeug“, Muster `ADR-0134` Teilfrage 2) | Gegenstand, Ausnahmen, Fail-closed, Grenzen (Raw-String, Heredoc, Rune, nachgestellter Kommentar), „Form, nicht Sinn“ — doppelt zu §Vertrag, §Ausgabe, Grenze 1–6 |
| `sensors/coverage-gate.md` | — | ganze Fassung doppelt zu §Vertrag (Fläche, Ausschluss, `db-package-lists-check.sh`) und §Kalibrierungs-Bindung (70 % → 80 %, `harness/mk/coverage.mk`) |
| `sensors/docs-check.md` | neuer Abschnitt §`make doc-tracked` (Einzel-Lauf, `--enable tracked`/`--disable`, Diagnose-Werkzeug, netzlos, `.git` im Mount) | „läuft im Bündel mit“ doppelt zu Grenze 3 |
| `sensors/fmt-check.md` | — | ganze Fassung doppelt zu §Vertrag, §Aufruf, §Exit-Codes, §Wer es aufruft, §Grenze 1, §Test; „nur die abweichende von zwei Dateien“ ist der Fall „gemischt“ der Test-Tabelle (Streich-Liste: `zwei`) |
| `sensors/generated-sync.md` | „kein `--user`-Workaround“ und „byte-gleich“ in §Vertrag | Stufe `proto-export`, Temp-Verzeichnis, Arbeitsbaum unberührt, Herkunfts-Anker — doppelt zu §Vertrag und §Bindung |
| `sensors/handbuch-public-doc-check.md` | Kennungsliste in §Vertrag; „netzlos und schnell“ in §Sperren; neuer Abschnitt §Tabellentest mit der Zahl **vier** (gemessen, siehe Widerspruch) | „fünf ausgenommenen Dateien“ (falsch, Befund 7); `*-abdeckung.md` als Sammelname — die Liste nennt die vier Namen; Rest doppelt zu §Vertrag, Grenze 1, 3, 4 |
| `sensors/kommentar-kennungen.md` | „Docker-only, netzlos“ in §Test | ganze Fassung sonst doppelt zu §Vertrag, §Kandidat, §Aufruf, §Exit-Codes, §Wer es aufruft, §Test |
| `sensors/meldungscodes-check.md` | „netzlos und schnell“ in §Sperren; neuer Abschnitt §Tabellentest (Fallliste, „bleibt Werkzeug“, Muster `ADR-0134` Teilfrage 2) | Prüfungen, Gegenstand, Fail-closed, „Mengen, nicht Sinn“ — doppelt zu §Vertrag, §Ausgabe, Grenze 1 |
| `sensors/pin-stale-all.md` | neuer Abschnitt §Tabellentest (Fallliste, Stub-`docker`, Fixture-Digests, `PROG=<Datei>`) | Muster, Vergleich, Ausgänge, Überschneidung, Grenzen, Sperren — doppelt; „wie `make pin-stale-*`“ steht als `tools/harness/pin-stale.sh` (P3 bis P6) im §Vertrag (Streich-Liste) |
| `sensors/sdk-public-doc-check.md` | „netzlos und schnell“ in §Sperren; neue Abschnitte §Vorstufe der Paket-Bauten (Voraussetzung der drei `make sdk-pack-*`, gemessen `harness/mk/sdk.mk` Zeilen 58, 84, 108) und §Tabellentest | „hängen weiterhin an der Vorgänger-Kante“ (Chronik, im Ist-Zustand neu formuliert); Rest doppelt zu §Vertrag |
| `sensors/suchlauf-nachmessen.md` | — | ganze Fassung doppelt zu §Vertrag, §Form der Zeile, §Ausgabe, §Grenze, §Test; `PLAN=<Datei>` steht als `PLAN=<Plan-Datei>` (Streich-Liste) |
| `targets/image-mutation.md` | Entfern-Befehl `docker rmi pg-change-feed-mutation:<TAG>` und „`harness/image-hash.txt`, `harness/image-hash.raw`, `:dev` unberührt“ in §Vertrag; Fallliste in §Test | Bau-Befehl, Eingabeprüfung, `$(error …)` — doppelt zu §Vertrag, §Aufruf, §Eingabeprüfung |
| `targets/schema-rollout.md` | — | ganze Fassung doppelt zu §Vertrag, §Ablauf (Nacharbeit-Tabelle nennt die vier Dateien statt `nacharbeit-*.sql`, Streich-Liste), §Erzeugnisse, §Belege |
| `targets/sdk-altserver.md` | — | ganze Fassung doppelt zu §Vertrag (Schritt-Tabelle B0 bis U, Ziel-Image, Arbeitsbaum unberührt), §Overrides, §Ausgänge, Grenze 6 |
| `targets/sdk-kompat.md` | — | ganze Fassung doppelt zu §Vertrag, §Versionen, §Overrides (`SDK_KOMPAT_NEU` mit `dist`, Streich-Liste), Grenze 5 (Netz) |

**Widerspruch, am Gegenstand entschieden.** „vier“ (Grenze 3) gegen „fünf“
(Fassung, Tabellentest) ausgenommene Dateien in `handbuch-public-doc-check.md`:
`tools/harness/handbuch-public-doc-check.sh` Zeilen 31–32 (`excluded=`) tragen vier
Namen; `grep -c '^expect 0 ".*ausgenommen"' tools/harness/run-handbuch-public-doc-check-tests.sh`
druckt `4`; `make test-handbuch-public-doc-check` Exit 0, Zeile
`run-handbuch-public-doc-check-tests: alle 41 Fälle bestanden` (2026-10-05).
Entscheidung: vier.

**Liefer-Punkt 2 — Befund-Punkt → Ausgang → Beleg** (Befund-Liste aus
`slice-harness-readme-zellen-kuerzen` §7):

| Punkt | Datei | Ausgang | Beleg |
|---|---|---|---|
| 2 `make image-cve` „kein `:latest` vorhanden“, Chronik | `targets/image.md` | nachgemessen, Chronik („löst die seit `slice-001` angekündigte … Zusage ein“, „bis zum ersten echten Release“) gestrichen | `make image-cve` 2026-10-05, Exit 0, Zeile `ghcr.io/pt9912/pg-change-feed:latest (debian 12.15) │ debian │ 0`; `gh release list`: `v0.7.0 Latest` |
| 3 `hub-description.yml` „bis zum ersten echten Release unbewiesen“ | `targets/workflows.md` | nachgemessen; der Satz zur Probe mit ungültigen Zugangsdaten (400/401) gestrichen, vom realen Erfolgslauf überholt | `gh run view 37258852180 --json jobs` (release.yml, `v0.7.0`): Job „Docker-Hub-Beschreibung synchronisieren“ `success`; `gh run list --workflow hub-description.yml`: fünf Läufe `success`, zuletzt 35483322058 |
| 4 Artefakt-Namen 0.2.x | `targets/sdk-pack.md` | nachgemessen; die Namen stehen als Form mit `<Version>` aus der Version-Datei, dazu die Messung | `ls sdks/*/dist` 2026-10-05: `PgChangeFeed.Client.0.6.1.nupkg`, `pgchangefeed-0.6.1-py3-none-any.whl`, `pgchangefeed-0.6.1.tar.gz`, `pgchangefeed-kotlin-0.6.1.jar`, `pgchangefeed-kotlin-0.6.1-sources.jar`; Version-Dateien je `0.6.1` |
| 5 Drift-Messung `make pin-stale-*` ohne Datum und Lauf | `targets/pin-stale.md` | nachgemessen, mit Datum und gedruckter Zeile; als Stand des Laufs gekennzeichnet | `make pin-stale-race`/`-pgtest`/`-dmigrate`/`-acheck` 2026-10-05, je Exit 0, je Zeile `OK …` (die alte Aussage „P3/P4/P5 zeigen Drift“ gilt nicht mehr) |
| 6 Python-HTTP-Satz in der Kotlin-Zeile | `targets/sdk-integration.md` | nachgemessen, Kopierfehler berichtigt (Kotlin- und C#-Abschnitt) | Marker `pgchangefeed-sdk-e2e:<sprache>-begin` je Runner: `run-sdk-csharp-…` Zeile 537 nur `csharp`, `run-sdk-kotlin-…` Zeile 532 nur `kotlin`, `run-sdk-python-…` Zeile 611 nur `python` |
| 7 „vier“ gegen „fünf“ | `sensors/handbuch-public-doc-check.md` | nachgemessen: vier | siehe Widerspruch oben |
| 8 Positionsverweis „`make schema-rollout`-Zeile oben“ | `targets/examples.md` | gestrichen, ersetzt durch Link auf `harness/targets/schema-rollout.md`; Chronik „seit `schema-rollout-zentrale-idempotenz-wache`“, „bleibt trotzdem bestehen“ im Ist-Zustand | „elf“ gemessen gegen `harness/targets/schema-rollout.md` (§Alles oder nichts, §Belege: „elf“) |
| 10a `baseline-verify.md` Grenze 3 nennt `harness/README.md` §Nicht behauptet | `sensors/baseline-verify.md` | nachgemessen: der Abschnitt existiert nicht (`grep -c 'Nicht behauptet' harness/README.md` → 0); Grenze 3 zeigt auf `make pin-stale-baseline` (P8). Dazu die Zahl „54 Dateien (Stand der Adoption)“ mit Lauf | `make baseline-verify` 2026-10-05, Exit 0, Zeile `baseline-verify: v6.14.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| 10b Kopfkommentar `.github/workflows/ci.yml` zählt sechs Gates, „(noch nicht existierenden) Release-Workflow“ | `.github/workflows/ci.yml` | nachgemessen; keine Aufzählung mehr (Verweis auf `GATE_CHECKS` und die Zeile `make gates`), Release-Workflows benannt; die Zeilenzahl bleibt, weil `sdk-csharp-release.yml` und `sdk-python-release.yml` „ci.yml Zeile 40f.“ zitieren (`tags-ignore:` steht weiter in Zeile 40); der Block trägt nur noch eine Kennung | `git grep -h 'GATE_CHECKS +='` → zehn Ziele; `ls .github/workflows/` trägt `release.yml` und drei `sdk-*-release.yml`; `make kommentar-kennungen DIFF=c38d3a6b` Exit 0 |
| 11 Chronik in `targets/test-integration.md` und `targets/image.md` | beide | im Ist-Zustand: test-integration „inhaltlich über den ursprünglichen MVP-Zuschnitt hinausgewachsen“, „den neuen `diagnose`-Sondermodus“, „ursprünglichen Container-Ende-Grenze“ (→ ersten), „bestehenden simulierten Neustart-Rundlauf“, „das bestehende Wecksignal … unverändert weiter“; image „weiterhin“, „jetzt als Multi-Arch“, „bestehende … unverändert“, „unveränderten … `:dev`-Pfad“; die zwei „real geprüft/bestätigt“-Sätze in `make image` sind als *übernommen* gekennzeichnet (Quelle: Kommentar am Ziel im `Makefile` Zeilen 36–45) | Lesung der Dateien; „seit slice-…“-Anker bleiben (zulässige Herkunfts-Form, `AGENTS.md` §3.7) |

Punkt 1 erledigte der Vorgänger, Punkt 9 ist Liefer-Punkt 1.

**Suchlauf** (`AGENTS.md` §3.13; Parent `c38d3a6b`, Stand `diff` nach diesem Slice;
Suchraum der ganze Baum ohne `docs/reviews/`, `done/`, `.harness/`):

```suchlauf
c38d3a6b 18 -F 'Fassung im Gate-Index' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 2 -F 'Fassung im Gate-Index' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 3 -F 'fünf ausgenommenen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 2 -F 'fünf ausgenommenen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 1 -F 'Nicht behauptet' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 0 -F 'Nicht behauptet' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 5 -F 'PgChangeFeed.Client.0.2' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 4 -F 'PgChangeFeed.Client.0.2' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 7 -F 'Python-HTTP-Abschnitt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 5 -F 'Python-HTTP-Abschnitt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 4 -F 'ersten echten Release' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 2 -F 'ersten echten Release' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 3 -F 'noch nicht existierenden' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 2 -F 'noch nicht existierenden' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
c38d3a6b 1 -F 'zeigen echten Drift' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
diff 0 -F 'zeigen echten Drift' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness'
```

Gefunden und nachgezogen: die 15 Dateien, `harness/README.md` (Bindung
`make doc-tracked`) und die Dateien der Befund-Liste. Verbleibende Treffer:
`Fassung im Gate-Index` in `open/slice-abgeleitete-dokumente-vorlagen-nachzug.md`
(Abgrenzung, die diesen Slice zitiert, zutreffend) und in einer Evidence-Datei des
Registers (Record); `fünf ausgenommenen` in `ADR-0143` (`Accepted`, eingefroren,
§1) und derselben Evidence-Datei; `Python-HTTP-Abschnitt` in Register-Records.

**Gemeldete Träger in fremden Dateien** (`AGENTS.md` §3.13, Frist: Closure dieses
Slice; außerhalb von §3 dieses Plans, deshalb nicht still mitgeändert):

- `spec/pflichtenheft.md` Zeilen 485, 491, 498 nennen die Artefakt-Namen 0.2.1/0.2.2
  (Befund 4 in der Spec; Rang 2, Änderung ist Spec-Arbeit).
- `harness/mk/sdk.mk` Zeile 55 und `sdks/csharp/Dockerfile` Zeile 63
  (Kommentare: `PgChangeFeed.Client.0.2.1.nupkg`).
- `.github/workflows/image-scan.yml` Zeile 24 und `.github/workflows/release.yml`
  Zeile 28 (Kommentare: „bis zum ersten echten Release“; Befund 2 und 3 im Kommentar).
- `tools/harness/run-sdk-kotlin-integration-tests.sh` Zeile 553 (Kommentar
  „Python-HTTP-Abschnitt der Folge-Slices“).

Nicht gefunden: kein weiterer Träger von `§Nicht behauptet`, `zeigen echten Drift`
oder der Positionsverweis „Zeile oben“.

**Außerhalb der Befund-Liste bemerkt, nicht geändert:** `sensors/ausgabe-kennungen-check.md`
§Grenze trägt die Nummer 6 zweimal; `sensors/generated-sync.md` §Vertrag trägt
Chronik („Bis `slice-generated-sync-tar-export` lief hier …“, „Vorher …“, „jetzt“);
`sensors/coverage-gate.md` §Vertrag „seit `slice-097`“ mitten im Satz.

**Sensoren dieses Laufs:** `make docs-check` Exit 0 (`d-check: 1737 Datei(en)
geprüft, 0 Befund(e)`); `make kommentar-kennungen DIFF=c38d3a6b` Exit 0 (kein
Kandidat); die zwei neuen Anker-Links (`#make-doc-tracked` in `harness/README.md`,
`#make-pin-stale-baseline` in `baseline-verify.md`) prüft `anchors`: an einer
Scratchpad-Kopie von `0cdeb963` mit je einem angehängten `X` am Anker endet d-check
(`--disable tracked`, die Kopie trägt kein `.git`) mit Exit 1 und zwei Befunden
`anchor-missing`, ohne Mutation mit Exit 0; `make test-handbuch-public-doc-check` Exit 0; `make fmt-check` entfällt
(der Diff trägt keine Go-Datei). `make gates` und `make suchlauf-nachmessen` nach dem
letzten Commit: Bericht der Sitzung.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1); läuft `slice-baseline-6-14-0-dokumente-nachziehen` zuerst, liest
dieser Slice dessen Änderung an `harness/targets/pin-stale.md` mit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Zusammenführung
  trägt in einem Review nicht alle 15 Dateien — dann Teilung nach
  `harness/sensors/` und `harness/targets/`.
- `in-progress` → `open` (blockiert): eine Befund-Aussage ist ohne Netz oder
  Release nicht messbar und lässt sich nicht als übernommen kennzeichnen —
  Frage an den Auftraggeber.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, und die Closure-Notiz in §7 trägt einen Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang.

- **Zusammenführen verliert eine Zusage:** ein Satz der Fassung steht nur dort und
  fällt beim Streichen weg. Gegenmaßnahme: Abgleich der Eigennamen (Make-Ziele,
  Pfade, Zahlen) der Fassung gegen den Endstand der Datei. — **Ausgang:** offen
  bis zur Closure.
- **Eine Aussage ist am Host nicht messbar** (Release-Stand, Registry): dann
  *übernommen* mit Quelle statt Messung. — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Ging der Gegenstand an einen anderen Slice oder entfiel er,
trägt diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund.

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Der Abschnitt selbst entfällt nie.**

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt eine Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`harness/`, `.github/workflows/ci.yml`) liegen in ihr.

**Vorgelagert — offene Beobachtungen sichten:** beim Start durch den Implementer
(Register `docs/plan/planning/observations/BEO-PGC/`). Bekannt bei der Anlage:
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, Anlass F-1) und
`BEO-PGC/bindung-spalte-uneinheitlich-tief` (offen, 1×; berührt, wenn die
Bindung `make doc-tracked` umgestellt wird).

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
