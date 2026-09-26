# Slice harness-suchlauf-nachmessen: Nachmess-Werkzeug für das Suchlauf-Feld der Slice-Pläne und Rücknahme der Rollout-Artefakte in den DB-Tier-Läufen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er geht `slice-transformationen-kern-rename`
voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5).

**Bezug:** [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen in Trägern; die Grenze „kein Sensor auf Prosa“ bleibt),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) und
[`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) (der
Schema-Rollout, dessen Erzeugnisse die DB-Tier-Läufe hinterlassen),
[`AGENTS.md`](../../../../AGENTS.md) §3.12 und §3.13, Architect-Verdikt
[`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§3.2 und §7 R3.

**Berührte Spec-Stellen:** — (Harness-Werkzeug und Test-Skript; keine
Spec-Stelle).

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Closure der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Suchlauf-Feld eines Slice-Plans ist maschinell nachmessbar: jede
Suchzeile steht als Zeile eines Codeblocks mit dem Etikett `suchlauf`, und
`make suchlauf-nachmessen PLAN=<Datei>` führt sie aus und druckt Soll und Ist.
Dazu stellt `tools/schema/apply-rollout.sh` die Rollout-Artefakte
`tools/schema/plan.yaml` und `tools/schema/down.sql` nach dem Lauf wieder her,
so dass `make test-store` und `make test-replication` den Arbeitsbaum
unverändert lassen.

**Form der Suchzeile** (ein Codeblock je Plan-Feld, eine Zeile je Messung; das
Etikett `suchlauf` kennzeichnet, was das Werkzeug ausführt — dieses Beispiel
trägt deshalb das Etikett `text`):

```text
<Stand> <erwartete Zeilenzahl> <Argumente von git grep>
```

`<Stand>` ist eine Commit-Kennung (der Parent, nie `HEAD`) oder das Wort `diff`
(der Stand des Arbeitsbaums beim Nachmessen). Die Plan-Datei ist immer aus dem
Suchraum ausgeschlossen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Gate.** Der Stand `diff` bewegt sich mit jedem Commit; die Messung
  gehört in den Lauf dessen, der sie braucht (Implementer, Reviewer, Verifier,
  Planner). Ein Gate über eine bewegliche Zahl färbte sich ohne Ursache rot.
- **Ein Parser über die vorhandenen Tabellenzellen der Pläne.** Die Befehle
  stehen dort mit dem Escape `\|`, der als `-E`-Muster kopiert null Treffer
  liefert, und ein Befund trägt mehrere Zahlen je Zelle in Fließtext; die Form
  ist nicht tragfähig lesbar. Bestehende Pläne wandern nicht um: ein Plan trägt
  die Blöcke, sobald ein Implementer sein Feld füllt.
- **Die Prüfung der Vollständigkeit von Suchraum und Suchmuster.** Das
  Werkzeug wiederholt eine deklarierte Messung (Befehl, Stand, Zahl); ob das
  Muster den Symbolnamen, das Zählwort und die Beschreibung trägt, bleibt eine
  Lese-Handlung des Reviewers ([`AGENTS.md`](../../../../AGENTS.md) §3.13
  §Suchform). Ein Sensor „jede Zahl trägt ihren Ursprung“ bleibt
  ausgeschlossen
  ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
  §Entscheidung 4).
- **Ein Mutations-Harness.** Docker-only-Pinnung eines Werkzeugs samt Laufzeit
  je Paket und äquivalente Mutanten; die Regel wirkt ohne es (Architect-Verdikt
  §3.4).
- **Andere Schreiber in committete Dateien.** Der Slice deckt die zwei Dateien,
  die der DB-Tier-Lauf verändert; ein weiterer Schreiber ist ein neuer Beleg
  in `BEO-PGC/test-schreibt-in-committete-datei`.

## 2. Definition of Done

- [ ] Das Werkzeug steht: `tools/harness/suchlauf-nachmessen.sh <Plan-Datei>`
      liest die `suchlauf`-Blöcke, führt je Zeile `git grep -n <Argumente>
      <Stand> -- . ':!<Plan-Datei>'` aus (Stand `diff`: der Arbeitsbaum), zählt
      die Trefferzeilen, druckt Befehl, Soll und Ist und endet mit Exit ≠ 0 bei
      jeder Abweichung; ein Stand `HEAD` und ein Plan ohne `suchlauf`-Block enden
      ebenfalls mit Exit ≠ 0 (leer ist nicht bestanden). `make
      suchlauf-nachmessen PLAN=<Datei>` ruft es auf; ohne `PLAN` bricht das
      Ziel mit `$(error …)` ab, bevor ein Befehl läuft. *Zu belegen durch:*
      `make test-suchlauf-nachmessen` (Tabellentest, netzlos) mit fünf Fällen —
      stimmt · weicht ab · Selbstverweis ausgeschlossen · `HEAD` abgelehnt ·
      kein Block —, je an seine Eingabe gebunden: die Mutation „Ausschluss der
      Plan-Datei entfernt“ färbt den Selbstverweis-Fall rot, „Zahlenvergleich
      umgekehrt“ den Abweichungs-Fall (Mutationen und gesehenes Rot im
      Bericht).
- [ ] Die Träger nennen das Werkzeug und seine Grenze:
      `harness/sensors/suchlauf-nachmessen.md` (Vertrag, Form der Zeile, Grenze
      „prüft Zahlen und Stände, nicht die Vollständigkeit von Suchraum und
      Muster“, kein Gate), eine Zeile im Werkzeug-Verzeichnis von
      [`harness/README.md`](../../../../harness/README.md) §Sensors,
      [`AGENTS.md`](../../../../AGENTS.md) §3.13 (§Suchform nennt Form und
      Aufruf), `.claude/commands/implement-slice.md` Schritt 18 (Aufruf nach
      jeder Fixrunde) und der Probe-Satz des Reviewer-Punkts „Zahl im Träger …“
      in `.harness/skills/reviewer.md` (Nachmessen mit dem Werkzeug, wo eine
      `suchlauf`-Zeile vorliegt). *Zu belegen durch:* Lesen der fünf Stellen und
      `make docs-check`.
- [ ] Die Rollout-Artefakte bleiben unverändert: `tools/schema/apply-rollout.sh`
      sichert `tools/schema/plan.yaml` und `tools/schema/down.sql` vor
      `make schema-rollout` in ein `mktemp`-Verzeichnis und stellt sie nach dem
      Lauf wieder her, auch bei einem Fehlschlag (Vorbild: der Guard-Test
      `tools/harness/run-schema-rollout-guard-test.sh`, der beide sichert und
      in `cleanup` zurückstellt); eine vor dem Lauf lokal geänderte Datei
      bleibt in diesem Zustand. *Zu belegen durch:* ein realer `make
      test-store`- und ein realer `make test-replication`-Lauf, danach `git
      status --short tools/schema` ohne Ausgabe; die Mutation „Wiederherstellung
      entfernt“ färbt die Prüfung rot (der Lauf lässt beide Dateien geändert,
      der Parent-Stand tut es: `BEO-PGC/test-schreibt-in-committete-datei`, 4
      Belege).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt für das Benutzerhandbuch — keine
      Betreiber-Oberfläche; die Doku-Träger stehen im zweiten Liefer-Punkt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten; ein Slice ohne Welle wird von ihr
      mitgeprüft).

**Umfang:** M — Schätzung des Architects, nicht gemessen
([`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§3.2: ein Skript von etwa 60 bis 100 Zeilen, ein Skript-Test mit mehreren
Fällen, ein Makefile-Ziel, eine Sensor-Doku, die Träger-Zeilen; dazu die
Rücknahme in `apply-rollout.sh`, drei Liefer-Punkte).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/suchlauf-nachmessen.sh` | neu | liest die Blöcke, `git grep` je Zeile, Soll/Ist, Exit-Code; die Plan-Datei ist immer ausgeschlossen. |
| `tools/harness/run-suchlauf-nachmessen-tests.sh` | neu | Tabellentest mit den fünf Fällen des ersten Liefer-Punkts, an die Eingabe gebunden. |
| `Makefile` | update | Ziele `suchlauf-nachmessen` (Pflicht-Argument `PLAN`) und `test-suchlauf-nachmessen`; Zeilen in `make help`. |
| `harness/sensors/suchlauf-nachmessen.md` | neu | Vertrag und benannte Grenze; der Name im Verzeichnis der Sensor-Dokus. |
| `harness/README.md` | update | eine Zeile in der Werkzeug-Tabelle §Sensors („kein Gate“). |
| `AGENTS.md` §3.13, `.claude/commands/implement-slice.md` Schritt 18, `.harness/skills/reviewer.md` | update | Träger nennen Form, Aufruf und Grenze des Werkzeugs. |
| `tools/schema/apply-rollout.sh` | update | Sicherung und Wiederherstellung von `plan.yaml` und `down.sql` (auch im Fehlerfall, `trap`). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Form des
Suchlauf-Felds und das Werkzeug, das es nachmisst“ und „welche Läufe
`tools/schema/plan.yaml` und `down.sql` verändern“; beide Stände gemessen; die
Befehle stehen im Codeblock, der Implementer trägt Stand und Trefferzahl ein):**

```text
git grep -n -i -E 'suchlauf|nachmess' -- AGENTS.md .claude .harness/skills harness Makefile tools docs/user
git grep -n -E 'plan\.yaml|down\.sql' -- tools harness Makefile docs/user
```

| Träger | Befund | Behandlung |
|---|---|---|
| Beschreibungen der Suchform (AGENTS.md, Commands, Skill, Sensor-Doku) | *(Implementer trägt ein)* | jede nennt Form und Aufruf des Werkzeugs; der Codeblock-Satz von §3.13 bleibt die Regel |
| Läufe, die `plan.yaml`/`down.sql` schreiben | *(Implementer trägt ein)* | `apply-rollout.sh` deckt `make test-store` und `make test-replication`; `tools/bench-lib.sh` und `tools/harness/run-integration-tests.sh` rufen `make schema-rollout` direkt (gemessen für `tools/bench-backfill.sh`: danach zeigt `git status --short` `M tools/schema/plan.yaml`, Closure von `welle-backfill-bestand`); ob die Rücknahme in `apply-rollout.sh` oder im Target `schema-rollout` selbst liegt, entscheidet der Implementer, jeder weitere Schreiber wird als Beleg gemeldet |

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/`
(WIP-Limit 1). Der Slice muss `done` sein, **bevor**
`slice-transformationen-kern-rename` startet (Start-Trigger dort): das
Suchlauf-Feld der zehn Pläne der Welle
[welle-transformationen](../welle-transformationen.md) wird mit dem Werkzeug
nachgemessen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Werkzeug,
  Träger-Zeilen und Rollout-Rücknahme nicht in einem Review tragen — der
  abtrennbare Teil ist der dritte Liefer-Punkt (Rücknahme in
  `apply-rollout.sh`), unabhängig vom Werkzeug, als eigener Slice.
- `in-progress` → `open` (blockiert): falls `git` im Werkzeug nur in einem
  gepinnten Image zulässig ist, das dieses Repo nicht führt (Architect-Frage zu
  [`AGENTS.md`](../../../../AGENTS.md) §3.1, siehe §6).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-suchlauf-nachmessen` grün +
ein realer `make test-store`- und `make test-replication`-Lauf mit leerem `git
status --short tools/schema` danach + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Docker-only für `git`** ([`AGENTS.md`](../../../../AGENTS.md) §3.1): die
  vorhandenen Skripte unter `tools/harness/` rufen `git` auf dem Host auf
  (`cd "$(git rev-parse --show-toplevel)"`, gelesen in `apply-rollout.sh` und
  `run-release-tag-info-tests.sh`); das Werkzeug folgt ihnen. *Erwartet, zu
  belegen durch:* Lesen der Skript-Köpfe am Start; ergibt sich eine Vorgabe, die
  dagegen spricht, ist es die Rückführung nach `open/`. **Ausgang:** *(bei
  Closure)*
- **Die Zahl streut mit dem Stand `diff`**: er bewegt sich mit jedem Commit, ein
  grüner Lauf gilt für den Arbeitsbaum des Laufs. *Erwartet, zu belegen durch:*
  die Sensor-Doku nennt den Stand im Aufruf-Beleg; kein Gate. **Ausgang:** *(bei
  Closure)*
- **Das Werkzeug erweckt den Eindruck, das Feld sei vollständig** (nur Zahlen
  und Stände sind geprüft). *Erwartet, zu belegen durch:* die benannte Grenze in
  Sensor-Doku und §3.13; der Reviewer prüft Suchraum und Muster weiter von Hand
  (`BEO-PGC/regel-weiter-als-ihr-sensor`, 3×, verkörpert teilweise). **Ausgang:** *(bei
  Closure)*
- **Die Rücknahme überschreibt eine gewollte lokale Änderung an `plan.yaml`.**
  *Erwartet, zu belegen durch:* ein Test mit vorab geänderter Datei — sie bleibt
  im Zustand vor dem Lauf. **Ausgang:** *(bei Closure)*
- **Die bestehenden Pläne tragen keine `suchlauf`-Blöcke.** *Erwartet, zu
  belegen durch:* der Implementer eines Slice überträgt sein Feld am Start in die
  Blöcke; das Werkzeug meldet einen Plan ohne Block mit Exit ≠ 0. **Ausgang:**
  *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); `tools/harness`, `tools/schema` und die Doku-Träger
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-25 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/test-schreibt-in-committete-datei` (4×, Ausgang: dieser Slice, dritter
Liefer-Punkt), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert,
21×, Suchlauf-Anteil: dieser Slice), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(verkörpert, 13×, der Befehl in der Tabellenzelle: der Codeblock),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 31×, Suchform),
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 8×),
`BEO-PGC/regel-weiter-als-ihr-sensor` (3×, verkörpert teilweise, einschlägig — die
benannte Grenze des Werkzeugs), `BEO-PGC/d-migrate-nacharbeit` (verkörpert, 7×,
der Rollout-Pfad wird berührt, nicht seine Objektklassen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
