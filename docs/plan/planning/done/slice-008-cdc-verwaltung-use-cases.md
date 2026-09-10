# Slice slice-008: CDC-Verwaltung als Use Cases (CFG-001/002/003/004)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-3.

**Bezug:** [`LH-FA-CFG-001`](../../../../spec/lastenheft.md)…004, [`ADR-0028`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`ARC-003`](../../../../spec/architecture.md), [`ARC-008`](../../../../spec/architecture.md)
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

**Ziel:** CDC-Verwaltung als Use Cases am Inbound-Port (EnableTable/DisableTable/Status/Liste je [`ADR-0028`](../../../../docs/plan/adr/README.md)) — die Aktivierung läuft nicht mehr per Runner-Seed-SQL, sondern als Use Case; das cdc.source_table-Modell trägt die Tabellen-Verwaltung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Vollständige Konfigurationsschicht (TOML/YAML) — **Klasse 2 (Bestand
  bleibt bewusst stehen):** die ENV-Minimalverdrahtung aus slice-007
  bleibt als Bestand; die vollständige Schicht ist anderer Vorgang mit
  eigenem Bedarfsmoment und folgt in späteren Wellen.

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

- [x] `EnableTableUseCase`/`DisableTableUseCase` am Inbound-Port —
      Teil-Beleg zu [`LH-FA-CFG-001`](../../../../spec/lastenheft.md)/002.
- [x] `GetStatusUseCase`/`ListTablesUseCase` (Kanon-Bezeichner je
      [`ADR-0028`](../../../../docs/plan/adr/README.md)) — Teil-Beleg zu
      [`LH-FA-CFG-003`](../../../../spec/lastenheft.md)/004.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für den ENV-Container-Vertrag (`compose.yaml`), falls
      berührt — geprüft: ENV-Vertrag unverändert, Item entfällt (Verifier-Stichprobe).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). — getragen durch die Welle-3-Closure (Repo arbeitet mit Wellen).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/inbound/verwaltung.go` | neu | EnableTable/DisableTable/Status/Liste je [`ADR-0028`](../../../../docs/plan/adr/README.md), Transport-Typen am Port je [`ADR-0042`](../../../../docs/plan/adr/README.md) |
| `internal/application/port/outbound/tableactivation.go` | neu | `TableActivationPort` als Fähigkeits-Port je [`ADR-0034`](../../../../docs/plan/adr/README.md) (Aktivierung außerhalb des Store-Ports) |
| `internal/application/usecase/{enable,disable,status,list}/*.go` | neu | Services + Command/Query/Result je [`ADR-0039`](../../../../docs/plan/adr/README.md); Status/Liste ergänzt (DoD-Deckung, §2) |
| `internal/adapters/driven/postgresstorage/tableactivation.go` + `queries/queries.go` | neu/update | Aktivierungs-Adapter mit Bindungs-/Publikations-Zugriff; Zustands-View (Published/Retained, Review F-2) |
| `internal/bootstrap/wiring.go` | update | Aktivierung als Use Case vor dem Stream-Start je [`ADR-0026`](../../../../docs/plan/adr/README.md) |
| `tools/harness/run-integration-tests.sh` | update | Aktivierung über die Use-Case-Fähigkeit statt Seed-SQL |
| `compose.yaml` | update | Seed-SQL-Aktivierung entfällt; Feed-Container aktiviert beim Start |
| `test/integration/mvp_test.go` | update | MVP-Test fährt die Use-Case-Abfolge; Zustands-Wächter (Disable → Retained) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): welle-2 schließt (slice-004…006 in `done/`); kein anderes Slice in
`in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Inbound-Port für EnableTable/DisableTable/Status/Liste lässt
  sich nicht in einem Slice als Use-Case-Trio liefern (mehr als drei
  Liefer-Punkte) — Rückzug mit erneuter Zerlegung (z. B. Enable/Disable
  getrennt von Status/Liste).
- `in-progress` → `open` (blockiert — Carveout?): Der Replication-Stand verhindert die Aktivierung über den Use Case
  (Compose/Scheitern am Seed-SQL-Ersatz) — Blocker, Priorität offen,
  Carveout-Prüfung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD vollständiges Häkchen + `make gates` grün + Review-Schluss
      ohne offenes HIGH-Finding; die Seed-SQL-Aktivierung ist aus dem
      Integrationstest-Prüfpfad entfernt (beobachtbar am Test-Diff).
- Lerneintrag §7: geschärfte Regel oder benannte Spec-Lücke —
  ohne ihn bleibt der Slice abgelegt, nicht fertig.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- (a) Aktivierung über Use Case vs. Runner-Seed-SQL — **Ausgang:**
  eingetreten (Träger: dieser Slice; der Seed-SQL-Prüfpfad entfällt).
- (b) Wirksamkeits-Grenze der Deaktivierung am laufenden Walsender —
  der Publication-Entzug wirkt am laufenden Stream erst nach dessen
  Neuaufbau (PostgreSQL-Verhalten); der Zustands-View (Kataloge) belegt
  die Bindung, nicht das Capture-Zeitverhalten — **Ausgang:** weiter
  offen → `BEO-PGC/walsender-wirksamkeit` im Register (Eintrag
  angelegt, Beleg `evidence/slice-008.md`).

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

- **Was hat funktioniert:** Die Rollen-Sequenz mit Übergabe-Artefakten
  (Report → Fixrunde beim selben Implementer → Verifier in frischem
  Kontext); die F-2-Entscheidung als echte Zustands-Trennung
  (`Published` am Port, `Retained` am Result) statt Grenz-Kommentar —
  der Reviewer-Pfad „verifizierbar" lief am realen Container
  (`TestMVPDisableRetainedState`) rot→grün der Reihenfolge-Fix ausgenommen.
- **Was ging anders als geplant:** Der Implementer-Lauf erweiterte §3 um
  8 Dateien über die geplante Liste hinaus (Review F-1, 8. Auftreten der
  Klasse); der Plan-Nachzug kam als Planner-Commit `81fc8f2` nach, nicht
  vom Implementer committet. F-2 brauchte eine Port-Erweiterung
  (`Published`) — die ADR-Listen-Erweiterung (F-5) folgt unten.
- **F-5-Closure-Vermerk (ADR-Listen):** `ListTablesUseCase` steht erstmals
  in ADR-0028 (dort „vorgesehen"), `TableActivationPort` und
  `GetStatusQuery`-`Publication`-Eingabe ergänzen die ADR-0039-Struktur
  — keine stillen ADR-Widersprüche (Listen sind offen, kanonische
  Quelle gewinnt); die ADRs bleiben unverändert (Accepted-immutable).
- **Steering-Loop-Eintrag:** nichts verkörpert — der Normalfall; die
  Plan-Erweiterungs-Klasse ist erst-registriert (unten), die
  Verkörperung als geschärfte Regel ist die Schwelle-Antwort der
  Welle-3-Closure (Lese-Schritt).
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/walsender-wirksamkeit/`
  neu angelegt (Beleg `evidence/slice-008.md`, 1×) ·
  `BEO-PGC/plan-nachzug/` neu angelegt (Erst-Registrierung der
  Review-Klasse, Beleg `evidence/slice-008.md`, 1×; Historie
  slice-001…007 benannt, nicht gezählt).
- **Folge-Slices:** keine — die Walsender-Wirksamkeit bleibt Register-
  Beobachtung (`BEO-PGC/walsender-wirksamkeit`), kein eigener Slice;
  die d-migrate-Retirement-Pflicht trägt slice-009/010-Trigger.
- **Risiken aus §6:** (a) eingetreten (Träger: dieser Slice;
  Seed-SQL-Prüfpfad entfällt) · (b) weiter offen →
  `BEO-PGC/walsender-wirksamkeit` im Register.
- **Drei Paarungen:** entfällt hier — das Repo arbeitet mit Wellen; die
  Welle-3-Closure prüft Anker · Folge-Slice · Register auch für diesen
  Slice.

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
(Application-Ports/Use-Cases, Store-Adapter, Bootstrap) sind Segmente
des GF-Baums der Modus-Deklaration (`PGC` Greenfield, Doc führt) —
jede erfüllt die Schwelle ≥ 2 von 3 Achsen: Strukturregeln im
Spec-Stratum committet (Konventionen-Dichte), Phase-Reife 3–4
(Spec-Straten committet, Code folgt), Evidenz-/Diskrepanz-Risiko
niedrig (Inventur prüft Code-Konformität gegen geführte Doku). Keine
Sub-Area zu grob („Backend" kommt nicht vor).

**Vorgelagert — offene Beobachtungen sichten:** Register gelesen (drei
Einträge zum Zeitpunkt der Planung): `d-migrate-nacharbeit` 1× (nicht
berührt — keine Schema-Änderung in diesem Slice), `adapter-fehler-
ausgang` 1× (randständig berührt — die Bootstrap-Verdrahtung erweitert
sich um den Aktivierungs-Zug, aber die Fehler-Ausgangs-Klasse trat im
Lauf nicht auf; unter der Schwelle), `a-check-null-abdeckung` verkörpert
(seit welle-1). Kein Eintrag erreicht mit diesem Slice 3× — keine Lücke.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

*Reiner GF-Hinweis genügt (siehe oben); kein Sub-Area-Block.*
