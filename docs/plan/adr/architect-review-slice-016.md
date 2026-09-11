# Architect-Review slice-016 — Verdikt zu Trigger-Audit (ADR-0043), Register-Bestätigung (BEO-PGC/d-migrate-nacharbeit, 4×) und Neuanlage-Prüfung (BEO-PGC/schema-rollout-fremdobjekte, 1×)

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-11.
**Eingang:** Slice-Plan
[`docs/plan/planning/done/slice-016-d-migrate-1.3.1-views-retirement.md`](../planning/done/slice-016-d-migrate-1.3.1-views-retirement.md)
§1–§8 · [`docs/reviews/review-slice-016.md`](../../reviews/review-slice-016.md)
(F-1 HIGH disponiert in `04590a4`, danach 0 HIGH) ·
[`docs/reviews/verify-slice-016.md`](../../reviews/verify-slice-016.md)
(DoD-Konformität bestätigt, VF-1/VF-2 LOW für Planner-Closure) ·
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (Accepted, permanent,
Re-Evaluierungs-Trigger) ·
`docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
(`observation.md`, `state.md`, `evidence/{slice-006,slice-010,slice-015,slice-016}.md`)
· `docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`
(`observation.md`, `state.md`, `evidence/slice-016.md`) · Baseline-Regelwerk
`modul-06-roadmap.md` §Das Beobachtungs-Register ·
`modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle (Tabelle „ohne
Wellen-Betrieb": Trigger-Audit bleibt, bei jeder Slice-Closure statt einmal
pro Welle — dieser Slice ist wellenlos).

**Ausgang:** Drei unabhängige Prüfungen, alle mit Verdikt:

1. **Trigger-Audit (ADR-0043):** Der Re-Evaluierungs-Trigger feuert
   **nicht**. ADR-0043 bleibt `Accepted`, `permanent`, unverändert, kein
   Folge-ADR fällig.
2. **Register-Bestätigung (BEO-PGC/d-migrate-nacharbeit, jetzt 4×):** Die
   `state.md`-Formulierung des Planners gibt den vollständig aufgelösten
   technischen Stand korrekt wieder, **ohne** die bereits verkörperte Regel
   (seit slice-015) zu widersprechen. Kein neuer Ausgang fällig, keine
   Korrektur nötig — bestätigt wie vorgelegt.
3. **Neuanlage-Prüfung (BEO-PGC/schema-rollout-fremdobjekte, 1×):**
   Sub-Area-Zuordnung und `observation.md`-Formulierung sind sauber. Ein
   **nicht-blockierender** Formhinweis zu `state.md` (Vokabular-Vermischung,
   siehe Zug 3) — keine Entscheidung nötig, keine Auswirkung auf diese
   Closure.

**Harte Regel eingehalten:** ADR-0043 (`Accepted`) wird von diesem Lauf
**nicht** inhaltlich geändert — der Trigger-Audit bestätigt sie nur. Weder
der Slice-Plan noch das Beobachtungs-Register werden von diesem Lauf
editiert — das ist Planner-Arbeit im Closure-Zug (Modul 8, Schritt 3a/3c:
der Architect liefert das Audit-Verdikt, der Planner trägt es ein).

---

## Zug 1 — Trigger-Audit (ADR-0043)

**Trigger-Text (wörtlich):** „d-migrate kann eine benötigte Operation
nicht ausdrücken — der Lauf blockiert (Exit 8) oder meldet einen Blocker,
**und** es gibt keine Ausweichform (Overlay, `--rename-*`, berichtete
manuelle Nacharbeit) — dann ist der Werkzeugwechsel als Folge-ADR mit
`supersedes` zu prüfen. Sonst `permanent`."

Der Trigger ist eine **Konjunktion**: beide Teilbedingungen müssen
gleichzeitig gelten — *nicht ausdrückbar* **und** *keine Ausweichform*.
Fehlt eine der beiden, feuert der Trigger nicht. Geprüft wird jeder
verbliebene `nacharbeit-*.sql`-Fall einzeln, nicht nur die zwei historisch
unter ADR-0043 geführten:

| Fall | Ausdrückbar? | Ausweichform vorhanden? | Trigger-Bedingung erfüllt? |
|---|---|---|---|
| `chk_change_operation` (CHECK, `raw-sql-text-drift`) | **ja, seit 1.3.0** (slice-015, deklarativ in `schema.yaml`) | entfällt — `nacharbeit-operation-check.sql` bereits gelöscht (slice-015) | **nein** — Fall gelöst |
| Drei Views (`active_tables`, `consumer_status`, `changes`) | **ja, seit 1.3.1** (dieser Slice: `source_dialect`+`columns:`, real dreifach reproduziert — Erstanlage und Folgelauf, Exit 0, Gegenprobe ohne `columns:` reproduziert den vorhergesagten `VIEW_SIGNATURE_UNKNOWN`) | entfällt — `nacharbeit-views.sql` bereits gelöscht (dieser Slice) | **nein** — Fall gelöst |
| Rollen-DDL (`nacharbeit-roles.sql`) | **kein Ausdrückbarkeits-Fall** — `CREATE ROLE` ist kein Tabellen-/View-Objekt und liegt außerhalb des Gegenstandsbereichs von `schema.yaml` (Kopfkommentar der Datei bestätigt das explizit); es gibt hier keine „benötigte Operation", die d-migrate ausdrücken soll | ja — DO-Block, unverändert seit slice-011, stabil in Betrieb | **nein** — Konjunktion greift nicht (erste Teilbedingung ist kategorial nicht einschlägig) |
| `cdc.metrics`/`cdc.heartbeat` (Views, `nacharbeit-observability.sql`/`nacharbeit-heartbeat.sql`) | **technisch ausdrückbar** — dieser Slice beweist gerade, dass d-migrate Views mit `source_dialect`+`columns:` deklarativ trägt; es ist **nicht dokumentiert**, dass eine Überführung an einer Werkzeug-Grenze scheitert, sondern schlicht **noch nicht angegangen** (slice-011-Abgrenzung, „anderer Vorgang") | ja — psql-Nacharbeit, real in Betrieb, grantet sich selbst an `cdc_reader` | **nein** — erste Teilbedingung nicht belegt („kann nicht" ≠ „wurde noch nicht überführt") |
| Foreign-Object-Blocker (`BEO-PGC/schema-rollout-fremdobjekte`, Exit 8 `DropView` auf `heartbeat`/`metrics` bei erneutem Rollout) | **kein Ausdrückbarkeits-Fall** — das ist **beabsichtigtes** Werkzeug-Verhalten: ADR-0043 Entscheidung Punkt 3 sagt wörtlich „destruktive Operationen bleiben default blockiert (Exit 8); ihre Zulassung ist ein bewusster, berichteter Entschluss, kein Default." Der Blocker ist die Sicherheits-Vorkehrung selbst, nicht ein Scheitern an einer benötigten Operation | — (Frage stellt sich nicht) | **nein** — kategorial kein Trigger-Kandidat |

Für **keinen** der sechs geprüften Fälle sind beide Teilbedingungen
gleichzeitig erfüllt. Bestätigt damit den Reviewer- (kein HIGH gegen
ADR-0043) und Verifier-Negativbefund („Re-Evaluierungs-Trigger feuert
nicht") unabhängig, und erweitert die Prüfung explizit auf die drei
verbliebenen `nacharbeit-*.sql`-Dateien sowie den neuen Fremdobjekte-Fund,
die beide im Reviewer-/Verifier-Kontext nicht als Trigger-Frage gestellt
wurden.

**Zur Formulierung der Aufgabenstellung:** Die erwartete Einordnung — dass
Rollen kein Tabellen-/View-Objekt sind und Observability/Heartbeat schlicht
noch nicht angegangen, nicht als „keine Ausweichform" dokumentiert sind —
wird durch die Kopfkommentare der Dateien selbst bestätigt (eigene Lektüre
im Volltext, nicht nur aus dem Slice-Plan übernommen). Ergänzend: Der
Foreign-Object-Fund ist eine **dritte**, in der Aufgabenstellung nicht
genannte Kategorie — kategorial ebenfalls kein Trigger-Kandidat, weil das
Exit-8-Verhalten dort keine Werkzeug-Grenze, sondern eine in der ADR selbst
beschriebene Schutzfunktion ist.

**Verdikt Zug 1:** Bestätigt — der Re-Evaluierungs-Trigger feuert nicht.
ADR-0043 bleibt `Accepted`, `permanent`, unverändert. Kein Folge-ADR.

---

## Zug 2 — Register-Bestätigung (BEO-PGC/d-migrate-nacharbeit, jetzt 4×)

Die Regel wurde bereits mit slice-015 bei 3× verkörpert (Zielort
`harness/README.md` §Sensors, `make schema-rollout`-Zeile, Herkunftsanker
`seit slice-015`) — dieser Slice ergänzt nur den vierten Beleg
(`evidence/slice-016.md`). Modul 6 verlangt hier **keine** erneute
Ausgangs-Zuweisung: Ein weiterer Beleg zu einem bereits `verkörpert`en
Eintrag löst keinen neuen Entscheidungsbedarf aus, er bestätigt nur die
Fortdauer des Musters.

**Geprüft:** Widerspricht die aktualisierte `state.md`-Formulierung der
weiterhin geltenden Regel?

- `state.md` sagt: „Zustand: verkörpert" (unverändert) — korrekt, kein
  Rückzug der bereits geschriebenen Regel.
- „Nachrichtlich: Beide ursprünglich betroffenen Fälle sind jetzt technisch
  aufgelöst … Die verkörperte Regel bleibt als Betriebsdisziplin bestehen"
  — das ist **kein Widerspruch**, sondern trennt zwei verschiedene Dinge
  sauber, in derselben Weise, wie das Architect-Verdikt zu slice-015 es
  bereits getan hat (siehe dort, Zug 2, „Zwei verschiedene Dinge"): Die
  Regel bewertet nicht, ob die *konkreten* zwei Fälle behoben sind, sondern
  ob aus dreimaliger Wiederholung eines *Musters* (Pin-Bump testet eine
  bekannte, unveränderte Ausweichform-Grenze routinemäßig real nach, ohne
  neue Information zu liefern) eine stehende Betriebsdisziplin wurde. Ein
  zukünftiger, neuer `nacharbeit-*`-Fall (etwa eine künftige Überführung von
  `cdc.metrics`/`cdc.heartbeat`, oder ein neuer, heute unbekannter
  d-migrate-Grenzfall) würde derselben Test-Kadenz-Regel unterliegen — die
  Regel ist über die zwei jetzt gelösten Einzelfälle hinaus wirksam. Das
  `Nachrichtlich`-Feld dokumentiert korrekt den *aktuellen technischen
  Stand* (beide Ausweichformen zurückgebaut), ohne die *prozedurale* Regel
  (wann ein Real-Test fällig wird) zu berühren oder zurückzunehmen.
- Der Zähler (4×, abgeleitet aus den vier `evidence/*.md`-Dateien) ist
  korrekt — eigene Nachzählung bestätigt vier Dateien.
- `observation.md` bleibt unverändert (unveränderlich ab Anlage, trägt
  weiterhin nur den ursprünglichen CHECK-Fall als Bezeichnung) — korrekt
  nach Modul 6, die Status-Entwicklung gehört in `state.md`, nicht in
  `observation.md`.

**Verdikt Zug 2:** Die vorgelegte `state.md`-Formulierung ist korrekt und
widerspruchsfrei. Keine Korrektur nötig, keine neue Ausgangs-Zuweisung
fällig — der bei slice-015 gesetzte Ausgang `verkörpert` bleibt unverändert
gültig.

---

## Zug 3 — Neuanlage-Prüfung (BEO-PGC/schema-rollout-fremdobjekte, 1×)

**Sub-Area-Zuordnung:** `Schemamigration (d-migrate; Sub-Area-Kürzel PGC)`
— identisch zur Schwester-Beobachtung `d-migrate-nacharbeit`, korrekt: Der
Fund betrifft dieselbe Sub-Area (Modus-Deklaration führt nur die eine
Default-Sub-Area `PGC`, Greenfield, für das gesamte Repo — keine
Ausdifferenzierung nötig, der Fund liegt eindeutig in `tools/schema/`).

**`observation.md`-Formulierung:** Faktisch korrekt und mit dem
Verifier-/Reviewer-Befund deckungsgleich — Exit 8,
`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`, `DropView` auf
`cdc.heartbeat`/`cdc.metrics`, außerhalb des neutralen Modells, real gegen
1.3.0 **und** 1.3.1 reproduziert (drei unabhängige Läufe: Implementer,
Reviewer im separaten Worktree gegen den alten Pin, Verifier). Der
ADR-0043-Bezug ist zutreffend gesetzt. `evidence/slice-016.md` benennt
korrekt, warum dieser Slice den Fund **nicht** behebt (anderer Vorgang,
§1-Klasse 3 — träfe `nacharbeit-observability.sql`/`nacharbeit-heartbeat.sql`,
nicht die Views-Ausweichform dieses Slices).

**Formhinweis, nicht blockierend:** `state.md` schreibt „Zustand: offen —
Ausgang: **weiter offen**". Modul 6 definiert für das
Beobachtungs-Register eine geschlossene Menge von **drei** Ausgängen —
`verkörpert` / `geplant` / `gestrichen` — und sagt ausdrücklich: „Nur zwei
der drei hängen an der Schwelle … Unterhalb der Schwelle ist `offen` der
Stand — dort ist er **kein Ausgang**, sondern der Normalzustand." Der
Begriff „weiter offen" gehört zu einer **anderen**, ebenfalls geschlossenen
Menge — den drei Risiko-Ausgängen aus Modul 5 §6 (`eingetreten` /
`entfallen` / `weiter offen`). Bei 1× (unter der Schwelle von 3×) ist für
das Register kein Ausgang fällig; die Formulierung „Ausgang: weiter offen"
vermischt die beiden Vokabulare und könnte einen späteren Leser oder ein
künftiges Prüfwerkzeug zu der Annahme verleiten, hier sei bereits ein
formaler Register-Ausgang gesetzt worden, obwohl keiner der drei legitimen
Register-Ausgänge gemeint ist. **Empfehlung an den Planner** (keine
Bedingung für die Closure dieses Slices — der Fund gehört zu einer anderen,
eigenständigen Beobachtung, nicht zu `d-migrate-nacharbeit`, und diese liegt
unter der Schwelle): `state.md` bei Gelegenheit auf „Zustand: offen" ohne
die Zeile „Ausgang: weiter offen" kürzen, den erläuternden Satz („kein
Träger bislang; eine Lösung bräuchte …") als reine Prosa ohne Ausgangs-Label
stehen lassen.

**Verdikt Zug 3:** Neuanlage korrekt — Sub-Area und Formulierung tragen.
Ein kosmetischer Formhinweis zu `state.md` (Vokabular-Vermischung
Register- vs. Risiko-Ausgänge), keine Entscheidung, keine Auswirkung auf
die Closure von slice-016.

---

## Disposition (Gesamt)

| Frage | Ergebnis |
|---|---|
| Feuert der ADR-0043-Re-Evaluierungs-Trigger? | Nein — für keinen der sechs geprüften Fälle (CHECK, Views, Rollen, Observability/Heartbeat, Foreign-Object-Blocker) ist die Konjunktion erfüllt |
| ADR-0043-Status | unverändert, `Accepted`, `permanent` |
| Folge-ADR nötig? | Nein |
| Register-Bestätigung `BEO-PGC/d-migrate-nacharbeit` (4×) | `state.md`-Formulierung korrekt, kein Widerspruch zur bereits verkörperten Regel, keine Korrektur nötig |
| Neuer Ausgang für `d-migrate-nacharbeit` fällig? | Nein — bereits `verkörpert` seit slice-015, vierter Beleg ändert daran nichts |
| Neuanlage `BEO-PGC/schema-rollout-fremdobjekte` (1×) | Sub-Area und Formulierung korrekt; Formhinweis zu `state.md` (nicht blockierend) |
| Auswirkung auf slice-016-Closure | Keine — beide Register-Prüfungen bestätigen den vorgelegten Stand |

---

## Beleg-Anker (Kurzfassung)

| Aussage | Beleg |
|---|---|
| CHECK ausdrückbar seit 1.3.0, Ausweichform zurückgebaut | slice-015-Closure, `BEO-PGC/d-migrate-nacharbeit/evidence/slice-015.md` |
| Views ausdrückbar seit 1.3.1, Ausweichform zurückgebaut | `review-slice-016.md` Negativbefund „DoD-Punkt 2 real erfüllt"; `verify-slice-016.md` Sensor-Tabelle (Erstanlage/Folgelauf/Gegenprobe); `BEO-PGC/d-migrate-nacharbeit/evidence/slice-016.md` |
| Rollen kein Tabellen-/View-Objekt | `tools/schema/nacharbeit-roles.sql:1-8` (Kopfkommentar, eigene Lektüre) |
| Observability/Heartbeat technisch ausdrückbar, nur noch nicht überführt | `tools/schema/nacharbeit-observability.sql:1-8`, `tools/schema/nacharbeit-heartbeat.sql:1-13` (Kopfkommentare, eigene Lektüre); Slice-Plan §1 (slice-011-Abgrenzung als „anderer Vorgang") |
| Foreign-Object-Blocker ist beabsichtigtes Sicherheitsverhalten, kein Ausdrückbarkeits-Fall | [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) Entscheidung Punkt 3 („destruktive Operationen bleiben default blockiert … kein Default"); `BEO-PGC/schema-rollout-fremdobjekte/observation.md` |
| Trigger-Wortlaut (Konjunktion) | [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) §Re-Evaluierungs-Trigger |
| Register-Ausgang bereits `verkörpert` seit slice-015 | `docs/plan/adr/architect-review-slice-015.md` Zug 2 |
| Zähler bei 4× | `BEO-PGC/d-migrate-nacharbeit/evidence/{slice-006,slice-010,slice-015,slice-016}.md` |
| Register-Ausgangs-Vokabular vs. Risiko-Ausgangs-Vokabular | Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register („Nur zwei der drei hängen an der Schwelle … unterhalb der Schwelle ist `offen` der Stand … kein Ausgang") vs. `modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst |
