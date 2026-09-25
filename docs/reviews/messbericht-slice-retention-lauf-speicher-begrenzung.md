# Messbericht: Retention-Lauf liest Kandidaten seitenweise — Nachmessung des Speichers

**Rolle:** Implementer (Modul 9) — Messbericht, **kein** Review und keine
Verifikation. **Datum:** 2026-09-25. **Gegenstand:** Nachmessung des Speichers des
Feed-Containers nach der Umsetzung von
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
(Bereinigungslauf liest Kandidaten seitenweise ohne Row Images) für
[`LH-FA-RET-002`](../../spec/lastenheft.md) bis
[`LH-FA-RET-004`](../../spec/lastenheft.md) und
[`LH-FA-CAP-009`](../../spec/lastenheft.md); Vergleichsbasis ist der
[Messbericht der Untersuchung](messbericht-slice-backfill-speicher-untersuchung.md)
(dort Reihen B und J). **Werkzeug:** `tools/bench-backfill-memory.sh` und
`tools/bench-scaling.sh` (Vertrag:
[`harness/targets/bench-backfill.md`](../../harness/targets/bench-backfill.md)); die
gedruckten Zeilen aller Läufe stehen unten, jede Zahl dieses Berichts nennt Lauf und
Zeile.

## 1. Ergebnis

1. **Der Speicher hängt nicht mehr an der Zahl der Changes in dem Maß, das die
   Untersuchung gemessen hat.** Die Spitze des Feed-Containers (`memory.peak`, gedruckt
   nach dem 60-s-Nachlauf) liegt bei 1.000.000, 2.000.000 und 3.000.000 Changes bei
   **14,9 bis 15,1**, **16,5 bis 16,9** und **17,3 bis 17,6 MiB** (zwei Läufe,
   `n` = 2 je Stufe; Läufe `20260925T183239Z` und `20260925T184503Z`, Zeilen „Zeilen in
   `cdc.change` nach dem Nachlauf“). Vor der Änderung standen dort 1.082,7 bis 1.273,5,
   2.269,2 bis 3.058,0 und 3.096,6 MiB (Messbericht der Untersuchung, Abschnitt 3.3,
   Reihen B und C). Die Zusage „Spitze im Bereich der Werte bei ausgeschalteter
   Bereinigung plus dem Bedarf einer Seite“ ist damit belegt, mit einer Einschränkung
   (Punkt 2).
2. **Ein Rest an Wachstum bleibt gemessen und unerklärt.** Die Werte ohne Bereinigung
   (Reihe J des Messberichts der Untersuchung, ein Lauf) liegen bei 13,3, 13,2 und
   13,7 MiB. Die Differenz der Nachmessung zu Reihe J steigt von 1,6 bis 1,8 MiB bei
   1.000.000 über 3,3 bis 3,7 bei 2.000.000 auf 3,6 bis 3,9 MiB bei 3.000.000 Changes
   (abgeleitet aus den Zeilen). Von 1.000.000 auf 3.000.000 Changes wachsen die Werte
   um 2,4 und 2,5 MiB (17,3 − 14,9 und 17,6 − 15,1), also etwa 1,3 Bytes je Change
   (abgeleitet: 2,5 MiB × 1.048.576 / 2.000.000); vor der Änderung waren es 1,03 bis 1,59 KiB je Change (Untersuchung,
   Abschnitt 3.3), das ist etwa das Eintausendfache. Die Ursache des Rests ist **nicht
   untersucht**; die Herleitung des Bedarfs einer Seite (etwa 1,5 MiB) sagt keinen Rest
   dieser Art voraus.
3. **Die Bereinigungs-Takte laufen bei 2.000.000 und 3.000.000 Changes weiter.** Die
   Zählung „Bereinigung gelaufen“ seit dem Start des Feeds steht nach dem Nachlauf bei
   20, 22 und 23 (Lauf 1) und 20, 23 und 23 (Lauf 2), fehlgeschlagen 0; der Nachlauf
   trägt in jedem der sechs Runs Garbage-Collection-Läufe (560 bis 2.467). Vor der
   Änderung blieben die Takte im dritten Run (2.000.000 Changes davor) aus: „Bereinigung
   gelaufen 5“, kein GC-Lauf im Nachlauf (Untersuchung, Abschnitt 3.3, Reihe C). Die
   Hypothese der Untersuchung, das Ausbleiben liege an der Dekodierung großer Listen im
   Feed-Prozess, ist damit nicht widerlegt, aber die Beobachtung tritt mit den seitenweise
   gelesenen Kandidaten nicht mehr auf (`n` = 2 Läufe).
4. **Skalierungs-Lauf** ([`LH-QA-PER-002`](../../spec/lastenheft.md)): `cdc_capture_lag`
   1,010 s, 1,004 s und 1,004 s bei den drei Lastenstufen, unter der bestehenden Grenze
   von 60 s ([`SPEC-013`](../../spec/pflichtenheft.md)); ein Lauf im verkürzten Modus
   gegen einen leeren `cdc.change`-Bestand (Abschnitt 4).
5. **Eingabeseiten-Mutationen** der Zusagen: Abschnitt 5.

## 2. Umgebung und Herkunft der Zahlen

**Host** (erste Ausgabezeile jedes Laufs): Linux 6.8.0-139-generic, Docker 29.8.1,
20 CPU, 33.362.599.936 Byte RAM; PostgreSQL 18 (`postgres:18-alpine`, Pin aus
`tools/bench-lib.sh`). **Feed-Image:** `ghcr.io/pt9912/pg-change-feed:dev`,
Image-ID `sha256:8c8dea000428448c0ff22a95f66061b32634beff80eeef1cf5c70e94c4a82a40`
(`make image` am Stand des Commits `0e6b1b30`; der Produktionscode dieses Standes ist der
des Slice-Diffs, danach änderte sich nur Doku). **Aufruf** beider
Läufe wie Reihe B der Untersuchung, damit der Vergleich dieselbe Anordnung hat:
`BENCH_MEM_STAGES=1000000 BENCH_FEED_ENV='GODEBUG=gctrace=1' BENCH_FEED_DOCKER_ARGS='--memory 6g'`,
drei Runs, ein frischer Feed-Container je Run, `cdc.change` vor Run 1 leer (Bestand vor
den Runs 0, 1.000.000 und 2.000.000 Changes, danach 3.000.000), schmale Zeilen (74 Bytes
`pg_column_size`).

**Herkunft.** *Gemessen*: gedruckt in einer Zeile der Läufe unten. *Abgeleitet*: aus
gemessenen Zahlen gerechnet, die Rechnung steht daneben. *Übernommen*: aus dem Messbericht
der Untersuchung (Reihen B, C und J, Läufe `20260925T130048Z`, `20260925T131359Z`,
`20260925T145726Z`), dort gedruckt; in diesem Lauf nicht nachgemessen. Die Läufe der
Nachmessung liefen nacheinander (nie zwei schwere Docker-Läufe zugleich); verfügbarer
Speicher laut `free -m` vor Lauf 1: 21.082 MB, vor Lauf 2: 20.350 MB, vor dem
Skalierungs-Lauf: 19.897 MB. Fremde Last auf dem Host ist nicht gemessen; während der
Läufe lief kein weiterer schwerer Docker-Lauf dieses Auftrags.

## 3. Nachmessung des Speichers

Spitze des Feed-Containers (`memory.peak` seit dem Start, gedruckt 60 s nach dem Run,
MiB); Bereich über die zwei Läufe (`n` = 2 je Stufe), daneben die Werte der Untersuchung:

| Changes in `cdc.change` nach dem Run | Nachmessung (Lauf 1; Lauf 2) | Bereich | Reihe J (Bereinigung aus) | Reihe B (vor der Änderung) |
|---|---|---|---|---|
| 1.000.000 (Run 1, 0 davor) | 15,1; 14,9 | 14,9 bis 15,1 | 13,3 | 1.082,7 |
| 2.000.000 (Run 2, 1.000.000 davor) | 16,9; 16,5 | 16,5 bis 16,9 | 13,2 | 2.269,2 |
| 3.000.000 (Run 3, 2.000.000 davor) | 17,6; 17,3 | 17,3 bis 17,6 | 13,7 | 3.096,6 |

Werte der Reihen J und B **übernommen** (Untersuchung, Abschnitt 3.3 bzw. Reihe J;
Reihe B Run 3: Spitze im Run, nicht im Nachlauf). Weitere gedruckte Größen der
Nachmessung (Lauf 1; Lauf 2):

| Größe | Run 1 | Run 2 | Run 3 |
|---|---|---|---|
| Kopierdauer in s | 91,6; 90,5 | 95,1; 96,0 | 102,0; 99,7 |
| `anon` Spitze im Run in MiB | 10,2; 10,0 | 12,2; 12,2 | 12,6; 12,8 |
| `memory.peak` unmittelbar nach dem Run in MiB | 13,6; 13,1 | 16,9; 16,5 | 17,6; 17,3 |
| GC-Läufe im Nachlauf | 560; 560 | 1.794; 1.655 | 2.454; 2.467 |
| Bereinigung gelaufen seit dem Start | 20; 20 | 22; 23 | 23; 23 |

Die Spitze im Run (`anon`) liegt bei leerem `cdc.change` (Run 1) bei 10,2 und 10,0 MiB, im
Bereich der Untersuchung bei leerem `cdc.change` (8,9 bis 10,5 MiB); bei 2.000.000 Changes
davor (Run 3) bei 12,6 und 12,8 MiB, also 2,4 und 2,8 MiB darüber. Der Bedarf einer Seite
(hergeleitet etwa 1,5 MiB,
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md))
ist aus dieser Messung nicht getrennt: sie enthält keinen Lauf mit einer anderen
Seitengröße, gegen den die Differenz zu bilden wäre.

**Zur Zusage der Nachmessung.** Die Erwartung war eine Form („hängt nicht an der Zahl der
Changes“, Werte im Bereich von Reihe J plus einer Seite), keine Schwelle. Erfüllt ist sie
für die Form: das Wachstum je Change liegt etwa um das Eintausendfache unter dem der
Untersuchung, und die Takte laufen weiter. Die Werte selbst liegen bei 1.000.000 Changes
im Rahmen „Reihe J plus etwa 1,5 MiB“ (Differenz 1,6 bis 1,8 MiB), bei 2.000.000 und
3.000.000 Changes darüber (3,3 bis 3,9 MiB); dieser Rest wird nicht als Beleg für „flach“
gezählt (Punkt 2 des Ergebnisses).

## 4. Skalierungs-Lauf (Regressions-Beleg der Live-Last)

`tools/bench-scaling.sh` im verkürzten Modus (Standard), ein Lauf, Feed-Image wie oben,
leerer `cdc.change`-Bestand: `cdc_capture_lag` je Lastenstufe **1,010192 s**, **1,003726 s**
und **1,003992 s** (Zeilen „Stufe klein / mittel / gross“ unten) gegen die bestehende
Grenze von 60 s ([`SPEC-013`](../../spec/pflichtenheft.md)). Der Lauf belegt, dass die
Änderung die Erfassung nicht in dieser Größenordnung stört; er belegt **nicht** das Verhalten
bei gefülltem `cdc.change` gleichzeitig mit Last (nicht gemessen). `make bench` als Ganzes
lief nicht (Auftrag: das Ganz-Target endet am Messhost bekannt rot bei
`tools/bench-source-impact.sh`).

## 5. Zusagen, Eingabeseiten-Mutationen und gesehenes Rot

Je Zusage eine benannte Änderung am geprüften Code (Eingabeseite), gefahren am Stand des
Commits `6d8e5d3f` (Use-Case-Mutationen 1 bis 8 und die der Übersetzung, ohne die spätere
Fortschrittsprüfung) bzw. `7572051c` (Fortschrittsprüfung, Lösch-Fehler, alle
Store-Tier-Mutationen) und danach zurückgenommen (`git checkout`); die
Zeilen „gesehenes Rot“ sind die gedruckten Zeilen des jeweiligen Laufs. Unit-Läufe: der
Docker-Aufruf von `make test` auf das Paket beschränkt (Race-Detector, netzlos); Store-Tier:
`make test-store` mit dem Aufruf auf die Tests `TestRetention…` beschränkt (Kopie des
Runners im Scratchpad, dieselben Images und derselbe Schema-Rollout).

**Use Case** (`internal/application/usecase/retention/service_test.go`):

| Zusage | Mutierte Eingabe | Gesehenes Rot |
|---|---|---|
| Jeder Lese-Aufruf trägt `limit` = `retention.PageSize` | Aufruf mit `PageSize+1` | `TestRunReadsEachPageAtPageSizeFromTheLastKey` |
| `after` ist die letzte Kennung der vorigen nichtleeren Seite | `after = page[0].ChangeID` | sechs Tests, u. a. `TestRunReadsEachPageAtPageSizeFromTheLastKey`, `TestRunReleasesSameSetAtEveryPageSize`, `TestRunDistinguishesEligibleChangesFromMixedSet` |
| `DeleteChanges` trägt nur die Freigaben dieser Seite | die ganze Seite statt der Freigaben | acht Tests, u. a. `TestRunDeletesOnlyTheReleasedIDsOfEachPage` |
| Nur die leere Seite beendet den Lauf | `len(page) < PageSize` beendet den Lauf | sieben Tests, u. a. `TestRunContinuesPastShortPages` |
| Positionen einmal je Lauf | Positionen je Seite neu gelesen | `TestRunReadsPositionsOncePerRun`, `TestRunDistinguishesEligibleChangesFromMixedSet` |
| `Deleted` zählt über alle Seiten | `deleted = len(eligible)` statt `+=` | `TestRunReleasesSameSetAtEveryPageSize`, `TestRunContinuesPastShortPages` |
| Die Freigabe trägt die Consumer-Position | `AllowsDeletion(…, consumerPositions[:0])` | sechs Tests, u. a. `TestRunReleasesSameSetAtEveryPageSize` |
| Ein Lese-Fehler ab Aufruf `n` hinterlässt die Löschungen davor | Lese-Fehler als Erfolg zurückgegeben | `TestRunStoreReadErrorFollowsSource`, `TestRunReadFailureLeavesPagesBefore` |
| Ein Lösch-Fehler an Seite `n` endet den Lauf | Lösch-Fehler als Erfolg zurückgegeben | `TestRunDeleteErrorFollowsEligibleSet`, `TestRunDeleteFailureStopsAtThatPage` |
| Eine Seite ohne Fortschritt endet als Fehler | Fortschrittsprüfung `last == after` abgeschaltet | kein Fehlschlag, sondern `panic: test timed out after 1m0s` in `TestRunRejectsPageWithoutProgress` (Rot durch Zeitüberschreitung) |

**Übersetzung** (`sqlexec/translate_test.go`): Commit-Zeitpunkt als `committedAt.Unix()`
statt `UnixNano()` → `TestReadRetentionCandidatesIssuesQueryAndTranslatesRows` rot. Die
Klassifikation der Treiber-Fehler und der Domänen-Fehler ist mit je einem Test gedeckt und
nicht einzeln mutiert.

**Store-Tier** (Abfrage `SelectRetentionCandidates`, Adapter, Rollen-Datei):

| Zusage | Mutierte Eingabe | Gesehenes Rot |
|---|---|---|
| Die Seite beginnt hinter `after` | `c.change_id >= $2` | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`: „der Durchlauf endet nicht nach 1000 Seiten (limit=1)“; `TestRetentionCandidateBehindCursorAppearsInNextRun`: „laufender Durchlauf = [c-20 c-20 c-20 …], wollen [c-20 c-30 c-40]“; `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames`: „Kandidaten-Seite ohne Fortschritt hinter "rp-tx-05000-5"“ |
| Nur Changes der Quelle | Quellfilter durch `$1::text IS NOT NULL` ersetzt (Parameter bleibt referenziert) | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`: „Kennung "tx4-b1" gehört nicht zur Quelle“ |
| Höchstens `limit` je Seite | `LIMIT GREATEST($3::int, 1000000)` | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`: „limit=1: Seite 1 trägt 5 Kandidaten“; `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames`: „Lese-Aufrufe = 2 ([10000 10000]), wollen 4“ |
| Aufsteigende Ordnung des Schlüssels | `ORDER BY c.change_id` entfernt | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`: „limit=1: 1 Kandidaten, wollen 5“; `TestRetentionCandidateBehindCursorAppearsInNextRun`: „laufender Durchlauf = [c-20 c-40], wollen [c-20 c-30 c-40]“; `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames`: „Deleted = 3084, unabhängige Zählung nennt 7710“ |
| Lesen unter `cdc_admin` gelingt, unter `cdc_capture`/`cdc_reader` nicht | `GRANT SELECT, DELETE ON cdc.transaction, cdc.change TO cdc_admin` auf `GRANT DELETE` gekürzt (`tools/schema/nacharbeit-roles.sql`) | `TestRetentionCandidatesRunUnderTheirRoles`: „ReadRetentionCandidates unter cdc_admin: … permission denied for table change (SQLSTATE 42501)“; `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames`: dieselbe Meldung |
| `limit` kleiner 1 endet als `ErrNonPositiveLimit` | Prüfung entfernt | `TestRetentionCandidatesRejectInvalidInput`: „limit=0: Fehler <nil>, wollen ErrNonPositiveLimit“ |
| Eine leere Quelle endet als `ErrEmptyIdentifier` | Prüfung entfernt (die `limit`-Prüfung bleibt) | `TestRetentionCandidatesRejectInvalidInput`: „leere Quelle: Fehler <nil>, wollen ErrEmptyIdentifier“ |
| Der Lauf löscht genau, was die Regel freigibt | `AllowsDeletion(…) \|\| true` im Use Case | `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames`: „Deleted = 25000, unabhängige Zählung nennt 7710“ |
| Zeitpunkt und Position je Kandidat gleich `ReadChanges` | nicht mutiert | kein gesehenes Rot: der Test vergleicht gegen `ReadChanges` und gegen die eingefügten Werte; die Zusage ist an der Eingabeseite (Spaltenwahl der Abfrage) nicht rot gesehen |

Zwei erste Versuche der Mutationen „Quellfilter entfernt“ und „`LIMIT` entfernt“ ließen den
Parameter der Abfrage unreferenziert und endeten mit `could not determine data type of
parameter $1 (SQLSTATE 42P18)`: rot aus einem technischen Grund, nicht wegen der Zusage;
die Zeilen oben stammen aus den wiederholten Läufen, in denen der Parameter referenziert
bleibt. Die Mutation `>=` ließ vor der Fortschrittsprüfung des Use Cases den Store-Tier-Lauf
bis zur Go-Zeitüberschreitung von 10 Minuten laufen (`TestRetentionRunOver…` „(9m59s)“);
seitdem endet er als Fehler.

## 6. Läufe des Slice und Befunde

**Läufe** (Zeilen unten): `make test` Exit 0 (42 Pakete `ok`); `make test-store` Exit 0
(`internal/bootstrap` ok, `postgresstorage` ok, `DB-Adapter-Coverage: 82.61% (gedeckt 879 von
1064 Statements; Profile gemergt: store,replication)`); `make test-replication` Exit 0 (dieselbe
Zeile); `make coverage-gate` Exit 0 (`Coverage 83.30% erfüllt Schwelle 80%`);
`make test-integration` Exit 0 (PostgreSQL 18; der Retention-Rundlauf „`RetentionOld` blieb
erhalten, solange ein Consumer zurückhing, und wurde nach Freigabe durch beide Consumer real
entfernt“ ohne geänderte Erwartung; das PostgreSQL-17-Leg läuft in CI). `make gates`: Abschnitt 8.

**Befunde.**

1. *Eine bestehende Testerwartung war zu ändern, ohne dass sich ihre Aussage ändert.*
   `TestRunStoreReadErrorFollowsSource` erwartete nach dem Lauf der Quelle `src-1`
   einen `DeleteChanges`-Aufruf und gab dem Fake keinen Bestand; ein Lauf über eine leere
   erste Seite endet ohne `DeleteChanges`. Die Erwartung („Lösch-Zug“) bleibt, der Fake trägt
   jetzt einen Change. Die übrigen Tests der Anforderungen `LH-FA-RET-002` bis
   `LH-FA-RET-004` liefen mit dem angepassten Fake ohne geänderte Erwartung.
2. *Fortschrittsprüfung im Use Case über den Plan hinaus* (Abschnitt 5, `>=`-Mutation):
   ohne sie endet ein Store, der den Cursor nicht beachtet, nicht.
3. *Neue Consumer während eines Laufs.* Die Frage des Plans (§6, Risiko „Ein Consumer
   bestätigt erstmals während eines längeren Laufs“) ist am Code beantwortet: `Positions` liest
   aus `cdc.consumer_position`; ein Consumer ohne Zeile dort erscheint nicht in der Rückgabe (`sqlexec.ReadConsumerPositions`) und blockiert die Freigabe nicht — das ist
   der bestehende Betriebs-Hinweis des Handbuchs. Die Vorwärts-Aussage von
   [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
   Festlegung 4 (Positionen wandern nur vorwärts) trägt für **bekannte** Consumer, nicht für
   einen neu bestätigenden. Das Fenster zwischen Lesen der Positionen und Löschen bestand mit
   dem Lauf vor der Änderung ebenso (Positionen, dann `ReadChanges`, dann `DeleteChanges` über
   die ganze Menge) und reicht jetzt bis zum Ende des Laufs; die Dauer des Laufs ist nach den
   Messungen des Architect-Verdikts (0,38 bis 2,0 s je 1.000.000 Changes für die Abfrage,
   übernommen) nicht länger als die des bisherigen Laufs. Das Mindestalter von 24 Stunden
   begrenzt die Kandidaten auf alte Changes. Der Slice ändert daran nichts (keine
   Vertragsänderung, Festlegung 4 wörtlich umgesetzt); eine Erneuerung der Positionen je
   Seite wäre eine Schärfung von Festlegung 4 (Positionen je Seite statt einmal) mit geänderter
   Erwartung in `TestRunReadsPositionsOncePerRun` und ist als Option benannt, nicht umgesetzt.
   Die Bewertung „akzeptiertes Negativ“ bleibt dem Architect (Plan: „vor dem Schließen
   vorlegen“).
4. *Bereinigungsdauer je Takt ungemessen.* Der Bericht enthält keine Dauer eines
   Bereinigungslaufs; die 0,38 bis 2,0 s je 1.000.000 Changes sind aus dem Architect-Verdikt
   übernommen (Einzelläufe, warmer Cache). Die Zeilen „Bereinigung gelaufen“ tragen keine
   Dauer.
5. *Der Zähler „Bereinigung gelaufen“ ist ein Gesamtzähler seit dem Start des Feeds*
   (20 bis 23): das Skript druckt ihn erst nach dem Nachlauf, nicht am Run-Ende. Das Wachstum
   im Nachlauf allein (hergeleitet fünf bis sechs je 60 s) ist damit nicht als Differenz
   gedruckt; als Beleg für laufende Takte tragen die Garbage-Collection-Läufe im Nachlauf
   (560 bis 2.467, in Reihe C vor der Änderung bei ausbleibenden Takten: keine) und der
   Gesamtzähler.

**Aufgeräumt.** `docker volume ls -q -f dangling=true | wc -l` vor den Läufen 34, nach
jedem Lauf 34 (kein Lauf hinterließ ein Volume); kein `docker volume prune`, kein
`docker system prune`; alle Container mit `docker rm -fv` (Bench-Skripte über ihren
Abräum-Pfad; ein durch Zeitüberschreitung abgebrochener Store-Tier-Lauf: `cdc-store-test-pg`
und das Netz `cdc-store-test` von Hand entfernt). `tools/schema/plan.yaml` und
`tools/schema/down.sql`, die jeder Rollout-Lauf umschreibt, sind nach jedem Lauf mit
`git checkout` zurückgenommen; nichts Erzeugtes ist committet.

## 7. Bewertung der Warn-Richtgröße

Trigger: „Eine Messung liegt vor“
([`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)).
**Ergebnis: der Wert von 4.000.000 geschätzten Zeilen
(`internal/application/usecase/backfill/warn.go`) bleibt unverändert.**

1. Das Speicher-Argument entfällt: der Speicher des Feed-Containers hängt nach der Änderung
   nicht mehr in dem Maß an der Zahl der Changes, die ein Backfill hinterlässt (Abschnitt 3);
   die Richtgröße bezog ihn ohnehin nicht ein.
2. Die Richtgröße folgt der Kopierdauer (Kopierrate mal Toleranz von 10 Minuten, auf eine
   Stelle abgerundet, Handbuch §9 „Grenzwerte“). Die Nachmessung druckt die Kopierdauer je
   Run über 1.000.000 Zeilen: 90,5 bis 102,0 s, das sind 9.802 bis 11.056 Zeilen/s (sechs Runs,
   Zeilen „Kopierdauer“ der Läufe 1 und 2); mal 600 s ergibt 5,9 bis 6,6 Millionen Zeilen
   (abgeleitet). Die Runs über 1.000.000 Zeilen der Untersuchung standen bei 9.522 bis 10.512
   (Reihe B) und 8.406 bis 8.756 Zeilen/s (Reihe C; gedruckt „Zeilen/s“, übernommen). Die
   Werte, aus denen der Code-Wert abgeleitet ist, stammen aus Stufen mit 200.000 Zeilen
   (4.504 bis 8.933 Zeilen/s, Handbuch §9).
3. Eine Anhebung wäre eine Änderung des Startwerts aus einer Messung, die nicht nach der
   Methode der Richtgröße läuft (Stufe mit 200.000 Zeilen, Median von drei Runs,
   `tools/bench-backfill.sh`); die Nachmessung dieses Slice ist eine Speicher-Messung, die
   Kopierdauer ihr Nebenwert. Der Wert bleibt eine Orientierung, keine Grenze
   (`AGENTS.md` §3.6: keine Schwellen-Änderung ohne ADR).

Das Handbuch nennt die Bewertung im Absatz zur Richtgröße („Die Richtgröße folgt der
Kopierdauer; der Speicher des Feed-Containers geht nicht ein“).

## 8. Gate

`make gates`: *(Zeile folgt nach dem Lauf)*

## 9. Gedruckte Zeilen

Die Ausgabe der Läufe, unverändert bis auf zwei Auslassungen: das Präfix
`bench-backfill-memory[<Lauf>]: ` bzw. `bench-scaling: ` am Zeilenanfang steht einmal in
der Überschrift, und die Hinweiszeile von `psql` zum Rollout („NOTICE: constraint … does not
exist, skipping“) fehlt.

### Lauf 1 — `20260925T183239Z`

`tools/bench-backfill-memory.sh`, Aufruf: `BENCH_MEM_STAGES=1000000 BENCH_FEED_ENV='GODEBUG=gctrace=1' BENCH_FEED_DOCKER_ARGS='--memory 6g'`, Ausgang 0.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:8c8dea000428448c0ff22a95f66061b32634beff80eeef1cf5c70e94c4a82a40
Stufen 1000000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '--memory 6g', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_1000000 — 1000000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 70 B (Textform)
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 5.9 MiB zu Beginn, Maximum 6.0, Minimum 5.9, Ende 5.9; anon Beginn 5.1, Maximum 5.1, Minimum 5.1, Ende 5.1 MiB
Stufe 1000000 narrow, Run 1/3 — Kopierdauer 91565 ms, 10921 Zeilen/s; Proben 352; memory.current Beginn 6.7, bei 25% 11.2, 50% 10.6, 75% 10.8, Ende 10.4, Spitze 12.3 MiB; anon Beginn 5.1, bei 25% 9.4, 50% 8.9, 75% 9.5, Ende 9.3, Spitze 10.2 MiB
Stufe 1000000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.4 MiB (anon 9.3, file 0.0, kernel 1.0), Prozess RssAnon 9.3 MiB, docker stats 10.3 MiB, VmHWM 25164 KiB, memory.peak seit Start des Feeds 13.6 MiB; im Run: GC-Läufe 1653, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.4 MiB zu Beginn, Maximum 13.7, Minimum 9.2, Ende 11.9; anon Beginn 9.3, Maximum 11.4, Minimum 8.1, Ende 10.9 MiB
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 1000000, VmHWM 26868 KiB, memory.peak seit Start des Feeds 15.1 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 560, größter Heap zu Beginn 7 MB, größtes Ziel 7 MB, letzter Heap nach GC 2 MB; seit dem Start des Feeds: Bereinigung gelaufen 20, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 155: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 1000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 11.2 MiB zu Beginn, Maximum 11.9, Minimum 10.7, Ende 10.7; anon Beginn 10.2, Maximum 10.2, Minimum 9.3, Ende 9.6 MiB
Stufe 1000000 narrow, Run 2/3 — Kopierdauer 95098 ms, 10515 Zeilen/s; Proben 365; memory.current Beginn 11.2, bei 25% 10.5, 50% 10.8, 75% 10.2, Ende 10.2, Spitze 14.7 MiB; anon Beginn 9.6, bei 25% 9.5, 50% 8.8, 75% 9.1, Ende 9.2, Spitze 12.2 MiB
Stufe 1000000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.2 MiB (anon 9.2, file 0.0, kernel 1.0), Prozess RssAnon 9.2 MiB, docker stats 12.0 MiB, VmHWM 27912 KiB, memory.peak seit Start des Feeds 16.9 MiB; im Run: GC-Läufe 2867, größter Heap zu Beginn 8 MB, größtes Ziel 8 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 11.8 MiB zu Beginn, Maximum 13.9, Minimum 11.2, Ende 11.9; anon Beginn 10.0, Maximum 11.4, Minimum 9.0, Ende 9.6 MiB
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 2000000, VmHWM 27912 KiB, memory.peak seit Start des Feeds 16.9 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 1794, größter Heap zu Beginn 7 MB, größtes Ziel 7 MB, letzter Heap nach GC 2 MB; seit dem Start des Feeds: Bereinigung gelaufen 22, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3202: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 2000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 11.1 MiB zu Beginn, Maximum 12.7, Minimum 10.7, Ende 10.7; anon Beginn 10.1, Maximum 10.1, Minimum 9.2, Ende 9.7 MiB
Stufe 1000000 narrow, Run 3/3 — Kopierdauer 102025 ms, 9802 Zeilen/s; Proben 391; memory.current Beginn 10.7, bei 25% 11.1, 50% 11.3, 75% 11.3, Ende 10.8, Spitze 15.4 MiB; anon Beginn 9.7, bei 25% 9.9, 50% 9.7, 75% 9.7, Ende 9.7, Spitze 12.6 MiB
Stufe 1000000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.8 MiB (anon 9.7, file 0.0, kernel 1.1), Prozess RssAnon 9.7 MiB, docker stats 10.8 MiB, VmHWM 28552 KiB, memory.peak seit Start des Feeds 17.6 MiB; im Run: GC-Läufe 4283, größter Heap zu Beginn 8 MB, größtes Ziel 9 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.8 MiB zu Beginn, Maximum 15.4, Minimum 10.8, Ende 11.6; anon Beginn 9.7, Maximum 12.0, Minimum 9.2, Ende 10.5 MiB
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 3000000, VmHWM 28552 KiB, memory.peak seit Start des Feeds 17.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 2454, größter Heap zu Beginn 7 MB, größtes Ziel 7 MB, letzter Heap nach GC 2 MB; seit dem Start des Feeds: Bereinigung gelaufen 23, fehlgeschlagen 0
Ende — Lauf 20260925T183239Z
```

### Lauf 2 — `20260925T184503Z`

Derselbe Aufruf, Ausgang 0.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:8c8dea000428448c0ff22a95f66061b32634beff80eeef1cf5c70e94c4a82a40
Stufen 1000000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '--memory 6g', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_1000000 — 1000000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 70 B (Textform)
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.3 MiB zu Beginn, Maximum 6.3, Minimum 6.3, Ende 6.3; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 1000000 narrow, Run 1/3 — Kopierdauer 90450 ms, 11056 Zeilen/s; Proben 347; memory.current Beginn 6.8, bei 25% 11.4, 50% 10.1, 75% 11.2, Ende 10.9, Spitze 12.5 MiB; anon Beginn 5.4, bei 25% 9.0, 50% 9.1, 75% 9.4, Ende 9.5, Spitze 10.0 MiB
Stufe 1000000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.5 MiB (anon 9.5, file 0.0, kernel 1.0), Prozess RssAnon 9.5 MiB, docker stats 10.5 MiB, VmHWM 24900 KiB, memory.peak seit Start des Feeds 13.1 MiB; im Run: GC-Läufe 1658, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.5 MiB zu Beginn, Maximum 13.7, Minimum 9.7, Ende 11.9; anon Beginn 9.5, Maximum 11.5, Minimum 8.1, Ende 10.8 MiB
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 1000000, VmHWM 26416 KiB, memory.peak seit Start des Feeds 14.9 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 560, größter Heap zu Beginn 7 MB, größtes Ziel 7 MB, letzter Heap nach GC 2 MB; seit dem Start des Feeds: Bereinigung gelaufen 20, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 149: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 1000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 10.7 MiB zu Beginn, Maximum 12.7, Minimum 10.7, Ende 11.5; anon Beginn 9.7, Maximum 10.5, Minimum 9.4, Ende 10.5 MiB
Stufe 1000000 narrow, Run 2/3 — Kopierdauer 96002 ms, 10416 Zeilen/s; Proben 369; memory.current Beginn 11.7, bei 25% 9.9, 50% 11.5, 75% 11.2, Ende 10.6, Spitze 14.6 MiB; anon Beginn 10.5, bei 25% 8.8, 50% 9.9, 75% 9.6, Ende 9.0, Spitze 12.2 MiB
Stufe 1000000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 12.8 MiB (anon 10.9, file 0.0, kernel 1.1), Prozess RssAnon 10.9 MiB, docker stats 12.5 MiB, VmHWM 27896 KiB, memory.peak seit Start des Feeds 16.5 MiB; im Run: GC-Läufe 2985, größter Heap zu Beginn 8 MB, größtes Ziel 8 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 12.5 MiB zu Beginn, Maximum 14.2, Minimum 10.9, Ende 13.8; anon Beginn 11.4, Maximum 11.7, Minimum 9.1, Ende 11.1 MiB
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 2000000, VmHWM 28184 KiB, memory.peak seit Start des Feeds 16.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 1655, größter Heap zu Beginn 7 MB, größtes Ziel 7 MB, letzter Heap nach GC 2 MB; seit dem Start des Feeds: Bereinigung gelaufen 23, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3162: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 2000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 10.3 MiB zu Beginn, Maximum 13.0, Minimum 10.3, Ende 11.4; anon Beginn 9.3, Maximum 10.7, Minimum 9.3, Ende 10.4 MiB
Stufe 1000000 narrow, Run 3/3 — Kopierdauer 99706 ms, 10029 Zeilen/s; Proben 383; memory.current Beginn 11.7, bei 25% 11.8, 50% 11.1, 75% 11.8, Ende 11.1, Spitze 15.5 MiB; anon Beginn 10.4, bei 25% 9.3, 50% 9.8, 75% 9.5, Ende 9.5, Spitze 12.8 MiB
Stufe 1000000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.6 MiB (anon 9.5, file 0.0, kernel 1.1), Prozess RssAnon 9.5 MiB, docker stats 10.6 MiB, VmHWM 28792 KiB, memory.peak seit Start des Feeds 17.3 MiB; im Run: GC-Läufe 4289, größter Heap zu Beginn 8 MB, größtes Ziel 9 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.6 MiB zu Beginn, Maximum 14.8, Minimum 10.6, Ende 11.4; anon Beginn 9.5, Maximum 11.6, Minimum 9.5, Ende 10.4 MiB
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 3000000, VmHWM 28792 KiB, memory.peak seit Start des Feeds 17.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 2467, größter Heap zu Beginn 7 MB, größtes Ziel 7 MB, letzter Heap nach GC 3 MB; seit dem Start des Feeds: Bereinigung gelaufen 23, fehlgeschlagen 0
Ende — Lauf 20260925T184503Z
```

### Skalierungs-Lauf

`tools/bench-scaling.sh` (verkürzter Modus), Ausgang 0.

```text
Umgebung wird aufgebaut (Modus: verkürzt)…
CREATE TABLE
INSERT 0 1
Stufe klein — Ziel 10/s über 10s, 100 Zeilen in 10s eingefügt (~10.0/s), cdc_capture_lag=1.010192s
Stufe mittel — Ziel 100/s über 15s, 1500 Zeilen in 15s eingefügt (~100.0/s), cdc_capture_lag=1.003726s
Stufe gross — Ziel 1000/s über 15s, 15000 Zeilen in 15s eingefügt (~1000.0/s), cdc_capture_lag=1.003992s
Ergebnis (LH-QA-PER-002) — alle drei SPEC-014-Lastenstufen real durchlaufen (Modus: verkürzt)
```

### Tests und Gates (Auszug der gedruckten Zeilen)

```text
make test            Exit 0 (42 Zeilen ok, keine FAIL)
make test-store      Exit 0
  ok  	github.com/pt9912/pg-change-feed/internal/bootstrap	3.336s
  ok  	github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage	7.409s	coverage: 78.9% of statements in ./internal/adapters/driven/postgresstorage, ./internal/adapters/driven/postgresack, ./internal/adapters/driven/postgressnapshot, ./internal/adapters/driving/replication/receive
  DB-Adapter-Coverage: 82.61% (gedeckt 879 von 1064 Statements; Profile gemergt: store,replication)
  db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%
make test-replication Exit 0
  DB-Adapter-Coverage: 82.61% (gedeckt 879 von 1064 Statements; Profile gemergt: store,replication)
  db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%
make coverage-gate   Exit 0
  total:	(statements)	83.3%
  coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%
make test-integration Exit 0
  run-integration-tests: Retention-Beleg — 'RetentionOld' (id=200) blieb erhalten, solange ein Consumer zurückhing (LH-FA-RET-004), und wurde nach Freigabe durch beide Consumer real entfernt (LH-FA-RET-003); 'RetentionYoung' (id=201, zu jung) blieb durchgehend erhalten
  run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 14 Go-Zeilen und 35 Bash-Zeilen
```
