# Slice sdk-sse-filter-phase-verbindung-haertung: Härtung der Filter-Phase der drei SDK-Tiers — eine zweite Dreiergruppe nach stehender Verbindung trägt das Negativ

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming,
Haupt-Bezug), [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4
(SSE nutzt die Query-Parameter `schema`/`table`),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2 (Mechanik der SDK-Realserver-Tiers),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen: erprobt gegen hergeleitet),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`). Vorgänger:
[`slice-sdk-sse-client-schema-table-filter-realserver`](../done/slice-sdk-sse-client-schema-table-filter-realserver.md)
(dort §6 „Negativ-Beleg über Abwesenheit", §7 „Folge-Slices") und dessen Verifikation
[`verifikation-slice-sdk-sse-client-schema-table-filter-realserver`](../../../reviews/verifikation-slice-sdk-sse-client-schema-table-filter-realserver.md)
(§5 Review-F-2, §7 Bedingung B-2).

**Berührte Spec-Stellen:** [`SPEC-021`](../../../../spec/pflichtenheft.md) (Drahtvertrag
des SSE-Endpunkts mit `schema`/`table`) — gelesen, nicht geändert.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Folge-Slice zu
[`slice-sdk-sse-client-schema-table-filter-realserver`](../done/slice-sdk-sse-client-schema-table-filter-realserver.md)
(Freigabe des Auftraggebers vom 2026-10-02: „Alles angehen"). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Filter-Phase der drei SDK-Realserver-Tiers (`make test-sdk-csharp-integration`,
`make test-sdk-kotlin-integration`, `make test-sdk-python-integration`) belegt ihr Negativ („der
gefilterte Client sieht die Changes einer anderen Tabelle nicht") erst **nach beobachteter
Verbindung**: nachdem alle drei Clients je eine Change der ersten Dreiergruppe empfangen haben
(`SEEN`, ein Client empfängt nur über eine stehende Verbindung), committet der Runner eine
**zweite Dreiergruppe**; erst wenn die zweite Gruppe an allen drei Clients angekommen ist,
beginnt das Ruhefenster, und der Runner prüft die Zusage in beiden Richtungen (Reinheit und
Vollständigkeit der zweiten Gruppe) gegen `cdc.changes`.

**Befund, aus dem der Slice entsteht** (Verifikation §5/§7 B-2, Closure-Notiz des Vorgängers §6
und §7 „F-2/F-3"): `READY` druckt, sobald die Konsumenten gestartet sind, nicht sobald ihre
Verbindungen stehen; der Runner committet danach sofort die erste Gruppe (B, zweites Schema, A).
Ein gefilterter Client, dessen Verbindung zwischen zwei Commits dieser Gruppe zustande kommt,
bleibt ohne Fremdwert, weil ihn die fremde Change nie erreichte, und kann trotzdem `SEEN`
erfüllen; bei einem Package-Defekt (`table` nicht auf dem Draht) bliebe das unentdeckt. Der
Vorgänger führte das als benannte, **hergeleitete** Grenze (F-2, in acht gesehenen Läufen nie
getroffen). Dazu F-3: `f1 = 1`, `f2 = 1`, `unfiltered = 3` belegen Reinheit der Auswahl, nicht
Vollständigkeit über eine Change hinaus.

**Prüfung des Vorschlags durch den Planner (Code-Lesung am Stand `91e46048`, nicht gefahren):**
Der Vorschlag trägt, mit drei Festlegungen, die der Plan trifft:

1. *Was `SEEN` belegt.* `SEEN` verlangt von F1 die Change der Tabelle A, von F2 die des zweiten
   Schemas, von U alle drei (`lib-sdk-filter-fixture.sh:122–133`, Testschleife in den drei
   Testdateien). Jeder Client hat also vor `SEEN` mindestens eine Change empfangen — seine
   Verbindung steht. Eine danach committete Gruppe erreicht alle drei Streams, und ein Package-
   Defekt zeigt sich an der zweiten Gruppe als Fremdwert (F1 ohne `table`: Change der Tabelle B).
2. *Die zweite Gruppe trägt einen eigenen Sentinel* (`<Sentinel>Second`, Umgebungswert
   `PGCHANGEFEED_FILTER_SENTINEL_SECOND`, vom Fixture gesetzt, in den drei Sprachen gelesen). Ohne
   ihn wäre „die zweite Gruppe ist angekommen" nur über Zählen gegen einen Schnappschuss zu
   erkennen, und ein Nachzügler der ersten Gruppe (Versuch ≥ 2) würde als zweite Gruppe
   gezählt. Mit eigenem Sentinel ist das Ruhefenster-Kriterium exakt: F1 hat `A`·Second, F2 hat
   die Change des zweiten Schemas·Second, U hat alle drei·Second. Das Ruhefenster beginnt erst
   danach (Wert unverändert 15 s). Der Test druckt dann eine Zeile `SEEN_SECOND` (Information
   für die Laufzeit; der Runner wertet `FILTER_RESULT` aus, nicht `SEEN_SECOND`; sein
   `SEEN`-Muster `^\s*SEEN\s*$` trifft sie nicht).
3. *Die Zahlen.* Bei `SEEN` im Versuch 1 (in allen acht gesehenen Läufen des Vorgängers der
   Fall) ergibt sich `f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6`. Bei `SEEN` erst im
   Versuch k > 1 sind die Zahlen der ersten Gruppe(n) größer und nicht festgelegt (ein Client
   kann eine frühere Gruppe verpasst haben). Der Runner prüft deshalb nicht `f1 == 2`, sondern
   die **zweite Gruppe exakt** — gelesen aus `cdc.changes` über den Sentinel der Zeile: F1 genau
   eine, F2 genau eine, U genau drei (je Tabelle eine) — und für die erste Gruppe wie bisher
   „mindestens eine je Auswahl, alle drei bei U". Das schließt die Lücke „Fremdwert 0 belegt
   Reinheit, nicht Vollständigkeit" für die zweite Gruppe deterministisch; die glatte Zahl 2/2/6
   ist der erwartete Befund bei Versuch 1 und steht als solcher in der DoD, nicht als Assertion.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release oder eine Versionsänderung der Packages** — jedes SDK-Release ist eine
  Freigabe des Auftraggebers; der Slice berührt keine Versionsdatei (`.csproj`,
  `pyproject.toml`, `build.gradle.kts`), kein Tag wird gesetzt.
- **SDK-Produktivcode** — geändert werden Test-Code der Tiers und Runner. Findet der Slice
  einen Package-Fehler, geht er an den Vorgänger-Slice
  [`slice-sdk-sse-client-schema-table-filter`](../done/slice-sdk-sse-client-schema-table-filter.md)
  zurück (Rückführung §4), keine stille Reparatur hier.
- **Server-Änderungen** — der Server filtert seit
  [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md); ein gefundener
  Server-Fehler wird gemeldet.
- **Ein neuer Zustellweg** — SSE bleibt die einzige Fläche dieser Phase; gRPC und NATS sind keine
  Fläche des Slice.
- **Die Routing-Phase der Tiers** (`tools/harness/lib-sdk-route-fixture.sh`) — sie ist nach
  Code-Lesung gleich aufgebaut (`READY` abwarten, Dreiergruppe committen, bis `SEEN`,
  `lib-sdk-route-fixture.sh:77–78` und `:138`), ihre Betroffenheit durch dieselbe Lücke ist
  **nicht untersucht** (hergeleitet aus der Gleichheit des Aufbaus). Der Slice schneidet sie
  nicht ein: anderer Gegenstand (Ziel-Auswahl statt Tabellen-Filter), eigene Mutationsmatrix.
  Der Planner meldet den Träger beim Anlegen an den Auftraggeber; ein Folge-Slice entsteht
  nur auf dessen Entscheid (Frist: Closure dieses Slice, §7 nennt dann die Adresse oder die
  Absage).
- **Das Ruhefenster-Maß und die Tier-Struktur** — Wert (15 s) und Form des Fensters bleiben;
  die Phasenzahl der Tiers bleibt dreizehn (die Härtung ist Teil der Filter-Phase, keine
  vierzehnte Phase). Die Hilfetexte in `harness/mk/sdk.mk` bleiben deshalb unverändert
  (Messung in §3, Suchlauf).
- **Eine Workflow- oder Compose-Änderung** — `compose.yaml` bleibt unverändert; kein Workflow
  ruft die drei `make`-Ziele auf (am Start zu messen, §3 Suchlauf),
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift dann nicht.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — Härtung im gemeinsamen Ablauf und C#-Tier.** In
      `tools/harness/lib-sdk-filter-fixture.sh` committet der Runner nach `SEEN` genau eine zweite
      Dreiergruppe (Reihenfolge B, zweites Schema, A; ID-Bereich getrennt von der ersten, z. B.
      Basis + 100; Sentinel `<Sentinel>Second`, vom Fixture als
      `PGCHANGEFEED_FILTER_SENTINEL_SECOND` an den Test gegeben) und wartet danach auf das
      Prozessende. Die Gegenlesung über `cdc.changes` nimmt beide Sentinels an und prüft die zweite
      Gruppe exakt: F1 genau eine Change (`<erstes Schema>.<Tabelle A>`), F2 genau eine (Schema des
      zweiten Schemas), U genau drei (je Tabelle eine); die erste Gruppe wie bisher (mindestens
      eine je Auswahl, bei U alle drei Tabellen). Der C#-Test (`SseFilterRealserverTests`,
      `PhaseEnvironment.cs`) wartet nach `SEEN` bis F1 die Change der Tabelle A·Second, F2 die des
      zweiten Schemas·Second und U alle drei·Second hat (Frist wie die positive Frist, 90 s),
      druckt `SEEN_SECOND` und beginnt **dann** das Ruhefenster; `FILTER_RESULT` behält Form und
      Felder und zählt beide Gruppen, die Fremdwerte werden vor den Assertions gedruckt
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12). *Zu belegen durch:* ein realer, grüner
      `make test-sdk-csharp-integration`-Lauf unmittelbar nach `make image`, im Bericht die
      gedruckte Zeile — erwartet bei `SEEN` im Versuch 1:
      `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`
      (bei Versuch k > 1 größere Zahlen der ersten Gruppe, dann nennt der Bericht den Versuch) —,
      dazu die Zeilen der Gegenlesung und die Schlusszeile des Runners; die dreizehn Phasen laufen
      mit unveränderter Erwartung grün.
- [x] **Liefer-Punkt 2 — Kotlin- und Python-Tier.** Dieselbe Form in
      `SseFilterRealserverTest.kt` und `PhaseEnvironment.kt` (Kotlin) und
      `test_sse_filter_realserver.py` (Python; ein Umgebungswert mehr, gelesen wie die
      bestehenden); die Runner brauchen nur die Aktualisierung ihrer Abdeckungs-Zeile und
      Schlusszeile (der Phasenablauf liegt im Fixture). *Zu belegen durch:* je ein realer, grüner
      `make test-sdk-kotlin-integration`- bzw. `make test-sdk-python-integration`-Lauf nach
      `make image`, im Bericht die gedruckte `FILTER_RESULT`-Zeile je Sprache (erwartet bei
      Versuch 1 dieselben Zahlen `f1=2 … f2=2 … unfiltered=6`).
- [x] **Liefer-Punkt 3 — Mutationsproben, davon der Härtungsbeweis.** Alle Mutationen laufen an
      Kopien im Scratchpad ([`AGENTS.md`](../../../../AGENTS.md) §3.1: Edit/Write, nie `sed -i` oder
      eine Umleitung auf den Arbeitsbaum), jede Probe nennt im Bericht Stelle, Instanz,
      Farbe und die gedruckte Zeile. Ablauf der Läufe seriell (feste Container-Namen der
      Compose-Umgebung); die Kopie eines Commits (`git clone`, anschließend `git checkout <Kennung>`)
      baut ihr Tier-Image unter einem eigenen Tag (`SDK_<SPRACHE>_INTEGRATION_IMAGE=…`,
      gemessen: die drei Runner lesen diesen Wert, `tools/harness/run-sdk-*-integration-tests.sh`),
      so dass das Image des Arbeitsbaums unberührt bleibt; die Unit-Stufe des Package-Baus
      wird an der Kopie so umgangen, wie es der Verifier des Vorgängers vormachte (gemessen
      dort: Python/Kotlin rot in der Filter-Phase, Unit-Stufe an der Kopie umgangen).
      (i) **Package-Mutation je Sprache** (SSE-Client setzt `table` nicht auf den Draht), rot in
      der Filter-Phase mit gedrucktem `f1_foreign > 0`; für C# ist das die bisher nicht erprobte
      Zelle. (ii) **Härtungsbeweis (C#), drei Läufe** — Arm A *ohne Härtung*: Kopie des Parent
      (`91e46048`), Arm B *mit Härtung*: Kopie des Implementer-Commits, dieselben drei
      Überlagerungen in beiden: (a) Package-Mutation wie in (i); (b) im Test der
      F1-Konsument startet erst 4 s nach `READY` (`await Task.Delay(4000)` vor dem Öffnen des
      Streams); (c) im Fixture werden die drei INSERTs der **ersten** Gruppe als drei getrennte
      `psql`-Aufrufe abgesetzt, mit 2 s nach B und 4 s nach dem zweiten Schema (B bei t, zweites
      Schema bei t+2 s, A bei t+6 s; F1 verbindet bei etwa `READY` + 4 s und verpasst damit B und
      die Change des zweiten Schemas, erhält aber A). Erwartet (**hergeleitet**, die Läufe
      messen): Arm A bleibt **grün** mit `f1=1 f1_foreign=0` — das ist das Falsch-Grün, das die
      Härtung schließen soll; Arm B wird **rot** mit `f1_foreign ≥ 1` (die Change der Tabelle B
      aus der zweiten Gruppe erreicht den verbundenen F1). Die Probe gilt nur, wenn die
      `RECEIVED_F1`-Zeilen von Arm A zeigen, dass F1 B der ersten Gruppe tatsächlich verpasst
      hat (sonst ist die Überlagerung wirkungslos und der Lauf kein Beleg). Dazu eine
      **Kontrolle**: Arm B mit den Überlagerungen (b) und (c), aber **ohne** (a) bleibt grün
      mit `f1=2 f1_foreign=0` (die Überlagerung allein färbt nicht rot). (iii) **Die drei bisher
      nicht erprobten Zellen der Matrix** (Closure des Vorgängers §7): C#-Package (= (i) für C#),
      Kotlin M3 (der Test übergibt an F2 das Schema `public` statt des zweiten Schemas: F2 erreicht
      die SEEN-Bedingung nie, der Runner endet mit „kein SEEN" nach der Versuchsgrenze), Python M2
      (der Test übergibt `schema` nicht an F2: `f2_foreign > 0`). Die übrigen Zellen
      (M1 je Sprache, M2 C#/Kotlin, M3 C#/Python, Package Kotlin/Python am Vorgänger-Stand
      gemessen) gelten für den neuen Stand als **hergeleitet** — die Härtung ändert an der
      Eingabeseite nichts —, und die Closure-Notiz führt die Matrix mit dem Ursprung je Zelle.
- [x] **Nur Test-Code und Runner.** `git diff --name-only 91e46048 -- sdks` nennt ausschließlich
      Pfade unter den drei Test-Verzeichnissen (`PgChangeFeed.Client.Integration/`,
      `src/integrationTest/`, `integration/`), keine Versionsdatei; kein Release, kein Tag;
      `make sdk-public-doc-check` endet mit Exit 0. *Zu belegen durch:* der Diff-Befehl im
      Bericht und der Suchlauf in §3.
- [x] **Die Abdeckung ist getragen.** Die Zeile je SDK im marker-gegrenzten Abschnitt von
      [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) (geschrieben von den drei
      Runnern, nicht von Hand) nennt die zweite Gruppe nach stehender Verbindung; Kennungen
      ([`LH-FA-SST-008`](../../../../spec/lastenheft.md),
      [`LH-FA-SST-009`](../../../../spec/lastenheft.md)) und Nachweis-Spalte unverändert.
      *Zu belegen durch:* `git diff` der Datei nach den drei Läufen zeigt genau die drei geänderten
      Zeilen, ein zweiter Lauf je Tier schreibt nichts, fremde Abschnitte bleiben; `make docs-check`
      Exit 0. Die Läufe der Tiers sind **seriell** und enden vor dem Commit der Datei (andere
      Runner und `make bench` schreiben Erzeugnisse im selben Arbeitsbaum; der Implementer
      committet ausschließlich eigene Pfade).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-sdk-sse-filter-phase-verbindung-haertung.md` liegt vor, kein
      offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      *Lesart:* kein HIGH; das eine MEDIUM (F-1, `grep -q` hinter `docker logs` unter `pipefail` in
      den Runnern) liegt **außerhalb des Diffs** (Bestand der Runner, nicht Gegenstand dieses Slice)
      und trägt die Adresse
      [`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../in-progress/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
      sowie `BEO-PGC/runner-grep-pipe-verfehlt-zeile` — „benannt, nicht in diesem Slice“, wie die
      Verifikation §7 B-2 es verlangt; die Review-Reports bleiben unverändert.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes je
      Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-sse-filter-phase-verbindung-haertung.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen `make test-sdk-*-integration`:
      Satz zur Filter-Phase um die zweite Gruppe nach beobachteter Verbindung ergänzt), die
      Kopf-Kommentare von `lib-sdk-filter-fixture.sh`, der drei Testklassen und der Runner
      tragen den Ist-Umfang; das Benutzerhandbuch bleibt unberührt (es beschreibt die Packages,
      nicht die Tiers).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7) — neues Verzeichnis oder
      weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Wellen-Betrieb für diesen wellenlosen Slice hier geprüft.

**Umfang:** M — Schätzung, nicht gemessen: eine Hilfsdatei, drei Testdateien samt Umgebungs-Hilfen,
drei Runner-Texte; der Anteil an Läufen überwiegt (drei Läufe unmutiert, drei Package-Läufe, zwei
Zellen-Läufe, drei Läufe des Härtungsbeweises: elf Tier-Läufe zu je mehreren Minuten).

**Voraussetzung:** `make image` ist ausgeführt (`compose.yaml` referenziert das lokal gebaute
Image, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); die Tier-Läufe setzen es voraus.
Kein anderer Runner und kein `make bench` läuft gleichzeitig (feste Container-Namen).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/lib-sdk-filter-fixture.sh` | update | Nach `SEEN` die zweite Dreiergruppe (Sentinel `${sentinel}Second`, ID-Bereich getrennt); Umgebungswert `PGCHANGEFEED_FILTER_SENTINEL_SECOND` an den Test-Container; Gegenlesung nimmt beide Sentinels (`IN (…)`, derzeit `new_data->>'name' = '$sentinel'` an zwei Stellen) und prüft die zweite Gruppe exakt (F1 1, F2 1, U 3, je Tabelle eine); Kopf-Kommentar, Fehlertexte und `SDK_FILTER_REPORT` nennen die zweite Gruppe. Ein Ablauf für die drei Sprachen (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). |
| `sdks/csharp/PgChangeFeed.Client.Integration/SseFilterRealserverTests.cs`, `PhaseEnvironment.cs` | update | Warten auf die zweite Gruppe nach `SEEN`, `SEEN_SECOND`, Ruhefenster danach; `Own`/`HasAllThree` je Sentinel; neuer Umgebungswert `FilterSentinelSecond`. |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/SseFilterRealserverTest.kt`, `PhaseEnvironment.kt` | update | dieselbe Form. |
| `sdks/python/pgchangefeed/integration/test_sse_filter_realserver.py` | update | dieselbe Form (`_SENTINEL_SECOND` aus der Umgebung). |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh` | update | Abdeckungs-Zeile im marker-gegrenzten Abschnitt und Schlusszeile („Filter-Belege") nennen die zweite Gruppe; Kopf-Kommentar. Aufrufe von `sdk_filter_phase` bleiben (ID-Basis, Sentinel je Sprache). |
| `docs/user/sdk-e2e-abdeckung.md` | Erzeugnis | die Runner schreiben ihre Abschnitte; nicht von Hand. |
| `harness/README.md` §Sensors (drei Zeilen), `harness/mk/sdk.mk` | update / prüfen | README: der Satz zur Filter-Phase; `sdk.mk`: Phasenzahl und Aufzählung unverändert (nichts nachzuziehen, solange Suchlauf-Zeile 9 unverändert 4 ist). |
| `sdks/*/Dockerfile` (Stufe `integration`), `compose.yaml`, Workflows | prüfen | keine Änderung erwartet; eine Änderung ist ein Plan-Nachzug. |

**Ansatz — Reihenfolge und Ruhefenster.** Der Test kennt zwei Marken: `SEEN` (erste Gruppe,
Verbindung belegt) und `SEEN_SECOND` (zweite Gruppe an allen drei Clients angekommen). Der Runner
reagiert auf `SEEN` mit dem Commit der zweiten Gruppe (Polling 0,2 s; Latenz von `SEEN` bis
Commit unter einer Sekunde, **hergeleitet** aus dem Polling-Intervall, im Lauf zu messen). Das
Ruhefenster beginnt nach `SEEN_SECOND`: derselbe Zustellweg hat dann beide Gruppen an U und die
gefilterten Treffer an F1/F2 ausgeliefert, die Abwesenheit der fremden Changes ist kein früher
Abbruch. Der Test wartet auf die Sentinel-Zeilen, nicht auf eine Zahl: ein Nachzügler der ersten
Gruppe verschiebt das Kriterium nicht.

**Ansatz — Härtungsbeweis (ii).** Beide Arme laufen an Kopien in einem Scratchpad-Verzeichnis
(Arm A: Parent `91e46048`; Arm B: der Implementer-Commit). Die Überlagerungen (a) bis (c) sind
Edit/Write an der Kopie; die Rücknahme ist das Verwerfen der Kopie. Kein Lauf mutiert den
Arbeitsbaum, kein `make image` mit mutiertem Baum ([`AGENTS.md`](../../../../AGENTS.md) §3.1,
[`harness/targets/image-mutation.md`](../../../../harness/targets/image-mutation.md)); das
Server-Image `:dev` bleibt unmutiert, nur das Tier-Image der Kopie trägt einen eigenen Tag. Die
Docker-Cache-Falle (`BEO-PGC/docker-cache-ueberspringt-tests-still`) gilt wie im Vorgänger: der
Beleg ist die gedruckte `FILTER_RESULT`-Zeile mit dem gezählten Fremdwert, nicht der Exit-Code des
Baus.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Filter-Phase committet eine
Dreiergruppe und endet bei `SEEN`", die Zahlen 1 und 3 in `FILTER_RESULT` und den Beschreibungen
(„alle drei", „drei Gruppen"), das Zählwort „dreizehn Phasen" (soll stehen bleiben); Parent ist
`91e46048`, gemessen am 2026-10-02 mit `make suchlauf-nachmessen`; die `diff`-Zeilen und die
Befunde trägt der Implementer nach; neue Dateien sind für den Stand `diff` mit `git add` im
Index):**

```suchlauf
91e46048 12 -n -F FILTER_RESULT -- harness tools sdks docs/user
91e46048 9 -n -F Dreiergruppe -- harness tools sdks docs/user
91e46048 20 -n -F 'drei Gruppen' -- harness tools sdks docs/user
91e46048 13 -n -F 'ohne Filter alle drei' -- harness tools sdks docs/user
91e46048 14 -n -F 'unfiltered=' -- harness tools sdks docs/user
91e46048 16 -n -F f1_foreign -- harness tools sdks docs/user
91e46048 2 -n -F "new_data->>'name'" -- tools/harness/lib-sdk-filter-fixture.sh
91e46048 7 -n -F Filter-Phase -- harness tools sdks docs/user .github
91e46048 4 -n -E 'dreizehn Phasen' -- harness tools
91e46048 0 -n -F SEEN_SECOND -- harness tools sdks docs/user
91e46048 0 -n -E 'SENTINEL_SECOND|FilterSentinelSecond' -- harness tools sdks docs/user
diff 13 -n -F FILTER_RESULT -- harness tools sdks docs/user
diff 30 -n -F Dreiergruppe -- harness tools sdks docs/user
diff 11 -n -F 'drei Gruppen' -- harness tools sdks docs/user
diff 12 -n -F 'ohne Filter alle drei' -- harness tools sdks docs/user
diff 14 -n -F 'unfiltered=' -- harness tools sdks docs/user
diff 16 -n -F f1_foreign -- harness tools sdks docs/user
diff 2 -n -F "new_data->>'name'" -- tools/harness/lib-sdk-filter-fixture.sh
diff 7 -n -F Filter-Phase -- harness tools sdks docs/user .github
diff 4 -n -E 'dreizehn Phasen' -- harness tools
diff 13 -n -F SEEN_SECOND -- harness tools sdks docs/user
diff 12 -n -E 'SENTINEL_SECOND|FilterSentinelSecond' -- harness tools sdks docs/user
```

| Träger | Messung am Parent (`91e46048`, 2026-10-02) | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Zeile 1 `FILTER_RESULT` (Abschlusszeile, Parser, Beschreibungen) | 12 Zeilen | Fixture (Parser/Regex), drei Testdateien, drei Abdeckungs-Texte, `harness/README.md`; die Form der Zeile bleibt, der Wert ändert sich — jede Trefferzeile gelesen, ob sie eine Zahl nennt. **Befund (Diff, 13 Zeilen, +1):** die Mehrzeile ist die Test-Kommentar-/Fixture-Zeile, die `FILTER_RESULT` beschreibt; keine der 13 Zeilen nennt einen festen Zahlenwert (Zahlen stehen nur als gedruckte Laufzeile im Bericht), nichts nachzuziehen. |
| Zeile 2 `Dreiergruppe` | 9 Zeilen — `lib-sdk-filter-fixture.sh` (3), `lib-sdk-route-fixture.sh` (3), `run-integration-tests.sh` (3) | nur die drei Treffer der Filter-Datei sind Gegenstand; die sechs fremden (Routing-Phase der Tiers, Server-Rundlauf) bleiben, **gemeldet** (§1). **Befund (Diff, 30 Zeilen):** 27 Zeilen in geänderten Dateien — Fixture 6 (+3), drei Runner je 4 (Kopf, Phasenkommentar, Abdeckungs-Text, Schlusszeile), `harness/README.md` 3 — und 3 im Erzeugnis `docs/user/sdk-e2e-abdeckung.md` (von den Runnern geschrieben); die sechs fremden (`lib-sdk-route-fixture.sh` 3, `run-integration-tests.sh` 3) unverändert 6. **Gemeldet, nicht geändert:** die Routing-Phase hat denselben Aufbau (Dreiergruppe nach `READY`, `SEEN`); Betroffenheit nicht untersucht (hergeleitet), Entscheid beim Auftraggeber. |
| Zeile 3 `drei Gruppen` | 20 Zeilen: Filter-Dateien (Fixture ×3, drei Testdateien je ×2) 9, davon fremd: Routing-Szenarien der drei Tiers (je ×3) 9, `docs/user/e2e-abdeckung.md` ×1 und `run-integration-tests.sh` ×1 | je Filter-Zeile lesen: „alle drei" meint die drei **Tabellen** der Gruppe und bleibt wahr; nachzuziehen nur, wo „genau eine Gruppe" gemeint ist; die elf fremden Treffer gehören nicht zum Gegenstand. **Befund (Diff, 11 Zeilen, −9):** gelesen — die neun Filter-Treffer des Parents sind auf null gesunken: die Fixture-Treffer (3) und die Beschreibungen sind auf die Zwei-Gruppen-Form umgeschrieben, in den drei Testdateien steht statt „alle drei Gruppen" nun „alle drei Tabellen" (die Meldung der ersten Frist meint die Tabellen der Dreiergruppe); verbleibend sind die elf fremden (Routing: `RouteScenario`/`route_scenario` ×9, `docs/user/e2e-abdeckung.md` 1, `run-integration-tests.sh` 1), unverändert. |
| Zeile 4 `ohne Filter alle drei` | 13 Zeilen (`sdk-e2e-abdeckung.md` ×3, `harness/README.md` ×3, Fixture ×1, drei Runner je ×2) | wie Zeile 3; die Abdeckungs-Texte der drei Runner nennen die zweite Gruppe. **Befund (Diff, 12 Zeilen, −1):** die Fixture-Zeile (Report-Text „alle drei Gruppen") trägt jetzt die Zahlen der zweiten Gruppe; die Abdeckungs-Texte der Runner (je 2) und `harness/README.md` (3) und der Erzeugnis-Träger (3) tragen „ohne Filter alle drei" weiter wahr (gemeint: alle drei Tabellen), die zweite Gruppe steht daneben. |
| Zeile 5 `unfiltered=` | 14 Zeilen | Parser und Beschreibungen; erwarteter Wert bei Versuch 1: 6. **Befund (Diff, 14 Zeilen, ±0):** Parser/Beschreibungen gelesen, keine nennt einen festen Wert; gemessen 6 in allen drei Läufen. |
| Zeile 6 `f1_foreign` | 16 Zeilen | Parser, Fehlertexte, Beschreibungen; Wert unverändert 0. **Befund (Diff, 16 Zeilen, ±0):** unverändert, gemessen 0 in allen drei unmutierten Läufen. |
| Zeile 7 Gegenlesung `new_data->>'name'` | 2 Zeilen in der Fixture-Datei | beide Stellen nehmen beide Sentinels. **Befund (Diff, 2 Zeilen):** beide Zeilen lesen `IN ('<Sentinel>', '<Sentinel>Second')` und liefern zusätzlich den Namen der Zeile (Zuordnung zur Gruppe). |
| Zeile 8 `Filter-Phase` | 7 Zeilen (`harness/README.md` ×3, Fixture ×1, drei Runner je ×1) | Texte, die die Phase beschreiben: README-Satz nachziehen. **Befund (Diff, 7 Zeilen, ±0):** die drei README-Zeilen tragen den Satz zur zweiten Gruppe, Fixture und die drei Runner-Phasenkommentare bleiben wahr; gefunden und nachgezogen, nichts offen. |
| Zeile 9 `dreizehn Phasen` | 4 Zeilen (drei Hilfetexte `harness/mk/sdk.mk`, ein Runner-Kopf) | **soll unverändert 4 bleiben** (keine vierzehnte Phase). **Befund (Diff):** 4, unverändert. |
| Zeilen 10 und 11 (neue Namen) | 0 Zeilen am Parent | Soll am Diff: `SEEN_SECOND` in Fixture-Kommentar, drei Testdateien; `SENTINEL_SECOND`/`FilterSentinelSecond` in Fixture, Kotlin-/C#-Umgebungs-Hilfe, Python-Test — der Implementer trägt die gezählten Werte ein. **Befund (Diff):** `SEEN_SECOND` 13 Zeilen: Fixture 4, `harness/README.md` 3, drei Testdateien je 2. `SENTINEL_SECOND|FilterSentinelSecond` 12 Zeilen: Fixture 2, C# `PhaseEnvironment.cs` 1 und Test 2, Kotlin `PhaseEnvironment.kt` 1 (die Kotlin-Testdatei liest `filterSentinelSecond` in Kleinschreibung und trifft das Muster nicht — Musterlücke, nicht Trägerlücke), Python 6. |

**Belege des Implementers (gemessen am 2026-10-02, Commit `26e6a500` plus Wortlaut-Korrektur der Frist-Meldung, jeder Lauf seriell nach `make image`; die DoD-Häkchen setzt die Planner-Closure):**

Unmutierte Läufe (zweiter Lauf je Tier schreibt `docs/user/sdk-e2e-abdeckung.md` nicht erneut; der erste Lauf änderte genau drei Zeilen):

```text
C#:     FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15   (SEEN nach 333ms, SEEN_SECOND nach 333ms, Versuch 1)
Kotlin: FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15   (SEEN nach 346ms, SEEN_SECOND nach 331ms, Versuch 1)
Python: FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15   (SEEN nach 350ms, SEEN_SECOND nach 379ms, Versuch 1)
```

Mutationen an Kopien im Scratchpad (Edit/Write, kein `sed -i`; Unit-Stufe des Package-Baus an der Kopie umgangen; eigener Image-Tag `pg-change-feed-mutation:hd-*`, danach `docker rmi`):

| Zelle | Stelle der Kopie | Farbe | gedruckte Zeile |
|---|---|---|---|
| Package C# | `PgChangeFeedSseClient.cs`: `table` nicht auf den Draht, `dotnet test` → `RUN true` | rot, Exit 2, Filter-Phase | `FILTER_RESULT f1=4 f1_foreign=2 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` |
| Package Kotlin | `PgChangeFeedSseClient.kt`: `"table" to (null as String?)`, Tests umgangen | rot, Exit 2 | `f1=4 f1_foreign=2 f2=2 f2_foreign=0 unfiltered=6` |
| Package Python | `sse_client.py`: `("table", None)`, `pytest` umgangen | rot, Exit 2 | `f1=4 f1_foreign=2 f2=2 f2_foreign=0 unfiltered=6` |
| Kotlin M3 | `SseFilterRealserverTest.kt`: F2 mit `filterSchemaA` | rot, Exit 2 | „kein SEEN" nach fünf Dreiergruppen |
| Python M2 | `test_sse_filter_realserver.py`: F2 ohne `schema` | rot, Exit 2 | `f1=2 f1_foreign=0 f2=6 f2_foreign=4 unfiltered=6` |
| Arm A (Parent `91e46048`, ohne Härtung; Überlagerungen a+b+c) | C#-Package-Mutation, F1 startet 4 s nach `READY`, erste Gruppe als drei `psql`-Aufrufe (B bei t, zweites Schema t+2 s, A t+6 s) | **grün (Falsch-Grün)**, Exit 0 | `FILTER_RESULT f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3 quiet_seconds=15`, SEEN nach 6454ms; f1=1 bei Package-Mutation und Gegenlesung `public.feed_e2e_sdkfilt_a`: F1 hat B der ersten Gruppe verpasst |
| Arm B (Stand mit Härtung; a+b+c) | dieselben Überlagerungen | rot, Exit 2 | `FILTER_RESULT f1=3 f1_foreign=1 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` |
| Kontrolle (Stand mit Härtung; b+c, ohne a) | ohne Package-Mutation | grün, Exit 0 | `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`, SEEN nach 6201ms, SEEN_SECOND nach 324ms |

## 4. Trigger

**Start** (`next` → `in-progress`): `make image` ist ausgeführt, Docker mit Netzzugang für die
Tier-Bauten, kein anderer Runner und kein `make bench` laufen (die Hintergrundläufe der aktuellen
Sitzung sind beendet), und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Der Vorgänger
[`slice-sdk-sse-client-schema-table-filter-realserver`](../done/slice-sdk-sse-client-schema-table-filter-realserver.md)
liegt in `done/`; Versionen der Packages unverändert.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Härtungsbeweis (ii) trennt sich als
  eigener Slice ab (`…-haertung-beweis`); Fixture-Härtung und die drei Tiers bleiben in diesem
  Slice.
- `in-progress` → `open` (blockiert): ein Package oder der Server verhält sich mit der zweiten
  Gruppe anders als `ADR-0133` (etwa: F1 bekommt trotz stehender Verbindung die zweite Gruppe
  nicht) — der Fund geht an den Vorgänger-Slice bzw. den Architect, kein Umgehen im Test; ein Rot
  geht nie als Anpassung der Erwartung in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
je ein realer, grüner Lauf der drei `make test-sdk-*-integration`-Ziele nach `make image` mit den
gedruckten `FILTER_RESULT`-Zeilen, Mutationsproben (i)/(iii) rot gesehen, Härtungsbeweis (ii)
(Arm A grün, Arm B rot, Kontrolle grün — oder jede Abweichung benannt), Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Laufzeit-Verlängerung durch die zweite Gruppe.** Erwartet: unter zwei Sekunden je Tier
  (Polling 0,2 s plus Zustellung ~0,3 s laut Vorgänger-Lauf „SEEN nach 335–369 ms"); das
  Ruhefenster bleibt 15 s. **Hergeleitet**, nicht gefahren. — **Ausgang:** *entfallen* —
  `SEEN`→`SEEN_SECOND` 105 bis 327 ms (Verifier, gemessen, Verifikation §6; Implementer-Läufe
  331 bis 379 ms, übernommen), unter der erwarteten Sekunde.
- **Flake der bestehenden NATS-Phase.** Der Vorgänger sah einmal einen Ausfall der **alten**
  NATS-Phase (F-4: Fenster des Runners 20 s gegen 15 s des Tests, hergeleitet). Die Härtung
  verlängert diese Phase nicht; sie wird hier nicht untersucht und nur als Hinweis geführt. —
  **Ausgang:** *eingetreten* — im ersten C#-Tier-Lauf des Reviewers (gemessen, Review F-1;
  die Ursache ist dort `docker logs … | grep -qF` unter `pipefail`, hergeleitet, nicht das
  Zeitfenster aus F-4 des Vorgängers). Adresse:
  [`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../in-progress/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
  und `BEO-PGC/runner-grep-pipe-verfehlt-zeile` (1×). Gegenzählung des Verifiers: 0 Ausfälle in
  fünf C#-, zwei Kotlin- und einem Python-Gesamtlauf (gemessen).
- **Die Überlagerungen des Härtungsbeweises treffen die Reihenfolge nicht.** Die Rechnung
  (B bei t, zweites Schema t+2 s, A t+6 s, F1-Verbindung ≈ `READY` + 4 s, Commit der ersten Gruppe
  ≤ 1 s nach `READY`) lässt je ≥ 1 s Spielraum auf beiden Seiten; Beleg ist die
  `RECEIVED_F1`-Liste von Arm A (nur A der ersten Gruppe). — **Ausgang:** *entfallen* —
  Arm A grün mit `f1=1`, SEEN nach 6454 ms (Implementer) bzw. 6376 ms (Verifier, gemessen):
  F1 verpasste B der ersten Gruppe. Der Beleg ist **indirekt** über die Zählung, nicht über eine
  abgelesene `RECEIVED_F1`-Liste (Review F-2; die Fixture druckt sie im grünen Lauf nicht); die
  Plan-Bedingung ist wörtlich nicht erfüllt, der Schluss trägt (Entscheid des Auftraggebers
  vom 2026-10-02: kein zusätzliches Drucken).
- **`SEEN` im Versuch > 1.** Die Zahlen der ersten Gruppe sind dann größer; die exakte Prüfung
  trägt die zweite Gruppe (Sentinel), nicht `f1 == 2`. — **Ausgang:** *entfallen* — der Runner
  zählt die zweite Gruppe über den Sentinel (Konstruktion); `SEEN` fiel in allen Läufen im
  Versuch 1 (Implementer, Review, Verifier).
- **Rest-Grenze der ersten Gruppe.** Für die erste Gruppe bleibt F-2 sinngemäß: ihr Negativ ist
  nicht belegt. Das ist gewollt — das Negativ trägt die zweite Gruppe; die erste trägt nur die
  Verbindung (`SEEN`). Ein Empfang belegt eine stehende Verbindung (**hergeleitet**: ein Client
  empfängt nur über seine Verbindung). — **Ausgang:** *entfallen* — Arm A grün, Arm B rot,
  Kontrolle grün (Verifier, gemessen).
- **Docker-Cache der Stufe `integration`** (`BEO-PGC/docker-cache-ueberspringt-tests-still`). —
  **Ausgang:** *entfallen* (kein Auftreten) — jede gezeigte Mutation druckte eine geänderte
  `FILTER_RESULT`-Zeile mit gezähltem Fremdwert (Verifier, gemessen).
- **Drei Sprachen, ein Randfall** (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). —
  **Ausgang:** *entfallen* (kein Auftreten) — Ablauf und Gegenlesung stehen einmal in der
  Fixture-Datei, derselbe Eingabesatz je Sprache.
- **Kein Sensor übersetzt die Integrationsprojekte**
  (`BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`). — **Ausgang:** *weiter offen* →
  Register (kein Auftreten in diesem Slice, keine `evidence/`; der Slice fuhr alle drei Tiers
  real, die Integrationsprojekte übersetzten und liefen).
- **Gleichzeitige Läufe im Arbeitsbaum.** Der Auftraggeber fährt SDK-Runner und `make bench` im
  Hintergrund; diese schreiben Erzeugnisse (z. B. `docs/user/sdk-e2e-abdeckung.md`,
  `docs/user/bench-abdeckung.md`) und belegen feste Container-Namen. — **Ausgang:**
  *eingetreten*, ohne Auswirkung auf die Belege: der Verifier startete
  im ersten Anlauf einen zweiten C#-Lauf neben dem laufenden und las veraltete `.ec`-Dateien
  einer früheren Sitzung (eigener Fehler, im Verifikations-Report benannt); beide Läufe
  kollidierten im Netz `cdc-feed-test`, die Ergebnisse stammen aus dem seriellen zweiten Anlauf.
  Eine Kollision mit den Hintergrundläufen des Auftraggebers ist in keinem Bericht genannt.
  Gegenmaßnahme trug: serielle Läufe, Commit nur eigener Pfade (`git status --short` im
  Echtrepo nach den Läufen leer, Verifier gemessen).
- **Kein Release, keine Versionsänderung.** — **Ausgang:** *entfallen* (Verifier, gemessen:
  `git diff --name-only 91e46048 HEAD -- sdks` nennt fünf Test-Pfade, keine Versionsdatei, kein
  Tag).

## 7. Closure-Notiz

Ursprung der Angaben ([`AGENTS.md`](../../../../AGENTS.md) §3.12): **gemessen** = vom Reviewer
([`review-slice-sdk-sse-filter-phase-verbindung-haertung`](../../../reviews/review-slice-sdk-sse-filter-phase-verbindung-haertung.md))
oder Verifier
([`verifikation-slice-sdk-sse-filter-phase-verbindung-haertung`](../../../reviews/verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md))
im eigenen Lauf; **übernommen** = aus dem Bericht des Implementers (§3 „Belege des Implementers“)
oder einem der beiden Berichte ohne Nachmessung des Planners; **hergeleitet** = nicht gefahren.
Der Planner hat keine Zahl dieser Notiz nachgemessen, außer dem Suchlauf-Feld und den
Register-Zählern (Dateien gezählt).

- **Was hat funktioniert:** Alle drei Tiers druckten unmutiert, je Lauf vom Verifier gefahren
  (gemessen), `FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`,
  jeweils `SEEN` im Versuch 1 (`SEEN` nach 361 ms C#, 340 ms Kotlin, 325 ms Python;
  `SEEN`→`SEEN_SECOND` 122 ms, 327 ms, 105 ms). Der Reviewer sah dieselbe Zeile in seinen Läufen
  (gemessen), die Implementer-Zeilen decken sich (übernommen; `SEEN_SECOND` nach 333, 331, 379 ms).
  Die Gegenlesung über `cdc.changes` (zweite Gruppe exakt: F1 1, F2 1, ohne Filter 3) bestand in
  allen Läufen; `docs/user/sdk-e2e-abdeckung.md` änderte genau drei Zeilen, ein zweiter Lauf je
  Tier schrieb nichts (gemessen, Verifier). Mutationsmatrix mit Ursprung je Zelle (gedruckte
  Zeilen in den beiden Berichten und in §3 „Belege des Implementers“):

  | Zelle | C# | Kotlin | Python |
  |---|---|---|---|
  | Package (`table` nicht auf den Draht) | rot in der Filter-Phase, `f1=4 f1_foreign=2` (Verifier, **erprobt**, Unit-Stufe an der Kopie umgangen) | rot, `f1=4 f1_foreign=2` (Verifier, **erprobt**) | rot, `f1=4 f1_foreign=2` (Review, **erprobt**) |
  | M3 (falsches Schema an F2) | **hergeleitet** | rot, „kein SEEN“ (Implementer, **übernommen**) | **hergeleitet** |
  | M2 (`schema` weg an F2) | **hergeleitet** | **hergeleitet** | rot, `f2=6 f2_foreign=4` (Implementer, **übernommen**) |
  | M1 (`table` weg) | **hergeleitet** | **hergeleitet** | **hergeleitet** |
  | Zeilenzahl der zweiten Gruppe falsch erwartet (Fixture, einmal) | — | — | rot, „F1 empfing 1 Change(s) der zweiten Gruppe“ (Review, **erprobt**; die Fixture ist für alle drei Sprachen dieselbe) |

  Die sieben hergeleiteten Zellen stützen sich darauf, dass die Eingabeseite der Mutation sich
  gegenüber dem Vorgänger nicht ändert (dort erprobt) und der Ablauf in der Fixture einmal steht;
  eine Messung am neuen Stand steht aus. Härtungsbeweis (C#, SSE-Package-Mutation plus F1
  verzögert plus erste Gruppe mit Abstand): Arm A (Parent `91e46048`, ohne Härtung) **grün**,
  `f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3`, `SEEN` nach 6376 ms (Verifier, **erprobt**;
  Implementer 6454 ms, übernommen) — das Falsch-Grün; Arm B (Härtung) **rot**,
  `f1=3 f1_foreign=1 f2=2 f2_foreign=0 unfiltered=6` (Verifier, **erprobt**); Kontrolle (Härtung,
  ohne Package-Mutation) **grün**, `f1=2 … unfiltered=6`, `SEEN` nach 6111 ms, `SEEN_SECOND` nach
  327 ms (Verifier, **erprobt**). **Arm A ist ein indirekter Beleg** (Review F-2): die Zählung
  `f1=1` und die gemessene Zeit bis `SEEN` zeigen, dass F1 die ersten Changes verpasste; eine
  abgelesene `RECEIVED_F1`-Liste druckt die Fixture im grünen Lauf nicht, und der Auftraggeber
  hat ein zusätzliches Drucken nicht gewollt. Die Überlagerung (c) fuhr der Verifier als `\! sleep`
  zwischen den INSERTs derselben `psql`-Sitzung statt als drei getrennte `psql`-Aufrufe (gleiche
  Wirkung: autocommit, gemessene 6 s bis `SEEN`).
- **Was ging anders als geplant:** (1) Review F-1 (MEDIUM): ein Ausfall der NATS-Phase im ersten
  C#-Lauf des Reviewers, Ursache `docker logs … | grep -qF` unter `pipefail` (hergeleitet,
  Reviewer-Messung 3 von 300 an einem Nachbau, vom Reviewer gemessen); liegt außerhalb des Diffs
  und in allen drei Runnern und in `run-integration-tests.sh` (35 Fundstellen, am Stand
  `a420e223` gemessen) — **Adresse:**
  [`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../in-progress/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
  und `BEO-PGC/runner-grep-pipe-verfehlt-zeile`; die Review-Reports bleiben unverändert. (2) F-2
  (LOW): der indirekte Beleg für Arm A bleibt (Entscheid des Hauptlaufs 2026-10-02: kein
  zusätzliches Drucken). (3) F-3 (INFO): die `SEEN_SECOND`-Wartebedingung ist innerhalb des
  15-s-Fensters nicht eigenständig falsifizierbar (hergeleitet, nicht gefahren); sie senkt das
  Flake-Risiko und ordnet die Gegenlesung, die Zusage trägt die zweite Gruppe — keine Aktion.
  (4) F-4 (INFO): die Routing-Phase hat denselben Aufbau; Adresse jetzt
  [`slice-sdk-routing-phase-verbindung-haertung`](../open/slice-sdk-routing-phase-verbindung-haertung.md)
  (Lesen am Stand `a420e223`: die Falsch-Grün-Lücke besteht für die drei Stream-Flächen, die
  HTTP-Fläche ist nicht betroffen — beides hergeleitet, im Folge-Slice zu messen). (5) F-5 (INFO):
  Kotlin-Zeilenlängen, kein Formatwerkzeug für Kotlin →
  `BEO-PGC/formatierungs-drift-ohne-gate` (5×). Package-Version und Release unberührt
  (Verifier: nur Test-Pfade unter `sdks`, keine Versionsdatei). Der Verifier benannte einen
  eigenen Fehler im ersten Anlauf (§6, „Gleichzeitige Läufe“).
- **Steering-Loop-Eintrag (Lerneintrag):** (1) Ein Negativ-Beleg über Abwesenheit braucht eine
  **Beobachtung, dass der Beobachter verbunden war**, bevor das Negativ-Fenster beginnt: `READY`
  heißt bei Stream-Clients „Konsument gestartet“, nicht „Verbindung steht“; die Vorbedingung wird
  durch eine Änderung beobachtet, die **vor** der Negativ-Gruppe beim Client ankommt (`SEEN`), und
  die Negativ-Gruppe selbst durch eine zweite Beobachtung abgeschlossen (`SEEN_SECOND`). Belegt am
  Arm-A/Arm-B-Paar: dieselbe Überlagerung ist ohne Härtung grün (Falsch-Grün), mit Härtung rot.
  Anwendung der verkörperten Regel `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, kein
  neuer Eintrag. (2) Ein `grep -q` als Empfänger einer Pipe unter `pipefail` in einem
  Hilfsskript ist eine Fehlerquelle, sobald der Erzeuger nach dem Treffer weiterschreibt: eine
  vorhandene Zeile kann verfehlt werden. Neuer Eintrag `BEO-PGC/runner-grep-pipe-verfehlt-zeile`
  (1×), Ausgang *geplant* → der Folge-Slice (a); nichts verkörpert, daher ohne „liegt in“. (3) Ein
  Beleg über die Zählung (statt über abgelesene Zeilen) trägt, wenn die Closure ihn als
  **indirekt** kennzeichnet.
- **Beobachtungs-Register (`../observations/`):** neu angelegt `BEO-PGC/runner-grep-pipe-verfehlt-zeile/`
  mit `evidence/slice-sdk-sse-filter-phase-verbindung-haertung.md` (Zähler 1×: ein gemessenes
  Vorkommen mit Anker, Review F-1; der Lauf des Implementers desselben Vorgangs ist keine zweite
  Datei; die Implementer-Meldung im Vorgänger-Slice F-4 ist wegen abweichender hergeleiteter
  Ursache und fehlendem Lauf-Beleg nicht gezählt). `BEO-PGC/formatierungs-drift-ohne-gate` um
  `evidence/slice-sdk-sse-filter-phase-verbindung-haertung.md` fortgeschrieben, Zähler **5×**
  (Review F-5, Nicht-Go-Fläche Kotlin wie der vierte Beleg; Auslegung des Triggers im state
  unverändert, Ausgang bleibt kein Gate, Trigger für Go nicht ausgelöst).
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, über dem Deckel): F-2 (LOW)
  und F-3 (INFO) sind Funde ≤ LOW vor dem Merge, nach dem Deckel kein `evidence/`, sondern die
  Finding-Kennungen in dieser Notiz. `BEO-PGC/pipe-maskiert-make-exit-code`: Nachbar der neuen
  Klasse (andere Hälfte der Pipe-Eigenschaft), kein Beleg. `BEO-PGC/docker-cache-ueberspringt-tests-still`,
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`,
  `BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`: kein Auftreten.
  `BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot`: unberührt; der Folge-Slice (a) behauptet
  keinen Zusammenhang.
- **Folge-Slices:** zwei Dateien in `open/` —
  [`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../in-progress/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
  (F-1) und
  [`slice-sdk-routing-phase-verbindung-haertung`](../open/slice-sdk-routing-phase-verbindung-haertung.md)
  (F-4, Routing-Phase). Reihenfolge: (a) vor (b), beide ändern `lib-sdk-route-fixture.sh`.
- **Risiken aus §6:** Laufzeit entfallen · NATS-Flake eingetreten → Adresse Register und
  Folge-Slice (a) · Überlagerungen entfallen (Arm A indirekt belegt) · `SEEN` im Versuch > 1
  entfallen · Rest-Grenze der ersten Gruppe entfallen · Docker-Cache kein Auftreten · Drei
  Sprachen kein Auftreten · Übersetzung der Integrationsprojekte weiter offen → Register ·
  Gleichzeitige Läufe eingetreten (eigener Fehler des Verifiers, ohne Auswirkung) · Kein Release
  entfallen. Routing-Phase: Adresse Folge-Slice (b).
- **Drei Paarungen:** Anker: Review F-1 und Verifikation §5/§6. Folge-Slices: (a) und (b) in
  `open/`. Register: `runner-grep-pipe-verfehlt-zeile` (neu), `formatierungs-drift-ohne-gate`
  (5×).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den Pfaden
`tools/harness/`, `sdks/*/` (nur Test-Verzeichnisse), `harness/` und `docs/user/` — eine
Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Treffer in §6:
`negativtest-ohne-bindung-an-seine-eingabe` (verkörpert; 24 `evidence/`-Dateien, bei der Closure
vom Planner nachgezählt), `docker-cache-ueberspringt-tests-still` (2),
`drei-sprachen-kopie-divergiert-am-randfall` (3 Dateien), `integrationsprojekt-uebersetzt-nicht-unbemerkt`
(1); bei der Closure neu: `runner-grep-pipe-verfehlt-zeile` (1), `formatierungs-drift-ohne-gate` (5).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
