# Slice slice-034: Black-Box-Status/Liste gegen `cdc.active_tables`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-11`](welle-11.md) — der Nachweis, dass CDC-Status/
-Liste real über die externe SQL-Sicht lesbar sind, ist `welle-11`s
Closure-Trigger (§3), kein Einzel-Slice-DoD.

**Bezug:** [`LH-FA-CFG-003`](../../../../spec/lastenheft.md),
[`LH-FA-CFG-004`](../../../../spec/lastenheft.md) — dieser Slice erweitert
die Testabdeckung für bestehende Verträge, ändert sie nicht. Kein aktives
ADR wird geändert.

**Berührte Spec-Stellen:** — (reine Testinfrastruktur, keine
Verhaltensänderung an einer Spec-Stelle).
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

**Ziel:** Ein neuer Testfall im Compose-Integrationstest liest den
CDC-Aktivierungsstatus **ausschließlich über die SQL-Sicht**
`cdc.active_tables` (rohes SQL, kein Import von
`status.NewGetStatusService`/`list.NewListTablesService`) und belegt real:
(a) eine aktivierte Tabelle erscheint in der Sicht (`LH-FA-CFG-003` Happy
Path, `LH-FA-CFG-004` Happy Path); (b) eine nie aktivierte Tabelle
erscheint nicht (`LH-FA-CFG-003` Boundary, `LH-FA-CFG-004` Boundary — bei
keiner aktivierten Tabelle liefert die Sicht eine leere Ergebnismenge).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`LH-FA-CFG-003`s Negative-Fall (Tabelle existiert nicht)** — Bestand
  bleibt bewusst stehen: `tools/schema/schema.yaml`s eigener Kommentar zu
  `active_tables` dokumentiert bereits, dass die Sicht bewusst auf
  Bindungs-/Versions-Zeilen beschränkt bleibt und keinen expliziten
  Fehlerpfad für eine nicht existierende Tabelle trägt (das ist Sache der
  Go-Use-Case-Schicht, `TableActivationPort`) — eine SQL-Abfrage gegen
  eine nicht existierende Quelltabelle liefert schlicht keine Zeile,
  keinen Fehler; dieser Unterschied zur internen API-Semantik ist
  architektonisch bereits entschieden, kein Testauftrag.
- **Unterscheidung „aktiviert" vs. „retained" (historische Daten nach
  Deaktivierung)** — anderer Vorgang: `active_tables` bildet diese
  Unterscheidung laut `schema.yaml`-Kommentar bewusst nicht nach, und
  `LH-FA-CFG-002` (Deaktivierung) hat aktuell ohnehin keinen
  Live-Zugriffsweg (`BEO-PGC/verwaltung-keine-sql-administration`) — ein
  Testfall dafür wäre gegenstandslos, bevor diese Fähigkeit existiert.
- **CLI-/Verwaltungs-Exposition des Status** (`LH-FA-SST-003`) — bereits
  als Teil der eigenen, vorgemerkten Feature-Welle
  „Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose" benannt;
  kein Bestandteil dieses reinen Testabdeckungs-Slices.

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

- [x] `LH-FA-CFG-003`/`004` erfüllt: neuer Testfall liest
      `cdc.active_tables` real über SQL (kein Go-Adapter-Import) für eine
      aktivierte und eine nie aktivierte Tabelle im selben Compose-Lauf —
      `TestMVPActiveTablesViewMatchesActivationState` in
      `test/integration/integration_test.go`. **Präzisierung (Verifier-
      Fund V-1, LOW):** black-box real belegt ist die Einzeltabellen-
      Instanz von `LH-FA-CFG-004`s Boundary (eine bestimmte, nie
      aktivierte Tabelle fehlt in der Sicht); die wörtliche Situation
      „keine Tabelle überhaupt aktiviert → leere Gesamtliste" ist am
      Compose-Stack strukturell nicht herstellbar (andere Tabellen sind
      per Verdrahtung aktiv) — bereits auf Go-Use-Case-Ebene belegt
      (`TestListTablesEmpty`, `internal/application/usecase/list/service_test.go:117`),
      kein Bestandteil dieses black-box-Testfalls.
- [x] Vertragstest belegt Übereinstimmung mit der internen
      Use-Case-Semantik für den geprüften Fall (aktiviert/nicht
      aktiviert) — dieselbe Aussage über beide Lesewege — derselbe
      Testfall hält die SQL-Sicht gegen `status.NewGetStatusService`.
- [x] `make gates` grün, `make test-integration` dreimal in Folge grün —
      beide Läufe lokal ausgeführt, siehe Bericht an den Reviewer.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-034.md`](../../../reviews/review-slice-034.md)
      (0 Findings). Verifikation in
      [`docs/reviews/verify-slice-034.md`](../../../reviews/verify-slice-034.md)
      (DoD eigenständig nachgeprüft, 1 LOW-Fund V-1 — siehe oben).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird — kein
      öffentlicher Vertrag berührt (reine Testabdeckung eines
      bestehenden, bereits dokumentierten Lesezugriffswegs); keine
      Doku-Änderung.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. — entfällt: `../reconciliation.md` existiert nicht (Repo ist Greenfield, kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Keine Beobachtung angefallen — reine Testabdeckung ohne neuen Fund.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der `welle-11`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | neuer Testfall: SQL-Lesung gegen `cdc.active_tables` für eine aktivierte und eine nie aktivierte Tabelle |
| `tools/harness/run-integration-tests.sh` | update | **Plan-Nachzug:** neuer Testfall ins `-run`-Muster aufgenommen — ohne diese Zeile liefe er nie in `make test-integration` (`BEO-PGC/test-runner-stiller-ausschluss`, im Slice-Auftrag bereits benannt, aber ursprünglich nicht als eigene §3-Zeile geführt). |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass die „aktiviert vs. retained"-Unterscheidung doch benötigt wird
  (z. B. weil sie leichter zu bauen ist als gedacht), gehört das zurück
  zur Zerlegung — die ursprüngliche Abgrenzung (§1) bleibt sonst nur auf
  dem Papier.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
dreimal in Folge grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Testisolation (welche Zeilen/Tabellen dieser Testfall gegen
  `cdc.active_tables` erwartet) könnte mit bereits im selben Compose-Lauf
  aktivierten Tabellen (`feed_mvp_flow`, `feed_mvp_full`,
  `feed_mvp_schema`) kollidieren, wenn nicht sauber gefiltert wird.
  **Ausgang: entfallen** — Reviewer und Verifier haben die
  Testisolation unabhängig voneinander real geprüft: `feed_mvp_full`
  (Happy Path) bleibt über den gesamten Compose-Lauf aktiviert, auch
  nachdem `TestMVPDisableRetainedState` später `feed_mvp_flow`
  deaktiviert; `feed_mvp_idle` (Boundary) trägt nie eine
  Bindungs-Zeile. Kein Kollisionsfall gefunden, dreifach real
  reproduziert.

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

- **Was hat funktioniert:** Der Implementer erinnerte sich selbst an
  `BEO-PGC/test-runner-stiller-ausschluss` (aus `slice-033`s Closure) und
  nahm den neuen Testfall proaktiv ins Runner-Skript-Muster auf, statt
  die Lücke real eintreten zu lassen — der Sichtungs-Schritt (§8) hat
  hier sichtbar gewirkt. Reviewer und Verifier bestätigten die
  Testisolations- und Mutationsprobe-Behauptungen unabhängig voneinander
  real, keine Übernahme aus dem Implementer-Bericht.
- **Was ging anders als geplant:** Der Verifier fand eine kleine
  Unschärfe (V-1, LOW) zwischen dem im Kopf zitierten
  `LH-FA-CFG-004`-Boundary-Kriterium und dem tatsächlich black-box
  belegten Umfang (Einzeltabellen-Abwesenheit statt „leere
  Gesamtliste") — per Planner-Korrektur in §2 präzisiert, kein
  Sachdefekt.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× —
  der Normalfall. `BEO-PGC/test-runner-stiller-ausschluss` bleibt
  unverändert bei 1× (dieser Slice hat die Lücke vermieden, nicht
  erneut belegt).
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen — siehe §2-Begründung.
- **Folge-Slices:** keine — `slice-035` steht bereits in `welle-11` §4
  als vorgesehener zweiter Slice dieser Welle, unabhängig von diesem
  Slice geplant, wird als nächster Schritt neu geschnitten.
- **Risiken aus §6:** *entfallen* — siehe §6 für Begründung.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-11` offen) —
  Prüfung läuft bei der `welle-11`-Closure.

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
benannt nicht gezählt — nicht Gegenstand dieses reinen Testabdeckungs-
Slices), `BEO-PGC/test-runner-stiller-ausschluss` (1×, weiter offen —
relevant für den Implementer: der neue Testfall muss im
`-run`-Muster des Runner-Skripts erfasst werden). Keiner der übrigen
Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
