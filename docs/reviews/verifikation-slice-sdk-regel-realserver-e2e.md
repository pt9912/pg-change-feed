# Verifikations-Report: slice-sdk-regel-realserver-e2e — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-sdk-regel-realserver-e2e.md`](review-slice-sdk-regel-realserver-e2e.md);
Hausform dieses Reports:
[`verifikation-slice-backfill-change-origin.md`](verifikation-slice-backfill-change-origin.md).

**Gegenstand:** Slice-Plan
`docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md`,
Diff-Range `300637cf..98c7f8e4` (5 Commits): Feature-Commit `ad3af754`
(Umsetzung, drei Tiers), `47efceee` (§3-Suchlauf-Nachtrag),
Review-Report `bd3b99f1` (1 HIGH F-1, 1 LOW F-2), `300602a1` (Kennungs-Link-
Fix im Review-Report — von außen dazwischengekommen, nur die Review-Datei
betroffen, harmlos), Fixrunde `98c7f8e4` (F-1/F-2 real geschlossen). Slice-
Parent für die Code-Diffs ist `71024045`. Dieser Lauf ändert keine Spec,
keinen Plan und keinen Code; er schreibt nur diesen Report (und eine
selbst gefahrene, exakt zurückgenommene Mutation, siehe §3.2).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make gates` | **EXIT=0** | baseline-verify OK; `docs-check`/d-check 1362 Datei(en), 0 Befund(e); `a-check` „gesamt: 0 Befund(e)"; `commit-traceability` OK; `generated-sync` OK (byte-gleich, Stufe `proto-export`); coverage-gate lief im selben Zug |
| `make commit-traceability RANGE=300637cf..HEAD` | **EXIT=0** | „commit-traceability: OK — 5 Commit(s) in "300637cf..HEAD", Betreffs ohne Struktur-ID" |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md` | **EXIT=0** | „suchlauf-nachmessen: 22 Zeilen stimmen" — alle 22 Zeilen (Parent `71024045` und Stand `diff`) selbst nachgemessen, keine Abweichung |
| `git diff --name-only 71024045 98c7f8e4 -- sdks` | 14 Pfade | ausschließlich unter den drei Test-Verzeichnissen (`PgChangeFeed.Client.Integration/`, `src/integrationTest/kotlin/…/integration/`, `pgchangefeed/integration/`); kein `.csproj`/`pyproject.toml`/`build.gradle.kts` |
| `make sdk-public-doc-check` | **EXIT=0** | „sdk-public-doc-check: keine interne Kennung unter sdks" |
| `make kommentar-kennungen DIFF=300637cf` | **EXIT=0**, 0 Kandidaten | kein Treffer (Diff berührt keine Go-Datei) |
| `make doc-trace` | **EXIT=0** | 80 Anforderungen, 1 Waise (`LH-FA-CFG-008`, unverändert erwartet); `LH-FA-CFG-007`-Zeile führt jetzt die Dimension `SDK-E2E` zusätzlich zu `E2E` |
| `git diff 71024045 98c7f8e4 -- docs/user/sdk-e2e-abdeckung.md` | 3 neue Zeilen | genau eine Regel-Zeile je Sprache (C#, Kotlin, Python), marker-gegrenzte Abschnitte sonst unverändert (nur Kosmetik: pro Sprachabschnitt jetzt eine eigene Tabellenkopfzeile statt einer gemeinsamen — Rendering-neutral, kein Befund) |
| `git diff 71024045 98c7f8e4 -- harness/mk/sdk.mk harness/README.md` | gelesen | `sdk.mk`: „vier Phasen" → „acht Phasen — vier ohne Regel, vier mit rename_column-Regel" an beiden Zielen; `README.md`: alle drei `make test-sdk-*-integration`-Zeilen nennen jetzt die Regel-Phasen und den Slice-Namen |
| eigene Mutationsprobe (§3.2) | rot wie erwartet, danach exakt zurückgenommen | siehe unten |

Kein Lauf durch eine Pipe gefiltert; jeder Exit-Code wurde unmittelbar nach
dem jeweiligen `make`-Aufruf gelesen.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 (Hilfsdatei + C#-Tier) | **erfüllt** | `tools/harness/lib-sdk-rule-fixture.sh` (gelesen, exakt eine Aufrufform, Poll mit Frist und `failed`-Fehlertext); vier neue C#-Testklassen; Feature-Commit-Message nennt alle drei Mutationen real am C#-Tier gefahren — der Umfang deckt sich mit dem, was ich für Kotlin unabhängig nachgefahren habe (§3.2) |
| 2 | Liefer-Punkt 2 (Kotlin-Tier) | **erfüllt** | Fixrunde `98c7f8e4` trägt jetzt Mutation (b) und (c) real gefahren und rot gesehen (gRPC-Fläche), zusätzlich zu Mutation (a) aus dem Feature-Commit; **von mir unabhängig ein viertes Datum nachgefahren** — Mutation (c) an der SSE-Regel-Fläche (weder Reviewer noch Fixrunde hatten dort getestet, nur an der gRPC-Fläche) — rot wie erwartet, siehe §3.2 |
| 3 | Liefer-Punkt 3 (Python-Tier) | **erfüllt** | dieselbe Fixrunde-Form; Reviewer hatte zusätzlich Mutation (b) am Python-gRPC-Weg selbst nachgefahren (Review-Report, rot: `AssertionError: assert '501' == 'PythonGrpcRuleSdkE2ESentinel'`) — deckungsgleich mit der jetzt im Plan dokumentierten Fixrunde-Messung |
| 4 | Nur Test-Code und Runner | **erfüllt** | eigener `git diff --name-only` (§1); `make sdk-public-doc-check` EXIT=0 |
| 5 | Abdeckung getragen | **erfüllt** | eigener Diff der Zieldatei (§1): genau drei neue Zeilen; `make doc-trace` EXIT=0 mit `SDK-E2E`-Dimension an `LH-FA-CFG-007` |
| 6 | `make gates` grün | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 7 | Review durchgeführt, Report vorliegt | **Report liegt vor, Kästchen bleibt formal falsch — V-1 (LOW)** | `review-slice-sdk-regel-realserver-e2e.md` existiert, 0 offenes HIGH/MEDIUM nach der Fixrunde; die Checkbox ist aber weiterhin `[ ]`, obwohl `.claude/commands/implement-slice.md` Schritt 21 verlangt, dass die Fixrunde sie „im selben Commit" mitsetzt, wenn sie den einzigen offenen HIGH-Punkt auflöst — siehe §5 |
| 8 | §3.13-Suchlauf, Feld trägt Gefundenes und Nichtgefundenes, beide Stände | **erfüllt** | eigener Lauf `make suchlauf-nachmessen`: 22/22, EXIT=0 (§1) |
| 9 | Doku-Update (`harness/README.md`, `harness/mk/sdk.mk`, Kopf-Kommentare) | **erfüllt** | §1, Diff gelesen |
| 10 | Closure-Notiz | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 11 | Reconciliation-Register „entfällt" | **korrekt** | Greenfield, `[x]` mit Begründung „entfällt: keine Reconciliation-Datei in diesem Repo" |
| 12 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 13 | Risiken §6 je ein Ausgang | **korrekt offen** | jede der elf Risiko-Zeilen trägt „**Ausgang:** *(bei Closure)*" |
| 14 | Drei Paarungen | **korrekt offen** | hängen an der Slice-Closure selbst (kein Welle-Bezug) |

Kein `[x]` ohne Beleg außer der einen benannten Formal-Abweichung (Zeile 7,
V-1); kein `[ ]`, das über die Rollen-Sequenz hinaus schon belegt wäre.

## 3. Kernaussagen, selbst geprüft

### 3.1 Plan-Text der Fixrunde gegen die Commit-Message

Die Rot-Belege im Fixrunde-Nachtrag (Plan-Datei Zeilen 356–390) sind
wörtlich identisch mit denen der Commit-Message `98c7f8e4` (Assertion-Texte,
Zeitfenster, Laufzeiten). Die Kostenklasse-Werte (`free -m`, dangling
Volumes) stehen jetzt als eigener Absatz mit Ursprung „gemessen" im Plan
(Zeile 383 ff.) — F-2 damit sachlich geschlossen. Beide Nachträge sind reine
Doku-Commits (57 Einfügungen/2 Löschungen in genau der Plan-Datei, `git show
98c7f8e4 --stat` gelesen); kein Code, kein Test wurde dabei angefasst.

### 3.2 Eigene Mutationsprobe — ein von Reviewer/Implementer nicht gefahrenes Datum

Um F-1 nicht nur textlich, sondern an einer neuen Stelle zu prüfen, habe ich
eine vierte, bislang ungetestete Kombination gefahren: **Mutation (c) an der
Kotlin-**SSE**-Regel-Fläche** (Reviewer und Fixrunde hatten (b)/(c) beide
Male nur an der **gRPC**-Regel-Fläche unabhängig erprobt — die übrigen drei
Flächen blieben bislang strukturell, nicht real belegt).

- Bearbeitet: `tools/harness/run-sdk-kotlin-integration-tests.sh`, Zeile
  354 (`SSE_RULE_IDENT`-Aufruf), `sql_kind` von `changes_renamed` auf
  `changes` geändert — per Read+Write, kein `sed -i` (`AGENTS.md` §3.1).
- `git diff --stat` vor dem Lauf: genau eine geänderte Zeile in dieser
  Datei.
- `make test-sdk-kotlin-integration` real ausgeführt (22:03:30–22:04:47,
  ~77 s) — **EXIT=2** (rot, wie erwartet), Log:
  `run-sdk-kotlin-integration-tests: SSE-Regel-Fläche (SPEC-021, ADR-0112)
  — die Identität (844-1) ist nicht real über den SQL-Lesezugriffsweg
  lesbar (count=0)` — dieselbe Fehlerklasse wie an der bereits erprobten
  gRPC-Fläche, jetzt zusätzlich an der SSE-Fläche bestätigt.
- Rücknahme: `git checkout -- tools/harness/run-sdk-kotlin-integration-tests.sh`;
  `git status --short` und `git diff --stat` danach beide leer.
- Dangling Docker Volumes vor/nach meinem Lauf: 39/39 (unverändert, wie im
  Fixrunde-Nachtrag behauptet — der Runner räumt seinen eigenen
  Compose-Stack per `compose down -v` ab).

Das bestätigt: die Assertion-Bindung an `changes_renamed` ist keine
Eigenschaft, die nur an der gRPC-Fläche zufällig hält, sondern strukturell
in `run_phase`s SQL-Variante verankert (dieselbe Funktion bedient alle vier
Flächen) — ein zusätzliches, unabhängiges Datenpunkt über das hinaus, was
Review und Fixrunde bereits geprüft hatten.

### 3.3 Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 8 | keine Änderung des Nachrichtenschemas, keine SDK-Code-Änderung erwartet | `git diff --name-only` zeigt ausschließlich Test-Verzeichnisse; `make sdk-public-doc-check` EXIT=0 | konform |
| [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 6 | „SDK-/Doku-Beleg: bestätigen, dass die SDK-Tests keine Image-Schlüssel voraussetzen" | genau das ist der Slice-Gegenstand; zwölf neue Testklassen/-dateien belegen die Bild-Prüfung über das opake Modell (`JsonElement`/`Any`) statt über einen festen Schlüssel | konform, Folgepflicht real eingelöst |
| [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung Festlegung 2 | Mechanik-Klasse der Realserver-Integrationstests (compose hoch, `integration`-Docker-Stufe, Phase je Fläche) | die drei Runner erweitern exakt diese bestehende Mechanik um vier weitere Phasen je Tier, keine neue Mechanik | konform |
| [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md)/[`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md) | Aufrufform von `cdc.set_transformation` (JSON-Regelform) | `lib-sdk-rule-fixture.sh` ruft `cdc.set_transformation('src-e2e', 'public', …, '{"kind": "rename_column", "column": "…", "to": "…"}')` — Literalform wie in [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) Folgepflicht 2 | konform |
| `AGENTS.md` §3.1 (Docker-only, kein `sed -i`) | Textänderung nur über Edit/Write oder Repo-Werkzeug | eigene Mutationsprobe (§3.2) per Read+Write, keine In-Place-Bearbeitung; Runner-/Fixture-Diff enthält keinen Host-Toolchain-Aufruf | konform |

## 4. Harte Regeln

- **§3.3** (git mv + Inhalt = zwei Commits) — n/a, kein Lifecycle-Move in
  diesem Diff.
- **§3.5** (ADR-Immutabilität) — `docs-check`/`a-check` 0 Befunde; keine
  `Accepted`-ADR im Diff berührt (nur Plan- und Testdateien).
- **§3.7** (Kommentarform) — `make kommentar-kennungen DIFF=300637cf`
  EXIT=0, 0 Kandidaten; der Kopf-Kommentar von
  `lib-sdk-rule-fixture.sh` trägt genau einen Anker
  (`slice-sdk-regel-realserver-e2e, ADR-0112`) plus eine BEO-Kennung, keine
  Kette im Sinne der Regel (vom Reviewer bereits geprüft, von mir per
  Sensor bestätigt).
- **§3.9** (Exit-Code direkt, nie durch Pipe) — jeder eigene Sensor-Lauf
  wurde einzeln ausgeführt und sein Exit-Code unmittelbar danach gelesen,
  vor jeder Filterung.
- **§3.12** Instanz A/B — die Kostenklasse-Werte im Fixrunde-Nachtrag tragen
  den Ursprung „gemessen"; die Assertion-Zitate sind wörtliche Log-Auszüge,
  keine Paraphrase.
- **§3.13** (Suchlauf) — eigener Lauf `make suchlauf-nachmessen`: 22/22,
  EXIT=0.
- **§3.15** (keine Ersatzhandlung nach verweigerter Aktion) — kein Hinweis
  auf eine verweigerte Aktion in diesem Diff; nicht einschlägig.

## 5. Befund

**V-1 (LOW) — DoD-Checkbox „Review durchgeführt" wurde im Fixrunde-Commit
nicht mitgesetzt, obwohl der Prozess das verlangt.**

- `pfad`: `docs/plan/planning/in-progress/slice-sdk-regel-realserver-e2e.md:217`
- `befund`: `.claude/commands/implement-slice.md` Schritt 21 sagt wörtlich:
  „Löst eine Fixrunde nach Reviewer-Findings einen bislang offenen
  DoD-Punkt auf (typischerweise „Review durchgeführt, … kein offenes
  HIGH"), wird die zugehörige Checkbox **im Fixrunden-Commit** mitgesetzt —
  nicht erst bei der Planner-Closure nachgetragen." Der Fixrunde-Commit
  `98c7f8e4` löst den einzigen offenen HIGH-Punkt (F-1) und den einzigen
  LOW-Punkt (F-2) real auf — 0 offene HIGH/MEDIUM-Findings danach —, setzt
  aber die Checkbox „Review durchgeführt, Report unter `docs/reviews/`
  liegt vor" nicht auf `[x]`. Zum Vergleich: im Referenz-Fall
  `slice-backfill-change-origin` wurde dieselbe Checkbox korrekt im
  Fixrunde-Commit gesetzt (`daa86db5`, siehe
  [`verifikation-slice-backfill-change-origin.md`](verifikation-slice-backfill-change-origin.md)
  §2 Zeile 6).
- `verifizierbar`: ja — `git show 98c7f8e4` (Diff der Plan-Datei) zeigt
  keine Änderung an dieser Zeile.
- **Einordnung:** dies ist die „sichere" Richtung des Fehlers (Checkbox
  untertreibt den Fortschritt, sie behauptet keine Fertigstellung, die
  nicht vorliegt — anders als F-1 selbst, das die gegenteilige, unsichere
  Richtung hatte). Kein DoD-Inhalt ist dadurch unbelegt; die Substanz hinter
  der Checkbox (Review fand statt, Report liegt vor, 0 offenes HIGH/MEDIUM)
  ist real gegeben und von mir unabhängig bestätigt (§1, §2 Zeile 7).
- **Blockiert nicht.** Der Planner zieht die Checkbox beim Closure-Nachzug
  mit; kein neuer Implementer-Lauf nötig.

Kein weiterer Befund. Kein offenes HIGH/MEDIUM aus dem Review-Report.

## 6. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (9) | **9 von 9 inhaltlich erfüllt**; eine (Zeile 7) mit formaler Abweichung V-1 (LOW), nicht blockierend |
| DoD §2 — übrige `[ ]`/entfällt-Zeilen (5) | **5 von 5 korrekt** (Closure-/Rollen-Sequenz) |
| Plan-vs-Code-Diff | deckungsgleich: Hilfsdatei + drei Runner + zwölf Testklassen/-dateien + drei Abdeckungs-Zeilen + zwei Doku-Träger — nichts außerhalb, kein SDK-Produktivcode |
| F-1 (Review, HIGH) | **real geschlossen** — alle drei Mutationen an allen drei Tiers rot gesehen; zusätzlich von mir ein viertes, bislang ungetestetes Datum (Kotlin SSE-Regel-Fläche, Mutation c) unabhängig bestätigt |
| F-2 (Review, LOW) | **real geschlossen** — Kostenklasse-Werte im Plan nachgetragen, mit Ursprung „gemessen" |
| Entscheidungen `ADR-0112`/`ADR-0110`/`ADR-0125`/`ADR-0126` | **konform** |
| Harte Regeln (§3.1, §3.3, §3.5, §3.7, §3.9, §3.12, §3.13, §3.15) | **erfüllt** |
| Gates | `make gates` EXIT=0, `make commit-traceability RANGE=300637cf..HEAD` EXIT=0, `make suchlauf-nachmessen` 22/22 |

## Verdikt

**Bereit für Planner-Closure.** Die DoD-Substanz ist durch eigene Messung
bestätigt — inklusive einer eigenen, von Reviewer und Implementer nicht
gefahrenen vierten Mutationsprobe (Kotlin SSE-Regel-Fläche), die real rot
wurde und danach exakt zurückgenommen wurde (`git status --short` leer).
Kein offenes HIGH oder MEDIUM. Der einzige Befund (V-1, LOW) betrifft eine
Checkbox, die den bereits real erbrachten Fortschritt untertreibt, nicht
überzeichnet — er blockiert die Closure nicht, sollte aber beim
Closure-Nachzug mitgezogen werden.

**Übergabe:** Verifier → Planner. Beim Closure-Nachzug: Checkbox „Review
durchgeführt" auf `[x]` setzen (V-1), Closure-Notiz mit
Steering-Loop-Lerneintrag schreiben, Beobachtungs-Register fortschreiben,
die elf §6-Risiken je mit einem Ausgang versehen, die drei Paarungen
(Anker · Folge-Slice `slice-transformationen-betriebsdoku` · Register)
tragen. Danach der reine `git mv` nach `done/` (`AGENTS.md` §3.3, Fall 2).
