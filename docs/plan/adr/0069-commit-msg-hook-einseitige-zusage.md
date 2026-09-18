# ADR-0069: `commit-msg`-Hook — einseitige Zusage statt exakter Spiegelung — Supersedes ADR-0062 (nur Punkt 3)

**Status:** Accepted — Supersedes [`ADR-0062`](0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(nur deren Entscheidung **Punkt 3**, „der Hook … spiegelt exakt die zwei
bestehenden Regeln"; die Punkte 1, 2 und 4 — das Standing-Gate bleibt die
durchsetzende Instanz, die Aktivierung bleibt lokaler Opt-in, `--no-verify`
bleibt ein gültiger Umgehungsweg — bleiben unverändert bestehen und werden
hier nicht wiederholt)

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug im
Konflikt-Pfad, Baseline-Regelwerk `modul-08-agentenrollen.md`
§Konflikt-Pfad als Rollen-Sequenz, Verdikt 2 — anderer Kontext als der
Reviewer-Lauf, der die beiden HIGH fand)

**Bezug:** [`ADR-0062`](0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(korrigierter Punkt 3), [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(bleibt für alles außer der Hook-Klausel bindend), Review zu `slice-073` <!-- d-check:status-provenance -->
(Übergabe-Artefakt, Findings F-1/F-2) und der Architect-Verdikt dieses Zugs
(commit-msg-Hook — einseitige Zusage), `.githooks/commit-msg` (das Artefakt),
`tools/harness/commit-traceability.sh` und `.d-check.yml` (`commits`-Abschnitt)
(die zwei Gate-Hälften), `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md` <!-- d-check:status-provenance -->
(Umsetzungs-Slice), `harness/README.md` §Traceability rules.

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0045` und `ADR-0062`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0062` (Accepted, 2026-09-14) macht einen lokalen, nicht-durchsetzenden
`commit-msg`-Hook zulässig und beschreibt ihn in Punkt 3 als Artefakt, das
„**exakt** die zwei bestehenden Regeln … spiegelt" (positiv: ≥ 1
`LH-*`/`ADR-*`-Kennung **im Betreff**; negativ: keine
`SPEC-*`/`ARC-*`-Struktur-ID **im Betreff**), inklusive der Merge-/Revert-Ausnahme.

Der Reviewer-Lauf von `slice-073` <!-- d-check:status-provenance --> hat an diesem Punkt drei
Dinge gemessen, die diese Beschreibung nicht trägt:

1. **Die positive Hälfte des Hooks prüft nicht den Betreff, sondern die
   ganze Message-Datei** (`grep -qE` über `$1`, `.githooks/commit-msg:57-64`).
   Das ist inhaltlich die richtige Weite — die positive Hälfte des Gating
   (`d-check`-Modul `commits`) liest „je Message", nicht „je Betreff"; eine
   betreff-scoped Prüfung wäre **strenger** als das Gate und würde einen
   Commit mit Kennung nur im Body fälschlich zurückweisen, den das Gate
   zulässt. Punkt 3 ist damit **in sich widersprüchlich**: er verlangt
   zugleich „im Betreff" und die Spiegelung einer Hälfte, die je Message
   liest — beides ist nicht gleichzeitig erfüllbar.
2. **„Spiegelt exakt" ist gemessen unwahr.** Das `d-check`-Modul liest die
   **bereinigte** Message (alles ab der ersten scissors-Zeile `^#.*>8`
   entfällt, jede `#`-Zeile entfällt), der Hook die **rohe** Datei. Drei
   Divergenz-Klassen sind am gepinnten Image (`--commit-msg`) und an je einem
   realen Commit mit aktivem Hook belegt — die Messung dieses Zugs
   reproduziert sie am Digest `sha256:18e9cd85…` aus `d-check.mk`:

   | Klasse | Eingabe (rohe Datei) | Modul/Gate | Hook |
   |---|---|---|---|
   | (a) Kennung nur in einer `#`-Kommentarzeile | `docs: betreff ohne kennung` + `# ADR-0045` | Befund (Exit 1) | lässt durch (Exit 0) |
   | (b) Kennung nur hinter der scissors-Zeile/im Verbose-Diff | `…` + `# ---- >8 ----` + `+ADR-0045` | Befund (Exit 1) | lässt durch (Exit 0) |
   | (c) Struktur-ID auf der Fortsetzungszeile des ersten Absatzes | `docs: betreff (ADR-0045)` + `SPEC-012 fortsetzung` | Shell-Sensor rot (`%s` fügt den Absatz zusammen) | lässt durch (Exit 0) |

   Alle drei laufen in die **laxe** Richtung; ein Fall, in dem der Hook einen
   Commit zurückweist, den das Gate zulässt, wurde weder real noch am Modul
   gefunden.
3. **Der Träger der Abweichung ist eine Plan-Zelle.** Die Begründung, die die
   message-weite positive Hälfte rechtfertigt, steht ausschließlich in §3 des
   Slice-Plans — einem Zeitdokument, das mit der Closure nach `done/` wandert
   und archiviert wird. Der Text einer Accepted-ADR wird durch eine Folge-ADR
   geändert ([`AGENTS.md`](../../../AGENTS.md) §3.5), nicht durch eine
   Plan-Zelle.

Der Reviewer hat beides als HIGH mit Rollen-Widerspruch in den Konflikt-Pfad
gegeben; ein Herabstufen wegen des begründeten Implementer-Widerspruchs ist
ausgeschlossen. Dieser Zug ist der Architect-Schritt der Sequenz.

## Entscheidung

Wir korrigieren `ADR-0062`s Punkt 3 und fassen die Zusage des Hooks als
**einseitige Approximation** — nicht als exakte Spiegelung. Die
Festschreibung hat zwei Hälften, und nur eine davon gilt:

1. **Sicherheits-Zusage (gilt).** Der Hook weist keinen Commit
   zurück, den das Standing-Gate zulässt: jede Hook-Rückweisung ist auch eine
   Gate-Rückweisung. Der Zuschnitt trägt das — beide Hook-Hälften lesen einen
   Text, der den Text ihrer Gate-Gegenhälfte enthält (die positive Hälfte die
   rohe Datei ⊇ die bereinigte Message; die Grenz-Hälfte die erste Zeile ⊆
   `%s`, den ersten Absatz) —, aber er trägt es nicht ausnahmslos: liest der
   Hook unter `--cleanup=scissors` hinter der scissors-Zeile, die `git` später
   verwirft, kann seine Betreff-Zeile Text sein, den das Gate nicht mehr
   liest. Diese Kante ist eng — sie setzt einen **leeren** Vor-scissors-Text
   voraus, den `git` im Default ohnehin abweist —, und ist bislang nicht real
   aufgetreten; sie wird hier **benannt statt behauptet**.
2. **Vollständigkeits-Zusage (gilt NICHT).** Der Hook ist **nicht**
   vollständig: er fängt nicht jede Verletzung, die das Standing-Gate fängt.
   Die drei oben benannten Divergenz-Klassen (a), (b), (c) bleiben offen und
   sind die festgeschriebene Grenze. Das Gate bleibt die durchsetzende
   Instanz (`ADR-0062` Punkt 1, unverändert).

Zugleich wird die **Selbstwidersprüchlichkeit** aus Punkt 3 aufgelöst, und
zwar in der sicheren Richtung: die positive Hälfte ist **message-weit** (die
ganze rohe Datei), die Grenz-Hälfte **betreff-scoped** (die erste Zeile des
ersten Absatzes). Die message-weite positive Hälfte ist die Zusage, die den
Hook nie strenger macht als das Gate; „im Betreff" war zu streng.

Die **Hook-Logik bleibt unverändert.** Insbesondere wird die
Bereinigungs-Semantik des Moduls **nicht** nachgebaut (§Verglichene
Alternativen, Option B).

### Warum die Approximation das Festschreibbare ist

- **Der Zweck, den Punkt 3 selbst nennt, ist einseitig.** Die Begründung der
  Merge-/Revert-Ausnahme lautet: „ohne diese Ausnahme würde der Hook jeden
  Merge-Commit fälschlich zurückweisen, den das bestehende Gate zulässt". Die
  Zusage zielt auf **Fehlrückweisung**, nicht auf Vollständigkeit. Diese eine
  Richtung ist in allen gemessenen Fällen erfüllt; Vollständigkeit war nie die
  Anforderung, nur die überzogene Formulierung „exakt".
- **Vollständigkeit wäre nur durch Nachbau zweier verschiedener Bereinigungen
  erreichbar** — der `#`-/scissors-Bereinigung des Moduls für die positive
  Hälfte und der `git`-`%s`-/Cleanup-Semantik für die Grenz-Hälfte. Letztere
  ist zur Hook-Zeit teils **unbekannt**: der Hook läuft vor
  `git`s Cleanup, und `--cleanup=` bzw. `commit.cleanup` übersteuern den
  Default. „Exakt" bliebe deshalb auch nach einem partiellen Nachzug unwahr;
  die Zusage müsste weiterhin einseitig formuliert werden — nur mit mehr Code.
- **Der Hook ist kein Gate und keine zweite Wahrheitsquelle** (Punkt 1). Eine
  treue Zweitimplementierung einer bereits gate-getragenen Regel invertiert
  die Schichtung und schafft eine zweite Drift-Fläche in einem Artefakt, dessen
  ganze Rechtfertigung „bash-only, keine Latenz" war.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0062` Punkt 3 bleibt stehen | keine neue ADR, kein Doku-Nachzug | der Text bleibt gemessen unwahr („exakt") und in sich widersprüchlich („im Betreff" **und** message-weite Spiegelung); der Träger der Abweichung bliebe eine Plan-Zelle, die mit der Closure archiviert wird — genau der `AGENTS.md` §3.5-Verstoß, den der Reviewer als HIGH führt |
| B — „exakte Spiegelung" ist das Ziel; der Hook wird nachgezogen (bereinigte Message vor der positiven Prüfung, Absatz-`%s`-Semantik für die Grenz-Hälfte) | schließt die drei Klassen; der Hook würde das Gate vollständiger abbilden | baut zwei divergierende Bereinigungs-Semantiken nach, von denen eine (`git`s Cleanup-Modus) zur Hook-Zeit teils unbekannt ist; „exakt" bliebe trotzdem unwahr (siehe oben), die Zusage müsste weiterhin einseitig lauten — mehr Verhaltens-Code in einem Nicht-Gate-Artefakt, ohne den entscheidenden Vorteil |
| C — `ADR-0062` vollständig neu fassen (ganze ADR superseded) | ein einziges, vollständiges Nachfolgedokument statt zweier zusammenzulesender ADRs | unverhältnismäßig — die Punkte 1, 2 und 4 (Gate ist durchsetzend, Opt-in, `--no-verify`) bleiben unverändert richtig; dupliziert stabilen Text, entgegen dem etablierten Muster enger Klausel-Korrekturen (`ADR-0048` → `ADR-0047`, `ADR-0063` → `ADR-0058`) |
| **D — Teil-Supersede: nur Punkt 3 wird auf die einseitige Zusage umgestellt; Hook-Logik unverändert (gewählt)** | die Zusage wird wieder wahr und **messbar** (drei benannte Klassen statt „exakt"); der Träger ist die Entscheidung, nicht eine Plan-Zelle; keine neue Verhaltens-Fläche in einem Nicht-Gate-Artefakt; folgt dem Repo-Muster | zwei ADRs (`0062`+`0069`) müssen zusammengelesen werden; der Hook verpasst die drei Klassen bewusst — die Ursprungsbeobachtung (Vorfall erst nach `git push` sichtbar) kann in genau diesen Klassen erneut auftreten |

## Konsequenzen

- Positiv: Die Zusage des Hooks ist wieder wahr und am gepinnten Image
  nachmessbar; die drei Klassen sind benannt statt von „exakt" überdeckt.
- Positiv: Der Träger der Abweichung ist ab hier die Entscheidung; die
  Plan-Zelle verliert ihre Rolle als Regel-Träger.
- Positiv: Das Standing-Gate bleibt unangetastet die durchsetzende Instanz
  (`ADR-0062` Punkt 1); es wird keine Schwelle gesenkt (`AGENTS.md` §3.6
  greift nicht — der Hook ist kein Gate).
- Negativ mit Grenze: Der Hook bleibt in den Klassen (a), (b), (c) blind; die
  Ursprungsbeobachtung (`BEO-PGC/commit-traceability-kein-vorab-hook`) kann
  in genau diesen Klassen erneut auftreten — das Gate fängt sie beim nächsten
  `make gates`, ggf. erst nach einem Push. Das ist die bewusst
  festgeschriebene Grenze, kein offenes Risiko mehr.
- Negativ: Der Hook läuft vor `git`s Cleanup; seine Sicht ist strukturell die
  rohe Datei. Jede künftige Vollständigkeits-Erwartung an ihn muss gegen die
  drei Klassen geprüft werden.
- Folgepflicht (Implementer-Zug): Slice-Plan §1 (Ziel-Wortweite), §2
  DoD-Zeile 2 („im Betreff" → message-weit) und §3 (Plan-Zelle trägt den
  Rang-Zeiger auf diese ADR statt der Begründung selbst) werden nachgezogen;
  DoD-Zeile 1 führt den Nachweis mit einem **isolierenden** Input
  (Kennung **und** Struktur-ID, z. B. `docs: update SPEC-012 table (ADR-0045)`)
  — nur er unterscheidet die Grenz-Hälfte von der positiven (Review F-3).
- Folgepflicht (Implementer-Zug): Der Einleitungskommentar von
  `.githooks/commit-msg` nennt die einseitige Zusage und trägt den
  Rang-Zeiger auf diese ADR; die Zeile in `harness/README.md`
  §Traceability rules („meldet denselben Verstoß") wird auf die einseitige
  Zusage nachgezogen (Review F-2).
- Folgepflicht (Implementer-Zug): Der Kommentar-Rang-Zeiger in `d-check.mk`
  („Gate statt commit-msg-Hook (`ADR-0045`, Alternativen)") wird auf den
  aktuellen Träger gezogen — die zitierten Alternativen tragen die von
  `ADR-0062` superseded Hook-Aussage (Review F-4).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `commits` via `--commit-msg <rohe Message-Datei>` (Digest aus `d-check.mk`) | positive Hälfte liest die **bereinigte** Message: Kennung nur in einer `#`-Zeile (a) oder hinter der scissors-Zeile (b) ⇒ Befund (Exit ≠ 0), der Hook lässt durch — die Klasse ist reproduzierbar | kein Gate — Messung (der Hook ist kein `make`-Ziel) |
| `tools/harness/commit-traceability.sh` über eine Range (`%s`) | Grenz-Hälfte liest den **ersten Absatz**: Struktur-ID auf der Fortsetzungszeile ⇒ Exit ≠ 0, der Hook lässt durch — Klasse (c) | `make commit-traceability` |
| `.githooks/commit-msg` | **Sicherheits-Zusage:** jede Rückweisung des Hooks ist auch eine Rückweisung des Gates (einseitig, im belegten Zuschnitt; benannte Kante: leerer Vor-scissors-Text unter `--cleanup=scissors`). **Nicht** vollständig: die Klassen (a)/(b)/(c) bleiben offen | kein Gate — lokaler Hook, außerhalb von `make gates`/`record-gates` (`ADR-0062` Punkt 1) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** eine der drei benannten Klassen (a)/(b)/(c) wird
im Hook geschlossen — dann die Liste neu messen und diese ADR per Folge-ADR
nachziehen; **(b)** der Hook weist einen Commit zurück, den das Standing-Gate
zulässt (Bruch der Sicherheits-Zusage — die einzige Richtung, die diese ADR
nicht toleriert) — dann unverzüglich Folge-ADR bzw. Hook-Korrektur;
**(c)** `ADR-0045` §Re-Evaluierungs-Trigger (a) tritt ein (das `d-check`-Modul
`commits` erhält eine verbotene-Muster-Option) — dann kann die Grenz-Hälfte an
das Modul delegiert werden, und die Klassen sind neu zu messen. Andernfalls
permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Review zu `slice-073` F-1/F-2 (zwei HIGH mit Rollen-Widerspruch gegen `ADR-0062` Punkt 3) → Konflikt-Pfad, Verdikt 2; korrigiert `ADR-0062` Punkt 3 auf die einseitige Zusage | Review zu `slice-073`, der Architect-Verdikt dieses Zugs <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | PENDING_COMMIT |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0069` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
