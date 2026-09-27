# Review-Report: slice-capture-leerlauf-quellbelege — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec-Stellen und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-capture-leerlauf-quellbelege` (wellenlos), Diff-Range `7367c483..HEAD` (12 Commits,
15 Dateien, +1424/−321). Implementer-Lauf: drei Lifecycle-Commits `a45d64c7`, `a827f497`, `dbc4dbe4`; `332199ba`
(Keepalive-Test `internal/adapters/driving/replication/receive/sourcekeepalive_test.go`, Anpassung
`tools/harness/run-replication-tests.sh`); `6ab20878` (Runner-Phase „Fehlerschwelle beendet den Container“ in
`tools/harness/run-integration-tests.sh`, regeneriertes `docs/user/e2e-abdeckung.md`); `cf39c0ed`
(`harness/README.md`); `3f688677` (Plan-Nachzug, Befund, Suchlauf-Feld, zwei Verweis-Umstellungen). Danach die
Planner-Nachzüge `348226cc` (Plan, DoD 2) und `92af01f5` (`harness/README.md`, ein Satzteil). Nicht Gegenstand, nur
als Kontext gelesen: der Architect-Commit `7305b578` und die Planner-Commits `068267db`, `ff31937a`. Produktionscode
ist im Diff nicht berührt (Test, Runner-Skripte, Doku).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen
ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Das Edit-Werkzeug stand im Lauf
nicht zur Verfügung; alle Mutationen liefen als `sed … Datei > Kopie` im Scratchpad (Ausgabe nach stdout), die Kopie
wurde im Testlauf über die Repo-Datei in den read-only eingehängten Baum gelegt (`-v Kopie:Ziel:ro`); nie `sed -i`,
nie Host-Interpreter, nie eine Umleitung auf eine Repo-Datei. Kein Aufruf einer Host-Toolchain (Go lief im
gepinnten Toolchain-Image).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-capture-leerlauf-quellbelege` (§1, §2 DoD als Bezug, §3 Plan mit Befund und Suchlauf-Feld, §6
  Risiken)
- [`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md) (Festlegung 1 und 2,
  Re-Evaluierungs-Trigger „Ein committeter Test der Quellseite entsteht“),
  [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md),
  [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-QA-REL-001`](../../spec/lastenheft.md), [`LH-QA-REL-003`](../../spec/lastenheft.md),
  [`LH-FA-CAP-009`](../../spec/lastenheft.md); [`SPEC-012`](../../spec/pflichtenheft.md),
  [`SPEC-013`](../../spec/pflichtenheft.md)
- Architect-Verdikt
  [`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
  (Kontext: die Klasse `storage` statt `replication` der Phase ist Codefehler, Träger
  `slice-wal-fehlerschwelle-ausgangsklasse`; die Phase führt die Klasse bewusst nicht als Zusage)
- `AGENTS.md` (§3.1, §3.2, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md` (`MR-000` bis `MR-003`)
- Vorherige Findings am Modul: [`review-slice-transformationen-start-reihenfolge`](review-slice-transformationen-start-reihenfolge.md),
  [`review-slice-backfill-slot-leerlauf-bestaetigung`](review-slice-backfill-slot-leerlauf-bestaetigung.md)

**Eigene Messungen:**

- **Läufe** (Exit-Codes ungefiltert gesichert): `make kommentar-kennungen DIFF=7367c483` Exit 0, kein Kandidat;
  `make fmt-check` „260 Go-Dateien geprüft, alle formatiert“; `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md` „10 Zeilen stimmen“ (Exit 0);
  `make commit-traceability RANGE=7367c483..HEAD` OK für 12 Commits, keine Struktur-ID im Betreff; die zwei
  Lifecycle-Moves `a45d64c7` und `dbc4dbe4` sind reine Renames (0 Zeilen, `git show -M --numstat`);
  `make gates` vor dem Commit dieses Reports Exit 0 (Log im Scratchpad).
- **Keepalive-Test, Grundlauf** (gepinnte Digests aus `run-replication-tests.sh` und `.github/workflows/e2e.yml`,
  je eigener Wegwerf-Container mit Standard-`wal_sender_timeout` 60 s, beide Läufe gleichzeitig auf demselben
  Docker-Host, Aufruf wie im Runner: `go test -count=1 -v -run '^TestSourceKeepalive…$'`). Gedruckt, PostgreSQL
  18.6: `Keepalive inmitten der Transaktion nach 11538 von 400000 Änderungen, ServerWALEnd 0/53CFB40 gleich
  Commit-LSN 0/53CFB40, bestätigt 0/53CFB40; Neustart ab 0/53CFB40 lieferte 400000 Änderungen mit Commit-LSN
  0/53CFB40`, `--- PASS (40.52s)`; PostgreSQL 17.11: `… nach 11075 von 400000 …, ServerWALEnd 0/513C5E0 gleich
  Commit-LSN 0/513C5E0 …`, `--- PASS (40.42s)`. Ohne `CDC_SOURCE_KEEPALIVE_TEST_DSN`: `--- SKIP`, Exit 0.
- **Runner-Phase:** ein voller `make test-integration` am Arbeitsbaum (Image `:dev` vom 2026-09-26 23:41, nach dem
  letzten Produktionscode-Commit `990608d3`; Exit 0). Die Ausgabezeile der Phase, gedruckt: „… lief der
  Feed-Container über eine Nullprobe von 12 s weiter; ein Schreiber auf die nicht aktivierte Tabelle
  feed_e2e_wal_stop_foreign erzeugte 16033424 B WAL, der Container endete 2 s danach (höchstens 90 s gewartet) mit
  Ausgang 1, die Abbruch-Zeile nannte einen Rückstand von 16035800 B (Fehlerschwelle 8388608 B),
  cdc.process_heartbeat trug die Klasse storage; nach dem Neustart ohne Konfigurationsdatei war die wartende
  Transaktion über cdc.changes lesbar“. Die Klasse `storage` ist damit im sechsten Lauf bestätigt (der Lauf des
  Reviewers ist der sechste; die fünf Läufe des Implementers sind **übernommen**, von mir nicht wiederholt).
  Mutationen an der Phase habe ich **nicht** gefahren (ein schwerer Lauf, siehe Verdikt); die Angabe zur
  `stopStream`-Mutation im Plan ist weiter **übernommen**.
- **Erzeugnis:** `docs/user/e2e-abdeckung.md` aus einer Erzeuger-Kopie im Scratchpad neu erzeugt (Funktionsblock
  und alle 39 `abdeckung_declare`-Zeilen des Runners mit unveränderten Zeilennummern, Go-Hälfte über
  `go test -run '^TestAbdeckungstabelleZeilen$'`): **byte-gleich** zur committeten Datei (`cmp`). Der volle Lauf
  meldet „E2E-Abdeckungstabelle unverändert“. Neue Zeile 68 trägt die Kennungen `LH-QA-REL-001`, `LH-QA-REL-003`;
  die 38 verschobenen Ort-Zeilen (`git diff -U0`: 39 hinzugefügt, 38 entfernt) sind die Folge von zwei eingefügten
  Kommentarzeilen im Kopf des Runners.
  `-run`-Abgleich: der Diff berührt weder `test/` noch ein `-run`-Muster des Runners (`git diff … | grep -c '^[+-].*-run'`
  gibt 0), es entsteht keine `func TestE2E*`; die Zählung bleibt die des Parent.
- **Mutationen am Keepalive-Test** (sieben Mutationen, acht Läufe, je an einer Scratchpad-Kopie; Eingabeseite des
  Tests, dann die Zusage-Seite): siehe Tabelle. Docker-Ausgangsstand 36 dangling Volumes, Endstand 36 (kein `prune`,
  kein neues Volume); die vier Wegwerf-PostgreSQL-Container (mit `docker rm -fv`) und das Netz sind entfernt.

| # | Ort in `sourcekeepalive_test.go` | Mutation | Ergebnis |
|---|---|---|---|
| M1 | `:235` | bestätigte Position `ServerWALEnd + 1 GiB` (Mutation der Plan-DoD) | rot an PostgreSQL 18.6 und 17.11: „Neustart des Streams lieferte die Transaktion nicht bis zum COMMIT (0 Änderungen gelesen)“ nach 40 s |
| M5 | `:235` | bestätigte Position `ServerWALEnd + 1` (ein Byte hinter der Commit-LSN) | rot an 18.6, dieselbe Meldung: die Gleichheit trägt, nicht erst „hinter dem WAL-Ende“ |
| M2 | `:236` | Bestätigung des Clients gestrichen | rot an 17.11: „confirmed_flush_lsn … erreicht die bestätigte Position nicht“ |
| M3 | `:30` | Stille ab START_REPLICATION 5 s statt 38 s | rot an 18.6 nach 7 s: „COMMIT nach 400000 Änderungen ohne Keepalive zwischen BEGIN und COMMIT“ |
| M7 | `:24` | 2000 statt 400000 Änderungen | rot an 18.6: „COMMIT nach 2000 Änderungen ohne Keepalive“ |
| M9 | `:229` | Behauptung `ServerWALEnd == BEGIN-Kopf` auf `+ 8` gedreht | rot an 17.11 an der Behauptung |
| M10 | `:254` | `START_REPLICATION` des Neustarts ab `confirmed_flush_lsn + 1` | rot an 17.11 (Neustart liefert nichts) |

## Findings

### F-1 — DoD 1 beschreibt einen Aufbau, den der Test nicht hat; §3 trägt den Umbau, die DoD nicht

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“; `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md:74-88` (DoD 1) gegen `:165`
  (§3-Zeile „Plan-Nachzug: ein **roher Protokoll-Client** statt des Stream-Adapters“) und
  `internal/adapters/driving/replication/receive/sourcekeepalive_test.go:201-266`
- `befund`: DoD 1 sagt: „eine offene Transaktion hält die Stream-Sitzung, währenddessen committet eine zweite
  Transaktion eine große Änderungsmenge …, ein Keepalive tritt inmitten der **ersten** Transaktion auf“. Der Test hat
  keine erste und keine zweite Transaktion: ein roher Client liest 38 s nichts, eine einzige Transaktion über 400.000
  Zeilen committet, der Keepalive steht zwischen BEGIN und COMMIT **dieser** Transaktion; der Stream-Adapter ist nicht
  beteiligt. Die §3-Zeile trägt den Umbau als „Plan-Nachzug“, der Wortlaut der DoD blieb, und die Verifikation liest
  die DoD. Die gemessene Aussage (Position des Keepalive = Commit-LSN, Neustart liefert vollständig) ist dieselbe wie in
  `ADR-0121` §Gemessen; die Aufbauform ist eine andere (dort hielt ein Stand-in des Adapters die Transaktion offen).
- `verifizierbar`: ja — Lesen von DoD 1 gegen `sourcekeepalive_test.go:213-223`
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-2 — DoD-Haken 1 steht auf `[x]`, Lauf und Mutation der Zeile „Zu belegen durch“ stehen in keinem committeten Träger

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B (ein Kriterium, das erst nach der Arbeit belegt werden kann, ist eine Zusage,
  bis der Anker steht)
- `pfad`: `docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md:74-88`; §6 (`:265-287`, alle Ausgänge
  „(bei Closure)“)
- `befund`: Die DoD verlangt je eine gedruckte Zeile mit Position und Änderungszahl an PostgreSQL 18 und 17 und die
  Mutation „Position `+ 1 GiB`“ rot. Der Plan trägt weder eine Zeile noch einen Lauf noch die Mutation; die Läufe
  stehen nur im Bericht des Implementers. Ich habe beide Läufe und die Mutation nachgefahren (Messungen oben, M1
  an beiden Versionen rot) — der Haken stimmt, sein Anker fehlt im Träger. Die Risiken „zeitabhängig“ und „17 wie 18“
  (§6) tragen „mehrere Läufe je Version ohne Ausfall“ nur in der Erwartungs-Form.
- `verifizierbar`: ja — Lesen von §2 und §6; `make test-replication` bzw. `PG_TEST_IMAGE` mit dem 17-Digest
- `klasse`: DoD-Haken ohne committeten Beleg-Anker

### F-3 — Suchlauf-Feld: der zweite Stand ist `diff`, zwei Zahlenreihen stehen nebeneinander, ein Muster ohne Grund für den Suchraum

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (§Suchform: jede Einschränkung mit Grund im Feld; Stände als Commit-Kennung),
  `harness/sensors/suchlauf-nachmessen.md` (`diff` bewegt sich mit jedem Commit)
- `pfad`: `docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md:219-230` (Block), `:234-235` (Tabelle)
- `befund`: Nachgemessen (`make suchlauf-nachmessen` Exit 0; zusätzlich `git grep` an Commit-Ständen, Plan-Datei
  ausgeschlossen): am Stand der Übergabe `3f688677` 33 · 47 · 15 · 7 — das sind die Zahlen der Tabellenzeilen; am
  Stand `92af01f5` 33 · 56 · 24 · 11 — das sind die `diff`-Zeilen des Blocks. Beide Reihen stimmen und sind
  redlich erklärt (Klammer in der Tabelle, Absatz vor dem Block). Der Block bindet den zweiten Stand als `diff`,
  also als Arbeitsbaum: er hat sich seit der Übergabe einmal bewegt (drei Zeilen vom Planner neu gemessen) und bewegt
  sich mit jedem weiteren Commit, der einen Träger der drei Muster anlegt — die ADR-Ergänzung des Architects (DoD 3)
  nennt „Quellseite“ und „inmitten einer Transaktion“ und färbt `diff 33` dann rot (hergeleitet, nicht gemessen).
  Das Vorbild `slice-transformationen-start-reihenfolge` bindet beide Stände als Commit-Kennung; Commit-Stände liegen
  hier vor (`3f688677`, `92af01f5`). Muster 5 (`nur zum Zug|Fehler bei Stream-Ende|regulär endete`) ist auf
  `-- internal` beschränkt, ohne dass das Feld den Grund für diese Einschränkung nennt (der Absatz danach lässt
  Muster 4 die übrigen Bäume tragen).
- `verifizierbar`: ja — `make suchlauf-nachmessen PLAN=<Plan-Datei>` nach dem nächsten Commit, der einen Träger anlegt
- `klasse`: Suchlauf-Feld bindet beweglichen Stand

### F-4 — Kommentare der Runner-Phase: eine Allaussage ohne Beleg, die Grenze ohne Rang-Zeiger, die Kopplung an die Phase davor nur teilweise genannt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Klassen Grenze und Kopplung), §3.12 Instanz B; Skill-HIGH „Kommentar trägt keine der
  Kommentar-Klassen“ (Unterpunkt Zusage — hier nicht erreicht, weil der Kommentar kein zugesagtes Verhalten des
  Skripts, sondern eine Aussage über das Produkt macht)
- `pfad`: `tools/harness/run-integration-tests.sh:3468-3469` („der Rückstand erreicht die Fehlerschwelle nur, wenn
  der Slot nichts bestätigt“), `:3544-3545` („seine Klasse steht in der Ausgabe des Runners und ist nicht Teil der
  Zusage dieser Phase“), `:3468` („Die Gegenseite der Phase oben“, `:3476-3477`), `:42-44` (Kopf)
- `befund`: (a) Der Satz „nur, wenn der Slot nichts bestätigt“ ist eine Allaussage; der Rückstand ist die Differenz
  aus dem WAL-Ende und `confirmed_flush_lsn` (`internal/adapters/driving/replication/receive/walretention.go`,
  `Measure`), er übersteigt die Schwelle immer dann, wenn die Bestätigung dem WAL-Ende nicht folgt; die Phase belegt
  einen Fall („die Sperre hält die Persistierung, die Bestätigung bleibt stehen“) und trägt keinen Anker für „nur“.
  (b) Die Grenze der Klasse steht im Skript ohne Zeiger auf die Norm (`ADR-0049`) oder den Träger; die Zeiger stehen
  im Plan und in `harness/README.md`, nicht am Ort der Grenze (der Träger-Slice führt diesen Kommentar in seiner
  Änderungs-Tabelle, `slice-wal-fehlerschwelle-ausgangsklasse.md:176`). (c) Die Phase nutzt aus der Phase davor
  `bf_wal_hold`, `WAL_WARN_BYTES`, `WAL_ERROR_BYTES`, `WAL_WAIT_SECONDS` und den Zeitstempel `wal_feed_started`
  (setzt sie neu); der Kommentar sagt „wie oben“ nur zur Konfigurationsdatei — ein Herausnehmen oder Umsortieren der
  Phase davor bricht diese hier (mit `set -u` laut, nicht still). (d) Der Kopf `:42-44` trägt jetzt zwei Kennungen
  (`ADR-0120`, `ADR-0049`) in einem Satz; das Werkzeug `make kommentar-kennungen` liest nur Go-Dateien, der Kopf ist
  Bestand mit vielen Kennungen.
- `verifizierbar`: nein (Lese-Handlung; (c) ja durch Umsortieren der Phasen)
- `klasse`: Kommentar-Aussage breiter als der Beleg / Grenze ohne Rang-Zeiger

### F-5 — Träger außerhalb des Diffs: ein Begründungssatz, den `ADR-0121` ersetzt hat, und die Handbuch-Zusage der Klasse

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger außerhalb des Diffs), Register `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:467-469` (Godoc von `confirmIdle`);
  `docs/user/benutzerhandbuch.md:466`
- `befund`: (a) `confirmIdle` begründet die Bedingung „keine offene Transaktion“ mit „das WAL-Ende liegt dann hinter
  Nachrichten, die noch nicht gespeichert sind“ — das ist der Begründungssatz, den `ADR-0121` Festlegung 1 ausdrücklich
  ersetzt („`ServerWALEnd` liegt dann hinter Nachrichten …“), und der neue Test bindet das Gegenteil an beiden
  Versionen (ServerWALEnd = Commit-LSN der offenen Transaktion). Der Satz ist Bestand, der Slice ändert keinen
  Produktionscode (§1); die Muster des Suchlauf-Feldes treffen ihn nicht (kein Muster enthält „hinter Nachrichten“
  oder „noch nicht gespeichert“). (b) Das Handbuch nennt „mit der Klasse `replication` (Ausgang 1)“; die Klasse ist
  gemessen `storage` (Lauf des Reviewers). Das ist im Plan (§3 Befund, Tabellenzeile 2), in `harness/README.md` und im
  Architect-Verdikt benannt und dem Träger `slice-wal-fehlerschwelle-ausgangsklasse` (Codefehler, das Handbuch bleibt
  danach wahr) übergeben; bis dahin ist die Zusage für die gehaltene Persistierung falsch.
- `verifizierbar`: ja — `git grep -n 'hinter Nachrichten, die noch nicht' -- internal`
- `klasse`: Träger-Nachzug (außerhalb des Diffs)

### F-6 — Was der Keepalive-Test bindet, was nicht, und drei Kleinigkeiten

- `kategorie`: INFO
- `quelle`: Maintainability; Register `BEO-PGC/test-runner-stiller-ausschluss`
- `pfad`: `internal/adapters/driving/replication/receive/sourcekeepalive_test.go:201-266`,
  `tools/harness/run-replication-tests.sh:205-226`
- `befund`: Gebunden (Mutationen oben): der Keepalive liegt inmitten der Transaktion (M3, M7), die Position ist die
  Commit-LSN (M9 direkt, M5 und M1 über die Lieferung), der Slot nimmt sie an (M2), der Neustart setzt an
  `confirmed_flush_lsn` an und liefert vollständig (M10, M1, M5). Nicht gebunden: das Verhalten des Adapters (das
  trägt `TestRunNoConfirmationInsideOpenTransaction` im Unit-Tier), mehr als ein Keepalive je Transaktion,
  `proto_version` 2 mit Streaming großer Transaktionen (Re-Evaluierungs-Trigger von `ADR-0121`), eine gleichzeitige
  zweite Quelltransaktion. Die vom Implementer genannte Lücke „ServerWALEnd == Commit-LSN nicht mutiert“ ist an der
  Behauptung geschlossen (M9) und hat eine zweite, unabhängige Bindung (M5: liegt die Position auch nur ein Byte
  hinter der Commit-LSN, liefert der Neustart nichts); mutiert werden kann nur der Test, nicht die Quelle. Kein stiller
  Ausschluss: ohne `CDC_SOURCE_KEEPALIVE_TEST_DSN` überspringt der Test (gemessen), der Runner verlangt sein
  `--- PASS` im eigenen Lauf (aus dem Wächter hergeleitet: ein Umbenennen oder ein Überspringen liefert keine
  PASS-Zeile und färbt das Tier rot; nicht als Lauf erprobt). Zeit: 38 s Stille ab START_REPLICATION, Abbruch der
  Quelle bei 60 s, rechnerisch 22 s Puffer (aus den Konstanten hergeleitet); gemessen liefen zwei gleichzeitige Läufe
  (PostgreSQL 18.6 und 17.11) je in rund 40 s grün, der Keepalive lag nach etwa 11.000 von 400.000 Änderungen. Das
  Flake-Risiko besteht auf einem so langsamen Runner, dass 400.000 Zeilen mehr als 38 s brauchen (dann liest der
  Client den COMMIT ohne Keepalive und der Test meldet das laut), nicht aus dem 22-s-Puffer. Das Tier
  `make test-replication` läuft nicht in `ci.yml`/`e2e.yml`; der Beleg lebt vom manuellen Lauf. Kleinigkeiten: die
  Meldung `:240` druckt `confirmed_flush_lsn` als `%x` ohne `0/`-Präfix, die bestätigte Position als `%s` (dieselbe
  Zeile, zwei Formen); die Skip-Meldung („laufen über make test-replication“) gilt für diesen Test nur im eigenen Lauf
  (der tier-weite `go test ./...` überspringt ihn ebenfalls); der Test-Godoc beschreibt die Mutation `+ 1 GiB` im
  Indikativ als Grenze der Aussage und trägt genau eine Kennung (`ADR-0121`).
- `verifizierbar`: ja — Mutationen M1 bis M10
- `klasse`: Bindung und Grenze des Tests benannt

### F-7 — Verweise der Planner-Nachzüge auf Dateien, die wandern

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13; derselbe Mechanismus, den `3f688677` am Register und in `slice-start-vorlauf-grenze`
  für diese Plan-Datei behoben hat
- `pfad`: `docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md:103`, `:201`, `:235`
  (Links auf `../open/slice-wal-fehlerschwelle-ausgangsklasse.md`);
  `docs/plan/planning/open/slice-wal-fehlerschwelle-ausgangsklasse.md:230` (Link auf
  `../in-progress/slice-capture-leerlauf-quellbelege.md`)
- `befund`: Beide Dateien tragen Lifecycle-Verzeichnisse im Linkziel; die eine geht nach `next/` (Kante in
  `welle-transformationen`), die andere nach `done/`; der jeweilige Move lässt den Link des anderen Dokuments
  brechen (`docs-check` Modul `links`). Nicht Gegenstand des Code-Reviews (Planner-Commits), nur gemeldet.
- `verifizierbar`: ja — `make docs-check` nach dem ersten Move
- `klasse`: Verweis auf wandernde Datei

## Negativbefunde

- geprüft, ohne Befund: **`tools/harness/run-replication-tests.sh`** (Schwerpunkt 1, Ablauf): der Lauf steht hinter
  dem Schwellen-Beleg (das WAL der Standard-Instanz, das der Keepalive-Test erzeugt, stört dessen Messung nicht), nutzt
  ausschließlich `docker run` mit dem gepinnten Toolchain-Image und dem DSN der Standard-Instanz, hat den `--- PASS`-Wächter
  mit Ausgabe bei Rot, und die Phase `measure` (Coverage) und der tier-weite Lauf setzen die neue Variable nicht
  (Überspringen, keine 38 s Wartezeit dort). Der Kommentar „drei letzte Läufe brauchen zusätzlich das `--- PASS`“
  stimmt mit dem Skript überein (Slot-Reserve, Schwellen, Keepalive).
- geprüft, ohne Befund: **Aufbau der Phase „Fehlerschwelle beendet den Container“** (Schwerpunkt 2): die Sperre
  (`ACCESS EXCLUSIVE` auf `cdc.change`, 180 s Haltedauer gegen höchstens etwa 150 s Phasenlänge in der
  ungünstigsten Reihenfolge, aus den Konstanten hergeleitet) steht, bevor die Transaktion eintrifft, und wird über
  `pg_locks` belegt; die wartende Persistierung wird in `pg_stat_activity` gelesen (Wartezustand `Lock`); die
  Nullprobe (12 s, mindestens zwei Prüf-Takte zu 5 s) bindet „die gehaltene Persistierung allein beendet nicht“ mit
  Startzeitpunkt-Vergleich; die Last (16033424 B) wird gegen die Schwelle gemessen, bevor gewartet wird (Abbruch mit
  Meldung, nicht ein grüner Lauf ohne Beleg); Ausgang, Zahl der Abbruch-Zeilen (`1`), Rückstand der Abbruch-Zeile über
  der Schwelle, „Lauf beendet mit Fehler“ und eine nicht leere Klasse in `cdc.process_heartbeat` sind je einzeln
  geprüft; die wartende Transaktion ist nach dem Neustart über `cdc.changes` gelesen (Persist-before-ACK am realen
  Prozess). Geteilter Zustand: die Sitzung endet mit `bf_hold_end` vor dem Neustart, das Temp-Verzeichnis wird
  entfernt, der Abbruch-Pfad räumt über den bestehenden `trap`; die nachfolgende Phase (Transformations-Rundläufe)
  startet gegen einen frisch erzeugten Container ohne Override; der volle Lauf lief nach der Phase grün durch. Die
  Reihenfolge (zwischen Leerlauf-Bestätigung und Transformations-Rundläufen) und die `abdeckung_declare`-Zeile
  (Anker in der Ausgabezeile) sind konsistent, das Erzeugnis byte-gleich.
- geprüft, ohne Befund: **Klasse `storage` als benannte Grenze** (Schwerpunkt 2 und 7): die Phase schreibt die Klasse
  nicht als Zusage fest (nur „nicht leer“), die Ausgabezeile nennt den gemessenen Wert, `harness/README.md` nennt die
  Abweichung von `ADR-0049` mit Träger, der Plan (DoD 2) hält „Grenze bleibt“ als zulässigen Ausgang. Die Herkunft im
  Code stimmt mit der Messung überein: `mergeStreamAndWALFaultOutcome` (`internal/bootstrap/wiring.go:1178`) gibt einen
  Stream-Fehler jeder Klasse vor dem WAL-Fehler zurück, und `stopStream` bricht den Kontext des Streams ab.
  Beurteilung der Entscheidung des Implementers: siehe Verdikt.
- geprüft, ohne Befund: **`harness/README.md`** (Schwerpunkt 3): die Zeilen `make test-replication` und
  `make test-integration` nennen Test, Aufbau, Beleg und Grenze und stimmen mit Test und Skript überein (Instanz,
  Variable, Digest-Wahl über `PG_TEST_IMAGE`, 400.000 Änderungen, Nullprobe 12 s, Ausgang 1, Klasse `storage` mit
  Träger); die Planner-Änderung `92af01f5` („Codefehler, Träger `slice-wal-fehlerschwelle-ausgangsklasse`“ statt
  „der Befund liegt beim Architect“) ist konsistent mit Verdikt, Plan (DoD 2, §3) und dem Träger-Slice. Die beiden
  umgestellten Verweise in `3f688677` (`beleg-nur-als-einmalige-reviewer-messung/state.md`,
  `slice-start-vorlauf-grenze.md`: Link → Kennung) sind inhaltlich unverändert. Die Träger in fremden Dateien
  (Handbuch, `wiring.go`-Kommentar an `mergeStreamAndWALFaultOutcome`,
  `walretention_slotgrowth_internal_test.go`) sind mit Adresse und Frist (Closure des Träger-Slice) gemeldet, nicht
  mitgeändert.
- geprüft, ohne Befund: **Zahlen und Ursprung** (§3.12, Schwerpunkt 4): die zehn Zeilen des Suchlauf-Blocks stimmen
  (Ausnahme der Form: F-3); die Zahlen der Tabellenzeilen („Übergabe“, 33 · 47 · 15 · 7) habe ich am Stand `3f688677`
  mit `git grep` nachgezählt; die Angaben „fünf Läufe“ und „zwei Sekunden“ im Befund-Absatz von §3 stehen als
  gemessen, mit den Läufen und der Ausgabezeile der Phase als Quelle; die Angabe in DoD 2 kennzeichnet sich als
  **übernommen**; mein Lauf bestätigt „2 s“ und `storage`. Die im Auftrag genannten Einzelwerte (123/125/372 s,
  40,3 s, 85,30 %, LSN-Werte, „16033448 B“) stehen in keinem committeten Träger des Diffs (Bericht des Implementers);
  mein Lauf nennt 16033424 B, eine Messung je Lauf, kein Widerspruch. DoD-Haken gegen den Ist-Zustand: DoD 1 (F-2),
  DoD 6 (Suchlauf-Feld vorhanden) und DoD 8 (`harness/README.md` nachgezogen, Handbuch unberührt) treffen zu; DoD 2
  bewusst offen, ebenso 3, 4, 5 und die Closure-Punkte.
- geprüft, ohne Befund: **Kommentare §3.7** (Schwerpunkt 5) in allen geänderten Zeilen (Go, Shell):
  `make kommentar-kennungen DIFF=7367c483` ohne Kandidat, `make fmt-check` ohne Abweichung; ein Herkunftsfeld je
  Go-Kommentar, keine Kette, keine Kompaktform, kein „ff.“; Suche in den hinzugefügten Zeilen von `*.go` und `*.sh`
  nach `würde`, `wäre`, `hätte`, `früher`, `zuvor`, `vorher`, `ursprüngl`, `bisher`, `jetzt`, `nicht mehr`,
  `nachträgl`, `statt`, `sonst`, `Vorfassung`, `zunächst`, `slice-`, `welle-`, `seit slice`: ein Treffer
  (`nicht mehr als die Fehlerschwelle` in einer Fehlermeldung, kein Kommentar). Einzige Kommentar-Findings: F-4.
- geprüft, ohne Befund: **Docker-only/§3.1, Traceability, Moves** (Schwerpunkt 6): die Skript-Änderungen rufen
  ausschließlich `docker`, `docker compose`, `psql` im Container und die Runner-Hilfen auf (`mktemp`, `printf`, `cat`
  in ein Temp-Verzeichnis außerhalb des Repos wie in der Phase davor); keine Host-Toolchain, kein `sed -i`, kein
  Interpreter, keine Umleitung auf eine Repo-Datei; das Erzeugnis wird vom bestehenden Schreib-Weg des Runners
  geschrieben (byte-gleich nachgewiesen). Alle 12 Betreffe tragen `LH-*`/`ADR-*`, keine `SPEC-*`/`ARC-*` im Betreff;
  Moves `a45d64c7` und `dbc4dbe4` rein, `a827f497` (Verantwortlich) eigener Commit; kein `//nolint` (§3.2).
- geprüft, ohne Befund: **`make test`-Auswirkung und Netz**: der Test überspringt ohne DSN (gemessen), die Datei
  kompiliert mit den vorhandenen Abhängigkeiten (`pglogrepl`, `pgconn`, `pgproto3`), `make test` bleibt netzlos.
- geprüft, ohne Befund: **Spec-Stratum, Zwei-Quellen-Drift, Handbuch**: `spec/`, `docs/plan/adr/` und
  `docs/user/benutzerhandbuch.md` nicht berührt; keine neue Betreiber-Oberfläche (keine `CDC_*`-Variable der
  Produktion, keine SQL-Funktion, kein Endpunkt; `CDC_SOURCE_KEEPALIVE_TEST_DSN` ist eine Test-Variable), also weder
  Handbuch-Zug noch Versionshistorie fällig. Produktionscode nicht berührt, kein Eingriff in die kritische Sektion
  (Persist-before-ACK).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Nachzug widerspricht dem Nachbarn im selben Träger · DoD-Haken ohne committeten
Beleg-Anker · Suchlauf-Feld bindet beweglichen Stand · Kommentar-Aussage breiter als der Beleg / Grenze ohne
Rang-Zeiger · Träger-Nachzug (außerhalb des Diffs) · Bindung und Grenze des Tests benannt · Verweis auf wandernde Datei

## Verdikt

**Merge-blockierend:** nein. Kein HIGH; der Keepalive-Test bindet an beiden PostgreSQL-Versionen (sieben
Mutationen, alle rot, drei an der Eingabeseite der Bestätigung), die Runner-Phase lief im vollen Lauf grün und ihr
Erzeugnis ist byte-gleich. Keine Fixrunde am Implementer: F-1 bis F-3 sind Plan-Text (Planner bzw. Plan-Nachzug),
F-4 ist Kommentar-Text in einem Skript, dessen Phase der Träger-Slice
`slice-wal-fehlerschwelle-ausgangsklasse` ohnehin anfasst (dessen Änderungs-Tabelle führt den Kommentar der Phase),
F-5 bis F-7 sind Träger außerhalb des Diffs. Die DoD-Zeile „Review durchgeführt“ ist im selben Commit wie dieser
Report auf `[x]` gezogen (Skill §DoD-Checkbox-Nachzug ohne Fixrunde).

**Frage 7 des Auftrags.** *Ist der Aufbau der Phase stabil genug für einen Beleg?* Ja, für einen Beleg mit der
benannten Grenze: fünf von fünf Läufen des Implementers (**übernommen**) und mein sechster (voller Lauf, Ende 2 s nach
der Last) tragen ihn; die Zeitfenster haben Reserve (Nullprobe 12 s gegen den 5-s-Takt, Wartegrenze 90 s gegen 2 s
gemessen, Last rund das Doppelte der Schwelle und vor dem Warten gemessen), und ein Fehlschlag der Vorbedingungen
endet als benannter Abbruch, nicht als grüner Lauf ohne Beleg. Was die Phase nicht trägt, ist eine Mutation der
Eingabeseite des Aufbaus (etwa „Sperre nicht gesetzt“): die habe ich nicht gefahren; aus dem Aufbau folgt, dass die
Phase ohne Sperre der Phase davor gleicht (Leerlauf-Bestätigung entlastet, der Container läuft weiter) und in die
90-s-Grenze läuft — hergeleitet, nicht erprobt. *War die Entscheidung des Implementers, die Klasse nicht
festzuschreiben, richtig?* Ja. Eine Prüfung auf `storage` schriebe einen Codefehler als Zusage fest (das Verdikt des
Architects nennt ihn Codefehler), eine Prüfung auf `replication` wäre rot; die Phase druckt die gemessene Klasse und
prüft „nicht leer“. Die Schärfung zur Zusage trägt der Träger-Slice (DoD 1 dort). Offen bleibt bis dahin, dass
nichts die gedruckte Klasse erzwingt — das ist die benannte Grenze, kein Mangel.

**Übergabe:** F-1 bis F-3 und F-7 an den Planner (Plan-Text, Suchlauf-Stände als Commit-Kennung, Links auf
wandernde Dateien). F-2 zusätzlich an den Verifier: Lauf und Mutation der DoD 1 habe ich nachgefahren (Zeilen und
Tabelle oben), der Verifier trägt sie in seinen Report. F-4 an den Träger-Slice
`slice-wal-fehlerschwelle-ausgangsklasse` (Kommentar der Phase, Frist: dessen Implementierung). F-5 an den Planner
(Träger außerhalb des Diffs, Frist: Closure des Träger-Slice; der `confirmIdle`-Satz ist neu benannt und hat noch
keine Adresse). F-6 ohne erwartete Aktion. Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den
Zähler. Der Report ist ein Lauf-Beleg und ersetzt keine Verifikation; die Läufe des Reviewers (zwei PostgreSQL-
Versionen, ein voller `make test-integration`) sind Nachmessungen, keine DoD-Abnahme.
