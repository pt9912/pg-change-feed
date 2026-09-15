# Architect-Verdikt: Supersede-Reichweite und Klassengrenze des Hook-Teilsupersedes (Trigger-Audit der `slice-073`-Closure)

**Rolle:** Architect (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Konflikt-Pfad als Rollen-Sequenz und §Rollen-Sequenz für eine Welle).
Träger dieses Zugs ist der **Trigger-Audit der Slice-Closure** von
`slice-073` (Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur
Schritt 2, ADR-Zweig; im wellenlosen Betrieb trägt ihn die Slice-Closure,
Modul 6 §Was der wellenlose Betrieb selbst auslöst).

**Beteiligte Rollen:** Verifier (hat die zwei Reviewer-Befunde unabhängig
reproduziert, `docs/reviews/verify-slice-073.md` §7) · Reviewer (hat F-5,
F-6, F-7 gemeldet, `docs/reviews/review-slice-073-fixrunde.md`) ·
**Architect (dieser Zug, entscheidet)** · Implementer (Folgearbeit,
siehe unten) · Planner (Closure: §7-Notiz, Register, §6-Risiko-Ausgänge,
drei Paarungen). Verifier und Validator waren am Konflikt selbst **nicht**
beteiligt; der Verifier-Beleg reist nur als Eingang mit.

**Übergabe-Artefakt:** `docs/reviews/review-slice-073-fixrunde.md`
(Commit `bac2bbc`) zusammen mit `docs/reviews/verify-slice-073.md`
(Commit `15d5ac1`).

**Rolleninhaber dieses Zugs:** der unabhängige Architect-Lauf, der dieses
Verdikt und die Folge-ADR [`ADR-0070`](../plan/adr/0070-supersede-reichweite-und-klassengrenze.md)
geschrieben hat.

**Datum:** 2026-09-15.

**Bezug:** [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
§Re-Evaluierungs-Trigger (b), §Entscheidung Punkt 3;
[`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
(in drei Klauseln korrigiert); [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
(bleibt bindend); neu: [`ADR-0070`](../plan/adr/0070-supersede-reichweite-und-klassengrenze.md).

---

## (1) `ADR-0062` §Re-Evaluierungs-Trigger (b) — **bestätigt: eingetreten und ausgegangen**

**Verdikt: bestätigt.** Der Trigger ist eingetreten, und sein Ausgang ist
`ADR-0069`; `ADR-0062` bleibt darüber hinaus für ihre Punkte 1, 2 und 4
bindend (`AGENTS.md` §3.5) — sie wird nicht geändert.

Begründung, mit einer Präzisierung:

- **Der Trigger ist eingetreten.** Sein beobachtbarer Kern lautet: „der Hook
  und das Standing-Gate weichen real auseinander". Das ist gemessen — vier
  Divergenz-Klassen (a), (b), (c) und die hier neu bestätigte vierte, alle in
  die **laxe** Richtung (Hook lässt durch, das Gate verwirft), belegt
  `docs/reviews/verify-slice-073.md` §6 und F-7.
- **Der Ausgang trägt.** Die vom Trigger verlangte Folge ist „Vereinheitlichung
  oder Rückbau des Hooks als Folge-ADR prüfen". Genau das ist `ADR-0069`: eine
  Accepted-Folge-ADR, die beide genannten Wege als Optionen ausweist und
  verwirft (Option B „exakte Spiegelung, Hook nachziehen", Option A „nichts
  tun / Rückbau") und eine dritte, dokumentierte Auflösung wählt („die
  Approximation ist das Festschreibbare", einseitige Zusage). Der Trigger
  verlangt die **Prüfung**, nicht einen der beiden Wege; die Prüfung ist mit
  dem Folge-ADR-Zug erfolgt.
- **Präzisierung (ändert den Ausgang nicht).** Der Klammer-Zusatz des Triggers
  nennt als Mechanismus „eine der beiden Seiten wird geschärft, die andere
  nicht nachgezogen". Der real eingetretene Mechanismus ist ein **anderer**:
  die Divergenz ist **strukturell** — der Hook liest die rohe Datei, das
  Modul die bereinigte Message; sie bestand von Beginn an, sie entstand nicht
  durch eine einseitige Schärfung. Die Trigger-*Bedingung* („weichen real
  auseinander") ist erfüllt; der *beispielhaft* genannte Mechanismus war
  nicht der ihre. Der Ausgang bleibt derselbe — die vom Trigger geforderte
  Prüfung ist erfolgt.

## (2) Drei Befunde am `ADR-0069`-Text — je ein Verdikt

### F-5 (MEDIUM) — Reichweite des Teil-Supersedes → **Folge-ADR, Lesart „zu weit gefasst"**

**Verdikt: Verdikt 2 (die Entscheidung wird per Folge-ADR `supersedes`d).**
Die Korrektur ist `ADR-0070` Klausel 1.

Begründung: `ADR-0062` Punkt 3 trägt **drei** Aussagen — (i) bash-only/kein
Docker-Aufruf, (ii) „spiegelt exakt die zwei bestehenden Regeln",
(iii) Merge-/Revert-Ausnahme. `ADR-0069`s Kopf-Satz nennt zur Bestimmung des
superseded Punktes allein (ii), schreibt aber „nur deren Entscheidung
**Punkt 3**" — die Apposition ist zweideutig. Die engere Lesart trägt und ist
gewählt: `ADR-0069` supersedes **nur die Spiegelungs-Aussage (ii)**;
`ADR-0062` Punkt 3 bleibt für (i) und (iii) in Kraft. Damit behalten beide
normativen Aussagen einen ausgesprochenen `Accepted`-Träger, und die zwei
Plan-Zeiger (`.md:81` bash-only, `.md:116` Merge-/Revert-Ausnahme) bleiben
**gültig** — kein Plan-Nachzug. Die weitere Lesart (Punkt 3 ganz abgelöst)
ist verworfen: sie ließe (i)/(iii) ohne ausgesprochenen Träger und machte die
zwei Plan-Zeiger veraltet, und sie müsste zwei unveränderte Aussagen ohne
Grund aus `ADR-0062` in eine neue ADR duplizieren.

### F-6 (INFO) — die benannte Sicherheits-Kante → **stehen lassen**

**Verdikt: der Text trägt; keine Änderung.** Kein Folge-ADR-Zug für F-6.

Begründung: Die Klausel ist ausdrücklich als „**benannt statt behauptet**"
formuliert (`.md:94-96`) — sie *behauptet* die Erreichbarkeit der Kante
nicht, sie benennt den einen Punkt, an dem die Sicherheits-Zusage brechen
könnte, und markiert ihn als vorsichtig. Der Bestätigungslauf hat sie in drei
Konstruktionen nicht erreicht gemacht und den Grund als **strukturell**
beschrieben — aber selbst benannt als **kein Unmöglichkeitsbeweis**
(`review-slice-073-fixrunde.md` F-6, `verify-slice-073.md` §6). Drei nicht
erreichte Konstruktionen plus Quelltext-Lesen tragen keine Streichung: die
Kante auf dieser Grundlage zu entfernen wäre die **Überbehauptung**, gegen
die `ADR-0069` gerade geschrieben ist. Die vorsichtige Benennung ist die
ehrliche Form und bleibt.

### F-7 (INFO) — die Klassen-Liste ist nicht erschöpfend → **Folge-ADR**

**Verdikt: Verdikt 2 (Folge-ADR).** Die Korrektur ist `ADR-0070` Klausel 2.

Begründung: `ADR-0069` §Entscheidung Punkt 2 sagt, die drei Klassen seien
„**die festgeschriebene Grenze**", und §Re-Evaluierungs-Trigger (a) ist auf
„einer der **drei**" geschlüsselt. Eine vierte, ebenfalls laxere Klasse ist
erreichbar und in diesem Zug am gepinnten Digest selbst gemessen (s. u.):
ein Commit, dessen rohe Message mit einer Leerzeile beginnt, auf die
`Merge branch 'vp'` folgt, unter `--cleanup=verbatim` — der Hook lässt durch
(Exit 0), das Modul meldet `commit-untraceable` (Exit 1). Die Klausel liest
sich als erschöpfend und ist es nicht — derselbe Fehlertyp wie das von
`ADR-0069` korrigierte „spiegelt exakt". `ADR-0070` fasst die Grenze als
**gemessen, nicht erschöpfend**, benennt die vierte Klasse und erweitert den
Trigger (a) auf „eine weitere Klasse wird gemessen".

## Eigene Messungen dieses Zugs

Nur zur Bestätigung von F-7; F-5 und F-6 sind Text-/Rollen-Fragen. Alle
Läufe außerhalb des Arbeitsbaums (`/tmp/arch0070/`), gepinntes Image
`ghcr.io/pt9912/d-check@sha256:18e9cd85…` (Digest aus `d-check.mk`),
Wegwerf-Repository **mit** `commits`-Konfiguration (Positiv-Kontrolle
schließt das inerte Modul aus).

| Messung | Exit (Modul) | Exit (Hook) | Richtung |
|---|---|---|---|
| vierte Klasse: rohe Message `\nMerge branch 'vp'`, `--cleanup=verbatim` | **1** (`commit-untraceable`) | **0** | laxer — Gate fängt, Hook lässt durch |
| Gegenprobe: derselbe Aufbau, Default-Cleanup (`git` strippt die Leerzeile) | 0 | 0 | beide grün — Klasse nur unter `verbatim` |
| Positiv-Kontrolle: Message ohne Kennung, `--commit-msg` | **1** (`commit-untraceable`) | — | Konfiguration nicht inert |

## Folgearbeit

| Posten | Träger | Was |
|---|---|---|
| `ADR-0070` (Folge-ADR, drei Klauseln) + Index-Zeile | **Architect (dieser Zug, erledigt)** | `Supersedes ADR-0069` (Supersede-Satz · Grenz-Klausel §Entscheidung Punkt 2 · Re-Eval (a)); Index-Zeile `ADR-0070`, Zeiger `(→ ADR-0070, teilweise)` an der `ADR-0069`-Zeile |
| Plan-Nachzug F-7 | **Implementer** | Slice-Plan `slice-073` §Bezug (`.md:20`), §1 (`.md:55-60`), §3 (`.md:151`) nennen „drei … Divergenz-Klassen" als Grenze → auf die **benannte, nicht als erschöpfend bewiesene** Grenze nachziehen (vierte Klasse nennen oder „nicht erschöpfend" aussprechen). **Nicht** nachzuziehen: Hook-Kommentar und `harness/README.md`-Zeile — beide nennen keine geschlossene Drei-Liste („vollstaendig ist er nicht" / „fängt **nicht** jeden Verstoß") und bleiben korrekt |
| Plan-Nachzug F-5 | **niemand** | entfällt — Punkt 3 bleibt in Kraft, die zwei Zeiger `.md:81`/`.md:116` sind korrekt |
| §7-Notiz, Register-Ausgang `verkörpert`/`seit slice-073`, §6-Risiko-Ausgänge, drei Paarungen, §4-Rückführungs-Grund | **Planner (Closure)** | Modul 5/6; der Register-Vorschlag `spiegelung-ist-approximation` (1×) und der §6-Ausgang zu Risiko 1 (Trigger (b) eingetreten → Ausgang `ADR-0069`/`ADR-0070`) sind Planner-Entscheid |

## Was dieser Zug NICHT tut

- **Kein** Verhaltens-Change: `.githooks/commit-msg`,
  `tools/harness/commit-traceability.sh`, `.d-check.yml` (`commits`-Abschnitt)
  und `GATE_CHECKS` bleiben unverändert.
- **Keine** Änderung an `ADR-0062` oder `ADR-0045` (Immutabilität,
  `AGENTS.md` §3.5); die Korrektur des `ADR-0069`-Textes ist `ADR-0070`.
- **Kein** Plan-Nachzug (Implementer) und **keine** Register-Einträge/§7-Notiz
  (Planner).
- **Kein** Reviewer-Skill-Patch (die Klasse tritt unter der Schwelle auf).

## Verantwortlicher dieses Zugs

**Rolleninhaber:** der unabhängige Architect-Lauf, der dieses Verdikt und
[`ADR-0070`](../plan/adr/0070-supersede-reichweite-und-klassengrenze.md)
geschrieben hat. Der nächste Zug ist der **Implementer** (Plan-Nachzug F-7),
dann die **Closure** durch den **Planner** (`git mv` nach `done/`, §7-Notiz,
Register-Ausgang, §6-Risiko-Ausgänge, drei Paarungen).
