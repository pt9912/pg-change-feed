# Slice slice-043: ChangeStorePort-Löschmethode und RunRetentionUseCase

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — erster Slice, liefert die eigentliche
Löschausführung, auf der `slice-044` (Hintergrund-Job/CLI-Trigger)
aufbaut.

**Bezug:**
[`LH-FA-RET-002`](../../../../spec/lastenheft.md),
[`LH-FA-RET-003`](../../../../spec/lastenheft.md),
[`LH-FA-RET-004`](../../../../spec/lastenheft.md),
[`ADR-0009`](../../../../docs/plan/adr/0009-change-store-outbound-port.md),
[`ADR-0011`](../../../../docs/plan/adr/0011-persist-before-ack.md),
[`ADR-0012`](../../../../docs/plan/adr/0012-at-least-once.md),
[`ADR-0014`](../../../../docs/plan/adr/0014-retention-domain-policy.md)
(alle nur umgesetzt — keine aktive ADR wird geändert).

**Berührte Spec-Stellen:**
[`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) (Safe Watermark
der Retention, bereits vorhanden — dieser Slice implementiert das dort
beschriebene Verfahren real). Der Verweis zeigt **aufwärts**: Die Spec
nennt diesen Slice nie (Baseline-Regelwerk
`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP),
`grundlagen-source-precedence.md` §ID-Schema als Klammer).

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

**Ziel:** `ChangeStorePort` (`internal/application/port/outbound/changestore.go`)
bekommt eine neue Löschmethode (z. B. `DeleteChangesBefore` oder
äquivalent — Implementer entscheidet Namen/Signatur und begründet im
Plan-Nachzug), implementiert vom `PostgresChangeStoreAdapter`. Ein neuer
`RunRetentionUseCase` (Application-Schicht, Muster analog zu
bestehenden Use-Cases wie `EnableTableUseCase`) liest für eine Quelle
alle bestätigten Consumer-Positionen, ruft `RetentionPolicy
.AllowsDeletion` je betrachtetem Change real auf und übergibt die
freigegebene Menge an die neue Port-Methode. Kein Aufrufer dieses
Use-Case in dieser Ebene (Hintergrund-Job/CLI) — das liefert
`slice-044`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Hintergrund-Job/CLI-Trigger, der `RunRetentionUseCase` real
  aufruft** — Folge-Slice `slice-044`; dieser Slice liefert den
  Use-Case und die Port-Fähigkeit, nicht den Auslösemechanismus
  (Schicht-Abgrenzung: Application/Adapter-Schicht hier, Composition
  Root/CLI in `slice-044`).
- **Sichtbarkeit blockierender Consumer** (`LH-FA-RET-005`) — Folge-Slice
  `slice-045`; ein anderer Liefer-Fokus (Observability, nicht
  Löschausführung selbst).
- **`cdc_storage_bytes`-Metrik** (`LH-FA-RET-006`) — Folge-Slice
  `slice-046`; berührt eine SQL-View, nicht den Go-Store-Port.
- **Konfigurierbarkeit von `RetentionPolicy.MinAge` zur Laufzeit** (z. B.
  über `CDC_CONFIG_FILE`, `slice-041`) — Bestand bleibt bewusst stehen:
  `RetentionPolicy` wird weiterhin so konstruiert, wie es der aufrufende
  Code heute vorsieht; eine Laufzeit-Konfigurationsanbindung ist ein
  anderer Vorgang, den `welle-13` §6 bereits ausschließt.

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

- [ ] `ChangeStorePort` trägt eine neue Löschmethode; der
      `PostgresChangeStoreAdapter` implementiert sie real gegen
      PostgreSQL (`cdc.change`-Zeilen werden tatsächlich entfernt,
      `make test-store` real belegt).
- [ ] `RunRetentionUseCase` (Application-Schicht) ruft für eine Quelle
      `RetentionPolicy.AllowsDeletion` je betrachtetem Change real auf
      (Alter, Change-Position, alle bestätigten Consumer-Positionen) und
      übergibt ausschließlich freigegebene Changes an die neue
      Port-Methode — mit Fake-Port-Test, der eine gemischte Menge
      (löschbar/nicht löschbar wegen Alter, löschbar/nicht löschbar wegen
      eines zurückhängenden Consumers) real unterscheidet.
- [ ] `LH-FA-RET-002` (Negative: ungültige Konfiguration liefert
      expliziten Fehler statt stiller Übernahme) real erfüllt — Test für
      den Fehlerpfad.
- [ ] `make gates` grün, `make test`/`make test-store` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: falls ein öffentlicher Vertrag entsteht (neue
      Port-Methode ist intern, kein CLI/SQL-Vertrag in diesem Slice) —
      Implementer prüft und begründet im Plan-Nachzug.
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
| `internal/application/port/outbound/changestore.go` | update | neue Löschmethode am `ChangeStorePort` |
| `internal/adapters/driven/postgresstorage/*.go` | update | `PostgresChangeStoreAdapter` implementiert die Löschmethode real |
| `internal/application/usecase/retention/*.go` (neu) | neu | `RunRetentionUseCase`, ruft `AllowsDeletion` real auf |
| `internal/application/port/outbound/consumerposition.go` (falls nötig) | prüfen | ob ein bestehender Port bereits alle Consumer-Positionen einer Quelle liefert, oder eine neue Lesefähigkeit nötig ist |
| `*_test.go` | neu/update | Fake-Port-Tests (Use-Case), reale PostgreSQL-Tests (Adapter) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` liegt flach unter
`docs/plan/planning/` (eröffnet), `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass die Ermittlung „alle bestätigten Consumer-Positionen einer
  Quelle" einen bisher nicht existierenden, größeren Lesezugriffsweg
  braucht als erwartet, gehört das zurück zur Zerlegung (eigener
  Vorbereitungs-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-store` grün
**und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Eine Löschoperation am `ChangeStorePort` könnte mit `ADR-0011`s
  Idempotenz-Pflicht (`PersistTransaction` ist deduplizierbar über die
  interne Transaktions-ID) kollidieren, wenn eine gelöschte Transaktion
  erneut persistiert werden müsste (Crash-Replay) — **Ausgang:** <bei
  Closure einzutragen>
- Es könnte bereits einen etablierten Lesezugriffsweg für „alle
  bestätigten Consumer-Positionen einer Quelle" geben (z. B. über
  `cdc.consumer_status`s Go-Pendant), oder er könnte fehlen und müsste
  neu gebaut werden — der Aufwand ist vor der Implementierung nicht
  exakt bekannt. **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/retention-keine-loeschausfuehrung` (0×,
benannt nicht gezählt — dieser Slice liefert den ersten und
architektonisch entscheidenden Baustein der Auflösung). Keiner der
übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
