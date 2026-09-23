# Review-Report: slice-sdk-kotlin-reale2e — 2026-09-23

**Review-Art:** Code — geprüft gegen Plan
([`slice-sdk-kotlin-reale2e`](../plan/planning/done/slice-sdk-kotlin-reale2e.md),
§2 DoD, §3 Plan + Plan-Nachzug + §3.13-Suchlauf-Feld, §6 Risiken),
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Accepted) §Entscheidung Festlegung 2, [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(Ort/Import-Grenze), [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
(benannter Bau-Kontext `proto`), `AGENTS.md` §3 Hard Rules (§3.2, §3.7,
§3.9, §3.11, §3.12, §3.13). Kein DoD-Abgleich als solcher — Verifier-Aufgabe
(Modul 11); wo eine DoD-Zeile eine **Zahl** behauptet, ist das Nachmessen
dieser Zahl trotzdem Reviewer-Aufgabe (`AGENTS.md` §3.12 Instanz A,
Reviewer-Skill HIGH „Zahl im Träger … driftend").

**Gegenstand:** Diff-Range `b58cb173..HEAD` — Commit `060aa62f` des
Slices `slice-sdk-kotlin-reale2e` (12 Dateien, 825 Insertions /
10 Deletions): Runner-Skript `tools/harness/run-sdk-kotlin-integration-tests.sh`,
Docker-Stufe `integration` (`sdks/kotlin/Dockerfile`), Gradle-SourceSet
`integrationTest` + Task (`build.gradle.kts`), fünf Integrations-Testdateien,
Make-Target `test-sdk-kotlin-integration` (`harness/mk/sdk.mk`),
Abdeckungs-Träger `docs/user/sdk-e2e-abdeckung.md` (Kotlin-Abschnitt),
`harness/README.md`-Zeile, Plan-Update (DoD-Häkchen, Plan-Nachzug,
§3.13-Feld).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09".
**Modell:** glm-5.3-flash (Claude-Agent-SDK-Subagent) · **Datum:** 2026-09-23.

**Eingangs-Kontext:**

- Plan (§1–§8), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  Festlegung 2 (Mechanik-Klasse), [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md)
  Festlegung 1/3/5 (Ort, Import-Grenze), [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
  Festlegung 2, [`welle-sdk-reale2e.md`](../plan/planning/welle-sdk-reale2e.md)
  (§1 Träger-Erzeugnis, §4 Reihenfolge, §6 Beobachtungen-Sicht)
- `spec/lastenheft.md` [`LH-FA-SST-009`](../../spec/lastenheft.md),
  [`LH-FA-SST-008`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md),
  [`LH-FA-CON-001`](../../spec/lastenheft.md); `spec/pflichtenheft.md`
  `SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-024`/`SPEC-027` (gelesen)
- Vorbilder (gelesen, nicht importiert):
  `tools/harness/run-sdk-python-integration-tests.sh` (Origin),
  `tools/harness/run-sdk-csharp-integration-tests.sh` (Spiegel),
  `sdks/csharp/PgChangeFeed.Client.Integration/*.cs`,
  [`docs/reviews/review-slice-sdk-csharp-reale2e.md`](review-slice-sdk-csharp-reale2e.md)
  (Lernklassen-Kette F-1 Suchlauf-Verfeinerung, F-2 SST-009-Zitierpflicht,
  F-3 Form-Fragment, F-7 Writer-Form-Grenze, F-8 Zahl-Ursprung)
- `compose.yaml` Container-Vertrag (Token-/Publikations-/Slot-Werte gegen
  die Runner-Umgebungsvariablen gehalten), `AGENTS.md` (Hard Rules)

---

## Findings

### F-1 — §3.13-Suchlauf-Feld nennt einen Träger, der nicht existiert

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug) · Reviewer-Skill HIGH
  „Beleg trägt seinen Satz nicht" (Variante: Adresse auf ein Artefakt,
  das es nicht gibt — Nachbar-Form `BEO-PGC/zitat-nennt-die-falsche-stelle`)
- `pfad`: `docs/plan/planning/done/slice-sdk-kotlin-reale2e.md:189`
- `befund`: Die Geprüft-Zeile des committeten Suchlauf-Felds nennt als
  Träger „`sdks/kotlin/README.md` §Status" — auf dieser Ebene existiert
  keine README (Inhalt von `sdks/kotlin/`: `dist/`, `Dockerfile`,
  `pgchangefeed-kotlin/`); der reale Träger ist
  `sdks/kotlin/pgchangefeed-kotlin/README.md` (dort §Status, auch so vom
  Benutzerhandbuch zitiert). Die Aussage der Zeile („trägt keine
  Teststrategie-Aussage über Realserver-Läufe") ist gegen die reale Datei
  nachfahrbar und hält, aber die Adresse des Suchlauf-Belegs ist
  unauflösbar — das Feld ist das Übergabe-Artefakt des Suchlaufs, seine
  Prüf-Angabe muss auflösbar sein.
- `verifizierbar`: ja — `ls sdks/kotlin/` (kein `README.md`); der
  getragene Check ist gegen `sdks/kotlin/pgchangefeed-kotlin/README.md`
  nachfahrbar.
- `klasse`: „Beleg trägt seinen Satz nicht (Adresse ohne Artefakt)"

### F-2 — Form-Fragment „FlaecheN" wortgleich aus dem C#-Spiegel übernommen

- `kategorie`: LOW
- `quelle`: Maintainability (Plan §6 Risiko 1:
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` — die
  übernommenen Form-Teile werden je separat auf Sprachreinheit geprüft)
- `pfad`: `tools/harness/run-sdk-kotlin-integration-tests.sh:6`
- `befund`: „die vier FlaecheN des Packages" — Großbuchstabe mitten im
  Wort, byte-identisch aus dem C#-Spiegel-Kopf übernommen (dort dieselbe
  Form; der C#-Review F-3 traf dort die Nachbarform „die vier
  Flaeche"). Die im §6-Risiko zugesagte separate Sprachreinheits-Prüfung
  je Form-Teil hat diesen Treffer nicht ausgesiebt; ihr Ergebnis ist
  noch nirgends getragen (Closure-Notiz §7 ausstehend).
- `verifizierbar`: nein — Text lesen; kein Sensor für Prosa-Orthografie.
- `klasse`: „Form-Fragment in der Vorbild-Kopie"

### F-3 — Writer-Fehlbestands-Pfad schreibt einen kopflosen Träger

- `kategorie`: LOW
- `quelle`: Maintainability (Writer-Form-Grenze, C#-Review F-7-Kette)
- `pfad`: `tools/harness/run-sdk-kotlin-integration-tests.sh:328-345`
- `befund`: Existiert `docs/user/sdk-e2e-abdeckung.md` nicht (oder
  verliert ihre Marker), schreibt `abdeckung_schreiben` ausschließlich
  den nackten Kotlin-Abschnitt — ohne Titel, ohne Erklärkopf, ohne
  Tabellenkopf. Der C#-Vorbild-Writer regeneriert in demselben Fall den
  vollen Kopf neu; die Kotlin-Form (beidseitiger Erhalt, F-7) hat den
  Fehlbestands-Pfad nicht mitübernommen. Degenerater Pfad — die Datei
  ist committet und existiert; der normale, idempotente Pfad ist real
  gemessen (Marker-Abschnitt korrekt angehängt, C#-Abschnitt
  byte-identisch erhalten).
- `verifizierbar`: ja — Ziel-Datei zeitweilig fortbewegen und den
  Writer-Pfad lesen/nachfahren; kein Gate.
- `klasse`: „Writer-Fehlbestands-Pfad ohne Kopf"

### F-4 — Keine In-Test-Frist in den drei Stream-Happy-Paths (Vorbild trägt 90 s)

- `kategorie`: INFO
- `quelle`: Maintainability (Vorbild-Abweichung ohne semantische Wirkung)
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/GrpcRealserverTest.kt:32-40`,
  `SseRealserverTest.kt:27-35`, `NatsRealserverTest.kt:27-35`
- `befund`: Die Kotlin-Happy-Paths warten mit `firstOrNull { … }` ohne
  eigene Frist; das C#-Vorbild trägt eine 90-s-Frist im Test
  (`ReceiveSentinelAsync`, `Assert.True` auf Deadline). Der Rot-Pfad
  verschiebt sich dadurch vollständig in die Runner-Fristen
  (Received-Schleife, Stopped-Warte) — kein falsch-grüner Pfad: ein
  hängender Test erreicht Ausgang 0 nie, der Runner färbt über
  „empfing keine …" bzw. „endete nicht mit Ausgang 0" real rot.
- `verifizierbar`: nein — Verhaltens-Vergleich zum Vorbild; kein Gate.
- `klasse`: „Vorbild-Frist nicht gespiegelt"

### F-5 — READY-Marker hängt an der Testmethoden-Reihenfolge innerhalb der Phase

- `kategorie`: INFO
- `quelle`: Maintainability (Robustheitshinweis, vom C#-Spiegel geerbt)
- `pfad`: `tools/harness/run-sdk-kotlin-integration-tests.sh:173-187`,
  `GrpcRealserverTest.kt` (READY nur im Happy-Path-Test)
- `befund`: Der Runner pollt auf READY, das nur der Happy-Path-Test
  druckt; der Reject-Test derselben Klasse druckt es nicht. JUnit
  Platforms Default-Methoden-Ordnung ist deterministisch, aber nicht
  als Deklarations-Reihenfolge garantiert — läuft der Reject-Test
  zuerst, bleibt READY aus und die Phase färbt rot („kein READY").
  Fehlschlag, nicht Falsch-Grün; dieselbe Struktur wie der geprüfte
  C#-Spiegel.
- `verifizierbar`: ja — ein Lauf mit umgekehrter Methoden-Ordnung
  scheitert sichtbar am READY-Poll; kein Gate.
- `klasse`: „Marker-Abhängigkeit von Testmethoden-Ordnung"

### F-6 — Suchlauf-Feld-Behandlungsspalte zur C#-Zeile unpräzise formuliert

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 („gemeldet statt still mitgeändert")
- `pfad`: `docs/plan/planning/done/slice-sdk-kotlin-reale2e.md:188`
- `befund`: Die Behandlungs-Spalte sagt „gezogen: die Zeile nennt jetzt
  den Kotlin-Abschnitt als real" — die C#-Zeile in `harness/README.md`
  wurde aber nicht geändert (Diff: unverändert); ihre Endklause
  „erweitert sich Slice für Slice um die Kotlin- und Python-HTTP-Abschnitte"
  bleibt als Verlaufs-Aussage wahr und ist zur Hälfte verbraucht
  dokumentiert. Die tatsächliche Behandlung ist „gemeldet, nicht
  gezogen" — der Zustand ist regelkonform, die Spalten-Formulierung
  liest sich als erfolgte Änderung.
- `verifizierbar`: ja — `git diff b58cb173..HEAD -- harness/README.md`
  zeigt die C#-Zeile als Kontext (unverändert).
- `klasse`: „Suchlauf-Feld-Behandlung unpräzise"

---

## Negativbefunde

- geprüft, ohne Befund: `sdks/kotlin/Dockerfile` (additive Stufe
  `integration`: `FROM build` — derselbe Pin `eclipse-temurin:21-jdk@sha256:085e…`
  wiederverwendet, kein neuer Pin; `:?`-Guard im CMD; `integrationTestClasses`
  nur in dieser Stufe, nicht in der `pack-export`-Closure von
  `make sdk-pack-kotlin` — kein stiller Mitlauf, real strukturell
  bestätigt über `tools/harness/sdk-pack-kotlin.sh --target pack-export`)
- geprüft, ohne Befund: `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`
  (eigener SourceSet + Task, **nicht** an `check` gehängt; die zwei
  Konfigurationen tragen die Plan-Nachzug-Begründung als Kommentar
  [Kopplung/Grenze-Form]; `showStandardStreams` trägt den Risiko-§6-Ausgang 3)
- geprüft, ohne Befund: `harness/mk/sdk.mk` (Target real, `.PHONY`,
  direkter `bash`-Aufruf ohne Pipe — §3.9)
- geprüft, ohne Befund: Runner-Disziplin `tools/harness/run-sdk-kotlin-integration-tests.sh`
  (`set -euo pipefail`; `exit 1` je Befund propagiert aus der
  Kommando-Substitution über `set -e`; `trap cleanup EXIT`; Inserts mit
  skript-kontrollierten Werten — keine Injektionsfläche; Phase-Auswahl
  explizit je Testklasse — kein stiller Ausschluss,
  `BEO-PGC/test-runner-stiller-ausschluss`-Disziplin)
- geprüft, ohne Befund: ID-/Sentinel-Kollisionen (Kotlin 440/450/460/470
  gegen C# 400/410/420/430 und Python 300/310/320 — disjunkt; Sentinels
  sprachspezifisch disjunkt), Token-/Publikations-/Slot-/Quellen-Werte
  gegen den `compose.yaml`-Container-Vertrag gehalten (`e2e-reader-token`,
  `e2e-admin-token`, `e2e-nats-stream-token`, `pub_pgc_e2e`,
  `slot_pgc_e2e`, `src-e2e`)
- geprüft, ohne Befund: Testbindung an die Eingabeseite je Phase
  (Token-Form: Ablehnungs-Tests fahren „no-such-token"/„wrong-token" als
  Eingabe, die Plan-Dokumentation trägt die reale Rot-Mutation mit dem
  Wortlaut der Runner-Fehlerzeile; Subjekt-Form über
  `buildSourceSubject` + Tabellen-/Sentinel-Filter; HTTP-Quelle/Publikation
  als Env-Eingabe) — inklusive der gson-JsonNull-Form: die Assertion
  `oldImage == null || isJsonNull` trägt die Aussage (kein Alt-Bild am
  INSERT) über beide gson-Formen und bleibt rot-fähig bei real
  vorhandenem Alt-Bild; sie schärft die C#-Spiegelform (`Assert.Null`),
  die auf dem gson-Baum falsch wäre, und deckt sich semantisch mit der
  Hausform der Unit-Tests (`oldImage?.isJsonNull != false`)
- geprüft, ohne Befund: Inline-Suppression — keine `@Suppress`-Form im
  neuen Kotlin-Baum (§3.2); Import-Grenze `ADR-0109` — keine
  `internal/**`/`cmd/**`/`gen/**`-Imports in `sdks/kotlin/**` (§3.11
  host-lokale Pfade: keine Fundstelle, `make docs-check` real grün:
  930 Dateien, 0 Befunde)
- geprüft, ohne Befund: `docs/user/sdk-e2e-abdeckung.md` (Kotlin-Abschnitt:
  je Zeile `LH-FA-SST-009` zitiert — C#-Review-F-2-Lernklasse angewendet;
  Linktiefe `../../spec/lastenheft.md` korrekt; C#-Abschnitt
  byte-identisch erhalten; Marker-Form dem Runner-Writer identisch)
- geprüft, ohne Befund: `harness/README.md`-Zeile
  `make test-sdk-kotlin-integration` (Target real im Makefile, Runner
  real, Zeile behauptet den Lauf nach dessen Ausführung laut Plan;
  Ursprungs-Formen „real gemessen" den Hauszeilen der Geschwister-Targets
  konform; C#-Zeilen-Endklause bleibt als Verlaufs-Aussage wahr — siehe
  F-6)
- geprüft, ohne Befund: eigener, breiterer §3.13-Suchlauf über beide
  Stände (`grep` nach `Realserver`/`Kotlin-Abschnitt`/`sdk-e2e-abdeckung`
  in `docs/`, `harness/`, `spec/`) — keine weitere Fundstelle, die durch
  diesen Zug falsch geworden wäre; die `done/`-Pläne der Kotlin-Flächen-Slices
  tragen ihre „kein Realserver-Test"-Aussagen sliced-scoped und bleiben
  als Records wahr; `harness/sensors/docs-check.md:145` listet den
  Träger als Coverage-Dimension (im C#-Zug nachgezogen) und wird hier
  nicht überholt; `spec/pflichtenheft.md` `SPEC-026`/`SPEC-027` tragen
  Paketier-/Prüf-Aussagen, keine E2E-Beleg-Aussagen — kein falsch
  werdender Träger
- geprüft, ohne Befund: Traceability/ID-Schema (Commit-Betreff
  `060aa62f` trägt `LH-FA-SST-009` + `ADR-0110`, keine `SPEC-*`/`ARC-*`
  im Betreff; alle neuen Kennungen MR-000-konform), Spec-Stratum
  (Lastenheft/Pflichtenheft unverändert), Docker-only (Runner läuft
  über `make`/Docker, kein lokales Toolchain-Install), §3.4/§3.5/§3.6/§3.8/§3.10
  (Architektur-Sicht, ADRs, Gates, Workflows: unberührt)

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht (Adresse
ohne Artefakt) · Form-Fragment in der Vorbild-Kopie · Writer-Fehlbestands-Pfad
ohne Kopf · Vorbild-Frist nicht gespiegelt · Marker-Abhängigkeit von
Testmethoden-Ordnung · Suchlauf-Feld-Behandlung unpräzise

**Nicht Gegenstand dieses Laufs:** der reale
`make test-sdk-kotlin-integration`-Lauf selbst (EXIT=0, die quotierten
`change_id`s `806-1`/`814-1`/`818-1`) — der Verifier fährt ihn; die
C#-Nachmessung (F-8 des C#-Reviews) zeigt die Klasse reproduzierbar bis
zur nächsten Bring-up-Ketten-Änderung.

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) blockiert; die Fix-Fläche ist
eine Feld-Korrektur im Plan (Träger-Adresse `sdks/kotlin/pgchangefeed-kotlin/README.md`
statt `sdks/kotlin/README.md`); F-2/F-3 können in derselben Fixrunde
mitlaufen oder als akzeptierte Reste dokumentiert werden. F-4/F-5/F-6
sind Hinweise ohne erwartete Aktion.

**Übergabe:** Findings gehen an den Implementer (Rückkante Review →
Implementer); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure
§7 und von dort in den Zähler. Die DoD-Zeile „Review durchgeführt" bleibt
offen — es kommt eine Fixrunde, der Nachzug läuft regulär über
Implementer-Schritt 21. Dieser Report selbst ist ein **Lauf-Beleg**; die
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).