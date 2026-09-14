# Welle 19 — Live-Change-Streaming — gRPC-Server-Streaming und HTTP/SSE — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-NNN` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-19
**Abschluss:** 2026-09-15
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-SST-008`](../../../../spec/lastenheft.md) trägt beide Zustellwege.
  gRPC-Server-Streaming ([`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md))
  und HTTP/Server-Sent-Events ([`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md))
  sind unabhängig funktionsfähig, gespeist von **einem** In-Prozess-`Broadcaster`
  hinter dem neuen Outbound Port `ChangeStreamPort`
  (`internal/application/port/outbound/changestream.go`); die
  Anforderung war bis dahin vollständig unimplementiert.
- `slice-069` — Adapter-Grundgerüst: `Broadcaster`
  (`internal/adapters/driven/grpcstream/`), gRPC-Driving-Adapter samt
  Auth-Interceptor (`internal/adapters/driving/grpc/`), Docker-only
  Protobuf-/buf-Codegen-Stufe (`make proto-generate`), additive
  `CDC_GRPC_ADDR`-Verdrahtung. Neuer
  [`SPEC-020`](../../../../spec/pflichtenheft.md) (Protobuf-/Nachrichtenschema,
  RPC-Methode, Stream-Semantik). Zustellsemantik aus
  [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md):
  `Publish` übergibt nicht-blockierend an eine begrenzte
  Empfangs-Warteschlange je Abonnent (Drop-Newest), der Erzeuger hält nie an.
- `slice-070` — Capture-Integration: `CaptureService.WithChangeStream` und
  die Publish-Kette `… → ACK Source → Notify (best effort) → Stream-Publish
  (best effort)`, fehlerisoliert; der Capture-kritische Pfad
  ([`ADR-0011`](../../adr/0011-persist-before-ack.md)) bleibt unberührt.
  [`ADR-0067`](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
  korrigiert die dritte Fitness-Function-Zeile von
  [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md).
- `slice-071` — gRPC-Beispiel-Client (`tools/harness/grpcclient/`) und realer
  E2E-Rundlauf gegen den laufenden Feed-Container: `READY` →
  `RECEIVED change_id=964-1` → `REJECTED code=Unauthenticated` ohne Token.
  [`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  trägt die begleitende Änderung der Maschinenform (`.a-check.yml`: Gruppe
  `tooling` + begrenzte Import-Kante).
- `slice-072` — HTTP/SSE-Endpunkt `GET /changes/stream` und realer
  E2E-Rundlauf: `READY` → `RECEIVED change_id=967-1` → `REJECTED code=401`.
  Bootstrap-Entkopplung (`changeStreamEnabled`: `CDC_GRPC_ADDR` **oder**
  `CDC_HTTP_ADDR`; beide Adressen leer → unverändertes Bestandsverhalten),
  Fire-and-Forget-Regressionstest, neuer
  [`SPEC-021`](../../../../spec/pflichtenheft.md) für die SSE-Drahtform.
- `harness/README.md` §Sensors/§Werkzeuge — die `make test-integration`-Zelle
  trägt beide Rundlauf-Abschnitte (`· seit slice-071`, `· seit slice-072`).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- **Ein `Broadcaster`, zwei Driving-Adapter.** Der zweite Zustellweg brauchte
  keinen neuen Server, keinen neuen Port und keine neue Secret-Klasse:
  dieselbe `withToken`-Middleware, dieselben Token-Klassen, derselbe
  `Broadcaster` als zweiter Abonnent — und der gRPC-Weg denselben
  `Bearer`-Wert über Metadata.
- **Die Port-Grenze trug ohne Nacharbeit.** `ChangeStreamPort` und der
  `Broadcaster` aus `slice-069` banden in `slice-070` und `slice-072` an,
  ohne ihren Vertrag zu ändern; getestet wird der `Broadcaster` in seinem
  eigenen Paket, die Adapter über je lokal deklarierte Fake-Subscriber
  (keine `adapters → adapters`-Kante).
- **Die E2E-Belege banden mehr als eine Sentinel-Zeile.** Die Fixrunde von
  `slice-071` bindet die **empfangene** `change_id` an `cdc.changes`; der
  Verifier hat eine gefälschte `change_id` rot gesehen. `slice-072`s Lauf
  prüft dieselbe Bindung über `sse_captured ≥ 1`.
- **Der Konflikt-Pfad trug zweimal real.** Der Selbstwiderspruch in
  [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  endete in [`ADR-0067`](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md),
  die Gate-Scope-Frage in
  [`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md) —
  beide als Folge-ADR mit partiellem `Supersedes`, keine stille
  Abweichung.
- **Fremde Mutationen je Zusage.** Der Verifier von `slice-072` sah drei
  eigene Mutationen rot (Oder-Bedingung → `&&`, `withToken` entfernt,
  `503`-Guard → `200`), der von `slice-071` eigene; beide Läufe sind
  unabhängig vom Implementer-Kontext.
- **Ein Compose-Lauf trägt beide Rundläufe.** Der Verifier von `slice-072`
  sah gRPC- **und** SSE-Rundlauf in **einem** grünen
  `make test-integration`-Lauf gegen **denselben** laufenden Feed-Container —
  das ist der gemeinsame Träger, den kein Einzel-DoD belegt.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **`ADR-0066` widersprach sich selbst.** Ihre dritte Fitness-Function-Zeile
  verlangte, dass `Capture()` auch bei einem nicht zurückkehrenden `Publish`
  nicht anhält, während dieselbe ADR den synchronen Aufruf festschreibt und
  die caller-seitige Zeit-Isolation ausdrücklich verwirft. Der Review fand
  das, der Architect-Zug bestätigte es und korrigierte die Zeile über
  [`ADR-0067`](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
  (partielles `Supersedes`); Registereintrag
  `BEO-PGC/fitness-function-gegen-eigene-entscheidung` (1×).
- **Der Gate-Scope brauchte eine Entscheidung.** Der neue Wegwerf-Client
  importiert den generierten Stub; der erste `make a-check`-Lauf war rot.
  Die erste Lösung (Ausnahme für `tools/**` im `composition_root`) war eine
  Gate-Lockerung ohne Träger (`AGENTS.md` §3.6) und wurde im Review als HIGH
  bestätigt; der Architect-Zug ersetzte sie über
  [`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  durch die Gruppe `tooling` + eine begrenzte Kante. Registereintrag
  `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (1×). `welle-19` §6 ist
  auf den ausgelieferten Stand nachgezogen — **kein** Eintrag in *Historische
  Trigger-Verschiebungen* (Begründung unter §Verifikation, *zwei
  Planner-Entscheidungen*).
- **`SPEC-018` → `SPEC-021`.** Die SSE-Drahtform lag zunächst als Erweiterung
  im [`SPEC-018`](../../../../spec/pflichtenheft.md)-Eintrag, dessen Intro sich
  auf die neun Port-gedeckten Fähigkeiten aus [`LH-FA-SST-006`](../../../../spec/lastenheft.md)
  abgrenzt — der neue Endpunkt ist keine davon und gehört zu einer anderen
  Anforderung. Bei der Closure von `slice-072` ist der Block in den neuen
  [`SPEC-021`](../../../../spec/pflichtenheft.md) herausgelöst (Review F-6,
  Verifier bestätigt).
- **Die Bootstrap-Entkopplung war nötig.** Ohne sie hätte der SSE-Weg
  `CDC_GRPC_ADDR` erzwungen; der `Broadcaster` entsteht seit `slice-072`,
  sobald **einer** der beiden Wege aktiv ist. Beide Adressen leer →
  unverändertes Bestandsverhalten (Mutation `||` → `&&` färbt genau die zwei
  Sub-Fälle rot).
- **Das *Mehr* der Welle war in §1 weiter gefasst als der operative Trigger
  in §3.** §1 verlangte „eine **einzelne** committed Änderung erreicht
  **gleichzeitig** über **beide** Wege einen verbundenen Client"; §3 verlangte
  zwei **getrennte** reale Rundläufe — genau die DoDs von `slice-071` und
  `slice-072`. Der Verifier hat die Spannung benannt (`verify-slice-072` §8).
  Entscheidung des Planners: §1 ist auf die **belegte** Aussage
  zurückgeschnitten — **gemeinsamer Träger**: eine Bootstrap-Bedingung, ein
  `Broadcaster`, beide Wege in einem Compose-Lauf gegen denselben laufenden
  Prozess —, und §3 benennt diesen Träger als eigene Bedingung; die
  **Gleichzeitigkeit** wird nicht mehr behauptet. Begründung: Der gemeinsame
  Träger ist die Aussage, die die Welle über die vier DoDs hinaus hinzufügt
  und die real belegt ist; die **Verteilungsbreite** eines `Publish`
  (n Abonnenten gleichzeitig) ist Eigenschaft des `Broadcaster` und in
  `make test` belegt (`internal/adapters/driven/grpcstream`), nicht Eigenschaft
  der Integration zweier Zustellwege. Ein Simultanbeleg wäre eine neue
  E2E-Form mit eigenem Aufwand und eigener Flake-Fläche für eine bereits
  unit-belegte Eigenschaft — also ein eigener Slice bei benanntem Bedarf, kein
  Bestandteil dieses Closure-Triggers. **Kein Folge-Slice**, weil kein Bedarf
  benannt ist (der Wellen-Trigger ist ohne ihn erfüllt).
- **Ein Registereintrag hat die Schwelle *mit* dieser Welle erreicht.** Der
  dritte Beleg (`evidence/slice-069.md`) hebt
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` auf 3×;
  der Ausgang (`verkörpert`, Zielorte `.claude/commands/implement-slice.md`
  Schritt 17 und `.harness/skills/reviewer.md`) ist bereits zugewiesen, die
  Anker-Paarung ist getragen. Die Welle trägt ihn deshalb nicht als eigenen
  Steering-Loop-Eintrag — siehe §Steering-Loop-Einträge und
  §Beobachtungs-Register.
- **Zwei §6-Risiken der Welle sind entfallen, weil die Parallelität nicht
  eintrat.** `slice-061` war bei Beginn von `slice-071`/`slice-072` längst
  geschlossen; die additiven Änderungen an `internal/bootstrap/wiring.go`,
  `harness/README.md` und `compose.yaml` trafen keine parallele Änderung.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

**Feststellung des Lese-Schritts: kein Registereintrag hat *mit dieser Welle*
die 3×-Schwelle neu erreicht und keinen Ausgang** — es gibt in dieser Closure
keinen Steering-Loop-Eintrag zu verkörpern und damit keinen
Planner → Architect → Planner-Zug (Modul 8 §Rollen-Sequenz für eine Welle,
Schritt 3b). Das ist das Ergebnis der Prüfung, keine Auslassung; die geprüften
Stände stehen unter §Beobachtungs-Register. Insbesondere:

- Der einzige Eintrag, der die Schwelle *mit* dieser Welle erreicht
  (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`, 3×),
  trägt seinen Ausgang `verkörpert` bereits (Zielorte siehe unten) — die Welle
  hat den Übertritt nicht erzeugt und verkörpert ihn nicht.
- Alle übrigen Einträge mit ≥ 3 Belegen tragen einen Ausgang (`verkörpert`,
  aufgelöst oder `geplant` mit Kennung); keiner steht ohne Ausgang über der
  Schwelle.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`.
In dieser Welle **neu angelegt**: `generierte-artefakte-ohne-sync-sensor`
(1×, `evidence/slice-069.md` — der erzeugte Protobuf-Code liegt committet im
Baum, kein Sensor hält ihn gegen seine `.proto`-Quelle),
`fitness-function-gegen-eigene-entscheidung` (1×, `evidence/slice-070.md`),
`gate-scope-erweiterung-ohne-adr-traeger` (1×, `evidence/slice-071.md`),
`aufschub-adresse-nimmt-sendung-nicht-an` (2×, `evidence/slice-071.md` +
`evidence/slice-072.md`).
**Fortgeschrieben**: `adapter-unittest-verdeckt-bootstrap-luecke`
(2×, `evidence/slice-072.md`) — der `503`-Pfad hat in der heutigen
Verdrahtung keinen Erzeuger; der Guard ist Grenz-Bedingung, kein toter Code.

**Die Schwelle mit dieser Welle erreicht, Ausgang bereits zugewiesen:**
`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3×,
`evidence/slice-069.md` ist der dritte Beleg) → Ausgang `verkörpert`, Zielorte
`.claude/commands/implement-slice.md` (Schritt 17, diff-skopierte
Selbstprüfung) und `.harness/skills/reviewer.md` (eigener benannter HIGH-Punkt),
Herkunfts-Anker `· seit slice-077`. **Anker-Paarung getragen**: beide Zielorte
existieren und tragen `seit slice-077` real.

**Über der Schwelle, Ausgang unverändert** (die Welle hat nur Belege
hinzugefügt, keinen neuen Übertritt erzeugt):
`slice-chronik-in-code-kommentar` (6×), `d-migrate-nacharbeit` (6×),
`report-nackte-id-ohne-link` (5×), `rollen-verdrahtung` (4×, aufgelöst seit
`slice-023`), `plan-vorlagen-defekt` (3×), `pipe-maskiert-make-exit-code`
(3×), `lese-doppelquelle` (3×), `handbuch-versionshistorie-uebersprungen` (3×),
`github-actions-unverifizierbar-lokal` (3×), `dod-checkbox-nachzug` (3×),
`dod-checkbox-nachzug-review-ohne-fixrunde` (3×), `a-check-null-abdeckung`
(3×, verkörpert seit `welle-1`), `commit-traceability-kein-vorab-hook` (3×,
Ausgang `geplant`, Träger `slice-073` in `open/`).

**Durchgesehen ohne Zähler-Änderung** (über die §8-Sichtungen und die §§7
dieser vier Slices): `adapter-fehler-ausgang` (2×, weiter offen),
`slice-pfad-als-link-in-berichten` (2×, weiter offen),
`slice-chronik-in-code-kommentar`, `rollen-verdrahtung`,
`a-check-null-abdeckung`, `commit-traceability-kein-vorab-hook`.
Unverändert, nicht von dieser Welle berührt: alle übrigen Einträge.

**Benannt, nicht gezählt — ein Punkt dieser Welle ohne Registereintrag:**
der unaufgelöste Vorlagenrest der §7-Anker-Zeile in `slice-070`, `slice-071`
und `slice-072` (drei geschlossene Vorgänge). Kein Gate fängt ihn
(`d-check`s `structure`-Regel verlangt „Steering-Loop", keine offenen Aufgaben
und nicht die Auflösung von Vorlagen-Zeilen); sichtbar wurde er erst bei der
Anker-Paarung dieser Closure, weil das Feld `liegt in` sie auslöst — mit
Platzhalter-Zielort hätte sie rot gemeldet. Die drei Reste sind entfernt
(`041dc0c`, Inhalts-Commit vor dem Wellen-`git mv`). Die Klasse ist mit
`plan-vorlagen-defekt` verwandt, aber von ihm **nicht** gedeckt: dessen
`observation.md` ist auf §2 und den Fill der Wellen-Eröffnung verengt (und
unveränderlich ab Anlage), der Zielort seiner Verkörperung
(`.claude/commands/plan-welle.md` §2-Form-Prüfung) reicht nicht in die
§7-Closure-Notiz. Ein eigener Eintrag ist **nicht** angelegt: Die
Register-Regel verbietet einen Eintrag, der eine Closure ohne Ausgang
übersteht, und die Zuordnung (Erweiterung des bestehenden Eintrags gegen einen
eigenen) samt Schärfung des Zielorts ist eine **Verkörperungs-Entscheidung**
(Planner → Architect → Planner, Modul 8 Schritt 3b). Sie ist hier benannt,
offen und mit dem Nachtrag in `plan-vorlagen-defekt/state.md` festgehalten.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

- `slice-077` (Handbuch-Betreiber-Stand) — wellenlos, aus `slice-072`s
  Review-Finding F-1; sein §2-DoD trägt die **drei** Netzwerk-Zugriffswege
  ausdrücklich mit, den SSE-Endpunkt eingeschlossen. Datei in `open/` —
  Folge-Slice-Paarung getragen.
- `slice-073` (Commit-msg-Git-Hook) — Träger des Ausgangs `geplant` von
  `commit-traceability-kein-vorab-hook`; nicht aus dieser Welle, aber in ihrem
  Register-Durchgang berührt. Datei in `open/`.
- `slice-076` (Coverage-Gate auf Reifestufe 40 %) — nicht aus dieser Welle;
  Träger der im Trigger-Audit benannten Reifestufe. Datei in `next/`.
- **Kein Folge-Slice aus dem *Mehr*:** ein Simultanbeleg ist nicht als Bedarf
  benannt — siehe §Was ging anders als geplant.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle vier Slices (`slice-069`…`slice-072`) liegen in `docs/plan/planning/done/`.
- `make gates` grün — Exit-Code je ungefiltert ermittelt und in einem **eigenen**
  Schritt geprüft, nie durch eine Pipe maskiert (`AGENTS.md` §3.9):
  - **Lauf 1** (Trigger prüfen, vor den Closure-Commits) — Exit **0**:
    `baseline-verify` grün (v6.5.0), `docs-check` 573 Dateien / 0 Befunde,
    `commit-traceability` OK (5 Commits, Betreffe ohne Struktur-ID),
    `a-check` `gesamt: 0 Befund(e)`, `coverage-gate` OK — Coverage **49,40 %**
    über der geltenden 35-%-Schwelle.
  - **Lauf 2** (nach allen Closure-Commits, Ergebnis-Notiz · Welle-§7-Zeiger ·
    Roadmap · `git mv` · Link-Reconciliation) — Exit **0**; Beleg siehe
    §Abschluss-Bemerkung unten.
- **`welle-19`s Closure-Trigger (§3)** — geprüft gegen die vier Bedingungen:
  - Slices in `done/` — erfüllt (vier Dateien, `git mv` je eigener Commit).
  - `make gates` grün — erfüllt (Lauf 1 und Lauf 2, Exit 0).
  - **Realer gRPC-E2E-Rundlauf** (`slice-071`) — erfüllt, belegt und
    zusammengefasst statt erneut erzeugt: `make test-integration` Exit 0 mit
    `READY` → `RECEIVED change_id=964-1 table=feed_e2e_full operation=INSERT
    new_image={"id":"261","name":"GrpcStreamE2ESentinel"}` → `REJECTED
    code=Unauthenticated`; die `change_id` ist gegen `cdc.changes` gebunden
    (`docs/reviews/verify-slice-071.md`).
  - **Realer SSE-E2E-Rundlauf** (`slice-072`) — erfüllt, ebenso zusammengefasst:
    derselbe `make test-integration`-Lauf des Verifiers zeigt `READY` →
    `RECEIVED change_id=967-1 table=feed_e2e_full operation=INSERT
    new_image={"id":"271","name":"SseStreamE2ESentinel"}` → `REJECTED code=401`;
    `sse_captured ≥ 1` hielt (`docs/reviews/verify-slice-072.md` §5).
  - **Gemeinsamer Träger** — erfüllt: **ein** Lauf führt beide Rundläufe gegen
    **denselben** Feed-Container; `changeStreamEnabled` (`wiring.go`) ist die
    eine Oder-Bedingung, derselbe `changeBroadcaster`-Zeiger geht an
    `apihttp.Server{Subscriber: …}` und an das gRPC-`Server{Subscriber: …}`
    (`docs/reviews/verify-slice-072.md` §4; Mutation `||` → `&&` rot).
  - Closure-Notiz — diese Datei.
- **Trigger-Audit der Welle** (Modul 6 §Wellen-Closure-Prozedur, Schritt 2 —
  drei Artefaktklassen, jede mit belegter Feststellung):
  - **Carveout:** `docs/plan/carveouts/` enthält ausschließlich `.gitkeep` —
    kein offener Carveout, kein rotes Gate, keine stehengebliebene Frist.
  - **bootstrap-aware Gate:** Die Welle führt keine neue Reifestufe ein und
    verbessert die Coverage-Zahl nicht (die neuen Belege liegen in
    `test/integration/`, `tools/harness/` und den Adapter-Paketen; `THRESHOLD`
    steht unverändert bei 35 %). **Feststellung, nicht Auslassung:** der
    deklarierte Hochschalt-Trigger des `coverage-gate` ist **nach wie vor
    fällig** — der real gemessene Ist-Stand (49,40 %) liegt über der nächsten
    5-%-Stufe (40 %) der geltenden 35-%-Schwelle. Die Entscheidung ist
    getroffen (`architect-verdict-coverage-gate-reifestufe`, Stufe 40 %), der
    **Träger ist `slice-076`** und liegt priorisiert in `next/`, noch nicht in
    `in-progress/`. Die Stufe ist damit nicht stehengeblieben, sondern
    adressiert (Trigger-Audit der `welle-18`-Closure → Architect-Zug →
    `slice-076`). **Benannt für den nächsten Trigger-Audit:**
    `slice-076` §1 bindet die Stufe **45 %** an ihren eigenen Trigger
    („Ist-Stand ≥ 45 % bei einem weiteren Coverage-Fortschritt") — der
    gemessene Stand erfüllt die Wert-Hälfte bereits (49,40 %). Mit dem
    Hochschalten auf 40 % wird 45 % „die nächste Stufe"; die 45er-Stufe ist
    dann fällig und braucht ihren eigenen Grün-/Rot-Beleg (§3.6 bleibt
    unberührt, es ist keine Senkung).
  - **ADR mit Re-Evaluierungs-Trigger:** [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
    (stream-internes Replay · tabellen-granulare Filterung) — **keiner
    eingetreten**: kein Bedarf benannt, kein Registereintrag auf der Schwelle,
    beide Punkte stehen unverändert in §6 dieser Welle.
    [`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
    (`Last-Event-ID`-Replay) — **nicht eingetreten**: `sse.go` liest den Header
    nicht und sendet keine `id:`-Zeile (Verifier-Negativbefund).
    [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
    (Replay · Zustellgarantie jenseits Fire-and-Forget) — **nicht
    eingetreten**: kein Bedarf an Rückstau-/Metrik-Zusage.
    [`ADR-0067`](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
    (caller-seitige Zeit-Isolation) — **nicht eingetreten**:
    `Publish` bleibt strukturell nicht-blockierend, der Port trägt die
    Nicht-Blockade weiter.
    [`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
    (der Wegwerf-Client braucht mehr als die eine begrenzte Kante · oder der
    generierte Stub wird aus `adapters` herausgelöst) — **nicht eingetreten**:
    `grpcclient` importiert genau ein Paket, `sseclient` keines; kein zweiter
    Bedarf benannt. [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
    — Trigger unverändert (kein Linter, Ist-Stand nicht bei 80 %); der
    Reifestufen-Träger steht daneben und ist kein ADR-Trigger (siehe oben).
    **Feststellung: kein ADR-Trigger dieser Welle ist fällig, keine Nacharbeit
    nötig.**
- **Zwei Planner-Entscheidungen mit Begründung** (beide ohne Folge-Slice):
  1. **`welle-19` §6, `.a-check.yml`** — nachgezogen auf den ausgelieferten
     Stand (Gruppe `tooling` + begrenzte Kante, Träger
     [`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md))
     und **kein** Eintrag in *Historische Trigger-Verschiebungen*: Das
     Drift-Log führt Umplanungen der Wellen-Sequenz (Trigger verschoben,
     präzisiert oder ersetzt; Slice oder Welle umgehängt). Hier ist nichts
     sequenziert worden — dieselben vier Slices in derselben Reihenfolge mit
     denselben Triggern; geändert hat sich eine **Out-of-Scope-Abgrenzung
     innerhalb** des Wellen-Plans, und zwar nicht durch eine Planungs-,
     sondern durch eine **Entscheidung** (Review-HIGH → Architect-Zug →
     `ADR-0068`), die dort vollständig dokumentiert ist. Ein Drift-Eintrag
     wäre ein zweites Log zur selben Änderung.
  2. **Das *Mehr* (§1 gegen §3)** — §1 auf die belegte Aussage zurückgeschnitten
     (gemeinsamer Träger), §3 um die Träger-Bedingung ergänzt, die
     Gleichzeitigkeit ausdrücklich **nicht** behauptet; **kein** Folge-Slice
     für einen Simultanbeleg. Begründung siehe §Was ging anders als geplant.
- **Drei Paarungen** (Anker · Folge-Slice · Register) — geprüft **nach** dem
  `git mv` dieser Welle-Plan-Datei nach `done/`:
  - **Anker** — diese Closure führt **keinen** Steering-Loop-Eintrag mit
    `liegt in`-Feld; die drei unaufgelösten Vorlagen-Zeilen, die das Feld
    nur zitierten, sind vorher entfernt. Zusätzlich geprüft (weil der
    Ausgang eines Fremd-Eintrags sie berührt): die zwei Zielorte von
    `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
    (`.claude/commands/implement-slice.md` Schritt 17,
    `.harness/skills/reviewer.md`) existieren und tragen `seit slice-077`.
  - **Folge-Slice** — `slice-073`, `slice-076` und `slice-077` existieren real
    als Dateien im Planning-Lifecycle (`open/` bzw. `next/`) — grün.
  - **Register** — alle in dieser Welle genannten Kennungen
    (`generierte-artefakte-ohne-sync-sensor`,
    `fitness-function-gegen-eigene-entscheidung`,
    `gate-scope-erweiterung-ohne-adr-traeger`,
    `aufschub-adresse-nimmt-sendung-nicht-an`,
    `adapter-unittest-verdeckt-bootstrap-luecke`,
    `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
    `commit-traceability-kein-vorab-hook`, `plan-vorlagen-defekt`) existieren
    als Verzeichnisse mit nicht leerem `evidence/` — grün.
  - **Gegenprobe zur Nachbarschaftsklasse:** kein Bericht und kein Slice-Plan
    dieser Welle adressiert `welle-19` als Markdown-Link mit festem
    Verzeichnis; die `**Welle:**`-Kopfzeilen und die Berichte zitieren die
    Kennung in Code-Form. `BEO-PGC/slice-pfad-als-link-in-berichten` bleibt
    damit bei 2×.

## Archivierung

Feststellung: Das Repo führt weiterhin **kein Archivierungs-Werkzeug** (kein
`archiv.zip`-Target, keine Stub-Erzeugung) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**. Die vier Slice-Dateien dieser Welle, ihre
Review-/Verifier-Reports und der Wellen-Plan bleiben vollständig in `done/`.
Dieselbe Feststellung trafen bereits `welle-14`, `welle-15`, `welle-16`,
`welle-17` und `welle-18`; die Vorprüfung des Sensor-Geltungsbereichs ist damit
weiterhin nicht ausgelöst (kein Sensor keilt auf `done/*.md`, der durch Stubs
eine Ebene tiefer erblinden würde).
