# Slice sdk-python-http-client-flaeche: Öffentliche HTTP-API-Client-Fläche (`SPEC-018`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md).

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

**Verantwortlich:** — bis zur Priorisierung.

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

- [ ] `sdks/python/pgchangefeed/src/pgchangefeed/http_client.py` (oder
      gleichwertiger Modulname) trägt eine öffentliche Client-Klasse mit
      einer Methode je der neun Port-gedeckten Fähigkeiten von
      [`SPEC-018`](../../../../spec/pflichtenheft.md) plus dem
      Changes-Lesen (`SPEC-022`) — Signatur, Request-/Response-Form und
      Fehler-Antwortform (`400`/`401`/`403`/`404`/`500` → typisierte
      Exception-Hierarchie, konsistent über alle Methoden) spiegeln
      `SPEC-018`/`SPEC-022` exakt. Bearer-Token wird bei Konstruktion
      übergeben (kein globaler State).
- [ ] Eigene Tests (`pytest`) decken je Fähigkeit mindestens den Happy Path
      und die Auth-Boundary (`401` fehlendes/unbekanntes Token, `403`
      `reader`-Token gegen einen `admin`-Endpunkt) ab — netzlos prüfbar
      (kein realer Server nötig, `httpx`-Mock-Transport analog dem
      `HttpMessageHandler`-Fake-Muster der C#-Fläche).
- [ ] Kein Import aus `internal/**`/`cmd/**` dieses Repos (`ADR-0107`
      §Entscheidung Festlegung 3, Import-Grenze) — real geprüft:
      `grep -rn "internal/\|cmd/" sdks/python/` liefert keinen Treffer.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die Python-HTTP-Oberfläche (analog dem C#-Eintrag,
      `ADR-0106` §Konsequenzen Folgepflicht 4 als Präzedenzfall) — getragen
      durch die bereits verkörperte Selbstprüf-Instruktion
      (`.claude/commands/implement-slice.md` Schritt 17) und den
      Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md)
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
- Die Fehler-Antwortform (`{"error": "<Klartext>"}`) lässt sich auf
  unterschiedliche Arten in Python abbilden (Exception-Hierarchie vs.
  Ergebnis-Tupel/`Result`-Typ) — eine falsche Wahl bindet spätere Consumer
  an ein API-Design, das ein Major-Bump bräuchte, um es zu ändern.
  **Ausgang:** weiter offen, entschieden beim Schreiben — `ADR-0107`
  Festlegung 4 bindet die PEP-440-Major-Boundary an Draht-Änderungen,
  nicht an dieses interne Design; ein API-Redesign bleibt vor `1.0.0`
  folgenlos möglich.
- Ein `httpx`-Mock-Transport für die Tests könnte reale Netzwerk-/
  Serialisierungs-Eigenheiten (z. B. Groß-/Kleinschreibung der
  JSON-Felder, Timeout-Verhalten) verdecken, die erst gegen einen echten
  Server auffielen. **Ausgang:** weiter offen — ein realer Rundlauf-Beleg
  bleibt `make test-integration`s bestehendem `tools/harness/httpclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import, `ADR-0068`); dieses SDK
  bekommt frühestens mit einem Folge-Slice einen eigenen
  Integrationsbeleg.

## 7. Closure-Notiz

<…>

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
