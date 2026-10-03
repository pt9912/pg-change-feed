# Review-Report: slice-pin-digests-aktualisieren-2026-10-b — 2026-10-03

**Review-Art:** Code — geprüft gegen Plan, `ADR-0146` (Festlegung 2: Index-Digest als Pin-Form;
Festlegung 6: Tag-/Major-Grenze), `ADR-0051` Entscheidung 7 und die `AGENTS.md` Hard Rules
(Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `pin-digests-aktualisieren-2026-10-b`, Diff-Range `7bc4aadd..HEAD`
(11 Commits, 13 Dateien): Pin-Hebungen `b010e244` (nats), `3912b9f0` (dotnet/sdk),
`21d1d168` (dotnet/runtime), `cc2cbddc` (temurin jdk), `f763bb1c` (temurin jre), `457d9fc6`
(python), `f6acabc0` (trivy), `1c298431` (postgres:17), Kommentar-Commit `a4134a96`
(postgres:18-Gewinnungsweg), Plan-Commits `726c085f`, `cfd818c7`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Keine
Mutationsprobe (der Diff trägt keine Zusage, die an einer Eingabeseite rot werden könnte: er
ändert Pin-Werte und Kommentare). Gate-Läufe ungefiltert, Exit je Lauf direkt gelesen
(`AGENTS.md` §3.9). Keine Berechtigungs-Verweigerung im Lauf.

**Eingangs-Kontext:**

- Slice-Plan `slice-pin-digests-aktualisieren-2026-10-b` (§1 Ziel/Abgrenzung, §2 DoD, §3
  Suchlauf und Stopp-Regel, §6 Risiken, §7 Implementer-Belege)
- `ADR-0146`, `ADR-0051`, `ADR-0058`, `ADR-0055` (nats-Pin, Accepted, unberührt)
- berührte Anforderungen: `LH-QA-POR-001`, `LH-QA-POR-003` (Scope)
- `AGENTS.md` (Hard Rules §3.1, §3.5, §3.6, §3.7, §3.9, §3.10, §3.12, §3.13), `harness/conventions.md`

---

## Eigene Messungen des Reviewers (Stand 2026-10-03, Host `x86_64`)

Schwerpunkt 1 und 2 (Index-Digest, Träger-Gleichheit), je Pin
`docker buildx imagetools inspect <image>@<Digest> --raw` (oberster `mediaType`) und
`docker buildx imagetools inspect <Tag> --format '{{.Manifest.Digest}}'`:

| Pin | Form des gepinnten Werts | Tag-Digest == Pin | Digest je Träger (`git grep`, ohne `docs/`) |
|---|---|---|---|
| `nats:2-alpine` | `oci.image.index.v1` | ja | 3 Treffer, 1 Digest (`compose.yaml`, `examples/compose.yaml`, `run-notify-tests.sh`) |
| `dotnet/sdk:10.0` | `docker.distribution.manifest.list.v2` | ja | 2 Treffer, 1 Digest |
| `dotnet/runtime:10.0` | `docker.distribution.manifest.list.v2` | ja | 5 Treffer, 1 Digest |
| `eclipse-temurin:21-jdk` | `oci.image.index.v1` | ja | 2 Treffer, 1 Digest |
| `eclipse-temurin:21-jre` | `oci.image.index.v1` | ja | 5 Treffer, 1 Digest |
| `python:3.14-slim` | `oci.image.index.v1` | ja | 1 Treffer, 1 Digest |
| `aquasec/trivy` (gegen `:latest`) | `oci.image.index.v1` | ja | 1 Treffer, 1 Digest |
| `postgres:17-alpine` | `oci.image.index.v1` | ja | 1 Treffer (`e2e.yml`), 1 Digest |

`postgres:18-alpine`: 8 Treffer, 1 Digest (`77f58511…`), `e2e.yml` PG18-Leg und `Makefile`
`PG_TEST_IMAGE` byte-gleich (`git grep -n 'postgres:18-alpine@'`). Alte Digest-Präfixe
(`065e8355|60a2b223|e6541e52|085eb93e|ca7551d4|caaf356f|62b1e65e|aa90e97e`) stehen in lebenden
Trägern nirgends mehr; Resttreffer ausschließlich in `Accepted` ADRs (`ADR-0055`, `ADR-0087`,
`ADR-0093`, `ADR-0146`) und im Plan. `make suchlauf-nachmessen PLAN=…`: Exit 0, „24 Zeilen
stimmen“. `git grep -n -i "manifest inspect\|amd64"` in lebenden Trägern: nur
`.claude/settings.json` (Berechtigung), `tools/harness/image-stale.sh` (echter Aufruf auf Major-Tag),
Plattform-Listen (`--platform linux/amd64,linux/arm64`, README, Spec) — kein Gewinnungsweg.
Weitere `image@sha256`-Referenzen außerhalb der acht (golang, distroless, uv, d-migrate, a-check,
`golang:1.27`) sind Achsen P1 bis P9; alle `OK` (siehe unten).

Versionen: `nats-server --version` am neuen Pin: v2.15.0 (bestätigt §7). Übrige Versionsangaben
aus §7 (Temurin, Python, .NET) nicht nachgemessen: **übernommen**.

Gelaufen und Exit 0 (je ungefiltert, eigener Lauf): `make gates`, `make fmt-check`, `make test`,
`make sdk-pack-python`, `make image-stale`, `make pin-stale-race`, `-pgtest`, `-dmigrate`,
`-acheck`, `-dcheck`, `-baseline`, `-actions` (jeweils `OK`, keine `DRIFT`), `make image-cve`
(debian 12.15: 0, gobinary: 0 Befunde; deckt die Aussage in §7), `make kommentar-kennungen
DIFF=7bc4aadd` (keine Ausgabe, kein Kandidat; Exit nur hinter `tail` gelesen, daher als Probe).
**Nicht gelaufen (Verifier-Aufgabe, in §7 übernommen):** `make sdk-pack-csharp`/`-kotlin`,
`make examples-csharp`/`-kotlin`, `make test-notify`, `make test-integration`,
`make test-sdk-*-integration`, `make test-replication`, `make test-store`, `make image`.
`git diff 7bc4aadd HEAD -- .github`: nur `e2e.yml`, Zeilen = Pin-Wert PG17 und zwei Kommentare;
Matrix-Struktur, PG18-Leg und `uses:`-Pins unverändert. Versionsdateien der SDKs
(`.csproj`, `pyproject.toml`, `build.gradle.kts`), `.d-check.yml`, `.a-check.yml`, `harness/mk/`:
unberührt (`git diff --stat` leer).

---

## Findings

### F-1 — Messzeilen in §7 sind keine gedruckten Zeilen; Alt-Digests gekürzt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz A; Plan §2 Liefer-Punkt 1 („Die gedruckten Zeilen mit
  vollen Digests und der Parent-Kennung stehen in §7“)
- `pfad`: `docs/plan/planning/in-progress/slice-pin-digests-aktualisieren-2026-10-b.md:424-433`
- `befund`: Die §7-Tabelle nennt Messbefehle und Ergebnisse als Tabelle; die Alt-Digests stehen
  gekürzt (`065e8355…ccc`), die gedruckte Zeile des Laufs fehlt. Ursprung (gemessen) und Befehl
  sind genannt, der Lauf als gedruckte Zeile nicht; ein Nachmessen der Alt-Form ist am
  Endstand nicht mehr möglich, weil kein Träger den alten Pin führt.
- `verifizierbar`: nein (Alt-Zustand nicht mehr im Baum; Neu-Seite von mir nachgemessen)
- `klasse`: Zahl im Träger ohne gedruckte Messzeile

### F-2 — Zwei Parent-Kennungen im Plan

- `kategorie`: INFO
- `quelle`: Maintainability (Plan-Lesbarkeit)
- `pfad`: Plan §3 (Parent `39e27242`, Suchlauf) und §7 (Parent `7bc4aadd`)
- `befund`: Der Suchlauf-Block misst am Parent `39e27242`, die Implementer-Belege nennen
  `7bc4aadd` als Parent; der Plan erklärt den Unterschied nicht. Die Zeilen sind
  nachmessbar und stimmen (24 von 24), die Lesung verlangt aber, beide Kennungen zu kennen.
- `verifizierbar`: nein
- `klasse`: Parent-Kennung uneinheitlich

### F-3 — Vorab-Messung P1 bis P9 fehlt, aber ehrlich vermerkt

- `kategorie`: INFO
- `quelle`: Plan §2 Liefer-Punkt 1 (DoD-Haken bleibt `[ ]`)
- `pfad`: Plan §7, Absatz „Die Achsen-Messung vor der Hebung ist **nicht** vollständig gefahren“
- `befund`: Der Vermerk ist wahr und benennt Ursache (Docker-Hub-Abruflimit, `make image-stale`
  Exit 2) und Umfang (nur nachher gemessen); der Haken ist nicht gesetzt. Die Schlussfolgerung
  „Drift der acht Pins je vorher mit `imagetools inspect` gemessen“ ist durch die Tabelle in §7
  getragen (alt → neu, Form des alten Werts) und trägt die Aussage nur für die **acht** Pins;
  für P1 bis P9 gibt es keine Vorher-Messung, sondern ausschließlich Nachher-`OK` (von mir
  reproduziert: alle Achsen `OK`, Exit 0). Die Planaussage „stehen auf dem Stand des
  Messtags“ (§1) ist damit am Endstand bestätigt, am Parent nicht gemessen.
- `verifizierbar`: ja (`make image-stale`, `make pin-stale-*` am Endstand — `OK`)
- `klasse`: Messung unvollständig vermerkt

### F-4 — Stopp-Regel: nats-Minor und Trivy-Minor

- `kategorie`: INFO
- `quelle`: Plan §3 Stopp-Regel; `ADR-0146` Festlegung 6
- `pfad`: Plan §7 (nats v2.14.6 → v2.15.0; Trivy 0.74.0 → 0.75.0)
- `befund`: Urteil nats: Die Stopp-Regel nennt „die Minor-Linie, die der Tag nennt“ und führt
  `nats` ausdrücklich als „2.x-Linie“; der Tag `2-alpine` nennt nur die Major-Linie 2.
  2.14 → 2.15 liegt innerhalb der genannten Linie, die Regel ist nicht ausgelöst, die
  Behandlung als Normalfall ist regelkonform. Dass eine Minor-Hebung das Laufzeitverhalten
  ändern kann, deckt §6 Risiko 1 und der Läufe-Nachweis (`make test-notify`, `make
  test-integration`, die der Verifier fährt). Trivy: kein Tag-Pin, Pin gegen `:latest`; die
  Minor-Hebung 0.74 → 0.75 ändert ein Prüfwerkzeug; mein Lauf `make image-cve` endet Exit 0 mit
  0 Befunden in beiden Zielen (deckt §7), der Parent-Stand der Befundzahl ist in §7 nicht
  genannt, ein Verhaltensvergleich gegen den Parent steht damit aus.
- `verifizierbar`: ja (`make image-cve`, Exit 0, 0/0 Befunde)
- `klasse`: Stopp-Regel-Auslegung (Linie vs. Tag)

### F-5 — Langläufe als übernommen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A/B; Plan §2 Liefer-Punkt 3 (Haken `[ ]`)
- `pfad`: Plan §7
- `befund`: Die Exit-0-Aussagen zu `make test-notify`, `make test-integration`,
  `make test-sdk-*-integration`, `make examples-*`, `make sdk-pack-csharp`/`-kotlin`,
  `make test-replication`/`-store` und `make image` sind vom Reviewer nicht nachgefahren
  (**übernommen**); der Post-Push-Lauf von `e2e.yml` und `upstream-drift.yml`
  (`AGENTS.md` §3.10) steht laut Plan aus und ist Closure-Pflicht.
- `verifizierbar`: ja (Verifier-Läufe, `gh run list --workflow e2e.yml`)
- `klasse`: Beleg übernommen, nicht nachgefahren

---

## Negativbefunde

- geprüft, ohne Befund: Index-Form und Tag-Gleichheit aller acht gehobenen Digests (Tabelle oben)
- geprüft, ohne Befund: Träger-Gleichheit (genau ein Digest je Image in lebenden Trägern;
  kein Mischstand, keine vergessene Kopie; `postgres:18-alpine` unverändert und byte-gleich in
  `Makefile` und `e2e.yml`)
- geprüft, ohne Befund: Kommentare zur Digest-Gewinnung (Dockerfile-Köpfe, `e2e.yml`, `Makefile`,
  `run-*-tests.sh`): wahr (`--format '{{.Manifest.Digest}}'` liefert den Index-Digest, am Pin
  nachgemessen), konsistent, im Indikativ, ohne Vorher/Nachher-Sprache; je Block höchstens eine
  Kennung; `make kommentar-kennungen DIFF=7bc4aadd` ohne Kandidat
- geprüft, ohne Befund: `.github/workflows/e2e.yml` (nur Pin-Wert und Kommentare)
- geprüft, ohne Befund: Gate-/Schwellen-Konfiguration (§3.6) und SDK-Versionsdateien unberührt
- geprüft, ohne Befund: Docker-only/§3.1 (nur Edit/Write-Änderungen an Werten und Kommentaren, kein
  Host-Werkzeug-Aufruf im Diff erkennbar), §3.9 (Pipe/Wrapper) im Diff nicht berührt
- geprüft, ohne Befund: Traceability (alle Commits nennen `ADR-0146`/`ADR-0051`, kein `SPEC-*`/
  `ARC-*` im Betreff; `make commit-traceability` in `make gates` grün)
- geprüft, ohne Befund: `spec/`, `docs/user/`, `README.md`, `AGENTS.md` ohne Digest (Suchlauf-Zeile 0)
- geprüft, ohne Befund: `ADR-0055`/`ADR-0087`/`ADR-0093` unberührt (§3.5)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Zahl im Träger ohne gedruckte Messzeile · Parent-Kennung
uneinheitlich · Messung unvollständig vermerkt · Stopp-Regel-Auslegung (Linie vs. Tag) ·
Beleg übernommen, nicht nachgefahren

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; F-1 (LOW) und F-2 bis F-5 (INFO) gehen ohne
Rückgabe-Pfeil an den Implementer weiter (Planner nimmt F-1/F-2 bei der Closure-Notiz mit,
der Verifier fährt die in F-5 genannten Läufe und den Post-Push-Lauf). Die Pin-Hebungen sind in
Form (Index-Digest), Träger-Gleichheit und Kommentar-Wahrheit am Endstand belegt.

**Übergabe:** Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
