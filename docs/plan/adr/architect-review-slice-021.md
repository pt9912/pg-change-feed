# Architect-Review slice-021 — Zugriffsweg für Consumer-Registrierung: ADR-0019 + ADR-0020 + ADR-0046 tragen die Frage bereits, kein neues ADR

**Rolle:** Architect (Modul 8). **Datum:** 2026-09-12.

**Eingang:** Slice-Plan
[`docs/plan/planning/done/slice-021-consumer-registrierung-zugriffsweg.md`](../planning/done/slice-021-consumer-registrierung-zugriffsweg.md)
§1–§8 · [`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md)
`LH-FA-CON-001.a`/`LH-FA-CON-004.a` · [`spec/lastenheft.md`](../../../spec/lastenheft.md)
`LH-FA-CON-001`, `LH-FA-SST-003`, `LH-FA-SST-005` · [`spec/architecture.md`](../../../spec/architecture.md)
§4 (Sequenz-Diagramme, insbesondere `LH-FA-REA-002` und `LH-FA-CON-004`),
`ARC-003`/`ARC-005` · [`ADR-0018`](0018-sql-driving-adapter.md) (Superseded
by `ADR-0046`) · [`ADR-0019`](0019-cli-driving-adapter.md) (Accepted,
`permanent`) · [`ADR-0020`](0020-http-grpc-optional.md) (Accepted, Trigger
„beobachtbarer API-Consumer-Bedarf") · [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)
(Accepted, `permanent`, benennt die offene SQL-Bridge-Lücke explizit) ·
`cmd/pg-change-feed/main.go` · `internal/bootstrap/wiring.go` ·
`internal/application/port/inbound/consumer.go` · `compose.yaml` ·
`docs/plan/planning/observations/BEO-PGC/` (alle elf Verzeichnisse, kein
Treffer zu HTTP/gRPC/API-Bedarf) · Baseline-Regelwerk `modul-08-agentenrollen.md`
§Rollen-Regeln, §Konflikt-Pfad als Rollen-Sequenz.

**Ausgang:** Kein neues ADR nötig. Der Zugriffsweg ist **CLI-Unterbefehl**
(z. B. `register-consumer <name>`), der `RegisterConsumerUseCase` direkt im
selben Go-Prozess über den Inbound Port aufruft — das ist keine neue
Entscheidung, sondern die Zusammenführung dreier bereits `Accepted`-ADRs,
die den Optionsraum bereits erschöpfend abdecken:

| Option | Bereits entschieden durch | Verdikt für slice-021 |
|---|---|---|
| **A — CLI-Unterbefehl** | [`ADR-0019`](0019-cli-driving-adapter.md): CLI verwendet dieselben Inbound Use Cases wie jeder andere Driving Adapter; physisch trivial, weil CLI im selben Go-Prozess läuft und den Port ohne Brücke direkt aufruft | **Gewählt** |
| B — Netzwerkschnittstelle (HTTP/gRPC) | [`ADR-0020`](0020-http-grpc-optional.md): „nicht zwingend zum MVP", nur bei beobachtbarem API-Consumer-Bedarf; ihr Re-Evaluierungs-Trigger ist nicht eingetreten — kein Eintrag im Beobachtungs-Register, keine Lastenheft-Änderung, die einen solchen Bedarf benennt | Ausgeschlossen für diesen Slice |
| C — schreibende SQL-Funktion | [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md): schreibende SQL-Funktionen müssen über Inbound Ports laufen, aber die dafür nötige physische Brücke (FDW/`dblink`/Extension mit Prozess-/Socket-Zugriff) existiert nicht und ist nirgends entschieden — die ADR benennt diese Lücke ausdrücklich als **offen** und verweist sie an „den Slice, der das `ADR-0018`/`ADR-0019`-Rest umsetzt" | Ausgeschlossen — würde eine ungelöste Architektur-Lücke voraussetzen, die dieses Slice laut eigenem §1 nicht lösen muss (es genügt, sie nicht zu wählen) |
| D — Direktes Schreiben der CDC-Tabellen (Status quo) | Ist der von `LH-FA-CON-001.a` benannte **Missstand** selbst, keine Option | Verworfen (Ausgangslage, kein Zugriffsweg) |

**Harte Regel eingehalten:** Keine der drei herangezogenen `Accepted`-ADRs
(`ADR-0019`, `ADR-0020`, `ADR-0046`) wird von diesem Lauf inhaltlich
geändert — dieser Bericht wendet sie nur an. Der Slice-Kopf (`Bezug:`-Feld)
wird von diesem Lauf **nicht** editiert — das ist Planner-Arbeit nach
Übergabe dieses Verdikts (Modul 8, Planner→Architect→Planner-Zug).

---

## 1. Warum das keine neue Entscheidung ist, sondern eine Anwendung

`ARC-005` benennt den Optionsraum für Driving Adapters bereits erschöpfend:
„Replication Stream, CLI, SQL-Funktionen/Views, später HTTP-/gRPC."
Replication Stream scheidet aus (reiner Capture-Kanal, `ADR-0006`, keine
Rolle für Consumer-seitige Aktionen). Von den verbleibenden drei
Kandidaten sind zwei bereits durch unabhängige, nicht miteinander in
Konflikt stehende `Accepted`-ADRs **aus anderem Anlass** geschlossen:

- `ADR-0020` schließt die Netzwerkschnittstelle nicht grundsätzlich aus,
  sondern macht ihre Öffnung von einem **beobachtbaren** Ereignis
  abhängig (API-Consumer-Bedarf). Ich habe das gesamte
  Beobachtungs-Register (`docs/plan/planning/observations/BEO-PGC/`, elf
  Verzeichnisse) durchgesehen: keines davon trägt einen HTTP-/gRPC-/
  API-Bezug. Der Trigger ist nicht eingetreten.
- `ADR-0046` schließt die SQL-Funktion für **schreibende** Zugriffe nicht
  aus Geschmack, sondern weil die Entscheidung, wie eine reine
  SQL-Funktion einen Go-Inbound-Port synchron aufrufen soll, physisch
  ungelöst ist — und die ADR sagt selbst, dass diese Lösung nicht ihr
  Gegenstand ist, sondern der eines späteren Slices. Dieser Slice muss
  diese Lücke nicht schließen, um sein eigenes Ziel (`LH-FA-CON-001.a`)
  zu erreichen; er muss nur vermeiden, sie als Voraussetzung zu wählen.

Damit bleibt CLI die einzige Option, die (a) bereits als Muster
`Accepted` ist, (b) ohne neue Infrastruktur umsetzbar ist und (c) keine
offene Architektur-Frage voraussetzt. Das ist eine Ableitung aus
bestehenden Entscheidungen, kein neues Urteil über einen bislang
unentschiedenen Sachverhalt — und genau das ist die Voraussetzung dafür,
mit einem Verdikt statt einem neuen ADR zu antworten (Modul 8: ein
Accepted-ADR wird nicht überschrieben, aber sie darf angewendet werden,
ohne dass die Anwendung selbst eine neue ADR wäre).

**Gegenprobe — wäre hier tatsächlich eine neue ADR nötig gewesen?** Ja,
wenn eine der drei folgenden Bedingungen zuträfe: (1) `ADR-0019`s
Geltungsbereich wäre explizit auf `LH-FA-SST-003` (Installation/Diagnose/
Administration) begrenzt und Consumer-Registrierung fiele erkennbar
nicht darunter — trifft nicht zu, `ADR-0019`s Entscheidungssatz ist
kanalgenerisch formuliert und ihr Re-Evaluierungs-Trigger ist `permanent`,
nicht an einzelne LH-IDs gebunden. (2) Der Optionsraum wäre durch die
drei herangezogenen ADRs nicht erschöpft — trifft nicht zu, `ARC-005`
listet exakt die vier genannten Kandidaten, drei sind durch ADRs
abgedeckt, der vierte (Replication Stream) ist strukturell irrelevant.
(3) Eine der drei ADRs hätte einen fälligen Re-Evaluierungs-Trigger —
geprüft für `ADR-0020` (kein Register-Treffer, kein Lastenheft-Change),
`ADR-0046` (die Lücke ist als offen benannt, aber ihr Fehlen ist genau
das Contra-Argument gegen Option C, kein Trigger für eine Neubewertung
von `ADR-0046` selbst — die ADR sagt korrekt voraus, dass sie offen
bleibt, bis eine Brücke existiert), `ADR-0019` (`permanent`, kein
Ereignis definiert). Keine der drei Bedingungen trifft zu.

## 2. Konkrete Anwendung — Unterbefehl-Namen, Prozess-Lebenszyklus im Container

Diese Sektion beantwortet die im Slice-Plan offene Frage „wie" der
Zugriffsweg konkret aussieht — Implementierungsdetail, kein zusätzlicher
ADR-Gegenstand:

- **Unterbefehl statt Flag.** `cmd/pg-change-feed/main.go` kennt heute
  zwei Sondermodi (`--version`, `--healthcheck`, Zeilen 21–46), die früh
  prüfen und beenden, bevor der Capture-Loop (`bootstrap.Run`) startet.
  Ein dritter Sondermodus `register-consumer <name>` reiht sich dort ein:
  liest dieselben Verdrahtungs-Vorbedingungen (`bootstrap.ConfigFromEnv`),
  baut die Verdrahtung bis zum `RegisterConsumerUseCase`, ruft `Register`
  auf, druckt das Ergebnis (Consumer-Kennung, `AlreadyRegistered`-Flag)
  und beendet sich — **ohne** je `bootstrap.Run` (den Dauerbetrieb) zu
  erreichen. Ein Unterbefehl (nicht ein Flag mit Wert) passt hier besser,
  weil er eine Befehlsfamilie eröffnet, die `slice-022` mit
  `acknowledge-consumer` fortsetzt.
- **Container-Betrieb.** Derselbe Image-Tag wie der Daemon
  (`ghcr.io/pt9912/pg-change-feed:dev`, `compose.yaml`), dieselbe
  `CDC_SOURCE_DSN`-Umgebungsvariable — heute die gemeinsame Instanz-DSN
  (siehe `BEO-PGC/rollen-verdrahtung`, vom Slice in §6 als unverändertes
  Risiko benannt) —, aufgerufen als **einmaliger, kurzlebiger Lauf**
  (`docker run --rm <image> register-consumer <name>` bzw.
  `docker compose run --rm feed register-consumer <name>`), nicht als
  Dauerdienst in `compose.yaml`. Das Runtime-Image ist distroless (kein
  Shell); die exec-Form von `CMD`/`ENTRYPOINT` ruft das Binary direkt mit
  Argumenten auf — dasselbe Muster, das der bestehende
  `--healthcheck`-Aufruf bereits nutzt (Kommentar in `main.go`,
  Zeilen 26–35).
- **Exit-Codes** folgen dem bestehenden Muster (`0` Erfolg, `1`
  Verdrahtungs-/Domänenfehler, `2` falsche Argumente, siehe die
  bestehende `unbekanntes Argument`-Behandlung, Zeile 44) — Implementer-
  Entscheidung im Detail, kein ADR-Gegenstand.

## 3. Anwendung auf den Bestätigungs-Fall (`slice-022`, informativ)

Die hier angewendete Mechanismus-Entscheidung (CLI-Unterbefehl über den
Inbound Port, kein SQL-Funktions-Weg, keine Netzwerkschnittstelle) ist
kanalgenerisch (`ADR-0019`) und gilt für `AcknowledgeConsumerUseCase`
identisch — voraussichtlich `acknowledge-consumer <consumer-id>
<position>`. `slice-021` schließt diesen Fall in §1 ausdrücklich aus
(Folge-Slice `slice-022`); dieser Hinweis erspart `slice-022` lediglich
einen erneuten Architect-Rundlauf für dieselbe Mechanismus-Frage — die
konkrete Verdrahtung von `AcknowledgeConsumerUseCase` bleibt
`slice-022`s eigene Arbeit.

## 4. Beobachtung für den Planner (kein Blocker für diesen Slice)

`spec/architecture.md` §4 trägt bereits ein Sequenzdiagramm für
`LH-FA-CON-004` (Consumer-ACK, Kanal „CLI/SQL" — konsistent mit dem hier
bestätigten CLI-Weg; der SQL-Anteil bleibt dort aspirational, solange
`ADR-0046`s Lücke offen ist, und diese ADR ändert das Diagramm nicht).
Für `LH-FA-CON-001` (Registrierung) existiert dagegen **kein**
Sequenzdiagramm — der Abschnitt zeigt Capture, Lesen, ACK und
Tabellen-Aktivierung, aber keine Registrierung. Das ist keine
Inkonsistenz, die dieses Verdikt auflösen muss (die Sicht ist
meilensteinfrei und wird bei Bedarf fortgeschrieben, `AGENTS.md` §3.4),
aber eine Beobachtung, die der Planner einordnen sollte: entweder als
Ergänzung des `§3`-Plans dieses Slices (ein weiteres Diagramm neben den
bestehenden vier) oder als eigener, benannter Nachtrag. Ich treffe hier
keine Entscheidung darüber, welcher Weg — das ist Planungsarbeit, keine
Architektur-Entscheidung.

---

## Disposition (Gesamt)

| Frage | Ergebnis |
|---|---|
| Trägt eine bestehende `Accepted`-ADR die Zugriffsweg-Frage bereits? | Ja — `ADR-0019` (CLI-Muster) + `ADR-0020` (Netzwerkschnittstelle ausgeschlossen ohne Bedarf) + `ADR-0046` (SQL-Funktion physisch ungelöst) schließen den von `ARC-005` benannten Optionsraum erschöpfend |
| Neues ADR nötig? | Nein |
| Gewählter Zugriffsweg | CLI-Unterbefehl (`register-consumer <name>`), ruft `RegisterConsumerUseCase` direkt über den Inbound Port |
| `ADR-0020`-Trigger eingetreten? | Nein — kein Register-Treffer, keine Lastenheft-Änderung |
| `ADR-0046`-Lücke durch diesen Slice berührt? | Nein — dieser Slice wählt bewusst nicht den SQL-Weg und muss die Lücke daher nicht schließen |
| Auswirkung auf `slice-021`-§4-Trigger (zu groß?) | Keine — CLI-Unterbefehl löst den „Netzwerkschnittstelle mit eigenem Protokoll"-Rückführungs-Trigger nicht aus |
| Auswirkung auf `slice-021`-Closure | Keine Blocker — DoD-Punkt „ADR entschieden" ist mit diesem Verdikt erfüllt, Planner trägt den Bezug im Slice-Kopf nach |

---

## Beleg-Anker (Kurzfassung)

| Aussage | Beleg |
|---|---|
| CLI ruft dieselben Inbound Use Cases wie jeder Driving Adapter, `permanent` | [`ADR-0019`](0019-cli-driving-adapter.md) §Entscheidung, §Re-Evaluierungs-Trigger |
| Netzwerkschnittstelle nur bei beobachtbarem API-Consumer-Bedarf | [`ADR-0020`](0020-http-grpc-optional.md) §Entscheidung, §Re-Evaluierungs-Trigger |
| Kein Register-Eintrag zu HTTP/gRPC/API-Bedarf | `docs/plan/planning/observations/BEO-PGC/` (eigene Durchsicht aller elf Verzeichnisnamen) |
| SQL-Funktions-Bridge physisch ungelöst, Lücke ausdrücklich offen und an diesen Slice verwiesen | [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md) §Kontext („Was diese ADR nicht löst"), §Konsequenzen |
| `ARC-005` benennt den vollständigen Optionsraum | `spec/architecture.md` Zeile 58 (`ARC-005`-Tabellenzeile) |
| [`LH-FA-CON-001.a`](../../../spec/pflichtenheft.md)/[`LH-FA-CON-004.a`](../../../spec/pflichtenheft.md) benennen den Missstand (Direktschreiben umgeht Prüfung/Invariante) | `spec/pflichtenheft.md` (siehe Link in der linken Spalte) |
| `RegisterConsumerUseCase` bislang nicht verdrahtet | `internal/bootstrap/wiring.go` (kein Treffer für „Register") |
| Bestehendes Sondermodus-Muster in `main.go` (`--version`/`--healthcheck`) | `cmd/pg-change-feed/main.go:20-46` |
| Kein Sequenzdiagramm für `LH-FA-CON-001` in der Sicht | `spec/architecture.md` §4 (vier Diagramme: Persist-before-ACK, Lesen, ACK, Tabelle aktivieren — keines für Registrierung) |
