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

**Verantwortlich:** —
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

- [ ] **Reichweite je übersprungener Version, vor jedem Pin-Commit
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
- [ ] **Pins auf v0.23.2 und v0.86.1 (Liefer-Punkt 2).** `a-check.mk` trägt
      `A_CHECK_IMAGE ?= ghcr.io/pt9912/a-check@sha256:2368f7b3a84f1dc5d075edccfe2201e19947d12fcbc8eaf4df84ef162d94f422`,
      `d-check.mk` `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.86.1` und
      `DCHECK_DIGEST ?= sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`
      — je Werkzeug ein Commit, der sonst nichts ändert; Kommentare in den
      beiden Dateien, die einen Stand nennen, sind mit dem Suchlauf §3
      geprüft. *Zu belegen durch:* `docker buildx imagetools inspect` je Tag
      (gedruckter Digest gleich dem Release-Digest), `make pin-stale-acheck`
      und `make pin-stale-dcheck` ohne `DRIFT` (gedruckte Zeilen; braucht
      Netz), `git show --stat` je Pin-Commit (eine Datei).
- [ ] **`SPEC-040` gegen v0.86.1 gehalten (Liefer-Punkt 3).** Für jeden Satz
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
- [ ] `make gates` grün, Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — die
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
  am Stand des Slice und die Läufe aus Liefer-Punkt 1. Erwarteter Ausgang:
  entfallen, mit diesem Grund.
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

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Gegenstand:** <übernommen von `slice-<Kennung>` | entfallen: <Grund>>
  *(nur beim Ausgang ohne Arbeit; sonst Zeile löschen)*
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<KUERZEL>/<slug>` (<slice-kennung-a>, <slice-kennung-b>, <slice-kennung-c> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-<Kennung>.md` | `evidence/slice-<Kennung>.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-<Kennung> (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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
