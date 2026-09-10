# ADR-0045: Commit-Traceability als Standing-Gate

**Status:** Accepted

**Datum:** 2026-09-10

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-10)

**Bezug:** — (Prozess-ADR ohne Spec-Stratum; die getragene Regel steht in
[`AGENTS.md`](../../../AGENTS.md) §5 und in
[`harness/README.md`](../../../harness/README.md) §Traceability rules)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Traceability-Regel für Commit-Messagen — je Message mindestens eine
`LH-*`-/`ADR-*`-Kennung, nie eine Struktur-ID (`SPEC-…`/`ARC-…`) im
Betreff ([`AGENTS.md`](../../../AGENTS.md) §5,
[`harness/README.md`](../../../harness/README.md)
§Traceability rules) — ist als Regel verkörpert, aber von keinem Sensor
getragen. Die Verstoß-Klasse hat die Steering-Loop-Schwelle erreicht: die
Review-Berichte unter [`docs/reviews/`](../../reviews/) tragen drei
Auftreten der direkten Form (Commit ohne Vertrags-Kennung; F-1 im Bericht
vom 2026-09-09, F-2 im Bericht vom 2026-09-10 mit den Commits `3ea9244`
und `779bc64`) und ein weiteres Auftreten der Struktur-ID-Form
(Planner-Commits mit `SPEC-…` im Betreff, F-5 im Bericht vom 2026-09-09).
Die Regel-Form trägt seit ihrer ersten Schärfung den Anker in
`.claude/commands/implement-slice.md` und wirkte dennoch nur als
Review-Prüfpflicht — der Review sieht nur den jeweils geprüften Range,
und die Messagen bleiben historisch: Korrektur wirkt nur vorwärts.

Fähigkeits-Prüfung des d-check-Moduls `commits` (Handbuch-Kapitel
„Traceability-Kennungen in Commit-Messages prüfen", Stand des gepinnten
Images v0.74.1): **(a)** ein Commit ohne gültige Kennung wird als Befund
`commit-untraceable` gemeldet — ja; **(b)** verbotene Muster im Betreff —
nein, das Modul kennt nur positives `id-patterns` und das
`exempt-pattern` für Merge-/Revert-Betreffe, keine
verbotene-Muster-Option.

## Entscheidung

Wir wählen **ein Standing-Gate `commit-traceability` im
`make gates`-Bündel über die letzten fünf Commits** (`HEAD~5..HEAD`,
per `RANGE=base..head` überschreibbar), getragen von zwei Werkzeugen mit
geteilter Hälfte: die positive Hälfte (≥ 1 `LH-*`/`ADR-*` je Message)
prüft das d-check-Modul `commits` (Befund `commit-untraceable`,
konfiguriert im `commits`-Abschnitt der `.d-check.yml`); die Grenz-Hälfte
(keine Struktur-ID im **Betreff**) prüft
`tools/harness/commit-traceability.sh`. Beide Hälften laufen im selben
Target über dieselbe Range — eine Regel, je Hälfte ein Träger; ein
zweiter Sensor über dieselbe Hälfte wäre eine zweite Quelle für
denselben Zustand. **Kein commit-msg-Hook** — der Hook fängt nur die
Pending-Message des Committers, aber die Regel gilt für Commits jedes
Urhebers (Planner-, Zweitschreiber-, Allowlist-Commits), und der
Gate-Nachweis (record-gates) sieht einen Hook nicht.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; Regel bleibt Review-Prüfpflicht | kein Werkzeug-Aufwand | genau der belegte Zustand: die Klasse trat trotz Regel-Text und Skill-Zeile wiederholt auf; der Review deckt nur den jeweils geprüften Range |
| B — commit-msg-Hook (`d-check --commit-msg`) | fängt die Pending-Message vor dem Commit; keine Fenster-Problematik | Hook-Installation ist client-seitig nicht erzwingbar (kein Wächter); Commits außerhalb des Harness bleiben ungedeckt; der Gate-Nachweis-Stempel sieht den Hook nicht |
| C — CI-only Range-Prüfung (`doc-commits` je PR) | deckt jeden PR-Range vollständig | der Ort der erzwingenden Nachweise in diesem Repo ist das Gate-Bündel (Stop-Hook auf record-gates); ein nur-CI-Sensor läuft im lokalen Lauf nie und trägt die Regel nicht pro Commit |
| **D — Standing-Gate im Gate-Bündel (gewählt)** | mechanisch bei jedem `make gates`; deckt Commits jedes Urhebers im Fenster; d-check-Hälfte hermetisch und digest-gepinnt; die Betreff-Hälfte trägt ein kleiner Shell-Sensor, dessen einziges Werkzeug `git log` ist | Altbestand außerhalb des 5-Commit-Fensters bleibt unentdeckt (Beweis-Grenze, unten benannt); ein flacher Klon scheitert fail-closed an `HEAD~5` (Exit 2) |

## Konsequenzen

- Positiv: Jede künftige kennungslose Message und jeder Struktur-ID-Betreff
  im Fenster färbt `make gates` rot — die Klasse ist ab dem Inkrafttreten
  mechanisch getragen, und der Gate-Nachweis-Stempel schließt den Commit
  als Umweg aus.
- Positiv: Die Teilung der Hälften folgt der Fähigkeits-Prüfung: d-check
  trägt, was sein Modul ausdrückt; der Shell-Sensor trägt nur die Grenze,
  die das Modul heute nicht ausdrückt. Zwei Sensoren über dieselbe Hälfte
  wären zwei Quellen für dieselbe Regel und würden driften.
- Negativ mit Grenze: Das Fenster sind die letzten fünf Commits —
  Altbestand außerhalb des Fensters (hier: die bereits belegten Altlasten
  vor dem Inkrafttreten) bleibt historisch und wird nicht nachträglich
  rot; die Beweis-Kraft des Gates gilt nur vorwärts. Ein Review-Auffund
  von Verstößen außerhalb des Fensters ist der Trigger zur Fenster-Prüfung
  (unten), nicht zur stillen Verdopplung des Sensors.
- Negativ: `HEAD~5` setzt eine vollständige Historie voraus; auf einem
  flachen Klon bricht das Gate fail-closed mit Exit 2 — sichtbar, nicht
  still. Das `exempt-pattern` nimmt Merge-/Revert-Betreffe aus.
- Folgepflicht: `GATE_CHECKS += commit-traceability` in
  `harness/mk/doc-gate.mk`; Sensors-Zeile in `harness/README.md` §Sensors
  und Target-Zeile in `AGENTS.md` §4; die regeltragende Stelle in
  `harness/README.md` §Traceability rules trägt den Herkunfts-Anker der
  Verkörperung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| d-check Modul `commits` (DC-FA-COMMITS-001) + `tools/harness/commit-traceability.sh` | je Message der Range ≥ 1 `LH-*`-/`ADR-*`-Kennung (Befund `commit-untraceable`); kein Struktur-ID-Token im Betreff | `make commit-traceability` (im Gate-Bündel, Standing-Range `HEAD~5..HEAD`) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger: **(a)** das d-check-Handbuch führt eine
verbotene-Muster-Option für das `commits`-Modul ein (Changelog-Eintrag im
Handbuch-Kapitel zum Modul) — dann trägt das Modul beide Hälften, und der
Shell-Sensor wird durch die `.d-check.yml`-Konfiguration ersetzt
(Folge-Commit; Folge-ADR nur, wenn auch die Fenster- oder Bündel-Semantik
ändert); **(b)** ein Review-Auffund von Verstößen, die das 5-Commit-Fenster
bereits verlassen haben, ohne dass ein Sensor sie sah — dann Fenster-Größe
oder Träger-Kombination als Folge-ADR prüfen. Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-10 | Accepted — Anlass: drittes Auftreten der Traceability-Klasse (Review-Berichte `docs/reviews/`, F-1/F-2); Verkörperung durch den Architect-Lauf | [`docs/reviews/`](../../reviews/) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).