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
- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: Träger, die die geänderten Abschnitte zitieren, sind
      nachgezogen (Suchlauf, `AGENTS.md` §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/README.md`, `harness/conventions.md` | update | Liefer-Punkt 1 |
| `docs/plan/adr/README.md`, `docs/plan/carveouts/README.md` | update / neu | Liefer-Punkt 2 |
| `.d-check.yml` (zwei `structure`-Regeln), `harness/sensors/docs-check.md` Grenze | update | Liefer-Punkt 3 |
| `docs/plan/planning/README.md` | update (Plan-Nachzug) | Liefer-Punkt 2: der Kommentar `d-check:ignore` am Pfad `docs/plan/carveouts/done/` fällt — er ist wirkungslos (Messung unten) und die einzige Abweichung der Zeile von `planning/README.template.md` Zeile 47 |

### Ergebnis je Liefer-Punkt (Implementer, Stand vor dem Commit, Parent `3f51d3ed`)

**Liefer-Punkt 1 — `harness/README.md`.** Je Platzhalter, Ausgang und Quelle:

| Stelle | Ausgang | Quelle / Grund |
|---|---|---|
| Gate-Tabelle, Zeile `` `<make-target>` `` „volle Closure“ | gestrichen | Das Makefile führt kein Closure-Ziel (`grep -nE '^(fullbuild\|ci)[: ]' Makefile harness/mk/*.mk d-check.mk`: 0 Treffer); der Image-Digest ist ein lokaler Lauf-Beleg von `make image`, das in der Werkzeug-Tabelle steht ([`ADR-0103`](../../adr/0103-image-hash-lokal-statt-committet.md)) |
| Werkzeug-Tabelle, Musterzeilen `make <mover>`, `make <messung>`, `make <vorschau>` | gestrichen | Formbeispiele der Vorlage; die Werkzeug-Tabelle führt die realen Ziele dieser Arten |
| §Safety and scope boundaries, zweimal `<…>` | ausgefüllt, fünf Punkte | je Punkt die zitierte Quelle in der Zeile selbst: `spec/lastenheft.md` §MVP-Schnitt (Reihenfolge, Rollback, Neustart), `LH-FA-CFG-006`, `LH-QA-SEC-002` mit `ADR-0047`, die drei Kennungs-Gates aus §Sensors, `AGENTS.md` §3.1 |
| §Leseordnung, drei Platzhalter | ausgefüllt, vier Zeiger | Rang-Reihenfolge von `AGENTS.md` §2 (`README.md` Rang 7 als Überblick, `AGENTS.md` §3, `spec/lastenheft.md` Rang 1) und das Bei-Bedarf-Beispiel der Vorlage (`harness/conventions.md`) |

**Liefer-Punkt 1 — `harness/conventions.md`.**

| Stelle | Ausgang | Quelle / Grund |
|---|---|---|
| `**Extern (Lehrmaterial):** <Pfad oder URL>` | ausgefüllt | `https://github.com/pt9912/ai-harness-course` — das Repo des vendored Release-Assets (Zeile darunter) und der Abfrage in `tools/harness/pin-stale-baseline.sh` Zeile 26 |
| `**In-Repo (verkörperte Form):** <Pfade zu …>` | ausgefüllt | die Liste der abgeleiteten Dokumente in `harness/targets/pin-stale.md` §Bump-Ablauf Schritt 2, mit `docs/plan/carveouts/README.md` statt des Verzeichnisses |
| MR-000 `**Datum:** <Datum>` | `2026-09-09` | `git log --follow --format='%h %ad' --date=short -- harness/conventions.md`: ältester Commit `40c8c431` 2026-09-09 (legt die Datei an) |
| §Aufgelöste Adaptionen, Musterzeile `\<NNN\>` | `— \| —` | `harness/conventions/done/` existiert nicht (`ls`: nicht gefunden) |
| §Zusatzklassen, Musterzeile | fünf Klassen | die Bindung-Spalte beider Tabellen in `harness/README.md` §Sensors, ausgezählt: `LH-*` (3 Zeilen), Vertragsdatei `sensors/…`/`targets/…`, `seit slice-…`, `BEO-PGC/…` (1), `MR-*` (1); je Klasse ein Beispiel aus einer realen Zeile |
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
Verzeichnis und bleibt richtig. Gemeldet, nicht mitgeändert (ein anderer
Gegenstand als die bewegten Eigenschaften): `harness/sensors/docs-check.md`
§Bindung zählt `BEO-PGC/slice-pfad-als-link-in-berichten` mit „3×“, der
`state.md` des Eintrags mit 5× — Frist: die Closure dieses Slice, der Planner
zieht nach oder benennt den Träger mit Adresse.

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
  **Ausgang:** offen bis zur Closure.
- **Das erweiterte Muster meldet Bestand.** Gemessen bei der Anlage: der weite
  Ausdruck trifft in den zwei Dateiklassen eine Zeile, in Inline-Code (für die
  Regel unsichtbar). Nachgemessen vom Implementer: d-check mit dem weiten Muster
  über den Arbeitsbaum Exit 0, 0 `section-forbidden` (§3, Liefer-Punkt 3). —
  **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln. Ging der
Gegenstand an einen anderen Slice oder entfiel er, trägt diese Sektion die Zeile
`Gegenstand:` mit Kennung oder Grund.

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

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
