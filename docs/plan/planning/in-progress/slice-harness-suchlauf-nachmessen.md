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
Dazu stellt der Wrapper `tools/schema/rollout-restore.sh` die Rollout-Artefakte
`tools/schema/plan.yaml` und `tools/schema/down.sql` nach dem Lauf wieder her,
der `make schema-rollout` als Vorbedingung ruft (`tools/schema/apply-rollout.sh`
und sechs weitere direkte Aufrufer), so dass `make test-store` und
`make test-replication` den Arbeitsbaum unverändert lassen.

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

- [x] Das Werkzeug steht: `tools/harness/suchlauf-nachmessen.sh <Plan-Datei>`
      liest die `suchlauf`-Blöcke, führt je Zeile `git grep -n <Argumente>
      <Stand> -- . ':!<Plan-Datei>'` aus (Stand `diff`: der Arbeitsbaum), zählt
      die Trefferzeilen, druckt Befehl, Soll und Ist und endet mit Exit ≠ 0 bei
      jeder Abweichung; ein Stand `HEAD` und ein Plan ohne `suchlauf`-Block enden
      ebenfalls mit Exit ≠ 0 (leer ist nicht bestanden). `make
      suchlauf-nachmessen PLAN=<Datei>` ruft es auf; ohne `PLAN` bricht das
      Ziel mit `$(error …)` ab, bevor ein Befehl läuft. *Zu belegen durch:*
      `make test-suchlauf-nachmessen` (Tabellentest, netzlos) mit elf Fällen —
      stimmt · weicht ab · Selbstverweis ausgeschlossen · `HEAD` abgelehnt ·
      kein Block · Commit-Stand mit Pathspec · erlaubte Optionen · Optionen und
      Pathspec-Magic außerhalb der Allow-List · `git grep`-Fehler · nicht
      geschlossener Block · Zeilenform —, je an seine Eingabe gebunden: die
      Mutation „Ausschluss der Plan-Datei entfernt“ färbt den
      Selbstverweis-Fall rot, „Zahlenvergleich umgekehrt“ den Abweichungs-Fall,
      „Optionsprüfung entfernt“ den Fall, in dem ein `-O`-Kommando seine
      Marker-Datei anlegt (Mutationen und gesehenes Rot im Bericht).
- [x] Die Träger nennen das Werkzeug und seine Grenze:
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
- [x] Die Rollout-Artefakte bleiben unverändert: der Wrapper
      `tools/schema/rollout-restore.sh` sichert `tools/schema/plan.yaml` und
      `tools/schema/down.sql` vor dem Kommando in ein `mktemp`-Verzeichnis und
      stellt sie nach dem Lauf wieder her, auch bei einem Fehlschlag;
      `tools/schema/apply-rollout.sh` und die sechs weiteren direkten Aufrufer
      von `make schema-rollout` rufen es durch den Wrapper (Vorbild: der Guard-Test
      `tools/harness/run-schema-rollout-guard-test.sh`, der beide sichert und
      in `cleanup` zurückstellt); eine vor dem Lauf lokal geänderte Datei
      bleibt in diesem Zustand. *Zu belegen durch:* ein realer `make
      test-store`- und ein realer `make test-replication`-Lauf, danach `git
      status --short tools/schema` ohne Ausgabe; die Mutation „Wiederherstellung
      entfernt“ färbt die Prüfung rot (der Lauf lässt beide Dateien geändert,
      der Parent-Stand tut es: `BEO-PGC/test-schreibt-in-committete-datei`, 4
      Belege).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt für das Benutzerhandbuch — keine
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
| `tools/harness/run-suchlauf-nachmessen-tests.sh` | neu | Tabellentest mit den fünf Fällen des ersten Liefer-Punkts, an die Eingabe gebunden; in der Fixrunde nach dem Review um sechs Fälle erweitert (Commit-Stand mit Pathspec, erlaubte Optionen, Optionen und Pathspec-Magic außerhalb der Allow-List mit Marker-Datei, `git grep`-Fehler, nicht geschlossener Block, Zeilenform). |
| `Makefile` | update | Ziele `suchlauf-nachmessen` (Pflicht-Argument `PLAN`) und `test-suchlauf-nachmessen`; Zeilen in `make help`. |
| `harness/sensors/suchlauf-nachmessen.md` | neu | Vertrag und benannte Grenze; der Name im Verzeichnis der Sensor-Dokus. |
| `harness/README.md` | update | eine Zeile in der Werkzeug-Tabelle §Sensors („kein Gate“). |
| `AGENTS.md` §3.13, `.claude/commands/implement-slice.md` Schritt 18, `.harness/skills/reviewer.md` | update | Träger nennen Form, Aufruf und Grenze des Werkzeugs. |
| `tools/schema/apply-rollout.sh` | update | Sicherung und Wiederherstellung von `plan.yaml` und `down.sql` (auch im Fehlerfall, `trap`) — nicht im Skript selbst, sondern über `rollout-restore.sh` (nächste Zeile): das Skript ruft `make schema-rollout` durch den Wrapper. |
| `tools/schema/rollout-restore.sh` | neu (Nachzug) | Wrapper `rollout-restore.sh <Kommando…>`: sichert `plan.yaml` und `down.sql` in ein `mktemp`-Verzeichnis, führt das Kommando aus, stellt beide danach wieder her (`trap … EXIT`, auch bei Fehlschlag; lokal geänderte Datei bleibt so, fehlende bleibt fehlend), Exit-Code des Kommandos. Ort der Rücknahme: weder im Target `schema-rollout` — im Betrieb sind Report und Rollback-Artefakt sein Erzeugnis, das der Betreiber behält — noch in `apply-rollout.sh` allein — sechs weitere Skripte rufen `make schema-rollout` direkt (Aufrufer-Messung: Feld §3.13, Zeile 7/8). Ein Wrapper deckt alle sieben Aufrufer mit einer Stelle; §1 und der dritte DoD-Punkt nennen ihn als Ort der Rücknahme. |
| `tools/harness/run-integration-tests.sh`, `tools/harness/run-sdk-csharp-integration-tests.sh`, `tools/harness/run-sdk-kotlin-integration-tests.sh`, `tools/harness/run-sdk-python-integration-tests.sh`, `tools/bench-lib.sh`, `examples/bootstrap.sh` | update (Nachzug) | je ein `make schema-rollout`-Aufruf läuft durch `rollout-restore.sh`; die Kommentare zu Report und Rollback-Artefakt in `run-integration-tests.sh` und `apply-rollout.sh` nennen die Rücknahme. `tools/harness/run-schema-rollout-guard-test.sh` bleibt (sichert und stellt beide Dateien selbst wieder her, Aufrufer-Prüfung nimmt es aus). |
| `tools/harness/run-rollout-restore-tests.sh`, `Makefile` (`test-rollout-restore`) | neu (Nachzug) | Tabellentest für den Wrapper (schreibt · scheitert · lokale Änderung bleibt · fehlende Datei bleibt fehlend · ohne Kommando) und die Aufrufer-Prüfung im Repo: jede Shell-Zeile unter `tools/`/`examples/` mit `make … schema-rollout` trägt `rollout-restore.sh`. |
| `harness/targets/schema-rollout.md` | update (Nachzug) | Abschnitt „Erzeugnisse in Test-, Bench- und Beispiel-Läufen“: welche Läufe die Rücknahme tragen. |
| Ausschluss der Plan-Datei | Abweichung vom Plan-Wortlaut | Der Plan nennt `':!<Plan-Datei>'`; das Werkzeug schließt sie über ihren **Dateinamen in jedem Verzeichnis** aus (`:(exclude,glob)**/<Dateiname>`): der Plan liegt am Parent-Stand unter einem anderen Lifecycle-Verzeichnis, ein Pfad-Ausschluss zählte dort ihn selbst mit (im Test `stimmt`/`Selbstverweis` an einem Parent mit dem Plan unter `next/` gebunden). Der Pfad-Ausschluss allein bindet keinen Fall (Mutation ohne ihn: alle fünf Fälle grün). |
| Form der Zeile | Festlegung | `<Stand>` ist eine Hex-Kennung (7 bis 40 Ziffern) oder `diff`; ein alleinstehendes `--` trennt Optionen und Muster vom Pathspec (der Stand steht vor dem `--`, das Werkzeug setzt ihn); `'…'`/`"…"` gruppieren ohne Expansion (eigener Zerleger, kein `eval`). Eine Plan-Zeile führt keinen Kommando-Code aus: Optionen vor `--` stehen auf einer Allow-List (`-e <Muster>`, `-i -w -E -F -G -P -v -I -a -h -H -o -c -l -L -n`, `--and --or --not --all-match` und die Langformen der Muster-Art; jede andere Option, darunter `-O`, `--open-files-in-pager`, `-f`, `--no-index`, endet mit Exit 2), Pathspec-Magic ist auf `:!`, `:^` und `:(exclude\|glob\|literal\|icase\|top)` begrenzt, ein leeres Argument wird abgelehnt (Beleg: `make test-suchlauf-nachmessen`, Fall mit Marker-Datei). Kein Namens- oder `HEAD`-Stand: er bewegt sich. |
| `tools/schema/rollout-restore.sh`, `harness/targets/schema-rollout.md` | update (Fixrunde) | Die Grenzen des Wrappers stehen im Skriptkopf und im Vertrag: zwei gleichzeitige Aufrufe teilen dieselben zwei Dateien, ein `SIGKILL` lässt die Erzeugnisse verändert, die Aufrufer-Prüfung liest nur Einzelzeilen. |
| `harness/sensors/suchlauf-nachmessen.md`, `harness/README.md`, `Makefile` | update (Fixrunde) | Vertrag nennt die erlaubten Optionen, die Pathspec-Magic, das leere Argument, die getrennte Standardfehlerausgabe und die Host-Werkzeuge (`bash`, `git`, `realpath`); Testzahl elf in Sensor-Doku, README-Zeile und `make help`. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Form des
Suchlauf-Felds und das Werkzeug, das es nachmisst“ und „welche Läufe
`tools/schema/plan.yaml` und `down.sql` verändern“; beide Stände gemessen; die
Befehle stehen im Codeblock, der Implementer trägt Stand und Trefferzahl ein):**

```suchlauf
8717c4fb 14 -i -E 'suchlauf|nachmess|Suchform' -- AGENTS.md .claude .harness/skills harness Makefile tools docs/user spec README.md
diff 80 -i -E 'suchlauf|nachmess|Suchform' -- AGENTS.md .claude .harness/skills harness Makefile tools docs/user spec README.md
8717c4fb 9 -i -E 'Suchform|Suchlauf.*(Codeblock|Tabellenzelle)|(Codeblock|Tabellenzelle).*Suchlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 12 -i -E 'Suchform|Suchlauf.*(Codeblock|Tabellenzelle)|(Codeblock|Tabellenzelle).*Suchlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
8717c4fb 59 -E 'plan\.yaml|down\.sql' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 74 -E 'plan\.yaml|down\.sql' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
8717c4fb 35 -E 'make .*schema-rollout|apply-rollout' -- tools examples Makefile harness .github
diff 47 -E 'make .*schema-rollout|apply-rollout' -- tools examples Makefile harness .github
```

Suchraum: Zeile 1/2 die Träger der Regel und des Werkzeugs (Planungs-Bäume
nennen den Suchlauf als Handlung eigener Slices in über 60 Dateien und
beschreiben seine Form nicht — Zeile 3/4 misst das ganz-baum-weit auf die Form);
Zeile 5/6 und 7/8 ganzer Baum bzw. die Aufrufer-Bäume, ohne Berichte, Records
und Baseline. Stände: Parent `8717c4fb`, Diff der Arbeitsbaum des Laufs.

| Träger | Befund | Behandlung |
|---|---|---|
| Beschreibungen der Suchform (AGENTS.md, Commands, Skill, Sensor-Doku) | Zeile 1/2: Parent 14 Zeilen in 3 Dateien (`AGENTS.md` 9, `.harness/skills/reviewer.md` 2, `docs/user/benutzerhandbuch.md` 3); Diff 80 Zeilen in 9 Dateien. Die drei Handbuch-Zeilen (Nachmessung eigener Messläufe) beschreiben keine Suchform, unverändert. Neu: `.claude/commands/implement-slice.md` 5, `harness/README.md` 3, `harness/sensors/suchlauf-nachmessen.md` 10, `Makefile` 7, die zwei Werkzeug-Skripte 35; `AGENTS.md` 13, `.harness/skills/reviewer.md` 4. Zeile 3/4 (ganzer Baum, Form): Parent 9, Diff 12 — die neun Parent-Zeilen stehen in `AGENTS.md` (2) und in Register-`state.md`/`evidence` (7), nicht in den Regel-Trägern. Nicht gefunden: eine weitere Beschreibung der Suchform in `.claude/commands/*` außer `implement-slice.md` (Parent 0 Zeilen mit `suchlauf`), in `harness/sensors/*` und im Handbuch. | die vier Träger nennen Form und Aufruf (`AGENTS.md` §3.13, `implement-slice.md` Schritt 18, Reviewer-Skill, Sensor-Doku), `harness/README.md` das Werkzeug; der Codeblock-Satz von §3.13 bleibt die Regel. Register-Verweise (`state.md` in vier Einträgen) trugen einen Link auf `open/slice-harness-suchlauf-nachmessen.md`, der mit dem Lifecycle-Move bricht (`make docs-check`: `target-missing`); als Kennungs-Zitat ersetzt, Zustandstext unverändert — der Zustand („geplant“) bleibt Closure-Arbeit des Planners (gemeldet, Frist: Closure dieses Slice) |
| Läufe, die `plan.yaml`/`down.sql` schreiben | Zeile 7/8 (Aufrufer-Bäume): Parent 35 Zeilen in 15 Dateien, Diff 47 in 18 (neu: `rollout-restore.sh` 3, `run-rollout-restore-tests.sh` 3, `harness/sensors/suchlauf-nachmessen.md` 1; `harness/README.md` 3→4, `harness/targets/schema-rollout.md` 2→6). Direkte Aufrufer von `make schema-rollout` außerhalb von `apply-rollout.sh` (Parent): `tools/bench-lib.sh`, `tools/harness/run-integration-tests.sh`, `tools/harness/run-sdk-{csharp,kotlin,python}-integration-tests.sh`, `examples/bootstrap.sh` — sechs Dateien statt der zwei, die der Plan nennt — und `tools/harness/run-schema-rollout-guard-test.sh` (eigene Sicherung). Zeile 5/6 (ganzer Baum): Parent 59, Diff 74 Zeilen; die Zunahme sind die neuen Träger der Rücknahme. Nicht gefunden: ein weiterer Schreiber in `plan.yaml`/`down.sql` (`.github/workflows` ruft `make schema-rollout` nicht, `e2e.yml` nennt `apply-rollout.sh` in einem Kommentar). | Die Rücknahme liegt in `tools/schema/rollout-restore.sh`; `apply-rollout.sh` (deckt `make test-store` und `make test-replication`) und die sechs direkten Aufrufer gehen durch sie, der Guard-Test behält seine eigene Sicherung (§3 Zeilen oben). Die Aufrufer, die dieser Lauf nicht real gefahren hat (`make test-integration`, drei `make test-sdk-*-integration`, `make bench`, `make example-demo-up`), trägt die Aufrufer-Prüfung von `make test-rollout-restore` (jede Aufruf-Zeile geht durch den Wrapper) und `bash -n`; ein weiterer Schreiber ist nicht gefunden. |

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
  abtrennbare Teil ist der dritte Liefer-Punkt (Rücknahme der
  Rollout-Artefakte über den Wrapper `rollout-restore.sh`), unabhängig vom
  Werkzeug, als eigener Slice.
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
