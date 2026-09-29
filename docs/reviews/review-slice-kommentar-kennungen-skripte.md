# Review-Report: slice-kommentar-kennungen-skripte — 2026-09-29

**Review-Art:** Code — geprüft gegen den Slice-Plan
`slice-kommentar-kennungen-skripte` (Lifecycle-Ort:
`docs/plan/planning/in-progress/`), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
und die Hard Rules `AGENTS.md` §3.7/§3.12/§3.13/§3.9/§3.1 (Modul 10). Keine
DoD-Verifikation — das ist die Aufgabe des Verifiers (Modul 11).

**Gegenstand:** Diff `cd3a1c60..b16831d5` — die Lifecycle-Züge, `89d4fa3e`
(Werkzeug-Erweiterung + Tabellentests + Plan-Nachzug Messraum-Entscheidung),
die Tranchen `279e6416` (T1 Shell), `9cac13d0` (T2 YAML), `357b482c` (T3
Make), `58b9e415` (T4 SQL), `dd681e1f` (T5 Dockerfile) und `b16831d5`
(T2–T5-Rest, Träger-Nachzug Sensor-Vertrag/README, Suchlauf-Feld); nach
Fixrunde `754f0558`.

**Skill:** `.harness/skills/reviewer.md` (geschärft 2026-09-09) ·
**Modell:** GLM (Claude-Agent-SDK, Typ `reviewer`) · **Datum:** 2026-09-29

**Eingangs-Kontext:**

- Slice-Plan @ `b16831d5` (inkl. §3-Nachzug Messraum und §3.13-Suchlauf-Feld)
- `AGENTS.md` §3.7 (ein Anker; Ausnahme „Erzeugnis-Eingabe ist Abdeckung"),
  §3.12 (Zahlen-Ursprung), §3.13 (Suchlauf, Träger-Meldung), §3.9
  (Exit-Disziplin), §3.1 (Docker-only)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `harness/sensors/kommentar-kennungen.md`,
  `harness/sensors/suchlauf-nachmessen.md`
- `harness/conventions.md` (MR-000) · Baseline `v6.13.0` ·
  Vorgänger-Plan `docs/plan/planning/done/slice-code-kommentare-bereinigung.md`
  (Restmenge-Übergang)

---

## Findings

*(Stand des Erst-Reviews, Diff `cd3a1c60..b16831d5`; der Stand nach der
Fixrunde steht im Abschnitt „Fixrunde-Prüfung".)*

### F-1 — Sensor-Vertrag: Grenze 2 widerspricht der erweiterten Messform im selben Träger

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7 (Werkzeugvertrag); Skill-MEDIUM „Nachzug
  widerspricht dem Nachbarn im selben Träger"
- `pfad`: `harness/sensors/kommentar-kennungen.md:132-133` (Grenze 2) gegen
  `harness/sensors/kommentar-kennungen.md:47-53` (Kandidat, Zeilenkommentare)
- `befund`: Grenze 2 behauptet unverändert „**Nur Go-Kommentare.** Skripte,
  `Makefile`, `.sql`, `.yml` und Dockerfiles liest das Werkzeug nicht" —
  derselbe Abschnitt des Dokuments (Kandidat) definiert seit diesem Diff
  genau diese Formen samt Blockgrenze, und das Werkzeug zählt sie nachgemessen
  (135 Nicht-Go-Kandidaten am Stand `933ea5c0`). Der Test-Abschnitt (§Test)
  zählt ebenfalls nur die Go-Blockgrenzen auf; die drei neuen Tests
  (`TestRunLineForms`, `TestLineCommentMarker`, `TestLineBlocks`) und die
  Nicht-Go-Blockgrenze (nur-Marker-Zeile als Grenz-Marker) fehlen in der
  Aufzählung. Kein Gate fängt das.
- `verifizierbar`: ja — `make kommentar-kennungen` zählt Nicht-Go-Kandidaten;
  der Widerspruch ist eine Lese-Handlung im selben Dokument
- `klasse`: „Nachzug widerspricht dem Nachbarn im selben Träger"

### F-2 — Messinstrument-Wechsel in Tranche T1: Basismessung und Rückführungs-Schwelle gelten für das frühere Werkzeug

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 (Zahl im Träger — Ursprung: Instrument);
  Plan §4 Rückführungen (Schwelle „mehr als ~150 Kandidaten"); Plan §2
  Liefer-Punkt 2 („je Tranche Diff ohne Nicht-Kommentarzeile")
- `pfad`: `279e6416` — `tools/harness/kommentar-kennungen/main.go`
  (Grenz-Marker-Regel, +4 Code-Zeilen) und `main_test.go` (+Testfall); Plan
  §3-Nachzug (Basismessung-Satz)
- `befund`: Die Grenz-Marker-Regel (eine nur aus `#`/`--` bestehende Zeile
  beendet den Block) landete nicht im Werkzeug-Commit `89d4fa3e`, sondern in
  der Tranche T1 — die zugleich die Shell-Bereinigung trägt. Die Regel ändert
  das Messinstrument: nachgemessen an beiden Ständen mit beiden Instrumenten
  liefert `933ea5c0` 149 (Werkzeug `89d4fa3e`, Stand der Plan-Basismessung,
  bestätigt) und **160** (finales Werkzeug) — die Basismessung 149 liegt
  damit unter der Rückführungs-Schwelle (~150), derselbe Stand unter dem
  Instrument, das der Slice selbst adoptiert hat, liegt mit 160 darüber. Die
  Tranche-Zahl „62 → 0" (T1-Message) trägt das alte Instrument (sh 62 am
  Basismessungs-Stand; mit dem finalen Werkzeug 59), „Stand gesamt: 149 → 102"
  ist als altes Instrument konsistent (102 nachgemessen am T1-Baum mit
  `89d4fa3e`; finales Werkzeug 101) — der Plan trägt die Zahl 149 aber ohne
  Instrument-Label, und die Commit-Message von T1 behauptet „Nur
  Kommentarzeilen" bei enthaltenem Verhaltenswechsel des Messwerkzeugs (T5
  deklariert seinen gofmt-Ausgleich dagegen explizit). Die Grenz-Marker-Regel
  selbst ist notwendig (die Paragraphen-Struktur des gekürzten
  E2E-Runner-Kopfs nutzt `#`-allein als Absatztrenner) und im
  Sensor-Vertrag dokumentiert — unbestimmt bleibt nur die Zuordnung der
  Plan-Zahlen zum Instrument.
- `verifizierbar`: ja — `make kommentar-kennungen COUNT=1` an beiden Ständen
  mit beiden Werkzeugständen (149/160/102/101, Worktree-Klon im Temp)
- `klasse`: „Messwerkzeugwechsel in Tranche — Zahlen nicht komparabel"

### F-3 — DIFF-Modus sieht Nicht-Go-Kandidaten nicht: Pfadspec `*.go` im Aufrufer nicht nachgezogen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7 (Implementer-Probe); Skill-MEDIUM „Nachzug
  widerspricht dem Nachbarn im selben Träger" (Vertrag-Tabellenzeile gegen
  §Wer es aufruft)
- `pfad`: `tools/harness/kommentar-kennungen.sh:53` (`-- '*.go'`) gegen
  `harness/sensors/kommentar-kennungen.md` §Wer es aufruft (Implementer,
  Schritt 20)
- `befund`: Das Werkzeug liest seit diesem Diff Nicht-Go-Formen (der
  Tabellentest `TestRunLineForms` belegt den Diff-Modus sogar am
  Shell-Block), aber der Aufrufer erzeugt den Diff-Strom weiterhin mit
  `git diff -U0 … -- '*.go'` — nicht-Go-Änderungen erreichen die Messung im
  Diff-Modus nie. Mechanisch bestätigt: ein frisch hinzugefügter
  Shell-Block mit zwei Kennungen im Arbeitsbaum, `make kommentar-kennungen
  DIFF=b16831d5` → Exit 0, kein Kandidat. Der Implementer-Probe (Schritt 20,
  „vor der Übergabe … über die eigenen Änderungen") deckt damit genau die
  Formen nicht, die dieser Slice in die Messung nimmt; die
  Vertragstabelle dokumentiert die `*.go`-Eingabe, der §Wer-es-aufruft-Absatz
  verspricht aber unverändert die eigenen Änderungen. Die Asymmetrie entsteht
  mit diesem Diff (das Skript selbst ist unberührt) und ist nirgends als
  Grenze deklariert.
- `verifizierbar`: ja — reproduzierter Lauf (Exit 0 trotz Nicht-Go-Kandidat)
- `klasse`: „Messform erweitert, Aufrufer-Pfadspec nicht nachgezogen"

### F-4 — Plan-Breakdown summiert nicht: „Dockerfile 3" lässt die Beispiel-Dockerfiles weg

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 (Zahl im Träger); Plan §3-Nachzug
  (Basismessung-Satz)
- `pfad`: `docs/plan/planning/in-progress/slice-kommentar-kennungen-skripte.md`
  §3 (Messraum-Nachzug, Basismessung-Satz)
- `befund`: Die je-Form-Aufzählung „sh 62, yml 19, sql 15, mk 14, yaml 13,
  Makefile 4, Dockerfile 3" summiert zu 130, der Text nennt „davon 135
  Nicht-Go". Nachgemessen (Stand `933ea5c0`, Werkzeug `89d4fa3e`): die Totale
  149 (135 + 14 Go) stimmt, aber die Dockerfile-Familie zählt 8 Kandidaten
  (Wurzel-`Dockerfile` 3, `examples/Dockerfile` 1, `examples/csharp/Dockerfile`
  2, `examples/kotlin/Dockerfile` 2) — die fünf Beispiel-Dockerfile-Kandidaten
  fehlen in der Aufzählung. Die belastenden Totale (149/135/14) sind
  nachgemessen korrekt.
- `verifizierbar`: ja — Messung am Stand mit dem Stand-Werkzeug
- `klasse`: „Aufzählung trägt die Summe nicht"

### F-5 — Zeilenformen: false-positive-Klassen im Vertrag unbenannt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7; Sensor-Vertrag §Grenze (Ehrlichkeit der
  Grenz-Dokumentation)
- `pfad`: `harness/sensors/kommentar-kennungen.md` §Kandidat (Zeilenkommentare)
- `befund`: Der Go-Bullet nennt seine Abweichung explizit („Eine Kennung in
  einem Zeichenketten-Literal ist kein Kommentar und wird nicht gelesen");
  der Zeilenkommentar-Bullet nennt nur die nachgestellte Ausnahme. Die
  Gegenklasse — Zeilen, die wie Kommentare lesen, aber Inhalt sind — bleibt
  unbenannt: `#`-Zeilen in YAML-Blockskalaren (`key: |`), Zeilen in
  Shell-Heredocs und `--`-Zeilen in mehrzeiligen SQL-Zeichenketten liest das
  Werkzeug als Kommentarzeilen. Im Bestand zeigt sich das aktuell nicht
  (nach T1–T5: 0 Nicht-Go-Kandidaten, nachgemessen) — es ist eine
  vorsorgliche Dokumentationslücke, kein aktueller Befund.
- `verifizierbar`: nein — Lese-Handlung; ein Befund braucht erst ein
  entsprechendes Konstrukt im Baum
- `klasse`: „Grenze des Sensors unbenannt"

### F-6 — Blockgrenze je Form: nur-marker-Grenzfall nur für `#` getestet

- `kategorie`: LOW
- `quelle`: Skill-MEDIUM „fehlende Negativtests bei neuem öffentlichem
  Vertrag" (mild: die Grenzlogik ist marker-parametrisch)
- `pfad`: `tools/harness/kommentar-kennungen/main_test.go` (`TestLineBlocks`,
  `TestRunLineForms`)
- `befund`: Der nur-marker-Grenzfall (`#` allein beendet den Block) ist nur
  mit dem Shell-Marker getestet; für `--` (SQL) gibt es Treffer- und
  Nicht-Treffer-Fälle, aber keinen Grenzfall, für die Dockerfile-/Make-Form
  keinen Fortsetzungs-/Grenzzuschnitt. Da `lineBlocks` denselben Code für
  alle Marker ausführt, ist die Lücke schwach — der DoD-Wortlaut „Tabellentest
  je Form (Treffer, Blockgrenze, Nicht-Treffer)" trägt die Blockgrenze je
  Form aber nur für Shell voll.
- `verifizierbar`: ja — Testfall-Mutation (Grenz-Marker-Regel für `--`
  entfernen: kein Test wird rot)
- `klasse`: „Grenzfallestreuung unter Formen"

## Fixrunde-Prüfung (`754f0558`, 2026-09-29)

Je Finding nachgemessen bzw. nachgelesen; die Fixrunde bestätigt alle sechs:

- **F-1 — bestätigt (behoben).** Grenze 2 ist auf „Nur Code- und
  Konfig-Kommentare, keine Prosa" gestellt, nennt d-check `ids` als
  Prosa-Ebene und die Ausnahmen `gen/**`/`.harness/**` — konsistent mit dem
  Kandidat-Abschnitt; der §Test-Abschnitt zählt die Nicht-Go-Fälle auf
  (Form-Auswahl, Blockgrenze, Baum-Lauf, Diff-Modus). Nachgelesen am
  Stand `754f0558`.
- **F-2 — bestätigt (behoben).** Der Plan-Basismessung-Satz trägt jetzt beide
  Instrumente (149 mit `89d4fa3e`, **160** mit dem finalen Werkzeug am selben
  Stand — 160 gegen meine eigene Messung bestätigt), nennt die Grenz-Marker-
  Regel als Instrument-Anteil und führt den Schwellen-Ausgang (160 > ~150)
  für §7.
- **F-3 — bestätigt (behoben, Rest unten).** Die Pfadspec des Wrappers liest
  die sechs Nicht-Go-Formen mit (`.sh`, `.mk`, `.yml`, `.yaml`, `.sql`,
  `Makefile`, `Dockerfile`, `**/Dockerfile`). Mechanisch reproduziert: ein
  frisch hinzugefügter Shell-Block mit zwei Kennungen erscheint im
  DIFF-Lauf (`probe.sh:1-1  ADR-0001, ADR-0002`). Verbleibend: die
  Vertragszeile `DIFF=<Basis>` dokumentiert als Eingabe noch
  `git diff -U0 <Basis> -- '*.go'` — stumm gegen die erweiterte Wrapperspec
  (LOW-Rest, s. „Verbleibende Risiken").
- **F-4 — bestätigt (behoben).** Die unstimmige je-Form-Aufzählung ist
  entfernt; die Basismessung trägt nur noch die Totalen mit Instrument-Label
  (die per-Form-Zahlen dieses Reviews: sh 62/59, yml 19/21, sql 15/25, mk 14,
  yaml 13/15, Makefile 4, Dockerfile-Familie 8 — alt/finalem Instrument).
- **F-5 — bestätigt (behoben).** Die False-positive-Klassen (YAML-Blockskalare,
  Shell-Heredocs, SQL-Zeichenketten) sind in Grenze 2 benannt, mit Geltung
  („neue Fälle sind Befund") — der Ort, an den der Erst-Befund verwies.
- **F-6 — bestätigt (behoben).** `TestLineBlocks` trägt den nur-marker-
  Grenzfall jetzt zusätzlich mit SQL-Marker (Fall „nur-Marker-Zeile ist
  Grenz-Marker (SQL)").

**Sensors nach der Fixrunde (Exit direkt gelesen, `AGENTS.md` §3.9):**
`make test` Exit 0 · `make fmt-check` Exit 0 · `make suchlauf-nachmessen`
2/2 OK, Exit 0 · `make kommentar-kennungen COUNT=1` = 14, Exit 0 ·
DIFF-Probe (siehe F-3) — Nicht-Go-Kandidat wird gemeldet.

## Negativbefunde

- geprüft, ohne Befund: **Messstand** — `make kommentar-kennungen COUNT=1` =
  14 (Exit 0), Vollist 14 Kandidaten, **0 außerhalb `test/integration`**
  (nachgemessen, nicht übernommen).
- geprüft, ohne Befund: **Übergang der geerbten Restmenge** — 14 = 14; die
  Blöcke sind die `func TestE2E*`-Godocs mit Grenz-Vermerk am Block
  („Kennungs-Menge: … Abdeckung, keine Herkunft, kein Kürzen auf einen
  Anker", `test/integration/integration_test.go:185-204` nachgelesen); die
  Ausnahme steht bestandsmäßig in `AGENTS.md` §3.7 (nicht Teil dieses Diffs).
- geprüft, ohne Befund: **Abdeckungs-Daten unberührt** — `abdeckung_declare`:
  0 Treffer im Diff; `tools/schema/schema.yaml`: alle geänderten Zeilen sind
  `#`-Kommentarzeilen; die Tranchen T2–T4 sind vollständig kommentarzeilig
  (mechanisch geprüft je Commit: Residual leer).
- geprüft, ohne Befund: **Anker-Wahl (Stichproben gegen ADR-Dateien und
  Specs)** — alle im Diff behaltenen ADR-Kennungen existieren
  (`docs/plan/adr/`, mechanisch gegen Dateinamen geprüft); alle behaltenen
  `LH-*` gegen `spec/lastenheft.md`, alle `SPEC-*` gegen
  `spec/pflichtenheft.md` — 0 Fehltreffer. Stichproben mit geprüftem Anker:
  `examples/bootstrap.sh`
  ([ADR-0098](../plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)),
  `tools/harness/run-integration-tests.sh`
  ([ADR-0043](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) je Absatz,
  dazu die dort stehenden Anker 0026/0028/0030/0047/0120/0049 gegen
  Existenz),
  `Dockerfile`
  ([ADR-0060](../plan/adr/0060-grpc-streaming-mechanismus.md) proto-Stufe,
  [ADR-0071](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  coverage,
  [ADR-0042](../plan/adr/0042-transport-typen-am-port.md) build),
  `harness/mk/examples.mk`
  ([ADR-0090](../plan/adr/0090-beispiel-clients-volle-matrix.md) Festlegung 5,
  [ADR-0098](../plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
  Festlegung 2),
  `tools/schema/nacharbeit-roles.sql`
  ([LH-QA-SEC-001](../../spec/lastenheft.md) Rollenschnitt, 002 im selben
  Absatz), `.a-check.yml`
  ([ADR-0002](../plan/adr/0002-abhaengigkeitsrichtung.md)
  Abhängigkeitsrichtung, 0076/0079 gegen Existenz).
- geprüft, ohne Befund: **Chronik-/Vorher-Nachher-Sprache in +Zeilen** —
  Mustersuche über alle +Zeilen (inkl. Fixrunde `754f0558`): Treffer nur als
  Zustandsbeschreibung oder erlaubter Herkunfts-Anker („· seit
  slice-kommentar-kennungen-skripte" in der README-Sensors-Zeile).
- geprüft, ohne Befund: **Suchlauf-Feld** — `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-kommentar-kennungen-skripte.md`:
  2/2 Zeilen OK (506 @ `89d4fa3e`, 350 Arbeitsbaum), Exit 0 — auch nach der
  Fixrunde.
- geprüft, ohne Befund: **Sensors** — `make test` (Exit 0, ungefiltert
  gesichert; vor und nach der Fixrunde), `make fmt-check` (Exit 0),
  `make a-check` (0 Befunde, Exit 0 — berührt, weil `.a-check.yml` im Diff
  ist), `make kommentar-kennungen` (7 Läufe, s. o.).
- geprüft, ohne Befund: **Docker-only** — der Diff führt keine
  Host-Toolchain, keine in-place-Umschreibung und keine Umleitung auf
  Repo-Dateien ein; die Werkzeug-Tests laufen im gepinnten Toolchain-Container
  (`make test`-Lauf oben).

## Summary

| Kategorie | Anzahl (vor Fixrunde) | nach Fixrunde |
|---|---|---|
| HIGH | 0 | 0 |
| MEDIUM | 3 | 0 |
| LOW | 3 | 0 (Reste: „Verbleibende Risiken") |

**Finding-Klassen dieses Laufs:** „Nachzug widerspricht dem Nachbarn im
selben Träger" (1×) · „Messwerkzeugwechsel in Tranche — Zahlen nicht
komparabel" (1×, neu) · „Messform erweitert, Aufrufer-Pfadspec nicht
nachgezogen" (1×, neu) · „Aufzählung trägt die Summe nicht" (1×) · „Grenze
des Sensors unbenannt" (1×) · „Grenzfallestreuung unter Formen" (1×)

## Verdikt

**Merge-blockierend: nein (nach Fixrunde).** Alle drei MEDIUM- und drei
LOW-Findings sind durch `754f0558` behoben und von mir gegen die eigenen
Messungen bestätigt (Fixrunde-Prüfung oben); die Sensor-Läufe nach der
Fixrunde sind grün (Exit-Codes direkt gelesen).

**Verbleibende Risiken (nicht merge-blockierend):**

- Die Vertragszeile `DIFF=<Basis>` dokumentiert als Eingabe noch
  `git diff -U0 <Basis> -- '*.go'` — die Wrapperspec liest seit der Fixrunde
  mehr; ein Leser, der die Eingabe nachbaut, filtert die Nicht-Go-Formen
  wieder heraus (LOW-Rest; eine Zeile im Sensor-Vertrag, für §7).
- Pfadspec-Restlücke: ein künftiges `Dockerfile.<Variante>` (das Werkzeug
  matcht das Namens-Präfix) würde baumseitig gelesen, im DIFF-Modus aber von
  keiner Pathspec-Form erfasst; im Bestand existiert keine solche Datei.
- Die Zahl „11 gemessene Ketten, 5 Shell + 6 Makefile" (Plan §3-Nachzug,
  nachgestellte Kommentare) ist nicht unabhängig nachgemessen — das
  Zählprogramm für „Kette" ist nicht definiert; sie stützt keine Entscheidung.
- Die Restmenge von 14 ist eine Zustandsgröße: der nächste Godoc mit zwei
  Kennungen in `test/integration` bewegt sie; die Ausnahme in `AGENTS.md`
  §3.7 trägt den Umfang („Godocs der `func TestE2E*`"), ein Gate hält ihn
  nicht.

**Übergabe:** die Finding-Klassen gehen in die Slice-Closure §7 und von dort
in den Steering-Loop-Zähler (neu: „Messwerkzeugwechsel in Tranche — Zahlen
nicht komparabel", „Messform erweitert, Aufrufer-Pfadspec nicht nachgezogen",
„Aufzählung trägt die Summe nicht", „Grenze des Sensors unbenannt",
„Grenzfallestreuung unter Formen" — je 1×, erstes Auftreten). Dieser Report
ist ein Lauf-Beleg.
