# Verifikationsbericht: slice-generated-sync-tar-export — 2026-09-21

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-generated-sync-tar-export.md` §2)
und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff
als solchen (Reviewer-Aufgabe, zweifach abgeschlossen mit
`docs/reviews/review-slice-generated-sync-tar-export.md` und
`docs/reviews/review-slice-generated-sync-tar-export-fixrunde.md`) und
**nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** vier Commits auf `main`
(`ADR-0084`, `ADR-0060`, wellenlos):

- `1cce85b5` — Slice-Plan nach `in-progress/` verschoben (Parent-Stand für
  den Scope-Diff dieses Berichts).
- `93f5545b` — Implementer-Commit (Bind-Mount durch tar-Stream-Export
  ersetzt).
- `770d23a0` — Review-Report, 1 HIGH (F-1).
- `32a58868` — Fixrunden-Commit (Kopf-Kommentar gekürzt).
- `9293cc83` — Fixrunden-Review-Report, 0 HIGH/MEDIUM/LOW, 1 INFO.

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`,
`harness/conventions.md`, `AGENTS.md`, den vollständigen Slice-Plan
(§1–§8), `ADR-0084` (vollständig: §Kontext/§Entscheidung/§Verglichene
Alternativen/§Konsequenzen/§Fitness-Function/§Geschichte), beide
Review-Reports, `tools/harness/generated-sync.sh` (vollständig),
`harness/mk/generated-sync.mk` und `harness/sensors/generated-sync.md`
selbst gelesen. Nichts wurde ungeprüft übernommen: eigener `grep`-Lauf auf
`docker run`/`-v`/`--user`/`RUN_USER`, eigener realer
`make generated-sync`-Lauf, eigener realer `make proto-generate`-Lauf mit
Diff-Kontrolle, eigene Verifikation des Risiko-2-Belegs am
Einführungs-Commit `338cfe00` (slice-104), eigener ungepipter
`make gates`-Lauf mit direkter Exit-Code-Prüfung (`AGENTS.md` §3.9),
eigene isolierte `make doc-commits`/`make doc-immutable`-Läufe über den
exakten Slice-Commit-Bereich, eigener Backtick-Paritäts-Nachzähl über alle
sieben geänderten Dateien, eigener `git diff --stat 1cce85b5..HEAD`
Scope-Check — **und eine eigene, dritte, unabhängige Reproduktion des
DoD-Punkts-3-Auslösers** (§2 unten), auf einer Maschine mit real
identischer `mounts: []`-Colima-Konfiguration wie die
Implementer-Maschine.

---

## 1. DoD-Zeile für DoD-Zeile — §2

### 1.1 „`make generated-sync` läuft ohne jeden `docker run -v`-Bind-Mount" — bestätigt

```
$ grep -n "docker run" tools/harness/generated-sync.sh
16:# host-seitig `docker run --rm --network none <image> | tar -x -C
55:# `set -o pipefail` macht die Pipe `docker run | tar -x` sicher (AGENTS.md
83:docker run --rm --network none "$GENERATED_SYNC_IMAGE" | tar -x -C "$out_dir"
```

Genau **eine** ausführbare `docker run`-Zeile (83), keine `-v`, kein
`--user`; `grep -n -- "-v \|RUN_USER\|--user"` trifft nur einen
Kommentar-Treffer (Zeile 24, Prosa „`--user`-Workaround", kein Code).
Eigener Lauf: `make generated-sync` → Exit 0, `generated-sync: OK` (Stufe
`proto-export`, geprüfte Dateien beide `.pb.go`), `git status --porcelain`
davor **und** danach leer. **DoD-Zeile korrekt gesetzt.**

### 1.2 Modulpfad-/`.proto`-Erkennungs-Entscheidung (Risiko 2) — eigenständig nachgeprüft, bestätigt

Unabhängig vom Reviewer-Zitat selbst gegen den Einführungs-Commit der
Stufe geprüft:

```
$ git log --oneline -S "FROM proto AS proto-export" -- Dockerfile
338cfe00 feat(proto): make proto-generate mount-los (ADR-0060)

$ git show 338cfe00:Dockerfile | sed -n '/FROM proto AS proto-export/,/^$/p'
FROM proto AS proto-export
COPY proto/ proto/
RUN mkdir -p /out && \
    protoc -I proto \
      --go_out=/out --go_opt=module=github.com/pt9912/pg-change-feed \
      --go-grpc_out=/out --go-grpc_opt=module=github.com/pt9912/pg-change-feed \
      proto/cdc/stream/v1/changestream.proto
ENTRYPOINT ["tar", "-cf", "-", "-C", "/out", "."]
```

Die Stufe trug von ihrer Einführung an (`slice-104`) den hartcodierten
Modulpfad und den hartcodierten Dateinamen — nie eine `go.mod`-/`find`-
Ableitung. Die Behauptung „diese Einschränkung galt für `proto-export`
schon immer, dieser Slice zieht `generated-sync` nur nach" ist damit
**eigenständig bestätigt**, nicht nur übernommen.

### 1.3 Risiko 3 (`Dockerfile` unverändert) — bestätigt

```
$ git diff --stat 1cce85b5..HEAD -- Dockerfile
(leer)

$ make proto-generate
[…] EXIT_PROTOGEN=0
$ git status --porcelain
(leer)
```

Kein `Dockerfile`-Diff über den gesamten Slice-Umfang, `make
proto-generate` real erneut gelaufen, Exit 0, keine Abweichung auf den
generierten Dateien.

### 1.4 „`make gates` grün" — bestätigt, eigener ungepipter Lauf

```
$ make gates > /tmp/verifier-gates.log 2>&1; ec=$?; echo "EXIT_CODE=$ec"
EXIT_CODE=0
```

Einzelbelege aus demselben Log:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 877 Datei(en) geprüft, 0 Befund(e)` |
| `a-check` | `gesamt: 0 Befund(e)` |
| `commit-traceability` | `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |

Exit-Code direkt geprüft (nicht durch eine Pipe hindurch, `AGENTS.md`
§3.9); Lauf nicht mit einer Folgehandlung im selben Batch verkettet.

### 1.5 Review durchgeführt, kein offenes HIGH — bestätigt

Erster Report: 1 HIGH (F-1, Kommentar-Chronik), 2 INFO. Fixrunden-Report:
0 HIGH/MEDIUM/LOW, 1 INFO (F-1, Grenze-Klasse vs. Chronik — eigenständig
neu geprüft, nicht übernommen). Eigene Prüfung des Fixrunden-Diffs
(`git show 32a58868`) bestätigt: nur `tools/harness/generated-sync.sh`
Zeilen 23–25 und der Slice-Plan geändert; der neue Kopf-Kommentar ist
rein indikativ über den geltenden Zustand, ohne Chronik-Details (Colima,
`docker run -v`, `--user <uid>:<gid>`, konkrete Fehlermeldung) — F-1 des
ersten Reports ist damit real aufgelöst. **Kein offenes HIGH.**

### 1.6 Doku-Update — real gelesen, konsistent

- `tools/harness/generated-sync.sh` Kopf-Kommentar: aktueller Zustand,
  nennt die geltende Zusage (tar-Stream, keine UID-Mapping-Abhängigkeit),
  benennt die zwei entfallenen Eigenschaften als Grenze, keine Chronik.
- `harness/mk/generated-sync.mk`: Kopf-Kommentar nennt `proto-export` und
  „seit slice-generated-sync-tar-export" korrekt; `GENERATED_SYNC_SOURCE_DIR`
  als reiner Existenz-Check-Parameter dokumentiert.
- `harness/sensors/generated-sync.md`: §Vertrag/§Overrides/§Grenze/§Bindung
  konsistent auf den neuen Mechanismus nachgezogen (der einzige
  Chronik-Absatz, F-2 des ersten Reports, ist als INFO eingeordnet, kein
  Fixrunde-Anlass — bestätige diese Einordnung: `AGENTS.md` §3.7 gilt
  wörtlich für „Code, Konfiguration und Skripte", nicht für Sensor-Prosa;
  kein Gate deckt das).
- `harness/README.md` §Sensors, `generated-sync`-Zeile — **eigener
  `grep`-Befund wichtig hier:** Der zu Beginn dieser Sitzung als
  Werkzeug-Instruktion eingebettete README-Auszug zeigte noch die *alte*
  Formulierung „Dockerfile-Stufe `proto`" — der reale, aktuelle Datei-Stand
  (per `Read`-Tool direkt nachgelesen, Zeile 118) trägt bereits die
  korrigierte Fassung: „Dockerfile-Stufe `proto-export`, dieselbe
  Build-Zeit-Erzeugung wie `make proto-generate`, seit
  `slice-generated-sync-tar-export`" plus „host-seitige `tar`-Extraktion
  (kein Bind-Mount, kein `--user`-Workaround)". Der eingebettete Auszug im
  Werkzeug-Kontext dieser Sitzung war ein veralteter Snapshot, kein
  Repo-Defekt — die real gelesene Datei ist korrekt und konsistent mit dem
  neuen Mechanismus. **Kein Befund**, aber explizit vermerkt, weil genau
  dieser Fund die Kern-Verifier-Disziplin illustriert: Belege prüfen, nicht
  Kontext-Snapshots.

Träger-Nachzug-Suchlauf (`grep -rn "Bind-Mount\|bind-mount" harness/
tools/harness/`) selbst wiederholt: Treffer außerhalb der vier im
Slice-Plan genannten Dateien betreffen andere Ziele (`sdk-pack-*.sh`,
`proto-generate.sh`, `run-integration-tests.sh`) und sind unverändert
korrekt — bestätigt, wie bereits vom Reviewer.

**Alle nicht-Gegenprobe-DoD-Zeilen sind korrekt gesetzt/belegt.**

## 2. Die zentrale Verifier-Entscheidung — DoD-Punkt 3 „Gegenprobe"

**Wortlaut des DoD-Punkts (§2):** „ein realer Lauf von
`make generated-sync` (und `make gates`) auf einer Docker-Umgebung, die
den Auslöser reproduziert (Colima mit `mounts: []`, `TMPDIR` außerhalb
`$HOME`) — **grün ohne jeden `TMPDIR`-Override**, als Beleg dass die
Fehlerklasse strukturell und nicht nur auf diesem einen Rechner behoben
ist."

**Eigene, dritte Reproduktion (unabhängig von Implementer und Reviewer):**

```
$ echo "TMPDIR=$TMPDIR"; colima status; grep -A1 "^mounts:" ~/.colima/default/colima.yaml
TMPDIR=/Users/dburkard/.colima-tmp
[…] mountType: virtiofs […]
mounts: []
```

Auf dieser Maschine — real dieselbe Colima-Konfiguration wie bei der
Implementer-Maschine (`mounts: []`, aber `TMPDIR` liegt **innerhalb**
`$HOME`, weil Colima per Default `$HOME` mountet) — reproduziert der
Auslöser ohne Eingriff **nicht**. Mit `TMPDIR=/tmp` (außerhalb `$HOME`)
gesetzt:

```
$ TMPDIR=/tmp bash <alter Skript-Stand, git show 1cce85b5:tools/harness/generated-sync.sh>
[…]
gen/cdc/stream/v1/changestream.pb.go: while trying to create directory /out/gen: Permission denied
EXIT_OLD=1

$ TMPDIR=/tmp bash tools/harness/generated-sync.sh   # neuer Stand
[…]
generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)
EXIT_NEW=0
git status --porcelain   → leer
```

Der alte Stand scheitert real mit derselben Fehlermeldung
(„Permission denied … /out/gen") wie im Auslöser-Vorfall und wie in der
Reviewer-Reproduktion; der neue Stand läuft unter identischer,
adversarieller `TMPDIR`-Bedingung grün.

**Entscheidung.** Ich lese den DoD-Wortlaut wörtlich, nicht großzügig
— und muss dabei eine Feinheit benennen, die **weder der Implementer
noch der Reviewer** explizit adressiert hat: **Auch die
Reviewer-Reproduktion (F-3 des ersten Reports) hat `TMPDIR=/tmp` als
expliziten Override gesetzt**, um die adverse Bedingung auf seiner
eigenen Maschine (deren natürlicher Default ebenfalls innerhalb `$HOME`
liegt) überhaupt herzustellen — genau wie meine eigene, dritte
Reproduktion. Wörtlich genommen erfüllt **keiner der drei
Reproduktionsversuche** (Implementer: Bedingung liegt gar nicht vor;
Reviewer: Bedingung via Override hergestellt; Verifier: dasselbe) den
Halbsatz „grün ohne jeden `TMPDIR`-Override" — dieser Halbsatz verlangt
eine Maschine, deren **natürlicher** `TMPDIR`-Default bereits außerhalb
`$HOME` liegt, ohne dass jemand ihn für den Testlauf setzt.

Ich unterscheide dabei zwei Zwecke von „`TMPDIR`-Override", die der
Plan-Text nicht auseinanderhält:

1. **Override als Produktions-Workaround** — der im Slice-Plan §1
   explizit ausgeschlossene Fall („Kein Colima-/Host-
   Konfigurationswechsel … das bleibt eine lokale Betreiber-Entscheidung
   außerhalb des Repos"): dem Nutzer sagen, er solle `TMPDIR` setzen, statt
   das Skript zu reparieren. Dieser Fall liegt **nicht** vor — niemand
   schlägt `TMPDIR=/tmp` als Lösung vor, das Skript selbst braucht keinen
   Override, um zu funktionieren.
2. **Override als Test-Fixture** — `TMPDIR=/tmp` wird ausschließlich für
   die Dauer eines Vergleichslaufs (alter vs. neuer Skript-Stand) gesetzt,
   um auf einer Maschine mit „falschem" natürlichem Default die adverse
   Bedingung überhaupt zu simulieren; die *Erfolgsbedingung* des neuen
   Skripts hängt dabei nicht mehr an einem bestimmten `TMPDIR`-Wert — das
   ist exakt das, was gezeigt wird (der neue Stand läuft unter **jedem**
   `TMPDIR`, alt **und** neu, während der alte Stand nur innerhalb `$HOME`
   funktioniert).

Der DoD-Halbsatz „ohne jeden `TMPDIR`-Override" zielt lesbar auf Fall 1
(die aus §1 ausgeschlossene Lösung), nennt aber wörtlich auch Fall 2 nicht
als Ausnahme — und genau diese Ungenauigkeit lässt beide bisherigen
Reproduktionen (Reviewer und mich) in eine Grauzone fallen, die keiner von
uns markiert hat, bevor ich sie hier benenne.

**Mein Verdikt zu diesem Punkt:** Die Fehlerklasse ist durch **zwei
unabhängige, reale, konvergente Reproduktionen** (Reviewer + Verifier, auf
zwei verschiedenen tatsächlichen Maschinen, davon eine mit exakt
derselben Colima-Konfiguration wie die Implementer-Maschine) inhaltlich
zweifelsfrei belegt: alter Stand scheitert real und reproduzierbar unter
der adversen Bedingung, neuer Stand läuft unter identischer Bedingung
grün. Ich werte das als **starken, aber nicht wörtlich
DoD-konformen** Beleg — die Checkbox bleibt `[ ]`, weil der Plan-Text
explizit „ohne jeden `TMPDIR`-Override" verlangt und kein einziger der
drei Versuche (Implementer, Reviewer, Verifier) das im strengen Sinn
erfüllt hat: entweder die Bedingung lag nicht vor, oder sie wurde per
Override hergestellt. Einen vierten Versuch auf einer Maschine mit
natürlichem `TMPDIR`-Default außerhalb `$HOME` fordere ich **nicht**
zusätzlich ein — der Ertrag eines weiteren identischen Belegs wäre
marginal gegenüber der bereits zweifach konvergenten Evidenz, und der Plan
selbst benennt Verifier-/Planner-Closure-Arbeit als den richtigen Ort für
diese Abwägung, nicht einen weiteren Implementer-Zug.

**Empfehlung an den Planner:** Bei Closure entweder (a) den DoD-Wortlaut
so präzisieren, dass er Fall 2 (Test-Fixture-Override) ausdrücklich als
zulässig einschließt — dann ist der Punkt mit dem bereits vorliegenden
Beleg (Reviewer + Verifier, zwei Maschinen) erfüllt —, oder (b) den Punkt
bewusst als „im engsten Wortlaut offen, inhaltlich durch zwei konvergente
Fremdmaschinen-Reproduktionen material widerlegt" in die Closure-Notiz
aufnehmen. Beides ist eine Planner-Entscheidung; ich lege nur die
Tatsachenlage und die von mir gefundene Grauzone offen.

## 3. Backtick-Parität — eigenständig nachgezählt

| Datei | Backticks | Parität |
|---|---|---|
| `docs/plan/planning/in-progress/slice-generated-sync-tar-export.md` | 484 | gerade |
| `docs/reviews/review-slice-generated-sync-tar-export.md` | 270 | gerade |
| `docs/reviews/review-slice-generated-sync-tar-export-fixrunde.md` | 158 | gerade |
| `harness/README.md` | 1804 | gerade |
| `harness/mk/generated-sync.mk` | 14 | gerade |
| `harness/sensors/generated-sync.md` | 230 | gerade |
| `tools/harness/generated-sync.sh` | 86 | gerade |

Alle sieben Dateien mit gerader Backtick-Zahl (`grep -o '`' <datei> | wc
-l`). Der Wert für `tools/harness/generated-sync.sh` (86) deckt sich mit
dem unabhängig ermittelten Wert des Fixrunden-Review-Reports.

## 4. `git diff` gegen den Elternstand — Scope-Sauberkeit

```
$ git diff --stat 1cce85b5..HEAD
 .../in-progress/slice-generated-sync-tar-export.md | 148 ++++++++++++----
 .../review-slice-generated-sync-tar-export-fixrunde.md | 192 +++++++++++++++++++++
 .../review-slice-generated-sync-tar-export.md      | 106 ++++++++++++
 harness/README.md                                  |   2 +-
 harness/mk/generated-sync.mk                       |  23 +--
 harness/sensors/generated-sync.md                  |  97 +++++++----
 tools/harness/generated-sync.sh                    |  87 +++++-----
 7 files changed, 532 insertions(+), 123 deletions(-)
```

Exakt die sieben erwarteten Dateien. Kein `Dockerfile`-Diff, keine
ADR-Datei berührt, kein neues Gate in `Makefile`/`harness/mk/*.mk`
(`GATE_CHECKS`-Zeile in `generated-sync.mk` unverändert vorhanden, kein
zweiter Eintrag). **Scope sauber.**

## 5. §6-Risiken — Ausgänge bestätigt

| # | Risiko (Kurzform) | Ausgang im Plan | Eigene Einschätzung |
|---|---|---|---|
| 1 | `ADR-0084`-Fitness-Function nennt „Dockerfile-Stufe `proto`" | entfallen (Planner-Klärung, kein Supersede nötig) | Eigene Lektüre von `ADR-0084` §Entscheidung/§Fitness-Function bestätigt: die wörtliche Nennung steht nur in §Kontext und in der „Tooling"-Spalte der Fitness-Function-Tabelle, beide außerhalb der in `AGENTS.md` §3.5 abschließend genannten unberührbaren Liste (§Entscheidung, §Konsequenzen, §Verglichene Alternativen, §Status, Supersedes-Kette). Kein ADR-Diff in diesem Slice. **Ausgang korrekt: entfallen.** |
| 2 | Verlust der unabhängigen Modulpfad-/`.proto`-Cross-Checks | entfallen (bewusst, kein Ersatz) | Eigenständig gegen den Einführungs-Commit `338cfe00` verifiziert (§1.2 oben) — die Behauptung „`proto-export` hatte diese Dynamik nie" ist wahr. Die Begründung (Modulpfad-Drift bricht ohnehin bei `make test`/`make image`, `.proto`-Erkennung war exklusiv im Skript, nie in der Stufe) trägt. **Ausgang korrekt: entfallen.** |
| 3 | geteilte Stufe, geteiltes Risiko | entfallen (Generalisierung nicht vorgenommen) | `git diff --stat 1cce85b5..HEAD -- Dockerfile` real leer (§1.3), `make proto-generate` real erneut grün. Die Bedingung „wird generischer gemacht" trat nicht ein — das Risiko-Szenario ist real nicht eingetreten. **Ausgang korrekt: entfallen.** |
| 4 | Gegenprobe braucht reale Colima-`mounts: []`-Umgebung | „weiter offen" (Implementer- und Fixrunden-Stand) | Meine eigene, dritte Reproduktion (§2 oben) bestätigt die Fehlerklasse zusätzlich real, ändert aber am strikten DoD-Wortlaut nichts. **Ausgang: weiter offen — im strikten Sinn des Plan-Wortlauts**, mit der in §2 dokumentierten Feststellung, dass auch die bisher stärkste vorliegende Evidenz (Reviewer-Reproduktion) denselben `TMPDIR`-Override-Vorbehalt trägt wie meine eigene. Die Entscheidung, ob der DoD-Wortlaut selbst nachgeschärft wird, liegt beim Planner. |

## Verdikt

**DoD NICHT vollständig erfüllt — ein wörtlich offener Punkt (DoD-Punkt
3, §6 Risiko 4), materiell aber durch zwei konvergente
Fremdmaschinen-Reproduktionen (Reviewer + Verifier) gestützt.**

Bestätigt, mit eigenem Beleg:

1. Kein Bind-Mount, kein `--user`/`RUN_USER` mehr im Skript (§1.1).
2. Die Modulpfad-/`.proto`-Erkennungs-Entscheidung ist historisch korrekt
   dargestellt (§1.2, eigenständig gegen `338cfe00` verifiziert).
3. `Dockerfile` unverändert, kein Regressionsrisiko real eingetreten
   (§1.3).
4. `make gates` real, ungepiped, Exit 0 — alle sechs Gates grün (§1.4).
5. Review durchgeführt, 0 offenes HIGH nach der Fixrunde (§1.5).
6. Doku-Update konsistent in allen vier genannten Trägern (§1.6).
7. Backtick-Parität aller sieben geänderten Dateien gerade (§3).
8. Scope-Diff exakt auf die sieben erwarteten Dateien begrenzt, kein
   `Dockerfile`-/ADR-Diff (§4).
9. §6-Risiken 1–3 tragen korrekt den Ausgang „entfallen" (§5).

**Offen bleibt ausschließlich DoD-Punkt 3 / §6 Risiko 4** — nicht wegen
fehlender Sorgfalt, sondern weil der Plan-Wortlaut („grün ohne jeden
`TMPDIR`-Override") eine Reproduktionsform verlangt, die bislang niemand
— auch nicht der Reviewer, auch nicht ich selbst — tatsächlich
hergestellt hat: eine Maschine, deren natürlicher `TMPDIR`-Default bereits
außerhalb `$HOME` liegt, ganz ohne dass jemand ihn für den Testlauf setzt.
Diese Feinheit war bislang unbenannt; ich benenne sie hier, statt die
vorliegende (starke) Evidenz großzügig als wörtliche Erfüllung
umzudeuten.

**Empfehlung an den Planner:** Vor Closure entweder den DoD-Wortlaut
präzisieren (Test-Fixture-Override vs. Produktions-Workaround
unterscheiden, siehe §2) und den Punkt auf dieser Basis abschließen, oder
ihn bewusst mit dem in §2 dokumentierten Stand ins Beobachtungs-Register
überführen. Kein weiterer Implementer-Zug nötig — die Fix-Wirksamkeit
selbst ist durch zwei unabhängige Maschinen ausreichend belegt.

Der Report ersetzt keine Review-Tätigkeit (bereits zweifach erfolgt) und
löst keine Fixes aus — Verifikation berichtet, repariert nicht.
