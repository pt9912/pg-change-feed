# Verifikationsbericht: slice-054 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-054` §1 Ziel/Abgrenzung, §2 DoD, §4 Trigger, §6 Risiken, §8 Sub-Area)
und die bindenden `ADR-0055`/`ADR-0056` — nicht gegen Diff (Reviewer-Aufgabe,
bereits abgeschlossen: `docs/reviews/review-slice-054.md`, vollständig
gelesen) und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst —
kein MVP-Meilenstein-Slice, reiner Testablauf-Beleg ohne Produktionscode).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan, beide
ADRs vollständig, den Review-Report vollständig, den tatsächlichen Diff
(`git show`/`git diff` gegen `1cb1f06`) sowie den geänderten Code-Abschnitt
selbst (`tools/harness/run-integration-tests.sh`) — keine Implementer- oder
Reviewer-Behauptung wird ungeprüft übernommen. `make gates` und `make
test-integration` wurden in dieser Sitzung **eigenständig real ausgeführt**,
jeweils mit separat geprüftem Exit-Code (Ausgabe in eine Datei umgeleitet,
kein `| tail`/`| grep`, das den Exit-Code maskieren könnte). Zusätzlich habe
ich den `/subsz`-Mechanismus **unabhängig vom Reviewer** an einem eigenen,
frisch gestarteten NATS-Server mit identischem Digest-Pin real reproduziert
(§3 unten).

**Gegenstand:**
`docs/plan/planning/in-progress/slice-054-nats-boundary-beleg.md` zum Stand
`HEAD = 725814a`. Commits (chronologisch): `ec8b3d3` (`next`→`in-progress`,
reiner `git mv`, 0 Einfügungen/Löschungen real bestätigt), `1cb1f06`
(Boundary-Beleg-Implementierung, 55 Zeilen, ausschließlich Insertionen),
`725814a` (Review-Report, 0 HIGH/MEDIUM/LOW/1 INFO, DoD-Zeile „Review
durchgeführt" im selben Commit nachgezogen). Sequenz selbst geprüft:
Implementierung → Review, kein Self-Review (getrennter Kontext laut
Report-Kopf des Review-Reports), keine Rolle springt rückwärts ohne
Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `LH-FA-SST-007` Boundary real erfüllt: ohne abonnierten Test-Client eine Change erzeugen, danach real gegen `cdc.changes` belegen | **erfüllt, selbst reproduziert** | Eigener `make test-integration`-Lauf (voller Compose-Durchlauf, Exit 0, separat geprüft). Reale Log-Zeile: `run-integration-tests: NATS-Boundary-Beleg (LH-FA-SST-007) — Change (id=231, feed_mvp_full) entstand real ohne einen auf cdc.changes.src-mvp.public.feed_mvp_full abonnierten Client (belegt über NATS-Server-Monitor /subsz vor und nach der Change), blieb vollständig über cdc.changes lesbar, und der Feed-Container lief unverändert weiter`. Code real gelesen (`run-integration-tests.sh:1310-1325`): Vor-Check gegen `/subsz`, `INSERT id=231`, SQL-Lesebeleg gegen `cdc.changes` mit `count(*) = 1`-Prüfung, sonst `exit 1`. |
| 2 | Derselbe Lauf belegt real, dass `CaptureService.Capture()` nicht blockiert/fehlschlägt | **erfüllt, selbst reproduziert** | Im selben Lauf, unmittelbar nach dem Boundary-Insert: `feed_running=$(docker inspect --format '{{.State.Running}}' "$FEED_CONTAINER" …)`, bei `!= "true"` bricht das Skript mit `exit 1` ab (Zeile 1329-1332, real gelesen). Mein Lauf zeigt Exit 0 über den gesamten restlichen Testablauf hinweg (u. a. der nachfolgende `TestMVPSchemaChangeIncompatibleTypeChange`-Testfall lief noch danach durch) — der Feed-Container hat den Boundary-Fall also nicht nur überlebt, sondern produktiv weitergearbeitet. |
| 3 | `make gates` grün, `make test-integration` grün | **erfüllt, selbst reproduziert** | Beide Läufe eigenständig ausgeführt, Exit-Code jeweils **separat** geprüft (`make … > datei 2>&1; echo "EXIT=$?"`, keine Pipe): `make gates` → `EXIT=0` (`baseline-verify: v6.5.0 OK`, `coverage-gate: OK — Coverage 40.90% erfüllt Schwelle 35%`, `d-check: 428 Datei(en) geprüft, 0 Befund(e)` sowohl im vollen Lauf als auch in der Commit-Range, `commit-traceability: OK`, `a-check: gesamt: 0 Befund(e)` mit dem vorbestehenden, unveränderten Hinweis zu `tools/harness/natssub/main.go`). `make test-integration` → `EXIT=0`, voller Compose-Durchlauf inkl. Happy-Path- und Boundary-Beleg, endet mit `PASS`/`ok`. |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | `docs/reviews/review-slice-054.md` vollständig gelesen: 0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO (F-1, ohne erwartete Aktion). DoD-Zeile im selben Commit (`725814a`) korrekt nachgezogen. |
| 5 | Kein Doku-Update nötig | **erfüllt** | `git diff --stat 1cb1f06~1 1cb1f06` real geprüft: ausschließlich `tools/harness/run-integration-tests.sh` (55 Insertionen, 0 Deletionen) geändert — kein öffentlicher Vertrag (`spec/`, `docs/user/`) berührt. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt noch ausschließlich Platzhalter (`<…>`) — Planner-Closure-Arbeit, die laut Rollen-Sequenz (Modul 8) erst nach diesem Bericht beginnt. |
| 7 | Reconciliation-Register | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` real geprüft: existiert nicht — Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/` real aufgelistet — kein `slice-054`-Beleg in irgendeinem `evidence/`-Verzeichnis. Konsistent mit „offen, Planner-Closure-Arbeit". |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide Zeilen tragen noch wörtlich `<bei Closure einzutragen>`, real per `grep` bestätigt. Eigene Einschätzung als Vorschlag an den Planner unten (§4). |
| 10 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind erst nach dem `git mv` sinnvoll prüfbar. |

## 2. Plan-vs-Code-Diff

- **Datei-Liste (§3) exakt getroffen:** einzige geänderte Datei ist
  `tools/harness/run-integration-tests.sh` (`git show --stat 1cb1f06`:
  `1 file changed, 55 insertions(+)`), genau wie in §3 des Slice-Plans
  vorgesehen.
- **Kein Out-of-Scope-Punkt berührt (§1):** eigener `git diff 1cb1f06~1
  1cb1f06` real durchsucht — keine Änderung an
  `internal/application/port/outbound/changenotification.go`, keine
  Änderung an `internal/adapters/driven/natsnotify/`, kein neuer Eintrag in
  `compose.yaml`, kein Reconnect-Verhalten (der Testablauf prüft
  ausschließlich den Fall *nie verbunden gewesen*, keine
  Trennung-und-Wiederverbindung eines zuvor verbundenen Clients).
- **Neue Umgebungsvariable `NATS_CONTAINER`** (Default `cdc-test-nats`) ist
  eine reine Testskript-Variable nach demselben Muster wie
  `PG_CONTAINER`/`FEED_CONTAINER` — keine neue Compose-Wiring, sondern nur
  eine Referenz auf den bereits in `compose.yaml` (aus `slice-053`)
  bestehenden Service-Namen.
- **Kommentar-Disziplin (`AGENTS.md` §3.7):** eigener Blick auf den
  17-zeiligen Kommentar über der Boundary-Sektion — beschreibt den
  gegenwärtigen Mechanismus im Indikativ (Zusage: reale
  Server-Momentaufnahme; Abgrenzung: `$SYS.*` vs. `cdc.changes.*`), keine
  der drei HIGH-Kommentar-Klassen (verworfene Alternative im Konjunktiv,
  abwesenter Text, abgebrochener Satz) trifft zu.
- **Keine unbegründete Abweichung gefunden.**

## 3. `/subsz`-Mechanismus — unabhängige eigene Reproduktion

Der Reviewer hat den Mechanismus bereits real gegengeprüft; ich habe das
**ohne Bezug auf seine Ausgabe**, mit einem eigenen, frisch gestarteten
Container desselben Digest-Pins (`nats:2-alpine@sha256:065e8355c2…`)
wiederholt:

1. Leerlauf-Zustand: `/subsz?subs=1` gegen einen frisch gestarteten
   NATS-Server zeigt ausschließlich `$SYS.*`-Subjekte (60 Treffer, alle
   `$SYS`-Präfix, keiner mit `cdc.changes` überlappend).
2. Aktiver Subscriber: ein realer `nats sub`-Client auf exakt
   `cdc.changes.src-mvp.public.feed_mvp_full` (`natsio/nats-box`-Image) —
   `/subsz?subs=1` zeigt das Subjekt danach real im Output
   (`"subject": "cdc.changes.src-mvp.public.feed_mvp_full"`, Exit 0 des
   `grep`).
3. Nach `docker rm -f` des Subscriber-Containers: `/subsz?subs=1` zeigt das
   Subjekt **nicht mehr** (0 Treffer).

Damit ist unabhängig — nicht nur vom Reviewer, sondern von mir selbst am
laufenden Server — bestätigt: Der `grep -qF "\"subject\": \"$NATS_SUBJECT\""`
-Mechanismus in `run-integration-tests.sh` prüft real gegen den NATS-Server
selbst, nicht gegen eine Annahme über das Testablauf-Verhalten. Er kann eine
tatsächlich noch offene Subscription real erkennen und meldet sie
(begonnen ab Schritt 2 oben).

## 4. §6-Risiken — eigenes, unabhängiges Urteil (kein Ausgang eingetragen — Planner-Arbeit)

- **Risiko 1 — Rest-Subscriber aus vorherigem Testschritt könnte den
  Boundary-Fall unbemerkt verfälschen.** Real adressiert: Der Testablauf
  entfernt den Happy-Path-Subscriber-Container (`docker rm -f
  "$NATS_SUBSCRIBER_CONTAINER"`) **vor** dem Boundary-Abschnitt (Code real
  gelesen, Reihenfolge sequenziell im selben Bash-Skript), und prüft danach
  real gegen `/subsz` — zweimal, vor **und** nach der auslösenden Change.
  Mein eigener Lauf (§1 Punkt 1) zeigt keinen Abbruch an dieser Stelle, und
  meine unabhängige Reproduktion des Mechanismus (§3) bestätigt: Wäre ein
  Rest-Subscriber tatsächlich noch verbunden gewesen, hätte der Vor-Check
  ihn real gefunden und das Skript mit `exit 1` beendet. Meine Einschätzung:
  **trägt für den Ausgang „entfallen"** — die Prüfung ist strukturell
  selbst-detektierend; ein künftiger Rest-Subscriber würde den Testlauf
  selbst zum Scheitern bringen, nicht unbemerkt bleiben.
- **Risiko 2 — `Capture()`s Notify-Aufruf könnte bei fehlendem Subscriber
  einen propagierten Fehler erzeugen (Regressionsrisiko trotz
  `slice-052`s Unit-Test).** Real end-to-end widerlegt: Mein eigener Lauf
  zeigt den Feed-Container nach der Boundary-Change weiter laufend
  (`feed_running=true`, sonst hätte das Skript abgebrochen), und der
  nachfolgende Testfall (`TestMVPSchemaChangeIncompatibleTypeChange`) lief
  im selben Container-Lebenszyklus noch danach erfolgreich durch — kein
  stiller Absturz, keine Blockade. Meine Einschätzung: **trägt ebenfalls für
  den Ausgang „entfallen"** — die end-to-end-Prüfung bestätigt real, dass
  `ADR-0055` Punkt 4s Fehlerschluckung am tatsächlichen NATS-Client greift,
  nicht nur in der isolierten Unit-Test-Annahme aus `slice-052`.

Beide Einschätzungen sind Vorschläge zur eigenständigen Prüfung durch den
Planner, keine gesetzten Werte.

## 5. Hard Rules

- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `ec8b3d3` real per
  `git show --stat` bestätigt: 0 Einfügungen/Löschungen, reiner Rename.
- **3.7 (Kommentar-Disziplin):** siehe §2 oben — kein Verstoß gefunden.
- **3.1 (Docker-only):** Boundary-Beleg läuft ausschließlich über `docker
  exec`/`docker inspect` gegen bereits laufende, digest-gepinnte Container;
  keine lokale Toolchain-Installation.
- **3.8 (Action-Pinning):** nicht einschlägig — kein `.github/workflows/`-Diff.

## 6. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 10) — Slice liegt noch in `in-progress/`.
Beobachtungs-Register-Eintrag und Closure-Notiz (Planner-Closure-Arbeit,
beginnt erst nach diesem Bericht). Validierung gegen realen Bedarf (kein
Validator-Zug ausgelöst — reiner Testablauf-Beleg ohne neuen
Architektur-Sicht-Meilenstein).

## Verdikt

**DoD-Konformität: bestätigt** für alle fünf technisch/funktional prüfbaren
Punkte (1–5), jeweils selbst reproduziert (`make gates`, `make
test-integration`, direkte Code-Lektüre, unabhängige `/subsz`-Reproduktion).
Die verbleibenden fünf Punkte (Closure-Notiz, Beobachtungs-Register,
§6-Risiko-Ausgänge, Paarungen) sind **korrekt noch offen** — Planner-Closure-
Arbeit, die laut Rollen-Sequenz (Modul 8) erst nach diesem Bericht beginnt.

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Genau eine Datei
geändert, exakt wie in §3 des Slice-Plans vorgesehen; kein Out-of-Scope-Punkt
(§1) berührt — keine Änderung an `ChangeNotificationPort`/`natsnotify`,
keine neue Compose-Verdrahtung, kein Reconnect-Verhalten.

**`/subsz`-Mechanismus: unabhängig bestätigt.** Eigene Reproduktion an einem
frisch gestarteten Server desselben Digest-Pins (nicht nur Lektüre der
Reviewer-Aussage) zeigt: Subjekt erscheint real bei aktivem Subscriber und
verschwindet real unmittelbar nach dessen Entfernung. Mechanismus-basiert,
keine Annahme.

**§6-Risiken:** Für beide Risiken spricht meine eigene Prüfung für den
Ausgang „entfallen" — beide sind durch reale, selbst-detektierende Prüfungen
im Testablauf strukturell abgedeckt, kein verbleibendes Restrisiko über
einen künftigen Testablauf-Refactor hinaus. Vorschlag zur eigenständigen
Prüfung durch den Planner, kein gesetzter Wert.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: die beiden §6-Risiko-Ausgänge
(Vorschlag oben, beide „entfallen"), der Beobachtungs-Registereintrag (laut
§8 des Slice-Plans keine Treffer zu NATS/Notification/Messaging vor diesem
Slice — die INFO-Finding-Klasse des Reviews, „fixe Wartezeit als einziges
Fenster für Negativ-Beleg", ist eine der drei Closure-Quellen und gehört in
die Register-Prüfung), und der `git mv` nach `done/` mit den drei Paarungen
danach. Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
