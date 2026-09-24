# Verifikations-Report: slice-backfill-row-image-gemeinsam — 2026-09-24

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-backfill-row-image-gemeinsam.md`](review-slice-backfill-row-image-gemeinsam.md);
Hausform dieses Reports:
[`verifikation-slice-backfill-spec-nachzug.md`](verifikation-slice-backfill-spec-nachzug.md).

**Gegenstand:** `Slice-Plan` `slice-backfill-row-image-gemeinsam` (Welle
`welle-backfill-bestand`), Diff-Range `43f394a1..HEAD` (`0d5c3f1a`):
Implementer-Commits `14653d1d` … `f0d18d91` (zwei Lifecycle-Moves, ein
Meta-Commit, ein Code-Commit `ed76b32e`, Plan-Commits), Fixrunde
`dc633a97`/`cd046787`, Review-Report `0d5c3f1a`. 5 Dateien, 630 Insertions /
56 Deletions (`git diff --stat 43f394a1..HEAD`): der Slice-Plan, der
Review-Report, `mapper.go`, `rowimage.go`, `rowimage_test.go`. Dieser Lauf
ändert keine Spec, keinen Plan und keinen Code; er schreibt nur diesen Report.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf
gesondert gesichert (`make … > <log> 2>&1; echo …=$?`), die Logs danach gelesen.

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make gates` (vor dem Report-Commit) | **EXIT=0** | baseline-verify v6.9.0 OK (54 Dateien) · `docs-check` `971 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability` OK (5 Commits in `HEAD~5..HEAD`) · coverage-gate `coverage-gate: OK — Coverage 82.60% erfüllt Schwelle 80%` (`total: (statements) 82.6%`) · generated-sync OK (byte-gleich, Stufe `proto-export`) · a-check `gesamt: 0 Befund(e)` |
| `make test` (Race-Detector, Kopf) | **EXIT=0** | `ok` für `replication/mapper` und `domain/model`, kein `FAIL` im Log |
| `make a-check` (einzeln) | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` (einzeln) | **EXIT=0** | `total: (statements) 82.7%`, `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%`; `rowimage.go:26 BuildRowImage 90.9%`, `containsName 100.0%` |
| `make commit-traceability RANGE=43f394a1..HEAD` | **EXIT=0** | „OK — 10 Commit(s) in `43f394a1..HEAD`, Betreffs ohne Struktur-ID"; `git log --format=%s 43f394a1..HEAD \| grep -E 'SPEC-\|ARC-'` ohne Treffer |
| `make doc-commits RANGE=43f394a1..HEAD` | **EXIT=0** | 971 Dateien, 0 Befunde |
| `make doc-immutable RANGE=43f394a1..HEAD` | **EXIT=0** | 971 Dateien, 0 Befunde; kein `Accepted`-ADR im Diff (`git diff --name-only 43f394a1..HEAD \| grep -c 'docs/plan/adr/'` = 0) |
| `make doc-trace` (advisory) | **EXIT=0** | 80 Anforderung(en), 3 Waise(n): `LH-FA-CAP-009`, `LH-FA-CFG-007`, `LH-FA-CFG-008` — unverändert gegenüber dem Vorgänger-Slice, nur zur Kenntnis |

Die Coverage-Zahl schwankt zwischen zwei Läufen desselben Baums (82.60 % im
`make gates`-Lauf, 82.70 % im Einzellauf); beide sind gedruckte Zeilen konkreter
Läufe, beide über der Schwelle (V-3).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Row-Image-Konstruktion an genau einer Stelle; Mapper ruft sie, privates `rowImage` entfällt; `make a-check` grün | **erfüllt** | Konstruktionsstelle: `git grep -n "WriteByte('{')" HEAD -- 'internal/*.go'` → **ein** Treffer, `internal/domain/model/rowimage.go:31` (Parent: ein Treffer, `mapper.go:529`). `git grep -n 'json.Marshal' HEAD -- 'internal/*.go'` → 6 Treffer: `rowimage.go:37/41` (die eine Funktion), `natsstream/publisher.go:228` und `http/sse.go:114` (serialisieren die Nachricht, das Bild geht als `json.RawMessage` hinein), zwei Test-Treffer (`publisher_test.go:268`, `sse_test.go:422`, Letzterer ein Kommentar). Die gleichnamigen `rowImage` in `publisher.go:156` und `sse.go:48` sind 5-Zeilen-Funktionen (`len(image)==0` → `null`, sonst `json.RawMessage(image)`) — sie konstruieren nichts. `git grep -n BuildRowImage HEAD -- internal` (ohne Tests): genau die zwei Aufrufe `mapper.go:233/237`. `rowimage.go` importiert nur `bytes` und `encoding/json`; `make a-check` EXIT=0, `.a-check.yml` führt `domain` als eigenen Bereich |
| 2 | Byte-Gleichheit gegen den Parent; `mapper_test.go` ohne geänderte Erwartung; neuer Domänen-Test mit am Parent gemessenen Referenz-Bytes; `make test` | **erfüllt** | `git diff 43f394a1..HEAD -- internal/adapters/driving/replication/mapper/mapper_test.go \| wc -l` = 0. Referenz-Bytes selbst am Parent reproduziert (§3): `rowimage_test.go` des Kopfes, per Shim gegen das **alte** `rowImage` gefahren, grün. Fälle des DoD im Test vertreten: nil-Werteliste (`TestBuildRowImageNilValuesIsAbsent`), NULL, ausgeschlossene Spalte, `{}` (zwei Fälle), Maskierung (`<`, `>`, `&`, `"`, Umlaut, `\x01`, `\n`, `\t`, `\`, U+2028), ungültiges UTF-8, weniger und mehr Werte als Spalten. `make test` EXIT=0 |
| 3 | Kommentar-Träger folgt (Doc an der neuen Funktion nennt Zusage, Abwesenheits-Vertrag, Ausschluss; Mapper-Kommentare nachgezogen) | **erfüllt** | `rowimage.go:8–25` nennt Text-Stand ohne Typ-Interpretation, Abwesenheit für nil/fehlende Spalte, Ausschluss samt „Wert nie serialisiert", nil-Werteliste = kein Bild; Reinheit als Zusage (`rowimage.go:22–25`), der Aufrufer-Satz („gerufen wird sie vom Replication-Mapper") ist am Kopf wahr (zwei Aufrufe). `git grep -n rowImage HEAD -- internal/adapters/driving/replication/mapper/mapper.go` → 0 Treffer; `appendExcluded`-Kommentar (`mapper.go:471`) und `columnNames` (`mapper.go:515–516`) nennen `model.BuildRowImage`. §3.7-Kandidatenlauf: V-1 |
| 4 | `make gates` grün, Exit-Code ungefiltert gesichert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report unter `docs/reviews/` | **materiell erfüllt** (Aufrufer hakt im Closure) | `review-slice-backfill-row-image-gemeinsam.md` liegt vor; Summary Erstlauf 1 HIGH / 0 MEDIUM / 1 LOW / 5 INFO; F-1 (HIGH, Suchlauf-Summe) und F-2 (LOW, Doc-Kommentar) in der Fixrunde behoben, von mir selbst am Kopf nachgemessen (§3 Suchlauf, §4 Doc); F-3–F-7 INFO ohne Aktion. Kein offenes HIGH/MEDIUM/LOW. Die Rollen-Sequenz (Reviewer → Implementer → Verifier) ist eingehalten |
| 6 | §3.13-Suchlauf: Feld in §3 trägt Gefundenes **und** Nichtgefundenes je Träger, beide Stände gemessen | **erfüllt** | vier Zeilen im Plan §3; alle Zahlen nachgemessen (§3) |
| 7 | Doku-Update entfällt | **erfüllt** | `docs/user/`, `spec/`, `harness/` im Diff unberührt (`git diff --stat`); `git grep -n -iE 'row.?image' HEAD -- spec/architecture.md docs/user harness README.md AGENTS.md` mit Ort-/Konstruktions-Bezug: nur ein Treffer in `harness/README.md` (Zeile `make test-integration`, nennt „Row Image" als Inhalt eines Belegs, keinen Ort) |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 9 | Reconciliation-Register | **erfüllt** (`[x]`, „entfällt") | Greenfield, keine Reconciliation-Datei im Repo |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | vier Risiko-Zeilen tragen „**Ausgang:** *(bei Closure)*" |
| 12 | Drei Paarungen | **korrekt offen** | hängt an der Closure von `welle-backfill-bestand` (offen) |

Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt
wäre (Zeile 5 ist materiell erfüllt und wird beim Closure-Nachzug gehakt).

## 3. Plan-vs-Code-Diff (§3 + §3.13-Suchlauf-Feld)

**Geplant und geliefert:** alle vier Plan-Zeilen (§3) sind im Diff vertreten:
`rowimage.go` neu mit der Signatur
`BuildRowImage(columns []string, values []*string, excluded []string) ([]byte, error)`
und privatem `containsName`; `rowimage_test.go` neu (Referenz-Bytes,
Abwesenheit, Ausschluss „Wert nirgends", Reinheit/Aliasing/Nebenläufigkeit,
`BenchmarkBuildRowImage`); `mapper.go` update (`change` ruft die
Domänen-Funktion, `rowImage` samt `encoding/json`-Import entfallen, neuer
`columnNames` einmal je Änderung für beide Bilder, `containsColumn` bleibt —
Aufrufer `mapper.go:445` in `ExcludeColumn`); `mapper_test.go` unverändert.
**Ungeplant:** keine Datei außerhalb dieser Liste (Plan und Review-Report sind
Slice-Artefakte). **Nicht erfolgt:** nichts aus §1 „Ausdrücklich NICHT"
berührt — kein `origin`, keine Regel-Schnittstelle, keine geänderte Bildform,
kein Snapshot-Leser.

**Semantik-Vergleich alt gegen neu (Diff gelesen):** die Schleife ist Zeile für
Zeile dieselbe; einziger Unterschied ist `relation.Columns[i].Name` →
`columns[i]` und `containsColumn` → `containsName` (identischer Rumpf,
`mapper.go:495` gegen `rowimage.go:60`).

**Byte-Gleichheit selbst gemessen.** Ich habe das alte `rowImage` am Parent
gefahren: ein `git worktree` von `43f394a1` **außerhalb** des Repo-Baums (im
Scratchpad), darin zwei nicht committete Testdateien — eine Shim-Funktion
`BuildRowImage(columns, values, excluded)`, die eine `decode.Relation` baut und
das **alte** `rowImage` ruft, und die unveränderte `rowimage_test.go` des Kopfes
(nur `package model` → `package mapper`). Ergebnis: `make test` im Worktree
EXIT=0, und `go test -race -run BuildRowImage -v` (identisches Toolchain-Image
und Aufruf wie das `make test`-Rezept, auf ein Paket eingeschränkt) zeigt alle
Fälle von `TestBuildRowImageBytes` sowie
`TestBuildRowImageNilValuesIsAbsent`, `…ExcludedValueNowhere`,
`…PureAndConcurrent` als `PASS`. Die Referenz-Erwartungen sind damit echte
Parent-Charakterisierung. Worktree danach entfernt (`git worktree list` nennt nur
den Hauptbaum, `git status --short` leer).

**§3.13-Suchlauf-Feld — gegen beide Stände (`43f394a1` und `HEAD`) selbst
nachgemessen:**

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| 1 (Go-Code) | `git grep -n 'rowImage' 43f394a1 -- '*.go'` → **16** (`mapper.go` 5 an Z. 233/237/471/515/524, `natsstream/publisher.go` 4, `http/sse.go` 4, `http/readchanges.go` 2, `examples/nats-stream-client/format_test.go` 1); `… HEAD -- '*.go'` → **15** (dieselben elf fremden + 4 in `rowimage_test.go`, Z. 18/76/93/107); `grep -rn 'rowImage' --include=*.go .` → 15 (Arbeitsbaum sauber) | **bestätigt** (16 − 5 + 4 = 15) |
| 2 (Dokumente) | `git grep -n 'rowImage' <Stand> -- '*.md'` ohne die Plan-Datei: Parent und Kopf identisch — die vier ADR-Dateien, `welle-backfill-bestand.md:187`, `slice-transformationen-kern-rename.md:240`, `done/welle-18-results.md:40`, zwei `docs/reviews/review-slice-nats-drittstream-example-*.md`; kein Treffer in `spec/`, `docs/user/`, `harness/`, `README.md`, `AGENTS.md` (Ort-/Konstruktions-Suche, §2 Zeile 7) | **bestätigt**; die beiden Plan-Treffer bleiben wahr bzw. zitieren `ADR-0112` wörtlich |
| 3 (`Accepted` ADRs) | `git grep -n 'rowImage' <Stand> -- 'docs/plan/adr/*.md'`: Parent 10, Kopf 10; `0059` Z. 59/60/152/240, `0111` Z. 88/228, `0112` Z. 89/170/240/359 | **bestätigt**; `make doc-immutable` grün, keine ADR im Diff |
| 4 (zweite JSON-Erzeugung) | `json.Marshal`: je Stand 6 Treffer, Fundstellen wie im Plan (§2 Zeile 1); `WriteByte('{')`: je Stand ein Treffer (Parent `mapper.go:529`, Kopf `rowimage.go:31`). Zusätzlich `json.NewEncoder`/`RawMessage` über `*.go` am Kopf gelesen: nur Antwort-Encoder der HTTP-Handler und die Nachrichten-Strukturen; keine zweite Bild-Konstruktion | **bestätigt** |

Zahlen ohne auflösbaren Ursprung im Suchlauf-Feld: keine — jede trägt Befehl und
Stand (`gemessen`). Die Wirkungs-Zahlen in §6 „Leistung" sind als
**Review-Messung, nicht committet** gekennzeichnet (Instanz A: übernommen und
benannt).

**Fremde offene Pläne, die die bewegte Eigenschaft nennen:** die Pläne
`slice-transformationen-kern-rename` (Z. 25/50/74/161/165/240),
`slice-transformationen-backfill-pfad` (Z. 42/77) und
`slice-backfill-snapshot-reader` (Z. 56/114) führen die gemeinsame Funktion
unter dieser Bezeichnung mit „am Start gemessen"; keine dieser Stellen wird
durch den Diff falsch (Namen und Ort sind dort ausdrücklich noch offen). Der Plan
meldet die zwei `rowImage`-Nennungen und ändert sie nicht — regelkonform
(§3.13: gemeldet statt still mitgeändert).

**Lifecycle (`AGENTS.md` §3.3):** beide Moves sind reine Renames —
`git show --stat -M 14653d1d`: `{open => next}/slice-backfill-row-image-gemeinsam.md | 0`;
`git show --stat -M ffbd973f`: `{next => in-progress}/… | 0`, je „0 insertions(+),
0 deletions(-)". Die Inhaltsänderung „Verantwortlich gesetzt" (`47b9ea17`) liegt
in einem eigenen Commit **zwischen** den Moves.

## 4. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2 | Bild-Konstruktion ist **eine** Funktion, die WAL- und Backfill-Pfad gemeinsam aufrufen; „eine Stelle, die beide Adapter-Schichten importieren dürfen"; eine zweite JSON-Erzeugung wäre ein zweiter Träger | Ort `internal/domain/model` — `.a-check.yml` erlaubt `adapters → domain` und `app → domain`; `make a-check` grün; Suchlauf (§2 Zeile 1): eine Konstruktionsstelle; die Byte-Gleichheit ist am Parent reproduziert (§3) | konform |
| [`ADR-0016`](../plan/adr/0016-jsonb-row-images-mvp.md) | JSON-Row-Image, Text-Stand ohne Typ-Interpretation | `rowimage.go` serialisiert `*values[i]` (String) mit `json.Marshal`, keine Typumwandlung; Test „Maskierung" und „ungültiges UTF-8" tragen die `encoding/json`-Bytes | konform |
| [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 3 | Ausschluss im Bild, **vor** Serialisierung; kein Wert eines ausgeschlossenen Schlüssels | `rowimage.go:34` überspringt die Spalte, bevor `json.Marshal` für Name oder Wert läuft; `TestBuildRowImageExcludedValueNowhere` prüft, dass weder Schlüssel noch Wert im Bild stehen; Mapper reicht `binding.ExcludedColumns` an **beide** Bilder durch (`mapper.go:233/237`) | konform |
| [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 7 / Kopplung K1 der Welle | die Regelauswertung hängt an der **einen** gemeinsamen Funktion; Ausschluss gilt zuerst | Es gibt genau eine Bild-Funktion, in der die Auswertung ansetzen kann; die Schleife trägt Ausschluss-Prüfung vor jedem Wert-Zugriff. Kein Vorbau (keine Regel-Parameter, kein Interface). Die Signatur ist positional: ein Regel-Parameter änderte die Aufrufer — heute zwei in `mapper.go`; die Folge-Pläne sehen das ausdrücklich vor („reicht … die leere Regelmenge", `slice-transformationen-kern-rename` Z. 74/165). Folgepflicht 7 selbst (Backfill-Pfad trägt die Regelauswertung) ist der Kopplung K2 zugewiesen, nicht dieser Slice | konform (V-2) |

`ARC-001` (Domain Core: reine Funktionen): `BuildRowImage` liest nur ihre
Eingaben und hält keinen Zustand; `TestBuildRowImagePureAndConcurrent` deckt
Aliasing (Ergebnis-Mutation wirkt nicht auf ein zweites Ergebnis, Eingaben
unverändert) und acht nebenläufige Aufrufe unter `-race`.

## 5. Mutationen (Eingabeseite, selbst gesehen)

Je Mutation: Edit per `sed`, `make test` (Race-Detector), Exit-Code und
`FAIL`-Zeilen gelesen, danach `git checkout -- <Datei>`; `git diff --stat` und
`git status --short` nach jeder Rücknahme leer.

| # | Mutation | `make test` | Rote Tests |
|---|---|---|---|
| M1 | NULL-Guard aus `rowimage.go:34` entfernt (`values[i] == nil ||`) | **EXIT=2** | `TestBuildRowImageBytes` (Domäne) und `TestConsumeFullTransaction` (Mapper): nil-Pointer-Panic |
| M2 | Ausschluss-Prüfung aus `rowimage.go:34` entfernt (`containsName(excluded, column)`) | **EXIT=2** | Domäne `TestBuildRowImageBytes`/`…ExcludedValueNowhere`; Mapper u. a. `TestConsumeExcludedColumnAbsentFromRowImages`, `TestConsumeExcludedColumnSurvivesSchemaBump`, `TestExcludeColumnFiltersLiveBinding`, `TestIncludeColumnRestoresLiveBinding`, `TestAddBindingKeepsExclusionState`; zusätzlich Paket `internal/bootstrap` |
| M3 | Ausschluss-Durchreichung im Mapper für das alte Bild gekappt (`mapper.go:237`: `binding.ExcludedColumns` → `nil`) | **EXIT=2** | `TestConsumeExcludedColumnAbsentFromRowImages` |
| M4 | alter Bild-Aufruf liest `event.New` statt `event.Old` (`mapper.go:237`) | **EXIT=2** | `TestConsumeFullTransaction`, `TestConsumeExcludedColumnAbsentFromRowImages` |

Alle vier Mutationen färben rot; die Zusagen sind an ihrer Eingabeseite
gebunden, nicht nur an der Ausgabe. Baum nach M4: sauber, `mapper.go` und
`rowimage.go` byte-gleich dem Kopf.

## 6. Leistung / §6-Risiko

Plan §6 „Leistung im Backfill" nennt die Wirkung (**eine zusätzliche
Allokation je WAL-Änderung**, Ursprung „Code-Lesung", Messung als
„Review-Messung, nicht committet" gekennzeichnet). Vom Code gedeckt:
`mapper.go:232` ruft `columnNames(event.Relation)` **einmal** je Änderung und
teilt die Liste über beide Bild-Aufrufe (`mapper.go:233/237`);
`columnNames` legt genau ein `make([]string, len(relation.Columns))` an
(`mapper.go:517–523`), das ist eine Allokation je Änderung. `BuildRowImage`
selbst hat dieselbe Schleife wie das alte `rowImage` (§3, Semantik-Vergleich) —
keine zusätzliche Allokation im Rumpf. Die genannte Differenz +96 B je
Änderung passt rechnerisch zu sechs Spalten × 16 B Stringkopf. Die Zahlen
64 → 65 allocs/op und 1313 → 1409 B/op habe ich **nicht** nachgemessen — sie
stehen im Plan mit ihrem Ursprung und ohne Anspruch auf Beleg im Repo; das ist
nach §3.12 Instanz A zulässig. Der Ausgang des Risikos bleibt bei Closure.

Risiko „Byte-Abweichung durch Escaping": durch die Referenz-Bytes mit
`<`/`>`/`&` am Parent belegt (§3). Risiko „Ort der Funktion": `make a-check`
und der Import-Blick auf `rowimage.go` (nur Standardbibliothek).

## 7. Harte Regeln

- **§3.3** — Moves rein (§3 dieses Reports).
- **§3.5** — kein `Accepted`-ADR im Diff.
- **§3.7** — diff-skopierter Kandidatenlauf über die `+`-Zeilen der geänderten
  `*.go` (Chronik-/Vorher-Vokabular, Konjunktive): drei Treffer, davon zwei
  Testfall-Namen („NULL-Wert entfällt", „… entfällt samt Wert" — Indikativ über
  den Zustand, kein Befund) und ein Test-Kopfkommentar: V-1. Die Produktions-
  Kommentare (`BuildRowImage`, `containsName`, `columnNames`, `appendExcluded`)
  beschreiben den Ist-Zustand und tragen keine Slice-/Wellen-Kennung.
- **§3.11** — keine host-lokalen Pfade in den geänderten Dateien; `docs-check`
  0 Befunde (`hostpaths`).
- **§3.12** — Zahlen im Suchlauf-Feld tragen Befehl, Stand und „gemessen"; die
  Leistungs-Zahlen sind als übernommen gekennzeichnet (§6).
- **§3.13** — Suchlauf-Feld vollständig, beide Stände, Gefundenes und
  Nichtgefundenes je Träger (§3); die fremden offenen Pläne sind gemeldet.
- **§3.9** — alle Läufe dieses Reports mit separatem Exit-Code (§1, §5).
- **Commit-Traceability** — 10 Commits, jede Message nennt `LH-FA-CAP-009` und
  `ADR-0111`, keine `SPEC-`/`ARC-`-Kennung im Betreff (§1).

## 8. Beobachtungen (nicht blockierend)

- **V-1 (INFO) — Vorher-Formulierung im Test-Kopfkommentar.**
  `internal/domain/model/rowimage_test.go:9–11`: „gemessen am Bild-Erzeuger des
  Replication-Mappers vor seinem Umzug in die Domäne". Der Satz beschreibt einen
  Vorher-Zustand (`AGENTS.md` §3.7: Herkunft in ein auflösbares Feld); die
  Provenienz ist mit [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 2 verankert und von mir reproduziert (§3). Deckt sich mit
  Review-Finding F-3 (INFO, keine Aktion); Klasse „Vorher-Formulierung in
  Test-Provenienz".
- **V-2 (INFO) — Signatur ist positional.** Ein Regel-Parameter (Kopplung K1)
  ändert jeden Aufrufer. Heute sind es zwei Aufrufe in `mapper.go:233/237`;
  der künftige Backfill-Aufrufer kommt ohnehin neu hinzu. Die Folge-Pläne
  (`slice-transformationen-kern-rename` Z. 74/165) benennen die Anpassung. Keine
  Abweichung, kein Vorbau, kein Handlungsbedarf in diesem Slice; Beurteilung wie
  Review-Finding F-7.
- **V-3 (INFO) — Coverage-Rauschen.** Dieselbe Fläche misst 82.60 % (im
  `make gates`-Lauf) und 82.70 % (im Einzellauf); beide sind gedruckte Zeilen
  konkreter Läufe (§1). `BuildRowImage` trägt 90.9 % (die zwei
  `json.Marshal`-Fehlerzweige sind für String-Eingaben nicht erreichbar); das
  Gate kennt keine Ausnahme-Zeile (`AGENTS.md` §3.2), die Schwelle wird
  eingehalten.
- **V-4 (INFO) — zeitrelative Diff-Stand-Zelle.** Der Suchbefehl der Plan-Zeile 1
  für den „Diff-Stand" nennt `HEAD` (mit dem Zusatz „Arbeitsbaum sauber"). Er
  bleibt wahr, solange kein weiterer Commit `rowImage` in `*.go` bewegt; danach
  wäre `43f394a1..`-Bezug oder ein Commit-Anker stabiler. Kein Befund am
  gelieferten Stand.

## 9. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (7) | **7 von 7 erfüllt**, je mit eigenem Beleg |
| DoD §2 — „Review durchgeführt" | **materiell erfüllt** (Report liegt vor, Fixrunde selbst nachgemessen) — Haken beim Closure-Nachzug |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-/Rollen-Sequenz) |
| Plan-vs-Code-Diff | **deckungsgleich**, keine ungeplante Datei, kein stilles Streichen |
| Byte-Gleichheit | **am Parent reproduziert** (alte Funktion, Referenz-Erwartungen des Kopfes: grün) |
| Eine Konstruktionsstelle | **bestätigt** (ein `WriteByte('{')`, zwei Aufrufe von `BuildRowImage`) |
| Mutationen (Eingabeseite) | **4 von 4 rot**, Baum exakt zurückgenommen |
| Entscheidungen `ADR-0111`/`0016`/`0059`/`0112` | **konform** |
| Harte Regeln (§3.3, §3.5, §3.7, §3.9, §3.11, §3.12, §3.13) | **erfüllt** |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=43f394a1..HEAD`) |
| Gates | **`make gates` EXIT=0**, `make test` EXIT=0, `make coverage-gate` EXIT=0 (`82.7%`), `make a-check` EXIT=0 |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung
belegt: die Row-Image-Konstruktion liegt an genau einer Stelle in der Domäne,
der WAL-Pfad ruft sie, das Ergebnis ist am Parent-Stand byte-gleich reproduziert,
die Schichtenregel hält, und die vier Eingabeseiten-Mutationen färben rot. Die
Zahlen des §3.13-Suchlauf-Felds stimmen an beiden Ständen. Kein V-Befund mit
Blockade; V-1…V-4 sind INFO-Beobachtungen ohne erwartete Aktion.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: der Aufrufer hakt die
Zeile „Review durchgeführt", schreibt die Closure-Notiz mit Lerneintrag (die
Finding-Klassen des Reviews gehen dort ins Beobachtungs-Register), trägt die
Ausgänge der vier §6-Risiken (das Leistungs-Risiko mit dem gemessenen Zuschlag
samt Lauf) und die Welle-Paarungen; danach der reine `git mv` nach `done/`
(`AGENTS.md` §3.3, Fall 2).
