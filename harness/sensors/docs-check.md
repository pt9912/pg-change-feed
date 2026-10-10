# `make docs-check` — prüft Markdown-Doku auf kaputte Referenzen (d-check)

## Vertrag

Wird dieses Target rot, trägt die gescannte Markdown-Doku eine kaputte
Referenz, eine nackte Kennung, eine verbotene Referenzrichtung, einen
abweichenden Baseline-Pin, einen Struktur-Verstoß im Abschnitt, einen
host-lokalen Pfad oder ein Linkziel außerhalb des git-Index — je Modul der
Liste `modules:` in `.d-check.yml` eine Klasse.

**Wer was trägt.** Die **Festlegung** — welche Fläche gescannt wird, was jedes
der acht Module als Befund meldet und unter welchem Grund-Code, wie es an seinen
Randformen entscheidet, welche Abschnitte die `structure`-Regeln prüfen und was
die Ausgänge bedeuten — steht in
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge);
diese Datei wiederholt sie nicht. Die Werte (Muster, Klassen, Pfadlisten,
Spaltengrenzen) stehen in `.d-check.yml`, der Eingabe des Laufs. Diese Datei
sagt, wie ein Lauf zu lesen ist und was sein Grün nicht abdeckt. Ein
Schwester-Artefakt wird als blankes Repo-Wort mit relativem Pfad zitiert
([`ADR-0074`](../../docs/plan/adr/0074-zitationsform-schwester-repo-hausform.md)).

## Grenze — was das Grün nicht abdeckt

1. **Within-Spec-Ordnung** — `direction: no-downward` ist bewusst nicht
   gesetzt (die Lastenheft-Messmethoden delegieren an Pflichtenheft-
   Festlegungen, [`SPEC-012`](../../spec/pflichtenheft.md)/013/014); die Ordnung Vertrag > Technik > Sicht
   bleibt Review-Prüfpflicht. Heilbar durch Config, wenn die Delegationen
   entfallen.
2. **`codepaths` aus** — Pfade in Inline-Code werden nicht auf Existenz
   geprüft; Bedingung in `.d-check.yml` (alle referenzierten Pfade
   existieren). Heilbar.
3. **Opt-in-Module nicht im Bündel** — `planning`, `vcs`, `commits`,
   `reviews` laufen nur über ihre `doc-*`-Einzel-Targets mit
   `--enable`; `make gates` belegt sie nicht. Heilbar je Aktivierungs-
   bedingung (in `.d-check.yml` kommentiert). `tracked` ist seit
   `slice-d-check-tracked-modul` Teil des `modules:`-Bündels und läuft
   damit in `make docs-check`/`make gates` mit; sein eigenes
   `doc-tracked`-Einzel-Target bleibt zusätzlich isoliert aufrufbar.
4. **`MR-*` nicht linkpflichtig** — das `ids`-Muster deckt LH/SPEC/ARC/ADR,
   nicht MR; Adaptions-Verweise werden vom `tracked`-Modul auf ihren
   Getrackt-Status geprüft (seit dessen Aktivierung Teil des Bündels, kein
   Opt-in mehr) — nicht auf Linkpflicht. Permanent bis zur Muster-Erweiterung.
5. **Vendored Bestand ausgenommen** — `.harness/**` und `**/*.template.md`
   sind vom Scan ausgenommen; die Baseline selbst prüft `baseline-verify`,
   die Templates sind Referenz-Form. Permanent (Setzung).

6. **Verweisform auf wandernde Slice-Pläne — Textform, kein Link-Baum.** Die
   beiden `structure`-Bedingungen über `docs/reviews/**` und
   `docs/plan/planning/observations/**/observation.md` lesen den bereinigten
   Abschnittstext: ein Lifecycle-Pfad in Inline-Code oder im Fence ist
   unsichtbar — die zulässige Inline-Code-Form bleibt damit grün —, und
   Reference-Style-Links umgehen die Textform. `state.md` und `evidence/*.md`
   tragen keine Überschrift und haben deshalb keinen Abschnitts-Anker; dort
   trägt die Selbstprüfung im Closure-Schritt. Die Wellen-Form (flaches
   `planning/welle-<Kennung>.md` → `done/`) und der gleich-ordnerige Nachbar-Verweis
   (Quell-Seite, `links.resolve-from`) sind Nachbar-Klassen, nicht gedeckt.

7. **Die erzeugte E2E-Abdeckungstabelle trägt keine Symbol-Prüfung.** Die
   Tabelle `docs/user/e2e-abdeckung.md` ist ein Erzeugnis von
   `make test-integration`: die Go-Zeilen leitet das Testpaket per
   `go/parser` aus seinem eigenen Quelltext ab, die Bash-Zeilen deklariert
   jede Phase des Runners über einen Anker, geschrieben wird nur bei
   inhaltlicher Abweichung. Das Doku-Gate prüft an ihr die **Form** — die
   Linkpflicht des `ids`-Moduls auf der Kennungsspalte und die
   `structure`-Regel (Abschnitt, Spalten-Mindestbreiten) —, nicht den
   Nachweis. Vier Grenzen sind daran real gemessen, nicht angenommen:

   - **`ids` prüft den Link, nicht die Existenz.** Eine nackte Kennung ohne
     Link meldet `id-unlinked`; eine **verlinkte, erfundene** Kennung bleibt
     grün — das `ids`-Modul gleicht nicht gegen die Kennungen des
     Definitions-Dokuments ab.
   - **Inline-Code-Spans bleiben ungeprüft.** Eine Kennung in einem
     Code-Span ist auch ohne Link grün; die Beschreibungsspalte der Tabelle
     trägt deshalb **keine** Kennungen aus dem Doc-Kommentar (der Erzeuger
     entfernt sie samt umgebendem Span). Was nur im Kennungsverweis des
     Kommentars stand, bleibt über die Spalte `Ort` an seiner Quelle
     erreichbar.
   - **Kein Modul prüft einen Symbol- oder Funktionsnamen.** Ein Nachweis
     darf einen Funktionsnamen tragen, den das Paket nicht kennt, ohne dass
     ein Befund entsteht; `--trace` gibt die Requirements-Traceability-Matrix
     aus und ist keine Code-Prüfung (die Tabelle speist sie auch nicht — der
     Trace-Ausgang ist mit und ohne die Datei identisch), `codepaths` prüft
     Pfade und Zeilenbereiche statt Symbole.
   - **Die Bash-Hälfte ist deklariert, nicht abgeleitet.** Der
     Deklarations-Anker deckt die eine Richtung („die Tabelle behauptet eine
     Phase, die es nicht gibt" — der Lauf bricht ab), nicht die andere: eine
     neue Runner-Phase, die niemand deklariert, fehlt still in der Tabelle.
     Dieselbe Klasse führt `BEO-PGC/test-runner-stiller-ausschluss` für die
     `-run`-Muster (Register `verkörpert`): die Vollständigkeits-Hälfte trägt
     `TestRunnerFuehrtJedeE2EFunktionAus` in
     `test/integration/runner_vollstaendigkeit_test.go` — jede `func TestE2E*`
     steht in einem `-run`-Wert des Runners; die Phasen ohne Go-Testfunktion
     bleiben deklariert, nicht abgeleitet.

   Die tragende Garantie der Tabelle ist deshalb die **Ableitung aus dem
   Quelltext** (Go-Hälfte) und der **Anker je Phase** (Bash-Hälfte) — nicht
   das Doku-Gate.

8. **`hostpaths` — die gescannte Markdown-Fläche, nicht das Repo.** Welche
   Pfade das Modul meldet und welche nicht, steht in
   [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
   (Home-relative Formen am Werkzeug gemessen seit slice-dcheck-v0-82-0). Das
   Grün sagt deshalb nichts über drei Flächen: über die zwei, die die Regel
   deckt und das Modul nicht (**Fence** und **Tilde mit Benutzername**), über
   die Skriptkommentare in `Makefile`, `tools/**` und `harness/mk/**`
   (**Nicht-Markdown**, außerhalb des Scans) und über **relative** Pfade (die
   Zusage lautet nicht „kein Pfad verlässt das Repo"). Das Ventil
   `hostpaths.exempt-targets` ist verfügbar und nicht gesetzt
   ([`ADR-0072`](../../docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
   Entscheidung 2). Die Reichweite der Regel und diese Lücken stehen in
   `AGENTS.md` §3.11; der Wächter dort ist das Review, kein Gate. Dieser
   Abschnitt trägt, **was der Sensor deckt und was nicht**; die Reichweite der
   Regel steht in `AGENTS.md` §3.11. Träger:
   die Reichweite
   [`ADR-0075`](../../docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md)
   (Fläche) und
   [`ADR-0160`](../../docs/plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
   (Pfad-Klassen),
   die Aktivierung
   [`ADR-0072`](../../docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md).
   Die Aktivierung kennt **keinen**
   Ausschlussblock: kein `scope`, kein `ignore`, kein `exempt-paths`.
9. **Links mit umbrochenem Linktext sind ungeprüft — `anchors` und `links`.**
   Wie das Werkzeug einen Link mit Zeilenumbruch im Linktext oder im Ziel
   entscheidet, steht in
   [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
   (gemessen in `slice-spec-festlegungen-doku-gates`, M27 bis M30). Ungedeckt
   bleibt der Link mit umbrochenem Linktext: sein Ziel kann fehlen, ohne dass
   das Grün es zeigt. Der Wächter dieser Form ist das Review, kein Gate. Erster
   Fund: Review zu `slice-077`, Delta-Review, N-2.

10. **`trace:` ist kein Modul und läuft nicht in `docs-check`/`make gates`.**
    `.d-check.yml`s `trace:`-Block konfiguriert ausschließlich die
    Requirements Traceability Matrix hinter `--trace`/`--require-complete`
    (`make doc-trace`/`make doc-complete`, `harness/README.md` §Werkzeuge) —
    dasselbe Konfigurationsprinzip wie die Opt-in-Module aus Punkt 3, aber
    ohne selbst eines zu sein: `trace.requirements.id-pattern` weicht bewusst
    von d-checks generischem Default ab und ist auf die
    `LH-(FA|QA)-[A-Z]{3}-\d{3}`-Kennungskonvention dieses Repos verdrahtet,
    `trace.coverage` liest zusätzlich `docs/user/e2e-abdeckung.md`,
    `docs/user/bench-abdeckung.md`, `docs/user/ci-matrix-abdeckung.md` und
    `docs/user/sdk-e2e-abdeckung.md` als kuratierte Coverage-Dimensionen.
    Ein grünes `make docs-check`/`make gates`
    sagt über die RTM nichts aus — sie bleibt advisory, ihr Exit-Code steht
    unabhängig neben dem Gate · seit slice-d-check-trace-rtm.

11. **Die Gate-Index-Regel misst die Länge, nicht den Satz.** Sie hält jede
    Zelle der Spalten `Vertrag` und `Tut was` in `harness/README.md`
    §Sensors (Feedback-Gates) unter ihrer Höchstlänge; ob der Satz die Zeile
    trägt, ob die Bindung die ausführliche Datei unter `harness/sensors/` bzw.
    `harness/targets/` verlinkt und ob diese Datei den Inhalt führt, prüft sie
    nicht — das bleibt Review. Die Spalte `Bindung` hat keine Höchstlänge.
    Eine neue Tabelle im Abschnitt mit einem anderen Spaltennamen fällt aus
    der Regel
    ([`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge);
    gemessen vom Verifier an einem Klon mit einer dritten Tabelle, Spalte
    `Macht was`, 300 Zeichen) · seit slice-harness-readme-zellen-kuerzen.

12. **`matrix` — zwei angenommene Lücken und der Fence.** Der Kommentar
    `<!-- d-check:status-provenance -->` hebt einen Token-Befund auch in
    `spec/` auf, auch wenn er selbst in Inline-Code steht; dort ist er kein
    zulässiger Weg, und das Grün zeigt ihn nicht. Die ADRs 0039 und 0041 stehen unter `matrix.exempt-paths` und sind
    damit auch von den Regeln `adr → slice` und `adr → review` ausgenommen.
    Ein Token in einem Fence ist für `matrix` kein Treffer. Die Wirkung steht
    in
    [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge),
    die Annahme der zwei Lücken in
    [`ADR-0163`](../../docs/plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
    Entscheidung 2 und 3; der Wächter ist das Review jeder Spec-Änderung bzw.
    jeder Zitat-Korrektur an den zwei ADRs, kein Gate.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`docker run … d-check` über `scan.roots: ["."]` mit `scan.ignore`; die
Vollständigkeits-Zeile „N Datei(en) geprüft, 0 Befund(e)“ sagt etwas über
diesen Ausschnitt (nicht über das Repo).

## Ausgabe und Ausgänge

Zeile je Befund, Schlusszeile, Exit-Codes und wie ein Exit über `make`
ankommt, stehen in
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge);
ein Beleg zitiert die Schlusszeile.

Reparatur-Pfad: `make doc-repair` (konservativ, nur `id-unlinked`/
`target-missing`, `git apply --unidiff-zero`); Diagnose: `make doc-doctor`.

## Sperren

Wann der Lauf vor der Prüfung abbricht, steht in
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
(Exit 2). Was den Weg frei macht:

- **Konfigurationsfehler** (`d-check: error: …`) → den Schlüssel bzw. die
  Quelle in `.d-check.yml` berichtigen.
- **Kein git-Repository im Mount** → aus einem Klon mit `.git` aufrufen.

## `make doc-tracked`

Isolierter Einzel-Lauf des `tracked`-Moduls (Getrackt-Status auflösbarer,
existierender Link-/Bild-Ziele gegen den git-Index): `--enable tracked`, alle
anderen Module `--disable` (`d-check.mk`). Das Modul läuft im `modules:`-Bündel
mit und damit in `make docs-check`/`make gates` (Grenze 3); dieses Ziel steht
daneben als isoliertes Diagnose-Werkzeug, kein Gate, netzlos, braucht `.git` im
Mount.

## Bindung

[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge)
(Festlegung) · `harness/conventions.md` MR-000 (ID-Schema als Linkpflicht) ·
Closure-Notiz je `done/slice-*.md` (`.d-check.yml` §structure,
`seit slice-001`) · Decken-Regel
(Baseline-Regelwerk Modul 5/6, abgebildet in `.d-check.yml` §matrix) ·
Baseline-Pin (`harness/conventions.md` §Baseline) · Register-Spalten und
Verweisform auf wandernde Slice-Pläne (`BEO-PGC/slice-pfad-als-link-in-berichten`,
`seit slice-075`) · Form der erzeugten E2E-Abdeckungstabelle
([`LH-QA-POR-003`](../../spec/lastenheft.md) — die E2E-Kette ist ihr
Erzeuger; die `structure`-Regel sichert die Zeilenform, die `ids`-Linkpflicht
die Kennungsspalte) · Zellenlänge des Gate-Index (Baseline-Regelwerk
`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt,
`seit slice-harness-readme-zellen-kuerzen`) — `.d-check.yml` §structure · kein host-lokaler absoluter
oder Home-relativer Pfad in der Doku (`hostpaths` in `modules`, ohne Ausschlussblock —
[`ADR-0072`](../../docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
trägt die Aktivierung, die Reichweite der Regel
[`ADR-0075`](../../docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md)
und
[`ADR-0160`](../../docs/plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md),
die Zitationsform eines Schwester-Repos
[`ADR-0074`](../../docs/plan/adr/0074-zitationsform-schwester-repo-hausform.md),
die Zitat-Korrektur an immutablen Dokumenten
[`ADR-0073`](../../docs/plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md))
· Getrackt-Status auflösbarer, existierender Link-/Bild-Ziele gegen den
git-Index (`tracked` in `modules`, Konfiguration `.d-check.yml` §tracked,
`exempt-targets: []` — keine eigene ADR, Präzedenzmuster Commit `f9e5a3c`
(`structure`-Modul-Aktivierung, 2026-09-09), siehe
[`docs/reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md`](../../docs/reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md)
· seit slice-d-check-tracked-modul).
