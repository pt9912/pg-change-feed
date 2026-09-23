# Slice transformationen-antragsweg-usecase: Antragsweg, Verarbeitung — Use Cases `SetTransformation`/`RemoveTransformation` mit K1–K4, Regelstand-Port, Verdrahtung und Dauerhaftigkeit

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Happy Path,
Boundary), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (Administration
über SQL),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 1/3/6 (Konfliktfreiheit, Dauerhaftigkeit) und Folgepflicht 3,
[`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
(dauerhafter, tabellen-scoped Träger — Muster),
[`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 5
(Existenz der Spalte an der Quelle),
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Live-Reload), [`ADR-0028`](../../adr/0028-inbound-use-cases.md) (Inbound Use
Cases), [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach
Fähigkeiten),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(keine Domänenlogik in SQL).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md)
(Antrags-Datensatz, Fehlertexte, durch `slice-transformationen-spec-nachzug`),
[`ARC-002`](../../../../spec/architecture.md),
[`ARC-003`](../../../../spec/architecture.md),
[`ARC-004`](../../../../spec/architecture.md),
[`ARC-007`](../../../../spec/architecture.md) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein `set_transformation`-/`remove_transformation`-Antrag wird
verarbeitet und wirkt dauerhaft: zwei Use Cases prüfen die
Konfliktfreiheits-Invarianten K1–K4 (ein Verstoß endet `failed` mit dem
Fehlertext der Spec, der Regelstand bleibt unverändert), ein neuer Outbound
Port liefert die Regelstand-Ableitung und die Spaltenliste der Quelltabelle,
`applyAdministrationRequest` trägt die Regel in die laufende
`Assembler`-Bindung nach, und der aus den `applied`-Zeilen abgeleitete
Regelstand geht bei jedem Pfad, der eine Bindung anlegt (Prozessstart über
`activatedTableBindings`, Aktivierungs-Zweig), in die Bindung ein — Muster von
[`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Schema, Funktionen, Grants, Guard** — `antragsweg-schema`; dieser Slice
  setzt die Spalten `rule_name`/`rule_spec` voraus.
- **Wirkung im Assembler** — `kern-rename` (die Methoden
  `SetTransformation`/`RemoveTransformation` existieren dort; dieser Slice ruft
  sie).
- **Der zweite Regeltyp** — `map-value`; der Use Case prüft gegen den
  Regeltyp-Satz der Domäne und bleibt für einen weiteren Typ unverändert (das
  ist die Aussage von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 4, ihr Beleg steht dort).
- **Der Backfill-Pfad** — `backfill-pfad` liest den hier gelieferten
  Regelstand-Port.
- **Die Startreihenfolge** — `start-reihenfolge`; hier bleibt es bei der
  bestehenden Reihenfolge (Goroutine neben `stream.Run`).
- **Ein weiterer Lese- oder Anzeige-Weg des Regelstands** (Sicht,
  `diagnose`-Ausgabe) —
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6: nicht Teil dieser Entscheidung.

## 2. Definition of Done

- [ ] K1–K4 gegen Fakes, je an ihre Eingabe gebunden: (K1) ein bereits
      vergebener `rule_name` derselben Tabelle endet `failed`; (K2) eine zweite
      Spaltenregel auf derselben Quellspalte endet `failed`; (K3) ein Zielname,
      der einem anderen Zielnamen oder einem Spaltennamen der Quelltabelle
      (ausgeschlossene Spalten eingeschlossen) gleicht, endet `failed`; (K4)
      eine `column`, die an der Quelle nicht existiert, und ein
      `remove_transformation` gegen einen nicht geführten Namen enden `failed`;
      ein unbekannter `kind` und ein unbekannter Schlüssel in `rule_spec` enden
      `failed` (strikte Dekodierung); der Fehlertext entspricht der Spec, der
      Regelstand bleibt nach jedem `failed` unverändert. *Zu belegen durch:*
      `make test` und je Invariante eine Mutation der Prüfung, die den Test rot
      färbt (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert,
      6×).
- [ ] Regelstand-Ableitung und Spaltenliste im Store: der Adapter liefert den
      Regelstand je Tabelle einer Quelle aus den `applied`-Zeilen der zwei
      Arten in der Ordnung `requested_at`, bei gleichem Zeitstempel nach
      `administration_request_id` (`set_transformation` trägt ein,
      `remove_transformation` nimmt heraus), und die Spaltennamen der
      Quelltabelle aus dem Katalog; die Faltung der Zeilen ist eine reine
      Funktion an einer Stelle, die in der netzlos gemessenen Fläche liegt
      (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`, offen, 2×; der
      Ausschluss des Coverage-Gates trifft nur die Pakete `postgresstorage`,
      `postgresack` und `replication/receive` selbst, gelesen am Dockerfile).
      *Zu belegen durch:* `make test-store` (gleicher Zeitstempel, Reihenfolge,
      Set/Remove-Zyklus, Katalog-Lesart) und `make test` (Faltung).
- [ ] Verdrahtung und Dauerhaftigkeit: ein `set_transformation`-Antrag gegen
      eine aktivierte Tabelle wird `applied`, die laufende `Assembler`-Bindung
      trägt die Regel ohne Neustart; `remove_transformation` nimmt sie wieder
      heraus; der Prozessstart (`activatedTableBindings`) und der
      Aktivierungs-Zweig (Deaktivierung, dann Aktivierung) tragen den
      abgeleiteten Regelstand in die neue Bindung; die erneute Verarbeitung
      eines bereits nachgetragenen, noch `pending` stehenden Antrags ist
      idempotent (der Nachtrag ersetzt nach Namen, K1 prüft nur gegen
      `applied`-Zeilen); ein Antrag gegen eine Tabelle ohne laufende Bindung
      endet gemäß Spec. *Zu belegen durch:* `make test` (Whitebox in
      `internal/bootstrap` mit Fakes: Nachtrag, Prozessstart,
      Aktivierungs-Zweig, Wiederholung) und `make test-store` (Ableitung gegen
      reale Zeilen).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt — die Betreiber-Oberfläche ist mit
      `antragsweg-schema` entstanden und mit
      `slice-transformationen-betriebsdoku` adressiert; dieser Slice ändert
      keinen Vertrag, den das Handbuch beschreibt, ohne dass die Wirkung erst
      mit `e2e-wirkung` belegt ist.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/inbound/` (Datei Arbeitsname) | neu | `SetTransformationUseCase`/`RemoveTransformationUseCase` samt Commands ([`ADR-0028`](../../adr/0028-inbound-use-cases.md), Transport-Typen am Port). |
| `internal/application/port/outbound/` (Datei Arbeitsname) | neu | ein Outbound Port für Regelstand-Ableitung und Spaltenliste der Quelltabelle (Fähigkeits-Schnitt, [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md); [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) schärft [`ARC-004`](../../../../spec/architecture.md) um genau einen neuen Port). |
| `internal/application/usecase/settransformation/`, `internal/application/usecase/removetransformation/` (+ Tests) | neu | K1–K4, `failed`-Fehlertexte gemäß Spec, gegen Fakes; neue Pakete in der netzlos gemessenen Fläche. |
| `internal/domain/model/` (Arbeitsname) | update | Parser `rule_spec` → Regel (strikt) und reine Faltung `applied`-Zeilen → Regelstand; die Regeltypen kommen aus `kern-rename`. |
| `internal/adapters/driven/postgresstorage/` (Adapter, ggf. `sqlexec`/`queries`; + Tests) | update | Lesen der `applied`-Zeilen der zwei Arten und der Spaltennamen der Quelltabelle; dieselbe Adapter-Instanz wie die Aktivierung (im MVP eine Instanz, [`ARC-004`](../../../../spec/architecture.md)). |
| `internal/bootstrap/wiring.go` (+ Tests) | update | Zweige in `applyAdministrationRequest`, Feld in `administrationDeps`, `activatedTableBindings` (Regelstand neben Ausschlussstand), Aktivierungs-Zweig; der Fehlertext der Menge. |
| `internal/application/usecase/*` und `postgresstorage`-Test-Fixtures | prüfen | Build-Kontext für neue Pfade (`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad`, offen, 1×). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Menge der
Antragsarten in `applyAdministrationRequest`“, „die Menge der Pfade, die eine
Bindung anlegen und den Ausschlussstand mitführen“, „die Felder von
`administrationDeps`“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Stellen, die eine Bindung anlegen | `grep -rn 'TableBinding{' internal --include=*.go` | *(Implementer trägt ein)* | jede Stelle trägt den Regelstand mit oder begründet, warum nicht (Prozessstart `activatedTableBindings`, Aktivierungs-Zweig; der Seed aus `CDC_TABLES` in `parseTables` legt eine leere Bindung an, der Prozessstart überschreibt sie mit dem abgeleiteten Stand) |
| Stellen, die den Ausschlussstand mitführen | `grep -rn 'ExcludedColumns' internal --include=*.go` | *(Implementer trägt ein)* | jede Stelle, die `ExcludedColumns` liest oder setzt, wird auf einen Regelstand-Zwilling geprüft (Test-Fixtures eingeschlossen) |
| Beschreibung der Antragsarten im Doc-Kommentar von `applyAdministrationRequest` | Lesen des Doc-Kommentars und `grep -n 'Spalten-Antragsarten' internal/bootstrap/wiring.go` | *(Implementer trägt ein)* | Kommentar um die zwei Arten und den Regelstand ergänzen |
| Port-Übersichten in Doku | `grep -rn 'ColumnExclusionPort' docs spec harness internal` | *(Implementer trägt ein)* | Aufzählungen der Outbound Ports nachziehen, falls vorhanden |
| Fehlerklassen-Abbildung | Lesen von `classifyRunError` in `internal/bootstrap/wiring.go` | *(Implementer trägt ein)* | Antrags-Fehler bilden **nicht** über den Capture-Pfad ab (`failed`-Vermerk statt Prozessende); keine neue Abbildung nötig, im Bericht belegt |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `antragsweg-schema` in `done/` liegt
(die Spalten und Antragsarten bestehen im Store) und kein anderer Slice in
`in-progress/` liegt (WIP-Limit 1). Die Kopplung K3 der Welle
[welle-backfill-bestand](../welle-backfill-bestand.md) §5 ist mit
`antragsweg-schema` erfüllt; `applyAdministrationRequest` trägt zu diesem
Zeitpunkt bereits den Backfill-Zweig, dieser Slice erweitert ihn additiv.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Use Cases,
  Store-Adapter und Verdrahtung nicht in einem Review tragen — der abtrennbare
  Teil ist der zweite Liefer-Punkt (Store-Adapter für Regelstand und
  Spaltenliste, Tier `make test-store`) als eigener Slice mit Start nach der
  Use-Case-Hälfte.
- `in-progress` → `open` (blockiert): falls der Port-Schnitt (ein Port für
  Regelstand und Spaltenliste, wie
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  ihn zählt, gegen zwei Fähigkeiten nach
  [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md)) eine
  Architect-Entscheidung braucht, oder falls die Spaltenliste an der Quelle
  nicht ohne zusätzliche Rechte der Login-Identität von `CDC_ADMIN_DSN` lesbar
  ist (Rollen-Frage nach
  [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector), `make
test-store` und `make coverage-gate` real grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Der abgeleitete Stand und die laufende Bindung laufen auseinander**, wenn
  der Vermerk `applied` nach dem Nachtrag scheitert
  ([`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)-Muster;
  `processAdministrationRequests` meldet nur „Erfolg nicht vermerkt“).
  *Erwartet, zu belegen durch:* der Idempotenz-Test aus DoD Punkt 3
  (Wiederholung desselben Antrags ist folgenlos). **Ausgang:** *(bei Closure)*
- **K3 prüft gegen die Spaltenliste zum Antragszeitpunkt**; eine spätere
  Spalten-Erweiterung kann den Zielnamen kollidieren lassen. Das ist der Fall,
  den die Prüfung im Assembler (`kern-rename`) fängt und den `e2e-abhilfe` real
  belegt; hier wird er nicht verhindert. *Erwartet, zu belegen durch:*
  Kommentar am K3-Zweig, der die Grenze nennt, und der Verweis in §7.
  **Ausgang:** *(bei Closure)*
- **Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung.** Ab diesem Slice
  sind Regeln setzbar; ein Backfill-Run, der vor `backfill-pfad` läuft,
  lieferte die Rohform (Welle §5). Der Slice ändert das Backfill-Verhalten
  nicht. *Erwartet, zu belegen durch:* die Reihenfolge der Welle
  (`backfill-pfad` folgt diesem Slice unmittelbar) und die Benennung im
  Bericht. **Ausgang:** *(bei Closure: entfallen mit der Closure von
  `slice-transformationen-backfill-pfad`)*
- **Der Port-Schnitt weicht von der ADR-Zählung ab** („ein neuer Outbound
  Port“,
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Schärft; `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, offen, 1×).
  *Erwartet, zu belegen durch:* Begründung im Plan-Nachzug und Review-Prüfung;
  eine Abweichung wird als offene Frage geführt, nicht als Entscheidung.
  **Ausgang:** *(bei Closure)*
- **Store-Tests teilen Zustand** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): unskopierte `DELETE`/`DROP SCHEMA CASCADE` in einem Store-Test
  könnten die neuen Zeilen zerstören. *Erwartet, zu belegen durch:* skopierte
  Bereinigung je Test. **Ausgang:** *(bei Closure)*
- **Kosten und Rollen der Spaltenlisten-Abfrage**: sie läuft je
  `set_transformation`-Antrag über die Administrations-Verbindung
  (`cdc_admin`). *Erwartet, zu belegen durch:* `make test-store` unter der
  realen Rolle (gleiche Katalog-Lesart wie `ColumnExists`, gelesen am Bestand).
  **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Use Cases, Ports, Store-Adapter und Composition Root
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (verkörpert — der Regelstand
ist ein dauerhafter, tabellen-scoped Träger, der Assembler-Cache ist Laufzeit),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×, DoD Punkt
1), `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 2×,
einschlägig — DoD Punkt 2),
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×,
Plan-Zeile), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, Risiko §6),
`BEO-PGC/adapter-fehler-ausgang` (offen, 2×, gesichtet — Antrags-Fehler enden
`failed`, ein Retry ist nicht Teil),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, Risiko §6),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
