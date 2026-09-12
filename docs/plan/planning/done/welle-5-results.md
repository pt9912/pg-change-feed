# Welle 5 — Realer CDC-Capture-Lag — Closure-Notiz

**Welle:** welle-5
**Abschluss:** 2026-09-12
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-ADM-004`](../../../../spec/lastenheft.md) (messbarer CDC-Abstand) erfüllt real, nicht mehr nur
  genähert: der Quell-Commit-Zeitstempel läuft durch Decoder, Mapper,
  Domäne (slice-017) und Store-Adapter (slice-018) bis in
  `cdc.transaction.committed_at`; die Metrik `cdc_capture_lag` (slice-019,
  [`SPEC-013`](../../../../spec/pflichtenheft.md)) ersetzt die Persistenz-Zeit-Näherung `cdc_capture_lag_approx`.
- Ende-zu-Ende-Lasttest-Beleg real erbracht (`make test-integration`,
  künstliche `docker pause`-Verzögerung um einen Quell-Commit): Baseline
  ≈0,1–0,16 s, verzögert (1 s Pause) ≈1,15–1,25 s, stabil über mehrere
  unabhängige Läufe (Implementer, Reviewer, Verifier je eigenständig
  reproduziert).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der Ende-zu-Ende-Lasttest hat genau das getan, wofür er als *Mehr*
  gegenüber den einzelnen Slice-DoDs benannt war (Welle-Ziel §1): Der
  Implementer fand beim ersten Lauf eine nahezu verschwindende Verzögerung
  — statt das zu ignorieren, isolierte er die Ursache bis zu einem
  veralteten, vor slice-017/018 gebauten Container-Image und baute es neu,
  bevor er den Beleg als erbracht meldete. Kein einzelner Slice-DoD hätte
  diesen Bereitstellungs-Fehler gefunden.
- Reale Mutation-Tests und eigenständige Reproduktion durch Reviewer und
  Verifier (nicht nur Übernahme von Implementer-Behauptungen) trugen in
  allen drei Slices — inklusive einer eigenständig gegengeprüften
  Zeitzonen-Begründung (slice-017) und einer real getriggerten
  Fehlerprüfung (slice-019, `LAG_DELAY_SECONDS`-Obergrenze).
- [`ADR-0040`](../../adr/0040-clockport.md) (kein `time`-Import in `internal/domain`) wurde vom Implementer
  selbst erkannt und eingehalten, ohne dass der Auftrag sie explizit
  nannte — real bestätigt über alle drei Slices.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- Ein repo-weiter Kommentar-Bereinigungsauftrag (Hard Rule 3.7, vom
  Nutzer während der Arbeit an slice-018 konkret für `tools/schema/
  schema.yaml` und Quellcode scharf gestellt) betraf 18 Dateien
  außerhalb der ursprünglichen Slice-Pläne — eine gezielte Nachprüfung
  (`docs/reviews/review-comment-cleanup.md`) bestätigte die Änderungen
  als rein kommentar-ändernd und fachlich korrekt; Konsequenz: keine, war
  ein einmaliger Bereinigungsauftrag.
- `docs/user/benutzerhandbuch.md` (Plan-Nachzug slice-019) und
  `harness/image-hash.txt` (Plan-Nachzug slice-019, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)) lagen nicht
  in den ursprünglichen §3-Plänen — beide wurden im selben Lauf
  nachgetragen.

## Steering-Loop-Einträge

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Implementer-Workflow** geschärft: DoD-Checkboxen im §2 des Slice-Plans
  werden im selben Lauf abgehakt, sobald der Punkt materiell erfüllt ist —
  ein Sensor kann Checkbox-Wahrheit gegen freien DoD-Text nicht prüfen,
  die Regel bleibt Disziplin, nicht Gate
  — liegt in `.claude/commands/implement-slice.md` (Schritt 18
  Pre-completion-Checkliste, Schritt 21 Fixrunden-Nachzug).
  Auslöser: `BEO-PGC/dod-checkbox-nachzug` (slice-015, slice-016,
  slice-017 — 3×).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Zwei
Einträge mit neuem Ausgang in dieser Welle: `cdc-capture-lag-real`
(Ausgang von *weiter offen* auf **eingetreten** aktualisiert, Träger
`welle-5`, 2×) und `dod-checkbox-nachzug` (3×, **verkörpert** — siehe
Steering-Loop-Eintrag oben). Unter der Schwelle, zur Sichtung bei der
nächsten Slice-Planung: `adapter-fehler-ausgang` (1×), `lese-doppelquelle`
(2×), `rollen-verdrahtung` (1×), `schema-rollout-fremdobjekte` (1×),
`walsender-wirksamkeit` (1×). Bereits verkörpert, unverändert in dieser
Welle: `a-check-null-abdeckung` (3×), `d-migrate-nacharbeit` (4×),
`plan-nachzug` (2×), `plan-vorlagen-defekt` (3×).

## Folge-Slices

<!--
DERIVATIV: der Folge-Slice selbst ist eine Datei in `open/`; diese Liste
zeigt nur darauf. Deshalb braucht sie keinen eigenen Konsumenten — wohl
aber eine Deckung: jeder genannte Folge-Slice MUSS als Datei im
Planning-Lifecycle existieren (`open/`, `next/`, `in-progress/`, `done/` —
nicht nur `open/`, er kann bis zur Prüfung weitergewandert sein).
Folge-Slice-Paarung, geprüft am Ende von Schritt 3 der Closure-Prozedur.
Genannt ohne angelegt ist dieselbe Klasse wie ein halluziniertes Gate.
-->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine.

## Verifikation

<!--
Die Belege aus Schritt 1 der Closure-Prozedur. Keine Behauptung ohne
nachprüfbaren Anker (Hash, Lauf, Zahl).
-->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices (`slice-017`, `slice-018`, `slice-019`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Ende-zu-Ende-Lasttest (`make test-integration`, Welle-Closure-Trigger):
  Baseline ≈0,1–0,16 s, verzögert (1 s künstliche Pause) ≈1,15–1,25 s —
  real erbracht und von Implementer, Reviewer und Verifier unabhängig
  reproduziert (`docs/reviews/review-slice-019.md`,
  `docs/reviews/verify-slice-019.md`).
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig"/permanent bestätigt
  (`docs/plan/adr/architect-review-welle-5.md`).
