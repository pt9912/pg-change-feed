# Review-Report: Closure-Notizen welle-backfill-bestand — 2026-09-25

**Review-Art:** Closure-Note-Review (inferentieller Nachlauf zum Struktur-Gate) — geprüft
werden die Closure-Notizen (§7) der elf Slices der Welle `welle-backfill-bestand` und die
Results-Notiz der Welle gegen die drei Pflicht-Inhalte (a) konkretes Lernsignal, (b) konkretes
Folge-Slice, (c) konkrete Architektur-Beobachtung (`v6.9.0` · `regelwerk/modul-11-*.md`
§Schritt 5, wörtlich im Skill wiedergegeben). Für die Results-Notiz zusätzlich gegen
`welle-results.template.md` und `close-welle.md` Schritt 3. Kein DoD-Abgleich, keine fachliche
Bewertung der Slices (Verifier/Validator), keine Struktur-Prüfung (Struktur-Gate).

**Gegenstand:** elf Closure-Notizen unter `done/` (`ls docs/plan/planning/done` mit dem Präfix
`slice-backfill-`: 11 Dateien) und `docs/plan/planning/done/welle-backfill-bestand-results.md`,
Stand HEAD `eb3ba91e` (Baum sauber, gepusht).

**Skill:** `.harness/skills/closure-note-reviewer.md` (Status Accepted, Stand HEAD `eb3ba91e`).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Template `v6.9.0` · `templates/docs/plan/planning/slice.template.md` §Closure-Notiz und
  §2 DoD (Lerneintrag, Risiko-Ausgänge, drei Paarungen), Results-Template
  `templates/docs/plan/planning/welle-results.template.md`
- Lifecycle-Pflicht (`v6.9.0` · `regelwerk/modul-05-planning-harness.md` §Closure- und
  Lerneintrag-Regeln: ein Übergang nach `done/` ohne Lerneintrag ist Ablage, keine Closure),
  `.claude/commands/close-welle.md` Schritt 3, `.claude/commands/implement-slice.md` Schritt 23 und 24
- Struktur-Gate für denselben Stand: `make docs-check`/`make gates` Exit 0 (vom Auftraggeber
  gemessen, hier nicht doppelt gemeldet; der Lauf mit dem Report im Baum steht unter „Eigenständig
  durchgeführte Prüfungen“)
- Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt` (Ausgang je
  Register-Eintrag), `AGENTS.md` §3.7, §3.12, §3.13, Formvorbild
  `docs/reviews/review-slice-backfill-e2e.md`
- [`LH-FA-CAP-009`](../../spec/lastenheft.md), [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-/Planner-Bericht
übernommen; die Stichprobe der Results-Notiz umfasst mehr als die geforderten fünf Zahlen):

- **Zählung der Slices:** `ls docs/plan/planning/done | grep -c '^slice-backfill-'` druckt `11`; die elf
  Namen der Results-Notiz decken sich mit dem `ls`.
- **`make doc-trace`** (Exit 0): gedruckt „80 Anforderung(en), 2 Waise(n).“, die Waisen sind
  [`LH-FA-CFG-007`](../../spec/lastenheft.md) und [`LH-FA-CFG-008`](../../spec/lastenheft.md);
  [`LH-FA-CAP-009`](../../spec/lastenheft.md) trägt `E2E` und `ok`.
- **CI-Kennungen** (`gh run view`): Lauf `36108615045` (e2e), `36108615036` (ci), `36108615027` (examples)
  zum Push `c82d3333` je `success`; beide Legs „image + test-integration (PostgreSQL 17)“ und „(PostgreSQL 18)“
  `success`. Aus dem Log von `36108615045`: `Replay-Invariante: Snapshot-Position 28731112, 55 Backfill-Changes,
  WAL-Changes davor 412 und dahinter 131, 54 Zeilen im Quellstand` (PostgreSQL 17) und `… 56 …, 406 …, 166 …, 56 …`
  (PostgreSQL 18); `DB-Adapter-Coverage: 82.51% (gedeckt 873 von 1058 Statements; Profile gemergt:
  store,replication)` in beiden Legs; die Zeilen „Backfill-Happy-Path“ (5 Bestandszeilen, `READ changes=6`),
  „Backfill-Schema-Version-Beleg“, „Backfill-Startposition“, „Backfill-Boundary“, „Backfill-DDL-Fenster“,
  „Backfill-Negative“ und „Leerlauf-Bestätigung“ je Leg vorhanden. Die drei CI-Läufe, die die Slice-Notizen
  nennen (`36065957210` an `43137ebf`, `35988269802` an `455bcdef`, `36092733208` an `427f6d1b`), sind je
  `success` mit beiden Legs.
- **Zählerstände** (`ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence | wc -l`): 31, 21, 13 und 11 für
  die vier Mega-Einträge, außerdem 8, 8, 8, 7, 7, 5, 4 und sechsmal 3 für die übrigen Zeilen der Lese-Schritt-Tabelle
  — alle gleich den Zellen. `ls docs/plan/planning/observations/BEO-PGC | wc -l` druckt `113`; 39 Verzeichnisse
  tragen ab 3 Dateien `evidence/`; die erste Ausgangs-Zeile (`grep -m1 -o -i -E 'verkörpert|geplant|gestrichen|offen'`)
  liefert 34 verkörpert, 2 geplant, 2 gestrichen und einmal leer (`rollen-verdrahtung`, dort „eingetreten“) —
  wie in der Notiz. Alle 16 Einträge ab 3× mit einer `evidence/`-Datei eines `slice-backfill-*`-Slice stehen in
  der Tabelle (17 Zeilen: 16 plus `zusatzkontext-kopplung-breiter-als-dod-wortlaut` wegen nicht zugewiesenen
  Ausgangs).
- **Weitere Zahlen und Befehle der Trigger-Tabelle:** `sed -n 62p tools/schema/nacharbeit-administration.sql` druckt die
  `request_kind`-Menge mit fünf Werten; `git diff --name-only c82d3333..HEAD -- <Pfade der Notiz> | wc -l` druckt `0`
  (die Pfadliste schließt `tools/harness/` aus, dort ändert nur `db-coverage.sh` — in „Offener Nachlauf“ der Notiz
  benannt); `THRESHOLD ?= 80` in `harness/mk/coverage.mk` und `DB_COVERAGE_THRESHOLD:-80` in
  `tools/harness/db-coverage.sh`; `docs/plan/carveouts` trägt nur `.gitkeep`; `grep -rl 'make bench' .github/workflows`
  druckt 0 Treffer; `git tag -l` nennt die fünf Tags der Notiz; die drei SDK-Quellen tragen `0.2.0`; die Änderungshistorie
  des Handbuchs führt die Zeilen 1.44 bis 1.59 (16 Zeilen, `Version: 1.59`); die vier Folge-Slices liegen als Dateien in
  `open/`, zehn `slice-transformationen-*` ebenfalls; die Anker `seit welle-backfill-bestand` stehen in `AGENTS.md`,
  `.harness/skills/reviewer.md`, `harness/sensors/db-adapter-coverage.md` und `docs/plan/planning/observations/README.md`;
  Digest-Präfixe `sha256:862dfb04` (d-migrate, `Makefile`) und `7456ef82` (`e2e.yml`) stimmen; die Rechnung
  6.470 mal 600 = 3.882.000 und die Rundungsart (abgerundet auf eine Stelle) folgt der des Handbuchs (1.595–1.622).
- **Register-Paarung (c):** jede in den elf Notizen und in der Results-Notiz genannte Kennung `BEO-PGC/<slug>`
  (71 verschiedene in den Slice-Notizen) besitzt ein Verzeichnis mit nicht leerem `evidence/` — mit der einen
  benannten Ausnahme `limit-fortsetzung-innerhalb-einer-position`; `docker volume ls -q -f dangling=true | wc -l`
  druckt `34`, wie die Notiz.
- **`make gates`** mit dem Report im Baum, ungefiltert in eine Log-Datei, Exit direkt gesichert: Exit 0; gedruckt
  „baseline-verify: v6.9.0 OK — 54 Dateien“, „coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%“,
  „d-check: 1139 Datei(en) geprüft, 0 Befund(e)“, „commit-traceability: OK“, „generated-sync: OK“,
  „gesamt: 0 Befund(e)“ (a-check).
- **Nicht gefahren (Grenze):** `make pin-stale-dmigrate` (Netz, advisory), `tools/bench-backfill.sh` (Lauf
  `20260925T084803Z`: Zahlen 6.470 Zeilen/s, 433,0 MiB, 170,8 MiB sind nur in der Results-Notiz, kein Log im Repository),
  `make test-store`/`make test-replication` und die `make gates`-Zahlen der Notiz (Stand `32028d4f`) — gelesen, nicht
  nachgefahren. Der Lauf `36116700952` (e2e, Push `eb3ba91e`) lief zum Prüfzeitpunkt (`in_progress`).

---

## Findings

Jedes Finding folgt dem §Output-Schema des Closure-Note-Reviewer-Skills; die Spalte `Klasse` ist das stabile
Fehlermuster für den Steering-Loop-Zähler.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die vier Notizen `slice-backfill-e2e`, `-bench-richtgroesse`, `-sdk-origin` und `-slot-leerlauf-bestaetigung` sagen zum Validator-Schritt nichts (kein Lauf, kein „entfällt ausdrücklich“; `grep -c -i validator` je Datei 0). Sechs Vorgänger-Notizen (spec-nachzug, change-origin, snapshot-reader, run-usecase, run-store, sql-administration) verschieben die Validierbarkeit von [`LH-FA-CAP-009`](../../spec/lastenheft.md) auf „den Wellen-Beleg (`slice-backfill-e2e`)“ mit „Kein stilles Überspringen“; im E2E-Slice und in der Results-Notiz (kein Treffer für „Validator“) wird sie nicht aufgenommen. Der Folgeschritt (b) ist zugesagt und hat weder Träger noch Ausgang. | Closure-Inhaltspflicht (b) · `implement-slice` Schritt 23 („dann explizit sagen statt still überspringen“) | `docs/plan/planning/done/slice-backfill-e2e.md:399-494`; `docs/plan/planning/done/slice-backfill-bench-richtgroesse.md:361-480`; `docs/plan/planning/done/slice-backfill-sdk-origin.md:259-371`; `docs/plan/planning/done/slice-backfill-slot-leerlauf-bestaetigung.md:634-753`; Verweise: `slice-backfill-run-usecase.md:497-499`, `-change-origin.md:416-418`, `-spec-nachzug.md:339`, `-sql-administration.md:541-543`, `-snapshot-reader.md:495-497`, `-run-store.md:495-497` | nein — Floskel-/Auslassungs-Erkennung ist inferentiell; `grep -c -i validator` je Datei zeigt nur die Abwesenheit des Wortes | Zugesagter Folgeschritt ohne Träger |
| F-2 | LOW | Zwei von der Closure `slice-backfill-change-origin` gemeldete Ungenauigkeiten stehen zu HEAD unverändert und tragen weder einen Folge-Slice noch einen Register-Eintrag noch eine Zeile in der Results-Notiz: der Kommentar an der gRPC-Nachricht „mit denselben Feldern wie der Domain-Typ“ (der Domain-Typ trägt `origin`, die Live-Nachricht nach dem Pflichtenheft nicht) und der Kommentar am Feld `allowDestructive` „(nur bekannte Fremdobjekte blockieren)“ neben der Klasse „View-Signatur“. Die Notiz nennt beide selbst „gemeldet“ bzw. „ohne Träger-Slice“. | Closure-Inhaltspflicht (b) · `AGENTS.md` §3.13 (Meldung ist das Übergabe-Artefakt) | `docs/plan/planning/done/slice-backfill-change-origin.md:349-353` und `:428-431`; `proto/cdc/stream/v1/changestream.proto:12-13`; `tools/schema/rolloutguard/guard.go:53-54` | ja — `sed -n 12,13p proto/cdc/stream/v1/changestream.proto` und `sed -n 53,54p tools/schema/rolloutguard/guard.go` am Stand HEAD | Gemeldete Ungenauigkeit ohne Träger |
| F-3 | LOW | Die Results-Notiz führt die Spec-Lücke zu [`SPEC-008`](../../spec/pflichtenheft.md) (Klasse `transient` ohne tragendes Element) als „aufgelöst über die Architect-Entscheidung zur Wiederholungsform (Start-Trigger von `slice-capture-transient-wiederholung`)“. Die Entscheidung liegt nicht vor (`docs/plan/adr` endet bei [`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md), der Plan des Folge-Slice sagt „entscheidet der Architect in einer ADR vor dem Start“); das Template verlangt bei der Spec-Lücke ein `aufgelöst über` mit `LH-*`- oder ADR-Kennung. Die Lücke ist adressiert, nicht aufgelöst. | Results-Template §Steering-Loop-Einträge · Closure-Inhaltspflicht (b) | `docs/plan/planning/done/welle-backfill-bestand-results.md:203-209`; `docs/plan/planning/open/slice-capture-transient-wiederholung.md:56-58` | ja — `ls docs/plan/adr` und der Plan des Folge-Slice | Spec-Lücke als aufgelöst formuliert, Entscheidung ausstehend |
| F-4 | LOW | Die Streichung von `sdk-decoder-verhalten-am-neuen-feld-ungemessen` steht in der Results-Notiz als „in den veröffentlichten Packages ist keine strikte Dekoder-Einstellung gemessen“ — das liest sich als „nicht gemessen“, gemeint ist (Verdikt und `state.md`): ein `git grep` nach strikten Mustern druckt 0 Treffer an `sdk-csharp-v0.1.0` und `sdk-python-v0.1.0`. Die Begründung ist damit eine Quelltext-Suche, kein Dekoder-Lauf; die im Verdikt genannte Restgrenze (die `0.1.0`-Packages sind nicht erneut ausgeführt) steht in der Notiz nicht. | Closure-Inhaltspflicht (a) („weil X“ statt Behauptung) · `AGENTS.md` §3.12 Instanz B | `docs/plan/planning/done/welle-backfill-bestand-results.md:253-254` | ja — `git grep -n -i -E 'UnmappedMemberHandling\|MissingMemberHandling\|FAIL_ON_UNKNOWN\|Disallow\|strict' sdk-csharp-v0.1.0 -- sdks` und der Wortlaut der `state.md` des Eintrags | Streichungs-Begründung mehrdeutig, Restgrenze fehlt |
| F-5 | LOW | Drei Slice-Notizen adressieren einen Ausgangs-Kandidaten an den Lese-Schritt der Welle-Closure: „Fake scheitert ab Aufruf n, je Aufrufstelle eine Mutation“ (run-usecase), die wertfreie Bindung der Umrechnung einer Konstanten (bench, V-3) und „Filter-Eingabe verlangt Zeilen außerhalb des Filters“ (sql-administration). Die Results-Notiz nennt für `negativtest-ohne-bindung-an-seine-eingabe` nur „verkörpert; Deckel; Mutations-Harness verworfen“; keiner der drei Kandidaten hat einen benannten Ausgang, und die `state.md` des Eintrags trägt den Ausgangs-Kandidaten weiter als „nicht getroffen“ (Zeilen 58-62). Das Verdikt §3.4 sagt „kein neuer Träger“ ohne sie zu nennen. | Closure-Inhaltspflicht (b) · `close-welle.md` Schritt 3 (Ausgang je Eintrag in der `state.md`) | `docs/plan/planning/done/welle-backfill-bestand-results.md:232`; `docs/plan/planning/done/slice-backfill-run-usecase.md:427-437`; `docs/plan/planning/done/slice-backfill-bench-richtgroesse.md:399-411`; `docs/plan/planning/done/slice-backfill-sql-administration.md:486-487`; `docs/plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md:58-62` | ja — `grep -n 'Aufrufstelle\|nicht getroffen' docs/plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md` | Kandidat mit Adresse Lese-Schritt ohne benannten Ausgang |
| F-6 | LOW | Der Steering-Loop-Eintrag steht in zehn der elf Notizen als ein Sammelabsatz von 19 bis 36 Zeilen (`slice-backfill-sdk-origin`: 13), in dem Lernsignal, geschärfte Anwendung, Sensor-Umfang, Spec-Lücke und Kandidat ineinanderlaufen (größte: snapshot-reader 36, change-origin 31, sql-administration 31, run-usecase 29 Zeilen); der gesamte §7 umfasst je 96 bis 148 Zeilen. Die Substanz (a)/(b)/(c) ist in jeder Notiz vorhanden, aber nur durch vollständiges Lesen von der Chronik der Rollen-Kette zu trennen. | Closure-Inhaltspflicht (Klarheit; LOW-Klasse des Skills) | `docs/plan/planning/done/slice-backfill-snapshot-reader.md:416-451`; `docs/plan/planning/done/slice-backfill-change-origin.md:356-386`; `docs/plan/planning/done/slice-backfill-sql-administration.md:478-508`; `docs/plan/planning/done/slice-backfill-run-usecase.md:420-448`; übrige sechs Notizen gleiche Form (`slice-backfill-sdk-origin` ausgenommen) | nein — Nachvollziehbarkeit ist inferentiell | Lerneintrag im Sammelabsatz |
| F-7 | INFO | `slice-backfill-run-usecase` benennt die Kosten der Ausschluss-Lesung je Block bei großer Blockzahl als „nicht gemessen“ („F-10, benannt, nicht gezählt“); weder ein Register-Eintrag noch eine Zeile der Results-Notiz noch der Bench (Kopierdauer, Speicher) führt sie. Hinweis ohne erwartete Aktion des Implementers; Tracker-Nachtrag durch die Planning-Rolle, falls gewollt. | Closure-Inhaltspflicht (b) | `docs/plan/planning/done/slice-backfill-run-usecase.md:398-401` und `:462-464` | nein | Benannte Lücke ohne Adresse |
| F-8 | INFO | Die Results-Notiz fasst die Speicher-Zahlen enger als das Handbuch: „Spitze von 401,7 bis 467 MiB“ und „für einen Run über 1.000.000 Zeilen 1.544 MiB“; das Handbuch nennt dazu, dass die Spitze im Run eine Untergrenze ist (höchster gemessener Wert der Stufe 641,7 MiB, 20 s nach dem Run) und zwei Runs über 1.000.000 Zeilen (435 und 1.544 MiB). Die Aussage „Speicher ist ungeklärt“ trägt beides; der Folge-Slice `slice-backfill-speicher-untersuchung` steht als Datei. | `AGENTS.md` §3.12 Instanz A (Zahl mit Ursprung) | `docs/plan/planning/done/welle-backfill-bestand-results.md:124-130`; `docs/user/benutzerhandbuch.md:1652-1667` | ja — die zwei Textstellen nebeneinander lesen | Zusammenfassung enger als der Träger |
| F-9 | INFO | Zeitformen in Records: die Results-Notiz nennt die Prüfung dieser Closure „erfolgt nach diesem Zug“ und führt „Offener Nachlauf“ (erster e2e-Lauf mit Stufe 80); der Lauf `36116700952` zu `eb3ba91e` lief zum Prüfzeitpunkt noch. Die elf Slice-Notizen tragen Stand-Aussagen („die Folge-Slices der Welle liegen als Dateien in `open/`“, Welle „(offen)“), die zu HEAD nicht mehr gelten; die Results-Notiz nennt für die Folge-Slices den Ort `docs/plan/planning/open/` statt der Kennung allein. Records, keine Aktion; der Nachlauf-Beleg gehört Verifier und nächster Closure. | `AGENTS.md` §3.7 (Zustandsfelder) · `AGENTS.md` §3.10 | `docs/plan/planning/done/welle-backfill-bestand-results.md:280`, `:406-410`, `:438-443` | nein — der Lauf `36116700952` ist ein späterer Beleg | Stand-Aussage im Record |

## Markierung je Closure-Notiz (Prüfauftrag Modul 11 §Schritt 5)

(a) Lernsignal mit „weil X“, (b) konkretes Folge-Slice oder benannte Adresse, (c) beobachtbare
Architektur-Aussage; `ja` heißt konkret getragen, ein Finding-Verweis heißt mit Einschränkung.

| Notiz | (a) Lernsignal | (b) Folge | (c) Architektur-Beobachtung | Finding |
|---|---|---|---|---|
| `slice-backfill-spec-nachzug` | ja — Hedge „Warn-Spalte(n)“ wurde erst durch die Hedge-Suche gefunden, die „leer“-Suche traf nur einen Teil | ja — zwei Handbuch-Meldungen mit Adresse, Startposition an `slice-backfill-e2e` | ja — `SPEC-022` Tabellenzeile „Reihenfolge“ ohne Vorbehalt | F-6 |
| `slice-backfill-row-image-gemeinsam` | ja — das Suchlauf-Feld ist selbst ein Zahlen-Träger (17 gegen gemessen 16, Parent-Befehl mit `HEAD`) | ja — K1 an `slice-transformationen-kern-rename`, Handbuch-Übergaben | ja — die Kopplung hängt an einer positionalen Signatur | F-6 |
| `slice-backfill-change-origin` | ja — Lauf 6 nannte sich Beleg, die Mutation blieb grün, weil d-migrate den Blocker selbst abbricht | ja — Meldungen mit Adresse; ein Rest ohne Träger | ja — `DROP VIEW` verwirft die Rechteliste (abgeleitet, als nicht gemessen gekennzeichnet) | F-2, F-6 |
| `slice-backfill-snapshot-reader` | ja — Cast auf den Spaltennamen ist nicht die Ausgabefunktion von `pgoutput`, sieben von 86 Spalten weichen ab | ja — Übergaben mit Adressen, Risiken mit Ausgang | ja — Lesepfad im Text-Ergebnisformat, Unterpaket `snapshotlogic` | F-6 |
| `slice-backfill-run-usecase` | ja — der Fake scheiterte bei jedem Aufruf, deckte nur die erste Aufrufstelle | ja — Übergaben an fünf Pläne; Architect-Fragen F-8/F-9 mit Adresse | ja — `SchemaStorePort` als achter Pflicht-Port | F-5, F-6, F-7 |
| `slice-backfill-run-store` | ja — Test unter Superuser ließ den fehlenden `cdc_admin`-Grant unentdeckt | ja — Übergaben an vier Pläne, Architect-Frage adressiert (Verdikt: keine ADR) | ja — Rollenschnitt in [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) ohne Antrags-Queue | F-6 |
| `slice-backfill-sql-administration` | ja — Suchlauf nach Symbol fand den Zählwort-Rest nicht, Guard-Lauf führte die zugesagten Rechte nicht | ja — Start-Trigger zweier Slices erfüllt, Register-Adresse | ja — nur `backfill` ist an die Quelle gebunden, vier Antragsarten lesen ohne Quellbezug | F-5, F-6 |
| `slice-backfill-e2e` | ja — die sichtbar scheiternde DDL-Form belegt das Fenster nur für Fehler mit Meldung, die stille Rewrite-Form deckte den Verlust auf | ja — Replay-Tier als Architect-Entscheidung, später [`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md) | ja — `pg_publication_tables` wartet auf jede Sperre, `ACCESS EXCLUSIVE` blockiert die Slot-Anlage | F-1, F-6 |
| `slice-backfill-bench-richtgroesse` | ja — DoD über ein Ganz-Target hängt am Messhost (Exit 2 an `LH-QA-PER-001`) | ja — `slice-backfill-slot-leerlauf-bestaetigung`; Speicher-Lücke mit Adresse | ja — WAL-Rückstand hängt an WAL ohne Inhalt für die Publication, nicht am Backfill | F-1, F-5, F-6 |
| `slice-backfill-slot-leerlauf-bestaetigung` | ja — die Fitness-Function-Zeile war nicht erfüllbar (Store-Mutation färbt nicht rot), erst der Reviewer maß es | ja — [`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md), Grenzen mit Adresse | ja — Bestätigung im Leerlauf, Bedingung „keine offene Transaktion“ trägt nur der Unit-Test | F-1, F-6 |
| `slice-backfill-sdk-origin` | ja — Randwerte `null`/leer standen nirgends, drei Sprachen entschieden nach ihrer JSON-Bibliothek verschieden | ja — Release-Zug, Übergabe der drei Tags an den Betreiber | ja — `wireOrigin` in Kotlin, Fixtures handgeschrieben | F-1 |
| `welle-backfill-bestand-results` | ja — Steering-Loop-Einträge mit Auslöser, Zähler und Anker | ja — vier Folge-Slices als Dateien in `open/` | ja — Tier der Replay-Invariante, Rest-Grenzen | F-1, F-3, F-4, F-5, F-8, F-9 |

Kein HIGH: keine Notiz ist eine Floskel; jede trägt alle drei Pflicht-Inhalte mit Ursache
(„weil X“), benanntem Folge oder Adresse und beobachtbarer Aussage. Die Results-Notiz trägt Lerneinträge (neun
Steering-Loop-Punkte mit `liegt in` und Anker), den Lese-Schritt für alle 17 Einträge der Welle mit Ausgang,
den Zeiger aufs Register und vier Folge-Slices als Dateien.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `done/slice-backfill-spec-nachzug.md` §7 | geprüft, ohne Befund über F-6 hinaus: (a), (b), (c) konkret; Zähler `3×` und `27×` gegen `evidence/` gemessen |
| `done/slice-backfill-row-image-gemeinsam.md` §7 | geprüft, ohne Befund über F-6 hinaus: Validator „entfällt ausdrücklich“ mit Begründung, Träger-Übergaben mit Text der Empfänger |
| `done/slice-backfill-change-origin.md` §7 | geprüft, ohne Befund über F-2 und F-6 hinaus; die eine Zahl mit gemessenem Ursprung (16 gegen 18 Zeilen) ist an vier Ständen genannt |
| `done/slice-backfill-snapshot-reader.md` §7 | geprüft, ohne Befund über F-6 hinaus: Cast-Befund (sieben von 86 Spalten), Runner-Beleg des Vorgängers als aufgelöst geführt (`gh run view` bestätigt) |
| `done/slice-backfill-run-usecase.md` §7 | geprüft, ohne Befund über F-5, F-6, F-7 hinaus; der Lauf `35988269802` an `455bcdef` bestätigt |
| `done/slice-backfill-run-store.md` §7 | geprüft, ohne Befund über F-6 hinaus: Architect-Frage zu [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) und [`SPEC-029`](../../spec/pflichtenheft.md) ist im Verdikt `architect-verdict-backfill-schema-klasse-rollen` beantwortet (keine ADR) |
| `done/slice-backfill-sql-administration.md` §7 | geprüft, ohne Befund über F-5, F-6 hinaus; der Zählwort-Kandidat ist in `AGENTS.md` §3.13 verkörpert |
| `done/slice-backfill-e2e.md` §7 | geprüft, ohne Befund über F-1 und F-6 hinaus: der Lauf `36065957210` an `43137ebf` bestätigt, Grenzen (V-2 Mutationen der Negative-Phase nicht gefahren) benannt |
| `done/slice-backfill-bench-richtgroesse.md` §7 | geprüft, ohne Befund über F-1, F-5, F-6 hinaus; Rundung (abgerundet auf eine Stelle) und Zahlen mit Ursprung `übernommen`/`abgeleitet` gekennzeichnet |
| `done/slice-backfill-slot-leerlauf-bestaetigung.md` §7 | geprüft, ohne Befund über F-1, F-6 hinaus; der Lauf `36092733208` an `427f6d1b` bestätigt |
| `done/slice-backfill-sdk-origin.md` §7 | geprüft, ohne Befund über F-1 hinaus; die Versions-Aussage (drei Quellen `0.2.0`, kein Tag `sdk-*-v0.2.0`) gegen `git tag -l` bestätigt |
| Results: „Was wurde geliefert?“ und „Verifikation“ | geprüft, ohne Befund über F-8 hinaus: alle in „Eigenständig durchgeführte Prüfungen“ nachgefahrenen Zahlen und Befehle stimmen |
| Results: Lese-Schritt-Tabelle und Feststellung | geprüft, ohne Befund über F-4 und F-5 hinaus: 17 Zeilen, Zähler, Ausgänge gegen die `state.md` |
| Results: Folge-Slices und Paarungen (a) bis (c) | geprüft, ohne Befund: vier Dateien in `open/`, Anker auffindbar, `BEO-PGC/`-Kennungen mit `evidence/` (Ausnahme benannt) |
| Results: Sprache (`AGENTS.md` §3.7) | geprüft, ohne Befund über F-9 hinaus: Ist-Zustand, indikativ; keine Vorher-/Nachher-Formulierung über verworfene Alternativen, Herkunfts-Anker in der erlaubten Form |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zugesagter Folgeschritt ohne Träger · Gemeldete Ungenauigkeit ohne Träger ·
Spec-Lücke als aufgelöst formuliert, Entscheidung ausstehend · Streichungs-Begründung mehrdeutig, Restgrenze fehlt ·
Kandidat mit Adresse Lese-Schritt ohne benannten Ausgang · Lerneintrag im Sammelabsatz · Benannte Lücke ohne Adresse ·
Zusammenfassung enger als der Träger · Stand-Aussage im Record

## Verdikt

**Merge-blockierend:** ja — formal, wegen F-1 (MEDIUM): ein in sechs Records zugesagter Folgeschritt (Validator
zu [`LH-FA-CAP-009`](../../spec/lastenheft.md) am E2E-Beleg) ist in vier Notizen und in der Results-Notiz weder
erfüllt noch als „entfällt ausdrücklich“ geführt. Die inhaltliche Substanz der Closure trägt: keine der elf
Slice-Notizen und die Results-Notiz sind Floskeln, alle Pflicht-Inhalte (a), (b), (c) sind konkret, die
nachgemessenen Zahlen der Results-Notiz stimmen. F-2 bis F-6 (LOW) und F-7 bis F-9 (INFO) blockieren nicht. Die
Welle ist bereits gepusht; „blockierend“ heißt hier: die Closure gilt erst nach dem Nachtrag zu F-1 als auditierbar
geschlossen, eine Nacharbeit an Slices oder Code ist nicht verlangt.

**Übergabe:** F-1, F-3, F-4, F-5 und F-8 gehen an den Planner (Ergebnis-/Closure-Notiz; F-5 zusätzlich an den
Architect für die ausdrückliche Entscheidung der drei Kandidaten und die `state.md`); F-2 an den Planner (Träger:
Beobachtungs-Eintrag oder Folge-Slice-Zeile, Kommentar-Korrektur an Code und `.proto` ist Sache eines Implementer-Zuges
mit `make proto-generate`); F-6, F-7 und F-9 zur Kenntnis. Die **Finding-Klassen** gehen in den Zähler des Steering-Loops;
keine erreicht mit diesem Lauf 3×. Dieser Report ist ein **Lauf-Beleg**, ändert keine Closure-Notiz und ersetzt weder
Verifikation noch Validierung.
