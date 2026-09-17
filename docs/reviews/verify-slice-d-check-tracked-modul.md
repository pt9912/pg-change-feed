# Verifikations-Bericht: slice-d-check-tracked-modul — 2026-09-17

**Rolle:** Verifier (Baseline-Regelwerk `v6.9.0`
`regelwerk/modul-11-verification.md`) — „Bauen wir es richtig?" gegen **DoD und
Entscheidungen**, nicht gegen Plan-Maintainability (das war der Reviewer).

**Gegenstand:** `slice-d-check-tracked-modul` (`d-check`-Modul `tracked`
scharf geschaltet, Welle `welle-d-check`). Diff `314b6cc..1d9db37`, sieben
Commits: `7c9bb1b` (Modul aktiviert) · `a276838` (Träger-Nachzug) ·
`feacf52` (DoD-Checkboxen + Plan-Nachzug) · `550a959` (Review, 2 HIGH/1
MEDIUM/1 INFO) · `702e114` (Architect-Verdikt) · `5800a83` (Fixrunde) ·
`1d9db37` (Delta-Review Fixrunde, 0 HIGH).

**Plan:**
`docs/plan/planning/in-progress/slice-d-check-tracked-modul.md`
(vollständig, inkl. Plan-Nachzug `feacf52`/`5800a83`).

**Entscheidungen:** keine eigene ADR — der Slice ist bewusst ADR-frei
angelegt (§`Bezug`); die ADR-Frage aus §6 Risiko 1 ist über den nachträglichen
Architect-Zug
([`architect-verdict-slice-d-check-tracked-modul-adr-frage.md`](architect-verdict-slice-d-check-tracked-modul-adr-frage.md))
mit Ausgang „entfallen" beantwortet, tragender Präzedenzfall Commit `f9e5a3c`.

**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17.

**Eigener Eingabe-Kontext** (frisch, ohne Implementer-/Reviewer-/
Architect-Zusage): Plan §1–§8 · beide Review-Reports
([`review-slice-d-check-tracked-modul.md`](review-slice-d-check-tracked-modul.md),
[`review-slice-d-check-tracked-modul-fixrunde.md`](review-slice-d-check-tracked-modul-fixrunde.md))
· das Architect-Verdikt · `AGENTS.md` §3.5, §3.9, §3.12, §3.13 · `.d-check.yml`
· `harness/sensors/docs-check.md` · `harness/README.md` §Sensors ·
`.claude/agents/verifier.md`/`implementer.md` · `git log`/`git show` über
`314b6cc..HEAD` sowie den zitierten Präzedenzfall `f9e5a3c`. Alle Messungen
dieses Berichts sind **eigene** Läufe (Exit-Codes ungepiped, `AGENTS.md`
§3.9), keine Übernahme der Implementer-/Reviewer-Behauptung.

---

## A. Eigene Sensoren (Belege, keine Behauptungen)

| Sensor | Kommando | Ergebnis | Exit |
|---|---|---|---|
| `make gates` (repo-weit, HEAD `1d9db37`) | `make gates` (ungefiltert, Exit-Code direkt geprüft) | `docs-check` 878 Dateien/0 Befunde · `commit-traceability` 5 Commits OK · `coverage-gate` 83,40 % ≥ 80 % · `generated-sync` byte-gleich · `a-check` 0 Befunde | **0** |
| `make doc-tracked` (isoliert, HEAD) | `--enable tracked --disable <übrige>` | 878 Dateien, 0 Befunde — deckungsgleich mit dem Vollbündel | **0** |
| `make doc-commits RANGE=314b6cc..HEAD` | — | `--enable commits`, 878 Dateien, 0 Befunde | **0** |
| `make doc-immutable RANGE=314b6cc..HEAD` | — | `--enable vcs`, 878 Dateien, 0 Befunde | **0** |
| Commit-Betreffs der sieben Slice-Commits, einzeln gelesen (`git show -s --format=%s`) | — | alle sieben tragen `ADR-0045` im Betreff, keine `SPEC-*`/`ARC-*`-Kennung | — |
| Zitierter Präzedenzfall `f9e5a3c`, selbst gelesen (`git show -s`, `git show --stat`) | — | Datum/Charakterisierung stimmen mit dem in `harness/sensors/docs-check.md` zitierten Text überein: nur `.d-check.yml` (68 Zeilen), kein ADR-Bezug, kein `AGENTS.md`-Bezug | — |
| Repo-weiter Grep gegen die „sieben Module"-Aufzählung (`grep -rln "links, anchors, ids, matrix, versions" --include="*.md" .`, ohne `.harness/baseline/`) | — | 15 Fundstellen: die drei lebenden Träger (`harness/README.md`, `.claude/agents/{verifier,implementer}.md`) tragen bereits `tracked`; die übrigen 12 sind `Accepted`-ADRs (`0072`, `0089`) oder Records unter `docs/reviews/**` — beide Klassen sind unveränderliche Zeitdokumente (`AGENTS.md` §3.5 bzw. Reviewer-Skill „Records tragen keine §Geschichte") | — |
| `git status --porcelain` nach `make gates` | — | leer | — |

---

## B. DoD-Konformität Punkt für Punkt

### LP1 — `tracked`-Block + Bündel-Aufnahme

**Erfüllt.** `.d-check.yml:8`: `modules: [links, anchors, ids, matrix,
versions, structure, hostpaths, tracked]`. `.d-check.yml:213-220` trägt den
`tracked:`-Block mit `exempt-targets: []` und einer Begründung, warum leer
(„keine absichtlich untrackten, aber verlinkten Ziele im Bestand
(geprüft)"). `make docs-check` läuft **real** am aktuellen HEAD (nicht die
Planungs-Vorabmessung übernommen, `AGENTS.md` §3.12) — Exit 0, 878/0 im
Rahmen von `make gates` (§A).

### LP2 — `harness/sensors/docs-check.md` nachgezogen

**Erfüllt.** §Grenze Punkt 3 (`:44-49`) nennt `tracked` nicht mehr als
Opt-in, sondern als „seit `slice-d-check-tracked-modul` Teil des
`modules:`-Bündels"; die verbleibenden vier (`planning`, `vcs`, `commits`,
`reviews`) bleiben stehen. Punkt 4 (`:51-53`) ist konsistent mitgezogen
(„nicht mehr Opt-in"). Die §Bindung-Sektion (`:174-179`) trägt einen
`tracked`-Eintrag mit Modul, Konfigurationsort und — nach der Fixrunde —
dem **korrigierten** Präzedenz-Beleg (Commit `f9e5a3c`, nicht mehr
`ADR-0072`/`ADR-0075`; siehe §D unten). `.d-check.yml`s eigener
Kommentarblock (`:213-220`) ist Aktivierungsstand, nicht mehr
Beispiel-Kommentar.

### LP3 — `make gates` grün

**Erfüllt.** Eigener Lauf am HEAD `1d9db37`, ungepiped, Exit-Code direkt
geprüft: **Exit 0** (§A). Die welle-weite Zusatz-Bedingung (kombinierter
Lauf mit `slice-d-check-trace-rtm`) ist laut Plan **nicht** Teil dieser
Einzel-DoD, sondern der Welle-Closure — hier nicht geprüft, korrekt nicht
Teil dieses Slices.

### Review durchgeführt, Report liegt vor

**Erfüllt.** Ausgangs-Report `550a959` (2 HIGH, 1 MEDIUM, 1 INFO) →
Fixrunde `5800a83` → Delta-Review `1d9db37` (0 HIGH, 0 MEDIUM, nicht
merge-blockierend). Eigene Prüfung der drei Findings unten (§D) bestätigt,
dass die Fixrunde sie trägt, nicht nur behauptet.

### Doku-Update

**Erfüllt.** `harness/sensors/docs-check.md` (LP2) und `.d-check.yml`
eigener Kommentarblock — beide bestätigt oben. Zusätzlich (Plan-Nachzug
`feacf52`, außerhalb des ursprünglichen §3, aber im Slice-Umfang):
`harness/README.md:128` führt `make doc-tracked` korrekt in der
„kein Gate"-Werkzeuge-Tabelle, nicht in der Gates-Tabelle — eigene Prüfung
von `grep -n "GATE_CHECKS" Makefile harness/mk/*.mk` bestätigt, dass
`doc-tracked` nirgends an `GATE_CHECKS` hängt.

### Closure-Notiz mit Steering-Loop-Lerneintrag

**Offen — korrekt unchecked.** §7 des Plans trägt ausschließlich
Platzhalter (`<…>`). Zu diesem Zeitpunkt (Slice steht erst jetzt zur
Closure an) erwartungsgemäß noch nicht ausgefüllt — keine DoD-Verletzung,
sondern die nächste Planner-Aufgabe.

### Beobachtungs-Register fortgeschrieben

**Offen — korrekt unchecked.** Die im Plan (§2, §8) gestellte Frage („trägt
dieser Vorgang eine weitere `evidence/`-Datei zu
`BEO-PGC/gate-modul-abgeschaltet-trotz-regel` oder eine eigene
Beobachtung, oder keine?") ist noch nicht beantwortet. Eigene Prüfung:
`docs/plan/planning/observations/BEO-PGC/gate-modul-abgeschaltet-trotz-regel/`
trägt bislang genau die eine `hostpaths`-Evidenz aus `slice-078`; keine neue
Datei zu diesem Slice ist angelegt. Die im Plan selbst genannte
Unterscheidung (bei `tracked` war das Modul dokumentiert Opt-in, nicht
stillschweigend abgeschaltet wie bei `hostpaths`) ist ein plausibles
Argument für „andere Klasse, eigene oder keine Beobachtung" — das ist aber
eine Planner-Wertung, keine, die der Verifier vorwegnimmt.

### Jedes Risiko aus §6 trägt einen Ausgang

**Teilweise — korrekt unchecked, siehe §C unten.** Risiko 1 trägt einen
Ausgang mit tragendem Beleg (siehe §C); Risiko 2 und Risiko 3 tragen im
committeten Stand weiterhin ihre Platzhalter-Form
(`<weiter offen: … / entfallen: …>` bzw. `<Reviewer prüft … — eingetreten: …
/ entfallen: …>`) — nicht ausgefüllt. Die Fakten für beide liegen jedoch
bereits vor (siehe §C), sie sind nur noch nicht in den Plan geschrieben.

### Drei Paarungen

**Korrekt unchecked, delegiert.** Plan-Text weist das der Welle-Closure zu
(„geprüft von der Welle-Closure … nicht hier einzeln") — konsistent mit
dem repo-weiten Wellen-Betrieb.

---

## C. §6-Risiko-Ausgänge — eigene Prüfung

### Risiko 1 — „Aufnahme in `modules:` verlangt eine eigene ADR"

**Ausgang „entfallen" trägt.** Eigene, unabhängige Prüfung (nicht Übernahme
des Architect-Verdikts):

- Der zitierte Präzedenzfall `f9e5a3c` ist real und trägt die behaupteten
  Eigenschaften (§A: nur `.d-check.yml`, kein ADR-Bezug, keine
  `AGENTS.md`-Änderung im selben Commit).
- Der Vergleich mit `hostpaths` (`ADR-0072`/`ADR-0075`) ist korrekt
  gegenläufig: `ADR-0072` §Entscheidung Punkt 1 sagt wörtlich, der Träger
  sei „diese ADR, keine Aktennotiz" — das Dokument selbst benennt sich als
  notwendig, nicht als Beleg für „keine ADR nötig". Selbst gelesen, nicht
  übernommen.
- Die drei unterscheidenden Merkmale (reale Befunde am Aktivierungsstand,
  Korrektur an `Accepted`-Inhalten, neue Hard Rule) treffen auf `tracked`
  in **keinem** Fall zu (0 Befunde, kein ADR-Text berührt, kein
  `AGENTS.md`-Vorschlag) — das deckt sich mit `structure`, nicht mit
  `hostpaths`.
- Die drei Ausgänge (eingetreten/entfallen/weiter offen) sind eine
  geschlossene Menge (Modul 5); „entfallen" ist hier die korrekte Wahl, weil
  das Risiko („braucht eine eigene ADR") sich als **nicht zutreffend**
  herausgestellt hat, nicht als „noch ungeklärt" oder „real eingetreten".

**Einziger Vorbehalt (kein Blocker):** Der Architect-Zug kam laut Plan §4
*nach* statt *vor* Implementierungsbeginn (F-3 des Ausgangs-Reviews). Das
ist im Plan als Record stehen geblieben und explizit **nicht** rückwirkend
geheilt (Architect-Verdikt §7) — korrekt: die Sequenz-Verletzung selbst ist
kein DoD-Punkt, sondern ein Prozess-Fakt, den der Plan unverändert trägt.

### Risiko 2 — „`exempt-targets` bleibt bei Planung unbekannt"

**Faktenlage stützt „entfallen", aber der Plan trägt das noch nicht.**
Eigene Messung: `make doc-tracked` (isoliert) und `make docs-check`
(Vollbündel) liefern am HEAD identisch **0 Befunde** — der reale
`modules:`-Lauf weicht nicht vom isolierten Planungslauf ab. Das
Delta-Review bestätigt dasselbe für den Fixrunden-Stand. Die
Tatsachengrundlage für „entfallen" liegt vollständig vor; der Plan-Text
selbst ist aber noch nicht auf einen der beiden Ausgänge festgelegt
(Platzhalter im committeten Stand).

### Risiko 3 — „`docs-check.md`s Grenze-Aufzählung nicht vollständig
nachgezogen"

**Faktenlage stützt „entfallen", aber der Plan trägt das noch nicht.**
Eigene Prüfung der gesamten Datei `harness/sensors/docs-check.md` (nicht
nur Punkt 3): Punkt 4 ist konsistent mitgezogen, die §Bindung-Sektion trägt
den korrigierten `tracked`-Eintrag, keine weitere Fundstelle mit
veraltetem Opt-in-Stand. Auch hier: Tatsachenlage vollständig, Plan-Feld
noch Platzhalter.

**Einordnung:** Diese beiden offenen Platzhalter sind konsistent mit der
unchecked DoD-Checkbox „Jedes Risiko aus §6 trägt einen Ausgang" — kein
Widerspruch zwischen Plan-Zustand und Checkbox-Zustand. Sie sind
Closure-Arbeit, kein Verifikations-Defekt.

---

## D. Plan-vs-Code-Diff über den gesamten Slice-Lauf

- **Review F-1 (Träger-Nachzug, HIGH) → Fixrunde `5800a83`.** Eigene
  Gegenprobe über den **gesamten** Repo-Baum (nicht nur die drei vom Review
  benannten Dateien): `grep -rln "links, anchors, ids, matrix, versions"
  --include="*.md" .` (ohne `.harness/baseline/`) findet 15 Treffer; die
  drei lebenden Träger tragen bereits `tracked`, die übrigen 12 sind
  `Accepted`-ADRs oder Records unter `docs/reviews/**` — beides
  unveränderliche Zeitdokumente, kein offener Träger übersehen. F-1 ist
  vollständig geschlossen.
- **Review F-2 (Beleg trägt seinen Satz nicht, HIGH) → Architect-Verdikt +
  Fixrunde.** Eigene Lektüre von `ADR-0072` §Entscheidung Punkt 1 bestätigt:
  Das Dokument benennt sich selbst als notwendigen Träger, nicht als
  Präzedenz für „keine ADR nötig". Der Fixrunden-Text in
  `harness/sensors/docs-check.md:174-179` zitiert jetzt korrekt `f9e5a3c`
  statt `ADR-0072`/`ADR-0075`. F-2 ist geschlossen, mit dem **richtigen**
  Beleg, nicht nur formal.
- **Review F-3 (Trigger ohne Übergabe-Artefakt, MEDIUM) → Architect-Verdikt
  als nachgetragenes Artefakt.** Das Verdikt-Dokument selbst benennt sich
  als „der fehlende Zug, nachträglich vollzogen" (§7) und heilt die
  Sequenz-Verletzung ausdrücklich **nicht** rückwirkend — der Plan behält
  die ursprüngliche Formulierung als Record. Kein Widerspruch zu `AGENTS.md`
  §3.7 (Zustandsfeld ohne Chronik): Das §6-Risiko-1-Feld nennt Zustand
  (`entfallen`) und Beleg-Anker (Link), keine erzählte Historie.
- **F-4 (INFO, Präzedenzfall) → verkörpert.** Der `structure`-Präzedenzfall
  ist jetzt der tragende Beleg in `harness/sensors/docs-check.md` — eigene
  Prüfung des Commits (§A) bestätigt die Charakterisierung unabhängig.
- **Keine über die vier Findings hinausgehende Abweichung** zwischen Plan
  und committetem Code gefunden: `.d-check.yml`, `harness/README.md`,
  `harness/sensors/docs-check.md` entsprechen jeweils dem in §3 des Plans
  vorgesehenen Änderungsumfang; keine Änderung an
  `slice-d-check-trace-rtm` oder anderen Out-of-Scope-Bereichen.

---

## E. AGENTS.md §3.12 — Herkunft von Aussagen

- **„Acht `.d-check.yml`-Module"** (`.claude/agents/verifier.md:40-41`,
  `.claude/agents/implementer.md:44-45`, `harness/README.md:114`): Die Zahl
  steht in jedem Fall **in derselben Aussage** unmittelbar neben ihrer
  eigenen Enumeration (`links, anchors, ids, matrix, versions, structure,
  hostpaths, tracked`) — die Herleitung ist für den Leser selbst
  nachvollziehbar (abzählen), keine verdeckte Behauptung. Das deckt sich mit
  der im Repo gelebten Praxis (z. B. der Ausgangs-Review selbst: „sieben
  Module nennt: …" mit sofortiger Aufzählung). Kein §3.12-Verstoß.
- **`make doc-tracked`-Dateizahlen über die Stände:** 875 (Review-Lauf,
  `550a959`) → 876 (Architect-Verdikt, eigener Lauf, Stand `550a959`,
  explizit als „Selbst gemessen (nicht übernommen)" markiert, Differenz zum
  Review durch die hinzugekommene Review-Datei selbst erklärt) → 877
  (Fixrunde, eigener Lauf, Differenz ebenfalls erklärt) → 878 (dieser
  Bericht, eigener Lauf am HEAD `1d9db37`). Jede Zahl trägt ihren Lauf
  (Commit-Stand) und ihre Erklärung der Differenz zur Vorzahl — **gemessen**,
  nicht übernommen, in jedem der vier Berichte. Kein §3.12-Verstoß.
- **Der Präzedenz-Beleg „`f9e5a3c`"** in `harness/sensors/docs-check.md`
  trägt seinen Anker (Commit-Hash + Link auf das Architect-Verdikt) — eine
  Tatsachenbehauptung mit Beleg-Anker (Instanz B), keine ungeprüfte
  Übernahme.

---

## F. Entscheidungs-Konformität

- Keine `Accepted`-ADR wird durch diesen Slice inhaltlich verändert —
  `git diff 314b6cc..HEAD -- docs/plan/adr/` ist leer (eigene Prüfung).
- `make doc-immutable RANGE=314b6cc..HEAD` → Exit 0 (§A): keine `MR-*`-Kerndatei
  verändert.
- `make doc-commits RANGE=314b6cc..HEAD` → Exit 0 (§A): alle sieben Commits
  traceable.
- Der Architect-Zug selbst ist entscheidungskonform eingeordnet (Modul 8,
  Rückkante Review → Plan-Defekt, nicht einer der drei Konflikt-Pfad-Fälle im
  engeren Sinn) — eigene Lektüre von §6 des Verdikts bestätigt diese
  Einordnung als plausibel: keine `Accepted`-Entscheidung stand gegen eine
  Plan-Behauptung, der Plan hatte die Frage bewusst offengelassen.

---

## G. Offene Closure-Obliegenheiten (benannt, nicht ausgeführt)

1. **§6-Risiko-Ausgänge 2 und 3** — Faktenlage liegt vor (§C), Plan-Text
   noch Platzhalter. Empfehlung: „entfallen" für beide, mit Verweis auf
   diesen Bericht oder die bereits vorliegenden Messungen im Review/Delta-
   Review.
2. **Beobachtungs-Register-Entscheidung** — ob eine weitere `evidence/`-Datei
   zu `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` oder eine eigene
   Beobachtung entsteht, oder keine — Planner-Wertung, hier nicht
   vorweggenommen.
3. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag — noch nicht
   geschrieben.
4. **`git mv` nach `done/`** (eigener Commit nach dem Inhalts-Commit,
   `AGENTS.md` §3.3) — noch nicht erfolgt (Datei liegt weiterhin unter
   `in-progress/`).
5. **Drei Paarungen** — von der Welle-Closure `welle-d-check` zu prüfen,
   nicht von diesem Slice allein.

---

## H. Verdikt

**DoD-Konformität (Implementierungs-Teil, LP1–LP3 + Review + Doku-Update):
erfüllt.** Alle fünf `[x]`-Checkboxen sind am Artefakt belegt, nicht nur
behauptet: `.d-check.yml` trägt den `tracked`-Block korrekt, `make gates`
läuft am HEAD grün (eigener, ungepipter Exit-0-Lauf), die
Sensor-Dokumentation ist vollständig nachgezogen (repo-weite Gegenprobe,
nicht nur die drei vom Review benannten Dateien), und die zwei
Review-Runden schließen sauber (0 HIGH im Delta-Review, eigene Prüfung
aller drei behobenen Findings bestätigt inhaltliche Korrektheit, nicht nur
formalen Abschluss).

**§6-Risiko 1 („eigene ADR nötig"): Ausgang „entfallen" trägt.** Eigene,
vom Review/Architect-Verdikt unabhängige Prüfung des Präzedenzfalls
`f9e5a3c` und des Gegenbelegs `ADR-0072`/`ADR-0075` bestätigt die
Schlussfolgerung.

**Commit-Traceability: erfüllt.** Alle sieben Commits tragen `ADR-0045`,
keine `SPEC-*`/`ARC-*`-Kennung im Betreff; `doc-commits`/`doc-immutable`
über die volle Slice-Range beide Exit 0.

**AGENTS.md §3.12: kein Verstoß.** Zahlen tragen ihren Lauf/ihre
Enumeration in jedem geprüften Fall.

**Nicht closure-bereit — erwartungsgemäß.** Vier DoD-Checkboxen bleiben
unchecked: Closure-Notiz, Beobachtungs-Register-Entscheidung, vollständige
§6-Risiko-Ausgänge (2 von 3 fehlen), `git mv` nach `done/`. Das ist **keine
DoD-Verletzung**, sondern der erwartete Zustand unmittelbar vor der
Closure — dieser Bericht ist der Verifikations-Schritt, der dem
Planner-Closure-Zug vorausgeht (Modul 8). Die Faktenlage für die zwei
verbleibenden Risiko-Ausgänge liegt bereits vollständig vor (§C, §G).

**Übergabe an den Planner:** §G-Obliegenheiten, insbesondere die zwei
noch nicht geschriebenen Risiko-Ausgänge (Fakten liegen vor) und die
Closure-Notiz.

---

Der Bericht ist ein **Lauf-Beleg** dieses Verifier-Laufs (dieser Diff, diese
Sensoren, dieses Modell) und ersetzt weder Review noch Closure.
