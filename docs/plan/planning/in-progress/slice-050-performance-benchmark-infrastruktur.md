# Slice slice-050: Performance-Benchmark-Infrastruktur

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-14 — zweiter Slice, unabhängig von `slice-049`.

**Bezug:** [LH-QA-PER-001](../../../../spec/lastenheft.md),
[LH-QA-PER-002](../../../../spec/lastenheft.md),
[LH-QA-PER-003](../../../../spec/lastenheft.md),
[ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(Muster, kein Gate).

**Berührte Spec-Stellen:** — (Mess-Skript-Ergänzung, keine neue
Architektur-Sicht-Aussage).

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

**Ziel:** Drei eigenständige Bench-Skripte nach dem Stil von
`/Development/d-check/tools/bench-fixture.sh`
([ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(b)), gebündelt hinter einem gemeinsamen `make bench`-Target (analog
`/Development/d-check/Makefile`), das explizit **kein Gate** ist:
(1) `LH-QA-PER-001` Quell-Impact mit/ohne CDC (dieselbe Schreiblast auf
die Quelltabelle, einmal mit aktivem Replication-Slot/Capture-Prozess,
einmal ohne), (2) `LH-QA-PER-002` Skalierung über die drei
[SPEC-014](../../../../spec/pflichtenheft.md)-Lastenstufen (klein ≤10/s,
mittel 100/s×30min, groß 1000/s×60min) als feste Eingabeparameter,
(3) `LH-QA-PER-003` Batch- vs. Einzelabruf-Effizienz beim Lesen über
`cdc.changes`. Jedes Skript dokumentiert Aufwand/Ergebnis (Report-Datei
oder stdout), ohne einen Pass/Fail-Schwellenwert durchzusetzen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Aufnahme in `make gates`/`fullbuild` als Pass/Fail-Bedingung** —
  `ADR-0054` §(b) legt Benchmark ausdrücklich als dokumentierten Beleg
  fest, nicht als Gate; andere Disziplin als das Coverage-Gate.
- **Neue Lastenstufen-Definition** — `SPEC-014` legt die drei Stufen
  bereits fest; dieser Slice übernimmt sie als Eingabeparameter, ändert
  sie nicht.
- **Test-Coverage-Gate** — `slice-049`; andere Schicht (Build-/Gate-
  Infrastruktur statt eigenständiger Mess-Skripte).
- **`LH-QA-PER-004`-Ausbau** (Commit→CDC-Latenz über den bestehenden
  Lasttest-Beleg hinaus) — `ADR-0054`s Re-Evaluierungs-Trigger (c) benennt
  einen eigenen Folge-Vorgang, falls das je nötig wird; bleibt hier
  unverändert Bestand.

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

- [x] `LH-QA-PER-001` real erfüllt: Bench-Skript zeigt real einen
      messbaren Unterschied (oder dessen Abwesenheit) zwischen
      Schreiblast mit und ohne aktivem Replication-Slot/Capture-Prozess,
      dokumentiertes Ergebnis. Siehe §3 Plan-Nachzug — real gemessen:
      1000 Schreibtransaktionen ohne CDC 3812 ms, mit CDC 7104 ms
      (+86,4 %).
- [x] `LH-QA-PER-002` real erfüllt: Bench-Skript durchläuft real alle
      drei `SPEC-014`-Lastenstufen, dokumentiertes Ergebnis je Stufe.
      Siehe §3 Plan-Nachzug — alle drei Stufen (10/100/1000 pro
      Sekunde) real mit Ziel-Rate durchlaufen, `cdc_capture_lag` je
      Stufe gelesen.
- [x] `LH-QA-PER-003` real erfüllt: Bench-Skript vergleicht real
      Batch- vs. Einzelabruf beim Lesen über `cdc.changes`,
      dokumentiertes Ergebnis. Siehe §3 Plan-Nachzug — real gemessen:
      200 Zeilen Batch 72 ms vs. Einzelabruf 12874 ms (178,8×).
- [x] `make bench` startet alle drei Skripte, kein Gate (Aufnahme in
      `make gates` explizit unterlassen), `harness/README.md`
      §Werkzeuge trägt die neue Zeile.
- [x] `make gates` grün (unverändert, da kein Gate hinzukommt).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Siehe `docs/reviews/review-slice-050.md` — 0 HIGH/MEDIUM/LOW, 2 INFO,
      keine Fixrunde nötig (DoD-Checkbox-Nachzug ohne Fixrunde,
      `.harness/skills/reviewer.md`).
- [x] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — neues Verzeichnis `BEO-PGC/schema-rollout-braucht-compose-init/`.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). `welle-14` ist offen; an die Welle-14-Closure delegiert.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/bench-source-impact.sh` | neu | `LH-QA-PER-001`-Beleg |
| `tools/bench-scaling.sh` | neu | `LH-QA-PER-002`-Beleg |
| `tools/bench-batch-vs-single.sh` | neu | `LH-QA-PER-003`-Beleg |
| `Makefile` (bzw. `harness/mk/*.mk`) | update | `bench`-Target, kein Gate |
| `harness/README.md` | update | §Werkzeuge-Zeile für `make bench` |

### Plan-Nachzug (nach Code)

**Vierte Datei ergänzt, ungeplant:** `tools/bench-lib.sh` — gemeinsame
Umgebungs-Bausteine (PostgreSQL-/Feed-Container über `docker
network`/`docker run`, Schema-Rollout, Cleanup), von allen drei
Bench-Skripten gequellt. `ADR-0054` §(b) verlangt „je Beleg ein eigenes
Bench-Skript" — das bleibt gewahrt: geteilt ist nur der Umgebungsaufbau,
nicht die Messung selbst (jedes der drei Skripte bleibt für sich lauffähig,
lesbar und änderbar). Ohne die gemeinsame Datei hätte jedes der drei
Skripte denselben ca. 40-zeiligen Umgebungs-Aufbau dupliziert — genau die
Kopplung, die `ADR-0054` mit „ein Fix an einem Beleg riskiert, die anderen
zwei mitzubrechen" bei der **verworfenen** Option B (ein kombiniertes
Skript) meinte, hier aber auf den Aufbau bezogen, nicht auf die Messung.
Die Umgebung ist bewusst **nicht** `compose.yaml` (das die
Integrationstests nutzen): eigene Netz-/Container-Namen
(`pgc-bench-*`), damit ein `make bench`-Lauf nicht mit einem parallel
laufenden `make test-integration` um dieselben Namen konkurriert.

**Realer Fallstrick beim Aufbau (nicht in `.dockerignore`/Alpine-Klasse,
siehe unten):** Der erste `make schema-rollout`-Versuch gegen die
eigenständige Bench-Umgebung scheiterte real mit `POST_EXECUTE_DRIFT`
(Exit 5, `relation "cdc.source_table" does not exist"`) beim Anlegen der
ersten View. Ursache: `compose.yaml` mountet
`tools/schema/compose-init/01-cdc-schema.sql` als PostgreSQL-Init-Skript
(`CREATE SCHEMA IF NOT EXISTS cdc; ALTER ROLE postgres IN DATABASE cdc SET
search_path = cdc;`) — ohne diesen Mount landet die unqualifizierte
Tabellen-DDL des d-migrate-Rollouts in `public`, während die generierten
Views explizit `cdc.<table>` referenzieren. Behoben, indem
`bench::start_postgres` denselben `docker-entrypoint-initdb.d`-Mount
setzt wie `compose.yaml`. Neue Beobachtung dokumentiert:
[`BEO-PGC/schema-rollout-braucht-compose-init`](../observations/BEO-PGC/schema-rollout-braucht-compose-init/observation.md)
(1×, unter der Schwelle) — eine andere Fehlerklasse als die aus
`slice-049` bekannte `.dockerignore`/Alpine-`bash`-Beobachtung
(`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`): Diese
Skripte fügen **keine** neue Docker-Multi-Stage-Stufe hinzu (sie laufen
direkt per `bash` gegen bereits gebaute/geladene Images, wie
`tools/harness/run-integration-tests.sh`), der dort benannte Fallstrick
(`.dockerignore`-Ausnahme, fehlendes `bash` in der Alpine-Basis) ist
deshalb hier **nicht einschlägig** — die zugehörige `state.md` bleibt
unverändert bei 1× (weiter offen).

**Risiko 2 (Slice-Plan §6) — Lösung:** `tools/bench-scaling.sh`
unterscheidet zwei Modi: **Default** fährt stark verkürzte, aber reale
Dauern je Stufe (`klein`/`mittel`/`groß`: 10 s/15 s/15 s, je per Env-Var
override- bar) bei **unveränderter Ziel-Rate** (10/100/1.000 Änderungen
pro Sekunde, real über eine `INSERT … generate_series`-Anweisung je
Sekunde erzeugt — bei 1.000/s wären 1.000 einzelne
`docker exec`-Aufrufe pro Sekunde nicht durchhaltbar gewesen, die
serverseitige Mengen-Anweisung bleibt ein reales Schreib-Volumen in der
Zielrate). **`--full`** fährt die tatsächlichen `SPEC-014`-Dauern
(`mittel` 1.800 s, `groß` 3.600 s). Für die `klein`-Stufe legt
`SPEC-014` keine Dauer fest (nur die Rate ≤10/s) — die hier gewählte
volle Dauer (60 s) ist eine dokumentierte Annahme dieses Bench-Skripts,
keine Schärfung von `SPEC-014` selbst. Real durchlaufen (Default-Modus,
`make bench`): `klein` 100 Zeilen/10 s (~10,0/s), `mittel` 1.500
Zeilen/15 s (~100,0/s), `groß` 15.000 Zeilen/15 s (~1.000,0/s) —
`cdc_capture_lag` blieb in allen drei Stufen nahe 1 s.

**Risiko 1 (Streuung) — Einordnung vorweggenommen für §6:** Anders als
d-checks `bench-fixture.sh` (eine Kennzahl gegen eine feste
< 5 s-Schwelle, N=3-Läufe + Median nötig, weil ein einzelner Ausreißer
das Gate fälschlich rot färben könnte) tragen diese drei Skripte **keine**
Schwelle — Aufwand/Ergebnis wird dokumentiert, nicht durchgesetzt
(`ADR-0054` §(b)). `bench-source-impact.sh` vergleicht zusätzlich beide
Phasen **innerhalb desselben Laufs** gegen dieselbe Postgres-Instanz
unmittelbar nacheinander — Host-seitige Varianz (Docker-Overhead,
CPU-Kontention) wirkt auf beide Phasen ähnlich und wird im Differenzwert
weitgehend herausgekürzt. Details siehe §6.

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-14` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-049`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Drei
  eigenständige Bench-Skripte plus Makefile-Wiring nähern sich der
  Drei-Liefer-Punkte-Grenze — zeigt sich, dass eines der drei Skripte
  selbst schon einen eigenen Liefer-Punkt umfangreicher wird (z. B.
  `LH-QA-PER-002`s drei Lastenstufen brauchen eine eigene
  Fixture-Erzeugung je Stufe statt eines Parameters), gehört das zurück
  zur Zerlegung (ein Skript je Slice statt aller drei zusammen).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make bench` liefert real alle drei Belege **und**
`make gates` unverändert grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Benchmark-Ergebnisse könnten in einer geteilten/virtualisierten
  Docker-Umgebung (kein dediziertes Hardware-Budget) real streuen —
  dasselbe Problem, das d-checks `bench:`-Target über N=3-Läufe und
  Median statt Einzelmessung adressiert. **Ausgang: entfallen.**
  Begründung: Anders als d-checks Ein-Schwellen-Gate (< 5 s, ein
  Ausreißer kann fälschlich rot werden) tragen alle drei Skripte hier
  keine Pass/Fail-Schwelle (`ADR-0054` §(b)) — Streuung verfälscht kein
  Urteil, nur die absolute Zahl. `bench-source-impact.sh` vergleicht
  zudem beide Phasen innerhalb desselben Laufs gegen dieselbe Instanz
  unmittelbar nacheinander, was Host-seitige Varianz in der
  Differenzmessung weitgehend kürzt (siehe §3 Plan-Nachzug). Wer höhere
  statistische Sicherheit braucht, kann jedes Skript mehrfach aufrufen
  (kein technischer Hinderungsgrund) — das ist bewusst nicht in den
  Skripten erzwungen, weil es den Aufbau-/Laufzeit-Aufwand für einen
  reinen Dokumentations-Beleg unnötig verdreifachen würde.
- `LH-QA-PER-002`s „groß"-Lastenstufe (1.000/s × 60 Minuten,
  `SPEC-014`) könnte `make bench` für einen schnellen, wiederholten
  Implementer-/Reviewer-Lauf unpraktikabel lang machen. **Ausgang:
  eingetreten.** Gelöst innerhalb dieses Slices (kein Carveout, kein
  Folge-Slice nötig): `tools/bench-scaling.sh` trägt einen
  Default-Modus mit real stark verkürzten Dauern bei unveränderter
  Ziel-Rate sowie ein `--full`-Flag für die tatsächlichen
  `SPEC-014`-Dauern — siehe §3 Plan-Nachzug für die volle Begründung und
  die real gemessenen Default-Lauf-Ergebnisse.

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

- **Was hat funktioniert:** Das Kopiervorbild `/Development/d-check/Makefile`
  Zeile 84 + `/Development/d-check/tools/bench-fixture.sh` trug den
  Grundriss (Fixture/Umgebung aufbauen, real messen, dokumentiertes
  Ergebnis auf stdout), angepasst um die `ADR-0054`-Vorgabe „drei
  eigenständige Skripte statt einer Kennzahl gegen eine Schwelle". Die
  gemeinsame `tools/bench-lib.sh` hielt die drei Skripte trotz geteiltem
  Umgebungsaufbau unabhängig lauffähig — jedes lief einzeln vor dem
  gebündelten `make bench` real durch. Alle drei Skripte liefern
  plausible, klar interpretierbare reale Ergebnisse (CDC-Schreib-Overhead
  ~86 %, Batch/Einzelabruf-Faktor ~179×, alle drei Lastenstufen mit
  `cdc_capture_lag` nahe 1 s).
- **Was ging anders als geplant:** Ein realer Fallstrick beim Aufbau der
  eigenständigen (von `compose.yaml` unabhängigen) Bench-Umgebung: der
  fehlende `tools/schema/compose-init`-Mount ließ den ersten
  `make schema-rollout`-Versuch mit `POST_EXECUTE_DRIFT` (Exit 5)
  scheitern, weil die generierten Views explizit `cdc.<table>`
  referenzieren, während die unqualifizierte Tabellen-DDL ohne den
  `search_path`-Init in `public` gelandet wäre — siehe §3 Plan-Nachzug
  und die neue Beobachtung unten. Eine andere, in `slice-049` bereits
  bekannte Fallstrick-Klasse (`.dockerignore`/Alpine-`bash`,
  `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`) trat **nicht**
  erneut auf, weil diese Skripte keine neue Docker-Multi-Stage-Stufe
  einführen, sondern — wie `tools/harness/run-integration-tests.sh` —
  direkt per `bash` gegen bereits gebaute/geladene Images laufen; diese
  `state.md` bleibt unverändert bei 1×.
- **Steering-Loop-Eintrag:** keiner — dieser Slice liefert die in
  `ADR-0054` §(b) bereits entschiedene Bench-Infrastruktur, ohne einen
  Guide/Sensor über dieses Repo hinaus zu schärfen. Der Eintrag ist
  gezählt (Beobachtung unten), nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/schema-rollout-braucht-compose-init/`
  neu angelegt, Beleg `evidence/slice-050.md` — Zähler steht bei 1×
  (unter der Schwelle).
- **Folge-Slices:** keine.
- **Risiken aus §6:** eines entfallen (Streuung — keine Schwelle
  betroffen, Differenzmessung innerhalb desselben Laufs), eines
  eingetreten und innerhalb dieses Slices gelöst (groß-Stufe-Dauer →
  Default-/`--full`-Modus) — siehe §6. Der Verifier reproduzierte die
  Bench-Läufe eigenständig und maß real spürbar andere Werte (91,4 %
  statt 86,4 % CDC-Overhead; Faktor 128,2× statt 178,8×) — das bestätigt
  die reale Streuung unabhängig, ändert aber den Ausgang *entfallen*
  nicht, da kein Skript einen Pass/Fail-Schwellenwert trägt
  (`verify-slice-050.md`).
- **Drei Paarungen:** entfällt hier — dieser Slice gehört zu `welle-14`
  (offen); die Paarungen prüft die Welle-14-Closure.

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
Keine Treffer für `PGC` zu Performance/Benchmark.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
