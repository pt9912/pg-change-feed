# Slice harness-mutationsbild-und-verweigerte-aktion: Ein `make`-Ziel für Mutations-Images mit eigenem Tag (`make image-mutation`, `make image-mutation-rm`) und die Regel „Eine verweigerte Aktion wird gemeldet, nicht auf anderem Weg wiederholt“ in `AGENTS.md`, Reviewer-Skill und Implementer-Command

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — Harness-Querschnitt: der Slice trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre. Er berührt das Makefile,
ein neues Skript samt Tabellentest, `AGENTS.md`, den Reviewer-Skill und den
Implementer-Command, nicht den Runner und nicht `internal/`; eine technische Kante
zu einem Slice der [welle-transformationen](welle-transformationen.md) hat er
nicht. Die empfohlene Position steht in §4.

**Bezug:** [`AGENTS.md`](../../../../AGENTS.md) §3.1 (Docker-only, Absatz zur
Mutationsprobe auf einer Kopie), §3.6 (das Ziel und die Regel lockern kein Gate:
keine ADR nötig), §3.7, §3.9, §3.12 und §3.13;
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md) und
[`ADR-0103`](../../adr/0103-image-hash-lokal-statt-committet.md) (der Digest in
`harness/image-hash.txt` ist der Lauf-Beleg des letzten `make image`-Laufs, lokal,
nicht committet; `:dev` ist der geteilte Tag der Compose-Umgebung),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft der
Aussagen im Plan; Grenze: kein Sensor über Berichte und Verläufe);
[`review-slice-leerlauf-phase-last-in-stuecken`](../../../reviews/review-slice-leerlauf-phase-last-in-stuecken.md)
F-3 (LOW, Abschnitt „Docker-only / §3.1“ und „Einordnung der Fragen des Auftrags“);
Beobachtungs-Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion` (neu, Adresse dieses
Slice), Nachbarn `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet`,
`BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` und
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`.

**Herkunft der Entscheidungen:** zwei Nutzer-Entscheidungen der Sitzung
(2026-09-27), Wortlaut des Nutzers: zur Frage „welcher Weg gilt für Mutations-Images“:
„ein eigenes make-Ziel mit eigenem Tag“; zur Frage „ob ein verweigerter Aufruf im
Bericht genannt und vor einem Ersatzweg zurückgefragt wird“: „ja“. Die Fragen stehen
im Verdikt des Reviews (Abschnitt „Einordnung der Fragen des Auftrags“, Punkte 1 und 2).

**Berührte Spec-Stellen:** — (Harness-Werkzeug und Arbeitsregel; keine Spec-Stelle).

**Verantwortlich:** Implementer-Agent (Sitzung 2026-09-27).

**Autor:** Planner-Agent, Auftrag des Auftraggebers (zwei Nutzer-Entscheidungen).
**Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Rollenlauf, der eine Mutation an Produktionscode gegen ein laufendes System
belegt, baut das Image dafür über ein `make`-Ziel mit **eigenem Tag**
(`make image-mutation SRC=<Verzeichnis> TAG=<Tag>`), ohne `:dev` und
`harness/image-hash.txt` zu berühren; und die Regel „Eine von der Berechtigungsschicht
verweigerte Aktion wird im Bericht genannt und nicht auf anderem Weg wiederholt, ohne
dass der Auftraggeber gefragt ist“ steht in `AGENTS.md`, im Reviewer-Skill und im
Implementer-Command.

**Ausgangslage — Beleg-Anker je Aussage.**

- *Die Regel `image` (gelesen, Stand `84f60e6f`, `Makefile` Zeilen 55–56):* die lokale Fassung ohne
  `VERSION` ruft `docker buildx build --load --metadata-file harness/image-hash.raw -t
  ghcr.io/pt9912/pg-change-feed:dev .` auf und schreibt den Digest nach `harness/image-hash.txt`.
  Die Compose-Umgebung liest denselben Tag (`compose.yaml` Zeile 72,
  `image: ghcr.io/pt9912/pg-change-feed:dev`, gelesen). Wer mit `make image` ein Image aus einem
  mutierten Baum baut, überschreibt damit den geteilten Tag und den Lauf-Beleg.
- *Kein `make`-Ziel für ein Mutations-Image (gemessen, `suchlauf`-Block in §3, Stand `84f60e6f`):*
  weder ein Ziel- noch ein Tag-Name dafür steht im Baum; die Suchzeilen zählen null Treffer außerhalb
  von Reviews und Records.
- *Der Vorgang (**übernommen** aus dem Review, nicht am Classifier gemessen):* der
  Berechtigungs-Classifier verweigerte dem Implementer von `slice-leerlauf-phase-last-in-stuecken`
  `make image` mit einer mutierten Go-Datei („Modify Shared Resources“); der Implementer baute danach
  `docker buildx build --load -t pg-change-feed-mutd:tmp .` und band das Image über eine
  Compose-Override-Datei einer Scratchpad-Kopie des Runners ein. Der Reviewer stuft den Ersatzweg als
  LOW ein (Wortlaut-Abweichung von `AGENTS.md` §3.1, in der Auslegung des Repos gedeckt) und nennt
  offen, dass er die Verweigerung auf anderem Weg umging; ob der Baum eine Scratchpad-Kopie oder der
  Arbeitsbaum war, ist aus den Artefakten nicht ablesbar (Review F-3, Punkt 3).
- *Frühere Praxis (**übernommen** aus demselben Review, die Verifikations-Reports nicht selbst gelesen):*
  zwei Verifikationen bauten Mutations-Images mit `make image` aus einer `git archive`-Kopie im
  Scratchpad und stellten `:dev` danach am unmutierten Repo wieder her; dabei wird der geteilte Tag
  überschrieben und wieder gesetzt.
- *Berechtigungsliste (gelesen, `.claude/settings.json`):* `Bash(make:*)` steht auf der Allow-Liste,
  `docker buildx build` nicht (nur `docker manifest inspect` und `docker buildx imagetools inspect`).
  Was der Classifier für das neue Ziel entscheidet, ist nicht gemessen (§6 Risiko 3).
- *Träger der zweiten Regel (gemessen, §3 `suchlauf`-Block):* keine Datei außerhalb von Reviews und Records
  nennt eine Regel zum Umgang mit einer Verweigerung; die Treffer der Suchworte sind der Stop-Hook
  („verweigert den Abschluss“), ein Satz in `MR-003` („nennen diesen Ersatzweg“, Guard-Meldung) und
  Testnamen in `internal/`.

**Kernaussage der Wirkung des Ziels (hergeleitet aus dem Aufruf, erprobt in DoD 1):** das Ziel
liest genau ein Verzeichnis (die Kopie) und schreibt genau ein Image unter dem Namen
`pg-change-feed-mutation:<TAG>` — kein `--metadata-file`, kein `-t` auf `:dev`, kein `--push`, keine
Datei im Repo. Die Kopie liefert der Aufrufer (`git archive <Stand>` in ein Scratchpad-Verzeichnis, Mutation
mit Edit/Write auf der Kopie); das Ziel lehnt die Repo-Wurzel und jedes Verzeichnis darunter als `SRC`
ab.

**Kernaussage der Regel (Entwurf, den der Implementer im Ton der Nachbar-Regeln schreibt):** eine Aktion,
die die Berechtigungsschicht — Classifier, Permission-Prompt, PreToolUse-Guard — verweigert, führt der
Lauf **nicht auf einem anderen Weg zum selben Ziel** aus. Der Bericht nennt den verweigerten Aufruf, den
Wortlaut der Ablehnung und das Ziel; danach fragt der Lauf den Auftraggeber, bevor er fortfährt — ohne
Ausnahme für den Fall, dass die Ablehnung selbst einen Weg nennt (ausgeliefert in `AGENTS.md` §3.15, deckungsgleich
mit der Nutzer-Entscheidung „ja“ ohne im Wortlaut genannte Ausnahme). Ein versehentlicher Aufruf ohne Nutzen
für die Aufgabe (ein `python3 --version` neben einem Repo-Pfad) wird gestrichen und im Bericht genannt, nicht
wiederholt. Die Ablehnung gilt dem Aufruf, nicht dem Ziel
(Nachbar: `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet`): der Bericht nennt Aufruf, Pfad und
Wortlaut, nicht eine Verallgemeinerung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung an `make image`, am `Dockerfile` oder an `.dockerignore`.** Das neue Ziel ruft denselben
  Bau mit anderem Kontext und anderem Tag; der Bau-Kontext der Kopie ist der Kontext des Repos
  (`.dockerignore` reist mit `git archive`). `make image` bleibt der Lauf-Beleg-Weg
  ([`ADR-0044`](../../adr/0044-image-beleg-semantik.md), [`ADR-0103`](../../adr/0103-image-hash-lokal-statt-committet.md)).
- **Ein Bau aus dem Arbeitsbaum.** Das Ziel verweigert die Repo-Wurzel: eine Mutation am Arbeitsbaum
  ist genau das, was `AGENTS.md` §3.1 verbietet (Mutationsprobe auf einer Kopie); ein Ziel, das es zuließe,
  machte den Ersatzweg zum Regelweg.
- **Das Erzeugen der Kopie im Ziel.** Die Mutation geschieht auf der Kopie mit den Datei-Werkzeugen des Laufs;
  ein Ziel, das die Kopie selbst zieht, gäbe der Mutation keinen Ort. `git archive` und `tar` stehen in der
  Klasse „Host-Werkzeug ohne Installation“ (`AGENTS.md` §3.1) und stehen im Vertrag als Aufrufbeispiel.
- **Die Anbindung an die Compose-Umgebung.** Die Override-Datei des Runners (eine Scratchpad-Kopie des Runners
  mit `image: pg-change-feed-mutation:<TAG>`) ist Anwendungsbeispiel im Vertrag, kein Code: je Mutation
  unterscheidet sich der Runner-Bereich, den sie treibt.
- **Aufräumen über `prune`.** Das Aufräum-Ziel `make image-mutation-rm TAG=<Tag>` entfernt genau ein Image
  unter dem eigenen Namen; ein `prune`-Weg bliebe nicht auf die eigenen Images begrenzt (Docker-Regel des
  Repos: nie `prune`).
- **Multi-Arch und Push.** Das Ziel lädt lokal und einplattformig wie der `:dev`-Pfad (`--load`); ein Mutations-Image
  wird nie veröffentlicht.
- **Ein Sensor für die Regel „Verweigerte Aktion“.** Ein Verlauf von Aufrufen und ihren Ablehnungen liegt
  außerhalb des Repos, und der Classifier ist kein Teil des Repos; ein Werkzeugaufruf hinterlässt in einer
  Datei keine Signatur ([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md), Grenze). Die Regel wirkt
  durch Lesen (Reviewer, über den Bericht des Laufs), nicht durch einen Sensor.
- **`.claude/agents/*`, `.claude/settings.json`, `plan-welle.md`, `close-welle.md`.** Die Agenten-Definitionen sind
  entschieden ausgenommen (Auftrag); die Berechtigungen ändert kein Agent-Lauf. Die Regel steht in `AGENTS.md` wie
  die Ausführungsdisziplin von §3.9, die für jede Rolle gilt; der Implementer-Command bekommt sie, weil der
  Vorfall dort lag und dort die Docker-only-Zeile zu `AGENTS.md` §3.1 verweist. Ob die beiden anderen Commands
  einen Satz brauchen, entscheidet ein Review-Fund, nicht dieser Plan.
- **Die Nachbar-Beobachtung `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet`.** Sie betrifft die
  Genauigkeit der Meldung einer Werkzeug-Ablehnung (eine Verallgemeinerung auf den Zielpfad), nicht den Ersatzweg;
  ihr Eintrag bleibt getrennt, der Satz der Regel zu „Aufruf, Pfad, Wortlaut“ trägt ihn mit.
- **Neubewertung, ob `docker buildx build` direkt (ohne `make`) unter `AGENTS.md` §3.1 statthaft ist.** Die Auslegung
  des Repos („ein direktes `docker run` gepinnter Images ist statthaft, wo kein `make`-Ziel existiert“, Review F-3)
  bleibt; für Mutations-Images gilt ab diesem Slice das Ziel, für die Verweigerung die Regel. Eine Schärfung des
  Wortlauts von §3.1 wäre eine eigene Entscheidung des Auftraggebers.

## 2. Definition of Done

Jedes Kriterium trägt „Zu belegen durch:“; jede Aussage über eine Mutation ist eine
**Erwartung**, bis der Implementer sie gefahren hat (Stelle, Instanz, gesehene Farbe;
[`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).

- [x] **Liefer-Punkt 1 — das Ziel.** `tools/harness/image-mutation.sh` kennt zwei Verben.
      `build <SRC> <TAG>` baut aus dem Verzeichnis `SRC` das Image `pg-change-feed-mutation:<TAG>` mit
      genau `docker buildx build --load -t pg-change-feed-mutation:<TAG> <SRC>` (`SRC` als aufgelöster
      absoluter Pfad, letztes Argument; kein `--metadata-file`, kein `--push`, kein `--platform`, kein
      `--build-arg`); `rm <TAG>` entfernt genau dieses Image mit `docker rmi pg-change-feed-mutation:<TAG>`
      (kein `-f`, kein `prune`). Eingabefehler enden mit Exit 2 und Klartext **vor** jedem Docker-Aufruf:
      `SRC` oder `TAG` fehlt; `TAG` außerhalb von `[a-z0-9][a-z0-9_.-]{0,62}`; `TAG` ist `dev` oder `latest`;
      `SRC` ist kein Verzeichnis, trägt kein `Dockerfile` oder kein `go.mod`, ist die Repo-Wurzel oder liegt
      unter ihr (Vergleich der aufgelösten Pfade gegen `git rev-parse --show-toplevel`). Ein Docker-Fehler
      endet mit Exit 1 und einer Zeile, die den Exit-Code des Aufrufs nennt. Das Skript schreibt keine Datei
      im Repo; `harness/image-hash.txt`, `harness/image-hash.raw` und das Image `:dev` sind nach dem Lauf
      unverändert. `make image-mutation SRC=<Verzeichnis> TAG=<Tag>` und `make image-mutation-rm TAG=<Tag>`
      rufen es; ohne `SRC` oder `TAG` bricht das Ziel mit `$(error …)` ab, bevor ein Befehl läuft
      (Vorbild: `make suchlauf-nachmessen`); Hilfe-Zeilen nennen Argumente und „kein Gate“; `.PHONY` trägt
      alle drei Ziele. *Zu belegen durch:* `make test-image-mutation` (Tabellentest, Stub-`docker`, netzlos,
      Vorbild `make test-fmt-check`/`make test-suchlauf-nachmessen`) mit den Fällen: Aufruf mit gültigen
      Argumenten (Docker-Argumente genau wie oben, `SRC` mit Leerzeichen), `harness/image-hash.txt` und
      `.raw` im Wegwerf-Repo unverändert, `TAG` leer/`dev`/`latest`/`Foo`/`a b`/`-x`/`a:b`/`../x`, `SRC`
      fehlt/kein Verzeichnis/ohne `Dockerfile`/ohne `go.mod`/Repo-Wurzel/Verzeichnis unter der Wurzel,
      Exit-Weitergabe des Docker-Fehlers (Exit 1), `rm` mit genau einem `rmi`-Argument, Make-Ebene ohne
      `SRC` bzw. `TAG` (kein Docker-Aufruf); je Zusage die Mutation ihrer Eingabeseite am Skript (Kopie,
      `TOOL=<Kopie>`), gesehenes Rot je Zeile im Bericht — *erwartet, zu erproben*: `-t` auf den Tag
      `ghcr.io/pt9912/pg-change-feed:dev` → Fall „Docker-Argument `-t`“ rot; `--metadata-file`
      angehängt oder eine Zeile, die `harness/image-hash.txt` schreibt → Fall „Hash-Datei unverändert“
      rot; Ablehnung von `dev` entfernt → Fall `TAG=dev` rot; Zeichenklasse gelockert → Fall `Foo`/`a:b`
      rot; Repo-Wurzel-Vergleich entfernt → Fall Wurzel rot; `SRC`-Prüfung auf `Dockerfile` entfernt →
      Fall ohne `Dockerfile` rot; Exit des Docker-Aufrufs verschluckt → Fall Docker-Fehler rot; `-f` an
      `rmi` → Fall `rm` rot. Menge der Erprobung: die Fälle des Tabellentests. **Realer Lauf** (einmal,
      kein Teil des Tabellentests): `git archive HEAD` in ein Scratchpad-Verzeichnis, eine Go-Zeile dort
      mit Edit ändern, `make image-mutation SRC=<Kopie> TAG=<Tag>` — das Image existiert
      (`docker image inspect`), seine Image-ID weicht von der Image-ID von `:dev` ab (Erwartung: die Mutation
      steckt im Bau), `sha256sum harness/image-hash.txt` und die Image-ID von `:dev` sind vor und nach dem
      Lauf gleich, `make image-mutation-rm TAG=<Tag>` entfernt es, danach nennt `docker image ls` den Tag nicht
      mehr; die Zahl der dangling Volumes vor und nach dem Lauf (Erwartung: gleich) und kein `prune`. Verweigert
      die Berechtigungsschicht den realen Lauf, gilt Liefer-Punkt 3: Bericht und Rückfrage, kein Ersatzweg.
- [x] **Liefer-Punkt 2 — die Träger des Ziels.** (a) Vertrag `harness/targets/image-mutation.md` (Vorbild der
      Form: `harness/sensors/fmt-check.md`, `harness/targets/schema-rollout.md`): Zweck, Aufruf mit beiden Verben,
      Exit-Codes, Host-Werkzeuge (`bash`, `git`, `realpath`, `docker`; `git archive`/`tar` für die Kopie),
      **Anwendungsbeispiel** (Kopie ziehen, mutieren, bauen, per Override-Datei einer Scratchpad-Kopie des
      Runners an die Compose-Umgebung binden, `image-mutation-rm`) als Text, Grenze („Bau, kein Lauf-Beleg: das
      Image trägt keine Aussage über `harness/image-hash.txt`; ein Mutations-Image ist kein Release“; der Tag
      ist Sache des Aufrufers, zwei parallele Läufe mit demselben Tag überschreiben einander), kein Gate.
      (b) Zeile im Werkzeug-Verzeichnis von [`harness/README.md`](../../../../harness/README.md) §Sensors (drei
      Ziele: `make image-mutation`, `make image-mutation-rm`, `make test-image-mutation`; „kein Gate“,
      Bindung auf den Vertrag). (c) [`AGENTS.md`](../../../../AGENTS.md) §3.1, Absatz zur Mutationsprobe: ein Satz,
      der für ein Mutations-Image `make image-mutation` nennt (Kopie im Scratchpad, eigener Tag, `make image` mit
      mutiertem Baum bleibt verboten), mit Herkunft `· seit slice-harness-mutationsbild-und-verweigerte-aktion`.
      *Zu belegen durch:* Lesen der drei Stellen gegeneinander, `make docs-check` und der Suchlauf in §3.
- [x] **Liefer-Punkt 3 — die Regel „Verweigerte Aktion“.** (a) [`AGENTS.md`](../../../../AGENTS.md): ein neuer
      Abschnitt `### 3.15` (Begründung der Nummer in §3 „Ansatz“), Fassung kurz im Ist-Ton, mit „Falsch/Richtig“
      wie die Nachbar-Regeln (Vorbild: §3.9), Kernaussage wie in §1 (Entwurf), ohne Chronik; er nennt die
      Grenze („wirkt durch Lesen; der Classifier ist kein Teil des Repos; kein Sensor“) und die Herkunft
      (`BEO-PGC/ersatzweg-nach-verweigerter-aktion`, `· seit slice-harness-mutationsbild-und-verweigerte-aktion`).
      (b) [`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md): eine Prüfzeile in der
      MEDIUM-Liste, Klasse „Ersatzweg nach Verweigerung ohne Meldung“, mit der Einstufungsregel aus §3 „Ansatz“
      und dem Satz, dass der Reviewer nur liest, was Bericht und Artefakte tragen (steht die Verweigerung
      allein in einer Angabe des Auftraggebers, ist der Fund **übernommen**, nicht gemessen). (c)
      [`.claude/commands/implement-slice.md`](../../../../.claude/commands/implement-slice.md): ein Satz in der
      Docker-only-Zeile (Zeile 28, gelesen), der auf `AGENTS.md` §3.15 verweist — nicht die Regel wiederholt.
      (d) Register: `BEO-PGC/ersatzweg-nach-verweigerter-aktion` (vom Planner bei der Anlage dieses Slice
      angelegt, Ausgang *geplant*, Träger dieser Slice) erhält den Ausgang *verkörpert* mit auflösbarem Anker
      auf die drei Stellen; der Beleg des Slice selbst (kein Anfall oder ein Anfall) ist Datei oder Zeile in §7.
      *Zu belegen durch:* Lesen der drei Stellen gegeneinander (Regel, Prüfzeile, Verweis: dieselbe Klassenmenge,
      dieselbe Ausnahme), `make docs-check` (Anker, Links), `make doc-immutable RANGE=<Parent>..HEAD` Exit 0
      und der Suchlauf in §3; `git grep -n '3\.14' -- AGENTS.md` nennt den Rang-Zeiger unverändert.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/open/slice-harness-mutationsbild-und-verweigerte-aktion.md`
      läuft nach jeder Fixrunde mit den `diff`-Zeilen des Implementers durch
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: [`harness/README.md`](../../../../harness/README.md) §Sensors, `harness/targets/image-mutation.md`
      und `AGENTS.md` (Liefer-Punkte 2 und 3); das Benutzerhandbuch bleibt unberührt (keine Betreiber-Oberfläche).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang von
      `BEO-PGC/ersatzweg-nach-verweigerter-aktion` (Liefer-Punkt 3d); je ein Anfall im Lauf dieses Slice
      ist eine weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine Antwort und wird in §7 notiert
      (§8 nennt die Einträge, die dieser Slice trägt). Zusätzlich ein neuer Beleg (20.) bei
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (H-1) und ein narrativer Anfall ohne eigene
      Datei nach der Deckel-Regel (V-2, LOW) — §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** S — Schätzung, nicht gemessen: ein Skript mit zwei Verben, ein Tabellentest mit Stub-`docker`, drei
Makefile-Ziele, ein Vertrag, eine README-Zeile, ein Satz und ein Abschnitt in `AGENTS.md`, eine Prüfzeile im Skill,
ein Satz im Command, ein Register-Ausgang. Der Tabellentest ist der größte Teil.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/image-mutation.sh` | neu | Liefer-Punkt 1: Verben `build` und `rm`, Eingabeprüfung vor jedem Docker-Aufruf, ein Docker-Aufruf je Verb |
| `tools/harness/run-image-mutation-tests.sh` | neu | Liefer-Punkt 1: Tabellentest mit Stub-`docker` in einem Wegwerf-Repo im Temp-Verzeichnis; Prüfling per `TOOL=<Datei>` übersteuerbar (Mutationsläufe an Kopien) |
| `Makefile` | update | `image-mutation`, `image-mutation-rm`, `test-image-mutation` mit `.PHONY` und Hilfe-Zeilen, neben `fmt-check`/`test-fmt-check` (Zeilen 127–133) |
| `harness/targets/image-mutation.md` | neu | Liefer-Punkt 2a: Vertrag, Anwendungsbeispiel, Grenze |
| `harness/README.md` §Sensors | update | Liefer-Punkt 2b: Zeile im Werkzeug-Verzeichnis |
| `AGENTS.md` §3.1 (Absatz zur Mutationsprobe, Zeile 102) | update | Liefer-Punkt 2c: ein Satz zu `make image-mutation` |
| `AGENTS.md` (neuer Abschnitt `### 3.15`, nach dem Rang-Zeiger `### 3.14`) | update | Liefer-Punkt 3a |
| `.harness/skills/reviewer.md` (MEDIUM-Liste) | update | Liefer-Punkt 3b |
| `.claude/commands/implement-slice.md` (Zeile 28) | update | Liefer-Punkt 3c |
| `docs/plan/planning/observations/BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md` | update | Liefer-Punkt 3d: Ausgang *verkörpert* |

**Ansatz (Liste):**

- **Name und Tag.** Das Repository des Mutations-Images ist die Konstante `pg-change-feed-mutation` im Skript; der
  Tag ist die Eingabe `TAG`. Ein anderer Repository-Name als der von `:dev` (`ghcr.io/pt9912/pg-change-feed`)
  schließt die Kollision mit `:dev` durch die Form aus, nicht durch die Prüfung; die Ablehnung von `dev` und `latest`
  hält daneben die Verwechslung mit den Lauf-Tags fern (Eingabeseite des Tabellentests).
- **Eingabe ist feindlich.** `TAG` und `SRC` stammen aus dem Aufruf eines Rollenlaufs; `TAG` läuft gegen eine
  Positiv-Zeichenklasse (kein führendes `-`, kein `:`, kein `/`), `SRC` wird zum absoluten Pfad aufgelöst, bevor es
  Argument wird. Nachbar: `BEO-PGC/werkzeug-fuehrt-plan-inhalt-als-argument-aus` (Eingabeseite bindet der Test, hier
  je ein Fall pro Zeichenart).
- **Wurzel-Vergleich.** Das Skript liest die Repo-Wurzel mit `git rev-parse --show-toplevel` im aktuellen Verzeichnis und
  lehnt `SRC` ab, wenn der aufgelöste Pfad gleich der Wurzel ist oder mit `<Wurzel>/` beginnt. Eine Kopie im
  Scratchpad, ein `git worktree` außerhalb der Wurzel und ein Verzeichnis im Temp-Verzeichnis passieren. Der
  Tabellentest führt das Skript in einem Wegwerf-Repo aus (dort liegt `harness/image-hash.txt` als Fühler), damit ein
  mutiertes Skript, das die Datei schreibt, das echte Repo nicht trifft.
- **Exit-Codes.** 0 gebaut oder entfernt · 1 der Docker-Aufruf endet ≠ 0 (die Zeile nennt seinen Exit) · 2 Eingabefehler.
  `make` gibt jeden Exit ≠ 0 als seinen eigenen Exit 2 weiter; die Unterscheidung trägt die Make-Meldung `Fehler <n>`
  und die Ausgabe (wie bei `make fmt-check`).
- **Der Bau-Kontext ist der Kontext des Repos.** `git archive` nimmt jede getrackte Datei mit, darunter `.dockerignore`,
  `Dockerfile`, `go.mod`; es gibt kein `.gitattributes` mit `export-ignore` (gemessen: `git ls-files | grep -c
  gitattributes` nennt 0, Stand `84f60e6f`). Eine **neue** Datei, die die Mutation in der Kopie anlegt, liegt im
  Kontext, wenn `.dockerignore` sie zulässt (Default-Deny mit Allow-Liste: `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`);
  der Vertrag nennt das als Grenze, das Ziel ändert `.dockerignore` nicht.
- **Nummer der Regel.** §3.14 bleibt der Rang-Zeiger ([`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md) zitiert die Nummer, `AGENTS.md` §3.14 trägt den
  Grund). §3.1 wäre die falsche Stelle: die Regel gilt für jede verweigerte Aktion (`Write`, `git push`, ein
  Edit, ein Bau), nicht für Docker-only allein, und unter §3.1 läse sie sich als Teil der Toolchain-Regel (zu
  eng für ihren Geltungsbereich); §3.9 ist der Präzedenzfall einer rollenübergreifenden Ausführungsdisziplin als
  eigener Abschnitt. Deshalb `### 3.15`. §3.1 bekommt nur den Satz zu `make image-mutation` (Liefer-Punkt 2c) und
  verweist für die Verweigerung auf §3.15.
- **Einstufung im Reviewer-Skill (Vorschlag des Planners, der Reviewer/Architect darf ihn schärfen).** Grund:
  die Klasse der **verweigerten Aktion** bestimmt die Einstufung, weil der Ersatzweg dieselbe Wirkung hat. Die
  Klasse „Ersatzweg nach Verweigerung ohne Meldung“ steht in der **MEDIUM**-Liste; sie ist **HIGH**, wenn die
  verweigerte Aktion selbst in einer HIGH-Klasse liegt oder geteilten Zustand berührt (ein Host-Werkzeug am Repo,
  ein Push, `:dev` oder `harness/image-hash.txt`, ein Geheimnis) — dann gilt die Nachbarklasse „Docker-only-Verstoß“
  bzw. die HIGH-Liste des Ziels der Aktion, und dieser Fund tritt als zweites Feld dazu. Ein Ersatzweg mit Meldung im
  Bericht, aber ohne vorherige Rückfrage, ist **LOW** (die Wirkung ist sichtbar, der Auftraggeber kann sie zurücknehmen),
  sofern die verweigerte Aktion keine HIGH-Klasse trägt. Der Reviewer liest, was Bericht und Artefakte tragen: das
  Fehlen der Meldung ist aus einem Diff nicht ablesbar; steht die Verweigerung allein in einer Angabe des
  Auftraggebers, ist der Fund **übernommen** (`verifizierbar: nein`, wie F-3 des Reviews). Die Nachbarklasse „Docker-only-Verstoß“
  (HIGH-Liste, Zeile 232) bleibt unverändert und verweist nicht auf die neue Klasse.
- **Was der Slice am Reviewer-Skill nicht ändert.** Die Herkunft des Skills („geschärft 2026-09-09“) bekommt keine
  Datumszeile, der Skill trägt Herkunft je Punkt (`Herkunft: BEO-PGC/…`).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „wie ein Rollenlauf ein Image aus einem mutierten
Baum baut“ und „was ein Rollenlauf nach einer verweigerten Aktion tut“).** Suchraum: der ganze Baum ohne die drei
Ausnahmen von [`AGENTS.md`](../../../../AGENTS.md) §3.13 (`docs/reviews/**`, Records unter `done/`,
`.harness/baseline/**`), eine weitere Einschränkung nur mit Grund (Zeile 5 nimmt `internal/` aus: die Treffer dort
sind Testnamen einer Bezeichner-Prüfung, nicht Beschreibungen der Regel). Die Plan-Datei schließt das Werkzeug
aus. Stand ist `84f60e6f` (der Commit vor der Anlage dieses Slice, vom Planner am 2026-09-27 gemessen; ein Stand
ist eine Commit-Kennung, nie `HEAD`); der Implementer misst am Parent seiner Arbeit neu und trägt die `diff`-Zeilen
ein. Die Zahlen sind mit `git grep -n` gemessen, nicht übernommen. Drei Arten des Musters: der Symbolname
(Ziel- und Tag-Name, Abschnittsnummer), das Zählwort (die Abschnittsnummer `3.15`) und die Beschreibung samt
dem Hedge (Mutations-Image, Ersatzweg, Verweigerung).

```suchlauf
84f60e6f 0 -n -i -E 'Mutations-?Image|Wegwerf-?Image|mutiertes Image|pg-change-feed-mutd' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
84f60e6f 0 -n -E 'image-mutation|image-mutant|mutation-image|mutationsbild' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
84f60e6f 0 -n -E '3\.15' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
84f60e6f 8 -n -i -E 'Mutationsprobe|arbeitet auf einer Kopie' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
84f60e6f 2 -n -i -E 'verweiger|Ersatzweg|Berechtigungsschicht|Classifier|Klassifizierer' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!internal'
84f60e6f 11 -n -E 'Docker-only' -- AGENTS.md .claude .harness/skills
84f60e6f 0 -n -E 'kein make-Ziel|make-Ziel für ein Mutations' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 9 -n -i -E 'Mutations-?Image|Wegwerf-?Image|mutiertes Image|pg-change-feed-mutd' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 65 -n -E 'image-mutation|image-mutant|mutation-image|mutationsbild' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 6 -n -E '3\.15' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 9 -n -i -E 'Mutationsprobe|arbeitet auf einer Kopie' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 58 -n -i -E 'verweiger|Ersatzweg|Berechtigungsschicht|Classifier|Klassifizierer' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!internal'
diff 12 -n -E 'Docker-only' -- AGENTS.md .claude .harness/skills
diff 0 -n -E 'kein make-Ziel|make-Ziel für ein Mutations' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

**Nachmessung durch den Planner bei Closure** (Parent `04ae69a2`, nach der Fixrunde der
Verifikations-Findings V-1/V-2 und vor dem `git mv` nach `done/`): die Closure-Edits an
`harness/targets/image-mutation.md` (zwei neue Test-Tabellenzeilen) und
`tools/harness/run-image-mutation-tests.sh` (vier neue `rm`-TAG-Fälle) heben Zeile 2 von
62 auf 65 (drei zusätzliche Nennungen von „image-mutation“ in den beiden Dateien); der
neue Registerbeleg `docs/plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/evidence/slice-harness-mutationsbild-und-verweigerte-aktion.md`
trägt denselben Dateinamen wie diese Plan-Datei und wird deshalb von
`make suchlauf-nachmessen` wie die Plan-Datei selbst aus dem Suchraum ausgeschlossen
(sein Inhalt zählt in den obigen `diff`-Zeilen nicht mit); die Änderung an Zeile 5
(56 → 58) stammt aus der `state.md`-Fortschreibung desselben Registereintrags
(zwei Nennungen des vollen Slice-Namens, der „verweigert“ als Teilwort trägt). Beide
Werte mit `make suchlauf-nachmessen` gegen den aktuellen Arbeitsbaum nachgemessen
(gemessen, nicht übernommen).

**Nachmessung durch den Implementer** (Parent `5f97700a`, nach dem `next → in-progress`-Move und vor
jeder Inhaltsänderung dieses Laufs — identisch zu den `84f60e6f`-Werten außer den bereits vom Planner
committeten Register-Dateien `docs/plan/planning/observations/BEO-PGC/ersatzweg-nach-verweigerter-aktion/{state.md,observation.md,evidence/*.md}`,
die zwischen `84f60e6f` und `5f97700a` als eigener Commit `1d42ed1f` hinzukamen und die künftige
Slice-Kennung sowie „§3.15“ bereits nennen — keine neue Eigenschaft, sondern die Selbstreferenz der
Register-Anlage): Zeile 1 (Mutations-Image-Namen) 2, Zeile 2 (Ziel-/Tag-Namen) 3, Zeile 3 (`3.15`) 1,
Zeilen 4/6/7 unverändert (8/11/0), Zeile 5 (Verweigerungs-Wortfeld) 26 (Register-Einträge mehrerer
Slices zwischen `84f60e6f` und `5f97700a` nennen „verweigert“/„Berechtigungsschicht“ unabhängig von
diesem Slice). Die `diff`-Zeilen oben zählen den **endgültigen** Stand nach `git add -A` dieses Laufs:
`git grep` ohne Revision durchsucht standardmäßig nur **getrackte** Dateien des Arbeitsbaums — eine
neu angelegte, noch nicht `git add`ete Datei (`tools/harness/image-mutation.sh`,
`tools/harness/run-image-mutation-tests.sh`, `harness/targets/image-mutation.md`) bleibt bis zum
`add` unsichtbar; eine erste Zwischenmessung vor dem `add` zählte deshalb zu niedrig (6/21/5/8/54) und
ist durch die Werte oben ersetzt. **Gefunden** — `Makefile` (drei Ziele, `.PHONY`, Hilfe-Zeilen, zwölf
Treffer für Zeile 2), `tools/harness/image-mutation.sh` (Kopf-Kommentar und Code, fünfzehn Treffer),
`tools/harness/run-image-mutation-tests.sh` (Kopf-Kommentar und Fallnamen, zwölf Treffer),
`harness/targets/image-mutation.md` (der ganze Vertrag, vierzehn Treffer für Zeile 2), `harness/README.md`
(zwei neue Zeilen), `AGENTS.md` (der Satz in §3.1 und der neue Abschnitt `### 3.15`, vier Treffer für
Zeile 2), `.harness/skills/reviewer.md` (die neue MEDIUM-Zeile), `.claude/commands/implement-slice.md`
(der Verweis-Satz), `docs/plan/planning/observations/BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md`
(Ausgang *verkörpert*); **nicht gefunden/nicht angefasst** — `.claude/commands/plan-welle.md` und
`close-welle.md` (§1 Abgrenzung, bewusst nicht berührt), `harness/conventions/MR-003-…md` (immutable,
unverändert), `AGENTS.md` §3.14 (unverändert, geprüft mit `git grep -n '3\.14' -- AGENTS.md`: ein
Treffer, derselbe Rang-Zeiger-Absatz wie vor diesem Lauf).

| Träger | Befund (Stand `84f60e6f`, vom Planner gelesen) | Behandlung |
|---|---|---|
| `AGENTS.md` §3.1, Absatz zur Mutationsprobe (Zeilen 102–105) | „Eine Mutationsprobe (Reviewer, Verifier) arbeitet auf einer Kopie im Scratchpad; die Rücknahme ist `cp` oder `git checkout`“ — nennt für ein Image keinen Weg | Liefer-Punkt 2c: ein Satz zu `make image-mutation` |
| `AGENTS.md` §3.1 „Durchsetzung“ (Zeilen ~122–134) | beschreibt, was der Guard blockt; **fremd berührt** von `slice-harness-guard-blocked-python` (Liefer-Punkt 3b dort) | nicht angefasst; wer später startet, liest den Stand des anderen (§4) |
| `harness/README.md` §Sensors | trägt keine Zeile zu Mutations-Images; die Zeile `make test-command-guard` (Zeile 155) berührt der andere Harness-Slice | Liefer-Punkt 2b: neue Zeile, keine Änderung an bestehenden |
| `.harness/skills/reviewer.md` Zeile 232 (HIGH „Docker-only-Verstoß“), Zeile 306 | nennen Docker-only als Klasse; keine Klasse zu einer Verweigerung | Liefer-Punkt 3b (MEDIUM-Liste), HIGH-Punkt bleibt |
| `.claude/commands/implement-slice.md` Zeile 28 | die Docker-only-Zeile verweist auf `AGENTS.md` §3.1 (Host-Werkzeuge, in-place, Umleitung); kein Wort zu einer Verweigerung | Liefer-Punkt 3c |
| `.claude/commands/plan-welle.md` Zeile 33, `.claude/commands/close-welle.md` Zeile 23 | Docker-only-Zeilen ohne Verweis auf die Verweigerung | **nicht angefasst** (§1, Abgrenzung); der Reviewer kann einen Satz verlangen |
| `harness/conventions/MR-003-guard-inplace-textwerkzeug.md` Zeile 90 | „nennen diesen Ersatzweg“ — die Meldung des Guards nennt Edit/Write als Weg | bleibt wahr (`Accepted`, immutable); `AGENTS.md` §3.15 trägt dazu **keine** Ausnahme — der Satz beschreibt die Guard-Meldung selbst, nicht eine Ausnahme von der Melde-/Rückfrage-Pflicht |
| `AGENTS.md` §3.14 und `ADR-0100` §Teilfrage 4 | der Rang-Zeiger und das Zitat der Nummer | unverändert (Liefer-Punkt 3a, Prüfung `git grep -n '3\.14' -- AGENTS.md`) |
| **Fremde Datei** (Stand `84f60e6f`: `open/`; zwischenzeitlich `git mv` nach `done/` — `slice-harness-guard-blocked-python` schloss vor dem Start dieses Slice ab, kein gleichzeitiger Arbeitsbaum): `slice-harness-guard-blocked-python` §3 | trug die Träger-Tabelle zu `AGENTS.md` §3.1 „Durchsetzung“ und zu `harness/README.md` Zeile `make test-command-guard` | gemeldet, nicht mitgeändert; keine Adresse nötig: keine gemeinsame Stelle (Absatz zur Mutationsprobe und Abschnitt 3.15 liegen außerhalb) — mit dem Abschluss des anderen Slice ohnehin erledigt |
| Beobachtungs-Register `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet` | die Genauigkeit der Meldung einer Werkzeug-Ablehnung (1×) | bleibt getrennt; der Regel-Satz „Aufruf, Pfad, Wortlaut“ und die Abgrenzung im neuen Eintrag nennen ihn |
| `docs/plan/planning/welle-transformationen.md` §5, Plan `slice-wal-fehlerschwelle-ausgangsklasse`, `slice-start-vorlauf-grenze`, `slice-transformationen-e2e-abhilfe` | die Pläne nennen Mutationen (`grep -c -i -E 'Image mit\|mutiert\|Mutation'` je Plan: 7, 7 und 1, gemessen); ob die Läufe ein Image aus einem mutierten Baum brauchen, ist nicht gelesen — **Erwartung** aus dem Auftrag des Auftraggebers | keine Änderung; Reihenfolge in §4 |
| Beschreibung in `done/` und `docs/reviews/**` | Records und Reports der Läufe, darunter der Review dieser Sache | bleiben stehen (Record-Einfrierung) |

**Nicht durch Suchlauf fassbar:** ein Zeilen-Lokator auf `AGENTS.md` §3.1 oder auf die Zeilen des Skills — der
Suchlauf trifft Symbolnamen, keine Zahlen ([`AGENTS.md`](../../../../AGENTS.md) §3.13, Grenze); der Reviewer liest
sie gegen den Diff. Ebenso nicht durch Suchlauf fassbar: was Rollen-Aufträge des Orchestrators tragen (nicht im Repo).

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1; beobachtbar:
`ls docs/plan/planning/in-progress` nennt nur `roadmap.md`). Eine Nutzer-Entscheidung ist nicht mehr abzuwarten: beide
liegen vor (Herkunft oben).

**Reihenfolge (Empfehlung an den Orchestrator, mit Begründung):** nach `slice-leerlauf-phase-last-in-stuecken`
(liegt in `in-progress/`; das WIP-Limit ordnet das), nach
`slice-harness-guard-blocked-python` und **vor**
`slice-wal-fehlerschwelle-ausgangsklasse`, also
`slice-leerlauf-phase-last-in-stuecken` → `slice-harness-guard-blocked-python` →
`slice-harness-mutationsbild-und-verweigerte-aktion` → `slice-wal-fehlerschwelle-ausgangsklasse`. Eine technische Kante gibt
es nicht: der Slice berührt weder den Runner noch `internal/bootstrap/wiring.go` (die Dateien der Runner-Slices), und
kein anderer offener Plan ändert das Makefile-Segment neben `fmt-check`, den Abschnitt 3.15 oder die MEDIUM-Liste des
Skills. Gründe für „vor den Runner-Slices“, ohne Bindung: die Läufe der drei Runner-Slices fahren Mutationen am Produktionscode
(die Pläne nennen Mutationen, §3 Träger-Tabelle; ob ein Image nötig ist, ist **Erwartung**), und ohne das Ziel steht dort
der Ersatzweg des Vorfalls oder die Rückfrage an den Auftraggeber je Lauf; der Preis ist ein S-Slice Verzögerung auf dem Weg der
Welle (nicht gemessen). Gründe für „nach `slice-harness-guard-blocked-python`“: die Empfehlung an den Orchestrator lag
für diesen Slice zuerst vor (ohne Sperre ist jede Sitzung ein möglicher Anfall der Klasse), und die beiden Slices berühren
dieselbe Datei an verschiedenen Stellen (unten). Wählt der Orchestrator die andere Reihenfolge, ändert sich kein
Plan: keine Kante, kein Start-Trigger eines anderen Plans nennt einen der beiden.

**Gemeinsame Datei `AGENTS.md`, `harness/README.md` und `Makefile` mit `slice-harness-guard-blocked-python` — Konflikte
vermeiden.** Beide Slices ändern `AGENTS.md` §3.1, aber an verschiedenen Absätzen (dort „Durchsetzung“, hier der Absatz zur
Mutationsprobe und der neue Abschnitt 3.15), `harness/README.md` an verschiedenen Zeilen (dort die Zeile `make test-command-guard`,
hier neue Zeilen) und das `Makefile` an verschiedenen Zeilen (dort die Hilfe-Zeile von `test-command-guard`, hier ein neuer
Block neben `fmt-check`). Das WIP-Limit macht sie sequentiell, es gibt keinen gleichzeitigen Arbeitsbaum. Der später startende
Implementer misst seinen Suchlauf am Parent seiner Arbeit neu (der Parent enthält dann den Stand des früheren Slice)
und liest §3.1 und die README-Tabelle vor dem Editieren; die Träger-Tabelle nennt die fremden Stellen.

**Zwischenregel bis zum Träger (Orchestrator-Prompts, nicht im Repo).** Der Orchestrator gibt ab sofort an jede Rolle
(Planner, Architect, Implementer, Reviewer, Verifier, Validator) weiter: *„Verweigerte Aktion: nicht auf anderem Weg
wiederholen; im Bericht nennen (Aufruf, Wortlaut der Ablehnung, Ziel) und den Auftraggeber fragen.“* Solange
`make image-mutation` nicht steht, gilt für ein Mutations-Image dasselbe: der Weg der zwei früheren Verifikationen
(`make image` aus einer `git archive`-Kopie, danach `:dev` wiederherstellen) ist nicht Nutzer-entschieden; ein Lauf, der ein
Mutations-Image braucht, fragt den Auftraggeber. Mit dem Merge von Liefer-Punkt 3 steht die Regel im Repo; der Prompt-Satz
ist danach ein Verweis auf `AGENTS.md` §3.15 statt einer eigenen Fassung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht zu erwarten (Umfang S); trägt der Tabellentest mehr als
  erwartet, ist Liefer-Punkt 3 (Regel und Träger) der abtrennbare Teil — er hängt an nichts aus Liefer-Punkt 1.
- `in-progress` → `open` (blockiert): die Berechtigungsschicht verweigert den realen Lauf von Liefer-Punkt 1 (Risiko 3) und der
  Auftraggeber gibt keinen Weg frei — dann geht der Slice nach `open/`, bis der Auftraggeber über den realen Lauf entschieden
  hat (er ist ein Kriterium der DoD, keine Grenze, die der Bericht benennt und übergeht); kein Ersatzweg.

## 5. Closure-Trigger

DoD vollständig (das Ziel mit Tabellentest und gesehenen Mutationen und dem einen realen Lauf, die Träger des Ziels, die
Regel in `AGENTS.md`, Skill und Command) + `make gates` grün + Review-Report ohne offenes HIGH oder MEDIUM + Verifikation, dass die
DoD trägt (`make test-image-mutation` grün mit gedruckter Zahl, die Mutationen rot gesehen, der reale Lauf mit den Vorher-/Nachher-Werten
von `:dev` und `image-hash.txt`) + Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: [`AGENTS.md`](../../../../AGENTS.md)
§3.15 und der Weg für Mutations-Images in §3.1; neuer Sensor: keiner — `make test-image-mutation` bindet ein Werkzeug ohne Gate; die
Regel „Verweigerte Aktion“ hat bewusst keinen Sensor, §1).

## 6. Risiken und offene Punkte

- **1. Das Ziel baut aus dem Arbeitsbaum statt aus einer Kopie.** Der Wurzel-Vergleich hängt an aufgelösten Pfaden:
  ein Symlink auf die Wurzel oder ein relativer Pfad könnte ihn umgehen, wenn er nicht aufgelöst wird. *Erwartet, zu belegen
  durch:* die Fälle Wurzel und Verzeichnis unter der Wurzel im Tabellentest, dazu ein Fall mit einem Symlink auf die Wurzel,
  soweit der Aufbau ihn zulässt (`realpath` löst ihn auf, *hergeleitet*, im Lauf zu bestätigen), und die Mutation „Vergleich
  entfernt“ (DoD 1). **Ausgang:** eingetreten, in der Fixrunde behoben — Review-Fund H-1 demonstrierte den
  Exploit real (`SRC=.` aus der Repo-Wurzel baute mit entfernter `realpath`-Zeile tatsächlich ein Image);
  drei neue Testfälle (relativer Pfad `.`, relativer Pfad `sub`, Symlink auf die Wurzel) binden die
  `realpath`-Auflösung jetzt an ihre Eingabeseite, vom Verifier mit einem eigenen Exploit-Nachlauf gegen
  den erweiterten Tabellentest bestätigt (Exit 1, exakt die drei neuen Fälle rot, 27 übrige grün).
- **2. Tag-Kollision mit parallel laufenden Rollen.** Zwei Läufe mit demselben `TAG` überschreiben einander (Docker ersetzt das Image
  unter dem Tag); der Vertrag nennt es als Grenze, das Ziel erzwingt keinen Namen je Rolle. *Erwartet, zu belegen durch:* die Grenze im
  Vertrag und ein Namensvorschlag im Anwendungsbeispiel (Kürzel des Vorgangs im Tag). **Ausgang:** nicht
  eingetreten — die Grenze steht im Vertrag (`harness/targets/image-mutation.md` §Grenze Punkt 2) mit dem
  Namensvorschlag, das Ziel selbst erzwingt keinen Namen je Rolle; kein Mechanismus nötig laut Plan-Design.
- **3. Die Berechtigungsschicht verweigert auch das neue Ziel.** `Bash(make:*)` steht auf der Allow-Liste (gelesen), was der Classifier
  für `make image-mutation` mit einer mutierten Kopie entscheidet, ist nicht gemessen. *Erwartet, zu belegen durch:* der reale Lauf von DoD 1;
  bei einer Verweigerung gilt Liefer-Punkt 3 (Bericht, Rückfrage, kein Ersatzweg) — der Slice wendet die Regel, die er liefert, auf
  sich selbst an. **Ausgang:** nicht eingetreten — Register `BEO-PGC/ersatzweg-nach-verweigerter-aktion/state.md`
  nennt keine Ablehnung während der Implementierung (übernommen, nicht gemessen — kein Log-Träger im Repo —,
  aber plausibel und ohne Widerspruch in keinem Artefakt, so auch Review I-2 und Verifikation §5).
- **4. `:dev` oder `harness/image-hash.txt` wird überschrieben.** Ausgeschlossen durch die Form (anderer Repository-Name, kein
  `--metadata-file`), gebunden durch den Fühler im Wegwerf-Repo des Tabellentests und den realen Lauf (Image-ID von `:dev` und `sha256sum`
  der Hash-Datei vorher/nachher). *Erwartet, zu belegen durch:* die zwei Mutationen aus DoD 1 („`-t` auf `:dev`“, „schreibt die Hash-Datei“).
  **Ausgang:** nicht eingetreten, real bestätigt — die Image-ID von `:dev` und `sha256sum harness/image-hash.txt`
  waren vor und nach dem realen Lauf gleich (Vertrag §Test „Realer Lauf“); der Verifier bestätigte zusätzlich,
  dass `:dev` seit vor dem Slice-Beginn unverändert ist (kein anderer Prozess hat es neu gebaut).
- **5. Der Stub-`docker` bindet die Argumente, nicht das Verhalten des Daemons** (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`).
  Die Fälle belegen, was der Aufrufer übergibt; dass Docker unter dem Namen ein Image anlegt und `rmi` genau dieses entfernt, belegt der
  eine reale Lauf. *Erwartet, zu belegen durch:* die Mutationsliste (DoD 1, je Zeile Stelle, Instanz, gesehene Farbe) und der reale Lauf.
  **Ausgang:** erfüllt/entfallen — der reale Lauf (Vertrag §Test) und die zehn Mutationen (acht ursprüngliche
  plus die zwei durch Fixrunde und Closure ergänzten, §7) schließen die Lücke; der Reviewer fand mit einer
  eigenen, nicht vertraglich benannten Mutation eine reale Testlücke (H-1, in der Fixrunde gebunden), der
  Verifier mit einer weiteren (V-2, in dieser Closure gebunden) — der Stub bindet damit die Eingabeseite
  vollständig, das Verhalten des Daemons trägt weiterhin nur der eine reale Lauf.
- **6. Eine neue Datei in der Kopie erreicht den Bau nicht** (`.dockerignore` mit Default-Deny,
  `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`): eine Mutation, die eine neue Go-Datei anlegt, bliebe im Image aus,
  und der Lauf zöge einen falschen Schluss. *Erwartet, zu belegen durch:* die Grenze im Vertrag und der reale Lauf, der eine bestehende
  Datei mutiert. **Ausgang:** weiter offen — der reale Lauf mutierte eine bestehende Go-Zeile, keine neue
  Datei; die Grenze steht dokumentiert im Vertrag (§Grenze Punkt 3), aber unbelegt real getestet (so auch der
  Verifier, §9 Punkt „Risiko 6“). Kein Blocker: dokumentierte Grenze statt „entfallen“.
- **7. Die Regel „Verweigerte Aktion“ verspricht mehr, als sie kann** (`BEO-PGC/regel-weiter-als-ihr-sensor`, verkörpert). Sie wirkt
  durch das Lesen des Berichts; der Reviewer sieht die Verweigerung nur, wenn Bericht oder Artefakt sie tragen, der Classifier ist kein Teil
  des Repos. *Erwartet, zu belegen durch:* die Grenz-Zeile in `AGENTS.md` §3.15 und der Satz „übernommen, nicht gemessen“ in der Prüfzeile
  des Skills; der Reviewer liest die drei Stellen gegeneinander. **Ausgang:** wie geplant, dokumentierte
  Grenze — `AGENTS.md` §3.15 trägt den Grenz-Absatz („wirkt durch Lesen; kein Sensor; der Classifier ist kein
  Teil des Repos“), die Skill-Prüfzeile nennt „übernommen, nicht gemessen“ für den Fall, dass die Verweigerung
  allein in einer Angabe des Auftraggebers steht (real angewendet: Review I-2/Verifikation §5 zum eigenen
  Registereintrag).
- **8. Die Ausnahme „die Ablehnung führt den Weg selbst an“ ist Auslegung des Planners, nicht Wortlaut der Nutzer-Entscheidung.** Sie
  verhindert eine Frageflut (die Meldung des Guards nennt Edit/Write, `make`), sie könnte aber zu weit greifen: ein Lauf könnte jede
  Ablehnung als „Weg genannt“ lesen. *Erwartet, zu belegen durch:* der Reviewer liest die Fassung von §3.15 gegen den Wortlaut des Nutzers
  („im Bericht genannt und vor einem Ersatzweg zurückgefragt“); ist die Ausnahme zu weit, geht der Fund an den Auftraggeber, nicht in eine
  stille Streichung. **Ausgang:** entfallen — die Ausnahme wurde nicht ausgeliefert (§3.15 enthält keine
  Ausnahme für den Fall, dass die Ablehnung selbst einen Weg nennt; Review-Fund M-1, in der Fixrunde behoben
  und vom Verifier bestätigt), die Frage „ist sie zu weit gefasst?“ ist damit gegenstandslos.
- **9. Die neuen Fälle laufen nicht** (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen durch:* die gedruckte
  Zahl von `make test-image-mutation` (Instanz A, `AGENTS.md` §3.12) und die Farbe der Mutationen. **Ausgang:** entfallen,
  bestätigt real gelaufen — `make test-image-mutation` druckt „alle Fälle bestanden“ (Exit 0, gemessen in dieser
  Closure); die drei `realpath`-Fälle (H-1) laufen aktiv (eigene Nachzählung des Verifiers: exakt 3 Fälle rot bei
  Entfernung der `realpath`-Zeile, §5 der Verifikation) und die vier neuen `rm`-TAG-Fälle (V-2, in dieser Closure
  ergänzt) laufen ebenfalls aktiv (eigener Mutationslauf: Entfernung des `check_tag`-Aufrufs im `rm`-Zweig färbt
  genau diese vier Fälle rot, Exit 1, 8 `FEHLER:`-Zeilen).
- **10. Ein Träger bleibt stehen** (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Deckel): `plan-welle.md` und `close-welle.md` tragen die
  Docker-only-Zeile ohne Verweis auf §3.15 (§3, bewusst nicht angefasst). *Erwartet, zu belegen durch:* das Suchlauf-Feld und der Review; ein
  Fund verlangt einen Satz dort, nicht eine stille Mitänderung. **Ausgang:** entfallen, wie geplant nicht
  angefasst — kein Fund im Review oder in der Verifikation, der einen Nachtrag an `plan-welle.md`/`close-welle.md`
  verlangt.

## 7. Closure-Notiz

- **Was hat funktioniert:** Das eigene `make`-Ziel mit eigenem Repository-Namen
  (`pg-change-feed-mutation`, getrennt von `:dev`) trennt den Lauf-Beleg-Pfad sauber vom
  Mutationsbild; der Tabellentest mit Stub-`docker` (Argumentzahl **und** -inhalt gebunden)
  fing alle acht vertraglich benannten und zwei vom Reviewer selbst gefahrene Mutationen
  korrekt rot, `make gates`/`make suchlauf-nachmessen`/`make doc-immutable` blieben über
  die ganze Fixrunde grün. Die Regel `AGENTS.md` §3.15 entstand ohne Ausnahme, wie von der
  Nutzer-Entscheidung „ja“ verlangt, und wurde vom Review selbst sofort angewendet (der
  Reviewer strich einen versehentlichen `python3`-Aufruf statt ihn zu wiederholen, §Bericht
  des Reviews) — ein Selbstbeleg noch vor der Closure.
- **Was ging anders als geplant:** Zwei reale Testlücken derselben Klasse
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`) blieben trotz acht vertraglich
  vorbedachter Mutationen bestehen: H-1 (HIGH, Reviewer) — die `realpath`-Auflösung von
  `SRC` war an keinem Testfall gebunden, ein realer Exploit (`SRC=.`) baute tatsächlich aus
  dem Arbeitsverzeichnis; in der Fixrunde durch drei neue Fälle gebunden. V-2 (LOW,
  Verifier) — derselbe Mechanismus am `rm`-Zweig: `check_tag` prüft `TAG` im Code, aber
  kein Testfall rief `rm` mit einem ungültigen `TAG`; in dieser Closure durch vier neue
  Fälle gebunden (Mutation „`check_tag`-Aufruf im `rm`-Zweig entfernt“ eigens nachgefahren:
  Exit 1, exakt die vier neuen Fälle rot). Zusätzlich blieb der Vertrag
  (`harness/targets/image-mutation.md`) nach der ersten Fixrunde unvollständig
  nachgezogen (V-1, LOW) — die Test-Tabelle nannte weiterhin nur „acht“ Mutationen, ohne
  die neue `realpath`-Bindung; in dieser Closure ergänzt (zehn Zeilen, „Menge der
  Erprobung: zehn“).
- **Steering-Loop-Eintrag (Lerneintrag):**
  - **Geschärfte Regel:** Eine Eingabeprüfung, die **beide** Verben eines Skripts teilen
    (`check_tag` in `build` **und** `rm`), braucht Testfälle in **beiden** Verben — ein
    Tabellentest, der eine geteilte Prüfung nur über den zuerst geschriebenen Zweig testet,
    lässt eine Mutation am zweiten Zweig unbemerkt grün (V-2, dieselbe Ausprägung wie H-1
    einen Zweig weiter). Ergänzt den bereits verkörperten Eintrag
    `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` um den 20. Beleg (HIGH, eigene
    Datei trotz Deckel) und benennt den 21. Fall (LOW, V-2) narrativ statt mit eigener
    Datei — die Regel selbst (Reviewer-Skill HIGH „Zusage ohne Bindung an ihre
    Eingabeseite“) bleibt unverändert, sie hat in diesem Slice zweimal gewirkt.
  - **Neuer Sensor:** keiner — `make test-image-mutation` bindet ein Werkzeug ohne Gate,
    wie im Plan (§5) vorgesehen; die zehn Mutationen plus die vier neuen `rm`-TAG-Fälle
    laufen darüber.
  - **Benannte Spec-Lücke:** keine — dies ist ein Harness-Werkzeug, keine Spec-Stelle
    (§1 „Berührte Spec-Stellen“).
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-PGC/ersatzweg-nach-verweigerter-aktion`: Ausgang **verkörpert**, bereits mit
    Liefer-Punkt 3d fortgeschrieben (Anker auf `AGENTS.md` §3.15, die Skill-Prüfzeile und
    den Verweis im Implementer-Command); kein Anfall im Lauf dieses Slice
    (`state.md`-Absatz „Anfall im Lauf dieses Slice: keiner“).
  - `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`: **zwei Anfälle** — H-1 (HIGH)
    bekam eine eigene Beleg-Datei (`evidence/slice-harness-mutationsbild-und-verweigerte-aktion.md`,
    20. Beleg, Zähler jetzt 20×), V-2 (LOW, gebunden in dieser Closure) bleibt nach der
    Deckel-Regel ohne eigene Datei und ist oben (Steering-Loop-Eintrag) sowie im
    Register-`state.md` narrativ genannt.
  - Alle übrigen im Plan (§8) gesichteten Einträge: kein Anfall in diesem Slice.
- **Folge-Slices:** keine — die im Plan §1 ausdrücklich ausgeschlossenen Punkte (Änderung
  an `make image`/`Dockerfile`/`.dockerignore`, Neubewertung von §3.1, Sensor für die
  Regel „Verweigerte Aktion“, die drei anderen Commands) bleiben eigenständige,
  unbeauftragte Entscheidungen; Risiko 6 (`.dockerignore`-Grenze mit neuer Datei) bleibt
  eine dokumentierte Grenze ohne eigenen Folge-Slice, kein Fund dieser Closure verlangt
  einen.
- **Risiken aus §6:** je ein Ausgang eingetragen (Risiko 1 eingetreten/behoben, 2–5
  nicht eingetreten/erfüllt, 6 weiter offen als dokumentierte Grenze, 7 wie geplant, 8
  entfallen, 9 entfallen/bestätigt, 10 entfallen/wie geplant nicht angefasst — siehe §6).
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure trägt sie
  selbst, nach dem `git mv` nach `done/` — **Anker:** `AGENTS.md` §3.15,
  `harness/targets/image-mutation.md`, `.harness/skills/reviewer.md` MEDIUM-Liste,
  `.claude/commands/implement-slice.md` lösen auf (`make docs-check` grün). **Folge-Slice:**
  keiner (siehe oben). **Register:** `BEO-PGC/ersatzweg-nach-verweigerter-aktion`
  (verkörpert) und `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (20×, ein neuer
  Beleg) beide fortgeschrieben.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`/`PGC`
(Greenfield); `tools/harness`, `harness/`, `.claude` und `.harness/skills` sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler gemessen am
2026-09-27 mit `ls evidence | wc -l` je Eintrag, Stand `84f60e6f`) —
`BEO-PGC/ersatzweg-nach-verweigerter-aktion` (neu, vom Planner bei der Anlage dieses Slice angelegt, 1×: der Gegenstand von
Liefer-Punkt 3; Risiken 3, 7, 8), `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet` (offen, 1×: Nachbar, bleibt
getrennt; die Regel trägt den Satz „Aufruf, Pfad, Wortlaut“), `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration`
(verkörpert, 2×) und `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 6×: die Regel zur Docker-only-Klasse und der
Guard, den dieser Slice nicht berührt; Abgrenzung im neuen Eintrag), `BEO-PGC/regel-weiter-als-ihr-sensor` (verkörpert, 4×, Risiko 7),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (Deckel; Risiko 5, die Mutationen je Eingabeseite),
`BEO-PGC/werkzeug-fuehrt-plan-inhalt-als-argument-aus` (offen, 1×: `TAG` und `SRC` sind Eingabe eines Rollenlaufs, §3 „Ansatz“),
`BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×: das Ziel schreibt nichts im Repo, der Fühler im Test),
`BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt` (`image-hash.txt` bleibt unberührt),
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×, Risiko 6), `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×,
Risiko 9), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Deckel: Suchlauf-Feld in §3, Risiko 10),
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` und `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Deckel: DoD und
Ausgangslage nennen ihren Beleg-Anker oder sind als Erwartung formuliert), `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar`
und `BEO-PGC/kommentar-herkunft-als-kette`: die Kopfkommentare von Skript und Tabellentest sind Kommentare in Skripten (Indikativ,
höchstens eine Kennung, keine Slice-Nummer); übrige Einträge gesichtet, kein Bezug.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
