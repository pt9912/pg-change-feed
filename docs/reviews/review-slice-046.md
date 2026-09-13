# Review-Report: slice-046 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-046`, §1/§2/§3/§4/§6/§7/§8),
`welle-13`, den Architect-Verdikt
`docs/reviews/architect-verdict-retention-loeschausfuehrung.md` (Frage 2),
`ADR-0046` sowie `AGENTS.md` §3 Hard Rules (§3.1, §3.3, §3.5, §3.7) —
Rollentrennung Modul 8: diese Prüfung läuft gegen Plan/ADR/Hard Rules
(Maintainability), nicht gegen DoD (Verifier-Aufgabe).

**Gegenstand:** Commit `d4bfde9` (`feat(observability): cdc_storage_bytes-Metrik
in cdc.metrics (LH-FA-RET-006)`) gegen Elter `bc7bba0`; Diff:
`tools/schema/nacharbeit-observability.sql`,
`internal/adapters/driven/postgresstorage/roles_test.go`,
`test/integration/integration_test.go`,
`tools/harness/run-integration-tests.sh`, `docs/user/benutzerhandbuch.md`,
`docs/plan/planning/in-progress/slice-046-cdc-storage-bytes-metrik.md`
(Plan-Nachzug, DoD-Häkchen, §6/§7).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft
2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-046-cdc-storage-bytes-metrik.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug drei
  Punkte, §4 Trigger, §6 Risiken mit Ausgang, §7 Closure-Notiz, §8)
- `docs/reviews/architect-verdict-retention-loeschausfuehrung.md`
  (Frage 2 — View-Owner-Muster, `pg_relation_size()` PUBLIC-ausführbar,
  `ADR-0046`-Kategorie-C-Einordnung)
- `docs/plan/planning/done/slice-043-changestoreport-loeschmethode.md`,
  `.../slice-044-retention-hintergrundjob.md`,
  `.../slice-045-blockierende-consumer-sichtbarkeit.md` (Vorgänger-Slices
  derselben Welle; Referenz für Testmuster und für die Beobachtung
  `BEO-PGC/retention-keine-loeschausfuehrung`)
- `docs/plan/planning/observations/BEO-PGC/retention-keine-loeschausfuehrung/`
  und `.../BEO-PGC/test-runner-stiller-ausschluss/` (Register-Einträge,
  real gegen `observation.md`/`state.md`/`evidence/` geprüft)
- `spec/lastenheft.md` `LH-FA-RET-006` (Happy Path/Boundary)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git-mv + Inhalt = zwei Commits),
  §3.5 (ADR-Immutabilität), §3.7 (Kommentar-Disziplin —
  `BEO-PGC/slice-chronik-in-code-kommentar`, 3×/verkörpert vor diesem Lauf,
  besonders geprüft)
- `git log -1 --format='%s' d4bfde9` (Commit-Traceability)
- `make docs-check` (lokal ausgeführt: 362 Dateien, 0 Befunde)
- `make commit-traceability` (lokal ausgeführt gegen `HEAD~5..HEAD`: OK,
  keine Struktur-ID im Betreff)
- `grep -nE '^\+.*\b(slice-[0-9]{3}|welle-[0-9]+)\b'` über den vollständigen
  `.go`/`.sql`-Diff (vollständiger Durchlauf, kein Treffer)
- `grep`-Lauf gegen `` \[`(LH|ADR|SPEC|ARC|slice|BEO|CO|MR|RC)- `` auf
  hinzugefügten Zeilen der geänderten `.md`-Dateien (kein Treffer)
- `git diff` auf `tools/schema/nacharbeit-roles.sql` (unverändert — reale
  Prüfung der Implementer-Behauptung „keine neue `cdc_reader`-Rollenfläche")

---

## Findings

### F-1 — Beobachtungs-Register: gestrichen-Ausgang an falschen Prozessschritt delegiert

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Das
  Beobachtungs-Register — „Gestrichen ist an [die Schwelle] nicht
  gebunden — fällt die Ursache weg, bevor der Zähler 3 erreicht, wandert
  die Zeile mit Begründung in §Gestrichene Einträge."
- `pfad`: `docs/plan/planning/in-progress/slice-046-cdc-storage-bytes-metrik.md`
  (§7, Absatz „Beobachtungs-Register") vs.
  `docs/plan/planning/observations/BEO-PGC/retention-keine-loeschausfuehrung/state.md`
  (unverändert in diesem Diff)
- `befund`: Die Closure-Notiz stellt selbst fest: „mit diesem Slice sind
  alle vier in der Beobachtung benannten Lücken (Löschausführung,
  Hintergrundjob, Sichtbarkeit blockierender Consumer,
  `cdc_storage_bytes`-Metrik) geliefert" — und genau diese vier Lücken sind
  exakt die in `observation.md`/`state.md` genannten Auflösungskriterien
  („tatsächliche Löschausführung … und die Metrik `cdc_storage_bytes`
  bauen"). Trotzdem bleibt `state.md` in diesem Diff unverändert bei
  „Zustand: offen — Ausgang: weiter offen", und die Closure-Notiz verschiebt
  die Ausgangs-Zuweisung explizit auf den Lese-Schritt der
  `welle-13`-Closure. Der Lese-Schritt (Modul 6) liest jedoch nur Einträge,
  die 3× erreicht haben — dieser Eintrag steht bei 0× (`observation.md`:
  „Benannt, nicht gezählt", kein abgeschlossener Vorgang trägt bisher
  einen Beleg) und wird vom Lese-Schritt strukturell nie erfasst. Der
  Ausgang „gestrichen" hängt laut Regelwerk gerade nicht an dieser
  Schwelle, sondern wird sofort mit Begründung vergeben, sobald die
  Ursache wegfällt.
- `verifizierbar`: nein — kein Gate prüft die semantische Korrektheit
  eines Register-Ausgangs; die Register-Paarung (Modul 6, Schritt 3c)
  prüft nur Existenz/Belegtheit von Verzeichnissen, nicht ob der
  eingetragene Zustand zur Sachlage passt.
- `klasse`: „Register-Ausgang an falschen Prozessschritt delegiert
  (gestrichen vs. Lese-Schritt)"

### F-2 — Commit bündelt Implementierung und Plan-Doc-Closure entgegen etablierter Trennung

- `kategorie`: LOW
- `quelle`: Maintainability (kein Hard-Rule-Verstoß — `AGENTS.md` §3.3
  betrifft nur `git mv` + Inhaltsänderung, hier liegt kein `git mv` vor)
- `pfad`: Commit `d4bfde9` (gesamter Commit)
- `befund`: Anders als bei `slice-044`/`slice-045` (dort: ein `feat(...)`-Commit
  für den Code getrennt von einem eigenen `docs(planning): ... Closure-Inhalt`-
  Commit für DoD-Häkchen/§3/§6/§7) bündelt dieser Commit Code-Änderungen
  (SQL, Go-Tests, Test-Runner-Skript) und Plan-Doc-Closure-Inhalt
  (DoD-Häkchen, Plan-Nachzug, §6-Ausgänge, §7-Notiz) in einem einzigen
  Commit. Der Diff bleibt trotzdem klein genug für eine Review-Sitzung.
- `verifizierbar`: nein — reine Stil-/Konsistenzfrage, kein Gate prüft
  Commit-Granularität jenseits der Traceability-ID.
- `klasse`: „Commit-Granularität weicht von etablierter Trennung ab"

## Negativbefunde

- geprüft, ohne Befund: `tools/schema/nacharbeit-observability.sql` —
  `cdc_storage_bytes`-Zeile ist syntaktisch korrekter zusätzlicher
  `UNION ALL`-Zweig (`SELECT` ohne `FROM` ist in PostgreSQL für einen
  konstanten/funktionsbasierten Einzelwert gültig), Typkompatibilität mit
  den übrigen `UNION ALL`-Zweigen gegeben (`text`/`text`/`numeric`); der
  neue Kommentarblock beschreibt den Ist-Zustand mit auflösbaren
  Herkunfts-Ankern (`SPEC-009`, `LH-FA-RET-006`, `LH-FA-CAP-008`,
  `ADR-0046`), keine Chronik, kein Verweis auf verworfene Alternativen.
- geprüft, ohne Befund: Aggregationslogik — Plan-Nachzug Punkt 1 begründet
  nachvollziehbar, warum eine Einzeltabellen-Größe (`cdc.change`) statt
  einer Summe über mehrere `cdc`-Tabellen gewählt wurde (Wachstum
  proportional zum Erfassungsvolumen vs. Referenz-/Katalogdaten fester
  bzw. linear wachsender Größe); reine Beobachtung, kein Duplikat einer
  Domänenentscheidung (`ADR-0046` Kategorie C).
- geprüft, ohne Befund: `cdc_reader`-Rollenfläche real unverändert — `git diff
  bc7bba0..d4bfde9 -- tools/schema/nacharbeit-roles.sql` liefert keinen
  Treffer; das bestehende `GRANT SELECT ON cdc.metrics TO cdc_reader;`
  trägt die neue Zeile mit, da sie ein zusätzlicher `UNION ALL`-Zweig
  derselben View ist. Konsistent mit Architect-Verdikt Frage 2
  (`pg_relation_size()` PUBLIC-ausführbar, kein `SECURITY DEFINER`, kein
  direktes Privileg auf `cdc.change` nötig).
- geprüft, ohne Befund: `TestMetricsViewCarriesStorageBytes`
  (`internal/adapters/driven/postgresstorage/roles_test.go`) — vollständig
  und korrekt im gepushten Commit enthalten (Insert-Kette
  `source_table`/`schema_version`/`transaction`/`change`, Lesen von
  `cdc.metrics`, Assertion `> 0`); kein Rest-Schaden aus dem vom
  Implementer selbst gemeldeten `git checkout`-Zwischenfall erkennbar —
  der Testfall folgt demselben etablierten Muster wie
  `TestMetricsViewCarriesConsumerLag` unmittelbar davor in derselben Datei.
- geprüft, ohne Befund: `TestMVPMetricsCarriesStorageBytes`
  (`test/integration/integration_test.go`) taucht tatsächlich im
  `-run`-Filtermuster von `tools/harness/run-integration-tests.sh` (Zeile
  246) auf — die selbstgemeldete proaktive Aufnahme zur Vermeidung von
  `BEO-PGC/test-runner-stiller-ausschluss` ist real umgesetzt; das
  gesonderte zweite `-run`-Muster (`TestMVPSchemaChangeIncompatibleTypeChange`,
  Zeile 911) ist davon zu Recht nicht betroffen.
  `BEO-PGC/test-runner-stiller-ausschluss/state.md` bleibt unverändert bei
  1× (kein neues Auftreten), konsistent mit der Closure-Notiz.
- geprüft, ohne Befund: Commit-Traceability — `git log -1 --format=%s
  d4bfde9` trägt `LH-FA-RET-006`, keine `SPEC-*`/`ARC-*`-Struktur-ID im
  Betreff; `make commit-traceability` real ausgeführt (`HEAD~5..HEAD`):
  OK.
- geprüft, ohne Befund: ID-Link-Form (docs-check-Falle) — keine neuen
  Fließtext-Links mit Backticks im Link-Text in den geänderten
  `.md`-Dateien; `make docs-check` real ausgeführt: 362 Dateien, 0
  Befunde.
- geprüft, ohne Befund: Kommentar-Disziplin (§3.7,
  `BEO-PGC/slice-chronik-in-code-kommentar`, vor diesem Lauf
  3×/verkörpert) — vollständiger `grep`-Lauf gegen
  `slice-[0-9]{3}`/`welle-[0-9]+` über alle neuen/geänderten `.go`-/
  `.sql`-Zeilen dieses Diffs: kein Treffer. Alle neuen Kommentare
  (SQL-Kopf, Godoc in `roles_test.go`/`integration_test.go`) beschreiben
  Ist-Zustand bzw. tragen auflösbare Herkunfts-Anker, keine Chronik.
- geprüft, ohne Befund: DoD-Checkboxen (§2) — der Implementer-Bericht
  „bereits selbst nachgezogen" trifft real zu: alle Punkte außer „Review
  durchgeführt" und „drei Paarungen" sind in diesem Diff auf `[x]`
  gesetzt und mit konkreten Testfall-/Datei-Verweisen belegt (Gegensatz
  zur wiederholt aufgetretenen `BEO-PGC/dod-checkbox-nachzug`-Klasse —
  hier kein neues Auftreten).
- geprüft, ohne Befund: §6 Risiken — beide Risiken tragen einen Ausgang
  aus der geschlossenen Drei-Menge (beide „entfallen", mit Begründung);
  die Begründungen sind in sich stimmig (Bloat-Sichtbarkeit ist laut
  `LH-FA-RET-006`-Akzeptanzkriterium der gewollte Zustand, nicht ein
  Defekt; Architect-Verdikt real über
  `TestCdcReaderRoleReadsViewsNotBaseTables` bestätigt).
- geprüft, ohne Befund: Schicht-Grenzen (`AGENTS.md` §3.4) — keine
  Berührung, `spec/architecture.md` wird in diesem Diff nicht geändert.
- geprüft, ohne Befund: ADR-Immutabilität (§3.5) — keine ADR-Datei in
  diesem Diff geändert.
- geprüft, ohne Befund: Docker-only (§3.1) — keine lokale
  Toolchain-Installation in diesem Diff.
- geprüft, ohne Befund: §1 Ziel und Abgrenzung — drei Ausschlüsse mit
  Begründung, konsistent mit den vier Klassen (Bestand/anderer
  Vorgang/Schicht-Abgrenzung), keine erfundene Zusatzzahl.
- geprüft, ohne Befund: §8 Sub-Area-Prüfungen — einzige Sub-Area `*`/`PGC`,
  GF, konsistent mit `harness/conventions.md`-Modus-Deklaration; keine
  Änderung ggü. dem bereits vor diesem Lauf korrekten Zustand.

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Register-Ausgang an falschen
Prozessschritt delegiert (gestrichen vs. Lese-Schritt)" ·
„Commit-Granularität weicht von etablierter Trennung ab"

## Verdikt

**Merge-blockierend:** nein, aber F-1 (MEDIUM) sollte vor der finalen
Slice-Closure (`git mv` nach `done/`) behoben werden — keine Fixrunde am
Code, sondern eine Korrektur am Beobachtungs-Register selbst
(`state.md` von `BEO-PGC/retention-keine-loeschausfuehrung` auf den
Ausgang setzen, den die eigene Closure-Notiz bereits inhaltlich behauptet,
statt ihn an einen Prozessschritt zu delegieren, der diesen Eintrag
strukturell nie erreicht). Kein Rollen-Widerspruch, keine Konflikt-Sequenz
nach Modul 8 erforderlich — die Implementer-Aussage in der Closure-Notiz
widerspricht dem Regelwerk nicht inhaltlich, sondern in der
Prozess-Zuordnung; das ist eine Korrektur am Plan-Dokument, kein Konflikt
zwischen Rollen.

Alle repo-spezifisch besonders geprüften Punkte (Kommentar-Disziplin,
docs-check-Link-Falle, `cdc_reader`-Grant-Fläche, Testvollständigkeit nach
dem gemeldeten `git checkout`-Zwischenfall, Test-Runner-Filtermuster,
Commit-Traceability, DoD-Checkboxen, Aggregationslogik) sind ohne Befund.

**Übergabe:** F-1 geht an den Implementer zur Korrektur des
Register-Standes (keine Fixrunde am Code); F-2 ist ein reiner
Stil-Hinweis ohne erwartete Aktion. Dieser Report ist ein Lauf-Beleg und
wird über Läufe hinweg nicht wieder gelesen; die Summary-Zeile speist bei
Bedarf den Closure-Eintrag (Modul 5).
