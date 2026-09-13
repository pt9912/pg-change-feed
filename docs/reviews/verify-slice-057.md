# Verifikationsbericht: slice-057 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-057` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan-Nachzug, §4 Trigger, §6
Risiken) — nicht gegen Diff (Reviewer-Aufgabe, bereits abgeschlossen:
`docs/reviews/review-slice-057.md`, 0 Findings, vollständig gelesen) und
nicht gegen realen Bedarf (Validator, hier nicht ausgelöst — reines
Namensrelikt-Renaming ohne neuen Architektur-Sicht-Meilenstein).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan
(`docs/plan/planning/in-progress/slice-057-e2e-umgebung-umbenennung.md`)
und den Review-Report vollständig, führt den Diff `bbf5828..6184908`
eigenständig aus (`git diff`, `git diff --stat`, `git log`) und **fährt alle
Sensoren selbst** — keine Implementer- oder Reviewer-Behauptung wird
ungeprüft übernommen. Arbeitsverzeichnis zu Beginn und Ende dieser Sitzung
sauber (`git status`).

**Gegenstand:** `bbf5828..6184908` (sechs Commits: `feaea4f`, `13f968f`,
`e806b19`, `88a9839`, `c1d24be`, `6184908`), plus der Folge-Commit
`8f5ded6` (Review-Report-Nachzug, kein weiterer Umbenennungs-Inhalt). Die
Slice-Datei liegt weiterhin in `in-progress/` — §7 (Closure-Notiz) und die
letzten fünf DoD-Punkte sind laut Modul 8 §Rollen-Sequenz **Planner-
Closure-Arbeit**, die erst nach diesem Bericht beginnt; das ist kein
DoD-Mangel.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Alle in §1 genannten Bezeichner real umbenannt, `grep -rn "mvp"` (case-insensitive, außer den zwei §1-Ausnahmen) auf null Treffer | **erfüllt, selbst reproduziert** | `grep -rniE "mvp" compose.yaml test/integration/integration_test.go tools/harness/run-integration-tests.sh docs/user/benutzerhandbuch.md` selbst ausgeführt: genau drei Treffer — die zwei in §1 deklarierten Meilenstein-Ausnahmen (`integration_test.go:5` Godoc-Verweis auf den fachlichen „MVP-Schnitt", `run-integration-tests.sh:4` „ursprünglichen MVP-Zuschnitt") plus die neue Changelog-Zeile 1.12 im Benutzerhandbuch (historisch-erklärend, kein technischer Bezeichner). Kein weiterer `mvp`-Fund. |
| 2 | `make test-integration` läuft real vollständig durch, inhaltlich unverändert erfolgreich | **erfüllt, nach eigener Mehrfach-Reproduktion — mit einer festgehaltenen Flake-Beobachtung** | Drei eigene, vollständige Läufe, Exit-Code jeweils separat und ungefiltert geprüft (`AGENTS.md` §3.9): Lauf 1 Exit `2` (Abbruch im Retention-Lebenszyklus-Abschnitt, siehe unten), Lauf 2 Exit `0` (alle zehn `TestE2E*`-Subtests `PASS`, alle Bash-Abschnitte inkl. Retention-Lebenszyklus, CLI-Diagnose, SQL-Administration Live-Reload, Publication-Entzug, NATS Happy/Boundary/Negative grün), Lauf 3 Exit `0` (identisch vollständig). Siehe §3 unten für die Einordnung des einen roten Laufs. |
| 3 | `make gates`, `make test` grün | **erfüllt, selbst reproduziert** | Beide eigenständig ausgeführt, Exit-Code je separat geprüft: `make gates` Exit `0` (`baseline-verify` OK, `coverage-gate: OK — Coverage 41.00% erfüllt Schwelle 35%`, `d-check: 439 Datei(en) geprüft, 0 Befund(e)` sowohl im vollen als auch im `commits`-Modul, `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID`, `a-check: gesamt: 0 Befund(e)`). `make test` Exit `0`, alle Pakete `ok`, kein Fehlschlag. |
| 4 | Review durchgeführt, Report unter `docs/reviews/review-slice-057.md` liegt vor | **erfüllt** | Report vollständig gelesen: 0 HIGH/MEDIUM/LOW/INFO, neun Negativbefunde, Verdikt „nicht merge-blockierend". DoD-Zeile im selben Commit-Bereich (`8f5ded6`) korrekt nachgezogen, kein Self-Review (getrennte Kontexte laut Report-Kopf). |
| 5 | Doku-Update: `docs/user/benutzerhandbuch.md`-Beispiele, falls alte Namen | **erfüllt** | `git diff bbf5828..6184908 -- docs/user/benutzerhandbuch.md` real gelesen: Beispielwert `src-mvp`→`src-e2e` in §4 aktualisiert, `Version:`-Kopf 1.11→1.12, `Stand:` nachgezogen, neue Zeile in `### Änderungshistorie` im etablierten Muster (Version, Datum, Kurzbeschreibung inkl. Slice-Bezug) — konsistent mit `BEO-PGC/handbuch-versionshistorie-uebersprungen`. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **korrekt offen** | §7 trägt noch ausschließlich Platzhalter (`<…>`) — Planner-Closure-Arbeit nach Modul 8 §Rollen-Sequenz, beginnt erst nach diesem Bericht. Kein DoD-Mangel. |
| 7 | Reconciliation-Register — falls Inventur-Fund | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht — Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Kein `slice-057`-Beleg in irgendeinem `evidence/`-Verzeichnis unter `docs/plan/planning/observations/BEO-PGC/`; kein neues Verzeichnis. Konsistent mit „offen, Planner-Closure-Arbeit". Für die Closure vorzumerken: die in §3 unten festgehaltene Flake-Beobachtung ist bisher in keinem Registereintrag erfasst. |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen — eigenes Urteil siehe §3 unten** | Beide Zeilen in §6 tragen noch wörtlich `<bei Closure einzutragen>` — real per Lektüre bestätigt, kein Ausgang gesetzt. |
| 10 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind erst nach dem `git mv` sinnvoll prüfbar. |

**Zwischenbefund:** Alle fünf technisch/funktional real prüfbaren Punkte
(1–5) sind **selbst reproduziert erfüllt**. Die verbleibenden fünf Punkte
(6, 8, 9, 10, plus der bereits erledigte Reconciliation-Entfall) sind
**korrekt noch offen** — Planner-Closure-Arbeit, die erst nach diesem
Bericht beginnt. Kein DoD-Verstoß.

## 2. Sensor-Läufe (selbst ausgeführt)

**`grep -rniE "mvp" compose.yaml test/integration/integration_test.go tools/harness/run-integration-tests.sh docs/user/benutzerhandbuch.md`**
— 3 Treffer, alle erwartet (siehe DoD-Punkt 1 oben).

**`make gates`** (vollständiger Lauf, Exit `0`):

```
coverage-gate: OK — Coverage 41.00% erfüllt Schwelle 35%
d-check: 439 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 439 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
  Hinweis: tools/harness/natssub/main.go liegt in keiner Schicht (dokumentiertes,
  nicht-fatales Verhalten, unverändert seit vor diesem Slice)
```

**`make test`** (alle Pakete, Exit `0`) — kein Fehlschlag.

**`make test-integration`** — drei eigene, vollständige Läufe:

- Lauf 1: Exit `2`. Abbruch im „Retention-Lebenszyklus-Rundlauf"
  (`tools/harness/run-integration-tests.sh` Zeile 470f.): Die Zeile
  `RetentionLifecycle` (id=210) war nach 25s Wartezeit bereits entfernt,
  obwohl der Consumer seine Position noch nicht bestätigt hatte
  (`LH-FA-RET-004`-Assertion). `make: *** [Makefile:69: test-integration]
  Fehler 1`.
- Lauf 2: Exit `0`. Alle zehn `TestE2E*`-Subtests `PASS`. Alle
  Bash-Abschnitte grün, inklusive **desselben** Retention-Lebenszyklus-
  Abschnitts, der in Lauf 1 fehlschlug — real bestätigt durch die Log-Zeile
  `run-integration-tests: Retention-Beleg — 'RetentionOld' (id=200) blieb
  erhalten, solange ein Consumer zurückhing (LH-FA-RET-004) …`. Ebenso
  grün: Black-Box-CLI-Rundlauf, Verarbeitungsrückstand-Beleg, CLI-Diagnose
  (Normalbetrieb + Fehlerzustand), SQL-Administration Live-Reload
  (enable/disable), Publication-Entzug-Wirksamkeit, NATS
  Happy-/Boundary-/Negative-Beleg — alle unter den neuen `e2e`-Namen
  (`src-e2e`, `feed_e2e_full`, `feed_e2e_sql_admin`,
  `feed_e2e_walsender_timing`, `cli-e2e-*`-Consumer-Namen).
- Lauf 3: Exit `0`, identisch vollständig durchgelaufen.

## 3. Einordnung des roten Laufs — eigenes Urteil

Der eine rote Lauf betrifft die zeitkritische Assertion „Zeile bleibt über
mehr als zwei Lösch-Takte (`retentionInterval=10s`) erhalten, solange ein
Consumer nicht bestätigt hat" nach einer festen `sleep 25`. Drei
unabhängige Belege stützen die Einordnung als **vorbestehende,
umgebungsbedingte Zeitflake, nicht als Regression dieses Slice**:

1. **Der Diff ist mechanisch nachweisbar reines Renaming.** `git diff
   bbf5828..6184908 -- test/integration/integration_test.go
   tools/harness/run-integration-tests.sh compose.yaml | grep -E '^[+-]' |
   grep -v '^+++\|^---' | grep -viE 'mvp|e2e'` liefert **keine** Zeile —
   jede geänderte Zeile in den drei technischen Dateien nennt `mvp` oder
   `e2e`. Kein Zeitwert (`sleep 25`, `retentionInterval`), keine Assertion
   und kein Kontrollfluss wurde berührt.
2. **Zwei von drei eigenen Läufen liefen die identische Assertion grün
   durch**, im selben Codestand, in derselben Sitzung — konsistent mit
   einer last-/timing-abhängigen Flake statt einem deterministischen
   Defekt.
3. **Kein Treffer im Beobachtungs-Register** unter einer bereits bekannten
   Zeitflake-Klasse für diesen Abschnitt (`docs/plan/planning/
   observations/BEO-PGC/` durchsucht: `retention-keine-loeschausfuehrung`
   betrifft die frühere, vollständig fehlende Löschausführung, nicht
   Timing; kein Treffer für „RET-004", „Flake", „Timing", „Race" in diesem
   Kontext).

**Meine Einschätzung, als Hinweis an den Planner:** Diese Flake ist nicht
ursächlich mit `slice-057` verbunden und blockiert die DoD nicht (Punkt 2
ist durch die zwei grünen Vollläufe erfüllt) — sie sollte aber bei der
Closure als **neue** Beobachtung
(`BEO-PGC/retention-lebenszyklus-timing-flake` o. ä.) ins Register
aufgenommen werden, damit ein wiederholtes Auftreten gezählt wird statt
jedes Mal neu entdeckt zu werden. Das ist ein Vorschlag, kein DoD-Punkt
dieses Slices — die Registerpflege bleibt Planner-Closure-Arbeit.

## 4. §6-Risiken — eigenes, unabhängiges Urteil (kein Ausgang eingetragen — Planner-Arbeit)

- **Risiko 1 — übersehenes `mvp`-Literal in einem Heredoc.** Eigene,
  vollständige Grep-Prüfung über alle vier Zieldateien (§2 oben): keinerlei
  unerwarteter `mvp`-Treffer, auch nicht innerhalb der SQL-Heredocs in
  `run-integration-tests.sh` (eigene Lektüre der Heredoc-Blöcke bestätigt
  konsistente `e2e`-Umbenennung, z. B. `pgc_e2e_reader_login`,
  `pgc_e2e_capture_login`). Meine Einschätzung: **entfallen**, mit
  Begründung — real geprüft, kein Fund.
- **Risiko 2 — `slice-054`/`055` tragen noch alte `mvp`-Namen.** Eigene
  Prüfung: Beide liegen bereits in `done/`
  (`docs/plan/planning/done/slice-054-nats-boundary-beleg.md`,
  `docs/plan/planning/done/slice-055-nats-negative-reconnect-beleg.md`),
  geschlossen **vor** Beginn dieses Slice (`git log` zeigt
  `slice-055 in-progress -> done` chronologisch vor `slice-057
  Verantwortlich gesetzt`). `grep -liE "mvp"` gegen beide Dateien liefert
  keinen Treffer. Meine Einschätzung: **entfallen**, mit Begründung — der
  befürchtete Drift-Zwischenzustand ist nie eingetreten, beide Pläne waren
  bereits vor Aktivierung dieses Slice geschlossen und tragen keine
  `mvp`-Reste.

## 5. Plan-vs-Code-Diff

- **Datei-Liste (§3) exakt getroffen:** `git diff bbf5828..6184908
  --name-only` zeigt genau `compose.yaml`,
  `test/integration/integration_test.go`,
  `tools/harness/run-integration-tests.sh`, `docs/user/benutzerhandbuch.md`
  plus die Slice-Plan-Datei selbst (DoD-Nachzug/Dublette-Entfernung, kein
  Liefer-Punkt). Kein zusätzlicher, unbegründeter Dateizugriff.
- **Kein Out-of-Scope-Punkt (§1) berührt:** `git diff bbf5828..6184908
  --stat -- docs/plan/planning/done/ docs/reviews/` liefert **keine**
  Ausgabe — archivierte/geschlossene Dokumente sind real unberührt. `git
  diff bbf5828..6184908 --stat -- spec/lastenheft.md
  docs/plan/planning/in-progress/roadmap.md` liefert ebenfalls **keine**
  Ausgabe — der Meilenstein-Begriff „MVP" ist unberührt. Compose-
  Servicenamen `cdc-test-postgres`/`cdc-test-feed` real geprüft
  (`grep -n` in `compose.yaml`): unverändert, tragen kein `mvp`-Fragment.
- **Plan-Nachzug (§3) korrekt als solcher deklariert:** Die vier
  zusätzlich umbenannten Bezeichner-Gruppen (`tbl-mvp-*`/`sv-mvp-*` in
  `CDC_TABLES`, `TestMVP*`-Funktionsnamen/`mvpEnv`, `feed_mvp_sql_admin`/
  `feed_mvp_walsender_timing`/`-run`-Regex/`'MVP-Quelle'`,
  `Version:`-Kopf/Changelog) sind in §3 explizit als *Plan-Nachzug*
  benannt und real in genau diesem Umfang umgesetzt — keine unbenannte
  Abweichung.
- **Keine unbegründete Abweichung gefunden.**

## 6. Verhaltensgleichheit — eigene, unabhängige Prüfung

- `git diff bbf5828..6184908 -- test/integration/integration_test.go
  tools/harness/run-integration-tests.sh | grep -E '^[+-]' | grep -v
  '^+++\|^---' | grep -viE 'mvp|e2e'` liefert **keine** Zeile — jede
  geänderte Zeile in den beiden Testträgern ist eine reine Umbenennung.
- `compose.yaml`-Diff vollständig gelesen: nur die vier
  `CDC_*`-Bezeichner-Werte geändert, keine sonstige Konfigurationszeile.
- Alle zehn `TestE2E*`-Subtests laufen unter identischem Namen, identischer
  Reihenfolge und identischer Aussage wie zuvor `TestMVP*` (Reviewer-Beleg
  bereits geprüft, eigener Lauf 2/3 bestätigt dieselbe Zahl und Reihenfolge
  `PASS`).
- Keine geänderte Assertion, kein geänderter Kontrollfluss, kein
  geänderter Zeitwert im gesamten Diff.

**Verhaltensgleichheit bestätigt.**

## 7. Hard Rules

- **3.9 (Exit-Code-Regel):** Alle Sensor-Läufe dieser Sitzung liefen
  ungefiltert in eine Log-Datei, der Exit-Code wurde unmittelbar danach
  und ungekettet gesichert (`make …; echo $? > …`), niemals durch eine
  Pipe hindurch geprüft.
- **3.7 (Kommentar-Disziplin/Slice-Chronik-Verbot):** `git diff
  bbf5828..6184908 -- '*.go' '*.sh' | grep "^+" | grep -inE
  "slice-[0-9]+|welle-[0-9]+"` liefert **keinen** Treffer in Code/Skript.
  Die einzige neue Slice-Kennung im Diff steht in der
  Änderungshistorie-Zeile des Benutzerhandbuchs — dort als deklarierter
  Chronik-Träger zulässig.
- **3.3 (git mv + Inhaltsänderung = zwei Commits):** Kein `git mv` im
  geprüften Bereich — nicht einschlägig.
- **3.1 (Docker-only):** Alle Sensor-Läufe liefen über `make`-Targets im
  Container.
- **3.6 (Gates nicht ohne ADR lockern):** nicht einschlägig — keine
  Schwellen-Änderung.

## 8. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 10) — Slice liegt noch in `in-progress/`.
Closure-Notiz und Beobachtungs-Registereintrag (Planner-Closure-Arbeit,
beginnt erst nach diesem Bericht). Validierung gegen realen Bedarf (kein
Validator-Zug ausgelöst — reines Namensrelikt-Renaming, kein neuer
Architektur-Sicht-Meilenstein).

## Verdikt

**DoD-Konformität: bestätigt.** Alle fünf technisch/funktional prüfbaren
Punkte (1–5) sind selbst reproduziert erfüllt, inklusive dreifach
wiederholtem `make test-integration` mit im Diff nachweisbar unbeteiligter
Ursache für den einen roten Lauf (siehe §3). Die verbleibenden Punkte (6,
8, 9, 10) sind korrekt noch offen — Planner-Closure-Arbeit nach Modul 8
§Rollen-Sequenz; der Reconciliation-Punkt entfällt korrekt (kein
Brownfield-Bootstrap in diesem Repo).

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Datei-Liste (§3)
exakt getroffen, kein Out-of-Scope-Punkt (§1) berührt — `done/`,
`docs/reviews/`, der Meilenstein-Begriff „MVP" und die Compose-
Servicenamen sind real unberührt. Der Plan-Nachzug (§3) ist vollständig
und ohne unbenannte Erweiterung umgesetzt.

**Verhaltensgleichheit: bestätigt.** Jede geänderte Zeile in den drei
technischen Trägern ist eine reine Umbenennung; keine Assertion, kein
Kontrollfluss, kein Zeitwert wurde geändert.

**Flake-Hinweis an den Planner:** Ein Lauf von `make test-integration`
brach im Retention-Lebenszyklus-Abschnitt (`LH-FA-RET-004`-Assertion) ab,
zwei weitere Läufe desselben Codestands liefen vollständig grün durch. Die
Ursache liegt nachweislich nicht im Umbenennungs-Diff (§3); zur Closure
empfohlen: neue Beobachtungs-Registerzeile für diese Zeitflake, damit ein
wiederholtes Auftreten gezählt wird.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: beide §6-Risiken als
„entfallen" mit Begründung (Vorschlag oben), der Flake-Hinweis als neuer
Beobachtungs-Registereintrag, und der `git mv` nach `done/` mit den drei
Paarungen danach. Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
