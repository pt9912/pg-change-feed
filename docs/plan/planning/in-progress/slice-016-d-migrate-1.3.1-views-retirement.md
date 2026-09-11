# Slice slice-016: d-migrate-1.3.1-Views-Retirement

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reaktive Wartungsarbeit (Fix-Release liegt vor,
explizites d-migrate-Fix-Signal), keine Closure-Bedingung über die eigene
DoD hinaus, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht (Modul 6).

**Bezug:** [`ADR-0043`](../../../../docs/plan/adr/README.md)
(Re-Evaluierungs-Trigger: Ausweichform entfällt, wenn d-migrate die
Operation ausdrücken kann)

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-11.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** d-migrate-Pin auf v1.3.1 heben (Digest verifiziert:
`sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89`)
und die verbliebene Views-Ausweichform (`tools/schema/nacharbeit-views.sql`)
zurückbauen: d-migrate hat den in slice-015 gemeldeten Befund bestätigt und
behoben (Changelog 1.3.1 — Ursache war `ViewDefinition.sourceDialect`, das
in unserer `schema.yaml` fehlt und deshalb aus dem Fingerabdruck-Vergleich
nicht ausgeblendet wurde; der Report trägt jetzt außerdem den echten
Prozess-Exit statt `status: ok`/`exitCode: 0` bei Drift). Das ist das
explizite Fix-Signal, das `harness/README.md` §Sensors (seit slice-015)
als Bedingung für einen erneuten Real-Test der Views-Drift nennt — nicht
der reguläre Pin-Bump-Zyklus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--provenance-output`/`--migration-overlay` für den Dauerbetrieb
  einführen — Bestand bleibt bewusst stehen: d-migrate selbst nennt das
  für den stationären Rollout sinnvoll (verhindert ein `CREATE OR REPLACE
  VIEW` bei jedem Folgelauf, ohne Herkunft konservativ, aber kein Fehler),
  ist aber eine Convergence-Optimierung, kein Blocker für dieses
  Retirement — der Rollout läuft (idempotent per `CREATE OR REPLACE`
  bislang, nach der Überführung per generiertem `CREATE VIEW`/Diff)
  weiterhin korrekt, nur nicht minimal-diff. Eigene Bewertung des
  Rollout-Workflows nötig, deshalb Folge-Vorgang.
- Rollen-DDL (`nacharbeit-roles.sql`) — anderer Vorgang: unverändert seit
  slice-011/slice-015 kein `raw-sql-text-drift`-Fall, hier nicht berührt.

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] d-migrate-Pin auf v1.3.1 gehoben (`Makefile` `D_MIGRATE_IMAGE`),
      `make schema-validate` grün am neuen Pin.
- [ ] Die drei Views (`active_tables`, `consumer_status`, `changes`) real
      gegen den neuen Pin getestet, mit `source_dialect: postgresql` und
      `columns:`-Signatur (Anhang F.11, verhindert `VIEW_SIGNATURE_UNKNOWN`
      beim Ersetzungspfad) deklarativ in `tools/schema/schema.yaml`
      überführt — `schema migrate --execute` konvergiert (Exit 0), sowohl
      gegen eine leere DB (Erstanlage) als auch gegen eine bereits
      migrierte DB (Folgelauf, kein `ReplaceView`-Blocker).
      `tools/schema/nacharbeit-views.sql` gelöscht, `Makefile`
      `schema-rollout`-Target um den psql-Schritt gekürzt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update falls öffentlicher Vertrag berührt — `harness/README.md`
      §Sensors, `make schema-rollout`-Bindungszeile (die mit slice-015
      verkörperte Test-Kadenz-Regel referenziert die jetzt gelöste
      Views-Drift; Zeile entsprechend aktualisieren).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Makefile` | update | `D_MIGRATE_IMAGE`-Digest auf v1.3.1 gehoben ([`ADR-0043`](../../../../docs/plan/adr/README.md)) |
| `tools/schema/schema.yaml` | update | `views:`-Knoten mit den drei Views ergänzt (`query`, `source_dialect: postgresql`, `columns:`) |
| `tools/schema/nacharbeit-views.sql` | löschen | Bedingung eingetreten — d-migrate 1.3.1 konvergiert deklarativ |
| `Makefile` (`schema-rollout`-Target) | update | psql-Nacharbeit-Schritt für die Views entfernt |
| `harness/README.md` | update | `make schema-rollout`-Bindungszeile — mit slice-015 verkörperte Test-Kadenz-Regel referenziert jetzt die gelöste Views-Drift |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): kein anderes Slice in `in-progress/`
(WIP-Limit 1) — wellenlos, kein Welle-Trigger nötig.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Falls
  `source_dialect`/`columns:` allein nicht genügen und ein weiterer,
  bisher unbekannter View-Metadaten-Fall auftritt, der eine größere
  Schema-Modell-Änderung verlangt — Rückzug mit Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Falls der reale Test
  gegen den neuen Pin einen NEUEN, unbekannten Blocker zeigt (Regression
  gegenüber 1.3.0, nicht der gemeldete/behobene Fall) — Blocker,
  Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss ohne
  offenes HIGH-Finding; `make schema-rollout` (oder `make
  test-integration`, das die Kette real fährt) grün am neuen Pin, sowohl
  gegen eine leere DB (Erstanlage) als auch gegen eine bereits migrierte
  DB (Folgelauf) — beide real getestet, nicht nur behauptet.
- Lerneintrag §7: Beobachtungs-Register-Eintrag `BEO-PGC/d-migrate-nacharbeit`
  aktualisieren — die Views-Hälfte der Beobachtung löst sich damit auf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der reale Test könnte zeigen, dass `source_dialect`/`columns:` allein
  nicht genügen (z. B. weitere Metadaten-Lücke für den Ersetzungspfad
  nötig) — Retirement bliebe dann partiell oder verschöbe sich.
  Wird bei Closure bewertet.
- Der Digest-Bump selbst (Patch-Release, kleineres Risiko als 1.3.0) könnte
  dennoch einen unerwarteten Regressions-Fund gegen den bestehenden
  `schema.yaml`-Bestand auslösen — wird bei Closure bewertet.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührte Sub-Area
(Schemamigration/`tools/schema/`) ist ein Segment des GF-Baums der
Modus-Deklaration (`PGC` Greenfield, Doc führt) — erfüllt die Schwelle
≥ 2 von 3 Achsen. Nicht zu grob.

**Vorgelagert — offene Beobachtungen sichten (einziger Leser, da dieser
Slice wellenlos läuft):** Register gelesen (zehn Einträge zum Zeitpunkt
der Planung): `d-migrate-nacharbeit` 3× — bereits `verkörpert` (seit
slice-015, Test-Kadenz-Regel in `harness/README.md`) — dieser Slice ist
der erwartete technische Auflöser der Views-Hälfte; ein vierter Beleg
(`evidence/slice-016.md`) wird bei Closure ergänzt, ändert den bereits
zugewiesenen Ausgang nicht. Übrige neun ohne Bezug:
`a-check-null-abdeckung` (3×, verkörpert), `adapter-fehler-ausgang` (1×),
`cdc-capture-lag-real` (1×), `health-endpoint-heartbeat` (2×,
eingetreten), `lese-doppelquelle` (2×), `plan-nachzug` (2×, verkörpert),
`plan-vorlagen-defekt` (3×, verkörpert), `rollen-verdrahtung` (1×),
`walsender-wirksamkeit` (1×). Kein Eintrag erreicht mit diesem Slice neu
3× — keine Lücke.

**Modus-Begründungsblock — Umfang.** Reiner GF-Hinweis genügt (siehe oben);
kein Sub-Area-Block.
