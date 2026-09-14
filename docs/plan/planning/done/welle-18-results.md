# Welle 18 — Spaltenauswahl — Antrags-Queue-Erweiterung, Assembler-Filterung, E2E-Beleg — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-NNN` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-18
**Abschluss:** 2026-09-14
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Spaltenauswahl) trägt
  den vollständigen, in [`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md)
  entschiedenen Umsetzungspfad. Die Messmethode von
  `LH-QA-SEC-004` (Ausschluss sensibler Spalten) läuft ausschließlich über
  diese Anforderung und ist damit ebenfalls erstmals real belegt; die
  `ADR-0059` Folgepflichten — `ARC-005`-Sequenzdiagramm in
  [`spec/architecture.md`](../../../../spec/architecture.md), neuer
  [`SPEC-019`](../../../../spec/pflichtenheft.md) für die Feldform des
  erweiterten Antrags-Datensatzes — liegen vor.
- `slice-066` — Antrags-Seite: zwei neue Antragsarten
  `cdc.exclude_column`/`cdc.include_column` auf der bestehenden
  Antrags-Queue ([`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md));
  `cdc.administration_request` um `column_name` und die vierwertige
  `request_kind`-Menge erweitert; neuer Outbound Port mit `ColumnExists`,
  `ExcludeColumnUseCase`/`IncludeColumnUseCase`, Sentinel
  `ErrSourceColumnMissing`, `applyAdministrationRequest` um zwei
  `case`-Zweige. Real gegen PostgreSQL belegt: `applied` im Happy Path,
  `failed` samt Fehlertext gegen eine nicht existierende Spalte.
- `slice-067` — Wirkort: `TableBinding.ExcludedColumns` und die Filterung
  in `rowImage`/`change` **vor jeder Serialisierung** (`ADR-0059`
  Teilfrage 3 Option D); eine neue, unter `tablesMu` synchronisierte
  Live-Reload-Methode; der Schema-Bump-Pfad erhält den Ausschlussstand an
  **beiden** Voll-Schreibstellen (`observeRelation`, `AddBinding`), statt
  ihn zurückzusetzen. Konvergenz mit
  [`LH-FA-SCH-003`](../../../../spec/lastenheft.md) bei realer
  Spaltenlöschung ohne Sonderfall (reine Schichtung), Filterung über die
  bestehende `nil`-Abwesenheits-Kodierung
  ([`LH-FA-DAT-005`](../../../../spec/lastenheft.md)).
- `slice-068` — E2E-Beleg am **laufenden** Feed-Container (`make
  test-integration`, kein Gate): `SELECT cdc.exclude_column(...)` gegen
  eine über `cdc.enable_table` aktivierte, nicht in `CDC_TABLES` gelistete
  Tabelle, Poll auf `status = 'applied'`, danach eine erfasste Change ohne
  den ausgeschlossenen Spaltenschlüssel und ohne den Wert — ohne Neustart.
  Eine **Baseline**-Change vor dem Ausschluss trägt den Wert real (schließt
  „die Spalte war nie im Row Image" als Alternativerklärung aus), die
  nicht ausgeschlossene Spalte bleibt (trennt gezielten von totalem
  Filter), die vor dem Ausschluss erfasste Change bleibt unverändert
  lesbar. Negative-Beleg: `failed` samt Fehlertext.
- [`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  (Accepted) — `Supersedes ADR-0059` nur in der Dauerhaftigkeits-Aussage —
  und der wellenlose Folge-Slice `slice-075`.
- `harness/README.md` §Sensors: Zeile `make test-integration` um den neuen
  Rundlauf-Abschnitt ergänzt (inklusive der tragenden `LH-QA-SEC-004`-
  Kennung), Herkunfts-Anker `· seit slice-068`.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die Wiederverwendung der bestehenden Antrags-Queue trug die Erweiterung
  ohne neuen Mechanismus: zwei neue Antragsarten, ein neuer Outbound Port
  am bereits vorhandenen `cdc_admin`-Pool, dieselbe Live-Reload-Verdrahtung
  und derselbe `LISTEN`/`NOTIFY`-Pfad aus `ADR-0050`.
- Die Filterung nutzt die im Row Image bereits vorhandene
  `nil`-Abwesenheits-Kodierung — kein neuer Platzhalter, `LH-FA-DAT-005`s
  Negative-Kriterium damit ohne neue Kodierung erfüllt.
- Der Implementer fand selbst, dass **zwei** (nicht eine) Stellen eine
  vollständige `TableBinding` schreiben, und belegte mit Mutationen, dass
  beide einzeln tragend sind.
- Der E2E-Rundlauf brauchte weder Neustart noch neue Infrastruktur; die
  beiden Konstruktionen Baseline und Kontrollspalte schlossen beide
  Alternativerklärungen real aus. Implementer, Reviewer und Verifier
  setzten je eigene Mutationen und sahen sie rot — über vier unabhängige
  vollständige `make test-integration`-Läufe kein Flake.
- Der Reviewer fing den HIGH-Fund (Slice-Chronik in Produktionscode) vor
  dem Merge: die bereits verkörperte Verteidigungslinie
  (`BEO-PGC/slice-chronik-in-code-kommentar`, Reviewer-HIGH-Punkt) trug,
  und der datei-skopierte Enumerationslauf aus `.claude/commands/`
  förderte vier weitere Produktionscode-Stellen derselben Klasse zutage,
  die über den ursprünglichen Review-Befund hinausgingen.
- Zwei unabhängige Belege statt einer Behauptung beim Closure-Trigger: der
  Verifier hat den realen E2E-Rundlauf in eigener Sitzung reproduziert
  (`docs/reviews/verify-slice-068.md` §2/§3) und den Slice als
  DoD-konform und `welle-18`-closure-reif gemeldet.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Die `d-migrate`-Grenze traf auch eine Tabellenänderung.** Der Plan
  führte die `request_kind`-CHECK-Klausel deklarativ in
  `tools/schema/schema.yaml`; real trägt d-migrate 1.3.1 eine CHECK-Klausel
  an einer **bestehenden** Tabelle nicht (Exit 5, `POST_EXECUTE_DRIFT`, die
  bestehende Klausel entfällt beim Änderungsversuch) — die Spalte
  `column_name` konvergiert, die Klausel nicht. Konsequenz: die Klausel
  liegt in der etablierten Ausweichform
  `tools/schema/nacharbeit-administration.sql` (idempotent, `DROP
  CONSTRAINT IF EXISTS` vor `ADD CONSTRAINT`), die Spalte bleibt
  deklarativ. Der Befund ist eine **vierte Objektklasse** derselben
  Werkzeuggrenze und als Beleg in `BEO-PGC/d-migrate-nacharbeit`
  nachgetragen (jetzt 6×, Ausgang unverändert `verkörpert`).
  Zusatzbefund aus demselben Lauf: der Rollout des aktuellen Baums gegen
  eine Instanz auf dem Stand *davor* endet Exit 8
  (`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`); der Diff erzeugt eine
  zusätzliche destruktive Plan-Operation auf
  `administration_request.request_kind`, weil ein im Katalog liegendes
  Objekt aus dem deklarativen Modell genommen wurde.
- **Der Review deckte eine Lücke über den Plan hinaus auf — mit
  Entscheidung, nicht nur Benennung.** `slice-067` führte den
  Ausschlussstand als reine Laufzeit-Eigenschaft der `Assembler`-Bindung.
  Der Reviewer fand, dass er damit **keinen dauerhaften Träger** hat, und
  zwar enger als im §6-Risiko benannt: verloren geht er nicht nur beim
  Prozess-Neustart, sondern auch bei einem
  `cdc.disable_table`/`enable_table`-Zyklus ohne Neustart. Der
  Architect-Zug entschied das ausdrücklich als **Lücke**, nicht als
  zulässige Grenze. Konsequenz: [`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  (`Supersedes ADR-0059`, nur die Dauerhaftigkeits-Aussage), der
  wellenlose Folge-Slice `slice-075` und der Registereintrag
  `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (Zustand `geplant`,
  Träger `slice-075`, 1×). `welle-18` §6 trägt die Grenz-Benennung; der
  Closure-Trigger der Welle war davon ausdrücklich unberührt (er verlangt
  den Beleg des *laufenden* Pfads, nicht Dauerhaftigkeit).
- **`slice-068` §1 sagte mehr zu als der Wirkort hergibt.** Der Satz
  „künftige **und historische** Changes … tragen den ausgeschlossenen Wert
  nicht mehr" ist für den Historien-Teil unter `ADR-0059` Teilfrage 3
  Option D unerfüllbar: die Filterung liegt in der Row-Image-Konstruktion,
  `cdc.changes` ist eine reine Projektion einer persistierten Change.
  Konsequenz: der Reviewer stufte das als Finding gegen den **Text** ein
  (kein Implementer-Rückweg), der Planner hat §1 auf den belegbaren
  Wortlaut korrigiert.
- **Die Platzierung des neuen Rundlauf-Abschnitts wanderte.** Vom
  SQL-Administrations-Block hinter das Go-E2E-Tier — die dedizierte
  Tabelle wird selbst erst über `cdc.enable_table` aktiviert, und kein
  Go-Testfall dieses Tiers liest sie; damit bleibt der Abschnitt
  vollständig vor der Container-Ende-Grenze der beiden Schema-Negative-
  Funktionen und trägt keine Zustands-Überschneidung.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

**Feststellung des Lese-Schritts: kein Registereintrag hat *mit dieser Welle*
die 3×-Schwelle erreicht** — es gibt in dieser Closure keinen
Steering-Loop-Eintrag zu verkörpern, und damit auch keinen Planner → Architect
→ Planner-Zug (Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b). Das ist
das Ergebnis der Prüfung, keine Auslassung; die geprüften Stände stehen unter
§Beobachtungs-Register. Insbesondere:

- `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` steht bei **1×** mit dem
  Ausgang **`geplant`** (verkörpert wird er nicht hier, sondern durch die
  Umsetzung; Träger ist `slice-075`) — der Ausgang ist also schon zugewiesen,
  und die Welle trägt ihn nicht.
- Über der Schwelle, **Ausgang unverändert `verkörpert`** — diese Welle hat
  nur Belege hinzugefügt, keinen neuen Übertritt erzeugt, also auch keinen
  neuen Lese-Schritt ausgelöst:
  `BEO-PGC/d-migrate-nacharbeit` (`slice-066`, **6×**),
  `BEO-PGC/slice-chronik-in-code-kommentar` (`slice-066`, **6×**) und
  `BEO-PGC/report-nackte-id-ohne-link` (`slice-068`, **5×** — der
  Beleg-Nachtrag war in dessen `state.md` ausdrücklich auf die
  `slice-068`-Closure vertagt und ist mit ihr fällig geworden).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`.
In dieser Welle **neu angelegt**: `laufzeitzustand-ohne-dauerhaften-traeger`
(Ausgang `geplant`, Träger `slice-075`; Beleg
`evidence/review-slice-067.md`, angelegt vom Architect-Zug — Kennung bewusst
als **Klasse** formuliert, nicht als Instanz) und
`slice-pfad-als-link-in-berichten` (2×, Ausgang noch nicht zugewiesen,
unter der Schwelle). In dieser Welle **fortgeschrieben, weiter über der
Schwelle, Ausgang unverändert** — je ein Beleg, kein neuer Übertritt:
`d-migrate-nacharbeit` (6×, `evidence/slice-066.md` — vierte Objektklasse
„CHECK-Klausel an bestehender Tabelle") und
`slice-chronik-in-code-kommentar` (6×, `evidence/slice-066.md`).
**Durchgesehen ohne Zähler-Änderung** (über die §8-Sichtungen der drei
Slices): `adapter-fehler-ausgang` (2×, §6-Risiko in `slice-066` entfallen),
`rollen-verdrahtung` (aufgelöst seit `slice-023`),
`verwaltung-keine-sql-administration` (verkörpert seit `welle-12`),
`schema-evolution-nicht-dynamisch` (2×, aufgelöst seit `slice-033`),
`test-integration-retention-timing-flake` (1×),
`test-isolation-geteilter-zustand` (1×).
Unverändert, nicht von dieser Welle berührt: alle übrigen Registereinträge
aus `welle-17` und früher.

**Benannt, nicht gezählt — ein Punkt dieser Welle, ohne eigenen
abgeschlossenen Vorgang:** Die Verweisform auf den **Welle**-Plan.
`verify-slice-068.md` und
`architect-verdict-spaltenausschluss-dauerhaftigkeit.md` adressierten
`welle-18` als Markdown-Link mit festem Verzeichnis, ebenso die
`**Welle:**`-Kopfzeile der drei Slice-Pläne dieser Welle. Diese Closure hat
die Fälle aufgelöst (Kennungs-Zitierung bzw. ein Verzeichnis tiefer), so
dass kein Gate real rot lief. Sie sind der **verwandten** Klasse von
`BEO-PGC/slice-pfad-als-link-in-berichten` zuzuordnen, aber **kein neuer
Beleg**: die registrierte Beobachtung handelt von **Slice**-Plänen in
Berichten und Entscheidungen (2×, `slice-068`,
`architect-verdict-spaltenausschluss-dauerhaftigkeit`), während hier ein
**Welle**-Plan adressiert ist und die beiden betroffenen Berichte bereits
gezählten Vorgängen angehören. Träte die Klasse bei einem künftigen
Übergang erneut aus einem **Bericht** auf, wäre sie belegfähig.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keiner, der diese Welle fortsetzt — `welle-18` schließt vollständig mit ihren
drei Slices (`slice-066`, `slice-067`, `slice-068`). Der einzige genannte
Folge-Slice ist **`slice-075`** (dauerhafter Träger des Ausschlussstandes) aus
dem Architect-Verdikt zu `slice-067`; er ist **wellenlos** und ausdrücklich
**nicht Teil dieser Welle** (`welle-18` §6). Er liegt als Datei in
`docs/plan/planning/open/` — Folge-Slice-Paarung damit getragen.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices (`slice-066`, `slice-067`, `slice-068`) liegen in
  `docs/plan/planning/done/`.
- `make gates` grün — Planner-Lauf zur Closure, Exit-Code ungefiltert
  ermittelt und in einem eigenen Schritt geprüft, nicht durch eine Pipe
  maskiert (`AGENTS.md` §3.9): `baseline-verify` OK (v6.5.0, 54 Dateien),
  d-check 536 Dateien/0 Befunde, d-check `commits`-Modul über
  `HEAD~5..HEAD` 0 Befunde, `commit-traceability` OK (5 Commits, Betreffe
  ohne Struktur-ID), `a-check` gesamt 0 Befunde, `coverage-gate` OK —
  Coverage 45,80 % über der geltenden 35-%-Schwelle.
- **`welle-18`s Closure-Trigger (c)** — der reale E2E-Beleg aus `slice-068`
  ist grün, und er ist das *Mehr* gegenüber den Einzel-DoDs: `make
  test-integration` Exit 0 (voller Compose-Stack-Lauf, beide neuen
  Beleg-Zeilen in der Ausgabe), SQL-Antrag → Live-Reload → gefiltertes Row
  Image am **laufenden** Feed-Container, ohne Neustart. Dieser Beleg liegt
  **doppelt unabhängig** vor und wird hier zusammengefasst statt erneut
  erzeugt:
  - **Implementer/Reviewer** (`slice-068` §2/§5, `welle-18` §3): `make
    test-integration` Exit 0; Implementer und Reviewer haben je eigene
    Mutationen gesetzt und rot gesehen.
  - **Verifier** (`docs/reviews/verify-slice-068.md` §2/§3): eigener,
    unabhängiger `make test-integration`-Lauf Exit 0, dazu zwei **eigene**
    Mutationen (Happy Path und Negative) — beide erwartet rot,
    Rücknahme über `git checkout` mit unverändertem Blob-Hash belegt.
- **Trigger-Audit der Welle** (Modul 6 §Wellen-Closure-Prozedur, Schritt 2
  — drei Artefaktklassen):
  - **Carveouts:** `docs/plan/carveouts/` enthält ausschließlich
    `.gitkeep` — kein offener Carveout im Repo, keine roten Gates.
  - **Bootstrap-aware Gates:** keines ist **Teil dieser Welle** — diese Welle
    führt keine neue Reifestufe ein, und ihr Umfang verbessert die
    Coverage-Zahl nicht (die neuen Belege liegen im `test/integration`-Tier
    bzw. in `tools/harness/`, beide außerhalb von `./internal/...`+`./cmd/...`,
    `harness/sensors/coverage-gate.md` §Vertrag). Geprüft ist damit dieselbe
    wellen-skopierte Lesart wie in `welle-15`/`welle-16`/`welle-17`.
    **Feststellung, nicht Auslassung:** der deklarierte Hochschalt-Trigger
    des `coverage-gate` („nächste Coverage-Verbesserung schließt die Lücke
    zur nächsten 5-%-Stufe", `harness/sensors/coverage-gate.md`
    §Kalibrierungs-Bindung) ist unter einer **repo-weiten** Lesart fällig —
    der real gemessene Ist-Stand (45,80 %) liegt über der nächsten
    5-%-Stufe (40 %) der geltenden 35-%-Schwelle, und `THRESHOLD` steht seit
    `slice-049` unverändert. Diese Entscheidung ist der **Reifestufen-Zweig**
    (Modul 8 §Rollen-Sequenz für eine Welle, Schritt 2: Planner → Architect →
    Planner) und liegt damit beim Architect — sie gehört nicht in diese
    Planner-Closure und wird hier als offener Träger benannt.
  - **ADRs mit Re-Evaluierungs-Trigger:**
    [`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) trägt zwei
    Trigger; **keiner ist eingetreten** — die technische Brücke für
    synchrone SQL→Go-Aufrufe (FDW/`dblink`) existiert nicht, und für einen
    quellen-/musterweiten Spaltenausschluss (Teilfrage 2) gibt es weder
    einen Registereintrag auf der 3×-Schwelle noch eine geschärfte
    `LH-QA-SEC-004`. [`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
    trägt zwei Trigger; **keiner ist eingetreten** — es gibt keine
    Alters-Retention/Archivierung/Löschung auf `cdc.administration_request`
    (die Antrags-Historie ist damit weiterhin tragend), und weder
    `CDC_TABLES` noch die optionale YAML-Konfigurationsdatei hat ein Feld
    für initial ausgeschlossene Spalten bekommen.
    [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
    (Coverage-Gate) — keiner seiner Re-Evaluierungs-Trigger ist eingetreten
    (kein Linter, Ist-Stand nicht bei 80 %, `LH-QA-PER-004` nicht ausgebaut);
    der oben benannte Reifestufen-Träger steht daneben und ist kein
    ADR-Trigger.
    **Feststellung: kein ADR-Trigger dieser Welle ist fällig, keine
    Nacharbeit nötig.**
- Drei Paarungen (Anker · Folge-Slice · Register) — geprüft **nach** dem
  `git mv` dieser Welle-Plan-Datei nach `done/`:
  - **Anker** — diese Closure führt **keinen** Steering-Loop-Eintrag mit
    `liegt in`-Feld; nichts zu prüfen (die zwei vorhandenen, bereits
    verkörperten Klassen sind oben als *benannt, nicht gezählt* geführt).
  - **Folge-Slice** — `slice-075` existiert real als Datei in
    `docs/plan/planning/open/` (siehe §Folge-Slices) — grün.
  - **Register** — alle in dieser Welle berührten bzw. genannten
    Verzeichnisse (`d-migrate-nacharbeit`, `slice-chronik-in-code-kommentar`,
    `laufzeitzustand-ohne-dauerhaften-traeger`,
    `slice-pfad-als-link-in-berichten`, `adapter-fehler-ausgang`,
    `test-integration-retention-timing-flake`,
    `test-isolation-geteilter-zustand`) existieren mit nicht leerem
    `evidence/` — grün.

## Archivierung

Feststellung: das Repo führt weiterhin **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die drei Slice-Dateien dieser Welle, ihre
Review-/Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in
`done/`. Dieselbe Feststellung traf bereits `welle-14`/`welle-15`/
`welle-16`/`welle-17`.
