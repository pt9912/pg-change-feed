# Welle transformationen: Transformationen — erfasste Changes tragen vor der Persistierung die durch deklarative Regeln bestimmte Form, konfiguriert über die SQL-Antrags-Queue (`LH-FA-CFG-007`, `ADR-0112`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-transformationen-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** — (Rolleninhaber der Implementer-Rolle je Slice, gesetzt
beim Übergang `open` → `next`; geschnitten aus
[`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
§Folgepflichten durch den Planner). **Datum:** 2026-09-23.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Eine Tabelle trägt Transformationsregeln, und **jede** danach erfasste Change
trägt die transformierte Form — nicht die Rohform — an jedem Weg, der sie
liefert: SQL-Lesezugriff `cdc.changes`, `GET /changes`, gRPC-Stream, SSE-Stream
und NATS-Vollinhalts-Stream; eine auf eine Change nicht anwendbare Regel endet
sichtbar statt still. Das ist die Aussage der drei Akzeptanzkriterien von
[`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Happy Path, Boundary,
Negative). Der Mechanismus ist mit
[`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
entschieden: ein geschlossener Satz deklarativer Regeltypen (`rename_column`,
`map_value`), konfiguriert über zwei neue Antragsarten der Antrags-Queue
(`cdc.set_transformation`, `cdc.remove_transformation`), dauerhaft aus den
`applied`-Zeilen abgeleitet, im `Assembler` bei der Row-Image-Konstruktion
**nach** dem Spaltenausschluss und **vor** jeder Serialisierung und
Persistierung ausgewertet.

Das *Mehr* gegenüber den zehn Slice-DoDs: keiner der Slices belegt eines der
drei Kriterien allein. Der Happy Path braucht die Domäne, den `Assembler`, den
Antragsweg mit Dauerhaftigkeit und alle Zustellwege am laufenden
Feed-Container; die Boundary (Mehrdeutigkeit statisch ausgeschlossen, K1–K4)
liegt im Use Case und wird erst am komponierten System als `failed` mit Text
sichtbar; die Negative braucht die Prüfung im `Assembler`, die
Klassen-Abbildung, `diagnose`, die Startreihenfolge und die Abhilfe — das
**Abhilfe-Akzeptanzkriterium** (a)–(d) aus
[`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 ist ein Zustand des ganzen Bündels. Dazu kommt der
Träger-Zustand, den kein Einzel-DoD beobachtet: die Regelauswertung liegt an
genau **einer** Stelle (die gemeinsame Row-Image-Funktion, aufgerufen vom
WAL-Pfad **und** vom Backfill-Pfad), und
[`LH-FA-CFG-007`](../../../spec/lastenheft.md) ist im RTM-Lauf (`make
doc-trace`) nicht mehr Waise.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  trägt den Status `Accepted` — **bereits erfüllt**, geprüft an der Zeile
  `ADR-0112` im ADR-Index ([`docs/plan/adr/README.md`](../adr/README.md),
  Status-Spalte `Accepted`, Datum 2026-09-23).
- Das Lastenheft trägt [`LH-FA-CFG-007`](../../../spec/lastenheft.md) auf
  Transformationen begrenzt und das Routing als eigene Anforderung
  [`LH-FA-CFG-008`](../../../spec/lastenheft.md) — **bereits erfüllt**, geprüft
  an `spec/lastenheft.md` (Kopf: Version 0.13.0; Out-of-Scope-Satz von
  `LH-FA-CFG-007` nennt `LH-FA-CFG-008`).
- Kein weiterer Trigger nötig — die Welle kann sofort eröffnet werden. Die
  Kopplungen K1–K3 zur Welle
  [welle-backfill-bestand](done/welle-backfill-bestand.md) sind **Start-Trigger
  einzelner Slices** (§4, §5), kein Trigger der Welle: sie ist als Ganzes
  planbar, ihre Slices starten gestaffelt.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle zehn Slices in `done/`.
- `make gates` grün — der Exit-Code des Laufs wird ungefiltert gesichert und
  gesondert ausgewertet ([`AGENTS.md`](../../../AGENTS.md) §3.9).
- **Ein realer, grüner `make test-integration`-Lauf** mit den Belegen am
  laufenden Feed-Container — Happy Path (die per SQL beantragte Regel prägt die
  danach erfasste Change über `cdc.changes`, `GET /changes`, gRPC, SSE und
  NATS-Vollinhalt in derselben Form), Boundary (eine K1–K4-Verletzung endet
  `failed` mit Text, der Regelstand bleibt), Neustart-Festigkeit (`docker
  restart`), Ausschluss+Regel ([`LH-QA-SEC-004`](../../../spec/lastenheft.md):
  weder Quell- noch Zielname noch Wert im Image), ein Backfill-Bestand mit
  transformierter Form und das **Abhilfe-Akzeptanzkriterium (a)–(d)** aus
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 5 (die Regel wird nichtanwendbar, der Prozess endet mit Klasse
  `schema`; `cdc.remove_transformation` wird beantragt, während er steht; nach
  dem Neustart ist der Antrag `applied`, bevor die erste Transaktion der
  Tabelle assembliert wird; die zuvor nicht bestätigte Transaktion erscheint
  über `cdc.changes`) — kein Gate, aber Pflichtbeleg dieser Welle.
- **Reale, grüne Läufe der DB-Tiers:** `make test-store` (Regelstand-Ableitung
  aus den `applied`-Zeilen, Antrags-Funktionen, Grants), `make
  test-replication` (der `Assembler` und die gemeinsame Bild-Funktion am realen
  Stream; die DB-Adapter-Coverage prüft `make test-replication` gegen ihre
  Schwelle) und `make schema-rollout` zweimal hintereinander gegen dieselbe
  Ziel-Datenbank (Idempotenz) — die Fitness-Function-Zeilen von
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
  die nur das gebündelte System zeigt. Dazu der **Alt-Tag-Lauf** von
  `tools/harness/run-schema-rollout-guard-test.sh`
  ([`ADR-0114`](../adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 7): das Schema des jüngsten `v*`-Tags ausrollen, danach den
  Arbeitsbaum — Exit 0 zweimal, der Datenstand über `cdc.changes` lesbar; er
  zeigt die zwei nullable Spalten und die zwei Funktionen dieser Welle über einen
  Alt-Bestand (die Welle ändert keine bestehende View, siehe den
  §3.13-Suchlauf in §5 von
  [welle-backfill-bestand](done/welle-backfill-bestand.md)).
- **Die Fitness-Function-Tests aus
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  sind grün und tragen ihre Eingabe:** der Eigenschaftstest im `mapper`-Paket
  (Regeltyp × ausgeschlossene Spalte, beide Regeltypen aus der Domänen-Menge),
  der Determinismus-/Nebenläufigkeits-Test unter `-race`, `make a-check` (die
  Regeltypen in `internal/domain/**`).
- `make doc-trace` führt [`LH-FA-CFG-007`](../../../spec/lastenheft.md) nicht
  mehr unter den Waisen; der Träger ist die Zeile in
  [`docs/user/e2e-abdeckung.md`](../../user/e2e-abdeckung.md), die der Runner
  von `make test-integration` schreibt (die Datei ist ein Erzeugnis, kein
  Lauf-Beleg). [`LH-FA-CFG-008`](../../../spec/lastenheft.md) bleibt als Waise
  sichtbar — das Routing hat keine Umsetzung (§6).
- Die Regelauswertung liegt an genau einer Stelle: der Suchlauf über
  `internal/**` nach der Row-Image-Konstruktion und nach den Aufrufern der
  Regelauswertung (Befehl und Fundstellen im Closure-Bericht) findet eine
  Auswertungsstelle, die WAL-Pfad und Backfill-Pfad gemeinsam aufrufen.
- Das Benutzerhandbuch trägt den Abschnitt zur Regel-Konfiguration samt
  Rohform-Konsequenz (eine Regeländerung wirkt nicht rückwirkend),
  Abhilfe-Prozedur und Fehlerklassen-Zeile — jede Zahl und Wirkungs-Aussage mit
  ihrem Ursprung ([`AGENTS.md`](../../../AGENTS.md) §3.12); der SDK-Beleg
  (Row-Image-Schlüssel sind in den drei SDK-Modellen opak) steht im Bericht von
  `slice-transformationen-betriebsdoku`.
- Der **Lese-Schritt** des Beobachtungs-Registers ist gelaufen (Einträge bei 3×
  oder darüber, Modul 6).
- Closure-Notiz in `welle-transformationen-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-transformationen-spec-nachzug | Pflichtenheft (`LH-FA-CFG-007.a` beantwortet, `SPEC-019`, Regelform, `SPEC-002`, §4-Zeile `schema`) und Architektur-Sicht auf den beschlossenen Stand ziehen — ohne ADR-/Slice-Bezug | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`LH-FA-CFG-007.a`](../../../spec/pflichtenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 1 |
| slice-transformationen-kern-rename | Regeltyp `rename_column` in der Domäne; Regelstand, Set/Remove und Nichtanwendbarkeits-Prüfung im `Assembler`; Fitness-Tests (Eigenschaft, Determinismus, `-race`, a-check) | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 2 |
| slice-transformationen-antragsweg-schema | Antragsarten `set_transformation`/`remove_transformation`, Spalten `rule_name`/`rule_spec`, zwei SQL-Funktionen, Grants, Idempotenz-Guard, Rollen-Test, Domäne und Store-Lesen | [`LH-FA-ADM-001`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 3 |
| slice-transformationen-antragsweg-usecase | Use Cases `SetTransformation`/`RemoveTransformation` mit K1–K4, Regelstand-Port, Store-Adapter, Verdrahtung in `applyAdministrationRequest`, `activatedTableBindings` und Aktivierungs-Zweig | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 3, [`ADR-0065`](../adr/0065-spaltenausschluss-dauerhafter-traeger.md) |
| slice-transformationen-backfill-pfad | Regelauswertung im Backfill-Run, Fail-closed um den Regelstand, Nichtanwendbarkeit im Run, E2E-Beleg | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`LH-FA-CAP-009`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 7 |
| slice-transformationen-map-value | Regeltyp `map_value` in der Domäne — ohne Änderung an Antragsweg und Wirkort | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 4 |
| slice-transformationen-e2e-wirkung | `make test-integration`: Regel per SQL, alle Zustellwege, Boundary, Neustart, Ausschluss+Regel; [`LH-FA-CFG-007`](../../../spec/lastenheft.md) verlässt die Waisen | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5 |
| slice-transformationen-start-reihenfolge | offene Anträge werden vor `stream.Run` verarbeitet — Ordnung deterministisch und ohne Datenbank prüfbar | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5, Bedingung (c) |
| slice-transformationen-e2e-abhilfe | Nichtanwendbarkeit und das Abhilfe-Akzeptanzkriterium (a)–(d) am laufenden System | [`LH-FA-CFG-007`](../../../spec/lastenheft.md) Negative, [`LH-FA-ADM-003`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5 |
| slice-transformationen-betriebsdoku | Benutzerhandbuch: Regel-Konfiguration, Rohform-Konsequenz, Abhilfe, Fehlerklasse; SDK-Beleg (Bild-Schlüssel opak) | [`LH-FA-CFG-007`](../../../spec/lastenheft.md), [`LH-FA-SST-009`](../../../spec/lastenheft.md), [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 6 |

**Reihenfolge:** sequentiell in der Tabellen-Reihenfolge (WIP-Limit 1 je
Rolleninhaber, Baseline-Regelwerk `modul-05-planning-harness.md`); die Slices
dieser Welle und die der Welle
[welle-backfill-bestand](done/welle-backfill-bestand.md) teilen sich das WIP-Limit,
und jeder Slice startet, sobald sein Start-Trigger gilt — die technischen
Kanten stehen in §5. Die Reihenfolge der Tabelle folgt einer Regel: **erst die
Wirkung, dann der Beleg, dann die Doku** — und der Backfill-Pfad **unmittelbar
nach** dem Antragsweg, damit das Fenster, in dem Regeln setzbar sind, ein
Backfill aber noch die Rohform liefert, ein Slice lang ist (§5).

**Abweichungen vom Schnitt-Vorschlag** in
[`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
§Folgepflichten (dort ausdrücklich „der Planner formt daraus Welle und Slices“)
— fünf, je mit Grund:

1. **Folgepflicht 3 geteilt** in `antragsweg-schema` (Schema-Rollout, zwei
   Funktionen, Grants, Guard, Rollen-Test, Domäne und Store-Lesen — Tier `make
   test-store`, `make schema-rollout`) und `antragsweg-usecase` (Use Cases mit
   K1–K4, Port, Adapter, Verdrahtung — Tier `make test`, `make test-store`).
   Die Folgepflicht nennt mehr als drei Liefer-Punkte (Schema, Antragsart, Use
   Cases, Port, Verdrahtung); die Backfill-Welle teilte einen vergleichbaren
   Slice aus demselben Grund (Rückführung im Plan von
   `slice-backfill-sql-administration`). Zwischen beiden besteht ein benannter
   Zwischenzustand: die Funktion existiert, ein Antrag endet `failed` mit Text
   (Risiko in `antragsweg-schema` §6).

2. **Folgepflicht 5 geteilt und ein Vorab-Slice davor.** (i) `e2e-wirkung`
   (Happy Path, Boundary, Neustart, Ausschluss+Regel) und `e2e-abhilfe`
   (Nichtanwendbarkeit, Abhilfe (a)–(d)) — verschiedene Szenarien, das zweite
   endet den Prozess und hat eine eigene Container-Ende-Grenze im Runner. (ii)
   „Offene Anträge vor `stream.Run` verarbeiten“ ist ein **eigener Slice**
   (`start-reihenfolge`), nicht Teil des E2E-Slice:
   [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
   legt ihn in den E2E-Slice, **wenn** die heutige Startreihenfolge Kriterium
   (c) nicht trägt. Die Bedingung ist am Code entscheidbar und wird hier
   entschieden — `runAdministration` startet in `Run` per `go func()` vor
   `stream.Run`, zwischen beiden liegt keine Synchronisation (gelesen an
   `internal/bootstrap/wiring.go`); ein E2E-Lauf kann das Fehlen einer Race
   nicht beweisen. Der Zug ändert Produktivcode im Startpfad für **alle**
   Antragsarten und braucht einen eigenen Review und einen ordnungsprüfenden
   Test ohne Datenbank. Das weicht vom ADR-Wortlaut ab
   (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, offen, 1×) und ist
   als **Frage an den Auftraggeber** geführt, nicht als Entscheidung: die
   Alternative, den Zug in `e2e-abhilfe` zu bündeln, ist zulässig, verlängert
   aber den Slice um Produktivcode und trennt die Ordnung von ihrem Beleg
   nicht.

3. **Folgepflicht 6 zusammengelegt, Handbuch-Zug an das Ende.** SDK-Beleg und
   Betriebsdokumentation stehen in **einem** Slice (`betriebsdoku`), zuletzt,
   weil das Handbuch nur beschreiben darf, was `e2e-wirkung` und `e2e-abhilfe`
   belegt haben. Die Betreiber-Oberfläche (zwei SQL-Funktionen) entsteht aber
   schon in `antragsweg-schema`: dieser Slice **schiebt** die Beschreibung mit
   Adresse auf (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, offen, 2×),
   und `betriebsdoku` nennt den aufgeschobenen Gegenstand vollständig in seinem
   §2. Das weicht von der Regel der Welle
   [welle-backfill-bestand](done/welle-backfill-bestand.md) ab, wo der Handbuch-Zug
   im oberflächen-liefernden Slice liegt; Grund: bei den Transformationen liegt
   die Wirkung fünf Slices hinter der Oberfläche.

4. **Folgepflicht 4 (`map_value`) hinter den Antragsweg und den Backfill-Pfad
   gestellt, ihre Spec-Zeile in `spec-nachzug`.** Die Pflicht nennt „Domäne,
   Spec-Zeile, Tests“; die Spec führt beide Regeltypen von Anfang an (Doku
   führt, Modus Greenfield), also trägt `map-value` keine Spec-Zeile mehr. Die
   Reihenfolge beweist die Aussage „ohne Änderung an Antragsweg oder Wirkort“:
   der zweite Typ kommt, nachdem beide stehen, und der Diff-Stat des Slice
   belegt sie.

5. **Folgepflicht 7 (Backfill-Pfad) als eigener Slice, Start strenger als die
   Kopplung K2.** K2 der Welle
   [welle-backfill-bestand](done/welle-backfill-bestand.md) verlangt
   `slice-backfill-run-usecase`; dieser Slice startet zusätzlich erst nach
   `slice-backfill-e2e` (er trägt einen E2E-Beleg im Backfill-Runner) und
   trägt ein **Architect-Kurzverdikt** zur Nichtanwendbarkeit im Run: die
   Run-Fehlerklassen von
   [`ADR-0111`](../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5
   nennen `schema` nicht,
   [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
   Teilfrage 4 beschreibt die Nichtanwendbarkeit nur für den Erfassungspfad —
   eine Entscheidung, keine Auslegung; das Verdikt liegt mit
   [`ADR-0117`](../adr/0117-backfill-run-fehlerklasse-schema.md) vor (Start-Trigger
   im Plan).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Wird blockiert von:** keiner Welle;
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  ist `Accepted`, der Trigger (§2) ist erfüllt. Eine **benannte Kopplung**
  verbindet die Welle mit [welle-backfill-bestand](done/welle-backfill-bestand.md) —
  die Kanten K1–K3 der Welle
  [welle-backfill-bestand](done/welle-backfill-bestand.md) §5 sind hier
  Start-Trigger einzelner Slices:
  - **K1 — Row-Image-Funktion.** `kern-rename` startet erst, wenn
    `slice-backfill-row-image-gemeinsam` in `done/` liegt: die Regelauswertung
    hängt an der **einen** gemeinsamen Funktion, nicht an zwei Bild-Erzeugern.
  - **K2 — Backfill-Pfad.** `backfill-pfad` startet nach
    `slice-backfill-run-usecase` und — strenger — nach `slice-backfill-e2e`,
    nach `antragsweg-usecase`; das Architect-Kurzverdikt (§4, Abweichung 5)
    liegt vor ([`ADR-0117`](../adr/0117-backfill-run-fehlerklasse-schema.md)).
  - **K3 — gemeinsame Zeilen.** `antragsweg-schema` startet nach
    `slice-backfill-sql-administration`: beide Umsetzungen berühren die
    geschlossene `request_kind`-Menge in
    `tools/schema/nacharbeit-administration.sql`,
    [`SPEC-019`](../../../spec/pflichtenheft.md), `applyAdministrationRequest`
    und den Idempotenz-Guard; die Transformationen erweitern die dort stehende
    Menge additiv. Stand am 2026-09-25 (gemessen): `request_kind` trägt fünf
    Werte (`tools/schema/nacharbeit-administration.sql`, Zeile 62), die
    Transformationen erweitern auf sieben; die nächste freie Kennung im
    Pflichtenheft ist `SPEC-030` (`grep -o 'SPEC-0[0-9][0-9]'
    spec/pflichtenheft.md | sort -u | tail -3` druckt `SPEC-027`, `SPEC-028`,
    `SPEC-029`); jeder Slice misst beide Zahlen an seinem Start neu.
  - **Kanten zu zwei wellenlosen Slices.** `slice-harness-suchlauf-nachmessen`
    (Nachmess-Werkzeug für das Suchlauf-Feld der Slice-Pläne) geht
    `kern-rename` voraus: das Feld der zehn Pläne dieser Welle wird mit dem
    Werkzeug nachgemessen. `slice-capture-leerlauf-quellbelege` (Beleg
    „Fehlerschwelle erreicht, Container endet“ im Runner) geht `e2e-abhilfe`
    voraus: beide Slices tragen eine Container-Ende-Grenze im selben Runner.
    Beide sind ohne Welle geführt, weil sie keine Closure-Bedingung tragen, die
    von ihrer DoD verschieden wäre; die Kanten stehen als Start-Trigger in den
    beiden Plänen dieser Welle (§4 dort). Ein dritter wellenloser Slice,
    `slice-capture-transient-wiederholung`, und ein vierter,
    `slice-backfill-speicher-untersuchung`, tragen keine Kante zu dieser Welle.
  - **Spec-Kollision.** `spec-nachzug` startet nach
    `slice-backfill-spec-nachzug`: beide ändern
    [`SPEC-019`](../../../spec/pflichtenheft.md), die Kennungsvergabe im
    Pflichtenheft zählt fortlaufend je Datei.
  - **Start-Reihenfolge und Backfill.** `start-reihenfolge` startet nach
    `slice-backfill-sql-administration`: beide Slices ändern `Run` im Bereich
    der Administrations-Goroutine (Backfill-Zweig, Worker und Start-Abgleich
    stehen dann fest).
  - **Begründete Reihenfolge — Backfill zuerst.** Die Kanten zeigen alle in
    dieselbe Richtung: die Transformationen erweitern Stellen, die der Backfill
    anlegt (Welle [welle-backfill-bestand](done/welle-backfill-bestand.md) §5).
    Beide Wellen sind ohne einander eröffnet; jeder Slice dieser Welle wartet
    auf seine Kante, nicht auf die Closure der Backfill-Welle.
- **Blockiert:** keine Welle. Das Routing
  ([`LH-FA-CFG-008`](../../../spec/lastenheft.md)) hat keine Welle-Datei und
  keine Entscheidung; es hängt nicht an dieser Welle, weil
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 7 es ausdrücklich entkoppelt (§6).
- **Ereignis-Adresse.** „Die Closure dieser Welle“ ist eine Adresse, die
  eintreten kann: die Roadmap führt diese Welle unter *Offene Wellen*, und die
  Welle hat einen Closure-Trigger (§3). Die Paarungen der Slice-Closures (Anker
  · Folge-Slice · Register) laufen bei ihr, nicht bei der Closure der Welle
  [welle-backfill-bestand](done/welle-backfill-bestand.md).
- **Benanntes Zwischenzustands-Fenster.** Ab `antragsweg-usecase` sind Regeln
  setzbar; ein Backfill-Run vor `backfill-pfad` lieferte die Rohform, dieselbe
  Tabelle also zwei Formen. Das Fenster ist ein Slice lang (`backfill-pfad`
  folgt unmittelbar, §4 Reihenfolge); Risiko und Ausgang tragen `kern-rename`
  §6 und `antragsweg-usecase` §6.
- **Intern:** die Ordnung folgt der Tabelle in §4; die technischen Kanten (`A →
  B` heißt: `B` setzt `A` voraus) sind: `spec-nachzug` → jeder übrige Slice
  (die Spec führt); `slice-harness-suchlauf-nachmessen` → `kern-rename` →
  `antragsweg-schema` → `antragsweg-usecase` → `backfill-pfad` → `map-value` →
  `e2e-wirkung` → `start-reihenfolge` → `e2e-abhilfe` → `betriebsdoku` — jede
  Stufe braucht das lauffähige System der Vorstufe (die Regeltypen, der
  Antragsweg, der Backfill-Pfad, beide Typen, die Wirkung, die Ordnung, die
  Abhilfe); dazu die Kante `slice-capture-leerlauf-quellbelege` →
  `e2e-abhilfe` (Belegaufbau der Container-Ende-Grenze im Runner).

**Träger der Folgepflichten** — jede Pflicht aus
[`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
§Folgepflichten hat einen Träger oder eine benannte Adresse
(`BEO-PGC/adr-folgepflicht-ohne-traeger-slice` — eine adresslose Pflicht bleibt
liegen, bis sie ein Leser zufällig findet):

| Folgepflicht | Träger |
|---|---|
| 1 Spec-Nachzug | `slice-transformationen-spec-nachzug`; der Handbuch-Anteil an `betriebsdoku` (§4, Abweichung 3) |
| 2 Kern: Domäne und `Assembler`, `rename_column` | `slice-transformationen-kern-rename` |
| 3 Antragsweg und Dauerhaftigkeit | `slice-transformationen-antragsweg-schema` und `slice-transformationen-antragsweg-usecase` (§4, Abweichung 1) |
| 4 Regeltyp `map_value` | `slice-transformationen-map-value`; die Spec-Zeile trägt `spec-nachzug` (§4, Abweichung 4) |
| 5 E2E-Belege, Abhilfe-Akzeptanzkriterium, Startreihenfolge | `slice-transformationen-e2e-wirkung`, `slice-transformationen-start-reihenfolge`, `slice-transformationen-e2e-abhilfe` (§4, Abweichung 2) |
| 6 SDK-/Doku-Beleg | `slice-transformationen-betriebsdoku` |
| 7 Bindung künftiger Erzeugungspfade | `slice-transformationen-backfill-pfad` für den Backfill-Pfad; Adresse für jeden weiteren Change-Erzeuger: [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 7 selbst — die ADR, die ihn einführt, nennt die Regelauswertung |
| 8 Routing-ADR | kein Slice dieser Welle; Adresse: §6 dieser Welle und die offene Kennung [`LH-FA-CFG-008.a`](../../../spec/pflichtenheft.md) im Pflichtenheft (sichtbar als Waise in `make doc-trace`) |

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Routing** ([`LH-FA-CFG-008`](../../../spec/lastenheft.md)) — eigene
  Anforderung mit eigener Routing-ADR
  ([`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 7 und Folgepflicht 8); die Leitplanken (höchstens ein Ziel je
  Change, erster Treffer, keine Route auf das bestehende Ziel) stehen nur in
  der ADR. Die Welle legt weder ein Zielmodell noch eine Filter-Fähigkeit an
  gRPC/SSE an.
- **Weitere Regeltypen** (`set_constant` u. ä.) — der Satz ist geschlossen; ein
  dritter Typ braucht eine Folge-ADR (Re-Evaluierungs-Trigger: im
  Beobachtungs-Register dreimal auf denselben Bedarf getroffen).
- **Rohform-Speicherung, Rekonstruktion und Rückwirkung** — die Rohform ist
  nicht mehr vorhanden, eine Regeländerung wirkt nur auf künftige Changes, eine
  Rekonstruktion je Change ist nicht zugesagt
  ([`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6 und §Konsequenzen); ein Bedarf ist ein Re-Evaluierungs-Trigger
  mit `Supersedes` für Teilfrage 6.
- **Eine Regel-Sicht (View) und eine Anzeige des Regelstands in `diagnose`** —
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6: der Regelstand ist über `cdc.administration_request` lesbar,
  eine eigene Sicht ist nicht Teil.
- **Regeln je Consumer oder je Zustellweg, Konfiguration über YAML oder
  Umgebungsvariable** —
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 1 (Optionen B, C) und Teilfrage 6 (Option B) sind verworfen.
- **Consumer-Reset** — das Positionsmodell bleibt unverändert; ein Reset ist
  „explizit administrativ“
  ([`ADR-0013`](../adr/0013-consumer-domainkonzept.md)) und eine eigene
  Fähigkeit mit eigener Entscheidung.
- **Ein allgemeiner Recovery-Weg für Schema-Fehler** —
  `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×): die Welle
  liefert die Abhilfe für die Nichtanwendbarkeit einer Transformationsregel,
  nicht die Wiederinbetriebnahme nach einer inkompatiblen Typänderung oder
  einer entfernten Spalte.
- **Änderung des Nachrichtenschemas und der Proto-Artefakte** —
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 8: Live-Nachrichten
  ([`SPEC-020`](../../../spec/pflichtenheft.md)/[`SPEC-021`](../../../spec/pflichtenheft.md)/[`SPEC-024`](../../../spec/pflichtenheft.md))
  tragen zehn Felder, die HTTP-Antwort von
  [`SPEC-022`](../../../spec/pflichtenheft.md) dreizehn, `origin` inbegriffen —
  die Transformationen ändern keines; die Row Images bleiben JSON-Objekte mit
  String-Werten; `make generated-sync` ist nicht betroffen.
- **Kein neuer GitHub-Actions-Workflow und keine strukturelle
  Workflow-Änderung** — [`AGENTS.md`](../../../AGENTS.md) §3.10 greift nicht.
  Die Pipeline `e2e.yml` fährt `make test-integration` unverändert weiter; die
  Laufzeit des erweiterten Testpakets ist ein Risiko der E2E-Slices (§6 dort),
  keine Workflow-Änderung.
- **Kein Server-Release und kein SDK-Release** — `docs/user/version.md` und die
  Package-Versionen der SDKs bleiben durch die Slices dieser Welle unberührt
  (erwartet: keine SDK-Code-Änderung, `betriebsdoku`); die Version des
  Kotlin-Package hebt der wellenlose Slice `slice-sdk-kotlin-cloudsmith`, nicht
  ein Slice dieser Welle; die Änderungshistorie des Benutzerhandbuchs trägt
  eine Zeile.
- **Keine Schwellen-Senkung und keine neue Gate-Klasse** — die neuen
  Use-Case-Pakete liegen in der netzlos gemessenen Fläche des Coverage-Gates;
  eine Senkung bliebe per [`AGENTS.md`](../../../AGENTS.md) §3.6 ADR-pflichtig.
- **Lastenheft unverändert** — die Welle schärft Techniken im Pflichtenheft,
  nicht Anforderungen.

**Beobachtungs-Register — Eröffnungs-Sichtung (Modul 6 Eröffnungs-Schritt 2):**
Das Register (`docs/plan/planning/observations/README.md`) wurde durchgesehen.
Die Zähler sind die Zahl der Dateien unter `evidence/` (real ausgezählt am
2026-09-23 an `docs/plan/planning/observations/BEO-PGC/*/evidence/`). **Kein
Eintrag erreicht mit der Eröffnung dieser Welle 3×** — kein Slice ist allein
wegen einer Beobachtung nötig;
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` steht bereits vor
dieser Welle bei 3× (Ausgang nicht zugewiesen, andere Sub-Area: die
Beispiel-Client-Dockerfiles). Die Einträge bei 2×
(`aufschub-adresse-nimmt-sendung-nicht-an`,
`nachzug-laesst-ueberholten-text-stehen`, `test-runner-stiller-ausschluss`,
`vorab-bedingung-nach-umsetzung-geprueft`,
`kommentar-behauptet-nicht-getragenen-fehlerpfad`,
`db-gegenstand-enthaelt-netzlos-geprueften-code`,
`rollen-test-abdeckungsluecken`, `adapter-fehler-ausgang`) erreichen die
Schwelle mit dem nächsten Beleg; der Lese-Schritt der Welle-Closure liest sie.
Treffer je Sub-Area, mit dem Slice, der sie als Risiko trägt (Detail je Slice
in §8):

- `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×) —
  **einschlägig:** die Abhilfe nach einer nicht anwendbaren Regel ist der erste
  Administrationsweg nach einem `schema`-Abbruch, für **eine** Ursache;
  `e2e-abhilfe` belegt ihn, `betriebsdoku` beschreibt ihn; die allgemeine
  Recovery bleibt offen (§6).
- `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×) —
  **einschlägig:** die Abweichungen 2 und 5 in §4 und der Ort der Auswertung
  (`kern-rename` §6) weichen vom ADR-Wortlaut ab; sie sind als Fragen an den
  Auftraggeber geführt.
- `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×) und
  `BEO-PGC/aufschub-adresse-verfaellt` (verkörpert, 3×) — der Handbuch-Aufschub
  von `antragsweg-schema` trägt eine Adresse, deren §2 den Gegenstand deckt;
  „die Closure dieser Welle“ ist eine Adresse, die eintreten kann (§5).
- `BEO-PGC/adr-folgepflicht-ohne-traeger-slice` (offen, 1×) — die
  Träger-Tabelle in §5 gibt jeder Folgepflicht von
  [`ADR-0112`](../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  einen Träger oder eine benannte Adresse.
- `BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (offen, 1×) und
  `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×) —
  `backfill-pfad` (Architect-Kurzverdikt als Vorab-Bedingung im Start-Trigger)
  und `e2e-abhilfe` (Ordnung steht, bevor ihre Wirkung gemessen wird).
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert,
  [`AGENTS.md`](../../../AGENTS.md) §3.13, 26×) — betrifft jeden Slice; der
  Suchlauf steht als committetes Feld je Slice-Plan §3.
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 2×) — `spec-nachzug`
  und `betriebsdoku`.
- `BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×),
  `BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×), `BEO-PGC/rollen-verdrahtung`
  (eingetreten), `BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×) —
  `antragsweg-schema`: neue Spalten, CHECK-Menge, zwei Funktionen,
  Guard-Eintrag, Grants.
- `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (verkörpert) —
  `antragsweg-usecase`: der Regelstand ist ein dauerhafter, tabellen-scoped
  Träger, der `Assembler`-Cache ist Laufzeit.
- `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 2×),
  `BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert),
  `BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
  `antragsweg-usecase`, `kern-rename`: neuer Code verschiebt Nenner und Zähler
  der Coverage-Messungen; die Faltung der `applied`-Zeilen liegt in der netzlos
  gemessenen Fläche.
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
  `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×),
  `BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×),
  `BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×) —
  `e2e-wirkung`, `e2e-abhilfe`, `backfill-pfad`, `antragsweg-usecase`.
- `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×),
  `BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, 9×) — `kern-rename`,
  `backfill-pfad`, `start-reihenfolge`: Kommentare zu Fehlerpfaden sagen nur
  zu, was der Code trägt.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert, 3×), `BEO-PGC/handbuch-versionshistorie-uebersprungen`
  (verkörpert, 3×) — `antragsweg-schema` (Aufschub), `betriebsdoku` (Handbuch
  und Änderungshistorie).
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 13×),
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 6×),
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 8×),
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×) — Slices
  mit Messungen und Negativ-Belegen (`e2e-*`, `kern-rename`,
  `antragsweg-usecase`, `backfill-pfad`).
- `BEO-PGC/roadmap-kontext-verwaist-nach-zeilen-entfernung` (offen, 1×) —
  gesichtet bei der Roadmap-Fortschreibung dieser Eröffnung: der
  Platzhalter-Knoten `TRF` und sein Prosa-Absatz sind durch die Welle ersetzt,
  kein darüberstehender Gruppen-Header bleibt ohne Bezug.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 7×) — **nicht
  einschlägig:** kein Workflow-Zug in dieser Welle.
- Gesichtet, ohne Bezug zu dieser Welle: die übrigen Einträge des Registers.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-transformationen-results.md`
Zähler: `../observations/` (Beobachtungs-Register)
