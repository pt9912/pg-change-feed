# Review-Report: slice-016 — 2026-09-11

**Review-Art:** Code — geprüft gegen Slice-Plan +
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
(Maintainability).

**Gegenstand:** Commits `b86cf70`
(`feat(schema): d-migrate 1.3.1 Pin + Views-Retirement (ADR-0043)`) und
`4367809`
(`docs(schema): schema.yaml-Kommentare auf Zustand statt Chronik gekürzt
(ADR-0043)`) — `Makefile`, `tools/schema/schema.yaml`,
`tools/schema/nacharbeit-views.sql` (gelöscht), `tools/schema/plan.yaml`,
`tools/schema/down.sql`, `harness/README.md`, sowie der Slice-Plan selbst
(`docs/plan/planning/in-progress/slice-016-d-migrate-1.3.1-views-retirement.md`).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-11

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-016-d-migrate-1.3.1-views-retirement.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4–6 Trigger/Risiken, §8
  Sub-Area-Prüfung)
- [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (Schemamigrationen mit d-migrate; Re-Evaluierungs-Trigger: kein
  Werkzeugwechsel nötig, solange eine Ausweichform existiert)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
  (Register-Stand laut Slice-Plan §8: 3× — bereits verkörpert seit
  slice-015; die Views-Hälfte löst sich mit diesem Slice technisch auf)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.3 Move/Inhalt-Trennung,
  §3.7 Kommentar-Klassen) · `harness/conventions.md` (MR-000, genau eine
  Sub-Area `PGC`, Greenfield)
- Commit-Traceability (`AGENTS.md` §5, `ADR-0045`)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-015.md`
  (F-2, LOW: stale Kopfkommentar-Querverweis in `schema.yaml` nach
  Teiländerung) und `docs/reviews/verify-slice-015.md` (V-2, MEDIUM: drei
  `nacharbeit-*.sql`-Dateien mit toter „s. o."-Querreferenz auf die
  damals gelöschte `nacharbeit-operation-check.sql`)

---

## Findings

### F-1 — Drei `nacharbeit-*.sql`-Dateien tragen erneut eine tote „s. o."-Querreferenz — jetzt auf die von diesem Slice gelöschte `nacharbeit-views.sql`, 3. Auftreten derselben Klasse

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen — ein Kommentar beschreibt
  einen abwesenden Text/ein nicht mehr existierendes Ziel); Reviewer-Skill
  §HIGH „Kommentar trägt keine der Kommentar-Klassen"
- `pfad`: `tools/schema/nacharbeit-roles.sql:1-2`,
  `tools/schema/nacharbeit-observability.sql:1-2`,
  `tools/schema/nacharbeit-heartbeat.sql:1` — jeweils Kopfkommentar „…
  (`ADR-0043`, s. o. `tools/schema/nacharbeit-views.sql`): …"
  (`nacharbeit-heartbeat.sql` referenziert zusätzlich
  `nacharbeit-observability.sql`, das bleibt gültig)
- `befund`: Der Diff löscht `tools/schema/nacharbeit-views.sql` (Zweck des
  Slice), lässt aber drei andere, im Diff nicht berührte
  `nacharbeit-*.sql`-Dateien mit einem Kopfkommentar stehen, der genau
  diese Datei als „s. o." (Referenzpunkt für die
  „Berichtete manuelle Nacharbeit"-Rahmung) zitiert. Ein Leser, der dem
  Verweis folgt, findet nichts — die Datei existiert nicht mehr. Dieselbe
  Grund-Klasse trat bereits zweimal auf: `review-slice-015.md` F-2 (LOW,
  innerhalb derselben Datei) und `verify-slice-015.md` V-2 (MEDIUM,
  dieselben drei Dateien, damals mit Ziel `nacharbeit-operation-check.sql`)
  — dort wurde der Verweis offenbar auf `nacharbeit-views.sql`
  nachgezogen, jetzt bricht er am selben Muster ein drittes Mal. Kein Gate
  fängt das: `make docs-check`/d-check prüft laut `harness/README.md`
  §Sensors kaputte Referenzen in der Markdown-Doku, keine
  Kommentar-Querverweise innerhalb von `.sql`-Dateien (`make docs-check`
  lief in dieser Review-Sitzung grün, 168 geprüfte Dateien, 0 Befunde —
  bestätigt die Lücke).
- `verifizierbar`: nein — kein Gate prüft `.sql`-Kommentar-Querverweise.
- `klasse`: Stale Querverweis nach Löschung des Referenzziels (3.
  Auftreten — 1.: review-slice-015 F-2, LOW; 2.: verify-slice-015 V-2,
  MEDIUM, exakt dieselben drei Dateien mit anderem Zielnamen). Nach
  Baseline-Regelwerk `modul-08-agentenrollen.md` §Konflikt-Pfad
  (Drei-Auftreten-Schwelle) und `modul-10-review-harness.md` §Pflege ist
  das dritte Auftreten derselben Klasse der Punkt, an dem Klassifikation,
  ADR-/`AGENTS.md`-Schärfung und Gate-Bedarf geprüft werden — hier
  zusätzlich hochgestuft auf HIGH gegenüber den beiden Vorläufern, weil
  sich das identische Muster trotz zweifacher vorheriger Meldung ein
  drittes Mal reproduziert hat, an denselben drei Dateien.

---

## Negativbefunde

- geprüft, ohne Befund: **Pin-Änderung** — `Makefile:73` ändert
  ausschließlich den `D_MIGRATE_IMAGE`-Digest (64-Hex-`sha256`, formal
  valide, real gezogen und in dieser Sitzung gegen `schema-validate` und
  `schema-rollout` erfolgreich verwendet); keine Nebenänderung an
  `SCHEMA_SOURCE`/`SCHEMA_TARGET`/`SCHEMA_ROLLOUT_NETWORK`.
- geprüft, ohne Befund: **`schema.yaml`-Views-Deklaration inhaltsgleich
  zur gelöschten `nacharbeit-views.sql`** — Feldnamen, Reihenfolge und
  Typen der drei `columns:`-Signaturen (`active_tables`,
  `consumer_status`, `changes`) stimmen mit dem vorherigen
  `CREATE OR REPLACE VIEW`-SQL exakt überein (Spaltentypen gegen die
  physische Katalogform geprüft: `old_data`/`new_data` als `jsonb`
  entspricht der generierten `"old_data" JSONB`-Spalte, `committed_at`
  als `"timestamp with time zone"` entspricht `TIMESTAMP WITH TIME ZONE`
  aus der Tabellen-DDL); die drei `query:`-Blöcke sind textuell identisch
  zur gelöschten Datei, nur um das schema-qualifizierende `cdc.`-Präfix im
  YAML unverändert beibehalten.
- geprüft, ohne Befund: **§3 Plan-Nachzug vollständig** — jede im Diff
  geänderte Datei (`Makefile`, `schema.yaml`, `nacharbeit-views.sql`
  gelöscht, `plan.yaml`, `down.sql`, `harness/README.md`) ist in §3 der
  Slice-Plan-Datei mit Änderungsart und Begründung gelistet;
  `plan.yaml`/`down.sql` sind korrekt als „Plan-Nachzug … Nebenprodukt
  des Sensor-Laufs" gekennzeichnet.
- geprüft, ohne Befund: **DoD-Punkt 1 real erfüllt** — `make
  schema-validate` lief in dieser Sitzung gegen den neuen Pin grün (8
  Tabellen, 29 Spalten, 1 Index, 7 Constraints, „Validation passed").
- geprüft, ohne Befund: **DoD-Punkt 2 real erfüllt, beide Szenarien
  nachvollzogen** — eigene Reproduktion (Docker, PostgreSQL 18,
  `wal_level=logical`, Compose-Init-Skript): `make schema-rollout` gegen
  eine leere DB (Erstanlage) lief Exit 0, alle drei Views als `CreateView`
  gerendert. Ein zweiter `make schema-rollout`-Lauf gegen dieselbe,
  bereits migrierte DB rendert die drei Views korrekt als `ReplaceView`
  (kein `VIEW_SIGNATURE_UNKNOWN`) — die eigentliche DoD-Zusage („kein
  ReplaceView-Blocker") ist damit real bestätigt, nicht nur behauptet.
- geprüft, ohne Befund: **Nebenfinding-Behauptung „keine Regression
  gegenüber Pin 1.3.0" real geprüft** — eigene Reproduktion in einem
  separaten `git worktree` auf `b86cf70^` (alter Pin
  `sha256:d8dc38c…`, `schema.yaml` ohne `views:`-Knoten,
  `nacharbeit-views.sql` noch vorhanden): derselbe zweite
  `make schema-rollout`-Lauf gegen eine bereits migrierte DB blockiert
  dort ebenfalls mit Exit 8
  (`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`, `DropView` auf
  `cdc.heartbeat`/`cdc.metrics`) — der Blocker existierte unabhängig vom
  `views:`-Knoten und vom Pin bereits vorher. Die Implementer-Behauptung
  ist damit bestätigt: keine Regression durch diesen Slice. **Aber:**
  Dieser reale, reproduzierbare Blocker steht weder als Ausschluss in §1
  noch als Risiko in §6 des Slice-Plans, und es existiert kein
  Beobachtungs-Register-Eintrag dafür — reines Transparenz-/
  Dokumentationsdefizit, kein DoD-Blocker (`make test-integration` räumt
  die Umgebung laut Makefile-Kommentar vor jedem Lauf ab und trifft den
  Fall dadurch nie); für die Closure als Hinweis vermerkt, nicht als
  eigenständiges Finding gezählt, weil der aktuelle Diff daran nichts
  ändert und nichts verspricht, das dem widerspricht.
- geprüft, ohne Befund: **Hard Rule 3.7 in `tools/schema/schema.yaml` nach
  beiden Commits** — der zweite Commit (`4367809`) entfernt sämtliche
  bare `slice-NNN`-Nennungen (`slice-004`, `slice-012`, `slice-013`) und
  die „Retired (slice-015): … lebte bis d-migrate 1.2.0 als Ausweichform …
  Mit d-migrate 1.3.0 … real gegen einen frischen Rollout getestet
  (slice-015), kein E5/E012 mehr"-Chronik sowie den „Benannte Grenze — die
  drei Views: … Post-execute-Drift (Exit 5) … Ausweichform weiter nach dem
  Re-Evaluierungs-Trigger …"-Absatz vollständig; verbleibende Verweise
  (`Herkunft: docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit`)
  sind reine Rang-Zeiger ohne Chronik-Sprache. Eigene Volltextprüfung der
  finalen Datei: keine Reste von Forensik-Sprache
  („real erneut geprüft", „derselbe Post-execute-Drift wie …", Exit-Code-
  Nennungen vergangener Läufe) mehr vorhanden.
- geprüft, ohne Befund: **Kommentar-Klassen im `Makefile`
  (neu formulierter Views-Absatz)** — folgt demselben Muster, das
  `review-slice-015.md` bereits für das CHECK-Segment als konform
  bestätigt hat (Zustand + Rang-Zeiger auf ADR/Slice/Beleg, kein
  Konjunktiv über eine verworfene Alternative, kein abgebrochener Satz);
  „Changelog: fehlendes `ViewDefinition.sourceDialect`…" und „real …
  getestet (slice-016, Exit 0 in beiden Fällen)" sind Herkunfts-Anker im
  Dienst der aktuellen Zusage, keine Chronik-Erzählung über einen
  verworfenen Ansatz.
- geprüft, ohne Befund: **`harness/README.md`-Zeile inhaltlich korrekt**
  — behauptet weder, dass die Views weiterhin Ausweichform sind, noch
  enthält sie eine tote Referenz auf die gelöschte `nacharbeit-views.sql`;
  „seit slice-015" (Test-Kadenz-Regel) ist konsequent durch „seit
  slice-016" ersetzt, der Herkunfts-Anker-Konvention aus Modul 6
  entsprechend (wellenloser Slice → `seit slice-<NNN>` statt `seit
  welle-<NN>`); Aussagen zu Exit-0-Ergebnissen für beide getesteten
  Szenarien decken sich mit der eigenen Reproduktion oben.
- geprüft, ohne Befund: **Vollständigkeit der vier verbleibenden
  `nacharbeit-*.sql`-Dateien** — `nacharbeit-roles.sql`,
  `nacharbeit-observability.sql`, `nacharbeit-heartbeat.sql` sind
  inhaltlich unverändert und bleiben korrekt im `schema-rollout`-Target
  verdrahtet (Reihenfolge Rollen → Observability → Heartbeat unverändert);
  nur ihr Kopfkommentar trägt den in F-1 gemeldeten toten Verweis.
- geprüft, ohne Befund: **Hard Rule 3.3 (git mv + Inhalt = zwei Commits)**
  — kein `git mv` in diesem Diff (Slice bleibt in `in-progress/`, kein
  Lifecycle-Übergang); `nacharbeit-views.sql` wird gelöscht, nicht
  verschoben — ein Commit mit Inhaltsänderung ist hier korrekt. Der zweite
  Commit (`4367809`, reine Kommentar-Korrektur an derselben Datei) ist
  ebenfalls kein Move und braucht keine Trennung von einem `git mv`.
- geprüft, ohne Befund: **Docker-only** — alle Läufe (`schema validate`,
  `schema migrate`, psql-Nacharbeit) laufen über `docker run` im
  Makefile; kein lokales Toolchain-Install im Diff.
- geprüft, ohne Befund: **Suppression-Verbot** — kein
  `nolint`/`noqa`/`SuppressMessage`/`d-check:ignore`-Marker im Diff.
- geprüft, ohne Befund: **Traceability der Commit-Messages** — beide
  Commits tragen `ADR-0043` im Betreff, keine `SPEC-*`/`ARC-*`-Kennung im
  Betreff; `make commit-traceability` lief in dieser Sitzung gegen die
  Range `b86cf70^..4367809` grün.
- geprüft, ohne Befund: **`make gates`** — in dieser Sitzung erneut
  ausgeführt (`baseline-verify`, `docs-check`, `commit-traceability`,
  `a-check`): alle grün, 0 Befunde.
- geprüft, ohne Befund: **§6 Risiken vorab benannt und real bewertet** —
  beide Risiken (Metadaten-Lücke jenseits `columns:`/`source_dialect`;
  unerwarteter Regressions-Fund durch den Digest-Bump) sind mit diesem
  Diff *nicht* eingetreten — die eigene Reproduktion bestätigt Exit 0 für
  beide DoD-Szenarien ohne neuen Blocker; die formale Ausgangs-Zuweisung
  selbst ist Planner-Aufgabe bei Closure, hier nicht bewertet.
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung und Register-Sichtung** —
  GF-Sub-Area `PGC` korrekt referenziert; Register-Stand zum
  Planungszeitpunkt (zehn Einträge, `d-migrate-nacharbeit` 3× bereits
  verkörpert) korrekt wiedergegeben; kein Eintrag erreicht mit diesem
  Slice neu die 3×-Schwelle — die in F-1 gemeldete Klasse („Stale
  Querverweis nach Löschung des Referenzziels") stand zum Planungszeitpunkt
  noch nicht im Register (sie entsteht erst über die Review-/
  Verify-Reports von slice-015/slice-016) und konnte dort folgerichtig
  nicht auftauchen; das ist kein Fehler der §8-Sichtung selbst.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Ein HIGH-Finding (F-1).** Kein ADR-Verstoß, kein Traceability-/
ID-Schema-Verstoß, keine Suppression, kein Docker-only-Verstoß, keine
Zwei-Quellen-Drift zwischen Plan und Code. F-1 trifft die HIGH-Klasse
„Kommentar trägt keine der Kommentar-Klassen" (abwesenden
Referenz-Text), verschärft durch das dreifache Auftreten derselben
Grund-Klasse an denselben drei Dateien.

**Finding-Klassen dieses Laufs:** Stale Querverweis nach Löschung des
Referenzziels (3. Auftreten)

## Verdikt

**Merge-blockierend:** ja — F-1 ist ein offenes HIGH-Finding; der
Closure-Trigger des Slice-Plans (§5: „Review-Schluss ohne offenes
HIGH-Finding") ist damit noch nicht erfüllt.

**Kein Rollen-Widerspruch.** F-1 ist ein neuer Befund dieser
Review-Sitzung, vom Implementer bislang nicht kommentiert oder
bestritten — der Konflikt-Pfad aus Modul 8 (Sequenz mit
Übergabe-Artefakten über den Architect) ist dadurch *nicht* ausgelöst;
er würde erst greifen, wenn der Implementer diesem HIGH mit einer
Sachbehauptung widerspricht. Unabhängig davon ist die
Drei-Auftreten-Schwelle für die Finding-*Klasse* selbst erreicht (Modul 8
§Konflikt-Pfad nennt sie als Auslöser für die Sequenz-Pflicht; Modul 10
§Pflege nennt sie als Auslöser für Klassifikations-Schärfung/Gate-Bedarf)
— das ist eine andere Schwelle als der Rollen-Widerspruch und bezieht
sich auf die Wiederholung des Musters, nicht auf einen Dissens zwischen
Rollen.

**Zur eigentlichen technischen Frage (Pin-Bump, Views-Retirement,
Nebenfinding-Regression):** unproblematisch und real verifiziert in
dieser Sitzung — beide DoD-Szenarien (Erstanlage, Folgelauf mit
`ReplaceView`) laufen Exit 0 gegen den neuen Pin; die
Nebenfinding-Behauptung „keine Regression" ist durch eigene Reproduktion
gegen den alten Pin bestätigt. Die Views-Deklaration in `schema.yaml` ist
inhaltlich deckungsgleich mit der gelöschten Ausweichform. Das einzige
offene Problem dieses Diffs ist die redaktionelle Querreferenz in drei
unveränderten Nachbardateien (F-1).

**Übergabe:** F-1 geht an den Implementer zur redaktionellen Korrektur
der drei Kopfkommentare (`nacharbeit-roles.sql`,
`nacharbeit-observability.sql`, `nacharbeit-heartbeat.sql`) — lösbar ohne
neue Entscheidung, kein Rollen-Widerspruch. Zusätzlich empfiehlt sich für
die Slice-Closure (§7), die Finding-Klasse „Stale Querverweis nach
Löschung des Referenzziels" wegen des dritten Auftretens als
Steering-Loop-Kandidaten zu behandeln (Planner-Entscheidung: Register-
Eintrag, geschärfte Regel in `AGENTS.md`/Skill, oder ein Sensor auf
`.sql`-Kommentar-Querverweise) — das ist keine Entscheidung dieses
Reports. Die **Finding-Klasse** geht zusätzlich in die Slice-Closure §7
und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg**
und wird über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).

---

**Gate-Beleg:** `make gates` in dieser Review-Sitzung real ausgeführt
(grün, 0 Befunde); `make schema-validate` und zwei eigene
`make schema-rollout`-Läufe (Erstanlage + Folgelauf, Docker/PostgreSQL 18,
netzgebunden) real ausgeführt, inklusive einer Vergleichs-Reproduktion
gegen den alten Pin (`b86cf70^`, separater `git worktree`) für die
Nebenfinding-Prüfung. Arbeitsverzeichnis nach allen Läufen wieder sauber
(`tools/schema/plan.yaml`/`down.sql` auf den committeten Stand
zurückgesetzt, Test-Container/-Netz entfernt, Worktree entfernt).
