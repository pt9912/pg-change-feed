# Welle welle-sdk-grpc-administration-flaeche — SDK-Erweiterung um die gRPC-Verwaltungs-API — Closure-Notiz

**Welle:** welle-sdk-grpc-administration-flaeche
**Abschluss:** 2026-09-29
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- **Die drei SDK-Packages decken die gRPC-Verwaltungs-API und den
  Stream-Filter — die Drei-Sprachen-SDK-Matrix für `LH-FA-SST-009` ist
  vollständig.** Jedes Package trägt einen `PgChangeFeedAdministrationClient`
  mit allen elf `Administration`-RPCs (RegisterConsumer, AcknowledgeConsumer,
  GetConsumerPosition, RemoveConsumer, EnableTable, DisableTable,
  GetTableStatus, ListTables, RunRetention, ReadChanges, Diagnose) und einen
  typisierten Fehlerklasse-Pfad je gRPC-Status-Code; der bestehende
  Stream-Client trägt die optionalen `schema`/`table`-Filter-Parameter
  ([`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md))
  additiv, der parameterlose Aufruf bleibt regressionsgetestet:
  - C#: `PgChangeFeed.Client` — `PgChangeFeedAdministrationClient`,
    `PgChangeFeedGrpcException`-Hierarchie (`slice-sdk-csharp-grpc-administration-flaeche`),
    109 Tests grün im Pack-Lauf.
  - Python: `pgchangefeed` — `administration_client.py`, `PgChangeFeedGrpcError`
    + sechs Unterklassen (`slice-sdk-python-grpc-administration-flaeche`),
    130 Tests grün im frischen Container-Lauf.
  - Kotlin: `pgchangefeed-kotlin` — `PgChangeFeedAdministrationClient`
    (elf `suspend fun`), `sealed class PgChangeFeedGrpcException`
    (`slice-sdk-kotlin-grpc-administration-flaeche`), 19 Suites/89 Testfälle
    0 failures im `--no-cache`-Lauf der Verifikation, nach dem V-1-Fix
    `tests="7", failures="0"` für die TableTest-Klasse.
- **Form: die generierten Protobuf-Nachrichten laufen durch, kein
  DTO-Layer** — dieselbe, im Repo etablierte Entscheidung des Stream-Clients,
  je Sprache im Plan als Deviation dokumentiert; die Feld-Kongruenz zum
  `.proto` ist damit strukturell gebunden, nicht nur testseitig belegt
  (Verifikations-Reports beider Slices §4/§5).
- **Rechtsklassen und Fehlerform konform zu
  [`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
  [`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md),
  [`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md):**
  beide Verifikations-Reports gleichen je RPC die ADR-Rechtsklasse gegen
  `administrationRPCRoles` (`internal/adapters/driving/grpc/interceptor.go`)
  und die Nachrichtenschemata gegen das `.proto` ab — kein Widerspruch;
  die Fehlerform-Mapping-Tabelle (400/401/403/404/500 plus Fallback) ist
  je Sprache gespiegelt und testseitig belegt.
- **Betreiber-Doku:** `docs/user/benutzerhandbuch.md` nennt in beiden
  gRPC-Abschnitten („Zugriff über den gRPC-Change-Stream", „Zugriff über die
  gRPC-Verwaltungs-API") alle drei SDK-Packages namentlich als abdeckend
  samt Aufrufform der Filter — der Verwaltungs-API-Absatz schließt mit
  „Drei-Sprachen-SDK-Matrix … vollständig"; Versionshistorie 1.80/1.81
  nachgezogen. Die drei `sdks/*/README.md` sind nachgezogen.
- **Drei Slices in `done/`:** `slice-sdk-csharp-grpc-administration-flaeche`,
  `slice-sdk-python-grpc-administration-flaeche`,
  `slice-sdk-kotlin-grpc-administration-flaeche` — jeder mit Review (Report
  unter `docs/reviews/`), Fixrunde, Verifikations-Report (DoD je „ja") und
  ausgefüllter Closure-Notiz.

## Was hat funktioniert?

- **Das fachliche Vorbild (Beispiel-Clients) plus das zuerst liefernde
  Geschwister-SDK trugen die Design-Entscheidungen vollständig vor.** C#
  ging voraus, Python und Kotlin folgten ohne offene Designfrage — die
  Protoc-Namenskonflikt-Kenntnis (`AdministrationOuterClass` in Kotlin)
  wurde aus dem Beispiel-Client-Review übernommen statt erneut entdeckt.
- **Die Rollen-Sequenz Implementer → Reviewer → Fixrunde → Verifier fing
  reale Mängel in drei unabhängigen Läufen:** 2 HIGH (deutsche
  Fixture-Strings aus dem C#-Formvorbild; falsche Suchlauf-Zahl gegen die
  eigene Messung), je 1 MEDIUM je Verifikation (stiller Test in Kotlin,
  selbstwidersprüchlicher Docstring in Python) — alle real gezogen
  (`bd10c391`, `48e04899`, `3b381f5b`, `8b12198e`) und nachgemessen.
- **Der Docker-only-Pack-Lauf war das erste, untrügliche Signal der
  fehlenden zweiten `COPY --from=proto`-Zeile** — der Kotlin-Compiler
  brach mit „Unresolved reference 'administration'" ab, bevor irgendein
  stiller Defekt entstehen konnte; der Fix wurde in derselben
  Rundschrift analog zum C#-Vorbild nachgezogen.

## Was ging anders als geplant?

- **Beide Sprach-Dockerfiles brauchten die zweite `COPY --from=proto`-Zeile
  für `administration.proto`** — beim Kotlin-Slice real erst am
  Pack-Lauf gefunden (dort als Abweichung im Plan nachgetragen), das
  Python-Dockerfile trug sie gleich planmäßig plus einen dritten,
  so nicht geplanten `sed`-Schritt gegen eine interne Kennung, die das
  gRPC-Plugin aus dem `.proto`-Kommentar in die generierten Docstrings
  kopiert (`tests/test_public_text.py` verbietet sie je Python-Quelldatei).
- **Die beiden Verifikations-Reports trugen je ein MEDIUM, beide nach dem
  eigenen Slice-Closure-Commit entdeckt und vor dieser Closure gezogen:**
  Kotlin V-1 (eine `@Test`-Methode läuft still nie — Nicht-void-Rückgabe
  entzieht sich der JUnit-Entdeckung; `3b381f5b`) und Python V-1
  (selbstwidersprüchlicher Stream-Filter-Docstring in der vom Zug selbst
  berührten `grpc_client.py`; `8b12198e`). Die Register-Arbeit dazu
  übernimmt diese Closure, weil beide Reports nach den Slice-Closure-Commits
  liegen.
- **Die Welle fuhr zwei Slices zugleich in einem Arbeitsbaum** — der
  Kotlin-Closure-Commit hob die Handbuch-Zeilen des Python-Slices unbenannt
  mit hoch (Review F-4, INFO; die Inhaltsseite war korrekt). Die Klasse ist
  als neues Register-Verzeichnis geführt (siehe unten); Konsequenz für den
  nächsten Doppel-Slice: Offenlegung im Commit-Text oder Trennung der Commits.
- **Die optional geplante Python-Integrationstest-Datei entstand nicht** —
  laut Plan §3 bewusst an das Slice-Zeitbudget gebunden; der Verifikations-Report
  prüft das als konform (kein fehlender Umfang). Der Realserver-Beleg der
  Administration-Fläche bleibt für alle drei Sprachen Gegenstand der
  bestehenden `make test-sdk-*-integration`-Züge.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die
im Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln.

**Diese Closure verkörpert keine neue Regel — kein gelesener Eintrag ist in
dieser Welle über die 3×-Schwalte hinaus ohne Ausgang gestanden.** Der
Lese-Schritt im Einzelnen:

- **Mit Beleg-Wachstum, Ausgang bereits zugewiesen:** `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
  (3× → 4×, HIGH, erste Test-Fixture-Zeichenkette; Ausgang unverändert
  *verkörpert* → `.harness/skills/reviewer.md` HIGH-Punkt „Form-Vorbild-Kopie
  trägt ein sprachgebrochenes Wortfragment weiter" — er hat den Fund real
  getragen) und `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (33× → 34×,
  MEDIUM; Ausgang unverändert *verkörpert* → `AGENTS.md` §3.13).
- **Neu angelegt, unter der 3×-Schwelle:** `BEO-PGC/test-methode-lauft-still-nicht`
  (1× — die JUnit-Entdeckung überspringt still eine Nicht-void-`@Test`-Methode;
  Nachbar-Eintrag `test-runner-stiller-ausschluss` bleibt auf seinen
  Runner-Filter-Mechanismus begrenzt) und `BEO-PGC/arbeitsbaum-race-ohne-offenlegung`
  (1× — paralleler Arbeitsbaum, Commit-Text nennt die mitgezogene fremde
  Datei nicht).
- **Von den Slices benannt, unverändert unter der Schwelle:**
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (2× — in allen drei
  Slices als §6-Risiko geprüft, kein neuer Randfall),
  `BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung` (1× —
  `requires-python` unverändert), `BEO-PGC/test-runner-stiller-ausschluss`
  (2× — in dieser Welle nicht ausgelöst).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Zähler steht in [`../observations/`](../observations/)
(`BEO-PGC/<slug>/evidence/`). Was in dieser Welle 3× erreicht hat, steht
oben unter *Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei
im Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse
wie ein halluziniertes Gate.

- Keiner — die Welle benennt keinen neuen Folge-Slice. Die Realserver-Belege
  der Administration-Fläche laufen über die bestehenden Züge
  `make test-sdk-csharp-integration`, `make test-sdk-python-integration`,
  `make test-sdk-kotlin-integration`, ohne dass ein Slice-Plan sie neu
  benennt; die drei Paarungen sind unten geprüft.

**Drei Paarungen (Schritt 3, Abschluss):** (a) *Anker* — diese Closure
schafft kein neues `liegt in`-Feld; die beiden von ihr berührten verkörperten
Einträge tragen ihre Anker unverändert (`.harness/skills/reviewer.md`
HIGH-Punkt, `AGENTS.md` §3.13 mit „seit welle-20"), geprüft am Ist-Stand.
(b) *Folge-Slice* — kein genannter Folge-Slice, nichts nachzuweisen.
(c) *Register* — jede in Welle- und Slice-Plänen genannte Kennung
`BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/`
(`drei-sprachen-kopie-divergiert-am-randfall`, `sdk-python-untergrenze-ohne-anwender-begruendung`,
`arbeit-ueberholt-stehenden-traeger`, `formvorbild-kopie-traegt-deutsches-wortfragment-weiter`,
`test-methode-lauft-still-nicht`, `arbeitsbaum-race-ohne-offenlegung`) —
grün in allen drei Fällen.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl). Gate-Exit-Codes ungefiltert direkt gesichert
([`AGENTS.md`](../../../../AGENTS.md) §3.9).

| Kriterium (Welle-Datei §3) | Beleg |
|---|---|
| Alle drei Slices in `done/` | `ls docs/plan/planning/done \| grep -c 'slice-sdk-.*-grpc-administration-flaeche'` druckt `3` (`slice-sdk-csharp-…` seit `76afcad4`, Kotlin `87e47213`, Python `8c0d9f12`, je reiner `git mv`) |
| `make gates` grün | Exit `0` am Stand `2047ef96`: `d-check: 1415 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `coverage-gate: OK — Coverage 80.40% erfüllt Schwelle 80%`, `generated-sync: OK`, a-check `gesamt: 0 Befund(e)`, `baseline-verify: v6.9.0 OK — 54 Dateien`. Der erste Lauf desselben Tages endete Exit `2` mit 6 Befunden (4× Review-Pfad-Token in den neuen Register-Belegen, 2× Linkziel auf die gezogenen Slice-Pläne) — beide Klassen in `2047ef96` gezogen |
| `make sdk-pack-csharp` grün | Exit `0`, Erzeugnis `sdks/csharp/dist/PgChangeFeed.Client.0.2.1.nupkg` |
| `make sdk-pack-python` grün | Exit `0`, Erzeugnisse `sdks/python/dist/pgchangefeed-0.2.1-py3-none-any.whl` + `.tar.gz`; frischer Testlauf im Build-Image 130 passed (Verifikations-Report §1) |
| `make sdk-pack-kotlin` grün | Exit `0`, Erzeugnisse `sdks/kotlin/dist/pgchangefeed-kotlin-0.2.2.jar` + `-sources.jar`; frischer `--no-cache`-Testlauf 19 Suites/89 Testfälle/0 failures (Verifikations-Report §1) |
| Handbuch-Träger | beide gRPC-Abschnitte nennen alle drei SDK-Packages als abdeckend, kein offener Folge-Schritt zur gRPC-Verwaltungs-API; Versionshistorie 1.80/1.81 |
| Review-Artefakte | `review-sdk-csharp-grpc-administration-flaeche.md`, `review-sdk-python-grpc-administration-flaeche.md` (je 1 HIGH, gezogen), `review-sdk-kotlin-grpc-administration-flaeche.md` (1 HIGH, 2 LOW, 1 INFO), `review-fixrunde-welle-sdk-grpc-administration-flaeche.md` (0 Befunde), `verifikation-slice-sdk-kotlin-grpc-administration-flaeche.md` und `verifikation-slice-sdk-python-grpc-administration-flaeche.md` (DoD je „ja") |

### Schritt 2 — Trigger-Audit

- **Carveouts: 0 offen.** `find docs -iname "CO-*.md"` liefert keinen
  Treffer — dieses Repo hat aktuell keinen offenen Carveout.
- **Bootstrap-aware Gates.** `harness/mk/coverage.mk` führt
  `THRESHOLD ?= 80` (Endstufe, Rampe ausgeschöpft); der Gate-Lauf dieser
  Closure druckt `coverage-gate: OK — Coverage 80.40% erfüllt Schwelle 80%`.
  **0 offen.**
- **ADR-Re-Evaluierungs-Trigger** (gelesen je §Re-Evaluierungs-Trigger):
  [`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md) — kein Bedarf
  an Netzwerktrennung, kein benannter `LH-FA-REA-*`/`SPEC-022`-Bedarf über
  gRPC (der eigene `ReadChanges`-Bedarf ist mit
  [`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md) bereits als
  eigene Entscheidung bearbeitet), keine Rechtsklassen-Verfeinerung benannt;
  [`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md) — kein
  Stream-Form-Bedarf, kein neues Rechtsklassen-Bedürfnis, kein
  Wachstums-Limit; [`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md) —
  permanent; [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) —
  kein Mehr-Tabellen-Bedarf, kein Broadcaster-seitiger Filterdruck benannt.
  **0 offen.**
- **Beobachtungs-Register:** Lese-Schritt oben — kein Eintrag über der
  Schwelle ohne Ausgang, kein stiller Verbleib.

### Schritt 4 — Archivierung

Das Repo führt kein eigenes Archivierungs-Werkzeug: kein Make-Ziel in
`Makefile` oder `harness/mk/*.mk` und kein Skript unter `tools/`. Das
externe Werkzeug `ai-harness-init archive-welle` committet mit
fest einprogrammierten Messages ohne `LH-*`/`ADR-*`-Kennung
(`BEO-PGC/externes-werkzeug-committet-ohne-kennung`) — dieselbe
Bedingung, die die Closures von `welle-transformationen`,
`welle-backfill-bestand` und `welle-nats-drittstream` als nicht
eingetreten festhalten. Die Bedingung des Schrittes (ein im Repo
geführtes Werkzeug) ist nicht eingetreten; keine Handarchivierung. Die
Entscheidung, den Lauf des externen Werkzeugs für diese Welle
gleichwohl auszuführen, liegt beim Betreiber.

### Ausgang des Closure-Note-Reviews

Ein eigener Review-Lauf dieser Closure-Notiz in frischem Kontext
(`.harness/skills/closure-note-reviewer.md`) liegt zum Zeitpunkt dieses
Schreibens **nicht** vor. Diese Abweichung ist benannt, nicht
verschwiegen: die Prüfung obliegt der nächsten Rolle im Workflow, nicht
diesem Planner-Zug selbst (kein Self-Review, Baseline-Regelwerk
`modul-08-agentenrollen.md`).
