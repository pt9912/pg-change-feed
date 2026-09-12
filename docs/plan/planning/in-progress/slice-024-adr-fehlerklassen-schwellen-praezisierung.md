# Slice slice-024: ADR — Fehlerklassen-Trennung und Schwellen-Präzisierung (`SPEC-008`/`SPEC-013`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-7`](../welle-7.md) — der Ende-zu-Ende-Beleg über beide
Seiten der WAL-Rückstand-Schwelle (welle-7 §3) geht über die DoD dieses
Slice hinaus; dieser Slice liefert nur die Entscheidung, keine Laufzeit.

**Bezug:** Kein `LH-FA-*`/`LH-QA-*`-Punkt trägt diese Entscheidung direkt —
sie präzisiert eine bereits vertraglich gedeckte technische Zusage
(`SPEC-008`/`SPEC-013` sind Rang 2, „technisch fortschreibbar", siehe
`harness/README.md` §Source precedence); ein ADR darf die Spezifikation
direkt schärfen, ohne die Lastenheft-Change-Request-Zeremonie
(Baseline-Regelwerk `modul-03-spec.md`). Kein aktives ADR wird von diesem
Slice ersetzt.

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Zeile `replication`), [`SPEC-013`](../../../../spec/pflichtenheft.md)
(Schwellenwert-Tabelle) — beide werden inhaltlich präzisiert, nicht neu
erfunden.
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

**Ziel:** Eine ADR entscheidet zwei Fragen und präzisiert `SPEC-008`/`SPEC-013`
entsprechend: (a) welche `replication`-klassifizierten Fehler-Sentinels als
*Stream-Ordnungs-Verletzung* (hart abbrechend, unverändert) und welche als
*Transport-/Verbindungsstörung* (Kandidat für Schwellen-Überwachung) gelten
— Kandidaten aus dem Code-Bestand:
[`mapper.ErrChangeWithoutBegin`/`ErrCommitWithoutBegin`/`ErrBeginWithoutCommit`](../../../../internal/adapters/driving/replication/mapper/mapper.go)
(hart) gegen
[`receive.ErrReplication`/`outbound.ErrReplication`](../../../../internal/adapters/driving/replication/receive/receive.go)
(Schwellen-Kandidat); (b) die Initialwerte für die WAL-Rückstand-Schwellen
(`cdc_wal_retention_bytes`, `SPEC-009`), die `SPEC-013` heute nicht trägt,
obwohl `SPEC-008`s Prosa sie dort verortet.

**Vorschlag zur Prüfung durch den Architect (nicht Teil der Entscheidung
selbst — die trifft die ADR):** Warnschwelle 100 MiB, Fehlerschwelle 1 GiB.
Begründung des Vorschlags: `SPEC-013`s bestehende Zeitschwellen
(`cdc_capture_lag`, p95 ≤ 1 s · warn > 5 s · error > 60 s) nutzen runde
Zahlen mit deutlichem Warn/Error-Abstand (Faktor 12); 100 MiB/1 GiB spiegelt
denselben Stil (runde Zehnerpotenz, Faktor 10) für eine Byte-Größe, deren
Referenzpunkt der reale Risikofall ist: Ein wachsender WAL-Rückstand auf
einem inaktiven Replication-Slot verbraucht Plattenplatz unbegrenzt (kein
`max_slot_wal_keep_size` in PostgreSQL erzwungen bedeutet: keine Deckelung),
bis die Quelle selbst Schreibfehler durch volle Disk bekommt — 1 GiB ist als
Fehlerschwelle klein genug, um vor einer Disk-Erschöpfung zu reagieren, aber
groß genug, um normale kurzzeitige Lastspitzen nicht fälschlich zu meiden.
Wie `SPEC-013` selbst notiert: „über ADR schärfbar" — spätere reale Werte aus
Produktionsbetrieb ersetzen diesen Startwert per Folge-ADR, ohne diesen
Slice erneut zu öffnen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die tatsächliche Implementierung der Metrik und der Schwellen-Logik** —
  Folge-Slices `slice-025` (Metrik) und `slice-026` (Schwellen-Überwachung im
  Capture-Pfad) übernehmen das; dieser Slice liefert nur die Entscheidung,
  auf der beide aufbauen.
- **Änderung des Verhaltens für Stream-Ordnungs-Verletzungen** — Bestand
  bleibt bewusst stehen: `mapper.ErrChangeWithoutBegin` u. ä. bleiben
  hart abbrechend; die ADR bestätigt das explizit, ändert es nicht.
- **Eine Lastenheft-Change-Request-Zeremonie** — anderer Vorgang: `SPEC-008`/
  `SPEC-013` sind Rang-2-Dokumente („technisch fortschreibbar"), eine ADR
  darf sie direkt schärfen (Baseline-Regelwerk `modul-03-spec.md`); nur eine
  Änderung des Lastenhefts selbst bräuchte die schwerere Form.

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

- [x] ADR `Accepted`: entscheidet die Fehler-Sentinel-Trennung
      (Stream-Ordnungs-Verletzung vs. Transport-/Verbindungsstörung) und die
      WAL-Rückstand-Schwellenwerte (Warn/Fehler, in Bytes). Beleg:
      [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
      (Warn 100 MiB / Fehler 1 GiB, wie im Vorschlag oben).
- [x] `SPEC-008`-Zeile `replication` präzisiert: Aktionstext benennt beide
      Fehlerklassen-Unterarten getrennt (hart abbrechend vs. Schwellen-
      überwacht mit kontrollierter Fortsetzung).
- [x] `SPEC-013` um einen Eintrag für `cdc_wal_retention_bytes` ergänzt
      (Warn-/Fehlerschwelle in Bytes, „über ADR schärfbar" wie beim
      bestehenden `cdc_capture_lag`-Eintrag) — siehe Plan-Nachzug §3: im
      bestehenden `SPEC-013`-Eintrag (§3 Defaults) erweitert, nicht als
      neue §5-Zeile.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-024.md`](../../../reviews/review-slice-024.md)
      (0 HIGH/MEDIUM/LOW, 1 INFO, kein Merge-Blocker), verifiziert unabhängig
      durch [`docs/reviews/verify-slice-024.md`](../../../reviews/verify-slice-024.md)
      (1 LOW: V-1, siehe §7).
- [x] Kein Doku-Update jenseits von `spec/pflichtenheft.md` und dem ADR-Index
      nötig — dieser Slice ändert keinen Betreiber-sichtbaren Vertrag.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Geprüft: Datei existiert nicht (GF-Repo) — entfällt.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      `BEO-PGC/spec008-replication-luecke` → Ausgang bleibt *weiter offen*
      (Zähler unverändert 1×) — dieser Slice liefert erst die Entscheidung,
      die eigentliche Schließung erfolgt in `slice-026`.
      `BEO-PGC/dod-checkbox-nachzug-architect-pfad` neu angelegt (1×,
      weiter offen).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo **mit** Wellen-Betrieb (`welle-7` offen) — verschoben auf die welle-7-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/00NN-replication-fehlerklassen-schwellen.md` | neu | Architect-Entscheidung — Sentinel-Trennung + WAL-Rückstand-Schwellenwerte |
| `docs/plan/adr/README.md` | update | ADR-Index-Eintrag |
| `spec/pflichtenheft.md` (§4 `SPEC-008`, §5 `SPEC-013`) | update | Zeile `replication` präzisiert, `SPEC-013`-Tabelle um `cdc_wal_retention_bytes`-Schwellen ergänzt |
| `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md` | update | Ausgang auf *geplant* (`slice-025`/`slice-026`) gesetzt, sobald die ADR die Umsetzung verbindlich zuweist |

**Plan-Nachzug (Architect-Lauf, vor dem Sensor-Lauf eingetragen —
`modul-09-implementierung.md` §Plan-Nachzug im selben Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md` | neu | Nummer aus dem ADR-Index ermittelt (`0049`, nächste freie nach `0048`) statt des Platzhalters `00NN` oben |
| `spec/pflichtenheft.md` §3 (`SPEC-013`) statt §5 | Klarstellung | `SPEC-013` ist als ID in §3 *Defaults und Konstanten* definiert, nicht in §5 *Metriken*; §5 verweist bereits per Prosa auf diesen Eintrag als Initialwert-Quelle für WAL-Rückstand *und* Capture-Lag — die WAL-Rückstand-Schwelle erweitert deshalb den bestehenden `SPEC-013`-Eintrag in §3 (Name-Feld auf `CDC_THRESHOLDS` verallgemeinert), keine neue Zeile in §5; das löst genau die in §1 benannte Spec-interne Inkonsistenz |
| `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md` | **nicht geändert** | Die Beobachtungs-Register-Fortschreibung ist laut Aufgaben-Vorgabe Planner-Arbeit bei der Closure (§7), nicht Teil dieses Architect-Laufs — bleibt für die Slice-Closure offen |
| `**Bezug:**`/`**Schärft:**`-Felder der neuen ADR ohne `#anker`-Suffix | Klarstellung | Alle 48 bestehenden ADRs dieses Repos verlinken `LH-*`/`SPEC-*`/`ARC-*` ohne Anker (ganzes Dokument, nicht Abschnitt) — Anker-Auflösung ist laut Baseline-Regelwerk `grundlagen-referenz-richtung.md` eine noch nicht aktivierte Reifestufe; diese ADR folgt dem etablierten Bestand statt der Template-Kopfzeile wörtlich zu nehmen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer/Architect) frei — keine harte Abhängigkeit von
einem anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich beim
  Schreiben, dass die Sentinel-Trennung eine Code-Änderung braucht (z. B.
  weil die heutigen Sentinels nicht sauber in „hart" vs. „Schwellen-Kandidat"
  fallen), gehört diese Code-Änderung in `slice-026`, nicht in diese ADR —
  Rückführung zur Zerlegung, falls sich das erst während der Arbeit zeigt.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** ADR `Accepted` **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der vorgeschlagene Byte-Schwellenwert (100 MiB/1 GiB) ist eine
  Stil-Analogie zu `SPEC-013`s Zeitschwellen, keine aus Produktionsdaten
  abgeleitete Zahl — er könnte sich als zu grob oder zu fein erweisen, sobald
  `slice-025`/`slice-026` real dagegen testen. **Ausgang: eingetreten →
  slice-026** — `slice-026`s eigenes §6 führt denselben Punkt bereits als
  Risiko mit Prüfauftrag (dort: „könnte sich beim realen Ende-zu-Ende-Test
  als unpraktikabel erweisen"); die Kennung nimmt den Punkt an, `ADR-0049`
  benennt die Werte zusätzlich explizit als nicht-permanent
  (Re-Evaluierungs-Trigger).
- Die Sentinel-Trennung könnte beim genauen Hinsehen mehr als zwei Kategorien
  brauchen (z. B. ein dritter Fall, der weder eindeutig hart noch eindeutig
  Schwellen-Kandidat ist). **Ausgang: entfallen** — Reviewer und Verifier
  bestätigten unabhängig durch eigenes Lesen von `classifyRunError`
  (`internal/bootstrap/wiring.go:400-421`): Der `switch` bildet exakt fünf
  Sentinels sauber auf genau zwei Kategorien ab, kein dritter Fall existiert
  im heutigen Code-Bestand.

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

- **Was hat funktioniert:** Der Architect-Lauf hat die Sentinel-Trennung
  nicht nur behauptet, sondern selbst im Code verifiziert (vier Dateien
  gelesen, nicht nur den Plan-Vorschlag übernommen) — Reviewer und Verifier
  bestätigten das unabhängig durch eigenes Nachlesen von `classifyRunError`.
  Der vorgeschlagene Schwellenwert (100 MiB/1 GiB) wurde unverändert
  übernommen, mit nachvollziehbarer eigener Begründung ergänzt und explizit
  als vorläufig (Re-Evaluierungs-Trigger) markiert statt als endgültig
  ausgegeben.
- **Was ging anders als geplant:** Vier kleinere Abweichungen ggü. dem
  ursprünglichen §3-Plan, alle im Architect-Lauf als Plan-Nachzug dokumentiert
  und von Reviewer/Verifier einzeln als plausibel bestätigt (ADR-Nummer 0049
  statt Platzhalter, `SPEC-013`-Erweiterung in §3 statt neuer §5-Zeile,
  Register-Fortschreibung auf die Closure verschoben, ADR-Verweise ohne
  `#anker`-Suffix passend zum Bestand). Zusätzlich: Die DoD-Zeile „Review
  durchgeführt" blieb trotz abgeschlossenem Review zunächst auf `[ ]` —
  derselbe Symptom-Typ wie `BEO-PGC/dod-checkbox-nachzug`, aber über einen
  Pfad (Architect-Rolle), den die dortige verkörperte Regel nicht abdeckt.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× und
  wird verkörpert (siehe unten) — der Normalfall.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/spec008-replication-luecke`
  bleibt *weiter offen* (Zähler unverändert 1×) — dieser Slice liefert die
  Entscheidung, `slice-026` die tatsächliche Schließung.
  `BEO-PGC/dod-checkbox-nachzug-architect-pfad` neu angelegt, Beleg
  `evidence/slice-024.md` — Zähler steht bei 1×.
- **Folge-Slices:** `slice-025` (Metrik `cdc_wal_retention_bytes`) und
  `slice-026` (Schwellen-Überwachung mit kontrollierter Fortsetzung) — beide
  bereits Dateien in `open/` (welle-7 §4).
- **Risiken aus §6:** Byte-Schwellenwert-Risiko → *eingetreten* (`slice-026`
  übernimmt die reale Prüfung). Sentinel-Kategorien-Risiko → *entfallen*
  (Reviewer/Verifier bestätigten Vollständigkeit der Zweiteilung).
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-7` offen) — verschoben
  auf die welle-7-Closure (Modul 8 §Rollen-Sequenz für eine Welle).

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
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md`
§Modus-Deklaration pro Sub-Area) — kein feingranularerer Bereich betroffen,
da dieser Slice nur `spec/pflichtenheft.md` und einen neuen ADR ändert.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/README.md`). Treffer für `PGC`:
`BEO-PGC/spec008-replication-luecke` (1×, weiter offen — dieser Slice ist
die Antwort darauf), `BEO-PGC/walsender-wirksamkeit` (nicht einschlägig,
anderes Sub-Thema, siehe welle-7 §6 Out-of-Scope),
`BEO-PGC/rollen-test-abdeckungsluecken` (1×, nicht einschlägig — Test-Layer,
nicht Fehlerklassen). Keiner der Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`, siehe `harness/conventions.md`
§Modus-Deklaration pro Sub-Area).
