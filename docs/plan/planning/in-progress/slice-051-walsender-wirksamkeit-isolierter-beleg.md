# Slice slice-051: Publication-Entzug-Wirksamkeit — isolierter Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — Architect-Verdikt
([`docs/reviews/architect-verdict-walsender-wirksamkeit.md`](../../../reviews/architect-verdict-walsender-wirksamkeit.md))
und die vorausgehende Fork-Recherche fanden übereinstimmend: kein Mehr
über die eigene DoD hinaus (ein einzelner isolierter Testabschnitt, keine
Schema-/Domänen-Änderung), siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

**Bezug:** [LH-FA-CFG-002](../../../../spec/lastenheft.md) (CDC-Deaktivierung),
[ADR-0050](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antragsqueue, Live-Reload — nur referenziert, nicht geändert).

**Berührte Spec-Stellen:** — (reine Testabdeckungs-Ergänzung, keine neue
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

**Ziel:** Ein neuer, isolierter Testabschnitt in
`tools/harness/run-integration-tests.sh` belegt real, ob PostgreSQLs
bereits laufende logische Decoding-Session eine Tabelle sofort oder
verzögert ausfiltert, nachdem sie per `ALTER PUBLICATION ... DROP TABLE`
entfernt wurde — **ohne** den regulären `cdc.disable_table`-Antragsweg zu
nutzen, damit die App-seitige `Assembler`-Filterung als
Alternativerklärung ausgeschlossen ist (exakter Ablauf: Architect-Verdikt
oben, Abschnitt „Der isolierte Testansatz"). Je nach realem Ausgang
schließt dieser Slice `BEO-PGC/walsender-wirksamkeit` entweder mit
Ausgang *entfallen* (PostgreSQL filtert sofort) oder trägt einen
benannten Liefer-Punkt für den Fall, dass der Walsender real verzögert
liefert (z. B. Doku-Klarstellung, dass die App-seitige Filterung die
tragende Ebene ist).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ersatz des bestehenden „SQL-Administration Live-Reload-Belegs
  (disable)"** — der Architect-Verdikt stellt klar: beide Belege prüfen
  unterschiedliche Eigenschaften (App-seitige Idempotenz vs. reines
  PostgreSQL-Walsender-Timing); der bestehende Beleg bleibt unverändert
  bestehen.
- **Neuer Introspektions-Mechanismus** (`pg_logical_slot_peek_changes`
  o. ä.) — der Architect-Verdikt hat das geprüft und verworfen (Slot ist
  exklusiv an den laufenden Feed-Prozess gebunden, ein zweiter Slot
  würde die eigentliche Frage verfehlen).
- **Änderung an `Assembler`/`DisableTableUseCase`/`ADR-0050`** — dieser
  Slice ist reine Testabdeckung auf bereits bestehendem Verhalten, keine
  Verhaltensänderung. Ein etwaiger Grace-Wait oder eine
  Doku-Klarstellung (falls der Walsender real verzögert) ist der einzig
  mögliche Code-/Doku-Berührungspunkt und bleibt klein genug, um im
  selben Slice zu bleiben (siehe §2).

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

- [x] Neuer, isolierter Testabschnitt in `run-integration-tests.sh` real
      ausgeführt: dedizierte, über die reguläre `cdc.enable_table`-Kette
      aktivierte Tabelle, `ALTER PUBLICATION ... DROP TABLE` direkt per
      `psql` (ohne `cdc.disable_table`), neue Zeile eingefügt, reales
      Ergebnis gegen `cdc.changes` dokumentiert (erscheint/erscheint
      nicht).
- [x] `LH-FA-CFG-002` real belegt in Bezug auf die Walsender-Timing-Frage
      — Ergebnis eindeutig einer der beiden Auswertungen aus dem
      Architect-Verdikt zugeordnet.
- [x] Falls der Walsender real verzögert liefert: ein benannter
      Liefer-Punkt (Grace-Wait-Doku oder Klarstellung, dass die
      App-seitige Filterung die tragende Ebene ist) — falls PostgreSQL
      real sofort filtert: entfällt dieser Punkt ersatzlos (§1 „Keine
      Mindestzahl"-Prinzip sinngemäß auf DoD-Punkte übertragen; der
      Implementer trägt im Plan-Nachzug nach, welcher Fall eintrat).
      **Entfallen ersatzlos** — PostgreSQL filtert real sofort (siehe §3
      Plan-Nachzug); kein Liefer-Punkt.
- [x] `make gates` grün, `make test-integration` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Report:** `docs/reviews/review-slice-051.md` — 0 HIGH/MEDIUM/LOW,
      2 INFO, keine Fixrunde nötig.
- [x] Doku-Update: nur falls der Walsender real verzögert (siehe oben) —
      Implementer prüft und begründet im Plan-Nachzug. **Entfällt** — kein
      Verzögerungsfall eingetreten, kein Doku-Update nötig.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. **Entfällt** — Repo ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine `reconciliation.md` vorhanden.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | neuer isolierter Walsender-Timing-Testabschnitt |
| `docs/user/benutzerhandbuch.md` | update, falls Walsender real verzögert | Implementer entscheidet, siehe §2 |

### Plan-Nachzug (nach Implementierung)

- **Realer Ausgang:** PostgreSQL (18-alpine, `compose.yaml`) filtert eine
  per `ALTER PUBLICATION ... DROP TABLE` entzogene Tabelle an einer
  bereits laufenden Decoding-Session **sofort** aus — dreimal
  hintereinander real mit `make test-integration` reproduziert
  (identisches Ergebnis: die nach dem Entzug eingefügte Zeile `id=2`
  erscheint in keinem der drei Läufe in `cdc.changes`). Damit trifft die
  Verzögerungs-Annahme aus `BEO-PGC/walsender-wirksamkeit` für die
  geprüfte Version nicht zu — Fall "PostgreSQL filtert sofort" aus dem
  Architect-Verdikt. Kein Liefer-Punkt, kein Doku-Update nötig (§2).
- **Neue Tabelle statt Wiederherstellung:** `feed_mvp_walsender_timing`
  ist eine Wegwerf-Tabelle — kein späterer Abschnitt des Skripts liest
  oder schreibt sie, deshalb entfällt der in Architect-Verdikt Punkt 6
  genannte `ALTER PUBLICATION ... ADD TABLE`-Nachlauf ersatzlos (zweite
  der beiden vom Verdikt genannten Optionen).
- **Publication-Name** `pub_pgc_mvp` direkt im neuen Testabschnitt
  benannt (kein Alias existiert im Skript bisher; `CDC_PUBLICATION` aus
  `compose.yaml` trägt denselben Wert).
- **Testabschnitt-Platzierung:** direkt nach dem bestehenden
  „SQL-Administration Live-Reload-Beleg (disable)", vor dem
  Retention-Beleg — beide Abschnitte laufen am selben, weiterhin
  laufenden Feed-Container, ohne sich gegenseitig zu stören
  (unterschiedliche Tabellen, unterschiedliche IDs-Wertebereiche).
- **Auswertung ohne Abbruch bei beiden Ausgängen:** der neue
  Testabschnitt beendet den Lauf nicht mit `exit 1`, wenn die Zeile
  erscheint — das ist eine offene empirische Frage, kein bekannter
  Fehlerzustand; beide Ausgänge werden nur real geloggt (siehe §1
  Architect-Verdikt-Auswertung).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Architect-Verdikt liegt vor
(`docs/reviews/architect-verdict-walsender-wirksamkeit.md`),
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass der reale Walsender-Verzögerungs-Fall eine über eine
  Doku-Klarstellung/einen Grace-Wait hinausgehende Architektur-Änderung
  braucht, gehört das zurück zur Zerlegung (neuer Architect-Zug nötig,
  nicht mehr Testarbeit).
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

- Der Testabschnitt weist eine Abwesenheit nach (kein Poll auf ein
  eintretendes Ereignis, sondern ein reales, begrenztes `sleep` — Muster
  Architect-Verdikt Punkt 4) — eine zu kurze Wartezeit könnte einen real
  verzögerten Walsender fälschlich als „sofort filternd" auswerten.
  **Ausgang: entfallen.** Dieselbe Wartezeit (`sleep 3`) trägt bereits
  den bestehenden „SQL-Administration Live-Reload-Beleg (disable)" als
  ausreichend akzeptiert; drei unabhängige reale `make test-integration`-
  Läufe lieferten dasselbe Ergebnis (keine Zeile). Strukturelle Grenze
  bleibt bestehen (ein beliebig langsamer Walsender ist mit endlicher
  Wartezeit nie ausschließbar) — dieselbe Grenze gilt bereits für den
  bestehenden Beleg und ist damit kein neues, unadressiertes Risiko.
- Isolation gegenüber dem bestehenden „SQL-Administration Live-Reload-
  Beleg (disable)" — beide Abschnitte dürfen sich nicht gegenseitig
  stören (eigene, dedizierte Tabelle nötig, Architect-Verdikt Punkt 1).
  **Ausgang: entfallen.** Dedizierte Tabelle `feed_mvp_walsender_timing`
  (getrennt von `$ADMIN_TABLE`) implementiert und real verifiziert: beide
  Abschnitte liefen in allen drei Testläufen fehlerfrei nacheinander,
  keine gegenseitige Störung beobachtet.
- Bestätigt sich real die Verzögerungs-Annahme (Walsender liefert trotz
  entzogener Publication weiter), könnte der benannte Liefer-Punkt
  (Grace-Wait/Doku-Klarstellung) den realen Umfang unterschätzen — dann
  greift die in §4 vorab benannte Rückführung `in-progress` → `next`.
  **Ausgang: entfallen.** Die Verzögerungs-Annahme hat sich real nicht
  bestätigt (PostgreSQL filtert sofort, siehe §3 Plan-Nachzug) — die
  Rückführung greift nicht, der Kontingenzfall ist nicht eingetreten.

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

- **Was hat funktioniert:** Der isolierte Testansatz aus dem
  Architect-Verdikt ließ sich unverändert 1:1 umsetzen — dedizierte
  Tabelle über die reguläre `cdc.enable_table`-Kette, direkter
  `ALTER PUBLICATION ... DROP TABLE`-Aufruf per `psql` ohne
  `cdc.disable_table`, reale Wartezeit, Auswertung gegen `cdc.changes`.
  Drei unabhängige `make test-integration`-Läufe lieferten dasselbe,
  eindeutige Ergebnis: die Zeile erscheint nicht — PostgreSQL filtert
  eine entzogene Tabelle an einer bereits laufenden Decoding-Session
  sofort aus.
- **Was ging anders als geplant:** Nichts Wesentliches. Der im Slice-Plan
  vorab benannte Kontingenzfall (Walsender liefert verzögert, Rückführung
  `in-progress` → `next`) ist nicht eingetreten — der einfachere der
  beiden im Architect-Verdikt beschriebenen Ausgänge traf real zu.
- **Steering-Loop-Eintrag:** neuer Sensor — `run-integration-tests.sh`
  trägt jetzt einen isolierten, von der App-seitigen
  `Assembler`-Filterung entkoppelten Beleg für die
  Publication-Entzug-Wirksamkeit am laufenden Walsender — liegt in
  `tools/harness/run-integration-tests.sh` (Abschnitt
  „Publication-Entzug-Wirksamkeit — isolierter Beleg"). Auslöser:
  Architect-Verdikt
  [`docs/reviews/architect-verdict-walsender-wirksamkeit.md`](../../../reviews/architect-verdict-walsender-wirksamkeit.md)
  (Fork-Recherche identifizierte den bestehenden „SQL-Administration
  Live-Reload-Beleg (disable)" als strukturell unzureichend für diese
  Frage).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-051.md`
  in `BEO-PGC/walsender-wirksamkeit/` ergänzt — Ausgang **gestrichen**
  (Verzögerungs-Annahme real widerlegt, `state.md` aktualisiert; Zähler
  abgeleitet 2×: `evidence/slice-008.md`, `evidence/slice-051.md`).
  `BEO-PGC/test-isolation-geteilter-zustand` und
  `BEO-PGC/test-runner-stiller-ausschluss` bleiben unverändert bei 1× —
  dieser Slice hat gegen beide Risiken gearbeitet (dedizierte Tabelle,
  reiner Shell/SQL-Abschnitt außerhalb jedes `-run`-Filtermusters), aber
  keinen neuen Fund derselben Klasse ausgelöst.
- **Folge-Slices:** keine — die Beobachtung ist mit diesem Slice
  geschlossen (Ausgang gestrichen), kein offener Kontingenzfall.
- **Risiken aus §6:** alle drei mit Ausgang **entfallen** — siehe §6.
- **Drei Paarungen:** Repo ohne Wellen-Betrieb, hier geprüft (nach dem
  `git mv`, siehe unten): (a) Anker-Paarung — kein `liegt in`-Feld
  außerhalb des Steering-Loop-Eintrags oben; dessen Zielort
  `tools/harness/run-integration-tests.sh` existiert und trägt den neuen
  Abschnitt. (b) Folge-Slice-Paarung — keine Folge-Slices genannt,
  nichts zu prüfen. (c) Register-Paarung — `BEO-PGC/walsender-wirksamkeit`
  existiert als Verzeichnis, `evidence/` ist nicht leer
  (`slice-008.md`, `slice-051.md`).

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
Treffer für `PGC`: `BEO-PGC/walsender-wirksamkeit` (1×, dieser Slice
liefert den Testbeleg, der ihren Ausgang bestimmt), `BEO-PGC/test-isolation-geteilter-zustand`
(1×, direkt einschlägig — siehe §6), `BEO-PGC/test-runner-stiller-ausschluss`
(1×, jeder neue Abschnitt muss real mitlaufen). Keiner erreicht mit
diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
