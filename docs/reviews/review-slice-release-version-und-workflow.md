# Review-Report: release-version-und-workflow — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan (`docs/plan/planning/in-progress/release-version-und-workflow.md`),
`ADR-0051` Entscheidung 1/2/3 und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten).

**Gegenstand:** Commit `5e3c7b68` (Parent `75d17987`), Slice
`release-version-und-workflow`, Welle `welle-release-pipeline-adr-0051`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-version-und-workflow.md` (Slice-Plan)
- `docs/plan/planning/welle-release-pipeline-adr-0051.md` (Welle-Datei, §6 Out-of-Scope)
- `ADR-0051` (Accepted) — Entscheidung 1/2/3, Verglichene Alternativen 1–3
- `ADR-0103` (Accepted, Supersedes `ADR-0044` teilweise) — Digest-Beleg-Semantik
- `ADR-0044` (Accepted) — Image-Beleg-Semantik
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.8
  (Action-Pinning), §3.10 (Post-Push-Verifikationspflicht), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema)

---

## Findings

### F-1 — DoD-Checkboxen zitieren §7 als Beleg-Anker; §7 ist unausgefüllt

- `kategorie`: HIGH
- `quelle`: Hard-Rule „Beleg trägt seinen Satz nicht" (`.harness/skills/reviewer.md`
  §Klassifikation HIGH, `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`)
- `pfad`: `docs/plan/planning/in-progress/release-version-und-workflow.md:76-77`
  und `:84-86` (DoD-Zeilen, „siehe §7") gegen `:158-160` (§7 Closure-Notiz)
- `befund`: Zwei DoD-Checkbox-Zeilen behaupten reale Verifikation
  (`docker buildx imagetools inspect` gegen zwei lokale Registry-Container;
  Ruby-Stdlib-YAML-Parse + `bash -n` je `run:`-Schritt) und verweisen dafür
  explizit auf „siehe §7". §7 „Closure-Notiz" trägt im selben Commit,
  der diese Zeilen auf `[x]` setzt, weiterhin ausschließlich den
  unausgefüllten Platzhaltertext `*(wird bei Bearbeitung gefüllt.)*`. Es
  existiert im Diff auch sonst kein committetes Artefakt (Skript, Log,
  Protokoll), das die behauptete Verifikation nachvollziehbar macht — die
  einzige Spur ist die Commit-Message-Prosa. Ich habe die behaupteten
  Prüfungen unabhängig selbst nachgefahren (SemVer-Regex-Testfälle,
  `bash -n` auf allen fünf extrahierten `run:`-Blöcken, YAML-Parse via
  Ruby) und sie bestätigt sich als plausibel durchführbar — das ändert
  aber nichts daran, dass der im Slice-Plan genannte Beleg-Anker (§7) die
  Aussage nicht trägt, die er zu tragen behauptet.
- `verifizierbar`: nein — kein Gate prüft die Vollständigkeit einer
  Closure-Notiz; nur Lesen der Datei zeigt die Lücke.
- `klasse`: Beleg trägt seinen Satz nicht (Zitat auf leere Stelle)

### F-2 — SemVer-Regex akzeptiert numerische Prerelease-Identifier mit führender Null

- `kategorie`: MEDIUM
- `quelle`: `ADR-0051` Entscheidung 1 („validiert den Tag strikt gegen die
  SemVer-2.0-Form")
- `pfad`: `.github/workflows/release.yml:98` (`semver_re=...`)
- `befund`: Die Regex `^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$`
  behandelt jeden Prerelease-Identifier als generisches
  `[0-9A-Za-z-]+`-Muster, ohne zwischen rein-numerischen und
  alphanumerischen Identifiern zu unterscheiden. SemVer 2.0 §9 verbietet
  führende Nullen bei rein-numerischen Prerelease-Identifiern. Eigener
  Test (`bash -c '[[ "1.0.0-01" =~ $semver_re ]]'`) bestätigt: `1.0.0-01`
  und `1.0.0-1.02` werden fälschlich als gültig akzeptiert. Die Haupt- und
  Nebenversion-Komponenten selbst sind korrekt gegen führende Nullen
  geschützt (`01.0.0` wird korrekt abgelehnt, real getestet).
- `verifizierbar`: nein — kein Gate prüft diese Regex-Semantik; per
  Bash-Testfall reproduzierbar (siehe Befund).
- `klasse`: SemVer-Validierung unvollständig gegenüber Spezifikation

### F-3 — Stabilitätserkennung markiert Bindestrich im Build-Metadata-Teil fälschlich als Prerelease

- `kategorie`: MEDIUM
- `quelle`: `ADR-0051` Entscheidung 2 Option D (":latest nur bei stabilen,
  nicht-Prerelease Tags")
- `pfad`: `.github/workflows/release.yml:113-118` (`case "$version" in *-*) ...`)
- `befund`: Das Muster `*-*` prüft nur, ob der Versions-String irgendwo
  einen Bindestrich enthält — unabhängig davon, ob dieser vor oder nach
  einem `+` (Build-Metadata-Trenner) steht. Ein stabiler Tag mit
  Build-Metadata, die selbst einen Bindestrich enthält — SemVer 2.0s
  eigenes Spezifikationsbeispiel `1.0.0+exp-sha.5114f85` sowie
  `1.0.0+21AF26D3---117B344092BD` —, ist laut Spezifikation **kein**
  Prerelease, wird von diesem `case`-Ausdruck aber als
  `latest=false` behandelt. Eigener Test bestätigt: `1.0.0+build.5` (kein
  Bindestrich) ergibt korrekt `latest=true`, `1.0.0+exp-sha.5114f85`
  (Bindestrich nur in der Build-Metadata) ergibt fälschlich `latest=false`.
- `verifizierbar`: nein — kein Gate prüft diese Fallunterscheidung; per
  Bash-Testfall reproduzierbar (siehe Befund).
- `klasse`: SemVer-Stabilitätserkennung verwechselt Bindestrich-Vorkommen mit Prerelease

### F-4 — Keine Negativtests für die SemVer-Validierung/Stabilitätserkennung

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation MEDIUM
  („fehlende Negativtests bei neuem öffentlichem Vertrag")
- `pfad`: `.github/workflows/release.yml:94-118`
- `befund`: Die Tag-Validierungsregel und die Stabilitätserkennung sind
  der öffentliche Vertrag, der entscheidet, welche Tags einen Release
  auslösen und ob `:latest` gesetzt wird. Es existiert kein
  automatisierter, committeter Test dieser Logik (kein Skript unter
  `tools/`, kein Testfall) — die einzige Prüfung ist eine laut
  Commit-Message durchgeführte, nicht committete Ad-hoc-Probe. F-2 und
  F-3 wären mit einer minimalen, committeten Testtabelle (gültige/ungültige
  Beispiele wie in dieser Review-Anfrage vorgegeben) mechanisch aufgefallen.
- `verifizierbar`: nein — kein Gate erzwingt Tests für
  Workflow-Inline-Shell-Logik.
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag

### F-5 — ADR-0103s eigener Re-Evaluierungs-Trigger ist sachlich eingetreten, ohne Folge-ADR

- `kategorie`: MEDIUM
- `quelle`: `ADR-0103` §Re-Evaluierungs-Trigger
- `pfad`: `Makefile:168-178` (`ifdef VERSION` → `docker buildx build --push`)
  gegen `docs/plan/adr/0103-image-hash-lokal-statt-committet.md:153-161`
- `befund`: `ADR-0103` §Kontext (2) begründet die lokale, nicht-committete
  Form von `harness/image-hash.txt` explizit damit, dass „kein `docker
  push` existiert irgendwo im Repo" (real per `grep` geprüft, Stand vor
  diesem Slice). Ihr eigener §Re-Evaluierungs-Trigger nennt als
  Beispielinstanz „ein Lauf-Zweig braucht einen historischen, über
  Commits hinweg nachschlagbaren Image-Beleg — etwa ein
  `docker push`-Workflow (Archiv-Form wird real erfüllt)" und schreibt
  vor: „Dann wird der Beleg als Folge-ADR mit `supersedes` wieder
  gehoben." Dieser Slice fügt genau das hinzu: einen realen
  `docker buildx build --push`-Pfad, erreichbar über `release.yml`. Die
  neue `harness/README.md`-Formulierung selbst widerspricht sich nicht —
  sie grenzt sauber ab, dass der `:dev`-Pfad unverändert bleibt und der
  `VERSION`-Pfad seinen eigenen, nicht-`harness/image-hash.txt`-basierten
  Beleg (GitHub-Release-Beschreibungstext) bekommt. Ob das bereits die vom
  Trigger geforderte Folge-ADR ersetzt oder ob der Trigger formal als
  ausgelöst gilt und eine explizite Folge-ADR (oder ein expliziter
  „Trigger greift hier nicht, weil …"-Vermerk) fehlt, ist eine
  Architect-Entscheidung, keine Reviewer-Entscheidung — ich flagge den
  Befund, ohne ihn selbst aufzulösen.
- `verifizierbar`: nein — ADR-Trigger-Bewertung ist keine Gate-Prüfung.
- `klasse`: ADR-Re-Evaluierungs-Trigger eingetreten ohne Folge-ADR

### F-6 — DoD-Checkbox „Doku-Update harness/README.md" nicht nachgezogen trotz erledigter Arbeit im selben Commit

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/release-version-und-workflow.md:91-92`
- `befund`: Der Commit `5e3c7b68` aktualisiert `harness/README.md`
  tatsächlich (`make image`-Zeile, neue `release.yml`-Zeile) — die
  entsprechende DoD-Checkbox „Doku-Update für `harness/README.md`" bleibt
  im selben Commit auf `[ ]` stehen, obwohl die Arbeit bereits erledigt
  ist.
- `verifizierbar`: nein — Checkbox-Zustand ist keine Gate-Prüfung.
- `klasse`: DoD-Checkbox hinter dem tatsächlichen Fortschritt

### F-7 — `benutzerhandbuch.md`s „kein veröffentlichtes Image"-Satz wird nach dem ersten realen Release stale

- `kategorie`: INFO
- `quelle`: Maintainability / `AGENTS.md` §3.13 (vorausschauend)
- `pfad`: `docs/user/benutzerhandbuch.md:40-41`
- `befund`: „Es gibt aktuell kein veröffentlichtes Container-Image" ist
  heute noch wahr (kein realer Tag-Push fand statt, bewusst außerhalb
  dieses Slices). Sobald der erste reale Release-Tag gepusht wird, wird
  dieser Satz falsch, ohne dass ein Mechanismus ihn automatisch
  nachzieht. Die Welle plant einen eigenen Folge-Slice
  (`release-doku-releasing` für `docs/user/releasing.md`) — dieser Satz
  im Benutzerhandbuch ist eine andere Datei und nicht explizit als dessen
  Ziel genannt.
- `verifizierbar`: nein.
- `klasse`: vorausschauende Träger-Staleness

### F-8 — Zwei-lokale-Registries-Test kann die realen Tag-Strings nicht wörtlich abdecken; `docker.io/`-Präfix-Form selbst separat verifiziert

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `Makefile:171-172` (`-t ghcr.io/pt9912/pg-change-feed:$(VERSION)`,
  `-t docker.io/pt9912/pg-change-feed:$(VERSION)`)
- `befund`: Ein generischer lokaler `registry:2`-Container wird zwangsläufig
  über eine `host:port`-Adresse angesprochen, nicht über den Hostnamen
  `docker.io` (der bei Docker-Tooling fest auf Docker Hub verdrahtet ist)
  — die reale lokale Probe konnte also nicht wortgleich die im Makefile
  stehenden `docker.io/pt9912/pg-change-feed:...`-Tag-Strings verwenden,
  sondern nur die Multi-Tag-Push-in-einem-Build-Mechanik strukturell
  analog. Ich habe die Namensform selbst unabhängig gegen echtes
  Docker-Tooling geprüft (`docker buildx build --load -t
  docker.io/testuser/testrepo:1.2.3 .`): Docker normalisiert
  `docker.io/<user>/<repo>:<tag>` korrekt zu `<user>/<repo>:<tag>` — die
  im Makefile verwendete Form ist syntaktisch korrekt für einen
  Docker-Hub-Push, kein Fehler.
- `verifizierbar`: nein.
- `klasse`: Ersatz-Evidenz deckt nicht wortgleich den Produktionspfad

## Negativbefunde

- geprüft, ohne Befund: Action-Pins in `release.yml`
  (`actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1` # v7.0.1,
  `docker/setup-buildx-action@f87e5991a6d7451dcb8d9637bfbc97413f497069` # v4.4.1,
  `docker/login-action@dbcb813823bdd20940b903addbd779551569679f` # v4.6.0) —
  alle drei SHAs real gegen `git ls-remote --tags` der jeweiligen
  Upstream-Repos verifiziert, exakt passend zum genannten Tag; der
  `checkout`-Pin ist identisch mit dem bereits in `ci.yml`/`e2e.yml`
  verwendeten.
- geprüft, ohne Befund: `Makefile` `ifdef VERSION`-Zweig — `make -n image`,
  `make -n image VERSION=1.2.3`, `make -n image VERSION=1.2.3 LATEST=true`
  rendern jeweils exakt das erwartete Kommando (unveränderter `--load`-Pfad
  ohne `VERSION`; `--push` mit zwei bzw. vier `-t`-Flags mit/ohne `LATEST`).
- geprüft, ohne Befund: Behauptung im `release.yml`-Kopfkommentar,
  `ci.yml`/`e2e.yml` schlössen Tag-Pushes per `tags-ignore: ['**']` aus —
  real in beiden Dateien so vorgefunden.
- geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentar-Disziplin) an allen
  neuen Kommentarstellen (`Makefile`-Blockkommentar vor `ifdef VERSION`,
  beide `##`-Zielkommentare, `release.yml`-Kopfkommentar) — Zusage/
  Kopplung/Rang-Zeiger, kein Konjunktiv über verworfene Alternativen,
  keine Slice-/Wellen-Zahl als Begründungsanker im Sinn des verbotenen
  Musters.
- geprüft, ohne Befund: YAML-Struktur von `release.yml` (unabhängig via
  Ruby-Stdlib-`YAML.load_file` reproduziert) und `bash -n` auf allen fünf
  extrahierten `run:`-Blöcken — alle syntaktisch fehlerfrei.
- geprüft, ohne Befund: `make gates` — realer Lauf, Exit-Code direkt
  (ungepiped) geprüft: `0`. Alle Segmente grün — `baseline-verify` (v6.9.0,
  54 Dateien), `docs-check`/d-check (763 Dateien, 0 Befunde, volle
  Modul-Liste inkl. `hostpaths`), `commit-traceability` (5 Commits im
  Range, alle mit Struktur-ID), `coverage-gate` (82.80 % ggü. Schwelle
  80 %), `generated-sync` (Proto-Erzeugnis byte-gleich), `a-check`
  (0 Befunde).
- geprüft, ohne Befund: Traceability — Commit-Message nennt `ADR-0051`.
- geprüft, ohne Befund: `docs/user/version.md` — genau eine Zeile,
  `0.1.0`, gültiges SemVer 2.0 ohne führendes `v`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 4 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht (Zitat auf
leere Stelle) · SemVer-Validierung unvollständig gegenüber Spezifikation
(2×, F-2/F-3) · fehlende Negativtests bei neuem öffentlichem Vertrag ·
ADR-Re-Evaluierungs-Trigger eingetreten ohne Folge-ADR · DoD-Checkbox
hinter dem tatsächlichen Fortschritt · vorausschauende Träger-Staleness ·
Ersatz-Evidenz deckt nicht wortgleich den Produktionspfad

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) ist ein Rückgabe-Pfeil an den
Implementer: entweder §7 mit der tatsächlich durchgeführten Verifikation
füllen (Befehle, Ergebnisse) oder die „siehe §7"-Verweise aus den
DoD-Zeilen entfernen, bis der Beleg wirklich dort steht. Die vier
MEDIUM-Findings (F-2 bis F-5) sind kein Blocker für die Kern-Mechanik
(`make gates` bleibt grün, die Digest-Content-Mirror-Eigenschaft ist real
geprüft), gehören aber vor Closure entschieden — F-2/F-3 sind mechanisch
demonstrierte Abweichungen von der als „strikt" behaupteten SemVer-2.0-
Validierung, F-5 ist eine Architect-Frage (`ADR-0103`-Trigger), keine, die
der Implementer allein auflösen sollte.

**Übergabe:** Da mindestens ein HIGH-Finding eine Implementer-Fixrunde
auslöst, ziehe ich die DoD-Checkbox „Review durchgeführt" im Slice-Plan
**nicht** nach (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift
nur bei 0 HIGH bzw. wenn keine Rückgabe an den Implementer erfolgt — hier
nicht der Fall). Die Finding-Klassen gehen bei Slice-Closure ins
Beobachtungs-Register. Dieser Report ersetzt keine Verifikation gegen die
DoD — das bleibt Verifier-Aufgabe (Modul 11).
