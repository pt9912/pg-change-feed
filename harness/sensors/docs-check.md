# `make docs-check` — prüft Markdown-Doku auf kaputte Referenzen (d-check)

## Vertrag

Wird dieses Target rot, trägt die gescannte Markdown-Doku eine kaputte
Referenz: lokaler Link oder Heading-Anker ins Leere (`target-missing`,
`anchor-missing`), nackte Kennung ohne Link auf ihre Definition
(`id-unlinked`), verbotene Referenzrichtung zwischen Dokumentklassen
(`matrix-forbidden` / `matrix-inactive`), abweichender Baseline-Pin
(`version-stale`), Struktur-Verstoß im Abschnitt
(`section-cell-*`, `section-forbidden` — Register-Spalten,
Closure-Notiz-Guidance, die Verweisform auf wandernde Slice-Pläne in
Berichten und der Register-Identität, die erzeugte
E2E-Abdeckungstabelle und die Zellenlänge des Gate-Index), oder einen host-lokalen absoluten Pfad in Prosa oder
Inline-Code (`hostpath-forbidden` — ein Schwester-Artefakt wird als blankes
Repo-Wort mit relativem Pfad zitiert,
[`ADR-0074`](../../docs/plan/adr/0074-zitationsform-schwester-repo-hausform.md)).
Die Module und ihre Grenzen stehen in `.d-check.yml`; die Konfiguration ist die
Deklaration dieses Vertrags, nicht dieses Dokument.

Die `structure`-Regeln prüfen Abschnitts-Invarianten in ihren Trägerdateien:
den ADR-Index, die Pflichtenheft-Defaults, zwei Architektur-Tabellen, die
Closure-Notiz je `done/slice-*.md` (sie läuft seit der ersten Closure mit,
`· seit slice-001`), die Verweisform auf wandernde Slice-Pläne in Berichten
und Register-Identität, die erzeugte E2E-Abdeckungstabelle sowie den Gate-Index
in `harness/README.md` §Sensors (Feedback-Gates): dort ist eine Zelle der
Spalte `Vertrag` höchstens 220 und eine Zelle der Spalte `Tut was` höchstens
120 Zeichen lang, beide mindestens 1; die Regel misst beide Tabellen des
Abschnitts (Gate-Tabelle und Werkzeuge-Tabelle), weil sie die Spalte über den
Kopfzeilen-Namen findet · seit slice-harness-readme-zellen-kuerzen. Innerhalb
dieser Familie adressieren die vier Tabellen-Regeln und die
E2E-Abdeckungstabelle ihre Spalten über Mindestbreiten, die Gate-Index-Regel
über Mindest- und Höchstlänge; die
Closure-Notiz-Regel und die beiden Verweisform-Regeln tragen keinen
Spalten-Knoten und prüfen nur Abschnitt und Muster
(`non-empty`/`max-open-tasks`/`require-pattern`/`forbid-pattern`). Keine
Regel zählt Zeilen.

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
   `planning/welle-NN.md` → `done/`) und der gleich-ordnerige Nachbar-Verweis
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

8. **`hostpaths` — die gescannte Markdown-Fläche, nicht das Repo.** Das Modul
   meldet host-lokale **absolute** Pfade in `.md`-Dateien unter `scan.roots`,
   in Prosa und Inline-Code (`hostpath-forbidden`). Vier benannte Ränder:
   **Fenced-Code-Blöcke** prüft es nicht — einen Opt-out-Marker kennt es
   nicht; **relative** Pfade sind ungeprüft (die Zusage lautet nicht „kein
   Pfad verlässt das Repo"); die **Windows-Laufwerks- und UNC-Muster** sind
   fest, nicht konfigurierbar; und **Nicht-Markdown** ist ungelesen — die
   Skriptkommentare in `Makefile`, `tools/**` und `harness/mk/**` erreicht der
   Scan nicht. Dateien unter `scan.ignore` (`.harness/**`, `**/*.template.md`)
   liegen ebenfalls außerhalb. **Die Regel deckt die Fenced-Fläche voll, dieses
   Modul nicht** — ihre Reichweite und diese benannte Lücke stehen in
   `AGENTS.md` §3.11; der Wächter dort ist das Review, kein Gate. Dieser
   Abschnitt trägt, **was der Sensor deckt und was nicht**; die Reichweite der
   Regel steht in `AGENTS.md` §3.11. Träger:
   die Reichweite
   [`ADR-0075`](../../docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md),
   die Aktivierung
   [`ADR-0072`](../../docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md).
   Die Aktivierung kennt **keinen**
   Ausschlussblock: kein `scope`, kein `ignore`, kein `exempt-paths`.
9. **Mehrzeilige Markdown-Links sind ungeprüft — `anchors` und `links`.** Geht
   der **Linktext** oder das **Ziel** eines Links über einen Zeilenumbruch,
   melden beide Module nichts: einzeilig gestellt kippt dieselbe Mutation rot
   (`anchor-missing` bzw. `target-missing`), zweizeilig bleibt sie grün (Exit 0).
   In `slice-077` mit sechs Mutationen gemessen; repo-weit betrifft es neun
   Links in fünf Dateien (Stand `slice-077`), ihre Ziele waren von Hand prüfbar
   und in Ordnung. Der Wächter dieser Form ist das Review, kein Gate. Träger des
   Fundes: Review zu `slice-077`, Delta-Review, N-2.

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
    der Regel (gemessen vom Verifier an einem Klon: eine dritte Tabelle mit
    Spalte `Macht was` und 300 Zeichen bleibt grün; eine umbenannte Kopfzelle
    `Tut was` meldet `section-column-missing`). Gemessen an einer Kopie im
    Scratchpad: eine Zelle `Tut was` mit 121 Zeichen und eine Zelle `Vertrag`
    mit 222 Zeichen enden mit `section-cell-oversized` (Exit 2), eine leere
    Zelle `Tut was` mit `section-cell-undersized`, eine Zelle mit 120 Zeichen
    bleibt grün · seit slice-harness-readme-zellen-kuerzen.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`docker run … d-check` über `scan.roots: ["."]` mit `scan.ignore`; die
Vollständigkeits-Zeile „N Datei(en) geprüft, 0 Befund(e)“ sagt etwas über
diesen Ausschnitt (nicht über das Repo).

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Befunde im Ausschnitt |
| 1 | mindestens ein Befund (`Datei:Zeile  Ziel  Grund-Code`) |
| 2 | Nutzungs-/Konfigurationsfehler — kein Urteil über die Doku |

Reparatur-Pfad: `make doc-repair` (konservativ, nur `id-unlinked`/
`target-missing`, `git apply --unidiff-zero`); Diagnose: `make doc-doctor`.

## Sperren

- Exit 2 bei Config-Fehler (unbekannter Schlüssel in `.d-check.yml`,
  `versions.current-from` unlesbar) — kein stiller Rückfall auf Defaults.

## Bindung

`harness/conventions.md` MR-000 (ID-Schema als Linkpflicht) · Decken-Regel
(Baseline-Regelwerk Modul 5/6, abgebildet in `.d-check.yml` §matrix) ·
Baseline-Pin (`harness/conventions.md` §Baseline) · Register-Spalten und
Verweisform auf wandernde Slice-Pläne (`BEO-PGC/slice-pfad-als-link-in-berichten`,
3×, `seit slice-075`) · Form der erzeugten E2E-Abdeckungstabelle
([`LH-QA-POR-003`](../../spec/lastenheft.md) — die E2E-Kette ist ihr
Erzeuger; die `structure`-Regel sichert die Zeilenform, die `ids`-Linkpflicht
die Kennungsspalte) · Zellenlänge des Gate-Index (Baseline-Regelwerk
`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt,
`seit slice-harness-readme-zellen-kuerzen`) — `.d-check.yml` §structure · kein host-lokaler absoluter
Pfad in der Doku (`hostpaths` in `modules`, ohne Ausschlussblock —
[`ADR-0072`](../../docs/plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
trägt die Aktivierung, die Reichweite der Regel
[`ADR-0075`](../../docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md),
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

## Fassung im Gate-Index

Ausführliche Fassung der Index-Zeile aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher. Der Text darunter ist der wortgleich umgezogene Index-Text, kein eigener Vertrag: weicht er von dieser Datei ab, gilt [§Vertrag](#vertrag) mit den Abschnitten bis zu diesem.

### `make doc-tracked`

isolierter Einzel-Lauf des `tracked`-Moduls (Getrackt-Status auflösbarer, existierender Link-/Bild-Ziele gegen den git-Index, `--enable tracked` mit allen anderen Modulen `--disable`) — das Modul läuft bereits im `modules:`-Bündel mit und damit in `make docs-check`/`make gates`; dieses Ziel bleibt daneben als isoliertes Diagnose-Werkzeug bestehen, netzlos, braucht `.git` im Mount
