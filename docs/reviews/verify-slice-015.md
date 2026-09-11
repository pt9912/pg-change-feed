# Verifier-Report: slice-015 — 2026-09-11

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code), §6 (Risiko-Vorschlag, Ausgang bleibt
Planner-Entscheidung) und Entscheidungs-Konformität gegen
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
(Re-Evaluierungs-Trigger). Nicht geprüft: Diff gegen Plan/Hard Rules im
Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe, bereits erledigt,
siehe [`review-slice-015.md`](review-slice-015.md)), realer Bedarf
(Validator).

**Gegenstand:** `3cb0c8e` (Implementierung: Pin-Bump + CHECK-Retirement),
`0f00b52` (Review-Report), `ceaf464` (Review-Fixrunde F-1/F-2) — einziger
Implementierungscommit ist `3cb0c8e`; `ceaf464` ändert ausschließlich den
Slice-Plan (Prosa-Korrektur, kein Code).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren, inklusive einer eigenen,
frischen Testcontainer-Rollout-Probe gegen den gepinnten d-migrate-1.3.0-
Digest (unabhängig vom Implementer-Lauf, eigene `schema migrate`-/`schema
compare`-Aufrufe gegen eine selbst hochgefahrene Postgres-Instanz) und
einer eigenen erschöpfenden `--help`-Prüfung aller Subcommands gegen das
gepinnte Image. Alle selbst angelegten Testcontainer, -netze und
Temp-Dateien wurden nach dem Lauf vollständig entfernt (`docker compose
down` + manuelles `docker rm`/`network rm` für den außerhalb des
Compose-Projektnamens laufenden Container; `git status` nach dem Lauf
sauber — kein Rest im Arbeitsbaum).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-11

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand (`in-progress/slice-015-…md`, nach
  `ceaf464`)
- `review-slice-015.md` (F-1 MEDIUM, F-2 LOW, beide committet `0f00b52`,
  disponiert `ceaf464`, kein HIGH)
- [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) im
  Volltext (Re-Evaluierungs-Trigger: kein Werkzeugwechsel, solange eine
  Ausweichform existiert)
- `harness/conventions.md` (MR-000, Sub-Area `PGC`, Greenfield),
  `harness/README.md` (Sensors-Tabelle, `schema-rollout`-Bindungszeile)
- Code im Volltext: `Makefile` (Zeilen 73–120), `tools/schema/schema.yaml`,
  `tools/schema/nacharbeit-views.sql`, `tools/schema/nacharbeit-roles.sql`,
  `tools/schema/nacharbeit-observability.sql`,
  `tools/schema/nacharbeit-heartbeat.sql`, `tools/schema/plan.yaml`,
  `tools/schema/down.sql`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make schema-validate` | zieht den gepinnten Digest real aus der Registry (`Status: Downloaded newer image for …@sha256:d8dc38c…`), `Validation passed: 0 warning(s)`, 8 Tabellen/29 Spalten/7 Constraints | **0** |
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 164 Datei(en) geprüft, 0 Befund(e)` (voll und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | **0** |
| `make test-integration` (Compose, echter Rollout gegen frische Postgres-Instanz) | `schema migrate --execute` grün (CHECK real erzeugt), alle drei Views via `nacharbeit-views.sql` angelegt, 4/4 Go-Integrationstests grün | **0** |
| Eigene `docker run … --help` je Subcommand (Top-Level, `schema`, `validate`, `generate`, `compare`, `reverse`, `migrate`, `rollback`) gegen den gepinnten Digest, `grep -i sandbox` über die gesammelte Ausgabe | **kein einziges Vorkommen** von „sandbox" in irgendeiner Hilfe-Ausgabe | — |
| `docker run … --version` | `d-migrate version 1.3.0` — Digest und Versionsangabe stimmen überein | — |
| Eigener frischer Rollout: `schema.yaml` unverändert gegen leere DB, danach `schema migrate --dry-run` erneut gegen dieselbe (bereits migrierte) DB | Plan schlägt **keine** Operation für `chk_change_operation`/Tabelle `change` vor — nur die drei (dort nicht deklarierten) Views werden als `DropView` vorgeschlagen (Exit 8, `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`, erwartet, da die Views bewusst außerhalb von `schema.yaml` bleiben) | **8** (erwartet, Views-Drop-Bestätigung — kein Konstraint-Op für `change`) |
| Eigener Rollout mit den drei Views **testweise im `views:`-Knoten deklariert** (Autorenform, `active_tables` exemplarisch) gegen eine leere DB | `schema migrate --execute` bricht mit **exakt Exit 5** ab: `[ERROR] Post-execute compare detected drift; the target does not match the desired schema.` | **5** |
| Eigener `schema compare` (Standalone-Tool) der unveränderten `tools/schema/schema.yaml` gegen die bereits erfolgreich migrierte DB | zeigt `~ constraint chk_change_operation (check) -> chk_change_operation (check)` als "geändert" — **abweichend** von der `schema migrate`-Planungsebene (siehe Einordnung unten) | **1** |

**Nicht selbst neu gebaut:** `make image` (kein Gate, in diesem Diff nicht
berührt).

## Sandbox-Mode-Behauptung — eigenständig verifiziert

**Bestätigt: kein „Raw SQL Sandbox Mode"-Flag existiert.** Eigene,
erschöpfende `--help`-Prüfung (Top-Level plus alle sieben `schema`-
Subcommands) gegen den exakt gepinnten Digest
(`sha256:d8dc38cfe3315ab4a1df9f366b262700f402a531b62abd27eb836d3bc4c1e5cc`,
Version `1.3.0`) — kein Treffer für „sandbox" in irgendeiner
Hilfe-Ausgabe. Der einzige mit dem Changelog-Begriff „Post-compare drift
eliminated" tatsächlich zusammenhängende Mechanismus ist die *interne*
Post-execute-Vergleichsebene von `schema migrate --execute` (siehe
unten) sowie `--provenance-output`/`--migration-overlay` (explizit
Out-of-Scope, §1) — nicht ein eigenständiges Flag. Die Implementer-
Behauptung ist damit nicht nur plausibel, sondern **unabhängig
reproduziert**.

## CHECK-Konvergenz und Views-Grenze — eigenständig reproduziert

Zwei getrennte technische Ebenen, die im Diff nicht auseinandergehalten
werden, aber für das Verdikt wichtig sind:

1. **`schema migrate`-Planungs-/Execute-Ebene** (das, was der DoD-Punkt
   und `make schema-rollout` tatsächlich nutzen): Ein `--dry-run` gegen
   eine bereits mit `schema.yaml` migrierte DB schlägt **keine** Operation
   für `chk_change_operation` vor — die Engine erkennt den Constraint als
   konvergent. Ein frischer `--execute`-Lauf mit den drei Views
   **testweise im `views:`-Knoten deklariert** bricht reproduzierbar mit
   **Exit 5** („Post-execute compare detected drift") ab — exakt der vom
   Implementer berichtete Befund, von mir unabhängig in einer separaten
   Testcontainer-Instanz nachvollzogen (nicht dieselbe Session, andere
   Rollout-Reihenfolge, gleiches Ergebnis).
2. **`schema compare`-Standalone-Tool:** Hier zeigt sich, dass ein reiner
   Text-Vergleich Autorenform-YAML-gegen-Katalog auch für den bereits
   konvergenten `chk_change_operation` weiterhin eine kosmetische
   „geändert"-Markierung ausgibt (`before`/`after` identisch benannt, kein
   Ausdruckstext im Diff). Das ist **nicht** derselbe Codepfad wie die
   `migrate`-Planungsebene und wird vom Plan-Text auch nicht behauptet zu
   sein — die DoD-Formulierung spricht konkret von „Post-compare drift"
   im Rollout-Kontext (`schema migrate --execute`, Exit 0 vs. Exit 5),
   nicht vom Standalone-Compare. Kein Widerspruch zur DoD-Aussage, aber
   eine Nuance, die im Plan-Text nicht sichtbar ist: Wer nur `schema
   compare` benutzt, sähe weiterhin „Drift" beim CHECK-Ausdruck. Notiert,
   kein Blocker.

**Fazit:** Beide Kernaussagen aus DoD-Punkt 2 — CHECK konvergiert
deklarativ (Rollout-Ebene), Views bleiben Exit-5-Ausweichform — sind real
und unabhängig reproduziert, nicht nur behauptet.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Pin v1.3.0, `make schema-validate` grün | **bestätigt** | eigener Lauf, Digest-Pull aus Registry verifiziert den Digest im `Makefile`, `Validation passed` |
| 2 | CHECK+Views real gegen neuen Pin getestet, je Fall überführt/dokumentiert | **bestätigt** | eigene, unabhängige Reproduktion beider Teilergebnisse (siehe Abschnitt oben); Sandbox-Mode-Nichtexistenz eigenständig verifiziert |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, Exit 0, alle vier inneren Gates |
| 4 | Review durchgeführt, Report liegt vor | **materiell bestätigt, Checkbox-Inkonsistenz** | `review-slice-015.md` liegt vor (`0f00b52`), 0 HIGH, F-1/F-2 disponiert (`ceaf464`) — die DoD-Checkbox in §2 ist aber weiterhin `[ ]` (nicht abgehakt), obwohl die Arbeit real erledigt ist. Siehe Finding V-1 unten |
| 5 | Doku-Update falls öffentlicher Vertrag berührt | **bestätigt, mit Einschränkung** | `harness/README.md`-Bindungszeile für `schema-rollout` nennt korrekt keine `nacharbeit-*.sql`-Dateinamen — Item entfällt wie geprüft. **Aber:** drei *andere* Dateien tragen eine jetzt tote Querreferenz, siehe Finding V-2 |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin Platzhalter — Planner-Closure-Arbeit, läuft nach diesem Report |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo durchgehend GF, `harness/conventions.md`) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | keine neue `evidence/slice-015.md` unter `BEO-PGC/d-migrate-nacharbeit/` — Planner-Closure-Arbeit. Anzumerken: Der Zähler dieser Beobachtung stand vor diesem Slice bei 2× (CHECK slice-006, Views slice-010); dieser Slice liefert für den CHECK-Fall den *Auflösungs*-Beleg (Ausgang „eingetreten" mit Retirement), keinen dritten Treffer der offenen Klasse — die Views-Hälfte bleibt unter der Schwelle |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen, mit Verifier-Beobachtung** | Risiko 1 (partielles Retirement): eingetreten — Commit-Botschaft `3cb0c8e` benennt exakt diesen Fall explizit („Retirement damit teilweise"), von mir unabhängig bestätigt (CHECK gelöst, Views nicht). Risiko 2 (Regressions-Fund durch Pin-Bump): kein Regressions-Fund in `make gates`/`make test-integration`/meinen eigenen Rollout-Proben — Ausgang „entfallen" ist plausibel gedeckt. Beides bleibt **Vorschlag**, Ausgangs-Zuweisung ist Planner-Aufgabe |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | wellenlos — bei dieser Closure fällig, nicht in diesem Lauf; kein `git mv` bisher, Slice bleibt in `in-progress/` |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, nicht nur behauptet), 1 korrekt entfallen
(Item 7), 4 korrekt noch offen als Planner-Closure-Arbeit (Items 6, 8, 9,
10). Zwei eigene Zusatz-Findings (V-1, V-2), keiner davon DoD-Blocker im
Sinn eines materiell unbelegten Punkts — beide sind Konsistenz-/
Vollständigkeits-Lücken, die vor `git mv` nach `done/` behoben oder
bewusst benannt werden sollten.**

## Plan-vs-Code-Diff (Commit `3cb0c8e`, gegen Plan-§3)

```
git diff --stat 3cb0c8e~1..3cb0c8e -- Makefile tools/schema/
```

liefert genau die sechs Code-/Artefakt-Dateien, die §3 nennt: `Makefile`,
`tools/schema/schema.yaml`, `tools/schema/nacharbeit-operation-check.sql`
(gelöscht), `tools/schema/nacharbeit-views.sql`,
`tools/schema/plan.yaml`, `tools/schema/down.sql`. **Vollständige
Deckung — keine unangekündigte Datei.** Der Review-Report
(`docs/reviews/review-slice-015.md`, `0f00b52`) und die Plan-Prosa-Fixes
(`ceaf464`) zählen nicht als Liefer-Punkt (Baseline-Regelwerk, §2-Hinweis
im Plan). Die Makefile-Änderung deckt beide in §3 genannten Zeilen ab:
Digest-Bump *und* Entfernen des `nacharbeit-operation-check.sql`-
psql-Schritts im `schema-rollout`-Target — eigene Prüfung des aktuellen
`Makefile`-Inhalts bestätigt, dass der Schritt tatsächlich weg ist.

**Deckt sich der Diff mit §3? Ja, vollständig.**

## Befunde (eigene, aus diesem Lauf)

### V-1 — DoD-Checkbox „Review durchgeführt" bleibt unabgehakt, obwohl die Arbeit vorliegt

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-015-d-migrate-1.3.0-retirement.md:106`
- `befund`: Der Review-Report existiert, ist committet (`0f00b52`) und
  seine beiden Findings sind disponiert (`ceaf464`) — die Arbeit dieses
  DoD-Punkts ist real erledigt. Die Checkbox selbst wurde in keinem der
  beiden Folgecommits auf `[x]` gesetzt und steht weiterhin auf `[ ]`.
  Kein Sensor prüft DoD-Checkboxen gegen den tatsächlichen Zustand
  (verifizierbar: nein).
- **Für die Closure:** trivial zu beheben (Checkbox setzen), gehört aber
  vor den `git mv` nach `done/` — sonst dokumentiert der Plan selbst einen
  falschen Zwischenstand.

### V-2 — Drei Nacharbeit-Dateien tragen eine jetzt tote Querreferenz auf die gelöschte `nacharbeit-operation-check.sql`

- `kategorie`: MEDIUM
- `pfad`: `tools/schema/nacharbeit-observability.sql:2`,
  `tools/schema/nacharbeit-roles.sql:2`,
  `tools/schema/nacharbeit-heartbeat.sql:2` — jeweils Kopfkommentar „…
  (`ADR-0043`, s. o. `tools/schema/nacharbeit-operation-check.sql`): …"
- `befund`: Der Implementer-Commit (`3cb0c8e`) löscht
  `tools/schema/nacharbeit-operation-check.sql`, lässt aber drei *andere*,
  im Diff nicht berührte Nacharbeit-Dateien mit einem Kopfkommentar
  stehen, der auf genau diese Datei als „s. o." (Referenzpunkt für das
  Nacharbeit-Muster) verweist. Der Reviewer-Negativbefund „keine dangling
  Referenzen auf die gelöschte Datei … keine in aktiv geltender Doku"
  prüfte `harness/README.md`, `AGENTS.md` und das `schema-rollout`-Target
  selbst — nicht die Kopfkommentare der drei anderen, unveränderten
  `nacharbeit-*.sql`-Dateien, die ebenfalls „aktiv geltender" Code sind
  (sie laufen im selben `schema-rollout`-Target). Kein Gate erkennt das:
  `make gates`/d-check prüft laut `harness/README.md` §Sensors kaputte
  Referenzen in der **Markdown**-Doku, keine Kommentar-Querverweise
  innerhalb von `.sql`-Dateien — dieselbe Deckungslücke wie bei Review-
  Finding F-2 (dort: YAML-Kommentar), hier: drei SQL-Dateien, außerhalb
  des geprüften Diffs.
- `verifizierbar`: nein — kein Gate prüft `.sql`-Kommentar-Querverweise.
- `klasse`: Stale Querverweis nach Löschung des Referenzziels, aus dem
  Diff-Scope herausfallend (Review sah nur die geänderten Dateien) — 2.
  Auftreten derselben Grund-Klasse wie Review-F-2 (dort: „Stale
  Querverweis nach Teiländerung am Zielabsatz"), hier auf drei Dateien
  gleichzeitig.
- **Für die Closure:** kein DoD-Blocker (die drei Dateien bleiben
  inhaltlich korrekt und lauffähig, nur der Verweis ist tot), aber ein
  legitimer Kandidat für den Steering-Loop-Eintrag in §7 — sowohl als
  Korrektur-Nachzug als auch als Beleg dafür, dass ein Review-Scope, der
  sich auf den Diff beschränkt, Querverweise außerhalb des Diffs
  strukturell nicht sieht.

## Negativbefunde

- geprüft, ohne Befund: **[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Konformität** — Re-Evaluierungs-Trigger
  feuert nicht (Operation jetzt ausdrückbar für CHECK, Ausweichform
  weiterhin vorhanden für Views); kein Folge-ADR fällig — eigene
  Bestätigung deckungsgleich mit dem Reviewer-Negativbefund
- geprüft, ohne Befund: **Digest-Korrektheit** — `Makefile:73` trägt
  exakt den Digest, den `docker buildx`/`docker run` beim Pull auflöst
  (`sha256:d8dc38c…`), kein Drift zwischen Pin und tatsächlich gezogenem
  Image
- geprüft, ohne Befund: **Hard Rule 3.3** — kein `git mv` in diesem Range,
  ein Commit mit Inhaltsänderung ist korrekt (kein Lifecycle-Übergang)
- geprüft, ohne Befund: **Traceability** — `make gates` bestätigt alle
  Commits im Fenster mit `LH-*`/`ADR-*`-Bezug, kein `SPEC-*`/`ARC-*` im
  Betreff
- geprüft, ohne Befund: **Arbeitsbaum nach allen Sensor-/Probenläufen** —
  `git status` nach diesem Verifikationslauf sauber; alle selbst
  angelegten Testcontainer/-netze/-dateien entfernt

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (V-2) |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (`make schema-validate`, `make gates`, `make
test-integration`, plus eigene unabhängige d-migrate-Reproduktion beider
Kernaussagen aus DoD-Punkt 2), 1 Item korrekt entfallen
(Reconciliation-Register), 4 Items regulär noch offen als
Planner-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgang, drei Paarungen). **Kein DoD-Defekt im Sinn eines
unbelegten „bestätigt"-Punkts** — beide eigenen Findings (V-1, V-2) sind
Konsistenz-Lücken außerhalb dessen, was die DoD-Checkliste selbst als
Liefer-Punkt zählt, aber real und vor Closure zu beheben oder bewusst zu
benennen.

## Verdikt

**DoD-/Entscheidungs-Konformität: im Kern bestätigt, mit zwei offenen
Klein-Findings vor Closure.** Kein ADR-Verstoß, kein
Traceability-/ID-Schema-Verstoß, keine Halluzination in den
Implementer-/Reviewer-Behauptungen — insbesondere die „Sandbox Mode
existiert nicht"-Feststellung und die CHECK-/Views-Konvergenz-/Drift-
Behauptungen sind **eigenständig reproduziert**, nicht nur übernommen.

**Plan-vs-Code-Diff:** deckt sich vollständig mit §3 — keine
unangekündigte Datei, keine fehlende.

**Vor `git mv` nach `done/` zu klären (Planner):**

1. DoD-Checkbox Punkt 4 nachziehen (V-1, trivial).
2. V-2 disponieren — Korrektur-Nachzug der drei toten Querverweise, oder
   bewusst als benannte Grenze in §7 vermerken (nicht automatisch
   prüfbar, siehe oben).
3. §6-Risiken disponieren — Verifier-Empfehlung: Risiko 1 „eingetreten"
   (Teil-Retirement, CHECK gelöst/Views nicht, real belegt), Risiko 2
   „entfallen" (kein Regressions-Fund in irgendeinem real gefahrenen
   Sensor dieses oder des Verifikationslaufs) — **reine Empfehlung, kein
   gesetzter Ausgang.**
4. §7 Closure-Notiz, Beobachtungs-Register-Fortschreibung
   (`BEO-PGC/d-migrate-nacharbeit`, CHECK-Hälfte jetzt aufgelöst) und die
   drei Paarungen — reguläre Planner-Closure-Arbeit, in diesem Lauf nicht
   fällig.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; alle Testcontainer/-netze
und Temp-Dateien dieses Laufs wurden vollständig entfernt.

---

**Gate-Beleg:** `make schema-validate` und `make gates` in diesem Lauf,
beide Exit 0 (siehe Sensor-Tabelle oben).
