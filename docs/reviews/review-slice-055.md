# Review-Report: slice-055 — 2026-09-13

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `2987665` (Testcode) + `1a75087` (DoD-Nachzug) —
NATS-Negative-Beleg (Reconnect-Nachholen) in
`tools/harness/run-integration-tests.sh` + Slice-Plan-Checkboxen. Elter
`8c52c6b` (reiner `next→in-progress`-Move, keine Inhaltsänderung).

**Skill:** `.harness/skills/reviewer.md` @ `f07ba3d` (Stand 2026-09-13, vier
repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-055-nats-negative-reconnect-beleg.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §6 Risiken)
- `ADR-0055` (`docs/plan/adr/0055-nats-change-notification-wecksignal.md`) —
  Punkt 3 (leerer Payload), Punkt 4 (Fehlerklasse/kritischer Pfad,
  best-effort nach ACK, keine Retry-Logik für ein verpasstes Signal)
- `ADR-0056` (`docs/plan/adr/0056-nats-tabellen-granulares-subjekt.md`) —
  Vier-Token-Subjekt `cdc.changes.<source_id>.<schema>.<table>`, hier in
  `NATS_SUBJECT` (Bestand aus `slice-053`) korrekt verwendet
- `LH-FA-SST-007` (Negative-Akzeptanzkriterium)
- `tools/harness/natssub/main.go` (Wegwerf-Subscriber, genau ein
  `NextMsg` je Prozess)
- `compose.yaml` (NATS-Service, Digest-Pin, `-m 8222` Monitor-Port)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7), §6 Workflow
- `harness/conventions.md` (MR-000 ID-Schema)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-054.md`
  (akzeptierte Kommentar-Länge/-Form als Präzedenzfall)

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings. Zwei INFO-Beobachtungen unten.

### F-1 — `docker inspect`-Mechanismus real reproduziert und bestätigt belastbar

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:1401-1406`
- `befund`: Eigener Reproduktionsversuch außerhalb des Repos (Wegwerf-Container
  `alpine:3.20` auf einem eigenen Docker-Netz) bestätigt: `docker network
  disconnect` gefolgt von `docker inspect --format
  '{{json .NetworkSettings.Networks}}'` zeigt das Netz sofort als entfernt
  (`{}`), und ein Ping von einem zweiten Container auf dieselbe IP zeigt im
  selben Moment 100 % Paketverlust — `docker inspect` ist damit kein bloßer
  Absichts-Beleg, sondern deckt sich exakt mit dem realen, sofortigen
  Verbindungsabbruch auf Netzwerkebene. Der im Kommentar begründete Verzicht
  auf `/subsz` (NATS-eigene, Ping-basierte Erkennung im Minutenbereich) als
  Trennungs-Beleg ist damit nicht nur plausibel, sondern verifiziert.
- `verifizierbar`: ja — eigener `docker network disconnect`/`docker
  inspect`/`ping`-Testlauf gegen Wegwerf-Container, in diesem Review real
  ausgeführt.
- `klasse`: „Mechanismus-Realität des docker-inspect-Trennungsbelegs"

### F-2 — Notify-Signal-Unterscheidung zwischen verpasster und neuer Change ist architektonisch, nicht inhaltlich geprüft

- `kategorie`: INFO
- `quelle`: `ADR-0055` Punkt 3 (leerer Payload)
- `pfad`: `tools/harness/run-integration-tests.sh:1461-1468`
- `befund`: Der Test kann inhaltlich nicht unterscheiden, ob ein empfangenes
  Signal zu Change `id=241` oder `id=240` gehört — `ADR-0055` Punkt 3 legt den
  Payload bewusst leer an. Die Eindeutigkeit des Beweises stützt sich
  ausschließlich auf zwei architektonische Garantien, die dieser Slice nicht
  selbst prüft, sondern voraussetzt: Core NATS hält für ein Subjekt ohne
  aktiven Abonnenten keine Nachricht zurück (kein Nachliefern beim späteren
  Abonnieren), und ein einmal verpasstes Notify wird nicht erneut versucht
  (`ADR-0055` Punkt 4, „ein verpasstes einzelnes Signal hat … keine
  Konsequenz"). Beide Garantien sind bereits an anderer Stelle
  belegt/entschieden (ADR bzw. `natsnotify`-Whitebox-Test), nicht Gegenstand
  dieses Slice — kein Mangel, aber eine Kette, deren Glieder außerhalb dieses
  Diffs liegen.
- `verifizierbar`: ja — `make test-integration`, bereits real ausgeführt;
  die architektonischen Garantien selbst deckt der `natsnotify`-Whitebox-Test
  aus `slice-053`/`ADR-0055` Fitness Function ab.
- `klasse`: „Negativbeleg stützt sich auf externe architektonische Garantie
  statt auf inhaltliche Unterscheidung"

## Design-Entscheidung: frischer Subscriber-Prozess statt Reconnect desselben Containers

Geprüft, kein Finding. Der Implementer begründet im Code-Kommentar
(`tools/harness/run-integration-tests.sh:1333-1358`), warum die
Wiederverbindung über einen **neuen** Subscriber-Prozess statt über
`docker network connect` auf demselben Container hergestellt wird: Eine vom
Broker noch nicht als tot erkannte TCP-Verbindung könnte nach Heilung des
Netzpfads eine im Sendepuffer hängengebliebene Zustellung nachholen — genau
die Zweideutigkeit, die der Negativbeleg ausschließen soll. Ein neuer
Prozess kann strukturell nichts empfangen, das vor seiner eigenen
Subscription publiziert wurde (dieselbe Begründung wie beim Boundary-Beleg
aus `slice-054`).

Zwei Punkte dazu, beide **ohne** Finding:

- **Plan-Wortlaut vs. Umsetzung:** §1/DoD von `slice-055` sprechen von „real
  wiederverbinden". Der Testablauf verbindet die ursprünglich getrennte
  Test-Instanz nie im Wortsinn neu ans Netz zurück (der getrennte Container
  wird nach dem Negativbeleg entfernt, nicht rekonnektiert). Das ist eine
  bewusste, begründete Auslegung von „Wiederverbindung" als
  „ein Consumer ist danach wieder aktiv verbunden", nicht als „exakt dieselbe
  TCP-Sitzung wird geheilt" — und diese Auslegung deckt sich mit §1s
  explizitem Ausschluss „Automatische NATS-Client-Reconnect-Logik im
  Produktionscode": Ein literales Reconnect derselben Instanz hätte
  zwangsläufig das eingebaute `nats.go`-Reconnect-Verhalten der
  Client-Bibliothek ausgelöst, dessen Testen dieser Slice explizit nicht zum
  Ziel hat. Die gewählte Umsetzung bedient das Ziel (SQL-Nachholung
  belegen, keine stille NATS-Nachlieferung) ohne diese Grenze zu verletzen.
- **Fitness-für-Zweck:** Der Negativbeleg selbst (kein `RECEIVED` während der
  Trennung, verpasste Change nur über `cdc.changes` sichtbar) ist von dieser
  Designwahl unabhängig und bereits vor dem Subscriber-Wechsel vollständig
  erbracht (Zeilen 1394–1408 des Diffs). Der frische Subscriber dient nur dem
  Zusatzbeleg „Wecksignal funktioniert nach der Unterbrechung wieder normal",
  der über die reine DoD-Anforderung hinausgeht.

## Negativbefunde

- geprüft, ohne Befund: **Log-Beleg für Abwesenheit während der Trennung.**
  Der Check `printf '%s' "$nats_reconnect_before_output" | grep -qF
  "RECEIVED"` (Zeile ~1408) ist eine **explizite** Prüfung auf Abwesenheit
  des Frames, nicht nur eine implizite Annahme — bei einem Treffer bricht
  der Testlauf mit `exit 1` ab. Kombiniert mit dem realen
  `docker network disconnect`+`docker inspect`-Beleg (F-1) ist die
  Negativ-Aussage sowohl mechanistisch als auch am Log verankert.
- geprüft, ohne Befund: **Vollständigkeit des Nachhol-Belegs.** Die
  verpasste Change (`id=240`) wird explizit gegen `cdc.changes` per
  `source_id`/`table_name`/`new_data->>'id'` geprüft (`count(*) = 1`), bevor
  der neue Subscriber überhaupt startet — die SQL-Nachholung ist damit
  unabhängig vom späteren Reconnect-Beleg bewiesen.
- geprüft, ohne Befund: **Nur-neues-Signal-Prüfung (mit Einschränkung, siehe
  F-2).** `natssub` wartet strukturell auf genau **ein** `NextMsg` und
  beendet sich danach; zwei RECEIVED-Zeilen im selben Prozess sind
  unmöglich. Die inhaltliche Zuordnung „dieses Signal gehört zu `id=241`,
  nicht zu `id=240`" ist nicht am Payload nachweisbar (leerer Payload,
  `ADR-0055` Punkt 3), sondern folgt aus der Fire-and-Forget-Semantik von
  Core NATS und der fehlenden Retry-Logik (`ADR-0055` Punkt 4) — beides
  außerhalb dieses Diffs bereits entschieden/getestet.
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7).** Der neue
  26-zeilige Kommentar (Zeilen 1333–1358) beschreibt den gegenwärtigen
  Mechanismus und die Designentscheidung im Indikativ bzw. mit begründendem
  Konjunktiv über ein **hypothetisches Risiko einer verworfenen
  Design-Alternative** (nicht über einen früheren/entfernten Code-Zustand,
  da die Reconnect-Belieg-Logik neu ist) — dieselbe Form, die
  `docs/reviews/review-slice-054.md` F-1 bereits als zulässig eingestuft hat
  (Abgrenzung-Klasse: was der Beleg zeigt und was er bewusst nicht so tut).
  Keine `slice-<NNN>`/`welle-<NN>`-Referenz, keine der drei HIGH-Kommentar-
  Klassen trifft zu.
- geprüft, ohne Befund: **Out-of-Scope-Einhaltung (§1 des Slice-Plans).**
  `git show --stat 2987665` bestätigt genau eine geänderte Datei
  (`tools/harness/run-integration-tests.sh`, nur Insertionen); keine
  Änderung an `ChangeNotificationPort`/`internal/adapters/driven/natsnotify/`,
  keine neue Compose-Verdrahtung, kein Reconnect-Verhalten im
  Produktionscode.
- geprüft, ohne Befund: **Cleanup-Vollständigkeit.** Jeder `exit 1`-Pfad im
  neuen Block hat vor sich ein passendes `docker rm -f … || true` für den
  jeweils gerade aktiven Test-Container (`NATS_RECONNECT_BEFORE_CONTAINER`
  bzw. `_AFTER_CONTAINER`); keine verwaisten Container auf einem
  Fehlerpfad.
- geprüft, ohne Befund: **ID-Kollision.** `id=240`/`id=241` sind repo-weit
  einmalig (Grep gegen `feed_mvp_full`-Inserts bestätigt keine
  Wiederverwendung); `NATS_SUBJECT` wird unverändert aus dem bestehenden
  Happy-Path-Block übernommen und trägt bereits das Vier-Token-Schema aus
  `ADR-0056`.
- geprüft, ohne Befund: **Docker-only (`AGENTS.md` §3.1).** Ausschließlich
  `docker exec`/`docker run`/`docker inspect`/`docker network
  disconnect`/`docker logs` gegen bereits laufende, digest-gepinnte
  Container; keine lokale Toolchain-Installation.
- geprüft, ohne Befund: **Commit-Trailer/Traceability.** Weder `2987665`
  noch `1a75087` tragen `Co-Authored-By:`/`Claude-Session:`-Zeilen; beide
  Betreffs referenzieren `LH-FA-SST-007`/`ADR-0055`, keine `SPEC-*`/`ARC-*`-ID
  im Betreff.
- geprüft, ohne Befund: **Nackte ID-Referenzen in diesem Report und im
  Diff.** Alle `LH-*`/`ADR-*`-Vorkommen in diesem Report sind in
  Inline-Code-Backticks gesetzt (Form wie `docs/reviews/review-slice-054.md`,
  von `d-check` mit 0 Befunden akzeptiert); die geänderten Dateien selbst
  (Shell-Skript, Slice-Plan-Checkboxen) führen keine neuen Markdown-Prosa
  mit `LH-*`/`ADR-*`-IDs ein.
- geprüft, ohne Befund: **reale Sensor-Läufe, in diesem Review-Lauf selbst
  ausgeführt, Exit-Code separat geprüft (nicht durch Pipe maskiert).**
  `make gates` — Exit `0` (Kommando-Ende `echo "EXIT_CODE=$?"` separat
  geloggt): `baseline-verify` OK, `docs-check`/`d-check` 429 Dateien/0
  Befunde (voller Lauf und Commit-Range), `commit-traceability` OK (5
  Commits, Betreffs ohne Struktur-ID), `coverage-gate` 41,00 % ≥ 35 % OK,
  `a-check` 0 Befunde (plus dem vorbestehenden, unveränderten
  Abdeckungshinweis zu `tools/harness/natssub/main.go`, keine neue
  Schicht-Verletzung). `make test-integration` — im Hintergrund gestartet
  wegen des 120s-Timeouts der interaktiven Shell, `REAL_EXIT=0` separat in
  eine Log-Datei geschrieben (kein Wrapper-Fehlsignal); der neue
  Negative-Beleg-Block erscheint im Log mit seiner Erfolgszeile, die
  abschließende `TestMVPSchemaChangeIncompatibleTypeChange` lief ebenfalls
  `PASS`. Kein `git gc` nötig — `make gates` lief beim ersten Versuch sauber
  durch, das vom Implementer berichtete Repack-Problem trat in diesem
  Review-Lauf nicht auf.
- geprüft, ohne Befund: **`.harness/skills/reviewer.md`-Kontext-Eingang
  vollständig** (Slice-Plan, `ADR-0055`, `ADR-0056`, `LH-FA-SST-007`,
  `AGENTS.md` §3, `harness/conventions.md` MR-000, vorheriges Review am
  gleichen Modul).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Mechanismus-Realität des
docker-inspect-Trennungsbelegs (real reproduziert, verifiziert) ·
Negativbeleg stützt sich auf externe architektonische Garantie statt auf
inhaltliche Unterscheidung.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW. Beide INFO-Punkte
sind Beobachtungen ohne erwartete Aktion; F-1 bestätigt sogar aktiv die
Belastbarkeit der gewählten Beleg-Mechanik.

**`docker inspect`-Belastbarkeit:** Belastbar und real reproduziert (F-1) —
der sofortige Wegfall des Netzwerk-Eintrags in `docker inspect` fällt mit
dem realen, sofortigen Verbindungsabbruch auf Paketebene zusammen; kein
bloßer Absichts-Beleg.

**Design-Entscheidung „frischer Subscriber":** Sauber begründet und mit dem
expliziten Out-of-Scope-Punkt „keine Produktionscode-Reconnect-Logik" des
Slice-Plans konsistent — kein stiller Plan-Widerspruch, sondern eine
plausible Auslegung von „Wiederverbindung" innerhalb der im Plan selbst
gezogenen Grenze.

**Fixrunde:** nicht nötig. Nach `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde wird die DoD-Zeile „Review durchgeführt"
im Slice-Plan in diesem Commit mitgezogen. Die beiden Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler (jeweils Erstauftreten — keine Registrierung
durch den Reviewer selbst, das ist Planner-Arbeit bei Closure). Dieser
Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und wird über Läufe hinweg nicht wieder
gelesen. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11).
