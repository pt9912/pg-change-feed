# Slice implementieren (Harness)

Argument: $ARGUMENTS

Dieser Command führt die **Implementer**-Rolle (Modul 9) für *einen* Slice — innerhalb der
Rollen-Sequenz Planner → Architect → Implementer → Reviewer → Verifier → Validator →
Planner-Closure (Modul 8). **Rollen-Trennung ist Kontext-Trennung:** die nachgelagerten Rollen
(Review, Verifikation, Validation, Closure) laufen in **frischem Kontext** (Subagent / geleerter
Kontext), nie im Kontext, der den Code schrieb — sonst wiederholt sich derselbe blinde Fleck.
Keine Rolle springt rückwärts ohne Übergabe-Artefakt (Findings · Folge-ADR · Carveout, Modul 8).

Kanonische Quellen (vendored Regelwerk, `.harness/baseline/<tag>/regelwerk/`): Modul 9
(Implementierung), Modul 5 (Lifecycle), Modul 8 (Rollen), Modul 10 (Review), Modul 11
(Verifikation).

## Repo-lokale Adaptionen, die du beachten MUSST (ANPASSEN an dein Repo)

<!-- ANPASSEN: Dieser Block listet die Adaptionen DEINES Repos gegenüber der Baseline
     (dein `harness/conventions.md`, „MR-Block"). Der Bootstrap hat eine
     Durchsetzungsschicht emittiert (Stop-Hook, Command-Guard, Gate-Nachweis,
     Doc-Gate); die daraus folgenden, workflow-relevanten Adaptionen stehen unten.
     Ergänze/streiche nach deinem Repo. -->

Über das Regelwerk hinaus trägt dein Repo lokale Adaptionen gegenüber der Baseline. Lies den
Adaptions-Block („MR-Block") in `harness/conventions.md`; die workflow-relevanten (aus der
emittierten Durchsetzungsschicht):

- **Docker-only, kein Host-Toolchain.** Jeder Gate und jedes Tool läuft in einem gepinnten
  Docker-Image; erlaubte Host-Werkzeuge, das Verbot des in-place Text-Umschreibens
  (`sed -i`, `perl -pi`, Host-Interpreter) und das Verbot, den Text einer Repo-Datei per
  Umleitung zu schreiben (`>`, `>>`, `tee`, Heredoc) stehen in `AGENTS.md` §3.1. Rufe nur
  `make`-Targets auf. Verweigert die Berechtigungsschicht einen Aufruf, gilt `AGENTS.md`
  §3.15: melden, nicht auf anderem Weg wiederholen.
- **Gate-Nachweis + Stop-Hook.** `make gates` endet mit `record-gates`, das einen Content-Hash des
  Working Tree stempelt; der Stop-Hook verweigert den Abschluss, solange der aktuelle Tree nicht
  passt. **Jede Inhaltsänderung nach einem Gate-Lauf — inklusive jedes Commits und jedes `git mv`
  — macht den Stempel ungültig: `make gates` erneut laufen.** Ein Commit/Move ohne frischen
  Gate-Lauf lässt den Stop-Hook rot.
- **Strenges Doc-Gate (d-check, `.d-check.yml`).** Aktiv: `links`, `anchors`,
  `ids` (Kennungs-Linkpflicht für `LH-FA/QA-*`, `SPEC-*`, `ARC-*`, `ADR-*` —
  nackte Kennungen in Prosa brauchen einen Link auf ihr Definitions-Dokument;
  Inline-Code ist im Default `prose` nicht linkpflichtig; `MR-*` ist nicht
  linkpflichtig konfiguriert), `matrix` (Decken-Regel: spec ⊬ ADR/Slice,
  ADR ⊬ Slice; eine Referenz auf eine superseded ADR nur über
  `allow-supersede-lineage` bei deklariertem `Supersedes`), `versions`
  (Baseline-Pin) und `structure` (Register-Spalten). `codepaths` ist bewusst
  deaktiviert (Bedingung: alle referenzierten Pfade existieren).
  `docs/reviews/**` wird gescannt — Kennungen dort verlinken.
- **Neue Artefakte per `cp` aus den vendored Templates** (`.harness/baseline/<tag>/templates/…`),
  dann ausfüllen — keine handgeschriebenen oder repo-gepflegten Template-Kopien.
- **Commit via Message-Datei** (`git commit -F <datei>`): der Guard liest den Command-String
  quote-bewusst (ein Tool-Token in den Anführungszeichen einer Inline-Message blockt nicht), ein
  Heredoc-Text oder ein unbalanciertes Anführungszeichen (Apostroph) blockt aber wieder ein
  Tool-Token darin — die Message steht in einer Datei, nie inline.
- **Commit-Message-Kennungen.** **Struktur-IDs (`SPEC-*`, `ARC-*`) gehören NICHT in die
  Commit-Message** (`AGENTS.md` §5) — nur `LH-*`/`ADR-*`. Mechanisch
  getragen: `make commit-traceability` prüft je Message der letzten 5 Commits beide
  Grenzen (d-check `commits` + `tools/harness/commit-traceability.sh`,
  [`ADR-0045`](../../docs/plan/adr/0045-commit-traceability-standing-gate.md))
  — `make gates` färbt einen Verstoß rot, bevor der Review ihn sieht.

## Kontext lesen (Modul 9, Schritte 1–3)

1. `CLAUDE.md` lesen (falls dein Repo eines führt — das agentseitige Briefing).
2. `harness/README.md` lesen.
3. `AGENTS.md` lesen.
4. `harness/conventions.md` lesen.
5. Den Regelwerk-Index (`.harness/baseline/<tag>/regelwerk/README.md`) und das aufgabenrelevante
   Modul **on-demand** lesen (Source Precedence, committet vendored Baseline). Nicht den ganzen
   Baum laden.
6. Die als Argument übergebene Slice-Datei lesen.
7. Alle referenzierten ADRs und Anforderungen lesen.
8. Berichten: Slice-ID · LH-IDs · ADR-IDs · betroffene Komponenten · zu laufende Gates.

## Nach in-progress eintreten (Modul 5 Lifecycle + Modul 8 Übergabe)

9. Der Implementer erhält den Slice **in `in-progress/`** (Planner→Implementer-Übergabe,
   Modul 8; `next → in-progress` = „Implementer beginnt", Modul 5). Liegt er noch in `open/`,
   zuerst dorthin verschieben (`open → next → in-progress`); `open → next` setzt dabei das
   Kopf-Feld `Verantwortlich:`. Jedes `git mv` ist ein **reiner Move, getrennt vom Inhalt
   committet** (Hard Rule 3.3), und `next → in-progress` landet **auf dem Hauptzweig, vor der
   Arbeit** — der Branch entsteht danach. Reist der Move erst im PR mit, ist der Zustand
   zweigelokal, und `in-progress/` bleibt für alle anderen leer, bis die Arbeit fertig ist.
10. WIP-Limit = 1 pro Implementer (Modul 5): kein paralleles `in-progress/`.
11. Lifecycle-Rücksprungkanten (Modul 5), falls sich der Slice als falsch erweist: zu groß →
    `in-progress → next` (zurück zur Zerlegung); blockiert → `in-progress → open` (Carveout,
    Modul 7). Zurückführen ist Disziplin, kein Scheitern.
    **Rückführung mit fortgesetzter Lieferung:** Lieferst du trotz
    Rückführung den unabhängigen Teil noch im selben Lauf (statt den Slice komplett
    ruhen zu lassen), braucht die spätere Planner-Rückkehr nach `in-progress/` einen
    **eigenständigen Architect-Verdikt-Zug** (frischer Kontext) — keine Planner-
    Selbstbestätigung in derselben Session (Modul 8, Rollen-Trennung ist Kontext-
    Trennung). Der Rückkehr-Move landet außerdem vor jedem weiteren
    Produktions-Commit des fortgesetzten Teils, nicht danach.

## Plan vor Code (Modul 9, Schritt 4 — nicht optional)

12. **Den Ist-Zustand gegen den Slice-Plan messen, bevor du editierst** (`grep`/`diff`, nicht
    `edit`) — Geschwister-Slices lassen Pläne altern (gelöschte Pfade, verschobene
    Lifecycle-Dateien). Drift zuerst abgleichen; keinen veralteten Plan blind abarbeiten.
13. Die kleinste sinnvolle Änderung gegen die DoD planen. Erst planen, dann coden.

## Implementieren und gaten (Modul 9, Schritte 5–6)

14. Die kleinste sinnvolle Änderung implementieren.
    **Plan-Nachzug im selben Lauf:** Jede über den
    Slice-Plan hinausgehende Änderung wird **vor** dem Sensor-Lauf (15)
    in die §3-Tabelle des Slice-Plans eingetragen — neu gelieferte
    Dateien und Artefakte ebenso wie Nicht-Realisierungen geplanter
    Punkte; eine Reduktion steht als §3-Zeile mit Begründung im Plan,
    sie wird nicht still gestrichen. Der Plan-Commit gehört in denselben
    Lauf; ein Plan-Nachzug durch den Planner ist ein Review-Befund,
    kein Ablauf-Schritt.
15. Zuerst den engsten nützlichen Gate laufen lassen (z. B. eine Testdatei / ein Gate).
16. `make gates` laufen lassen.

**Plan-Defekt-Rücksprungkanten (Modul 9):** ein roter Sensor (15) oder rotes Gate (16) führt
zurück zum **Plan** (13) — den Plan verfeinern, nicht den Kontext neu lesen. Ein Rücksprung zu
Schritt 1 signalisiert einen Kontext-Defekt. Ein struktureller Fehlschnitt (zu groß / blockiert)
ist eine Lifecycle-Rücksprungkante (11).

## Pre-completion-Checkliste (Modul 9, Schritt 8 — letzte Handlung der Implementer-Rolle)

17. Doku, ADR-Index und README aktualisieren, falls ein öffentlicher Vertrag berührt ist.
    **Handbuch-Versionshistorie im selben Diff**
    (`BEO-PGC/handbuch-versionshistorie-uebersprungen`):
    Berührt dieser Lauf `docs/user/benutzerhandbuch.md` inhaltlich, zieht derselbe
    Diff **zwingend** den `Version:`-Kopf hoch und ergänzt eine neue Zeile in
    `### Änderungshistorie` — in Betreibersicht und **ohne Kennungen**. Schreib-
    regeln und Kandidatenlauf: Skill `.harness/skills/nutzerdoku-schreiben.md`;
    Fangnetz: `make handbuch-public-doc-check` (liest Kennungen, nicht Sinn).
    **Grenze:** wie bei Schritt 20 (§Grenze dieser Selbstprüfung) erste, nicht
    tragende Linie; die tragende ist der unabhängige Reviewer
    (`.harness/skills/reviewer.md`, eigener benannter HIGH-Punkt).

    **Neue Betreiber-Oberfläche zieht das Handbuch mit**
    (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`): Führt dieser Lauf eine neue
    Betreiber-Oberfläche ein — eine Umgebungsvariable des Feed-Containers
    (`CDC_*`, als Konstante in `internal/bootstrap/wiring.go`), eine
    administrative SQL-Funktion (`cdc.*`, in `tools/schema/`), eine
    Horch-Adresse oder einen Endpunkt (Adapter unter
    `internal/adapters/driving/`) —, so zieht **derselbe Diff**
    `docs/user/benutzerhandbuch.md` mit (§5 „Umgebungsvariablen des
    Feed-Containers", §4 „Aufgaben") **oder benennt den Aufschub mit
    Adresse**: eine Folge-Slice-ID, die die Doku nachholt. Ein Aufschub ohne
    Adresse ist keiner; die Versionshistorie folgt der Regel direkt darüber.
    Kandidatenlauf: `git diff --name-only <Basis> -- internal/bootstrap/ tools/schema/ internal/adapters/driving/` — trifft er eine neue Oberfläche, ohne dass `docs/user/benutzerhandbuch.md` im selben Diff liegt, ist sie unversorgt. **Grenze:** wie oben; die tragende Linie ist der unabhängige Reviewer (`.harness/skills/reviewer.md`, eigener benannter HIGH-Punkt).

    **Aufschub-Adresse deckt den Gegenstand · seit welle-transformationen**
    (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, 5×; Architect-Zug des
    Lese-Schritts, Modul 6): Ein benannter Aufschub trägt nur, wenn die Adresse
    die Sendung **annimmt** — zwei Pflichten zusätzlich zur bloßen Nennung der
    Folge-Slice-ID. (1) Vor dem Eintragen der Adresse: `git grep` der
    Kernbegriffe des aufgeschobenen Gegenstands (Funktionsnamen, Parameter,
    Aufrufform, betroffene Fälle) im Plan der Adresse — 0 Treffer heißt, die
    Adresse ist falsch gewählt oder ihr DoD-Ausschnitt zu eng gefasst; vor
    Handoff klären, nicht nur behaupten. (2) Der Gegenstand steht danach als
    **committeter Text** im Plan der Adresse (ein Übergabe-Block, z. B. in
    dessen §2), nicht nur als Meldung im eigenen Bericht — ein Bericht erreicht
    den, der die Adresse später abarbeitet, nicht, wenn dieser nur den Plan
    liest. Fünffach real aufgetreten (`slice-071`, `slice-072`,
    `slice-transformationen-antragsweg-schema`,
    `slice-transformationen-backfill-pfad`,
    `slice-antragsqueue-lesefehler-failed`), jedes Mal von Reviewer oder
    Verifier vor dem Merge gefunden, nie automatisiert. **Grenze:** wie oben;
    die tragende Linie ist der unabhängige Reviewer
    (`.harness/skills/reviewer.md`, eigener benannter HIGH-Punkt, Probe:
    `git grep` der Kernbegriffe im Plan der Adresse).
18. Die Pre-completion-Checkliste laufen: die DoD Punkt für Punkt **behaupten** und die
    **Sensor-Belege** anhängen — `make gates` **und die Nicht-Gate-Sensoren, die den Slice
    betreffen** (die dein Repo führt — z. B. ein Mutations-Sensor, wenn Wächter neu/geändert sind;
    ein Emit-/Integrations-Smoke, wenn der betroffene Pfad berührt ist). Modul 11 verlangt genau
    hier den Lauf: *„der Implementer-Agent läuft `make verify-*` **selbst** vor der
    ‚fertig'-Meldung"* — ein Sensor, der erst zur Wellen-Closure feuert, ist pro Slice keiner.
    **Ein nicht gelaufener Sensor ist ein Befund, kein Formfehler:** ihn wegzulassen ist eine
    Aussage („betrifft diesen Slice nicht"), die begründet werden muss. Kein Gate erzwingt das —
    der Stop-Hook deckt nur `make gates`. Das ist die *Behauptung* der Implementer-Rolle und die
    *Eingabe* des Verifiers — **nicht** das finale DoD-Urteil (Modul 11: „Behauptung ohne
    Bestätigung ist die häufigste Verifier-Lücke"; eine DoD-Verletzung ist eine Verifier-only-Klasse,
    unsichtbar für Review und Tests). Ausgeführte Sensors + Restrisiken berichten.
    **DoD-Checkbox-Nachzug im selben Lauf
    (`BEO-PGC/dod-checkbox-nachzug`):** „behaupten" heißt hier auch: jede zu
    diesem Zeitpunkt materiell erfüllte oder korrekt entfallene DoD-Zeile in §2 wird
    **in diesem Lauf** von `[ ]` auf `[x]` gesetzt, nicht nur im Bericht behauptet und
    der Planner-Closure zur Nacharbeit überlassen. Nur Punkte, die die Rollen-Sequenz
    zu diesem Zeitpunkt noch nicht durchlaufen haben (Review, Verifikation,
    Register-/Risiko-Ausgänge — Planner-Closure-Arbeit), bleiben regulär `[ ]`.
    Eine DoD-Zeile, die einen Beleg zusagt („zu belegen durch …", „je eine Mutation
    im Bericht"), wird erst abgehakt, wenn der **Plan-Text selbst** den Beleg trägt
    (Befehl, gedruckte Ausgabezeile, Lauf, bei Mutationen Stelle und Farbe) oder auf
    einen committeten Report verweist; der Bericht des Implementers ist kein Träger.
    Herkunft: `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` (5×) · seit welle-routing.
    **Suchlauf nachmessen:** Trägt der
    Plan ein Suchlauf-Feld (Codeblöcke mit dem Etikett `suchlauf`,
    [`AGENTS.md`](../../AGENTS.md) §3.13 §Suchform), läuft
    `make suchlauf-nachmessen PLAN=<Plan-Datei>` vor der „fertig"-Meldung und nach
    jeder Fixrunde; ein Exit ≠ 0 ist ein Befund (die Zahl im Plan nachziehen oder
    die Abweichung erklären), der Lauf steht im Bericht. **Grenze:** das Werkzeug
    prüft Zahlen und Stände, nicht die Vollständigkeit von Suchraum und Muster
    ([`harness/sensors/suchlauf-nachmessen.md`](../../harness/sensors/suchlauf-nachmessen.md)).
    **Format
    (`BEO-PGC/formatierungs-drift-ohne-gate`):** Trägt der Diff Go-Dateien, läuft
    `make fmt-check` (Docker-only, netzlos, Repo lesend gemountet; `gofmt -l` über alle
    Go-Dateien des Baums, Exit 0 = formatiert), vor dem „fertig“ und nach jeder Fixrunde.
    Eine gemeldete Datei wird nach der Ausgabe von `gofmt -d` korrigiert (nie mit einem
    Textwerkzeug wie `sed` am Quelltext; der Aufruf von `gofmt -d` steht im Vertrag).
    Die gedruckte Zeile des Laufs steht im Bericht. **Grenze:** wie Schritt 20; die
    tragende Linie ist der Reviewer (`.harness/skills/reviewer.md`, LOW). Werkzeug,
    kein Gate ([`harness/sensors/fmt-check.md`](../../harness/sensors/fmt-check.md)).
19. **Zu jedem neuen oder geänderten Wächter die rot färbende Mutation benennen**
    (`AGENTS.md` §3.6). Ein grüner Gate-Lauf belegt nur, dass nichts *bricht* — nicht, dass
    der Wächter greift. Pro Zusage also: *welche Änderung am geprüften Code müsste diesen
    Test rot machen, und wurde sie einmal gesehen?* Wo die Antwort dauerhaft interessant
    ist, gehört sie in den Mutations-Sensor deines Repos (falls vorhanden); wo sie einmalig ist, in
    den Bericht. **Keine Antwort ist ein Befund**, kein Formfehler.
    **Die Richtung der Mutation**
    (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`):
    Der Satz oben sagt **dass** mutiert wird, nicht **wo**. *Eine Zusage ist
    nur dann gebunden, wenn der Test an ihrer **Eingabeseite** rot werden kann: mutiere den
    **Eingabewert**, nicht nur die Ausgabeseite. Wer nur den Fake oder den Rückgabewert mutiert,
    prüft den Fake — die Aussage bleibt grün, egal was der Adapter mit der Eingabe tut. Fehlt die
    Mutation der Eingabeseite, ist die Zusage **grün ohne Aussage**: ein Befund, kein Formfehler.*
    **Enumerations-Pflicht statt Erinnerung** (dieselbe Form wie Schritt 20): **je Zusage eine
    benannte Eingabeseiten-Mutation** — die Liste lautet *Zusage · mutierte Eingabe · gesehenes
    Rot*, und wo sie leer bleibt, steht der Grund. Die bloße **Zahl** der gefahrenen Mutationen
    trägt nicht. Der Finder-Träger derselben Regel ist der HIGH-Unterpunkt
    „Zusage ohne Bindung an ihre Eingabeseite" in `.harness/skills/reviewer.md`.
20. **Jeden in diesem Lauf neu geschriebenen oder geänderten Kommentar gegen `AGENTS.md` §3.7
    prüfen** (Code, Konfiguration, Skripte). Die Probe: beschreibt der Satz den **Ist-Zustand**
    (indikativ, auflösbar), oder trägt er eine Slice-Nummer als Begründung, ein „(… , entschieden)"
    ohne Anker-Form, oder einen Konjunktiv über eine verworfene Alternative bzw. eine noch nicht
    existierende künftige Änderung (**„sobald Slice X das tut …"**)? Herkunft steht nur als **ein**
    auflösbares Feld in den in `AGENTS.md` §3.7 genannten Formen — alles andere ist Zustand,
    keine Chronik, und wird vor der Übergabe umformuliert statt mitgeschleift.
    **Enumerations-Pflicht statt Erinnerung** (`BEO-PGC/slice-chronik-in-code-kommentar`):
    ein visueller Scan übersieht Fundstellen; deshalb vor der Übergabe
    zusätzlich einen **diff-skopierten** (nicht repo-weiten) Kandidatenlauf gegen genau die in
    diesem Lauf geänderten `.go`-/`tools/schema/*.sql`-Dateien ausführen, z. B.:
    `git diff --name-only <Basis> -- '*.go' 'tools/schema/*.sql' | xargs -r grep -nE
    'slice-[0-9]+|welle-[0-9]+|vor diesem [Ss]lice|nach diesem [Ss]lice|seit diesem [Ss]lice'`.
    Repo-weit liefe derselbe Lauf gegen die etablierte, zulässige Testfall-Provenienz-Zitierform
    (Godoc-Kommentare wie „TestXyz trägt/deckt … aus `review-slice-NNN.md` F-x") und würde in
    Rauschen ertrinken; deshalb nur der Diff dieses
    Laufs. Jeder Treffer bekommt **eine** Probe: Begründet der Satz, **warum ein Testfall/
    Regressionsfall existiert** (Subjekt: der Test — zulässige Provenienz), oder **warum sich der
    Produktionscode aktuell so verhält** (Subjekt: die Funktion/der Code-Pfad — Chronik,
    unzulässig)? Nur Letzteres wird umformuliert. Kein Sensor/Gate dafür — die Unterscheidung
    ist ein Satz-Subjekt-Urteil, kein Zeichenkettenmuster; dieser Schritt bleibt Disziplin.
    **Konjunktiv über die verworfene Alternative** (`BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar`):
    derselbe Lauf auf die hinzugefügten Kommentarzeilen des Diffs, in Go (`//`) und in Skripten
    (`#`; die Kommentare dort sind transliteriert, „waere“, „wuerde“) —
    `git diff -U0 <Basis> -- '*.go' '*.sh' '*.awk' | grep -nE '^\+.*(//|#).*(wäre|waere|würde|wuerde|hielte|hätte|haette|sonst|statt|\b(trüge|truege|bliebe|ließe|liesse|könnte|koennte|müsste|muesste|bräuchte|braeuchte|läge|laege|stünde|stuende|käme|kaeme|wären|waeren|gäbe|gaebe|ginge|fände|faende|dürfte|duerfte|hieße|brächte|braechte)\b)'`.
    Die Konjunktiv-II-Formen hinter `sonst|statt` tragen eine Wortgrenze (`\b`).
    Jeder Treffer bekommt ein Urteil: die Zusage der Stelle im
    Indikativ (zulässig) oder die Beschreibung einer verworfenen Alternative (umformulieren);
    Mutationsbeschreibungen in Test-Godocs („sonst bliebe der Test grün“) und normale Zweige
    („sonst auf stdout“, `else`-Zweige in Shell-Kommentaren) sind zulässig; Treffer ohne
    Kommentar (`${#var}`, `$#`) sind keine.
    **Grenze dieser Selbstprüfung:** Dieser Schritt
    läuft im selben Kontext, der den Kommentar geschrieben hat — genau die Konstellation,
    vor der Modul 8 §Kernidee warnt („wer geschrieben hat, reviewt nicht“). Er bleibt Pflicht,
    weil er die Zahl der Fixrunden senkt, ist aber **nicht** die tragende
    Verteidigungslinie — das ist der unabhängige Reviewer (`.harness/skills/reviewer.md`,
    eigener benannter HIGH-Punkt). Ein Auftreten trotz gelaufenem
    Schritt 20 ist kein Beleg für einen defekten Prozess, solange der Reviewer den Fall vor
    Merge fängt — das ist der Regelfall, für den die Rollentrennung sorgt.
    **Herkunft als ein Feld**
    ([`AGENTS.md`](../../AGENTS.md) §3.7): neben dem Chronik-Kandidatenlauf läuft
    der diff-skopierte Kandidatenlauf des Werkzeugs über die eigenen Änderungen —
    `make kommentar-kennungen DIFF=<Basis>` (neue Dateien vorher `git add`;
    `TESTS=exclude` trennt die Nicht-Test-Dateien). Ein Kandidat ist ein
    Kommentarblock mit mindestens zwei verschiedenen Kennungen oder „ff.“; er
    wird auf **eine** Kennung als Rang-Zeiger umformuliert, eine Kopplung nennt
    die mitzuändernde Stelle (Datei, Funktion) statt einer Kennungsreihe. Der
    Chronik-Kandidatenlauf sucht Slice-/Wellen-Nummern, das Werkzeug zählt
    Kennungen und liest kein Satz-Subjekt
    ([`harness/sensors/kommentar-kennungen.md`](../../harness/sensors/kommentar-kennungen.md)).
    Das Werkzeug prüft die **Form**, nicht die Wahrheit: eine Spec-Aussage in
    eigenen Worten hinter einer einzigen Kennung fängt es nicht, ein Lauf ohne
    Kandidat ist keine Konformitätsaussage. Es gilt die Grenze der Selbstprüfung oben.

Hier endet die Implementation. Die übrigen Rollen laufen in **getrennten Kontexten** (Modul 8).

## Übergaben an nachgelagerte Rollen (Modul 8 → 10 → 11)

21. **→ Reviewer (Code-Review, Modul 10):** den Diff + Plan-Verweis an einen **unabhängigen**
    Reviewer übergeben (`.harness/skills/reviewer.md`, frischer Kontext — kein Selbst-Review). Er
    kategorisiert Findings (HIGH/MEDIUM/LOW/INFO) in einen Report unter `docs/reviews/` und prüft
    den Diff gegen **Plan + ADR + Hard Rules** (nicht die DoD). HIGH/MEDIUM auflösen; ein HIGH mit
    Rollen-Konflikt folgt Modul 8 §Konflikt-Pfad (Sequenz mit Übergabe-Artefakten, nie
    „herabstufen, weil der Implementer widerspricht").
    **Fixrunden-Checkbox-Nachzug:** Löst eine Fixrunde nach
    Reviewer-Findings einen bislang offenen DoD-Punkt auf (typischerweise „Review
    durchgeführt, … kein offenes HIGH"), wird die zugehörige Checkbox **im
    Fixrunden-Commit** mitgesetzt — nicht erst bei der Planner-Closure nachgetragen.
    **Re-Review nach der Fixrunde:** Ändert die Fixrunde Produktionslogik oder eine
    Norm (Spec-Zeile, ADR-Wortlaut), oder hat nach der Fixrunde kein anderer Kontext
    sie ausgeführt, wird ein Re-Review verlangt, bevor die Closure läuft. Eine
    Fixrunde, die nur Text, Kommentare oder Test-Code ändert und die der Verifier in
    frischem Kontext ausgeführt hat (Lauf, Mutation), braucht keines; der Verifier
    nennt im Report, was er ausgeführt hat. Herkunft:
    `BEO-PGC/fixrunde-ohne-reviewer-lesung` (3×) · seit welle-routing.
22. **→ Verifier (Modul 11):** in getrenntem Kontext die DoD-/Spec-Behauptung und den
    Plan-vs-Code-Diff **bestätigen**, dazu ADR-Konformität. Das fängt, was Tests übersehen und der
    Reviewer nicht sieht (DoD-Verletzung).
23. **→ Validator (Modul 8):** falls der Slice End-Nutzer-Wert liefert, gegen den realen Bedarf
    validieren („das Richtige bauen"). Meist n/a bei interner Wartung — dann explizit sagen statt
    still überspringen. Verschiebt der Slice den Schritt auf den Wellen-Beleg, nennt die Notiz
    den Abschnitt „Validator-Feststellung (Modul 8)" der Results-Notiz der Welle als Adresse.

## Closure — Planner-Rolle (Modul 8 + Modul 5)

24. Erst wenn der Review konform **und** die Verifikation die DoD bestätigt hat, schließt der
    **Planner**: die Closure-Notiz mit einem **Steering-Loop-Eintrag** schreiben (geschärfte Regel ·
    neuer Sensor · benannte Spec-Lücke — Modul 5: der `→ done`-Übergang verlangt einen Lerneintrag,
    nicht nur grüne Gates), dann den Slice `in-progress → done` verschieben (`git mv`, eigener
    Commit, getrennt vom Inhalt — Hard Rule 3.3). Ein rotes Gate erreicht `done/` **nur** mit
    dokumentiertem Carveout (Modul 7), nie als stilles Rot. **Jedes offene Risiko aus dem Slice-Plan
    bekommt dabei genau einen von drei Ausgängen** (Modul 5): *eingetreten* → Carveout oder
    Folge-Slice mit ID · *entfallen* → gestrichen **mit Begründung** · *weiter offen* → wandert ins
    Beobachtungs-Register (Schritt 25). Ein Slice geht nicht nach `done/`, während ein Risiko ohne
    Ausgang dasteht.
    **§7-Vorlagenrest beim Füllen entfernen**
    (`BEO-PGC/vorlagenrest-in-closure-notiz`; mechanisch getragen durch die
    `structure`-Regel auf `done/slice-*.md` §7 in `.d-check.yml`): Die §7 wird aus
    der Vorlage kopiert und behält deren Guidance für die Anker-Zeile — die
    Teil-Zeile `— liegt in …`, die `Auslöser:`-Platzhalter-Zeile und den kursiven
    Ausfüll-Hinweis. Wird mit dem Slice nichts verkörpert (der Normalfall), fällt
    die Guidance **ersatzlos** weg; nur eine echte Verkörperung trägt
    `— liegt in <Zielort>` mit Herkunfts-Anker. Kandidatenlauf vor dem `git mv`,
    auf dem **rohen** Text (die `structure`-Regel liest den bereinigten, in dem
    Backtick-Spans geleert sind):
    `grep -nE 'Wurde mit diesem Slice nichts verkörpert|Auslöser: .BEO-<(NNN|KUERZEL)>|— liegt in .<' <slice-datei>`
    — jeder Treffer in §7 wird ersatzlos entfernt oder zu einem echten Zielort gefüllt.
    **Grenze:** die tragende Linie ist die `structure`-Regel (`.d-check.yml`).
25. **Das Beobachtungs-Register fortschreiben** (`docs/plan/planning/observations/`, Modul 6) —
    der **Schreib**-Schritt, und er hängt an der Closure, nicht an der Implementation. Für jede
    Beobachtung aus der Closure-Notiz: führt das Register die Klasse schon, dann die vorhandene
    Kennung `BEO-<KUERZEL>/<slug>` **zitieren** und eine weitere Datei in ihrem `evidence/` anlegen
    — wer neu formuliert, spaltet eine Klasse in zwei Pfade, und keiner der beiden erreicht je 3×.
    Sonst ein neues Verzeichnis `BEO-<KUERZEL>/<slug>/` mit `observation.md` und `state.md` anlegen
    — Kürzel aus der Modus-Deklaration nachschlagen, nicht erfinden; das Register ist zugleich die
    Vergabestelle für den `<slug>`-Teil. Der Beleg ist **formgebunden**: `evidence/slice-<Kennung>.md`,
    kein Freitext, eine Datei je Auftreten. Geschrieben wird er **vor** dem `git mv` — die
    Slice-Datei liegt dann noch nicht in `done/`, und das ist richtig so, weil Move und Inhalt
    getrennt committen (Hard Rule 3.3). Der Zähler wird **nicht gesetzt**, er ist die Zahl der
    Evidence-Dateien und **folgt** aus ihnen. **Bei null Beobachtungen** bleibt die Ablage
    unverändert und die Closure-Notiz trägt den Satz *keine Beobachtung angefallen*: das Auslassen
    ist keine Antwort. Erreicht ein Eintrag **mit diesem Slice** 3× (die Zahl seiner
    Evidence-Dateien), wandert er in die Steering-Loop-Einträge der laufenden Welle-Closure
    (`/close-welle`); läuft keine Welle, löst die Slice-Closure den Lese-Schritt selbst aus, und der
    Herkunfts-Anker lautet dann `seit slice-<Kennung>` statt `seit welle-<Kennung>`.

    **Verweisform auf wandernde Slice-Pläne**
    (`BEO-PGC/slice-pfad-als-link-in-berichten`; mechanisch getragen für
    `docs/reviews/**` und `observation.md` durch die `structure`-Regeln in
    `.d-check.yml`): Ein Slice-Plan wandert
    (`open/` → `next/` → `in-progress/` → `done/`). Ein Markdown-Link mit festem
    Verzeichnis (`…/planning/in-progress/slice-<Kennung>.md`) löst am Ist-Ort auf und
    bricht erst beim nächsten `git mv` (`links` → `target-missing`). Stattdessen
    die Kennung zitieren (`slice-<Kennung>`) oder einen Inline-Code-Pfad. Kandidatenlauf
    auf den **rohen** Zeilen der von diesem Lauf berührten Dokumente:
    `git diff --name-only <Basis> -- 'docs/**' | xargs -r grep -nE '\]\([^)]*(open|next|in-progress)/slice-'`
    — jeder Treffer wird zitiert statt verlinkt. **Grenze:** `state.md` und
    `evidence/*.md` tragen keine Überschrift und damit keinen Abschnitts-Anker;
    die `structure`-Regeln greifen dort nicht, dieser Schritt ist dort die
    einzige Linie.

Gates nicht überspringen. Keine Erfolgsmeldung ohne Command-Ausgabe.
