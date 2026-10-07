# ADR-0161: Baseline-Pins in ADRs und MR-Einträgen bleiben auf dem Stand ihrer Abfassung

**Status:** Accepted — Supersedes [`ADR-0157`](0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(nur Entscheidung 4: der Pin-Commit und der normalisierte `cmp` je MR-Datei;
Entscheidung 1 bis 3 und 5 bleiben in Kraft) · Supersedes
[`ADR-0159`](0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(nur Entscheidung 2 „MR-Pins: Datei-`cmp` und Referent-Messung“ und der Teil der
Folgepflicht, der `.claude/agents/verifier.md` und `harness/targets/pin-stale.md`
die Referent-Messung am Pin-Commit vorschreibt; Entscheidung 1, 3 bis 5 bleiben)
· Supersedes [`ADR-0160`](0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
(nur Entscheidung 2 „Die Teile von `ADR-0157`, die bleiben“; der Leer-Test aus
Entscheidung 1 bleibt und gilt für den Form-Commit dieser ADR, Entscheidung 3
bleibt unberührt)

**Datum:** 2026-10-07

**Autor:** pt9912 (Architect-Rolle, Modul 8; den Weg „Pins einfrieren“ hat der
Auftraggeber gewählt, diese ADR trägt seine Begründung, die Grenze und die
Messungen)

**Bezug:** [`ADR-0156`](0156-versions-gate-nimmt-done-records-aus.md) (dieselbe
Einfrierung nach Zeit für `docs/plan/planning/done/**`; diese ADR ist ihre Folge
und lässt sie unverändert) · [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Zitat-Korrektur der Form) · [`ADR-0051`](0051-cicd-pipeline-github-actions.md)
Entscheidung 7 (Pin-Inventar P8, Bump als bewusster Bootstrap-Vorgang) ·
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) (gemessen/hergeleitet)
· `AGENTS.md` §3.5, §3.6, §3.12 · `.d-check.yml` Blöcke `versions:` und `vcs:`
· `harness/targets/pin-stale.md` §Bump-Ablauf · Slice
`slice-harness-baseline-v6-16-0`, Abschnitt „Halt vor MR-001 und den
ADR-Zeilen“ (Anlass).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0156` und `ADR-0157`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Das d-check-Modul `versions` meldet jeden Verweis in die vendored Baseline,
dessen Versions-Segment vom adoptierten Stand (`harness/conventions.md`
§Baseline) abweicht, als `version-stale`. Ausgenommen sind bisher
`harness/conventions/done/**`, `docs/reviews/**` und
`docs/plan/planning/done/**` (`ADR-0156`). `Accepted`-ADRs und aktive
MR-Einträge sind nicht ausgenommen, obwohl beide ab `Accepted` bzw. ab ihrer
`Datum`-Zeile immutabel sind (`AGENTS.md` §3.5, `ADR-0157` Entscheidung 3).

**Folge bei jedem Bump** (übernommen aus der Commit-Folge; die Kennungen nennen
`ADR-0156` und `ADR-0157`): Die Pins dieser Dateien werden per Zitat-Korrektur
nachgezogen — `996e6231`, `54c6e632`, `cd364fc5` (u. a. `ADR-0095`), der
MR-Pin-Commit `5d8855d9`. Jede solche Korrektur verlangt die Referent-Messung
(`ADR-0157` Entscheidung 1 (b)), an MR-Einträgen zusätzlich Pin-Commit,
Teil-Ranges und Datei-`cmp` (`ADR-0157` Entscheidung 4, `ADR-0159`
Entscheidung 2, `ADR-0160`).

**Anlass, gemessen** (Slice-Plan, Abschnitt „Halt vor MR-001 und den
ADR-Zeilen“, und eigener Lauf V0 unten): Beim Bump v6.14.1 → v6.16.0 meldet
`make docs-check` sechs `version-stale`, die vier MR-Einträge, `ADR-0095`
Zeile 107 und `ADR-0160` Zeile 278. Drei davon lassen sich nach den geltenden
Regeln **nicht** nachziehen:

- `MR-001` → `grundlagen-referenz-richtung.md#spec-straten-mehr-als-ein-spec-dokument`:
  `cmp 1`, der Abschnitt wächst in v6.16.0 (*übernommen* aus dem Slice-Plan).
- `ADR-0095` Zeile 107 → `templates/.d-check.yml`: `cmp 1` (*übernommen*).
- `ADR-0160` Zeile 278 ist ein Messbeleg: sie beschreibt, welchen Pfad ein
  simulierter Pin-Commit umstellte. Ein Nachziehen fälschte den Beleg.

Eine Umstellung bei geändertem Referenten ist keine Zitat-Korrektur
(`ADR-0157` Entscheidung 1). Der geltende Weg wäre je Zeile eine Folge-ADR, die
nur einen Tag im Pfad neu fasst — das Werkzeug, das `ADR-0073` Option C und
`ADR-0157` Option B als falsch verworfen haben. Das Muster tritt bei jedem Bump
wieder auf.

**Was die Baseline dazu sagt.** Das Regelwerk verlangt das Gegenteil des
Einfrierens: „nicht die Form wechseln, damit nichts mehr rotten kann, sondern
das Rotten **sichtbar** machen“, der Versions-Sensor prüft **jeden** Pin
(`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln,
Spiegelstrich „Zwei Rot-Quellen, ein Prinzip“, Stand v6.16.0); die
MR-Vorlage wiederholt es für das Feld `Ersetzt-Baseline-Regel`. Dieses Repo
weicht davon ab. Die Abweichung ist deshalb eine Adaption und braucht einen
`MR`-Eintrag (Folgepflicht, Entscheidung 6).

**Ein Befund aus der Messung** (gemessen, M-L): Die vier MR-Links in die
Baseline stehen mit **mehrzeiligem** Linktext. `links` und `anchors` sehen sie
nicht. Ein Linkziel auf eine fehlende Datei unter v6.16.0 an genau dieser Stelle
in `MR-001` bleibt grün, dieselbe Probe als einzeiliger Link in `MR-001` meldet
`target-missing` (M7). Die Zusage „den Formcheck erledigt das Doku-Gate“
(`modul-02-harness-bootstrap.md` §Freshness-Audit, Stand v6.16.0) und das
Contra von `ADR-0157` Option E („verliert die Link-/Anker-Prüfung“) hatten für
den Bestand also kein Objekt. Das erklärt auch, warum der gebrochene Anker in
`MR-001` (`ADR-0159` Entscheidung 2, Review F-9) nie rot war.

## Entscheidung

1. **Baseline-Pins in ADRs und MR-Einträgen nennen den Stand ihrer Abfassung.**
   Ein Bump zieht sie nicht nach. Für den Bestand gilt der Stand, den der Pin am
   Commit dieser ADR trägt (`ADR-0095` und `ADR-0160`: v6.14.1, `MR-001` bis
   `MR-004`: v6.14.1; gemessen mit `git grep`, Fitness Function). Ein neuer Pin
   nennt den Stand, der beim Schreiben adoptiert ist.

2. **`versions.exempt-paths` nimmt `docs/plan/adr/[0-9]*.md` und
   `harness/conventions/MR-[0-9]*.md` auf**, mit der Begründung von
   `ADR-0156`: Ein `version-stale` in diesen Dateien entsteht durch einen Bump
   außerhalb der Datei, nicht durch ihren Text. Die Ausnahme ist **pfad-**,
   nicht **status**gebunden. `versions.exempt-paths` ist eine datei-weite
   Glob-Liste ohne Bedingung auf den Inhalt (`d-check --print-config`, v0.82.0,
   gelesen). Eine Trennung nach `Accepted` ginge nur über eine Liste der
   Dateinamen, die mit jeder ADR wächst. Gewinnen würde sie wenig: zwei der 160
   ADR-Dateien tragen `Proposed` (gemessen mit `grep`), und eine ADR wird hier
   in der Regel im Commit ihres Schreibens `Accepted`. Der Index
   `docs/plan/adr/README.md` fällt nicht unter `[0-9]*` und bleibt geprüft (M3).

3. **`links` und `anchors` bleiben ohne Ausnahme.** Bricht ein Link aus einer ADR
   oder einem MR-Eintrag, weil ein altes Tag-Verzeichnis gelöscht wird, gilt die
   **Form-Korrektur** nach `ADR-0073` und `ADR-0156` Entscheidung 3: Aus dem
   Markdown-Link wird Inline-Code mit dem Pfad ab Repo-Wurzel samt alter Version
   und Anker; der Linktext bleibt als Prosa stehen, kein Wort der Aussage ändert
   sich. Form für ein Feld `Ersetzt-Baseline-Regel`:

   ```text
   vorher:  [`<datei>.md` §<Abschnitt>](../../.harness/baseline/<alt-tag>/regelwerk/<datei>.md#<anker>)
   nachher: `<datei>.md` §<Abschnitt> (`.harness/baseline/<alt-tag>/regelwerk/<datei>.md#<anker>`)
   ```

   Die Korrektur fällt **einmal** je Link an, beim Löschen seines Tags. Danach ist
   das Verweisgerüst eingefroren. Sie gilt auch für die MR-Links, die `links`
   heute nicht sieht (Kontext): Ein toter Link bleibt eine falsche Adresse, und
   ein Werkzeug, das mehrzeilige Linktexte liest, färbte ihn bei einem fremden
   Anlass rot. Den Referenten verändert die Korrektur nicht, die Adresse bleibt
   byte-gleich (gemessen an den vier MR-Verweisen, roh `cmp 0`).

4. **Der Pin-Commit entfällt. An seine Stelle tritt der Form-Commit, nur wenn er
   gebraucht wird.** Zeigt ein Link in einem MR-Eintrag in das Tag, das ein Bump
   löscht, steht seine Form-Korrektur in einem eigenen Commit `F`. Dieser Commit
   ändert **nur** MR-Dateien, seine Message nennt `ADR-0073` und diese ADR, und
   er liegt **vor** dem Commit, der das alte Tag löscht. Der Verifier prüft:
   - `make doc-immutable` in den Teil-Ranges `B..F~1` und `F..H`, mit dem
     Leer-Test `teilrange` aus `ADR-0160` Entscheidung 1 (unverändert, `P` dort
     ist hier `F`);
   - am Form-Commit je MR-Datei den `cmp` nach Form-Normalisierung, der grün ist
     genau dann, wenn sich außer Linkklammern, Backticks, dem relativen Präfix
     `../` und Leerraum nichts bewegt hat:

     ```bash
     formnorm() { sed -E 's#\]\((\.\./)*\.harness/#.harness/#g' | tr -d '][()`' | tr -d '[:space:]'; }
     for f in $(git diff --name-only "$F~1" "$F" -- 'harness/conventions/MR-[0-9]*.md' 'harness/conventions/**/MR-[0-9]*.md'); do
       cmp -s <(git show "$F~1:$f" | formnorm) <(git show "$F:$f" | formnorm); echo "$f cmp $?"
     done
     ```

     Je MR-Datei steht eine Zeile da, jede mit `cmp 0`; ohne Zeile ist der
     Commit falsch bestimmt.

   An ADRs gibt es keinen `vcs`-Lauf. Ihre Form-Korrektur steht in einem eigenen
   Commit mit `ADR-0073` und je ADR einer §Geschichte-Zeile (`AGENTS.md` §3.5).

5. **Die Referent-Messung der MR-Einträge wird Prüfauftrag des
   Adaptions-Durchgangs, nicht Bedingung eines Commits.** Im Bump-Ablauf
   Schritt 1 (`harness/targets/pin-stale.md`, Durchgang durch die Adaptionen)
   misst der Messende je MR-Eintrag die Einheit seines Verweises in
   `Ersetzt-Baseline-Regel` **zwischen den Tags des Bumps**, nicht zwischen dem
   eingefrorenen Tag und dem neuen. Pfad und Anker kommen aus dem eingefrorenen
   Verweis, das Tag-Segment ist `<alt-tag>` bzw. `<neu-tag>`. Gemessen wird an
   einem Stand mit beiden Bäumen:

   ```text
   make zitat-vergleich ARGS="<stand> .harness/baseline/<alt-tag>/<pfad> '<anker>' <stand> .harness/baseline/<neu-tag>/<pfad> '<anker>'"
   ```

   Gemessen wird roh, bei `cmp 1` zusätzlich mit Tag-Paar `<alt-tag>:<neu-tag>`.
   - `cmp 0` (roh oder normalisiert): Die ersetzte Regel hat sich im Bump nicht
     bewegt. Kein Prüfauftrag.
   - `cmp 1` normalisiert oder Exit 2 (Anker oder Datei fehlt im neuen Tag): ein
     **Prüfauftrag** an den Durchgang. Er bekommt einen der fünf Ausgänge aus
     `modul-02-harness-bootstrap.md` §Freshness-Audit, mit Grund, im Plan des
     Bumps. Bleibt die Adaption gültig, bleibt die MR-Datei unberührt. Jeder
     andere Ausgang ist ein Nachfolge-Eintrag, wie bisher.

   Diese Messung ersetzt den Formcheck, den das Regelwerk dem Doku-Gate
   zuschreibt (Kontext). Sie ist strenger: Sie misst den Inhalt der Einheit,
   nicht nur, ob der Anker existiert.

   **Für `MR-001` in diesem Bump:** `cmp 1` normalisiert (übernommen aus dem
   Slice-Plan). Die Aussage der Adaption, der Datei-Name des Rang-2-Dokuments,
   berührt der neue Absatz nicht. Ausgang *bleibt gültig*, die MR-Datei bekommt
   nur die Form-Korrektur aus Entscheidung 3. Der Satz „Abschnitte 1–7“ in
   `MR-001` hängt am Folge-Slice zum Delta-Punkt „Neu 2“ (Festlegungen der
   Harness-Werkzeuge) und ist kein Pin-Thema. Ob der Satz dort einen
   Nachfolge-Eintrag braucht, entscheidet jener Slice.

6. **Die Abweichung von der Baseline wird als `MR-005` deklariert**
   (Folgepflicht). `Ersetzt-Baseline-Regel` ist
   `grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln,
   Spiegelstrich „Zwei Rot-Quellen, ein Prinzip“, die Versions-Hälfte. Die
   Adaption: Pins in ADRs und MR-Einträgen frieren ein, siehe Entscheidung 1
   bis 5. Sie ist `permanent`, ihr Auflösungs-Trigger ist der
   Re-Evaluierungs-Trigger (a) dieser ADR.

7. **Was weiter nachgezogen wird.** Alle Baseline-Verweise außerhalb der
   ausgenommenen Pfade bleiben auf den adoptierten Stand gepinnt, und jeder Bump
   zieht sie nach: `AGENTS.md`, `harness/README.md`, `harness/conventions.md`
   (§Baseline ist `current-from`), `harness/sensors/*.md`, `harness/targets/*.md`,
   `.claude/agents/*.md`, die Symlinks unter `.claude/rules/`, `.harness/skills/*.md`,
   der ADR-Index, die Pläne in `open/`, `next/`, `in-progress/`, die flachen
   Welle-Dateien und das Beobachtungs-Register. `versions` prüft davon alle
   Markdown-Dateien unter `scan.roots` (gemessen für `AGENTS.md`, den Index und
   einen Plan in `in-progress/`: M1, M3, M4). Die Symlinks liest es nicht, das
   trägt der Bump-Ablauf Schritt 1. `.harness/skills/**` liegt unter
   `scan.ignore` und wird ebenfalls nicht geprüft (*hergeleitet* aus der
   Konfiguration, nicht gefahren); dort trägt allein die Liste des Bump-Ablaufs.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: weiter je Bump nachziehen, bei geändertem Referenten Folge-ADR | `versions` bleibt über ADRs und MR-Einträge scharf | An drei Zeilen dieses Bumps nur per Folge-ADR möglich, die einen Tag im Pfad neu fasst (`ADR-0073` Option C, `ADR-0157` Option B: falsches Werkzeug). `ADR-0160` Zeile 278 ließe sich gar nicht korrekt nachziehen (Messbeleg). Jeder Bump erzeugt wieder Pin-Commits, Teil-Ranges und Referent-Messungen an Dateien, die sich inhaltlich nicht ändern |
| B — Zeilen-Marker `d-check:ignore` je Pin | punktgenau; alle anderen Pins bleiben geprüft | Der Marker ist ein Eingriff in `Accepted`-ADRs und MR-Einträge (Text, kein Gerüst). Jeder neue Pin braucht ihn beim nächsten Bump wieder. Dass der Marker für `versions` wirkt, ist *übernommen* aus `--print-config`, nicht gefahren |
| C — Ausnahme nur für `Accepted`-ADRs, als Liste der Dateinamen | ein falscher Pin in einer `Proposed`-ADR bliebe sichtbar | Die Liste wächst mit jeder ADR, und ein vergessener Eintrag macht den Bump wieder rot. Der Gewinn ist klein: 2 von 160 ADRs `Proposed` |
| D — Pins in ADRs und MR-Einträgen ohne Version schreiben (`<tag>`) oder ohne Pfad (Hausform) | kein Bump berührt sie | Weicht von der MR-Vorlage ab, die den Link mit Anker verlangt; der Bestand müsste einmal umgeschrieben werden; ein Verweis ohne Version sagt nicht, welche Fassung die Entscheidung las |
| **E — Einfrieren per Pfad, `links`/`anchors` scharf, Form-Korrektur einmal je gelöschtem Tag, Referent-Messung im Adaptions-Durchgang (gewählt)** | dieselbe Begründung und derselbe Mechanismus wie `ADR-0156`; kein Pin-Commit mehr; die drei Haltestellen dieses Bumps verschwinden ohne Eingriff in ihren Text; der Prüfauftrag an MR-Einträgen misst den Inhalt der ersetzten Regel, nicht nur ihren Anker | Ein falscher Pin in einer neuen ADR oder einem neuen MR-Eintrag fällt `versions` nicht mehr auf (M2, M6; §Konsequenzen). Eine weitere Adaption (`MR-005`) |

## Konsequenzen

- **Positiv:** Ein Baseline-Bump berührt `Accepted`-ADRs und MR-Einträge nur
  noch dort, wo ein Link nach dem Löschen bricht, einmal je Link. Heute sind das
  die vier MR-Links; aus ADRs zeigt kein Markdown-Link in die Baseline
  (gemessen). Der Pin-Commit, sein Datei-`cmp` und die Referent-Messung als
  Bedingung einer Zitat-Korrektur an MR-Einträgen entfallen.
- **Positiv:** Die Haltestellen `ADR-0095` Zeile 107 und `ADR-0160` Zeile 278
  bleiben unverändert und sind kein Befund mehr (V1). `ADR-0095` bekommt keine
  §Geschichte-Zeile.
- **Die bisherigen Zitat-Korrekturen bleiben stehen** (`996e6231`, `54c6e632`,
  `cd364fc5`, `5d8855d9`). Ein Rückbau auf den Stand der Abfassung wäre eine
  weitere Änderung an eingefrorenen Dateien ohne Gewinn.
- **Negativ, was das Gate verliert:**
  - *Gemessen:* Einen Inline-Pin mit falscher Version in einer ADR-Datei (M2)
    oder in einem MR-Eintrag (M5, M6) meldet `versions` nicht mehr. Ohne die
    Ausnahme meldet es denselben Pin (M2b).
  - *Hergeleitet:* Das gilt auch für eine `Proposed`-ADR und für einen Pin, der
    schon beim Schreiben falsch war. Anders als bei `ADR-0156` gibt es keine
    Vorstufe, auf der `versions` scharf ist, denn eine ADR entsteht in
    `docs/plan/adr/`. Leser bleibt der Reviewer am Diff (`AGENTS.md` §3.12,
    „Verfasser einer ADR“). Das ist ein akzeptiertes Negativ: Ein Pin in einer
    ADR nennt die Fassung, die die Entscheidung las, und wird nicht als
    Arbeitsanweisung gelesen.
  - *Gemessen, bleibt scharf:* einzeiliger Link aus einem MR-Eintrag in ein
    gelöschtes Tag (M7), Link aus einer ADR auf eine fehlende Datei (M8) und
    einen fehlenden Anker (M9), Pins in `AGENTS.md`, im ADR-Index und in einem
    Plan unter `in-progress/` (M4, M3, M1).
- **Negativ:** Ein Leser einer ADR folgt einem eingefrorenen Pfad, den es im
  Baum nicht mehr gibt. Er löst in der Git-Historie auf (`git show
  <commit>:<pfad>`). Das gilt seit `ADR-0156` schon für Records.
- **Negativ, benannte Lücke:** Der Form-`cmp` übersieht eine Änderung, die nur
  Leerraum, Linkklammern, Backticks oder `../` betrifft, etwa ein entferntes
  Leerzeichen zwischen zwei Wörtern. Ob nur die Form bewegt wurde, bleibt
  daneben Urteil am Diff (`ADR-0073`).
- **Folgepflicht (Implementer, `slice-harness-baseline-v6-16-0`):**
  1. `.d-check.yml`, Block `versions:`: `exempt-paths` nach Entscheidung 2,
     Kommentar mit Anker auf diese ADR; Block `vcs:`: der Kommentar nennt den
     Form-Commit statt der Pin-Umstellung.
  2. Form-Commit `F` an `MR-001` bis `MR-004` nach Entscheidung 3/4, vor dem
     Löschen von v6.14.1; die Form-Korrektur der `done/`-Zeile nach `ADR-0156`
     Entscheidung 3 in einem eigenen Commit.
  3. `MR-005` nach Entscheidung 6, mit Zeile im Index von `harness/conventions.md`.
  4. `harness/targets/pin-stale.md`: Absatz „MR-Pins in einem eigenen Commit“
     ersetzt durch Entscheidung 1, 4, 5; Schritt 1 nennt die Referent-Messung
     des Adaptions-Durchgangs.
  5. `.claude/agents/verifier.md` und `.claude/agents/implementer.md`: Pin-Commit
     durch Form-Commit ersetzt, `cmp` nach Entscheidung 4, die Referent-Messung
     am Pin-Commit gestrichen.
  6. `AGENTS.md` §3.5, Absatz **Beleg**: Verweis „für einen MR-Eintrag
     `ADR-0159` Entscheidung 2“ ersetzt; ein Satz zum Einfrieren.
  7. `harness/targets/zitat-vergleich.md`: „MR-Pin“ und „Pin-Commit“ an den
     drei Stellen nachgezogen.
  8. Im Slice-Plan: Ausgang für die drei Haltestellen und für das Risiko
     „Referent von MR-001 bewegt sich“ (*eingetreten*, diese ADR); der Beleg
     des Adaptions-Durchgangs für `MR-001` bis `MR-004` nach Entscheidung 5.
  Der ADR-Index trägt diese ADR bereits.

## Fitness Function (falls maschinell prüfbar)

**Instanz aller Zeilen:** d-check v0.82.0
(`ghcr.io/pt9912/d-check@sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`,
der Pin aus `d-check.mk`), `docker run --rm --network none -v <Klon>:/repo:ro`,
gegen einen Klon von `18bbdc98` im Scratchpad (v6.14.1 und v6.16.0 im Baum).
Die Konfiguration ist per Kopie um Entscheidung 2 ergänzt, ab V1. Jede Mutation
steht an **einer** Stelle, meist als angehängte Zeile in der genannten Datei,
und wurde danach zurückgesetzt. Der Lauf im Arbeitsbaum mit der echten
`.d-check.yml` ist eine **Erwartung** an den Implementer und noch nicht
erprobt.

| Lauf | Zustand | Gesehen |
|---|---|---|
| V0 | Konfiguration unverändert | Exit 1, `6 Befund(e)`, alle `version-stale`: `ADR-0095`:107, `ADR-0160`:278, `MR-001`:12, `MR-002`:13, `MR-003`:13, `MR-004`:19 |
| V1 | + Entscheidung 2 | Exit 0, `0 Befund(e)` |
| V2 | V1, v6.14.1 per `git rm -r` entfernt | Exit 1, `1 Befund(e)`: `target-missing` in `done/slice-harness-guard-blocked-python.md`:272. Die vier MR-Links in das entfernte Tag melden nichts (mehrzeiliger Linktext) |
| V3 | V2 + Form-Korrektur nach Entscheidung 3 an `MR-001` bis `MR-004` und an der `done/`-Zeile | Exit 0, `0 Befund(e)` |
| M-L | V2-Konfiguration, Linkziel in `MR-001` auf `v6.16.0/regelwerk/gibtsnicht.md` umgebogen, mehrzeiliger Linktext | grün für `MR-001` (nur der `done/`-Befund aus V2): `links` sieht den mehrzeiligen Link nicht |
| M1 | V3 + Inline-Pin v6.13.0 im Plan unter `in-progress/` | rot: 1 `version-stale` |
| M2 | V3 + neue ADR-Datei `0161-probe.md` mit Inline-Pin v6.13.0 | grün: 0 Befunde (der Verlust) |
| M2b | wie M2, Konfiguration ohne Entscheidung 2 | rot: 7 `version-stale`, darunter die Probe-Datei |
| M3 | V3 + Inline-Pin v6.13.0 im ADR-Index `README.md` | rot: 1 `version-stale` (`[0-9]*` lässt den Index aus) |
| M4 | V3 + Inline-Pin v6.13.0 in `AGENTS.md` | rot: 1 `version-stale` |
| M5 | V3 + Inline-Pin v6.13.0 in `MR-001` | grün: 0 Befunde (der Verlust) |
| M6 | V3 + neue Datei `harness/conventions/MR-005-probe.md` mit Inline-Pin v6.13.0 | grün: 0 Befunde (der Verlust) |
| M7 | V3 + einzeiliger Link aus `MR-001` in das entfernte v6.14.1 | rot: 1 `target-missing` |
| M8 | V3 + Link aus `ADR-0160` auf eine fehlende Datei unter v6.16.0 | rot: 1 `target-missing` |
| M9 | V3 + Link aus `ADR-0160` mit fehlendem Anker in `regelwerk/README.md` | rot: 1 `anchor-missing` |

**Form-Commit, gemessen im selben Klon:** `B` = `18bbdc98`. Commit 1 trägt
Konfiguration, `done/`-Zeile und das Löschen von v6.14.1, Commit `F` nur die
Form-Korrektur an `MR-001` bis `MR-004`.

| Tooling | Regel | Make-Target |
|---|---|---|
| d-check `versions` | `exempt-paths` enthält `docs/plan/adr/[0-9]*.md` und `harness/conventions/MR-[0-9]*.md`; außerhalb bleibt jeder Pin geprüft (M1, M3, M4) | `make docs-check` (in `make gates`) |
| d-check `links`/`anchors` | keine Ausnahme; Linkbruch und fehlender Anker aus ADR und MR-Eintrag bleiben rot (M7–M9). Grenze: mehrzeiliger Linktext wird nicht gelesen (M-L) | `make docs-check` (in `make gates`) |
| `formnorm`-`cmp` (Entscheidung 4), Befehl aus dem `bash`-Block | an `F` je MR-Datei `cmp 0` (vier Zeilen). Mutation an je **einer** Stelle, gefahren als Shell-`cmp` gegen eine Kopie: ein Wort in `MR-002` („unabhängig“ → „abhängig“) `cmp 1`; der Tag in `MR-001` (v6.14.1 → v6.16.0, also Pin statt Form) `cmp 1`. Dass jede andere Wortänderung ebenso fällt, ist *hergeleitet* | kein Make-Target (Verifier-Lauf) |
| `teilrange` (`ADR-0160` Entscheidung 1) um `F` | `B..F~1`: „enthält 1 Commit(s), Lauf“, `0 Befund(e)`, Exit 0 — grün. `F..HEAD` (`HEAD` = `F`): „leer …, kein Lauf“, Exit 0. Gegenprobe volle Range `make doc-immutable RANGE=B..F`: `4 Befund(e)`, je MR-Datei `core-drift-vcs`, make-Exit 2 | kein Make-Target (Verifier-Lauf) |
| `make zitat-vergleich`, Adresse vor und nach der Form-Korrektur | Form-Commit vor dem Löschen (eigener Zweig ab `18bbdc98`): je MR-Verweis `vergleich roh: F~1:<alt-adresse> <-> F:<alt-adresse> cmp 0`, viermal. Der Referent bleibt also gleich (Entscheidung 3) | `make zitat-vergleich` (kein Gate) |
| Bestand (`git grep` am Stand `18bbdc98`) | Pins nach `pin-pattern` in `docs/plan/adr/[0-9]*.md`: 2 (`ADR-0095`:107, `ADR-0160`:278); in MR-Einträgen: 4, je einer pro Datei. Markdown-Links in die Baseline aus ADRs: 0, aus MR-Einträgen: 4, alle mehrzeilig. Links in das Tag v6.14.1 im Baum ohne `.harness/`: 5 (vier MR, eine `done/`-Zeile; die Treffer unter `docs/reviews/` stehen als Inline-Code, V2) | — |

## Re-Evaluierungs-Trigger

Beobachtbare Trigger:

- **(a)** d-check bietet für `versions` eine Bedingung auf den Inhalt einer
  Datei, etwa auf die Status-Zeile. Dann wird die Ausnahme auf `Accepted`
  verengt; das schließt den Verlust an `Proposed`-ADRs (M2).
- **(b)** Ein eingefrorener Pin führt einen Leser nachweislich in die Irre,
  etwa eine Arbeit, die einer alten Fassung folgt. Dann ist die Annahme des
  akzeptierten Negativs widerlegt; es folgt eine Folge-ADR.
- **(c)** Die Baseline gibt die Regel „Zwei Rot-Quellen, ein Prinzip“ für die
  Versions-Hälfte auf oder nimmt immutable Dokumente aus. Dann wird `MR-005`
  im Adaptions-Durchgang *gegenstandslos*.

Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-07 | Accepted. Anlass: drei `version-stale` beim Bump auf v6.16.0, die keine Zitat-Korrektur zulassen (`MR-001`, `ADR-0095`:107, `ADR-0160`:278); Weg vom Auftraggeber gewählt | `slice-harness-baseline-v6-16-0` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0161` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
