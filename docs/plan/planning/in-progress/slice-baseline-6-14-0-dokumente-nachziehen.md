# Slice baseline-6-14-0-dokumente-nachziehen: die aus der Baseline abgeleiteten Repo-Dokumente entsprechen v6.14.0, und der nächste Bump trägt einen Vergleichs-Schritt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice; die Roadmap führt wellenlose Arbeit nicht (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht, Regel „Wellenlose Arbeit
erscheint nicht in der Roadmap“).

**Bezug:** [`ADR-0051`](../../adr/README.md) (Baseline-Bump, Pin-Inventar P8, Entscheidung 7),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen).
Keine `LH-*`-Anforderung ist berührt: der Slice ändert Harness-Dokumente, nicht das Produkt.

**Berührte Spec-Stellen:** — (keine; `harness/conventions.md` §Baseline nennt den
Stand v6.14.0 und wird nur gelesen).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag). **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass:** Der Bump von v6.13.0 auf v6.14.0 hat bisher Verweise und Pins
umgestellt (Commits `f5b2840a`, `675246dd`, `9b360010`), nicht die vom Repo
geführten, aus der Baseline abgeleiteten Dokumente geprüft. Gemessen am Parent
`9b360010` (Suchlauf, §3): `docs/plan/planning/README.md` trägt zwei Mal den
Platzhalter `welle-<NN>-results` (die Vorlage `welle-<Kennung>-results`) und
nicht die Klausel der Vorlage-Zeile `done/` („oder Gegenstand an einen anderen
Slice übergegangen oder entfallen …“). Beide Abweichungen sind **älter als der
Bump** (die Vorlage ist zwischen 6.13.0 und 6.14.0 an diesen Stellen
unverändert; Angabe des Auftraggebers, **übernommen**, im Slice nachzumessen
durch `diff` der beiden vendored Stände) und nie entdeckt worden — der Bump-Ablauf
hat keinen Vergleichs-Schritt.

**Ziel:** Die aus der Baseline abgeleiteten Repo-Dokumente sind gegen die
vendored Vorlagen und das Regelwerk v6.14.0 abgeglichen (je Dokument „entspricht“
oder „Abweichung mit Beleg“, die Abweichungen der Planungs-README behoben), und
der Bump-Ablauf trägt einen festgeschriebenen Vergleichs-Schritt, der den
nächsten Bump davor bewahrt, dass das Repo hinter der Baseline zurückbleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Roadmap-Inhalt jenseits der Regel** (OTLP-Metrik-Export, TLS/Mehrfach-Token,
  Bench-Umstellung PER-001, SDK-TLS, Baseline-Bump als Zeilen). Das Regelwerk
  (`modul-06-roadmap.md`, „Wellenlose Arbeit erscheint nicht in der Roadmap“)
  verlangt für wellenlose Slices keine Zeile; der Slice prüft nur, ob eine
  *Welle* oder ein *Meilenstein-Beleg* der Regel nach fehlt, und notiert sonst
  „kein Nachzug nötig“ mit Beleg. Eine Roadmap-Neuordnung wäre ein anderer Vorgang.
- **Arbeit von Agent `ab8c3e73293b11d22`** (Reviewer-Skill und sechs Reports auf
  v6.14.0; läuft parallel). Der Slice **verifiziert** sein Ergebnis (Liefer-Punkt 2),
  er wiederholt es nicht: Doppelarbeit prüfte denselben Kontext zweimal, ohne
  zweiten Leser.
- **SDK- und Server-Releases** — kein Bezug zum Gegenstand; Releases brauchen
  eigene Freigabe.
- **Ein neues Gate für den Vergleich**, falls Liefer-Punkt 3 auf eine
  Verfahrensregel statt auf ein Werkzeug fällt. Entsteht ein Gate, ist die
  Aufnahme in `make gates` eine Schwellen-/Gate-Entscheidung (`AGENTS.md` §3.6) und
  geht als **Übergabe an den Architect** (ADR), nicht in diesen Slice.
- **Abweichungen, die der Abgleich in Liefer-Punkt 2 über `docs/plan/planning/README.md`
  hinaus findet** (z. B. in `harness/conventions.md`, ADR-Index, Carveout-Verzeichnis):
  sie werden als Befund mit Beleg und als Folge-Slice benannt, nicht mitgeändert —
  sonst wächst der Slice über die drei Liefer-Punkte.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Alle Beleg-Angaben dieser Liste sind **Zusagen** („zu belegen durch …“): der
Planungsstand hat keinen der Läufe gefahren
([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) Instanz B).

- [ ] **Planungs-README und Roadmap (Liefer-Punkt 1).**
      `docs/plan/planning/README.md` gegen
      `.harness/baseline/v6.14.0/templates/docs/plan/planning/README.template.md`
      abgeglichen und angeglichen: (a) die Zeile `done/` trägt die Klausel „oder
      Gegenstand an einen anderen Slice übergegangen oder entfallen: §7 nennt
      Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer (Baseline-Regelwerk
      `modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
      übernimmt)“; (b) `welle-<NN>-results.md` heißt an beiden Stellen
      `welle-<Kennung>-results.md` (MR-002); jede verbleibende, **repo-spezifische**
      Abweichung (z. B. der `d-check:ignore`-Kommentar am Carveout-`done/`-Pfad)
      steht mit Grund im Bericht. Der Beleg ist `diff` Vorlage ↔ Repo-Datei nach
      Normalisierung der Platzhalter, zu belegen durch die gedruckte Diff-Länge.
      `roadmap.md`: Struktur gegen die Vorlage gemessen (Auftraggeber: nur Inhalt
      weicht ab, **übernommen**); Befund „kein Nachzug nötig“ mit Beleg
      (`modul-06-roadmap.md`: wellenlose Arbeit erscheint nicht in der Roadmap;
      erwartet: es fehlt keine Welle und kein Meilenstein-Beleg), sonst nur das
      Nachzuziehende, das die Regel verlangt.
- [ ] **Verifikation der Reviewer-Dateien und Abgleich der übrigen Träger
      (Liefer-Punkt 2).** (a) Die Ergebnisse von Agent `ab8c3e73293b11d22` gemessen:
      `.harness/skills/reviewer.md` ↔ `reviewer.template.md`,
      `.harness/skills/closure-note-reviewer.md` ↔ `closure-note-reviewer.template.md`,
      `.claude/agents/reviewer.md` und `verifier.md` ↔ Regelwerk
      `modul-10-review-harness.md` v6.14.0, die sechs Reports ↔
      `review-report.template.md`; je Datei „entspricht“ oder „Abweichung mit Beleg“.
      (b) `modul-15-observability.md` (Audit-Span-Schema: ein Pflicht-Feld, dessen
      Wert die Quelle nicht liefert, bleibt Pflicht und wird als „nicht bekannt“
      samt Quelle gekennzeichnet, nicht `0`/`false`/„keine Rolle“) gegen
      `.claude/hooks/span-emit.sh` und die Doku zum Audit-Span-Schema des Repos
      geprüft; **erwartet** (aus der Lesung des Wrappers hergeleitet, nicht
      gemessen): der Wrapper erzeugt keinen Span selbst, er ruft den Träger
      `.harness/state/bin/ai-harness-init span-emit`, das Schema liegt also im
      Träger, nicht im Repo — dann „Repo erfüllt bereits / nicht Gegenstand dieses
      Repos“ mit Beleg, sonst Anpassung. (c) Abgleich `harness/README.md` ↔
      `harness/README.template.md`, `harness/conventions.md` ↔
      `conventions.template.md`, `docs/plan/adr/README.md` ↔ `adr/README.template.md`,
      `docs/plan/carveouts/` ↔ Carveout-Vorlage, Slice- und Welle-Vorlage ↔ Praxis:
      nur Bestandsaufnahme, je Dokument „entspricht“ oder „Abweichung mit Beleg“,
      Abweichungen als benannte Folgepunkte (§1, vierter Ausschluss).
- [ ] **Bump-Vergleichs-Schritt (Liefer-Punkt 3).** Festgelegt, ob der Vergleich
      Repo-Dokument ↔ Vorlage ein **Werkzeug** (Docker-only, netzlos,
      Normalisierung von Versionsstrings und Platzhaltern) oder eine
      **Verfahrensregel** ist, mit Begründung im Plan-Nachzug; mindestens als
      Pflichtschritt im Bump-Ablauf festgeschrieben: „vor dem Löschen der alten
      Baseline `diff` alte ↔ neue Baseline und jede inhaltliche Änderung gegen die
      Repo-Dokumente prüfen“ — in `harness/sensors/baseline-verify.md` und/oder an
      der Stelle des Bump-Ablaufs (Regelwerk `modul-02-harness-bootstrap.md`
      §Freshness-Audit, `harness/targets/pin-stale.md` Abschnitt
      `make pin-stale-baseline` — die Index-Zeile in `harness/README.md` ist ein
      Satz und bleibt ≤ 120 Zeichen, `structure`-Regel in `.d-check.yml`; der
      Implementer nennt die Stelle nach Lesung). Entsteht ein Gate, endet der Slice
      an dieser Stelle mit der Übergabe an den Architect; die Verfahrensregel ist
      der Mindestumfang und braucht keine ADR.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make suchlauf-nachmessen PLAN=` mit diesem Plan Exit 0 am Endstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durch den Verifier (Belege, nicht Behauptung); Liefer-Punkt 2
      prüft er unabhängig von Agent `ab8c3e73293b11d22`.
- [ ] Doku-Update: [`harness/targets/pin-stale.md`](../../../../harness/targets/pin-stale.md)
      bzw. [`harness/sensors/baseline-verify.md`](../../../../harness/sensors/baseline-verify.md),
      falls Liefer-Punkt 3 `make pin-stale-baseline` oder `make baseline-verify`
      berührt (die Index-Zeile in [`harness/README.md`](../../../../harness/README.md)
      §Sensors nur, wenn sich ihr Kurzsatz ändert); gemeldete Träger fremder Dateien mit der
      Closure nachgezogen (§3 Suchlauf, [`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, solange die Roadmap unter *Offene Wellen* keine Welle führt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/planning/README.md` | update | Klausel in der Zeile `done/`, `welle-<Kennung>-results.md` an zwei Stellen (Liefer-Punkt 1) |
| `docs/plan/planning/in-progress/roadmap.md` | kein Eingriff erwartet | Struktur-/Regelprüfung, Befund im Bericht (Liefer-Punkt 1) |
| `.harness/skills/reviewer.md`, `.harness/skills/closure-note-reviewer.md`, `.claude/agents/reviewer.md`, `.claude/agents/verifier.md`, `docs/reviews/` (sechs Reports) | nur lesen | Verifikation der Arbeit von Agent `ab8c3e73293b11d22` (Liefer-Punkt 2a) |
| `.claude/hooks/span-emit.sh`, `.claude/settings.json` | nur lesen | Abgleich Audit-Span-Schema (Liefer-Punkt 2b) |
| `harness/README.md`, `harness/conventions.md`, `docs/plan/adr/README.md`, `docs/plan/carveouts/` | nur lesen | Bestandsaufnahme gegen die Vorlagen (Liefer-Punkt 2c) |
| `harness/sensors/baseline-verify.md` und/oder `harness/targets/pin-stale.md` (Abschnitt `make pin-stale-baseline`) | update | Pflichtschritt „diff alte ↔ neue Baseline“ (Liefer-Punkt 3) |
| `tools/harness/…` (Vergleichs-Werkzeug) | neu, **nur falls** die Entscheidung auf Werkzeug fällt | Docker-only, netzlos; Gate-Aufnahme geht an den Architect |

**Entscheidung Werkzeug vs. Verfahrensregel (Vorab-Position, vom Implementer zu
messen).** Ein Vergleich Repo-Dokument ↔ Vorlage ist nach Normalisierung
mechanisierbar, aber die Repo-Dokumente sind **ausgefüllte** Vorlagen: der Inhalt
weicht beabsichtigt ab (Roadmap, Konventionen), nur die *Struktur* und feste
Klauseln sollen übereinstimmen. Ein generischer `diff` meldete deshalb Rauschen;
ein Strukturvergleich (Überschriften, feste Klauseln) wäre ein neues Werkzeug mit
eigenem Pflegeaufwand und der Gefahr, Pflichterfüllung statt Wahrheit zu prüfen.
Die Änderung zwischen zwei **Baseline-Ständen** (`diff` alt ↔ neu) ist dagegen
klein und von einer Rolle lesbar. Position: **Verfahrensregel** (Pflichtschritt
im Bump-Ablauf) als Mindestumfang; ein Werkzeug nur, wenn der Implementer am
Bestand (Liefer-Punkt 2c) zeigt, dass der Schritt ohne Werkzeug nicht trägt —
dann Übergabe an den Architect. Die Position ist hergeleitet, nicht gemessen.

**Reihenfolge:** (1) Startmessung (Suchlauf unten, Parent als Commit-Kennung neu
setzen); (2) `diff` v6.13.0 ↔ v6.14.0 der Vorlagen/Regelwerk-Dateien, soweit die
alte Baseline noch lesbar ist (Git-Historie, ggf. `git show <Commit>:.harness/baseline/<alter Tag>/…`),
um die Behauptung „Vorlage unverändert“ nachzumessen; (3) Liefer-Punkt 1;
(4) Liefer-Punkt 2; (5) Liefer-Punkt 3; (6) `make gates`.

**Suchlauf (§3.13 der Regeln, [`AGENTS.md`](../../../../AGENTS.md)).** Bewegte
Eigenschaften: (A) die Schreibweise des Welle-Platzhalters in
`docs/plan/planning/README.md` (Symbol `welle-<NN>-results`); (B) die Klausel der
Zeile `done/` (Beschreibung „gegenstand an einen anderen Slice übergegangen“,
Hedge: kein Treffer heißt: fehlt); (C) der Baseline-Pin `v6.13.0` in getrackten
Dateien (Zählwort-Art: Versionsstring); (D) das Audit-Span-Schema
(`span-emit`, `Audit-Span`, `requirement.id`, `slice.id`). Parent ist der Stand
`9b360010…` (`git rev-parse HEAD` am Planungsstand; nie `HEAD`). Suchraum: ganzer
Baum außer `docs/reviews`, `.harness/baseline` (vendored, hält seinen Stand),
`docs/plan/planning/done` (Records) und `docs/plan/planning/observations`
(Register hält Belege mit dem Stand ihrer Zeit); Zeilen 1 und 5 sind auf die
Planungs-README bzw. die Vorlage eingeschränkt, weil nur dort der Träger liegt.
**Parent-Zeilen gemessen** (Planungsstand 2026-10-05); Zeilen mit `diff` setzt der
Implementer am Endstand und liest jede Trefferzeile. **Erwartung am Endstand**
(hergeleitet, nicht gemessen): Zeile 1 sinkt von 7 auf 5 (die zwei Treffer in der
Planungs-README gehen weg; die übrigen — `close-welle.md` ×2, `implement-slice.md`,
`reviewer.md`, `AGENTS.md`, mit `welle-<NN>` im Wortlaut — sind eigene Träger und
werden gelesen, ob sie MR-002 widersprechen); Zeile 3 bleibt 1
(`docs/plan/adr/0095…`, `Accepted`, unberührbar); Zeile 4 steigt auf 1; Zeile 2
sinkt auf 0.

```suchlauf
9b3600104203cba0f1eeea63902517066a98c26f 7 -n -F 'welle-<NN>' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
9b3600104203cba0f1eeea63902517066a98c26f 2 -n -F 'welle-<NN>-results' -- docs/plan/planning/README.md
9b3600104203cba0f1eeea63902517066a98c26f 1 -n -E 'v6\.13\.0' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
9b3600104203cba0f1eeea63902517066a98c26f 0 -n -F 'Gegenstand an einen anderen Slice übergegangen' -- docs/plan/planning/README.md
9b3600104203cba0f1eeea63902517066a98c26f 1 -n -F 'Gegenstand an einen anderen Slice übergegangen' -- .harness/baseline/v6.14.0/templates/docs/plan/planning/README.template.md
9b3600104203cba0f1eeea63902517066a98c26f 7 -n -E 'span-emit|Audit-Span|requirement\.id|slice\.id' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
```

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1; am Planungsstand steht dort `slice-sdk-tls-optionen`), **und** Agent
`ab8c3e73293b11d22` hat seine Arbeit an Reviewer-Skill und Reports committet
(sonst liest Liefer-Punkt 2a einen halben Stand), **und** der Implementer hat die
Startmessung gefahren: die Suchlauf-Zeilen aus §3 am Arbeitsstand neu gemessen
(neue Commit-Kennung als Parent, `make suchlauf-nachmessen PLAN=` Exit 0 nach
Anpassung; jede Abweichung mit Ursache im Bericht).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Abgleich in
  Liefer-Punkt 2c findet so viele Abweichungen, dass ihre Behebung ein vierter
  Liefer-Punkt würde, oder Liefer-Punkt 3 verlangt ein Werkzeug (dann Teilung in
  „Dokument-Nachzug“ und „Vergleichs-Werkzeug“).
- `in-progress` → `open` (blockiert — Carveout?): das Ergebnis von Agent
  `ab8c3e73293b11d22` ist nicht committet oder widerspricht der Vorlage so, dass
  eine Entscheidung nötig ist; kein Carveout, weil dann kein Gate rot ist — Frage
  an den Auftraggeber.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, `make suchlauf-nachmessen PLAN=` mit diesem Plan endet am Endstand mit
Exit 0 (die `diff`-Zeilen gesetzt, jede Trefferzeile gelesen), und die
Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder
benannte Spec-Lücke; Kandidat: der Pflichtschritt „diff alte ↔ neue Baseline“ im
Bump-Ablauf als geschärfte Regel mit Träger, oder die benannte Lücke, dass kein
Sensor ausgefüllte Vorlagen gegen ihre Vorlage vergleicht). Ein Gate, das am Stand
der Closure rot ist, geht nur mit dokumentiertem Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Agent `ab8c3e73293b11d22` läuft parallel** und schreibt in
  `.harness/skills/` und `docs/reviews/`; ein Abgleich vor seinem Commit liest
  einen Zwischenstand, ein gleichzeitiger Eingriff erzeugt Konflikte. Gegenmaßnahme:
  Start erst nach seinem Commit (§4), der Slice liest diese Dateien nur. —
  **Ausgang:** offen bis zur Closure (eingetreten / entfallen).
- **Das Abgleich-Werkzeug wird zu groß:** der Strukturvergleich ausgefüllter
  Vorlagen wächst zu einem eigenen Programm mit Pflegeaufwand und einem Gate, das
  Form statt Wahrheit prüft. Gegenmaßnahme: Verfahrensregel als Mindestumfang, Werkzeug
  nur mit Beleg und Übergabe an den Architect (§3). — **Ausgang:** offen bis zur
  Closure (entfallen: Verfahrensregel genügt · eingetreten: Folge-Slice mit
  Architect-Übergabe).
- **Die Aussage „Vorlage zwischen 6.13.0 und 6.14.0 unverändert“ ist übernommen**
  (Auftraggeber), nicht gemessen; stimmt sie nicht, stammt eine Abweichung doch aus
  dem Bump. Gegenmaßnahme: Reihenfolge Schritt (2) misst sie per `diff` der
  Git-Stände. — **Ausgang:** offen bis zur Closure.
- **Der Bestandsabgleich in Liefer-Punkt 2c findet mehr, als der Slice trägt.**
  Gegenmaßnahme: nur Bestandsaufnahme, Abweichungen als benannte Folge-Slices
  (§1). — **Ausgang:** offen bis zur Closure (weiter offen: → Folge-Slice in `open/`).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*Der Plan füllt diese Sektion nicht; sie wird bei der Closure vor dem
`git mv` nach `done/` geschrieben.*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration führt **eine** Sub-Area (`*`, Kürzel `PGC`, Greenfield); die
berührten Pfade (`docs/plan/planning/`, `harness/`, `.harness/skills/`,
`.claude/`) liegen alle in ihr, es entsteht keine neue Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) für die berührte Sub-Area
durchgegangen; Treffer ohne Zähler-Nachmessung (**übernommen**, nicht gezählt):
`zitat-korrektur-reichweite-abschnitte-kurzform` (Beleg
`evidence/slice-harness-baseline-v6-13-0.md`, betrifft den vorigen Bump und die
Reichweite von Zitat-Korrekturen — vom Implementer zu lesen, bevor Liefer-Punkt 3
den Bump-Ablauf berührt). Sonst keine Treffer, die der Planungsstand benennen kann.

**Modus:** alle berührten Sub-Areas GF (`*`/`PGC`, Greenfield).
