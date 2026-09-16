# Review-Report: slice-093 — **Delta 2** (Fixrunde `e140363` + Verifikations-Korrektur `9f9bffc`) · 2026-09-16

**Review-Art:** Delta-Review der **zweiten** Fixrunde — geprüft gegen **Plan und
Entscheidungen**, gegen die Findings des ersten Delta-Reports
(`docs/reviews/review-slice-093-delta.md`, dort D-1…D-5) und gegen die Zusagen der
beiden Commit-Texte. Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-10-review-harness.md` §Drei Review-Arten. DoD-/Spec-Konformität
(Verifier), §6-Ausgänge, Register-Zähler und die drei Paarungen sind **nicht**
Gegenstand.

**Gegenstand:**

| Pfad | `e140363` | `9f9bffc` |
|---|---|---|
| `internal/bootstrap/roles_rollout_file_internal_test.go` | **+41/−16** | **+1/−1** |
| `.dockerignore` | — | **+4/−2** |

**Nummerierung:** Die `D-N` dieses Reports zählen **neu**; Bezug auf den ersten
Delta-Report steht als „1. Delta, D-N".

**Skill:** `.harness/skills/reviewer.md` @ `HEAD` (`9f9bffc`).

---

## Findings

### D-1 (2. Delta) — Die Einelementigkeit der Ausnahme ist an eine Begründung gehängt, die sie nicht trägt

- `kategorie`: **LOW**
- `quelle`: §3.12 Instanz B · `LH-QA-SEC-001`
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:299-300`
- `befund`: Der Satz lautete „Ausgenommen ist allein der Schema-USAGE-Grant — **und
  die Ausnahme ist einelementig, weil (6a) jedes Schema-Objekt auf `USAGE`
  festhält**." Gemessen trifft die **heutige** Datei zu (jede Rolle hält genau
  **ein** Schema-Objekt, `schema cdc`, nur `[usage]`). Der **`weil`-Satz** trägt die
  Einelementigkeit aber nicht: (6a) bindet die **Privileg-Dimension**, nicht die
  **Zahl** der Schema-Objekte. Eine Datei mit zwei Schema-`USAGE`-Grants passiert
  den Test **grün**, und die Ausnahme von (7) überspringt dann **zwei** Elemente:

  ```text
  M9 GRANT USAGE ON SCHEMA public TO cdc_reader;  → --- PASS: TestRolloutDateiTraegtDieRechteDerVerdrahtung
  ```
- `verifizierbar`: **ja** — eigener Dump + M9 (`TEST-EXIT=0`), 1. Delta `D-2`-Messung
- `klasse`: „benannte Ausnahme enger als die Ausnahme im Code" — **derselbe Vorgang**
  wie F-2 und 1. Delta `D-2` (der Zähler bewegt sich nicht)

### D-2 (2. Delta) — Das Zahlwort „zwei" und „genau diese Dateien" greifen auf sechs Negationen zu

- `kategorie`: LOW
- `quelle`: §3.7 (Zahl/Zustand am Ist-Zustand) · Geschwister von `V-2` der Verifikation
- `pfad`: `.dockerignore:3-4`
- `befund`: Die von `9f9bffc` **neu geschriebene** Kopfzeile sagte „Die zwei
  Negationen unten VERGRÖSSERN den Kontext … sie heben den Ausschluss für **genau
  diese Dateien** auf". Unter dem Kommentar stehen **sechs** Negations-Zeilen
  (`grep -c '^!' .dockerignore` → **6**); welches Paar gemeint ist, sagt erst `:7-8`
  — „genau diese Dateien" ist eine **Vorwärts**-Deixis auf Namen, die zwei Sätze
  später kommen. Die **Richtungs**-Aussage selbst ist richtig (eigene Messung).
- `verifizierbar`: **ja** — `grep -c '^!' .dockerignore`; `git show 03fc53a -- .dockerignore`
- `klasse`: „Zahlwort/Deixis ohne Deckung an der Fundstelle" — dieselbe Klasse wie `V-2`

### D-3 (2. Delta) — Die Begründung von (6) nennt (7) als Träger der Trennung; gemessen trägt (6)

- `kategorie`: INFO
- `quelle`: §3.12 Instanz B
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:263-269`, Fundstelle `:267`
- `befund`: „Sie lässt die übrigen Zeilen der Datei stehen; **aufgehoben ist die
  Trennung, die (7) prüft**". Der erste Halbsatz hält; der zweite ist eine
  Attributions-Nuance: **gemessen kommt der Fatal allein aus (6)** — mit
  deaktiviertem Block (6) ist die Probe **grün**, (7) erreicht den Fall nicht:

  ```text
  (6) AUS, GRANT ALL ON ALL TABLES IN SCHEMA cdc angehängt → TEST-EXIT=0, kein Fatal
  (6) AN,  dieselbe Mutation                                → TEST-EXIT=1 (:273, Regel (6))
  ```
- `klasse`: „Regel-Begründung weiter als die Messung" — dieselbe Klasse wie 1. Delta `D-4`

### D-4 (2. Delta) — „jedes zu diesem Zeitpunkt existierende Objekt des Schemas" ist weiter als der PostgreSQL-Gegenstand

- `kategorie`: INFO
- `quelle`: §3.12 Instanz B (PostgreSQL-Behauptung) · `LH-QA-SEC-001`
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:264-266`
- `befund`: Die **Korrektur** von 1. Delta `D-4` ist gegen PostgreSQL **richtig** —
  eigene Messung (PG 17, lokaler Container): ein **nach** dem Grant erzeugtes Objekt
  ist nicht gedeckt, nach `ALTER DEFAULT PRIVILEGES` ist es gedeckt. Zu weit ist nur
  das **Subjekt**: `ALL TABLES IN SCHEMA` deckt **tabellenartige** Objekte, **nicht
  jedes Objekt** des Schemas — eine **vor** dem Grant erzeugte Sequenz bleibt
  ungedeckt (dafür bräuchte es `ALL SEQUENCES`):

  ```text
  P1 Tabelle a (vor Grant)                    | t      P3 Sequenz q                     | f
  P2 View v (vor Grant)                       | t      P4 Tabelle c (nach Grant)        | f
  P6 Tabelle d (nach ALTER DEFAULT PRIVILEGES)| t      P5 View w (nach Grant)           | f
  ```
- `klasse`: „PostgreSQL-Behauptung weiter als ihr Gegenstand"

---

## Negativbefunde

- **1. Delta `D-1` (HIGH) ist geschlossen.** Die Absent-Formen sind aus dem
  Gegenstand weg; der neue Satz steht im **Präsens** und ist **wahr** — die eine
  Zusage („ein nicht aufgeführtes, von einer schreibenden Rolle gehaltenes Objekt
  wird gefangen") ist durch eigene Mutation belegt
  (`GRANT SELECT ON cdc.table_schema TO cdc_reader` → `TEST-EXIT=1`, Regel (7)).
- **Alle sieben aufgeführten Mutationen färben wirklich rot.** 8 eigene Läufe
  (netzlos, gepinnte Toolchain): Baseline `0`, dann siebenmal `1` — jede mit **ihrer**
  Regel in der Fatal-Zeile (`:215` (1) · `:224` (2) · `:239` (3) · `:204` (0) ·
  `:259` (5) · `:273` (6) · `:289` (6a)). Die Liste im Testkopf und die Liste in §7
  sind **dieselben sieben in derselben Reihenfolge**.
- **Die neue Prüfung (6a) bindet:** `GRANT CREATE ON SCHEMA cdc TO cdc_reader` →
  `TEST-EXIT=1`; `GRANT ALL ON SCHEMA cdc TO cdc_reader` → ebenso. Die zwei
  Regressionen, die der erste Delta-Review als **grün** gemessen hatte, sind **rot**.
- **(6a) löst nicht falsch aus:** doppeltes `GRANT USAGE ON SCHEMA cdc` → `0`;
  `GRANT USAGE ON SCHEMA public TO cdc_reader` → `0`.
- **Die PG-Hälften des (6a)-Satzes stimmen** (PG 17): `GRANT CREATE ON SCHEMA` gibt
  das Schema-`CREATE` und **keinen** Objektzugriff; `GRANT ALL ON SCHEMA` = `USAGE`+`CREATE`.
- **`V-2`/`V-3` der Verifikation sind richtig und verschieben nichts:** die
  Aufzählung führt genau **drei**; die Kontext-Messung (BuildKit, `COPY . .`,
  netzlos, Kopie außerhalb des Baums) ergibt **171** Einträge **mit** der Negation,
  **170** ohne.
- **Die Adresse „§7 der Closure-Notiz von `slice-093`" löst auf** — sie zeigt auf die
  **existierende** Sektion eines existierenden Artefakts, und die Exit-Codes sind
  unabhängig nachgefahren und **stimmen**. **Kein Finding.** Die Pflicht daraus ist
  die der Verifikation (`V-4`): jene §7-Fassung gehört **mit** der Closure in den Baum.
- **1. Delta `D-5` ist geschlossen, ohne eine Slice-Nummer zu erfinden:** die
  Nicht-Anker-Form ist aus **dieser** Datei weg; `nacharbeit-roles.sql:101` unberührt.
- **Hygiene, Traceability, Umfang:** kein `//nolint`, kein host-lokaler Pfad;
  `tools/schema/`, `harness/mk/`, `Makefile`, `Dockerfile`, `THRESHOLD` von beiden
  Commits **nicht** berührt; Betreffe mit `ADR-0085`, keine Struktur-ID;
  `docs-check` 771/0.
- **Kein neuer A-Klassen-Satz.** Jeder neue/geänderte Satz wurde gegen **Zahl ·
  Anker · Deixis · §-Verweis · Chronik-/Absent-Text · PostgreSQL-Behauptung** geprüft.

## Eigene Messungen (Exit-Codes je direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git show --numstat --format="" e140363` / `9f9bffc` | 0 | +41/−16 · +4/−2 und +1/−1 |
| `git status --porcelain -- internal .dockerignore` | 0 | **leer** — Gegenstand unverändert |
| **8 LP2-Läufe** (Container, `--network none`) | Baseline **0**, Proben je **1** | **7 × rot**, jede mit ihrer Regel |
| 2 Negativ-Proben | je **0** | keine falsche Auslösung |
| 2 Block-Disabling-Proben | je **0** | (7) erreicht die Fälle nicht → **D-3** |
| Dump des geparsten Bestands (3 Rollen) | 0 | je **ein** Schema-Objekt, nur `[usage]` |
| Kontext-Messung (`COPY . .`, BuildKit) | 0 | **171** mit / **170** ohne die Negation |
| `make docs-check` | **0** | 771 Dateien, 0 Befunde |
| `make commit-traceability` | **0** | OK, 5 Commits |
| PG-17-Probe (`has_*_privilege`, 10 Fragen) | 0 | s. D-4 |
| `grep -c '^!' .dockerignore` | 0 | **6** → D-2 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **0** |
| MEDIUM | **0** |
| LOW | 2 |
| INFO | 2 |

## Verdikt

**Trägt die zweite Fixrunde? — Ja, sie trägt.** Die vier Adressen des ersten
Delta-Reports sind erledigt und **nachgemessen**, nicht übernommen: `D-1`
(Absent-Text weg, neuer Satz wahr), `D-2` (**(6a) bindet**, die zwei grünen
Regressionen sind rot), `D-3` (Adresse löst auf und trägt die sieben Exit-Codes),
`D-4` (Subjekt auf die Messung gezogen, die PG-Nebenaussage ist richtig), `D-5`
(Nicht-Anker-Form ohne erfundene Nummer entfernt). Die zwei
Verifikations-Korrekturen sind korrekt.

**Merge-blockierend: nein** — 0 HIGH, 0 MEDIUM. Die zwei LOW sind **Weiten von
Sätzen** ohne offene least-privilege-Lücke (D-1) bzw. ohne semantische Auswirkung
(D-2); die zwei INFO sind Präzisierungen.

**Rückgabe-Pfeil Reviewer → Implementer: nicht erforderlich.** Die vier Adressen,
*falls* die Runde sie mitnehmen will (dann greift die §2-Zeile erneut — eine
Planner-Entscheidung, keine Reviewer-Auflage): `:299-300` (D-1) · `.dockerignore:3-4`
(D-2) · `:267` (D-3) · `:265` (D-4).

**Was nicht geprüft wurde:** `make test` (`-race`, ganzer Baum) ·
`make test-store`/`-replication`/`-integration`/`-notify` · `make coverage-gate`/
`make image`/`make gates` als Ganzes · DoD/LP1–LP3, §6-Ausgänge, Paarungen.
