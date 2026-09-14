# Verifikationsbericht: slice-063 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-063` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger inkl. Nachtrag,
§6 Risiken, §8 Sub-Area) und die bindende
[`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
(Supersedes `ADR-0058`, nur Entscheidung 3) — nicht gegen Diff
(Reviewer-Aufgabe, bereits abgeschlossen ohne Fixrunde:
`docs/reviews/review-slice-063.md`, vollständig gelesen, aber als Kontext,
nicht als Ersatz für eigene Prüfung übernommen) und nicht gegen realen Bedarf
(Validator — hier nicht ausgelöst, `slice-063` ist kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan, die
vollständige `ADR-0064` (inkl. `ADR-0058` Entscheidung 3 als Kontext), den
vollständigen Review-Report, den vollständigen Blocker-Report
(`docs/reviews/blocker-slice-063.md`), den tatsächlichen Diff seit `265eda8`
(reiner `next→in-progress`-Move) bis `HEAD`, den vollständigen neuen
Skript-Abschnitt in `tools/harness/run-integration-tests.sh`
(Zeilen ~1554–1649), den `harness/README.md`-Diff, `spec/lastenheft.md`
(`LH-QA-OPS-005`-Wortlaut selbst, nicht nur die ADR-Paraphrase), den
ADR-Index (`docs/plan/adr/README.md`), das Beobachtungs-Register
(`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`) und
`docs/plan/planning/welle-17.md`. `make gates` und `make test-integration`
wurden in dieser Sitzung **eigenständig real ausgeführt** — kein
Implementer- oder Reviewer-Beleg ungeprüft übernommen.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-063-upgrade-sicherheit-schema-rollout-zyklus.md`
zum Stand `HEAD = a73efb2`. Zwei Commits seit `265eda8`:

- `94d7530` — Implementierung: neue Upgrade-Sicherheits-Phase in
  `tools/harness/run-integration-tests.sh` (nach `ADR-0064`s korrigiertem
  Container-Tausch-Mechanismus), `harness/README.md`-Update,
  DoD-Häkchen für Implementierungs-/Sensor-Punkte
- `a73efb2` — Review-Report (0 HIGH/MEDIUM/LOW), zieht die
  DoD-Review-Zeile nach

Diesem Diff ging ein realer Blocker-Durchlauf voraus (Rückführung
`in-progress → open`, `docs/reviews/blocker-slice-063.md`) und ein
unabhängiger Architect-Zug (`ADR-0064`), der den Mechanismus korrigierte und
`slice-063` zurück nach `in-progress` (Kopf-Feld-Nachzug auf `ADR-0064`)
brachte — beide Vorgänge liegen vor `265eda8` bzw. sind dessen Ursache und
damit außerhalb des hier geprüften Implementierungs-Diffs, aber vollständig
gelesen als Kontext.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `LH-QA-OPS-005` erfüllt: neue Phase, `ADR-0064`-Mechanismus, Datenstand identisch lesbar | **erfüllt, selbst reproduziert** | Code gelesen (Zeilen ~1554–1649): `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` ersetzt den `ADR-0058`-Dreischritt vollständig, kein Rest davon im Diff. Eigener `make test-integration`-Lauf: Log-Zeile `„Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005, ADR-0064) belegt"` mit Positionswerten vor/nach. `spec/lastenheft.md` selbst geprüft (Zeile 1150–1154): Messmethode verlangt `„Datenstand vor/nach Upgrade identisch lesbar (LH-QA-REL-001)"` — kein Migrationsschritt im Wortlaut; `ADR-0064`s Kernbehauptung damit unabhängig bestätigt. |
| 2 | Real bestätigt: `--force-recreate` erzeugt neue Instanz, `postgres`/`nats` unberührt via `--no-deps` | **erfüllt, selbst reproduziert** | Eigener Lauf: Feed-Container-ID wechselte real (`5d637c9e…` → `504f4867…`), Skript prüft das hart (`upgrade_feed_id_after = upgrade_feed_id_before` → `exit 1`) und vergleicht `postgres`/`nats`-IDs ebenso hart vor/nach — im eigenen Lauf keine Abweichung, kein `exit 1` ausgelöst. Log zeigt `Container cdc-test-feed Recreated`, keine entsprechende Zeile für `postgres`/`nats`. |
| 3 | Health-Poll nach dem Tausch analog `LH-QA-REL-001` | **erfüllt, selbst reproduziert** | Code gelesen: identische Schleifenstruktur (60×1s, `docker inspect --format '{{.State.Health.Status}}'`) wie der bestehende `docker restart`-Rundlauf. Eigener Lauf durchlief den Poll ohne Timeout (kein `exit 1` an dieser Stelle, Folgezeilen liefen). |
| 4 | `make gates` grün | **erfüllt, selbst reproduziert** | Eigener, vollständiger, ungefiltert ausgeführter Lauf, Exit-Code unmittelbar danach in eigenem Schritt geprüft: `0` (§2 unten). |
| 5 | `make test-integration` grün mit neuer Phase sichtbar im Log | **erfüllt, selbst reproduziert** | Eigener, vollständiger, ungefiltert ausgeführter Lauf, Exit-Code unmittelbar danach in eigenem Schritt geprüft: `0`; Log-Zeile der neuen Phase real vorhanden (§2 unten). |
| 6 | Review durchgeführt, Report liegt vor | **erfüllt** | `docs/reviews/review-slice-063.md` vollständig gelesen: 0 HIGH/0 MEDIUM/0 LOW/0 INFO, keine Fixrunde. Zentrale Prüfpunkte (Mechanismus, Platzierung, Container-ID-Wechsel, Datenstand) real gegengeprüft (§3/§4 unten), nicht nur Reviewer-Aussage übernommen. |
| 7 | Doku-Update `harness/README.md` §Sensors | **erfüllt, selbst reproduziert** | `git diff 265eda8..HEAD -- harness/README.md` real gelesen: bestehende `make test-integration`-Zeile trägt neuen Satz zur Upgrade-Sicherheits-Phase, korrekt gegen `LH-QA-OPS-005`/`ADR-0064` referenziert (`Supersedes ADR-0058 Entscheidung 3`), `· seit slice-063` im etablierten Muster. Kein neues Gate, kein neues Target — Tabellenstruktur unverändert. |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich Platzhalter (`<…>`), real per Volltext-Lektüre bestätigt — Planner-Arbeit nach diesem Bericht. |
| 9 | Reconciliation-Register — entfällt | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` real geprüft: existiert nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/evidence/` real gelistet: enthält ausschließlich `slice-016.md` — kein `slice-063`-Beleg angelegt, obwohl §4-Nachtrag und §8 des Slice-Plans den realen Zweitfund bereits beschreiben. `state.md` zeigt weiterhin Zähler „1× — unter der 3×-Schwelle". Konsistent mit dem offenen DoD-Punkt: die formale Registereintragung ist Closure-Arbeit, nicht Implementierungs-Arbeit. |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Vier §6-Einträge real gezählt: nur der erste (`BEO-PGC/schema-rollout-fremdobjekte`) trägt einen ausformulierten Ausgang (`eingetreten → ADR-0064`); die drei übrigen (Container-Stopp-Seiteneffekt, Timing-Flake-Interferenz, kein echter Versionswechsel) tragen noch wörtlich `<bei Closure zu füllen>` — real per Volltext-Lektüre bestätigt. |
| 12 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/` (real per Verzeichnis-Listung bestätigt); die Paarungen suchen in `done/` und sind vor dem `git mv` nicht sinnvoll prüfbar — Slice trägt `Welle: welle-17`, DoD verweist korrekt auf die Welle-17-Closure. |

**Ergebnis §1:** Alle sieben implementierungs-/reviewbezogenen DoD-Punkte
(1–7) sind real erfüllt und selbst reproduziert, nicht nur behauptet. Die
fünf Closure-Punkte (8–12) sind korrekt noch offen und wurden **nicht** vom
Implementer oder Reviewer vorweggenommen — die DoD-Checkbox-Trennung (7
abgehakt: 6 Implementierung/Sensor/Doku + Review, 5 offen: sämtliche
Closure-Pflichten) ist sauber.

## 2. Sensor-Läufe (selbst ausgeführt)

**`make gates`** — vollständiger, ungefiltert ausgeführter Lauf, Ausgabe in
eigene Log-Datei umgeleitet, Exit-Code unmittelbar danach in einem eigenen,
nicht-gepipten Schritt geprüft (`AGENTS.md` §3.9):

```
coverage-gate: OK — Coverage 44.40% erfüllt Schwelle 35%
d-check: 505 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 505 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
  Hinweis: tools/harness/httpclient/main.go, tools/harness/natssub/main.go
  liegen in keiner Schicht (unverändert bekannt, keine neue Fundstelle)
```

Exit-Code: **0**.

**`make test-integration`** (voller Compose-Lauf, gegen den echten
Feed-Container) — ebenso ungefiltert, Exit-Code separat geprüft:

```
Container cdc-test-feed Recreate
Container cdc-test-feed Recreated
Container cdc-test-feed Starting
Container cdc-test-feed Started
run-integration-tests: Upgrade-Sicherheits-Rundlauf (LH-QA-OPS-005, ADR-0064)
belegt — realer Container-Tausch über $COMPOSE up -d --force-recreate
--no-deps (5d637c9e22d4… -> 504f486771c4…), postgres/nats unberührt
(--no-deps), Feed-Container danach healthy, Datenstand vor dem Tausch
(id=250, Position 30803840) identisch lesbar, danach eingefügte Zeile
(id=251) weiterhin erfasst (Position 30805120)
=== RUN   TestE2ESchemaChangeDropColumn
--- PASS: TestE2ESchemaChangeDropColumn (0.26s)
--- PASS: TestE2ESchemaChangeIncompatibleTypeChange (0.26s)
```

Exit-Code: **0**. Container-ID-Wechsel real beobachtet (Präfix identisch mit
dem Log-Wert), kein `postgres`/`nats`-Recreate-Log — beide unberührt.
`TestE2ESchemaChangeDropColumn` PASSt direkt im Anschluss, ohne dass der
vorangegangene Container-Tausch dessen eigenen Recovery-Mechanismus störte.

## 3. Platzierungsentscheidung — eigenständig geprüft, nicht nur Review-Urteil übernommen

`TestE2ESchemaChangeDropColumn` (Zeile 919 in
`test/integration/integration_test.go`) beendet den Erfassungspfad des Feed-
Containers real dauerhaft: Das umgebende Skript prüft danach hart
`feed_running != "false"` → `exit 1` — der Test **verlangt selbst**, dass
der Container nach seinem Lauf nicht mehr läuft, und behebt das über einen
eigenen, spezifischen Recovery-Block (Replication-Slot-Neuanlage,
Schema-Version-Nachtrag, expliziter Neustart, eigener Health-Poll,
Zeilen ~1670–1700+).

Läge die Upgrade-Sicherheits-Phase **nach** diesem Test statt davor, träfe
ihr eigener `--force-recreate`-Aufruf entweder auf einen Container, dessen
Recovery-Zustand (neu angelegter Slot, nachgetragene Schema-Version) durch
den Container-Tausch redundant überschrieben oder mit ihm verwechselt
würde, oder die Upgrade-Phase liefe gegen einen Container, den der
DropColumn-Test bewusst als „nicht mehr laufend" hinterlassen will — beides
kollidierte mit einer der beiden Testabsichten. Die gewählte Platzierung
**vor** `TestE2ESchemaChangeDropColumn` vermeidet das strukturell: Die
Upgrade-Phase läuft gegen den „normalen", noch unversehrten Feed-Container,
bevor dessen Recovery-Mechanismus überhaupt ins Spiel kommt.

**Eigenes Urteil:** Die Platzierungsentscheidung trägt, unabhängig
nachvollzogen anhand des realen Recovery-Codes und des realen Testlaufs
(`TestE2ESchemaChangeDropColumn` PASSt unverändert direkt nach der neuen
Phase) — nicht nur aus der Review-Begründung übernommen.

## 4. `ADR-0064` — eigenständig gegen `spec/lastenheft.md` geprüft

`ADR-0064`s Kernargument (§Kontext, „Was verlangt `LH-QA-OPS-005`
tatsächlich?") behauptet, der Lastenheft-Wortlaut verlange keinen
Migrationsschritt. Real gegengeprüft gegen `spec/lastenheft.md` Zeilen
1150–1154: Der Wortlaut lautet exakt „Upgrade-Test in der Testumgebung:
Datenstand vor/nach Upgrade identisch lesbar (`LH-QA-REL-001`)" — kein
Migrationsschritt-Bezug, keine Schema-Wandel-Anforderung. `ADR-0064`s
Behauptung trägt, unabhängig bestätigt gegen die Rang-1-Quelle
(Lastenheft), nicht nur gegen die ADR-eigene Paraphrase.

`ADR-0064` selbst wurde als `Accepted`-Datei in genau einem Commit angelegt
(`c19657a`, per `git log --follow` bestätigt) — keine nachträgliche
inhaltliche Überschreibung. `ADR-0058` selbst ebenfalls unverändert
(einziger Commit `30d07d8`) — die Korrektur lief korrekt über den
Folge-ADR-Weg (`Supersedes`, nur Entscheidung 3), Hard Rule `AGENTS.md`
§3.5 eingehalten. Der ADR-Index (`docs/plan/adr/README.md`) trägt die
Supersedes-Beziehung bidirektional (`ADR-0058`-Zeile: `„→ ADR-0063/0064”`;
`ADR-0064`-Zeile: `„Supersedes ADR-0058, teilweise”`).

## 5. §2 DoD-Häkchen — Trennung Implementierung/Review vs. Closure

Real per Volltext-Lektüre bestätigt: Die sieben abgehakten Punkte decken
ausschließlich Implementierung (Punkte 1–3), Sensor-Läufe (4–5), Review (6)
und Doku-Update (7) — keiner der fünf Closure-Punkte (Closure-Notiz,
Reconciliation-Entfall, Beobachtungs-Register, Risiko-Ausgänge, drei
Paarungen) ist vorzeitig abgehakt. Der Review-Report-Commit (`a73efb2`) zog
laut eigenem Diff ausschließlich die Review-Zeile nach (`[ ]`→`[x]` plus
zwei angehängte Zeilen mit dem Report-Zeiger) — keine andere DoD-Zeile
wurde in diesem Commit berührt.

## 6. §4-Nachtrag und Bezug-Feld — real konsistent

Der Kopf trägt `ADR-0064` als maßgeblich, `ADR-0058` als „nur noch Kontext"
— konsistent mit dem tatsächlichen Supersedes-Umfang (nur Entscheidung 3).
§4 „Nachtrag — was tatsächlich eintrat" beschreibt die reale Rückführung
`in-progress → open` korrekt im Indikativ über den eingetretenen Zustand,
mit Zeiger auf `docs/reviews/blocker-slice-063.md` — real gelesen, deckt
sich mit dem dortigen Bericht (Exit 8, vier statt zwei Fremdobjekte, drei
verworfene Umgehungen).

## 7. Welle-17-Erfüllbarkeit — nur zur Einordnung, kein Closure-Urteil dieser Rolle

`docs/plan/planning/welle-17.md` §3 verlangt `slice-062`, `slice-063`,
`slice-064`, `slice-065` **alle vier** in `done/`, zusätzlich einen real
belegten grünen `e2e.yml`-Matrix-Lauf über beide PostgreSQL-Legs nach einem
echten Push auf den Hauptzweig, sowie `make gates` grün und eine
Closure-Notiz.

Real geprüft (`ls docs/plan/planning/done/`, `ls
docs/plan/planning/{in-progress,next,open}/`):

- `slice-062` liegt bereits in `done/`.
- `slice-063` liegt noch in `in-progress/` — dieser Bericht bestätigt seine
  DoD-Konformität für den aktuellen Implementierungs-/Review-Stand, macht
  ihn aber nicht selbst zu `done/` (Planner-Arbeit).
- `slice-064` und `slice-065` liegen **beide noch in `open/`** — nicht
  einmal in `next/` oder `in-progress/`.

**Eigenes Urteil:** Selbst wenn `slice-063` nach diesem Bericht sofort
Closure durchläuft und nach `done/` wandert, wird `welle-17`s
Closure-Trigger dadurch **nicht** erfüllbar — zwei der vier Pflicht-Slices
(`slice-064`, `slice-065`) haben ihre Implementierung noch nicht einmal
begonnen (`open/`, kein `Verantwortlich:`-Zug, kein WIP-Claim). Der
Closure-Trigger fordert zusätzlich einen realen grünen `e2e.yml`-Matrix-Lauf
über beide PostgreSQL-Legs — dieser Lauf hängt laut `welle-17.md` §5 selbst
von `slice-064` ab, das die Matrix erst parametrisiert. `slice-063`s
Closure ist eine notwendige, aber bei weitem keine hinreichende Bedingung
für die Welle-17-Closure.

## 8. Hard Rules

- **3.3 (`git mv` + Inhaltsänderung = zwei Commits):** `265eda8` real als
  Elter bestätigt (reiner `next→in-progress`-Move, nicht Bestandteil des
  geprüften Bereichs).
- **3.5 (Accepted-ADR-Immutabilität):** `ADR-0058` unverändert (§4 oben),
  `ADR-0064` selbst in genau einem Commit angelegt, keine nachträgliche
  Überschreibung.
- **3.7 (Kommentar-/Chronik-Disziplin):** Der neue Phasen-Kommentar
  (Zeilen 1554–1565) beschreibt den geltenden Mechanismus samt Begründung
  im Indikativ, referenziert `BEO-PGC/schema-rollout-fremdobjekte` und
  `docs/reviews/blocker-slice-063.md` als Beleg-Anker, nicht als Chronik
  ohne Träger — folgt demselben, bereits etablierten Muster wie die
  bestehende `BEO-PGC/walsender-wirksamkeit`-Referenz im selben Skript.
- **3.9 (Exit-Code nie gepiped):** in dieser Sitzung durchgehend beachtet
  (§2) — `make gates`/`make test-integration` jeweils in eine Log-Datei
  umgeleitet (kein Pipe zwischen Lauf und Exit-Code-Prüfung), Exit-Code
  unmittelbar danach in einem eigenen, ungeketteten Bash-Aufruf geprüft.

## 9. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 12) — Slice liegt noch in `in-progress/`.
Closure-Notiz, Beobachtungs-Register-Neueintrag und §6-Risiko-Ausgänge
(Planner-Closure-Arbeit, beginnt laut Rollen-Sequenz Modul 8 erst nach
diesem Bericht). Die Welle-17-Closure selbst (§7 oben ordnet nur ein, urteilt
nicht abschließend — hängt an zwei noch nicht begonnenen Slices).
Validierung gegen realen Bedarf: **kein Validator-Zug ausgelöst** —
`slice-063` ist kein MVP-Meilenstein-Slice.

Dieses Repo führt keine `make doc-commits`-/`make doc-immutable`-Targets
(`AGENTS.md` §4 listet nur real existierende Targets). Die Traceability-
Prüfung je Commit läuft hier über `make commit-traceability`
(Bestandteil von `make gates`, §2 oben bestätigt: 5 Commits im Range,
0 Befunde); ADR-Immutabilität ist eine Hard Rule (`AGENTS.md` §3.5, §4 oben
real gegengeprüft) ohne eigenes Sensor-Target.

## Verdikt

**DoD-Konformität: bestätigt** für alle sieben implementierungs-/
reviewbezogenen Punkte (1–7), jeweils selbst reproduziert (`make gates`,
`make test-integration`, eigenständige Prüfung von Container-ID-Wechsel,
Datenstand-Beleg, Platzierungsentscheidung und `ADR-0064`s
Lastenheft-Argument gegen `spec/lastenheft.md` selbst). Die fünf
verbleibenden Closure-Punkte (8–12) sind korrekt noch offen und wurden
nicht vorweggenommen.

**`ADR-0064`-Konformität: bestätigt.** Der implementierte Mechanismus
(`$COMPOSE up -d --force-recreate --no-deps pg-change-feed`) entspricht
`ADR-0064`s Entscheidung wortgetreu; die ADR selbst ist unverändert
`Accepted`, `ADR-0058` unverändert, Supersedes-Beziehung bidirektional im
Index verankert. `ADR-0064`s Kernbehauptung (kein Migrationsschritt im
Lastenheft-Wortlaut gefordert) wurde eigenständig gegen `spec/lastenheft.md`
gegengeprüft und bestätigt.

**Platzierungsentscheidung: eigenständig bestätigt**, nicht nur aus dem
Review-Urteil übernommen (§3 oben).

**Sensor-Läufe:** `make gates` — Exit-Code **0** (in Log-Datei umgeleitet,
danach in eigenem Schritt geprüft). `make test-integration` — Exit-Code
**0** (ebenso), Container-Tausch real beobachtet
(`5d637c9e22d4…` → `504f486771c4…`), `postgres`/`nats` unberührt,
Datenstand-Beleg im Log vollständig vorhanden.

**Welle-17: nicht erfüllbar durch diesen Slice allein.** `slice-064` und
`slice-065` liegen noch in `open/`, nicht einmal begonnen — die
Welle-Closure bleibt fern, unabhängig vom Closure-Zeitpunkt von `slice-063`.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: der erste §6-Risiko-Eintrag
(`BEO-PGC/schema-rollout-fremdobjekte`) trägt bereits einen Ausgang
(„eingetreten → `ADR-0064`"), die drei übrigen sind noch offen zu füllen
(Container-Stopp-Seiteneffekt trat im real ausgeführten Lauf nicht ein;
Timing-Flake-Interferenz trat im real ausgeführten Lauf nicht ein; „kein
echter Versionswechsel" bleibt strukturell weiter offen — Ausgang „weiter
offen" naheliegend, aber Planner-Urteil). Der Beobachtungs-Registereintrag
für `BEO-PGC/schema-rollout-fremdobjekte` braucht bei Closure eine zweite
`evidence/`-Datei (`slice-063`), aktuell trägt das Verzeichnis nur
`evidence/slice-016.md`. Welle-17 bleibt nach dieser Closure weiterhin
offen — zwei Slices (`slice-064`, `slice-065`) sind noch nicht begonnen.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
