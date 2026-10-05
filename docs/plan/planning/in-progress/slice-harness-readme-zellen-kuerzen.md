# Slice harness-readme-zellen-kuerzen: die Index-Zeilen der `harness/README.md` bleiben kurz, der ausführliche Vertrag lebt in Sensor- und Target-Dateien, und eine d-check-Regel hält die Zellenlänge

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

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (CI/CD-Pipeline über GitHub Actions, mit Pin-Inventar),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen:
beim Umzug behalten Zahlen ihren Ursprung). Keine `LH-*`-Anforderung ist berührt:
der Slice ändert Harness-Dokumente und eine Doku-Gate-Regel, nicht das Produkt.

**Berührte Spec-Stellen:** — (keine; die Quellen-Rangfolge in `harness/README.md`
bleibt unberührt).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag). Reihenfolge zu
`baseline-6-14-0-dokumente-nachziehen`: dieser Slice zuerst, der Bump-Slice
danach (§1, §4).

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**.

**Anlass (gemessen am Parent `4db7cf95`, 2026-10-05, Befehle in §3).**
`harness/README.md` hat 136.970 Byte in 265 Zeilen (Wert des Auftraggebers,
**übernommen**; `wc -c`/`wc -l` am Planungsstand bestätigen 136970 und 265); die
Vorlage `.harness/baseline/v6.14.0/templates/harness/README.template.md` hat
11 KB in 210 Zeilen (**übernommen**). 15 Zeilen haben mehr als 2.000 Zeichen
(gemessen, Suchlauf Zeile 1); die längste ist die Zeile `make test-integration`
der Tabelle „Werkzeuge“ mit 26.286 Zeichen (**übernommen**, im Slice
nachzumessen). Die Zellen der zweiten Spalte (Sensors „Vertrag“, Werkzeuge
„Tut was“) sind in 57 Zeilen länger als 120 Zeichen und in 55 länger als 220
(gemessen, Suchlauf Zeilen 2 und 3; beide Tabellen zusammen, Backtick-Zeilen).
Die Baseline-Regel (Vorlage `harness/sensors/gate.template.md`, Regelwerk
`grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt) verlangt:
die Index-Zeile bleibt kurz und verlinkt eine Sensor- oder Target-Datei, sobald
ein Target mehr braucht als einen Satz; ob der Überhang unter der Tabelle steht
oder in die Zelle gedrängt wurde, ist dieselbe Sache. Heute stehen unter
`harness/targets/` 5 Dateien (gemessen, Suchlauf Zeile 4: je ein Titel) und
unter `harness/sensors/` 14 Dateien; etwa 40 der 59 Werkzeug-Zeilen haben
keine eigene Datei (**übernommen** vom Auftraggeber; die Zeilenzahl 59 ist im
Slice nachzuzählen).

**Ziel:** Jede Index-Zeile der `harness/README.md` trägt höchstens einen Satz
(„Tut was“ ≤ 120, „Vertrag“ ≤ 220 Zeichen), der ausführliche Inhalt steht
wortgleich in der verlinkten Datei, und `make docs-check` hält die Grenze.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Inhaltliche Korrektur der umgezogenen Aussagen.** Der Umzug ändert keinen
  Wortlaut und keine Zahl; fällt eine Aussage als falsch oder veraltet auf, geht
  sie in eine Befund-Liste (Closure-Notiz, ggf. Beobachtung), nicht still in den
  neuen Text — sonst prüfte derselbe Diff Umzug und Berichtigung zugleich
  ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)).
- **Die Tabellen „Source precedence“, „Guides“ und die Abschnitte „Traceability
  rules“, „Minimal agent workflow“.** Sie sind kurz und tragen die Index-Funktion;
  die neue Regel adressiert nur den Abschnitt „Sensors (Feedback-Gates)“.
- **Streichen von Index-Zeilen.** Der Gate-Index steht einmal in dieser Datei
  (`AGENTS.md` §4: Zeile je real existierendem Target); gekürzt wird die Zelle,
  nicht die Zeile.
- **Neuschnitt der bestehenden fünf Target-Dateien und 14 Sensor-Dateien.**
  Sie bleiben Bestand; ein Target mit vorhandener Datei verlinkt sie und die
  Mega-Zelle wird dort nur ergänzt, soweit ihr Inhalt dort noch fehlt.
  **Geändert durch Plan-Nachzug 1 (§3 Umsetzung):** die bestehende Datei trägt
  die ganze Zelle wortgleich in einem Abschnitt `## Fassung im Gate-Index`;
  ersetzt wird weiterhin nichts.
- **Eine eigene Baseline-Vorlage für `harness/targets/`.** Die Baseline führt
  keine; der Aufbau folgt den bestehenden Dateien (Vertrag · Aufruf · Grenzen ·
  „Kein Gate“).
- **Die Prüfung der Baseline-Abgleiche des Bump-Slices.** Der Slice
  `baseline-6-14-0-dokumente-nachziehen` liest `harness/README.md` nur (Abgleich
  gegen die Vorlage, Liefer-Punkt 2c) und zitiert ihre Struktur; er ändert dort
  höchstens die Zeile `make pin-stale-baseline`. Reihenfolge: dieser Slice
  zuerst oder danach — nie gleichzeitig (WIP-Limit 1); steht der Bump-Slice
  davor, liest er die ungekürzte Datei, steht er danach, die gekürzte; die
  Zeile `make pin-stale-baseline` ist dann kurz und er ändert sie dort.

**Keine Mindestzahl.** Die vier Klassen des Ausschlusses sind ein Suchraster,
keine Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben dieser Liste sind
**Zusagen** („zu belegen durch …“): der Planungsstand hat keinen der Läufe
gefahren ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
Instanz B).

- [x] **Zuschnitt und Umzug (Liefer-Punkt 1).** Je Ziel oder Ziel-Gruppe der
      Tabelle „Werkzeuge“ (z. B. `make pin-stale-*`, `make test-sdk-*-integration`,
      `make examples-*` und `make example-*`, die `make test-*`-Tier-Läufe,
      `make bench`, `make proto-generate`, die `.github/workflows/*`) liegt eine
      Datei unter `harness/targets/` (eine Gruppe darf eine Datei teilen; Aufbau
      wie die bestehenden Target-Dateien); Gate-Zeilen der Tabelle „Sensors“
      verlinken `harness/sensors/<name>.md` (neu, wo die Datei fehlt). Der
      ausführliche Inhalt der Zelle ist **umgezogen**: Wortlaut und Zahlen samt
      Herkunftsangabe erhalten, nichts erfunden, nichts weggelassen. Gegenprobe:
      die Zeichenzahl der Mega-Zelle gegen die Summe der neuen Absätze (Abweichung
      nur durch Entfernen von Tabellen-Syntax und Verlinkung, zu belegen durch die
      gedruckten Zahlen) und ein Suchlauf der Eigennamen (Make-Ziele, Dateipfade,
      Schlüsselzahlen) in alter Zelle und neuer Datei. Aussagen, die dabei als
      falsch oder veraltet auffallen, stehen in einer Befund-Liste in §7.
- [x] **Zellen kürzen (Liefer-Punkt 2).** In `harness/README.md` ist jede Zelle
      „Tut was“ ≤ 120 und jede Zelle „Vertrag“ der Sensors-Tabelle ≤ 220 Zeichen
      (gemessen wie dargestellt, inkl. Markdown-Syntax); die Bindung-Zelle
      verlinkt die Datei in Link-Form (kein Pfad im Fließtext); die Index-Zeilen
      bleiben vollständig (Zeilenzahl der Tabellen vorher = nachher). Die Tabellen
      „Source precedence“, „Guides“ und die Prosa-Abschnitte sind unverändert
      (`git diff` zeigt keine Änderung dort). Träger, die Zeilenzahl oder Position
      der README zitieren, sind nachgezogen (Suchlauf unten, §3.13 der Regeln).
- [x] **Regel, die die Grenze hält (Liefer-Punkt 3).** `.d-check.yml` trägt unter
      `structure:` eine neue Regel für `harness/README.md`, Abschnitt
      `## Sensors (Feedback-Gates)`, Spalten `Tut was` (max 120) und `Vertrag`
      (max 220), min jeweils 1 — **nach** dem Umbau eingeführt, nicht davor
      (sonst rot). Zu messen und im Bericht festzuhalten: ob die Regel beide
      Tabellen des Abschnitts erfasst (Spalte per Kopfzeilen-Name; erwartet: ja,
      **hergeleitet**, nicht gemessen) und ob eine absichtlich überlange Zelle auf
      einer Kopie im Scratchpad `make docs-check` rot färbt (Mutationsprobe mit
      gesehener Farbe, Stelle und Instanz benannt). `harness/sensors/docs-check.md`
      §Bindung/Gegenstand nennt die neue Regel (einschließlich ihrer Grenze: sie
      misst Länge, nicht ob der Satz die Zeile trägt). Eine Verschärfung braucht
      keine ADR (`AGENTS.md` §3.6 betrifft die Lockerung).

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [x] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make suchlauf-nachmessen PLAN=` mit diesem Plan Exit 0 am Endstand.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durch den Verifier (Belege, nicht Behauptung): Wortlaut-Gegenprobe
      und Mutationsprobe der Regel unabhängig nachgefahren.
- [ ] Doku-Update: `harness/README.md` ist der Gegenstand;
      [`AGENTS.md`](../../../../AGENTS.md) §4 und die Träger aus dem Suchlauf
      bleiben konsistent (gemeldete Träger fremder Dateien mit der Closure
      nachgezogen).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/targets/*.md` (neu, nach Gruppen) | neu | Umzug der ausführlichen Zellen der Tabelle „Werkzeuge“ (Liefer-Punkt 1) |
| `harness/sensors/*.md` (neu, wo für eine Gate-Zeile keine Datei besteht) | neu | Umzug des Überhangs der Sensors-Zellen (Liefer-Punkt 1) |
| bestehende `harness/sensors/*.md` und `harness/targets/*.md` (15 Dateien, Zuordnung unten) | update | Abschnitt `## Fassung im Gate-Index` mit der Zelle wortgleich (Liefer-Punkt 1; Plan-Nachzug 1 unter §3 Umsetzung) |
| `harness/README.md` | update | Zellen auf Ein-Satz-Länge, Bindung verlinkt die Datei (Liefer-Punkt 2) |
| `.d-check.yml` | update | neue `structure:`-Regel, nach dem Umbau (Liefer-Punkt 3) |
| `harness/sensors/docs-check.md` | update | §Bindung/Gegenstand nennt die Regel (Liefer-Punkt 3) |
| Träger aus dem Suchlauf (z. B. Zitate der README-Struktur) | update, nur falls ein Treffer wahr zu halten ist | Nachzug, §3.13 |

**Ansatz.**

- Gruppen vor dem Schnitt festlegen: die 59 Werkzeug-Zeilen nach Familie sortieren,
  je Familie einen Dateinamen (`pin-stale.md`, `sdk-integration.md`, `examples.md`,
  `tier-tests.md`, `workflows.md`, `proto-generate.md` o. ä.), die Zuordnung als
  Tabelle in den Bericht.
- Je Zelle umziehen, nicht neu schreiben: der Zelltext wird Absatz bzw. Abschnitt
  der Zieldatei (Aufbau Vertrag · Aufruf · Grenzen · „Kein Gate“, soweit die Zelle
  das trägt); die Zelle der README behält einen Satz plus Verweise.
- Ein Target mit bestehender Datei (die fünf in `harness/targets/`, die 14 in
  `harness/sensors/`) verlinkt sie; fehlt in ihr Inhalt der Zelle, wird er
  ergänzt, nicht ersetzt.
- Wortlaut-Gegenprobe je Zeile: Zeichenzahl vorher/nachher und Suchlauf der
  Eigennamen. Auffälliges (veraltete Zahl, widersprechende Aussage) in die
  Befund-Liste.
- Die d-check-Regel erst nach dem letzten Kürzen setzen.

**Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
Eigenschaften: (A) Länge der Zeilen und Zellen der `harness/README.md`; (B) die
Menge der Target-Dateien; (C) Träger, die die README und ihren Abschnitt
„Sensors“ zitieren (Symbolnamen `harness/README.md`, `§Sensors`; Zeilen-Lokatoren
und Abschnittsnummern fängt `grep` nicht — Lese-Handlung des Reviewers). Parent
ist der Stand `4db7cf95` (nie `HEAD`). Suchraum der Zeilen 5 und 6: ganzer Baum
außer `docs/reviews`, `.harness/baseline`, `docs/plan/planning/done`
(Records), `docs/plan/planning/observations` (Register hält Belege mit dem Stand
ihrer Zeit) und dieser Plan-Datei. **Parent-Zeilen gemessen** (2026-10-05);
Zeilen mit `diff` setzt der Implementer am Endstand und liest jede Trefferzeile.
**Erwartung am Endstand** (hergeleitet, nicht gemessen): Zeile 1 sinkt auf 0,
Zeilen 2 und 3 sinken auf 0, Zeile 4 steigt um die neuen Dateien; die Zeilen 5 und
6 bleiben im Wesentlichen stehen (Zitate nennen die Datei, nicht ihre Länge),
jede Trefferzeile wird gelesen, ob sie eine Länge, Position oder einen Zelleninhalt
behauptet, der sich bewegt hat.

**Startmessung (Implementer, 2026-10-05).** `make suchlauf-nachmessen PLAN=` mit
diesem Plan am Stand `b61412ac` mit den Parent-Zeilen `4db7cf95`: sechs Zeilen
`OK`, Exit 0. Seit `4db7cf95` änderte `d0e781b2` drei Zeilen der
`harness/README.md` (die drei `make test-sdk-*-integration`-Zeilen, TLS-Phasen),
ohne eine der Schwellen zu kreuzen: dieselben sechs Befehle am Stand `b61412ac`
liefern 15 · 57 · 55 · 5 · 113 · 48 (gemessen, `git grep … b61412ac -- …`). Der
Parent ist deshalb auf `b61412ac` nachgezogen, die Soll-Zahlen bleiben; der
Ausschluss der Plan-Datei nennt ihren Ort in `in-progress/` (das Werkzeug
schließt sie ohnehin aus).

**Endstand (`diff`, gemessen am Arbeitsbaum vor dem Plan-Commit; die Plan-Datei
ist aus dem Suchraum ausgeschlossen).** Zeile 1 sinkt auf 0 (längste Zeile der
README jetzt 790 Zeichen). Zeile 2 sinkt auf **6**, nicht auf 0 wie erwartet:
das Muster misst die zweite Spalte beider Tabellen mit 120, die Spalte `Vertrag`
darf aber 220 tragen (§2 Liefer-Punkt 2) — die sechs Treffer sind Gate-Zeilen
(`make a-check`, `make commit-traceability`, `make coverage-gate`,
`make generated-sync`, `make handbuch-public-doc-check`,
`make meldungscodes-check`), keine Werkzeug-Zeile. Zeile 3 sinkt auf 0. Zeile 4
steigt von 5 auf 18 (13 neue Target-Dateien). Zeile 5 steigt von 113 auf 148 und
Zeile 6 von 48 auf 81; jede neue Trefferzeile ist gelesen: 30 sind der
Einleitungssatz „Ausführliche Fassung der Index-Zeile(n) aus
[`harness/README.md` §Sensors]…“ der 15 neuen und der 15 ergänzten Dateien, die
übrigen stehen in `harness/sensors/docs-check.md` (§Vertrag, Grenze 11, §Bindung)
und in der neuen Regel in `.d-check.yml`. Die stehenden Treffer sind gelesen;
die, die eine bewegte Eigenschaft behaupten, stehen unter „Träger-Meldungen“ in
§3 Umsetzung.

```suchlauf
b61412ac838199028c81a33b217b7991b52189d4 15 -n -E '^.{2001,}$' -- harness/README.md
b61412ac838199028c81a33b217b7991b52189d4 57 -n -E '^\| `[^|]*\| [^|]{121,} \| ' -- harness/README.md
b61412ac838199028c81a33b217b7991b52189d4 55 -n -E '^\| `[^|]*\| [^|]{221,} \| ' -- harness/README.md
b61412ac838199028c81a33b217b7991b52189d4 5 -n -E '^# `' -- harness/targets/*.md
b61412ac838199028c81a33b217b7991b52189d4 113 -n -F 'harness/README.md' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations :!docs/plan/planning/in-progress/slice-harness-readme-zellen-kuerzen.md
b61412ac838199028c81a33b217b7991b52189d4 48 -n -E 'Sensors \(Feedback-Gates\)|§Sensors' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations :!docs/plan/planning/in-progress/slice-harness-readme-zellen-kuerzen.md
diff 0 -n -E '^.{2001,}$' -- harness/README.md
diff 6 -n -E '^\| `[^|]*\| [^|]{121,} \| ' -- harness/README.md
diff 0 -n -E '^\| `[^|]*\| [^|]{221,} \| ' -- harness/README.md
diff 18 -n -E '^# `' -- harness/targets/*.md
diff 148 -n -F 'harness/README.md' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations :!docs/plan/planning/in-progress/slice-harness-readme-zellen-kuerzen.md
diff 81 -n -E 'Sensors \(Feedback-Gates\)|§Sensors' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations :!docs/plan/planning/in-progress/slice-harness-readme-zellen-kuerzen.md
```

### Umsetzung (Implementer, 2026-10-05)

**Nachgemessene Anlass-Werte (§1).** Die Zeile `make test-integration` hatte am
Parent 26.286 Zeichen (gemessen, `git show b61412ac:harness/README.md | sed -n
147p | awk '{print length($0)}'`). Die Tabelle „Werkzeuge“ hatte 58 Zeilen, davon
3 Vorlagen-Platzhalter (`make <mover>`, `make <messung>`, `make <vorschau>`);
von den 55 echten verlinkten 22 schon eine Datei unter `harness/sensors/` oder
`harness/targets/`, 33 keine (gemessen, `git show b61412ac:harness/README.md |
awk 'NR>=131 && NR<=185' | grep -cE '\]\((sensors|targets)/'`). Die Tabelle
„Sensors“ hat 12 Zeilen, davon 1 Platzhalter (`<make-target>`). Die README hat
nach dem Umbau 34.679 Byte in 265 Zeilen (gemessen, `wc -c -l`).

**Zuordnung Gruppe → Datei** (63 geänderte Zeilen; unverändert blieben
`make baseline-verify`, `make docs-check`, `make a-check` — Vertrag schon ≤ 220,
Datei schon verlinkt — und die vier Platzhalter):

| Gruppe (Index-Zeilen) | Datei | Art |
|---|---|---|
| `make image`, `make image-stale`, `make image-cve` | `harness/targets/image.md` | neu |
| `make pin-stale-race`/`-pgtest`/`-dmigrate`/`-acheck`, `-dcheck`, `-baseline`, `-actions` | `harness/targets/pin-stale.md` | neu |
| `make doc-trace`, `make doc-ci-matrix` | `harness/targets/doc-trace.md` | neu |
| `make proto-generate` | `harness/targets/proto-generate.md` | neu |
| `make test`, `make test-store`, `make test-replication`, `make test-notify` | `harness/targets/tier-tests.md` | neu |
| `make test-integration` | `harness/targets/test-integration.md` | neu |
| die sieben `.github/workflows/*.yml`-Zeilen | `harness/targets/workflows.md` | neu |
| `make test-command-guard` | `harness/targets/command-guard.md` | neu |
| `make schema-validate` | `harness/targets/schema-validate.md` | neu |
| `make bench` | `harness/targets/bench.md` | neu |
| `make examples-csharp`, `make examples-kotlin`, `make example-run-*`, `make example-demo-*`, `make example-transformation-demo` | `harness/targets/examples.md` | neu |
| `make sdk-pack-csharp`/`-python`/`-kotlin` | `harness/targets/sdk-pack.md` | neu |
| `make test-sdk-kotlin`/`-csharp`/`-python-integration` | `harness/targets/sdk-integration.md` | neu |
| `make commit-traceability` | `harness/sensors/commit-traceability.md` | neu |
| `make gates` | `harness/sensors/gates.md` | neu |
| `make coverage-gate` · `make generated-sync` · `make doc-tracked` | `harness/sensors/coverage-gate.md` · `generated-sync.md` · `docs-check.md` | ergänzt |
| `make sdk-public-doc-check` + Tabellentest · `make handbuch-public-doc-check` + Tabellentest · `make ausgabe-kennungen-check` + Tabellentest · `make meldungscodes-check` + Tabellentest | die gleichnamigen Dateien unter `harness/sensors/` | ergänzt |
| `make pin-stale-all` + Tabellentest · `make suchlauf-nachmessen` + Tabellentest · `make kommentar-kennungen` + Tabellentest · `make fmt-check` + Tabellentest | die gleichnamigen Dateien unter `harness/sensors/` | ergänzt |
| `make image-mutation` + Tabellentest · `make schema-rollout` · `make test-sdk-kompat` · `make test-sdk-altserver` | `harness/targets/image-mutation.md` · `schema-rollout.md` · `sdk-kompat.md` · `sdk-altserver.md` | ergänzt |

**Plan-Nachzug — drei Abweichungen vom Ansatz, mit Grund.**

1. *Bestehende Dateien: ganze Zelle statt „soweit fehlt“.* Der Ansatz sah vor,
   in eine bestehende Datei nur den Teil der Zelle zu übernehmen, der dort
   fehlt. Ob ein Satz „schon da“ ist, ist ein Urteil über Wortlaut-Nähe und lässt
   sich am Diff nicht nachprüfen; ein Eigennamen-Abgleich (Backtick-Tokens und
   Zahlen der Zelle gegen die Datei) fand in 13 der 22 Zellen fehlende Tokens.
   Jede bestehende Datei trägt deshalb am Ende einen Abschnitt
   `## Fassung im Gate-Index` mit der Zelle wortgleich; nichts wurde ersetzt.
   Folge: Inhalt, der dort schon stand, steht nun zweimal in derselben Datei
   (dieselbe Doppelung, die vorher zwischen README und Datei bestand) —
   Befund 9 unten.
2. *Eine Datei für `make gates`.* Die Zelle (230 Zeichen) zählte die zehn
   Gates auf; die Aufzählung steht wortgleich in `harness/sensors/gates.md`, die
   Zelle sagt „alle inneren Gates dieser Tabelle“.
3. *Absätze.* Beim Umzug ist eine Zelle ein Absatz; vor „Zusätzlich“ und
   „Seit slice-“ (nach einem Punkt) beginnt ein neuer Absatz, damit die
   26-KB-Zelle lesbar bleibt. Weitere Abweichungen vom Wortlaut sind nur
   Tabellen-Syntax und Verlinkung: `\|` → `|` und die Link-Präfixe relativ zum
   neuen Ort (`../` → `../../`, `sensors/`, `targets/`, `mk/`, `conventions/`
   → `../…/`). Die Bindung-Zelle der README bleibt, sie bekommt den Link auf die
   Datei vorn hinter „kein Gate,“ (bzw. vor „· seit“); die Bindung steht in den
   neuen Target-Dateien zusätzlich als Zeile **Bindung:**.

**Gegenprobe Wortlaut (Liefer-Punkt 1).** Je geänderter Index-Zeile wird die
Zelle am Parent mit derselben Umformung (Absatz, `\|`, Link-Präfix) gebildet und
jeder Absatz als ganze Zeile in der Datei gesucht, die die neue Bindung-Zelle
zuerst verlinkt; gedruckt werden Zeichenzahl alt/umgeformt und Absatzzahl. Der
Befehl (Host-Werkzeuge `bash`, `git`, `awk`, `sed`, `grep`; im Repo-Wurzel
ausführen, Argument der Parent):

```bash
P=b61412ac; T=$(mktemp -d); git show "$P:harness/README.md" > "$T/alt"; fail=0; n=0
for nr in $(seq 1 "$(wc -l < "$T/alt")"); do
  alt=$(sed -n "${nr}p" "$T/alt"); neu=$(sed -n "${nr}p" harness/README.md)
  [ "$alt" = "$neu" ] && continue; n=$((n+1))
  ziel=$(printf '%s\n' "$neu" | grep -oE '\]\((sensors|targets)/[^)]*\.md\)' | head -1 | sed -E 's/^\]\(//; s/\)$//')
  zelle=$(printf '%s\n' "$alt" | awk -F' \\| ' '{ if (NF==4) print $2 " | " $3; else print $2 }')
  form=$(printf '%s\n' "$zelle" | sed -E 's/\\\|/|/g; s#\]\(\.\./#](../../#g; s#\]\((sensors|targets|mk|conventions)/#](../\1/#g')
  ok=OK; while IFS= read -r abs; do [ -z "$abs" ] || grep -qFx -- "$abs" "harness/$ziel" || ok=FEHLT
  done < <(printf '%s\n' "$form" | sed -E 's/\. (Zusätzlich|Seit slice-)/.\n\1/g')
  echo "$ok $nr $ziel alt=${#zelle} form=${#form}"; [ $ok = OK ] || fail=1
done; echo "geaenderte Zeilen: $n, Fehler: $fail"; rm -rf "$T"
```

Gedruckt am Arbeitsbaum (Auszug; alle 63 Zeilen `OK`): `OK 147
targets/test-integration.md alt=26149 form=26314`, `OK 183
targets/sdk-integration.md alt=6461 form=6488`, `OK 166 targets/bench.md
alt=4473 form=4500`, `OK 124 sensors/gates.md alt=230 form=230`, letzte Zeile
`geaenderte Zeilen: 63, Fehler: 0`. Die Differenz form−alt ist je Zelle +3 je
umgeschriebenem Link-Präfix und −1 je `\|` (Beispiel `make example-run-go`:
alt=1148, form=1145, drei `\|`). Weil jeder Absatz als ganze Zeile gefunden
wird, sind alle Eigennamen (Make-Ziele, Pfade, Zahlen, Herkunftsangaben) der
alten Zelle in der Datei; zusätzlich ist jede neue Datei byte-gleich mit einem
im Scratchpad aus den Zellen erzeugten Entwurf (`cmp`, 13 Target-, 2
Sensor-Dateien) und jede ergänzte Datei byte-gleich mit „alter Inhalt +
Leerzeile + Abschnitt“ (`cmp`, 15 Dateien). Mutationsprobe der Gegenprobe (Kopie
im Scratchpad, `harness/targets/bench.md`: „vier eigenständige Bench-Skripte“ →
„vier Bench-Skripte“): `FEHLT 166 targets/bench.md`, `Fehler: 1`, Exit 1.

**Zellen und Zeilen (Liefer-Punkt 2).** Längste Zelle `Vertrag`: 162 Zeichen
(`make commit-traceability`); längste Zelle `Tut was`: 113 Zeichen
(`make pin-stale-all`) (gemessen, `awk -F' \\| '` über den Abschnitt, Länge der
zweiten Spalte). Zeilen der beiden Tabellen vorher = nachher = 70 (gemessen:
Zeilen des Abschnitts, die mit `` | ` `` beginnen, an Parent und Arbeitsbaum).
`git diff b61412ac -- harness/README.md` zeigt zwei Hunks, Zeilen 117–124 und
131–185 (63 Zeilen entfernt, 63 eingefügt); Source precedence, Guides und die
Prosa sind unberührt.

**Regel und Mutationsprobe (Liefer-Punkt 3).** `.d-check.yml` §structure trägt
die Regel für `harness/README.md`, Abschnitt `## Sensors (Feedback-Gates)`,
Spalten `Vertrag` (1–220) und `Tut was` (1–120). Gemessen, ob sie beide Tabellen
erfasst: ja — dieselbe Regel auf die **ungekürzte** README (Kopie im
Scratchpad, Klon am Stand `b61412ac`) meldete 63 Befunde
`section-cell-oversized`, darunter die Werkzeug-Zeilen 131 bis 185 (Spalte
`Tut was`) und die Gate-Zeilen 117 bis 124 (Spalte `Vertrag`); eine Tabelle ohne
die Spalte meldet nichts. Mutationen auf der Kopie mit gekürzter README und
Regel, je `make docs-check` im Klon:

| Stelle | Mutation | gesehene Farbe |
|---|---|---|
| Zeile 131 `make image`, Spalte `Tut was` | Zelle auf 122 und auf 121 Zeichen verlängert (zwei Läufe) | rot: `section-cell-oversized … hat 122 Zeichen, erlaubt sind 120` bzw. `… hat 121 Zeichen, erlaubt sind 120`, je Exit 2 |
| Zeile 131, Spalte `Tut was` | Zelle auf 120 Zeichen | grün, 0 Befunde, Exit 0 |
| Zeile 120 `make sdk-public-doc-check`, Spalte `Vertrag` | Zelle auf 222 Zeichen | rot: `… Spalte "Vertrag" hat 222 Zeichen, erlaubt sind 220`, Exit 2 |
| Zeile 132 `make image-stale`, Spalte `Tut was` | Zelle geleert | rot: `section-cell-undersized … hat 0 Zeichen, verlangt sind 1`, Exit 2 |

Ohne Mutation endet derselbe Lauf im Klon mit `0 Befund(e)`, Exit 0.
`harness/sensors/docs-check.md` nennt die Regel in §Vertrag, als Grenze 11 (sie
misst Länge, nicht ob der Satz die Zeile trägt) und in §Bindung.

**Befund-Liste (Vorlage für §7; nicht berichtigt, der Umzug ändert keinen
Wortlaut).**

1. `make proto-generate`: die Zelle trug ein unmaskiertes `|` im Code-Span
   (`docker run … <image> | tar -x -C .`); die Tabellenzeile zerfiel damit in
   vier Zellen. In der Datei ist der Text kein Tabelleninhalt mehr, der Fehler
   ist mit dem Umzug weg.
2. `make image-cve` (`targets/image.md`): „scheitert bis zum ersten echten
   Release strukturell am fehlenden Ziel-Image (real geprüft — kein
   `ghcr.io/pt9912/pg-change-feed:latest` vorhanden)“ — ein Server-Release ist
   seither erschienen (`make test-sdk-altserver` nennt das Image `:0.5.0`); die
   Aussage ist vermutlich veraltet. Dazu Chronik-Sprache („bis hierhin
   unimplementierte Zusage“).
3. `.github/workflows/hub-description.yml` (`targets/workflows.md`): „der volle
   Erfolgspfad … bleibt bis zum ersten echten Release unbewiesen“ — derselbe
   Stand wie Befund 2.
4. `make sdk-pack-*` (`targets/sdk-pack.md`): „real erzeugt“ nennt
   `PgChangeFeed.Client.0.2.1.nupkg`, `pgchangefeed-0.2.1` und
   `pgchangefeed-kotlin-0.2.2`; `make test-sdk-kompat` misst Packages 0.6.0 —
   die Artefakt-Namen sind ein alter Lauf-Beleg.
5. `make pin-stale-*` (`targets/pin-stale.md`): „Real ausgeführt: P3/P4/P5
   zeigen echten Drift …, P6 ist aktuell“ — Messung ohne Datum und Lauf.
6. `make test-sdk-kotlin-integration` (`targets/sdk-integration.md`): „trägt den
   Python-HTTP-Abschnitt im Abdeckungs-Träger seit
   slice-sdk-python-http-reale2e“ — im Kotlin-Runner unplausibel (der Satz
   steht so in der C#-Zelle mit Kotlin- und Python-Abschnitt); vermutlich
   Kopierfehler.
7. `make handbuch-public-doc-check` nennt „die vier von Runnern geschriebenen
   `*-abdeckung.md`“, `make test-handbuch-public-doc-check` „die fünf
   ausgenommenen Dateien“ — zwei Zahlen für die Ausnahme-Menge; in
   `harness/sensors/handbuch-public-doc-check.md` stehen beide jetzt in
   derselben Datei (Grenze 3 „Die vier Erzeugnisse“ gegen die Fassung im
   Gate-Index). Seit der Fixrunde sagt der Einleitungssatz der Fassung, dass
   bei Abweichung §Vertrag gilt; die inhaltliche Bereinigung ist Folgearbeit.
8. Positionsverweise auf die alte Tabelle stehen jetzt in einer anderen Datei:
   `make example-run-go` „(siehe Zeile darunter)“ — in `targets/examples.md` ist
   es der nächste Abschnitt, stimmt also; `make example-demo-up` „(siehe
   `make schema-rollout`-Zeile oben)“ — in `targets/examples.md` steht kein
   solcher Abschnitt, gemeint ist die README-Zeile.
9. Die 15 ergänzten Dateien tragen ihren Inhalt teils doppelt (eigener Vertrag
   und `## Fassung im Gate-Index`); das Zusammenführen ist eine inhaltliche
   Entscheidung je Datei und nicht Teil dieses Slice. Bis dahin ordnet der
   Einleitungssatz jeder Fassung den Rang: bei Abweichung gilt §Vertrag der
   Datei mit den Abschnitten davor. Die inhaltliche Bereinigung der
   Doppelungen ist Folgearbeit.
10. Vorbestand, nicht durch diesen Slice: `harness/sensors/baseline-verify.md`
    Grenze 3 verweist auf `harness/README.md` §Nicht behauptet, einen Abschnitt,
    den die README nicht führt (Inline-Code, kein Link, deshalb kein
    `anchor-missing`); `.github/workflows/ci.yml` Kopfkommentar zählt die
    Gates von `make gates` mit sechs auf, `harness/sensors/gates.md` nennt zehn.
11. Weitere Chronik-Sprache im umgezogenen Wortlaut (Review F-6):
    `harness/targets/test-integration.md` „inhaltlich über den ursprünglichen
    MVP-Zuschnitt hinausgewachsen“ und `harness/targets/image.md` „jetzt als
    Multi-Arch-Manifestliste“ — nicht berichtigt (§1).

**Träger-Meldungen (§3.13; fremde Dateien, gemeldet, nicht mitgeändert; Frist:
Closure dieses Slice).**

- `tools/coverage-gate.sh` Zeile 2–3: „aktuelle Schwelle und Historie in
  harness/README.md §Sensors“ — die Rampe (70 % → 80 %) steht jetzt in
  `harness/sensors/coverage-gate.md` (§Fassung im Gate-Index) und die Stufe in
  `harness/mk/coverage.mk`; die README-Zelle verweist auf `coverage.mk`.
- `examples/http-client/consumer_test.go` Zeile 34–35,
  `examples/http-client/request_test.go` Zeile 64 und
  `examples/csharp/http-client/HttpClient.Tests/TablesClientTests.cs` Zeile 10:
  „dieselbe Grenze wie beim `natsnotify`-Adapter, `harness/README.md` §Sensors,
  `make test-notify`“ — die Loopback-Grenze steht jetzt in
  `harness/targets/tier-tests.md` (die README-Zeile verlinkt sie).
- [`AGENTS.md`](../../../../AGENTS.md) §4: „dort stehen Target, Vertrag und
  Bindung (inkl. ADR-Links, Schwellen, Carveout-Verweise) vollständig“ — der
  Index ist vollständig, der ausführliche Vertrag steht in der verlinkten Datei.
- `docs/plan/planning/open/slice-baseline-6-14-0-dokumente-nachziehen.md`
  (Liefer-Punkt 3, §3): „`harness/README.md` (Zeile `make pin-stale-baseline`)“ —
  die Zeile ist jetzt ein Satz, der ausführliche Text steht in
  `harness/targets/pin-stale.md` (Abschnitt `make pin-stale-baseline`); der
  Bump-Slice ändert dort.
- [`AGENTS.md`](../../../../AGENTS.md) §3.14: „die zentrale Wache …
  (`tools/schema/rolloutguard`, siehe `harness/README.md` §Sensors,
  `make schema-rollout`-Zeile)“ — die README-Zeile nennt die Wache nicht mehr,
  sie steht in `harness/targets/schema-rollout.md` (die Zeile verlinkt sie).
- `harness/mk/coverage.mk` Zeile 11–14: „Kalibrierungs-Bindung
  (harness/README.md §Sensors …) … die uebrigen Traeger nennen nur die Rampe
  (Einstieg 70 %, Endstufe 80 %)“ — die README-Zelle nennt die Rampe nicht mehr;
  sie steht in `harness/sensors/coverage-gate.md`.
- Die `Accepted`-ADRs, die eine Zeile der Tabelle zitieren (z. B.
  [`ADR-0044`](../../adr/0044-image-beleg-semantik.md) „`make image`-Werkzeuge-Zeile
  in `harness/README.md` tragen die korrigierte …“), frieren den Stand ihrer
  Zeit ein; die Zeilen bestehen weiter, kein Nachzug (`AGENTS.md` §3.5).

### Fixrunde (Review `review-slice-harness-readme-zellen-kuerzen`, 2026-10-05)

- **F-1 (MEDIUM):** jede der 15 Fassungen trägt hinter dem Einleitungssatz
  „Der Text darunter ist der wortgleich umgezogene Index-Text, kein eigener
  Vertrag: weicht er von dieser Datei ab, gilt §Vertrag (Link `#vertrag`) mit
  den Abschnitten bis zu diesem.“ — alle 15 Dateien führen `## Vertrag`, der
  Anker ist derselbe. Der umgezogene Wortlaut ist unverändert (Gegenprobe
  unten); Befunde 7 und 9 nennen die Bereinigung als Folgearbeit. Die Aussage
  „jede ergänzte Datei byte-gleich mit alter Inhalt + Abschnitt“ (§3 oben) gilt
  damit für den Stand `76c3b9e0`; seither ist je Datei dieser eine Satz
  hinzugekommen.
- **F-2 (LOW):** `AGENTS.md` §3.14 und `harness/mk/coverage.mk` stehen unter
  Träger-Meldungen; nicht geändert, Frist Closure.
- **F-3 (LOW):** §1 verweist am Ausschluss „Neuschnitt“ auf Plan-Nachzug 1.
- **F-4 (INFO):** die Bindung von `make doc-tracked` zeigt auf
  `§Fassung im Gate-Index` statt `§Bindung`.
- **F-5 (INFO):** Grenze 11 in `harness/sensors/docs-check.md` nennt, dass
  `Bindung` keine Höchstlänge hat, und kennzeichnet den Satz über eine neue
  Tabelle als hergeleitet.
- **F-6 (INFO):** Befund 11.
- **F-7 (INFO):** in der Zeile `make doc-trace` steht der Datei-Link hinter
  „wie `make image-stale`“; der Kopf `**Bezug:**` verlinkt
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) auf die ADR-Datei
  mit ihrem Titel.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1; am Planungsstand steht dort `slice-sdk-tls-optionen`, das README-Zeilen
ändert — dessen Closure geht voran, sonst Konflikte im selben Tabellenblock),
**und** der Slice `baseline-6-14-0-dokumente-nachziehen` ist entweder done oder
bewusst hinter diesem eingeordnet (Reihenfolge benannt, §1), **und** der
Implementer hat die Startmessung gefahren: die Suchlauf-Zeilen am Arbeitsstand
neu gemessen (neue Commit-Kennung als Parent, Exit 0 nach Anpassung; jede
Abweichung mit Ursache im Bericht).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Umzug umfasst so
  viele Gruppen, dass ein Slice sie nicht in einem Review trägt — dann Teilung nach
  Tabelle (Sensors / Werkzeuge) oder nach Familie, mit der d-check-Regel im
  letzten Teil-Slice.
- `in-progress` → `open` (blockiert): ein paralleler Slice ändert Zeilen der
  README und lässt sich nicht vorziehen oder abwarten — Frage an den Auftraggeber;
  kein Carveout, weil kein Gate rot ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit
Exit 0 (die `diff`-Zeilen gesetzt, jede Trefferzeile gelesen), und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; Kandidat: die neue `structure`-Regel als Sensor mit der
benannten Grenze, dass sie Länge misst und nicht, ob der Satz die Zeile trägt).
Ein Gate, das am Stand der Closure rot ist, geht nur mit dokumentiertem Carveout
nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang.

- **Große Umzugsmenge → Wortlaut-Verlust:** rund 70 Zellen, davon eine mit
  26 KB; beim Umkopieren fällt ein Satz, eine Zahl oder ein Herkunftsvermerk weg.
  Gegenmaßnahme: Zeichenzahl-Gegenprobe je Zelle und Suchlauf der Eigennamen
  (§2 Liefer-Punkt 1), Review prüft Stichprobe Wort für Wort. —
  **Ausgang:** offen bis zur Closure.
- **Abhängigkeit zu `baseline-6-14-0-dokumente-nachziehen`:** beide lesen oder
  berühren `harness/README.md` (der Bump-Slice den Abgleich und die Zeile
  `make pin-stale-baseline`). Gegenmaßnahme: Reihenfolge ist benannt (§1 und §4),
  kein gleichzeitiger Lauf (WIP-Limit 1). — **Ausgang:** offen bis zur Closure.
- **Parallele Slices ändern README-Zeilen:** `slice-sdk-tls-optionen` läuft in
  `in-progress/` und kann Zeilen der Tabellen berühren (neue Target-Beschreibung);
  ein Rebase in einen umgebauten Block erzeugt Konflikte oder lässt eine lange
  Zelle zurück, die die neue Regel rot färbt. Gegenmaßnahme: Start erst nach dessen
  Closure (§4); die Regel kommt zuletzt. — **Ausgang:** offen bis zur Closure.
- **Die Schwellen 120/220 sind Vorschlag** (angelehnt an andere Register-Regeln in
  `.d-check.yml`, nicht gemessen als passend): zu knapp führt zu Pflichterfüllung
  durch Satzzerhacken. Gegenmaßnahme: der Implementer nennt die längste verbleibende
  Zelle; Review liest, ob der Satz die Zeile trägt. — **Ausgang:** offen bis zur
  Closure.
- **Der Auftraggeber-Wert „Regel adressiert beide Tabellen“ ist hergeleitet:**
  adressiert die Regel nur die erste Tabelle des Abschnitts, ist die Werkzeuge-Tabelle
  ungeschützt. Gegenmaßnahme: Mutationsprobe auf der Kopie (§2 Liefer-Punkt 3);
  sonst eine zweite Regel mit eigenem Abschnitt oder eine Unterüberschrift. —
  **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Ging der Gegenstand an einen anderen Slice oder entfiel er,
trägt diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund.

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**.

**Der Abschnitt selbst entfällt nie.**

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt **eine** Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`harness/`, `.d-check.yml`) liegen alle in ihr, es entsteht keine
neue Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) für die berührte Sub-Area:
vom Planungsstand nicht durchgesehen; der Implementer hat es zum Start gesichtet
(2026-10-05, gemessen: `grep -l -i -E 'zelle|zellen|README\.md|structure|register-spalt|Zeilenl'
docs/plan/planning/observations/BEO-PGC/*/observation.md`, 15 Dateien, jede
gelesen; Zähler = Zahl der Dateien in `evidence/`). Treffer zum Gegenstand:

- `BEO-PGC/bindung-spalte-uneinheitlich-tief` — **offen, 1×**: die
  Bindung-Spalte der Sensors-Tabelle verweist unterschiedlich tief (ADR inline
  oder nur über die Sensor-Datei). Dieser Slice berührt sie: jede gekürzte Zeile
  verlinkt ihre Datei, die ADR-Links der Zelle bleiben stehen; die Asymmetrie
  bleibt bestehen (Formfrage, kein Gegenstand dieses Slice). Mit diesem Slice
  nicht erneut aufgetreten, kein zweites Auftreten.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, Deckel) und
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, Deckel) — die Klasse
  der Befunde 2 bis 5 und der Träger-Meldungen in §3 Umsetzung; kein neuer
  Mechanismus nötig, der Suchlauf trägt sie.
- `BEO-PGC/regel-weiter-als-ihr-sensor` (verkörpert, teilweise) — Nachbarklasse
  zur Grenze der neuen Regel (sie misst Länge, nicht den Satz).

Kein Eintrag zur Zellenlänge der `harness/README.md` selbst; keiner erreicht mit
diesem Slice 3×.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
