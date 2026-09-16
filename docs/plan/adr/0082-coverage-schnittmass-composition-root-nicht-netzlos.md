# ADR-0082: Coverage 80 % — Schnittmaß; Composition Root netzlos nicht prüfbar

**Status:** Accepted — **kein** Supersedes. Diese ADR führt keine Festlegung
der [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
ab. Sie **wendet** deren Gegenstands-Regel auf einen Fall an, den jene ADR
nicht im Blick hatte — der Composition Root liegt **innerhalb** des
Messgegenstands, nicht in einem der drei ausgenommenen Pakete —, und
**bestätigt** Gegenstand, Endstufe und Rampe im Übrigen unverändert.

**Datum:** 2026-09-16

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Planner-Zug, der
`welle-20` eröffnet hat, und als der Implementer-/Verifier-Lauf des
postgresack-Naht-Slice, dessen Zahlen hier unabhängig nachgerechnet werden —
Modul 8 §Rollen-Regeln)

**Bezug:** [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Gegenstand, Endstufe, Punkt 5) ·
[`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md) ·
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(Rampen-Mechanik, unberührt) ·
[`ADR-0080`](0080-nahtform-pgconn-adapter-treiberhuelle.md) (Nahtform der
Adapter — hier **nicht** auf den Composition Root übertragen) ·
[`ADR-0026`](0026-composition-root.md) (der Composition Root als der eine Ort
der Verdrahtung) · [`ADR-0030`](0030-testpyramide.md) (Unit-Tier gegen
DB-gestützten Tier) · `AGENTS.md` §3.1 (Docker-only), §3.5 (Accepted-ADRs
sind immutable), §3.6 (Schwellen nur per ADR), §3.7 (Ist-Zustand) ·
`harness/mk/coverage.mk` (`THRESHOLD`) · `tools/coverage-gate.sh` ·
`tools/harness/db-coverage.sh` (die Zählbasis dieses Zugs) · `Dockerfile`
(Stufe `coverage`) · `internal/bootstrap/wiring.go` (`Run`, `Diagnose`,
`Healthcheck`, `RegisterConsumer`, `AcknowledgeConsumer`) ·
`internal/adapters/driven/postgresstorage/store.go` (der `Ping`-Riegel) ·
`cmd/pg-change-feed/main.go` · `docs/plan/planning/welle-20.md` ·
`docs/plan/planning/observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/observation.md`
· `docs/plan/planning/observations/BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code/observation.md`

**Schärft:** — (Prozess-/Test-ADR ohne Spec-Stratum, wie `ADR-0054`,
`ADR-0071`, `ADR-0077`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass: eine Welle, die auf ihr eigenes Schnittmaß wartet.**
`welle-20` ist eröffnet; ihr Start-Trigger liegt in `done/`, ihre
Slice-Liste führt einen Slice, und ihr eigener Text benennt die Lücke:
*„Als Nächstes zu schneiden — der Schnitt folgt der Messung. … dann entstehen
die Test-Slices nach dem Größenmaß aus dem Schnittvorschlag des Verdikts
(Statement-Anteil)."* Der Architect-Zug, der dieses Maß liefert, hat nie
stattgefunden. Diese ADR ist er.

**(2) Die Messung ist nachgerechnet, nicht übernommen — und sie schwankt.**
Aus dem Profil der `coverage`-Stufe des gepinnten Images, dedupliziert über
die Block-Position (dieselbe Zählbasis, die `tools/harness/db-coverage.sh`
beschreibt), in **zwei Läufen über denselben Quelltext**:

| Größe | Lauf 1 | Lauf 2 |
|---|---|---|
| Nenner (Summe der Statement-Zahlen) | **1903** | **1903** — exakt der Gate-Wert |
| gedeckt | **1369** | **1371** |
| Ist-Stand | **71,94 %** | **72,04 %** (das Gate druckt 71,9 % / 72,0 %) |
| ungedeckt | **534** | **532** |
| für die Endstufe 80 % nötig (1523) | **+154** | **+152** |

Die beiden Läufe unterscheiden sich in **genau einer** Funktion:
`runWALRetentionCheck` (`internal/bootstrap/wiring.go:981`) trägt in Lauf 1
14 von 16, in Lauf 2 16 von 16 Statements — der Takt-Zweig seiner Schleife
feuert nur, wenn der Tick vor dem Kontext-Ende liegt. **Der Nenner ist stabil,
der Zähler nicht.** Alle Zahlen dieser ADR sind auf den **ungünstigen** Lauf
gerechnet (1369 gedeckt, 534 ungedeckt, +154 nötig); wo beide Läufe
auseinandergehen, nennt der Text beide Enden. Die Fortschritts-Aussage dieser
Welle bleibt die **Quote**, nicht eine eingefrorene Statement-Zahl — die
Klasse `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` steht genau dafür.

**(3) Warum `Run` dasteht, wie es dasteht — der `Ping`-Riegel.** `Run`
(`internal/bootstrap/wiring.go:414`) ist die Verdrahtungsfunktion des
Composition Root; im Profil trägt sie **0 von 198** Statements. Der Grund ist
kein Testmangel, sondern die Bauform: ihr erster Zug nach dem Log-Aufbau ist
`postgresstorage.New(ctx, cfg.CaptureDSN, …)`, und jeder der
`postgresstorage.New*`-Konstruktoren ist `pgxpool.New` **plus
`pool.Ping(ctx)`** (`store.go:35`); `receive.NewStream`,
`receive.NewWALRetentionChecker`, `NewAdministrationListener` und
`nats.Connect` verlangen ebenso einen lebenden Dienst. Netlos erreichbar sind
deshalb **9 von 198** Statements — genau der Fehlerpfad des ersten
Konstruktors: der Log-Aufbau, der fehlschlagende Verbindungsaufbau, der
`return err` und die Fehler-Zeile im `defer`. Alles Übrige liegt hinter dem
Riegel.

**(4) Die Messung auf Paket- und Funktionsebene.**

*(4a) Die sechs dienstgebundenen Funktionen — 357 der ungedeckten Statements
(67 %):*

| Funktion | ungedeckt / gesamt | netzlos erreichbar |
|---|---|---|
| `Run` (`internal/bootstrap`) | **198 / 198** | **9** (Fehlerpfad) |
| `Diagnose` (`internal/bootstrap`) | 73 / 81 | 0 (über den bereits gedeckten Anteil hinaus) |
| `main` (`cmd/pg-change-feed`) | **49 / 49** | **≈45** (Argument-Dispatch, Konfigurations-Fehlerpfade) |
| `Healthcheck` (`internal/bootstrap`) | 14 / 22 | 0 |
| `RegisterConsumer` (`internal/bootstrap`) | 13 / 17 | 0 |
| `AcknowledgeConsumer` (`internal/bootstrap`) | 10 / 14 | 0 |
| **Summe** | **357** | **≈54** |

*(4b) Der Tail — 175 (Lauf 2) bzw. 177 (Lauf 1) ungedeckte Statements
(≈33 %), über rund 30 Pakete; größter Einzelposten 8 Statements:*

| Paket | ungedeckt | darunter |
|---|---|---|
| `internal/bootstrap` (netzlos prüfbarer Rest) | **26 – 28** | `mergeConfig` 6, `applyAdministrationRequest` 6, `splitQualifiedName` 4, `activatedTableBindings` 4, `processAdministrationRequests` 4, `runWALRetentionCheck` 0 – 2 |
| `internal/adapters/driving/http` | 29 | `streamChangesHandler` 5, `Start` 4, `runRetentionHandler` 4, `removeConsumerHandler` 4 |
| `internal/adapters/driving/replication/mapper` | 25 | `observeRelation` 8, `Consume` 6 |
| `internal/adapters/driving/grpc/streamv1` (erzeugt) | 25 | Getter 10, Client-/Handler-Stubs 5, Proto-Runtime-Interna ≈10 |
| `internal/adapters/driven/postgresstorage/sqlexec` | 17 | Fehlerzweige in `translate.go` |
| `internal/adapters/driving/replication/decode` | 14 | `Decode` 8, `tupleValues`/`oldTupleValues` je 3 |
| 20 weitere Pakete (Use-Cases, `natsnotify`, `domain/model`, `telemetry`, `grpc`, `postgresstorage/mapper`) | 39 | je 1–5 |

**Der Composition Root trägt den Hebel nicht.** `internal/bootstrap` hat 336
ungedeckte Statements — aber **308 davon** liegen in den fünf dienstgebundenen
Funktionen aus (4a); netzlos prüfbar sind aus diesem Paket nur **26–28**. Wer
den Schnitt auf `internal/bootstrap` legt, legt ihn auf den unbeweglichen
Teil.

**(5) Die Decke — benannt, nicht verschwiegen.** Der Messgegenstand ist eine
**Paket**-Partition ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 1, [`ADR-0080`](0080-nahtform-pgconn-adapter-treiberhuelle.md)
§Kontext (5)); die Eigenschaft, die ihn trägt, ist **code**-granular („Code,
dessen Testlauf einen externen Dienst voraussetzt"). Für die drei
ausgenommenen Pakete fallen beide Ebenen zusammen; für `internal/bootstrap`
und `cmd/pg-change-feed` **nicht** — beide laufen netzlos und enthalten
zugleich die sechs Funktionen aus (4a). Daraus folgt:

| Größe | Wert |
|---|---|
| Decke **ohne** jede Test-Arbeit an den sechs (die sechs bleiben stehen) | 1546/1903 = **81,2 %** |
| Decke **mit** ihren netzlos erreichbaren Präfixen (≈54) | 1600/1903 = **84,1 %** |
| Realistische Decke (Tail ohne die ≈10 Proto-Runtime-Interna) | 1590/1903 = **83,6 %** |
| Bedarf (Endstufe 80 %) | **1523** |

**Beide Läufe ergeben dieselbe Decke** — die Schwankung aus (2) verschiebt
sich zwischen Tail und gedeckter Menge und hebt sich auf. **80 % ist
erreichbar — aber nicht durch `Run`.** Die Welle trägt auf dem Tail: die
**165 bis 167** realistisch erreichbaren Tail-Statements decken die fehlenden
**152 bis 154** mit **13 Statements Puffer** — an beiden Enden der Schwankung
derselbe Abstand (0,7 pp). Mit den Prozess-Rand-Präfixen wächst der Puffer
auf **63 Statements (3,3 pp)**. Die Decke bestätigt zugleich den Rot-Beleg,
den die Welle verlangt: bei `THRESHOLD=85` (1618 Statements) bleibt selbst die
realistische Decke rot — die Stufe prüft real.

**(6) Zwei offene Beobachtungen warten auf genau diese Aussage.**
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` steht bei 2× und nennt
als möglichen Träger „einen schlanken Bootstrap-Smoke-Test (`ConfigFromEnv` +
`Run()`-Teilaufruf gegen einen minimalen ENV-Satz, der prüft, dass kein
Adapter-Konstruktor einen `nil`-Use-Case erhält)" — mit dem Vermerk, dass
„das erst ein Architect-Zug entscheidet". Und
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` hält fest, die
naheliegende Antwort auf die Gegenstands-Frage wäre „eine Neudefinition des
Gegenstands (nur Code, dessen Test eine Verbindung braucht)", und das sei
„eine **Entscheidung** über eine Messfläche … als Architect-Frage behandelt,
nicht als Notiz". Dieser Zug beantwortet die erste (Festlegung 5) und
**vertagt** die zweite ausdrücklich (Festlegung 4).

## Entscheidung

Wir wählen: **Das Schnittmaß dieser Welle ist der Tail; der Composition Root
bleibt unangetastet.** Fünf Festlegungen:

**1. `Run` ist netzlos nicht prüfbar — und es wird keine Naht für sie
gezogen.** Gemessen (Kontext (3)): 9 von 198 Statements sind netzlos
erreichbar; die einzige Tür zu den übrigen 189 ist eine Naht am Composition
Root. Beide Kandidaten werden verworfen:

- **Konstruktor-Injektion** (die `postgresstorage.New*`-Aufrufe als
  Funktionswerte, ein Test liefert Fakes): sie machte den Composition Root zum
  Mock-Ziel — geprüft würde, dass `Run` in einer Reihenfolge aufruft, was der
  Test selbst gestellt hat. Ihre Rechtfertigung wäre allein die Zahl; das ist
  der Goodhart-Fall, den [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Option D selbst verwirft.
- **Zerlegung `build`/`serve`** (Konstruktion von Lebenszyklus trennen): sie
  ist design-fähig — `runHeartbeat`, `runAdministration`,
  `runWALRetentionCheck` sind genau dieses Muster und **gedeckt** —, aber sie
  bewegt nur die Lebenszyklus-Hälfte: **83 der 198** Statements
  (Goroutinen-Start/-Stopp, Shutdown-Reihenfolge, Ergebnis-Merge, die
  HTTP-/gRPC-Server). Die Bau-Hälfte (115) bleibt hinter dem `Ping`-Riegel.
  Voll ausgeschöpft ergäbe sie 1452/1903 = **76,3 %** — **unter** der
  Endstufe. Sie ist damit kein Hebel für diese Welle; ihr Nutzen wäre ein
  Design-Nutzen (prüfbarer Lebenszyklus), nicht dieser.

**2. Das Schnittmaß ist der Tail (175–177); der Prozess-Rand ist der
Puffer.** Die Test-Arbeit dieser Welle greift **nicht** auf
`internal/bootstrap` als Ganzes, sondern auf die in Kontext (4b) benannten
Paket-Cluster. Load-bearing ist der Tail; die netzlos erreichbaren Präfixe
der sechs dienstgebundenen Funktionen (Kontext (4a), ≈54 — allen voran der
Argument-Dispatch von `main`, der über einen Re-Exec-Harness **ohne
Produktionsänderung** prüfbar ist) sind der Puffer, nicht die Last.

**3. Die Endstufe 80 % und der Messgegenstand bleiben unverändert.** 80 % ist
über dem heutigen Gegenstand erreichbar (Kontext (5)); es gibt **keinen**
Grund, den Gegenstand zu schneiden oder die Schwelle zu senken. Die
Rest-Spanne ist benannt, nicht weggerechnet: 357 Statements (18,8 % des
Nenners) sind netzlos nicht erreichbar, und die Decke liegt deshalb bei
83,6 %, nicht bei 100 %.

**4. Der Composition Root bleibt im Gegenstand — seine Decke ist die benannte
Grenze.** `internal/bootstrap` und `cmd/pg-change-feed` erfüllen die
Ausschluss-Eigenschaft aus [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 1 **nicht** („ein Paket, dessen Testlauf einen externen Dienst
voraussetzt"): ihr Testlauf läuft netzlos, einzelne Fälle überspringen. Sie
bleiben also im Gegenstand. Die Alternative — die Partition **code**-granular
zu fassen — ist **nicht** Teil dieser Entscheidung: sie ist für 80 % nicht
nötig und wäre ein eigener, design-getragener Vorgang mit eigener Beweislast
(siehe §„Die benannte Grenze" und §Re-Evaluierungs-Trigger (c)).

**5. Der vorgeschlagene Bootstrap-Smoke-Test ist netzlos nicht baubar.** Ein
`Run()`-Teilaufruf kann die Verdrahtung nicht beobachten: nach dem Log-Aufbau
steht der `Ping`-Riegel, `Run` kehrt im Fehlerpfad zurück, bevor irgendein
Adapter mit einem Use-Case versehen wird. Die Lücke, die
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` beschreibt, hat ihren
Träger bereits **außerhalb** des Unit-Tiers: der Integrations-Rundlauf hat sie
real gefunden und der Folge-Zug sie geschlossen. Ein netzloser Träger
existiert nicht — die Beobachtung ist mit dieser ADR insoweit beantwortet.

### Warum das keine Herabstufung des Composition Root ist

Der Composition Root ist der eine Ort, an dem die konkreten Adapter
zusammengesetzt werden ([`ADR-0026`](0026-composition-root.md)); dass sein
Bau-Teil netzlos nicht erreichbar ist, ist die Eigenschaft dieses Ortes, nicht
ein Defekt seiner Prüfung. Sein Beleg liegt dort, wo er hingehört: im
DB-gestützten Tier (`make test-store`/`make test-replication` fahren `Run`
real) und im Integrations-Rundlauf. Der Unit-Tier kann ihn nicht tragen — und
soll es nicht, denn ein Composition Root, der gegen Fakes prüft, prüft die
Fakes.

### Die benannte Grenze: Paket-Partition gegen code-granulare Eigenschaft

Der Gegenstand ist eine Paket-Partition, weil das Verfahren sie ist
(`-coverpkg` und der `Dockerfile`-Filter greifen auf Paketnamen,
[`ADR-0080`](0080-nahtform-pgconn-adapter-treiberhuelle.md) §Kontext (5)).
Die tragende Eigenschaft ist dagegen eine des Codes. Solange ein Paket
*ganz* dienstgebunden ist (die drei ausgenommenen), deckt die Partition die
Eigenschaft. `internal/bootstrap` und `cmd/pg-change-feed` sind **gemischt**:
ihr netzloser Teil trägt den Gegenstand zu Recht, ihr dienstgebundener Teil
lastet als Leergewicht darin. Daraus entsteht die Decke aus Kontext (5) —
benannt, nicht behoben. Eine code-granulare Partition (etwa: Messung über
diejenigen Block-Positionen, deren Test keine Verbindung braucht) würde die
Decke heben; sie ist eine Änderung am **Messmechanismus** und damit ein
eigener Vorgang, der eine eigene Entscheidung braucht — nicht Beigabe dieser
Welle und für die 80 % nicht nötig.

### Was daraus für die Slices folgt — Maß und Zahl, kein Schnitt

Der Planner schneidet; diese ADR liefert das Maß. Die Test-Arbeit dieser Welle
teilt sich in Cluster, die je einzeln lieferbar sind (jeder hebt die Quote
real):

| Cluster | Pakete | ungedeckt | netzlos erreichbar |
|---|---|---|---|
| **A — Prozess-Rand** | `cmd/pg-change-feed` (`main`-Dispatch), `internal/bootstrap` (`Run`-Fehlerpfad) | ≈54 | ≈54 |
| **B — Reine Übersetzung** | `replication/decode`, `replication/mapper`, `postgresstorage/sqlexec`, `postgresstorage/mapper` | 60 | 60 |
| **C — Zustell- und Betriebs-Rand** | `driving/http`, `driving/grpc`, `driven/natsnotify`, `driving/grpc/streamv1` | 62 | ≈52 |
| **D — Anwendungs-Kern und Bootstrap-Rest** | Use-Cases, `domain/model`, `telemetry`, `bootstrap`-Rest | 53 – 55 | 53 – 55 |
| **Summe** | | **229 – 231** | **≈219 – 221** |

**Zahl: vier Slices — einer je Cluster; bei Schicht-Überschreitung wird
Cluster D an seiner Grenze geteilt (fünf).** Die Last tragen Cluster B, C und
D zusammen (175 – 177 ungedeckt, ≈165 – 167 erreichbar); allein lassen sie
**13 Statements Puffer** — dünner als eine einzelne unerreichbare
Fehlerverzweigung. Cluster A ist deshalb **kein Beiwerk**: er ist der Puffer,
der die Welle von 0,7 pp auf 3,3 pp Abstand zur Schwelle hebt. Der Puffer ist
billig und trägt sich selbst — `main`s Argument-Dispatch ist eine reale
Zusage (die Sondermodi sind im Prozess-Vertrag benannt), nicht
Coverage-Theater.

**Der Weg über `internal/bootstrap` als „Hebel" ist widerlegt.** Die Welle
notiert dort 327 ungedeckte Statements und erwägt zwei Slices; 308 davon
liegen in den fünf dienstgebundenen Funktionen aus (4a) und bewegen sich
durch keinen Test-Slice. Aus diesem Paket sind **26–28 Statements** netzlos
prüfbar — sie sind ein Teil von Cluster D, kein eigener Hebel und kein eigener
Slice.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — Composition Root per Konstruktor-Injektion prüfbar machen | alle 198 Statements von `Run` würden netzlos erreichbar; die Endstufe wäre bequem erreichbar | prüft die Fakes, nicht das Produkt; macht den einen Ort, der die konkreten Adapter zusammensetzt, zum Mock-Ziel — genau der Goodhart-Fall, den [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Option D verwirft; Produktionscode-Refaktorierung, die die Welle sich selbst verboten hat |
| B — `Run` in `build`/`serve` zerlegen und nur `serve` fahren | design-fähig (dasselbe Muster wie `runHeartbeat`/`runAdministration`); hebt die Decke ohne Messänderung | bewegt nur 83 der 198 Statements → 76,3 %, **unter** der Endstufe; die Bau-Hälfte bleibt hinter dem `Ping`-Riegel; löst das Schnittmaß-Problem dieser Welle also nicht — bleibt als design-getragener Vorgang möglich, nicht als ihr Weg |
| C — Gegenstand code-granular neu fassen (nur Code, dessen Test eine Verbindung braucht) | beseitigt das Leergewicht an seiner Ursache; hebt die Decke; würde beide Richtungen der Gegenstands-Frage (auch die Verdünnungs-Hälfte aus [`ADR-0080`](0080-nahtform-pgconn-adapter-treiberhuelle.md)) mit einem Schlag ordnen | Änderung am **Messmechanismus** (Profil-Filter auf Block-Positionen statt Paketnamen) — nicht nötig für 80 %; §3.6-Gegenstand (Messfläche), also eigene Entscheidung mit eigener Beweislast; die Welle hat ihn ausdrücklich ausgeschlossen |
| D — Schnitt auf `internal/bootstrap` als „Hebel" (der Plan der Welle) | der Welle-Text nennt das Paket bereits; ein Slice, ggf. zwei | **widerlegt durch die Messung**: 308 der 336 ungedeckten Statements liegen in fünf Funktionen, die kein Test-Slice bewegt; der Schnitt legte die Welle auf den unbeweglichen Teil |
| E — Endstufe unter die erreichbare Decke senken | die Zahl wäre sofort erreichbar | die Endstufe steht ([`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md), [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) und ist **nicht** unerreichbar (Kontext (5)) — eine Senkung wäre eine §3.6-Schwellenänderung **ohne** den sie rechtfertigenden Grund |
| **F — Schnittmaß = Tail, Prozess-Rand als Puffer; Composition Root bleibt stehen (gewählt)** | erreicht 80 % real (Decke 83,6 %, Puffer 13–63 Statements); kein Produktionscode, keine Messänderung, keine Schwellenänderung; §3.6 unberührt; die Arbeit liegt dort, wo sie die Quote real bewegt | die Decke bei 83,6 % ist eine **benannte** Grenze, kein behobener Zustand; die Endstufe ist damit nicht mehr beliebig steigerbar, und der Rot-Beleg bei 85 % ist dauerhaft rot |

**Fazit:** F. A und D sind widerlegt (A durch das Design-Argument, D durch die
Messung), B greift zu kurz (76,3 %), C und E wären Entscheidungen ohne Anlass
— C nicht nötig, E durch die Messung verboten. F erreicht die Endstufe mit
benannter Grenze und ohne Eingriff in Gegenstand oder Schwelle.

## Konsequenzen

- Positiv: Die Welle bekommt ihr Schnittmaß — und es liegt nicht dort, wo ihr
  Text es vermutete. Der Planner kann schneiden, ohne die Messung zu
  wiederholen.
- Positiv: Kein Produktionscode, kein Messmechanismus, keine Schwelle wird
  berührt. §3.6 und die Out-of-Scope-Liste der Welle bleiben unverletzt.
- Positiv: Der Rot-Beleg bei `THRESHOLD=85`, den die Welle als
  Closure-Kriterium verlangt, ist vorab erklärt — er ist bei der heutigen
  Decke **notwendig** rot und damit ein echter Beleg, kein Zufall.
- Negativ mit Grenze: Die Endstufe 80 % ist die **vorletzte** Stufe, die
  dieser Gegenstand trägt. Darüber liegt die Decke bei 83,6 %; wer mehr will,
  braucht die code-granulare Partition (Option C) — eine eigene Entscheidung.
- Negativ mit Grenze: Der Puffer ist dünn, wenn Cluster A entfällt (13
  Statements, ≈0,7 pp), und er ist selbst nicht stabil: der Nenner der Messung
  ist eine Zustandsgröße, und der Zähler schwankt bereits über denselben
  Quelltext um 2 Statements (Kontext (2)). Die Welle sollte Cluster A nicht
  als „nice to have" führen.
- Negativ mit Grenze: Die 25 Statements des erzeugten `streamv1`-Pakets sind
  teils Proto-Runtime-Interna (`String`, `ProtoMessage`, `Descriptor`,
  `rawDescGZIP`, `init`) — ein Test, der *sie* abdeckt, prüft erzeugten Code
  gegen sich selbst. Sie sind in der realistischen Decke bereits abgezogen;
  der Slice, der Cluster C trägt, soll die Getter und die Client-/Handler-Stubs
  prüfen, nicht die Runtime-Interna.
- Folgepflicht (Planner-Zug): den Schnitt ziehen (vier Cluster, Cluster D ggf.
  an der Schichtgrenze geteilt); die §4-Zeile der Welle berichtigen, die
  `internal/bootstrap` als „Hebel" führt; die zwei verwandten Beobachtungen
  nachziehen — `adapter-unittest-verdeckt-bootstrap-luecke` (Träger-Frage
  beantwortet: netzlos nicht baubar, Träger ist der DB-/Integrations-Tier) und
  `db-gegenstand-enthaelt-netzlos-geprueften-code` (ihre Frage nach einer
  code-granularen Partition ist mit Festlegung 4 **offen gelassen**, nicht
  beantwortet — sie bekommt hier keine Kennung, sondern eine Fortsetzung).
- Folgepflicht (Implementer-Zug, je Test-Slice): **test-only**. Der Diff eines
  Slice dieser Welle enthält ausschließlich `*_test.go` (und, für Cluster A,
  keinen Produktionscode außerhalb der Testdatei); keine Zeile `internal/**`
  oder `cmd/**` außerhalb von Tests ändert sich, um Coverage zu gewinnen.
- Folgepflicht (Welle-Closure): Der Beleg ist der Gate-Lauf selbst — grün bei
  `THRESHOLD=80`, rot bei `THRESHOLD=85`, beides mit eigenem Grün-/Rot-Beleg
  dieses Laufs.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `make coverage-gate` (Stufe `coverage`, gepinntes Toolchain-Image) | **Gesamt-Coverage** der netzlos prüfbaren Fläche ≥ `THRESHOLD`; bei `THRESHOLD=80` grün, bei `THRESHOLD=85` rot | `make coverage-gate` (in `make gates`) |
| `git diff --name-only` des Slice-Commits | **Test-only:** jede geänderte Datei endet auf `_test.go`; keine Änderung an Produktionscode | — (Disziplin, Review-Gegenstand; kein Sensor) |
| Zählbasis-Nachweis des Laufs | Der Nenner der Messung (1903 zum Stand dieser ADR) und der Zähler (±2 über denselben Quelltext) sind **Zustandsgrößen**; die Zählbasis (Dedup über die Block-Position) ist dieselbe wie in `tools/harness/db-coverage.sh` | — (Beleg je Lauf) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger, sonst permanent:

**(a)** Ein Lauf zieht **doch** eine Naht am Composition Root (Konstruktor-
Injektion, `build`/`serve`-Zerlegung oder eine andere Form) — dann ist die
Aussage „9 von 198 netzlos erreichbar" ungültig und die Decke neu zu bemessen;
festgelegt ist hier nur, dass **diese Welle** sie nicht braucht.

**(b)** `internal/bootstrap` oder `cmd/pg-change-feed` ändert seinen
Dienst-Bedarf als **Paket** (etwa: der netzlos laufende Testrest zieht um) —
dann greift [`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Trigger (a): die namentliche Liste in `Dockerfile` und `DB_COVERAGE_PKGS` ist
nachzuziehen und die Rampe gegen ihren neuen Nenner zu prüfen.

**(c)** Die **code-granulare Partition** (Option C) wird entschieden — dann
fällt die Decke aus Kontext (5) weg und die Rampe ist neu zu bemessen.

**(d)** Die netzlos erreichbaren Präfixe der sechs Funktionen (Kontext (4a))
ändern ihre Form — etwa weil `pgxpool.New` lazy wird, ein Konstruktor seinen
`Ping` verliert, oder die CLI-Sondermodi umgebaut werden —, dann sind die
≈54 Puffer-Statements und Cluster A neu zu bemessen.

**(e)** Der Nenner wächst über die realistische Decke hinaus (neuer
dienstgebundener Code innerhalb des Gegenstands) — dann ist zu prüfen, ob die
Endstufe 80 % noch erreichbar ist; ist sie es nicht, ist **das** der Anlass
für Option C oder E, nicht ein Test-Slice.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-16 | Accepted — Anlass: `welle-20` steht ohne Schnittmaß („der Schnitt folgt der Messung"). Entscheidet: `Run` ist netzlos nicht prüfbar (9 von 198 Statements erreichbar, `Ping`-Riegel) und **keine** Naht wird gezogen; Schnittmaß ist der Tail (175–177) mit dem Prozess-Rand (≈54) als Puffer; Endstufe 80 % ist über dem unveränderten Gegenstand erreichbar (Decke 83,6 %); der Composition Root bleibt im Gegenstand; der vorgeschlagene Bootstrap-Smoke-Test ist netzlos nicht baubar. Bestätigt `ADR-0071`/`ADR-0077`/`ADR-0078` unverändert, ohne Supersedes | eigene Nachrechnung des `coverage`-Profils in zwei Läufen (Dedup über die Block-Position, Nenner 1903 = Gate-Wert, Zähler-Schwankung ±2), `docs/plan/planning/welle-20.md`, `docs/plan/planning/observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/observation.md`, `docs/plan/planning/observations/BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code/observation.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0082` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
