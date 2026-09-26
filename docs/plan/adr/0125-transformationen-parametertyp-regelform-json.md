# ADR-0125: Transformationen — die Regelform kommt als `json`-Parameter in `cdc.set_transformation` (Supersedes ADR-0112, teilweise)

**Status:** Accepted — Supersedes [`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
in genau zwei Stellen: den Parametertyp `rule_spec jsonb` in der Festlegung von
Teilfrage 1 (die Signatur von `cdc.set_transformation`) und den Satz „ein
JSON-Parameter je Regel ist in SQL nur als `jsonb` durchreichbar“ in der
Contra-Zelle der Option D. Alles Übrige von `ADR-0112` bleibt in Kraft, insbesondere
die **Spalte** `rule_spec jsonb`, der Regelsatz, K1–K4 und die Validierung
vollständig in Go.

**Datum:** 2026-09-26

**Autor:** Architect-Agent (Modul 8), Architect-Zug zu einer Abweichung im
Slice `slice-transformationen-antragsweg-schema`; jede Tatsachenaussage trägt
ihren Beleg-Anker oder ist als hergeleitet gekennzeichnet (siehe §Kontext).

**Bezug:** [`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Konfiguration der
Transformation; Haupt-Bezug),
[`LH-FA-ADM-001`](../../../spec/lastenheft.md) (Administration über SQL),
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(teilweise superseded — Haupt-Bezug),
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (Schemamigrationen,
Nacharbeit-Form für Funktionen),
[`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md) (keine
Domänenlogik in SQL),
[`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue)

**Schärft:** [`SPEC-019`](../../../spec/pflichtenheft.md) — der Antrags-Datensatz
und der Absatz „Transformations-Antragsarten“: die Regelform kommt als
`json`-Parameter, die Spalte bleibt `jsonb`.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md) ist
`Accepted` und unberührbar (`AGENTS.md` §3.5). Seine Festlegung von Teilfrage 1
nennt `cdc.set_transformation(source_id, schema_name, table_name, rule_name,
rule_spec jsonb)`; `SPEC-019` übernimmt „reichen `rule_spec` als `jsonb` durch“. Mit
einem `jsonb`-Parameter ist der zweite `make schema-rollout` nicht ausführbar. Die
Berichtigung ändert den Referenten (den Parametertyp) und ist keine Zitat-Korrektur.

### Gemessen

**Ein Scratch-Baum aus `git archive HEAD` (Stand `f67e2916`), eine frische
PostgreSQL-18-Instanz (`PG_TEST_IMAGE` und `D_MIGRATE_IMAGE` aus dem `Makefile`),
zwei aufeinanderfolgende `make schema-rollout`, Parametertyp je Lauf geändert;
der Arbeitsbaum des Repositories ist unberührt (die Läufe liegen nicht im
Repository).** Gedruckt:

| Parametertyp von `p_rule_spec` | Lauf 1 | Lauf 2 |
|---|---|---|
| `json` | `make`-Exit 0 | `make`-Exit 0 |
| `jsonb` | `make`-Exit 0 | `make`-Exit 2 (`make: *** [Makefile:278: schema-rollout] Fehler 5`); im Report `executionError":"ERROR: function set_transformation(text, text, text, text, json) does not exist","recoverability":"FULL_ROLLBACK_CONFIRMED"` |
| `text` (Guard-Eintrag auf `in:text` gestellt) | `make`-Exit 0 | `make`-Exit 0 |

Im `jsonb`-Lauf nimmt die Wache den Guard-Eintrag `in:json` an (Meldung
„bekannte Fremdobjekt-Blocker … `--allow-destructive`“): der Fehler liegt in
`--execute`, nicht bei der Wache. d-migrate 1.3.1 rendert den Abbau einer Funktion mit
`json`- oder `jsonb`-Parameter als `DROP FUNCTION … (…, json)`; für die
`jsonb`-Funktion existiert diese Signatur nicht (die Aussage zur
Report-Schreibweise `in:json` für beide Typen ist im Plan des Slice gemessen und
hier durch die angenommene Wache bestätigt; die Rendering-Regel selbst ist aus dem
Report-Fehler hergeleitet).

**Aufrufformen je Parametertyp** (dieselben Scratch-Instanzen nach den zwei
Rollouts, `SELECT cdc.set_transformation('s','public','t','r', <Form>)` unter dem
Superuser; die sieben Formen der Tabelle sind die gesamte Messung, eine
Aussage über weitere Formen folgt nicht daraus):

| Form | `json` | `jsonb` | `text` |
|---|---|---|---|
| Literal `'{"kind":"rename_column"}'` | angenommen | angenommen | angenommen |
| `…::json` | angenommen | „function … does not exist“ | „function … does not exist“ |
| `…::jsonb` | „function … does not exist“ | angenommen | „function … does not exist“ |
| `jsonb_build_object('kind','x')` | „function … does not exist“ | angenommen | „function … does not exist“ |
| `'{oops'` (ungültiges JSON) | `invalid input syntax for type json`, keine Zeile | dasselbe | dasselbe |
| `'[1]'` (JSON ohne Objekt) | angenommen, Zeile `pending` | dasselbe | dasselbe |
| `NULL` | angenommen, Zeile `pending` | dasselbe | dasselbe |

Die Spalte `rule_spec` ist in allen Varianten `jsonb`: das neutrale Modell trägt
`rule_spec: { type: json }` in `tools/schema/schema.yaml`, und der Report gibt
`"rule_spec" JSONB` aus (`tools/schema/plan.yaml`, `CREATE TABLE
"administration_request"`).

### Konstraints

- Der Rollout ist idempotent ([`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md),
  [`LH-QA-OPS-005`](../../../spec/lastenheft.md)): ein zweiter Lauf gegen ein
  migriertes Ziel endet mit Exit 0.
- Keine Domänenlogik in SQL ([`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)):
  die Funktion schreibt den Antrag; Regelform, Regeltyp und Spalte prüft Go
  (`SPEC-019`, `SPEC-030`).

## Entscheidung

Wir wählen **den Parametertyp `json` für die Regelform in `cdc.set_transformation`**;
die Funktion schreibt den Wert als `jsonb` in die Spalte `rule_spec`. Zwei
Festlegungen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `jsonb`-Parameter wie in `ADR-0112` | Wortlaut von `ADR-0112`/`SPEC-019` gilt; `::jsonb` und `jsonb_build_object(…)` werden angenommen | der zweite `make schema-rollout` endet mit Exit 5 (gemessen); die Idempotenz von `ADR-0043` fällt |
| B — die Nacharbeit-Funktionen aus dem Abbau-Plan von d-migrate nehmen (Tooling-Weg, `--routine-capability 'function:enabled=false'` o. ä.) | der Parametertyp bliebe `jsonb` | am Slice ungeklärt: `--plan-only` endete mit dem Schalter weiter mit Exit 8 (übernommen aus dem Plan des Slice, nicht wiederholt); ein Eingriff in `make schema-rollout`, Wache und Guard, der auch die sechs übrigen Funktionen berührt, für einen einzelnen Parametertyp |
| C — `text`-Parameter, Go parst | Rollout gemessen grün; Guard-Schreibweise wie bei den übrigen Funktionen | kein Gewinn gegenüber `json`: ungültiges JSON scheitert gemessen ebenso schon beim Aufruf (Cast in der Funktion), und ein `::json`-Wert wird zusätzlich abgelehnt; die Validierung nach Go zu verlagern verlangte eine `text`-Spalte, also eine Änderung von `SPEC-019`, Schema und Store ohne fachlichen Anlass |
| **D — `json`-Parameter, Spalte `jsonb` (gewählt)** | Rollout gemessen grün (Exit 0 und 0); Spalte, `SPEC-030` und Store unverändert; die Zusage „Antrag wird angenommen, Prüfung in Go“ gilt wie mit `jsonb`: `NULL`, JSON-`null` und JSON ohne Objekt werden angenommen und in Go abgelehnt | ein `jsonb`-typisierter Wert wird abgelehnt; Fehlermeldung „function … does not exist“ statt eines Typhinweises |

### Festlegung 1 — Parametertyp und Spalte

`cdc.set_transformation(source_id text, schema_name text, table_name text,
rule_name text, rule_spec json)`; die Funktion schreibt `rule_spec::jsonb` in die
Spalte `rule_spec` (Typ `jsonb`, nullable). Ein Literal und ein `::json`-Wert
werden angenommen, ein `jsonb`-typisierter Wert (`::jsonb`, das Ergebnis von
`jsonb_build_object`) wird mit „function … does not exist“ abgelehnt und braucht
den Cast `::json`. Syntaktisch ungültiges JSON scheitert beim Aufruf, ohne
Antrags-Zeile; alles, was gültiges JSON ist — auch `NULL`, JSON-`null` und ein
Wert ohne Objekt —, wird als Antrag angenommen und in Go abgelehnt (`SPEC-019`,
Fehlertext-Tabelle). Diese Festlegung ersetzt den Parametertyp in Teilfrage 1
von `ADR-0112`.

### Festlegung 2 — Die Begründung „nur als `jsonb` durchreichbar“ entfällt

Ein JSON-Parameter ist in SQL als `json`, `jsonb` oder `text` durchreichbar
(die Tabelle oben); die Wahl des Typs hängt an der Rollout-Mechanik der
Nacharbeit-Funktionen ([`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md)), nicht am
Datentyp der Spalte. Diese Festlegung ersetzt den Satz in der Contra-Zelle der
Option D von `ADR-0112`.

## Konsequenzen

- Positiv: der zweite Rollout endet mit Exit 0; Spalte, Store-Lesen
  (`COALESCE(rule_spec::text, '')`), `SPEC-030` und der Guard-Eintrag
  `in:json` bleiben, wie sie sind.
- Negativ (Aufruferkomfort): wer die Regelform als `jsonb`-Wert übergibt, schreibt
  `…::json`; der Fehlertext nennt den Typ nicht als Ursache. Das Betreiber-Handbuch
  trägt die Aufrufform (Folgepflicht 2).
- Negativ (Grenze, benannt): die Messung ist an PostgreSQL 18 mit dem gepinnten
  d-migrate-Image (`D_MIGRATE_IMAGE`) gemessen; für PostgreSQL 17 liegt der zweite
  Rollout nicht vor (hergeleitet: das Rendering des Abbaus liegt in d-migrate, nicht
  in der PostgreSQL-Version).
- Folgepflicht 1: **`SPEC-019` Nachzug** — der Absatz „Transformations-Antragsarten“
  nennt den Parametertyp und die Aufrufform (in demselben Commit wie diese ADR).
- Folgepflicht 2: **Betreiber-Handbuch** — `slice-transformationen-betriebsdoku`
  nennt in „Transformationsregel konfigurieren“ die Aufrufform der Regelform
  (Literal oder `::json`, ein `jsonb`-Wert braucht den Cast).
- Folgepflicht 3: **Slice `slice-transformationen-antragsweg-schema`** — Plan §3
  (Zeile `nacharbeit-administration.sql`) und §6 (Punkt „Der Parameter `rule_spec`
  ist `json`“) verweisen auf diese ADR statt auf eine offene Architect-Frage; das
  Risiko trägt den Ausgang „eingetreten, entschieden durch `ADR-0125`“. Die
  Folgepflicht 3 von `ADR-0112` (Antragsweg und Dauerhaftigkeit) bleibt in Umfang
  und Beleg unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Rollout-Skript gegen eine Wegwerf-PostgreSQL (`tools/harness/run-schema-rollout-guard-test.sh`, Lauf 2 und Lauf 5) | der zweite `make schema-rollout` und der Rollout des Arbeitsbaums über den Stand des jüngsten `v*`-Tags enden mit Exit 0; Mutation (Parametertyp `jsonb`) → zweiter Rollout `Fehler 5` (im Scratch-Baum gemessen, mit dem Skript selbst **nicht** wiederholt; dass Lauf 2 rot färbt, ist hergeleitet aus `set -e` bei der Zuweisung des Laufs) | — (kein `make`-Ziel, `bash tools/harness/run-schema-rollout-guard-test.sh`) |
| Go-Test, Paket `rolloutguard` | der Guard trägt die Schreibweise `set_transformation(in:text,in:text,in:text,in:text,in:json)` und lehnt `in:jsonb` ab (Mutation `in:json` → `in:jsonb` färbt drei Tests rot, im Plan des Slice gemessen) | `make test` |
| Go-Test, Paket `bootstrap` | `TestAdministrationDateiTraegtDieFunktionsRechte` liest die Signatur `cdc.set_transformation(text, text, text, text, json)` aus `nacharbeit-administration.sql` und bindet `GRANT`/`REVOKE` daran | `make test` |
| Go-Test, reale PostgreSQL | Store-Round-Trip: `TestAdministrationRequestTransformationRequestsCarryRuleAndNotify` ruft `cdc.set_transformation(…, $5::json)` und liest die Zeile mit `rule_spec` zurück | `make test-store` |

## Re-Evaluierungs-Trigger

- **Eine d-migrate-Version rendert den Abbau einer Funktion mit `jsonb`-Parameter
  als `jsonb`** (`make pin-stale-dmigrate` meldet Pin-Drift; die Messung der
  Tabelle wiederholen): eine Folge-ADR kann den Parametertyp auf `jsonb` stellen.
- **Die Nacharbeit-Funktionen wandern in das neutrale Modell** (Trigger von
  [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md)): den Parametertyp neu
  bestimmen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-26 | Accepted — Architect-Zug (Vollmacht des Auftraggebers); die Rollout-Läufe der drei Parametertypen und die Aufrufformen im Scratch-Baum nachgemessen, jede Tatsachenaussage an einer gedruckten Zeile oder als hergeleitet gekennzeichnet | [`LH-FA-CFG-007`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0125` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
