# Welle sdk-reale2e: SDK-Realserver-E2E — die zwölf Zustellweg-Flächen der drei SDK-Packages tragen reale Server-Belege nach (`LH-FA-SST-009`, `ADR-0110`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-sdk-reale2e-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (wellenlos priorisiert, geschnitten aus
[`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2 — der etablierte Mechanismus — auf die zwei
weiteren Sprachen gespiegelt und um den HTTP-Retrofit der bestehenden
Flächen ergänzt, Nutzerentscheidung 2026-09-23, dokumentiert als
Roadmap-Vorschau-Zeile). **Datum:** 2026-09-23.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die drei SDK-Packages (`sdks/csharp/**` `PgChangeFeed.Client` 0.2.0,
`sdks/kotlin/**` `pgchangefeed-kotlin` 0.2.0, `sdks/python/**`
`pgchangefeed` 0.1.0) decken seit ihren Vollabdeckungs-Wellen zusammen
**zwölf** Zustellweg-Flächen (3 Sprachen × 4 Wege: HTTP, gRPC `SPEC-020`,
SSE `SPEC-021`, NATS-Vollinhalt `SPEC-024`) — jede einzelne Fläche ist
aber bislang **nur netzlos gegen Fakes** abgenommen worden
(`FakeCallInvoker`/`FakeHttpMessageHandler`/`FakeNatsClient` in
`PgChangeFeed.Client.Tests`, `FakeHttpTransport`/`FakeGrpcStreamTransport`/
`FakeSseTransport`/`FakeNatsStreamTransport` im Kotlin-Baum,
`httpx.MockTransport` im Python-Baum; die Ausnahme ist der
Python-Realserver-Test der drei neuen Flächen aus
`welle-sdk-python-vollabdeckung`). [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/Folgepflicht 1 hat die verschärfte Test-Pflicht
(realer, laufender Server zusätzlich zu Unit-Tests) bereits als Klasse
etabliert und für die drei Python-Neuflächen real umgesetzt — diese Welle
spiegelt denselben Mechanismus auf die übrigen neun Flächen (C#- und
Kotlin-Bäume vollständig, Python-HTTP-Fläche als Retrofit) in **einem
Bündel**.

Das *Mehr* gegenüber drei isolierten Slice-DoDs:
[`LH-FA-SST-009`](../../../spec/lastenheft.md)s Zielmatrix („drei SDKs ×
vier Zustellwege") ist erst vollständig belegt, wenn **jede der zwölf
Flächen** einen realen, laufenden Server-Beleg trägt — drei unabhängige
Slice-Closures belegen je nur ihre Sprache; keine einzelne DoD beobachtet,
dass die Matrix **als Ganzes** geschlossen ist. Dazu kommt der
parallele Abdeckungs-Träger
`docs/user/sdk-e2e-abdeckung.md` (stabile Abdeckungs-Deklaration der
SDK-Realserver-Belege im Muster von
[`e2e-abdeckung.md`](../../../docs/user/e2e-abdeckung.md), Erzeugnis der
SDK-Runner, kein Lauf-Beleg) samt seinem `trace.coverage`-Eintrag (Label
`SDK-E2E`) in `.d-check.yml` — die RTM-Sichtbarkeit der zwölf Belege ist
ein Zustand, den erst die Welle-Closure gegen die volle Matrix prüft;
kein Einzel-DoD trägt diese Prüfung. Der Belegbestand der
Server-Ebenen (`make test-integration`/`test-notify`/`test-replication`)
bleibt unberührt — dieser Träger ist der vierte seiner Art, nicht eine
Erweiterung des ersten.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln.

- `welle-sdk-python-vollabdeckung` liegt in `done/` — **bereits erfüllt**,
  real geprüft 2026-09-23 (`ls docs/plan/planning/done/ | grep
  sdk-python-vollabdeckung` listet `welle-sdk-python-vollabdeckung.md`
  und `welle-sdk-python-vollabdeckung-results.md`; Roadmap-Log-Zeile
  trägt Abschluss 2026-09-23). Damit ist der Mechanismus der
  verschärften Test-Pflicht ([`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  §Entscheidung Festlegung 2/Folgepflicht 1) real gebaut, real gelaufen
  und bewährt (`make test-sdk-python-integration`, drei reale
  Flächen-Läufe mit SQL-Gegenprüfung und Ablehnungs-Belegen —
  `harness/README.md` §Werkzeuge, Zeile `make test-sdk-python-integration`).
- Der Retrofit der HTTP-Flächen und die Spiegelung auf C#/Kotlin waren in
  den drei Vollabdeckungs-Wellen bewusst Out-of-Scope (belegt: die
  Out-of-Scope-Formeln in `welle-sdk-csharp-vollabdeckung.md`/`welle-sdk-kotlin-vollabdeckung.md`
  §6 — „kein Sprung in eine strengere Klasse ohne fachlichen Grund", und
  `welle-sdk-python-vollabdeckung.md` §6 — die verschärfte Pflicht galt
  ausdrücklich nur für die neuen Flächen). Die Nutzerentscheidung
  2026-09-23 (Roadmap-Vorschau-Zeile, in dieser Welle-Datei übernommen)
  liefert den fachlichen Grund für den Sprung: Parität der reale
  Server-Ebene über die volle Drei-Sprachen-Matrix.
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

- Alle drei Slices in `done/`.
- `make gates` grün (netzloser Gate-Satz, unverändert Docker-only).
- **Ein realer, grüner Lauf je Sprache**: `make test-sdk-csharp-integration`
  und `make test-sdk-kotlin-integration` je mit den vier Flächen-Phasen
  ihres Runners, `make test-sdk-python-integration` mit der erweiterten
  HTTP-Phase — gegen eine reale, laufende Server-Instanz; kein Gate
  (dieselbe Klasse wie `make test-integration`), aber Pflichtbeleg dieser
  Welle.
- Der Abdeckungs-Träger `docs/user/sdk-e2e-abdeckung.md` existiert real,
  trägt die Zeilen **aller** zwölf Flächen-Belege (3 Sprachen × 4 Wege)
  und ist über den `trace.coverage`-Eintrag (Label `SDK-E2E`) in
  `.d-check.yml` für `make doc-trace` sichtbar — beides real im
  `.d-check.yml`-Bestand.
- Kein Tag-Push (kein `sdk-*-v*`-Tag): Version-Hebung ist hier **nicht**
  nötig — der Testbestand dieser Welle ändert keine öffentliche
  API-Fläche und kein Package-Artefakt; ein Release-Zug bleibt
  Betreiber-Handlung ([`AGENTS.md`](../../../AGENTS.md) §3.10) und liegt
  außerhalb dieser Welle.
- Closure-Notiz in `welle-sdk-reale2e-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine.

| Slice | Titel | Bezug |
|---|---|---|
| slice-sdk-csharp-reale2e | C#-SDK: vier Realserver-Phasen (HTTP, gRPC, SSE, NATS-Vollinhalt); führt den C#-Integrationstest-Runner und den Abdeckungs-Träger samt `trace.coverage`-Eintrag ein | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung Festlegung 2 |
| slice-sdk-kotlin-reale2e | Kotlin-SDK: dieselbe Mechanik, vier Flächen-Phasen | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung Festlegung 2 |
| slice-sdk-python-http-reale2e | Python-SDK: HTTP-Fläche in den bestehenden Realserver-Runner nachziehen; erweitert den Träger um den Python-Abschnitt und prüft die Matrix-Vollständigkeit | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung Festlegung 2 |

**Reihenfolge:** Sequentiell in dieser Reihenfolge (csharp → kotlin →
python-http). Drei Gründe:

1. **Der Träger gehört zum ersten Slice (csharp), nicht zum letzten.**
   `docs/user/sdk-e2e-abdeckung.md` ist ein **Erzeugnis der Runner** —
   jedes der drei Runner-Skripte schreibt seinen eigenen Matrix-Abschnitt
   idempotent (Marker-gegrenzt, Muster
   [`e2e-abdeckung.md`](../../../docs/user/e2e-abdeckung.md): „Erzeugt
   von …, stabile Abdeckungs-Deklaration, kein Lauf-Beleg"). Die Datei
   entsteht real mit dem ersten Runner; der `trace.coverage`-Eintrag
   (Label `SDK-E2E`) kommt im selben Zug, weil er die RTM auf eine Datei
   verdrahtet, die von da an real existiert — und ihre Zeilen von da an
   wahr sind (sie trägt nur Zeilen real existierender Runner-Phasen, nie
   einen Beleg, der nicht gelaufen ist). Beim letzten Slice zu liegen
   würde die Datei zu einem retrospektiven Artefakt machen: ihre
   C#-/Kotlin-Abschnitte wären von bereits abgeschlossenen Läufen
   generiert — ein Nachtrags-Generieren wäre kein Lauf-Beleg des
   aktuellen Laufs und würde genau die §3.12-Klasse
   („Beleg trägt seinen Satz nicht") erzeugen, die diese Welle vermeiden
   will. Die drei Vorbilder des Trägermusters tragen ihre Entscheidung
   je am Einführungs-Ort ([`ADR-0104`](../adr/0104-benchmark-schwellen-per-001-002-003.md)/[`ADR-0105`](../adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md));
   der `trace.coverage`-Eintrag folgt dort der Datei, nicht der Welle.
2. **C# und Kotlin bauen je neue Infrastruktur, Python erweitert
   bestehende.** Der C#-Slice führt Runner-Form, Docker-Stufe und
   Make-Target erstmals für die Sprache ein (Muster des
   Python-Runners); der Kotlin-Slice ist die wortwörtliche
   Spiegelung desselben Musters auf den zweiten Baum — er profitiert
   real davon, dass die Form einmal durchlaufen ist. Der Python-Slice
   ist der kleinste: eine vierte Phase in einem bereits bewährten
   Skript plus Träger-Erweiterung plus Matrix-Vollständigkeits-Prüfung.
3. **Die Vollständigkeits-Prüfung gehört zum Schluss.** Der letzte Slice
   (python-http) trägt als DoD die Prüfung, dass der Träger die volle
   Matrix deklariert — sie ist erst nach den ersten beiden Slices
   überhaupt formulierbar, und die Welle-Closure prüft dasselbe Mehr
   (§3, Träger-Kriterium) gegen den abgeschlossenen Stand.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine andere Welle.
- Wird blockiert von: keiner anderen Welle;
  [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
  ist bereits `Accepted`, und der Trigger (§2) ist real erfüllt.
- **Geschwister-Slices:** der C#- und der Kotlin-Slice verändern
  disjunkte Bäume (`sdks/csharp/**` vs. `sdks/kotlin/**`, je eigene
  Runner-Skripte) — sie sind technisch unabhängig voneinander. Der
  vorgeschlagene sequentielle Zuschnitt folgt dem WIP-Limit 1 je
  Rolleninhaber (Modul 5), nicht einer technischen Kette. Der
  Python-HTTP-Slice hängt technisch von keinem der beiden ab (erweitert
  ein eigenes, bestehendes Skript), gehört aber in derselben Reihenfolge
  an das Ende, weil seine DoD die Matrix-Vollständigkeit des Trägers
  prüft (§4, Grund 3).
- Intern: strikt sequentiell, csharp → kotlin → python-http (Träger
  wächst je Slice, siehe §4).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1.

- **Eröffnungs-Entscheidung — kein eigener ADR für diese Welle.** Die
  Prüfung nach Modul 4 (§Kernidee: ADR trägt das „weil"; `AGENTS.md`
  §3.6: Gate-Lockerungen brauchen einen ADR) ergibt für jede
  ADR-auslösende Klasse „tritt nicht ein":
  1. **Keine Vertragsänderung** — die Flächen existieren bereits real
     (drei SDK-Packages, zwölf Flächen, committet und gepackt); die
     Welle ändert keine `LH-*`-Anforderung, keine `SPEC-*`-Festlegung
     und keinen öffentlichen API-Vertrag; die berührten
     Messmethoden-Annahmen der Spec-Stellen (`SPEC-018`, `SPEC-020`,
     `SPEC-021`, `SPEC-024`) bleiben unverändert — es kommen Belege
     einer höheren Test-Ebene hinzu, die Methode, die sie prüfen, ändert
     sich nicht.
  2. **Kein Gate, keine Schwellen-Änderung** — die drei (bzw. zwei
     neuen) Make-Targets bleiben Werkzeug, brauchen DB-Zugang/Docker/
     Netz und hängen nicht an `GATE_CHECKS` (dieselbe Klasse wie
     `make test-integration` und das bestehende
     `make test-sdk-python-integration`, das dieselbe Begründung trägt).
  3. **Der Mechanismus ist bereits entschieden und real gebaut** —
     [`ADR-0110`](../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
     §Entscheidung Festlegung 2/Folgepflicht 1 etablierte die
     Teststrategie-Klasse (Realserver-Integrationstest für SDK-Flächen,
     Delegation des Mechanismus an den umsetzenden Zug); der umsetzende
     Zug entschied die Architektur im Slice
     (`slice-sdk-python-grpc-client-flaeche.md` §3, in `done/`:
     eigenständiges Skript, geteilte Compose-Umgebung, SDK als Prüfling)
     und baute sie real. Diese Welle spiegelt denselben Mechanismus auf
     zwei weitere Bäume und erweitert dasselbe Skript um eine Phase —
     keine neue Entscheidungsklasse, sondern der bewährte Stand als
     Vorlage.
  4. **Das Trägermuster hat seine Form** — drei Instanzen
     ([`e2e-abdeckung.md`](../../../docs/user/e2e-abdeckung.md),
     [`bench-abdeckung.md`](../../../docs/user/bench-abdeckung.md),
     [`ci-matrix-abdeckung.md`](../../../docs/user/ci-matrix-abdeckung.md));
     von ihnen trugen zwei einen ADR, weil die *Abfrage-Architektur*
     ([`ADR-0105`](../adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md):
     REST-API-Abfrage) bzw. die *Schwellen* ([`ADR-0104`](../adr/0104-benchmark-schwellen-per-001-002-003.md))
     die Entscheidung trugen — nicht das Muster selbst. Der erste
     Präzedenzfall (`e2e-abdeckung.md`, vom Server-Runner generiert)
     kam ohne eigenen ADR aus; der vierte Instanz-Zug ist dem ersten
     strukturell identisch (generierte Abdeckungs-Deklaration des
     eigenen Runner-Laufs, kein Lauf-Beleg).
  5. **Die Nutzerentscheidung trägt die Roadmap** — die
     Vorschau-Zeile (Trigger, drei Slices, Träger, `.d-check.yml`-Eintrag)
     ist der dokumentierte Auftrag; die Welle-Datei (dieses Dokument, §1)
     trägt das „weil". Wo die drei Vollabdeckungs-Wellen den Retrofit
     bewusst Out-of-Scope ließen, benennt diese Welle-Datei die Umkehrung
     ausdrücklich (§2, zweiter Punkt) — dieselbe Begründungs-Disziplin,
     die dort die Abgrenzung trug.
     **Gegenprobe — wann der ADR nachzuholen wäre:** eine Aufnahme
     eines der neuen Targets in `make gates`/`GATE_CHECKS`, eine
     Änderung einer Spec-Messmethode oder ein Mechanismus-Wechsel
     gegenüber der `ADR-0110`-Form (z. B. separater Compose-Stack,
     Wegwerf-Client statt SDK-Prüfling) — jede dieser Bedingungen ist
     ein Re-Evaluierungs-Anker dieser Entscheidung (kein ADR heißt
     nicht: kein Prüfpunkt).
- **Keine Änderung an `spec/pflichtenheft.md`/`spec/lastenheft.md`** —
  `SPEC-027`s Deckungs-Aussage („deckt HTTP-API, gRPC-Stream, SSE und
  NATS-Vollinhalt") bleibt durch diese Welle richtig; die Beleg-Ebene
  kommt hinzu, der Vertrag ändert sich nicht. Ein Träger-Nachzug
  (Modul-5-Sinn) trifft hier keinen falsch-werdenden aktiven Träger —
  real geprüft (2026-09-23): weder `spec/pflichtenheft.md` noch die drei
  SDK-READMEs noch die `PgChangeFeedClientOptions`-Doku behaupten eine
  Teststrategie oder eine „Fakes-only"-Aussage; die falsch werdenden
  Formeln stehen in Records (`done/**`), die unberührt bleiben.
- **Keine Aufnahme in `make gates`** — die neuen Targets brauchen
  DB-Zugang/Docker/Netz (dieselbe strukturelle Grenze wie
  `make test-integration`); eine Gate-Aufnahme wäre eine eigene
  Schwellen-/Gate-Entscheidung ([`AGENTS.md`](../../../AGENTS.md) §3.6).
- **Kein `examples/**`-Bezug** — die Beispiel-Clients bleiben Doku mit
  Bau-Bindung ([`ADR-0087`](../adr/0087-beispiel-clients-csharp-kotlin.md),
  [`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md)); diese
  Welle baut keinen neuen Beispiel-Client und berührt die bestehenden
  nicht.
- **`spec/architecture.md`/`.a-check.yml`** bleiben unberührt —
  `sdks/**` liegt außerhalb der Server-Komponentensicht und wird von
  a-check nicht gelesen (dieselbe Formel wie die drei Vollabdeckungs-Wellen).
- **Kein Version-Bump der SDK-Packages, kein Publish-Workflow-Änderung** —
  die `.csproj`-`<Version>`, `pyproject.toml` `[project] version` und
  `build.gradle.kts`-Versionen sowie die drei `sdk-*-release.yml`-Workflows
  bleiben unverändert; der Realserver-Testbestand ist kein
  Artefakt-Vertrag.
- **`docs/user/e2e-abdeckung.md` bleibt unverändert** — ihr Erzeuger
  bleibt Eigentum des Server-Runners
  `tools/harness/run-integration-tests.sh`; der neue Träger ist ein
  paralleler, eigenständiger Abschnitts-Arbeitsraum der SDK-Runner
  (Roadmap-Vorschau-Zeile, übernommen).
- **Keine Erweiterung auf eine vierte SDK-Sprache** — unverändert
  `ADR-0106`/`ADR-0107` §Re-Evaluierungs-Trigger 1.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
vollständig durchgesehen. Relevante Treffer:

- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert,
  [`AGENTS.md`](../../../AGENTS.md) §3.13, 23×) — dieser Eintrag trägt
  bereits SDK-Belege derselben Sub-Area-Familie; die bewegte Eigenschaft
  dieser Welle („welche Flächen tragen einen Realserver-Beleg") wird
  voraussichtlich in SDK-READMEs und `harness/README.md`-Zeilen
  beschrieben — der Suchlauf je Slice (beide Stände gemessen, Parent und
  Diff) prüft die README-Status-Abschnitte und die
  Werkzeuge-Zeilen ausdrücklich mit.
- `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
  (offen, 2×) — **einschlägig**: die C#-/Kotlin-Runner-Skripte und
  Dockerfile-Kommentare entstehen als Kopie des Python-Vorbilds
  (`tools/harness/run-sdk-python-integration-tests.sh`,
  `sdks/python/Dockerfile`); ein drittes Auftreten ist genau die
  Prüfpunkt-Klasse, die der Eintrag benennt — beide Slices tragen als
  Risiko (§6 ihrer Pläne), dass die übernommene Form (Skript-Kopf,
  Phasen-Dokumentation) auf verbliebene deutsche Wortfragmente geprüft
  wird, statt die Übereinstimmung mit dem Vorbild allein als Beleg zu
  werten.
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, unter der
  Schwelle) — betrifft die `-run`-Musterabdeckung des Server-E2E-Runners
  (`test/integration/integration_test.go`); **nicht** einschlägig für die
  neuen SDK-Runner — sie folgen dem Python-Muster der expliziten
  Testdatei-Auswahl je Phase (Umgebungsvariable mit laut scheiterndem
  Guard, kein stiller Ausschluss des Rests); die Deklarations-Hälfte
  (eine neue Runner-Phase, die niemand deklariert, fehlt still in der
  Tabelle) deckt der Träger ab, der je Phase eine Zeile trägt — beide
  Hälften sind in den Slice-Plänen verankert, kein dritter Treffer.
- `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (offen, 3×,
  Architect-Entscheidung) — gesichtet; die Slices berühren READMEs
  höchstens im Suchlauf-Nachzug, nicht als Schreib-Ziel; kein
  Handlungsbedarf hier.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert,
  [`AGENTS.md`](../../../AGENTS.md) §3.12 Instanz A) — der Träger
  `docs/user/sdk-e2e-abdeckung.md` trägt je Abschnittszeile einen
  Runner-Verweis als `Ort` (Form-Vorbild `e2e-abdeckung.md`, aber ohne
  `Datei:Zeile`-Anker — die Runner-Phase-Auswahl ist klassen-stabil, ein
  Zeilenanker driftete bei jeder Runner-Einfügung); je Slice schreibt sein
  eigener Runner den Abschnitt neu
  (idempotent, aus derselben Messung, die ihn belegt) — dieselbe
  Disziplin, die die Python-Vorbild-Kette in ihrer Closure geschärft hat.
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 2×) — der
  Python-Runner-Kopf trägt die Phasen-Liste der Flächen („gRPC, SSE,
  NATS"); der HTTP-Retrofit-Slice zieht sie nach (vierte Phase), ohne
  die überholte Form still stehen zu lassen — im Plan des Slices als
  DoD-Verankerung benannt.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 8×,
  Schwelle erreicht) — die neuen Make-Target-Zeilen in
  `harness/README.md` behaupten Verhalten (Beleg-Aufbau je Phase), das
  der umsetzende Slice real gefahren haben muss, bevor die Zeile
  entsteht (AGENTS.md §4 — kein Träger nennt ein Target, das es nicht
  gibt; umgekehrt: eine Zeile über ein Target trägt ihren realen Lauf).
- `BEO-PGC/github-actions-unverifizierbar-lokal` (offen, 3×) —
  **nicht** einschlägig: diese Welle ändert keinen GitHub-Actions-Workflow
  (die drei `sdk-*-release.yml` bleiben unverändert, keine
  Strukturänderung, kein neuer Workflow); `AGENTS.md` §3.10 greift nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**.
Die beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/`
auflösen, nicht vom Schreibort.

Ergebnis: [welle-sdk-reale2e-results.md](done/welle-sdk-reale2e-results.md)
Zähler: [observations/](observations/README.md) (Beobachtungs-Register)