# Slice slice-038: CLI-Diagnose-Befehl

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-12 — dritter und letzter Slice, schließt
`BEO-PGC/verwaltung-keine-sql-administration` ab (siehe `welle-12.md` §1/§3).

**Bezug:**
[`LH-FA-SST-003`](../../../../spec/lastenheft.md) (CLI),
[`LH-FA-ADM-002`](../../../../spec/lastenheft.md) (Betriebsstatus),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbare Fehlerzustände),
[`LH-FA-ADM-004`](../../../../spec/lastenheft.md) (messbarer CDC-Abstand),
[`LH-FA-ADM-005`](../../../../spec/lastenheft.md) (sichtbarer
Verarbeitungsrückstand).

**Berührte Spec-Stellen:** — (reine CLI-Ergänzung auf bereits bestehenden
SQL-Lese-Views, keine Verhaltensänderung an einer Architektur-Sicht-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

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

**Ziel:** Ein neuer CLI-Diagnose-Befehl (`cmd/pg-change-feed/main.go`, Muster
identisch zum bestehenden `--healthcheck`) liest über `cfg.ReaderDSN`
(`ADR-0047`) die bereits bestehenden, `cdc_reader`-gegrantenen SQL-Lese-Views
`cdc.heartbeat` (`tools/schema/nacharbeit-heartbeat.sql`) und `cdc.metrics`
(`tools/schema/nacharbeit-observability.sql`) und gibt eine
menschenlesbare Zusammenfassung aus, die `LH-FA-SST-003`s Boundary
("deckt mindestens die Status-/Diagnoseabfragen ab, die `LH-FA-ADM-002`
… `005` nennen") real erfüllt:
- **`LH-FA-ADM-002`** (Betriebsstatus) und **`LH-FA-ADM-003`** (sichtbare
  Fehlerzustände) — aus `cdc.heartbeat` (`age_seconds`, `error_class`).
- **`LH-FA-ADM-004`** (messbarer CDC-Abstand) — aus `cdc.metrics`,
  Zeile `cdc_capture_lag`.
- **`LH-FA-ADM-005`** (sichtbarer Verarbeitungsrückstand) — aus
  `cdc.metrics`, Zeilen `cdc_consumer_lag{consumer}` (je Consumer).

Kein neuer Domänentyp, kein neuer Port: alle vier Werte stehen bereits als
Lese-Views bereit (`LH-FA-SST-002`, seit `welle-11`/früher); dieser Slice
liefert ausschließlich den fehlenden CLI-Zugriffsweg darauf.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue SQL-Views oder -Spalten** — `cdc.heartbeat`/`cdc.metrics` decken
  alle vier `LH-FA-ADM-00N`-Signale bereits vollständig ab (siehe oben);
  dieser Slice ist reiner Lesezugriffsweg, keine Datenmodell-Änderung
  (Schicht-Abgrenzung: CLI-Schicht, nicht Schema-Schicht).
- **Maschinenlesbares Ausgabeformat (JSON o. ä.)** — `LH-FA-SST-003`
  verlangt nur, dass die Abfragen "über eine CLI erfolgen" können, kein
  bestimmtes Format; `cdc.metrics` selbst ist bereits das
  maschinenlesbare Long-Format (`LH-FA-SST-004`) und über SQL direkt
  erreichbar (`LH-FA-SST-002`) — ein zweites, CLI-eigenes
  Maschinenformat wäre eine Doppelung ohne neuen Anforderungs-Träger.
  Bleibt Bestand: der bestehende SQL-Zugriffsweg auf dieselben Views.
- **Schwellenwert-Interpretation (healthy/degraded/unhealthy,
  `SPEC-007`)** — dieselbe bewusste Nicht-Entscheidung wie bei
  `cdc.heartbeat`/`cdc.metrics` selbst (siehe deren Kommentare): die
  Views liefern Rohwerte, keine Klassifikation; der neue Befehl gibt sie
  unverändert weiter, trifft keine neue Schwellenwert-Entscheidung, die
  ein ADR bräuchte (`AGENTS.md` §3.6).
- **Ein Datenbankzugriff außerhalb `cdc_reader`s Grant-Fläche**
  (z. B. `pg_stat_replication`, Relationsgrößen) — dieselbe bewusste
  Nicht-Abdeckung wie in `nacharbeit-observability.sql` dokumentiert
  (`cdc_wal_retention_bytes`, `cdc_storage_bytes` u. a.); ein anderer
  Vorgang, der die Least-Privilege-Fläche erweitern müsste.

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

- [ ] Neuer CLI-Befehl (`cmd/pg-change-feed/main.go`, z. B. `diagnose` —
      Implementer entscheidet den exakten Namen und begründet ihn im
      Plan-Nachzug) liest `cdc.heartbeat` + `cdc.metrics` über
      `cfg.ReaderDSN` und gibt `age_seconds`/`error_class`
      (`LH-FA-ADM-002`/`003`), `cdc_capture_lag`
      (`LH-FA-ADM-004`) und `cdc_consumer_lag{consumer}` je Consumer
      (`LH-FA-ADM-005`) menschenlesbar aus.
- [ ] `LH-FA-SST-003` real erfüllt: ein Integrationstest (Ergänzung in
      `tools/harness/run-integration-tests.sh`, analog zum bestehenden
      `--healthcheck`-Testabschnitt) ruft den neuen Befehl gegen den
      laufenden Compose-Feed-Container per `docker exec` auf und
      bestätigt, dass alle vier Signale in der Ausgabe erscheinen —
      sowohl im Normalbetrieb als auch (mindestens für ADM-003) in einem
      erkennbar von Normalbetrieb unterscheidbaren Fehlerzustand
      (Boundary-Kriterium von `LH-FA-ADM-002`/`003`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `README.md`/`harness/README.md`, falls dort die
      vorhandenen CLI-Befehle aufgezählt sind (Implementer prüft und
      begründet im Plan-Nachzug, ob eine Stelle existiert, die den neuen
      Befehl nennen muss).
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
| `cmd/pg-change-feed/main.go` | update | neuer Subcommand, Muster identisch zu `--healthcheck`/`register-consumer` |
| `internal/bootstrap/wiring.go` | update | neue Funktion (Muster `Healthcheck`) — öffnet kurzlebige `ReaderDSN`-Verbindung, liest `cdc.heartbeat` + `cdc.metrics`, formatiert die Ausgabe |
| `tools/harness/run-integration-tests.sh` | update | realer E2E-Beleg des neuen Befehls gegen den laufenden Feed-Container |
| `README.md` oder `harness/README.md` | update, falls zutreffend | Implementer prüft, ob eine Stelle die vorhandenen CLI-Befehle aufzählt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-036` und `slice-037` liegen in
`done/` (unabhängig davon laut `welle-12` §5 auch parallel startbar,
tatsächlich aber sequenziell nach diesem Slice-Zuschnitt), `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  der neue Befehl plus sein E2E-Beleg zusammen mehr als drei Liefer-Punkte
  oder mehr als zwei Schichten in einer Review-Sitzung nicht mehr prüfbar
  machen, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün (inkl. des neuen CLI-Diagnose-Belegs) **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein realer Fehlerzustand (`LH-FA-ADM-003`s Boundary: „unterscheidbar von
  normalem Betrieb") im laufenden Compose-Stack künstlich zu erzeugen,
  könnte schwieriger sein als angenommen (bisherige Fehlerklassen-Tests,
  z. B. `TestMVPSchemaChangeIncompatibleTypeChange`, brauchen einen
  eigenen, isolierten Feed-Container-Lauf) — **Ausgang:** <bei Closure
  einzutragen>
- `cdc.metrics`s `cdc_consumer_lag`-Zeilen existieren nur für Consumer mit
  mindestens einer bestätigten Position (`WHERE cs.acknowledged_position
  IS NOT NULL`); ein frisch registrierter, noch nie bestätigender
  Consumer erscheint dort nicht — die CLI-Ausgabe könnte das
  fälschlich als „kein Rückstand" statt „noch nie gemessen" lesen lassen.
  **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/verwaltung-keine-sql-administration` (0×,
benannt nicht gezählt — dieser Slice liefert den dritten und letzten
Baustein der Auflösung, siehe `welle-12` §3 Closure-Trigger) und
`BEO-PGC/github-actions-unverifizierbar-lokal` (1×, thematisch nicht
berührt von diesem Slice — reine CDC-Fähigkeit, kein CI/CD-Bezug). Keiner
der Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
