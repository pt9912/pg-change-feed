# Slice generated-sync-tar-export: `make generated-sync` ohne Bind-Mount (Wiederverwendung der `proto-export`-Erzeugung)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — kein Closure-Kriterium jenseits der eigenen DoD.

**Bezug:** [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
(Festlegung 1 — Sync-Gate für den Protobuf-Code, „erzeugt in ein
Temp-Verzeichnis … vergleicht dort", **nicht** festgelegt: die konkrete
Erzeugungs-Mechanik), [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(der Generator selbst). Dieser Slice ändert **keine** `Accepted`-ADR-Entscheidung
inhaltlich, sondern liefert Werkzeug-Verhalten innerhalb der dort getroffenen
Entscheidung — siehe §6 Risiko 1 zur offenen Frage, ob die
Fitness-Function-Tabelle von `ADR-0084` davon berührt ist.

**Berührte Spec-Stellen:** — (reine Werkzeug-/Sensor-Mechanik, kein
Lastenheft-/Pflichtenheft-/Architektur-Bezugspunkt; `make generated-sync`
selbst trägt keine `SPEC-*`/`ARC-*`-Kennung, ebenso wie `ADR-0084` selbst
kein Spec-Stratum berührt).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** Claude Code (direkt vom Nutzer beauftragt, Anlass: ein realer
`make gates`-Fehlschlag von `generated-sync` auf einem Rechner mit
Colima-`mounts: []` — `docker run --user <uid>:<gid> -v <TMPDIR-Pfad>:/out`
liefert dort `Permission denied`, weil das UID-Mapping außerhalb von `$HOME`
nicht greift; Diagnose real reproduziert). **Datum:** 2026-09-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** `tools/harness/generated-sync.sh` erzeugt sein Vergleichs-Erzeugnis
ohne `docker run -v .../out:/out`-Bind-Mount — durch Wiederverwendung
derselben Build-Zeit-Erzeugung, die `make proto-generate` bereits nutzt
(Dockerfile-Stufe `proto-export`, `ENTRYPOINT ["tar", "-cf", "-", "-C",
"/out", "."]`, host-seitige Extraktion `docker run --rm --network none
<image> | tar -x -C <Temp-Verzeichnis>`), statt selbst `protoc` mit
manuellen Bind-Mount-Argumenten aufzurufen. Der Bind-Mount ist die
Fehlerquelle des Auslösers (Colima-UID-Mapping außerhalb `$HOME`) — die
tar-Stream-Extraktion braucht keinen Mount und ist damit strukturell immun
gegen diese Fehlerklasse, nicht nur auf diesem einen Rechner umschifft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Umbau von `make proto-generate` selbst in seinem beobachtbaren
  Verhalten** — es bleibt bei „ein `.proto`, ein Modulpfad, Erzeugnis
  identisch". Falls sich beim Schreiben zeigt, dass die `proto-export`-Stufe
  generischer werden muss (dynamische `.proto`-Erkennung/Modulpfad-Ableitung
  **innerhalb** der Stufe, siehe §3 und §6 Risiko 3), ist das eine
  Verallgemeinerung der Stufe, kein Verhaltenswechsel von
  `make proto-generate`.
- **Keine inhaltliche Änderung an `ADR-0084` oder `ADR-0060`** — beide sind
  `Accepted` und nach `AGENTS.md` §3.5 immutabel. Erweist sich die
  Fitness-Function-Tabelle von `ADR-0084` (nennt „Dockerfile-Stufe `proto`")
  durch den Stufen-Wechsel als inhaltlich berührt, ist eine Supersede-ADR
  ein **eigener, vorgelagerter Architect-Zug** — nicht Bestandteil dieses
  Implementer-Slice (siehe §6 Risiko 1).
- **Kein Colima-/Host-Konfigurationswechsel** (`TMPDIR`-Override,
  `mounts:`-Eintrag in `~/.colima/default/colima.yaml`) — das bleibt eine
  lokale Betreiber-Entscheidung außerhalb des Repos; dieser Slice löst das
  Problem strukturell im Sensor selbst, unabhängig vom Docker-Backend des
  Aufrufers.
- **Kein neues Gate und keine Schwellen-Änderung** — `generated-sync` bleibt
  unverändert in `GATE_CHECKS`, nur seine interne Erzeugungs-Mechanik ändert
  sich (`AGENTS.md` §3.6 unberührt).

## 2. Definition of Done

- [ ] `make generated-sync` läuft ohne jeden `docker run -v`-Bind-Mount —
      Erzeugung über dieselbe `proto-export`-Stufe/`tar`-Extraktion wie
      `make proto-generate`, Vergleich weiterhin in einem Temp-Verzeichnis,
      Arbeitsbaum bleibt unverändert (Fitness-Kriterium `ADR-0084`
      Festlegung 1 unverändert erfüllt).
- [ ] Die beiden bisherigen Eigenschaften — unabhängige Modulpfad-Ableitung
      aus `go.mod` und dynamische `.proto`-Dateierkennung (`find … -name
      '*.proto'`) — bleiben erhalten, **oder** ihr Wegfall ist im Bericht
      explizit benannt und begründet (kein stiller Verlust, siehe §6
      Risiko 2).
- [ ] Gegenprobe: ein realer Lauf von `make generated-sync` (und `make
      gates`) auf einer Docker-Umgebung, die den Auslöser reproduziert
      (Colima mit `mounts: []`, `TMPDIR` außerhalb `$HOME`) — grün ohne
      jeden `TMPDIR`-Override, als Beleg dass die Fehlerklasse
      strukturell und nicht nur auf diesem einen Rechner behoben ist.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `tools/harness/generated-sync.sh` Kopf-Kommentar (nennt
      aktuell den Bind-Mount-Mechanismus explizit, Zeilen 11–22) und ggf.
      `harness/README.md` §Sensors (`generated-sync`-Zeile), falls die
      Beschreibung dort den alten Mechanismus wiederholt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      Repo ohne Wellen-Betrieb für diesen (wellenlosen) Slice, hier direkt
      geprüft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/generated-sync.sh` | update | `docker build --target proto` + manueller `protoc`-Aufruf mit Bind-Mount (Zeilen 70–76) ersetzt durch `docker build --target proto-export` + `docker run --rm --network none <image> \| tar -x -C "$out_dir"` (Muster `tools/harness/proto-generate.sh` Zeile 31f, aber Extraktion in ein Temp-Verzeichnis statt `.`). `RUN_USER`/`--user`-Logik entfällt (kein Mount mehr, dieselbe Begründung wie beim seinerzeitigen `proto-generate`-Umbau in `slice-104`). |
| `Dockerfile` (Stufe `proto-export`) | ggf. update | Nur falls die Verallgemeinerung nötig ist (siehe §6 Risiko 2/3): der `RUN`-Schritt (Zeile 62–66) leitet Modulpfad und `.proto`-Dateiliste selbst ab (`awk` gegen das bereits kopierte `go.mod`, `find proto -name '*.proto'`) statt sie hartcodiert zu tragen — dieselbe Dynamik, die `generated-sync.sh` heute im Skript trägt, wandert in die Stufe, die jetzt **beide** Ziele (`proto-generate` und `generated-sync`) gemeinsam nutzen. |
| `tools/harness/generated-sync.sh` (Kopf-Kommentar) | update | Beschreibung des Mechanismus (aktuell: „Der Generator … schreibt in ein Temp-Verzeichnis, der Baum haengt als `:ro`-Bind-Mount im Container", Zeilen 11–14) auf den neuen, mount-losen Mechanismus nachgezogen — `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist). |
| `harness/README.md` §Sensors (`generated-sync`-Zeile), `harness/mk/generated-sync.mk` (Kopf-Kommentar) | ggf. update | Falls dort der Bind-Mount-Mechanismus explizit genannt ist (Trägernachzug, `AGENTS.md` §3.13) — Implementer-Suchlauf (`grep -rn "Bind-Mount\|bind-mount" harness/ tools/harness/`) vor Abschluss. |
| Testdatei | — | Kein eigenständiges Unit-Test-Ziel; der Beleg ist der reale `make generated-sync`-Lauf selbst (Sensor-Skript, kein Go-Testpaket) — analog dem bestehenden Muster für `tools/harness/proto-generate.sh`. |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn priorisiert (Verantwortlich
gesetzt) — kein weiterer Slice als Vorbedingung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Schreiben zeigt, dass die Dockerfile-Verallgemeinerung (§6 Risiko 3)
  selbst ein größerer, riskanter Umbau der von zwei Zielen geteilten Stufe
  wird — dann zurück zur Zerlegung (ein Slice für die Stufen-Generalisierung,
  ein zweiter für die Skript-Umstellung).
- `in-progress` → `open` (blockiert — Carveout?): falls die Architect-Rolle
  (§6 Risiko 1) entscheidet, dass der Stufen-Wechsel eine Supersede-ADR zu
  `ADR-0084` braucht — dann blockiert bis diese ADR `Accepted` ist.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reale Gegenprobe gegen den Auslöser
(DoD-Punkt 3) bestanden + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Risiko 1 — [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)-Berührung.**
  Die Fitness-Function-Tabelle dieser ADR nennt
  „Dockerfile-Stufe `proto`" als Erzeugungsmechanismus für `generated-sync`.
  Ein Wechsel auf `proto-export` könnte eine inhaltliche Änderung einer
  `Accepted`-ADR sein (`AGENTS.md` §3.5 — unberührbar sind §Entscheidung,
  §Konsequenzen, §Verglichene Alternativen, §Status; die
  Fitness-Function-Tabelle ist als Teil der Entscheidungsdarstellung zu
  behandeln, nicht als reine Zitat-Korrektur nach `ADR-0073`). **Ausgang:**
  weiter offen — vor der Implementierung klärt der Implementer-Zug mit der
  Architect-Rolle, ob eine Supersede-ADR zu `ADR-0084` nötig ist (die
  Kern-Entscheidung „Temp-Verzeichnis, kein Baum-Schreiben, Diff im Befund"
  bleibt unverändert erfüllt — nur die Fitness-Function-Zeile nennt einen
  anderen Stufennamen).
- **Risiko 2 — Verlust der unabhängigen Modulpfad-Cross-Check.** Aktuell
  leitet das Skript den Modulpfad selbst aus `go.mod` ab (`generated-sync.sh`
  Zeile 47) und vergleicht implizit gegen den im Dockerfile hartcodierten
  Wert (`module=github.com/pt9912/pg-change-feed`, `Dockerfile` Zeile 64+65)
  — ein Drift zwischen beiden fiele heute auf. Verschiebt sich die Ableitung
  vollständig in die Dockerfile-Stufe (beide Ziele lesen dieselbe Quelle),
  entfällt dieser Cross-Check. **Ausgang:** zu entscheiden beim Schreiben —
  entfällt bewusst (Begründung: ein Modulpfadwechsel fiele ohnehin an
  anderer Stelle auf, `make test`/`make image` bräche ohne den korrekten
  Importpfad) oder wird durch einen expliziten, weiterhin unabhängigen
  Zweit-Check ersetzt (z. B. das Skript prüft nach dem Build zusätzlich,
  dass die generierten Dateien tatsächlich unter dem aus `go.mod`
  abgeleiteten Pfad liegen).
- **Risiko 3 — geteilte Stufe, geteiltes Risiko.** Wird die
  `proto-export`-Stufe generischer gemacht (dynamische `.proto`-Erkennung,
  `go.mod`-Ableitung im `RUN`-Schritt statt fester Argumente), nutzen
  `make proto-generate` **und** `make generated-sync` künftig exakt dieselbe
  Stufe — ein Fehler in der Verallgemeinerung träfe beide Ziele gleichzeitig.
  **Ausgang:** abzudecken durch einen Regressionsbeleg im Bericht: nach dem
  Umbau `make proto-generate` real erneut laufen lassen, Diff gegen den
  committeten Stand muss leer bleiben (kein unbeabsichtigter
  Verhaltenswechsel für das bestehende Ziel).
- **Risiko 4 — Gegenprobe braucht eine reale Colima-Umgebung mit
  `mounts: []`.** Der DoD-Punkt „Gegenprobe" (§2) ist nur auf einem Rechner
  mit genau dieser Docker-Backend-Konfiguration direkt beobachtbar. **Ausgang:**
  weiter offen bis zum realen Lauf — ersatzweise (falls kein solcher
  Rechner verfügbar ist) mindestens der Nachweis, dass `generated-sync.sh`
  nach dem Umbau keine `docker run -v`-Zeile mehr enthält (`grep -n "docker
  run" tools/harness/generated-sync.sh` zeigt keine `-v`-Option mehr) als
  schwächerer, aber netzloser Ersatzbeleg.

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird beim Abschluss ergänzt>
- **Was ging anders als geplant:** <wird beim Abschluss ergänzt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss ergänzt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss
  ergänzt>
- **Folge-Slices:** keine aus diesem Slice selbst erwartet.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <wird beim Abschluss ergänzt — Repo ohne
  Wellen-Betrieb für diesen Slice, hier direkt geprüft>

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `tools/harness/` +
`Dockerfile` (Sensor-/Generator-Mechanik) — bereits mehrfach berührt
(`generated-sync.sh`, `proto-generate.sh`, die Dockerfile-Stufen `proto`/
`proto-export`), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`grep`-Suche nach `Bind-Mount`/`generated-sync`/`proto-export`/`Colima` über
`docs/plan/planning/observations/**`) — kein Treffer; diese Beobachtung
(Colima-UID-Mapping außerhalb `$HOME`) ist neu und wird bei Closure als
eigener Eintrag angelegt (§2 DoD-Punkt Beobachtungs-Register), falls sie sich
über diesen einen Vorfall hinaus als wiederkehrend erweist — beim ersten
Auftreten reicht die Nennung in diesem Slice-Plan.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
