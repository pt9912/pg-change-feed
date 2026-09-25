# Architect-Verdikt: Retention — neuer Consumer im laufenden Lauf und Transaktion über der Seitengrenze

**Rolle:** Architect (Modul 8)

**Anlass:** Zwei offene Bewertungen aus dem Review von
`slice-retention-lauf-speicher-begrenzung`, die der Plan dem Architect zuweist:
F-3 (ein Consumer bestätigt erstmals oder wird zurückgesetzt, während ein Lauf
läuft; Plan §6, sechstes Risiko) und F-13 a (eine Transaktion liegt über zwei
Seiten). Der Rollenwechsel Implementer → Architect → Planner folgt Modul 8
§Konflikt-Pfad.

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der
Implementer-Lauf des Slice)

**Datum:** 2026-09-25

**Bezug:** [`LH-FA-RET-003`](../../spec/lastenheft.md),
[`LH-FA-RET-004`](../../spec/lastenheft.md),
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
(Festlegung 4 und 6), [`ADR-0014`](../plan/adr/0014-retention-domain-policy.md),
[`ADR-0029`](../plan/adr/0029-domain-invarianten.md) (Regel 2 und 5),
`architect-verdict-retention-lauf-speicher-begrenzung` und
`review-slice-retention-lauf-speicher-begrenzung` (F-3, F-13 a; beide im
Verzeichnis `docs/reviews/`); der Slice ist als Kennung genannt, nicht als
Pfad-Link (ein Slice wechselt die Lifecycle-Ablage).

**Erzeugte Artefakte dieses Zugs:** dieses Dokument. Keine neue ADR; Plan, Spec,
Handbuch und Code bleiben unberührt.

---

## Verdikt

**F-3 — akzeptiertes Negativ, keine Vertragsänderung, keine neue ADR.** Die
Bewertung trägt der Code, nicht die Laufdauer: ein Consumer ohne Positions-Zeile
ist für die Bereinigung in **jedem** Lauf unsichtbar, nicht nur in einem
laufenden. Das Fenster „Position gelesen — Lauf löscht" ist der letzte Ausschnitt
einer Schutzlücke, die mit der Registrierung beginnt und erst mit der ersten
Bestätigung endet; sie hat keine Obergrenze, und der Slice verlängert nur ihren
Rest um die Laufdauer.

**F-13 a — akzeptiertes Negativ, keine Vertragsänderung, keine neue ADR.**
Alter und Commit-Position sind Eigenschaften der Transaktion; die Freigabe
fällt für alle Changes einer Transaktion gleich aus. Ein Consumer, für den die
Transaktion zählt, bekommt sie nie teilweise gelöscht.

**Empfehlung an den Implementer/Planner** (ohne Code-Änderung): die beiden
Risiko-Ausgänge im Plan auf „akzeptiertes Negativ, Verdikt
`architect-verdict-retention-neue-consumer-und-seitengrenze`" setzen; im Handbuch
den Satz „gilt seine Position erst im nächsten Durchlauf" um die Folge ergänzen
(Wortlaut unten, §Folge-Arbeit).

---

## F-3 — Consumer erstmals bestätigt oder zurückgesetzt im laufenden Lauf

**Was der Code tut** (Beleg: `internal/application/usecase/retention/service.go`
am HEAD, `Run`; `git show 8cd39719:internal/application/usecase/retention/service.go`
für den Stand davor): Positionen und Uhr werden je einmal zu Beginn gelesen, die
Freigabe je Kandidat ist `AllowsDeletion(age, position, consumerPositions)`
(`internal/domain/model/retention.go`). Vor dem Slice lag dieselbe Reihenfolge
vor (Positionen, alle Changes lesen, alles löschen); die Sicherheitsaussage
„Positionen wandern nur vorwärts" trägt für **bekannte** Consumer
([`ADR-0029`](../plan/adr/0029-domain-invarianten.md) Regel 2; `Advance` in
`internal/domain/model/consumer.go` lehnt eine frühere Position mit
`ErrPositionRegression` ab).

**Neuer Consumer.** `ConsumerStatePort.Positions` liest `cdc.consumer_position`
(`queries.SelectConsumerPositionsBySource`); dessen Kommentar sagt, ein Consumer
ohne Zeile „blockiert ihre Retention nicht". Ein frisch registrierter Consumer
trägt keine Zeile (Handbuch, „Startposition eines neuen Consumers": `offset` 0,
`acknowledged` `false`, gemessen im E2E-Lauf laut Handbuch). Er ist damit in jedem
Lauf vor seiner ersten Bestätigung ungeschützt, unabhängig davon, wann der Lauf
startet; das Handbuch nennt es im „Betriebs-Hinweis (Consumer-Bindung)". Bestätigt
er während eines Laufs erstmals, behandelt ihn dieser Lauf wie einen, der eine
Sekunde vor der Bestätigung gestartet wäre: das ist kein neuer Fall, sondern
derselbe, mit Zeitpunkt statt Wartezeit. Der Schutz beginnt mit dem nächsten Lauf
und gilt dann für Changes hinter seiner bestätigten Position.

**Zurücksetzen.** Einen Weg, die Position eines Consumers zu senken, gibt es im
Vertrag nicht: `acknowledge` scheitert an `ErrPositionRegression`; `POST
/consumers/remove` (Handbuch, Tabelle der Operationen) entfernt den Consumer samt
Position, ein erneutes Registrieren ist ein neuer Consumer (dieselbe Lage wie
oben). Ein Rücksetzen per direktem SQL durch `cdc_admin` liegt außerhalb der
Zusage von [`ADR-0029`](../plan/adr/0029-domain-invarianten.md) („regulär nur
vorwärts"); das Fenster dafür ist dieselbe Laufdauer und dieselbe Klasse.

**Mindestalter begrenzt das Risiko.** `AllowsDeletion` gibt nur frei, wenn das
Alter das Mindestalter erreicht (`retentionMinAge` = 24 h,
`internal/bootstrap/wiring.go`; [`LH-FA-RET-003`](../../spec/lastenheft.md)). Die
Uhr wird einmal zu Beginn gelesen: sie liegt früher als jede spätere, das Alter
wird höchstens unterschätzt, eine Löschung kommt nie früher. Betroffen sind damit
nur Changes, die zum Lauf-Beginn seit mindestens 24 h bestehen **und** die alle
bekannten Consumer bereits bestätigt haben.

**Größe des Fensters (abgeleitet, nicht am Feed gemessen).** Die
Bereinigungsdauer am Feed-Container ist ungemessen (Messbericht Abschnitt 6,
Befund 4). Aus den Einzelläufen des Architect-Verdikts (`n` = 1, PostgreSQL 18.6,
1.000.000 Changes; Zeilen in
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
§Gemessen) folgt: Lesen 0,38 bis 2,0 s; Löschen einer vollen Seite von 10.000
21,5 bis 242,2 ms, bei 100 vollen Seiten 2,2 bis 24,2 s; ein Lauf mit
**vollständiger** Löschung höchstens etwa 26 s (Summe, abgeleitet), im Regelbetrieb
mit wenig neu freigegebenen Changes 0,38 bis 2,0 s. Der Stand vor dem Slice lag
bei etwa 3,4 s (1,9 bis 2,0 s Lesen plus 1,5 bis 1,6 s Löschen, gleiche Quelle,
Summe abgeleitet); das Fenster wächst damit im Extremlauf um das Siebenfache, im
Regelbetrieb nicht. Der Extremlauf ist der erste Lauf nach einem großen Bestand, der
das Mindestalter erreicht. Eine eigene Messung am Feed ist nicht nötig: die
Bewertung hängt nicht an der Zahl, sondern am Ursprung der Lücke (siehe Verdikt);
sie trägt auch bei dem Zehnfachen der abgeleiteten Obergrenze.

**Alternative „Positionen je Seite neu lesen"** (Schärfung von
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
Festlegung 4, Folge-ADR mit `Supersedes`): eine kleine Abfrage je Seite, 100
zusätzliche je Lauf bei 1.000.000 Changes (hergeleitet). Sie verkürzt das Fenster
von der Laufdauer auf die Dauer einer Seite (höchstens etwa 0,3 s, abgeleitet), sie
schließt es nicht: zwischen Lesen der Positionen und Löschen bleibt ein Rest wie
schon im Stand davor, ohne Sperre über beide Schritte. Der Gewinn betrifft eine
Lücke, die vor der ersten Bestätigung unbegrenzt offen ist; die Kosten sind eine
weitere ADR, ein geänderter Test (`TestRunReadsPositionsOncePerRun`) und ein
weiterer Review-Lauf vor dem Release. Nicht gewählt.

---

## F-13 a — Transaktion über zwei Seiten

**Eigenschaften der Transaktion.** Die Kandidatenabfrage liest Position und Zeit
aus der Zeile der Transaktion (`queries.SelectRetentionCandidates`:
`t.commit_position`, `t.committed_at` aus `cdc.transaction`, verbunden über
`transaction_id`). Alle Changes einer Transaktion tragen damit dieselbe Position
und dieselbe Zeit; bei gleichen Eingaben (Positionen und Uhr, je einmal je Lauf
gelesen) fällt `AllowsDeletion` für alle gleich aus. Bei einem Backfill trägt der
ganze Run eine Position und einen Zeitpunkt (`blockBuilder.at` wird je Run
gesetzt, `internal/application/usecase/backfill/service.go`; der Schreiber legt
alle Blöcke in einer Datenbank-Transaktion an, `backfillwriter.go`), die
Blöcke `0bf-<run>-<n>` werden ebenfalls gemeinsam freigegeben oder gemeinsam
zurückgehalten.

**Teilweise gelöscht — wann?** Nur wenn ein Lauf zwischen zwei Seiten endet
(Fehler, Beendigung) oder sich die Eingaben des nächsten Laufs geändert haben. Die
Löschmenge eines Laufs entspricht der des Stands vor dem Slice; sie liegt in
Reihenfolge des Schlüssels `change_id` (Text `<Transaktion>-<Sequenz>`, also nicht
die Lese-Ordnung), der Zwischenstand ist ein Ausschnitt dieser Menge.

- **Bekannte Consumer** (Positions-Zeile vorhanden): entweder hat jeder die
  Position der Transaktion erreicht — dann gehört ihm keine Löschung mehr —, oder
  einer liegt davor, und dann wird **keiner** ihrer Changes freigegeben (`Before`
  in `AllowsDeletion`). Eine Transaktion wird für diese Consumer nie zur Hälfte
  entzogen. Der Rest wird im nächsten Lauf gelöscht: Positionen wandern nur
  vorwärts, das Alter wächst.
- **Lese-Lücke an der Transaktionsmitte:** eine bestätigte Position ist eine
  Commit-Position (`SourcePosition`: Quelle und `Offset`), sie trägt keine
  Change-Sequenz. Ein Consumer beginnt an einer Transaktionsgrenze (`from` und
  `to` sind Commit-Positionen, `SelectChanges`); eine Bestätigung deckt alle
  Changes bis einschließlich dieser Position (Handbuch: ein `LIMIT` kann innerhalb
  einer Position nicht fortsetzen). Wer über den Schlüsselvergleich im SQL-Zugriff
  fortsetzt, hat die Position der Transaktion noch nicht bestätigt; ist er
  bekannt, hält er sie damit zurück.
- **Nicht bekannte Consumer:** ohne Positions-Zeile ungeschützt (F-3); sie sehen
  während eines Laufs eine alte Transaktion mit fehlender erster Hälfte statt gar
  keiner. Das ist dieselbe Klasse wie in F-3 — der Stream eines ungeschützten
  Consumers hat an seinem Anfang schon ganze Transaktionen verloren; eine
  Teil-Transaktion fügt keine neue Zusage hinzu, die verletzt würde.

**Zwei Prädikate, die den Schnitt tragen** (gelesen, nicht am Lauf gemessen): die
Wahrheit „Alter und Position sind Transaktions-Eigenschaften" trägt das Schema
(`cdc.transaction` mit `commit_position`, `committed_at`, `UNIQUE (transaction_id,
sequence)` an `cdc.change`, `tools/schema/schema.yaml`) und die Abfrage. Die
Fitness-Function im Store-Tier von
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
prüft „Position und Zeit je Kandidat gleich dem Datensatz von `ReadChanges`";
eine Aussage über eine Transaktion über zwei Seiten ist keine dort belegte
Aussage, sondern hergeleitet aus der Abfrage.

**Nicht Gegenstand:** F-13 b (Zeile mit unbekanntem `origin` ist jetzt Kandidat
statt Lauf-Abbruch) — kein Sicherheitsbefund, der Plan führt ihn als INFO.

---

## Folge-Arbeit

Der Planner zieht nach (`AGENTS.md` §3.13); dieses Dokument ändert keinen Plan.

| Träger | Nachzug |
|---|---|
| Plan `slice-retention-lauf-speicher-begrenzung` §6 | Risiken „Consumer bestätigt erstmals während eines Laufs" und „Ein Abbruch zwischen zwei Seiten": Ausgang „akzeptiertes Negativ, Verdikt `architect-verdict-retention-neue-consumer-und-seitengrenze`"; die Begründung „nicht breiter geworden" durch die abgeleiteten Grenzen oben ersetzen (Regelbetrieb 0,38 bis 2,0 s; Extremlauf etwa 26 s statt etwa 3,4 s) |
| Handbuch, Abschnitt „Aufbewahrung (Retention)", „Betriebs-Hinweis (Consumer-Bindung)" | Satz „Bestätigt ein Consumer erstmals, während ein Durchlauf läuft, gilt seine Position erst im nächsten Durchlauf" um die Folge ergänzen: *Changes, die dieser Durchlauf nach dem Lesen der Positionen löscht, sind für diesen Consumer verloren; sein Schutz beginnt mit dem nächsten Durchlauf.* Und ein Satz zur Nicht-Atomarität: *Ein Durchlauf kann eine Transaktion in mehreren Schritten löschen; ein Consumer mit bestätigter Position sieht davon nichts, ein Consumer ohne Bestätigung kann eine alte Transaktion währenddessen unvollständig lesen.* |
| Test | keine neue Test-Pflicht: die Aussagen sind am Code hergeleitet; der Slice trägt bereits `TestRunReadsPositionsOncePerRun` und die Abbruch-Tests |

Kein neuer Träger für eine Messung: die Laufdauer am Feed bleibt ungemessen
(akzeptiert, weil die Bewertung nicht von ihr abhängt); trifft
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
Trigger (a) ein (Lauf erreicht den Takt), entsteht dort die Messung samt der Frage,
ob Festlegung 4 dann neu bemessen wird.

---

## Auftraggeber-Frage

**Keine offen.** Beide Ausgänge sind akzeptierte Negative mit Begründung am Code;
eine ADR entsteht nicht.
