# Verifikationsbericht: slice-055 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-055` §1 Ziel/Abgrenzung, §2 DoD, §4 Trigger, §6 Risiken, §8 Sub-Area)
und die bindenden `ADR-0055`/`ADR-0056` — nicht gegen Diff (Reviewer-Aufgabe,
bereits abgeschlossen: `docs/reviews/review-slice-055.md`, vollständig
gelesen) und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst —
kein MVP-Meilenstein-Slice, reiner Testablauf-Beleg ohne Produktionscode).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan, beide
ADRs vollständig, den Review-Report vollständig, den tatsächlichen Diff
(`git show`) sowie den geänderten Code-Abschnitt selbst
(`tools/harness/run-integration-tests.sh:1333-1490`) — keine Implementer-
oder Reviewer-Behauptung wird ungeprüft übernommen. `make gates` und `make
test-integration` wurden in dieser Sitzung **eigenständig real ausgeführt**,
jeweils mit separat geprüftem Exit-Code (Ausgabe in eine Datei umgeleitet,
`echo "ECHTER_EXIT=$?"` an das Log angehängt, kein `| tail`/`| grep`, das den
Exit-Code maskieren könnte). Zusätzlich habe ich den
`docker network disconnect`/`docker inspect`-Mechanismus **unabhängig vom
Reviewer** an einem eigenen, frisch erzeugten Wegwerf-Netz mit zwei
`alpine:3.20`-Containern real reproduziert (§3 unten), inklusive eines
gleichzeitigen `ping`-Belegs, den der Reviewer nicht geführt hat.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-055-nats-negative-reconnect-beleg.md`
zum Stand `HEAD = a550520`. Commits (chronologisch): `8c52c6b`
(`next`→`in-progress`, reiner `git mv`), `2987665` (Negative-Beleg-
Implementierung, 154 Zeilen, ausschließlich Insertionen), `1a75087`
(DoD-Nachzug für die in diesem Lauf real erfüllten Punkte), `a550520`
(Review-Report, 0 HIGH/MEDIUM/LOW/2 INFO, DoD-Zeile „Review durchgeführt" im
selben Commit nachgezogen). Sequenz selbst geprüft: Implementierung →
DoD-Nachzug → Review, kein Self-Review (getrennter Kontext laut Report-Kopf
des Review-Reports), keine Rolle springt rückwärts ohne Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `LH-FA-SST-007` Negative-Beleg real erfüllt: real trennen, während der Trennung Changes erzeugen, danach real wiederverbinden, Changes ausschließlich über `cdc.changes` sichtbar | **erfüllt, selbst reproduziert** | Eigener `make test-integration`-Lauf (voller Compose-Durchlauf, `ECHTER_EXIT=0`, separat geprüft). Reale Log-Zeile: `run-integration-tests: NATS-Negative-Beleg (LH-FA-SST-007, Reconnect-Nachholen) — Test-Subscriber real vom Compose-Netz getrennt (belegt über docker inspect), verpasste Change (id=240) blieb ohne jedes Wecksignal (Log-Beleg) und wurde ausschließlich über cdc.changes nachgeholt; ein frischer Wiederverbindungs-Subscriber empfing für eine neue Change (id=241) real ein Signal …`. Code real gelesen (Zeilen 1363–1429): `docker network disconnect`, SQL-Beleg gegen `cdc.changes` mit `count(*) = 1`-Prüfung für `id=240`, sonst `exit 1`. |
| 2 | Testablauf belegt real (Log-Beleg), dass während der Trennung **kein** Wecksignal ankommt | **erfüllt, selbst reproduziert** | Zeile 1417–1422 real gelesen: `nats_reconnect_before_output=$(docker logs …)`; `if printf '%s' "$nats_reconnect_before_output" | grep -qF "RECEIVED"; then exit 1; fi` — explizite Negativ-Prüfung, kein impliziter Schluss. Mein Lauf zeigt keinen Abbruch an dieser Stelle. |
| 3 | `make gates` grün, `make test-integration` grün | **erfüllt, selbst reproduziert** | Beide Läufe eigenständig ausgeführt, Exit-Code jeweils **separat** geprüft (`make … > datei 2>&1; echo "ECHTER_EXIT=$?"` ans Log angehängt, keine Pipe). `make gates` → `ECHTER_EXIT=0` (`coverage-gate: OK — Coverage 40.90% erfüllt Schwelle 35%`, `d-check: 430 Datei(en) geprüft, 0 Befund(e)` sowohl im vollen Lauf als auch in der Commit-Range `HEAD~5..HEAD`, `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID`, `a-check: gesamt: 0 Befund(e)` mit dem vorbestehenden, unveränderten Hinweis zu `tools/harness/natssub/main.go`). `make test-integration` → `ECHTER_EXIT=0`, voller Compose-Durchlauf inkl. Happy-Path-, Boundary- und Negativ-Beleg, endet mit `PASS`/`ok` (`TestMVPSchemaChangeIncompatibleTypeChange` lief danach noch erfolgreich durch). |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | `docs/reviews/review-slice-055.md` vollständig gelesen: 0 HIGH, 0 MEDIUM, 0 LOW, 2 INFO (F-1, F-2, beide ohne erwartete Aktion). DoD-Zeile im selben Commit (`a550520`) korrekt nachgezogen. |
| 5 | Kein Doku-Update nötig | **erfüllt** | `git show --stat 2987665` real geprüft: ausschließlich `tools/harness/run-integration-tests.sh` (154 Insertionen, 0 Deletionen) geändert — kein öffentlicher Vertrag (`spec/`, `docs/user/`) berührt. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt noch ausschließlich Platzhalter (`<…>`) — Planner-Closure-Arbeit, die laut Rollen-Sequenz (Modul 8) erst nach diesem Bericht beginnt. |
| 7 | Reconciliation-Register | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` real geprüft: existiert nicht — Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`); DoD-Zeile trägt die Begründung bereits explizit. |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/` real aufgelistet — kein `slice-055`-Beleg in irgendeinem `evidence/`-Verzeichnis, keine der beiden Review-Finding-Klassen (Mechanismus-Realität `docker inspect`, Negativbeleg stützt sich auf externe architektonische Garantie) bereits registriert. Konsistent mit „offen, Planner-Closure-Arbeit". |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide Zeilen tragen noch wörtlich `<bei Closure einzutragen>`, real per Lektüre bestätigt. Eigene Einschätzung als Vorschlag an den Planner unten (§4). |
| 10 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind erst nach dem `git mv` sinnvoll prüfbar. |

## 2. Plan-vs-Code-Diff

- **Datei-Liste (§3) exakt getroffen:** einzige geänderte Datei ist
  `tools/harness/run-integration-tests.sh` (`git show --stat 2987665`:
  `1 file changed, 154 insertions(+)`), genau wie in §3 des Slice-Plans
  vorgesehen.
- **Kein Out-of-Scope-Punkt berührt (§1):** eigener `git show 2987665`
  vollständig durchgesehen — keine Änderung an
  `internal/application/port/outbound/changenotification.go`, keine
  Änderung an `internal/adapters/driven/natsnotify/`, kein neuer Eintrag in
  `compose.yaml`. Kein Reconnect-Verhalten im Produktionscode: Der Testablauf
  ruft an keiner Stelle in `nats.go`s eingebaute Reconnect-Logik ein,
  sondern trennt und ersetzt den **Test-Subscriber-Container** komplett
  (neuer Prozess statt Reconnect derselben Instanz — Design-Entscheidung, die
  der Reviewer bereits geprüft und für plan-konsistent befunden hat; ich habe
  die Begründung im Kommentar Zeilen 1334–1362 selbst nachvollzogen und
  teile das Verdikt: ein neuer Prozess kann strukturell nichts empfangen, was
  vor seiner eigenen Subscription publiziert wurde — das macht den
  Negativ-Beleg eindeutig, ohne die produktive `nats.go`-Reconnect-Logik
  überhaupt zu berühren).
- **Kein „nie verbunden"-Fall:** Der Testablauf lässt den Subscriber vor der
  Trennung real abonnieren und über `/subsz?subs=1` bestätigen (Zeile
  1391–1396) — das ist der *verbunden-gewesen*-Fall aus §1, nicht der
  *nie-verbunden*-Fall aus `slice-054`. Klare Abgrenzung eingehalten.
- **Kommentar-Disziplin (`AGENTS.md` §3.7):** eigener Blick auf den
  26-zeiligen Kommentar (Zeilen 1334–1362) — beschreibt den gegenwärtigen
  Mechanismus und die Designentscheidung im Indikativ bzw. mit
  begründendem Konjunktiv über ein **hypothetisches Risiko einer verworfenen
  Design-Alternative**, nicht über einen abwesenden/früheren Code-Zustand;
  keine der drei HIGH-Kommentar-Klassen trifft zu.
- **Keine unbegründete Abweichung gefunden.**

## 3. `docker network disconnect`/`docker inspect`-Mechanismus — unabhängige eigene Reproduktion

Der Reviewer hat den Mechanismus bereits real gegengeprüft (F-1 des Review-
Reports); ich habe das **ohne Bezug auf seine Ausgabe**, mit einem eigenen,
frisch erzeugten Wegwerf-Netz und zwei `alpine:3.20`-Containern wiederholt:

1. Zwei Container (`vnc_a`, `vnc_b`) auf einem eigenen Docker-Netz
   (`vnc_net`). `docker inspect vnc_a` zeigt vor der Trennung das Netz mit
   IP-Adresse; ein `ping` von `vnc_b` auf `vnc_a`s IP zeigt `0% packet loss`.
2. `docker network disconnect vnc_net vnc_a` ausgeführt.
3. Unmittelbar danach: `docker inspect vnc_a --format
   '{{json .NetworkSettings.Networks}}'` liefert `{}` — das Netz ist sofort
   und vollständig aus der Auskunft verschwunden.
4. Im **selben Moment** ein erneuter `ping` von `vnc_b` auf dieselbe IP:
   `100% packet loss`, sofort.

Damit ist unabhängig — nicht nur vom Reviewer, sondern von mir selbst an
einem eigenen Wegwerf-Aufbau — bestätigt: `docker inspect`s Wegfall des
Netzwerk-Eintrags fällt exakt mit einem realen, sofortigen
Verbindungsabbruch auf Paketebene zusammen. Es ist kein bloßer
Absichts-Beleg, sondern ein Zustand, der sich mit dem realen Netzwerkverhalten
deckt — dieselbe Schlussfolgerung wie F-1 des Reviews, jetzt zusätzlich mit
einem eigenen `ping`-Beleg gestützt, den der Review-Report nicht führte.

## 4. §6-Risiken — eigenes, unabhängiges Urteil (kein Ausgang eingetragen — Planner-Arbeit)

- **Risiko 1 — Reale NATS-Verbindungstrennung im Compose-Netz könnte
  technisch aufwändiger sein als ein reiner Prozess-Stopp.** Real
  widerlegt: `docker network disconnect` gegen den Test-Subscriber-Container
  funktionierte sowohl im Implementer-Lauf als auch in meiner unabhängigen
  Reproduktion (§3) beim ersten Versuch, ohne Sonderbehandlung oder
  zusätzliche Untersuchung. Der Mechanismus ist robust und deterministisch
  (sofortiger Wegfall im `docker inspect`, kein Warten auf serverseitige
  Ping-Erkennung nötig). Meine Einschätzung: **trägt für den Ausgang
  „entfallen"** — die befürchtete zusätzliche Komplexität ist nicht
  eingetreten, der gewählte Ansatz (Trennung des Test-Containers statt des
  NATS-Servers) hat direkt funktioniert.
- **Risiko 2 — Test könnte fälschlich einen SQL-Nachhol-Beleg zeigen, obwohl
  in Wahrheit doch ein NATS-Signal nachgeliefert wurde.** Real adressiert:
  Der Testablauf prüft explizit und getrennt vom SQL-Beleg, dass **kein**
  `RECEIVED`-Frame im Log des getrennten Subscribers erscheint (Zeile
  1417–1422, `grep -qF "RECEIVED"` → `exit 1` bei Treffer), **bevor** der SQL-
  Nachhol-Beleg überhaupt geprüft wird (Zeile 1424–1429). Beide Prüfungen
  sind damit unabhängig voneinander geführt, nicht nur eine Schlussfolgerung
  aus der anderen. Zusätzlich eine reale Wartezeit (`sleep 5`, Zeile 1415)
  vor der Log-Prüfung, die einem verspätet doch zugestellten Frame Zeit
  ließe, aufzutauchen — trat in meinem Lauf nicht ein. Meine Einschätzung:
  **trägt ebenfalls für den Ausgang „entfallen"** — die Log-Prüfung ist
  strukturell unabhängig vom SQL-Beleg und selbst-detektierend; ein
  künftiges Nachliefern würde den Testlauf selbst zum Scheitern bringen,
  nicht unbemerkt bleiben. (Die verbleibende Einschränkung — inhaltliche
  Unterscheidung `id=240` vs. `id=241` beruht auf externen architektonischen
  Garantien statt auf Payload-Inhalt — ist bereits als Review-Finding F-2
  erfasst und eine separate Beobachtungs-Register-Kandidatin, kein
  unaufgelöster Teil dieses Risikos.)

Beide Einschätzungen sind Vorschläge zur eigenständigen Prüfung durch den
Planner, keine gesetzten Werte.

## 5. Hard Rules

- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `8c52c6b` real per
  `git show --stat` bestätigt: reiner `next→in-progress`-Move ohne
  Inhaltsänderung, laut Review-Report-Kopf.
- **3.7 (Kommentar-Disziplin):** siehe §2 oben — kein Verstoß gefunden.
- **3.1 (Docker-only):** Negativ-Beleg läuft ausschließlich über `docker
  run`/`docker exec`/`docker inspect`/`docker network disconnect`/`docker
  logs`/`docker rm` gegen bereits laufende, digest-gepinnte Container; keine
  lokale Toolchain-Installation.
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
test-integration`, direkte Code-Lektüre, unabhängige `docker network
disconnect`/`docker inspect`-Reproduktion samt `ping`-Beleg). Die
verbleibenden fünf Punkte (Closure-Notiz, Beobachtungs-Register,
§6-Risiko-Ausgänge, Paarungen) sind **korrekt noch offen** — Planner-Closure-
Arbeit, die laut Rollen-Sequenz (Modul 8) erst nach diesem Bericht beginnt.

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Genau eine Datei
geändert, exakt wie in §3 des Slice-Plans vorgesehen; kein Out-of-Scope-Punkt
(§1) berührt — keine Änderung an `ChangeNotificationPort`/`natsnotify`, keine
neue Compose-Verdrahtung, kein Reconnect-Verhalten im Produktionscode, kein
„nie verbunden"-Fall.

**`docker network disconnect`/`docker inspect`-Mechanismus: unabhängig
bestätigt.** Eigene Reproduktion an einem frischen Wegwerf-Netz (nicht nur
Lektüre der Reviewer-Aussage) zeigt: Der Netzwerk-Eintrag verschwindet real
sofort aus `docker inspect`, zeitgleich mit einem realen `ping`-Ausfall
(100 % Paketverlust). Mechanismus-basiert, keine Annahme.

**§6-Risiken:** Für beide Risiken spricht meine eigene Prüfung für den
Ausgang „entfallen" — beide sind durch reale, strukturell
selbst-detektierende Prüfungen im Testablauf abgedeckt, kein verbleibendes
Restrisiko über einen künftigen Testablauf-Refactor hinaus. Vorschlag zur
eigenständigen Prüfung durch den Planner, kein gesetzter Wert.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: die beiden §6-Risiko-Ausgänge
(Vorschlag oben, beide „entfallen"), der Beobachtungs-Register-Eintrag (die
beiden Review-Finding-Klassen — Mechanismus-Realität des
`docker-inspect`-Trennungsbelegs, Negativbeleg stützt sich auf externe
architektonische Garantie — sind die dritte Closure-Quelle, `neuer Sensor`/
`geschärfte Regel` je nach Planner-Urteil), und der `git mv` nach `done/` mit
den drei Paarungen danach. Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
