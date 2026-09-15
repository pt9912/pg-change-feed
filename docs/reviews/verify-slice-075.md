# Verifikationsbericht: slice-075 — 2026-09-15

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-075` §1–§8) und die bindenden Entscheidungen
([`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
vollständig; [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) in
der Dauerhaftigkeits-Aussage superseded, sonst gültig,
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)). Nicht gegen den Diff
als solchen (Reviewer-Aufgabe;
[`review-slice-075`](review-slice-075.md) vollständig gelesen, aber nur als
Kontext) und nicht gegen realen Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8,
Stand `HEAD = 31e698c`), die vollständige
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md), das
[Architect-Verdikt](architect-verdict-spaltenausschluss-dauerhaftigkeit.md),
den Review-Report, die drei Vorgänger `slice-066`/`slice-067`/`slice-068` in
`docs/plan/planning/done/` sowie den tatsächlichen Diff `6be714d^..HEAD`
(der Slice-Anteil ist genau `6be714d`, 13 Pfade). Die Sensoren wurden in
dieser Sitzung **eigenständig real ausgeführt** — Exit-Code je in einem
eigenen, ungepipten Schritt (`AGENTS.md` §3.9); kein Beleg des Implementers
oder Reviewers wurde ungeprüft übernommen. **Fünf** Mutationen sind eigene,
am Code gesetzt.

**Wichtige Randbedingung:** Dieser Slice fällt in die **Befundlage F-3** des
Reviews selbst. Der Reviewer hat den realen `make test-integration`-Lauf
**nicht** nachgefahren (sein Report sagt das ausdrücklich). Genau diese
Lücke schließt dieser Bericht — der Kern-Beleg des Slice (Neustart) wurde
hier **erstmals unabhängig** gefahren und zusätzlich in einem eigenen
E2E-Mutationslauf als nicht-vakuum belegt (§3, Mutation E).

**Gegenstand:** `slice-075` zum Stand `HEAD = 31e698c`. Innerhalb
`6be714d^..HEAD` liegen neben dem Slice-Commit `6be714d` der Review-Commit
`31e698c` und sechs unabhängige Commits anderer Züge (`7ef760b`, `f90ba96`,
`b37cd9f`, `09c211c`, `87a2b9d`, `17596fa`) — sie sind **nicht** Gegenstand.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst — die implementierungs-/review-/dokubezogenen Zeilen sind
Prüfgegenstand, die Closure-Zeilen **müssen** offen bleiben.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Ausschlussstand aus den `applied`-Zeilen abgeleitet, bei **jedem** Bindungs-Anlegen mitgeführt — Startpfad **und** Aktivierungs-Zweig; Belege in `internal/bootstrap/` | **erfüllt, selbst reproduziert + rot gesehen** | Drei Tests existieren real (`administration_internal_test.go:742/786/843`). Eigener `make test` Exit **0**. Eigene Mutationen A/B (§3) reißen genau diese Tests einzeln auf. |
| 2 | Lesefähigkeit am Outbound Port (`ARC-004`, Zuschnitt `ADR-0034`), gegen reale PostgreSQL erprobt inkl. deterministischer Reihenfolge | **erfüllt, selbst reproduziert + rot gesehen** | Methode sitzt am bestehenden `ColumnExclusionPort` (`columnexclusion.go`), kein zweiter Port. Eigener `make test-store` Exit **0** (u. a. `postgresstorage` 4,18 s real). Eigene Mutationen C/D (§3) reißen den Store-Test einzeln auf. |
| 3 | E2E-Beleg (`LH-QA-SEC-004`): Neustart lässt einen zuvor beantragten Ausschluss wirksam | **erfüllt, selbst reproduziert + rot gesehen** | Eigener `make test-integration` Exit **0**; die Zeile *Spaltenausschluss-Neustart-Beleg* steht real in der Ausgabe (`run-integration-tests.sh:510`). Eigene E2E-Mutation E (§3) macht **genau diese** Zusage rot. |
| 4 | `make gates` grün | **erfüllt, selbst reproduziert** | Eigener Lauf Exit **0** (§2). |
| 5 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | [`review-slice-075`](review-slice-075.md) vorhanden; das Häkchen zieht der Review-Commit `31e698c` in **genau diesem** Commit nach (§6). Verdikt 0 HIGH / 0 MEDIUM / 1 LOW / 2 INFO deckungsgleich mit §2. |
| 6 | Doku-Update: `harness/README.md`-Werkzeuge-Zeile um den Beleg-Baustein; `SPEC-019`-Fließtext zur Bedeutung von `applied` plus Historie-Zeile | **erfüllt** | `harness/README.md:132` trägt den neuen Abschnitt; `spec/pflichtenheft.md:339-352` trägt den Fließtext, `:491` die Historie-Zeile. |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<bei Closure>`-Platzhalter (gelesen). Planner-Arbeit nach diesem Bericht. |
| 8 | Reconciliation-Register — entfällt (Greenfield) | **korrekt offen, Entfall trägt** | `docs/plan/planning/reconciliation.md` existiert real **nicht**; Repo durchgehend GF (`harness/conventions.md` `*`/`PGC`). |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger/evidence/` trägt real **nur** `review-slice-067.md`; `evidence/slice-075.md` fehlt noch — gehört in die Closure. Der Zieleintrag steht real (Zustand `geplant`, `…/evidence/review-slice-067.md` = 1×). |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Alle drei §6-Einträge tragen wörtlich `<bei Closure zuzuweisen>`; keines vorzeitig geschlossen. |
| 11 | Die drei Paarungen | **korrekt offen** | `slice-075` liegt real in `in-progress/`; das Kopf-Feld sagt `ohne Welle`. Die Prüfung verweist korrekt auf die nächste Wellen-Closure (Modul 6: die Wellen-Closure liest auch wellenlose Slices). |

**Ergebnis §1:** Die sechs implementierungs-/review-/dokubezogenen DoD-Punkte
(1–6) sind real erfüllt; die Punkte 1–3 wurden **selbst reproduziert** und
über fünf eigene Mutationen als einzeln tragend belegt. Die fünf
Closure-Punkte (7–11) sind korrekt noch offen und nicht vorweggenommen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungekettenen** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | d-check 583 Dateien / 0 Befunde; `--enable commits --range HEAD~5..HEAD` 0 Befunde; `commit-traceability: OK`; a-check gesamt 0 Befunde; `coverage-gate: OK — 49.30 % ≥ 35 %`; baseline-verify grün |
| `make test` | **0** | `go test -race ./...` im gepinnten Toolchain-Container; alle Pakete `ok`, `internal/bootstrap` dabei |
| `make test-store` | **0** | reale PostgreSQL im Testcontainer; `postgresstorage` 4,18 s `ok` (inkl. `TestTableActivationExcludedColumnsDerivesAppliedColumnRequests`) |
| `make image` | **0** | **reproduziert exakt** `sha256:9b7937d1…` — denselben Digest, den der Slice-Diff in `harness/image-hash.txt` committet. Der Digest ist kein Staleness-Beweis (`ADR-0044`), aber „derselbe Baum → derselbe Digest" ist eine unabhängige Korroboration der Reihenfolge `make image` → `make test-integration`. |
| `make test-integration` | **0** | voller Compose-Stack-Lauf; u. a. die Zeile *Spaltenausschluss-Rundlauf …* **und** die Kern-Zeile *Spaltenausschluss-Neustart-Beleg (`ADR-0065`, `LH-QA-SEC-004`) — der reale Container-Neustart ließ den dauerhaften Ausschlussstand wirksam werden: die danach erfasste Change (id=3) trägt `secret` nicht im Row Image und den Wert `ColumnAfterRestartSentinel` nirgends* |
| `make doc-commits RANGE=6be714d^..HEAD` | **0** | 583 Dateien, 0 Befunde — jeder Commit des Fensters ist kennungstragend |
| `make doc-immutable RANGE=6be714d^..HEAD` | **0** | 0 Befunde — keine `Accepted`-ADR überschrieben |

Zweiter, abschließender `make gates`-Lauf **mit dieser Bericht-Datei im Baum**: Exit **0**, d-check **584** Dateien / 0 Befunde — der Bericht selbst bricht weder `docs-check` noch die seit dem jüngsten Architect-Zug verkörperte Verweisform-Regel (`.d-check.yml`).

`git status --porcelain` nach allen Sensor- und Mutationsläufen: **leer** (der
lokal reproduzierte Rollout-Report `tools/schema/plan.yaml` wurde auf den
committeten Stand zurückgesetzt).

## 3. Mutations-Stichprobe — DoD-Zusagen real rot gesehen (Verifier-only-Nachweis)

Fünf **eigene** Mutationen, jede einzeln, jede mit anschließender Rücknahme
und Blatt-Prüfung (`git hash-object` = `git rev-parse HEAD:<pfad>`;
`git status --porcelain` nach jeder Rücknahme leer).

| # | Mutation | Ort | Erwartung | Ergebnis |
|---|---|---|---|---|
| A | Startpfad liest den Stand nicht (`ExcludedColumns` durch leere Map ersetzt) | `internal/bootstrap/wiring.go:351-354` | `make test` rot am Startpfad-Test | **rot**, Exit 2 — `TestActivatedTableBindingsCarriesExcludedColumns`: `ExcludedColumns = [], wollen [secret]` |
| B | Aktivierungs-Zweig liest den Stand nicht | `wiring.go:1171-1174` | `make test` rot am Zyklus-Test | **rot**, Exit 2 — `TestProcessAdministrationRequestsDisableEnableCycleRestoresExclusion`: `Row Image nach dem disable/enable-Zyklus = {"id":"1","secret":"geheim"}, wollen ohne den ausgeschlossenen Schlüssel secret` (dazu der `MarksFailedWhenExclusionReadFails`-Test) |
| C | `ORDER BY` ohne den Zweitschlüssel | `queries/queries.go` (`SelectAppliedColumnRequests`) | `make test-store` rot, weil die Reihenfolge bei gleichem `requested_at` kippt | **rot**, Exit 2 — `TestTableActivationExcludedColumnsDerivesAppliedColumnRequests`: `Ausschlussstand public.orders_exclusion_derivation = [], wollen [secret]` |
| D | `include_column` als No-op | `postgresstorage/tableactivation.go` (`ExcludedColumns`) | `make test-store` rot an der Einschluss-Wirkung | **rot**, Exit 2 — derselbe Test: `… = [secret], wollen keinen Eintrag (include_column hat den letzten Namen genommen)` |
| E | **E2E-Kern:** Startpfad ohne den Stand, Bild neu gebaut (`make image`), voller Compose-Rundlauf | `wiring.go:351-354` | rot **genau** am Neustart-Beleg, nicht am Live-Reload | **rot** — `make test-integration` Exit 2; der *Spaltenausschluss-Rundlauf* (Live-Reload, `run-integration-tests.sh:436`) blieb **grün**, die Kern-Zeile *Neustart-Beleg* brach: `die nach dem Neustart erfasste Change (id=3, feed_e2e_column_exclusion) trägt den ausgeschlossenen Spaltenschlüssel secret weiterhin im Row Image (der Ausschlussstand wurde beim Prozessstart nicht wiederhergestellt)` |

Nach Mutation E: `wiring.go` und `harness/image-hash.txt` per `git checkout`
zurückgenommen, `make image` erneut gefahren — das `:dev`-Image trägt wieder
`sha256:9b7937d1…` (derselbe Digest wie committet), `git status --porcelain`
leer.

**Ergebnis §3:** Alle vier Kern-Zusagen sind **einzeln** tragend. Mutation E
ist der stärkste Beleg und schließt zugleich die Kern-Lücke des Reviews: Der
Neustart-Beleg ist **nicht vakuum**. Er hängt am Startpfad
(`activatedTableBindings`) und nicht am Live-Reload-Pfad — die zwei Hälften
sind unabhängig sensibel (Mutation E färbt nur die Neustart-Hälfte rot, die
Live-Reload-Hälfte bleibt grün). Zugleich beweist E, dass der `docker restart`
real einen Prozess-Neustart auslöst: ohne Neustart bliebe der Stand im
laufenden Assembler stehen und der Beleg wäre auch mutiert grün.

## 4. Entscheidungs-Konformität — `ADR-0065` gegen den Code

Alle vier Festlegungen am Diff nachgeprüft:

1. **Tabellen-scoped, ein Mechanismus für beide Auslöser.** Beide Pfade, die
   eine Bindung anlegen, tragen denselben abgeleiteten Stand:
   `activatedTableBindings` (`wiring.go:346-370`, je Quelle) und der
   Aktivierungs-Zweig (`wiring.go:1171-1179`). Beide rufen **dieselbe**
   Fähigkeit; Mutationen A/B zeigen beide Aufrufer als einzeln tragend. Die
   `Assembler`-Bindung bleibt Laufzeit-Cache.
2. **Abgeleitet aus der Antrags-Historie, keine zweite Quelle.** Der einzige
   neue Zugriff ist ein reines `SELECT`
   (`SelectAppliedColumnRequests`): geprüft — kein `INSERT`/`UPDATE`/`DELETE`/
   `ALTER`/`CREATE` in den Produktionsdateien des Diffs. `tools/schema/` ist
   im Fenster unberührt (`git diff 6be714d^..HEAD -- tools/schema/schema.yaml
   tools/schema/nacharbeit-administration.sql` leer): **kein neues
   Schema-Objekt, kein zweiter Schreibpfad**. Die Reihenfolge
   `ORDER BY requested_at, administration_request_id` ist total geordnet
   (die ID ist PK); Mutation C zeigt den Zweitschlüssel als tragend.
3. **`applied` heißt dauerhaft vermerkt.** Der `SPEC-019`-Fließtext
   (`spec/pflichtenheft.md:339-352`) spricht genau die vier Hälften aus:
   dauerhaft vermerkt · `applied`-Zeilen als einzige Herkunft ·
   `requested_at` mit Zweitschlüssel · `applied` statt `failed` ohne laufende
   Bindung. Zur Attribution („Planner-/Architect-Zug" in `ADR-0065` §Folgepflicht
   vs. Slice-DoD): Der Slice schreibt den **in der ADR festgelegten** Wortlaut
   aus, entscheidet inhaltlich nichts, und §3 benennt die Abweichung
   ausdrücklich statt still — die Sicht des Reviews trägt, ich schließe mich an.
4. **Wirkort und Rückkanal unverändert; In-Prozess-Erhalt unangetastet.**
   `internal/adapters/driving/replication/mapper/**` ist über das ganze
   Fenster **nicht** im Diff (`git diff --stat 6be714d^..HEAD` leer):
   `AddBinding`-Merge (`mapper.go:410-412`) und `setSchemaVersion`
   (`mapper.go:422-431`) tragen unveränderten Quelltext. Die Rückkanal-Hälfte
   (`applied`/`failed`/`pending`) bleibt geschlossen.

`ADR-0065` §Fitness Function nennt zwei Zeugen — `make test` (Bindungs-Neuaufbau
und `disable`/`enable`-Zyklus) und `run-integration-tests.sh` (Neustart) —,
beide sind real belegt (§2/§3).

## 5. Plan-vs-Code-Diff

- **§3-Tabelle deckungsgleich.** Alle sieben genannten Pfade sind im Diff; die
  sechs weiteren Pfade sind durch die fünf §3-Nachzüge gedeckt
  (`queries.go`-Unterpaket, `administrationrequest_test.go`, zwei
  Use-Case-Fakes, `spec/pflichtenheft.md`, `harness/image-hash.txt`) — keine
  unerklärte Datei, keine fehlende.
- **Der `queries.go`-Pfad-Nachzug trägt:** die SQL-Texte liegen real im
  Unterpaket (`queries/queries.go`), nicht in `tableactivation.go`.
- **Der „je Quelle"-Nachzug trägt:** die Fähigkeit liest **je Quelle**
  (`ExcludedColumns(ctx, source) (map[string][]string, error)`), beide
  Aufrufer bedienen dieselbe Rückgabe — genau die von `ADR-0065` bezifferte
  Form („eine Abfrage je Quelle, nicht je Tabelle").
- **Der `SPEC-019`-Nachzug trägt:** der Fließtext steht in `SPEC-019` selbst
  (vor `SPEC-020`), plus Historie-Zeile in §7 (`:491`).
- **Kein Fremd-Schema-Objekt, kein Fremd-Wirkort:** `compose.yaml` und
  `tools/schema/nacharbeit-administration.sql` sind nicht im Diff.

## 6. Urteile zu F-1/F-2/F-3 und zur DoD-Häkchen-Trennung

**F-1 (fehlender `· seit slice-NNN`-Vermerk am neuen Beleg-Baustein) —
eigenständig als LOW bestätigt, nicht blockierend.** Unabhängige Zählung:
`harness/README.md` trägt 15 `· seit slice-NNN`-Vermerke; der neue Abschnitt
(„Zusätzlich ein Spaltenausschluss-Neustart-Beleg …") ist der **einzige** der
„Zusätzlich …"-Bausteine ohne ihn (`harness/README.md:132`). Die Herkunft ist
über die im selben Zellinhalt zitierte `ADR-0065` **auflösbar** — es fehlt die
Gleichförmigkeit, nicht die Information (`AGENTS.md` §3.7 verlangt *ein*
auflösbares Herkunftsfeld, nicht die `seit`-Schreibweise). Kein Gate fängt es.
Ein Ein-Zeilen-Nachzug `· seit slice-075` ist möglich, aber nicht erzwungen.

**F-2 (Lesefehler des Standes im Aktivierungs-Zweig hinterlässt „aktiviert,
aber ungebunden") — eigenständig als INFO bestätigt, fail-closed.** Am Code
nachgezogen: der Aktivierungs-Zweig liest den Stand **nach** dem erfolgreichen
`EnableTableUseCase` (`wiring.go:1146-1179`); scheitert der Lese-Zug, endet der
Antrag `failed`, während Bindungs-Zeile und Publication real geschrieben sind.
Das ist **keine** Verletzung von `ADR-0065` Festlegung 3: die sagt „gegen eine
existierende Spalte **ohne laufende Bindung** endet *nicht* `failed`" — das ist
der **Erfolgs**-Fall, nicht der **Lese-Fehler**-Fall. Die Form ist überdies
vorbestehend (derselbe Zustand entsteht beim `Registered`-/
`CurrentVersion`-Fehlschlag); die Richtung ist fail-closed (lieber keine
Bindung als eine ohne den geführten Ausschluss), und
`TestProcessAdministrationRequestsMarksFailedWhenExclusionReadFails` pinnt sie
(im eigenen Mutationslauf B mitrot). Kein Handlungsbedarf am Diff.

**F-3 (die Nicht-Vakuum-Aussage hängt an einem Lauf ohne abgelegten Beleg) —
der Aussage fehlt kein Beleg mehr; ich habe ihn ersetzt.** Das Review konnte
F-3 nur per **Code-Vergleich** `6be714d^` ↔ `6be714d` stützen (der Eltern-Stand
liest in `activatedTableBindings` keinen Stand — am `git show` reproduziert:
die Signatur trägt keinen `ColumnExclusionPort`, kein `ExcludedColumns` im
Startpfad). Dieser Bericht liefert die stärkere Form: **Mutation E** fährt denselben Gedanken am
**laufenden** E2E-Rundlauf und sieht den Neustart-Beleg **real rot**, während
der Live-Reload-Beleg grün bleibt. Damit ist die Nicht-Vakuum-Hälfte
erstklassig belegt — ein abgelegtes Artefakt des verlorenen roten Erstlaufs ist
**nicht** mehr nötig. F-3 trägt als **benannte Grenze**: was fehlt, ist allein
die Reproduzierbarkeit jenes einen roten Laufs (das `:dev`-Image wurde danach
ersetzt); die substanzielle Zusage ist unabhängig bestätigt. Nicht
blockierend; die §3-Klammer „— siehe Bericht" könnte beim Closure-Nachzug auf
den Code-Vergleich (oder auf diesen Bericht) adressiert werden, erzwungen ist
das nicht.

**DoD-Häkchen-Trennung — korrekt.** Per `git show` über die zwei Commits
geprüft: `6be714d` (Implementer) hakt genau die vier Liefer-/Gate-Zeilen
(Nachweis 1–4) plus das Doku-Update ab; `31e698c` (Reviewer) hakt **genau**
die Review-Zeile ab und fügt `docs/reviews/review-slice-075.md` mit dem
DoD-Nachzug hinzu; **kein** Häkchen sitzt auf einer Closure-Zeile. Die sechs
Closure-Punkte bleiben über beide Commits offen.

## 7. Handbuch-Regel und Verweisform-Regel (verkörperte Hard Rules)

- **Kein neues Betreiber-Oberflächen-Element im Diff.** Geprüft gegen die seit
  `slice-077` verkörperte HIGH-Regel
  (`.harness/skills/reviewer.md` §Klassifikation „Neue Betreiber-Oberfläche
  ohne Handbuch-Zug"): der Slice-Diff führt **kein** `CDC_*`-Envar
  (`compose.yaml` nicht im Diff), **keine** neue administrative `cdc.*`-Funktion
  (`tools/schema/**` nicht im Diff), **keine** Horch-Adresse und **keinen**
  Endpunkt. Die Regel hängt am Diff, der die Oberfläche **einführt** — dieser
  tut das nicht. Die vorbestehende Lücke (`cdc.exclude_column`/
  `cdc.include_column` fehlen im Handbuch) ist der 3×-verkörperte Eintrag
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` mit
  Remediation-Träger `slice-077`; ein zweiter Beleg hier wäre doppelte Zählung
  derselben Klasse für eine nicht angefasste Oberfläche. Bestätigt.
- **Verweisform-Regel (`.d-check.yml`) nicht berührt.** Der Diff fügt keinen
  Markdown-Link auf einen Slice-Plan mit festem Lifecycle-Verzeichnis hinzu
  (dieselbe Regel, die der jüngste Architect-Zug verkörpert hat); `make gates`
  läuft mit `docs-check 583 Dateien / 0 Befunde` grün. Dieser Bericht hält
  sich ebenfalls daran (§Selbst-Check).

## 8. Bemerkungen (Verifier-only, beide nicht blockierend)

**B-1 — Der Review-Report zitiert eine Register-Kennung, die es nicht gibt
(INFO).** `review-slice-075.md:236` nennt
`BEO-PGC/handbuch-nachgezogen-bei-neuer-betreiber-oberflaeche`; die reale
Register-Kennung ist
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (das
`nicht-` fehlt). Modul 6 §Lese-Schritt verlangt, vorhandene `BEO-<KUERZEL>`
zu **zitieren**, damit der Zähler nicht zwei Namen trennt — und die
Register-Paarung (c) prüft, dass ein in einer Closure-Notiz genanntes
Verzeichnis existiert. Die Closure muss den **korrekten** Pfad zitieren; die
Kennung im Review-Report bricht beim Nachziehen der Closure die Paarung nicht,
weil die Closure den Pfad selbst setzt — aber sie sollte jetzt schon bekannt
sein. Nicht Teil der Slice-DoD (Review-Artefakt), daher hier nur benannt.

**B-2 — DoD-Wortlaut „simulierter Container-Neustart" vs. realem `docker
restart` (INFO).** §2/§1 des Plans nennen den Beleg einen „simulierten
Container-Neustart"; das Skript fährt einen **realen** `docker restart`
(`run-integration-tests.sh:452`, Block-Kommentar „echten Prozess-Neustart").
Die Repo-eigene Vokabel für denselben Befehl ist „Simulierter
Container-Neustart" (der Black-Box-Rundlauf nennt denselben `docker restart`
so, `run-integration-tests.sh:906`), also ist **kein** DoD-Verstoß
gegeben — die gelieferte Evidenz ist **mindestens so stark**, nicht schwächer.
Nur die Wortgleichheit zwischen DoD und Skript ist nicht hergestellt; wer
schärfen will, kann sie beim Closure-Nachzug angleichen. Erzwungen ist das
nicht.

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/wiring.go` (beide Aufrufer, beide
  Fehlerpfade; die neuen Tests binden die Zusagen — fünf eigene Mutationen)
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/` (neue
  Abfrage + Auswertung; kein Schreibpfad, kein Schema-Objekt; Zweitschlüssel
  und `include_column` einzeln rot)
- geprüft, ohne Befund: `internal/adapters/driving/replication/mapper/`
  (nicht im Fenster — In-Prozess-Erhalt aus `ADR-0065` Festlegung 4
  unangetastet)
- geprüft, ohne Befund: `internal/application/port/outbound/` und die zwei
  Use-Case-Fakes (Fähigkeit am bestehenden Zuschnitt, kein zweiter Port)
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (Position vor
  der Container-Ende-Grenze und vor dem Upgrade-Tausch, drei Zusicherungen
  inkl. Kontrolle gegen ein leeres Row Image; kein geteilter Zählraum)
- geprüft, ohne Befund: `spec/pflichtenheft.md` (Präzisierung, keine
  Spec-Stratum-Verletzung) und `harness/README.md` (Beleg-Baustein korrekt)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (kein neues
  Oberflächen-Element im Diff; die vorbestehende Lücke trägt `slice-077`)
- geprüft, ohne Befund: `harness/image-hash.txt` und `.d-check.yml`
  (Digest reproduzierbar aus demselben Baum; Referenzform-Regel eingehalten)

## Verdikt

**DoD-Konformität:** bestätigt — die sechs implementierungs-/review-/
dokubezogenen DoD-Punkte sind erfüllt und die drei Kern-Zusagen (Startpfad,
Aktivierungs-Zweig, Neustart) in dieser Sitzung **selbst** und über fünf
eigene Mutationen als einzeln tragend nachgewiesen; die fünf Closure-Punkte
sind korrekt offen.

**Entscheidungs-Konformität:** bestätigt — die Umsetzung trägt alle vier
Festlegungen von `ADR-0065`; keine zweite Quelle, kein neues Schema-Objekt,
kein zweiter Schreibpfad, der In-Prozess-Erhalt unangetastet.

**Plan-vs-Code-Diff:** deckungsgleich — keine unerklärte und keine fehlende
Datei; alle fünf §3-Nachzüge tragen.

**Keine DoD-Verletzung, kein blockierender Befund.** Die drei
Review-Findings tragen in meiner eigenen Beurteilung; F-3 verliert seine
Beleg-Lücke durch den E2E-Mutationslauf dieses Berichts. Zwei benannte,
nicht blockierende Bemerkungen (B-1 Register-Kennung im Review-Report, B-2
DoD-Wortlaut) gehen als Kontext an den Planner-Zug; die Closure-Notiz, die
Risiko-Ausgänge, das Beobachtungs-Register und die drei Paarungen bleiben
Planner-Arbeit.
