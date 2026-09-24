# Architect-Verdikt: Backfill — Schema-Version der Backfill-Changes, Fehlerklasse `schema` im Run, Rollenschnitt der Antrags-Queue

**Rolle:** Architect (Modul 8)

**Anlass:** Drei Übergaben aus den Reviews der Backfill-Slices, die der Reviewer
ausdrücklich als **Entscheidungen** an den Architect gereicht hat (kein
Implementer-Pfeil, keine Fixrunde): die Findings F-8 und F-9 aus dem Review des
Run-Use-Case-Slice (`docs/reviews/review-slice-backfill-run-usecase.md`) und F-1
aus dem Review des Run-Store-Slice samt Verifikation
(`docs/reviews/review-slice-backfill-run-store.md`,
`docs/reviews/verifikation-slice-backfill-run-store.md`). Der Rollenwechsel
Planner → Architect → Planner folgt Modul 8 §Konflikt-Pfad.

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als die
Implementer-Läufe und die Reviewer-Läufe der beiden Slices)

**Datum:** 2026-09-24

**Bezug:** [`LH-FA-SCH-005`](../../spec/lastenheft.md),
[`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`LH-FA-CFG-007`](../../spec/lastenheft.md),
[`LH-QA-SEC-001`](../../spec/lastenheft.md),
[`LH-QA-SEC-002`](../../spec/lastenheft.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md),
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[`SPEC-019`](../../spec/pflichtenheft.md),
[`SPEC-029`](../../spec/pflichtenheft.md); die Slices sind als Kennung genannt,
nicht als Pfad-Link (ein Slice wechselt die Lifecycle-Ablage, ein Pfad-Link
bräche mit)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md)
  (Proposed, `Supersedes ADR-0111` — nur der akzeptierte-Negativ-Punkt
  „Schema-Version-Verweis")
- [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md)
  (Proposed, `Supersedes ADR-0111` — nur die Klassen-Aufzählung in Teilfrage 5)
- dieses Dokument; **keine** ADR zum dritten Punkt (Begründung dort)

Den ADR-Index (`docs/plan/adr/README.md`) setzt der Aufrufer; die Pläne, die
Spec und der Code bleiben in diesem Zug unberührt.

---

## Verdikt

**F-8 — Lücke im Wortlaut, keine Lücke in der Sache; per Folge-ADR benannt
(Verdikt 2/3).** Der Run liest die Version vor dem Snapshot, sie wechselt allein
im Erfassungspfad; an den realen Adaptern gemessen zeigen die Fälle „nach der
Aktivierung" und „nach der letzten Relation-Nachricht, ohne WAL-Change" eine
Referenz, die hinter dem Bild liegt bzw. auf eine Version ohne Spaltenform zeigt.
Das ist harmlos, weil kein Consumer-Weg die Spaltenform je Version liefert — aber
es steht nicht dort, wo ein Leser es sucht. Entscheidung: die Referenz ist ein
**Marker** (Version zum Run-Start), keine Katalog-Prüfung, kein Abbruch, keine
neue Version. Träger: `ADR-0116`.

**F-9 — Übergabepunkt ohne Entscheidung; per Folge-ADR entschieden
(Verdikt 2/3).** Ein Run, dessen Regelmenge auf die Snapshot-Spalten nicht
anwendbar ist, endet `failed` mit Klasse **`schema`** — dieselbe Ursache, dieselbe
Klasse wie im Erfassungspfad —, run-lokal, vor der ersten Kopie. Abhilfe: Regel
ändern, neuer Antrag. Träger: `ADR-0117`; der Start-Trigger des Backfill-Pfads der
Transformationen ist mit diesem Dokument **plus dem Status `Accepted` von
`ADR-0117`** erfüllt (siehe §Auftraggeber-Frage).

**F-1 — Träger fehlte, keine Aussage war falsch; keine ADR, Nachzug in die Spec
(Verdikt 3/3: die Lockerung ist legitim, aber undokumentiert und wird
nachgezogen).** Keine der drei `Accepted`-ADRs sagt etwas Falsches über die
Rechte auf `cdc.administration_request` (§Befund 3): sie treffen dazu **keine**
Aussage. Der Grant `SELECT`, `UPDATE` für `cdc_admin` ist im Code gelandet und
richtig; es fehlt der Ort, an dem das Pflichtenheft die Rechte der Antrags-Queue
führt.

---

## Befundlage

### 1 — Schema-Version der Backfill-Changes (F-8)

Gemessen am 2026-09-24: PostgreSQL **17.11** und **18.6** (Digests wie in
`ADR-0116` §Gemessen), Wegwerf-Container, Rollout des Arbeitsbaums am Stand
`77fc6e43`, die realen Adapter unter dem Use Case des Runs. Zwei Tabellen mit je
zwei Zeilen, nach der Aktivierung `ALTER TABLE … ADD COLUMN`, keine WAL-Change,
dann Run:

| Fall | Referenz der Backfill-Changes | Bild |
|---|---|---|
| Aktivierung, dann `ADD COLUMN` | Version 1, **keine** Zeile in `cdc.table_schema` | trägt die neue Spalte |
| Spaltenform `[id, a]` nachgetragen, dann `ADD COLUMN` | Version 1, Spaltenform `[id, a]` | trägt die neue Spalte `b` |

Beide Läufe endeten `completed` mit 2 Zeilen, identisch auf 17.11 und 18.6;
`CurrentVersion` blieb Version 1. **Reproduziert** (Frage b). Die Optionen und ihre
Nachteile stehen in `ADR-0116`; der Ausschlag: der Snapshot trägt keine Typ-OIDs
(`TableSnapshot.Columns()` liefert Namen), eine vom Run angelegte Version bliebe
der `TableBinding` des `Assembler` unbekannt (Rückwärtsverlauf der Versionen, aus
dem Code hergeleitet), und ein Abbruch träfe den Regelfall „aktivieren, dann
sofort Backfill".

**Der WAL-Anteil der Probe ist nachgestellt** (`RegisterVersion` wie in
`Assembler.observeRelation`), der Replication-Strom lief nicht: die Abhilfe
„eine Zeile ändern, dann Antrag" (`ADR-0116` Festlegung 3) ist **erwartet**, nicht
gemessen, und keine Zusage.

### 2 — Fehlerklasse `schema` im Run (F-9)

Am Code belegt (Symbolnamen, keine Läufe): `model.NewErrorClass` akzeptiert alle
sieben Klassen, `BackfillRun.Fail` schreibt `<Klasse>: <Ursache>` in
`error_message`, die Tabellenspalte trägt keine Prüfbedingung auf die Klasse —
`schema` im Run braucht **keinen** Modell-, Schema- oder Port-Diff. Die
Nichtanwendbarkeit ist nach `ADR-0112` Teilfrage 2 eine Eigenschaft von
Regelmenge × Spaltenmenge (`map_value` lässt jeden anderen Wert unverändert), also
einmal je Run gegen `Columns()` prüfbar, vor der Schreibtransaktion. Ein
Run-Abbruch verliert nichts (eine Transaktion, Rollback, keine Bestätigung); der
Halt des Erfassungspfads bei Nichtanwendbarkeit gründet auf Persist-before-ACK
und gilt für den Run nicht.

### 3 — Rollenschnitt der Antrags-Queue (F-1)

**Messung** (2026-09-24, PostgreSQL 18.6 Alpine, Digest von `PG_TEST_IMAGE`,
Wegwerf-Container, Rollout des Arbeitsbaums am Stand `77fc6e43`;
`has_table_privilege` je Rolle und Recht, gedruckt):

| Rolle | `administration_request` | `backfill_run` |
|---|---|---|
| `cdc_admin` | `SELECT` ja · `UPDATE` ja · `INSERT` nein · `DELETE` nein | `SELECT` ja · `INSERT` ja · `UPDATE` nein · `DELETE` nein |
| `cdc_capture` | keines | `SELECT` ja · `UPDATE` ja · `INSERT` nein · `DELETE` nein |
| `cdc_reader` | keines | keines (liest über die View) |

Der Stand deckt sich mit `tools/schema/nacharbeit-roles.sql` und mit der Grants-Tabelle
von `SPEC-029` für `backfill_run`. Lesen der drei Accepted-Stellen:

- **`ADR-0047`** — die Aufrufer-Tabelle („Verdrahtung") nennt die
  Administrations-Verarbeitung der Antrags-Queue nicht; sie ist mit
  `ADR-0050` später hinzugekommen. Die ADR sagt dazu nichts aus.
- **`ADR-0050`** — nennt die Rolle `cdc_admin` für die Goroutine, keine Grants.
- **`ADR-0113`** — die Grants-Tabelle steht ausdrücklich „auf `cdc.backfill_run`";
  der Satz „`cdc_admin` trägt kein `UPDATE`" gilt in diesem Zusammenhang. Der Verweis
  „Zuordnung `ADR-0047`, Tabelle ‚Verdrahtung'" für die Goroutine ist als Beleg für die
  **Rolle** ungenau (der Beleg ist `ADR-0050`), aber keine Rechte-Aussage; das ist eine
  Zitat-Ungenauigkeit ohne Folge und wird nicht korrigiert.

`SPEC-019` (Antrags-Datensatz) trägt keine Rollenrechte; das Pflichtenheft führt
die Grants **je Tabelle** in der Spec-Stelle der Tabelle (Vorbild `SPEC-029`), es hat
keine eigene Rollen-Spezifikation (Suchlauf `grep -n 'cdc_capture\|cdc_admin\|cdc_reader' spec/*.md`: die Treffer, die
eine Rolle meinen, stehen nur in `SPEC-029`; die übrigen sind der Kanalname
`cdc_administration` und die Kennzahl `cdc_capture_lag`).

**Entscheidung: kein ADR.** Eine ADR hielte eine Entscheidung fest, die eine
`Accepted`-Aussage ändert oder ergänzt; hier ändert nichts eine Aussage. Der
Grant folgt aus `ADR-0050` (die Goroutine läuft als `cdc_admin` und liest/vermerkt
Anträge) und aus `ADR-0113` (die Annahme vermerkt den Antrag in ihrer
Transaktion); er ist die notwendige Bedingung ihrer Umsetzung, keine neue
Wahl. Eine ADR würde ein Dokument mehr erzeugen, ohne dass eine Zusage anders
wird.

**Akzeptiertes Negativ (kurz begründet).** `UPDATE` steht auf Tabellenebene: `cdc_admin`
könnte auch `request_kind`/`column_name` einer `applied`-Zeile umschreiben, die den
Ausschlussstand trägt (`ADR-0065`). Ein Spalten-Grant `UPDATE (status, error_message)`
wäre feiner; er kostet eine Änderung der Grant-Zeile, der Rollen-Tests und der
Rollout-Datei für einen Angreifer, der bereits `cdc.include_column` aufrufen kann —
dieselbe Begründung wie `ADR-0113` Festlegung 1 „Akzeptiertes Negativ". Wird ein
Spalten-Grant gewollt: eigene Entscheidung.

**Akzeptiertes Negativ (Ursache der Latenz).** Der Compose-Aufbau fährt alle drei DSNs
als Superuser; ein fehlender Grant fällt dort nicht auf. Der Beleg liegt bei den
Login-Tests im Store-Tier (`roles_test.go`, `backfillroles_test.go`,
`administration_roles_internal_test.go` — Login `IN ROLE cdc_admin`/`cdc_capture`),
nicht bei `make test-integration`. Das genügt: eine Compose-Umgebung mit
Rollen-Logins wäre ein Neuaufbau der Integrationsumgebung für eine Klasse, die der
Store-Tier bereits belegt.

---

## Folge-Arbeit je Slice

Jede Pflicht hat einen Träger; die Pläne ändert dieses Dokument nicht — der Planner
zieht sie nach (`AGENTS.md` §3.13).

| Slice (Kennung) | Nachzug | Herkunft |
|---|---|---|
| `slice-backfill-sql-administration` | **Spec-Zug:** `SPEC-019` bekommt einen Absatz „Grants" im Muster von `SPEC-029` (Text unten); „Berührte Spec-Stellen" ändert sich von „gelesen, nicht geändert" auf „`SPEC-019` geändert"; DoD-Punkt, Suchlauf (§3.13) über `SPEC-019` und Handbuch §2/§4 (Handbuch trägt den Rechteschnitt seit 1.47 bereits, der Zug liest ihn, überschreibt ihn nicht) | Punkt 3, F-1 |
| `slice-backfill-sql-administration` | **Spec-Satz zur Version:** in `LH-FA-CAP-009.a`, Absatz „Markierung", ein Satz (Wortlaut in `ADR-0116` Folgepflicht 3); ein Satz im Handbuch-Abschnitt zum Backfill (Folgepflicht 2) — der Slice schreibt diesen Abschnitt ohnehin neu | Punkt 1, `ADR-0116` |
| `slice-backfill-e2e` | Start-Trigger „Architect-Verdikt zur Schema-Version": mit diesem Dokument erfüllt (der Übergangs-Commit `next` → `in-progress` nennt es) **und** `ADR-0116` `Accepted`; der Beleg der Bild-Form prüft an der Version nur die Zeile aus `cdc.schema_version` zum Antrag (`ADR-0116` Festlegung 4, Fitness Function), nicht die Übereinstimmung mit den Bild-Spalten; der Grenze-Kommentar an `currentVersion` (Folgepflicht 1, ein Kommentar-Diff in `service.go`, kein Verhaltens-Diff) gehört in diesen Slice, weil kein früherer Slice `service.go` mehr anfasst | Punkt 1 |
| `slice-transformationen-backfill-pfad` | Start-Trigger „Architect-Kurzverdikt zur Nichtanwendbarkeit im Run": mit diesem Dokument erfüllt **und** `ADR-0117` `Accepted` (eine `Accepted`-ADR unberührt lassen ist Bedingung der Rückführung `in-progress` → `open`, die der Plan nennt — sie entfällt, weil die Folge-ADR bereits vorliegt); DoD-Punkt „Fail-closed und Nichtanwendbarkeit" nennt Klasse `schema`, Prüfung einmal je Run vor `Begin` (`ADR-0117` Festlegung 2), Regelstand-Wechsel `configuration` (Festlegung 5); **Spec-Zug:** Zeile `schema` von `SPEC-008` und ein Satz in `LH-FA-CAP-009.a` (`ADR-0117` Folgepflicht 1) — der Plan sagt „Doku-Update: entfällt" und „Spec gelesen", beides ändert sich; der Kommentar „Die Klasse `schema` vergibt der Run nicht" an `classifyError` entfällt mit der Umsetzung | Punkt 2 |
| `slice-transformationen-e2e-abhilfe` | die Abhilfe des Runs (Regel ändern, neuer Antrag; `ADR-0117` Festlegung 4) als zweiter Rundlauf neben dem bestehenden Abhilfe-Kriterium des Erfassungspfads — Empfehlung an den Planner, kein Zwang | Punkt 2 |

**Wortlaut-Vorschlag `SPEC-019`, Absatz „Grants"** (Rollen nach der Zuordnung der
DSN-Verdrahtung; der Slice übernimmt oder schärft ihn):

| Rolle | Recht auf `cdc.administration_request` | Träger |
|---|---|---|
| `cdc_admin` | `SELECT`, `UPDATE` | die Administrations-Verarbeitung: offene Anträge lesen, den Ausgang (`applied`/`failed`) vermerken, die dauerhaften Stände aus den `applied`-Zeilen ableiten; die Annahme eines Backfills vermerkt in derselben Transaktion |
| `cdc_capture` | **keines** | — |
| `cdc_reader` | **keines** | — |

Niemand trägt `INSERT` oder `DELETE`: Anträge legen ausschließlich die
SQL-Funktionen an (`SECURITY DEFINER`, unter den Rechten ihres Eigentümers).

---

## Auftraggeber-Frage (einzige)

**Annahme von `ADR-0116` und `ADR-0117`** (beide `Proposed`). Ohne `Accepted` bleibt
der jeweilige Start-Trigger offen: eine `Proposed`-ADR ist keine Entscheidung, auf
die ein Slice bauen darf. Alles Übrige dieses Verdikts (Punkt 3) braucht keine
Annahme; es ist ein Nachzug in die Spec.
