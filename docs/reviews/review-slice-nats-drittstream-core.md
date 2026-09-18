# Review — `slice-nats-drittstream-core` (Commit `64a138b`)

**Review-Art:** Code und Doku gegen Plan, Entscheidung und Hard Rules
(`.harness/skills/reviewer.md`, Modul 10).

**Gegenstand:** `git show 64a138b` (17 Dateien, +1190/−109) gegen
`docs/plan/planning/in-progress/slice-nats-drittstream-core.md` (Plan) und
[`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(`Accepted`).

**Datum:** 2026-09-18.

**Selbst ausgeführt (Exit-Code jeweils direkt und ungepiped, `AGENTS.md`
§3.9):**

- `make gates` — Exit 0 (baseline-verify OK · docs-check 0 Befunde über 713
  Dateien · coverage-gate 82,70 % ≥ 80 % · commit-traceability OK · a-check
  0 Befunde · generated-sync OK).
- `make test-integration` — Exit 0, frischer Lauf; die neue Phase lief real
  durch (`RECEIVED change_id=980-1 table=feed_e2e_full operation=INSERT
  new_image={"id":"281",…,"NatsStreamE2ESentinel"}`, vorausgehend `REJECTED`
  mit falschem Token). Arbeitsbaum danach unverändert, keine
  `cdc`-Container zurückgelassen.
- `nats.go v1.53.1` am Modul-Cache-Pin aufgeschlagen (`connectProto()`), weil
  der Review eine Verhaltensaussage über die Bibliothek prüft (F-1).

**Verdikt:** HIGH-Findings offen → Fixrunde nötig. Die DoD-Checkbox „Review
durchgeführt" ist **nicht** nachgezogen, der Slice bleibt in `in-progress/`.

---

## HIGH

### F-1 — Kommentar behauptet eine Credential-Präzedenz, die die Bibliothek umgekehrt implementiert

- **klasse:** Tatsachenbehauptung gegen die Messung driftend (`AGENTS.md`
  §3.12 Instanz B, §3.7); Nachbar-Präzedenz
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
- **pfad:** `examples/.env:37-38` (Kommentar), `examples/.env:40` (URL)
- **befund:** Der Kommentar sagt, „ein expliziter Verbindungs-Options-Wert
  schlägt dabei einen eingebetteten URL-Token". Gemessen an `nats.go
  v1.53.1`: `connectProto()` liest `o.Token` nur im `else`-Zweig
  (`nc.current.URL.User == nil`); trägt die URL Userinfo, gewinnt **immer**
  die URL, die Option wird gar nicht gelesen. Die Bibliothek dokumentiert
  das selbst (`Token()` „… when a token is not included directly in the
  URLs"). Die Aussage ist invertiert; der Demo-Fall funktioniert nur, weil
  beide Werte identisch sind — der Defekt ist die falsche Regel im Text.
  Dieselbe Formulierung trägt die Begründung des §6-Risikos und wäre
  Vorbild für den Folge-Slice.
- **verifizierbar:** nein — kein Gate liest Kommentar-Prosa; die Probe ist
  das Aufschlagen des Pins.

### F-2 — §3.13: der Zählwert der Zugangsdaten-Klasse ist im selben, vom Diff geänderten File nicht nachgezogen

- **klasse:** Arbeit überholt stehenden Träger (`AGENTS.md` §3.13)
- **pfad:** `internal/bootstrap/config_file.go:59`
- **befund:** Der Diff zieht den Kommentar über
  `forbiddenFileCredentialKeys` auf „die **drei** Token-Schlüssel" (:34) und
  die Liste auf sieben Einträge (:43-46), lässt aber 15 Zeilen darunter den
  `fileConfig`-Kommentar unverändert auf „Die DSNs, die **zwei**
  Token-Schlüssel und `nats_url` bleiben env-var-exklusiv" stehen — die
  Datei widerspricht sich nach dem Commit selbst. Zugleich falsifiziert das
  die Plan-Aussage (§6), der Suchlauf habe „alle drei lebenden Träger … alle
  drei in diesem Diff nachgezogen"; das Suchmuster „zwei Token-Schlüssel"
  trifft genau diese Stelle.
- **verifizierbar:** nein — kein Gate prüft eine Freitext-Zahl über Träger
  hinweg.

### F-3 — §3.13: zweite bewegte Eigenschaft („was `make test-integration` prüft") ohne Suchlauf

- **klasse:** Arbeit überholt stehenden Träger (`AGENTS.md` §3.13)
- **pfad:** `harness/README.md:135` (`make test-integration`-Zeile; nicht im
  Diff)
- **befund:** Die Zeile enumeriert jede Rundlauf-Phase mit `· seit
  slice-<NNN>` und wurde bei jedem vergleichbaren Rundlauf im selben Commit
  mitgezogen — belegt: SSE-Rundlauf `89d31d1`, gRPC-Rundlauf `b835dde`,
  `GET /changes`-Rundlauf `93ac92e`. Dieser Commit fügt eine vollständige
  neue Phase hinzu und lässt die Zeile unverändert; sie kennt weder
  `natsstreamsub` noch `cdc.stream.>`.
- **verifizierbar:** nein — Lesen/`grep` über die Träger, kein Sensor.

### F-4 — §3.13-Meldung nennt die superseded Stelle und lässt den heute tragenden Träger aus

- **klasse:** Zitat nennt die falsche Stelle; §3.13-Meldung unvollständig
- **pfad:** `docs/plan/planning/in-progress/slice-nats-drittstream-core.md`
  §6; gemessene Träger:
  `docs/plan/adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md:80,129,165,179,209`,
  `docs/plan/adr/0091-zugangsdaten-klasse-sechs-schluessel.md:201`
- **befund:** Die Meldung nennt als „zusätzlichen Fund" nur `ADR-0091` und
  verweist für den Wächter auf `ADR-0089`. `ADR-0089` ist durch
  [**`ADR-0092`**](../plan/adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  (Accepted, „Supersedes [`ADR-0089`](../plan/adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)") ersetzt; der heute tragende Träger ist
  `ADR-0092`, und er sagt in seiner aktiven Fitness-Function-Zeile genau die
  jetzt falsche Sechs. Dieselbe Zahl steht bei `ADR-0092:165` und `:80`.
  `ADR-0091` führt zusätzlich eine zweite veraltete Zeile (:201; der Test
  iteriert inzwischen sieben). Folge: die vorgeschlagene Remediation
  „Folge-ADR (`Supersedes ADR-0091`)" würde `ADR-0092`s falsche
  Fitness-Function-Zeile stehen lassen. Vierte Spur derselben Klasse:
  `docs/plan/adr/README.md:106` („sechs Schlüssel").
- **verifizierbar:** nein — `grep` über die Träger; die Probe ist das
  Nachzählen.

---

## MEDIUM

### F-5 — Der Retry-Bereich des neuen Rundlaufs kollidiert mit dem Wertebereich einer früheren Phase

- **klasse:** Zähl-/ID-Bereich gegen fremden Testfall-Bereich driftend
- **pfad:** `tools/harness/run-integration-tests.sh:2440-2441` (Schleife,
  IDs 281…285) gegen `:1989`/`:1992` (`HTTP_READ_ID=285`, INSERT in
  `public.feed_e2e_full`), Kommentare `:2389-2390` und `:2437`
- **befund:** `HTTP_READ_ID=285` wird in der HTTP-Phase real in dieselbe
  Tabelle eingefügt; der späteste Retry-Versuch fügt erneut `id=285` ein.
  Unter `set -euo pipefail` beendet der primär-schlüsselverletzende
  `psql`-Aufruf den gesamten Runner mit einem rohen Duplikat-Fehler, statt
  der beabsichtigten Diagnose. Beide Kommentare behaupten das Gegenteil.
  Gemessen sind 281–284 frei, 285 nicht.
- **verifizierbar:** teilweise — mechanisch prüfbar (`grep` über die
  `INSERT … VALUES`-Zeilen), kein Gate.

### F-6 — Der neue Start-Fehlervertrag ist an keinem seiner Aufruforte gebunden

- **klasse:** Zusage ohne Bindung an ihre Eingabeseite
- **pfad:** `internal/bootstrap/wiring.go:310-314`,
  `internal/bootstrap/config_file.go:224-226`; Test nur
  `internal/bootstrap/changestream_internal_test.go:73`
- **befund:** `validateNatsStreamTokenRequiresURL` ist als reine Funktion
  vollständig getestet. Kein Test bindet jedoch die beiden Aufruforte:
  entfernt man den Aufruf aus `ConfigFromEnv` **oder** aus `mergeConfig`,
  bleibt `make test` grün — obwohl die Handbuch-Aussage „nur
  `CDC_NATS_STREAM_TOKEN` gesetzt … startet der Container nicht" genau an
  diesen Aufruforten hängt.
- **verifizierbar:** ja — Mutieren des Aufruforts färbt keinen Test rot.

---

## LOW

### F-7 — Der Leerwert-Teil der Vorbild-Validierung fehlt

- **klasse:** Abgrenzung weicht vom benannten Vorbild ab, ohne Begründung
- **pfad:** `internal/adapters/driven/natsstream/publisher.go:212-219` gegen
  `internal/adapters/driven/natsnotify/notify.go:123-131`
- **befund:** Der Reserved-Zeichen-Guard ist byte-gleich; die Divergenz liegt
  woanders: `natsnotify.Notify` lehnt zusätzlich `sourceID == ""` sowie
  `schema == "" || table == ""` ab, `natsstream.publish` prüft nur die
  beiden Namen je Token. Ein Change mit leerem Schema/Tabelle erzeugte ein
  verkürztes Subjekt statt übersprungen zu werden — erreichbar nur, wenn der
  Mapper einen leeren Relationsnamen liefert (praktisch nicht).
- **verifizierbar:** ja.

---

## INFO

- **F-8** — `ADR-0100` nennt für den Empfangsbeleg das Target `make test` —
  eine netzgebundene Regel, die `make test`s `--network none` strukturell
  nicht tragen kann; geliefert ist der Beleg über `make test-integration`.
  Transparent, aber die ADR-Zeile bleibt unscharf. → Verifier/Architect.
- **F-9** — Plan-§3-Zeile `test/integration/integration_test.go | update`
  wurde nicht geliefert; die DoD trägt den Rundlauf ausdrücklich über
  `compose.yaml` + Belegträger, geliefert ist genau das. Die §3-Zeile ist
  damit überholt, ohne dass die Abweichung notiert wurde.
- **F-10** — Das §6-Risiko „Wecksignal-Demo bricht" ist quellgestützt
  plausibel (Go-Client: URL-Userinfo ohne Passwort → Token), für die
  C#-/Kotlin-Wecksignal-Clients aber **nicht** belegt (`make example-run-*`
  ist kein Gate, `make examples-*` fährt keinen echten NATS). Das §6-Risiko
  bleibt richtigerweise weiter offen; die Erwartung „real gegen beide
  Demo-Clients geprüft" ist noch nicht eingelöst.
- **F-11** — Die DoD-Satzhälfte „ohne Token abgelehnt" ist nicht ausgeführt;
  geprüft wird ausschließlich der **falsche** Token. Die Zusicherung ist
  wahr (Core-NATS-`--auth`), die Hälfte aber unbelegt. → Verifier.

---

## Negativbefunde (geprüft, ohne Befund)

- `internal/adapters/driven/natsstream/` — kein Zeilen-Diff an `Broadcaster`,
  `ChangeStreamPort`, `CaptureService`; der Publish-Pfad kann nicht in den
  kritischen Pfad propagieren (`Broadcaster.Publish` ist
  nicht-blockierend/Drop-Newest, `conn.Publish` blockiert nicht). Der
  Broadcaster **schließt den Empfänger-Kanal bewusst nie**
  (`broadcaster.go:61-65`), der `for/select` kann also nicht heiß laufen;
  `if change == nil { continue }` ist defensiver als die gRPC-/SSE-Konsumenten.
- `internal/bootstrap/` — der Publisher entsteht genau unter beiden
  Bedingungen; `changeStreamEnabled` kann nicht `true` sein, während der
  Broadcaster nil bleibt; die Abschaltreihenfolge deadlockt nicht und lässt
  keine Goroutine zurück; `validateNatsStreamTokenRequiresURL` hängt auf
  beiden Ladepfaden.
- `compose.yaml`, `examples/compose.yaml` — `--auth`-Werte und Token-Env
  konsistent; der 8222-Monitor-Healthcheck ist unberührt; `make test-notify`
  startet seinen eigenen NATS-Container ohne `--auth` und ist unberührt.
- `tools/harness/natsstreamsub/`, `tools/harness/natssub/` — echte Clients
  ohne Mock; die `RECEIVED`-Assertion prüft Inhalt (Tabelle, Operation,
  Spaltenwert im `new_image`), die `cdc.changes`-Gegenprobe ist
  aussagekräftig, die Phase liegt wie kommentiert vor dem Container-Tausch.
- `spec/pflichtenheft.md`, `docs/user/benutzerhandbuch.md`,
  `docs/user/e2e-abdeckung.md` — die Sieben-Schlüssel-Zahl ist über die drei
  lebenden Träger konsistent; `SPEC-024` deckt sich mit dem Code; Handbuch
  `Version: 1.27 → 1.28` mit neuer Historie-Zeile.
- `AGENTS.md` §3.1/§3.2/§3.7 auf dem neuen Produktionscode — kein
  host-lokales Toolchain, kein `//nolint`, keine Chronik in
  Produktionskommentaren.
- Traceability — Commit-Betreff trägt `LH-FA-SST-008` und `ADR-0100`;
  `make commit-traceability` OK.
