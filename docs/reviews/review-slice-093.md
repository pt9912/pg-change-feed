# Review-Report: slice-093 — Cluster D2 (Bootstrap-Rest und Telemetrie) · 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen** (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten). DoD-/Spec-Konformität
(Verifier), §6-Risiko-Ausgänge, Register und die drei Paarungen (Planner-Closure) sind
**nicht** Gegenstand.

**Gegenstand:** `014f29c` (Parent `3ed5311`) — **ein** Commit, **8** Dateien, +1023/−20:
**6** `*_test.go` (4 neu, 2 geändert) und **zwei** Dateien ohne Test-Endung
(`.dockerignore`, `harness/sensors/coverage-gate.md`). Kein Produktcode, `THRESHOLD`,
ADRs, keine Naht. Arbeitsbaum vor und nach allen Messungen sauber.

**Skill:** `.harness/skills/reviewer.md` @ `014f29c` · **Datum:** 2026-09-16.

---

## Findings

### F-1 — Die `.dockerignore`-Ausnahme verletzt die `test-only`-Folgepflicht aus `ADR-0082` im Wortlaut — und sie ist die einzige Stelle, an der sie das tut

- `kategorie`: **HIGH**
- `quelle`: `ADR-0082` §Konsequenzen („Folgepflicht (Implementer-Zug, je Test-Slice):
  **test-only**. Der Diff eines Slice dieser Welle enthält ausschließlich `*_test.go` …")
  und ihre Fitness-Function-Zeile „`git diff --name-only` des Slice-Commits | **Test-only**
  … — (Disziplin, **Review-Gegenstand**; kein Sensor)" · `.harness/skills/reviewer.md`
  HIGH „ADR-Verstoß"
- `pfad`: `.dockerignore:5-9` (Kommentar) und `:16` (`!tools/schema/nacharbeit-roles.sql`)
- `befund`: Der Diff ändert eine **Build-Konfigurationsdatei** — die Folgepflicht ist im
  **Wortlaut** verletzt; im **Zweck** („keine Zeile `internal/**`/`cmd/**` um Coverage zu
  gewinnen") **nicht**. Der Implementer nennt die Abweichung selbst offen; die ADR weist
  das Urteil dem Review zu. Selbst gemessen:
  - **Notwendig, nicht beliebig:** ohne die Ausnahme fällt der LP2-Test im Gate real um
    (eigener Lauf: `docker build --target coverage` mit der Parent-`.dockerignore` →
    `--- FAIL: TestRolloutDateiTraegtDieRechteDerVerdrahtung … Rollout-Datei
    /src/tools/schema/nacharbeit-roles.sql nicht lesbar: no such file or directory`,
    Bau-EC **1**).
  - **So eng wie möglich:** eigener Kontextvergleich: Kontext **169 → 170 Dateien**,
    Differenz **genau eine Zeile**. Kein Verzeichnis, kein Muster, keine zweite Ausnahme.
  - **Das Image ist unberührt:** `buildx`-Digest `sha256:84bdca56…` — identisch mit dem in
    `harness/image-hash.txt` geführten Beleg; zusätzlich der inhalts-entscheidende
    Vergleich (`ADR-0044`): extrahiertes `/pg-change-feed` beider Images **sha256
    `bbbf7135…`**.
  - **Präzedenz, kein neuer Mechanismus:** `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`
    (1×, `slice-049`) beschreibt genau diesen Fall für `tools/coverage-gate.sh`.
  - **Ein Weg ohne Build-Konfiguration existiert — er ist teurer:** den Arbeitsbaum in die
    `coverage`-Stufe bind-mounten statt `COPY . .`. Er ändert den **Messmechanismus** der
    Stufe (Bau-Zeit → Lauf-Zeit) und damit ein Gate; die `.dockerignore`-Zeile ist die
    kleinere Änderung. Das ist eine Abwägung, keine Offensichtlichkeit — und sie gehört
    dem Architect, nicht dem Review.
  - **Nebenbeobachtung:** `.dockerignore:3` beschreibt die Ausnahmen als „Stage"-gebunden;
    gemessen sind sie **kontextweit**. Die tragende Aussage („Runtime-Image unberührt")
    ist wahr, die Formulierung „Stage" nicht ganz.
- `verifizierbar`: **nein** — die ADR sagt selbst, dass hierfür kein Sensor existiert
- `klasse`: „ADR-Folgepflicht-Abweichung, notwendig und minimal — undokumentiert"

### F-2 — LP2 bindet drei Grant-Familien scharf und lässt die zwei tragenden Grants derselben Datei ungeprüft: drei reale Regressionen bleiben im Gate grün

- `kategorie`: **MEDIUM**
- `quelle`: DoD `LP2` („an das **reale Artefakt** gebunden") · Kopf-Kommentar des Tests
  (`roles_rollout_file_internal_test.go:29-31`) · `LH-QA-SEC-001`…`003`
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:148-202` gegen
  `tools/schema/nacharbeit-roles.sql:54`, `:73-77`, `:105`
- `befund`: Der Test liest die **echte** Datei (Pfad über `runtime.Caller(0)`, Kommentare
  entfernt, GRANT-/CREATE-ROLE-Regex) und bindet die Zusage **an ihrer Eingabeseite** —
  vier eigene Mutationen färben ihn rot. Seine **Objektliste ist aber eine handverlesene
  Teilmenge**, und fünf Regressionen bleiben **grün**:
  - **`GRANT USAGE ON SCHEMA cdc` entfernt** (`:54`) → **grün**. Das ist der Grant, ohne
    den *alle* drei Rollen kein Objekt im Schema berühren können — die Vorbedingung jedes
    anderen Grants der Datei.
  - **`cdc.retention_blockers` aus dem Reader-Grant entfernt** (`:105`) → **grün**; der
    Test prüft **drei** Views, die Datei grantet **vier**.
  - **der dynamische `GRANT CREATE ON DATABASE` (DO-Block, `:73-77`) entfernt** → **grün**;
    die Regex trifft nur literale `^GRANT`-Zeilen.
  - zwei Least-Privilege-Regressionen ebenfalls grün (`GRANT ALL ON ALL TABLES … TO cdc_reader`;
    ein zusätzlicher Grant auf `cdc.metrics`).
- `verifizierbar`: **ja** — fünf eigene Mutationsproben in der Rollout-Datei, `go test
  -count=1 -run TestRolloutDatei ./internal/bootstrap/`, **EC=0** in allen fünf Fällen
- `klasse`: „gebundene Zusage auf handverlesener Objektliste" — **verwandt** mit
  `BEO-PGC/rollen-test-abdeckungsluecken` (dessen Punkt 1 schließt dieser Diff im
  Mechanismus) und mit `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`

### F-3 — Der Mutations-Beleg des Telemetrie-Tests trägt seinen Satz nicht: der Test färbt rot, aber nicht aus dem genannten Grund

- `kategorie`: LOW
- `quelle`: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (3×, Schwelle erreicht) ·
  `AGENTS.md` §3.12 Instanz B
- `pfad`: `internal/adapters/driven/telemetry/slog_levels_internal_test.go:25-28`
- `befund`: Der Kommentar sagt, die Mutation lasse die Warn-Zeile `level` `INFO` tragen und
  der Test breche daran. Eigene Probe: der Test bricht **rot** — aber die Zeile wird **gar
  nicht geschrieben**: der Handler steht auf `slog.LevelWarn`, ein INFO-Record wird
  gefiltert, und der Test fällt an `json.Unmarshal` („Log-Zeile ist kein gültiges JSON:
  unexpected end of JSON input"); die Zeile `decoded["level"] != c.wantLevel` läuft nie.
  Die **Wirkung** ist wahr, die **Begründung** nicht.
- `verifizierbar`: **ja** — mit der genannten Rumpf-Mutation, EC=1
- `klasse`: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` — vierter Vorgang, wenn die
  Closure ihn dateit

### F-4 — „vorher war sie es nirgends" trägt nicht: der Integrations-Tier liest den Heartbeat-Grant real

- `kategorie`: INFO
- `quelle`: `BEO-PGC/rollen-test-abdeckungsluecken` (Punkt 1) · §3.12 Instanz B
- `pfad`: Commit `014f29c` gegen `tools/harness/run-integration-tests.sh:344-365`
- `befund`: Für **diese** Regression existiert ein realer Leser außerhalb des Gate-Bündels:
  der Integrations-Lauf rollt `nacharbeit-roles.sql` aus, lässt den Feed-Container den
  Heartbeat über `CDC_ADMIN_DSN` schreiben und verlangt den Compose-Status `healthy`; ohne
  den Schreibpfad altert der Heartbeat und der Lauf endet mit Exit 1. Korrekt ist die
  **engere** Formulierung des Testkopfes („bliebe **dort** unsichtbar"). Der
  Integrations-Lauf wurde **nicht** gefahren — Lese-Kette am Skript, keine Messung.
- `klasse`: `BEO-PGC/rollen-test-abdeckungsluecken` — **Präzisierung**, kein neues Vorkommen

### F-5 — `1523/1903 = 80,00 %` ist die gedruckte Rundung, nicht die Rechnung; und die zitierte Gate-Zeile nennt ihre Schwelle nicht

- `kategorie`: INFO
- `quelle`: §3.12 Instanz A
- `pfad`: Commit `014f29c` gegen `harness/mk/coverage.mk` (`THRESHOLD ?= 70`),
  `tools/coverage-gate.sh:52`
- `befund`: Eigene Messung der Zählbasis: **1903** Nenner, **1523** gedeckt → **80,0315 %**,
  gedruckt `80.0%`. „80,00 %" ist die `%.2f`-Ausgabe des Gate-Skripts über die
  **ein-stellig gedruckte** Zahl — zwei Größen für eine Zelle. Zweitens nennt die zitierte
  Zeile die **Schwelle** nicht (`… erfüllt Schwelle 70%`; `THRESHOLD` steht unverändert auf
  **70**). Für §7 gehören **beide** Wellen-Belege hinein — eigener Lauf: `make coverage-gate
  THRESHOLD=80` → `OK — Coverage 80.00% erfüllt Schwelle 80%`, EC 0; `THRESHOLD=85` →
  `FAIL — Coverage 80.00% unter Schwelle 85%`, Bau-EC 1.
- `verifizierbar`: **ja**
- `klasse`: — (Träger ist die git-Historie; keine Zählung)

### F-6 — Der Register-Eintrag `test-integration-retention-timing-flake` (3×) nennt genau den Block, den dieser Diff festfährt — gegenläufig, kein Vorkommen

- `kategorie`: INFO
- `quelle`: `BEO-PGC/test-integration-retention-timing-flake` · `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
- `pfad`: `…/test-integration-retention-timing-flake/evidence/slice-091.md:10`, `:23`, `:41` — **außerhalb** des Diffs
- `befund`: Der Eintrag nennt als ersten von „zwei unruhigen Blöcken" den „Takt-Zweig von
  `runWALRetentionCheck` (`wiring.go:991.5,992.13`)" — genau den Block, den der neue Test
  deterministisch fährt (vier eigene Läufe: `count > 0`). Die Evidence-Datei ist
  **unveränderlich ab Merge** und wird zu Recht nicht angefasst; die Beobachtung selbst ist
  nicht widerlegt. Zu melden: der Eintrag steht auf der Schwelle, und **einer seiner zwei
  benannten Träger fällt weg** — Material für den Lese-Schritt.
- `verifizierbar`: **ja**
- `klasse`: kein Vorkommen (gegenteilig)

---

## Negativbefunde

- **die Zahlen des Commits — selbst nachgemessen, blockweise.** Profil `3ed5311`: **410**
  ungedeckt / **1903**, `78.5%`; Profil `014f29c`: **380** / **1903**, `80.0%`.
  Blockdifferenz: **+30 Statements in 27 Blöcken, −0**. `internal/bootstrap` **598/290/308**
  (Parent **598/262/336** — exakt `336 → 308`), `telemetry` **6/6/0** (Parent `6/4/2`).
  Die `+28`/`+2` und „Nenner unverändert 1903" halten **punktgenau**; kein anderes Paket
  bewegt sich. Vier eigene Läufe: 4 × **1523**, 4 × `80.0%`.
- **§6-Risiko 3 (Coverage-Theater) nicht eingetreten** — 20 eigene Mutationsproben. Die
  fünf **grünen** Proben sind kein Theater-Fund, sondern **F-2**.
- **die 14 im Diff genannten Mutationen sind alle real rot** — jede selbst gefahren. Nur
  eine **Begründung** trägt nicht (F-3).
- **die Kern-Bindung ist echt.** Der Test liest `tools/schema/nacharbeit-roles.sql` über
  `filepath.Join(filepath.Dir(runtime.Caller(0)), …)` — unabhängig vom Arbeitsverzeichnis,
  im Gate real aufgelöst; er entzieht und erteilt **nichts** (`grep -n "REVOKE\|GRANT \|Exec("`
  → keine SQL-Ausführung).
- **der Gate-Rot-Beleg endet nicht am Test.** Mit „SELECT aus dem Heartbeat-Grant entfernt"
  fällt der **Bau der `coverage`-Stufe** real um (Bau-EC 1). Die Kern-Aussage von LP2 hält
  **end-to-end**.
- **die Sensor-Doku-Korrektur trägt — und macht nichts Neues falsch.** Der Satz „in zwei
  Blöcken gemessen" ist durch diese Arbeit tatsächlich falsch geworden (vier Läufe, Block
  durchweg `count > 0`); die Korrektur setzt den Herkunfts-Anker (`slice-093`), markiert die
  neue Zahl als **abgeleitet** mit ihrer Rechnung, schließt die **Deixis**-Falle („die
  Messung **jenes** Test-Stands, nicht die des geltenden"), ersetzt „hier `71.9%`" durch
  `z. B. 71.9%`, und beide `§`-Verweise lösen auf.
- **das verbleibende Band kann das Gate bei 80 nicht kippen** — auch der ungünstigste Fall
  (1522 = 79,979 %) druckt `80.0%`.
- **Nebenläufigkeit und Laufzeit.** `make test` (`-race`, `--network none`) → **EC=0**,
  32 × `ok`, 0 × `FAIL`, 0 × `DATA RACE`; die neuen Tests **10 ×** unter `-race` grün.
- **`make gates` grün, Exit separat gelesen** → **0**: sechs Checks, `d-check` **764/0`,
  danach Baum leer.
- **§3.7, §3.11, §3.2** — kein Produktionscode-Kommentar geändert, kein Slice-Vokabular in
  Produktionskommentaren, kein host-lokaler Pfad, kein `//nolint`.
- **Traceability und Betreiber-Oberfläche** — Betreff mit [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md), ohne Struktur-ID;
  keine neue Kennung, keine neue Oberfläche.
- **[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) „test-only" im Übrigen** — 6 der 8 Pfade `_test.go`; die Sensor-Doku ist
  durch §3 des Plans ausdrücklich vorgesehen; kein `git mv`, keine Naht, kein `THRESHOLD`,
  keine Accepted-ADR berührt. Die vierte Stelle ist **F-1**.
- **Planungs-Dokumente außerhalb des Diffs** — `welle-20.md` §4 trägt **Soll**-Zahlen und
  den ADR-Stand; kein Ist-Stand in einem stehenden Träger.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-only 014f29c^..014f29c` | 0 | **8** Pfade, **2** ohne `_test.go` |
| Profil `3ed5311` / `014f29c` | 0 | 410/1903 `78.5%` · 380/1903 `80.0%` |
| Block-Differenz Parent → Diff | 0 | **+30** in **27** Blöcken, **−0** |
| 4 × `make coverage-gate THRESHOLD=80` | 0 | 4 × **1523**, 4 × `80.0%`, `erfüllt Schwelle 80%` |
| `make coverage-gate THRESHOLD=85` | **1** | `FAIL — Coverage 80.00% unter Schwelle 85%` |
| `make gates` (Default 70) | 0 | 6 Checks grün, `d-check` 764/0, Baum leer |
| Bau mit `.dockerignore`-Ausnahme | 0 | `buildx`-Digest `sha256:84bdca56…` = geführter Beleg; Binär-sha256 `bbbf7135…` |
| Bau **ohne** Ausnahme (Parent-`.dockerignore`) | **1** | LP2-Test fällt im Gate |
| Kontextvergleich | 0 | **169 → 170** Dateien, Differenz = genau diese Datei |
| `make test` (`-race`, `--network none`) | 0 | 32 × `ok`, 0 × FAIL, 0 × DATA RACE |
| neue Tests `-race -count=10` | 0 | 2 × `ok`, kein Flake |
| 4 × Bootstrap-Profil, Block `991.5,992.13` | 0 | 4 × `count > 0` |
| **20 Mutationsproben** | 1 je roter | **14 × rot**, **5 × grün** (F-2), 1 verworfen |

## Antwort auf die Schwerpunkte

**(1) Die `.dockerignore`-Zeile — Stellungnahme.** (a) **So eng wie möglich**: +1 Datei, keine
Verzeichnis-, keine Muster-Ausnahme. (b) Der Weg über den Build-Kontext ist **richtig
gewählt**; ohne die Ausnahme fällt der LP2-Test im Gate real um (gemessen). Ein Weg **ohne**
Build-Konfiguration existiert (Bind-Mount, Muster `make proto-generate`), ist aber die
**größere** Änderung, weil er den Messmechanismus der Stufe verschiebt — eine Abwägung für
den Architect. (c) Das Image **ändert sich nicht** (Digest **und** Binär-sha256 identisch,
geprüft). **Urteil:** die Abweichung trägt in der Sache und ist die einzige Stelle, die
`ADR-0082`s `test-only`-Wortlaut bricht — sie ist damit **undokumentiert lockernd** und
gehört als Artefakt nachgezogen (Modul 8 §Konflikt-Pfad, drittes Verdikt: *Lockerung
legitim, aber undokumentiert*). **Kein** Implementer-Pfeil.

**(2) LP2 — der wertvollste Teil des Diffs, und er trägt.** Der Test liest das reale
Artefakt, färbt bei vier eigenen Mutationen rot, und der rote Fall endet am **Gate**. Die
Objektliste ist handverlesen → **F-2**.

**(3) Die Sensor-Doku-Korrektur trägt und macht nichts Neues falsch** — Anker gesetzt,
Ableitung markiert, Deixis-Falle geschlossen, `§`-Verweise lösen auf.

**(4) Prüf-Kraft statt Zeilen.** 20 Proben: 14 rot, 5 grün (F-2), 1 verworfen. Risiko 3
nicht eingetreten; Risiko 2 nicht eingetreten; **Risiko 4** (`arbeit-ueberholt`) real und
richtig bedient.

**(5) §3.12 — welche Zahl wohin.** **gemessen, mit Lauf:** `380`/`1903`/`80.0%` (Lauf
`slice-093`), `336`/`308`/`2`/`0`, `1523` in vier Läufen. **abgeleitet, nach §7 mit
Rechnung:** `+30`, `1 Statement = 0,05 pp`, `80,03 %` exakt gegen `80,00 %` gedruckt.
**Mit Lauf und Schwelle:** der Grün-Beleg bei **80** und der Rot-Beleg bei **85** (beide
gefahren).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **1** |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

## Verdikt

**Der Diff trägt — in der Sache und in der Zahl.** Test-only im Zweck, die Zahlen halten
punktgenau, der Kern (LP2) ist an das reale Artefakt gebunden und färbt das **Gate** real
rot, `make gates` und `make test` grün, Prüf-Kraft mit 20 Proben belegt.

**Merge-blockierend:** **ja — aber nicht für den Implementer-Code.** F-1 ist ein HIGH nach
der geschlossenen Liste; die Abweichung ist notwendig, minimal und image-neutral, also ist
sein Artefakt ein **Architect-/Planner-Nachzug**, **kein** Code-Nachzug. F-2 braucht einen
**kleinen Implementer-Nachzug** (zwei Assertions + eine benannte Grenze für den dynamischen
Grant). F-3 ist ein Kommentarsatz und geht mit. **Fixrunde:** ja (F-2), und die deckt nach
Slice-Plan §2 ein **Delta-Review** ab. **Rückgabe-Pfeil:** F-2 und F-3 (beide klein).
F-1 gehört dem Architect; F-4/F-5/F-6 der Closure bzw. dem `welle-20`-Lese-Schritt.

**DoD-Häkchen „Review durchgeführt":** bleibt **offen** — es kommt eine Fixrunde.

**Nicht gefahren:** `make test-store`, `make test-replication`, `make test-integration`
(F-4 stützt sich auf eine Lese-Kette, nicht auf einen Lauf) und `make image` — die
Image-Frage wurde über `buildx`-Digest **und** Binär-sha256 entschieden.
