# Verifikations-Report: slice-sdk-sse-filter-phase-verbindung-haertung — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + ADR-Konformität +
Plan-vs-Code-Diff + Gates, in frischem Kontext. Review-Artefakt:
[`review-slice-sdk-sse-filter-phase-verbindung-haertung.md`](review-slice-sdk-sse-filter-phase-verbindung-haertung.md)
(0 HIGH; F-1 MEDIUM außerhalb des Diffs, F-2 LOW, F-3 bis F-5 INFO). Formvorbild:
[`verifikation-slice-sdk-sse-client-schema-table-filter-realserver.md`](verifikation-slice-sdk-sse-client-schema-table-filter-realserver.md).

**Gegenstand:** Slice-Plan
[`slice-sdk-sse-filter-phase-verbindung-haertung`](../plan/planning/done/slice-sdk-sse-filter-phase-verbindung-haertung.md)
(wellenlos), Diff `e900e5c3..HEAD` (`bfea64f0`): Implementer-Commits `26e6a500`, `d0ed33d0`, Review-Commit `bfea64f0`.
Bezug: [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4,
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2,
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md).

Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an Kopien im
Scratchpad (`git archive` je Kopie, `git init`, eigener Image-Tag `pg-change-feed-mutation:vf-<name>` über
`SDK_CSHARP_INTEGRATION_IMAGE`/`SDK_KOTLIN_INTEGRATION_IMAGE`; Mutation per `sed … > Temp` und `cp` in die Kopie, kein `sed -i`
am Repo; die fünf Mutations-Images danach mit `docker rmi` ohne `-f` entfernt, `docker images | grep -c mutation` = 0).
`git status --short` im Echtrepo war nach den Tier-Läufen, nach den Mutationen und nach den Sensoren leer. Keine Verweigerung
der Berechtigungsschicht ([`AGENTS.md`](../../AGENTS.md) §3.15 nicht ausgelöst).

**Eigener Fehler, benannt:** Mein erster Anlauf der Tier-Läufe las veraltete `.ec`-Dateien einer früheren Sitzung im
selben Scratchpad (Pfad ohne frisches Unterverzeichnis) und startete zusätzlich einen zweiten C#-Lauf neben dem laufenden —
die beiden kollidierten im Netz `cdc-feed-test` (Exit 2, „network … not found“). Ich habe beide Läufe beendet und alle
drei Tier-Läufe in einem frischen Verzeichnis sauber seriell wiederholt; die Ergebnisse unten stammen ausschließlich aus
diesem zweiten Anlauf.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | gedruckte Zeile |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.13.0 OK — 54 Dateien`; `d-check: 1561 Datei(en) geprüft, 0 Befund(e)` (docs-check und `commits`); `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `generated-sync: OK`; `sdk-public-doc-check: keine interne Kennung unter sdks`; a-check `gesamt: 0 Befund(e)` |
| `make test` | Exit 0 | kein Go-Code im Diff (Gegenprobe) |
| `make fmt-check` | Exit 0 | `fmt-check: 323 Go-Dateien geprüft, alle formatiert` (deckt nur Go) |
| `make kommentar-kennungen DIFF=e900e5c3` | Exit 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-sse-filter-phase-verbindung-haertung.md` | Exit 0 | `suchlauf-nachmessen: 22 Zeilen stimmen` |
| `make docs-check` | Exit 0 | `d-check: 1561 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make sdk-public-doc-check` | Exit 0 | `keine interne Kennung unter sdks` |
| `make doc-commits RANGE=e900e5c3..HEAD` | Exit 0 | `d-check: 1561 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=e900e5c3..HEAD` | Exit 0 | `d-check: 1561 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-trace` | Exit 0 | `80 Anforderung(en), 0 Waise(n).` (gemessen in diesem Lauf) |

## 2. Die drei Tier-Läufe (selbst gefahren, unmutiert, seriell)

Image: `ghcr.io/pt9912/pg-change-feed:dev` war geladen (Erzeugt 2026-10-01 18:00). `make image` habe ich nicht gefahren:
`git diff 91e46048 HEAD --name-only` nennt keinen Server-Code (nur Tests, Runner, Doku). Dass das geladene Image aus genau
diesem Baum stammt, ist **übernommen**, nicht nachgebaut.

| Tier | Exit | gedruckte Zeile |
|---|---|---|
| `make test-sdk-csharp-integration` | 0 | `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`; SEEN nach 361 ms (Versuch 1), SEEN_SECOND nach 122 ms |
| `make test-sdk-kotlin-integration` | 0 | `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`; SEEN nach 340 ms (Versuch 1), SEEN_SECOND nach 327 ms |
| `make test-sdk-python-integration` | 0 | `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`; SEEN nach 325 ms (Versuch 1), SEEN_SECOND nach 105 ms |

Die Gegenlesung über `cdc.changes` ist Teil des Runners („zweite Gruppe: F1 1, F2 1, ohne Filter 3“ steht in der Report-Zeile
aller drei Läufe). Die Zahlen entsprechen der Erwartung des Plans bei SEEN im Versuch 1. `git status --short` nach den
drei Läufen leer: `docs/user/sdk-e2e-abdeckung.md` wurde nicht erneut geschrieben (Idempotenz). Der bekannte NATS-Phasen-Flake
trat in diesen drei Läufen nicht auf; zusammen mit den vier weiteren C#-Läufen der Mutationsmatrix (§3) liefen **fünf**
C#-Gesamtläufe, jeweils mit durchlaufener NATS-Phase (die Filter-Phase ist die letzte), ohne Ausfall.

## 3. Mutationen (auf Kopien; Farbe mit gedruckter Zeile, je einzeln, Stand HEAD `bfea64f0`)

| Nr. | Sprache | Stelle der Kopie | Farbe | gedruckte Zeile |
|---|---|---|---|---|
| P (Package, **vom Review nicht gefahren**) | C# | `Sse/PgChangeFeedSseClient.cs`: `if (table is not null && table.Length < 0)` (kein `table` auf dem Draht); Dockerfile `dotnet test` → `RUN true` | **rot**, Exit 2, in der Filter-Phase | `FILTER_RESULT f1=4 f1_foreign=2 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` |
| P (Package, **vom Review nicht gefahren**) | Kotlin | `PgChangeFeedSseClient.kt`: `"table" to (null as String?)`; Dockerfile `test` → `RUN true`, `build` → `build -x test` | **rot**, Exit 2, in der Filter-Phase | `FILTER_RESULT f1=4 f1_foreign=2 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` |
| P (Package) | Python | — | — | vom Review gefahren (`f1=4 f1_foreign=2`), von mir **übernommen**, nicht nachgefahren |
| Arm A (ohne Härtung, Parent `91e46048`; Überlagerungen a+b+c) | C# | Package-Mutation wie oben; F1-Konsument startet 4 s später (`Task.Delay(4000)` im Task); in der Fixture `\! sleep 2` nach B und `\! sleep 4` nach dem zweiten Schema in derselben `psql`-Sitzung (autocommit, je INSERT eine eigene Transaktion) | **grün (Falsch-Grün)**, Exit 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15`, SEEN nach 6376 ms (Versuch 1) |
| Arm B (mit Härtung, HEAD; a+b+c) | C# | dieselben Überlagerungen | **rot**, Exit 2 | `FILTER_RESULT f1=3 f1_foreign=1 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` (Runner: „fremde Changes am Filter-Client (f1_foreign=1, f2_foreign=0)“) |
| Kontrolle (HEAD; b+c ohne a) | C# | ohne Package-Mutation | grün, Exit 0 | `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`, SEEN nach 6111 ms, SEEN_SECOND nach 327 ms |

Lesart:

- **Härtungsbeweis erprobt, nicht übernommen:** Arm A grün mit `f1=1`, Arm B rot mit `f1_foreign=1`, Kontrolle grün — wie vom Plan
  erwartet und wie der Implementer es berichtet. Dass die Überlagerung greift (F1 verpasst B der ersten Gruppe), zeigen die
  gemessenen 6376 ms bis SEEN und `f1=1` unter der Package-Mutation: ein verbundener F1 ohne `table` hätte beide `public`-Tabellen
  gesehen. Das ist wie im Review F-2 ein **indirekter** Beleg über die Zählung; gedruckte `RECEIVED_F1`-Zeilen gibt die Fixture im
  grünen Lauf nicht aus. Ich teile Reviews Einschätzung: der Schluss trägt, die Plan-Bedingung „RECEIVED_F1-Zeilen von Arm A“ ist wörtlich
  nicht erfüllt.
- **Abweichung der Überlagerung (c) vom Plan-Wortlaut:** Der Plan nennt „drei getrennte `psql`-Aufrufe“; ich habe `\! sleep` zwischen
  den INSERTs derselben Sitzung benutzt. Wirkung gleich (autocommit, gemessene 6 s bis SEEN), die Form ist eine andere; der Beleg ist die
  gedruckte Zeit, nicht der Wortlaut.
- **Package-Zellen:** C# und Kotlin sind rot **mit** gedrucktem Fremdwert 2 (nicht über „kein SEEN“) — die im Plan verlangte Form.
- **Nicht von mir gefahren:** Kotlin M3 („kein SEEN“) und Python M2 (`f2_foreign > 0`) stehen nur im Implementer-Bericht (Plan §3,
  Beleg-Tabelle) — **übernommen**; das Review fuhr sie ebenfalls nicht. Wie im Vorgänger-Slice gelten sie in der Closure-Notiz als
  *übernommen*, nicht als erprobt, solange sie niemand nachfährt. Dasselbe gilt für Python-Package (Review).

## 4. DoD-Abgleich, Zeile für Zeile (Plan §2; Häkchen setzt die Planner-Closure)

| DoD-Zeile | Befund |
|---|---|
| Liefer-Punkt 1 — Härtung im gemeinsamen Ablauf und C#-Tier | **erfüllt**: `lib-sdk-filter-fixture.sh` committet nach SEEN genau eine zweite Gruppe (B, zweites Schema, A; `id_base + 100`; Sentinel `<Sentinel>Second`, Umgebungswert `PGCHANGEFEED_FILTER_SENTINEL_SECOND`), die Gegenlesung nimmt beide Sentinels (`IN (…)`) und prüft die zweite Gruppe exakt (F1 1, F2 1, U 3, je Tabelle eine). C#-Test wartet auf SEEN_SECOND, druckt es und beginnt dann das Ruhefenster (Quelltext gelesen). Grüner Lauf mit erwarteter Zeile: §2 |
| Liefer-Punkt 2 — Kotlin und Python | **erfüllt**: beide Läufe grün, gedruckte Zeile gleich `f1=2 … unfiltered=6` (§2) |
| Liefer-Punkt 3 — Mutationsproben (i)/(ii)/(iii) | (i) C#, Kotlin **erprobt rot**, Python Review (übernommen); (ii) Arm A/Arm B/Kontrolle **erprobt wie erwartet** (§3); (iii) C#-Package erprobt, Kotlin M3 und Python M2 nur im Implementer-Bericht (**übernommen**). Die übrigen Zellen der Matrix sind laut Plan hergeleitet |
| Nur Test-Code und Runner | **erfüllt**: `git diff --name-only 91e46048 HEAD -- sdks` nennt fünf Pfade (`PhaseEnvironment.cs`, `SseFilterRealserverTests.cs`, `PhaseEnvironment.kt`, `SseFilterRealserverTest.kt`, `test_sse_filter_realserver.py`), alle unter den Integrationstest-Verzeichnissen, keine Versionsdatei, kein Server-Code; kein Tag; `make sdk-public-doc-check` Exit 0 |
| Die Abdeckung ist getragen | **erfüllt**: `git diff e900e5c3 HEAD -- docs/user/sdk-e2e-abdeckung.md` = genau drei Zeilen geändert (je Runner die Filter-Zeile, Kennungen `LH-FA-SST-008`/`-009` unverändert; 3 +/3 −); zweiter Lauf je Tier schreibt nichts (§2); `make docs-check` Exit 0; `make doc-trace` 80/0 |
| `make gates` grün | **erfüllt** (§1, Exit 0 direkt gesichert) |
| Review durchgeführt, kein offenes HIGH/MEDIUM | **Review vorhanden, Zeile nicht ohne Planner-Entscheid abhakbar**: F-1 MEDIUM ist offen. Er liegt außerhalb des Diffs (Bestand der drei Runner); die Hauptlauf-Entscheidung benennt einen eigenen Folge-Slice. Dieser braucht für die Zeile eine Adresse (Slice-Kennung oder Beobachtungs-Eintrag), sonst bleibt „kein offenes MEDIUM“ wörtlich verletzt |
| §3.13-Suchlauf | **erfüllt**: Feld in §3 trägt Gefundenes und Nichtgefundenes je Träger, `suchlauf-nachmessen` Exit 0 (22 Zeilen) |
| Doku-Update | **erfüllt**: `harness/README.md` §Sensors (3 Zeilen, Satz zur zweiten Gruppe), Kopf-Kommentare von Fixture, Testklassen, Runnern; `harness/mk/sdk.mk` unverändert (Suchlauf „dreizehn Phasen“ weiter 4); Handbuch unberührt |
| Closure-Notiz, Beobachtungs-Register, §6-Ausgänge, drei Paarungen | **offen** — Closure-Handlungen des Planners, noch nicht fällig vor dem Verdikt (§6) |

**Entscheidungs-Konformität:** [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4 (SSE nutzt `schema`/`table`) ist am Realserver in zwei
Gruppen belegt (Auswahl rein und vollständig für die zweite Gruppe). [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2: Prüfling
ist das Package, Gegenlesung unabhängig über `cdc.changes`. **Plan-vs-Code:** keine Abweichung des Umfangs gefunden (Fixture, drei Runner, drei
Testdateien, zwei `PhaseEnvironment`-Hilfen, `harness/README.md`, Erzeugnis, Plan); kein Dockerfile-, Compose- oder Workflow-Diff;
kein Produktivcode, keine Version. Die Mehrzeile in §3 „Zeile 8 `Filter-Phase`“ ±0 ist stimmig mit dem Diff.

## 5. Bewertung der Review-Funde

- **F-1 (MEDIUM, `docker logs | grep -qF` unter `pipefail`):** Außerhalb des Diffs, nicht gegen diesen Slice zu wenden. In meinen fünf C#-Gesamtläufen
  trat der Ausfall nicht auf (kein Gegenbeleg gegen die 1-%-Rate der Probe im Review). Die Entscheidung „eigener Folge-Slice“ ist mit §3.9-Nähe
  konsistent; sie braucht eine Adresse (§7).
- **F-2 (LOW):** meine Arm-A-Läufe bestätigen den indirekten Beleg (6376 ms, `f1=1`); die Entscheidung „indirekter Beleg bleibt“ ist tragbar,
  wenn die Closure-Notiz ihn als *indirekt* kennzeichnet.
- **F-3 bis F-5:** zur Kenntnis; F-4 (Routing-Phase) ist Planner-Meldung mit Frist, siehe §7.

## 6. Register-Fortschreibung für die Closure (melden, nicht ändern)

- **NATS-Phasen-Flake (Register-Eintrag fällig):** Unter [`docs/plan/planning/observations/BEO-PGC/`](../plan/planning/observations/BEO-PGC/) gibt es keinen Eintrag zu ihm
  (`ls` der 141 Verzeichnisse, Namen geprüft). Vorkommen mit Anker: (1) **gemessen:** erster C#-Tier-Lauf des Reviewers dieses Slice, Fehlertext
  „der Ablehnungs-Beleg blieb aus (REJECTED token-rejected fehlt)“ ([Review F-1](review-slice-sdk-sse-filter-phase-verbindung-haertung.md); die Rate
  3 von 300 an einer Probe ist im Review gemessen, die Ursache SIGPIPE dort **hergeleitet**). (2) **übernommen:** Lauf des Implementers dieses Slice —
  Review F-1 nennt „meldete dasselbe Symptom“, der Plan §3 trägt dazu keinen auflösbaren Anker. (3) **übernommen/hergeleitet:**
  [Review des Realserver-Slices F-4](review-slice-sdk-sse-client-schema-table-filter-realserver.md) (Implementer-Meldung, Ursache dort als Zeitfenster 20 s gegen 15 s
  hergeleitet — **abweichende Ursache** gegenüber F-1 hier, also nicht ohne Weiteres dasselbe Ereignis). Gemessene Vorkommen mit Anker: **1**; mit den
  übernommenen: **bis zu 3**, davon zwei ohne Lauf-Beleg. Gegenzählung: **0 Ausfälle in 5 C#-Gesamtläufen dieses Verifiers** (und 0 in 2 Kotlin-, 1 Python-Lauf).
  Die Mutation-Läufe sind keine Flake-Vorkommen. Ein Eintrag wäre mit `evidence/`-Dateien zu (1) und den zwei übernommenen zu führen, gekennzeichnet als gemessen bzw. übernommen.
- [`negativtest-ohne-bindung-an-seine-eingabe`](../plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md): der Slice wendet die verkörperte Regel an (Negativ
  erst nach beobachteter Verbindung); der Lerneintrag der Closure-Notiz (Plan §7) ist die Antwort, Zählerstand beim Planner nachzulesen.
- [`docker-cache-ueberspringt-tests-still`](../plan/planning/observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md): **kein Auftreten** — alle gezeigten Mutationen druckten eine geänderte Zeile mit Fremdwert.
- [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md): **kein Auftreten** in diesem Slice (ein Ablauf in der Fixture); eine Vermeidung.
- [`integrationsprojekt-uebersetzt-nicht-unbemerkt`](../plan/planning/observations/BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt/state.md): **kein Auftreten**; alle drei Integrationsprojekte wurden in meinen Läufen gebaut und liefen.
- [`formatierungs-drift-ohne-gate`](../plan/planning/observations/BEO-PGC/formatierungs-drift-ohne-gate/state.md): Review F-5 (Kotlin-Zeilenlängen, INFO, `make fmt-check` deckt Kotlin nicht) ist ein weiteres Nicht-Go-Vorkommen
  derselben Klasse; die Auslegungsfrage aus dem Vorgänger-Verifikationsreport (§6 dort) bleibt beim Planner.
- Plan §6: Risiko „Laufzeit-Verlängerung“ → **entfallen** (SEEN→SEEN_SECOND 105–327 ms gemessen, unter der erwarteten Sekunde); „Überlagerungen treffen die Reihenfolge nicht“ → **entfallen** (Arm A 6376 ms, `f1=1`); „SEEN im Versuch > 1“ → **entfallen** (in allen Läufen dieses Verifiers fiel SEEN im Versuch 1); „Docker-Cache“ → kein Auftreten; „Flake NATS“ → **eingetreten** im Review-Lauf (Kette oben), Ausgang dort zu führen; „Kein Release“ → entfallen.

## 7. Verdikt

**bestanden** — mit den Bedingungen B-1 bis B-3 für die Closure.

- **B-1 (Closure, Pflicht):** Die Closure-Notiz führt die Mutationsmatrix mit Ursprung je Zelle: C#-/Kotlin-Package, Arm A/B/Kontrolle **erprobt (Verifier, 2026-10-02)**;
  Python-Package, Kotlin M3, Python M2 **übernommen** (Review bzw. Implementer-Bericht); alle weiteren Zellen **hergeleitet**. Arm A wird als *indirekter* Beleg
  (Zählung `f1=1`, 6376 ms) gekennzeichnet, nicht als abgelesene `RECEIVED_F1`-Liste.
- **B-2 (Closure, Pflicht):** Review-F-1 (MEDIUM) bekommt eine Adresse (Folge-Slice oder Beobachtungs-Eintrag mit Evidence); erst damit ist die DoD-Zeile „kein offenes
  HIGH/MEDIUM“ ehrlich als „benannt, nicht in diesem Slice“ abhakbar. Die Meldung zur Routing-Phase (Plan §1, Review F-4) nennt bis zur Closure die Adresse oder die Absage des Auftraggebers.
- **B-3 (Closure, Pflicht):** Die Closure-Notiz zählt den NATS-Flake nur mit dem, was ein Anker trägt (gemessen 1, übernommen bis zu 2), und kennzeichnet die abweichenden Ursachen-Herleitungen als nicht gleich.
- Offen für den Planner (keine Bedingung): DoD-Häkchen (hier gesetzt: **keine**), Lerneintrag, Zuordnung von F-5 zu `formatierungs-drift-ohne-gate`.

`make docs-check` nach Anlage dieses Reports: im Commit-Schritt (Exit 0 verlangt).
