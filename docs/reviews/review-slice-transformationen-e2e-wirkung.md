# Review-Report: slice-transformationen-e2e-wirkung — 2026-09-26

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec-Stellen und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-transformationen-e2e-wirkung` (Welle `welle-transformationen`), Diff-Range
`c246ba4f..HEAD` (7 Commits, 5 Dateien, +688/−56 einschließlich der drei Lifecycle-Commits `631748d2`, `c73d4292`,
`38c7b3bc`; Baum sauber, nicht gepusht). Inhalt: `test/integration/transformation_e2e_test.go` (neu, 321 Zeilen),
`tools/harness/run-integration-tests.sh` (+288/−2), `harness/README.md` (2 Zeilen), `docs/user/e2e-abdeckung.md`
(Erzeugnis), der Plan.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen
ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Das Edit-Werkzeug stand im
Lauf nicht zur Verfügung; alle Mutationen liefen als `sed … Datei > Kopie` bzw. `awk … > Kopie` im Scratchpad.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-transformationen-e2e-wirkung` (§1, §2 DoD-Wortlaut als Bezug, §3 Plan mit Abweichungen und
  Suchlauf-Feld, §6)
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Folgepflicht 5, §Fitness
  Function „Realer Rundlauf“), [`ADR-0030`](../plan/adr/0030-testpyramide.md) (E2E-Tier),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
  [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md) (Ausschluss nach Neustart)
- [`SPEC-019`](../../spec/pflichtenheft.md) (K1–K4, Fehlertexte)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../spec/lastenheft.md),
  [`LH-FA-CFG-005`](../../spec/lastenheft.md), [`LH-FA-SST-002`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md)
- `AGENTS.md` (§3.1, §3.2, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md` (`MR-000` bis `MR-003`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-transformationen-map-value.md`

**Eigene Messungen:**

- **Läufe** (Exit-Codes ungefiltert gesichert und einzeln gelesen): `make gates` am HEAD Exit 0; `make fmt-check`
  258 Go-Dateien formatiert; `make kommentar-kennungen DIFF=c246ba4f` Exit 0, kein Kandidat;
  `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-transformationen-e2e-wirkung.md` „16 Zeilen
  stimmen“; `RANGE=c246ba4f..HEAD make commit-traceability` OK für 7 Commits; `make doc-trace` am Arbeitsbaum:
  „80 Anforderung(en), 1 Waise(n)“, die Waise ist `LH-FA-CFG-008`.
- **Erzeugnis `docs/user/e2e-abdeckung.md`:** Ich habe den Erzeuger des Runners (Zeilen 95–217, Funktionen
  `abdeckung_declare`, `abdeckung_render`, `abdeckung_go_zeilen_lesen`, `abdeckung_schreiben`) in einer Kopie mit
  erhaltenen Zeilennummern gegen den HEAD-Quelltext laufen lassen (Go-Hälfte real über `go test -run
  '^TestAbdeckungstabelleZeilen$'`, Bash-Hälfte über die `abdeckung_declare`-Zeilen des Runners). Ergebnis: 16 Go-
  und 38 Bash-Zeilen, `cmp` gegen die committete Datei **byte-gleich**. Jede `Ort`-Angabe der 54 Zeilen zeigt auf
  eine `func TestE2E<Name>(`-Zeile bzw. auf eine Runner-Zeile, die den Anker enthält (alle 54 geprüft).
- **`-run`-Abdeckung:** `grep -rhoE '^func TestE2E…' test/integration` ergibt 16 Funktionen; jede steht in genau
  einem der vier E2E-`-run`-Muster des Runners (Zeilen 427, 3156, 3864, 3959), die zwei neuen im ersten.
- **Unmutierte Wiederholung gegen eine eigene Wegwerf-Umgebung:** Die Umgebung entstand aus den Zeilen 1–402 des
  committeten Runners (PostgreSQL, Schema-Rollout, NATS, Feed-Container, Compose-Namen wie im Runner, der
  Implementer-Stack lief nicht). Darin: beide Go-Tests `PASS` (1,5 s und 3,2 s); Phase „Transformationen-Happy-Path“
  (Zeilen 3464–3685 des Runners, mit `bf_*`-Hilfen aus 2727–2860) exit 0 in 17,5 s; Phase „Transformationen-Neustart
  und Ausschluss“ (3686–3746, mit den `TF_RULE_*`-Zeilen 3588–3589 davor) exit 0 in 13,7 s. Den vollen
  `make test-integration`-Lauf habe ich **nicht** gefahren; die Wechselwirkung der beiden Go-Tests (sie laufen im
  ersten `go test`-Aufruf am Anfang des Runners) mit den späteren Retention-/Backfill-Phasen ist damit von mir nicht
  nachgemessen.
- **Mutationen** (15, Tabelle unten): an Kopien im Scratchpad, je eigene Tabellennamen, gegen dieselbe
  Wegwerf-Umgebung; kein `sed -i`, kein Host-python, kein `cat >>` auf Repo-Dateien. Umgebung mit
  `docker compose down -v --remove-orphans` abgeräumt. Dangling Volumes vorher 34, nachher 34, kein prune;
  `git status` am Ende sauber.

| # | Ort | Mutation (Eingabeseite) | Ergebnis |
|---|---|---|---|
| G1 | Go, K3 zweite Form | `"to":"id"` → `"to":"idx"` | rot: „Antrag endete applied … erwartet failed“ |
| G2 | Go, Gegenprobe | `validSpec` `"to":"memo"` → `"to":"customer_name"` | rot: „Gegenprobe endete failed“ |
| G3 | Go | `setRule(… "status_lesbar" …)` entfernt | rot in beiden Tests (Alt-Bild, Change nach den Anträgen) |
| G5 | Go, K1 | Regelname `kundenname` → `kundenname2` | rot: „endete applied“ |
| G8 | Go | `REPLICA IDENTITY FULL` entfernt | rot: Alt-Bild leer |
| R1 | Runner Phase 1 | `tf_set … status_lesbar` entfernt | rot beim gRPC-Bild (`status:"o"`) |
| R2 | Runner Phase 1 | eingefügter Stream-Wert `o` → `x` | rot beim gRPC-Bild |
| R13 | Runner Phase 1 | Persistenz-Gleichheit gegen `change_id = 'nope'` | rot (`erwartet '1', gelesen '0'`) |
| R14 | Runner Phase 1 | dasselbe für `GET /changes` | rot |
| R15 | Runner Phase 1 | `tf_remove … kundenname` entfernt | rot bei „Rohform nach dem Entfernen“ |
| R3 | Runner Phase 2 | erster `docker restart` entfernt | rot: „Startzeit … blieb gleich“ |
| R4 | Runner Phase 2 | `tf_exclude … name` → `status` | rot bei „Bild mit Ausschluss und Regeln vor dem Neustart“ |
| R5 | Runner Phase 2 | `tf_no_value` gegen den vorhandenen Wert `open` | rot (`erwartet '0', gelesen '1'`) — die Gegenprobe des Wert-Prüfers beißt |
| R6 | Runner Phase 2 | erste `tf_remove … kundenname` entfernt | rot bei „Rohform nach dem Entfernen“ |
| R10 | Runner Phase 2 | **zweiter** `docker restart` entfernt | **grün, exit 0** — siehe F-2 |

## Findings

### F-1 — Runner-Kommentar nennt die Alternative im Konjunktiv

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7; Skill-HIGH „Kommentar trägt keine der Kommentar-Klassen“ (Skopus: Tests und Runner
  `tools/harness/*.sh`, Konjunktiv über die verworfene Alternative)
- `pfad`: `tools/harness/run-integration-tests.sh:3701-3703`
- `befund`: Der Kommentar vor dem ersten `docker restart` der Phase „Neustart und Ausschluss“ lautet: „Der Neustart ist
  ein echter Prozess-Neustart desselben Containers; ohne die Ableitung des Regelstands beim Prozessstart trüge die
  nächste Change die Rohform.“ Der zweite Halbsatz beschreibt einen Zustand, den es nicht gibt, im Konjunktiv. Die
  Kommentar-Klasse (Zusage, Kopplung, Abgrenzung, Rang-Zeiger, Grenze) trägt er nicht; was die Stelle zusagt, steht
  im ersten Halbsatz und im letzten Satz („Der Beleg des Neustarts ist die Startzeit des Containers“). Der Skill
  nennt genau diese Form („würde diese verzögern“) ausdrücklich unter HIGH, nicht unter INFO. Ein weiteres
  Konjunktiv-/Vorher-Nachher-Muster steht in den hinzugefügten Zeilen nicht (Suche über alle hinzugefügten
  Kommentarzeilen von Runner und Go-Datei).
- `verifizierbar`: nein (Lese-Handlung; `make kommentar-kennungen` prüft nur Kennungsform)
- `klasse`: Kommentar mit Konjunktiv über die verworfene Alternative (Runner)

### F-2 — Die Zusage „nach dem zweiten Neustart“ ist an den zweiten Neustart nicht gebunden

- `kategorie`: MEDIUM
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite“;
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5 und
  [`LH-QA-SEC-004`](../../spec/lastenheft.md) (Plan DoD Punkt 2: „vor und nach dem Neustart“)
- `pfad`: `tools/harness/run-integration-tests.sh:3731` (zweiter `docker restart`), `:3745` (Ausgabezeile) gegen
  `:3704-3707` (Startzeit-Vergleich des ersten)
- `befund`: Nur der erste `docker restart` trägt den Beleg „Startzeit davor ≠ Startzeit danach“. Vor dem zweiten steht
  kein Vergleich. Gemessen (R10): Entfällt die Zeile `docker restart "$FEED_CONTAINER"` an Stelle 3731 (das
  folgende `bf_await_healthy` bleibt stehen), endet die Phase mit exit 0 und druckt weiter „…vor und nach dem
  zweiten Neustart weder name noch customer_name noch den Wert“. Die Aussage „Ausschluss **und** Regel überleben den
  Neustart“ hängt an genau diesem zweiten Neustart: die Kombination aus Ausschluss und `rename_column`-Regel wird nur
  hier nach einem Neustart gelesen (der Beleg zu `ADR-0065` trägt den Ausschluss ohne Regel).
- `verifizierbar`: ja — Mutation wie beschrieben (`awk` streicht die zweite Zeile `docker restart "$FEED_CONTAINER"`)
- `klasse`: Zusage ohne Bindung an die Eingabeseite (zweiter Neustart)
- Einstufung: Der Skill führt die Klasse unter HIGH. Ich stufe MEDIUM ein, weil der Aufruf wörtlich im Skript steht und
  bei einem Fehlschlag unter `set -e` endet; die Lücke ist ein fehlender zweiter Startzeit-Vergleich, kein
  Falsch-Grün einer vorhandenen Assertion. Wer die Skill-Regel wörtlich anwendet, hebt auf HIGH; das Verdikt ändert
  sich dadurch nicht, weil F-1 die Fixrunde ohnehin auslöst.

### F-3 — Der Go-Doc-Kommentar sagt „in genau einem Feld“, die K4-Form `remove` weicht in mehreren ab

- `kategorie`: LOW
- `quelle`: `ADR-0083` (Aussage breiter als ihre Messung)
- `pfad`: `test/integration/transformation_e2e_test.go:247-249`, Fall `{"K4 Regelname nicht geführt", …}` `:284-285`
- `befund`: „Jeder Negativfall weicht in genau einem Feld von der Gegenprobe ab, die `applied` endet.“ Fünf der sechs
  Fälle sind `set_transformation`-Anträge, die in Regelname, Spalte oder Zielname von `validSpec` abweichen (G1, G2,
  G5 färben rot). Der sechste ist ein `remove_transformation`-Antrag: andere Antragsart, kein `rule_spec`. Seine
  Gegenprobe (Entfernen einer gesetzten Regel endet `applied`) läuft nur im `t.Cleanup` von `removeRulesAtEnd` und
  wird dort nur bei nicht fehlgeschlagenem Test als Fehler gemeldet. Die Aussage trägt für die Set-Fälle, für den
  Remove-Fall nur indirekt.
- `verifizierbar`: nein (Lese-Handlung)
- `klasse`: Test-Aussage breiter als die Messung

### F-4 — Das Fenster READY → erste Zeile trägt die Timeouts der drei Clients, gemessen nur auf dem Entwicklerhost

- `kategorie`: LOW
- `quelle`: Maintainability (Flake-Risiko); Plan §6 nennt „Neustart-Beleg misst Timing“, dieses Fenster nicht
- `pfad`: `tools/harness/run-integration-tests.sh:3610-3612`, `:3620-3634`; `tools/harness/natsstreamsub/main.go:93`
  (`NextMsg(30 * time.Second)`), `tools/harness/grpcclient/main.go:52` und `tools/harness/sseclient/main.go:41`
  (je 60 s)
- `befund`: Die drei Clients starten gleichzeitig; ihr READY wird nacheinander abgewartet (bis 90 s je Client), erst
  danach entsteht die erste Zeile. Der NATS-Client wartet ab seinem READY 30 s, gRPC und SSE 60 s ab Stream-Öffnung.
  Das Fenster ist die Differenz der Übersetzungszeiten der drei `go run`-Aufrufe (kalter `GOCACHE` je Container).
  Gemessen auf dem Host dieses Reviews (20 Kerne): `go build` kalt 5,30 s (`natsstreamsub`), 4,85 s (`sseclient`),
  6,15 s (`grpcclient`); im unmutierten Lauf sowie in R1 und R2 (Zeilen-Id 11) genügte der erste Einfüge-Versuch. Auf dem
  Zwei-Kern-Runner von `e2e.yml` ist die Differenz nicht gemessen. Die Wiederholung bis sechs Zeilen fängt eine
  späte Registrierung am Broadcaster, nicht einen abgelaufenen Client-Timeout.
- `verifizierbar`: nein (Messung an der Pipeline `e2e.yml` steht aus)
- `klasse`: Fenster-Kopplung an Client-Timeouts

### F-5 — Der Beleg-Anker im Suchlauf-Feld nennt einen nicht committeten Bericht

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-e2e-wirkung.md:161` (Suchlauf-Zeile „Jede neue
  `TestE2E*`-Funktion …“), `:160` („gedruckt vom Runner“)
- `befund`: „Der Lauf von `make test-integration` fuhr beide (`--- PASS`, gedruckt, siehe Bericht)“ nennt als Anker den
  Bericht des Implementers, keinen committeten Träger. Die Zahlen daneben sind an beiden Ständen nachmessbar und
  stimmen (Stat 651/40, 50 → 54 Tabellenzeilen, 16/38). Der Lauf selbst ist im Repo nicht belegt; die DoD-Haken 1 bis
  3 stehen mit `[x]` auf „zu belegen durch: ein realer, grüner `make test-integration`-Lauf“.
- `verifizierbar`: nein
- `klasse`: Beleg-Anker nicht committet

### F-6 — Fremder Träger `welle-transformationen.md:168` schreibt diesem Slice das bereits Erreichte zu; die Meldung nennt keine Frist

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger in einer fremden Datei: melden, Frist nennen)
- `pfad`: `docs/plan/planning/welle-transformationen.md:168` (Tabellenzeile des Slice), Plan §3 Zeile „Waisen-Aussage im
  RTM-Träger“
- `befund`: Die Zeile führt als Wirkung dieses Slice „`LH-FA-CFG-007` verlässt die Waisen“. Der Befund des Plans stimmt:
  `git show fc0b8d38:docs/user/e2e-abdeckung.md | grep -c 'CFG-007'` ergibt 1 (die Zeile der Phase
  „Backfill-Regelstand“); `make doc-trace` am Parent-Stand habe ich nicht gefahren, am HEAD führt es
  `LH-FA-CFG-008` als einzige Waise. Die Sätze in Zeile 53 und 124–127 beschreiben den Zustand des ganzen Bündels und
  bleiben wahr. Das committete Feld sagt „gemeldet (Bericht)“ und nennt weder Frist noch Adresse (die Regel: Closure
  des meldenden Slice, Planner zieht nach oder benennt den Träger mit Adresse); die Meldung liegt nur im Bericht.
- `verifizierbar`: ja — `git grep -n 'verlässt die Waisen' -- docs/plan/planning/welle-transformationen.md`
- `klasse`: Träger-Nachzug (fremde Datei, Meldung ohne Frist)

### F-7 — Alt-Bild und UPDATE/DELETE sind nur über `cdc.changes` belegt

- `kategorie`: INFO
- `quelle`: `ADR-0112` Folgepflicht 5 (Zustellwege), Plan §3
- `pfad`: `tools/harness/grpcclient/main.go:72`, `tools/harness/sseclient/main.go:73`,
  `tools/harness/natsstreamsub/main.go:105` (Ausgabe nur `new_image`)
- `befund`: Die vier Wegwerf-Clients drucken das Neu-Bild; die Phase belegt daher die Form auf den fünf Wegen für
  `INSERT`. Beide Bilder von INSERT, UPDATE und DELETE trägt allein der Go-Test über `cdc.changes`. Die README-Zeile
  ordnet das richtig zu („Bilder aller Operationen“ hängt an `TestE2ETransformationRulesShapeBothImages`). Kein Fehler,
  eine Grenze der Aussage „alle Wege dieselbe Form“.
- `verifizierbar`: ja
- `klasse`: Aussagegrenze

### F-8 — Phase 2 liest Variablen, die Phase 1 definiert

- `kategorie`: INFO
- `quelle`: Maintainability; Plan §6 „Geteilter Zustand zwischen Phasen“
- `pfad`: `tools/harness/run-integration-tests.sh:3588-3589` (`TF_RULE_RENAME`, `TF_RULE_MAP` in Phase 1) gegen
  `:3699-3730` (Verwendung in Phase 2)
- `befund`: Die Regelformen stehen im Block der ersten Phase und werden von der zweiten gelesen; kein Kommentar nennt
  die Kopplung. Fällt Phase 1 weg oder wandert sie, endet Phase 2 unter `set -u` laut mit „unbound variable“, nicht
  still. Die Tabellen der beiden Phasen sind getrennt (`feed_e2e_transform`, `feed_e2e_transform_restart`), der
  Regelstand ist es auch.
- `verifizierbar`: ja (Phase 2 allein ohne die zwei Zeilen läuft nicht)
- `klasse`: nicht benannte Kopplung zwischen Phasen

### F-9 — Die Tabellenzeilen wurden durch Anpassung des Quelltexts, nicht durch einen Lauf in Einklang gebracht

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Erzeugnis), Maintainability
- `pfad`: Commit `66f60c8b` (`test/integration/transformation_e2e_test.go`, −1 Zeile im Kommentar)
- `befund`: In `87281d08` und `21a34696` trägt das Erzeugnis `…transformation_e2e_test.go:162` und `:250`, die Funktionen
  standen an den Zeilen 163 und 251 (`git show 87281d08:…`). `66f60c8b` kürzt den Doc-Kommentar von
  `removeRulesAtEnd` um eine Zeile, damit die Zeilen stimmen. Am HEAD ist das Erzeugnis byte-gleich dem, was der
  Runner schreiben würde (Messung oben); für Bisect-Läufe über die Zwischenstände gilt das nicht. Die Änderung ist
  rein ein Kommentar.
- `verifizierbar`: ja
- `klasse`: Erzeugnis an Zwischenständen nicht im Einklang

### F-10 — Die Abweichung „Boundary in Go statt im Runner“ nennt ihren Grund nicht

- `kategorie`: INFO
- `quelle`: Plan §3, `ADR-0030`
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-e2e-wirkung.md:145-146`
- `befund`: Die Zeile zu `integration_test.go` trägt „nicht realisiert“ samt Grund. Die Verschiebung der K1–K4-Boundary
  aus dem Runner in die Go-Datei steht als Satz in der Runner-Zeile, ohne Grund. Die Abweichung ist redlich:
  `transformation_e2e_test.go` importiert nur Standardbibliothek, spricht das laufende System über pgx an
  (`cdc.set_transformation`, `cdc.remove_transformation`, `cdc.administration_request`, `cdc.changes`), vergleicht die
  Texte exakt (`message != violation.want`) gegen die sechs Zeilen von [`SPEC-019`](../../spec/pflichtenheft.md)
  (`spec/pflichtenheft.md:753-758`) und läuft im ersten `go test`-Aufruf des Runners. Kein Go-Import interner
  Anwendungslogik, damit E2E-Tier im Sinn von [`ADR-0030`](../plan/adr/0030-testpyramide.md).
- `verifizierbar`: ja
- `klasse`: Abweichung ohne Grund im Plan-Feld

## Negativbefunde

- geprüft, ohne Befund: **Wirkung auf den fünf Wegen** (Schwerpunkt 1). Jeder Weg wird gegen das persistierte Bild
  **derselben `change_id`** gehalten (`tf_expect_form`: Form **und** `new_data = <Bild>` an `change_id`; R13 färbt
  rot). `GET /changes` liest den ganzen Bereich, zählt 1 + n + 2 Zeilen, vergleicht die `change_id`-Menge gegen
  `cdc.changes` und jedes Bild gegen die Zeile (R14 rot). Ein älteres Roh-Bild kann keinen Weg grün färben: die
  Stream-Clients haben kein Replay, die Zeilen entstehen nach dem `applied` der Regeln, und die Form-Prüfung verlangt
  `customer_name`/`open` (R1 und R2 rot). Die Wiederholung „bis zu sechs gleichförmige Zeilen“ ist an die Eingabeseite
  gebunden: jede eingefügte Zeile trägt denselben Quellwert, und die Gesamtzahl der Zeilen mit der Form muss gleich
  `tf_stream_rows` sein.
- geprüft, ohne Befund: **Neustart und Ausschluss** (Schwerpunkt 2), soweit gebunden: Startzeit-Vergleich des ersten
  Neustarts (R3 rot), Regelstand nach dem Neustart, Rücknahme (R6 rot), erneutes Setzen desselben Regelnamens nach
  dem Entfernen (K1 schlägt nicht an), Ausschluss der Spalte mit Regel (R4 rot), Schlüsselmenge `{"status":"open"}`
  schließt Quell- und Zielnamen aus, `tf_no_value` beißt (R5). Rückräumen: beide Phasen nehmen ihre Regeln am Ende
  zurück, beide Go-Tests im `t.Cleanup`; bei einem Fehlschlag endet der Runner (`bf_fail` → `exit 1`, Trap räumt die
  Compose-Umgebung ab), ein Phasen-Übersprung mit stehendem Regelstand ist damit nicht möglich. Lücke: F-2.
- geprüft, ohne Befund: **Go-Tests K1–K4** (Schwerpunkt 3): beide K3-Formen, beide K4-Formen, Klartext samt Adresse
  gegen `spec/pflichtenheft.md:753-758` wörtlich gleich, Regelstand nach den abgelehnten Anträgen unverändert (Change
  trägt `customer_name`/`open`/`note` roh), Gegenprobe `applied` mit Wirkung ab der nächsten Change; G1, G2, G3, G5,
  G8 färben rot. `null`-Werte im Bild würden als leerer Wert unter dem Schlüssel erscheinen und die
  `DeepEqual`-Erwartung brechen.
- geprüft, ohne Befund: **Abdeckungs-Erzeuger** (Schwerpunkt 3/4): `abdeckungsZeilen` liest jede `.go`-Datei des
  Verzeichnisses, die neuen Funktionen tragen `LH-FA-CFG-007` genau einmal im Doc-Kommentar, beide Bash-Phasen tragen
  ihn im Anker; Erzeugnis byte-gleich (Messung oben), 16/38 Zeilen stimmen mit der Ausgabe des Runners
  („E2E-Abdeckungstabelle aus 16 Go-Zeilen und 38 Bash-Zeilen“) überein.
- geprüft, ohne Befund: **Träger-Nachzüge** (Schwerpunkt 4): `harness/README.md` Zeile `make test-integration` nennt
  vier Belege in der Reihenfolge Happy Path, Bilder aller Operationen, Boundary, Neustart mit Ausschluss und stimmt mit
  dem Diff (zwei Phasen, zwei Testfunktionen); die Zählung „vier“ ist die der Aufzählung. Herkunfts-Anker
  `· seit slice-transformationen-e2e-wirkung` ein Feld. Zeile `make doc-trace`: „gemessen 2026-09-26 … 80
  Anforderungen, **1 Waise** — `LH-FA-CFG-008`“ trägt Datum, Werkzeug und Stand und ist von mir nachgemessen. Kopfkommentar
  des Runners nennt die beiden Rundläufe. Handbuch unberührt, benannter Aufschub `slice-transformationen-betriebsdoku`
  steht im Plan (§1, §6); keine neue Betreiber-Oberfläche, keine `CDC_*`-Variable.
- geprüft, ohne Befund: **Suchlauf-Feld** (§3.13): alle 16 Zeilen an beiden Ständen gleich; die Trefferaufschlüsselung
  im Plan stimmt (12 Treffer des Phasen-Symbols in drei Dateien, 19 → 21 `func TestE2E`, 1 → 5 Zeilen
  `LH-FA-CFG-007` im Erzeugnis, 7/7 `new_image`). Suchraum ganzer Baum ohne Records, Muster aus Symbol, Zählwort und
  Beschreibung. Eigene Gegensuche nach Trägern der Phasenliste (`Leerlauf-Bestätigung`, `Waise`, `Bash-Zeilen`,
  `e2e-abdeckung`) findet nichts Weiteres außerhalb von Records; einziger Träger außerhalb des Feldes ist F-6.
- geprüft, ohne Befund: **Kommentare §3.7** (Schwerpunkt 5): `make kommentar-kennungen DIFF=c246ba4f` ohne Kandidat,
  `make fmt-check` 258 Dateien formatiert; ein Herkunftsfeld je Go-Kommentar (`LH-FA-CFG-007`), keine Kette, keine
  Kompaktform, kein „ff.“; keine Spec-Wiedergabe über das hinaus, was der Test prüft. Einziges Konjunktiv-Muster:
  F-1.
- geprüft, ohne Befund: **Zahlen und Ursprung** (Schwerpunkt 6): die im Repo genannten Zahlen (80 Anforderungen, 1
  Waise, 16/38/54 Zeilen, Stat 651/40, 19/21, 12/12) tragen Ursprung oder sind nachgemessen und stimmen an beiden
  Ständen. Die im Auftrag genannten Größen „349 s“ und „85,30 %“ stehen nicht im Diff und wurden nicht geprüft. Die
  beiden Phasen messe ich mit 17,5 s und 13,7 s, die zwei Go-Tests mit 4,7 s (Wegwerf-Umgebung, nicht der volle
  Lauf); die Laufzeit-Zunahme von `make test-integration` liegt damit unter einer Minute, und `e2e.yml` trägt
  `timeout-minutes: 60`.
- geprüft, ohne Befund: **Docker-only und Erzeugnis-Schreiben** (Schwerpunkt 7): die Skript-Änderung ruft nur `docker`,
  `grep`, `sed` (ohne `-i`, nur Filter auf stdin), `cut`, `seq`, `sleep`, `awk` (lesend) und `date` auf; keine
  Host-Toolchain, keine Umleitung auf Repo-Dateien. Das Erzeugnis schreibt weiter allein der bestehende Weg
  `abdeckung_schreiben` (Temp-Datei, `cmp`, `mv`, nur bei Abweichung); der Diff ändert ihn nicht. Kein `//nolint`,
  keine Coverage-Ausnahme (§3.2).
- geprüft, ohne Befund: **Traceability/Lifecycle**: alle sieben Betreffe tragen `LH-*`/`ADR-*`, keine
  `SPEC-*`/`ARC-*` im Betreff (`make commit-traceability` OK); `631748d2` und `38c7b3bc` sind reine Renames (0
  Zeilen), die Inhaltsänderung `c73d4292` liegt in einem eigenen Commit dazwischen (§3.3).
- geprüft, ohne Befund: **Handbuch-Regeln des Skills**: `docs/user/benutzerhandbuch.md` im Diff unberührt, damit weder
  Versionshistorie noch Handbuch-Zug fällig.
- geprüft, ohne Befund: Produktivcode: keiner im Diff (`internal/`, `cmd/`, `tools/schema`, `spec/`, `.github/`
  unberührt); `.github/workflows/e2e.yml` fährt das Paket unverändert, §3.10 greift nicht.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Kommentar mit Konjunktiv über die verworfene Alternative (Runner) · Zusage ohne
Bindung an die Eingabeseite (zweiter Neustart) · Test-Aussage breiter als die Messung · Fenster-Kopplung an
Client-Timeouts · Beleg-Anker nicht committet · Träger-Nachzug (fremde Datei, Meldung ohne Frist) · Aussagegrenze ·
nicht benannte Kopplung zwischen Phasen · Erzeugnis an Zwischenständen nicht im Einklang · Abweichung ohne Grund im
Plan-Feld

## Verdikt

**Merge-blockierend:** nein für die Wirkung — die Belege binden an ihre Eingabeseite (14 von 15 Mutationen färben
rot), beide Go-Tests und beide Runner-Phasen liefen in meiner Wegwerf-Umgebung grün, kein Produktivcode im Diff. Der
Skill klassifiziert F-1 als HIGH; damit ist eine kleine Fixrunde am Implementer nötig: F-1 (ein Kommentar, drei
Zeilen), F-2 (ein zweiter Startzeit-Vergleich oder eine ehrliche Zusage in der Ausgabezeile), F-3 bis F-5 nach dessen
Ermessen. Die DoD-Zeile „Review durchgeführt“ bleibt offen und wird bei Schritt 21 des Implementer-Ablaufs
nachgezogen.

**Übergabe:** Findings an den Implementer. F-5 zusätzlich an den Verifier (der reale, grüne
`make test-integration`-Lauf ist von mir nicht wiederholt), F-6 an den Planner (Träger `welle-transformationen.md:168`,
Frist: Closure dieses Slice), F-4 an den Planner für `e2e.yml`-Beobachtung beim nächsten realen Lauf. Die
Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report ist ein Lauf-Beleg und ersetzt
keine Verifikation.
