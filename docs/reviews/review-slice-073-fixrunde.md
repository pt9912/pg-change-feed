# Review-Report: slice-073 — Fixrunde — 2026-09-15

**Review-Art:** Code — Bestätigungslauf zu einer Fixrunde nach eigenem
Vorbefund (`docs/reviews/review-slice-073.md`: F-1 HIGH, F-2 HIGH, F-3 LOW,
F-4 INFO). **Gegen die Rollen-Falle dieses Laufs:** Jede behauptete Änderung
ist am **Artefakt** nachgeprüft (Diff, Datei-Ist-Stand, sha256 über den
Verhaltens-Code) und jede Zusage in einem **eigenen** Wegwerf-Repository neu
gemessen — nichts ist aus der Commit-Message, dem Verdikt oder dem
Implementer-Bericht übernommen (Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-08-agentenrollen.md` §Rollen-Regeln).

**Gegenstand:** die vier Commits seit dem Vorbefund — `73a96ef`
(Slice-Plan: §1/§2/§3), `78d4007` (`.githooks/commit-msg`-Kommentar +
`harness/README.md`), `41e409f` (`d-check.mk`-Zeiger), `167130f` (Kopf-Feld
`Bezug`). **Kontext, nicht Gegenstand** (sie tragen die Entscheidung des
Konflikt-Pfads): `cf5315a`, `35f152e`
([`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)),
`48f6b24`, `d6139f5` (Architect-Verdikt).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
  vollständig (§Status, §Entscheidung Punkt 1/2, §Verglichene Alternativen,
  §Konsequenzen, §Fitness Function, §Re-Evaluierungs-Trigger)
- [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  §Entscheidung Punkt 1–4 vollständig (Punkt 3 als abgelöster Wortlaut mit
  seinen **drei** Aussagen gelesen), §Status;
  [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
  §Entscheidung
- eigener Vorbefund `docs/reviews/review-slice-073.md` (F-1 bis F-4, Urteile 1–8)
- `docs/reviews/architect-verdict-commit-msg-hook-einseitige-zusage.md`
- `.githooks/commit-msg` (Ist-Stand, mit Kommentaren), `harness/README.md`
  §Traceability rules, `d-check.mk`, `docs/plan/adr/README.md`,
  Slice-Plan `docs/plan/planning/in-progress/slice-073-commit-msg-git-hook.md`
- `docs/plan/planning/observations/BEO-PGC/commit-traceability-kein-vorab-hook/state.md`
- `AGENTS.md` §3.5/§3.7/§3.9, `harness/conventions.md` (MR-000)
- eigene Messungen: gepinntes Image `ghcr.io/pt9912/d-check@sha256:18e9cd85…`
  (Digest aus `d-check.mk`), Wegwerf-Repositories `/tmp/review073/repo3`,
  `--print-mk`

---

## Findings

### F-1 (Vorbefund HIGH) — Abweichung von `ADR-0062` Punkt 3 — geschlossen, getragen durch `ADR-0069`

- `kategorie`: HIGH (Vorbefund, jetzt geschlossen)
- `quelle`: [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
  §Status/§Entscheidung · `AGENTS.md` §3.5
- `pfad`: `docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md:3-8` ·
  `docs/plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md`
  (unberührt) · Slice-Plan `:12-27` (`Bezug`), `:44-60` (§1), `:98-106` (§2),
  `:148` (§3)
- `befund`: Der Träger der Abweichung ist von der Plan-Zelle in die
  **Entscheidung** gewandert. `ADR-0069` ist Accepted, nennt sich selbst
  Teil-Supersede von `ADR-0062` Punkt 3 und formuliert die Zusage des Hooks
  neu (einseitig, message-weite positive Hälfte, Grenz-Hälfte betreff-scoped,
  Hook-Logik unverändert). `ADR-0062`s Datei ist per `git diff 517eee2 HEAD --`
  **unberührt** (Immutabilität gewahrt, `AGENTS.md` §3.5); der ADR-Index trägt
  bei der `ADR-0062`-Zeile den Zeiger „(→ `ADR-0069`)". Der Plan ist an allen
  drei vom Verdikt genannten Stellen nachgezogen: §1 nennt die Prüf-Weite
  **message-weit** und die drei offenen Klassen, §2 DoD-Zeile 2 ebenfalls
  message-weit, §3 trägt statt der Begründung einen **Rang-Zeiger** auf
  `ADR-0069`. Keine Plan-Stelle behauptet mehr den Betreff-Scope der positiven
  Hälfte (grep über die vier Betreff-Stellen: `:47`, `:99`, `:110` — alle
  betreffen die **Grenz-Hälfte**, die zu Recht betreff-scoped ist).
- `verifizierbar`: nein — kein Gate prüft die Übereinstimmung von ADR-Text und
  Artefakt; die Carrier-Frage ist eine Rollen-Frage.
- `klasse`: „Abweichung von einer Accepted-ADR ohne Folge-ADR — Träger ist ein
  Zeitdokument statt der Entscheidung" (Vorbefund, jetzt geschlossen)

### F-2 (Vorbefund HIGH) — „Spiegelt exakt" — geschlossen: die einseitige Zusage steht an vier Stellen, die Hook-Logik ist unverändert

- `kategorie`: HIGH (Vorbefund, jetzt geschlossen)
- `quelle`: [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
  §Entscheidung Punkt 1/2 · `AGENTS.md` §3.7
- `pfad`: `.githooks/commit-msg:5-9` und `:52-57` · `harness/README.md:165` ·
  `d-check.mk:52-56` · Slice-Plan `:55-60`
- `befund`: Alle vier Träger nennen jetzt die einseitige Zusage statt „denselben
  Verstoß" bzw. „exakt": der Einleitungskommentar („weist keinen Commit
  zurueck, den das Standing-Gate zulaesst; vollstaendig ist er nicht"), der
  Kommentar an der positiven Hälfte („macht seine Weite laxer als die des
  Moduls — die einseitige Zusage traegt das"), die README-Zeile („meldet vorab
  in den Klassen, in denen er liest — er fängt **nicht** jeden Verstoß … und
  weist keinen Commit zurück, den das Gate zulässt") und Plan §1 (drei offene
  Divergenz-Klassen). **Verhaltens-Code unverändert, selbst geprüft:**
  `git show <rev>:.githooks/commit-msg | grep -vE '^\s*#' | grep -vE '^\s*$' |
  sha256sum` ergibt für `517eee2`, `78d4007` und `HEAD` **denselben** Hash
  (`1a2d0ebe…`); die Änderung ist rein kommentierend. Dasselbe für `41e409f`:
  `git show 41e409f -- d-check.mk` enthält keine Nicht-Kommentarzeile.
- `verifizierbar`: ja — Mutations-/sha256-Vergleich am Skript (s. o.),
  `make gates` für die README-Zeile.
- `klasse`: „Spiegelung ist Approximation — belegte Divergenz-Klasse zwischen
  lokalem Hook und dem Modul, das die Regel trägt" (Vorbefund, jetzt
  geschlossen)

### F-3 (Vorbefund LOW) — Nachweis ohne Isolierung — geschlossen, in **beiden** Richtungen

- `kategorie`: LOW (Vorbefund, jetzt geschlossen)
- `quelle`: Slice-Plan §2 DoD-Zeile 1 · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-13-quality-gates.md` (Mutations-Beleg)
- `pfad`: Slice-Plan `:98-103` · `.githooks/commit-msg:45-50`, `:58-65`
- `befund`: DoD-Zeile 1 nennt jetzt den **isolierenden** Input
  (`docs: update SPEC-012 table (ADR-0045)`) und die Aussage, dass ihn allein
  die Grenz-Hälfte zurückweist. Eigene Messung mit **aktivem** Hook
  (`core.hooksPath`, Wegwerf-Repository) bestätigt beide Hälften: der
  isolierende Input wird real abgewiesen (Exit 1, Fehlertext nennt die
  Grenz-Hälfte), die Gegenprobe ohne Kennung ebenfalls (Exit 1, Fehlertext
  nennt die positive Hälfte). Eigene Mutationsmatrix über fünf Proben × vier
  Hook-Varianten: **Grenz-Hälfte aus** → isolierender Input kippt 1→0, während
  die Gegenprobe „Struktur-ID ohne Kennung" rot bleibt (1) — die positive
  Hälfte passiert den isolierenden Input also wirklich; **positive Hälfte aus**
  → die Gegenprobe ohne Kennung kippt 1→0, der isolierende Input bleibt 1;
  **Ausnahme aus** → der Merge-Commit kippt 0→1. Die Isolierung trägt damit in
  beide Richtungen (Grenz-Hälfte von der positiven unterscheidbar und
  umgekehrt).
- `verifizierbar`: ja — Mutations-Lauf am Skript (s. Urteil 3 der
  Sensoren-Tabelle).
- `klasse`: „Nachweis ohne Isolierung der geprüften Hälfte" (Vorbefund, jetzt
  geschlossen)

### F-4 (Vorbefund INFO) — Zeiger auf die Alternativen-Tabelle — geschlossen

- `kategorie`: INFO (Vorbefund, jetzt geschlossen)
- `quelle`: `AGENTS.md` §3.7 (Rang-Zeiger) ·
  [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
- `pfad`: `d-check.mk:52-56`
- `befund`: Der Kommentar über dem `commit-traceability`-Block nennt als Träger
  der Hook-Politik jetzt `ADR-0069` („die Hook-Politik trägt `ADR-0069`
  (Supersedes `ADR-0062` Punkt 3)") statt der Alternativen-Tabelle, deren
  Hook-Aussage korrigiert ist. Die Feststellung des Verdikts, dass
  `commit-traceability` **Adopter-Inhalt in einem generierten Fragment** ist,
  ist eigenständig bestätigt: `docker run … --print-mk` am gepinnten Digest
  emittiert `doc-check`, `doc-trace`, …, `doc-commits`, `doc-planning`,
  `doc-tracked`, `doc-targets`, `doc-structure`, `doc-usage`, `doc-help` — ein
  `commit-traceability`-Target ist **nicht** darunter (76 emittierte Zeilen
  gegen 88 im Repo-Fragment).
- `verifizierbar`: ja — `--print-mk` am gepinnten Digest.
- `klasse`: „Kommentar-Verweis auf eine teilweise abgelöste ADR ohne Zeiger"
  (Vorbefund, jetzt geschlossen)

### F-5 — Der Teil-Supersede beschreibt Punkt 3 schmaler, als er greift; zwei Plan-Zeiger stehen weiter auf dem abgelösten Punkt

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
  §Status („nur deren Entscheidung Punkt 3, ‚der Hook … spiegelt exakt die
  zwei bestehenden Regeln'") gegen
  [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
  §Entscheidung Punkt 3 (**drei** Aussagen: bash-only · exakte Spiegelung ·
  Merge-/Revert-Ausnahme) · `AGENTS.md` §3.7 (Rang-Zeiger)
- `pfad`: `docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md:3-8` ·
  Slice-Plan `:78-81` und `:110-116`
- `befund`: `ADR-0069` schreibt „nur deren Entscheidung **Punkt 3**" superseded,
  zitiert zur Bestimmung dieses Punktes aber allein die Spiegelungs-Aussage;
  Punkt 3 trägt daneben **bash-only/kein Docker-Aufruf** und die
  **Merge-/Revert-Ausnahme** der positiven Hälfte. Keine der beiden Aussagen
  wird in `ADR-0069`s §Entscheidung erneut ausgesprochen (die einzige
  Nennung von „bash-only, keine Latenz" steht als Begründungs-Zitat in
  §Verglichene Alternativen, die Ausnahme nur in §Kontext und im
  §Warum-die-Approximation-Absatz); sie sind aus Punkt 1 („keine
  Hook-Rückweisung, die das Gate nicht auch wirft") bzw. „Die Hook-Logik
  bleibt unverändert" **ableitbar**, nicht ausgesprochen. Der Slice-Plan
  zitiert für beide Aussagen weiter `ADR-0062` **Entscheidung Punkt 3**: in §1
  („der Hook bleibt bash-only (`ADR-0062` Entscheidung Punkt 3)") und in DoD
  Zeile 4 (Merge-/Revert-Ausnahme, „`ADR-0062` Entscheidung Punkt 3"). Je nach
  Lesart ist damit entweder der Supersede-Satz zu weit gefasst (dann bleiben
  die Plan-Zeiger gültig) — oder Punkt 3 ist ganz abgelöst (dann haben zwei
  normative Aussagen keinen ausgesprochenen Accepted-Träger mehr und die zwei
  Plan-Zeiger sind veraltet). Beide Lesarten führen zu einer kleinen Korrektur;
  entschieden werden muss sie beim Träger der Entscheidung, nicht im Plan.
- `verifizierbar`: nein — kein Gate prüft die Reichweite eines Supersedes
  gegen den Wortlaut des abgelösten Punktes; grep- und diff-prüfbar.
- `klasse`: „Supersede-Reichweite schmaler beschrieben als der abgelöste Punkt
  — Zeiger auf den abgelösten Punkt bleibt stehen"

### F-6 — Die benannte Sicherheits-Kante ist in drei Konstruktionen nicht erreichbar; der Grund ist strukturell

- `kategorie`: INFO
- `quelle`: [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
  §Entscheidung Punkt 1 („benannt statt behauptet") · §Fitness Function
- `pfad`: `docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md:86-96` ·
  `.githooks/commit-msg:5-6` · Slice-Plan `:55-57`
- `befund`: Die Kante („liest der Hook unter `--cleanup=scissors` hinter der
  scissors-Zeile, die `git` später verwirft, kann seine Betreff-Zeile Text sein,
  den das Gate nicht mehr liest") habe ich in drei Konstruktionen **nicht**
  erreichbar gemacht: (1) Kommentar-only-Vor-scissors-Text + Struktur-ID hinter
  der scissors → Hook weist zurück (Exit 1), derselbe Commit ohne Hook hat aber
  `commit-untraceable` (Exit 1) — das Modul bereinigt die `#`-Zeilen weg, die
  gespeicherte Message wird leer, der Treffer ist rot; (2) derselbe Aufbau
  hinter der scissors mit leerem Vor-scissors-Text → **git selbst** bricht mit
  „Commit aufgrund leerer Beschreibung ab" (Exit 1), bevor der Hook zählt;
  (3) `git commit -v --cleanup=scissors` mit unberührtem Editor-Ergebnis →
  ebenfalls git-Abbruch „leere Beschreibung". Der Grund ist strukturell und
  nicht bloß Häufigkeit: der Hook liest nur dann hinter die scissors-Zeile,
  wenn **alle** Vor-scissors-Zeilen leer oder `#` sind — genau dann reduziert
  die gespeicherte Message sich auf Leer-/Kommentarzeilen, und das Modul färbt
  rot (verifiziert: `commit-untraceable` mit der `commits`-Konfiguration des
  Repos). In jeder gemessenen Konstruktion ist eine Hook-Rückweisung damit auch
  eine Gate-Rückweisung; die unkonditionierte Fassung der Sicherheits-Zusage im
  Hook-Kommentar und in Plan §1 ist durch die Messung gedeckt, obwohl die ADR
  sie vorsichtshalber konditioniert. Grenze dieser Feststellung: sie beruht auf
  drei Konstruktionen plus dem Quelltext-Lesen von
  `cleanCommitMessage`/`commitSubject` und der Git-Cleanup-Modi — sie ist
  **kein** Unmöglichkeitsbeweis über alle Konfigurationen.
- `verifizierbar`: ja — dieselben Konstruktionen sind mit dem gepinnten Image
  (`--commit-msg`, `--range`) und einem Wegwerf-Repository nachstellbar.
- `klasse`: „Benannte Kante ohne gemessenen Erreichbarkeitsbeleg (hier: nicht
  erreichbar, Kante bleibt vorsichtig benannt)"

### F-7 — Die festgeschriebene Grenz-Liste ist nicht erschöpfend: eine vierte, nur unter `--cleanup=verbatim` erreichbare laxere Klasse existiert

- `kategorie`: INFO
- `quelle`: [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
  §Entscheidung Punkt 2 („Die drei oben benannten Divergenz-Klassen (a), (b),
  (c) bleiben offen und sind die festgeschriebene Grenze") · §Re-Evaluierungs-
  Trigger (a) („eine der **drei** benannten Klassen")
- `pfad`: `docs/plan/adr/0069-commit-msg-hook-einseitige-zusage.md:97-101`,
  `:195-196`
- `befund`: Eigene Messung (Wegwerf-Repository, aktiver Hook, `commits`-Config
  des Repos): ein Commit, dessen rohe Message mit einer **Leerzeile** beginnt,
  auf die `Merge branch 'vp'` folgt, unter `--cleanup=verbatim` — der Hook lässt
  ihn durch (Exit 0, er überspringt die Leerzeile und greift die Ausnahme), die
  gespeicherte Message behält die Leerzeile, das Modul liest als Betreff `""`,
  nimmt die Ausnahme nicht an und meldet `commit-untraceable` (Exit 1): wieder
  eine **laxere** Klasse (Hook grün, Gate rot), und zwar außerhalb der drei
  benannten. Sie ist nur mit einem ausdrücklichen `--cleanup=verbatim` (oder
  `commit.cleanup=verbatim`) erreichbar — unter dem Default-Cleanup strippt
  `git` die Leerzeile, gespeicherter `%s` ist `Merge branch 'vp'`, und beide
  Seiten sind grün (im Vorbefund als H2 verworfen, hier nur als Listenfrage
  erneut gemessen). Der Hook folgt in dieser Klasse git's eigener
  `%s`-Normalisierung; das Modul ist der Ausreißer.
- `verifizierbar`: ja — gepinntes Image (`--range`) gegen einen so erzeugten
  Commit.
- `klasse`: „Grenz-Liste einer ADR als erschöpfend formuliert, obwohl eine
  weitere erreichbare Klasse existiert"

---

## Urteile zu den gestellten Prüfpunkten

**1. F-1 — trägt die Korrektur, und sagt der Plan noch etwas Unwahres?**
**Getragen, ja** (s. F-1): die Entscheidung ist der Träger, die Immutabilität
von `ADR-0062` ist gewahrt, der Index zieht den Zeiger. Zur Frage in **beide**
Richtungen: Der Plan behauptet **nichts Untertriebenes** (die drei Klassen sind
benannt, die message-weite positive Hälfte steht in §1, §2 und §3, die
betreff-scoped Grenz-Hälfte ist überall als solche bezeichnet) und **nichts
Übertriebenes** mehr an der Stelle, die den Vorbefund trug — die einzige
Reserve sind die zwei `ADR-0062`-Punkt-3-Zeiger aus F-5. **Zum entfernten
Satz im `Bezug`-Feld:** Die Streichung von „ohne `ADR-0062` widerspräche dieser
Slice `ADR-0045` stillschweigend" ist die **richtige** Anwendung der
Prosa-Regel, kein Herkunfts-Verlust. Die Regel verlangt Indikativ über den
Zustand statt Konjunktiv über die verworfene Alternative (`AGENTS.md` §3.7),
und dieselbe Tatsache steht jetzt positiv im Feld („sie korrigiert
`ADR-0045`s Klausel ‚Kein commit-msg-Hook' und erlaubt diesen Slice erst");
die `ADR-0069`-Ergänzung fügt den aktuellen zweiten Träger hinzu. Der Link auf
das Gegenprüfungs-Verdikt war nicht der Träger dieser Aussage — das Dokument
ist **nicht** verwaist (Urteil 5).

**2. F-2 — trägt die neue Formulierung, und stimmt die Zusage?** **Ja, beide
Hälften.** Die drei Klassen sind am Artefakt wiederzufinden (Plan §1, ADR §Kontext
und §Fitness Function, Hook-Kommentar als „laxer als das Modul"), und die
Richtung stimmt: eigene Messung in einem Wegwerf-Repository mit aktivem Hook
über fünf Proben × vier Varianten (Urteil 3) findet **keinen** Fall, in dem der
Hook strenger ist; die drei Klassen (a)/(b)/(c) sind reproduzierbar. Die im
letzten Lauf verworfene H2-Klasse ist auch hier **keine** Laxheits-Klasse unter
Default-Cleanup (gemessen: git strippt die Leerzeile, beide grün); sie taucht
nur unter `--cleanup=verbatim` auf (F-7). Die von `ADR-0069` selbst benannte
Kante habe ich eigenständig zu erreichen versucht und **nicht** erreicht (F-6)
— damit ist die Sicherheits-Zusage nicht nur formuliert, sondern in allen
gemessenen Konstruktionen erfüllt. Die README-Zeile trägt beide Hälften in
einem Satz (nicht vollständig · keine Fehlrückweisung) und bleibt bei
`core.hooksPath`/`--no-verify` korrekt.

**3. F-3 — isoliert der Nachweis jetzt in beide Richtungen?** **Ja** — eigene
Mutationsmatrix (Urteil 2/3 der Sensoren-Tabelle, Wegwerf-Repository): der
isolierende Input `docs: update SPEC-012 table (ADR-0045)` wird real abgewiesen
und kippt **nur** mit der Grenz-Hälfte-Mutation auf 0; die Gegenprobe ohne
Kennung kippt **nur** mit der Positiv-Mutation; die Gegenprobe „Struktur-ID
ohne Kennung" bleibt unter beiden Einzelmutationen rot (sie deckt beide
Hälften) und die Merge-Probe kippt nur mit der Ausnahme-Mutation. Die beiden
vom Implementer genannten Mutationen habe ich selbst gesetzt und gesehen; die
dritte (Ausnahme) zusätzlich.

**4. F-4 — Zeiger gezogen, und stimmt die Behauptung zu `--print-mk`?** **Ja,
beides** (s. F-4). Der Zeiger nennt `ADR-0069`; das gepinnte Image emittiert
kein `commit-traceability`-Target (`--print-mk`: 76 Zeilen, `doc-*`-Familie,
kein `commit-traceability`) — der Target ist Adopter-Inhalt im generierten
Fragment, wie das Verdikt sagt.

**5. Verwaist das Gegenprüfungs-Verdikt?** **Nein — bestätigt, die Lesung des
Coordinators trägt.** `grep -rn` über die nicht-vendored Doku findet
`architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md` in
vier Dateien: [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(§Autor, §Geschichte), Slice-Plan (dreimal: Kopf-Feld `Autor` Zeile 37,
§1 Zeile 71, §8 Zeile 208),
`BEO-PGC/commit-traceability-kein-vorab-hook/state.md` (Ausgangs-Beleg) und im
Vorbefund `docs/reviews/review-slice-073.md`. Die Streichung betraf nur das
Kopf-Feld `Bezug`; das Dokument bleibt über den Register-Eintrag **und** über
den Plan erreichbar.

**6. Kein neues Finding aus der Fixrunde selbst — Negativbefunde unten.**
Vorläufig benannt, weil **nicht** Teil dieses Reports: §4 des Plans nennt als
Rückführungs-Bedingung, dass die Kongruenz-Frage „zurück zur Zerlegung" gehört,
„bevor der Hook geschrieben wird". Diese Bedingung ist mit den drei Klassen
faktisch entscheidbar geworden (sie ist eingetreten oder — mit `ADR-0069`
Option B — bewusst nicht gewählt), und Modul 5 verlangt den *Grund* im
Nachhinein. Der Plan trägt dazu (noch) nichts; das ist **Planner-Arbeit** in
§7 und bleibt hier unangetastet, wie §6-Risiko-Ausgänge und das
Beobachtungs-Register.

---

## Negativbefunde

- geprüft, ohne Befund: `.githooks/commit-msg` (Verhaltens-Code sha256-gleich zu
  `517eee2`; die neuen Kommentare tragen Zusage/Kopplung/Abgrenzung, keine
  Chronik, kein Konjunktiv über eine verworfene Alternative — `AGENTS.md` §3.7)
- geprüft, ohne Befund: `harness/README.md:165` (beide Hälften der Zusage,
  Opt-in, Nicht-Ersatz, `--no-verify`; `make docs-check` grün)
- geprüft, ohne Befund: `d-check.mk` (nur Kommentarzeilen geändert; kein
  Verhaltens-Change, Target-Block unberührt)
- geprüft, ohne Befund: Slice-Plan §1/§2/§3 (message-weite positive Hälfte,
  isolierender Nachweis, Rang-Zeiger; keine Betreff-Scope-Behauptung mehr für
  die positive Hälfte)
- geprüft, ohne Befund: ADR-Immutabilität (`ADR-0062` unverändert, `ADR-0045`
  unverändert, `ADR-0069` Accepted mit Supersedes-Zeiger; Index-Zeilen für alle
  drei vorhanden)
- geprüft, ohne Befund: Muster-Äquivalenz der Regexe (Stichproben
  `ADR-045`/`LH-FA-CFG-005.a`/`LH-FA-CFG-5`/`LH-QA-POR-001`: identische
  Exit-Codes von Hook und Modul)
- geprüft, ohne Befund: Nicht-durchsetzend (kein Gate-Ziel, kein Docker,
  `--no-verify` umgeht, Opt-in) — am Artefakt unverändert bestätigt

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Supersede-Reichweite schmaler beschrieben als
der abgelöste Punkt — Zeiger auf den abgelösten Punkt bleibt stehen" (neu) ·
„Benannte Kante ohne gemessenen Erreichbarkeitsbeleg" (neu) · „Grenz-Liste
einer ADR als erschöpfend formuliert" (neu). Die vier Vorbefunde aus
`docs/reviews/review-slice-073.md` sind **geschlossen**.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. Die beiden HIGH des Vorbefunds sind
substanziell geschlossen, nicht kosmetisch: der Träger der Abweichung ist von
der Plan-Zelle in eine Accepted-Folge-ADR gewandert (F-1), und die Zusage des
Hooks ist an allen vier Trägern auf die einseitige Form umgestellt, während der
Verhaltens-Code nachweislich unverändert ist (F-2, sha256). Der Nachweis der
Grenz-Hälfte isoliert jetzt in beide Richtungen (F-3), der Zeiger aus F-4 ist
gezogen, und die Feststellung zur generierten Fragment-Herkunft trifft zu
(Urteil 4).

**Übergabe:** **Keine Rückgabe an den Implementer** — kein Finding dieses Laufs
trägt einen Reviewer→Implementer-Pfeil. F-5s entscheidende Hälfte ist eine
Architect-Frage (Reichweite des Teil-Supersedes; der Plan-Zeiger folgt der
Antwort), F-6 und F-7 sind Feststellungen ohne erwartete Aktion an das Artefakt
— die Auflösung von F-6/F-7 (falls gewünscht) liegt beim Träger der ADR, nicht
im Hook. Der Slice kann aus Reviewer-Sicht an die **Verifier**-Rolle übergeben
werden (DoD-/Spec-Konformität, Modul 11).

**DoD-Nachzug:** Da 0 HIGH offen sind und kein Finding einen
Implementer-Rückgabe-Pfeil trägt, zieht dieser Report die Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des Slice-Plans
selbst auf `[x]` nach, mit Verweis auf diesen Report, im selben Commit
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). **Nur
diese eine Zeile** (Grenze des Skill-Abschnitts): Closure-Notiz,
Beobachtungs-Register, §6-Risiko-Ausgänge und die drei Paarungen bleiben
Planner-Arbeit.

---

## Ausgeführte Sensoren (Exit-Codes dieses Laufs)

| Aufruf | Exit-Code | Beleg |
|---|---|---|
| `make gates` | 0 | `docs-check` 588 Dateien / 0 Befunde · `commit-traceability` OK · `a-check` 0 · `coverage-gate` 49,30 % ≥ 35 % · `baseline-verify` OK |
| `make commit-traceability` | 0 | `HEAD~5..HEAD`: d-check 0 Befunde · „Betreffs ohne Struktur-ID" |
| Verhaltens-Code-Vergleich des Hooks | — | `sha256` der Nicht-Kommentar-Zeilen für `517eee2`/`78d4007`/`HEAD` identisch (`1a2d0ebe…`) |
| eigene Mutationsmatrix (5 Proben × 4 Hook-Varianten) | 1 / 1 / 1 / 0 / 0 | Wegwerf-Repository, aktiver Hook; Varianten: Original · Grenz-Hälfte aus (A 1→0) · positive Hälfte aus (C 1→0) · Ausnahme aus (Merge 0→1) |
| zwei reale Commits mit dem isolierenden Input | 1 / 1 | Hook aktiv, kein Commit angelegt; Fehlertexte nennen je die geprüfte Hälfte |
| Kante-Konstruktionen E1/E2/E3 | Hook 1 · Gate 1 (E1) bzw. git-Abbruch (E2/E3) | Kante nicht erreichbar (F-6) |
| vierte Klassen-Messung (`--cleanup=verbatim`) | Hook 0 · Gate 1 | Wegwerf-Repository, `commits`-Config des Repos (F-7) |
| `--print-mk` am gepinnten Digest | 0 | 76 Zeilen, kein `commit-traceability` (F-4) |

Alle Messungen und Mutationen liefen ausschließlich außerhalb des Arbeitsbaums
(`/tmp/review073/`); `git status --porcelain` ist vor dem Anlegen dieses Reports
leer bis auf die DoD-Zeile dieses Commits.
