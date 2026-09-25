# Slice sdk-public-doc-check-gate: Der Kennungs-Wächter der SDK-Dateien (`make sdk-public-doc-check`) wird ein Gate in `make gates`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (offizielle
Client-Bibliotheken; ihre Kommentare und Fehlertexte erreichen Anwender über die
Pakete), [`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Umfang der Packages), [`AGENTS.md`](../../../../AGENTS.md) §3.6 (ein Gate braucht
eine ADR als Träger — die bestehenden Gates tragen
[`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md),
[`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md),
[`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)).

**Berührte Spec-Stellen:** — (Harness-Werkzeug; keine Spec-Stelle).

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure von `slice-sdk-readme-nutzerdoku`
(Verifikation V-4). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make sdk-public-doc-check` läuft in jedem `make gates`. Ein Rückfall —
eine interne Kennung (`SPEC-`/`ADR-`/`ARC-`/`LH-FA-`/`LH-QA-`, Slice-/Welle-Name)
in einer Datei unter `sdks/` — färbt das Gate beim Commit rot, nicht erst beim
manuellen `make sdk-pack-*` oder beim Release.

**Ausgangslage (gemessen 2026-09-25, `slice-sdk-readme-nutzerdoku`):** Das Ziel ist
netzlos und ein `grep`; `time make sdk-public-doc-check` druckt `real 0m0,013s`,
`time make test-sdk-public-doc-check` `real 0m0,195s`. Die drei
`sdk-*-release.yml` fahren `make sdk-pack-*` vor jedem Publish, und die Pack-Ziele
hängen vom Wächter ab — ein Rückfall wird nicht ausgeliefert, fällt aber nicht bei
jedem Push auf. Die Begründung „ein weiteres Gate ändert die Gate-Liste“ in
`harness/mk/sdk.mk` und `harness/README.md` (§Sensors, Zeile
`make sdk-public-doc-check`) nennt den Aufwand einer Änderung der Gate-Liste; als
Grund gegen ein netzloses, 0,013 s schnelles Gate trägt er nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Erweiterung der Muster oder des Prüfumfangs** — der Wächter behält Muster und
  Prüfumfang; eine andere Reichweite ist eine andere Entscheidung (die Python-Stubs
  prüft `test_public_text.py` im Bau, sie entstehen erst zur Bauzeit).
- **Ein Wächter für die C#- und Kotlin-README-Beispiele** — ein anderer Vorgang
  (Bau-Kontexte, siehe `slice-sdk-readme-nutzerdoku` §1).
- **Der Tabellentest als Gate** — offen für die Entscheidung der ADR; der Prüflauf
  ist der Gegenstand.

## 2. Definition of Done

- [ ] ADR (Architect, `Accepted`, Index in `docs/plan/adr/README.md`): der
      Kennungs-Wächter der SDK-Dateien ist ein Gate; Festlegung zum Tabellentest,
      zur Sensor-Bindung und zur Reichweite (Lauf im Repo-Wurzel-`make gates`,
      netzlos). *Zu belegen durch:* die ADR selbst.
- [ ] `harness/mk/sdk.mk` hängt das Ziel an `GATE_CHECKS`; `harness/sensors/sdk-public-doc-check.md`
      (Vertrag, Grenze, Bindung, wie die übrigen Sensor-Dateien),
      `harness/README.md` §Sensors (Zeile in der Gate-Tabelle und in `make gates`),
      die Begründung „kein Gate“ in `harness/mk/sdk.mk` und `harness/README.md`
      entfernt bzw. auf den Ist-Zustand gezogen (Suchlauf §3.13 über
      `git grep -n 'sdk-public-doc-check' -- harness AGENTS.md docs`).
- [ ] Beleg der Eingabeseite: eine Kennung in eine Datei unter `sdks/` eingefügt
      färbt `make gates` rot (Exit ≠ 0, Meldung des Wächters gedruckt), die
      Rücknahme grün.
- [ ] `make gates` grün (Exit-Code ungefiltert gesichert,
      [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow, kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo
      **mit** Wellen von der nächsten Welle-Closure geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/00NN-….md`, `docs/plan/adr/README.md` | neu / update (Architect) | Träger des Gates ([`AGENTS.md`](../../../../AGENTS.md) §3.6); der Slice startet erst mit der `Accepted` ADR. |
| `harness/mk/sdk.mk` | update | `GATE_CHECKS += sdk-public-doc-check`; der Kommentar beschreibt den Ist-Zustand. |
| `harness/sensors/sdk-public-doc-check.md` | neu | Sensor-Vertrag wie `generated-sync.md`: Gegenstand, Muster, Ausnahmen, Grenze, Sperren, Bindung. |
| `harness/README.md` | update | §Sensors: Gate-Zeile, Liste der Gates unter `make gates`, die Zeile `make sdk-public-doc-check` wechselt von „kein Gate“ zur Bindung an die ADR. |

## 4. Trigger

**Start** (`next` → `in-progress`): die ADR des Architects ist `Accepted`; kein anderer
Slice liegt in `in-progress/` (WIP-Limit 1); die Änderungen an `harness/mk/sdk.mk` und
`harness/README.md` durch parallele Züge sind gelandet (Ist-Zustand am Start lesen).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die Sensor-Bindung
  mehr als die vier Dateien der Tabelle §3 berührt.
- `in-progress` → `open` (blockiert): wenn der Architect die Aufnahme verwirft
  (dann geht dieser Slice mit `Gegenstand: entfallen` nach `done/`, und die
  Begründung „ein Rückfall fällt beim Pack, nicht beim Push auf“ steht als
  tragende Begründung in `harness/mk/sdk.mk`).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review-Report liegt vor und ist aufgelöst,
Closure-Notiz mit Steering-Loop-Eintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Wächter färbt `make gates` an Bau-Ausgaben unter `sdks/` rot** —
  `make gates` läuft an einem Arbeitsbaum mit unversionierten Dateien unter `sdks/`
  (`dist/`, `obj/`, `bin/` u. a.; der Tabellentest belegt die Ausnahmen an
  Beispieldateien). *Zu belegen durch:* ein Lauf von `make gates` am Arbeitsbaum nach
  `make sdk-pack-*`. **Ausgang:** offen bis zum Start.

## 7. Closure-Notiz

*(wird bei der Closure durch den Planner gefüllt)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `harness/` (Werkzeug-Fragment und
Sensor-Doku); die Modus-Deklaration führt nur die Default-Sub-Area `*` (`PGC`,
Greenfield) — keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-PGC` durchgegangen; Treffer:
`gate-scope-erweiterung-ohne-adr-traeger` (offen, 2×: ein Gate braucht seinen ADR-Träger,
der Slice startet deshalb erst mit der ADR),
`intern-kennungen-in-ausgelieferten-texten` (offen, 1×: der Gegenstand des Wächters).

**Modus:** alle berührten Sub-Areas GF.
