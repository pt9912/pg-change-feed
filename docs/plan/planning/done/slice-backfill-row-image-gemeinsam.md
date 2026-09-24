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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      Report `docs/reviews/review-slice-backfill-row-image-gemeinsam.md`
      (Erstlauf 1 HIGH, 0 MEDIUM, 1 LOW, 5 INFO); F-1 und F-2 in der Fixrunde
      (`cd046787`, `dc633a97`) behoben und vom Reviewer am Kopf nachgeprüft,
      kein offenes HIGH/MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — kein öffentlicher Vertrag berührt (Byte-Gleichheit); die Suche in §3 belegt es.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke). *(§7: geschärfte Anwendung der
      Zahl-im-Träger-Regel auf Suchlauf-Felder; benannte Lücke.)*
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. *(Eine weitere `evidence/`-Datei
      in `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, ein neues
      Verzeichnis `BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet`;
      siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen). *(Vier Ausgänge in §6, je am Ort; siehe §7.)*
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
| Go-Code und Kommentare | Parent (Commit `43f394a1`, Stand vor diesem Slice): `git grep -n 'rowImage' 43f394a1 -- '*.go'`; Diff-Stand: `git grep -n 'rowImage' HEAD -- '*.go'` (Arbeitsbaum sauber, gleiche Zahl wie `grep -rn 'rowImage' --include=*.go .`) | Parent: **16** Treffer (gemessen) — 5 in `mapper.go` (Aufrufe Z. 233/237, Kommentar Z. 471, Definition samt Doc Z. 515/524), 4 in `natsstream/publisher.go`, 6 in `driving/http/{sse,readchanges}.go` (4 in `sse.go`, 2 in `readchanges.go`), 1 in `examples/nats-stream-client/format_test.go`. Diff-Stand: **15** Treffer (gemessen) — alle 5 `mapper.go`-Treffer weg; die übrigen 11 bleiben (Summe 4+6+1) — sie gehören zu **zwei anderen, gleichnamigen** Funktionen (`natsstream.rowImage`, `http.rowImage`), die ein **bereits gebautes** Bild (`[]byte`) als `json.RawMessage` in die Wire-Nachricht einbetten und kein Bild konstruieren. Neue Treffer: 4, alle aus der Testvariable `rowImageColumns` in `rowimage_test.go` (Z. 18/76/93/107) | mapper-Treffer nachgezogen; die elf anderen nicht ändern (andere Funktion, andere Schicht; die Namensgleichheit ist eine gemeldete Beobachtung, kein Teil dieses Slice) |
| Dokumente (Spec, Handbuch, Harness) | Parent: `git grep -n 'rowImage' 43f394a1 -- '*.md'`; Diff-Stand: `git grep -n 'rowImage' HEAD -- '*.md'` | Parent und Diff tragen dieselben Treffer außerhalb dieser Plan-Datei: `docs/plan/planning/welle-backfill-bestand.md` Z. 187 (K1, zitiert [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) „in `rowImage`"), `docs/plan/planning/open/slice-transformationen-kern-rename.md` Z. 240 (nennt das Entfallen des privaten `rowImage` bereits als kommend), `docs/plan/planning/done/welle-18-results.md` Z. 40 (Record), `docs/reviews/review-slice-nats-drittstream-example-*.md` (Records; meinen `natsstream.rowImage`). Kein Treffer in `spec/`, `docs/user/`, `harness/`, `README.md`, `AGENTS.md`. Ergänzend `grep -rnE 'Row-?Image-Konstruktion' --include=*.md .`: benennt den *Ort* nirgends als Mapper-Datei außer den `Accepted` ADRs | keine Änderung: die beiden Plan-Treffer (fremde offene Pläne) bleiben wahr bzw. zitieren [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) wörtlich — in den Bericht gemeldet, nicht still mitgeändert; Records/`Accepted` unberührbar |
| `Accepted` ADRs, die den Ort als Kontext nennen (`mapper.rowImage`) | `git grep -n 'rowImage' 43f394a1 -- 'docs/plan/adr/*.md'` bzw. mit `HEAD` | Parent und Diff identisch (gemessen an beiden Ständen): `0059` (Z. 59/60/152/240), `0111` (Z. 88/228), `0112` (Z. 89/170/240/359) — 10 Treffer, jeder nennt `rowImage` als den Ort zum Entscheidungszeitpunkt | nicht ändern: der Ort steht dort als Stand zum Entscheidungszeitpunkt; nur eine reine Verweisgerüst-Korrektur wäre nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) zulässig ([`AGENTS.md`](../../../../AGENTS.md) §3.5) |
| Zweite JSON-Erzeugung für Bilder | `git grep -n 'json.Marshal' 43f394a1 -- 'internal/*.go'` bzw. mit `HEAD` (Fundstellen nach Bild-Bezug lesen) | Je Stand 6 Treffer (gemessen). Parent: `mapper.go` (2× in `rowImage`), `natsstream/publisher.go` (`json.Marshal(toStreamMessage(change))`), `http/sse.go` (`json.Marshal(toStreamChange(change))`), dazu zwei Test-Treffer (`natsstream/publisher_test.go` Z. 268, `http/sse_test.go` Z. 422 — Letzterer ein Kommentar). Diff-Stand: `model/rowimage.go` (2×, die eine Funktion), dieselben zwei Wire-Stellen und dieselben zwei Test-Treffer. Gelesen: die beiden Wire-Stellen serialisieren die Nachricht und betten das Bild als `json.RawMessage` ein — keine Bild-Konstruktion; die Test-Treffer konstruieren kein Bild. Ergänzend `git grep -n "WriteByte('{')" <Stand> -- 'internal/*.go'`: je Stand ein Treffer (Parent `mapper.go` Z. 529, Diff-Stand `model/rowimage.go` Z. 31) | Fundstellen, die ein Row Image bauen, gehören in die eine Funktion |

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
  Referenz-Bytes mit diesen Zeichen. **Ausgang:** *entfallen* — nicht
  eingetreten: die Referenz-Bytes mit `<`, `>`, `&` (Fall „Maskierung" in
  `TestBuildRowImageBytes`) sind am Parent-Stand `43f394a1` von Reviewer und
  Verifier je unabhängig in einem temporären Worktree gegen das alte `rowImage`
  reproduziert (Review §Eigenständig durchgeführte Prüfungen; Verifikation §3),
  und die Reviewer-Mutation „Wert ohne `json.Marshal`" (M7) färbt den
  Maskierungs-Fall rot.
- **Ort der Funktion.** Die Domäne setzt voraus, dass die Eingabe ohne
  `decode`-Typen auskommt. *Erwartet, zu belegen durch:* `make a-check` und der
  Import-Blick auf die neue Datei. **Ausgang:** *entfallen* — nicht
  eingetreten: `make a-check` endet Exit 0 mit „gesamt: 0 Befund(e)" (Review
  und Verifikation §1, je eigener Lauf), und `rowimage.go` importiert nur `bytes`
  und `encoding/json` (Verifikation §2, Zeile 1).
- **Kopplung an die Transformations-Umsetzung** (Welle §5, K1): die Funktion
  ist ihr Ansatzpunkt; ein Signatur-Zuschnitt, der jeden Aufrufer bei einer
  Erweiterung ändern müsste, verschöbe Kosten in die zweite Welle. *Erwartet, zu
  belegen durch:* Review liest die Signatur mit dieser Frage. **Ausgang:**
  *eingetreten* in der erwarteten, begrenzten Form → Folge-Slice
  `slice-transformationen-kern-rename`: die Signatur ist positional, ein
  Regel-Parameter ändert die zwei Aufrufstellen in `Assembler.change` (je Bild
  eine) und die künftige Backfill-Aufrufstelle (Review F-7, Verifikation V-2,
  beide INFO — kein Vorbau, keine Abweichung). Die Übergabe steht im Plan des
  Empfängers (§3-Suchlauf, Zeile „Aufrufer der gemeinsamen Funktion", am Stand
  `89053d3b` nachgemessen); kein Carveout.
- **Leistung im Backfill.** Ein Backfill ruft die Funktion je Zeile; sie darf
  nicht mehr allozieren als das alte `rowImage`. *Erwartet, zu belegen durch:*
  ein `go test -bench` gegen den Parent-Stand oder eine begründete
  Nicht-Messung im Bericht. *Wirkung (Ursprung: Code-Lesung):* der
  WAL-Adapter bildet einmal je Änderung eine Namensliste (`columnNames`,
  geteilt von neuem und altem Bild) — das ist **eine zusätzliche Allokation je
  WAL-Änderung**; `BuildRowImage` selbst allokiert nicht mehr als das alte
  `rowImage`. *Review-Messung, nicht committet* (Vergleich Parent/Kopf,
  `go test -bench` je 300000×, 3 Läufe, 6 Spalten, eine ausgeschlossen):
  64 → 65 allocs/op, 1313 → 1409 B/op, ns/op im Rauschen; der Lauf liegt in
  keinem Artefakt des Repos. *Ergebnis:* tragbar, keine Vermeidung — eine
  Vermeidung der Namensliste ist ausdrücklich nicht Teil dieses Slice.
  **Ausgang:** *entfallen* — die Bedingung des Risikos („die Funktion allokiert
  mehr als das alte `rowImage`") ist nicht eingetreten: `BuildRowImage` hat
  denselben Rumpf (Review, zeilenweiser Vergleich; Verifikation §3), und ein
  Backfill-Aufrufer trägt den Zuschlag der Namensliste nicht. Bestehen bleibt
  die Wirkung im WAL-Adapter, als tragbar entschieden und ohne Folge-Slice:
  **eine zusätzliche Allokation je WAL-Änderung** (`columnNames`, Ursprung
  Code-Lesung; der Verifier bestätigte am Code, dass `columnNames` genau ein
  `make([]string, …)` je Änderung anlegt und die Liste über beide Bild-Aufrufe
  geteilt wird). Die Zahlen 64 → 65 allocs/op und 1313 → 1409 B/op sind
  **Review-Messung, nicht committet** (Review F-4); der Verifier hat sie nicht
  nachgemessen, er hat nur die Differenz von 96 B gegen sechs Spalten × 16 B
  Zeichenketten-Kopf rechnerisch plausibilisiert.

## 7. Closure-Notiz

- **Was hat funktioniert:** der Zug blieb eine Verschiebung einer Funktion mit
  Tests, ohne Rückführung (§4). Die Byte-Gleichheit ist von zwei Rollen
  unabhängig am Parent-Stand `43f394a1` belegt: Reviewer und Verifier fuhren
  je in einem temporären Worktree die Referenz-Erwartungen des Kopfes gegen das
  **alte** `rowImage` — alle Fälle grün (Review §Eigenständig durchgeführte
  Prüfungen; Verifikation §3). Die Zusagen sind an ihrer Eingabeseite
  gebunden: der Review setzte acht Mutationen, die Verifikation vier, alle
  rot (Review §Mutationen; Verifikation §5). Der Review (Erstlauf 1 HIGH ·
  0 MEDIUM · 1 LOW · 5 INFO, Summary des Reports) und die Verifikation
  („Bestätigt", sieben `[x]`-Zeilen je mit eigenem Beleg, `make gates`
  EXIT=0 im Verifier-Lauf) liefen im frischen Kontext; die
  Suchlauf-Zahlen maß der Verifier an beiden Ständen nach (vier
  Plan-Zeilen, alle bestätigt).
- **Was ging anders als geplant:** die Fixrunde nach dem Review zog drei
  Träger: F-1 (die Summe des Suchlauf-Feldes und der Parent-Bezug des
  Befehls, `cd046787`), F-2 (der Doc-Kommentar von `BuildRowImage` stellte den
  Backfill-Aufrufer als Bestand dar, `dc633a97`) und F-4 (die Wirkung der
  Namensliste als Zahl mit Ursprung im Plan, `cd046787`). Der Plan wuchs
  im Lauf um die Zeile für `columnNames` und die Ausnahme, dass
  `containsColumn` im Mapper bleibt (§3); beide Abweichungen stehen mit
  Begründung in §3, der Diff ist zu ihnen wahr (Verifikation §3).
- **Steering-Loop-Eintrag (Lerneintrag):** geschärfte Anwendung der
  verkörperten Regel [`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`):
  *das §3.13-Suchlauf-Feld ist selbst ein Zahlen-Träger und trägt dieselbe
  Pflicht — der Parent-Stand steht als Commit-Kennung, nie als `HEAD` oder
  „vorher"; jede Gesamtzahl ist gedruckt gemessen (`… | wc -l`) oder als
  abgeleitet gekennzeichnet, nie aus den Teilzahlen addiert als Messung
  ausgegeben.* Beleg: das Feld nannte 17 gegen gemessen 16 mit stimmenden
  Summanden (Review F-1), und sein Parent-Befehl nannte `HEAD` (Verifikation
  V-4 hält fest, dass ein `HEAD`-Bezug im Diff-Stand-Befehl nur bis zum
  nächsten Commit wahr bleibt). Kein neuer Sensor: der Reviewer fand beide
  Stellen durch Nachmessen, kein Gate liest Prosa-Zahlen (die Begründung
  führt die Klasse im Register). **Benannte Lücken** (kein Träger im Repo, keine
  Spec-Lücke): (a) die Allokations-Messung Parent gegen Kopf ist über kein
  `make`-Ziel wiederholbar — `BenchmarkBuildRowImage` liegt im Domänen-Test,
  der Vergleich gegen den Parent-Stand lief ad hoc in temporären Worktrees;
  deshalb tragen die Zahlen im Plan den Ursprung „Review-Messung, nicht
  committet" (§6). (b) Die Kopplung K1 hängt an einer positionalen Signatur;
  ob der Regel-Parameter als weiteres Positionsargument oder anders eintritt,
  entscheidet `slice-transformationen-kern-rename`, nicht dieser Slice.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** — weitere Datei
    `evidence/slice-backfill-row-image-gemeinsam.md` (Review F-1 und F-4,
    Verifikation V-4); Zähler **14×** (Datei-Anzahl unter `evidence/`, real
    ausgezählt). Ausgang bleibt verkörpert — kein neuer Schwellen-Übertritt,
    kein Lese-Schritt der Welle-Closure.
  - **`BEO-PGC/subagent-write-ablehnung-als-zielpfad-sperre-gemeldet`** —
    neues Verzeichnis (Kürzel `PGC`), Zähler **1×**, offen: ein
    Reviewer-Subagent meldete den Zielpfad des Reports als gesperrt, nachdem
    ein `Write` auf einen Scratchpad-Entwurf abgelehnt worden war; der Zielpfad
    war nicht gesperrt (Beleg: Commit `0d5c3f1a`; Ursache vom Auftraggeber im
    Transkript geprüft, vom Planner nicht nachgemessen).
  - **F-2 (Kommentar stellt eine Zusage der Welle als bestehenden Zustand
    dar)** — benannt, nicht angelegt: erstes Auftreten dieser Form. Gegen
    `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (2×) geprüft und
    nicht dorthin gezählt: jene Klasse betrifft einen zugesagten **Fehlerpfad**,
    der das Verhalten umkehrt; hier war das Verhalten der Funktion richtig, nur
    der Aufrufer stand im Präsens. Ein zweites Auftreten legt die Klasse an.
  - **F-6 (namensgleiche `rowImage` in `natsstream` und `http`)** — kein
    Befund: der Review bewertet die Verwechslungsgefahr als durch den Diff
    verringert; der Plan meldet die Namensgleichheit selbst (§3), keine
    Beobachtung angefallen.
  - **F-3, F-5, F-7, V-1 bis V-3** — INFO ohne erwartete Aktion (zulässige
    Test-Provenienz, dritte Kopie einer Sechs-Zeilen-Namenssuche,
    Kopplungs-Beurteilung, Coverage-Rauschen 82.60 % gegen 82.70 % in zwei
    gedruckten Läufen desselben Baums); keine Beobachtung angefallen.
- **Träger-Übergaben (§3.13):** die zwei Handbuch-Meldungen des Vorgängers
  standen nur als Text im Plan von `slice-backfill-spec-nachzug`; sie sind
  jetzt in den Plänen der Empfänger eingetragen (am Stand `89053d3b`
  nachgemessen): `slice-backfill-change-origin` trägt die zwei
  `origin`-Stellen (SQL-Beispiel unter „Änderungen lesen", Feldliste unter
  „Changes lesen"), `slice-backfill-sql-administration` die zwei
  Fortsetzungs-Idiome samt der Nennung von `backfill` (0 Treffer im Handbuch);
  beide mit `Version:`-Kopf und Änderungshistorie. Die Übergabe der Kopplung
  K1 trägt der Plan von `slice-transformationen-kern-rename`.
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice ist eine interne
  Refaktorierung ohne End-Nutzer-Wert; das Ergebnis ist byte-gleich zum
  Parent, ein Nutzer-Bedarf ist daran nicht zu validieren. Kein stilles
  Überspringen.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):**
  eine getrennte Rolle im frischen Kontext, kein Schritt der Planner-Closure;
  der Skill prüft Slices in `done/` und greift daher erst nach dem `git mv` —
  hier nicht ausgeführt.
- **Folge-Slices:** keine neuen — die Folge-Slices der Welle
  [welle-backfill-bestand](../welle-backfill-bestand.md) liegen als Dateien in `open/`; die Kopplung K1 übernimmt
  `slice-transformationen-kern-rename` (Start-Trigger jetzt erfüllt), die
  Handbuch-Übergaben gehen an `slice-backfill-change-origin` und
  `slice-backfill-sql-administration`.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Byte-Abweichung durch
  Escaping) **entfallen**; Risiko 2 (Ort der Funktion) **entfallen**;
  Risiko 3 (Kopplung K1) **eingetreten** in der begrenzten Form, übergeben an
  `slice-transformationen-kern-rename`; Risiko 4 (Leistung) **entfallen**
  für die Funktion, die Wirkung im WAL-Adapter (eine zusätzliche Allokation
  je WAL-Änderung, Review-Messung, nicht committet) ist als tragbar
  entschieden.
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure. (a) Anker: der Lerneintrag
  verkörpert nichts neu, er schärft die Anwendung einer bestehenden Regel
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A, am Ort existent); (b) Folge-Slice: keiner neu; (c) Register:
  beide genannten Kennungen existieren als Verzeichnis, und ihre
  `evidence/`-Verzeichnisse tragen die neue Datei (14× bzw. 1×, real
  ausgezählt).

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
