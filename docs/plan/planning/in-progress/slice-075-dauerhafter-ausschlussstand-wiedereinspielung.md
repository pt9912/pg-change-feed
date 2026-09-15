# Slice slice-075: Dauerhafter Ausschlussstand — Wiedereinspielung über Neustart und Bindungs-Zyklus

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist ausschließlich die eigene
DoD (ein einzelner Slice, der eine im Review von `slice-067` gefundene Lücke
schließt), kein repo-weites *Mehr* wie bei `welle-13`/`welle-14`/`welle-15`
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).
Ausdrücklich **nicht** Teil des `welle-18`-Closure-Triggers — `welle-18`
schließt über den realen E2E-Beleg aus `slice-068`, siehe
`docs/reviews/architect-verdict-spaltenausschluss-dauerhaftigkeit.md`.

**Bezug:** [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Haupt-Bezug —
Spaltenausschluss), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (deren
Zusage hängt an der Dauerhaftigkeit des Ausschlussstandes),
[`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
(bindend — entscheidet den dauerhaften, tabellen-scoped Träger),
[`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) (Mechanismus,
Wirkort und Rückkanal — in der Dauerhaftigkeits-Aussage durch `ADR-0065`
superseded, in allen übrigen Teilfragen unverändert gültig),
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue als einziger Schreibpfad administrativer Zustände),
[`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt der
neuen Lesefähigkeit).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md)
(Antrags-Datensatz — seine `applied`-Zeilen sind die Herkunft des Standes),
[`ARC-004`](../../../../spec/architecture.md) (die Spalten-Prüfungs-Fähigkeit
bekommt eine Lesefähigkeit für den Ausschlussstand).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Architect-Zug — dieser Plan trägt die Adresse der
Verdikt-Auflage aus
`docs/reviews/architect-verdict-spaltenausschluss-dauerhaftigkeit.md`; die
endgültige Planung führt der Planner, einschließlich der Wellen-Zuordnung).
**Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Den Ausschlussstand über die Prozesslebensdauer hinaus tragen
(`ADR-0065`): den Stand einer Tabelle aus den `applied`-Zeilen der beiden
Spalten-Antragsarten in `cdc.administration_request` ableiten
(`exclude_column` trägt ein, `include_column` nimmt heraus, Reihenfolge
`requested_at` mit deterministischem Zweitschlüssel), ihn über eine neue
Lesefähigkeit am Outbound Port (`ARC-004`) bereitstellen und bei **jedem**
Anlegen einer Bindung mitführen — im Prozessstart
(`activatedTableBindings`) **und** im Aktivierungs-Zweig der
Antrags-Verarbeitung (`AddBinding`). Belegt wird das durch einen Unit-Test
über den Bindungs-Neuaufbau und einen `disable`/`enable`-Zyklus sowie durch
einen realen E2E-Beleg: ein simulierter Container-Neustart im bestehenden
Rundlauf (`tools/harness/run-integration-tests.sh`) lässt einen zuvor
beantragten Ausschluss wirksam.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Boot-Zeit-Ausschlussfeld in `CDC_TABLES`/der optionalen YAML-Datei** —
  `ADR-0059` Teilfrage 1 erklärt es für nicht ausgeschlossen, aber nicht zum
  Gegenstand; `ADR-0065` führt es als eigenen Re-Evaluierungs-Trigger (ein
  zweiter Herkunftsort des Standes bräuchte eine eigene Zusammenführungs-
  Festlegung). Es gibt keinen benannten Bedarf dafür.
- **Quellen-/musterweiter Ausschluss über mehrere Tabellen hinweg** —
  `ADR-0059` Teilfrage 2 schließt Option C bewusst aus; dieser Slice belegt
  weiterhin ausschließlich die pro-Tabelle-Granularität.
- **Ausschluss zusätzlich als Publication-Spaltenliste** — `ADR-0065`
  §Verglichene Alternativen (Option C) verwirft den zweiten Wirkort: er
  verdoppelte die Auswertung (SQL und `Assembler`) und berührte die
  Konvergenz mit `LH-FA-SCH-003`.
- **Wiederholung des Live-Reload-Belegs aus `slice-068`** — dieser Slice
  ergänzt den Neustart- und den Bindungs-Zyklus-Beleg; der Beleg des
  laufenden Pfads bleibt bei `slice-068` (`welle-18` §3), ein zweiter Lauf
  desselben Pfades prüfte nichts Neues.
- **Eine dritte Zustands-Kategorie „ausgeschlossen" vs. „real gelöscht"** —
  `ADR-0059` Teilfrage 4 (Option B) bleibt verworfen; die Boundary-Klausel
  von `LH-FA-CFG-005` verlangt ausdrücklich nur das Verhalten aus
  `LH-FA-SCH-003`.

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

- [x] Der Ausschlussstand einer Tabelle wird aus den `applied`-Zeilen der
      beiden Spalten-Antragsarten abgeleitet und bei jedem Anlegen einer
      Bindung mitgeführt — im Prozessstart (`activatedTableBindings`) und im
      Aktivierungs-Zweig (`AddBinding`). Beleg: Test in
      `internal/bootstrap/`, der den Bindungs-Neuaufbau und einen
      `disable`/`enable`-Zyklus real nachbildet und danach prüft, dass der
      ausgeschlossene Spaltenname im Row Image fehlt
      (`TestActivatedTableBindingsCarriesExcludedColumns`,
      `TestProcessAdministrationRequestsDisableEnableCycleRestoresExclusion`,
      dazu `TestProcessAdministrationRequestsMarksFailedWhenExclusionReadFails`);
      `make test` Exit 0, beide Zusagen einzeln rot gesehen (Mutationen 1/2).
- [x] Die Ableitung liegt hinter einer neuen Lesefähigkeit am Outbound Port
      (`ARC-004`, Fähigkeits-Zuschnitt nach `ADR-0034`), gegen die reale
      PostgreSQL erprobt (`make test-store`) — inklusive der deterministischen
      Reihenfolge bei gleichem `requested_at`. Beleg:
      `TestTableActivationExcludedColumnsDerivesAppliedColumnRequests`,
      `make test-store` Exit 0; Zweitschlüssel und `include_column`-Wirkung
      einzeln rot gesehen (Mutationen 3/4).
- [x] E2E-Beleg (`LH-QA-SEC-004`): ein simulierter Container-Neustart im
      bestehenden Rundlauf (`tools/harness/run-integration-tests.sh`) lässt
      einen zuvor per `cdc.exclude_column` beantragten Ausschluss wirksam —
      ein danach eingefügter Change trägt den Spaltenwert nicht in
      `cdc.changes`. Beleg: `make test-integration` Exit 0 (Zeile
      *Spaltenausschluss-Neustart-Beleg* in der Ausgabe).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-075.md`, `.harness/skills/reviewer.md`) —
      Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6),
      kein Self-Review (Modul 8).
- [x] Doku-Update: `harness/README.md` §Werkzeuge, Zeile `make
      test-integration`, um den neuen Beleg-Baustein ergänzt; dazu der
      `SPEC-019`-Fließtext zur Bedeutung des `applied`-Wertes (Folgepflicht
      aus `ADR-0065`, Planner-/Architect-Zug) plus Historie-Zeile in §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `PGC`), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Erwartet: `evidence/slice-075.md` in `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/` (der Eintrag trägt diesen Slice als Auslöser, siehe §8).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der Closure der Welle, der dieser Slice zugeordnet wird.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/columnexclusion.go` | update | neue Lesefähigkeit für den abgeleiteten Ausschlussstand einer Tabelle (Fähigkeits-Zuschnitt, `ADR-0034`) |
| `internal/adapters/driven/postgresstorage/tableactivation.go` | update | Implementierung der Lesefähigkeit gegen `cdc.administration_request` (dieselbe Instanz wie die Spalten-Prüfung) |
| `internal/adapters/driven/postgresstorage/queries.go` | update | Abfrage der `applied`-Zeilen der beiden Spalten-Antragsarten in `requested_at`-Ordnung mit deterministischem Zweitschlüssel |
| `internal/bootstrap/wiring.go` | update | `activatedTableBindings` und der Aktivierungs-Zweig tragen den abgeleiteten Stand in die Bindung |
| `internal/bootstrap/*_test.go` | update | Belege für Bindungs-Neuaufbau und `disable`/`enable`-Zyklus |
| `tools/harness/run-integration-tests.sh` | update | Neustart-Beleg des Ausschlussstandes im bestehenden Rundlauf |
| `harness/README.md` | update | Werkzeuge-Zeile `make test-integration` um den neuen Beleg-Baustein ergänzt |

**Plan-Nachzug (Implementer, 2026-09-15) — Pfad der SQL-Texte.** Die Tabelle
nennt `internal/adapters/driven/postgresstorage/queries.go`; die SQL-Texte
liegen im Unterpaket (`internal/adapters/driven/postgresstorage/queries/queries.go`,
`ADR-0042`-Paketstruktur). Die neue Abfrage steht dort, nicht im Adapter.

**Plan-Nachzug (Implementer, 2026-09-15) — Form der Lesefähigkeit.** Die
Fähigkeit liest **je Quelle** (`ExcludedColumns(ctx, source) (map[string][]string, error)`,
Schlüssel `schema.table`) statt je Tabelle: `ADR-0065` beziffert den
Startpfad-Zusatzaufwand ausdrücklich mit „eine Abfrage je Quelle, nicht je
Tabelle", und beide Aufrufer — `activatedTableBindings` (alle Tabellen der
Quelle) und der Aktivierungs-Zweig (eine Tabelle) — bedient dieselbe Rückgabe.
Der Ausschlussstand einer Tabelle ohne geführten Namen trägt **keinen**
Map-Eintrag, keine leere Liste.

**Plan-Nachzug (Implementer, 2026-09-15) — drei weitere Test-Orte.** Über
`internal/bootstrap/*_test.go` hinaus: (a)
`internal/adapters/driven/postgresstorage/administrationrequest_test.go`
trägt den realen PostgreSQL-Beleg der Ableitung samt deterministischem
Zweitschlüssel (DoD-Punkt 2, `make test-store`); (b)/(c)
`internal/application/usecase/excludecolumn/service_test.go` und
`.../includecolumn/service_test.go` tragen die neue Methode in ihren
`ColumnExclusionPort`-Fakes nach — die Fähigkeit sitzt am bestehenden
Port-Zuschnitt (`ADR-0034`), ihre Fakes müssen den erweiterten Vertrag
erfüllen, ohne ihn zu nutzen.

**Plan-Nachzug (Implementer, 2026-09-15) — `spec/pflichtenheft.md`.** Der
DoD-Punkt *Doku-Update* nennt neben `harness/README.md` den
`SPEC-019`-Fließtext zur Bedeutung von `applied`; die §3-Tabelle führte ihn
nicht. Der Satz steht jetzt in `SPEC-019` (dauerhaft vermerkter Stand, eine
Herkunft, Wirksamkeit sobald die Tabelle erfasst wird) samt Zeile in §7
Historie. **Abweichung zum ADR-Wortlaut, benannt:** `ADR-0065`s Folgepflicht
etikettiert diesen Text als „Planner-/Architect-Zug"; der Slice-DoD ordnet
ihn diesem Lauf zu. Inhaltlich entscheidet der Slice nichts — er schreibt den
in `ADR-0065` festgelegten Wortlaut der Bedeutung aus.

**Plan-Nachzug (Implementer, 2026-09-15) — `harness/image-hash.txt`.** Der
Diff ändert Build-Kontext-Dateien (`internal/**`); `compose.yaml`
referenziert das lokal gebaute Image ohne eigenen `build:`-Block, `make
test-integration` lief deshalb erst nach einem `make image` gegen den neuen
Digest (der alte Lauf traf noch den Vorstand und färbte den Neustart-Beleg
rot — siehe Bericht). Der neue Digest reist mit dem Diff mit
(`ADR-0044`).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `ADR-0065` liegt Accepted vor;
`slice-067` liegt in `done/` (sein §6-Risiko trägt den Ausgang
*eingetreten → slice-075*, der Planner-Zug zu diesem Risiko-Ausgang ist
gelaufen); WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Ableitung, Lesefähigkeit und die zwei Beleg-Tiers zusammen mehr als drei
  Liefer-Punkte oder mehr als zwei Schichten in einer Review-Sitzung nicht
  mehr prüfbar machen, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Die Antrags-Historie trägt
  den Stand nicht (z. B. weil eine Bereinigungs- oder Altersgrenze auf
  `cdc.administration_request` liegt oder die Reihenfolge nicht
  deterministisch herstellbar ist) — dann greift der Re-Evaluierungs-Trigger 1
  aus `ADR-0065`: der Träger wird per Folge-ADR neu entschieden, statt hier
  ein unvollständiges Modell zu bauen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration` real
grün (der Neustart-Beleg ist darin enthalten) **und** Closure-Notiz
geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `requested_at` trägt `current_timestamp` (Transaktionszeit): ein
  `exclude_column` und ein `include_column` für dieselbe Spalte in **einer**
  Transaktion tragen denselben Zeitstempel — die Ableitung braucht einen
  deterministischen Zweitschlüssel, sonst ist die Reihenfolge zufällig und der
  Stand hängt von der Ausführungsreihenfolge ab. — **Ausgang: entfallen** —
  `administration_request_id` (Primärschlüssel) ist der Zweitschlüssel; der
  Verifier hat die Mutation „`ORDER BY` ohne ihn“ real rot gesehen
  (`make test-store`, Exit 2, Reihenfolge kippte auf Einfüge-Reihenfolge).
- Der Prozessstart liest eine weitere Quelle je Quelle; ein Lesefehler dort
  endet in der Startfehlerklasse des bestehenden Pfads (`storage`,
  `SPEC-008`) — eine neue Startabbruch-Bedingung, die der Slice benennen muss.
  — **Ausgang: entfallen** — der Slice hat sie benannt: ein Lesefehler
  führt in den bestehenden Startabbruch-Pfad (`storage`); Review F-2 hat die
  Form als vorbestehend und fail-closed bestätigt, per Test gepinnt.
- Der Neustart-Beleg teilt Zustand mit den übrigen Abschnitten des langen
  Compose-Rundlaufs (`BEO-PGC/test-isolation-geteilter-zustand`, 1×, weiter
  offen — Musterrisiko für geteilten Testzustand) und trifft dieselbe
  Poll-Familie wie `BEO-PGC/test-integration-retention-timing-flake` (1×,
  weiter offen). — **Ausgang: entfallen** — über sechs unabhängige
  vollständige Läufe (Implementer, Reviewer, Verifier) kein Flake und keine
  Zustandsüberschneidung; beide Registereinträge bleiben bei 1× (kein
  zweiter Beleg).

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

- **Was hat funktioniert:** Die in `ADR-0065` entschiedene Ableitung ließ
  sich ohne neues Schema-Objekt und ohne zweiten Schreibpfad umsetzen: eine
  neue **Lesefähigkeit** am bestehenden `ColumnExclusionPort`, ein `SELECT`,
  zwei Aufrufer (Prozessstart und Aktivierungs-Zweig) — und beide Auslöser
  (Neustart und `disable`/`enable`-Zyklus) fallen auf **denselben**
  Mechanismus zusammen. Der reale Beleg ist stark: Der Implementer sah den
  Neustart-Beleg zunächst **rot**, weil sein Lauf gegen das Vorstands-Image
  lief; der Verifier hat die Nicht-Vakuumität mit einer eigenen Mutation am
  E2E-Tier bestätigt — dabei blieb der Live-Reload-Beleg **grün**, nur der
  Neustart-Beleg kippte. Die Trennung beider Zusagen ist damit real gezeigt.
- **Was ging anders als geplant:** (a) Der erste Compose-Lauf lief gegen das
  Vorstands-Image (`make image` fehlte davor) und lieferte damit versehentlich
  den Rot-Beleg — im Plan-Nachzug benannt. (b) Der `SPEC-019`-Fließtext wurde
  in diesem Lauf ergänzt; `ADR-0065` etikettiert ihn als
  „Planner-/Architect-Zug" — die Abweichung ist benannt, der Review hat sie
  als tragfähig bestätigt und einen eigenen `SPEC-*`-Eintrag für **nicht**
  nötig gehalten (die Fähigkeit ist eine interne Go-Schnittstelle). (c) Der
  neue Beleg-Baustein in `harness/README.md` trug als einziger keinen
  `· seit slice-075`-Vermerk (Review F-1) — hier nachgezogen.
- **Steering-Loop-Eintrag:** keiner neu verkörpert — der Slice führt keine
  wiederkehrende Klasse ein. Benannt statt gezählt: die Nicht-Vakuum-Aussage
  des Neustart-Belegs stützt sich auf einen Rot-Lauf, dessen Protokoll nicht
  abgelegt ist (Review F-3); tragend sind die **strukturelle** Begründung
  (der Elternstand liest keinen Ausschlussstand) und der eigene Rot-Beleg des
  Verifiers.
- **Beobachtungs-Register (`../observations/`):** kein **neuer** Beleg —
  `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` war der Träger dieser
  Arbeit; sein Ausgang wechselt mit dieser Closure von `geplant` auf
  **`verkörpert`** (Zielort ist der gelieferte Mechanismus, Anker
  `· seit slice-075`). **Kein** zusätzliches `evidence/`-Dokument dort: der
  Slice ist die *Behebung*, kein zweites Auftreten („ein Vorgang zählt
  einmal"), ein zweiter Beleg würde den Zähler falsch auf 2× heben.
  `evidence/slice-075.md` in `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/`>
- **Folge-Slices:** keiner.
- **Risiken aus §6:** alle drei entfallen — siehe §6.
- **Drei Paarungen:** Repo **ohne** laufende Welle (die Roadmap ist leer;
  dieser Slice ist wellenlos) — hier geprüft: **Anker** — kein
  `liegt in`-Feld in dieser Closure, Paarung entfällt. **Folge-Slice** —
  keiner genannt. **Register** — der zitierte Eintrag
  `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` existiert mit nicht
  leerem `evidence/`.
  Closure der Welle, der dieser Slice zugeordnet wird.

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
Treffer mit Bezug zu diesem Slice: `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`
(1×, offen — **dieser Slice ist sein Träger**; die Beobachtung bleibt bis zur
Umsetzung unter der Schwelle offen), `BEO-PGC/test-isolation-geteilter-zustand`
(1×, weiter offen — siehe §6) und `BEO-PGC/test-integration-retention-timing-flake`
(1×, weiter offen — siehe §6). Keiner der drei erreicht mit diesem Slice 3×.
Keine weiteren Treffer für `Assembler`/`TableBinding`/Antrags-Queue über die
bereits in `slice-066`/`slice-067`/`slice-068` gesichteten hinaus
(`BEO-PGC/schema-evolution-nicht-dynamisch` bleibt verkörpert, ohne neuen
Bezug).

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
