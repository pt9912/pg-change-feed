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

- [x] **Planungs-README und Roadmap (Liefer-Punkt 1).**
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
- [x] **Verifikation der Reviewer-Dateien und Abgleich der übrigen Träger
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
- [x] **Bump-Vergleichs-Schritt (Liefer-Punkt 3).** Festgelegt, ob der Vergleich
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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
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
- [ ] Ruhe-Marker unter *Offene Wellen* in `docs/plan/planning/in-progress/roadmap.md`
      nach dem `git mv` dieses Slice nach `done/` wieder eingesetzt (Wortlaut aus
      `d79b7ebd`), sofern `in-progress/` dann keinen Slice trägt; die Suchlauf-Zeile
      `diff 0 … 'Nichts in Arbeit'` nachgezogen (Review F-3).
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
| `docs/plan/planning/in-progress/roadmap.md` | update *(Plan-Nachzug: statt „kein Eingriff erwartet“)* | Ruhe-Marker unter *Offene Wellen* entfernt: er stand, während `in-progress/` diesen Slice trägt (`modul-06-roadmap.md` §Roadmap-Struktur, Bullet *Offene Wellen*); Beleg unter *Belege, Liefer-Punkt 1* (Liefer-Punkt 1) |
| `.harness/skills/reviewer.md`, `.harness/skills/closure-note-reviewer.md`, `.claude/agents/reviewer.md`, `.claude/agents/verifier.md`, `docs/reviews/` (sechs Reports) | nur lesen | Verifikation der Arbeit von Agent `ab8c3e73293b11d22` (Liefer-Punkt 2a) |
| `.claude/hooks/span-emit.sh`, `.claude/settings.json` | nur lesen | Abgleich Audit-Span-Schema (Liefer-Punkt 2b) |
| `harness/README.md`, `harness/conventions.md`, `docs/plan/adr/README.md`, `docs/plan/carveouts/` | nur lesen | Bestandsaufnahme gegen die Vorlagen (Liefer-Punkt 2c) |
| `harness/targets/pin-stale.md` (Abschnitt `make pin-stale-baseline`, neuer Unterabschnitt *Bump-Ablauf: Vergleich vor dem Löschen der alten Baseline*) | update | Pflichtschritte Delta, Stichprobe gegen den Bestand, Ergebnis je Dokument (Liefer-Punkt 3) |
| `harness/sensors/baseline-verify.md` (§Grenze, Punkt 2) | update | ein Satz mit Verweis auf den Bump-Ablauf; der Vertrag bleibt (Liefer-Punkt 3) |
| `tools/harness/…` (Vergleichs-Werkzeug) | **nicht realisiert** *(Plan-Nachzug)* | die Entscheidung fiel auf die Verfahrensregel, Begründung unter *Belege, Liefer-Punkt 3* |

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

**Startmessung (Implementer, 2026-10-05).** Die sechs Zeilen am Stand
`d79b7ebd54ba353e484f654a16c5abb0859b3743` (Übergang nach `in-progress/`, Parent
dieses Laufs) neu gemessen: gedruckt 7 · 2 · 1 · 0 · 1 · 7, gleich den Zahlen am
Planungsstand `9b360010…` — die 33 Commits dazwischen berühren keine der sechs
Eigenschaften; keine Abweichung. Der Parent der Zeilen ist damit `d79b7ebd…`.
Startbedingungen: `in-progress/` trug am Parent nur `roadmap.md` und diesen Slice
(`ls docs/plan/planning/in-progress`); die Arbeit von Agent `ab8c3e73293b11d22`
ist committet in `f5b2840a` (Report-Gerüst aus der Baseline, Repo-Kopie der
Review-Vorlage entfernt), `675246dd` (Reviewer-Skills) und `9b360010` (sechs
Reports); `484d20ec` ist der Pin-Wechsel des Bumps selbst, nicht seine Arbeit —
gemessen mit `git show --stat` je Commit, Arbeitsbaum am Parent sauber
(`git status`). Neu hinzu kommen zwei bewegte Eigenschaften dieses Laufs: (E) der
Ruhe-Marker der Roadmap (Wortlaut `Nichts in Arbeit`) und (F) der Bump-Ablauf
(Beschreibung `Stichprobe gegen den Bestand` unter `harness/`). Endstand
(`diff`-Zeilen, jede Trefferzeile gelesen): Zeile 1 sinkt auf 5 — die übrigen
Treffer `.claude/commands/close-welle.md` ×2, `.claude/commands/implement-slice.md`,
`.harness/skills/reviewer.md` und `AGENTS.md` tragen `welle-<NN>` als Platzhalter
der Herkunfts-Anker-Form `seit welle-<NN>` und sind fremde Träger (Befund 5
unten, gemeldet, nicht mitgeändert); Zeile 3 bleibt 1 (`ADR-0095` §Geschichte,
`Accepted`); Zeile 6 bleibt 7, alle in `.claude/hooks/span-emit.sh` und
`.claude/settings.json`.

```suchlauf
d79b7ebd54ba353e484f654a16c5abb0859b3743 7 -n -F 'welle-<NN>' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
diff 5 -n -F 'welle-<NN>' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
d79b7ebd54ba353e484f654a16c5abb0859b3743 2 -n -F 'welle-<NN>-results' -- docs/plan/planning/README.md
diff 0 -n -F 'welle-<NN>-results' -- docs/plan/planning/README.md
d79b7ebd54ba353e484f654a16c5abb0859b3743 1 -n -E 'v6\.13\.0' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
diff 1 -n -E 'v6\.13\.0' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
d79b7ebd54ba353e484f654a16c5abb0859b3743 0 -n -F 'Gegenstand an einen anderen Slice übergegangen' -- docs/plan/planning/README.md
diff 1 -n -F 'Gegenstand an einen anderen Slice übergegangen' -- docs/plan/planning/README.md
d79b7ebd54ba353e484f654a16c5abb0859b3743 1 -n -F 'Gegenstand an einen anderen Slice übergegangen' -- .harness/baseline/v6.14.0/templates/docs/plan/planning/README.template.md
d79b7ebd54ba353e484f654a16c5abb0859b3743 7 -n -E 'span-emit|Audit-Span|requirement\.id|slice\.id' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
diff 7 -n -E 'span-emit|Audit-Span|requirement\.id|slice\.id' -- . :!docs/reviews :!.harness/baseline :!docs/plan/planning/done :!docs/plan/planning/observations
d79b7ebd54ba353e484f654a16c5abb0859b3743 1 -n -F 'Nichts in Arbeit' -- docs/plan/planning/in-progress/roadmap.md
diff 0 -n -F 'Nichts in Arbeit' -- docs/plan/planning/in-progress/roadmap.md
d79b7ebd54ba353e484f654a16c5abb0859b3743 0 -n -F 'Stichprobe gegen den Bestand' -- harness
diff 1 -n -F 'Stichprobe gegen den Bestand' -- harness
```

### Belege (Implementer)

Arbeitsverzeichnis der Normalisierungen: Scratchpad des Laufs (Kopien, keine
Repo-Datei geschrieben).

**Schritt (2) — Delta v6.13.0 ↔ v6.14.0, gemessen.** Die alte Baseline aus der
Git-Historie gelesen (`git archive 990f1a0e^ .harness/baseline/v6.13.0`, der
Commit `990f1a0e` entfernte sie), die neue aus `HEAD`. `diff -rq` meldet 32
verschiedene Dateien: 26 im `regelwerk/`, 5 Vorlagen, `SHA256SUMS`. Nach
Normalisierung der Versionsstrings (`sed 's/v6\.13\.0/vX/g; s/6\.13\.0/X/g'` bzw.
`6.14.0`, je Kopie) bleiben im `regelwerk/` drei Dateien mit Inhalt: 
`modul-10-review-harness.md` (Kein Stil-Polizist; kein HIGH/MEDIUM ohne
Failure-Szenario; `pfad` als wörtliches Kurzzitat), `modul-15-observability.md`
(Pflicht-Feld ohne Wert der Quelle bleibt Pflicht, „nicht bekannt“ mit Quelle),
`README.md` (Stand-Zeile Kurs-Welle 153 → 156); die normalisierte Diff-Ausgabe
hat 24 Zeilen. Vorlagen: `reviewer.template.md` (LOW mit Konventions-Anker, zwei
neue Grenzen, `pfad` als Kurzzitat; 17 Diff-Zeilen),
`closure-note-reviewer.template.md` (`pfad` als Kurzzitat),
`review-report.template.md` (Spalte Pfad als Kurzzitat), `AGENTS.template.md` und
`conventions.template.md` (nur die Release-URL). **`README.template.md` und
`roadmap.template.md` der Planung sind byte-gleich** (`cmp`, Exit 0). Damit ist
die übernommene Aussage „Vorlage zwischen 6.13.0 und 6.14.0 unverändert“
nachgemessen: die beiden Abweichungen der Planungs-README sind älter als der Bump.

**Liefer-Punkt 1.** `diff <Vorlage> docs/plan/planning/README.md`: am Parent 20
Zeilen, nach der ersten Runde 10 Zeilen, nach Verifikation V-1 13 Zeilen. Übrig
bleiben drei Stellen, alle repo-spezifisch: der ausgefüllte Projektname im Titel
(V-1), der Template-Hinweis-Block (die Vorlage verlangt, ihn beim Kopieren zu löschen)
und der Kommentar `d-check:ignore` am Pfad `docs/plan/carveouts/done/` — das
Verzeichnis existiert nicht (`git ls-files docs/plan/carveouts` nennt nur
`.gitkeep`), ohne den Kommentar meldete `links` ein fehlendes Ziel.
Roadmap: Überschriften gegen `roadmap.template.md` (`diff` über `grep -E '^#'`)
ohne Unterschied, Exit 0. Regelprüfung nach `modul-06-roadmap.md`: keine flache
Welle-Datei (`ls docs/plan/planning/*.md` nennt nur `README.md`), *Offene Wellen*
ohne Zeiger — Bijektion erfüllt; *Abgeschlossene Wellen* gegen die
`welle-*-results.md` unter `done/`: 36 zu 36, `comm -3` leer; einziger
Meilenstein M1 *erreicht* mit Beleg; wellenlose Slices brauchen keine Zeile.
**Abweichung:** der Ruhe-Marker stand, obwohl `in-progress/` diesen Slice trägt
— er steht unverändert seit `d38cf9e7` (2026-09-23, `git blame`), über alle
seither beanspruchten Slices hinweg; die übernommene Aussage „nur Inhalt weicht
ab“ trifft damit nicht zu. Behoben: der Marker ist entfernt. **Übergabe an die
Closure:** nach dem `git mv` dieses Slice nach `done/` trägt `in-progress/`
keinen Slice mehr, und der Marker kommt zurück (Wortlaut aus `d79b7ebd`);
Befund 4 unten.

**Liefer-Punkt 2a — Reviewer-Dateien.** Je Änderung des Deltas (Schritt 2):

| Datei | Ergebnis | Beleg |
|---|---|---|
| `.harness/skills/reviewer.md` | entspricht | alle vier Änderungsstellen der Vorlage übernommen (`git show 675246dd`): LOW mit Konventions-Anker, „Kein Stil-Polizist“, „Kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario“, `pfad` als Kurzzitat; Report-Gerüst zeigt auf die vendored Vorlage (`f5b2840a`) |
| `.harness/skills/closure-note-reviewer.md` | entspricht | `pfad` als Kurzzitat übernommen (`675246dd`) |
| `.claude/agents/reviewer.md`, `.claude/agents/verifier.md` | entspricht | Report-Gerüst `.harness/baseline/v6.14.0/templates/docs/reviews/review-report.template.md` (`f5b2840a`); die Regeln aus `modul-10-review-harness.md` tragen sie über den Skill, keine Zeile `Datei:Zeile` mehr (`grep -n 'Datei:Zeile'` über die vier Dateien: 0 Treffer) |
| sechs Reports (`9b360010`) | Gliederung entspricht; **Abweichung mit Beleg** | Überschriften `##`–`####` gegen die Vorlage: die drei Review-Reports ohne Unterschied, die drei Verifikations-Reports mit je ein bis zwei zusätzlichen Abschnitten (`Offen / übernommen`, `Nicht neu gefahren (übernommen)`, `Plan-DoD gegen Belege`, `Offene Punkte`); Tabellenköpfe Findings, Negativbefunde, Summary in allen sechs je 1×; die Finding-Kennungen alt ↔ neu gleich (`git show 9b360010^:` gegen `9b360010:`), bis auf `verifikation-slice-bench-source-impact-absolut.md` (alt ohne Kennungen, neu `V-1`…`V-5`); Felder ohne Quelle stehen als „nicht erhoben“ (`git grep -c -i 'nicht erhoben'`: 9 · 8 · 6 · 6 · 5 · 1). **Kopfzeile mit falscher Herkunft** (Fixrunde F-4): alle sechs tragen nach `9b360010` `**Skill:** … @ 675246dd` (`git show 9b360010 \| grep -E '^[-+]\*\*Skill'`: drei Zeilen ersetzt, drei neu); `675246dd` ist committet 2026-10-05 07:29:03, der früheste Commit je Report liegt zwischen 2026-10-04 20:43:07 und 2026-10-05 07:17:05 (`git log --format=%ad -- <Report> \| tail -1`) — die Zeile nennt einen Skill-Stand, den keiner der Läufe benutzt hat. Die Reports gehören zu Slices in `done/`; das Neusetzen schreibt Records rückwirkend um (Befund 1) |

**Liefer-Punkt 2b — Audit-Span-Schema.** `.claude/hooks/span-emit.sh` erzeugt
keinen Span: es ruft `.harness/state/bin/ai-harness-init span-emit` und endet in
jedem Zweig mit 0; der Träger ist gitignoriert (`git check-ignore -v` →
`.harness/.gitignore:6:state/`), das Repo führt kein Schema
(`git grep` nach `Audit-Span|requirement\.id|slice\.id` außerhalb von Baseline,
Records und diesem Plan: 0 Treffer neben `span-emit`). Ergebnis: **nicht
Gegenstand dieses Repos** — die Änderung in `modul-15-observability.md` trifft den
Träger. Nebenbefund 2: die Kommentare des Wrappers zitieren Kennungen und Dateien
des Träger-Projekts.

**Liefer-Punkt 2c — Bestandsaufnahme** (Überschriften per
`diff <(grep -E '^#{1,4} ' <Vorlage>) <(grep -E '^#{1,4} ' <Datei>)`,
Platzhalter per
`grep -n -E '<…>|<z\. B\.|<zuerst|<dann|<bei Bedarf|<mover>|<messung>|<vorschau>|<make-target>|<Pfad oder URL>|<Datum>|<repo-spezifischer|<NNN>|<Pfade zu'`
über `harness/README.md` und `harness/conventions.md` — das ist das Muster des
Laufs; die frühere Angabe an dieser Stelle war gekürzt, Fixrunde F-2):

| Dokument | Ergebnis | Beleg |
|---|---|---|
| `harness/README.md` ↔ `README.template.md` | **Abweichung mit Beleg** | Überschriften gleich (Exit 0); stehengebliebene Platzhalter: Zeile `` `<make-target>` `` (§Sensors), Zeilen `make <mover>`/`<messung>`/`<vorschau>` (Werkzeuge), §Safety and scope boundaries zweimal `<…>`, §Leseordnung drei Platzhalter |
| `harness/conventions.md` ↔ `conventions.template.md` | **Abweichung mit Beleg** | Überschriften gleich (Exit 0); Platzhalter: §Adoptierte Konventions-Quellen `<Pfad oder URL>` und `<Pfade zu deinen …>`, MR-000 `**Datum:** <Datum>`, Musterzeilen in §Zusatzklassen und §Glossar |
| `docs/plan/adr/README.md` ↔ `adr/README.template.md` | **Abweichung mit Beleg** | Abschnitt `## Konventionen` fehlt (Regeln stehen als Vorspann, ohne die Regel zum Feld `**Schärft:**`, das alle 155 ADR-Dateien tragen: `git grep -l -F '**Schärft:**'`); Spalten `Datum`/`Datei` statt `Bezug` |
| `docs/plan/carveouts/` ↔ `carveouts/README.template.md` | **Abweichung mit Beleg** | kein `README.md` (`git ls-files docs/plan/carveouts`: nur `.gitkeep`); es gibt keinen Carveout |
| Slice-Vorlage ↔ Praxis | entspricht | `##`-Überschriften der beiden Pläne in `open/` gegen `slice.template.md`: 0 Diff-Zeilen; dieser Plan: nur der bedingte Block `### Sub-Area:` fehlt (alle Sub-Areas GF) |
| Welle-Vorlage ↔ Praxis | entspricht | `welle-routing.md` gegen `welle.template.md`: Exit 0; `welle-routing-results.md` gegen `welle-results.template.md`: zwei zusätzliche, repo-spezifische Abschnitte (`Validator-Feststellung (Modul 8)` — Adresse aus `.claude/commands/implement-slice.md` Schritt 23 —, `Offene Punkte (mit Adresse)`) |

Keine der vier Abweichungen stammt aus dem Bump (die betroffenen Vorlagen haben
im Delta keine inhaltliche Änderung; `conventions.template.md` nur die URL, die
`harness/conventions.md` schon trägt). Sie werden hier nicht behoben (§1, letzter
Ausschluss); **Folge-Slice vorgeschlagen:** `slice-abgeleitete-dokumente-vorlagen-nachzug`
(Gegenstand: die vier Zeilen oben mit Abweichung). Die Datei legt der Planner bei
der Closure in `open/` an (§6, vierter Punkt), mit diesem Gegenstand als Text in
ihrem §1.

**Liefer-Punkt 3 — Entscheidung: Verfahrensregel.** Gemessen in diesem Lauf: (a)
das normalisierte Delta zweier Baseline-Stände ist klein (24 Zeilen Regelwerk,
fünf Vorlagen mit 4 bis 17 Diff-Zeilen) und von einer Rolle in einem Durchgang
lesbar; (b) die Abweichungen dieses Slice fängt das Delta allein **nicht** — die
Planungs-README-Vorlage ist byte-gleich, die Abweichung lag im Bestand; das
Regelwerk benennt genau diesen Fall (`modul-02-harness-bootstrap.md`
§Freshness-Audit, Punkt *Eine Stichprobe gegen den Bestand*); (c) die
Bestandsprüfung trägt mit drei Befehlsformen, deren Reichweite verschieden ist:
Überschriften-`diff` und Platzhalter-`grep` fanden die Abweichungen aus 2c, aber
**nicht** die der Planungs-README (am Parent `d79b7ebd`: Überschriften-`diff`
Exit 0, Platzhalter-Muster 0 Treffer); die fand nur der volle `diff` gegen die
Vorlage (Nachweis unter *Fixrunde*, F-1). Ein voller `diff` meldet bei einem
inhaltlich gefüllten Dokument dagegen vor allem beabsichtigten Inhalt
(`diff README.template.md harness/README.md | wc -l`: 175 Zeilen bei 265 Zeilen
der Datei); deshalb schreibt die Regel ihn nur für Dokumente vor, deren Vorlage im
Delta steht, und immer für die Planungs-README, deren Zeilen feste Klauseln der
Vorlage sind (Fixrunde F-1). Ein Werkzeug brächte über die drei Befehle hinaus
nur eine Liste „abgeleitetes Dokument ↔ Vorlage“ mit eigenem Pflegeaufwand; ein
Gate darauf prüfte Form statt Wahrheit. Festgeschrieben in `harness/targets/pin-stale.md`,
Unterabschnitt *Bump-Ablauf: Vergleich vor dem Löschen der alten Baseline* (drei
Schritte: Delta, Stichprobe gegen den Bestand, Ergebnis je Dokument), mit Verweis
aus `harness/sensors/baseline-verify.md` §Grenze Punkt 2. Die Index-Zeile in
`harness/README.md` bleibt unverändert (ihr Kurzsatz gilt weiter). Kein Gate,
keine ADR. Der Unterabschnitt berührt die Zitat-Korrektur nicht
(`BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` gelesen): Schritt 3
schreibt Records und `Accepted`-ADRs nicht um. Überschneidung: `pin-stale.md` und
`baseline-verify.md` stehen auch in §3 von `slice-harness-targets-inhalt-bereinigen`;
der neue Unterabschnitt ist kein Abschnitt `## Fassung im Gate-Index` und keiner
der dort genannten Befunde.

**Befunde dieses Laufs** (gemeldet, nicht mitgeändert):

1. **Records rückwirkend umgeschrieben.** `9b360010` setzt sechs Reports zu
   Slices in `done/` aus der neuen Vorlage neu; `modul-02-harness-bootstrap.md`
   §Freshness-Audit verlangt für wiederkehrende Vorlagen (Review-Report): „Neue
   Instanzen folgen der neuen Form, bestehende werden nicht rückwirkend
   umgeschrieben.“ Finding-Kennungen erhalten (bis auf die neu vergebenen
   `V-1`…`V-5` eines Verifikations-Reports), 35 Felder „nicht erhoben“, und die
   Kopfzeile `**Skill:** … @ 675246dd` nennt einen Skill-Stand, der nach allen
   sechs Läufen entstand (2a, F-4) — die Lauf-Herkunft ist damit verfälscht.
   **Die Entscheidung ist beim Auftraggeber offen** (stehen lassen mit benannter
   Ausnahme oder `git revert 9b360010`); dieser Slice revertiert nichts.
2. **Fremde Kennungen im Wrapper.** `.claude/hooks/span-emit.sh` zitiert
   `LH-FA-10`, `ADR-0022 Festlegung 5`, `LH-QA-01`, `ADR-0011 Festlegung 6`,
   `test/span-emit-wrapper.bats` und `harness/tools/full-smoke.sh` — Kennungen
   und Dateien des Träger-Projekts; in diesem Repo ist `ADR-0022` der
   Filesystem-Spool, `ADR-0011` Persist-before-ACK, beide Dateien fehlen
   (`ls`). Die Datei ist emittiert (`53420f6a`), Gegenstand des Träger-Projekts.
3. **Bestandsabweichungen** der vier Dokumente aus 2c → vorgeschlagener
   Folge-Slice oben.
4. **Ruhe-Marker ohne Träger.** Keine Datei unter `.claude/commands/` und kein
   Sensor nennt das Setzen oder Entfernen des Markers (`git grep -i
   'Nichts in Arbeit|Ruhe-Marker'` außerhalb der Records: nur der
   Bedienhinweis der Roadmap); bis `f4112fe9` geschah es je Slice per Commit,
   seit `d38cf9e7` nicht mehr. Kandidat für das Beobachtungs-Register bei der
   Closure.
5. **Nummern-Platzhalter `welle-<NN>`/`slice-<NNN>` in fremden Trägern**
   (gemessen am Parent `d79b7ebd`, Ausnahmen wie im Suchlauf; Fixrunde F-5):
   - `welle-<NN>`: 5 Zeilen (Zeile 1 des Suchlaufs) — `.claude/commands/close-welle.md`
     69 und 87, `.claude/commands/implement-slice.md` 352,
     `.harness/skills/reviewer.md` 62, `AGENTS.md` 299.
   - `slice-<NNN>` (`git grep -n -F 'slice-<NNN>'`): 13 Zeilen in 9 Dateien —
     `.claude/commands/close-welle.md` (Zeile 87, dieselbe Zeile wie oben),
     `.claude/commands/implement-slice.md`, `.d-check.yml`,
     `.harness/skills/reviewer.md`, `harness/conventions.md`,
     `harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md` und drei
     `Accepted`-ADRs (`0083`, `0084`, `0086`, unberührbar).
   - Form ohne spitze Klammern `welle-NN`/`slice-NNN`
     (`git grep -n -E 'welle-NN|slice-NNN'`): 10 Zeilen, darunter
     `.d-check.yml` 242, `harness/sensors/docs-check.md` 72,
     `test/integration/integration_test.go` 1360,
     `.claude/commands/implement-slice.md` 242/359/361, `.harness/skills/reviewer.md`
     65 und zwei `Accepted`-ADRs (`0094`, `0099`).

   Nicht jede Stelle ist eine Anker-Form, die `MR-002` überholt:
   `.harness/skills/reviewer.md` 62 beschreibt das Chronik-Muster (Slice- oder
   Wellen-Nummer als Begründung im Produktionscode) und trifft auch die
   nummerierten Kennungen im Bestandsschutz `slice-001`–`slice-105`;
   `MR-002` selbst nennt `slice-<NNN>` als frühere Vergabe. Ob eine Stelle
   nachgezogen wird, ist je Stelle zu urteilen. Gemeldet, nicht geändert; Frist:
   Closure dieses Slice, der Planner zieht nach oder benennt die Träger mit
   Adresse (`AGENTS.md` §3.13).

### Fixrunde (Review `review-slice-baseline-6-14-0-dokumente-nachziehen`, `d5f2c858`)

Nachweise am Parent `d79b7ebd` auf Kopien im Scratchpad
(`git show d79b7ebd:<Datei> > <Kopie>`).

- **F-1 (HIGH) — Regel findet die Abweichung der Planungs-README nicht.** Die
  Regel in `harness/targets/pin-stale.md` (Schritt 2) verlangt jetzt den vollen
  `diff` gegen die versions-normalisierte Vorlage, jede Abweichung beurteilt, für
  jedes Dokument, dessen Vorlage im Delta steht oder von einer Regel des Deltas
  berührt ist, und **immer** für `docs/plan/planning/README.md`; die Rotation ist
  gestrichen. Die Begründung in *Liefer-Punkt 3* ist berichtigt (Überschriften und
  Platzhalter fanden die Planungs-README nicht). **Nachweis:**
  `sed 's/v6\.14\.0/vX/g'` auf Vorlage und Kopie, dann
  `diff <Vorlage> <Kopie>`: Exit 1, 20 Zeilen, 4 Abschnitte (`3,7d2` Hinweis-Block,
  `24c19` Zeile `done/`, `37c32` und `46,47c41,42` je `welle-<NN>-results`);
  `grep -c -F 'Gegenstand an einen anderen Slice übergegangen'` auf die
  Diff-Ausgabe: 1, `grep -c -F 'welle-<NN>-results'`: 2 — beide Abweichungen
  gefunden. Gegenprobe der zwei anderen Formen an derselben Kopie:
  Überschriften-`diff` Exit 0, Muster aus F-2 0 Treffer.
- **F-2 (HIGH) — Platzhalter-Muster der Regel enger als gemessen.** Die Regel
  trägt jetzt genau das Muster des Laufs (`<…>|<z\. B\.|<zuerst|<dann|<bei
  Bedarf|<mover>|<messung>|<vorschau>|<make-target>|<Pfad oder URL>|<Datum>|<repo-spezifischer|<NNN>|<Pfade
  zu`) und zusätzlich die Platzhalter der Vorlage selbst
  (`grep -n -F -f <(grep -o -E '<[^<>]+>' <Vorlage> | sort -u)`); 2c nennt das
  Muster jetzt vollständig. **Nachweis** am Parent: das Muster trifft
  `harness/README.md` 14 Zeilen und `harness/conventions.md` 11 Zeilen, darunter
  die vier Stellen `conventions.md` 56 (`<Pfad oder URL>`), 86 (`<Datum>`),
  `harness/README.md` 125 (`` `<make-target>` ``) und 263 (`<zuerst — …>`). Die
  Vorlagen-Form allein trifft 125 nicht (die Vorlage trägt dort `make fullbuild`),
  deshalb stehen beide Formen in der Regel.
- **F-3 (LOW) — Ruhe-Marker als Closure-Pflicht.** Eigener Closure-Punkt in §2
  und Risiko in §6.
- **F-4 (LOW) — Kopfzeile der Reports.** In 2a als Abweichung mit Beleg
  aufgenommen; Befund 1 nennt sie und hält die Entscheidung beim Auftraggeber
  offen. Nichts revertiert.
- **F-5 (LOW) — Befund 5 unvollständig.** Ergänzt um `slice-<NNN>` (13 Zeilen,
  9 Dateien) und die Form `welle-NN`/`slice-NNN` (10 Zeilen); die Aussage zu
  `.harness/skills/reviewer.md` 62 ist berichtigt. Nur gemeldet.
- **F-6 (INFO) — Grenze der Stichprobe.** Die Regel nennt sie in einem Satz:
  ohne Bump trägt die Stichprobe weiter der Drift-Audit nach `AGENTS.md` §1.
- **Verifikation V-1 (MEDIUM) — Titel-Platzhalter der Planungs-README.** Zeile 1
  trug `# Planning — <Projektname>`; ersetzt durch `# Planning — pg-change-feed`.
  Ergänzung zum F-1-Nachweis: die zweite Platzhalter-Form der Regel
  (`grep -n -F -f <(grep -o -E '<[^<>]+>' README.template.md | sort -u)`, vier
  Platzhalter der Vorlage: `<Kennung>`, `<Platzhalter>`, `<Projektname>`,
  `<welle-id>`) trifft am Parent `d79b7ebd` 2 Zeilen, darunter Zeile 1
  `<Projektname>` — die Regel hätte V-1 gefunden; Liefer-Punkt 1 hatte die Zeile
  übersehen. Am Endstand: zweite Form 3 Zeilen (29 `<welle-id>`, 32 und 41
  `welle-<Kennung>`), alle drei wortgleich mit der Vorlage — die Vorlage trägt sie
  als Form der Konvention, kein auszufüllender Platzhalter (der volle `diff`
  zeigt an diesen Zeilen keinen Unterschied); erste Form (Muster des Laufs)
  0 Treffer. Kein weiterer Platzhalter.

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
- **Der Ruhe-Marker bleibt nach der Closure weg.** Dieser Slice hat ihn entfernt,
  weil `in-progress/` ihn trägt; kein Command und kein Sensor setzt ihn zurück
  (Befund 4), und genau so stand er seit `d38cf9e7` falsch. Gegenmaßnahme:
  eigener Closure-Punkt in §2 (Review F-3). — **Ausgang:** offen bis zur Closure
  (entfallen: Marker nach dem `git mv` gesetzt, Suchlauf-Zeile nachgezogen).

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
