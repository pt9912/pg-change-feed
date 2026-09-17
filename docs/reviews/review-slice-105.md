# Review-Report: slice-105 — 2026-09-17

**Review-Art:** Code — geprüft gegen Plan + Konventionen (Modul 10
§Drei Review-Arten); dieser Slice trägt keine Code-Änderung im engeren
Sinn, sondern eine Baseline-Hebung + Adaptions-Eintrag — Prüfmaßstab bleibt
Plan/Entscheidungen, nicht DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 3c5e18a..HEAD` (Commits `47d86b2`, `94de7bb`,
`9b66dc1`, `9edad9b`)

**Skill:** `.harness/skills/reviewer.md` @ Accepted (2026-09-09-Fassung,
unverändert seit Eröffnungs-Review)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-105-baseline-v6.9.0-materialisieren.md`
  (vollständig gelesen, §1–§8)
- `harness/conventions.md` (MR-000, MR-001, jetzt MR-002)
- `AGENTS.md` §3.5 (ADR-Immutabilität), §3.13 (Träger nachziehen)
- `ADR-0073` (Zitat-Korrektur an immutablen Dokumenten)
- `ADR-0051` (P8-Zeile, Kontext des Pin-Inventars)
- `v6.9.0` · `regelwerk/grundlagen-source-precedence.md` §Vergabe: woher
  die nächste Kennung kommt (real materialisiert nachgelesen)

---

## Findings

### F-1 — `.d-check.yml`-Ausnahme für `docs/reviews/**` ist breiter als ihr genanntes Vorbild

- `kategorie`: MEDIUM
- `quelle`: Maintainability
- `pfad`: `.d-check.yml:61` (`versions.exempt-paths`)
- `befund`: Die neue Ausnahme wird mit „dieselbe Begründung wie für
  `harness/conventions/done/**`" gerechtfertigt. Das Vorbild ist aber enger:
  eine Adaptions-Datei wird erst durch ihren `git mv` nach `done/` fixiert —
  bis dahin steht sie unter aktiver Versions-Prüfung. `docs/reviews/**` hat
  keine solche Übergabe-Schwelle: Eine Datei ist ab dem ersten Commit, der sie
  anlegt, von der `versions`-Prüfung ausgenommen. Ein frisch geschriebener
  Report, der (per Tippfehler oder Copy-Paste aus einem älteren Report) einen
  inzwischen veralteten `.harness/baseline/vX.Y.Z/…`-Pfad zitiert, obwohl er
  eigentlich den zum Schreibzeitpunkt aktuellen Stand meinte, würde nicht mehr
  rot markiert — anders als bei `conventions/done/**`, wo der Schutz bis zum
  Abschluss aktiv bleibt und erst danach endet. Real reproduziert: ein
  testweise unter `docs/reviews/**` eingefügter `v6.5.0`-Pfad bleibt grün
  (nur `target-missing`/`section-missing` aus anderen Modulen schlagen an,
  `version-stale` nicht), während derselbe Pfad außerhalb der Ausnahme
  (`harness/README.md`) korrekt rot färbt.
- `verifizierbar`: ja — reproduziert (zwei Testeinfügungen, in und außerhalb
  der Ausnahme, `make docs-check` je einmal ausgeführt und wieder
  zurückgenommen).
- `klasse`: Exemption ohne Reifegrenze

### F-2 — Vorbestehender kaputter Anker in `MR-001`, unberührt von diesem Diff

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md:12`
- `befund`: Der Link `grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument`
  zeigt auf eine Überschrift, die in dieser Datei nie existierte — sie liegt
  in `grundlagen-referenz-richtung.md` (`#### Spec-Straten: mehr als ein
  Spec-Dokument`). Gegen `v6.9.0` und, per `git show 3c5e18a:…`, auch schon
  gegen `v6.5.0` bestätigt: der Fehler ist nicht durch diesen Slice
  entstanden, dieser Diff hat nur den Versions-Segment im Pfad nachgezogen
  (korrekt, kein Zitat-Fehler-Anlass laut Diff-Umfang). `docs-check` kann ihn
  strukturell nicht fangen (`.harness/**` ist vom Scan ausgenommen,
  `harness/sensors/docs-check.md` §Grenze Punkt 5). Kein Handlungsbedarf in
  diesem Slice — Kandidat für eine künftige Zitat-Korrektur (`ADR-0073`).
- `verifizierbar`: nein — kein Gate deckt Anker in `.harness/**`-Zielen.
- `klasse`: Anker zeigt auf falsche Datei (vorbestehend)

## Negativbefunde

- geprüft, ohne Befund: SHA256 des `v6.9.0`-Release-Assets — eigener Download
  (`curl` gegen das GitHub-Release), `sha256sum` ergibt
  `8a4e0aaf597a9c67404cb7a350a6fba992f0c98195e011e025c073660ee55cce`,
  identisch zur Implementer-Behauptung.
- geprüft, ohne Befund: `make baseline-verify` — grün, „v6.9.0 OK — 54
  Dateien"; `.harness/baseline/` enthält ausschließlich `v6.9.0/`.
- geprüft, ohne Befund: `sha256sum -c SHA256SUMS` aus
  `.harness/baseline/v6.9.0/` heraus — alle Zeilen `OK`.
- geprüft, ohne Befund: die zwei Zitat-Korrekturen
  (`architect-verdict-aufschub-adresse-verfaellt.md`,
  `review-slice-041-fixrunde-2.md`) — reine Pfad-Segment-Änderung
  `v6.5.0`→`v6.9.0`, Referent (dieselbe Regelwerk-/Templates-Datei,
  Zielanker existiert real in `v6.9.0`) unverändert; kein mitverändertes
  Wort außerhalb des Pfads.
- geprüft, ohne Befund: Präzision der `.d-check.yml`-Ausnahme selbst (nicht
  ihre Breite, siehe F-1) — modul-lokal auf `versions` beschränkt, wirkt
  nicht auf `links`/`anchors`.
- geprüft, ohne Befund: repo-weiter Gegencheck
  `grep -rln "v6\.5\.0" . --exclude-dir=.git` — keine lebenden Träger
  außerhalb `docs/reviews/**` (jetzt ausgenommen), `docs/plan/planning/done/**`
  (Records), `docs/plan/adr/0051-…md` (geprüft, s. u.) und der Plan-Datei
  von `slice-105` selbst (mit `d-check:ignore`-Markern annotiert).
- geprüft, ohne Befund: `ADR-0051` P8-Zeile (`Kurs-Baseline v6.5.0`) — Status
  `Accepted`, reine Pin-Inventur-Momentaufnahme zum Entscheidungszeitpunkt,
  kein Link/Zitat auf die vendorte Baseline; `AGENTS.md` §3.5 greift korrekt
  (kein Zitat-Fehler-Anlass).
- geprüft, ohne Befund: `MR-002` — Wortlaut-Zitate „Welle- und
  Slice-Kennungen sind Namen, nicht Nummern — unabhängig von der
  Schreiberzahl" und „Schreiber ist, was committet: ein Mensch, ein Agent,
  ein Automat" stimmen zeichengenau mit der materialisierten
  `v6.9.0`-Datei überein; der Anker `#vergabe-woher-die-nächste-kennung-kommt`
  löst auf; Geltungsbereich „`slice-001`–`slice-105`" ist im Dateikopf, in
  §Geltungsbereich und in §Adaption konsistent (105 ist die letzte
  nummerierte Kennung, korrekt eingeschlossen).
- geprüft, ohne Befund: `.claude/agents/architect.md`,
  `.claude/agents/reviewer.md`, `.harness/skills/closure-note-reviewer.md`,
  `AGENTS.md`, `harness/sensors/baseline-verify.md`,
  `harness/conventions/MR-001-…md` — jeweils genau eine Zeile,
  Versions-Segment im Pfad geändert, kein weiterer Inhalt berührt.
- geprüft, ohne Befund: Out-of-Scope-Disziplin — `git diff --name-only`
  gegen die Lifecycle-Verzeichnisse zeigt ausschließlich `slice-105` selbst;
  keine der 104 bestehenden Slice-Dateien umbenannt oder inhaltlich
  angefasst; keine Umsetzung der Wellen-Funde 129/132/134/135/136/137
  (nur zur Kenntnis in §7 vermerkt).
- geprüft, ohne Befund: Commit-Traceability — `RANGE=3c5e18a..HEAD make
  commit-traceability` grün, alle vier Commits tragen `ADR-0045` im Betreff,
  keine `SPEC-*`/`ARC-*`-Kennung im Betreff.
- geprüft, ohne Befund: `make gates` — vollständiger Lauf grün
  (`baseline-verify`, `docs-check`, `commit-traceability`, `coverage-gate`,
  `generated-sync`, `a-check`), ungepiped geprüft.
- geprüft, ohne Befund: eigener §3.13-Gegencheck — kein weiterer Träger
  gefunden, der eine „54 Dateien"/Versions-Zahl oder eine Slice-Zählung
  (104/105) hält, die durch diesen Diff falsch geworden wäre;
  `harness/sensors/baseline-verify.md`s „54 Dateien"-Beispielzahl bleibt
  zufällig identisch zur neuen Zählung (real gegengeprüft), keine Drift.
- geprüft, ohne Befund: §8-Sichtungsschritt (offene Beobachtungen) — eigener
  Gegen-Grep über `docs/plan/planning/observations/`; die fünf Treffer auf
  „baseline" sind ausschließlich generische „Baseline-Regelwerk
  `modul-XX-…`"-Zitate, kein Treffer zum Versions-/Namenskonventions-Thema.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Exemption ohne Reifegrenze · Anker zeigt
auf falsche Datei (vorbestehend)

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. F-1 ist eine Abwägungsfrage
(dokumentiert, mit eigenem Reproduktions-Beleg), keine Verletzung einer
Hard Rule oder eines ADRs; eine Fixrunde am Implementer ist dafür nicht
zwingend — Kandidat für einen künftigen, engeren Zuschnitt der Ausnahme
(z. B. Beschränkung auf bereits bestehende Dateien zum Zeitpunkt der
Einführung) oder eine benannte Akzeptanz des Trade-offs, zu entscheiden bei
der Slice-Closure. F-2 ist eine vorbestehende, außerhalb des Diff-Umfangs
liegende Fundstelle ohne Handlungsbedarf in diesem Slice.

**Übergabe:** Da keine Fixrunde am Implementer erforderlich ist, zieht
dieser Report selbst die DoD-Checkbox „Review durchgeführt" im Slice-Plan
nach (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde),
im selben Commit wie dieser Report. Die Finding-Klassen gehen in die
Slice-Closure §7 und von dort in den Zähler des Beobachtungs-Registers.
Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen. Verifikation (DoD-/Spec-Konformität) bleibt Verifier-Aufgabe
(Modul 11, frischer Kontext).
