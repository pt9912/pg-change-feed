# Welle welle-d-check — `d-check`-Erweiterung: Getrackt-Status und Requirements-Traceability-Matrix — Closure-Notiz

**Welle:** welle-d-check
**Abschluss:** 2026-09-17
**Verantwortlich:** pt9912

## Was wurde geliefert?

- `tracked`-Modul in `.d-check.yml` `modules:`-Liste aktiv — Getrackt-Status
  auflösbarer, existierender Link-/Bild-Ziele gegen den git-Index läuft jetzt
  in `make docs-check`/`make gates`, nicht mehr nur im isolierten
  `make doc-tracked`-Einzel-Target (`slice-d-check-tracked-modul`).
- `--trace`/Requirements Traceability Matrix verdrahtet:
  `trace.requirements.id-pattern` gegen alle 76 `### LH-*`-Überschriften in
  `spec/lastenheft.md` (0 Nicht-Treffer), `trace.coverage` auf
  `docs/user/e2e-abdeckung.md` — 7 statt 9 reale Waisen, advisory
  (`make doc-trace`, `slice-d-check-trace-rtm`).
- Beide Fähigkeiten gleichzeitig in derselben `.d-check.yml` — der
  welle-eigene Closure-Trigger — ohne Config-Kollision (`verify-welle-d-check.md`).

## Was hat funktioniert?

- Die reale Vorabmessung bei Wellen-Eröffnung (§1 der Welle-Datei) hat sich
  über beide Slices hinweg exakt bestätigt — keine Überraschung beim vollen
  Implementer-Lauf.
- Der Architect-Zug als expliziter Slice-Trigger (`slice-d-check-tracked-modul`
  §4) hat eine echte architektonische Frage (ADR-Notwendigkeit) offengehalten,
  statt sie implizit vorwegzunehmen — auch wenn der Zug im Ist-Lauf
  übersprungen wurde (siehe unten), war die Trigger-Konstruktion selbst
  richtig.
- Der Lese-Schritt dieser Wellen-Closure (Schritt 3a/3b) hat zwei
  Beobachtungen mit 3×-Schwelle sauber verarbeitet, inklusive einer, die seit
  `welle-20` unbearbeitet lag (`BEO-PGC/regel-weiter-als-ihr-sensor`) — der
  Wellen-Betrieb ist damit der einzige Ort, an dem sie überhaupt gelesen
  wurde.

## Was ging anders als geplant?

- **`slice-d-check-tracked-modul`** übersprang seinen eigenen
  Start-Trigger (Architect-Zug vor `next→in-progress`) — der Implementer
  beantwortete die ADR-Frage selbst, mit falschem Präzedenzfall
  (`ADR-0072`/`ADR-0075` statt `structure`-Modul-Commit `f9e5a3c`). Führte zu
  einer vollen Fixrunde: nachträglicher Architect-Zug → Implementer-Fix →
  Delta-Review → Verifikation. Konsequenz: neue Beobachtung
  `BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (1×).
- **Die Verkörperungs-Fixrunde selbst** (Schritt 3b-Folge, Textänderung an
  `AGENTS.md` §3.13) enthielt ihrerseits ein HIGH-Finding derselben Klasse,
  die sie gerade in `.harness/skills/reviewer.md` einbettete
  (`zitat-nennt-die-falsche-stelle`, 4. Beleg) — eine falsche Zählung und
  Zuschreibung beim Zitieren des Reviews zu `slice-096` F-2. Direkt korrigiert
  gegen die primäre Quelle, unabhängig delta-reviewt.
- **`slice-d-check-trace-rtm`** lief dagegen ohne jede Rückkante — 0
  HIGH-Findings im Review, DoD vom Verifier ohne Einschränkung bestätigt.

## Steering-Loop-Einträge

- **`AGENTS.md` §3.13** geschärft: neuer Grenz-Absatz (Symbolname- vs.
  Zahlen-Lücke im Träger-Suchlauf, Reviewer als Schließer) und eine
  Trägerpflicht für das Suchlauf-Ergebnis im Slice-Plan selbst, statt nur im
  Lauf-Bericht — liegt in `AGENTS.md §3.13`. Auslöser: `BEO-PGC/regel-weiter-als-ihr-sensor`
  (`slice-078`, `slice-079`, `slice-096` — 3×; nur die dritte Manifestation
  verkörpert, die ersten beiden bleiben bewusst ohne Sensor).
- **`.harness/skills/reviewer.md`** geschärft: der bestehende HIGH-Punkt
  „Beleg trägt seinen Satz nicht" trägt jetzt die Verweis-Form (Zitat einer
  Stelle eines anderen Dokuments) und eine explizite Gegenprobe-Pflicht am
  Original — liegt in `.harness/skills/reviewer.md`. Auslöser:
  `BEO-PGC/zitat-nennt-die-falsche-stelle` (`slice-090`, `slice-102`,
  `slice-d-check-tracked-modul` — 3×; ein vierter Beleg traf denselben
  Verkörperungs-Commit unmittelbar).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`docs/plan/planning/observations/`](../observations/)
(BEO-PGC-Verzeichnisse). Was in dieser Welle 3× erreicht hat, steht oben
unter *Steering-Loop-Einträge*. Neu angelegt während dieser Welle (unter der
Schwelle, 1×): `BEO-PGC/start-trigger-ohne-uebergabe-artefakt`,
`BEO-PGC/anforderung-ohne-erkennbaren-nachweis`,
`BEO-PGC/konfigurierter-pfad-ohne-umbenennungs-schutz`.

## Folge-Slices

Keine unmittelbar ausgelöst. Kandidat für eine künftige Slice-Planung: die
inhaltliche Prüfung der 7 realen RTM-Waisen (`BEO-PGC/anforderung-ohne-erkennbaren-nachweis`)
— echte Lücke oder fehlende Zitierung, pro Anforderung zu klären.

## Verifikation

- Slice-Verifikationen: `docs/reviews/verify-slice-d-check-tracked-modul.md`,
  `docs/reviews/verify-slice-d-check-trace-rtm.md` — beide DoD erfüllt.
- Welle-weiter Verifikations-Beleg (Schritt 1, das *Mehr* über die
  Einzel-DoDs hinaus): `docs/reviews/verify-welle-d-check.md` — kombinierter
  `make gates`-Lauf grün, `tracked` und `trace:` gleichzeitig aktiv, keine
  Config-Kollision, `make doc-trace` Exit 0 (76 Anforderungen, 7 Waisen),
  `make doc-tracked` Exit 0 (892 Dateien, 0 Befunde) im `trace:`-aktiven
  Stand.
- Architect-Verdikte: `docs/reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md`,
  `docs/reviews/architect-verdict-welle-d-check-lese-schritt.md`.
- `make gates`: grün über den gesamten Wellen-Verlauf, Exit-Code jeweils
  direkt und ungepiped geprüft (`AGENTS.md` §3.9).
