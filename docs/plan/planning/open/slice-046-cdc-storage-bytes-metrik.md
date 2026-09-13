# Slice slice-046: cdc_storage_bytes-Metrik

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — vierter Slice, unabhängig von den übrigen drei
(reine View-Nacharbeit, analog zu `cdc.metrics`).

**Bezug:** [`LH-FA-RET-006`](../../../../spec/lastenheft.md).

**Berührte Spec-Stellen:** — (reine SQL-View-Ergänzung, keine neue
Architektur-Sicht-Aussage).

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

**Ziel:** `cdc.metrics` (`tools/schema/nacharbeit-observability.sql`)
bekommt eine neue Zeile `cdc_storage_bytes` über `pg_relation_size('cdc.change')`
(oder eine geeignete Aggregation mehrerer `cdc`-Tabellen — Implementer
entscheidet und begründet im Plan-Nachzug). Der Architect-Verdikt
([`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`](../../../../docs/reviews/architect-verdict-retention-loeschausfuehrung.md))
bestätigt vorab: das etablierte View-Owner-Muster (View läuft mit den
Rechten ihres Eigentümers, `GRANT SELECT` nur auf die fertige View an
`cdc_reader`) trägt das bereits — keine Rollen-Erweiterung nötig, da
`pg_relation_size()` eine reguläre, `PUBLIC`-ausführbare Funktion ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Rollen-/Grant-Änderung an `cdc_reader`** — der Architect-Verdikt
  schließt das ausdrücklich aus; die neue Metrik-Zeile kommt über das
  bestehende View-`GRANT`, keine neue Berechtigung.
- **Weitere, laut `nacharbeit-observability.sql` bereits als „nicht
  abgedeckt" benannte Metriken** (`cdc_changes_pending`,
  `cdc_errors_total`, `cdc_wal_retention_bytes` als `cdc.metrics`-Zeile
  statt eigenem Log) — anderer Vorgang, hier nicht angefragt;
  `cdc_wal_retention_bytes` existiert bereits als eigene, strukturierte
  Log-Ausgabe (`slice-020`/`ADR-0049`), nicht als `cdc.metrics`-Zeile.
- **Löschausführung/Sichtbarkeit blockierender Consumer** — `slice-043`/
  `044`/`045`; dieser Slice liefert ausschließlich die Speicher-Metrik.

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

- [ ] `cdc.metrics` trägt eine neue `cdc_storage_bytes`-Zeile, real
      gegen PostgreSQL getestet (ein numerischer Wert > 0 nach dem
      Einfügen von Testdaten).
- [ ] `LH-FA-RET-006` real erfüllt: ein Integrationstest liest die neue
      Metrik über `cdc.metrics` (analog zum bestehenden
      `cdc_capture_lag`-Testmuster).
- [ ] Bestätigt: `cdc_reader`s Grant-Fläche bleibt unverändert (kein
      neuer direkter Grant außerhalb des View-`GRANT SELECT`).
- [ ] `make gates` grün, `make test-integration` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` §„Metriken lesen"
      nennt die neue Zeile.
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
| `tools/schema/nacharbeit-observability.sql` | update | neue `cdc_storage_bytes`-Zeile in `cdc.metrics` |
| `test/integration/integration_test.go` | update | Testfall für die neue Metrik |
| `docs/user/benutzerhandbuch.md` | update | neue Metrik dokumentiert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-043`/`044`/`045`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht zu
  erwarten bei einer einzelnen Metrik-Zeile — falls doch, wäre das ein
  Zeichen für eine unerwartet komplexe Aggregation (mehrere Tabellen,
  Indizes).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `pg_relation_size('cdc.change')` könnte nach einer realen Löschung
  (`slice-043`/`044`) durch PostgreSQLs Tabellen-Bloat (gelöschte
  Tupel, kein automatisches `VACUUM FULL`) einen irreführend hohen Wert
  zeigen, der den tatsächlich freigegebenen Platz nicht widerspiegelt.
  **Ausgang:** <bei Closure einzutragen>
- Der Architect-Verdikt zum View-Owner-Muster wurde vor der
  tatsächlichen Implementierung getroffen — eine reale Prüfung könnte
  eine bisher unbekannte PostgreSQL-Versions-/Berechtigungs-Eigenheit
  aufdecken. **Ausgang:** <bei Closure einzutragen>

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
benannt nicht gezählt). Keiner erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
