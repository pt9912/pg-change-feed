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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-026.md`](../../../reviews/review-slice-026.md)
      (0 HIGH — Prioritätsgarantie explizit geprüft inkl. `-race`-Läufen —,
      1 MEDIUM F-1, kein Merge-Blocker), F-1 per Fixrunde behoben
      (`29a49ea`), unabhängig bestätigt durch
      [`docs/reviews/verify-slice-026.md`](../../../reviews/verify-slice-026.md)
      (eigene Rot-Grün-Gegenproben zu Prioritätsgarantie und
      Zero-Value-Fallback, 1× LOW V-1, siehe §7).
- [x] Doku-Update für `docs/user/benutzerhandbuch.md` §6 Fehlerklassen
      (Zeile `replication` — Aktion ändert sich real von „Sichtbarer Fehler"
      auf „Schwellen-Überwachung; kontrollierte Fortsetzung"). Beleg:
      `docs/user/benutzerhandbuch.md` §4 (WAL-Rückstand prüfen), §6
      (Fehlerklassen-Tabelle), §9 (Grenzwerte), Änderungshistorie 1.4.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
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
      Erledigt: `BEO-PGC/spec008-replication-luecke` → *eingetreten*
      (`evidence/slice-026.md`).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo **mit** Wellen-Betrieb — `slice-026` ist der letzte Slice von `welle-7`; die Paarungen laufen bei der unmittelbar folgenden welle-7-Closure.

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

**Plan-Nachzug (Fixrunde nach Review, [`review-slice-026.md`](../../../reviews/review-slice-026.md) F-1 behoben):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/wiring.go` | update | F-1: Fallback-Logik (`if warnBytes <= 0 { warnBytes = walRetentionWarnBytes }`, dieselbe Form für `errorBytes`) aus der Inline-Stelle in `Run` in die eigene Funktion `resolveWALRetentionThresholds` extrahiert — reiner Wert-Vergleich ohne PostgreSQL-Zugriff, dadurch ohne Testcontainer testbar |
| `internal/bootstrap/walretention_internal_test.go` | update | F-1: `TestResolveWALRetentionThresholdsDefaultsToSpec013` belegt Zero-Value- und negative Config-Felder → `SPEC-013`-Startwerte (100 MiB/1 GiB), nicht 0; `TestResolveWALRetentionThresholdsKeepsPositiveOverride` belegt, dass ein gesetzter Override unverändert bleibt. Rot färbende Mutation (`<= 0` → `>= 0`) real gesehen und zurückgesetzt |

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
  stillschweigend fortsetzen lassen — das genaue Gegenteil der Absicht.
  **Ausgang: entfallen** — Reviewer und Verifier bestätigten unabhängig
  voneinander (eigene Race-Analyse, eigene Rot-Grün-Gegenprobe an
  `mergeStreamAndWALFaultOutcome`, `go test -race`), dass ein WAL-Fault
  einen echten Stream-Ordnungs-Fehler strukturell nicht maskieren kann —
  der Mapper-Sentinel-Pfad prüft nie `ctx.Err()` und gibt seinen Fehler
  unbedingt zurück.
- Der in `slice-024` vorgeschlagene Byte-Schwellenwert könnte sich beim
  realen Ende-zu-Ende-Test als unpraktikabel erweisen (z. B. zu schnell oder
  zu langsam erreichbar für einen reproduzierbaren Test).
  **Ausgang: entfallen** — gelöst über einen Test-Override
  (`Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes`, 32 KiB/512 KiB
  statt 100 MiB/1 GiB): derselbe Produktionscode-Pfad, kalibriert auf
  testbare Zeiträume, Produktions-Startwerte unverändert. Die separate
  Frage, ob 100 MiB/1 GiB die richtigen *realen* Werte sind, trägt bereits
  `ADR-0049`s eigener Re-Evaluierungs-Trigger — kein neuer Slice/BEO nötig.

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

- **Was hat funktioniert:** Die Prioritäts-Garantie
  (`mergeStreamAndWALFaultOutcome`) — der sicherheitskritische Kernpunkt
  dieses Slices — wurde dreimal unabhängig bewiesen statt einmal behauptet:
  Implementer (Rot-Grün am Prioritäts-Code), Reviewer (eigene Race-Analyse
  plus `-race`-Läufe), Verifier (eigene Rot-Grün-Gegenprobe plus statische
  Analyse des Mapper-Fehlerpfads). Derselbe Musterablauf wie bei
  `slice-025`, jetzt am sicherheitsrelevantesten Punkt der ganzen Welle.
- **Was ging anders als geplant:** Der Review fand eine MEDIUM-Lücke
  (fehlender Test für den Zero-Value-Fallback der neuen
  `Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes`-Felder — ein
  Produktionslauf hätte sonst bei jedem Start sofort mit Schwelle `0`
  abbrechen können), behoben per kurzer Fixrunde mit eigener
  Rot-Grün-Verifikation. Außerdem, im Gespräch mit dem Nutzer während der
  Closure-Vorbereitung: Der als „Ende-zu-Ende-Test" bezeichnete
  `TestWALRetentionThresholdEndToEnd` erfüllt zwar welle-7 §3s Anforderung
  (beide Seiten der Schwelle real in einem Lauf), ruft aber `bootstrap.Run`
  in-process auf — nach `ADR-0030`s Testpyramide ist das ein
  **Integrationstest**, kein Black-Box-E2E-Test. Dasselbe gilt für die
  beiden älteren `*_endtoend_test.go`-Dateien (`slice-022`, `welle-6`) und
  für `make test-integration`, das zudem seit `welle-3` nicht mehr über
  MVP-Scope hinausgewachsen ist. Der Nutzer hat dafür eine eigene Welle
  vorgemerkt (Roadmap *Nächste Wellen*, direkt nach `welle-7`) — kein
  Nachtrag an diesem Slice, da die Benennung zum Zeitpunkt der Arbeit dem
  hier verlangten *funktionalen* Ende-zu-Ende-Beleg (beide Schwellen-Seiten
  in einem Lauf) korrekt entsprach; nur die Verwechslung mit der
  *Test-Architektur*-Bedeutung aus `ADR-0030` war die Lücke.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× und
  wird verkörpert — der Normalfall. Die vom Verifier notierte
  DoD-Checkbox-Lücke (V-1) ist wie bei `slice-024`/`slice-025` keine neue
  Beobachtung, sondern derselbe strukturelle Punkt: die Review-Zeile kann
  erst nach dem Review getickt werden.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/spec008-replication-luecke`
  → **eingetreten**, `evidence/slice-026.md` neu angelegt. Damit ist die
  Kette geschlossen: `slice-020` fand die Spec-vs-Code-Lücke,
  `slice-024`/`025`/`026` (welle-7) bauen die tatsächliche Umsetzung.
- **Folge-Slices:** keine neuen — die Testing-Welle (Black-Box-E2E,
  Integrationstest-Nachzug) ist als Vorschau-Zeile in der Roadmap
  vermerkt, noch nicht als Slice geschnitten.
- **Risiken aus §6:** beide *entfallen* (Sentinel-Trennung strukturell
  bewiesen; Test-Reproduzierbarkeit über Test-Override gelöst) — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb — `slice-026` ist der
  letzte Slice von `welle-7`; die Paarungen laufen bei der unmittelbar
  folgenden welle-7-Closure (Modul 8 §Rollen-Sequenz für eine Welle).

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
