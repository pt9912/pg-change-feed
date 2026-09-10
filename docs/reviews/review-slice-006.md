# Review-Report: slice-006 Implementer-Diff — 2026-09-10

**Review-Art:** Diff-Review (Implementer-Range `c286339..f79af17`; im Range zusätzlich
`3ea9244` — Allowlist-Commit „Vorgabe pt9912", Bewertung siehe F-2) — *wogegen*:
Slice-Plan §1/§2/§3/§6 (Plan-Treue, Bewertung der gemeldeten Entscheidungen),
ADR-Bezüge ([`ADR-0023`](../plan/adr/README.md), 0026, 0029, 0043, 0044), Hard Rules
(`AGENTS.md` §3.1 Docker-only, §3.3 mv/content-Trennung, §3.7 Kommentar-Klassen),
Traceability (LH-*/ADR-* je Commit, keine Struktur-IDs, keine superseded-Referenzen),
d-check/a-check-Fixes des Laufs, Test-Qualität gegen den MVP-Schnitt-Wortlaut
([`spec/lastenheft.md` §1](../../spec/lastenheft.md)), Implementer-Risiken (a)–(c),
Kennungs-Prüfung der Plan-Kopf-Zeile. Keine DoD-Prüfung — das ist der Verifier
(Modul 11).

**Gegenstand:** `3950e77` (schema.yaml-Erstlieferung) · `779bc64` (Ruhe-Marker) ·
`2d07c40` (compose + Rollout-Erstversatz) · `8d23352` (MVP-Integrationstest,
`.a-check.yml`-Glob) · `f79af17` (`.PHONY`-Fix) · `3ea9244` (`.claude/settings.json`)
— der Range-Start `c286339` ist der reine `git mv` (next → in-progress, §3.3 sauber).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, vier repo-spezifische
HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst:
`docs/reviews/review-report.template.md` (Form wie `review-slice-005.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Diff `git diff c286339..f79af17` (13 Dateien, +754/−6: `tools/schema/schema.yaml`
  + `plan.yaml`/`down.sql`/`nacharbeit-operation-check.sql`/`compose-init/`,
  `compose.yaml`, Makefile, `tools/harness/run-integration-tests.sh`,
  `test/integration/mvp_test.go`, `.a-check.yml`, `.claude/settings.json`,
  `harness/README.md`-Werkzeuge-Zeile, Roadmap-Ruhe-Marker)
- `docs/plan/planning/in-progress/slice-006-integrationstest-umgebung.md` (§1–§8;
  am Range-Base `c286339` und am Range-Head **unverändert** — der Plan trägt die
  Plan-Nachzüge aus `fc5ce9f`/`31e1b41` bereits)
- `docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md` (Erstversatz, §5
  Testloader-Grenze, Re-Evaluierungs-Trigger) · [`ADR-0044`](../plan/adr/README.md) (Lauf-Beleg) ·
  [`ADR-0026`](../plan/adr/README.md) (Composition Root) · [`ADR-0029`](../plan/adr/README.md) (Regel 1/3/6) ·
  [`ADR-0023`](../plan/adr/README.md)/`SPEC-008` (Fehlerklassen), `docs/plan/planning/welle-2.md`,
  `docs/plan/planning/in-progress/roadmap.md`
- `spec/lastenheft.md` (MVP-Schnitt §1, [`LH-FA-CAP-001`](../../spec/lastenheft.md)…004/008,
  [`LH-FA-REA-002`](../../spec/lastenheft.md)/004.a/005, [`LH-FA-CFG-001`](../../spec/lastenheft.md),
  [`LH-QA-POR-003`](../../spec/lastenheft.md), [`LH-QA-OPS-001`](../../spec/lastenheft.md)/002),
  `spec/pflichtenheft.md` ([`LH-QA-REL-001.a`](../../spec/pflichtenheft.md), [`LH-FA-REA-004.a`](../../spec/pflichtenheft.md),
  Kennungs-Prüfung `PH-*` — F-11), `internal/adapters/driven/postgresstorage/
  schema.sql` (Überführungs-Quelle) und `store.go` (Schreibpfad)
- Commit-Messagen der Range; `AGENTS.md` §3/§5; `harness/conventions.md` (MR-000);
  `.a-check.yml` (composition_root); `.claude/settings.json` samt
  `.claude/hooks/*` (Existenz geprüft)

**Gate- und Probe-Läufe (am Range-Head `f79af17`, alles über `make`/gepinnte
Digests):** `baseline-verify` OK (54 Dateien) · `d-check` 95 Dateien/0 Befunde ·
`a-check` 0 Befunde (die Glob-Erweiterung deckt `test/integration/**`) ·
`make schema-validate` grün (5 Tabellen, 21 Spalten, 5 Constraints, 1 Index) ·
`make test` grün (`.PHONY`-Fix belegt: das Target fährt das Rezept,
`test/integration` kompiliert und skippt ohne DSN) · **`make test-integration`
grün am Range-Head** (Compose frisch → Rollout inkl. Nacharbeit-Schritt →
MVP-Test gegen die Compose-DB; Re-Lauf lässt `plan.yaml`/`down.sql`
byte-identisch — der Rollout-Beleg ist deterministisch) · eigene
`schema generate`-Probe (F-1 unten). `make image` nicht ausgeführt —
keine Build-Kontext-Datei des Ranges berührt den Dockerfile-Pfad.

---

## Findings

### F-1 — Schema-Drift DDL ↔ neutrales Modell: `committed_at` verliert `NOT NULL` (eigene Probe belegt)

- `kategorie`: HIGH
- `quelle`: [`ADR-0043`](../plan/adr/README.md) Folgepflicht („Vor dieser Überführung bleibt
  die DDL die Quelle; Abweichung zwischen DDL und YAML wäre zwei Quellen für
  dieselbe Schema-Form **und ein Befund**") · Skill-HIGH-Klasse
  „Zwei-Quellen-Drift"
- `pfad`: `tools/schema/schema.yaml:95` (`committed_at: { type: datetime,
  timezone: true, default: current_timestamp }` — ohne `required: true`) gegen
  `internal/adapters/driven/postgresstorage/schema.sql` (`committed_at
  timestamptz NOT NULL DEFAULT now()`) und den committeten Rollout-Beleg
  `tools/schema/plan.yaml` (Statement: `"committed_at" TIMESTAMP WITH TIME
  ZONE DEFAULT CURRENT_TIMESTAMP` — ohne NOT NULL)
- `befund`: Die Überführung behauptet am Commit `3950e77`, die generated DDL
  trage „dieselbe Form wie die handgeschriebene DDL (Typen, Keys, Checks,
  UNIQUEs, Index)". Eigene `schema generate`-Probe am Range-Head: die erzeugte
  DDL trägt `committed_at` **ohne** `NOT NULL` — die handgeschriebene DDL
  trägt die Spalte NOT NULL. Die Abweichung ist genau der Fall, den
  [`ADR-0043`](../plan/adr) als Befund deklariert, und sie ist im committeten
  Beleg (`plan.yaml`) selbst lesbar. Verhaltensnah, aber nicht verhaltens-
  gleich: der Store-Schreibpfad schreibt `committed_at` nie
  (`store.go:95-99` — nur drei Parameter, die Spalte kommt aus dem DEFAULT),
  unter dem Rollout-Schema wäre eine externe NULL-Einfügung aber
  möglich, die die DDL-Form ausschließt. Auch die Constraint-Namen
  weichen ab (DDL unnamed UNIQUE → Katalog-Autoname; YAML benannt) —
  derselbe Drift-Typ, sekundär.
- `verifizierbar`: ja — `make schema-validate` grün trotz Abweichung (kein
  Sensor deckt die DDL↔YAML-Deckung); Probe
  `schema generate --target postgresql | grep committed_at`; `grep -n
  committed_at` in DDL, YAML und `plan.yaml`
- `klasse`: Zwei-Quellen-Drift der Schema-Form (Überführungs-Abweichung,
  [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Befund-Fall)

### F-2 — Zwei Commits des Ranges ohne `LH-*`-/`ADR-*`-Kennung (2./3. Auftreten — die Klasse erreicht 3×)

- `kategorie`: HIGH
- `quelle`: `harness/README.md` §Traceability rules („PRs/Commits **müssen**
  mindestens eine `LH-*` oder `ADR-*`-ID nennen") · Skill-HIGH-Klasse
  „Traceability-/ID-Schema-Verstoß" · `AGENTS.md` §5
- `pfad`: Commit `3ea9244` (`.claude/settings.json` — Body nennt nur „Vorgabe
  pt9912", keine `LH-*`/`ADR-*`-Kennung) · Commit `779bc64` (Ruhe-Marker —
  Body nennt `slice-006` und den Regelwerk-Pfad, keine `LH-*`/`ADR-*`-Kennung)
- `befund`: Vier der sechs Commits des Ranges tragen je mindestens eine
  Vertrags-Kennung; diese beiden nicht. Der Marker-Commit berührt die
  Roadmap-Regeln (Modul 6), der Allowlist-Commit die Lauf-Bindung — beide
  berühren Verträge ohne die Kennung des berührten Vertrags zu nennen. Damit
  steht die direkte Form (review-slice-005 F-1, 1. Auftreten) beim **dritten**
  Auftreten — der Steering-Loop-Zähler ist erreicht (1× notieren · 2×
  Symptom · 3× Lücke); die Verkörperung (Gate/Conventions-Nachzug) ist
  über den **Architect** fällig, nicht still wiederholbar. Die Messagen
  bleiben historisch; Korrektur wirkt nur vorwärts.
- `verifizierbar`: ja — `git log --format=%B c286339..f79af17` gegen
  `harness/README.md` §Traceability rules
- `klasse`: Commit ohne Vertrags-Kennung (3. Auftreten — Sequenz-Pflicht fällig)

### F-3 — Plan-Erweiterung ohne Plan-Nachzug (6. Auftreten — laufende Sequenz)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3 · Modul 5 („Wer später mitnimmt …, hat den Plan
  **geändert**, nicht nur ergänzt") · laufende Konflikt-Sequenz (Modul 8 —
  seit review-slice-003 F-1; 5. Auftreten review-slice-005 F-2)
- `pfad`: `docs/plan/planning/in-progress/slice-006-integrationstest-umgebung.md`
  §3 (vier Zeilen; keine Zeile für `.a-check.yml` oder `harness/README.md`)
  gegen `.a-check.yml` (composition_root um `test/integration/**` erweitert)
  und `harness/README.md` (Werkzeuge-Zeile `test-integration`), beide in
  `8d23352`
- `befund`: Zwei gelieferte Dateien des Ranges fehlen in §3; ihre Begründung
  lebt nur im Commit-Text („ohne Eintrag fiel der Pfad aus der Abdeckung und
  der Lauf meldete wrong-direction"). Beide Berührungen sind öffentliche
  Verträge (Maschinenform der Architekturprüfung; Werkzeuge-Tabelle). Die
  Fixes selbst sind sachlich richtig (siehe Design-Entscheidungen), aber die
  Klasse steht beim sechsten Auftreten — die Konflikt-Sequenz läuft weiter;
  der Nachzug geht als Übergabe-Artefakt an den **Planner**.
- `verifizierbar`: ja — Datei-Menge gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (6. Auftreten — Sequenz läuft)

### F-4 — Nacharbeit-Brücke `chk_change_operation`: Rückbau-/Pin-Hebungs-Prüfung ohne terminierten Ausgangs-Träger

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0043`](../plan/adr/README.md) Re-Evaluierungs-Trigger („es gibt
  *keine* Ausweichform (… berichtete manuelle Nacharbeit)" — die Ausweichform
  existiert hier, also feuert der Werkzeugwechsel-Trigger nicht) · Modul 5
  („offene Risiken werden bei Closure aufgelöst") · Implementer-Risiko (a)
- `pfad`: `tools/schema/nacharbeit-operation-check.sql` ·
  `tools/schema/schema.yaml:14-33` (Grenze benannt) · `Makefile` (Kommentar
  über die E012-Grenze und die Runner-Räumung) · gegen den Slice-Plan §6
  (zwei Risiken, keine Nacharbeit-Brücke) und `docs/plan/planning/open/`
  (leer — kein Folge-Slice)
- `befund`: Die Ausweichform ist sauber *berichtete* Nacharbeit (Konvention,
  Rollout-Kette, DDL-Gegenprobe — siehe Design-Bewertung), aber die
  Rücknahme der Nacharbeit ist nirgends terminiert: die Retirement-Prüfung
  (Pin-Hebung bei neuem d-migrate-Release, Rückbau von
  `nacharbeit-operation-check.sql` und der Ausweich-Kette im
  `schema-rollout`-Target) steht nur in Datei-Kommentaren und Commit-Text.
  Kein Plan-§6-Risiko trägt sie, keine Folge-Slice-ID nimmt sie an, und kein
  Sensor beobachtet neue d-migrate-Releases (`make image-stale` deckt nur
  die Base-Images des OCI-Builds, braucht Netz, deckt `D_MIGRATE_IMAGE`
  nicht). Der Auftraggeber meldet, d-migrate löse die Konvergenz *heute
  schon* — genau deshalb ist die Brücke ohne Anker still ausreifbar.
  Verankerung: §6-Risiko im Plan mit Ausgang „weiter offen" →
  Beobachtungs-Register (`BEO-PGC/<slug>`, Beleg je Slice-Closure), oder
  Folge-Slice-ID in `open/` — beides Planner-Arbeit vor der Closure.
- `verifizierbar`: ja — `grep -n nacharbeit` im Plan/§7 und Register gegen
  die Datei-Menge
- `klasse`: Befristete Grenze ohne Ausgangs-Träger (Rückbau ungeterminiert)

### F-5 — Beleg-Kette des Pflicht-Reports endet vor der Nacharbeit

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0043`](../plan/adr/README.md) Regel 3 („Der Report ist Pflicht
  und ist der Beleg des Rollouts") · Maintainability
- `pfad`: `tools/schema/plan.yaml` (`postUpVerified: true`,
  `postUpFingerprint: eafe26…`, 5 CreateTable-Operationen) gegen
  `tools/schema/nacharbeit-operation-check.sql` (Constraint entsteht
  *nach* dem Report-Lauf, im zweiten docker run desselben Targets)
- `befund`: Der Rollout-Beleg beschreibt den DB-Stand, den die Nacharbeit
  anschließend ändert: zur Beleg-Zeit existiert `chk_change_operation`
  nicht (eigener Lauf zeigt das NOTICE „constraint … does not exist,
  skipping"); der Report und das `down.sql`-Artefakt kennen es nicht.
  Ein Leser des Belegs sieht einen vollständigen, verifizierten Stand, der
  nicht der Endzustand ist — die Brücke ist nur über die Makefile-Kette
  rekonstruierbar. Für den MVP-Rollout tragbar (Constraint-Aufnahme ist
  Teil der kommittierten Kette), aber die Beleg-Lücke gehört zur selben
  Grenze wie F-4 und übersteht die Brücke nicht still.
- `verifizierbar`: ja — Reihenfolge im `schema-rollout`-Rezept gegen
  `postUpVerified`; `down.sql` ohne den Constraint (deckt trotzdem
  vollständig, da DROP TABLE die Constraints mitnimmt)
- `klasse`: Beleg beschreibt abwesenden Endzustand (Nacharbeit nach Report)

### F-6 — Verdrahtungs-/Aktivierungs-Grenze des MVP-Claims ohne Ausgangs-Träger

- `kategorie`: MEDIUM
- `quelle`: [`LH-QA-POR-003`](../../spec/lastenheft.md) (MVP-Abnahme „vollständige
  Testumgebung") · Modul 5 §6 (jedes Risiko braucht einen Ausgang) ·
  Implementer-Risiko (b)
- `pfad`: `compose.yaml:44-53` (Feed-Container: `command:
  ["--version"]`, Binary `0.1.0-bootstrap`, unverdrahtet) ·
  `test/integration/mvp_test.go:70-88` (Aktivierung per direktem SQL:
  `CREATE PUBLICATION` + Bindungs-Zeilen in den CDC-Referenztabellen) ·
  gegen Slice-Plan §6 (nimmt die Grenze nicht an) und `open/` (leer)
- `befund`: Der MVP-Integrationstest belegt den MVP-Schnitt auf
  Adapter-/Service-Ebene (reale Treiber, reale Store-Adapter, echtes
  Capture-Persist-before-ACK) — die Produkt-Laufzeit trägt den Beleg nicht:
  der Feed-Container ist Image-Vertrag-Smoke (`--version`, Exit 0), und das
  „CDC aktivieren" des Ablaufs läuft am Test-Rand per SQL, nicht durch eine
  Produkt-Fähigkeit. Beide Grenzen sind in `compose.yaml`-Kommentar und
  Commit-Message offen benannt — aber sie tragen **keinen Ausgang**: kein
  §6-Risiko im Plan, kein Bootstrap-/Verdrahtungs-Slice in `open/`. Und
  `welle-2.md` §1 hängt den M1-Claim („der MVP-Schnitt ist belegt; der
  Meilenstein M1 kann erreicht werden") an genau diesen Test, ohne die
  Grenze zu nennen — der Meilenstein-Beleg ist ein Adapter-Ebenen-Beleg.
  Vor der Closure benennen: Ausgang im Plan §6 (Kennung des Slices, der die
  Verdrahtung annimmt) oder Welle-§1-Formulierung schärfen.
- `verifizierbar`: ja — `grep -n version cmd/pg-change-feed/main.go` (Stub),
  Test-Aktivierungs-SQL gegen [`LH-FA-CFG-001`](../../spec/lastenheft.md)-Oberfläche
  (nicht MVP, kein Träger); `ls open/` leer
- `klasse`: Grenze benannt ohne Ausgangs-Träger (slice-005-F-4-Muster)

### F-7 — Datei-Abschluss: fünf neue Dateien ohne Zeilenumbruch (Klasse bleibt MEDIUM)

- `kategorie`: MEDIUM
- `quelle`: Maintainability — Klasse seit review-slice-005 F-6 auf MEDIUM-Stufe
  (drittes Auftreten dort; „Wiederholung eines Musters, das schon zweimal
  LOW war")
- `pfad`: `compose.yaml` · `tools/schema/schema.yaml` ·
  `tools/schema/nacharbeit-operation-check.sql` ·
  `tools/schema/compose-init/01-cdc-schema.sql` ·
  `tools/harness/run-integration-tests.sh` — je letzte Zeile ohne `\n`
  (`.a-check.yml` ebenfalls, Vorbestand); `mvp_test.go`, `down.sql`,
  `plan.yaml` schließen sauber
- `befund`: Fünftes/Sechstes Auftreten derselben Klasse — diesmal in fünf von
  acht neuen Dateien eines Commits; die Klasse ist ab hier MEDIUM-Träger
  im Steering-Loop.
- `verifizierbar`: ja — `tail -c1 <Datei> | od -c`
- `klasse`: Datei-Abschluss (Klasse MEDIUM, weiter auflaufend)

### F-8 — Feed-Container-Smoke ohne Beleg-Abgleich: `:dev`-Tag kann alt sein

- `kategorie`: LOW
- `quelle`: [`ADR-0044`](../plan/adr/README.md) (Digest ist Lauf-Beleg; innerhalb
  desselben Builder-Laufs ist ein Digest-Vergleich gültiger Befund) ·
  Maintainability
- `pfad`: `compose.yaml:49` (`image: ghcr.io/pt9912/pg-change-feed:dev`) ·
  `tools/harness/run-integration-tests.sh:44-49` (Smoke prüft nur Exit 0)
- `befund`: Der Runner prüft, dass der Feed-Container Exit 0 endet, aber
  nicht, dass `:dev` das geladene Image des letzten `make image`-Laufs ist
  (Abgleich gegen `harness/image-hash.txt` fehlt). Ein veralteter Tag läuft
  als „Compose-Anteil des Lauf-Belegs" grün durch den Smoke.
- `verifizierbar`: ja — Smoke mit veraltetem `:dev` bleibt grün
- `klasse`: Beleg-Anteil ohne Deckungsprüfung

### F-9 — Plan §1 nennt „Test-Consumer"; Compose trägt keinen Consumer-Dienst

- `kategorie`: LOW
- `quelle`: Slice-Plan §1 (Ziel: „PostgreSQL + PG Change Feed + Test-Consumer")
- `pfad`: `compose.yaml` (zwei Services: `postgres`, `pg-change-feed`) gegen
  den Plan-§1-Wortlaut
- `befund`: Das „Changes lesen" des MVP-Ablaufs läuft in-process über den
  Store-Port im Toolchain-Container, nicht als Compose-Service. Der Ziel-
  Satz listet drei Umgebungsteile, die Lieferung trägt zwei; die Lesart
  „Test-Consumer = der lesende Test" ist vertretbar, steht aber nirgends.
- `verifizierbar`: ja — `grep -n "Test-Consumer"` im Plan gegen
  `compose.yaml`-Service-Menge
- `klasse`: Ziel-Wortlaut vs. Lieferung (Lesart unbenannt)

### F-10 — CDC-Aktivierung per direktem SQL: MVP-marker-konsistent, Grenze steht nur im Commit-Text

- `kategorie`: INFO
- `quelle`: [`LH-FA-CFG-001`](../../spec/lastenheft.md) (kein `MVP: ja`-Marker;
  Idempotenz-/Negative-Pfade der Aktivierung sind kein MVP-Träger) ·
  `welle-2.md` §6 (SQL-/CLI-Adapter „nicht MVP")
- `pfad`: `test/integration/mvp_test.go:113-125` (Aktivierung: Publication +
  Bindungs-Zeilen per SQL)
- `befund`: Die Aktivierung umgeht keine MVP-Fähigkeit — [`LH-FA-CFG-001`](../../spec/lastenheft.md)
  trägt keinen MVP-Marker, und welle-2 deklariert die SQL-/CLI-Adapter als
  spätere Wellen. Der MVP-Schnitt-Wortlaut („CDC aktivieren") ist als
  Umgebungs-Schritt erfüllt. Hinweis ohne erwartete Aktion am Diff; die
  Ausgangs-Frage ist F-6.
- `verifizierbar`: ja — MVP-Marker-Lage im Lastenheft §1 gegen die
  Aktivierungs-Statements
- `klasse`: Aktivierungs-Grenze am MVP-Rand (marker-konsistent)

### F-11 — `PH-DEP-002`/`PH-TST-001` existieren nicht (Vorbestand, im Range nicht berührt)

- `kategorie`: INFO
- `quelle`: Skill-HIGH-Klasse „Traceability-/ID-Schema-Verstoß" (Erstauftreten
  review-slice-002 F-2; Vorbestand-Fassung review-slice-005 F-16)
- `pfad`: `docs/plan/planning/in-progress/slice-006-integrationstest-umgebung.md:12`
  gegen `spec/pflichtenheft.md` und `spec/lastenheft.md` (kein Treffer für
  `PH-*` in beiden Straten)
- `befund`: Beide Kennungen existieren in keinem Stratum — der Plan-Verweis
  löst nicht auf. Der Implementer-Diff berührt die Plan-Datei nicht (daher
  kein „berührt und unkorrigiert" wie in slice-005); der Befund ist
  unverändert Planner-Sache: Plan-Korrektur fällig, und kein Sensor deckt
  das `PH-*`-Präfix (d-check prüft Referenzen, nicht ID-Existenz gegen das
  ID-Schema — die Lücke bleibt benannt).
- `verifizierbar`: ja — `grep -rn "PH-" spec/` ohne Treffer
- `klasse`: Undeklariertes ID-Präfix (Vorbestand, unberührt im Range)

### F-12 — `welle-2.md` §3: Closure-Trigger trägt noch die Vorlagen-Platzhalter

- `kategorie`: INFO
- `quelle`: Modul 6 (Welle schließt durch beobachtbare Closure-Kriterien) ·
  Vorbestand, entdeckt beim Lesen des Kontexts
- `pfad`: `docs/plan/planning/welle-2.md` §3 („<z.B. Alle Slices done.>" u.
  a. — unfilled) · §2 Bullet endet mit Vorlagen-Glitch „**bereits
  eingetreten**.>"
- `befund`: Die Welle, die slice-006 einsammelt, trägt keinen benannten
  Closure-Trigger — für die anstehende welle-2-Closure (M1-Claim) ist das
  die größte offene Form-Lücke im Planning-Stratum. Vorbestand (Anlage
  `d5010e5`), nicht Teil des Implementer-Diffs; Planner-Sache vor der
  Welle-Closure.
- `verifizierbar`: ja — Lese der Welle-Datei §2/§3 gegen die Vorlage
- `klasse`: Vorlagen-Platzhalter in eröffneter Welle (Vorbestand)

## Design-Entscheidungen des Implementers — Bewertung (Prüf-Fragen)

- **(a) Store-Seite ohne Consumer-/Capture-State-Tabellen in `schema.yaml`:**
  *Konsistent mit [`ADR-0043`](../plan/adr/README.md), kein Befund.* Die ADR nennt die
  Consumer-/Betriebs-Tabellen ausdrücklich als „trägt die DDL noch nicht"
  (Kontext) und legt den Ausbau auf einen eigenen Slice (§4); die
  Überführung ist ausdrücklich die DDL „Stand slice-004" als Quelle, und der
  YAML-Kommentar trägt die Grenze mit Spec-Verweis nach — die Grenze ist
  dieselbe wie in der DDL, nicht eine neu erfundene. Kein Stratum-Verstoß,
  kein stiller Umfang. (Die Abweichung im Umfang ist F-1, nicht (a).)
- **(b) Feed-Container als Image-Vertrag-Smoke:** *Sachlich richtig getragen,
  aber ohne Ausgangs-Träger* — **F-6**. Der Smoke selbst ist sauber gebaut
  (Exit-Code-Prüfung, Abgrenzung im Kommentar, kein CDC-Runtime-Claim), und
  die Aktivierungs-Lesart ist MVP-marker-konsistent (F-10). Was fehlt, ist
  der Plan-Ausgang, nicht der Code.
- **(c) Rollout-Erstversatz:** *Konform.* [`ADR-0043`](../plan/adr/README.md) §4 nennt
  den Compose-Rollout als Ersteinsatz; Pflicht-Report (`tools/schema/plan.yaml`,
  `status ok`, 5 Operationen, `postUpVerified`) und Rollback-Artefakt
  (`tools/schema/down.sql`) sind committet und deterministisch (Re-Lauf im
  Review: byte-identisch); `make schema-validate` als Vorlauf hängt am
  Rezept (`schema-rollout: schema-validate`). `make schema-rollout` hängt
  korrekt an keinem Gate (ADR-§Fitness). Grenzen: F-4/F-5.
- **(d) Nacharbeit `chk_change_operation` (d-migrate 1.2.0, `raw-sql-text-drift`):**
  *Ausweichform korrekt nach dem Re-Evaluierungs-Trigger gewählt.* Der
  Trigger verlangt den Werkzeugwechsel nur, wenn es **keine** Ausweichform
  gibt — „berichtete manuelle Nacharbeit" existiert hier und ist berichtet
  (schema.yaml, Makefile, `nacharbeit-operation-check.sql`); der
  Werkzeugwechsel bleibt damit korrekt **Prüfauftrag an die Architect-Rolle**,
  nicht Entscheid des Implementer-Laufs. Der DROP-Guard macht den Schritt
  wiederholbar; die Runner-Kette räumt vor dem Rollout ab, so dass der
  E012-Fall (zweiter Rollout gegen bestückte Instanz) im automatisierten
  Pfad nicht auftritt — die Grenze ist im Makefile-Kommentar benannt
  (Implementer-Risiko (a): getragen). Offen: F-4 (Retirement-Verankerung),
  F-5 (Beleg-Kette).
- **(e) d-check/a-check-Befund-Fixes des Laufs:** *alle sachlich richtig.*
  `.a-check.yml`-Glob-Erweiterung: die Verdrahtungs-Tests tragen
  Composition-Root-Verhalten ([`ADR-0026`](../plan/adr/README.md)-Konsequenz „der Aufbau ist
  als Test reproduzierbar"); ohne Eintrag fiel der Pfad aus der Abdeckung —
  a-check 0 Befunde am Head. `.PHONY`-Fix (`f79af17`): das schattierte
  `test`-Target war still tot; der Fix trägt die Ursache im Kommentar
  (Klasse Kopplung) und `make test` läuft im eigenen Probe-Lauf. Der
  Ruhe-Marker-Fix (`779bc64`) ist die deklarierte Redundanz in der richtigen
  Richtung (Marker steht genau dann, wenn `in-progress/` leer ist — Modul 6).
- **(f) kein `build:`-Block in `compose.yaml`:** *Konform mit [`ADR-0044`](../plan/adr/README.md)*
  — der Zweit-Build würde das geladene Tag vom Lauf-Beleg entkoppeln; die
  Abgrenzung steht als Kommentar-Klasse Abgrenzung/Kopplung in der Datei.
  Grenze: F-8 (Beleg-Abgleich).
- **(g) `D_MIGRATE_RUN_USER` (uid 10001):** *Grenze dokumentiert, getragen* —
  der Makefile-Kommentar nennt die Image-uid und den Host-uid-Ausweg
  (Kommentar-Klasse Kopplung); der eigene Erstlauf belegt den Pfad. Kein
  Befund.
- **(h) [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)/`SPEC-008` am Rollout-Fehlerpfad:** die Fehlerklassen-ADR
  bindet an Port-Kontrakte und bleibt unberührt; der Rollout-Pfad trägt
  sichtbare Fehler (`docker run`-Exit, `psql -v ON_ERROR_STOP=1`). Kein
  Befund. [`ADR-0029`](../plan/adr/README.md) Regel 1/Idempotenz: der Diff trägt keinen
  Produkt-Code; der Test verdrahtet Persist-before-ACK über den echten
  Service am realen Treiber. Unberührt.

## Test-Qualität — MVP-Schnitt am realen Pfad

`make test-integration` am Range-Head: **grün** (eigener Lauf; Kette
Compose frisch → Rollout → Nacharbeit → Tests). Der MVP-Test deckt die
Schnitt-Abfolge vollständig ab, und zwar an den Zusagen, die die
Vorgänger-Slices an Tests trugen:

| MVP-Schritt | Träger im Test |
|---|---|
| PostgreSQL starten | Compose-`postgres`, Digest-Pin, `wal_level=logical` |
| CDC aktivieren | direktes SQL: Publication + Bindungs-Zeilen in `cdc.source`/`source_table`/`schema_version` (F-10) |
| INSERT/UPDATE/DELETE | drei getrennte Quelltransaktionen am realen Treiber |
| Changes lesen | `ReadChanges` über den echten Store-Adapter (Polling mit Zeitgrenze) |
| Reihenfolge | strenge Positionsordnung (`LH-FA-CAP-004`), Sequenz je Transaktion, Wiederlesen in derselben Ordnung (`LH-FA-REA-004.a`/005), Bereich hinter letzter Position leer (`LH-FA-REA-002`) |
| Inhalt | INSERT-Neu-Bild, DELETE-Alt-Bild (Schlüsselspalte), UPDATE-Alt-Bild abwesend bei Default-Identity — **Gegenprobe** mit `REPLICA IDENTITY FULL` trägt beide Spalten (`LH-FA-CAP-008`, [`ADR-0016`](../plan/adr/README.md)) — das trägt den slice-006-§6-Risiko-(b)-Ausgang am Beleg |

Die §6-Risiken des Plans sind damit am Beleg: (a) reale Ordnung
(streige Positionsordnung) und (b) Row-Images (beide Identity-Formen).
Die Cleanup-Kette droppt Slot (mit Retry, erst nach dem Stream-Ende),
Publication und Tabelle; der Runner räumt Compose und Netz in jedem
Ausgang ab (der slice-004-F-5-Hygiene-Defekt ist nicht wiederholt). Der
Test skippt ohne DSN sauber und läuft damit auch treiberfrei in `make test`.

## Implementer-Risiken — Bewertung

- **(a) Nacharbeit-Brücke:** Grenze benannt und im automatisierten Pfad
  getragen (Runner-Räumung) — **F-4** (kein terminierter Ausgang) und
  **F-5** (Beleg endet vor der Nacharbeit) bleiben als Plan-/Beleg-Lücken.
- **(b) Feed-Container als Smoke:** **F-6** — benannt in Code und Commit,
  aber ohne Ausgangs-Träger; die M1-Verkabelung in `welle-2.md` §1 nennt
  die Grenze nicht.
- **(c) `D_MIGRATE_RUN_USER`:** dokumentiert (Makefile-Kommentar, uid
  10001), eigener Lauf belegt — kein Befund.

## Zweitschreiber-Frage — `3ea9244` (Allowlist)

Der Commit liegt im Range, ist aber kein Implementer-Produktcode (`.claude/
settings.json` + Hook-Verweise, Dateien existieren). Bewertung: **zulässige
Konfig-Arbeit auf Nutzer-Vorgabe** — mit zwei Residuen: (1) der Commit trägt
keine Vertrags-Kennung (F-2); (2) die Allowlist wirkt für jeden Agenten-Lauf
und trägt ihre Abwägung (was bewusst *nicht* gelistet ist) nur im
Commit-Text — die Datei selbst trägt keine Kommentar-Ebene; die Abwägung
gehört in eine Adresse, die der Lauf liest, nicht in die git-Historie
([`AGENTS.md`](../../AGENTS.md) §3.7, Abwägung gehört in die ADR — hier keine
entscheidungstragende Abwägung, aber der „bewusst nicht in der Liste"-Satz
ist eine Zusage, die nur historisch lesbar ist). INFO-Stufe, kein
Blockier-Befund.

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff (abgesehen von F-1
  und F-2)** — kein ADR-Verstoß auf Layer-/Tool-Ebene (a-check 0 Befunde;
  `make test-integration` verdrahtet real über den Composition-Root-Aufbau),
  kein Sicherheits-Anti-Pattern (Credentials bleiben im Compose-Test-Netz;
  `plan.yaml` scrubbed das Target-Passwort `***`; psql-Lauf mit ro-Mount),
  kein Korrektheitsfehler im kritischen Pfad (der Test belegt
  Persist-before-ACK am realen Treiber; der Adapter-Code ist unberührt),
  keine Gate-Suppression, keine Norm nur im Template-Kommentar, kein
  Chronik-tragendes Zustandsfeld im Diff (der Ruhe-Marker-Fix entfernt
  Chronik-Richtung), **kein Docker-only-Verstoß** (alle Läufe über
  make/docker; Digeste gepinnt: d-migrate `8d1433…`, postgres `63bdc97…`,
  Toolchain `cf6fca…`; `go mod download` im Container)
- geprüft, ohne Befund: **Traceability der übrigen vier Implementer-Commits**
  — jeder trägt mindestens eine `LH-*`-/`ADR-*`-Kennung; alle genannten IDs
  existieren ([`LH-FA-CAP-001`](../../spec/lastenheft.md)…003/004/008, [`LH-FA-REA-002`](../../spec/lastenheft.md)/004.a/005,
  [`LH-FA-CFG-001`](../../spec/lastenheft.md), [`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0026`](../plan/adr)/0043/0044)
- geprüft, ohne Befund: **Struktur-IDs in Commit-Messagen** — keine
  `SPEC-*`-/`ARC-*`-Kennung im Range (die slice-005-F-5-Klasse zählt hier
  **nicht** weiter); `SPEC-*`-Nennungen in Datei-Kommentaren
  (`schema.yaml`, `compose-init`, `nacharbeit`) sind zulässige Rang-Zeiger
- geprüft, ohne Befund: **superseded-ADR-Referenzen** — die neuen Referenzen
  nennen nur Accepted-ADRs; keine Referenz auf 0038/0039 als tragenden Anker
- geprüft, ohne Befund: **Kommentar-Klassen in den neuen Dateien (§3.7)** —
  `compose.yaml`, `Makefile`-Blöcke, `run-integration-tests.sh`,
  `schema.yaml`, `nacharbeit-operation-check.sql`, `01-cdc-schema.sql`,
  `mvp_test.go`, `.a-check.yml` — Indikativ über den geltenden Zustand,
  Klassen Zusage/Kopplung/Abgrenzung/Grenze/Rang-Zeiger; kein
  Konjunktiv über verworfene Alternativen, kein abwesender Text
- geprüft, ohne Befund: **Runner-Hygiene** — `down -v` in jedem Ausgang
  (trap EXIT) und vor dem Start, Readiness-Wait, Feed-Smoke, Read-only
  Mount für den psql-Schritt, Modul-Cache-Volume erhalten; gepinnte
  Digests konsistent (Toolchain/PostgreSQL wie slice-004/005)
- geprüft, ohne Befund: **§3-Datei-Menge** — `compose.yaml`,
  `test/integration/*`, Makefile-Verkabelung und `tools/schema/*` decken die
  vier §3-Zeilen (inkl. Plan-Nachzüge); darüber hinaus geliefert:
  `.a-check.yml`, `harness/README.md` (F-3), `.claude/settings.json`
  (Zweitschreiber-Abschnitt), `roadmap.md` (Marker-Übergang, Modul-6-konform)
- geprüft, ohne Befund: **§3.3 im Range** — `c286339` ist ein reiner
  `git mv` (R100, kein Inhaltsanteil); keine mv/content-Mischung im Range
- geprüft, ohne Befund: **Rollout-Beleg-Determinismus** — der
  `test-integration`-Re-Lauf des Reviews lässt `plan.yaml`/`down.sql`
  byte-identisch; die committeten Belege dirtien nicht
- geprüft, ohne Befund: **Gate-Läufe am Range-Head** — `baseline-verify` OK
  (54 Dateien) · `d-check` 95 Dateien/0 Befunde · `a-check` 0 Befunde ·
  `make schema-validate` grün · `make test` grün · `make test-integration`
  grün (Nacharbeit-Schritt im Lauf sichtbar; `make image` nicht berührt —
  keine Build-Kontext-Änderung des Ranges)
- geprüft, ohne Befund: **Beobachtungs-Register** — keine neue Beobachtung
  im Diff angefallen; die Register-Pflichten sind Closure-Sache (§7 steht
  aus)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 5 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zwei-Quellen-Drift der Schema-Form
(`committed_at` NOT NULL, Probe belegt) · Commit ohne Vertrags-Kennung (3.
Auftreten — Sequenz-Pflicht fällig, 2 Commits) · Plan-Erweiterung ohne
Plan-Nachzug (6. Auftreten — Sequenz läuft) · Befristete Nacharbeit ohne
terminierten Ausgang (Retirement ungeterminiert) · Beleg-Kette endet vor der
Nacharbeit (Report deckt Endzustand nicht) · Verdrahtungs-/Aktivierungs-Grenze
ohne Ausgangs-Träger (Feed-Smoke, Bootstrap-Stub, kein Folge-Slice) ·
Datei-Abschluss (5 neue Dateien, Klasse MEDIUM) · Beleg-Abgleich fehlt
(`:dev`-Tag) · Ziel-Wortlaut vs. Lieferung (Test-Consumer) · Aktivierung per
SQL am MVP-Rand (marker-konsistent) · undeklarierte ID-Präfixe (Vorbestand
`PH-*`) · Vorlagen-Platzhalter in welle-2 §3 (Vorbestand)

**Sequenz-Beobachtung (Steering-Loop):** die Klasse „Commit ohne
Vertrags-Kennung" erreicht mit diesem Lauf das **dritte Auftreten**
(review-slice-005 F-1 → hier 3ea9244, 779bc64) — Verkörperung über den
**Architect** fällig (Gate oder Conventions-Eintrag), parallel zur weiterhin
nicht gelaufenen Struktur-ID-Sequenz aus review-slice-003.

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig, die
MVP-Schnitt-Abfolge ist am realen Treiber belegt (Reihenfolge, Inhalt,
Wiederlesen, Bereichs-Grenze, beide Identity-Formen), der Rollout-Erstversatz
ist [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-konform mit committetem Pflicht-Report und Rollback-Artefakt,
alle drei Gates laufen grün, und der MVP-Integrationstest-Claim trägt seine
Grenzen — im Code und Commit-Text offen benannt.

**Blockierend für Closure:** ja, in vier Punkten, bevor der Slice nach
`done/` geht:

1. **F-1 (Schema-Drift):** die Überführungs-Abweichung (`committed_at`
   NOT NULL; nachrichtlich auch die Constraint-Namen) ist vor der Closure zu
   schließen — entweder Korrektur in `schema.yaml` (`required: true`) mit
   Erneuerung von Report/Rollback-Beleg, oder die Abweichung wird als
   berichtete Grenze derselben Nacharbeit-Klasse getragen (dann mit
   F-4-Verankerung). Der [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Befund-Fall verlangt eine der beiden
   Formen, nicht den Stillstand.
2. **F-4 (Retirement-Verankerung):** die Nacharbeit-Brücke braucht einen
   Ausgangs-Träger **vor** der Closure — §6-Risiko mit Ausgang „weiter
   offen" → `BEO-PGC/<slug>` im Beobachtungs-Register oder Folge-Slice-ID in
   `open/`. Übergabe-Artefakt: Plan-Nachzug über den **Planner**.
3. **F-6 (Verdrahtungs-/Aktivierungs-Grenze):** derselbe Mechanismus —
   Ausgang im Plan §6 oder Welle-§1-Formulierung schärfen (M1-Beleg auf
   Adapter-Ebene benennen); ohne Ausgang geht der Slice mit einem
   ungeterminierten Claim in `done/`.
4. **F-2 (Kennungspflicht, 3. Auftreten):** wirkt nur vorwärts, löst aber
   die Sequenz-Pflicht aus — Übergabe-Artefakt an den **Architect**; die
   Closure §7 trägt die Klasse als wiederkehrende Finding-Klasse
   (Summary-Zeile dieses Reports).

F-3 (Plan-Nachzug `.a-check.yml`/`harness/README`), F-5 (Beleg-Kette) und
F-7 (Datei-Abschluss) gehen als Vor-Closure-Nacharbeit bzw. in die
Closure §7. F-8/F-9 sind LOW ohne Blockier-Charakter. F-10–F-12 gehen an
Planner bzw. in die Closure §7. DoD- und Spec-Konformität prüft der Verifier
separat (Modul 11) — insbesondere die vollständigen DoD-Häkchen und der
Beleg „`make gates` grün" an der finalen Range-Fassung.