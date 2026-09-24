# ADR-0117: Backfill — Fehlerklasse `schema` für eine im Run nicht anwendbare Transformationsregel, Klassenmenge des Runs (Supersedes ADR-0111, Teilfrage 5 teilweise)

**Status:** Accepted — Supersedes [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
in **Teilfrage 5, teilweise** (genau ein Satzteil: die Aufzählung der fünf
Fehlerklassen im Absatz „Fehlerklassen und Heartbeat", siehe §Entscheidung);
alles Übrige von `ADR-0111` bleibt in Kraft. [`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
wird **nicht** ersetzt, sondern in seiner Folgepflicht 7 für den Run ausgefüllt.

**Datum:** 2026-09-24

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch den Befund F-9 des
Review-Reports zum Run-Use-Case (`docs/reviews/`) und durch den Start-Trigger des
Backfill-Pfads der Transformationen, der ein Kurzverdikt zur Nichtanwendbarkeit im
Run verlangt)

**Bezug:** [`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Transformation vor
Persistierung; Negative: keine still unveränderte Auslieferung),
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill des Bestands),
[`LH-FA-SCH-004`](../../../spec/lastenheft.md) (inkompatible Änderungen,
sichtbarer Fehler), [`LH-QA-REL-001`](../../../spec/lastenheft.md) (kein
Datenverlust, Persist-before-ACK), [`LH-QA-SEC-004`](../../../spec/lastenheft.md)
(nur zur Abgrenzung), [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
(teilweise superseded — Haupt-Bezug),
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Regeln, Nichtanwendbarkeit im Erfassungspfad, Folgepflicht 7),
[ADR-0023](0023-fehlerklassifikation.md) (sieben Fehlerklassen),
[ADR-0012](0012-at-least-once.md) (Persist-before-ACK)

**Schärft:** [`SPEC-008`](../../../spec/pflichtenheft.md) (Zeile `schema`:
Bedingung um die Nichtanwendbarkeit einer Regel ergänzt, im Erfassungspfad **und**
im Run), [`LH-FA-CAP-009.a`](../../../spec/pflichtenheft.md) (Absätze
„Fail-closed vor dem Commit" und „Sichtbarkeit und Fehler des Runs"),
[`SPEC-029`](../../../spec/pflichtenheft.md) (`error_message` trägt die Klasse aus
`SPEC-008` — unverändert wahr)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) und
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md) sind
`Accepted` und unberührbar (`AGENTS.md` §3.5). Die eine nennt für Run-Fehler fünf
Klassen und sagt „run-lokal", die andere beschreibt die Nichtanwendbarkeit einer
Regel nur für den Erfassungspfad. Der Run des Bestands trägt dieselbe
Regelauswertung (`ADR-0112` Folgepflicht 7); welche Klasse er bei
Nichtanwendbarkeit trägt, steht in keiner der beiden. Die Ergänzung ändert den
Referenten und ist keine Zitat-Korrektur.

### Befunde am Bestand (Code-Anker als Symbolnamen, Stand `77fc6e43`)

- **Die fünf Klassen von `ADR-0111`.** Teilfrage 5, Absatz „Fehlerklassen und
  Heartbeat": „`permission`, `configuration`, `storage`, `transient`,
  `replication`" — beschreibend für die Ursachen, die der Run bis dahin kannte.
- **Die Abbildung im Use Case.** `classifyError`
  (`internal/application/usecase/backfill/service.go`) bildet die fünf
  Snapshot-Sentinels, die Storage-Sentinels, `ErrTableNotActivated`,
  `ErrExclusionStateChanged` und `ErrSchemaVersionUnknown` (auf `configuration`)
  ab; ein nicht erkannter Fehler endet als `internal`; der Doc-Kommentar hält fest,
  dass der Run `schema` nicht vergibt.
- **Der Domänentyp trägt alle sieben.** `model.NewErrorClass`
  (`internal/domain/model/errorstate.go`) akzeptiert `schema`;
  `BackfillRun.Fail` prüft die Klasse dort und schreibt sie als Präfix in
  `error_message` („`<Klasse>: <Ursache>`", `internal/domain/model/backfillrun.go`);
  `cdc.backfill_run.error_message` ist freier Text ohne Prüfbedingung auf die
  Klasse (`tools/schema/schema.yaml`, `SPEC-029`).
- **Die Nichtanwendbarkeit im Erfassungspfad** (`ADR-0112` Teilfrage 4): eine
  Regel ist nicht anwendbar, wenn ihre `column` in der Relation der Change fehlt
  oder ein Zielname mit einer Spalte der Relation kollidiert; der Erfassungspfad
  endet sichtbar mit der Klasse `schema` — Grund: keine Transaktion wird
  persistiert oder bestätigt, sonst ginge eine Change unverändert oder gar nicht
  hinaus (`ADR-0012`).
- **Die Nichtanwendbarkeit ist eine Eigenschaft von Regelmenge und Spaltenmenge,
  nicht einer Zeile.** `ADR-0112` Teilfrage 2: `map_value` lässt „jeden anderen
  Wert" unverändert, ein Wert, der im Image fehlt, bleibt abwesend; keine der beiden
  Regeln scheitert an einem Zeilenwert. Sie ist deshalb einmal je Run gegen die
  Spaltenliste des Snapshots (`TableSnapshot.Columns()`) prüfbar, bevor eine Zeile
  gelesen wird.
- **Ein Run verliert beim Abbruch nichts.** Er schreibt in **eine** Transaktion,
  die bei jedem Fehler zurückgerollt wird (`ADR-0111` Teilfrage 4), hält keine
  Position im Sinne des Erfassungspfads und bestätigt nichts.

### Konstraints

- `SPEC-008`/`ADR-0023`: die Menge der Klassen ist geschlossen (sieben); keine
  achte Klasse.
- `ADR-0111` Teilfrage 5: Run-Fehler sind **run-lokal** — weder Heartbeat-Fehler-
  zustand noch Halt des Erfassungspfads. Dieser Satz bleibt.
- `LH-FA-CFG-007` Negative: eine nicht anwendbare Regel darf nicht in stiller
  Rohform ausliefern; `LH-QA-SEC-004`: Ausschluss vor Regel (`ADR-0112`
  Teilfrage 5) — unberührt.

## Entscheidung

Wir wählen **Klasse `schema` für einen Run, dessen Regelmenge auf die Spalten des
Snapshots nicht anwendbar ist — run-lokal, vor der ersten Kopie, mit der Abhilfe
„Regel ändern, neuer Antrag"**. Fünf Festlegungen.

**Was diese ADR von `ADR-0111` ersetzt** — genau einen Satzteil: in Teilfrage 5,
Absatz „Fehlerklassen und Heartbeat", die Klammer „(`permission`, `configuration`,
`storage`, `transient`, `replication`)" (gilt jetzt: Festlegung 1). Der Satz
„Run-Fehler … sind **run-lokal**: sie setzen weder den Heartbeat-Fehlerzustand
noch stoppen sie den Capture-Pfad" bleibt wörtlich in Kraft und gilt für jede
Klasse.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: die Nichtanwendbarkeit bleibt unklassifiziert und endet als `internal` (Rückfall von `classifyError`) | keine Änderung | `internal` heißt nach `SPEC-008` „unerwarteter interner Fehler; Restart-Strategie": falsch für einen erwartbaren, vom Betreiber behebbaren Zustand; `diagnose` und `error_message` sagten nicht, dass eine Regel die Ursache ist |
| B — Klasse `configuration` | eine Regel ist Konfiguration; die Klasse ist im Run bereits für Zustandsabweichungen vergeben | dieselbe Ursache hieße je Pfad anders (`schema` im Erfassungspfad, `configuration` im Run): ein Betreiber, der nach der Klasse sucht, sieht zwei Bezeichnungen für dasselbe Ereignis; die Regel war beim Antrag gültig (K4 prüft die Spalte), erst der Stand der Relation macht sie unanwendbar — `configuration` bedeutet „ungültig gesetzt" (`SPEC-008`), hier ist nichts ungültig gesetzt |
| **C — `schema` (gewählt)** | dieselbe Ursache, dieselbe Klasse in beiden Pfaden (`ADR-0112` Teilfrage 4); keine neue Klasse (`ADR-0023`); der Domänentyp und die Spalte `error_message` tragen sie bereits — kein Schema-, Port- oder Modell-Diff | der Run trägt jetzt sechs Ursachen-Klassen statt fünf; das Wort „`schema`" meint im Run „Regelmenge passt nicht zur Form der Tabelle", nicht nur Relation-Metadaten (durch die Ergänzung der `SPEC-008`-Zeile aufgelöst, Folgepflicht 1) |
| D — die Regel überspringen und den Bestand roh kopieren | Run läuft durch | die still unveränderte Auslieferung, die `LH-FA-CFG-007` Negative und `ADR-0112` Teilfrage 4 Option A ausschließen; Rohdatenaustritt |
| E — Zeilen mit nicht anwendbarer Regel auslassen, Run `completed` | Run läuft durch | Bestand unvollständig bei Status `completed`: ein stiller Verlust, den `ADR-0112` Teilfrage 4 Option B ausschließt |
| F — den Erfassungspfad stoppen, wie bei der Nichtanwendbarkeit einer WAL-Change | gleiche Wirkung in beiden Pfaden | der Run hat keine Position und keine Bestätigung, ein Abbruch verliert nichts; ein `cdc.backfill_table`-Antrag würde die Erfassung der ganzen Quelle stoppen, obwohl der Erfassungspfad sein eigenes Kriterium hat und dieselbe Ursache bei der nächsten Change der Tabelle selbst meldet; widerspricht dem Satz „run-lokal" von `ADR-0111` |

### Festlegung 1 — Klassenmenge des Runs

Ein Run endet `failed` mit genau einer Klasse aus der Menge von `SPEC-008`. Die
Ursachen, die er vergibt: `permission`, `configuration`, `storage`, `transient`,
`replication` (`ADR-0111`) und **`schema`** — eine Regel des Regelstands der
Tabelle ist auf die Spalten des Snapshots nicht anwendbar (Definition:
`ADR-0112` Teilfrage 4, gegen `TableSnapshot.Columns()` statt gegen die Relation
einer WAL-Change). `internal` bleibt der Rückfall für einen nicht erkannten
Fehler. Die Abbildung von `ErrSchemaVersionUnknown` auf `configuration` (eine
aktivierte Tabelle ohne Versions-Zeile; `TableActivationAdapter.Register` schreibt
beide in einer Transaktion) bleibt bestätigt.

### Festlegung 2 — Wirkort und Zeitpunkt der Prüfung

Die Prüfung liegt im Use Case des Runs, nutzt **dieselbe** Prüffunktion der
Domäne wie der Erfassungspfad (`ADR-0112` Folgepflicht 7, `ADR-0111` Teilfrage 2:
eine Funktion für beide Pfade) und läuft **einmal je Run**, nachdem der Snapshot
seine Spalten liefert und **bevor** die Schreibtransaktion geöffnet und die erste
Zeile gelesen wird. Ergebnis: der Run endet `failed`, `error_message` beginnt mit
`schema: ` und nennt Regelname und Spalte; der Snapshot ist geschlossen, es
besteht keine Schreibtransaktion, es entsteht keine Change.

### Festlegung 3 — Run-lokal, Sichtbarkeit

Der Fehler ist run-lokal (`ADR-0111` Teilfrage 5, unverändert): er setzt weder den
Heartbeat-Fehlerzustand noch stoppt er den Erfassungspfad. Sichtbar ist er in
`cdc.backfill_status` und in `diagnose` (Run `failed`, Klasse im Fehlertext). Der
Erfassungspfad endet bei derselben Ursache erst nach seinem eigenen Kriterium
(`ADR-0112` Teilfrage 4), wenn eine Change der Tabelle bei ihm ankommt — nicht
durch den Run.

### Festlegung 4 — Abhilfe

Die Abhilfe für den Run ist die Änderung des Regelstands — eine Regel entfernen
(`cdc.remove_transformation`) oder durch eine passende ersetzen
(`cdc.set_transformation` nach dem Entfernen, `ADR-0112` K1) — und danach ein
**neuer Antrag** `cdc.backfill_table`. Ein `failed`-Run wird nicht fortgesetzt
(`ADR-0111` Teilfrage 4), hinterlässt keine Change (Festlegung 2) und ist nicht
„aktiv" (`queued`/`running`): der neue Antrag ist zulässig. Ein Prozessneustart
ist für den neuen Run nicht nötig, sofern der Run seinen Regelstand über den
Port neu liest — **erwartet**, zu belegen im E2E-Beleg des Backfill-Pfads der
Transformationen.

### Festlegung 5 — Abgrenzung: Regelstand-Wechsel im Lauf

Wechselt der Regelstand zwischen zwei Lesungen des Runs (Regel gesetzt, dann
entfernt, oder umgekehrt), endet der Run wie bei einem Wechsel des Ausschlussstands
mit `configuration` (`ErrExclusionStateChanged`-Abbildung, `ADR-0111` Teilfrage 4
Fail-closed): der **Zustand** wechselt, es ist nicht eine Regel auf eine Form
nicht anwendbar. Die Klasse `schema` bleibt der Nichtanwendbarkeit vorbehalten.

## Konsequenzen

- Positiv: dieselbe Ursache trägt in beiden Erzeugungspfaden dieselbe Klasse; der
  Betreiber sucht nach einem Wort.
- Positiv: kein Modell-, Schema- oder Port-Diff — der Domänentyp und die
  Fehlertext-Spalte tragen `schema` bereits; die Änderung ist eine Zeile in
  `classifyError`, ein Sentinel-Fehler der Regelprüfung und die Prüfung selbst.
- Positiv: der Erfassungspfad bleibt von einem Antrag unberührt; ein Antrag kann
  die Erfassung der Quelle nicht stoppen.
- Negativ: der Run nennt `schema` auch für einen Fall, den der Erfassungspfad
  (noch) nicht erreicht hat; der Betreiber sieht einen Run-Fehler früher als den
  Halt des Erfassungspfads — das ist die Absicht (der Run entdeckt, was die
  nächste Change ohnehin meldete).
- **Akzeptiertes Negativ** (kurz begründet, keine Folgepflicht): *ein Wechsel des
  Regelstands zwischen der Prüfung (Festlegung 2) und dem Bau eines späteren
  Blocks* — wird von der Fail-closed-Prüfung je Block und vor dem Commit
  gefangen (Festlegung 5), nicht von der einmaligen Prüfung.

### Folgepflichten

Jede Pflicht hat einen Träger im Architect-Verdikt zu dieser ADR
(Verzeichnis `docs/reviews/`, Verdikt „backfill-schema-klasse-rollen“); **die
genannten Träger ändert diese ADR nicht** — der Planner zieht sie nach
(`AGENTS.md` §3.13).

1. **Spec:** die Zeile `schema` von `SPEC-008` nennt die Nichtanwendbarkeit einer
   Regel („im Erfassungspfad und im Run"); `LH-FA-CAP-009.a` sagt in einem Satz,
   dass eine im Run nicht anwendbare Regel den Run `failed` mit Klasse `schema`
   beendet, ohne Change und ohne den Erfassungspfad zu berühren.
2. **Backfill-Pfad der Transformationen:** die Prüfung (Festlegung 2) mit
   Sentinel und Abbildung in `classifyError` (Kommentar „vergibt `schema` nicht"
   entfällt), die Regelstand-Abweichung (Festlegung 5) auf `configuration`,
   Negativtests je Eingabe gebunden.
3. **E2E-Beleg:** nicht anwendbare Regel → Run `failed`/`schema`, keine Change,
   Erfassungspfad läuft weiter; nach Abhilfe (Festlegung 4) neuer Run `completed`.
4. **Handbuch:** Betriebshinweis zur Abhilfe.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Fakes) | eine Regel, deren Spalte nicht in `Columns()` steht, und eine Regel, deren Zielname mit einer Spalte kollidiert, enden den Run `failed` mit Klasse `schema` ohne `Begin` der Schreibtransaktion; je Negativfall eine Mutation der Prüfung färbt den Test rot; ein Regelstand-Wechsel endet `configuration` | `make test` |
| Go-Test (Fakes) | der Erfassungspfad-Halt und der Heartbeat-Fehlerzustand bleiben aus (kein Aufruf an `HeartbeatPort`) | `make test` |
| Review-Prüfpflicht | die Prüffunktion der Regelanwendbarkeit steht an einer Stelle der Domäne, beide Pfade rufen sie | — (kein Gate) |

## Re-Evaluierungs-Trigger

- **Der Erfassungspfad ändert die Behandlung der Nichtanwendbarkeit** (z. B.
  Überspringen statt Halt): Klasse und Abhilfe des Runs neu paaren.
- **Ein dritter Regeltyp** (`ADR-0112`) hat eine Zeilen-Nichtanwendbarkeit
  (Wert-abhängig): die Prüfung je Run reicht dann nicht — Folge-ADR.
- **Parallele Runs oder Fortsetzung mit Checkpoint** (`ADR-0111`-Ausbaustufe):
  Wirkort und Zeitpunkt der Prüfung neu bewerten.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-24 | Proposed — Architect-Vorschlag zum Review-Befund F-9 und zum Start-Trigger des Backfill-Pfads der Transformationen | Review-Report zum Run-Use-Case, Befund F-9 (`docs/reviews/`) |
| 2026-09-24 | Accepted — Annahme durch den Auftraggeber samt Klasse schema für eine im Run nicht anwendbare Regel, run-lokal | [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0117` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
