# Review-Report: release-binary-version-injektion — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/release-binary-version-injektion.md`),
`ADR-0051` und `AGENTS.md` Hard Rules, insbesondere §3.7
(Kommentar-Disziplin) (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `fe805a31` ("feat(release): --version zeigt die
injizierte Release-Version ([`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md))"), Slice
`release-binary-version-injektion`, ohne Welle.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen, seither um
mehrere weitere HIGH-Klassen ergänzt).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-binary-version-injektion.md` (Slice-Plan, §1–§8)
- `ADR-0051` (Accepted) — CI/CD-Pipeline über GitHub Actions
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.10
  (Post-Push-Verifikationspflicht), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema)
- `docs/reviews/review-slice-release-version-und-workflow.md` (Nachbar-Review
  derselben Sub-Area, als Stil- und Präzedenzvergleich)

---

## Findings

### F-1 — Versions-Injektionsmechanismus ohne automatisierten Regressionstest

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation MEDIUM
  („fehlende Negativtests bei neuem öffentlichem Vertrag")
- `pfad`: `Dockerfile:118-124` (`ARG VERSION=dev`, `-ldflags -X
  main.version=$VERSION`), `Makefile:46-57` (`ifdef VERSION`-Zweig)
- `befund`: Der gesamte Injektionspfad (`release.yml` → `make image
  VERSION=…` → `--build-arg VERSION` → `Dockerfile`-`ARG`/`-ldflags -X` →
  `main.version`) trägt keinen committeten, automatisierten Test. Die
  einzige Absicherung ist die im Slice-Plan/der Commit-Message
  dokumentierte, nicht committete Ad-hoc-Probe (`docker buildx build
  --load --build-arg VERSION=…`). Genau dieser Mechanismus ist der
  Gegenstand des Bugs, den dieser Slice behebt (`--version` zeigte seit
  Ersteinführung unbemerkt einen falschen Wert) — ohne Regressionstest
  bleibt eine künftige, versehentliche Entfernung des `-X`-Flags oder der
  `ARG`-Zeile aus demselben Grund unbemerkt: Ich habe empirisch bestätigt
  (siehe Negativbefunde), dass eine solche Regression **lautlos** wäre —
  der Build bricht nicht ab, `main.version` bleibt nur unverändert. Kein
  bestehendes Gate (`make gates`, `make test`, `make test-integration`)
  prüft den gebauten Binary-Wert gegen den übergebenen `VERSION`-Parameter.
- `verifizierbar`: nein — kein Gate erzwingt einen Test für diesen
  Docker-Build-Mechanismus; die Lücke ist über den bestehenden
  `make test-integration`-Rundlauf (baut bereits `make image` und startet
  den Feed-Container) real schließbar, aber aktuell ungenutzt.
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag

### F-2 — §6-Risiko 1 klassifiziert eine tatsächlich eingetretene, bewusst akzeptierte Verhaltensänderung als „entfallen"

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/release-binary-version-injektion.md:117-123`
  (Abschnitt bereits vor diesem Commit von Planner-Hand gefüllt, in diesem
  Diff unverändert)
- `befund`: Das Risiko beschreibt eine reale, bereits verifiziert
  eingetretene Verhaltensänderung (`:dev`-Build zeigt jetzt `dev` statt
  `0.2.0-verdrahtung`) und trägt trotzdem den Ausgang „entfallen als
  Risiko". „Entfallen" liest sich eher als „ist nicht eingetreten"/„wurde
  gegenstandslos" — treffender wäre „eingetreten, akzeptiert" gewesen, da
  die beschriebene Änderung real stattfand und lediglich als kein
  Informationsverlust bewertet wird. Kein Merge-Blocker, da §6-Ausgänge
  laut DoD ohnehin erst zur Closure endgültig geprüft werden und dieser
  Text nicht Teil des vorliegenden Commits ist.
- `verifizierbar`: nein.
- `klasse`: Ausgangsklasse eines Risikos leicht unscharf

## Negativbefunde

- geprüft, ohne Befund: `const` → `var` (`cmd/pg-change-feed/main.go:25`)
  — Go-Linker-Verhalten eigenständig empirisch reproduziert (nicht nur
  vertraut): ein Minimal-Programm mit `const version = "0.2.0-verdrahtung"`
  und `go build -ldflags="-X main.version=injected"` **kompiliert
  fehlerfrei (Exit 0)**, das erzeugte Binary zeigt aber weiterhin
  `0.2.0-verdrahtung` — der Linker-Flag hat bei einer Konstante **lautlos
  keine Wirkung** (kein Compile-Fehler, kein Warnhinweis sichtbar in
  Standard-Ausgabe). Nach `const` → `var` liefert exakt derselbe Aufruf
  korrekt `injected`. Die DoD-Begründung „Linker-Flag `-X` kann nur
  Variablen, keine Konstanten, überschreiben" ist damit nicht nur
  plausibel, sondern empirisch bestätigt — und die Änderung war
  zwingend notwendig, keine kosmetische Wahl.
- geprüft, ohne Befund: die zwei im Commit/DoD behaupteten Build-Ergebnisse
  eigenständig nachgefahren — `docker buildx build --load --build-arg
  VERSION=0.1.1-test -t <tag> .` gefolgt von `docker run --rm <tag>
  --version` liefert real `pg-change-feed 0.1.1-test`; derselbe Build ohne
  `--build-arg` (Dockerfile-Default `dev`, entspricht dem `VERSION`-losen
  `make image`-Zweig) liefert real `pg-change-feed dev`. Beide Ergebnisse
  stimmen mit den im Slice-Plan/der Commit-Message genannten Werten
  überein.
- geprüft, ohne Befund: Seiteneffekt der `var`-Änderung im restlichen Code
  — `grep -rn "\bversion\b" cmd/pg-change-feed/` findet ausschließlich die
  Deklaration selbst, den `fmt.Printf`-Aufruf im `--version`-Zweig und den
  Vergleichswert in `main_test.go:132` (`TestVersionMeldetDenLieferstandOhneDiagnose`,
  liest die package-level-Variable zur Testlaufzeit, unbeeinflusst vom
  Linker-Flag, da `go test` nicht über `-ldflags -X` läuft) — kein
  `const`-Block, kein Array-Größen-Kontext, keine Stelle, die eine
  Konstante zwingend voraussetzt.
- geprüft, ohne Befund: `AGENTS.md` §3.7 (Kommentar-Disziplin) an beiden
  neuen/geänderten Kommentarstellen (`cmd/pg-change-feed/main.go:18-24`,
  `Dockerfile:118-120`) — beide indikativ, keine verworfene Alternative im
  Konjunktiv, kein abwesender Text, keine Slice-/Wellen-Chronik im
  Produktionscode-Kommentar; beide tragen Kopplungs- und
  Rang-Zeiger-Information ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)/[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md)/[`ADR-0103`](../plan/adr/0103-image-hash-lokal-statt-committet.md)) statt Vorher/Nachher-
  Erzählung.
- geprüft, ohne Befund: §1 „Ausdrücklich NICHT in diesem Slice" — die drei
  Ausschlüsse (kein `git describe`-Fallback, keine rückwirkende Korrektur
  veröffentlichter Images, keine Versionsangabe in `diagnose`/
  `--healthcheck`) sind plausibel begründet; `grep -rn "version"
  internal/bootstrap/diagnose*.go internal/bootstrap/healthcheck*.go`
  bestätigt: beide Pfade referenzieren aktuell keine Versionsinformation,
  der Ausschluss erweitert also tatsächlich keinen bestehenden Vertrag.
- geprüft, ohne Befund: §3 Plan-Tabelle gegen den realen Diff — exakt die
  drei geplanten Dateien (`main.go`, `Dockerfile`, `Makefile`) geändert,
  kein zusätzlicher, ungeplanter Umfang; DoD zählt ≤ 3 Liefer-Punkte ein.
- geprüft, ohne Befund: `AGENTS.md` §3.13 (Träger-Nachzug) —
  `docs/user/benutzerhandbuch.md` dokumentiert `--version` nirgends
  inhaltlich (kein Treffer für `--version` im Handbuch); der bereits
  bekannte Altfall (Kopf-Feld „Software-Version" zeigte ebenfalls
  `0.2.0-verdrahtung`) ist laut Änderungshistorie-Zeile 1.32 bereits in
  einem separaten, vorherigen Schritt auf einen Verweis auf
  `docs/user/version.md` umgestellt — keine Überschneidung, kein
  liegengebliebener Träger.
- geprüft, ohne Befund: end-to-end-Verdrahtung — `release.yml:97` ruft
  real `make image VERSION="$version" LATEST="$latest"` auf; der
  `VERSION`-Zweig des `Makefile` übergibt `--build-arg VERSION=$(VERSION)`
  an genau den in diesem Commit geänderten Dockerfile-Pfad. Der
  `VERSION`-lose `:dev`-Zweig (`Makefile:56`) übergibt weiterhin keinen
  `--build-arg` und nutzt damit den Dockerfile-Default `dev` — unverändert
  gegenüber vor diesem Commit.
- geprüft, ohne Befund: andere Docker-Build-Aufrufer derselben `build`-Stufe
  (`tools/harness/generated-sync.sh`, `tools/harness/proto-generate.sh`) —
  keiner übergibt `VERSION` explizit; da `ARG VERSION=dev` einen Default
  trägt, bricht keiner der beiden Aufrufe.
- geprüft, ohne Befund: Makefile-Tab-Einrückung der neuen `--build-arg`-Zeile
  (`Makefile:49`) — Tab-Präfix konsistent mit den Nachbarzeilen.
- geprüft, ohne Befund: Traceability — Commit-Betreff/-Body nennt
  `ADR-0051`, kein `SPEC-*`/`ARC-*`-Präfix im Betreff.
- geprüft, ohne Befund: `make gates` — eigenständig, ungepiped
  nachgefahren (nicht nur die Implementer-Angabe übernommen): Exit-Code
  unmittelbar geprüft, `0`. `generated-sync`/`a-check` liefen sichtbar mit
  im Log, keine Befunde.
- geprüft, ohne Befund: Spec-Stratum, Docker-only, Zwei-Quellen-Drift,
  Suppression-Verbot — keine der vier repo-spezifischen HIGH-Klassen
  greift; keine neue Betreiber-Oberfläche eingeführt (`VERSION` als
  Build-Parameter existierte bereits vor diesem Slice, `ADR-0051`).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** fehlende Negativtests bei neuem
öffentlichem Vertrag · Ausgangsklasse eines Risikos leicht unscharf

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. F-1 (MEDIUM) ist ein realer,
actionable Befund (der Injektionsmechanismus kann exakt auf die Art
regredieren, die dieser Slice gerade behebt, ohne dass ein Gate es fängt),
aber kein Blocker für die Kern-Mechanik: Die `const`→`var`-Notwendigkeit
ist empirisch bestätigt, beide behaupteten Build-Ergebnisse sind
eigenständig reproduziert, `make gates` läuft grün, kein Seiteneffekt im
restlichen Code, keine Kommentar-Disziplin-Verletzung. F-2 (INFO) ist eine
Nuance in bereits vor diesem Commit geschriebenem Plan-Text, kein
Merge-relevanter Fund.

**Übergabe:** Da kein HIGH vorliegt und keine Rückgabe an den Implementer
erfolgt, ziehe ich die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan selbst auf `[x]` nach, mit
Verweis auf diesen Report, im selben Commit, der diesen Report anlegt
(Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde"). F-1 geht bei
Slice-Closure ins Beobachtungs-Register; ob F-1 im Rahmen dieses Slices
oder als Folge-Slice aufgelöst wird, ist eine Planner-/Architect-
Entscheidung. Dieser Report ersetzt keine Verifikation gegen die DoD —
das bleibt Verifier-Aufgabe (Modul 11).
