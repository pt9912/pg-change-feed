# Review-Report: slice-014 — 2026-09-10

**Review-Art:** Code — geprüft gegen Slice-Plan + ADRs (Maintainability).

**Gegenstand:** Implementer-Commits von slice-014, `d663da1..HEAD` —
`e499e42` (Bootstrap-Logger, `slog.SetDefault`, `CDC_LOG_LEVEL`),
`7e35fa7` (Driven-Adapter `postgresstorage`/`postgresack` protokollieren
strukturiert), `4dcb493` (Replication-Stream-Adapter, `receive.go`),
`2b177f7` (Compose: `CDC_LOG_LEVEL`-ENV-Vertrag), `fe42546`
(Image-Hash-Beleg).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-014-strukturiertes-logging.md`
  (§1 Abgrenzung: Domain bleibt frei; §5 Closure-Trigger: stdout-Diff als
  Beleg der JSON-Struktur)
- `spec/lastenheft.md` §`LH-QA-OPS-004` (Strukturiertes Logging)
- `spec/architecture.md` §1 (`ARC-004` Outbound Ports, `ARC-006` Driven
  Adapters, `ARC-007` Bootstrap), §2 (Schichten-Tabelle,
  Import-Constraints), §3 (`ARC-011` Telemetrie-Backend — „Metriken und
  strukturierte Logs … über den Outbound Port substituierbar")
- [`ADR-0001`](../plan/adr/0001-hexagonale-architektur.md)
  (Hexagonale Architektur, Domain-Grenze), [`ADR-0024`](../plan/adr/0024-observability-ausserhalb-der-domain.md)
  (Observability außerhalb der Domain: `MetricsPort`/`EventSinkPort`
  durch Driven Adapters, Alternative B „globale Telemetrie-Singletons"
  explizit verworfen), [`ADR-0026`](../plan/adr/0026-composition-root.md)
  (Composition Root: Alternative B „globale Singletons" für die
  Verdrahtung von Driving/Inbound/Application/Outbound/Driven verworfen)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.2 Suppression-Verbot,
  §3.3 Move/Inhalt-Trennung, §3.7 Kommentar-Klassen) ·
  `harness/conventions.md` (MR-000/MR-001, genau eine Sub-Area `PGC`,
  Greenfield) · `.a-check.yml` (Layer-Edges,
  `composition_root: internal/bootstrap/**, cmd/**, test/integration/**`)
- `docs/plan/planning/observations/` (Register-Stand laut Slice-Plan §8:
  neun Einträge zum Planungszeitpunkt, keiner betrifft Logging)

---

## Findings

### F-1 — Globaler `slog`-Singleton für strukturiertes Logging umgeht die in ADR-0024 entschiedene Port/Driven-Adapter-Form

- `kategorie`: HIGH
- `quelle`: [`ADR-0024`](../plan/adr/0024-observability-ausserhalb-der-domain.md)
  (Entscheidung: „Logging-/Metrics-Frameworks bleiben Infrastruktur.
  MetricsPort und EventSinkPort werden durch Driven Adapters
  implementiert."; Alternative B „globale Telemetrie-Singletons" explizit
  verworfen — Contra: „versteckte Abhängigkeit; Verdrahtung und
  Testdoubles intransparent") · `spec/architecture.md` §3 `ARC-011`
  („Telemetrie-Backend … Metriken und strukturierte Logs … über den
  Outbound Port substituierbar")
- `pfad`: `internal/bootstrap/wiring.go:213-214` (`slog.SetDefault(newLogger(...))`,
  danach paketweite `slog.InfoContext`/`ErrorContext`-Aufrufe) ·
  `internal/adapters/driven/postgresack/ack.go:37,50` ·
  `internal/adapters/driven/postgresstorage/store.go:40,55` ·
  `internal/adapters/driven/postgresstorage/heartbeat.go:43,55` ·
  `internal/adapters/driven/postgresstorage/consumerstate.go:46,60` ·
  `internal/adapters/driven/postgresstorage/tableactivation.go:54,141,227,240` ·
  `internal/adapters/driving/replication/receive/receive.go:128,224,238,283,286`
- `befund`: `ADR-0024` wurde ausdrücklich für `LH-QA-OPS-003`/`LH-QA-OPS-004`
  geschrieben (Kontext-Abschnitt zitiert beide) und entscheidet für
  strukturiertes Logging **und** Metriken dieselbe Form: Weiterleitung
  über `MetricsPort`/`EventSinkPort`, implementiert durch Driven Adapters
  — nicht über einen globalen Telemetrie-Singleton, den dieselbe ADR als
  Alternative B mit genau der Begründung verwirft, die hier zutrifft
  (keine Testdouble-Substitution möglich, Abhängigkeit auf jeden
  Adapter-Aufrufer verteilt statt an einer Port-Grenze sichtbar).
  `spec/architecture.md` `ARC-011` bestätigt das für strukturierte Logs
  namentlich: Sie sollen „über den Outbound Port substituierbar" sein.
  Der Diff baut stattdessen `slog.SetDefault` in der Composition Root und
  lässt alle Driven-/Driving-Adapter über die paketweiten `slog`-Funktionen
  gegen diesen globalen Zustand schreiben — kein `EventSinkPort`, kein
  Driven Adapter dafür existiert im Repo (bestätigt: kein Treffer für
  `EventSinkPort`/`MetricsPort` in `internal/`). Die Commit-Begründung
  (`e499e42`: „ein Logger je Adapter-Konstruktor ist kein Bestandteil
  dieses Verdrahtungsstands") verteidigt die Wahl ausschließlich gegen
  `ADR-0026` (Composition Root) — die dort verworfenen „globalen
  Singletons" betreffen jedoch die *fachliche* Verdrahtung (Driving
  Adapters, Inbound Ports, Application Services, Outbound Ports, Driven
  Adapters selbst als Objekte im Abhängigkeitsgraph), nicht ein
  Cross-Cutting-Infrastrukturdetail wie einen Logger; `slog.SetDefault` +
  paketweite Aufrufe ist dafür das stdlib-Idiom und verletzt `ADR-0026`
  nicht. Die tatsächlich einschlägige Entscheidung ist `ADR-0024`, und die
  wird an keiner Stelle des Diffs oder der Commit-Nachrichten adressiert
  — die vom Implementer selbst gestellte Design-Frage zielt auf die
  falsche ADR.
- `verifizierbar`: nein — kein Gate modelliert eine Pflicht, Telemetrie
  über einen Outbound Port zu leiten; `.a-check.yml` prüft nur
  Layer-Kanten zwischen `domain`/`ports`/`app`/`adapters`, nicht die
  Nutzung von Cross-Cutting-Infrastruktur innerhalb der Adapter-Schicht
- `klasse`: Globaler Telemetrie-Singleton statt Port/Driven-Adapter (1.
  Auftreten)

### F-2 — Keine automatisierte Prüfung der Log-Feldstruktur

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `v6.5.0` ·
  `modul-08-agentenrollen.md` §Welche Rolle braucht welche
  Artefaktklasse — DoD-/Spec-Konformität ist Verifier-Gegenstand, nicht
  Reviewer-Gegenstand
- `pfad`: `internal/bootstrap/wiring_test.go` (nur `parseLogLevel`/Config
  getestet); kein Test in den geänderten Adapter-Dateien prüft
  Feldnamen/-Struktur der JSON-Log-Ausgabe
- `befund`: Der Implementer benennt die Lücke selbst; der Slice-Plan
  akzeptiert sie ausdrücklich als Closure-Beleg („die JSON-Log-Struktur
  ist am realen Integrationstest-Lauf beobachtbar (stdout-Diff)", §5).
  Gegen den Plan ist das damit kein Befund — ob der manuelle stdout-Diff
  als Closure-Beleg tatsächlich vorliegt und ausreicht, ist eine
  Verifier-Frage (DoD-Konformität), keine Reviewer-Frage.
- `verifizierbar`: ja — `make test`/Integrationstest-Log-Ausgabe
  (Verifier-Prüfung)
- `klasse`: Fehlender automatisierter Struktur-Test bei neuem Vertrag,
  im Plan vorab akzeptiert (1. Auftreten)

---

## Negativbefunde

- geprüft, ohne Befund: **`ADR-0001`/`ADR-0024`-Domain-Grenze** — kein
  `slog`-Aufruf in `internal/domain/` oder
  `internal/application/usecase/` (`grep -rn "slog\." internal/domain
  internal/application/usecase` liefert keinen Treffer); Logging bleibt
  auf `internal/adapters/**` und `internal/bootstrap` beschränkt, wie in
  §1 des Slice-Plans zugesagt.
- geprüft, ohne Befund: **sensible Daten in Log-Zeilen** — kein DSN,
  Passwort, Connection-String oder `conninfo`-Fragment in einem der
  neuen `slog`-Aufrufe (`grep -iE "dsn|password|secret|connstr|conninfo"`
  auf den Diff liefert keinen Treffer); geloggte Felder sind IDs,
  Offsets, Tabellennamen, Fehlerklassen — keine Rohwerte aus
  Verbindungsdaten.
- geprüft, ohne Befund: **Feldnamen-Konsistenz** — durchgängig
  `snake_case` (`consumer_id`, `transaction_id`, `table_id`, `start_lsn`,
  `error`, `class`, `publication`, `source`, `slot`), keine
  Mischformen zwischen den fünf Commits.
- geprüft, ohne Befund: **`a-check`-Konformität** — alle geänderten
  Dateien liegen in ihren deklarierten Layern
  (`adapters`/`composition_root`); keine neue Kante verletzt
  `.a-check.yml` (das Werkzeug modelliert allerdings keine
  Cross-Cutting-Telemetrie-Pflicht, siehe F-1).
- geprüft, ohne Befund: **Docker-only** — kein lokales Toolchain-Install
  in den fünf Commits.
- geprüft, ohne Befund: **Suppression-Verbot** — keine `nolint`/`noqa`/
  `SuppressMessage`-Marker im Diff `d663da1..HEAD`.
- geprüft, ohne Befund: **Kommentar-Klassen (`AGENTS.md` §3.7)** — die
  neuen/erweiterten Kommentare in `wiring.go`, `ack.go`, `store.go`,
  `heartbeat.go`, `consumerstate.go`, `receive.go` sind indikativ, nennen
  den geltenden Zustand und den Grund (Kopplungs-/Zusage-Klasse), kein
  Konjunktiv über eine verworfene Alternative, kein abgebrochener Satz.
- geprüft, ohne Befund: **Traceability der fünf Commits** — jede Message
  trägt mindestens eine `LH-*`-/`ADR-*`-Kennung
  (`LH-QA-OPS-004`/`ADR-0026`/`ADR-0044`), keine Struktur-ID
  (`SPEC-*`/`ARC-*`) in einem Betreff.
- geprüft, ohne Befund: **`tools/schema/plan.yaml`** — die im
  Arbeitsverzeichnis unverändert dirty gebliebene Datei (Ziel-Host
  `cdc-store-test-pg` → `cdc-test-postgres`) ist ein generierter
  `schema migrate --report`-Beleg (`Makefile:112`), stammt aus keinem der
  fünf slice-014-Commits und ist mit `tools/schema` als
  Rollout-Artefakt, nicht als Logging-Änderung, plausibel unabhängig von
  diesem Slice; kein übersehener Change dieses Diffs.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Globaler Telemetrie-Singleton statt
Port/Driven-Adapter (1. Auftreten) · Fehlender automatisierter
Struktur-Test bei neuem Vertrag, im Plan vorab akzeptiert (1. Auftreten)

## Verdikt

**Merge-blockierend:** ja, wegen F-1 (HIGH) — der Slice-Plan verlangt
als Closure-Trigger „Review-Schluss ohne offenes HIGH-Finding" (§5); der
`in-progress → done`-Übergang bleibt gesperrt, bis F-1 einen Ausgang
trägt. F-1 ist ein HIGH mit Rollen-Widerspruch: Die Commit-Begründung
verteidigt die gewählte Form explizit (gegen `ADR-0026`), der Reviewer
stellt eine Verletzung einer anderen, direkt einschlägigen ADR
(`ADR-0024`) fest. Das greift den Konflikt-Pfad aus Modul 8 §Konflikt-Pfad
als Rollen-Sequenz: Die Klärung braucht ein Architect-Verdikt als
Übergabe-Artefakt — eines der drei legitimen Verdikte (Plan/ADR-Bezug im
Slice war unvollständig und wird korrigiert · `ADR-0024` wird per
Folge-ADR gelockert/`supersedes`d, mit begründetem Verglichene-Alternativen-
Eintrag für Cross-Cutting-Infra-Logging · die vorliegende Form ist eine
legitime, aber bisher undokumentierte Lockerung und wird per Folge-ADR
nachgezogen) — nicht eine Herabstufung, weil der Implementer
widerspricht. F-2 ist kein Closure-Stopp; sie geht als Hinweis an den
Verifier (Modul 8 §Welche Rolle braucht welche Artefaktklasse).

**Zur eigentlichen Design-Frage (`ADR-0026` vs. globaler Singleton):**
unproblematisch. `ADR-0026` verwirft „globale Singletons" für die
*fachliche* Verdrahtung (Driving Adapters, Inbound Ports, Application
Services, Outbound Ports, Driven Adapters als Objekte im
Abhängigkeitsgraph) — mit der Begründung „versteckte Abhängigkeiten,
unklare Lebenszyklen, Tests brauchen globalen Zustand". Ein
Cross-Cutting-Infrastrukturdetail wie ein Logger, per
`slog.SetDefault` einmal in der Composition Root gesetzt und über die
stdlib-paketweiten Funktionen genutzt, ist kategorial etwas anderes: Es
trägt keine fachliche Lebenszyklus- oder Identitäts-Semantik und ist das
in `log/slog` selbst vorgesehene Idiom für genau diesen Zweck. `ADR-0026`
ist durch diesen Diff **nicht** verletzt. Das eigentliche Problem liegt
in `ADR-0024`, die der Implementer in seiner Selbstprüfung nicht
herangezogen hat (F-1).

**Übergabe:** F-1 geht an Implementer/Architect (Rollen-Widerspruch,
Konflikt-Pfad Modul 8) — kein Self-Review-Ausgang. F-2 geht als Hinweis
an den Verifier. Die **Finding-Klassen** gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist
ein **Lauf-Beleg** und wird über Läufe hinweg nicht wieder gelesen. Der
Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).

---

**Gate-Beleg:** `make gates` nach diesem Report-Commit (Lauf
2026-09-10); Ergebnis im Commit-Text vermerkt.
