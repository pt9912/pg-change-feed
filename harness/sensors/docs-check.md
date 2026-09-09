# `make docs-check` — prüft Markdown-Doku auf kaputte Referenzen (d-check)

## Vertrag

Wird dieses Target rot, trägt die gescannte Markdown-Doku eine kaputte
Referenz: lokaler Link oder Heading-Anker ins Leere (`target-missing`,
`anchor-missing`), nackte Kennung ohne Link auf ihre Definition
(`id-unlinked`), verbotene Referenzrichtung zwischen Dokumentklassen
(`matrix-forbidden` / `matrix-inactive`), abweichender Baseline-Pin
(`version-stale`), Struktur-Verstoß in Register-Spalten
(`section-cell-*`). Die Module und ihre Grenzen stehen in `.d-check.yml`;
die Konfiguration ist die Deklaration dieses Vertrags, nicht dieses Dokument.

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
Baseline-Pin (`harness/conventions.md` §Baseline) · Register-Spalten
(`.d-check.yml` §structure, fünfte Regel auskommentiert bis zur ersten
Closure).
