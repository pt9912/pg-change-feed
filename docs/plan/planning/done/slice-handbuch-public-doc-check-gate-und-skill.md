# Slice handbuch-public-doc-check-gate-und-skill: Ein Gate für kennungsfreie Nutzerdokumentation, die ADR dazu und ein kleiner Skill „Nutzerdokumentation schreiben“

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

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (Dokumentation für
Betreiber; Haupt-Bezug, die Nutzerdokumentation ist die Betreiber-/Integrator-Sicht),
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
(`Accepted`; die Entscheidung dieses Slice, Ergebnis des
[Architect-Verdikts](../../../reviews/architect-verdict-handbuch-public-doc-check-gate-und-skill.md)),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (Vorbild: ein
netzloses `grep`-Gate für „keine interne Kennung“, hier für `docs/user/`),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen: der Skill-Inhalt „Ursprungs-Wort statt Berichts-Kennung“),
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) (`Accepted`, unberührbar;
fordern die Zeile in der Änderungshistorie des Handbuchs — die ADR dieses Slice
schärft sie, ändert sie nicht). Anlass: das Benutzerhandbuch trug interne Kennungen
und Chronik; die Bereinigung ist committet
(`7d1c4611`, `de0d899a`, `983e209a`), ein Wächter, der den Rückfall fängt, fehlt.

**Berührte Spec-Stellen:** keine geändert. Gelesen:
[`spec/pflichtenheft.md`](../../../../spec/pflichtenheft.md) Zeile 1075 (Beispiel-Zeile
„die Änderungshistorie des Handbuchs trägt die Zeile“) — bleibt wahr, weil das Kapitel
bleibt; ein Träger, der nicht mitzuändern ist, im Suchlauf in §3 geführt.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, auf Freigabe des Auftraggebers (zwei Aussagen, wörtlich):
(1) Auf die Frage „Brauchen wir dafür eine eigene Rolle, Command oder Skill
(Handbuch-Schreiber)?“ lautete der Vorschlag: zuerst ein Gate (Muster
`make sdk-public-doc-check`, [`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md)), dann ein kleiner Skill statt einer Rolle, eine
Rolle nur bei erneutem Drift; Antwort des Auftraggebers: „ja, nachdem die anderen offenen
Slices abgearbeitet sind“ (sie sind abgearbeitet). (2) „Die History kann ohne Kennung
auskommen“ — das Kapitel `### Änderungshistorie` bleibt, ohne Kennungen.
**Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** Drei zusammengehörige Teile, weil sie dieselbe Regel tragen („Nutzerdokumentation
ist Betreiber-/Integrator-Sicht und trägt keine interne Kennung“) und sich gegenseitig
brauchen — das Gate fängt, der Skill verhindert, die ADR legt fest, was gilt:

- **(A) Ein Gate.** Ein netzloser, `grep`-basierter Wächter, Teil von `make gates`, der die
  Nutzerdokumente unter `docs/user/` (Reichweite: unten) auf interne Kennungen und auf Links
  in `docs/reviews/` und `docs/plan/` prüft. Form wie `sdk-public-doc-check`:
  Skript `tools/harness/<name>.sh`, Sensor-Vertrag `harness/sensors/<name>.md`,
  Make-Ziel mit `GATE_CHECKS +=` (Heimat: `harness/mk/doc-gate.mk`, wo die
  Dokumentations-Gates liegen; der Architect bestätigt), Tabellentest
  `make test-<name>` (Werkzeug, kein Gate — wie `test-sdk-public-doc-check`), Zeile in
  `harness/README.md` §Sensors (Gate-Index, [`AGENTS.md`](../../../../AGENTS.md) §4) und in
  der `make gates`-Aufzählung. Arbeitsname `handbuch-public-doc-check`; der endgültige
  Name ist Teil der Architect-Entscheidung.
- **(B) Eine ADR.** Eine neue Schwelle ist nach [`AGENTS.md`](../../../../AGENTS.md) §3.6 eine
  Entscheidung, kein PR-Kommentar. Die ADR (Muster `ADR-0134`) legt fest: Reichweite,
  Muster, Ausnahmen, Gate-Aufnahme, Beziehung zum Standard und zu `ADR-0087`/`-0088`/`-0090`
  (`Schärft:` — die Accepted-ADRs bleiben unberührt, §3.5).
- **(C) Ein kleiner Skill und die Umstellung zweier Prüfpunkte.** Ein Skill
  „Nutzerdokumentation schreiben“ (Ort und Inhalt: unten) und die Umstellung der beiden
  Stellen, die die Handbuch-Historie verlangen: `.claude/commands/implement-slice.md`
  Schritt 17 und `.harness/skills/reviewer.md` (HIGH-Punkt „Handbuch-Versionshistorie
  nicht fortgeschrieben“). **Die Pflicht „Version hochzählen + Zeile in der
  Änderungshistorie“ bleibt**; sie wird ergänzt um „ohne Kennungen, Betreibersicht“ und den
  Verweis auf Skill und Gate.

**Die Architect-Entscheidung liegt vor** (Startbedingung §4 erfüllt):
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) ist
`Accepted`, das [Verdikt](../../../reviews/architect-verdict-handbuch-public-doc-check-gate-und-skill.md)
nennt die Entscheidungen Q1 bis Q6. Die Vorschläge unten sind Eingang der Entscheidung; wo
sie abweichen, gilt die ADR (Reichweite: Hybrid, Muster `P`/`L`, keine Ausnahme).

**Vorschläge des Planners als Eingang für den Architect (gemessen am Stand `7e993efd`,
2026-10-02, Zahlen im Suchlauf-Feld §3):**

1. *Reichweite.* Jede Datei unter `docs/user/` wurde per `git grep` auf Kennungen
   (`LH-(FA|QA)-`, `ADR-[0-9]`, `SPEC-[0-9]`, `ARC-[0-9]`, `slice-`, `welle-`, `BEO-`,
   `MR-[0-9]`, `CO-[0-9]`) eingestuft:

   | Datei | Treffer | Einstufung (Vorschlag) |
   |---|---|---|
   | `benutzerhandbuch.md` | 0 | **im Gate** — Nutzerdokumentation, Gegenstand des Anlasses |
   | `benutzerhandbuch-standard.md` | 0 | **im Gate** — Nutzerdoku (kein Record), trägt die Form des Handbuchs; kostet nichts, schützt den Standard vor Rückfall |
   | `version.md` | 0 | **im Gate** (Entscheidung der ADR: kostet nichts, das Handbuch verweist auf die Datei) |
   | `releasing.md` | 26 (7 Links nach `../plan/adr/`) | **nicht im Gate** — Betreiber-/Maintainer-Doku zum Release-Prozess, die ihre Entscheidungen ([`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md), `-0106` …) und eine kennungsgespickte Änderungshistorie bewusst trägt; entschieden: dauerhaft ausgenommen, kein Folge-Slice für eine Bereinigung |
   | `e2e-abdeckung.md`, `bench-abdeckung.md`, `ci-matrix-abdeckung.md`, `sdk-e2e-abdeckung.md` | 74 · 4 · 3 · 21 | **nicht im Gate** — von Runnern geschriebene Erzeugnisse; die Kennungen *sind* ihr Inhalt (Spec-Kennung je Zeile, RTM-Eingabe von `make doc-trace`) |

   **Entschieden: Hybrid.** Geprüft werden `benutzerhandbuch.md`, `benutzerhandbuch-standard.md`
   und `version.md`; ausgenommen sind `releasing.md` und die vier `*-abdeckung.md`. Jede andere
   `*.md` unter `docs/user/` und jede genannte, aber fehlende Datei färbt den Lauf rot
   (Exit 2, Vollständigkeitsprüfung). Damit werden beide Fehlerbilder laut: die neue
   Nutzerdatei bleibt nicht lautlos unbewacht, die neue Erzeugnis-Datei färbt nicht
   unbemerkt falsch rot (`sdk-public-doc-check` Grenze 2), sondern erzwingt die Klassifikation.
2. *Muster.* Kennungen wie oben; dazu Links, deren Ziel `docs/reviews/` oder `docs/plan/`
   ist (Muster `docs/(reviews|plan)/` und relative Formen `\.\./(reviews|plan)/`). Gemessen:
   0 solcher Links im Handbuch (`git grep`, §3). Entschieden (ADR, Festlegung 3): ERE (kein
   `grep -P`) mit geschärftem Wortrand — `slice-`/`welle-` treffen nur am Zeilenanfang oder
   nach einem Zeichen außerhalb von Buchstabe, Ziffer, `_` und `-` (`byte-slice-x` trifft
   nicht), `MR-`/`CO-` verlangen drei Ziffern; die Muster heißen dort `P` (Kennungen) und
   `L` (Links) und stehen wörtlich im Suchlauf in §3.
3. *Ausnahmen / Befehlsbeispiele.* Im Handbuch tragen Befehlsbeispiele keine Kennung
   (0 Treffer gemessen; Tag-Muster wie `sdk-csharp-v*` treffen keines der Muster). Vorschlag:
   **keine Ausnahme**, auch nicht in Fenced-Blöcken — ein Handbuch-Beispiel mit Kennung ist
   selbst der Fehler. **Entschieden: keine Ausnahme**; eine Ausnahmeform (Opt-out-Marker)
   braucht eine eigene ADR. **Folge für den Implementer:** der neue Satz in
   `benutzerhandbuch-standard.md` zur Kennungsfreiheit nennt **keine Beispiel-Kennung**
   („keine internen Kennungen“, nicht „z. B. …“ mit einer Kennung), sonst färbt er das
   Gate rot.
4. *Beziehung zum Standard.* [`benutzerhandbuch-standard.md`](../../../user/benutzerhandbuch-standard.md)
   Zeile 229 („Changelog oder Änderungshistorie“) und Zeile 350 (`## 11. Änderungshistorie`)
   verlangen das Kapitel — es **bleibt**. Zu ergänzen ist an beiden Stellen die
   Kennungsfreiheit (Historie ohne Kennungen, Betreibersicht). Das ist Nutzerdoku, kein
   Record: die Anpassung gehört in den Slice (Liefer-Punkt 4).
5. *Beziehung zu `ADR-0087`/`-0088`/`-0090`.* Sie fordern die Zeile in der
   Änderungshistorie (Fundstellen im Suchlauf, §3); die Forderung bleibt wahr. Die neue ADR
   schärft sie um „ohne Kennungen“ (`Schärft:`), keine Supersession.
6. *Gate-Fläche und CI.* `make gates` läuft in `ci.yml` unverändert; kein Workflow ändert
   sich, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift deshalb nicht (am Start zu messen,
   §6).
7. *Skill-Ort.* Bestand: `.harness/skills/` trägt `reviewer.md` und
   `closure-note-reviewer.md` (beide mit Kopf „Status/Bezug/Gilt für“), `.claude/` trägt
   Agenten, Commands, Hooks und Rules, **kein** `skills/`-Verzeichnis. Vorschlag:
   `.harness/skills/nutzerdoku-schreiben.md` (gleiche Form wie die beiden Vorgänger),
   eingetragen in `harness/README.md` §Guides (dort steht der Reviewer-Skill) und aus
   `implement-slice.md` Schritt 17 sowie dem Reviewer-Skill verlinkt. Der Skill ersetzt
   **nicht** den Schritt 17, er trägt den Inhalt, den Schritt 17 heute inline führt
   (siehe Risiko Zeilen-Gate, §6).

**Skill-Inhalt (Vorgabe, Form-Prüfung im Review):** Ist-Zustand statt Chronik; keine
internen Kennungen und keine Links nach `docs/plan/`/`docs/reviews/`;
Betreiber-/Integrator-Sicht (was tut der Nutzer, nicht wie ist es gebaut); bei Messungen
das Ursprungs-**Wort** (gemessen/übernommen/abgeleitet, [`AGENTS.md`](../../../../AGENTS.md)
§3.12) und den Lauf-Hinweis **ohne** Berichts-Kennung; Version hochzählen und je Version eine
Zeile in `### Änderungshistorie`, in Betreibersicht und ohne Kennungen; Verweis auf das Gate
als Fangnetz („das Gate liest Kennungen, nicht Sinn“). **Der Skill ändert `AGENTS.md`
nicht** (Regeln dort sind Hard Rules; der Skill verweist auf sie, ergänzt keine).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Rolle „Handbuch-Schreiber“** — der Auftraggeber hat „nur bei erneutem Drift“
  gesagt; Beleg wäre ein weiterer Rückfall trotz Gate und Skill, nicht jetzt.
- **Eine Änderung des Handbuch-Inhalts** — das Handbuch ist bereinigt (0 Treffer am Stand
  `7e993efd`); der Slice berührt `benutzerhandbuch.md` nicht und zählt keine Version hoch.
  Findet das Gate am echten Baum einen Treffer, den die Bereinigung übersah, geht er als
  Befund an den Planner, keine stille Mitänderung.
- **`releasing.md` und die Erzeugnisse bereinigen** — andere Dokumentklasse (Vorschlag 1);
  gemeldet, nicht geändert.
- **`AGENTS.md`, `spec/` oder Accepted-ADRs ändern** — `AGENTS.md` ist Hard-Rules-Quelle,
  die Spec-Zeile 1075 bleibt wahr, Accepted-ADRs sind nach §3.5 unberührbar (die neue ADR
  schärft).
- **Ein Release oder eine Versionsänderung** — kein Tag, keine Versionsdatei.
- **Ein neuer Workflow** — das Gate hängt in `GATE_CHECKS`, `ci.yml` fährt es mit.

## 2. Definition of Done

- [x] **Liefer-Punkt 0 — Architect-Entscheidung.** Erledigt:
      [`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
      liegt `Accepted` vor (Reichweite, Muster, Ausnahmen, Gate-Aufnahme, Standard Z. 229/350,
      `Schärft:` `ADR-0087`/`-0088`/`-0090`, Skill-Ort), samt
      [Architect-Verdikt](../../../reviews/architect-verdict-handbuch-public-doc-check-gate-und-skill.md)
      und Index-Zeile in `docs/plan/adr/README.md`. Die ADR und ihr Index-Eintrag gehören
      zum Diff des Architect-Zugs, nicht zu dem des Implementers.
- [x] **Liefer-Punkt 1 — Das Gate.** Skript, Sensor-Vertrag, Make-Ziel samt
      `GATE_CHECKS +=` (in `harness/mk/doc-gate.mk`), Tabellentest samt Werkzeug-Ziel
      `make test-handbuch-public-doc-check` (kein Gate) und `harness/README.md` nach der
      ADR: Gate-Zeile in §Sensors, **Zeile des Tabellentests** (Werkzeug, kein Gate) und die
      `make gates`-Aufzählung. Reichweite ist der **Hybrid**
      ([`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Festlegung 1):
      geprüft `benutzerhandbuch.md`, `benutzerhandbuch-standard.md`, `version.md`;
      ausgenommen `releasing.md` und die vier `*-abdeckung.md`; jede andere `*.md` unter
      `docs/user/` und jede genannte, aber fehlende Datei endet mit Exit 2.
      Der Tabellentest deckt je Kennungsart einen Treffer (Exit 1), je Link-Klasse
      (`docs/reviews/`, `docs/plan/`, relativ) einen Treffer, eine saubere Datei, die
      Wortrand-Fälle ohne Treffer (`byte-slice`, `CO-2 Emissionen`, `Slice-1`), eine
      unklassifizierte `.md` (Exit 2), eine fehlende genannte Datei (Exit 2) und die
      ausgenommenen Dateien (Exit 0). *Zu belegen
      durch:* (a) `make test-<name>` Exit 0, die gedruckte Schlusszeile im Bericht;
      (b) **Mutationsprobe**: eine Kennung (z. B. `ADR-0134`) in eine **Kopie** von
      `docs/user/benutzerhandbuch.md` im Scratchpad einfügen (Edit/Write, nie am
      Arbeitsbaum, [`AGENTS.md`](../../../../AGENTS.md) §3.1) und das Skript mit dem
      Kopie-Verzeichnis als Wurzel aufrufen → Exit 1 mit gedruckter Zeile; ebenso ein Link
      nach `docs/plan/`; Rücknahme = Kopie verwerfen; Stelle, Instanz und Farbe im Bericht
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12); (c) **grüner Lauf am echten Baum**:
      `make <name>` Exit 0 ohne Treffer; (d) `make gates` Exit 0 **mit** dem neuen Gate
      (Exit-Code ungefiltert gesichert und gesondert ausgewertet, §3.9), das Gate steht in
      der Aufzählung der Lauf-Ausgabe.
- [x] **Liefer-Punkt 2 — Der Skill.** `.harness/skills/nutzerdoku-schreiben.md` (Ort laut
      ADR) existiert, trägt den Inhalt aus §1, ändert `AGENTS.md` nicht und ist aus
      `harness/README.md` §Guides, `implement-slice.md` Schritt 17 und dem Reviewer-Skill
      verlinkt. *Zu belegen durch:* `test -f`, `git grep -n nutzerdoku-schreiben` nennt die
      drei Verweise, `git diff --name-only <Basis> -- AGENTS.md` ist leer, `make docs-check`
      Exit 0 (Links lösen auf).
- [x] **Liefer-Punkt 3 — Umstellung der zwei Prüfpunkte (Regeltexte).** Diese Änderungen
      an Regeltexten sind **durch die Freigabe des Auftraggebers gedeckt** („ja, nachdem …“
      für Gate, ADR, Skill und Umstellung der Prüfpunkte) und im Slice als solche benannt:
      `.claude/commands/implement-slice.md` Schritt 17 und `.harness/skills/reviewer.md`
      (HIGH-Punkt „Handbuch-Versionshistorie nicht fortgeschrieben“). Beide behalten die
      Pflicht „`Version:` hochzählen **und** eine Zeile in `### Änderungshistorie`“,
      ergänzt um „ohne Kennungen, Betreibersicht“ und den Verweis auf Skill und Gate;
      keine Abschwächung der Pflicht. Der Kandidatenlauf wandert mit dem Inhalt in den Skill
      (beide Befehle stehen dort vollständig); Schritt 17 verweist auf ihn. *Zu belegen
      durch:* `git diff` beider Dateien im Bericht; im Diff steht die Pflicht weiter
      (`git grep -n "Änderungshistorie" -- .claude .harness/skills` hat weiterhin je einen
      Treffer), die Ergänzung und der Verweis sind ergänzt; Zeilenzahl von
      `implement-slice.md` nach dem Eingriff **nicht größer** als am Start (Risiko
      Zeilen-Gate, §6).
- [x] **Liefer-Punkt 4 — Der Standard.** `docs/user/benutzerhandbuch-standard.md` Zeile 229
      und Kapitel `## 11. Änderungshistorie` (Zeile 350) tragen die Kennungsfreiheit der
      Historie (Betreibersicht), das Kapitel bleibt. *Zu belegen durch:* `git diff` der Datei;
      `git grep -nF "Änderungshistorie" -- docs/user/benutzerhandbuch-standard.md` liefert
      weiter zwei Treffer; das neue Gate läuft grün über die Datei.
- [x] **Nur diese Pfade.** Der Diff berührt `docs/user/benutzerhandbuch.md`, `AGENTS.md`,
      `spec/` und Accepted-ADRs nicht. *Zu belegen durch:*
      `git diff --name-only <Basis> -- docs/user/benutzerhandbuch.md AGENTS.md spec` ist leer;
      für die ADRs `0087`/`0088`/`0090` ebenfalls.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert
      ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-handbuch-public-doc-check-gate-und-skill.md` liegt vor,
      kein offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8
      des Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review
      (Modul 8). **Abweichung, ehrlich festgehalten:** der Report sah den Stand vor der
      Fixrunde `28a5aafe`; ein Re-Review danach fand nicht statt (§7).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-handbuch-public-doc-check-gate-und-skill.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (Zeile des neuen Gates, `make gates`-
      Aufzählung) und §Guides (Skill), ADR-Index; das Benutzerhandbuch bleibt unberührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7). Kandidat des Lerneintrags: ein
      Rückfall in Nutzerdoku war bisher nur über Lesen (Reviewer) gefangen — das Gate ist der
      erste Sensor für diese Klasse; die benannte Grenze („liest Kennungen, nicht Sinn“) ist
      Teil des Eintrags.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo
      (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7) — kein Anfall ist
      ebenfalls eine Antwort und wird in §7 notiert (Nachbarn:
      `BEO-PGC/handbuch-versionshistorie-uebersprungen`,
      `BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

**Umfang:** S bis M — Schätzung, nicht gemessen: ein Skript, ein Tabellentest, ein
Sensor-Vertrag, eine ADR, ein Skill, vier Textstellen in Regeln/Standard.

**Voraussetzung:** Die ADR (Liefer-Punkt 0) ist `Accepted` (erfüllt:
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)).
Der Skill nennt die Muster `P`/`L` nicht wörtlich (sie gehören in den Sensor-Vertrag). Kein weiterer offener Slice mit demselben Gegenstand liegt vor
(`slice-meldungscodes-statt-interner-kennungen` in `open/` betrifft Meldungstexte des Codes,
nicht Nutzerdoku — am Start zu bestätigen).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0143-…md`, `docs/plan/adr/README.md` | erledigt (Architect-Zug, `30a060c9`) | Architect-Entscheidung (Liefer-Punkt 0); `Schärft:` `ADR-0087`/`-0088`/`-0090`; nicht Teil des Implementer-Diffs |
| `tools/harness/handbuch-public-doc-check.sh` (Name laut ADR) | neu | `grep`-Wächter nach dem Muster von `sdk-public-doc-check.sh`; optionales erstes Argument = Wurzel für den Tabellentest |
| `tools/harness/run-handbuch-public-doc-check-tests.sh` | neu | Tabellentest (Muster `run-sdk-public-doc-check-tests.sh`) |
| `harness/mk/doc-gate.mk` | update | Ziele `…-check` und `test-…`, `GATE_CHECKS +=` (Heimat laut ADR) |
| `harness/sensors/handbuch-public-doc-check.md` | neu | Vertrag, Ausgänge, Grenze, Sperren, Bindung (Muster `sdk-public-doc-check.md`) |
| `harness/README.md` | update | §Sensors: Gate-Zeile, `make gates`-Aufzählung, Zeile des Tabellentests; §Guides: Skill |
| `.harness/skills/nutzerdoku-schreiben.md` (Ort laut ADR) | neu | der Skill |
| `.claude/commands/implement-slice.md` Schritt 17 | update (Regeltext) | Pflicht bleibt, Ergänzung + Verweis; Inhalt zum Skill verlagern, Zeilen nicht mehren |
| `.harness/skills/reviewer.md` (HIGH-Punkt, Z. ~75) | update (Regeltext) | dieselbe Ergänzung + Verweis |
| `docs/user/benutzerhandbuch-standard.md` Z. 229, 350 | update (Nutzerdoku) | Kennungsfreiheit der Historie; Kapitel bleibt |
| `docs/user/benutzerhandbuch.md`, `AGENTS.md`, `spec/` | prüfen, nicht ändern | siehe „Nur diese Pfade“ |
| `.github/workflows/*` | prüfen | keine Änderung erwartet (§3.10) |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Nutzerdokumentation trägt
keine Kennung; die Historie ist kennungsfrei“; Zählwort „Kennung“, Symbolname
„Änderungshistorie“; Parent ist `30a060c9` (Zeilen mit Muster `P`/`L` der ADR und die
ADR-Zeile; die übrigen Zeilen bleiben bei `7e993efd`, deren Gegenstand sich nicht
bewegt hat), nachgemessen am 2026-10-02 mit `make suchlauf-nachmessen`; die `diff`-Zeilen
und die Befunde trägt der Implementer nach; neue Dateien sind für den Stand `diff` mit
`git add` im Index):**

```suchlauf
30a060c9 0 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/benutzerhandbuch.md
30a060c9 0 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/benutzerhandbuch-standard.md
30a060c9 0 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/version.md
30a060c9 26 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/releasing.md
30a060c9 4 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/bench-abdeckung.md
30a060c9 3 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/ci-matrix-abdeckung.md
30a060c9 74 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/e2e-abdeckung.md
30a060c9 21 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/sdk-e2e-abdeckung.md
30a060c9 0 -n -E 'docs/(reviews|plan)/|\.\./(reviews|plan)/' -- docs/user/benutzerhandbuch.md docs/user/benutzerhandbuch-standard.md docs/user/version.md
7e993efd 2 -n -F Änderungshistorie -- docs/user/benutzerhandbuch-standard.md
7e993efd 2 -n -F Änderungshistorie -- .claude .harness/skills
30a060c9 15 -n -F Änderungshistorie -- docs/plan/adr
7e993efd 1 -n -F Änderungshistorie -- spec
7e993efd 0 -n -F nutzerdoku-schreiben -- .
diff 0 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/benutzerhandbuch.md
diff 0 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/benutzerhandbuch-standard.md
diff 0 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/version.md
diff 26 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/releasing.md
diff 4 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/bench-abdeckung.md
diff 3 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/ci-matrix-abdeckung.md
diff 74 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/e2e-abdeckung.md
diff 21 -n -E '\b(LH-(FA|QA|RB)-|(ADR|SPEC|ARC)-[0-9]|BEO-[A-Za-z]|(MR|CO)-[0-9]{3})|(^|[^[:alnum:]_-])(slice|welle)-[a-z0-9]' -- docs/user/sdk-e2e-abdeckung.md
diff 0 -n -E 'docs/(reviews|plan)/|\.\./(reviews|plan)/' -- docs/user/benutzerhandbuch.md docs/user/benutzerhandbuch-standard.md docs/user/version.md
diff 2 -n -F Änderungshistorie -- docs/user/benutzerhandbuch-standard.md
diff 5 -n -F Änderungshistorie -- .claude .harness/skills
diff 15 -n -F Änderungshistorie -- docs/plan/adr
diff 1 -n -F Änderungshistorie -- spec
diff 11 -n -F nutzerdoku-schreiben -- .
```

| Träger | Messung am Parent (`7e993efd`, 2026-10-02) | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Kennungstreffer je Nutzerdokument (Zeilen 1–8) | Handbuch 0, Standard 0, `version.md` 0, `releasing.md` 26, vier Erzeugnisse (bench 4, ci-matrix 3, e2e 74, sdk-e2e 21) | Die Null-Zeilen sind die Bedingung für den grünen Lauf des Gates am echten Baum; `releasing.md` und die Erzeugnisse bleiben außerhalb der Reichweite (Vorschlag 1) und ihre Zahlen am Diff **unverändert**. Befund am Diff: Implementer. |
| Links nach `docs/plan/`/`docs/reviews/` (Zeile 9) | 0 im Handbuch und im Standard | Soll am Diff weiter 0. Die 7 Links in `releasing.md` und je einer in zwei Erzeugnissen liegen außerhalb. |
| `Änderungshistorie` im Standard (Zeile 10) | 2 (Z. 229, 350) | beide Stellen bekommen die Kennungsfreiheit (Liefer-Punkt 4); Soll am Diff weiter 2 (Kapitel bleibt). |
| `Änderungshistorie` in Regeln (Zeile 11) | 2 Trefferzeilen (`implement-slice.md` 1, `reviewer.md` 1) | Die Pflicht bleibt; keiner der beiden Treffer darf verschwinden. Befund am Diff: Implementer. |
| `Änderungshistorie` in den drei Accepted-ADRs (Zeile 12) | 8 (`0087` ×3, `0088` ×2, `0090` ×3) | **unberührbar** (§3.5); die 8 Trefferzeilen der drei Accepted-ADRs bleiben, dazu 7 in `ADR-0143`: Soll in Zeile 12 ist **15** (gemessen am Stand `30a060c9`, in dem die ADR liegt); am Diff bleibt 15. |
| `Änderungshistorie` in `spec/` (Zeile 13) | 1 (`pflichtenheft.md` Z. 1075, Beispiel-Zeile) | bleibt wahr, kein Nachzug; Soll am Diff weiter 1. |
| Skill-Name (Zeile 14) | 0 | Soll am Diff: Treffer in Skill-Datei, `harness/README.md`, `implement-slice.md`, `reviewer.md`, ADR (der Implementer trägt den gezählten Wert ein). |

**Befunde am Diff (Implementer, gemessen mit `make suchlauf-nachmessen`, 28 Zeilen stimmen,
Exit 0):** Gefunden: (1) Kennungstreffer: Handbuch, Standard und `version.md` je 0, die fünf
ausgenommenen Dateien unverändert 26 · 4 · 3 · 74 · 21; Links 0. (2) `Änderungshistorie` im
Standard bleibt 2 (die Ergänzung an Z. 229 steht in derselben Zeile, die an Kapitel 11 als
eigener Absatz ohne das Wort). (3) `Änderungshistorie` in `.claude`/`.harness/skills`: **5**
statt 2 — die zwei Pflicht-Treffer bleiben (`implement-slice.md` 1, `reviewer.md` 1), dazu
drei Zeilen im neuen Skill; kein Treffer verschwunden. (4) Accepted-ADRs und `ADR-0143`
15, `spec/` 1, unverändert. (5) Skill-Name: 7 Trefferzeilen außerhalb des Plans
(`implement-slice.md` 1, `reviewer.md` 1, `harness/README.md` 2, Sensor-Vertrag 1,
Architect-Verdikt 1, Review-Report 1 — der Report entstand nach der ersten Messung mit 6; der Verifikations-Report nennt den Namen
ebenfalls und macht die Zahl in der Closure zu 8; die drei fortgeschriebenen Register-Zustände
nennen den Namen je einmal und machen sie zu 11 — Planner nachgemessen); die Skill-Datei selbst nennt ihren Namen nicht, und die ADR bricht
den Dateinamen über einen Zeilenumbruch (kein Treffer). Nicht gefunden: weitere Träger der
Aussage „Handbuch-Versionshistorie“ (`git grep -il Versionshistorie`, ohne `docs/reviews`,
`done/`, Baseline) außer den bearbeiteten Dateien; gemeldet, nicht geändert (Accepted-ADRs
und Register sind fremde Träger): `ADR-0058`, `ADR-0090` (unberührbar) und die Register-Einträge
`BEO-PGC/handbuch-versionshistorie-uebersprungen` und
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` — Hinweis an den Planner
(Frist: Closure): die Prüfpunkte des ersten Eintrags tragen jetzt die Ergänzung „ohne
Kennungen“.

## 4. Trigger

**Start** (`next` → `in-progress`): die Architect-Entscheidung (Liefer-Punkt 0) ist als
[`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
`Accepted` gelandet (erfüllt; [Verdikt](../../../reviews/architect-verdict-handbuch-public-doc-check-gate-und-skill.md)).
Kein anderer Slice liegt in
`in-progress/` (WIP-Limit 1). Das Handbuch ist am Start weiter kennungsfrei
(Suchlauf-Zeile 1 gemessen am dann aktuellen Stand, sonst zuerst Befund an den Planner).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): (C) trennt sich als eigener Slice
  ab (`…-skill-und-pruefpunkte`); Gate, ADR und Standard bleiben hier.
- `in-progress` → `open` (blockiert): der Architect verwirft die Reichweite oder die Aufnahme
  in `make gates`, oder das Gate färbt am echten Baum rot durch einen Treffer, den der
  Slice nicht bereinigen darf (Handbuch-Inhalt ist ausgeschlossen) — der Fund geht an den
  Planner, kein Verstecken in einer Ausnahmeliste.

## 5. Closure-Trigger

DoD vollständig (inkl. ADR `Accepted`), Review ohne offenes HIGH/MEDIUM, `make gates` grün
(Exit-Code ungefiltert) mit dem neuen Gate, Mutationsprobe rot gesehen und grüner Lauf am
echten Baum, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Falsch-positive bei Befehlsbeispielen und Wortrand.** Entschieden
  ([`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
  Festlegung 3): der Wortrand-Schnitt (`byte-slice-x`, `CO-2 Emissionen` ohne Treffer) steht
  im Muster `P`; keine Ausnahme für Befehlsbeispiele. Am Stand `30a060c9` 0 Treffer im
  Handbuch (gemessen, Suchlauf Zeile 1). **Restgrenzen** (im Sensor-Vertrag zu führen):
  Groß-/Kleinschreibung zählt (`Slice-1` trifft nicht), Chronik-Sprache ohne Kennung bleibt
  grün. — **Ausgang:** *entfallen als Fehlalarm-Risiko, Grenzen benannt:* Muster `P` und
  Tabellentest mit Wortrand-Fällen tragen es (42 Fälle, Verifier gemessen; Mutation
  `{3}` zu `{9}` rot). Restgrenzen im Sensor-Vertrag: Groß-/Kleinschreibung,
  `ADR 0143` mit Leerzeichen, `--slice-x`; eine GitHub-URL mit `docs/plan/` trifft das
  Linkmuster und ist als Falsch-Positiv gewollt.
- **Gate-Fläche zu klein oder zu groß.** Zu klein: eine neue Nutzerdatei bleibt unbewacht;
  zu groß: eine Erzeugnis-Datei färbt falsch rot (`sdk-public-doc-check` Grenze 2). —
  **Ausgang:** *entschieden: Hybrid*
  ([`ADR-0143`](../../adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)
  Festlegung 1; Listen plus
  Vollständigkeitsprüfung, Exit 2 für unklassifiziert/fehlend); **entfallen:** der
  Tabellentest deckt die unklassifizierte `.md`, die fehlende genannte Datei und die
  ausgenommenen Dateien (Verifier: Mutation des Klassifikationszweigs rot).
- **Skill-Drift.** Der Skill kann von Schritt 17 und dem Reviewer-Punkt auseinanderlaufen
  (zwei Orte für „was ins Handbuch gehört“). Gegenmaßnahme: Schritt 17 und Reviewer
  verweisen auf den Skill, statt den Inhalt zu duplizieren; Verweis-Existenz im Suchlauf
  (Zeile 14). — **Ausgang:** *weiter offen* → `BEO-PGC/handbuch-versionshistorie-uebersprungen`
  im Register (Träger der Pflicht sind jetzt Skill, Schritt 17 und Reviewer-Punkt; die
  Verweise stehen, ob sie auseinanderlaufen, zeigt erst ein späterer Diff).
- **Zeilen-/Dateigröße-Gate für `implement-slice.md` kommt** (Auftraggeber-Vorgabe,
  Memory-Eintrag; übernommen, nicht gemessen): eine Textkürzung reicht dann nicht, es
  braucht eine strukturelle Lösung. Deshalb Liefer-Punkt 3: Inhalt von Schritt 17 in den
  Skill verlagern, Zeilenzahl nicht mehren. — **Ausgang:** *entfallen:* `implement-slice.md`
  hat 369 Zeilen gegenüber 373 am Parent (Verifier nachgemessen mit `wc -l`).
- **Regeltext-Änderung ohne eigene Freigabe.** Änderungen an `implement-slice.md` und
  `reviewer.md` sind Regeltexte; sie stützen sich auf die Freigabe des Auftraggebers („ja,
  nachdem …“ für Gate, ADR, Skill und Umstellung der Prüfpunkte), wörtlich im Kopf dieser
  Datei. Weitergehende Änderungen (z. B. Abschwächung der Pflicht) sind nicht gedeckt. —
  **Ausgang:** *entfallen:* die Änderung ist durch die Freigabe gedeckt („ja, nachdem die
  anderen offenen Slices abgearbeitet sind“); Review und Verifier lasen die Pflicht im
  Diff (beide Pflichten bleiben, keine Abschwächung).
- **Das Gate fängt Kennungen, nicht Chronik.** Eine Historie-Zeile in Chronik-Sprache ohne
  Kennung („früher … jetzt …“) bleibt grün; das ist die Lese-Hälfte (Skill, Reviewer).
  Grenze im Sensor-Vertrag zu führen. — **Ausgang:** *weiter bestehende Grenze, keine
  Aktion:* das Gate liest Kennungen, nicht Sinn; die Grenze steht im Sensor-Vertrag und im
  Skill, die Lese-Hälfte bleibt bei Reviewer und Verifier.
- **Kein Workflow ändert sich.** `ci.yml` fährt `make gates`; die Aufnahme in `GATE_CHECKS`
  ändert keinen Workflow ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht — am Start
  mit `git diff --name-only <Basis> -- .github` zu messen). — **Ausgang:** *weiter offen*
  → CI-Beobachtung: kein Workflow ändert sich (Verifier: `git diff --name-only … -- .github`
  leer), aber der erste CI-Lauf nach dem Push trägt das Gate erstmals auf dem Runner
  (`ci.yml`, Schritt „Gates“; das Gate braucht `bash`, `git`, `grep`, `find`, `sort`,
  `mktemp`, `sed`). Anker: Verifikations-Report §6. — **Ausgang (nachgetragen):**
  *entfallen:* der Lauf `ci` 37056062289 am Push-Commit `abf71ace` endete `success`
  (gemessen vom Hauptlauf mit `gh run watch`, 2026-10-02); `e2e` 37056062371 und `examples`
  37056062284 waren ebenfalls `success`.
- **Kein Release, keine Versionsänderung.** — **Ausgang:** *entfallen:* kein Tag, keine
  Versionsdatei berührt (Verifier).

## 7. Closure-Notiz

Ursprung der Angaben ([`AGENTS.md`](../../../../AGENTS.md) §3.12): **gemessen** = von der
genannten Rolle im eigenen Lauf; **übernommen** = aus einem Bericht ohne Nachmessung der
nennenden Rolle. Der Planner hat keine Zahl dieser Notiz nachgemessen, außer dem
Suchlauf-Feld (`make suchlauf-nachmessen`, 28 Zeilen) und den Register-Zählern (Dateien gezählt).

- **Was hat funktioniert:**
  - Gate, Tabellentest, Sensor-Vertrag, Skill, Umstellung der Prüfpunkte und Standard liegen
    vor; `make gates` Exit 0 mit dem Gate in der Lauf-Ausgabe (Verifier gemessen:
    `handbuch-public-doc-check: keine interne Kennung in 3 Nutzerdokumenten unter docs/user/`).
  - Tabellentest: 42 Fälle bestanden (Verifier gemessen; der Review sah vor der Fixrunde 31).
  - Sieben Einzelmutationen an Kopien im Scratchpad, alle rot (Verifier gemessen): echte
    Kennung und `docs/plan/`-Link, totes Link-Muster, `|| true` im Lesepfad, verschluckter
    `find`-Fehler, geänderter Meldungstext, `grep -a` zu `grep -n`, `{3}` zu `{9}`. Die
    Aussage gilt für diese Stellen in Skript und Tabellentest; Mutationen anderer Stellen
    (Wortrand `slice-`) sind nicht gefahren, dort tragen die Fälle `byte-slice`,
    `CO-2 Emissionen`, `Slice-1`.
  - `implement-slice.md`: 373 Zeilen am Parent, 369 am Stand der Closure (Verifier
    nachgemessen); die Pflicht „Version hochzählen und Zeile in der Änderungshistorie“ bleibt.
  - Muster-Messung am Baum (Planner/Architect gemessen, Suchlauf Zeilen 1–8): Handbuch,
    Standard und `version.md` je 0, `releasing.md` 26, die vier Erzeugnisse 4 · 3 · 74 · 21;
    Links nach `docs/plan/`/`docs/reviews/` 0.
- **Was ging anders als geplant:**
  - **Fixrunde ohne Re-Review (ehrliche Abweichung).** Der Review-Report (0 HIGH, 0 MEDIUM,
    4 LOW, sieben Findings F-1 bis F-7 samt drei INFO) sah den Stand vor `28a5aafe`. Die
    Fixrunde `28a5aafe` behob die vier LOW (Gate fail-closed bei Lesefehlern, Tabellentest
    bindet Meldungstext und `LH-RB-`-Zweig, Plan-DoD berichtigt). Ein Re-Review danach fand
    nicht statt. Der Verifier las das Skript (94 Zeilen) vollständig und mutierte alle
    Lesefehler-Zusagen einzeln (grep ≥ 2, `find`, NUL-Byte; sieben Mutationen rot) und
    hielt einen weiteren Reviewer-Durchgang nicht für nötig. Die Review-Reports bleiben
    unverändert; der Haken „Review durchgeführt“ trägt die Abweichung als Zusatz.
  - **Lesefehler-Semantik nur im Sensor-Vertrag.** `ADR-0143` nennt Exit 2 für
    Klassifikationsfehler, nicht für Lesefehler; das Gate verhält sich strenger als der
    Wortlaut (Verifier: Abweichung mit Deckung im Zweck, kein Verstoß). Die Festlegung
    steht im Sensor-Vertrag; die Review-Frage 1 (soll `ADR-0143` sie festlegen, auch für
    `sdk-public-doc-check`) bleibt als Architect-Frage im Folge-Slice (siehe dort §4).
  - **Das SDK-Gate ist fail-open (Verifier gemessen):** Kopie des `sdks`-Baums, `README.md`
    mit Kennung und `chmod 000`: `grep` meldet „Keine Berechtigung“, das Skript
    `keine interne Kennung` und Exit 0 (Ursache `xargs … grep || true`, `find` mit
    `2>/dev/null`). Das Vorbild des neuen Gates trug die Lücke; beim Kopieren des Musters
    wurde sie mitkopiert und im Review (F-1) gefunden.
  - **Zitat-Korrektur nach Lifecycle-Wechsel.** `ADR-0143` nennt den Pfad des Plans unter
    `planning/open/` (Kopf, Geschichte); nach dem Wechsel nach `done/` ist das eine
    Zitat-Korrektur nach `ADR-0073` (eigener Commit, genau eine neue Geschichte-Zeile,
    Wortlaut der Entscheidung unberührt).
  - Die Suchlauf-Zeile zum Skill-Namen steht in der Closure bei 11 statt 7: der
    Verifikations-Report nennt den Namen (8) und die drei fortgeschriebenen Register-Zustände
    je einmal (11); Planner nachgemessen.
- **Steering-Loop-Eintrag (Lerneintrag):** (1) Ein Rückfall in Nutzerdokumentation war bisher
  nur über Lesen (Reviewer) gefangen; das Gate ist der erste **Sensor** dieser Klasse, mit
  der benannten Grenze: es **liest Kennungen, nicht Sinn** (Chronik-Sprache ohne Kennung
  bleibt grün; Skill und Reviewer tragen die Lese-Hälfte). (2) Ein Gate ist erst dann ein
  Gate, wenn ein **Lesefehler** es rot färbt: ein Muster, das man als Vorbild kopiert, gibt
  seine Lücken (`|| true` hinter `grep`) weiter; die geschärfte Regel: ein neuer
  `grep`-Wächter trennt Exit 1 (kein Treffer) von Exit ≥ 2 (Lesefehler → Exit 2) und bindet
  beides im Tabellentest — gefunden durch den Review, nicht durch den Implementer. (3)
  **Benannte Spec-Lücke:** die Lesefehler-Semantik der Wächter-Gates steht in keiner
  `Accepted` ADR (`ADR-0134`, `ADR-0143`); Träger ist der Sensor-Vertrag, die
  Architect-Frage steht im Folge-Slice. — liegt in `harness/sensors/handbuch-public-doc-check.md` (Ausgangstabelle, Exit 2) und `tools/harness/handbuch-public-doc-check.sh`.
  Auslöser: `BEO-PGC/handbuch-versionshistorie-uebersprungen` (3×, verkörpert),
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3×, verkörpert).
- **Beobachtungs-Register (`../observations/`):** Keine neue Beobachtung angelegt, drei Einträge
  fortgeschrieben (ohne neue `evidence/`-Datei, keine Zähl-Inflation):
  `BEO-PGC/handbuch-versionshistorie-uebersprungen` (Träger der Prüfpunkte: Skill,
  Schritt 17, Reviewer; „ohne Kennungen“), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (der Inhalt „neue Betreiber-Oberfläche“ liegt weiter in Schritt 17 und im Skill) und
  `BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (das Gate als Fangnetz für
  `docs/user/`, Grenze „liest Kennungen, nicht Sinn“; Zähler bleibt **1×**).
- **Folge-Slices:**
  [`sdk-public-doc-check-lesefehler-fail-closed`](slice-sdk-public-doc-check-lesefehler-fail-closed.md)
  — ist eine Datei in `open/` (Architect-Frage, ob eine Folge-ADR nötig ist, steht in
  dessen §4).
- **Risiken aus §6:** Falsch-positive entfallen (Grenzen benannt) · Gate-Fläche entfallen ·
  Skill-Drift weiter offen → Register · Zeilen-Gate entfallen (369) · Regeltext-Änderung
  entfallen (durch Freigabe gedeckt) · Gate fängt Kennungen, nicht Chronik: bestehende
  Grenze · Workflow weiter offen (erster CI-Lauf mit dem Gate auf dem Runner, Anker
  Verifikations-Report §6; Ergebnis trägt der Hauptlauf nach; `AGENTS.md` §3.10 berührt den
  Slice nicht, weil kein Workflow geändert wird) · kein Release entfallen.
- **Offene Entscheidungen des Auftraggebers:** keine.
- **Drei Paarungen:** Anker: Review F-1 bis F-7 und Verifikation §7 Punkte 1 bis 6. Folge-Slice:
  `sdk-public-doc-check-lesefehler-fail-closed`. Register: `handbuch-versionshistorie-uebersprungen`,
  `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`, `intern-kennungen-in-ausgelieferten-texten`
  (fortgeschrieben).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den Pfaden
`tools/harness/`, `harness/`, `.harness/skills/`, `.claude/commands/`, `docs/user/` und
`docs/plan/adr/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Treffer in §6/§2:
`handbuch-versionshistorie-uebersprungen` (Klasse, deren Prüfpunkte dieser Slice umstellt),
`zwei-quellen-drift-handbuch-gegen-pflichtenheft`,
`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (Nachbarn). Der Implementer
zählt die `evidence/`-Dateien bei der Closure nach.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
