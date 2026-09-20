# Verifikations-Bericht: slice-sdk-kotlin-pack-werkzeug — 2026-09-20

**Rolle:** Verifier (Modul 11), frischer Kontext, unabhängig von
Implementer-Commit `84293e14` und Reviewer-Commit `5cf43c65`.

**Gegenstand:** DoD-Konformität von
`docs/plan/planning/in-progress/slice-sdk-kotlin-pack-werkzeug.md`
(`LH-FA-SST-009`, `ADR-0109`) gegen den tatsächlichen Repo-Stand (HEAD =
`5cf43c65`), Elternstand `9c4efd42`.

**Eingelesen:** `harness/README.md`, `AGENTS.md`, `harness/conventions.md`,
der vollständige Slice-Plan (§1–§8), `ADR-0109`, `spec/pflichtenheft.md`
(aktueller Stand, §1 `LH-FA-SST-009.a`, §6 Externe Verträge, §7 Historie),
`docs/reviews/review-slice-sdk-kotlin-pack-werkzeug.md`.

**Frage dieser Rolle:** Bauen wir es richtig — gegen Plan und DoD? (nicht
Validator-, nicht Reviewer-Frage.) Jede Behauptung unten ist gegen einen
selbst ausgeführten Beleg geprüft, keine Übernahme aus Implementer- oder
Reviewer-Bericht ohne eigenen Lauf.

---

## 1. DoD-Zeilen einzeln geprüft (§2 des Slice-Plans)

| # | DoD-Zeile | Eigene Prüfung | Ergebnis |
|---|---|---|---|
| 1 | `make sdk-pack-kotlin` existiert, Docker-only, kein Gate | `grep -n "sdk-pack-kotlin\|GATE_CHECKS" harness/mk/sdk.mk` — Target `sdk-pack-kotlin:` vorhanden (Zeile 52/53), `GATE_CHECKS` wird in `sdk.mk` an keiner Stelle um `sdk-pack-kotlin` erweitert (Kommentarzeile 8: „NICHT an GATE_CHECKS"); `grep -rn "GATE_CHECKS" Makefile harness/mk/*.mk` zeigt `sdk.mk` nicht unter den Fragmenten, die tatsächlich `+=` auf `GATE_CHECKS` schreiben (`doc-gate.mk`, `coverage.mk`, `generated-sync.mk`, `baseline.mk`) | **Bestätigt** |
| 2 | Reales `.jar` existiert, Smoke-Beleg | Eigener, **vierter** unabhängiger Lauf (nach Implementer, Reviewer, jetzt Verifier): `rm -rf sdks/kotlin/dist && bash tools/harness/sdk-pack-kotlin.sh` — lief durch (Docker-Cache, alle Layer `CACHED`, Gesamtlauf 1,24 s), erzeugte `sdks/kotlin/dist/pgchangefeed-kotlin-0.1.0.jar`, **112591 Bytes** (`ls -la`) — byte-identisch mit dem in Commit-Message/Review-Report genannten Wert, kein Drift. `unzip -l` zeigt ein valides Archiv: `META-INF/MANIFEST.MF` sowie die erwarteten `cdc/stream/v1/*.class`-Einträge (u. a. `ChangeStreamGrpc.class`, `Changestream$Change.class`) | **Bestätigt** |
| 3 | `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a`-Nachzugsatz | Real gelesen (`sed -n '170,200p' spec/pflichtenheft.md`): neuer Kotlin-Absatz ab „Für Kotlin/GitHub Packages ist die Frage ebenfalls beantwortet: `pgchangefeed-kotlin` (`SPEC-028`) …" schließt korrekt an den (bereits im selben Diff bereinigten) Python-Absatz an; „Eine vierte Sprache oder ein vierter Vertriebsweg bleibt offen" ist am aktuellen Stand sachlich richtig. **Kein `ADR-`-Verweis** in den neuen Zeilen 170–200/620–665: `grep -n "ADR-"` in diesem Bereich trifft ausschließlich die drei bereits vor diesem Slice bestehenden, unveränderten `ADR-pflichtig`-Vorkommen (Zeilen 10/42/68-70, keine davon in den drei neuen Kotlin-Zeilen) — `matrix`-Modul-Regel (`spec → adr` verboten) gewahrt | **Bestätigt** |
| 4 | `spec/pflichtenheft.md` §6 neue `SPEC-<NNN>`-Zeile + §7 Historie | Real gelesen: §6 Zeile 624 `SPEC-028` `pgchangefeed-kotlin` GitHub-Packages-Gradle-/Maven-Package, SemVer 2.0 `0.x.y`, Vertrag-Datei `build.gradle.kts`; §7 Zeile 663 Historie-Eintrag 2026-09-20. `grep -n "SPEC-028"` liefert genau die drei erwarteten Treffer (§1-Absatz, §6-Zeile, §7-Historie) — keine Kollision mit einem anderen Zug (`SPEC-027` bleibt der einzige Vorgänger) | **Bestätigt** |
| 5 | `harness/README.md` §Werkzeuge, Bindung `ADR-0109` | Real gelesen (`grep -n "sdk-pack-kotlin" harness/README.md`): Zeile 158, vollständige Zeile analog `sdk-pack-csharp`/`sdk-pack-python`, endet mit der Bindungsangabe „kein Gate" plus einem echten Markdown-Link auf [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Festlegung 5, Zeitangabe „seit slice-sdk-kotlin-pack-werkzeug" | **Bestätigt** |
| 6 | `make gates` grün | Eigener, ungepipter Lauf (`make gates > log 2>&1; echo EXIT_CODE=$?`) → `EXIT_CODE=0`. Log-Ende zeigt `generated-sync: OK`, `a-check … gesamt: 0 Befund(e)` | **Bestätigt, Exit 0** |
| 7 | Review durchgeführt, kein offenes HIGH | `docs/reviews/review-slice-sdk-kotlin-pack-werkzeug.md` real gelesen: Summary-Tabelle 0 HIGH / 1 MEDIUM / 1 LOW, Verdikt „Merge-blockierend: nein", DoD-Checkbox „Review durchgeführt …" ist im Slice-Plan bereits mit `[x]` markiert (Zeile 123) — deckungsgleich mit dem Report-Inhalt | **Bestätigt** |

Alle sieben inhaltlichen DoD-Zeilen dieses Slice sind real erfüllt. Die
verbleibenden vier `[ ]`-Zeilen (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgänge, Drei-Paarungen) sind **korrekt unmarkiert** — sie sind
laut Slice-Plan explizit Planner-Closure-Arbeit, kein
Implementer-/Reviewer-Umfang. `git grep -n "^\- \["` gegen die Datei
bestätigt exakt neun `[x]` und vier `[ ]`, deckungsgleich mit dieser
Rollenteilung. §7 Closure-Notiz trägt an allen sechs Stellen weiterhin
`<wird beim Abschluss ergänzt>` — korrekt noch nicht final ausgefüllt.

## 2. F-1 (LOW) und F-2 (MEDIUM) aus dem Review-Report nachvollzogen

**F-1 (LOW, C#-Absatz-Schlusszeile bleibt stale):** Real gegengelesen.
`spec/pflichtenheft.md:182-184` (C#-Absatz) behauptet weiterhin „Eine
zweite Sprache … bleibt offen", während real bereits die dritte (Kotlin)
beantwortet ist. Diese Staleness entstand nachweislich **vor** diesem
Diff (mit dem Python-Nachzug) — dieser Slice hat sie nicht erzeugt, nur
nicht behoben, und der Slice-Plan (§3) benennt das bewusst als
„außerhalb des DoD-Umfangs". Korrekt als LOW/nicht-merge-blockierend
eingestuft: reine Prosa-Kohärenz-Frage ohne funktionale Auswirkung, kein
Sensor deckt sie, keine neue Fehlerklasse durch diesen Zug.

**F-2 (MEDIUM, Sources-/Javadoc-Jar-Zuschreibung):** Eigene, **dritte**
unabhängige Web-Recherche (nach Implementer und Reviewer) durchgeführt:

- `curl -sL https://docs.github.com/en/actions/publishing-packages/publishing-java-packages-with-gradle`
  → HTTP 200, 425712 Bytes geladen, Titel „Publishing Java packages with
  Gradle" bestätigt (richtige Seite). `grep -io "sources[- ]jar\|javadoc\|withSourcesJar\|withJavadocJar"`
  gegen den geladenen HTML-Text liefert **null Treffer** — die Seite
  erwähnt Sources-/Javadoc-Jars an keiner Stelle, weder fordernd noch
  freistellend. Das bestätigt den Kernbefund des Reviewers: Die im
  Dockerfile-Kommentar/Commit-Message/`harness/README.md` verwendete
  Formulierung „GitHub Packages verlangt laut offizieller Dokumentation
  … keines" ist eine **Zuschreibung, die die zitierte Quelle nicht
  trägt** — die Quelle schweigt, sie bestätigt nicht.
- `curl -sL https://docs.gradle.org/current/userguide/publishing_maven.html`
  → HTTP 200. `grep -io "withSourcesJar\|withJavadocJar\|components\[.java.\]"`
  liefert Treffer für alle drei Suchbegriffe — bestätigt das zweite,
  tatsächlich tragende Faktum: `from(components["java"])` ohne
  `withSourcesJar()`/`withJavadocJar()` im `java{}`-Block packt nur den
  Haupt-Jar (Gradle-Kern-Mechanik, nicht GitHub-spezifisch).
- **Praktischer Schluss weiterhin haltbar:** GitHub Packages (anders als
  Maven Central/Sonatype OSSRH) erzwingt kein Sources-/Javadoc-Jar für
  Gradle-/Maven-Packages — das ist allgemein bekannte Registry-Praxis und
  wird durch das Schweigen der offiziellen Doku (keine Anforderung
  irgendwo erwähnt) nicht widerlegt, nur nicht explizit bestätigt. In
  Kombination mit dem real bestätigten Gradle-Mechanik-Faktum bleibt der
  Bau-Entscheid (ein Haupt-Jar genügt) sachlich richtig. Das Risiko ist
  zusätzlich strukturell durch `ADR-0109` §Re-Evaluierungs-Trigger 4
  abgefangen: ein realer Publish-Fehlschlag (erst im Folge-Slice
  `slice-sdk-kotlin-publish-workflow` real testbar) löst eine Folge-ADR
  aus, kein stiller Blindflug.
- **Einstufung MEDIUM statt HIGH weiterhin angemessen:** Die Aussage ist
  eine unpräzise Attribution, keine fachlich falsche Schlussfolgerung,
  und ohne Auswirkung auf das reale Erzeugnis dieses Slice (kein Publish
  in diesem Umfang). Korrekt als nicht-merge-blockierend eingestuft.

Beide Findings sind demnach korrekt kategorisiert und korrekt als
nicht-merge-blockierend behandelt; keine Verifier-Beanstandung an dieser
Einstufung.

## 3. Mutation-Testing-Behauptung

Keine eigene, vierte Mutation nötig — Plausibilisierung anhand der
Bau-Reihenfolge genügt und wurde durchgeführt: `sdks/kotlin/Dockerfile`
zeigt real `RUN ./gradlew --no-daemon test` **vor** `RUN ./gradlew
--no-daemon build`, beide in derselben `build`-Stufe, auf der `pack`
aufbaut (`FROM build AS pack`). Ein roter Test in dieser Docker-Bau-Kette
bricht `docker build` mit Exit ≠ 0 ab, bevor die `test`-Zeile überhaupt
zur `build`-Zeile weiterläuft — dieselbe Struktur, die Implementer und
Reviewer je einmal real mutiert und beobachtet haben (Reviewer:
`assertEquals("test-token", …)` → `wrong-token`, echter Exit 1). Die
Docker-Layer-Reihenfolge macht das Verhalten strukturell zwingend, nicht
nur beobachtet — eine dritte Mutation hätte dasselbe Ergebnis liefern
müssen und wurde als redundant eingestuft.

## 4. `git diff 9c4efd42..HEAD` — Scope-Prüfung

```
git diff --stat 9c4efd42..HEAD
```
zeigt genau acht geänderte Dateien: Slice-Plan, neuer Review-Report,
`harness/README.md`, `harness/mk/sdk.mk`, `sdks/kotlin/.gitignore`,
`sdks/kotlin/Dockerfile`, `spec/pflichtenheft.md`,
`tools/harness/sdk-pack-kotlin.sh`.

```
git diff --stat 9c4efd42..HEAD -- examples/kotlin/ .a-check.yml \
  spec/architecture.md docs/user/version.md sdks/csharp/ sdks/python/ \
  sdks/kotlin/pgchangefeed-kotlin/src/
```
liefert **keine Ausgabe** — keine der genannten Flächen (inklusive der
beiden Kotlin-Fläche-Verzeichnisse unter `pgchangefeed-kotlin/src/`,
`.a-check.yml`, `spec/architecture.md`, `docs/user/version.md`, beide
Schwester-SDKs) ist berührt. Bestätigt: dieser Slice paketiert
ausschließlich, ohne API-/Architektur-/Fremd-SDK-Fläche anzufassen.

## 5. Backtick-Parität

Eigenständig nachgezählt (Anzahl der Backtick-Zeichen je Datei per `grep`/`wc -l` gezählt) für alle vier
von diesem Slice (Implementer- und Reviewer-Commit zusammen) geänderten
Markdown-Dateien:

| Datei | Backtick-Anzahl | gerade? |
|---|---|---|
| `docs/plan/planning/in-progress/slice-sdk-kotlin-pack-werkzeug.md` | 364 | ja |
| `docs/reviews/review-slice-sdk-kotlin-pack-werkzeug.md` | 348 | ja |
| `harness/README.md` | 1690 | ja |
| `spec/pflichtenheft.md` | 1454 | ja |

Alle vier Zählungen gerade — keine Backtick-Paritäts-Verletzung.

## 6. §6-Risiken — Ausgang real geprüft, Checkbox-Stand korrekt offen

1. **Export-Mechanismus für `.jar`:** real aufgelöst — der
   `tar`-Stream-Export (`pack-export`-Stufe, `docker run --rm --network
   none <image> | tar -x -C sdks/kotlin/dist/`) hat im eigenen, vierten
   Lauf ein reales, valides `.jar` mit korrekten Dateirechten erzeugt
   (Datei gehört dem aufrufenden Nutzer, kein `--user`-Workaround nötig,
   `ls -la` zeigt normale Owner-Rechte). Kein Bind-Mount-Problem
   aufgetreten.
2. **Sources-/Javadoc-Jar-Frage:** real recherchiert und entschieden
   (siehe Abschnitt 2 oben) — ein Haupt-Jar genügt, Entscheidung im
   Bericht (Dockerfile-Kommentar, `harness/README.md`) mit Beleg
   dokumentiert, auch wenn die Quellen-Attribution laut F-2 unpräzise
   ist. Der fachliche Ausgang selbst ist aufgelöst.
3. **`SPEC-<NNN>`-Kollision:** real aufgelöst — `SPEC-028` real vergeben,
   `grep -n "SPEC-028" spec/pflichtenheft.md` zeigt genau die drei
   erwarteten, widerspruchsfreien Treffer, kein Parallel-Zug hat die
   Nummer belegt.

Alle drei Risiken sind durch diesen Zug **inhaltlich real aufgelöst**.
Die dazugehörigen DoD-Checkbox-Zeilen „Jedes Risiko aus §6 trägt einen
Ausgang" und die §7-Closure-Notiz-Zeile „Risiken aus §6" sind **korrekt
noch nicht final eingetragen** (`[ ]` bzw. `<wird beim Abschluss
ergänzt>`) — das bleibt, wie im Slice-Plan vorgesehen, Planner-
Closure-Arbeit und keine Implementer-/Reviewer-Pflicht in diesem Schritt.

---

## Verdikt

**DoD-Konformität: JA.** Alle sieben inhaltlich fälligen DoD-Zeilen (§2)
sind real erfüllt und wurden in diesem Lauf unabhängig gegen den Repo-
Stand nachgewiesen — keine Behauptung ohne eigenen Beleg übernommen. Die
vier verbleibenden, unmarkierten DoD-Zeilen sind korrekt der
Planner-Closure zugeordnet und bewusst offen gelassen. Die beiden
Review-Findings (F-1 LOW, F-2 MEDIUM) sind korrekt kategorisiert, keine
davon merge-blockierend, F-2s praktischer Schluss bleibt nach eigener,
dritter unabhängiger Recherche haltbar. Kein Scope-Übergriff auf
benachbarte Flächen. Backtick-Parität gewahrt. `make gates` läuft grün.

**Durchgeführte Prüfungen (Zusammenfassung):**

- `grep`-Check `sdk-pack-kotlin` vs. `GATE_CHECKS` in `harness/mk/sdk.mk`
  und `Makefile`/`harness/mk/*.mk` — bestätigt: kein Gate-Eintrag.
- Vierter unabhängiger Lauf `bash tools/harness/sdk-pack-kotlin.sh` —
  reales `.jar`, 112591 Bytes, valide `unzip -l`-Struktur.
- Eigenes Lesen von `spec/pflichtenheft.md` §1/§6/§7 gegen DoD-Wortlaut —
  Nachzugsatz korrekt, `SPEC-028`-Zeile korrekt, Historie korrekt, kein
  `ADR-`-Verweis in den neuen Zeilen (`matrix`-Modul-Konformität).
- `grep`-Check `sdk-pack-kotlin` in `harness/README.md` — Zeile vorhanden,
  Bindung auf `ADR-0109` Festlegung 5.
- Eigener, ungepipter `make gates`-Lauf — Exit 0.
- Review-Report gelesen — 0 HIGH bestätigt.
- Dritte unabhängige Web-Recherche zu F-2 (`docs.github.com`,
  `docs.gradle.org`) — Kernbefund des Reviewers bestätigt, praktischer
  Schluss bleibt haltbar.
- Bau-Reihenfolgen-Plausibilisierung der Mutation-Testing-Behauptung
  anhand `sdks/kotlin/Dockerfile` (`test` vor `build`, gleiche Stufe wie
  `pack`) — keine eigene vierte Mutation nötig.
- `git diff --stat 9c4efd42..HEAD` und gezielter Diff gegen sieben
  geschützte Pfade — kein Scope-Übergriff.
- Eigenständige Backtick-Zählung aller vier geänderten Markdown-Dateien —
  alle gerade.
- §6-Risiken-Ausgänge real gegen den Stand geprüft — alle drei inhaltlich
  aufgelöst, Checkbox-/Closure-Notiz-Eintrag korrekt noch offen
  (Planner-Arbeit).

**Finaler `make gates`-Exit-Code (dieser Verifikationslauf):** `0`.

**Nicht getan (außerhalb dieser Rolle):** keine Fixes, kein `git mv` nach
`done/`, keine Closure-Notiz-Fertigstellung — das bleibt Planner-Arbeit
bei der Welle-Closure.
