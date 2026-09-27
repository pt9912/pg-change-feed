# Verifikations-Report: slice-transformationen-betriebsdoku — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-betriebsdoku.md`](review-slice-transformationen-betriebsdoku.md).
Formvorbild dieses Reports: [`verifikation-slice-backfill-e2e.md`](verifikation-slice-backfill-e2e.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`); der
Report folgt dem Formvorbild, verschlankt auf den Umfang eines reinen Doku-Slice (kein Code-Diff,
keine Mutationen von Eingabeseiten sinnvoll — der Gegenstand ist Prosa und ein Zitat).

**Gegenstand:** `git diff b80f45cf..994b887a` — dritter und letzter Slice von
`welle-transformationen`: Feature-Commit `97de226c` (Handbuch-Abschnitt, SDK-Suchlauf), Review
`c621bf36` (1 HIGH F-1, 1 MEDIUM F-2), Fixrunde `994b887a` (beide Findings behoben, DoD-Checkbox
nachgezogen). Zwei Dateien geändert: `docs/plan/planning/in-progress/slice-transformationen-betriebsdoku.md`,
`docs/user/benutzerhandbuch.md`. Dieser Lauf ändert weder Code noch Plan, Spec oder Doku; er
schreibt nur diesen Report.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Stand aller Läufe: `HEAD` = `994b887a`, Arbeitsbaum sauber (`git status` vor und nach den Läufen).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make docs-check` | **EXIT=0** | `d-check: 1364 Datei(en) geprüft, 0 Befund(e)` |
| `make gates` | **EXIT=0** | `d-check: 1364 Datei(en) geprüft, 0 Befund(e)` (docs-check) · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` (coverage-gate/baseline-verify liefen im selben Aufruf ebenfalls grün, keine eigene Log-Zeile im gezeigten Ausschnitt nötig — Gesamt-Exit 0 ist maßgeblich) |
| `make doc-commits RANGE=b80f45cf..994b887a` | **EXIT=0** | `d-check: 1364 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=b80f45cf..994b887a` | **EXIT=0** | `d-check: 1364 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-transformationen-betriebsdoku.md` | **EXIT=0** | 8 Zeilen `OK`, u. a. `soll=11 ist=11 diff 11 -E 'exclude_column\|enable_table' -- docs/user`; `suchlauf-nachmessen: 8 Zeilen stimmen` |

Kein Code-Diff, keine Mutation nötig: der Gegenstand ist ein Prosa-Abschnitt und ein Zitat; die
Prüfung ist Textvergleich gegen die genannten Quellen (§2).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | SDK-Beleg: Row-Image-Schlüssel opak, kein Zugriff auf festen Schlüssel im SDK-Produktivcode | **erfüllt** | `git diff --stat c621bf36..994b887a -- sdks/` leer — die Fixrunde ändert keinen SDK-Code, der Suchlauf des Feature-Commits (137/137 an beiden Ständen, vom Reviewer bereits Fundstelle-für-Fundstelle gelesen) bleibt unberührt; von mir per `make suchlauf-nachmessen` erneut auf 137/137 bestätigt (§1) |
| 2 | Betriebsdokumentation im Benutzerhandbuch (neuer Abschnitt §4, K1–K4, Wirkung, Dauerhaftigkeit, Nichtanwendbarkeit+Abhilfe, Backfill-Form, Reihenfolge der Aufrufe) | **erfüllt** | Abschnitt „Transformationsregel konfigurieren“ vorhanden (`docs/user/benutzerhandbuch.md:328ff`); F-1 (Fehlertext-Zitat) und F-2 (drittes Beispiel) aus dem Review behoben, beide selbst gegen ihre Quellen nachgelesen (§4) |
| 3 | Jede Zahl/Wirkungs-Aussage trägt ihren Ursprung (§3.12); Änderungshistorie trägt eine Zeile | **erfüllt** | Zeile 1.68 in der Änderungshistorie vorhanden, im selben Diff wie `Version: 1.68` im Kopf; der korrigierte Fehlertext trägt weiterhin „(gemessen im Review-Report …)“ als Ursprungs-Anker |
| 4 | `make gates` grün, Exit-Code ungefiltert gesichert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH | **erfüllt** | `review-slice-transformationen-betriebsdoku.md` liegt vor (1 HIGH, 1 MEDIUM); beide in der Fixrunde behoben, von mir Zeile-für-Zeile nachgemessen (§4) — kein offenes HIGH/MEDIUM |
| 6 | §3.13-Suchlauf: Feld trägt Gefundenes und Nichtgefundenes, beide Stände gemessen | **erfüllt** | 8 Zeilen des Suchlauf-Felds selbst nachgefahren (§1); die Fixrunde hat die vierte Zeile korrekt von 9 auf 11 nachgezogen (§5) |
| 7 | Doku-Update: der Slice ist das Doku-Update | **erfüllt** | einziger inhaltlicher Diff-Gegenstand ist `docs/user/benutzerhandbuch.md`; `harness/README.md` unberührt (`git diff --stat b80f45cf..994b887a` zeigt nur die Plan-Datei und das Handbuch) |
| 8 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt weiterhin `*(zu tragen bei Closure)*` |
| 9 | Reconciliation-Register — entfällt | **korrekt offen/entfällt** | Greenfield, wie in allen Slices dieses Repos |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht, nicht Sache der Fixrunde |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | alle vier Zeilen tragen weiterhin `*(bei Closure)*` |
| 12 | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | hängen an der Closure von `welle-transformationen`, wie im Plan benannt |

`[x]` sind acht Zeilen (Nr. 1–7 plus Reconciliation Nr. 9), `[ ]` vier (Nr. 8, 10, 11, 12) —
deckungsgleich mit den zwölf Aufzählungspunkten in §2 des Plans (`grep -c '^\- \[.\]'` = 12). Kein
`[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus (Closure/Planner) belegt wäre.

## 3. Fixrunde: Findings des Reviews selbst nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Zitierter PostgreSQL-Fehlertext | `docs/reviews/review-slice-transformationen-antragsweg-schema.md:101-102` gelesen: „`{"kind":"x"}`::jsonb`` „function cdc.set_transformation(unknown, unknown, unknown, unknown, jsonb) does not exist“ — Wortlaut exakt identisch mit `docs/user/benutzerhandbuch.md:337` nach der Fixrunde. `git grep -n "text, text, text, text, jsonb" -- '*.md'` findet die alte, falsche Zeichenkette nur noch im Review-Report selbst (als zitierter Fehlbefund), nicht mehr im Handbuch; `git grep -n "unknown, unknown, unknown, unknown, jsonb"` findet genau zwei Treffer: Review-Report (Zitat des korrekten Werts) und Handbuch (die korrigierte Zeile) | **behoben**, Wortlaut-Identität bestätigt |
| F-2 (MEDIUM) Unvollständige Übernahme „Zeilen, die kein Antrag sind“ | `internal/adapters/driven/postgresstorage/administrationrequest_test.go:1439-1442` gelesen: vier Fälle — „Schemaname ist leer“, „Tabellenname ist leer“, „Spaltenname ist leer“ (`exclude_column`), „Spaltenname ist leer“ (`include_column`). Handbuch-Absatz (`docs/user/benutzerhandbuch.md:406-409`) nennt jetzt „leeres Schema, leerer Tabellenname, leere Spalte bei `exclude_column`/`include_column`“ — deckt alle vier Testfälle in der vom Plan verlangten Drei-Beispiel-Form (die beiden Spalten-Fälle als ein Beispiel gebündelt, wie der Plan es selbst formuliert: „`exclude_column`/`include_column` mit leerer Spalte“) | **behoben**, deckungsgleich mit Plan-Wortlaut und Testfällen |

Keine neuen Findings durch die Fixrunde selbst eingeführt: die einzige inhaltliche Änderung neben
den zwei Korrekturen ist die Versionshebung 1.67→1.68 mit einer neuen Änderungshistorie-Zeile, die
in Form und Inhalt (Kennungen, Ursprungs-Angabe „gemessen“/„belegt durch“) dem Muster der
elf vorangehenden Zeilen (1.57–1.67) folgt — kein Chronik-Verstoß gegen §3.7 (das ist die
Änderungshistorie-Tabelle selbst, deren etablierte Form genau diese Rückschau trägt, keine Prosa
außerhalb davon).

## 4. Zitat- und Beleg-Prüfung (eigenständig, nicht dem Fixrunden-Commit geglaubt)

- `docs/user/benutzerhandbuch.md:337`: `` `function cdc.set_transformation(unknown, unknown, unknown, unknown, jsonb) does not exist` `` — Zeichen-für-Zeichen-Abgleich mit
  `docs/reviews/review-slice-transformationen-antragsweg-schema.md:101-102`: identisch.
- `docs/user/benutzerhandbuch.md:406-409`: „verwirft (leeres Schema, leerer Tabellenname, leere
  Spalte bei `exclude_column`/`include_column`), endet `failed` …“ — deckt die im Plan (§2 Punkt 2,
  Übergabe aus `antragsqueue-lesefehler-failed`) explizit verlangten drei Beispiele vollständig.
- Änderungshistorie-Zeile 1.68 (`docs/user/benutzerhandbuch.md`, letzte Zeile der Tabelle): nennt
  `LH-FA-CFG-007`, `ADR-0112`, `ADR-0125`, den Slice-Namen und exakt die zwei behobenen Findings —
  keine über den Fixrunden-Diff hinausgehende Behauptung.

## 5. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Alle acht Zeilen des committeten Suchlauf-Blocks selbst über `make suchlauf-nachmessen`
ausgeführt (§1) — Ergebnis: **8/8 stimmen**, inklusive der von der Fixrunde geänderten Zeile
(„Aufzählungen der Antragsarten im Handbuch“, Diff-Stand jetzt 11 statt 9 Treffer). Die
Fixrunde hat damit selbst korrekt erkannt und nachgezogen, dass sie durch ihre eigene
Textänderung (den ergänzten Absatz und die neue Änderungshistorie-Zeile, beide mit
`exclude_column`/`include_column`) die gezählte Eigenschaft aus §3.13 erneut bewegt hat — ein
Suchlauf-Nachzug innerhalb der Fixrunde selbst, korrekt ausgeführt.

**Eigene Nicht-Fund-Suche** (unabhängig vom Plan-Feld, ob die Fixrunde eine *weitere*, im Plan
nicht geführte Eigenschaft bewegt hat, ohne sie nachzuziehen): Die Fixrunde ändert genau drei
Textstellen (Versionskopf, Fehlertext-Zeile, Absatz „Zeilen, die kein Antrag sind“, Änderungshistorie).
Geprüft, ob eine dieser vier Stellen eine weitere im Repo beschriebene Eigenschaft bewegt:

- Der korrigierte Fehlertext (`unknown` statt `text`) erscheint sonst nirgends im Repo als
  falsche Zitat-Kette (`git grep -rn "text, text, text, text, jsonb"` nach der Fixrunde: 0 Treffer
  im Baum außerhalb des Review-Reports, der den *falschen* Wortlaut als *Zitat des Fundes*
  historisch korrekt weiterträgt — kein Nachzugsbedarf dort, ein Review-Report ist Lauf-Beleg und
  immutable im Sinne der Auditierbarkeit, nicht im Sinne von §3.5).
- Der ergänzte dritte Beispiel-Fall (`exclude_column`/`include_column` mit leerer Spalte) ist kein
  neues Konzept, sondern die vom Plan bereits verlangte, jetzt nachgetragene Vervollständigung —
  keine weitere Beschreibungsstelle im Repo zählt diese Beispiele um eine Eigenschaft, die über das
  bereits im Plan geführte Suchlauf-Feld hinausginge (`grep -rln "leeres Schema, leerer
  Tabellenname" docs/` liefert nur die eine Handbuch-Stelle).

Kein weiterer, vom Plan-Feld nicht erfasster Nachzugsbedarf gefunden.

## 6. Entscheidungs-Konformität (Stichprobe der berührten ADRs)

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) Folgepflicht 2 | `rule_spec` ist `json`, kein `jsonb`; Fehlertext bei `::jsonb`-Aufruf dokumentieren | Handbuch-Absatz korrekt mit dem realen PostgreSQL-Fehlertext (`unknown`-Parametertyp), Cast `::json` als Abhilfe genannt | konform |
| [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 8 | keine Änderung am Nachrichtenschema durch diesen Slice | `git diff --stat b80f45cf..994b887a` berührt kein `proto/`, kein `gen/`; `make generated-sync` in `make gates` grün (§1) | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B | jede Tatsachenbehauptung trägt einen Beleg-Anker oder ist als Zusage/übernommen gekennzeichnet | Fehlertext-Zeile trägt „(gemessen im Review-Report …)“ als Anker, der jetzt tatsächlich zutrifft (§4) | konform (vorher: F-1-Verstoß, jetzt behoben) |
| [`AGENTS.md`](../../AGENTS.md) §3.13 | bewegte Eigenschaft im Suchlauf-Feld nachgezogen | §5 — die Fixrunde hat die betroffene Zeile korrekt von 9 auf 11 aktualisiert | konform |

## 7. Plan-vs-Code-Diff

`git diff --stat c621bf36..994b887a`: zwei Dateien — `docs/plan/planning/in-progress/slice-transformationen-betriebsdoku.md`
(+12/−7: DoD-Checkbox-Nachzug, Suchlauf-Feld-Nachzug) und `docs/user/benutzerhandbuch.md` (+8/−3:
Versionskopf, Fehlertext-Zeile, Absatz-Ergänzung, Änderungshistorie-Zeile). Nichts Ungeplantes:
beide Änderungen sind exakt die vom Review verlangten Korrekturen plus ihr notwendiger Nachzug im
Plan (Suchlauf-Feld, DoD-Checkbox). Kein SDK-Code, kein Produktivcode, keine ADR, keine Spec im
Diff — deckt sich mit dem in AGENTS.md §3.10 nicht einschlägigen Charakter dieses reinen
Doku-Slice (kein Workflow-Zug).

## 8. Harte Regeln

- **§3.1** — keine Host-Toolchain; alle Läufe über `make`.
- **§3.5** — kein Eingriff in eine `Accepted`-ADR; `make doc-immutable RANGE=b80f45cf..994b887a` EXIT=0 (§1).
- **§3.7** — Änderungshistorie-Zeile 1.68 folgt der etablierten Chronik-Form der Tabelle selbst (keine Prosa-Chronik außerhalb davon); geprüft, kein Verstoß.
- **§3.9** — Exit-Codes ungefiltert gesichert (`; ec=$?` je Lauf), nie durch eine Pipe.
- **§3.11** — kein host-lokaler absoluter Pfad im Diff; `make docs-check` (hostpaths-Modul in `make gates`) grün.
- **§3.12** — F-1 (falscher Zitat-Wortlaut, Instanz B) behoben und nachgemessen (§4); Änderungshistorie-Zeile 1.68 trägt korrekte Ursprungs-Kennzeichnung.
- **§3.13** — Suchlauf-Feld an allen 8 Zeilen nachgefahren (§5), die von der Fixrunde selbst bewegte Eigenschaft korrekt nachgezogen.
- **Commit-Traceability** — 3 Commits (`97de226c`, `c621bf36`, `994b887a`), jede Message nennt `LH-FA-CFG-007`/`ADR-0112` u. a.; `make commit-traceability` (Teil von `make gates`) grün.

## 9. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (8) | **8 von 8 erfüllt**, je mit eigenem Beleg |
| DoD §2 — übrige `[ ]`-Zeilen (4) | **4 von 4 korrekt offen** (Closure-/Planner-Sequenz) |
| Review-Findings F-1 (HIGH), F-2 (MEDIUM) | **beide behoben**, je selbst gegen ihre Quelle nachgemessen (§3, §4) — kein offenes HIGH/MEDIUM |
| §3.13-Suchlauf | **8/8 bestätigt**, inkl. der von der Fixrunde selbst nachgezogenen Zeile; kein weiterer Nachzugsbedarf gefunden (§5) |
| Plan-vs-Code-Diff | **deckungsgleich**; nichts Ungeplantes |
| Entscheidungen | [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md), [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) **konform** |
| SDK-Beleg (DoD 1) | unberührt von der Fixrunde, `sdks/` nicht im Diff (§2 Nr. 1) |
| Harte Regeln | **erfüllt** (§8) |
| Gates | `make gates` **EXIT=0**, `make docs-check` **EXIT=0**, `make doc-commits`/`make doc-immutable` **EXIT=0** |

## Verdikt

**Bestätigt — bereit für Planner-Closure.** Beide Review-Findings (F-1 HIGH, F-2 MEDIUM) sind in
der Fixrunde korrekt behoben; ich habe beide Korrekturen selbst gegen ihre genannten Quellen
(Review-Report-Zeile 101-102, Testdatei-Zeilen 1439-1442) nachgelesen und Zeichen-für-Zeichen-Identität
festgestellt — keine Behauptung ungeprüft übernommen. Der §3.13-Suchlauf ist an allen acht Zeilen
selbst nachgefahren (`make suchlauf-nachmessen`, 8/8), inklusive der von der Fixrunde selbst
korrekt nachgezogenen vierten Zeile; eine eigene Nicht-Fund-Suche über die Fixrunden-Änderungen
hinaus findet keine weitere übersehene bewegte Eigenschaft. `make gates`, `make docs-check`, `make
doc-commits` und `make doc-immutable` laufen am Stand `994b887a` mit Exit 0, jeweils ungepiped
geprüft. Der SDK-Beleg aus dem ursprünglichen Implementer-Lauf (137/137, vom Reviewer bereits
Fundstelle-für-Fundstelle gelesen) ist von der Fixrunde unberührt geblieben (`sdks/` nicht im
Fixrunden-Diff). Die DoD-Tabelle ist konsistent: acht Implementer-/Fixrunden-eigene Punkte `[x]`
mit Beleg, vier Closure-/Planner-Punkte korrekt `[ ]`. Kein neues HIGH/MEDIUM, keine weitere
Fixrunde nötig.

**Welle-Schließbarkeit:** Mit der Closure dieses Slice liegen alle zehn Slices von
`welle-transformationen` in `done/` (neun bereits verifiziert vor diesem Lauf, dieser Slice ist der
zehnte und letzte laut `docs/plan/planning/welle-transformationen.md` §4-Tabelle). Das ist genau
eines der in `welle-transformationen.md` §3 genannten Closure-Kriterien („Alle zehn Slices in
`done/`“) — die **übrigen** Kriterien dieses Abschnitts (realer grüner `make test-integration`-Lauf,
`make test-store`/`make test-replication`, Alt-Tag-Lauf, Fitness-Function-Tests, `make doc-trace`,
Lese-Schritt des Beobachtungs-Registers, Closure-Notiz der Welle) sind **nicht** Gegenstand dieser
Slice-Verifikation und von mir nicht geprüft — das bleibt Sache der Planner-geführten
Welle-Closure-Prozedur selbst, mehrere davon sind laut dem Wellen-Plan bereits mit
Architect-Verdikten oder vorangegangenen Slice-Closures erfüllt (§3 dort, „Erfüllt:“-Vermerke).
Der Planner kann nach der Closure dieses Slice direkt die Welle-Closure-Prüfung nach
`welle-transformationen.md` §3 einleiten.

**Übergabe:** Verifier → Planner. Offen bleibt für den Planner: reiner `git mv` dieses Slice nach
`done/` nach Ergänzung von Closure-Notiz (§7 des Slice-Plans), Beobachtungs-Register-Fortschreibung
und den vier §6-Risiko-Ausgängen — danach die Welle-Closure-Prüfung nach `welle-transformationen.md`
§3 mit ihren eigenen, breiteren Kriterien.
