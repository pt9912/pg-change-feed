# Review-Report: slice-090 — **Delta** (Fixrunde `c2bc08f`) · 2026-09-16

**Review-Art:** Delta-Review eines Fix-Commit — geprüft gegen **Plan und
Entscheidungen** sowie gegen die Findings der ersten Runde
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier, Modul 11),
§6-Risiko-Ausgänge, Beobachtungs-Register und die drei Paarungen
(Planner-Closure) sind **nicht** Gegenstand dieses Reports.

**Gegenstand:** `c2bc08f` (`fix(harness): Review-Fixrunde slice-090 …`),
Parent `0122eb1`, ein Commit. `git diff --name-status c2bc08f^..c2bc08f`:
`M .github/workflows/ci.yml`, `M docs/plan/planning/in-progress/slice-090-sync-gate-protobuf.md`,
`M harness/README.md`, `A harness/sensors/generated-sync.md` (111 Z.),
`M tools/harness/generated-sync.sh` — `git diff --stat`: **5 Dateien, +136/−11**.
Kein Produktionscode, `Makefile`/`Dockerfile` unberührt, Arbeitsbaum sauber. Der
Gate-Skript-Stand ist an `HEAD` byte-gleich dem an `c2bc08f`
(`git diff --stat c2bc08f..HEAD -- tools/harness/generated-sync.sh
harness/sensors/generated-sync.md` → leer), die Mutations-Messungen laufen daher
am Prüfgegenstand.

**Skill:** `.harness/skills/reviewer.md` @ `HEAD` — die fünf repo-spezifischen
HIGH-Unterpunkte gehören zum Prüfraster · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/reviews/review-slice-090.md` (F-1…F-5, die Findings der ersten Runde)
- `ADR-0084` §Entscheidung **vollständig** — insbesondere die Zählung
  „**Vier Festlegungen**" (`:126`) und ihre Selbstverweise (`:215`, `:267`,
  `:297`, `:301`), `ADR-0060`, `ADR-0044`, `ADR-0045`
- `AGENTS.md` §3.1, §3.3, §3.6, §3.7, §3.9, §3.11, §3.12, §4, §6 ·
  `harness/README.md` §Sensors **samt** Kommentar-Block (Z. 60–109) ·
  `harness/conventions.md` (MR-000)
- Form-Vorbilder `harness/sensors/coverage-gate.md` (180 Z.) und
  `harness/sensors/a-check.md` (60 Z.), Nachbar `docs-check.md`,
  `baseline-verify.md`
- **Plan-Stand VOR der Fixrunde** (`git diff --stat 9029c05 0122eb1 -- <Plan>`
  → leer) — gegen diesen Stand wurde verglichen, nicht gegen den nachgezogenen
- Sibling-Kontext zur Schnittmenge: `docs/reviews/verify-slice-090.md`
  (Kopfzeilen und Messtabelle gelesen, um Doppelführung zu vermeiden; die
  eigenen Messungen dieses Reports sind **unabhängig** gefahren)

---

## Findings

### D-1 — Die Bindung nennt die falsche Festlegung — `ADR-0084` zählt anders

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz A/B (übernommener Wert aus einem Bericht
  ohne Nachmessen) · Source Precedence Rang 4 (`ADR-*` ist der Beleg)
- `pfad`: `harness/sensors/generated-sync.md:108-109` (gegen
  `docs/plan/adr/0084-sync-gate-fuer-generierte-artefakte.md:126`, `:215`,
  `:267`, `:297`, `:301`)
- `befund`: Das Sensor-Dokument schreibt „(Festlegung 1 = der Baum wird nicht
  geschrieben; **Festlegung 2 = der Befund nennt Datei und Zeile**)". Gemessen
  am ADR: es führt „**Vier Festlegungen**" — **Festlegung 1** ist „Der
  Protobuf-Code bekommt ein Gate" und trägt die zwei **Bedingungen**
  (Temp-Verzeichnis / Datei-und-Zeile) als Unterpunkte; **Festlegung 2** ist
  „Die E2E-Abdeckungstabelle bekommt heute kein Gate" und wird vom ADR selbst so
  zitiert (`:215` „(Festlegung 2), ihr Schnitt ist Planner-Arbeit", `:267`,
  `:297`). „Datei und Zeile" ist damit **Festlegung 1, zweite Bedingung** —
  genau die Form, die der erste Report in F-2 als `quelle` richtig nennt
  (`review-slice-090.md:75`). Die Nummerierung des Sensor-Dokuments ist aus der
  *Kopfzeile* von `review-slice-090.md:22` übernommen (dort ebenso
  „Festlegung 1 = … Festlegung 2 = …", und ebenso in `verify-slice-090.md:6`) —
  ein Bericht ist **Lauf-Beleg**, kein Zitat-Träger; übernommen wurde er hier
  ungeprüft in einen **stehenden** Träger.
- `verifizierbar`: nein — kein Gate prüft Abschnitts-Nummern innerhalb eines
  zitierten Dokuments; `make docs-check` prüft Link-Ziele, nicht Zitat-Nummern
  (gemessen: `738 Datei(en), 0 Befund(e)`)
- `klasse`: „Zitat nennt die falsche Stelle — Nummerierung aus einem Bericht
  übernommen" (Nachbar-Klasse von
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, nicht dieselbe: Zitat
  statt Größe)

### D-2 — „… der des Kontext-Diffs eine Zeile voraus" hält gegen die Messung nicht

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz B (Tatsachenbehauptung ohne Beleg-Anker) ·
  §3.7 (eine Aussage über die Ausgabe beschreibt, was ist) · `ADR-0084`
  Festlegung 1, zweite Bedingung (der Befund nennt die Zeile)
- `pfad`: `harness/sensors/generated-sync.md:58-60` (gegen `:53-57`)
- `befund`: Der Absatz leitet `<N>` korrekt aus dem `diff -U0`-Hunk-Kopf her
  (F-2-Fix, siehe Negativbefunde) und schließt dann: „Die Nummerierung ist damit
  der des **Kontext**-Diffs **eine Zeile voraus**, dessen Hunk-Kopf unverändert
  bleibt". Nachgemessen im Einfügungs-Zweig (`-N,0`), also in dem Zweig, den der
  Satz beschreibt: der Abstand ist **1** nur, wenn die Abweichung auf **Zeile
  2** liegt; bei Zeile 4/5/40/226 ist er **3**, bei Zeile 1 ist er **0**.
  Auch der Nachbarsatz „systematisch ein bis drei Zeilen darunter" (`:54-55`)
  trifft die Zeile-1-Lage nicht (Abstand 0). Der Satz nennt keinen Beleg-Anker
  und ist auch nicht als „erwartet" formuliert; er ist die einzige Aussage des
  Dokuments, die einer Messung widerspricht.
- `verifizierbar`: ja — eine Mutation/eine gelöschte Zeile je Lage, dann Befund
  aus `generated-sync.sh` gegen den Kontext-Diff-Kopf desselben Logs stellen
- `klasse`: „Nummerierungs-Relation im Träger driftet gegen die Messung"

### D-3 — Die Exit-Tabelle definiert Code 1 nicht überschneidungsfrei

- `kategorie`: LOW
- `quelle`: Form-Vorbild `harness/sensors/a-check.md:42` („1 = mindestens ein
  Befund"), `AGENTS.md` §3.7
- `pfad`: `harness/sensors/generated-sync.md:39-43` (gegen `:99-103`)
- `befund`: Die Zeile „1 | Abweichung: Inhalt, fehlende Datei im Baum, oder ein
  gekennzeichnetes Erzeugnis ohne Quelle" liest sich als Definition des Codes.
  Gemessen liefern **drei** weitere Lagen ebenfalls `1`: ein fehlgeschlagener
  Stufen-Build (`EC=1`, Meldung des Builds), ein falscher Modulpfad (`EC=1`,
  `--go_out: …: generated file does not match prefix "example.com/other"`) und
  ein Quellverzeichnis außerhalb des Baums (`EC=1`, `Could not make proto path
  relative`). „§Sperren" nennt den Build-Fall als „Exit ungleich 0" (`:100-103`)
  — die Tabelle daneben führt `1` weiter als Abweichung. Ein Konsument, der `1`
  als „Drift" liest, deutet einen Build-/Generator-Fehler als Inhaltsbefund.
- `verifizierbar`: ja — Build-Stufe absichtlich brechen bzw.
  `GENERATED_SYNC_MODULE` verstellen und den Exit-Code lesen
- `klasse`: „Exit-Tabelle definiert einen Code nicht überschneidungsfrei"

### D-4 — Der Quell-Override nennt seine Vorbedingung nicht: der Pfad muss im Baum liegen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Grenze) — die Sektion des Dokuments, die genau
  dafür existiert (`:71-87`), führt sie nicht · `review-slice-090.md` F-3
- `pfad`: `harness/sensors/generated-sync.md:67` (Override-Tabelle) und `:78-79`
  (§Grenze 2)
- `befund`: Die Zeile sagt „ersetzt die **verglichene Quelle** (Default `proto`)"
  und §Grenze 2 „Mit gesetztem `GENERATED_SYNC_SOURCE_DIR` vergleicht der Lauf
  eine andere Quelle". Gemessen gilt der Ersatz **nur für Pfade im Arbeitsbaum**:
  der Baum hängt als `-v "$repo_root":/src:ro` im Container, ein Verzeichnis
  außerhalb ist dort unsichtbar — `GENERATED_SYNC_SOURCE_DIR` auf ein auf dem
  Host existierendes Verzeichnis außerhalb → `EC=1`,
  `… : warning: directory does not exist` / `Could not make proto path
  relative: …: No such file or directory`. Die Richtung ist die sichere (rot),
  aber die dokumentierte Übergabe trägt ihre Bedingung nicht. **Schnittmenge
  offen benannt:** `verify-slice-090.md` V-4 führt dieselbe Beobachtung — eine
  Fundstelle genügt, nicht zwei.
- `verifizierbar`: ja — Override auf ein Verzeichnis außerhalb des Baums
- `klasse`: „Dokumentierte Übergabe ohne ihre Vorbedingung"

### D-5 — Der Skriptkopf ist mitten im Absatz neu umbrochen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Form des Kommentars, stilistisch)
- `pfad`: `tools/harness/generated-sync.sh:26-28`
- `befund`: Die Teilersetzung hat den Absatz nur bis „beginnt." neu umbrochen;
  die Folgezeile ist stehen geblieben und reißt den Satz auf eine
  40-Zeichen-Zeile: `# beginnt. Die gepinnte Stufe laeuft ohne` /
  `` # `--no-cache-filter`: ihren Layer-Cache … `` . Der Satz bricht **nicht**
  ab (die Aussage bleibt lesbar) — es ist ein Umbruch-Artefakt, keine
  Aussage-Drift. `:20-25` ist inhaltlich korrekt (siehe Negativbefunde).
- `verifizierbar`: nein (kein Sensor für Kommentar-Umbruch)
- `klasse`: „Teilersetzung lässt den Rest des Absatzes stehen"

### D-6 — §Grenze führt eine Sperre; die Sektion hat einen Sperren-Abschnitt daneben

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Form) · Form-Vorbilder: `coverage-gate.md`
  §Grenze (nur Deckungsgrenzen), `a-check.md` §Sperren (Image-Digest)
- `pfad`: `harness/sensors/generated-sync.md:88-94` (gegen `:99-103`)
- `befund`: Punkt 4 „Der Generator-Build braucht Netz auf kaltem Layer-Cache"
  beschreibt **keine** Deckungsgrenze des Grün („was das Grün nicht abdeckt",
  `:71`), sondern eine Vorbedingung des Laufs — dieselbe Gegenstands-Klasse wie
  „§Sperren" (`:101-103`: Docker-Daemon, `Dockerfile` mit Stufe `proto`). Die
  Punkte 1, 2, 3 und 5 sind Deckungsgrenzen und passen; Punkt 4 mischt die
  Klasse. Der Inhalt selbst ist richtig (siehe Negativbefunde, `#6–#9 CACHED`).
- `verifizierbar`: nein
- `klasse`: „Grenze-Abschnitt trägt eine Sperre"

### D-7 — Die nachgezogene Plan-Zeile nennt eine Zelle, die der Diff nicht ändert

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 / §3.12 Instanz B (eine Aussage über den eigenen
  Diff) · `harness/README.md` §Sensors Kommentar-Block `:98`
- `pfad`: `docs/plan/planning/in-progress/slice-090-sync-gate-protobuf.md:141`
  (gegen `harness/README.md:118`)
- `befund`: Die in dieser Fixrunde neu geschriebene Zeile der Änderungstabelle
  sagt: „die **Target-Zelle wird zum Link** auf
  `harness/sensors/generated-sync.md`". Gemessen steht der Link in der
  **Bindung**-Spalte (`:118`, drittes Feld); die Target-Zelle bleibt
  `` `make generated-sync` `` ohne Link. Alle **vier** Nachbarzeilen mit
  Sensor-Dokument tun dasselbe (`baseline-verify` `:113`, `docs-check` `:114`,
  `a-check` `:115`, `coverage-gate` `:117`) — die *Umsetzung* folgt der Praxis,
  der *Satz* wiederholt die Formulierung des Kommentar-Blocks (`:98`, seit
  `40c8c43`, außerhalb dieses Diffs). Damit hält die neue Plan-Aussage nicht;
  sie ist die Art Aussage, die der Verifier im Plan-vs-Code-Diff prüft.
- `verifizierbar`: ja — Zeile `:118` gegen den Diff des Fix-Commit
- `klasse`: „Plan-Aussage über den eigenen Diff benennt die falsche Stelle"

### D-8 — Die neue Vertragszelle assertiert einen Ist-Zustand

- `kategorie`: INFO
- `quelle`: `harness/README.md` §Sensors Kommentar-Block `:66-67`
  („Vertrag: was prüft das Gate (was wäre verletzt, wenn es rot wird)") und
  `:72-73` („Lauf-Wahrheit pro Commit liegt in CI, nicht hier")
- `pfad`: `harness/README.md:118`
- `befund`: Die Zelle beginnt mit „Der committete Protobuf-/gRPC-Code **ist**
  byte-gleich der Ausgabe des gepinnten Generators …". Die fünf Nachbarzellen
  benennen ihren Gegenstand (Verb oder Substantiv: „prüft", „verifiziert",
  „Go-Test-Coverage gegen `THRESHOLD`"); diese Zelle formuliert eine
  Zustandsaussage über den Baum. Beide Lesarten sind vertretbar — als Invariante
  liest sie sich wie „was wäre verletzt, wenn es rot wird"; als Ist-Satz
  widerspräche sie dem Verbot des Lauf-Status in dieser Tabelle. **Kein Drift
  gemessen:** `make gates` ist grün (`generated-sync: OK`), die Zelle ist an
  diesem Commit wahr.
- `verifizierbar`: ja (der Ist-Stand), nein (die Form)
- `klasse`: „Vertragszelle assertiert einen Ist-Zustand statt der Verletzung"

### D-9 — Eine dritte, schwächere Fassung der Gate-Liste ist stehen geblieben (Träger außerhalb des Diffs)

- `kategorie`: INFO
- `quelle`: `.harness/skills/reviewer.md` HIGH-Unterpunkt „Zahl im Träger …",
  *Träger außerhalb des Diffs bleiben INFO* · Klasse von
  `review-slice-090.md` F-4
- `pfad`: `.claude/agents/implementer.md:42`, `.claude/agents/verifier.md:38`
- `befund`: `ci.yml` ist auf sechs Namen nachgezogen (gemessen: deckungsgleich
  mit `GATE_CHECKS`). Zwei Rollen-Briefings beschreiben `make gates` weiterhin
  als „baseline-verify + d-check über `.d-check.yml`" — zwei von sechs Zielen,
  ohne Zähl-Aussage (schwächere Form als F-4, das „deckt … ab" sagte). Der
  Wortlaut „Deine repo-spezifischen Sensoren (neben `make gates`)" macht die
  Liste nicht ausdrücklich erschöpfend; als stehen gebliebene Zweitfassung
  derselben Aussage ist sie es faktisch.
- `verifizierbar`: ja — `make -p | grep '^GATE_CHECKS'` gegen die beiden Zeilen
- `klasse`: `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Träger
  außerhalb des Diffs)

---

## Negativbefunde

- **geprüft, ohne Befund: die F-2-Korrektur trägt — gemessen in neun
  Stellungen.** Der Befund nennt genau die erste abweichende Stelle, geprüft
  gegen eine unabhängige Auswertung (Python-Vergleich der Zeilenlisten, fehlende
  Zeile = Unterschied): Mutation Z. 1/102/260 → 1/102/260, Löschung
  Z. 1/40/226-230 → 1/40/226, Anhang am Ende → 261, Einfügung in der Mitte → 51
  — **9× deckungsgleich**. Ebenso am realen Drift (`string probe_feld = 11` in
  der `.proto`): `FAIL … Zeile 47` und vier Hunks; der erste Kontext-Hunk-Kopf
  ist `@@ -44,6 +44,7 @@`, die eingefügte Zeile ist die vierte darin = 47 — die
  Zahl ist die genaue Stelle.
- **geprüft, ohne Befund: die `N+1`-Regel des Skriptkopfs beschreibt den
  Code.** Einfügungs-Hunks (`-N,0`) gemessen: `@@ -0,0 +1 @@` → 1,
  `@@ -1,0 +2 @@` → 2, `@@ -2,0 +3 @@` → 3, `@@ -39,0 +40 @@` → 40,
  `@@ -225,0 +226 @@` → 226 — je `N+1`. Löschungs- und Änderungs-Hunks
  (Zähler > 0) nennen den Start unverändert (u. a. `@@ -261 +260,0 @@` → 261).
  Die Grenzen halten: leere committete Datei (`@@ -0,0 +1,260 @@`) → `Zeile 1`;
  fehlender Zeilenumbruch am Ende → `Zeile 260` mit
  `\ Kein Zeilenumbruch am Dateiende.`; kein `@@`-Kopf ist nicht erreichbar,
  weil `cmp -s` vorher die Abweichung feststellt.
- **geprüft, ohne Befund: die F-1-Korrektur.** `harness/README.md:118` trägt
  **227 Zeichen** (nachgemessen, wie die Commit-Message behauptet; vorher 440),
  liegt damit in der Reihe `a-check` 123 / `docs-check` 101 / `baseline-verify`
  91 / `coverage-gate` 286 / `commit-traceability` 295; die Zelle ist ein Satz
  und trägt keinen Absatz-Überhang mehr. `harness/sensors/generated-sync.md`
  existiert (111 Z., gegen `coverage-gate.md` 180, `a-check.md` 60), ist aus der
  Vertrags-Zelle erreichbar, und `make docs-check` ist grün (**738 Datei(en),
  0 Befund(e)**) — die neuen Verweise (`../../docs/plan/adr/0084-…`,
  `../../spec/lastenheft.md`) lösen auf.
- **geprüft, ohne Befund: die übernommenen Aussagen des Sensor-Dokuments
  halten.** Nachgemessen: `GENERATED_SYNC_SOURCE_DIR` auf eine unveränderte
  Quelle **in gleicher relativer Lage** im Baum → `EC=0`, Erfolgstext
  `Quelle: proto-q/cdc/stream/v1/changestream.proto` (§Grenze 2,
  Override-Tabelle); `TMPDIR` innerhalb des Baums → `EC=1` mit zwei Rot-Zeilen
  über die **eigenen** Temp-Dateien, `trap` räumt, `git status --porcelain`
  danach leer (§Grenze 3); Defaults `pg-change-feed:proto-sync`, `proto`,
  `module`-Zeile aus `go.mod`, Aufrufer-`uid:gid` gegen `:38-48` des Skripts;
  Pins `protobuf-dev=31.1-r1`, `protoc-gen-go@v1.36.12`,
  `protoc-gen-go-grpc@v1.6.2` gegen `Dockerfile:33-36`; `:ro`-Bind-Mount und
  `--network none` gegen `:66-72`; „in-place, kein Prüf-Schritt" gegen
  `Makefile:104-109` (`-v "$(CURDIR)":/src` ohne `:ro`, `--go_out=.`). Die
  beiden Zahlen von §Grenze 4 (Cache-Lauf) sind durch `#6 CACHED` … `#9 CACHED`
  im heutigen Build gestützt; der kalte Build ist in `verify-slice-090.md` #17
  unabhängig reproduziert.
- **geprüft, ohne Befund: `ci.yml` — nur zwei Textstellen, keine Semantik.**
  Die Kommentar-freie Fassung gediffed (`grep -v '^\s*#'`, `c2bc08f^` gegen
  `c2bc08f`) zeigt **genau eine** geänderte Zeile — den `name:` des Schritts
  (Z. 66); `run: make gates` ist unberührt, ebenso Stufen, Matrix,
  Abhängigkeiten und jede `uses:`-Zeile (§3.8). Kein Gate und kein anderer
  Träger referenziert den Schrittnamen (`grep -rn 'Gates (baseline-verify'` →
  nur `ci.yml:66` und die README-Zeile `:119`). Die Commit-Message-Aussage hält.
- **geprüft, ohne Befund: §3.7 — Kommentar-Klassen.** Der Skriptkopf trägt
  Zusage (Befund/Zeile), Kopplung (Bind-Mount, Aufrufer-uid,
  Modul-Layout-Relation), Abgrenzung („Nicht Gegenstand: dass der erzeugte Code
  kompiliert …") und Rang-Zeiger (`ADR-0084` Festlegung 1, `ADR-0060`); keine
  Slice-/Wellen-Chronik im Mechanik-Teil, kein Vorher/Nachher. Die neue Datei
  führt den Herkunfts-Anker in der zulässigen Form (`seit slice-090`, `:111`);
  „Gemessen in diesem Zug" (`:91`) löst über dieselbe Zeile auf slice-090 auf.
  Als Fund bleibt D-5 (Umbruch), nicht die Klasse.
- **geprüft, ohne Befund: §3.11 und §3.12 Instanz A.** Keine host-lokalen
  absoluten Pfade in den hinzugefügten Zeilen (Präfix-Grep leer;
  `make docs-check` grün, `hostpaths`-Modul eingeschlossen). Das Delta führt
  **keine** neue Messzahl in einem Träger ein; die vorhandenen Zahlen des neuen
  Dokuments (§Grenze 4) tragen Ursprung, Lauf-Marker und die gedruckten Zeilen.
  Keine Betreiber-Oberfläche (`CDC_*`, `cdc.*`, Adresse/Endpunkt) im Delta,
  `docs/user/` unberührt — die beiden Handbuch-HIGH-Regeln lösen nicht aus.
- **geprüft, ohne Befund: Bindung, §3.6, §3.3, §3.1.**
  `make -p | grep '^GATE_CHECKS'` → sechs Einträge (`baseline-verify
  coverage-gate docs-check commit-traceability generated-sync a-check`); die
  neue Doku ist kein Gate und senkt keine Schwelle (kein
  `THRESHOLD`/`matrix:`/`timeout-minutes` im Diff); kein `git mv`; alles läuft
  über `make`/Docker.
- **geprüft, ohne Befund: die neue Plan-Zeile zu `ci.yml`.** Sie sagt „**Nur
  zwei Textstellen** … Keine Stufe, keine Matrix, keine Abhängigkeit, kein
  Semantik-Wechsel" — deckungsgleich mit der Messung; die Aussage über
  `generated-sync` als sechstem Namen ist deckungsgleich mit `GATE_CHECKS`. Im
  **vor** der Fixrunde gültigen Plan war der `ci.yml`-Punkt **nicht** aufgeführt
  (er ist eine Adresse der stehen gebliebenen Zweitfassung, F-4).
- **geprüft, ohne Befund: Hygiene des Delta-Range.** Keine Lauf-Artefakte,
  keine Vorlagen-Reste; Betreff nennt `ADR-0084` und keine Struktur-ID;
  Arbeitsbaum nach allen Läufen sauber; `make commit-traceability` grün
  (`5 Commit(s)`, Betreffs ohne Struktur-ID) — der Delta-Commit selbst trägt
  seine Kennung.
- **geprüft, ohne Befund: `make gates` real — Exit 0** (Log in Datei, Exit-Code
  separat in eigener Datei gelesen und in einem eigenen Schritt ausgewertet,
  §3.9): `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK —
  Coverage 74.70% erfüllt Schwelle 70%` · `d-check: 738 Datei(en) geprüft,
  0 Befund(e)` · `commit-traceability: OK` · `generated-sync: OK — byte-gleich`
  · `a-check … gesamt: 0 Befund(e)`; danach `git status --porcelain` leer.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-status c2bc08f^..c2bc08f` | **0** | 5 Dateien (1 `A`, 4 `M`), +136/−11 |
| `make generated-sync` (Original-Baum) | **0** | `OK — byte-gleich`; Status leer; 1,17 s |
| `docker build --target proto …` (2. Lauf) | **0** | `#6`–`#9 CACHED` → Cache-Aussage §Grenze 4 |
| Klon: 9 Mutationen (Ändern/Löschen/Anhängen, Z. 1…260) | **1 (9×)** | Befund == unabhängige Auswertung |
| Klon: Einfügungs-Zweig, gelöschte Z. 1/2/3/4/5/40/226 | **1 (7×)** | Befund 1/2/3/4/5/40/226; Kontext-Kopf-Abstand **0/1/2/3/3/3/3** → **D-2** |
| Klon: Änderungs-Leiter Z. 1/2/3/4/5/50 | **1 (6×)** | Abstand **0/1/2/3/3/3** → **D-2** |
| Klon: committete Datei leer | **1** | `(… Zeile 1)` — `-0,0`-Grenze hält |
| Klon: kein Zeilenumbruch am Ende | **1** | `(… Zeile 260)` — `\ Kein Zeilenumbruch …` |
| Klon: realer Drift (`probe_feld = 11`) | **1** | `Zeile 47` + 4 Hunks — Klasse real gefangen |
| Klon: Override auf unveränderte Quelle, gleiche Lage | **0** | `Quelle: proto-q/…` — Grün-Pfad real |
| Klon: Override auf Verzeichnis **außerhalb** des Baums | **1** | `Could not make proto path relative: …: No such file or directory` → **D-4** |
| Klon: `TMPDIR` im Baum | **1** | zwei Rot-Zeilen über die eigenen Temp-Dateien; Status danach leer |
| Klon: `GENERATED_SYNC_MODULE=example.com/other` | **1** | `generated file does not match prefix …` → **D-3** |
| Klon: Stufe `proto` absichtlich gebrochen | **1** | Build-Fehler, **keine** Skript-`FAIL`-Zeile → **D-3** |
| `make -p \| grep '^GATE_CHECKS'` | **0** | 6 Einträge |
| Zell-Länge `harness/README.md:118` | **0** | Vertrag 227; Link in der Bindung-Spalte, Target unverlinkt → **D-7/D-8** |
| `ci.yml` Kommentar-frei gediffed (`c2bc08f^` vs `c2bc08f`) | **1** | genau **eine** Zeile: der Schrittname |
| `make docs-check` | **0** | `738 Datei(en), 0 Befund(e)` |
| `make gates` (Log in Datei, Exit separat gelesen) | **0** | sechs Checks; danach Status leer |
| `git diff --stat c2bc08f..HEAD -- <Skript> <Doku>` | **0** | leer — Prüfgegenstand = `HEAD`-Stand |

Der Gate-Lauf und seine Auswertung sind **zwei Schritte**: `make gates` lief
ungepiped in eine Log-Datei, der Exit-Code in eine eigene Datei, gelesen in
einem eigenen Werkzeug-Aufruf; kein Folgekommando hing an Pipe oder Wrapper
(§3.9).

## Antwort auf die Schwerpunkte

**(1) Das neue Sensor-Dokument — Form trägt, zwei Aussagen nicht.** Aufbau und
Zuschnitt folgen `coverage-gate.md` (Vertrag · Ausgabe und Ausgänge · Grenze ·
Sperren · Bindung), die vier Überhang-Arten des Kommentar-Blocks
(Deckungsgrenze, Ausgabe-Bedeutung, Exit-Codes, Sperren) sind abgedeckt, alle
Zahlen und Mechanik-Aussagen wurden nachgemessen und halten — **mit zwei
Ausnahmen**, beide in derselben Datei: die Festlegungs-Nummerierung (D-1) und
der Voraus-Satz (D-2). Die vom Skriptkopf getragenen Mechanik-Aussagen stehen im
Dokument erneut (wie in `coverage-gate.md` §Vertrag üblich) — kein Fund.

**(2) Die F-2-Korrektur — trägt, gemessen.** Die `diff -U0`-Herleitung ist
richtig, auch für Einfügungs- und Lösch-Hunks; die `N+1`-Regel des Skriptkopfs
beschreibt den Code (neun Stellungen deckungsgleich, plus realer Drift). Der
einzige Punkt, an dem eine Aussage nicht trägt, ist der **zusätzliche** Satz im
Sensor-Dokument, nicht der Skriptkopf (D-2).

**(3) §3.7 — die Kommentare tragen Klassen.** Zusage, Kopplung, Abgrenzung,
Rang-Zeiger je Stelle; keine Chronik, keine abwesende Sprache. Ein
Umbruch-Artefakt aus der Teilersetzung (D-5, stilistisch).

**(4) `ci.yml` — zwei Textstellen, keine Semantik.** Kommentar-frei gediffed
bleibt genau der Schrittname; Stufe, Matrix, Abhängigkeiten und `uses:` sind
unberührt, kein Träger referenziert den Namen. Die Commit-Message-Aussage hält.

**(5) Ohne Zirkularität geprüft.** Verglichen wurde gegen den Plan-Stand **vor**
der Fixrunde. Die neue `ci.yml`-Zeile und die neue
`harness/sensors/generated-sync.md`-Zeile der Änderungstabelle beschreiben den
Diff korrekt; die neue README-Zeile tut es nicht (D-7).

**Was nicht geprüft wurde:** DoD/LP1–LP3 (Verifier), die §6-Risiko-Ausgänge und
die drei Paarungen (Planner-Closure), der Register-Zähler (keine Register-Datei
im Delta), und die ersten fünf Nachbar-Reviews (Klassen nur über den
Anschluss-Report `review-slice-090.md` geführt). Der reale CI-Lauf des neuen
Gates ist mit `AGENTS.md` §3.10 weiterhin offen — dieser Diff ändert den
Workflow nur im `name:` und in Kommentarzeilen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Zitat nennt die falsche Stelle — Nummerierung
aus einem Bericht übernommen" (D-1) · „Nummerierungs-Relation im Träger driftet
gegen die Messung" (D-2) · „Exit-Tabelle definiert einen Code nicht
überschneidungsfrei" (D-3) · „Dokumentierte Übergabe ohne ihre Vorbedingung"
(D-4) · „Teilersetzung lässt den Rest des Absatzes stehen" (D-5) ·
„Grenze-Abschnitt trägt eine Sperre" (D-6) · „Plan-Aussage über den eigenen Diff
benennt die falsche Stelle" (D-7) · „Vertragszelle assertiert einen Ist-Zustand"
(D-8) · `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (D-9, Träger
außerhalb des Diffs).

**Register-Hinweis (Modul 6):** D-1 und D-2 sind **ein Vorgang** (`slice-090`)
— sie zählen den Zähler ihrer Klassen je **einmal**, auch wenn sie in drei
Artefakten desselben Zuges sichtbar sind (`review-slice-090.md:22`,
`verify-slice-090.md:6`, Sensor-Dokument). Der Beleg gehört an `slice-090`; der
`state.md`-Ausgang bleibt Sache des Lese-Schritts der `welle-20`-Closure.

## Verdikt

**Merge-blockierend:** nein — **0 HIGH**.

**Trägt die Fixrunde?** **F-1 und F-2 tragen — gemessen, nicht gelesen.** Die
Zell-Länge ist von 440 auf 227 gefallen, das Sensor-Dokument existiert und ist
erreichbar, und die Zeilenherleitung nennt in neun Stellungen plus im realen
Drift-Fall die erste abweichende Stelle. Die Fixrunde führt aber **neuen Text in
einem neuen stehenden Träger** ein, und zwei Sätze darin halten nicht: die
Festlegungs-Nummerierung (D-1) und der Voraus-Satz (D-2). Beide sind
Ein-Satz-Korrekturen in `harness/sensors/generated-sync.md` — der Datei, die
diese Fixrunde selbst angelegt hat.

**Zweite Runde: ja, aber eine kurze.** Für D-1 und D-2 ist ein **Rückgabe-Pfeil
Reviewer → Implementer** nötig (MEDIUM, beide im Kern-Erzeugnis dieses Zuges);
D-3 bis D-7 sind Form-Findings in denselben zwei bis drei Dateien
(`harness/sensors/generated-sync.md`, `tools/harness/generated-sync.sh`, die
eine Plan-Zeile) und können im selben Lauf mitgehen. D-8 und D-9 sind INFO: D-8
ohne Drift (Form-Abwägung), D-9 ist eine Adresse außerhalb des Diffs
(Planner-/Briefing-Seite), nicht Schuld dieses Zugs.

**DoD-Häkchen „Review durchgeführt, Report unter `docs/reviews/` liegt vor":**
bleibt **offen** — die Fixrunde geht zurück an den Implementer
(§DoD-Checkbox-Nachzug: die Regel greift nur, wenn *keine* Fixrunde nötig ist).

**Übergabe:** Rückgabe-Pfeil Reviewer → Implementer mit D-1…D-7; die
Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser
Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation (Modul 11).
