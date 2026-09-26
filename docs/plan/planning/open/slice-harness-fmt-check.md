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

**Verantwortlich:** — (noch nicht priorisiert).

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

- [ ] **Liefer-Punkt 1 — das Werkzeug.** `make fmt-check` (Aufruf über
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
- [ ] **Liefer-Punkt 2 — die Bestandsdateien.** Die sechs Dateien aus §3 sind
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
- [ ] **Liefer-Punkt 3 — die Träger.** [`.claude/commands/implement-slice.md`](../../../../.claude/commands/implement-slice.md)
      Schritt 18, Absatz „Format“: der Aufruf ist `make fmt-check` statt des
      Docker-Befehls; der Satz „Bestandsdateien außerhalb des Diffs bleiben
      unberührt“ entfällt (der Bestand ist formatiert), die Regel „nach der
      Ausgabe von `gofmt -d` korrigieren, nie mit einem Textwerkzeug“ bleibt.
      `.harness/skills/reviewer.md` LOW-Zeile: die Probe ist der Lauf von
      `make fmt-check`. *Zu belegen durch:* Lesen der beiden Stellen und der
      Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=<Plan-Datei>` läuft nach jeder Fixrunde
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors
      (Liefer-Punkt 1); das Benutzerhandbuch bleibt unberührt (keine
      Betreiber-Oberfläche).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang
      von `BEO-PGC/formatierungs-drift-ohne-gate` (`state.md`: Werkzeug
      geliefert, Adresse aufgelöst; Trigger der Gate-Aufnahme bleibt stehen) und
      eine weitere `evidence/`-Datei, falls ein Auftreten anfällt; kein Anfall
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
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
| `internal/adapters/driven/postgresstorage/queries/queries.go` | update | ein Doc-Kommentar: `''` (leerer Text) wird von `gofmt` zu einem typografischen Anführungszeichen — der Kommentar wird umformuliert. |
| `internal/adapters/driving/replication/mapper/transformation_test.go` | update | Ausrichtung eines Feldes in einem Literal. |
| `internal/adapters/driving/replication/receive/seam_test.go` | update | Ausrichtung eines Map-Eintrags. |
| `internal/application/port/outbound/log_test.go` | update | vier einzeilige Methoden werden auf je drei Zeilen umgebrochen (+8 Zeilen). |
| `internal/application/usecase/retention/service_test.go` | update | Ausrichtung eines Feld-Kommentars. |
| `test/integration/integration_test.go` | update | ein Doc-Kommentar mit einem Backtick-Paar; die Hunk-Zeilenzahl ist unverändert. |

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
```

| Träger | Befund | Behandlung |
|---|---|---|
| `.claude/commands/implement-slice.md` Schritt 18 | drei Zeilen mit `gofmt` (Absatz „Format“), ein Satz „Bestandsdateien außerhalb des Diffs bleiben unberührt“ (Zeile 1 und 2 des Feldes) | wird auf `make fmt-check` umgeschrieben (Liefer-Punkt 3); der Bestands-Satz entfällt. |
| `.harness/skills/reviewer.md` LOW-Zeile | eine Zeile mit `gofmt -l` (im Absatz **LOW**, Satz „eine Go-Datei des Diffs, die `gofmt -l` im gepinnten Toolchain-Image meldet“) | Probe wird `make fmt-check`. |
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
  sagt dasselbe und ist gofmt-stabil (`make fmt-check`). **Ausgang:** *(bei
  Closure)*
- **Das Werkzeug meldet grün, ohne etwas geprüft zu haben** (ein falsch
  gemounteter Pfad, ein leerer Suchraum). *Erwartet, zu belegen durch:* die
  Zählung der Go-Dateien mit Exit 2 bei null; die Mutation „Zählung entfernt“
  färbt den Fall rot (DoD 1). **Ausgang:** *(bei Closure)*
- **Der Format-Commit vermischt sich mit Inhalt.** *Erwartet, zu belegen durch:*
  ein eigener Commit; der Reviewer liest den Diff gegen die Ausgabe von
  `gofmt -d`. **Ausgang:** *(bei Closure)*
- **Die Zeilen-Lokatoren des Erzeugnisses `docs/user/e2e-abdeckung.md`
  verschieben sich.** *Erwartet, zu belegen durch:* die Hunk-Zeilenzahl in
  `integration_test.go` bleibt gleich (`@@ -1236,7 +1236,7 @@`, gemessen
  am Stand `7b70b34a`); am Start neu messen; verschöbe sie sich, ist die
  Meldung der Träger (§3). **Ausgang:** *(bei Closure)*
- **Ein Auftreten trotz gelaufenem Werkzeug** — der Trigger der Gate-Aufnahme
  (`BEO-PGC/formatierungs-drift-ohne-gate`, 3×). *Erwartet:* nicht in diesem
  Slice belegbar; der Ausgang ist eine Kenntnis mit benanntem Trigger.
  **Ausgang:** *(bei Closure)*
- **Das Werkzeug wird als Gate gelesen.** *Erwartet, zu belegen durch:* die Zeile
  in `harness/README.md` trägt „kein Gate“, das Ziel steht in keinem
  Gate-Bündel (`make gates` unverändert). **Ausgang:** *(bei Closure)*
- **Host-Werkzeuge außerhalb von Docker und `make`** (`bash` im Aufrufer;
  `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`, 1×).
  *Erwartet, zu belegen durch:* der Vertrag nennt sie (dieselbe Klasse wie
  `make suchlauf-nachmessen`); ein zweites Auftreten ist die Architect-Frage zu
  [`AGENTS.md`](../../../../AGENTS.md) §3.1. **Ausgang:** *(bei Closure)*
- **Ein in-place schreibendes Textwerkzeug bei der Formatierung**
  (`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, 3×). *Erwartet, zu
  belegen durch:* die sechs Hunks folgen der Ausgabe von `gofmt -d` mit dem
  Edit-Werkzeug; das Werkzeug selbst mountet lesend. **Ausgang:** *(bei
  Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu tragen bei Closure — erwartet: neuer Sensor
  `Makefile:fmt-check` ohne Gate, Schritt 18 und Reviewer-Probe nachgezogen;
  Auslöser `BEO-PGC/formatierungs-drift-ohne-gate`; das Feld `liegt in` steht nur,
  wenn wirklich verkörpert; ohne Eintrag kein `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure —
  `state.md` von `BEO-PGC/formatierungs-drift-ohne-gate`: Werkzeug geliefert;
  der Trigger der Gate-Aufnahme bleibt)*
- **Folge-Slices:** keine erwartet; der Trigger der Gate-Aufnahme ist eine Kenntnis
  (§1), kein Slice.
- **Risiken aus §6:** *(je ein Ausgang, zu tragen bei Closure)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt
  die drei Paarungen (Anker · Folge-Slice · Register), nach dem `git mv` nach
  `done/`.

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
