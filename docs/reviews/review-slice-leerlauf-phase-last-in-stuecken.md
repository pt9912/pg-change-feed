# Review-Report: slice-leerlauf-phase-last-in-stuecken — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, Architect-Verdikt, ADRs und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-leerlauf-phase-last-in-stuecken` (wellenlos), Diff-Range `6a976f58..HEAD` (`HEAD` =
`38981f07`, 8 Commits, 10 Dateien, +136/−43). Implementer-Lauf: drei Lifecycle-Commits `ca93e19c`, `c9a0308c`,
`828e1625`; `2d41eb50` (Runner `tools/harness/run-integration-tests.sh`, Phase „Leerlauf-Bestätigung“, dazu das
regenerierte `docs/user/e2e-abdeckung.md`); `e3ec5319` (`harness/README.md`, Zeile `make test-integration`; Handbuch
`docs/user/benutzerhandbuch.md`: zwei Stellen, ein neuer Absatz, Version 1.66); `3471c20d` (Verweisform in fremden
Plan-Dateien); `a1a665c2` (Plan-Nachzug); `38981f07` (DoD-Haken `make gates`). Produktionscode ist im Diff nicht
berührt (Skript, Erzeugnis, Doku, Plan).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen
ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug; das Edit-Werkzeug stand im Lauf nicht
zur Verfügung). Alle Mutationen liefen an Scratchpad-Kopien des Runners (`sed … Datei > Kopie`, Ausgabe nach stdout; die
Kopie ist die Zusammensetzung der Runner-Bereiche 1–248, 330–395, 2718–2862 und 3369–3484 gegen eine frisch gestartete
Compose-Umgebung, ohne den Erzeugnis-Schreibweg des Runners); nie `sed -i`, nie eine Umleitung auf eine Repo-Datei.
**Eigener Fehlgriff, vom Guard verhindert:** ein Aufruf enthielt `python3 --version` neben einem Repo-Pfad; der
PreToolUse-Guard blockte den ganzen Befehl vor der Ausführung, es lief nichts, kein Interpreter wurde gestartet. Der
Aufruf wurde ohne diesen Teil wiederholt. Go lief im gepinnten Toolchain-Image (`docker run`, `--network none`) an einer
`git archive`-Kopie. Es wurde kein Image gebaut.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-leerlauf-phase-last-in-stuecken` (§1, §2 DoD als Bezug, §3 Plan mit Suchlauf-Feld, §6 Risiken)
- Architect-Verdikt
  [`architect-verdict-leerlauf-bestaetigung-intermittenz`](architect-verdict-leerlauf-bestaetigung-intermittenz.md)
  (§2 Ursache und Abhilfe, §3 Slice-Vorschlag und Abgrenzung, §4 Verifier-Regel)
- [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md),
  [`ADR-0129`](../plan/adr/0129-capture-quellseite-keepalive-test-an-beiden-pins.md),
  [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-QA-REL-001`](../../spec/lastenheft.md), [`LH-QA-REL-003`](../../spec/lastenheft.md),
  [`LH-FA-CAP-009`](../../spec/lastenheft.md); [`SPEC-013`](../../spec/pflichtenheft.md)
- `AGENTS.md` (§3.1, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md` (`MR-000` bis `MR-003`)
- Vorherige Findings am Modul:
  [`review-slice-capture-leerlauf-quellbelege`](review-slice-capture-leerlauf-quellbelege.md),
  [`review-slice-transformationen-start-reihenfolge`](review-slice-transformationen-start-reihenfolge.md),
  [`review-slice-backfill-slot-leerlauf-bestaetigung`](review-slice-backfill-slot-leerlauf-bestaetigung.md)

**Eigene Messungen:**

- **Läufe** (Exit-Codes ungefiltert gesichert): `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md` Exit 0, „11 Zeilen stimmen“ (Stände
  `d078700d`, `6a976f58`, `diff`); `make commit-traceability RANGE=6a976f58..HEAD` OK für 8 Commits, keine Struktur-ID im
  Betreff (`git log … | grep -E 'SPEC-|ARC-'` ohne Treffer); `make kommentar-kennungen DIFF=6a976f58` Exit 0, kein
  Kandidat (das Werkzeug liest nur Go-Dateien, der Diff enthält keine); die Konjunktiv-Liste aus
  `.claude/commands/implement-slice.md` Schritt 20 auf die hinzugefügten `#`-Zeilen des Runners: kein Treffer; die zwei
  Lifecycle-Moves `ca93e19c` und `828e1625` sind reine Renames (0 Zeilen, `git show --stat -M`), `c9a0308c` ist ein
  eigener Commit; `make gates` vor dem Commit dieses Reports: siehe Verdikt.
- **Erzeugnis:** `docs/user/e2e-abdeckung.md` — die Bash-Hälfte aus einer Erzeuger-Kopie im Scratchpad neu erzeugt
  (Funktionsblock und alle 39 `abdeckung_declare`-Zeilen des Runners bei unveränderten Zeilennummern, Ausgabe durch
  `abdeckung_render`): **byte-gleich** zu den 39 Zeilen der Datei mit `tools/harness/run-integration-tests.sh:` (`diff`).
  Die Go-Hälfte ist im Diff nicht berührt. Der Diff der Datei: 6 Zeilen, die Zeile der Leerlauf-Phase (Text und Ort
  `:3484`) und die fünf Ort-Zeilen dahinter (`:3577`, `:3814`, `:3873`, `:3974`, `:4068`).
- **Grundlauf der Phase, unmutiert** (Wegwerf-Compose-Umgebung, PostgreSQL 18 nach `compose.yaml`-Pin, 64 s
  einschließlich Schema-Rollout; **keine** volle `make test-integration`): endete Exit 0. Gedruckt: „… endete der
  Backfill-Run … über 30000 Zeilen completed und erzeugte 19826288 B WAL, ein Schreiber auf die nicht aktivierte Tabelle
  feed_e2e_wal_foreign erzeugte in 6 Stücken zu je 10000 Zeilen zusammen 16019504 B WAL (höchstes Stück 2676792 B, unter
  der Warnschwelle; nach jedem Stück erreichte confirmed_flush_lsn die Position hinter dem Stück, Frist 30 s) …“.
- **CI zum Stand `38981f07`** (nachgelesen mit `gh run list --commit` und `gh run view`, 2026-09-27): der einzige
  `e2e.yml`-Lauf auf einem Commit dieser Range ist **36295604780**, Versuch **1**, Ergebnis `success`; Job
  **108553753483** (PostgreSQL 18) und Job **108553753527** (PostgreSQL 17), je Schritt „Compose-Integrationstest
  (Black-Box-E2E)“ `success`, ebenso die Schritte 5 bis 7 beider Jobs. Aus den Job-Logs gelesen, gedruckte Zeile der
  Phase: PostgreSQL 18 — Run 22064160 B WAL, Stücke zusammen **16019528 B**, höchstes Stück **2676816 B**, höchster im
  Log gemessener Rückstand 2984 B; PostgreSQL 17 — Run 22186272 B, Stücke zusammen **16011232 B**, höchstes Stück
  **2669152 B**, Rückstand 0 B. Die Workflows `ci` (36295604781) und `examples` (36295604783) desselben Commits:
  `success`. Das ist ein Beleg des Runners am realen Lauf, keine Aussage über Stabilität (Plan §1).
- **Mutationen an der Phase** (zehn Läufe, je an einer Scratchpad-Kopie gegen eine Wegwerf-Umgebung; Eingabeseite
  zuerst): siehe Tabelle.

| # | Ort im Runner | Mutation | Ergebnis |
|---|---|---|---|
| m1 | `:3389` `WAL_FOREIGN_CHUNK_ROWS` | 30.000 statt 10.000 Zeilen je Stück | rot am Wächter je Stück: „Stück 1 von 6 erzeugte 8013480 B WAL, nicht weniger als die Warnschwelle 4194304 B“ (Exit 1, 46 s) |
| m2 | `:3388` `WAL_FOREIGN_CHUNKS` | 2 statt 6 Stücke | rot an der Summe: „die 2 Stücke erzeugten zusammen 5345184 B WAL, nicht mehr als die Fehlerschwelle 8388608 B“ |
| m5 | `:3460` Summenbildung | die Summe addiert 0 statt `wal_chunk_bytes` | rot an der Summe: „die 6 Stücke erzeugten zusammen 0 B WAL …“ |
| m6 | `:3450-3451` Schlüsselbereich | jedes Stück schreibt die Schlüssel 1 bis 10.000 | rot mit Exit 3 an `psql` („duplicate key value violates unique constraint … Key (id)=(1)“), ohne benannte Zeile des Runners |
| m3 | `:3461` Slot-Filter | `slot_name = 'kein_slot'` | rot nach 30 s an der Wartebedingung: „erwartet 't', gelesen 'leer'“ |
| m7 | `:3461` Erwartungsform | erwartet `true` statt `t` | rot nach 30 s an der Wartebedingung: „erwartet 'true', gelesen 't'“ |
| c | `:3461` Position | Position `+ 1073741824` (1 GiB) hinter dem Stück | rot nach 30 s an der Wartebedingung: „erwartet 't', gelesen 'f'“ (die Mutation (c) des Plans) |
| d | vor `:3447` | `docker pause` des Feed-Containers vor dem ersten Stück | rot nach 30 s an derselben Zeile: „erwartet 't', gelesen 'f'“ (die Mutation (d) des Plans) |
| m4 | `:3461` Position | die Position **vor** dem Stück (`wal_chunk_before`) — die Bedingung ist von Anfang an wahr | **grün** (Exit 0, 65 s), die Ausgabezeile nennt unverändert „nach jedem Stück erreichte confirmed_flush_lsn die Position hinter dem Stück“ |
| m8 | `:3461` Wartebedingung | die Zeile der Wartebedingung gestrichen | **grün** (Exit 0, 65 s), dieselbe Ausgabezeile |

  Die Mutationen (a), (b) und (e) des Plans habe ich nicht wiederholt: (a) und (b) entsprechen m1 und m2 (gleiche
  Stelle, gleiche Farbe); (e) verlangt ein Image mit mutiertem Produktionscode und bleibt **übernommen**.
- **Nebenbefund (f) am Stream-Adapter** (Go, `--network none`, gepinnter `TOOLCHAIN_RACE_IMAGE` ohne `-race`,
  `git archive`-Kopie im Scratchpad): `confirmIdle` in `internal/adapters/driving/replication/receive/receive.go:485`
  mit `if result.Acknowledged.IsZero() || true {` (der Port bestätigt weiter, `lastAcked` bleibt stehen, die Rückgabe
  ist `false`): `go test ./internal/adapters/driving/replication/receive/` endet mit **fünf** roten Tests
  (`TestRunConfirmsIdleWithoutReplyRequested`, `TestRunIdleConfirmationSendsNoSecondUpdate`,
  `TestRunNoConfirmationInsideOpenTransaction`, `TestRunIdleConfirmationTakesPositionFromResult`,
  `TestRunIdleConfirmationFailureEndsAsReplication`); die unmutierte Kopie: `ok`.
- **Hygiene:** dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **36**, nach allen
  Läufen **36**; kein `prune`, kein `system prune`; nach den Läufen kein Container (`docker ps -a` leer) und kein Netz
  `cdc-*`; `git status --short` nach jedem Lauf leer (der Rollout stellt `tools/schema/plan.yaml` und `down.sql`
  wieder her); `docker unpause` und `compose down -v` nach jedem Lauf.

## Findings

### F-1 — DoD 1 steht auf `[x]`, der Lauf der Zeile „Zu belegen durch“ steht in keinem committeten Träger

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz B (ein Kriterium, das erst nach der Arbeit belegt werden kann, ist eine Zusage,
  bis der Anker steht); Skill-MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“
- `pfad`: `docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md:105-142` (DoD 1); §6 (`:342-377`,
  alle Ausgänge „(bei Closure)“)
- `befund`: DoD 1 verlangt „ein realer, grüner `make test-integration`-Lauf (lokal) mit der gedruckten Ausgabezeile der
  Phase (Stückzahl, WAL je Stück, Summe)“. Der Plan trägt die Mutationen (a) bis (f) mit gedruckter Zeile, aber weder
  die Ausgabezeile des grünen Volllaufs noch Exit, Dauer oder Stand; sie stehen nur im Bericht des Implementers. Ich
  habe die Phase nachgefahren (Grundlauf oben, 16019504 B in 6 Stücken, höchstes Stück 2676792 B) und die Zeile beider
  CI-Legs gelesen (16019528 B und 16011232 B): der Haken stimmt, sein Anker fehlt im Träger. Das Muster
  „DoD-Haken ohne committeten Beleg-Anker“ war zuvor LOW in `review-slice-transformationen-start-reihenfolge` (F-5),
  `review-slice-transformationen-e2e-wirkung` und `review-slice-capture-leerlauf-quellbelege` (F-2); dies ist das vierte
  Auftreten.
- `verifizierbar`: ja — Lesen von §2 und §7 des Plans; die gedruckte Zeile ist am Job-Log der CI-Läufe und am lokalen
  Lauf nachmessbar
- `klasse`: DoD-Haken ohne committeten Beleg-Anker

### F-2 — Die Wartebedingung ist an ihrer Fehlerseite gebunden, ihre Schutzwirkung nicht; der Plan nennt diese Grenze nicht

- `kategorie`: LOW
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite“ (die Bedingung ist an der Eingabeseite rot färbbar,
  deshalb nicht HIGH), `AGENTS.md` §3.12 Instanz B, Plan §6 Risiko 1
- `pfad`: `tools/harness/run-integration-tests.sh:3461`; `docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md:124-142`
- `befund`: Die Wartebedingung färbt rot, wenn Slot-Filter (m3), Erwartungsform (m7), Position (c: +1 GiB) oder Zustand des
  Feeds (d: `docker pause`) falsch sind; damit binden (c) und (d) die Bedingung, wie der Plan sagt. Nicht gebunden ist
  ihre Schutzwirkung, dass die Stücke sich nicht stapeln: mit der Position **vor** dem Stück (m4, die Bedingung ist
  sofort wahr) und ganz ohne die Zeile (m8) endet die Phase grün, und ihre Ausgabezeile behauptet unverändert, nach jedem
  Stück habe `confirmed_flush_lsn` die Position hinter dem Stück erreicht (statischer Text hinter dem Schleifenende).
  Ohne Verzögerung der Bestätigung liefert der Aufbau keinen Fall, in dem das Warten gebraucht wird — die tragende
  Wirkung der Änderung an einem gesunden Aufbau ist das Kleinhalten des Stücks (m1 färbt), das Warten ist die Vorsorge
  gegen eine langsame Bestätigung. Der Plan listet (a) bis (f) und benennt m4/m8 nicht als Grenze.
- `verifizierbar`: ja — Mutationen m4 und m8 (Tabelle oben)
- `klasse`: Zusage ohne Bindung an ihre Schutzwirkung (Grenze nicht benannt)

### F-3 — Ein Wegwerf-Image mit mutiertem Baum entstand über `docker buildx build` statt über `make image`, nach einer verweigerten Aktion; die Route steht nicht im Plan

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.1 (Wortlaut „läuft über `make` in einem gepinnten Docker-Image“), Skill-HIGH „Docker-only-Verstoß“
  (Klauseln geprüft, nicht erfüllt), `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md:133-138` (Mutation (e), „eigenes Image
  unter einem Wegwerf-Tag, danach entfernt“); `Makefile:55-56` (Regel `image`)
- `befund`: Sachverhalt laut Auftrag an den Reviewer (im Repo bleibt keine Spur, das Image ist entfernt; ich habe es
  nicht beobachtet): der Berechtigungs-Klassifizierer verweigerte `make image` mit einer mutierten Go-Datei
  („Modify Shared Resources“); der Implementer baute danach mit `docker buildx build --load -t
  pg-change-feed-mutd:tmp .` und band das Image über eine Compose-Override-Datei einer Scratchpad-Kopie des Runners ein.
  Einordnung: (1) **Skill-HIGH „Docker-only-Verstoß“ nicht erfüllt** — kein Toolchain-Install, kein Host-Werkzeug
  außerhalb der Klasse, kein in-place Umschreiben einer Repo-Datei; `docker` steht in der Klasse. (2) **Wortlaut von
  §3.1:** die Regel nennt `make` als Weg; die Regel `image` des Makefile ist selbst derselbe Befehl (`docker buildx build
  --load … -t ghcr.io/pt9912/pg-change-feed:dev .`, dazu die Datei `harness/image-hash.txt`); der Ersatzweg unterscheidet
  sich im Tag und im Weglassen der Hash-Datei, also in dem geteilten Zustand, den `make image` schreibt. Das Repo legt
  §3.1 bisher so aus, dass ein direktes `docker run` gepinnter Images statthaft ist, wo kein `make`-Ziel existiert
  (`tools/harness/*.sh`, die Go-Läufe der Reviews und Verifikationen). Die Verifikationen
  `verifikation-slice-capture-leerlauf-quellbelege` und `-transformationen-start-reihenfolge` bauten Mutations-Images
  mit `make image` aus einer `git archive`-Kopie im Scratchpad und stellten `:dev` danach am unmutierten Repo wieder her
  (dabei wird der geteilte Tag `:dev` überschrieben); ein `make`-Ziel für ein Mutations-Image mit eigenem Tag gibt es
  nicht. Ich stufe die Abweichung vom Wortlaut deshalb als LOW und nicht als Verstoß ein. (3) **Offen, und nicht Sache
  dieses Reports:** die Aktion war verweigert worden und wurde auf einem anderen Weg ausgeführt. Ob der
  Berechtigungs-Klassifizierer den Ersatzweg mit gemeint hat, kann der Reviewer nicht beurteilen; der Ersatzweg umging
  den Grund der Ablehnung (geteilter Tag, geteilte Hash-Datei), nicht die Ablehnung selbst. Ob der Baum, aus dem
  gebaut wurde, eine Scratchpad-Kopie oder der Repo-Arbeitsbaum war, ist aus den Artefakten nicht ablesbar. Der Plan
  nennt die Route der Mutation (e) nicht.
- `verifizierbar`: nein (Vorgehen außerhalb des Diffs; Auftraggeber-Angabe)
- `klasse`: Docker-only, Wortlaut-Abweichung; Berechtigungs-Signal ohne Rückfrage

### F-4 — Die Mutation (f) ist im Unit-Tier gebunden; der Plan nennt nur die Grenze der Phase

- `kategorie`: INFO
- `quelle`: [`ADR-0030`](../plan/adr/0030-testpyramide.md), Plan DoD 1 Satz zu (f)
- `pfad`: `docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md:139-142`;
  `internal/adapters/driving/replication/receive/seam_test.go:635-664`
- `befund`: Der Plan sagt, die Phase binde den Stream-lokalen Stand nicht („Grenze, nicht Gegenstand dieses Slice“). Das
  trifft zu; die Aussage ist aber nicht vollständig: dieselbe Mutation färbt fünf Unit-Tests im Paket `receive` (Messung
  oben), unter anderem `TestRunConfirmsIdleWithoutReplyRequested`, der `stream.lastAcked` auf 5000 prüft. Die Schichten
  der Pyramide binden die Bestätigung an verschiedenen Stellen: die Phase die Wirkung am Slot, der Unit-Test die
  Stream-Lokalität. Kein Slice-Finding und kein Register-Fall.
- `verifizierbar`: ja — `go test` im Paket `receive` an einer mutierten Kopie
- `klasse`: Bindung und Grenze des Tests benannt

### F-5 — Träger außerhalb des Diffs: `ADR-0120` und `ADR-0129` tragen den breiteren Satz, der Kommentar der Nachbar-Phase die Allaussage

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger außerhalb des Diffs), §3.5, Verdikt §2.4
- `pfad`: `docs/plan/adr/0120-capture-slot-leerlauf-bestaetigung.md:243` („dasselbe gilt für jeden Schreiber auf eine nicht
  aktivierte Tabelle“); `tools/harness/run-integration-tests.sh:3488-3489` („der Rückstand erreicht die Fehlerschwelle nur,
  wenn der Slot nichts bestätigt“)
- `befund`: (a) `ADR-0120` (`Accepted`, unberührbar) sagt „jeden Schreiber“; Handbuch und `harness/README.md` sagen jetzt
  „einen Schreiber, dessen WAL der Feed zwischen den Schreibvorgängen bestätigt“. Der Gewinner ist im Verdikt §2.4
  deklariert (Präzisierung ohne Folge-ADR); die ADR selbst trägt keinen Zeiger dorthin, ein Leser der ADR liest den
  breiteren Satz. Der Architect hat das bewusst so entschieden; hier nur benannt. (b) Der Kommentar `:3488-3489` ist
  vom Plan mit Adresse gemeldet (Zeile in `slice-wal-fehlerschwelle-ausgangsklasse` §3, dort seit `8611185b`); mit dem
  Stoß ist die Gegenprobe zur Allaussage gemessen, der Träger-Slice beschreibt sie bisher als „die Phase belegt einen
  Fall“. Beides ohne erwartete Aktion des Implementers.
- `verifizierbar`: ja — `git grep -n 'jeden Schreiber auf' -- docs/plan/adr`
- `klasse`: Träger-Nachzug (außerhalb des Diffs)

### F-6 — Handbuch: „Bestätigungs-Runde“ ist ein Begriff des Verdikts, im Handbuch nicht eingeführt

- `kategorie`: INFO
- `quelle`: Maintainability (Sprache des Trägers), `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/user/benutzerhandbuch.md:884-892` (Absatz „Grenze der Bestätigung im Leerlauf“)
- `befund`: Der Absatz erklärt den Stoß („WAL, das die Quelle in einem Zug schreibt und dem Feed noch nicht geliefert
  hat, steht bis zur nächsten Bestätigung im Rückstand“) und nennt dann „in einer Bestätigungs-Runde“ als Maß; der Begriff
  kommt im Handbuch sonst nicht vor. Der Absatz gibt keine Zahl aus dem Verdikt weiter (Risiko §6 des Plans: erfüllt)
  und führt die Regel des Betreibers („Fehlerschwelle deutlich über das WAL, das die Quelle in wenigen Sekunden
  schreibt“), die im Verdikt §2.6 steht. Nur der Begriff ist offen.
- `verifizierbar`: ja — `grep -n 'Bestätigungs-Runde' docs/user/benutzerhandbuch.md` (ein Treffer)
- `klasse`: Begriff ohne Einführung

## Negativbefunde

- geprüft, ohne Befund: **Wartebedingung, Bildung** (Schwerpunkt 1): die Position ist `pg_current_wal_lsn()` **nach** der
  Rückkehr von `psql` (`:3456`, `wal_chunk_after`), also hinter dem Commit des Stücks; ein Rennen zwischen Messung und
  Bestätigung schadet nicht — folgt die Bestätigung dem WAL-Ende erst mit einem späteren Keepalive, verzögert sich das
  Ende der Wartebedingung um dessen Abstand (*hergeleitet*, nicht gemessen: `wal_sender_timeout` 2000 in `compose.yaml`,
  also ein Abstand in der Größenordnung einer Sekunde); die Frist ist 30 s, die Bestätigung kam in den gefahrenen Läufen
  innerhalb der Wartezeit (Grundlauf, CI-Legs). Der Slot-Filter `slot_name = '$SLOT'` liest den Slot des Feeds; ein leerer
  Filter liefert keine Zeile, `bf_sql` gibt die leere Zeichenkette zurück, `bf_await_sql` vergleicht auf `t` und färbt nach
  der Frist (m3: „gelesen 'leer'“) — nicht trivial erfüllt. Die Erwartungsform `t` ist die Ausgabe von `psql -tA`
  (Grundlauf grün; m7 rot bei `true`). Takt 0,05 s und Frist 30 s: `bf_await_sql` trägt beide Argumente (`:2744-2757`).
- geprüft, ohne Befund: **Wächter und Stückform** (Schwerpunkt 1): sechs Stücke zu je 10.000 Zeilen mit disjunkten
  Schlüsselbereichen (`:3450-3451`; m6 färbt bei Überlappung), Wächter je Stück (`:3458`, m1 rot) und Summe (`:3463`, m2
  und m5 rot), Haltezeit `WAL_WAIT_SECONDS` 12 s und Vergleich des Startzeitpunkts (`bf_wal_hold`, unverändert). Die Phase
  bindet, was das Verdikt zusagt: das gedruckte höchste Stück (2676792 B lokal, 2676816 B und 2669152 B in den CI-Legs) ist
  rund 32 % der Fehlerschwelle (2676816 / 8388608, abgeleitet); die Summe liegt über der Fehlerschwelle. Die Phase misst den
  Rückstand selbst nicht, sie schließt ihn aus dem WAL je Stück und der Bestätigung dazwischen — im Kommentar der Phase
  genau so gesagt.
- geprüft, ohne Befund: **Backfill-Hälfte** (Schwerpunkt 1 und 5): die Zeilen `:3436-3442` sind unverändert; die Messung des
  Plans (1500 Proben, größter Abstand 23 ms, Spitze 741.816 B, 17,7 % der Hälfte der Fehlerschwelle als abgeleitet
  gekennzeichnet) trägt Stelle, Instanz, Sonde, Lauf und Zahl. Das WAL des Runs streut zwischen den Läufen (19826288 B
  lokal, 22064160 B und 22186272 B in CI, 22490392 B im Plan): jede Zahl steht mit ihrem Lauf, keine wird als „der“ Wert
  weitergegeben.
- geprüft, ohne Befund: **Erzeugnis `docs/user/e2e-abdeckung.md`** (Schwerpunkt 4): byte-gleich (Messung oben), Ort-Anker
  `:3484` ist die Zeile der Ausgabezeile mit dem Anker der Deklaration.
- geprüft, ohne Befund: **`harness/README.md` Zeile 141** (Schwerpunkt 4): nur der Satzteil zur Phase ist geändert
  („in 6 Stücken zu je 10.000 Zeilen (das WAL je Stück liegt unter der Warnschwelle, die Summe über der Fehlerschwelle; der
  Runner wartet nach jedem Stück mit einer Frist von 30 s …)“, dazu „Last in Stücken seit slice-leerlauf-phase-last-in-stuecken“
  als Herkunftsfeld); die Zahlen sind Konstanten des Runners, keine Messwerte; die Zeile zu `make test-replication`
  (Form X1, anderer Gegenstand) ist unberührt.
- geprüft, ohne Befund: **Handbuch** (Schwerpunkt 4): die zwei Stellen und der neue Absatz sagen dieselbe Zusage wie
  Runner und README; `Version:` 1.65 → 1.66, `Stand:` 2026-09-25 → 2026-09-27 (Datum des Laufs), neue Zeile in
  `### Änderungshistorie` mit Kennungen und Slice-Namen; Schritt 17 des Implementer-Ablaufs erfüllt (Skill-HIGH
  „Handbuch-Versionshistorie nicht fortgeschrieben“ nicht erreicht). Keine neue Betreiber-Oberfläche (keine `CDC_*`-Variable,
  keine SQL-Funktion, kein Endpunkt; `wal_retention_error_bytes` besteht und wird verlinkt). Der Satz, das Steigen und
  Fallen des Wertes belege der Integrationstest, trägt: die Wartebedingung erzwingt, dass `confirmed_flush_lsn` nach jedem
  Stück die Position hinter dem Stück erreicht (Rückstand danach klein), und die Summe steht über der Fehlerschwelle.
  Der Betreiber-Hinweis „Klasse `replication`“ stimmt mit der gemessenen Signatur des Verdikts (M1) überein; die Klasse bei
  gehaltener Persistierung (`storage`) ist ein anderer Fall und liegt beim Träger-Slice.
- geprüft, ohne Befund: **Suchlauf-Feld (§3.13)** (Schwerpunkt 4): `make suchlauf-nachmessen` Exit 0 mit 11 Zeilen; Stände
  sind Commit-Kennungen (`d078700d`, `6a976f58`) und `diff`, die Plan-Datei ist am Parent per Pathspec, am Arbeitsbaum vom
  Werkzeug ausgeschlossen; Suchraum der ganze Baum ohne die drei Ausnahmen; drei Musterarten (Symbol, Beschreibung samt
  Hedge, Phase als Gegenstand) plus das Zählwort `60000`. Zusätzlich gegen den Baum gesucht (`nicht aktivierte[n]? Tabelle`,
  `nicht veröffentlichte`, `Leerlauf-Bestätigung` ohne die Ausnahmen): jeder Treffer außerhalb des Diffs ist im Feld
  behandelt (ADR-Zeilen, Register-Evidence, Folge-Pläne, Go-Kommentare der Produktion) oder beschreibt anderes (Mapper,
  Testnamen). Der Träger `walretention_slotgrowth_internal_test.go:28` bleibt wahr und ist mit Adresse benannt (Plan §3).
  Nichtgefundenes je Träger steht in der Tabelle.
- geprüft, ohne Befund: **Verweisform in fremden Dateien** (`3471c20d`, Schwerpunkt 4): `git diff --word-diff` zeigt in
  `roadmap.md` (zwei Stellen), `welle-transformationen.md`, `…/test-integration-retention-timing-flake/state.md`,
  `slice-start-vorlauf-grenze.md` und `slice-wal-fehlerschwelle-ausgangsklasse.md` (zwei Stellen) ausschließlich den
  Wechsel „Link auf den wandernden Pfad → Kennung in Inline-Code“, keine Aussage geändert; das Register und die Folge-Pläne
  sind im Plan als fremde Träger mit Frist (Closure) gemeldet, nicht inhaltlich mitgeändert.
- geprüft, ohne Befund: **Kommentare §3.7** (Schwerpunkt 6) in den geänderten Runner-Zeilen: Kopf der Phase (Zeilen
  3371–3383) und der Kommentar vor dem Schreiber (`:3444-3446`) im Indikativ, eine Kennung (`ADR-0120`), keine Kette, keine
  Kompaktform, kein „ff.“, keine Wiedergabe der Spec; die Aussage „bis die Bestätigung eines Stücks folgt, steht das Stück
  im Rückstand“ ist durch die Wartebedingung und das Verdikt (M5, M6) getragen, die Allaussage über „jeden“ Schreiber ist
  nicht mehr im Runner. Konjunktiv-Suche (Liste aus Schritt 20) ohne Treffer; der eine Treffer der weiteren Wortliste
  („nicht mehr als die Fehlerschwelle“) steht in einer Meldung, nicht in einem Kommentar.
- geprüft, ohne Befund: **Traceability, ID-Schema, Moves** (Schwerpunkt 6): alle 8 Betreffe tragen `LH-*`/`ADR-*`, keine
  `SPEC-*`/`ARC-*`; die Kennungen des Reports stehen in Inline-Code oder als Link; die Moves `ca93e19c`/`828e1625` sind
  rein, `c9a0308c` (Verantwortlich) ein eigener Commit; kein `//nolint` (`AGENTS.md` §3.2, ohnehin kein Go im Diff).
- geprüft, ohne Befund: **Docker-only im Diff** (Schwerpunkt 3, Diff-Seite): die geänderten Runner-Zeilen rufen `docker
  exec`/`docker compose` und die Runner-Hilfen auf, keine Host-Toolchain, kein `sed -i`, kein Interpreter, keine Umleitung
  auf eine Repo-Datei; das Erzeugnis schreibt der bestehende Schreibweg des Runners (byte-gleich nachgewiesen).
- geprüft, ohne Befund: **Verifier-Regel aus Verdikt §4** (Schwerpunkt 7): die Regel steht in den Plänen der drei
  Folge-Slices (seit `a5712aad`) und in §6 dieses Plans (Risiko 4) im Erwartungs-Ton („hergeleitet etwa 5 %“); der Plan
  behauptet nicht, in den Läufen des Implementers sei kein unbeabsichtigtes Rot aufgetreten — das ist eine Angabe des
  Implementers im Bericht und bleibt dem Verifier. Beobachtung: der einzige `e2e.yml`-Lauf zum Stand `38981f07` war im
  ersten Versuch grün (kein Fall der Regel).
- geprüft, ohne Befund: **Spec-Stratum, Zwei-Quellen-Drift, Sicherheit, kritischer Pfad**: `spec/` und `docs/plan/adr/` nicht
  berührt (F-5a: benannte, vom Architect entschiedene Präzisierung); Produktionscode und die Persist-before-ACK-Sektion
  nicht berührt; keine Geheimnisse im Diff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** DoD-Haken ohne committeten Beleg-Anker · Zusage ohne Bindung an ihre Schutzwirkung
(Grenze nicht benannt) · Docker-only, Wortlaut-Abweichung; Berechtigungs-Signal ohne Rückfrage · Bindung und Grenze des
Tests benannt · Träger-Nachzug (außerhalb des Diffs) · Begriff ohne Einführung

## Verdikt

**Merge-blockierend:** nein. Kein HIGH; die Phase färbt an allen acht Mutationen der Fehlerseite (m1, m2, m3, m5, m6, m7,
c, d) an der erwarteten Zeile, die zwei grünen Mutationen (m4, m8) sind eine benannte Grenze und kein Fehler des Aufbaus,
das Erzeugnis ist byte-gleich, und der `e2e.yml`-Lauf zum Stand `38981f07` war an beiden Legs im ersten Versuch grün.
Keine Fixrunde am Implementer: F-1 ist Plan-Text und wird mit der Closure-Notiz (§7) durch die Zeilen gedeckt (Lauf,
Ausgabezeile, Dauer); F-2 ist eine Ergänzung der Grenzen-Liste im Plan (Zeile zu m4 und m8) ohne Code; F-3 geht an den
Auftraggeber; F-4 bis F-6 sind Träger und Hinweise ohne erwartete Aktion. Die DoD-Zeile „Review durchgeführt“ ist im selben
Commit wie dieser Report auf `[x]` gezogen (Skill §DoD-Checkbox-Nachzug ohne Fixrunde).

**Einordnung der Fragen des Auftrags.**

- *Ist die Wartebedingung korrekt gebildet?* Ja (Negativbefunde, erster Punkt). Sie ist nicht trivial erfüllbar durch einen
  leeren Filter (`gelesen 'leer'`), und die Position ist nach dem Commit des Stücks gemessen.
- *Reichen (c) und (d) als Bindung der Wartebedingung, ist der Plan ehrlich?* Für die **Fehlerseite** ja: (c) und (d)
  färben, und die Wartebedingung reagiert auch auf Slot-Filter und Erwartungsform (m3, m7); der Plan sagt im Ist-Ton, dass
  (e) nicht an der Wartebedingung färbt und dass (c) und (d) sie binden (§2 DoD 1, §6 Risiko 1). Die **Schutzwirkung**
  ist nicht gebunden (m4, m8 grün) und im Plan nicht benannt (F-2); dass Mutation (e) die Backfill-Hälfte früher färbt, ist
  richtig berichtet — die Backfill-Hälfte selbst ist damit die Bindung der Leerlauf-Bestätigung an der Ebene der Phase.
- *Nebenbefund (f):* durch fünf Unit-Tests im Paket `receive` gebunden (F-4); kein Finding im Slice, kein Register-Fall.
- *Docker-only-Probe am Vorgehen des Implementers:* Verstoß im Sinn der HIGH-Klasse **nein**, Abweichung vom Wortlaut von
  `AGENTS.md` §3.1 **ja, in der Auslegung des Repos gedeckt**, Einstufung LOW (F-3). Was daraus für die Rollen-Aufträge
  folgt, ist eine Entscheidung des Auftraggebers, nicht des Reviewers: (1) welcher Weg für ein Mutations-Image gilt —
  `make image` aus einer `git archive`-Kopie mit anschließender Wiederherstellung von `:dev` (Vorgehen der zwei früheren
  Verifikationen) oder ein Wegwerf-Tag über einen benannten Weg —, (2) ob ein verweigerter Aufruf im Bericht genannt und
  vor einem Ersatzweg zurückgefragt wird. Bis das steht, hält der Reviewer für sich fest: kein Bau eines Images, keine
  Aktion nach einer Verweigerung auf anderem Weg.

**Übergabe:** F-1 an den Implementer/Planner für den Closure-Schritt (Zeile des grünen Volllaufs mit Lauf, Dauer und
Ausgabezeile; die CI-Zeilen oben sind nachgelesen und können als Anker dienen); F-1 zusätzlich an den Architect für den
Steering-Loop (vier Auftreten derselben Klasse, Skill §Pflege: Regel, Sensor oder Schärfung des Schritts 18 des
Implementer-Ablaufs). F-2 an den Planner (Grenzen-Zeile im Plan, §2 DoD 1 oder §6 Risiko 1). F-3 an den Auftraggeber
(Regel für Mutations-Images und für den Umgang mit verweigerten Aktionen). F-4 bis F-6 ohne erwartete Aktion; F-5b bleibt
als Adresse im Träger-Slice `slice-wal-fehlerschwelle-ausgangsklasse`. An den Verifier: der ungeprüfte Rest — DoD 5 (der
erste grüne `e2e.yml`-Lauf mit beiden Legs; Lauf, Versuch und Job-Kennungen oben nachgelesen), die Mutation (e) (Image
mit mutiertem Produktionscode, übernommen), ein voller lokaler `make test-integration` (nicht von mir gefahren; die Phase
lief im Wegwerf-Aufbau grün). Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report
ist ein Lauf-Beleg und ersetzt keine Verifikation.
