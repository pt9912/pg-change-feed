# Slice code-kommentare-kennungen: Herkunft im Code-Kommentar ist ein Feld — Regelschärfung, Kandidaten-Werkzeug, Erstbeleg

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
`slice-transformationen-map-value` voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5).

**Bezug:** [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen in Trägern; die Grenze „kein Sensor auf Prosa“ bleibt),
[`AGENTS.md`](../../../../AGENTS.md) §3.7 (ein Kommentar beschreibt, was da
ist) und §3.13, Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein
Kommentar trägt, Architect-Verdikt
[`architect-verdict-slice-chronik-in-code-kommentar`](../../../reviews/architect-verdict-slice-chronik-in-code-kommentar.md)
(kein Textmuster-Sensor über Satz-Subjekte),
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (Anker der Datei des
Erstbelegs).

**Berührte Spec-Stellen:** — (Harness-Regel und Werkzeug; keine Spec-Stelle).

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Auftrag des Auftraggebers zur Kennungsdichte in
Code-Kommentaren. **Datum:** 2026-09-26.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Kommentar im Go-Code trägt seine Herkunft als **ein**
auflösbares Feld: die Regel steht in [`AGENTS.md`](../../../../AGENTS.md) §3.7
in konkreter Form, `make kommentar-kennungen` listet die Kommentarblöcke, die
sie verletzen (Kandidaten), Schritt 20 des Implementer-Ablaufs und der
Reviewer-Skill rufen es als Probe auf, und die Datei des Erstbelegs ist
bereinigt.

**Ausgangslage — gemessen, Stand `7b70b34a`, Befehle im Suchlauf-Feld (§3).**
Der Code hält §3.7 nicht ein: der Godoc von `Publish` in
`internal/application/port/outbound/changestream.go` nennt in **einem** Absatz
zwei Entscheidungen und drei Anforderungen (eine davon mit „ff.“) und gibt
dazu Inhalt der Spec in eigenen Worten wieder. Die Baseline sagt: Herkunft als
ein auflösbares Feld, nie als Absatz (`grundlagen-harness-dateien.md` §Was ein
Kommentar trägt). Der Bestand: 2233 Kommentarzeilen mit einer Kennung in 243
Go-Dateien (125 Nicht-Test-Dateien), davon 31 Zeilen mit drei oder mehr
Kennungen und 8 mit „ff.“ (der Auftrag nannte 9; die Zahl im Plan ist die
gemessene). Die 2233 sind eine **Obergrenze des Suchraums, kein
Verstoß-Zählwert**: eine Zeile mit genau einer Kennung ist die zulässige Form.
Die Größenordnung der Kandidaten (ein Kommentarblock mit mindestens zwei
verschiedenen Kennungen oder „ff.“ nach einer Kennung) ist **erwartet, nicht
belegt**: ein nicht committeter Wegwerf-Prototyp mit `go/parser`-Kommentargruppen
(Lauf am Stand `7b70b34a`, mit keinem Repo-Befehl wiederholbar) zählte 396
Blöcke in Nicht-Test-Dateien und 192 in Testdateien bei 1474 Kommentarblöcken
mit mindestens einer Kennung; die Basis-Messung des Werkzeugs (DoD 2) ersetzt
diese Zahlen.

**Zielform der Regel** (der Wortlaut ist Sache des Implementers, der Inhalt ist
gesetzt; kein `MR-*`, weil die Baseline die Regel bereits trägt und der Code
sie verletzt — keine Adaption): (1) höchstens **eine** Kennung je Kommentar, als
Rang-Zeiger auf die Norm, deren Umsetzung die Stelle ist; keine Kette, keine
Kompaktform (`…003/005`), kein „ff.“. (2) Der Kommentar gibt keine Aussage der
Spec oder einer ADR in eigenen Worten wieder; er trägt, was die **Stelle**
zusagt, koppelt, abgrenzt oder nicht leistet (die fünf Klassen), und verweist
für die Norm auf den einen Anker. (3) Eine Kopplung nennt die mitzuändernde
**Stelle** (Datei, Funktion), nicht eine Reihe von Kennungen.

**Zielform des Erstbelegs:** der Godoc von `Publish` trägt das Verhalten der
Stelle — verteilt an die zum Aufrufzeitpunkt registrierten Abonnenten, blockiert
nie, jeder Abonnent trägt eine begrenzte Warteschlange und Überlauf wird für ihn
verworfen, keine Zustellgarantie, ein Fehler entsteht nur aus einem ungültigen
Aufruf oder einem beendeten `ctx` und beeinflusst Persistierung und Bestätigung
nicht — und höchstens einen Anker
([`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)); der Satz zur
Nachvollziehbarkeit wird ein kurzer Verweis auf den Lesezugriffsweg ohne
Wiederholung der Spec.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Gate.** Ein Lauf über den Bestand färbte sich rot, bis
  `slice-code-kommentare-bereinigung` durch ist, und die Aufnahme in `make gates`
  braucht eine ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6, §4). Trigger für
  die Frage, nicht Umfang: die Bereinigung endet mit Kandidatenzahl 0, und ein
  Kommentar mit Kennungs-Kette erreicht trotz gelaufenem Werkzeug den Review —
  dann Architect-Frage mit ADR-Vorschlag.
- **Ein Textmuster-Sensor über Satz-Subjekte oder Chronik.** Das Verdikt
  [`architect-verdict-slice-chronik-in-code-kommentar`](../../../reviews/architect-verdict-slice-chronik-in-code-kommentar.md)
  schließt ihn aus: „Testfall-Provenienz oder Produktionsverhalten-Chronik“ ist
  ein Satz-Subjekt-Urteil. Das Werkzeug zählt verschiedene Kennungen je
  Kommentarblock und entscheidet kein Subjekt; der Chronik-Kandidatenlauf in
  Schritt 20 bleibt unverändert.
- **Die Bereinigung des Bestands** über die Datei des Erstbelegs hinaus —
  `slice-code-kommentare-bereinigung`. Mit ihr wüchse der Diff auf Hunderte
  Kommentarblöcke, die ein Review nicht trägt; ohne das Werkzeug fehlte der
  Bereinigung die Messung.
- **Kommentare außerhalb von Go** (Skripte, Makefile, `.sql`, `.yml`,
  Dockerfiles): das Werkzeug liest Go-Kommentargruppen; die anderen
  Kommentarformen sind eine Folgetranche mit benannter Adresse in
  `slice-code-kommentare-bereinigung` §1.
- **Die Wahrheit eines Kommentars.** Das Werkzeug prüft die Form (Zahl der
  Kennungen), nicht, ob ein Satz zutrifft oder eine Spec-Aussage in eigenen
  Worten wiedergibt; das ist Lese-Handlung des Reviewers
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12, §3.13).
- **`sdks/`.** Kennungen in Kommentaren der SDK-Bäume prüft
  `make sdk-public-doc-check` (strenger: keine Kennung); das Werkzeug liest
  `sdks/` nicht.

## 2. Definition of Done

- [ ] **Liefer-Punkt 1 — die Regel und ihre Träger.**
      [`AGENTS.md`](../../../../AGENTS.md) §3.7 trägt die Konkretisierung der
      Zielform aus §1 (eine Kennung je Kommentar, keine Kette, keine
      Kompaktform, kein „ff.“, keine Spec-Wiedergabe in eigenen Worten,
      Kopplung nennt die Stelle);
      `.claude/commands/implement-slice.md` Schritt 20 ruft den
      diff-skopierten Kandidatenlauf des Werkzeugs auf (der bestehende
      Chronik-Kandidatenlauf bleibt daneben) und nennt seine Grenze;
      `.harness/skills/reviewer.md` trägt einen Lese-Unterpunkt „Herkunft als
      mehrere Felder, Kette, ‚ff.‘ oder Spec-Wiederholung“ mit dem Lauf als
      Probe. Einstufung: Vorschlag MEDIUM für einen im Diff neu geschriebenen
      Kommentar (Form ohne falsche Aussage); eine falsch wiedergegebene
      Spec-Aussage bleibt unter dem bestehenden HIGH-Punkt „Kommentar trägt
      keine der Kommentar-Klassen“; der Implementer entscheidet die Einstufung
      mit Begründung im Bericht. *Zu belegen durch:* Lesen der drei Stellen,
      `make docs-check` und `make suchlauf-nachmessen` über das Feld in §3.
- [ ] **Liefer-Punkt 2 — das Werkzeug.** `make kommentar-kennungen
      [PATHS=<Pfade>] [COUNT=1] [TESTS=exclude|only] [DIFF=<Basis>]` ist
      netzlos und Docker-only (gepinntes Toolchain-Image, `--network none`, Repo
      lesend gemountet), kein Gate und nicht in `make gates`. Es liest die
      Kommentarblöcke der `.go`-Dateien unter `PATHS` (Standard: der Baum;
      ausgenommen `gen/`, `sdks/`, `.git`, `.harness`) und meldet einen
      **Kandidaten**, wenn ein Block mindestens zwei verschiedene Kennungen der
      vier Arten (`ADR-…`, `LH-FA-…`/`LH-QA-…`, `SPEC-…`, `ARC-…`) trägt — eine
      Kompaktform wie `…-003/005` zählt als zwei — oder „ff.“ nach einer
      Kennung. Ausgabe je Kandidat: `Datei:Zeile-Zeile` und die Kennungen. Exit 1
      bei mindestens einem Kandidaten, 0 sonst, 2 bei Eingabefehler; `COUNT=1`
      druckt nur die Zahl (Exit 0), `TESTS=exclude|only` trennt Nicht-Test- und
      Testdateien. Mit `DIFF=<Basis>` meldet es nur Blöcke, die eine seit
      `<Basis>` hinzugefügte Zeile überlappen (Eingabe: ein `git diff -U0`-Strom,
      die Pipe läuft unter `bash` mit `set -o pipefail`,
      [`AGENTS.md`](../../../../AGENTS.md) §3.9): Bestandskandidaten färben den
      Lauf eines Implementers nicht, solange die Bereinigung läuft. Es gibt
      keinen Ausnahme-Pfad (keine Marker, keine Ausnahmeliste,
      [`AGENTS.md`](../../../../AGENTS.md) §3.2). *Form (Empfehlung, am Start zu
      bestätigen):* ein Go-Programm `tools/harness/kommentar-kennungen/` —
      `go/parser`-Kommentargruppen treffen die Blockgrenzen exakt, eine Kennung
      in einem Zeichenketten-Literal ist kein Kommentar; nur Standardbibliothek;
      der Tabellentest ist ein Go-Test im selben Paket (läuft unter `make test`).
      Die Alternative `bash` + `git` wie `make suchlauf-nachmessen` ist
      verworfen: ohne Parser sind Blockgrenzen und Literale nicht sauber
      erkennbar (hergeleitet, nicht erprobt). *Zu belegen durch:* der
      Tabellentest mit den Fällen — eine Kennung (kein Kandidat) · zwei
      verschiedene Kennungen in einem Block · dieselbe Kennung zweimal (kein
      Kandidat) · Kompaktform · „ff.“ · Kennung in einem Zeichenketten-Literal
      (kein Kandidat) · Endkommentar hinter Code · Testdatei · Direktive
      `//go:build` · Diff-Modus (Block überlappt eine hinzugefügte Zeile /
      überlappt keine) · Exit-Codes und `COUNT=1` —, je an seine Eingabe
      gebunden; die Mutationen der Eingabeseite je einmal rot gesehen:
      Schwelle „mindestens zwei“ auf „mindestens drei“, Kompaktform als eine
      Kennung gezählt, Überlappungsprüfung des Diff-Modus invertiert, Kennung
      aus einem Zeichenketten-Literal mitgezählt. Die Selbstanwendung
      (`make kommentar-kennungen PATHS=tools/harness/kommentar-kennungen` endet
      mit 0 Kandidaten). Die Basis-Messung `make kommentar-kennungen COUNT=1`
      (gesamt, `TESTS=exclude`, `TESTS=only`) mit ihrem Lauf im Bericht — Ursprung
      „gemessen“ mit Stand und Befehl
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12). Der Vertrag
      `harness/sensors/kommentar-kennungen.md` (Definition des Kandidaten, Form
      der Aufrufe, Exit-Codes, Grenze: prüft die Form, nicht die Wahrheit; ein
      grüner Lauf sagt nicht „die Kommentare sind konform“; kein Gate) und eine
      Zeile in [`harness/README.md`](../../../../harness/README.md) §Sensors,
      Werkzeuge; `make a-check` grün (Gruppe `tooling`, nur Standardbibliothek).
- [ ] **Liefer-Punkt 3 — der Erstbeleg.** Alle Kandidaten in
      `internal/application/port/outbound/changestream.go` sind bereinigt
      (mindestens drei Blöcke, gelesen: `ErrChangeStream`, `ChangeStreamPort`,
      `Publish`; die Zahl am Start gemessen und im Bericht genannt); der Godoc
      von `Publish` trägt die Zusagen aus §1 vollständig und höchstens den Anker
      [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md); nur
      Kommentarzeilen ändern sich. *Zu belegen durch:* `make
      kommentar-kennungen PATHS=internal/application/port/outbound/changestream.go`
      mit 0 Kandidaten; der Diff der Datei ohne Zeile, die nicht mit `//`
      beginnt (Befehl und Zahl 0 im Bericht); `make test` grün; `gofmt -l` über
      die Datei im gepinnten Toolchain-Image ohne Ausgabe (Schritt 18 des
      Implementer-Ablaufs; `make fmt-check` gibt es erst nach
      `slice-harness-fmt-check`); der Reviewer liest, dass keine Zusage
      verloren ist.
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
- [ ] Doku-Update: [`AGENTS.md`](../../../../AGENTS.md) §3.7 und
      [`harness/README.md`](../../../../harness/README.md) §Sensors (Liefer-Punkte
      1 und 2); das Benutzerhandbuch bleibt unberührt (keine
      Betreiber-Oberfläche).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis (Arbeitsname `BEO-PGC/kommentar-herkunft-als-kette`, Beleg:
      dieser Slice) oder eine weitere `evidence/`-Datei; kein Anfall ist
      ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** M — Schätzung, nicht gemessen: eine Regel mit drei Trägern, ein
neues Go-Programm samt Tabellentest, ein Makefile-Ziel, ein Vertragsdokument
und eine bereinigte Datei.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.7 | update | Konkretisierung der Herkunft als ein Feld (Liefer-Punkt 1); der Satz „Herkunft in **ein** auflösbares Feld“ steht dort bereits, die Konkretisierung nennt die drei Fälle (Kette, Kompaktform und „ff.“, Spec-Wiedergabe). |
| `.claude/commands/implement-slice.md` Schritt 20 | update | diff-skopierter Aufruf des Werkzeugs neben dem Chronik-Kandidatenlauf; der Satz „Herkunft steht nur als **ein** auflösbares Feld“ bleibt. |
| `.harness/skills/reviewer.md` | update | Lese-Unterpunkt neben „Kommentar trägt keine der Kommentar-Klassen“ und „Slice-/Wellen-Chronik“. |
| `tools/harness/kommentar-kennungen/` (`main.go`, `main_test.go`) | neu | Programm und Tabellentest (Form: Empfehlung in §2). |
| `tools/harness/kommentar-kennungen.sh` | neu | Aufrufer: Docker-Aufruf, der `DIFF`-Strom unter `bash` mit `set -o pipefail`. |
| `Makefile` | update | Ziel `kommentar-kennungen` mit `.PHONY` und Hilfe-Zeile; kein Eintrag in den Gate-Zielen. |
| `harness/sensors/kommentar-kennungen.md` | neu | Vertrag, Definition des Kandidaten, Grenze. |
| `harness/README.md` | update | Zeile unter Werkzeuge. |
| `internal/application/port/outbound/changestream.go` | update | Erstbeleg (Liefer-Punkt 3). |
| `.a-check.yml` | prüfen | die Gruppe `tooling` deckt `tools/harness/**`; das Programm importiert nur die Standardbibliothek — keine Änderung erwartet. |

- **Diff-Modus:** die Hunk-Köpfe `@@ -a,b +c,d @@` des `git diff -U0`-Stroms
  liefern die hinzugefügten Zeilenbereiche je Datei; ein Block, dessen
  Zeilenbereich einen davon schneidet, gehört zum Diff. Ein Block, den der Diff
  nur berührt, ohne eine Zeile zu ändern, gehört nicht dazu.

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaften: „ein Kommentar
trägt seine Herkunft als ein Feld“ (Go-Bestand, Regel und die drei Träger, die
sie beschreiben). Suchraum: der ganze Baum für die Go-Zeilen, die vier
Träger-Wurzeln für die Regel; ausgenommen sind `docs/reviews/**`, die Records
unter `done/` und `.harness/baseline/**`. Das Muster trägt die vier Kennungsarten
(Symbolnamen), das Zählwort („ff.“, drei Kennungen) und die Beschreibung („ein
auflösbares Feld“, die Regelnummer); Stand ist der Parent `7b70b34a`, der
Implementer ergänzt die Zeilen mit Stand `diff`:

```suchlauf
7b70b34a 2233 -E '//.*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.go'
7b70b34a 243 -l -E '//.*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.go'
7b70b34a 125 -l -E '//.*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.go' ':!*_test.go'
7b70b34a 31 -E '//.*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3}).*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3}).*(ADR-[0-9]{4}|LH-(FA|QA)-[A-Z]{3}-[0-9]{3}|SPEC-[0-9]{3}|ARC-[0-9]{3})' -- '*.go'
7b70b34a 8 -E '//.* ff\.' -- '*.go'
7b70b34a 2 -E 'auflösbares Feld' -- AGENTS.md .claude .harness/skills harness
7b70b34a 6 -E '§3\.7' -- AGENTS.md .claude .harness/skills harness
```

| Träger | Befund | Behandlung |
|---|---|---|
| Go-Kommentare (Bestand) | 2233 Zeilen mit einer Kennung in 243 Dateien, 31 mit drei oder mehr, 8 mit „ff.“ (Zeilen 1 bis 5 des Feldes) | Die Bereinigung ist `slice-code-kommentare-bereinigung`; hier nur die Datei des Erstbelegs. Die Zahlen sind ein Suchraum, kein Verstoß-Zählwert. |
| `AGENTS.md` §3.7 und §3.13 | „Herkunft in **ein** auflösbares Feld“ (§3.7, eine Fundstelle); §3.7 wird an sechs Stellen der Wurzeln zitiert (Zeile 7 des Feldes) | §3.7 wird konkretisiert; §3.13 bleibt (andere Regel: bewegte Eigenschaft eines Trägers); die Zitate anderer Träger von §3.7 (`.harness/skills/closure-note-reviewer.md`, `harness/sensors/coverage-gate.md`) meinen Zustandsfelder und bleiben unberührt. |
| `.claude/commands/implement-slice.md` Schritt 20 | trägt den Satz „Herkunft steht nur als **ein** auflösbares Feld“ und den Chronik-Kandidatenlauf | Ergänzung um den Aufruf; kein anderer Schritt nennt die Kennungszahl. |
| Beobachtungs-Register `BEO-PGC/slice-chronik-in-code-kommentar` (`state.md`) | führt „kein mechanischer Sensor“ als Ausgang | Träger der Planner-Closure: ein Verweis auf das Werkzeug (anderes Muster: Kennungen je Block, kein Satz-Subjekt) und die Abgrenzung; fremde Datei, deshalb Meldung, keine stille Änderung. |
| `sdks/` | Kennungen dort deckt `make sdk-public-doc-check` | unberührt. |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/`
(WIP-Limit 1). **Reihenfolge (Empfehlung an den Orchestrator):** 1. dieser
Slice, 2. `slice-harness-fmt-check`, 3. `slice-antragsqueue-lesefehler-failed`,
danach die Welle [welle-transformationen](../welle-transformationen.md) in ihrer
Tabellen-Reihenfolge, `slice-code-kommentare-bereinigung` nach
`slice-transformationen-e2e-abhilfe`. Dieser Slice muss `done` sein, **bevor**
`slice-transformationen-map-value` startet (Kante, dort §4): der Implementer
jedes folgenden Slices der Welle läuft Schritt 20 mit dem Werkzeug, und ein
Slice, der Kommentare schreibt, bevor die Regel steht, vergrößert den Bestand,
den die Bereinigung wegräumt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Regel, Werkzeug
  und Erstbeleg nicht in einem Review tragen — der abtrennbare Teil ist der
  Diff-Modus des Werkzeugs (Liefer-Punkt 2) als eigener Slice mit Start vor
  `slice-transformationen-map-value`; Regel, Grundmodus und Erstbeleg bleiben.
- `in-progress` → `open` (blockiert): (a) die Basis-Messung zeigt, dass die
  Definition „mindestens zwei verschiedene Kennungen“ mehr als die Hälfte der
  Kommentarblöcke mit Kennung trifft (der Prototyp zählte 588 von 1474, etwa 40
  %; erwartet, nicht belegt) — die Regel widerspräche dann der gelebten Praxis:
  Architect-Frage (ist ein Rang-Zeiger mit zwei Ankern zulässig?), kein Anpassen
  der Definition an das Ergebnis; (b) das Programm baut nur mit Netz oder mit
  einem Modul-Download (dann Architect-Frage zur Form: `bash` + `git`).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Review-Report ohne offenes HIGH oder
MEDIUM + Verifikation, dass die DoD trägt + Closure-Notiz mit Lerneintrag
geschrieben (geschärfte Regel: [`AGENTS.md`](../../../../AGENTS.md) §3.7; neuer
Sensor: `make kommentar-kennungen`, ohne Gate).

## 6. Risiken und offene Punkte

- **Die Definition des Kandidaten trifft legitime Kommentare** (eine Kopplung
  oder Abgrenzung mit zwei Ankern). *Erwartet, zu belegen durch:* eine
  Stichprobe der Basis-Messung (30 Kandidaten quer über die Pakete; Auswahl,
  Klasse und Urteil „Verstoß“ oder „legitim“ je Kandidat im Bericht); ein hoher
  Anteil legitimer Fälle ist die Rückführung nach `open/` (§4). **Ausgang:**
  *(bei Closure)*
- **Das Werkzeug erweckt den Eindruck, ein grüner Lauf heiße „Kommentare
  konform“** (`BEO-PGC/regel-weiter-als-ihr-sensor`, 3×, Restrisiko: Werkzeug
  grün, Muster trifft nicht). *Erwartet, zu belegen durch:* der Vertrag nennt die
  Grenze im ersten Absatz (Form, nicht Wahrheit; eine Spec-Wiedergabe in eigenen
  Worten erkennt es nicht), und der Reviewer-Unterpunkt nennt den Lauf als
  Probe, nicht als Beleg. **Ausgang:** *(bei Closure)*
- **Der Slice wiederholt den verworfenen Chronik-Sensor**
  (`BEO-PGC/slice-chronik-in-code-kommentar`, 9×, Verdikt: kein Sensor).
  *Erwartet, zu belegen durch:* der Vertrag nennt die Abgrenzung (Kennungen je
  Block, kein Satz-Subjekt, keine Slice-/Wellen-Nummer), und der Reviewer liest,
  dass das Werkzeug keine der beiden auswertet. **Ausgang:** *(bei Closure)*
- **Die Regel wird enger gelesen als gemeint** — „höchstens eine Kennung“
  verdrängt eine nötige Kopplung oder streicht Zusagen mit. *Erwartet, zu
  belegen durch:* die Konkretisierung nennt die Kopplung als **Stelle** (Datei,
  Funktion) und lässt die fünf Klassen unberührt; der Reviewer liest den
  Erstbeleg gegen den Code (keine Zusage verloren,
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, 5×). **Ausgang:**
  *(bei Closure)*
- **Das Programm braucht Netz oder einen Modul-Download.** *Erwartet, zu belegen
  durch:* der Lauf mit `--network none` und nur der Standardbibliothek; sonst
  die Rückführung nach `open/` (§4 (b)). **Ausgang:** *(bei Closure)*
- **Die Zahlen des Plans bewegen sich mit jedem Commit** (Zustandsgröße,
  [`AGENTS.md`](../../../../AGENTS.md) §3.12). *Erwartet, zu belegen durch:* jede
  Zahl im Bericht trägt Befehl, Stand und Lauf; die Zahlen dieses Plans stehen
  mit dem Stand `7b70b34a`. **Ausgang:** *(bei Closure)*
- **Der `DIFF`-Strom läuft über Host-`git` und `bash`** (dieselbe Klasse wie
  `make suchlauf-nachmessen`; `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`,
  1×). *Erwartet, zu belegen durch:* der Vertrag nennt die Host-Werkzeuge; ein
  zweites Auftreten der Klasse ist die Architect-Frage zu
  [`AGENTS.md`](../../../../AGENTS.md) §3.1. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu tragen bei Closure — erwartet: geschärfte Regel
  [`AGENTS.md`](../../../../AGENTS.md) §3.7 und neuer Sensor `Makefile:kommentar-kennungen`
  ohne Gate; das Feld `liegt in` steht nur, wenn sie wirklich verkörpert sind;
  ohne Eintrag kein `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure —
  erwartet: `BEO-PGC/kommentar-herkunft-als-kette` neu angelegt, Beleg
  `evidence/slice-code-kommentare-kennungen.md`, Ausgang verkörpert, Abgrenzung
  zu `BEO-PGC/slice-chronik-in-code-kommentar`; dazu der Verweis dort)*
- **Folge-Slices:** `slice-code-kommentare-bereinigung` (Bereinigung des
  Bestands, Datei in `open/`).
- **Risiken aus §6:** *(je ein Ausgang, zu tragen bei Closure)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt
  die drei Paarungen (Anker · Folge-Slice · Register), nach dem `git mv` nach
  `done/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); `tools/harness`, `harness/`, `AGENTS.md` und die
Ablauf-Träger sind keine eigenen Sub-Areas — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
gemessen am 2026-09-26 mit `ls evidence | wc -l` je Eintrag):

- `BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, 9×): Abgrenzung in §6;
  Schritt 20 wird ergänzt, der Chronik-Lauf bleibt.
- `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (verkörpert, 5×) und
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (verkörpert, 3×): der
  Reviewer-Punkt „Kommentar trägt keine der Kommentar-Klassen“ bleibt, der neue
  Unterpunkt steht daneben.
- `BEO-PGC/regel-weiter-als-ihr-sensor` (teilweise verkörpert, 3×): das
  Restrisiko „Werkzeug grün, Muster trifft nicht“ ist ein Risiko in §6.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 14×) und
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32×): das
  Suchlauf-Feld in §3 mit Codeblöcken.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23×): jede
  Zahl dieses Plans trägt Ursprung und Stand.
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 16×): die
  Mutationen der Eingabeseite in Liefer-Punkt 2.
- `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (offen, 1×):
  Risiko in §6.
- `BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen, 1×): abgegrenzt,
  `sdks/` bleibt bei `make sdk-public-doc-check`.
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
