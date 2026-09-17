# Slice slice-d-check-trace-rtm: `d-check --trace`/Requirements Traceability Matrix verdrahten

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-d-check` — Closure-Bedingung ist ein **kombinierter**
`make gates`-Lauf mit diesem Slice **und** `slice-d-check-tracked-modul`
gleichzeitig aktiv in derselben `.d-check.yml` (siehe
[`welle-d-check.md`](../welle-d-check.md) §1/§3); das geht über die eigene
DoD dieses Slice hinaus.

**Bezug:** kein einschlägiges `LH-*`, `ADR-*` oder `CO-*` — reine
`d-check`-Konfigurationserweiterung ohne Vertrags- oder Entscheidungsgegenstand.
Die RTM selbst **liest** `LH-*`-Kennungen aus `spec/lastenheft.md`, ohne sie
zu ändern — kein Rang-4-Verweis (ADR) oder Rang-5-Verweis (Roadmap) *auf* die
Spec, sondern ein Lesewerkzeug außerhalb der Source-Precedence-Kette.

**Berührte Spec-Stellen:** — (das `id-pattern` **liest** die bestehende
Heading-Form von `spec/lastenheft.md`, ändert sie nicht; kein Schreibzugriff
auf eine Spec-Stelle).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** `d-check --trace` wird über einen `trace:`-Block in `.d-check.yml`
auf die Kennungskonvention dieses Repos verdrahtet
(`trace.requirements.id-pattern: 'LH-(FA|QA)-[A-Z]{3}-\d{3}'` — **geprüft**
gegen alle 76 `### LH-*`-Überschriften in `spec/lastenheft.md`, 0
Nicht-Treffer, Default-Format `headings` passt ohne `table`-Konfiguration).
Zusätzlich liest ein `trace.coverage`-Block `docs/user/e2e-abdeckung.md` als
kuratierte Coverage-Dimension, damit Anforderungen, deren einziger Nachweis
ein E2E-Test ist (nicht ADR oder Slice), nicht fälschlich als `WAISE`
erscheinen. **Real gemessen** (Planungslauf 2026-09-17, gepinnter Digest
`v0.75.0`/`sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`,
Kopie des Bestands): ohne `trace.coverage` `76 Anforderung(en), 9 Waise(n)`;
mit `trace.coverage` auf `docs/user/e2e-abdeckung.md` `76 Anforderung(en),
7 Waise(n)` — zwei Anforderungen (darunter mindestens eine mit
ausschließlichem E2E-Nachweis) wechseln von falscher Waise zu `ok`. `--trace`
selbst bleibt **advisory** (Exit 0 auch mit Waisen); `--require-complete`
wird in diesem Slice **nicht** aktiviert (siehe Abgrenzung und
`welle-d-check.md` §6 Out-of-Scope — real gemessen: Exit 1 bei den
verbleibenden 7 Waisen).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`doc-complete`/`--require-complete` als Gate.** *Es wäre ein anderer
  Vorgang* — real gemessen bricht es sofort mit 7 Waisen
  (`LH-FA-CFG-006`, `LH-FA-CON-002`, `LH-FA-DAT-002`, `LH-FA-DAT-003`,
  `LH-FA-SST-001`, `LH-FA-SST-005`, `LH-QA-REL-004`), die nie geprüft wurden.
  Bestandsschutz-Argument: erst sehen, ob diese sieben echte Lücken sind oder
  nur fehlende Zitierung, dann entscheiden — das ist mehr als „RTM
  verdrahten" und ein eigener Folge-Vorgang.
  `doc-trace` bleibt advisory, wie `make image-stale` (`harness/README.md`
  §Werkzeuge, kein Gate).
- **`trace.requirements.modality`.** *Es wäre ein anderer Vorgang* — dieses
  Repo trägt keine durchgängige RFC-2119-Modalverb-Konvention in den
  Lastenheft-Absätzen, die eine automatische Klassifikation zuverlässig
  trüge; eine Fehlklassifikation wäre schlimmer als keine Spalte.
- **`trace.cross-consistency`.** *Es wäre ein anderer Vorgang* — dieses Repo
  pflegt keine zweite, unabhängige Anforderung→Design-Tabelle außerhalb der
  RTM selbst; es gibt keine zwei Sichten, die auseinanderlaufen könnten.
- **Sanierung der sieben realen Waisen.** *Bestand bleibt bewusst stehen* —
  ob sie echte Lücken sind (kein Test, keine ADR) oder nur fehlende
  Zitierung in einem bestehenden Artefakt, ist eine inhaltliche Prüfung pro
  Anforderung, kein Konfigurations-Vorgang; ggf. eigener Folge-Slice.
- **Versionsanhebung des gepinnten `d-check`-Digests.** *Bestand bleibt
  bewusst stehen* — `--trace` seit `v0.22.0`, `trace.coverage` seit
  `v0.41.0` (Benutzerhandbuch §11), beide weit unterhalb des gepinnten
  `v0.75.0`.

## 2. Definition of Done

- [x] **LP1:** `.d-check.yml` trägt `trace.requirements.id-pattern:
      'LH-(FA|QA)-[A-Z]{3}-\d{3}'`. `make doc-trace` läuft real gegen den
      Bestand (Exit 0, advisory) und zeigt die volle RTM auf stdout — nicht
      die Planungs-Vorabmessung übernehmen, sondern den Lauf im
      Implementer-Kontext erneut ausführen (`AGENTS.md` §3.12).
- [x] **LP2:** `.d-check.yml` ergänzt `trace.coverage: [{files:
      [docs/user/e2e-abdeckung.md], label: E2E}]`. `make doc-trace` erneut
      laufen lassen: Coverage-Spalte erscheint, Waisenzahl sinkt gegenüber
      LP1 (real geprüft, nicht angenommen) — Vergleich beider Läufe im
      Bericht dokumentieren.
- [x] **LP3:** Träger-Nachzug: `harness/README.md` §Werkzeuge bekommt eine
      Zeile für `make doc-trace` (advisory, analog `make image-stale`) mit
      Bindung auf diesen Slice; `harness/sensors/docs-check.md` erwähnt die
      jetzt konfigurierte RTM-Quellen-Abweichung vom Default (§Grenze/
      Bindung, je nachdem wo es besser trägt). Explizit vermerken: `doc-complete`
      ist **nicht** in `GATE_CHECKS`/`make gates` aufgenommen (Begründung:
      §1 Abgrenzung, reale 7-Waisen-Messung).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report [`review-slice-d-check-trace-rtm.md`](../../../reviews/review-slice-d-check-trace-rtm.md)
      (0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO), nicht merge-blockierend — DoD-Nachzug
      ohne Fixrunde (Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde).
- [ ] `make gates` grün — Exit-Code direkt und ungepiped geprüft (`AGENTS.md`
      §3.9).
- [x] Doku-Update: `harness/README.md` §Werkzeuge (siehe LP3);
      `harness/sensors/docs-check.md` falls dort der bessere Trägerort ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-PGC/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; kein Zähler wird gesetzt, er folgt aus den Dateien. Keine
      Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert. Insbesondere die 7 realen Waisen als eigene, benannte
      Beobachtung erwägen (nicht Risiko dieses Slice, da außerhalb seines
      Lieferumfangs, aber ein Fund, der sonst nirgends steht).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      geprüft von der **Welle-Closure** (`welle-d-check`), nicht hier
      einzeln (Repo mit Wellen-Betrieb für diese Welle).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` | update | `trace.requirements`- und `trace.coverage`-Block (§1 Ziel) |
| `harness/README.md` | update | `make doc-trace`-Zeile unter §Werkzeuge (advisory, LP3) |
| `harness/sensors/docs-check.md` | ggf. update | RTM-Konfigurationsabweichung vom Default vermerken (LP3) |
| `docs/plan/planning/observations/BEO-PGC/` | ggf. neues Verzeichnis oder `evidence/`-Ergänzung | siehe §8 Sichtungs-Schritt und §7 Closure-Notiz |

## 4. Trigger

**Start** (`next` → `in-progress`): sofort verfügbar — keine offene
Architektur-Frage wie beim Schwester-Slice (`slice-d-check-tracked-modul`),
da `trace`/`--require-complete` **nicht** ins verpflichtende Gate-Bündel
wandert (advisory bleibt advisory, kein neuer Gate-Vertrag).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Implementer-Lauf herausstellt, dass die 7 realen Waisen doch in diesem
  Slice mitbehandelt werden müssen (z. B. weil ein Reviewer die
  Auslassung als unvollständig bewertet) — dann zwei Lieferwerte
  (Konfiguration *und* Waisen-Auflösung), zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): falls `trace.requirements.
  id-pattern` beim vollen Lauf (statt der isolierten Planungsmessung mit
  reduziertem Baum) eine andere Anforderungszahl liefert als die gemessenen
  76 (z. B. weil eine `.template.md`-Datei unerwartet mitgelesen wird) —
  dann Blocker, Ursache klären, ggf. Carveout.

## 5. Closure-Trigger

DoD vollständig (§2) **und** `make gates` grün **und** Closure-Notiz mit
Lerneintrag geschrieben — kein Datum, keine externe Bedingung. Die
welle-weite Bedingung (kombinierter Lauf mit `slice-d-check-tracked-modul`)
steht zusätzlich in `welle-d-check.md` §3 und wird von der Welle-Closure
geprüft, nicht von diesem Slice allein.

## 6. Risiken und offene Punkte

- **`trace.coverage`s `files:`-Liste veraltet, sobald `docs/user/
  e2e-abdeckung.md` umbenannt oder verschoben wird** (das Erzeugnis von
  `make test-integration`, `BEO-PGC/generierte-artefakte-ohne-sync-sensor`).
  — **Ausgang:** <weiter offen: → BEO-PGC im Register, als Kopplung
  vermerkt / entfallen: Pfad ist stabil, kein Änderungsanlass erkennbar>
- **Die Waisenzahl (9 → 7) verschiebt sich beim vollen Implementer-Lauf**,
  weil andere `d-check`-Module (z. B. `matrix`, `ids`) den gescannten Baum
  anders beschneiden als die isolierte Planungsmessung. — **Ausgang:**
  <entfallen: Implementer-Lauf bestätigt identische Zahlen / weiter offen:
  Abweichung dokumentiert, Ursache nicht in diesem Slice auflösbar>
- **Ein zukünftiger Reviewer oder Verifier liest `--trace`s Advisory-Status
  falsch als bereits durchgesetztes Gate** (Verwechslungsgefahr zwischen
  `doc-trace` und `doc-complete`). — **Ausgang:** <entfallen: LP3s expliziter
  Vermerk in `harness/README.md` genügt / weiter offen: zusätzliche
  Reviewer-Skill-Ergänzung nötig>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln.

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor geschärft/ergänzt: was genau>
  — liegt in `<…>`. Auslöser: `BEO-<KUERZEL>/<slug>` (<Belege> — N×).
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <von der Welle-Closure `welle-d-check` geprüft>

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Dieser Slice berührt ausschließlich
die Default-Sub-Area `*` (Kürzel `PGC`, `harness/conventions.md` §Modus-
Deklaration pro Sub-Area) — reine Harness-Werkzeug-Konfiguration
(`.d-check.yml`, `harness/README.md`), kein Code, kein eigener
Sub-Area-Kandidat. Die Schwelle-≥2-Prüfung entfällt damit auf die bereits
deklarierte Default-Sub-Area; keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`grep -rli "trace\|rtm\|traceability.matrix\|waise" docs/plan/planning/observations/`,
2026-09-17) — **keine Treffer** für eine bestehende Beobachtung zu genau
diesem Gegenstand (Requirements Traceability Matrix, Waisen-Anforderungen).
Verwandt, aber nicht deckungsgleich:

- `BEO-PGC/regel-weiter-als-ihr-sensor` (**offen**, 3×, Schwelle erreicht) —
  derselbe *Familientyp* Beobachtung (Verhältnis Zusage/Sensor-Reichweite),
  anderer Gegenstand (`hostpaths`/`coverage-gate`, nicht `trace`); der
  fällige Lese-Schritt gehört der Closure von `welle-d-check`
  (Modul 6 Closure-Schritt 3a), nicht diesem Slice.
- `BEO-PGC/adr-folgepflicht-ohne-traeger-slice` und
  `BEO-PGC/plan-nachzug` — beide betreffen unvollständige Nachverfolgung
  zwischen Entscheidung und Umsetzung, ein anderer Mechanismus als die
  hier verdrahtete RTM (die RTM **liest** den bestehenden Zustand, ändert
  ihn nicht); kein direkter Bezug.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Default,
`harness/conventions.md`) — reiner Konfigurations-/Doku-Slice ohne
Produktionscode-Berührung, kein eigener Begründungsblock pro Sub-Area nötig.
