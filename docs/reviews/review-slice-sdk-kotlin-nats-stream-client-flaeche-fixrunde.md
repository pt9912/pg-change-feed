# Review-Report: slice-sdk-kotlin-nats-stream-client-flaeche (Fixrunde) — 2026-09-22

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md`),
`ADR-0109` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Frisches, unabhängiges Review der **Fixrunde** —
kein Nachvollzug des Vorgänger-Reports, sondern eine eigene Prüfung mit
eigenem Diff-Lesen, eigenem Suchlauf und eigenem `make gates`-Lauf.

**Gegenstand:** Commit `6b279311` ("fix(sdk-kotlin): Fixrunde review
F-1/F-2 Chronik-Kommentar + Dockerfile-Nachzug (`LH-FA-SST-009`,
`ADR-0109`)"), isoliert per `git show 6b279311` geprüft (3 Dateien, 12
Insertions / 10 Deletions) — nicht der breitere Range
`c8c9e3ae..6b279311`, der zusätzlich den zwischenzeitlichen
Review-Report-Commit `49e35ab9` enthält. Vorgänger-Commit: `c8c9e3ae`
(„feat(sdk-kotlin): NATS-Vollinhalts-Client-Fläche + Version-Hebung
0.2.0"), dessen Review (`docs/reviews/review-slice-sdk-kotlin-nats-stream-client-flaeche.md`)
1 HIGH (F-1) + 1 MEDIUM (F-2) fand.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. „Slice-/Wellen-Chronik in Produktionscode-Kommentar",
`AGENTS.md` §3.13 Träger-Nachzug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-22.

**Eingangs-Kontext:**

- `harness/README.md`, `AGENTS.md`, `harness/conventions.md`,
  `.harness/skills/reviewer.md`
- `docs/reviews/review-slice-sdk-kotlin-nats-stream-client-flaeche.md`
  (Vorgänger-Report, F-1/F-2 im Detail)
- `docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md`
  (vollständig, §1–§8)
- `ADR-0109` (Accepted) §Entscheidung Festlegung 1 (letzter Absatz)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.9
  (Exit-Code-Disziplin), §3.13 (Träger-Nachzug)

---

## Findings

### F-1 (Vorgänger-HIGH) — Status: aufgelöst

- `kategorie`: — (kein offenes Finding mehr, Prüfergebnis dokumentiert)
- `quelle`: `AGENTS.md` §3.7 / Reviewer-Skill HIGH „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar"
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:97-100`
- `befund`: Der beanstandete Kommentarabsatz wurde von „Diese
  Version-Hebung (`0.1.0` -> `0.2.0`) trägt außerdem den SSE-Client-Fläche
  (`slice-sdk-kotlin-sse-client-flaeche`, `SPEC-021`) — …" ersetzt durch
  „Version `0.2.0` deckt beide zuletzt gelieferten Client-Flächen ab, SSE
  (`SPEC-021`) und NATS-Vollinhalt (`SPEC-024`, `ADR-0109` Festlegung 1) —
  additive, rückwärtskompatible Erweiterung, SemVer-Minor, keine
  ADR-pflichtige Ausnahme." Real gegen die Datei gelesen (nicht nur den
  Diff-Ausschnitt übernommen): kein Arrow-Muster (`X -> Y`) mehr, keine
  Slice-Kennung mehr, keine Verweis auf einen Slice-Plan-Abschnitt mehr —
  ausschließlich `SPEC-021`/`SPEC-024`/`ADR-0109` als Anker. Reine
  Zustandsaussage, keine Vorher/Nachher-Erzählung.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Chronik in diesem
  Repo (Reviewer ist die tragende Instanz).
- `klasse`: „Slice-/Wellen-Chronik in Produktionscode-Kommentar" — **behoben**.

### F-2 (Vorgänger-MEDIUM) — Status: aufgelöst

- `kategorie`: — (kein offenes Finding mehr, Prüfergebnis dokumentiert)
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug)
- `pfad`: `sdks/kotlin/Dockerfile:77`, `sdks/kotlin/Dockerfile:104`
- `befund`: Beide Zeilen real per `grep -n "0\.1\.0\|0\.2\.0"
  sdks/kotlin/Dockerfile` geprüft — beide tragen jetzt wörtlich
  „`build/libs/pgchangefeed-kotlin-0.2.0.jar`". Kein `0.1.0` mehr in dieser
  Datei. Der Fix trägt exakt die zwei im Vorgänger-Report benannten
  Zeilennummern, keine dritte Stelle in derselben Datei übersehen.
- `verifizierbar`: ja — `grep -n "0\.1\.0" sdks/kotlin/Dockerfile` liefert
  keinen Treffer (in diesem Review-Lauf ausgeführt).
- `klasse`: „Arbeit überholt stehenden Träger" — **behoben**.

## Eigener, dritter unabhängiger Suchlauf (nicht nur Implementer-Bericht übernommen)

Fünf Formulierungsvarianten repo-weit gefahren, unabhängig vom
Implementer-Suchlauf (Ergebnis in der Fixrunden-Commit-Message) und vom
Vorgänger-Reviewer-Suchlauf:

1. `grep -rn "pgchangefeed-kotlin-0\.1\.0\|0\.1\.0\.jar" --include="*.md" --include="*.mk" --include="*.sh" --include="*.kts" --include="Dockerfile" --include="*.yml" --include="Makefile" .`
2. `grep -rn "0\.1\.0" … | grep -i kotlin` (breiterer, nicht auf den Jar-Namen
   beschränkter Treffer)
3. `grep -rn "HTTP-API und gRPC-Stream" --include="*.md" .` (die stale
   Abdeckungs-Formulierung, die dieses Slice ablösen sollte)
4. `grep -rn "BEIDE Testflächen\|zwei Testflächen\|drei Testflächen" .`
5. `grep -rln "pgchangefeed-kotlin-0\.1\.0" .` (Dateiliste, um jede Fundstelle
   einzeln zu klassifizieren)

**Ergebnis:** Alle verbleibenden Treffer außerhalb von `sdks/kotlin/Dockerfile`
gehören zu einer von drei erwarteten, unproblematischen Klassen — keine
davon behauptet fälschlich einen aktuellen Ist-Stand:

- **Records/Historie** (`docs/plan/planning/done/**`,
  `docs/reviews/review-*`, `docs/reviews/verifikation-*`,
  `spec/pflichtenheft.md`s Änderungshistorie-Zeilen, `docs/user/benutzerhandbuch.md`s
  Versionshistorie-Zeilen, `docs/plan/planning/observations/BEO-PGC/**`)
  — diese Dateien sind laut Aufgabenstellung „erwartungsgemäß unverändert
  korrekt": Sie beschreiben einen historischen Zustand zu einem benannten
  Zeitpunkt/Lauf, nicht den heutigen Ist-Stand. Beispiel:
  `docs/plan/planning/welle-sdk-kotlin-vollabdeckung.md:50` nennt
  `pgchangefeed-kotlin-0.1.0.jar` explizit als „Stand dieser Eröffnung"
  (2026-09-21, §2 Trigger) — ein datierter Trigger-Zustand, kein
  Live-Zustandsfeld.
- **Illustrative Beispiele ohne Ist-Stand-Anspruch**
  (`docs/user/releasing.md:225` „z. B. `sdk-kotlin-v0.1.0`",
  `docs/plan/adr/0109-…md:386` „Start bei `0.x.y` (z. B. `0.1.0`)",
  `tools/harness/run-sdk-kotlin-release-tag-info-tests.sh` — Test-Fixture-Werte
  für die Tag-Validierungslogik, kein Bezug zur aktuellen Paket-Version).
- **Die im Slice-Plan selbst stehende Chronik** (§6 Risiken „Version-Hebung
  `0.1.0` → `0.2.0`", §2 DoD-Punkt „von `0.1.0` auf `0.2.0`") — ein
  Planungsdokument trägt legitim eine Vorher/Nachher-Beschreibung des
  eigenen Vorgangs; die HIGH-Klasse „Slice-/Wellen-Chronik" adressiert
  ausdrücklich **Produktionscode**-Kommentare, nicht Planungsartefakte.

Kein weiterer, bislang unbenannter Fund. Die Behauptung der Fixrunden-Commit-Message
„Erweiterter Repo-weiter grep-Suchlauf fand keine weitere übersehene
Stelle" ist damit durch einen dritten, unabhängigen Lauf bestätigt, nicht
nur übernommen.

## Vollständige Neulektüre des Fixrunden-Diffs

`git show 6b279311` komplett gelesen (3 Dateien, 12/10 Zeilen):

- `docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md`:
  ausschließlich die DoD-Checkbox „Review durchgeführt" von `[ ]` auf `[x]`
  mit Report-Verweis und Fixrunden-Zusammenfassung — inhaltlich zutreffend
  (siehe F-1/F-2 oben), keine weitere Zeile geändert.
- `sdks/kotlin/Dockerfile`: ausschließlich die zwei `0.1.0.jar` →
  `0.2.0.jar`-Kommentarwörter, keine `RUN`/`COPY`/`FROM`-Zeile berührt —
  kein funktionaler Bau-Schritt geändert.
- `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`: ausschließlich der
  Kommentarblock über `plugins {}`, keine `version`-Zeile, keine
  Dependency-Koordinate, kein `publishing`-Block berührt (per
  `git diff c8c9e3ae..6b279311 -- sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`
  bestätigt: der Diff enthält ausschließlich `//`-Kommentarzeilen).

Kein neuer Fehler in diesem Diff gefunden. Keine funktionale Code-Änderung
an der NATS-Client-Implementierung selbst (`nats/**`,
`NatsStreamTransport`, `PgChangeFeedNatsStreamClient`) — der Umfang der
Fixrunde bleibt exakt auf die zwei benannten Findings begrenzt.

## Backtick-Parität (selbst nachgezählt)

| Datei | Backticks | Gerade? |
|---|---|---|
| `docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md` | 332 | ja |
| `sdks/kotlin/Dockerfile` | 112 | ja |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | 168 | ja |

Alle drei geänderten Dateien: gerade Backtick-Anzahl, paarig.

## `make gates`

Ungefiltert laufen lassen, Exit-Code direkt (kein Pipe/Wrapper dazwischen)
geprüft: **`0`**. U. a. `generated-sync: OK`, `a-check: gesamt: 0
Befund(e)`.

## Negativbefunde

- geprüft, ohne Befund: F-1 real aufgelöst — kein Arrow-Muster, keine
  Slice-Kennung im neuen Kommentar, ausschließlich `ADR-*`/`SPEC-*`-Anker.
- geprüft, ohne Befund: F-2 real aufgelöst — beide genannten
  Dockerfile-Zeilen (77, 104) tragen `0.2.0.jar`, kein verbleibender
  `0.1.0`-Treffer in dieser Datei.
- geprüft, ohne Befund: dritter, unabhängiger repo-weiter Suchlauf (fünf
  Formulierungsvarianten) — kein weiterer, bislang unbenannter Fund; alle
  verbleibenden `0.1.0`-Treffer sind Records/Historie, illustrative
  Beispiele oder legitime Planungsdokument-Chronik.
- geprüft, ohne Befund: Umfang der Fixrunde — ausschließlich die zwei
  benannten Findings behoben, keine funktionale Code-Änderung an der
  NATS-Client-Implementierung, keine weitere, nicht mandatierte Änderung.
- geprüft, ohne Befund: Backtick-Parität aller drei geänderten Dateien.
- geprüft, ohne Befund: Commit-Traceability — Betreff nennt
  `LH-FA-SST-009`/`ADR-0109`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.
- geprüft, ohne Befund: DoD-Checkbox-Text „kein offenes HIGH/MEDIUM mehr"
  — Aussage stimmt mit dem realen Prüfergebnis dieses Reviews überein.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine neuen — beide Vorgänger-Findings
(„Slice-/Wellen-Chronik in Produktionscode-Kommentar", „Arbeit überholt
stehenden Träger") bestätigt behoben.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW, 0 INFO. Beide
Vorgänger-Findings (F-1 HIGH, F-2 MEDIUM) sind real, an genau den
benannten Stellen behoben; ein eigener, dritter unabhängiger Suchlauf mit
fünf Formulierungsvarianten fand keinen weiteren, bislang unbenannten
Fund; die Fixrunde hat den Umfang des ursprünglichen Diffs nicht
überschritten (keine funktionale Code-Änderung, ausschließlich
Kommentartext in zwei Dateien plus die DoD-Checkbox).

**Übergabe:** Keine weitere Fixrunde nötig. Kein Rollen-Widerspruch, keine
Konflikt-Pfad-Sequenz über den Architect (Modul 8) — dieser Lauf endet
ohne Rückgabe-Pfeil an den Implementer.

**DoD-Checkbox-Nachzug:** entfällt — die DoD-Zeile „Review durchgeführt"
wurde bereits im Fixrunden-Commit (`6b279311`) auf `[x]` gesetzt, mit
korrektem Verweis auf den Vorgänger-Report und einer zutreffenden
Fixrunden-Zusammenfassung. Dieser Report ist der zweite, bestätigende
Lauf-Beleg zu dieser Checkbox, kein eigener Nachzug-Anlass
(Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde — hier lag ohnehin
bereits eine Fixrunde vor, deren regulärer Nachzug bereits erfolgt ist).

Dieser Report ersetzt keine Verifikation gegen die volle DoD — das bleibt
Verifier-Aufgabe (Modul 11).
