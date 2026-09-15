# Review-Report: slice-074 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-074`, Diff-Range `c553ddc..972851e` — die sechs Commits
`8e1240f`, `7e006ae`, `3f07ac8`, `c1ae6da`, `399370e`, `972851e`. `c553ddc`
ist der reine `next → in-progress`-Move und ausdrücklich **nicht** Teil des
Gegenstands; er ist der Elter von `8e1240f`. Berührt sind genau sechs Pfade
(`git diff --stat c553ddc..972851e`): `test/integration/integration_test.go` ·
`tools/harness/run-integration-tests.sh` · `docs/user/e2e-abdeckung.md` ·
`.d-check.yml` · `harness/sensors/docs-check.md` · der Slice-Plan (§2/§3).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14 — die Regel „Neue Betreiber-Oberfläche ohne Handbuch-Zug" liegt
**vor** diesem Implementer-Lauf; auf diesen Diff kommt keine der vier
repo-spezifischen HIGH-Regeln zur Anwendung).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-074-e2e-abdeckungstabelle.md`
  vollständig (§1–§8, einschließlich des §3-Nachzugs aus `8e1240f`/`399370e`)
- [`LH-QA-POR-003`](../../spec/lastenheft.md) (der Bezug des Slice) und die von
  der Tabelle adressierten Kennungen [`LH-QA-SEC-001`](../../spec/lastenheft.md),
  [`LH-QA-OPS-005`](../../spec/lastenheft.md),
  [`LH-QA-REL-001`](../../spec/lastenheft.md)
- `.d-check.yml` (alle acht `structure`-Regeln, `ids`-Muster, `modules`,
  `scan`), `d-check.mk` (gepinntes Image), `harness/sensors/docs-check.md`,
  `harness/mk/doc-gate.mk`
- `AGENTS.md` §3 (Hard Rules, darunter §3.1 Docker-only, §3.6, §3.7, §3.9) und
  §5/§6 · `harness/conventions.md` (MR-000 ID-Schema)
- Beobachtungs-Register (gemergter Stand): `BEO-PGC/generierte-artefakte-ohne-sync-sensor`
  (1×, `*`/PGC, offen) · `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`
  (1×, `*`/PGC, offen) · `BEO-PGC/test-runner-stiller-ausschluss` (1×, offen) ·
  `BEO-PGC/lese-doppelquelle` · `BEO-PGC/d-migrate-nacharbeit` ·
  `BEO-PGC/slice-chronik-in-code-kommentar`
- `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung und §Offene Risiken werden bei Closure aufgelöst ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register ·
  `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz ·
  `regelwerk/modul-10-review-harness.md` §Drei Review-Arten
- Fremdquelle zum Gate: gepinntes Image
  `ghcr.io/pt9912/d-check@sha256:18e9cd857f8db3569526d1f9a3cbeba8af51e9f2dd84c17a22444028b977c3da`
  — als **Prüfgegenstand** verwendet (eigene Stichproben, siehe §Messungen),
  nicht als Modulbeschreibung

---

## Findings

### F-1 — Die neue Prosa über die `structure`-Regeln trifft für drei der acht Regeln nicht zu; die Kennungs-Mindestbreite der neuen Regel fängt unverlinkte Kennungen gemessen nur bis zwei Stück

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist") ·
  Maintainability · eigene Messung am gepinnten Image
- `pfad`: `harness/sensors/docs-check.md:17-23` (bes. `:21-23`) ·
  `.d-check.yml:180-182`
- `befund`: Die neue Zusammenfassung schließt mit „… sowie die erzeugte
  E2E-Abdeckungstabelle — **jede** über ihren Abschnitt und ihre
  **Spalten-Mindestbreiten**, keine über eine Zeilenzahl". Die fünfte Regel
  (Closure-Notiz je `docs/plan/planning/done/slice-*.md`) und die sechste und
  siebte (Verweisform auf wandernde Slice-Pläne in `docs/reviews/**` und
  `observations/**/observation.md`) tragen **keinen** `table`-Knoten: sie
  arbeiten mit `non-empty`, `max-open-tasks`, `require-pattern` und
  `forbid-pattern`. Die Aussage beschreibt damit drei der acht Regeln falsch.
  Die zweite Aussage (`.d-check.yml:180-182`, „eine Kennungsspalte ohne Link
  auf ihr Definitions-Dokument bleibt darunter") ist an ihrer eigenen
  Zahlengrundlage **gemessen unwahr**: die Mindestbreite ist 40 Zeichen, eine
  Zelle mit **drei** unverlinkten Kennungen ist 43 Zeichen lang und löst
  keinen Befund aus (siehe §Messungen, Probe C). Der tragende Wächter für
  diesen Fall ist das `ids`-Modul (`id-unlinked`), nicht die Mindestbreite —
  die Aussage schreibt ihm eine Wirkung zu, die er nur für bis zu zwei
  Kennungen hat. Beide Sätze sind mit diesem Diff neu.
- `verifizierbar`: ja — `.d-check.yml` §structure gegen die zitierten
  Schlüssel lesen; die zweite Hälfte ist am gepinnten Image reproduzierbar
  (`.d-check.yml` mit `modules: [structure]` über eine Zelle mit drei
  unverlinkten Kennungen).
- `klasse`: „Zusage über die `structure`-Regelfamilie überdeckt den Config"
  (1× — im Skill-Katalog noch keine Klasse; betrifft die Prosa über eine
  Gate-Konfiguration, nicht ihr Verhalten)

### F-2 — Der §8-Sichtungs-Schritt lässt den Register-Eintrag aus, der genau dieses Erzeugnis führt; §7 erwartet für ihn keinen Beleg

- `kategorie`: MEDIUM
- `quelle`: `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Zwei Schritte
  vor der Modus-Begründung (Sichtungs-Schritt, in jedem Slice-Plan
  unbedingt) · `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register
- `pfad`: `docs/plan/planning/in-progress/slice-074-e2e-abdeckungstabelle.md:442`
  (§8-Schluss „Kein weiterer Eintrag des Registers berührt die Sub-Area dieses
  Slice") · `:406-445` (§8-Sichtungsliste) · `:365-371` (§7-Erwartung zum
  Beobachtungs-Register) · `:43-53` (§Ziel) ·
  `docs/plan/planning/observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/state.md`
- `befund`: Das Register führt unter demselben Sub-Area-Kürzel `PGC` den
  Eintrag `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (1×, offen) mit
  wörtlich dieser Beobachtung: „Ein Teil des Repos besteht aus **erzeugten
  Artefakten, die committet im Baum liegen** … Für keines davon existiert ein
  Sensor, der das committete Artefakt gegen seine Quelle hält … und `make
  gates` bliebe grün." Der Diff legt genau ein solches Artefakt an; die neue
  `structure`-Regel prüft seine Zeilenform, nicht seinen Abgleich gegen den
  Quelltext. Der Eintrag steht weder in der §8-Liste noch in der
  §7-Erwartung (die nur `BEO-PGC/test-runner-stiller-ausschluss` als zweiten
  Träger nennt), und die §8-Schlusszeile behauptet seine Nicht-Berührung.
  Dieselbe Lücke trägt die §Ziel-Zusage „per Konstruktion": sie deckt die
  Richtung *Quelle geändert, Erzeuger nicht gelaufen* nicht — eine veraltete,
  aber wohlgeformte Tabelle bleibt in `make docs-check` und in `make gates`
  grün (die `section-missing`-Kopplung sichert nur die **Existenz**).
- `verifizierbar`: ja — `ls docs/plan/planning/observations/BEO-PGC/` gegen die
  §8-Liste lesen; die Grün-Bleibt-Aussage ist über `.d-check.yml` §structure
  (Kopfzeilen + Mindestbreiten, keine Quelltext-Bindung) zeichengleich
  nachvollziehbar.
- `klasse`: „Sichtungs-Schritt lässt einen Register-Eintrag der eigenen
  Sub-Area aus" (1× — Nachbar-Klasse zu `BEO-PGC/generierte-artefakte-ohne-sync-sensor`,
  von der dieser Slice ein neuer Träger wäre)

### F-3 — Die Beschreibungsspalte ist nicht kennungsfrei: Register-Pfad und Lebenszyklus-Bezug überleben die Normierung

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 · Plan §Ziel/§8
  („die Tabelle trägt **keine** Slice-/Wellen-Chronik") gegen
  `harness/sensors/docs-check.md:71-76` („die Beschreibungsspalte … trägt
  deshalb **keine** Kennungen")
- `pfad`: `docs/user/e2e-abdeckung.md:17` (`BEO-PGC/lese-doppelquelle`) ·
  `:27` („im Slice-Plan (§6) benannten Fällen") ·
  `test/integration/integration_test.go:1213`
  (`abdeckungKennungMuster`) · `tools/harness/run-integration-tests.sh`
  (`abdeckungKopf`)
- `befund`: Der Erzeuger entfernt `LH-*`, `SPEC-*`, `ARC-*`, `ADR-*`,
  `slice-NNN`, `welle-NN` samt umgebendem Inline-Code-Span — die Kennungsform
  `BEO-<KUERZEL>/<slug>` aus MR-000 liegt außerhalb des Musters, und die
  Wortform „Slice-Plan" ist kein Kennungs-Treffer. Beide stehen in der
  erzeugten Tabelle in der vierten Spalte. Der Nachbarfix `399370e` begründet
  die Entfernung für genau diese Klasse (Provenienz der Testfälle:
  `slice-030`) — die Normierung ist damit nach ihrer eigenen Begründung nicht
  vollständig. Kein Gate fängt das: `ids` kennt `BEO-*` nicht als Muster.
- `verifizierbar`: ja — die Spalte 4 der Datei gegen das Muster in
  `integration_test.go:1213` prüfen.
- `klasse`: „Normierungs-Zusage deckt eine Provenienz-Form nicht" (1×)

### F-4 — Die Stabilitäts-Zusage im Kopf nennt eine der beiden Änderungsursachen der Datei

- `kategorie`: LOW
- `quelle`: Maintainability · Slice-Plan §6 Risiko 4 gegen die Kopf-Zeile
- `pfad`: `tools/harness/run-integration-tests.sh:84-95` (erzeugt den Kopf) ·
  `docs/user/e2e-abdeckung.md:8-10` · Slice-Plan `:332-335` (§6 Risiko 4)
- `befund`: Der Kopf sagt „sie ändert sich mit den Nachweis-Deklarationen,
  nicht mit jedem Lauf". §6 Risiko 4 desselben Plans nennt die zweite
  Ursache: jede Einfügung oberhalb einer `func TestE2E*` **oder** oberhalb
  einer `abdeckung_declare`-Zeile verschiebt die `Ort`-Zelle, die Datei
  ändert sich also auch ohne geänderte Nachweis-Deklaration (37 von 37 Zeilen
  binden an `Datei:Zeile`). Die DoD-Zeile (b) bleibt davon unberührt und ist
  ehrlich — ein zweiter Lauf bei unverändertem Testbestand liefert
  byte-gleiche Ausgabe (§Messungen) —, der Kopf unterschlägt aber den im
  eigenen §6 benannten Fall.
- `verifizierbar`: ja — die Zeilenbindung steht in jeder Tabellenzeile; jede
  Einfügung oberhalb verschiebt sie.
- `klasse`: „Stabilitäts-Zusage enumeriert ihre Ursachen unvollständig" (1×)

### F-5 — Die Go-Hälfte liest das Verzeichnis der Testdatei, der Runner führt über `./test/integration/...` (rekursiv) aus

- `kategorie`: LOW
- `quelle`: Beobachtungs-Register `BEO-PGC/test-runner-stiller-ausschluss`
  (dieselbe Klasse, andere Richtung) · Maintainability
- `pfad`: `test/integration/integration_test.go:1284-1294`
  (`abdeckungsZeilen`/`os.ReadDir(verzeichnis)`) ·
  `tools/harness/run-integration-tests.sh:153`, `:396`
- `befund`: Beide `go test`-Aufrufe des Runners zielen auf
  `./test/integration/...` und erfassen damit auch Unterpakete; die Ableitung
  liest nur die Geschwisterdateien des Verzeichnisses, in dem die
  Erzeuger-Datei liegt. Eine `func TestE2E*` in einem Unterpaket würde real
  ausgeführt und erschiene in keiner Tabellenzeile — dieselbe stille
  Ausschluss-Klasse, die der Slice auf der `-run`-Seite ausdrücklich offen
  lässt (§1), hier auf der Ableitungs-Seite. Heute kein Fall (das Paket
  besteht aus einer Datei); die Grenze ist nirgends benannt.
- `verifizierbar`: ja — `grep -rn "^func TestE2E" test/integration/` zeigt
  heute eine Datei; die Scope-Differenz ist zeichengleich lesbar.
- `klasse`: „Ableitungs-Scope enger als der Ausführungs-Scope" (1×)

### F-6 — Das neue Erzeugnis hat keinen eingehenden Verweis aus einem lesenden Knoten

- `kategorie`: INFO
- `quelle`: Maintainability — zuständig: Planner (`docs/user/*` ist
  Rang-6-Quelle in `harness/README.md` §Source precedence)
- `pfad`: `README.md:24` (verweist auf `docs/user/benutzerhandbuch.md`) ·
  `harness/README.md:28` (`docs/user/*` entlinkt) · `docs/user/`
- `befund`: `docs/user/e2e-abdeckung.md` ist aus keiner Doku verlinkt: nicht
  aus `README.md`, nicht aus dem Benutzerhandbuch, nicht aus
  `harness/README.md`, und `docs/user/` führt keinen Index. Ein Modul, das
  verwaiste Dateien meldet, ist nicht aktiv (`modules: [links, anchors, ids,
  matrix, versions, structure]`). Der Plan schließt nur den
  §Sensors-Eintrag aus („kein neues Target, kein neues Gate") — über einen
  lesenden Kanten-Zug sagt er nichts.
- `verifizierbar`: ja — `grep -rn "e2e-abdeckung" --include=*.md .` außerhalb
  der Datei selbst und des Plans.
- `klasse`: „Erzeugnis ohne eingehenden Verweis" (1×)

### F-7 — Der Deklarations-Anker prüft das Vorkommen, nicht die Lage

- `kategorie`: INFO
- `quelle`: Slice-Plan §1 („mit einem wörtlichen Anker aus der Phase") ·
  Slice-Plan §6 Risiko 1
- `pfad`: `tools/harness/run-integration-tests.sh:99-116`
  (`abdeckung_declare`, `awk 'NR > ab && index($0, anker)…'`)
- `befund`: Gesucht wird die **erste** Zeile hinter dem Deklarations-Aufruf,
  die den Ankertext enthält — nicht die Phase selbst. Ein Kommentar oberhalb
  der Phase, der den Ankertext wiederholt, verschiebt die `Ort`-Zelle still;
  die Prüfung belegt Existenz, nicht Identität. Zweite, heute falllose
  Beobachtung derselben Stelle: der Anker geht als `awk -v`-Wert durch die
  Escape-Verarbeitung von `awk` — ein Backslash im Anker wäre dort ein
  zweiter Sonderfall (heute trägt kein Anker einen). Alle 24 Anker lösen
  heute auf die abschließende Echo-Zeile ihrer Phase auf (§Messungen).
- `verifizierbar`: nein — kein Gate; die Auflösung ist nachrechenbar (siehe
  §Messungen).
- `klasse`: „Anker-Prüfung prüft Vorkommen, nicht Lage" (1×)

## Negativbefunde

- geprüft, ohne Befund: **`tools/harness/run-integration-tests.sh` — die
  Kommando-Substitutions-Falle ist vollständig geschlossen.** Alle 24
  Deklarations-Aufrufe sind auf ihren Shell-Metazeichen-Gehalt geprüft: kein
  `$`, kein Backslash; genau ein Aufruf (`SSE-Stream-Rundlauf`, `:2140`) trägt
  Backticks und beide sind escaped (`\``) — die erzeugte Zelle trägt den
  Inline-Code-Span (Zeile 49 der Tabelle). `ABDECKUNG_KOPF` ist
  single-quoted, das `printf`-Format von `abdeckung_render` ebenfalls; die
  einzigen doppelt gequoteten Fremdtexte sind die Deklarations-Argumente, und
  sie sind Autorentext, kein abgeleiteter Text. **Abgeleiteter Text** (der
  Doc-Kommentar des Go-Pakets) gelangt ausschließlich als Daten in
  Variablen/`read`/`printf` — es gibt kein `eval`, kein `sh -c`, keine
  Re-Expansion einer Variable. Kein Erzeuger-Pfad führt Fremdtext als Code
  aus.
- geprüft, ohne Befund: **`docs/user/e2e-abdeckung.md` reproduziert einen
  frischen Lauf byte-gleich** (§Messungen) — Go-Hälfte real im
  Toolchain-Container, Bash-Hälfte und Rendering nachgerechnet.
- geprüft, ohne Befund: **`.d-check.yml` §structure** — die achte Regel ist
  live (Abschnitts-Mismatch ⇒ `section-missing`), die
  `section-missing`-Kopplung bei fehlender Datei ist real (Exit 1), der
  Regel-Kommentar benennt die Nicht-Prüfung „falscher Nachweis" korrekt, und
  `docs/user/e2e-abdeckung.md` ist bei `make docs-check` grün.
- geprüft, ohne Befund: **`test/integration/integration_test.go`** — die
  Ableitung liest den Quelltext (`go/parser` über die AST-Deklarationen),
  nicht den Lauf; der Abbruch-Wächter ist real gemessen (§Messungen). Kein
  `//nolint`, kein Produktionspfad berührt.
- geprüft, ohne Befund: **die §1-Abgrenzungen des Plans halten** —
  `git diff --name-only c553ddc..972851e` führt weder `harness/README.md`
  noch `AGENTS.md` noch `d-check.mk`; `internal/**` und `tools/schema/**`
  sind unberührt; `codepaths` ist in `.d-check.yml:210` weiterhin
  auskommentiert; kein neues Skript, kein neues Binary, keine Lauf-Spalte.
- geprüft, ohne Befund: **`AGENTS.md` §3.7 in den neuen Kommentaren** — kein
  Slice-/Wellen-Chronik-Satz, kein „früher stand hier", kein abwesender Text
  und kein abgebrochener Satz in den neuen Kommentaren von
  `integration_test.go`, im Runner oder in `harness/sensors/docs-check.md`;
  die Herkunfts-Anker (`· seit slice-001`, `· seit slice-075`) stehen in der
  zulässigen Form.
- geprüft, ohne Befund: **`harness/sensors/docs-check.md` jenseits von F-1** —
  die drei gemessenen Grenzen (`ids` nur Link, Code-Spans ungeprüft,
  `--trace` unabhängig von der Datei) sind unabhängig nachgestellt und
  bestätigt (§Messungen).
- geprüft, ohne Befund: **`docs/plan/planning/in-progress/slice-074-e2e-abdeckungstabelle.md`
  §8-Registerstand** für `BEO-PGC/test-runner-stiller-ausschluss` — 1×,
  `evidence/slice-033.md`, Zustand `offen`; die Plan-Angabe stimmt. Die
  Sichtung des Eintrags `generierte-artefakte-ohne-sync-sensor` fehlt (F-2).

## Messungen dieses Laufs (jeweils mit Exit-Code)

| Messung | Kommando | Ergebnis |
|---|---|---|
| Doc-Gate | `make docs-check` (Exit direkt, ungepiped) | Exit `0` — „600 Datei(en) geprüft, 0 Befund(e)" |
| Inneres Gate-Bündel | `make gates` (Exit direkt, ungepiped) | Exit `0` (coverage 49,30 % ≥ 35 %, commit-traceability `HEAD~5..HEAD` OK, a-check „gesamt: 0 Befund(e)") |
| Race-Suite | `make test` | Exit `0`, `test/integration` ok unter `-race` |
| Erzeuger, Go-Hälfte | `docker run … go test -v -run '^TestAbdeckungstabelleZeilen$' ./test/integration/...` | Exit `0`, 13 `ABDECKUNG|`-Zeilen (13 `func TestE2E*` im Paket) |
| Mutation 1 — `func TestE2E*` ohne Spec-Kennung (Doc-Kommentar von `TestE2ERetentionBlockersViewShowsFurthestBehindConsumer`, zwei Kennungen entfernt) | derselbe Aufruf | Exit `1` — „keine Spec-Kennung im Doc-Kommentar — der Nachweis wäre keiner Anforderung zugeordnet"; Tabelle byte-identisch (sha `854cd778…` vor und nach); Mutation per `git checkout` zurückgenommen, `git status --porcelain` leer |
| Bash-Hälfte, unabhängig nachgerechnet | 24 Deklarationen geparst, Anker wie `awk` gesucht, Zeilen gegen die Datei gehalten | alle 24 `Ort`-Zeilen identisch |
| Erzeugnis gegen frischen Lauf | Go-Hälfte real + Bash-Hälfte + Rendering nachgebaut, sha verglichen | sha gleich (`854cd778f0ea9c0b64aa088974a66e8b6b6302acdf537774865d72be938b2364`) — **IDENTISCH** |
| Probe A — Datei fehlt | `d-check` @ `sha256:18e9…` über `.d-check.yml` mit `modules: [structure]`, Ziel fehlt | Exit `1`, `section-missing` — „Regel trifft keine Datei" |
| Probe B — Abschnitt falsch | dito, Datei vorhanden, Selektor umbenannt | Exit `1`, `section-missing` — „kein Abschnitt passt auf den Selektor" |
| Probe C — drei unverlinkte Kennungen in der Kennungsspalte | dito, Zelle 43 Zeichen ≥ `cell-min-chars: 40` | Exit `1` — aber **nur** `section-cell-undersized` der Nachweis-Spalte; die Kennungsspalte löste keinen Befund aus (Grundlage von F-1) |
| Probe D — `ids`-Grenzen | dito, `modules: [ids]`, verlinkte **erfundene** Kennung / Kennung im Code-Span / nackte Kennung | erfundene (verlinkte) Kennung **grün**; Code-Span **grün**; nackte Kennung ⇒ `id-unlinked`, Exit `1` — beide Grenzen aus `harness/sensors/docs-check.md` bestätigt |
| Probe E — `make doc-trace` mit/ohne die neue Datei | zweimal gelaufen, Datei dafür kurz entfernt und zurückgelegt | `diff` der beiden Ausgaben leer (Exit `0`) — die Tabelle speist die RTM nicht |

Wegwerf-Dateien lagen außerhalb des Arbeitsbaums (`/tmp/`); nach allen
Mutationen ist `git status --porcelain` leer.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Zusage über die `structure`-Regelfamilie
überdeckt den Config" (F-1, 1×) · „Sichtungs-Schritt lässt einen
Register-Eintrag der eigenen Sub-Area aus" (F-2, 1×) · „Normierungs-Zusage
deckt eine Provenienz-Form nicht" (F-3, 1×) · „Stabilitäts-Zusage enumeriert
ihre Ursachen unvollständig" (F-4, 1×) · „Ableitungs-Scope enger als der
Ausführungs-Scope" (F-5, 1×) · „Erzeugnis ohne eingehenden Verweis" (F-6, 1×)
· „Anker-Prüfung prüft Vorkommen, nicht Lage" (F-7, 1×).

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. Die beiden MEDIUM liegen **außerhalb
des Erzeugnisses**: F-1 sind zwei Sätze Prosa über eine Gate-Konfiguration,
F-2 ist eine Sichtungs-/Erwartungs-Zeile des Slice-Plans; die Tabelle selbst
ist byte-gleich zu einem frischen Lauf, alle Sensoren sind grün, und keine
der beiden Korrekturen ändert Verhalten. Die Abweichung von der Regel „HIGH
und MEDIUM blockieren typischerweise" wird damit begründet, nicht still
entschieden.

**Der Gegenstand trägt.** Die drei tragenden Zusagen des Plans sind gegen den
Diff geprüft und halten: (1) Die Go-Hälfte **leitet** aus dem Quelltext ab —
`go/parser` über die AST-Deklarationen des Pakets, nicht aus dem, was
gelaufen ist; eine Mutation, die einer Testfunktion die Spec-Kennung nimmt,
beendet den Erzeuger sichtbar (Exit `1`) und lässt die Datei unangetastet.
(2) Die Bash-Hälfte ist an **jeder** Phase deklariert; alle 24 Anker lösen
heute auf, und keine Phase der Kette bleibt ohne Zeile. (3) Geschrieben wird
nur bei Abweichung, und ein Abbruch fasst die Datei nicht an — der einzige
Schreibpfad (`abdeckung_schreiben`, `:167-188`) steht als letzte Anweisung
des Runners, alle Abbrüche liegen davor; die Idempotenz-Zusage ist
byte-genau nachgerechnet. Die Kommando-Substitutions-Falle des
SSE-Aufrufs ist geschlossen, und der Fix ist **vollständig**: geprüft sind
alle Deklarations-Argumente (kein `$`, kein Backslash, ein Backtick-Paar,
escaped), der single-quoted Kopf, das `printf`-Format und sämtliche Pfade, auf
denen abgeleiteter Text fließt — er fließt ausschließlich als **Daten**.

**Zu den vier §6-Risiko-Einschätzungen des Implementers:**

1. **R1 „Bash-Hälfte deklariert, nicht abgeleitet" — weiter offen: bestätigt,
   ohne Einschränkung.** Die Anker-Prüfung deckt die Vorwärtsrichtung
   (Vorkommen) und ist selbst lage-blind (F-7); die Rückwärtsrichtung — eine
   neue Phase, die niemand deklariert — hat keinen Wächter. Für den heutigen
   Stand ist die Rückwärtsrichtung *nicht* eingetreten: die 24 Deklarationen
   decken jede Phase mit eigener Abschluss-Zeile; die beiden nicht separat
   repräsentierten Echo-Zeilen (`:935`, `:952`) gehören zur kombinierten
   Retention-Phase, die die Zeile zu `:967` trägt.
2. **R2 „Normalisierung eingetreten, benannt" — bestätigt, und die
   Reichweite ist größer als benannt.** Die Fragmente stehen real in der
   Datei (`:24` „Supersedes Entscheidung 1", „Teilfrage 4 akzeptierter
   Präzedenzfall" — beide ohne ihren ADR-Verweis), ebenso die in derselben
   Normierung entstandene Verkürzung. Zusätzlich stehen dort eine
   Register-Kennung und eine Wortform der Lebenszyklus-Familie (F-3): die
   Normierung ist nach ihrer eigenen Begründung nicht vollständig.
3. **R3 „Kurzformen entlastet" — bestätigt.** Die real vorgefundenen Formen
   sind im Erzeugnis aufgelöst (`LH-FA-CAP-001…003` → drei Zeilen-Kennungen;
   `LH-FA-CFG-003/004` → zwei; `LH-FA-RET-002…006` → fünf), und der Code
   bricht bei einer nicht auflösbaren Form sichtbar ab (`abdeckungsKurzform`:
   Ziffernanzahl ungleich drei, fehlender Abschluss-Span, nicht aufsteigende
   oder zu weite Spanne) — „sichtbarer Abbruch statt Raten" trägt.
4. **R4 „`Ort` bindet an Quellzeilen" — eingetreten: bestätigt.** Alle 37
   Zeilen binden auf `Datei:Zeile`; für den Go-Anteil ist die Zeile aus dem
   AST, für den Bash-Anteil aus dem Anker gesucht. Der Preis ist real, und
   die Stabilitäts-Zusage im Datei-Kopf nennt ihn nicht (F-4).

**Übergabe:** F-1 geht zurück an den **Implementer** (Rückkante
Review → Implementer) — zwei Sätze in Dateien, die dieser Diff angelegt hat,
ohne Verhaltenswirkung, aber gemessen unwahr; die Rückkante ist damit
**geöffnet**, eine kleine Fixrunde folgt. F-2 geht an den **Planner**
(Rückkante Review → Plan bei Plan-Defekt): die Korrektur liegt in §7/§8 des
Slice-Plans und im Register-Beleg — sie ist Planner-Arbeit bei der Closure und
öffnet die Implementer-Kante nicht. F-3 bis F-5 sind benannte Grenzen
innerhalb des Diffs (kein Handlungsbedarf an einer Zeile, die nicht ohnehin
angefasst wird), F-6 und F-7 sind Hinweise ohne erwartete Aktion.

**DoD-Nachzug:** **nein.** Die Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" bleibt in §2 des Slice-Plans offen, weil F-1 einen
Reviewer→Implementer-Rückgabe-Pfeil trägt (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde: der Nachzug greift nur, wenn keine
Fixrunde folgt). Nach der Fixrunde zieht Schritt 21 des
Implementer-Workflows die Zeile regulär nach.

**Die Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den
Zähler (Modul 10 §Pflege); alle sieben sind Erstauftreten (1×). Dieser Report
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt). Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11).
