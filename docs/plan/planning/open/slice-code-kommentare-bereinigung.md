# Slice code-kommentare-bereinigung: Bereinigung der Kennungs-Ketten in Go-Kommentaren — Nicht-Test-Code in Tranchen nach Paketen, Test-Code zuletzt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — Harness-Querschnitt: der Slice trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre. Er startet nach
`slice-code-kommentare-kennungen`, `slice-harness-fmt-check` und
`slice-transformationen-e2e-abhilfe` (Start-Trigger §4;
[welle-transformationen](../welle-transformationen.md) §5).

**Bezug:** [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen in Trägern), [`AGENTS.md`](../../../../AGENTS.md) §3.7
(Regel in der Konkretisierung von `slice-code-kommentare-kennungen`) und §3.13,
Architect-Verdikt
[`architect-verdict-slice-chronik-in-code-kommentar`](../../../reviews/architect-verdict-slice-chronik-in-code-kommentar.md)
(Testfall-Provenienz bleibt zulässig).

**Berührte Spec-Stellen:** — (nur Kommentare in Go-Dateien; keine Spec-Stelle).

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Auftrag des Auftraggebers zur Kennungsdichte in
Code-Kommentaren. **Datum:** 2026-09-26.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Go-Kommentare des Baums tragen ihre Herkunft als ein Feld: die
Kandidatenzahl von `make kommentar-kennungen` sinkt in Tranchen nach Paketen auf
0 oder auf eine begründete, benannte Restmenge mit Adresse — Nicht-Test-Code
zuerst, Test-Code als letzte Tranche —, und nur Kommentare ändern sich.

**Voraussetzung:** `slice-code-kommentare-kennungen` liefert Regel
([`AGENTS.md`](../../../../AGENTS.md) §3.7), Werkzeug und die Basis-Messung;
dieser Slice wendet sie an. Die Zahlen unten sind **gemessen** am Stand
`d13ab81e` (2026-09-26) mit `make kommentar-kennungen COUNT=1 TESTS=exclude
PATHS=<Suchraum>` (T1 bis T7) bzw. `… TESTS=only PATHS=<Suchraum>` (die Teile von
T8): 597 Kandidaten gesamt, davon 400 in Nicht-Test- und 197 in Testdateien
(`make kommentar-kennungen COUNT=1` mit `TESTS=exclude`, `TESTS=only` und ohne;
die Summen der Tranchen sind 400 und 197, abgeleitet). Sie sind eine
Zustandsgröße ([`AGENTS.md`](../../../../AGENTS.md) §3.12): die Basis-Messung
am Start ersetzt sie, und jede Tranche misst vor und nach.

**Tranchen** (je Tranche ein oder mehrere Commits; `PATHS` ist das Argument des
Werkzeugs):

| Tranche | Suchraum (`PATHS`) | Kandidaten (gemessen, Stand `d13ab81e`) |
|---|---|---|
| T1 | `internal/application/port` | 77 |
| T2 | `internal/domain` | 41 |
| T3 | `internal/application/usecase` | 36 |
| T4 | `internal/adapters/driving` | 76 |
| T5 | `internal/adapters/driven` | 91 |
| T6 | `internal/bootstrap`, `cmd` | 57 |
| T7 | `tools`, `examples` | 22 |
| T8 | alle `*_test.go`, in fünf Teil-Commits: `internal/adapters/driven` (70), `internal/bootstrap` und `cmd` (37), `internal/domain`, `internal/application` (37), `internal/adapters/driving` (32), `test/integration` (21) | 197 |

**Regeln je Kandidat.** Die Sätze, die eine Zusage, Kopplung, Abgrenzung oder
Grenze der **Stelle** tragen, bleiben; gekürzt wird die Herkunft auf **einen**
Anker (den, der die Norm trägt) und eine Wiedergabe von Spec-Inhalt in eigenen
Worten wird zum Verweis. Eine Zusage, die der Code nicht trägt, wird nicht
stillschweigend gestrichen, sondern im Bericht gemeldet (Befund für den
Reviewer, `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`). Testfall-
Provenienz („`TestXyz` trägt … aus …“) ohne Kennungs-Kette bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung an Code.** Nur Kommentarzeilen ändern sich; der Nachweis steht in
  der DoD (Diff ohne Nicht-Kommentarzeile, Ausnahme Endkommentare einzeln
  aufgezählt). Ein Fund im Code ist ein Befund, kein Mitnehmen.
- **Kommentare außerhalb von Go** — Skripte (`tools/**/*.sh`), `Makefile`,
  `harness/mk`, `.sql`, `.yml`, Dockerfiles, Markdown: das Werkzeug liest
  Go-Kommentargruppen, für die anderen Formen gibt es keine Blockgrenze und
  keine Messung. Sie sind der Folge-Slice `slice-kommentar-kennungen-skripte`
  (Datei in `open/`, die der Planner bei der Closure dieses Slice anlegt, §7);
  die Regel gilt dort ab sofort, das Werkzeug fehlt.
- **Erzeugter Code** (`gen/`) und die SDK-Bäume (`sdks/`): erzeugt bzw. durch
  `make sdk-public-doc-check` strenger gedeckt.
- **Ein Suchlauf über Slice-/Wellen-Chronik in Produktionskommentaren.** Ein
  Chronik-Befund im selben Block wird mit umformuliert; eine eigene Messung dazu
  ist `BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, Verdikt: kein
  Sensor).
- **Die Wahrheitsprüfung der Zusagen.** Der Slice kürzt Herkunft; ob ein Satz
  zutrifft, prüft der Reviewer je Tranche (Lese-Handlung).
- **Ein Gate.** Siehe `slice-code-kommentare-kennungen` §1: die Frage stellt sich
  nach diesem Slice, mit ADR.

## 2. Definition of Done

- [ ] **Liefer-Punkt 1 — Nicht-Test-Code (T1 bis T7).** Je Tranche steht die
      Kandidatenzahl vor und nach der Änderung im Bericht, gemessen mit
      `make kommentar-kennungen COUNT=1 TESTS=exclude PATHS=<Suchraum der
      Tranche>` (Befehl, Lauf und Zahl, Ursprung „gemessen“,
      [`AGENTS.md`](../../../../AGENTS.md) §3.12); nach T7 endet
      `make kommentar-kennungen COUNT=1 TESTS=exclude` mit 0 oder mit der
      benannten Restmenge (Liefer-Punkt 3). Nur Kommentare ändern sich. *Zu
      belegen durch:* je Tranche der Diff (`git diff -U0 <Basis der Tranche> --
      '*.go'`) ohne Zeile, die nicht mit `//` nach Leerraum beginnt — jede
      Endkommentar-Zeile (Code mit `// …` in derselben Zeile) einzeln
      aufgezählt, ihr Code-Teil unverändert; `make test` grün nach jeder Tranche
      (Exit direkt, [`AGENTS.md`](../../../../AGENTS.md) §3.9); `make fmt-check`
      ohne Ausgabe (gofmt formt Doc-Kommentare um: ein Kommentar, der es
      ändert, wird so umformuliert, dass er stabil bleibt); `make a-check`
      grün; die Zusagen der Kommentare bleiben (der Reviewer liest je Tranche
      Stichproben gegen den Code, Zahl der Stichproben im Bericht).
- [ ] **Liefer-Punkt 2 — Test-Code (T8).** Dasselbe Nachweis-Paket mit
      `TESTS=only` in fünf Teil-Commits; der Teil `test/integration` ist der
      letzte und schließt mit einem realen, grünen `make test-integration`-Lauf,
      der die Datei [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md)
      neu schreibt (Erzeugnis: sie trägt Zeilen-Lokatoren auf
      `test/integration/*.go`, die Kommentar-Kürzungen verschieben). *Zu belegen
      durch:* je Teil-Commit die Zahlen vor und nach und der Diff-Nachweis wie in
      Liefer-Punkt 1; am Ende `make kommentar-kennungen COUNT=1` (gesamt) mit
      0 oder der Restmenge; der `make test-integration`-Lauf mit seiner
      gedruckten Zeile; `git diff` der Erzeugnis-Datei zeigt nur verschobene
      Lokatoren.
- [ ] **Liefer-Punkt 3 — Restmenge und Träger.** Jeder nach T8 verbleibende
      Kandidat steht im Bericht und in §7 mit `Datei:Zeile`, Klasse (Zusage ·
      Kopplung · Abgrenzung · Rang-Zeiger · Grenze) und Grund; die Gesamtzahl
      des Werkzeugs ist 0 oder gleich der Restmenge; ein Kandidat mit zwei
      Ankern, den der Implementer als konform begründet, ist ein Fall für die
      Regel (Meldung an den Planner: Konkretisierung in
      [`AGENTS.md`](../../../../AGENTS.md) §3.7), kein Ausnahme-Eintrag. Träger,
      die die Bewegung beschreiben (§3-Suchlauf), sind nachgezogen oder mit
      Adresse gemeldet. *Zu belegen durch:* der Lauf des Werkzeugs, die Liste im
      Bericht und das Suchlauf-Feld in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — je einer nach T7 und nach T8 (zwei
      Reports; ein Report über den ganzen Diff trüge ihn nicht, §4);
      Rollenwechsel nach Schritt 8 des Minimal Agent Workflow
      ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=<Plan-Datei>` läuft nach jeder Fixrunde
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update entfällt: nur Kommentare ändern sich; das Benutzerhandbuch
      bleibt unberührt (keine Betreiber-Oberfläche).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — eine weitere
      `evidence/`-Datei in `BEO-PGC/kommentar-herkunft-als-kette` (angelegt von
      `slice-code-kommentare-kennungen`) oder ein neuer Eintrag; kein Anfall ist
      ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** L — Schätzung des Umfangs; die Zahl der Kandidaten ist gemessen:
597 (Stand `d13ab81e`, `make kommentar-kennungen COUNT=1`, §1) in bis zu 243
Go-Dateien (gemessen: Dateien mit einer Kennung in einer Kommentarzeile, Stand
`7b70b34a`), acht Tranchen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Go-Dateien der Tranchen T1 bis T7 (§1) | update (nur Kommentare) | Kandidaten nach den Regeln je Kandidat (§1). |
| `*_test.go` der Tranche T8 (§1) | update (nur Kommentare) | dieselbe Regel; Testfall-Provenienz ohne Kennungs-Kette bleibt. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis, neu geschrieben | trägt 15 Zeilen-Lokatoren auf `test/integration/*.go` (gemessen, Suchlauf Zeile 4); der Runner von `make test-integration` schreibt sie neu, wenn die Kommentare der Dateien gekürzt sind. |
| `harness/sensors/kommentar-kennungen.md` | prüfen | keine Zahl des Bestands im Vertrag; ein Stand-Satz dort wäre Chronik. |

- **Reihenfolge und Commits:** je Tranche mindestens ein Commit; der Commit nennt
  die Tranche und die Zahlen vor und nach. Nach jeder Tranche laufen `make test`,
  `make fmt-check` und `make a-check` (Exit direkt, §3.9).
- **Messung:** `make kommentar-kennungen COUNT=1 PATHS=<Suchraum>` vor und nach
  jeder Tranche, dazu `TESTS=exclude` bzw. `TESTS=only`.
- **Übergabe aus `slice-antragsqueue-lesefehler-failed` (Kommentar mit ungenauer
  Herkunftsangabe, Tranche T5).** Der Godoc von `rejectionMessage`
  (`internal/adapters/driven/postgresstorage/sqlexec/translate.go`) nennt „die
  Reihenfolge der Prüfungen des Konstruktors“. Der Konstruktor
  `NewAdministrationRequest` (`internal/domain/model/administrationrequest.go`) prüft
  Kennung, Quelle, Schema und Tabelle in einem Ausdruck und nennt einen Grund; die
  Reihenfolge dieser vier gehört `rejectionMessage`, gebunden durch die Tabelle von
  [`SPEC-019`](../../../../spec/pflichtenheft.md) und die Paar-Fälle in
  `translate_test.go`. Die Wiedergabe der Reihenfolge im Godoc wird in T5 auf diese
  Zuordnung umformuliert; die Zusage der Stelle (Klartext je Grund in der Reihenfolge
  der Tabelle) bleibt. Der Fund ist die Verifikation V-2 des Slice
  `slice-antragsqueue-lesefehler-failed`.

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaft: „die Kommentare des
Go-Baums tragen ihre Herkunft als ein Feld“ und ihre Träger — Zeilen-Lokatoren
auf Dateien, deren Kommentare sich ändern. Stand ist der Plan-Stand `7b70b34a`;
der Parent am Start ist ein anderer Commit, der Implementer misst neu und ergänzt
die Zeilen mit Stand `diff`:

```suchlauf
7b70b34a 2233 -E '//.*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.go'
7b70b34a 125 -l -E '//.*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.go' ':!*_test.go'
7b70b34a 8 -E '//.* ff\.' -- '*.go'
7b70b34a 15 -E '\.go:[0-9]+' -- AGENTS.md README.md harness spec docs/user
```

| Träger | Befund | Behandlung |
|---|---|---|
| `docs/user/e2e-abdeckung.md` | 14 Zeilen tragen Lokatoren `test/integration/<Datei>.go:<Zeile>` (Zeile 4 des Feldes, `git grep -c`: 14 in dieser Datei) | ein Erzeugnis: der Lauf von `make test-integration` schreibt es neu (Liefer-Punkt 2); nicht von Hand ändern. |
| `harness/sensors/coverage-gate.md` | eine Zeile (Zeile 249) trägt den Lokator `main.go:44` (Zeile 4 des Feldes, `git grep -c`: 1 in dieser Datei) — ein Aufruf in `cmd/pg-change-feed` (Ziel-Datei am Start prüfen), den T6 berührt | fremde Datei: Meldung an den Planner mit Adresse, wenn T6 die Zeile verschiebt; der Lokator wird dort nachgezogen (Frist: die Closure dieses Slice). |
| Weitere Living Docs mit Zeilen-Lokatoren auf Go-Dateien | keine Fundstelle in `AGENTS.md`, `README.md`, `spec`, `docs/user` außer der ersten Zeile (dasselbe Feld) | Suchraum am Start neu messen. |
| Zitate von Kommentar-Text in Living Docs | nicht gemessen — ein Muster über Kommentar-Text ist nicht ableitbar | Lese-Handlung des Reviewers je Tranche (§3.13 Grenze). |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-code-kommentare-kennungen` und
`slice-harness-fmt-check` liegen in `done/` (das Werkzeug ist die Messung, `make
fmt-check` die Format-Probe der Tranchen), `slice-transformationen-e2e-abhilfe`
liegt in `done/`, und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).
Grund für die dritte Bedingung: die Test-Tranche T8 kürzt Kommentare in den
Testdateien, die die beiden E2E-Slices der Welle erweitern
(`test/integration/integration_test.go`, der Runner-Aufruf in
`tools/harness/run-integration-tests.sh` liegt nicht in Go und bleibt außen), und
ihr `make test-integration`-Lauf schreibt das Erzeugnis
`docs/user/e2e-abdeckung.md` neu; ein früherer Start ließe die E2E-Slices auf
gekürzten Kommentaren aufbauen und den Lauf doppelt fahren. **Reihenfolge
(Empfehlung an den Orchestrator):** 1. `slice-code-kommentare-kennungen`, 2.
`slice-harness-fmt-check`, 3. `slice-antragsqueue-lesefehler-failed`, dann die
Welle [welle-transformationen](../welle-transformationen.md), 4. dieser Slice
nach `slice-transformationen-e2e-abhilfe`; zu
`slice-transformationen-betriebsdoku` gibt es keine Kante (nur Doku). Will der
Auftraggeber die Nicht-Test-Tranchen früher, schneidet der Planner den Plan
**vor dem Start** entlang T7|T8 (siehe die erste Rückführung): die Tabelle und
die DoD tragen den Schnitt bereits (Liefer-Punkte 1 und 2).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erwartet, nicht
  ausgeschlossen — ein Review trägt acht Tranchen mit 597 Kandidaten (gemessen,
  Stand `d13ab81e`) unter Umständen nicht. Bedingung: der Reviewer nennt im Report
  eines Review-Laufs den Diff nicht tragbar. Schnitt entlang der Liefer-Punkte:
  **B1** = T1 bis T7 (Nicht-Test-Code) und **B2** = T8 (Test-Code) samt
  Restmenge und dem Erzeugnis-Lauf; die bis dahin gelandeten Tranchen bleiben in
  B1.
- `in-progress` → `open` (blockiert): (a) eine Klasse von Kommentaren, die zwei
  Anker legitim braucht, ohne dass §3.7 sie entscheidet — Architect-Frage zur
  Konkretisierung, kein Ausnahme-Eintrag im Werkzeug; (b) ein `make
  test-integration`-Lauf in T8, der aus Gründen außerhalb des Diffs rot ist —
  Carveout nach Modul 7 statt stilles Rot.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + zwei Review-Reports ohne offenes HIGH oder
MEDIUM (nach T7 und nach T8) + Verifikation, dass die DoD trägt +
`make kommentar-kennungen COUNT=1` mit 0 oder der benannten Restmenge + ein
realer, grüner `make test-integration`-Lauf (T8) + Closure-Notiz mit Lerneintrag
geschrieben; der Planner legt in der Closure den Folge-Slice
`slice-kommentar-kennungen-skripte` in `open/` an (§7).

## 6. Risiken und offene Punkte

- **Zusagen gehen beim Kürzen verloren**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, 5×). *Erwartet, zu
  belegen durch:* der Reviewer liest je Tranche eine Stichprobe gegen den Code
  (Zahl der Stichproben im Bericht); eine nicht getragene Zusage wird gemeldet,
  nicht gestrichen. **Ausgang:** *(bei Closure)*
- **Ein in-place schreibendes Textwerkzeug am Repo bei 597 Kandidaten**
  (`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`, 4×, verkörpert:
  der PreToolUse-Guard blockt `sed -i`, `perl -i`, `awk -i inplace` und einen
  Host-Interpreter auf Repo-Pfaden, die Grenze steht in `MR-003`; Umleitungen
  und ein Skript über den Baum liest er nicht). *Erwartet, zu belegen durch:*
  die Änderung läuft mit dem Edit-Werkzeug je Datei, nie mit `sed -i`,
  `perl -pi` oder einem Skript über den Baum; der Diff-Nachweis („nur
  Kommentarzeilen“) und `make test` sind die Probe. **Ausgang:** *(bei Closure)*
- **Der Diff ändert Code statt nur Kommentare.** *Erwartet, zu belegen durch:*
  der Diff-Nachweis je Tranche (Liefer-Punkt 1 und 2), `make test` und `make
  a-check`. **Ausgang:** *(bei Closure)*
- **`gofmt` formt einen gekürzten Doc-Kommentar um** (Zeichenpaare werden zu
  typografischen Anführungszeichen). *Erwartet, zu belegen durch:* `make
  fmt-check` nach jeder Tranche; ein Fund heißt: der Kommentar wird
  gofmt-stabil umformuliert. **Ausgang:** *(bei Closure)*
- **Die Restmenge wird zur Ausnahmeliste.** *Erwartet, zu belegen durch:* jede
  Restzeile trägt Klasse und Grund, das Werkzeug kennt keinen Marker und keine
  Liste ([`AGENTS.md`](../../../../AGENTS.md) §3.2 gilt sinngemäß); der Reviewer
  liest die Liste. **Ausgang:** *(bei Closure)*
- **Gekürzte Kommentare ersetzen die Kennungs-Kette durch Vorher-Nachher-Sprache
  oder Chronik** (`BEO-PGC/slice-chronik-in-code-kommentar`, 9×). *Erwartet, zu
  belegen durch:* der Kandidatenlauf aus Schritt 20 des Implementer-Ablaufs
  (Chronik) auf den Dateien der Tranche und der Reviewer-Punkt „Kommentar trägt
  keine der Kommentar-Klassen“. **Ausgang:** *(bei Closure)*
- **Der `make test-integration`-Lauf in T8 ist rot oder flackert**
  (`BEO-PGC/test-integration-retention-timing-flake`, 3×). *Erwartet, zu belegen
  durch:* ein Lauf mit gedruckter Laufzeit; ein rotes Ergebnis außerhalb des
  Diffs ist ein Carveout (§4), kein stilles Rot. **Ausgang:** *(bei Closure)*
- **Die Zahlen der Tranchen bewegen sich mit jedem Commit** (Zustandsgröße,
  [`AGENTS.md`](../../../../AGENTS.md) §3.12). *Erwartet, zu belegen durch:* jede
  Zahl im Bericht trägt Befehl, Stand und Lauf; die Zahlen dieses Plans sind
  am Stand `d13ab81e` gemessen (§1). **Ausgang:** *(bei Closure)*
- **Die Größe trägt einen Review nicht.** *Erwartet, nicht ausgeschlossen:* zwei
  Review-Läufe (nach T7, nach T8); die Rückführung §4 ist der Ausgang, falls der
  Reviewer den Diff nicht tragbar nennt. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu tragen bei Closure — geschärfte Regel · neuer
  Sensor · benannte Spec-Lücke; erwartet: die Antwort auf die Trigger-Frage der
  Gate-Aufnahme in `slice-code-kommentare-kennungen` §1 und, falls Fälle mit zwei
  Ankern auftraten, die Konkretisierung von
  [`AGENTS.md`](../../../../AGENTS.md) §3.7; ohne Eintrag kein `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure)*
- **Restmenge:** *(zu tragen bei Closure — `Datei:Zeile`, Klasse, Grund je
  Kandidat, oder „0“)*
- **Folge-Slices:** `slice-kommentar-kennungen-skripte` (Kommentare außerhalb von
  Go — Skripte, `Makefile`, `harness/mk`, `.sql`, `.yml`, Dockerfiles): der
  Planner legt die Datei bei dieser Closure in `open/` an, damit der Punkt eine
  Adresse hat, die ihn annimmt.
- **Risiken aus §6:** *(je ein Ausgang, zu tragen bei Closure)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt
  die drei Paarungen (Anker · Folge-Slice · Register), nach dem `git mv` nach
  `done/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield), quer über alle Pakete; Adapter, Domäne, Composition Root,
Werkzeuge und Tests sind keine eigenen Sub-Areas — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
gemessen am 2026-09-26 mit `ls evidence | wc -l` je Eintrag):

- `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (verkörpert, 5×),
  `BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, 9×),
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (verkörpert, 3×):
  Risiken in §6 (Zusagen, Chronik).
- `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 4×):
  Risiko in §6; der Guard blockt die Flag-Formen, und der Slice ändert Hunderte
  Kommentare und ist der Ort, an dem ein Umleitungs- oder Skript-Weg (vom
  Guard ungelesen) eintreten kann.
- `BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×): der
  `make test-integration`-Lauf in T8 schreibt `docs/user/e2e-abdeckung.md`
  absichtlich (Erzeugnis); der Wrapper deckt die Rollout-Artefakte, nicht diese
  Datei — kein Auftreten.
- `BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×): Risiko in
  §6.
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 13×),
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32×),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23×): das
  Suchlauf-Feld in §3 und die Zahlen mit Ursprung.
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
