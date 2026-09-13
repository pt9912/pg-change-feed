# Review-Report: slice-039 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-039`, §1/§2/§3/§4/§6/§8)
und `ADR-0051` (Accepted, Entscheidungen 2/5/9 zentraler Prüfmaßstab dieses
Laufs) sowie `AGENTS.md` §3 Hard Rules (§3.1, §3.7, §3.8) und `ADR-0045`
(Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `37d0d5a1fda253e2f70984a5793a84c5887cda32`
(`feat(ci): CI-Workflow und Dependabot anlegen (ADR-0051)`) — neu:
`.github/workflows/ci.yml`, `.github/dependabot.yml`; geändert:
`harness/README.md` (§Aktueller Lauf-Status),
`docs/plan/planning/in-progress/slice-039-ci-workflow-dependabot.md`
(DoD-Häkchen 1–4/6, Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-039-ci-workflow-dependabot.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug, §4
  Trigger, §6 Risiken, §8)
- `docs/plan/adr/0051-cicd-pipeline-github-actions.md` (vollständig —
  zentraler Prüfmaßstab, insb. Entscheidungen 2, 5, 9)
- `docs/plan/adr/0045-commit-traceability-standing-gate.md` (bindet die
  Dependabot-Commit-Message-Form)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.8
  (Action-Pinning)
- `tools/harness/commit-traceability.sh`, `.d-check.yml` §`commits`
  (`id-patterns`) — geprüft, ob `[ADR-0051]` als Dependabot-Commit-Prefix
  beide Gate-Hälften erfüllt
- `Makefile` (`mod-download`, `test`, `gates`-Kette) — geprüft, ob der
  Workflow ausschließlich Docker-only-Targets aufruft
- Alle neuen/geänderten Dateien vollständig gelesen, nicht nur die
  Implementer-Zusammenfassung: `.github/workflows/ci.yml`,
  `.github/dependabot.yml`, `harness/README.md`-Diff (vollständig),
  `slice-039`-Diff (vollständig)
- Zum Vergleich herangezogen (Bestand, nicht Teil des Diffs, nur
  Einordnung): `/Development/d-check/.github/workflows/ci.yml`,
  `/Development/d-check/.github/dependabot.yml`
- `docs/reviews/review-slice-037.md` (Format-Vorlage)

---

## Findings

### F-1 — Checkout-Action-Pin real verifiziert

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.8
- `pfad`: `.github/workflows/ci.yml:47`
- `befund`: Der SHA
  `3d3c42e5aac5ba805825da76410c181273ba90b1` wurde gegen die GitHub-API
  geprüft (`GET /repos/actions/checkout/git/refs/tags/v7.0.1` und
  `GET /repos/actions/checkout/commits/<sha>`): Der Tag `v7.0.1` zeigt
  exakt auf diesen (signiert verifizierten) Commit. Kein erfundener Pin —
  §3.8 ist real erfüllt, nicht nur textlich plausibel.
- `verifizierbar`: ja (GitHub-API, oben ausgeführt)
- `klasse`: „Action-Pin real gegen Upstream verifiziert" (Negativbefund
  im positiven Sinn — hier als INFO dokumentiert, weil es über die reine
  Plan-Konformität hinaus einen zusätzlichen, in dieser Sitzung
  eigenständig erbrachten Beleg trägt)

### F-2 — Plan-Nachzug behauptet ein "sauberes" `yamllint`-Ergebnis, das sich nicht reproduzieren lässt

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Plan-Konformität — DoD-Item 3 „Beide
  YAML-Dateien real syntaktisch valide", Plan-Nachzug §3, Punkt „YAML-
  Realprüfung")
- `pfad`: `docs/plan/planning/in-progress/slice-039-ci-workflow-dependabot.md`
  (Plan-Nachzug, Punkt 3) i.V.m. `.github/workflows/ci.yml:1,27,47` und
  `.github/dependabot.yml:14`
- `befund`: Der Plan-Nachzug behauptet für den `yamllint`-Lauf: „sauber
  bis auf eine Kommentar-Abstand-Warnung, behoben". Ein realer
  `yamllint`-Lauf (Default-Konfiguration — das Repo führt keine
  `.yamllint`-Datei, also genau die Konfiguration, die das benannte
  Werkzeug ohne weitere Angaben nutzt) gegen exakt die committeten
  Dateien liefert stattdessen: `ci.yml:1:1` fehlender `---`-Dokumentstart
  (Warning), `ci.yml:27:1` `truthy value`-Warnung (der von der eigenen
  Nachzug-Notiz an anderer Stelle beschriebene PyYAML-„Norway
  Problem"-Effekt auf den `on:`-Key), `ci.yml:47:81` Zeile zu lang (81 >
  80 Zeichen, **Error**, an der SHA+Tag-Kommentar-Zeile des
  Checkout-Pins) und `dependabot.yml:14:1` derselbe fehlende
  Dokumentstart. Keiner dieser vier Befunde ist eine
  Kommentar-Abstand-Warnung; keiner ist als behoben dokumentiert. Die
  Diskrepanz wurde in dieser Sitzung durch einen eigenständigen
  `yamllint`-Lauf gegen den committeten Stand hergestellt, nicht aus dem
  Implementer-Bericht übernommen.
- `verifizierbar`: ja — `yamllint .github/workflows/ci.yml
  .github/dependabot.yml` (kein Repo-Sensor, aber genau das im DoD-Item
  selbst benannte Werkzeug, ohne Sonderkonfiguration reproduzierbar)
- `klasse`: „Plan-Nachzug-Verifikationsnachweis nicht reproduzierbar"
  (erstes Auftreten dieser Klasse — verwandt, aber nicht identisch mit
  den bereits verkörperten Mustern `BEO-PGC/plan-nachzug` [Umfangs-Drift
  der §3-Datei-Liste] und `BEO-PGC/dod-checkbox-nachzug` [Checkbox nicht
  nachgezogen]: hier ist die Checkbox korrekt gesetzt und der Umfang
  stimmt, aber der im Nachzug behauptete Werkzeug-*Befund* selbst weicht
  vom real reproduzierbaren Ergebnis ab)

### F-3 — Folge-Slice `slice-040` noch nicht als Datei angelegt

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form:
  Slice (Out-of-Scope-Disziplin, Klasse 1 „Ein Folge-Slice übernimmt es")
- `pfad`: `docs/plan/planning/in-progress/slice-039-ci-workflow-dependabot.md`
  §1, Ausschlusspunkt 2
- `befund`: §1 benennt `slice-040` als Adresse für
  `release.yml`/`hub-description.yml`/`image-scan.yml`/`upstream-drift.yml`/
  `version.md`. Im gesamten Planning-Lifecycle (`open/`, `next/`,
  `in-progress/`, `done/`) existiert aktuell keine Datei mit diesem
  Namen — die Adresse hat noch keine Sendung angenommen. Für diesen
  Slice selbst ist das keine Verletzung (die Datei entsteht planbar
  erst mit dem nächsten Slice-Anlauf, `ADR-0051`s Folgepflicht nennt
  beide Slices ausdrücklich als Zwei-Slice-Aufteilung); der Hinweis ist
  ein Merkposten für die Closure/Folge-Planung, kein aktueller Mangel.
- `verifizierbar`: ja (`find docs/plan/planning -iname "*slice-040*"` —
  kein Treffer, Stand dieser Sitzung)
- `klasse`: „Folge-Slice-Adresse noch nicht angelegt" (Hinweis ohne
  erwartete Aktion an diesem Slice)

## Negativbefunde

- geprüft, ohne Befund: **`permissions: {}` auf Workflow-Ebene,
  `contents: read` auf Job-Ebene** — beide exakt wie im DoD-Item und in
  `ADR-0051` gefordert gesetzt (`ci.yml:19`, `ci.yml:25`).
- geprüft, ohne Befund: **`tags-ignore: ['**']` verhindert Doppelläufe
  bei Tag-Pushes.** Kombiniert mit `branches: ['**']` unter `push:` — die
  Kombination Branch-Filter + `tags-ignore` ist die dokumentierte
  GitHub-Actions-Idiomatik (Branch-Achse und Tag-Achse werden unabhängig
  ausgewertet); `actionlint` (realer GHA-Parser, in dieser Sitzung
  eigenständig gegen `ci.yml` ausgeführt: `docker run --rm -v "$PWD":/repo
  -w /repo rhysd/actionlint:latest -color .github/workflows/ci.yml`)
  meldet dazu keinen Befund.
- geprüft, ohne Befund: **Docker-only (`AGENTS.md` §3.1).** Jeder `run:`-
  Schritt in `ci.yml` ruft ausschließlich ein bestehendes `make`-Target
  auf (`make gates`, `make mod-download`, `make test`); kein `docker
  run`/`docker build` direkt im Workflow, keine Inline-Shell-Logik, die
  eines dieser Ziele umgeht oder dupliziert. Alle drei Targets existieren
  im `Makefile`.
- geprüft, ohne Befund: **`fetch-depth: 0` ist für `commit-traceability`
  (`ADR-0045`, Default-Range `HEAD~5..HEAD`) tatsächlich notwendig** — ein
  flacher Checkout (Default-Tiefe 1) ließe `tools/harness/commit-traceability.sh`
  fail-closed mit Exit 2 abbrechen (Skript verifiziert `HEAD~5` explizit
  vor dem eigentlichen Lauf).
- geprüft, ohne Befund: **`make gates` ist vollständig netzlos.**
  `baseline-verify`, `docs-check`/`d-check`, `a-check` und
  `commit-traceability` laufen ohne Netzzugriff (`--network none` bzw.
  reine Git-/Dateisystem-Operationen) — `make mod-download` ist der
  einzige netzbedürftige Schritt und steht im Workflow korrekt separat
  vor `make test`, nicht innerhalb von `make gates`.
- geprüft, ohne Befund: **`dependabot.yml` deckt exakt `gomod` und
  `github-actions` ab, wöchentlich, `commit-message.prefix`
  `[ADR-0051]`, kein `docker`-Ecosystem** — wie in §1/§2 des Slice-Plans
  und `ADR-0051` Entscheidung 9 gefordert.
- geprüft, ohne Befund: **`[ADR-0051]`-Prefix erfüllt beide Hälften des
  Commit-Traceability-Gates.** Positive Hälfte: `.d-check.yml`
  `commits.id-patterns` enthält `ADR-\d{4}`, `ADR-0051` matcht. Grenz-
  Hälfte: `tools/harness/commit-traceability.sh` prüft nur auf
  `(SPEC|ARC)-[0-9]{3}` im Betreff — `[ADR-0051]` triggert das nicht.
- geprüft, ohne Befund: **`dependabot.yml` ist schema-valid gegen das
  offizielle SchemaStore-Schema** — in dieser Sitzung eigenständig
  reproduziert (Docker-Container, `jsonschema` gegen
  `dependabot-2.0.json`), nicht nur aus dem Plan-Nachzug übernommen:
  Ergebnis `schema-valid`.
- geprüft, ohne Befund: **`actionlint` gegen `ci.yml` real ausgeführt,
  keine Befunde** (siehe oben) — deckt sich mit der Behauptung im
  Plan-Nachzug, im Gegensatz zum `yamllint`-Befund (F-2).
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7), keine
  Slice-/Wellen-Chronik in den neuen YAML-Dateien.** `grep -n
  "slice-\|welle-\|Slice\|Welle" .github/workflows/ci.yml
  .github/dependabot.yml` liefert keinen Treffer; alle Kommentare
  beschreiben den Ist-Zustand (Zusage/Kopplung/Abgrenzung), keinen
  verworfenen Alternativtext.
- geprüft, ohne Befund: **§3.8-Nummerierung konsistent** — `AGENTS.md`
  führt §3.1–§3.8 lückenlos, §3.8 direkt nach §3.7 angehängt, Beispielform
  (`Falsch`/`Richtig`) mit dem tatsächlich verwendeten Pin identisch.
- geprüft, ohne Befund: **`harness/README.md`-Diff korrekt verlinkt und
  am richtigen Ort.** `../.github/workflows/ci.yml` und
  `../docs/plan/adr/0051-…md` lösen relativ zu `harness/README.md`
  korrekt auf; die Änderung sitzt an der bereits vor diesem Slice
  bestehenden "Aktueller Lauf-Status"-Zeile (keine neue Sensors-Tabellen-
  Zeile für ein Nicht-Gate, wie im Plan-Nachzug begründet) — keine
  Zwei-Quellen-Drift zur Sensors-Tabelle.
- geprüft, ohne Befund: **Kein Accepted-ADR verändert (`AGENTS.md`
  §3.5), kein Gate gelockert (§3.6).** `ADR-0051` bleibt `Accepted` und
  unverändert; `make gates`-Inhalt ist unangetastet, der Workflow
  automatisiert nur seinen Aufruf.
- geprüft, ohne Befund: **§1 Out-of-Scope, Klasse-Zuordnung korrekt.**
  Alle drei Ausschlusspunkte tragen eine Begründung und fallen in die
  Klassen „Folge-Slice" (Punkt 1 zur Sicht "Umfang: Folge-Slice" der ADR
  selbst, Punkt 2 mit `slice-040` als Adresse) bzw. „anderer Vorgang"
  (Punkt 3, realer GitHub-Lauf ist Nutzer-Handlung nach Merge).
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung.** Einzige berührte
  Sub-Area `*`/`PGC`, korrekt GF; Beobachtungs-Register-Sichtung
  dokumentiert und nachvollziehbar leer für CI/CD-spezifische Treffer
  (eigenständig gegen den vollständigen `BEO-PGC/*`-Bestand geprüft — 20
  Einträge, keiner CI/CD-spezifisch).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Plan-Nachzug-Verifikationsnachweis
nicht reproduzierbar" (1×, erstes Auftreten) · „Action-Pin real gegen
Upstream verifiziert" (1×, positiver Beleg, keine Fehler-Klasse) ·
„Folge-Slice-Adresse noch nicht angelegt" (1×, Hinweis ohne erwartete
Aktion) — kein Steering-Loop-Eintrag fällig, keine Klasse erreicht 3×.

## Verdikt

**Merge-blockierend:** nein im Sinne eines Rollback — der Commit ist
bereits gepusht, und das einzige MEDIUM-Finding (F-2) trägt **keinen
Rollen-Widerspruch**: Es ist eine Beobachtung gegen den eigenen
Plan-Nachzug-Text, keine bestrittene Implementer-Aussage. Die
Architect-Sequenz aus Modul 8 greift hier nicht (kein HIGH, kein
Widerspruch, keine dritte Wiederholung derselben Klasse). Erwartete
Reaktion: Der Implementer korrigiert den Plan-Nachzug-Text (entweder die
`yamllint`-Warnungen real beheben und den Nachzug entsprechend
aktualisieren, oder den Nachzug ehrlich auf den tatsächlichen
Werkzeug-Output umschreiben) vor der Closure — die inhaltlichen
YAML-Dateien selbst sind laut `actionlint` und dem eigenständig
reproduzierten `jsonschema`-Lauf funktional unbeeinträchtigt.

**Zur zentralen Prüffrage dieses Laufs (Action-Pinning, `AGENTS.md`
§3.8):** Bestätigt — der Checkout-SHA ist real, kein erfundener Pin
(F-1). **Zur Docker-only-/Gate-Automatisierungs-Frage (`ADR-0051`
Entscheidung 2):** Bestätigt — der Workflow ruft ausschließlich
bestehende, netzlose `make`-Targets auf, kein Gate wird umgangen oder
dupliziert. **Zur Dependabot-Traceability-Frage (`ADR-0045`,
Entscheidung 9):** Bestätigt — `[ADR-0051]` erfüllt beide Gate-Hälften,
real gegen die Regex-Muster geprüft.

**Übergabe:** Ein MEDIUM-Finding (Plan-Nachzug-Honesty), zwei INFO-
Hinweise. Der Implementer entscheidet über Annahme oder Begründung des
MEDIUM-Findings; keine Architect-Sequenz erforderlich. Dieser Report ist
ein Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen — die
Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5). Er ersetzt
keine Verifikation — DoD-/Spec-Konformität (inkl. der übrigen offenen
DoD-Punkte: Review-Häkchen selbst, Closure-Notiz, Register, Risiko-
Ausgänge, drei Paarungen) prüft der Verifier separat.
