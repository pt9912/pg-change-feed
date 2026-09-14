# Review-Report: slice-062 — 2026-09-14

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-062` — ein Implementierungs-Commit `ed7418c`, Elter
`7292870` (`ADR-0063`, Architect-Zug), davor `96dbca8` (reiner
`next→in-progress`-Move). `git show --stat ed7418c` bestätigt genau vier
geänderte Dateien: `docs/plan/planning/in-progress/slice-062-…md`,
`harness/README.md`, `test/integration/integration_test.go`,
`tools/harness/run-integration-tests.sh` — kein Griff in
`internal/adapters/driving/replication/mapper/mapper.go` (`ADR-0063`
erklärt diese Datei ausdrücklich als unverändert bestehenden Code, nicht
Gegenstand der ADR).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft
2026-09-13: vier repo-spezifische HIGH-Regeln plus Slice-/Wellen-Chronik-
und Handbuch-Versionshistorie-Regel)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-062-e2e-schema-drop-column-metadaten-erweiterbarkeit.md`
  (vollständig gelesen, inkl. §6/§8)
- `docs/plan/adr/0058-testansatz-fuenf-luecken.md` (vollständig gelesen —
  bindend für `LH-FA-DAT-006`/Entscheidung 2, für `LH-FA-SCH-003` nur noch
  über `ADR-0063` bindend)
- `docs/plan/adr/0063-lh-fa-sch-003-testform-korrektur.md` (vollständig
  gelesen — bindend für `LH-FA-SCH-003`s Testform, inkl. dem exakt
  vorgegebenen Code-Block für `TestE2ESchemaChangeDropColumn`)
- `internal/adapters/driving/replication/receive/receive.go`
  (`ensureSlot`) — geprüft, ob ein neu angelegter Replication-Slot beim
  Container-Neustart tatsächlich ohne Replay startet
- `tools/schema/schema.yaml` (`table_schema`/`schema_version`-Abschnitte)
  — Constraint-Prüfung für die direkten SQL-Eingriffe im Recovery-Block
- `compose.yaml` (`CDC_TABLES`, `SLOT`) — Abgleich der im Recovery-Block
  verwendeten `SCHEMA_TABLE_ID=tbl-e2e-schema` gegen den realen
  Container-Vertrag
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.5, 3.7, 3.9), §5
  Traceability-Regeln
- `harness/conventions.md` (MR-000 ID-Schema)

---

## Der zentrale Prüfpunkt: Recovery-Mechanismus zwischen `TestE2ESchemaChangeDropColumn` und `TestE2ESchemaChangeIncompatibleTypeChange`

**Behauptung des Implementers:** Ein bloßer `docker start` reicht zwischen
den beiden Container-beendenden Testfunktionen nicht — zwei unabhängige
Poison-Zustände (Replication-Slot-Replay der bereits verarbeiteten
`ADD COLUMN removable`-Transaktion; veraltete, `removable` weiterhin
listende Schema-Version) wurden real gegen den Compose-Stack gefunden und
über Slot-Neuanlage + Schema-Version-Nachtrag behoben.

**Eigene Prüfung, nicht die Selbstauskunft übernommen:**

1. **`git show ed7418c` vollständig gelesen**, insbesondere der neue
   Recovery-Block in `tools/harness/run-integration-tests.sh` (Zeilen
   ~1580–1650 des neuen Stands).
2. **Slot-Scope real geprüft.** `SLOT=slot_pgc_e2e` ist repo-weit ein
   einziger, globaler Slot für die eine Quelle `src-e2e` dieses
   E2E-Compose-Stacks (`grep -n SLOT` liefert nur eine Definitionsstelle,
   `compose.yaml`s `CDC_SLOT`-Vertrag trägt denselben Wert) — der Drop
   betrifft also die gesamte Erfassung dieses Stacks, nicht nur
   `feed_e2e_schema`. Das ist unkritisch, weil `run-integration-tests.sh`
   streng sequenziell läuft: Jede vorangegangene Testphase (Retention-
   Lebenszyklus, NATS-Belege, HTTP-Rundlauf, Black-Box-CLI-Rundlauf) hat
   ihre Daten bereits real über `cdc.changes` gelesen und bestätigt,
   **bevor** dieser Slice-Code läuft — bereits verarbeitete Changes stehen
   dauerhaft in `cdc.change` (unabhängig vom Slot-Zustand); ein Slot-Drop
   entfernt keine bereits geschriebenen Zeilen, sondern nur die
   Redelivery-Position für künftige Reconnects. Es gibt zum Zeitpunkt des
   Drops keinen unbestätigten Rückstand aus früheren, bereits
   abgeschlossenen Testphasen.
3. **`ensureSlot` real gelesen** (`internal/adapters/driving/replication/receive/receive.go:242`):
   Ein neu angelegter Slot startet am `ConsistentPoint` (aktueller
   WAL-Stand zum Anlagezeitpunkt), **kein** Replay davor liegender WAL.
   Das bestätigt technisch, warum Slot-Neuanlage statt bloßem
   `docker start` die beschriebene doppelte Wiedereinspielung
   (`ADD COLUMN removable`-Transaktion, dann erneut die `DROP
   COLUMN`-Transaktion) tatsächlich verhindert — kein Zufallstreffer,
   sondern eine dokumentierte PostgreSQL-Eigenschaft.
4. **Schema-Version-Nachtrag real gegen die Constraints geprüft**
   (`tools/schema/schema.yaml` §`table_schema`): Primärschlüssel ist
   `(schema_version_id, ordinal_position)`, kein Zwang zu lückenloser
   Ordinalfolge (nur `CHECK (ordinal_position >= 1)`). Die neue Version
   `tbl-e2e-schema-v4` kopiert `id`/`name`/`amount`/`extra` unverändert
   aus `tbl-e2e-schema-v3` und lässt nur `removable` aus — exakt die reale
   Post-Drop-Spaltenmenge, ohne Primärschlüssel- oder
   Fremdschlüssel-Konflikt; `tbl-e2e-schema-v3` (weiterhin referenziert
   von `cdc.change.schema_version` für die historische id=3-Zeile,
   Boundary-Assertion von `TestE2ESchemaChangeDropColumn`) bleibt
   unangetastet. `SCHEMA_TABLE_ID=tbl-e2e-schema` stimmt exakt mit dem in
   `compose.yaml`s `CDC_TABLES`-Vertrag konfigurierten Wert überein — kein
   geratener Wert.
5. **Alternative geprüft, ob die vorhandene SQL-Administration
   (`cdc.enable_table`/`disable_table`, `ADR-0050`) statt direkter
   SQL-Eingriffe genutzt werden könnte:** Nein — dieser Mechanismus läuft
   asynchron über die Administrations-Goroutine **im laufenden
   Feed-Container**; der Container ist an dieser Stelle des Skripts aber
   gerade beendet. Der direkte SQL-Eingriff ist damit nicht nur eine von
   mehreren Optionen, sondern der einzig verfügbare Weg — die
   Delegation in `ADR-0063` §Konsequenzen an „Implementierungsdetail des
   umsetzenden Slices" ist damit sachlich getragen.
6. **Real und wiederholt reproduziert (eigener, unabhängiger Nachvollzug,
   nicht nur der Implementer-Log-Ausschnitt):** Zwei vollständige,
   unabhängige `make test-integration`-Läufe (siehe unten) — beide Exit 0,
   beide mit identischem Ablauf: `TestE2ESchemaChangeDropColumn` PASS →
   Recovery-Block registriert real `tbl-e2e-schema-v4` → Health-Poll
   erfolgreich → `TestE2ESchemaChangeIncompatibleTypeChange` PASS ohne den
   vom Implementer beschriebenen artefaktbedingten `schema`-Fehler. Kein
   Flake in beiden Läufen.

**Verdikt zu diesem Punkt: korrekt, nicht nur zufällig einmal
funktionierend.** Der Mechanismus stützt sich auf real geprüftes
System-Verhalten (`ensureSlot`-Semantik, PK/FK-Constraints des
Schema-Version-Modells, Abwesenheit eines lauffähigen
Administrations-Alternativwegs bei gestopptem Container) und ist über zwei
eigene volle Läufe reproduziert.

## Weitere Findings

Keine HIGH/MEDIUM/LOW-Findings. Zwei INFO unten.

### F-1 — Direkter SQL-Eingriff in `cdc.schema_version`/`cdc.table_schema` ist ein Sonderfall ohne Produktions-Äquivalent

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh` (Recovery-Block nach
  `TestE2ESchemaChangeDropColumn`)
- `befund`: Der Recovery-Mechanismus manipuliert `cdc.schema_version`/
  `cdc.table_schema` direkt per `psql`, weil es für den Fall „Feed-Container
  nach `relationOther`-Fehler dauerhaft beendet, Schema-Version veraltet"
  keinen laufenden Administrationsweg gibt (`cdc.enable_table` braucht den
  laufenden Container). Das ist für diesen Testharness korrekt gelöst,
  benennt aber implizit eine Betriebslücke: Ein echter Produktionsbetrieb
  hätte nach einem vergleichbaren `relationOther`-Abbruch denselben
  Nachtrag-Bedarf, aber kein dokumentiertes Werkzeug dafür.
- `verifizierbar`: nein — kein Gate-Lauf bestätigt die Abwesenheit eines
  Produktions-Recovery-Wegs; die Beobachtung selbst ist der Beleg.
- `klasse`: „Schema-Version-Nachtrag nach relationOther-Abbruch ohne
  Produktions-Administrationsweg"

Hinweis an die Planner-Rolle zur Erwägung, ob dies neben dem bereits
verkörperten `BEO-PGC/schema-evolution-nicht-dynamisch` einen eigenen
Beobachtungs-Register-Eintrag verdient — die beiden Beobachtungen sind
verwandt (Schema-Version-Pflege), aber nicht deckungsgleich: die
verkörperte Beobachtung betraf die Abwesenheit *jeder* dynamischen
Re-Versionierung (behoben `seit slice-033`), diese hier betrifft die
Abwesenheit eines *Recovery*-Wegs nach einem bewusst harten Abbruch.

### F-2 — Einmaliger `docker inspect`-Running-Check ohne Retry direkt nach dem crash-auslösenden Testlauf

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:1564-1568` (Check nach
  `TestE2ESchemaChangeDropColumn`)
- `befund`: Der Check `feed_running=$(docker inspect …)` läuft als
  Einzelabfrage ohne Retry-Schleife unmittelbar nach dem `go
  test`-Aufruf — theoretisch ein schmales Zeitfenster zwischen dem
  Heartbeat-Schreiben (worauf `awaitHeartbeatErrorClass` pollt) und dem
  tatsächlichen Prozess-/Container-Exit. Dieses Muster ist aber
  **kein neues Risiko dieses Slice**: derselbe Einzelabfrage-Stil ohne
  Retry steht bereits an über einem Dutzend Stellen im selben Skript
  (`grep -n feed_running` zeigt denselben Aufbau seit früheren Slices).
  Zwei eigene, vollständige `make test-integration`-Läufe zeigten keinen
  Flake an dieser Stelle.
- `verifizierbar`: ja — wiederholte `make test-integration`-Läufe; bislang
  kein Fehlschlag beobachtet.
- `klasse`: „Einzelabfrage-Container-Status-Check ohne Retry"

Kein erwarteter Reviewer-Schritt aus beiden INFO-Punkten — Hinweis an die
Planner-Rolle zur Erwägung im Rahmen der Slice-Closure (§7/
Beobachtungs-Register).

## Negativbefunde

- geprüft, ohne Befund: `TestE2ESchemaChangeDropColumn` gegen `ADR-0063`s
  exakt vorgegebenen Code-Block — Zeile-für-Zeile identisch übernommen
  (Funktionskörper unverändert; nur der Godoc-Kommentar darüber ist neu
  hinzugefügt, was `ADR-0063` nicht ausschließt).
- geprüft, ohne Befund: `TestE2EChangeTableMetadataExtensibility` — folgt
  `ADR-0058` Entscheidung 2 punktgenau (Tabelle `feed_e2e_full`, additive
  `jsonb DEFAULT NULL`-Spalte, `t.Cleanup`-Rückbau, Happy-Path-Reread vor
  der Erweiterung, Boundary-Reread danach, Nicht-Sichtbarkeit über
  `cdc.changes`); keine ID-Kollision mit anderen `feed_e2e_full`-Testfällen
  (IDs 500/501 sind im gesamten Testpaket sonst ungenutzt, andere
  Testfälle auf derselben Tabelle bleiben im Bereich ≤ 241). Keine
  Kollateral-Änderung durch die Recovery-Arbeit: Die Funktion steht
  unverändert in der vorderen Testgruppe, ihr `-run`-Muster-Eintrag wandert
  nicht.
- geprüft, ohne Befund: `-run`-Musterzeilen — `TestE2EChangeTableMetadataExtensibility`
  korrekt in der vorderen (Container-lebt-noch)-Gruppe ergänzt,
  `TestE2ESchemaChangeDropColumn` korrekt **nicht** dort, sondern als
  eigener Aufruf nach der bisherigen Container-Ende-Grenze — Quelltext-
  Deklarationsreihenfolge (nach `TestE2ESchemaChangeAddColumn`, vor
  `TestE2EHeartbeatHealthy`/`TestE2ESchemaChangeIncompatibleTypeChange`)
  und `-run`-Gruppierung laufen nicht auseinander.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — die
  beiden neuen Testfunktionen tragen Test*-Godoc mit dem Testfall selbst
  als Satzsubjekt (zulässige Provenienz-Form, kein Slice-/Wellen-Chronik-
  Verstoß); die neuen Kommentarblöcke in `run-integration-tests.sh`
  beschreiben den geltenden Mechanismus samt Begründung
  (Kopplungs-/Zusage-Klasse), verankert an `ADR-0063`, keine „früher
  stand hier"-Sprache, kein abgebrochener Satz. Keine `slice-062`-Nennung
  im Produktionscode- oder Skript-Kommentar selbst (nur in der
  Commit-Message und in `harness/README.md`s etabliertem
  `· seit slice-<NNN>`-Herkunfts-Anker-Muster).
- geprüft, ohne Befund: `harness/README.md` — `make test-integration`-Zeile
  referenziert korrekt `ADR-0063` (Supersedes `ADR-0058` Entscheidung 1)
  für `TestE2ESchemaChangeDropColumn` und weiterhin `ADR-0058` Entscheidung 2
  für `TestE2EChangeTableMetadataExtensibility` — keine stehengebliebene
  Alleinreferenz auf `ADR-0058` für den korrigierten Testfall; trägt
  `· seit slice-062` im etablierten Muster.
- geprüft, ohne Befund: Slice-Plan-Kopf (`Bezug:`, §1 Ziel, DoD, §3
  Plan-Tabelle) — sauber auf `ADR-0063` nachgezogen, `ADR-0058`-Referenzen
  bleiben nur dort stehen, wo sie weiterhin bindend sind (Entscheidung 2,
  `LH-FA-DAT-006`).
- geprüft, ohne Befund: §6 Risiken — der real während der Umsetzung
  gefundene, nicht vorab benannte Poison-Zustand ist als eigener,
  inhaltlich vollständiger Risiko-Eintrag nachgetragen (Ursache, Befund,
  Behebung, Sonderfall-Einordnung gegenüber dem bestehenden
  `cdc.process_heartbeat`-Präzedenzfall), Ausgang korrekt offen für die
  Closure gelassen (`<bei Closure zu füllen>`).
- geprüft, ohne Befund: Scope-Treue — `git show --stat` bestätigt genau
  die im Plan §3 genannten drei Produktions-/Doku-Dateien plus die
  Plan-Datei selbst; kein Griff in `mapper.go`, `compose.yaml`,
  `.a-check.yml` oder andere Spec-Dateien.
- geprüft, ohne Befund: Traceability — Commit-Betreff trägt
  `LH-FA-SCH-003 … (slice-062)` und `ADR-0063` im Betreff/Body, keine
  `SPEC-*`/`ARC-*`-Kennung im Betreff.
- geprüft, ohne Befund: `make gates` — Exit-Code direkt geprüft (nicht
  gepiped, `AGENTS.md` §3.9), Exit 0: `baseline-verify` (v6.5.0, 54
  Dateien), `docs-check` (2×, inkl. Commit-Range, je 495 Dateien/0 Befunde),
  `commit-traceability` (5 Commits, alle mit Struktur-ID im Betreff
  ausgeschlossen), `coverage-gate` (44.40 % ≥ 35 % Schwelle), `a-check`
  (0 Befunde; zwei bereits bekannte schichtlose Dateien unverändert
  gemeldet, keine neue).
- geprüft, ohne Befund: `make test-integration` — **zwei** eigene,
  vollständige, unabhängige Läufe, jeweils Exit-Code direkt geprüft (nicht
  gepiped): Lauf 1 Exit 0, Lauf 2 Exit 0. Beide Läufe zeigen
  `--- PASS: TestE2EChangeTableMetadataExtensibility`, `--- PASS:
  TestE2ESchemaChangeDropColumn`, die Recovery-Log-Zeile mit
  `tbl-e2e-schema-v4` und anschließend `--- PASS:
  TestE2ESchemaChangeIncompatibleTypeChange` — kein Flake über beide Läufe
  hinweg. Nach Lauf 1 keine verwaisten `cdc-test-*`-Container (`docker ps
  -a` leer).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Schema-Version-Nachtrag nach
relationOther-Abbruch ohne Produktions-Administrationsweg ·
Einzelabfrage-Container-Status-Check ohne Retry

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW; beide INFOs
erwarten keine Implementer-Aktion.

**Übergabe:** Keine Fixrunde nötig. Gemäß `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde zieht dieser Report die DoD-Zeile
„Review durchgeführt …" im selben Commit, der ihn anlegt, auf `[x]` nach.
Beide INFO-Findings gehen als Hinweis an die Planner-Rolle für die
Slice-Closure (§7/Beobachtungs-Register-Erwägung), keine Rückkante an den
Implementer. Der Recovery-Mechanismus selbst ist geprüft **korrekt** —
technisch begründet (nicht nur empirisch einmalig funktionierend) und
über zwei unabhängige, eigene volle `make test-integration`-Läufe ohne
Flake reproduziert. Dieser Report ist ein Lauf-Beleg (Audit: dieser Diff,
dieser Skill, dieses Modell, dieses Verdikt) und wird über Läufe hinweg
nicht wieder gelesen. Verifikation gegen DoD/Spec bleibt Aufgabe des
Verifiers (Modul 11).
