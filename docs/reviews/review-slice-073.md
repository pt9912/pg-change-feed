# Review-Report: slice-073 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-073`, Commit `517eee2`. Der Slice-Diff ist genau dieser
Commit (3 Pfade: `.githooks/commit-msg` neu · `harness/README.md` eine Zeile ·
Slice-Plan §2/§3); sein unmittelbarer Elter ist `1cf4a61` (reiner
`next→in-progress`-Move). Der danach liegende Zug `4240b36`
(Register-Verweisform in `BEO-PGC/commit-traceability-kein-vorab-hook/state.md`)
ist Rollen-Arbeit eines anderen Kontexts und **nicht** Teil des Gegenstands.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14 — die Regel `Neue Betreiber-Oberfläche ohne Handbuch-Zug` liegt
**vor** diesem Implementer-Lauf).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md`
  vollständig (§1–§8), einschließlich des §3-Nachzugs aus `517eee2`
- [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  vollständig — §Entscheidung Punkt 1–4, §Verglichene Alternativen,
  §Konsequenzen, §Fitness Function, §Re-Evaluierungs-Trigger (a)/(b)
- [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
  vollständig (nur die Hook-Klausel ihrer §Entscheidung ist superseded; die
  positive Hälfte „je Message", die Grenz-Hälfte „im Betreff", Fenster und
  Werkzeug-Teilung bleiben bindend)
- `tools/harness/commit-traceability.sh`, `.d-check.yml` (`ids`-Muster,
  `commits`-Abschnitt), `harness/mk/doc-gate.mk`, `d-check.mk`
- `harness/README.md` §Traceability rules (die neue Zeile), `AGENTS.md` §3
  (Hard Rules) und §5, `harness/conventions.md` (MR-000)
- `docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md`
  (die vier Design-Vorgaben und ihre Herkunft)
- `docs/user/benutzerhandbuch.md` (§1 Zielgruppe, §4 Aufgaben, §5 ENV)
- Fremdquelle zum Modul: `d-check`s `internal/hexagon/core/rules/commits.go`
  (`CheckCommitMessage`, `cleanCommitMessage`, `commitSubject`) — gelesen, um
  das Verhalten des **gepinnten** Images zu verstehen, nicht als
  Prüfgegenstand

---

## Findings

### F-1 — Die positive Hälfte prüft die ganze Message-Datei, `ADR-0062` Punkt 3 fixiert für sie den Betreff; Träger der Abweichung ist eine Plan-Zelle

- `kategorie`: HIGH
- `quelle`: [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  §Entscheidung Punkt 3 („positiv: ≥ 1 `LH-*`/`ADR-*`-Kennung **im Betreff**")
  · [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
  §Entscheidung („die positive Hälfte (≥ 1 `LH-*`/`ADR-*` **je Message**)")
  · `AGENTS.md` §3.5 und Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz
  (Verdikt 3: „Lockerung legitim, aber undokumentiert")
- `pfad`: `.githooks/commit-msg:57-64` (Prüfung) ·
  `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md:130`
  (Plan-Nachzug) · `:40` (§1-Ziel) · `:88` (DoD-Zeile 2)
- `befund`: Der Hook prüft die positive Hälfte per `grep -qE` über die **rohe**
  Message-Datei, während `ADR-0062` Punkt 3 für sie den **Betreff** fixiert.
  Die Messung des Implementers trägt (das d-check-Modul `commits` akzeptiert
  eine Kennung, die nur im Body steht — von mir am gepinnten Image
  nachgestellt), und die Abweichung liegt in der Richtung, die Punkt 3 selbst
  als Zweck nennt („kein Auseinanderlaufen vom Standing-Gate"); ihr Träger ist
  aber ausschließlich eine **Zelle** der §3-Tabelle des Slice-Plans, eines
  Zeitdokuments, das mit der Closure nach `done/` wandert und archiviert wird.
  §1 („wenn sein Betreff keine `LH-*`-/`ADR-*`-Kennung trägt") und die
  DoD-Zeile 2 („ganz ohne `LH-*`-/`ADR-*`-Kennung im Betreff") tragen
  unverändert den Betreff-Scope, sind also nach dem eingetretenen Nachzug
  inhaltlich unrichtig über das Artefakt. Der Text der Accepted-ADR ist dabei
  in sich widersprüchlich: ihr Punkt 3 verlangt „im Betreff" **und** die
  Spiegelung der d-check-Positiv-Hälfte, die ihrerseits je Message liest —
  beide Anforderungen sind nicht gleichzeitig erfüllbar.
- `verifizierbar`: nein — kein Gate prüft die Übereinstimmung von ADR-Text und
  Hook-Semantik; die Messung ist per `docker run … --commit-msg <datei>` am
  gepinnten Image reproduzierbar, die Carrier-Frage ist eine Rollen-Frage.
- `klasse`: „Abweichung von einer Accepted-ADR ohne Folge-ADR — Träger ist ein
  Zeitdokument statt der Entscheidung"

### F-2 — „Spiegelt exakt" ist gemessen unwahr: der Hook lässt in drei belegten Klassen durch, was das Standing-Gate danach verwirft

- `kategorie`: HIGH
- `quelle`: [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  §Entscheidung Punkt 3 („spiegelt **exakt** die zwei bestehenden Regeln") ·
  §Re-Evaluierungs-Trigger (b) · Slice-Plan §6 Risiko 1 · `harness/README.md:165`
  („meldet denselben Verstoß") · Modul `commits` des gepinnten Images,
  dessen Quelltext unter
  `d-check`s `internal/hexagon/core/rules/commits.go`
  (`cleanCommitMessage`, `commitSubject`) liegt und dessen Modul-ADR als
  Fitness Function festhält: „`#`-Kommentar-only-Kennung zählt **nicht**
  (uniforme Bereinigung)" — die ID dieses Modul-ADRs gehört dem
  d-check-Repo und wird hier deshalb **nicht** zitiert (die Nummer ist in
  diesem Repo belegt)
- `pfad`: `.githooks/commit-msg:28-32` (Betreff-Extraktion) · `:39-41`
  (Muster) · `:57-64` (positive Hälfte) · `harness/README.md:165` ·
  `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md:159-165`
  (§6 Risiko 1)
- `befund`: Die Begründung des §3-Nachzugs setzt „das d-check-Modul liest die
  gesamte Message" mit „der Hook darf die ganze Datei greppen" gleich; das
  Modul liest die **bereinigte** Message (alles ab der ersten
  scissors-Zeile `^#.*>8` und jede `#`-Zeile entfällt), der Hook die **rohe**
  Datei — die Differenz ist genau der Kommentar-/scissors-Bereich. Am
  gepinnten Image (`--commit-msg` auf genau die rohe Datei, die der Hook
  sieht) und an je einem **realen Commit mit aktivem Hook**
  (`core.hooksPath .githooks`, Commit angelegt mit Exit 0, das Gate danach
  rot) sind drei Klassen belegt: (a) Kennung nur in einer `#`-Kommentarzeile
  (`git commit -m '<Betreff>' -m '# ADR-…'`, die Zeile bleibt unter git's
  Default-Cleanup `whitespace` in der gespeicherten Message, das Modul
  bereinigt sie weg); (b) Kennung nur hinter der scissors-Zeile
  (`git commit -v` bzw. `-F <datei> --cleanup=scissors` — das Modul bricht am
  Cleanup an der scissors ab); (c) `SPEC-*` auf der Fortsetzungszeile des
  ersten Absatzes (git's `%s` fügt den Absatz zusammen, der Shell-Sensor
  färbt rot, der Hook sieht nur die erste Zeile). In allen drei Klassen ist
  der Hook grün und `make commit-traceability` rot. Die Richtung ist in allen
  belegten Klassen **laxer**; ein Fall, in dem der Hook einen Commit
  zurückweist, den das Gate zulässt, wurde weder real noch im
  Modul-Message-Modus gefunden — die Zusage der §3-Zelle trägt also, die
  ADR-Formulierung „exakt" und die README-Zeile „denselben Verstoß" nicht.
- `verifizierbar`: ja — `make commit-traceability` bzw. `make gates` über die
  betroffene Range; alle drei Klassen sind als realer Commit mit aktivem Hook
  reproduzierbar (keine Umgehung des Hooks nötig) und zusätzlich per
  `docker run … --commit-msg <rohe Message-Datei>` am gepinnten Image
  (Digest aus `d-check.mk`) isoliert.
- `klasse`: „Spiegelung ist Approximation — belegte Divergenz-Klasse zwischen
  lokalem Hook und dem Modul, das die Regel trägt"

### F-3 — Der reale Nachweis der Grenz-Hälfte isoliert sie nicht: derselbe Input wird auch von der positiven Hälfte zurückgewiesen

- `kategorie`: LOW
- `quelle`: Slice-Plan §2 DoD-Zeile 1 („real ausgeführt, nicht nur am
  Skripttext behauptet") · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-13-quality-gates.md` (Mutations-Beleg als Nachweisform)
- `pfad`: `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md:84-86`
- `befund`: Der für die Grenz-Hälfte geführte Nachweis benutzt einen Betreff
  mit `SPEC-*` **ohne** Vertrags-Kennung (`.githooks/commit-msg:45-50` und
  `:58-63` weisen ihn beide zurück). Setzt man die Grenz-Hälfte im Skript auf
  `if false` (eigene Mutation dieses Laufs), bleibt der Nachweis-Input rot —
  er unterscheidet die beiden Hälften damit nicht. Der Nachweis selbst bleibt
  gültig, weil der beobachtete Fehlertext die Grenz-Hälfte benennt; es fehlt
  die Isolierung gegen einen künftigen Ausfall genau dieser Hälfte.
- `verifizierbar`: ja — Mutations-Lauf am Skript (Grenz-Hälfte aus → der
  isolierende Input mit Kennung **und** Struktur-ID kippt von Exit 1 auf 0).
- `klasse`: „Nachweis ohne Isolierung der geprüften Hälfte"

### F-4 — Der Kommentar über dem Gate-Target verweist weiter auf die Alternativen-Tabelle, deren Hook-Aussage `ADR-0062` korrigiert hat

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7 (Rang-Zeiger) ·
  [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  §Status (nur die Klausel „Kein commit-msg-Hook" ist superseded)
- `pfad`: `d-check.mk:54-55` (nicht im Diff dieses Slice)
- `befund`: Die Zeile „Gate statt commit-msg-Hook (`ADR-0045`, Alternativen)"
  zitiert als Begründung genau die Tabelle, deren Aussage `ADR-0062` in der
  Hook-Klausel korrigiert hat; wer dort nachliest, findet die abgelöste Lage
  ohne Zeiger auf die Korrektur. Die getragene Aussage bleibt richtig (das
  Gate ist weiterhin die durchsetzende Instanz, `ADR-0062` Punkt 1) — es fehlt
  der Rang-Zeiger, nicht die Regel.
- `verifizierbar`: nein — `make docs-check` prüft Referenzen, nicht die
  Aktualität eines Begründungs-Zeigers.
- `klasse`: „Kommentar-Verweis auf eine teilweise abgelöste ADR ohne Zeiger"

---

## Urteile zu den gestellten Prüfpunkten

**1. Trägt die Aufteilung der Prüf-Weiten? — In der Richtung ja, in der Form
nein.** Drei getrennte Fragen, drei Antworten:

- *Ist die Messung des Implementers richtig?* **Ja.** Am gepinnten Image
  (`ghcr.io/pt9912/d-check@sha256:18e9cd85…`, Digest aus `d-check.mk`) weist
  das Modul `commits` einer Message mit Kennung nur im Body **keinen** Befund
  zu (`--commit-msg` auf eine Datei „Betreff ohne Kennung / Body mit
  `ADR-0045`": Exit 0), und `tools/harness/commit-traceability.sh` ist rein
  Betreff-scoped (`git log --format='%h %s'`, nur die Grenz-Hälfte). Ein
  Betreff-scoped positiver Hook-Check wäre damit tatsächlich **strenger** als
  das Gate — die Begründung des §3-Nachzugs trifft den Punkt.
- *Ist die Abweichung damit gedeckt?* **Nein.** `ADR-0062` Punkt 3 fixiert für
  die positive Hälfte den Betreff; ein Accepted-ADR-Text wird durch eine
  Folge-ADR geändert (`AGENTS.md` §3.5), nicht durch eine Plan-Zelle. Der Text
  ist dabei in sich widersprüchlich (siehe F-1) — die Auflösung dieses
  Widerspruchs ist eine Architect-Entscheidung mit Übergabe-Artefakt
  (Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad, Verdikt 2 oder 3), und
  `§1`/DoD-Zeile 2 des Plans brauchen danach denselben Nachzug, den §3
  bekommen hat.
- *Trägt die Aufteilung technisch?* **Nein, nicht als „exakt"** — die
  Betreff/Ganze-Datei-Teilung löst genau ein Divergenz-Muster ab (die Kennung
  im Body) und lässt drei andere stehen (F-2). Der Kern: das Modul liest die **bereinigte**
  Message, der Hook die **rohe** Datei; die Teilung „positive Hälfte = ganze
  Datei, Grenz-Hälfte = Betreff" beschreibt den Unterschied nicht.

**2. Die vier realen Nachweise — alle vier selbst reproduziert.** In einem
Wegwerf-Repository (`/tmp/review073/repo`, Kopie des Skripts unter
`core.hooksPath .githooks`, nicht der Arbeitsbaum):

| Nachweis | Eingabe | Hook-Exit | Beobachtung |
|---|---|---|---|
| `SPEC-`/`ARC-`-Rückweisung | `docs: update SPEC-012 table` | 1 | Meldung nennt die Grenz-Hälfte namentlich |
| Rückweisung ohne Kennung | `docs: gar keine Kennung` | 1 | Meldung nennt die positive Hälfte |
| konformer Commit | `docs: konform (ADR-0045)` | 0 | Commit abgeschlossen |
| Merge-Commit | `Merge branch 'feature'` (`git merge --no-ff`) | 0 | Commit abgeschlossen |

**3. Die drei rot gesehenen Mutationen — alle drei, je an einer Kopie des
Hooks, ohne Eingriff in den Arbeitsbaum:**

| Mutation | Ort | Ergebnis |
|---|---|---|
| Grenz-Hälfte aus (`if false`) | `.githooks/commit-msg:45` | isolierender Input (Kennung **und** `SPEC-012`) kippt 1 → 0 |
| positive Hälfte aus (`if false`) | `:58` | „gar keine Kennung" kippt 1 → 0 |
| Ausnahme aus (`if true`) | `:57` | `Merge branch 'feature'` kippt 0 → 1 |

Die Paarung Mutation ↔ Nachweis ist in beiden Richtungen dieselbe, die die
DoD nennt; die Konfundierung aus F-3 trifft nur den *ersten* Nachweis in seiner
Originalform, nicht die Mutation.

**4. Nicht-durchsetzend — bestätigt, an beiden Artefakten.** `grep` über
`Makefile`, `harness/mk/*.mk`, `d-check.mk` und `.github/workflows/` findet
keinen Hook-Aufruf; `GATE_CHECKS` bleibt `docs-check`, `commit-traceability`,
`a-check`, `coverage-gate`, `baseline-verify` — das Standing-Gate ist
unverändert im Bündel, der Hook in keinem. Kein `docker` im Hook
(`grep -n docker .githooks/commit-msg` leer; die einzige Nennung steht im
Kommentar als Abgrenzung). `git commit --no-verify -m "docs: gar keine
Kennung ueberhaupt"` läuft mit Exit 0 durch. Ein frisches Repository mit
vorhandener Hook-Datei, aber ohne `core.hooksPath`, committet eine ID-lose
Message ebenfalls mit Exit 0 (Opt-in bestätigt). Der Hook ersetzt das
Standing-Gate nicht — `ADR-0062` Punkt 1 und die README-Zeile tragen dieselbe
Aussage.

**5. Handbuch — kein neuer Zug; ich teile die Einschätzung des Implementers.**
Das Handbuch richtet sich laut §1 an „Betreiber und Integratoren", seine
Oberflächen sind §5 (ENV des Feed-Containers), §4 (Aufgaben) und §2 (Zugriff
und Rollen); ein lokaler Git-Hook wirkt auf den Entwickler im Repository,
nicht auf den Betrieb — er erscheint in keiner der vier Aufzählungen der
HIGH-Regel (`CDC_*`-Variable, `cdc.*`-Funktion, Horch-Adresse, Endpunkt) und
ist kein öffentlicher Vertrag dieses Produkts. Die Handbuch-Zeile bleibt
unberührt, also greift auch die Versionshistorie-Regel nicht; die neue
README-Zeile ist der richtige Träger (Harness-Einstieg). Kein Finding.

**6. Kommentar-Disziplin (`AGENTS.md` §3.7) — ohne Befund.** `grep -nE
"slice-|welle-|frueher|früher|vorher|ohne .* waere"` über `.githooks/commit-msg`
ist leer; kein Kommentar beschreibt eine verworfene Alternative im Konjunktiv
oder einen abwesenden Text. Die sechs Kommentarblöcke tragen je eine der
Klassen: Zusage („Exit != 0 bricht den Commit ab"), Kopplung („spiegelt die
zwei Haelften … `id-patterns` und `exempt-pattern` aus `.d-check.yml`"),
Abgrenzung („kein Gate-Ziel und kein Ersatz fuer `make commit-traceability`
(`ADR-0062`)"), Grenze („ohne lesbare Message-Datei … bricht hier nicht ab")
und Rang-Zeiger (`AGENTS.md` §5, `harness/README.md` §Traceability rules).
Der Satz „das d-check-Modul `commits` liest die gesamte Message" ist als
Zusage richtig, als Begründung der *Datei*-Weite irreführend — das steht in
F-2, nicht hier.

**7. §2 DoD-Häkchen — korrekt getrennt.** Gesetzt sind die sechs Punkte, die
der Implementer liefern konnte (1–4 die realen Nachweise, 5 die
Aktivierungs-Doku als README-Zeile, 6 `make gates`, 8 das Doku-Update, 9 der
Reconciliation-Entfall mit Begründung); offen bleiben Review (7),
Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge und die drei Paarungen
— genau die Arbeit nach diesem Report. Kein Häkchen greift einem anderen
Punkt vor (insbesondere ist der Merge-Nachweis als eigener Punkt geführt, wie
das gegengeprüfte Architect-Verdikt ihn gefordert hatte). Die Nachweise halten
der Prüfung stand (Urteil 2); die DoD-Zeile 2 formuliert allerdings weiterhin
„im Betreff" und beschreibt damit nach dem §3-Nachzug nicht mehr, was der Hook
tut (F-1).

**8. Zusatz-Hypothesen zum Modul `commits` (Quelltext-Lesung) — zwei
bestätigt, eine widerlegt; Folgen für §6.** Die vermuteten Divergenzklassen habe ich nicht übernommen,
sondern am **gepinnten Image** (Message-Modus `--commit-msg` auf genau die
rohe Datei, die der Hook liest; das Image kennt den Modus, `--help` weist ihn
als „Modul commits: eine Commit-Message aus <datei>" aus) **und** an realen
Commits in einem Wegwerf-Repository (`/tmp/review073/repo`, **nicht** dieser
Baum) mit **aktivem** Hook gemessen — die Kurz-`sha`s der Tabelle gehören
diesem Wegwerf-Repository:

| Fall | Hook | Modul `--commit-msg` | Gate (Range) am realen Commit | Ergebnis |
|---|---|---|---|---|
| Kennung nur in einer `#`-Kommentarzeile | 0 | 1 | 1 — `commit-untraceable`, Commit-Exit 0 | Divergenz, real erreichbar |
| Kennung nur hinter der scissors-Zeile/im Verbose-Diff | 0 | 1 | 1 — `commit-untraceable`, Commit-Exit 0 | Divergenz, real erreichbar |
| Leerzeile vor `Merge branch 'x'` (Default-Cleanup) | 0 | 1 | **0** — git strippt die Leerzeile, `%s` = `Merge branch 'x'` | **widerlegt** (s. u.) |
| Leerzeile vor `Merge branch 'x'` (`--cleanup=verbatim`) | 0 | 1 | 1 | Randfall **zugunsten** des Hooks |
| Kennung im Body, Betreff ohne Kennung | 0 | 0 | 0 | Implementer-Messung bestätigt |
| gar keine Kennung | 1 | 1 | 1 | Deckung |
| Muster-Stichproben `ADR-045` / `LH-FA-CFG-005.a` / `LH-FA-CFG-5` / `LH-QA-POR-001` | 1 / 0 / 1 / 0 | 1 / 0 / 1 / 0 | — | Muster-Äquivalenz, kein Befund |

Zwei Divergenzklassen (Kommentarzeile · scissors/Verbose-Diff) treffen damit
zu und sind an realen Commits erreichbar; die dritte vermutete (Leerzeile vor
`Merge …`) ist im Default-Pfad **widerlegt**, und wo sie auftritt
(`--cleanup=verbatim`, ein ausdrückliches „nicht bereinigen"), ist der Hook
**git-treuer** als das Modul: git selbst setzt den Betreff auf
`Merge branch 'x'`, ein Delegieren an den `--commit-msg`-Modus hätte diesen
Commit fälschlich verworfen. Diese Klasse steht deshalb **nicht** als Finding,
sondern als benannter Randfall **zugunsten** des Hooks. Die
Subjekt-Variablen-Frage (eine Variable, zwei nötige Semantiken) bleibt als
Mechanismus bestätigt, mit einem zusätzlichen Effekt **zuungunsten** des
Hooks: `SPEC-*` auf der Fortsetzungszeile des ersten Absatzes (Realfall
`521716b` desselben Wegwerf-Repositories, s. Klasse (c) in F-2). Die Regexe des Hooks sind zu den
`id-patterns`/dem `exempt-pattern` des Moduls äquivalent — eigene Stichproben
mit identischen Exit-Codes in beiden Richtungen (Tabellenzeile).

**Zur Folge für §6:** die dort als „weiter offen" geführte
Doppel-Implementierung ist nicht mehr nur das von `ADR-0062` §Konsequenzen
bewusst akzeptierte Restrisiko — sie ist **belegt divergierend**, ohne dass
eine der beiden Seiten geschärft wurde, und zwar auf drei real erreichbaren
Pfaden. `ADR-0062` §Re-Evaluierungs-Trigger (b) („der Hook und das
Standing-Gate weichen real auseinander") ist damit faktisch eingetreten und
braucht beim Trigger-Audit der Closure einen Ausgang (Vereinheitlichung oder
Rückbau als Folge-ADR, Modul 6 §Trigger-Audit). Die *Härte* des Ergebnisses
bleibt benannt: die Abweichungen laufen in die **laxe** Richtung, die Zusage
der §3-Zelle („weist keinen Commit zurück, den das Standing-Gate zulässt")
ist in keinem Fall gebrochen — der Preis ist die Wirkung des Hooks als
Frühwarnung in genau diesen Klassen.

---

## Negativbefunde

- geprüft, ohne Befund: `.githooks/commit-msg` (Modus `100755`, kein Docker,
  keine Chronik im Kommentar, `--no-verify`-Weg unverändert; die vier
  Nachweise und die drei Mutationen aus Urteil 2/3 laufen wie im Plan benannt)
- geprüft, ohne Befund: `harness/README.md:165` (Opt-in, Nicht-Ersatz,
  `--no-verify` — alle drei Zusagen am Artefakt bestätigt; die Formulierung
  „denselben Verstoß" trägt F-2, nicht diese Zeile als eigener Befund)
- geprüft, ohne Befund: `Makefile` / `harness/mk/*.mk` / `d-check.mk` /
  `.github/workflows/` (kein Hook-Aufruf, `GATE_CHECKS` unverändert, keine
  Workflow-Berührung)
- geprüft, ohne Befund: `tools/harness/commit-traceability.sh` (unberührt; die
  Grenz-Hälfte bleibt betreff-scoped, `%s`-Semantik unverändert)
- geprüft, ohne Befund: `docs/plan/adr/0045-*` und `docs/plan/adr/0062-*`
  (nicht im Diff; die Immutabilität ist gewahrt, der Supersedes-Zeiger liegt
  im Index und ist `make docs-check`-grün)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (keine neue
  Betreiber-Oberfläche im Diff, keine inhaltliche Änderung — die
  HIGH-Regel ist nicht ausgelöst, s. Urteil 5)
- geprüft, ohne Befund: Slice-Plan §1–§8 (Out-of-Scope mit vier begründeten
  Ausschlüssen, §6-Risiken mit Adresse auf `ADR-0062`s Trigger, §8-Sichtung
  mit drei beantworteten Einträgen — die §1/§2-Wortweite trägt F-1, die
  §6-Statusfrage trägt Urteil 8)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Abweichung von einer Accepted-ADR ohne
Folge-ADR — Träger ist ein Zeitdokument statt der Entscheidung" · „Spiegelung
ist Approximation — belegte Divergenz-Klasse zwischen lokalem Hook und dem
Modul, das die Regel trägt" · „Nachweis ohne Isolierung der geprüften Hälfte" ·
„Kommentar-Verweis auf eine teilweise abgelöste ADR ohne Zeiger"

## Verdikt

**Merge-blockierend:** ja — zwei HIGH. Beide hängen an **derselben** Stelle,
`ADR-0062` §Entscheidung Punkt 3, aber an verschiedenen Hälften: F-1 an ihrem
Wortlaut (positive Hälfte „im Betreff" vs. ganze Message-Datei), F-2 an ihrer
Behauptung („spiegelt exakt" — drei belegte Divergenzklassen, s. Urteil 8). Was **nicht**
blockiert und ausdrücklich bestätigt bleibt: der Hook ist
nicht-durchsetzend in jeder geprüften Hinsicht (Urteil 4), er ersetzt das
Standing-Gate nicht, er ruft kein Docker, er lässt sich per `--no-verify`
umgehen, die Aktivierung ist Opt-in, und die §3-Zusage „weist keinen Commit
zurück, den das Standing-Gate zulässt" trägt in allen gemessenen Fällen.

**Übergabe:** F-1 und F-2 gehen als **HIGH mit Rollen-Widerspruch** in den
Konflikt-Pfad (Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz): der
Reviewer kann den Widerspruch zwischen ADR-Wortlaut und Begründung nicht
selbst auflösen, und ein Herabsetzen wegen des begründeten Implementer-Widerspruchs
ist ausgeschlossen. Der nächste Zug ist der **Architect** mit den Verdikten 2
oder 3 (Folge-ADR bzw. Nachtrag zu `ADR-0062` Punkt 3, plus Plan-Nachzug in
§1 und DoD-Zeile 2) — erst danach entscheidet sich, ob der Implementer die
Hälften-Semantik nachziehen muss (falls das Verdikt die „exakte" Spiegelung
aufrechterhält, sind die drei Klassen aus F-2/Urteil 8 der Auftrag). F-3 und F-4
sind isoliert und brauchen keinen Rückgabe-Pfeil; sie gehen mit diesem Report
in die Closure.

**DoD-Nachzug:** **kein** Nachzug in diesem Report — zwei HIGH offen, die
Fixrunde ist damit nicht ausgeschlossen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde greift nur bei 0 HIGH bzw. wenn alle
Findings ohne Reviewer→Implementer-Pfeil weitergereicht werden). Die Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des
Slice-Plans bleibt offen. Der Report ersetzt keine Verifikation — DoD-/Spec-
Konformität prüft der Verifier separat (Modul 11).

---

## Ausgeführte Sensoren (Exit-Codes dieses Laufs)

| Aufruf | Exit-Code | Beleg |
|---|---|---|
| `make gates` | 0 | zwei Läufe: vor diesem Report 584 Dateien, mit ihm 585 — beide 0 Befunde · `commit-traceability` OK (5 Commits) · `a-check` 0 Befunde · `coverage-gate` 49,30 % ≥ 35 % · `baseline-verify` OK |
| `make commit-traceability` | 0 | `HEAD~5..HEAD`: d-check 0 Befunde · „Betreffs ohne Struktur-ID" |
| vier Hook-Nachweise (eigene) | 1 / 1 / 0 / 0 | Wegwerf-Repository `/tmp/review073/repo`; s. Urteil 2 |
| drei Mutationen (eigene) | 1→0 / 1→0 / 0→1 | Kopien des Skripts; s. Urteil 3 |
| `--commit-msg`-Messungen am Image | 1 / 1 / 1 / 0 / 1 | Urteil 8; rohe Message-Dateien, Digest `sha256:18e9cd85…` aus `d-check.mk` |
| drei Divergenzklassen (eigene, mit **aktivem** Hook) | Commit 0 / Gate 1 (je) | Wegwerf-Repository: (a) `2c362ff` (b) `d1b4b32` (c) `521716b`; je Hook grün, das Gate danach rot |
| „Leerzeile vor `Merge …`" Default vs. `--cleanup=verbatim` | Commit 0 / 0 | Wegwerf-Repository: `7c85cd6` Gate 0 (H2 widerlegt) · `b0dc4dc` Gate 1 nur unter `verbatim` |
| `git commit --no-verify` bzw. Repo ohne `core.hooksPath` | 0 / 0 | Urteil 4 (Umgehungsweg und Opt-in) |

Kein Lauf dieses Reports hat den Arbeitsbaum verändert (Messungen und
Mutationen ausschließlich in `/tmp/review073/`); `git status --porcelain` ist
vor dem Anlegen dieses Reports leer.
