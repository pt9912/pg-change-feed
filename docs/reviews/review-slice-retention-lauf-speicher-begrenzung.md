# Review-Report: slice-retention-lauf-speicher-begrenzung — 2026-09-25

**Review-Art:** Code — der Diff ändert Produktionscode im Retention-Pfad (Löschung von Changes):
`ChangeStorePort.ReadRetentionCandidates`, den Use Case `RunRetentionService`, Adapter, Abfrage und
Übersetzung, vier Test-Fakes, Store-Tier-Tests, den Nachzug der Träger (Handbuch 1.62, Pflichtenheft,
Sensor-Dokumente, `harness/README.md`) und einen Messbericht; geprüft gegen Plan, ADRs und
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Hauptprüfgegenstand ist die Sicherheit der
Löschung („es wird nie früher oder mehr gelöscht als vorher“). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-retention-lauf-speicher-begrenzung` (wellenlos), Diff-Range
`8cd39719..b2484add` (HEAD `b2484add`, Arbeitsbaum vor dem Report sauber). Slice-Inhalt sind
`6d8e5d3f` (Umsetzung), `d49e248f` und `7572051c` (Fortschrittsprüfung), `0e6b1b30` (Plan-Nachzug),
`2636ca46` (Nachmessung, Messbericht, Handbuch, Pflichtenheft), `a583afe5` (Gate-Zeilen im
Messbericht), `b2484add` (DoD-Haken) sowie die drei Plan-Commits `e486beb8`, `8a68d2b2`, `ff387a69`.
**CI-Stand am Commit `b2484add`** (`gh run list --commit b2484add6cabd42781794ad68c3727f4fa85978e`,
abgerufen nach Abschluss der Läufe): `ci` success (Lauf 36178837163, 2m41s), `examples` success
(Lauf 36178837307, 1m55s), `e2e` success (Lauf 36178837299, 12m47s).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Bindung, Kommentar-Klassen,
Handbuch-Zug). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-retention-lauf-speicher-begrenzung` (§1 Ziel und Abgrenzung, §2 DoD, §3 Plan samt
  Suchlauf-Feld, §4 Trigger, §6 Risiken, §7 Release-Aussage), Architect-Verdikt
  [`architect-verdict-retention-lauf-speicher-begrenzung`](architect-verdict-retention-lauf-speicher-begrenzung.md)
- [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) (`Accepted`),
  [`ADR-0014`](../plan/adr/0014-retention-domain-policy.md),
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
  [`ADR-0053`](../plan/adr/0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md)
- [`LH-FA-RET-002`](../../spec/lastenheft.md), [`LH-FA-RET-003`](../../spec/lastenheft.md),
  [`LH-FA-RET-004`](../../spec/lastenheft.md), [`LH-FA-CAP-009`](../../spec/lastenheft.md);
  [`LH-FA-RET-004.a`](../../spec/pflichtenheft.md) im Pflichtenheft
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`)
- Vorbericht am selben Modul: `review-slice-backfill-speicher-untersuchung`; Report-Gerüst
  `docs/reviews/review-report.template.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; ein schwerer Docker-Lauf
zugleich; kein `prune`; dangling Volumes `docker volume ls -q -f dangling=true | wc -l` **34 vor und
34 nach** allen Läufen und Mutationen; keine eigenen Container oder Netze zurückgeblieben;
`tools/schema/plan.yaml` und `tools/schema/down.sql` nach jedem Rollout-Lauf per `git checkout`
zurückgenommen):

- **Sicherheit der Löschung am Code nachgelesen** (Stand `b2484add` gegen `git show 8cd39719:…`):
  - *Freigabe je Kandidat:* `command.Policy.AllowsDeletion(age, candidate.Position, consumerPositions)`
    mit `age = now.Sub(candidate.CommittedAt)` — dieselben drei Eingaben wie am Parent
    (`record.CommittedAt`, `record.Position`, Consumer-Positionen). Die Wanduhr wird jetzt vor der
    ersten Seite gelesen (am Parent nach dem Lesen aller Changes): das Alter wird höchstens
    unterschätzt, eine Löschung kommt nie früher.
  - *Kandidatenmenge:* `SelectRetentionCandidates` filtert `t.source_id = $1`, kein Prädikat außer dem
    Cursor; der Parent-Weg `SelectChanges` verband zusätzlich `cdc.source_table` per `JOIN`. Die
    Fremdschlüssel-Spalte `change.source_table_id` ist `required: true` mit `references`
    (`tools/schema/schema.yaml:153–156`), der `JOIN` filterte deshalb keine Zeile: die Kandidatenmenge
    ist gleich der bisherigen Bereinigungsmenge-Grundlage.
  - *Cursor:* Schlüssel ist der Primärschlüssel `change_id` (`type: text`, `schema.yaml:148`,
    `primary_key: [change_id]:172`); `c.change_id > $2` und `ORDER BY c.change_id` laufen in **einer**
    Anweisung unter derselben Sortierung, die Seitenfolge ist ein Keyset über einen eindeutigen Schlüssel
    ohne `OFFSET`; das Löschen der bereits gelesenen Seite verschiebt den Cursor nicht. Eine Kennung
    kann nur „hinter“ den Cursor sortieren, wenn sie neu committet wird (dann erscheint sie im nächsten
    Lauf, ist aber im laufenden Lauf ohnehin jünger als das Mindestalter oder nicht gesehen); ein
    Umsortieren einer schon gesehenen Kennung ist bei einer deterministischen Sortierung der Datenbank
    nicht möglich (das Schema trägt keine Spalten-`COLLATE`-Klausel, `grep -i collat tools/schema/*`
    leer). Die Ordnung `0bf-…` vor WAL-Kennungen ist für die Korrektheit des Keyset nicht erforderlich.
  - *Löschung:* `DeleteChanges` übergibt die Kennungen als **ein** Array-Parameter (`change_id =
    ANY($1)` mit `[]string`, `store.go:200–214`), nicht als Einzelparameter: das
    Parameter-Limit von 65.535 greift bei 10.000 Kennungen je Seite nicht; Adapter und Abfrage sind im
    Diff unverändert, die Waisen-Bereinigung (`DeleteOrphanedTransactions`) prüft `NOT EXISTS` je
    berührter Transaktion und lässt eine über zwei Seiten liegende Transaktion bis zur Löschung ihrer
    letzten Change-Zeile stehen.
  - *Fehlerpfade:* jeder Lese- und Lösch-Fehler kehrt als Fehler zurück (`return RunRetentionResult{},
    err`), kein stilles Weitermachen; ein Kontext-Abbruch trifft den nächsten Lese-Aufruf der Seite und
    endet als Fehler der Klasse `storage`; `runRetentionCleanup` loggt den Fehler und tickt weiter.
    Die Fortschrittsprüfung `last == after` steht seit `7572051c` **vor** der Löschung der Seite.
  - *Rollen:* `tools/schema/nacharbeit-roles.sql` ist im Diff unverändert; das Lesen der Kandidaten
    braucht kein neues Recht (Store-Tier-Test unter echten Logins, unten).
- **Gates am Stand `b2484add`:** `make test` Exit 0 (alle Pakete `ok`, darunter
  `usecase/retention` und `internal/bootstrap`); `make test-store` (PostgreSQL 18) Exit 0, gedruckt
  `DB-Adapter-Coverage: 82.61% (gedeckt 879 von 1064 Statements; Profile gemergt: store,replication)`
  und `db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%`; `make coverage-gate`
  Exit 0, gedruckt `total: (statements) 83.3%` und `coverage-gate: OK — Coverage 83.30% erfüllt
  Schwelle 80%`; das Profil `/out/coverage.out` des gebauten Images per Awk über die Block-Position
  dedupliziert ausgezählt: **gedeckt 2167 von 2601 = 83,31 %** (gleich der Zeile in
  `harness/sensors/coverage-gate.md`); `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make docs-check`
  Exit 0, `d-check: 1175 Datei(en) geprüft, 0 Befund(e)`.
- **Mutationen der Eingabeseite selbst gefahren** (je mit `git checkout` zurückgenommen):
  - *Use Case* (`make test`-Aufruf auf das Paket beschränkt, Race-Detektor, netzlos): **M1**
    Fortschrittsprüfung abgeschaltet → kein Test-Fehlschlag, sondern `panic: test timed out after 40s`
    (`-timeout 40s`); **M2** `AllowsDeletion(…)` durch `true || AllowsDeletion(…)` → neun Tests rot
    (u. a. `TestRunDistinguishesEligibleChangesFromMixedSet`, `TestRunReleasesSameSetAtEveryPageSize`);
    **M3** Positionen je Seite neu gelesen → `TestRunReadsPositionsOncePerRun`,
    `TestRunDistinguishesEligibleChangesFromMixedSet` rot; **M4** `limit = PageSize+1` →
    `TestRunReadsEachPageAtPageSizeFromTheLastKey` rot; **M5** Cursor auf `""` gesetzt → Zeitüberschreitung
    (wie M1); **M6** die ganze Seite statt der Freigaben an `DeleteChanges` → neun Tests rot; **M7**
    Alter mit vertauschtem Vorzeichen → acht Tests rot.
  - *Store-Tier* (Kopie von `tools/harness/run-store-tests.sh` im Scratchpad, dieselben Images und
    derselbe d-migrate-Rollout, Aufruf auf `TestRetention…` beschränkt): **S1** `c.change_id >= $2` →
    `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames`: „Kandidaten-Seite ohne
    Fortschritt hinter "rp-tx-05000-5"“, `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`: „der
    Durchlauf endet nicht nach 1000 Seiten (limit=1)“, `TestRetentionCandidateBehindCursorAppearsInNextRun`:
    „laufender Durchlauf = [c-20 c-20 …]“; **S2** Quellfilter durch `$1::text IS NOT NULL` ersetzt →
    `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`: „Kennung "tx4-b1" gehört nicht zur Quelle“
    (der 25.000-Changes-Test bleibt dabei grün: in seiner Datenbank steht nur seine Quelle); **S3**
    `LIMIT GREATEST($3::int, 1000000)` → „Lese-Aufrufe = 2 ([10000 10000]), wollen 4“ und „limit=1:
    Seite 1 trägt 5 Kandidaten“; **S4** `ORDER BY` entfernt → „Deleted = 3084, unabhängige Zählung
    nennt 7710“, „limit=1: 1 Kandidaten, wollen 5“; **S5** `t.commit_position + 1` →
    `TestRetentionCandidateCarriesPositionAndTimeOfReadChanges` rot (Offset 1001 statt 1000; der
    25.000-Changes-Test bleibt grün, siehe F-7); **S6** `t.committed_at + interval '1 hour'` → beide
    Tests rot („Deleted = 7500, unabhängige Zählung nennt 7710“). Alle sechs Store-Tier-Läufe
    reproduzieren die Zeilen des Messberichts (Abschnitt 5), soweit dort gedruckt (S1, S3, S4 wörtlich).
- **Nachmessung selbst gefahren** (ein Lauf, `tools/bench-backfill-memory.sh` nach Vertrag von
  `harness/targets/bench-backfill.md`, `BENCH_MEM_STAGES=1000000`, `BENCH_FEED_ENV='GODEBUG=gctrace=1'`,
  `BENCH_FEED_DOCKER_ARGS='--memory 6g'`, Feed-Image
  `sha256:8c8dea000428448c0ff22a95f66061b32634beff80eeef1cf5c70e94c4a82a40` — dieselbe ID wie im
  Messbericht; Lauf `20260925T193207Z`, Exit 0; `free -m` zu Beginn des Reviews 19.482 MB verfügbar;
  Host Linux 6.8.0-139-generic, 20 CPU, 33.362.599.936 Byte RAM). Gedruckt, Run 1/2/3 (0, 1.000.000,
  2.000.000 Changes vor dem Run): `memory.peak seit Start des Feeds` nach dem Nachlauf **31,8 / 16,6 /
  16,8 MiB**; `anon` Spitze im Run 10,3 / 12,7 / 12,5 MiB; Kopierdauer 96.260 / 103.730 / 110.181 ms
  (10.389 / 9.640 / 9.076 Zeilen/s); „Bereinigung gelaufen“ seit dem Start 22 / 23 / 24, fehlgeschlagen
  0; Zeilen in `cdc.change` nach dem Nachlauf 1.000.000 / 2.000.000 / 3.000.000; `OOMKilled false`,
  Neustarts 0. Der Wert 31,8 MiB in Run 1 trägt `file` 16,5 MiB (`cgroup Run-Ende: memory.current
  26,6 MiB (anon 9,1, file 16,5, kernel 1,0)`), siehe F-9.
- **Nachrechnungen gegen die gedruckten Zeilen des Messberichts** (Zeilen extrahiert, dann gerechnet;
  alle bestätigt): Differenzen zu Reihe J 15,1−13,3 = 1,8 und 14,9−13,3 = 1,6 / 16,9−13,2 = 3,7 und
  16,5−13,2 = 3,3 / 17,6−13,7 = 3,9 und 17,3−13,7 = 3,6 MiB; Wachstum 17,3−14,9 = 2,4 und
  17,6−15,1 = 2,5 MiB; 2,5 MiB × 1.048.576 / 2.000.000 = 1,31 Byte je Change; Kopierraten
  1.000.000 / 91,565 s = 10.921, 1.000.000 / 102,025 s = 9.802, 1.000.000 / 90,450 s = 11.056 Zeilen/s,
  × 600 s = 5,88 bis 6,63 Millionen; „Anon-Spitze im Run“ Run 3 minus Run 1: 12,6−10,2 = 2,4 und
  12,8−10,0 = 2,8 MiB; Zähler „Bereinigung gelaufen“ 20/22/23 und 20/23/23; die Streuung der Reihen
  B/C/J des Messberichts der Untersuchung (1.082,7; 1.273,5; 2.269,2; 3.058,0; 3.096,6; 2.751,7;
  13,3; 13,2; 13,7) an dessen Zeilen 159–161 und 303 bestätigt.
- **Suchlauf des Plans (§3) an beiden Ständen nachgefahren** (`git grep -n … 8cd39719 -- <Wurzeln>`
  und `git grep -n … -- <Wurzeln>` im Arbeitsbaum, dazu `5b1f7762`): Befehl 1 Parent 298 (`internal`
  263, `test` 24, `tools` 2, `spec` 7, `docs/user` 2, `harness` 0), Diff 309 (274/24/2/7/2/0);
  Befehl 2 Parent 13, Diff 6 (alle `docs/user`); Befehl 3 Parent 87 (25/5/5/16/9/0/27), Diff 85
  (18/7/7/17/9/0/27); `5b1f7762`: 298 / 12 / 80. **Jede Zahl des Feldes stimmt.** Die Ergänzung „Ein
  Lokator, der mitgewandert wäre“ nachgeprüft: `coverage-gate.md` nennt `:1287.4,1288.1` und
  `:1183.5,1184.13`; `wiring.go` Zeile 1183–1184 und 1287–1288 tragen an `8cd39719` und am Diff-Stand
  denselben Inhalt (`git diff --stat 8cd39719..HEAD -- internal/bootstrap/wiring.go`: 4 Einfügungen,
  4 Löschungen, Zeilenzahl gleich). Zusätzlich gesucht: `git grep -n -i -E
  'mem_limit|memory:|--memory|OOMKilled|Speicherbedarf|Arbeitsspeicher'` über das Repository ohne
  Records — keine weiteren Träger der bewegten Eigenschaft.
- **Release-Aussage nachgeprüft:** `git show v0.1.0:…/service.go`, `v0.1.1:…`, `v0.1.2:…` tragen in
  Zeile 63 `ReadChanges(ctx, outbound.ChangeQuery{Source: command.Source})`, `wiring.go` Zeile 172
  `const retentionInterval = 10 * time.Second`; `git ls-tree -r --name-only <Tag> | grep -c
  usecase/backfill` ist 0 in allen drei Tags; `git diff v0.1.2 8cd39719 -- internal/application/usecase/retention
  internal/application/port/outbound/changestore.go` ist leer; `git tag -l 'v*'` nennt nur diese drei
  Tags. Herleitung 10 Changes/s × 86.400 s = 864.000 Changes; × 1,03 bis 1,59 KiB = 0,85 bis 1,31 GiB
  (gerechnet, als hergeleitet gekennzeichnet).
- **Roter Zwischencommit gemessen:** `git archive d49e248f` in ein Scratchpad-Verzeichnis, `go test
  -race ./internal/application/usecase/retention/` im gepinnten Toolchain-Image: `--- FAIL:
  TestRunRejectsPageWithoutProgress … gelöschte Menge = [c1 c1], wollen genau die erste Seite [c1]`
  (F-5).
- **Nicht gefahren (Grenze):** `make bench` als Ganzes (Auftrag), `make test-integration` und
  `make test-replication` (die CI-Läufe `ci`/`e2e` an `b2484add` sind grün, der Retention-Rundlauf
  des E2E-Laufs liegt darin), Nachmessungen mit anderer Seitengröße, mit breiten Zeilen und mit
  gleichzeitiger Live-Last bei gefülltem `cdc.change`.
- **Nebenläufe anderer Rollen:** keine; der Report-Commit nimmt ausschließlich diese Datei auf.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Träger sagen ohne Einschränkung, der Speicher „hänge nicht an der Zahl der Changes“, und tragen im selben Träger eine Messung, in der die Spitze von 1.000.000 auf 3.000.000 Changes um 2,4 und 2,5 MiB wächst (Bereiche 14,9–15,1 < 16,5–16,9 < 17,3–17,6, nicht überlappend) samt „Ursache … nicht untersucht“. Im Handbuch steht der Satz „hängt nicht an der Zahl der Changes“ (Zeile 1676) rund zwanzig Zeilen vor dem Wachstumssatz (Zeile 1696–1699); dieselbe unbedingte Aussage tragen §4 (Zeile 476–477, 665–666), Zeile 1648, das Pflichtenheft (Punkt 4) und der Godoc von `PageSize`. | [`LH-FA-RET-004.a`](../../spec/pflichtenheft.md) · Reviewer-Skill „Nachzug widerspricht dem Nachbarn im selben Träger“ · `AGENTS.md` §3.12 | `docs/user/benutzerhandbuch.md:1676` gegen `:1696`; `spec/pflichtenheft.md:141`; `internal/application/usecase/retention/service.go:46` | ja — Lesen der Zeilen; die Zahlen an Messbericht Abschnitt 3 und am Review-Lauf `20260925T193207Z` | Nachzug widerspricht dem Nachbarn |
| F-2 | MEDIUM | Die Plan-Klausel „wächst die Spitze mit der Zahl der Changes … geht der Slice nicht nach `done/`“ (DoD Punkt „Nachmessung“, §2 Zeile 200–201) hat im Messbericht keine ausgesprochene Auswertung: der Wortlaut ist an den zwei Läufen des Implementers erfüllt (Bereiche steigen ohne Überlappung), der Bericht nennt den Rest „unerklärt“ und die DoD-Zeile steht auf `[x]`; die Zeile nennt sich zugleich „Orientierung, keine Schwelle“. Die gedruckten Zeilen tragen eine Deutung, die der Bericht nicht zieht: siehe die Empfehlung unter „Bewertung der Plan-Klausel“. Die abgeleitete Zahl „1,3 Bytes je Change“ gilt zwei Endpunktpaaren, deren Anstieg an der Stufe 1.000.000 → 2.000.000 hängt. | Slice-Plan §2/§6 · `AGENTS.md` §3.12 Instanz A und B | `docs/plan/planning/in-progress/slice-retention-lauf-speicher-begrenzung.md:186–201`; `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md:30–39` | nein — Entscheidung, keine Gate-Frage; die Messung: Nachmessung des Reviewers und Zeilen des Messberichts | Plan-Klausel ohne dokumentierte Auswertung |
| F-3 | MEDIUM | Die Frage „neuer Consumer während eines Laufs“ ist vom Implementer am Code beantwortet, die vom Plan verlangte Bewertung des Architect liegt nicht vor (Plan §6: „der Implementer legt die Frage dem Architect vor, bevor er sie als akzeptiert schließt“; Messbericht Abschnitt 6, Befund 3: „bleibt dem Architect“). Das Fenster „Positionen gelesen — Änderung geprüft — gelöscht“ reicht mit den Seiten bis zum Ende des Laufs; die Aussage „nicht breiter geworden“ stützt sich auf übernommene Datenbank-Zeiten der Abfrage (`n` = 1, Verdikt), die Dauer eines Laufs am Feed ist nach Befund 4 ungemessen, und der Nachlauf-Beleg löscht nichts (alle Changes der Nachmessung sind jünger als das Mindestalter). Der Handbuch-Satz „gilt seine Position erst im nächsten Durchlauf“ (Zeile 679–680) nennt die Folge (Löschung in diesem Lauf) nicht. | [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) Festlegung 4 · [`LH-FA-RET-004`](../../spec/lastenheft.md) · Slice-Plan §6 | `docs/plan/planning/in-progress/slice-retention-lauf-speicher-begrenzung.md:442–462`; `internal/application/usecase/retention/service.go:70–78`; `docs/user/benutzerhandbuch.md:679` | nein — die Bewertung ist eine Architect-Entscheidung; das Verhalten ist am Code nachlesbar, die Laufdauer ungemessen | Offene Bewertung eines Sicherheitsfensters |
| F-4 | LOW | Die Fortschrittsprüfung `last == after` (Endlosschleifen-Schutz) ist an ihrer Eingabeseite nur durch eine Zeitüberschreitung rot zu färben: Mutation M1 (Prüfung abgeschaltet) und M5 (Cursor auf `""`) enden in `panic: test timed out`, nicht in einer Assertion; der Fake hängt bei jedem Umlauf an `removed` und `deleteCalls` an. Der Test bindet die Zusage, sein Rot ist aber ein Hänger von der Länge des Go-Timeouts (Vorgabe von `make test`: 10 Minuten). | Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“ (Randfall: gebunden, Rot nur durch Zeitablauf) | `internal/application/usecase/retention/service_test.go:619` | ja — Mutation `if false && last == after` und `go test -timeout 40s` | Rot nur durch Zeitüberschreitung |
| F-5 | LOW | Der Commit `d49e248f` (auf `main`, gepusht) ist rot: `TestRunRejectsPageWithoutProgress` schlägt an ihm fehl (`gelöschte Menge = [c1 c1], wollen genau die erste Seite [c1]`), der Folgecommit `7572051c` behebt es. Ein `git bisect` über diesen Bereich hält an `d49e248f` an. | Maintainability | Commit `d49e248f`, `internal/application/usecase/retention/service.go` | ja — `git archive d49e248f` und `go test` auf das Paket | Roter Zwischencommit auf `main` |
| F-6 | LOW | Das Handbuch beschreibt den Stand der Vorversionen mit Vorher-Nachher-Sprache: „am Stand vor der seitenweisen Lesung gemessen 1,03 bis 1,59 KiB je Change“ und der Absatz „Vorversionen“ nennen einen früheren Stand des Codes statt den Ist-Zustand einer benannten Version; die Aussage über `v0.1.0` bis `v0.1.2` ist an sich belegt (Prüfung oben). | `AGENTS.md` §3.7 · Memory-Regel „keine Chronik/Forensik in Doku-Prosa“ · Reviewer-Skill „Zustandsfeld/Träger trägt Chronik“ | `docs/user/benutzerhandbuch.md:1710–1717` | nein — Lesen | Vorher-Nachher-Sprache in Handbuch-Prosa |
| F-7 | INFO | Der 25.000-Changes-Test bindet die Grenze „Consumer-Position gleich Change-Position“ nicht: die Transaktion an Position 3.500 (`g` = 3.500, Commit-Zeitpunkt seedAt + 4.500 Minuten) ist jünger als die 24-h-Grenze der festen Uhr, die Mutation S5 (`commit_position + 1`) bleibt dort grün. Die Grenze tragen der Einheitentest (`c4` an Position 200, Consumer 200) und der Adapter-Test gegen `ReadChanges` (S5 rot). | Maintainability | `internal/bootstrap/retention_pages_test.go:55–58, 83–88` | ja — Mutation S5 | Grenzwert ohne Bindung auf dieser Test-Ebene |
| F-8 | INFO | Zwei Kommentare beschreiben den Fehlerpfad der Fortschrittsprüfung mit dem ausgeschlossenen Ausgang: „endet als Fehler der Klasse `storage`, nicht als Endlosschleife“ (Godoc von `Run`) und „statt dieselbe Seite endlos zu lesen“ (Test-Godoc). Subjekt ist der Zustand des Codes, kein Konjunktiv über eine verworfene Alternative; die Formulierung liegt an der Grenze der Kommentar-Klasse „Abgrenzung“. | `AGENTS.md` §3.7 · Reviewer-Skill „Kommentar trägt keine der Kommentar-Klassen“ | `internal/application/usecase/retention/service.go:61`; `service_test.go:618` | nein — Lesen | Kommentar an der Klassengrenze |
| F-9 | INFO | `memory.peak` (Kennzahl von Bericht, Handbuch und Plan-Erwartung) schließt den Seiten-Cache der cgroup ein: im Review-Lauf trägt Run 1 `file` 16,5 MiB bei `anon` 9,1 MiB und `memory.peak` 31,8 MiB, in den zwei Run-1-Fenstern des Implementers steht dort 14,9 und 15,1 MiB (`memory.current` zu Beginn 5,9 und 6,3 MiB gegen 22,3 MiB im Review-Lauf; die Ursache des Cache-Anteils ist nicht untersucht). Die Kennzahl streut damit um Größen, die mit dem Bedarf des Go-Prozesses nichts zu tun haben; der Messbericht nennt die Spalte `anon` als Prozessgröße, wertet aber `memory.peak`. | Maintainability · `harness/targets/bench-backfill.md` §Grenzen dieses Skripts | Review-Lauf `20260925T193207Z` Run 1; `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md:83–87` | ja — der genannte Lauf | Kennzahl streut mit Cache-Anteil |
| F-10 | INFO | Die gedruckten Zeilen des Messberichts trennen den Seiten-Anteil, den der Bericht als „nicht getrennt“ führt: die Grundlinie (`anon`, Ruhe 10 s nach dem Start) liegt bei leerem `cdc.change` bei 5,1 und 5,4 MiB (Run 1, beide Läufe) und bei 1.000.000 und 2.000.000 Changes davor bei 10,1 bis 10,7 MiB Maximum (Run 2 und 3, beide Läufe und Review-Lauf: 10,2 / 10,1 / 10,5 / 10,7 / 10,7 / 10,6, Maximum je Fenster) — ein Sprung von etwa 5 MiB **einmal**, danach flach. Das ist rund das Dreifache der hergeleiteten 1,5 MiB je Seite (abgeleitet; Anteile von Heap-Überhang der Speicherbereinigung und Treiber-Puffern sind darin nicht getrennt). | [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) §Was diese ADR nicht behauptet | `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md:104–107` und Zeilen „Grundlinie“ Abschnitt 9 | ja — Zeilen „Grundlinie“ der drei Läufe | Aussage schwächer als die gedruckten Zeilen |
| F-11 | INFO | Zwei Zitier-Ungenauigkeiten im Messbericht, beide ohne Folge für die Zahlen: (a) Ergebnis 1 nennt als Wert vor der Änderung bei 3.000.000 Changes nur „3.096,6 MiB“; Reihe C hat 2.751,7 MiB, beide sind „Spitze im Run“ und nicht nach dem Nachlauf gedruckt (Abschnitt 3 trägt den Hinweis für Reihe B); (b) Abschnitt 2 sagt, der Produktionscode des Image-Stands `0e6b1b30` sei der des Slice-Diffs und „danach änderte sich nur Doku“ — `2636ca46` hat den Kommentar an `retentionInterval` in `wiring.go` umgebrochen (Kommentar, kein ausführbarer Code). | `AGENTS.md` §3.12 | `docs/reviews/messbericht-slice-retention-lauf-speicher-begrenzung.md:25–27, 61–62` | ja — `git diff 0e6b1b30..HEAD -- internal/bootstrap/wiring.go` | Zitat trägt seinen Satz nicht ganz |
| F-12 | INFO | Punkt 4 von [`LH-FA-RET-004.a`](../../spec/pflichtenheft.md) trägt eine Ressourcen-Aussage („Arbeitsspeicher … hängt an der Seitengröße (10.000 Kandidaten)“), zu der das Lastenheft keine Anforderung führt (`grep -n -i -E 'arbeitsspeicher|speicherbedarf|RAM' spec/lastenheft.md` leer). Die ADR trägt sie mit `Schärft: LH-FA-RET-004.a`, der Wortlaut ist der des Architect-Vorschlags; die Grenze zwischen „präzisieren“ und „erweitern“ liegt hier beim Planner. | Reviewer-Skill „Spec-Stratum-Verstoß“ (nicht ausgelöst, Randfall) | `spec/pflichtenheft.md:141–144` | nein — Lesen | Technik-Stratum an der Präzisierungsgrenze |
| F-13 | INFO | Zwei Verhaltensdifferenzen der neuen Lesung gegenüber `ReadChanges`, beide ohne Auswirkung auf die Sicherheit der Löschung: (a) eine Transaktion kann über zwei Seiten liegen und wird dann in zwei Datenbank-Transaktionen entleert; ein Abbruch dazwischen oder eine neu erscheinende Consumer-Position im nächsten Lauf hinterlässt eine Transaktion, der der Anfang fehlt — ADR und Handbuch nennen als Folge nur „ein Präfix der Löschmenge“; (b) die Kandidatenabfrage führt `Change.Origin` nicht mehr durch `mapper.ToChange`: eine Zeile mit unbekanntem `origin` stoppte den Lauf am Parent mit einem Fehler, sie ist jetzt Kandidat. | [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) Festlegung 6 | `internal/adapters/driven/postgresstorage/queries/queries.go:76–90`; `internal/application/usecase/retention/service.go:80–101` | nein — Lesen | Folge der Seitengrenze nicht benannt |

## Bewertung der Plan-Klausel (F-2, Empfehlung)

Die Entscheidung trifft der Verifier bzw. Planner; der Reviewer stützt sie auf gedruckte Zeilen.

- *Wortlaut:* „Trifft das nicht zu — die Spitze wächst mit der Zahl der Changes oder die Takte bleiben
  aus —, steht der Befund im Bericht und der Slice geht nicht nach `done/`.“ Die Takte bleiben nicht aus
  (Implementer 20–23, Review-Lauf 22–24, fehlgeschlagen 0). Die Spitze steigt in den zwei Läufen des
  Implementers ohne Überlappung (14,9–15,1 → 16,5–16,9 → 17,3–17,6 MiB): der Wortlaut ist erfüllt.
- *Was die Zeilen tragen:* (i) Der Sprung 1.000.000 → 2.000.000 Changes fällt mit der Änderung des
  **Startzustands** zusammen: Run 1 beginnt mit leerem `cdc.change` (Grundlinie `anon` 5,1 bis 5,4 MiB),
  Run 2 und 3 mit gefüllten Seiten ab dem ersten Takt (Grundlinie-Maximum 10,1 bis 10,7 MiB, sechs
  Fenster aus drei Läufen, zwischen 1.000.000 und 2.000.000 davor **flach**); die Größe des Sprungs
  (rund 5 MiB in der Grundlinie, 2,4 bis 2,8 MiB in der `anon`-Spitze des Runs) ist der Seiten-Anteil, kein
  Anstieg je Change. (ii) Der Schritt 2.000.000 → 3.000.000 beträgt 0,7 und 0,8 MiB (Implementer, 16,9 →
  17,6; 16,5 → 17,3) und 0,2 MiB im Review-Lauf (16,6 → 16,8); die Streuung **desselben** Schritts über
  drei Läufe ist 0,2 bis 0,8 MiB (3.000.000 Changes: 16,8; 17,3; 17,6). Der Review-Lauf liegt bei
  3.000.000 unter dem gesamten Bereich des Implementers. (iii) Die Größenordnung: 1,3 Byte je Change
  gegen 1.055 bis 1.628 Byte je Change am Stand vor der Änderung (1,03 bis 1,59 KiB, Untersuchung) sind
  0,08 bis 0,12 % (abgeleitet). Die Eigenschaft, die der Slice beseitigen sollte — lineare Abhängigkeit
  im KiB-Bereich je Change —, ist nicht mehr vorhanden.
- *Empfehlung (begründet, nicht entschieden):* Die Klausel löst **der Sache nach nicht** aus, dem
  **Wortlaut nach** schon. Der Verifier bzw. Planner sollte die Lesart in der Closure ausdrücklich
  festhalten, statt die DoD-Zeile still auf `[x]` zu halten, und die abgeleitete Angabe „1,3 Bytes je
  Change“ nicht als gemessene Steigung führen: sie ist die Differenz zweier Endpunkte über eine Stufe,
  deren Anteil der Sprung im Startzustand ist. Die Grenze der Aussage: drei Stufen bis 3.000.000
  Changes, schmale Zeilen, `n` = 2 (Implementer) plus `n` = 1 (Review); über 3.000.000 Changes hinaus
  ist nichts gemessen.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Sicherheit der Löschung (`retention/service.go`, Reihenfolge, Cursor, Positionen, Uhr, Fehlerpfade) | geprüft, ohne Befund über F-3, F-13 hinaus: jede Freigabe läuft durch `AllowsDeletion` mit denselben Eingaben wie am Parent; Uhr früher gelesen = Alter höchstens unterschätzt; Lese- und Lösch-Fehler enden den Lauf; Fortschrittsprüfung vor der Löschung; sieben Use-Case-Mutationen reproduziert, davon fünf mit Assertion-Rot (M2, M3, M4, M6, M7), zwei durch Zeitüberschreitung (M1, M5, F-4) |
| Kandidatenabfrage und Sortierung (`queries.go`, Kollation, Cursor, Parameter) | geprüft, ohne Befund: Keyset über den eindeutigen Primärschlüssel `change_id` in einer Anweisung; das Schema trägt keine Spalten-Kollation; der wegfallende `JOIN` auf `cdc.source_table` filterte wegen `required: true` und Fremdschlüssel keine Zeile; Parameter-Typen (`$1` Quelle, `$2` Text, `$3` Zahl) tragen; die Ordnung `0bf-…` vor WAL-Kennungen ist für das Keyset nicht erforderlich |
| Löschung (`DeleteChanges`, Array-Parameter, Waisen-Bereinigung) | geprüft, ohne Befund: `= ANY($1)` mit einem `[]string`, kein Parameter-Limit bei 10.000 Kennungen; im Diff unverändert; eine über Seiten geteilte Transaktion bleibt bis zur letzten Change-Zeile stehen |
| Adapter, Übersetzung, Fehlerklassen (`store.go`, `sqlexec/translate.go`) | geprüft, ohne Befund: `limit < 1` → `ErrNonPositiveLimit`, leere Quelle → `ErrEmptyIdentifier` (Store-Tier-Test rot ohne die Prüfungen laut Messbericht), Treiber-Fehler in Klasse `storage`, Domänen-Fehler unklassifiziert; `UnixNano` wie `ReadChanges` (`translate_test.go` bindet `Unix` gegen `UnixNano`), keine Zeitzonen-Übersetzung |
| Rollenrechte (`tools/schema/nacharbeit-roles.sql`, `cdc_admin`) | geprüft, ohne Befund: Datei im Diff unverändert; `TestRetentionCandidatesRunUnderTheirRoles` liest unter einer `cdc_admin`-Login-Identität und scheitert unter `cdc_capture`/`cdc_reader` mit SQLSTATE 42501; der 25.000-Changes-Test läuft unter einer `cdc_admin`-Login-Identität ohne Superuser |
| Vier Test-Fakes und Use-Case-Tests (`service_test.go`, drei Schnittstellen-Fakes) | geprüft, ohne Befund über F-4 hinaus: der prüfende Fake begrenzt je Aufruf und zählt; die Erwartungen von `LH-FA-RET-002` bis `004` bleiben (eine Fake-Ausstattung geändert, Befund 1 des Messberichts trägt); Menge bei Seitengröße 1, 2, 3, N, größer N gegen handgeschriebene Erwartung |
| Store-Tier-Tests und ihre Lage (`retentioncandidates_test.go`, `internal/bootstrap/retention_pages_test.go`) | geprüft, ohne Befund über F-7 hinaus: die Begründung „Dateiname-Ordnung vor `DROP SCHEMA`“ trägt (`DROP SCHEMA IF EXISTS cdc CASCADE` steht nur in `store_test.go:50` und `tableactivation_test.go:39`, beide sortieren nach `retentioncandidates_test.go`; die übrigen Treffer sind Kommentare), die Lage des Use-Case-Tests im Paket der Verdrahtung folgt aus `.a-check.yml` (`make a-check` Exit 0); die Bindung „nie mehr als vorher“ beruht auf drei Zählungen (`Deleted` gegen unabhängige Zählung, verbleibende Changes, freigegebene Changes, die noch stehen) und ist an M2/S4/S6 rot |
| Messbericht und Nachmessung | geprüft, ohne Befund über F-2, F-9, F-10, F-11 hinaus: alle nachgerechneten Zahlen stimmen (siehe Prüfungen oben, u. a. 1,31 Byte je Change, 5,88 bis 6,63 Millionen, Delta zu Reihe J); Ursprung je Zahl und Lauf genannt; Aufräum-Schritte benannt (dangling Volumes 34, kein `prune`); die Bewertung „Erwartung für die Form belegt, nicht als Zahl; Rest unerklärt“ ist ehrlich und im Bericht offen geführt |
| Bewertung der Richtgröße (Messbericht Abschnitt 7) | geprüft, ohne Befund: 9.802 bis 11.056 Zeilen/s aus den gedruckten Kopierdauern nachgerechnet, ×600 s = 5,9 bis 6,6 Millionen; Review-Lauf 9.076 bis 10.389 Zeilen/s (5,4 bis 6,2 Millionen); der Wert 4.000.000 bleibt, die Begründung („folgt der Kopierdauer, Speicher geht nicht ein“) trägt |
| Handbuch (Version 1.62, `### Änderungshistorie`, Abschnitte Retention, Backfill, Grenzwerte) | geprüft, ohne Befund über F-1, F-3, F-6 hinaus: `Version:` 1.61 → 1.62 mit Zeile 1.62; Zahlen gemessen/abgeleitet gekennzeichnet (Werte gegen den Messbericht nachgerechnet); Anwender-Sprache; die neue Aussage zu Seiten, Nicht-Atomarität und einmal gelesenen Positionen stimmt mit dem Code; keine neue Betreiber-Oberfläche (kein Zug nötig); Wegfall der Bemessungsregel „2 KiB je Change plus 64 MiB“ und des `--memory`-Hinweises entspricht dem Plan |
| `spec/pflichtenheft.md`, `harness/README.md`, `harness/sensors/*.md`, `harness/targets/bench-backfill.md` | geprüft, ohne Befund über F-12 hinaus: Punkt 4 ohne ADR- und Slice-Bezug, Änderungshistorie-Zeile vorhanden; Nenner 1064 (692 · 32 · 130 · 210 = 1064) und 2601 (2167 gedeckt, 83,31 %) mit Lauf und gedruckter Zeile, beide am Review nachgezählt bzw. nachgemessen; die Differenz 34 als abgeleitet gekennzeichnet; `harness/README.md` nennt die Kandidaten-Seiten-Tests in der Zeile `make test-store` (wahr); der Halbsatz in `bench-backfill.md` ist entfernt |
| Suchlauf-Feld des Plans (§3), beide Stände | geprüft, ohne Befund: alle Zahlen an `8cd39719`, `5b1f7762` und am Arbeitsbaum bestätigt; „Nicht gefunden“ trägt (`spec/architecture.md` sechs Zeilen `retention|bereinigung`, keine Aussage zur Lesung); Lokatoren in `coverage-gate.md` stehen unverändert (zeilenneutraler Kommentar-Umbruch) |
| Release-Aussage (Plan §7) und Reichweite (`v0.1.0`–`v0.1.2`) | geprüft, ohne Befund: die Zeile `ReadChanges(ctx, outbound.ChangeQuery{Source: command.Source})` steht in allen drei Tags, kein Backfill-Paket in den Tags, `git diff v0.1.2 8cd39719` leer; 0,85 bis 1,31 GiB nachgerechnet und als hergeleitet gekennzeichnet |
| Kommentare (`AGENTS.md` §3.7, diff-skopiert) und Hard Rules | geprüft, ohne Befund über F-8 hinaus: kein `slice-*`/`welle-*` in Produktionscode-Kommentaren, keine Suppression, kein host-lokaler Pfad, Docker-only (alle Gates und Mutationen über `make`/Repo-Skripte bzw. Docker-Aufruf des Toolchain-Images), Commit-Betreffs ohne `SPEC-*`/`ARC-*`, jeder Slice-Commit nennt `LH-*` und `ADR-*`; die zwei `git mv`-Lifecycle-Commits `e486beb8` und `ff387a69` sind reine Renames (`git show --stat -M`: 0 Einfügungen, 0 Löschungen), die Inhaltsänderung `8a68d2b2` steht in einem eigenen Commit |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 7 |

**Finding-Klassen dieses Laufs:** Nachzug widerspricht dem Nachbarn · Plan-Klausel ohne dokumentierte
Auswertung · Offene Bewertung eines Sicherheitsfensters · Rot nur durch Zeitüberschreitung · Roter
Zwischencommit auf `main` · Vorher-Nachher-Sprache in Handbuch-Prosa · Grenzwert ohne Bindung auf dieser
Test-Ebene · Kommentar an der Klassengrenze · Kennzahl streut mit Cache-Anteil · Aussage schwächer als die
gedruckten Zeilen · Zitat trägt seinen Satz nicht ganz · Technik-Stratum an der Präzisierungsgrenze ·
Folge der Seitengrenze nicht benannt

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Die Sicherheit der Löschung trägt: jede Freigabe läuft durch
`AllowsDeletion` mit denselben Eingaben wie am Parent, die Kandidatenmenge ist gleich der bisherigen, die
Wanduhr wird höchstens früher gelesen (Löschung nie früher), das Keyset über `change_id` ist auf beiden
Seiten der Löschung stabil, die Löschung übergibt einen Array-Parameter, Fehlerpfade enden den Lauf, und
die Eingabeseiten-Mutationen (sieben am Use Case, sechs an der Abfrage) färben Tests rot — der
25.000-Changes-Test bindet „nie mehr, nie weniger als eine unabhängige SQL-Zählung“ (M2, S4, S6). Alle
Gates sind am Stand `b2484add` grün (eigene Läufe und CI `ci`/`examples`/`e2e` success). Die drei MEDIUM
sind keine Fehler des Löschpfads: F-1 ist ein Widerspruch zwischen Aussagen im selben Träger (Wortlaut),
F-2 und F-3 sind offene Bewertungen, die der Plan selbst dem Verifier bzw. dem Architect zuweist. Sie
sind Bedingungen des `done/`-Übergangs, nicht des Codes: **vor `done/` und vor dem Server-Release
`v0.2.0`** sollten (a) der Architect die Frage „neuer Consumer während eines Laufs“ samt der Folge der
Seitengrenze (F-13 a) bewertet und in Plan §6 als Ausgang festgehalten haben, (b) Verifier bzw. Planner
die Lesart der Plan-Klausel (F-2) ausgesprochen haben, (c) der Handbuch- und Pflichtenheft-Wortlaut
(F-1) mit dieser Lesart übereinstimmen. Die vom Nutzer gestellte Frage nach dem roten
Zwischencommit `d49e248f`: kein Sicherheitsproblem (der Zustand liegt zwischen zwei Commits desselben
Arbeitsgangs, `HEAD` ist grün), aber eine `git bisect`-Falle (F-5, LOW).

**DoD-Zeile „Review durchgeführt“:** bleibt offen. F-1 verlangt einen Wortlaut-Nachzug an
Handbuch, Pflichtenheft und `PageSize`-Godoc, der von der Entscheidung zu F-2 abhängt und über einen
Rückgabe-Pfeil an den Implementer läuft (Fixrunde); die Checkbox wird dann regulär nachgezogen.

**Übergabe:** F-1 an den Implementer nach der Entscheidung zu F-2 (Fixrunde, Wortlaut); F-2 an
Verifier bzw. Planner (Lesart der Klausel; Empfehlung oben); F-3 und F-13 a an den Architect (Bewertung
des Fensters und der Transaktions-Zerschneidung, Ausgang in Plan §6); F-4, F-5, F-6 nach Ermessen des
Implementers; F-7 bis F-13 zur Kenntnis (F-12 an den Planner). Die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg**
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt); er ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
