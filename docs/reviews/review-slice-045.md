# Review-Report: slice-045 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-045`, §1/§2/§3/§4/§6/§8),
`welle-13`, den Architect-Verdikt
`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`, `ADR-0046`
sowie `AGENTS.md` §3 Hard Rules (§3.1, §3.5, §3.7) — Rollentrennung Modul 8:
diese Prüfung läuft gegen Plan/ADR/Hard Rules (Maintainability), nicht gegen
DoD (Verifier-Aufgabe).

**Gegenstand:** Commit `134bc0a` (`feat(retention): cdc.retention_blockers
zeigt blockierende Consumer (LH-FA-RET-005)`) gegen Elter `5d36fea`; Diff:
`tools/schema/schema.yaml`, `tools/schema/nacharbeit-roles.sql`,
`tools/schema/plan.yaml`, `tools/schema/down.sql` (d-migrate-generiert),
`internal/adapters/driven/postgresstorage/roles_test.go`,
`test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh`,
`docs/user/benutzerhandbuch.md`,
`docs/plan/planning/in-progress/slice-045-blockierende-consumer-sichtbarkeit.md`
(Plan-Nachzug, DoD-Häkchen).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft
2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-045-blockierende-consumer-sichtbarkeit.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug fünf
  Punkte, §4 Trigger, §6 Risiken, §8)
- `docs/reviews/architect-verdict-retention-loeschausfuehrung.md`
  (vollständig — Frage 1/Frage 2, View-Owner-Muster, `ADR-0046`-Einordnung)
- `docs/plan/planning/done/slice-043-changestoreport-loeschmethode.md`,
  `docs/plan/planning/done/slice-044-retention-hintergrundjob.md`
  (Grundlage, `RetentionPolicy.AllowsDeletion`, `ConsumerStatePort`-Lesart)
- `internal/domain/model/retention.go` (`AllowsDeletion`), `position.go`
  (`SourcePosition.Before`/`Compare`, `Offset uint64`) — reale Prüfung, ob
  die View-Logik (`DISTINCT ON` nach kleinster `acknowledged_position`) der
  Domain-Freigabebedingung entspricht, statt sie zu duplizieren
- `spec/lastenheft.md` `LH-FA-RET-005` (Happy Path/Boundary/Negative)
- `tools/schema/schema.yaml` bestehende Views (`consumer_status`,
  `active_tables`, `changes`) als Vergleichsmuster (Subquery-Form,
  `source_dialect`/`columns:`-Pflicht)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (ADR-Immutabilität), §3.7
  (Kommentar-Disziplin — `BEO-PGC/slice-chronik-in-code-kommentar`, 3× vor
  diesem Lauf, besonders geprüft)
- `git log -1 --format='%s' 134bc0a` (Commit-Traceability: `LH-*`/`ADR-*`
  vorhanden, keine `SPEC-*`/`ARC-*`-Struktur-ID im Betreff)
- `make docs-check` (lokal ausgeführt: 357 Dateien, 0 Befunde)
- `grep -nE '^\+.*\b(slice-[0-9]{3}|welle-[0-9]+)\b'` über den vollständigen
  `.go`/`.sql`/`.yaml`-Diff (vollständiger Durchlauf, kein Treffer)

---

## Findings

Keine HIGH-/MEDIUM-/LOW-Findings in diesem Lauf.

### F-1 — LH-FA-RET-005-Boundary hängt an zwei Views, im DoD nur eine genannt

- `kategorie`: INFO
- `quelle`: `LH-FA-RET-005` (Boundary: „mehrere blockierende Consumer …
  alle einzeln erkennbar")
- `pfad`: `docs/plan/planning/in-progress/slice-045-blockierende-consumer-sichtbarkeit.md:79-87`
  (DoD-Punkt 2) vs. Plan-Nachzug Punkt 1 (Zeilen 140–151)
- `befund`: `cdc.retention_blockers` liefert per `DISTINCT ON
  (cp.source_id)` genau **eine** Zeile je Quelle (den am weitesten
  zurückliegenden Consumer). Die Boundary „mehrere blockierende Consumer …
  alle einzeln erkennbar" wird laut Plan-Nachzug Punkt 1 nicht von dieser
  View, sondern von der bereits bestehenden `cdc.consumer_status`
  abgedeckt (Architect-Verdikt, Frage 1, bereits vor diesem Slice
  festgestellt). Das DoD-Häkchen §2 Punkt 2 formuliert „`LH-FA-RET-005`
  real erfüllt" jedoch ausschließlich mit Verweis auf den neuen
  Integrationstest gegen `cdc.retention_blockers` — die Mitwirkung von
  `consumer_status` an der Boundary-Erfüllung ist nur im Plan-Nachzug,
  nicht im DoD-Punkt selbst sichtbar. Das ist kein Plan-Defekt (die
  Argumentation im Plan-Nachzug ist in sich stimmig und deckt sich mit dem
  Architect-Verdikt), sondern ein Hinweis für den nächsten Rollenwechsel.
- `verifizierbar`: ja — `TestMVPRetentionBlockersViewShowsFurthestBehindConsumer`
  deckt nur den Happy-Path/den Fall „ein Consumer je Quelle sichtbar"; ein
  eigener Boundary-Test mit **mehreren** gleichzeitig zurückliegenden
  Consumern derselben Quelle existiert in diesem Diff nicht (erwartungsgemäß,
  da laut Plan `consumer_status` diesen Fall trägt — dort ebenfalls ohne
  neuen Testfall in diesem Diff, weil `consumer_status` unverändert bleibt).
- `klasse`: „Akzeptanzkriterium über zwei Artefakte verteilt, DoD nennt nur
  eines" — Hinweis für Verifier-Rolle (Modul 8, Reviewer prüft nicht gegen
  DoD/Spec-Konformität).

## Geprüft, ohne (weiteren) Befund

- **View-Design (`cdc.retention_blockers`):** reine Projektion, keine
  Duplikat-Logik zu `RetentionPolicy.AllowsDeletion`. Real geprüft: Die
  Freigabebedingung verlangt, dass **alle** übergebenen Consumer-Positionen
  an oder hinter der Change-Position liegen — die bindende Grenze ist damit
  mathematisch die **kleinste** bestätigte Position über alle Consumer
  einer Quelle. Genau das liefert `DISTINCT ON (cp.source_id) ORDER BY
  cp.source_id, cp.acknowledged_position ASC, c.consumer_id ASC` (Tie-Break
  `consumer_id`). `SourcePosition.Offset` ist `uint64`, `Before`/`Compare`
  vergleichen numerisch auf `Offset` — die SQL-`bigint`-Ordnung auf
  `acknowledged_position` bildet dieselbe Ordnung ab, keine abweichende
  Vergleichssemantik.
- **`backlog`-Subquery:** identisches Muster wie das bestehende
  `consumer_status.latest_commit_position` (korrelierte Unterabfrage gegen
  `cdc.transaction`, `max(commit_position)` je `source_id`) — keine neue
  Abfrageform, wie im Plan-Nachzug behauptet.
- **„Consumer nie bestätigt"-Randfall (§6 Risiko 1):** `INNER JOIN`
  statt `LEFT JOIN` auf `cdc.consumer_position` — ein nie bestätigender
  Consumer trägt konsequent keine Zeile, dieselbe Abwesenheits-Lesart wie
  `ConsumerStatePort.Positions`. Kein Sonderfall, der fälschlich „kein
  Blocker" oder einen irreführenden Wert zeigen könnte; Boundary real
  getestet (der bereits weiter bestätigende Consumer erscheint nicht als
  Zeile).
- **`cdc_reader`-Grant-Erweiterung (`nacharbeit-roles.sql`):** konsistent
  mit dem Architect-Verdikt (Frage 2, `ADR-0046`): `GRANT SELECT` auf die
  neue View selbst, kein Grant auf eine Basistabelle
  (`cdc.consumer_position`/`cdc.consumer`/`cdc.transaction` bleiben ohne
  direkten `cdc_reader`-Zugriff) — dasselbe View-Owner-Muster (Definer-
  Semantik), das `cdc.metrics`/`cdc.heartbeat`/die drei bestehenden Views
  bereits tragen. Keine neue Fähigkeit, keine Rollen-*Erweiterung* im
  Sinne einer neuen Zugriffs-Kategorie — real regressionsgetestet
  (`TestCdcReaderRoleReadsViewsNotBaseTables`, neue Zeile gegen
  `cdc.retention_blockers`, Erfolg erwartet; die bestehenden
  Basistabellen-Verweigerungen bleiben unverändert in derselben
  Testfunktion).
- **Kommentar-Disziplin (§3.7, `BEO-PGC/slice-chronik-in-code-kommentar`,
  Zähler vor diesem Lauf bei 3×/verkörpert):** vollständiger `grep`-Lauf
  gegen `slice-[0-9]{3}`/`welle-[0-9]+` über alle neuen/geänderten
  `.go`-/`.sql`-/`.yaml`-Zeilen dieses Diffs — kein Treffer. Alle neuen
  Kommentare (Schema-Kopf, View-`description`, `nacharbeit-roles.sql`,
  Godoc in `integration_test.go`/`roles_test.go`) beschreiben den
  Ist-Zustand bzw. tragen Herkunfts-Anker (`LH-FA-RET-005`,
  `RunRetentionUseCase`, `RetentionPolicy.AllowsDeletion` als
  Kopplungs-Referenz) — keine Chronik, kein Verweis auf eine verworfene
  Alternative oder abwesenden Text.
- **Consumer-Cleanup-Fix im Testfall (Plan-Nachzug Punkt 5):** `t.Cleanup`
  ruft `ConsumerStatePort.Remove` für beide im Testfall registrierten
  Consumer nach jedem Lauf auf — isoliert von den über
  `register-consumer`/`acknowledge-consumer` geführten Consumern des
  externen Black-Box-Rundlaufs (unterschiedliche Consumer-IDs,
  unterschiedlicher Registrierungsweg: direkter `ConsumerStatePort`-Aufruf
  vs. `docker exec`). Der Fix greift genau an der vom Implementer
  benannten Ursache: `Positions(mvpSource)` liest alle bestätigten
  Positionen der Quelle unabhängig vom anlegenden Testfall — ohne Cleanup
  hätten liegen gebliebene Positionen den bestehenden
  `slice-044`-Retention-Beleg dauerhaft blockiert. Keine erkennbaren
  Seiteneffekte auf andere Testfälle: Die entfernten Consumer-IDs
  (`retention-view-behind`/`retention-view-ahead`) sind ausschließlich in
  diesem neuen Testfall registriert, keine andere Testfunktion referenziert
  sie.
- **Schicht-Grenzen (`AGENTS.md` §3.4, `spec/architecture.md`):** keine
  Berührung — `spec/architecture.md` wird in diesem Diff nicht geändert,
  `ARC-005` (SQL-Funktionen/Views) bleibt als bestehender Anker unverändert
  zitiert.
- **ADR-Immutabilität (§3.5):** keine ADR-Datei in diesem Diff geändert;
  der Architect-Verdikt bestätigt bestehende `Accepted`/`permanent`-ADRs
  (`ADR-0009`, `ADR-0011`, `ADR-0012`, `ADR-0014`, `ADR-0029`, `ADR-0046`),
  keine davon wird überschrieben oder widersprochen.
- **Docker-only (§3.1):** keine lokale Toolchain-Installation in diesem
  Diff.
- **Traceability:** Commit-Betreff trägt `LH-FA-RET-005`, keine
  `SPEC-*`/`ARC-*`-Struktur-ID im Betreff.
- **ID-Link-Form (docs-check-Falle):** keine neuen Fließtext-Links mit
  Backticks im Link-Text in diesem Diff (`grep` gegen
  `` \[`(LH|ADR|SPEC|ARC|slice|BEO|CO|MR|RC)- `` auf hinzugefügten Zeilen:
  kein Treffer); `make docs-check` real ausgeführt: 357 Dateien, 0 Befunde.
- **Zwei-Quellen-Drift:** `tools/schema/schema.yaml`-Kopfkommentar (Zeile
  32, „Die drei Views …") bleibt unverändert stehen und beschreibt weiter
  korrekt nur die drei `LH-FA-SST-002`-Views; die neue View bekommt einen
  eigenen, klar abgegrenzten Absatz („Eine vierte View … trägt einen
  anderen Vertrag") — keine widersprüchliche Doppelaussage.
- **§1 Ziel und Abgrenzung / CLI-Erweiterung bewusst unterlassen:**
  Ausschluss trägt eine Begründung (Slice-Größe, kein Informationsgewinn
  ggü. direktem SQL-Zugriff über `CDC_READER_DSN`) und ordnet sich sauber
  der Klasse „anderer Vorgang/Bestand bleibt bewusst stehen" zu — konsistent
  mit dem im Slice-Ziel selbst als „Implementer-Entscheidung" vorgesehenen
  Spielraum.
- **§8 Sub-Area-Prüfungen:** einzige Sub-Area `*`/`PGC`, GF, konsistent mit
  `harness/conventions.md` Modus-Deklaration. Beobachtungs-Register-Treffer
  (`BEO-PGC/retention-keine-loeschausfuehrung` 0×,
  `BEO-PGC/d-migrate-nacharbeit` 5×) real gegen die Evidence-Verzeichnisse
  geprüft (`ls .../evidence/`) — Zählerstände stimmen mit der Plan-Aussage
  überein.

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Akzeptanzkriterium über zwei Artefakte
verteilt, DoD nennt nur eines" (1×, INFO — kein Steering-Loop-Zähler-Eintrag
erforderlich, reiner Hinweis an die Verifier-Rolle).

## Verdikt

**Merge-blockierend:** nein. Kein HIGH-, kein MEDIUM-Finding. Die
`cdc_reader`-Grant-Erweiterung ist konsistent mit dem Architect-Verdikt
(View-Owner-Muster, keine neue Zugriffs-Kategorie); das View-Design
dupliziert `RetentionPolicy.AllowsDeletion` nicht; die Kommentar-Disziplin
(§3.7, dreimal zuvor getroffene Fehlerklasse) ist in diesem Diff vollständig
eingehalten (kein Treffer im vollständigen `grep`-Lauf); der real gemeldete
Test-Cleanup-Fix ist sauber isoliert und ohne erkennbare Seiteneffekte.

Kein Rollen-Widerspruch, keine Konflikt-Sequenz nach Modul 8 erforderlich.

**Übergabe:** keine Fixrunde nötig. Das einzige Finding (F-1, INFO) geht
als Hinweis an die Verifier-Rolle — ob `LH-FA-RET-005`s Boundary-Kriterium
über die Kombination `cdc.retention_blockers` + bestehende
`cdc.consumer_status` tatsächlich vollständig gedeckt ist, ist eine
DoD-/Spec-Konformitätsfrage (Modul 8, außerhalb des Reviewer-Kontexts).
Dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht wieder
gelesen; die Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5).
