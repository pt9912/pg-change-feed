# Review-Report: release-doku-releasing — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/release-doku-releasing.md`), `ADR-0051`
(derivativ, kein eigener Entscheidungspunkt) und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe
(Modul 11).

**Gegenstand:** Commit `b019b058` ("docs(user): releasing.md für den realen
Release-Prozess (`ADR-0051`)"), Slice `release-doku-releasing`, **letzter**
Slice der Welle `welle-release-pipeline-adr-0051`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln; HIGH-Nachbarform „Zitat nennt die
falsche Stelle" seit welle-d-check gefasst, im Vorgänger-Slice
`release-hub-description` zweimal angewendet).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-doku-releasing.md` (Slice-Plan, §1–§3, §6)
- `ADR-0051` (Accepted) — derivativ, Summe der vier vorangehenden Slices
- `AGENTS.md` §3.7 (Kommentar-Disziplin), §3.10 (Post-Push-Verifikationspflicht),
  §3.11 (kein host-lokaler absoluter Pfad), §3.12 (Herkunft von Aussagen),
  §4 (kein behauptetes Gate ohne Deckung in `harness/README.md`)
- `harness/conventions.md` (MR-000 ID-Schema)
- `harness/sensors/docs-check.md` (Grenzen des `hostpaths`-/`tracked`-Moduls)
- Nachbar-Reviews `docs/reviews/review-slice-release-hub-description.md` und
  `docs/reviews/review-slice-release-version-und-workflow.md` (Präzedenz für
  die HIGH-Klasse „Zitat nennt die falsche Stelle" und für den
  `release-tag-info.sh`-Testbarkeits-Standard)

**Eigenständig nachgefahrene Prüfungen (nicht vom Implementer-Bericht
übernommen):**

- Alle sechs im Dokument genannten Tag-Beispiele erneut real gegen
  `tools/harness/release-tag-info.sh` ausgeführt: `v0.1.0` →
  `version=0.1.0 latest=true`, `v1.2.0-rc.1` → `version=1.2.0-rc.1
  latest=false`, `v2.0.0+build.5` → `version=2.0.0+build.5 latest=true`
  (alle drei Exit 0); `v1.0.0-01`, `v1.0`, `1.0.0` → jeweils Exit 1 mit dem
  im Dokument behaupteten Fehlerbild. Deckt sich exakt mit
  `docs/user/releasing.md` §3.
- `docs/user/release.yml`, `hub-description.yml`, `image-scan.yml`,
  `upstream-drift.yml` vollständig gelesen (nicht nur die vom Implementer
  zitierten Ausschnitte) und Satz für Satz gegen `docs/user/releasing.md`
  §4/§5 abgeglichen: Trigger (`push: tags: ['v*']`), Reihenfolge
  (Tag-Validierung → `version.md`-Abgleich → `make image
  VERSION=…`/`LATEST=…` → GitHub-Release mit Digest → `hub-description`-Job
  über `needs: release`/`uses:`/`secrets: inherit`), Docker-Hub-Token-Scope-
  Hinweis und die beiden advisory-Workflows stimmen jeweils überein.
- `ci.yml`/`e2e.yml` gelesen: beide tragen `tags-ignore: ['**']` — die
  Behauptung „kein Doppellauf" in `releasing.md` §4 ist zutreffend.
- Alle relativen Links aus `docs/user/releasing.md` heraus rechnerisch
  aufgelöst (`os.path.normpath` von `docs/user/` aus) — alle fünf Ziele
  (`version.md`, `../../.github/workflows/{release,image-scan,upstream-drift}.yml`,
  `../../README.md`) existieren real im Baum.
- Realer `make docs-check`-Lauf auf dem Arbeitsbaum nach dem Commit,
  Exit-Code ungepiped geprüft: `0` (775 Datei(en) geprüft, 0 Befund(e) —
  volle Modul-Liste inkl. `hostpaths`, `links`, `anchors`, `tracked`,
  `structure`). Bestätigt unabhängig, dass die beiden Link-Umstellungen in
  `harness/README.md`/`AGENTS.md` keine Regel auslösen.
- Realer `make commit-traceability RANGE=b019b058~1..b019b058`-Lauf, Exit
  `0` — Commit-Betreff trägt `ADR-0051`, keine verbotene `SPEC-*`/`ARC-*`-
  Kennung im Betreff.
- Kein vollständiger `make gates`-Lauf eigenständig wiederholt
  (Coverage-Gate/`a-check`/`generated-sync`/`baseline-verify` sind für
  einen reinen Doku-Diff ohne Go-/Konfigurationsänderung nicht
  aussagekräftig neu zu prüfen) — die beiden diff-relevanten Teilgates
  (`docs-check`, `commit-traceability`) sind oben eigenständig bestätigt;
  die pauschale „`make gates` grün"-Aussage bleibt in diesem Umfang eine
  vom Implementer übernommene, nicht vollständig eigenständig
  nachgefahrene Behauptung (fällt strukturell in den Verifier-Kontext).

---

## Findings

### F-1 — Falsch zitierte Abschnitts-Stelle in `docs/user/releasing.md` §2

- `kategorie`: HIGH
- `quelle`: Hard-Rule „Beleg trägt seinen Satz nicht" / Nachbarform „Zitat
  nennt die falsche Stelle" (`.harness/skills/reviewer.md` §Klassifikation
  HIGH, `BEO-PGC/zitat-nennt-die-falsche-stelle`)
- `pfad`: `docs/user/releasing.md:32`
- `befund`: Satz: „Diese Datei ist die Quelle der Wahrheit für die
  Version: der Release-Workflow (§3) gleicht sie gegen den gesetzten
  Git-Tag ab und bricht bei jeder Abweichung ab, bevor irgendein Login,
  Build oder Push läuft." Der zitierte Abschnitt **§3 „Einen Release
  auslösen"** (Zeilen 41–60) beschreibt ausschließlich die Tag-Form und
  ihre SemVer-2.0-Validierung über `tools/harness/release-tag-info.sh` —
  er erwähnt `docs/user/version.md` und den Abgleich gegen den Tag mit
  keinem Wort. Der Vergleich gegen `version.md` steht tatsächlich in
  **§4 „Was beim Release automatisch passiert"**, Punkt 2: „`docs/user/
  version.md` gegen den Tag abgleichen — Abbruch bei Abweichung (§2)." —
  jener Punkt zitiert korrekt zurück auf §2, aber §2 selbst zeigt mit
  „(§3)" auf die falsche Stelle. Ich habe beide Abschnitte im Original
  aufgeschlagen, nicht nur die Zusammenfassung: §3 enthält keine
  `version.md`-Erwähnung. Der genannte Anker trägt die Aussage nicht.
  (Der Nachbarsatz zwei Zeilen darunter, „… dann den passenden Tag setzen
  (§3)" auf Zeile 37, zitiert dagegen korrekt — §3 beschreibt tatsächlich
  das Setzen des Tags.)
- `verifizierbar`: nein — kein Gate prüft Zitat-Ziel-Treue innerhalb eines
  Doku-Trägers; nur Aufschlagen der zitierten Stelle zeigt die Abweichung.
- `klasse`: Zitat nennt die falsche Stelle

### F-2 — `docs/user/releasing.md` fehlt in README.md-Doku-Übersicht

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `README.md:22-27`
- `befund`: Die „Siehe:"-Liste im Projekt-`README.md` verlinkt
  `docs/user/benutzerhandbuch.md`, aber nicht die neue
  `docs/user/releasing.md` — ein Leser, der nur das Top-Level-`README.md`
  liest, findet den Release-Prozess nicht darüber. Dieselbe Lücke besteht
  bereits für `docs/user/version.md`, `bench-abdeckung.md`,
  `e2e-abdeckung.md` und `ci-matrix-abdeckung.md` (keine davon ist dort
  verlinkt) — die Liste war schon vor diesem Commit nicht vollständig für
  `docs/user/*`, dieser Slice führt also kein neues Muster ein, sondern
  setzt ein bestehendes fort. Kein Blocker, da der Plan (§1) diese Datei
  nicht als Änderungsziel nennt.
- `verifizierbar`: nein.
- `klasse`: Neue Doku-Datei nicht in Top-Level-README verlinkt

## Negativbefunde

- geprüft, ohne Befund: alle sechs Tag-Beispiele — eigenständig erneut
  gegen `tools/harness/release-tag-info.sh` gefahren, Ergebnis deckt sich
  exakt mit `docs/user/releasing.md` §3.
- geprüft, ohne Befund: Existenz und Pfad-Auflösung aller in
  `docs/user/releasing.md` genannten Dateien/Workflows/Targets
  (`docs/user/version.md`, `.github/workflows/{release,image-scan,
  upstream-drift,hub-description}.yml`, `README.md`,
  `make test-release-tag-info`, `make image-cve`, `make image`) — alle
  real vorhanden, alle relativen Links rechnerisch aufgelöst.
- geprüft, ohne Befund: inhaltliche Übereinstimmung von §4/§5 mit dem
  tatsächlichen Inhalt von `release.yml`, `hub-description.yml`,
  `image-scan.yml`, `upstream-drift.yml` — Trigger, Reihenfolge,
  Job-Verdrahtung (`needs`/`uses`/`secrets: inherit`) und die advisory-
  Eigenschaft der beiden Begleit-Workflows stimmen jeweils überein; keine
  Abweichung zwischen behauptetem und tatsächlichem Workflow-Verhalten
  gefunden außer F-1 oben (Abschnitts-internes Zitat, keine Abweichung
  zum Workflow-Inhalt selbst).
- geprüft, ohne Befund: Docker-Hub-Token-Scope-Hinweis
  (`read/write/delete`, `403 Forbidden` bei nur `read/write`, Verweis auf
  `packaging/dockerhub/README.md` §Transport im Schwester-Repo d-check) —
  wortgleich mit dem bereits im Vorgänger-Slice reviewten Kopfkommentar
  von `.github/workflows/hub-description.yml:37-41`, keine Drift.
- geprüft, ohne Befund: `ci.yml`/`e2e.yml` `tags-ignore: ['**']` — beide
  Dateien real gelesen, Behauptung „kein Doppellauf" zutreffend.
- geprüft, ohne Befund: `AGENTS.md` §3.11 (kein host-lokaler absoluter
  Pfad) — die Schwester-Repo-Zitation läuft in Hausform
  („Schwester-Repo d-check … `packaging/dockerhub/README.md` §Transport"),
  kein Host-Pfad; eigener `make docs-check`-Lauf (`hostpaths`-Modul im
  Bündel) bestätigt 0 Befunde.
- geprüft, ohne Befund: Umstellung von `harness/README.md`
  Source-Precedence-Zeile 6 und `AGENTS.md` §2 Zeile 6 auf einen echten
  Verzeichnis-Link — analog zur bereits bestehenden Zeile 4
  (`[docs/user/](../docs/user/)` bzw. `[docs/user/](docs/user/)`, gleiches
  Muster wie `[docs/plan/adr/](../docs/plan/adr/)`); die im entfernten
  `<!-- d-check:ignore -->`-Kommentar genannte Bedingung („im frischen
  Repo selten vorhanden") ist mit sieben Dateien unter `docs/user/`
  tatsächlich widerlegt. Eigener `make docs-check`-Lauf nach der Änderung:
  0 Befunde — kein durch den Ignore-Kommentar vormals unterdrückter Fund
  wird durch die Umstellung auf einen echten Link sichtbar.
- geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentar-Disziplin) in den
  beiden geänderten Zeilen — der entfernte HTML-Kommentar wird nicht durch
  einen neuen, potenziell fehlerhaften Kommentar ersetzt; die verbleibende
  Zeile trägt keinen Kommentar mehr, der zu prüfen wäre.
- geprüft, ohne Befund: `AGENTS.md` §4 (kein behauptetes Gate ohne
  Deckung) — `docs/user/releasing.md` nennt `release.yml`,
  `image-scan.yml`, `upstream-drift.yml` durchgängig konsistent mit ihrer
  in `harness/README.md` §Sensors dokumentierten „kein Gate"-Einordnung;
  keine der drei Dateien wird als Gate bezeichnet. `make
  test-release-tag-info` wird nicht als Gate behauptet, sondern nur als
  netzlos ausführbarer Test benannt — der Umstand, dass dieses (und
  `make test-dockerhub-token`) in `harness/README.md` §Sensors/Werkzeuge
  nicht als eigene Zeile geführt wird, ist keine Neuerung dieses Diffs
  (bereits in den beiden Vorgänger-Slices `release-version-und-workflow`/
  `release-hub-description` so eingeführt und dort ohne Beanstandung
  reviewt/verifiziert).
- geprüft, ohne Befund: Plan-Selbstkonsistenz — die AGENTS.md-Änderung war
  im ursprünglichen Slice-Plan-§3 nicht vorab benannt; der Commit trägt im
  selben Zug einen Plan-Nachzug (§3-Tabellenzeile „`AGENTS.md` §2 Zeile 6
  | update (Plan-Nachzug)"), bevor die DoD-Zeile auf `[x]` gesetzt wird —
  entspricht der etablierten Praxis, einen Nachzug im Plan selbst
  sichtbar zu machen statt ihn nur im Diff verschwinden zu lassen.
- geprüft, ohne Befund: Plan §1 „Ausdrücklich NICHT in diesem Slice" —
  weder eine Behauptung eines bereits erfolgten echten Release noch eine
  Änderung an `docs/user/benutzerhandbuch.md` im Diff enthalten.
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt `ADR-0051`;
  eigener `make commit-traceability`-Lauf über die Range Exit 0.
- geprüft, ohne Befund: `docs/user/releasing.md` Version-/Stand-Kopf und
  `### Änderungshistorie`-Tabellenform — strukturell analog
  `benutzerhandbuch.md` (gleiche Spaltenköpfe `Version | Datum |
  Änderung`); das Fehlen eines eigenen „Software-Version"-Felds ist
  sachlich begründet (der Release-Prozess ist nicht an eine einzelne
  Software-Version gebunden) und keine DoD-Abweichung, da die DoD-Zeile
  nur „analog" verlangt, nicht identisch.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zitat nennt die falsche Stelle · Neue
Doku-Datei nicht in Top-Level-README verlinkt.

## Verdikt

**Merge-blockierend:** ja — F-1 ist eine weitere Instanz der bereits
etablierten HIGH-Klasse „Zitat nennt die falsche Stelle" (nach zwei
Vorkommen im unmittelbaren Vorgänger-Slice `release-hub-description` nun
ein drittes/viertes Vorkommen innerhalb derselben Welle). Die zugrunde
liegende Aussage selbst ist sachlich zutreffend — nur der Abschnitts-
Anker zeigt auf die falsche Stelle (§3 statt §4). Mechanisch billig zu
beheben: „(§3)" auf Zeile 32 durch „(§4)" ersetzen. F-2 (INFO) ist kein
Blocker und setzt lediglich eine bereits vor diesem Slice bestehende
Lücke fort, ohne sie zu verschärfen.

**Übergabe:** Da F-1 (HIGH) eine Fixrunde auslöst, ziehe ich die
DoD-Checkbox „Review durchgeführt" im Slice-Plan **nicht** nach
(Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift nur bei 0 HIGH).
Die Finding-Klasse „Zitat nennt die falsche Stelle" geht bei
Slice-Closure ins Beobachtungs-Register (`BEO-PGC/zitat-nennt-die-
falsche-stelle`, jetzt mit zwei weiteren Fundstellen aus derselben Welle
— dieselbe Klasse trat in `release-hub-description` bereits zweimal auf,
hier ein weiteres Mal; die Häufung innerhalb einer einzigen Welle ist für
die Steering-Loop-Zählung selbst relevant und wird bei Wellen-Closure
mitgeführt). Dieser Report ersetzt keine Verifikation gegen die DoD —
das bleibt Verifier-Aufgabe (Modul 11).
