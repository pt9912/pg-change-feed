# Review-Report: slice-harness-baseline-v6-14-1 — 2026-10-06

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `11a5bac5..4045dc4f` (14 Commits, darunter `60cc0ed1` mit
[`ADR-0156`](../plan/adr/0156-versions-gate-nimmt-done-records-aus.md) und
`e2666499` mit dem Entfernen von v6.14.0), dazu auf Wunsch des Auftraggebers der
während des Laufs gelandete Commit `b6c5b419` (Regel-Symlinks unter `.claude/rules`)

**Skill:** `.harness/skills/reviewer.md` @ b6c5b419
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-baseline-v6-14-1` (Stand `4045dc4f`)
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P8, Entscheidung 7)
- [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) (Zitat-Korrektur, §Entscheidung 1 und 4)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen)
- [`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md) (korrigierte `Accepted`-ADR)
- [`ADR-0156`](../plan/adr/0156-versions-gate-nimmt-done-records-aus.md) (neu)
- `harness/targets/pin-stale.md` §Bump-Ablauf, Schritte 1–3
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.5, §3.7, §3.9, §3.12, §3.13, §3.15)
- Register-Einträge `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`,
  `BEO-PGC/exemption-ohne-reifegrenze`, `BEO-PGC/record-rueckwirkend-umgeschrieben`
- keine `LH-*`-Kennung berührt

**Nachgefahren (Exit direkt, Ausgaben im Scratchpad des Laufs):**

- Release `v6.14.1` von `pt9912/ai-harness-course` (veröffentlicht
  2026-10-06T05:16:10Z) per `gh release download`: `sha256sum -c SHA256SUMS`
  → `lab-regelwerk.zip: OK`, Exit 0, Wert `98862525…7dc7bc` gleich dem Plan.
  Entpackt (`unzip`, 54 Dateien) gegen `.harness/baseline/v6.14.1`:
  `diff -r` meldet nur `SHA256SUMS` als zusätzliche Datei im Repo (die erzeugt
  `vendor-baseline`). Die Repo-`SHA256SUMS` gegen den entpackten Baum:
  `sha256sum -c` Exit 0, 54 Zeilen `: OK` bei 54 Zeilen. `git ls-files` 55.
- `make baseline-verify`: Exit 0, gedruckt `baseline-verify: v6.14.1 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)`.
- Normalisiertes Delta neu erzeugt (alter Baum per `git archive 11a5bac5`, neuer aus dem Release-Asset):
  `templates/` Exit 1, 46 Diff-Zeilen, 4 Dateien, 22 Zeilen `<`/`>`;
  `regelwerk/` Exit 1, 10 Diff-Zeilen, 2 Dateien, 4 Zeilen; roh `diff -rq` 32 Zeilen.
  Gleich den Zahlen im Plan. Inhalt Zeile für Zeile gegen die
  Tabelle in Schritt 1 gelesen: stimmt.
- Stichprobe Nachzug (b): `harness/conventions.md` Z. 98 gleich Z. 107 der
  v6.14.1-`conventions.template.md`. Nachzug (c): das neue Muster aus
  `implement-slice.md` trifft die `Auslöser:`-Zeile der Slice-Vorlage alt und
  neu (Exit 0, Zeile 194), eine gefüllte Zeile nicht (Exit 1).
- `make suchlauf-nachmessen PLAN=<Plan>`: Exit 0, `14 Zeilen stimmen`. Die
  Abweichungen von der Planner-Erwartung sind begründet und nachgemessen:
  Zeile 1 (10) = `ADR-0095` §Geschichte Z. 156 + zwei Register-Records + sieben
  Zeilen in `ADR-0156` (`git grep -c` 7); Zeile 6 (2) = Record-Zeilen 277 und 337.
- [`ADR-0156`](../plan/adr/0156-versions-gate-nimmt-done-records-aus.md), Fitness Function:
  V0 auf einem Klon von `aa0a44bc` reproduziert (Exit 1, `2 Befund(e)`, Zeilen
  277 und 337); M1 und M2 auf einem Klon von `4045dc4f` mit der echten
  `.d-check.yml` gefahren (M1: Pin in `in-progress/` → Exit 1, 1 `version-stale`;
  M2: Pin in `done/` → Exit 0, 0 Befunde). Die Messung zu Zeile 277:
  `git grep -c … d79b7ebd -- <v6.14.0-Pfad>` druckt 1 (Exit 0), mit dem
  v6.14.1-Pfad Exit 1. Entscheidung 2 („ältere Versionen trifft die Suche in
  `done/` gar nicht“): `git grep -nE 'baseline/v[0-9]+\.[0-9]+\.[0-9]+/'` in
  `done/` ohne v6.14.x → 0.
- Referenten der Zitat-Korrekturen, je Tag roh und normalisiert verglichen
  (`cmp`): `templates/.d-check.yml` 0/0, `README.template.md` (Planung) 0/0,
  `MR-NNN-titel.template.md` 0/0, `harness/README.template.md` 0/0 (10962 Byte,
  210 Zeilen), `regelwerk/modul-13-quality-gates.md` **1**/0 (nur der
  Versionsstring; Anker `guard-haertung` vorhanden).
- `make doc-commits RANGE=11a5bac5..HEAD`: Exit 0, 0 Befunde.
  `make doc-immutable RANGE=11a5bac5..HEAD`: **Exit 2**, 4 Befunde (F-3).
- Symlinks (Nachtrag des Auftraggebers): siehe F-1, F-4.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Suchlauf trägt den Satz von Liefer-Punkt 2 („Jede lebende Nennung von `v6.14.0` zeigt auf v6.14.1“) nicht: `git grep` überspringt getrackte Symlinks (Modus 120000) im Baum- und im Arbeitsbaum-Modus. Gemessen: `git grep -c 'v6\.14\.0' 4045dc4f -- .claude/rules` Exit 1, obwohl `git cat-file -p 4045dc4f:.claude/rules/modul-05-planning-harness.md` den Pfad `../../.harness/baseline/v6.14.0/regelwerk/…` druckt; am Stand `4045dc4f` zeigen vier Regel-Symlinks ins Leere (`find . -xtype l` → 4). Behoben ist der Gegenstand in `b6c5b419`, der Plan nennt aber weder die Symlinks noch den Commit (Liste in Liefer-Punkt 2, Tabelle §3, Suchlauf ohne Messzeile für Symlink-Ziele), und der Haken `[x]` steht auf einem Beleg, der diese Träger nicht sieht. | `AGENTS.md` §3.13 (Suchform: ganzer Baum) · §3.12 | `docs/plan/planning/in-progress/slice-harness-baseline-v6-14-1.md` · „Suchraum der ganze Baum ohne `.harness/baseline/**`“ | ja — `git ls-files -s \| awk '$1==120000'` mit `git cat-file -p` je Blob am Stand `4045dc4f` | Beleg trägt seinen Satz nicht |
| F-2 | MEDIUM | Die Zitat-Korrektur an `ADR-0095` sitzt in §Verglichene Alternativen, das `AGENTS.md` §3.5 und `ADR-0073` §Entscheidung 1 wörtlich als unberührbar führen; zulässig ist sie nur nach der Kurzform und der Lesart des Architect-Audits (`9f1eb320`, §4.2/§7). Der Referent ist gemessen unverändert (`templates/.d-check.yml` roh `cmp` 0). Es ist das dritte Auftreten derselben Konfliktklasse, damit ist der Konflikt-Pfad über den Architect Pflicht (Modul 8). Die Folge-Slice-Datei, die §6 für die Closure zusagt, liegt noch nicht im Lifecycle: `git grep` nach `abschnitte-kurzform` in `open/`, `next/` und `in-progress/` trifft nur diesen Plan. | `AGENTS.md` §3.5 · [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1 · `v6.14.1` · `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz | `docs/plan/adr/0095-review-klasse-exempt-status-check.md` · „Baseline-Vorlage `.harness/baseline/v6.14.1/templates/.d-check.yml`“ | nein — Auslegungsfrage zweier Norm-Texte, kein Gate | Zitat-Korrektur-Reichweite: Abschnitte-Liste gegen Kurzform |
| F-3 | MEDIUM | Die Pin-Umstellung in `MR-001` bis `MR-004` (`5d8855d9`) ändert den Core von Dateien, die das `vcs`-Modul append-only hält: `make doc-immutable RANGE=11a5bac5..HEAD` endet mit Exit 2, vier `core-drift-vcs`. Der Plan führt die MR-Dateien als gewöhnliche Pins (Liefer-Punkt 2), nicht als Zitat-Korrektur an einer immutablen Datei, und der Commit nennt nur `ADR-0051`, nicht `ADR-0073`. Dasselbe Muster liegt beim vorigen Bump vor (`make doc-immutable RANGE=484d20ec~1..484d20ec` ebenfalls Exit 2, 4 Befunde). Folge: Der Verifier-Lauf nach seinem Briefing (`doc-immutable` je Slice-Umfang) ist bei jedem Bump rot, ohne dass ein Träger die Ausnahme benennt. | [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 3 (Belegform) · `.d-check.yml` Block `vcs:` | `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md` · „§Spec-Straten](../../.harness/baseline/v6.14.1/“ (gleichartig MR-002 bis MR-004) | ja — `make doc-immutable RANGE=11a5bac5..HEAD` | Zitat-Korrektur an immutabler Datei ohne Beleg |
| F-4 | LOW | Weder der Bump-Ablauf in `harness/targets/pin-stale.md` noch `make docs-check` sehen ein Symlink-Ziel: der Ablauf kennt keinen Schritt für getrackte Symlinks auf den alten Tag, und d-check bleibt bei vier hängenden Symlinks grün (Klon von `4045dc4f`: Exit 0, `1751 Datei(en) geprüft, 0 Befund(e)`). Es ist das zweite Auftreten: beim Bump auf v6.14.0 kamen die Symlinks ebenfalls in einem eigenen, in keinem Plan genannten Commit nach (`2d51b380`). Die Datei liegt außerhalb des Diffs. Meldung an den Planner mit Frist Closure dieses Slice (`AGENTS.md` §3.13). | `AGENTS.md` §3.13 (Träger in fremder Datei) | `harness/targets/pin-stale.md` · „Bump-Ablauf: Vergleich vor dem Löschen der alten Baseline“ | ja — `find . -xtype l` nach dem Löschen des alten Tags | Bump-Ablauf ohne Symlink-Ziele |
| F-5 | LOW | `ADR-0156` stützt die stehenbleibenden Korrekturen auf „der Referent ist je Zeile `cmp`-gleich“ und lässt in Entscheidung 3 das Umziehen eines Linkziels nur bei gleichem Referenten „(`cmp`)“ zu. Der Plan hat aber an versions-normalisierten Kopien gemessen. Für `regelwerk/modul-13-quality-gates.md` (Link in `done/slice-harness-guard-blocked-python.md`) ist der rohe `cmp` 1 und nur der normalisierte 0. Die Aussage nennt die Instanz ihrer Messung nicht vollständig. | `AGENTS.md` §3.12 („Verfasser einer ADR“) | `docs/plan/adr/0156-versions-gate-nimmt-done-records-aus.md` · „der Referent ist je Zeile `cmp`-gleich“ | ja — `cmp` roh gegen normalisiert auf beiden Tags | ADR-Aussage breiter als ihre Messung |
| F-6 | INFO | `ADR-0156` ist kein weiteres Auftreten von `BEO-PGC/exemption-ohne-reifegrenze`, sondern die Form mit Reifegrenze, die der Eintrag empfiehlt: Die Ausnahme greift erst nach dem `git mv` nach `done/`, ein falscher Pin in `in-progress/` bleibt rot (M1 nachgefahren). Die ADR nennt den Eintrag nicht. Ihr Satz „dieselbe Begründung und derselbe Mechanismus wie bei `docs/reviews/**`“ zieht ausgerechnet das Vorbild heran, das der Eintrag als Fall ohne Reifegrenze führt. | `BEO-PGC/exemption-ohne-reifegrenze` | `docs/plan/adr/0156-versions-gate-nimmt-done-records-aus.md` · „Dieselbe Begründung und derselbe Mechanismus wie bei `docs/reviews/**`“ | nein | Ausnahme-Begründung über ein weiteres Vorbild |
| F-7 | INFO | Die fünf Record-Korrekturen (`54c6e632`) genügen `ADR-0073`: Die Referenten sind gleich (siehe Nachgefahren), die Haken sind unberührt, und der Anker löst auf. Drei davon stehen in `[x]`-Punkten und lesen sich jetzt anachronistisch, etwa der v6.14.0-Bump-Slice „gegen `.harness/baseline/v6.14.1/…` abgeglichen“. Mit `ADR-0156` wären sie für `versions` nicht mehr nötig gewesen; `ADR-0156` §Konsequenzen lässt sie bewusst stehen. | [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 4 · `BEO-PGC/record-rueckwirkend-umgeschrieben` | `docs/plan/planning/done/slice-baseline-6-14-0-dokumente-nachziehen.md` · „`.harness/baseline/v6.14.1/templates/docs/plan/planning/README.template.md`“ | nein | Record-Korrektur ohne Gate-Not |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.harness/baseline/v6.14.1/` gegen das Release-Asset (sha256, Entpacken, `SHA256SUMS`, `make baseline-verify`) | geprüft, ohne Befund |
| Commit-Folge (`4da92b66` nur 55 A, `e2666499` nur 55 D, Zitat-Korrekturen je Datei-Klasse getrennt, `AGENTS.md` §3.3) | geprüft, ohne Befund |
| Bump-Ablauf Schritte 1–3 im Plan (Befehl und gedruckte Zahl, Delta neu erzeugt und gleich, Adaptionen-Durchgang) | geprüft, ohne Befund über F-1 und F-4 hinaus |
| `harness/conventions.md` §Baseline (Stand, Datum, Release-URL, Stand-Zeile gleich `regelwerk/README.md` Z. 3) und MR-000 (Nachzug b) | geprüft, ohne Befund |
| `.claude/commands/implement-slice.md` Schritt 24 (Nachzug c, Probe alt/neu/gefüllt) | geprüft, ohne Befund |
| `AGENTS.md` §1 (Release-URL) und §3.7 (Vorlagen-Klammer entfernt, Beispieltext und Regel unverändert) | geprüft, ohne Befund |
| `.claude/agents/{architect,reviewer,verifier}.md`, `.harness/skills/{reviewer,closure-note-reviewer}.md`, `harness/sensors/baseline-verify.md` (Messzeile gleich dem eigenen Lauf) | geprüft, ohne Befund |
| `.d-check.yml` Block `versions:` gegen `ADR-0156` Folgepflicht (Liste wörtlich gleich; Kommentar mit einer Kennung, Klassen Zusage/Abgrenzung, ohne Chronik, `AGENTS.md` §3.7; widerspricht dem `matrix`-Kommentar „kein Freibrief“ nicht, weil er sich auf `versions` beschränkt) | geprüft, ohne Befund |
| `ADR-0156` (Begründung, ehrlich benannte Grenze, gemessen/hergeleitet/übernommen getrennt, Fitness-Function V0/M1/M2 nachgefahren, Messung zu Zeile 277, Index-Zeile) | geprüft, ohne Befund über F-5 und F-6 hinaus |
| Form-Korrektur `docs/reviews/architect-verdict-aufschub-adresse-verfaellt.md` Z. 114 (`cc9ecb27`): Link → Inline-Code, gleicher Referent, alter Tag bleibt. Der Wegfall von `../../` ist die wurzelrelative Form, die Inline-Code-Pfade im Repo tragen, und entspricht der Zitier-Form der Report-Vorlage („Tag + Pfad in Inline-Code“). Das ist die „Form einer gebrochenen Referenz“ nach `ADR-0073` §Entscheidung 1 | geprüft, ohne Befund |
| `b6c5b419` (vier Symlinks auf v6.14.1, Modus 120000 bleibt, Ziele existieren, `find . -xtype l` → 0) | geprüft, ohne Befund am Commit selbst; Plan-Träger siehe F-1 |
| Suchlauf (`make suchlauf-nachmessen`, 14/14; Abweichungen von der Erwartung begründet und nachgemessen) | geprüft, ohne Befund über F-1 hinaus (Zahlen stimmen, der Suchraum hat eine Lücke) |
| Traceability (`make doc-commits RANGE=11a5bac5..HEAD` 0 Befunde) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Zitat-Korrektur-Reichweite: Abschnitte-Liste gegen Kurzform · Zitat-Korrektur an immutabler Datei ohne Beleg · Bump-Ablauf ohne Symlink-Ziele · ADR-Aussage breiter als ihre Messung · Ausnahme-Begründung über ein weiteres Vorbild · Record-Korrektur ohne Gate-Not

## Verdikt

**Merge-blockierend:** ja. F-1 (HIGH) braucht eine Fixrunde am Plan: Die
Symlinks und `b6c5b419` gehören in Liefer-Punkt 2 und §3, und eine Messung der
Symlink-Ziele gehört an beide Stände. F-3 färbt den `doc-immutable`-Lauf des
Verifiers rot. F-2 läuft als dritter gleicher Konflikttyp den Konflikt-Pfad über
den Architect; Übergabe-Artefakt ist dieses Finding, das Verdikt ein Artefakt
(Folge-ADR zur Abschnitte-Liste von `ADR-0073`).

**DoD-Haken „Review durchgeführt“:** nicht nachgezogen, weil eine Fixrunde
nötig ist (Skill §DoD-Checkbox-Nachzug ohne Fixrunde). Der Haken wird regulär in
der Fixrunde gesetzt.

**Lauf-Notiz (`AGENTS.md` §3.15):** Der PreToolUse-Guard verweigerte meinen
Aufruf `sed -i 's/v6\.14\.0/vX/g'` auf Kopien im Scratchpad. Wortlaut: „In-place
text tools (sed -i, perl -i, awk -i inplace) … are blocked, also on a scratch
copy (AGENTS.md Hard Rule 3.1) … write to stdout: sed s/a/b/ file >
/path/to/scratch-copy“. Danach lief die Normalisierung in der Form, die
`AGENTS.md` §3.1 und die Ablehnung selbst nennen (`sed … Datei > Kopie` im
Scratchpad). Keine Repo-Datei wurde von einem Host-Werkzeug geschrieben.

**Übergabe:** Findings gehen an den Implementer (F-1, F-3), F-2 an den
Architect, F-4 an den Planner (fremde Datei, Frist Closure). Die
Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ist ein
Lauf-Beleg und ersetzt die Verifikation nicht.
