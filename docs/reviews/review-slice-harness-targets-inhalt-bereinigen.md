# Review-Report: slice-harness-targets-inhalt-bereinigen — 2026-10-05

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (Maintainability), nicht gegen die DoD

**Gegenstand:** Diff `c38d3a6b..e6e336bc`, die Commits 0cdeb963 (Sensor- und Target-Dateien, `harness/README.md`, Kopfkommentar `.github/workflows/ci.yml`) und e6e336bc (Plan: Anker-Gegenprobe)

**Skill:** `.harness/skills/reviewer.md` @ 7fff2cac
**Modell:** Opus 5.5 (`claude-opus-5-5`) · **Datum:** 2026-10-05

> **Zitier-Form.** Dieser Report friert ein; er zitiert Kennungen, nicht
> Adressen (`slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag +
> Pfad in Inline-Code).

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-targets-inhalt-bereinigen` (§1 Abgrenzung, §2 Liefer-Punkte 1 und 2, §3 Umsetzung mit Tabellen, Widerspruch, Suchlauf und gemeldeten Trägern, §6 Risiken)
- Befund-Liste in §7 von `slice-harness-readme-zellen-kuerzen`
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.10, §3.12, §3.13
- keine `LH-*`-ID berührt

**Eigene Proben dieses Laufs (gemessen, 2026-10-05, Stand e6e336bc):**

1. **Liefer-Punkt 1, Satz für Satz** an sechs Dateien, alte Fassung aus `git show c38d3a6b:<pfad>`: ganz gestrichen `harness/sensors/coverage-gate.md`, `harness/targets/sdk-kompat.md`, `harness/targets/schema-rollout.md`; eingearbeitet `harness/sensors/handbuch-public-doc-check.md`, `harness/sensors/generated-sync.md`, `harness/targets/image-mutation.md`. Je Satz der Fassung die Aussage im verbleibenden Vertrag gesucht und gelesen (nicht nur Eigennamen): Fläche, Ausschluss, Listen-Gleichheit und Rotfärbung vor dem Bau (coverage-gate §Vertrag, §Kalibrierungs-Bindung); A1 bis A5, Gast über Basis-Konstruktoren, Bau nur übersetzt und `docker run` statt `RUN`, `SDK_KOMPAT_NEU`/`REGISTRY_VERSION`, Netz und Gate-Aufnahme nur per ADR (sdk-kompat §Vertrag, §Versionen, §Overrides, §Ausgänge, Grenze 5); Precheck, Wache, View-Vorlauf, Report und Rollback-Artefakt, vier Nacharbeit-Schritte, zwei Images, stdin, Artefakt-Verzeichnis, Idempotenz, sechs Läufe (schema-rollout); Muster, Ausnahmen, Exit 2 für neue Datei, Lese-Hälfte beim Reviewer (handbuch); byte-gleich per `cmp`, beide `gen/`-Bäume, Temp-Verzeichnis (generated-sync, gegen `tools/harness/generated-sync.sh` und `Dockerfile` gelesen); Bau- und Entfern-Befehl, `kein -f`/`kein prune`, `$(error …)`, Testfälle (image-mutation). Keine Aussage ging verloren; die Streich-Liste (7 Einträge) ist an den genannten Ersatzstellen gelesen und trägt.
2. **„vier gegen fünf“ am Gegenstand:** `excluded=` in `tools/harness/handbuch-public-doc-check.sh` trägt vier Namen; `docs/user/` führt genau diese vier `*-abdeckung.md`; der Tabellentest trägt vier Fälle `expect 0 "… ausgenommen"` und einen fünften verwandten Fall `check 2 "fehlende ausgenommen-genannte Datei"` (Exit 2, kein Ausnahme-Fall). `make test-handbuch-public-doc-check` Exit 0, Zeile `run-handbuch-public-doc-check-tests: alle 41 Fälle bestanden`. Die Entscheidung „vier“ trägt.
3. **Liefer-Punkt 2, Stichproben:** `make pin-stale-acheck` Exit 0, Zeile `OK          A_CHECK_IMAGE (ghcr.io/pt9912/a-check:latest) == sha256:e8208764…`; die Aussage in `harness/targets/pin-stale.md` ist mit Datum gebunden und als „Stand dieses Laufs“ gekennzeichnet. `make image-cve`: „Exit 1 bei mindestens einem Befund“ und „`GHCR_USERNAME`/`GHCR_PASSWORD` optional“ gegen das Ziel im `Makefile` gelesen (`--exit-code 1`, `$(if $(GHCR_USERNAME),…)`); die zwei „real geprüft“-Sätze in `make image` stehen als *übernommen* mit Quelle, der Makefile-Kommentar trägt beide. Punkt 6: die Marker `pgchangefeed-sdk-e2e:<sprache>-begin` je Runner gelesen (csharp, kotlin, python je nur die eigene Sprache); Kotlin- und C#-Zeile in `harness/targets/sdk-integration.md` sagen das jetzt.
4. **`.github/workflows/ci.yml`:** `git diff c38d3a6b..e6e336bc -- .github` ändert ausschließlich Kommentarzeilen (Filter auf Nicht-Kommentarzeilen: leer); `tags-ignore:` steht weiter in Zeile 40, die Zitate „ci.yml Zeile 40f.“ in `sdk-csharp-release.yml` und `sdk-python-release.yml` lösen auf; `harness/sensors/gates.md` existiert; zehn `GATE_CHECKS +=` in den Make-Fragmenten. `make kommentar-kennungen DIFF=c38d3a6b` und `make kommentar-kennungen PATHS=.github/workflows/ci.yml` je Exit 0 ohne Kandidat.
5. **README-Link** `make doc-tracked` → `sensors/docs-check.md#make-doc-tracked`: der Abschnitt §`make doc-tracked` existiert in `harness/sensors/docs-check.md`; Zellenlänge und Anker prüft `make docs-check` im Gate-Lauf.
6. `make suchlauf-nachmessen PLAN=<Plan>` Exit 0, `suchlauf-nachmessen: 16 Zeilen stimmen`. Dazu eigene Suche nach den Befund-Aussagen in anderer Wortform (`unbewiesen`, `0.2.x`-Artefaktnamen aller drei Sprachen) über den Baum ohne `docs/reviews/`, `done/`, `.harness/`.
7. `git grep -n 'Fassung im Gate-Index'`: Treffer nur in Records (`done/`, `docs/reviews/`, Register-Evidence), im Plan selbst und in `open/slice-abgeleitete-dokumente-vorlagen-nachzug.md` (Abgrenzung) — keine verwaiste Verweisstelle in `harness/`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Liste der gemeldeten fremden Träger ist unvollständig: `.github/workflows/hub-description.yml` trägt im Kopfkommentar genau die Aussage von Befund 3, die dieser Slice in `harness/targets/workflows.md` als überholt gestrichen hat (Probe mit ungültigen Zugangsdaten, „der volle Erfolgspfad bleibt bis zum ersten echten Lauf … unbewiesen“); das Suchmuster `ersten echten Release` trifft „ersten echten Lauf“ nicht, und der Plan sagt „Nicht gefunden“. Dieselbe Klasse in `sdk-csharp-release.yml` und `sdk-python-release.yml` („Vorhandensein bleibt bis zum ersten echten Tag-Push unbewiesen“), während der Gate-Index die Secrets durch grüne Läufe belegt nennt. Failure-Szenario: die Closure gilt die Träger-Meldung als vollständig, und der Kommentar am Workflow widerspricht weiter dem berichtigten Vertrag. | `AGENTS.md` §3.13 (Suchform: Beschreibung samt Hedge; Meldung fremder Träger) | `.github/workflows/hub-description.yml` · `zum ersten echten Lauf mit realen Repository-Secrets unbewiesen` | nein — kein Sensor liest Träger außerhalb des Diffs; `git grep -n unbewiesen -- .github` zeigt die Stellen | Träger-Meldung unvollständig (Suchmuster trifft die Wortform nicht) |
| F-2 | LOW | Der neue Abschnitt §Tabellentest verweist auf „Grenze 7“; diese Nummer gibt es nur, weil §Grenze die 6 zweimal trägt, was der Plan als „bemerkt, nicht geändert“ führt. Wer die Doppelnummer berichtigt, macht aus der Wächter-Logik-Grenze die 8, und der neue Verweis zeigt dann auf nichts. | Skill HIGH „Beleg trägt seinen Satz nicht“ (Nachbar-Form Zitat; hier trägt der Zeiger heute, gebunden an einen bekannten Defekt) | `harness/sensors/ausgabe-kennungen-check.md` · `Teilfrage 2; Grenze 7).` | ja — Aufschlagen der Stelle (`^6\.` zweimal in §Grenze) | Abschnitts-Zeiger auf fehlerhafte Nummerierung |
| F-3 | INFO | §Tabellentest nennt eine Fallliste, die die Fälle Unterverzeichnis, Nicht-`.md`-Datei, fehlende geprüfte Datei, NUL-Byte und Lesefehler nicht nennt, und trägt die Zahl „41 Fälle“ als datierte Messung im Vertrag; in `harness/sensors/sdk-public-doc-check.md` hat derselbe Lauf die Fallliste am Skript nachgemessen. Die Zahl trägt ihren Lauf (§3.12 erfüllt), bewegt sich aber mit jedem neuen Fall. | `AGENTS.md` §3.12 | `harness/sensors/handbuch-public-doc-check.md` · `alle 41 Fälle bestanden` | ja — `make test-handbuch-public-doc-check` | Fallliste im Vertrag kürzer als der Test |
| F-4 | INFO | In bearbeiteten Dateien bleibt Chronik-Sprache außerhalb der Befund-Liste stehen: `harness/targets/sdk-integration.md` (alle drei Abschnitte „baut die neue `integration`-Docker-Stufe“, „Pins aus der bestehenden `build`-Stufe“) und, vom Implementer selbst benannt, `harness/sensors/generated-sync.md` §Vertrag. Nicht Gegenstand von Liefer-Punkt 2; Hinweis an den Planner. | `AGENTS.md` §3.7 | `harness/targets/sdk-integration.md` · `Pins aus der bestehenden` | nein | Chronik-Sprache in Harness-Vertrag |
| F-5 | INFO | Implementer- und Reviewer-Lauf teilen dasselbe Scratchpad-Verzeichnis der Sitzung; dort liegen `gegenprobe.sh`, Gate-Logs und Commit-Messages des Implementers. Dieser Lauf hat keines davon als Beleg benutzt (alle Proben oben eigen gefahren), die Kontext-Trennung hängt aber an dieser Disziplin. Hinweis an den Auftraggeber. | `v6.14.0` · `regelwerk/modul-08-agentenrollen.md` §Rollen-Regeln | Scratchpad der Sitzung · `gegenprobe.sh` | nein | Rollen-Artefakte im geteilten Scratchpad |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Liefer-Punkt 1, sechs Dateien Satz für Satz (coverage-gate, sdk-kompat, schema-rollout, handbuch-public-doc-check, generated-sync, image-mutation) | geprüft, ohne Befund — kein Inhalt verloren, Streich-Liste trägt |
| übrige neun Dateien mit entfernter Fassung (ausgabe-kennungen-check, docs-check, fmt-check, kommentar-kennungen, meldungscodes-check, pin-stale-all, sdk-public-doc-check, suchlauf-nachmessen, sdk-altserver) | geprüft per Stichwort-Abgleich der Fassungs-Aussagen gegen den Vertrag und Lesen der neuen Abschnitte; Grenz-Verweise (Grenze 1, 3, 4, 5, 6) aufgeschlagen — ohne Befund außer F-2, F-3 |
| Widerspruch „vier gegen fünf“ | geprüft am Skript, am Verzeichnis und am Tabellentest, ohne Befund |
| Liefer-Punkt 2, Tabelle Punkt → Ausgang → Beleg | geprüft (Punkte 2, 3, 5, 6, 8, 10a, 10b, 11 gelesen; 5 und 2 nachgemessen bzw. am Makefile gelesen), ohne Befund; gemessene Aussagen tragen Datum und Zeile, übernommene sind gekennzeichnet |
| `.github/workflows/ci.yml` | geprüft, nur Kommentarzeilen geändert, keine strukturelle Workflow-Änderung (§3.10 nicht ausgelöst), Zeilen-Zitate lösen auf, ohne Befund |
| `harness/README.md` | geprüft, eine Bindung-Zelle geändert, Anker existiert, ohne Befund |
| Abschnittsstruktur, verwaiste Anker, `Fassung im Gate-Index`-Reste in `harness/` | geprüft, ohne Befund |
| Suchlauf und gemeldete fremde Träger | geprüft — Zahlen stimmen (16 Zeilen); Vollständigkeit siehe F-1 |
| Docker-only, Suppression, Traceability der Commits | geprüft, ohne Befund (Commits nennen `ADR-0051`/`ADR-0083`) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Träger-Meldung unvollständig (Suchmuster trifft die Wortform nicht) · Abschnitts-Zeiger auf fehlerhafte Nummerierung · Fallliste im Vertrag kürzer als der Test · Chronik-Sprache in Harness-Vertrag · Rollen-Artefakte im geteilten Scratchpad

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 betrifft fremde Dateien
und geht als Meldung an den Planner (Frist: Closure dieses Slice, `AGENTS.md`
§3.13); F-2 nimmt der Implementer an oder begründet es, ohne Fixrunden-Pflicht;
F-3 bis F-5 sind Hinweise.

**Übergabe:** Findings an den Implementer und den Planner; die Finding-Klassen
gehen in die Slice-Closure §7. Die DoD-Zeile „Review durchgeführt“ ist im selben
Commit nachgezogen (Skill §DoD-Checkbox-Nachzug ohne Fixrunde). Dieser Report
ersetzt keine Verifikation.
