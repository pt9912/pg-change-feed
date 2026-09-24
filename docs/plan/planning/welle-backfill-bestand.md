# Welle backfill-bestand: Backfill des Bestands — der Tabellenbestand einer aktivierten Tabelle wird als Backfill erkennbar über den bestehenden Lesezugriffsweg lesbar (`LH-FA-CAP-009`, `ADR-0111`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-backfill-bestand-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (Rolleninhaber der Implementer-Rolle je Slice, gesetzt
beim Übergang `open` → `next`; geschnitten aus
[`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
§Folgepflichten 2 durch den Planner). **Datum:** 2026-09-23.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der Bestand einer aktivierten Tabelle wird auf ausdrückliche Auslösung hin
über **denselben Lesezugriffsweg** lesbar, den auch WAL-erfasste Changes
tragen — erkennbar als Backfill (`origin = 'backfill'`), lückenlos an den
laufenden WAL-Pfad angeschlossen und ohne Verlust bei Unterbrechung. Das ist
die Aussage der drei Akzeptanzkriterien von
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Happy Path, Boundary,
Negative); der Mechanismus ist mit
[`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) entschieden:
eine Bulk-Copy des Tabellenbestands in dem `REPEATABLE READ`-Snapshot eines je
Run angelegten temporären logischen Slots, in **einer** Store-Transaktion
committet, als `INSERT` mit dem neuen Feld `origin`, ausgelöst über die
SQL-Funktion `cdc.backfill_table`.

Das *Mehr* gegenüber den zehn Slice-DoDs: keiner der Slices belegt eines der
drei Kriterien allein. Der Happy Path braucht das Feld `origin` (Store, View,
`GET /changes`), den Snapshot-Leser, den Run mit atomarem Schreiber und die
SQL-Auslösung am **laufenden Feed-Container**; die Boundary (Überlappungsfenster
zum WAL-Pfad: weder stille Lücke noch unbegrenzte Dopplung) liegt in der
Zusammenschaltung von Slot-Position `X` und WAL-Stream und ist erst am
komponierten System prüfbar; die Negative (Unterbrechung, Neubeginn ohne
Verlust) braucht Start-Abgleich, Run-Zustand und den realen Prozessabbruch.
Dazu kommt der Träger-Zustand, den kein Einzel-DoD beobachtet: die
Row-Image-Konstruktion liegt an genau **einer** Stelle
([`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) §Fitness
Function, Review-Prüfpflicht), und [`LH-FA-CAP-009`](../../../spec/lastenheft.md)
ist im RTM-Lauf (`make doc-trace`) nicht mehr Waise — beides ist ein Zustand
des ganzen Bündels, zu belegen bei der Welle-Closure.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) trägt den
  Status `Accepted` — **bereits erfüllt**, geprüft an der Zeile
  `ADR-0111` im ADR-Index
  ([`docs/plan/adr/README.md`](../adr/README.md), Status-Spalte
  `Accepted`, Datum 2026-09-23).
- [`LH-FA-CAP-009`](../../../spec/lastenheft.md) steht im Lastenheft als
  Anforderung mit den drei Akzeptanzkriterien — **bereits erfüllt**, geprüft
  an `spec/lastenheft.md` §`LH-FA-CAP-009`.
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden; die
  Reihenfolge der Slices steht in §4.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle zehn Slices in `done/`.
- `make gates` grün — der Exit-Code des Laufs wird ungefiltert gesichert und
  gesondert ausgewertet ([`AGENTS.md`](../../../AGENTS.md) §3.9).
- **Ein realer, grüner `make test-integration`-Lauf** mit den Backfill-Belegen
  am laufenden Feed-Container — Happy Path (Bestand über `cdc.changes` und
  `GET /changes`, `origin = 'backfill'`), Boundary (nebenläufige Schreiber
  während des Runs samt Replay-Invariante; leere Tabelle als Randfall) und
  Negative (`docker kill` mitten im Run, `interrupted`, erneuter Antrag) —;
  kein Gate, aber Pflichtbeleg dieser Welle.
- **Reale, grüne Läufe der DB-Tiers:** `make test-replication` (Snapshot-Träger,
  Bild-Parität) und `make test-store` (Run-Zustand,
  Annahme-Transaktion, atomarer Schreiber, Antrags-Funktion, Rollen), und `make schema-rollout` zweimal
  hintereinander gegen dieselbe Ziel-Datenbank (Idempotenz) — die
  Fitness-Function-Zeilen von
  [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md), die nur
  das gebündelte System zeigt.
- `make doc-trace` führt [`LH-FA-CAP-009`](../../../spec/lastenheft.md) nicht
  mehr unter den Waisen; der Träger ist die Zeile in
  [`docs/user/e2e-abdeckung.md`](../../user/e2e-abdeckung.md), die der Runner
  von `make test-integration` schreibt (die Datei ist ein Erzeugnis, kein
  Lauf-Beleg).
- Die Row-Image-Konstruktion liegt an genau einer Stelle: der Suchlauf über
  `internal/**` nach der JSON-Bild-Erzeugung (Befehl und Fundstellen im
  Closure-Bericht) findet eine Konstruktionsstelle, die WAL-Pfad und
  Backfill-Pfad gemeinsam aufrufen.
- Die Sichtbarkeits-Grenze (Backfill ist ein Zustandsabzug für Consumer vor
  `X`), die Lese-Regel „Bestandsabzug ohne `Limit`", die gemessene
  Startposition eines frisch registrierten Consumers, die gemessene
  Warn-Richtgröße und die Toleranz als „Startwert, Setzung ohne Messung" stehen im
  Benutzerhandbuch — jede Zahl mit ihrem Ursprung
  ([`AGENTS.md`](../../../AGENTS.md) §3.12).
- Closure-Notiz in `welle-backfill-bestand-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-backfill-spec-nachzug | Pflichtenheft (`LH-FA-CAP-009.a` beantwortet, `SPEC-002`/`SPEC-019`/`SPEC-022`, neue `SPEC-029`) und Architektur-Sicht (Backfill-Sequenz, ohne ADR-/Slice-Bezug) auf den beschlossenen Stand ziehen | [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Folgepflicht 1 |
| slice-backfill-row-image-gemeinsam | Row-Image-Konstruktion als eine gemeinsame, reine Funktion der Domäne; der WAL-Pfad ruft sie, byte-gleiches Ergebnis | [`LH-FA-CAP-008`](../../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2 |
| slice-backfill-change-origin | Feld `origin` (`wal` \| `backfill`, `NULL` ≙ `wal`) in Domäne, Store, View `cdc.changes` und `GET /changes` | [`LH-FA-DAT-006`](../../../spec/lastenheft.md), [`LH-FA-REA-001`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2/8 |
| slice-backfill-snapshot-reader | Outbound Port `TableSnapshotPort` und Driven Adapter: temporärer Slot mit Export-Snapshot, Import, Cursor-Blöcke, Spaltenliste, GUC-Parität | [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`LH-FA-CAP-008`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 1 |
| slice-backfill-run-usecase | Domäne `BackfillRun`, `BackfillTableUseCase`, Fähigkeits-Ports für Annahme, Run-Zustand und atomaren Schreiber, Fail-closed-Prüfung, Wecksignal — gegen Fakes | [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 3/4/6 |
| slice-backfill-run-store | Tabelle `cdc.backfill_run`, Grants, Postgres-Adapter für Annahme, Run-Zustand und den einen atomaren Schreiber (Ordnung, `change_id`, `origin`) | [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`LH-FA-CAP-004`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4/6 |
| slice-backfill-sql-administration | Antragsart `backfill`, `cdc.backfill_table`, Worker-Schleife mit Start-Aufnahme und Wecksignal, Start-Abgleich, `cdc.backfill_status`, `diagnose`, Idempotenz-Guard, Handbuch der Auslösung | [`LH-FA-ADM-001`](../../../spec/lastenheft.md), [`LH-FA-SST-003`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5 |
| slice-backfill-e2e | `make test-integration`: Happy Path, Boundary, Negative; Startposition eines frisch registrierten Consumers gemessen | [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1 |
| slice-backfill-bench-richtgroesse | Bench der Kopierdauer je Tabellengröße, daraus die Warn-Richtgröße; Auswertung der zwei Warnungen im Use Case des Runs, sichtbar über View und `diagnose` | [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3, [`ADR-0113`](../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3, [`ADR-0054`](../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) |
| slice-backfill-sdk-origin | `origin` in den drei SDK-HTTP-Lesemodellen; Package-Versionen | [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8 |

**Reihenfolge:** sequentiell in der Tabellen-Reihenfolge (WIP-Limit 1 je
Rolleninhaber, Baseline-Regelwerk `modul-05-planning-harness.md`), außer `slice-backfill-sdk-origin`, der nach `slice-backfill-change-origin`
technisch unabhängig ist und aus Ordnungsgründen am Ende steht. Die
technischen Kanten stehen in §5.

**Abweichungen vom Schnitt-Vorschlag** in
[`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
§Folgepflichten 2 (dort ausdrücklich „Planner entscheidet final") — vier, je
mit Grund:

1. **`S1` geteilt** in `row-image-gemeinsam` und `change-origin`. Die
   Verschiebung der Bild-Konstruktion ist ein Refactor mit Byte-Gleichheits-
   Pflicht am WAL-Pfad — ein Fehler dort verändert **jede** erfasste Change —,
   `origin` ist eine additive Erweiterung von Schema, View und Lesepfaden;
   beide zusammen sprengen den ≤3-Liefer-Punkte-Rahmen eines einzelnen
   Reviews. Außerdem ist die gemeinsame Funktion der Punkt, an dem die
   Transformations-Umsetzung ([`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md))
   ansetzt (§5): ein eigener Slice macht diese Kante benennbar.
2. **`S3` geteilt** in `run-usecase` (Domäne, Use Case, Ports, Fakes — netzlos)
   und `run-store` (Tabelle, Grants, Postgres-Adapter — Store-Tier). Ein Slice
   mit Domäne, Use Case, zwei Ports, zwei Adaptern, Schema und Grants trägt
   mehr als drei Liefer-Punkte.
3. **`S5` geteilt** in `e2e` und `bench-richtgroesse`. Der Bench braucht ein
   lauffähiges System und liefert eine Zahl, aus der Produktions-Code
   (die Warnschwelle) entsteht — ein anderer Liefer-Gegenstand als die
   E2E-Belege; die Zahl darf nicht vor ihrer Messung im Code stehen
   (`BEO-PGC/geschaetzter-wert-als-grenze`).
4. **Handbuch-Zug verteilt** statt in `S0` gebündelt: `S0` trägt die
   führenden Spec-Straten (Pflichtenheft, Architektur-Sicht — Doku führt, die
   Modus-Deklaration in `harness/conventions.md` ist Greenfield); das
   Benutzerhandbuch beschreibt Betreiber-Oberfläche und darf keine Funktion
   nennen, die es noch nicht gibt. Jeder Slice, der eine Betreiber-Oberfläche
   liefert, zieht den Handbuch-Abschnitt seiner Oberfläche selbst nach
   (`change-origin`: Lesen; `sql-administration`: Auslösung, Sichtbarkeits-
   Grenze, `SELECT`-Grant, Lesen ohne `Limit`; `e2e`: gemessene Startposition;
   `bench-richtgroesse`: Warn-Richtgröße).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Wird blockiert von:** keiner Welle;
  [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) ist
  `Accepted`, der Trigger (§2) ist erfüllt. Den Status `Accepted` von
  [`ADR-0113`](../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  verlangen zwei Slice-Start-Trigger (`slice-backfill-run-usecase`,
  `slice-backfill-bench-richtgroesse`), kein Wellen-Start-Trigger: die Eröffnung
  hängt nicht daran.
- **Blockiert:** Slices der Welle
  [welle-transformationen](welle-transformationen.md) — die Umsetzung der
  Transformationen
  ([`LH-FA-CFG-007`](../../../spec/lastenheft.md),
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md))
  hängt über **benannte Kopplungen** an dieser Welle; die Bedingungen sind
  Start-Trigger einzelner Slices der Transformations-Welle:
  - **K1 — Row-Image-Funktion.** Der Kern der Transformations-Umsetzung
    ([`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
    Folgepflicht 2, Auswertung „in `rowImage` nach dem Ausschluss")
    startet erst, wenn `slice-backfill-row-image-gemeinsam` in `done/` liegt:
    die Regelauswertung hängt an der **einen** gemeinsamen Funktion, nicht an
    zwei Bild-Erzeugern.
  - **K2 — Backfill-Pfad.**
    [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
    Folgepflicht 7 („Ein Backfill-Pfad … trägt dieselbe Regelauswertung; die
    zugehörige ADR nennt sie") wird als eigener Slice der Transformations-
    Umsetzung geführt, nach `slice-backfill-run-usecase`; er trägt (a) den
    Regelstand im Run (Auswertung im Bild-Bau), (b) die Erweiterung der
    Fail-closed-Prüfung um den Regelstand und (c) einen E2E-Beleg (Backfill-
    Change trägt die transformierte Form). [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
    nennt Transformationen nicht; die Bindung trägt allein
    [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md).
  - **K3 — gemeinsame Zeilen.** Beide Umsetzungen berühren dieselben Stellen:
    die geschlossene `request_kind`-Menge (`tools/schema/nacharbeit-administration.sql`),
    [`SPEC-019`](../../../spec/pflichtenheft.md), `applyAdministrationRequest`
    und den Idempotenz-Guard. Der Antragsweg-Slice der Transformationen
    (Folgepflicht 3) startet nach `slice-backfill-sql-administration` in
    `done/` und erweitert die dort stehende Menge additiv.
  - **Begründete Reihenfolge — Backfill zuerst.** Beide Umsetzungen sind ohne
    einander startbar; die Kanten K1–K3 zeigen alle in dieselbe Richtung: die
    Transformationen erweitern Stellen, die der Backfill anlegt, nicht
    umgekehrt. Käme die Transformations-Umsetzung zuerst, müsste der Backfill
    den Regelstand schon im ersten Entwurf des Run-Use-Case und der
    Fail-closed-Prüfung tragen (mehr Umfang in `slice-backfill-run-usecase`
    und ein Entwurf gegen einen noch ungebauten Regelstand-Port) und die
    `request_kind`-Menge würde zweimal umgeschrieben.
- **Intern:** die Ordnung folgt der Tabelle in §4; die technischen Kanten
  (`A → B` heißt: `B` setzt `A` voraus) sind:
  - `spec-nachzug` → jeder übrige Slice (die Spec führt);
  - `row-image-gemeinsam` → `snapshot-reader` (die Bild-Parität vergleicht den
    Backfill-Pfad mit dem WAL-Pfad über dieselbe Funktion) und → `run-usecase`
    (der Run baut das Bild über sie);
  - `snapshot-reader` → `run-usecase` (der Run nutzt `TableSnapshotPort`);
  - `run-usecase` → `run-store` (die Adapter implementieren die dort
    festgelegten Fähigkeits-Ports);
  - `change-origin` → `run-store` (der atomare Schreiber setzt
    `origin = 'backfill'`);
  - `run-store` → `sql-administration` → `e2e` → `bench-richtgroesse`
    (jeder Schritt braucht das lauffähige System der Vorstufe);
  - `change-origin` → `sdk-origin` (sonst unabhängig von den übrigen Slices).

**Träger der Folgepflichten** — jede Pflicht aus
[`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
§Folgepflichten hat einen Träger oder eine benannte Adresse
(`BEO-PGC/adr-folgepflicht-ohne-traeger-slice` — eine adresslose Pflicht
bleibt liegen, bis sie ein Leser zufällig findet):

| Folgepflicht | Träger |
|---|---|
| 1 Spec-Nachzug | `slice-backfill-spec-nachzug`; der Benutzerhandbuch-Anteil verteilt auf `change-origin`, `sql-administration`, `e2e`, `bench-richtgroesse` (§4, Abweichung 4) |
| 2 Umsetzung (S1–S6) | die zehn Slices dieser Welle |
| 3 Nicht Teil (Checkpoint, Parallelisierung, HTTP-/CLI-Auslösung, Live-Zustellung, Filter `origin`, Auslösung als Option von `enable`) | §6 dieser Welle; Adresse: die Re-Evaluierungs-Trigger in [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) |
| 4 Kommentar-Nachzug `model.Change` | `slice-backfill-change-origin` (DoD) |
| 5 Idempotenz-Guard | `slice-backfill-sql-administration` (DoD) |
| 6 Cursor-Form (Priorität nach `S5`) | kein Slice dieser Welle; Adresse: `BEO-PGC/limit-fortsetzung-innerhalb-einer-position` (offen) und der Re-Evaluierungs-Trigger von [`ADR-0081`](../adr/0081-changes-lesen-ueber-die-http-api.md); bis dahin trägt die Doku die Lese-Regel „Bestandsabzug ohne `Limit`" (`slice-backfill-sql-administration`) |

**Träger der Folgepflichten von
[`ADR-0113`](../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)**
(Rollenschnitt, Aufnahme, Warn-Kriterium):

| Folgepflicht | Träger |
|---|---|
| 1 Spec-Nachzug (Grants, zwei Warn-Spalten, `applied` = „angenommen", Sicht) | `slice-backfill-spec-nachzug` |
| 2 Port-Schnitt der Annahme, `Request` endet mit `Admit`, Vorbedingungen erneut prüfen | `slice-backfill-run-usecase` (Start-Trigger: `ADR-0113` `Accepted`) |
| 3 Grants, `estimated_rows` nullable, zwei Warn-Spalten, Annahme-Adapter, Rollen-Tests | `slice-backfill-run-store` |
| 4 Verarbeitungs-Zweig, Worker-Schleife mit Start-Aufnahme und Signal, View und `diagnose` | `slice-backfill-sql-administration` |
| 5 Warn-Auswertung, Toleranz-Konstante, Handbuch | `slice-backfill-bench-richtgroesse` (Start-Trigger: `ADR-0113` `Accepted`) |
| 6 `reltuples` = `−1` an PostgreSQL 17 und 18 | `slice-backfill-snapshot-reader` |
| 7 Neustart-Beleg der `queued`-Zeile | `slice-backfill-e2e` |

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Fortsetzen mit Checkpoint und Block-Transaktionen** — die Ein-Transaktions-
  Form gilt für den Erstumfang ([`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 4 Option B, Festlegung 3); der Re-Evaluierungs-Trigger (Kopierdauer
  über der Betriebs-Toleranz) führt zu einer Folge-ADR.
- **Parallelisierung** — Ausbaustufe laut
  [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (ein Worker, ein Run zugleich).
- **HTTP-Endpunkt und CLI-Auslösung** — im Erstumfang genügt
  `cdc.backfill_table` samt `cdc.backfill_status` und `diagnose`
  ([`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Festlegung 2); die Auslösung ohne SQL-Zugang ist additiv über denselben
  Use Case und ein eigener Slice.
- **Live-Zustellung von Backfill-Changes** über gRPC, SSE und NATS-Vollinhalt —
  die Live-Wege tragen ausschließlich WAL-Changes, ihre Schemata
  ([`SPEC-020`](../../../spec/pflichtenheft.md), [`SPEC-021`](../../../spec/pflichtenheft.md),
  [`SPEC-024`](../../../spec/pflichtenheft.md)) und die Proto-Artefakte bleiben
  unberührt, `make generated-sync` ist nicht betroffen; das Wecksignal nach dem
  Run-Commit läuft über den bestehenden `ChangeNotificationPort`.
- **Filter `origin` an `GET /changes`** — das Feld ist lesbar, die
  serverseitige Filterung eine Folge-Entscheidung (Filter-Grammatik nach
  [`ADR-0081`](../adr/0081-changes-lesen-ueber-die-http-api.md)).
- **Consumer-Reset** — das Positionsmodell bleibt unverändert; der Reset ist
  „explizit administrativ" ([`ADR-0013`](../adr/0013-consumer-domainkonzept.md))
  und eine eigene Fähigkeit mit eigener Entscheidung
  ([`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Festlegung 1).
- **Cursor-Form für `Limit` innerhalb einer Position** — Re-Evaluierung von
  [`ADR-0081`](../adr/0081-changes-lesen-ueber-die-http-api.md); diese Welle
  trägt nur die Lese-Regel in der Doku. Hinweis für eine spätere Welle: alle
  Blöcke eines Runs teilen die Position `X`, ein Bestand jenseits von `Limit`
  Zeilen ist über `GET /changes` nur ohne `Limit` lesbar.
- **Retention-Paging** — `RunRetentionService` liest alle Changes einer Quelle
  in den Speicher; die Grenze besteht unabhängig vom Backfill
  (Re-Evaluierungs-Trigger von
  [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md)).
- **Instanz-zu-Instanz-Migration** — Out-of-Scope im Lastenheft.
- **Transformationen und Routing**
  ([`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`LH-FA-CFG-008`](../../../spec/lastenheft.md),
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)) —
  die Transformationen sind Gegenstand von
  [welle-transformationen](welle-transformationen.md), das Routing hat weder
  Welle noch Umsetzung; die Kopplung steht in §5.
- **Kein neuer GitHub-Actions-Workflow und keine strukturelle Workflow-
  Änderung** — [`AGENTS.md`](../../../AGENTS.md) §3.10 greift nicht. Die
  Package-Versionen der drei SDKs heben im Slice `sdk-origin` als Datei-
  Änderung; ein Tag-Push und damit ein Release bleibt Betreiber-Handlung
  außerhalb dieser Welle. Die Pipeline `e2e.yml` fährt `make test-integration`
  unverändert weiter — die Laufzeit des erweiterten Testpakets ist ein
  Risiko des Slice `e2e` (§6 dort), keine Workflow-Änderung.
- **Kein Server-Release** — `docs/user/version.md` bleibt unberührt; die
  Änderungshistorie des Benutzerhandbuchs trägt je berührtem Slice eine Zeile.
- **Keine Schwellen-Senkung und keine neue Gate-Klasse** — das neue DB-Paket
  aus `slice-backfill-snapshot-reader` (`postgressnapshot`) ist nach
  [`ADR-0071`](../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 1 nicht Gegenstand des Coverage-Gates, weil sein Testlauf einen externen
  Dienst voraussetzt; die Regel greift ohne Textänderung, der Slice zieht die
  namentlichen Listen nach (Trigger (a)); eine Senkung bliebe per
  [`AGENTS.md`](../../../AGENTS.md) §3.6 ADR-pflichtig.
- **Lastenheft unverändert** — die Welle schärft Techniken im Pflichtenheft,
  nicht Anforderungen.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt
2):** Das Register (`docs/plan/planning/observations/README.md`) wurde
durchgesehen. Die Zähler sind die Zahl der Dateien unter `evidence/` (real
ausgezählt am 2026-09-23). **Kein Eintrag erreicht mit der Eröffnung dieser
Welle 3×** — kein Slice ist allein wegen einer Beobachtung nötig. Die Einträge
bei 2× (`vorab-bedingung-nach-umsetzung-geprueft`, `nachzug-laesst-ueberholten-text-stehen`,
`test-runner-stiller-ausschluss`, `db-gegenstand-enthaelt-netzlos-geprueften-code`,
`coverage-stage-dockerignore-blockiert-tooling`, `rollen-test-abdeckungsluecken`)
erreichen die Schwelle mit dem nächsten Beleg; der Lese-Schritt der
Welle-Closure liest sie. Treffer je Sub-Area, mit dem Slice, der sie als
Risiko trägt (Detail je Slice in §8):

- `BEO-PGC/limit-fortsetzung-innerhalb-einer-position` (offen, 0×, „benannt,
  nicht gezählt") — **Bezug** dieser Welle: der Backfill macht die Klasse zur
  Regel (jeder Run ist eine Position mit vielen Changes). Träger in der
  Doku: `slice-backfill-sql-administration`; Cursor-Form: §6.
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert,
  [`AGENTS.md`](../../../AGENTS.md) §3.13, 26×) — betrifft jeden Slice; der
  Suchlauf steht als committiertes Feld je Slice-Plan §3.
- `BEO-PGC/adr-folgepflicht-ohne-traeger-slice` (offen, 1×) — betrifft diese
  Welle als Ganzes: die Träger-Tabelle in §5 gibt jeder Folgepflicht von
  [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) einen
  Träger oder eine benannte Adresse; die Kopplung zu
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 7 steht in §5.
- `BEO-PGC/geschaetzter-wert-als-grenze` (offen, 1×) — `slice-backfill-sql-administration`
  und `slice-backfill-bench-richtgroesse`: die geschätzte Zeilenzahl
  (`pg_class.reltuples`) und die Warn-Richtgröße dürfen an keinem Träger zur
  Grenze werden; die Richtgröße entsteht erst aus einer Messung.
- `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×) —
  `slice-backfill-run-usecase` und `slice-backfill-bench-richtgroesse`: der Status
  `Accepted` von
  [`ADR-0113`](../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  steht als Start-Trigger, nicht als Rückführungs-Bedingung; die Zuordnung des
  DB-Pakets im `snapshot-reader` ist vor dem Start entschieden.
- `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 2×),
  `BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert),
  `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (offen, 2×),
  `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
  `slice-backfill-snapshot-reader`, `slice-backfill-run-usecase`,
  `slice-backfill-run-store`: ein neues Paket verschiebt Nenner und Zähler
  der beiden Coverage-Messungen; die Build-Kontext-Prüfung für neue Pfade ist
  je Slice ein DoD-Punkt.
- `BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×),
  `BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×),
  `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (verkörpert) —
  `slice-backfill-run-store`, `slice-backfill-sql-administration`: neue
  Tabelle, CHECK-Menge, Funktion, Guard-Eintrag; der `queued`-Run trägt einen
  dauerhaften Träger.
- `BEO-PGC/lese-doppelquelle` (verkörpert, 3×) — `slice-backfill-change-origin`:
  View und `ReadChanges` tragen dieselbe Spaltenmenge.
- `BEO-PGC/rollen-verdrahtung` (eingetreten), `BEO-PGC/rollen-test-abdeckungsluecken`
  (offen, 2×) — `slice-backfill-run-store`, `slice-backfill-sql-administration`:
  neue Grants an `cdc_capture`, `cdc_admin`, `cdc_reader`.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (verkörpert,
  3×), `BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, 3×) —
  jeder Slice, der Betreiber-Oberfläche liefert, zieht Handbuch **und**
  Änderungshistorie nach.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 13×),
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 6×),
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 8×),
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×) —
  Slices mit Messungen und Negativ-Belegen (`e2e`, `bench-richtgroesse`,
  `run-usecase`).
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
  `BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×),
  `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×),
  `BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×) —
  `slice-backfill-e2e`, `slice-backfill-run-store`.
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 2×) —
  `slice-backfill-spec-nachzug` (die Überschrift „… offen" von
  `LH-FA-CAP-009.a` und die Antragsarten-Aufzählungen).
- `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (offen, 3×,
  Ausgang **nicht** zugewiesen — bereits vor dieser Welle bei 3×; kein
  vierter Beleg nötig), `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
  (verkörpert, 3×), `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme`
  (verkörpert, 3×) — `slice-backfill-sdk-origin`.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 7×) —
  **nicht einschlägig**: kein Workflow-Zug in dieser Welle.
- Gesichtet, ohne Bezug zu dieser Welle: die übrigen Einträge des Registers.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-backfill-bestand-results.md`
Zähler: `../observations/` (Beobachtungs-Register)
