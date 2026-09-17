# Welle welle-d-check: `d-check`-Erweiterung — Getrackt-Status und Requirements-Traceability-Matrix

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-d-check-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Namensform** nach
[`harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md`](../../../harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md)
— die erste **namensbasierte** Welle dieses Repos (`welle-1`…`welle-20` bleiben
nummeriert, Bestandsschutz).

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** —. **Datum:** 2026-09-17.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

`d-check` (bereits vendored gepinnt, `d-check.mk`, `DCHECK_DIGEST =
sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`,
Tag `v0.75.0`) wird um **zwei unabhängige, additive Fähigkeiten** erweitert,
die beide in derselben `.d-check.yml` verankert werden:

1. **`tracked`-Modul aktiviert** — jedes auflösbare, existierende
   Link-/Bild-Datei-Ziel muss im git-Index getrackt sein (`target-untracked`
   fängt eine Umgebungs-Drift, die beim Erzeuger grün bleibt, auf einem
   frischen Klon aber `target-missing` wäre). Konkrete Kandidaten:
   `docs/user/e2e-abdeckung.md` (generiert von `make test-integration`, aber
   committet) und `gen/cdc/stream/v1/*.pb.go`/`*_test.go` (generiert von
   `make proto-generate`, committet) — beide vielfach verlinkt, beide könnten
   künftig durch eine zu breite `.gitignore`-Regel unbemerkt herausfallen.
2. **`--trace`/Requirements Traceability Matrix inkl. `trace.coverage`** —
   `spec/lastenheft.md` nutzt bereits die von `--trace` erwartete
   ATX-Heading-Form (`### LH-FA-CFG-001 — …`); das `id-pattern`
   `'LH-(FA|QA)-[A-Z]{3}-\d{3}'` passt **geprüft** gegen alle 76 Anforderungs-
   Überschriften (`grep -oE '^### LH-[A-Za-z0-9.-]+' spec/lastenheft.md`, 0
   Nicht-Treffer). `trace.coverage` liest zusätzlich
   `docs/user/e2e-abdeckung.md` als kuratierte Coverage-Dimension, damit
   E2E-gedeckte, aber ADR-/Slice-lose Anforderungen nicht fälschlich als
   Waise erscheinen.

**Das Mehr, das die Welle rechtfertigt** (statt zwei unabhängiger,
wellenloser Slices): ein **kombinierter** `make gates`-Lauf, in dem **beide**
neuen Fähigkeiten gleichzeitig aktiv sind (`tracked` im `modules:`-Bündel,
`trace`-Block in derselben Datei) — das prüft eine Interaktion, die keine der
beiden Einzel-DoDs allein abdeckt: eine mögliche Config-Kollision in
derselben `.d-check.yml` (Schema-Validierung, YAML-Struktur,
Modul-Reihenfolge) wird nur sichtbar, wenn beide Blöcke gleichzeitig geladen
werden.

**Real gemessen (Planungslauf 2026-09-17, gepinnter Digest
`sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`,
gegen eine Kopie des Bestands, `--network none`):**

- `tracked` allein (`--enable tracked`, alle anderen Module `--disable`):
  `872 Datei(en) geprüft, 0 Befund(e)` — beide genannten Kandidatendateien
  sind heute getrackt, keine Vorab-Sanierung nötig.
- `--trace` ohne `trace.coverage`: `76 Anforderung(en), 9 Waise(n)`.
- `--trace` mit `trace.coverage` auf `docs/user/e2e-abdeckung.md`:
  `76 Anforderung(en), 7 Waise(n)` — zwei Anforderungen sind ausschließlich
  über die E2E-Tabelle gedeckt und wären ohne den Coverage-Block falsche
  Waisen.
- `--trace --require-complete` (ohne `modality`-Filter, damit gatet **jede**
  Waise): **Exit 1**, 7 verbleibende Waisen
  (`LH-FA-CFG-006`, `LH-FA-CON-002`, `LH-FA-DAT-002`, `LH-FA-DAT-003`,
  `LH-FA-SST-001`, `LH-FA-SST-005`, `LH-QA-REL-004`) — der reale Beleg dafür,
  dass `doc-complete` **heute** scharf geschaltet sofort `make gates` rot
  schießen würde, ohne dass diese sieben je geprüft wurden. Trägt §Out-of-Scope
  unten.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Der gepinnte `d-check`-Digest (`v0.75.0`,
  `sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`)
  unterstützt beide Fähigkeiten bereits — **erfüllt**, real geprüft: `tracked`
  seit `v0.37.0`, `--trace` seit `v0.22.0`, `trace.coverage` seit `v0.41.0`
  (Änderungshistorie, Benutzerhandbuch §11), alle drei weit unterhalb des
  gepinnten Standes. Kein Versions-Sprung nötig.
- `welle-20` liegt in `done/` — **erfüllt**; `in-progress/` ist zum Zeitpunkt
  der Eröffnung leer (WIP-Limit frei, geprüft 2026-09-17).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Beide Slices (`slice-d-check-tracked-modul`, `slice-d-check-trace-rtm`)
  liegen in `done/`.
- **Ein kombinierter `make gates`-Lauf ist grün**, während `.d-check.yml`
  gleichzeitig den `tracked`-Eintrag in `modules:` **und** den vollständigen
  `trace:`-Block (inkl. `trace.coverage`) trägt — das *Mehr* dieser Welle.
- `make doc-trace` läuft gegen denselben Stand ohne Konfigurationsfehler
  (Exit 0, advisory) und zeigt die reale RTM mit Coverage-Spalte.
- Closure-Notiz in `welle-d-check-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| `slice-d-check-tracked-modul` | `tracked`-Modul aktivieren (Getrackt-Status von Link-/Bild-Zielen) | — (Konfigurationsstruktur, kein Vertrags-Bezug) |
| `slice-d-check-trace-rtm` | `--trace`/RTM inkl. `trace.coverage` für `docs/user/e2e-abdeckung.md` verdrahten | — (Konfigurationsstruktur, kein Vertrags-Bezug) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Blockiert:** nichts.
- **Wird blockiert von:** nichts — beide Fähigkeiten sind im bereits
  gepinnten Digest vorhanden (siehe §2).
- **Beide Slices teilen sich eine Datei** (`.d-check.yml`) — das ist die
  eigentliche Abhängigkeit dieser Welle: Reihenfolge der beiden Slices ist
  beliebig, aber der zweite committende Slice sieht den Config-Stand des
  ersten und muss additiv bleiben (kein Überschreiben fremder Blöcke).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **`doc-complete`/`--require-complete` als neues Gate in `GATE_CHECKS`.**
  *Es wäre ein anderer Vorgang* — real gemessen (§1) meldet der Bestand
  **heute** 7 reale Waisen; sie scharf zu schalten bräche `make gates` sofort,
  ohne dass eine dieser sieben Anforderungen bislang bewusst geprüft wurde.
  `doc-trace` bleibt advisory (wie `make image-stale`); ob und wie die sieben
  Waisen aufgelöst werden (echte Lücke, fehlende ADR-/Slice-Zitierung, oder
  ein weiteres `trace.coverage`-Ziel), ist ein **eigener** Folge-Vorgang.
- **Eine dritte, in dieser Welle nicht genannte `d-check`-Fähigkeit**
  (`trace.modality`, `trace.cross-consistency`, `pins`, `sources`, …). *Es
  wäre ein anderer Vorgang* — diese Welle liefert genau die zwei im
  Nutzer-Dialog recherchierten Fähigkeiten, keine Erkundung des restlichen
  `d-check`-Funktionsumfangs.
- **Eine Versionsanhebung des gepinnten `d-check`-Digests.** *Bestand bleibt
  bewusst stehen* — §2 belegt, dass der bereits gepinnte Stand (`v0.75.0`)
  beide Fähigkeiten trägt; ein Versions-Sprung wäre ein unabhängiger,
  unbegründeter Zusatzschritt mit eigenem Diff-Risiko (Änderungshistorie
  zwischen `v0.75.0` und der Handbuch-Referenzversion `v0.76.3` nennt u. a.
  `vcs`/`commits`-Verhaltensänderungen, die diese Welle nicht braucht und
  nicht bewertet).
- **Auflösung von `BEO-PGC/regel-weiter-als-ihr-sensor`** (3×, Schwelle
  erreicht, Ausgang seit `welle-20`s Schließung noch nicht zugewiesen — siehe
  §8 Sichtungs-Hinweis unten). *Es wäre ein anderer Vorgang*: Der fällige
  Lese-Schritt (Modul 6 Closure-Schritt 3a) gehört in die **Closure** dieser
  Welle, nicht in einen ihrer Liefer-Slices — er bewegt keinen der beiden
  hier gelieferten Umfänge und ist kein Liefer-Punkt.

## 7. Closure-Notiz

<!--
BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md §Verwendung,
Schritt 5) und darf deshalb nichts Tragendes halten. Erst nach
Welle-Abschluss fuellen — bis dahin bleiben die beiden Zeiger unten
Platzhalter, kein realer Link (sonst meldet `make docs-check` schon jetzt
`target-missing`, da `welle-d-check-results.md` erst bei Closure entsteht).
-->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: [`welle-d-check-results.md`](welle-d-check-results.md)
Zähler: [`../observations/`](../observations/)
