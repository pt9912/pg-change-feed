# Welle welle-1 — MVP-Grundlage — Closure-Notiz

**Welle:** welle-1
**Abschluss:** 2026-09-09
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- Go-Modul mit minimalem Einstiegspunkt; Build-/Deploy-Vertrag real
  (`make image` grün, Beleg `harness/image-hash.txt` = `sha256:447eab36…`
  am 1.27-Toolchain) — Teil-Beleg [`LH-QA-POR-003`](../../../../spec/lastenheft.md),
  [`LH-QA-OPS-001`](../../../../spec/lastenheft.md).
- Domänenkern: neun Modelle je [`ARC-001`](../../../plan/adr/0004-cdc-domainmodell.md)
  mit Invarianten in Konstruktoren ([`ADR-0029`](../../../plan/adr/0029-domain-invarianten.md)),
  `ClockPort` ([`ADR-0040`](../../../plan/adr/0040-clockport.md)) — Teil-Beleg
  [`LH-FA-DAT-001`](../../../../spec/lastenheft.md)/[`LH-FA-DAT-004`](../../../../spec/lastenheft.md),
  [`LH-FA-CAP-004`](../../../../spec/lastenheft.md)/[`LH-FA-CAP-005`](../../../../spec/lastenheft.md).
- Capture-Persist-Pfad: `CaptureInboundPort`, `ChangeStorePort`,
  `ReplicationAckPort`, `CaptureService` mit der Persist-before-ACK-Ordnung —
  Beleg [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md)
  (Ordnung an einer Stelle, kein ACK bei Persistenzfehler, Crash-Fenster
  getestet), [`LH-QA-REL-002`](../../../../spec/lastenheft.md).
- Gate-Ausbau in diesem Zug: a-check als drittes Gate im Bündel
  ([`ADR-0041`](../../../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)),
  Closure-Notiz-Wächter (structure-Regel 5) und reviews-Modul aktiviert,
  `image-stale` um die Major-Drift-Erkennung ergänzt, Go-Toolchain-Bump
  1.26 → 1.27 (bewusster Commit nach `make image-stale`-Fund).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die Rollen-Sequenz mit Übergabe-Artefakten lief erstmals vollständig
  über drei Slices (Architect-Review vor `next`, Review/Verifier in
  frischem Kontext mit selbst gefahrenen Belegen) — die Findings jeder
  Stufe flossen in die nächste ein.
- Die Konflikt-Sequenz (Modul 8) wurde bei ihrem dritten Auftreten
  gezogen und endete in verkörperten Regeln ([`ADR-0042`](../../../plan/adr/0042-transport-typen-am-port.md)),
  Commit-Message-Regel im Implementer-Briefing).
- Docker-only-Kette ohne Host-Toolchain: alle Go-Belege über den
  gepinnten Toolchain-Container, Arbeitsbaum read-only, Belege über
  stdout.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Der Erstversuch von slice-002 landete über den Auto-Push der Umgebung
  auf origin** (`e509c09`, ohne Review-Fixes) und erzeugte eine Divergenz —
  per Force-Push mit Ours-Auflösung aufgelöst (Vorgabe pt9912; der Verlauf
  ist in der slice-002-Closure-Notiz dokumentiert). Konsequenz gezogen:
  der Verlaufskorrektur-Schnitt muss vor jedem Push stehen.
- **`image-hash.txt` driftete vom HEAD** (V-1 aus verify-slice-001): der
  Beleg trug den Digest eines früheren Baums, nachdem Zeilenenden-Fixes
  den Build-Kontext geändert hatten. Konsequenz: geschärfte Regel im
  `harness/README.md` (Beleg gilt am HEAD, Re-Build vor Closure) und
  Re-Build vor dieser Closure.
- **Die ARC-IDs in Commit-Messagen** (Klasse B, 3×) zogen die
  Konflikt-Sequenz nach sich; die Regel ist jetzt im Implementer-Briefing
  verkörpert (`.claude/commands/implement-slice.md`, `· seit slice-003`).

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Sensor `a-check` ergänzt (Abdeckung):** die Layer-Globs matchen mit
  slice-003 über den vollen Baum (`domain`, `ports`, `app` + `cmd` als
  Composition Root) — die Null-Abdeckung ist entkräftet.
  — liegt in `.a-check.yml` (Layer-Globs voll besetzt) ·
  Auslöser: `BEO-PGC/a-check-null-abdeckung` (slice-001, slice-002,
  slice-003 — 3×).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/a-check-null-abdeckung/`
(`state.md`: Ausgang **verkörpert**, `· seit welle-1` — die Layer-Globs sind
voll besetzt; die Beobachtungsklasse kann wieder auftreten, wenn neue
Layer-Globs ohne Content aktiviert werden). Was in dieser Welle **3×**
erreicht hat, steht oben unter *Steering-Loop-Einträge*.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — welle-1 ist vollständig geschlossen. Welle 2 (reale
PostgreSQL-Integration, [`ADR-0010`](../../../plan/adr/0010-postgresql-cdc-store.md)-Adapter,
[`PH-TST-001`](../../../../spec/pflichtenheft.md)-Umfang) wird bei ihrer Eröffnung
geschnitten; die offenen Risiko-Ausgänge (a)/(b) aus slice-003 tragen
ihren Termin.

## Verifikation

- `make gates` grün am Abschluss-Stand (baseline-verify OK — 54 Dateien ·
  d-check 80 Dateien, 0 Befunde · a-check 0 Befunde, im Gate-Bündel).
- `make image` grün am HEAD; `docker run --rm
  ghcr.io/pt9912/pg-change-feed:dev --version` →
  `pg-change-feed 0.1.0-bootstrap` (Exit 0).
- Trigger-Audit (Schritt 2): Carveouts 0 offen (kein Bestand) ·
  bootstrap-aware Gates: Stufen aktuell (a-check, structure-Regel 5,
  reviews-Modul — dokumentiert in `.d-check.yml`/[`ADR-0041`](../../../plan/adr/README.md)) ·
  ADR-Re-Evaluierungen: keine fälligen (Kette 0038→0039→0042 und
  0036→0041 aktuell; [`ADR-0021`](../../../plan/adr/README.md)/
  [`ADR-0022`](../../../plan/adr/README.md)-Trigger nicht eingetreten).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug** (`archiv.zip`-
Target existiert nicht) — die Archivierungs-Bedingung ist in diesem Zug
**nicht eingetreten**; die drei Slice-Dateien und dieser Welle-Plan bleiben
vollständig in `done/`. Vor der ersten tatsächlichen Archivierung gilt die
Prüfpflicht aus dem Closure-Command: Geltungsbereich der Sensoren
(`structure`-Regel 5 keilt auf `done/slice-*.md` — bei Stubs in einem
Unterverzeichnis anzupassen) und Link-/ID-Pflichten im Stub prüfen.

