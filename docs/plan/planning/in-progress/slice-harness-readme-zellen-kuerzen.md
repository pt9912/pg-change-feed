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

**Bezug:** [`ADR-0051`](../../adr/README.md) (Harness-Pflege und Pin-Inventar),
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

- [ ] **Zuschnitt und Umzug (Liefer-Punkt 1).** Je Ziel oder Ziel-Gruppe der
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
- [ ] **Zellen kürzen (Liefer-Punkt 2).** In `harness/README.md` ist jede Zelle
      „Tut was“ ≤ 120 und jede Zelle „Vertrag“ der Sensors-Tabelle ≤ 220 Zeichen
      (gemessen wie dargestellt, inkl. Markdown-Syntax); die Bindung-Zelle
      verlinkt die Datei in Link-Form (kein Pfad im Fließtext); die Index-Zeilen
      bleiben vollständig (Zeilenzahl der Tabellen vorher = nachher). Die Tabellen
      „Source precedence“, „Guides“ und die Prosa-Abschnitte sind unverändert
      (`git diff` zeigt keine Änderung dort). Träger, die Zeilenzahl oder Position
      der README zitieren, sind nachgezogen (Suchlauf unten, §3.13 der Regeln).
- [ ] **Regel, die die Grenze hält (Liefer-Punkt 3).** `.d-check.yml` trägt unter
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

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make suchlauf-nachmessen PLAN=` mit diesem Plan Exit 0 am Endstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
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

```suchlauf
4db7cf958264705b9d18e727e48f95361095bad4 15 -n -E '^.{2001,}$' -- harness/README.md
4db7cf958264705b9d18e727e48f95361095bad4 57 -n -E '^\| `[^|]*\| [^|]{121,} \| ' -- harness/README.md
4db7cf958264705b9d18e727e48f95361095bad4 55 -n -E '^\| `[^|]*\| [^|]{221,} \| ' -- harness/README.md
4db7cf958264705b9d18e727e48f95361095bad4 5 -n -E '^# `' -- harness/targets/*.md
4db7cf958264705b9d18e727e48f95361095bad4 113 -n -F 'harness/README.md' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations :!docs/plan/planning/open/slice-harness-readme-zellen-kuerzen.md
4db7cf958264705b9d18e727e48f95361095bad4 48 -n -E 'Sensors \(Feedback-Gates\)|§Sensors' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations :!docs/plan/planning/open/slice-harness-readme-zellen-kuerzen.md
```

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
vom Planungsstand nicht durchgesehen — der Implementer sichtet es zum Start auf
Einträge zu Zellenlänge, `harness/README.md` und Register-Regeln und nennt Treffer
mit Zähler-Stand (offene Aufgabe, keine Aussage über den Bestand).

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
