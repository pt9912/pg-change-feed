# Verifikationsbericht: slice-080 — 2026-09-15

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-080` §2) und die bindenden Entscheidungen
([`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 — die Messung, ihr Träger, ihre Schwelle — und Punkt 4 — „die
Coverage" ohne Subjekt = die Unit-Zahl;
[`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a) — bootstrap-aware Rampe; `AGENTS.md` §3.1, §3.6, §3.9). **Nicht** gegen
den Diff als solchen (Reviewer-Aufgabe — die drei Review-Reports wurden als
Kontext gelesen, **nicht** als Beleg übernommen) und **nicht** gegen realen
Bedarf (Validator, hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan
(`in-progress/`, Stand `HEAD`), die beiden ADRs, die drei Review-Reports, den
Diff `9fa9be3..HEAD` und die berührten Artefakte. **Alle** Zahlen dieses
Berichts stammen aus eigenen, in dieser Sitzung gefahrenen Läufen und eigener
Profil-Deduplizierung; **kein** Beleg des Implementers, des Reviewers oder der
Planner-Closure wurde übernommen. Exit-Codes je in eigenem, ungepiptem Schritt
gelesen (`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung getrennt beauftragt.
Profil- und Rechenartefakte lagen außerhalb des Arbeitsbaums
(`/tmp/v80/…`); die eine Baumänderung, die `make test-store` über den
d-check-Rollout an `tools/schema/plan.yaml` hinterlässt, wurde real
zurückgenommen (`git status --porcelain` leer).

**Gegenstand:** `slice-080`. Der Anteil am Range sind `efc1f23` (der Bau),
`95ccfb4` (Review), `93a5cad` (Fixrunde F-1/F-2/F-3), `c82ef36`
(Delta-Review + DoD-Review-Zeile), `8101656` (zwei Register-Einträge,
Planner-/Implementer-Zug) und `1f0c4ac` (slice-082 angelegt — die Adresse für
den vorbestehenden roten Tier-Lauf, **während dieses Laufs** vom Planner
gesetzt; s. §6).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Zeile (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1a | Die beiden Träger-Läufe erzeugen je ein `-coverprofile`; beide werden gemergt und als **eine** Zahl ausgegeben — das Merge-Verfahren ist benannt | **erfüllt** | Eigene Läufe: `make test-store` legt `store.coverprofile` ab (**8** Dateien, alle `postgresstorage/*.go`, **432** Profilzeilen, **610** Statements); `run-replication-tests.sh measure` legt `replication.coverprofile` ab (**3** Dateien: `postgresack/ack.go`, `receive/receive.go`, `receive/walretention.go`, **264** rohe Zeilen, **356** rohe Statements). Der `measure`-Lauf druckt **eine** Zahl (`DB-Adapter-Coverage: 75.25% …`). Das Merge-Verfahren ist im Skript-Kopf und in `harness/sensors/db-adapter-coverage.md` §Zählbasis benannt (Dedup über die Block-Position, „gedeckt = mindestens ein Vorkommen `count > 0`"). |
| 1b | Die **Herkunft der Zahl** steht dabei (Lauf, Profil, Deduplizierung) — Zählbasis-Regel analog `slice-079` | **erfüllt** | Der Ausgabesatz nennt die gemergten Profile (`Profile gemergt: store,replication`); die Zählbasis steht in `db-coverage.sh` (Kopf + `awk`-Merge) und in `harness/sensors/db-adapter-coverage.md` §Zählbasis Punkt 1–3 (Regel · Nenner · Verfahren). |
| 1c | **Eigene** Auswertung reproduziert **593/788 = 75,25 %** | **erfüllt** | Eigener `awk`-Merge über die Block-Position über **meine** beiden Profile: `covered=593 total=788 pct=75.2538`. Deckungsgleich mit der gedruckten Zahl. |
| 1d | Die zwei Profile sind **disjunkt** (Merge ohne Doppelzählung); Nenner **788 = 610+155+23** unabhängig belegbar | **erfüllt** | Dateimengen-Schnitt via `comm -12` **leer** (8 Store-Dateien, 3 Replication-Dateien; keine gemeinsame Datei). Nenner unabhängig: `postgresstorage` **610** (Store-Profil) + `postgresack` **23** + `replication/receive` **155** (Replication-Profil, dedup: 132 Positionen × 2) = **788**. Im Store-Profil steht **keine** Zeile für `postgresack`/`receive` (Go-Warnung „no packages being tested depend on matches for pattern ./…postgresack"/…replication/receive"); kein `mapper`-Eintrag. |
| 2a | Erster Wert gemessen, Stufe = abgerundet auf volle 5-%-Stufe, Endstufe 80, **keine Senkung** | **erfüllt** | Ist-Stand **75,2538 %** (eigener Lauf) → `floor(75,2538/5)*5 = 75`; `DB_COVERAGE_THRESHOLD=${DB_COVERAGE_THRESHOLD:-75}` in `tools/harness/db-coverage.sh:54`; Endstufe **80** in `harness/sensors/db-adapter-coverage.md` §Kalibrierungs-Bindung. **Keine Senkung:** die DB-Schwelle ist **neu** (kein Vorwert); die Unit-Schwelle (`harness/mk/coverage.mk`, im Gate `COVERAGE_THRESHOLD=65`) ist von diesem Diff **unberührt**. |
| 2b | **Ein** beweglicher Träger (`DB_COVERAGE_THRESHOLD`); die Zahl heißt nie „die Coverage" | **erfüllt** | Repo-weiter `grep` (ohne vendored Baseline, ohne Review-Reports): **genau eine** Zuweisung, `db-coverage.sh:54`. README, Sensor-Doku und Workflow nennen den Namen bzw. verweisen darauf, tragen den Wert nicht. Keine nackte Verwendung „die Coverage" für die DB-Zahl in den neuen Artefakten; die Zahl trägt überall „DB-Adapter-Coverage". |
| 3a | Ein **realer** Lauf zeigt die Zahl, Exit direkt gelesen und ungepiped | **erfüllt** | Eigener `run-replication-tests.sh measure` → **Exit 0**, `DB-Adapter-Coverage: 75.25% (gedeckt 593 von 788 Statements; Profile gemergt: store,replication)` · `db-coverage: OK — DB-Adapter-Coverage 75.25% erfuellt Schwelle 75%`. |
| 3b | `make gates` bleibt **unverändert** grün — kein Container dort, keine zweite Schwelle im Bündel | **erfüllt** | Eigener `make gates` → **Exit 0** (ungepiped): baseline-verify `v6.5.0` OK (54 Dateien) · d-check 664 Dateien / 0 Befunde · commit-traceability OK (5 Commits) · coverage-gate OK 69,70 % (`--build-arg COVERAGE_THRESHOLD=65`) · a-check 0 Befunde. `DB_COVERAGE_THRESHOLD` kommt im Gate-Bündel **nicht** vor; `Makefile`/`harness/mk/**` sind unberührt. |
| 4 | `make gates` grün (§2-Zeile) | **erfüllt** | s. 3b — dieselbe Messung. |
| 5 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | `docs/reviews/review-slice-080.md` (0 HIGH / 0 MEDIUM / 3 LOW / 2 INFO) und `docs/reviews/review-slice-080-fixrunde.md` (0 HIGH / 0 MEDIUM / 0 LOW / 2 INFO) liegen real vor; die §2-Review-Zeile ist auf `[x]`. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 ist reiner Platzhalter (`<…>`), gelesen. Planner-Closure-Arbeit. |
| 7 | Reconciliation-Register fortgeschrieben, **falls** ein Inventur-Fund aufgelöst wird | **entfällt** | `docs/plan/planning/reconciliation.md` existiert real **nicht**; das Repo ist durchgehend GF (`harness/conventions.md` §Modus-Deklaration `*`/`PGC`). Die Entfall-Zeile trägt. |
| 8 | Beobachtungs-Register fortgeschrieben — kein Zähler gesetzt, er folgt aus den Dateien | **erfüllt** | Zwei Verzeichnisse real angelegt und committet (`8101656`): `BEO-PGC/roter-test-ohne-leser/` (evidence `slice-080.md` → abgeleitet **1×**) und `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg/` (evidence `review-slice-079.md`, `review-slice-080.md` → abgeleitet **2×**). Kein Zähler-Feld gesetzt (geprüft). Beide `state.md` tragen `offen`. |
| 9 | Jedes §6-Risiko trägt einen Ausgang (eingetreten / entfallen / weiter offen) | **korrekt offen — mit V-1** | Alle drei §6-Einträge tragen wörtlich `<bei Closure>`; keines vorzeitig geschlossen. Bewertung der Empfehlungen des Implementers in §5 (V-1: „aufgelöst" und „dünn" sind **keine** Elemente der geschlossenen Drei-Menge). |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen, hier nicht zuständig** | Das Repo **hat** Wellen (`welle-20` ist offen, ihre Datei flach unter `docs/plan/planning/`), also trägt die nächste Wellen-Closure die Paarungen (§2 letzte Zeile, `v6.5.0` · `regelwerk/modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 3) — nicht diese Slice-Closure. |

**Ergebnis §1:** Die **zehn** gesetzten Zeilen (1a–5, 8) sind real erfüllt; die
Belege 1a–3b, 4 in dieser Sitzung **selbst** gefahren, die Nenner- und
Profil-Aussage über eine **eigene** Deduplizierung nachgerechnet. Die
Planner-Posten (6, 9, 10) sind korrekt offen; die Entfall-Zeile (7) trägt.
Die Beobachtungen (8) sind real angelegt.

---

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener, ungepipter Schritt)

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | alle inneren Gates grün, Nachweis-Stempel zuletzt (s. 3b) |
| `make test-store` (`DB_COVERAGE_DIR=/tmp/v80/cov`) | **0** | Teilzahl `75.90% (463 von 610)` — nur Store-Profil, keine Schwellen-Prüfung |
| `run-replication-tests.sh measure` | **0** | `DB-Adapter-Coverage: 75.25% (593 von 788)` · `db-coverage: OK … 75%` |
| `make test-replication` (ohne Argument = `both`, **Parität**) | **2** (`make`), Skript **1** | measure-Phase `OK 75.25%`, **danach** Tier-`FAIL`; `make: *** [Makefile:67: test-replication] Fehler 1` |
| **Unveränderter** Runner `9fa9be3:tools/harness/run-replication-tests.sh` | **1** | `--- FAIL: TestWALRetentionThresholdEndToEnd` · `FAIL …/internal/bootstrap`; `42P01` auf `cdc.administration_request` **und** `cdc.process_heartbeat` |
| `make doc-commits RANGE=9fa9be3..HEAD` | **0** | 664 Dateien / 0 Befunde — jeder Commit des Fensters kennungstragend |
| `make doc-immutable RANGE=9fa9be3..HEAD` | **0** | 664 Dateien / 0 Befunde — keine Immutabilitäts-Verletzung |

**Eigene Profil-Auswertung** (Deduplizierung über die Block-Position,
„gedeckt = mindestens ein Vorkommen `count > 0`"):

| Größe | eigener Wert |
|---|---|
| Store-Profil: Dateien / Zeilen / Statements | 8 · 432 · **610** (kein `mapper`) |
| Replication-Profil: Dateien / rohe Zeilen / rohe Statements | 3 · 264 · 356 |
| Replication-Profil: deduplizierte Positionen / Statements | **132** (jede genau **2×**) · **178** |
| gemergte Positionen / Statements | 564 · **788** |
| davon gedeckt | **593 → 75,2538 %** |
| Statements für genau 75,00 % | **591** → **Marge 2 Statements** |

---

## 3. Die Zahlbasis — abschließend

**(a) Reproduziert die eigene Auswertung 593/788 = 75,25 %?** **Ja.**
Mein eigener `awk`-Merge über die Block-Position über meine **beiden**
Profile ergibt `covered=593 total=788 pct=75.2538` — identisch mit der vom
Skript gedruckten Zahl.

**(b) Sind die zwei Profile disjunkt (Merge ohne Doppelzählung)?** **Ja.**
Die Mengen der berührten Quell-Dateien schneiden sich nicht (`comm -12` leer:
8 `postgresstorage`-Dateien gegen `postgresack/ack.go` + 2
`receive`-Dateien). Die beiden Läufe messen verschiedene Testbestände und
partitionieren das Subjekt; ihr Merge ist die **Vereinigung**. Die
Deduplizierung im Replication-Profil ist dennoch nötig und wird von genau
einer Regel getragen: die 132 Positionen erscheinen je **2×** (einmal mit
`count`, einmal mit 0); „nur das erste Vorkommen" würde `receive` auf 0/155
drücken.

**(c) Ist der Nenner 788 = 610+155+23 unabhängig belegbar?** **Ja.** Aus
meinen eigenen Profilen: `postgresstorage` **610** (nur Store-Profil),
`postgresack` **23**, `replication/receive` **155** (Replication-Profil,
dedup). Summe **788** — deckungsgleich mit der Zahl, die `ADR-0071` für die
drei ausgenommenen Pakete nennt. Der Nenner entsteht aus dem Profil, nicht aus
einer gepflegten Konstante.

**Kalibrierung:** 75,2538 % → `floor(75,2538/5)*5 = 75`. Endstufe **80**.
**Ein** beweglicher Träger (`db-coverage.sh:54`). **Keine** Senkung: der
Wert ist neu, die Unit-Schwelle unberührt. **Ein** beweglicher Ort,
nachgeprüft per `grep`.

---

## 4. Die Teilung des Trägers (Fixrunde F-3)

| Prüfung | Verdikt | Beleg |
|---|---|---|
| Der **Mess-Exit** teilt sich mit **keinem fremden** Exit | **trägt** | `measure` allein → **Exit 0** (`OK 75.25%`). Der Tier-Lauf ist der **einzige** rote Bestandteil und liegt in einer **eigenen** Phase; im `both`-Lauf (`make test-replication` ohne Argument) endet der Lauf **rot** (Exit 2), die Messphase davor druckt ihr eigenes `OK`. Der frühere Defekt — `make test-replication` als **ein** Schritt, dessen Tier-`FAIL` den Mess-Exit verschluckt — ist behoben. Kein `|| true` in den neuen Zeilen. |
| Ohne Argument ist die **Parität** gewahrt | **trägt** | `MODE=${1:-both}`; `measure|both` und `tier|both` laufen; der letzte Exit (tier) wird zum Skript-Exit. Eigener Lauf `make test-replication` → Skript-Exit **1**, `make`-Exit **2**, Ausgabe zeigt Mess-`OK` **und** Tier-`FAIL`. Der Unbekannt-Modus-Guard liefert Exit 2 (`case`-Zweig, vor dem Aufbau). |
| Die deklarierte Ausnahme vom `make`-Aufruf im Workflow-Kopf | **ehrliche Deklaration, benannte Grenze — V-2** | Der Kopf nennt die zwei Direkt-Aufrufe `run-replication-tests.sh measure|tier` als „Einzige Ausnahme von den `make`-Aufrufen" samt Grund („ein Make-Target je Phase gibt es nicht"). Die **Substanz** von `AGENTS.md` §3.1 (Docker-only, kein Host-Toolchain) bleibt intakt: das Skript ruft dasselbe gepinnte Toolchain-Image. **Rest:** der Kopf-Satz „Kein Workflow-Schritt enthaelt Inline-Shell-Logik, die eines dieser Ziele umgeht" widerspricht jetzt den zwei Zeilen daneben — und zwei Wege zum selben Skript sind eine Divergenz-Fläche (s. §6). |

---

## 5. Bewertung der §6-Risiko-Empfehlungen (V-1) und der Kalibrierungs-Marge

Der Plan trägt für **jedes** §6-Risiko noch wörtlich `<bei Closure>` — die
Ausgänge sind offen. Der Implementer hat in der Übergabe drei Ausgänge
**vorgeschlagen**; ich bewerte sie gegen die geschlossene Drei-Menge des
Slice-Plans (`eingetreten / entfallen / weiter offen`) und gegen die
Profilbefunde:

- **Risiko 3** — „`mapper` liegt in beiden Messungen; die Zahlen überlappen":
  vorgeschlagen „**nicht eingetreten**". **Trägt und ist formrein
  `entfallen`** („nicht eingetreten" ist die umgangssprachliche Form von
  „entfallen"). Eigener Beleg: das Store-Profil trägt **null** `mapper`-Zeilen;
  das DB-Subjekt ist `./…/postgresstorage` als **exaktes** Paketmuster (kein
  `/...`), `mapper` fällt heraus; der Unit-Gegenstand führt `mapper` und
  schließt die drei Pakete aus. Die Subjekte sind **verschieden**, die zwei
  Zahlen überlappen nicht.
- **Risiko 1** — „ein gemeinsamer Nenner könnte eine Zahl erzeugen, die keinen
  der beiden Gegenstände beschreibt": vorgeschlagen „**aufgelöst**". **„Aufgelöst"
  ist kein Element der Drei-Menge** — formrein ist der Ausgang **entfallen**
  (das Risiko ist nicht eingetreten: der benannte gemeinsame Nenner, die
  disjunkte Partition, beschreibt die **Vereinigung** des Subjekts, mein
  `comm -12`-leerer Schnitt und 788 = 610+178 tragen das). Der Ausgang gehört
  bei Closure auf `entfallen` mit dieser Begründung umgeschrieben, **nicht**
  auf „aufgelöst".
- **Risiko 2** — „die erste Kalibrierung könnte sehr niedrig liegen": 
  vorgeschlagen „**dünn** (593 gegen 591 = 2 Statements Marge)". Auch „dünn"
  ist **kein** Element der Drei-Menge. Der Ist-Stand (75,25 %) ist **nicht**
  „sehr niedrig" → formrein **entfallen**; die **2-Statement-Marge** (mein
  eigener Nachweis: `floor(0,75*788) = 591`, gedeckt 593) ist eine **benannte
  Grenze** des Floor-Mechanismus — sie gehört als solche festgehalten
  (Register oder Sensor-Doku §Grenze), nicht als Ausgangs-Wort. Dass die Marge
  klein ist, ist zulässig (bootstrap-aware: die Stufe folgt der Messung) und
  kein Defekt.

**Ergebnis §5 (V-1):** Zwei der drei vorgeschlagenen Ausgänge sind formunrein.
Weil die Drei-Menge eine **geschlossene** ist (`v6.5.0` ·
`regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst), muss die Closure für Risiko 1 und 2 das Wort **`entfallen`** mit
Begründung setzen — sonst steht zu einem Risiko ein Ausgang, der keiner der
drei ist. Risiko 3 trägt.

---

## 6. Der vorbestehende rote Tier-Lauf (der wichtigste Befund)

| Frage | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|
| (a) **Kann** der Fehler aus diesem Diff stammen? | **Nein** | `git diff --name-only 9fa9be3..HEAD` zeigt **kein** `internal/**`, `cmd/**`, `test/**`; der Gegenstand dieses Slice ist Messung/Runner/Workflow. Der Fehlschlag liegt in `internal/bootstrap` — von diesem Diff unberührt. |
| (b) Ist er wirklich **vorbestehend**? | **Ja** | Eigener Lauf des **unveränderten** Runners `9fa9be3:tools/harness/run-replication-tests.sh` → **Exit 1**, derselbe `--- FAIL: TestWALRetentionThresholdEndToEnd` mit `relation "cdc.administration_request"/"cdc.process_heartbeat" does not exist (SQLSTATE 42P01)`. Mechanik am Artefakt bestätigt: `internal/bootstrap/walretention_endtoend_test.go` fährt `DROP SCHEMA cdc CASCADE` + eigenes `ApplySchema`; `administration_endtoend_test.go:37` dokumentiert selbst, dass dieser Neuaufbau die Tabellen **nicht** trägt; `internal/bootstrap/wiring.go:616` liest `cdc.administration_request` (seit `ADR-0050`). |
| (c) Ist die Konsequenz korrekt benannt? | **Ja, mit Präzisierung** | Die Artefakte benennen: `measure` (die Messung) und `tier` sind **getrennte** Schritte mit je eigenem Exit; der Tier-Schritt ist bis zum Fixture-Nachzug rot (Sensor-Doku §Grenze Punkt 6, Workflow-Kopf, Plan §3). **Präzisierung:** nach der F-3-Teilung ist der **Mess-Schritt selbst grün** (Exit 0, `OK 75.25%`) und **separat** beobachtbar; rot ist der **Tier-Schritt** und damit der **Gesamtlauf** (nicht-blockierend). „Die Zahl selbst ist nicht grün" gilt also **nicht** für den Mess-Schritt, sondern für den Workflow-Lauf als Ganzes. Wer „die DB-Zahl in CI" beurteilt, muss den Mess-Schritt lesen — der Lauf allein färbt grün→rot. Das sollte die Closure präzise festhalten. |

**Adresse des Befunds:** Der Reviewer hatte (F-4) benannt, dass dem roten Lauf
eine **Adresse** fehlt. Sie existiert jetzt: **`slice-082`
(`docs/plan/planning/open/slice-082-replication-fixture-nachzug.md`, Commit
`1f0c4ac`)** — angelegt **während dieses Laufs** vom Planner, mit Bezug auf
`ADR-0050`, `ADR-0071` Punkt 3 und `BEO-PGC/roter-test-ohne-leser`. Der
Slice-Plan verweist bereits in §3/§4 auf `BEO-PGC/roter-test-ohne-leser`.
Damit ist die Adresse **gesetzt** (die Closure braucht sie nur zu zitieren).

---

## 7. Entscheidungs- und Spec-Konformität am Diff

| Prüfung | Verdikt | Beleg |
|---|---|---|
| `ADR-0071` **Punkt 3** umgesetzt (eigene Messung in `test-store`/`test-replication`, `-coverprofile`, eigene Schwelle, Träger `e2e.yml`, nicht `make gates`) | **erfüllt** | s. §1 (1a–3b) |
| `ADR-0071` **Punkt 4** umgesetzt („die Coverage" ohne Subjekt = die Unit-Zahl; die DB-Zahl trägt ihr Subjekt **immer**) | **erfüllt** | s. §1 (2b) |
| `Makefile` / `harness/mk/**` **unberührt** — kein zweites Gate-Target, keine zweite Schwelle im Bündel | **erfüllt** | `git diff --name-only 9fa9be3..HEAD` führt weder `Makefile` noch `harness/mk/**`; `grep` findet `DB_COVERAGE_THRESHOLD` nur in `db-coverage.sh` |
| `internal/**`, `cmd/**`, `test/**` **unberührt** (der Slice misst, er ändert keine Tests) | **erfüllt** | Diff-Namensliste führt keinen dieser Pfade |
| `spec/**` **unberührt** | **erfüllt** | dito |
| Die **ADRs** unberührt | **erfüllt** | dito (Rang-4-Dateien nicht im Diff) |
| `ADR-0071` §Fitness-Function, Zeile 2 („neues Target (geplant)") | **erfüllt in der Substanz, Wortlaut-Abweichung — V-3** | Die operative Klausel (§Entscheidung Punkt 3: Messung „in `test-store`/`test-replication`") ist wörtlich umgesetzt; die Fitness-Function-Zelle nennt ein „neues Target", das der Slice bewusst **nicht** anlegt (`Makefile` unberührt). Die Bindung existiert (`harness/sensors/db-adapter-coverage.md`). Buchhaltung, kein Defekt — die Closure sollte die Zelle nachziehen oder die Abweichung benennen. |

**Diff-Umfang (18 Dateien, `9fa9be3..HEAD`):** `.github/workflows/e2e.yml`,
der §2/§3-Nachzug im Slice-Plan, `harness/README.md` (zwei Werkzeug-Zeilen),
`harness/sensors/coverage-gate.md` (§Grenze), `harness/sensors/db-adapter-coverage.md`
(neu), `tools/harness/db-coverage.sh` (neu), `tools/harness/run-store-tests.sh`,
`tools/harness/run-replication-tests.sh`, die drei Review-Reports, die zwei
Register-Verzeichnisse (7 Dateien) und `open/slice-082-…md`.

---

## 8. Findings (V-Nummern)

### V-1 — Formunreine §6-Risiko-Ausgänge „aufgelöst" / „dünn"

- `kategorie`: LOW
- `quelle`: `v6.5.0` · `regelwerk/modul-05-planning-harness.md` §Offene
  Risiken werden bei Closure aufgelöst (geschlossene Drei-Menge)
- `pfad`: `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md`
  §6 (noch `<bei Closure>`) — Übergabe-Empfehlungen „aufgelöst" (Risiko 1)
  und „dünn" (Risiko 2)
- `befund`: Von den drei vorgeschlagenen Ausgängen sind zwei **keine**
  Elemente der geschlossenen Menge `{eingetreten, entfallen, weiter offen}`.
  Formrein ist für Risiko 1 und 2 **`entfallen`** (mit Begründung); Risiko 3
  („nicht eingetreten") trägt als `entfallen`. Die 2-Statement-Marge ist eine
  **benannte Grenze**, kein Ausgangs-Wort.
- `verifizierbar`: ja — mein `floor(0,75*788) = 591` gegen gedeckt 593; der
  leere Datei-Schnitt für Risiko 1 und 3
- `klasse`: „Risiko-Ausgang außerhalb der geschlossenen Drei-Menge"

### V-2 — Zwei Wege zum selben Skript; widersprüchlicher Kopf-Satz

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.1 (Docker-only); `ADR-0051` (Workflow ruft
  bestehende Targets); `v6.5.0` · `regelwerk/grundlagen-harness-dateien.md`
  §Was ein Kommentar trägt (Ist-Zustand)
- `pfad`: `.github/workflows/e2e.yml:118-122`, Kopf `:13-28`
- `befund`: Zwei der Schritte rufen `bash tools/harness/run-replication-tests.sh
  measure|tier` direkt (make-Aufruf umgangen), während der Store-Schritt
  `make test-store` nutzt. Die Ausnahme ist **deklariert** und die
  Docker-only-Substanz intakt — aber der Kopf-Satz „Kein Workflow-Schritt
  enthaelt Inline-Shell-Logik, die eines dieser Ziele umgeht" steht daneben in
  Widerspruch, und eine spätere Rezeptur-Änderung an `make test-replication`
  erreicht die beiden CI-Schritte nicht. Zwei Wege, zwei Divergenz-Punkte.
- `verifizierbar`: teilweise — `grep` zeigt kein Make-Target je Phase
- `klasse`: „`make`-only-Invariante deklariert gelockert (benannte Grenze)"

### V-3 — `ADR-0071`s Fitness-Function-Zelle „neues Target" wörtlich nicht angelegt

- `kategorie`: INFO
- `quelle`: `ADR-0071` §Fitness Function (Zeile 2, Spalte Make-Target),
  §Konsequenzen („eigenes Target")
- `pfad`: `ADR-0071` §Fitness Function vs. Slice-Plan §3 („kein neues
  Gate-Target") und `.github/workflows/e2e.yml:115-122`
- `befund`: Die Messung hängt an den **bestehenden** Targets
  `make test-store`/`make test-replication` plus dem neuen Hilfsskript
  `db-coverage.sh` — kein **neues** Target. Vereinbar mit der operativen
  Klausel (Punkt 3), nicht mit der Fitness-Function-Zelle. Buchhaltung.
- `verifizierbar`: ja — `grep db-coverage Makefile harness/mk/` findet nichts
- `klasse`: „ADR-Fitness-Function-Zeile wörtlich nicht eingelöst"

**Negativbefunde** (geprüft, ohne Befund): `db-coverage.sh` — Merge/Dedup/
Zählbasis reproduzieren **593/788 = 75,25 %**, der Nenner ist profil-geboren,
alle Ausgänge (0/1/2, Teilzahl) verhalten sich wie dokumentiert; die
Schwelle steht an **einem** Ort. `run-store-tests.sh`/`run-replication-tests.sh`
— `postgresstorage` läuft **genau einmal**, der Messlauf steht vor dem
Tier-Lauf, kein stiller Paket-Ausschluss. Die zwei neuen Register-Einträge
— Pfad-Kennungen korrekt abgeleitet, kein Zähler-Feld, `state.md` `offen`.
`harness/sensors/coverage-gate.md` §Grenze — die `mapper`-Aussage ist konsistent
mit den Profilen. Kein host-lokaler Pfad, keine Chronik, kein Vorlagen-Rest in
den neuen Artefakten.

---

## 9. Verdikt

**DoD-Konformität:** Die gesetzten §2-Zeilen sind **erfüllt** (mit eigenen
Läufen belegt); die Planner-Posten sind **korrekt offen**; die Entfall-Zeile
trägt. **Kein Häkchen fällt.**

**Entscheidungs-Konformität:** `ADR-0071` Punkte 3 und 4 sind umgesetzt;
`Makefile`/`harness/mk/**`, `internal/**`, `cmd/**`, `test/**`, `spec/**` und
die ADRs sind unberührt. **Keine** DoD- oder Entscheidungs-Verletzung.
**V-1/V-2/V-3 sind nicht blockierend** (V-1 ist eine Closure-Formpflicht,
V-2/V-3 sind benannte Grenzen/Buchhaltung).

**Übergabe an den Planner (Closure-Obliegenheiten — hier benannt, nicht
ausgeführt):**

1. **§6-Risiko-Ausgänge** setzen — für Risiko 1 und 2 formrein **`entfallen`**
   (mit Begründung), für Risiko 3 `entfallen`; die **2-Statement-Marge** als
   benannte Grenze festhalten (V-1).
2. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag schreiben; die drei
   **Finding-Klassen** der zwei Reviews in den Zähler geben.
3. **`git mv` nach `done/`** (reiner Move-Commit; Inhalt davor).
4. **Drei Paarungen** — Träger ist die **nächste Wellen-Closure** (`welle-20`
   ist offen); die Slice-Closure trägt sie **nicht**.
5. **Adresse** des roten Tier-Laufs zitieren: **`slice-082`** (liegt real in
   `open/`, Commit `1f0c4ac`) — gesetzt.
6. Die zwei **Register-Einträge** (`roter-test-ohne-leser` 1×,
   `mechanismus-erklaerung-ohne-werkzeugbeleg` 2×) sind angelegt; §2-Zeile
   nachziehen.
7. **V-2** (zwei Wege zum Skript / Kopf-Satz) und **V-3**
   (Fitness-Function-Zelle) als benannte Grenzen in §7 festhalten oder
   adressieren.

Dieser Bericht ist die Verifikation gegen DoD/Entscheidungen (Modul 11); er
ersetzt **nicht** das Review gegen Plan/ADR (Modul 10) und **nicht** die
Validierung gegen realen Bedarf (Modul 8/11).
