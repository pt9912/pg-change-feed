# Architect-Review welle-6 — Trigger-Audit (Carveout · Bootstrap-aware Gate · ADR-Re-Evaluierung) und Beobachtungs-Register-Lese-Schritt (`BEO-PGC/rollen-verdrahtung`, 3×)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-12.

**Eingang:** [`docs/plan/planning/welle-6.md`](../planning/done/welle-6.md) §1–§6,
§Vermerk für den Trigger-Audit bei Closure ·
[`docs/plan/planning/done/slice-021-consumer-registrierung-zugriffsweg.md`](../planning/done/slice-021-consumer-registrierung-zugriffsweg.md),
[`docs/plan/planning/done/slice-022-consumer-bestaetigung-zugriffsweg.md`](../planning/done/slice-022-consumer-bestaetigung-zugriffsweg.md)
(je §6–§8) ·
[`docs/plan/adr/architect-review-slice-021.md`](architect-review-slice-021.md) ·
[`docs/reviews/verify-slice-022.md`](../../reviews/verify-slice-022.md)
(V-1, V-2) · [`ADR-0019`](0019-cli-driving-adapter.md) (Accepted,
`permanent`) · [`ADR-0020`](0020-http-grpc-optional.md) (Accepted, Trigger
„beobachtbarer API-Consumer-Bedarf") · [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)
(Accepted, `permanent`) · `spec/lastenheft.md` `LH-FA-SST-005`…007,
`LH-QA-SEC-001`…003 (Versionshistorie 0.3.0→0.5.0, Commits `9936e82`,
`9f5030d`) · `docs/plan/carveouts/` (nur `.gitkeep`) ·
`harness/conventions.md` §Modus-Deklaration (`PGC`, Greenfield, gesamtes
Repo) · `docs/plan/planning/observations/BEO-PGC/rollen-verdrahtung/`
(`observation.md`, `state.md`, `evidence/{slice-011,slice-021,slice-022}.md`) ·
Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur
(Closure-Schritt 2 Trigger-Audit, Schritt 3a/3b Lese-Schritt/Verkörperung) ·
`modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle.

**Ausgang:** Zwei unabhängige Züge, beide mit Verdikt:

1. **Trigger-Audit der Welle (Modul 6, Closure-Schritt 2):** Carveout:
   **0 aktiv.** Bootstrap-aware Gate: **0 fällig** (durchgehend
   Greenfield). ADR-Re-Evaluierung: `ADR-0020`s Trigger ist gegen
   `LH-FA-SST-006` **wörtlich eingetreten** — Re-Evaluierung durchgeführt,
   Verdikt **bestätigt** (kein Folge-ADR; Begründung in Zug 1c).
   `LH-FA-SST-007` löst **keinen** ADR-Trigger aus (unabhängig bestätigt,
   deckt sich mit `verify-slice-022.md` V-2). `ADR-0019`/`ADR-0046`:
   beide `permanent`, unverändert bestätigt.
2. **Beobachtungs-Register — Lese-Schritt (Modul 6, Closure-Schritt
   3a/3b):** `BEO-PGC/rollen-verdrahtung` erreicht mit
   `evidence/slice-022.md` 3×. Ausgang: **geplant** — Folge-Slice
   [`slice-023`](../planning/done/slice-023-rollen-spezifische-dsn-verdrahtung.md)
   (Skelett von diesem Lauf angelegt, in `open/`; Planner vervollständigt
   §2–§8 vor `next/`).

**Harte Regel eingehalten:** `ADR-0019`, `ADR-0020`, `ADR-0046` (alle
`Accepted`) werden von diesem Lauf **nicht** inhaltlich geändert — der
Trigger-Audit bestätigt sie nur; kein Wortlaut wurde angetastet. Das
Beobachtungs-Register-`state.md` von `BEO-PGC/rollen-verdrahtung` wird von
diesem Lauf **nicht** editiert — die Fortschreibung auf den zugewiesenen
Ausgang (`geplant → slice-023`) ist Planner-Arbeit bei der `welle-6`-
Closure-Notiz (Modul 8: 3b läuft Planner→Architect→Planner, der
*Rückweg* zum Planner schreibt). `welle-6.md` selbst und die beiden
Slice-Dateien in `done/` werden von diesem Lauf **nicht** editiert — die
restliche Welle-Closure (Schritte 1, 3c–6: Trigger-Prüfung, Closure-Notiz,
`git mv`, Paarungen, Archivierung, Roadmap-Fortschreibung) bleibt
Planner-Arbeit.

---

## Zug 1 — Trigger-Audit der Welle

### 1a. Carveout (Modul 7)

`docs/plan/carveouts/` enthält ausschließlich `.gitkeep` — kein aktiver
Carveout im Repo, damit auch keiner mit Bezug zu `welle-6`. Beide
Verifier-Läufe dieser Welle (`verify-slice-021.md`, `verify-slice-022.md`)
liefen `make gates` eigenständig mit Exit 0. **Feststellung: 0 offen.**

### 1b. Bootstrap-aware Gate (Modul 13)

`harness/conventions.md` §Modus-Deklaration führt weiterhin genau eine
Sub-Area (`*`/`PGC`, Greenfield) für das gesamte Repo; beide Slice-Pläne
bestätigen in §8 unabhängig „Alle berührten Sub-Areas GF". Keine
Brownfield-/Hybrid-Sub-Area, keine Reifestufe zum Hochschalten.
**Feststellung: 0 fällig** — die Frage stellt sich in einem durchgehend
Greenfield-Repo nicht. Keines der vier bestehenden Gates
(`baseline-verify`, `docs-check`, `a-check`, `commit-traceability`) ist
selbst gestuft; diese Welle führt kein fünftes ein.

### 1c. ADR-Re-Evaluierung — Schwerpunkt `ADR-0020` gegen `LH-FA-SST-006`

**Wortlaut-Prüfung.** `ADR-0020`s Re-Evaluierungs-Trigger: „Beobachtbarer
Bedarf eines API-Consumers — sichtbar als Anforderung im Lastenheft-Change
oder als Eintrag im Beobachtungs-Register." `LH-FA-SST-006` (Commit
`9936e82`, Version 0.3.0→0.4.0) fügt genau das ein: eine **`muss`**-Anforderung
(„Eine konkrete HTTP-/gRPC-API **muss** die bestehenden Lese- und
Verwaltungsfähigkeiten … über einen Netzwerkzugriffsweg bereitstellen"),
nicht mehr das weiche `LH-FA-SST-005` („**soll** … ermöglichen"). Das ist
eine Anforderung im Lastenheft-Change im wörtlichen Sinn des Triggers.
**Der Trigger ist eingetreten — Re-Evaluierung ist fällig und wird hiermit
durchgeführt.**

**Was die Re-Evaluierung ergibt: `ADR-0020` bleibt unverändert, kein
Folge-ADR.** Der scheinbar naheliegende Schluss „Trigger eingetreten →
Folge-ADR mit `supersedes`" trägt hier nicht, aus drei Gründen:

1. **`ADR-0020`s Entscheidungssatz ist bereits mit diesem Fall
   kompatibel.** Die Entscheidung lautet „HTTP/gRPC können **später** als
   zusätzliche Driving Adapters ergänzt werden; sie gehören **nicht
   zwingend** zum MVP" — das ist keine Verneinung von HTTP/gRPC, sondern
   eine Terminierungs-/Priorisierungs-Entscheidung. `LH-FA-SST-006`
   fordert, dass die API **irgendwann** entsteht („muss … bereitstellen"),
   sagt aber nichts über MVP-Zugehörigkeit oder Zeitpunkt — im Gegenteil,
   `LH-FA-SST-006`s eigenes Out-of-Scope verweist Protokoll-/Endpunktwahl
   ausdrücklich an eine **künftige** ADR/Spezifikation. Es gibt keinen
   Widerspruch zwischen „muss irgendwann als zusätzlicher Adapter
   existieren" (`LH-FA-SST-006`) und „nicht zwingend zum MVP, als
   zusätzlicher Adapter später ergänzbar" (`ADR-0020`).
2. **`ADR-0020`s eigener Konsequenzen-Abschnitt sagt genau diesen Fall
   voraus** und benennt die richtige Reaktion: „Negativ: ein früh
   auftretender API-Bedarf braucht einen **neuen Entscheidungsanlass**."
   Das ist nicht „diese ADR wird falsch", sondern „eine neue,
   *eigenständige* Entscheidung wird nötig" — und diese neue Entscheidung
   ist die Protokoll-/Wireformat-Wahl (HTTP vs. gRPC, Auth-Verfahren,
   Endpunktschnitt), die `LH-FA-SST-006`s eigenes Out-of-Scope-Feld selbst
   als künftigen ADR-Gegenstand benennt — **nicht** eine Korrektur von
   `ADR-0020`. Eine Folge-ADR mit `Supersedes ADR-0020` würde eine
   Entscheidung ersetzen, die durch `LH-FA-SST-006` gar nicht widerlegt
   wird; die richtige neue ADR (wenn ein Slice die API tatsächlich baut)
   steht **neben** `ADR-0020`, so wie `ADR-0019` (CLI-Wahl) neben
   `ADR-0018`/`ADR-0046` (SQL-Wahl) steht — parallele Anwendungen
   desselben von `ARC-005` benannten Optionsraums, kein Sequenz-Verhältnis.
3. **Kein Slice dieser Welle baut die API.** `welle-6.md` §Vermerk selbst
   stellt fest: der Ausgang dieser Prüfung hat „keine Auswirkung auf
   `slice-021`/`slice-022`, die bereits auf CLI festgelegt sind". Es gibt
   damit auch keinen unmittelbaren Umsetzungsdruck, der eine sofortige
   Protokollentscheidung erzwingen würde — anders als bei `slice-021`, wo
   der Architect tatsächlich *anwenden* musste (Zugriffsweg wählen), gibt
   es hier nichts anzuwenden, weil kein Slice den API-Zugriffsweg gerade
   baut.

**Verdikt 1c(i): `ADR-0020` erneut bestätigt, Wortlaut unverändert, kein
Folge-ADR.** Re-Evaluierung fand statt (Trigger-Eintritt korrekt erkannt,
nicht übergangen); ihr Ergebnis ist Bestätigung, nicht Supersession. Für den
Planner als Vormerkung (kein Blocker, analog `architect-review-slice-021.md`
§4): Sobald ein Slice `LH-FA-SST-006` tatsächlich umsetzt, braucht **dieser**
Slice eine eigene, neue Architect-Entscheidung (Protokoll-/Endpunktwahl) —
keine Änderung an `ADR-0020`.

**`LH-FA-SST-007` (NATS-Benachrichtigung, Commit `9f5030d`) — unabhängige
Prüfung, deckt sich mit `verify-slice-022.md` V-2.** Eigene Durchsicht aller
ADR-Dateien (`grep -ril` über `nats|benachrichtig|publish.?subscribe|pub/sub`
in `docs/plan/adr/*.md`) liefert **keinen Treffer** — kein bestehendes ADR
trägt einen Wortlaut, der auf einen Publish-/Subscribe- oder
Benachrichtigungsweg zielt. `ADR-0006` (Replication Stream als Driving
Adapter) betrifft die Capture-Seite (Quelle → CDC-Store), nicht einen
Consumer-seitigen Push-Kanal; `ADR-0025` (Domain Events) ist ein
Domänen-internes Konzept ohne Trigger-Bezug zu einem externen
Transportweg; `ADR-0019`/`ADR-0020` sind kanalgenerisch bzw.
API-Bedarfs-spezifisch, keiner zielt auf Pub/Sub. **Feststellung: kein
ADR-Trigger eingetreten.** Bestätigt unabhängig vom Verifier-Fund, mit
eigener Belegquelle.

**Prozessbeobachtung (kein ADR-Gegenstand):** Beide Lastenheft-CRs
(`9936e82`, `9f5030d`) landeten während laufender `slice-021`/`slice-022`-
Commit-Fenster auf `main`, ohne Bezug zum jeweiligen Slice. Das ist das
zweite und dritte Auftreten desselben Musters „Lastenheft-CR im laufenden
Slice-Diff-Fenster" (`verify-slice-021.md` V-1, `verify-slice-022.md`
V-2) — an sich kein ADR-Trigger-Gegenstand, aber ein Hinweis für den
Planner: Erreicht dieses Muster mit einem vierten Auftreten die
3×-Schwelle als eigene Beobachtungs-Register-Beobachtung (bislang nicht
als `BEO-PGC/*`-Eintrag geführt, weil es kein Slice-internes, sondern ein
Roadmap-/Lastenheft-Pflege-Muster ist), wäre ein eigener Registereintrag
angezeigt. Das ist keine Entscheidung dieses Laufs — nur eine Weitergabe
der bereits im Verifier-Bericht angelegten Beobachtung.

**Sonstige berührte ADRs.** `ADR-0019` (CLI-Driving-Adapter, `permanent`):
beide Slices wenden sie kanalgenerisch an (bereits durch
`architect-review-slice-021.md` §1 hergeleitet); ihr Trigger ist
`permanent` und nicht ereignisgebunden — **bestätigt, unverändert.**
`ADR-0046` (SQL-Rollenspaltung, `permanent`, Trigger „technische
FDW-/`dblink`-Brücke eingeführt"): Keine dieser Welle zugehörige Änderung
führt eine solche Brücke ein (beide Slices sind CLI-Zugriffswege, kein
SQL-Artefakt geändert, wie `verify-slice-022.md` „Entscheidungs-Konformität"
bereits real geprüft hat) — **bestätigt, unverändert.**

---

## Zug 2 — Beobachtungs-Register, Lese-Schritt (`BEO-PGC/rollen-verdrahtung`)

### Zähler und Belege

Drei `evidence/`-Dateien: `slice-011.md` (Welle 3, Erstfund — DDL-Rollen
existieren, Verdrahtung nutzt sie nicht), `slice-021.md` (`welle-6`,
CLI-Registrierungsweg reproduziert dasselbe Muster), `slice-022.md`
(`welle-6`, CLI-Bestätigungsweg reproduziert es ein drittes Mal). Eigene
Nachzählung: `find
docs/plan/planning/observations/BEO-PGC/rollen-verdrahtung/evidence/`
liefert genau drei Dateien — Schwelle real erreicht, nicht nur behauptet.
`state.md` trägt korrekt „Zähler (abgeleitet): 3×" mit Ausgang bislang
`weiter offen` — dem Lese-Schritt dieser Welle-Closure zugeordnet, exakt
der Punkt, an dem dieser Zug greift (`verify-slice-022.md` V-1 hatte diese
Zuordnung bereits empfohlen, nicht selbst vorweggenommen).

### Prüfung: trägt „verkörpert" (geschärfte Regel statt Fix)?

Geprüft und **verneint** — anders als bei `BEO-PGC/dod-checkbox-nachzug`
(welle-5, verkörpert) liegt hier keine Workflow-Disziplin-Lücke vor,
sondern eine **physische Verdrahtungslücke im Code**:

- **Die Ursache ist kein wiederholtes Vergessen einer Handlung, sondern
  ein nicht gebauter Konfigurationspfad.** `internal/bootstrap/wiring.go`
  hat schlicht keinen Mechanismus, verschiedenen Adaptern verschiedene
  Rollen-DSNs zu geben — jeder neue Aufrufer (Capture, Store,
  `RegisterConsumer`, `AcknowledgeConsumer`) bekommt zwangsläufig `cfg.DSN`,
  weil es keine zweite Konfigurationsquelle gibt. Eine Workflow-Zeile
  („denk daran, die Rolle zu binden") ändert diesen Zustand nicht — der
  nächste Implementer hätte weiterhin nur `cfg.DSN` zur Verfügung und
  könnte die Regel gar nicht befolgen, ohne selbst genau die
  Bootstrap-Änderung vorzunehmen, die hier ansteht.
- **Der Aufwand ist slice-groß, nicht zeilen-groß.** Wie in der
  Aufgabenstellung zutreffend vorgezeichnet: Eine rollen-spezifische
  DSN-Verdrahtung berührt `internal/bootstrap/wiring.go` **und** jeden dort
  gebauten Adapter-Konstruktor zugleich (Capture-Stream, Store-Adapter,
  beide CLI-Wege, potenziell künftige Wege) — mehrere Komponenten, ein
  neuer Konfigurationsvertrag (`compose.yaml`, Betreiberdoku), eigene
  Tests gegen echte PostgreSQL-Rollen. Das ist exakt die Grössenordnung,
  für die der Harness Slices vorsieht, keine Ein-Zeilen-Workflow-Ergänzung.
- **Eine geschärfte Regel würde das bestehende Defizit nicht auflösen.**
  Alle drei Ausgänge (`slice-011`, `slice-021`, `slice-022`) bestehen
  weiterhin unverändert fort, auch nach einer Regel für *künftige*
  Zugriffswege — das Muster wäre „begrenzt", nicht „behoben". `Modul 6`
  verlangt an der 3×-Schwelle einen Ausgang, der die *bereits* eingetretene
  Wiederholung trägt, nicht nur ihr Fortschreiten verlangsamt.

### Prüfung: trägt „gestrichen"?

Nein — offensichtlich verneint, wie in der Aufgabenstellung selbst
vorgezeichnet: Das zugrunde liegende Verdrahtungsmuster (gemeinsame
Instanz-DSN) besteht unverändert fort; nichts an dieser Welle hat es
aufgelöst oder unmöglich gemacht. Eine Streichung „mit Begründung, warum
die Beobachtung nicht mehr auftreten kann" wäre eine Falschaussage.

### Prüfung: trägt „geplant"?

**Ja — das ist der tragfähige Ausgang.** Modul 6 verlangt für `geplant`
eine „Kennung eines Slice oder der Welle, die sie schreibt" — nicht nur
einen Vorsatz. Geprüfte Kriterien:

- **Ist die Arbeit real umsetzbar, ohne eine noch offene Architektur-Frage
  vorauszusetzen?** Ja — die drei Rollen existieren bereits als DDL
  (`slice-011`), sind einzeln per `SET ROLE` testbar; es fehlt „nur" die
  Bootstrap-Verdrahtung und der erweiterte Konfigurationsvertrag. Keine
  neue ADR nötig, um zu beginnen.
- **Ist der Umfang slice-groß (≤ 3 Liefer-Punkte, ≤ 2 Schichten) oder
  braucht es eine Welle?** Voraussichtlich ein einzelner, aber
  substantieller Slice: Bootstrap-Schicht + Konfigurations-/Doku-Schicht
  (zwei Schichten), Liefer-Punkte grob: (1) Konfigurationsvertrag
  erweitert, (2) alle Adapter-Konstruktoren auf Rollen-DSN umgestellt und
  real getestet, (3) Betreiberdoku/`compose.yaml` nachgezogen — genau an
  der Obergrenze von drei; sollte sich beim Ausplanen mehr ergeben
  (z. B. weil ein Migrationspfad für Bestandskonfiguration nötig wird),
  ist die im Slice-Entwurf selbst benannte Rückführung `in-progress→next`
  die vorgesehene Antwort, keine stillschweigende Ausweitung.
- **Existiert eine belastbare Kennung?** Ja —
  [`slice-023`](../planning/done/slice-023-rollen-spezifische-dsn-verdrahtung.md),
  von diesem Lauf als Skelett in `open/` angelegt (§1 Ziel/Abgrenzung,
  §3 Datei-Kandidaten, §4/§5 Trigger, §8 Sub-Area-Vorprüfung bereits
  gefüllt; §2/§6/§7 bewusst als Platzhalter markiert — Planner
  vervollständigt vor dem Übergang nach `next/`, das ist Planungsarbeit,
  keine Architektur-Entscheidung). Diese Datei macht die Register-Paarung
  (b) bei der nächsten Closure, die sie zitiert, prüfbar: die Kennung
  löst im Planning-Lifecycle auf.

**Verdikt Zug 2: Ausgang `geplant` → `slice-023`.** Herkunfts-Anker für die
Register-`state.md`-Fortschreibung (Planner-Arbeit): `seit welle-6`.

---

## Disposition (Gesamt)

| Frage | Ergebnis |
|---|---|
| Aktiver Carveout mit Bezug zu `welle-6`? | Nein — 0 aktiv im gesamten Repo |
| Reifestufe (Bootstrap-aware Gate) hochzuschalten? | Nein — durchgehend Greenfield, keine Stufe vorhanden |
| `ADR-0020`-Trigger gegen `LH-FA-SST-006` eingetreten? | **Ja**, wörtlich — Re-Evaluierung durchgeführt, Verdikt: **bestätigt**, kein Folge-ADR (Entscheidungssatz widerspruchsfrei mit `LH-FA-SST-006`; künftige Protokollwahl ist eine neue, parallele ADR, keine Supersession) |
| `LH-FA-SST-007` löst einen bestehenden ADR-Trigger aus? | Nein — eigenständig geprüft (Volltextsuche über alle ADR-Dateien), deckungsgleich mit `verify-slice-022.md` V-2 |
| `ADR-0019`/`ADR-0046` fällig? | Nein — beide `permanent`, in dieser Welle nur angewendet |
| `BEO-PGC/rollen-verdrahtung` (3×) — Ausgang | **geplant** → [`slice-023`](../planning/done/slice-023-rollen-spezifische-dsn-verdrahtung.md) (Skelett angelegt) |
| Warum nicht „verkörpert"? | Physische Verdrahtungslücke im Bootstrap-Code, kein Workflow-Disziplin-Defizit; eine Regel löst das bestehende, dreifach reproduzierte Defizit nicht auf |
| Auswirkung auf die übrige `welle-6`-Closure | Keine Blocker — Schritte 1, 3c–6 bleiben Planner-Arbeit |

---

## Beleg-Anker (Kurzfassung)

| Aussage | Beleg |
|---|---|
| 0 aktive Carveouts | `docs/plan/carveouts/` (nur `.gitkeep`) |
| Durchgehend Greenfield, keine Reifestufe | `harness/conventions.md` §Modus-Deklaration; `slice-021`/`slice-022` §8 |
| [`LH-FA-SST-006`](../../../spec/lastenheft.md) ist eine `muss`-Anforderung, Version 0.3.0→0.4.0 | `spec/lastenheft.md` §Versionshistorie 0.4.0 |
| `ADR-0020`-Trigger-Wortlaut, Konsequenzen-Vorwegnahme | [`ADR-0020`](0020-http-grpc-optional.md) §Re-Evaluierungs-Trigger, §Konsequenzen |
| Kein ADR mit Pub/Sub-/Benachrichtigungs-Trigger | eigene Volltextsuche `docs/plan/adr/*.md` (keine Treffer für `nats\|benachrichtig\|publish.?subscribe\|pub/sub`) |
| `ADR-0019`/`ADR-0046` `permanent`, unberührt | jeweiliges ADR §Re-Evaluierungs-Trigger; `verify-slice-022.md` §Entscheidungs-Konformität |
| Zähler `rollen-verdrahtung` = 3× | `evidence/{slice-011,slice-021,slice-022}.md`, `state.md` |
| Least-Privilege ist Soll-Anforderung, bereits als DDL+Test erfüllt, Laufzeit-Verdrahtung offen | [`LH-QA-SEC-001`](../../../spec/lastenheft.md)…003; `observation.md` |
| Folge-Slice-Skelett angelegt | [`slice-023`](../planning/done/slice-023-rollen-spezifische-dsn-verdrahtung.md) |
| Präzedenzfall „verkörpert" zum Vergleich (Workflow-Disziplin, nicht Architektur-Lücke) | `architect-review-welle-5.md` Zug 2 |
