# Slice slice-048: Retention-Lebenszyklus — kombinierter E2E-Rundlauf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — dieselbe Begründung wie bei `slice-047`: kein Mehr
über die eigene DoD hinaus, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht (Modul 6). Nutzerwunsch (2026-09-13): ein
kombinierter End-zu-Ende-Rundlauf über alle vier `welle-13`-Fähigkeiten, den
`slice-047` §1 bewusst als „ein eigener Vorgang" ausgeschlossen hatte.

**Bezug:** [`LH-FA-RET-002`](../../../../spec/lastenheft.md)…`006`.

**Berührte Spec-Stellen:** — (Testabdeckungs-Erweiterung auf bereits
bestehenden Fähigkeiten, keine neue Architektur-Sicht-Aussage).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein neuer, zusammenhängender Abschnitt in
`tools/harness/run-integration-tests.sh` spielt die vollständige
Retention-Kette real in einer Kette durch, statt sie wie bisher über vier
getrennte Slice-Belege zu prüfen: (1) eine Change-Zeile wird zurückdatiert
und ein Consumer bleibt zurück → `cdc.retention_blockers` zeigt diesen
Consumer real als aktuellen Blocker; (2) der Consumer bestätigt über die
Position hinweg → `cdc.retention_blockers` zeigt für diese Quelle keinen
Blocker mehr (leeres Ergebnis); (3) die Löschausführung entfernt die Zeile
real, ohne Neustart; (4) `cdc_storage_bytes` bleibt über die gesamte Kette
hinweg ein real abfragbarer, numerischer Wert. Das *Mehr* gegenüber den
vier Einzel-Slice-Belegen (`slice-043`…`046`): keiner von ihnen zeigt den
Blocker-*Übergang* (blockierend → nicht mehr blockierend) für **dieselbe**
Zeile in **derselben** Kette — jeder einzelne Beleg prüft nur seinen
eigenen Ausschnitt isoliert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`cdc_storage_bytes` als Beleg für eine Größenreduktion nach der
  Löschung** — `slice-046` §6 Risiko 1 hat bereits real begründet, dass
  PostgreSQL ohne `VACUUM`/`VACUUM FULL` die physische Tabellengröße nach
  einem `DELETE` nicht verkleinert (Bloat ist der korrekte, nicht
  irreführende Zustand). Ein Testfall, der eine Verkleinerung erwartet,
  wäre flaky (abhängig vom Autovacuum-Zeitpunkt) und würde genau die
  Fehllesart erzeugen, die `slice-046` bewusst vermied. Dieser Slice prüft
  nur, dass die Metrik **durchgehend abfragbar und numerisch** bleibt.
- **Neue Berechnungslogik oder neue SQL-Sichten** — reine Verkettung
  bereits bestehender, unveränderter Bausteine (`cdc.retention_blockers`,
  `cdc.metrics`, `RunRetentionUseCase`/`runRetentionCleanup`).
- **CLI-Diagnose-Erweiterung** — bereits `slice-047`; dieser Slice bleibt
  auf die reine SQL-Ebene (`psql` gegen die Compose-Postgres) beschränkt,
  keine Überschneidung mit `slice-047`s `docker exec`-CLI-Belegen.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Neuer, zusammenhängender Abschnitt in `run-integration-tests.sh`
      zeigt real die vollständige Blocker-Übergangs-Kette für **eine**
      Zeile: `cdc.retention_blockers` zeigt den zurückhängenden Consumer,
      dann — nach dessen Bestätigung — keinen Blocker mehr für diese
      Quelle, dann die reale Löschung der Zeile. Siehe Plan-Nachzug.
- [x] `cdc_storage_bytes` bleibt über die gesamte Kette hinweg real
      abfragbar und numerisch (kein Fehler, kein `NULL`).
- [x] `make gates` grün, `make test-integration` grün. Beide real
      ausgeführt (Ausgaben im Implementer-Bericht); `make test-integration`
      dreimal in Folge grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Siehe `docs/reviews/review-slice-048.md` — 0 HIGH, 1 MEDIUM (F-1,
      ohne Reviewer→Implementer-Rückgabe-Pfeil), keine Fixrunde nötig.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` §„Aufbewahrung
      (Retention)" nennt den kombinierten Rundlauf als Testbeleg, falls
      das über den bereits dokumentierten Mechanismus hinausgeht
      (Implementer prüft und begründet im Plan-Nachzug). **Geprüft:**
      nein — kein neuer Betriebs-Aspekt entsteht, siehe Plan-Nachzug
      Punkt 4.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — keine Beobachtung angefallen.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — beide entfallen.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | neuer, zusammenhängender Retention-Lebenszyklus-Abschnitt |
| `docs/user/benutzerhandbuch.md` | update, falls zutreffend | Testbeleg-Erwähnung, Implementer entscheidet |

### Plan-Nachzug (nach Implementierung)

Regeln dieser Sektion: Implementierungsentscheidungen, die über die Tabelle
oben hinausgehen — Platzierung des neuen Abschnitts, Isolation von Zeile und
Consumer, und ein realer Fund zum Sichtbarkeits-Vertrag von
`cdc.retention_blockers`, der die geplante Form von DoD-Punkt 1 verändert
hat.

**1. Platzierung: unmittelbar vor dem bestehenden „Zustand 1 — kein
Blocker"-Abschnitt, nicht am Ende des Skripts.** `cdc.retention_blockers`
trägt je Quelle höchstens eine Zeile (`DISTINCT ON`), solange irgendein
Consumer eine bestätigte Position gegen sie trägt (`INNER JOIN
cdc.consumer_position`, `tools/schema/schema.yaml`) — **unabhängig vom
Rückstandswert**. Sobald `CLI_CONSUMER`/`BACKLOG_CONSUMER` (beide später im
Skript) einmal bestätigt haben, trägt `cdc.retention_blockers` für
`src-mvp` bis zum Skriptende dauerhaft eine Zeile; eine „leere"
Rückkehr wäre dort nicht mehr real erreichbar. Der einzige Punkt im Skript,
an dem für `src-mvp` real noch keine bestätigte Consumer-Position
existiert, ist unmittelbar nach der Feed-Container-Bereitschaft und vor der
`CLI_CONSUMER`-Registrierung — exakt die Ausgangslage, die der bestehende
„Zustand 1"-Abschnitt bereits für sich in Anspruch nimmt (siehe dessen
Kommentar). Der neue Abschnitt läuft deshalb direkt davor, in derselben
Ausgangslage, und stellt sie am eigenen Ende real wieder her (Punkt 3).

**2. Realer Fund — „kein Blocker" ist ausschließlich Abwesenheit jeder
bestätigten Position, kein Rückstand von 0; DoD-Punkt 1 dementsprechend
umgesetzt.** Die View-Definition trägt keine Bedingung auf den berechneten
`backlog`-Wert — sie liefert exakt eine Zeile je Quelle, für die
`cdc.consumer_position` mindestens eine Zeile trägt, unabhängig davon, ob
dieser Consumer noch tatsächlich zurückhängt. Ein Consumer, der über die
Position der eigenen Test-Zeile hinweg bestätigt, verschwindet dadurch
**nicht** aus der View — er erscheint weiter, nur mit `backlog = 0`. Die
geplante Lesart von DoD-Punkt 1 ging von einer sich beim Aufholen selbst
leerenden Sicht aus; das entspricht nicht dem real ausgelieferten
Vertrag (`slice-045`, real durch drei `make test-integration`-Läufe hier
gegengeprüft, nicht nur gelesen). Der Abschnitt behandelt diesen Unterschied
nicht als Bug — die View macht exakt das sichtbar, was
`RunRetentionUseCase`/`ConsumerStatePort.Positions` bereits als
Abwesenheits-Lesart etabliert (`slice-045` Plan-Nachzug Punkt 3): ein
Consumer ohne jede bestätigte Position blockiert nicht, und nur die
Abwesenheit jeder Position macht die Quelle „unblockiert" sichtbar. Der
neue Abschnitt bildet deshalb nach dem Sichtbarkeits-Nachweis
(`cdc.retention_blockers` zeigt real `cli-e2e-lifecycle-consumer` als
aktuellen Blocker) und der realen Bestätigung über die Position der
Test-Zeile hinweg einen expliziten dritten Schritt: die bestätigte Position
des eigens für diesen Test registrierten Consumers wird real entfernt
(`DELETE FROM cdc.consumer_position WHERE consumer_id = …`) — kein
CLI-Unterbefehl dafür existiert, dieselbe Klasse direkter administrativer
SQL-Schreibzugriff, die das Skript bereits für `cdc.process_heartbeat` und
`cdc.transaction.committed_at` nutzt. Wichtig: diese Entfernung **löst**
die reale Löschung nicht aus — sie ist bereits durch die Bestätigung
selbst freigegeben (`Positions()` liefert ab diesem Zeitpunkt nur noch
Positionen ≥ der Zeilen-Position, unabhängig davon, ob die Zeile in
`cdc.consumer_position` später entfernt wird); die Entfernung dient
ausschließlich der Sichtbarkeits-Prüfung selbst (DoD-Punkt 1, „keinen
Blocker mehr … leeres Ergebnis") und der Wiederherstellung der
Ausgangslage für den nachfolgenden „Zustand 1"-Abschnitt (Punkt 3). Beide
Zwischenzustände sind real per SQL geprüft, nicht nur die Endzustände:
`blocker_before` (Consumer sichtbar) und `blocker_after` (keine Zeile mehr)
in `run-integration-tests.sh`.

**3. Isolation — eigene Zeile (id=210 auf `feed_mvp_full`), eigener
Consumer (`cli-e2e-lifecycle-consumer`), reale Wiederherstellung der
Ausgangslage danach (löst §6 Risiko 1 auf).** `id=210` kollidiert mit
keinem bestehenden Wertebereich auf `feed_mvp_full` (id=1/7 aus den
Go-Integrationstests, 90/91 Lasttest-Beleg, 95/96 `CLI_CONSUMER`, 120/121
`BACKLOG_CONSUMER`, 200/201 `RetentionOld`/`RetentionYoung`).
`cli-e2e-lifecycle-consumer` ist von `CLI_CONSUMER`/`BACKLOG_CONSUMER`
namentlich getrennt. Der zurückdatierte Timestamp
(`current_timestamp - interval '25 hours'`, dieselbe Technik wie beim
bestehenden Retention-Beleg) und die Bestätigung auf die real gemessene
älteste Position der Quelle (`min(commit_position) FROM cdc.transaction
WHERE source_id = 'src-mvp'`, real vorhanden aus den vorangegangenen
Go-Integrationstests/dem Lasttest-Beleg) berühren keine der später
laufenden Sektionen. Nach der realen Entfernung der Consumer-Position
(Punkt 2) trägt `cdc.consumer_position` für `src-mvp` wieder exakt null
Zeilen — dieselbe Ausgangslage, die der nachfolgende „Zustand 1"-Abschnitt
voraussetzt; dessen eigener Kommentar wurde entsprechend ergänzt. Real
durch drei aufeinanderfolgende `make test-integration`-Läufe bestätigt,
inklusive des unveränderten Bestehens des nachfolgenden „Zustand 1"- und
des bereits bestehenden `RetentionOld`/`RetentionYoung`-Belegs.

**4. Kein Doku-Update — reiner Testbeleg auf bereits vollständig
dokumentiertem Mechanismus.** `docs/user/benutzerhandbuch.md` dokumentiert
`cdc.retention_blockers` (§„Blockierende Consumer erkennen", `slice-045`),
`cdc_storage_bytes` (§„Metriken lesen", `slice-046`) und die
`diagnose`-Erweiterung (§„Diagnose ausführen", `slice-047`) bereits
vollständig — dieser Slice führt keine neue Betriebs-Aussage ein, sondern
verkettet drei bereits real gelieferte und dokumentierte Fähigkeiten in
einem zusätzlichen Testbeleg. Kein neuer Abschnitt, keine neue
Changelog-Zeile.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-047` liegt in `done/` (WIP-Limit
1 je Implementer — beide Slices ändern `run-integration-tests.sh`, daher
seriell statt parallel), `Verantwortlich:` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  der kombinierte Abschnitt eine eigene, isolierte Feed-Container-Instanz
  braucht (wie bei Fehlerklassen-Tests), statt im bestehenden Lauf
  mitzulaufen, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der neue kombinierte Abschnitt teilt sich den Compose-Lauf mit dem
  bestehenden `slice-044`-Retention-Beleg (`RetentionOld`/`RetentionYoung`,
  dieselben Consumer `CLI_CONSUMER`/`BACKLOG_CONSUMER`) — eine eigene,
  isolierte Zeile/Position ist nötig, sonst könnten sich beide Abschnitte
  gegenseitig verfälschen (`BEO-PGC/test-isolation-geteilter-zustand`).
  **Ausgang: entfallen.** Eigene, isolierte Zeile (id=210 auf
  `feed_mvp_full`, kollidiert mit keinem bestehenden Wertebereich) und
  eigener, neu registrierter Consumer (`cli-e2e-lifecycle-consumer`,
  namentlich getrennt von `CLI_CONSUMER`/`BACKLOG_CONSUMER`), zusätzlich
  real am eigenen Ende wieder entfernt, um die Ausgangslage des
  nachfolgenden „Zustand 1"-Abschnitts unverändert zu halten
  (Plan-Nachzug Punkt 3) — real durch drei `make test-integration`-Läufe
  bestätigt, kein gegenseitiges Verfälschen beobachtet.
- Der Retention-Hintergrundzug läuft alle 10 Sekunden
  (`retentionInterval`, `slice-044`) — der neue Abschnitt braucht
  ausreichend Poll-Zeit für sowohl den Blocker-Übergang als auch die
  anschließende Löschung, sonst wäre der Testlauf flaky. **Ausgang:
  entfallen.** Dieselbe Technik wie beim bestehenden Retention-Beleg: eine
  feste Wartezeit über mehr als zwei Lösch-Takte (25s) belegt die
  Abwesenheit einer verfrühten Löschung, ein anschließender Poll (bis zu
  60×1s) belegt die reale Löschung nach Freigabe. Real dreimal in Folge
  grün, keine Flakiness beobachtet.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Die Isolations-Disziplin aus §6 (eigene
  Zeile, eigener Consumer, reale Rücknahme am eigenen Ende) ließ den neuen
  Abschnitt direkt vor den bestehenden „Zustand 1"-Beleg setzen, ohne dessen
  Voraussetzung (keine bestätigte Consumer-Position für `src-mvp`) zu
  verletzen — real durch drei aufeinanderfolgende `make
  test-integration`-Läufe bestätigt, alle bestehenden Belege (Black-Box-CLI-
  Rundlauf, Rückstands-Beleg, CLI-Diagnose, bestehender Retention-Beleg)
  liefen dabei unverändert grün mit.
- **Was ging anders als geplant:** Die geplante Lesart von DoD-Punkt 1
  („nach Bestätigung zeigt die View keinen Blocker mehr") ging von einem
  sich selbst leerenden View-Ergebnis beim Aufholen aus. Real trägt
  `cdc.retention_blockers` je Quelle weiterhin eine Zeile, solange
  irgendein Consumer eine bestätigte Position gegen sie trägt — unabhängig
  vom Rückstandswert. „Kein Blocker" ist ausschließlich die Abwesenheit
  jeder bestätigten Position (siehe Plan-Nachzug Punkt 2), nicht ein
  Rückstand von 0. Kein Substanz-Mangel des bestehenden Vertrags
  (`slice-045`) — reiner Unterschied zwischen der ursprünglich angenommenen
  und der real ausgelieferten Semantik, gelöst durch einen zusätzlichen,
  transparent dokumentierten dritten Schritt (reale Entfernung der
  Test-Consumer-Position nach der Bestätigung).
- **Steering-Loop-Eintrag:** *(kein Eintrag verkörpert — der Normalfall.)*
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen. Beide vorab gesichteten Treffer (`BEO-PGC/test-isolation-
  geteilter-zustand`, `BEO-PGC/test-runner-stiller-ausschluss`) bleiben bei
  1× — der erste ist mit diesem Slice *entfallen* (§6), nicht eingetreten;
  der zweite betrifft neue Go-Testfunktionen und war für diesen rein
  SQL/Shell-basierten Slice nicht einschlägig (§1 Ausschluss).
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide *entfallen* — siehe §6.
- **Drei Paarungen:** Wellenlos — Prüfung nach dem `git mv` nach `done/`
  durch die nächste Rolle (kein `liegt in`-Feld in diesem Slice, kein
  neuer Folge-Slice, keine neue Registerzeile — nur die Register-Existenz
  der beiden zitierten, bereits bestehenden Beobachtungen bleibt zu
  bestätigen).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/test-isolation-geteilter-zustand` (1×, direkt
einschlägig für diesen Slice, siehe §6), `BEO-PGC/test-runner-stiller-ausschluss`
(1×, jeder neue Abschnitt muss real mitlaufen). Keiner erreicht mit
diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
