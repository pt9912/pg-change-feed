# Welle 8 — Black-Box-E2E und Integrationstest-Nachzug — Closure-Notiz

**Welle:** welle-8
**Abschluss:** 2026-09-12
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (`slice-027`): der erste
  echte Black-Box-E2E-Test des Repos — `register-consumer`/
  `acknowledge-consumer` laufen ausschließlich als externer `docker exec`-
  Prozess gegen den laufenden Feed-Container, inklusive echtem Container-
  Neustart und einem Beleg, dass keine bereits bestätigte Änderung erneut
  gelesen wird.
- [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…`003` (`slice-028`,
  **teilweise**): der Compose-Integrationstest belegt jetzt real
  PostgreSQLs serverseitige Durchsetzung der Rollentrennung
  (`cdc_reader`/`cdc_capture` mit/ohne `REPLICATION`-Attribut) — nicht
  aber, dass die tatsächlichen Replication-Stream-/ACK-Adapter selbst
  rollenbeschränkt sind (Adapter-Ebene bleibt offen,
  `BEO-PGC/rollen-test-abdeckungsluecken` jetzt 2×).
- `test/integration/mvp_test.go` → `integration_test.go` umbenannt,
  Makefile-Helptext und `harness/README.md` nachgezogen — der Name trägt
  jetzt den tatsächlichen, seit `welle-3` gewachsenen Scope.
- Ein neuer, unabhängig vom Testauftrag entstandener Fund: `Retention`
  (`LH-FA-RET-002`…`006`) hat nur Entscheidungslogik, keine
  Löschausführung — registriert als
  `BEO-PGC/retention-keine-loeschausfuehrung`, adressiert in einer eigenen,
  bereits vorgemerkten Feature-Welle (Roadmap *Nächste Wellen*).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Reviewer und Verifier haben bei `slice-028` eine vom Implementer zu stark
  formulierte DoD-Zeile gefunden (F-1) und dieselbe Einordnung einer
  vorbestehenden Superuser-Verdrahtung (F-2, unverändert seit `slice-023`)
  unabhängig voneinander getroffen — der Reviewer eskalierte den Fund
  bewusst **nicht** zu HIGH, nachdem er selbst nachgeprüft hatte, dass
  keine der HIGH-Kategorien zutrifft und keine bestehende Zusage
  widersprochen wird. Der Verifier bestätigte anschließend unabhängig, dass
  die vom Planner vorgenommene DoD-Korrektur akkurat war.
- Der Implementer verifizierte seine eigene Plan-Annahme (`slice-028`)
  empirisch an einem Scratch-Postgres, bevor er sie übernahm — ersparte
  eine unnötige Erweiterung von `run-replication-tests.sh`.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- `slice-028`s ursprünglicher DoD-Text behauptete, `BEO-PGC/rollen-test-
  abdeckungsluecken` Punkt (2) vollständig zu schließen — real erreicht
  wurde nur die PostgreSQL-Server-Ebene, nicht die Adapter-Ebene.
  Konsequenz: DoD/§1/§6 korrigiert, die Beobachtung bleibt mit 2× im
  Register offen, statt sie stillschweigend als erledigt zu markieren.
- Während der Closure-Vorbereitung entstand aus einem Nutzerwunsch
  ("E2E-Tests für alle implementierten Fähigkeiten") ein deutlich größeres
  Folgeprogramm: eine repo-weite Fähigkeiten-Bestandsaufnahme (75
  Anforderungen, ~55–58 vollständig implementiert) plus fünf neu
  vorgemerkte Wellen (E2E-Abdeckung CDC-Kernpfad, E2E-Abdeckung Verwaltung/
  Observability, Retention-Löschausführung, E2E-Abdeckung Retention,
  Performance-Benchmarks & Test-Coverage-Gate) — vor `walsender-
  wirksamkeit`/`LH-FA-SST-007` eingereiht. Das ist kein Scope dieser
  Welle, aber ihre unmittelbare Fortsetzung.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle 3× — der Normalfall.
`BEO-PGC/rollen-test-abdeckungsluecken` steht jetzt bei 2× (weiterhin unter
der Schwelle, zur Sichtung bei der nächsten Slice-Planung).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Mit
neuem Stand in dieser Welle: `rollen-test-abdeckungsluecken` (1× → 2×,
weiter offen — `slice-028` schloss nur die PostgreSQL-Ebene). Neu angelegt
in dieser Welle: `retention-keine-loeschausfuehrung` (0×, weiter offen —
gefunden bei einer Fork-Recherche, kein abgeschlossener Vorgang trägt bisher
einen Beleg). Unter der Schwelle, unverändert: `adapter-fehler-ausgang`
(1×), `dod-checkbox-nachzug-architect-pfad` (1×), `lese-doppelquelle` (2×),
`schema-rollout-fremdobjekte` (1×), `test-isolation-geteilter-zustand`
(1×), `walsender-wirksamkeit` (1×). Bereits verkörpert, unverändert:
`a-check-null-abdeckung` (3×), `d-migrate-nacharbeit` (4×),
`dod-checkbox-nachzug` (3×), `plan-nachzug` (2×), `plan-vorlagen-defekt`
(3×). Bereits eingetreten, unverändert: `cdc-capture-lag-real` (2×),
`health-endpoint-heartbeat` (2×), `rollen-verdrahtung` (4×),
`spec008-replication-luecke` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — die fünf als nächstes vorgemerkten Wellen (E2E-Abdeckungsprogramm,
Retention-Löschausführung, Performance-Benchmarks & Test-Coverage-Gate)
sind bislang nur Vorschau-Zeilen in der Roadmap (*Nächste Wellen*), noch
nicht als Slices geschnitten.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Beide Slices (`slice-027`, `slice-028`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Black-Box-E2E-Beleg (repo-weiter Verifikations-Beleg, Modul 8
  Closure-Schritt 1): neuer Abschnitt „Black-Box-CLI-Rundlauf" in
  `tools/harness/run-integration-tests.sh`, dreimal in Folge grün real
  gegen den Compose-Stack — Rundlauf ausschließlich über `docker exec`
  gegen die CLI, inklusive echtem Container-Neustart und Beleg, dass keine
  bereits bestätigte Änderung erneut gelesen wird.
- Rollen-DSN-Beleg (`slice-028`): PostgreSQL-Server-Ebene real geprüft
  (Schreib-Ablehnung `cdc_reader`, `REPLICATION`-Attribut-Durchsetzung),
  Adapter-Ebene **nicht** erreicht — ehrlich als `BEO-PGC/rollen-test-
  abdeckungsluecken` (2×) im Register festgehalten, nicht stillschweigend
  als vollständig behauptet.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig". Kein Carveout in dieser Welle. Kein
  bootstrap-aware Gate berührt. `ADR-0047`s Re-Evaluierungs-Trigger (vierte
  CDC-Rolle · Secrets-Management-System) ist nicht eingetreten — bestätigt.
  `ADR-0030` ist permanent, kein Trigger fällig.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — kein
  Steering-Loop-Eintrag mit `liegt in` in dieser Welle, nichts zu prüfen.
  Folge-Slice — keiner genannt, nichts zu prüfen. Register — jede zitierte
  `BEO-PGC/*`-Kennung existiert als Verzeichnis; alle Verzeichnisse tragen
  ein nicht leeres `evidence/`, **mit einer benannten, transparenten
  Ausnahme:** `BEO-PGC/retention-keine-loeschausfuehrung` (neu in dieser
  Welle) hat 0 Belege — sein eigenes `observation.md` benennt das explizit
  unter „Benannt, nicht gezählt" (Fund aus einer Fork-Recherche, kein
  abgeschlossener Vorgang trägt bisher einen Beleg). Kein fabrizierter
  Eintrag, keine stille Lücke.

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die beiden Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
Vor der ersten tatsächlichen Archivierung gilt die Prüfpflicht aus dem
Closure-Command: Geltungsbereich der Sensoren (`structure`-Regel 5 keilt
auf `done/slice-*.md` — bei Stubs in einem Unterverzeichnis anzupassen) und
Link-/ID-Pflichten im Stub prüfen.
