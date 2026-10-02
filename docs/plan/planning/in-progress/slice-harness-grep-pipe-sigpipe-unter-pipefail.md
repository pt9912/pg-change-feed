# Slice harness-grep-pipe-sigpipe-unter-pipefail: Runner-Skripte lesen Container-Logs ohne `grep -q` in der Pipe — eine vorhandene Zeile wird nicht mehr verfehlt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Client-Bibliotheken; die
drei SDK-Realserver-Tiers sind Träger ihres Nachweises),
[`LH-QA-POR-003`](../../../../spec/lastenheft.md) (der E2E-Lauf von
`tools/harness/run-integration-tests.sh` ist der Nachweis-Weg),
[`ADR-0030`](../../adr/0030-testpyramide.md) (Testpyramide, E2E-Tier),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung
Festlegung 2 (Mechanik der SDK-Realserver-Tiers),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen:
gemessen gegen hergeleitet). Die Nachbarregel
[`AGENTS.md`](../../../../AGENTS.md) §3.9 (Exit-Code einer Pipe) behandelt eine andere
Hälfte derselben Shell-Eigenschaft (dort maskiert eine Pipe einen roten Exit, hier färbt sie
einen vorhandenen Treffer falsch). Ursprung des Slice: Review F-1 in
[`review-slice-sdk-sse-filter-phase-verbindung-haertung`](../../../reviews/review-slice-sdk-sse-filter-phase-verbindung-haertung.md)
und Verifikation §5/§6 in
[`verifikation-slice-sdk-sse-filter-phase-verbindung-haertung`](../../../reviews/verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md).

**Berührte Spec-Stellen:** — (Werkzeug-Skripte unter `tools/harness/`; keine Spec-Stelle wird
geändert).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Folge-Slice aus der Closure von
`slice-sdk-sse-filter-phase-verbindung-haertung` (Freigabe des Auftraggebers vom 2026-10-02:
„Alles angehen"). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Warte- und Prüfschleifen der Runner-Skripte lesen den Ausgang eines Containers
(`docker logs …`) so, dass eine **vorhandene** Zeile nicht verfehlt wird: kein `grep -q` (und
keine Option, die `grep` beim ersten Treffer beenden lässt) als Empfänger einer Pipe unter
`set -o pipefail`.

**Befund, aus dem der Slice entsteht** (Ursprung: Review F-1 des Vorgängers, vom Reviewer
gemessen; der Planner hat nichts davon nachgemessen):

- Im ersten C#-Tier-Lauf des Reviewers endete die NATS-Phase mit „der Ablehnungs-Beleg blieb aus
  (REJECTED token-rejected fehlt)“, obwohl die Zeile im ausgegebenen Container-Log stand; der
  zweite Lauf war grün (**gemessen** vom Reviewer, Fehlertext des Laufs).
- Mechanismus (**hergeleitet** vom Reviewer, teilweise gemessen): `grep -q` beendet sich beim
  ersten Treffer; `docker logs` schreibt die restlichen Zeilen in die geschlossene Pipe, erhält
  SIGPIPE und beendet sich mit einem Fehlerstatus; unter `pipefail` ist dann der Status der
  **Pipeline** der des ersten fehlgeschlagenen Glieds, `if` liest „kein Treffer“. Messung des
  Reviewers: ein beendeter Container mit einer Treffer-Zeile und 20 Folgezeilen,
  `docker logs | grep -qF` unter `pipefail`, **300 Aufrufe → 3 Fehlschläge** (1 %, **vom Reviewer
  gemessen**, nicht nachgemessen).
- Die Form steht im Bestand nicht nur in der Ablehnungs-Schleife: am Stand `a420e223`
  **35 Zeilen** der Form `docker logs … | grep -q…` in sechs Dateien (Suchlauf unten), dazu
  weitere `… | grep -q…`-Pipes mit anderem Erzeuger (114 Zeilen in 15 Dateien; Auswahl nach
  Lesen, siehe §3).

**Lösungsrichtung — im Plan an Hand von Lesen entschieden, am Implementer-Lauf zu bestätigen:**
Es gibt drei Formen; gelesen wurde `tools/harness/run-sdk-csharp-integration-tests.sh` (Schleife
um Zeile 246), `tools/harness/lib-sdk-filter-fixture.sh` (Zeilen 113, 136, 163) und
`tools/harness/run-integration-tests.sh` (Zeilen 2570, 2584, 3181, 3195):

1. *`grep` ohne `-q`, Ausgabe nach `/dev/null`* (`docker logs … | grep -F "READY" >/dev/null`).
   `grep` liest bis zum Dateiende, `docker logs` schreibt vollständig, kein SIGPIPE; der Status
   bleibt 0 bei Treffer, 1 ohne Treffer, die Pipeline-Form und der `if`-Rumpf bleiben. Ein
   Zeichen-Austausch je Stelle, kein neuer Hilfsprozess. **Gewählt** (kleinster Schnitt,
   mechanisch an 35 Stellen anwendbar).
2. *Ausgabe zuerst in eine Variable lesen* (`out=$(docker logs … 2>/dev/null || true)`,
   dann `grep -q … <<<"$out"`). Trägt ebenso, braucht aber je Stelle eine Variable und ein
   `|| true` (unter `set -e` bricht ein fehlschlagendes `docker logs` sonst das Skript ab); der
   Bestand nutzt diese Form bereits dort, wo das Log mehrfach gelesen wird (`rn_expect_end`,
   `run-integration-tests.sh` Zeile 5320 ff.). **Zurückgestellt** für Stellen, an denen die
   Variable ohnehin gebraucht wird.
3. *`set +o pipefail` um die Schleife.* Verworfen: schaltet die Prüfung für alle Glieder ab, auch
   für ein wirklich fehlschlagendes `docker logs`.

Die Entscheidung gilt bis zur Messung (§2 Liefer-Punkt 1); bleibt die Form 1 an der Reproduktion
rot, entscheidet der Implementer begründet auf Form 2 und trägt das im Plan nach.

**Verhältnis zur Beobachtung `BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot`** (ein einmaliger
roter Lauf des Workflows `e2e` in der Phase „Routing-Nichtanwendbarkeit und Abhilfe (b)“):
**Zusammenhang unbewiesen.** Gelesen am Stand `a420e223`: `rn_abhilfe` und `rn_expect_end`
(`run-integration-tests.sh` Zeilen 5296 bis 5345) tragen **keine** `… | grep -q…`-Pipe; das
Log wird dort in eine Variable gelesen. Auch die Hilfen `bf_await_healthy`, `bf_await_applied` und
`bf_await_sql` (Zeilen 3427 bis 3510) lesen keine Container-Logs durch eine Pipe. Aus diesem
Lesen folgt: **an dieser Phase trägt `run-integration-tests.sh` die Form nicht**, und der Slice
behauptet keine Ursache für den roten Lauf. Die offene Frage dieses Eintrags bleibt dort.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release oder eine Versionsänderung** der Packages oder des Servers — jeder Release ist
  eine Freigabe des Auftraggebers; der Slice ändert Skripte unter `tools/`, keinen
  Produktivcode.
- **Server-Code und SDK-Code** — Gegenstand sind Hilfsskripte; die Tests der Tiers (C#, Kotlin,
  Python) bleiben unberührt.
- **Die Routing-Phase der SDK-Tiers als Verbindungs-Härtung** — die Lücke „READY heißt nicht
  verbunden“ ist ein anderer Gegenstand (Folge-Slice
  [`slice-sdk-routing-phase-verbindung-haertung`](../open/slice-sdk-routing-phase-verbindung-haertung.md));
  dieser Slice berührt dort nur die zwei `grep -q`-Zeilen von `lib-sdk-route-fixture.sh`, und nur
  in der Form des Austauschs.
- **Andere `… | grep -q…`-Pipes mit kleinem, einmaligem Erzeuger** (`echo`, `printf`, `git
  rev-parse`, `psql -c` mit einer Ausgabezeile): ein Erzeuger, der seine gesamte Ausgabe
  schreibt, bevor `grep` liest, bekommt kein SIGPIPE. Der Slice liest jede der 114 Zeilen auf
  „Erzeuger schreibt mehr als einen Puffer oder schreibt nach dem Treffer weiter“ und ändert nur
  die, auf die das zutrifft (§3 Tabelle); die übrigen bleiben **mit Begründung** stehen. Eine
  Pauschal-Umschreibung wäre ein Eingriff ohne Befund.
- **Ein neues Gate oder ein Linter** — Gate-Aufnahme verlangt eine ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6/§4); die Frage nach einem Form-Wächter (`git
  grep`) ist eine Architect-Frage, die der Slice **stellt**, nicht beantwortet (§7).
- **`make bench` und andere Runner** außerhalb der sechs Dateien — kein Befund dort (Suchlauf).

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — Reproduktion und Fix der Form `docker logs … | grep -q…`.** Vor dem Fix
      reproduziert der Implementer die Fehlerrate in einem Wegwerf-Skript im Scratchpad (nicht im
      Repo): ein beendeter Container mit einer Treffer-Zeile und 20 Folgezeilen, 1000 Aufrufe der
      Schleifenform des Bestands (`docker logs "$c" 2>/dev/null | grep -qF "$marker"` unter
      `set -euo pipefail`). Erwartet (**hergeleitet** aus der Messung des Reviewers, 1 % bei 300
      Aufrufen): einstellige Zahl bis zweistellige Zahl Fehlschläge; die gemessene Zahl und das
      Skript stehen im Bericht. Danach ersetzt er in allen 35 Fundstellen der sechs Dateien die
      Pipe-Form durch die gewählte Form (`-q` entfällt, Ausgabe nach `/dev/null`) und fährt
      dieselbe Schleife mit der Fix-Form 1000-mal: erwartet **0** Fehlschläge. *Zu belegen
      durch:* beide gedruckten Zahlen im Bericht; das Wegwerf-Skript und seine Ausgabe, mit der
      Docker-Version. (Eine Null-Zahl bei 300 Aufrufen wäre bei einer wahren Rate von 1 % mit
      Wahrscheinlichkeit 0,99³⁰⁰ ≈ 5 % zufällig, **abgeleitet**; deshalb 1000 Aufrufe, ≈ 4·10⁻⁵.)
- [x] **Liefer-Punkt 2 — Belege am Fix, Mutationsprobe, Auswahl der übrigen Pipes.**
      (i) **Mutationsprobe an einer Kopie im Scratchpad** ([`AGENTS.md`](../../../../AGENTS.md)
      §3.1: Edit/Write, nie `sed -i`): die Fix-Form einer Stelle kehrt zu `-qF` zurück; dieselbe
      1000er-Schleife gegen die Zeile der Kopie zeigt wieder Fehlschläge (rot) — die Schleife
      unterscheidet also. (ii) **Auswahl der übrigen 79 Zeilen** (114 minus 35) der weiteren
      `| grep -q…`-Pipes: jede Zeile in §3 mit Erzeuger und Entscheid „ändern“ oder „bleibt“ und
      Grund. (iii) **Reale Läufe nach dem Fix**, seriell, nach `make image`:
      `make test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`,
      `make test-sdk-python-integration` und `make test-integration` je Exit 0, im Bericht die
      Schlusszeilen; ein Lauf belegt hier **keine Flake-Freiheit**, nur dass die Ersetzung die
      Phasen nicht bricht (die Fehlerrate trägt die Schleife aus Punkt 1, nicht die Tier-Läufe).
- [ ] **Liefer-Punkt 3 — Folgesatz im Regelwerk und Beobachtungs-Register.** Der Slice entscheidet
      (Planner bei Closure, begründet), ob der Mechanismus einen Satz in
      [`AGENTS.md`](../../../../AGENTS.md) §3.9 („Geschärft“-Absatz) oder im Implementer-Ablauf
      (`.claude/commands/implement-slice.md`) bekommt, und trägt `BEO-PGC/runner-grep-pipe-verfehlt-zeile`
      mit dem Ausgang fort (Anker: dieser Slice). *Zu belegen durch:* der Diff der geänderten Datei
      oder die Begründung der Absage in §7.
- [ ] **Nur Skripte, kein Produktivcode.** `git diff --name-only <Parent>` nennt ausschließlich
      Dateien unter `tools/harness/` (die sechs genannten und, soweit Punkt 2 (ii) „ändern“
      entscheidet, die in §3 benannten), `docs/plan/` und `docs/reviews/`; kein `internal/`,
      kein `cmd/`, kein `sdks/`, keine Versionsdatei.
- [ ] **GitHub-Actions-Beleg** ([`AGENTS.md`](../../../../AGENTS.md) §3.10): der Workflow `e2e`
      führt `run-integration-tests.sh` auf dem Runner aus; der Slice ändert den Workflow nicht,
      aber den Runner. Der erste reale Post-Push-Lauf von `e2e` nach dem Merge wird gelesen
      (`gh run list --workflow e2e.yml`) und das Ergebnis im Bericht genannt; rot ohne Bezug zur
      Änderung wird gemeldet, nicht verschwiegen. Das Risiko bleibt bis dahin *weiter offen* (§6).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-harness-grep-pipe-sigpipe-unter-pipefail.md` liegt vor, kein
      offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes je
      Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-harness-grep-pipe-sigpipe-unter-pipefail.md`
      endet mit Exit 0.
- [x] Doku-Update: Kopf-Kommentare der geänderten Skripte tragen den Ist-Umfang;
      `harness/README.md` §Sensors und das Benutzerhandbuch bleiben unberührt (kein neues Target,
      keine Nutzerfläche) — am Diff zu belegen (Suchlauf).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Wellen-Betrieb für diesen wellenlosen Slice hier geprüft.

**Umfang:** S bis M — Schätzung, nicht gemessen: 35 mechanische Ersetzungen in sechs Dateien,
dazu die Lesung von 79 weiteren Zeilen und die Läufe (vier reale Läufe zu je mehreren Minuten,
zwei 1000er-Schleifen).

**Voraussetzung:** `make image` ist ausgeführt (die Tier-Läufe und `make test-integration`
setzen das geladene `:dev`-Image voraus, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md));
kein anderer Runner und kein `make bench` läuft gleichzeitig (feste Container-Namen).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | 21 Fundstellen `docker logs … \| grep -q…` (Zeilen 1783 bis 3351); darunter die nicht in einer Schleife stehenden Prüfungen (Zeilen 2570, 3181 u. a.): dort wäre ein verfehlter Treffer ein falsches Rot, kein Wiederholen. Gelesen: Einrückung null (**hergeleitet** aus der Einrückung, der Implementer liest den Kontext). |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh` | update | je 3 Fundstellen (READY, RECEIVED, `$reject_marker`); die Ablehnungs-Schleife ist die vom Reviewer gemessene. |
| `tools/harness/lib-sdk-filter-fixture.sh`, `tools/harness/lib-sdk-route-fixture.sh` | update | 3 bzw. 2 Fundstellen (READY, SEEN, SEEN_SECOND); die Dateien ändert auch der Folge-Slice `slice-sdk-routing-phase-verbindung-haertung` — Reihenfolge der Slices beachten (WIP-Limit 1). |
| weitere `\| grep -q…`-Pipes mit anderem Erzeuger (79 Zeilen in 15 Dateien) | prüfen | Erzeuger und Ausgabemenge lesen; „ändern“ nur bei Erzeuger mit Ausgabe über einem Puffer oder Schreiben nach dem Treffer (z. B. `docker exec … psql` mit Mehrzeilen-Ausgabe, `git log`, `ls -R`); sonst „bleibt“ mit Grund. Das Ergebnis steht als Tabelle im Bericht und im Befund-Feld unten. |
| `harness/README.md`, `AGENTS.md` | prüfen | siehe Liefer-Punkt 3; kein Target, kein Gate. Unverändert (Implementer, Liefer-Punkt 3 gehört dem Planner). |
| `docs/user/e2e-abdeckung.md` | update (Erzeugnis) | `make test-integration` schreibt die Datei neu; die Zeilen-Lokatoren der Runner-Zeilen verschieben sich um die vier neuen Kommentarzeilen im Kopf von `run-integration-tests.sh`. Kein Handschnitt. |
| Kopf-Kommentare der sechs Skripte | update | je ein Absatz mit dem Ist-Zustand der Prüfform (kein `-q`, `grep` liest bis zum Ende); Zusage/Kopplung nach `AGENTS.md` §3.7. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Empfänger einer Pipe unter
`pipefail` ist `grep` mit `-q`“; Parent ist `a420e223`, gemessen am 2026-10-02 mit
`make suchlauf-nachmessen`; die `diff`-Zeilen und die Befunde trägt der Implementer nach):**

```suchlauf
a420e223 35 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- harness tools examples test sdks
a420e223 21 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- tools/harness/run-integration-tests.sh
a420e223 9 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- tools/harness/run-sdk-csharp-integration-tests.sh tools/harness/run-sdk-kotlin-integration-tests.sh tools/harness/run-sdk-python-integration-tests.sh
a420e223 5 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- tools/harness/lib-sdk-filter-fixture.sh tools/harness/lib-sdk-route-fixture.sh
a420e223 114 -n -E '\| *grep +-[a-zA-Z]*q' -- harness tools examples test sdks Makefile .github
a420e223 67 -l -F pipefail -- harness tools examples test
a420e223 0 -n -F 'docker logs' -- harness/README.md AGENTS.md .claude .harness/skills
diff 0 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- harness tools examples test sdks
diff 0 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- tools/harness/run-integration-tests.sh
diff 0 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- tools/harness/run-sdk-csharp-integration-tests.sh tools/harness/run-sdk-kotlin-integration-tests.sh tools/harness/run-sdk-python-integration-tests.sh
diff 0 -n -E 'docker logs.*\| *grep +-[a-zA-Z]*q' -- tools/harness/lib-sdk-filter-fixture.sh tools/harness/lib-sdk-route-fixture.sh
diff 79 -n -E '\| *grep +-[a-zA-Z]*q' -- harness tools examples test sdks Makefile .github
diff 69 -l -F pipefail -- harness tools examples test
diff 0 -n -F 'docker logs' -- harness/README.md AGENTS.md .claude .harness/skills
```

| Träger | Messung am Parent (`a420e223`, 2026-10-02) | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Zeile 1: `docker logs … \| grep -q…`, ganzer Baum außer Doku | 35 Zeilen in sechs Dateien | Soll am Diff: 0. Gelesen, ob eine Zeile ein anderes Ziel als einen Treffer-Test hat (`-c`, `-o` mit Zählung). |
| Zeilen 2 bis 4: Aufteilung nach Datei | 21 · 9 · 5 | wie Zeile 1; die Summe 35 ist die Gegenprobe. |
| Zeile 5: jede `\| grep -q…`-Pipe | 114 Zeilen in 15 Dateien (davon 78 in `run-integration-tests.sh`, je 4 in den drei SDK-Runnern, 4 in `lib-sdk-filter-fixture.sh`, 3 in `lib-sdk-route-fixture.sh`) | Soll am Diff: 114 minus die Zahl der in Punkt 2 (ii) mit „ändern“ entschiedenen; der Rest steht mit Grund in der Tabelle im Bericht. |
| Zeile 6: Dateien, die `pipefail` setzen | 67 | Gegenprobe zum Umfang: die sechs Dateien setzen es (`run-integration-tests.sh` Zeile 87, `run-sdk-csharp-integration-tests.sh` Zeile 54); die Lib-Dateien erben es vom Aufrufer. |
| Zeile 7: Träger, die die Form beschreiben (Regelwerk, Agenten-Briefing, Skills) | 0 Zeilen | Soll am Diff: 0, außer Liefer-Punkt 3 entscheidet für einen Satz; dann trägt der Befund die Fundstelle. |

**Befunde und Belege des Implementers (gemessen am 2026-10-02, Docker 29.8.2 build 7fc2dff, Parent
`3bcafd41`; Skripte im Scratchpad, nicht im Repo):**

| Messung | Gedruckte Zeile |
|---|---|
| Reproduktion vor dem Fix: beendeter Container (`busybox`), Treffer-Zeile `REJECTED token-rejected: x` plus 20 Folgezeilen, `docker logs "$c" 2>/dev/null \| grep -qF "$marker"` unter `set -euo pipefail`, 1000 Aufrufe | `form=q calls=1000 failures=10` (1,0 %; ein zweiter Lauf an der Mutationskopie: 9 von 1000) |
| Fix-Form `docker logs … \| grep -F "$marker" >/dev/null`, dieselbe Schleife, 1000 Aufrufe | `form=nq calls=1000 failures=0` |
| Fix-Form aus der geänderten Datei gelesen (die Prüfzeile `READY` von `tools/harness/lib-sdk-filter-fixture.sh`, per `eval` gefahren), 1000 Aufrufe | `calls=1000 failures=0` |
| Mutationsprobe (Kopie der Fixture im Scratchpad, die `READY`-Zeile zurück auf `grep -qF`), dieselbe Schleife | `calls=1000 failures=9` (rot gesehen: die Schleife unterscheidet) |
| Grenze der `printf '%s' "$var" \| grep -q…`-Form: Variable mit Treffer in der ersten Zeile, 300 Aufrufe je Größe | 8 KiB: 0 Fehlschläge; 48 KiB: 0; 200 KiB: 300 (der Pipe-Puffer von Linux, 64 KiB, ist **übernommen**; die Schwelle liegt zwischen 48 und 200 KiB gemessen) |
| Reale Läufe nach dem Fix, seriell, `:dev` als `ghcr.io/pt9912/pg-change-feed:dev` geladen vorgefunden (`make image` nicht neu gelaufen: kein Server-Code im Diff) | `make test-integration` Exit 0, Schlusszeile `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen`; `make test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`, `make test-sdk-python-integration` je Exit 0, Schlusszeile je `… Filter-Belege (ADR-0133 Teilfrage 4) grün — drei Tabellen in zwei Schemas …` |

Ein Lauf belegt hier keine Flake-Freiheit (DoD); die Rate trägt die 1000er-Schleife. Der gesamte
Lauf-Log von `make test-integration` ist 38 214 Byte lang, die der drei Tier-Läufe 7,5 bis 8,1 KB
(**gemessen** mit `wc -c`); eine Variable aus diesen Ausgaben liegt weit unter 64 KiB.

**Die zwei Einmal-Stellen** (`run-integration-tests.sh` Parent-Zeilen 2570 und 3181, außerhalb einer
Schleife): beide sind das **Negativ** der Filter-Rundläufe („Change einer fremden Tabelle darf nicht
ankommen“, `sleep 3`, dann `docker logs … \| grep -qF "RECEIVED"` als Leck-Test, Treffer = `exit 1`).
Dort färbt der Mechanismus **nicht** rot, sondern **grün**: ein vorhandener Treffer (ein Leck des
Filters), den `grep -q` durch SIGPIPE verwirft, bliebe unbemerkt. Die Form ist wie bei den
Schleifen ausgetauscht; die Richtung der Fehlwirkung ist die entgegengesetzte (verfehlter
Treffer = verdecktes Leck). Die Stellen 2584 und 3195 (dieselbe Prüfung in der Wiederhol-Schleife)
sind Schleifen-Stellen.

**Entscheid über die übrigen 79 Zeilen** (Parent-Zeilen; **bleibt** in allen 79, **0** „ändern“).
Gemeinsamer Grund, soweit nicht anders genannt: der Erzeuger ist `printf '%s'`/`echo`/`printf '%s\n'`
auf eine **bereits eingelesene Shell-Variable** (Befehlsersetzung), kein laufender Prozess. Ein
Builtin schreibt die Variable in einem Zug; SIGPIPE kann nur entstehen, wenn die Variable den
Pipe-Puffer (64 KiB) übersteigt, was die Messung oben (Schwelle zwischen 48 und 200 KiB) und die
Größen der Läufe (38 KB für den ganzen Lauf-Log) ausschließen — **hergeleitet** aus Erzeuger und
Lauf-Log, nicht je Variable gemessen.

| Datei | Zeilen (Parent) | Anzahl | Erzeuger | Entscheid |
|---|---|---|---|---|
| `tools/harness/run-integration-tests.sh` | 328, 339, 743, 915, 1025, 1060, 1064, 1311, 1316, 1320, 1324, 1358, 1367, 1846, 1875, 1939, 1948, 1966, 2117, 2121, 2125, 2129, 2133, 2137, 2148, 2153, 2163, 2168, 2234, 2256, 2307, 2477, 2605, 2609, 2714, 2718, 2722, 2726, 2730, 2741, 2760, 2778, 2782, 2855, 2902, 2954, 3092, 3216, 3220, 3372, 3607, 3656, 3670, 3712, 3806, 3807, 3898 | 57 | `printf`/`echo` einer Variable aus `psql`, `docker exec`, den Wegwerf-Clients oder der `diagnose`-Ausgabe (je wenige Zeilen) | bleibt (gemeinsamer Grund); 3607: das letzte `grep -qF` liest aus `tail -n1`, eine Zeile |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `…-kotlin-…`, `…-python-…` | 274 · 272 · 308 | 3 | `printf '%s' "$test_output"` (Log eines Test-Containers, wenige KB) | bleibt |
| `tools/harness/lib-sdk-filter-fixture.sh`, `lib-sdk-route-fixture.sh` | 201 · 177 | 2 | `printf '%s' "$result_line"`, eine Zeile | bleibt |
| `tools/harness/run-command-guard-tests.sh` | 97, 99, 446, 447 | 4 | `printf '%s' "$out"`, Hook-JSON einer Zeile | bleibt |
| `tools/harness/run-fmt-check-tests.sh` | 38, 45, 175 | 3 | `printf` einer Variable; 175: `grep -A1` mit zwei Ausgabezeilen | bleibt |
| `tools/harness/run-dockerfile-from-tests.sh`, `run-image-mutation-tests.sh`, `run-kommentar-kennungen-tests.sh` | 53 · 72 · 78 | 3 | `printf` einer Variable | bleibt |
| `tools/harness/run-suchlauf-nachmessen-tests.sh` | 65, 102, 150, 182 | 4 | `printf` einer Variable | bleibt |
| `tools/schema/rollout.sh` | 130 | 1 | `printf` einer Variable (Ausgabe der Wache, wenige Zeilen) | bleibt |
| `sdks/kotlin/Dockerfile` | 92 | 1 | `printf` des Gradle-`--dry-run`-Plans; `sh` ohne `pipefail` | bleibt |
| `examples/compose.yaml` | 60 | 1 | `wget -q -O - …/healthz`, Antwort `ok` | bleibt |

Summe 57 + 3 + 2 + 4 + 3 + 3 + 4 + 1 + 1 + 1 = 79 (abgeleitet), gleich der Messung `diff 79` im
Suchlauf-Block. **Nicht Gegenstand, aber gelesen und gemeldet:** `tools/harness/pin-stale-actions.sh`
(Zeile 43), `pin-stale-baseline.sh` (26) und `pin-stale-dcheck.sh` (31) lesen
`github_api_get … \| grep -m1 '"tag_name"'` unter `set -uo pipefail` — derselbe Mechanismus
(`grep` endet beim ersten Treffer, der Erzeuger schreibt weiter), mit `-m1` statt `-q`, deshalb
von der Suche dieses Slice nicht erfasst; `make pin-stale-*` ist kein Gate. Die Betroffenheit
ist **nicht gemessen**; Meldung an den Planner (Frist: Closure dieses Slice).

## 4. Trigger

**Start** (`open` → `next` → `in-progress`): kein anderer Slice liegt in `in-progress/`
(WIP-Limit 1); `make image` ist ausgeführt; Docker mit Netzzugang für die Tier-Bauten; kein
anderer Runner und kein `make bench` laufen. Der Slice braucht den Vorgänger
`sdk-sse-filter-phase-verbindung-haertung` in `done/` (er ändert `lib-sdk-filter-fixture.sh`);
der Folge-Slice `slice-sdk-routing-phase-verbindung-haertung` startet **nach** diesem (beide
ändern `lib-sdk-route-fixture.sh`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Lesung der 79 weiteren Pipes
  ergibt mehr als zehn „ändern“-Entscheide — sie trennen sich als eigener Slice ab; die 35
  Fundstellen bleiben hier.
- `in-progress` → `open` (blockiert): die Reproduktion zeigt **keine** Fehlschläge in 1000
  Aufrufen (der Mechanismus der Herleitung trägt nicht, oder die Docker-Version verhält sich
  anders) — dann gilt der Befund des Reviewers als nicht reproduziert, der Slice hält an, der
  Fund geht an den Planner/Architect (andere Ursache zu suchen); eine Änderung ohne Befund
  ist kein Fix.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
beide 1000er-Zahlen gedruckt (vor dem Fix > 0, nach dem Fix 0), Mutationsprobe rot gesehen,
die vier realen Läufe grün, der erste Post-Push-Lauf von `e2e` gelesen, Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Fehlerrate ist an der Reproduktion nicht zu sehen.** Dann trägt der Mechanismus nicht
  (siehe Rückführung). — **Ausgang:** *zu entscheiden bei Closure*.
- **Das Ersetzen ändert eine Schleife still.** Eine Stelle, die `-q` für Exit-Status-Zwecke braucht
  (`grep -q` mit anschließendem `|| …` in einer Bedingung), behält mit `>/dev/null` dieselbe
  Semantik; eine Stelle, die `-q` mit `-c`/`-o` kombiniert, tut es nicht. — **Ausgang:** *zu
  entscheiden bei Closure*; Gegenmaßnahme: jede Zeile gelesen (Suchlauf-Befund), vier reale Läufe.
- **GitHub-Actions-Lauf unbewiesen** (`e2e` führt den Runner auf dem Runner aus,
  [`AGENTS.md`](../../../../AGENTS.md) §3.10). — **Ausgang:** *weiter offen* bis der erste
  Post-Push-Lauf gelesen ist; dann *entfallen* oder *eingetreten*.
- **Die Aussage „1 %“ ist die eines Reviewers.** Sie ist hier **übernommen**, nicht nachgemessen;
  die Zahl der 1000er-Schleife des Implementers ersetzt sie in der Closure-Notiz. — **Ausgang:**
  *zu entscheiden bei Closure*.
- **Überschneidung mit `slice-sdk-routing-phase-verbindung-haertung`** in
  `lib-sdk-route-fixture.sh`. — **Ausgang:** *entfallen* durch die Reihenfolge der Slices (§4);
  bei gleichzeitigem Bedarf Merge-Konflikt, kein stilles Überschreiben.
- **Der rote `e2e`-Lauf in der Routing-Abhilfe-Phase** (`BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot`)
  ist nicht Gegenstand. — **Ausgang:** *weiter offen* → Register (dort); der Slice behauptet keine
  Ursache und schließt keine aus.
- **Kein Release, keine Versionsänderung.** — **Ausgang:** *zu entscheiden bei Closure*;
  *entfallen*, solange `git diff` nur Skripte und Planungs-Dokumente zeigt.

## 7. Closure-Notiz

Wird bei Closure gefüllt (Vorlage: Closure-Notiz des Vorgängers
`slice-sdk-sse-filter-phase-verbindung-haertung`). Mindestinhalt: die gedruckten Zahlen der
1000er-Schleifen (Ursprung gemessen, mit Docker-Version), die Tabelle der 79 Entscheide,
die Fundstelle des Satzes in Liefer-Punkt 3 oder dessen Absage, die Antwort auf die
Architect-Frage „Form-Wächter für `… | grep -q…` unter `pipefail`?“ (gestellt, nicht
beantwortet), `BEO-PGC/runner-grep-pipe-verfehlt-zeile` mit dem Zählerstand, den der Slice
vorfindet (**1×** am Anlegen, gemessen), und der Ausgang des Post-Push-Laufs von `e2e`.
Lerneintrag-Richtung (vom Implementer zu bestätigen): ein `grep -q` als Empfänger einer Pipe
unter `pipefail` ist eine Fehlerquelle, sobald der Erzeuger nach dem Treffer weiterschreibt —
ein Hilfsskript, das ein Gate speist, liest die Ausgabe zu Ende.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem Pfad
`tools/harness/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Treffer:
`BEO-PGC/runner-grep-pipe-verfehlt-zeile` (neu mit der Closure des Vorgängers angelegt,
Zähler **1×**), `BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot` (**1×**, unbewiesener
Zusammenhang, siehe §1), `BEO-PGC/pipe-maskiert-make-exit-code` (verkörperte Nachbarregel
`AGENTS.md` §3.9, **4×**, Zähler dort nachzulesen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
