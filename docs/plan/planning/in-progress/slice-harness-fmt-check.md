# Slice harness-fmt-check: `make fmt-check` — Format-Werkzeug ohne Gate und Formatierung der sechs Bestandsdateien

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — Harness-Querschnitt: der Slice trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre. Er geht
`slice-transformationen-e2e-wirkung` voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5).

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.1 (Docker-only), §3.2 (kein
Linter, kein `//nolint`), §3.6 (ein Gate braucht eine ADR) und §3.9;
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft
der Zahlen im Bericht); Entscheidung des Architect-Zugs im Beobachtungs-Register:
`BEO-PGC/formatierungs-drift-ohne-gate` (`state.md`, Ausgang „Werkzeug ohne
Gate“).

**Berührte Spec-Stellen:** — (Harness-Werkzeug; keine Spec-Stelle).

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent, Auftrag des Auftraggebers (Architect-Zug zu
`BEO-PGC/formatierungs-drift-ohne-gate`). **Datum:** 2026-09-26.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make fmt-check` meldet jede Go-Datei des Baums, die `gofmt -l` im
gepinnten Toolchain-Image nicht als formatiert führt (Docker-only,
`--network none`, Exit ≠ 0 bei Abweichung, kein Gate); die sechs
Bestandsdateien sind formatiert, und Schritt 18 des Implementer-Ablaufs ruft das
Ziel statt des Docker-Befehls.

**Ausgangslage — gemessen, Stand `7b70b34a`, Befehle im Suchlauf-Feld (§3) und
in §3.** `git ls-files '*.go'` nennt 252 Dateien; `gofmt -l` über sie und
`gofmt -l .` im Toolchain-Image melden dieselben sechs Dateien (Liste in §3).
`gofmt -d` druckt sechs Hunks (86 Zeilen), fünf mit gleicher Zeilenzahl, einer
(`internal/application/port/outbound/log_test.go`) mit acht Zeilen mehr.
`gofmt -l` endet auch bei Abweichung mit Exit 0: das Ziel wertet die Ausgabe
aus, nicht den Exit-Code des Formatierers.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Aufnahme in `make gates`.** Ein Gate braucht eine ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6, §4). Die Entscheidung des
  Architect-Zugs (`BEO-PGC/formatierungs-drift-ohne-gate`, `state.md`): der
  Fund kam in allen drei Vorgängen von einem Leser vor dem Merge, Schwere LOW,
  je eine kleine Zahl Dateien — die kleinste tragende Form ist Schritt plus
  Werkzeug. **Kenntnis, nicht Umfang:** Trigger für die Aufnahme als Gate ist ein
  weiteres Auftreten, das der Reviewer trotz gelaufenem Schritt findet; dann
  Architect-Frage mit ADR-Vorschlag.
- **Ein Linter oder ein `golangci`-Profil.** Das Repo führt keinen Linter
  ([`AGENTS.md`](../../../../AGENTS.md) §3.2); `gofmt` ist der Formatierer, kein
  Lint-Profil.
- **Ein schreibender Pfad im Werkzeug** (`gofmt -w`, `--user`-Workaround). Das
  Ziel mountet lesend; die sechs Bestandsdateien folgen der Ausgabe von
  `gofmt -d` per Edit-Werkzeug, nie mit einem in-place schreibenden
  Textwerkzeug.
- **Formatierung anderer Sprachen** (Markdown, YAML, Shell, C#, Kotlin,
  Python): das Repo führt für sie keinen Formatierer im Werkzeugsatz.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — das Werkzeug.** `make fmt-check` (Aufruf über
      `tools/harness/fmt-check.sh` mit optionalem Verzeichnis-Argument, damit der
      Test ein Wegwerf-Verzeichnis nimmt) führt `gofmt -l` über alle Go-Dateien
      unter der Wurzel im gepinnten `TOOLCHAIN_IMAGE` aus
      (`docker run --rm --network none`, Repo lesend gemountet), druckt die
      abweichenden Dateien und endet mit Exit 1 bei mindestens einer, mit 0
      sonst und mit 2 bei einem Formatierer-Fehler (Syntaxfehler in einer Datei,
      Docker-Fehler) oder bei einem Verzeichnis ohne Go-Datei (leer ist nicht
      bestanden: ein falsch gemounteter Pfad meldete sonst grün). Nicht in
      `make gates`. *Zu belegen durch:* `make test-fmt-check` (Tabellentest
      gegen ein Wegwerf-Verzeichnis, netzlos) mit den Fällen formatierte Datei ·
      unformatierte Datei (genannt, Exit 1) · Datei in einem Unterverzeichnis ·
      Datei mit Syntaxfehler (Exit 2) · Verzeichnis ohne Go-Datei (Exit 2), je an
      seine Eingabe gebunden: die Mutation „nur der Exit-Code des Formatierers
      gilt“ (`gofmt -l` endet bei Abweichung mit Exit 0) färbt den Fall
      „unformatiert“ rot, „nur die oberste Ebene gelesen“ den
      Unterverzeichnis-Fall, „Zählung der Go-Dateien entfernt“ den Fall ohne
      Go-Datei (Mutationen und gesehenes Rot im Bericht). Der Vertrag
      `harness/sensors/fmt-check.md` (Form, Exit-Codes, Grenze: `gofmt` formt
      auch Doc-Kommentare um — ein Paar `''` oder ein Backtick-Paar in einem
      Doc-Kommentar wird zu einem typografischen Anführungszeichen, wer es
      anders meint, schreibt es anders; kein Import-Sortieren jenseits `gofmt`,
      kein Lint; kein Gate, Trigger der Gate-Aufnahme als Kenntnis) und eine
      Zeile in [`harness/README.md`](../../../../harness/README.md) §Sensors,
      Werkzeuge.
- [x] **Liefer-Punkt 2 — die Bestandsdateien.** Die sechs Dateien aus §3 sind
      formatiert, `make fmt-check` endet am Baum mit Exit 0. Der Format-Commit
      ist ein eigener Commit ohne Verhaltensänderung: vier Hunks sind
      Layout-Änderungen (Ausrichtung, Zeilenumbruch), **zwei ändern Zeichen in
      einem Doc-Kommentar** (`queries.go`, `integration_test.go`: `gofmt -d`
      ersetzt ein Zeichenpaar durch ein typografisches Anführungszeichen) —
      dort wird der Kommentar so umformuliert, dass `gofmt` ihn nicht mehr
      ändert und er dasselbe sagt (`''` ist der leere Text, das Backtick-Paar in
      `…003`, `/004`, `...005` ist ein Tippfehler des Kommentars). Angewandt mit
      dem Edit-Werkzeug nach der Ausgabe von `gofmt -d`, nie mit einem
      Textwerkzeug. *Zu belegen durch:* der Lauf von `make fmt-check` (Exit 0),
      `make test` grün, und der Diff des Commits gegen die Ausgabe von `gofmt -d`
      am Start (die Hunks am Start neu gemessen: 6 Dateien, 6 Hunks — die
      Zahlen aus §1 gelten für Stand `7b70b34a`).
- [x] **Liefer-Punkt 3 — die Träger.** [`.claude/commands/implement-slice.md`](../../../../.claude/commands/implement-slice.md)
      Schritt 18, Absatz „Format“: der Aufruf ist `make fmt-check` statt des
      Docker-Befehls; der Satz „Bestandsdateien außerhalb des Diffs bleiben
      unberührt“ entfällt (der Bestand ist formatiert), die Regel „nach der
      Ausgabe von `gofmt -d` korrigieren, nie mit einem Textwerkzeug“ bleibt.
      `.harness/skills/reviewer.md` LOW-Zeile: die Probe ist der Lauf von
      `make fmt-check`. *Zu belegen durch:* Lesen der beiden Stellen und der
      Suchlauf in §3.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=<Plan-Datei>` läuft nach jeder Fixrunde
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Liefer-Punkt 1); das Benutzerhandbuch bleibt unberührt (keine
      Betreiber-Oberfläche).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang
      von `BEO-PGC/formatierungs-drift-ohne-gate` (`state.md`: Werkzeug
      geliefert, Adresse aufgelöst; Trigger der Gate-Aufnahme bleibt stehen) und
      eine weitere `evidence/`-Datei, falls ein Auftreten anfällt; kein Anfall
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** S — Schätzung, nicht gemessen: ein Skript, ein Makefile-Ziel, ein
Tabellentest, ein Vertragsdokument, sechs Hunks in sechs Dateien.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/fmt-check.sh` | neu | Aufrufer: `docker run --rm --network none` mit dem Repo (oder dem Verzeichnis-Argument) lesend gemountet, `gofmt -l`, Auswertung der Ausgabe, Zählung der Go-Dateien, Exit-Codes. |
| `tools/harness/run-fmt-check-tests.sh` | neu | Tabellentest gegen ein Wegwerf-Verzeichnis (Vorbild `tools/harness/run-suchlauf-nachmessen-tests.sh`). |
| `Makefile` | update | Ziele `fmt-check` und `test-fmt-check` mit `.PHONY` und Hilfe-Zeile; kein Eintrag in den Gate-Zielen. |
| `harness/sensors/fmt-check.md` | neu | Vertrag, Exit-Codes, Grenze, Kenntnis zum Trigger der Gate-Aufnahme. |
| `harness/README.md` | update | Zeilen unter Werkzeuge für `make fmt-check` und `make test-fmt-check`. |
| `.claude/commands/implement-slice.md` Schritt 18 | update | Absatz „Format“ auf `make fmt-check`. |
| `.harness/skills/reviewer.md` | update | LOW-Zeile: die Probe ist der Lauf von `make fmt-check`. |
| `harness/sensors/fmt-check.md` | Nachzug | trägt zusätzlich den Aufruf von `gofmt -d` (lesend, im Toolchain-Image): Schritt 18 verweist für die Korrekturvorlage auf den Vertrag statt den Docker-Befehl zu wiederholen; der Vertrag nennt je Zusage die Mutation der Eingabeseite (Tabelle §Test), auch für die drei Eingabefehler des Aufrufers. |
| Schritt-18-Absatz „Format“ | Nachzug | der Verweis auf den Vertrag `harness/sensors/fmt-check.md` steht im Absatz (Rang-Zeiger statt Wiederholung). |
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | ein Doc-Kommentar: `''` (leerer Text) wird von `gofmt` zu einem typografischen Anführungszeichen — der Kommentar wird umformuliert. **Nachzug:** der Block trägt zwei Kennungen (`ADR-0124`, `ADR-0014`), `make kommentar-kennungen DIFF=<Basis>` meldet ihn als Kandidat; `ADR-0014` wird durch die Stelle ersetzt (`RetentionPolicy` in `internal/domain/model/retention.go`, [`AGENTS.md`](../../../../AGENTS.md) §3.7 Kopplung statt Kennungsreihe). |
| `internal/adapters/driving/replication/mapper/transformation_test.go` | update | Ausrichtung eines Feldes in einem Literal. |
| `internal/adapters/driving/replication/receive/seam_test.go` | update | Ausrichtung eines Map-Eintrags. |
| `internal/application/port/outbound/log_test.go` | update | vier einzeilige Methoden werden auf je drei Zeilen umgebrochen (+8 Zeilen). |
| `internal/application/usecase/retention/service_test.go` | update | Ausrichtung eines Feld-Kommentars. |
| `test/integration/integration_test.go` | update | ein Doc-Kommentar mit einem Backtick-Paar; die Hunk-Zeilenzahl ist unverändert. |
| `tools/harness/fmt-check.sh` | Fixrunde | Review F-1: der Kommentar zu „leer ist nicht bestanden“ nennt die Zusage (Exit 2 ohne Go-Datei, Exit 0 sagt: mindestens eine Datei geprüft) statt der verworfenen Alternative; F-4: die Zählung nimmt `! -type d` statt `-type f` und gleicht sich damit `gofmt` an (ein `.go`-Symlink zählt); F-6: „netzlos“ im Kopf heißt „der Container läuft ohne Netz“. |
| `tools/harness/run-fmt-check-tests.sh` | Fixrunde | F-3: Fall 11 hält die Argumente des Docker-Aufrufs mit einem Stub-`docker` fest (`--network none`, Mount `<absoluter Pfad>:/src:ro`, Image); F-4: Fall 8b bindet den `.go`-Symlink; F-6: der Kopf trennt Container ohne Netz und Image-Zugriff des Daemons im Fall Docker-Fehler. |
| `harness/sensors/fmt-check.md` | Fixrunde | F-1: „Leer ist nicht bestanden“ und Grenze 2 im Indikativ; F-3: vier Mutationszeilen (Netz, `:ro`, Image, Symlink) mit der Menge der Erprobung; F-4: Grenze 5 (`:` im Pfad, gemessen: Exit 2 mit „too many colons“) und Grenze 6 (Symlink, Symlink auf fehlendes Ziel gemessen: Exit 2); F-6: die Aussage zum Netz im Abschnitt Test. |
| `Makefile`, `harness/README.md` | Fixrunde | F-6: die Hilfe-Zeile und die README-Zeile von `make test-fmt-check` sagen „Container ohne Netz“ statt „netzlos“; die README-Zeile nennt die zwei neuen Fälle. |

**Fixrunde — gemessene Läufe (Stand: Arbeitsbaum der Fixrunde, Exit ungepiped gesichert):**
`make test-fmt-check` Exit 0 („run-fmt-check-tests: alle Fälle bestanden“), `make
fmt-check` Exit 0 („fmt-check: 254 Go-Dateien geprüft, alle formatiert“), `make test`
Exit 0, `make a-check` Exit 0 („gesamt: 0 Befund(e)“), `make coverage-gate` Exit 0
(„Coverage 85.00% erfüllt Schwelle 80%“), `make kommentar-kennungen DIFF=cfaf4c5f
COUNT=1` Exit 0 (Ausgabe 0), `make gates` Exit 0, `make suchlauf-nachmessen` Exit 0
(„10 Zeilen stimmen“). Mutationen der Eingabeseite an Kopien des Aufrufers (je
Lauf `TOOL=<Kopie> make test-fmt-check`, alle Exit 2, der Vertrag §Test nennt sie
mit der Menge der Erprobung): `--network none` entfernt und `--network host` ·
`:ro` entfernt · fest eingetragenes Image · Zählung mit `-type f` — je rot am
genannten Fall. **Träger-Meldung an den Planner (Frist: Closure dieses Slice):**
`docs/plan/planning/welle-transformationen.md` Zeilen 304 bis 306 (Sätze über die
„sechs Bestandsdateien“ über einen Zeilenumbruch; Muster des Suchlaufs blind gegen
den Umbruch, Kontrolle im Suchlauf-Feld mit dem Nebenwort `Bestandsdatei` und dem
Muster `der sechs`);
der Satz bleibt nach `done/` wahr, die Meldung ist Kenntnis, keine Änderung.

**Ist-Zustand am Start (gemessen, Stand `17cb4eb3`):** `git ls-files '*.go'` nennt
**254** Dateien (der Plan nannte 252 am Stand `7b70b34a`; die Differenz sind die zwei
Go-Dateien von `slice-code-kommentare-kennungen`, `tools/harness/kommentar-kennungen/`).
`gofmt -l` über die Dateiliste und `gofmt -l .` melden dieselben sechs Dateien;
`gofmt -d` druckt sechs Hunks in 86 Zeilen — die Zahlen aus §1 gelten unverändert.
Der Hunk in `integration_test.go` ist `@@ -1236,7 +1236,7 @@`, die Zeilen-Lokatoren
von `docs/user/e2e-abdeckung.md` verschieben sich nicht (der Diff des Format-Commits
ändert dort eine Zeile, `@@ -1239 +1239 @@`).

- **Fremder Träger, gemeldet statt geändert:** `docs/user/e2e-abdeckung.md`
  trägt Zeilen-Lokatoren `test/integration/*.go:<Zeile>` (Erzeugnis von `make
  test-integration`). Der Hunk in `integration_test.go` ersetzt sieben Zeilen
  durch sieben (`@@ -1236,7 +1236,7 @@`, `gofmt -d` am Stand `7b70b34a`), die
  Lokatoren verschieben sich nicht; am Start neu messen. Verschöbe er sie, ist
  der nächste Lauf von `make test-integration` der Träger (in der Welle
  `slice-transformationen-e2e-wirkung`).

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaft: „wie prüft der
Implementer die Formatierung“ — Symbolname `gofmt`, die Beschreibung des
Bestands-Vorbehalts und das künftige Ziel `fmt-check`. Suchraum: die
Ablauf-Träger und der Werkzeug-Baum; ausgenommen sind `docs/reviews/**`, die
Records unter `done/`, `.harness/baseline/**` und das Beobachtungs-Register (es
zitiert den Befehl in seinen Beleg-Dateien; der Ausgang wird in der Closure
geschrieben). Stand ist der Parent `7b70b34a`; der Implementer ergänzt die
Zeilen mit Stand `diff`:

```suchlauf
7b70b34a 4 -E 'gofmt' -- .claude/commands .harness/skills
7b70b34a 1 -E 'Bestandsdateien außerhalb des Diffs' -- .claude .harness/skills harness AGENTS.md
7b70b34a 0 -E 'fmt-check' -- .claude .harness/skills harness AGENTS.md README.md Makefile tools
17cb4eb3 2 -E 'sechs Bestandsdateien|Bestandsdateien, die' -- docs/plan/planning/open docs/plan/planning/next
diff 3 -E 'gofmt' -- .claude/commands .harness/skills
diff 0 -E 'Bestandsdateien außerhalb des Diffs' -- .claude .harness/skills harness AGENTS.md
diff 38 -E 'fmt-check' -- .claude .harness/skills harness AGENTS.md README.md Makefile tools
diff 0 -E 'sechs Bestandsdateien|Bestandsdateien, die' -- docs/plan/planning/open docs/plan/planning/next
17cb4eb3 3 -E 'Bestandsdatei' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/welle-transformationen.md
diff 1 -E 'Bestandsdatei' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/welle-transformationen.md
17cb4eb3 2 -E 'der sechs' -- docs/plan/planning/welle-transformationen.md
diff 2 -E 'der sechs' -- docs/plan/planning/welle-transformationen.md
```

Die drei Zeilen 1 bis 3 mit Stand `7b70b34a` sind am Start des Slice nachgemessen
und stimmen (Exit 0); die Zeilen mit Stand `17cb4eb3` (Parent des Slice) und die
mit Stand `diff` (Arbeitsbaum) messen die fremden Träger, die von den „sechs
Bestandsdateien“ sprechen. Das Muster der Zeilen 4 und 8 ist zeilenweise und
trifft nur Sätze, die die Wörter in einer Zeile tragen: je ein Satz in
`slice-transformationen-e2e-wirkung` und `slice-antragsqueue-lesefehler-failed`
(am Parent zwei Treffer). Das Nebenwort `Bestandsdatei` in den Zeilen 9 und 10
trifft zusätzlich den Satz über einen Zeilenumbruch in
`docs/plan/planning/welle-transformationen.md` (Zeile 305, „Bestandsdateien“;
am Parent drei Treffer). Die Zeilen 11 und 12 messen dieselbe Datei mit dem
Muster `der sechs`, das den Zeilenumbruch nicht braucht: Zeile 304 („Formatierung
der sechs“) und Zeile 306 („eine der sechs;“), an beiden Ständen gleich — das
sind drei Träger, nicht zwei.

Die Closure zog die zwei Sätze in fremden Dateien nach (Frist erfüllt): der Satz in
`slice-transformationen-e2e-wirkung` §4 begründete eine Kante, die mit `done/` dieses
Slice erfüllt ist, und der Satz in `slice-antragsqueue-lesefehler-failed` §4 nannte
die Dateien des Bestands; beide nennen jetzt den Ist-Zustand (formatiert,
`make fmt-check` Exit 0), deshalb steht das `diff`-Soll der Zeilen 8 und 10 auf 0
und 1. Die Sätze in `welle-transformationen` (Zeilen 304 bis 306) beschreiben den
Umfang dieses Slice (Formatierung der sechs Bestandsdateien, Kante zu
`e2e-wirkung`) und bleiben nach `done/` wahr; sie tragen keine offene Kante und
bleiben unverändert. Zeile 7 steht bei 38 statt 37: die Closure trägt die
Herkunft `seit slice-harness-fmt-check` in Schritt 18 des Implementer-Ablaufs ein
(eine Zeile mehr).

| Träger | Befund | Behandlung |
|---|---|---|
| `.claude/commands/implement-slice.md` Schritt 18 | drei Zeilen mit `gofmt` (Absatz „Format“), ein Satz „Bestandsdateien außerhalb des Diffs bleiben unberührt“ (Zeile 1 und 2 des Feldes) | wird auf `make fmt-check` umgeschrieben (Liefer-Punkt 3); der Bestands-Satz entfällt. Am Ende: drei Zeilen mit `gofmt` (`make fmt-check` erklärt, `gofmt -d` als Korrekturvorlage), null Zeilen mit dem Bestands-Satz. |
| `.harness/skills/reviewer.md` LOW-Zeile | eine Zeile mit `gofmt -l` (im Absatz **LOW**, Satz „eine Go-Datei des Diffs, die `gofmt -l` im gepinnten Toolchain-Image meldet“) | Probe wird `make fmt-check` (Ende: null Zeilen mit `gofmt` in `.harness/skills`, ein Link auf den Vertrag). |
| `.claude/hooks/pretooluse-command-guard.sh` | ein Kommentar erwähnt `gofmt` als Beispiel eines Kommando-Segments (nicht Teil des Feldes: Suchraum `.claude/commands`) | unberührt. |
| Beobachtungs-Register `BEO-PGC/formatierungs-drift-ohne-gate` | Zustand „Werkzeug offen, Adresse `slice-harness-fmt-check`“ | Träger der Planner-Closure: `state.md` auf „Werkzeug geliefert“; fremde Datei, deshalb Meldung. |
| `harness/README.md` | keine Zeile zu einem Format-Ziel (Zeile 3 des Feldes: 0 Treffer am Parent) | Zeilen kommen mit Liefer-Punkt 1. |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/`
(WIP-Limit 1). **Reihenfolge (Empfehlung an den Orchestrator):** 1.
`slice-code-kommentare-kennungen`, 2. dieser Slice, 3.
`slice-antragsqueue-lesefehler-failed`, danach die Welle
[welle-transformationen](../welle-transformationen.md); eine technische Kante zu
`slice-code-kommentare-kennungen` gibt es nicht. Dieser Slice muss `done` sein,
**bevor** `slice-transformationen-e2e-wirkung` startet (Kante, dort §4):
`slice-transformationen-e2e-wirkung` und `slice-transformationen-e2e-abhilfe`
erweitern `test/integration/integration_test.go`, eine der sechs
Bestandsdateien (gelesen an ihren §3-Tabellen); ein Implementer, der sie ändert,
sähe sonst `gofmt` auf Bestand melden.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  kleiner Zug. Zeigt `gofmt -d` am Start mehr als die sechs Hunks, ist der
  abtrennbare Teil die Formatierung der Bestandsdateien (Liefer-Punkt 2) als
  eigener Slice, mit dem Werkzeug davor.
- `in-progress` → `open` (blockiert): `gofmt -l` über den Baum meldet Dateien,
  deren Format ein Generator vorgibt (erzeugter Code), oder das gepinnte Image
  trägt `gofmt` nicht mit demselben Ergebnis wie am Start gemessen — dann
  Architect-Frage zum Suchraum, kein Ausnahme-Pfad im Werkzeug.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Review-Report ohne offenes HIGH oder
MEDIUM + Verifikation, dass die DoD trägt (`make fmt-check` am Baum mit Exit 0,
`make test-fmt-check` grün) + Closure-Notiz mit Lerneintrag geschrieben (neuer
Sensor: `make fmt-check`, ohne Gate).

## 6. Risiken und offene Punkte

- **`gofmt` ändert den Inhalt eines Doc-Kommentars** (zwei der sechs Hunks: ein
  Zeichenpaar wird zu einem typografischen Anführungszeichen). *Erwartet, zu
  belegen durch:* der Reviewer liest die zwei Kommentare gegen die Stelle
  (`''` ist im SQL der leere Text, in `queries.go`); die umformulierte Fassung
  sagt dasselbe und ist gofmt-stabil (`make fmt-check`). **Ausgang:** eingetreten
  und behandelt. Zwei Hunks ändern einen Doc-Kommentar (`queries.go`,
  `integration_test.go`); Review und Verifikation lasen die umformulierten
  Kommentare gegen den Code und fanden sie wahr (`''` ist im SQL der leere Text,
  die Zeichenliste des Tests ist `{"…", "...", "/"}`), `make fmt-check` endet am
  Baum mit Exit 0 (gedruckt: „254 Go-Dateien geprüft, alle formatiert“, Lauf dieser
  Closure).
- **Das Werkzeug meldet grün, ohne etwas geprüft zu haben** (ein falsch
  gemounteter Pfad, ein leerer Suchraum). *Erwartet, zu belegen durch:* die
  Zählung der Go-Dateien mit Exit 2 bei null; die Mutation „Zählung entfernt“
  färbt den Fall rot (DoD 1). **Ausgang:** entfallen. Die Mutation „Zählung
  entfernt“ färbte „leeres Verzeichnis“ und „Verzeichnis ohne Go-Datei“ rot, im Lauf
  des Reviewers (M3) und des Verifiers (M3, gedruckt), die Zählung ist an
  `gofmt` angeglichen (`! -type d`, ein `.go`-Symlink zählt). Zwei Grenzen bleiben
  im Vertrag benannt und enden geschlossen mit Exit 2, nie mit 0: ein `:` im
  Verzeichnispfad und ein Symlink auf ein fehlendes Ziel (Grenzen 5 und 6).
- **Der Format-Commit vermischt sich mit Inhalt.** *Erwartet, zu belegen durch:*
  ein eigener Commit; der Reviewer liest den Diff gegen die Ausgabe von
  `gofmt -d`. **Ausgang:** entfallen. Der Format-Commit `32637052` ist ein eigener
  Commit; der Verifier reduzierte ihn mit `git show -w` auf zwei Kommentare und
  einen Umbruch (`log_test.go`), der Reviewer las den Diff gegen `gofmt -d` (6 Hunks
  in 86 Zeilen am Start).
- **Die Zeilen-Lokatoren des Erzeugnisses `docs/user/e2e-abdeckung.md`
  verschieben sich.** *Erwartet, zu belegen durch:* die Hunk-Zeilenzahl in
  `integration_test.go` bleibt gleich (`@@ -1236,7 +1236,7 @@`, gemessen
  am Stand `7b70b34a`); am Start neu messen; verschöbe sie sich, ist die
  Meldung der Träger (§3). **Ausgang:** entfallen. Der Hunk ist `@@ -1239 +1239
  @@` (eine Zeile ersetzt eine Zeile), `git diff 17cb4eb3..HEAD --
  docs/user/e2e-abdeckung.md` ist leer (Verifikation §3).
- **Ein Auftreten trotz gelaufenem Werkzeug** — der Trigger der Gate-Aufnahme
  (`BEO-PGC/formatierungs-drift-ohne-gate`, 3×). *Erwartet:* nicht in diesem
  Slice belegbar; der Ausgang ist eine Kenntnis mit benanntem Trigger.
  **Ausgang:** weiter offen als Kenntnis. Im Slice nicht belegbar: Review und
  Verifikation fanden kein Auftreten trotz gelaufenem Werkzeug; der Trigger der
  Gate-Aufnahme steht unverändert im Vertrag (§Kein Gate) und im `state.md` des
  Register-Eintrags.
- **Das Werkzeug wird als Gate gelesen.** *Erwartet, zu belegen durch:* die Zeile
  in `harness/README.md` trägt „kein Gate“, das Ziel steht in keinem
  Gate-Bündel (`make gates` unverändert). **Ausgang:** entfallen. Beide
  README-Zeilen und der Vertrag tragen „kein Gate“; `git diff 17cb4eb3..HEAD --
  Makefile` trägt keine Zeile zu `GATE_CHECKS`, die Ausgabe von `make gates` kein
  `fmt-check` (Review und Verifikation je gemessen).
- **Host-Werkzeuge außerhalb von Docker und `make`** (`bash` im Aufrufer;
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`, 1×).
  *Erwartet, zu belegen durch:* der Vertrag nennt sie (dieselbe Klasse wie
  `make suchlauf-nachmessen`); ein zweites Auftreten ist die Architect-Frage zu
  [`AGENTS.md`](../../../../AGENTS.md) §3.1. **Ausgang:** entfallen. Der Vertrag
  nennt `bash`, `git`, `realpath` und `docker` als Host-Werkzeuge der Klasse „Host-Werkzeug
  ohne Installation“ (`AGENTS.md` §3.1, seit dem Architect-Zug `17cb4eb3`), das
  Werkzeug ruft kein Host-`gofmt`; die Architect-Frage ist beantwortet, ein
  drittes Auftreten fiel nicht an (Register-Zähler bleibt 2×).
- **Ein in-place schreibendes Textwerkzeug bei der Formatierung**
  (`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, 3×). *Erwartet, zu
  belegen durch:* die sechs Hunks folgen der Ausgabe von `gofmt -d` mit dem
  Edit-Werkzeug; das Werkzeug selbst mountet lesend. **Ausgang:** entfallen. Der Diff
  trägt keine Kommandozeile mit `sed -i`, `perl -pi` oder `awk -i` (Review und
  Verifikation je gemessen, der Verifier nennt eine Prosa-Nennung im Review-Report);
  die Mutationsproben liefen an Kopien im Scratchpad, das Original blieb
  unverändert (`git status --short` leer).

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Leser-Kette fand, was der Implementer-Lauf nicht
  fand: der Review nennt 1 HIGH, 0 MEDIUM, 2 LOW und 3 INFO, die Verifikation 0 HIGH,
  0 MEDIUM, 1 LOW und 3 INFO (übernommen aus den Reports). Der Reviewer fuhr 17
  Mutationen am Aufrufer, 14 rot, drei grün: M11 (`:ro` allein) und M14 (`--network
  none`) waren ungebunden (F-3), M16 ist ein unerreichbarer Zweig; die Fixrunde band
  M11 und M14 mit einem Stub-`docker` (Fall 11). Der Verifier fuhr 12 Mutationen, alle
  rot (übernommen aus Review und Verifikation §4). (2) Der Format-Commit `32637052`
  ist verhaltensgleich: sechs Dateien, sechs Hunks in 86 Zeilen (am Start je von
  Implementer, Reviewer und Verifier gemessen, übernommen), der Verifier reduzierte den
  Diff mit `git show -w` auf zwei Kommentare und einen Umbruch; die Zeilen-Lokatoren
  von `docs/user/e2e-abdeckung.md` bewegten sich nicht (`git diff` leer). (3) Das
  Werkzeug endet am Baum mit Exit 0 („fmt-check: 254 Go-Dateien geprüft, alle
  formatiert“, gemessen bei dieser Closure; `git ls-files '*.go' | wc -l` gibt 254).
  (4) Das Suchlauf-Feld trug: `make suchlauf-nachmessen` meldete bei Review und
  Verifikation je 10 stimmende Zeilen (Exit 0, übernommen), der Verifier fuhr alle
  zehn von Hand nach; das Feld trägt bei dieser Closure 12 Zeilen (Lauf bei dieser
  Closure, unten).
- **Was ging anders als geplant:** (1) Die Zahl der Go-Dateien: der Plan nannte 252
  (übernommen aus der Messung am Stand `7b70b34a`), am Start und bei der Closure sind
  es 254 (gemessen); die Differenz sind die zwei Go-Dateien von
  `slice-code-kommentare-kennungen` (§3, Ist-Zustand am Start). (2) Der Block in
  `queries.go` trug zwei Kennungen; der Ersatz von `ADR-0014` durch die Stelle
  (`RetentionPolicy` in `internal/domain/model/retention.go`) war Nachzug, weil
  `make kommentar-kennungen` ihn als Kandidat meldete (§3). (3) Über den Plan hinaus
  entstanden `gofmt -d` im Vertrag (die Korrekturvorlage), vier Mutationszeilen je
  Zusage (Netz, `:ro`, Image, Symlink), Fall 8b (Symlink) und Fall 11 (Docker-Argumente
  per Stub); Grund: F-3 und F-4 des Reviews. (4) Die Fixrunde meldete „Tabelle
  wiederhergestellt“, die Tabelle in §3 blieb zerrissen; der Verifier fand es (V-1),
  diese Closure behob es (Leerzeile entfernt, sechs plus vier Zeilen stehen wieder in
  der Tabelle). (5) Die Coverage-Zahl ist an diesem Stand nicht lauf-stabil: der Plan
  nennt 85.00 % (Lauf der Fixrunde), der Verifier 85.10 % und 85.00 % in zwei Läufen
  (übernommen aus der Verifikation §1); die Streuung ist 0,10 Prozentpunkte
  (abgeleitet), die Zahl ist Beleg des jeweiligen Laufs, keine Zustandsgröße. (6) Die
  Prosa unter dem Suchlauf-Block nummerierte die Zeilen nach der Umordnung falsch (V-2);
  sie ist nachgezogen, das Feld trägt zwei Zeilen mit dem Muster `der sechs`, das den
  Zeilenumbruch in `welle-transformationen.md` nicht braucht (V-3).
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Geschärfte Regel.*
  `.claude/commands/implement-slice.md` Schritt 18 (Absatz „Format“) ruft `make
  fmt-check` statt des Docker-Befehls und verweist für die Korrekturvorlage auf den
  Vertrag; `.harness/skills/reviewer.md` trägt in der LOW-Zeile die Probe `make
  fmt-check` (der Lauf misst den ganzen Baum, ein Bestands-Vorbehalt entfällt) · seit
  slice-harness-fmt-check. Herkunft: `BEO-PGC/formatierungs-drift-ohne-gate`. *(b) Neuer
  Sensor, kein Gate.* `make fmt-check` und `make test-fmt-check` — liegen in
  `harness/sensors/fmt-check.md`, `harness/README.md` §Sensors (zwei Werkzeug-Zeilen) und
  im `Makefile` · seit slice-harness-fmt-check. Ein Gate braucht eine ADR ([`AGENTS.md`](../../../../AGENTS.md)
  §3.6); der Trigger der Aufnahme (ein Auftreten trotz gelaufenem Schritt 18) steht als
  Kenntnis im Vertrag und im Register. *(c) Benannte Lücken, je mit Adresse.* Erstens:
  `make fmt-check` misst den ganzen Baum, nicht den Diff — Exit 0 ist kein Beleg dafür,
  dass der Diff eine Formatierung trägt, die er hätte tragen müssen (Vertrag, Grenze
  „Grün ist kein Beleg für den Diff“; Adresse: der Reviewer liest den Lauf im Bericht
  nach). Zweitens: zwei Pfad-Grenzen enden geschlossen mit Exit 2 statt mit einer
  Meldung — ein `:` im Verzeichnispfad und ein Symlink auf ein fehlendes Ziel (Vertrag,
  Grenzen 5 und 6; kein Repo-Pfad betroffen — kein Symlink mit `.go`-Endung im Baum,
  übernommen aus dem Review, dort mit `git ls-files -s` gemessen). Drittens: Formatierung anderer Sprachen (Markdown, YAML, Shell, C#,
  Kotlin, Python) hat keinen Formatierer im Werkzeugsatz (§1, Nicht-Umfang; Adresse:
  keine, bis ein Fund sie verlangt). Viertens: die Gate-Frage (Trigger im `state.md` von
  `BEO-PGC/formatierungs-drift-ohne-gate`). *(d) Gelernt, bei 1× keine Regel:* eine
  Fixrunde, die ein Finding zur Form eines Dokuments behebt, meldet die Behebung nur
  glaubwürdig, wenn sie den Kontext um die geänderte Stelle nachliest (`git diff -U20`);
  `make docs-check` sieht keine zerrissene Markdown-Tabelle (V-1). Das Auftreten steht
  als Deckel-Fall unten; ein drittes ohne Widerspruch der Aussagen wäre der Anlass, die
  Form als eigenen Register-Eintrag abzuspalten. *(e) Finding-Klassen des Reviews und
  der Verifikation:* Kommentar beschreibt die verworfene Alternative (F-1) · Einfügung
  zerreißt eine Tabelle (F-2, V-1) · Isolations-Flags des Aufrufs ohne Testbindung (F-3)
  · Zählregel weicht von `gofmt` bei Symlinks ab (F-4) · Suchmuster blind gegen
  Zeilenumbruch (F-5, V-3) · „netzlos“ trägt für den Docker-Fehler-Fall nur mittelbar
  (F-6) · Prosa nummeriert Feld-Zeilen falsch (V-2) · Zahl nicht lauf-stabil (V-4); ihre
  Zuordnung zu den Register-Zählern steht im nächsten Punkt.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-harness-fmt-check.md`, Zähler = Zahl der Dateien. *Neuer Beleg:*
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` **4×** (F-1, HIGH, daher
  Datei; der erste Fund nach der Verkörperung, von der Skopus-Klausel als HIGH
  eingeordnet und vor dem Merge gefunden; Ausgang unverändert **verkörpert**). *`state.md`
  fortgeschrieben ohne neue Datei:* `BEO-PGC/formatierungs-drift-ohne-gate` (Werkzeug
  geliefert, Schritt und Werkzeug verkörpert, Trigger der Gate-Aufnahme unverändert
  als Kenntnis, Zähler 3×; der Trigger ist nicht eingetreten: Review und Verifikation
  nennen keinen Format-Befund). *Deckel-Fälle ohne Datei, Finding-Kennung hier*
  (verkörpert ab 10×, vor dem Merge von Reviewer bzw. Verifier gefunden, Schwere ≤ LOW,
  bekannter Träger-Typ): F-3 (LOW, Aufrufer eines Werkzeugs; Ausprägung: die Argumente
  des Docker-Aufrufs sind die Eingabeseite der Isolations-Zusage) zu
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (Deckel bei 14×); F-2 und V-1 (LOW,
  Einfügung in einen stehenden Träger, beide ohne Widerspruch der Aussagen; V-1: die
  Fixrunde behob F-2 nicht wirksam) zu `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
  (Deckel bei 10×); F-5 und V-3 (INFO, das Muster des Suchlaufs ist zeilenweise, die
  Sätze in `welle-transformationen.md` laufen über einen Umbruch) zu
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Deckel bei 32×); V-2 und V-4 (INFO, Prosa
  nummeriert Zeilen des Feldes falsch, Coverage-Zahl nicht lauf-stabil) zu
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Deckel bei 23×). *Kein Anfall:*
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (Review und Verifikation nennen
  keinen Fund; die Mutationen liefen an Kopien; sein Guard-Ausbau ist Gegenstand eines
  eigenen Slice, dessen Plan in `open/` angelegt wird),
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (der Vertrag nennt
  `bash`, `git`, `realpath` und `docker` nach der Klasse in [`AGENTS.md`](../../../../AGENTS.md)
  §3.1; das Werkzeug ruft kein Host-`gofmt`; Zähler bleibt 2×),
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (F-1 ist eine verworfene
  Alternative, keine nicht getragene Zusage; F-6 nennt einen Wortlaut, dessen Zusage für
  den Container wahr ist). *Lese-Schritt:* aus diesem Slice erreicht neu kein Eintrag die
  Schwelle ohne Ausgang.
- **Folge-Slices:** keine. Der Trigger der Gate-Aufnahme ist eine Kenntnis (§1) mit
  Architect-Frage als Adresse, kein Slice. Start-Bedingungen der nächsten Slices,
  gelesen: `slice-transformationen-e2e-wirkung` verlangt diesen Slice in `done/` (dort §4;
  mit dem Move erfüllt, der Satz ist auf den Ist-Zustand gezogen: die Datei
  `integration_test.go` ist formatiert, Schritt 18 lässt `make fmt-check` nach dem
  Diff grün); `slice-antragsqueue-lesefehler-failed` verlangt kein anderes Slice in
  `in-progress/` (dort §4, WIP-Limit 1; mit dem Move erfüllt, in `in-progress/` liegt
  danach nur die Roadmap), der Satz über die „sechs Bestandsdateien“ ist gezogen (Frist der
  Meldung: diese Closure, erfüllt); `welle-transformationen` (Zeilen 304 bis 306)
  beschreibt den Umfang dieses Slice und bleibt unverändert wahr.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Eingetreten:* `gofmt` ändert
  einen Doc-Kommentar (behandelt, Kommentare gegen den Code gelesen). *Entfallen:*
  Werkzeug grün ohne geprüft zu haben · Format-Commit vermischt sich mit Inhalt · Zeilen-
  Lokatoren verschieben sich · Werkzeug als Gate gelesen · Host-Werkzeuge · in-place
  schreibendes Textwerkzeug. *Weiter offen als Kenntnis:* ein Auftreten trotz gelaufenem
  Werkzeug (Trigger der Gate-Aufnahme).
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt die
  drei Paarungen, nach dem `git mv` nach `done/`. *Anker:* `.claude/commands/implement-slice.md`
  Schritt 18 trägt `seit slice-harness-fmt-check`; die zwei Zeilen `make fmt-check` und
  `make test-fmt-check` in `harness/README.md` tragen es ebenfalls; `harness/sensors/fmt-check.md`
  und die zwei `Makefile`-Ziele existieren (geprüft mit `git grep` und `ls`). *Folge-Slice:* keiner
  genannt (die Ereignis-Adresse des Gate-Triggers ist eine Architect-Frage, sie kann
  eintreten: ein Fund des Reviewers trotz gelaufenem Schritt 18). *Register:* jede genannte
  Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/`
  (geprüft mit `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); `tools/harness`, `harness/` und sechs Go-Dateien in
fünf Paketen sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
gemessen am 2026-09-26 mit `ls evidence | wc -l` je Eintrag):

- `BEO-PGC/formatierungs-drift-ohne-gate` (verkörpert als Schritt, Werkzeug offen,
  3×): der Gegenstand dieses Slice; Trigger der Gate-Aufnahme als Kenntnis (§1).
- `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (offen, Schwelle
  erreicht, 3×): die Bestandsdateien werden mit dem Edit-Werkzeug nach `gofmt -d`
  korrigiert (Risiko in §6); die Spur eines Fehlgriffs der Klasse ist eine
  `gofmt`-Abweichung im Diff.
- `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (offen, 1×):
  Risiko in §6.
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 16×): die
  Mutationen der Eingabeseite in Liefer-Punkt 1.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 14×),
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32×),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23×): das
  Suchlauf-Feld in §3 und die Zahlen mit Ursprung.
- `BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×): das Werkzeug
  mountet lesend und schreibt nichts; kein Auftreten möglich.
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
