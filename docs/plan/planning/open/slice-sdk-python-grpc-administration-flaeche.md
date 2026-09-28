# Slice slice-sdk-python-grpc-administration-flaeche: Python-SDK — Administration-Client + Stream-Filter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`.

**Welle:** welle-sdk-grpc-administration-flaeche.

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md).

**Berührte Spec-Stellen:** — (kein SPEC-/ARC-Eintrag über `LH-FA-SST-009`
hinaus).

**Verantwortlich:** —.

**Autor:** Planner-Agent. **Datum:** 2026-09-28.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Python-SDK-Package `pgchangefeed` bekommt einen
`PgChangeFeedAdministrationClient` mit allen elf `Administration`-RPCs, und
der bestehende `grpc_client.py`s `stream_changes()` bekommt die optionalen
`schema`/`table`-Filter-Parameter (`ADR-0133`) — analog zu
`examples/grpc-client` (Go, bereits vollständig umgesetzt, direktes
fachliches Vorbild, da Python typmäßig näher an Go als an C#/Kotlin liegt —
keine generierten Stubs sind hier committet, sie entstehen wie beim
Beispiel-Client-Vorbild `grpc_client.py` im Bau aus der `.proto`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Server-seitige Änderungen** — anderer Vorgang, siehe C#-Geschwister-Slice
  §1 für dieselbe Begründung.
- **C#-/Kotlin-SDK** — eigene Slices derselben Welle.
- **Neue Beispiel-Clients** — `examples/grpc-client` ist bereits fertig,
  dieser Slice liest ihn nur als Vorbild.
- **`get_diagnose()`-SDK-Methode (HTTP)** — dieselbe Fähigkeit steht bereits
  über `diagnose()` (gRPC) im Scope; HTTP-Diagnose ist Gegenstand einer
  eigenen, hier nicht geplanten Folge-Welle.
- **Änderung der Python-Untergrenze (`requires-python`)** — die neue
  Administration-Client-Klasse braucht keine neuere Sprachfunktion als der
  bestehende `grpc_client.py`; `BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung`
  (1×, offen) bleibt unberührt, solange `requires-python`/das Basis-Image
  unverändert bleiben.
- **Breaking Change an `stream_changes()`** — die neuen Parameter sind
  optional/additiv (Default `None`/leer), die bestehende Signatur ohne
  Filter bleibt aufrufbar.

## 2. Definition of Done

- [ ] `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle
      elf RPCs mit typisierten Requests/Responses (`dataclasses`, Formvorbild
      `models.py`) und einer typisierten Fehlerklasse für die gRPC-Status-Codes
      ab (Formvorbild `exceptions.py`); Unit-Tests je RPC.
- [ ] `ADR-0133` erfüllt: `stream_changes()` trägt optionale `schema`/`table`-
      Parameter, `None`/leer = ungefiltert (Regressionstest für den
      parameterlosen Aufruf).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte nennen Python
      jetzt mit der vollen Fläche; Versionshistorie nachgezogen.
      `sdks/python/README.md` nachgezogen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen sind getragen — von der nächsten Welle-Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/src/pgchangefeed/administration_client.py` | neu | elf RPC-Methoden, Formvorbild `examples/grpc-client/{consumer,tables_admin,retention,changes,diagnose}.go` |
| `sdks/python/pgchangefeed/src/pgchangefeed/models.py` | update | typisierte Request-/Response-`dataclasses`, 1:1 aus `proto/cdc/administration/v1/administration.proto` |
| `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py` | update | Fehlerklasse für die gRPC-Status-Codes, Formvorbild bestehende HTTP-Fehlerklasse |
| `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py` | update | `stream_changes()` um optionale `schema`/`table`-Parameter erweitern |
| `sdks/python/pgchangefeed/tests/test_administration_client.py` | neu | je RPC ein Happy-/Boundary-/Negative-Fall, Regressionstest für `stream_changes()` ohne Filter |
| `sdks/python/pgchangefeed/integration/test_grpc_administration_realserver.py` | neu (optional) | Realserver-Beleg, Formvorbild `test_grpc_realserver.py`/`test_grpc_rule_realserver.py`, nur falls im Slice-Zeitbudget |
| `docs/user/benutzerhandbuch.md` | update | beide gRPC-Abschnitte, Versionshistorie |
| `sdks/python/README.md` | update | Administration-Client dokumentieren |

## 4. Trigger

**Start** (`next` → `in-progress`): Welle `welle-sdk-grpc-administration-flaeche`
eröffnet (erfüllt).

**Rückführungen:**

- `in-progress` → `next` (zu groß): mehr als drei Fixrunden nötig.
- `in-progress` → `open` (blockiert): ein Server-Fehler blockiert einen
  bestimmten RPC — Carveout auf den betroffenen RPC.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`** (2×, offen) —
  dieselbe Randfall-Divergenz-Gefahr wie beim C#-Geschwister-Slice (§6
  dort), hier speziell: Python liest ein `bytes`-Feld (`old_image`/
  `new_image`) möglicherweise anders als C#/Kotlin bei Abwesenheit
  (`None` vs. leeres `bytes`-Objekt) — real gegen die generierten Stubs
  prüfen, nicht annehmen. — **Ausgang:** weiter offen: → Register (unter
  der Schwelle).
- **Server-Fehler analog zum Kotlin-Beispiel-Client-Fund** (`EnableTable`
  über gRPC, separat gefixt) — siehe C#-Geschwister-Slice §6. —
  **Ausgang:** eingetreten: Blockade melden | entfallen: Fix bereits
  gepusht | weiter offen.
- **Kein generierter Stub committet** (anders als C#/Kotlin, deren
  Build-Systeme die Stubs in einem eigenen Verzeichnis erzeugen) — die
  Administration-Client-Klasse muss denselben Bau-Kontext (`--build-context
  proto=proto`, `ADR-0090` Festlegung 2) wie `grpc_client.py` bereits nutzt,
  korrekt für die zweite `.proto`-Quelle (`administration.proto`) erweitern.
  — **Ausgang:** eingetreten: Dockerfile-Korrektur im selben Slice | entfallen:
  Bau-Kontext deckt bereits beide Quellen ab (real prüfen, siehe
  `sdks/python/Dockerfile`) | weiter offen.

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird bei Closure gefüllt>
- **Was ging anders als geplant:** <wird bei Closure gefüllt>
- **Steering-Loop-Eintrag:** <wird bei Closure gefüllt, falls einer entsteht>
- **Beobachtungs-Register (`../observations/`):** <wird bei Closure gefüllt>
- **Folge-Slices:** <wird bei Closure gefüllt, falls einer entsteht>
- **Risiken aus §6:** <wird bei Closure gefüllt — siehe §6>
- **Drei Paarungen:** von der Welle-Closure (`welle-sdk-grpc-administration-flaeche`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** eine Sub-Area, `sdks/python/`
(bestehendes, aktiv gepflegtes Package mit eigenem Release-Workflow).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer unter der 3×-Schwelle, als Risiko in §6 aufgenommen:
`drei-sprachen-kopie-divergiert-am-randfall` (2×),
`sdk-python-untergrenze-ohne-anwender-begruendung` (1×, nur relevant, falls
dieser Slice `requires-python`/das Basis-Image anfasst — laut §1 nicht
geplant). Kein Treffer ≥ 3× offen und einschlägig.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (neue Dateien in
einer bestehenden, bereits Greenfield geführten Sub-Area).
