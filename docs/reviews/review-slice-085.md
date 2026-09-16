# Review-Report: slice-085 — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `slice-085`
(`docs/plan/planning/in-progress/slice-085-receive-naht.md`), Diff
`95480a5..1abc7c5` (drei Commits: `610751c` Naht · `506c811` Plan-Nachzug und
Null-Befund · `1abc7c5` Digest-Beleg). Sechs Dateien:
`internal/adapters/driving/replication/receive/{seam.go,seam_test.go,receive.go,walretention.go}`,
`harness/image-hash.txt`, der Slice-Plan selbst. **Kein** anderer Träger ist
berührt — `Dockerfile` (Stufe `coverage`), `tools/harness/db-coverage.sh`,
`harness/mk/coverage.mk`, `internal/bootstrap/wiring.go`, `.a-check.yml`,
`spec/**` sind diff-frei (eigene Prüfung, §Eigene Messungen).

**Skill:** `.harness/skills/reviewer.md` @ Repo-Stand (vier repo-spezifische
HIGH-Regeln) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-085-receive-naht.md`
  vollständig (§1–§8)
- [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
  (Nahtform: Hülle + fake-fähige Fläche, **im Paket**; Festlegung 1/2/3/4/5,
  §Fitness Function, §„Die benannte Grenze", Re-Evaluierungs-Trigger (a)/(e))
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3/4/5 · [`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  (dreiteiliger Nachweis, hier als Null-Befund) ·
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Digest ist
  Lauf-Beleg) · [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md)
  (Option C) · [`ADR-0006`](../plan/adr/0006-replication-stream-driving-adapter.md) ·
  [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md)
- `AGENTS.md` §3.1, §3.6, §3.7, §3.9, §3.11, §5 · `harness/conventions.md`
  `MR-000` (ID-Schema)
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/`
  (namentlich `db-gegenstand-enthaelt-netzlos-geprueften-code` (1×),
  `zahl-in-traeger-driftet-gegen-die-messung` (2×))
- `docs/reviews/review-slice-084.md` / `verify-slice-084.md`
  (Geschwister-Naht, vorherige Findings am gleichen Modul) ·
  `docs/plan/planning/done/slice-084-postgresack-naht.md` (die Vergleichs-Form)
- `harness/sensors/db-adapter-coverage.md`, `harness/sensors/coverage-gate.md`,
  `tools/harness/db-coverage.sh`, `Dockerfile` (Stufe `coverage`)

---

## Findings

### F-1 — Der Slice bewegt die **Zustandsgröße** des Gegenstands (659 → 691) und lässt den Träger stehen; die Klasse ist damit beim dritten Vorgang

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was
  da ist") · [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
  §Fitness Function (der Null-Befund ist „Bringschuld des Laufs und
  Review-Gegenstand") · Klasse `zahl-in-traeger-driftet-gegen-die-messung`
  (Register **2×**, aus `slice-081`/`slice-084` — dieser Vorgang wäre der
  **dritte**)
- `pfad`: `harness/sensors/db-adapter-coverage.md:65` („Der gemergte Nenner ist
  **659 Statements** (`postgresstorage` 472 · `postgresack` 32 ·
  `replication/receive` **155**) — die **Zustandsgröße** dieses Gegenstands")
  gegen `:81` (Kalibrierungs-Zeile: „**477 von 650** Statements = 73,38 %
  (die geltende Größe des Gegenstands: §Zählbasis)")
- `befund`: Der DB-Gegenstand trägt nach diesem Diff **691** Statements mit
  **532** gedeckten (76,99 %), `replication/receive` **187** statt 155 — eigene
  Messung des gemergten Profils (§Eigene Messungen). Das Sensor-Dokument führt
  den Nenner seit der `slice-084`-Fixrunde ausdrücklich als **Zustandsgröße**
  (nicht als lauf-gebundenen Beleg) und nennt daneben 650, wo §Zählbasis 659
  sagt; beide Werte sind nach diesem Zug falsch. Der Plan erklärt
  `harness/sensors/**` unter „**Nicht angefasst**" (§3) — die Drift ist damit
  Folge der erklärten Grenze, genau wie bei `review-slice-084` F-1 (dort
  ebenfalls **LOW**). **Unterschied zu `084`:** dort war der Träger eine Zahl,
  die der Zug *falsch gemacht* hatte; hier ist der Nenner vom Vorgänger
  ausdrücklich zum **Zustand** erklärt worden — wer ihn bewegt, muss ihn
  ziehen. Die Kategorie bleibt dennoch LOW und wird nicht auf MEDIUM gehoben:
  die Schwellen-Regel des Skills („Muster, das schon zweimal LOW war") ist
  nicht erfüllt — die Klasse war einmal HIGH (`slice-081` F-1) und einmal LOW
  (`slice-084` F-1). **Folge für den Lese-Schritt:** das Register erreicht mit
  diesem Vorgang **3×**; ein Ausgang (verkörpert/geplant/gestrichen) ist fällig.
  *Nebenbefund, nicht diesem Vorgang zuzurechnen:* die 650 in `:81` widersprach
  schon nach der `084`-Fixrunde der eigenen §Zählbasis (659) — sie ist von
  dort stehengeblieben, wird von diesem Zug aber zusätzlich überholt.
- `verifizierbar`: ja — das gemergte `-coverprofile` aus
  `make test-store` + `make test-replication`, per `awk` über die Block-Position
  dedupliziert (§Eigene Messungen, Zeile „DB-Adapter-Coverage gemergt")

### F-2 — Die „Benannte Grenze" nennt die Verdünnung als **Nenner-Wachstum**, nicht als **Anteil** — die Größe, an der `ADR-0080` Trigger (a) hängt

- `kategorie`: INFO
- `quelle`: [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
  §„Die benannte Grenze" + Re-Evaluierungs-Trigger **(a)** („wird die
  Verdünnung **materiell**") · `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`
- `pfad`: `docs/plan/planning/in-progress/slice-085-receive-naht.md:245-257`
  (§3, „Benannte Grenze")
- `befund`: §3 meldet `replication/receive` 155 → 187 (+32 Statements) und
  „153 von 187 gedeckten" — beides gemessen bestätigt. Die Zahl, die den
  Trigger (a) entscheidet, steht dort nicht: **netzlos** (ohne jede
  PostgreSQL-Verbindung) gedeckt sind vorher **11 von 155 (7,10 %)**, nachher
  **112 von 187 (59,89 %)**; im Gegenstand insgesamt steigt der netzlos
  gedeckte Anteil der gedeckten Statements von **41/491 (8,35 %)** auf
  **142/532 (26,69 %)** (§Eigene Messungen). `slice-084` hatte +9 Statements
  gebracht, dieser Zug bringt +32 plus 101 netzlos gedeckter Alt-Statements —
  die Verdünnung ist mit diesem Vorgang in derselben Größenordnung wie das
  reale Messergebnis (von den 153 gedeckten Statements des größten
  Gegenstands-Pakets brauchen nur **41** eine PostgreSQL-Instanz). Kein Defekt
  des Diff (die Grenze ist in `ADR-0080` entschieden und der Implementer durfte
  die Zahl nicht „heilen"); die Handlung liegt bei der Closure: die Anteilszahl
  als Datum des Triggers (a) führen **und** `evidence/slice-085.md` in
  `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` anlegen (Zähler dann
  **2×** — siehe Schwerpunkt 5).
- `verifizierbar`: ja — zwei `go test -coverpkg=… ./…/receive`-Läufe (mit und
  ohne `CDC_REPLICATION_TEST_DSN`), Block-Position-dedupliziert

### F-3 — Der führende Kommentar in `seam_test.go` ist syntaktisch ein **zweiter** Package-Doc-Kommentar (Wiederholung von `review-slice-084` F-3)

- `kategorie`: INFO
- `quelle`: Maintainability (Doc-Kommentar-Konvention; kein Gate)
- `pfad`: `internal/adapters/driving/replication/receive/seam_test.go:1-8`
  gegen `internal/adapters/driving/replication/receive/receive.go:1-9`
- `befund`: `seam_test.go` erklärt `package receive` (paket-intern, nötig für
  die unexportierten Naht-Bezeichner) und trägt seine Test-Erklärung ohne
  Leerzeile unmittelbar vor der Paket-Klausel — damit ist sie eine zweite
  Package-Doc-Kandidatin neben der in `receive.go`. Es ist dieselbe Form, die
  `review-slice-084` F-3 an `postgresack/seam_test.go` benannt hat: die
  Wiederholung ist strukturell bedingt (paket-interner Test) und sachlich
  unvermeidlich, die **Form** ist es nicht. Kein Defekt, kein Gate, keine
  erwartete Aktion.
- `verifizierbar`: ja — `go list -f '{{.Doc}}' ./internal/adapters/driving/replication/receive`
  zeigt den `receive.go`-Text; `go vet` ist still

### F-4 — Mutations-Überlebender: die netzlose Suite ist gegen die **SQL-Prädikate** der Katalogabfragen blind

- `kategorie`: INFO
- `quelle`: Maintainability (Test-Stärke der Fake-Seite) gegen `slice-085` §6
  Risiko 3 („Der Fake könnte grün sein, ohne etwas zu prüfen")
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:245-247`
  (`slotLSNQuery`) gegen `seam_test.go:353-358`
  (`TestSlotLSNQueryNamesTheSlot`)
- `befund`: Entfernt man ` AND slot_type = 'logical'` aus `slotLSNQuery`,
  bleibt die netzlose Suite **grün** (Exit 0, eigene Probe M3) — der Fake
  ignoriert den SQL-Text, und der einzige Test gegen die Abfrage prüft zwei
  `strings.Contains` (Slot-Name, `confirmed_flush_lsn`), nicht das Prädikat.
  Dasselbe gilt für `ensurePublication`s `WHERE pubname = …`: der Fake
  beantwortet jede Abfrage gleich. Die reale Schicht fängt das nicht
  zwangsläufig (die realen Tests legen Slot und Publication an, also *existieren*
  beide) — die Lücke ist damit keine Verhaltens-Lücke, sondern eine
  **Sensitivitäts**-Lücke der neuen Testfläche. Keine erwartete Aktion; benannt,
  damit „die Fake-Seite fährt die Verklebung" (Liefer-Punkt 2) nicht als
  „sie prüft die Abfragen" gelesen wird.
- `verifizierbar`: ja — Wegwerf-Kopie, Mutation M3 (§Eigene Messungen)

---

## Negativbefunde

- geprüft, ohne Befund: **kein Verhaltens-Change** — Zeile-für-Zeile-Vergleich
  von `receive.go`/`walretention.go` gegen den Vorstand: die Fehlertexte sind
  wörtlich gleich (`"%w: leeres CopyData"`, `"%w: XLogData: %v"`,
  `"%w: Keepalive: %v"`, `"%w: Keepalive-Antwort: %v"`,
  `"%w: confirmed_flush_lsn %q: %v"` — letzterer über `parseLSN(what, text)`
  mit `what` = `confirmed_flush_lsn` bzw. `ConsistentPoint`, wörtlich die zwei
  alten Stellen), der Dispatch (`XLogData`/`PrimaryKeepalive`/leeres CopyData)
  ist deckungsgleich, die Keepalive-Antwort trägt weiter **drei** LSN-Felder aus
  `s.lastAcked`, `defer s.conn.Close` wurde zu `defer s.session.Close` bei
  identischem Feldwert, `reconnectAfterError` setzt die Verbindung wie zuvor neu.
  Die vier `default`-Zweige (`Run`-Typ-Switch, Byte-ID-Switch) sind der alte
  Durchfall, jetzt explizit. Belegt zusätzlich durch die eigenen Mutationen
  (M1/M2 rot, M3/M4 netzlos grün, M4 real rot) und die realen Tiere.
- geprüft, ohne Befund: **die reale Verdrahtung ist unverändert grün** —
  `make test-store` **Exit 0**, `make test-replication` **Exit 0**, `make test`
  **Exit 0** (31 Pakete `ok`, kein `FAIL`), `make gates` **Exit 0** — je
  ungepiped, Gate und Folgehandlung getrennt (§3.9).
- geprüft, ohne Befund: **keine abgeschwächte Zusicherung, keine Maskierung** —
  `stream_test.go` ist **byte-identisch** (SHA256 `7b5d6a04…` an beiden Enden
  der Range, `git diff` leer); `*_test.go`-Diff zeigt genau **eine neue** Datei
  (`seam_test.go`), **keine** gelöschte; kein neues `t.Skip`/`SkipNow`, kein
  `//nolint`/`#noqa`, kein `|| true` im Diff (Muster-Grep leer).
- geprüft, ohne Befund: **der öffentliche Rand ist unverändert** — Mengenvergleich
  der exportierten Symbole beider Bäume ergibt für die Produktionsfläche exakt
  dieselbe Menge (`NewStream`, `NewWALRetentionChecker`, `Config`, `Stream`,
  `WALRetentionChecker`, `ErrConfiguration`, `ErrReplication`); die einzigen
  Neuzugänge sind `Test*`-Funktionen. `NewStream`, `Stream.Conn/Run/Assembler/
  BindCapture`, `Measure`, `Close` tragen ihre alten Signaturen; die Naht selbst
  (`driverSession`, `connSession`, `newStreamOnSession`,
  `newWALRetentionCheckerOnSession`) ist unexportiert — „kein neuer öffentlicher
  Rand", wie `ADR-0080` Festlegung 5 es für dieses Paket zusagt.
- geprüft, ohne Befund: **die Naht deckt sich wörtlich mit `ADR-0080`
  Festlegung 5** — sieben Operationen (`IdentifySystem`,
  `CreateReplicationSlot`, `StartReplication`, `SendStandbyStatusUpdate`,
  `Exec → []*pgconn.Result`, `ReceiveMessage → pgproto3.BackendMessage`,
  `Close`), Hülle `connSession{conn *pgconn.PgConn}`, Operationen in der
  **Wertform** des Treibers (Festlegung 2: kein Draht-Typ, keine
  Anwendungsfall-Sprache), konkreter Typ nur in Hülle und Dial
  (`connectReplication`, Festlegung 3), `var _ driverSession = connSession{}`
  liegt **im Produktionspfad** (`seam.go`, nicht im Test) — genau die
  Fitness-Function-Zeile der ADR.
- geprüft, ohne Befund: **die Fläche ist fake-fähig, und der Fake erfüllt
  dieselbe Schnittstelle** — `var _ driverSession = connSession{}` (Paket) und
  `var _ driverSession = (*fakeSession)(nil)` (Test) kompilieren beide; der
  Fake liefert `pgproto3.CopyData`/`CopyDone`/`ErrorResponse` und
  `[]*pgconn.Result` und braucht keinen exportierten Konstruktor. Eigene
  Gegenprobe an der Hülle: eine **nicht delegierende** `connSession.Exec`
  (M4) bleibt netzlos **grün** und färbt den realen Tier **rot** (7 Fälle) —
  siehe Schwerpunkt 4.
- geprüft, ohne Befund: **die Liste der netzlos prüfbaren Teile aus
  `ADR-0080` Festlegung 5 ist abgedeckt** — Empfangs-Schleife (Byte-ID-Dispatch, Keepalive-Antwort aus
  `lastAcked`, Kontext-Ende, leeres CopyData), Slot-Auflösung (bestehend →
  `confirmed_flush_lsn`; fehlend → `CreateReplicationSlot`), Publication-Grenze,
  Katalog-Zeilen-Übersetzung, Rückstands-Messung: je mindestens ein Test in
  `seam_test.go` (26 Testfälle).
- geprüft, ohne Befund: **kein anderer Träger berührt** — `git diff --name-status`
  über die Range listet sechs Pfade; die Gegenstands-Listen sind diff-frei
  (`Dockerfile`, `tools/harness/db-coverage.sh` → `DB_COVERAGE_PKGS` und
  `DB_COVERAGE_THRESHOLD`, `harness/mk/coverage.mk` → `THRESHOLD` bleibt 70,
  beide Endstufen 80 %, **keine Schwellen-ADR fällig**), ebenso
  `internal/bootstrap/wiring.go`, `.a-check.yml`, `spec/**`. Kein Paketwechsel,
  kein Unterpaket — der Gegenstand bleibt, wie er ist.
- geprüft, ohne Befund: **der Null-Befund (a)/(b)/(c) ist gemessen, nicht
  behauptet** — (a) Unit-Nenner **1903** (1369 gedeckt, 71,94 % dedup / Stufe
  71.90 %) an **beiden** Enden der Range selbst gemessen, und das Unit-Profil
  enthält **keine** Zeile aus `postgresack`/`replication/receive` (`k_ab = 0`,
  Grep = 0); `postgresstorage` 472 und `postgresack` 32 byte-stabil.
  (b) Paket-Diff mit **Range und Pathspec**, beide Richtungen: vier Dateien
  unter `…/receive/`, die Gegenrichtung `-- internal/ ':!…/receive/'` **leer**
  (nachgeprüft, ebenso über die volle Range `95480a5..1abc7c5`). (c) kein
  Testfall entfernt, die realen Läufe grün.
- geprüft, ohne Befund: **der Plan zieht die Lehren aus `review-slice-084`
  F-2** — §3(b) zitiert seinen Beleg jetzt mit Range **und** Pathspec **und**
  die Gegenrichtung; §3(c) nennt die `--name-status`-Form. Die zitierte Form
  trägt die Aussage (nachgeprüft), nicht nur die Aussage die Form.
- geprüft, ohne Befund: **`ADR-0044`-Form** — der Zug ändert
  Build-Kontext-Dateien, `harness/image-hash.txt` ist der deklarierte Träger,
  der Commit `1abc7c5` ist ein eigener `chore(image)`-Commit mit `ADR-*` im
  Betreff. Eigene Gegenprobe: ein frischer
  `docker buildx build --load --metadata-file /tmp/…` dieses Baums liefert
  **denselben** Digest `sha256:84bdca56…` wie der getrackte Träger; der Baum
  blieb unberührt. Das ist ein gültiger Befund im selben Builder — **kein**
  Inhalts- oder Umgebungs-Beweis (`ADR-0044` §1/§3).
- geprüft, ohne Befund: **Hygiene** — kein Lauf-Artefakt in einem der drei
  Commits (Datei-Listen: Naht, Plan, `harness/image-hash.txt`);
  `tools/schema/plan.yaml` liegt in **keinem** Commit und wurde nach den
  Schema-fahrenden Zielen dieses Laufs zurückgenommen (`git status` leer);
  keine Trailer (`Co-Authored-By`/`Signed-off-by`/„Generated with"), Betreffe
  mit `ADR-*` und ohne `SPEC-*`/`ARC-*` (`commit-traceability: OK — 5
  Commit(s), Betreffs ohne Struktur-ID`); keine host-lokalen Pfade (d-check
  inkl. `hostpaths`-Modul: 712 Dateien, 0 Befunde); keine unaufgelösten
  `BEO-`/`CO-`/`MR-`-Kennungen im Diff; Historie **linear auf `main`**
  (0 Merges, je Commit ein Elternteil, `git rev-parse --abbrev-ref HEAD` =
  `main`).
- geprüft, ohne Befund: **Vorlagen-Reste** — §2 trägt keine doppelte DoD-Zeile
  und keinen Platzhalter; die vier `BEDIENHINWEIS`-Kommentare und die
  `<…>`/`<bei Closure>`-Felder in §6/§7 stehen **vor** dem `git mv` nach
  `done/` (dort fängt der `forbid-pattern` der fünften `structure`-Regel sie);
  dieser Stand ist erwartet und deckt sich mit `review-slice-084`
  (Negativbefunde).
- geprüft, ohne Befund: **Slice-/Wellen-Chronik in Produktionscode** — weder
  `seam.go` noch `receive.go`/`walretention.go` nennt `slice-<NNN>` oder
  `welle-<NN>`; alle Herkunfts-Anker sind `ADR-*`/`LH-*`/`SPEC-*`. Die
  `slice-0NN`-Nennungen des Pakets stehen ausschließlich in `stream_test.go`
  (vorbestehend, Subjekt ist der Testfall — zulässige Provenienz).
- geprüft, ohne Befund: **keine Zwei-Quellen-Drift neu erzeugt** — im Gegenteil:
  `slotLSNQuery`/`parseLSN` **beseitigen** eine bestehende Doppelung (das
  SQL-Literal für `confirmed_flush_lsn` stand vorher in `receive.go` **und**
  `walretention.go`; die LSN-Wrapping-Stelle ebenso). `standbyStatus` liegt
  nun — wie `connSender`/`connSession` — je Adapter-Paket einmal; mangels
  Adapter→Adapter-Import (`.a-check.yml`) ist das die vorgesehene Form
  (`ADR-0080` Trigger (d) greift erst beim **dritten** Adapter).
- geprüft, ohne Befund: **§8-Sichtung** — der Plan nennt
  `adapter-fehler-ausgang` (2×) und `spec008-replication-luecke` (2×) als
  „benachbart, kein Treffer" bzw. „kein Treffer"; die Registerstände sind am
  Verzeichnis nachgezählt (beide 2×, keine erreicht 3×) — die Einschätzung
  trägt. Die **neue** Berührung (`db-gegenstand-enthaelt-netzlos-geprueften-code`)
  war bei der Plan-Anlage noch nicht sichtbar — sie entstand erst mit `084`;
  sie gehört in §7 (offen, siehe F-2), nicht in §8.

---

## Eigene Messungen (Exit-Codes ungepiped, Docker-only)

| Lauf | Exit | Ausgabe / Ergebnis |
|---|---|---|
| `make gates` (Log-Datei, Exit danach separat gelesen) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · d-check 712 Dateien / 0 Befunde (inkl. `hostpaths`) · d-check `commits` 0 Befunde · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `coverage-gate: OK — Coverage 71.90% erfüllt Schwelle 70%` · a-check 0 Befunde |
| `make test` (Race-Detector, netzlos) | **0** | **31 Pakete `ok`**, kein `FAIL`, 0 `--- FAIL` |
| `make test-store` | **0** | `…/postgresstorage 5.589s` |
| `make test-replication` (frischer store- **und** replication-Lauf) | **0** | `DB-Adapter-Coverage: 76.99% (gedeckt 532 von 691 Statements; Profile gemergt: store,replication)` · `db-coverage: OK … Schwelle 70%` · Tier-Phase: 33 × `ok`, 0 × `FAIL` |
| Dedup `merged.coverprofile` (eigene `awk`, Block-Position) | — | `postgresack 32/32` · `postgresstorage 347/472` · `receive 153/187` (`receive.go` 124/155, `seam.go` 7/7, `walretention.go` 22/25) · **GESAMT 532/691 = 76,99 %** |
| Dedup `/out/coverage.out` der `coverage`-Stufe **am Parent** `95480a5` (eigener `docker build --target coverage` über `git archive 95480a5`) | 0 | **1369/1903 = 71,94 %**; `postgresack`/`replication/receive`-Zeilen = **0** |
| Dedup `/out/coverage.out` der `coverage`-Stufe **am HEAD** `1abc7c5` | 0 | **1369/1903 = 71,94 %** (Stufe druckt 71.90 %) — identisch, `k_ab = 0` beidseitig |
| DB-Messphase **am Parent** `95480a5` (reale PostgreSQL, `-coverpkg`) | 0 | `receive 112/310` → dedup **112/155** · `postgresack 32/64` → **32/32** |
| Paket-Coverage **netzlos** (ohne DSN), Parent | 0 | `receive` **11/155 = 7,10 %** (`receive.go` 6/131, `walretention.go` 5/24) |
| Paket-Coverage **netzlos** (ohne DSN), HEAD | 0 | `receive` **112/187 = 59,89 %** (`receive.go` 92/155, `walretention.go` 20/25, `seam.go` **0/7**) |
| Paket-Coverage **netzlos** (ohne DSN), HEAD, `postgresack` | 0 | **30/32 = 93,75 %** (deckt die Zahl aus `review-slice-084` Schwerpunkt 4) |
| **M1** — Keepalive-Antwort-Bedingung invertiert (`if replyRequested` statt `if !replyRequested`) | **1** | `--- FAIL: TestRunAnswersKeepaliveWithAcknowledgedPosition` · `TestRunKeepaliveWithoutReplyRequestedSendsNothing` · `TestRunReportsKeepaliveReplyFailure` |
| **M2** — `parseXLogData` liefert den **rohen** Puffer (Byte-ID nicht abgeschnitten) | **1** | `--- FAIL: TestParseXLogDataCarriesPayload` |
| **M3** — `slotLSNQuery` ohne `AND slot_type = 'logical'` | **0** | **Überlebender** (F-4): netzlos grün |
| **M4** — Hülle delegiert `Exec` **nicht** (`return nil, nil`) | netzlos **0** / real **1** | netzlos grün (die Fake-Tests fahren die Hülle nie) · real **7 Fälle rot**: `TestStreamTranslatesRealChanges`, `TestStreamOpenTransactionNotConsumable`, `TestStreamTruncateUnsupported`, `TestStreamKeepaliveReportsAcknowledgedPosition`, `TestWALRetentionMeasuresGrowingBytes`, `TestWALRetentionMeasureReconnectsAfterConnectionLoss`, `TestStreamRestartsOnExistingSlot` |
| Kontrolle zu M4: unveränderte Basis im **selben** Harness | 0 | `ok …/receive 5.181s`, `ok …/postgresack` — die sieben Fälle sind der Mutation zuzurechnen |
| Exportierte Symbolmengen Parent vs. HEAD (Produktionsfläche) | — | identisch; Neuzugänge ausschließlich `Test*` |
| `docker buildx build --load --metadata-file /tmp/…` dieses Baums | **0** | `containerimage.digest = sha256:84bdca56…` — identisch mit `harness/image-hash.txt`; getrackter Träger unberührt (`git status` leer) |
| `git diff --name-only 95480a5..1abc7c5 -- internal/ ':!…/receive/'` | — | **leer** (Gegenrichtung des Paket-Diffs) |

Alle Mutationen liefen in **Wegwerf-Kopien** (`git archive` unter `/tmp`), der
Arbeitsbaum wurde nie mutiert; das Schema-Lauf-Artefakt `tools/schema/plan.yaml`
ist zurückgenommen, `git status` ist am Ende leer. **Exit-Code-Hinweis:** die
M4-Netzlos-/Real-Zahlen stammen aus je einem eigenen, ungepipeten `go test`; die
Rohform über `make` würde wegen GNU make **2** melden (der Implementer nennt
„Exit 2" — dieselbe Zahl, anderer Träger).

---

## Antwort auf die Schwerpunkte

1. **Kein Verhaltens-Change — die tragende Zusage hält.** Die reale Verdrahtung
   geht **unverändert** durch alle Tiere (`make test-store`/`make test-replication`/
   `make test` je Exit 0), `stream_test.go` ist **byte-identisch**, keine
   Testdatei entfernt, kein neues `t.Skip`, kein `|| true`, kein `//nolint`. Die
   Fehlertexte, der Dispatch und die Keepalive-Form sind wörtlich die alten; die
   Produktionsfläche exportiert exakt dieselbe Symbolmenge wie zuvor. Die Fakes
   sind **Zusatz, nicht Ersatz**: die M4-Probe zeigt, dass die netzline Suite
   eine **nicht delegierende Hülle** nicht sieht und der reale Tier sie fängt.
2. **Der Null-Befund (a)/(b)/(c) trägt — alle drei Teile selbst gemessen, und
   die „vor"-Spalte ist keine Übernahme.** (a) Unit-Nenner **1903** mit **1369**
   gedeckten an **beiden** Enden der Range real gemessen (eigener
   `coverage`-Stage-Bau am Parent `95480a5` **und** am HEAD) — Δ 0, und das
   Unit-Profil enthält keine Zeile der zwei ausgenommenen Pakete (`k_ab = 0` für
   die Unit-Fläche). (b) Der Paket-Diff ist auf **ein** Paket isoliert, mit Range
   **und** Pathspec in **beide** Richtungen (die Gegenrichtung ist leer), und die
   Gegenstands-Listen sind diff-frei — die Form, die `review-slice-084` F-2
   vermisst hatte, ist hier gezogen. (c) kein Testfall entfernt, die realen
   Läufe grün. Die Plan-Tabelle stimmt Zeile für Zeile: 691/532/76,99 %,
   `postgresstorage` 472, `postgresack` 32, `receive` 187 — und die „vor"-Werte
   659/491/74,51 % sowie `receive` **112/155** sind am Parent-Commit
   **nachgemessen**, nicht aus dem Vorgängerbericht übernommen.
3. **Die Hülle trägt.** `driverSession` spricht die sieben Operationen in der
   **Wertform** des Treibers aus; `connSession` hält `*pgconn.PgConn` als Feld
   und delegiert 1:1; `var _ driverSession = connSession{}` steht im
   **Produktionspfad**. Die Zusicherung, die `slice-081` nicht halten konnte,
   trägt hier in der von `ADR-0080` Festlegung 1 gemessenen Form: Hülle **und**
   Fake erfüllen dieselbe Schnittstelle, und der konkrete Typ reist nirgends
   über die Naht (Dial und `Conn()` bleiben die einzigen Träger).
4. **Die fünfte Gegenprobe ist reproduziert** — sie ist die Aussage, auf der die
   ganze Arbeitsteilung dieses Slice ruht. Eigene Probe M4 (nicht delegierendes
   `connSession.Exec`, `return nil, nil`): die neuen netzlosen Tests bleiben
   **grün**, der reale Tier färbt **genau sieben** Fälle **rot** — reproduziert,
   zusätzlich mit unveränderter Kontrolle im selben Harness abgesichert (Exit 0).
   Beide Mutationen des Implementers, die ich nicht wiederholt habe, sind
   konsistent mit meinen: die Slot- und die Katalog-Zeilen-Mutation fallen in
   dieselbe Klasse wie M1/M2 (rot). **Wächter-Rolle:** die Hülle fängt der reale
   Tier, nicht der Fake — das ist genau die Verortung aus §1.
5. **Die Verdünnung ist dieselbe Beobachtung — aber nicht derselbe Betrag.**
   *Urteil:* Es ist **eine** Beobachtung
   (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`); die Kennung ist
   die Identität, und ein zweiter Eintrag mit neuem Slug wäre genau das, was die
   Register-Regel („ohne Kennung zählt eine Umformulierung als zweite
   Beobachtung, und keine erreicht je 3×") verbietet. `slice-084` und `slice-085`
   sind **zwei Vorgänge** der **einen** Klasse — der Zähler geht mit
   `evidence/slice-085.md` auf **2×**. Zugleich ist das „kein neuer Befund"
   nicht „kein neues Gewicht": `084` hat **+9** Statements in ein Paket mit 32
   gestellt, `085` stellt **+32** plus **101** netzlos gedeckter Alt-Statements
   in das **größte** Paket — der netzlos gedeckte Anteil springt von **7,10 %**
   auf **59,89 %** in `receive` und von **8,35 %** auf **26,69 %** im ganzen
   Gegenstand; von den 153 gedeckten `receive`-Statements brauchen nur **41**
   eine PostgreSQL-Instanz. Damit ist dies der Vorgang, an dem `ADR-0080`
   Trigger (a) („wird sie materiell") seine erste belastbare Zahl bekommt. Der
   Implementer darf die Zahl **nicht heilen** (die Grenze ist entschieden) — aber
   der Anteil gehört in die Closure-Notiz und in den Register-Beleg (F-2).
6. **Hygiene sauber, Historie linear auf `main`.** Drei Commits (Naht · Plan ·
   Digest), kein Lauf-Artefakt, keine Trailer, Betreffe mit `ADR-*` und ohne
   Struktur-IDs, `commit-traceability` OK, d-check inkl. `hostpaths` 0 Befunde,
   keine unaufgelösten Kennungen, 0 Merges. Reste, die **vor** der Closure
   erwartet sind: die vier `BEDIENHINWEIS`-Blöcke und die §6/§7-Platzhalter —
   sie gehören dem `done/`-Sensor, der sie beim `git mv` fängt. Einzige
   Ausnahme: der Digest-Beleg ist verifiziert (derselbe Digest aus einem
   frischen Bau dieses Baums).

---

## Formvergleich `slice-084` ↔ `slice-085`

Dieselbe Form, zwei Größenordnungen — die Abweichungen im Einzelnen:

| | `084` (`postgresack`) | `085` (`receive`) | begründet? |
|---|---|---|---|
| Schnittstelle | `standbySender`, 1 Operation | `driverSession`, 7 Operationen | **ja** — `ADR-0080` Festlegung 5 nennt beide Wort für Wort; die Fläche ist die gemessene (18 Symbole) |
| Hülle | `connSender` | `connSession` | ja — dieselbe Form, eigener Name je Paket (kein Adapter-Import) |
| paket-interner Einstieg | `newOnSender` | `newStreamOnSession` **und** `newWALRetentionCheckerOnSession` | **ja** — `receive` trägt zwei Naht-Träger (`Stream`, `WALRetentionChecker`); ein zweiter Einstieg ist der kleinste Weg, auch den Checker netzlos zu fahren |
| Zusicherung | `var _ standbySender = connSender{}` | `var _ driverSession = connSession{}` | ja — beide im Produktionspfad (`ADR-0080` §Fitness Function) |
| reine Funktionen | `ackLSN`, `standbyStatus`, `replicationClass` | `parseXLogData`, `parseKeepalive`, `parseLSN`, `firstRow`, `standbyStatus`, `slotLSNQuery` | ja — und **neu**: `slotLSNQuery`/`parseLSN` **entdoppeln** zwei Stellen, die vorher in `receive.go` **und** `walretention.go` standen |
| §3(b)-Beleg | ohne Range/Pathspec (`084` F-2) | **Range + Pathspec + Gegenrichtung** | ja — die Lehre ist gezogen, das ist eine **Verbesserung** gegenüber dem Geschwister |
| Test-Datei-Kommentar | Package-Doc-Kandidat (`084` F-3, INFO) | derselbe Fall (`085` F-3) | **nein** — dieselbe Form wiederholt, sachlich unvermeidlich, harmlos |
| `harness/sensors/**` | „Nicht angefasst" → Zahl driftete (`084` F-1, LOW, in der Fixrunde gezogen) | „Nicht angefasst" → Zahl driftet wieder (`085` F-1, LOW) | **nein** — und schwerer: die Fixrunde hat den Nenner ausdrücklich zur **Zustandsgröße** erklärt; wer ihn bewegt, muss ihn ziehen |

Zwei echte Abweichungen (beide begründet und beide Verbesserungen: der Beleg in
§3(b), die entdoppelten Helfer), eine **unbegründete** (der stehengelassene
Sensor-Nenner — F-1), eine folgenlose Formwiederholung (F-3).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** „Träger-Zahl driftet gegen die bewegte
Zustandsgröße" (F-1, Klasse `zahl-in-traeger-driftet-gegen-die-messung` —
**dritter Vorgang**, Register erreicht 3×) · „Verdünnung als Nenner-Wachstum
statt als Anteil berichtet" (F-2, Klasse
`db-gegenstand-enthaelt-netzlos-geprueften-code` — **zweiter Vorgang**, 2×) ·
„Test-Datei trägt Package-Doc-Kommentar" (F-3, zweiter Vorgang der Klasse aus
`review-slice-084` F-3) · „netzlose Suite blind gegen SQL-Prädikate" (F-4, neu).

## Verdikt

**Merge-blockierend:** **nein** — 0 HIGH, 0 MEDIUM. Die tragenden Zusagen des
Slices halten gegen jede eigene Messung: kein Verhaltens-Change (byte-identische
reale Tests, unveränderte Fehlertexte, unveränderter öffentlicher Rand), der
Null-Befund in allen drei Teilen **an beiden Enden der Range** nachgemessen, die
Nahtform wörtlich nach `ADR-0080`, der Digest form-korrekt und reproduziert. Die
fünfte Gegenprobe des Implementers ist reproduziert — die Hülle hängt am realen
Tier, nicht am Fake. Der eine LOW-Befund ist ein **Nachzug an einem lebenden
Träger**, kein Substanz-Fehler.

**Rückgabe-Pfeil an den Implementer: ja (klein).** F-1 (den Nenner in
`harness/sensors/db-adapter-coverage.md` von 659/155 auf **691/187** ziehen;
dabei die stehengebliebene 650 der Kalibrierungs-Zeile mitnehmen) — ein
Ein-Zeilen-Nachzug, der in denselben Lauf mit der noch offenen §7-Closure gehen
kann. **Kein** Rückweg zur Plan-Korrektur, **kein** Architect-Zug: F-1 ist eine
Zahl, kein Schnitt; die Verdünnungs-Frage (a) ist mit `ADR-0080` entschieden und
mit Trigger versehen. F-2/F-3/F-4 tragen keinen Pfeil.

**DoD-Häkchen:** Das Kästchen „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" wird **nicht** von mir nachgezogen — mit dem Rückgabe-Pfeil ist die
Bedingung des Skill-Nachzugs („keine Fixrunde nötig") nicht erfüllt; der
reguläre Nachzug läuft über den Implementer-Lauf (Schritt 21).

**Für die Closure §7 mitzunehmen:** (1) die vier Finding-Klassen oben; (2) der
Register-Beleg `evidence/slice-085.md` in
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (**2×**) — und mit F-1
zugleich `evidence/slice-085.md` in
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**3×**, Lese-Schritt
fällig); (3) die Anteilszahl der Verdünnung (7,10 % → 59,89 % in `receive`;
8,35 % → 26,69 % im Gegenstand) als Datum des `ADR-0080`-Triggers (a); (4) die
drei §6-Risiken mit Ausgang — Risiko 2/3/4 sind mit diesem Lauf **entfallen**
(kein Verhaltens-Change; die Fake-Seite prüft real, M1/M2 rot, F-4 benennt die
Grenze; kein Transfer), Risiko 1 war vorab gestrichen.

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt); er ersetzt **keine** Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11, anderer Eingabe-Kontext).
