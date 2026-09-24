# Review-Report: slice-backfill-row-image-gemeinsam — 2026-09-24

**Review-Art:** Code — der Diff verschiebt die Row-Image-Konstruktion als reine
Funktion in die Domäne und lässt den WAL-Pfad sie rufen; geprüft gegen Plan,
ADRs und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein
DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `43f394a1..f0d18d91`, sieben Commits `14653d1d` …
`f0d18d91`; vier Dateien: `internal/domain/model/rowimage.go` (neu),
`internal/domain/model/rowimage_test.go` (neu),
`internal/adapters/driving/replication/mapper/mapper.go` (geändert),
Slice-Plan `slice-backfill-row-image-gemeinsam` (Lifecycle `open` → `next` →
`in-progress`, Inhalt). Nachprüfung der Fixrunde gegen `HEAD` (Commits
`dc633a97`, `cd046787`), siehe §Nachprüfung.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Zusage-ohne-Eingabeseite, Kommentar-Chronik).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-row-image-gemeinsam` (§1 Ziel, §3 Plan und
  Suchlauf, §6 Risiken) und Welle `welle-backfill-bestand` (§5 Kopplung K1)
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  (Accepted) Teilfrage 2 (Row-Image-Parität: eine Funktion),
  [`ADR-0016`](../plan/adr/0016-jsonb-row-images-mvp.md),
  [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 3,
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 2 und 7 (Kopplungs-Bezug, kein Umfang dieses Slice)
- [`LH-FA-CAP-008`](../../spec/lastenheft.md),
  [`LH-FA-CFG-005`](../../spec/lastenheft.md),
  [`LH-FA-DAT-005`](../../spec/lastenheft.md),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md),
  [`LH-FA-CAP-009`](../../spec/lastenheft.md);
  [`SPEC-002`](../../spec/pflichtenheft.md),
  [`ARC-001`](../../spec/architecture.md),
  [`ARC-006`](../../spec/architecture.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.7, §3.9, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen):

- **Zeilenweiser Vergleich** des alten `rowImage`
  (`git show 43f394a1:internal/adapters/driving/replication/mapper/mapper.go`)
  gegen `internal/domain/model/rowimage.go`: der Rumpf ist bis auf
  `relation.Columns`/`column.Name` → `columns`/`column` und
  `containsColumn` → `containsName` identisch. Gleiche Semantik bei
  `nil`-Werteliste (`nil, nil`), leerer Liste (`{}`), `nil`-Wert,
  ausgeschlossener Spalte (Guard vor dem Marshal, Wert nie serialisiert),
  weniger Werten als Spalten (`i >= len(values)`), mehr Werten
  (unbeachtet). Escaping: beide Fassungen nutzen `json.Marshal` für Name und
  Wert, kein `SetEscapeHTML`, kein Encoder — die Standard-Maskierung (`<`,
  `>`, `&`, U+2028) ist unverändert; Trennzeichen und Reihenfolge
  identisch. Einziger Aufrufer-Unterschied: `columnNames(event.Relation)`
  läuft jetzt vor dem `nil`-Test der Werteliste — `event.Relation` ist an
  dieser Stelle bereits über `event.Relation.QualifiedName()` dereferenziert,
  kein neuer Nil-Pfad.
- **Charakterisierung am Parent-Stand selbst reproduziert:** temporärer
  `git worktree` von `43f394a1` außerhalb des Repo-Verzeichnisses, darin eine
  nicht committete Testdatei im Paket `mapper`, die die neun Byte-Fälle, den
  `nil`-Fall, den Ausschluss-Fall und den Reinheits-/Nebenläufigkeitsfall aus
  `rowimage_test.go` gegen das **alte** `rowImage` fährt (Docker, `go test
  -race`, Toolchain-Image des Makefile): alle Fälle grün, Exit 0. Die
  Referenz-Bytes der Tests sind damit echte Parent-Charakterisierung. Der
  Worktree ist entfernt (`git worktree list` zeigt nur das Repo).
- **Mutationen der Eingabeseite selbst gesetzt** (je Edit, `make test`, exakt
  zurückgenommen): acht Mutationen, alle rot, siehe §Mutationen.
- **Gates am Stand `f0d18d91`:** `make test` (Race-Detector) Exit 0;
  `make a-check` Exit 0 („gesamt: 0 Befund(e)"); `make coverage-gate` Exit 0
  („Coverage 82.60% erfüllt Schwelle 80%"); `make gates` Exit 0. Die Domäne
  importiert nur `bytes` und `encoding/json`.
- **Reste der alten Konstruktion:** `git grep -n 'json.Marshal' HEAD --
  'internal/*.go'` ohne Test-Treffer: `model/rowimage.go` (2×, die eine
  Funktion), `natsstream/publisher.go:228` und `http/sse.go:114`
  (Wire-Nachricht, bettet das fertige Bild als `json.RawMessage` ein);
  `git grep -n "WriteByte('{')" HEAD -- internal`: ein Treffer
  (`model/rowimage.go:31`). Keine zweite Bild-Konstruktion.
- **Zahlen des Suchlauf-Feldes** (Plan §3) mit `git grep` an `43f394a1` und
  am Kopf nachgemessen: siehe F-1 und §Nachprüfung.
- **Leistung** an beiden Ständen gemessen: siehe F-4.
- **Commit-Struktur:** `git show --stat -M` über `14653d1d` und `ffbd973f`:
  reine Renames ohne Zeilenänderung; `47b9ea17` (Feld `Verantwortlich`) und die
  Inhalts-Commits liegen getrennt davon — `AGENTS.md` §3.3 erfüllt. Betreffs
  tragen `LH-FA-CAP-009` und `ADR-0111`, keine `SPEC-`/`ARC-`-Kennung.

### Mutationen (Eingabeseite, selbst ausgeführt)

Datei nach jeder Mutation per `git checkout` zurückgenommen; Endstand
`git diff` leer, `git status` sauber.

| # | Mutation | Ort | `make test` |
|---|---|---|---|
| M1 | Mapper reicht dem New-Bild `nil` statt `binding.ExcludedColumns` | `mapper.go` (`change`) | rot: `TestConsumeExcludedColumnAbsentFromRowImages`, `TestConsumeExcludedColumnSurvivesSchemaBump`, `TestExcludeColumnFiltersLiveBinding`, weitere |
| M5 | dasselbe für das Old-Bild | `mapper.go` (`change`) | rot: `TestConsumeExcludedColumnAbsentFromRowImages` |
| M6 | Old-Bild aus `event.New` gebaut | `mapper.go` (`change`) | rot: `TestConsumeFullTransaction`, `TestConsumeExcludedColumnAbsentFromRowImages` |
| M2 | NULL-Guard `values[i] == nil` entfernt | `rowimage.go` | rot: `TestBuildRowImageBytes/NULL-Wert_entfällt`, `TestConsumeFullTransaction` (Panik) |
| M3 | Längen-Guard `i >= len(values)` entfernt | `rowimage.go` | rot: `TestBuildRowImageBytes/leere_Werteliste…` (Index-Panik), `TestConsumeRelationBackfillsMissingTableSchema` |
| M4 | `containsName` liefert immer `false` | `rowimage.go` | rot: sechs Mapper-Ausschluss-Tests |
| M7 | Wert ohne `json.Marshal` (rohe Anführungszeichen) | `rowimage.go` | rot: Maskierungs- und UTF-8-Fall |
| M8 | Trennzeichen `,` → `;` | `rowimage.go` | rot: `TestConsumeFullTransaction` und weitere |

Die Zusagen (Ausschluss, NULL-Abwesenheit, Längen-Grenze, Escaping) sind an
ihrer Eingabeseite gebunden — im Domänen-Test (Funktionseingabe) wie im
Mapper-Test (Durchreichung von `ExcludedColumns` je Bild).

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Suchlauf-Feld (Plan §3): Gesamtzahl „17 Treffer" gegen die Messung, und der Parent-Befehl nennt `HEAD`

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz A (Klasse „Zahl im Träger ohne Ursprung —
  oder gegen die Messung driftend" der HIGH-Liste des Skills); zugleich „Beleg
  trägt seinen Satz nicht"
- `pfad`: Slice-Plan `slice-backfill-row-image-gemeinsam` §3, Zeile
  „Go-Code und Kommentare" des Suchlauf-Feldes (Stand `f0d18d91`, Zeile 124)
- `befund`: Das Feld nennt „Parent: 17 Treffer", zerlegt in 5 + 4 + 6 + 1;
  `git grep -n 'rowImage' 43f394a1 -- '*.go'` liefert 16 Zeilen — die
  Summanden stimmen, die Summe nicht. Der als Parent-Befehl genannte
  `git grep -n 'rowImage' HEAD -- '*.go'` liefert am Kopf 15 (11 fremde plus 4
  in `rowimage_test.go`), nicht den Parent-Stand; er reproduziert weder 17
  noch 16.
- `verifizierbar`: ja — `git grep -n 'rowImage' 43f394a1 -- '*.go' | wc -l`
  → 16
- `klasse`: Zahl im Träger driftet gegen die Messung

### F-2 — `BuildRowImage`-Doc: „der Backfill-Pfad ruft dieselbe Funktion" im Indikativ Präsens

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Zusage nur als Zusage) / §3.12 Instanz B
- `pfad`: `internal/domain/model/rowimage.go:22-25` (Stand `f0d18d91`)
- `befund`: Der Kommentar sagt, der WAL-Pfad **und** der Backfill-Pfad rufen die
  Funktion; einen Backfill-Pfad gibt es im Baum noch nicht
  (`git grep -n BuildRowImage f0d18d91 -- internal`: ein Aufrufer, der
  Mapper). Eine Zusage der Welle steht als bestehender Zustand; der Bezug auf
  `ADR-0111` Teilfrage 2 ist der tragende Rang-Zeiger.
- `verifizierbar`: ja — `git grep -n BuildRowImage <Stand> -- internal`
- `klasse`: Kommentar stellt Zusage als bestehenden Zustand dar

### F-3 — Test-Kopfkommentar: „vor seinem Umzug in die Domäne"

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7 (abwesender Code als Bezugsgröße); Abgrenzung
  Testfall-Provenienz des Skills (Subjekt ist der Test, zulässig)
- `pfad`: `internal/domain/model/rowimage_test.go:6-12`
- `befund`: Der Kopfkommentar begründet die Erwartungswerte mit ihrer Messung am
  Bild-Erzeuger „vor seinem Umzug"; der Erzeuger in dieser Form existiert nicht
  mehr. Die Herkunft ist ein Testfall-Beleg (Subjekt: die Erwartungswerte) und
  damit zulässig; die Aussage ist wahr — die Reproduktion dieses Reviews
  bestätigt sie gegen den Parent.
- `verifizierbar`: nein
- `klasse`: Vorher-Formulierung in Test-Provenienz (kein Zähler)

### F-4 — Leistung: +1 Allokation und +96 B je WAL-Änderung durch `columnNames`

- `kategorie`: INFO
- `quelle`: Plan §6 „Leistung im Backfill"; `AGENTS.md` §3.12 (Ursprung der
  Zahlen)
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:232-233`
  und `:515-523` (Stand `f0d18d91`)
- `befund`: Selbst gemessen (Go-Benchmark in zwei temporären Worktrees,
  `-benchtime 300000x -count 3`, 6 Spalten, eine ausgeschlossen, beide Bilder
  je Änderung): Parent 64 allocs/op, 1313 B/op, 3324–3736 ns/op; Kopf
  (`columnNames` + 2× `BuildRowImage`) 65 allocs/op, 1409 B/op, 3330–3558
  ns/op — +1 Allokation, +96 B (6 Spalten × 16 B Slice-Header), Zeit im
  Rauschen. `BuildRowImage` selbst allokiert wie das alte `rowImage` (gleicher
  Rumpf); der Zuschlag entsteht nur im WAL-Adapter, der Backfill-Pfad trägt ihn
  nicht. Die Implementer-Zahlen (43→44 Allokationen, 1642→1739 B) standen in
  keinem committeten Träger. Bewertung: tragbar — eine Vermeidung (Namensliste
  je `decode.Relation` im `Assembler` vorhalten oder eine zweite Eingabeform)
  brächte Zustand und eine zweite Signatur gegen rund 1,5 % der Allokationen;
  Aufwand größer als Nutzen. Der Ausgang des §6-Risikos „Leistung" trägt bei
  Closure den gemessenen Zuschlag samt Lauf.
- `verifizierbar`: ja — `go test -bench` je Stand (nicht als Make-Target
  geführt)
- `klasse`: Zahl ohne auflösbaren Träger (außerhalb des Diffs, INFO)

### F-5 — Dritte Kopie der Namens-Suche (`containsName`)

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/domain/model/rowimage.go:57-66`;
  `internal/adapters/driving/replication/mapper/mapper.go:492-500`
  (`containsColumn`, bleibt für `ExcludeColumn`);
  `internal/adapters/driven/postgresstorage/sqlexec/translate.go:204-212`
  (`containsColumn`)
- `befund`: Die lineare Namens-Suche steht jetzt dreimal (zwei
  `containsColumn` im Bestand, neu `containsName`). Ein Reuse der
  Domänen-Fassung im Mapper bräuchte eine exportierte Domänen-Funktion für
  sechs Zeilen; ein Schichtenbruch entstünde nicht (`adapters → domain` ist
  erlaubt), wohl aber eine neue öffentliche Oberfläche ohne zweiten
  fachlichen Grund. Das Paket `slices` wird im Bestand nirgends benutzt
  (`grep -rl '"slices"' internal` → leer). Duplikat vertretbar.
- `verifizierbar`: nein
- `klasse`: kleine Helfer-Duplikate

### F-6 — Namensgleiche `rowImage` in `natsstream` und `http` (Verständnis-Risiko)

- `kategorie`: INFO
- `quelle`: Maintainability; `AGENTS.md` §3.13 (Suchlauf-Ergebnis)
- `pfad`: `internal/adapters/driven/natsstream/publisher.go:156`;
  `internal/adapters/driving/http/sse.go:48`
- `befund`: Beide Funktionen betten ein fertiges Bild als `json.RawMessage` ein
  und konstruieren nichts; sie sind paketprivat und liegen in anderen
  Schichten. Der Diff **verringert** die Verwechslungsgefahr (das
  konstruierende `rowImage` ist entfallen, der Konstruktionsort trägt jetzt den
  unterscheidbaren Namen `BuildRowImage`). Kein Befund; der Plan (§3) meldet
  die Namensgleichheit selbst als Beobachtung.
- `verifizierbar`: nein
- `klasse`: Namensgleichheit ohne Funktionsgleichheit

### F-7 — Kopplung K1: Signatur trägt die Regelauswertung ohne Vorbau

- `kategorie`: INFO
- `quelle`: Plan §6 „Kopplung an die Transformations-Umsetzung";
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 2 (Auswertung „in `rowImage` nach dem Ausschluss")
- `pfad`: `internal/domain/model/rowimage.go:27`;
  `internal/adapters/driving/replication/mapper/mapper.go:232-241`
- `befund`: Die Funktion iteriert Spalten in Relation-Reihenfolge, wendet den
  Ausschluss an und schreibt erst dann — der Ansatzpunkt „nach dem Ausschluss"
  liegt in ihr, nicht in einer Kopie. Ein weiterer Parameter (Regeln) ändert
  eine Signatur mit zwei Aufrufstellen im Mapper (je Bild eine; `columns` ist
  bereits einmal je Änderung berechnet) plus die künftige
  Backfill-Aufrufstelle; die Nichtanwendbarkeits-Prüfung gegen die
  Relation-Spalten braucht `columns`, das der Mapper schon hält. Kein Vorbau im
  Diff, keine Regel-Schnittstelle. Der Wortlaut der `Accepted` ADR nennt den Ort
  noch als `rowImage`; die Welle-Kopplung K1 („die eine gemeinsame Funktion")
  gilt unverändert, der Umzug ist im Suchlauf gemeldet, nicht in der ADR
  nachgezogen (richtig nach `AGENTS.md` §3.5).
- `verifizierbar`: nein
- `klasse`: Kopplung an Folge-Welle (Beurteilung, keine Abweichung)

---

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model/rowimage.go` — Rumpf gegen das
  alte `rowImage` zeilenweise identisch in Semantik und Bytes; rein (kein
  Zustand, kein geteilter Speicher — `TestBuildRowImagePureAndConcurrent` mit
  `-race`), importiert nur Standardbibliothek; keine Slice-/Wellen-Nummer in
  einem Produktionscode-Kommentar (`AGENTS.md` §3.7; einziger Vorbehalt F-2).
- geprüft, ohne Befund: `internal/domain/model/rowimage_test.go` — die
  Referenz-Bytes stammen aus dem Parent (am Parent-Stand reproduziert); die
  Zusagen sind an ihrer Eingabeseite gebunden (acht Mutationen rot, darunter
  NULL-Guard, Längen-Guard, Ausschluss-Durchreichung je Bild,
  `containsName`); kein Fake, dessen Rückgabe geprüft würde. Sonderzeichen
  stehen in Hex-Escapes; keine host-lokalen Pfade (`AGENTS.md` §3.11).
- geprüft, ohne Befund: `internal/adapters/driving/replication/mapper/mapper.go`
  — `rowImage` und `encoding/json`-Import entfallen; `containsColumn` bleibt
  begründet (Aufrufer `ExcludeColumn`); `appendExcluded`-Kommentar nennt
  `model.BuildRowImage`; `columnNames`-Kommentar beschreibt den Ist-Zustand;
  `mapper_test.go` im Diff unberührt, die bestehenden Tests laufen ohne
  geänderte Erwartung grün.
- geprüft, ohne Befund: Schichten — `make a-check` Exit 0; die Domäne
  importiert nichts aus anderen Schichten; der Mapper importiert
  `internal/domain/model` (`adapters → domain`). `AGENTS.md` §3.1 Docker-only:
  alle Läufe über `make` bzw. `docker run` des gepinnten Images; kein Host-Go.
- geprüft, ohne Befund: Slice-Plan — Abweichungen (`containsColumn` bleibt im
  Mapper, `containsName` als privater Helfer, neuer `columnNames`) stehen in §3
  Tabelle und sind zum Diff wahr; die übrigen Zahlen des Suchlauf-Feldes
  nachgemessen und bestätigt: `mapper.go` 5 Treffer (Z. 233/237/471/515/524 am
  Parent), `natsstream` 4, `http` 6, `examples` 1, elf fremde Treffer am Kopf,
  zehn `.md`-Treffer der `Accepted` ADRs (`0059`: 4, `0111`: 2, `0112`: 4) mit
  den genannten Zeilen, `welle-backfill-bestand` Z. 187,
  `slice-transformationen-kern-rename` Z. 240, `welle-18-results` Z. 40,
  `json.Marshal`-Fundstellen und `WriteByte('{')` (ein Treffer je Stand) — bis
  auf die Summe aus F-1 alle bestätigt. Kein Treffer in `spec/`, `docs/user/`,
  `harness/`, `README.md`, `AGENTS.md` — §3.13 erfüllt; die fremden offenen
  Pläne sind gemeldet, nicht mitgeändert.
- geprüft, ohne Befund: Lifecycle und Commit-Struktur (Moves rein, Betreffs ohne
  `SPEC-`/`ARC-`-Kennung, jede Message mit `LH-FA-CAP-009`/`ADR-0111`);
  `AGENTS.md` §3.5 (keine `Accepted` ADR im Diff), §3.6 (keine Gate-Lockerung),
  §3.11.
- geprüft, ohne Befund: `docs/user/` — im Diff unberührt; keine neue
  Betreiber-Oberfläche, keine Handbuch-Versionshistorie berührt.
- Schwerpunkte laut Skill ohne Anwendungsfall im Diff: Suppression,
  Sicherheits-Anti-Pattern (der Ausschluss-Pfad ist funktional unverändert und
  mutationsgebunden), Workflow-Dateien, Spec-Stratum (Spec im Diff
  unberührt), Handbuch, Formvorbild-Kopie.

## Nachprüfung der Fixrunde (Kopf `cd046787`)

Die Commits `dc633a97` (Doc-Kommentar) und `cd046787` (Plan-Suchlauf) haben F-1,
F-2 und F-4 adressiert. Geprüft am Kopf, nicht aus dem Bericht übernommen:

- **F-1 behoben.** Das Feld nennt jetzt „Parent: **16** Treffer (gemessen)" und
  „Diff-Stand: **15** Treffer (gemessen)" mit den Befehlen
  `git grep -n 'rowImage' 43f394a1 -- '*.go'` (Parent) und
  `git grep -n 'rowImage' HEAD -- '*.go'` (Diff-Stand, „Arbeitsbaum sauber").
  Nachgemessen: Parent 16, Kopf 15, `grep -rn 'rowImage' --include=*.go .` 15;
  `json.Marshal` je Stand 6 Treffer (gemessen, wie im Feld genannt);
  `WriteByte('{')` je Stand ein Treffer (Parent `mapper.go` Z. 529, Kopf
  `model/rowimage.go` Z. 31); ADR-Treffer je Stand 10. Der Parent-Befehl trägt
  jetzt den Parent-Commit statt `HEAD`. Der Diff-Stand-Befehl mit `HEAD` ist
  eine Aussage über den jeweiligen Kopf; sie ist mit dem Zusatz „Arbeitsbaum
  sauber" gekennzeichnet und am aktuellen Kopf wahr.
- **F-2 behoben.** Der Doc-Kommentar lautet jetzt: „Sie ist die eine
  Konstruktionsstelle für Row Images (`ADR-0111` Teilfrage 2): jeder Pfad, der
  ein Row Image erzeugt, ruft sie; gerufen wird sie vom Replication-Mapper." Die
  Aussage ist wahr (`git grep -n BuildRowImage HEAD -- internal` ohne Tests:
  zwei Aufrufe in `mapper.go` Z. 233/237) und stellt die Backfill-Nutzung nicht
  mehr als Bestand dar. Der Diff `f0d18d91..HEAD` unter `internal/` besteht aus
  genau dieser Kommentar-Änderung (3 Zeilen), keine Verhaltensänderung.
- **F-4 gezogen.** Der Plan trägt jetzt den Zuschlag mit Wirkung: „eine
  zusätzliche Allokation je WAL-Änderung", `BuildRowImage` allokiert nicht mehr
  als das alte `rowImage`, „64 → 65 allocs/op, 1313 → 1409 B/op, ns/op im
  Rauschen" — deckungsgleich mit der Messung dieses Reviews.
- Verhalten unberührt: die Kommentar-Änderung in `rowimage.go` ändert keine
  Zeile Code; die Mutationsergebnisse und die Byte-Gleichheit gelten
  unverändert.

## Summary

Zählung der Findings des Erstlaufs (Stand `f0d18d91`); Nachprüfung siehe oben.

| Kategorie | Anzahl | davon nach Nachprüfung offen |
|---|---|---|
| HIGH | 1 | 0 |
| MEDIUM | 0 | 0 |
| LOW | 1 | 0 |
| INFO | 5 | 0 (ohne erwartete Aktion) |

**Finding-Klassen dieses Laufs:** Zahl im Träger driftet gegen die Messung ·
Kommentar stellt Zusage als bestehenden Zustand dar · Vorher-Formulierung in
Test-Provenienz · Zahl ohne auflösbaren Träger · kleine Helfer-Duplikate ·
Namensgleichheit ohne Funktionsgleichheit · Kopplung an Folge-Welle

## Verdikt

**Erstlauf (Stand `f0d18d91`):** merge-blockierend nach Skill-Liste — F-1 fällt
in die HIGH-Klasse „Zahl im Träger gegen die Messung"; Fixrunde am Implementer
nötig, F-2 (LOW) in derselben Runde.

**Nach der Fixrunde (Kopf `cd046787`):** F-1 und F-2 sind behoben und am Kopf
nachgemessen, F-4 ist im Plan mit Wirkung geführt; F-3 und F-5 bis F-7 sind INFO
ohne erwartete Aktion. **Keine offenen HIGH/MEDIUM/LOW, nicht mehr
merge-blockierend.** Der Diff ist verhaltensgleich zum Parent (Semantik
zeilenweise, Bytes am Parent reproduziert, acht Mutationen rot),
schichtenkonform und ohne Vorbau. Die Fixrunde lief mit dem Rückgabe-Pfeil
Reviewer → Implementer; die DoD-Zeile „Review durchgeführt" wird deshalb
regulär beim Implementer-Nachzug (Schritt 21) gezogen, nicht von diesem
Report (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug, Grenze). Offen
bleibt bei Closure: der Ausgang des §6-Risikos „Leistung" trägt den gemessenen
Zuschlag samt Lauf.

**Übergabe:** Finding-Klassen in die Slice-Closure §7. Dieser Report ist ein
**Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und
ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier
separat (Modul 11).
