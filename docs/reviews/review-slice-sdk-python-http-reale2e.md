# Review-Report: slice-sdk-python-http-reale2e — 2026-09-23

**Review-Art:** Code — der Diff (`git diff d99768f2..HEAD`, Implementations-Commit
`040991ad`) gegen den [Slice-Plan](../plan/planning/in-progress/slice-sdk-python-http-reale2e.md),
die ADRs und die Hard Rules (`AGENTS.md` §3) geprüft. **Nicht** gegen die
DoD — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Commit `040991ad` ([`LH-FA-SST-009`](../../spec/lastenheft.md), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)) —
[slice-sdk-python-http-reale2e](../plan/planning/in-progress/slice-sdk-python-http-reale2e.md)

**Skill:** `.harness/skills/reviewer.md` ( geschärft 2026-09-09: vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** glm-5.3-flash:cloud[1m] · **Datum:** 2026-09-23

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei,
> eine Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link.
> Ein `pfad`-Feld auf den **geprüften Gegenstand** ist davon nicht
> betroffen — es zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [Slice-Plan slice-sdk-python-http-reale2e](../plan/planning/in-progress/slice-sdk-python-http-reale2e.md) (§2 DoD, §3 Plan + Plan-Nachzug + §3.13-Suchlauf-Feld, §6 Risiken, §8)
- [ADR-0110](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Mechanik-Klasse, Festlegung 2/Folgepflicht 1) · [ADR-0107](../plan/adr/0107-python-pypi-zweites-sdk-package.md) (Ort/Import-Grenze) · [ADR-0057](../plan/adr/0057-http-grpc-api.md) (HTTP-Server-Vertrag, benutzt nicht erweitert)
- [`SPEC-018`](../../spec/pflichtenheft.md) (Draht), [`LH-FA-SST-006`](../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../spec/lastenheft.md)/[`LH-FA-SST-009`](../../spec/lastenheft.md)/[`LH-FA-CON-001`](../../spec/lastenheft.md)
- `AGENTS.md` §3 (Hard Rules), `harness/conventions.md` (MR-000/MR-001)
- Parent-Stand: `d99768f2` (beide Stände gemessen, `AGENTS.md` §3.13)
- Vorbilder: C#-/Kotlin-Runner (`tools/harness/run-sdk-csharp-integration-tests.sh`,
  `run-sdk-kotlin-integration-tests.sh`) und deren HTTP-Realserver-Tests

---

## Findings

### F-1 — Planpunkt `harness/mk/sdk.mk` still gestrichen

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Code-Review gegen Plan, Modul 10)
- `pfad`: docs/plan/planning/in-progress/slice-sdk-python-http-reale2e.md:145 (Plan §3) · harness/mk/sdk.mk:118
- `befund`: Der Plan §3 führt `harness/mk/sdk.mk` als Update-Punkt („der
  Kommentar-Block des Targets `test-sdk-python-integration` zieht um die
  HTTP-Phase nach"); der Diff enthält keine Änderung an dieser Datei, und
  das Auslassen ist nirgends deklariert — kein Plan-Nachzug-Eintrag, §7
  Closure-Notiz noch leer. Der Kommentar-Block des Targets trägt keine
  Phasen-Enumeration, sodass nichts falsch geworden ist; der committete
  Plan verspricht trotzdem eine Änderung, die der Zug still auslässt.
- `verifizierbar`: ja — `git diff d99768f2..HEAD -- harness/mk/sdk.mk` ist
  leer (Lese-Handlung, kein Gate).
- `klasse`: Planpunkt still gestrichen

### F-2 — Zwei harness/README-Zeilen tragen die bewegte Eigenschaft als ausstehende Zusage weiter

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Arbeit überholt stehenden Träger)
- `pfad`: harness/README.md:160 (Kotlin-Zeile) · harness/README.md:161 (C#-Zeile)
- `befund`: Die Kotlin-Zeile schließt „erweitert sich um den Python-HTTP-
  Abschnitt (Folge-Slice)", die C#-Zeile „erweitert sich Slice für Slice um
  die Kotlin- und Python-HTTP-Abschnitte" — beide Erweiterungen existieren
  seit diesem Zug; die Zeilen lesen als ausstehende Zusage weiter. Das
  committete §3.13-Feld führt nur die Python-Zeile von §Werkzeuge
  (harness/README.md:162), nicht diese beiden Nachbarn derselben Tabelle.
- `verifizierbar`: ja — Träger-Lese-Handlung; kein Gate fängt sie.
- `klasse`: arbeit-ueberholt-stehenden-traeger (bisher 3× dokumentiert:
  slice-091/093/094 — dieser Fall wäre das vierte Auftreten)

### F-3 — Vorher-referierende Zeitform im neuen Test-Docstring

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen) · Gedächtnis-Regel „nur
  Ist-Zustand, kein Vorher/Nachher"
- `pfad`: sdks/python/pgchangefeed/integration/test_http_realserver.py:8-10
- `befund`: Der Docstring begründet die Phase mit „the HTTP surface is the
  one delivery-path surface of this package whose origin-acceptance
  predated that rule" — sachlich durch [`[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) gedeckt (dessen
  §Prüfung-Zeile nennt nur gRPC/SSE/NATS als Regel-Adressaten), aber die
  Zeitform referenziert den Vorher-Zustand statt den Ist-Zustand; die drei
  Geschwister-Docstrings tragen keine solche Klausel. Grenzfall: ADR-Anker
  vorhanden, als Kopplung formuliert.
- `verifizierbar`: nein — Lese-Handlung.
- `klasse`: Chronik-Sprache im Kommentar (Grenzfall)

### F-4 — Wortfragment „degenerater" aus dem Kotlin-Plan weitergetragen

- `kategorie`: LOW
- `quelle`: Maintainability · offene Beobachtung
  `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
- `pfad`: docs/plan/planning/in-progress/slice-sdk-python-http-reale2e.md:163
- `befund`: Der Plan-Nachzug schreibt „degenerater Pfad mit Kopf-
  Regeneration" — wortgleiches Fragment aus
  `slice-sdk-kotlin-reale2e` §7 (docs/plan/planning/done/slice-sdk-kotlin-reale2e.md:359,
  „degenerater Pfad, in der Fixrunde"). Drittes Auftreten der offenen
  Lernklasse, diesmal in Plan-Prosa statt im Code-Kommentar; §8s
  Deklaration „nicht einschlägig — dieser Slice kopiert kein Form-Vorbild"
  hält für diesen Fall nicht. Der Runner selbst trägt das Fragment nicht
  (`grep -rn degenerater tools/` — kein Treffer).
- `verifizierbar`: nein — Lese-Handlung (`grep -rn degenerater docs/plan/`).
- `klasse`: formvorbild-kopie-traegt-deutsches-wortfragment-weiter — bei
  der Closure als **drittes Auftreten** in den Zähler; Modul 10 §Pflege
  (Kategorie schärfen oder Regel orten) wird dort fällig.

### F-5 — Stream-förmige Fehlerdiagnose der HTTP-Phase

- `kategorie`: LOW
- `quelle`: Maintainability (Slice-Plan §6, Risiko 1 — deklariert; realer
  Lauf grün)
- `pfad`: tools/harness/run-sdk-python-integration-tests.sh:269-271 (Fehler-
  text) · 220-242 (Insert-Fenster)
- `befund`: Die HTTP-Phase fährt die stream-förmige Fire-and-Forget-Insert-
  Schleife fort (id_base 330, `HTTP_SENTINEL`, bis zu 5 committete Zeilen
  in `feed_e2e_full` ohne Empfänger), und ihr Fehlschlagstext nennt
  Stream-Semantik („der Test empfing keine der committeten Aenderungen …
  ueber den Stream"), auch wenn der Registrations-Rundlauf scheitert.
  Funktionell grünpfad-korrekt (realer Lauf EXIT=0, DoD belegt); das
  Residuum ist die irreführende Diagnose im Fehlerpfad plus toter
  Insert-Arbeit im HTTP-Lauf.
- `verifizierbar`: ja — ein roter HTTP-Lauf würde den Text real zeigen
  (`make test-sdk-python-integration` mit mutierter Registrierung).
- `klasse`: Form-Residuum der Phasen-Parametrisierung

### F-6 — Zitat-Basis der Mechanik-Klasse reicht über den ADR-Wortlaut hinaus

- `kategorie`: INFO
- `quelle`: [`[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
- `pfad`: tools/harness/run-sdk-python-integration-tests.sh:2-7 · harness/README.md:162
- `befund`: Runner-Kopf und README-Zeile zitieren „[`[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung
  Festlegung 2/Folgepflicht 1" auch für die HTTP-Fläche; der ADR-Wortlaut
  (§Prüfung-Zeile) nennt nur „jede neue gRPC-/SSE-/NATS-Client-Fläche".
  Die Anwendung auf HTTP ist die im Slice-Plan §Bezug deklarierte
  Übertragung der Mechanik-Klasse; die C#/Kotlin-README-Zeilen zitieren
  denselben Anker für ihre HTTP-Phasen — etablierte Praxis, kein Verstoß.
  Sollte die ADR HTTP nennen sollen, wäre das ein neuer ADR mit
  `Supersedes`, kein in-place-Eingriff (`AGENTS.md` §3.5).
- `verifizierbar`: nein — Lese-Handlung.
- `klasse`: — (Zitier-Basis, kein Fehlermuster)

## Negativbefunde

- geprüft, ohne Befund: **Plan-Nachzug und Suchlauf-Feld** — alle drei
  Plan-Nachzug-Punkte real im Diff (`received_grep`/`sql_kind`-
  Parametrisierung, Träger-Writer mit beidseitigem Rest-Erhalt und
  Kopf-Regeneration, `API_TOKEN_ADMIN`/`API_TOKEN_READER` im Runner-Kopf);
  das §3.13-Feld ist committet.
- geprüft, ohne Befund: **Verhaltensgleichheit der drei Stream-Phasen**
  (gegen `git show d99768f2:tools/harness/run-sdk-python-integration-tests.sh`
  nachgemessen) — alle drei Stream-Tests drucken `RECEIVED change_id=…`
  als erstes Token; die neuen `received_grep`-Parameter sind gegenüber der
  alten festen Form (`RECEIVED .*table=…`) um die Erst-Token-Annahme
  schärfer, nicht breiter; die ident-Extraktion ist äquivalent
  (`change_id` ist der erste Schlüssel); Reject-Marker, Exit-0- und
  Feed-Running-Prüfung unverändert — kein Verhaltenstausch.
- geprüft, ohne Befund: **HTTP-Phase** — eigene Sentinel-/ID-Skala
  (`HTTP_SENTINEL`, 330; PK-Bereiche 301–305/311–315/321–325/331–335
  disjunkt), `REJECTED status=401`-Marker, SQL-Variante `consumer` gegen
  `cdc.consumer`, explizite Testdatei-Auswahl (`PGCHANGEFEED_TEST_FILE`,
  Dockerfile-`:?`-Guard — kein stiller Ausschluss,
  `BEO-PGC/test-runner-stiller-ausschluss` nicht einschlägig real).
- geprüft, ohne Befund: **Eingabeseiten-Bindung der HTTP-Testdatei** —
  fünf Env-Variablen (`PGCHANGEFEED_HTTP_ADDR`, `…_API_TOKEN_ADMIN`,
  `…_API_TOKEN_READER`, `…_SOURCE_ID`, `…_HTTP_PUBLICATION`), admin/reader-
  Paar wie der C#/Kotlin-Vorbild-Test (spiegelgleiche Bindungsform geprüft:
  keine Tabellen-Inhalts-Assertion auch dort); der 401-Test bindet an der
  Eingabe (Token „no-such-token"; bei gültigem Token bliebe die Ausnahme
  und der REJECTED-Marker aus — der Lauf färbt rot;
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` nicht einschlägig);
  die dokumentierte Mutation im DoD ist strukturell in der richtigen
  Richtung. Kein `# noqa`/`# type: ignore`/`//nolint` im neuen Code.
- geprüft, ohne Befund: **Träger** — 12 Zeilen gezählt (4 C# + 4 Kotlin +
  4 Python = 3 Sprachen × 4 Wege); `[LH-FA-SST-009](../../spec/lastenheft.md)` in jeder Zeile;
  Linktiefe `../../spec/lastenheft.md` korrekt; der Python-Abschnitt ist
  **byte-identisch** mit dem Runner-Generator (mechanisch verglichen —
  idempotente Form real); Parent-Stand ohne Python-Marker (`git show
  d99768f2:…` — 0 Treffer); die C#/Kotlin-Abschnitte sind im Diff rein
  additions — beidseitiger Erhalt real; keine Phantom-Zeile (alle vier
  Nachweis-Dateien existieren).
- geprüft, ohne Befund: **harness/README-Zeile** — Phasen-Liste mit
  `pgchangefeed.http_client`, cdc.consumer-Gegenprüfung und Träger-Hinweis
  vertreten; die Verschärfung „die SDK-Clients" → „die Stream-Clients"
  korrekt gegen die alte Fassung.
- geprüft, ohne Befund: **§3.12 Instanz A — Zahlen tragen Ursprung** —
  die „79/2"-Zahl in harness/README §doc-trace unverändert richtig,
  mechanisch nachgemessen: `make doc-trace` EXIT=0, „79 Anforderung(en),
  2 Waise(n)" (Exit-Code ungepiped gesichert, `AGENTS.md` §3.9); die
  DoD-Zahlen (gRPC 804-1, SSE 807-1, NATS 813-1,
  consumer_id=python-sdk-e2e-…) sind laufgebunden mit Ursprung („der Lauf
  trug") und der Form der gedruckten Runner-Zeile; die consumer_id-Form
  (12 Hex-Zeichen) hält gegen den Test-Code (`uuid.uuid4().hex[:12]`).
- geprüft, ohne Befund: **Zitate** — [`[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)s §Prüfung-Zeile verankert
  die Mechanik-Klasse an Festlegung 2/Folgepflicht 1 (die Runner-Zitier-
  kette folgt dem ADR-eigenen Anker, siehe F-6); „[`SPEC-018`](../../spec/pflichtenheft.md) Negative" ist
  durch das Pflichtenheft gedeckt (401 bei unbekanntem Bearer-Token);
  [`SPEC-027`](../../spec/pflichtenheft.md)s „real"-Aussage bindet an die Paketierung, nicht an einen
  Wire-Beleg — kein falsch werdender Träger im Pflichtenheft (die
  Suchlauf-Behauptung hält).
- geprüft, ohne Befund: **Suchlauf an beiden Ständen** — Parent-Träger
  ohne Python-Marker, Parent-README-Zeile mit Dreier-Phasen-Liste,
  Parent-Runner-Kopf mit Dreier-Form: die Befund-Zellen des §3.13-Felds
  stimmen mit beiden Ständen (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-
  messung` nicht einschlägig).
- geprüft, ohne Befund: **Runner-Kopf-Nachzug**
  (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`) — die Dreier-Form ist
  gezogen, die vierte Phase steht im Kopf; die übrigen Kopf-Aussagen
  (kein Server-E2E-Beleg, Marker-Disziplin) unverändert richtig.
- geprüft, ohne Befund: **compose-/Dockerfile-Vertrag** —
  `e2e-admin-token`/`e2e-reader-token`/`pub_pgc_e2e` stimmen mit
  `compose.yaml` (`CDC_API_TOKEN_ADMIN`/`CDC_API_TOKEN_READER`/
  `CDC_PUBLICATION`) überein; der Dockerfile-`COPY pgchangefeed/integration`
  trägt den Ordner — die Plan-Zeile „keine Änderung erwartet" hält.
- geprüft, ohne Befund: **Commit-Traceability** — `make commit-traceability`
  EXIT=0 (5 Commits, Betreffs ohne Struktur-ID); kein `SPEC-*` im Betreff.
- geprüft, ohne Befund: **Suppression/Hostpfade** — keine Treffer im neuen
  Code (`AGENTS.md` §3.2, §3.11).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Planpunkt still gestrichen ·
arbeit-ueberholt-stehenden-traeger (4. Auftreten) · Chronik-Sprache im
Kommentar (Grenzfall) · formvorbild-kopie-traegt-deutsches-wortfragment-
weiter (3. Auftreten — Steuerungs-Schwelle) · Form-Residuum der
Phasen-Parametrisierung

## Verdikt

**Merge-blockierend:** ja — die beiden MEDIUM-Findings gehen an den
Implementer (F-1: Deklaration des Auslassens im Plan-Nachzug/§7 oder der
versprochene Nachzug selbst; F-2: Zeilen-Nachzug der beiden
harness/README-Nachbarn oder benannter Aufschub mit Adresse). Beide sind
innerhalb des noch offenen Slice-Fensters (§7 Closure-Notiz steht aus)
billig auflösbar; die drei LOW-Findings sind dem Implementer zur freien
Entscheidung mitgegeben, das INFO ist ohne erwartete Aktion.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen**
gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler —
F-4 erreicht die Dreifach-Schwelle der offenen Beobachtung
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (Modul 10
§Pflege: Kategorie prüfen), F-2 verlängert `BEO-PGC/arbeit-ueberholt-
stehenden-traeger` auf ein viertes Auftreten. Die DoD-Checkbox „Review
durchführt" bleibt offen (Fixrunde angezeigt — kein Selbst-Nachzug,
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Dieser
Report ist ein Lauf-Beleg; die Verifikation gegen die DoD prüft der
Verifier separat (Modul 11).