# Slice abgeleitete-dokumente-vorlagen-nachzug: Harness-Einstieg, Konventionen, ADR-Index und Carveout-Ablage entsprechen ihren Vorlagen v6.14.0

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle
braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P8,
Baseline-Bump). Keine `LH-*`-Anforderung ist berührt: der Slice ändert
Harness-Dokumente, nicht das Produkt.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner, Closure von `slice-baseline-6-14-0-dokumente-nachziehen`).
**Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass.** Die Bestandsaufnahme von `slice-baseline-6-14-0-dokumente-nachziehen`
(§3, *Liefer-Punkt 2c*, gemessen am Stand `d79b7ebd`) nennt vier Dokumente mit
**Abweichung mit Beleg** gegen ihre Vorlage, alle älter als der Bump:
`harness/README.md` (Platzhalter `` `<make-target>` ``, `make <mover>`/`<messung>`/`<vorschau>`,
§Safety and scope boundaries zweimal `<…>`, §Leseordnung drei Platzhalter),
`harness/conventions.md` (`<Pfad oder URL>`, `<Pfade zu deinen …>`, MR-000
`**Datum:** <Datum>`, Musterzeilen in §Zusatzklassen und §Glossar),
`docs/plan/adr/README.md` (kein Abschnitt `## Konventionen`, keine Regel zum Feld
`**Schärft:**`) und `docs/plan/carveouts/` (kein `README.md`). Dazu kommt ein
Träger der Kennungsform nach `MR-002`, den die Closure jenes Slice nicht
nachziehen konnte: das `forbid-pattern` der zwei `structure`-Regeln gegen
Slice-Pfad-Links in `.d-check.yml` trifft nur `slice-[0-9]{3}`.

**Ziel:** Die vier Dokumente tragen keinen Vorlagen-Platzhalter mehr und die
Abschnitte ihrer Vorlage (oder eine benannte, repo-spezifische Abweichung mit
Grund), und die Wächter gegen Slice-Pfad-Links decken die Namens-Kennung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Inhalt der Sensor- und Target-Dateien unter `harness/sensors/` und
  `harness/targets/`** (Abschnitte `## Fassung im Gate-Index`, veraltete Aussagen,
  `baseline-verify.md` Grenze 3) — trägt `slice-harness-targets-inhalt-bereinigen`.
  Beide Slices berühren `harness/README.md`: dieser Slice die Vorlagen-Platzhalter
  (§Sensors-Musterzeile, Werkzeug-Musterzeilen, §Safety, §Leseordnung), jener nur
  einen Index-Kurzsatz, der sich durch eine Berichtigung ändert.
- **`MR-002` auflösen.** Die Vorlage `conventions.template.md` v6.14.0 nennt im
  ID-Schema `slice-<Kennung>`; ob damit der Auflösungs-Trigger von `MR-002`
  eingetreten ist, entscheidet die Änderung an `harness/conventions.md`
  §Adaptions-Block — ein anderer Vorgang, Architect-Zug.
- **Die Regex der E2E-Abdeckungstabelle** (`abdeckungKennungMuster` in
  `test/integration/integration_test.go`, `slice-\d{3}|welle-\d{1,2}`) — Bestand
  bleibt: kein `TestE2E*`-Kommentar trägt eine Namens-Kennung (gemessen bei der
  Closure des Auslöser-Slice, `git grep` über `test/integration`, 0 Treffer), und
  der Kommentar beschreibt die Regex richtig; eine Änderung wäre Produkt-Testcode.
- **Records und `Accepted`-ADRs** mit `slice-<NNN>`/`welle-NN` — eingefroren
  (`AGENTS.md` §3.5).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [x] **Harness-Einstieg und Konventionen (Liefer-Punkt 1).** `harness/README.md`
      und `harness/conventions.md` tragen keinen Platzhalter der Vorlage; jede
      Musterzeile ist ausgefüllt oder gestrichen, und eine bewusst stehende Zeile
      trägt ihren Grund. Zu belegen durch beide Platzhalter-Formen der Regel in
      `harness/targets/pin-stale.md` (*Bump-Ablauf*, Schritt 2) mit 0 Treffern
      bzw. jedem Treffer mit Grund.
- [x] **ADR-Index und Carveout-Ablage (Liefer-Punkt 2).** `docs/plan/adr/README.md`
      trägt den Abschnitt `## Konventionen` der Vorlage samt der Regel zu
      `**Schärft:**`; die Spalten `Datum`/`Datei` statt `Bezug` stehen als
      repo-spezifische Abweichung mit Grund oder sind angeglichen.
      `docs/plan/carveouts/README.md` existiert aus `carveouts/README.template.md`.
      Zu belegen durch den Überschriften-`diff` gegen die Vorlage (Exit 0).
- [x] **Slice-Pfad-Wächter (Liefer-Punkt 3).** Das `forbid-pattern` der zwei
      `structure`-Regeln in `.d-check.yml` trifft auch einen Link auf
      `…/in-progress/slice-<Name>.md`; Grenze (5) im Kommentar und in
      `harness/sensors/docs-check.md` ist gestrichen. Vor der Umsetzung
      **Übergabe an den Architect**: ob die Erweiterung eine ADR braucht
      (`AGENTS.md` §3.6, `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`).
      Zu belegen durch eine Mutation (ein solcher Link in einer Kopie) mit Befund.
- [x] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0. —
      Beleg: Verifikation
      [`verify-slice-abgeleitete-dokumente-vorlagen-nachzug.md`](../../../reviews/verify-slice-abgeleitete-dokumente-vorlagen-nachzug.md)
      (`make docs-check` Exit 0 am Stand `e54b0322`); `make gates` Exit 0 nach
      dem Commit des Berichts (`1466964e`, übernommen aus der Rückmeldung des
      Verifiers); `make gates` nach dem letzten Commit der Closure, Exit im
      Bericht der Sitzung.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: Träger, die die geänderten Abschnitte zitieren, sind
      nachgezogen (Suchlauf, `AGENTS.md` §3.13).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/README.md`, `harness/conventions.md` | update | Liefer-Punkt 1 |
| `docs/plan/adr/README.md`, `docs/plan/carveouts/README.md` | update / neu | Liefer-Punkt 2 |
| `.d-check.yml` (zwei `structure`-Regeln), `harness/sensors/docs-check.md` Grenze | update | Liefer-Punkt 3 |
| `docs/plan/planning/README.md` | update (Plan-Nachzug) | Liefer-Punkt 2: der Kommentar `d-check:ignore` am Pfad `docs/plan/carveouts/done/` fällt — er ist wirkungslos (Messung unten) und die einzige Abweichung der Zeile von `planning/README.template.md` Zeile 47 |
| `harness/README.md` (zwei Bindung-Zellen der SDK-Integrationstests) | update (Fixrunde F-1) | nackte Slice-Kennung vor `· seit slice-…` gestrichen; die Zelle trägt die Herkunft in der deklarierten Form |
| `harness/sensors/docs-check.md` §Bindung, `observations/BEO-PGC/slice-pfad-als-link-in-berichten/state.md` | update (Fixrunde F-7) | Zähler „3×“ entfernt (der Zähler ist abgeleitet); Reparatur-Pfad `slice-NNN` → `slice-<Kennung>` wie im `hint` |

### Ergebnis je Liefer-Punkt (Implementer, Stand vor dem Commit, Parent `3f51d3ed`)

**Liefer-Punkt 1 — `harness/README.md`.** Je Platzhalter, Ausgang und Quelle:

| Stelle | Ausgang | Quelle / Grund |
|---|---|---|
| Gate-Tabelle, Zeile `` `<make-target>` `` „volle Closure“ | gestrichen | Das Makefile führt kein Closure-Ziel (`grep -nE '^(fullbuild\|ci)[: ]' Makefile harness/mk/*.mk d-check.mk`: 0 Treffer); der Image-Digest ist ein lokaler Lauf-Beleg von `make image`, das in der Werkzeug-Tabelle steht ([`ADR-0103`](../../adr/0103-image-hash-lokal-statt-committet.md)) |
| Werkzeug-Tabelle, Musterzeilen `make <mover>`, `make <messung>`, `make <vorschau>` | gestrichen | Formbeispiele der Vorlage; die Werkzeug-Tabelle führt die realen Ziele dieser Arten |
| §Safety and scope boundaries, zweimal `<…>` | ausgefüllt, fünf Punkte | je Punkt die zitierte Quelle in der Zeile selbst: `spec/lastenheft.md` §MVP-Schnitt (Reihenfolge, Rollback, Neustart), `LH-FA-CFG-006`, `LH-QA-SEC-002` mit `ADR-0047`, die drei Kennungs-Gates aus §Sensors (Punkt 4 auf deren Gegenstand eingeengt, Fixrunde F-5), `AGENTS.md` §3.1 |
| §Leseordnung, drei Platzhalter | ausgefüllt, vier Zeiger | Die **Zeiger** stammen aus `AGENTS.md` §2 (Rang 7 `README.md`, Rang 1 `spec/lastenheft.md`) und aus den Beispielen der Vorlage (`AGENTS.md` Hard Rules, `harness/conventions.md` bei Bedarf). Die **Reihenfolge** ist eigene Wahl und folgt nicht den Rängen: ein neuer Mensch braucht zuerst den Überblick über den Gegenstand (`README.md`), dann die Regeln, die jede Änderung bindet (`AGENTS.md` §3), dann den Vertrag im Detail (`spec/lastenheft.md`); die Konventionen braucht erst, wer eine Struktur- oder Kennungsfrage hat |

**Liefer-Punkt 1 — `harness/conventions.md`.**

| Stelle | Ausgang | Quelle / Grund |
|---|---|---|
| `**Extern (Lehrmaterial):** <Pfad oder URL>` | ausgefüllt | `https://github.com/pt9912/ai-harness-course` — das Repo des vendored Release-Assets (Zeile darunter) und der Abfrage in `tools/harness/pin-stale-baseline.sh` Zeile 26 |
| `**In-Repo (verkörperte Form):** <Pfade zu …>` | ausgefüllt | die Liste der abgeleiteten Dokumente in `harness/targets/pin-stale.md` §Bump-Ablauf Schritt 2, mit `docs/plan/carveouts/README.md` statt des Verzeichnisses |
| MR-000 `**Datum:** <Datum>` | `2026-09-09` | `git log --follow --format='%h %ad' --date=short -- harness/conventions.md`: ältester Commit `40c8c431` 2026-09-09 (legt die Datei an) |
| §Aufgelöste Adaptionen, Musterzeile `\<NNN\>` | `— \| —` | `harness/conventions/done/` existiert nicht (`ls`: nicht gefunden) |
| §Zusatzklassen, Musterzeile | sechs Klassen (Fixrunde F-1) | die Bindung-Spalte beider Tabellen in `harness/README.md` §Sensors, ausgezählt (Befehl und Zahlen unter *Fixrunde*, F-1): LH (2 Zellen), Vertragsdatei `sensors/…`/`targets/…` (66), Herkunfts-Anker `seit slice-…`/`seit welle-…` (40), `BEO-PGC/…` (1), `MR-*` (1), Link auf `AGENTS.md` §n (1); je Klasse ein Beispiel aus einer realen Zeile. Die zwei nackten Slice-Kennungen ohne `seit` sind gestrichen: beide Zellen trugen dieselbe Kennung schon als `· seit slice-…` (Verifikation V-2) |
| §Glossar, Musterzeile | gestrichen, Satz mit Grund | Produkt-Begriffe führt `spec/lastenheft.md` §6 Glossar |

**Platzhalter-Prüfung** (beide Formen der Regel in `harness/targets/pin-stale.md`
§Bump-Ablauf Schritt 2, Arbeitsbaum): Form B (ältere Vorlagen)
`harness/README.md` 5, `harness/conventions.md` 7, `docs/plan/adr/README.md` 0,
`docs/plan/carveouts/README.md` 2 Treffer; Form A (Platzhalter der Vorlage)
11 / 15 / 2 / 3 Treffer. Jeder verbleibende Treffer ist **Notation**, kein
Platzhalter, und bleibt: die Kennungsformen `CO-<NNN>`, `MR-<NNN>`,
`ADR-<NNNN>`, `SPEC-<NNN>`, `<PREFIX>-FA-*`, `slice-<Kennung>`, `LH-FA-*.<Buchstabe>`,
`PCF-<S><NNNN>` (Schema-Angaben in Inline-Code), `.harness/baseline/<tag>/`
(der Tag wechselt mit jedem Bump), `harness/sensors/<target>.md` und
`<LH-*>` (Wortlaut der Vorlage), die HTML-Kommentare der Vorlage und die
Anker `<a id="mr-…"></a>` der MR-Zeilen (Form A trifft `</a>`).

**Liefer-Punkt 2.** `docs/plan/adr/README.md` trägt `## Konventionen` mit den
vier Punkten der Vorlage (`<PREFIX>` als `LH`) und einem fünften zu den
Spalten: `Datum`/`Datei` statt `Bezug` bleiben **repo-spezifisch** — den Bezug
trägt das `**Schärft:**`-Feld jeder ADR (`git grep -l -E '^\*\*Schärft:\*\*' -- 'docs/plan/adr/0*.md' | wc -l`: 155 von 155 Dateien),
und die `structure`-Regel in `.d-check.yml` hält `ID` auf acht Zeichen, also
ohne Link. `docs/plan/carveouts/README.md` per `cp` aus
`carveouts/README.template.md`, ausgefüllt: keine aktiven und keine
aufgelösten Carveouts (`git log --all -- 'docs/plan/carveouts/CO-*'`: leer),
der Trigger-Audit ohne Welle aus `modul-06-roadmap.md` ergänzt.
Gliederung, `<Projektname>` der Vorlage in einer Scratchpad-Kopie durch
`PG Change Feed` ersetzt:
`diff <(grep -E '^#{1,4} ' <Vorlage>) <(grep -E '^#{1,4} ' <Datei>)` —
Exit 0 für alle vier Dokumente (`harness/README.md`, `harness/conventions.md`,
`docs/plan/adr/README.md`, `docs/plan/carveouts/README.md`).
**`d-check:ignore` in `docs/plan/planning/README.md`:** wirkungslos, gemessen
mit d-check (Pin aus `d-check.mk`) an einer Kopie im Scratchpad ohne den
Kommentar: Modulliste der `.d-check.yml` mit und ohne `carveouts/README.md`
Exit 0, 0 Befunde; mit `--enable codepaths` 132 Befunde mit und ohne Kommentar,
keiner in `planning/README.md`. Der Kommentar fällt.

**Liefer-Punkt 3** (Verdikt `architect-verdict-slice-pfad-waechter-namensform`,
Variante (a)): beide `forbid-pattern` lauten
`'\]\([^)]*(open|next|in-progress)/slice-[a-z0-9]'`; Grenze (5) im Kommentar
von `.d-check.yml` und der Satz zur Nummern-Form in
`harness/sensors/docs-check.md` §Grenze Punkt 6 sind gestrichen, `· seit
slice-075` und die Grenzen (1) bis (4) bleiben. Mutation, d-check
(`ghcr.io/pt9912/d-check@sha256:b4b8756b…`) gegen eine Kopie des Arbeitsbaums
im Scratchpad:

| Lauf | Muster | Inhalt | Exit | `section-forbidden` |
|---|---|---|---|---|
| Bestand | weit | Arbeitsbaum | 0 | 0 |
| Mutation | weit | je ein Link `[Plan](…/in-progress/slice-abgeleitete-dokumente-vorlagen-nachzug.md)` in `docs/reviews/architect-verdict-slice-pfad-waechter-namensform.md` und in `observations/BEO-PGC/a-check-null-abdeckung/observation.md` | 1 | 2 (je Datei einer) |
| Mutation | alt (`slice-[0-9]{3}`, `.d-check.yml` aus `3f51d3ed`) | dieselben zwei Links | 0 | 0 |
| Nummern-Form | weit | ein Link `…/next/slice-105-x.md` im Review-Bericht | 1 | 1 (dazu `target-missing`) |

Stellen: eine Instanz je Dateiklasse; Instanz: d-check-Lauf über die Kopie,
nicht ein Test.

**Suchlauf (`AGENTS.md` §3.13).** Bewegte Eigenschaften: Reichweite des
Slice-Pfad-Musters (M1, M2), Existenz der Carveout-Ablage (M3), der
`d-check:ignore`-Kommentar (M4), die Platzhalter der vier Dokumente (M5).
Suchraum: ganzer Baum ohne `docs/reviews/**`, Records unter `done/` und
`.harness/baseline/**`.

```suchlauf
3f51d3ed 6 -E 'slice-\[0-9\]\{3\}|Nummern-Form|Namens-Kennung (ist )?nicht gedeckt|Namens-Form' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
diff 0 -E 'slice-\[0-9\]\{3\}|Nummern-Form|Namens-Kennung (ist )?nicht gedeckt|Namens-Form' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
3f51d3ed 5 -E 'Grenze \(5\)|Grenze 5|Punkt 6' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
diff 5 -E 'Grenze \(5\)|Grenze 5|Punkt 6' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
3f51d3ed 0 -E 'carveouts/README|carveouts/\.gitkeep|Carveout-(Index|Ablage)' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
diff 1 -E 'carveouts/README|carveouts/\.gitkeep|Carveout-(Index|Ablage)' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
3f51d3ed 1 -F 'd-check:ignore (done/' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
diff 0 -F 'd-check:ignore (done/' -- . :(exclude)docs/reviews :(exclude,glob)**/done/** :(exclude).harness/baseline
3f51d3ed 15 -E '<make-target>|<mover>|<messung>|<vorschau>|<…>|<zuerst|<dann|<bei Bedarf|<Pfad oder URL>|<Pfade zu|<Datum>|<repo-spezifischer|<z\. B\.|MR-\\<NNN\\>' -- harness/README.md harness/conventions.md docs/plan/adr/README.md docs/plan/carveouts/README.md
diff 0 -E '<make-target>|<mover>|<messung>|<vorschau>|<…>|<zuerst|<dann|<bei Bedarf|<Pfad oder URL>|<Pfade zu|<Datum>|<repo-spezifischer|<z\. B\.|MR-\\<NNN\\>' -- harness/README.md harness/conventions.md docs/plan/adr/README.md docs/plan/carveouts/README.md
```

Gefunden: M1 am Parent die sechs Zeilen der zwei Träger (`.d-check.yml`
Grenze (5) und beide Muster, `harness/sensors/docs-check.md` Punkt 6), am
Arbeitsbaum keine. M3 am Arbeitsbaum die neue Zeile in `harness/conventions.md`.
**Nicht gefunden:** kein weiterer Träger der Nummern-Grenze — die Fundstellen
von `slice-pfad-als-link-in-berichten` (`.claude/commands/implement-slice.md`
Schritt 25, `state.md` und `evidence/` des Eintrags,
`harness/sensors/docs-check.md` §Vertrag und §Bindung; Suche
`git grep -n -i -E 'slice-pfad|wandernde[n]? Slice|Lifecycle-Verzeichnis'`)
nennen die Dateiklassen, nicht die Kennungsform. M2 trifft an beiden Ständen
fünf fremde Zeilen (`ADR-0053`, zwei `evidence/`-Dateien,
`internal/bootstrap/otlp_internal_test.go`) — andere Gegenstände. Kein Träger
beschreibt die Carveout-Ablage als leer oder den ADR-Index als ohne
Konventionen; `harness/targets/pin-stale.md` nennt `docs/plan/carveouts/` als
Verzeichnis und bleibt richtig. Den hier zunächst nur gemeldeten Zähler „3×“
in `harness/sensors/docs-check.md` §Bindung zieht die Fixrunde nach (F-7).

### Fixrunde (Review `review-slice-abgeleitete-dokumente-vorlagen-nachzug`, Commit `a2fb6cf0`)

- **F-1 (HIGH), angenommen.** Auszählung der Bindung-Spalte, je Stand (Parent
  `3f51d3ed` mit den vier gestrichenen Musterzeilen, Arbeitsbaum ohne sie);
  Zellen = letzte Spalte jeder Datenzeile zwischen `## Sensors` und
  `## Traceability`, je Muster die Zahl der Zellen mit Treffer:

  ```text
  awk '/^## Sensors/,/^## Traceability/' harness/README.md | grep -E '^\| ' \
    | grep -vE '^\| (Target|---)' | awk -F' \\| ' '{print $NF}' > <Scratchpad>/b.txt
  grep -cE '<Muster>' <Scratchpad>/b.txt
  ```

  Gedruckt (`3f51d3ed` / Arbeitsbaum): Zellen 70 / 66 · `LH-(FA|QA)-` 2 / 2 ·
  `AGENTS\.md` 1 / 1 · `BEO-PGC/` 1 / 1 · `MR-[0-9]` 1 / 1 ·
  `seit (slice|welle)-` 40 / 40 · `(sensors|targets)/` 66 / 66 ·
  nackte Kennung `(^|, )slice-[a-z0-9-]+ ·` 2 / 0. Die übrigen Formen der
  Spalte sind die kanonischen Klassen der Vorlage (`ADR-*`, `kein Gate`,
  Reproduzierbarkeit über Modul 14 in der Zeile `make image-stale`).
  `harness/conventions.md` §Zusatzklassen führt jetzt sechs Klassen, mit
  „Hard-Rule-Bindung“ (`AGENTS.md` §3.7 bei `make kommentar-kennungen`); die
  nackte Kennung ist nicht deklariert, sondern in den Zellen von
  `make test-sdk-kotlin-integration` und `make test-sdk-csharp-integration`
  gestrichen — beide trugen dieselbe Kennung schon als `· seit slice-…`. Die
  Zellenlängen-Regel misst `Vertrag`/`Tut was`, nicht `Bindung`;
  `make docs-check` bleibt bei 0 Befunden.
- **F-2 (MEDIUM), angenommen.** Die Regel steht nur noch in §Konventionen
  des ADR-Index, mit der Ausnahme der Zitat-Korrektur nach `ADR-0073` und dem
  Verweis auf `AGENTS.md` §3.5; der Kopf-Absatz verweist auf §Konventionen.
- **F-3 (LOW), angenommen.** `docs/plan/carveouts/README.md` stellt die Achse
  voran: das Repo führt Wellen, die Welle-Closure liest auch die wellenlosen
  Slices seit der letzten Welle (Baseline-Regelwerk `modul-06-roadmap.md`
  §Wann Arbeit eine Welle braucht). Ohne offene Welle läuft der Trigger-Audit
  zusätzlich bei der Closure eines wellenlosen Slice — Beleg im Bestand:
  `ADR-0070` nennt als Autor-Anlass den „Trigger-Audit der Slice-Closure“
  (`git grep -n 'Trigger-Audit der Slice-Closure' -- docs/plan/adr`), und die
  DoD dieses Slice lässt die Slice-Closure ihre Paarungen tragen, „solange die
  Roadmap unter *Offene Wellen* keine Welle führt“.
- **F-4 (LOW), angenommen.** Die Quelle der Leseordnung in §3 trennt Zeiger
  (`AGENTS.md` §2, Vorlage) und Reihenfolge (eigene Wahl, mit Grund).
- **F-5 (INFO), angenommen.** Safety-Punkt 4 nennt die drei geprüften
  Nutzerdokumente und schließt die erzeugten `*-abdeckung.md` aus
  (`harness/sensors/handbuch-public-doc-check.md`).
- **F-6 (INFO), angenommen.** Das Argument „zweite Quelle“ ist gestrichen;
  tragend bleiben die `structure`-Regel (`ID` acht Zeichen, Link in `Datei`)
  und das `**Schärft:**`-Feld in jeder ADR-Datei (155 von 155, §3).
- **F-7 (INFO), angenommen.** `harness/sensors/docs-check.md` §Bindung
  nennt den Eintrag zählerfrei (`· seit slice-075` bleibt); `state.md` des
  Eintrags gibt den Reparatur-Pfad als `slice-<Kennung>` wieder.
  `git grep -n -F 'slice-NNN' -- docs/plan/planning/observations/BEO-PGC/slice-pfad-als-link-in-berichten`:
  `3f51d3ed` 3 Zeilen, danach 2 — die zwei in `observation.md` bleiben, die
  Identitäts-Datei ist ab Anlage unveränderlich (Baseline-Regelwerk
  `modul-06-roadmap.md` §Das Beobachtungs-Register).
- **F-8 (INFO), keine Änderung.** Befund-Liste, Hinweis auf die
  Baseline-Vorlage: MR-000 führt `BEO-<NNN>`, wie
  `templates/harness/conventions.template.md` Zeile 107 der Baseline
  `v6.14.0`; Register und Regelwerk (`modul-06-roadmap.md`) führen
  `BEO-<KUERZEL>/<slug>`. Die Abweichung kommt aus der Vorlage — Sache des
  Freshness-Audits (`AGENTS.md` §1) bzw. einer Meldung an die Baseline.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1); läuft `slice-harness-targets-inhalt-bereinigen` zuerst, liest
dieser Slice dessen Stand von `harness/README.md` mit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Architect verlangt
  für Liefer-Punkt 3 eine ADR — dann geht Liefer-Punkt 3 in einen eigenen Slice.
- `in-progress` → `open` (blockiert): eine Leseordnung oder §Safety-Grenze lässt
  sich nicht aus dem Bestand ableiten — Frage an den Auftraggeber.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, und die Closure-Notiz in §7 trägt einen Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ausgefüllte Leseordnung und §Safety erfinden Inhalt**, den kein Träger trägt.
  Gegenmaßnahme: jede Zeile zeigt auf eine bestehende Regel oder Datei; die
  Quelle je Zeile steht in §3 (Liefer-Punkt 1). —
  **Ausgang:** entfallen — jede Zeile der Leseordnung und jeder der fünf
  Safety-Punkte ist an ihrer Quelle nachgeprüft, vom Reviewer (Review, Proben 1
  und 2) und in frischem Kontext vom Verifier (Verifikation §2, Zeilen
  „§Safety trägt an der Quelle“ und „§Leseordnung“); kein Punkt ohne Quelle.
  Punkt 4 ist in der Fixrunde auf den Gegenstand des Gates eingeengt (F-5).
- **Das erweiterte Muster meldet Bestand.** Gemessen bei der Anlage: der weite
  Ausdruck trifft in den zwei Dateiklassen eine Zeile, in Inline-Code (für die
  Regel unsichtbar). Nachgemessen vom Implementer: d-check mit dem weiten Muster
  über den Arbeitsbaum Exit 0, 0 `section-forbidden` (§3, Liefer-Punkt 3). —
  **Ausgang:** entfallen — `make docs-check` meldet am Stand `e54b0322` 0
  Befunde, ebenso der Bestandslauf des Verifiers am Klon (Verifikation §1 und
  §4); die Negativkontrollen dort zeigen, dass das Muster `roadmap.md` in
  `in-progress/`, Links nach `done/` und Inline-Code nicht trifft. Einen neuen
  Link auf einen wandernden Plan meldet die Regel, wie sie soll.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln. Ging der
Gegenstand an einen anderen Slice oder entfiel er, trägt diese Sektion die Zeile
`Gegenstand:` mit Kennung oder Grund.

- **Geliefert:** `harness/README.md` und `harness/conventions.md` tragen keinen
  Platzhalter ihrer Vorlage mehr (Liefer-Punkt 1; jeder verbleibende Treffer
  beider Platzhalter-Formen ist Notation, §3), `docs/plan/adr/README.md` trägt
  `## Konventionen` samt der Regel zu `**Schärft:**` und der Spalten-Abweichung
  mit Grund, `docs/plan/carveouts/README.md` existiert aus der Vorlage
  (Liefer-Punkt 2; Überschriften-`diff` Exit 0 für alle vier Dokumente), und
  beide `structure`-Regeln gegen Slice-Pfad-Links treffen die Namens-Kennung
  (Liefer-Punkt 3, Architect-Verdikt `architect-verdict-slice-pfad-waechter-namensform`
  Variante (a), keine ADR). Review
  `docs/reviews/review-slice-abgeleitete-dokumente-vorlagen-nachzug.md`
  (`a2fb6cf0`; F-1 HIGH, F-2 MEDIUM, F-3/F-4 LOW, F-5 bis F-8 INFO), Fixrunde
  `e54b0322` (F-1 bis F-7 angenommen, F-8 ohne Änderung), Verifikation
  `docs/reviews/verify-slice-abgeleitete-dokumente-vorlagen-nachzug.md`
  (`1466964e`; DoD bestätigt, V-1/V-2 INFO, H-1 an die Closure).
- **Was hat funktioniert:** Die Übergabe an den Architect *vor* der Umsetzung
  von Liefer-Punkt 3 (DoD) hat die Frage `AGENTS.md` §3.6 entschieden, bevor ein
  Diff entstand; das Verdikt nannte die Grenze zur Klasse
  `gate-scope-erweiterung-ohne-adr-traeger` mit Beleg. Die Mutation des
  Wächters lief an drei Stellen mit je anderen Formen (Implementer: Link ohne
  Anker in beiden Dateiklassen; Reviewer: `open/`, `next/`; Verifier:
  Anker-Suffix, Name mit Ziffer-Präfix, dazu Negativkontrollen) — jede rot mit
  neuem und grün mit altem Muster. Jede Zeile von §Safety und §Leseordnung ist
  von zwei Rollen an ihrer Quelle nachgeprüft.
- **Was ging anders als geplant:** (1) **Eine Auszählung stand als gemessen
  da, ohne gemessen zu sein** (Review F-1, HIGH): §3 leitete die Zusatzklassen
  aus „der Bindung-Spalte beider Tabellen …, ausgezählt“ ab, nannte aber weder
  Befehl noch gedruckte Zahlen; die Nachzählung des Reviewers ergab `LH-*` in
  2 statt 3 Zellen und zwei Bindungsformen ohne Klasse (Hard-Rule-Link,
  nackte Slice-Kennung). Die Fixrunde trägt Befehl und Zahlen je Stand
  (`3f51d3ed` / Arbeitsbaum), deklariert die Hard-Rule-Bindung und streicht
  die nackten Kennungen; der Verifier zählte mit eigener Methode dieselben
  Zahlen. (2) **Der Nachzug widersprach dem Nachbarn im selben Träger**
  (F-2, MEDIUM): `## Konventionen` des ADR-Index nannte die Immutabilität ohne
  die Ausnahme, die der Kopf-Absatz derselben Datei trägt; die Regel steht
  jetzt einmal, mit der Zitat-Korrektur nach `ADR-0073`. (3) Der Plan-Nachzug
  der Fixrunde ließ in §3 „angeglichen“ neben dem richtigen „gestrichen“ des
  Fixrunde-Absatzes stehen (V-2, INFO); bei der Closure berichtigt. (4) Zwei
  Quellenangaben trugen nur einen Teil ihrer Aussage (F-3 Achse der
  Baseline-Regel, F-4 Reihenfolge der Leseordnung); in der Fixrunde behoben.
- **Re-Review:** nicht nötig — die Fixrunde `e54b0322` änderte nur Text (Plan,
  `harness/README.md`, `harness/conventions.md`, ADR-Index, Carveout-README,
  `harness/sensors/docs-check.md`, `state.md` eines Register-Eintrags), keine
  Logik und kein Muster, und der Verifier hat F-1 bis F-8 in frischem Kontext
  selbst nachgefahren (Verifikation §3, Auszählung mit eigener Methode)
  (`.claude/commands/implement-slice.md` Schritt 21).
- **Review- und Verifikations-Pflichten in der Closure:** V-2 → §3 Zeile der
  Zusatzklassen sagt „gestrichen“. H-1 → Ausgänge in §6. V-1 → siehe
  *Befunde*.
- **Befunde und Folgearbeit:**
  - **V-1 (INFO), Freitext in der Bindung-Spalte:** neben den Klassen trägt die
    Spalte erläuternden Text („braucht Netz“ in 6 Zellen, „(ein Wächter
    verhindert eine Handlung, er prüft kein Ergebnis)“, „wie `make
    image-stale`“, „DB-Adapter-Coverage (…)“) — die Begründung zu `kein Gate`,
    keine Bindungsform. Bestand, nicht dieser Diff; die Vorlage entscheidet
    nicht, ob er in die Spalte gehört. Keine Folgearbeit: es gibt keine Regel,
    gegen die er verstieße. Kein Register-Beleg: die nächste Klasse
    `BEO-PGC/bindung-spalte-uneinheitlich-tief` (1×) betrifft die Zitiertiefe
    der Spalte, nicht Begleittext — ein Beleg dort spaltete die Klasse in zwei
    Formen; *benannt, nicht gezählt*.
  - **F-8 (INFO), Baseline-Vorlage führt `BEO-<NNN>`:** MR-000 in
    `harness/conventions.md` folgt
    `templates/harness/conventions.template.md` Zeile 107 der Baseline
    `v6.14.0`; das Regelwerk (`modul-06-roadmap.md` §Das Beobachtungs-Register)
    und das Register führen `BEO-<KUERZEL>/<slug>`. Die Abweichung liegt in der
    Vorlage. Adresse: **Meldung an das Kurs-Projekt**
    (`pt9912/ai-harness-course`, repo-extern; Übergabe an den Auftraggeber im
    Bericht dieser Sitzung). Im Repo kein Folge-Slice — eine lokale Abweichung
    von der Vorlage wäre eine Adaption (`MR-*`); der Freshness-Audit des
    nächsten Bumps (`harness/targets/pin-stale.md` §Bump-Ablauf) liest MR-000
    gegen die dann gültige Vorlage.
  - **`docs/plan/carveouts/.gitkeep`:** neben der neuen `README.md` nicht mehr
    nötig, bleibt aber — das Repo hält `.gitkeep` auch neben Inhalt in
    `docs/plan/adr/`, `docs/reviews/` und `harness/conventions/` (Review,
    Negativbefunde); kein Konventions-Anker, keine Folgearbeit.
  - **Architect-Verdikt, Register:** kein Beleg für
    `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` — das Verdikt stuft
    Liefer-Punkt 3 als Verfeinerung innerhalb der bestehenden Bindung ein
    (akzeptiertes Negativ); der Eintrag bleibt bei 2×.
  - Keine Folgearbeit im Repo offen: der gemeldete Träger (Zähler „3×“ in
    `harness/sensors/docs-check.md` §Bindung) ist in der Fixrunde nachgezogen
    (F-7).
- **Steering-Loop-Eintrag:** Lernpunkt ohne neue Regel im Text und ohne neuen
  Sensor: *Wer aus einer Auszählung eine Liste ableitet, schreibt den Befehl und
  die gedruckten Zahlen je Stand vor die Liste — „ausgezählt“ ohne beides ist
  eine übernommene, keine gemessene Aussage.* Die Regel stand
  (`AGENTS.md` §3.12 Instanz A: jede Zahl trägt ihren Ursprung, eine Messung
  ihren Lauf; Reviewer-Skill HIGH „Beleg trägt seinen Satz nicht“: die Zählung
  nachfahren); der Lauf schrieb das Ergebnis, nicht die Messung, und erst die
  Nachzählung des Reviewers fand die Abweichung. Gezählt unter
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert in
  `.harness/skills/reviewer.md`, HIGH). Benannte Spec-Lücke: keine.
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-abgeleitete-dokumente-vorlagen-nachzug.md` in
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht/` angelegt (Review F-1, HIGH,
  daher Datei trotz Deckel; F-4 im selben Vorgang zählt nicht ein zweites Mal;
  Zähler 23×, Ausgang unverändert verkörpert) und in
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen/` (Review F-2, MEDIUM, daher
  Datei trotz Deckel; V-2 im selben Vorgang; Zähler 20×, Ausgang unverändert
  verkörpert). Ohne Register-Eintrag: F-3 (Baseline-Regel enger zitiert als
  ihre Achse, LOW, in der Fixrunde behoben, keine wiederkehrende Klasse im
  Register), V-1 (siehe *Befunde*). Kein Eintrag erreicht mit diesem Slice die
  Schwelle 3× neu; kein Lese-Schritt fällig.
- **Folge-Slices:** keine.
- **Risiken aus §6:** je ein Ausgang, siehe §6 — beide entfallen (jede Zeile an
  ihrer Quelle nachgeprüft; das weite Muster meldet am Bestand 0 Befunde).
- **Drei Paarungen:** nach dem `git mv` gemessen — (a) *Anker:* kein Eintrag
  dieser Notiz trägt das Feld `liegt in` (nichts verkörpert), kein Gegenstand
  der Paarung; (b) *Folge-Slice:* keiner genannt; (c) *Register:* die vier
  genannten Verzeichnisse (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`,
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`,
  `BEO-PGC/bindung-spalte-uneinheitlich-tief`) existieren, jedes mit nicht
  leerem `evidence/` (23, 20, 2, 1 Dateien). Ergebnis: getragen. Durch den Move
  brach kein Verweis (`git grep` nach dem `in-progress`-Pfad dieses Plans: drei
  Treffer, alle als Inline-Code in Mutations-Tabellen — dieser Plan §3, das
  Architect-Verdikt und der Verifikationsbericht —, kein Link). Der Ruhe-Marker
  der Roadmap steht wieder, Wortlaut gleich `15646fcb` (`diff` Exit 0);
  `in-progress/` trägt nur `roadmap.md`.
- **Gates der Closure:** `make docs-check` vor dem Inhalts-Commit Exit 0;
  `make gates` und `make suchlauf-nachmessen` nach dem letzten Commit stehen im
  Bericht der Sitzung, nicht in diesem Plan.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Der Abschnitt selbst entfällt nie.**

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt eine Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`harness/`, `docs/plan/adr/`, `docs/plan/carveouts/`,
`.d-check.yml`) liegen in ihr.

**Vorgelagert — offene Beobachtungen sichten:** beim Start durch den Implementer
(Register `docs/plan/planning/observations/BEO-PGC/`). Bekannt bei der Anlage:
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (offen, 2×; Liefer-Punkt 3
wäre das dritte Auftreten, wenn die Erweiterung ohne ADR-Klärung käme).
Sichtung beim Start (Implementer, Register am Stand `3f51d3ed`, Verzeichnisnamen
mit `vorlag|platzhalt|template|baseline|abgeleitet` und die offenen Einträge):
`gate-scope-erweiterung-ohne-adr-traeger` bleibt bei 2× — das Architect-Verdikt
`architect-verdict-slice-pfad-waechter-namensform` stuft Liefer-Punkt 3 als
Verfeinerung ohne ADR ein, kein Auftreten der Klasse;
`report-nicht-aus-baseline-vorlage` (offen, 1×) betrifft Review-Reports, nicht
die vier Dokumente; `plan-vorlagen-defekt` und `vorlagenrest-in-closure-notiz`
sind verkörpert. Kein weiterer Treffer.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
