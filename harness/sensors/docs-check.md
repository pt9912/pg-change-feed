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
Berichten und der Register-Identität und die erzeugte
E2E-Abdeckungstabelle). Die Module und ihre Grenzen stehen in `.d-check.yml`;
die Konfiguration ist die Deklaration dieses Vertrags, nicht dieses Dokument.

Die `structure`-Regeln prüfen Abschnitts-Invarianten in ihren Trägerdateien:
den ADR-Index, die Pflichtenheft-Defaults, zwei Architektur-Tabellen, die
Closure-Notiz je `done/slice-*.md` (sie läuft seit der ersten Closure mit,
`· seit slice-001`), die Verweisform auf wandernde Slice-Pläne in Berichten
und Register-Identität sowie die erzeugte E2E-Abdeckungstabelle. Innerhalb
dieser Familie adressieren die vier Tabellen-Regeln und die
E2E-Abdeckungstabelle ihre Spalten über Mindestbreiten; die
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
3. **Opt-in-Module nicht im Bündel** — `planning`, `tracked`, `vcs`,
   `commits`, `reviews` laufen nur über ihre `doc-*`-Einzel-Targets mit
   `--enable`; `make gates` belegt sie nicht. Heilbar je Aktivierungs-
   bedingung (in `.d-check.yml` kommentiert).
4. **`MR-*` nicht linkpflichtig** — das `ids`-Muster deckt LH/SPEC/ARC/ADR,
   nicht MR; Adaptions-Verweise werden nur vom `tracked`-Modul geprüft
   (opt-in). Permanent bis zur Muster-Erweiterung.
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
     `-run`-Muster.

   Die tragende Garantie der Tabelle ist deshalb die **Ableitung aus dem
   Quelltext** (Go-Hälfte) und der **Anker je Phase** (Bash-Hälfte) — nicht
   das Doku-Gate.

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
die Kennungsspalte) — `.d-check.yml` §structure.
