# Slice backfill-row-image-gemeinsam: Row-Image-Konstruktion als eine gemeinsame, reine Funktion der Domäne — der WAL-Pfad ruft sie, das Ergebnis bleibt byte-gleich

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-008`](../../../../spec/lastenheft.md) (Row-Image-Abwesenheit),
[`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Spaltenausschluss, Wert nirgends im Image),
[`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Ausschluss), [`LH-FA-DAT-005`](../../../../spec/lastenheft.md) (Abwesenheits-Vertrag),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2 (Row-Image-Parität: **eine** Funktion für WAL- und
Backfill-Pfad), [`ADR-0016`](../../adr/0016-jsonb-row-images-mvp.md) (JSONB-Row-Images), [`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 3
(Ausschluss im Bild).

**Berührte Spec-Stellen:** [`SPEC-002`](../../../../spec/pflichtenheft.md) (Row Images als JSON-Objekt) —
gelesen, nicht geändert; [`ARC-001`](../../../../spec/architecture.md) (Domain Core: reine Funktionen) und
[`ARC-006`](../../../../spec/architecture.md) (Driven darf Driving nicht importieren).

**Verantwortlich:** Implementer-Agent, 2026-09-24.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Row-Image-Konstruktion — ein JSON-Objekt über die gesendeten
Spaltenwerte in Relation-Reihenfolge, `NULL`/unverändertes TOAST/ausgeschlossene
Spalte entfallen — liegt als **eine reine Funktion** in der Domäne
(`internal/domain/model`); der WAL-Pfad ruft sie anstelle des privaten
`rowImage` des Driving-Adapter-Mappers, das Ergebnis ist **byte-gleich**.

**Warum in der Domäne.** `rowImage` steht heute im Driving-Adapter-Paket
`internal/adapters/driving/replication/mapper` und nimmt einen
`decode.Relation`. Der Backfill-Pfad (Use Case und Driven Adapter) darf dieses
Paket nicht importieren: die Schichtenregel der Architektur-Sicht
(`spec/architecture.md` §2) lässt [`ARC-006`](../../../../spec/architecture.md) keine Driving Adapters und
[`ARC-002`](../../../../spec/architecture.md) keine konkreten Adapter importieren. Die Domäne ([`ARC-001`](../../../../spec/architecture.md)) ist die
Stelle, die beide Seiten importieren dürfen (`.a-check.yml`: Kanten `adapters →
domain`, `app → domain`); [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) fordert wörtlich „eine Stelle, die beide
Adapter-Schichten importieren dürfen" — die Domäne erfüllt das. Die Eingabe ist
neutral (Spaltennamen, Werte als `[]*string`, Ausschluss-Liste), kein
Adapter-Typ. Signatur und Typ der Spalten-Eingabe legt der Implementer fest.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`origin`** — eigener Slice (`change-origin`); dieser Slice ändert kein
  Feld, kein Schema, keinen Lesepfad.
- **Regelauswertung/Transformationen** ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)) — die Funktion ist der
  Ort, an dem die Transformations-Umsetzung ansetzt (Welle §5, K1); dieser Slice
  baut keine Regel-Schnittstelle vor (kein Umfang ohne Anforderung). Er hält die
  Funktion klein und rein, damit ein weiterer Parameter die Zahl der
  Aufrufstellen nicht vergrößert.
- **Änderung der Bildform** — Schlüssel-Reihenfolge, Escaping und
  Abwesenheits-Vertrag bleiben; ein anderes Bild wäre ein Vertragsbruch für jeden
  Consumer ([`LH-FA-CAP-008`](../../../../spec/lastenheft.md)).
- **Der Snapshot-Leser und `col::text`** — `snapshot-reader` (die Funktion
  bekommt dort ihren zweiten Aufrufer).

## 2. Definition of Done

- [x] Die Row-Image-Konstruktion liegt an genau einer Stelle: die
      JSON-Erzeugung für Row Images kommt ausschließlich aus der
      Domänen-Funktion; der Mapper des Replication-Pfads ruft sie, das
      private `rowImage` entfällt. *Zu belegen durch:* Suchlauf über
      `internal/**` nach der Bild-Erzeugung (Befehl und Fundstellen im
      Bericht) — eine Konstruktionsstelle; `make a-check` grün (die Domäne
      importiert nichts aus anderen Schichten).
- [x] Byte-Gleichheit gegen den Parent-Stand: die bestehenden Mapper-Tests
      (`internal/adapters/driving/replication/mapper/mapper_test.go`) laufen
      ohne geänderte Erwartungswerte grün, und ein neuer Domänen-Test trägt
      Referenz-Bytes, die **am Parent-Stand mit dem alten `rowImage` gemessen**
      sind (Fälle: `nil`-Werteliste → kein Bild; `NULL`-Wert; ausgeschlossene
      Spalte; leeres Bild `{}`; Zeichen, die `encoding/json` maskiert, etwa `<`,
      `>`, `&`, Anführungszeichen, Umlaute, Steuerzeichen; weniger Werte als
      Spalten). *Zu belegen durch:* `make test` mit
      Race-Detector; `git diff` der Mapper-Testdatei zeigt keine geänderte
      Erwartung.
- [x] Der Kommentar-Träger folgt: der Doc-Kommentar des Bildes steht an der
      neuen Funktion und nennt Zusage, Abwesenheits-Vertrag und Ausschluss
      ([`AGENTS.md`](../../../../AGENTS.md) §3.7); Kommentare im Mapper, die `rowImage` nennen, sind
      nachgezogen.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — kein öffentlicher Vertrag berührt (Byte-Gleichheit); die Suche in §3 belegt es.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/rowimage.go` | neu | die eine Funktion `BuildRowImage(columns []string, values []*string, excluded []string) ([]byte, error)`; rein, keine geteilten Zustände; Eingabe neutral (Spaltennamen als `[]string`, kein Adapter-Typ). Ein privater Helfer `containsName` gehört dazu. |
| `internal/domain/model/rowimage_test.go` | neu | Referenz-Bytes vom Parent-Stand (Fälle laut DoD), Abwesenheit der nil-Werteliste, Ausschluss (Wert nirgends), Reinheit/Aliasing/Nebenläufigkeit (`-race`), `BenchmarkBuildRowImage`. |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `Assembler.change` ruft die Domänen-Funktion; `rowImage` und der `encoding/json`-Import entfallen; neuer kleiner Helfer `columnNames` (Relation → Namensliste, einmal je Änderung für beide Bilder). `containsColumn` **bleibt**: `Assembler.ExcludeColumn` (Doppelt-Ausschluss-Prüfung) braucht es, nicht `appendExcluded`/`removeExcluded`. Der Kommentar an `appendExcluded` nennt `model.BuildRowImage` statt `rowImage`. |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | unverändert erwartet | trägt die Byte-Gleichheit; eine geänderte Erwartung wäre ein Plan-Nachzug. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „wo die Row-Image-Konstruktion liegt und wie sie heißt"; beide Stände gemessen: Parent und Diff):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Go-Code und Kommentare | `grep -rn 'rowImage' --include=*.go .` (Parent: `git grep -n 'rowImage' HEAD -- '*.go'`) | Parent: 17 Treffer — 5 in `mapper.go` (Aufrufe Z. 233/237, Kommentar Z. 471, Definition samt Doc Z. 515/524), 4 in `natsstream/publisher.go`, 6 in `driving/http/{sse,readchanges}.go`, 1 in `examples/nats-stream-client/format_test.go`. Diff: alle 5 `mapper.go`-Treffer weg; die übrigen 11 bleiben — sie gehören zu **zwei anderen, gleichnamigen** Funktionen (`natsstream.rowImage`, `http.rowImage`), die ein **bereits gebautes** Bild (`[]byte`) als `json.RawMessage` in die Wire-Nachricht einbetten und kein Bild konstruieren. Neue Treffer: nur `rowImageColumns` in `rowimage_test.go` (Testvariable) | mapper-Treffer nachgezogen; die elf anderen nicht ändern (andere Funktion, andere Schicht; die Namensgleichheit ist eine gemeldete Beobachtung, kein Teil dieses Slice) |
| Dokumente (Spec, Handbuch, Harness) | `grep -rn 'rowImage' --include=*.md .` (Parent: `git grep -n 'rowImage' HEAD -- '*.md'`) | Parent und Diff tragen dieselben Treffer außerhalb dieser Plan-Datei: `docs/plan/planning/welle-backfill-bestand.md` Z. 187 (K1, zitiert [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) „in `rowImage`"), `docs/plan/planning/open/slice-transformationen-kern-rename.md` Z. 240 (nennt das Entfallen des privaten `rowImage` bereits als kommend), `docs/plan/planning/done/welle-18-results.md` Z. 40 (Record), `docs/reviews/review-slice-nats-drittstream-example-*.md` (Records; meinen `natsstream.rowImage`). Kein Treffer in `spec/`, `docs/user/`, `harness/`, `README.md`, `AGENTS.md`. Ergänzend `grep -rnE 'Row-?Image-Konstruktion' --include=*.md .`: benennt den *Ort* nirgends als Mapper-Datei außer den `Accepted` ADRs | keine Änderung: die beiden Plan-Treffer (fremde offene Pläne) bleiben wahr bzw. zitieren [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) wörtlich — in den Bericht gemeldet, nicht still mitgeändert; Records/`Accepted` unberührbar |
| `Accepted` ADRs, die den Ort als Kontext nennen (`mapper.rowImage`) | dieselbe Suche, Treffer unter `docs/plan/adr/` | Parent und Diff identisch: `0059` (Z. 59/60/152/240), `0111` (Z. 88/228), `0112` (Z. 89/170/240/359) — 10 Treffer, jeder nennt `rowImage` als den Ort zum Entscheidungszeitpunkt | nicht ändern: der Ort steht dort als Stand zum Entscheidungszeitpunkt; nur eine reine Verweisgerüst-Korrektur wäre nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) zulässig ([`AGENTS.md`](../../../../AGENTS.md) §3.5) |
| Zweite JSON-Erzeugung für Bilder | `grep -rn 'json.Marshal' internal --include=*.go` (Fundstellen nach Bild-Bezug lesen) | Parent: `mapper.go` (2× in `rowImage`), `natsstream/publisher.go` (`json.Marshal(toStreamMessage(change))`), `http/sse.go` (`json.Marshal(toStreamChange(change))`). Diff: `model/rowimage.go` (2×, die eine Funktion), dieselben zwei Wire-Stellen. Gelesen: die beiden Wire-Stellen serialisieren die Nachricht und betten das Bild als `json.RawMessage` ein — keine Bild-Konstruktion. Ergänzend `grep -rn "WriteByte('{')" internal --include=*.go`: ein Treffer (`model/rowimage.go`) | Fundstellen, die ein Row Image bauen, gehören in die eine Funktion |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `spec-nachzug` in `done/` liegt
und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet —
  Verschiebung einer Funktion mit Tests; wächst der Zug um eine
  Regel-Schnittstelle, ist das der Aufschub-Fall des ersten Ausschlusses.
- `in-progress` → `open` (blockiert): falls `make a-check` die Domäne als Ort
  ablehnt oder die Byte-Gleichheit nicht herzustellen ist (ein
  Referenz-Byte-Unterschied ohne Verhaltensänderung an der Quelle) — dann
  Architect-Frage zum Ort statt einer Anpassung der Referenz.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Byte-Abweichung durch Escaping.** `encoding/json` maskiert `<`, `>`, `&`
  standardmäßig; ein anderer Encoder oder ein `SetEscapeHTML(false)` änderte
  jedes Bild mit solchen Zeichen. *Erwartet, zu belegen durch:* die
  Referenz-Bytes mit diesen Zeichen. **Ausgang:** *(bei Closure)*
- **Ort der Funktion.** Die Domäne setzt voraus, dass die Eingabe ohne
  `decode`-Typen auskommt. *Erwartet, zu belegen durch:* `make a-check` und der
  Import-Blick auf die neue Datei. **Ausgang:** *(bei Closure)*
- **Kopplung an die Transformations-Umsetzung** (Welle §5, K1): die Funktion
  ist ihr Ansatzpunkt; ein Signatur-Zuschnitt, der jeden Aufrufer bei einer
  Erweiterung ändern müsste, verschöbe Kosten in die zweite Welle. *Erwartet, zu
  belegen durch:* Review liest die Signatur mit dieser Frage. **Ausgang:**
  *(bei Closure)*
- **Leistung im Backfill.** Ein Backfill ruft die Funktion je Zeile; sie darf
  nicht mehr allozieren als das alte `rowImage`. *Erwartet, zu belegen durch:*
  ein `go test -bench` gegen den Parent-Stand oder eine begründete
  Nicht-Messung im Bericht. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Domäne und Replication-Mapper sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, 9×, der Kommentar der
neuen Funktion trägt keine Slice-/Wellen-Chronik),
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×, der
Kommentar sagt nur zu, was die Funktion trägt),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, einschlägig:
der Ort der Funktion ist eine Auslegung des ADR-Wortlauts „eine Stelle, die
beide Adapter-Schichten importieren dürfen" — die Wahl der Domäne ist in §1
begründet und im Review prüfbar).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
