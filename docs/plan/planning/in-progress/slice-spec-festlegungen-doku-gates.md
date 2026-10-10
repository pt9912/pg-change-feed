# Slice spec-festlegungen-doku-gates: Festlegung von docs-check im Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(Entscheidung 7: Baseline-Aktualisierung als bewusster Bootstrap-Vorgang). ADR-Kette
von `make docs-check` (gemessen bei der Rückführung, siehe §4):
[`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
[`ADR-0074`](../../adr/0074-zitationsform-schwester-repo-hausform.md),
[`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md),
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
(Entscheidung 3, Home-relative Pfade in `hostpaths`) für `hostpaths`;
[`ADR-0094`](../../adr/0094-review-matrixklasse-kennung-statt-adresse.md),
[`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md),
[`ADR-0097`](../../adr/0097-observation-matrixklasse-review-verboten.md),
[`ADR-0099`](../../adr/0099-slice-welle-review-regel-zurueckgenommen.md) für
`matrix`;
[`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md),
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
(Block `versions:`) für `versions`;
[`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
(`Schärft:`-Kante zur Festlegung, Provenance-Marker, `matrix.exempt-paths`).
Keine `LH-*`-Anforderung ist berührt: die
Festlegung gilt einem Harness-Werkzeug, nicht dem Produkt.

**Berührte Spec-Stellen:** `spec/pflichtenheft.md` §7 „Festlegungen der
Harness-Werkzeuge“ (neue Zeile; die Kennung vergibt der Implementer).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Neuschnitt von
`slice-spec-festlegungen-harness-werkzeuge`). **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**; die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Herkunft.** Neuschnitt von `slice-spec-festlegungen-harness-werkzeuge` am
2026-10-07 (dort §1 „Planänderung“): der Pilot legt §7 „Festlegungen der
Harness-Werkzeuge“ im Pflichtenheft an und trägt die Festlegungen von
`make zitat-vergleich` und des Leer-Tests der Teil-Range; die Festlegungen der
übrigen Gates gehen je Gruppe an einen eigenen Slice. **Übergabe an diesen
Slice:** die Delta-Punkte R2, R7, R9, T6 und T8 des Bumps auf Baseline v6.16.0
(`slice-harness-baseline-v6-16-0`, Abschnitt „Bump-Ablauf — Belege“) für das
Gate `make docs-check` (alle Module der Liste `modules:` in `.d-check.yml`,
samt `hostpaths`).

**Planänderung bei der Rückführung `in-progress → next` (2026-10-10, §4).**
Bis dahin trug dieser Slice zusätzlich `make commit-traceability`,
`make baseline-verify` und `make doc-immutable` (samt der Übergabe F-5/F-13 aus
`slice-spec-festlegungen-harness-werkzeuge`). Die drei gehen mit ihren
Übergabe-Punkten an `slice-spec-festlegungen-commit-baseline-gates`; der Text
der Übergabe steht dort in §2.

**Ziel:** `spec/pflichtenheft.md` §7 trägt für `make docs-check` eine
Festlegung mit eigener `SPEC-<NNN>` (je Modul was als Treffer gilt, Randform,
Ausgänge); der Vertrag `harness/sensors/docs-check.md` nennt die Kennung, statt
Schwelle und Randform selbst zu tragen, und die Zeile in `harness/README.md`
§Sensors trägt die Bindung „Spec-Kennung“.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Struktur von §7, die `Schärft:`-Konvention im ADR-Index und die
  `MR-001`-Ablösung** — Gegenstand von `slice-spec-festlegungen-harness-werkzeuge`
  (Pilot); dieser Slice fügt Zeilen in eine bestehende Tabelle ein.
- **`make commit-traceability`, `make baseline-verify` und `make doc-immutable`
  samt der Übergabe F-5/F-13** — Folge-Slice
  `slice-spec-festlegungen-commit-baseline-gates` (§2 dort, Übergabe-Block);
  abgetrennt bei der Rückführung, weil die Quellen der vier Gates zusammen
  eine Review-Sitzung sprengen (§4).
- **`structure` in einen eigenen Slice abtrennen** — Bestand bleibt bewusst
  zusammen: `docs-check` ist ein Gate mit einem Vertrag; eine Festlegung ohne
  `structure` ließe `harness/sensors/docs-check.md` halb mit Verweis, halb mit
  eigener Randform zurück, und keiner der zwei Teile wäre für sich lieferbar
  (Schnitt nach Lieferwert, nicht nach Modul).
- **Gates anderer Gruppen** — Folge-Slices:
  `slice-spec-festlegungen-kennungs-gates` (sdk-, handbuch-public-doc-check,
  ausgabe-kennungen-check, meldungscodes-check),
  `slice-spec-festlegungen-code-gates` (a-check, generated-sync),
  `slice-spec-festlegungen-coverage-gates` (coverage-gate, db-adapter-coverage),
  `slice-spec-festlegungen-pruefer-hooks` (Prüfer und Hooks ohne Gate).
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die `Schärft:`-Kante entsteht
  über eine Architect-ADR (§2, §6).
- **Das Verhalten eines Gates ändern** — ein anderer Vorgang: die Festlegung
  beschreibt, was das Werkzeug heute entscheidet; weicht der Vertrag vom
  Werkzeug ab, ist das ein Befund mit eigenem Slice (§6).
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `spec/pflichtenheft.md`,
  `harness/sensors/docs-check.md`, `harness/README.md` und eine ADR.

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
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [x] **Festlegung in §7 (Liefer-Punkt 1).** `spec/pflichtenheft.md` §7 trägt
      eine Zeile mit eigener `SPEC-<NNN>` für `make docs-check`: je Modul der
      Liste `modules:` in `.d-check.yml` (links, anchors, ids, matrix,
      versions, structure, hostpaths, tracked) was als Treffer gilt, samt der
      Reichweite von `hostpaths`, in der Form, die der Pilot für §7 festlegt
      (R9). *Zu belegen durch:* Gegenprobe der Quellen und Anschluss-Frage
      (§3), Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) und
      `make docs-check` Exit 0.
      *Beleg:* Kennung `SPEC-040` (nächste freie: `git grep -ohE 'SPEC-[0-9]{3}'`
      am Stand `942ebf3f` endet bei `SPEC-039`), Commit `9e43f968`; Tabelle,
      Absätze „Zu `SPEC-040` — …“ und Historie-Zeile wie beim Pilot.
      Gegenprobe (76 Zeilen), Anschluss-Frage und Messungen M0 bis M23 in §3;
      Suchlauf in §3, `make suchlauf-nachmessen` Exit 0
      („16 Zeilen stimmen“); `make docs-check` am Arbeitsbaum nach `9e43f968`
      mit diesem Plan: Exit 0, `d-check: 1830 Datei(en) geprüft, 0 Befund(e)`.
      Die drei Befunde aus §6 hat `ADR-0163` entschieden; die zwei Sätze zu
      Provenance-Marker und `matrix.exempt-paths` stehen nach ihrem Wortlaut in
      `SPEC-040` (Commit nach `deecc466`).
- [x] **Vertrag und Index verweisen (Liefer-Punkt 2).**
      `harness/sensors/docs-check.md` nennt die Kennung der Festlegung und
      trägt Schwelle und Randform nicht mehr selbst (R2, R7, T8); die Zeile
      `make docs-check` in `harness/README.md` §Sensors trägt die Bindung
      „Spec-Kennung“ (T6).
      *Beleg:* Commit `9e43f968`. Der Vertrag verweist in §Vertrag („Wer was
      trägt“), §Grenze 8, 9, 11, §Ausgabe und §Bindung auf `SPEC-040`; die
      Befund-Code-Tabelle, die Zellgrenzen 220/120 und die Exit-Tabelle stehen
      nicht mehr in ihm (einzelne Grund-Codes nennt er weiter als Kontext,
      Verifikation A-2) (Suchlauf: `höchstens (220|120)` und „Deklaration dieses
      Vertrags“ im Vertrag 0). `harness/README.md` Zeile `make docs-check`,
      Spalte Bindung: Vertragsdatei · `SPEC-040`.
- [x] **Bedingung vor der Closure — `Schärft:`-Kante.** Eine Architect-ADR
      (Rollenwechsel, Baseline-Regelwerk `modul-08-agentenrollen.md`) stellt die
      Kante der ADRs aus **Bezug** zur neuen Kennung her, nach dem
      Muster, das die ADR des Pilots setzt; Status `Accepted`, Index
      nachgezogen. Kein Liefer-Punkt: der Implementer liefert die Kennungen, der
      Architect die Kante.
      *Beleg:* Commit `deecc466`,
      [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md),
      Status `Accepted`, `Schärft:` auf `SPEC-040`, Zeile im ADR-Index. Ihre
      Folgepflicht 1 und 2 ist im Commit nach `deecc466` umgesetzt (Spec-Absatz
      „Randformen der Verweis-Module“, §3, §6).
- [x] `make gates` grün, Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
      *Beleg:* `make gates > <Log im Scratchpad> 2>&1; echo $?` am Stand
      `780cd2b5`: Exit 0; gedruckt u. a. `coverage-gate: OK — Coverage 83.30%
      erfüllt Schwelle 80%`, `d-check: 1830 Datei(en) geprüft, 0 Befund(e)`,
      `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`,
      `generated-sync: OK`. Nach dem Commit dieses Belegs erneut gelaufen
      (Stempel), Ergebnis im Bericht des Implementers.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §7 | update | eine Zeile, `make docs-check` mit acht Modulen (Liefer-Punkt 1) |
| `harness/sensors/docs-check.md` | update | Verweis statt Schwelle und Randform (Liefer-Punkt 2) |
| `harness/README.md` §Sensors | update | Bindung „Spec-Kennung“ in der Zeile `make docs-check` (Liefer-Punkt 2) |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | `Schärft:`-Kante (Bedingung vor der Closure) — erledigt: [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) (`deecc466`), Folgepflicht in `64c68933` |
| `spec/pflichtenheft.md` §8 Historie | update (Plan-Nachzug, Implementer) | eine Zeile für `SPEC-040`, wie der Pilot für `SPEC-038`/`SPEC-039` |
| `harness/sensors/docs-check.md` §Sperren, §Bindung | update (Plan-Nachzug, Implementer) | Sperren je mit „was zu tun ist“ (Vorlage `gate.template.md`), zweite Sperre „kein git-Repository im Mount“ (gemessen, M1 unten); Bindung nennt `SPEC-040` und den Herkunfts-Anker der Closure-Notiz-Regel, den der alte Vertrags-Absatz trug |
| `AGENTS.md` §3.11, `.d-check.yml` Kommentare | **nicht geändert**, mit Grund | Suchlauf unten: §3.11 nennt die Vertragsdatei für „was der Sensor deckt und was nicht“, und das trägt §Grenze Punkt 8 weiter; die Kommentare begründen Werte der Eingabe und bleiben richtig |

Der Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) ist beim Start zu
messen; bewegte Eigenschaft ist der Ort von Schwelle und Randform von
`make docs-check` (Träger außerhalb des Vertrags: `AGENTS.md` §3.11,
`.d-check.yml` Kommentare).

**Lese-Umfang.** *Übernommen* aus der Messung des Implementers vor jeder
Änderung (Stand `b2451d70`): die Blöcke der acht Module in `.d-check.yml` etwa
293 Zeilen, die zu lesenden Quellen für `docs-check` zusammen etwa 1130 Zeilen
(der Pilot `slice-spec-festlegungen-harness-werkzeuge` etwa 790). *Gemessen*
vom Planner am Stand `cffa45be`: `wc -l .d-check.yml` = 347 (Datei samt
`vcs:`, `reviews:`, `commits:`, `trace:`); Block `structure:` Zeilen 137–296
(`grep -n '^[a-z-]*:' .d-check.yml`), also 160 Zeilen;
`wc -l harness/sensors/docs-check.md` = 247; die zehn ADR-Dateien aus
**Bezug** zusammen 2209 Zeilen (`wc -l`, ganze Dateien; welche Abschnitte die 1130
des Implementers zählen, ist nicht nachgemessen).

**Gegenprobe der Quellen und Anschluss-Frage** (Beleg zu Liefer-Punkt 1;
Regel: `.claude/commands/plan-welle.md` Schritt 6 · seit
slice-spec-festlegungen-harness-werkzeuge,
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`): Vor dem
Review trägt dieser Abschnitt eine Tabelle mit einer Zeile je normativem Satz
der Quellen — der Vertrag aus Liefer-Punkt 2, die Kommentare der Blöcke in
`.d-check.yml`, und die Entscheidungen der ADRs aus **Bezug** — und ihrem Ausgang *steht in
`SPEC-<NNN>` · bleibt im Vertrag (Grund) · entfällt (Grund)*; darunter je
Ausgang und Randfall der neuen Zeilen die Folge für den Anwender (besteht oder
nicht, welcher Fall gewinnt). Eine Folge, die keine Quelle trägt, entscheidet
der Architect vor der Closure, nicht eine Lesung im Auftrag.

**Messungen am Werkzeug** (Implementer, d-check-Digest aus `d-check.mk`,
`docker run --rm --network none -v <Baum>:/repo:ro <Digest>` direkt, ohne
`make`, außer wo genannt). M0 lief an einer `git archive`-Kopie ohne `.git` im
Scratchpad. **Vorfall bei M1:** M1 lief am Arbeitsbaum (`942ebf3f`), weil ein
`cd` in den Klon fehlschlug und der Befehl ohne Abbruch weiterlief. Die
Mutationen schrieben je eine Zeile per Umleitung `>>` an `README.md` und an
`.d-check.yml` — Text in eine Repo-Datei per Umleitung, verletzte Regel
[`AGENTS.md`](../../../../AGENTS.md) §3.1 („Text-Umschreiben im Repo ist Sache
der Datei-Werkzeuge des Laufs“). Beide wurden im selben Befehl sofort per
`git checkout` zurückgesetzt, `git status` danach leer; im Diff des Slice
steht keine Spur. Kandidat für das Register (Eintrag macht der Planner bei der
Closure): `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`. M2 bis M23
liefen an einem Klon von `942ebf3f` im Scratchpad, M24 bis M33 (Fixrunde) an
einem Klon von `72e57998`, M34 bis M43 (Verifikation) an einem Klon von
`e4f609d8`, je Mutation einzeln, danach `git checkout -- .`.
Gedruckt ist der Grund-Code bzw. die Summenzeile.

| # | Mutation | Ergebnis |
|---|---|---|
| M0 | Baum ohne `.git` | Exit 2, `d-check: error: kein lesbares git-Repository unter /repo` |
| M1a | Link auf eine fehlende Datei | Exit 1, `target-missing`, `1 Befund(e)`; über `make docs-check` Exit 2 |
| M1b | unbekannter Schlüssel in `.d-check.yml` | Exit 2, `d-check: error: .d-check.yml: yaml: unmarshal errors` |
| M2 | `require-pattern` der §7-Regel fehlt | Exit 1, `section-pattern-missing` |
| M3 | offene Aufgabe `- [ ]` in §7 | Exit 1, `section-tasks-open` |
| M4 | §7 leer | Exit 1, `section-empty` (und `section-pattern-missing`) |
| M5 | Kopfzelle `Titel` im ADR-Index umbenannt | Exit 1, `section-column-missing` („keine Tabelle des Abschnitts trägt eine Kopfzelle“) |
| M6 | Überschrift „3. Defaults und Konstanten“ umbenannt | Exit 1, `section-missing` („kein Abschnitt passt auf den Selektor“) |
| M7 | Kennung in Inline-Code ohne Link | Exit 0 |
| M8 | Kennung nackt in Prosa | Exit 1, `id-unlinked` |
| M9 | Kennung mit `<!-- d-check:ignore -->` in der Zeile | Exit 0 |
| M10 | ADR-Kennung in Inline-Code in der Architektur-Sicht | Exit 1, `matrix-forbidden` (Token) |
| M11 | Link aus der Architektur-Sicht in das ADR-Verzeichnis | Exit 1, `matrix-forbidden` |
| M12 | `slice-099` und `slice-abc-name` in der Architektur-Sicht | Exit 1, `matrix-forbidden` nur für `slice-099` |
| M13 | M10 mit `<!-- d-check:status-provenance -->` in der Zeile | Exit 0 (Befund, §6) |
| M14 | Token `docs/reviews/…` in einer ADR | Exit 1, `matrix-forbidden` adr → review |
| M15 | M14 mit `<!-- d-check:status-provenance -->` | Exit 0 (Befund, §6) |
| M16 | `slice-099` in der ADR 0039 (unter `matrix.exempt-paths`) | Exit 0 (Befund, §6) |
| M17 | Link auf ein fehlendes Heading und auf eine HTML-`id` | Exit 1, `anchor-missing` nur für das fehlende Heading |
| M18 | Link mit URL-Schema auf ein Ziel, das es nicht gibt | Exit 0 |
| M19 | `~/<Verzeichnis>/x` in Inline-Code | Exit 1, `hostpath-forbidden` |
| M20 | Tilde in URL, `~/.config/x`, nackte Tilde, `~<Benutzer>/x` | Exit 0 |
| M21 | M19 mit `<!-- d-check:ignore -->` | Exit 1, `hostpath-forbidden` |
| M22 | Pfad in die vendored Baseline mit dem Tag `v6.14.0` in Inline-Code | Exit 1, `version-stale` |
| M23 | M22 mit `<!-- d-check:ignore -->` | Exit 0 |
| M24 | offene Aufgabe `- [ ]` in §7 von `done/release-image-scan.md` (geschlossen, Name ohne `slice-`) | Exit 0 (Review F-1, R1) |
| M25 | `[SPEC-040](harness/conventions.md)` in `README.md` | Exit 0 (Review F-2, R3) |
| M26 | `[LH-QA-REL-001.a](spec/lastenheft.md)` in `README.md` | Exit 0 (Review F-2, R3) |
| M27 | Link, Ziel umbrochen: `[Text](harness/conventions` / `.md)` (zusammengesetzt existiert es) | Exit 1, `target-missing` (Review F-4, R8) |
| M28 | Link, Ziel-Anker umbrochen: `…conventions.md#gibt-es` / `-nicht)` | Exit 1, `anchor-missing` |
| M29 | Link, Linktext umbrochen, Ziel fehlt | Exit 0 |
| M30 | Link, Linktext umbrochen, Anker fehlt | Exit 0 |
| M31 | `ADR-0094` und `slice-099` in einem Fence der Architektur-Sicht | Exit 0 (Review F-7, R13) |
| M32 | Token `ADR-0038` (Status `Superseded`) in Inline-Code im Slice-Plan in `in-progress/`; in einer `observation.md`; in einem Review-Bericht mit `docs/reviews/*.md` aus `matrix.exempt-paths` entfernt | je kein `matrix-inactive`; im dritten Lauf 57 `matrix-inactive`, alle an Links, keiner an einem Token (Review F-3, R7) |
| M33 | `ADR-0094` in Inline-Code und der Marker als Text in Inline-Code in derselben Zeile der Architektur-Sicht | Exit 0 — der Marker wirkt auch in Inline-Code |
| M34 | Basislauf direkt, stdout und stderr getrennt in Dateien | Exit 0; stdout leer, stderr `d-check: 1835 Datei(en) geprüft, 0 Befund(e)` (Verifikation A-1) |
| M35 | Link auf eine fehlende Datei, direkt, Ströme getrennt | Exit 1; stdout `README.md:66` TAB `gibt-es-nicht.md` TAB `target-missing` TAB `Linkziel existiert nicht`; stderr `… 1 Befund(e)` |
| M36 | M35 über `make docs-check > Datei 2>&1` | Exit 2; Reihenfolge: `docker run …`, Summenzeile, Befundzeile, `make: *** [d-check.mk:20: docs-check] Fehler 1` |
| M37 | drei fehlende Linkziele über `make docs-check 2>&1`, fünf Läufe | je Exit 2, je Zeile 2 die Summenzeile `… 3 Befund(e)`, je letzte Zeile die Fehlerzeile von `make` |
| M38 | `~/<Verzeichnis>/x` in Inline-Code, direkt | Exit 1; stdout `README.md:66` TAB der Pfad TAB `hostpath-forbidden` — drei Felder, kein Text |
| M39 | unbekannter Schlüssel auf oberster Ebene von `.d-check.yml` | Exit 2; stdout leer; stderr zwei Zeilen: `d-check: error: .d-check.yml: yaml: unmarshal errors:` und `  line 348: field unbekannt not found` |
| M40 | unbekannter Schlüssel unter `scan:` | Exit 2; stderr dieselbe Form, zweite Zeile `  line 4: field unbekannt not found { … }` |
| M41 | `versions.current-from` auf eine fehlende Datei | Exit 2; stderr eine Zeile `d-check: error: versions.current-from nicht lesbar (…)` |
| M42 | unbekanntes Modul `gibtsnicht` in `modules:` | Exit 2; stderr eine Zeile `d-check: error: unbekanntes Modul "gibtsnicht" in der Konfiguration — gültig: …` (Verifikation A-5) |
| M43 | M39 über `make docs-check 2>&1` | Exit 2; die zwei Zeilen des Fehlers, danach `make: *** … Fehler 2` |

**Gegenprobe der Quellen** (Implementer). Ausgänge: **S** = steht in
`SPEC-040` · **V** = bleibt im Vertrag (Grund) · **E** = entfällt (Grund).
Quellen: der Vertrag am Stand `942ebf3f` (A), die Kommentare der Modul-Blöcke
in `.d-check.yml` (B), die Entscheidungen der ADRs aus **Bezug** (C). Ein Satz
ohne Wirkung auf das Werkzeug (Begründung, Verfahren, einmaliger Vorgang)
steht mit E und Grund.

| # | Satz der Quelle | Ausgang |
|---|---|---|
| A1 | Vertrag: je Klasse ein Grund-Code (`target-missing`, `anchor-missing`, `id-unlinked`, `matrix-forbidden`/`-inactive`, `version-stale`, `section-cell-*`, `section-forbidden`, `hostpath-forbidden`) | S (Tabelle „Befund je Modul“, ergänzt um die gemessenen `structure`-Codes M2 bis M6 und `target-untracked`); V nur die Klassen als Satz |
| A2 | Schwester-Artefakt als blankes Repo-Wort | V — Zitierregel der Doku, kein Prüfgegenstand; das Modul meldet den Pfad, nicht die Form |
| A3 | Module und Grenzen in `.d-check.yml`, Konfiguration ist Deklaration des Vertrags | S umformuliert: `.d-check.yml` ist die Eingabe, die Festlegung sagt, was ein Modul daraus als Befund meldet |
| A4 | Trägerdateien der `structure`-Regeln | S („Die Regeln gelten …“); die Closure-Notiz-Regel in der Fixrunde auf `done/slice-*.md` enger gefasst, wie die Quelle („je `done/slice-*.md`“) und das Werkzeug (M24) |
| A5 | Closure-Notiz-Regel `seit slice-001` | V — Herkunfts-Anker, in §Bindung (die Spec trägt keine Slice-Kennung) |
| A6 | Gate-Index: `Vertrag` ≤ 220, `Tut was` ≤ 120, beide ≥ 1 | E — Werte der Eingabe (`.d-check.yml`); die Regelart (Zelle unter Mindest-/über Höchstlänge) steht in S |
| A7 | die Regel misst beide Tabellen, Spalte über den Kopfzeilen-Namen | S |
| A8 | `seit slice-harness-readme-zellen-kuerzen` | V (Grenze 11) |
| A9 | Tabellen-Regeln über Mindestbreiten, Gate-Index über Mindest- und Höchstlänge, drei Regeln ohne Spalten-Knoten | S als Regelarten (Tabelle · Muster · Inhalt); E die Zuordnung je Regel — Werte der Eingabe |
| A10 | Keine Regel zählt Zeilen | S |
| A11 | Grenze 1: keine Richtung innerhalb der Spec, Review-Prüfpflicht | V; S „Innerhalb der Spec-Straten gibt es keine Regel“ |
| A12 | Grenze 2: `codepaths` aus | V; S (Fläche) |
| A13 | Grenze 3: Opt-in-Module nur über `doc-*`; `tracked` im Bündel | V; S „es laufen genau die Module der Liste“ |
| A14 | Grenze 4: `MR-*` nicht linkpflichtig, `tracked` prüft Getrackt-Status | V; S (Kennung ohne Muster) |
| A15 | Grenze 5: `.harness/**`, `**/*.template.md` ausgenommen | V; S (Fläche, `scan.ignore`) |
| A16 | Grenze 6: Verweisform als Textform, bereinigter Text, Reference-Style umgeht, `state.md`/`evidence` ohne Überschrift, Nachbar-Klassen | V; S der bereinigte Text (Muster) |
| A17 | Grenze 7: `ids` prüft Link, nicht Existenz; Inline-Code ungeprüft; kein Modul prüft Symbole; Bash-Hälfte deklariert | V alle vier; S die ersten zwei (Randformen der Verweis-Module) |
| A18 | Grenze 8: gemeldet werden absolute und Home-relative Pfade (drei Formen) in Prosa und Inline-Code | S (`hostpaths`) |
| A19 | Grenze 8: Tilde mit Benutzername deckt die Regel, das Modul nicht | V; S (kein Treffer, M20) |
| A20 | Grenze 8: `~/.config/…`, nackte Tilde, Tilde in URL meldet das Modul nicht | S (M20) |
| A21 | Grenze 8: Ventil verfügbar, nicht gesetzt | V; S („ausgenommen ist nichts“) |
| A22 | Grenze 8: Fence, relative Pfade, Nicht-Markdown, `scan.ignore` ungeprüft | V; S Fence und relativer Pfad |
| A23 | Grenze 8: Windows-Laufwerks-/UNC-Muster fest | S |
| A24 | Grenze 8: Regel strenger als Modul, Wächter Review, Träger-ADRs, kein Ausschlussblock | V |
| A25 | Grenze 9: mehrzeiliger Link wird nicht gemeldet | S, in der Fixrunde nach M27 bis M30 berichtigt: nur ein umbrochener Linktext wird nicht gemeldet, ein umbrochenes Ziel ist ein Befund; V die Lücke (umbrochener Linktext), Wächter Review |
| A26 | Grenze 9: einzeilig rot, zweizeilig grün (Exit 0), neun Links in fünf Dateien | E — Messung an einem älteren d-check (`slice-077`), am gepinnten Werkzeug nur für den Linktext wahr; ersetzt durch M27 bis M30 |
| A27 | Grenze 10: `trace:` kein Modul, RTM advisory | V; S („`trace:` ist kein Modul“) |
| A28 | Grenze 11: misst Länge, nicht Satz; `Bindung` ohne Höchstlänge | V |
| A29 | Grenze 11: Tabelle mit anderem Spaltennamen fällt aus der Regel; umbenannte Kopfzelle → `section-column-missing` | S (M5); V Verweis |
| A30 | Grenze 11: 121/222 Zeichen → `section-cell-oversized` „(Exit 2)“, leer → `-undersized`, 120 grün | E — Messprotokoll der Werte; die Regel steht in S, „Exit 2“ ist der Exit von `make` (S, Ausgänge; M1a) |
| A31 | Ausschnitt: Vollständigkeits-Zeile sagt etwas über den Ausschnitt | V; S (Ausgänge) |
| A32 | Exit 0/1/2 | S |
| A33 | `make doc-repair`, `make doc-doctor` | V |
| A34 | Sperre: Config-Fehler → Exit 2, kein stiller Rückfall | S; V als Sperre mit Abhilfe |
| A35 | `make doc-tracked` | V |
| A36 | Bindung | V, ergänzt um `SPEC-040` |
| B1 | `scan.ignore`: `.harness/**` ist tool-interne Ablage | S der Glob; E der Grund (Kommentar der Eingabe) |
| B2 | `ids`: je Stratum ein Muster, `.a`-Muster vor dem Basismuster | E — am Werkzeug ohne beobachtbare Wirkung: ein Link auf ein beliebiges Ziel genügt (M25, M26); S nur „ohne Link ist ein Treffer“ |
| B3 | `matrix`: Decken-Regel; `no-downward` nicht gesetzt, Delegation des Lastenhefts | S die Wirkung; E die Begründung (Kommentar) |
| B4 | Klasse `slice` eng gefasst, `*` quert kein `/` | E — Wert der Eingabe; S das Token nur mit Nummer (M12) |
| B5 | Klasse `welle` für die Abgrenzung | E — Wert der Eingabe |
| B6 | keine Regel `slice`/`welle` → `review` | S |
| B7 | `status.forbidden`, `allow-supersede-lineage`, `supersede-fields` | S — `matrix-inactive` nur an einem Link zwischen Dateien einer Klasse, ein Token ist kein Treffer, eine Datei ohne Klasse ist weder Quelle noch Ziel (M32; `docs/reviews/architect-verdict-matrix-inactive-nur-link.md`, `037a9d32`, Messungen T1, T2, L1–L7, Z1, Z2) |
| B8 | `matrix.exempt-paths`: Alt-ADRs, `docs/reviews`, `welle-3-results` | S — die Datei ist als Quelle ganz ausgenommen, als Ziel nicht (M16; [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 3); E die Gründe |
| B9 | `versions`: Pin trägt die Version aus `current-from` | S |
| B10 | `versions.exempt-paths`: Records, ADRs, MR-Einträge, pfadgebunden; nur `versions`, `links`/`anchors` weiter | S |
| B11 | `structure`: Spalte über Kopfzeilen-Namen, Zell-Länge wie dargestellt inkl. Markdown-Syntax | S |
| B12 | Gate-Index-Regel und ihre Grenze | S; V (Grenze 11) |
| B13 | §7-Regel: aktiviert mit der ersten Closure, `forbid-pattern` gegen Vorlagenrest | S die Regel; V der Anker; E der Grund |
| B14 | Verweisform-Regeln samt vier Grenzen | S (bereinigter Text); V Grenze 6 |
| B15 | E2E-Regel: Mindestbreiten, kein Maximum, Link trägt `ids`, rotet bei fehlender Datei | S (`section-missing` „keine Datei“); V Grenze 7 |
| B16 | `tracked`: existierendes, nicht getracktes Ziel; `exempt-targets: []` | S |
| C1 | 0072 E1: Aktivierung als Entscheidung mit Träger | E — Grund der Aktivierung, kein Prüfverhalten |
| C2 | 0072 E2: kein `scope`, kein `exempt-paths`, kein Zeilen-Marker; `ids`/`versions` kennen `d-check:ignore`, `hostpaths` nicht; Geltungsbereich globaler Scan | S (M9, M21, M23; „ausgenommen ist nichts“) |
| C3 | 0072 E2: ein nicht korrigierbarer Befund braucht eine Folge-ADR | E — Verfahren, kein Werkzeug |
| C4 | 0072 E3: Zitationsform besitzer-qualifiziert | E — abgelöst durch 0074 |
| C5 | 0072 E4, E5: 31 Stellen korrigieren, §3.11-Entwurf | E — einmaliger Vorgang bzw. Regel in `AGENTS.md` §3.11 |
| C6 | 0072 Config-Block: `hostpaths` in `modules:`, kein `hostpaths:`-Knoten | S |
| C7 | 0072 Entwurf §3.11 „Was der Sensor deckt“ | E — abgelöst durch 0075 E2 und 0160 E3; die geltende Fassung steht in S |
| C8 | 0072 §3.11-Entwurf: „Modul-Semantik steht einmal in `harness/sensors/docs-check.md`“ | E — überholt durch diesen Slice; Träger ist eine `Accepted`-ADR, bleibt als Geschichte ([`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md) Entscheidung 4) |
| C9 | 0074 E1–E3: Hausform, Link-Variante, Rest aus 0072 bleibt | V (A2) |
| C10 | 0075 E1: Regel deckt Fences, Sensor nicht | S (Fence kein Treffer); V Grenze 8 |
| C11 | 0075 E2, E3: Fassung 2, Lokator-Disposition | E — abgelöst durch 0160 E3 bzw. kein Werkzeug |
| C12 | 0160 E3: Gegenstand der Regel und Ausnahmen | S die Modul-Hälfte („Was der Sensor deckt“, Fassung 3); E die Regel-Hälfte (`AGENTS.md` §3.11) |
| C13 | 0160 E4: Ventil `hostpaths.exempt-targets` ungenutzt | S |
| C14 | 0094: Klasse `review`, Regel adr → review, Token fängt Pfad ohne Link | S (Link oder Token, M14) |
| C15 | 0094: kein Provenance-Marker-Escape für adr → review | S — in diesem Absatz abgelöst: der Marker hebt den Token-Befund jeder Regel auf, ein Link bleibt ein Befund (M13, M15; `ADR-0163` Entscheidung 2) |
| C16 | 0094: die Regel prüft eine Adresse, nicht die Semantik einer Umformulierung | E — Review-Prüfpflicht, kein Werkzeug |
| C17 | 0095: `docs/reviews/*.md` nur von der Status-Prüfung ausgenommen, Regel adr → review bleibt scharf | S — als Quelle ganz ausgenommen, als Ziel Gegenstand der Regel (M14, M16); für `docs/reviews/*.md` deckt der Wortlaut die Wirkung, weil keine Regel `review` als Quelle hat (`ADR-0163` Entscheidung 3). Der Satz in §Kontext über die „ausgehenden Token“ ist keine Entscheidung und bleibt als Geschichte; die Status-Prüfung liest nur Links (Architect-Verdikt `architect-verdict-matrix-inactive-nur-link`, §1) |
| C18 | 0097: Klasse `observation`, Regel observation → review | S (Regel der Eingabe, Befund `matrix-forbidden`) |
| C19 | 0099: keine Regel slice/welle → review, Klasse `welle` bleibt | S |
| C20 | 0156 E1, E2: `done/**` von `versions` ausgenommen, `links`/`anchors` ohne Ausnahme | S |
| C21 | 0156 E3, E4: Form-Korrektur eines Links in ein altes Tag | E — Doku-Verfahren (Zitat-Korrektur), kein Werkzeug |
| C22 | 0161 E2: ADRs und MR-Einträge von `versions` ausgenommen, pfadgebunden, Index bleibt geprüft | S |
| C23 | 0161 E3: `links`/`anchors` ohne Ausnahme, Form-Korrektur | S der erste Teil; E die Form-Korrektur (Verfahren) |
| C24 | 0161 E7: was ein Bump nachzieht; `versions` prüft alle `.md` unter `scan.roots`; Symlinks und `.harness/skills/**` liest es nicht | S `scan.roots`/`scan.ignore`; E die Liste des Bumps (Träger `harness/targets/pin-stale.md`) |

**Anschluss-Frage** (je Ausgang und Randfall die Folge für den Anwender):

| Fall | Folge |
|---|---|
| Befund und Fehler über `make` | beide Exit 2 von `make`; die letzte Zeile ist die Fehlerzeile von `make`, es gilt die Summenzeile (`<M> Befund(e)`) bzw. die Zeile `d-check: error: …` (M1a, M1b, M36, M43) |
| Befund-Zeilen und Summenzeile | Befunde auf stdout, Summe auf stderr (M34, M35); über `make` mit `2>&1` die Summe vor den Befunden, in fünf Läufen gleich (M37) |
| Konfigurationsfehler | kein Urteil über die Doku, kein Rückfall auf Defaults; der Lauf prüft nichts (M1b) |
| Baum ohne `.git` | Exit 2 vor jeder Prüfung (M0); ein `git archive`-Export ist kein prüfbarer Baum |
| Kennung mit Link auf ein anderes Dokument als ihr Definitions-Dokument | kein Treffer (M25, M26); ein falsches Ziel bleibt dem Review |
| Link mit Zeilenumbruch | im Linktext kein Befund (M29, M30), im Ziel ein Befund, auch wenn das zusammengesetzte Ziel existiert (M27, M28) |
| Kennung in Inline-Code | `ids` kein Treffer (M7), `matrix` Treffer (M10) — die zwei Module lesen Inline-Code verschieden |
| `d-check:ignore` | hebt `ids` und `versions` für die Zeile auf (M9, M23), `hostpaths` nicht (M21) und `matrix` nicht (*übernommen* aus `ADR-0163` Fitness Function Zeile 6) |
| `d-check:status-provenance` | hebt in seiner Zeile den Token-Befund von `matrix` für jede Regel auf (M13, M15); ein Link in derselben Zeile bleibt `matrix-forbidden` (`ADR-0163` Entscheidung 2). In `spec/` ist der Marker kein zulässiger Weg, das Werkzeug lässt ihn trotzdem zu; Wächter ist das Review (ebenda, Grenze) |
| Datei unter `matrix.exempt-paths` als Quelle | weder Status-Prüfung noch Regeln greifen (M16); als Ziel bleibt sie Gegenstand der Regel (M14). Für `ADR-0039`/`ADR-0041` angenommen, wer einen Pfad neu aufnimmt, nennt die Regel-Ausnahme in seiner ADR (`ADR-0163` Entscheidung 3) |
| Tabelle ohne die Spalte vs. Abschnitt ohne die Spalte | erste fällt aus der Regel, zweiter `section-column-missing` (M5) |
| Regel trifft keinen Abschnitt | `section-missing` (M6), auch wenn die Datei fehlt (B15) |

**Suchlauf** ([`AGENTS.md`](../../../../AGENTS.md) §3.13). Bewegte
Eigenschaft: der Ort von Befund-Codes, Schwellen und Randformen von
`make docs-check` (Vertrag → `SPEC-040`). Suchraum: der ganze Baum ohne
`docs/reviews/**`, `docs/plan/planning/done/**` und `.harness/baseline/**`;
Parent `942ebf3f`, `diff` ist der Arbeitsbaum nach dem Liefer-Commit.

```suchlauf
942ebf3f 26 -F 'sensors/docs-check.md' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 27 -F 'sensors/docs-check.md' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 3 -F 'Modul-Semantik' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 5 -F 'Modul-Semantik' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 3 -E 'höchstens (220|120)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 2 -E 'höchstens (220|120)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 1 -F 'Deklaration dieses Vertrags' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 0 -F 'Deklaration dieses Vertrags' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 9 -F 'hostpath-forbidden' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 8 -F 'hostpath-forbidden' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 2 -E 'acht (docs-check-)?Module' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 5 -E 'acht (docs-check-)?Module' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 15 -F '7-festlegungen-der-harness-werkzeuge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 25 -F '7-festlegungen-der-harness-werkzeuge' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
942ebf3f 0 -F 'SPEC-040' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 35 -F 'SPEC-040' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Gefunden: `harness/sensors/docs-check.md` nennen 26 Zeilen; bewegt hat sich
keine, weil die Abschnitte des Vertrags und die Nummern von §Grenze bleiben —
`AGENTS.md` §3.11 („was der Sensor deckt und was nicht, führt
`harness/sensors/docs-check.md` aus seiner Sicht“) bleibt wahr, Punkt 8 trägt
die Lücken zwischen Regel und Modul weiter und verweist für die Formen auf
`SPEC-040`; die Zustandsfelder im Register und `harness/targets/doc-trace.md`
zeigen auf §Grenze (Punkte 8 und 10 bleiben). „Modul-Semantik“ dreimal in
`Accepted`-ADRs: `ADR-0072` Zeile 198 sagt, die Modul-Semantik stehe einmal in
`harness/sensors/docs-check.md` — das ist mit diesem Slice überholt, die Datei
ist eingefroren ([`AGENTS.md`](../../../../AGENTS.md) §3.5); **gemeldet** an
den Architect (Frage in §6), nicht geändert. `ADR-0160` Zeile 252 nennt
„§Grenze Punkt 8: die Formen nach …“ — die Formen stehen jetzt in `SPEC-040`,
Punkt 8 verweist dorthin; eingefroren, mit gemeldet. „höchstens 220/120“: die
Zeile im Vertrag ist entfallen, die übrigen zwei Treffer (`bench-backfill`)
gehören nicht zur Eigenschaft. „hostpath-forbidden“ 9 → 8: der Vertrag nannte
den Code zweimal und nennt ihn nicht mehr, `SPEC-040` nennt ihn einmal; die
übrigen sechs Treffer stehen in `Accepted`-ADRs (Messprotokolle) und bleiben.
Nicht gefunden: ein
Träger außerhalb des Vertrags, der die Schwellen 220/120 oder die Exit-Codes
von `make docs-check` nennt; `.claude/commands/implement-slice.md` („Strenges
Doc-Gate“) und `.harness/skills/reviewer.md` beschreiben die Module, nicht den
Ort ihrer Festlegung, und bleiben richtig.

Nachgemessen nach `deecc466` (Arbeitsbaum mit der Folgepflicht aus
`ADR-0163`): die `diff`-Zeilen oben sind auf diesen Stand gezogen. Alle
Zuwächse gegenüber dem Stand des Liefer-Commits kommen aus `ADR-0163` selbst
(`sensors/docs-check.md` +1, „Modul-Semantik“ +2, der Anker +1, `SPEC-040`
+18; je `git grep -c` am Arbeitsbaum, nach Datei gelesen) — kein weiterer
Träger. Nach der Fixrunde: der Anker und `SPEC-040` je +2, beide im Vertrag
(neuer Grenze-Punkt 12 und §Sperren); die übrigen Zahlen unverändert.

**Umfang gemessen vor dem Schreiben** (Implementer, zur zweiten
Rückführungs-Bedingung in §4, `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`).
Satzenden (`.`, `!`, `?`, `:` gefolgt von Leerraum, ohne Fences) als
Näherung der normativen Sätze, mit derselben `awk`/`grep`-Zählung je Quelle:
`docs-check`: Vertrag 103, Kommentare `.d-check.yml` Zeilen 1–305 90,
Entscheidungen der ADRs 253 (0072 59, 0074 18, 0075 32, 0160 E3/E4 40, 0094
17, 0095 6, 0097 13, 0099 13, 0156 26, 0161 E2/E3/E7 29) — zusammen 446.
Pilot am Parent `fcff30ec`: `harness/targets/zitat-vergleich.md` 108,
Entscheidungen 0158 39, 0159 85, 0160 E1/E2 etwa 37, 0161 E4/E5 etwa 48 —
zusammen etwa 317. Verhältnis etwa 1,4 (*abgeleitet*), dasselbe wie das
Zeilen-Verhältnis, das der Planner bei der Rückführung abgewogen hat (§4). Die
Bedingung ist an dieser Messung **nicht neu eingetreten**; die Gegenprobe
oben hat 76 Zeilen. Ob sie eine Review-Sitzung übersteigt, sagt der
Review-Report (zweite Hälfte der Bedingung).

**Fixrunde** (Review-Report `docs/reviews/review-slice-spec-festlegungen-doku-gates.md`,
Commit `72e57998`; 5 HIGH, 3 MEDIUM, 2 LOW, 2 INFO):

- **F-1 (HIGH), erledigt.** `SPEC-040` nennt für die Closure-Notiz-Regel die
  Slice-Pläne `slice-*.md` direkt unter `done/` und sagt, dass ein geschlossener
  Plan mit anderem Namen nicht darunter fällt (M24); §Bindung des Vertrags
  ebenso. A4 nachgezogen.
- **F-2 (HIGH), erledigt.** Zeile `ids` und Randformen: ohne Link ist ein
  Treffer, das Ziel des Links prüft der Lauf nicht; der Satz „das erste Muster
  gilt …“ ist gestrichen (M25, M26). B2 und Anschluss-Frage nachgezogen.
- **F-3 (HIGH), angehalten.** Frage an den Architect in §6 (`ADR-0095`
  §Kontext sagt Tokens zu, das Werkzeug meldet nur Links, M32). **Danach
  erledigt** nach dem Architect-Verdikt
  `architect-verdict-matrix-inactive-nur-link` (`037a9d32`), Wortlaut aus
  dessen §3 in `SPEC-040` (Zeile `matrix-inactive`, Randformen, Historie) und
  im Vertrag (Grenze 12).
- **F-4 (HIGH), erledigt.** `SPEC-040` und Vertrag Grenze 9 nach M27 bis M30:
  umbrochener Linktext kein Befund, umbrochenes Ziel ein Befund. Die Aussage
  aus `slice-077` steht nicht mehr als Messung da (A25, A26).
- **F-5 (HIGH), nicht behebbar, erfasst.** Siehe F-6.
- **F-6 (MEDIUM), erledigt.** Vorfall bei M1 mit Weg (`>>`), verletzter Regel
  (`AGENTS.md` §3.1) und Rücksetzung oben eingetragen; Register-Kandidat
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` für die Closure.
- **F-7 (MEDIUM), erledigt.** Vertrag §Grenze Punkt 12: Marker in `spec/`,
  Regel-Ausnahme der ADRs 0039/0041, Token im Fence (M31; der Fence-Satz auch
  in `SPEC-040`).
- **F-8 (MEDIUM), erledigt.** §3-Zeile der Architect-ADR auf „erledigt“, C8
  mit `ADR-0163` Entscheidung 4.
- **F-9 (LOW), erledigt für §6; fremde Träger gemeldet.** §6 vierter Punkt
  belegt zusätzlich, dass keine Zeile unter `spec/` den Marker und zugleich
  ein Token `adr`/`slice` trägt:
  `git grep -n status-provenance -- spec/ | grep -cE 'ADR-[0-9]{4}|slice-[0-9]{3}'`
  druckt `0` (Arbeitsbaum nach `b36af1f1`); die fünf
  Pläne in `open/` stehen als Meldung an den Planner mit Frist Closure.
- **F-10 (LOW), erledigt.** Vertrag Grenze 8, §Ausgabe und §Sperren
  wiederholen die Randformen und Exit-Ursachen nicht mehr, sie verweisen auf
  `SPEC-040` und tragen nur noch, was das Grün nicht deckt bzw. was den Weg
  frei macht.
- **F-11 (INFO), erledigt.** Die Historie-Zeile von `SPEC-040` in §8 nennt
  Marker und `matrix.exempt-paths`; keine neue Zeile.
- **F-12 (INFO).** Keine Handlung; die zweite Rückführungs-Bedingung tritt
  von Seiten des Reviews nicht ein.

**Re-Review** (`docs/reviews/review-slice-spec-festlegungen-doku-gates-re-review.md`,
Commit `b36af1f1`):

- **N-1 (HIGH), erledigt.** Der Beleg zu F-9 oben nennt jetzt den Befehl aus
  §6 und seine gedruckte Zeile `0` statt `git grep -c status-provenance --
  spec/` (der druckt `spec/pflichtenheft.md:2` und trägt den Satz nicht).
- **N-2 (LOW), bekannt, nicht geändert.** Die Überschrift von Grenze 12 im
  Vertrag („zwei angenommene Lücken, der Fence und der Status nur am Link“)
  ist der Wortlaut aus §3.2 des Verdikts
  `architect-verdict-matrix-inactive-nur-link`; das Zählwort deckt vier
  Lücken nicht eindeutig. Der Ausgang steht bei der Closure an.
- **N-3, N-4.** Keine Handlung.

**Verifikation** (`docs/reviews/verify-slice-spec-festlegungen-doku-gates.md`,
Commit `e4f609d8`):

- **A-1, erledigt.** Absatz „Ausgänge“ von `SPEC-040` und die Historie-Zeile
  auf das gemessene Verhalten gezogen (M34 bis M43): Befundzeilen auf stdout,
  Summen- und Fehlerzeile auf stderr; `hostpaths` mit drei Feldern; Fehler des
  YAML-Parsers mit zweiter Zeile `  line <n>: …`; über `make` mit `2>&1` die
  Summenzeile vor den Befunden und die Fehlerzeile von `make` am Ende — der
  Beleg ist die Summen- bzw. `d-check: error`-Zeile, nicht die letzte Zeile.
  Vertrag §Ausgabe nachgezogen („Summenzeile“ statt „Schlusszeile“).
- **A-2, erledigt.** Beleg von Liefer-Punkt 2 nennt „die Befund-Code-Tabelle“.
- **A-5, erledigt.** Unbekanntes Modul in `modules:` als Exit-2-Fall in der
  Tabelle (M42).
- **A-3, A-4.** An die Closure (Planner), keine Handlung hier.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): `slice-spec-festlegungen-harness-werkzeuge`
liegt in `done/` — §7 und das Muster der `Schärft:`-Kante existieren; danach
Priorisierung durch den Auftraggeber, `in-progress/` trägt keinen Slice
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Festlegungen der acht `docs-check`-Module
  (`.d-check.yml` Zeile `modules:`, gezählt bei der Anlage) sprengen eine Review-Sitzung — dann `docs-check` allein und die zwei übrigen
  Gates in einem zweiten Slice. *Eingetreten, siehe Rückführungs-Grund unten.*
- `in-progress` → `next` (nach der Rückführung, zweites Mal zu groß): die
  Gegenprobe der Quellen (§3) für `docs-check` allein überschreitet eine
  Review-Sitzung, oder der Review-Report meldet, dass er `structure` nicht in
  derselben Sitzung prüfen konnte — dann `structure` (Festlegung der Regeln
  und Vertrags-Abschnitt) in einen eigenen Slice, und dieser Slice verweist
  für `structure` auf dessen Kennung.
- `in-progress` → `open` (blockiert): der Architect entscheidet, dass die Kante
  für diese Gruppe eine andere Form braucht als im Pilot, und die ADR liegt
  nicht vor.

**Rückführungs-Grund `in-progress` → `next` (2026-10-10).** Der Implementer hielt
vor jeder Änderung an; die Bedingung oben trat ein. Seine Messung (*übernommen*,
Stand `b2451d70`):

- `.d-check.yml` aktiviert acht Module (links, anchors, ids, matrix, versions,
  structure, hostpaths, tracked); die Konfiguration hat etwa 293 Zeilen,
  `structure` allein etwa 160 Zeilen und 9 Regeln.
- Verträge: `harness/sensors/docs-check.md` 247 Zeilen,
  `baseline-verify.md` 53, `commit-traceability.md` 11.
- Quell-ADRs: 0072, 0074, 0075, 0094, 0095, 0097, 0099, 0156, 0160
  (Entscheidung 3), 0161, 0045 und 0162 (Entscheidung 3).
- Zu lesen: für `docs-check` allein etwa 1130 Quellzeilen, für alle vier Gates
  etwa 1330; der Pilot hatte etwa 790.

Die Bedingung nannte „die zwei übrigen Gates“; zum Zeitpunkt der Rückführung
waren es drei, denn `make doc-immutable` kam mit der Closure des Pilots hinzu
(Übergabe F-5/F-13). **Neuschnitt:** dieser Slice behält `make docs-check` mit
allen acht Modulen; `make commit-traceability`, `make baseline-verify` und
`make doc-immutable` samt Übergabe F-5/F-13 gehen an
`slice-spec-festlegungen-commit-baseline-gates` (Datei in `open/`).
**`structure` bleibt hier** (Planner-Entscheidung): `docs-check` ist ein Gate
mit einem Vertrag und einer `SPEC`-Zeile; ein Schnitt nach Modul ließe den
Vertrag halb verweisend, halb selbsttragend zurück und wäre ein Schnitt nach
Bauteil statt nach Lieferwert. Die 1130 Zeilen liegen etwa 1,4-mal über dem
Pilot (*abgeleitet*); die Gefahr, dass das zu viel ist, trägt die zweite
Rückführungs-Bedingung oben und das Risiko in §6.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

Die zwei Liefer-Punkte sind abgehakt mit Beleg, die Architect-ADR ist
`Accepted`, `make gates` endet mit Exit 0, der Review-Report liegt vor und ist
aufgelöst, die Closure-Notiz trägt den Lerneintrag und jedes Risiko aus §6
seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Architect-Frage zur `Schärft:`-Kante.** Die ADRs aus **Bezug** sind
  `Accepted` und tragen ihr `Schärft:`-Feld unveränderlich; ob eine
  Sammel-ADR für `docs-check` oder eine ADR je Modul die Kante herstellt — und
  ob sie mit der Kante von `slice-spec-festlegungen-commit-baseline-gates`
  zusammengeht —, entscheidet der Architect. *Zu belegen durch:* die ADR.
  **Stand (Implementer, nach dem Liefer-Commit):** keine ADR liegt vor
  (`ADR-0162` ist die jüngste, sie schärft `SPEC-038`/`SPEC-039`). Offen für
  den Architect: (1) die Kante von `ADR-0072`, `ADR-0074`, `ADR-0075`,
  `ADR-0160` (E3, E4), `ADR-0094`, `ADR-0095`, `ADR-0097`, `ADR-0099`,
  `ADR-0156` und `ADR-0161` (E2, E3, E7) zu `SPEC-040` — Sammel-ADR oder je
  Modul; (2) ob dieselbe ADR auch `vcs:`/`make doc-immutable` des Folge-Slice
  trägt, dessen Festlegung es noch nicht gibt; (3) die drei Befunde der
  nächsten Zeile.
  **Ausgang: entfallen** — Grund:
  [`ADR-0163`](../../adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
  (`deecc466`) beantwortet die Frage: eine ADR für `SPEC-040` mit Herkunft je
  Teil, `make doc-immutable` nicht dort, sondern in der eigenen Architect-ADR
  von `slice-spec-festlegungen-commit-baseline-gates` (Entscheidung 1).
- **Umfang auch nach dem Neuschnitt an der Grenze.** `docs-check` allein liegt
  bei etwa 1130 Quellzeilen gegen etwa 790 im Pilot (*übernommen*, §4); der
  Block `structure:` trägt 160 Zeilen (gemessen, §3). *Zu belegen durch:* der
  Review-Report — prüft er die Festlegung in einer Sitzung, ist das Risiko
  entfallen; sonst greift die zweite Rückführungs-Bedingung in §4.
- **Vertrag und Werkzeug weichen ab**, sobald die Randform als Festlegung
  formuliert wird (etwa eine Reichweite von `hostpaths`, die
  [`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md) anders
  nennt als der Vertrag). *Zu belegen durch:* Gegenlesen der Festlegung gegen
  Vertrag und ADR-Kette; ein Befund wird ein eigener Slice.
  **Stand (Implementer):** drei Befunde, gemessen (§3, Messungen), nicht in
  `SPEC-040` geschrieben, weil die Lesung eine Entscheidung ist
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5): (a) der Kommentar
  `<!-- d-check:status-provenance -->` hebt einen Token-Befund von `matrix` auf,
  auch für adr → review (M15) und spec → adr (M13); `ADR-0094` §Entscheidung
  schließt den Marker für adr → review aus. (b) Eine Datei unter
  `matrix.exempt-paths` prüft das Modul als **Quelle** gar nicht, auch nicht auf
  die Regeln (M16, `slice-099` in der ADR 0039 bleibt grün); `ADR-0095` nennt
  nur die Status-Prüfung. (c) `ADR-0072` Zeile 198 nennt
  `harness/sensors/docs-check.md` als einzigen Ort der Modul-Semantik, und
  `ADR-0160` Zeile 252 die Formen in §Grenze Punkt 8; beide sind eingefroren
  und mit `SPEC-040` überholt. Ob (a) und (b) ein eigener Slice werden (Werkzeug
  oder Konfiguration) oder die Festlegung sie aufnimmt, und ob (c) eine
  Folge-ADR braucht, entscheidet der Architect.
  **Ausgang der drei Befunde: entfallen** — Grund: `ADR-0163`. (a) Entscheidung
  2 löst den Absatz von `ADR-0094` und den Satz von `ADR-0097` ab: Token mit
  Marker ist Herkunft, ein Link bleibt verboten; `SPEC-040` trägt den Satz.
  (b) Entscheidung 3 legt die Wirkung fest (als Quelle ganz ausgenommen) und
  nimmt sie für `ADR-0039`/`ADR-0041` an; `SPEC-040` trägt den Satz. (c)
  Entscheidung 4: Ort der Modul-Semantik ist `SPEC-040`, die zwei Sätze
  bleiben als Geschichte, keine Folge-ADR.
- **Die Spec darf nicht auf ADRs zeigen** (gemessen im Pilot: `matrix-forbidden`
  für Link und Kennung im Inline-Code). Die Festlegung trägt ihren Inhalt
  selbst. *Zu belegen durch:* `make docs-check` Exit 0 **und** keine Zeile
  unter `spec/`, die den Marker trägt und zugleich ein Token der Klassen `adr`
  oder `slice` — seit `ADR-0163` Entscheidung 2 hebt der Marker einen
  Token-Befund auch in `spec/` auf (M13), auch wenn er selbst in Inline-Code
  steht (M33), das Grün allein belegt die Aussage nicht mehr; Wächter ist dazu
  das Review des Spec-Diffs (Review F-9). Gemessen am Arbeitsbaum der
  Fixrunde: `git grep -n status-provenance -- spec/` trifft eine Zeile
  (`SPEC-040`, Randformen, der Marker als Text in Inline-Code) und den
  Historie-Eintrag; `git grep -n status-provenance -- spec/ | grep -cE 'ADR-[0-9]{4}|slice-[0-9]{3}'`
  druckt `0`.
  **Meldung an den Planner (Review F-9, Frist: Closure dieses Slice):** dieselbe
  Beleg-Zeile „Die Spec darf nicht auf ADRs zeigen … `make docs-check` Exit 0“
  steht in fünf Plänen unter `open/` — `slice-spec-festlegungen-code-gates`,
  `slice-spec-festlegungen-commit-baseline-gates`,
  `slice-spec-festlegungen-coverage-gates`,
  `slice-spec-festlegungen-kennungs-gates`,
  `slice-spec-festlegungen-pruefer-hooks`. Fremde Träger, hier nicht geändert;
  der Planner zieht sie nach oder benennt sie mit Adresse.
- **`matrix-inactive` gegen eine Quell-ADR (Review F-3) — angehalten, Frage an
  den Architect.** `SPEC-040` sagt „eine Datei verweist auf ein Ziel“ mit
  verbotenem Status. Das Werkzeug meldet nur einen **Link** und nur aus einer
  Datei, die einer `matrix`-Klasse angehört (Review R5–R7; M32: ein Token in
  einem Slice-Plan, einer `observation.md` und einem Review-Bericht bleibt
  grün, die 57 Befunde des dritten Laufs hängen alle an Links). `ADR-0095`
  §Kontext sagt dagegen, mit der Klasse `review` würden die „ausgehenden
  `ADR-\d{4}`-Token gegen `matrix.status.forbidden` geprüft“ — ein Wortlaut,
  der mehr zusagt als das gepinnte Werkzeug leistet. Die Zeile in `SPEC-040`
  ist deshalb **nicht** umgeschrieben ([`AGENTS.md`](../../../../AGENTS.md)
  §3.5); der Architect entscheidet, ob die Festlegung auf „Link aus einer
  Datei einer Klasse“ geht und wie mit dem Satz in `ADR-0095` umzugehen ist.
  **Ausgang: entfallen** — Grund: Architect-Verdikt
  `docs/reviews/architect-verdict-matrix-inactive-nur-link.md` (`037a9d32`):
  die Festlegung hat falsch behauptet, `matrix-inactive` meldet nur einen Link
  zwischen Dateien einer Klasse; keine Folge-ADR, `ADR-0095` bleibt
  unverändert. Die Zeile in `SPEC-040`, der Satz in den Randformen, die
  Historie-Zeile und Grenze 12 des Vertrags tragen den Wortlaut aus §3 des
  Verdikts.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Gegenstand:** <übernommen von `slice-<Kennung>` | entfallen: <Grund>>
  *(nur beim Ausgang ohne Arbeit; sonst Zeile löschen)*
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<KUERZEL>/<slug>` (<slice-kennung-a>, <slice-kennung-b>, <slice-kennung-c> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-<Kennung>.md` | `evidence/slice-<Kennung>.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-<Kennung> (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/` und
der ADR-Index; die Modus-Deklaration führt nur die Default-Sub-Area `*`
(Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** bei Anlage gelesen:
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (bei Anlage
2×, mit `slice-spec-festlegungen-harness-werkzeuge` 3× und verkörpert) trifft
Liefer-Punkt 1 — eine Festlegung, die eine Randform offen lässt, wäre ein
Auftreten nach der Verkörperung; Prüfschritt in §3.

Nachgeholt vom Implementer am gemergten Stand `942ebf3f` (alle Einträge unter
`observations/BEO-PGC/` mit Zustand `offen`, Zähler = Zahl der Dateien in
`evidence/`). Die einzige Sub-Area ist `*` (`PGC`); getroffen heißt: der
Eintrag betrifft Festlegung, Vertrag, das Doku-Gate oder die Übergabe dieses
Slice.

- `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (2×) — trifft §4: die
  zweite Rückführungs-Bedingung ist vor dem Schreiben gemessen (§3, „Umfang
  gemessen vor dem Schreiben“), nicht erst bei der Closure. Mit diesem Slice
  kein Auftreten.
- `BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open` (1×) — trifft
  Liefer-Punkt 1: `SPEC-040` nennt die gemessenen Formen ohne Befund
  (mehrzeiliger Link, Inline-Code bei `ids`, `d-check:ignore`, Fence und
  Tilde mit Benutzername bei `hostpaths`); zwei fail-open-Formen (§6, Befund a
  und b) stehen als Frage an den Architect, nicht unbenannt.
- `BEO-PGC/rollen-uebergabe-ohne-committetes-artefakt` (2×) — trifft die
  `Schärft:`-Kante: die Fragen an den Architect stehen in §6 dieses Plans, nicht
  nur im Bericht; ein dritter Beleg entstünde, wenn eine Vorgabe des Architect
  ohne Artefakt umgesetzt würde.
- `BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor` (1×) — Randform von
  `ids`, die die Quellen nicht tragen; nicht in `SPEC-040` aufgenommen (keine
  Quelle, kein Messfall in diesem Slice), benannt für den Review.
- `BEO-PGC/bindung-spalte-uneinheitlich-tief` (1×) — trifft Liefer-Punkt 2:
  die Zeile `make docs-check` trägt Vertragsdatei und Spec-Kennung inline,
  wie die Zeile `make zitat-vergleich`.
- `BEO-PGC/gate-prueft-existenz-nicht-passung` (2×) — gilt Gates, die
  Kennungsmengen vergleichen; `ids` prüft den Link, nicht die Existenz (A17) —
  dieselbe Klasse auf der Doku-Seite, als Randform in `SPEC-040` benannt; kein
  Auftreten.

Keine der übrigen offenen Einträge betrifft Spec, Harness-Vertrag oder
Doku-Gate; keiner erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
