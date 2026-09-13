# Slice slice-042: Benutzerhandbuch — SQL-Administration nachdokumentiert

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — reine Doku-Nacharbeit, keine Closure-Bedingung, die
über die eigene DoD hinausgeht (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht). `welle-12` (`slice-036`/`037`) ist bereits
geschlossen; dieser Slice trägt nur die zum Zeitpunkt jener Closure
übersehene Benutzerhandbuch-Aktualisierung nach.

**Bezug:**
[`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (SQL-Administration),
[`LH-FA-CFG-002`](../../../../spec/lastenheft.md) (Live-Deaktivierung),
[`ADR-0050`](../../../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(nur dokumentiert — keine aktive ADR wird geändert).

**Berührte Spec-Stellen:** — (reine Benutzerdoku-Ergänzung, keine
Verhaltensänderung an einer Spec-Stelle). Der Verweis zeigt **aufwärts**:
Die Spec nennt diesen Slice nie (Baseline-Regelwerk
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

**Ziel:** `docs/user/benutzerhandbuch.md` §4 bekommt zwei fehlende
Abschnitte: „Tabelle live aktivieren" (`SELECT cdc.enable_table(...)`,
Ergänzung zum bestehenden „Tabelle aktivieren"-Abschnitt, der nur den
Start-Zeitpunkt-Weg über `CDC_TABLES`/die Konfigurationsdatei beschreibt)
und ein neuer Abschnitt „Tabelle deaktivieren"
(`SELECT cdc.disable_table(...)`) — beide mit Voraussetzung
(`cdc_admin`-Mitgliedschaft), Vorgehen und real beobachtbarem Ergebnis
(Antrags-Queue-Verarbeitung durch die Administrations-Goroutine, kein
Neustart nötig), analog zum bestehenden Abschnitts-Stil. Changelog-Eintrag
1.7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Code-/Verhaltensänderung jeder Art** — `cdc.enable_table`/
  `cdc.disable_table` existieren bereits vollständig und getestet
  (`slice-036`, `slice-037`, beide in `done/`); dieser Slice trägt
  ausschließlich die zum damaligen Zeitpunkt übersehene
  Benutzerhandbuch-Aktualisierung nach, keine neue Fähigkeit
  (Schicht-Abgrenzung: Doku-Schicht, kein Produkt-Code).
  `AGENTS.md` §3.7 gilt unverändert.
- **Consumer-Verwaltung über SQL** — `welle-12` §6 hat das bereits
  ausdrücklich als Out-of-Scope benannt (keine schreibende
  Consumer-Registrierung/-ACK über SQL); es existiert nichts, was hier
  nachdokumentiert werden könnte.
- **`spec/pflichtenheft.md`/`spec/architecture.md`-Änderungen** — beide
  sind bereits vollständig (`slice-037`s Sequenzdiagramm-Korrektur,
  `LH-FA-ADM-001`/`LH-FA-CFG-002` im Lastenheft); nur das Benutzerhandbuch
  hatte die Lücke.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Neuer Abschnitt „Tabelle live aktivieren" in
      `docs/user/benutzerhandbuch.md` §4, direkt nach dem bestehenden
      „Tabelle aktivieren"-Abschnitt: Voraussetzung (`cdc_admin`-
      Mitgliedschaft, `CDC_ADMIN_DSN`), Vorgehen
      (`SELECT cdc.enable_table('<source_id>', '<schema>', '<tabelle>')`),
      Ergebnis (Antrag landet in `cdc.administration_request`, die
      Administrations-Goroutine des laufenden Prozesses verarbeitet ihn
      ohne Neustart, `status` wird `applied`/`failed`).
- [x] Neuer Abschnitt „Tabelle deaktivieren" (analoges Format,
      `SELECT cdc.disable_table(...)`, Ergebnis: Erfassung endet für die
      Tabelle, der Prozess läuft unverändert weiter).
- [x] Änderungshistorie: Versionsfeld → 1.7, neuer Eintrag mit Bezug auf
      `LH-FA-ADM-001`/`LH-FA-CFG-002`, `ADR-0050`, `slice-036`/`037`.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-042.md`](../../../reviews/review-slice-042.md)
      (0 HIGH, 1 MEDIUM, 2 LOW), Fixrunde in Commit `9cae8f9`, bestätigt
      in [`docs/reviews/review-slice-042-fixrunde.md`](../../../reviews/review-slice-042-fixrunde.md).
      Verifikation in
      [`docs/reviews/verify-slice-042.md`](../../../reviews/verify-slice-042.md)
      (DoD eigenständig nachgeprüft, keine Rückführung nötig).
- [x] Doku-Update — entfällt zusätzlich: dieser Slice **ist** das
      Doku-Update, kein weiterer öffentlicher Vertrag wird berührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — keine Beobachtung angefallen; `BEO-PGC/verwaltung-keine-sql-administration` ist bereits seit `welle-12` verkörpert und wird durch diesen reinen Doku-Nachtrag nicht erneut berührt.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — entfallen.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Wellenlos — Prüfung läuft hier, siehe §7.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | zwei neue §4-Abschnitte, Changelog 1.7 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei reiner Doku-Ergänzung von zwei Abschnitten — falls doch,
  wäre das ein Zeichen, dass die Abschnitte selbst weitere, hier nicht
  vorgesehene Inhalte anziehen.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die neue Doku könnte den realen Verhalten-Detail (Antrags-Queue,
  asynchrone Verarbeitung statt sofortiger Wirkung) ungenau oder
  irreführend vereinfacht darstellen, wenn der Implementer den
  `slice-037`-Code nicht erneut liest, sondern nur aus der Erinnerung an
  diesen Plan schreibt. **Ausgang: entfallen.** Der Implementer las
  `slice-037`, `tools/schema/nacharbeit-administration.sql` und
  `ADR-0050` real erneut, statt aus der Erinnerung zu schreiben; der
  Reviewer und der Verifier bestätigten unabhängig voneinander, dass
  das asynchrone Antrags-Framing (Status `pending`→`applied`/`failed`,
  kein „sofort aktiv") sachlich korrekt gegenüber der realen
  Implementierung ist.

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

- **Was hat funktioniert:** Der Implementer las die reale Implementierung
  (`slice-037`, `tools/schema/nacharbeit-administration.sql`) statt aus
  der Erinnerung an diesen Plan zu schreiben — das machte die neue Doku
  von Anfang an sachlich korrekt (asynchrones Antrags-Framing, exakte
  Funktionssignaturen).
- **Was ging anders als geplant:** Der Reviewer fand 1 MEDIUM (fehlende
  Tabellen-/`REPLICA IDENTITY`-Voraussetzung in den neuen Abschnitten)
  und 2 LOW (uneinheitliches Poll-Format, Prosa-Bruch in der
  Changelog-Zeile), alle drei in der Fixrunde behoben und vom Reviewer
  bestätigt. Zwei Docs-Check-Reparatur-Commits waren nötig
  (`a98e099`), weil eine Kennungs-Verlinkung mit verschachtelten
  Backticks in Link-Klammern (`` [`ID`](pfad) ``) von d-check in diesem
  Report abgelehnt wurde, obwohl dasselbe Muster andernorts im Repo
  funktioniert — behoben durch bare ID-Linktext ohne Backticks.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3×.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen.
- **Folge-Slices:** keine.
- **Risiken aus §6:** *entfallen* — siehe §6.
- **Drei Paarungen (wellenlos, hier geprüft):** (a) Anker — kein
  Steering-Loop-Eintrag mit `liegt in`-Feld, nichts zu prüfen.
  (b) Folge-Slice — keiner genannt, nichts zu prüfen. (c) Register —
  keine neue/geänderte Beobachtung in diesem Slice, nichts zu prüfen.
  Alle drei grün, kein Rot.

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
Keine Treffer für `PGC`, die spezifisch Benutzerdokumentation betreffen.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
