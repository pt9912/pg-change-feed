# Slice slice-074: E2E-Abdeckungstabelle je Spec-Kennung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine beobachtbare Closure-Bedingung, die
von der DoD dieses Slice verschieden wäre: `make test-integration` erzeugt
die Datei, `make gates` hält sie gegen Drift, und die drei realen Belege
(Erzeugung, Idempotenz, Rot-Beleg) stehen in seiner eigenen DoD. Ein
Wellen-Trigger schriebe damit die DoD ab (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht (Modul 6)).

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (die E2E-Kette
`make test-integration` ist die Quelle des Erzeugnisses). Keine funktionale
Anforderung ändert ihr Verhalten; die Tabelle macht die bereits geführten
E2E-Nachweise (u. a. zu [`LH-QA-REL-001`](../../../../spec/lastenheft.md),
[`LH-QA-SEC-001`](../../../../spec/lastenheft.md),
[`LH-QA-OPS-005`](../../../../spec/lastenheft.md)) je Kennung sichtbar.
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md) trägt die
Beleg-Semantik erzeugter Artefakte.

**Berührte Spec-Stellen:** — (kein `SPEC-NNN` und kein `ARC-NNN` berührt: die
E2E-Kette ist Messmethode des Lastenhefts, kein Element der Technik- oder
Sicht-Strate).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/e2e-abdeckung.md` liegt als **erzeugtes Artefakt** vor —
eine Tabelle, die die E2E-Nachweise des Repos je Spec-Kennung führt
(`Spec-Kennung | Nachweis | Ort | Kurzbeschreibung`) und die der E2E-Lauf
selbst schreibt: `make test-integration` erzeugt sie, ein Commit erfolgt nur
bei inhaltlicher Änderung (Vorbild für erzeugte Artefakte mit Beleg-Charakter:
`harness/image-hash.txt` aus `make image`, `tools/schema/plan.yaml` als
Pflicht-Report von `make schema-rollout`; Beleg-Semantik erzeugter Artefakte
in [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)). Damit verschwindet
die Drift-Richtung *„die Tabelle behauptet einen Nachweis, den es nicht
gibt"* **per Konstruktion** — und ebenso die Gegenrichtung *„ein Nachweis
existiert, fehlt aber in der Tabelle"*.

**Erzeugungsweg — entschieden, mit seinen zwei Hälften.** Der Erzeuger ist
eine Erweiterung der **bestehenden** E2E-Kette, kein neues Werkzeug. Die
Go-Hälfte (die `func TestE2E*` in `test/integration/integration_test.go`)
leitet das **Go-Testpaket selbst** aus seinem eigenen Quelltext ab
(`go/parser`, stdlib): Funktionsname und Zeile aus dem AST, Spec-Kennungen und
Kurzbeschreibung aus dem Doc-Kommentar der Funktion — der Zeilensatz entsteht
damit nicht aus dem, *was gelaufen ist*, sondern aus dem, *was im Paket
steht*; eine Testfunktion ohne Spec-Kennung im Doc-Kommentar lässt den
Erzeuger **sichtbar abbrechen**, statt sie stillschweigend wegzulassen. Die
Bash-Hälfte (die Phasen des Runners: Rollen-DSN-Verifikation, Lasttest,
Black-Box-CLI-Rundlauf, Diagnose, Retention-Lebenszyklus, Upgrade-Rundlauf)
**deklariert** der Runner an Ort und Stelle je Phase — mit einem wörtlichen
Anker aus der Phase; findet der Runner diesen Anker nicht mehr, bricht er
sichtbar ab. Zusammengesetzt (feste Reihenfolge: Go-Zeilen nach Quellzeile,
dann Bash-Zeilen nach Runner-Zeile) und geschrieben wird die Datei **vom
Runner auf dem Host** — der Toolchain-Container läuft gegen
`-v "$(pwd)":/src:ro` und kann nicht schreiben; geschrieben wird nur bei
inhaltlicher Abweichung (Temp-Datei + `cmp`), damit ein Lauf mit unverändertem
Testbestand die Datei nicht anfasst.

**Drift-Schutz — `d-check`, kein eigenes Skript.** Die tragende Garantie ist
die **Ableitung**: Die Zeilen entstehen aus den Nachweis-Deklarationen selbst,
nicht neben ihnen. Das bestehende Doku-Gate kommt dazu: `ids` (läuft in
`make gates`) erzwingt auf der Kennungsspalte den Link auf das
Definitionsdokument, eine neue `structure`-Regel die Zeilenform (Abschnitt,
Spalten-Mindestbreiten). Was das Gate dabei **nicht** trägt, ist real gemessen
und nicht aus der Modulbeschreibung geschlossen: `ids` prüft den Link, nicht
die Existenz der Kennung in `spec/lastenheft.md` (eine verlinkte, erfundene
Kennung bleibt grün), und Kennungen in Inline-Code-Spans bleiben ungeprüft.
Die Beschreibungsspalte trägt deshalb **keine** Kennungen aus dem
Doc-Kommentar (der Erzeuger entfernt sie): die Aussage steht über den
`Ort`-Verweis auf die Testfunktion ohnehin an ihrer Quelle.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein eigener Lauf-Beleg** (Protokoll, Ergebnis-/OK-Spalte je Lauf,
  Zeitstempel) — der reale Lauf ist über den CI-Lauf belegt (`AGENTS.md`
  §3.10), und die Log-Ausschnitte stehen in den Slice- und Wellen-Abschlüssen;
  eine Lauf-Spalte änderte die Datei bei jedem Lauf und hebelte die
  Idempotenz-Zusage der DoD aus. Die Tabelle ist eine **stabile
  Abdeckungs-Deklaration**, kein Lauf-Beleg — und das ist eine Entscheidung,
  kein Versäumnis.
- **Ein eigenes Erzeuger-Skript** (oder ein Werkzeug neben der E2E-Kette) —
  ein zweiter Träger müsste die Phasen des Runners kennen und driftete gegen
  ihn; genau die Doppelquelle, die `BEO-PGC/lese-doppelquelle` als Klasse
  führt. Eine Erweiterung eines bestehenden Artefakts ist kein neues Skript.
- **Die `codepaths`-Aktivierung** (Prüfung von `datei:zeile` in Inline-Code,
  heute in `.d-check.yml` aus) — real gemessen liefert sie global aktiviert
  rund 50 Befunde, fast alle in `docs/reviews/**` (historische Lauf-Belege mit
  zu ihrem eigenen Ort relativen Pfaden) plus zwei Pseudo-Pfade (`./cmd/`,
  `./internal/`); sie braucht eine bewusst gescopte Konfiguration und ist eine
  **Verschärfung des repo-weiten Gates**, die nach `AGENTS.md` §3.6 ihren
  eigenen ADR-Pfad hat. Eigener Vorgang. Die `Ort`-Spalte steht deshalb in
  Inline-Code in `datei:zeile`-Form — genau die Form, die das Modul prüft,
  damit der Folge-Vorgang sie ohne Umbau greifen kann.
- **Ein Vollständigkeits-Sensor für den stillen Ausschluss einer
  Testfunktion** (`BEO-PGC/test-runner-stiller-ausschluss`: eine Funktion, die
  in keinem `-run`-Muster des Runners auftaucht, wird nie ausgeführt) — das
  ist ein Wächter über die `-run`-Muster des Runners, nicht über die
  Abdeckungstabelle; eigener Vorgang, die Beobachtung ist seine Adresse und
  bleibt offen. Dieser Slice schließt sie nicht und behauptet es nicht.
- **Produktionscode unter `internal/**` und Schemamigrationen unter
  `tools/schema/**`** — Schicht-Abgrenzung: das Erzeugnis dokumentiert die
  E2E-Kette; berührt sind `test/integration/**` (Testcode),
  `tools/harness/**` (Runner), `.d-check.yml`, `harness/sensors/**` und
  `docs/user/**`. `internal/**` und `tools/schema/**` bleiben damit
  unberührt.
- **Ein Ort für den Ist-Zustand der Sensordoku** *(mitgenommen, nicht
  ausgeschlossen):* `harness/sensors/docs-check.md` wird ohnehin um die neuen
  Grenzen erweitert; drei Zeilen darunter steht die Aussage, die fünfte
  `structure`-Regel sei „auskommentiert bis zur ersten Closure" — sie ist seit
  der ersten Closure aktiv. Die Sensordoku, auf deren Grenzen sich dieser
  Slice stützt, nennt mit diesem Slice den Ist-Zustand ihres Configs
  (`AGENTS.md` §3.7).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — das Erzeugnis.** `docs/user/e2e-abdeckung.md` liegt im Repo:
eine Tabelle `Spec-Kennung | Nachweis | Ort | Kurzbeschreibung`, je Zeile
mindestens eine Spec-Kennung als Link auf ihr Definitionsdokument, `Ort` als
`Datei:Zeile` in Inline-Code.

- [ ] Die Datei liegt vor; ihr Inhalt entspricht der Ausgabe eines frischen
      `make test-integration`-Laufs (nicht einem von Hand gepflegten Stand),
      und sie nennt in ihrem Kopf ihren Erzeuger und ihre Stabilitäts-Zusage
      (stabile Abdeckungs-Deklaration, kein Lauf-Beleg).
- [ ] Jede Zeile trägt mindestens eine verlinkte Spec-Kennung; der
      Nachweis der E2E-Kette ist auf beiden Trägern vertreten (eine
      `func TestE2E*` und mindestens eine Bash-Phase des Runners).

**Liefer-Punkt 2 — der Erzeuger in der bestehenden Kette** (Go-Testpaket für
die Go-Hälfte, Runner für die Bash-Hälfte, Zusammensetzung und Schreiben im
Runner; kein neues Skript, kein neues Binary).

- [ ] **Beleg (a) — Erzeugung:** ein realer `make test-integration`-Lauf auf
      leerem Stand (Datei vorher entfernt) legt `docs/user/e2e-abdeckung.md`
      an; der Exit-Code des Laufs wird direkt und ungepiped festgestellt
      (`AGENTS.md` §3.9).
- [ ] **Beleg (b) — Idempotenz:** ein zweiter realer Lauf mit unverändertem
      Testbestand lässt die Datei inhaltsgleich (`git diff --exit-code`,
      zusätzlich `cmp` gegen die Datei vor dem Lauf) — kein Churn im Diff.
- [ ] **Beleg (c) — Rot-Beleg beider Richtungen:** eine real hinzugefügte
      `func TestE2E*`-Funktion erscheint als neue Zeile in der Datei, ihr
      reales Entfernen nimmt die Zeile wieder heraus; beides mit dem
      tatsächlichen Lauf gezeigt, nicht behauptet.
- [ ] Die beiden Abbruch-Wächter greifen real: eine `func TestE2E*` ohne
      Spec-Kennung im Doc-Kommentar bzw. ein entfernter Deklarations-Anker
      einer Bash-Phase lässt den Erzeuger sichtbar fehlschlagen (Exit ≠ 0,
      Datei unverändert) — je einmal gezeigt.
- [ ] `go test -race` grün (der Erzeuger läuft im bestehenden Testpaket).

**Liefer-Punkt 3 — der deklarierte Drift-Schutz.**

- [ ] `.d-check.yml` trägt die `structure`-Regel für die Tabelle
      (Abschnitts-Adressierung über Kopfzeilen-Namen, Mindestbreiten je
      Spalte); `make docs-check` ist mit der neuen Regel grün, und die Grenze
      der Regel ist benannt (sie fängt leere/kaputte Zeilen, nicht einen
      falschen Nachweis).
- [ ] Die Grenzen der beiden tragenden Doku-Regeln sind **real gemessen und
      notiert**, nicht angenommen: `ids` erzwingt den Link auf das
      Definitionsdokument, **nicht** die Existenz der Kennung in
      `spec/lastenheft.md` (eine verlinkte, erfundene Kennung ist grün), und
      Kennungen in Inline-Code-Spans bleiben ungeprüft (beides am
      gepinnten Image real gezeigt, nicht aus der Modulbeschreibung
      geschlossen). Die tragende Garantie der Tabelle ist deshalb die
      **Ableitung aus dem Quelltext** und nicht das Doku-Gate.
- [ ] `harness/sensors/docs-check.md` benennt die Grenzen dieses Erzeugnisses:
      keine Symbol-/Funktionsnamen-Prüfung (`--trace` ist ausdrücklich keine
      Code-Prüfung, `codepaths` prüft Pfade und Zeilenbereiche), die zwei
      gemessenen `ids`-Grenzen (nur Link, keine Existenz; Code-Spans
      ungeprüft), `codepaths` aus (Folge-Vorgang) und die deklarierte
      Bash-Hälfte; und die Aussage zur fünften `structure`-Regel nennt den
      Ist-Zustand (aktiv seit der ersten Closure).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/sensors/docs-check.md` (Liefer-Punkt 3) und
      `docs/user/e2e-abdeckung.md` sind die berührten Dokumente; ein neuer
      Eintrag in `harness/README.md` §Sensors ist **nicht** fällig — kein neues
      Target, kein neues Gate. Der Implementer bestätigt das oder begründet
      eine Abweichung im Plan-Nachzug.
- [ ] `AGENTS.md` §4/`harness/README.md` §Sensors bleiben unverändert —
      geprüft, nicht nur angenommen (der Slice verschärft kein bestehendes
      Gate; die `structure`-Regel ist eine Register-Invariante, keine neue
      Schwelle).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — hier
      geprüft **von der Slice-Closure selbst**: die Roadmap führt derzeit keine
      offene Welle, es gibt also keine Welle-Closure, die sie einsammeln
      könnte (Baseline-Regelwerk `modul-06-roadmap.md` §Was der wellenlose
      Betrieb selbst auslöst).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/e2e-abdeckung.md` | neu (erzeugt) | das Erzeugnis; Erststand aus einem realen `make test-integration`-Lauf, danach nur noch vom Lauf fortgeschrieben |
| `test/integration/integration_test.go` | update | Erzeuger der Go-Hälfte: leitet Funktionsname, Spec-Kennung, Kurzbeschreibung und Quellzeile der `func TestE2E*` per `go/parser` aus der eigenen Datei ab, bricht ohne Spec-Kennung sichtbar ab, gibt die Zeilen in fester Reihenfolge auf stdout aus; läuft als eigener `-run`-Aufruf |
| `tools/harness/run-integration-tests.sh` | update | ruft den Erzeuger der Go-Hälfte auf, deklariert die Bash-Hälfte je Phase (Anker, Anzeigename, Spec-Kennungen, Kurzbeschreibung), prüft die Anker wörtlich, setzt beide Hälften in fester Reihenfolge zusammen und schreibt `docs/user/e2e-abdeckung.md` nur bei inhaltlicher Abweichung |
| `.d-check.yml` | update | `structure`-Regel für die Tabelle (Abschnitt + Spalten-Mindestbreiten) — der `ids`-Abschnitt bleibt unverändert, er trägt die Linkpflicht schon |
| `harness/sensors/docs-check.md` | update | neue Grenzen (keine Symbol-Prüfung; `codepaths` aus → Folge-Vorgang; Bash-Hälfte deklariert statt abgeleitet) und der Ist-Zustand der `structure`-Regeln |

**Nachzug aus der Implementierung** (`implement-slice.md` Schritt 14) — was
über diesen Plan hinausging oder von ihm abwich:

| Punkt | Was gilt | Begründung |
|---|---|---|
| Go-Hälfte liest **das Paket**, nicht nur die eigene Datei | Der Erzeuger parst alle `*.go` neben der Testdatei (Reihenfolge Datei, dann Quellzeile) | Eine `func TestE2E*` in einer zweiten Datei desselben Pakets fiele sonst still aus der Tabelle — dieselbe Klasse wie ein `-run`-Muster, das eine Funktion nicht trifft (`BEO-PGC/test-runner-stiller-ausschluss`), die dieser Slice sonst neu erzeugte |
| Kurzbeschreibung: **erster Absatz, darin der erste Satz** | Satzende = erster Punkt außerhalb eines Inline-Code-Spans, dem Leerraum oder Textende folgt und dem keine Ziffer vorausgeht; danach Aufräum-Regeln für die Lücken, die entfernte Kennungen lassen | Der Plan ließ die Ableitung offen. Der Absatz-Bezug hält eine nummerierte Aufzählung aus der Zelle (eine reine Satzregel zöge den ersten Aufzählungspunkt mit herein); beide Regeln stehen als Kommentar am Erzeuger |
| **Dritter Abbruch-Wächter**: unpaarige Inline-Code-Zeichen | Bricht der Doc-Kommentar bzw. die Kurzbeschreibung ein Backtick-Zeichen übrig, endet der Erzeuger sichtbar | Beim Bauen real gefeuert (der Kurzform-Tail wurde zunächst um ein Zeichen zu kurz entfernt): ein halbierter Span bleibt sonst stumm und formatiert den Rest der Zeile um |
| Kurzform-Auflösung deckt **`` `NNN` ``-umschlossene** Nachbarnummern | `` `…-002`…`006` `` und `` `…-003`/`004` `` lösen auf; eine Einleitung ohne dreistellige Ziffernfolge ist eine Auslassung im Fließtext und bleibt stehen | Die real vorgefundenen Formen; ein Abbruch bei jeder Auslassung träfe auch Prosa |
| Die `structure`-Regel ist die **achte** des Configs | Zählung folgt den vorhandenen Regel-Kommentaren; die Mindestbreiten binden an die Darstellung | Der Plan nennt keine Nummer |
| `section-missing` bei fehlender Datei | Die Regel rotet, solange `docs/user/e2e-abdeckung.md` nicht im Baum liegt („Regel trifft keine Datei — das Gate liefe leer") | Real gemessen: die Datei muss committet sein, das Erzeugnis ist kein Nur-Lauf-Artefakt |
| **Nicht-Realisierung** der Aussage über die fünfte `structure`-Regel | `harness/sensors/docs-check.md` trägt die Aussage „auskommentiert bis zur ersten Closure" **nicht** (gemessen mit `grep`); `.d-check.yml` führt die Regel bereits als „AKTIVIERT mit der ERSTEN Closure" | Der Plan beschrieb einen Stand, den `slice-075` bereits nachgezogen hat. Statt einer Korrektur nennt die Sensordoku jetzt **den Ist-Zustand aller `structure`-Regeln** — dieselbe Pflicht, anderer Träger |

**Nicht in dieser Liste, mit Begründung:** `harness/README.md` §Sensors
(kein neues Target, kein neues Gate — die Tabelle ist ein Erzeugnis der
bestehenden E2E-Kette), `d-check.mk` (kein neues Target), `AGENTS.md` (keine
neue Hard Rule; die Erzeuger-Disziplin ist Implementer-Wissen und steht im
erzeugenden Code und in dieser Plan-Datei — verkörpert wird sie erst, wenn der
Zähler sie verlangt, §7). Beide bestätigt-gemessen: `git diff --name-only`
führt weder `harness/README.md` noch `AGENTS.md` noch `d-check.mk`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei (1 je Rolleninhaber der
Implementer-Rolle) — `in-progress/` trägt keinen Slice; und eine Umgebung, in
der die Compose-Kette real läuft (Docker, `make test-integration` aufrufbar):
die Belege (a)–(c) sind reale Läufe, kein Review-Argument.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass die
  Bash-Hälfte ohne ein zweites Werkzeug nicht mechanisch prüfbar wird (also
  Anker-Prüfung und Zusammensetzung im Runner nicht ohne fragile
  Textableitung gelingen), ist der Schnitt zu breit — dann zurück zur
  Zerlegung in *Go-Hälfte* und *Bash-Hälfte* als zwei Slice.
- `in-progress` → `open` (blockiert — Carveout?): Erweist sich, dass die
  Erzeugung ohne ein eigenständiges Programm nicht geht (die abgeleiteten
  Zeilen sind im Runner nicht zusammensetzbar, ohne den Runner zu einem
  zweiten Werkzeug umzubauen), kollidiert das mit der Auflage *kein eigenes
  Skript*: Blocker mit Carveout-Entscheidung, nicht stiller Umbau.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** die drei realen Belege
geführt: (a) ein `make test-integration`-Lauf erzeugt die Datei, (b) ein
zweiter Lauf mit unverändertem Testbestand lässt sie inhaltsgleich, (c) eine
real hinzugefügte und eine real entfernte `func TestE2E*` ändern sie in beide
Richtungen. Dazu: **der committete Stand entspricht der Ausgabe eines frischen
Laufs** (die Datei ist Erzeugnis, kein handgepflegter Nachzug) **und** die
Closure-Notiz ist geschrieben.

Kein Ablaufdatum. Die Erzeugung hängt an der Umgebung, nicht am Kalender.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Bash-Hälfte ist deklariert, nicht abgeleitet.** Der Anker-Check deckt
  nur die Vorwärtsrichtung (*„die Tabelle behauptet eine Phase, die es nicht
  gibt"*); eine neue Runner-Phase, die niemand deklariert, fehlt still in der
  Tabelle — dieselbe Klasse, die `BEO-PGC/test-runner-stiller-ausschluss` für
  die `-run`-Muster führt. — **Ausgang:** <bei Closure zuzuweisen>
- **Kurzbeschreibungen gehen durch eine Normalisierung.** Der Doc-Kommentar
  einer Testfunktion nennt Kennungen (teils in Kurz- und Bereichsform) und
  trägt Inline-Code-Spans; roh in eine Tabellenzelle kopiert, kann sein Text
  die Spans unbalanced machen und Kennungen außerhalb einer Spans setzen —
  dann meldet `ids` `id-unlinked`. Der Erzeuger entfernt Kennungen deshalb. Die
  Aussage, die nur im Kennungsverweis des Kommentars lag, steht danach nicht
  mehr in der Tabelle — erreichbar bleibt sie über die Spalte `Ort` an ihrer
  Quelle. — **Ausgang:** <bei Closure zuzuweisen>
- **Die Kurzformen der Doc-Kommentare.** Die E2E-Kommentare nennen Kennungen
  teils in Kurzform (`LH-FA-CAP-001`…`003` als Bereich, `LH-FA-CFG-003`/`004`
  als Nachbarnummern); der Erzeuger muss sie auflösen, sonst ist die Tabelle
  unvollständig. Die Antwort ist *sichtbarer Abbruch* statt Raten — ob sie
  alle im Repo benutzten Formen trägt, entscheidet der reale Lauf. —
  **Ausgang:** <bei Closure zuzuweisen>
- **`Ort` bindet an Quellzeilen.** Jede Änderung oberhalb einer Testfunktion
  verschiebt die Zeilennummer und damit den Diff der Datei; das ist gewollt
  (erzeugt, nicht handgepflegt), erzeugt aber einen Doku-Diff in PRs, die nur
  Testcode ändern. — **Ausgang:** <bei Closure zuzuweisen>
- **Die Belege brauchen mehrere volle Compose-Läufe**, während ein anderer
  Vorgang parallel arbeitet; der Rot-Beleg (c) verlangt einen realen Eingriff
  in `test/integration/**` und seine reale Rücknahme. Ein Eingriff, der im
  Arbeitsbaum liegen bleibt, wandert in einen fremden Commit. — **Ausgang:**
  <bei Closure zuzuweisen>

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <bei Closure>
- **Was ging anders als geplant:** <bei Closure>
- **Steering-Loop-Eintrag:** <bei Closure — erwartet: die drei neuen
  `docs-check`-Grenzen sind benannte Spec-/Sensor-Grenzen, kein neuer Zähler;
  verkörpert wird in diesem Slice nichts außerhalb von
  `harness/sensors/docs-check.md`>.
- **Beobachtungs-Register (`../observations/`):** <bei Closure — erwartet:
  kein neues Verzeichnis; `evidence/slice-074.md` in
  `BEO-PGC/test-runner-stiller-ausschluss` als **zweiter Träger derselben
  Klasse** (die Deklarations-Hälfte der Tabelle; Präzedenz: slice-011s zweiter
  Träger in `BEO-PGC/lese-doppelquelle`), Zähler dann 2× — weiter unter der
  Schwelle, Ausgang bleibt *weiter offen*. Keine Beobachtung angefallen ist
  ebenfalls eine Antwort und wird hier notiert.>
- **Folge-Slices:** keine — die `codepaths`-Aktivierung ist in §1 als eigener
  **Vorgang** benannt (Gate-Verschärfung mit eigenem ADR-Pfad), nicht als
  Slice-Kennung zugesagt; eine Kennung, die es nicht gibt, wäre eine leere
  Adresse.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** hier geprüft — die Roadmap führt derzeit keine offene
  Welle, die Prüfung trägt damit die Slice-Closure selbst (Baseline-Regelwerk
  `modul-06-roadmap.md` §Was der wellenlose Betrieb selbst auslöst).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md`
§Modus-Deklaration pro Sub-Area) — sie deckt das gesamte Repo und damit auch
`test/integration/**`, `tools/harness/**`, `.d-check.yml` und `docs/user/**`.
Eine feinere Sub-Area ist nicht zu bilden: „Test-Infrastruktur" ist keine
deklarierte Sub-Area, und das Register führt den Runner unter demselben Kürzel
(`BEO-PGC/test-runner-stiller-ausschluss`, dort als „Test-Infrastruktur
(`tools/harness/run-integration-tests.sh`; Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)"). Kein Differenzierungs-Bedarf, kein Block mit vermischten
Modi.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; vier Einträge mit
Bezug zu diesem Vorhaben:

- `BEO-PGC/lese-doppelquelle` (3×, **verkörpert**) — „zwei Quellen für
  denselben Zustand" ist genau die Gefahr einer handgepflegten Tabelle. Der
  Slice antwortet auf die Klasse, ohne sie erneut zu berühren: die Tabelle
  entsteht aus **einer** Quelle (`Nachweis-Deklaration ↔ Zeile`) statt neben
  einer zweiten. Die Beobachtung gehört der Sub-Area *Store-Adapter/Lesepfad*,
  die dieser Slice nicht anfasst — **kein neuer Beleg**; die Übernahme des
  Musters ist keine neue Berührung.
- `BEO-PGC/d-migrate-nacharbeit` (5×, verkörpert) — erzeugte Artefakte, die an
  einer Werkzeug-Grenze eine Nacharbeit tragen: der Präzedenzfall dafür, dass
  ein Erzeugnis seine Ausweichform **benennen** muss. Die Grenze dieses
  Erzeugnisses ist die deklarierte Bash-Hälfte und die fehlende
  Symbol-Prüfung — sie steht als §6-Risiko und als Grenze in
  `harness/sensors/docs-check.md`. **Kein neuer Beleg**, kein neues Auftreten.
- `BEO-PGC/test-runner-stiller-ausschluss` (1×, **offen**) — die direkteste
  Berührung: eine Testfunktion außerhalb jedes `-run`-Musters läuft still nie.
  Für eine *erzeugte* Tabelle heißt das: käme der Zeilensatz aus dem, **was
  gelaufen ist**, verschwände eine exkludierte Funktion spurlos aus der
  Tabelle — der Erzeuger leitet den Zeilensatz deshalb aus dem **Quelltext**
  ab, und eine exkludierte Funktion erscheint weiterhin. Der
  Vollständigkeits-Sensor über die `-run`-Muster ist ausdrücklich **nicht**
  Teil dieses Slice (§1). Bei Closure: **zweiter Träger derselben Klasse** (die
  Deklarations-Hälfte der neuen Tabelle, §6 Risiko 1) als
  `evidence/slice-074.md` — Präzedenz: der zweite Träger in
  `BEO-PGC/lese-doppelquelle` (slice-011); Zähler dann 2×, weiter unter der
  Schwelle, Ausgang bleibt *weiter offen*.
- `BEO-PGC/slice-chronik-in-code-kommentar` (3×, verkörpert) — berührt diesen
  Slice an einer benennbaren Stelle: Erzeuger-Kommentare und Doku-Prosa
  beschreiben den Ist-Zustand, keine Entstehungsgeschichte (`AGENTS.md` §3.7),
  und die Tabelle trägt **keine** Slice-/Wellen-Chronik. Die verkörperte
  Instruktion (`.harness/skills/reviewer.md`, `implement-slice.md` Schritt 20)
  wirkt im Implementer-Lauf — **kein neuer Beleg**, kein neues Auftreten.

Kein weiterer Eintrag des Registers berührt die Sub-Area dieses Slice; kein
Eintrag erreicht mit diesem Slice die Schwelle 3×, und es wird kein neues
Verzeichnis angelegt. Die Antwort je Eintrag steht oben — *keine Treffer* wäre
hier falsch gewesen, `test-runner-stiller-ausschluss` ist ein Treffer.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`) — kein Modus-Begründungsblock.
Die vier Pflichtkriterien sind für diesen Slice dennoch beantwortet, weil die
GF-Setzung hier trägt und nicht nur behauptet: **Konventionen-Dichte** hoch
(`harness/conventions.md` MR-000, `.d-check.yml`-Module, Sensors-Bindung in
`harness/README.md`), **Phase-Reife** hoch (drei Spec-Straten, ADR-Index,
Planning-Lifecycle und 73 geschlossene bzw. laufende Slices liegen vor),
**Evidenz-/Diskrepanz-Risiko** niedrig für dieses Erzeugnis (die Aussage wird
nicht gegen einen Code-Bestand inventarisiert, sondern aus ihm *erzeugt*;
das Restrisiko ist die deklarierte Bash-Hälfte, §6), **Reconciliation-Aufwand**
null — kein Brownfield-Bestand, kein Graduation-Trigger.
