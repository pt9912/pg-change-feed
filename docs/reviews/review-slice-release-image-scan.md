# Review-Report: release-image-scan — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/release-image-scan.md`), `ADR-0051`
Entscheidung 6 und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `1141786c` (Parent `ea631152`), Slice
`release-image-scan`, Welle `welle-release-pipeline-adr-0051`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-image-scan.md` (Slice-Plan)
- `docs/plan/planning/welle-release-pipeline-adr-0051.md` (Welle-Datei, §5/§6)
- `ADR-0051` (Accepted) — Entscheidung 6, Verglichene Alternativen 4,
  Fitness Function
- `docs/plan/planning/done/release-version-und-workflow.md` und dessen
  Review (`docs/reviews/review-slice-release-version-und-workflow.md`) —
  Vorgänger-Slice derselben Welle, als Sorgfalts-Referenz
- `AGENTS.md` §3.1 (Docker-only), §3.5 (ADR-Immutabilität), §3.7
  (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9 (Exit-Code-Disziplin),
  §3.12 (Herkunft von Aussagen)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)

---

## Findings

### F-1 — `image-scan.yml` authentifiziert nicht gegen GHCR; CVE-Scan kann nach dem ersten Release strukturell anders (und dauerhaft) scheitern, als der Slice-Plan annimmt

- `kategorie`: MEDIUM
- `quelle`: `ADR-0051` Entscheidung 6 (Scan gegen „das publizierte GHCR-
  `:latest`-Image") / Maintainability
- `pfad`: `Makefile:49-50` (`image-cve`-Rezept, kein Registry-Login);
  `.github/workflows/image-scan.yml` (kein `docker/login-action`- oder
  `docker login ghcr.io`-Schritt) gegen `docs/plan/planning/in-progress/release-image-scan.md:113-120`
  (§6 Risiko-Ausgang)
- `befund`: GitHub Container Registry legt ein neu über `GITHUB_TOKEN`
  gepushtes Package standardmäßig als **privat** an; die
  Repository-Sichtbarkeit vererbt sich nicht automatisch, sondern
  verlangt einen einmaligen, manuellen Schritt in den
  Paket-Einstellungen (oder eine „Inherit access"-Verknüpfung). Weder
  `image-scan.yml` noch `make image-cve` enthalten einen
  Registry-Login-Schritt — der `docker run`-Aufruf zieht das Image
  ausschließlich anonym. Der Slice-Plan führt den heute real
  reproduzierten `DENIED`-Fehler ausschließlich auf das fehlende Image
  zurück und schließt daraus: „löst sich strukturell mit dem ersten
  echten Release, kein weiterer Implementierungsschritt nötig" (§6,
  Zeile 119-120). Bleibt das künftige Package privat — was der
  GHCR-Default ist, sofern niemand es manuell ändert —, bliebe der exakt
  gleiche `DENIED`-Fehler nach dem ersten Release bestehen, nur aus einem
  anderen Grund (fehlende Auth statt fehlendes Image). Anders als bei den
  strukturell analogen `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`-Secrets
  (Welle-Datei §6 Out-of-Scope, dort explizit als externer,
  auftraggeber-gebundener Schritt benannt) ist die GHCR-Paket-Sichtbarkeit
  nirgends als offener, externer Vorbedingungs-Schritt genannt.
- `verifizierbar`: nein — hängt von einer GitHub-Repository-/
  Paket-Einstellung ab, die von hier aus nicht einsehbar ist und erst
  nach dem ersten echten Release real geprüft werden kann.
- `klasse`: fehlender externer Vorbedingungs-Schritt nicht benannt (GHCR-Paket-Sichtbarkeit)

### F-2 — Kopfkommentar von `image-scan.yml` zitiert den Slice-Plan über seinen aktuellen Lifecycle-Pfad statt über eine stabile Kennung

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Zitierform-Konvention, vgl.
  `docs/reviews/review-report.template.md` §Zitier-Form: „Kennung, nicht
  Adresse")
- `pfad`: `.github/workflows/image-scan.yml:22`
- `befund`: Der Kommentar verweist auf
  „`docs/plan/planning/in-progress/release-image-scan.md` §6". Dieser
  Slice-Plan wandert bei Closure per `git mv` nach `done/` (ggf. tiefer
  unter `done/<welle>/`) — der zitierte Pfad existiert danach nicht mehr.
  Kein Sensor deckt diese Stelle: `docs-check`s Link-/Anchor-/
  Tracked-Module parsen Markdown-Linksyntax, keine YAML-Kommentare
  (`.d-check.yml` `scan.roots: ["."]`, aber die Module selbst arbeiten
  gegen `.md`-Dateien). Die Zeile darüber (Zeile 19, Verweis auf die
  Welle-Datei) ist von derselben Fragilität betroffen, folgt aber einem
  bereits an anderer Stelle des Repos (`release.yml:27`) real vorhandenen
  Muster — diese neue Zeile (22) fügt eine zusätzliche, engere
  Lifecycle-Bindung hinzu (Slice- statt Welle-Pfad, mit häufigerem
  Zustandswechsel).
- `verifizierbar`: nein — kein Gate prüft YAML-Kommentare gegen den
  Datei-Bestand; nur Nachschlagen nach der nächsten Slice-Closure zeigt
  die Lücke.
- `klasse`: Zitat auf lifecycle-abhängigen Pfad statt stabiler Kennung

### F-3 — `docker run` ruft Trivy ohne `--image-src remote` auf; die implizite Fallback-Kette erzeugt eine irreführende Vier-Fehler-Meldung

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `Makefile:50`
- `befund`: Trivys `image`-Subcommand probiert standardmäßig
  `docker,containerd,podman,remote` in dieser Reihenfolge (`--image-src`,
  real per `--help` verifiziert). Da der `docker run`-Aufruf keinen
  Docker-/Containerd-/Podman-Socket mountet, scheitern die ersten drei
  Quellen immer strukturell (real reproduziert: „failed to connect to the
  docker API at unix:///var/run/docker.sock", „containerd socket not
  found", „no podman socket found"), bevor die vierte (`remote`) den
  eigentlich relevanten Fehler liefert (`GET
  https://ghcr.io/token?...: DENIED`). Das Ergebnis ist heute korrekt,
  aber implizit: `--image-src remote` würde denselben, für dieses Ziel
  einzig sinnvollen Pfad explizit machen und die drei irrelevanten
  Fehlermeldungen aus jedem CI-Log entfernen, ohne das Verhalten zu
  ändern.
- `verifizierbar`: ja — `make image-cve` real ausgeführt, FATAL-Block mit
  allen vier Teilfehlern bestätigt den Befund unmittelbar.
- `klasse`: implizite Werkzeug-Fallback-Kette statt expliziter Quellenangabe

### F-4 — Tippfehler „saegte" statt „sagte" im Kopfkommentar

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `.github/workflows/image-scan.yml:14`
- `befund`: „ein Scan dagegen saegte nichts ueber das aus" — gemeint ist
  vermutlich „sagte … aus" (aussagen); „saegte" (von „sägen") ergibt
  keinen Sinn im Satz.
- `verifizierbar`: nein — reiner Text, kein Gate prüft Kommentar-Prosa.
- `klasse`: Tippfehler

### F-5 — `ADR-0051`s Fitness-Function-Tabelle nennt `make image-cve` weiterhin als „noch nicht implementiert"

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahl/Aussage im Träger driftet
  gegen die Messung) — Träger liegt **außerhalb** dieses Diffs
- `pfad`: `docs/plan/adr/0051-cicd-pipeline-github-actions.md:291`
- `befund`: Die Zeile „`make image-cve` (noch nicht implementiert — siehe
  Folgepflicht)" ist seit diesem Commit sachlich überholt — das Target
  existiert real. `ADR-0051` ist `Accepted` und laut `AGENTS.md` §3.5
  inhaltlich immutable (keine Zitat-Korrektur-Situation, da nicht
  Zitat-/Verweisgerüst, sondern Sachaussage); dieser Diff hat die Datei
  korrekt nicht angefasst. Bleibt als bekannte, unadressierte
  Text-Staleness stehen — Träger außerhalb des Diffs, deshalb INFO statt
  höherer Kategorie (Skill §Klassifikation: „Träger außerhalb des Diffs
  haben nur einen Leser — die Messung — und bleiben INFO").
- `verifizierbar`: ja, durch Lesen der Zeile gegen `Makefile:46-50`.
- `klasse`: Zahl/Aussage im Träger driftet gegen die Messung (außerhalb Diff)

### F-6 — Risiko-Ausgang „löst sich strukturell … kein weiterer Implementierungsschritt nötig" ist eine Tatsachenbehauptung über einen unbelegten Zukunftszustand

- `kategorie`: INFO (Zuständigkeit Verifier, nicht Reviewer)
- `quelle`: `AGENTS.md` §3.12 Instanz B — laut Hard-Rule-Text explizit dem
  Verifier zugewiesen („der Verifier für Instanz B: ‚Prüfe die Belege,
  nicht die Behauptung'")
- `pfad`: `docs/plan/planning/in-progress/release-image-scan.md:118-120`
- `befund`: Siehe F-1 — dieselbe Beobachtung, hier als Formulierungsfrage
  gerahmt: Die Aussage ist als abgeschlossene Feststellung formuliert,
  nicht als Zusage/Erwartung, obwohl sie von einer heute unbeweisbaren
  externen Bedingung (künftige GHCR-Paket-Sichtbarkeit) abhängt. Ich
  flagge dies als INFO mit Verweis auf die zuständige Rolle, statt es
  selbst als Reviewer-Finding zu werten (Skill §„Was dieser Skill NICHT
  macht": „Keine Verifikation gegen DoD").
- `verifizierbar`: nein.
- `klasse`: Tatsachenbehauptung über unbewiesenen Zukunftszustand

## Negativbefunde

- geprüft, ohne Befund: Digest-Pin `TRIVY_IMAGE` (`Makefile:48`) — real
  gegen `docker buildx imagetools inspect aquasec/trivy:0.74.0` verifiziert:
  Index-Digest `sha256:62b1e65e8869bc4b4c6aa4fa2b21595256c7c2f6018a9d9ad61caf87187c1969`
  ist byte-identisch mit dem im Makefile gepinnten Wert; der Kommentar
  „v0.74.0" ist korrekt.
- geprüft, ohne Befund: reale Ausführung von `make image-cve` — Trivy lädt
  real seine Vulnerability-DB (114,61 MiB, sichtbarer Download-Fortschritt)
  und scheitert danach am `remote`-Pfad mit `GET
  https://ghcr.io/token?scope=repository%3Apt9912%2Fpg-change-feed%3Apull&service=ghcr.io:
  DENIED` — exakt der im Slice-Plan und in der Commit-Message behauptete
  Fehlschlag, kein Kommandozeilen- oder Image-Namens-Fehler.
- geprüft, ohne Befund: YAML-Struktur von `image-scan.yml` — real via
  Ruby-Stdlib-`YAML.load_file` in einem Container geparst: `"on":`
  korrekt gequotet (kein YAML-1.1-Boolean-Problem), `schedule`/
  `workflow_dispatch` beide vorhanden, `permissions: {}` auf Workflow-,
  `permissions: {contents: read}` auf Job-Ebene, ein Job, zwei Steps.
- geprüft, ohne Befund: Action-Pinning — `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1
  # v7.0.1` ist identisch mit dem bereits in `ci.yml`/`e2e.yml`/
  `release.yml`/`examples.yml` verwendeten Pin (`AGENTS.md` §3.8).
- geprüft, ohne Befund: `AGENTS.md` §3.1 (Docker-only) — der einzige
  Ausführungs-Schritt in `image-scan.yml` ist `run: make image-cve`, kein
  Inline-`docker run` außerhalb `make`.
- geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentar-Disziplin) — Zusage-/
  Rang-Zeiger-Klasse an allen neuen Kommentarstellen (`Makefile`-Zeile 47,
  `image-scan.yml`-Kopfkommentar bis auf den in F-4 genannten Tippfehler),
  kein Konjunktiv über verworfene Alternativen, keine Slice-/Wellen-Zahl
  als alleiniger Begründungsanker.
- geprüft, ohne Befund: `harness/README.md` — die entfernte „Nicht
  behauptet (geplant): `make image-cve`"-Zeile hat repo-weit keine
  verbleibende Referenzstelle, die jetzt ins Leere zeigt (per `grep`
  geprüft; verbleibende Treffer liegen ausschließlich in eingefrorenen
  `done/`- bzw. `docs/reviews/`-Lauf-Belegen zu `slice-001`, die per
  Definition nicht mehr aktualisiert werden).
- geprüft, ohne Befund: `make gates` — realer Lauf, Exit-Code direkt
  (ungepiped) geprüft: `0`. Alle Segmente grün — `baseline-verify`
  (v6.9.0, 54 Dateien), `docs-check`/d-check (766 Dateien, 0 Befunde,
  volle Modul-Liste inkl. `hostpaths`), `commit-traceability` (5 Commits
  im Range, alle mit Struktur-ID), `coverage-gate` (82,70 % ggü. Schwelle
  80 %), `generated-sync` (Proto-Erzeugnis byte-gleich), `a-check`
  (0 Befunde). Repo-Status danach sauber (`git status --short` leer).
- geprüft, ohne Befund: Traceability — Commit-Message nennt `ADR-0051`,
  kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `harness/conventions.md` MR-002 (Slice-Kennungen
  sind Namen) — `release-image-scan` ist bereits ein Name, keine Nummer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** fehlender externer
Vorbedingungs-Schritt nicht benannt (GHCR-Paket-Sichtbarkeit) · Zitat auf
lifecycle-abhängigen Pfad statt stabiler Kennung · implizite
Werkzeug-Fallback-Kette statt expliziter Quellenangabe · Tippfehler ·
Zahl/Aussage im Träger driftet gegen die Messung (außerhalb Diff) ·
Tatsachenbehauptung über unbewiesenen Zukunftszustand

## Verdikt

**Merge-blockierend:** ja — die beiden MEDIUM-Findings (F-1, F-2) sind kein
Blocker für die Kern-Mechanik dieses Slice (`make gates` bleibt grün, das
Target läuft real und scheitert real aus dem behaupteten Grund), berühren
aber die Frage, ob `ADR-0051` Entscheidung 6 nach dem ersten echten
Release tatsächlich funktionsfähig sein wird (F-1) bzw. ob ein Verweis
diesen Slice überlebt (F-2). Beides sollte vor Welle-Closure entschieden
sein, nicht erst beim ersten realen Scan-Fehlschlag nach Release auffallen.
Die beiden LOW-Findings (F-3, F-4) sind kleinere Politur ohne
Blockier-Charakter.

**Übergabe:** Da MEDIUM-Findings vorliegen, die eine Rückmeldung an den
Implementer bzw. eine bewusste Architect-/Planner-Entscheidung nahelegen
(F-1 könnte je nach GHCR-Sichtbarkeits-Entscheidung sogar eine
ADR-Ergänzung nach sich ziehen), ziehe ich die DoD-Checkbox „Review
durchgeführt" im Slice-Plan **nicht** nach (Skill-Regel
„DoD-Checkbox-Nachzug ohne Fixrunde" greift nur bei 0 HIGH **und** keiner
Implementer-Rückgabe — hier liegt mit F-1/F-2 eine substanzielle
Rückmeldung vor, auch ohne HIGH). Die Finding-Klassen gehen bei
Slice-/Welle-Closure ins Beobachtungs-Register. F-6 geht als Hinweis an
den Verifier (Instanz-B-Zuständigkeit laut `AGENTS.md` §3.12). Dieser
Report ersetzt keine Verifikation gegen die DoD — das bleibt
Verifier-Aufgabe (Modul 11).
