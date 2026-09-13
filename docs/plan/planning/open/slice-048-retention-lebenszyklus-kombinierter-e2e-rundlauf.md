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

**Verantwortlich:** —.

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

- [ ] Neuer, zusammenhängender Abschnitt in `run-integration-tests.sh`
      zeigt real die vollständige Blocker-Übergangs-Kette für **eine**
      Zeile: `cdc.retention_blockers` zeigt den zurückhängenden Consumer,
      dann — nach dessen Bestätigung — keinen Blocker mehr für diese
      Quelle, dann die reale Löschung der Zeile.
- [ ] `cdc_storage_bytes` bleibt über die gesamte Kette hinweg real
      abfragbar und numerisch (kein Fehler, kein `NULL`).
- [ ] `make gates` grün, `make test-integration` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` §„Aufbewahrung
      (Retention)" nennt den kombinierten Rundlauf als Testbeleg, falls
      das über den bereits dokumentierten Mechanismus hinausgeht
      (Implementer prüft und begründet im Plan-Nachzug).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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
  **Ausgang:** <bei Closure einzutragen>
- Der Retention-Hintergrundzug läuft alle 10 Sekunden
  (`retentionInterval`, `slice-044`) — der neue Abschnitt braucht
  ausreichend Poll-Zeit für sowohl den Blocker-Übergang als auch die
  anschließende Löschung, sonst wäre der Testlauf flaky. **Ausgang:** <bei
  Closure einzutragen>

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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
