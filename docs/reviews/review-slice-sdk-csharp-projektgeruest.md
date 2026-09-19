# Review-Report: slice-sdk-csharp-projektgeruest — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-projektgeruest.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten).

**Gegenstand:** Diff-Range `8606c8a1..f43f03f0` (Welle-Eröffnung bis
Implementer-Commit), Slice `slice-sdk-csharp-projektgeruest`,
Welle `welle-sdk-csharp-lh-fa-sst-009`.
Vier Commits: `810944df` (open→next, reiner Move), `4a19a505`
(Verantwortlich gesetzt, nur Slice-Datei), `348a7cf8` (next→in-progress,
reiner Move), `f43f03f0` (Inhalt: neuer Baum `sdks/csharp/`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-projektgeruest.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0106` (Accepted) — Festlegung 1/2/3/5, §Kontext Bindung 4
  (a-check liest kein C#)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.11 (kein host-lokaler Pfad), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `examples/csharp/{Dockerfile,Directory.Packages.props,http-client/*.csproj,grpc-client/*.csproj}`
  als Vorbild-Referenz (`ADR-0106` §Kontext Ist-Stand)
- `docs/reviews/review-slice-release-multi-arch-image.md` — Formvorbild
  für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message
übernommen):**

- `docker build --no-cache -f sdks/csharp/Dockerfile sdks/csharp` real
  ausgeführt (nicht nur die gecachte Variante): Exit-Code direkt (ungepiped)
  geprüft, `0`. `dotnet build` 0 Warnings/0 Errors, `dotnet test`: 5/5
  Tests grün — bestätigt die Implementer-Behauptung in der Commit-Message
  wortgleich. Test-Image danach mit `docker rmi` entfernt.
- `make docs-check` real ausgeführt, Exit-Code direkt geprüft: `0`
  (`d-check: 797 Datei(en) geprüft, 0 Befund(e)`).
- `git show 810944df --stat`/`git show 348a7cf8 --stat` geprüft: beide
  zeigen ausschließlich den Rename `{open,next} => {next,in-progress}`,
  0 Insertions/Deletions — reine Moves. `git show 4a19a505 --stat`
  geprüft: ändert ausschließlich die Slice-Datei selbst (1 Zeile) — das
  3-Commit-Muster (Move · Inhalt · Move) ist sauber getrennt, keine
  Move+Inhalt-Vermischung.
- `grep -rn "internal/\|cmd/pg-change-feed" sdks/` — kein Treffer (Exit 1).
  `grep -n "PackageReference\|ProjectReference"` im
  `PgChangeFeed.Client.csproj` — kein Treffer; nur das Test-Projekt trägt
  ein `ProjectReference` auf die eigene Geschwister-`.csproj` (zulässig,
  kein privater Repo-Baum).
- `cat PgChangeFeed.Client.csproj` Feld für Feld gegen die DoD-Liste
  geprüft: `PackageId=PgChangeFeed.Client` ✓, `Version=0.1.0` ✓,
  `Description` ✓, `Authors` ✓, `PackageLicenseExpression=MIT` ✓ (`LICENSE`
  im Repo-Root real gegengelesen — MIT, Copyright pt9912 2026) ✓,
  `PackageProjectUrl`/`RepositoryUrl` ✓ (beide auf dieses Repo),
  `PackageReadmeFile=README.md` ✓ (referenziert über `../README.md` mit
  explizitem `PackagePath`).
- `sdks/csharp/Directory.Packages.props` gegen `examples/csharp/Directory.Packages.props`
  diff-artig verglichen: exakt dieselben drei xUnit-Versionen
  (`Microsoft.NET.Test.Sdk 18.10.1`, `xunit 2.9.3`,
  `xunit.runner.visualstudio 4.0.0`) — keine neue, eigenständig zu
  bewertende Fremdabhängigkeit, wie der Plan-Nachzug in §3 behauptet.
- Suchlauf nach host-lokalen absoluten Pfaden (Wurzel-Segment-Präfixe
  eines Entwicklerrechners bzw. Windows-Laufwerksmuster, `AGENTS.md`
  §3.11) über `sdks/csharp/` und die Slice-Datei — kein Treffer.
- `git diff 8606c8a1..HEAD --stat -- .a-check.yml harness/README.md spec/
  AGENTS.md harness/conventions.md docs/user/version.md` — leer; bestätigt
  `ADR-0106` §5 „Was diese ADR nicht ändert" wortgleich für diesen Slice.
- `PgChangeFeedClientOptionsTests.cs` gelesen: trägt ein explizites
  `using Xunit;` — der vom Implementer berichtete erste rote
  Docker-Build (ImplicitUsings deckt `Xunit` nicht ab) ist in der finalen
  Fassung sauber behoben, kein Rest-Fehler, kein zusätzlicher globaler
  Using-Eintrag im `.csproj` nachgezogen.
- `git status --short` nach allen eigenen Läufen leer — kein
  Test-Artefakt im Arbeitsbaum zurückgelassen.

---

## Findings

### F-1 — Deutsches Wortfragment in einem englischsprachigen öffentlichen XML-Doc-Kommentar

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs:5`
- `befund`: Der Klassen-Doc-Kommentar der öffentlichen Klasse
  `PgChangeFeedClientOptions` ist durchgängig Englisch, enthält aber das
  deutsche Wort „unstrittige" mitten im englischen Satz („this is the one
  unstrittige, shared configuration denominator identified while writing
  this project skeleton"). Dieser Kommentar ist Teil der öffentlichen
  API-Dokumentation eines NuGet-Packages (IntelliSense-Tooltip für
  Konsumenten) — ein Sprachbruch ist dort sichtbarer als in einem
  internen Kommentar.
- `verifizierbar`: nein — reine Stilfrage, kein Gate prüft
  Kommentar-Sprache.
- `klasse`: Sprachbruch in öffentlichem Doc-Kommentar (Erstauftreten)

## Negativbefunde

- geprüft, ohne Befund: Docker-only-Disziplin (`AGENTS.md` §3.1) —
  `sdks/csharp/Dockerfile` trägt `dotnet restore`/`build`/`test`
  ausschließlich im gepinnten Container; kein Skript, kein
  `Makefile`-Target und keine Doku-Stelle unter `sdks/` setzt einen
  Host-`dotnet`-Aufruf voraus; `.gitignore`-Kommentar benennt den
  Host-Bau ausdrücklich nur als Schutz gegen ein Versehen, nicht als
  vorgesehenen Pfad.
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als zwei/drei getrennte
  Commits (`AGENTS.md` §3.3) — die beiden reinen Lifecycle-Moves
  (`810944df`, `348a7cf8`) tragen 0 Insertions/Deletions, der
  Inhalts-Commit dazwischen (`4a19a505`) ändert ausschließlich das
  `Verantwortlich:`-Feld der Slice-Datei; das gewählte 3-Commit-Muster
  ist eine zulässige, sauber getrennte Erweiterung des Grundmusters.
- geprüft, ohne Befund: Import-Grenze (`ADR-0106` Festlegung 2, §Kontext
  Bindung „Import-Grenze, hier ohne Ausnahme") — kein Treffer für
  `internal/`/`cmd/pg-change-feed` unter `sdks/`, keine `PackageReference`/
  `ProjectReference` auf einen privaten Repo-Baum.
- geprüft, ohne Befund: NuGet-Metadaten-Vollständigkeit des `.csproj`
  gegen die DoD-Liste — jedes verlangte Feld einzeln geprüft (siehe oben),
  keine Lücke.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — die
  Dockerfile-/`.csproj`-/`.cs`-Kommentare beschreiben durchgängig den
  geltenden Zustand oder eine Abgrenzung zu Folge-Slices („Sache von
  `slice-sdk-csharp-pack-werkzeug`"); keine Konjunktiv-Begründung über
  eine verworfene Alternative, kein abwesender Text. Der bare
  Slice-Name im Dockerfile-Kopfkommentar (`slice-sdk-csharp-projektgeruest`)
  folgt demselben, bereits akzeptierten Muster wie
  `examples/csharp/Dockerfile`s „(…, slice-102)" — kein neues Muster,
  keine Vorher/Nachher-Sprache, deshalb kein HIGH nach der
  Chronik-in-Produktionscode-Regel.
- geprüft, ohne Befund: `README.md` dupliziert nicht die kanonische
  Draht-Doku — es verweist auf `SPEC-018`/`SPEC-020` bzw. das
  Repo-Root-`README.md`, statt Endpunkte/Protokolldetails selbst zu
  beschreiben; ist durchgängig Englisch (Auftragsvorgabe).
- geprüft, ohne Befund: host-lokale absolute Pfade (`AGENTS.md` §3.11) —
  kein Treffer unter `sdks/csharp/` oder der Slice-Datei.
- geprüft, ohne Befund: Träger, die `ADR-0106` §5 ausdrücklich als
  „bleibt unberührt" benennt (`.a-check.yml`, `harness/README.md`,
  `spec/**`, `AGENTS.md`, `harness/conventions.md`,
  `docs/user/version.md`) — tatsächlich in diesem Diff unangetastet.
- geprüft, ohne Befund: Scope-Treue gegen §1 „Ausdrücklich NICHT in
  diesem Slice" — kein `make sdk-pack-csharp`, kein Publish-Workflow,
  kein Umbau von `examples/csharp/`, kein Pflichtenheft-Träger-Nachzug
  im Diff.
- geprüft, ohne Befund: Traceability — Commit-Betreff `f43f03f0` nennt
  `LH-FA-SST-009` und `ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff;
  `harness/conventions.md` MR-002 (Slice-Kennungen sind Namen) —
  `slice-sdk-csharp-projektgeruest` ist bereits ein Name.
- geprüft, ohne Befund: reale Docker-Build- und Testausführung — 5/5
  Tests grün, Build ohne Warnungen (die einzige Docker-Meldung ist
  `InvalidBaseImagePlatform`, ein plattform-bedingter Buildx-Hinweis
  wegen amd64-gepinntem Image auf arm64-Host, identisch zum bereits
  bestehenden `examples/csharp/Dockerfile`-Verhalten — kein Befund, kein
  neues Muster).
- geprüft, ohne Befund: `make gates`-Nachweis — nicht separat erneut
  gefahren (dieser Review-Lauf beschränkt sich auf `make docs-check`
  als engsten einschlägigen Sensor für einen reinen Doku-/neuen-Sprach-
  Baum-Diff, plus den realen Docker-Build als Gegenprobe zur
  Implementer-Behauptung); `make docs-check` real grün, siehe oben.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Sprachbruch in öffentlichem
Doc-Kommentar (Erstauftreten).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. F-1 ist eine einmalige
stilistische Beobachtung ohne semantische Auswirkung (LOW), keine
Fixrunde erforderlich.

**Übergabe:** Keine Rückmeldung an den Implementer nötig — F-1 kann bei
Gelegenheit (z. B. im nächsten Slice, der dieselbe Datei berührt)
mitgezogen werden, ist aber keine eigene Fixrunde wert. Da 0 HIGH
vorliegen (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde"), ziehe ich
die DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" im Slice-Plan selbst auf `[x]` nach, im selben Commit wie
diesen Report. Dieser Report ersetzt keine Verifikation gegen die DoD —
das bleibt Verifier-Aufgabe (Modul 11).
