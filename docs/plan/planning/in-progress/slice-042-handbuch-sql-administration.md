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

- [ ] Neuer Abschnitt „Tabelle live aktivieren" in
      `docs/user/benutzerhandbuch.md` §4, direkt nach dem bestehenden
      „Tabelle aktivieren"-Abschnitt: Voraussetzung (`cdc_admin`-
      Mitgliedschaft, `CDC_ADMIN_DSN`), Vorgehen
      (`SELECT cdc.enable_table('<source_id>', '<schema>', '<tabelle>')`),
      Ergebnis (Antrag landet in `cdc.administration_request`, die
      Administrations-Goroutine des laufenden Prozesses verarbeitet ihn
      ohne Neustart, `status` wird `applied`/`failed`).
- [ ] Neuer Abschnitt „Tabelle deaktivieren" (analoges Format,
      `SELECT cdc.disable_table(...)`, Ergebnis: Erfassung endet für die
      Tabelle, der Prozess läuft unverändert weiter).
- [ ] Änderungshistorie: Versionsfeld → 1.7, neuer Eintrag mit Bezug auf
      `LH-FA-ADM-001`/`LH-FA-CFG-002`, `ADR-0050`, `slice-036`/`037`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update — entfällt zusätzlich: dieser Slice **ist** das
      Doku-Update, kein weiterer öffentlicher Vertrag wird berührt.
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
  diesen Plan schreibt. **Ausgang:** <bei Closure einzutragen>

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
Keine Treffer für `PGC`, die spezifisch Benutzerdokumentation betreffen.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
