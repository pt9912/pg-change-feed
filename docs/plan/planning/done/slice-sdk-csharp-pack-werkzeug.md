# Slice sdk-csharp-pack-werkzeug: `make sdk-pack-csharp` + Träger-Nachzug Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-lh-fa-sst-009](welle-sdk-csharp-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 4 (Build-Mechanismus: Docker-only bauen), §Konsequenzen
Folgepflicht 2/3 (Träger-Nachzug Pflichtenheft/`harness/README.md`).

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(wird mit diesem Slice für C#/NuGet aufgelöst — die Kennung selbst bleibt
bestehen, ihr Satz „ist offen, ADR-pflichtig" wird durch `ADR-0106` und die
jetzt reale Paketierbarkeit falsch, `AGENTS.md` §3.13), `spec/pflichtenheft.md`
§6 Externe Verträge (neue Zeile für das Package).

**Verantwortlich:** Implementer-Agent (dietmar.burkard@nerdware.dev), ab
2026-09-19.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0106` §Konsequenzen
Folgepflicht 1/2/3). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, netzlos **nicht** prüfbares Werkzeug-Ziel
`make sdk-pack-csharp` (Docker-only, `dotnet build`/`dotnet test`/
`dotnet pack`, gepinntes `mcr.microsoft.com/dotnet/sdk`-Image) erzeugt ein
reales `.nupkg` aus `sdks/csharp/PgChangeFeed.Client/` — dazu der
Träger-Nachzug, den `ADR-0106` §Konsequenzen Folgepflicht 2/3 fordert:
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (§1) und `spec/pflichtenheft.md`
§6 (Externe Verträge, neue `SPEC-<NNN>`-Zeile für das Package) sowie
`harness/README.md` §Werkzeuge (die reale `make sdk-pack-csharp`-Zeile).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **NuGet.org-Veröffentlichung** — `slice-sdk-csharp-publish-workflow`
  übernimmt das; dieser Slice erzeugt das `.nupkg`-Artefakt, veröffentlicht
  es nicht (`ADR-0106` Festlegung 4: „Bauen und `pack` bleiben Docker-only
  … Veröffentlichung … braucht einen eigenen, Netz-bindenden Workflow").
- **Neues Gate** — `make sdk-pack-csharp` bleibt außerhalb von
  `GATE_CHECKS`/`make gates` (`ADR-0106` Festlegung 4, dieselbe
  Begründung wie `make examples-csharp`: Paket-Bezug braucht Netz).
- **Änderung der beiden Fläche-Slices** — dieser Slice paketiert, was
  `slice-sdk-csharp-http-client-flaeche`/`slice-sdk-csharp-grpc-client-flaeche`
  bereits geliefert haben; kein neuer API-Umfang.
- **`docs/user/version.md`-Änderung** — die Server-Version bleibt
  eigenständig geführt (`ADR-0106` Festlegung 3/5).

## 2. Definition of Done

- [x] `make sdk-pack-csharp` existiert (Docker-only, kein Gate — analog
      `make examples-csharp`): baut, testet (`dotnet test` gegen die
      Test-Projekte beider Fläche-Slices) und paketiert
      (`dotnet pack -c Release`) `sdks/csharp/PgChangeFeed.Client/` im
      gepinnten `mcr.microsoft.com/dotnet/sdk`-Image; ein roter Test bricht
      den `docker build` mit Exit ≠ 0 ab (Muster
      `harness/mk/examples.mk`). Das `.nupkg` wird über einen `docker cp`-
      oder Bind-Mount-Export aus dem Build in ein lokales Verzeichnis
      (z. B. `sdks/csharp/dist/`, `.gitignore`t) abgelegt — analog dem
      Extraktionsmuster von `make proto-generate`
      (`tools/harness/proto-generate.sh`, host-seitiger Export statt
      Bind-Mount/`--user`-Workaround).
- [x] Real ausgeführt: ein `.nupkg` mit Dateiname
      `PgChangeFeed.Client.0.1.0.nupkg` existiert nach dem Lauf und ist
      als Smoke-Beleg im Bericht dieses Slice genannt (Datei-Existenz,
      keine Behauptung — `AGENTS.md` §3.12 Instanz B).
- [x] `spec/pflichtenheft.md` §1 trägt bei
      [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) einen
      Nachzug-Satz: für C#/NuGet ist die Sprachmatrix-/Vertriebsweg-Frage
      durch [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
      beantwortet und das Package real paketierbar; eine zweite Sprache
      oder ein zweiter Vertriebsweg bleibt offen (`AGENTS.md` §3.13, kein
      Streichen der Kennung — sie bleibt die Adresse für künftige
      Sprachen/Vertriebswege).
- [x] `spec/pflichtenheft.md` §6 Externe Verträge bekommt eine neue
      `SPEC-<NNN>`-Zeile für `PgChangeFeed.Client` (System: NuGet-Package;
      Version: SemVer 2.0, `0.x.y`; Vertrag-Datei: Verweis auf
      `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` als
      Metadaten-Quelle, analog der bestehenden `SPEC-023`-Zeile für die
      Beispiel-Client-Werkzeugketten) sowie §7 Historie-Zeile.
- [x] `harness/README.md` §Werkzeuge bekommt die reale
      `make sdk-pack-csharp`-Zeile (kein Gate, Bindung auf
      [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)) —
      erst jetzt zulässig, weil das Ziel jetzt real existiert
      (`AGENTS.md` §4).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md`
      (0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO — keine Fixrunde nötig,
      DoD-Checkbox-Nachzug ohne Fixrunde laut Reviewer-Skill).
- [x] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben) — entfällt
      als eigener Punkt, da bereits oben als DoD-Kriterium geführt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-lh-fa-sst-009](welle-sdk-csharp-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/mk/sdk.mk` (neu, Arbeitsname) oder `harness/mk/examples.mk` (erweitert) | neu/update | `sdk-pack-csharp`-Target, Docker-only, kein `GATE_CHECKS`-Eintrag. |
| `sdks/csharp/Dockerfile` | update | zusätzliche Bau-/Export-Stufe für `dotnet pack`, Extraktion analog `make proto-generate`. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) | update | Nachzug-Satz: C#/NuGet nicht mehr offen, zweite Sprache/Vertriebsweg bleibt offen. |
| `spec/pflichtenheft.md` §6 | update | neue `SPEC-<NNN>`-Zeile für `PgChangeFeed.Client`. |
| `spec/pflichtenheft.md` §7 Historie | update | Historie-Zeile für die neue `SPEC-<NNN>` und den `LH-FA-SST-009.a`-Nachzug. |
| `harness/README.md` §Werkzeuge | update | reale `make sdk-pack-csharp`-Zeile. |

**Plan-Nachzug (Implementer-Zug):**

- Die `spec/pflichtenheft.md` §1-Nachzug-Zeile bei `LH-FA-SST-009.a` nennt
  `ADR-0106` **nicht** namentlich — abweichend vom ursprünglichen
  DoD-Wortlaut oben, der die Nennung vorsah. Grund: `make docs-check`s
  `matrix`-Modul verbietet mechanisch jede Referenz `spec → adr`
  (`.d-check.yml` `matrix.rules: {from: spec, to: adr, allow: false}`,
  Decken-Regel) — real getroffen: ein erster Versuch mit Markdown-Link auf
  `ADR-0106` erzeugte den Befund `spec/pflichtenheft.md:183
  ../docs/plan/adr/0106-....md matrix-forbidden Referenz spec → adr ist
  nicht erlaubt`. Die Nachzug-Zeile trägt den fachlichen Inhalt (Frage
  beantwortet, `SPEC-026`, real paketierbar) ohne den ADR-Bezug; die
  Beziehung selbst steht bereits umgekehrt in `ADR-0106`s eigenem
  `Schärft:`-Feld. Dieselbe Regel griff versucht auch für die neue
  `SPEC-026`-Zeile in §6 — dort ebenfalls ohne ADR-Verweis geschrieben,
  analog der bestehenden `SPEC-023`/`SPEC-020`/`SPEC-017`-Zeilen, die auch
  keinen ADR-Bezug tragen.
- `harness/mk/sdk.mk` (neu angelegt, nicht `harness/mk/examples.mk`
  erweitert) — ein eigenständiges Fragment hält das SDK-Pack-Werkzeug von
  den Beispiel-Werkzeugketten getrennt (unterschiedliche ADRs: `ADR-0106`
  vs. `ADR-0087`/`ADR-0090`), analog wie `harness/mk/*.mk` bereits je
  Belang aufgeteilt ist.
- `sdks/csharp/Dockerfile` bekommt zwei neue Stufen (`pack`, `pack-export`)
  statt einer einzelnen — `pack-export` trägt ausschließlich das
  `ENTRYPOINT`, ohne eigenen `RUN`-Schritt, exakt das Muster von
  `proto`/`proto-export` im Wurzel-Dockerfile.
- Ein roter Test wurde real erzeugt und geprüft (mutierte Testerwartung in
  `PgChangeFeedClientOptionsTests.cs`, danach exakt auf den committeten
  Stand zurückgesetzt — `git diff --stat` bestätigt keine Restdifferenz):
  `docker build` bricht mit Exit 1 an der `dotnet test`-Stufe ab, kein
  `.nupkg` wird exportiert (§6 Risiko 1 damit zusätzlich zur reinen
  Erfolgs-Bestätigung auch am Fehlerpfad real geprüft).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-csharp-http-client-flaeche`
**und** `slice-sdk-csharp-grpc-client-flaeche` in `done/` liegen (siehe
Welle-Plan §4 Reihenfolge — `ADR-0106` Festlegung 1: ein Package aus beiden
Flächen).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Make-Target plus zwei Doku-Nachzüge.
- `in-progress` → `open` (blockiert — Carveout?): `dotnet pack` scheitert
  strukturell am gepinnten SDK-Image (z. B. fehlende `net10.0`-Pack-
  Unterstützung) — unwahrscheinlich, `dotnet pack` ist Standard-Tooling
  jeder `dotnet`-SDK-Version.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reales `.nupkg` als Smoke-Beleg +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der Export-Mechanismus für das `.nupkg` aus dem Docker-Bau (analog
  `make proto-generate`s `tar`-Stream-Export) ist neu für diesen
  Artefakt-Typ — ein Bind-Mount-Workaround oder ein falsch gesetztes
  `--user` könnte das Artefakt mit falschen Dateirechten oder gar nicht
  auf den Host bringen. **Ausgang:** entfallen/aufgelöst — der
  `tar`-Stream-Export (`docker run --rm --network none <image> | tar -x
  -C .`, `AGENTS.md` §3.9 Pipe-Disziplin) trug beim ersten Versuch real
  und wurde danach dreifach unabhängig bestätigt (Implementer-,
  Reviewer-, Verifier-Lauf; der Verifikationsbericht §1.1 zählt insgesamt
  sechs eigenständige Bau-Läufe über den Zyklus, inklusive einer eigenen,
  vom Reviewer unabhängigen Mutations-Probe) — kein Bind-Mount, kein
  `--user`-Workaround nötig, keine Dateirechte-Auffälligkeit in keinem der
  Läufe (`git status --porcelain` je leer/`.gitignore`-konform nach dem
  Export).
- Die neue `SPEC-<NNN>`-Nummer in `spec/pflichtenheft.md` §6/§2 muss
  fortlaufend vergeben werden (aktuell zuletzt `SPEC-024`) — ein
  Parallel-Slice, der ebenfalls eine neue `SPEC-*`-Nummer vergibt, könnte
  zu einer Kollision führen. **Ausgang:** entfallen — `SPEC-026` real
  kollisionsfrei vergeben; Reviewer und Verifier haben je unabhängig
  `grep -c "SPEC-026" spec/*.md` gefahren (drei Treffer, ausschließlich in
  `spec/pflichtenheft.md`: §1-Fließtext, §6-Zeile, §7-Historie; `0` in
  `spec/lastenheft.md`/`spec/architecture.md`) — genau eine
  Definitionsstelle, keine Zweitvergabe.

## 7. Closure-Notiz

- **Was hat funktioniert:** Ein glatter Durchlauf ohne Fixrunde — der
  Reviewer fand 0 HIGH/MEDIUM/LOW und 1 rein kosmetisches INFO
  (`docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md`, F-1: Form der
  `SPEC-026`-Vertrag-Datei-Spalte weicht sachlich begründet vom
  Nachbarzeilen-Muster ab, kein Fix nötig), der Verifier bestätigte „DoD
  erfüllt" mit vollständig eigenständigen Läufen — inklusive einer
  eigenen, vom Reviewer unabhängigen Mutations-Probe (andere Testdatei,
  andere Assertion) mit identischem Ergebnis: ein roter Test verhindert
  das `.nupkg`-Artefakt strukturell, nicht nur im vom Reviewer geprüften
  Einzelfall. Der neue Export-Mechanismus für einen `.nupkg` (Docker-Bau
  → `tar`-Stream-Export, analog `make proto-generate`) trug beim ersten
  Versuch real, ohne Bind-Mount-/`--user`-Nacharbeit. Der
  DoD-Checkbox-Nachzug ohne Fixrunde (`BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde`)
  griff erneut korrekt: der Reviewer zog die DoD-Zeile „Review
  durchgeführt" im selben Commit selbst nach.
- **Was ging anders als geplant:** Der ursprüngliche DoD-Wortlaut (§2)
  sah eine namentliche `ADR-0106`-Nennung in der neuen
  `LH-FA-SST-009.a`-Nachzugzeile und der neuen `SPEC-026`-Zeile vor. Ein
  erster Versuch mit Markdown-Link erzeugte real den Befund
  `matrix-forbidden Referenz spec → adr ist nicht erlaubt`
  (`.d-check.yml` `matrix.rules`) — der Implementer schrieb beide Zeilen
  ohne ADR-Bezug (Plan-Nachzug, §3) und beließ die fachliche Aussage
  unverändert. Reviewer und Verifier prüften das unabhängig voneinander
  gegen die bereits im Dokument etablierten §3-/§7-Kopfregeln
  („kein ADR- und kein Slice-Verweis", die Beziehung steht bereits
  umgekehrt in `ADR-0106`s `Schärft:`-Feld) und bestätigten: keine
  Ad-hoc-Notlösung, sondern dieselbe Struktur-Regel wie bei
  `SPEC-010`/`SPEC-017`/`SPEC-020`/`SPEC-023`/`SPEC-024`.
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine geschärfte Regel.
  Zwei mögliche Beobachtungs-Kandidaten geprüft (siehe unten) — beide ohne
  neue Handlung.
- **Beobachtungs-Register (`../observations/`):** keine neue
  Beleg-Datei angelegt, zwei Kandidaten geprüft:
  - a) Der Auftrag benannte einen möglichen dritten Fall in dieser Welle,
    in dem Reviewer **und** Verifier unabhängig voneinander denselben
    `docs-check`-`id-unlinked`-Stolperstein beim Schreiben ihrer eigenen
    Berichte getroffen hätten (Klasse `BEO-PGC/report-nackte-id-ohne-link`,
    bereits verkörpert seit `slice-063`, zuletzt Beleg 7 bei
    `slice-sdk-csharp-grpc-client-flaeche`). Eigene Prüfung: weder der
    Review- noch der Verifikationsbericht dieses Slice nennen
    `id-unlinked` an irgendeiner Stelle, und anders als beim
    Vorgänger-Slice (dortiger Fix-Commit `899a3f80`) existiert für diesen
    Slice **kein** separater Fix-Commit zwischen Report-Entwurf und
    -Commit (`git log` zeigt genau je einen Commit für Review-
    (`69d2dc00`) und Verifikationsbericht (`500fdc85`), keinen weiteren).
    Die behauptete dritte Instanz lässt sich aus den vorliegenden
    Artefakten **nicht** unabhängig bestätigen — Feststellung ohne neue
    Beleg-Datei; die bereits verkörperte Regel (`AGENTS.md` §3.9, 7
    bestehende Belege) bleibt unverändert und ausreichend, ein
    unbestätigter Vorgang wird nicht als Beleg eingetragen
    (`AGENTS.md` §3.12).
  - b) Der Implementer musste eine im ursprünglichen Plan vorgesehene
    ADR-Referenz streichen, weil `docs-check`s `matrix`-Modul sie
    strukturell verbietet (siehe „Was ging anders als geplant" oben).
    `grep` über `docs/plan/planning/observations/*/*/observation.md`
    nach „verbietet"/„verboten"/„matrix" fand zwei benachbarte, aber
    andere Klassen (`BEO-PGC/fitness-function-gegen-eigene-entscheidung`:
    ADR widerspricht sich selbst; `BEO-PGC/regel-weiter-als-ihr-sensor`:
    eine Zusage ist weiter gefasst als ihr Sensor) — keine trifft „Plan
    nimmt eine Referenzform an, die ein bestehender Sensor bereits
    verbietet". Entscheidung: **keine neue Beobachtung** — einmaliger
    Fehltritt, ohne Folgekosten über den bestehenden, bereits
    verkörperten Plan-Nachzug-Mechanismus (§3 dieses Slice-Plans)
    korrekt aufgefangen: Der Implementer dokumentierte die Abweichung
    inline, Reviewer und Verifier bestätigten sie unabhängig als
    regelkonform. Kein wiederkehrendes Muster, kein Gate-Umweg, keine
    Fixrunde — die Fehlerklasse bräuchte ein zweites Auftreten, um eine
    eigene Beobachtung zu rechtfertigen.
- **Folge-Slices:** `slice-sdk-csharp-publish-workflow` — bereits als
  Datei in `open/` vorhanden.
- **Risiken aus §6:**
  - „Export-Mechanismus für `.nupkg` neu, Bind-Mount-/`--user`-Risiko" —
    **Ausgang: entfallen/aufgelöst** — dreifach unabhängig real bestätigt
    (Implementer, Reviewer, Verifier), keine Dateirechte-Probleme
    aufgetreten.
  - „Neue `SPEC-<NNN>`-Nummer könnte kollidieren" — **Ausgang: entfallen**
    — `SPEC-026` real kollisionsfrei vergeben, mehrfach gegengeprüft
    (`grep -c` durch Reviewer und Verifier, je genau drei Treffer
    ausschließlich in `spec/pflichtenheft.md`).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-lh-fa-sst-009](welle-sdk-csharp-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Zwei Sub-Areas: `sdks/csharp/`
(bereits eröffnet, GF) und `spec/pflichtenheft.md`/`harness/README.md`
(Default-Sub-Area `*`/`PGC`, Greenfield, `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/plan-nachzug`/`BEO-PGC/adr-folgepflicht-ohne-traeger-slice` geprüft
(Namen legen einen Bezug nahe): beide betreffen andere Fehlerklassen
(Slice-Plan-Nachzug bzw. ADR-Folgepflicht ohne jeden Träger-Slice) — dieser
Slice **ist** der Träger-Slice, kein Treffer, der ihn selbst beträfe.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
