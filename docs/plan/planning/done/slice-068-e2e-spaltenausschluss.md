# Slice slice-068: E2E-Beleg — Spaltenausschluss am laufenden Feed-Container

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-18`](../welle-18.md) — letzter Slice, baut auf
`slice-066` (Antrag existiert) **und** `slice-067` (Filterwirkung existiert)
auf; sein grüner Lauf ist `welle-18`s Closure-Trigger (§3).

**Bezug:** [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Haupt-Bezug —
Happy Path und Negative; die Boundary-Klausel ist ausdrücklich nicht
Gegenstand dieses Slice, sie verweist auf `LH-FA-SCH-003` und ist dort
belegt),
[`ADR-0059`](../../../../docs/plan/adr/0059-spaltenauswahl-mechanismus.md)
(nur umgesetzt — keine aktive ADR wird geändert, `ADR-0059` bleibt
`Accepted`).

**Berührte Spec-Stellen:** `—` (dieser Slice erweitert einen bestehenden
Testrundlauf, berührt keine neue Spec-Stelle über `LH-FA-CFG-005` selbst
hinaus).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make test-integration`
(`tools/harness/run-integration-tests.sh`) um einen realen
Spaltenausschluss-Rundlauf erweitern, analog zum bestehenden
Schema-Evolution-Rundlauf: `SELECT cdc.exclude_column(...)` gegen eine
bereits aktivierte Tabelle im laufenden Feed-Container, Poll auf
`status = 'applied'` (dasselbe Muster wie beim bestehenden
`cdc.enable_table`-Live-Reload-Beleg), dann realer Beleg, dass **künftige**
Changes gemäß `LH-FA-CFG-005`s Happy Path den ausgeschlossenen Wert nicht
mehr tragen — der Schlüssel fehlt im Row Image, der Wert steht weder in
`new_data` noch in `old_data` (Happy Path). Eine **vor** dem Ausschluss
erfasste Change bleibt über `cdc.changes` unverändert lesbar; ein
rückwirkendes Entfernen aus bereits persistierten Changes ist mit dem in
`ADR-0059` Teilfrage 3 Option D entschiedenen Wirkort (Filterung in der
Row-Image-Konstruktion, vor jeder Serialisierung) nicht verbunden und
deshalb auch nicht zugesagt. Zusätzlich ein
Negative-Beleg: `cdc.exclude_column` gegen eine nicht existierende Spalte,
Antrag landet real `failed` mit Fehlertext (Negative).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Boundary-Fall reale Spaltenlöschung nach Ausschluss** — bereits als
  Unit-/Konvergenz-Test in `slice-067` gedeckt (`ADR-0059` Teilfrage 4);
  ein zusätzlicher E2E-Beleg für denselben, bereits durch reine Schichtung
  garantierten Pfad liefert keinen neuen Erkenntniswert und würde den
  bereits langen Compose-Stack-Lauf weiter verlängern, ohne ein
  Akzeptanzkriterium zu erfüllen, das nicht schon gedeckt ist.
- **Quellen-/musterweiter Ausschluss über mehrere Tabellen hinweg** —
  Bestand bleibt bewusst außen vor, `ADR-0059` Teilfrage 2 schließt Option
  C bewusst aus; dieser Slice belegt ausschließlich die pro-Tabelle-
  Granularität.
- **`cdc.include_column`-eigener E2E-Rundlauf** — bleibt bewusst Bestand:
  Der Happy-Path-Rundlauf dieses Slice belegt bereits den vollständigen
  Antrags-/Live-Reload-/Filterungs-Pfad für `exclude_column`;
  `include_column` teilt denselben Code-Pfad (`ADR-0059`s symmetrische
  Entscheidung, zwei Antragsarten auf derselben Queue) und ist bereits in
  `slice-066`/`slice-067`s Unit-Tests gedeckt — ein zusätzlicher
  Compose-Stack-Rundlauf für die exakt spiegelbildliche Operation liefert
  keinen neuen E2E-spezifischen Erkenntniswert.

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

- [x] `LH-FA-CFG-005` Happy Path real belegt: `SELECT
      cdc.exclude_column(...)` gegen eine bereits aktivierte Tabelle im
      laufenden Feed-Container, Poll auf `status = 'applied'`, danach ein
      neuer Change ohne den ausgeschlossenen Spaltenschlüssel in
      `new_data`; **zusätzlich** (über den Happy Path hinaus, als
      Abgrenzungsbeleg) bleibt ein bereits vor dem Ausschluss erfasster
      Change über `cdc.changes` unverändert lesbar — die Anforderung sagt
      nur „künftige Changes" zu, ein rückwirkendes Entfernen ist mit dem
      Wirkort aus `ADR-0059` Teilfrage 3 Option D nicht verbunden (siehe
      §1). Beleg:
      `tools/harness/run-integration-tests.sh` (neuer Abschnitt, real
      grün: `make test-integration` Exit 0, beide neuen Belege in der
      Ausgabe; Auslegung des historischen Changes in §3 Plan-Nachzug).
- [x] `LH-FA-CFG-005` Negative real belegt: `cdc.exclude_column` gegen eine
      nicht existierende Spalte, Antrag landet `failed` mit Fehlertext
      (`ErrSourceColumnMissing`-Pfad aus `slice-066`). Beleg: derselbe
      Skript-Abschnitt, Poll auf `status = 'failed'` und Prüfung des
      Fehlertexts.
- [x] `make gates` grün, `make test-integration` real grün (voller
      Compose-Stack-Lauf). Beleg: `make test-integration` Exit 0 (voller
      Compose-Stack-Lauf, beide neuen Abschnitte enthalten), `make gates`
      Exit 0.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Report: `docs/reviews/review-slice-068.md`, Verdikt 0 HIGH, 1 MEDIUM
      (F-1 an den Planner — Plan-Text, kein Reviewer→Implementer-Pfeil),
      1 LOW, 3 INFO, keine Fixrunde.
- [x] Doku-Update: `harness/README.md` §Sensors, Zeile `make
      test-integration` um den neuen Rundlauf-Abschnitt ergänzt (Muster
      der bestehenden Zeile, die jeden Rundlauf-Baustein einzeln nennt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb (`welle-18` offen) — Prüfung läuft bei der `welle-18`-Closure. Dieser Slice ist zusätzlich `welle-18`s Closure-Trigger selbst (§3 der Welle-Datei) — sein grüner `make test-integration`-Lauf ist das *Mehr*, das die Welle schließt.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | update | neuer Rundlauf-Abschnitt (Happy Path + Negative) |
| `harness/README.md` | update | Sensors-Zeile `make test-integration` um neuen Baustein ergänzt |

**Plan-Nachzug (Implementer, 2026-09-14) — Platzierung.** Der Abschnitt
liegt nach dem Go-E2E-Tier und vor dem Lasttest-Beleg, nicht direkt neben
dem SQL-Administration-Live-Reload-Beleg. Grund: die dedizierte Tabelle
(`feed_e2e_column_exclusion`) wird selbst erst über `cdc.enable_table`
aktiviert; kein Go-Testfall dieses Tiers liest sie (jeder Testfall bindet
seine eigene Tabelle), und der Rundlauf bleibt so vollständig vor der
Container-Ende-Grenze der beiden Schema-Negative-Funktionen. Das
Zustands-Überschneidungs-Risiko aus §4 (Muster
`BEO-PGC/test-isolation-geteilter-zustand`) trägt der Abschnitt damit
nicht: er legt keine Tabelle an, die ein anderer Abschnitt liest, und
liest keine, die ein anderer anlegt.

**Plan-Nachzug (Implementer, 2026-09-14) — Auslegung des „historischen
Changes".** Der DoD-Punkt „ein bereits vor dem Ausschluss erfasster
historischer Change über `cdc.changes` gemäß Akzeptanzkriterien geprüft"
ist gegen den gewählten Wirkort auszulegen: `ADR-0059` Teilfrage 3
Option D filtert in der Row-Image-Konstruktion, und `cdc.changes` ist eine
reine Projektion über `cdc.change` — ein **vor** dem Ausschluss
persistierter Change kann den Wert nicht mehr verlieren; ihn „ohne den
Wert" zu prüfen verlangte ein rückwirkendes Umschreiben, das dieselbe
Option D ausschließt. Der Abschnitt prüft deshalb beide Hälften: die
**nach** dem Ausschluss erfasste Change ist über `cdc.changes` **ohne**
Spaltenschlüssel und **ohne** den Wert lesbar (`LH-FA-CFG-005` Happy Path,
„künftige Changes von `t`"), und die **vor** dem Ausschluss erfasste Change
bleibt **unverändert** lesbar — dieselbe Lesart, die
`TestE2ESchemaChangeAddColumn` für die vor einer Schemaänderung erfasste
Change trägt. Der §1-Satz „künftige **und historische** Changes … den
ausgeschlossenen Wert nicht mehr tragen" trifft für den historischen Teil
nicht zu und ist insoweit ungenau; die DoD-Formulierung ist die
auslegbare.

**Plan-Nachzug (Implementer, 2026-09-14) — Baseline und Kontrolle.** Der
Abschnitt setzt zwei zusätzliche Beobachtungen, die der Plan nicht
ausdrücklich nennt: eine Change **vor** dem Ausschluss trägt den
Spaltenwert real (sie trennt „Wert abwesend" nach dem Ausschluss von „die
Spalte war nie Teil des Row Image"), und die **nicht** ausgeschlossene
Spalte bleibt in der nach dem Ausschluss erfassten Change (Kontrolle gegen
ein leeres Row Image als Alternativerklärung). Beide sind Teil desselben
Happy-Path-Belegs, kein vierter Liefer-Punkt.

**Plan-Nachzug (Implementer, 2026-09-14) — `tools/schema/plan.yaml`.** Der
`make test-integration`-Lauf führt `make schema-rollout` aus; dessen
Pflicht-Report (`tools/schema/plan.yaml`) wird dabei geschrieben und trägt
die Ziel-DSN des Laufs. Der Unterschied zum Stand davor ist diese eine
Ziel-Zeile (`cdc-test-postgres:5432/cdc`); das Rollback-Artefakt
`tools/schema/down.sql` bleibt unverändert. Die Datei reist mit dem Lauf
mit, wie in den vorigen Rollout-Commits.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-066` **und** `slice-067` liegen
in `done/` (der Rundlauf braucht Antrag, Verarbeitung und Filterung
zusammen), WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Happy-Path- und Negative-Beleg zusammen mehr als drei Liefer-Punkte
  ergeben (z. B. weil der Compose-Stack-Rundlauf unerwartet umfangreiche
  Hilfsfunktionen im Skript braucht), gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Der historische-Changes-
  Teil des Happy-Path-Belegs lässt sich am laufenden Compose-Stack nicht
  sauber von den bestehenden Schema-Evolution-/Walsender-Testabschnitten
  isolieren (Zustands-Überschneidung, `BEO-PGC/test-isolation-geteilter-
  zustand`-Muster) — dann Carveout statt eines flackernden Tests.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make test-integration` real grün (voller
Compose-Stack-Lauf, beide neuen Abschnitte enthalten) **und** `make gates`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der bestehende Compose-Stack-Lauf ist bereits lang (mehrere
  Administrations-, Schema-Evolution- und Walsender-Rundläufe im selben
  Skript); ein weiterer Abschnitt könnte eine bereits bekannte Timing-
  Empfindlichkeit verschärfen (`BEO-PGC/test-integration-retention-
  timing-flake`, 1×, weiter offen — anderer Rundlauf-Bereich, aber
  dieselbe Skript-Familie und dasselbe Poll-Muster). — **Ausgang:
  entfallen** — im realen Lauf nicht eingetreten; über vier unabhängige
  vollständige Läufe (Implementer, Reviewer 2×, Verifier) kein Flake.
  `BEO-PGC/test-integration-retention-timing-flake` bleibt unverändert bei
  1× (kein zweiter Beleg).
- Der historische-Changes-Teil des Happy-Path-Belegs setzt voraus, dass
  ein vor dem Ausschluss erfasster Change im selben Testlauf bereits
  existiert, ohne durch einen späteren Cleanup-Schritt eines anderen
  Rundlauf-Abschnitts überschrieben zu werden (Musterrisiko wie
  `BEO-PGC/test-isolation-geteilter-zustand`, 1×, weiter offen). —
  **Ausgang: entfallen** — der Abschnitt liest die Zeile unmittelbar nach
  ihrer Erfassung, die Retention-Schwelle liegt bei 24 h und die
  Retention-Abschnitte fassen eine andere Tabelle an; die Auslegung des
  „historischen" Teils ist zusätzlich in §1 korrigiert (F-1 des Reviews:
  die Anforderung sagt nur „künftige Changes" zu).

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

- **Was hat funktioniert:** Der Rundlauf nutzt die bestehende
  Antrags-Queue-Mechanik unverändert — der reale Beleg am laufenden
  Container brauchte keinen Neustart und keine neue Infrastruktur. Zwei
  Konstruktionen machen den Beleg belastbar: eine **Baseline** (Zeile *vor*
  dem Ausschluss, Sentinel im `new_data`) schließt die Alternativerklärung
  „die Spalte war nie im Row Image" real aus, und die **Kontrollspalte**
  `name` trennt gezielten von totalem Filter. Implementer, Reviewer und
  Verifier haben je eigene Mutationen gesetzt und rot gesehen; die zwei
  Hälften (Happy Path, Negative) sind einzeln tragend.
- **Was ging anders als geplant:** (a) Der §1-Satz des Plans sprach von
  „künftige **und historische** Changes" — der Historien-Teil ist unter
  `ADR-0059` Teilfrage 3 Option D unerfüllbar (Filterung zur Bauzeit,
  `cdc.changes` ist reine Projektion); der Reviewer stufte das als Finding
  gegen den **Text** ein, der Planner hat §1 korrigiert. (b) Die Platzierung
  wanderte vom SQL-Administrations-Abschnitt weg hinter das Go-E2E-Tier
  (die eigene Tabelle wird selbst erst über die Queue aktiviert) — im
  §3-Nachzug begründet, kein Scope-Creep.
- **Steering-Loop-Eintrag:** keiner neu verkörpert. Zwei Beobachtungen
  wurden in dieser Closure **benannt statt gezählt** (kein eigener Beleg,
  weil beide keinen eigenen abgeschlossenen Vorgang tragen): (1) das
  wiederholte Auftreten der Sequenzierungs-Klasse aus `AGENTS.md` §3.9 beim
  Planner-Koordinator während des `open→next`-Übergangs dieses Slice
  (Gate rot gesehen, Folgehandlung trotzdem ausgeführt) — im Registereintrag
  `BEO-PGC/report-nackte-id-ohne-link` vermerkt; (2) der Konventionsfehler
  im Architect-Verdikt, einen Slice-Pfad fest mit seiner Lifecycle-Ablage
  zu verlinken (bricht beim nächsten Übergang) — Repo-Konvention ist die
  Kennungs-Zitierung.
- **Beobachtungs-Register (`../observations/`):** kein neuer Eintrag, kein
  neuer Beleg aus diesem Slice selbst — die zwei obigen Punkte sind
  benannt, nicht gezählt.
- **Folge-Slices:** `slice-075` (dauerhafter Ausschlussstand, `open/`,
  wellenlos) aus dem Architect-Verdikt zu `slice-067` — **nicht** Teil
  dieser Welle.
- **Risiken aus §6:** beide entfallen — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-18` offen) —
  Prüfung läuft bei der `welle-18`-Closure.

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
Treffer mit Bezug zur Test-Integration-Skript-Familie:
`BEO-PGC/test-integration-retention-timing-flake` (1×, weiter offen —
anderer Rundlauf-Bereich, siehe §6), `BEO-PGC/test-isolation-geteilter-
zustand` (1×, weiter offen — Musterrisiko für geteilten Testzustand,
siehe §6). Keiner der beiden erreicht mit diesem Slice 3× erstmals; beide
bleiben unter der Schwelle offen. Keine weiteren Treffer für
`Assembler`/`TableBinding`/Antrags-Queue/Schema-Evolution über die bereits
in `slice-066`/`slice-067` gesichteten hinaus.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
