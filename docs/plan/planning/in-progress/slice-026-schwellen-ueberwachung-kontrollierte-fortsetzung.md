# Slice slice-026: Schwellen-Überwachung mit kontrollierter Fortsetzung im Capture-Pfad

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-7`](../welle-7.md) — dieser Slice liefert den letzten
Baustein für welle-7 §3 Closure-Trigger (Ende-zu-Ende-Beleg über beide
Seiten der Schwelle); die Welle schließt mit diesem Slice.

**Bezug:** Präzisiert `SPEC-008`/`SPEC-013` gemäß der ADR aus `slice-024`
(Accepted, sobald `slice-024` `done` ist) — kein separater `LH-FA-*`/
`LH-QA-*`-Punkt trägt diese Zusage direkt, siehe `slice-024`s `Bezug`-Feld.

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Zeile `replication`), [`SPEC-013`](../../../../spec/pflichtenheft.md)
(WAL-Rückstand-Schwellen) — dieser Slice macht die durch `slice-024`
präzisierten Zeilen im Code wirksam.
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Capture-Prozess vergleicht die Metrik `cdc_wal_retention_bytes`
(`slice-025`) periodisch gegen die in `slice-024`s ADR festgelegten
Warn-/Fehlerschwellen: unterhalb der Warnschwelle läuft der Betrieb
unverändert fort; zwischen Warn- und Fehlerschwelle läuft er fort, aber
sichtbar (Log-Warnung); oberhalb der Fehlerschwelle klassifiziert er als
`replication`-Fehler und bricht über den bestehenden Abbruchpfad ab —
genau die von `SPEC-008` geforderte „Überwachung über Schwellen; kontrollierte
Fortsetzung". Ein Ende-zu-Ende-Test erzeugt real einen wachsenden
WAL-Rückstand und belegt beide Seiten der Schwelle in einem Lauf (welle-7 §3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Verhalten für Stream-Ordnungs-Verletzungen** (`mapper.ErrChangeWithoutBegin`
  u. ä.) — Bestand bleibt bewusst stehen: Sie bleiben hart abbrechend,
  unabhängig vom WAL-Rückstand (welle-7 §1/§6, `slice-024`s ADR bestätigt
  das explizit).
- **Die Metrik-Erhebung selbst** — Folge-Slice `slice-025` liefert
  `cdc_wal_retention_bytes` bereits fertig; dieser Slice konsumiert sie nur.
- **Weiterleitung an ein externes Alerting-System** — anderer Vorgang
  (welle-7 §6 Out-of-Scope); dieser Slice reagiert im Prozess selbst
  (Log, kontrollierte Fortsetzung/Abbruch).

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

- [x] `SPEC-008`-Zeile `replication` (Transport-/Verbindungsstörungs-Anteil,
      nach `slice-024`s ADR-Trennung) erfüllt: unterhalb der Warnschwelle
      unveränderte Fortsetzung, zwischen Warn-/Fehlerschwelle sichtbare
      Fortsetzung, oberhalb der Fehlerschwelle kontrollierter Abbruch über
      den bestehenden `replication`-Klassifikationspfad. Beleg:
      `internal/bootstrap/wiring.go` (`classifyWALRetention`,
      `runWALRetentionCheck`, `mergeStreamAndWALFaultOutcome`) —
      Schwellen-Vergleich hängt am periodischen WAL-Rückstand-Check
      (`slice-025`) und löst bei Überschreiten der Fehlerschwelle
      `stopStream` + den bestehenden `outbound.ErrReplication`-Klassifikations-
      pfad (`classifyRunError`, `reportFault`) aus.
- [x] Ende-zu-Ende-Test (welle-7 §3): künstlich erzeugter, wachsender
      WAL-Rückstand durchläuft real beide Seiten der Schwelle in einem
      Testlauf, dreimal in Folge grün (`make test-replication` oder
      Äquivalent). Beleg:
      `internal/bootstrap/walretention_endtoend_test.go`
      (`TestWALRetentionThresholdEndToEnd`), `make test-replication` dreimal
      in Folge grün (Testfixture-Schwellen 32 KiB/512 KiB über
      `Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes`,
      Produktions-Startwerte aus `SPEC-013` unverändert).
- [x] Stream-Ordnungs-Verletzungen bleiben nachweislich unverändert hart
      abbrechend (Regressionstest — kein neuer Fortsetzungspfad greift dort
      versehentlich). Beleg:
      `internal/bootstrap/walretention_internal_test.go`
      (`TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`) — rot
      färbende Mutation (Priorität in `mergeStreamAndWALFaultOutcome`
      umgedreht) real gesehen und wieder zurückgesetzt.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `docs/user/benutzerhandbuch.md` §6 Fehlerklassen
      (Zeile `replication` — Aktion ändert sich real von „Sichtbarer Fehler"
      auf „Schwellen-Überwachung; kontrollierte Fortsetzung"). Beleg:
      `docs/user/benutzerhandbuch.md` §4 (WAL-Rückstand prüfen), §6
      (Fehlerklassen-Tabelle), §9 (Grenzwerte), Änderungshistorie 1.4.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
      Plan-Nachzug: Datei existiert in diesem Repo nicht (`PGC` ist
      Greenfield, kein Brownfield-Bootstrap) — Item entfällt.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      `BEO-PGC/spec008-replication-luecke` → *eingetreten* (dieser Slice
      schließt die Lücke real, `evidence/slice-026.md`), sofern der
      Ende-zu-Ende-Test wirklich beide Seiten der Schwelle belegt.
      Plan-Nachzug: Die Bedingung ist erfüllt (E2E-Test belegt real beide
      Seiten der Schwelle, dreimal grün) — das Eintragen selbst
      (`evidence/slice-026.md`, `state.md`) bleibt Closure-Arbeit (§7,
      Modul 6 „Eingetragen wird bei der Slice-Closure") und ist hier bewusst
      nicht vorweggenommen.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` (`classifyRunError`, `Run`) | update | Schwellen-Vergleich gegen `cdc_wal_retention_bytes` (`slice-025`) einhängen; nur der Transport-/Verbindungsstörungs-Zweig aus `slice-024`s ADR bekommt den neuen Fortsetzungspfad |
| `internal/adapters/driving/replication/receive/receive.go` oder Äquivalent (Ort aus `slice-025`) | update | periodischer Schwellen-Check ruft die Metrik ab und meldet den Zustand an den Capture-Loop |
| `docs/user/benutzerhandbuch.md` §6 | update | Zeile `replication` auf den neuen Ist-Zustand nachgezogen |
| Ende-zu-Ende-Test (Ort: Implementer entscheidet — Testcontainer-Suite für Replication) | neu | künstlich wachsender WAL-Rückstand, beide Seiten der Schwelle in einem Lauf |

Der Implementer erweitert diese Liste im ersten Lauf, sobald `slice-024`s
ADR und `slice-025`s konkreter Expositionsweg vorliegen.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-024` und `slice-025` liegen in
`done/` — Priorisiert, `Verantwortlich:` gesetzt, WIP-Limit (1 je
Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  der Schwellen-Vergleich und der Ende-zu-Ende-Test mehr als zwei Schichten
  brauchen (z. B. weil `classifyRunError` grundlegend umgebaut werden muss),
  gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): `slice-024` oder
  `slice-025` sind noch nicht `done`.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** der Ende-zu-Ende-Test
dreimal in Folge grün **und** Closure-Notiz geschrieben (inkl.
`BEO-PGC/spec008-replication-luecke` → *eingetreten*).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein Schwellen-Vergleich, der versehentlich auch den Stream-Ordnungs-
  Verletzungs-Zweig erreicht, würde echte Dateintegritäts-Korruption
  stillschweigend fortsetzen lassen — das genaue Gegenteil der Absicht. Der
  Regressionstest aus §2 ist der Gegenbeleg. **Ausgang:** <bei Closure
  einzutragen>
- Der in `slice-024` vorgeschlagene Byte-Schwellenwert könnte sich beim
  realen Ende-zu-Ende-Test als unpraktikabel erweisen (z. B. zu schnell oder
  zu langsam erreichbar für einen reproduzierbaren Test).
  **Ausgang:** <bei Closure einzutragen — ggf. Folge-ADR>

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
Treffer für `PGC`: `BEO-PGC/spec008-replication-luecke` (1×, weiter offen —
dieser Slice ist die zugewiesene Kennung, die die Lücke schließt),
`BEO-PGC/walsender-wirksamkeit` (nicht einschlägig),
`BEO-PGC/rollen-test-abdeckungsluecken` (nicht einschlägig). Keiner der
übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
