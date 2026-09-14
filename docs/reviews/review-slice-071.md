# Review-Report: slice-071 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-071`, Diff `3425bfb..6862667` (= `HEAD`) — zwei Commits:
`b835dde` (Client + E2E-Rundlauf + Gate-Scope, fünf Dateien, +261/−2) und
`6862667` (Plan-/DoD-Nachzug im selben Slice-Plan, +43/−4). Der Elter-Commit
`3425bfb` ist der reine `next → in-progress`-Move; `git diff --name-only
3425bfb..HEAD` nennt genau die sechs erwarteten Pfade — `.a-check.yml`,
`compose.yaml`, `harness/README.md`, `tools/harness/grpcclient/main.go`,
`tools/harness/run-integration-tests.sh`,
`docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md`. Kein
Fremd-Commit im Range.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14, 21:14 — die Schärfung `Neue Betreiber-Oberfläche ohne
Handbuch-Zug` liegt **vor** dem Implementer-Lauf 22:18).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-14.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md`
  vollständig (§1–§8, inkl. §3-Nachzug und DoD-Stand aus `6862667`), plus
  seine Vorfassung in `2621d89` (Wellen-Eröffnung `welle-19`, 07:15:44)
- [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) vollständig —
  insbesondere Teilfrage 2 (Port-/Adapter-Design, kein Adapter-Import),
  Teilfrage 3 (Fire-and-Forget ohne Replay), Teilfrage 4 (Authn), Teilfrage 5
  (Adapter-Platzierung), Teilfrage 6 (Bootstrap), §Fitness Function,
  §Slice-Schnitt-Empfehlung (Slice C = dieser Slice)
- [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  vollständig — §Entscheidung Festlegungen 1–4, §Verhältnis zu `ADR-0061`,
  §Konsequenzen (stiller Verlust ist zugesagte Semantik)
- [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
  vollständig — §Entscheidung, §Konsequenzen, §Fitness Function,
  §Re-Evaluierungs-Trigger
- [`ADR-0026`](../plan/adr/0026-composition-root.md),
  [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) §Teilfrage 4/§Testabdeckung/
  §Fitness Function, [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Fitness Function, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md),
  [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
- `spec/architecture.md` §1/§2 (Zeilen `ARC-005`/`ARC-007`),
  [`SPEC-020`](../../spec/pflichtenheft.md) (Nachrichtenschema, Zustellgarantie,
  Erzeuger-Blockade, Authentifizierung, Aktivierung),
  [`SPEC-019`](../../spec/pflichtenheft.md) (Abgrenzung — Antrags-Datensatz, **nicht**
  der Stream), [`LH-FA-SST-008`](../../spec/lastenheft.md) (Happy Path/
  Boundary/Negative), [`LH-FA-SST-006`](../../spec/lastenheft.md),
  [`LH-FA-SST-007`](../../spec/lastenheft.md)
- `AGENTS.md` §3.1/§3.5/§3.6/§3.7/§3.9/§5/§6; `harness/conventions.md`
  (`MR-000`); `harness/sensors/a-check.md`; `a-check.mk`; `.dockerignore`
- `docs/plan/planning/done/slice-069-grpc-streaming-adapter-grundgeruest.md`
  (§3 Handbuch-Abgrenzung), `…/done/slice-070-grpc-capture-integration.md`,
  [`review-slice-069`](review-slice-069.md), [`review-slice-070`](review-slice-070.md)
- `.a-check.yml` samt seiner Vorgeschichte (`174ae63`, `8d23352`, `b835dde`),
  `internal/adapters/driving/grpc/server.go` und `…/server_test.go`,
  `…/interceptor.go`, `tools/harness/httpclient/main.go`,
  `tools/harness/natssub/main.go`
- Beobachtungs-Register `BEO-PGC`: `a-check-null-abdeckung` (verkörpert),
  `fitness-function-gegen-eigene-entscheidung` (1×),
  `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (3×,
  verkörpert, Träger `slice-077`), `plan-nachzug` (verkörpert),
  `slice-chronik-in-code-kommentar` (verkörpert),
  `commit-traceability-kein-vorab-hook` (3×)

---

## Findings

### F-1 — `composition_root: [… , "tools/**"]` schaltet einen roten Architektur-Befund für den ganzen Werkzeug-Baum grün, ohne die Rolle zu tragen, die der Schlüssel deklariert

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.6 (Hard Rule: „Jede Schwellen-Senkung … Architekturregel
  ist ein ADR, kein PR-Kommentar") · [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
  §Entscheidung (`:57-59`) und §Re-Evaluierungs-Trigger (`:129`) ·
  [`ADR-0026`](../plan/adr/0026-composition-root.md) §Entscheidung ·
  `spec/architecture.md:79` (Zeile `ARC-007`) · Plan §3 (`:107`, `:111-122`)
- `pfad`: `.a-check.yml:28-34` gegen
  `docs/plan/adr/0041-a-check-maschinenform-architekturpruefung.md:57-59`,
  `:114`, `:129-130`; `spec/architecture.md:60`, `:79`;
  `docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md:107`,
  `:111-122`
- `befund`: Die Deklaration nimmt den **gesamten** `tools/**`-Baum in den
  Schlüssel auf, dessen deklarierte Bedeutung „Verdrahtung konkreter Adapter"
  ist (`.a-check.yml:20-21`, Bezug `ADR-0026`). Der Anlass-Pfad erfüllt sie
  nicht: `tools/harness/grpcclient/main.go:25` importiert genau **ein**
  Paket (`internal/adapters/driving/grpc/streamv1`, der erzeugte Stub) und
  verdrahtet nichts; `tools/harness/httpclient`/`natssub` importieren kein
  Paket unter `internal/`. Der Nachzug begründet die Aufnahme mit „derselbe
  testseitige Verdrahtungs-Konsument wie `test/integration`"
  (`:107`, `:116-118`) — falsch: `test/integration/*.go` importiert sieben
  interne Pakete (`postgresstorage`, `inbound`, `outbound`, vier Use-Cases,
  `domain/model`), also echte Verdrahtung, der Client keines davon. Der
  Kommentar im Diff beruft sich auf „`ADR-0041` §Entscheidung:
  `composition_root` ist deklarativer Stand der `.a-check.yml`, kein
  Schichten-Edge" — die Prämisse stimmt wörtlich, die daraus gezogene
  Schlussfolgerung („deshalb dieselbe Ausnahme") hat in `ADR-0041` keine
  Deckung: die Klausel, die Änderungen von der ADR-Pflicht ausnimmt, nennt
  „Änderungen an Schichten-Globs und Edges" (`:58`) und der
  Re-Evaluierungs-Trigger wiederholt genau diese zwei (`:129`) —
  `composition_root` steht dort in der Definition des Standes, nicht in der
  Ausnahme.
- `verifizierbar`: **ja** — zwei Proben, in diesem Review real ausgeführt
  (§Mutations- und Probenläufe): `make a-check` mit der **Vor**-Fassung der
  Datei auf demselben Baum endet Exit 2 mit genau dem im Nachzug zitierten
  Befund `wrong-direction: (ohne Schicht) -> adapters (…/streamv1)`; eine
  Probe-Datei unter `tools/harness/grpcclient/`, die
  `internal/application/usecase/capture` importiert, endet Exit 0 — die
  Aufnahme gewährt dem ganzen Baum das Importrecht der Composition Root.
- `klasse`: „Gate-Scope-Erweiterung ohne die Rolle, die der Schlüssel deklariert"

### F-2 — Erfolgsmeldung und Sensor-Zelle behaupten „vollständiger Inhalt" und eine Bestätigung *der empfangenen* Änderung; die Assertion prüft weder das eine noch das andere

- `kategorie`: LOW
- `quelle`: Maintainability · [`LH-FA-SST-008`](../../spec/lastenheft.md)
  Happy Path („den vollständigen Change-Inhalt")
- `pfad`: `tools/harness/run-integration-tests.sh:1861` (Regex) und `:1879`
  (SQL) gegen `:1891` (Erfolgszeile); `harness/README.md:132` (`· seit
  slice-071`-Satz)
- `befund`: Die Assertion pinnt `table`, `operation` und **einen** Wert der
  Spalte `name` als Teilzeichenkette in `new_image`; `feed_e2e_full` trägt
  zwei Spalten (`run-integration-tests.sh:189`), und die SQL-Quelle filtert
  auf `table_name` + `new_data->>'name'` — die im Client gedruckte
  `change_id` wird nirgends gegen `cdc.changes` gehalten. Die Erfolgszeile
  sagt „empfing … eine Änderung mit vollständigem Inhalt (real in
  `cdc.changes` lesbar)", die README-Zelle „die empfangene Änderung wird
  zusätzlich real über `cdc.changes` bestätigt" — beides ist breiter als die
  Assertion. Die tragende Vollständigkeits-Aussage liegt eine Ebene tiefer
  und ist dort rot-fähig (`internal/adapters/driving/grpc/server_test.go:198-207`
  prüft alle zehn Nachrichtenfelder und beide Row Images); der E2E-Tier
  belegt den realen Transport, nicht die Feldvollständigkeit.
- `verifizierbar`: ja — Assertion und SQL sind zeichengleich lesbar; eine
  Gegenprobe wäre eine Mutation des Mappers (zweiter Spaltenschlüssel
  entfällt), die den E2E-Lauf grün ließe.
- `klasse`: „Erfolgsmeldung behauptet mehr als die Assertion prüft"

### F-3 — Der Kopf des Slice-Plans zitiert `SPEC-019` als die vom Slice berührte Spec-Stelle; der gRPC-Stream steht in `SPEC-020`

- `kategorie`: LOW
- `quelle`: Maintainability (`AGENTS.md` §5 — Kennungen adressieren die
  Zusage, die sie nennen) · `harness/conventions.md` (`MR-000`)
- `pfad`:
  `docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md:13-15`
- `befund`: Der Abschnitt *Berührte Spec-Stellen* führt „das in `slice-069`
  festgelegte `SPEC-019`" als Bezugs-Stelle. `SPEC-019` ist der
  Antrags-Datensatz `cdc.administration_request` (eingeführt von `131fd98`),
  der gRPC-Stream ist `SPEC-020` (eingeführt von `7c02d15`, dem
  Grundgerüst-Slice). Der Fehler stammt aus der Wellen-Eröffnung
  (`2621d89`, 07:15:44) und steht außerhalb der Nachzug-Hunks — er steht
  trotzdem im Review-Gegenstand, weil der Plan die halbe Prüfgrundlage des
  Slice ist und der Slice noch in `in-progress/` liegt.
- `verifizierbar`: ja — `git log -S "### SPEC-020" -- spec/pflichtenheft.md`
  → `7c02d15`; `git log -S "### SPEC-019" …` → `131fd98`.
- `klasse`: „Spec-Stelle mit fremder Kennung adressiert"

### F-4 — Der Aufschub für die Betreiber-Dokumentation nennt `slice-072` als Adresse; die Sendung nimmt `slice-077` an

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Ziel-Form: Slice (Out-of-Scope-Klasse 1: „Die Adresse muss die Sendung
  annehmen") · `.harness/skills/reviewer.md` §Klassifikation (HIGH „Neue
  Betreiber-Oberfläche ohne Handbuch-Zug" — hier **nicht** erfüllt, s. u.)
- `pfad`:
  `docs/plan/planning/in-progress/slice-071-grpc-beispielclient-e2e.md:138-145`;
  `docs/plan/planning/open/slice-077-handbuch-betreiber-stand.md:52-58`;
  `docs/plan/planning/open/slice-072-http-sse-streaming-endpunkt.md` (§1,
  kein Handbuch-Bezug)
- `befund`: Der in `6862667` **neu hinzugefügte** Absatz *Handbuch-Assessment*
  verweist die Operator-Dokumentation des Streamings auf die
  „Operator-Erreichbarkeit (`slice-072`)". `slice-077` führt in seinem Ziel
  „die drei fehlenden Umgebungsvariablen-Gruppen in §5 (HTTP/JSON-API,
  gRPC-Stream)" und grenzt die Beispiel-Clients ausdrücklich auf
  `slice-071`/`slice-072` ab; `slice-072`s Plan nennt `benutzerhandbuch`
  nirgends. `slice-077` lag seit 21:01:42 in `open/`, der zitierte Absatz
  entstand 22:18:24 — die Adresse war zum Schreibzeitpunkt bereits die
  falsche, und der Implementer war durch §8 des eigenen Plans auf die
  Register-Sichtung verpflichtet. Die Klasse „Neue Betreiber-Oberfläche ohne
  Handbuch-Zug" ist dennoch **nicht** ausgelöst: die Variable
  `CDC_GRPC_ADDR` wurde von `slice-069` eingeführt (dort bereits als dritter
  Beleg gezählt, `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/evidence/slice-069.md`),
  und dieser Diff setzt sie nur im E2E-Compose-Vertrag — ein benannter
  Aufschub existiert, er zeigt nur auf den falschen Empfänger.
- `verifizierbar`: ja — `ls docs/plan/planning/open/`, `slice-077` §1 gegen
  `slice-072` §1; Commit-Zeitstempel `8423b28` (21:01) gegen `6862667` (22:18).
- `klasse`: „Aufschub-Adresse nimmt die Sendung nicht an"

### F-5 — Die Negative-Hälfte des E2E-Belegs ist ohne rot gesehenes Gegenbeispiel zugesagt

- `kategorie`: INFO
- `quelle`: [`LH-FA-SST-008`](../../spec/lastenheft.md) Negative ·
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) Teilfrage 4
- `pfad`: `tools/harness/grpcclient/main.go:75-107`;
  `tools/harness/run-integration-tests.sh:1851-1855`, `:1865-1868`;
  Plan §3 (`:132-137`)
- `befund`: Der Client verlangt `Unauthenticated` über zwei Ausgänge
  (Aufruf-Fehler, erster `Recv`-Fehler) und meldet `REJECTED` nur, wenn
  `status.Code(err)` genau `Unauthenticated` ist; jeder andere Ausgang —
  akzeptierter Stream, falscher Status, Fristablauf — endet mit Ausgang 1,
  der Runner verlangt zusätzlich `REJECTED code=Unauthenticated` **und**
  Client-Ausgang 0. Der Wächter kann damit nicht leerlaufen. Was fehlt, ist
  allein die Gegenprobe (der Implementer benennt sie offen): ein Server, der
  tokenlose Streams annimmt, existiert in dieser Compose-Umgebung nicht, und
  das Verhalten selbst ist auf Unit-Ebene rot-fähig getragen —
  `server_test.go:125-156` prüft drei Metadata-Grenzen gegen den echten
  Interceptor, `…:112-118` macht einen hängenden Stream zum Fehlschlag
  (ein entfernter Interceptor färbt sie also rot).
- `verifizierbar`: nein — die Aussage ist die Abwesenheit eines Belegs; die
  Sensibilität der Unit-Fassung ist dagegen lesbar belegt.
- `klasse`: „Negativ-Hälfte ohne rot gesehenes Gegenbeispiel"

### F-6 — `READY` ist ein lokaler Marker, nicht die Bestätigung, dass der Server den Stream angenommen hat

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.7 — ein Kommentar beschreibt,
  was da ist) · [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  §Entscheidung Festlegung 4
- `pfad`: `tools/harness/grpcclient/main.go:57-64` gegen
  `tools/harness/run-integration-tests.sh:1775-1781`;
  Plan §3 (`:123-131`)
- `befund`: Der generierte Client-Aufruf kehrt zurück, nachdem der Stream
  lokal angelegt und die Anfrage gesendet ist — `grpc.NewClient` verbindet
  erst beim ersten RPC; „Die Bereitschafts-Zeile folgt dem Aufbau der
  Streaming-Verbindung" (`main.go:62-63`) ist damit früher datiert, als das
  Wort „Aufbau" nahelegt. Die tragende Annahme des Rundlaufs
  („Registrierung des Empfängers am Broadcaster läuft serverseitig asynchron
  dazu", `run-integration-tests.sh:1793-1798`) ist dadurch **nicht**
  falsch, sondern eher konservativ — der Puffer-Loop ist die richtige
  Antwort auf ein Fenster, das nicht schmaler ist als beschrieben.
- `verifizierbar`: nein — Formulierungs-Aussage; der Mechanismus ist an
  `grpc-Go`-Aufrufform und `server.go:104-107` (Subscribe im
  Handler-Rumpf) nachvollzogen.
- `klasse`: „Bereitschafts-Zeile als serverseitiges Signal gelesen"

---

## Zum `.a-check.yml`-Verdikt — Kernfrage dieses Auftrags

**Verdikt: die Aufnahme von `tools/**` ist keine zulässige Ergänzung des
deklarativen Standes allein — sie braucht eine Architect-Entscheidung (ADR
als Träger oder ein Verdikt, das den Gegenstand der Ausnahme benennt).
`AGENTS.md` §3.6 greift.** Drei Prüfungen, alle am Text geführt:

**(a) Was `ADR-0041` wörtlich über `composition_root` sagt.** Der Satz lautet
vollständig (`:57-59`): „Die `.a-check.yml` ist ihr deklarativer Stand
(Schichten und Edges, `composition_root`); Änderungen an Schichten-Globs und
Edges sind Änderungen dieser Datei, keine neuen ADRs, solange die
§2-Constraints unverändert abgebildet bleiben." Der Trigger-Abschnitt
wiederholt die Ausnahme ein zweites Mal (`:129`): „Änderungen an
Schichten-Globs und Edgen sind `.a-check.yml`-Änderungen, keine ADRs."
`composition_root` wird also **einmal als Teil des Standes genannt und
zweimal aus der ADR-freien Änderungsklasse ausgelassen**. Zwei Auslassungen
an zwei Stellen sind keine Nachlässigkeit, die man weglesen darf; die
Implementer-Formulierung zitiert die tragende Hälfte und lässt die
einschränkende weg — genau der Fehlertyp, den
`BEO-PGC/fitness-function-gegen-eigene-entscheidung` beschreibt.

**Sind §2-Constraints berührt?** Nicht über die sieben Zeilen der Tabelle —
`tools/**` ist keine Komponente. Wohl aber über die Zeile, deren Maschinenform
`composition_root` **ist**: `spec/architecture.md:79` und `:60` führen
`ARC-007` als „Bootstrap … kennt als **einzige** Komponente konkrete Adapter".
`composition_root` ist nicht „der Rest", sondern die Pfad-Menge, die diese
Zeile abbildet; `ADR-0041`s eigene Fitness-Function-Zeile nennt denselben
Bezug (`:114`: „`composition_root` außerhalb der Schichten"). Für
Nicht-Komponenten-Pfade ist die Aufnahme also keine Abbildung von §2,
sondern deren Erweiterung — und `test/integration/**` rechnet dieselbe
Erweiterung bereits (dazu (c)).

**(b) Ist `tools/**` als Ganzes eine sinnvolle Einheit? Nein — enger wäre
`tools/harness/**`, und die Einheit ist nicht einmal das eigentliche
Problem.** Der hinzugefügte Kommentar begründet die Aufnahme mit „die
testseitigen Wegwerf-Clients des Harness" (`.a-check.yml:28`), der Glob
umfasst aber auch `tools/schema/**` (die d-migrate-Werkzeugkette) und die
Bench-Skripte — beides **keine** Wegwerf-Clients, beides laut Kommentar nicht
gerechtfertigt. Fachlich schwerer wiegt die Klasse: `test/integration/**`
verdrahtet real, der Client **konsumiert** über das Netz. Die richtige
Ausnahme wäre eine, die die Rolle „Client des erzeugten Protokoll-Stubs"
trägt; `composition_root` ist die unbeschränkteste Rolle dieses Modells.
Die Probe dieses Reviews zeigt das Maß: eine Datei unter
`tools/harness/grpcclient/`, die einen Produktions-Use-Case importiert,
bleibt grün (Exit 0, ohne jeden Hinweis). Nebenbefund desselben Laufs: der
Hinweis „3 gescannte Dateien liegen in keiner Schicht und bleiben
ungeprüft" (`harness/sensors/a-check.md` Grenze 1) verstummt mit der
Aufnahme ersatzlos — der Diff tauscht ein sichtbares „ungeprüft" gegen ein
stilles „deklariert".

**(c) Trägt der Präzedenzfall `test/integration/**` den Verzicht auf ein
ADR? Nur halb — als Mechanik, nicht als Klasse.** Der Präzedenzfall ist real
und wurde in einem Slice-Commit ohne ADR vollzogen (`8d23352`,
`test(integration): …`, Begründung „die Verdrahtungs-Tests tragen denselben
Aufbau wie der Bootstrap", Bezug `ADR-0026`). Er stützt genau eine Aussage:
**eine** `composition_root`-Erweiterung ist in diesem Repo als Datei-Änderung
praktiziert worden. Er trägt aber nicht die Gleichsetzung, auf die sich
dieser Nachzug stützt: `test/integration/**` ist ein Verdrahtungs-Konsument
im Wortsinn, der Client ist es nicht (a). Und `ADR-0041` ist **nach** dem
Präzedenzfall entstanden (09-09 vs. 09-10) und hat `composition_root` in der
Ausnahmeklausel gerade nicht mitgeführt — der Präzedenzfall ist keine
Entscheidung über die Klausel, er beschreibt sie stillschweigend fort.
Die Reviews der Slices, die ihn erwähnen
([`review-slice-007`](review-slice-007.md), [`review-slice-008`](review-slice-008.md),
`architect-review-welle-1.md`), behandeln `composition_root` durchweg als
**Faktum** („`test/integration` ist `composition_root`"), nie als
entschiedene Ausnahme.

**Was der Architect-Zug ohne Rückfrage vorfindet.** Die Frage, entscheidbar
in einem Wort, mit den drei zulässigen Antworten des Repos:

1. *Zulässige Ergänzung des Standes* → dann ist `ADR-0041`s Ausnahmeklausel
   (`:58`, `:129`) nachzuziehen, damit sie ihren Gegenstand nennt — sie
   steht heute enger als die Praxis; `composition_root` gehört in beide
   Sätze. Träger: Folge-ADR mit `Supersedes ADR-0041` (nur die Klausel) oder
   ein Verdikt, das in `.a-check.yml`/`harness/sensors/a-check.md` zitiert
   wird.
2. *Regel gilt, Erweiterung ist eine Lockerung* → ADR nach §3.6 **oder**
   Rücknahme der Zeile; fällt die Zeile, braucht der Client einen anderen
   Träger für seinen Stub-Import (Gegenstand eines Folge-Schnitts, nicht
   dieses Reports).
3. *Erweiterung zulässig, aber falsch zugeschnitten* → ADR, die den
   Gegenstand der Ausnahme benennt (Rolle „testseitiger
   Adapter-Konsument" / engerer Glob), plus Plan-Korrektur.

Kein Verdikt ist „Zeile bleibt, Kommentar bleibt": der Kommentar zitiert
`ADR-0041` für eine Aussage, die dort nicht steht.

**Warum HIGH und nicht MEDIUM:** die Änderung ist kein Gestaltungsdetail,
sondern der Vollzug einer Gate-Lockerung durch Konfiguration
(`make a-check` rot → grün für einen ganzen Baum). `AGENTS.md` §3.6 macht
genau das zur ADR-Pflicht, und der Reviewer-Skill führt „ADR-Verstoß (Layer,
Tool, Hard Rule)" als HIGH. `ADR-0060`s eigene Fitness-Function-Zeile („`.a-check`
unverändert") ist demgegenüber **nicht** verletzt — sie ist auf die zwei
neuen Adapter-Pakete des `ADR`-Gegenstands bezogen, und deren Globs sind
tatsächlich unberührt; siehe Negativbefunde.

---

## Negativbefunde

- **geprüft, ohne Befund: Ordnung des Rundlaufs (Auftrags-Punkt 1).** Der
  additive Block steht `run-integration-tests.sh:1757-1891`, der
  Upgrade-Sicherheits-Container-Tausch beginnt `:1893`, die
  Container-Ende-Grenze von `TestE2ESchemaChangeDropColumn` liegt `:2004` und
  die von `TestE2ESchemaChangeIncompatibleTypeChange` `:2099`. Der Rundlauf
  liegt damit **vollständig** vor beiden Schema-Negative-Aufrufen und vor dem
  Containertausch — die Bedingung, die `harness/README.md:132` für den
  HTTP-Beleg formuliert, ist für den gRPC-Beleg identisch erfüllt und im
  Blockkopf (`:1768-1772`) benannt.
- **geprüft, ohne Befund: der reale Beleg ist nicht nur eine Zeile
  (Auftrags-Punkt 2).** Die Assertion verlangt Tabelle, Operation **und** den
  Wert der Spalte `name` **im** `new_image` (`:1861`), der Client druckt
  `new_image=%s` aus `change.NewImage` (`main.go:70-71`), und `feed_e2e_full`
  ist über den ganzen Lauf aktiviert (`ADMIN_TABLE`/`COLUMN_TABLE` sind
  andere Tabellen, `:289`, `:1100`). Die ID-Prüfung 260–265 kollidiert mit
  keiner anderen Insert-Stelle (`:1474`, `:1527`, `:1631`, `:1682` nutzen
  230/231/240/241; der nachfolgende Upgrade-Block 250/251). Abweichung nur
  im Umfang der Zusage — F-2.
- **geprüft, ohne Befund: `ADR-0066`-Semantik und der bounded re-emit-Loop
  (Auftrags-Punkt 4).** Der Verlust der ersten committeten Änderung ist die
  zugesagte Semantik: `Subscribe()` wird im Handler-Rumpf registriert
  (`server.go:104-107`), die Bereitschafts-Zeile des Clients entsteht
  vor dem Server-Kontakt (F-6), und `ADR-0066` Festlegung 4 erklärt „Ein
  nicht verbundener oder langsamer lesender Empfänger verliert die
  betroffenen Nachrichten ersatzlos" zur Zusage, `SPEC-020` wiederholt sie.
  Der Loop verdeckt deshalb **keinen** Defekt: er wiederholt keinen
  zugesagten Nachliefer-Mechanismus (den es nicht gibt), sondern committet
  erneut und belegt damit, dass der Stream nach der Registrierung real
  trägt. Ein Defekt wäre nur, wenn die Registrierung selbst verspätet oder
  unzuverlässig wäre — dafür gibt es keinen Anhalt; die Loop-Grenze (`5`
  Versuche) liegt unter der Client-Frist (`60 s`).
- **geprüft, ohne Befund: `compose.yaml` (Auftrags-Punkt 5).** `CDC_GRPC_ADDR:
  ":9090"` steht in derselben Doppelpunkt-Form wie `CDC_HTTP_ADDR: ":8090"`
  (`:111` gegen `:106`), der Kommentar benennt Form und Grund (Netz-Alias,
  Bindung auf allen Interfaces), die Token-Klassen werden unverändert
  wiederverwendet (`CDC_API_TOKEN_READER` `:107`; der Runner übergibt
  `$HTTP_TOKEN_READER` = `e2e-reader-token`, `:68`, `:1783`). Kein
  Widerspruch zu `SPEC-020` („Aktivierung optional über `CDC_GRPC_ADDR`").
- **geprüft, ohne Befund: `harness/README.md` §Sensors (Auftrags-Punkt 6).**
  Der Satz steht in der `make test-integration`-Zelle (`:132`), trägt
  `· seit slice-071` in der Form aller Nachbarsätze und ändert die
  Bindungsspalte nicht („kein Gate, `ADR-0030`, `LH-QA-POR-003`") — keine
  neue Gate-Behauptung, kein neues Target.
- **geprüft, ohne Befund: Kommentar-Disziplin (Auftrags-Punkt 7).**
  `grep -nE "slice-[0-9]+|welle-[0-9]+"` über `.a-check.yml`, `compose.yaml`,
  `tools/harness/grpcclient/main.go` und die **hinzugefügten** Zeilen von
  `run-integration-tests.sh`: kein Treffer (die vier Treffer der Dateien
  liegen sämtlich in Bestands-Kommentaren — `compose.yaml:91`, `:132`;
  Runner `:20`, `:111`, `:225`, `:990`, `:1049`, `:1899`). Keine
  `//nolint`/`#noqa`-Suppression (`AGENTS.md` §3.2), kein Build-/Test-Pfad
  außerhalb von `make` (`§3.1`) — der Client läuft per `go run` im gepinnten
  Toolchain-Container, die Modul-Cache-Volume ist dieselbe wie beim
  bestehenden `httpclient` (`:1778-1783` gegen `:1722-1730`).
- **geprüft, ohne Befund: Scope-Fidelity (Auftrags-Punkt 8).** Kein Pfad
  unter `internal/**` im Range; `test/integration/**` unberührt; der neue
  Client importiert genau ein internes Paket und liegt außerhalb jeder
  Schicht. `tools/harness/grpcclient/main.go` ist über `.dockerignore` vom
  Build-Kontext ausgeschlossen (`*`, nur `cmd/`, `internal/`, `go.mod`,
  `go.sum`, `tools/coverage-gate.sh` sind Freistellungen) — ein
  `harness/image-hash.txt`-Nachzug ist für diesen Diff nicht fällig.
- **geprüft, ohne Befund: `ADR-0060`s Fitness-Function-Zeile „`.a-check`
  unverändert".** Sie ist auf den Gegenstand der ADR bezogen (die zwei neuen
  Pakete liegen in `adapters`/`ports`), und diese Globs sind im Diff
  tatsächlich unberührt. Ihre Begründungshälfte („kein neuer
  Hexagon-Schichten-Edge nötig") gilt weiter — die verletzte Aussage ist
  die der F-1, nicht diese. Dasselbe gilt für die gleichlautenden Zeilen in
  `ADR-0057`/`ADR-0061`.
- **geprüft, ohne Befund: die Auth-Konstanten im Client.** `main.go:35-39`
  setzt `authorizationMetadataKey`/`bearerPrefix` als dritte Fassung
  derselben Wertform (`internal/adapters/driving/grpc/interceptor.go:18-22`,
  HTTP-Middleware). Für einen **externen** Client ist die eigene Fassung
  richtig — er darf die unexportierten Konstanten des Servers weder
  importieren noch sollte er es; der Kommentar benennt die Kopplung
  (`SPEC-020`).
- **geprüft, ohne Befund: Traceability und ID-Schema.** Beide Betreffs
  tragen `LH-FA-SST-008`, `b835dde` zusätzlich `ADR-0060`; kein `SPEC-*`/
  `ARC-*` im Betreff; die drei Commit-Präfixe (`feat(harness)`, `docs(planning)`)
  folgen dem Bestand. `make commit-traceability` meldet im Gate-Lauf 5
  Commits, Betreffe ohne Struktur-ID (Exit 0).
- **geprüft, ohne Befund: DoD-Häkchen-Trennung im Nachzug `6862667`.** Gezogen
  werden genau vier Zeilen (Happy Path/Negative, `compose.yaml`+Runner,
  `make gates`, README-Doku); „Review durchgeführt", Closure-Notiz,
  Beobachtungs-Register, Risiko-Ausgänge und die drei Paarungen bleiben offen.
  Keine vorweggenommene Closure-Zeile, kein Übergriff in eine andere Rolle.
- **geprüft, ohne Befund: `tools/harness/grpcclient/main.go`,
  `tools/harness/run-integration-tests.sh` (additiver Block),
  `compose.yaml` (neuer Block), `.a-check.yml` (neuer Block), `harness/README.md`
  — über die oben benannten Punkte hinaus kein Befund.**

---

## Mutations- und Probenläufe (alle selbst ausgeführt)

Jeder Lauf ungefiltert in eine eigene Log-Datei, Exit-Code danach in einem
**eigenen, ungeketteten** Schritt ermittelt (`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify `v6.5.0` OK (54 Dateien); d-check 559 Dateien/0 Befunde (links/anchors/ids/matrix/versions/structure) + `commits`-Modul 0 Befunde; commit-traceability 5 Commits; a-check 0 Befunde; `coverage-gate: OK — 48.40 % ≥ 35 %` |
| `make test` | **0** | 28 Pakete `ok`, kein `FAIL` (`grpc`, `grpcstream`, `capture`, `bootstrap` grün) |
| `make a-check` (Stand `HEAD`) | **0** | `gesamt: 0 Befund(e)` — **und kein Abdeckungs-Hinweis mehr**; vor dem Diff stand dort „3 gescannte Datei(en) liegen in keiner Schicht" |
| `make a-check` + Vor-Fassung der `composition_root` (ohne `tools/**`) | **2** | `tools/harness/grpcclient/main.go:25: wrong-direction: (ohne Schicht) -> adapters (…/streamv1)` + derselbe Hinweis über drei Dateien — der rote Erstlauf des Nachzugs ist damit unabhängig reproduziert |
| `make a-check` + Probe-Datei `tools/harness/grpcclient/probe_scope.go` (Import `internal/application/usecase/capture`) | **0** | die Aufnahme gewährt dem ganzen Baum unbeschränktes Importrecht, ohne Hinweis — Beleg zu F-1 (b) |
| `bash -n tools/harness/run-integration-tests.sh` | **0** | Skript-Syntax des geänderten Runners |

**Rücknahme und Blatt-Identität:** `.a-check.yml` wurde aus der Sicherung
zurückgespielt, `git hash-object` = `44846bc5…` = `git rev-parse
HEAD:.a-check.yml`; die Probe-Datei ist gelöscht; `git status --porcelain`
ist leer. `tools/harness/grpcclient/main.go` hat unverändert
`c4ed07e8…`.

**Nicht ausgeführt und warum:** `make test-integration` (der E2E-Tier, in dem
der geänderte Runner läuft) — er setzt `make image` und einen vollständigen
Compose-Aufbau voraus und liegt in der Verifikations-Schicht; die Läufe
dieses Reports prüfen den Diff gegen Plan und Entscheidungen, nicht die DoD.
Die Aussagen zu F-2 und F-5 sind deshalb am **Code** geführt, nicht am Lauf.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Gate-Scope-Erweiterung ohne die Rolle, die
der Schlüssel deklariert · Erfolgsmeldung behauptet mehr als die Assertion
prüft · Spec-Stelle mit fremder Kennung adressiert · Aufschub-Adresse nimmt
die Sendung nicht an · Negativ-Hälfte ohne rot gesehenes Gegenbeispiel ·
Bereitschafts-Zeile als serverseitiges Signal gelesen

## Verdikt

**Merge-blockierend:** ja — F-1. Die Änderung an `.a-check.yml` ist der
Vollzug einer Gate-Lockerung, deren Träger in `ADR-0041` nicht benannt ist;
`AGENTS.md` §3.6 und §3.5 verlangen dafür ein Architect-Artefakt, keine
stille Konfigurationszeile. Die LOW/INFO-Findings blockieren nicht.

**Übergabe:** F-1 geht als **rückfragefrei entscheidungsfähige Frage** an die
**Architect-Rolle** (drei Antworten in §„Zum `.a-check.yml`-Verdikt"
ausformuliert, keine Code-Arbeit in diesem Report). Nach Modul 8
§Konflikt-Pfad ist bei der Antwort „Lockerung legitim, aber undokumentiert"
der Weg A → P → I vorgesehen (Folge-ADR plus Plan-Korrektur) — deshalb ist
hier **kein** Verdikt „Implementer-Rückkante entfällt" möglich. F-2 (eine
Zeile im Runner und die Sensor-Zelle), F-3/F-4 (Plan-Wortlaut) laufen als
Rückkante Reviewer → Implementer/Planner. Die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler; für F-1 ist
`BEO-PGC/a-check-null-abdeckung` der verwandte, **verkörperte** Eintrag — der
neue Befund betrifft dessen Gegenrichtung (der Hinweis verstummt, statt
Befunde zu liefern) und ist eine eigene Klasse. Dieser Report selbst ist ein
**Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und wird über Läufe hinweg nicht wieder gelesen. Er ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(`v6.5.0` · `regelwerk/modul-11-*.md`).

**DoD-Checkbox-Nachzug:** **kein Nachzug.** Die Zeile „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" bleibt im Slice-Plan offen. Die
Bedingung des Skills (§DoD-Checkbox-Nachzug ohne Fixrunde) ist nicht erfüllt:
mit F-1 läuft der Konflikt-Pfad des Moduls 8 weiter und kann nach der
Architect-Antwort eine Änderung an `.a-check.yml` und damit einen weiteren
Review-Lauf auslösen; F-2/F-3/F-4 liegen zudem in Dateien, die ein
Implementer-Lauf trägt. Der Nachzug gehört in den regulären Mechanismus
(Schritt 21 des Implementer-Workflows), nicht in diesen Report.
