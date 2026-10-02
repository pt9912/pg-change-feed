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

- [ ] **Liefer-Punkt 1 — Härtung im gemeinsamen Ablauf und C#-Tier.** In
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
- [ ] **Liefer-Punkt 2 — Kotlin- und Python-Tier.** Dieselbe Form in
      `SseFilterRealserverTest.kt` und `PhaseEnvironment.kt` (Kotlin) und
      `test_sse_filter_realserver.py` (Python; ein Umgebungswert mehr, gelesen wie die
      bestehenden); die Runner brauchen nur die Aktualisierung ihrer Abdeckungs-Zeile und
      Schlusszeile (der Phasenablauf liegt im Fixture). *Zu belegen durch:* je ein realer, grüner
      `make test-sdk-kotlin-integration`- bzw. `make test-sdk-python-integration`-Lauf nach
      `make image`, im Bericht die gedruckte `FILTER_RESULT`-Zeile je Sprache (erwartet bei
      Versuch 1 dieselben Zahlen `f1=2 … f2=2 … unfiltered=6`).
- [ ] **Liefer-Punkt 3 — Mutationsproben, davon der Härtungsbeweis.** Alle Mutationen laufen an
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
- [ ] **Nur Test-Code und Runner.** `git diff --name-only 91e46048 -- sdks` nennt ausschließlich
      Pfade unter den drei Test-Verzeichnissen (`PgChangeFeed.Client.Integration/`,
      `src/integrationTest/`, `integration/`), keine Versionsdatei; kein Release, kein Tag;
      `make sdk-public-doc-check` endet mit Exit 0. *Zu belegen durch:* der Diff-Befehl im
      Bericht und der Suchlauf in §3.
- [ ] **Die Abdeckung ist getragen.** Die Zeile je SDK im marker-gegrenzten Abschnitt von
      [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) (geschrieben von den drei
      Runnern, nicht von Hand) nennt die zweite Gruppe nach stehender Verbindung; Kennungen
      ([`LH-FA-SST-008`](../../../../spec/lastenheft.md),
      [`LH-FA-SST-009`](../../../../spec/lastenheft.md)) und Nachweis-Spalte unverändert.
      *Zu belegen durch:* `git diff` der Datei nach den drei Läufen zeigt genau die drei geänderten
      Zeilen, ein zweiter Lauf je Tier schreibt nichts, fremde Abschnitte bleiben; `make docs-check`
      Exit 0. Die Läufe der Tiers sind **seriell** und enden vor dem Commit der Datei (andere
      Runner und `make bench` schreiben Erzeugnisse im selben Arbeitsbaum; der Implementer
      committet ausschließlich eigene Pfade).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter
      `docs/reviews/review-slice-sdk-sse-filter-phase-verbindung-haertung.md` liegt vor, kein
      offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes je
      Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-sse-filter-phase-verbindung-haertung.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen `make test-sdk-*-integration`:
      Satz zur Filter-Phase um die zweite Gruppe nach beobachteter Verbindung ergänzt), die
      Kopf-Kommentare von `lib-sdk-filter-fixture.sh`, der drei Testklassen und der Runner
      tragen den Ist-Umfang; das Benutzerhandbuch bleibt unberührt (es beschreibt die Packages,
      nicht die Tiers).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7) — neues Verzeichnis oder
      weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
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
```

| Träger | Messung am Parent (`91e46048`, 2026-10-02) | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Zeile 1 `FILTER_RESULT` (Abschlusszeile, Parser, Beschreibungen) | 12 Zeilen | Fixture (Parser/Regex), drei Testdateien, drei Abdeckungs-Texte, `harness/README.md`; die Form der Zeile bleibt, der Wert ändert sich — jede Trefferzeile gelesen, ob sie eine Zahl nennt |
| Zeile 2 `Dreiergruppe` | 9 Zeilen — `lib-sdk-filter-fixture.sh` (3), `lib-sdk-route-fixture.sh` (3), `run-integration-tests.sh` (3) | nur die drei Treffer der Filter-Datei sind Gegenstand; die sechs fremden (Routing-Phase der Tiers, Server-Rundlauf) bleiben, **gemeldet** (§1) |
| Zeile 3 `drei Gruppen` | 20 Zeilen: Filter-Dateien (Fixture ×3, drei Testdateien je ×2) 9, davon fremd: Routing-Szenarien der drei Tiers (je ×3) 9, `docs/user/e2e-abdeckung.md` ×1 und `run-integration-tests.sh` ×1 | je Filter-Zeile lesen: „alle drei" meint die drei **Tabellen** der Gruppe und bleibt wahr; nachzuziehen nur, wo „genau eine Gruppe" gemeint ist; die elf fremden Treffer gehören nicht zum Gegenstand |
| Zeile 4 `ohne Filter alle drei` | 13 Zeilen (`sdk-e2e-abdeckung.md` ×3, `harness/README.md` ×3, Fixture ×1, drei Runner je ×2) | wie Zeile 3; die Abdeckungs-Texte der drei Runner nennen die zweite Gruppe |
| Zeile 5 `unfiltered=` | 14 Zeilen | Parser und Beschreibungen; erwarteter Wert bei Versuch 1: 6 |
| Zeile 6 `f1_foreign` | 16 Zeilen | Parser, Fehlertexte, Beschreibungen; Wert unverändert 0 |
| Zeile 7 Gegenlesung `new_data->>'name'` | 2 Zeilen in der Fixture-Datei | beide Stellen nehmen beide Sentinels |
| Zeile 8 `Filter-Phase` | 7 Zeilen (`harness/README.md` ×3, Fixture ×1, drei Runner je ×1) | Texte, die die Phase beschreiben: README-Satz nachziehen |
| Zeile 9 `dreizehn Phasen` | 4 Zeilen (drei Hilfetexte `harness/mk/sdk.mk`, ein Runner-Kopf) | **soll unverändert 4 bleiben** (keine vierzehnte Phase) |
| Zeilen 10 und 11 (neue Namen) | 0 Zeilen am Parent | Soll am Diff: `SEEN_SECOND` in Fixture-Kommentar, drei Testdateien; `SENTINEL_SECOND`/`FilterSentinelSecond` in Fixture, Kotlin-/C#-Umgebungs-Hilfe, Python-Test — der Implementer trägt die gezählten Werte ein |

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
  Ruhefenster bleibt 15 s. **Hergeleitet**, nicht gefahren. — **Ausgang:** *zu entscheiden bei
  Closure* (Messung: `SEEN`→`SEEN_SECOND` im Lauf).
- **Flake der bestehenden NATS-Phase.** Der Vorgänger sah einmal einen Ausfall der **alten**
  NATS-Phase (F-4: Fenster des Runners 20 s gegen 15 s des Tests, hergeleitet). Die Härtung
  verlängert diese Phase nicht; sie wird hier nicht untersucht und nur als Hinweis geführt. —
  **Ausgang:** *weiter offen* → bei Auftreten in den elf Läufen `evidence/` am passenden
  Register-Eintrag, sonst *kein Auftreten*.
- **Die Überlagerungen des Härtungsbeweises treffen die Reihenfolge nicht.** Die Rechnung
  (B bei t, zweites Schema t+2 s, A t+6 s, F1-Verbindung ≈ `READY` + 4 s, Commit der ersten Gruppe
  ≤ 1 s nach `READY`) lässt je ≥ 1 s Spielraum auf beiden Seiten; Beleg ist die
  `RECEIVED_F1`-Liste von Arm A (nur A der ersten Gruppe). — **Ausgang:** *zu entscheiden bei
  Closure*; trifft die Reihenfolge nicht, werden die Verzögerungen angepasst und der Lauf
  wiederholt, die Abweichung steht im Bericht.
- **`SEEN` im Versuch > 1.** Die Zahlen der ersten Gruppe sind dann größer; die exakte Prüfung
  trägt die zweite Gruppe (Sentinel), nicht `f1 == 2`. — **Ausgang:** *entfallen*, wenn der
  Runner die zweite Gruppe über den Sentinel zählt (Konstruktion); der Bericht nennt den Versuch.
- **Rest-Grenze der ersten Gruppe.** Für die erste Gruppe bleibt F-2 sinngemäß: ihr Negativ ist
  nicht belegt. Das ist gewollt — das Negativ trägt die zweite Gruppe; die erste trägt nur die
  Verbindung (`SEEN`). Ein Empfang belegt eine stehende Verbindung (**hergeleitet**: ein Client
  empfängt nur über seine Verbindung). — **Ausgang:** *zu entscheiden bei Closure* (entfallen,
  wenn Arm A/Arm B wie erwartet färben).
- **Docker-Cache der Stufe `integration`** (`BEO-PGC/docker-cache-ueberspringt-tests-still`). —
  **Ausgang:** *zu entscheiden bei Closure*; Beleg ist die gedruckte Zeile `FILTER_RESULT` mit
  gezähltem Fremdwert, nicht der Exit-Code des Baus.
- **Drei Sprachen, ein Randfall** (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). —
  **Ausgang:** *zu entscheiden bei Closure*; Gegenmaßnahme: Ablauf und Gegenlesung einmal in der
  Fixture-Datei, derselbe Eingabesatz je Sprache.
- **Kein Sensor übersetzt die Integrationsprojekte**
  (`BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`). — **Ausgang:** *weiter offen* →
  Register (der Slice fährt alle drei Tiers real).
- **Gleichzeitige Läufe im Arbeitsbaum.** Der Auftraggeber fährt SDK-Runner und `make bench` im
  Hintergrund; diese schreiben Erzeugnisse (z. B. `docs/user/sdk-e2e-abdeckung.md`,
  `docs/user/bench-abdeckung.md`) und belegen feste Container-Namen. — **Ausgang:** *zu
  entscheiden bei Closure*; Gegenmaßnahme: serielle Läufe, Commit nur eigener Pfade
  (`git add <Pfad>`, nie `-A`/`-a`).
- **Kein Release, keine Versionsänderung.** — **Ausgang:** *zu entscheiden bei Closure*;
  *entfallen*, solange `git diff` unter `sdks/` nur Test-Pfade zeigt (DoD „Nur Test-Code und
  Runner").

## 7. Closure-Notiz

Wird bei Closure gefüllt (Vorlage: Closure-Notiz des Vorgängers). Mindestinhalt: die gedruckten
`FILTER_RESULT`-Zeilen je Sprache mit Ursprung (gemessen/übernommen/hergeleitet,
[`AGENTS.md`](../../../../AGENTS.md) §3.12), die Mutationsmatrix (zwölf Zellen plus Arm A/B/
Kontrolle) mit Ursprung je Zelle, der Versuch, in dem `SEEN` fiel, und die gemessene Zeit
`SEEN`→`SEEN_SECOND`. Lerneintrag-Richtung (vom Implementer zu bestätigen): ein Negativ
braucht eine **beobachtete** Vorbedingung — bei Stream-Clients heißt „gestartet" nicht
„verbunden"; die Vorbedingung wird durch eine Änderung beobachtet, die **vor** der Negativ-Gruppe
beim Client ankommt. Anwendung der verkörperten Regel
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`; Entscheid über einen neuen Eintrag beim
Planner. Meldung an den Auftraggeber (offen): gleicher Aufbau der Routing-Phase (§1).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den Pfaden
`tools/harness/`, `sdks/*/` (nur Test-Verzeichnisse), `harness/` und `docs/user/` — eine
Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Treffer in §6:
`negativtest-ohne-bindung-an-seine-eingabe` (verkörpert; Vorgänger zählte 24 `evidence/`-Dateien,
**übernommen**, nicht nachgezählt), `docker-cache-ueberspringt-tests-still`,
`drei-sprachen-kopie-divergiert-am-randfall`, `integrationsprojekt-uebersetzt-nicht-unbemerkt`
(Stände beim Start aus dem Register zu lesen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
