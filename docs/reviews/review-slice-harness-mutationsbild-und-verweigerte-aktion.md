# Review-Report: slice-harness-mutationsbild-und-verweigerte-aktion — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, die zwei Nutzer-Entscheidungen der Sitzung 2026-09-27 (Weg 1: „ein
eigenes make-Ziel mit eigenem Tag“; Weg 2: „ja“ zur Regel „verweigerter Aufruf wird gemeldet und vor einem
Ersatzweg zurückgefragt“, **ohne** Ausnahme) und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-harness-mutationsbild-und-verweigerte-aktion` (wellenlos), Diff-Range
`5f97700a..HEAD`: Lifecycle `b345e326`/`b5c17be5`/`5f97700a` (davon `5f97700a` = Basis, nicht Teil des
geprüften Diffs), Implementierung `b3a23a7e` (Skript + Tabellentest + Makefile-Ziele + `AGENTS.md` §3.1/§3.15 +
Reviewer-Skill + Implementer-Command + Träger + Register-Nachzug + Plan-Fixes). 10 geänderte Dateien, kein
Go-Code im Diff.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere HIGH-/MEDIUM-Klassen
ergänzt — darunter die von diesem Slice selbst neu eingeführte MEDIUM-Klasse „Ersatzweg nach Verweigerung ohne
Meldung“. **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle Mutationen liefen an
Scratchpad-Kopien von `tools/harness/image-mutation.sh` (`TOOL=<Kopie>` gegen den unveränderten Tabellentest);
nie `sed -i`, nie eine Umleitung auf eine Repo-Datei — Kopien per `sed 's/…/…/' Datei > Kopie` im
Scratchpad-Verzeichnis erzeugt.

**Eigene, während des Reviews vom PreToolUse-Guard verweigerte Aktion (AGENTS.md §3.15, Selbstanwendung):**
Ein Werkzeugaufruf-Batch, der eine überflüssige Zeile `python3 -c "print('noop-not-used')"` neben mehreren
`sed`-Aufrufen enthielt, wurde vom Guard vollständig geblockt (Meldung: „This repository is make/Docker-only …
Use a repo file with the Edit/Write tools or a repo tool behind make; to try a change on a copy, write to
stdout: sed s/a/b/ file > /path/to/scratch-copy“). Die Python-Zeile war ohne Nutzen für die Aufgabe (Rest eines
vorherigen Entwurfs); ich habe sie gestrichen, **nicht** auf einem anderen Weg wiederholt, und die verbliebenen
`sed`-Aufrufe (bereits der zugesagte, nicht verbotene Weg — Ausgabe nach stdout auf eine Scratchpad-Kopie) im
nächsten Batch ohne die Python-Zeile erneut gestellt. Das fällt unter die Ausnahme in `AGENTS.md` §3.15 letzter
Satz („Ein versehentlicher Aufruf ohne Nutzen für die Aufgabe wird gestrichen und im Bericht genannt, nicht
wiederholt“) — keine Rückfrage nötig, hier als Meldung nachgetragen.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-harness-mutationsbild-und-verweigerte-aktion` (§1 Ziel/Ausgangslage/Kernaussagen, §2 DoD,
  §3 Plan mit Suchlauf-Feld und Träger-Tabelle, §6 Risiken)
- `AGENTS.md` §3.1 (Docker-only, Mutationsprobe), §3.6, §3.7, §3.9, §3.12, §3.13, neu §3.15
- `harness/conventions/MR-003-guard-inplace-textwerkzeug.md` (unverändert, `Accepted`, immutable — zitiert im
  Plan als Beleg für den Satz „nennen diesen Ersatzweg“)
- `ADR-0044`, `ADR-0083`, `ADR-0100` (Rang-Zeiger §3.14)
- Beobachtungs-Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion` (Träger dieses Slice),
  `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet` (Nachbar, bewusst getrennt gehalten)
- Zwei Nutzer-Entscheidungen der Sitzung 2026-09-27 (Wortlaut im Plan-Kopf zitiert)
- `docs/reviews/review-slice-leerlauf-phase-last-in-stuecken.md` F-3 (Ursprungsfall der Regel)

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

- **`make gates`**: voller Lauf, Exit 0 — `baseline-verify` (54 Dateien), `docs-check` (1349 Dateien, 0 Befunde,
  zweimal: einmal isoliert, einmal als Teil des Gate-Laufs), `--enable commits` (0 Befunde), `commit-traceability`
  (5 Commits, keine Struktur-ID im Betreff), `coverage-gate` (85.30 % gegen Schwelle 80 %), `generated-sync`
  (byte-gleich), `a-check` (0 Befunde). Exit-Code direkt und ungefiltert gelesen (`AGENTS.md` §3.9).
- **`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-mutationsbild-und-verweigerte-aktion.md`**:
  Exit 0, „14 Zeilen stimmen“ — alle 14 Zeilen (7 `84f60e6f`-Baseline, 7 `diff`) einzeln nachgerechnet, keine
  Abweichung, auch nicht bei den fünf vom Implementer selbst als korrigiert gemeldeten Werten (9/62/6/9/56).
- **`make docs-check`** (isoliert): 1349 Dateien, 0 Befunde — Links (u. a. auf
  `review-slice-leerlauf-phase-last-in-stuecken.md` und `harness/targets/image-mutation.md`) lösen auf.
- **`make kommentar-kennungen DIFF=5f97700a`**: Exit 0, keine Ausgabe — kein Kandidat (kein Go im Diff, wie vom
  Implementer behauptet und hier bestätigt).
- **`make doc-immutable RANGE=5f97700a..HEAD`**: Exit 0, 1349 Dateien, 0 Befunde — keine `Accepted`-ADR verändert,
  `MR-003` unverändert.
- **`make test-image-mutation`**: Exit 0, „alle Fälle bestanden“ (kein gedruckter Zähler — deckungsgleich mit dem
  Vorbild-Verhalten von `make test-fmt-check`/`make test-suchlauf-nachmessen`, die ebenfalls keine Zahl drucken;
  kein Formfehler). Eigene Nachzählung der `expect`-Aufrufe im Skript: 18 benannte Assertionen plus 2 eingebettete
  Fühler-Prüfungen (Hash-Dateien unverändert) über die vom Implementer/Vertrag genannten sieben Fallgruppen
  (gültige Argumente, Hash-Dateien unverändert, TAG-Klasse mit 8 Varianten, SRC-Fehler mit 7 Varianten,
  Docker-Fehler, `rm`, Make-Ebene) verteilt. Ich finde nirgends im Diff, im Commit oder in den Trägern eine
  Behauptung von „14 Fälle“ für diesen Tabellentest — die einzige „14“ in diesem Slice ist die Zahl der
  Suchlauf-Zeilen (siehe oben), ein anderer Gegenstand. Kein Befund, nur Klarstellung.
- **8 vom Vertrag (`harness/targets/image-mutation.md` §Test) benannte Mutationen selbst nachgefahren**, alle
  acht rot wie zugesagt: `-t` auf `ghcr.io/pt9912/pg-change-feed:dev` → „gültige Argumente“ Zeile 5 rot;
  `--metadata-file harness/image-hash.raw` angehängt → Argumentzahl-Check rot; `dev`/`latest`-Ablehnung entfernt
  → `TAG 'dev'`/`TAG 'latest'` rot; Zeichenklasse auf „alles“ gelockert → `TAG '../x'` rot (Docker-Argument
  enthält den nicht abgefangenen Wert); Wurzel-Vergleich entfernt → `SRC unter der Wurzel` rot (baut real aus
  einem Unterverzeichnis der Wegwerf-Wurzel); Dockerfile-Prüfung entfernt → `SRC ohne Dockerfile` rot;
  Docker-Exit-Auswertung entfernt → `Docker-Fehler (build)` **und** `Docker-Fehler (rm)` rot; `-f` an `docker rmi`
  angehängt → `rm`-Argumentzahl/-Zeile rot. Reproduziert exakt die im Vertrag genannten acht Mutationen.
- **2 eigene, im Vertrag/Plan nicht genannte Mutationen** (Auftrag verlangte mindestens zwei zusätzliche):
  1. **Doppeltes `--load`** (`docker buildx build --load --load -t …`) → Argumentzahl-Check rot (7 statt 6
     Zeilen). Erwartbar, bestätigt die Bindung der Argumentzahl-Prüfung.
  2. **`realpath`-Auflösung von `SRC` entfernt** (`src_abs="$src"` statt `src_abs=$(realpath -- "$src")`) →
     **der gesamte Tabellentest bleibt grün** (`make test-image-mutation` mit dieser Mutation: Exit 0, „alle
     Fälle bestanden“). Grund: **jeder** Testfall übergibt `SRC` bereits als absoluten Pfad (`$srcdir`, `$wrepo`,
     `$wrepo/sub`, `$tmp/…` sind alle unter `$tmp`, das `mktemp -d` bereits absolut liefert) — kein einziger Fall
     füttert einen **relativen** `SRC`-Wert. Siehe Finding H-1.

## Findings

### H-1 — „SRC als aufgelöster absoluter Pfad“ ist eine ungebundene Zusage: kein Testfall füttert einen relativen `SRC`-Wert, und das Entfernen der `realpath`-Auflösung hebelt den Wurzel-Schutz real aus

- `kategorie`: HIGH
- `quelle`: Skill-HIGH „Zusage ohne Bindung an ihre Eingabeseite — grün ohne Aussage“; DoD Liefer-Punkt 1 („SRC
  als aufgelöster absoluter Pfad“); Plan §6 Risiko 1 („ein relativer Pfad könnte ihn umgehen, wenn er nicht
  aufgelöst wird“, *„zu belegen durch: … ein Fall mit einem Symlink auf die Wurzel“*); `AGENTS.md` §3.1
  (Mutationsprobe „arbeitet auf einer Kopie“ — genau der Schutz, den dieses Ziel für Images durchsetzen soll)
- `pfad`: `tools/harness/image-mutation.sh:67` (`src_abs=$(realpath -- "$src") || exit 2`); `tools/harness/run-image-mutation-tests.sh`
  (kein Testfall mit relativem oder Symlink-`SRC` — jeder Aufruf übergibt einen bereits absoluten Pfad unter
  `$tmp`/`$wrepo`)
- `befund`: Der Vertrag und die DoD versprechen: „`SRC` wird vor dem Aufruf zu einem absoluten Pfad aufgelöst“ —
  das ist die einzige Verteidigung gegen einen relativen `SRC`, der versehentlich mit der Repo-Wurzel identisch
  ist (`SRC=.` von der Wurzel aus aufgerufen). Ich habe eine Kopie des Skripts erzeugt, in der `src_abs=$(realpath
  -- "$src")` durch `src_abs="$src"` ersetzt ist (die `realpath`-Auflösung entfällt vollständig, alle
  Folgeprüfungen bleiben unverändert), und sie gegen den unveränderten Tabellentest laufen lassen
  (`TOOL=<Kopie> bash tools/harness/run-image-mutation-tests.sh`): **Exit 0, „alle Fälle bestanden“** — die
  Mutation ist für den ganzen Tabellentest unsichtbar. Anschließend habe ich den Exploit real gefahren: in einem
  frischen Wegwerf-Git-Repo (mit `Dockerfile`/`go.mod`) ruft die **unveränderte** Originalfassung mit
  `SRC=.` korrekt `image-mutation: SRC '.' ist die Repo-Wurzel oder liegt unter ihr` auf und bricht mit Exit 2 ab
  — genau die Zusage aus `AGENTS.md` §3.1 („eine Mutation baut nie aus dem Arbeitsbaum“). Dieselbe Eingabe gegen
  die **mutierte** Fassung (ohne `realpath`) läuft dagegen unbeanstandet durch `docker buildx build --load -t
  pg-change-feed-mutation:exploittag .` und baut real ein Image aus dem Arbeitsverzeichnis (Exit 0, Image
  existierte danach unter `docker image ls`, von mir wieder mit `docker rmi` entfernt). Der Wurzel-Vergleich
  selbst ist korrekt geschrieben (String-Vergleich zweier bereits aufgelöster Pfade); die **einzige** Verteidigung
  gegen einen relativen `SRC` ist die `realpath`-Zeile, und genau die hat keinen Testfall, der sie an ihrer
  Eingabeseite bindet — Plan-Risiko 1 benennt exakt dieses Szenario („ein relativer Pfad könnte ihn umgehen“) und
  verlangt explizit einen Symlink-Testfall („*zu belegen durch* … ein Fall mit einem Symlink auf die Wurzel“); ein
  Symlink- oder auch nur ein einfacher relativer-Pfad-Testfall fehlt im gelieferten Tabellentest vollständig.
  `make image-mutation SRC=<Verzeichnis> TAG=<Tag>` wird uneingeschränkt mit `SRC` aus einem
  Rollenlauf-Aufruf gefüttert (`$(SRC)` im Makefile ohne eigene Normalisierung) — ein Aufrufer, der versehentlich
  einen relativen statt absoluten Pfad übergibt (z. B. `SRC=.` oder `SRC=../kopie`), verlässt sich auf genau die
  ungetestete Zeile.
- `verifizierbar`: ja — Mutation `src_abs="$src"` gegen `TOOL=<Kopie>`: Tabellentest bleibt grün (Exit 0); realer
  Exploit-Lauf in einem Wegwerf-Git-Repo mit `SRC=.`: Original lehnt ab (Exit 2, „Repo-Wurzel“), Mutation baut
  real (Exit 0, Image erzeugt)
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (relativer `SRC`/`realpath`)

### M-1 — Der Plan behauptet an zwei Stellen (§1, §3-Träger-Tabelle) weiterhin die im Auftrag verworfene Ausnahme „führt die Ablehnung selbst einen Weg an, gilt dieser Weg“, obwohl das ausgelieferte `AGENTS.md` §3.15 keine Ausnahme enthält

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“; Skill-HIGH „Beleg trägt seinen
  Satz nicht“ (Nachbarklasse — die Plan-Zeile zitiert `AGENTS.md` §3.15 als Träger eines Satzes, den §3.15 nicht
  trägt); Nutzer-Entscheidung „ja“ (ohne im Wortlaut genannte Ausnahme)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-mutationsbild-und-verweigerte-aktion.md:98-106`
  („Kernaussage der Regel (Entwurf)“ — „Führt die Ablehnung selbst einen Weg an … gilt dieser Weg; jeder andere
  braucht die Freigabe“), `:366` (Träger-Tabelle, Zeile zu `MR-003`: „die Ausnahme „Ablehnung führt den Weg selbst
  an“ in `AGENTS.md` §3.15 trägt genau diesen Satz“) — gegen `AGENTS.md:660-694` (§3.15, ausgeliefert, **ohne**
  jede Ausnahmeklausel)
- `befund`: Ich habe `AGENTS.md` §3.15 wörtlich gelesen: Es steht dort **keine** Ausnahme für den Fall, dass „die
  Ablehnung selbst einen Weg nennt“ — der Text verlangt unbedingt, dass der Lauf erst den Auftraggeber fragt,
  bevor er einen anderen Weg versucht, mit der einzigen (unveränderten) Ausnahme für einen „versehentlichen Aufruf
  ohne Nutzen für die Aufgabe“. Das ist korrekt und deckungsgleich mit der Nutzer-Entscheidung „ja“ ohne im
  Wortlaut genannte Ausnahme (der Auftrag dieses Reviews bestätigt das explizit als Implementer-Entscheidung).
  Der Plan selbst wurde dabei aber **nicht** nachgezogen: Zeile 98–106 (als „Entwurf“ deklariert, insofern mit
  geringerem Gewicht) beschreibt die Ausnahme weiterhin als Teil der Kernaussage, und **Zeile 366** — eine
  Tabellenzelle der Träger-Tabelle, keine als Entwurf gekennzeichnete Prosa — behauptet als Tatsache, dass „die
  Ausnahme … in `AGENTS.md` §3.15 trägt genau diesen Satz“. Das ist falsch: die genannte Stelle (§3.15) trägt
  diesen Satz nicht. Beide Textstellen liegen in derselben Datei, die der Implementer im selben Commit an anderer
  Stelle (der „Nachmessung durch den Implementer“-Absatz direkt unterhalb der Träger-Tabelle, Zeilen ~333–356)
  bearbeitet hat, ohne die stehen gebliebene, jetzt widersprüchliche Zeile 366 zu korrigieren oder auf die
  Abweichung zu verweisen. Die Wirkung ist auf den Plan-Text begrenzt — die ausgelieferte Regel selbst ist korrekt
  und entscheidungskonform —, aber ein späterer Leser des Plans (z. B. bei einer Closure-Notiz oder einem
  Steering-Loop-Rückblick) würde aus Zeile 366 einen falschen Schluss über den Inhalt von §3.15 ziehen.
- `verifizierbar`: ja — beide Textstellen direkt lesbar; Gegenprobe an `AGENTS.md:660-694`
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger (Plan vs. ausgelieferte Regel)

### I-1 — Plan-Risiko 8 bleibt konsistent mit M-1, aber seine Ausgangsfrage ist durch die Implementierung bereits beantwortet

- `kategorie`: INFO
- `quelle`: Plan §6 Risiko 8 („Die Ausnahme … ist Auslegung des Planners, nicht Wortlaut der Nutzer-Entscheidung“)
- `pfad`: `docs/plan/planning/in-progress/slice-harness-mutationsbild-und-verweigerte-aktion.md:462-466`
- `befund`: Risiko 8 fragt, ob die geplante Ausnahme zu weit greift, und verweist die Antwort auf den Reviewer.
  Da die Ausnahme in der ausgelieferten Fassung von §3.15 gar nicht existiert (siehe M-1), ist die im Risiko
  gestellte Frage gegenstandslos geworden — der Ausgang bei Closure sollte „entfallen“ lauten (die Ausnahme wurde
  nicht implementiert), nicht „bei Closure zu bewerten, ob sie zu weit greift“. Kein eigener Fixbedarf, da die
  Ausgang-Spalte ohnehin erst bei Closure gefüllt wird; als Hinweis für den Planner/die Closure-Notiz.
- `verifizierbar`: ja — Text von Risiko 8 gegen §3.15 gelesen
- `klasse`: Risiko-Ausgang durch Implementierung bereits beantwortet, noch nicht nachgezogen

### I-2 — Beobachtungsregister-Eintrag „Anfall im Lauf dieses Slice: keiner“ und Risiko 3 stimmen mit dem eigenen Review-Verlauf überein

- `kategorie`: INFO
- `quelle`: `docs/plan/planning/observations/BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md`; Plan §6 Risiko 3
- `pfad`: state.md, Absatz „Anfall im Lauf dieses Slice“
- `befund`: Die Behauptung, die Berechtigungsschicht habe während der Implementierung (`Bash(make:*)` erlaubt
  `make image-mutation`) keinen Aufruf verweigert, kann ich aus dem Repo allein nicht unabhängig nachmessen (kein
  Log-Träger) — sie ist `verifizierbar: nein`/**übernommen**, wie die Skill-Regel für diese Klasse selbst festhält.
  Sie ist plausibel (Allow-Liste enthält `Bash(make:*)`, kein Docker-Unterkommando-Pattern nötig, da `make`
  aufgerufen wird) und widerspricht keinem Artefakt. Während meines eigenen Reviews trat dagegen eine Verweigerung
  auf (oben, „Eigene … verweigerte Aktion“) — sie betraf eine geschachtelte `python3`-Zeile in meinem eigenen
  Mutationstest-Kommando, nicht den Implementer-Lauf, und ändert nichts an der Bewertung des Diffs.
- `verifizierbar`: nein (übernommen)
- `klasse`: Übernommene Angabe des Auftraggebers/Registers, kein Repo-Beleg

## Negativbefunde

- geprüft, ohne Befund: **`tools/harness/image-mutation.sh` — Docker-Aufruf-Form.** `build` ruft exakt
  `docker buildx build --load -t pg-change-feed-mutation:<TAG> <SRC>` (6 Argumente, `SRC` als letztes, aufgelöst),
  kein `--metadata-file`, kein `--push`, kein `--platform`, kein `--build-arg`; `rm` ruft exakt `docker rmi
  pg-change-feed-mutation:<TAG>` (2 Argumente), kein `-f`, kein `prune` — durch den Stub-`docker` (exakte
  Zeilenzahl **und** Zeileninhalt je Argument) und durch 8+2 eigene Mutationen bestätigt.
- geprüft, ohne Befund: **TAG-Zeichenklasse und reservierte Namen.** `^[a-z0-9][a-z0-9_.-]{0,62}$`, `dev`/`latest`
  abgelehnt, alle acht im Tabellentest genannten Schlechtfälle (leer, `dev`, `latest`, `Foo`, `a b`, `-x`, `a:b`,
  `../x`) enden mit Exit 2 vor jedem Docker-Aufruf (`no_docker_call`-Prüfung).
- geprüft, ohne Befund: **SRC-Eingabeprüfung (außer der `realpath`-Lücke aus H-1).** Fehlt/kein
  Verzeichnis/ohne `Dockerfile`/ohne `go.mod`/Repo-Wurzel/Verzeichnis unter der (bereits absoluten) Wurzel — alle
  Fälle enden mit Exit 2, kein Docker-Aufruf.
- geprüft, ohne Befund: **Exit-Code-Disziplin.** Ein Docker-Fehler (`build` und `rm`) endet mit Exit 1 und nennt
  den Exit-Code in der Meldung; jeder Eingabefehler endet mit Exit 2 **vor** jedem Docker-Aufruf; die Make-Ziele
  brechen ohne `SRC`/`TAG` mit `$(error …)` ab, bevor überhaupt ein Rezept läuft (drei eigene Läufe: `make
  image-mutation` ohne `SRC`, ohne `TAG`, `make image-mutation-rm` ohne `TAG` — jeweils Exit 2, kein Rezept
  sichtbar in der Ausgabe).
- geprüft, ohne Befund: **Fühler `harness/image-hash.txt`/`.raw`.** Der Fühler prüft den Dateiinhalt nach dem Lauf
  unabhängig vom Stub — der Stub fängt nur den `docker`-Aufruf ab, ein hypothetisches zusätzliches
  Datei-Schreibkommando im Skript liefe real gegen die Wegwerf-Wurzel (`cd "$wrepo"` vor dem Aufruf) und würde vom
  Fühler erfasst; das Skript selbst schreibt an keiner Stelle in eine Datei.
- geprüft, ohne Befund: **Stub-Bindung.** `args_count`/`args_line` binden Argumentzahl **und** -inhalt exakt (nicht
  nur „enthält“); die doppelte-`--load`-Mutation zeigt das (Argumentzahl 7 statt 6 schlägt fehl, obwohl Inhalt und
  Reihenfolge der übrigen Argumente unverändert blieben).
- geprüft, ohne Befund: **`harness/targets/image-mutation.md`.** Vollständiger Vertrag (Zweck, Aufruf, Host-
  Werkzeuge, Eingabeprüfung, Exit-Codes, Anwendungsbeispiel als Text ohne echten Bau, Grenze, Test-Tabelle,
  realer Lauf mit Beleg-Ankern); die Compose-Override im Anwendungsbeispiel ist reiner Text, kein ausführbarer
  Schritt — deckungsgleich mit dem Plan („Anwendungsbeispiel … als Text“).
- geprüft, ohne Befund: **`harness/README.md` §Sensors.** Zwei neue Zeilen (`make image-mutation`/`-rm`,
  `make test-image-mutation`), „kein Gate“ korrekt benannt, Link auf den Vertrag löst auf (`make docs-check`).
- geprüft, ohne Befund: **`AGENTS.md` §3.1, Absatz zur Mutationsprobe.** Der neue Satz beschreibt exakt das
  getestete Verhalten (eigener Repository-Name, `:dev`/`harness/image-hash.txt` unberührt), Herkunfts-Anker
  korrekt.
- geprüft, ohne Befund: **`.harness/skills/reviewer.md` neue MEDIUM-Klasse.** Selbst angewendet auf drei
  hypothetische Fälle: (a) verweigerter `docker push` ohne Meldung/Rückfrage, Ersatzweg über ein anderes
  Kommando → HIGH (Push berührt geteilten Zustand, zweites Feld neben der HIGH-Liste des Ziels der Aktion); (b)
  verweigertes `Write` auf eine reine Scratchpad-Datei ohne Meldung, Ersatzweg über `Bash echo >` → MEDIUM
  (weder HIGH-Klasse noch geteilter Zustand); (c) verweigerter Aufruf, im Bericht genannt, aber ohne vorherige
  Rückfrage direkt mit einem anderen Weg fortgefahren → LOW (sofern die Aktion selbst keine HIGH-Klasse trägt).
  Die Einstufungslogik ist widerspruchsfrei und deckt alle drei Fälle eindeutig ab.
- geprüft, ohne Befund: **`.claude/commands/implement-slice.md`.** Der neue Verweis-Satz sitzt am Ende der
  Docker-only-Zeile (Zeile 32–33, unmittelbar nach dem Satz „Rufe nur `make`-Targets auf.“), wiederholt die Regel
  nicht (keine Falsch/Richtig-Beispiele, keine Begründung, keine Grenze), sondern verweist nur — an der einzigen
  Stelle, die der Plan dafür vorsieht.
- geprüft, ohne Befund: **Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md`.** Ausgang „verkörpert“
  korrekt begründet, Anker auf die drei Stellen (Regel, Prüfzeile, Verweis) auflösbar, kein neuer `evidence/`-
  Eintrag nötig (kein Anfall).
- geprüft, ohne Befund: **Fremde-Datei-Link-Bereinigung.** Kein verbleibender Markdown-Link
  (`](…slice-harness-guard-blocked-python.md)`) außerhalb von `done/`/`docs/reviews/**`; alle verbleibenden
  Erwähnungen sind reine Kennungs-Zitate, keine Links.
- geprüft, ohne Befund: **Traceability, ID-Schema, Moves.** Alle vier Commits tragen `(ADR-0083)`, keine
  `SPEC-*`/`ARC-*`; die beiden Lifecycle-Commits `b345e326`/`5f97700a` sind reine Renames („Kein Inhalt
  geaendert“); kein `//nolint` (kein Go im Diff).
- geprüft, ohne Befund: **§3.7-Kommentar-Klassen.** Beide neuen Skript-Kopfkommentare im Indikativ, höchstens ein
  Rang-Zeiger (`AGENTS.md` §3.1) bzw. eine Testfall-Provenienz-Nennung („Nachbild: run-fmt-check-tests.sh Fall
  11“ — zulässige Form), keine Kette, kein „ff.“, keine Vorher/Nachher-Sprache; `make kommentar-kennungen
  DIFF=5f97700a` bestätigt 0 Kandidaten.
- geprüft, ohne Befund: **Spec-Stratum, Zwei-Quellen-Drift, Suppression, kritischer Pfad.** `spec/` nicht
  berührt; keine neue `Accepted`-ADR verändert (`make doc-immutable` grün); kein `//nolint`; kein
  Produktionscode/`internal/` im Diff; kein Geheimnis.
- geprüft, ohne Befund: **Dangling-Docker-Ressourcen.** `docker volume ls -f dangling=true -q | wc -l` vor und
  nach allen eigenen Docker-Läufen dieses Reviews: 36 (unverändert); kein `prune` verwendet; alle selbst gebauten
  Wegwerf-Images (Stub-Läufe sowie der reale Exploit-Bau) wieder mit `docker rmi` entfernt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite (relativer `SRC`/`realpath`) ·
Nachzug widerspricht dem Nachbarn im selben Träger (Plan vs. ausgelieferte Regel) · Risiko-Ausgang durch
Implementierung bereits beantwortet, noch nicht nachgezogen · Übernommene Angabe des Auftraggebers/Registers,
kein Repo-Beleg

## Antwort auf die Ausnahme-Klausel-Frage (Kernauftrag dieses Reviews)

**`AGENTS.md` §3.15 enthält keine Ausnahme** — der Wortlaut verlangt unbedingt, dass ein Lauf nach jeder
Verweigerung meldet und den Auftraggeber fragt, bevor er einen anderen Weg versucht; die einzige verbliebene
Ausnahme betrifft ausschließlich einen „versehentlichen Aufruf ohne Nutzen für die Aufgabe“ (streichen statt
wiederholen, ohne Rückfragepflicht) und ist keine Ausnahme von der Melde-/Rückfrage-Pflicht selbst, sondern eine
Sonderregel für einen Fall, in dem gar kein Ersatzweg verfolgt wird. Das ist **konsistent mit der
Nutzer-Entscheidung „ja“** in der vom Planner referenzierten, exceptionslosen Lesart. Die Regel ist außerdem in
sich stimmig: Falsch/Richtig-Beispiele treffen den Regelkern, die Grenze („wirkt durch Lesen; kein Sensor; der
Classifier ist kein Teil des Repos“) ist vorhanden und verweist korrekt auf `ADR-0083`, und die Herkunft
(`BEO-PGC/ersatzweg-nach-verweigerter-aktion` · seit slice-harness-mutationsbild-und-verweigerte-aktion) löst auf.

**Der Plan selbst wurde dabei nicht vollständig nachgezogen** (M-1): §1 („Kernaussage der Regel (Entwurf)“) und
insbesondere die Träger-Tabellenzeile zu `MR-003` in §3 behaupten weiterhin bzw. als Tatsache, dass die
Ausnahme „führt die Ablehnung selbst einen Weg an, gilt dieser Weg“ in `AGENTS.md` §3.15 stehe — das ist seit der
Implementierung falsch. Dies ist ein Dokumentationsbefund am Plan, kein Mangel an der ausgelieferten Regel.

## Verdikt

**Merge-blockierend: ja, H-1.** Der Ersatzweg-Mechanismus (Skript, Makefile-Ziele, Träger, Regel) ist im
gelieferten Umfang solide: die acht vertraglich benannten Mutationen und zwei eigene, nicht genannte Mutationen
(doppeltes `--load`) fingen den Tabellentest korrekt rot; `make gates`, `make suchlauf-nachmessen` und `make
docs-check` bestätigen unabhängig grün. Der Befund H-1 betrifft eine **echte, real demonstrierte Sicherheitslücke
in der Testabdeckung** des zentralen Schutzmechanismus dieses Ziels (der Arbeitsbaum-Schutz aus `AGENTS.md`
§3.1): kein Testfall bindet die `realpath`-Auflösung von `SRC` an ihre Eingabeseite, und das genau von Plan-Risiko
1 benannte Szenario (relativer Pfad umgeht den Wurzel-Schutz) ist real reproduzierbar, sobald diese eine Zeile
fehlt oder anders implementiert wird — was der aktuelle Tabellentest nicht bemerken würde.

**Fixrunde nötig:** Implementer ergänzt mindestens einen Testfall mit relativem `SRC` (z. B. `SRC=.` bzw.
`SRC=sub` aus dem Wegwerf-Repo heraus aufgerufen) und — soweit der Testaufbau es zulässt (Plan-Risiko 1 nennt
das als offene Frage) — einen Symlink-Fall, sodass ein Entfernen/Fehlverhalten der `realpath`-Zeile den
Tabellentest rot färbt; danach `make test-image-mutation` erneut mit der in H-1 beschriebenen Mutation gegen die
Kopie laufen lassen und das neue Rot im Bericht/Vertrag nachtragen. M-1 (MEDIUM) wird im selben Zug behoben:
Plan-Zeile 366 (Träger-Tabelle) korrigieren oder auf die tatsächliche, exceptionslose Fassung von §3.15 verweisen;
optional Zeile 98–106 („Entwurf“) nachziehen oder explizit als überholten Entwurf kennzeichnen. I-1 ist ein
Hinweis für die Closure-Notiz (Risiko 8 „entfallen“ statt „bei Closure zu bewerten“), kein eigener Fixbedarf. I-2
ist Kontext ohne erwartete Aktion.

Die DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" bleibt **offen** (kein Nachzug in
diesem Commit) — die Fixrunde an H-1/M-1 führt regulär über Schritt 21 des Implementer-Ablaufs zurück zum
Reviewer (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde greift hier nicht, da H-1 offen
bleibt).

**Übergabe:** H-1 und M-1 an den Implementer (Testfall relativer `SRC`/Symlink ergänzen; Plan-Zeilen 98–106/366
korrigieren); I-1 an den Planner für die künftige Closure-Notiz.
