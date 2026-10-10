# Slice gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1: Die Pins von a-check und d-check anheben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden ist; die Roadmap führt wellenlose Arbeit nicht
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(Pin-Inventar P6 `A_CHECK_IMAGE` und P7 `DCHECK_IMAGE`/`DCHECK_DIGEST`,
Entscheidung 7: ein Pin-Update ist ein bewusster Digest-Commit),
[`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md)
(Reichweite von `make a-check`),
[`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(Verfeinerung gegen Erweiterung der a-check-Regel),
[`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
(Kante von `SPEC-040`, Re-Evaluierungs-Trigger (a) bis (c)) samt den
Quell-ADRs ihrer Herkunftstabelle
([`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
[`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md),
[`ADR-0094`](../../adr/0094-review-matrixklasse-kennung-statt-adresse.md),
[`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md),
[`ADR-0097`](../../adr/0097-observation-matrixklasse-review-verboten.md),
[`ADR-0099`](../../adr/0099-slice-welle-review-regel-zurueckgenommen.md),
[`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md),
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md),
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)),
[`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
(Kante von `SPEC-039`, Leer-Test vor `make doc-immutable`),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft
von Aussagen). Keine `LH-*`-Anforderung ist berührt: der Slice ändert zwei
Harness-Werkzeuge, ihre Festlegung und Verträge, nicht das Produkt.

**Berührte Spec-Stellen:**
[`SPEC-040`](../../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
(gegen d-check v0.86.1 gehalten, nachgezogen, wo das Verhalten sich ändert) ·
[`SPEC-039`](../../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
(nur gehalten: ihr Ausgang hängt am Leerfall von `vcs`). Der Verweis zeigt
**aufwärts**: Die Spec nennt diesen Slice nie (Baseline-Regelwerk
`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag). **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ausgangslage (Herkunft je Angabe, [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)):**

- **Pins am Stand `da913c16`** (gemessen vom Planner, `grep -n` auf die beiden
  Dateien): `a-check.mk` Zeile 9 `A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check@sha256:e8208764b119c606c92f82722813386277a65b12812d23b6107ea7a14dc25da1`
  (v0.20.0, reiner Digest-Pin ohne Tag); `d-check.mk` Zeile 6 und 7
  `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.82.0`,
  `DCHECK_DIGEST ?= sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`.
  Weitere lebende Pins beider Werkzeuge gibt es nicht: Suchlauf §3, Zeilen 1
  bis 4.
- **Ziel-Stände und übersprungene Versionen** (gemessen vom Planner am
  2026-10-10 mit `gh release list -R pt9912/a-check` bzw. `-R pt9912/d-check`
  und `gh release view <tag> --json body`):
  - a-check **v0.23.2**, Release-Digest
    `sha256:2368f7b3a84f1dc5d075edccfe2201e19947d12fcbc8eaf4df84ef162d94f422`
    (Index-Digest, `linux/amd64` + `linux/arm64`); übersprungen 0.21.0,
    0.22.0, 0.23.0, 0.23.1.
  - d-check **v0.86.1**, Release-Digest
    `sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`;
    übersprungen 0.83.0, 0.84.0-rc.1 (Vorabversion), 0.84.0, 0.85.0, 0.86.0.
  - Beide Digests sind *übernommen* aus der Release-Notiz; der Implementer
    misst sie am Tag nach (`docker buildx imagetools inspect`).
- **Changelogs** (vom Planner gelesen: `CHANGELOG.md` beider Repos per
  `gh api …/contents/CHANGELOG.md`, Abschnitte 0.21.0–0.23.2 bzw.
  0.83.0–0.86.1). Kandidaten für eine Reichweiten-Änderung, **ungeprüft** —
  das Urteil je Version ist Liefer-Punkt 1:
  - a-check: **0.21.0** neuer Optionalblock `shapes` (Befunde `shape-*`);
    **0.22.0** Dialekte `gomod`/`json` für `shapes`, `shape-*`-Meldung
    einzeilig; **0.23.0** Pin ist ein Index-Digest; **0.23.1**, **0.23.2**
    Image- und Doku-Änderungen ohne Prüfverhalten laut Notiz. `.a-check.yml`
    trägt keinen Block `shapes` (`grep -n shapes .a-check.yml`: kein Treffer).
  - d-check: **0.83.0** `targets.authority-disjoint` (opt-in); **0.84.0**
    Multi-Plattform-Index, ein Digest für GHCR und Docker Hub; **0.85.0**
    `reviews` ändert zwei Defaults („nicht rein additiv“),
    `vcs.ignore-link-targets`, `planning.closure.recursive`/`skip-pattern`,
    `structure[].skip-pattern` (opt-in); **0.86.0** `--manual`,
    `skip-allows-empty` (opt-in), `reviews.match: name` deckt nur den längsten
    Namen („nicht rein additiv“); **0.86.1** Abhängigkeit, „keine
    Verhaltensänderung“. Die Liste `modules:` in `.d-check.yml` (Zeile 8) führt
    `reviews`, `targets`, `planning` und `vcs` nicht; `make doc-planning`,
    `make doc-targets` und `make doc-immutable` aktivieren je eines davon
    einzeln (`d-check.mk`, gemessen mit `grep -n -- '--enable' d-check.mk`).
- **Läufe mit den neuen Ständen** — *übernommen* (Orchestrator, nur lesend, am
  Stand vor `11da9df8`):
  `docker run --rm --network none -v "$PWD:/repo:ro" ghcr.io/pt9912/d-check:v0.86.1`
  → Exit 0, `1833 Datei(en) geprüft, 0 Befund(e)`; a-check v0.23.2 per Digest
  gegen `/src` → `gesamt: 0 Befund(e)`. Null Befunde belegen **keine** gleiche
  Reichweite (`harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs).
- **`SPEC-040`** gibt laut
  [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
  Entscheidung 1 das Werkzeug mit der Eingabe `.d-check.yml` wieder; ihre
  Messungen liefen mit v0.82.0 (Fitness Function dort). „Weicht das Werkzeug
  von `SPEC-040` ab, ist das Werkzeug falsch. Ändert sich die Festlegung, wird
  `SPEC-040` fortgeschrieben, und diese ADR bleibt ihre Kante.“ Ihre
  Re-Evaluierungs-Trigger: (a) d-check bietet einen Schlüssel, der
  `matrix.exempt-paths` auf die Status-Prüfung beschränkt oder den Marker je
  Regel abschaltbar macht; (b) ein Review findet den Marker in `spec/` oder
  einen neuen `exempt-paths`-Eintrag ohne Angabe; (c) eine Festlegung in
  `SPEC-040` widerspricht einer Quell-ADR. Die Changelogs 0.83.0–0.86.1 nennen
  keine Änderung an `matrix` (gelesen, nicht gemessen).
- **a-check hat keine Festlegung** im Pflichtenheft (`grep -n a-check
  spec/pflichtenheft.md`: kein Treffer). Sie ist Gegenstand von
  `slice-spec-festlegungen-code-gates` (Datei in `open/`, §1 „Ziel“).

**Ziel:** `a-check.mk` pinnt a-check **v0.23.2** und `d-check.mk` d-check
**v0.86.1** mit den Release-Digests, je Werkzeug in einem eigenen Pin-Commit;
vor jedem Pin-Commit steht im Plan je übersprungener Version „keine
Reichweiten-Änderung“ oder die Änderung mit ihrem Artefakt;
[`SPEC-040`](../../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
ist gegen das Verhalten von v0.86.1 gehalten und, wo es sich ändert,
nachgezogen; `make gates` ist grün.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Festlegung von `make a-check` im Pflichtenheft schreiben** — Folge-Slice
  `slice-spec-festlegungen-code-gates` (Datei in `open/`; sein Ziel nennt „was
  ein verletzter Edge ist, welche Ausnahmen zählen“). Er beschreibt das
  Werkzeug am Pin seiner Arbeit; liegt dieser Slice vorher in `done/`, ist das
  v0.23.2. Die Übergabe steht als Zeile in dessen §6 (Risiko „Werkzeug-Stand
  der a-check-Festlegung“), mit diesem Plan committet.
- **Ein opt-in-Schlüssel der neuen Stände setzen** (`shapes` in `.a-check.yml`;
  `targets.authority-disjoint`, `vcs.ignore-link-targets`,
  `planning.closure.recursive`/`skip-pattern`, `structure[].skip-pattern`,
  `skip-allows-empty` in `.d-check.yml`) — ein anderer Vorgang: jeder fügt eine
  Prüfung hinzu oder nimmt etwas aus und braucht seine eigene ADR (Muster
  [`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
  Entscheidung 1). `targets.authority-disjoint` gehört inhaltlich zu
  `slice-harness-gate-index-werkzeug-teile` (Datei in `open/`, Gate-Index in
  Teilen), der hier nicht vorgegriffen wird.
- **Ein Modul in die Liste `modules:` aufnehmen** (`reviews`, `targets`,
  `planning`, `vcs`) — anderer Vorgang aus demselben Grund; die
  „nicht rein additiven“ Änderungen an `reviews` (0.85.0, 0.86.0) treffen
  `make docs-check` nur, wenn das Modul läuft. Ob ein Einzel-Ziel sie trifft,
  misst Liefer-Punkt 1.
- **Die Nennungen von v0.82.0 und `d28e9437` in `Accepted`-ADRs**
  ([`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md),
  [`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md),
  [`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md),
  [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md))
  und die Zeilen P6/P7 in
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) — Bestand bleibt
  bewusst stehen: Mess- und Stand-Angaben ihres Entscheidungszeitpunkts, keine
  lebenden Pins ([`AGENTS.md`](../../../../AGENTS.md) §3.5).
- **Die Baseline heben** (`make pin-stale-baseline`) — anderer Vorgang mit
  eigenem Bump-Ablauf (`harness/targets/pin-stale.md` §Bump-Ablauf); d-check
  0.84.0 nennt intern einen Baseline-Pin v6.17.0, das betrifft das Werkzeug-Repo,
  nicht dieses.
- **Kein Produkt-Code** — Schicht-Abgrenzung: der Slice berührt `a-check.mk`,
  `d-check.mk`, `spec/pflichtenheft.md` §7 (nur `SPEC-040`, falls nötig
  `SPEC-039`), die Verträge `harness/sensors/a-check.md` und
  `harness/sensors/docs-check.md` (nur, falls eine bewegte Eigenschaft dort
  steht) und diesen Plan.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Alle Beleg-Angaben dieser Liste sind **Zusagen** („zu belegen durch …“): der
Planungsstand hat außer den in §1 genannten Lesungen keinen Lauf im Repo
gefahren ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
Instanz B).

- [x] **Reichweite je übersprungener Version, vor jedem Pin-Commit
      (Liefer-Punkt 1).** Ein committeter Abschnitt „Reichweite
      a-check v0.21.0–v0.23.2 und d-check v0.83.0–v0.86.1 — Belege“ in diesem
      Plan trägt je Version (a-check 0.21.0, 0.22.0, 0.23.0, 0.23.1, 0.23.2;
      d-check 0.83.0, 0.84.0, 0.85.0, 0.86.0, 0.86.1; die Vorabversion
      0.84.0-rc.1 mit 0.84.0) „keine Reichweiten-Änderung“ mit Grund oder die
      Änderung, gehalten gegen die ADR mit der Reichweite
      ([`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md)/[`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
      für a-check, die Herkunftstabelle von
      [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
      Entscheidung 1 für `make docs-check`,
      [`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md) für das Modul `commits`
      in `make commit-traceability`), mit ihrem Artefakt. Gemessen wird mit
      dem neuen Digest auf der Kommandozeile (`A_CHECK_IMAGE=…`,
      `DCHECK_DIGEST=…`), ohne die `.mk`-Dateien zu ändern: `make a-check`,
      `make docs-check`, `make commit-traceability`, dazu die Einzel-Ziele
      `make doc-immutable` (leere und nicht leere Range, Ausgang nach
      [`SPEC-039`](../../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)),
      `make doc-planning` und `make doc-targets` mit altem und neuem Digest
      (Exit, Befundzahl). Eine **Erweiterung** hat vor dem Pin-Commit ihre ADR
      oder ihr Verdikt (Architect-Zug, Modul 8); ohne sie landet der Pin nicht.
      *Zu belegen durch:* den Abschnitt mit Befehl, Stand und gedruckter Zeile
      je Lauf; der Commit des Abschnitts liegt vor den Pin-Commits
      (`git log --oneline`).
- [x] **Pins auf v0.23.2 und v0.86.1 (Liefer-Punkt 2).** `a-check.mk` trägt
      `A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check@sha256:2368f7b3a84f1dc5d075edccfe2201e19947d12fcbc8eaf4df84ef162d94f422`,
      `d-check.mk` `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.86.1` und
      `DCHECK_DIGEST ?= sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`
      — je Werkzeug ein Commit, der sonst nichts ändert; Kommentare in den
      beiden Dateien, die einen Stand nennen, sind mit dem Suchlauf §3
      geprüft. *Zu belegen durch:* `docker buildx imagetools inspect` je Tag
      (gedruckter Digest gleich dem Release-Digest), `make pin-stale-acheck`
      und `make pin-stale-dcheck` ohne `DRIFT` (gedruckte Zeilen; braucht
      Netz), `git show --stat` je Pin-Commit (eine Datei).
- [x] **`SPEC-040` gegen v0.86.1 gehalten (Liefer-Punkt 3).** Für jeden Satz
      von
      [`SPEC-040`](../../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
      den eine Version aus Liefer-Punkt 1 berührt, eine Zeile *Satz → gilt
      unverändert · nachgezogen (neuer Wortlaut) · Frage an den Architect*,
      mit den drei Prüfungen aus `.claude/commands/plan-welle.md` Schritt 6
      (Festlegung im Pflichtenheft): (a) Gegenprobe gegen die Quell-ADRs der
      Herkunftstabelle von `ADR-0163`, (b) Anschluss-Frage je geändertem
      Ausgang oder Randfall, (c) Messung am Werkzeug — je solcher Satz ein Lauf
      mit v0.86.1 (Befehl, konstruierte Eingabe in einem Wegwerf-Verzeichnis,
      gesehene Ausgabe) oder *übernommen*/*hergeleitet*. Dazu je
      Re-Evaluierungs-Trigger (a) bis (c) von `ADR-0163` eine Zeile
      „eingetreten“ (mit Architect-Zug) oder „nicht eingetreten“ (mit Grund).
      Ändert sich ein Satz, steht der Nachzug in `spec/pflichtenheft.md` mit
      einer Zeile der Historie-Tabelle; stützt sich eine Änderung auf keine
      Quell-ADR, entscheidet der Architect vor der Closure (`ADR-0163`
      Entscheidung 1, Trigger (c)). Dieselbe Zeile für `SPEC-039`, soweit ihr
      Ausgang am Leerfall von `vcs` hängt. *Zu belegen durch:* die Tabelle in
      §3 und `make docs-check` Exit 0 nach dem Nachzug.
- [x] `make gates` grün, Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: [`review-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md`](../../../reviews/review-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md)
      (1 MEDIUM, 1 LOW, 3 INFO; nicht merge-blockierend). Verifikation
      **bestanden**:
      [`verify-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md`](../../../reviews/verify-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — die
      Roadmap führt unter *Offene Wellen* keine Welle, also trägt sie die
      Slice-Closure selbst (nach dem `git mv`).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf, inklusive
einer Testdatei-Zeile mit der Akzeptanzkriterien-ID in `Begründung`
(Modul 9 §Minimal Agent Workflow). -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| dieser Plan, neuer Abschnitt „Reichweite a-check v0.21.0–v0.23.2 und d-check v0.83.0–v0.86.1 — Belege“ | update | Urteil je Version mit Befehl, Stand und gedruckter Zeile (Liefer-Punkt 1), eigener Commit vor den Pin-Commits |
| `a-check.mk` (`A_CHECK_IMAGE`) | update | Release-Digest v0.23.2 (Liefer-Punkt 2), eigener Commit |
| `d-check.mk` (`DCHECK_IMAGE`, `DCHECK_DIGEST`) | update | Tag v0.86.1 und Release-Digest (Liefer-Punkt 2), eigener Commit |
| `spec/pflichtenheft.md` §7, `SPEC-040` (und Historie-Zeile) | update, nur falls Liefer-Punkt 1 eine Änderung zeigt | Nachzug der Festlegung (Liefer-Punkt 3); `SPEC-039` nur gehalten |
| dieser Plan, Tabelle „`SPEC-040` gegen v0.86.1“ | update | Gegenprobe, Anschluss-Frage, Messung am Werkzeug, Trigger von `ADR-0163` (Liefer-Punkt 3) |
| `harness/sensors/docs-check.md`, `harness/sensors/a-check.md` | update, nur falls der Suchlauf dort eine bewegte Eigenschaft findet | Träger-Nachzug ([`AGENTS.md`](../../../../AGENTS.md) §3.13); der Vertrag verweist auf die Festlegung, statt sie zu doppeln (`BEO-PGC/vertrag-doppelt-die-festlegung`, §8) |
| `docs/plan/planning/open/slice-spec-festlegungen-code-gates.md` §6 | update (mit diesem Plan committet) | Übergabe-Zeile „Werkzeug-Stand der a-check-Festlegung“ (§1 Abgrenzung) |
| `spec/pflichtenheft.md` §7 | **nicht realisiert** (Implementer) | `SPEC-040` und `SPEC-039` gelten gegen v0.86.1 unverändert (Abschnitt „`SPEC-040` gegen v0.86.1“); keine Historie-Zeile |
| `harness/sensors/docs-check.md`, `harness/sensors/a-check.md` | **nicht realisiert** (Implementer) | der Suchlauf findet dort keine bewegte Eigenschaft (Ergebnis am Stand `diff`) |
| dieser Plan, Tabelle „`SPEC-040` gegen v0.86.1“ | im Commit des Reichweiten-Abschnitts (Implementer) | die Messung lief vor den Pin-Commits; Ansatz Schritt 5 trägt danach nur den Suchlauf und die Häkchen |

**Ansatz — Commit-Folge:**

1. **Messung vor dem Pin** (Liefer-Punkt 1, neue Digests auf der
   Kommandozeile), Belege in diesen Plan, eigener Commit.
2. **Architect-Zug**, nur falls Liefer-Punkt 1 eine Erweiterung oder einen
   eingetretenen Trigger von `ADR-0163` zeigt (Planner → Architect → Planner,
   Modul 8): ADR oder Verdikt vor dem betroffenen Pin-Commit.
3. **Pin-Commit a-check**, danach `make a-check` und `make pin-stale-acheck`.
4. **Pin-Commit d-check**, danach `make docs-check` und `make pin-stale-dcheck`.
   Die beiden Pin-Commits sind voneinander unabhängig; die Reihenfolge ist
   frei.
5. **`SPEC-040`** gehalten oder nachgezogen (Liefer-Punkt 3), Träger nach dem
   Suchlauf, eigener Commit; danach `make gates`.

**Teil-Range für den Verifier.** Ändert kein Commit dieses Slice eine
`Accepted`-ADR oder einen MR-Eintrag, läuft `make doc-immutable` über die
ganze Range; sonst um den betreffenden Commit nach
[`SPEC-039`](../../../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge).
Erwartet ist der erste Fall (hergeleitet aus §1: kein Gegenstand in
`docs/plan/adr/` oder `harness/conventions/`).

**Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13).** Bewegte
Eigenschaften: der Stand beider Werkzeuge (Digest, Tag) und, falls Liefer-Punkt
3 etwas nachzieht, die Aussagen von `SPEC-040`. Suchraum der ganze Baum ohne
`.harness/baseline/**`, `docs/reviews/**` und `done/**` (Records). Muster:
Symbolname (alte und neue Digests, Tags), Beschreibung (`SPEC-040`,
Summenzeile `Datei(en) geprüft`), dazu `Index-Digest` (0.23.0/0.84.0: der Pin
ist ein Index-Digest; die Vergleichs-Bibliothek `tools/harness/lib-pin-compare.sh`
nennt den Begriff bereits). Gemessen vom Planner am Parent `da913c16`
(Plan-Datei lag dort noch nicht vor):

```suchlauf
da913c16 1 -n 'e8208764' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 0 -n '2368f7b3' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 19 -n 'v0\.82\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 3 -n 'd28e9437' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 0 -n '3e0b9779' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 0 -n 'v0\.86\.1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 39 -n 'SPEC-040' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 6 -nE 'Datei\(en\) geprüft' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
da913c16 35 -n 'Index-Digest' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 1 -n 'e8208764' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 2 -n '2368f7b3' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 18 -n 'v0\.82\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 2 -n 'd28e9437' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 1 -n '3e0b9779' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 1 -n 'v0\.86\.1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 39 -n 'SPEC-040' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 6 -nE 'Datei\(en\) geprüft' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 36 -n 'Index-Digest' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
```

Verteilung (gemessen, `git grep … | cut -d: -f2 | sort | uniq -c`): Zeile 1 —
`a-check.mk` 1; Zeile 3 — `d-check.mk` 1, [`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
7, [`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
2, [`ADR-0162`](../../adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
3, [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
3, zwei Evidence-Dateien des Registers je 1, `slice-harness-gate-index-werkzeug-teile`
1 („d-check ≥ v0.82.0“, eine Untergrenze, bleibt wahr); Zeile 4 —
`d-check.mk` 1, `ADR-0160` 1, `ADR-0161` 1; Zeile 7 — `ADR-0163` 18,
`spec/pflichtenheft.md` 8, `harness/sensors/docs-check.md` 8,
`harness/README.md` 1, Register 4; Zeile 8 — `ADR-0157` 1, `ADR-0160` 2,
Register 2, `harness/sensors/docs-check.md` 1. **Erwartet** (hergeleitet, nicht
gemessen): Zeile 1 → 1 (die Übergabe-Zeile in
`slice-spec-festlegungen-code-gates` §6), Zeile 2 → 2 (`a-check.mk`, dieselbe
Übergabe-Zeile), Zeile
3 → 18, Zeile 4 → 2, Zeile 5 → 1 und Zeile 6 → 1 (`d-check.mk`); Zeile 7 bis 9
nach dem Ergebnis von Liefer-Punkt 3. Die `diff`-Zeilen trägt der Implementer
nach, jede Abweichung mit Grund.

**Ergebnis am Stand `diff`** (Implementer, nach dem Pin-Commit d-check,
Verteilung mit `git grep … | cut -d: -f1 | sort | uniq -c`, Plan-Datei
ausgeschlossen): Zeilen 1 bis 6 wie erwartet — Zeile 1 nur
`slice-spec-festlegungen-code-gates`, Zeile 2 `a-check.mk` und dieselbe Datei,
Zeile 3 die Accepted-ADRs 0160 (7), 0161 (2), 0162 (3), 0163 (3), zwei
Evidence-Dateien und `slice-harness-gate-index-werkzeug-teile` (Untergrenze),
Zeile 4 ADR 0160 und 0161, Zeilen 5 und 6 nur `d-check.mk`. Zeile 7 und 8
unverändert (39, 6): `SPEC-040` gilt unverändert, kein Nachzug. Zeile 9 steht
bei 36 statt 35: der Zuwachs ist die Übergabe-Zeile in
`slice-spec-festlegungen-code-gates` §6, die mit diesem Plan committet wurde
(`6884fb6b`), nicht ein Träger dieser Arbeit. **Gefunden, nicht nachgezogen:**
keine Fundstelle verlangt einen Nachzug — die Nennungen in den Accepted-ADRs
und den Evidence-Dateien sind Stand-Angaben ihres Zeitpunkts (§1 Abgrenzung),
`harness/sensors/docs-check.md` trägt Summenzeile und Messungen mit Herkunft
(„gemessen in …“, „seit slice-dcheck-v0-82-0“), keinen Stand-Pin.
`harness/sensors/a-check.md` nennt keinen Stand. **Nicht gefunden:** keine
weitere Nennung eines der beiden Stände außerhalb von `a-check.mk` und
`d-check.mk`.

### Reichweite a-check v0.21.0–v0.23.2 und d-check v0.83.0–v0.86.1 — Belege

Beleg zu Liefer-Punkt 1. **Herkunft:** gemessen vom Implementer am
2026-10-10 an einem Klon von `d496d0ee` im Scratchpad
(`git clone --no-hardlinks`), vor jedem Pin-Commit. „alt“ ist der Pin am Stand
(`A_CHECK_IMAGE=ghcr.io/pt9912/a-check@sha256:e8208764…`,
`DCHECK_DIGEST=sha256:d28e9437…`), „neu“ der Release-Digest
(`…a-check@sha256:2368f7b3…`, `sha256:3e0b9779…`), je auf der Kommandozeile;
die `.mk`-Dateien sind dabei unverändert. Über `make` heißt: `make -s -C <Klon>
<Ziel> <Variable>=<Wert> > <Datei> 2>&1`, Exit von `make`.

**Digests am Tag** (`docker buildx imagetools inspect <Tag>`, gedruckte Zeile
`Digest:`, beide `MediaType: application/vnd.oci.image.index.v1+json` mit
`linux/amd64` und `linux/arm64`):

- `ghcr.io/pt9912/a-check:v0.23.2` → `sha256:2368f7b3a84f1dc5d075edccfe2201e19947d12fcbc8eaf4df84ef162d94f422`
  — gleich dem Release-Digest aus §1.
- `ghcr.io/pt9912/d-check:v0.86.1` → `sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`
  — gleich dem Release-Digest aus §1.

**Läufe am Bestand, alt gegen neu** (über `make`; „gleich“ heißt: Ausgabe
byte-gleich, bei den zwei roten `doc-immutable`-Läufen nach `sort` byte-gleich,
weil Summen- und Befundzeile in `2>&1` die Stellung tauschen):

| Ziel | alt | neu | Vergleich |
|---|---|---|---|
| `make a-check` | Exit 0, `gesamt: 0 Befund(e)` | Exit 0, `gesamt: 0 Befund(e)` | gleich |
| `make docs-check` | Exit 0, `d-check: 1848 Datei(en) geprüft, 0 Befund(e)` | Exit 0, dieselbe Zeile | gleich |
| `make commit-traceability` | Exit 0, `… 0 Befund(e)` | Exit 0, `… 0 Befund(e)` | gleich |
| `make doc-planning` | Exit 0, `… 0 Befund(e)` | Exit 0, `… 0 Befund(e)` | gleich |
| `make doc-targets` | Exit 0, `… 0 Befund(e)` | Exit 0, `… 0 Befund(e)` | gleich |
| `make doc-structure`, `make doc-tracked` | je Exit 0, `… 0 Befund(e)` | je Exit 0, `… 0 Befund(e)` | gleich |
| `make doc-immutable RANGE=d496d0ee..d496d0ee` (leer) | Exit 2, `d-check: error: Range-Leerfall "d496d0ee".."d496d0ee" — Basis und Spitze benennen denselben Commit, es wurde nichts geprüft` | Exit 2, dieselbe Zeile | gleich |
| `make doc-immutable RANGE=dc04087e..d496d0ee` (nicht leer, ohne MR-Commit) | Exit 0, `… 0 Befund(e)` | Exit 0, `… 0 Befund(e)` | gleich |
| `make doc-immutable RANGE=3f88e0c2~1..3f88e0c2` (Form-Commit) | Exit 2, 4 × `core-drift-vcs` (MR-001 bis MR-004) | Exit 2, dieselben 4 | gleich nach `sort` |
| `make doc-immutable RANGE=8e00e831~1..8e00e831` (Umzugs-Commit) | Exit 2, 1 × `core-drift-vcs` „immutable Datei gelöscht oder umbenannt“ | Exit 2, derselbe | gleich nach `sort` |
| `make doc-immutable RANGE=d496d0ee..dc04087e` (umgekehrt) | Exit 2, `d-check: error: Range-Leerfall "d496d0ee".."dc04087e" — 0 Commits, es wurde nichts geprüft` | Exit 2, dieselbe Zeile | gleich |
| `make doc-immutable RANGE=gibtsnicht..d496d0ee` | Exit 2, `d-check: error: Range-Basis "gibtsnicht" nicht auflösbar: reference not found` | Exit 2, dieselbe Zeile | gleich |
| `make doc-immutable STAGED=1` | Exit 0, `… 0 Befund(e)` | Exit 0 | gleich |

**Läufe an Mutationen, alt gegen neu.** Null Befunde am Bestand belegen keine
gleiche Reichweite (`harness/targets/pin-stale.md` §Bump eines
Gate-Werkzeugs); deshalb je Werkzeug eine Mutation, die rot werden muss, am
Klon, danach `git reset --hard d496d0ee`:

| Mutation | Ziel | alt | neu |
|---|---|---|---|
| `internal/domain/model/zz_mut.go` importiert `…/internal/adapters/driven/grpcstream` | `make a-check` | Exit 2, `wrong-direction: domain -> adapters`, `gesamt: 1 Befund(e)` | byte-gleich |
| leerer Commit mit Betreff `chore: ohne Kennung` | `make doc-commits RANGE=HEAD~1..HEAD` | Exit 2, 1 × `commit-untraceable` | gleich nach `sort` |
| Slice-Plan in `in-progress/` gelöscht; getrennt davon der Ruhe-Marker der Roadmap von `d496d0ee~1` zurück | `make doc-planning` | je Exit 0, 0 Befund(e) | gleich |
| Zeile `make gibtsnicht` in die Gate-Tabelle von `harness/README.md` | `make doc-targets` | Exit 0, 0 Befund(e) | gleich |

Die Mutationen der `d-check`-Module, die in `make docs-check` laufen, stehen im
nächsten Abschnitt (M0 bis M64, je alt und neu); sie liefen vor den
Pin-Commits und stehen deshalb im selben Commit wie dieser Abschnitt.
**Grenze:** `make doc-planning` und `make doc-targets` laufen ohne eigenen
Block in `.d-check.yml` mit den Defaults des Moduls; an keiner der drei
Mutationen wurden sie rot, in keinem der beiden Stände. Ihre Reichweite ist
damit nicht durch ein gesehenes Rot belegt, nur ihre Gleichheit am Bestand und
an den Mutationen. Beide sind keine Gates.

**Modul `reviews`, das in keinem Ziel läuft** (Risiko 1 aus §6):
`git grep -n -e '--enable reviews' -- Makefile '*.mk' tools .github` → kein
Treffer, Exit 1; `git grep -c -e '--disable reviews' -- d-check.mk` →
`d-check.mk:7` (alle sieben Einzel-Ziele mit Modul-Liste);
`.d-check.yml` Zeile 8 `modules: [links, anchors, ids, matrix, versions,
structure, hostpaths, tracked]`. Ein direkter Lauf mit dem Modul
(`docker run --rm --network none -v <Klon>:/repo:ro <Digest> --enable reviews`)
gegen den Bestand: alt Exit 0, `… 0 Befund(e)`; neu Exit 1, `… 132 Befund(e)`,
alle `review-missing`. Das ist die angekündigte „nicht rein additive“ Änderung;
sie erreicht kein Ziel dieses Repos, solange `reviews` nicht aktiviert ist
(eine Aktivierung ist nach §1 ein anderer Vorgang), und sie trifft den
Bestand bei der Aktivierung mit 132 Befunden. **Zuordnung (gemessen im Review
am Klon von `d496d0ee`, `--enable reviews`):** v0.82.0 und v0.84.0 Exit 0,
`0 Befund(e)`; v0.85.0 und v0.86.1 Exit 1, `132 Befund(e)`. Die 132 gehen
damit auf v0.85.0 zurück; die Änderung an `reviews.match` in 0.86.0 ist im
Changelog genannt, ihr Anteil an den 132 ist nicht gemessen. Die Folge bei
Aktivierung trägt `BEO-PGC/modul-reviews-inaktiv-aktivierung-trifft-bestand`
(§7).

**Einzel-Ziele außerhalb der Liste von Liefer-Punkt 1** (*übernommen* aus dem
Review, nicht vom Implementer gemessen): `make doc-trace`, `make doc-complete`,
`make doc-doctor` und `make doc-repair` stehen in `d-check.mk`; mit dem neuen
Digest sind sie alt gegen neu byte-gleich (je Exit 0; 90, 90, 2 und 1
Zeilen).

**Schlüssel-Gerüst, alt gegen neu.** `--print-config` beider d-check-Stände,
`diff`: neu sind nur die Zeilen zu `vcs.ignore-link-targets`,
`planning.closure.recursive`/`skip-pattern`/`skip-allows-empty`,
`structure[].skip-pattern`/`skip-allows-empty`, die Erklärung von `reviews`
(`promise-pattern`, `match`, `require-promises`, `recursive`, `skip-pattern`,
`skip-allows-empty`) und `targets.authority-disjoint: false`; am Block
`matrix` ändert sich keine Zeile. `grep -n -e
'authority-disjoint\|ignore-link-targets\|skip-pattern\|skip-allows-empty\|recursive'
.d-check.yml` → kein Treffer. a-check `--print-config`, `diff`: neu ist nur der
kommentierte Block `shapes`; `grep -n shapes .a-check.yml` → kein Treffer.

**Urteil je Version.** „keine Reichweiten-Änderung“ heißt: kein Gate und kein
Einzel-Ziel dieses Repos meldet oder lässt mit dem neuen Stand anderes durch;
gehalten gegen
[`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md) und
[`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(a-check: Schichten-Edges, `tooling`/`examples` unverändert in
`.a-check.yml`), die Herkunftstabelle von
[`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
Entscheidung 1 (`make docs-check`) und
[`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md) (Modul
`commits`).

| Version | Änderung laut Changelog | Urteil und Grund |
|---|---|---|
| a-check 0.21.0 | Optionalblock `shapes`, Befunde `shape-*`; Baseline des Werkzeug-Repos | keine Reichweiten-Änderung: `.a-check.yml` trägt keinen Block `shapes`; `make a-check` am Bestand und an der Mutation byte-gleich |
| a-check 0.22.0 | Dialekte `gomod`/`json` für `shapes`; `shape-*`-Meldung einzeilig | keine Reichweiten-Änderung: nur `shapes`, nicht konfiguriert |
| a-check 0.23.0 | Multi-Arch, Pin ist ein Index-Digest | keine Reichweiten-Änderung: Distribution; der Pin ist der Index-Digest (Digest oben) |
| a-check 0.23.1 | CVE-Scan je Plattform, README und Handbuch | keine Reichweiten-Änderung: kein Prüfverhalten |
| a-check 0.23.2 | Go 1.27.2 (drei CVE der Standardbibliothek) | keine Reichweiten-Änderung: kein Prüfverhalten; Läufe byte-gleich |
| d-check 0.83.0 | `targets.authority-disjoint` (opt-in) | keine Reichweiten-Änderung: Schlüssel nicht gesetzt, Default `false` (`--print-config`); `make doc-targets` gleich |
| d-check 0.84.0-rc.1, 0.84.0 | Multi-Plattform-Index, Docker-Hub-Spiegel mit demselben Digest, Abhängigkeiten | keine Reichweiten-Änderung: Distribution und Abhängigkeiten; alle Läufe gleich |
| d-check 0.85.0 | `reviews`: zwei Defaults (nicht rein additiv), neue `reviews`-Schlüssel; `vcs.ignore-link-targets`, `planning.closure.recursive`/`skip-pattern`, `structure[].skip-pattern` (opt-in); Go 1.27.2 | keine Reichweiten-Änderung eines Ziels: `reviews` läuft in keinem Ziel (oben, 132 Befunde nur bei `--enable reviews`); die opt-in-Schlüssel sind nicht gesetzt; `make doc-immutable` in sechs Ranges und mit `STAGED=1`, `make doc-planning` und `make docs-check` samt M0 bis M64 gleich |
| d-check 0.86.0 | `--manual`; `skip-allows-empty` (opt-in); `reviews.match: name` deckt nur den längsten Namen (nicht rein additiv); `--suggest-config` | keine Reichweiten-Änderung eines Ziels: kein Ziel ruft `--manual` oder `--suggest-config`, `skip-allows-empty` nicht gesetzt, `reviews` nicht aktiv und ohne `match` |
| d-check 0.86.1 | `golang.org/x/net` (CVE), „keine Verhaltensänderung“ | keine Reichweiten-Änderung: alle Läufe gleich |

**Ergebnis:** keine Erweiterung, kein Architect-Zug vor den Pin-Commits nötig
(Ansatz Schritt 2 entfällt).

### `SPEC-040` gegen v0.86.1

Beleg zu Liefer-Punkt 3, Prüfungen nach `.claude/commands/plan-welle.md`
Schritt 6. **Herkunft:** gemessen vom Implementer am 2026-10-10 am Klon von
`d496d0ee` im Scratchpad; je Mutation einzeln, danach `git checkout -- .` und
`git clean -fd`; je Mutation zwei Läufe direkt, ohne `make`:
`docker run --rm --network none -v <Klon>:/repo:ro <Digest>` mit dem Digest
von v0.82.0 (`d28e9437…`) und von v0.86.1 (`3e0b9779…`), stdout, stderr und
Exit je in eine eigene Datei. **Vergleich:** in jedem der 67 Läufe je Stand
(M34 ist der Basislauf ohne Mutation; M1b und M39 sowie M19 und M38 sind je
dieselbe Mutation) sind Exit, stderr und stdout (als Menge, nach `sort`)
beider Stände gleich; die Spalte zeigt die Ausgabe von v0.86.1. **Ausgenommen
ist die Zeile „Basislauf mit `--enable reviews`“:** sie gehört nicht zu den 67
Läufen, und dort unterscheiden sich die Stände (v0.82.0 Exit 0, v0.86.1 Exit 1
mit 132 Befunden). Die Nummern
M0 bis M43 folgen den Messungen, auf die sich
[`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
stützt (Plan von `slice-spec-festlegungen-doku-gates`, §3); M19b und M44 bis
M64 sind neu; M35 ist hier M1a, M36, M37 und M43 liefen über `make` (unter
der Tabelle).

| # | Mutation | v0.86.1 (gleich v0.82.0, außer in der Zeile `--enable reviews`) | Satz von `SPEC-040` |
|---|---|---|---|
| M34 | keine (Basislauf) | Exit 0; stdout leer; stderr `d-check: 1848 Datei(en) geprüft, 0 Befund(e)` | Fläche; Ausgänge (Summenzeile auf stderr, Exit 0) |
| M62 | `.harness/mut/x.md` mit Link auf eine fehlende Datei | Exit 0, `1848 Datei(en)` | Fläche: `.harness/**` aus `scan.ignore` |
| M63 | `docs/x.template.md` mit Link auf eine fehlende Datei | Exit 0, `1848 Datei(en)` | Fläche: `**/*.template.md` |
| M45 | neue, ungetrackte `unverfolgt.md`, Link aus `README.md` darauf | Exit 1, `target-untracked`, `1849 Datei(en)` | Fläche (neue `.md` gezählt); `tracked` |
| M0 | Baum ohne `.git` (`git archive`) | Exit 2, `d-check: error: kein lesbares git-Repository unter /repo: …` | Eingabe braucht `.git`; Exit 2 |
| M44 | Pfad `gibt/es/nicht.go` in Inline-Code | Exit 0 | `codepaths` nicht in der Liste |
| — | Basislauf mit `--enable reviews` (Abschnitt oben; nicht gleich: v0.82.0 Exit 0, 0 Befunde; nicht unter den 67 Läufen) | Exit 1, 132 × `review-missing` | Block ohne Modul (`reviews:`) wirkt nicht, solange das Modul nicht läuft |
| M1a | Link auf eine fehlende Datei | Exit 1; stdout `README.md:65` TAB `gibt-es-nicht.md` TAB `target-missing` TAB `Linkziel existiert nicht` | `links`; Ausgänge (vier Felder) |
| M18 | Link mit URL-Schema auf ein Ziel, das es nicht gibt | Exit 0 | `links`: URL-Schema kein Gegenstand |
| M17 | Link auf ein fehlendes Heading und Link auf eine HTML-`id` | Exit 1, `anchor-missing` nur für das Heading | `anchors` |
| M61 | Links auf vorhandenes Heading und vorhandene HTML-`id` | Exit 0 | `anchors` |
| M8 | Kennung `ADR-0041` nackt in Prosa | Exit 1, `id-unlinked` | `ids` |
| M7 | dieselbe Kennung in Inline-Code | Exit 0 | `ids`: Inline-Code kein Treffer |
| M25 | `SPEC-040` verlinkt auf `harness/conventions.md` | Exit 0 | `ids`: Ziel des Links nicht abgeglichen |
| M26 | `LH-QA-REL-001.a` verlinkt auf `spec/lastenheft.md` | Exit 0 | `ids`: Ziel des Links nicht abgeglichen |
| M46 | `MR-002` und `BEO-PGC/gibt-es-nicht` in Prosa | Exit 0 | `ids`: Kennung ohne Muster |
| M10 | `ADR-0041` in Inline-Code in der Architektur-Sicht | Exit 1, `matrix-forbidden` | `matrix`: Token in Inline-Code |
| M11 | Link aus der Architektur-Sicht ins ADR-Verzeichnis | Exit 1, `matrix-forbidden` | `matrix-forbidden` als Link |
| M14 | Token `docs/reviews/x.md` in einer ADR | Exit 1, `matrix-forbidden` | Regel `adr → review` |
| M31 | `ADR-0094` und `slice-099` in einem Fence der Architektur-Sicht | Exit 0 | `matrix`: Fence kein Treffer |
| M12 | `slice-099` und `slice-abc-name` in der Architektur-Sicht | Exit 1, `matrix-forbidden` nur für `slice-099` | Token `slice` nur mit Nummer |
| M47a | Link aus der Architektur-Sicht auf `AGENTS.md` | Exit 0 | Datei ohne Klasse als Ziel |
| M47b | `ADR-0041` und `docs/reviews/x.md` in Inline-Code in `harness/README.md` | Exit 0 | Datei ohne Klasse als Quelle |
| M48 | Link aus der Architektur-Sicht auf das Lastenheft | Exit 0 | keine Regel innerhalb der Spec-Straten |
| M49 | Link aus dem Slice-Plan in `in-progress/` auf einen Review-Bericht | Exit 0 | `slice` darf auf `review` verweisen |
| M32a, M32b | Token `ADR-0038` (Status `Superseded`) in Inline-Code im Slice-Plan; in einer `observation.md` | je Exit 0 | `matrix-inactive` nur an Links |
| M32c | `docs/reviews/*.md` aus `matrix.exempt-paths` entfernt | Exit 1, 57 × `matrix-inactive`; jede der 57 Zeilen trägt einen Link `](<Ziel>` (Zählung je Befundzeile gegen die Quellzeile) | `matrix-inactive` nur an Links; `exempt-paths` als Quelle |
| M16 | `slice-099` in ADR 0039 (unter `matrix.exempt-paths`) | Exit 0 | `exempt-paths`: als Quelle ganz ausgenommen |
| M50 | Link aus der Architektur-Sicht auf ADR 0039 | Exit 1, `matrix-forbidden` und `matrix-inactive` | `exempt-paths`: als Ziel Gegenstand der Regel |
| M51 | `allow-supersede-lineage: true` gestrichen | Exit 1, 16 × `matrix-inactive` | `allow-supersede-lineage` (mit dem Schlüssel: Basislauf 0) |
| M22 | Pfad in die vendored Baseline mit Tag `v6.14.0` in Inline-Code in `README.md` | Exit 1, `version-stale` | `versions` |
| M53 | derselbe Pfad in einer ADR | Exit 0 | `versions.exempt-paths` |
| M52 | Link auf eine fehlende Datei in einer ADR | Exit 1, `target-missing` | `versions.exempt-paths` nimmt `links` nicht aus |
| M9 | M8 mit `<!-- d-check:ignore -->` | Exit 0 | `d-check:ignore` für `ids` |
| M23 | M22 mit `<!-- d-check:ignore -->` | Exit 0 | `d-check:ignore` für `versions` |
| M21 | `~/<Verzeichnis>/x` in Inline-Code mit `<!-- d-check:ignore -->` | Exit 1, `hostpath-forbidden` | `hostpaths` kennt keinen Zeilen-Marker |
| M55 | M10 mit `<!-- d-check:ignore -->` | Exit 1, `matrix-forbidden` | `d-check:ignore` wirkt auf `matrix` nicht |
| M13 | M10 mit `<!-- d-check:status-provenance -->` | Exit 0 | Provenance-Marker hebt Token auf |
| M15 | M14 mit dem Marker | Exit 0 | Marker für jede Regel (`adr → review`) |
| M33 | `ADR-0094` und der Marker als Text, beide in Inline-Code, Architektur-Sicht | Exit 0 | Marker hebt Token auf |
| M54 | M11 mit dem Marker | Exit 1, `matrix-forbidden` | Link in derselben Zeile bleibt Befund |
| M27 | Linkziel über einen Zeilenumbruch, zusammengesetzt vorhanden | Exit 1, `target-missing` | Randform Ziel umbrochen |
| M28 | Anker über einen Zeilenumbruch | Exit 1, `anchor-missing` | Randform Ziel umbrochen |
| M29, M30 | Linktext umbrochen, Ziel bzw. Anker fehlt | je Exit 0 | Randform Linktext umbrochen |
| M5 | Kopfzelle `Titel` im ADR-Index umbenannt | Exit 1, `section-column-missing` | `structure`: Spalte über den Namen |
| M56a | Titel-Zelle im ADR-Index über 80 Zeichen | Exit 1, `section-cell-oversized` | `cell-max-chars` |
| M56b | Status-Zelle `X` | Exit 1, `section-cell-undersized` | `cell-min-chars` |
| M6 | Überschrift „3. Defaults und Konstanten“ umbenannt | Exit 1, `section-missing` | Regel trifft keinen Abschnitt |
| M2 | „Steering-Loop“ in einem `done/slice-*.md` ersetzt | Exit 1, `section-pattern-missing` | `require-pattern` |
| M58 | Satz „Wurde mit diesem Slice nichts verkörpert“ in Prosa in §7 | Exit 1, `section-forbidden` | `forbid-pattern` |
| M57 | derselbe Satz in Inline-Code | Exit 0 | bereinigter Abschnittstext |
| M4 | §7 leer | Exit 1, `section-empty` und `section-pattern-missing` | `non-empty` |
| M3 | offene Aufgabe `- [ ]` in §7 | Exit 1, `section-tasks-open` | `max-open-tasks` |
| M24 | offene Aufgabe in §7 von `done/release-image-scan.md` | Exit 0 | Regel gilt nur `done/slice-*.md` |
| M19, M38 | `~/<Verzeichnis>/x` in Inline-Code | Exit 1; stdout `README.md:65` TAB der Pfad TAB `hostpath-forbidden` | `hostpaths` Home-relativ; drei Felder |
| M64 | `~/<Datei>` und `~/<Verzeichnis>` ohne Schrägstrich in Inline-Code | Exit 1, 2 × `hostpath-forbidden` | `hostpaths` Home-relativ, beide Formen |
| M19b | `<Host-Wurzel>/x` in Inline-Code | Exit 1, `hostpath-forbidden` | Präfixliste |
| M59 | ein Laufwerks- und ein UNC-Muster in Inline-Code | Exit 1, 2 × `hostpath-forbidden` | Windows-Laufwerk, UNC |
| M20 | Tilde in URL, `~/.config/x`, nackte Tilde, `~<Benutzer>/x` | Exit 0 | kein Treffer |
| M60 | Home-relativer und absoluter Pfad in einem Fence; relativer Pfad in Inline-Code | Exit 0 | Fence und relativer Pfad kein Treffer |
| M1b, M39 | unbekannter Schlüssel auf oberster Ebene | Exit 2; stdout leer; stderr `d-check: error: .d-check.yml: yaml: unmarshal errors:` und `  line 348: field unbekannt not found` | Exit 2, zwei Zeilen |
| M40 | unbekannter Schlüssel unter `scan:` | Exit 2; zweite Zeile `  line 4: field unbekannt not found { … }` | YAML-Fehler mit Zeilennummer |
| M41 | `versions.current-from` auf eine fehlende Datei | Exit 2, `d-check: error: versions.current-from nicht lesbar (…)` | Exit 2 |
| M42 | unbekanntes Modul `gibtsnicht` in `modules:` | Exit 2, `d-check: error: unbekanntes Modul "gibtsnicht" in der Konfiguration — gültig: …` | Exit 2 |

**Über `make`, nur v0.86.1:** drei fehlende Linkziele, `make -s docs-check
DCHECK_DIGEST=sha256:3e0b9779… > <Datei> 2>&1`, **20 Läufe**: je Exit 2; die
Summenzeile 14-mal vor der ersten und 6-mal nach der letzten Befundzeile, nie
dazwischen; die letzte Zeile 20-mal `make: *** [d-check.mk:20: docs-check]
Fehler 1`. M43 (M39 über `make`, ein Lauf): Exit 2, die zwei Zeilen des
Fehlers, danach `make: *** [d-check.mk:20: docs-check] Fehler 2`. Beides
deckt die Sätze „über `make` jeder Exit ≠ 0 als Exit 2“, „letzte Zeile die
Fehlerzeile von `make`“ und „Reihenfolge nicht festgelegt“.

**Nicht gemessen** (*hergeleitet*, ohne Lauf mit v0.86.1): „gemessen wird die
Zelle mit Markdown-Syntax“, „eine Tabelle ohne Spalte des Namens fällt aus der
Regel“, „keine Regel zählt Zeilen“ und die Aufzählung, wem die
`structure`-Regeln gelten (Werte der Eingabe `.d-check.yml`, Datei unverändert),
„`hostpaths.prefixes` ist der Default“ (kein Block `hostpaths:` in
`.d-check.yml`, `grep -n '^hostpaths:' .d-check.yml` ohne Treffer). Das
Changelog 0.83.0–0.86.1 nennt für `structure` nur den opt-in-Schlüssel
`skip-pattern`.

**Ergebnis je Satz:** jeder Satz von `SPEC-040` **gilt unverändert**; keine
Version aus Liefer-Punkt 1 berührt einen Satz (Urteil je Version oben), und
keiner der 67 Läufe zeigt einen Unterschied zwischen v0.82.0 und v0.86.1.
`spec/pflichtenheft.md` bleibt unverändert, keine Historie-Zeile.

- **(a) Gegenprobe gegen die Quell-ADRs:** entfällt als neue Prüfung — kein
  Satz ändert sich, die Gegenprobe von
  [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
  Entscheidung 1 bleibt unberührt.
- **(b) Anschluss-Frage:** kein Ausgang und kein Randfall ändert sich (Exit,
  stdout, stderr gleich in allen 67 Läufen und in den Läufen über `make`);
  keine neue Folge für den Anwender.
- **(c) Messung am Werkzeug:** die Tabelle oben.

**Re-Evaluierungs-Trigger von `ADR-0163`:**

- (a) **nicht eingetreten** — `--print-config` v0.86.1 gegen v0.82.0 ändert am
  Block `matrix` keine Zeile (kein Schlüssel für `exempt-paths` oder den
  Marker je Regel); M13, M15, M33, M54, M55 gleich in beiden Ständen.
- (b) **nicht eingetreten** — der Slice ändert weder `spec/` noch
  `matrix.exempt-paths`; `git grep -n 'd-check:status-provenance' -- spec/`
  trifft zwei Zeilen von `spec/pflichtenheft.md` (2062, 2216), beide nennen den
  Marker in Inline-Code, keine trägt ihn als Kommentar.
- (c) **nicht eingetreten** — keine Festlegung von `SPEC-040` ändert sich.

**`SPEC-039`** gilt unverändert: Ihr Ausgang hängt am Leerfall von
`make doc-immutable` — mit v0.86.1 wie mit v0.82.0 Exit 2 bei gleicher Basis
und Spitze (`Range-Leerfall … benennen denselben Commit`) und bei der
umgekehrten Range (`Range-Leerfall … 0 Commits`), Exit 2 bei nicht
auflösbarer Basis, Exit 0 in einer nicht leeren Range ohne MR-Commit, Exit 2
mit `core-drift-vcs` um den Form- und den Umzugs-Commit (Tabelle oben).

### Pin-Commits — Belege

Beleg zu Liefer-Punkt 1 (Reihenfolge) und 2; gemessen vom Implementer am
2026-10-10 am Arbeitsbaum. `git log --oneline d496d0ee..17611a23`:

```text
17611a23 chore(harness): d-check auf v0.86.1 gepinnt, Index-Digest 3e0b9779 (ADR-0051, ADR-0163)
0845d926 chore(harness): a-check auf v0.23.2 gepinnt, Index-Digest 2368f7b3 (ADR-0051, ADR-0041)
05c14aac docs(plan): Reichweite a-check v0.21.0-v0.23.2 und d-check v0.83.0-v0.86.1 vor den Pin-Commits gemessen (ADR-0051, ADR-0163)
```

- `git show --stat 0845d926` → `a-check.mk | 2 +-`, `1 file changed`;
  `git show --stat 17611a23` → `d-check.mk | 4 ++--`, `1 file changed`.
- Nach `0845d926`: `make a-check` Exit 0, `gesamt: 0 Befund(e)`;
  `make pin-stale-acheck` Exit 0,
  `OK          A_CHECK_IMAGE (ghcr.io/pt9912/a-check:latest) == sha256:2368f7b3a84f1dc5d075edccfe2201e19947d12fcbc8eaf4df84ef162d94f422`.
- Nach `17611a23`: `make docs-check` Exit 0,
  `d-check: 1848 Datei(en) geprüft, 0 Befund(e)`; `make pin-stale-dcheck`
  Exit 0, `OK          DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.86.1) == sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`
  und `OK          DCHECK_IMAGE Tag-Frische (v0.86.1) == neuester Release`.
  Kein `DRIFT`, kein `UNBESTIMMT`: der Freshness-Vergleich trägt den
  Index-Digest (Risiko 4 aus §6).
- Kommentare in `a-check.mk` und `d-check.mk`, die einen Stand nennen: keine
  (`grep -n 'v0\.' a-check.mk` → kein Treffer; in `d-check.mk` nur Zeile 6,
  der Pin selbst).
- `make gates` am Stand `bf8e1762` (letzter Commit mit Pin, Plan und
  Suchlauf), `make gates > <Datei> 2>&1; echo $? > <Datei>`: Exit 0; gedruckt
  u. a. `baseline-verify: v6.16.0 OK — 54 Dateien …`,
  `coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`, zweimal
  `d-check: 1848 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)`
  (a-check mit `2368f7b3…`), `commit-traceability: OK — 5 Commit(s) in
  "HEAD~5..HEAD", Betreffs ohne Struktur-ID`, `generated-sync: OK — …`.

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): die Releases a-check v0.23.2 und d-check
v0.86.1 sind veröffentlicht (gemessen, §1); `in-progress/` trägt außer
`roadmap.md` keinen Slice (WIP-Limit 1; am Parent `da913c16` gemessen mit
`ls`); `Verantwortlich:` ist beim `open → next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Liefer-Punkt 1 zeigt
  für **beide** Werkzeuge eine Erweiterung mit eigenem Architect-Zug, oder
  Liefer-Punkt 3 verlangt mehr als einen Nachzug in `SPEC-040` (etwa eine
  Folge-ADR auf eine Quell-ADR, Trigger (c) von `ADR-0163`, **und** einen
  Vertrags-Nachzug). Dann wird nach Werkzeug geschnitten: ein Slice
  a-check v0.23.2, ein Slice d-check v0.86.1 samt `SPEC-040`; beide einzeln
  lieferbar, weil die Pins unabhängig sind.
- `in-progress` → `open` (blockiert): `docker buildx imagetools inspect`
  oder `make pin-stale-acheck`/`make pin-stale-dcheck` zeigt einen anderen
  Digest als die Release-Notiz (Tag nicht vertrauenswürdig, der Slice wartet
  auf einen geklärten Release); oder eine Erweiterung aus Liefer-Punkt 1
  bekommt vom Architect kein Artefakt, bevor der Pin landen müsste.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die drei Liefer-Punkte der DoD sind abgehakt mit Beleg; `make gates` endet auf
dem Commit nach dem letzten Nachzug mit Exit 0; jede Erweiterung aus
Liefer-Punkt 1 hat ihr Artefakt, und jeder eingetretene Trigger von
[`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
seinen Architect-Zug; der Review-Report liegt vor und ist aufgelöst; die
Closure-Notiz (§7) trägt den Lerneintrag und jedes Risiko aus §6 seinen
Ausgang.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ein Gate-Lauf bewegt sich über eine „nicht rein additive“ Änderung.**
  d-check 0.85.0 und 0.86.0 ändern Defaults von `reviews`; das Modul steht
  nicht in `modules:` und in keinem `--enable` von `d-check.mk` (§1, gelesen).
  Läuft es doch irgendwo mit (ein Einzel-Ziel, ein Workflow), wird ein Slice
  rot, den bisher ein längerer Report-Name deckte. *Zu belegen durch:*
  `git grep -n -- '--enable reviews\|reviews' -- Makefile '*.mk' tools .github`
  am Stand des Slice und die Läufe aus Liefer-Punkt 1. Ausgang: **weiter
  offen** — das Modul läuft in keinem Ziel, bei Aktivierung meldet es am
  Bestand 132 × `review-missing` (v0.84.0: 0; Zuordnung in §3); die Folge
  steht im Register unter `BEO-PGC/modul-reviews-inaktiv-aktivierung-trifft-bestand`
  (§7).
- **`SPEC-040` beschreibt v0.86.1 nicht mehr.** Ein Default oder eine Randform
  eines aktiven Moduls ändert sich, ohne dass das Changelog es als Änderung
  nennt (die Changelogs sind gelesen, nicht gemessen). *Zu belegen durch:* die
  Messung am Werkzeug in Liefer-Punkt 3 (c) für die Sätze, die eine Version
  berührt, und den vollen `make docs-check`-Lauf; ein Satz ohne Lauf steht als
  *hergeleitet*. Ausgang: nachgezogen in Liefer-Punkt 3, oder Architect-Zug
  nach `ADR-0163` Trigger (c).
- **Trigger (a) von `ADR-0163` tritt ein** — ein neuer Schlüssel beschränkt
  `matrix.exempt-paths` oder macht den Provenance-Marker je Regel abschaltbar.
  Die Changelogs nennen keinen (§1, gelesen); `d-check --print-config` mit
  v0.86.1 zeigt die Schlüssel. *Zu belegen durch:* diesen Lauf in Liefer-Punkt
  3. Ist er eingetreten: Folge-ADR durch den Architect, nicht in diesem Slice.
- **Der Index-Digest bricht den Freshness-Vergleich.** Ab a-check 0.23.0 und
  d-check 0.84.0 ist der Pin ein Index-Digest. `tools/harness/lib-pin-compare.sh`
  vergleicht Index-Digests (Kopfkommentar, gelesen); ob `make pin-stale-acheck`
  und `make pin-stale-dcheck` ihn nutzen, ist nicht gemessen. *Zu belegen
  durch:* beide Läufe nach dem Pin-Commit ohne `DRIFT`/`UNBESTIMMT`. Ausgang
  bei Rot: Folge-Slice am Werkzeug, der Pin bleibt (die Release-Digests sind
  aus der Notiz und `imagetools inspect` doppelt belegt).
- **Die a-check-Festlegung beschreibt einen alten Stand.** Schreibt
  `slice-spec-festlegungen-code-gates` seine Messungen vor dem Pin-Commit
  dieses Slice, beschreibt er v0.20.0. *Zu belegen durch:* die Übergabe-Zeile
  in dessen §6 (mit diesem Plan committet). Erwarteter Ausgang: entfallen,
  wenn dieser Slice zuerst schließt; sonst eingetreten mit Adresse
  `slice-spec-festlegungen-code-gates`.

Kein Workflow unter `.github/workflows/` ist berührt; [`AGENTS.md`](../../../../AGENTS.md)
§3.10 greift nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Geliefert:** Reichweiten-Abschnitt mit Urteil je Version vor den
  Pin-Commits (`05c14aac`; Liefer-Punkt 1); Pins a-check v0.23.2
  (`0845d926`) und d-check v0.86.1 (`17611a23`), je eine Datei (Liefer-Punkt 2);
  `SPEC-040` und `SPEC-039` gegen v0.86.1 gehalten, ohne Nachzug, Trigger (a)
  bis (c) von [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
  nicht eingetreten (Liefer-Punkt 3; Belege `bf8e1762`, `606b7560`). Review
  (1 MEDIUM, 1 LOW, 3 INFO, `c1e26ae7`), keine Fixrunde, kein Re-Review:
  alle Funde sind Plan-Text der Closure. Verifikation in einem Lauf
  **bestanden** (`4eb479d0`). Validator: entfällt — Harness-Werkzeug ohne
  End-Nutzer-Wert, kein MVP-Slice.
- **Was hat funktioniert:** Die Regel „Reichweite vor dem Pin-Commit“
  (`harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs) hat getragen:
  zehn Versionen mit Urteil, 67 Vergleichsläufe alt gegen neu und Mutationen
  je Werkzeug vor dem Pin; kein Pin landete ohne Messung. Die Messung fand die
  einzige Änderung, die ein Gate-Lauf nicht zeigt — das Modul `reviews`
  meldet mit dem neuen Stand am Bestand 132 Befunde, in keinem Ziel dieses
  Repos, weil es nirgends läuft. Null Befunde im Gate-Lauf hätten das nicht
  gezeigt. Reviewer und Verifier haben unabhängig dieselbe Messung gefahren
  (v0.82.0 und v0.84.0: 0, v0.85.0 und v0.86.1: 132).
- **Was ging anders als geplant:**
  (1) **Erwarteter Risiko-Ausgang trug nicht.** §6 setzte für die „nicht rein
  additive“ Änderung vorab *entfallen* („läuft heute nicht mit“). Die Messung
  zeigte eine latente Folge: bei Aktivierung von `reviews` 132 × `review-missing`.
  „Läuft heute nicht“ beantwortet die Frage nach heute, nicht die nach der
  Aktivierung; Review F-1 (MEDIUM). Der Ausgang ist *weiter offen*, siehe unten.
  (2) **Die Tabellenüberschrift war breiter als ihre Zeilen** (F-2) und
  **die Version, auf die eine Messung zurückgeht, stand ohne Messung im Text**
  (F-3: 0.85.0/0.86.0, gemessen 0.85.0). Beides im Plan berichtigt.
  (3) **Die Messliste der Einzel-Ziele war nicht aus `d-check.mk` hergeleitet**
  (F-4): `doc-trace`, `doc-complete`, `doc-doctor`, `doc-repair` fehlten.
  Der Reviewer maß sie (byte-gleich); sie stehen jetzt im Plan, als
  *übernommen* gekennzeichnet.
- **Review- und Verifikations-Pflichten in der Closure:**
  - Review F-1 / Verifikation Abweichung 1 → **weiter offen**, ins Register:
    `BEO-PGC/modul-reviews-inaktiv-aktivierung-trifft-bestand` (neu, Zähler
    1×), mit der Zahl **132** und den Ständen in `observation.md`, nicht nur im
    archivierbaren Reichweiten-Abschnitt. Frist der Sichtung: jede
    Slice-Planung (§8) und jede Aktivierung des Moduls.
  - Review F-2 (LOW) / Verifikation Abweichung 2 → **nachgezogen**: Kopf der
    `SPEC-040`-Tabelle auf „außer in der Zeile `--enable reviews`“
    eingeschränkt, die Zeile nennt v0.82.0 Exit 0 und ist von den 67 Läufen
    ausgenommen (§3). Klasse „Nachzug widerspricht Nachbar im selben Träger“ =
    `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`: Deckel bei 10× (vor dem
    Merge vom Reviewer gefunden, ≤ LOW, bekannter Träger-Typ Slice-Plan) —
    keine `evidence/`-Datei, hier mit Finding-Kennung geführt.
  - Review F-3 (INFO) / Verifikation Abweichung 3 → **nachgezogen**: die
    132 Befunde sind v0.85.0 zugeordnet (gemessen), der Anteil von 0.86.0 ist
    als nicht gemessen gekennzeichnet (§3). Klasse „Ursprung einer Aussage
    nicht gekennzeichnet“ = `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`:
    Deckel bei 23× (INFO, vor dem Merge gefunden, Träger-Typ Slice-Plan) —
    keine Datei, hier mit Finding-Kennung geführt.
  - Review F-4 (INFO) → **nachgezogen**: die vier Ziele stehen im Plan (§3),
    als *übernommen* gekennzeichnet. Die Klasse „Messliste unvollständig gegen
    Zielliste“ ist neu: `BEO-PGC/bump-messliste-unvollstaendig-gegen-zielliste`
    (Zähler 1×, offen).
  - Review F-5 (INFO) → **hingenommen, mit Grund**: der Plan benennt die
    Grenze selbst („nicht durch ein gesehenes Rot belegt“), `make doc-planning`
    und `make doc-targets` sind keine Gates, und kein Satz von `SPEC-040`
    hängt an ihnen. Ein Rot-Beleg bräuchte eine Konfiguration, die dieses
    Repo nicht führt; das wäre ein anderer Vorgang (§1).
  - Verifikation Hinweis 5 (67 Läufe nicht vollständig nachgefahren, Stichprobe
    deckungsgleich) → **hingenommen**: Hinweis, kein Befund; der Reviewer fuhr
    eine zweite Stichprobe, beide gleich.
- **Steering-Loop-Eintrag:** benannte Lücke, nicht verkörpert: der Bump-Vertrag
  (`harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs) sagt nicht, wie
  eine „nicht rein additive“ Änderung an einem **inaktiven** Modul zu führen
  ist — das Urteil „keine Reichweiten-Änderung eines Ziels“ ist wahr und
  trägt die latente Folge nicht; er sagt auch nicht, woraus die Liste der
  Einzel-Ziele einer Messung stammt (aus `d-check.mk`, nicht aus dem
  Changelog). Beides steht je 1× im Register (siehe unten) und wird bei der
  Schwelle dem Architect vorgelegt; mit diesem Slice ist weder eine Regel
  geschärft noch ein Sensor gebaut. Neuer Sensor: keiner. Nicht
  verkörpert, deshalb ohne `liegt in`.
- **Lese-Schritt (wellenlos, Baseline-Regelwerk `modul-06-roadmap.md`
  §Wann Arbeit eine Welle braucht):** kein Eintrag erreicht mit diesem Slice
  3× ohne Ausgang (gemessen: jedes Verzeichnis mit mindestens drei
  `evidence/`-Dateien trägt in `state.md` einen Ausgang; `adapter-fehler-ausgang`
  und `rollen-verdrahtung` führen ihn als *umgesetzt* bzw. *eingetreten*).
  `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (3×, verkörpert in
  `harness/targets/pin-stale.md`): diese Arbeit folgte der Regel — Urteil je
  Version vor dem Pin-Commit —, kein viertes Auftreten.
  `BEO-PGC/vertrag-doppelt-die-festlegung` (2×, offen): kein Vertrag und keine
  Festlegung nachgezogen, kein drittes Auftreten.
- **Trigger-Audit:** Carveout — 0 aktiv (`docs/plan/carveouts/` trägt nur
  `README.md`). Bootstrap-aware Gate — `make coverage-gate` steht auf der
  Endstufe 80 (`harness/mk/coverage.mk`, `THRESHOLD ?= 80`), nichts fällig.
  ADR — die Re-Evaluierungs-Trigger von `ADR-0163` (a) bis (c) sind nicht
  eingetreten (Verifikation §4, selbst gemessen); beide Pins sind aktuell
  (`make pin-stale-acheck`, `make pin-stale-dcheck` ohne `DRIFT`, Verifikation
  §2), die Pin-Zeilen P6/P7 von [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
  bleiben Stand-Angaben ihres Entscheidungszeitpunkts (§1). Hard Rule — kein
  Auflösungs-Trigger in `AGENTS.md` eingetreten. Ergebnis: nichts fällig.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-PGC/modul-reviews-inaktiv-aktivierung-trifft-bestand/` neu angelegt,
    Beleg `evidence/slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md`
    (Review F-1, MEDIUM) — Zähler 1×, offen.
  - `BEO-PGC/bump-messliste-unvollstaendig-gegen-zielliste/` neu angelegt,
    Beleg `evidence/slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md`
    (Review F-4, INFO) — Zähler 1×, offen.
  - Unter dem Deckel, ohne Datei: F-2
    (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, 27×) und F-3
    (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, 30×), je vor dem
    Merge gefunden, ≤ LOW bzw. INFO, bekannter Träger-Typ.
  - Ohne Register-Eintrag: F-5 (hingenommen, oben) — vom Plan selbst benannte
    Grenze, kein Mangel.
  - Die in §8 gesichteten Einträge `vertrag-doppelt-die-festlegung`,
    `pin-ohne-inventar-eintrag-driftet-unsichtbar` und
    `messwerkzeug-grenze-unbenannt-fail-open`: kein Auftreten.
- **Folge-Slices:** keiner. Die Übergabe-Zeile an
  `slice-spec-festlegungen-code-gates` (§6, „Werkzeug-Stand der a-check-Festlegung“)
  liegt dort seit `6884fb6b`; die Datei ist in `open/`.
- **Risiken aus §6:** fünf Risiken, je ein Ausgang:
  - „Ein Gate-Lauf bewegt sich über eine ‚nicht rein additive‘ Änderung“ →
    *weiter offen*: Register `BEO-PGC/modul-reviews-inaktiv-aktivierung-trifft-bestand`
    (132 Befunde bei Aktivierung von `reviews`; heute kein Ziel betroffen).
  - „`SPEC-040` beschreibt v0.86.1 nicht mehr“ → *entfallen*: 67 Vergleichsläufe
    und der volle `make docs-check`-Lauf ohne Unterschied, Verifikation 6
    Mutationen und 12 Läufe über `make` gleich; `spec/pflichtenheft.md` blieb
    unverändert.
  - „Trigger (a) von `ADR-0163` tritt ein“ → *entfallen*: `--print-config`
    beider Stände ändert am Block `matrix` keine Zeile (Plan §3, Verifikation §4).
  - „Der Index-Digest bricht den Freshness-Vergleich“ → *entfallen*: beide
    `pin-stale`-Läufe Exit 0 ohne `DRIFT`/`UNBESTIMMT` (Plan §3, Verifikation §2).
  - „Die a-check-Festlegung beschreibt einen alten Stand“ → *entfallen*: dieser
    Slice schließt vor `slice-spec-festlegungen-code-gates`, der Pin steht
    auf v0.23.2; die Übergabe-Zeile steht in dessen §6.
- **Archivierung:** keine — die Praxis wellenloser Slice-Closures in diesem
  Repo legt kein `done/slice-<Kennung>-archiv.zip` an; der Plan bleibt flach in
  `done/`. Der Reichweiten-Abschnitt (§3) und die Zahl 132 bleiben damit lesbar;
  die Zahl steht zusätzlich im Register.
- **Drei Paarungen:** nach dem `git mv` (`2deb6c6a`)
  gemessen — (a) *Anker:* der Steering-Loop-Eintrag trägt kein Feld `liegt in`
  (mit diesem Slice ist nichts verkörpert), kein Gegenstand der Paarung.
  (b) *Folge-Slice:* kein Folge-Slice genannt; die Übergabe-Adresse
  `slice-spec-festlegungen-code-gates` existiert in `open/` (`ls`) und trägt die
  Übergabe-Zeile in §6. (c) *Register:* die genannten Verzeichnisse existieren,
  jedes mit nicht leerem `evidence/` (`ls evidence | wc -l`):
  `modul-reviews-inaktiv-aktivierung-trifft-bestand` 1,
  `bump-messliste-unvollstaendig-gegen-zielliste` 1,
  `nachzug-laesst-ueberholten-text-stehen` 27,
  `zahl-in-traeger-driftet-gegen-die-messung` 30,
  `gate-scope-erweiterung-ohne-adr-traeger` 3,
  `vertrag-doppelt-die-festlegung` 2,
  `pin-ohne-inventar-eintrag-driftet-unsichtbar` 2,
  `messwerkzeug-grenze-unbenannt-fail-open` 2. Ergebnis: getragen. Durch den
  Move brach kein Verweis (`git grep -n 'in-progress/slice-gate-werkzeuge'`
  außerhalb von `docs/reviews/`: kein Treffer). Der Ruhe-Marker der Roadmap
  steht wieder, Wortlaut gleich `cffa45be`; `in-progress/` trägt nur
  `roadmap.md`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `a-check.mk`,
`d-check.mk` (Pins), `spec/pflichtenheft.md` §7 (`SPEC-040`) und gegebenenfalls
`harness/sensors/` (Verträge). Die Modus-Deklaration in
`harness/conventions.md` führt nur die Default-Sub-Area `*` (Kürzel `PGC`,
Greenfield); alle Pfade fallen unter sie. Keine Ausdifferenzierung nötig: der
Slice bewegt in allen Pfaden dieselbe Eigenschaft (den Stand der beiden
Gate-Werkzeuge und was er an ihren Festlegungen bewegt).

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-PGC/` durchgegangen (160 Verzeichnisse,
gemessen mit `ls | wc -l` am 2026-10-10, Stand `da913c16`), gefiltert nach
Pin, Bump, Werkzeug, Gate, Spec, Festlegung, Digest und Messung. Treffer mit
Bezug, Zähler als Zahl der `evidence/`-Dateien (gemessen mit `ls`):

- `gate-scope-erweiterung-ohne-adr-traeger` — 3×, *verkörpert* in
  `harness/targets/pin-stale.md` §Bump eines Gate-Werkzeugs. Diese Regel
  trägt Liefer-Punkt 1; ein weiteres Auftreten wäre ein Pin-Commit ohne
  Reichweiten-Urteil.
- `vertrag-doppelt-die-festlegung` — **2×**, offen (Belege
  `slice-spec-festlegungen-harness-werkzeuge`, `slice-spec-festlegungen-doku-gates`).
  Berührt, falls Liefer-Punkt 3 `SPEC-040` nachzieht und
  `harness/sensors/docs-check.md` dieselbe Aussage selbst trägt: ein drittes
  Auftreten machte den Eintrag zur Lücke. Deshalb steht der Vertrag in der
  §3-Tabelle als Prüfstelle und im Suchlauf (Zeile 7).
- `mechanismus-erklaerung-ohne-werkzeugbeleg` — 4×, *verkörpert*, für
  Festlegungen über ein Werkzeug in `.claude/commands/plan-welle.md` Schritt 6,
  Punkt (c) „Messung am Werkzeug“ (Commit `da913c16`); trägt Liefer-Punkt 3
  (c).
- `pin-ohne-inventar-eintrag-driftet-unsichtbar` — **2×**, offen. Nicht
  berührt: beide Pins stehen im Inventar (P6, P7 in
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)); kein drittes
  Auftreten.
- `implementierung-weicht-von-adr-wortlaut-ab` — 4×, *verkörpert* in
  [`AGENTS.md`](../../../../AGENTS.md) §3.5. Bezug: weicht das Werkzeug von
  einer Quell-ADR von `SPEC-040` ab, ist das eine Frage an den Architect
  (Trigger (c) von `ADR-0163`), keine Regel im Plan.
- `messwerkzeug-grenze-unbenannt-fail-open` (2×, offen) — kein Bezug: kein
  nachbildendes Messwerkzeug wird geändert.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF.
