# Welle release-pipeline-adr-0051: Release-Pipeline gemäß ADR-0051 vollständig umsetzen

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, direkt beauftragt). **Datum:** 2026-09-19.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) (Accepted,
2026-09-13) legt fünf Teilentscheidungen für eine vollständige
Release-Pipeline fest; zwei ihrer Träger existieren bereits
(`ci.yml`, `.github/dependabot.yml`, aus dem archivierten
`slice-039`), die übrigen fünf Artefakte (`docs/user/version.md`,
`release.yml`, `image-scan.yml`, `upstream-drift.yml`,
`hub-description.yml`) sowie die Erweiterung von `make image` um
einen Versions-/Tag-Parameter und die neuen Pin-Freshness-Targets
für P3–P9 fehlen. Das *Mehr* gegenüber fünf isolierten Slice-DoDs:
`ADR-0051` ist als **eine** Entscheidung erst eingelöst, wenn alle
Teile zusammen bestehen und ineinandergreifen — `hub-description.yml`
läuft als `needs: release`-Job aus `release.yml`, und
`docs/user/releasing.md` (Ziel dieser Welle) beschreibt den
**vollständigen** realen Prozess, nicht einen Teilausschnitt. Ein
einzelner fertiger Slice ohne die anderen wäre nur ein Fragment der
Entscheidung, keine abgeschlossene Umsetzung.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) ist `Accepted`
  (bereits erfüllt, 2026-09-13).
- Direkter Auftraggeber-Auftrag, die verbleibenden fünf Artefakte der
  Entscheidung jetzt real umzusetzen ("Release-Pipeline jetzt wirklich
  bauen" / "alles — leg eine Welle mit slices dafür an").

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle fünf Slices in `done/`.
- `make gates` grün.
- `docs/user/releasing.md` beschreibt den nach Abschluss der vier
  vorangehenden Slices real existierenden Prozess (kein Plan-Text über
  einen noch nicht gebauten Mechanismus) — das ist das *Mehr*: erst wenn
  `release.yml`, `image-scan.yml`, `upstream-drift.yml` und
  `hub-description.yml` real bestehen und ineinandergreifen, kann dieser
  letzte Slice sie wahrheitsgemäß beschreiben.
- Closure-Notiz in `welle-release-pipeline-adr-0051-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| release-version-und-workflow | `docs/user/version.md`, `make image`-Versionsparameter, `release.yml` (Tag → Build → GHCR+Docker-Hub-Push → GitHub-Release) | [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 1/3/4 |
| release-image-scan | `image-scan.yml` — Trivy-CVE-Scan gegen das publizierte GHCR-`:latest`-Image, advisory | [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 6 |
| release-upstream-drift | `upstream-drift.yml` — Pin-Freshness über das Neun-Achsen-Inventar P1–P9, fail-open, neue Make-Targets für P3–P9 | [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 |
| release-hub-description | `hub-description.yml` — Docker-Hub-Beschreibungs-Sync nach erfolgreichem Release | [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 8 |
| release-doku-releasing | `docs/user/releasing.md` — Betreiber-/Maintainer-Dokumentation des jetzt real existierenden Release-Prozesses | [`ADR-0051`](../adr/0051-cicd-pipeline-github-actions.md) (derivativ, kein eigener Entscheidungspunkt) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- `release-hub-description` wird blockiert von `release-version-und-workflow`
  (läuft als `needs: release`-Job **aus** `release.yml` — ohne dessen
  Job-Namen kein sinnvoller `needs`-Bezug).
- `release-doku-releasing` wird blockiert von allen vier vorangehenden
  Slices (dokumentiert den vollständigen realen Endzustand, kein
  Zwischenstand).
- `release-image-scan` und `release-upstream-drift` sind voneinander und
  von `release-version-und-workflow` inhaltlich unabhängig — die
  Reihenfolge unter ihnen ist Planungsentscheidung (WIP-Limit 1), keine
  fachliche Abhängigkeit.
- Blockiert keine andere, bereits laufende Welle — *Offene Wellen* führte
  vor dieser Eröffnung keinen Eintrag.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Ein tatsächlicher Release** (`git tag v<SemVer>` + `git push --tags`) —
  eine irreversible, extern sichtbare Aktion (öffentlicher Registry-Push,
  öffentliches GitHub-Release), die nur nach expliziter, gesonderter
  Rückfrage beim Auftraggeber läuft, nie als Teil einer Slice-DoD dieser
  Welle. Diese Welle liefert den **Mechanismus**, nicht seine erste
  Anwendung.
- **Anlage der `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`-Repository-Secrets** —
  eine externe, kontobezogene Handlung im GitHub-Repo-Settings, die nur der
  Auftraggeber selbst ausführen kann; kein technischer Slice-Gegenstand.
  `release.yml` referenziert die Secret-Namen (`ADR-0051` Entscheidung 2),
  ihr tatsächliches Vorhandensein bleibt bis zum ersten echten Release
  unbewiesen (`AGENTS.md` §3.10 analog).
- **Ein echter Alt-Image-vs-Neu-Image-Upgrade-Test**
  (`BEO-PGC/kein-echter-versionswechsel-upgrade-test`) — dessen
  Re-Evaluierungs-Trigger
  ([`ADR-0064`](../adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
  Re-Evaluierungs-Trigger 1: „Release-Pipeline/Tags abgeschlossen, erster
  Git-Tag gesetzt") wird durch diese Welle nur zur Hälfte erfüllt — die
  Pipeline existiert danach real, ein erster echter Tag ist aber (siehe
  erster Punkt oben) bewusst nicht Teil dieser Welle. Bleibt ein
  Folge-Vorgang nach dem ersten echten Release.
- **Erweiterung von `make image` um den Versions-/Tag-Parameter** und die
  **neuen P3–P9-Pin-Freshness-Make-Targets** laufen als Teil ihrer
  jeweiligen Slices (`release-version-und-workflow` bzw.
  `release-upstream-drift`), nicht als eigene Slices dieser Welle.

## 7. Closure-Notiz

*(wird bei Closure gefüllt.)*

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: <Zeiger auf `welle-<Kennung>-results.md`, Geschwister im Ruheort `done/`>
Zähler: <Zeiger aufs Beobachtungs-Register, eine Ebene über dem Ruheort>
