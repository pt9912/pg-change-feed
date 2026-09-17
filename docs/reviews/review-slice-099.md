# Review-Report: slice-099 — 2026-09-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten): geprüft wird der Diff gegen den Slice-Plan, `ADR-0087`
und `ADR-0090`, nicht gegen die DoD (das ist Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 4f79639..HEAD` — Commits `a62f071`, `e034020`,
`e892a67` (Slice `slice-099`, Plan
`docs/plan/planning/in-progress/slice-099-kotlin-sprachwurzel-http-client.md`).

**Skill:** `.harness/skills/reviewer.md` @ Accepted, Schärfung 2026-09-09 ·
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-099-kotlin-sprachwurzel-http-client.md`
  (§1 Abgrenzung, §2 DoD)
- `ADR-0087` (Sprach-Wurzel-Ort, Bau-Kontext, Digest-Pinning-Tabelle,
  Werkzeug-Ziel-Form), `ADR-0090` (Umfang/Festlegung 4/5/7), `ADR-0073`
  (Zitat-Korrektur-Klasse)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (ADR-Immutabilität), §3.8
  (Action-Pinning), §3.9 (Exit-Code-Disziplin), §3.12/§3.13 (Beleg trägt
  seinen Satz / überholte Träger)
- `harness/conventions.md` (MR-000 ID-Schema)

---

## Findings

### F-1 — `ADR-0087`s Digest-Pinning-Tabelle enthält einen fehlerhaften, zu kurzen sha256-Wert; der Dockerfile-Wert ist real geprüft korrekt

- `kategorie`: HIGH
- `quelle`: `ADR-0087` §Entscheidung Festlegung 3 (Digest-Pinning-Tabelle);
  Klasse „Zahl im Träger ohne Ursprung — oder gegen die Messung driftend"
  (`.harness/skills/reviewer.md`, `AGENTS.md` §3.12)
- `pfad`: `docs/plan/adr/0087-beispiel-clients-csharp-kotlin.md:254` (Tabellenzeile
  „Kotlin | `eclipse-temurin:21-jdk` | …") gegen `examples/kotlin/Dockerfile:16`
- `befund`: Nachgemessen (`echo -n <hex> | wc -c` sowie live
  `docker manifest inspect eclipse-temurin:21-jdk`, amd64/linux-Eintrag):
  der in `ADR-0087` Zeile 254 genannte Digest
  `sha256:085eb93e049c7397f725bd8be31c4f52ba4777a75168e851428508234aa224e`
  hat **63** Hex-Zeichen statt der für sha256 zwingenden 64 — ein
  strukturell ungültiger Digest, dem im Vergleich zum Dockerfile-Wert genau
  ein `d` fehlt (`…c4f52ba…` statt `…c4fd52ba…`). Der in
  `examples/kotlin/Dockerfile:16` tatsächlich gepinnte Digest
  `sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e`
  ist **64 Hex-Zeichen lang und deckt sich exakt** mit dem live abgefragten
  amd64/linux-Manifest-Digest von `eclipse-temurin:21-jdk` (heute gemessen,
  `docker manifest inspect`). Die Diskrepanz ist real, kein Lesefehler des
  Implementers — sie sitzt in `ADR-0087` selbst, nicht im Diff dieses Slice.
  Der Bau des Slice ist davon **nicht** betroffen: er verwendet den
  korrekten, aktuell auflösbaren Digest. Betroffen ist die Verlässlichkeit
  der `ADR-0087`-Tabelle als Beleg — ein künftiger Zug, der ihren Wert statt
  eines eigenen Re-Measurements kopiert, würde einen ungültigen Digest
  pinnen. **Einordnung nach `AGENTS.md` §3.5/`ADR-0073`:** Die Tabelle steht
  innerhalb von `ADR-0087` §Entscheidung (Festlegung 3) — laut `ADR-0073`
  Punkt 1 ist genau dieser Abschnitt vom Zitat-Korrektur-Kanal
  **ausgeschlossen** („keine Zitat-Korrektur, wenn sie … §Entscheidung …
  berührt"); eine Zitat-Korrektur deckt zudem nur Verweis-/Zitatgerüst
  (host-lokale Pfade, Linkziele, Zeilen-Lokatoren) bei unverändertem
  Referenten — ein gemessener Digest-*Wert* ist kein Verweis-Gerüst, sondern
  Teil der Entscheidungssubstanz selbst. Eine Korrektur gehört damit **nicht**
  in einen In-place-Edit, sondern in eine engräumige Folge-ADR mit
  `Supersedes ADR-0087` (oder ein Architect-Verdikt, das diese enge Korrektur
  trägt) — dies ist keine Aufgabe dieses Slices oder seines Reviews, das
  Finding ist die Adresse dafür.
- `verifizierbar`: ja — `echo -n "<hex-teil>" | wc -c` bzw.
  `docker manifest inspect eclipse-temurin:21-jdk` gegen den amd64/linux-Eintrag
- `klasse`: Zahl im Träger ohne Ursprung — oder gegen die Messung driftend

## Negativbefunde

- geprüft, ohne Befund: `examples/kotlin/Dockerfile` — beide Stufen
  digest-gepinnt (`@sha256:…`, nicht nur Tag), Build-Kontext ist
  `examples/kotlin/` (nicht die Repo-Wurzel), Wurzel-`Dockerfile`/
  `.dockerignore` unangetastet.
- geprüft, ohne Befund: `examples/kotlin/gradle/wrapper/gradle-wrapper.properties`
  — `distributionSha256Sum` gesetzt (Prüfsummen-Pinning vorhanden).
- geprüft, ohne Befund: `examples/kotlin/{build.gradle.kts,settings.gradle.kts,
  http-client/build.gradle.kts}` — Kotlin-Gradle-Plugin (`2.4.20`) und
  Testabhängigkeiten (`kotlin-test-junit5:2.4.20`,
  `junit-platform-launcher:6.1.3`) exakt versionsgepinnt, Herkunft
  (Maven-Central-Metadata, Messdatum) im Kommentar genannt.
- geprüft, ohne Befund: `.github/workflows/examples.yml` — der C#-Job
  (`examples` → `examples-csharp`) ist eine reine Umbenennung samt
  Kommentar-Update; seine Steps (`actions/checkout@3d3c42e5aac…` #v7.0.1,
  `make examples-csharp`) sind inhaltlich unverändert. Der neue Job
  `examples-kotlin` nutzt dieselbe gepinnte Checkout-SHA (§3.8 erfüllt), läuft
  unabhängig (kein `needs:`) und ruft ausschließlich `make examples-kotlin`
  auf (keine Ziel-umgehende Inline-Logik).
- geprüft, ohne Befund: `.a-check.yml` — `git diff 4f79639..HEAD -- .a-check.yml`
  ist leer, wie im Slice-Plan §1 ausgeschlossen.
- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-023` Zeile *Sprachen und
  Umfang* — der neue Wortlaut („volle Matrix", benannter Zusatzkontext, Stubs
  im Bau statt committet) zitiert keine `ADR-*`-ID direkt und deckt inhaltlich
  exakt `ADR-0090` Festlegung 1–3, ohne zu über- oder untertreiben.
- geprüft, ohne Befund: `docs/plan/planning/observations/BEO-PGC/
  adr-folgepflicht-ohne-traeger-slice/state.md` — die Begründung, dass die
  Auflösung der konkreten Manifestation (`SPEC-023`-Nachzug) kein neues
  Auftreten der allgemeinen Beobachtungsklasse ist und den Zähler (1×) nicht
  bewegt, trägt: Ein Schließen ist per Registerregel keine Wiederholung.
- geprüft, ohne Befund: `examples/kotlin/http-client/src/test/kotlin/**` —
  netzlos, aussagekräftig (Flag-/Env-Parsing inkl. Fehlerpfade, URL-Aufbau
  inkl. RFC-3986-Escaping, Fehlerpfad ohne erreichbaren Host); keine
  „grün ohne Aussage"-Symptomatik erkennbar.
- geprüft, ohne Befund: `examples/kotlin/http-client/src/main/kotlin/**` —
  `TablesClient` sendet den `reader`-Token korrekt als `Authorization: Bearer`,
  kein Secret-Leak, keine Injection-Fläche; Kommentare bleiben im Indikativ
  über den Ist-Zustand, keine Chronik-Sprache.
- geprüft, ohne Befund: Out-of-Scope-Disziplin — kein SSE/NATS/gRPC in
  Kotlin, kein Kotlin-Bezug in `examples/csharp/**`, kein neues Gate
  (`examples-kotlin` ist Werkzeug, nicht in `GATE_CHECKS`/`harness/README.md`
  §Sensors).
- geprüft, ohne Befund: Commit-Traceability aller drei Commits
  (`a62f071`, `e034020`, `e892a67`) — jeder trägt `ADR-0087, ADR-0090` im
  Footer, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: §3.13-Gegencheck (`grep` nach weiteren „nur
  Go/HTTP-Familie"-Restaussagen) — die verbleibenden Treffer liegen in
  eingefrorenen Records (`done/`, `docs/reviews/**`, ADR-Geschichte), die
  nach `AGENTS.md` §3.5/`ADR-0073` nicht nachgezogen werden; `harness/README.md`
  §Werkzeuge und `docs/user/benutzerhandbuch.md` sind bereits aktuell.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Zahl im Träger ohne Ursprung — oder gegen
die Messung driftend

## Verdikt

**Merge-blockierend:** nein — das einzige Finding (F-1) betrifft einen
Fehler in `ADR-0087` selbst (außerhalb dieses Diffs), nicht den von diesem
Slice gelieferten Code; der im Diff gepinnte Digest ist real korrekt und
aktuell auflösbar. Es liegt außerhalb der Implementer-Zuständigkeit dieses
Slices (ADR-Korrektur ist Architect-/Planner-Arbeit, `AGENTS.md` §3.5) — eine
Fixrunde am Implementer ist nicht nötig. Empfehlung: F-1 als Adresse an die
Architect-Rolle weiterreichen (engräumige Folge-ADR `Supersedes ADR-0087`
für die eine Tabellenzelle, oder Architect-Verdikt).

**DoD-Checkbox-Nachzug:** Da keine Fixrunde am Implementer erfolgt, wird die
Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" im
Slice-Plan im selben Commit wie dieser Report auf `[x]` gezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** Die Finding-Klasse geht in die Slice-Closure §7 und von dort
in den Zähler des Beobachtungs-Registers. Dieser Report ist ein Lauf-Beleg
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und wird
über Läufe hinweg nicht erneut gelesen. Er ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
