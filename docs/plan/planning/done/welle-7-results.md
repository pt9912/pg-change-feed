# Welle 7 — Replication-Schwellen-Überwachung — Closure-Notiz

**Welle:** welle-7
**Abschluss:** 2026-09-12
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
  (`slice-024`) trennt die `replication`-Fehlerklasse in Stream-Ordnungs-
  Verletzung (bleibt hart abbrechend) und Transport-/Verbindungsstörung
  (Schwellen-Kandidat) und setzt WAL-Rückstand-Schwellen (Warn 100 MiB,
  Fehler 1 GiB) — präzisiert `SPEC-008`/`SPEC-013` in
  `spec/pflichtenheft.md` direkt per ADR, ohne Lastenheft-CR.
- Die Metrik `cdc_wal_retention_bytes` (`slice-025`) misst periodisch
  gegen `pg_replication_slots`, bewusst außerhalb der `cdc.metrics`-View
  (Least-Privilege-Grenze von `cdc_reader` bleibt unberührt), mit
  eigenem Reconnect-Mechanismus bei Verbindungsausfall.
- Der Capture-Prozess (`slice-026`) reagiert jetzt auf den Schwellen-
  Vergleich: unterhalb der Warnschwelle unveränderte Fortsetzung, zwischen
  Warn-/Fehlerschwelle sichtbare Fortsetzung (Log-Warnung), oberhalb der
  Fehlerschwelle kontrollierter Abbruch über den bestehenden
  `replication`-Klassifikationspfad — genau `SPEC-008`s „Überwachung über
  Schwellen; kontrollierte Fortsetzung". Eine strukturelle
  Prioritäts-Garantie (`mergeStreamAndWALFaultOutcome`) stellt sicher,
  dass dieser neue Fortsetzungspfad eine echte Stream-Ordnungs-Verletzung
  niemals maskieren kann.
- `TestWALRetentionThresholdEndToEnd` belegt real, mit künstlich erzeugtem
  WAL-Rückstand, beide Seiten der Schwelle in einem Lauf — das war das
  *Mehr* dieser Welle gegenüber den einzelnen Slice-DoDs.
- `BEO-PGC/spec008-replication-luecke` ist geschlossen (Ausgang
  *eingetreten*): die seit `slice-020` benannte Spec-vs-Code-Lücke ist
  real gebaut, nicht mehr nur beschrieben.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die Rot-Grün-Disziplin trug über die ganze Rollenkette hinweg, am
  deutlichsten an der sicherheitskritischen Prioritäts-Garantie in
  `slice-026`: Implementer, Reviewer und Verifier haben sie jeweils
  unabhängig mit eigener Rot-Grün-Gegenprobe bewiesen, statt eine
  Behauptung weiterzureichen.
- Der Architect-Lauf in `slice-024` verifizierte die Sentinel-Trennung
  selbst im Code, statt den Slice-Plan-Vorschlag blind zu übernehmen —
  Reviewer und Verifier bestätigten das unabhängig.
- Die Least-Privilege-Grenze aus `tools/schema/nacharbeit-observability.sql`
  wurde bewusst eingehalten: `cdc_wal_retention_bytes` läuft prozessintern,
  nicht über `cdc.metrics`/`cdc_reader` — dieselbe Begründung, mit der der
  Heartbeat-Mechanismus dort bereits einer reinen Lese-View ausweicht.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Der Expositionsweg der neuen Metrik war im Slice-Plan noch offen
  (Tabellenspalte analog `cdc.process_heartbeat` vs. strukturiertes Log) —
  der Implementer entschied sich für Log, mit Begründung im Plan-Nachzug
  (`slice-025` liest den Wert ohnehin im selben Prozess für die
  Schwellen-Entscheidung, kein zusätzliches Schema/keine zusätzliche
  Schreibrolle nötig).
- Beide Code-Slices (`slice-025`, `slice-026`) bekamen im Review je eine
  MEDIUM-Lücke (fehlender Reconnect bei dauerhaftem Verbindungsausfall;
  fehlender Test für den Zero-Value-Schwellenwert-Fallback) — beide über
  eine kurze, gezielte Fixrunde mit eigener Rot-Grün-Verifikation behoben,
  kein Merge-Blocker.
- Im Gespräch mit dem Nutzer während der Closure-Vorbereitung von
  `slice-026` fiel eine Terminologie-Lücke auf: Der als „Ende-zu-Ende-Test"
  bezeichnete `TestWALRetentionThresholdEndToEnd` (und ebenso die beiden
  älteren `*_endtoend_test.go`-Dateien aus `welle-6`/`slice-022`) ruft die
  Produktionsfunktion in-process auf — nach `ADR-0030`s Testpyramide ist
  das ein **Integrationstest**, kein Black-Box-**E2E**-Test; `make
  test-integration` ist zudem seit `welle-3` inhaltlich nicht über
  MVP-Scope hinausgewachsen. Konsequenz: eine eigene Testing-Welle
  (Black-Box-E2E über die bestehende CLI, Integrationstest-Nachzug) ist in
  der Roadmap vorgemerkt — direkt nach dieser Welle, vor den beiden bereits
  zuvor vorgemerkten Wellen (Nutzerentscheidung).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle 3× — der Normalfall.
`BEO-PGC/spec008-replication-luecke` schloss stattdessen direkt (Ausgang
*eingetreten* bei 2×, kein Schwellen-Mechanismus nötig — siehe
Beobachtungs-Register unten).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Mit
neuem Ausgang in dieser Welle: `spec008-replication-luecke` (2×,
**eingetreten** → `slice-026`, siehe „Was wurde geliefert?" oben). Neu
angelegt in dieser Welle: `dod-checkbox-nachzug-architect-pfad` (1×,
weiter offen — bei `slice-024` gefunden, zur Sichtung bei der nächsten
Slice-Planung). Unter der Schwelle, unverändert: `adapter-fehler-ausgang`
(1×), `lese-doppelquelle` (2×), `rollen-test-abdeckungsluecken` (1×),
`schema-rollout-fremdobjekte` (1×), `test-isolation-geteilter-zustand`
(1×), `walsender-wirksamkeit` (1×). Bereits verkörpert, unverändert:
`a-check-null-abdeckung` (3×), `d-migrate-nacharbeit` (4×),
`dod-checkbox-nachzug` (3×), `plan-nachzug` (2×), `plan-vorlagen-defekt`
(3×). Bereits eingetreten, unverändert: `cdc-capture-lag-real` (2×),
`health-endpoint-heartbeat` (2×), `rollen-verdrahtung` (4×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — die als nächstes vorgemerkte Testing-Welle (Black-Box-E2E,
Integrationstest-Nachzug) ist bislang nur eine Vorschau-Zeile in der
Roadmap (*Nächste Wellen*), noch nicht als Slice geschnitten.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices (`slice-024`, `slice-025`, `slice-026`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Ende-zu-Ende-Beleg (repo-weiter Verifikations-Beleg, Modul 8
  Closure-Schritt 1): `internal/bootstrap/walretention_endtoend_test.go`
  (`TestWALRetentionThresholdEndToEnd`), real gegen PostgreSQL, dreimal in
  Folge grün — künstlich erzeugter, wachsender WAL-Rückstand durchläuft
  beide Seiten der Schwelle in einem Lauf (Warn: sichtbare Fortsetzung;
  Fehler: kontrollierter Abbruch mit `outbound.ErrReplication`), ein
  Regressionstest bestätigt zugleich, dass eine Stream-Ordnungs-Verletzung
  davon strukturell unberührt bleibt.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig". Kein Carveout in dieser Welle. Kein
  bootstrap-aware Gate berührt. `ADR-0049`s eigener Re-Evaluierungs-Trigger
  (reale Produktionsdaten oder ein Ende-zu-Ende-Test zeigen falsche
  Schwellenwerte) ist **nicht** eingetreten: `slice-026`s Test lief mit
  skalierten Testfixture-Werten (32 KiB/512 KiB), nicht den realen
  Produktions-Startwerten (100 MiB/1 GiB) — die Frage, ob Letztere in der
  Praxis passen, bleibt am ADR hängen, nicht an dieser Welle.

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die drei Slice-Dateien, ihre Review-/
Verifier-Reports und die ADR sowie dieser Welle-Plan bleiben vollständig
in `done/`. Vor der ersten tatsächlichen Archivierung gilt die Prüfpflicht
aus dem Closure-Command: Geltungsbereich der Sensoren (`structure`-Regel 5
keilt auf `done/slice-*.md` — bei Stubs in einem Unterverzeichnis
anzupassen) und Link-/ID-Pflichten im Stub prüfen.
