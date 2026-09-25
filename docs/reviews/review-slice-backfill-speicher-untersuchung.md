# Review-Report: slice-backfill-speicher-untersuchung — 2026-09-25

**Review-Art:** Code — der Diff liefert eine Messreihe zum Speicher des Feed-Containers bei
Backfill-Runs (zwei Bench-Skripte, `BENCH_FEED_ENV`/`BENCH_FEED_DOCKER_ARGS`, Vertrag
`harness/targets/bench-backfill.md`, `harness/README.md`), zwei Messberichte, Handbuch §9/§4 (Version
1.60), den Folge-Slice `slice-retention-lauf-speicher-begrenzung` und den Plan-Nachzug samt
Suchlauf-Feld; geprüft gegen Plan, ADRs und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten).
Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-speicher-untersuchung` (wellenlos), Diff-Range
`493a28ad..4f94f900` (HEAD `4f94f900`). Slice-Inhalt sind `fb04c86b` (Werkzeuge), `b9cf9729`
(Messberichte), `7459afbb` (Handbuch), `5a141059` (Folge-Slice), `28768e7e` und `4f94f900`
(Plan-Nachzug) sowie die zwei reinen `git mv`-Commits `446e3ca6` und `37e825e2` (0 Einfügungen und
Löschungen); `ec7f76e2` (`slice-sdk-kotlin-cloudsmith`, nur Plan) und die Architect-Verdikt-Commits
im Range sind nicht Slice-Inhalt und nicht geprüft.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Kommentar-Klassen, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-speicher-untersuchung` (§1 Ziel und Abgrenzung, §3 Plan samt
  Suchlauf-Feld, §4 Trigger, §6 Risiken) und der Folge-Slice `slice-retention-lauf-speicher-begrenzung`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b),
  [`ADR-0014`](../plan/adr/0014-retention-domain-policy.md)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md), [`LH-FA-RET-002`](../../spec/lastenheft.md),
  [`LH-FA-RET-004`](../../spec/lastenheft.md); [`SPEC-005`](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-sdk-origin.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **Ursache am Code nachgelesen (Stand `4f94f900`):** `RunRetentionService.Run`
  (`internal/application/usecase/retention/service.go:63`) ruft `ReadChanges` mit
  `ChangeQuery{Source: command.Source}` ohne Limit auf, danach eine Schleife über alle Records;
  `PostgresChangeStoreAdapter.ReadChanges` übergibt `query.Limit` (`*int`, nicht gesetzt = `NULL`)
  als `LIMIT $6` (`postgresstorage/store.go:153–163`); `SelectChanges` liefert `old_data` und
  `new_data` (`queries/queries.go:47–70`); `sqlexec.ReadChanges` liest alle Zeilen in eine Liste
  (`sqlexec/translate.go:25–65`); `retentionInterval` = 10 s (`bootstrap/wiring.go:174`), der Takt
  läuft in `runRetentionCleanup` (`wiring.go:1212`). Es gibt kein Paging und kein Limit, das der
  Bericht übersähe; `ReadChanges` hat außer der Retention nur den Use Case `readchanges` und die
  HTTP-Schicht als Aufrufer. Der Retention-Lauf steht im Tag `v0.1.2`
  (`git show v0.1.2:internal/application/usecase/retention/service.go`, Zeile 63 gleich) — der
  Defekt ist vorbestehend und nicht auf den Backfill beschränkt; `git diff --stat
  493a28ad..4f94f900 -- internal cmd` ist leer (kein Produktionscode im Diff).
- **Kernmessung selbst nachgefahren** (`make image`-Stand des Implementers, Feed-Image-ID
  `sha256:969b7fad…cf4e` in der gedruckten Kopfzeile; je ein Lauf je Bedingung, `tools/bench-backfill-memory.sh`
  nach Vertrag, `BENCH_MEM_STAGES=100000`, `BENCH_FEED_ENV=GODEBUG=gctrace=1`, schmal, ein schwerer
  Docker-Lauf zugleich, `free -m` vorher 14,8 GB verfügbar):
  - *Default-Image, drei Runs* (Lauf `20260925T160407Z`, Exit 0): `memory.peak` nach dem Nachlauf
    135,4 / 242,2 / 352,7 MiB bei 100.000 / 200.000 / 300.000 Changes in `cdc.change` (1,39 / 1,24 /
    1,16 KiB je Change, abgeleitet; Reihe A des Berichts: 133,8 / 245,0 / 449,2 MiB); Spitze des
    Prozesses im Run 1 bei leerem `cdc.change` 10,1 MiB `anon`; `Bereinigung gelaufen 14`.
  - *Image ohne Bereinigung `pg-change-feed:bench-noret`* (Lauf `20260925T161254Z`, Exit 0): bei
    denselben 100.000 / 200.000 / 300.000 Changes `memory.peak` 13,0 / 13,4 / 12,3 MiB, `anon` im
    Nachlauf konstant 8,9 bis 9,1 MiB, gedruckt „Bereinigung gelaufen 0, fehlgeschlagen 0“ —
    Schalter 1 des Berichts reproduziert; die Bereinigung erzeugt den Speicher.
  - *`--memory 64m`* (Lauf `20260925T162100Z`, `BENCH_MEM_RUNS=1`): Exit 1, gedruckt „Feed-Container
    nicht mehr lesbar: Status exited, Exit 137, OOMKilled true“ im Nachlauf; nach dem Lauf weder
    Container noch Netz `pgc-bench-*` (bestätigt die Zeilen L1/L2 und den Exit-Code, den der
    Bericht nicht druckt).
  - *Beobachtung im Run bei nicht leerem `cdc.change`* (dieselbe Default-Messung): Run 2 mit
    100.000 Changes davor `anon` Spitze im Run 123,4 MiB, Run 3 mit 200.000 davor 237,9 MiB (Proben
    bei 25 % bzw. 50 % der Zeilen) — siehe F-3.
- **Nachrechnungen gegen das Zeilen-Dokument** (Zeilen extrahiert, dann gerechnet): elf Zahlen der
  Tabellen aus den gedruckten Zeilen bestätigt (`memory.peak` der Reihen A, B, C, D, H, J, K; die
  19 Einheiten der Grundlinie, `anon` 5,0 bis 5,5 MiB und `memory.current` 5,8 bis 6,8 bzw. 18,9 bis
  22,4 MiB; Reihe E/J `anon`-Spitze 8,9 bis 10,6 MiB, Nachlauf 8,0 bis 10,0 MiB; Reihe D
  55,7 / 84,1 / 82,9 % nach 120 s; Reihe N Tabelle 3.8; Richtgrößen-Rechnung 4.000.000 × 1,03/1,57/
  2,65/4,19 KiB = 3,9/6,0/10,1/16,0 GiB; 4.000.000 × 2 KiB + 64 MiB = 7,7 GiB). Abweichungen: F-1, F-2.
- **Suchlauf des Plans** an beiden Ständen nachgefahren (`git grep … 493a28ad` und `git grep … 4f94f900`):
  Befehl 1 Parent 58/13/36/111/24 = 242, Diff 71/23/36/111/39 = 280; Befehl 2 11 und 11; die alten
  Zahlen 7 und 2, „nicht untersucht“ 1 und 0 — jede Zahl des Feldes stimmt. Zusätzlich gesucht:
  `git grep -n -i -E 'Arbeitsspeicher|\bRAM\b|memory|OOM|Speicherlimit'` über `README.md`, `docs/user`,
  `spec`, `harness/README.md`, `compose.yaml`, `examples`, `sdks/*/README.md` — keine Träger außer
  Handbuch und `spec/pflichtenheft.md` Zeilen 72 und 389 (Capture-Transaktionspuffer).
- **Werkzeuge:** `bash -n tools/bench-backfill-memory.sh` ohne Meldung; `git diff --check` meldet
  nur zwei Zeilen mit Leerzeichen am Zeilenende im Zeilen-Dokument (gedruckte Ausgabe, Zeilen 451 und
  466); `docker stats --no-stream` gegen einen nur angelegten Container druckt `0B / 0B` (Grundlage
  von F-4); Datei-Modus des neuen Skripts 100755, `tools/bench-backfill.sh` bleibt 100644 wie am Parent.
- **Umgebung:** dangling Volumes (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34
  nach allen Läufen; kein `prune`; keine eigenen Container oder Netze zurückgeblieben;
  `tools/schema/plan.yaml` und `tools/schema/down.sql` nach den Läufen per `git checkout`
  zurückgenommen.
- **Nicht gefahren (Grenze):** die Reihen E bis J und N (Stunden Laufzeit; die Werte stammen aus dem
  Zeilen-Dokument, nachgerechnet, nicht nachgemessen); `make bench` als Ganzes (endet am Messhost
  bekannt rot); der Neubau der Varianten-Images (`bench-noret` und Geschwister sind lokal, nicht
  committet — ihr Inhalt ist nur über die gedruckte Zeile „Bereinigung gelaufen 0“ belegt).
- **Nebenläufe anderer Rollen:** während des Reviews standen im Index und Arbeitsbaum fremde
  Änderungen (ein neuer ADR-Entwurf und ein Architect-Verdikt); der Report-Commit nimmt
  ausschließlich diese Datei auf.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Handbuch-Satz „Ein Wert von `GOGC=25` … senkt die Spitze um etwa 18 bis 21 % (`n` = 3)“ und der Bericht („18 bis 21 % weniger“, „nur um 18 bis 21 %“) tragen die Messung nicht: Reihe K gegen Reihe A ergibt 224,5 gegen 273,9 MiB (−18,0 %), **415,5 gegen 468,6 MiB (−11,3 %)** und 477,4 gegen 605,9 MiB (−21,2 %); der mittlere der drei Werte fehlt in der Spanne (die Rechnung ergibt 11 bis 21 %). | `AGENTS.md` §3.12 Instanz A · Reviewer-Skill „Zahl im Träger gegen die Messung driftend“ | `docs/user/benutzerhandbuch.md:1706-1708`, `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md:220-224` und `:299-300` | ja — die drei Paare aus dem Zeilen-Dokument (Reihe A Stufe 200.000, Reihe K) dividieren | Zahl im Träger gegen die Messung driftend |
| F-2 | HIGH | Die abgeleiteten Werte „je Change“ der beiden Runs 3 (2.000.000 Changes vor dem Run, Spitzen 3.096,6 und 2.751,7 MiB) stehen mit 1,51 und 1,35 KiB in der Tabelle; die Rechnung Spitze durch 2.000.000 ergibt **1,585 und 1,409 KiB**. Damit ist der Höchstwert der 15 Werte 1,59 statt der in Bericht (§1, §6) und Handbuch („1,03 bis 1,57 KiB“, „oberhalb der gemessenen Höchstwerte von 1,57“) genannten 1,57. Zusätzlich steht die Spitze von Run 3 im Handbuch in der Zeile „2.000.000“ unter der Kopfzeile „Zahl der Changes, die **nach dem Run** in `cdc.change` stehen“ — nach Run 3 stehen dort 3.000.000. | `AGENTS.md` §3.12 Instanz A · Reviewer-Skill „Zahl im Träger gegen die Messung driftend“ | `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md:154`, `:159`; `docs/user/benutzerhandbuch.md:1690-1701` und `:1710-1713` | ja — 3.096,6 × 1024 / 2.000.000 rechnen; die Zeile „Zeilen in cdc.change nach dem Nachlauf: 3000000“ der Reihe B, Run 3 lesen | Zahl im Träger gegen die Messung driftend |
| F-3 | MEDIUM | Der neue §4-Absatz sagt, der Speicherbedarf wachse mit der Zahl der Changes „nach dem Run, nicht während er läuft“; der Nachbar-Bullet in §9 und der Bericht (Ergebnis 2) beschränken die flache Spitze im Run auf ein **leeres** `cdc.change`. Gemessen liegt die Spitze im Run bei nicht leerem `cdc.change` hoch: Reihe B Run 2 `anon` 1.460,8 MiB im Run (1.000.000 Changes davor), Run 3 3.088,1 MiB; im Review-Lauf Run 2 123,4 MiB und Run 3 237,9 MiB (100.000 bzw. 200.000 davor). Zwei Aussagen im selben Träger, keine verweist auf die andere. | Reviewer-Skill „Nachzug widerspricht dem Nachbarn im selben Träger“ · `AGENTS.md` §3.12 Instanz B | `docs/user/benutzerhandbuch.md:473-478` gegenüber `:1676-1679` | ja — Zeilen der Reihe B, Run 2 und 3 („Kopierdauer … anon … Spitze“) | Nachzug widerspricht dem Nachbarn |
| F-4 | MEDIUM | Der Vertrag sagt für **jedes** Skript, das `bench::start_feed` ruft, „ein Kill durch die Grenze ist ein Messergebnis“. Getragen ist das nur in `bench-backfill-memory.sh` (Exit 1 mit `OOMKilled`, im Review-Lauf bestätigt). In `bench-backfill.sh` liefert `feed_mem_mib` für einen beendeten Container `0.0` (gemessen: `0B / 0B`), die Statusabfrage in `run_backfill` läuft dann bis `RUN_TIMEOUT_S` (1800 s), ohne den Kill zu nennen; `BENCH_FEED_DOCKER_ARGS` macht diesen Pfad mit diesem Diff erst erreichbar, und der Bericht führt die Eigenschaft in §8 als „vorhanden, nicht Teil dieses Slice“. | `harness/targets/bench-backfill.md` (Zusage) · Reviewer-Skill „Zusage ohne Bindung an ihren Pfad“ (MEDIUM: unklare Fehlerbehandlung am Rand) | `harness/targets/bench-backfill.md:88-92`, `tools/bench-backfill.sh:88-97` und `:150-175` | ja — `BENCH_FEED_DOCKER_ARGS="--memory 64m"` gegen `tools/bench-backfill.sh` | Zusage ohne Bindung an ihren Pfad |
| F-5 | MEDIUM | Der Trigger des Folge-Slice („muss `done` sein, bevor ein Server-Release veröffentlicht wird, der den Backfill trägt **und dessen Handbuch den Speicherbedarf nicht als bekannte Grenze führt**“) ist mit diesem Diff unerfüllbar: das Handbuch führt die Grenze jetzt. Der Messbericht (§7) empfiehlt dagegen „den Änderungs-Slice vor dem Release umsetzen“ und stellt die Entscheidung dem Nutzer. Zwei Aussagen über dieselbe Vorbedingung; der Trigger trägt die Empfehlung nicht. | `AGENTS.md` §3.12 Instanz B · Reviewer-Skill „Nachzug widerspricht dem Nachbarn“ | `docs/plan/planning/open/slice-retention-lauf-speicher-begrenzung.md:129-135` gegenüber `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md:345-356` | ja — beide Stellen lesen | Trigger trägt seine Empfehlung nicht |
| F-6 | LOW | `gc_summary` druckt „keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)“ immer, wenn das Fenster keine GC-Zeile enthält, auch bei gesetztem `GODEBUG`: gedruckt in Reihe B Run 3 (Feed-Umgebung `GODEBUG=gctrace=1`) und in jedem Nachlauf des Review-Laufs ohne Bereinigung. Die Zeile nennt eine Ursache, die der Code an dieser Stelle nicht prüft. | Reviewer-Skill „Kommentar-Klassen“ (Ausgabe-Text als Zusage) · Maintainability | `tools/bench-backfill-memory.sh:194` | ja — Lauf mit gesetztem `GODEBUG` und ohne Bereinigung | Ausgabetext nennt eine ungeprüfte Ursache |
| F-7 | LOW | Zwei Handbuch-Sätze tragen Chronik bzw. Planung statt Ist-Zustand: „Die **früher** übernommene Zahl von 1.544 MiB … ist damit vereinbar“ und „Die Bereinigung … zu begrenzen ist Gegenstand einer eigenen Änderung; **bis dahin** gilt die Bemessung oben“ (ohne Adresse der Änderung). | `AGENTS.md` §3.7 · Reviewer-Skill „Zustandsfeld/Träger trägt Chronik“ | `docs/user/benutzerhandbuch.md:1717-1725` | nein — Lese-Handlung | Vorher-Nachher-Sprache in Handbuch-Prosa |
| F-8 | LOW | Das Risiko „Die Abfrage selbst wächst mit der Zahl der Changes“ nennt als Beobachtung „bei 3.000.000 Zeilen … ausbleibende Bereinigungs-Takte“ und als Ursache die Sortierung in PostgreSQL; der Bericht beobachtet das Ausbleiben ab **2.000.000** Changes vor dem Run (Reihen B und C, Run 3) und nennt die Ursache „nicht belegt“ (§3.3, §4). Zahl und Ursache stehen im Plan fester als in der Quelle. | `AGENTS.md` §3.12 Instanz B | `docs/plan/planning/open/slice-retention-lauf-speicher-begrenzung.md:156-159` | ja — Messbericht §3.3 lesen | Ursache im Plan fester als in der Messung |
| F-9 | LOW | Das Handbuch gibt die Aussage „gleich, ob ein Backfill oder die laufende Erfassung“ die Changes geschrieben hat, ohne Ursprung; der Bericht führt sie unter „Ursache (belegt)“ (Ergebnis 1) und benennt erst in §9, dass die Live-Erfassung als Quelle „aus dem Code gelesen, nicht gemessen“ ist. Der Code trägt die Aussage (Retention liest alle Changes der Quelle), gemessen ist sie nicht. | `AGENTS.md` §3.12 Instanz B | `docs/user/benutzerhandbuch.md:1670-1672`, `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md:22-24` und `:385-386` | nein — Messung mit Live-Last fehlt | Aussage ohne Ursprung im Träger |
| F-10 | LOW | Zwei Fehlerpfad-Reste in `bench-backfill-memory.sh`: (a) endet `run_one` mit `return 1` (Run `failed`/`interrupted`/Zeitgrenze), bleibt die mit `mktemp` angelegte Zeitreihen-Datei stehen (die Falle löscht nur die Fenster-Dateien); (b) die Meldung „Zeile 58: …/memory.current: Datei oder Verzeichnis nicht gefunden“ erscheint trotz `2>/dev/null` in der Ausgabe (gedruckt in Reihe L1, Zeile 449 des Zeilen-Dokuments). | Maintainability | `tools/bench-backfill-memory.sh:25`, `:147`, `:163-168`, `:58` | ja — Lauf mit `--memory 64m` (b), Lauf mit Run-Abbruch (a) | Fehlerpfad-Rest im Werkzeug |
| F-11 | INFO | Reihe N (ein Feed-Container über alle Stufen, `--memory 10g`) fährt Run 3 der Stufe 1.000.000 mit 2.330.000 Changes davor und zeigt 60 s nach dem Run 2.448,4 MiB (`anon` nicht auf 60 MiB gefallen); die Reihen B und C (frischer Container je Run) zeigen bei 2.000.000 davor ein Ausbleiben der Takte mit `anon` 62 bzw. 55 MiB. Der Bericht deutet das Ausbleiben als „nicht belegt“, nennt die Reihe N nicht als Gegenprobe und den Unterschied (frischer Container je Run gegen einen Container über alle Stufen) nicht. | Maintainability | `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md:171-179` gegenüber `:242-257` | ja — Zeilen der Reihen B, C und N | Gegenprobe zum ungeklärten Befund nicht genannt |
| F-12 | INFO | Der Retention-Lauf mit unbegrenzter Lesung steht im Tag `v0.1.2`; die Empfehlung (§7) und der Folge-Slice-Trigger rahmen die Entscheidung als Frage des Backfill-Releases, das Mindestalter von 24 Stunden lässt bei laufender Erfassung dieselbe Menge wachsen. Der Bericht nennt den Tag nicht. Zuständig: Planner (Priorisierung des Folge-Slice). | `AGENTS.md` §3.13 (Träger außerhalb des Diffs) | `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md:345-356` | ja — `git show v0.1.2:internal/application/usecase/retention/service.go` | Vorbestehender Defekt im Release-Rahmen nicht benannt |
| F-13 | INFO | Zwei Einschränkungen des Werkzeugs stehen im Bericht, nicht im Vertrag: der cgroup-Pfad ist auf `system.slice/docker-<id>.scope` (cgroup v2, systemd-Treiber) festgelegt — auf einem Host mit anderem Treiber endet das Skript sauber mit Exit 1, der Vertrag nennt die Vorbedingung nicht; und die Exit-Codes der Mutationen (§8) stehen nicht in den gedruckten Zeilen (der Review-Lauf bestätigt Exit 1 für `--memory 64m`). | Maintainability | `harness/targets/bench-backfill.md:60-62` und `:120`, `tools/bench-backfill-memory.sh:135` | ja — Skript auf einem cgroupfs-Host | Werkzeug-Vorbedingung nicht im Vertrag |
| F-14 | INFO | Die vom Implementer gemeldeten Eigenverstöße (`sed -i` und ein Host-Python-Heredoc auf der eigenen neuen Datei `tools/bench-backfill-memory.sh`) haben keinen Rückstand im Repository: die Datei parst (`bash -n`), trägt kein CR und keine Nicht-ASCII-Steuerzeichen, `git status` war vor den Review-Läufen sauber, kein Python-Artefakt im Diff. Der Befund ist ein Prozessverstoß gegen `AGENTS.md` §3.1 ohne Wirkung auf das Ergebnis. | `AGENTS.md` §3.1 | `tools/bench-backfill-memory.sh` | nein — Prozess, kein Gate | Docker-only-Verstoß ohne Rückstand |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Ursache des Speichers (Code: `retention/service.go`, `postgresstorage/store.go`, `queries/queries.go`, `sqlexec/translate.go`, `bootstrap/wiring.go`) | geprüft, ohne Befund: die Beschreibung des Berichts trifft den Code (kein Limit, kein Paging, beide Row Images, Takt 10 s, Kommentar an `retentionInterval` wahr); Schalter „Bereinigung aus“ am Review-Lauf reproduziert (13,0/13,4/12,3 MiB gegen 135,4/242,2/352,7 MiB bei gleichen 100.000/200.000/300.000 Changes, „Bereinigung gelaufen 0“ gegen 14) — plausibler Beleg für „notwendig“ und, mit dem Code, für „hinreichend“ |
| Vorbestand und Produktionscode | geprüft, ohne Befund: `git diff --stat 493a28ad..4f94f900 -- internal cmd` leer; die Varianten-Images sind nicht committet; der Defekt ist vorbestehend (`v0.1.2` trägt den Retention-Lauf), nicht durch diesen Diff eingeführt |
| `tools/bench-backfill-memory.sh` (Exit-Codes, Falle, Docker-Hygiene, cgroup-Auslesung) | geprüft, ohne Befund über F-6, F-10, F-13 hinaus: `set -euo pipefail`, Falle auf `EXIT` und `TERM` (`bench::cleanup` entfernt Container mit `-v` und das Netz; nach drei Läufen und einem Exit-1-Lauf kein `pgc-bench-*` übrig, dangling Volumes 34 vor und nach), `kill -TERM $$` aus dem Unter-Shell-Pfad erreicht die Falle (gemessen: Exit 1 nach OOM), Guard „Zähler nicht lesbar“ meldet und endet Exit 1 (Reihe O), kein `eval`, keine Host-Toolchain |
| `tools/bench-lib.sh` (`BENCH_FEED_ENV`, `BENCH_FEED_DOCKER_ARGS`) | geprüft, ohne Befund: `-e "$entry"` je Eintrag als eigenes Argument, `${BENCH_FEED_DOCKER_ARGS:-}` bewusst unquotiert (Wortteilung gewollt, im Kommentar benannt), kein `eval`, keine Injektionsfläche über die Umgebung des Aufrufers hinaus; `"${extra[@]}"` bei leerem Array unter `set -u` in Bash 5.2 ohne Fehler (gelaufen) |
| `tools/bench-backfill.sh` (Grundlinie, Zeilen in `cdc.change`, 60-s-Probe) | geprüft, ohne Befund über F-4 hinaus: drei Ergänzungen, Reihe N zeigt die neuen Felder (Grundlinie 6,1 MiB, 0 Changes; „60 s nach dem letzten Lauf“); Messung der Kopierdauer unverändert |
| `harness/targets/bench-backfill.md`, `harness/README.md` | geprüft, ohne Befund über F-4, F-13 hinaus: Ist-Zustand ohne Chronik, Ablauf und Variablen-Tabelle stimmen mit dem Skript überein (Defaults, Fenster 10/60/120 s, Exit-1-Bedingung), die Zeile `make bench` nennt das zweite Skript als Werkzeug außerhalb von `make bench`; Makefile-Kommentare unverändert |
| Handbuch §9 „Grenzwerte“, §4 „Bestand als Backfill überführen“, „Aufbewahrung (Retention)“ | geprüft, ohne Befund über F-1, F-2, F-3, F-7, F-9 hinaus: `Version:` 1.59 → 1.60 mit Zeile 1.60 in `### Änderungshistorie` (HIGH-Klasse „Handbuch-Versionshistorie“ nicht ausgelöst); Ursache, `--memory 64m` → `OOMKilled`/Exit 137 (Reihen L1, L2 und Review-Lauf), Bemessung „2 KiB je Change (schmal), 4,5 KiB (breit), plus 64 MiB“ ist ausdrücklich „abgeleitet … kein Nachweis einer Obergrenze“ (2 KiB liegt über dem korrigierten Höchstwert 1,59, 4,5 über 4,19); 7,7 GiB für 4.000.000 Zeilen nachgerechnet; keine neue Betreiber-Oberfläche; Anwender-Sprache trägt (Reihenbezeichnungen und `memory.peak` stehen im Bericht verlinkt) |
| Messberichte (Zahlen mit Ursprung, Reihen A–O, Selbstkonsistenz) | geprüft, ohne Befund über F-1, F-2, F-11 hinaus: jede Zahl nennt Reihe und Lauf, „übernommen“ trägt nur 1.544 MiB (Lauf `20260924T233628Z`) und ist so gekennzeichnet; 11 Tabellenzahlen gegen das Zeilen-Dokument nachgerechnet (Liste oben), die Reihen-Zeilen tragen den Wortlaut der Ausgabe (zwei Auslassungen benannt, `git diff --check` meldet nur zwei Zeilenende-Leerzeichen der gedruckten Ausgabe); Grenzen (§9) und Bestand (§10, 34 dangling Volumes, gegen den Review-Lauf gehalten) offen benannt |
| Folge-Slice `slice-retention-lauf-speicher-begrenzung` (Ziel, DoD, Umfang, Architect-Frage) | geprüft, ohne Befund über F-5, F-8 hinaus: Ziel und Abgrenzung tragen (Backfill-Blockgrenze ausdrücklich nicht die Ursache, Richtgröße erst nach der Messung nach der Änderung, andere `ReadChanges`-Aufrufer benannt), DoD-Punkt 1 bindet die Seitengrenze an ihre Eingabe (Seitengröße 1 gegen unbegrenzt), Rückführung „Vertrag von `ChangeStorePort` verlangt eine ADR“ benannt; §3.13-Tabelle bewusst leer („Implementer trägt ein“, Slice in `open/`) |
| Suchlauf-Feld des Plans (§3), beide Stände | geprüft, ohne Befund: alle 14 gemessenen Zahlen (Befehl 1 je Wurzel, Befehl 2, alte Zahlen, „nicht untersucht“) an `493a28ad` und `4f94f900` bestätigt; die Zeile „Nicht gefunden“ trägt (keine weiteren Träger für Speicherzahlen in `spec`, `docs/user`, `README.md`, SDK-READMEs, `compose.yaml`); der Register-Träger `state.md` wird als gemeldet geführt und trägt die Zeile „die Ursache ist nicht untersucht“ (Planner-Closure) |
| Kommentare (`AGENTS.md` §3.7, diff-skopiert) | geprüft, ohne Befund über F-6, F-7 hinaus: Kommentare in `bench-backfill-memory.sh` und `bench-lib.sh` tragen Zusage, Kopplung oder Abgrenzung im Indikativ, keine Slice-/Wellen-Nummer, kein Konjunktiv über verworfene Alternativen |
| Commits und Hard Rules | geprüft, ohne Befund über F-14 hinaus: alle Slice-Commits nennen `LH-*`/`ADR-*`, kein Betreff trägt `SPEC-*`/`ARC-*`, die zwei `git mv`-Commits sind rein (Rename ohne Inhalt), kein host-lokaler absoluter Pfad in Skript, Vertrag, Handbuch, Berichten (Suche nach Wurzel-Segmenten leer), keine Suppression, `git status` nach dem Review bis auf fremde Nebenläufe sauber |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 5 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Zahl im Träger gegen die Messung driftend · Nachzug widerspricht
dem Nachbarn · Zusage ohne Bindung an ihren Pfad · Trigger trägt seine Empfehlung nicht ·
Ausgabetext nennt eine ungeprüfte Ursache · Vorher-Nachher-Sprache in Handbuch-Prosa · Ursache im
Plan fester als in der Messung · Aussage ohne Ursprung im Träger · Fehlerpfad-Rest im Werkzeug ·
Gegenprobe zum ungeklärten Befund nicht genannt · Vorbestehender Defekt im Release-Rahmen nicht
benannt · Werkzeug-Vorbedingung nicht im Vertrag · Docker-only-Verstoß ohne Rückstand

## Verdikt

**Merge-blockierend:** ja — die zwei HIGH-Findings (F-1, F-2) sind Zahlen in einem Handbuch-Träger und
in den Berichten, die gegen die eigene Messung driften; sie sind an den gedruckten Zeilen
nachrechenbar (11,3 % statt 18 %-Untergrenze, 1,585 statt 1,51 KiB je Change, Höchstwert 1,59 statt
1,57 KiB). Die Substanz trägt: die Ursache (unbegrenzte Lesung des Retention-Laufs) ist am Code
gelesen und am Review-Lauf reproduziert (Bereinigung an gegen aus bei gleicher Zahl der Changes:
135,4/242,2/352,7 gegen 13,0/13,4/12,3 MiB), die Bemessung „2 KiB je Change plus 64 MiB“ liegt auch
über dem korrigierten Höchstwert, kein Produktionscode ist berührt, der Suchlauf stimmt an beiden
Ständen. F-3 bis F-5 sind Widersprüche zwischen Nachbar-Aussagen (Handbuch §4 gegen §9, Trigger
gegen Empfehlung) und sollten in derselben Fixrunde gezogen werden. Die DoD-Zeile „Review
durchgeführt“ bleibt offen, bis die Fixrunde läuft.

**Übergabe:** Findings gehen an den Implementer (F-1 bis F-5 Pflicht, F-6 bis F-10 nach Ermessen);
F-11 bis F-13 an den Planner bzw. Architect (Gegenprobe und Priorisierung des Folge-Slice, Frage an
den Architect zum `ChangeStorePort`-Vertrag läuft dort), F-14 ist ein Prozess-Hinweis. Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report
selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt); der
Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
