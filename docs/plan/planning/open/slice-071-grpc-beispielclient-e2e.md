# Slice 071: gRPC-Beispiel-Client und E2E-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-19.

**Bezug:** [LH-FA-SST-008](../../../../spec/lastenheft.md) (Haupt-Bezug),
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md).

**Berührte Spec-Stellen:** — (Beispiel-Client und E2E-Rundlauf validieren
gegen das in `slice-069` festgelegte
[SPEC-019](../../../../spec/pflichtenheft.md), ändern keinen Spec-Inhalt).

**Verantwortlich:** —.

**Autor:** pt9912 (Rolleninhaber: Planner-Lauf, 2026-09-14). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Wegwerf-Beispiel-Client (`tools/harness/grpcclient/`) und die
`compose.yaml`-/`run-integration-tests.sh`-Erweiterung belegen real gegen den
laufenden Feed-Container: eine committed Änderung erreicht einen
verbundenen gRPC-Stream-Client mit vollständigem Inhalt, und ein
Verbindungsversuch ohne gültiges Token wird abgelehnt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Port/Broadcaster/Server-Grundgerüst und Capture-Integration** —
  bereits geliefert durch `slice-069`/`slice-070` (Bestand); dieser Slice
  belegt nur, dass sie zusammen funktionieren.
- **HTTP/SSE-Beispiel-Client und -E2E-Beleg** — anderer Vorgang:
  `slice-072` liefert einen eigenständigen Beleg für den zweiten
  Zustellweg; beide Belege sind laut
  [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Slice-Schnitt-Empfehlung voneinander unabhängige Nachweisformen für
  zwei verschiedene Konsumenten desselben `Broadcaster` und werden
  deshalb nicht in einem Slice zusammengefasst.
- **Stream-internes Replay und tabellen-granulare Filterung** — Bestand
  bleibt bewusst stehen, dieselbe Vertagung wie in `slice-069` §1
  begründet; ein E2E-Beleg dafür kann es folgerichtig nicht geben, solange
  die Fähigkeit selbst nicht existiert.
- **Produktions-Client-Bibliothek oder SDK** — anderer Vorgang: Der
  Beispiel-Client ist Wegwerf-Werkzeug außerhalb der Produktionsschichten
  (analog `tools/harness/natssub/`/`tools/harness/httpclient/`), keine
  öffentlich unterstützte Klient-Bibliothek.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [LH-FA-SST-008](../../../../spec/lastenheft.md) Happy-Path/Negative
      erfüllt, Test referenziert: `make test-integration` — ein Wegwerf-
      Beispiel-Client (`tools/harness/grpcclient/`) empfängt eine reale
      committed Änderung mit vollständigem Inhalt über den laufenden
      Feed-Container, **und** ein Verbindungsversuch ohne gültiges Token
      wird real abgelehnt (`gRPC Unauthenticated`) — beide Belege in
      derselben Rundlauf-Erweiterung
      ([ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) Fitness
      Function).
- [ ] `compose.yaml` exponiert `CDC_GRPC_ADDR`;
      `run-integration-tests.sh` fährt den gRPC-Rundlauf als Teil des
      bestehenden Compose-Integrationstests.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update `harness/README.md` §Sensors (`make test-integration`
      Zeile: neuer gRPC-Rundlauf-Satz analog zum bestehenden HTTP-API-Satz).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/grpcclient/` | neu | Wegwerf-Beispiel-Client, analog `tools/harness/natssub/`/`tools/harness/httpclient/` |
| `compose.yaml` | update | Port-Exposition `CDC_GRPC_ADDR` |
| `tools/harness/run-integration-tests.sh` | update | Rundlauf-Erweiterung: gRPC-Stream-Client verbindet, empfängt Change, negativer Token-Test |
| `harness/README.md` | update | `make test-integration`-Sensor-Zeile um gRPC-Rundlauf-Satz ergänzt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-070` liegt in `done/` — ohne die
Capture-Integration veröffentlicht `CaptureService` nie einen `Change`
über den `Broadcaster`, ein E2E-Beleg hätte kein reales Ereignis; WIP-Limit
frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Wenn Happy-Path-
  und Negative-Beleg sich als zwei unabhängig aufwändige Rundläufe
  herausstellen, die nicht in derselben Compose-Erweiterung belegbar sind.
- `in-progress` → `open` (blockiert — Carveout?): Wenn der Compose-Stack
  den gRPC-Port aus infrastrukturellen Gründen (z. B. Netzwerk-Alias-
  Kollision) nicht zusätzlich exponieren kann.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) + PR gemerged + `make gates` grün + realer
`make test-integration`-Rundlauf grün + Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `harness/README.md` wird von diesem Slice und parallel von `slice-061`
  berührt (unterschiedliche Zeilen — neue gRPC-Sensor-Satz vs. neue
  HTTP-API-Sensor-Zeile); Merge-Konflikt-Risiko, falls `slice-061` bei
  Start dieses Slice noch nicht gemerged ist. — **Ausgang:** wird bei
  Closure zugewiesen.
- Der Compose-Stack exponiert mit `CDC_GRPC_ADDR` einen weiteren Port;
  Netzwerk-Alias- oder Port-Kollision mit einem bereits belegten
  Compose-Dienst ist nicht ausgeschlossen, bevor real getestet wurde. —
  **Ausgang:** wird bei Closure zugewiesen.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Der berührte Pfad
(`tools/harness/`, `compose.yaml`, `harness/README.md`) qualifiziert keine
eigene Sub-Area gegenüber dem in `harness/conventions.md`
§Modus-Deklaration deklarierten Default `*` (Kürzel `PGC`) — bleibt unter
dem Default.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/`) durchgegangen, Treffer für `PGC`:

- `BEO-PGC/slice-chronik-in-code-kommentar` — 5×, **verkörpert**.
  Implementer-Warnung: keine Slice-/Wellen-Verweise in
  Produktionscode-Kommentaren des Wegwerf-Clients (Godoc-Testfall-
  Provenienz ist dort ohnehin nicht einschlägig, da kein `_test.go`).
- `BEO-PGC/commit-traceability-kein-vorab-hook` — 3×, Ausgang steht noch
  aus (`welle-16`-Closure zuständig). Implementer-Warnung:
  Commit-Betreffs vor `git commit` gegen `SPEC-\d+` prüfen.

Keine weiteren Treffer für die in diesem Slice berührten Pfade.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF — Default-Sub-Area `*` (Kürzel `PGC`), siehe
`harness/conventions.md` §Modus-Deklaration pro Sub-Area. Kein eigener
Sub-Area-Block nötig.
