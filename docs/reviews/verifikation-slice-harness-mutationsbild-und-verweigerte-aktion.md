# Verifikations-Report: slice-harness-mutationsbild-und-verweigerte-aktion — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität (Nutzer-
Entscheidungen „ein eigenes make-Ziel mit eigenem Tag" und „ja" zur Regel „verweigerter Aufruf wird
gemeldet und vor einem Ersatzweg zurückgefragt", ohne Ausnahme) + Plan-vs-Code-Diff + Gates.
Review-Artefakt:
[`review-slice-harness-mutationsbild-und-verweigerte-aktion.md`](review-slice-harness-mutationsbild-und-verweigerte-aktion.md)
(Commit `d08a18e8`; 1 HIGH/1 MEDIUM/3 INFO, Fixrunde `c61142cd`/`3a5aed9a`/`4f8f17d4`). Formvorbild dieses
Reports:
[`verifikation-slice-harness-guard-blocked-python.md`](verifikation-slice-harness-guard-blocked-python.md).

**Gegenstand:** Slice-Plan `slice-harness-mutationsbild-und-verweigerte-aktion` (wellenlos), Diff-Range
`5f97700a..HEAD` — 5 Commits (`b3a23a7e` Implementierung, `d08a18e8` Review, `c61142cd`/`3a5aed9a`
Fixrunde, `4f8f17d4` DoD-Haken-Nachzug); zusätzlich 3 reine Lifecycle-Move-Commits davor
(`b345e326`/`b5c17be5`/`5f97700a`, außerhalb `5f97700a..HEAD`, aber Teil der 8 Commits ggü.
`origin/main`). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Alle
eigenen Mutationen liefen an Scratchpad-Kopien (`TOOL=<Kopie>` gegen den unveränderten Tabellentest);
kein `sed -i`, keine Umleitung auf eine Repo-Datei — Kopien per `sed 's/…/…/' Datei > Kopie` oder per
Write-Werkzeug im Scratchpad erzeugt. `git status --short` war nach jedem eigenen Lauf leer.

**Repo-Zustand:** `HEAD` ist 8 Commits vor `origin/main` (nicht mein Zutun, Auftrag verlangt kein Push).
Ich habe in diesem Lauf **nicht gepusht**.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| **H-1-Exploit gegen die Fixrunde** (Scratchpad-Kopie mit `src_abs="$src"` statt `src_abs=$(realpath -- "$src")`) gegen den **erweiterten** Tabellentest | **Exit 1, genau 3 Fälle rot** | 6 `FEHLER:`-Zeilen (2 je Fall), exakt die drei neuen Fälle „SRC relativ '.'"/„SRC relativ 'sub'"/„SRC ist Symlink auf die Repo-Wurzel" — alle 27 übrigen Aufrufe grün |
| Dieselbe Exploit-Kopie gegen den **unveränderten** (Pre-Fixrunde) Skript-Vertrag entfällt (nicht mehr im Baum); stattdessen: unveränderter Original-Tabellentest gegen unverändertes Skript | **Exit 0** | „run-image-mutation-tests: alle Fälle bestanden" |
| `make test-image-mutation` | **Exit 0** | „run-image-mutation-tests: alle Fälle bestanden"; eigene Nachzählung: 20 `expect`-, 15 `run`-, 15 `no_docker_call`-Aufrufe im Skript |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-mutationsbild-und-verweigerte-aktion.md` | **Exit 0** | „suchlauf-nachmessen: 14 Zeilen stimmen" — alle 14 Zeilen `OK`, Zahlen unverändert ggü. der Vor-Fixrunde-Messung (62/56/12/9/9/6), wie von der Fixrunde behauptet |
| `git grep -n "Ausnahme" -- <Plan-Datei>` | 7 Treffer | Zeilen 102/103 (jetzt „ohne Ausnahme … ausgeliefert in `AGENTS.md` §3.15"), 215 (DoD-Beleg-Satz, unrelated), 308 (§3.13-Ausnahmen, unrelated), 367 (Träger-Tabelle, jetzt korrigiert: „§3.15 trägt dazu **keine** Ausnahme"), 463/466 (Risiko 8, bewusst der Closure überlassen) |
| `AGENTS.md` §3.15 wörtlich gelesen (Zeilen 660–695) | **keine Ausnahmeklausel** | Wortlaut: „danach fragt der Lauf den Auftraggeber, bevor er fortfährt" — unbedingt, mit der einzigen (unveränderten) Ausnahme für einen „versehentlichen Aufruf ohne Nutzen für die Aufgabe" |
| `make gates` | **Exit 0** | alle sechs inneren Gates grün: `baseline-verify: v6.9.0 OK — 54 Dateien`, `docs-check`/`commit-traceability` (d-check `1350 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID`), `coverage-gate: OK — Coverage 85.20% erfüllt Schwelle 80%`, `generated-sync: OK — byte-gleich`, `a-check: gesamt: 0 Befund(e)` |
| `make doc-immutable RANGE=5f97700a..HEAD` | **Exit 0** | `d-check: 1350 Datei(en) geprüft, 0 Befund(e)` — `MR-003` unverändert bestätigt |
| `git diff --name-only 5f97700a..HEAD -- internal/ cmd/ gen/ test/` | **leer** | 0 Zeilen — kein Go-/Produktionscode im Diff |
| `git diff --name-only d08a18e8..HEAD` | 2 Dateien | Plan-Datei + `tools/harness/run-image-mutation-tests.sh` — die Fixrunde berührt **keine** weitere Datei, insbesondere nicht `harness/targets/image-mutation.md` (siehe V-1) |
| **10 eigene Mutationen** an Skript-Kopien (§4), davon 2 im Auftrag/Review/Fixrunde nicht genannt | 8 rot wie zugesagt, 1 rot (novel), 1 grün (novel, echte Test-Lücke) | siehe Tabelle |
| Dangling-Docker-Volumes vor/nach allen eigenen Läufen | **36 / 36** | `docker volume ls -qf dangling=true \| wc -l`, kein `prune` |

Nicht gefahren: ein realer `make image-mutation`-Bau (Docker-Bau, schwerer Lauf) — der bereits von
Implementer und Reviewer dokumentierte reale Lauf gilt als **übernommen** (§9 Punkt „Realer Lauf"), weil
das aktuelle `:dev`-Image-ID (`sha256:a8ef48aab93a…`, `docker image inspect
ghcr.io/pt9912/pg-change-feed:dev`, Created `2026-09-26T23:41:45`) mit dem im Vertrag dokumentierten
Vorher-/Nachher-Wert übereinstimmt und **vor** dem Beginn dieses Slice (Implementierung ab
2026-09-27T11:04) gebaut wurde — kein anderer Prozess hat `:dev` seither neu gebaut, ein zweiter,
schwerer Bau war für diese Bestätigung nicht nötig.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 10 `[x]`-Zeilen (Liefer-Punkte 1–3, `make gates`, Review, §3.13-Suchlauf, Doku-Update,
Reconciliation-entfällt, Beobachtungs-Register), 3 `[ ]`-Zeilen (Closure-Notiz, Risiken-Ausgänge, drei
Paarungen — korrekt offen, Planner-Pflicht bei Closure). Ich setze keinen Haken (Planner).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — das Ziel (`[x]`) | **bestätigt** | `tools/harness/image-mutation.sh` ruft exakt `docker buildx build --load -t pg-change-feed-mutation:<TAG> <SRC>` bzw. `docker rmi pg-change-feed-mutation:<TAG>` (eigen gelesen, Zeilen 82/103); alle Eingabefehler enden vor jedem Docker-Aufruf mit Exit 2 (eigene Mutationen m1–m8 + notagcheck/nodircheck, §4); `realpath`-Bindung an relativen/Symlink-`SRC` **jetzt gebunden** — eigener Exploit-Nachlauf bestätigt exakt 3 rote Fälle, 27 grüne (§1) |
| 2 | Liefer-Punkt 2 — die Träger des Ziels (`[x]`) | **bestätigt, mit einer Lücke** | (a) `harness/targets/image-mutation.md` vollständig (Zweck, Aufruf, Host-Werkzeuge, Eingabeprüfung, Exit-Codes, Anwendungsbeispiel, Grenze, Test-Tabelle, realer Lauf) — **aber** die Test-Tabelle nennt weiterhin nur „die acht oben genannten Mutationen" und trägt **keine** Zeile zur `realpath`-Bindung/H-1-Fix (V-1); (b) `harness/README.md` §Sensors zwei neue Zeilen, Link löst auf; (c) `AGENTS.md` §3.1 Satz zu `make image-mutation`, korrekter Herkunfts-Anker |
| 3 | Liefer-Punkt 3 — die Regel „Verweigerte Aktion" (`[x]`) | **bestätigt** | `AGENTS.md` §3.15 wörtlich ohne Ausnahmeklausel (§1), deckungsgleich mit der Nutzer-Entscheidung „ja"; `.harness/skills/reviewer.md` MEDIUM-Zeile „Ersatzweg nach Verweigerung ohne Meldung" mit korrekter Eskalationslogik (eigen gegen drei hypothetische Fälle geprüft, §5); `.claude/commands/implement-slice.md` Verweis-Satz an der Docker-only-Zeile; Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md` Ausgang „verkörpert" mit auflösbarem Anker |
| 4 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, kein offenes HIGH/MEDIUM (`[x]`) | **bestätigt** | H-1 und M-1 beide durch Fixrunde aufgelöst und von mir unabhängig nachgemessen (§5); Checkbox in eigenem Commit `4f8f17d4` unmittelbar nach der Fixrunde gesetzt — das ist der in `implement-slice.md` Schritt 21 „Fixrunden-Checkbox-Nachzug" vorgesehene Weg (die Nachzug-Regel selbst verlangt keinen zweiten Reviewer-Lauf, sondern die Fixrunden-Selbstbestätigung mit anschließender unabhängiger Verifikation — genau diese hier), kein Self-Review-Verstoß |
| 6 | §3.13-Suchlauf (`[x]`) | **bestätigt** | `make suchlauf-nachmessen` Exit 0, 14/14 (§1); Zahlen unverändert ggü. Vor-Fixrunde, korrekt begründet (die drei neuen Testfälle treffen kein Suchmuster) |
| 7 | Doku-Update `harness/README.md` §Sensors, `harness/targets/image-mutation.md`, `AGENTS.md` (`[x]`) | **bestätigt, mit derselben Lücke wie Zeile 2** | Benutzerhandbuch (`docs/user/`) korrekt unberührt (keine Betreiber-Oberfläche); `harness/targets/image-mutation.md` fehlt der H-1-Nachtrag (V-1) |
| 8 | Closure-Notiz (`[ ]`) | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 9 | Reconciliation-Register — entfällt (`[x]`) | **bestätigt** | keine Reconciliation-Datei im Repo (Greenfield) |
| 10 | Beobachtungs-Register fortgeschrieben (`[x]`) | **bestätigt** | `BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md`: Ausgang „verkörpert", Anker auf die drei Stellen auflösbar, „Anfall im Lauf dieses Slice: keiner" — plausibel und mit keinem Artefakt im Widerspruch (`verifizierbar: nein`, wie bereits vom Reviewer als übernommen eingestuft) |
| 11 | Jedes Risiko aus §6 (`[ ]`) | **korrekt offen** | alle zehn Zeilen tragen „**Ausgang:** *(bei Closure)*"; Vorschläge in §9 dieses Reports |
| 12 | Drei Paarungen (`[ ]`) | **korrekt offen** | Plan §7: „die Slice-Closure trägt sie selbst" |

Kein `[x]` ohne Beleg.

## 3. Plan-vs-Code-Diff

**Vollständiger Diff** (`git diff 5f97700a..HEAD --stat`): 11 Dateien, 957 Zeilen hinzugefügt, 30
entfernt. Kein Go-Code (`git diff --name-only … -- internal/ cmd/ gen/ test/` leer, §1).

| Plan-Zeile (§3) | Ist im Diff |
|---|---|
| `tools/harness/image-mutation.sh` (neu) | `A`, zwei Verben `build`/`rm`, genau die zugesagten Docker-Aufrufe, Eingabeprüfung vor jedem Aufruf (§1 gelesen) |
| `tools/harness/run-image-mutation-tests.sh` (neu) | `A` in `b3a23a7e`, `M` in `c61142cd` (drei neue Fälle 4b/4c: relativer `SRC`, Symlink) |
| `Makefile` | `M`, drei Ziele `image-mutation`/`-rm`/`test-image-mutation` mit `.PHONY` und Hilfe-Zeilen, `$(error …)` vor SRC/TAG-Fehlen |
| `harness/targets/image-mutation.md` (neu) | `A`, vollständiger Vertrag — **nicht** von der Fixrunde nachgezogen (V-1) |
| `harness/README.md` §Sensors | `M`, zwei neue Zeilen |
| `AGENTS.md` §3.1 (Mutationsprobe) | `M`, ein Satz zu `make image-mutation` |
| `AGENTS.md` (neuer Abschnitt §3.15) | `A`, ohne Ausnahmeklausel (§1, bestätigt) |
| `.harness/skills/reviewer.md` | `M`, neue MEDIUM-Klasse |
| `.claude/commands/implement-slice.md` | `M`, ein Verweis-Satz |
| `docs/plan/planning/observations/BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md` | `M`, Ausgang „verkörpert" |
| Plan-Datei selbst | `M` in `4f8f17d4`/`3a5aed9a` (DoD-Haken, Ausnahme-Korrektur) |

Kein Betreiber-Weg, kein Produktionscode, kein `.github`-Workflow im Diff.

## 4. Eigene Mutationen (Eingabeseite, Scratchpad-Kopien)

### 4a. Acht bereits vom Vertrag/Review benannte Mutationen (unabhängig nachgefahren)

| # | Mutation | Erwartet | Gesehen |
|---|---|---|---|
| m1 | `-t` auf `ghcr.io/pt9912/pg-change-feed:dev` | rot | **rot** — `gültige Argumente` Zeile 5 |
| m2 | `--metadata-file harness/image-hash.raw` angehängt | rot | **rot** — Argumentzahl-Check (6 statt 7) |
| m3 | `dev`/`latest`-Ablehnung entfernt | rot | **rot** — `TAG 'dev'`/`TAG 'latest'` |
| m4 | Zeichenklasse auf `.*` gelockert | rot | **rot** — `TAG '../x'` (u. a.) |
| m5 | Wurzel-Vergleich entfernt | rot | **rot** — `SRC ist Repo-Wurzel`/`SRC unter der Wurzel` |
| m6 | Dockerfile-Prüfung entfernt | rot | **rot** — `SRC ohne Dockerfile` |
| m7 | Docker-Exit-Auswertung (build) entfernt | rot | **rot** — `Docker-Fehler (build)` |
| m8 | `-f` an `docker rmi` angehängt | rot | **rot** — `rm`-Argumentzahl/-Zeile |

Alle acht syntaktisch validiert (`bash -n`) und einzeln gegen den unveränderten Tabellentest gefahren —
deckungsgleich mit den vom Review berichteten Ergebnissen (nicht dem Bericht geglaubt, selbst gefahren).

### 4b. Zwei eigene, im Auftrag/Review/Fixrunde nicht genannte Mutationen

| # | Mutation | Erwartet | Gesehen | Einordnung |
|---|---|---|---|---|
| notagcheck | `check_tag`-Aufruf im `rm`-Zweig entfernt (TAG-Regex/reservierte-Namen-Prüfung nur noch für `build`, nicht für `rm`) | offen | **grün** — `run-image-mutation-tests: alle Fälle bestanden` (kein Testfall ruft `rm` mit einem ungültigen `TAG`) | **reale Testlücke, neu gefunden** — siehe V-2 |
| nodircheck | `if [ ! -d "$src" ]` (Verzeichnis-Existenz) entfernt | rot | **rot** — Fall „SRC kein Verzeichnis" schlägt fehl (Meldung wechselt auf „trägt kein Dockerfile", weil `realpath` einen nicht existierenden Pfad dennoch auflöst und die nächste Prüfung greift; Exit weiterhin 2, Fall trotzdem rot) | Bindung besteht — kein Fund |

`notagcheck` ist eine reale, unabhängig von H-1 bestehende Lücke: `rm` validiert `TAG` im Code
(`check_tag "$tag"` steht in Zeile 100 des gelieferten Skripts) — kein Testfall im Tabellentest ruft
`rm` mit einem Wert, der gegen die Zeichenklasse oder die reservierten Namen verstößt, deshalb bleibt
eine versehentliche Entfernung dieser Prüfung unbemerkt. Sicherheitsrelevanz gering: `rm` operiert
ausschließlich auf dem eigenen Namensraum `pg-change-feed-mutation:<TAG>`, eine Kollision mit dem
geteilten Tag `ghcr.io/pt9912/pg-change-feed:dev` ist über den Repository-Namen bereits ausgeschlossen
— siehe V-2 (LOW, kein HIGH wie H-1, weil kein Arbeitsbaum-/Lauf-Beleg-Schaden möglich ist).

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| H-1 (HIGH) `realpath`-Bindung fehlte, Wurzel-Schutz real umgehbar | Exploit-Kopie (`src_abs="$src"`) gegen den **erweiterten** Tabellentest: Exit 1, exakt die drei neuen Fälle rot (6 `FEHLER:`-Zeilen), alle 27 übrigen grün; derselbe Exploit gegen den unveränderten Vertrag: nicht mehr reproduzierbar, da der reale Aufruf jetzt vor dem Docker-Aufruf abbricht (Fälle binden die Zeile) | **aufgelöst, bestätigt** |
| M-1 (MEDIUM) Plan behauptet eine im Auftrag verworfene Ausnahme | `AGENTS.md` §3.15 wörtlich gelesen: keine Ausnahme; Plan-Zeilen 102/103 jetzt „ohne Ausnahme … ausgeliefert"; Plan-Zeile 367 (vormals 366) jetzt „§3.15 trägt dazu **keine** Ausnahme — der Satz beschreibt die Guard-Meldung selbst" | **aufgelöst, bestätigt** |
| I-1 (INFO) Risiko 8 durch Implementierung bereits beantwortet | Plan-Zeilen 463–467 (Risiko 8) unverändert seit dem Review — der Hinweis „Ausgang sollte 'entfallen' lauten" ist **noch nicht** in eine Ausgang-Zelle übernommen, korrekt (Ausgang wird erst bei Closure gefüllt, Plan-Design) | **bestätigt, Hinweis an Planner in §9 wiederholt** |
| I-2 (INFO) Register-Angabe „kein Anfall" übernommen, nicht gemessen | `state.md` Absatz „Anfall im Lauf dieses Slice: keiner" bleibt plausibel, kein Artefakt widerspricht | **bestätigt, kein Aktionsbedarf** |

Kein offenes HIGH/MEDIUM aus dem Review.

## 6. Neue eigene Findings

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **`harness/targets/image-mutation.md` fehlt der H-1-Nachtrag.** Das Review-Verdikt verlangt: „danach `make test-image-mutation` erneut … laufen lassen und das neue Rot im Bericht/Vertrag nachtragen." Die Fixrunde hat das neue Rot im **Commit-Bericht** von `c61142cd` nachgetragen (Beleg-Text mit den sechs `FEHLER:`-Zeilen), aber **nicht im Vertrag** — die Test-Tabelle in `harness/targets/image-mutation.md` nennt weiterhin nur „die acht oben genannten Mutationen" ohne Zeile zur `realpath`-Bindung. `git diff --name-only d08a18e8..HEAD` zeigt, dass diese Datei in der Fixrunde nicht berührt wurde. Kein DoD-Bruch (die DoD-Formulierung von Liefer-Punkt 2a listet keine geschlossene Testfall-Menge, die diesen Nachtrag zwingend verlangt), aber eine reale Dokumentations-Drift ggü. der eigenen „Menge der Erprobung"-Aussage im Vertrag (jetzt untererfasst) | `harness/targets/image-mutation.md:114-138`, Review-Verdikt Zeile „das neue Rot im Bericht/Vertrag nachtragen" | ja — `git diff --name-only d08a18e8..HEAD`, Lesen |
| V-2 | LOW | **`rm` validiert `TAG` im Code, aber kein Testfall bindet das.** Eine Mutation, die `check_tag` aus dem `rm`-Zweig entfernt, bleibt für den kompletten Tabellentest unsichtbar (Exit 0, „alle Fälle bestanden") — analog zur Klasse, die H-1 für `build`/`realpath` aufgedeckt hat, hier für `rm`/`check_tag`. Sicherheitsrelevanz gering (der Namensraum `pg-change-feed-mutation:<TAG>` ist vom geteilten `:dev`-Tag bereits durch den Repository-Namen getrennt; ein ungültiger `TAG` bei `rm` könnte höchstens einen harmlosen, nicht existierenden Image-Namen an `docker rmi` übergeben) — deshalb LOW, nicht HIGH wie H-1. Kandidat für einen neunten Testfall (`rm` mit ungültigem `TAG`, z. B. `dev`/`Foo`) in einer künftigen Fixrunde oder als benannte Grenze im Vertrag | `tools/harness/run-image-mutation-tests.sh` (kein Fall ruft `rm` mit einem Wert, der gegen `TAG_RE` oder die reservierten Namen verstößt); eigene Mutation `notagcheck` (§4b) | ja — eigener Mutationslauf |
| V-3 | INFO | **10 eigene Mutationen gefahren** (8 vertragsbenannte + 2 neue: `notagcheck` deckt eine reale Lücke auf, `nodircheck` bestätigt bestehende Bindung) — zusätzlich zum H-1-Exploit-Nachlauf. Alle syntaktisch validiert (`bash -n`) vor dem Lauf | §4 | ja — eigener Lauf |
| V-4 | INFO | **Der Prozessschritt „DoD-Haken im eigenen Commit nach der Fixrunde setzen"** (`4f8f17d4`) ist der in `implement-slice.md` Schritt 21 vorgesehene „Fixrunden-Checkbox-Nachzug" — kein Self-Review-Verstoß, die unabhängige Bestätigung liefert regulär diese Verifikation | `.claude/commands/implement-slice.md:270-273` | ja — Lesen |

Kein HIGH, kein MEDIUM in dieser Verifikation — V-1/V-2 sind beide LOW und nicht merge-blockierend;
beide sind reale, aber begrenzte Lücken (Dokumentations-Nachtrag bzw. eine zusätzliche, sicherheitsarme
Testfall-Lücke derselben Klasse wie H-1), keine Wiederholung des bereits behobenen H-1-Kernproblems.

## 7. Entscheidungs-Konformität

- **Nutzer-Entscheidung 1 („ein eigenes make-Ziel mit eigenem Tag"):** `make image-mutation
  SRC=<Verzeichnis> TAG=<Tag>` baut unter dem eigenen Repository-Namen `pg-change-feed-mutation`,
  getrennt vom Lauf-Beleg-Pfad `:dev`/`harness/image-hash.txt` — bestätigt durch Code-Lesen (§1),
  8+2 eigene Mutationen (§4) und den übernommenen realen Lauf (§1). Konform.
- **Nutzer-Entscheidung 2 („ja", ohne im Wortlaut genannte Ausnahme):** `AGENTS.md` §3.15 enthält
  wörtlich keine Ausnahme für den Fall, dass die Ablehnung selbst einen Weg nennt — die einzige
  verbliebene Ausnahme betrifft einen „versehentlichen Aufruf ohne Nutzen für die Aufgabe" und ist keine
  Ausnahme von der Melde-/Rückfrage-Pflicht selbst. Konform (M-1 korrekt aufgelöst).
- **`AGENTS.md` §3.6:** keine ADR nötig (kein Gate berührt, kein Gate gelockert). Konform.
- **`AGENTS.md` §3.9:** `make gates` in diesem Lauf ungefiltert, Exit-Code direkt geprüft, vor jeder
  Folgehandlung ausgewertet. Konform.
- **`AGENTS.md` §3.7:** neue Skript-Kopfkommentare im Indikativ, höchstens ein Rang-Zeiger je Block
  (bereits vom Reviewer per `make kommentar-kennungen DIFF=5f97700a` mit 0 Kandidaten bestätigt; kein
  Go-Code in diesem Diff, daher kein eigener Nachlauf nötig).
- **`AGENTS.md` §3.13-Suchlauf:** 14/14 Zeilen bestätigt, Zahlen unverändert ggü. der Vor-Fixrunde-
  Messung (§1).
- **`harness/conventions/MR-003-…` (immutable):** unverändert (`make doc-immutable` Exit 0). Konform.
- **§3.1 Docker-only, Mutationsprobe:** kein Host-Interpreter, keine Umleitung im Diff; die eigene
  Verifikation dieses Reports lief ausschließlich mit Scratchpad-Kopien und Stub-`docker` bzw. dem bereits
  gebauten `:dev`-Image (kein neuer schwerer Docker-Bau nötig).

## 8. Register und Träger — Lese-Prüfung

- **`BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md`:** Ausgang „verkörpert" korrekt begründet,
  Anker auf `AGENTS.md` §3.15, die Prüfzeile in `.harness/skills/reviewer.md` und den Verweis in
  `.claude/commands/implement-slice.md` auflösbar (§2). Zähler „1× (evidence/…)" korrekt (der Ursprungsfall
  aus `slice-leerlauf-phase-last-in-stuecken`).
- **`BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet`:** bewusst getrennt gehalten, der
  Regel-Satz „Aufruf, Pfad, Wortlaut" (§3.15) nennt ihn als Nachbar — korrekt, kein Vermischen.
- **`.harness/skills/reviewer.md` neue MEDIUM-Klasse:** Eskalationslogik selbst gegen drei hypothetische
  Fälle geprüft — (a) verweigerter `docker push` ohne Meldung, Ersatzweg über ein anderes Kommando → HIGH
  (Push berührt geteilten Zustand); (b) verweigertes `Write` auf eine reine Scratchpad-Datei, Ersatzweg
  über `Bash echo >` → MEDIUM; (c) verweigerter Aufruf gemeldet, aber ohne vorherige Rückfrage fortgesetzt
  → LOW (sofern keine HIGH-Klasse). Widerspruchsfrei, deckt alle drei Fälle eindeutig ab.
- **Plan-Träger-Tabelle (§3):** die „fremde Datei"-Zeile zu `slice-harness-guard-blocked-python` korrekt
  als erledigt vermerkt (der andere Slice ist inzwischen nach `done/` verschoben, bestätigt über den
  vorherigen Verifikations-Report desselben Slice).

## 9. Verdikt

**DoD bestätigt:** ja, in der Substanz — alle zehn `[x]`-Zeilen tragen einen realen, von mir unabhängig
nachgemessenen Beleg; die drei `[ ]`-Zeilen sind korrekt der Planner-Closure vorbehalten. Zwei LOW-Funde
(V-1: Vertrag nicht mit dem H-1-Nachtrag aktualisiert; V-2: eine zweite, sicherheitsarme Testlücke
derselben Klasse bei `rm`/`check_tag`) sind real, aber nicht DoD-verletzend und nicht merge-blockierend.

**Plan-vs-Code:** keine unbenannte Abweichung; der Diff beschränkt sich exakt auf die im Plan §3
zugesagten Dateien, kein Produktionscode berührt.

**Entscheidungs-Konformität:** beide Nutzer-Entscheidungen exakt umgesetzt — eigenes Ziel mit eigenem
Tag (Entscheidung 1), Regel ohne Ausnahme (Entscheidung 2, M-1 korrekt aufgelöst).

**Review-Findings:** H-1 (HIGH) und M-1 (MEDIUM) beide unabhängig nachgemessen und bestätigt aufgelöst
(exakter Exploit-Nachlauf, wörtliche Lesung von §3.15 und der Plan-Korrekturen); I-1/I-2 bestätigt ohne
Aktionsbedarf über das in §9 Genannte hinaus.

**Mutationen:** 10 eigene Mutationen (8 vertragsbenannte + 2 neue) plus der exakte H-1-Exploit-Nachlauf,
alle wie erwartet — mit einer neu gefundenen, sicherheitsarmen Testlücke (V-2).

**Gates:** `make gates`, `make suchlauf-nachmessen`, `make test-image-mutation`, `make doc-immutable
RANGE=5f97700a..HEAD` — alle Exit 0 im eigenen Lauf. Dangling-Docker-Volumes unverändert (36/36), kein
`prune`, kein Push.

### Übergabe an den Planner

1. **DoD-Haken:** die zehn `[x]`-Zeilen sind belegt (Tabelle §2); die drei `[ ]`-Zeilen (Closure-Notiz,
   Risiken-Ausgänge, drei Paarungen) bleiben regulär bei der Closure zu setzen.
2. **Risiken §6 des Plans — Ausgang-Vorschläge:**
   - **Risiko 1** (Symlink/relativer Pfad umgeht Wurzel-Schutz): **eingetreten, in der Fixrunde behoben**
     — H-1 demonstrierte den realen Exploit, die drei neuen Testfälle binden ihn jetzt (§1/§5 dieses
     Reports).
   - **Risiko 2** (Tag-Kollision bei parallelen Läufen): **nicht eingetreten**, Grenze im Vertrag
     dokumentiert (Namensvorschlag im Anwendungsbeispiel), kein Mechanismus im Ziel selbst nötig laut
     Plan-Design.
   - **Risiko 3** (Berechtigungsschicht verweigert das neue Ziel): **nicht eingetreten** — Register
     `BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md` bestätigt „keine Ablehnung während der
     Implementierung" (übernommen, plausibel, kein Widerspruch).
   - **Risiko 4** (`:dev`/`harness/image-hash.txt` überschrieben): **nicht eingetreten**, real bestätigt
     (Image-ID von `:dev` unverändert seit vor dem Slice-Beginn, §1).
   - **Risiko 5** (Stub bindet nur Argumente, nicht Docker-Verhalten): **erfüllt/entfallen** — der reale
     Lauf (übernommen, §1) plus die Mutationsliste schließen die Lücke.
   - **Risiko 6** (`.dockerignore`-Default-Deny blockiert eine neue Datei in der Kopie): **weiter offen**
     — der reale Lauf mutierte eine bestehende Go-Zeile, keine neue Datei; die Grenze steht dokumentiert
     im Vertrag (§Grenze Punkt 3), aber unbelegt real getestet. Empfehlung: als dokumentierte Grenze
     akzeptieren (kein Blocker), nicht als „entfallen" verbuchen.
   - **Risiko 7** (Regel „Verweigerte Aktion" verspricht mehr als ihr Sensor): **wie geplant, dokumentierte
     Grenze** — `AGENTS.md` §3.15 §Grenze-Absatz vorhanden und korrekt.
   - **Risiko 8** (Ausnahme zu weit gefasst): **entfallen** — die Ausnahme wurde nicht ausgeliefert (M-1),
     die Frage „ist sie zu weit?" ist damit gegenstandslos (deckt sich mit I-1 des Reviews).
   - **Risiko 9** (neue Fälle laufen nicht): **entfallen, bestätigt real gelaufen** — eigene Nachzählung
     (20 `expect`, 15 `run`) und der Exploit-Nachlauf zeigen, dass die drei neuen Fälle tatsächlich aktiv
     binden.
   - **Risiko 10** (Träger bleibt stehen — `plan-welle.md`/`close-welle.md`): **entfallen, wie geplant
     nicht angefasst**, kein Fund im Review oder dieser Verifikation, der einen Nachtrag dort verlangt.
3. **V-1 (LOW):** `harness/targets/image-mutation.md` §Test um eine Zeile zur `realpath`-Bindung
   ergänzen (Mutation „realpath-Auflösung von SRC entfernt" → Fälle „SRC relativ '.'/'sub'"/„SRC ist
   Symlink auf die Repo-Wurzel"), „Menge der Erprobung: acht" auf elf korrigieren — kleiner Nachtrag,
   kein Blocker für Closure.
4. **V-2 (LOW):** neuer Testfall (oder benannte Grenze im Vertrag) für `rm` mit ungültigem `TAG` — reale,
   aber sicherheitsarme Lücke derselben Klasse wie H-1; guter Kandidat für den Steering-Loop-Lerneintrag
   der Closure-Notiz („geschärfte Regel": jede Eingabeprüfung, die in **beiden** Verben (`build`/`rm`)
   auftritt, braucht Testfälle in **beiden** Verben, nicht nur im zuerst geschriebenen).
5. **Drei Paarungen:** bleiben an der Slice-Closure selbst hängen (Plan §7, korrekt).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
