# Slice slice-009: Consumer-Verwaltung (CON-001…006)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-3.

**Bezug:** [`LH-FA-CON-001`](../../../../spec/lastenheft.md)…006, [`ADR-0013`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`ARC-004`](../../../../spec/architecture.md), [`ARC-006`](../../../../spec/architecture.md)
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.
**Autor:** pt9912. **Datum:** 2026-09-09.

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

**Ziel:** Consumer-Verwaltung real: ConsumerStatePort implementieren (PostgreSQL-Adapter auf cdc.consumer/cdc.consumer_position), Registrierung/Position/Entfernung als Use Cases, monotoner ACK.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Performance-Tuning — nicht MVP.
- **Verdrahtung der Consumer-Use-Cases in `internal/bootstrap`** —
  Klasse 2 (Bestand bleibt bewusst stehen): Die Runtime konsumiert
  Positionen noch nicht; dieser Slice liefert Port + Adapter + Use
  Cases, kein Composition-Root-Zug. Kein Folge-Slice mit Kennung, da
  noch keiner geschnitten ist (Review F-10, Architect-Verdikt
  2026-09-10) — der Bedarf wird beim nächsten Schneiden innerhalb
  welle-3 oder einer Folge-Welle bewertet.

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

- [ ] `PostgresConsumerStateAdapter` real — Teil-Beleg zu
      [`LH-FA-CON-003`](../../../../spec/lastenheft.md).
- [ ] Consumer-Use-Cases am Inbound-Port (Registrierung/Position/ACK/Entfernung)
      — Teil-Beleg zu [`LH-FA-CON-001`](../../../../spec/lastenheft.md)/004/006.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für den `make test-store`-Sensor-Vertrag
      (`harness/README.md`), falls berührt — geprüft: Rollout-Schritt
      nachgetragen (F-9, Fixrunde `c10722e`).
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
| `tools/schema/schema.yaml` | update | `consumer`/`consumer_position`-Tabellen in die d-migrate-Kette je [`ADR-0043`](../../../../docs/plan/adr/README.md); Spalten-Rename `position` → `acknowledged_position` (reserviertes Wort, Katalog-Kanonisierung führte den textlichen Rollout-Vergleich sonst als Drift) |
| `internal/application/port/outbound/consumerstate.go` | neu | `ConsumerStatePort` (Fähigkeits-Port je [`ADR-0034`](../../../../docs/plan/adr/README.md)); eigene Fehlerklasse `ErrConsumerStateStorage` + Registrierungs-Grenze `ErrConsumerUnregistered` statt Mitnutzung der ChangeStore-`storage`-Klasse (Review F-4) |
| `internal/application/port/inbound/consumer.go` | neu | Vier Use-Case-Ports (Registrierung/Position/ACK/Entfernung), Transport-Typen am Port je [`ADR-0042`](../../../../docs/plan/adr/README.md), [`ADR-0028`](../../../../docs/plan/adr/README.md) |
| `internal/application/usecase/{register,acknowledge,position,remove}/*.go` | neu | Services je [`ADR-0028`](../../../../docs/plan/adr/README.md); **`usecase/reset` nicht realisiert** — kein `LH-FA-CON-*` fordert Reset, die Monotonie (Store-Commit-Vergleich) verhindert bereits, dass ein Reset als ACK durchkommt (`ADR-0013`); `ResetConsumerUseCase` bleibt in [`ADR-0028`](../../../../docs/plan/adr/README.md) „vorgesehen", ohne Träger in diesem Slice — Klasse 2 (Bestand bleibt bewusst stehen) |
| `internal/adapters/driven/postgresstorage/consumerstate.go` + `queries/queries.go` | neu/update | Adapter mit gesperrtem Lese + Domänen-Vergleich (`Advance`) für den Monotonie-Vertrag; FK-Ende für ACK ohne Registrierung |
| `tools/harness/run-store-tests.sh` | update | Rollout-Kette (Schema + `search_path` + `make schema-rollout`) vor dem Testlauf; Bereitschafts-Wächter auf echte Abfrage gehärtet |
| `harness/README.md` | update | Sensors-Zeile `make test-store` nennt den Rollout-Schritt (F-9) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-008 liegt in `done/`; kein anderes Slice in `in-progress/`
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Registrierung + Position + Entfernung als Use Cases sprengen drei
  Liefer-Punkte — Rückzug mit Zerlegung (Registrierung getrennt von
  Position/ACK).
- `in-progress` → `open` (blockiert — Carveout?): d-migrate-Blocker (CHECK-Ausdruck, BEO-PGC/d-migrate-nacharbeit)
  verhindert die consumer-Tabellen-DDL — Blocker, Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss
      ohne offenes HIGH-Finding; monotoner ACK am realen Pfad belegt
      (beobachtbar am Test-Diff der Store-Tests).
- Lerneintrag §7: geschärfte Regel oder benannte Spec-Lücke —
  ohne ihn bleibt der Slice abgelegt, nicht fertig.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- (a) Consumer-State-Tabellen waren in slice-004-DDL bewusst nicht
  getragen (Ports außerhalb des Store-Adapters) — **Ausgang:**
  eingetreten; Träger ist dieser Slice (DDL-Erweiterung über die
  d-migrate-Rollout-Kette).

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

- **Was hat funktioniert:** Der Monotonie-Vertrag (gesperrter Lese +
  Domänen-Vergleich) trug beim ersten Konkurrenz-Test tatsächlich —
  die Mutations-Probe (FOR-UPDATE entfernt → Test rot) bestätigt eine
  echte Sperre, keine behauptete; die Rollout-Kette (Schema-Erweiterung
  über d-migrate) lief ohne neuen Nacharbeit-Bedarf durch den
  `chk_change_operation`-Blocker hindurch — die Consumer-Tabellen
  tragen keinen problematischen CHECK-Ausdruck.
- **Was ging anders als geplant:** Der Implementer-Lauf erweiterte §3
  um sieben Dateien über die geplante Liste hinaus (Review F-1, 9.
  Auftreten der Klasse — Konflikt-Sequenz gelaufen, siehe unten);
  `usecase/reset` wurde nicht realisiert (kein CON-Bedarf,
  [`ADR-0028`](../../../../docs/plan/adr/README.md) nennt ihn nur „vorgesehen"); ein Spalten-Rename (`position` →
  `acknowledged_position`) war wegen eines reservierten Worts nötig.
- **F-2-Closure-Vermerk (ADR-Liste):** `GetConsumerPositionUseCase` und
  `RemoveConsumerUseCase` stehen erstmals im Code, ohne in der
  [`ADR-0028`](../../../../docs/plan/adr/README.md)-Liste geführt zu sein — kein stiller Widerspruch (die
  Liste „Vorgesehen sind" ist offen, [`ADR-0034`](../../../../docs/plan/adr/README.md)/[`ADR-0042`](../../../../docs/plan/adr/README.md) tragen die
  Fähigkeits-Port- und Transport-Struktur mit); Architect-Verdikt
  2026-09-10: Closure-Vermerk reicht, konsistent mit slice-008 F-5.
  Zweites Auftreten der Klasse „ADR-Liste still erweitert" — ein
  drittes erreicht die Register-Schwelle.
- **Steering-Loop-Eintrag:** Implementer-Workflow geschärft (Pflicht-
  Plan-Nachzug vor dem Sensor-Lauf) — liegt in
  `.claude/commands/implement-slice.md §Implementieren und gaten`.
  Auslöser: `BEO-PGC/plan-nachzug` (slice-008, slice-009 — 2×
  seit Erst-Registrierung; Konflikt-Sequenz-Pflicht war bereits bei 9
  Gesamt-Auftreten der zugrundeliegenden Klasse erreicht, siehe
  `docs/reviews/review-slice-009.md` F-1 und den Architect-Verdikt-
  Commit `1d19738`).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-009.md`
  in `BEO-PGC/plan-nachzug/` ergänzt (Zähler 2×, Ausgang bereits vom
  Architect auf *verkörpert* gesetzt) — keine weitere Beobachtung
  angefallen.
- **Folge-Slices:** keine — die Verdrahtungs-Boundary (§1, Klasse 2)
  trägt keinen Folge-Slice; der Bedarf wird beim nächsten Schneiden
  bewertet.
- **Risiken aus §6:** (a) eingetreten (Träger: dieser Slice;
  DDL-Erweiterung über die d-migrate-Rollout-Kette vollzogen).
- **Drei Paarungen:** entfällt hier — das Repo arbeitet mit Wellen;
  die Welle-3-Closure prüft Anker · Folge-Slice · Register auch für
  diesen Slice.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührten Sub-Areas
(Application-Ports/Use-Cases, Store-Adapter, Schemamigration) sind
Segmente des GF-Baums der Modus-Deklaration (`PGC` Greenfield, Doc
führt) — jede erfüllt die Schwelle ≥ 2 von 3 Achsen: Strukturregeln im
Spec-Stratum committet, Phase-Reife 3–4, Evidenz-/Diskrepanz-Risiko
niedrig. Keine Sub-Area zu grob.

**Vorgelagert — offene Beobachtungen sichten:** Register gelesen (fünf
Einträge zum Zeitpunkt der Planung): `d-migrate-nacharbeit` 1×
(berührt — Schema-Erweiterung läuft über dieselbe Rollout-Kette, der
CHECK-Blocker selbst ist von diesem Slice nicht betroffen), `adapter-
fehler-ausgang` 1× (nicht berührt), `a-check-null-abdeckung`
verkörpert (seit welle-1), `walsender-wirksamkeit` 1× (nicht berührt
— kein Replication-Code geändert), `plan-nachzug` 1× zum
Planungszeitpunkt (mit diesem Slice auf 2× — Auflösung s. u.). Kein
Eintrag erreicht mit diesem Slice 3× — keine Lücke.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

*Reiner GF-Hinweis genügt (siehe oben); kein Sub-Area-Block.*
