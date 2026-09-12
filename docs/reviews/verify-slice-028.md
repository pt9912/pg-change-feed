# Verifier-Report: slice-028 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte, aktueller Stand nach Planner-Korrektur `ee627bd`), §3
(Plan-vs-Code inkl. Plan-Nachzug), §6 (Risiko-Ausgänge), §8
(Sub-Area-Prüfung), sowie Entscheidungs-Konformität gegen
[`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) und die
Reviewer-Findings aus [`review-slice-028.md`](review-slice-028.md) (1
MEDIUM F-1, 1 INFO F-2). Nicht geprüft: Diff gegen Plan/Hard Rules im
Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe, bereits erledigt),
realer Bedarf (Validator — hier nicht einschlägig, reine
Testinfrastruktur ohne Verhaltensänderung).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen oder ausgeführt: `tools/harness/run-integration-tests.sh`
(neuer Abschnitt „Rollen-DSN-Verifikation gegen den Compose-Stack",
Volltext), `tools/schema/nacharbeit-roles.sql` (Grants je Rolle),
`compose.yaml` (Env-Block, Healthcheck), `docs/reviews/review-slice-023.md`
INFO-1 (Volltext), `git log --follow -- compose.yaml`/`git show a32a2c9`,
`git show --stat -M` für den reinen `git mv`, `grep -rn mvp_test\.go` über
das gesamte Repo, `make gates` selbst gestartet, `make test-integration`
**dreimal** selbst gestartet (nicht nur die Implementer-/Reviewer-Behauptung
übernommen — der Reviewer hatte den Lauf ausdrücklich nicht reproduziert),
`git status`/`git diff` am Ende sauber.

**Gegenstand:** `632ddca` (git mv), `ba508ed` (Rollen-DSN-Verifikation),
`9a84407` (MVP-Sprachgebrauch-Nachzug), `f8ce37b` (DoD-Nachzug),
`ffe9ac8` (Review-Report), `ee627bd` (Planner-Korrektur nach F-1).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-028-integrationstest-nachzug-rollen-mvp.md`, nach `ee627bd`)
- `docs/reviews/review-slice-028.md` (1 MEDIUM F-1, 1 INFO F-2)
- `tools/harness/run-integration-tests.sh` (Volltext, 470 Zeilen, inkl. des
  neuen Abschnitts Zeilen 71–147)
- `tools/schema/nacharbeit-roles.sql` (Grant-Tabelle je Rolle)
- `compose.yaml` (Env-Block `pg-change-feed`, Zeilen 53–74)
- `docs/reviews/review-slice-023.md` (INFO-1, Volltext)
- `git log --follow -- compose.yaml`, `git show a32a2c9 --stat`
- `docs/plan/planning/observations/BEO-PGC/rollen-test-abdeckungsluecken/`
  (`observation.md`, `state.md`, `evidence/slice-023.md`)
- `harness/README.md` (`make test-integration`-Zeile), `Makefile`
  (Helptext), `AGENTS.md` (Grep auf „MVP")
- `test/integration/integration_test.go` (Paket-Kommentar)
- `git log`/`git show`/`git diff` über den vollen Commit-Verlauf dieses
  Slice

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (2×, vor und nach den drei Integrationsläufen) | `baseline-verify`: 54 Dateien OK · `d-check` Standardlauf: 250 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD`: 250 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-integration` (Lauf 1) | Rollen-DSN-Verifikation „belegt" (drei Login-Identitäten korrekt akzeptiert/abgelehnt), vier Go-Integrationstests `PASS`, Lasttest-Beleg, Black-Box-CLI-Rundlauf belegt | **0** |
| `make test-integration` (Lauf 2) | identischer grüner Ablauf | **0** |
| `make test-integration` (Lauf 3) | identischer grüner Ablauf | **0** |
| `git status`/`git diff` (am Ende dieses Laufs) | sauber, keine Restspur | — |

Drei unabhängige grüne Läufe bestätigen die DoD-Behauptung „dreimal in
Folge grün" **selbst** — der Reviewer hatte diesen Lauf ausdrücklich nicht
reproduziert (`review-slice-028.md`, Negativbefund-Zeile „nicht
reproduziert in dieser Review-Sitzung"). Keine Anzeichen von Flakiness
über die drei Läufe (Timing-Werte für `cdc_capture_lag` variierten nur
geringfügig, alle drei innerhalb der erwarteten Toleranzen).

## Prüfpunkt 0 — Ist die Planner-Korrektur (`ee627bd`) selbst akkurat?

**Ja, akkurat — eigenständig am Code nachvollzogen, nicht nur an der
Reviewer-Aussage.**

Eigene Volltext-Lektüre von `tools/harness/run-integration-tests.sh`
Zeilen 71–147 (der komplette neue Abschnitt):

- Der Abschnitt legt drei ephemere Login-Identitäten per rohem SQL an
  (`CREATE ROLE ... LOGIN PASSWORD ... IN ROLE cdc_reader/cdc_admin/cdc_capture`)
  und führt genau drei Prüfungen über `docker exec "$PG_CONTAINER" psql`
  aus: (1) ein `INSERT` über die `cdc_reader`-Identität, erwartet
  `permission denied` (Zeile 112–120); (2) eine Replication-Protokoll-
  Verbindung (`replication=database`, `IDENTIFY_SYSTEM;`) über die
  `cdc_admin`-Identität, erwartet `permission denied to start WAL sender`
  (Zeile 122–131); (3) dieselbe Replication-Verbindung über die
  `cdc_capture`-Identität mit explizit gesetztem `REPLICATION`-Attribut,
  erwartet Erfolg (Zeile 133–142).
- **Kein einziger Aufruf** von `receive.NewStream`, `postgresack.New`,
  keines der internen Go-Adapter-Pakete (`internal/...`) und kein
  Go-Testbinary, das diese Konstruktoren nutzen würde, taucht in diesem
  Abschnitt auf — die Prüfung läuft ausschließlich über `docker exec …
  psql` gegen PostgreSQL direkt. Bestätigt per `grep -n "receive\.\|postgresack\."
  tools/harness/run-integration-tests.sh` (kein Treffer im gesamten
  Skript).
- Die Aussage „nur PostgreSQL-Ebene, nicht Adapter-Ebene" trifft die
  Sache **präzise** — sie ist weder zu stark noch zu schwach formuliert.
  Zu stark wäre sie, wenn sie suggerierte, mehr als PostgreSQLs
  serverseitige Durchsetzung sei geprüft (das ist nicht der Fall — die
  drei Login-Identitäten sind eigens angelegte Test-Fixtures, keine
  Verbindung des laufenden Feed-Containers). Zu schwach wäre sie, wenn
  die Prüfung tatsächlich noch weniger zeigte als „PostgreSQL verweigert
  die Verbindung/den Schreibzugriff" — das ist nicht der Fall: Die drei
  Assertions sind real, laufen gegen den echten Compose-Stack (nicht
  gemockt) und prüfen tatsächlich das für `ADR-0047` Kontext-Befund 2
  relevante Verhalten (Attribut wird nicht über Mitgliedschaft vererbt —
  die `cdc_capture`-Kontrastprüfung setzt das Attribut deshalb explizit
  auf der Login-Rolle, nicht nur die Mitgliedschaft).
- Eigenständig gegen `tools/schema/nacharbeit-roles.sql` geprüft: `cdc_reader`
  trägt ausschließlich `GRANT SELECT` auf die drei Lese-Views (Zeile 98),
  kein `INSERT`-Recht auf `cdc.change` — die erwartete Ablehnung ist
  grant-seitig gedeckt. `cdc_admin` ist `NOLOGIN` ohne `REPLICATION`-
  Attribut (Zeile 46) — die erwartete Ablehnung der Replication-Verbindung
  ist rollenseitig gedeckt. `cdc_capture` trägt `REPLICATION` auf
  Rollenebene (Zeile 43), aber das Skript setzt es zusätzlich explizit auf
  der Login-Identität (`ALTER ROLE ... REPLICATION`, Zeile 108) — konsistent
  mit dem im Skript-Kommentar selbst benannten Nicht-Vererbungs-Verhalten.

**Verdikt Prüfpunkt 0:** Die Korrektur ist nicht nur vorhanden, sondern
präzise. Sie ist auch nicht zu vorsichtig — es gibt keinen Hinweis, dass
die tatsächliche Testabdeckung schwächer wäre, als „PostgreSQL-Ebene,
nicht Adapter-Ebene" impliziert.

## Prüfpunkt F-2 — Superuser-Verdrahtung, vorbestehend und korrekt klassifiziert?

**Ja, eigenständig bestätigt.**

- `compose.yaml` Zeilen 62–64: `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/
  `CDC_READER_DSN` zeigen alle drei identisch auf
  `postgres://postgres:postgres@cdc-test-postgres:5432/cdc?sslmode=disable`
  — real bestätigt, kein Lesefehler des Reviewers.
- `git log --follow --oneline -- compose.yaml` zeigt `a32a2c9` als
  **jüngsten** Commit, der `compose.yaml` berührt — vor allen vier
  slice-028-Commits. `git show a32a2c9 --stat` bestätigt `compose.yaml`
  in der Dateiliste dieses Commits. Keiner der vier slice-028-Commits
  (`632ddca`/`ba508ed`/`9a84407`/`f8ce37b`) fasst `compose.yaml` an
  (eigene `git show --stat` je Commit, oben bereits protokolliert).
- `docs/reviews/review-slice-023.md` INFO-1 (Zeilen 155–173) dokumentiert
  exakt denselben Zustand mit derselben Begründung (`ADR-0047`
  Kontext-Befund 1) bereits am `slice-023`-Commit — Wortlaut und Klasse
  („Least-Privilege nur in Unit-nahen Tests, nicht im End-to-End-Lauf
  belegt") stimmen mit `review-slice-028.md` F-2 überein.

**Verdikt Prüfpunkt F-2:** Kein neuer Befund dieses Slices, korrekt als
INFO ohne Eskalationsbedarf eingestuft. Keine Regression.

## Prüfpunkt — Widerspruch zwischen §8 (nicht nachgezogen) und §2/§6 (korrigiert)?

`docs/plan/planning/in-progress/slice-028-integrationstest-nachzug-rollen-mvp.md`
§8, Block „Vorgelagert — offene Beobachtungen sichten" (Zeile 281–286)
trägt weiterhin den vor der Planner-Korrektur formulierten Satz „dieser
Slice schließt Punkt (2) davon" — unverändert seit dem Anlegen des Slice,
nicht von `ee627bd` nachgezogen.

**Einschätzung: legitime Momentaufnahme, kein Defekt, aber ein
benennbarer Klein-Befund (V-1, siehe unten).** Nach Baseline-Regelwerk
`modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung
laufen die zwei vorgelagerten Prüfungen (Sub-Area-Wahl, offene
Beobachtungen sichten) **beim Anlegen** des Slice — sie sind eine
Planungs-Absicht vor der Umsetzung, keine laufend nachgezogene
Statusfläche. Der Slice-Kopf selbst trägt bereits den korrekten Verweis
(„real erreicht wurde nur die PostgreSQL-Ebene, nicht die Adapter-Ebene —
siehe §2 DoD und §7 Closure-Notiz für die tatsächliche Deckung", §1) und
weist damit explizit auf die tragende, aktuelle Quelle hin. §2 (DoD) und
§6 (Risiken) — die beiden Abschnitte, die laut Regelwerk die tatsächliche
Deckung tragen — sind korrekt und konsistent aktualisiert; §6 trägt sogar
bereits das neue „weiter offen"-Risiko mit Verweis auf denselben Zähler-
Stand (2×), den §8 nennt. Kein Sensor und keine nachgelagerte Rolle liest
§8 als Nachweis der tatsächlichen Deckung — dafür sind §2/§6/§7
zuständig, und die sind stimmig.

Dennoch: Ein Leser, der nur §8 liest (etwa bei der `welle-8`-Closure, die
laut Regelwerk primär das Beobachtungs-Register selbst konsultiert, nicht
§8 einzelner Slices), könnte den veralteten Satz für die Deckungsaussage
halten. Das ist der Gegenstand von V-1.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt, nach `ee627bd`)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-QA-SEC-001`…`003` teilweise erfüllt, Punkt (2) von `BEO-PGC/rollen-test-abdeckungsluecken` **nicht vollständig** geschlossen (Planner-Korrektur) | **bestätigt, Korrektur akkurat** | Prüfpunkt 0 oben — eigene Volltext-Lektüre, kein Aufruf von `receive.NewStream`/`postgresack.New`, Grants gegen `nacharbeit-roles.sql` verifiziert, drei eigene grüne Läufe |
| 2 | `mvp_test.go` → `integration_test.go`, Doku/Helptext nachgezogen | **bestätigt** | `git show --stat -M 632ddca`: reiner Rename (0/0), Similarity 100 % erkannt. `grep -rn mvp_test\.go` repo-weit: nur noch in historischen `done/`-/Review-/Verify-Dateien und im eigenen Vor-/Nachher-Text von `slice-028` selbst — keine lebende Referenz. `Makefile` Zeile 56, `harness/README.md` Zeile 129 tragen den neuen Wortlaut. `AGENTS.md` trägt keine „MVP"-Erwähnung (eigener `grep`, kein Treffer) — „kein Änderungsbedarf" bestätigt |
| 3 | `make gates` grün, `make test-integration` dreimal grün | **bestätigt, dreifach selbst ausgeführt** | Sensor-Tabelle oben: drei eigene `make test-integration`-Läufe grün, `make gates` zweimal grün mit 0 Befunden |
| 4 | Review durchgeführt, Report liegt vor, kein Self-Review | **bestätigt** | `docs/reviews/review-slice-028.md` existiert (`ffe9ac8`), 1 MEDIUM/1 INFO, Verdikt „merge-blockierend: ja (F-1)" — durch Planner-Korrektur (`ee627bd`) behoben; Checkbox korrekt `[x]` mit Verweis auf die Korrektur |
| 5 | Doku-Update `harness/README.md`/`AGENTS.md` | **bestätigt** | `git show 9a84407 -- harness/README.md`: aktualisierte Zeile deckungsgleich mit tatsächlichem Skript-Inhalt (auch die Rollen-DSN-Verifikation korrekt beschrieben); `AGENTS.md` unberührt, korrekt begründet |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin den Platzhalter-Vorlagentext — Planner-Arbeit, noch nicht fällig vor `git mv` |
| 7 | Reconciliation-Register, falls einschlägig | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (eigene Prüfung); Sub-Area `*`/`PGC` durchgehend Greenfield |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt zurückgestellt, DoD-Zeile bereits vorbereitet** | `evidence/` unter `BEO-PGC/rollen-test-abdeckungsluecken/` trägt aktuell nur `slice-023.md` (eigene `ls`); `state.md` steht noch bei 1× — konsistent mit „Zurückgestellt auf die Closure", DoD-Zeile benennt bereits korrekt den geplanten `evidence/slice-028.md`-Beleg und den neuen Zählerstand (2×) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **materiell bereits erfüllt, Checkbox nicht nachgezogen — siehe V-2** | Alle drei §6-Risiken tragen bereits einen der drei zulässigen Ausgänge (zwei „entfallen" mit Begründung, ein „weiter offen" aus F-1) — die DoD-Checkbox selbst steht weiterhin `[ ]` |
| 10 | Drei Paarungen getragen | **korrekt offen, turnusgemäß bei `welle-8`-Closure** | Slice-Kopf „Welle: welle-8" — Prüfung läuft laut DoD-Zeile selbst bei der nächsten Welle-Closure, nicht bei dieser Einzel-Slice-Closure |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (inkl. eigener Code-Lektüre der Rollen-Verifikation,
eigener Grant-Gegenprobe, dreifachem eigenem `make test-integration`-Lauf
und eigener `git log --follow`-Prüfung des F-2-Ursprungs), 1 Item korrekt
entfallen (Reconciliation-Register, Greenfield), 3 Items regulär offen als
Planner-/Welle-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
Paarungen bei `welle-8`), 1 Item materiell erfüllt mit
Checkbox-Diskrepanz (V-2, wie in `verify-slice-024`…`027` bereits
wiederkehrend).**

## Plan-vs-Code-Diff (§3, gegen `632ddca`/`ba508ed`/`9a84407`/`f8ce37b`)

| Plan-Zeile | Behauptung | Abgleich |
|---|---|---|
| `tools/harness/run-integration-tests.sh` — neuer Abschnitt Rollen-Verifikation | Zugriff mit falscher Rolle real gegen den Compose-Stack zurückgewiesen | **stimmt, mit der von Reviewer/Planner präzisierten Reichweite** — `ba508ed`, Zeilen 71–147 |
| `run-replication-tests.sh` unverändert (Plan-Nachzug) | Compose-Instanz trägt die Gruppenrollen bereits über `schema-rollout`, keine Erweiterung der `wal_level=logical`-Testcontainer-Kette nötig | **stimmt** — `git show --stat` über alle vier Commits: Datei wird von keinem berührt (bereits vom Reviewer bestätigt, hier eigenständig wiederholt) |
| `test/integration/mvp_test.go` → neuer Dateiname | Umbenennung nach tatsächlichem Scope | **stimmt** — reiner Rename, Zeile 2 oben |
| `Makefile` | Helptext-Zeile aktualisiert, Target-Name unverändert | **stimmt** — Zeile 56 |
| `harness/README.md`, `AGENTS.md` | Sensors-Tabelle/Erwähnungen nachgezogen | **stimmt** — `AGENTS.md` korrekt unberührt (keine MVP-Erwähnung vorhanden) |
| `docs/plan/planning/welle-8.md` (Plan-Nachzug) | Dateipfad-Kreuzverweis nachgezogen | **stimmt** — Zeile 22, `integration_test.go` |

`git show --stat` über alle vier Commits zeigt keine unangekündigte
fünfte Datei außerhalb der oben genannten. Out-of-Scope-Disziplin (§1)
eigenständig geprüft: keine neuen Rollen/DDL (`nacharbeit-roles.sql`
unverändert in allen vier Commits — eigene `git show --stat`), kein
zweites Black-Box-Aufrufmuster erfunden (der Rollen-Abschnitt nutzt
ausschließlich `docker exec … psql`, kein neuer CLI-Rundlauf), keine
zusätzliche WAL-Rückstand-Prüfung gegen den Compose-Stack.

## §6 Risiken — eigene Einschätzung zum Verifikationsstand

Alle drei Risiken tragen bereits einen der drei zulässigen Ausgänge, und
die Begründungen sind plausibel und mit dem Diff vereinbar:

1. „Rollenbeschränkte Login-Test-Identitäten könnten zusätzliche
   PostgreSQL-Rollen-Berechtigungen brauchen" → **entfallen** — plausibel:
   Der tatsächliche Ansatz (Prüfung gegen die Compose-Instanz, die die
   Gruppenrollen bereits real trägt) kam ohne neue Berechtigungen aus;
   eigene Lektüre von `nacharbeit-roles.sql` zeigt keine Änderung in
   diesem Slice.
2. „Umbenennung könnte weitere Erwähnungen übersehen" → **entfallen** —
   eigene `grep -rn mvp_test\.go` bestätigt: keine lebende Referenz
   übersehen.
3. „`BEO-PGC/rollen-test-abdeckungsluecken` Punkt (2) nur auf
   PostgreSQL-Ebene geschlossen" (aus F-1) → **weiter offen** — korrekt,
   deckt sich mit Prüfpunkt 0 oben.

Kein Risiko ohne Ausgang, keiner der drei Ausgänge außerhalb der
geschlossenen Menge (eingetreten/entfallen/weiter offen).

## Negativbefunde

- geprüft, ohne Befund: **Rollen-Verifikation, technische Korrektheit** —
  eigene Grant-Gegenprobe gegen `nacharbeit-roles.sql`: `cdc_reader` ohne
  `INSERT`-Recht, `cdc_admin` ohne `REPLICATION`-Attribut, `cdc_capture`
  mit `REPLICATION`-Attribut — alle drei erwarteten Ablehnungen/Erfolge
  sind grant-/rollenseitig gedeckt, nicht zufällig.
- geprüft, ohne Befund: **Kein Aufruf interner Adapter-Konstruktoren** —
  eigener `grep` auf `receive\.\|postgresack\.` im gesamten Skript, kein
  Treffer im neuen Abschnitt.
- geprüft, ohne Befund: **Reiner `git mv`** — `git show --stat -M
  632ddca`: 0 Insertions/0 Deletions, Similarity-Erkennung als Rename,
  kein Inhalt im selben Commit (`AGENTS.md` §3.3 eingehalten).
- geprüft, ohne Befund: **F-2-Herkunft** — eigene `git log --follow`-
  Prüfung bestätigt `a32a2c9` (slice-023) als Ursprung, kein slice-028-
  Commit berührt `compose.yaml`.
- geprüft, ohne Befund: **`make gates`** — zwei eigene Läufe, 0 Befunde
  in allen vier Gates.
- geprüft, ohne Befund: **`make test-integration`, dreimal** — drei
  eigene, unabhängige Läufe, alle grün, keine Flakiness-Anzeichen.
- geprüft, ohne Befund: **Doku-Konsistenz** `harness/README.md`,
  `Makefile`, `AGENTS.md` — Wortlaut gegen tatsächlichen Diff
  gegengelesen, keine Überzeichnung, keine übersehene Stelle.
- geprüft, ohne Befund: **Traceability** — alle sechs Commit-Betreffs
  ohne `SPEC-*`/`ARC-*`, `LH-QA-SEC-001` referenziert; `make gates`
  (inkl. `commit-traceability`) grün.
- geprüft, ohne Befund: **Out-of-Scope-Disziplin §1** — alle vier
  genannten Ausschlüsse im Diff tatsächlich unberührt.
- geprüft, ohne Befund: **Beobachtungs-Register-Konsistenz** — DoD-Text
  und `state.md` stimmen überein (1× bisher, 2× geplant nach Closure),
  kein Widerspruch zwischen angekündigtem und tatsächlichem Zählerstand.

## Eigene Befunde

### V-1 — §8 „Vorgelagert — offene Beobachtungen sichten" nicht an die Planner-Korrektur nachgezogen

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-028-integrationstest-nachzug-rollen-mvp.md:281-286`
  (§8, Satz „dieser Slice schließt Punkt (2) davon") vs. §1/§2/§6 (nach
  `ee627bd`: Punkt (2) **nicht vollständig**, nur PostgreSQL-Ebene)
- `befund`: Siehe Prüfpunkt oben — legitime Momentaufnahme der
  Planungs-Absicht vor der Umsetzung (Baseline-Regelwerk `modul-05` §Zwei
  Schritte vor der Modus-Begründung), aber ein Leser, der ausschließlich
  §8 liest, bekäme eine überholte Aussage. Kein Sensor prüft diesen
  Abschnitt gegen §2/§6; die tragenden Abschnitte sind korrekt.
- `verifizierbar`: ja — Zeile 283 wurde von `ee627bd` nicht verändert
  (`git show ee627bd` zeigt nur §1/§2/§6-Hunks).
- **Für die Closure:** kein Blocker. Empfehlung: bei der Closure (§7) den
  Satz in §8 auf den tatsächlichen Deckungsgrad nachziehen, oder einen
  Verweis auf §2/§7 ergänzen — Konsistenz-Kosmetik, keine
  Sachkorrektur.

### V-2 — DoD-Checkbox „Jedes Risiko aus §6 trägt einen Ausgang" nicht nachgezogen, obwohl materiell erfüllt

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-028-integrationstest-nachzug-rollen-mvp.md:136`
  (Checkbox weiterhin `[ ]`) vs. §6 (alle drei Risiken tragen bereits
  einen Ausgang)
- `befund`: Derselbe wiederkehrende Musterbefund wie in
  `verify-slice-024.md`…`verify-slice-027.md` (V-1 dort): die DoD-Zeile
  ist materiell bereits erfüllt, die Checkbox wird erst bei der formalen
  Closure gesetzt. Diese Finding-Klasse ist bereits in
  `BEO-PGC/dod-checkbox-nachzug` als **verkörpert** geführt (3× erreicht,
  `seit welle-5`) — dieser Fall trägt keinen eigenen Zähler-Beitrag (ein
  Vorgang zählt einmal, Modul 6).
- `verifizierbar`: ja — Checkbox bleibt `[ ]` bis zum aktuellen HEAD.
- **Für die Closure:** kein Blocker — der Planner setzt die Checkbox vor
  `git mv` nach `done/` (materiell gedeckt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 (V-1, V-2) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst nachgeprüft, 1 Item korrekt entfallen (Reconciliation-Register,
Greenfield), 3 Items regulär offen als Planner-/Welle-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register-Fortschreibung, drei Paarungen bei
`welle-8`), 1 Item materiell erfüllt mit Checkbox-Diskrepanz (V-2). Kein
DoD-Defekt im Sinn eines unbelegten „bestätigt"-Punkts.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt.** Die Planner-Korrektur
(`ee627bd`) ist akkurat — eigene Volltext-Lektüre des neuen
Rollen-Verifikations-Abschnitts bestätigt, dass er ausschließlich
PostgreSQLs serverseitige `REPLICATION`-Attribut-Durchsetzung über
ephemere Test-Identitäten prüft, nirgends `receive.NewStream`/
`postgresack.New` aufruft und vom laufenden Feed-Container entkoppelt
bleibt — die Formulierung „nur PostgreSQL-Ebene, nicht Adapter-Ebene" ist
weder zu stark noch zu schwach. F-2 (Superuser-Verdrahtung) ist
eigenständig als vorbestehend seit `slice-023` (`a32a2c9`) bestätigt,
keine Regression dieses Slices. Der nicht nachgezogene §8-Satz (V-1) ist
eine legitime Planungs-Momentaufnahme mit einem kleinen
Konsistenz-Risiko für isolierte Leser, kein Sachfehler — §2/§6/§7 tragen
die tatsächliche, korrekte Deckung. Kein viertes, bisher unbemerktes
Sachproblem gefunden (eigene Prüfung von Grant-Tabelle, Rename-Erkennung,
Doku-Nachzug, Out-of-Scope-Disziplin, `git log --follow`-Herkunft von
F-2).

**Plan-vs-Code-Diff:** vollständige Deckung, keine unangekündigte
Abweichung. Out-of-Scope-Punkte aus §1 gewahrt.

**Closure-Bereitschaft: ja, für die gelieferte Substanz — reguläre
Planner-/Welle-Closure-Schritte stehen noch aus, keiner davon ein Defekt
an bereits gelieferter Substanz:**

1. Closure-Notiz §7 schreiben (Was hat funktioniert / anders als geplant /
   Steering-Loop-Eintrag / Beobachtungs-Register / Folge-Slices /
   §6-Risiken-Zusammenfassung).
2. Beobachtungs-Register: `evidence/slice-028.md` unter
   `BEO-PGC/rollen-test-abdeckungsluecken/` anlegen, `state.md` auf 2×
   fortschreiben (DoD-Zeile benennt dies bereits korrekt vorbereitet).
3. Drei Paarungen: laufen turnusgemäß mit der bevorstehenden
   `welle-8`-Closure — nicht bei dieser Einzel-Slice-Closure.
4. Empfohlen, nicht blockierend: V-1 (§8-Satz) und V-2
   (Risiko-Ausgänge-Checkbox) vor `git mv` nach `done/` nachziehen.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (`git status`/`git diff` am
Ende sauber).

---

**Gate-Beleg:** `make gates` zweimal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `d-check` Standardlauf: 250
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 250
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits;
`a-check`: 0 Befunde). `make test-integration` **dreimal** vollständig
selbst ausgeführt (3× grün, keine Anzeichen von Flakiness) — die
DoD-Behauptung „dreimal in Folge grün" ist damit erstmals in dieser
Sitzung eigenständig reproduziert (der Reviewer hatte dies ausdrücklich
nicht getan). `git status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
