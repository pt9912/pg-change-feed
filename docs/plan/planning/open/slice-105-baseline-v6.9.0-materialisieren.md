# Slice slice-105: Baseline v6.9.0 materialisieren + Namenskonvention-Adaption

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die über die eigene
DoD dieses Slice hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht). Dieses Repo läuft durchgehend wellenlos
(siehe `slice-104` u. a.).

**Bezug:** `harness/conventions.md` §Baseline / MR-000 (Adoptions-Erklärung —
dieser Slice hebt den dort referenzierten Stand von `v6.5.0` auf `v6.9.0`;
keine ADR ist einschlägig, die Baseline-Adoption selbst ist kein
ADR-Gegenstand). Dies ist **Folge-Slice 1 von 3** einer Migrations-Recherche
(`/tmp/regelwerk-migration-v6.5.0-zu-v6.9.0.md`, Kurs-Wellen 129–137,
`pt9912/ai-harness-course`) — die beiden anderen (Gate-Index-Konsolidierung
aus Welle 129, DoD-Rot-vor-Grün-Regel aus Welle 135) sind eigene,
nachfolgende Slices mit noch zu vergebenden Namen (siehe §1 Abgrenzung).

**Berührte Spec-Stellen:** — (geprüft: `grep -n "baseline\|MR-\|slice-" spec/lastenheft.md spec/pflichtenheft.md spec/architecture.md` liefert keine Treffer zum Harness-Regelwerk selbst; reiner Harness-/Prozess-Vorgang ohne Vertrags- oder Draht-Bezug).

**Verantwortlich:** —.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-17.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die vendored Kurs-Baseline dieses Repos wird von `v6.5.0` auf
`v6.9.0` gehoben (`.harness/baseline/v6.9.0/{regelwerk,templates}/`, netzlos <!-- d-check:ignore (Ziel-Stand dieses Migrationsplans, noch nicht adoptiert — kein Baseline-Pin-Fund) -->
materialisiert, `SHA256SUMS`-geprüft aus dem Release-Asset
`https://github.com/pt9912/ai-harness-course/releases/download/v6.9.0/lab-regelwerk.zip`),
`harness/conventions.md` §Baseline zeigt danach auf den neuen Stand, und die
Baseline-Namenskonvention für Slice-/Welle-Kennungen bei Mehr-Schreiber-Betrieb
(Kurs-Wellen 130/131, `grundlagen-source-precedence.md` §Vergabe) wird als
eigener `MR-*`-Adaptions-Eintrag repo-lokal verankert: Dieses Repo ist —
unabhängig davon, dass ein einzelner Mensch es lenkt — strukturell ein
**Mehr-Schreiber-Repo**, weil mehrere, teils parallel gestartete
Claude-Agenten (Architect, Implementer, Reviewer, Verifier, Planner) jeweils
eigenständig committen; die Baseline definiert **„Schreiber ist, was
committet: ein Mensch, ein Agent, ein Automat"**. Bestehende `slice-001`–
`slice-104` bleiben nummeriert (Bestandsschutz, siehe Abgrenzung); ab dem
nächsten neu angelegten Slice nach diesem hier gilt die Namensform.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Umbenennung von `slice-001`–`slice-104` (und aller darauf verweisenden
  Artefakte) auf die neue Namensform.** *Bestand bleibt bewusst stehen.*
  Diese 104 Kennungen sind massenhaft querverwiesen — ADRs (`Bezug:`-Felder),
  das Beobachtungs-Register (`evidence/slice-<NNN>.md`), Closure-Notizen,
  Review-/Verify-Reports, Commit-Messages. Eine Umbenennung wäre gegenüber
  dem erreichbaren Nutzen (Formkonsistenz) unverhältnismäßig disruptiv und
  bräche `git log --follow`-Ketten sowie Hunderte Anker. Nutzerentscheidung
  vom 2026-09-17, im neuen `MR-*`-Eintrag dokumentiert.
- **Gate-Index-Konsolidierung (Kurs-Welle 129: `AGENTS.md` §4 auf Regel +
  Zeiger kürzen, `harness/README.md` §Sensors bleibt einzige Volltabelle).**
  *Ein Folge-Slice übernimmt es* — eigener, noch zu vergebender Slice-Name
  (Migrations-Schritt 2 der Recherche). Eigene Schicht (`AGENTS.md`-Struktur
  statt Baseline-Materialisierung), eigene Beweislast (Vollständigkeits-Grep
  über `.github/`), gehört nicht in denselben Zug wie das reine
  Baseline-Update.
- **DoD-Rot-vor-Grün-Regel für sicherheits-/korrektheitskritische
  Testbehauptungen + E2E-Gate-Typ-Abgrenzung (Kurs-Welle 135).** *Ein
  Folge-Slice übernimmt es* — eigener, noch zu vergebender Slice-Name
  (Migrations-Schritt 3 der Recherche). Inhaltliche Regel-Neuformulierung
  in `AGENTS.md`/`.harness/skills/reviewer.md`, unabhängig vom
  Baseline-Materialisierungs-Mechanismus dieses Slice.
- **Inhaltliche Prüfung/Anwendung der übrigen Wellen-Funde (132, 134, 136)
  der Recherche.** *Es wäre ein anderer Vorgang* bzw. *Bestand bleibt bewusst
  stehen* — laut Recherche für keine davon akuter Handlungsbedarf in diesem
  Repo (132: bereits informell gelebte Praxis, kein Datei-Delta erzwungen;
  134: reine Vorlagen-Klarstellung, betrifft `harness/README.md` als fertige
  Datei nicht; 136: reine Verständnis-Klärung des Freshness-Audits selbst,
  keine betroffene Datei in diesem Repo). Sie werden mit der neuen Baseline
  automatisch mitgeliefert und in §7 zur Kenntnis vermerkt — kein eigener
  Umsetzungsschritt.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **LP1 — Baseline `v6.9.0` netzlos materialisiert, Integrität geprüft,
      genau ein `<tag>`-Verzeichnis bleibt bestehen.** Release-Asset
      `lab-regelwerk.zip` von
      `https://github.com/pt9912/ai-harness-course/releases/download/v6.9.0/lab-regelwerk.zip`
      laden, `SHA256SUMS` **vor** dem Entpacken gegen den Asset-Hash prüfen
      (Vorgehen: `modul-02-harness-bootstrap.md` §Bootstrap — in der **neuen**
      v6.9.0-Fassung nachlesen, sobald sie vorliegt, falls sich der
      Bootstrap-Ablauf zwischen den Ständen geändert hat), nach
      `.harness/baseline/v6.9.0/{regelwerk,templates}/` entpacken. <!-- d-check:ignore (Ziel-Stand dieses Migrationsplans, noch nicht adoptiert — kein Baseline-Pin-Fund) -->
      **Geklärt (Recherche dieses Plans, siehe `tools/harness/baseline-verify.sh`):**
      Der Sensor verlangt zwingend **genau ein** Verzeichnis direkt unter
      `.harness/baseline/` (`dirs=("$base"/*/)`, `exit 1` bei `${#dirs[@]}`
      ≠ 1) — er verträgt **keine** zwei parallel vendorten Stände, auch nicht
      unter einem Unterordner wie `.harness/baseline/done/`, denn ein solcher
      Unterordner zählt selbst als zweites `*/`-Element und lässt den Sensor
      ebenso rot laufen wie zwei gleichrangige `<tag>`-Verzeichnisse. Die in
      der Recherche vorgeschlagene Option „`git mv` nach
      `.harness/baseline/done/`" **löst das Problem nicht** — nur eine
      Entfernung aus `.harness/baseline/` selbst tut es. Richtige Form daher:
      `.harness/baseline/v6.5.0/` vollständig entfernen (`git rm -r`), sobald
      `v6.9.0` steht — kein Datenverlust, `git log`/`git show
      <commit>:.harness/baseline/v6.5.0/...` trägt den alten Stand vollständig
      und dauerhaft. `make baseline-verify` grün gegen den einzigen
      verbleibenden Stand `v6.9.0`.
- [ ] **LP2 — `harness/conventions.md` §Baseline zeigt auf `v6.9.0`, keine
      lebende Erwähnung von `v6.5.0` bleibt zurück.** `Stand:` auf `v6.9.0`,
      `Datum der Adoption:` auf das Datum dieses Slice-Laufs; die
      Wellen-Registerzeile im Zeiger-Kommentar (`Kurs-Welle 137 · <Datum>`)
      nachgetragen. Beleg: `grep -rn "baseline/v6.5.0\|Stand.*v6.5.0"` über
      das Repo liefert **keine** Treffer außerhalb `docs/plan/planning/done/`
      und `git`-History (historisch korrekte, unveränderliche Records —
      `AGENTS.md` §3.5 analog: Records ändern sich nicht rückwirkend).
- [ ] **LP3 — neuer `MR-*`-Eintrag für die Namenskonvention-Entscheidung,
      verlinkt in `harness/conventions.md`.** Nächste freie Nummer ermitteln
      (Adaptions-Block-Tabelle in `harness/conventions.md` führt aktuell nur
      `MR-001` — die nächste freie ist `MR-002`, vom Implementer-Lauf am
      dann aktuellen Stand erneut zu bestätigen, falls zwischenzeitlich ein
      weiterer `MR-*` gezogen wurde). Eigene Datei
      `harness/conventions/MR-<NNN>-<slug>.md`, kopiert aus der
      Eintrags-Vorlage (`.harness/baseline/v6.9.0/templates/harness/conventions/MR-NNN-titel.template.md`, <!-- d-check:ignore (Ziel-Stand dieses Migrationsplans, noch nicht adoptiert — kein Baseline-Pin-Fund) -->
      falls unverändert gegenüber `v6.5.0/templates/harness/conventions/MR-NNN-titel.template.md`
      — real vergleichen, nicht annehmen). Pflichtfelder: Datum,
      Geltungsbereich (`harness/conventions.md` §Aktive Adaptionen; alle
      künftig neu angelegten Slice-/Welle-Plan-Dateien), Ersetzt-Baseline-Regel
      (Link mit Anker auf `grundlagen-source-precedence.md#vergabe-woher-die-naechste-nummer-kommt`
      im **neuen** `v6.9.0`-Pfad — genaue Anker-Form real gegen die
      materialisierte Datei prüfen, nicht raten), Adaption (wörtliches Zitat
      der Baseline-Regel „Welle- und Slice-Kennungen sind Namen, nicht
      Nummern — unabhängig von der Schreiberzahl" und „Schreiber ist, was
      committet: ein Mensch, ein Agent, ein Automat" — Zitat-Wortlaut real
      gegen die materialisierte `v6.9.0`-Datei verifizieren, LP1 liefert die
      Quelle erst), Begründung (Mehr-Agenten-Schreiber-Modell dieses Repos
      trotz einem menschlichen Lenker — Architect/Implementer/Reviewer/
      Verifier/Planner committen eigenständig, teils parallel gestartet, z. B.
      die parallele Architect+Verifier-Sequenz bei `slice-101`), Entscheidung
      (bestehende `slice-001`–`slice-104` bleiben nummeriert — Bestandsschutz,
      siehe §1 Abgrenzung —, ab dem nächsten neu angelegten Slice nach
      `slice-105` gilt: Name trägt das Präfix eines vorhandenen Ankers
      (`LH-*`, `ADR-*`, `CO-*`), wenn einer existiert, sonst ein freier Slug),
      Auflösungs-Trigger (permanent — die Entscheidung ist eine dauerhafte
      Konvention, kein Provisorium; kein bekanntes Re-Evaluierungs-Ereignis
      identifiziert). Zeile in der Aktive-Adaptionen-Tabelle in
      `harness/conventions.md` ergänzt.
- [ ] `make gates` grün — ungepiped geprüft, Exit-Code direkt ausgewertet
      (`AGENTS.md` §3.9), als eigener, abgeschlossener Schritt **vor** jedem
      `git push`.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/` liegt vor
      (Modul 11, frischer Kontext).
- [ ] Doku-Update für den gehobenen Baseline-Stand siehe LP1–LP3 — kein
      weiterer öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls
      dieser Slice einen Inventur-Fund auflöst** — **Entfällt.** Repo ist
      Greenfield (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = GF),
      keine `docs/plan/planning/reconciliation.md` vorhanden.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im
      Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der
      nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/v6.9.0/regelwerk/**` | neu | Netzlos materialisiertes Regelwerk aus dem `v6.9.0`-Release-Asset (LP1). | <!-- d-check:ignore (Ziel-Stand dieses Migrationsplans, noch nicht adoptiert — kein Baseline-Pin-Fund) -->
| `.harness/baseline/v6.9.0/templates/**` | neu | Netzlos materialisierte Templates aus demselben Asset (LP1). | <!-- d-check:ignore (Ziel-Stand dieses Migrationsplans, noch nicht adoptiert — kein Baseline-Pin-Fund) -->
| `.harness/baseline/v6.9.0/SHA256SUMS` | neu | Prüfsummen-Datei des neuen Stands, Grundlage für `make baseline-verify` (LP1). | <!-- d-check:ignore (Ziel-Stand dieses Migrationsplans, noch nicht adoptiert — kein Baseline-Pin-Fund) -->
| `.harness/baseline/v6.5.0/**` | löschen | Entfernt, weil `tools/harness/baseline-verify.sh` genau ein `<tag>`-Verzeichnis unter `.harness/baseline/` verlangt — Historie bleibt vollständig in `git` erhalten (LP1). |
| `harness/conventions.md` | update | §Baseline: `Stand: v6.9.0`, `Datum der Adoption:` heutiges Datum, Wellen-Registerzeile; neue Zeile in §Aktive Adaptionen für den neuen `MR-*`-Eintrag (LP2, LP3). |
| `harness/conventions/MR-<NNN>-<slug>.md` | neu | Adaptions-Eintrag für die Namenskonvention-Entscheidung (LP3). |
| `harness/image-hash.txt` | prüfen, nicht erwartet | Kein Build-Kontext berührt (`.harness/`/`harness/` sind keine `Dockerfile`-Build-Kontext-Dateien) — `make image` läuft nicht erneut nötig; real prüfen, nicht annehmen (§3.12). |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt aktuell keinen
Slice (geprüft: `ls docs/plan/planning/in-progress/` zeigt nur `roadmap.md`,
keine Slice-Datei — WIP-Limit frei) und dieser Slice hat keine Abhängigkeit
zu einer laufenden Welle oder einem anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `open` (blockiert — Carveout?): wenn die
  `SHA256SUMS`-Prüfung des `v6.9.0`-Release-Assets fehlschlägt (Hash-
  Mismatch gegen den vom Release-Autor veröffentlichten Wert, oder das Asset
  ist nicht erreichbar/verändert) — dann ist die Baseline nicht vertrauenswürdig
  materialisierbar, und der Slice geht als Blocker zurück, bis das Asset
  erneut geprüft werden kann (kein Carveout: es gibt kein akzeptables
  „vorläufig ungeprüft").
- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn
  `make baseline-verify` sich als strukturell unfähig erweist, einen
  Versionswechsel überhaupt abzubilden (z. B. wenn das Skript den Pfad
  `v6.5.0` fest verdrahtet statt über den Glob zu ermitteln — laut Recherche
  dieses Plans **nicht** der Fall, das Skript ist bereits generisch über
  `$base/*/`; sollte sich das bei der Umsetzung als falsch erweisen, ist ein
  Sensor-Umbau ein eigener, größerer Vorgang und gehört nicht mehr in diesen
  Slice) — oder wenn sich während der Umsetzung herausstellt, dass die
  Namenskonvention-Adaption (LP3) tatsächlich eine inhaltliche Prüfung
  weiterer Kurs-Wellen (129/135) voraussetzt, die dieser Slice bewusst
  ausklammert.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** (1) `make baseline-verify` meldet grün gegen
genau ein `<tag>`-Verzeichnis (`v6.9.0`), `wc -l < SHA256SUMS` zeigt die neue
Dateizahl des `v6.9.0`-Bundles. (2) `grep -rn "baseline/v6.5.0"` über das
gesamte Repo liefert außerhalb von `done/`-Records und `git`-History keine
Treffer mehr — **und** `make gates` insgesamt grün.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in
§7. Naheliegender Kandidat: eine geschärfte Fassung des Bootstrap-Vorgehens
(`modul-02-harness-bootstrap.md` §Freshness-Audit) — dieser Slice selbst ist
das Ergebnis eines Freshness-Audits, das *nachträglich* im Nutzer-Dialog eine
Welle (130/131) fand, die die ursprüngliche automatisierte Recherche als
„nicht anwendbar" fehleingeordnet hatte (Kategorienfehler: Ein-Schreiber-Repo
wurde mit Ein-Mensch-Lenker verwechselt). Das ist ein Kandidat für einen
Beobachtungs-Register-Eintrag oder eine geschärfte Freshness-Audit-Anleitung
— welcher Kandidat trägt, entscheidet der Umsetzungslauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **SHA256SUMS-Prüfung des Release-Assets schlägt fehl.** Der veröffentlichte
  Hash des `lab-regelwerk.zip`-Assets für `v6.9.0` stimmt nicht mit dem
  heruntergeladenen Artefakt überein, oder das Release ist zwischenzeitlich
  zurückgezogen/verändert worden. Muss vor jedem Entpacken geprüft werden
  (LP1). — *Ausgang bei Closure einzutragen.*
- **`make baseline-verify` verträgt strukturell keinen Versionswechsel.**
  Die Recherche dieses Plans hat das Skript bereits gelesen
  (`tools/harness/baseline-verify.sh`) und als generisch über `$base/*/`
  befunden — kein hartkodierter `v6.5.0`-Pfad. Das Risiko besteht dennoch
  fort, bis der reale Lauf gegen `v6.9.0` bestätigt, dass kein anderer,
  bisher unbemerkter Pfad-Bezug (z. B. in `.d-check.yml` oder einem anderen
  Sensor) den alten Tag literal referenziert. — *Ausgang bei Closure
  einzutragen.*
- **Freshness-Audit findet zusätzliche relevante Wellen, die die
  ursprüngliche Recherche übersah — wie beim Namenskonvention-Fund selbst.**
  Die Migrations-Recherche (`/tmp/regelwerk-migration-v6.5.0-zu-v6.9.0.md`)
  hatte Welle 130/131 fälschlich als „nicht anwendbar (Ein-Schreiber-Repo)"
  eingeordnet — ein Kategorienfehler, der erst im Nutzer-Dialog aufgedeckt
  wurde. Es ist nicht auszuschließen, dass ein weiterer Fund derselben Art
  (falsch als „nicht anwendbar" markiert) in den Wellen 129/132/133/134/136/137
  steckt und erst beim realen Materialisieren/Lesen der `v6.9.0`-Fassung
  auffällt. — *Ausgang bei Closure einzutragen*; ein hier neu gefundener
  Punkt braucht ggf. einen weiteren Folge-Slice mit eigenem Namen (die neue
  Konvention gilt für ihn bereits, siehe §1).

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

*Wird beim Umsetzungslauf gefüllt — noch offen (Plan-Anlage, kein Code, keine
Implementierung).*

**Zur Kenntnis (aus der Migrations-Recherche, keine eigene Umsetzung dieses
Slice):** Kurs-Welle 132 (Tests-Zeile an Akzeptanzkriterien der
Requirement-ID gebunden — bereits informell gelebte Praxis) · Welle 133
(kursinternes Mirror-Tooling, nicht anwendbar) · Welle 134
(`README.template.md` referenziert `gate.template.md` als Kopiervorlage —
`harness/README.md` dieses Repos ist eine fertige Datei, keine Vorlage) ·
Welle 136 (Klarstellung: Append-only-Klausel gilt pro Artefakt-Klasse, nicht
pro Vorlagen-Datei — reine Verständnis-Klärung ohne betroffene Datei hier) ·
Welle 137 (neuer sechster Lifecycle-Übergang `open|next → done` für Slices,
deren Gegenstand ein anderer Slice übernimmt — `open/`/`next/` sind aktuell
leer, kein akuter Anwendungsfall, aber ab sofort bekannt für künftige
Planer-Läufe, sobald mehrere offene Slices zusammengelegt werden).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind ausschließlich
`.harness/` (vendorte Baseline) und `harness/` (Konventionen, Adaptions-
Einträge) — die repo-weite Default-Sub-Area `*`/`PGC` aus der
Modus-Deklaration in [`harness/conventions.md`](../../../../harness/conventions.md).
Keine feinere Sub-Area ist im Repo deklariert, die diese Pfade eigens führt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`grep -rli "baseline\|mr-\|namenskonvention\|mehrere.*schreib\|schreiberzahl"
docs/plan/planning/observations/` sowie gezielt
`grep -rli "kennung|namenskonvention|schreiber" docs/plan/planning/observations/`)
— **keine Treffer**, die die Baseline-Versionierung oder die Slice-/
Welle-Namenskonvention betreffen; die breiteren Treffer der ersten Suche sind
Fehltreffer auf das allgegenwärtige Wort „Kennung" in unverwandtem Kontext.
Keine Treffer sind ebenfalls eine Antwort.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**. Ein
Block ist unten dennoch angelegt (statt nur des Kurzhinweises), weil dieser
Slice selbst eine Adaptions-Entscheidung mit Dauerwirkung trifft und die
Konventionen-Dichte/Phase-Reife-Einschätzung für spätere Regelwerk-Hebungen
als Referenz taugt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

### Sub-Area: `*`/`PGC` (Default)

- **Modus:** GF
- **Konventionen-Dichte:** hoch — `harness/conventions.md` MR-000 verankert
  die Baseline-Adoption bereits als Regel mit Stand-Feld; der
  Adaptions-Block-Mechanismus (`MR-<NNN>`, eigene Datei + Index-Zeile) ist
  mit `MR-001` bereits einmal gelebt.
- **Phase-Reife:** Phase 5 — die Sub-Area trägt bereits einen laufenden
  Baseline-Bootstrap-Zyklus, einen Adaptions-Block und mehrere Sensoren
  (`baseline-verify`, `docs-check`); keine Erstanlage.
- **Evidenz-/Diskrepanz-Risiko:** niedrig — GF, Doc führt; keine
  Register-Treffer zu diesem Mechanismus (siehe oben). Das einzig konkrete
  Risiko ist die Erst-Anwendbarkeit der neuen Namenskonvention selbst, nicht
  ein Bestands-Drift.
- **Reconciliation-Aufwand:** entfällt (GF, kein Brownfield-Bestand); kein
  Trigger einer der vier Klassen einschlägig, da keine Sub-Area BF/Hybrid ist.
