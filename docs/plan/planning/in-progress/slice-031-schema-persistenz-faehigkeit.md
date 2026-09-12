# Slice slice-031: Schema-Persistenz-Fähigkeit — `TableSchema`-Modell, `SchemaStorePort`, Adapter (ohne Live-Verdrahtung)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-10`](../welle-10.md) — der Nachweis, dass die
`ADR-0015`-Folgepflicht real eingelöst ist, ist `welle-10`s Closure-Trigger
(§3), kein Einzel-Slice-DoD; dieser Slice liefert nur die
Persistenz-Grundlage, ohne die sich `LH-FA-SCH-004`/`005` selbst noch nicht
schließen (das leisten `slice-032`/`slice-033`).

**Bezug:** [`SPEC-004`](../../../../spec/pflichtenheft.md) (technologie-
unabhängiges `TableSchema`/`SchemaVersion`-Modell), `ADR-0015` (nur
gelesen — keine aktive ADR wird geändert, `ADR-0015` bleibt `Accepted`).

**Berührte Spec-Stellen:** [`SPEC-004`](../../../../spec/pflichtenheft.md)
— die Kennung; die Architektur-Sicht zeigt `SchemaStorePort (ARC-004)`
bereits im Sequenzdiagramm zu `LH-FA-CFG-001.a` als vorgesehenen Port.

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

**Ziel:** Die Persistenz-Grundlage für `ADR-0015`s Option C bauen, ohne den
laufenden Erfassungspfad zu berühren: ein `TableSchema`-Domänenmodell
(Spaltenmenge inkl. Typ-/OID-Information je Version), ein neuer Outbound
Port `SchemaStorePort` (aktuelle Schema-Version einer Tabelle lesen, neue
Version registrieren, `TableSchema` zu einer Version lesen) und ein
Postgres-Adapter, der ihn über eine neue Tabelle `cdc.table_schema`
(ausgerollt über d-migrate, `ADR-0043`) real persistiert. Architect-Skizze:
[`docs/reviews/architect-verdict-slice-030-adr-0015.md`](../../../reviews/architect-verdict-slice-030-adr-0015.md)
§Umsetzungsskizze, Schritte 1–3.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Verdrahtung in `Assembler.Consume`** (dynamische Re-Versionierung bei
  echten `*decode.Relation`-Ereignissen) — Folge-Slice `slice-032`
  übernimmt das; dieser Slice liefert nur die Fähigkeit, keine
  Verhaltensänderung am laufenden Erfassungspfad.
- **Spalten-Typ-/OID-Auswertung im Decoder und Fehlerklasse `schema` für
  inkompatible Typänderungen** — Folge-Slice `slice-033` übernimmt das;
  `TableSchema` hält hier nur die Datenstruktur, keine
  Vergleichs-/Kompatibilitätslogik.
- **Änderung der bestehenden `EnableTable`-Erstaktivierung** — Bestand
  bleibt bewusst stehen: `EnableTableService`/`InsertSchemaVersion` bleiben
  unverändert (Version 1 bei Erstaktivierung bleibt korrekt, siehe
  Architect-Skizze Schritt 6); der neue Port wird nur additiv verdrahtet,
  nicht in den bestehenden Aktivierungspfad eingebaut.
- **CLI-/Verwaltungs-Exposition der Schema-Historie** — kein Lastenheft-
  Kriterium verlangt das in diesem Umfang; anderer Vorgang, falls je
  gebraucht.

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

- [ ] `TableSchema`-Domänenmodell (`internal/domain/model/`) trägt eine
      Spaltenmenge mit Typ-/OID-Information je `SchemaVersionID`; Unit-
      getestet (`make test`).
- [ ] `SchemaStorePort` (`internal/application/port/outbound/`) neu
      definiert (`CurrentVersion`, `RegisterVersion`, `TableSchema` lesen
      — exakte Signatur Implementer-Entscheidung im Plan-Nachzug); ein
      Postgres-Adapter implementiert ihn real gegen eine neue Tabelle
      `cdc.table_schema` (`tools/schema/schema.yaml`, ausgerollt über
      d-migrate). Real gegen PostgreSQL getestet (`make test-store`-Muster):
      Schreiben, Lesen, Round-Trip.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `harness/README.md` §Sensors/`AGENTS.md`, falls ein
      neuer Sensor/Vertrag entsteht — Implementer entscheidet und begründet
      im Plan-Nachzug.
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
| `internal/domain/model/table_schema.go` | neu | `TableSchema`-Domänenmodell: Spaltenmenge inkl. Typ-/OID je Version |
| `internal/application/port/outbound/schemastore.go` | neu | `SchemaStorePort`-Vertrag |
| `internal/adapters/driven/postgresstorage/` | update | Postgres-Adapter für `SchemaStorePort` |
| `tools/schema/schema.yaml` | update | neue Tabelle `cdc.table_schema` |
| `internal/bootstrap/wiring.go` | update | `SchemaStorePort`-Adapter verdrahten (additiv, keine Verhaltensänderung am laufenden Pfad) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Domänenmodell, Port und Adapter zusammen mehr als drei Liefer-Punkte
  oder mehr als zwei Schichten in einer Review-Sitzung nicht mehr prüfbar
  machen, gehört das zurück zur Zerlegung (z. B. Adapter als eigener
  Folge-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Adapter-Test gegen reale
PostgreSQL grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die genaue Form von `TableSchema` (welche Typ-Repräsentation trägt eine
  Spalte — PostgreSQL-OID roh oder ein bereits übersetzter,
  technologieunabhängiger Typ-Enum) könnte sich erst in `slice-033`
  (Typ-Auswertung) als unpassend erweisen und eine Nacharbeit an diesem
  Slice erzwingen. **Ausgang:** <bei Closure einzutragen>
- Eine neue Tabelle `cdc.table_schema` im neutralen Schema könnte mit
  d-migrates View-Signatur-Handling (`ADR-0043`, bekannte Historie bei
  `views:`) in Konflikt geraten, obwohl es sich um eine reguläre Tabelle
  handelt. **Ausgang:** <bei Closure einzutragen>

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
Treffer für `PGC`: `BEO-PGC/schema-evolution-nicht-dynamisch` (1×, weiter
offen — dieser Slice liefert den ersten Baustein der Auflösung, ohne sie
selbst zu vollenden). Keiner der übrigen Treffer erreicht mit diesem Slice
3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
