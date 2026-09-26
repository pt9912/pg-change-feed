# Review-Report: slice-harness-fmt-check — 2026-09-26

**Review-Art:** Code — der Diff liefert das Werkzeug `make fmt-check` (Aufrufer, Tabellentest, Vertrag,
zwei Make-Ziele, zwei README-Zeilen), den Format-Commit der sechs Bestandsdateien und den Nachzug der
Träger (Schritt 18 des Implementer-Ablaufs, LOW-Zeile im Reviewer-Skill); geprüft gegen Plan, die
Architect-Entscheidung „Werkzeug ohne Gate“ im Register-Eintrag `BEO-PGC/formatierungs-drift-ohne-gate`,
[ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.1/§3.6/§3.7/§3.9/
§3.12/§3.13 und die Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-harness-fmt-check` (ohne Welle, Harness-Querschnitt), Diff-Range
`17cb4eb3..95ff5ab2` (7 Commits, 14 Dateien, +391/−31; Baum sauber, nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ in der Form des Diff-Stands
`95ff5ab2` (mit der geänderten LOW-Zeile, die dieser Lauf zugleich prüft).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

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

- Slice-Plan `slice-harness-fmt-check` (§1 Ziel und Abgrenzung, §2 Liefer-Punkte als Bezug, §3 Plan,
  Nachzug-Tabelle, Ist-Zustand am Start, Suchlauf-Feld und Träger-Tabelle, §4 Rückführungen, §6 Risiken);
  Register-Eintrag `BEO-PGC/formatierungs-drift-ohne-gate` (`state.md`, Ausgang „Werkzeug ohne Gate“)
- [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen in Trägern)
- `AGENTS.md` (Hard Rules §3.1 mit der Klasse „Host-Werkzeug ohne Installation“ und dem in-place-Verbot,
  §3.2, §3.6, §3.7, §3.9, §3.12, §3.13, §4), `harness/conventions.md` (`MR-000`/`MR-001`)
- Baseline `v6.9.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-code-kommentare-kennungen.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; Scratchpad der Sitzung):

- **Läufe am Stand `95ff5ab2`** (einer zugleich; `free -m` vorher 13,4 GB verfügbar; dangling Volumes
  `docker volume ls -qf dangling=true | wc -l` 34 vor und 34 nach allen Läufen; kein `prune`;
  `git status --short` danach leer):
  `make test-fmt-check` Exit 0, gedruckt „run-fmt-check-tests: alle Fälle bestanden“;
  `make fmt-check` Exit 0, gedruckt „fmt-check: 254 Go-Dateien geprüft, alle formatiert“;
  `make test` (Race) Exit 0; `make a-check` Exit 0, „gesamt: 0 Befund(e)“;
  `make coverage-gate` Exit 0, „coverage-gate: OK — Coverage 85.00% erfüllt Schwelle 80%“;
  `make gates` (einmal) Exit 0, in der Ausgabe kein Vorkommen von `fmt-check` (das Ziel ist in
  keinem Gate-Bündel; `git diff` von `Makefile` trägt keine Zeile zu `GATE_CHECKS`);
  `make commit-traceability RANGE=17cb4eb3..95ff5ab2` Exit 0, „7 Commit(s) … Betreffs ohne Struktur-ID“;
  `make kommentar-kennungen DIFF=17cb4eb3 COUNT=1` Exit 0, Ausgabe **0**.
- **Dateizählung:** `git ls-files '*.go' | wc -l` **254** (am Stand `17cb4eb3` und `95ff5ab2`, `git ls-tree`
  ebenso; am Stand `7b70b34a` 252 — Differenz zwei neue Dateien: `tools/harness/kommentar-kennungen/main.go` und
  `main_test.go`; `find . -name '*.go'` außerhalb `.git` ebenfalls 254; kein tracked `.go` unter `testdata/`,
  `vendor/` oder mit `_`-Präfix, keine Symlink-Datei mit `.go`-Endung im Baum (`git ls-files -s`, Modus
  120000: 0). Die Zählung des Werkzeugs (254) stimmt damit.
- **Randfälle des Aufrufers** (Wegwerf-Verzeichnisse im Scratchpad, echter Docker-Lauf): Pfad mit Komma,
  `testdata/`, `_`-Verzeichnis und ein Verzeichnis `-dash` (absolut, relativ als `-dash` und `./-dash`)
  Exit 1 mit genannter Datei — `gofmt` und Zählung nehmen dieselben Dateien, Optionsinjektion über den
  Verzeichnisnamen tritt nicht ein (`realpath --` macht den Pfad absolut, `-v "$dir"` gequotet); Pfad mit
  Doppelpunkt: `docker: invalid spec … too many colons`, Exit 2 (F-4); `.go`-Symlink auf eine Datei: `gofmt -l`
  nennt ihn, die Zählung (`-type f`) nicht (F-4). Das `[ -d ]`-Vorprüfung fängt ein fehlendes Verzeichnis
  vor `docker run` ab (Test „kein Verzeichnis“, Mutation M9 rot).
- **Eingabeseiten-Mutationen am Werkzeug** (je eine Kopie im Scratchpad per `sed` ohne `-i` nach stdout
  erzeugt, `TOOL=<Kopie> make test-fmt-check`, Exit ungepiped; das Original blieb unverändert, `git status`
  sauber):

  | Nr. | Mutation an `fmt-check.sh` | Ergebnis (rote Fälle) |
  |---|---|---|
  | M1 | Auswertung der Ausgabe entfernt (nur Exit-Code des Formatierers) | rot: unformatierte Datei · gemischt · Unterverzeichnis · Eingabe unverändert · Pfad mit Leerzeichen |
  | M2 | `gofmt -l ./*.go` statt `.` | rot: dieselben fünf (der Pfad trägt `./`; die reine Baumtiefe färbt M2b) |
  | M2b | `gofmt -l *.go` statt `.` | rot: nur Unterverzeichnis |
  | M3 | Zählung der Go-Dateien entfernt | rot: leeres Verzeichnis · Verzeichnis ohne Go-Datei |
  | M4 | Exit des Formatierers ignoriert | rot: Syntaxfehler |
  | M5 | Docker-Exit 125 durchgereicht | rot: Docker-Fehler |
  | M6 | `-w` ohne `:ro` | rot: Eingabe wurde umgeschrieben |
  | M7 | `"$dir"` ohne Anführungszeichen | rot: Pfad mit Leerzeichen |
  | M8 | Argumentzahl-Prüfung entfernt | rot: zu viele Argumente |
  | M9 | Verzeichnis-Prüfung entfernt | rot: kein Verzeichnis |
  | M10 | Prüfung der Variablen `TOOLCHAIN_IMAGE` entfernt | rot: `TOOLCHAIN_IMAGE` fehlt |
  | M12 | Zählung nimmt Dateien mit Punkt am Namensanfang | rot: Verzeichnis ohne Go-Datei |
  | M13c | Ausgabe trägt zusätzlich jede Go-Datei | rot: formatierte Datei · gemischt: formatierte genannt · Syntaxfehler |
  | M15 | Exit 1 bei Abweichung entfernt | rot: dieselben fünf wie M1 |
  | M11 | `:ro` allein entfernt | **grün** (F-3) |
  | M14 | `--network none` entfernt | **grün** (F-3) |
  | M16 | Prüfung des Exit von `realpath` entfernt | grün (unerreichbar: das Verzeichnis existiert nach der Vorprüfung) |

  Alle zwölf Zeilen der Mutationstabelle im Vertrag (§Test) sind damit an genau den dort genannten
  Fällen rot gesehen; die drei grünen Läufe stehen nicht in der Tabelle.
- **Format-Commit `32637052`:** `git diff 17cb4eb3..95ff5ab2 -- internal test` gelesen: vier Dateien
  ausschließlich Whitespace/Ausrichtung (`transformation_test.go`, `seam_test.go`, `service_test.go`) bzw.
  Umbruch von vier einzeiligen Methoden in je drei Zeilen (`log_test.go`, gleiche Anweisungen); zwei
  Kommentarblöcke geändert (`queries.go`, `integration_test.go`), keine Nicht-Kommentar-Zeile mit
  Inhaltsänderung. Am Start (`32637052^`, Kopie der sechs Dateien im Container): `gofmt -d` druckt **6**
  Hunks in **86** Zeilen (`@@ -76`, `-399`, `-257`, `-14,10 +14,18`, `-51`, `-1236,7`), wie der Plan sagt;
  `git diff -U0` des Commits gegen `integration_test.go` trägt genau `@@ -1239 +1239 @@`, `docs/user/
  e2e-abdeckung.md` ist im Range unverändert (`git diff` leer) — die Zeilen-Lokatoren verschieben sich nicht.
  `make test` (Race) grün, `make fmt-check` Exit 0.
- **Die drei umformulierten Kommentare gegen den Code:** `queries.go`: die Abfrage lautet `c.change_id > $2`
  mit `$2` als Text-Kennung; `''` ist im SQL der leere Text, jede Kennung ist größer — „leerer Text: ab dem
  Anfang“ ist wahr. `RetentionPolicy` steht in `internal/domain/model/retention.go` (Typ Zeile 14,
  `AllowsDeletion` Zeile 35) und trägt die Freigabe je Change; der Verweis ersetzt `ADR-0014` (der Block
  trägt danach eine Kennung, `ADR-0124`) und ist die Kopplungsform nach `AGENTS.md` §3.7 — bestätigt.
  `integration_test.go`: der Kommentar nennt `…003`, `/004`, `...005`; die Zeichenliste im Code ist
  `{"…", "...", "/"}` — das Backtick-Paar war ein Tippfehler des Kommentars, die neue Form stimmt.
- **Suchlauf-Feld:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, gedruckt „suchlauf-nachmessen: 8
  Zeilen stimmen“. Von Hand mit `git grep` nachgefahren (ohne das Werkzeug, Ist gleich Soll): Zeile 1
  (`7b70b34a`, `gofmt`) 4; Zeile 2 (`Bestandsdateien außerhalb des Diffs`) 1; Zeile 3 (`fmt-check`, Parent) 0;
  Zeile 4 (`17cb4eb3`, fremde Träger, ohne die Plan-Datei) 2 (`slice-antragsqueue-lesefehler-failed` Zeile 243,
  `slice-transformationen-e2e-wirkung` Zeile 166, je ein Satz); Zeile 5 (`diff`, `gofmt`) 3 (drei Zeilen in
  Schritt 18); Zeile 6 (`diff`, Bestands-Satz) 0; Zeile 7 (`diff`, `fmt-check`) 37. Die zwei gemeldeten
  fremden Träger sprechen von den „sechs Bestandsdateien“ — die Meldung trifft zu (F-5 nennt einen dritten,
  vom Muster nicht erfassten Träger).
- **Träger gelesen:** Schritt 18 „Format“ in `.claude/commands/implement-slice.md` (Aufruf `make fmt-check`,
  „Bestandsdateien außerhalb des Diffs“ entfallen — richtig, denn der Lauf misst den ganzen Baum und der
  Bestand ist formatiert; die Regel „nach `gofmt -d` korrigieren, nie mit `sed`“ bleibt; Nummern 19–21
  unverändert); LOW-Zeile im Reviewer-Skill; Vertrag `harness/sensors/fmt-check.md`. Zeilen-Lokatoren auf
  `reviewer.md`/`implement-slice.md` außerhalb der Records: keine (`git grep -n -E
  'reviewer\.md:[0-9]|implement-slice\.md:[0-9]' -- . ':!docs/reviews'` leer). Die Beleg-Befehle des Vertrags
  (`gofmt -d` im Toolchain-Image, mit `$PWD`-Mount) habe ich gefahren: Exit 0, leere Ausgabe für eine
  formatierte Datei; auf den Startstand der sechs Dateien Ausgabe wie oben.
- **Hard Rules:** `git diff 17cb4eb3..95ff5ab2 | grep -nE 'sed -i|perl -pi|awk -i'` ohne Treffer;
  `.a-check.yml` unberührt, `make a-check` Exit 0; `tools/` nicht in der gemessenen Coverage-Fläche;
  Commit-Struktur: die zwei Lifecycle-Moves `2d053bba` (open → next) und `47926f65` (next → in-progress) sind reine
  Renames (`git show -M --stat`: 0 Zeilen), der Inhalt steht in eigenen Commits, der Format-Commit ist ein eigener Commit
  ohne Werkzeug-Änderung.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Ein Skript-Kommentar begründet die Regel „leer ist nicht bestanden“ mit der verworfenen Alternative im Präteritum: „ein falsch gemounteter Pfad meldete sonst grün“. Das ist die Form „Ohne X wäre …“ (`AGENTS.md` §3.7: Konjunktiv über die verworfene Alternative statt Indikativ über den Zustand); der Skopus des Reviewer-Skills nennt Skripte in `tools/harness/*.sh` ausdrücklich. Derselbe Halbsatz steht im Vertrag (Abschnitt „Exit-Codes“). | `AGENTS.md` §3.7; Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen“ | `tools/harness/fmt-check.sh:36-39`; `harness/sensors/fmt-check.md:59-61` | ja — `grep -nE 'sonst grün' tools/harness/fmt-check.sh harness/sensors/fmt-check.md` (zwei Treffer) | Kommentar beschreibt die verworfene Alternative |
| F-2 | LOW | Der Nachzug „Ist-Zustand am Start“ steht als Absatz mitten in der Datei-Tabelle von §3: die Kopfzeile und die ersten acht Zeilen der Tabelle stehen davor, die sechs Zeilen `queries.go` … `integration_test.go` danach ohne Kopf und ohne Trennzeile. Sie rendern als Fließtext, die Änderungs-Tabelle des Plans ist zerrissen. | Maintainability; Nachzug widerspricht dem Nachbarn im selben Träger (Reviewer-Skill; hier Form, kein Widerspruch der Aussagen) | `docs/plan/planning/in-progress/slice-harness-fmt-check.md:163-176` | ja — Lesen des Kontexts (`git diff 47926f65..95ff5ab2 -- <Plan-Datei>`); `make docs-check` meldet es nicht | Einfügung zerreißt eine Tabelle |
| F-3 | LOW | Zwei Isolations-Zusagen des Vertrags sind ohne Bindung: `--network none` (Vertrag und README „netzlos“) und `:ro` allein (Vertrag „das Verzeichnis lesend gemountet“, „schreibt nichts“) färben keinen Fall rot — M14 und M11 laufen grün; `:ro` ist nur zusammen mit `-w` gebunden (M6). Der Plan sagt, der Vertrag nenne „je Zusage die Mutation der Eingabeseite“; für diese beiden steht keine Zeile in der Tabelle. Abgrenzung: hier **fehlt** die Abdeckung, sie steht nicht ungebunden — deshalb LOW, kein HIGH; die Zusagen betreffen die Isolation des Aufrufs, nicht das Verdikt des Werkzeugs. | Maintainability; fehlende Negativtests bei neuem öffentlichem Vertrag (Reviewer-Skill), Plan §3 (Nachzug „je Zusage“) | `tools/harness/fmt-check.sh:59`; `tools/harness/run-fmt-check-tests.sh`; `harness/sensors/fmt-check.md` §Test (Tabelle) | ja — M11 und M14 an einer Kopie im Scratchpad, `make test-fmt-check` Exit 0 | Isolations-Flags des Aufrufs ohne Testbindung |
| F-4 | INFO | Zwei Randfälle liegen außerhalb dessen, was der Vertrag zusagt („dieselben Dateien, die `gofmt` liest“): ein `.go`-Symlink wird von `gofmt -l` genannt, von der Zählung (`find -type f`) nicht — ein Baum nur aus Symlinks endet mit Exit 2 „keine Go-Datei“, ein gemischter zählt zu wenig; ein Verzeichnis mit `:` im Pfad scheitert an `-v <Pfad>:/src:ro` mit Exit 2 (die Meldung nennt „Exit 125“). Beide enden geschlossen (Exit 2, nie fälschlich 0); der Repo-Baum trägt keinen Symlink mit `.go`-Endung, seine Wurzel keinen Doppelpunkt. | `harness/sensors/fmt-check.md` §Aufruf („dieselben, die `gofmt` liest“) | `tools/harness/fmt-check.sh:41`, `:59` | ja — Wegwerf-Verzeichnisse `co:lon` und `sym/l.go` (gefahren) | Zählregel weicht von `gofmt` bei Symlinks ab |
| F-5 | INFO | Das Suchmuster der vierten Zeile des Suchlauf-Felds (`sechs Bestandsdateien\|Bestandsdateien, die`) ist zeilenweise; `docs/plan/planning/welle-transformationen.md:304-305` und `:306` sprechen von den sechs Bestandsdateien über einen Zeilenumbruch („Formatierung der sechs“ / „Bestandsdateien)“, „eine der sechs;“ / „Bestandsdateien“) und fehlen im Ergebnis und in der Träger-Tabelle. Beide Sätze bleiben nach `done/` wahr (sie beschreiben den Umfang des Slice, keine offene Kante); die Aussage „zwei fremde Träger“ ist damit die Zahl der Treffer des Musters, nicht der Träger. | `AGENTS.md` §3.13 (Suchform: Zählwort und Beschreibung, Vollständigkeit ist Lese-Handlung) | `docs/plan/planning/in-progress/slice-harness-fmt-check.md:199-213`; `docs/plan/planning/welle-transformationen.md:304-306` | ja — `git grep -n -E 'Bestandsdatei' -- docs/plan/planning/welle-transformationen.md` | Suchmuster blind gegen Zeilenumbruch |
| F-6 | INFO | Der Test-Fall „Docker-Fehler“ benutzt ein Image auf `fmt-check-test.invalid`; der Kopf des Testskripts nennt den Test „netzlos“ (`--network none` gilt für den Container, nicht für den Image-Zugriff des Docker-Daemons — dessen Namensauflösung bleibt hergeleitet, nicht gemessen). Der Fall ist grün und rot färbbar (M5). | Maintainability | `tools/harness/run-fmt-check-tests.sh:132-135` | nein | „netzlos“ trägt für den Docker-Fehler-Fall nur mittelbar |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/fmt-check.sh` (Exit-Codes 0/1/2, Auswertung der Ausgabe statt des Exit-Codes, Zählung, Mount, Quoting) | geprüft, ohne Befund über F-1, F-3, F-4 hinaus: `gofmt -l` endet bei Abweichung mit Exit 0 und das Werkzeug wertet die Ausgabe (M1 rot); Docker-Exit 125 und jeder Exit ≠ 0/1 werden zu 2 (M5 rot); Syntaxfehler Exit 2 (M4 rot); Zählung 254 stimmt mit `git ls-files`; `gofmt` und Zählung nehmen `testdata`, `_`-Verzeichnisse, Punkt-Dateien einheitlich; `--network none`, Mount `:ro` mit absolutem `realpath`-Pfad, `$dir` gequotet (M7 rot), keine Options- oder Argumentinjektion über Verzeichnisnamen (`-dash`, Komma, Leerzeichen gefahren); `TOOLCHAIN_IMAGE` ist ein gepinnter Digest (`Makefile:176`) und wird vor dem Lauf geprüft (M10 rot); das Host-`git` wird nur für die Repo-Wurzel ohne Argument gebraucht und im Vertrag genannt |
| `tools/harness/run-fmt-check-tests.sh` | geprüft, ohne Befund über F-3, F-6 hinaus: zehn Fälle mit unterscheidbaren Erwartungen (Exit und Muster, `refuse` für die nicht genannte Datei), Eingabe byte-gleich per `cmp`, `trap` räumt das Temp-Verzeichnis, Prüfling per `TOOL` übersteuerbar, keine Kennungen in den Kommentaren, keine Chronik-Sprache; 12 von 12 Tabellenzeilen rot gesehen |
| `harness/sensors/fmt-check.md` | geprüft, ohne Befund über F-1, F-3 hinaus: Vertrag, Exit-Tabelle, Grenze (Formatierung, nicht Semantik; `gofmt` formt Doc-Kommentare um; nur Go; „Grün ist kein Beleg für den Diff“), Trigger der Gate-Aufnahme als Kenntnis mit „trotz gelaufenem Schritt 18“, Host-Werkzeuge nach der Klasse in `AGENTS.md` §3.1; die Mutationstabelle nennt je Zusage Mutation und roten Fall und stimmt mit meinen Läufen überein (Menge der Erprobung: die Fälle des Tabellentests) |
| `Makefile` (zwei Ziele), `harness/README.md` (zwei Zeilen) | geprüft, ohne Befund: `.PHONY`, Hilfe-Zeilen, kein Eintrag in `GATE_CHECKS` (`make gates` unverändert, Ausgabe ohne `fmt-check`), README-Zeile trägt „kein Gate“, Exit-Semantik samt „über `make` Exit 2“, Host-Werkzeuge und den Trigger-Verweis |
| Format-Commit `32637052` und die sechs Go-Dateien | geprüft, ohne Befund: nur Whitespace/Umbruch/Kommentar, keine Verhaltensänderung, `gofmt -d` am Start 6 Hunks/86 Zeilen wie im Plan, Zeilen-Lokatoren von `docs/user/e2e-abdeckung.md` unbewegt, drei Kommentartexte wahr gegen den Code, die geänderten Blöcke tragen höchstens eine Kennung (`make kommentar-kennungen DIFF=17cb4eb3` 0) |
| `.claude/commands/implement-slice.md` Schritt 18, `.harness/skills/reviewer.md` LOW-Zeile | geprüft, ohne Befund: Aufruf `make fmt-check`, Rang-Zeiger auf den Vertrag statt Docker-Befehl-Wiederholung, Grenze „erste, nicht tragende Linie“ bleibt, Bestands-Satz entfällt zu Recht, Schrittnummern unverändert; die LOW-Zeile nennt `make fmt-check` und den Vertrag; keine Zeilen-Lokatoren außerhalb der Records betroffen |
| Plan `slice-harness-fmt-check` (Nachzug, Ist-Zustand, Suchlauf-Feld, Träger-Meldungen) | geprüft, ohne Befund über F-2, F-5 hinaus: Drift 252 → 254 mit richtiger Ursache und Herkunft (Stand, Befehl) gemessen und stimmend, die Träger-Tabelle nennt Gefundenes und Nichtgefundenes je Träger, beide Stände; die Meldung der zwei fremden Träger trifft zu und trägt Frist (Planner-Closure); DoD-Haken stimmen mit dem Diff, die offenen Punkte (Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) bleiben `[ ]` |
| Hard Rules, Commit-Struktur, Handbuch | geprüft, ohne Befund: Docker-only (Host-Werkzeuge nur `bash`, `git`, `realpath`, `docker`), kein `sed -i`/`perl -pi`/`awk -i` im Diff, keine Suppression, `make a-check`/`make coverage-gate`/`make gates` grün, `tools/` nicht in der Coverage-Fläche, sieben Betreffs mit `ADR-0083` ohne `SPEC-`/`ARC-`, Format-Commit rein; das Benutzerhandbuch bleibt unberührt und muss es (keine Betreiber-Oberfläche) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kommentar beschreibt die verworfene Alternative · Einfügung zerreißt eine
Tabelle · Isolations-Flags des Aufrufs ohne Testbindung · Zählregel weicht von `gofmt` bei Symlinks ab ·
Suchmuster blind gegen Zeilenumbruch · „netzlos“ trägt für den Docker-Fehler-Fall nur mittelbar

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (HIGH nach der Liste des Reviewer-Skills): ein Kommentar im neuen
Skript und der Vertrag begründen eine Regel mit der verworfenen Alternative. Die Korrektur betrifft zwei
Halbsätze; das Werkzeug selbst ist in Ordnung: Exit-Codes, Auswertung der Ausgabe, Zählung, Mount und
Quoting halten meinen Mutationen stand (14 von 17 rot, die drei grünen sind F-3 und ein unerreichbarer
Zweig), der Format-Commit ist verhaltensgleich und die drei geänderten Kommentare sind wahr gegen den Code,
alle Läufe (`make test-fmt-check`, `make fmt-check`, `make test`, `make a-check`, `make coverage-gate`,
`make gates`) grün. F-2 und F-3 gehen in derselben Fixrunde mit; F-4 bis F-6 sind Hinweise ohne erwartete
Aktion, F-5 zusätzlich Meldung an den Planner (Träger `welle-transformationen`, Frist: Closure dieses Slice).

**Übergabe:** Findings gehen an den Implementer (Fixrunde); die **Finding-Klassen** gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler. Die DoD-Zeile „Review durchgeführt“ im Slice-Plan bleibt offen
(eine Fixrunde folgt; sie wird bei Schritt 21 des Implementer-Ablaufs nachgezogen). Dieser Report selbst ist
ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über Läufe
hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
