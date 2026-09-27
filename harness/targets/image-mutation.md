# `make image-mutation` / `make image-mutation-rm` — Mutations-Image mit eigenem Tag

## Vertrag

`make image-mutation SRC=<Verzeichnis> TAG=<Tag>` baut aus einer vom Aufrufer
bereitgestellten Kopie ein Image unter einem **eigenen** Repository-Namen
(`pg-change-feed-mutation:<TAG>`), getrennt vom Lauf-Beleg-Pfad
([`ADR-0044`](../../docs/plan/adr/0044-image-beleg-semantik.md),
[`ADR-0103`](../../docs/plan/adr/0103-image-hash-lokal-statt-committet.md)):
kein `--metadata-file`, kein `-t` auf `ghcr.io/pt9912/pg-change-feed:dev`,
kein `--push`, keine Datei im Repo. `make image-mutation-rm TAG=<Tag>`
entfernt genau dieses eine Image (`docker rmi`, kein `-f`, kein `prune`).

Der Bau ist genau:

```text
docker buildx build --load -t pg-change-feed-mutation:<TAG> <SRC>
```

`SRC` wird vor dem Aufruf zu einem absoluten Pfad aufgelöst und steht als
letztes Argument. Beide Ziele rufen `tools/harness/image-mutation.sh`
(Verben `build`/`rm`); jede Eingabeprüfung läuft **vor** dem jeweiligen
Docker-Aufruf, ein Docker-Fehler endet mit Exit 1 und einer Zeile, die
seinen Exit-Code nennt.

**Was dieses Ziel nicht tut:** Es zieht die Kopie nicht selbst (das bleibt
Sache des Aufrufers, siehe §Anwendungsbeispiel), es ändert `make image`, das
`Dockerfile` und `.dockerignore` nicht, es bindet die Compose-Umgebung nicht
an, und es räumt nie über `prune` auf.

**Kein Gate.** Das Ziel steht in keinem Gate-Bündel (`make gates`); die
Aufnahme als Gate braucht eine ADR ([`AGENTS.md`](../../AGENTS.md) §3.6, §4).

## Aufruf

```text
make image-mutation SRC=<Verzeichnis> TAG=<Tag>
make image-mutation-rm TAG=<Tag>
tools/harness/image-mutation.sh build <SRC> <TAG>   # direkter Aufruf
tools/harness/image-mutation.sh rm <TAG>
```

Ohne `SRC` bzw. `TAG` bricht das Makefile-Ziel mit `$(error …)` ab, bevor
irgendein Befehl läuft (Vorbild: `make suchlauf-nachmessen`).

Host-Werkzeuge: `bash`, `git` (der Wurzel-Vergleich), `realpath` und `docker`
([`AGENTS.md`](../../AGENTS.md) §3.1, Klasse „Host-Werkzeug ohne
Installation“). `git archive` und `tar`, mit denen der Aufrufer die Kopie
zieht, gehören derselben Klasse an; dieses Skript ruft sie nicht selbst auf.

## Eingabeprüfung (vor jedem Docker-Aufruf, Exit 2)

- `SRC` oder `TAG` fehlt.
- `TAG` verletzt die Zeichenklasse `[a-z0-9][a-z0-9_.-]{0,62}`.
- `TAG` ist `dev` oder `latest` (die Lauf-Beleg-Tags).
- `SRC` ist kein Verzeichnis.
- `SRC` trägt kein `Dockerfile`.
- `SRC` trägt kein `go.mod`.
- `SRC` ist — nach Auflösung beider Pfade — die Repo-Wurzel
  (`git rev-parse --show-toplevel`) oder liegt unter ihr: eine Mutation baut
  nie aus dem Arbeitsbaum ([`AGENTS.md`](../../AGENTS.md) §3.1).

## Exit-Codes

| Exit | Bedeutung |
|---|---|
| 0 | gebaut (`build`) bzw. entfernt (`rm`) |
| 1 | der Docker-Aufruf endet ≠ 0 (die Zeile nennt seinen Exit-Code) |
| 2 | Eingabefehler (siehe oben) oder unbekanntes Verb |

## Anwendungsbeispiel

```text
# 1. Kopie ziehen (Host-Werkzeuge, Klasse "ohne Installation")
git archive HEAD | tar -x -C /pfad/im/scratchpad

# 2. Mutation auf der Kopie, mit Edit/Write des Laufs (AGENTS.md §3.1)
#    — z. B. eine Zeile in cmd/pg-change-feed/main.go ändern.

# 3. Bauen
make image-mutation SRC=/pfad/im/scratchpad TAG=mein-vorgang-1

# 4. An eine Compose-Umgebung binden — eine Override-Datei einer
#    Scratchpad-Kopie des Runners, nicht die committete compose.yaml:
#    services: { pg-change-feed: { image: pg-change-feed-mutation:mein-vorgang-1 } }

# 5. Aufräumen
make image-mutation-rm TAG=mein-vorgang-1
```

## Wer es aufruft

- **Reviewer, Verifier, Implementer** — jeder Rollenlauf, der eine Mutation
  gegen ein laufendes System belegen muss und dafür ein Image braucht
  ([`AGENTS.md`](../../AGENTS.md) §3.1, Absatz zur Mutationsprobe).

## Grenze

1. **Bau, kein Lauf-Beleg.** Das erzeugte Image trägt keine Aussage über
   `harness/image-hash.txt`; ein Mutations-Image ist kein Release und wird
   nie veröffentlicht (kein `--push`, kein Multi-Arch).
2. **Der Tag ist Sache des Aufrufers.** Zwei gleichzeitige Läufe mit
   demselben `TAG` überschreiben einander (Docker ersetzt das Image unter
   dem Tag) — ein Namensvorschlag: ein Kürzel des Vorgangs im Tag.
3. **Der Bau-Kontext ist der Kontext des Repos.** `git archive` nimmt jede
   getrackte Datei mit, darunter `.dockerignore`; eine **neue** Datei, die
   die Mutation in der Kopie anlegt, erreicht den Bau nur, wenn
   `.dockerignore`s Allow-Liste sie zulässt (Default-Deny).
4. **Kein Sensor für die Ablehnung eines Aufrufs.** Verweigert die
   Berechtigungsschicht den Aufruf dieses Ziels selbst, gilt
   [`AGENTS.md`](../../AGENTS.md) §3.15: melden, den Auftraggeber fragen,
   nicht auf anderem Weg wiederholen.

## Test

`make test-image-mutation` (`tools/harness/run-image-mutation-tests.sh`)
fährt einen Tabellentest mit einem Stub-`docker` in einem Wegwerf-Repo im
Temp-Verzeichnis (netzlos, kein echter Bau); die Argumente des
Docker-Aufrufs hält der Stub fest — der Fall belegt, was der Aufrufer
übergibt, nicht, dass der Daemon es einhält (Vorbild:
[`harness/sensors/fmt-check.md`](../sensors/fmt-check.md) §Test). Je Zusage
die Mutation ihrer Eingabeseite, gesehenes Rot:

| Zusage | Mutation am Werkzeug | Fall, der rot wird |
|---|---|---|
| genau `-t pg-change-feed-mutation:<TAG>`, kein `:dev` | `-t` auf `ghcr.io/pt9912/pg-change-feed:dev` gesetzt | gültige Argumente (Docker-Argument Zeile 5) |
| kein `--metadata-file`, kein zusätzliches Argument | `--metadata-file harness/image-hash.raw` angehängt | gültige Argumente (Argumentzahl 6, Zeilen 4–6) |
| `TAG` `dev`/`latest` abgelehnt | die Ablehnung von `dev`/`latest` entfernt | `TAG 'dev'` · `TAG 'latest'` |
| `TAG`-Zeichenklasse durchgesetzt | die Zeichenklasse auf „alles“ gelockert | `TAG 'Foo'` · `TAG 'a b'` · `TAG '-x'` · `TAG 'a:b'` · `TAG '../x'` |
| `SRC` darf nicht die Repo-Wurzel oder darunter sein | der Wurzel-Vergleich entfernt | `SRC ist Repo-Wurzel` · `SRC unter der Wurzel` |
| `SRC` muss ein `Dockerfile` tragen | die Dockerfile-Prüfung entfernt | `SRC ohne Dockerfile` |
| ein Docker-Fehler endet mit Exit 1 und nennt den Exit-Code | die Auswertung des Docker-Exit-Codes entfernt | `Docker-Fehler (build)` |
| `rm` ruft `docker rmi` ohne `-f` | `-f` an `docker rmi` angehängt | `rm` (Docker-Argument Zeile 2, Argumentzahl) |

Menge der Erprobung: die acht oben genannten Mutationen, je einmal gegen
eine Kopie des Werkzeugs gelaufen (`TOOL=<Kopie> bash
tools/harness/run-image-mutation-tests.sh`); jede zeigte den genannten Fall
rot.

**Realer Lauf (einmal, kein Teil des Tabellentests):** `git archive HEAD` in
ein Scratchpad-Verzeichnis, eine Go-Zeile dort mit Edit geändert (`var
version = "dev"` → ein anderer Wert in `cmd/pg-change-feed/main.go`), `make
image-mutation SRC=<Kopie> TAG=mutprobe1` gebaut — die Image-ID
(`sha256:2e5c9410…`) weicht von der Image-ID von `:dev`
(`sha256:a8ef48aab9…`) ab, `sha256sum harness/image-hash.txt` und die
Image-ID von `:dev` waren vor und nach dem Lauf gleich, `harness/image-hash.raw`
existierte weder vorher noch nachher, `make image-mutation-rm TAG=mutprobe1`
entfernte das Image danach (`docker image ls` nennt den Tag nicht mehr), die
Zahl der dangling Volumes blieb bei 36 (vorher/nachher gemessen, kein
`prune`).
