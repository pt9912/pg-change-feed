# ADR-0084: Ein Sync-Gate nur für das Erzeugnis mit einer netzlosen, deterministischen Quelle

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-09-16

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als die Implementer-Läufe,
die die Erzeugnisse von Hand zurückgenommen und die Tabelle neu geschrieben
haben — Modul 8 §Rollen-Regeln). Der Zug ist der **vorgezogene Lese-Schritt**
der laufenden `welle-20`-Closure (Modul 6 §Wellen-Closure-Prozedur, Schritt 3).

**Bezug:** [`ADR-0060`](0060-grpc-streaming-mechanismus.md) (der
Protobuf-Generator, dessen Folgepflicht das neue Gate einlöst),
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (der Rollout mit
Pflicht-Report und Rollback-Artefakt),
[`ADR-0044`](0044-image-beleg-semantik.md) (der Image-Digest ist ein
**Lauf-Beleg**, kein Inhalts-Fingerabdruck),
[`ADR-0030`](0030-testpyramide.md) (Tier-Zuschnitt — der volle
Integrations-Lauf),
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md) §(a) und
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(bootstrap-aware Gates — das Muster für ein Gate, das nicht von Tag eins an
die Endstufe fährt),
`AGENTS.md` §3.1 (Docker-only) · §3.6 (Schwellen nur per ADR) · §4 (die
Gate-Tabelle) · `harness/README.md` §Sensors (Bindung je Target) ·
`Makefile` (`proto-generate`, `schema-rollout`), `Dockerfile` (Stufen `deps`,
`proto`, `coverage`), `harness/mk/coverage.mk`, `harness/mk/*.mk` ·
`tools/harness/run-integration-tests.sh` (der Erzeuger der
E2E-Abdeckungstabelle, `abdeckung_declare`/`abdeckung_go_zeilen_lesen`/
`abdeckung_schreiben`), `test/integration/integration_test.go`
(`TestAbdeckungstabelleZeilen`), `tools/schema/apply-rollout.sh`,
`tools/schema/plan.yaml` · `.d-check.yml` (die `structure`-Regel auf
`docs/user/e2e-abdeckung.md`) ·
`docs/plan/planning/observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/`
<!-- d-check:status-provenance -->
(die vier Belege) · `docs/reviews/review-slice-069.md` (F-6), <!-- d-check:status-provenance -->
`docs/reviews/verify-slice-069.md`, `docs/reviews/review-slice-074.md` <!-- d-check:status-provenance -->

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie `ADR-0043`,
`ADR-0044`, `ADR-0054`, `ADR-0071`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Eintrag steht bei 4×, und seine vier Belege zeigen drei Erzeugnisse
plus ein viertes.** Die Klasse: ein Teil des Repos besteht aus **erzeugten
Artefakten, die committet im Baum liegen**; für keines existiert ein Sensor,
der das Erzeugnis gegen seine Quelle hält. Die Beleg-Vorgänge sind `slice-069` <!-- d-check:status-provenance -->
(Protobuf-Code, als INFO eingeordnet), `slice-074` (die neue <!-- d-check:status-provenance -->
E2E-Abdeckungstabelle), `slice-082` (`tools/schema/plan.yaml`, bei **jedem** <!-- d-check:status-provenance -->
Rollout-Lauf verändert) und `slice-086` (derselbe Gegenstand erneut, weiter von <!-- d-check:status-provenance -->
Hand behandelt). `slice-086`s Beleg nennt daneben `harness/image-hash.txt` als <!-- d-check:status-provenance -->
„dieselbe Familie".

Damit sind **vier** Kandidaten zu prüfen, nicht drei:

| Erzeugnis | Quelle | Erzeuger |
|---|---|---|
| `internal/adapters/driving/grpc/streamv1/*.pb.go` (2 Dateien) | `proto/cdc/stream/v1/changestream.proto` | `make proto-generate` (Dockerfile-Stufe `proto`) |
| `docs/user/e2e-abdeckung.md` | `test/integration/integration_test.go` **und** `tools/harness/run-integration-tests.sh` | der Integrations-Runner |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | `tools/schema/schema.yaml` | `make schema-rollout` (`schema migrate --execute`) |
| `harness/image-hash.txt` | das gebaute Image | `make image` |

**(2) Der Prüf-Gedanke des Eintrags ist „Generator laufen lassen, `git diff
--exit-code` gegen das committete Erzeugnis".** Ob er trägt, entscheidet sich
an drei Eigenschaften des jeweiligen Generators — **netzlos**, **billig**,
**deterministisch**. Gelesen am Baum, nicht gemessen (Gegenstand dieser
Prüfung ist der Quelltext der Generatoren, kein Lauf):

- **Protobuf (netzlos im Lauf, deterministisch, eine Quelle).** Das Rezept von
  `make proto-generate` ruft `protoc` mit `--network none` auf; der erzeugende
  Container ist die Dockerfile-Stufe `proto`, die `protobuf-dev=31.1-r1` sowie
  `protoc-gen-go@v1.36.12` und `protoc-gen-go-grpc@v1.6.2` **gepinnt** trägt.
  Der Lauf ist damit netzlos; der **Build** dieser Stufe braucht Netz für
  `apk add` und die beiden `go install` — dasselbe Profil wie die heute schon
  von `make gates` gefahrene `coverage`-Stufe (deren `deps`-Vorgänger
  `go mod download` ausführt und die bei jedem Lauf mit
  `--no-cache-filter coverage` neu ausgewertet wird). Aus der gepinnten Quelle
  und dem gepinnten Generator ist das Erzeugnis eine **Funktion der Eingabe** —
  kein Umgebungsanteil geht in den Inhalt ein.
- **E2E-Abdeckungstabelle (zwei Quellen, eine davon nur im Lauf).** Der
  Erzeuger liegt **im** Integrations-Runner, zwischen den Phasen. Die
  **Go-Hälfte** ist für sich netzlos und billig: `TestAbdeckungstabelleZeilen`
  ist ein reiner Quelltext-Parser (`go/parser` über `test/integration/`,
  `runtime.Caller`) und braucht keine Datenbank. Die **Bash-Hälfte** dagegen
  ist **positional**: `abdeckung_declare` trägt `BASH_LINENO[0]` — die
  Aufrufzeile des Deklarations-Aufrufs **im Runner** — und sucht danach die
  erste spätere Zeile, die den Anker wörtlich trägt. Diese Positionen
  entstehen nur, während der Runner läuft; die Tabelle schreibt er nur bei
  inhaltlicher Abweichung (`cmp`). Ein voller Nachlauf ist der volle
  Integrations-Lauf (Compose-Umgebung, zuvor `make image`, Minuten).
- **`plan.yaml`/`down.sql` (nicht netzlos, nicht deterministisch über
  Umgebungen).** Das Rezept braucht `--network $(SCHEMA_ROLLOUT_NETWORK)` **und
  ein lebendes Ziel**; der geschriebene Report trägt seinen **Ziel-DSN** und
  einen `execution`-Block. Die drei Aufrufer übergeben drei verschiedene
  Container-Hosts — `cdc-store-test-pg` (`run-store-tests.sh`),
  `cdc-repl-test-pg` (`run-replication-tests.sh`), `cdc-test-postgres`
  (Compose/`run-integration-tests.sh`) —, und die committete Fassung trägt
  `cdc-test-postgres`. **Ein `git diff --exit-code` auf dieser Datei meldete
  damit die Umgebung als Drift**, nicht eine Abweichung von der Quelle; die
  Datei ist ein **Nebenprodukt des Sensor-Laufs** (so führen sie die
  Slice-Pläne selbst: „kein separat verfasster Inhalt") und wird bis heute bei
  jedem Commit von Hand zurückgenommen.
- **`harness/image-hash.txt` (kein Sync-Gegenstand per Entscheidung).**
  [`ADR-0044`](0044-image-beleg-semantik.md) erklärt den Digest ausdrücklich
  zum **Lauf-Beleg** eines `make image`-Laufs und **nicht** zum
  Inhalts-Fingerabdruck: ein Digest-Vergleich über Umgebungen oder Läufe
  entscheidet Staleness nicht. Ein Sync-Gate über dieser Datei wäre ein
  Widerspruch zu einer angenommenen ADR.

**(3) Die drei Größen, die den Ausschlag geben, sind orthogonal.** *Netzlos*
ist eine Eigenschaft des Generator-Laufs (und auf kaltem Cache des Builds),
*deterministisch* eine des Erzeugnisses über Umgebungen (eine Quelle und ein
gepinnter Generator gegen einen Umgebungsanteil wie den Ziel-DSN), *billig*
eine des marginalen Aufwands in `make gates` (ein Build mit warmem Cache und
ein Protokoll-Lauf gegen einen vollen Testlauf). Nur ein Kandidat erfüllt alle
drei.

## Entscheidung

Wir wählen: **Ein Sync-Gate bekommt nur das Erzeugnis mit einer netzlosen,
deterministischen und einzigen Quelle.** Vier Festlegungen:

**1. Der Protobuf-Code bekommt ein Gate.** Ein neues `make`-Ziel (Vorschlag:
`generated-sync`, der Name gehört dem umsetzenden Zug) läuft den
**gepinnten** Generator gegen die **committete** `.proto`-Quelle und
vergleicht das Ergebnis mit dem committeten Erzeugnis. Es hängt an
`GATE_CHECKS` und läuft damit in `make gates`. Zwei Bedingungen gehören zum
Gate, nicht in seinen Aufrufer:

- **Das Gate schreibt den Arbeitsbaum nicht.** Es erzeugt in ein
  **Temp-Verzeichnis** (mit derselben Modul-Layout-Relation, die `--go_out` mit
  `--go_opt=module=…` herstellt) und vergleicht dort; das Ziel
  `make proto-generate` schreibt **in-place** in den Bind-Mount und ist damit
  **nicht** als Prüf-Schritt geeignet — ein Gate, das erst schreibt und dann
  `git diff` liest, lässt den Baum schmutzig zurück und ist beim zweiten Lauf
  grün.
- **Der Befund nennt den Diff.** Ein rotes Gate zeigt, **welche** Datei und
  welche Zeile abweichen; „Erzeugnis nicht synchron" ohne Datei ist ein
  Befund, den der nächste Lauf nicht auflösen kann.

**Was dieses Gate prüft — und nur das:** die **Paarung** „committetes
Erzeugnis = Ausgabe des gepinnten Generators aus der committeten Quelle". Es
prüft **nicht**, dass der erzeugte Code kompiliert oder sich richtig verhält
(das tun `make test` und `make image`), und **nicht**, dass die `.proto` den
Draht-Vertrag richtig beschreibt (das prüfen die Belege zu `LH-FA-SST-008`).

**2. Die E2E-Abdeckungstabelle bekommt heute kein Gate — die Bauform, die
eines trüge, wird benannt.** Der Grund ist die **Zwei-Quellen-Positionalität**
aus Kontext (2): ein Gate über der Go-Hälfte allein wäre ein **halber**
Wächter — eine neue oder entfernte Phase im Runner verschiebt die `Ort`-Spalte
der Bash-Hälfte, und der halbe Wächter bliebe grün. Ein halber Wächter, dessen
Grün nur die halbe Zusage deckt, ist die Klasse, die dieses Repo an anderer
Stelle als benannte Lücke führt; sie wird hier **nicht** eingeführt.

Die Bauform, die das Gate trägt, ist die **eine Quelle**: Erzeuger und Prüfer
rufen **denselben** Code-Pfad für **beide** Hälften. Heute geht das nicht, weil
die Deklarationen zwischen den Phasen liegen und ihre Position aus dem Lauf
selbst stammt; eine Extraktion müsste den Runner so umbauen, dass die
Deklarationen (samt Position) aus **einer** Funktion kommen, die beide Seiten
lesen. Das ist ein **eigener Vorgang mit eigenem Schnitt** — ein
**Folge-Slice-Vorschlag**, keine Beigabe dieses Zuges. Eine Kennung nennt
dieser Zug **nicht**: eine Adresse, die die Sendung nicht annehmen kann, ist
keine.

**3. `plan.yaml`/`down.sql` bekommen kein Gate — die Begründung ist
strukturell, nicht ökonomisch.** Nicht nur teuer: ein Diff-Check wäre
**falsch-positiv**, weil das Erzeugnis einen Umgebungsanteil trägt (Kontext
(2)). Ein Sync setzte eine **kanonische Ziel-Umgebung** voraus, die es nicht
gibt; ohne sie ist die committete Fassung die Ausgabe *eines* der drei
Aufrufer. Die Datei bleibt damit, was sie ist: ein **Nebenprodukt des
Sensor-Laufs**, das committet wurde. Die tiefere Frage — ob ein Lauf-Artefakt
überhaupt in den Baum gehört (statt ignoriert und durch einen geprüften
Auszug ersetzt zu werden) ist **nicht** Gegenstand dieser ADR; sie wird
benannt, nicht verdeckt (siehe §Was diese ADR nicht entscheidet).

**4. `harness/image-hash.txt` ist kein Sync-Gegenstand.** Die Semantik aus
[`ADR-0044`](0044-image-beleg-semantik.md) — Lauf-Beleg, kein
Inhalts-Fingerabdruck — schließt einen Diff-Check gegen die Quelle aus. Ein
solcher Sensor wäre eine stille Umdeutung einer angenommenen Entscheidung; wer
ihn will, braucht zuerst eine Folge-ADR zu `ADR-0044`.

### Was das Gate nicht fangen kann — benannt

- **Protobuf:** eine **falsche** `.proto` (die Quelle ist synchron, aber die
  Aussage über den Draht stimmt nicht) — der Beleg liegt im Integrations- und
  Unit-Tier, nicht hier. Und: ein **bewusst** von Hand nachbearbeitetes
  Erzeugnis — das Gate macht daraus einen roten Lauf, was der gewünschte
  Befund ist.
- **E2E-Abdeckungstabelle:** ihre Veraltung **zwischen** zwei
  Integrations-Läufen, in **beiden** Hälften — die `Ort`-Spalte verschiebt
  sich mit jeder Einfügung oberhalb einer `func TestE2E*` **oder** einer
  Runner-Phase. Der Wächter bleibt bis zur Extraktion die Disziplin, die sie
  bis heute getragen hat; `docs-check`s `structure`-Regel sichert die
  **Existenz** der Datei, nicht ihren Inhalt.
- **`plan.yaml`/`down.sql`:** sie altern still gegen `tools/schema/schema.yaml`
  — sichtbar würde das erst an einem Rollout, der aus der Quelle neu schreibt.
  Der committete Stand ist zudem das Erzeugnis **einer** Umgebung.
- **`harness/image-hash.txt`:** Staleness gegen den Baum ist per `ADR-0044`
  **keine** Aussage, die der Beleg trägt; wer sie prüft, prüft etwas anderes
  (Inhalts-Streit → sha256 des extrahierten Binaries).

### Was diese ADR nicht entscheidet

- **Nicht die Disposition von `plan.yaml`/`down.sql`** (ob ein Lauf-Artefakt
  im Baum bleibt, ignoriert oder durch einen geprüften Auszug ersetzt wird) —
  das ist eine eigene Entscheidung mit eigenem Adressaten (Planner/Implementer,
  `Makefile`-Vertrag) und würde den Pflicht-Report-Vertrag aus
  [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) berühren.
- **Nicht die Extraktion des E2E-Erzeugers** — sie ist benannt
  (Festlegung 2), ihr Schnitt ist Planner-Arbeit.
- **Nicht die Coverage-Schwelle und nicht die Zusammenstellung von `make
  gates`** darüber hinaus; §3.6 ist unberührt (ein neues Gate ist keine
  Lockerung).
- **Kein Carveout.** Es gibt keinen roten Lauf, den diese ADR auflöst.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; die Bindung bleibt Disziplin | kein Aufwand, kein neuer Gate-Lauf | die Klasse steht bei 4×; `slice-082` hat gezeigt, dass ein Lauf die committete Datei **planmäßig** verändert und der Implementer sie jedes Mal von Hand zurücknehmen muss — Disziplin ist genau das, was die vier Belege als nicht tragend ausgewiesen haben <!-- d-check:status-provenance --> |
| B — Gates für alle vier Erzeugnisse | maximale Deckung, ein Muster | für `plan.yaml` **falsch-positiv** (Umgebungsanteil), für die E2E-Tabelle nur als halber Wächter baubar, für `image-hash.txt` ein Widerspruch zu `ADR-0044`; ein Gate, dessen Grün mehr behauptet als es prüft, ist teurer als keines |
| C — nur die **Go-Hälfte** der E2E-Tabelle prüfen (billig, deckt die häufigste Drift) | die Einfügung oberhalb einer `func TestE2E*` ist die beobachtete Drift (`slice-086`: alle 13 Verweise um +4) | eine neue **Phase** im Runner verschiebt die `Ort`-Spalte ebenso und bliebe grün — der Wächter deckte die halbe Zusage; die Zwei-Quellen-Form ist genau der Defekt, gegen den `slice-074` die Anker gebaut hat | <!-- d-check:status-provenance -->
| D — `plan.yaml` per kanonischem Ziel-DSN deterministisch machen | die Datei wäre vergleichbar | sie braucht weiter ein **lebendes** Ziel je Gate-Lauf (Compose), und der `execution`-Block bleibt lauf-gebunden; die Änderung berührt den `ADR-0043`-Vertrag des Pflicht-Reports |
| E — `plan.yaml`/`down.sql` aus dem Baum nehmen (ignorieren) | beendet die Hand-Rücknahme und die stille Alterung | ändert den Pflicht-Report-Vertrag und die Rolle der Datei als Beleg-Quelle (die Slice-Pläne zitieren sie); eine eigene Entscheidung, nicht Beigabe |
| **F — Gate nur für den Protobuf-Code; die drei übrigen mit benannter Begründung entschieden (gewählt)** | das einzige Erzeugnis mit einer netzlosen, deterministischen Quelle bekommt einen echten Wächter; die drei übrigen bekommen je ihren Grund statt einer stillen Lücke; kein Widerspruch zu `ADR-0044`, kein halber Wächter | drei Vorgänge der Klasse bleiben ohne Sensor; die E2E-Tabelle wartet auf einen Folge-Slice; ein neuer Build-Schritt im Gate-Bündel |

**Fazit:** F. B verletzt zwei Zusagen, C deckt die halbe Zusage, D und E wären
eigene Vorgänge mit fremdem Vertrag, A lässt eine 4×-Klasse ohne Wächter.
F gibt dem einen tragfähigen Kandidaten sein Gate und den drei anderen ihren
Grund.

## Konsequenzen

- Positiv: Die älteste der drei Lücken (`slice-069`, als INFO eingeordnet und <!-- d-check:status-provenance -->
  ins Register verwiesen statt still geschlossen) bekommt einen mechanischen
  Wächter; die Bindung des Protobuf-Erzeugnisses an seine Quelle ist keine
  Disziplin mehr.
- Positiv: Die Entscheidung **gegen** einen Sensor ist an drei Stellen
  begründet, nicht verschwiegen — jede der drei trägt einen eigenen Grund
  (falsch-positiv, halber Wächter, widersprochene ADR).
- Positiv: `AGENTS.md` §3.6 ist unberührt — ein neues Gate ist keine
  Schwellen-Senkung.
- Negativ mit Grenze: `make gates` bekommt einen weiteren Build-Schritt (warm
  gerechnet billig — die `coverage`-Stufe wertet bei jedem Lauf den vollen
  Testlauf neu aus; auf **kaltem** Cache braucht der `proto`-Build Netz für
  `apk add`/`go install`, wie die `deps`-Stufe heute schon).
- Negativ mit Grenze: Die drei übrigen Erzeugnisse bleiben ohne Sensor; ihre
  Veraltung ist **benannt** (§Was das Gate nicht fangen kann), nicht behoben.
  Für die E2E-Tabelle ist die Bauform des künftigen Wächters festgelegt (eine
  Quelle für beide Hälften) — die Form, nicht die Umsetzung.
- Folgepflicht (Implementer-Zug, Tooling + Doku, **kein** Produkt-Code): das
  neue Ziel im `Makefile`-Umfeld anlegen (Temp-Verzeichnis + Vergleich +
  Diff-Ausgabe), an `GATE_CHECKS` hängen, und seine Bindung nachtragen:
  `harness/README.md` §Sensors (Zeile mit Vertrag und Bindung `ADR-0084` ·
  `ADR-0060`) und `AGENTS.md` §4 §Quality Gates (die Tabelle listet auf, sie
  definiert nicht).
- Folgepflicht (Planner-Zug): den **Folge-Slice-Vorschlag** aus Festlegung 2
  schneiden — Extraktion des E2E-Erzeugers, damit Erzeuger und Prüfer denselben
  Code-Pfad lesen —, wenn der Planner ihn aufnimmt. Ohne ihn bleibt die
  Tabelle bewusst bei „Disziplin + `structure`-Regel".
- Folgepflicht (Planner-Zug, Register): der Eintrag erhält im Lese-Schritt der
  `welle-20`-Closure seinen Ausgang **`verkörpert`** für den Protobuf-Träger
  (Zielort: das neue Gate, Träger-Kennung `ADR-0084`) und **`geplant`** für die
  E2E-Tabelle (Kennung = der schreibende Folge-Slice) — `plan.yaml` und
  `image-hash.txt` sind mit dieser ADR **entschieden** (kein Gate, Grund
  benannt). Sobald das Gate gebaut ist, trägt die Zeile den Herkunfts-Anker des
  schreibenden Vorgangs (`seit welle-20` bzw. `seit slice-<NNN>`).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| neues Ziel (Vorschlag `generated-sync`) — Dockerfile-Stufe `proto` (gepinnt), `protoc` mit `--network none` | Der committete Protobuf-Code ist **byte-gleich** der Ausgabe des gepinnten Generators aus der committeten `.proto`-Quelle; erzeugt wird in ein Temp-Verzeichnis, verglichen mit dem Baum, Befund mit Diff | `make generated-sync` (in `make gates`) |
| — (Disziplin, **kein** Sensor) | Die E2E-Abdeckungstabelle wird von ihrem Erzeuger geschrieben; ihr Nachlauf ist der volle Integrations-Lauf. `docs-check`s `structure`-Regel sichert ihre **Existenz**, nicht ihren Inhalt | `make test-integration` (kein Gate) |
| — (entschieden **gegen** einen Sensor) | `plan.yaml`/`down.sql` tragen einen Umgebungsanteil (Ziel-DSN, `execution`) — ein Diff-Check meldete die Umgebung als Drift; `harness/image-hash.txt` ist per `ADR-0044` ein Lauf-Beleg | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger:

**(a)** Der E2E-Erzeuger wird extrahiert (eine Quelle für beide Hälften) —
dann ist das Gate aus Festlegung 2 baubar und wird gebaut; diese ADR ist dann
in **dieser** Klausel zu schärfen (Folge-ADR), nicht zu dehnen.

**(b)** `plan.yaml` bekommt ein **kanonisches** Ziel (oder der Report verliert
seinen Umgebungsanteil) — dann ist die Begründung aus Festlegung 3 weggefallen,
und der Kandidat ist neu zu prüfen. Ebenso, wenn ein Rollout-Lauf die
committete Datei **nicht** mehr verändert.

**(c)** Der Generator des Protobuf-Codes verliert seine Pinnung (Tag statt
Pin, ungepinnte Plugin-Version, Generator mit Umgebungsanteil) — dann ist das
Gate nicht mehr deterministisch, und `ADR-0060`s Folgepflicht ist vorrangig zu
schließen.

**(d)** Ein Gate-Lauf wird durch das neue Ziel **unzumutbar teuer**
(Beobachtung: `make gates` überschreitet die Grenze, ab der es nicht mehr vor
jedem Handoff läuft) — dann greift das Muster aus
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md) §(a): das Gate
bekommt eine Einstiegsstufe (Prüfung im PR, nicht lokal) **mit** einem
Hochschalt-Trigger, nicht stiller Abschaltung.

Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-16 | Accepted — Anlass: `BEO-PGC/generierte-artefakte-ohne-sync-sensor` steht bei 4× (`slice-069`, `-074`, `-082`, `-086`). Entscheidet: ein Sync-Gate **nur** für den Protobuf-Code (netzlos im Lauf, pin-genau, deterministisch, eine Quelle; Erzeugung in ein Temp-Verzeichnis, damit das Gate den Baum nicht schreibt); die E2E-Abdeckungstabelle **ohne** Gate (Zwei-Quellen-Positionalität — die Bauform des künftigen Wächters ist benannt und als Folge-Slice-Vorschlag geführt); `plan.yaml`/`down.sql` **ohne** Gate (Umgebungsanteil im Erzeugnis → Diff-Check wäre falsch-positiv); `harness/image-hash.txt` **ohne** Gate (Lauf-Beleg per `ADR-0044`). Löst die Folgepflicht des Generators aus `ADR-0060` ein und berührt `ADR-0043`/`ADR-0044` nicht | die vier Beleg-Dateien des Eintrags, `docs/reviews/review-slice-069.md` F-6, `docs/reviews/verify-slice-069.md` <!-- d-check:status-provenance -->, `Makefile`/`Dockerfile` (gelesen, nicht gemessen) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0084` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
