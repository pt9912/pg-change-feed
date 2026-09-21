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

**Verantwortlich:** Implementer-Agent, 2026-09-21.

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

- [x] `make generated-sync` läuft ohne jeden `docker run -v`-Bind-Mount —
      Erzeugung über dieselbe `proto-export`-Stufe/`tar`-Extraktion wie
      `make proto-generate`, Vergleich weiterhin in einem Temp-Verzeichnis,
      Arbeitsbaum bleibt unverändert (Fitness-Kriterium `ADR-0084`
      Festlegung 1 unverändert erfüllt). Real geprüft: `grep -n "docker run"
      tools/harness/generated-sync.sh` zeigt nur noch eine Aufruf-Zeile ohne
      `-v`; `make generated-sync` läuft grün, `git status --porcelain` danach
      leer.
- [x] Die beiden bisherigen Eigenschaften — unabhängige Modulpfad-Ableitung
      aus `go.mod` und dynamische `.proto`-Dateierkennung (`find … -name
      '*.proto'`) — Wegfall ist im Bericht (Implementer-Zug
      2026-09-21) explizit benannt und begründet: **beide entfallen
      bewusst**, siehe §6 Risiko 2 Ausgang.
- [x] Gegenprobe: ein realer Lauf von `make generated-sync` (und `make
      gates`) auf einer Docker-Umgebung, die den Auslöser reproduziert
      (Colima mit `mounts: []`, `TMPDIR` außerhalb `$HOME`) — grün ohne
      jeden `TMPDIR`-Override, als Beleg dass die Fehlerklasse
      strukturell und nicht nur auf diesem einen Rechner behoben ist.
      **Nicht erreicht** auf der Implementer-Maschine (§6 Risiko 4 Ausgang) —
      nur der schwächere, netzlose Ersatzbeleg geliefert; bleibt bis zu einem
      realen Repro-Lauf offen. **Nachtrag (Fixrunde F-1, 2026-09-21):** der
      Reviewer (`docs/reviews/review-slice-generated-sync-tar-export.md`
      F-3) hat auf seiner eigenen Maschine die exakte Auslöser-Kombination
      (Colima `mounts: []` **und** `TMPDIR=/tmp` außerhalb `$HOME`) real
      hergestellt und beide Skript-Stände dagegen laufen lassen: der **alte**
      Stand (`git show 1cce85b5:tools/harness/generated-sync.sh`) scheitert
      dort real mit „Permission denied" beim Schreiben nach `/out`, der
      **neue** Stand läuft unter identischem `TMPDIR=/tmp` grün
      (`generated-sync: OK`, `git status --porcelain` leer). Das ist ein
      direkterer Beleg als der oben genannte Ersatzbeleg — die
      Implementer-Maschine selbst reproduziert die Kombination weiterhin
      nicht (§6 Risiko 4 Ausgang bleibt unverändert stehen). Ob dieser
      Fremdmaschinen-Beleg die Checkbox trägt oder ein eigener Repro-Lauf
      nötig bleibt, entscheidet Verifier-/Planner-Closure-Arbeit.
      **Planner-Entscheidung bei Closure (2026-09-21, siehe §7 Closure-Notiz
      für die volle Begründung):** Checkbox auf `[x]` gesetzt. Der Verifier
      (`docs/reviews/verifikation-slice-generated-sync-tar-export.md` §2)
      hat eine dritte, unabhängige Reproduktion beigebracht (eigene
      Maschine, identische Colima-Konfiguration wie die
      Implementer-Maschine, `TMPDIR=/tmp`-Override) und dabei präzise
      benannt: **keiner** der drei Versuche (Implementer, Reviewer,
      Verifier) erfüllt den Wortlaut „ohne jeden `TMPDIR`-Override" streng
      wörtlich — bei Implementer lag die Bedingung gar nicht vor, bei
      Reviewer und Verifier wurde sie per Override hergestellt. Der Planner
      liest den Wortlaut-Zweck (§1 „Ausdrücklich NICHT" — kein
      Colima-/Host-Konfigurationswechsel **als ausgeliefertes
      Lösungs-Rezept**) und entscheidet: ein `TMPDIR`-Override, der
      ausschließlich als **Test-Fixture** die adverse Bedingung auf einer
      Maschine ohne sie herstellt, ist nicht der Fall, den der Wortlaut
      ausschließen wollte — ausgeschlossen war ein Override **als
      Produktions-Workaround** (dem Betreiber sagen, er solle `TMPDIR`
      setzen, statt das Skript zu reparieren). Dieser Fall liegt nicht vor:
      Weder Reviewer noch Verifier empfehlen `TMPDIR` als Lösung, beide
      zeigen den neuen Skript-Stand unter **beiden** Werten grün (Default
      innerhalb `$HOME` bei der Implementer-Maschine, `TMPDIR=/tmp` bei
      Reviewer/Verifier) — exakt der Beleg, dass der neue Mechanismus
      TMPDIR-wert-unabhängig ist, während der alte Stand unter identischer
      `TMPDIR=/tmp`-Bedingung real und reproduzierbar scheitert (zwei
      unabhängige Maschinen, konvergentes Ergebnis). Gegenprobe: bestanden.
- [x] `make gates` grün. Real gelaufen 2026-09-21, Exit 0 (siehe Bericht).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      `docs/reviews/review-slice-generated-sync-tar-export.md`: 1 HIGH
      (F-1, Kommentar-Chronik in `tools/harness/generated-sync.sh`) —
      Fixrunde 2026-09-21 behoben (Kopf-Kommentar auf die geltende Zusage
      gekürzt, Chronik entfernt); F-2/F-3 (beide INFO) sind kein
      Fixrunde-Anlass. Kein offenes HIGH mehr.
- [x] Doku-Update: `tools/harness/generated-sync.sh` Kopf-Kommentar,
      `harness/mk/generated-sync.mk` Kopf-Kommentar,
      `harness/sensors/generated-sync.md` und `harness/README.md` §Sensors
      (`generated-sync`-Zeile) auf den neuen, mount-losen Mechanismus
      nachgezogen (Träger-Nachzug-Suchlauf `grep -rn "Bind-Mount\|bind-mount"
      harness/ tools/harness/` lief, Treffer außerhalb dieser vier Dateien
      betreffen andere Ziele/Skripte und sind unverändert korrekt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. — Planner-Closure-Arbeit
      (Modul 8/Modul 5), siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — Beleg 8
      in `BEO-PGC/slice-chronik-in-code-kommentar/evidence/
      slice-generated-sync-tar-export.md` (F-1: erstes Vorkommen der
      Klasse in einem Shell-Skript; F-2: zweiter Beleg der Doku-Prosa-
      Grenze). Die Gegenprobe-Wortlaut-Frage (§7) erhält **keinen** neuen
      Registereintrag — Erstauftreten, per Closure aufgelöst statt offen
      gelassen; Begründung in §7. Siehe §7 für beide Entscheidungen im
      Detail.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen). — Risiko 1/2/3: entfallen (unverändert). Risiko 4:
      eingetreten, real behoben (Planner-Entscheidung, siehe DoD-Punkt 3
      oben und §7).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      Repo ohne Wellen-Betrieb für diesen (wellenlosen) Slice, hier direkt
      geprüft: Anker `ADR-0084`/`ADR-0060` (§Bezug, unverändert), keine
      Folge-Slices, Beobachtungs-Register siehe oben. Siehe §7.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/generated-sync.sh` | update | `docker build --target proto` + manueller `protoc`-Aufruf mit Bind-Mount (Zeilen 70–76) ersetzt durch `docker build --target proto-export` + `docker run --rm --network none <image> \| tar -x -C "$out_dir"` (Muster `tools/harness/proto-generate.sh` Zeile 31f, aber Extraktion in ein Temp-Verzeichnis statt `.`). `RUN_USER`/`--user`-Logik entfällt (kein Mount mehr, dieselbe Begründung wie beim seinerzeitigen `proto-generate`-Umbau in `slice-104`). |
| `Dockerfile` (Stufe `proto-export`) | ggf. update | Nur falls die Verallgemeinerung nötig ist (siehe §6 Risiko 2/3): der `RUN`-Schritt (Zeile 62–66) leitet Modulpfad und `.proto`-Dateiliste selbst ab (`awk` gegen das bereits kopierte `go.mod`, `find proto -name '*.proto'`) statt sie hartcodiert zu tragen — dieselbe Dynamik, die `generated-sync.sh` heute im Skript trägt, wandert in die Stufe, die jetzt **beide** Ziele (`proto-generate` und `generated-sync`) gemeinsam nutzen. |
| **Plan-Nachzug (Implementer-Zug, 2026-09-21):** `Dockerfile` | **nicht ausgeführt** | Die Zeile darüber ist konditional („ggf.", „nur falls nötig") — die Entscheidung beim Schreiben fällt auf **nicht generalisieren**: `proto-export` bleibt unverändert (ein `.proto`, ein hartcodierter Modulpfad, exakt wie vor diesem Slice). Begründung siehe §6 Risiko 2 Ausgang; die Alternative (Stufe generischer machen) hätte das Risiko-3-Szenario (geteilte Stufe, geteiltes Risiko, Regressionsbeleg nötig) real ausgelöst — vermieden, weil der Nutzen (zwei entfallende Cross-Checks wiederherstellen) den zusätzlichen Umbau-Umfang an einer von zwei Zielen gemeinsam genutzten Stufe nicht rechtfertigt. `git diff --stat Dockerfile` bleibt in diesem Lauf leer. |
| `tools/harness/generated-sync.sh` (Kopf-Kommentar) | update | Beschreibung des Mechanismus (aktuell: „Der Generator … schreibt in ein Temp-Verzeichnis, der Baum haengt als `:ro`-Bind-Mount im Container", Zeilen 11–14) auf den neuen, mount-losen Mechanismus nachgezogen — `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist). |
| `harness/README.md` §Sensors (`generated-sync`-Zeile), `harness/mk/generated-sync.mk` (Kopf-Kommentar) | ggf. update | Falls dort der Bind-Mount-Mechanismus explizit genannt ist (Trägernachzug, `AGENTS.md` §3.13) — Implementer-Suchlauf (`grep -rn "Bind-Mount\|bind-mount" harness/ tools/harness/`) vor Abschluss. |
| Testdatei | — | Kein eigenständiges Unit-Test-Ziel; der Beleg ist der reale `make generated-sync`-Lauf selbst (Sensor-Skript, kein Go-Testpaket) — analog dem bestehenden Muster für `tools/harness/proto-generate.sh`. |

**Plan-Nachzug (Planner-Zug, `next` → `in-progress`, 2026-09-21) — Klärung
Risiko 1 vorab:** Die Planner-Rolle hat
[`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) erneut
gegen `AGENTS.md` §3.5 gelesen und **entscheidet: keine Supersede-ADR
nötig.** Begründung:

- §Entscheidung selbst (der unberührbare Kern nach `AGENTS.md` §3.5) legt
  die Erzeugungs-Mechanik **nicht** fest — Festlegung 1 spricht durchgehend
  vom „gepinnten Generator" gegen die „committete `.proto`-Quelle", mit
  genau zwei Bedingungen: „das Gate schreibt den Arbeitsbaum nicht" (Temp-
  Verzeichnis) und „der Befund nennt den Diff". Beide Bedingungen bleiben
  nach dem Wechsel auf `proto-export` unverändert erfüllt — die tar-Stream-
  Extraktion landet weiterhin in einem Temp-Verzeichnis, nicht im Baum.
- Die wörtliche Nennung „Dockerfile-Stufe `proto`" steht an genau zwei
  Stellen: einmal in §Kontext (2) — dort beschreibt sie **`make
  proto-generate`s** Mechanik, nicht die von `generated-sync` — und einmal
  in der Fitness-Function-Tabelle (Spalte „Tooling"). Keine der beiden
  Stellen gehört zur in `AGENTS.md` §3.5 abschließend benannten
  unberührbaren Liste (§Entscheidung, §Konsequenzen, §Verglichene
  Alternativen, §Status, Supersedes-Kette).
- Die Fitness-Function-**Regel**-Spalte — die tatsächlich bindende Aussage
  („byte-gleich … erzeugt in ein Temp-Verzeichnis, verglichen mit dem Baum,
  Befund mit Diff") — bleibt vom Stufenwechsel unberührt; nur die
  „Tooling"-Spalte nennt danach einen anderen Stufennamen als den zum
  Entscheidungszeitpunkt vorgesehenen. Das ist eine veraltete
  Beleg-/Beispielangabe innerhalb einer nicht-unberührbaren Sektion, keine
  Änderung der Entscheidung selbst — anders als die vier in §Entscheidung
  unberührbaren Sektionen, die dieser Slice inhaltlich nicht anfasst.
- Dieser Slice ändert `ADR-0084` selbst **nicht** (kein Commit auf die
  ADR-Datei). Die Fitness-Function-Zeile bleibt wörtlich stehen und nennt
  weiterhin „Dockerfile-Stufe `proto`" — historisch korrekt als das zum
  Entscheidungszeitpunkt vorgesehene Tooling, nicht mehr deckungsgleich mit
  dem seit diesem Slice tatsächlich laufenden Mechanismus. Das ist benannt,
  nicht verdeckt: künftige Leser dieser ADR finden die Auflösung über den
  Verweis in diesem Slice-Plan und ggf. über einen
  späteren, eigenständigen `AGENTS.md` §3.5-konformen Nachzug (neue Zeile in
  `ADR-0084`s §Geschichte, kein Überschreiben — Entscheidung eines
  künftigen Zuges, nicht Bestandteil dieses Slice).
- **Kein Dispatch eines eigenen Architect-Zugs** — die Lektüre der
  einschlägigen Textstellen (`§Entscheidung` vs. `§Kontext`/Fitness-Function)
  ist eindeutig genug für eine begründete Planner-Entscheidung; ein
  Architect-Zug käme bei derselben Textgrundlage zur selben Schlussfolgerung.

**Folge für §4/§6:** Die Rückführung `in-progress` → `open` aus §4
(„Architect-Rolle entscheidet, dass der Stufen-Wechsel eine Supersede-ADR
braucht") entfällt damit als Vorbedingung — der Implementer-Zug startet
ohne diese Blockade. §6 Risiko 1 erhält seinen Ausgang **entfallen** mit
Verweis auf diesen Absatz (nicht erst bei Closure).

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
  **Entfallen** durch die Planner-Klärung vom 2026-09-21 (§3 „Plan-Nachzug",
  §6 Risiko 1) vor dem Übergang nach `in-progress` — diese Rückführung wird
  nicht mehr erwartet, bleibt hier stehen als Beleg, dass sie vorab benannt
  und dann aufgelöst wurde, statt stillschweigend zu verschwinden.

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
  entfallen (Planner-Klärung 2026-09-21, siehe §3 „Plan-Nachzug") — die
  wörtliche Nennung „Dockerfile-Stufe `proto`" steht ausschließlich in
  §Kontext (dort über `make proto-generate`, nicht über `generated-sync`)
  und in der Fitness-Function-„Tooling"-Spalte, beide außerhalb der in
  `AGENTS.md` §3.5 abschließend benannten unberührbaren Liste
  (§Entscheidung, §Konsequenzen, §Verglichene Alternativen, §Status,
  Supersedes-Kette); §Entscheidung selbst legt die Erzeugungs-Mechanik
  bewusst nicht fest und bleibt durch den Stufenwechsel unberührt. Dieser
  Slice ändert `ADR-0084` selbst nicht — die jetzt leicht veraltete
  „Tooling"-Nennung bleibt stehen und ist über diesen Slice-Plan auflösbar;
  kein Supersede-ADR nötig, kein Architect-Zug dispatcht.
- **Risiko 2 — Verlust der unabhängigen Modulpfad-Cross-Check.** Aktuell
  leitet das Skript den Modulpfad selbst aus `go.mod` ab (`generated-sync.sh`
  Zeile 47) und vergleicht implizit gegen den im Dockerfile hartcodierten
  Wert (`module=github.com/pt9912/pg-change-feed`, `Dockerfile` Zeile 64+65)
  — ein Drift zwischen beiden fiele heute auf. Verschiebt sich die Ableitung
  vollständig in die Dockerfile-Stufe (beide Ziele lesen dieselbe Quelle),
  entfällt dieser Cross-Check. **Ausgang: entfallen (Implementer-Zug
  2026-09-21).** Beide bisherigen Eigenschaften entfallen bewusst, kein
  Zweit-Check ersetzt sie:
  - Modulpfad-Cross-Check: entfällt mit der oben genannten Begründung — ein
    tatsächlicher `go.mod`-Modulpfadwechsel bricht die Importpfade im
    gesamten Baum und fällt damit bei `make test`/`make image` auf, auch
    ohne dass `generated-sync` ihn zusätzlich meldet.
  - Dynamische `.proto`-Dateierkennung: entfällt ebenfalls bewusst — diese
    Eigenschaft trug bisher **ausschließlich** `generated-sync.sh` (per
    `find`); die Dockerfile-Stufe `proto-export`, die `make proto-generate`
    bereits seit slice-104 nutzt, nennt ihre `.proto`-Datei bereits namentlich
    im `RUN`-Schritt und kannte diese Dynamik nie. Eine neue, zweite
    `.proto`-Datei würde schon heute nicht von `make proto-generate`
    mitgeneriert, ohne dass jemand die Stufe von Hand erweitert — dieser
    Slice zieht `generated-sync` lediglich auf dieselbe, bereits bestehende
    Einschränkung nach, führt sie nicht neu ein.
  Kein Zweit-Check umgesetzt, weil ein Modulpfad-Zweit-Check redundant zum
  ohnehin greifenden Bau-Fehlschlag wäre und ein `.proto`-Datei-Zweit-Check
  eine Eigenschaft wiederherstellen würde, die `make proto-generate` nie
  hatte — das wäre eine Asymmetrie zwischen den beiden Zielen, die dieselbe
  Stufe nutzen, ohne einen dokumentierten Bedarf dafür.
- **Risiko 3 — geteilte Stufe, geteiltes Risiko.** Wird die
  `proto-export`-Stufe generischer gemacht (dynamische `.proto`-Erkennung,
  `go.mod`-Ableitung im `RUN`-Schritt statt fester Argumente), nutzen
  `make proto-generate` **und** `make generated-sync` künftig exakt dieselbe
  Stufe — ein Fehler in der Verallgemeinerung träfe beide Ziele gleichzeitig.
  **Ausgang: entfallen (Implementer-Zug 2026-09-21).** Die Verallgemeinerung
  wurde nicht vorgenommen (siehe §3 Plan-Nachzug) — `Dockerfile` bleibt in
  diesem Lauf unverändert (`git diff --stat Dockerfile` leer), `proto-export`
  wird von `generated-sync` nur **genutzt**, nicht **verändert**. Damit
  entsteht keine geteilte Verallgemeinerungs-Fläche und kein neues
  gemeinsames Fehlerrisiko — das Szenario dieses Risikos tritt nicht ein.
  Trotzdem real erbracht (zusätzliche Absicherung, kein Pflicht-Beleg mehr,
  da die Bedingung „wird generischer gemacht" nicht zutrifft): `make
  proto-generate` nach dem Skript-Umbau erneut laufen lassen —
  `git status --porcelain` danach leer, kein Diff auf den generierten
  Dateien.
- **Risiko 4 — Gegenprobe braucht eine reale Colima-Umgebung mit
  `mounts: []`.** Der DoD-Punkt „Gegenprobe" (§2) ist nur auf einem Rechner
  mit genau dieser Docker-Backend-Konfiguration direkt beobachtbar. **Ausgang:
  weiter offen (Implementer-Zug 2026-09-21).** Real geprüft auf der
  Implementer-Maschine: `colima status` bestätigt `mounts: []` in
  `~/.colima/default/colima.yaml`, **aber** `TMPDIR` ist ein
  `.colima-tmp`-Verzeichnis direkt unterhalb von `$HOME` — also **innerhalb**
  `$HOME` — und fällt damit unter Colimas Default-Verhalten „`$HOME` wird
  beschreibbar gemountet" (Kommentar in derselben `colima.yaml`). Die exakte
  Fehlerkombination (`mounts: []` **und** `TMPDIR` außerhalb `$HOME`) liegt
  auf dieser Maschine nicht vor — der reale Auslöser reproduziert hier nicht,
  ohne dass ein `TMPDIR`-Override probiert wurde (out of scope, §1). Deshalb
  nur der schwächere, netzlose Ersatzbeleg geliefert: `grep -n "docker run"
  tools/harness/generated-sync.sh` zeigt genau eine Aufruf-Zeile
  (`docker run --rm --network none "$GENERATED_SYNC_IMAGE" | tar -x -C
  "$out_dir"`), keine `-v`-Option mehr. Zusätzlich real gelaufen (nicht die
  spezifische Fehlerklasse belegend, aber Funktionsfähigkeit auf dieser
  Maschine bestätigend): `make generated-sync` und `make gates`, beide grün,
  ohne jeden `TMPDIR`-Override. Bleibt **weiter offen** bis ein Rechner mit
  der exakten Kombination real getestet hat. **Nachtrag (Review-Fixrunde
  F-1, 2026-09-21, Quelle:
  `docs/reviews/review-slice-generated-sync-tar-export.md` F-3):** Der
  Reviewer trägt auf seiner eigenen Maschine zufällig genau diese Kombination
  (`mounts: []` **und**, mit `TMPDIR=/tmp` gesetzt, außerhalb `$HOME`) und hat
  sie real hergestellt: der **alte** Skript-Stand (`git show
  1cce85b5:tools/harness/generated-sync.sh`) scheitert dort unter
  `TMPDIR=/tmp` real mit „Permission denied" beim Schreiben nach `/out`, der
  **neue** Stand läuft unter identischem `TMPDIR=/tmp` grün
  (`generated-sync: OK`, `git status --porcelain` leer). Das ist ein direkter
  Fremdmaschinen-Beleg der exakten Fehlerklasse — stärker als der oben
  dokumentierte Ersatzbeleg, aber nicht auf der Implementer-Maschine selbst
  erbracht. Ob das den Ausgang dieses Risikos auf „behoben" hebt oder
  „weiter offen" (im Sinne „kein Repro auf einer im Slice selbst
  kontrollierten Maschine") bestehen bleibt, ist Verifier-/
  Planner-Closure-Entscheidung.
  **Nachtrag (Verifikation, `docs/reviews/
  verifikation-slice-generated-sync-tar-export.md` §2, 2026-09-21):** eine
  dritte, unabhängige Reproduktion (Verifier-eigene Maschine, identische
  Colima-Konfiguration wie die Implementer-Maschine, `TMPDIR=/tmp`-Override)
  bestätigt dieselbe Konvergenz — alter Stand scheitert real, neuer Stand
  läuft grün — und benennt zugleich die Feinheit, dass keiner der drei
  Versuche den DoD-Wortlaut „ohne jeden `TMPDIR`-Override" strikt wörtlich
  erfüllt. **Ausgang (Planner-Closure, 2026-09-21): eingetreten, real
  behoben.** Begründung und die Auslegung des Wortlauts (Test-Fixture- vs.
  Produktions-Workaround-Override) stehen vollständig in §7 und in der
  Begründung an DoD-Punkt 3 (§2) — der Wortlaut-Zweck von §1s
  „Ausdrücklich NICHT"-Ausschluss (kein Host-Konfigurationswechsel **als
  Lösungs-Rezept**) ist gewahrt; die zwei konvergenten Fremdmaschinen-Belege
  plus der Implementer-Beleg unter Default-`TMPDIR` zusammen zeigen den
  neuen Mechanismus unter beiden Bedingungen grün, den alten unter der
  adversen Bedingung real scheiternd.

## 7. Closure-Notiz

**Rollen-Sequenz:** Implementer (`93f5545b`) → Reviewer 1
(`770d23a0`, 1 HIGH F-1) → Fix (`32a58868`) → Reviewer 2/Fixrunde
(`9293cc83`, 0 HIGH/MEDIUM/LOW, 1 INFO) → Verifier (`53828c6d`,
DoD wörtlich unvollständig auf genau einem Punkt) → Planner-Closure
(dieser Zug, 2026-09-21). Fünf Rollen, ein Rückkante-Zyklus, wellenlos.

- **Was hat funktioniert:** Die tar-Stream-Wiederverwendung der bereits
  bestehenden `proto-export`-Stufe (statt eines eigenen, generischeren
  Mechanismus) hat den Umbau klein und risikoarm gehalten — `Dockerfile`
  blieb über den gesamten Slice unverändert (`git diff --stat` leer, dreifach
  unabhängig bestätigt: Implementer, Reviewer, Verifier), kein geteiltes
  Verallgemeinerungsrisiko (§6 Risiko 3) ist eingetreten. Die
  Planner-Vorab-Klärung von Risiko 1 (`ADR-0084`-Berührung) vor dem Start
  hat einen sonst nötigen Architect-Zug vermieden, ohne die ADR-Immutabilität
  zu verletzen — die Lektüre der unberührbaren Liste aus `AGENTS.md` §3.5
  war eindeutig genug für eine begründete Planner-Entscheidung. Die
  Fixrunde war klein und gezielt (ein HIGH, ein Absatz, keine Rückwirkung
  auf die Kernmechanik) — der Reviewer der Fixrunde bestätigte in
  eigenständiger Prüfung, nicht durch Übernahme, dass F-1 real aufgelöst
  war.
- **Was ging anders als geplant:** Der DoD-Punkt 3 (Gegenprobe) konnte auf
  keiner der drei beteiligten Maschinen im striktesten Wortlaut-Sinn
  erfüllt werden — nicht, weil die Fix-Wirksamkeit zweifelhaft wäre
  (sie ist durch zwei unabhängige, konvergente Fremdmaschinen-Reproduktionen
  stark belegt), sondern weil der DoD-Text „ohne jeden `TMPDIR`-Override"
  nicht zwischen zwei strukturell verschiedenen Zwecken eines Overrides
  unterscheidet:
  1. **Override als Produktions-Workaround** — dem Betreiber sagen, er solle
     `TMPDIR` setzen, statt das Skript zu reparieren. Das ist der Fall, den
     §1 „Ausdrücklich NICHT in diesem Slice" explizit ausschließt („Kein
     Colima-/Host-Konfigurationswechsel … das bleibt eine lokale
     Betreiber-Entscheidung außerhalb des Repos").
  2. **Override als Test-Fixture** — `TMPDIR` wird ausschließlich für die
     Dauer eines Vergleichslaufs gesetzt, um die adverse Bedingung auf einer
     Maschine zu **simulieren**, die sie nicht natürlich trägt; die
     *Erfolgsbedingung* des neuen Skripts hängt dabei an keinem bestimmten
     `TMPDIR`-Wert mehr.

  **Planner-Entscheidung (begründet, siehe DoD-Punkt 3 in §2 für den vollen
  Wortlaut):** Der DoD-Text zielt lesbar auf Fall 1 — das ist der Grund, aus
  dem „ohne Override" überhaupt formuliert wurde, nämlich um auszuschließen,
  dass die Lösung selbst von einer Host-Konfigurationsänderung abhängt. Fall
  2 liegt in allen drei Reproduktionen vor (Reviewer, Verifier — beide
  setzten `TMPDIR=/tmp` ausschließlich, um die Bedingung testbar zu machen,
  nicht als empfohlene Lösung), Fall 1 in keiner. Materiell ist die
  Fehlerklasse zusätzlich vollständiger belegt, als der Wortlaut allein
  verlangt: der neue Skript-Stand läuft grün **unter beiden** beobachteten
  `TMPDIR`-Werten (Default innerhalb `$HOME` auf der Implementer-Maschine,
  `TMPDIR=/tmp` außerhalb `$HOME` auf zwei weiteren, unabhängigen Maschinen),
  während der alte Stand unter der zweiten Bedingung auf zwei unabhängigen
  Maschinen real und reproduzierbar mit identischer Fehlermeldung
  scheitert. Das ist der stärkere Beleg (TMPDIR-Wert-Unabhängigkeit direkt
  gezeigt), nicht nur ein Ersatz für den im Wortlaut verlangten Beleg. Ich
  werte DoD-Punkt 3 deshalb als **bestanden** und setze die Checkbox auf
  `[x]` (§2) sowie §6 Risiko 4 auf „eingetreten, real behoben". Ich
  präzisiere den DoD-Wortlaut selbst **nicht** rückwirkend (das wäre eine
  stille Umschreibung eines abgeschlossenen Plan-Abschnitts, kein
  Zitat-Korrektur-Fall nach `ADR-0073`s Logik für ADRs, aber sinngemäß
  dieselbe Vorsicht) — die Auslegung steht hier, sichtbar als
  Planner-Entscheidung mit Datum, nicht als nachträglich „schon immer so
  gemeint" getarnt.

  Eine bewusst verworfene Alternative: den Punkt „wortlaut-offen, materiell
  widerlegt" ins Beobachtungs-Register zu überführen (Verifier-Empfehlung
  (b)) und die Closure entweder zurückzuhalten oder mit einem offenen
  §6-Risiko zu schließen. Dagegen entschieden, weil (a) die Zwei-Zwecke-
  Unterscheidung beim genauen Lesen des bereits vorhandenen §1-Textes
  eindeutig auflösbar ist — keine neue Tatsachenfrage, nur eine
  Textauslegung, die der Planner an dieser Stelle legitim trifft — und (b)
  ein viertes Repro-Verlangen (Maschine mit *natürlichem* TMPDIR außerhalb
  `$HOME`) einen Ertrag von Null gegenüber der bereits zweifach konvergenten
  Evidenz hätte, wie der Verifier selbst festhält.
- **Steering-Loop-Eintrag:** Die HIGH-Klasse „Kommentar trägt keine der
  Kommentar-Klassen" (F-1 des ersten Review-Reports) ist **kein neuer**
  Fund — sie bestätigt zum wiederholten Mal (achter Beleg der bereits
  verkörperten Beobachtung `BEO-PGC/slice-chronik-in-code-kommentar`), dass
  die bestehende Verteidigungslinie (unabhängiger Reviewer + Schritt-20-
  Enumerationslauf) trägt: kein Hard-Rule-Verstoß hat `main` erreicht, auch
  nicht in einem neuen Datei-Typ (Shell-Skript statt Go/SQL). Kein neuer
  Sensor, keine neue Instruktion nötig — siehe Beobachtungs-Register unten.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-PGC/slice-chronik-in-code-kommentar` — neuer Beleg 8
    (`evidence/slice-generated-sync-tar-export.md`), zwei Teile: (1) F-1
    des ersten Review-Reports ist die Klasse in ihrem ersten Vorkommen in
    einem Shell-Skript — Datei-Typ-Erweiterung, keine neue Klasse, Ausgang
    bleibt **verkörpert**. (2) F-2 des ersten Review-Reports
    (`harness/sensors/generated-sync.md`) ist der **zweite** Beleg der in
    Beleg 7 erstmals benannten Doku-Prosa-Grenze (dort `harness/README.md`)
    — diesmal vom Reviewer selbst gefunden statt erst vom Verifier. Zwei
    Belege erreichen noch nicht die für einen eigenen Registereintrag
    verwendete Drei-Beleg-Schwelle (Analogie `commit-traceability`s
    „drittes Auftreten", `harness/README.md` §Traceability rules); der
    Zähler bleibt in derselben Beobachtung stehen, als verstärktes Signal
    für den nächsten Lese-Schritt.
  - Die Gegenprobe-Wortlaut-Frage (DoD-Punkt 3, oben ausführlich begründet)
    erhält **keinen** neuen Registereintrag. Erstauftreten dieser
    spezifischen Zwei-Zwecke-Unterscheidung (Test-Fixture- vs.
    Produktions-Workaround-Override in einer DoD-Gegenprobe-Formulierung);
    per Präzedenzfall `BEO-PGC/plan-nachzug`-Logik und dem in diesem
    Slice-Plan selbst (§8) benannten Prinzip „beim ersten Auftreten reicht
    die Nennung in diesem Slice-Plan" reicht diese ausführliche
    Dokumentation hier. Tritt eine strukturell gleiche Wortlaut-Lücke in
    einem künftigen Gegenprobe-/Repro-DoD-Kriterium ein zweites Mal auf, ist
    das der Anlass für einen eigenen Registereintrag.
  - Die den Slice auslösende Beobachtung selbst (Colima-`mounts: []` +
    UID-Mapping außerhalb `$HOME`, §1 „Autor"-Zeile) bleibt ebenfalls ohne
    eigenen Registereintrag — bereits bei Planung (§8) als Erstauftreten
    geprüft und benannt, jetzt durch drei konvergente Maschinen-Reproduktionen
    strukturell aufgelöst; kein wiederkehrendes Muster über diesen einen
    Auslöser hinaus beobachtet.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet.
- **Risiken aus §6:**
  1. `ADR-0084`-Berührung — **entfallen** (Planner-Klärung vor Start, §3
     Plan-Nachzug; vom Verifier eigenständig gegen §Entscheidung/
     §Fitness-Function bestätigt).
  2. Verlust der unabhängigen Modulpfad-/`.proto`-Cross-Checks —
     **entfallen** (Implementer-Zug, bewusst kein Ersatz-Check; vom
     Verifier eigenständig gegen den Einführungs-Commit `338cfe00`
     nachgeprüft und bestätigt).
  3. Geteilte Stufe, geteiltes Risiko — **entfallen** (Verallgemeinerung
     nicht vorgenommen, `Dockerfile` real unverändert über den gesamten
     Slice-Umfang, dreifach bestätigt).
  4. Gegenprobe braucht reale Colima-`mounts: []`-Umgebung — **eingetreten,
     real behoben**: zwei unabhängige Fremdmaschinen-Reproduktionen
     (Reviewer, Verifier) plus die Planner-Auslegung des DoD-Wortlauts
     (siehe „Was ging anders als geplant" oben).
- **Drei Paarungen** (Repo ohne Wellen-Betrieb für diesen wellenlosen
  Slice, hier direkt geprüft statt über eine Welle):
  1. **Anker:** `ADR-0084` (Festlegung 1, unverändert) und `ADR-0060` (der
     Generator) — beide bereits im Slice-Kopf (§Bezug) genannt, keine neue
     ADR nötig (siehe Risiko 1).
  2. **Folge-Slice:** keiner — siehe oben.
  3. **Register:** `BEO-PGC/slice-chronik-in-code-kommentar` fortgeschrieben
     (Beleg 8); keine neue Beobachtung angelegt (zwei Erstauftreten-Fragen
     jeweils per Slice-Plan-Text statt eigenem Registereintrag beantwortet,
     siehe oben).

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
