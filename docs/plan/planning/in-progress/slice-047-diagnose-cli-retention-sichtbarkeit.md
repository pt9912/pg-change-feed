# Slice slice-047: Diagnose-CLI-Erweiterung für Retention-Sichtbarkeit

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Eröffnungs-Recherche (Fork) zur ursprünglich in der
Roadmap vorgemerkten Welle „E2E-Abdeckung — Retention" fand real: Die
SQL-View-Black-Box-Ebene ist bereits vollständig geliefert
(`TestMVPRetentionBlockersViewShowsFurthestBehindConsumer`,
`TestMVPMetricsCarriesStorageBytes`, beide `slice-045`/`046`, beide reine
SQL-Lesetests ohne Go-Domain-Import — exakt die Testebene, die `welle-11`
für `cdc.active_tables`/`cdc.consumer_status` erst nachliefern musste). Die
einzige real verbleibende Lücke ist die fehlende `docker exec`-CLI-
Sichtbarkeit (`internal/bootstrap/wiring.go`s `Diagnose`-Funktion zeigt
nichts zu Retention) — ein einzelner Slice ohne Closure-Bedingung jenseits
seiner eigenen DoD, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**Bezug:** [`LH-FA-SST-003`](../../../../spec/lastenheft.md) (CLI-Diagnose),
[`LH-FA-RET-005`](../../../../spec/lastenheft.md),
[`LH-FA-RET-006`](../../../../spec/lastenheft.md) (Sichtbarkeit, jetzt auch
über die CLI statt nur über SQL).

**Berührte Spec-Stellen:** — (CLI-Ausgabe-Erweiterung auf bereits
bestehenden Daten, keine neue Architektur-Sicht-Aussage).

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

**Ziel:** Der bestehende `diagnose`-CLI-Befehl
(`internal/bootstrap/wiring.go`, `slice-038`) bekommt einen zusätzlichen
Abschnitt: den ggf. aktuell blockierenden Consumer je Quelle (aus
`cdc.retention_blockers`, `slice-045`) und den aktuellen
`cdc_storage_bytes`-Wert (aus `cdc.metrics`, `slice-046`) — dieselbe
Lese-Disziplin wie die bestehenden Abschnitte (direkte SQL-Abfrage über
den `CDC_READER_DSN`-Pool, kein neuer Port). Ein neuer, externer
`docker exec`-Diagnose-Beleg in `tools/harness/run-integration-tests.sh`
(Muster: bestehender „CLI-Diagnose-Beleg (Normalbetrieb)"-Abschnitt)
zeigt real: kein Blocker (leere Quelle), dann ein realer Blocker
(zurückhängender Consumer), dann der reale `cdc_storage_bytes`-Wert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Berechnungslogik** — die CLI liest ausschließlich bereits
  bestehende SQL-Sichten (`cdc.retention_blockers`, `cdc.metrics`); keine
  neue Domain-/Use-Case-Logik entsteht.
- **Automatische Warnung/Alarmierung bei blockierenden Consumern** —
  bereits in `slice-045` §1 als „anderer Vorgang" ausgeschlossen, gilt
  unverändert.
- **Kombinierter End-zu-Ende-Rundlauf über alle vier `welle-13`-
  Fähigkeiten hinweg** (Consumer blockiert → sichtbar → bestätigt →
  gelöscht → Storage-Bytes reflektiert das) — der bestehende
  `slice-044`-Retention-Beleg in `run-integration-tests.sh` deckt die
  Löschausführung selbst bereits ab; dieser Slice fügt nur die
  CLI-Sichtbarkeit hinzu, keinen neuen kombinierten Testfall. Bestand
  bleibt bewusst stehen: ein eigener Vorgang, falls ein Betriebs-Bedarf
  dafür entsteht.

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

- [x] `diagnose`-CLI zeigt real den aktuell blockierenden Consumer je
      Quelle (aus `cdc.retention_blockers`) und den `cdc_storage_bytes`-Wert
      (aus `cdc.metrics`) — real gegen mindestens einen Zustand ohne
      Blocker und einen mit realem Blocker getestet. Siehe Plan-Nachzug.
- [x] `LH-FA-SST-003` real erweitert: ein externer `docker exec`-Beleg in
      `tools/harness/run-integration-tests.sh` zeigt beide Zustände in der
      `diagnose`-Ausgabe, ohne den laufenden Feed-Container zu beenden
      (analog zum bestehenden CLI-Diagnose-Beleg-Muster). Siehe Plan-Nachzug.
- [x] `make gates` grün, `make test-integration` grün. Beide real ausgeführt
      (Ausgaben im Implementer-Bericht).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-047.md`](../../../reviews/review-slice-047.md)
      (0 HIGH, 1 MEDIUM, 1 LOW — kein Fixrunden-Bedarf, beide Findings direkt
      behoben, Commit `fb6173d`). Verifikation in
      [`docs/reviews/verify-slice-047.md`](../../../reviews/verify-slice-047.md)
      (DoD eigenständig nachgeprüft, keine Rückführung nötig; vier
      Nacharbeits-Punkte VF-1…VF-4, hier nachgezogen).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` nennt die erweiterte
      `diagnose`-Ausgabe (Abschnitt „Aufbewahrung (Retention)"). Neuer
      Beispiel-Block in „Diagnose ausführen" (beide Zustände) plus
      Changelog-Zeile 1.10 (ursprünglich als 1.8 eingetragen; Verifikation
      VF-2 fand zwei rückwirkend fehlende Zeilen bei `slice-045`/`046`
      — siehe Beobachtungs-Register unten).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — zwei Register-Berührungen (3×-Verkörperung, neue Beobachtung).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — beide entfallen.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Wellenlos — Prüfung läuft bei diesem Slice erst nach dem `git mv` nach `done/` (AGENTS.md §3.3), also bei der Closure, nicht im Implementer-Lauf.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (`Diagnose`) | update | neuer Retention-Abschnitt (Blocker, Storage-Bytes) |
| `tools/harness/run-integration-tests.sh` | update | neuer `docker exec`-Diagnose-Beleg für Retention |
| `docs/user/benutzerhandbuch.md` | update | erweiterte `diagnose`-Ausgabe dokumentiert |

### Plan-Nachzug (nach Implementierung)

Regeln dieser Sektion: Implementierungsentscheidungen, die über die Tabelle
oben hinausgehen — Platzierung der beiden realen Zustände im
Compose-Lauf und die SQL-Form der neuen Abfragen.

**1. Zwei reale Zustände, kein neuer Backdating-Aufwand.** Der Plan nannte
einen „zurückhängenden Consumer" als Muster für den realen Blocker-Zustand,
analog zu `slice-044`s Backdating-Technik. Die Implementierung braucht diese
Technik nicht: Der bereits bestehende `CLI-Diagnose-Beleg
(Normalbetrieb)`-Abschnitt in `tools/harness/run-integration-tests.sh` läuft
zu einem Zeitpunkt, an dem `CLI_CONSUMER` real einen von Null verschiedenen
Rückstand trägt (siehe dessen eigener Kommentar: `BACKLOG_CONSUMER`s zweite
Bestätigung lief unmittelbar davor und ist dort verlässlich 0). Da
`cdc.retention_blockers` je Quelle den Consumer mit der kleinsten
bestätigten Position wählt, ist `CLI_CONSUMER` an genau dieser Stelle real
der blockierende Consumer — ein zweiter, eigens konstruierter Rückstand war
nicht nötig; die bestehende Assertion-Gruppe wurde um zwei weitere
Prüfungen ergänzt (Blocker-Zeile, `cdc_storage_bytes`-Zeile).

**2. „Kein Blocker"-Zustand vor jeder Consumer-Bestätigung, nicht über eine
eigene Quelle.** `cdc.retention_blockers` trägt für `src-mvp` erst ab der
ersten Consumer-Bestätigung gegen diese Quelle eine Zeile. Die beiden
Consumer, die `TestMVPRetentionBlockersViewShowsFurthestBehindConsumer`
direkt über den `ConsumerStatePort`-Adapter registriert und bestätigt hatte,
sind zu diesem Zeitpunkt bereits über `t.Cleanup` entfernt (siehe deren
Funktionskommentar) — der neue Beleg läuft deshalb unmittelbar nach dem
`exec_feed`-Funktionsdefinitionspunkt und vor der ersten
`register-consumer`/`acknowledge-consumer`-Zeile des Black-Box-Rundlaufs,
wo `cdc.retention_blockers` für `src-mvp` real keine Zeile trägt.

**3. `retention_blockers`-Abfrage mit `WHERE source_id = $1`, dieselbe
Ein-Zeilen-Erwartung wie die View selbst.** `cdc.retention_blockers` trägt
laut ihrer eigenen `DISTINCT ON (source_id)`-Definition höchstens eine Zeile
je Quelle; `Diagnose` fragt sie mit derselben Quellen-Bindung ab, die die
Funktion bereits für `cdc.heartbeat` trägt (`source model.SourceID`,
Parameter `$1`). `pgx.ErrNoRows` ist die einzige erwartete Abwesenheitsform
(dieselbe Lesart wie beim Betriebsstatus-Zweig oben in derselben Funktion),
kein Fehlerausgang.

**4. `backlog` als `*int64` gelesen, defensiv gegen einen strukturell nicht
erreichbaren Fall.** Die Spalte ist eine Subtraktion über eine korrelierte
`max(commit_position)`-Unterabfrage; ein `NULL`-Ergebnis wäre nur möglich,
wenn eine Quelle ohne jede committete Transaktion trotzdem eine bestätigte
Consumer-Position trägt — strukturell nicht erreichbar, aber der Scan liest
defensiv über einen Zeiger statt mit einem Lesefehler zu enden, dieselbe
Disziplin wie beim bestehenden `cdc_consumer_lag`-Zweig in derselben
Funktion.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` liegt in `done/`
(`slice-045`/`046` liefern `cdc.retention_blockers`/`cdc.metrics`),
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei einer reinen CLI-Ausgabe-Erweiterung auf bereits
  bestehenden SQL-Sichten — falls doch, wäre das ein Zeichen für eine
  unerwartet komplexe Ausgabeformat-Änderung.
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

- Eine Quelle ohne aktuellen Blocker (`cdc.retention_blockers` liefert
  keine Zeile) könnte die CLI-Ausgabe fälschlich als Fehlerzustand statt
  als „kein Blocker" zeigen — dieselbe Klasse Randfall wie `slice-038`s
  §6 Risiko 2 und `slice-045`s §6 Risiko 1. **Ausgang: entfallen.** Der
  `pgx.ErrNoRows`-Zweig behandelt die Abwesenheit explizit als eigenen Fall
  (dieselbe Lesart wie beim Betriebsstatus-Zweig) und gibt „kein Blocker
  (kein Consumer hat je gegen diese Quelle bestätigt)" aus, kein
  Fehlerausgang. Real bestätigt: `make test-integration` zeigt den Text vor
  jeder Consumer-Bestätigung, Prozess-Ausgang bleibt 0.
- Die neue Retention-Sektion könnte das bestehende, stabile
  `diagnose`-Ausgabeformat so verändern, dass `run-integration-tests.sh`s
  bereits bestehende Text-Assertions gegen `LH-FA-ADM-002`…`005`
  (Normalbetrieb/Fehlerzustand-Belege) brechen. **Ausgang: entfallen.** Die
  neuen Zeilen kommen ausschließlich als zusätzliche, angehängte Ausgabe
  nach dem bestehenden `cdc_consumer_lag`-Block — kein bestehender
  Ausgabetext wurde verändert oder verschoben. Real bestätigt: derselbe
  `make test-integration`-Lauf zeigt alle bereits bestehenden Assertions
  (Normalbetrieb, Fehlerzustand-Boundary) unverändert grün, zusätzlich zu
  den beiden neuen Retention-Zuständen.

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

- **Was hat funktioniert:** Der real bereits existierende Rückstand von
  `CLI_CONSUMER` im bestehenden `CLI-Diagnose-Beleg (Normalbetrieb)`-Abschnitt
  ließ sich direkt als realer Blocker-Beleg wiederverwenden — ohne die in
  `slice-044` etablierte Backdating-Technik nachzubauen, siehe Plan-Nachzug
  Punkt 1. Die „kein Blocker"-Prüfung nutzt denselben Effekt in die
  Gegenrichtung: der Zeitpunkt vor jeder `register-consumer`/
  `acknowledge-consumer`-Zeile des Black-Box-Rundlaufs trägt real keine
  Zeile in `cdc.retention_blockers`, ohne eine eigene Quelle oder einen
  eigenen Consumer anlegen zu müssen. Beide Zustände liefen im ersten
  vollständigen `make test-integration`-Lauf nach dem obligatorischen
  `make image`-Rebuild sofort grün.
- **Was ging anders als geplant:** Der erste `make test-integration`-Lauf
  schlug an der neuen „kein Blocker"-Prüfung fehl, weil der Compose-Stack
  noch das alte `ghcr.io/pt9912/pg-change-feed:dev`-Image ohne die neue
  `Diagnose`-Erweiterung führte (`compose.yaml` trägt bewusst keinen
  `build:`-Block) — kein Code-Fehler, sondern ein übersprungener
  `make image`-Lauf vor dem Testlauf, wie `harness/README.md` §Werkzeuge es
  für Build-Kontext-Änderungen vorschreibt. Nach `make image` lief derselbe
  Testlauf real grün. Der Image-Digest änderte sich entsprechend
  (`harness/image-hash.txt`), der Digest-Commit ist Teil dieses Slice.
- **Steering-Loop-Eintrag:** Reviewer-Skill geschärft: Zieht der Reviewer
  im eigenen Verdikt „keine Fixrunde nötig", zieht er die DoD-Checkbox
  „Review durchgeführt" im selben Commit selbst nach — liegt in
  `.harness/skills/reviewer.md §DoD-Checkbox-Nachzug ohne Fixrunde`.
  Auslöser: `BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde`
  (`slice-045`, `slice-046`, `slice-047` — 3×).
- **Beobachtungs-Register (`../observations/`):** Zwei Register-Berührungen
  während des Lebenszyklus dieses Slice: (1) `evidence/slice-047.md` in
  `BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde/` ergänzt — Zähler
  erreichte damit 3×, Architect-Zug
  ([`docs/reviews/architect-verdict-dod-checkbox-review-ohne-fixrunde.md`](../../../reviews/architect-verdict-dod-checkbox-review-ohne-fixrunde.md))
  setzte den Ausgang auf *verkörpert* (siehe Steering-Loop-Eintrag oben).
  (2) Während der direkten Korrektur der Reviewer-Findings F-1/F-2 (Commit
  `fb6173d`) wurde eine neue, eigenständige Beobachtung
  `BEO-PGC/handbuch-versionshistorie-uebersprungen` angelegt (2×,
  `slice-045`, `slice-046` — unter der Schwelle): beide Slices hatten
  `benutzerhandbuch.md` real erweitert, dabei aber die Änderungshistorie-
  Tabelle nicht fortgeschrieben; nachträglich mit den fehlenden Zeilen
  1.8/1.9 geschlossen.
- **Folge-Slices:** keine.
- **Risiken aus §6:** beide *entfallen* — siehe §6.
- **Drei Paarungen:** Wellenlos — Prüfung läuft bei diesem Slice erst nach
  dem `git mv` nach `done/`, nicht in diesem Implementer-Lauf. Vorab
  feststellbar: kein `liegt in`-Feld in dieser Notiz (nichts verkörpert),
  kein Folge-Slice genannt, keine neue Beobachtungs-Register-Zeile — alle
  drei Paarungen sind damit vor der eigentlichen Prüfung bereits vakuos
  erfüllbar.

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
Treffer für `PGC`: `BEO-PGC/test-runner-stiller-ausschluss` (1×, ein neuer
Testfall im `run-integration-tests.sh`-Textmuster muss real ins
Ausgabe-Matching aufgenommen werden), `BEO-PGC/test-isolation-geteilter-zustand`
(1×, kein direkter Bezug — dieser Slice fügt keinen neuen Zustands-
teilenden Testfall hinzu). Keiner erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
