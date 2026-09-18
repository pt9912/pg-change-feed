# Welle 13 — Retention-Löschausführung — Closure-Notiz

**Welle:** welle-13
**Abschluss:** 2026-09-13
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- `LH-FA-RET-002`…`004` (Aufbewahrung, zeit- und consumer-basiert): eine
  reale Löschausführung existiert jetzt — `ChangeStorePort.DeleteChanges`
  (atomar, inkl. Bereinigung verwaister `cdc.transaction`-Zeilen) und
  `RunRetentionUseCase`/`RunRetentionService` rufen
  `RetentionPolicy.AllowsDeletion` (`ADR-0014`) real auf (`slice-043`).
- Ein Hintergrundzug `runRetentionCleanup`
  (`internal/bootstrap/wiring.go`) macht die Löschausführung im laufenden
  Prozess wirksam, mit `SystemClockAdapter` als erster
  Produktionsimplementierung von `ADR-0040`s `ClockPort` (`slice-044`).
- `LH-FA-RET-005` (Sichtbarkeit blockierender Consumer): neue View
  `cdc.retention_blockers` zeigt je Quelle den aktuell die Löschung
  blockierenden Consumer (`slice-045`).
- `LH-FA-RET-006` (Kontrolle des Datenwachstums): `cdc.metrics` trägt eine
  neue `cdc_storage_bytes`-Zeile über `pg_relation_size('cdc.change')`
  (`slice-046`).
- `BEO-PGC/retention-keine-loeschausfuehrung` (real fehlende
  Löschausführung, gefunden bei einer repo-weiten
  Fähigkeiten-Bestandsaufnahme vor `welle-13`s Eröffnung) ist damit
  aufgelöst — siehe unten.
- Ein während `slice-044` real gefundenes, außerhalb des ursprünglichen
  Welle-Scopes liegendes Sicherheits-Loch (`cdc_admin` ohne `DELETE`-Grant
  auf `cdc.transaction`/`cdc.change`) wurde über einen eigenen
  Architect-Zug geschlossen: [`ADR-0053`](../../adr/0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md)
  (`Supersedes ADR-0047`, auf die Rollen-Zuweisungstabelle begrenzt).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der vorab eingeholte Architect-Verdikt zur Retention-Löschausführung hielt
  über alle vier Slices: keine neue ADR nötig für die Domain-/Port-Ebene,
  das View-Owner-Muster trug `cdc.retention_blockers` und
  `cdc_storage_bytes` beide ohne `cdc_reader`-Rollenerweiterung — real
  durch je einen eigenen Rollen-Test bestätigt (`slice-045`, `slice-046`).
- Das etablierte E2E-Backdating-Muster (`UPDATE
  cdc.transaction.committed_at`, statt reale Zeit abzuwarten) trug
  wiederholt und ließ `slice-044`s Löschausführungs-Beleg deterministisch
  und schnell bleiben.
- `slice-043` und `slice-045` fanden je einen realen, vor dem Review
  entdeckten Bug (verwaiste `cdc.transaction`-Zeilen nach Löschung;
  liegen gebliebene Test-Consumer, die den bestehenden
  `slice-044`-Retention-Beleg dauerhaft blockiert hätten) — beide durch
  saubere Ursachen-Isolation (Rot-Grün-Vergleich gegen unveränderten
  `HEAD`) statt Symptom-Behandlung behoben.
- Der Implementer von `slice-046` zog alle DoD-Checkboxen bis auf die
  strukturell erst nach dem Rollenwechsel erfüllbare „Review
  durchgeführt"-Zeile bereits selbst nach — sichtbare Wirkung der in
  `slice-044` geschärften Implementer-Instruktion
  (`.claude/commands/implement-slice.md` Schritt 20).

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Zwei wellenlose Architect-Züge wurden nötig, die `welle-13`s ursprüngliche
  Planung nicht vorsah: (1) der `cdc_admin`-Rollen-Grant-Konflikt
  (`ADR-0053`, Modul 8 §Konflikt-Pfad Verdikt 2, ausgelöst durch einen
  realen, welle-13 §6 explizit als Out-of-Scope benannten Fund während
  `slice-044`), (2) die Steering-Loop-Verkörperung für
  `BEO-PGC/slice-chronik-in-code-kommentar` (3×, drittes Vorkommen als
  Reviewer-HIGH in `slice-044`s Review) — beide dokumentiert in `welle-13.md`
  §6 bzw. `slice-044`s Plan-Nachzug.
- Eine neue, bei `slice-045` erstmals eigenständig registrierte
  Symptom-Klasse (`BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde`) trat
  ein zweites Mal bei `slice-046` auf (2×, unter der Schwelle): Die bereits
  verkörperte `BEO-PGC/dod-checkbox-nachzug`-Regel deckt weder den
  Implementer-eigenen Lauf (Review kommt strukturell danach) noch — anders
  als bislang angenommen — jedes saubere Review ohne Fixrunde ab. Kein
  Architect-Zug ausgelöst (unter 3×), aber ein Muster, das bei einer dritten
  Wiederholung fällig würde.
- `slice-046`s Reviewer stufte sein eigenes Finding F-1
  (Beobachtungs-Register-Ausgang für `retention-keine-loeschausfuehrung`)
  fälschlich als sofort fällig ein, weil er `welle-13.md` §3 (dort bereits
  explizit als Wellen-Closure-Trigger benannt) nicht in seinem
  Eingangs-Kontext führte — der Verifier korrigierte das eigenständig
  gegen den Wellen-Plan, keine Fixrunde nötig.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle-Closure neu 3× — der Normalfall.
`BEO-PGC/slice-chronik-in-code-kommentar` erreichte 3× bereits während
`slice-044` selbst (eigener, wellenloser Architect-Zug vor dieser
Welle-Closure, der Architect-Verdikt zur Slice-Chronik in
Code-Kommentaren)
und wird hier nicht erneut verarbeitet.
`BEO-PGC/retention-keine-loeschausfuehrung` wird unten als eigener,
direkter Auflösungsfall (nicht über die 3×-Schwelle) behandelt.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. In
dieser Welle **verkörpert** (direkt, nicht über die 3×-Schwelle, analog zu
`verwaltung-keine-sql-administration`/`schema-evolution-nicht-dynamisch`):
`retention-keine-loeschausfuehrung` (0×, seit welle-13). In dieser Welle
**neu angelegt**: `dod-checkbox-nachzug-review-ohne-fixrunde` (1× → 2×,
`slice-045`/`slice-046`), `architect-verdikt-rollen-scope-luecke` (1×,
`slice-044`, wellenloser Fund). Bereits während dieser Welle verkörpert
(wellenloser Architect-Zug, nicht Teil dieses Closure-Schritts):
`slice-chronik-in-code-kommentar` (3×, `seit slice-044`). Unter der
Schwelle, unverändert von dieser Welle: `dod-checkbox-nachzug-architect-pfad`
(1×), `commit-traceability-kein-vorab-hook` (2×),
`github-actions-unverifizierbar-lokal` (1×), `rollen-test-abdeckungsluecken`
(2×), `schema-rollout-fremdobjekte` (1×), `test-isolation-geteilter-zustand`
(1×), `test-runner-stiller-ausschluss` (1×, `slice-046` vermied bewusst ein
neues Auftreten), `walsender-wirksamkeit` (1×). Bereits verkörpert,
unverändert: `a-check-null-abdeckung` (3×), `d-migrate-nacharbeit` (5×),
`dod-checkbox-nachzug` (3×), `lese-doppelquelle` (3×), `plan-nachzug` (2×),
`plan-vorlagen-defekt` (3×), `schema-evolution-nicht-dynamisch` (2×),
`verwaltung-keine-sql-administration` (0×, `seit welle-12`),
`architect-verdikt-ablageort-uneinheitlich` (0×, gestrichen). Bereits
eingetreten, unverändert: `adapter-fehler-ausgang` (2×),
`cdc-capture-lag-real` (2×), `health-endpoint-heartbeat` (2×),
`rollen-verdrahtung` (4×), `spec008-replication-luecke` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — `welle-13` schließt vollständig mit ihren vier Slices; keine
Fortsetzung wurde als Folge-Slice angelegt.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle vier Slices (`slice-043`, `slice-044`, `slice-045`, `slice-046`) in
  `done/`.
- `make gates` grün (Planner-Lauf zur Closure, Commit `1564fb9` und danach
  unverändert).
- Alle in `welle-13` §3 genannten Closure-Kriterien in einem einzigen,
  realen `make test-integration`-Lauf gemeinsam bestätigt: eine
  freigegebene Change-Zeile (`'RetentionOld'`, id=200) wurde real entfernt,
  nachdem beide Consumer sie freigaben (`LH-FA-RET-003`), eine zu junge
  Zeile (`'RetentionYoung'`, id=201) blieb durchgehend erhalten
  (`LH-FA-RET-004`); `TestMVPRetentionBlockersViewShowsFurthestBehindConsumer`
  PASS (`LH-FA-RET-005`); `TestMVPMetricsCarriesStorageBytes` PASS
  (`LH-FA-RET-006`).
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig". Kein offener Carveout im Repo. Kein
  bootstrap-aware Gate berührt. `ADR-0053`s Re-Evaluierungs-Trigger
  (Sicherheitsvorfall/Compliance-Anforderung, oder eine vierte CDC-Rolle)
  ist nicht eingetreten — bleibt unverändert `Accepted`.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — kein
  Steering-Loop-Eintrag mit `liegt in` in dieser Welle-Closure selbst
  (die `slice-chronik-in-code-kommentar`-Verkörperung trägt ihren eigenen,
  bereits geprüften Anker aus dem wellenlosen Architect-Zug, siehe
  `slice-044`s Closure-Notiz), nichts Neues zu prüfen. Folge-Slice — keiner
  genannt, nichts zu prüfen. Register — alle in dieser Welle berührten
  Verzeichnisse (`retention-keine-loeschausfuehrung`,
  `dod-checkbox-nachzug-review-ohne-fixrunde`,
  `architect-verdikt-rollen-scope-luecke`) existieren mit nicht leerem
  `evidence/` (`retention-keine-loeschausfuehrung`: 0 Belege, bereits im
  eigenen `observation.md` unter „Benannt, nicht gezählt" transparent
  ausgewiesen — direkte Auflösung ohne 3×-Beleg, keine stille Lücke; die
  beiden übrigen: real vorhandene `evidence/*.md`-Dateien).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; alle vier Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
