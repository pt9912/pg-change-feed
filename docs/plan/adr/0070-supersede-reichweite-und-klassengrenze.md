# ADR-0070: Supersede-Reichweite und Klassengrenze präzisiert — Supersedes ADR-0069 (nur drei Klauseln)

**Status:** Accepted — Supersedes [`ADR-0069`](0069-commit-msg-hook-einseitige-zusage.md)
in genau **drei** Klauseln: dem Supersede-Satz in deren §Status, der
Grenz-Klausel in deren §Entscheidung Punkt 2 und deren §Re-Evaluierungs-Trigger
(a). Alles Übrige aus `ADR-0069` — die einseitige Zusage selbst
(§Entscheidung Punkt 1, Sicherheits-Zusage), die benannte Sicherheits-Kante,
„Die Hook-Logik bleibt unverändert", die Alternativen, die Folgepflichten —
bleibt unverändert bestehen und wird hier nicht wiederholt.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle, Trigger-Audit der Slice-Closure von
`slice-073`, Baseline-Regelwerk `modul-08-agentenrollen.md` §Konflikt-Pfad <!-- d-check:status-provenance -->
als Rollen-Sequenz und §Rollen-Sequenz für eine Welle, ADR-Zweig)

**Bezug:** [`ADR-0069`](0069-commit-msg-hook-einseitige-zusage.md) (in drei
Klauseln korrigiert), [`ADR-0062`](0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(Punkt 3 bleibt für zwei Aussagen in Kraft), [`ADR-0045`](0045-commit-traceability-standing-gate.md)
(bleibt bindend), `docs/reviews/review-slice-073-fixrunde.md` (F-5, F-7), <!-- d-check:status-provenance -->
`docs/reviews/verify-slice-073.md` (§7), `.githooks/commit-msg`, <!-- d-check:status-provenance -->
`.d-check.yml` (`commits`-Abschnitt), `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md` <!-- d-check:status-provenance -->
(Umsetzungs-Slice).

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0045`, `ADR-0062`
und `ADR-0069`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0069`](0069-commit-msg-hook-einseitige-zusage.md) (Accepted,
2026-09-15) stellt die Zusage des lokalen `commit-msg`-Hooks auf die
einseitige Form um und erklärt im Kopf „Supersedes `ADR-0062` (nur deren
Entscheidung **Punkt 3**, „der Hook … spiegelt exakt die zwei bestehenden
Regeln")". An zwei Stellen trägt dieser Text die Aussage nicht, die er tragen
soll. Beide sind im Bestätigungslauf der Fixrunde gemessen
(`docs/reviews/review-slice-073-fixrunde.md` F-5, F-7) und in der <!-- d-check:status-provenance -->
Verifikation unabhängig reproduziert (`docs/reviews/verify-slice-073.md` §7). <!-- d-check:status-provenance -->

1. **Die Reichweite des Teil-Supersedes ist zweideutig.** `ADR-0062` Punkt 3
   trägt drei Aussagen: (i) „der Hook ist bash-only (kein Docker-Aufruf,
   keine Latenz-Regression für `git commit`)"; (ii) „spiegelt **exakt** die
   zwei bestehenden Regeln"; (iii) die **Merge-/Revert-Ausnahme** der
   positiven Hälfte (`exempt-pattern: '^(Merge |Revert )'`). `ADR-0069`s
   Kopf-Satz nennt zur Bestimmung des abgelösten Punktes allein die
   Spiegelungs-Aussage (ii), schreibt aber „nur deren Entscheidung **Punkt
   3**" — die Apposition ist zweideutig: meint sie den ganzen Punkt oder
   dessen Aussage (ii)? Keine der Aussagen (i)/(iii) wird in `ADR-0069`
   §Entscheidung erneut ausgesprochen; sie sind aus Punkt 1 („keine
   Hook-Rückweisung, die das Gate nicht auch wirft") und aus „Die
   Hook-Logik bleibt unverändert" **ableitbar**, nicht ausgesprochen. Der
   Slice-Plan zitiert beide weiterhin über `ADR-0062` Entscheidung Punkt 3
   (§1 `.md:81` bash-only, DoD-Zeile 4 `.md:116` Merge-/Revert-Ausnahme).
   Bei Lesart „Punkt 3 ganz abgelöst" stünden diese zwei normativen Aussagen
   ohne ausgesprochenen `Accepted`-Träger und die zwei Plan-Zeiger wären
   veraltet.

2. **Die Klassengrenze liest sich als erschöpfend, ist es nicht.**
   `ADR-0069` §Entscheidung Punkt 2 sagt „Die drei oben benannten
   Divergenz-Klassen (a), (b), (c) bleiben offen und **sind die
   festgeschriebene Grenze**"; §Re-Evaluierungs-Trigger (a) spricht von
   „einer der **drei** benannten Klassen". Eine vierte, ebenfalls **laxere**
   Divergenz-Klasse ist erreichbar: ein Commit, dessen rohe Message mit
   einer **Leerzeile** beginnt, auf die `Merge branch 'vp'` folgt, unter
   `--cleanup=verbatim` — der Hook lässt ihn durch (Exit 0, er überspringt
   die Leerzeile und greift die Ausnahme), die gespeicherte Message behält
   die Leerzeile, das Modul liest als Betreff `""`, nimmt die Ausnahme nicht
   an und meldet `commit-untraceable` (Exit 1). Die Klasse ist am gepinnten
   Digest gemessen (dieser Zug, §Fitness Function); unter dem
   Default-Cleanup strippt `git` die Leerzeile, `%s` ist
   `Merge branch 'vp'`, beide Seiten sind grün — sie ist allein mit
   ausdrücklichem `--cleanup=verbatim` (oder `commit.cleanup=verbatim`)
   erreichbar. Der Hook folgt in dieser Klasse git's eigener
   `%s`-Normalisierung; das Modul ist der Ausreißer.

Der Bestätigungslauf hat beide Stellen als `MEDIUM` (F-5) bzw. `INFO` (F-7)
an den Träger der Entscheidung gegeben; ein Accepted-ADR-Text wird nur durch
eine Folge-ADR geändert (`AGENTS.md` §3.5). Dieser Zug ist der Architect-Zug
des Trigger-Audits.

## Entscheidung

Wir korrigieren zwei Klauseln des `ADR-0069`-Textes; die einseitige Zusage
selbst bleibt unverändert.

1. **Reichweite des Teil-Supersedes — eng.** `ADR-0069` supersedes `ADR-0062`
   **nur in der Spiegelungs-Aussage (ii)** („spiegelt exakt die zwei
   bestehenden Regeln"). `ADR-0062` **Punkt 3 bleibt für zwei Aussagen in
   Kraft**: der Hook ist und bleibt **bash-only (kein Docker-Aufruf)**, und
   er trägt die **Merge-/Revert-Ausnahme** der d-check-Positiv-Hälfte. Der
   Slice-Plan bleibt mit seinen zwei Zeigern auf `ADR-0062` Entscheidung
   Punkt 3 **korrekt**; sie sind nicht nachzuziehen.

2. **Klassengrenze — benannt, nicht erschöpfend.** Die drei Klassen (a), (b),
   (c) sind die **gemessene, reproduzierte** Grenze; sie sind **nicht** als
   erschöpfend bewiesen. Die vierte Klasse — führende Leerzeile vor einem
   `Merge …`-Betreff unter `--cleanup=verbatim` — wird als weitere benannte,
   laxere Klasse geführt. `ADR-0069` §Re-Evaluierungs-Trigger (a) wird auf
   „eine der **benannten** Klassen wird geschlossen **oder eine weitere
   Divergenz-Klasse wird gemessen**" erweitert.

Die **Sicherheits-Zusage** (`ADR-0069` §Entscheidung Punkt 1) bleibt
unverändert. Ebenso bleibt die dort benannte Sicherheits-Kante (leerer
Vor-scissors-Text unter `--cleanup=scissors`) stehen: sie ist ausdrücklich
„benannt statt behauptet"; drei nicht erreichte Konstruktionen sind **kein**
Unmöglichkeitsbeweis und tragen keine Streichung.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0069`s Supersede-Satz und Klassengrenze bleiben stehen | keine Folge-ADR, kein Doku-Nachzug | der Supersede-Satz bleibt zweideutig — bei Lesart „Punkt 3 ganz abgelöst" hätten (i)/(iii) keinen ausgesprochenen `Accepted`-Träger und zwei Plan-Zeiger zeigten auf einen abgelösten Punkt; die Grenz-Klausel bliebe eine erschöpfende Behauptung über eine Liste, die messbar unvollständig ist (`F-7`) — genau der Fehler, den `ADR-0069` an „spiegelt exakt" korrigiert |
| B — Punkt 3 vollständig abgelöst belassen und die Aussagen (i)/(iii) in dieser ADR **neu aussprechen** (weite Lesart) | ein einziger Träger für die ganze Hook-Politik | dupliziert zwei unveränderte Aussagen in eine neue ADR; macht die zwei Plan-Zeiger veraltet und vergrößert den Eingriff — stabile Aussagen wandern ohne Grund aus `ADR-0062` heraus |
| C — nur die Klassengrenze korrigieren, die Supersede-Reichweite unverändert lassen | kleinerer Eingriff | die Zweideutigkeit des Supersede-Satzes bleibt; ob (i)/(iii) von `ADR-0062` oder von `ADR-0069` getragen sind, bliebe offen — der `F-5`-Anlass bleibt ungelöst |
| **D — Reichweite auf die Spiegelungs-Aussage verengen (eng) UND die Klassengrenze als gemessen-nicht-erschöpfend fassen — eine Folge-ADR für beide Klauseln (gewählt)** | beide normativen Aussagen (i)/(iii) behalten einen ausgesprochenen `Accepted`-Träger (`ADR-0062` Punkt 3); die zwei Plan-Zeiger bleiben gültig; die Klassengrenze wird wieder wahr und messbar (vier benannte Klassen statt einer erschöpfenden Drei-Liste); eine Folge-ADR statt zweier | zwei ADRs (`0062`+`0069`+`0070`) müssen zusammengelesen werden; der Hook bleibt in vier statt drei Klassen blind — die Ursprungsbeobachtung kann dort erneut auftreten |

## Konsequenzen

- Positiv: Die zwei normativen Aussagen (i) bash-only und (iii)
  Merge-/Revert-Ausnahme haben einen ausgesprochenen `Accepted`-Träger
  (`ADR-0062` Punkt 3); die zwei Plan-Zeiger (`.md:81`, `.md:116`) bleiben
  gültig.
- Positiv: Die Klassengrenze ist wieder wahr — benannt und am gepinnten
  Digest reproduzierbar (vier Klassen), nicht als erschöpfend behauptet.
- Positiv: Keine Schwelle gesenkt, kein Verhaltens-Code geändert
  (`AGENTS.md` §3.6 greift nicht — der Hook ist kein Gate); die Hook-Logik
  bleibt unverändert (`ADR-0069`).
- Negativ mit Grenze: Der Hook bleibt in den **vier** Klassen blind; die
  Ursprungsbeobachtung (`BEO-PGC/commit-traceability-kein-vorab-hook`) kann
  in genau diesen Klassen erneut auftreten — das Gate fängt sie beim
  nächsten `make gates`, ggf. erst nach einem Push.
- Folgepflicht (Implementer-Zug): der Slice-Plan nennt die drei Klassen als
  Grenze in §Bezug (`.md:20`), §1 (`.md:55-60`) und §3 (`.md:151`) — die
  Formulierung wird auf die **benannte, nicht als erschöpfend bewiesene**
  Grenze nachgezogen (die vierte Klasse nennen oder „nicht erschöpfend"
  aussprechen). Hook-Kommentar und `harness/README.md`-Zeile nennen keine
  geschlossene Drei-Liste („vollstaendig ist er nicht" bzw. „fängt **nicht**
  jeden Verstoß") und bleiben unverändert.
- Hinweis: `F-5` braucht **keinen** Plan-Nachzug — Punkt 3 bleibt in Kraft.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `d-check` Modul `commits` via `--range`/`--commit-msg` (Digest aus `d-check.mk`) | vierte Klasse (führende Leerzeile + `Merge …`, `--cleanup=verbatim`) ⇒ Exit ≠ 0 (`commit-untraceable`); unter Default-Cleanup ⇒ Exit 0 — die Grenze ist reproduzierbar | kein Gate — Messung |
| `.githooks/commit-msg` | dieselbe rohe Message ⇒ Exit 0 (laxere Richtung) | kein Gate — lokaler Hook, außerhalb von `make gates` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** eine der benannten Klassen wird im Hook
geschlossen **oder es wird eine weitere Divergenz-Klasse gemessen** — dann
die Liste neu messen und per Folge-ADR nachziehen; **(b)** der Hook weist
einen Commit zurück, den das Standing-Gate zulässt (Bruch der
Sicherheits-Zusage — die einzige Richtung, die `ADR-0069` nicht toleriert) —
dann unverzüglich Folge-ADR bzw. Hook-Korrektur; **(c)** `ADR-0045`
§Re-Evaluierungs-Trigger (a) tritt ein (das d-check-Modul `commits` erhält
eine verbotene-Muster-Option) — dann sind die Klassen neu zu messen.
Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: Trigger-Audit der `slice-073`-Closure (Modul 6 §Wellen-Closure-Prozedur Schritt 2, ADR-Zweig); korrigiert `ADR-0069`s Supersede-Reichweite und Klassengrenze | [`docs/reviews/review-slice-073-fixrunde.md`](../../reviews/review-slice-073-fixrunde.md) F-5/F-7, [`docs/reviews/verify-slice-073.md`](../../reviews/verify-slice-073.md) §7 | <!-- d-check:status-provenance -->

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0070` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
