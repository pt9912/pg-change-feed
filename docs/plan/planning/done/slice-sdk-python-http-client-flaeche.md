# Slice sdk-python-http-client-flaeche: Öffentliche HTTP-API-Client-Fläche (`SPEC-018`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-lh-fa-sst-009](welle-sdk-python-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (die neun Port-gedeckten
Fähigkeiten, die diese Fläche als Consumer anspricht),
[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 1 (Umfang — HTTP-API only), [`ADR-0057`](../../adr/0057-http-grpc-api.md)
(HTTP/JSON-API-Vertrag, wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md)
(Endpunkte, Token-Header-Form — das SDK benutzt diese Festlegungen,
verändert sie nicht), [`SPEC-022`](../../../../spec/pflichtenheft.md)
(Changes-Lesen).

**Verantwortlich:** Implementer-Agent (direkt beauftragt).

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0107` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Python-API-Fläche für alle neun
Port-gedeckten Fähigkeiten von [`SPEC-018`](../../../../spec/pflichtenheft.md)
(`RegisterConsumer`, `AcknowledgeConsumer`, `GetConsumerPosition`,
`RemoveConsumer`, `EnableTable`, `DisableTable`, `GetStatus`,
`ListTables`, `RunRetention`) sowie das Changes-Lesen
(`GET /changes`, [`SPEC-022`](../../../../spec/pflichtenheft.md),
[`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)) — Bearer-
Token-Auth (`reader`/`admin`), eigene, netzlos prüfbare Tests
(`httpx`-Mock/Transport-Fake statt eines realen Servers).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **gRPC/SSE/NATS-Vollinhalts-Stream** — `ADR-0107` Festlegung 1 grenzt v1
  ausdrücklich auf HTTP-API ein (Welle-Plan §6 Out-of-Scope); anders als
  bei der C#-Welle gibt es dafür in dieser Welle **keinen** Folge-Slice —
  ein künftiges Python-gRPC-Package bleibt ein eigener, künftiger Zug
  (`ADR-0107` §Re-Evaluierungs-Trigger 2).
- **Diagnose/Health-Endpunkte** — `SPEC-018` grenzt sie ausdrücklich aus
  den neun Port-gedeckten Fähigkeiten aus; kein SDK-Umfang.
- **Änderung von `SPEC-018` selbst** — das SDK benutzt den bestehenden
  Draht-Vertrag, verlangt keine Vertragsänderung (`ADR-0107` §Kontext
  Bindung, Analogie zu `SPEC-023`s „Verhältnis zum Draht").
- **Ein realer Rundlauf-Beleg gegen einen laufenden Server** — dieses
  SDK bekommt frühestens mit einem Folge-Slice einen eigenen
  Integrationsbeleg (analog dem C#-SDK); dieser Slice prüft ausschließlich
  netzlos gegen einen HTTP-Transport-Fake.

## 2. Definition of Done

- [x] `sdks/python/pgchangefeed/src/pgchangefeed/http_client.py` (oder
      gleichwertiger Modulname) trägt eine öffentliche Client-Klasse mit
      einer Methode je der neun Port-gedeckten Fähigkeiten von
      [`SPEC-018`](../../../../spec/pflichtenheft.md) plus dem
      Changes-Lesen (`SPEC-022`) — Signatur, Request-/Response-Form und
      Fehler-Antwortform (`400`/`401`/`403`/`404`/`500` → typisierte
      Exception-Hierarchie, konsistent über alle Methoden) spiegeln
      `SPEC-018`/`SPEC-022` exakt. Bearer-Token wird bei Konstruktion
      übergeben (kein globaler State).
- [x] Eigene Tests (`pytest`) decken je Fähigkeit mindestens den Happy Path
      und die Auth-Boundary (`401` fehlendes/unbekanntes Token, `403`
      `reader`-Token gegen einen `admin`-Endpunkt) ab — netzlos prüfbar
      (kein realer Server nötig, `httpx`-Mock-Transport analog dem
      `HttpMessageHandler`-Fake-Muster der C#-Fläche).
- [x] Kein Import aus `internal/**`/`cmd/**` dieses Repos (`ADR-0107`
      §Entscheidung Festlegung 3, Import-Grenze) — real geprüft:
      `grep -rn "internal/\|cmd/" sdks/python/` liefert keinen Treffer.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Fixrunde durchlaufen (1 HIGH, 1 MEDIUM,
      `docs/reviews/review-slice-sdk-python-http-client-flaeche.md`);
      Fixrunden-Nachprüfung bestätigt beide Findings real behoben, kein
      neuer Fund, kein offenes HIGH/MEDIUM.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die Python-HTTP-Oberfläche (analog dem C#-Eintrag,
      `ADR-0106` §Konsequenzen Folgepflicht 4 als Präzedenzfall) — getragen
      durch die bereits verkörperte Selbstprüf-Instruktion
      (`.claude/commands/implement-slice.md` Schritt 17) und den
      Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
      Zusätzlich, proaktiv vor Review (vermeidet die C#-HIGH-Klasse
      `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` F-1):
      `sdks/python/README.md` §Status wurde im selben Zug nachgezogen — der
      Absatz behauptete noch „is added by a follow-up release" für genau
      die HTTP-Client-Fläche, die dieser Slice real liefert.
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
      [welle-sdk-python-lh-fa-sst-009](welle-sdk-python-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/src/pgchangefeed/http_client.py` (Arbeitsname) | neu | öffentliche API-Fläche für die neun Port-gedeckten Fähigkeiten + Changes-Lesen. |
| `sdks/python/pgchangefeed/src/pgchangefeed/models.py` (Arbeitsname) | neu | typisierte Request-/Response-Datenklassen (spiegeln `SPEC-018`/`SPEC-022` JSON-Schemas). |
| `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py` (Arbeitsname) | neu | typisierte Fehlerform für `400`/`401`/`403`/`404`/`500`. |
| `sdks/python/pgchangefeed/tests/test_http_client.py` (Arbeitsname) | neu | Happy Path je Fähigkeit (`SPEC-018`/`022` Referenz), Auth-Boundary `401`/`403`, netzlos über `httpx`-Mock-Transport. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für die Python-HTTP-Oberfläche, im selben Zug. |

**Ansatz:** Draht-Kenntnis kommt **direkt** aus
[`spec/pflichtenheft.md`](../../../../spec/pflichtenheft.md) §2
(`SPEC-018`/`SPEC-022`) und im Zweifel dem Go-Server-Code
(`internal/adapters/driving/http/`, nur zur Gegenprobe gelesen, nicht
importiert) — **nicht** aus einem `examples/python/`-Vorbild, das es nicht
gibt (`ADR-0107` §Kontext „Was das ändert"). Das C#-SDK
(`sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs`) dient
als **Struktur**-Vorbild (Aufbau der Fähigkeits-Methoden, Fehler-Mapping-
Muster), nicht als Draht-Vorbild — die tatsächliche Byte-/Feld-Form wird
gegen `SPEC-018`/`SPEC-022` selbst geprüft, nicht gegen den C#-Code
übernommen.

**Plan-Nachzug (Implementer, vor dem Gate-Lauf):**

- **Träger-Nachzug-Suchlauf durchgeführt (`AGENTS.md` §3.13):** vor dem
  Gate-Lauf `grep -rn "follow-up\|added by" sdks/python/README.md
  sdks/python/pgchangefeed/src/pgchangefeed/*.py` ausgeführt (Suchlauf
  gegen die bewegte Eigenschaft „Python-HTTP-Client-Fläche existiert",
  nicht gegen den eigenen Diff). Gefunden: `sdks/python/README.md:9`
  behauptete unverändert „is added by a follow-up release" für genau die
  Fläche, die dieser Slice liefert — dieselbe Fehlerklasse, die beim
  C#-Geschwister-Slice erst im Review auffiel (F-1,
  `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md`). Im selben
  Zug nachgezogen (§2 DoD „Doku-Update"), **vor** dem ersten Review-Lauf
  dieses Slices — anders als beim C#-Vorbild, wo der Implementer-Zug das
  README noch nicht gefunden hatte und die Fixrunde nötig war. Der zweite
  Treffer (`__init__.py`s „a follow-up release" über gRPC/SSE/NATS) ist
  keine überholte Aussage: gRPC/SSE/NATS bleiben nach `ADR-0107`
  Festlegung 1 tatsächlich außerhalb dieses Pakets, kein Drift.
- **Import-Grenzen-Grep-Kollision real geprüft und vermieden:** ein erster
  Entwurf des `http_client.py`-Docstrings nannte den Go-Server-Pfad
  wörtlich als Gegenprobe-Hinweis — das hätte den in §2 geforderten
  `grep -rn "internal/\|cmd/" sdks/python/`-Nulltreffer selbst verletzt
  (ein Kommentar, kein Import, aber derselbe Grep-Treffer). Umformuliert
  ohne den wörtlichen Pfad; `grep`-Lauf danach real 0 Treffer.

**Fixrunde nach Review (`docs/reviews/review-slice-sdk-python-http-client-flaeche.md`, 1 HIGH, 1 MEDIUM):**

- **F-1 (HIGH) behoben:** `sdks/python/pgchangefeed/pyproject.toml:22-24`
  behauptete weiterhin „die HTTP-Client-Flaeche selbst folgt erst mit
  slice-sdk-python-http-client-flaeche" — der vorherige Träger-Nachzug-
  Suchlauf (`grep -rn "follow-up\|added by" sdks/python/README.md
  sdks/python/pgchangefeed/src/pgchangefeed/*.py`) deckte weder die
  deutsche Formulierung „folgt erst" noch `pyproject.toml` (Glob schloss
  es aus) ab. Kommentar auf den Ist-Zustand gehoben (httpx wird von der
  jetzt vorhandenen HTTP-API-Client-Fläche genutzt), ohne Slice-Namen
  (`AGENTS.md` §3.7). Erweiterter Suchlauf danach real ausgeführt:
  `grep -rniE "follow-up|added by|folgt erst|folgt mit|noch nicht|steht
  noch aus" sdks/python/` — drei Treffer: `README.md:9` und
  `__init__.py:7` (beide bereits korrekt, beziehen sich ausschließlich
  auf gRPC/SSE/NATS außerhalb dieses Pakets, kein Drift) sowie
  `Dockerfile:8` (bezieht sich auf eine `pack`-Stufe für einen anderen,
  noch nicht gestarteten Folge-Slice — bleibt zutreffend, dieser Slice
  fügt keine `pack`-Stufe hinzu). Keine weitere Fundstelle.
- **F-2 (MEDIUM) dokumentiert, nicht per Amend korrigiert:** Die
  Commit-Message von `314b51f6` nennt „23 pytest-Tests" für eine
  Kategorisierung (Happy Path/Auth-Boundary/Statuscode-Mapping/malformter
  Body), die real nur **20** trägt — die „23" stimmt nur als Gesamtsumme
  inkl. der 3 vorbestehenden, unter keine der genannten Kategorien
  fallenden `ClientOptions`-Tests aus `test_options.py`
  (Vorgänger-Slice `slice-sdk-python-projektgeruest`). Nachgezählt:
  `test_http_client.py` = 20 Tests (Happy Path je der zehn Fähigkeiten,
  Auth-Boundary `401`/`403`, Statuscode-Mapping `400`/`404`/`500`/
  unerwartet, malformter `2xx`-Body), `test_options.py` = 3 Tests
  (vorbestehend, unrelated), Summe = 23 — deckt sich mit dem realen
  Docker-Testlauf „23 passed". Der Commit selbst wird nicht per
  `--amend` korrigiert (bereits gelandet, bereits reviewt) — diese
  Klarstellung ist der Beleg-Träger für die korrekte Zahl (`AGENTS.md`
  §3.12).

**Fixrunden-Nachprüfung durch den Reviewer**
(`docs/reviews/review-slice-sdk-python-http-client-flaeche.md`
§Fixrunden-Nachprüfung): beide Findings real und korrekt behoben — eigener
erweiterter Suchlauf des Reviewers bestätigt unabhängig keine weitere
Drift-Stelle für F-1; F-2-Zählung nachgemessen und bestätigt. Kein neuer
Fund. `make gates` erneut grün.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-python-projektgeruest`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Schreiben zeigt, dass Fehler-Antwortform und Datenklassen-Design für
  neun Fähigkeiten mehr als drei Liefer-Punkte brauchen — dann Aufteilung
  nach Fähigkeits-Gruppen (Consumer-Verwaltung / Tabellen-Verwaltung /
  Retention-und-Lesen), analog der C#-Fläche.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-python-projektgeruest`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Kein Python-Referenz-Client existiert — die Draht-Treue lässt sich nur
  gegen die Spec-Tabellen (und im Zweifel den Go-Server-Code) prüfen,
  nicht gegen ein bereits laufendes Vorbild** (`ADR-0107` §Kontext „Was das
  ändert", die zentrale, für diese Welle spezifische Risiko-Asymmetrie
  gegenüber der C#-Welle: Dort senkte das bestehende
  `examples/csharp/http-client`-Vorbild dieses Risiko bereits vor dem
  Schreiben; hier gibt es kein `examples/python/http-client`-Äquivalent).
  Ein Python-Team muss jede Feldbezeichnung, jeden Statuscode und jede
  Fehler-Antwortform direkt aus dem Klartext-Vertrag erschließen — eine
  falsch gelesene Spec-Zeile (z. B. Groß-/Kleinschreibung eines
  JSON-Felds, ein übersehener Statuscode) bliebe ohne Gegenprobe
  unentdeckt, bis ein realer Server-Rundlauf sie aufdeckt. **Ausgang:
  weiter offen** — dieses Risiko ist strukturell, nicht durch diesen Slice
  allein auflösbar; es wird durch besonders sorgfältiges Gegenlesen von
  `SPEC-018`/`SPEC-022` gegen jede Methode gemindert, aber ein realer
  Rundlauf-Beleg bleibt einem Folge-Slice vorbehalten (analog dem
  C#-Vorbild, das denselben Beleg ebenfalls erst später über
  `make test-integration`s Wegwerf-Client bekam). Dies ist der Grund, aus
  dem `ADR-0107` den Erst-Scope bewusst kleiner schneidet als beim
  C#-Package (§Entscheidung Festlegung 1) — die geringere Fläche senkt die
  Zahl der ungeprüften Annahmen, hebt das Risiko aber nicht auf null.
  **Ausgang: weiter offen, strukturell (unverändert) — real gemindert
  durch eine zweite, unabhängige Gegenprobe.** Der Reviewer prüfte jedes
  Feld jeder Dataclass in `models.py` Feld für Feld gegen `SPEC-018`/
  `SPEC-022` (keine Abweichung); die Verifikation hat dieselbe
  Feld-für-Feld-Gegenprobe ein zweites Mal, unabhängig und ohne Rückgriff
  auf den Review-Text, komplett neu durchgeführt
  (`docs/reviews/verifikation-slice-sdk-python-http-client-flaeche.md`
  §2) — ebenfalls keine Abweichung, für alle zehn Fähigkeiten plus die
  zwölf `Change`-Felder inklusive `old_image`/`new_image`-Namensgebung und
  `ListTablesResponse.retained`. Zwei komplette, unabhängige
  Durchläufe ohne Abweichung senken das Risiko real — sie ersetzen aber
  weiterhin nicht den einzigen Beleg, den ein Referenz-Client oder ein
  realer Server-Rundlauf liefern könnte: eine Feld-für-Feld-Lektüre kann
  eine falsch gelesene Spec-Zeile nicht entdecken, wenn beide Leser
  dieselbe Zeile gleich falsch läsen. Der reale Rundlauf-Beleg bleibt
  einem Folge-Slice vorbehalten.
- Die Fehler-Antwortform (`{"error": "<Klartext>"}`) lässt sich auf
  unterschiedliche Arten in Python abbilden (Exception-Hierarchie vs.
  Ergebnis-Tupel/`Result`-Typ) — eine falsche Wahl bindet spätere Consumer
  an ein API-Design, das ein Major-Bump bräuchte, um es zu ändern.
  **Ausgang: entschieden während der Umsetzung.** Der Implementer wählte
  eine Exception-Hierarchie (`PgChangeFeedError` und typisierte
  Unterklassen je Statuscode plus `PgChangeFeedMalformedResponseError`)
  mit einem zentralen `_handle`-Helfer, der sowohl den Erfolgs- als auch
  den Fehlerpfad aller zehn Methoden trägt (Review und Verifikation
  bestätigen unabhängig: kein methodenspezifischer Zweitpfad umgeht ihn).
  `ADR-0107` Festlegung 4 bindet die PEP-440-Major-Boundary an
  Draht-Änderungen, nicht an dieses interne Design — ein API-Redesign
  (z. B. hin zu einem `Result`-Typ) bleibt vor `1.0.0` folgenlos möglich.
- Ein `httpx`-Mock-Transport für die Tests könnte reale Netzwerk-/
  Serialisierungs-Eigenheiten (z. B. Groß-/Kleinschreibung der
  JSON-Felder, Timeout-Verhalten) verdecken, die erst gegen einen echten
  Server auffielen. **Ausgang:** weiter offen (unverändert) — ein realer
  Rundlauf-Beleg bleibt `make test-integration`s bestehendem
  `tools/harness/httpclient` vorbehalten (Wegwerf-Client, kein
  SDK-Import, `ADR-0068`); dieses SDK bekommt frühestens mit einem
  Folge-Slice einen eigenen Integrationsbeleg.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Implementer hatte den §3.13-Träger-Nachzug-
  Suchlauf bereits vor dem Erstreview korrekt aus dem C#-Geschwister-Fund
  gelernt und proaktiv angewandt (§3 Plan-Nachzug) — zwei Stellen
  (`README.md`, `__init__.py`) wurden dadurch bereits vor dem ersten
  Review korrigiert, anders als beim C#-Zyklus, wo derselbe README-Fund
  erst dem Reviewer auffiel. Die Fehler-Antwortform (zentraler
  `_handle`-Helfer für Erfolgs- **und** Fehlerpfad) hat beide bereits im
  C#-Review gefundenen Findings-Klassen (methodenspezifischer Zweitpfad,
  malformter 2xx-Erfolgskörper) von vornherein vermieden — Review und
  Verifikation bestätigen das unabhängig voneinander mit einer
  Negativbefund-Zeile. `make gates` blieb über den gesamten Zyklus
  (Implementer, Reviewer-Erstlauf, Fixrunden-Nachprüfung, Verifikation,
  diese Planner-Closure) durchgehend grün.
- **Was ging anders als geplant:** Trotz des proaktiven, aus dem
  C#-Vorkommen abgeleiteten Suchlaufs fand der Reviewer eine dritte,
  strukturell identische Trägerstelle (F-1, HIGH): `pyproject.toml:22-24`
  behauptete weiterhin, die HTTP-Client-Fläche „folge erst" mit genau
  diesem Slice — der Suchlauf hatte weder die Datei (Glob schloss sie
  aus) noch die deutsche Formulierung „folgt erst" (Wortmuster deckte nur
  „follow-up"/„added by") erfasst. Fixrunde behob den Fund
  (`d0da688b`), Fixrunden-Nachprüfung bestätigte unabhängig keine vierte
  Stelle. Zusätzlich trat in der §Fixrunden-Nachprüfung des
  Review-Reports selbst der wiederholt aufgetretene `id-unlinked`-
  Docs-Check-Stolperstein erneut auf (nackte `ADR-0107`-Erwähnung in
  einem zitierten Docstring-Ausschnitt, ohne Backticks) — vom Coordinator
  über einen eigenen `docs-check`-Lauf gefunden und in einem eigenen
  Commit (`7dd0ca68`) behoben, zwischen Fixrunden-Nachprüfung und
  Verifikation. Die Verifikation selbst hielt zusätzlich eine
  Beobachtung fest, die kein DoD-Verstoß ist, aber für künftige
  SDK-Slices bewusst gehalten werden sollte: Die Auth-Boundary-Tests
  (`401`/`403`) laufen nur **repräsentativ** über eine Methode
  (`remove_consumer`), nicht je der zehn Methoden einzeln — begründet
  über den zentralen `_handle`-Helfer, den ausnahmslos alle zehn Methoden
  durchlaufen, also gibt es keinen methodenspezifischen Zweitpfad, der
  einen Einzeltest bräuchte. Der Verifier stufte das als DoD-konform ein,
  exakt dieselbe Struktur trägt bereits das C#-Geschwister-Package
  (`PgChangeFeedHttpClientAuthBoundaryTests.cs`) — eine repo-konsistente,
  aber bewusst nicht vollständige Testabdeckungsform, die ein künftiger
  Reviewer/Verifier nicht neu herleiten muss.
- **Steering-Loop-Eintrag:** Zwei Beobachtungs-Register-Einträge
  fortgeschrieben, kein neuer Sensor:
  - `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (18. Beleg,
    `evidence/slice-sdk-python-http-client-flaeche.md`): F-1 zeigt, dass
    ein aus einem vorigen Vorkommen kopierter §3.13-Suchlauf zwei
    unabhängige Lücken zugleich tragen kann — eine Datei-Glob-Lücke
    (`pyproject.toml` nicht im Glob) und eine Sprach-/Wortmuster-Lücke
    (nur englische Formulierungen gesucht, die tatsächliche Fundstelle
    stand auf Deutsch). Beide Achsen (welche Dateien, welche
    Formulierungen) müssen unabhängig breit genug gewählt werden; ein aus
    dem letzten Fund übernommenes Muster garantiert keine Vollständigkeit
    am nächsten, strukturell ähnlichen Vorgang. Kein neuer Sensor — die
    Regel (`AGENTS.md` §3.13) bleibt Lese-/Suchlauf-Disziplin, kein
    Gate-fähiges Muster.
  - `BEO-PGC/report-nackte-id-ohne-link` (8. Beleg,
    `evidence/slice-sdk-python-http-client-flaeche.md`): dritter realer
    Beleg für den frühen Fang-Zeitpunkt „zwischen Reviewer- und
    Verifier-Zug" (nach Beleg 7, `slice-sdk-csharp-grpc-client-flaeche`)
    statt erst bei einem roten Verifier-`make gates`-Lauf (Beleg 6). Kein
    neuer Lese-Schritt — die Regel (`AGENTS.md` §3.9, verkörpert seit
    `slice-063`) trägt weiterhin.
- **Folge-Slices:** `slice-sdk-python-pack-werkzeug`,
  `slice-sdk-python-publish-workflow` — beide offen in
  [welle-sdk-python-lh-fa-sst-009](welle-sdk-python-lh-fa-sst-009.md).
- **Risiken aus §6:**
  - „Kein Python-Referenz-Client — Draht-Treue nur gegen Spec-Tabellen
    prüfbar" — **Ausgang: weiter offen, strukturell (unverändert)**, real
    gemindert durch eine zweite, unabhängige Feld-für-Feld-Gegenprobe
    (Reviewer + Verifikation, beide ohne Abweichung); kein realer
    Server-Rundlauf-Beleg bisher.
  - „Fehler-Antwortform-Design könnte binden" — **Ausgang: entschieden
    während der Umsetzung** (Exception-Hierarchie inkl. zentralem
    `_handle`-Helfer für Erfolgs- und Fehlerpfad); bleibt vor `1.0.0`
    folgenlos änderbar.
  - „`httpx`-Mock-Transport könnte reale Eigenheiten verdecken" —
    **Ausgang: weiter offen (unverändert)**, ein realer Rundlauf-Beleg
    bleibt einem Folge-Slice vorbehalten.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-python-lh-fa-sst-009](welle-sdk-python-lh-fa-sst-009.md)
  (noch offen — Pack-Werkzeug und Publish-Workflow stehen aus) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/python/` — mit
`slice-sdk-python-projektgeruest` bereits eröffnet (GF, siehe dessen §8),
keine erneute Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert) betrifft diesen Slice über den DoD-Punkt „Doku-Update"
oben; kein weiterer Treffer für diese Sub-Area (`grep`-Suche wie im
Welle-Plan §6 dokumentiert).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-python-projektgeruest`).
